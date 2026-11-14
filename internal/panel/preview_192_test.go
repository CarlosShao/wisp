package panel

// Criteria for ticket 192 AC#1 and AC#2 (artifact preview). Every criterion
// here answers one of the four questions the ticket's AC#1 asks, or pins the
// injection face AC#2 names. The rulers are:
//
//	go test -count=1 -run 'TestPreview192' ./internal/panel/
//
// and the red proof that they can actually fail is
// .scratch/wisp/probes/192/redproof-192.sh (mutants applied to an OUT-OF-REPO
// export tree; the evidence file names the assertion each mutant reddens).
//
// What these criteria deliberately do NOT do: they do not fake the artifacts
// face. Every fixture goes through a REAL memory.Store (memory.Open +
// PutArtifact), so the name guard, the directory and the file shapes under
// test are the production ones (AGENTS 1.3: no mock in place of the real
// thing). The risk engine in the quote-door criteria is likewise the real
// C25 engine, and the egress half is proven by the real CheckText/Inspect.

import (
	"bytes"
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"log/slog"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/risk"
)

// --- fixtures -----------------------------------------------------------------

func preview192Store(t *testing.T) *memory.Store {
	t.Helper()
	store, err := memory.Open(t.TempDir(), memory.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

func preview192Resolver(t *testing.T, store *memory.Store, limits PreviewLimits) *PreviewResolver {
	t.Helper()
	r, err := NewPreviewResolverWithLimits(store, memory.ValidArtifactName, limits)
	if err != nil {
		t.Fatalf("NewPreviewResolverWithLimits: %v", err)
	}
	return r
}

func preview192Put(t *testing.T, store *memory.Store, name string, data []byte) {
	t.Helper()
	if err := store.PutArtifact(context.Background(), name, data); err != nil {
		t.Fatalf("PutArtifact(%q): %v", name, err)
	}
}

func preview192PNG(n int) []byte {
	data := make([]byte, n)
	copy(data, []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a})
	for i := 8; i < n; i++ {
		data[i] = byte(i * 7)
	}
	return data
}

// preview192LargeText is a valid-UTF-8 file big enough to be truncated, with
// multi-byte characters so a byte-wise cut that ignores UTF-8 shows up as
// mojibake rather than passing silently.
func preview192LargeText() []byte {
	var b bytes.Buffer
	for i := 0; b.Len() < 300_000; i++ {
		fmt.Fprintf(&b, "第%06d行：这是一段用于预览截断与续读的内容，带多字节字符。\n", i)
	}
	return b.Bytes()
}

// preview192NoticeRe reads the numbers out of the first-view notice, i.e. the
// exact sentence a human sees. Parsing the display (not an internal field) is
// the point: the numbers on screen are the ones that must be true.
var preview192NoticeRe = regexp.MustCompile(`共 (\d+) 字节 · 显示开头 (\d+) 与结尾 (\d+) 字节，中间 (\d+) 字节未显示`)

// preview192ResumeRe reads the resume pointer off the notice's second line.
var preview192ResumeRe = regexp.MustCompile(`\[wisp 续读\] (\S+) · 每窗最多 (\d+) 字节`)

// --- AC#1 question 1: which classes ------------------------------------------

func TestPreview192FourClassesAndTheRestRefused(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{})
	ctx := context.Background()

	cases := []struct {
		name   string
		data   []byte
		want   PreviewKind
		render PreviewRender
		status int
		reason string
		why    string
	}{
		{
			name: "artifact-html.bin",
			data: []byte("<!DOCTYPE html><html><body><p>hi</p></body></html>"),
			want: PreviewKindHTML, render: PreviewRendersInline, status: 200,
			why: "HTML is one of the four classes; its bytes are sniffed, not the extension",
		},
		{
			name: "artifact-image.bin",
			data: preview192PNG(2048),
			want: PreviewKindImage, render: PreviewRendersInline, status: 200,
			why: "PNG is the composer's own image table, reused",
		},
		{
			name: "artifact-pdf.bin",
			data: []byte("%PDF-1.4\n1 0 obj\n<<>>\nendobj\ntrailer\n<<>>\n%%EOF\n"),
			want: PreviewKindPDF, render: PreviewRendersHostPDF, status: 200,
			why: "PDF is classified and served, but rendering depends on the WebView2 runtime",
		},
		{
			name: "artifact-text.bin",
			data: []byte("星期一 09:41 任务输出\n第二行\n"),
			want: PreviewKindText, render: PreviewRendersInline, status: 200,
			why: "plain UTF-8 text is the fourth class",
		},
		{
			name: "artifact-xml.bin",
			data: []byte(`<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg"><text>x</text></svg>`),
			want: PreviewKindText, render: PreviewRendersInline, status: 200,
			why: "SVG is XML, and it is deliberately NOT the image class: an SVG can carry script, so it is shown as text",
		},
		{
			name:   "artifact-video.bin",
			data:   append([]byte{0, 0, 0, 0x18}, []byte("ftypisom\x00\x00\x02\x00isomiso2")...),
			status: 415,
			reason: "视频",
			why:    "video is stored by the composer but is not one of AC#1's four classes",
		},
		{
			name:   "artifact-binary.bin",
			data:   []byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x1f, 0x7f, 0x00, 0x10, 0x11},
			status: 415,
			reason: "四类",
			why:    "unrecognised bytes are refused by name, never guessed",
		},
		{
			name:   "artifact-utf16.bin",
			data:   append([]byte{0xff, 0xfe}, []byte("h\x00i\x00\n\x00")...),
			status: 415,
			reason: "UTF-8",
			why:    "UTF-16 text is refused rather than rendered as mojibake",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			preview192Put(t, store, tc.name, tc.data)
			doc := r.Lookup(ctx, PreviewURL(tc.name))
			if tc.want == "" {
				if doc.Status != tc.status {
					t.Fatalf("%s: status = %d (%s), want %d - %s", tc.name, doc.Status, doc.Reason, tc.status, tc.why)
				}
				if doc.Kind != "" {
					t.Fatalf("%s: refused artifact still reported kind %q - %s", tc.name, doc.Kind, tc.why)
				}
				if len(doc.Body) != 0 {
					t.Fatalf("%s: a refusal carried %d byte(s) of body", tc.name, len(doc.Body))
				}
				if tc.reason != "" && !strings.Contains(doc.Reason, tc.reason) {
					t.Fatalf("%s: refusal %q does not name %q - %s", tc.name, doc.Reason, tc.reason, tc.why)
				}
				if doc.TotalBytes != int64(len(tc.data)) {
					t.Fatalf("%s: refusal reported TotalBytes=%d, want %d (the size is part of the honest no)",
						tc.name, doc.TotalBytes, len(tc.data))
				}
				return
			}
			if doc.Status != 200 {
				t.Fatalf("%s: status = %d (%s), want 200 - %s", tc.name, doc.Status, doc.Reason, tc.why)
			}
			if doc.Kind != tc.want {
				t.Fatalf("%s: kind = %q, want %q - %s", tc.name, doc.Kind, tc.want, tc.why)
			}
			if doc.Render != tc.render {
				t.Fatalf("%s: render = %q, want %q - %s", tc.name, doc.Render, tc.render, tc.why)
			}
			if !doc.Complete {
				t.Fatalf("%s: a class that is rendered whole reported Complete=false", tc.name)
			}
			if tc.want != PreviewKindText && !bytes.Equal(doc.Body, tc.data) {
				t.Fatalf("%s: body is not the artifact bytes (%d vs %d byte(s))", tc.name, len(doc.Body), len(tc.data))
			}
		})
	}
}

func TestPreview192ClassComesFromBytesNotFromTheName(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{})
	ctx := context.Background()
	png := preview192PNG(256)
	names := []string{"looks-like-text.txt", "looks-like-a-page.html", "no-extension", "looks-like-pdf.pdf"}
	for _, name := range names {
		preview192Put(t, store, name, png)
		doc := r.Lookup(ctx, PreviewURL(name))
		if doc.Status != 200 || doc.Kind != PreviewKindImage {
			t.Fatalf("%s: kind = %q status = %d (%s), want image/200: a name is a claim, the bytes decide",
				name, doc.Kind, doc.Status, doc.Reason)
		}
		if doc.ContentType != "image/png" {
			t.Fatalf("%s: contentType = %q, want image/png", name, doc.ContentType)
		}
	}
}

// --- AC#1 question 2: bytes, not a path --------------------------------------

func TestPreview192RefusesPathShapedNamesBeforeItOpensAnything(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{})
	ctx := context.Background()
	preview192Put(t, store, "real.png", preview192PNG(64))

	hostile := []string{"", "..", "../real.png", "a/b.png", `a\b.png`, "C:real.png", "real.png.", ".", "a:b"}
	for _, name := range hostile {
		doc := r.Lookup(ctx, PreviewURL(name))
		if doc.Status != 400 {
			t.Fatalf("%q: status = %d (%s), want 400", name, doc.Status, doc.Reason)
		}
		if doc.Kind != "" || doc.TotalBytes != 0 {
			t.Fatalf("%q: the refusal reached the filesystem (kind=%q size=%d); the guard must decide before any open",
				name, doc.Kind, doc.TotalBytes)
		}
	}
	// The one shape whose refusal text is the store's own (a path spelled into
	// a bare name) must say so, so the user is not told "no such artifact".
	for _, name := range []string{"../real.png", `a\b.png`} {
		doc := r.Lookup(ctx, PreviewURL(name))
		if !strings.Contains(doc.Reason, "bare file names only") {
			t.Fatalf("%q: reason %q is not the artifacts store's own guard verdict", name, doc.Reason)
		}
	}
}

func TestPreview192TheResponseNeverCarriesTheArtifactsPath(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{})
	ctx := context.Background()
	preview192Put(t, store, "artifact-a.txt", []byte("一段产物内容\n"))
	preview192Put(t, store, "artifact-b.png", preview192PNG(128))

	for _, name := range []string{"artifact-a.txt", "artifact-b.png"} {
		doc := r.Lookup(ctx, PreviewURL(name))
		if doc.Status != 200 {
			t.Fatalf("%s: status = %d (%s)", name, doc.Status, doc.Reason)
		}
		haystack := string(doc.Body) + strings.Join(doc.Headers, "\n") + doc.Reason + doc.ContentType
		if strings.Contains(haystack, store.ArtifactsDir()) {
			t.Fatalf("%s: the response carries the artifacts directory %q - the browser must never learn where the bytes live",
				name, store.ArtifactsDir())
		}
		if strings.Contains(haystack, `\`) && strings.Contains(haystack, ":") && strings.Contains(haystack, "artifact") {
			// A drive-letter path shape in a response about an artifact is the
			// same leak in another spelling.
			t.Fatalf("%s: the response looks like it carries a filesystem path: %q", name, haystack)
		}
	}
}

func TestPreview192NonRegularAndLinkedFilesAreRefused(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{})
	ctx := context.Background()

	if err := os.MkdirAll(store.ArtifactsDir()+string(os.PathSeparator)+"folder.txt", 0o700); err != nil {
		t.Fatalf("mkdir fixture: %v", err)
	}
	doc := r.Lookup(ctx, PreviewURL("folder.txt"))
	if doc.Status != 403 || !strings.Contains(doc.Reason, "不是普通文件") {
		t.Fatalf("directory: status = %d reason = %q, want 403 naming the shape", doc.Status, doc.Reason)
	}

	target := store.ArtifactsDir() + string(os.PathSeparator) + "real-target.txt"
	if err := os.WriteFile(target, []byte("target\n"), 0o600); err != nil {
		t.Fatalf("write link target: %v", err)
	}
	link := store.ArtifactsDir() + string(os.PathSeparator) + "link.txt"
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("this machine will not let the test create a symlink (%v); the link shape is unmeasurable here", err)
	}
	doc = r.Lookup(ctx, PreviewURL("link.txt"))
	if doc.Status != 403 || !strings.Contains(doc.Reason, "链接") {
		t.Fatalf("link: status = %d reason = %q, want 403 naming the link", doc.Status, doc.Reason)
	}
}

// --- AC#1 questions 3 and 4: truncation, and what the user sees ---------------

func TestPreview192TextFirstViewSaysWhatItHidAndWhereTheRestIs(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{})
	ctx := context.Background()
	full := preview192LargeText()
	preview192Put(t, store, "big-log.txt", full)

	doc := r.Lookup(ctx, PreviewURL("big-log.txt"))
	if doc.Status != 200 || doc.Kind != PreviewKindText {
		t.Fatalf("status = %d kind = %q (%s)", doc.Status, doc.Kind, doc.Reason)
	}
	if doc.Complete {
		t.Fatalf("a %d-byte text came back Complete=true; truncation did not happen", len(full))
	}
	body := string(doc.Body)
	m := preview192NoticeRe.FindStringSubmatch(body)
	if m == nil {
		t.Fatalf("the first-view notice does not state total/head/tail/elided; body starts %q", preview192Head(body, 160))
	}
	total, _ := strconv.Atoi(m[1])
	head, _ := strconv.Atoi(m[2])
	tail, _ := strconv.Atoi(m[3])
	elided, _ := strconv.Atoi(m[4])
	if total != len(full) {
		t.Fatalf("notice total = %d, want %d (the whole file, not what was shown)", total, len(full))
	}
	if head+tail+elided != total {
		t.Fatalf("notice numbers do not add up: %d + %d + %d != %d", head, tail, elided, total)
	}
	if head != int(doc.NextOffset) {
		t.Fatalf("the notice shows head %d but the pointer resumes at %d", head, doc.NextOffset)
	}
	if head != previewHeadBytes {
		t.Fatalf("head = %d, want %d - the configured window, not a coincidence", head, previewHeadBytes)
	}
	// The tail may lose up to utf8.UTFMax-1 bytes to the rune-boundary cut, and
	// it may never lose more: a smaller tail would mean the boundary logic ate
	// characters it did not need to.
	if tail > previewTailBytes || tail < previewTailBytes-(utf8.UTFMax-1) {
		t.Fatalf("tail = %d, want %d (or up to %d fewer for the rune cut)",
			tail, previewTailBytes, utf8.UTFMax-1)
	}
	if !bytes.Contains(doc.Body, full[:head]) {
		t.Fatal("the first view does not contain the file's own head bytes")
	}
	if !bytes.HasSuffix(doc.Body, full[len(full)-tail:]) {
		t.Fatal("the first view does not end with the file's own tail bytes")
	}
	rm := preview192ResumeRe.FindStringSubmatch(body)
	if rm == nil {
		t.Fatalf("the first view carries no resume pointer; body starts %q", preview192Head(body, 160))
	}
	if rm[1] != PreviewURLAt("big-log.txt", int64(head), PreviewMaxTextWindow) {
		t.Fatalf("resume URL = %q, want %q", rm[1], PreviewURLAt("big-log.txt", int64(head), PreviewMaxTextWindow))
	}
	if strings.Contains(body, store.ArtifactsDir()) {
		t.Fatal("the pointer leaked the artifacts directory")
	}
	if !utf8.Valid(doc.Body) {
		t.Fatal("the first view is not valid UTF-8: the head/tail cut split a character")
	}
	if full[doc.NextOffset]&0xC0 == 0x80 {
		t.Fatal("the resume offset points into the middle of a UTF-8 character")
	}
}

func TestPreview192ThePointerWalksEveryByteOfTheFile(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{})
	ctx := context.Background()
	full := preview192LargeText()
	preview192Put(t, store, "walk.txt", full)

	var got []byte
	offset := int64(0)
	steps := 0
	for {
		steps++
		if steps > 64 {
			t.Fatalf("the pointer did not finish the file in 64 windows (offset %d of %d)", offset, len(full))
		}
		doc := r.Lookup(ctx, PreviewURLAt("walk.txt", offset, r.limits.MaxWindow))
		if doc.Status != 200 {
			t.Fatalf("window at %d: status = %d (%s)", offset, doc.Status, doc.Reason)
		}
		if doc.Offset != offset {
			t.Fatalf("window at %d came back reporting Offset=%d", offset, doc.Offset)
		}
		if !utf8.Valid(doc.Body) {
			t.Fatalf("window at %d is not valid UTF-8: a cut split a character", offset)
		}
		got = append(got, doc.Body...)
		if doc.Complete {
			break
		}
		if doc.NextOffset <= offset {
			t.Fatalf("window at %d returned NextOffset=%d: the pointer does not advance, so the tail is unreachable",
				offset, doc.NextOffset)
		}
		offset = doc.NextOffset
	}
	if !bytes.Equal(got, full) {
		t.Fatalf("walking the pointer recovered %d byte(s), want %d - the truncation pointer does not find the whole file",
			len(got), len(full))
	}
	if steps < 4 {
		t.Fatalf("the walk took %d window(s) for %d bytes at %d per window: the fixture is not exercising multi-window reads",
			steps, len(full), r.limits.MaxWindow)
	}
}

func TestPreview192SmallTextComesBackWholeWithNoEnvelope(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{})
	ctx := context.Background()
	small := []byte("短文本，只有一行。\n")
	preview192Put(t, store, "small.txt", small)
	doc := r.Lookup(ctx, PreviewURL("small.txt"))
	if doc.Status != 200 || !doc.Complete {
		t.Fatalf("status = %d complete = %t (%s)", doc.Status, doc.Complete, doc.Reason)
	}
	if !bytes.Equal(doc.Body, small) {
		t.Fatalf("body = %q, want the file's own bytes with no notice glued on", doc.Body)
	}
	if doc.NextOffset != 0 {
		t.Fatalf("a whole file reported a resume pointer %d", doc.NextOffset)
	}
}

func TestPreview192EmptyArtifactIsAnEmptyTextFile(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{})
	ctx := context.Background()
	preview192Put(t, store, "empty.txt", nil)
	doc := r.Lookup(ctx, PreviewURL("empty.txt"))
	if doc.Status != 200 || doc.Kind != PreviewKindText || !doc.Complete {
		t.Fatalf("status = %d kind = %q complete = %t (%s)", doc.Status, doc.Kind, doc.Complete, doc.Reason)
	}
	if doc.TotalBytes != 0 || len(doc.Body) != 0 {
		t.Fatalf("empty artifact reported %d byte(s) of %d", len(doc.Body), doc.TotalBytes)
	}
}

func TestPreview192BinaryOverTheCapIsRefusedNotTruncated(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{ImageMax: 1024})
	ctx := context.Background()
	preview192Put(t, store, "big.png", preview192PNG(4096))
	doc := r.Lookup(ctx, PreviewURL("big.png"))
	if doc.Status != 413 {
		t.Fatalf("status = %d (%s), want 413", doc.Status, doc.Reason)
	}
	if len(doc.Body) != 0 {
		t.Fatalf("an over-cap image came back with %d byte(s) of body - a half image is not a preview", len(doc.Body))
	}
	if len(doc.Headers) != 0 {
		t.Fatalf("a refusal carried %d header line(s); nothing may be rendered from it", len(doc.Headers))
	}
	for _, want := range []string{"4096", "1024"} {
		if !strings.Contains(doc.Reason, want) {
			t.Fatalf("refusal %q does not name %q (both numbers are the honest no)", doc.Reason, want)
		}
	}
}

func TestPreview192UnknownOrMalformedWindowParametersAreRefused(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{})
	ctx := context.Background()
	preview192Put(t, store, "w.txt", []byte("0123456789"))
	for _, q := range []string{"?mist=1", "?offset=-1", "?offset=abc", "?limit=0", "?limit=999999", "?offset=1&offset=2"} {
		doc := r.Lookup(ctx, PreviewURL("w.txt")+q)
		if doc.Status != 400 {
			t.Fatalf("%q: status = %d (%s), want 400 - an option nobody honoured must not read as honoured",
				q, doc.Status, doc.Reason)
		}
	}
	doc := r.Lookup(ctx, "/index.html")
	if doc.Status != 400 {
		t.Fatalf("a non-preview path: status = %d (%s), want 400", doc.Status, doc.Reason)
	}
}

// --- AC#2: the injection face ------------------------------------------------

func TestPreview192HTMLIsSandboxedSoAPreviewCannotSpeak(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{})
	ctx := context.Background()
	preview192Put(t, store, "page.html", []byte("<!DOCTYPE html><html><body><script>parent.postMessage({source:'panel-composer',method:'panel.message.send'},'*')</script></body></html>"))
	doc := r.Lookup(ctx, PreviewURL("page.html"))
	if doc.Status != 200 || doc.Kind != PreviewKindHTML {
		t.Fatalf("status = %d kind = %q (%s)", doc.Status, doc.Kind, doc.Reason)
	}
	policy := ""
	for _, h := range doc.Headers {
		if v, ok := strings.CutPrefix(h, "Content-Security-Policy: "); ok {
			policy = v
		}
	}
	if policy == "" {
		t.Fatalf("the HTML preview carries no CSP; headers = %q", doc.Headers)
	}
	if !strings.HasPrefix(policy, "sandbox") {
		t.Fatalf("CSP = %q: it must start with the sandbox directive, or the previewed page keeps this window's origin and can postMessage as the composer",
			policy)
	}
	if strings.Contains(policy, "allow-scripts") {
		t.Fatalf("CSP = %q: allow-scripts would let a previewed page run script", policy)
	}
	if !strings.Contains(policy, "default-src 'none'") {
		t.Fatalf("CSP = %q: default-src 'none' is what stops the document from loading anything", policy)
	}
	if !strings.Contains(strings.Join(doc.Headers, "\n"), "X-Content-Type-Options: nosniff") {
		t.Fatalf("headers = %q: a document without nosniff can be re-sniffed into another type", doc.Headers)
	}
}

func TestPreview192PDFDeclaresTheHostDependencyInsteadOfClaimingPixels(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{})
	ctx := context.Background()
	preview192Put(t, store, "doc.pdf", []byte("%PDF-1.7\n1 0 obj\n<<>>\nendobj\n%%EOF\n"))
	doc := r.Lookup(ctx, PreviewURL("doc.pdf"))
	if doc.Status != 200 || doc.Kind != PreviewKindPDF {
		t.Fatalf("status = %d kind = %q (%s)", doc.Status, doc.Kind, doc.Reason)
	}
	if doc.Render != PreviewRendersHostPDF {
		t.Fatalf("render = %q, want %q: WebView2's PDF renderer is runtime-dependent, so this stack must not claim it",
			doc.Render, PreviewRendersHostPDF)
	}
	joined := strings.Join(doc.Headers, "\n")
	if !strings.Contains(joined, "Content-Type: application/pdf") {
		t.Fatalf("headers = %q, want application/pdf", doc.Headers)
	}
	if !strings.Contains(joined, "Content-Disposition: inline") {
		t.Fatalf("headers = %q, want an inline disposition so a host that can render it does", doc.Headers)
	}
}

func TestPreview192PreviewIsNotAMessageRoute(t *testing.T) {
	// (a) The preview file cannot reach the model: its imports and literals
	// are parsed, not trusted.
	src, err := os.ReadFile("preview_192.go")
	if err != nil {
		t.Fatalf("read preview_192.go: %v", err)
	}
	file, err := parser.ParseFile(token.NewFileSet(), "preview_192.go", src, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse preview_192.go: %v", err)
	}
	forbidden := []string{
		"github.com/CarlosShao/wisp/internal/agent",
		"github.com/CarlosShao/wisp/internal/tools",
		"github.com/CarlosShao/wisp/internal/llm",
	}
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		for _, bad := range forbidden {
			if path == bad {
				t.Fatalf("preview_192.go imports %s: the preview stack must not be able to reach the model", path)
			}
		}
	}
	for _, literal := range []string{"MethodMessageSend", "HandleMessageRequest", "ComposerDispatch"} {
		if bytes.Contains(src, []byte(literal)) {
			t.Fatalf("preview_192.go names %s: a preview must never be a route into a model-bound message", literal)
		}
	}
	// (b) The value type carries nothing a message route consumes.
	typ := reflect.TypeOf(PreviewDocument{})
	for i := 0; i < typ.NumField(); i++ {
		ft := typ.Field(i).Type
		if ft == reflect.TypeOf(ComposerRequest{}) || ft == reflect.TypeOf(&ComposerRequest{}) {
			t.Fatalf("PreviewDocument.%s is a ComposerRequest: preview output must not be a request envelope", typ.Field(i).Name)
		}
	}
	// (c) And the composer pipe refuses a preview method: the only inbound
	// door the page has does not answer "preview", while it does answer the
	// message route (the positive control that the door under test is alive).
	if _, err := ParseComposerRequest(`{"method":"panel.preview","source":"panel-composer","requestId":"rid-1"}`); err == nil {
		t.Fatal("ParseComposerRequest accepted a preview method: preview would be a second inbound door")
	}
	if _, err := ParseComposerRequest(`{"method":"` + MethodMessageSend + `","source":"panel-composer","requestId":"rid-2"}`); err != nil {
		t.Fatalf("positive control failed: the message route is not even parseable (%v), so the refusal above proves nothing", err)
	}
}

func TestPreview192QuoteDoorStampsAndTheEgressGateThenCatches(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{})
	ctx := context.Background()
	secret := "这一段内容来自产物，只应该留在面板里"
	preview192Put(t, store, "quoted.txt", []byte("前缀\n"+secret+"\n后缀\n"))

	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	scope := prov.OpenScope("preview-192")
	defer func() { _ = scope.Close() }()

	got, err := r.QuoteForModel(ctx, prov, "preview-192", "quoted.txt", secret)
	if err != nil {
		t.Fatalf("QuoteForModel: %v", err)
	}
	if got != secret {
		t.Fatalf("QuoteForModel returned %q, want the quote unchanged", got)
	}
	// The mark is live: the same string now trips the egress scan on a real
	// channel (the exfil half of the same pipe), naming the artifact door.
	hit, ok := prov.CheckText("preview-192", risk.ChHTTP, secret)
	if !ok {
		t.Fatal("after the quote the C25 egress scan found nothing: the stamp is not live")
	}
	if hit.SrcTool != PreviewSourceTool {
		t.Fatalf("hit source = %q, want %q", hit.SrcTool, PreviewSourceTool)
	}
	if hit.Origin != "artifact:quoted.txt" {
		t.Fatalf("hit origin = %q, want artifact:quoted.txt", hit.Origin)
	}
	if hit, ok := prov.Inspect("preview-192", "notify", map[string]any{"text": secret}); !ok || hit.SrcTool != PreviewSourceTool {
		t.Fatalf("Inspect on the notify channel: hit = %+v ok = %t, want a hit from %q", hit, ok, PreviewSourceTool)
	}
}

func TestPreview192QuoteDoorRefusesTextThatIsNotInTheArtifact(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{})
	ctx := context.Background()
	preview192Put(t, store, "quoted.txt", []byte("产物里真实存在的这一行\n"))
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	scope := prov.OpenScope("preview-192-neg")
	defer func() { _ = scope.Close() }()

	laundered := "页面自己编的一句话，产物里没有"
	if _, err := r.QuoteForModel(ctx, prov, "preview-192-neg", "quoted.txt", laundered); err == nil {
		t.Fatal("the quote door stamped text the artifact does not contain: the page could choose what the C25 index believes")
	}
	if _, ok := prov.CheckText("preview-192-neg", risk.ChHTTP, laundered); ok {
		t.Fatal("the refused quote still left a mark")
	}
}

func TestPreview192QuoteDoorRefusesWhatItCannotStamp(t *testing.T) {
	store := preview192Store(t)
	r := preview192Resolver(t, store, PreviewLimits{QuoteMax: 64})
	ctx := context.Background()
	big := bytes.Repeat([]byte("x"), 256)
	preview192Put(t, store, "small.txt", []byte("内容\n"))
	preview192Put(t, store, "big.txt", big)
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	scope := prov.OpenScope("preview-192-refuse")
	defer func() { _ = scope.Close() }()

	cases := []struct {
		name     string
		marker   TaintMarker
		scopeID  string
		artifact string
		quoted   string
		want     string
	}{
		{"no engine", nil, "preview-192-refuse", "small.txt", "内容", "C25"},
		{"no scope", prov, "  ", "small.txt", "内容", "scope"},
		{"empty quote", prov, "preview-192-refuse", "small.txt", "  ", "空的"},
		{"bad artifact name", prov, "preview-192-refuse", "../small.txt", "内容", "不合格"},
		{"missing artifact", prov, "preview-192-refuse", "nope.txt", "内容", "读不到"},
		{"over the check cap", prov, "preview-192-refuse", "big.txt", "xxxx", "上限"},
	}
	for _, tc := range cases {
		_, err := r.QuoteForModel(ctx, tc.marker, tc.scopeID, tc.artifact, tc.quoted)
		if err == nil {
			t.Fatalf("%s: the door accepted a quote it cannot stamp", tc.name)
		}
		if !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("%s: refusal %q does not name %q", tc.name, err.Error(), tc.want)
		}
	}
}

func TestPreview192SourceToolIsTheContractName(t *testing.T) {
	if PreviewSourceTool != risk.SrcTaskOutput {
		t.Fatalf("PreviewSourceTool = %q, but the contract's source name is %q: a preview mark must not invent a source",
			PreviewSourceTool, risk.SrcTaskOutput)
	}
}

func TestPreview192ResolverRefusesToBeBuiltWithoutItsGuards(t *testing.T) {
	store := preview192Store(t)
	if _, err := NewPreviewResolver(nil, memory.ValidArtifactName); err == nil {
		t.Fatal("a resolver with no artifacts directory was built")
	}
	if _, err := NewPreviewResolver(store, nil); err == nil {
		t.Fatal("a resolver with no name guard was built: that is a read-any-path door")
	}
	if _, err := NewPreviewResolver(emptyArtifactsSource{}, memory.ValidArtifactName); err == nil {
		t.Fatal("a resolver whose artifacts directory is empty was built")
	}
}

// emptyArtifactsSource is the one synthetic value in this file: it exists only
// to give NewPreviewResolver an empty directory string, which no production
// store can produce.
type emptyArtifactsSource struct{}

func (emptyArtifactsSource) ArtifactsDir() string { return "  " }

// preview192Head is a width-limited excerpt for failure messages.
func preview192Head(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
