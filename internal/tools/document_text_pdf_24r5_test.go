package tools

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 24 / leg 24-r5: the PDF branch of doc.read now extracts through the
// whitelisted pure-Go parser github.com/ledongthuc/pdf instead of answering the
// DEFERRED(D34-1) stub. These tests pin the four shapes the ticket names: a text
// PDF, a multi-page PDF against the caller's byte budget, a scan (no text layer)
// and a broken file. The fixture writer below only emits bytes; it is not a
// parser and reads nothing back.

// r5Esc escapes a literal string for a PDF parenthesised string.
func r5Esc(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `(`, `\(`, `)`, `\)`)
	return r.Replace(s)
}

// r5PDF builds a minimal uncompressed one-object-per-page PDF whose pages each
// draw a single text show op. An empty page text draws no Tj at all, which is
// how a scan-without-text-layer is represented here. Returns the file path.
func r5PDF(t *testing.T, pageTexts ...string) string {
	t.Helper()

	const (
		catalogObj = 1
		pagesObj   = 2
	)
	n := len(pageTexts)
	if n == 0 {
		t.Fatal("r5PDF needs at least one page")
	}
	fontObj := 3 + 2*n

	buf := &bytes.Buffer{}
	buf.WriteString("%PDF-1.4\n")
	offsets := map[int]int64{}
	writeObj := func(num int, body string) {
		offsets[num] = int64(buf.Len())
		fmt.Fprintf(buf, "%d 0 obj\n%s\nendobj\n", num, body)
	}

	kids := make([]string, n)
	for i := range pageTexts {
		kids[i] = fmt.Sprintf("%d 0 R", 3+2*i)
	}
	writeObj(catalogObj, fmt.Sprintf("<< /Type /Catalog /Pages %d 0 R >>", pagesObj))
	writeObj(pagesObj, fmt.Sprintf("<< /Type /Pages /Kids [%s] /Count %d >>",
		strings.Join(kids, " "), n))

	contentStream := func(text string) string {
		body := "BT\n/F1 12 Tf\n72 720 Td\n"
		if text != "" {
			body += "(" + r5Esc(text) + ") Tj\n"
		}
		body += "ET"
		return fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(body), body)
	}
	for i, txt := range pageTexts {
		pageObj := 3 + 2*i
		contentObj := pageObj + 1
		writeObj(pageObj, fmt.Sprintf(
			"<< /Type /Page /Parent %d 0 R /Resources << /Font << /F1 %d 0 R >> >> "+
				"/MediaBox [0 0 612 792] /Contents %d 0 R >>", pagesObj, fontObj, contentObj))
		writeObj(contentObj, contentStream(txt))
	}
	writeObj(fontObj, "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>")

	xrefOff := int64(buf.Len())
	fmt.Fprintf(buf, "xref\n0 %d\n0000000000 65535 f \n", fontObj+1)
	for num := 1; num <= fontObj; num++ {
		fmt.Fprintf(buf, "%010d 00000 n \n", offsets[num])
	}
	fmt.Fprintf(buf, "trailer\n<< /Size %d /Root %d 0 R >>\nstartxref\n%d\n%%%%EOF\n",
		fontObj+1, catalogObj, xrefOff)

	path := filepath.Join(t.TempDir(), "fixture.pdf")
	if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
		t.Fatalf("write fixture pdf: %v", err)
	}
	return path
}

// Test24R5PDFTextLayerExtracts is the ticket's "text PDF" acceptance line: a real
// PDF with a text layer comes back as text, never as a refusal.
func Test24R5PDFTextLayerExtracts(t *testing.T) {
	got, err := pdfText{}.ExtractText(context.Background(), r5PDF(t, "Hello from Wisp doc.read"), 0)
	if err != nil {
		t.Fatalf("a text-layer PDF must extract, got error: %v", err)
	}
	if !strings.Contains(got.Text, "Hello from Wisp doc.read") {
		t.Errorf("extracted text = %q, want the drawn string", got.Text)
	}
	if got.Runs == 0 {
		t.Errorf("a page with a text run must report a non-zero Runs, got 0")
	}
	if got.Truncated {
		t.Error("a one-page document under the default budget must not be truncated")
	}
	// Reverse criterion: the retired stub answered ErrUnsupported. Nothing on this
	// path may still route through that seam.
	if _, err := (pdfText{}).ExtractText(context.Background(), r5PDF(t, "x"), 0); err != nil {
		t.Errorf("no error expected, got %v", err)
	}
}

// Test24R5PDFMultiPageConcatenatesAndHonoursTheBudget pins the multi-page shape
// plus the C15/D15(3) contract: the extractor caps at the caller's byte budget
// and says so with Truncated, and never spills on its own.
func Test24R5PDFMultiPageConcatenatesAndHonoursTheBudget(t *testing.T) {
	long := func(m string) string { return m + " " + strings.Repeat("body", 40) }
	path := r5PDF(t, long("FIRSTPAGE"), long("SECONDAGE"), long("THIRDPAGE"))

	full, err := pdfText{}.ExtractText(context.Background(), path, 0)
	if err != nil {
		t.Fatalf("multi-page extraction: %v", err)
	}
	if full.Truncated {
		t.Error("the full read under the default budget must not be truncated")
	}
	if full.Runs != 3 {
		t.Errorf("Runs = %d, want one per page (3)", full.Runs)
	}
	i1 := strings.Index(full.Text, "FIRSTPAGE")
	i2 := strings.Index(full.Text, "SECONDAGE")
	i3 := strings.Index(full.Text, "THIRDPAGE")
	if i1 < 0 || i2 < 0 || i3 < 0 {
		t.Fatalf("all three pages must appear, got %q", full.Text)
	}
	if !(i1 < i2 && i2 < i3) {
		t.Errorf("pages must be concatenated in page order, got indexes %d/%d/%d", i1, i2, i3)
	}

	// A budget smaller than one page must stop early and flag it, and must never
	// hand back more text than the caller allowed.
	limited, err := pdfText{}.ExtractText(context.Background(), path, 40)
	if err != nil {
		t.Fatalf("budgeted extraction: %v", err)
	}
	if int64(len(limited.Text)) > 40 {
		t.Errorf("budget violated: %d bytes under a 40 byte limit", len(limited.Text))
	}
	if !limited.Truncated {
		t.Error("stopping at the budget must be reported as Truncated, not as a clean read")
	}
	if !strings.Contains(limited.Text, "FIRSTPAGE") {
		t.Errorf("the truncated head must still be the first page, got %q", limited.Text)
	}
}

// Test24R5PDFScanIsMarkedNotGuessed pins the ticket's scanned-PDF rule: a PDF
// with no text layer owes its reader an explicit marker naming the missing text
// layer and the deferred OCR, and must never come back as an empty success.
func Test24R5PDFScanIsMarkedNotGuessed(t *testing.T) {
	path := r5PDF(t, "", "") // two pages, neither draws any text

	got, err := pdfText{}.ExtractText(context.Background(), path, 0)
	if err == nil {
		t.Fatalf("a no-text-layer PDF must not read as a success, got %+v", got)
	}
	if got.Text != "" || got.Runs != 0 {
		t.Errorf("the refused extraction must carry no text, got %+v", got)
	}
	if !strings.Contains(err.Error(), scannedPDFReason) {
		t.Errorf("error = %v, want it to carry %q", err, scannedPDFReason)
	}
	if !strings.Contains(err.Error(), "扫描件") {
		t.Errorf("the marker must say this looks like a scan, got %v", err)
	}
	if strings.Contains(err.Error(), "DEFERRED(D34-1)") {
		t.Errorf("the retired stub marker must be gone from every path, got %v", err)
	}

	// And through the tool: IsError, with the reason in the answer (no silent
	// empty result reaching the model).
	dir := filepath.Dir(path)
	p, canon := r2aCanon(t, []string{dir}, path)
	res, err := docRead{d: DocDeps{Paths: p}}.Execute(context.Background(),
		r2aParams(t, map[string]any{"path": canon}), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.IsError {
		t.Errorf("a scan must come back IsError, got %+v", res)
	}
	if !strings.Contains(res.Text, scannedPDFReason) {
		t.Errorf("the answer must name the no-extractable-text reason, got %q", res.Text)
	}
}

// Test24R5PDFBrokenFailsClosed covers the "坏 PDF -> error" rule: an unreadable
// or non-PDF file is a parse failure, never laundered into a deferral or an
// empty result.
func Test24R5PDFBrokenFailsClosed(t *testing.T) {
	dir := t.TempDir()

	mk := func(name, body string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
		return p
	}
	cases := []struct {
		name    string
		path    string
		wantSub string
	}{
		{"garbage", mk("garbage.pdf", "this is not a pdf at all"), "不是有效的 PDF"},
		{"header-only", mk("short.pdf", "%PDF-1.7"), "不是有效的 PDF"},
		{"missing", filepath.Join(dir, "absent.pdf"), "无法打开文件"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := pdfText{}.ExtractText(context.Background(), tc.path, 0)
			if err == nil {
				t.Fatalf("a broken pdf must error, got %+v", got)
			}
			if got.Text != "" {
				t.Errorf("a failed extraction must carry no text, got %q", got.Text)
			}
			if !strings.Contains(err.Error(), tc.wantSub) {
				t.Errorf("error = %v, want it to name %q", err, tc.wantSub)
			}
			if strings.Contains(err.Error(), "DEFERRED") || errorsIsUnsupported(err) {
				t.Errorf("a broken file must not be reported as a deferral: %v", err)
			}
		})
	}
}

func errorsIsUnsupported(err error) bool {
	return strings.Contains(err.Error(), ErrUnsupported.Error())
}

// Test24R5DocReadEndToEndOnPDF is the acceptance line "doc.read: text PDF": the
// dispatch, the taint Origin and the extraction all line up.
func Test24R5DocReadEndToEndOnPDF(t *testing.T) {
	path := r5PDF(t, "wisp ticket twenty-four pdf body")
	dir := filepath.Dir(path)
	p, canon := r2aCanon(t, []string{dir}, path)

	res, err := docRead{d: DocDeps{Paths: p}}.Execute(context.Background(),
		r2aParams(t, map[string]any{"path": canon}), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if res.IsError {
		t.Fatalf("a text PDF must succeed, got %q", res.Text)
	}
	if !strings.Contains(res.Text, "wisp ticket twenty-four pdf body") {
		t.Errorf("body = %q, want the drawn text", res.Text)
	}
	if res.Origin != canon {
		t.Errorf("Origin = %q, want the canonical path %q (taint marker)", res.Origin, canon)
	}

	// 还原核: the format is dispatched as a real reader, so the "unsupported
	// format" refusal no longer fires for .pdf.
	s, err := NewExtractorSet(DefaultExtractors()...)
	if err != nil {
		t.Fatal(err)
	}
	if _, ext, ok := s.Lookup(canon); !ok || ext != "pdf" {
		t.Errorf("Lookup must resolve the fixture to pdf, got ext %q ok %v", ext, ok)
	}
}

// Test24R5DeferralMarkersAreGone is the reverse criterion on the source text:
// the PDF deferral marker must no longer appear anywhere in the implementation,
// while the two markers this leg did not touch must survive.
func Test24R5DeferralMarkersAreGone(t *testing.T) {
	src, err := os.ReadFile("document_text.go")
	if err != nil {
		t.Fatalf("read the implementation: %v", err)
	}
	text := string(src)
	if strings.Contains(text, "D34-1") {
		t.Error("DEFERRED(D34-1) must be deleted now that the PDF branch extracts")
	}
	for _, keep := range []string{"DEFERRED(xlsx, D34)", "DEFERRED(OCR, D34)"} {
		if !strings.Contains(text, keep) {
			t.Errorf("%q must stay: SPEC-12 §5 still registers it", keep)
		}
	}
}
