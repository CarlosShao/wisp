package agent

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
)

// Ticket 23 AC#2 - the C7 ImagePart -> LLM seam.
//
// What is under test is exactly one hop: a tool's answer carries image bytes
// (screen.capture is the production source), and those bytes must reach the
// mock LLM as an llm.ImagePart - resolved through Request.ResolveBytes and
// encoded onto the wire - while the bytes_ref itself never touches the disk,
// the database rows or a log line.
//
// The tool provider below is the host-bridge stand-in: the golden fixture
// (internal/agent/testdata/golden, frozen - not touched by this ticket) names
// its call "echo", so the declared entry is echo while its OUTCOME is what
// screen.capture's ActionResult becomes once the host bridge hands it to the
// loop: a text line plus the in-memory PNG. The real screen.capture tool lives
// in internal/tools, outside this leg's write face, so this test asserts the
// seam half that the loop and the adapters own.
//
// Honesty note (the AC's own words): the assertion is against the WIRE the
// genuine openai-chat adapter built (harness bodies), not against a function
// mock. If ResolveBytes is not wired, the adapter refuses the request and the
// task fails - so a green here cannot be a false green.

const (
	ticket23Mime = "image/png"
	ticket23Alt  = "桌面截图"
)

// ticket23PNG is NOT a real PNG: only the bytes matter for the seam (they are
// unique enough that a disk/db/log leak is greppable).
func ticket23PNG() []byte {
	b := []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a}
	for i := 0; i < 96; i++ {
		b = append(b, byte(0x40+i))
	}
	return b
}

// screenCaptureProvider answers the fixture's one call with screen.capture's
// outcome shape (text metadata + memory-only bytes).
type screenCaptureProvider struct {
	mu      sync.Mutex
	calls   int
	outcome ToolOutcome
	// text is what the tool reports to the model beside the image.
	text string
}

func (p *screenCaptureProvider) Tools(context.Context) ([]ToolInfo, error) {
	return []ToolInfo{{
		Name:        "echo",
		Description: "screen.capture stand-in: 截取前台窗口并让模型看到像素",
		Resident:    true,
		RiskLevel:   RiskL0,
		Parameters:  json.RawMessage(`{"type":"object","properties":{}}`),
	}}, nil
}

func (p *screenCaptureProvider) Execute(_ context.Context, _ ToolRequest) (ToolOutcome, error) {
	p.mu.Lock()
	p.calls++
	p.mu.Unlock()
	out := p.outcome
	out.Text = p.text
	return out, nil
}

func (p *screenCaptureProvider) count() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.calls
}

// TestScreenCaptureImageReachesMockLLM23 is AC#2's positive half: the wire body
// of the SECOND request (the one that follows the tool result) carries the PNG
// as an image part, in the data-URL shape the openai-chat adapter emits, and
// the C7 history carries an ImagePart whose BytesRef is an in-memory handle.
func TestScreenCaptureImageReachesMockLLM23(t *testing.T) {
	png := ticket23PNG()
	p := &screenCaptureProvider{
		text: `{"mime":"image/png","bytes_len":104,"drm":false}`,
		outcome: ToolOutcome{Images: []ImageAttachment{{
			MimeType: ticket23Mime, Data: png, Alt: ticket23Alt,
		}}},
	}
	recs := &ticket23Logs{}
	h := newHarness(t, "tool-then-text", withTools(p), withLogger(recs.logger()))

	res := h.run("截个图看看屏幕上是什么")
	if res.Status != StatusCompleted {
		t.Fatalf("fixture did not run: status=%s msg=%s err=%v", res.Status, res.Message, res.Err)
	}
	if p.count() != 1 {
		t.Fatalf("tool dispatches = %d, want the one this fixture makes", p.count())
	}

	// --- C7 half: the loop's own history carries the ImagePart -------------
	hist := h.loop.History()
	var img llmImagePartView
	found := false
	for _, m := range hist {
		for _, part := range m.Content {
			tr, ok := part.(llmToolResult)
			if !ok {
				continue
			}
			var text []string
			for _, inner := range tr.Content {
				switch c := inner.(type) {
				case llmText:
					text = append(text, c.Text)
				case llm.ImagePart:
					found = true
					img = llmImagePartView{mime: c.MimeType, ref: c.BytesRef, alt: c.Alt}
				}
			}
			if strings.Join(text, "\n") != p.text {
				t.Errorf("tool result text = %q, want the tool's own %q", strings.Join(text, "\n"), p.text)
			}
		}
	}
	if !found {
		t.Fatal("the tool result in history carries no llm.ImagePart: the screenshot never joined the C7 conversation")
	}
	if img.mime != ticket23Mime || img.alt != ticket23Alt {
		t.Errorf("ImagePart = {%s %s}, want {%s %s}", img.mime, img.alt, ticket23Mime, ticket23Alt)
	}
	if !strings.HasPrefix(img.ref, "mem:") {
		t.Fatalf("ImagePart.BytesRef = %q, want an in-memory handle (mem:...), never a path", img.ref)
	}

	// --- wire half: the second request body carried the bytes --------------
	bodies := h.requestBodies()
	if len(bodies) < 2 {
		t.Fatalf("provider requests = %d, want 2 (tool turn + after-tool turn)", len(bodies))
	}
	body := string(bodies[1])
	want := "data:" + ticket23Mime + ";base64," + base64.StdEncoding.EncodeToString(png)
	if !strings.Contains(body, want) {
		t.Fatalf("the screenshot's bytes never reached the wire.\nbody=%s", trunc(body, 1200))
	}
	if strings.Contains(body, img.ref) {
		t.Errorf("the bytes_ref %q leaked onto the wire body (the adapter must resolve it, not forward it)", img.ref)
	}
	if !strings.Contains(body, `"type":"image_url"`) {
		t.Errorf("wire body has no image_url part: %s", trunc(body, 1200))
	}

	// --- privacy half: no ref in the log records --------------------------
	for _, r := range recs.snapshot() {
		line := r.msg
		for k, v := range r.attr {
			line += " " + k + "=" + strings.TrimSpace(strings.ReplaceAll(
				strings.TrimSpace(jsonScalar(v)), "\n", " "))
		}
		if strings.Contains(line, img.ref) {
			t.Errorf("the bytes_ref leaked into a log record: %s", line)
		}
	}

	// --- privacy half: no PNG on disk under the artifacts dir -------------
	assertNoPNGOnDisk23(t, h.loop.opt.Config.ArtifactsDir, png)
}

// TestImageBytesRefNeverPersistedInJournal23 is AC#2's database half: the
// tool_call rows the loop writes carry arguments and status, never the image
// handle (a DB row is exactly where a bytes_ref must not live).
func TestImageBytesRefNeverPersistedInJournal23(t *testing.T) {
	png := ticket23PNG()
	p := &screenCaptureProvider{
		text: `{"mime":"image/png","bytes_len":104,"drm":false}`,
		outcome: ToolOutcome{Images: []ImageAttachment{{
			MimeType: ticket23Mime, Data: png, Alt: ticket23Alt,
		}}},
	}
	j := &fakeJournal{}
	h := newHarness(t, "tool-then-text", withTools(p), withJournal(j))

	res := h.run("截个图看看屏幕上是什么")
	if res.Status != StatusCompleted {
		t.Fatalf("fixture did not run: status=%s msg=%s", res.Status, res.Message)
	}

	var ref string
	for _, m := range h.loop.History() {
		for _, part := range m.Content {
			tr, ok := part.(llmToolResult)
			if !ok {
				continue
			}
			for _, inner := range tr.Content {
				if ip, ok := inner.(llm.ImagePart); ok {
					ref = ip.BytesRef
				}
			}
		}
	}
	if ref == "" {
		t.Fatal("no ImagePart in history: nothing to assert about persistence")
	}
	for _, row := range j.all() {
		blob := strings.Join([]string{
			row.Tool, row.ArgsJSON, row.RiskLevel, row.CorrelationID,
			row.Decision, row.Outcome, row.ErrorClass,
		}, "\x00")
		if strings.Contains(blob, ref) {
			t.Errorf("tool_call row carries the image handle %q: %+v", ref, row)
		}
	}
}

// TestImageBudgetRefusalIsStated23: an image that does not fit the loop's
// in-memory budget is DROPPED WITH A SENTENCE, not silently and not by writing
// it somewhere else. The model must not be able to believe it saw pixels that
// never left the tool.
func TestImageBudgetRefusalIsStated23(t *testing.T) {
	huge := bytes.Repeat([]byte{0x7f}, maxImageBytes+1)
	p := &screenCaptureProvider{
		text: `{"mime":"image/png","bytes_len":1,"drm":false}`,
		outcome: ToolOutcome{Images: []ImageAttachment{{
			MimeType: ticket23Mime, Data: huge, Alt: ticket23Alt,
		}}},
	}
	h := newHarness(t, "tool-then-text", withTools(p))
	res := h.run("截个图看看屏幕上是什么")
	if res.Status != StatusCompleted {
		t.Fatalf("fixture did not run: status=%s msg=%s", res.Status, res.Message)
	}

	var resultText string
	for _, m := range h.loop.History() {
		for _, part := range m.Content {
			tr, ok := part.(llmToolResult)
			if !ok {
				continue
			}
			for _, inner := range tr.Content {
				if _, ok := inner.(llm.ImagePart); ok {
					t.Fatal("an over-budget image was attached as an ImagePart: the budget is not enforced")
				}
				if t, ok := inner.(llmText); ok {
					resultText += t.Text
				}
			}
		}
	}
	if !strings.Contains(resultText, "图像") || !strings.Contains(resultText, "未送入模型") {
		t.Fatalf("the dropped image was not stated to the model: %q", resultText)
	}
	body := string(h.requestBodies()[1])
	if strings.Contains(body, base64.StdEncoding.EncodeToString(huge[:32])) {
		t.Fatal("the over-budget image reached the wire anyway")
	}
}

// TestImageStoreIsMemoryOnly23 pins the store's two promises: refs are opaque
// memory handles (not paths), and Clear really drops the bytes (session end).
func TestImageStoreIsMemoryOnly23(t *testing.T) {
	s := newImageStore(1 << 20)
	ref, err := s.put(ticket23Mime, ticket23PNG(), ticket23Alt)
	if err != nil {
		t.Fatalf("put: %v", err)
	}
	if !strings.HasPrefix(ref, "mem:") {
		t.Fatalf("ref = %q, want the mem: scheme", ref)
	}
	if filepath.IsAbs(ref) || strings.Contains(ref, `\`) || strings.Contains(ref, "/") {
		t.Fatalf("ref %q looks like a filesystem path; the bytes must never have a path", ref)
	}
	mime, data, err := s.resolve(ref)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if mime != ticket23Mime || !bytes.Equal(data, ticket23PNG()) {
		t.Fatalf("resolve returned {%s %d bytes}, want the stored image", mime, len(data))
	}
	s.clear()
	if _, _, err := s.resolve(ref); err == nil {
		t.Fatal("resolve after clear succeeded: the bytes outlived the session")
	}
	if ref2, err := s.put(ticket23Mime, ticket23PNG(), ticket23Alt); err != nil || ref2 == "" {
		t.Fatalf("put after clear: ref=%q err=%v", ref2, err)
	}
}

// TestImageStoreRejectsOverBudget23 is the store's own bound (the loop test
// above is the model-facing half).
func TestImageStoreRejectsOverBudget23(t *testing.T) {
	s := newImageStore(64)
	if _, err := s.put(ticket23Mime, bytes.Repeat([]byte{1}, 65), "alt"); err == nil {
		t.Fatal("put accepted an image above the store budget")
	}
	if _, err := s.put(ticket23Mime, bytes.Repeat([]byte{1}, maxImageBytes+1), "alt"); err == nil {
		t.Fatal("put accepted an image above the per-image ceiling")
	}
}

// ---------------------------------------------------------------------------
// small local helpers (kept in this file: the test is the only consumer)

type llmImagePartView struct {
	mime string
	ref  string
	alt  string
}

// ticket23Logs is this ticket's own slog capture (self-contained: it must not
// borrow another ticket's fixture, because a change there would silently
// weaken the leak assertion here). Every record is rendered for a substring
// scan, exactly what a log/diagnostics bundle would see.
type ticket23Logs struct {
	mu   sync.Mutex
	recs []capturedView
}

type capturedView struct {
	msg  string
	attr map[string]any
}

func (l *ticket23Logs) Enabled(context.Context, slog.Level) bool { return true }

func (l *ticket23Logs) Handle(_ context.Context, r slog.Record) error {
	attrs := map[string]any{}
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})
	l.mu.Lock()
	l.recs = append(l.recs, capturedView{msg: r.Message, attr: attrs})
	l.mu.Unlock()
	return nil
}

func (l *ticket23Logs) WithAttrs([]slog.Attr) slog.Handler { return l }
func (l *ticket23Logs) WithGroup(string) slog.Handler      { return l }
func (l *ticket23Logs) logger() *slog.Logger               { return slog.New(l) }

func (l *ticket23Logs) snapshot() []capturedView {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]capturedView(nil), l.recs...)
}

func jsonScalar(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case error:
		return x.Error()
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

// assertNoPNGOnDisk23 walks the artifacts directory and fails if the PNG bytes
// (or a base64 of them) sit anywhere on disk.
func assertNoPNGOnDisk23(t *testing.T, dir string, png []byte) {
	t.Helper()
	b64 := base64.StdEncoding.EncodeToString(png)
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		b, rerr := os.ReadFile(path)
		if rerr != nil {
			return nil
		}
		if bytes.Contains(b, png) || strings.Contains(string(b), b64) {
			t.Errorf("the screenshot's bytes were written to disk at %s", path)
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("walk artifacts dir: %v", err)
	}
}
