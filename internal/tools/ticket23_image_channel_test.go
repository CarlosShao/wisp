package tools

// Ticket 23 AC#2, the TOOLS half of the C7 image seam.
//
// a67e1d335 landed the loop+llm half (ImagePart in the history, wire shapes per
// protocol, memory-only `mem:<n>` handles). It left three holes on this side of
// the package boundary, all three named in
// docs/evidence/s1/23-c7-imagepart-llm-seam.md §7:
//
//  1. Result had no field that can carry pixels, so screen.capture's PNG died at
//     `renderSystemResult` and only the metadata reached the model;
//  2. the bridge's agent.ToolOutcome literal dropped anything image-shaped, so
//     even a tool that DID return pixels would have them dropped one layer up;
//  3. system.go's text told the model "C7 ImagePart 到 LLM 的接缝本仓当前不存在，
//     所以模型看不到画面" - true before a67e1d335, a lie after it.
//
// The legs below are written against the seam, not against an implementation
// detail: they assert what the model receives (an image with these pixels) and
// what must NOT receive it (the text, the tool_call row, the audit log).

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/risk"
)

// pixels23 is a payload shaped like a real capture (PNG signature + a marker
// that is unmistakable in any encoding) and long enough that an accidental
// %q/%x/%v of it is unmistakable too.
var pixels23 = append([]byte("\x89PNG\r\n\x1a\n"),
	[]byte("ticket23-pixels-marker-0123456789-abcdefghij")...)

// leakForms23 renders the payload in the spellings a careless implementation
// would use: raw, base64 (the data-URL form), hex, and Go's decimal-byte list.
func leakForms23() []string {
	return []string{
		string(pixels23),
		base64.StdEncoding.EncodeToString(pixels23),
		hex.EncodeToString(pixels23),
		fmt.Sprint(pixels23),
	}
}

// assertNoPixels23 fails when any of the given blobs contains the pixels in any
// spelling. `where` names the sink so a failure says which one leaked.
func assertNoPixels23(t *testing.T, where, blob string) {
	t.Helper()
	for _, form := range leakForms23() {
		if strings.Contains(blob, form) {
			t.Errorf("%s 里出现了截图像素（%v）：截图只许走内存图像通道", where, form)
			return
		}
	}
}

// image23Journal records what the bridge booked, so the "pixels are not in the
// database" half is asserted on the REAL row the bridge writes rather than on
// the absence of a call.
type image23Journal struct {
	mu   sync.Mutex
	rows []image23Row
}

// image23Row is the handful of columns these legs care about, copied out of
// memory.ToolCall at booking time: a column this test never heard of cannot
// smuggle pixels into the assertion below without someone reading the copy and
// noticing.
type image23Row struct {
	Tool     string
	ArgsJSON string
	Outcome  string
}

func (j *image23Journal) StartTaskLog(context.Context, memory.TaskLog) error { return nil }

func (j *image23Journal) FinishTaskLog(context.Context, string, string, string, int64, int64, int64, string) error {
	return nil
}

func (j *image23Journal) InsertToolCall(_ context.Context, tc memory.ToolCall) (int64, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.rows = append(j.rows, image23Row{
		Tool: tc.Tool, ArgsJSON: tc.ArgsJSON, Outcome: tc.Outcome,
	})
	return int64(len(j.rows)), nil
}

func (j *image23Journal) DecideToolCall(context.Context, int64, string, *int64) error { return nil }

func (j *image23Journal) FinishToolCall(context.Context, int64, string, string) error { return nil }

func (j *image23Journal) booked() []image23Row {
	j.mu.Lock()
	defer j.mu.Unlock()
	return append([]image23Row(nil), j.rows...)
}

// bridge23Image composes the REAL bridge over the REAL eleven-row family and one
// fake host backend: the gate is the only stub, and it says "yes" so a leg that
// wants the pixels is not accidentally measuring a refusal.
func bridge23Image(t *testing.T, f *fake23Backend, j agent.Journal,
	logged *[]string,
) (*Bridge, map[string]Entry) {
	t.Helper()
	reg := NewRegistry()
	byName := map[string]Entry{}
	for _, e := range BuiltinSystemEntries(SystemDeps{Bridge: &SystemBridge{Backend: f}}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
		byName[e.Tool.Name()] = e
	}
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	var logf func(string, ...any)
	if logged != nil {
		logf = func(format string, args ...any) {
			*logged = append(*logged, fmt.Sprintf(format, args...))
		}
	} else {
		logf = func(string, ...any) {}
	}
	b := New(Options{
		Registry:   reg,
		Paths:      NewPathCanonicalizer(nil, nil),
		Provenance: prov,
		Journal:    j,
		Gate: GateFuncs{
			Window:   func(context.Context, Decision) (Answer, string) { return AnswerAllow, "" },
			Approval: func(context.Context, Decision) (Answer, string) { return AnswerAllow, "" },
		},
		Logf: logf,
	})
	return b, byName
}

// Test23ScreenCaptureCarriesPixelsOnTheImageChannel is the core AC#2 assertion
// on this side: the bytes the backend produced arrive as an Image (bytes +
// mime + alt) on the tool's Result, and the TEXT keeps only the metadata.
func Test23ScreenCaptureCarriesPixelsOnTheImageChannel(t *testing.T) {
	f := &fake23Backend{reply: func(a Action, _ map[string]any) (ActionResult, error) {
		if a != ActionScreenCapture {
			return ActionResult{Detail: "unexpected " + string(a)}, nil
		}
		return ActionResult{
			Bytes: pixels23, Mime: "image/png", DRM: false, Target: "screen",
		}, nil
	}}
	_, byName := systemRoster23(t, f)
	res := run23(t, byName["screen.capture"], `{"target":"full"}`)

	if len(res.Images) != 1 {
		t.Fatalf("screen.capture 的 Result 带 %d 张图，截到的那一张要在这条通道上（模型看不到就没有 AC#2）", len(res.Images))
	}
	img := res.Images[0]
	if !bytes.Equal(img.Bytes, pixels23) {
		t.Errorf("图像通道里的字节不是后端截到的那份（%d vs %d 字节）", len(img.Bytes), len(pixels23))
	}
	if img.Mime != "image/png" {
		t.Errorf("图像通道的 Mime=%q，要后端给的 image/png", img.Mime)
	}
	if !strings.Contains(img.Alt, "screen") {
		t.Errorf("图像通道的 Alt=%q，要说出截的是哪块屏（C7 的 alt）", img.Alt)
	}
	// The text still carries the metadata, and still must not carry pixels.
	for _, want := range []string{`"mime":"image/png"`, `"drm":false`} {
		if !strings.Contains(res.Text, want) {
			t.Errorf("元数据文本读不到 %q：%q", want, res.Text)
		}
	}
	assertNoPixels23(t, "Result.Text", res.Text)
	if strings.Contains(res.Text, "mem:") || strings.Contains(res.Text, "data:image") {
		t.Errorf("Result.Text 里出现了句柄/内联数据（工具层不该知道句柄形状）：%q", res.Text)
	}
}

// Test23ScreenCaptureImageCrossesTheBridge pins the SECOND hole: the bridge's
// ToolOutcome must forward the image, or the loop never sees it (the loop half
// landed at a67e1d335 and reads ToolOutcome.Images).
func Test23ScreenCaptureImageCrossesTheBridge(t *testing.T) {
	f := &fake23Backend{reply: func(Action, map[string]any) (ActionResult, error) {
		return ActionResult{Bytes: pixels23, Mime: "image/png", Target: "记事本"}, nil
	}}
	b, _ := bridge23Image(t, f, nil, nil)
	out, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "t23-img", CorrelationID: "t23-img-c1", CallID: "call-1",
		Name: "screen.capture", Args: json.RawMessage(`{"target":"window","title":"记事本"}`),
	})
	if err != nil {
		t.Fatalf("bridge.Execute 返回宿主错误：%v", err)
	}
	if out.IsError {
		t.Fatalf("bridge 把截图判成错误：%q", out.Text)
	}
	if out.RiskLevel != "L1" {
		t.Errorf("screen.capture 的 RiskLevel=%q，D34 是 L1", out.RiskLevel)
	}
	if len(out.Images) != 1 {
		t.Fatalf("ToolOutcome 带 %d 张图，bridge 把工具的图丢了（loop 永远收不到）", len(out.Images))
	}
	got := out.Images[0]
	if !bytes.Equal(got.Data, pixels23) {
		t.Errorf("ToolOutcome 的图像字节不是工具给的那份（%d vs %d）", len(got.Data), len(pixels23))
	}
	if got.MimeType != "image/png" {
		t.Errorf("ToolOutcome 的 MimeType=%q，要 image/png", got.MimeType)
	}
	if !strings.Contains(got.Alt, "记事本") {
		t.Errorf("ToolOutcome 的 Alt=%q，要带上窗口标题（C7 的 alt 就是给人看的那句）", got.Alt)
	}
	assertNoPixels23(t, "ToolOutcome.Text", out.Text)
}

// Test23CapturePixelsNeverReachRowOrLog pins the privacy half of D34's
// screen.capture note ("只存内存 bytes_ref，不落盘、不进日志、不进 DB") at the
// layer that does the booking: the tool_call row and the audit sink must be
// pixel-free even though the same call DID carry the pixels in memory.
func Test23CapturePixelsNeverReachRowOrLog(t *testing.T) {
	var logged []string
	j := &image23Journal{}
	f := &fake23Backend{reply: func(Action, map[string]any) (ActionResult, error) {
		return ActionResult{Bytes: pixels23, Mime: "image/png", Target: "screen"}, nil
	}}
	b, _ := bridge23Image(t, f, j, &logged)
	out, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "t23-priv", CorrelationID: "t23-priv-c1", CallID: "call-2",
		Name: "screen.capture", Args: json.RawMessage(`{"target":"full"}`),
	})
	if err != nil || out.IsError {
		t.Fatalf("前置失败：截图调用没跑通 err=%v IsError=%v Text=%q", err, out.IsError, out.Text)
	}
	// Positive control: the pixels really were carried by this call, so the
	// absence below is about the SINKS, not about a screenshot that never happened.
	if len(out.Images) != 1 {
		t.Fatalf("前置失败：这次调用没有图像，下面的「没进库/日志」就没有意义")
	}
	rows := j.booked()
	if len(rows) != 1 {
		t.Fatalf("tool_call 行写了 %d 条，要 1 条（否则排查的是别的东西）", len(rows))
	}
	assertNoPixels23(t, "tool_call.ArgsJSON", rows[0].ArgsJSON)
	for i, line := range logged {
		assertNoPixels23(t, fmt.Sprintf("审计日志第 %d 行", i), line)
	}
}

// Test23NonCaptureResultsCarryNoPhantomImage is the negative control: only the
// bytes branch fills the image channel. A read verb must not grow an image, and
// neither must a capture whose backend returned no bytes (a refusal path).
func Test23NonCaptureResultsCarryNoPhantomImage(t *testing.T) {
	f := &fake23Backend{reply: func(a Action, _ map[string]any) (ActionResult, error) {
		switch a {
		case ActionSystemGet:
			return ActionResult{Report: map[string]any{"os": "windows"}}, nil
		case ActionScreenCapture:
			return ActionResult{Mime: "image/png", DRM: true, Target: "screen"}, nil
		}
		return ActionResult{Detail: "ok"}, nil
	}}
	_, byName := systemRoster23(t, f)

	get := run23(t, byName["system.get"], `{}`)
	if len(get.Images) != 0 {
		t.Errorf("system.get 的 Result 带了 %d 张图，读状态不产生图像（凭空出现图像会污染上下文）", len(get.Images))
	}
	empty := run23(t, byName["screen.capture"], `{"target":"full"}`)
	if len(empty.Images) != 0 {
		t.Errorf("后端没返回像素时 Result 仍带 %d 张图（空图会让模型以为看到了画面）", len(empty.Images))
	}
}

// Test23CaptureAltNamesDRMBlackFrame pins the one annotation AC#2 asks this
// layer for: a DRM-protected surface legitimately comes back black, and BOTH
// the alt (what the model reads beside the image) and the text must say so
// instead of letting the model describe a black rectangle as the user's screen.
func Test23CaptureAltNamesDRMBlackFrame(t *testing.T) {
	f := &fake23Backend{reply: func(Action, map[string]any) (ActionResult, error) {
		return ActionResult{Bytes: pixels23, Mime: "image/png", DRM: true, Target: "screen"}, nil
	}}
	_, byName := systemRoster23(t, f)
	res := run23(t, byName["screen.capture"], `{"target":"full"}`)
	if len(res.Images) != 1 {
		t.Fatalf("drm=true 的截图没有进图像通道（黑帧也是它该看到的画面）")
	}
	if alt := res.Images[0].Alt; !strings.Contains(alt, "DRM") && !strings.Contains(alt, "黑帧") {
		t.Errorf("drm=true 时 Alt=%q 没说明黑帧（模型会把纯黑当成桌面真黑）", alt)
	}
	if !strings.Contains(res.Text, "DRM") && !strings.Contains(res.Text, "黑帧") {
		t.Errorf("drm=true 时文本 %q 没说明黑帧", res.Text)
	}
}
