package tools

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// Ticket 24 / limb 24-r2a: the zero-dependency half (docx via archive/zip +
// encoding/xml, search.content) plus the TextExtractor seam. The PDF position
// was a stub when this limb landed; since 24-r5 it is the real ledongthuc/pdf
// reader, and 24-r6 replaced the two stub-era PDF assertions below with the
// post-extraction criteria. Every reverse criterion the ticket names is pinned
// here: a broken docx never reads as success, a broken or missing PDF fails
// closed by name, an out-of-allowlist root is refused, a blacklisted file never
// appears even in a listing, and a cancelled scan comes back marked truncated
// instead of clean.

// r2aDocxBody is a hand-written WordprocessingML body: three paragraphs, a
// tab, a hard break, a preserved-space run and one property element that must
// NOT become text.
const r2aDocxBody = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
	`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body>` +
	`<w:p><w:r><w:t xml:space="preserve">第一段 hello</w:t></w:r></w:p>` +
	`<w:p><w:r><w:t>第二段</w:t><w:tab/><w:t xml:space="preserve"> tail</w:t></w:r>` +
	`<w:r><w:br/><w:br/><w:t>第三段</w:t></w:r></w:p>` +
	`<w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:t>第四段</w:t></w:r></w:p>` +
	`</w:body></w:document>`

// r2aDocx writes a minimal .docx package - the fixture is built with the same
// stdlib zip writer the reader is paired with, so there is no binary blob to
// ship - and returns its path.
func r2aDocx(t *testing.T, body string, extra map[string]string) string {
	t.Helper()
	name := filepath.Join(t.TempDir(), "sample.docx")
	f, err := os.Create(name)
	if err != nil {
		t.Fatalf("create fixture: %v", err)
	}
	zw := zip.NewWriter(f)
	w, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatalf("create body entry: %v", err)
	}
	if _, err := w.Write([]byte(body)); err != nil {
		t.Fatalf("write body: %v", err)
	}
	for n, v := range extra {
		e, err := zw.Create(n)
		if err != nil {
			t.Fatalf("create extra entry %s: %v", n, err)
		}
		if _, err := e.Write([]byte(v)); err != nil {
			t.Fatalf("write extra entry %s: %v", n, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close file: %v", err)
	}
	return name
}

func r2aParams(t *testing.T, v any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal params: %v", err)
	}
	return b
}

// r2aCanon builds the C26 resolver over the given allowed dirs and canonicalizes
// one path through it, so the test drives the same spelling the gate judged.
func r2aCanon(t *testing.T, allowed []string, raw string) (*PathCanonicalizer, string) {
	t.Helper()
	p := NewPathCanonicalizer(allowed, nil)
	if p == nil {
		t.Fatal("NewPathCanonicalizer returned nil")
	}
	c, err := p.Canonicalize(raw)
	if err != nil {
		t.Fatalf("C26 canonicalize %q: %v", raw, err)
	}
	return p, c
}

// ---------------------------------------------------------------------------
// docx extraction (stdlib archive/zip + encoding/xml, zero new dependency)
// ---------------------------------------------------------------------------

func Test24R2ADocxExtractsRunsTabsBreaks(t *testing.T) {
	path := r2aDocx(t, r2aDocxBody, map[string]string{"[Content_Types].xml": "<Types/>"})
	got, err := docxText{}.ExtractText(context.Background(), path, 0)
	if err != nil {
		t.Fatalf("ExtractText: %v", err)
	}
	if got.Runs != 5 {
		t.Errorf("Runs = %d, want 5 (five w:t runs in the fixture)", got.Runs)
	}
	if !strings.HasPrefix(got.Text, "第一段 hello\n第二段\t tail\n") {
		t.Errorf("run/tab order wrong: %q", got.Text)
	}
	if !strings.Contains(got.Text, "\n第三段\n") || !strings.HasSuffix(got.Text, "第四段\n") {
		t.Errorf("break/paragraph newlines wrong: %q", got.Text)
	}
	if strings.Contains(got.Text, "center") {
		t.Error("a w:jc property value leaked into the text")
	}
	if got.Truncated {
		t.Error("Truncated set on a document that fit the budget")
	}
}

func Test24R2ADocxHonoursByteBudget(t *testing.T) {
	path := r2aDocx(t, r2aDocxBody, nil)
	got, err := docxText{}.ExtractText(context.Background(), path, 12)
	if err != nil {
		t.Fatalf("ExtractText: %v", err)
	}
	if !got.Truncated {
		t.Error("Truncated must be true when the budget stopped the walk")
	}
	if int64(len(got.Text)) > 12 {
		t.Errorf("text is %d bytes, over the 12-byte budget", len(got.Text))
	}
}

func Test24R2ADocxBrokenPackagesFailClosed(t *testing.T) {
	cases := []struct {
		name    string
		write   func(t *testing.T) string
		wantSub string
	}{
		{
			name: "not a zip",
			write: func(t *testing.T) string {
				p := filepath.Join(t.TempDir(), "fake.docx")
				if err := os.WriteFile(p, []byte("this is not a zip container"), 0o600); err != nil {
					t.Fatal(err)
				}
				return p
			},
			wantSub: "不是有效的 OOXML 包",
		},
		{
			name: "zip without the body part",
			write: func(t *testing.T) string {
				name := filepath.Join(t.TempDir(), "nobody.docx")
				f, err := os.Create(name)
				if err != nil {
					t.Fatal(err)
				}
				zw := zip.NewWriter(f)
				w, err := zw.Create("[Content_Types].xml")
				if err != nil {
					t.Fatal(err)
				}
				if _, err := w.Write([]byte("<Types/>")); err != nil {
					t.Fatal(err)
				}
				if err := zw.Close(); err != nil {
					t.Fatal(err)
				}
				if err := f.Close(); err != nil {
					t.Fatal(err)
				}
				return name
			},
			wantSub: "缺少 word/document.xml",
		},
		{
			name:    "malformed body xml",
			write:   func(t *testing.T) string { return r2aDocx(t, "<w:document><w:p><w:r><w:t>oops", nil) },
			wantSub: "解析失败",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := docxText{}.ExtractText(context.Background(), tc.write(t), 0)
			if err == nil {
				t.Fatalf("a broken package must not succeed; got text %q", got.Text)
			}
			if !strings.Contains(err.Error(), tc.wantSub) {
				t.Errorf("error = %v, want it to name %q", err, tc.wantSub)
			}
			if got.Text != "" {
				t.Errorf("a failed extraction must carry no text, got %q", got.Text)
			}
		})
	}
}

// Test24R2ADocxNoTextIsNotASilentSuccess pins the ticket's "no silent empty
// result" rule (ticket:20-21) on both sides of the seam: the reader reports zero
// runs, and doc.read turns that into an error instead of an empty success.
func Test24R2ADocxNoTextIsNotASilentSuccess(t *testing.T) {
	empty := `<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p/></w:body></w:document>`
	path := r2aDocx(t, empty, nil)
	got, err := docxText{}.ExtractText(context.Background(), path, 0)
	if err != nil {
		t.Fatalf("the package itself is valid: %v", err)
	}
	if got.Runs != 0 || strings.TrimSpace(got.Text) != "" {
		t.Fatalf("fixture must yield zero text, got runs=%d text=%q", got.Runs, got.Text)
	}
	p, canon := r2aCanon(t, []string{filepath.Dir(path)}, path)
	res, err := docRead{d: DocDeps{Paths: p}}.Execute(context.Background(),
		r2aParams(t, map[string]any{"path": canon}), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.IsError {
		t.Errorf("a zero-text document must come back IsError, got %q", res.Text)
	}
	if !strings.Contains(res.Text, scannedPDFReason) {
		t.Errorf("the answer must carry the no-extractable-text reason %q, got %q", scannedPDFReason, res.Text)
	}
}

// ---------------------------------------------------------------------------
// TextExtractor seam: format dispatch, the PDF branch (24-r5 real reader,
// re-pinned by 24-r6; xlsx remains refused, DEFERRED(xlsx, D34))
// ---------------------------------------------------------------------------

func Test24R2AExtractorSetDispatch(t *testing.T) {
	s, err := NewExtractorSet(DefaultExtractors()...)
	if err != nil {
		t.Fatalf("NewExtractorSet: %v", err)
	}
	if formats := strings.Join(s.Formats(), ","); formats != "docx,pdf" {
		t.Errorf("Formats() = %q, want docx,pdf", formats)
	}
	if _, ext, ok := s.Lookup(`C:\x\a.docx`); !ok || ext != "docx" {
		t.Errorf(`Lookup(C:\x\a.docx) = ext %q ok %v`, ext, ok)
	}
	if _, ext, ok := s.Lookup(`C:\x\A.PDF`); !ok || ext != "pdf" {
		t.Errorf("Lookup must be case-insensitive: ext %q ok %v", ext, ok)
	}
	// DEFERRED(xlsx, D34): no extractor is registered, so it is refused.
	if _, _, ok := s.Lookup(`C:\x\budget.xlsx`); ok {
		t.Error("xlsx must have no extractor until SPEC-12 §5 releases it")
	}
	if _, _, ok := s.Lookup(`C:\x\notes`); ok {
		t.Error("an extension-less path must resolve to no extractor")
	}
	if _, err := NewExtractorSet(docxText{}, docxText{}); err == nil {
		t.Error("a duplicate extractor for one format must be refused at wiring time")
	}
}

// Test24R6PDFBranchIsARealReaderNotAStub replaces the 24-r2a assertion that the
// PDF position answered ErrUnsupported with a DEFERRED(D34-1) marker: 24-r5
// landed the real ledongthuc/pdf reader, so the true shape at this seam is 有效
// PDF 抽出正文、坏 PDF fail-closed 具名报错 — never ErrUnsupported, never the
// retired marker.
func Test24R6PDFBranchIsARealReaderNotAStub(t *testing.T) {
	// 有效 PDF: the seam extracts the drawn body instead of refusing.
	path := r5PDF(t, "ticket 24 pdf seam body")
	got, err := pdfText{}.ExtractText(context.Background(), path, 0)
	if err != nil {
		t.Fatalf("a text-layer PDF must extract, got error: %v", err)
	}
	if !strings.Contains(got.Text, "ticket 24 pdf seam body") || got.Runs == 0 {
		t.Errorf("extraction must carry the drawn body and non-zero runs, got %+v", got)
	}
	if got.Truncated {
		t.Error("a one-page fixture under the default budget must not be truncated")
	}

	// 坏 PDF（打不开）: a named fail-closed error, not the retired stub seam.
	miss, err := pdfText{}.ExtractText(context.Background(), filepath.Join(t.TempDir(), "absent.pdf"), 0)
	if err == nil {
		t.Fatalf("a missing PDF must error, got %+v", miss)
	}
	if errors.Is(err, ErrUnsupported) {
		t.Errorf("a broken file must not be laundered into ErrUnsupported: %v", err)
	}
	if !strings.Contains(err.Error(), "无法打开文件") {
		t.Errorf("error = %v, want it to name the open failure", err)
	}
	if strings.Contains(err.Error(), "DEFERRED(D34-1)") {
		t.Errorf("the retired deferral marker must be gone from the error path: %v", err)
	}
	if miss.Text != "" || miss.Runs != 0 {
		t.Errorf("a failed extraction must carry no text, got %+v", miss)
	}
}

// Test24R6DocReadPDFSucceedsAndBrokenPDFFailsClosed replaces the 24-r2a refusal
// assertion (deferred format ⇒ IsError carrying the DEFERRED(D34-1) marker).
// Post-extraction the true criteria are: 有效 PDF 经 doc.read 端到端成功并带上
// C25 打标 Origin; 坏 PDF ⇒ IsError 具名解析失败, which is neither a success
// nor a "该格式暂未实现" deferral.
func Test24R6DocReadPDFSucceedsAndBrokenPDFFailsClosed(t *testing.T) {
	good := r5PDF(t, "doc.read end to end pdf body")
	p1, canon1 := r2aCanon(t, []string{filepath.Dir(good)}, good)
	res, err := docRead{d: DocDeps{Paths: p1}}.Execute(context.Background(),
		r2aParams(t, map[string]any{"path": canon1}), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if res.IsError {
		t.Fatalf("a text-layer PDF must read as success, got %q", res.Text)
	}
	if !strings.Contains(res.Text, "doc.read end to end pdf body") {
		t.Errorf("body = %q, want the drawn text", res.Text)
	}
	if res.Origin != canon1 {
		t.Errorf("Origin = %q, want the canonical path %q (the C25 taint seam)", res.Origin, canon1)
	}

	dir := t.TempDir()
	bad := filepath.Join(dir, "report.pdf")
	if err := os.WriteFile(bad, []byte("%PDF-1.7 fake"), 0o600); err != nil {
		t.Fatal(err)
	}
	p2, canon2 := r2aCanon(t, []string{dir}, bad)
	badRes, err := docRead{d: DocDeps{Paths: p2}}.Execute(context.Background(),
		r2aParams(t, map[string]any{"path": canon2}), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !badRes.IsError {
		t.Fatalf("a broken PDF must never read as success: %+v", badRes)
	}
	if !strings.Contains(badRes.Text, "解析失败") || !strings.Contains(badRes.Text, "不是有效的 PDF") {
		t.Errorf("the refusal must name the parse failure, got %q", badRes.Text)
	}
	if strings.Contains(badRes.Text, "DEFERRED(D34-1)") || strings.Contains(badRes.Text, "该格式暂未实现") {
		t.Errorf("a broken file must not be reported as a deferral, got %q", badRes.Text)
	}
	if badRes.Origin != canon2 {
		t.Errorf("Origin = %q, want the canonical path %q (the taint seam rides errors too)", badRes.Origin, canon2)
	}
}

func Test24R2ADocReadRefusesBlacklistedPathBeforeOpening(t *testing.T) {
	path := r2aDocx(t, r2aDocxBody, nil)
	p, canon := r2aCanon(t, []string{filepath.Dir(path)}, path)
	consulted := false
	d := DocDeps{
		Paths: p,
		Classify: func(string) risk.Class {
			consulted = true
			return risk.ClassA
		},
	}
	res, err := docRead{d: d}.Execute(context.Background(),
		r2aParams(t, map[string]any{"path": canon}), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.IsError || !strings.Contains(res.Text, "敏感清单") {
		t.Errorf("an A-list document must be refused, got %+v", res)
	}
	if !consulted {
		t.Error("the classifier was never consulted, so the veto has no teeth")
	}
	if strings.Contains(res.Text, "hello") {
		t.Error("the refused document leaked its text")
	}
}

func Test24R2ADocReadFailClosedWithoutResolver(t *testing.T) {
	res, err := docRead{d: DocDeps{}}.Execute(context.Background(),
		r2aParams(t, map[string]any{"path": `C:\x\a.docx`}), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.IsError || !strings.Contains(res.Text, "fail-closed") {
		t.Errorf("no C26 resolver must fail closed, got %+v", res)
	}
}

func Test24R2ADocReadExtractsADocxEndToEnd(t *testing.T) {
	path := r2aDocx(t, r2aDocxBody, nil)
	p, canon := r2aCanon(t, []string{filepath.Dir(path)}, path)
	progress := ""
	res, err := docRead{d: DocDeps{Paths: p}}.Execute(context.Background(),
		r2aParams(t, map[string]any{"path": canon}), func(delta string) { progress = delta })
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if res.IsError {
		t.Fatalf("docx read failed: %s", res.Text)
	}
	if !strings.Contains(res.Text, "第四段") || !strings.Contains(res.Text, "第一段 hello") {
		t.Errorf("text = %q, want the fixture paragraphs", res.Text)
	}
	if res.Origin != canon {
		t.Errorf("Origin = %q, want the canonical path %q (the C25 taint seam)", res.Origin, canon)
	}
	if progress == "" {
		t.Error("onUpdate never streamed a progress line")
	}
}

// ---------------------------------------------------------------------------
// search.content
// ---------------------------------------------------------------------------

// r2aCorpus writes a tiny corpus and returns its root. Five of the six .txt/.md
// files hold the needle; one is only there to be excluded by the blacklist.
func r2aCorpus(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	sub := filepath.Join(root, "notes")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		filepath.Join(root, "a.txt"):         "first line\nSECRET needle here\nthird\n",
		filepath.Join(root, "b.md"):          "header\nAnother NEEDLE in caps\n",
		filepath.Join(root, "bin.txt"):       "needle\x00binary\n",
		filepath.Join(sub, "c.log"):          "sub dir needle hit\n",
		filepath.Join(sub, "skip.docx"):      "needle inside a document is not scanned by extension\n",
		filepath.Join(root, "blacklist.txt"): "needle in a blacklisted file\n",
	}
	for name, body := range files {
		if err := os.WriteFile(name, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func r2aScan(d SearchDeps, root string, limit int) *searchScan {
	return &searchScan{d: d, needle: "needle", limit: limit, textExt: d.exts()}
}

func Test24R2ASearchContentFindsMatches(t *testing.T) {
	root := r2aCorpus(t)
	d := SearchDeps{Classify: func(string) risk.Class { return risk.ClassNone }}
	out, err := r2aScan(d, root, 10).run(context.Background(), root)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(out.matches) != 4 {
		t.Fatalf("matches = %d, want 4 (a.txt, b.md case-insensitive, notes/c.log, blacklist.txt): %+v",
			len(out.matches), out.matches)
	}
	if out.hitResultCap || out.hitFileCap || out.aborted != "" {
		t.Errorf("a finished scan must report no cap: %+v", out)
	}
	if out.excluded != 0 {
		t.Errorf("excluded = %d, want 0 with an empty blacklist", out.excluded)
	}
	var joined string
	for _, m := range out.matches {
		joined += filepath.Base(m.path) + ":" + m.snippet + "|"
	}
	if !strings.Contains(joined, "a.txt:SECRET needle here") {
		t.Errorf("missing the a.txt hit: %s", joined)
	}
	if !strings.Contains(joined, "b.md:Another NEEDLE in caps") {
		t.Errorf("the match is not case-insensitive: %s", joined)
	}
	if strings.Contains(joined, "skip.docx") {
		t.Error("a non-text extension was scanned")
	}
	// bin.txt holds a NUL byte: the content veto must catch it.
	if strings.Contains(joined, "bin.txt") {
		t.Error("a binary file was searched")
	}
}

// Test24R2ASearchContentNeverListsBlacklistedFiles is the AC "blacklist file
// never appears even in listings" (ticket:25-26): exclusion happens before the
// file is opened, and the name is never printed.
func Test24R2ASearchContentNeverListsBlacklistedFiles(t *testing.T) {
	root := r2aCorpus(t)
	d := SearchDeps{Classify: func(p string) risk.Class {
		if strings.Contains(p, "blacklist.txt") {
			return risk.ClassA
		}
		return risk.ClassNone
	}}
	out, err := r2aScan(d, root, 10).run(context.Background(), root)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if out.excluded != 1 {
		t.Errorf("excluded = %d, want 1", out.excluded)
	}
	if len(out.matches) != 3 {
		t.Errorf("matches = %d, want 3 with one file excluded", len(out.matches))
	}
	for _, m := range out.matches {
		if strings.Contains(m.path, "blacklist") {
			t.Errorf("a blacklisted path appears in the results: %q", m.path)
		}
		if strings.Contains(m.snippet, "blacklisted") {
			t.Errorf("a blacklisted file leaked its text: %q", m.snippet)
		}
	}
}

func Test24R2ASearchContentCapsAndSkips(t *testing.T) {
	root := r2aCorpus(t)
	none := func(string) risk.Class { return risk.ClassNone }

	limited, err := r2aScan(SearchDeps{Classify: none}, root, 1).run(context.Background(), root)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(limited.matches) != 1 || !limited.hitResultCap {
		t.Errorf("a limit of 1 must return one match and set hitResultCap: %+v", limited)
	}

	fileCapped, err := r2aScan(SearchDeps{MaxFiles: 1, Classify: none}, root, 10).run(context.Background(), root)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !fileCapped.hitFileCap {
		t.Errorf("MaxFiles=1 must stop the scan with hitFileCap: %+v", fileCapped)
	}

	big, err := r2aScan(SearchDeps{MaxFileBytes: 4, Classify: none}, root, 10).run(context.Background(), root)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(big.matches) != 0 {
		t.Errorf("every fixture file is over 4 bytes, want zero matches, got %+v", big.matches)
	}
}

// Test24R2ASearchContentCancelMarksPartialResult is ticket:27 - a scan cut short
// reports what it found and says truncated, it never pretends to be complete.
func Test24R2ASearchContentCancelMarksPartialResult(t *testing.T) {
	root := r2aCorpus(t)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	d := SearchDeps{Classify: func(string) risk.Class {
		cancel()
		return risk.ClassNone
	}}
	out, err := r2aScan(d, root, 10).run(ctx, root)
	if err != nil {
		t.Fatalf("a cancelled scan must not fail the call: %v", err)
	}
	if out.aborted == "" {
		t.Fatal("aborted is empty: a cancelled scan was reported as a clean one")
	}
}

func Test24R2ASearchContentRefusesOutsideAllowlist(t *testing.T) {
	inside := t.TempDir()
	outside := r2aCorpus(t)
	p := NewPathCanonicalizer([]string{inside}, nil)
	if p == nil {
		t.Fatal("nil canonicalizer")
	}
	res, err := searchContent{d: SearchDeps{Paths: p}}.Execute(context.Background(),
		r2aParams(t, map[string]any{"path": outside, "query": "needle"}), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.IsError || !strings.Contains(res.Text, "fail-closed") {
		t.Errorf("a root outside the allowlist must be refused fail-closed, got %+v", res)
	}
	if strings.Contains(res.Text, "SECRET") {
		t.Error("the refused scan leaked result lines")
	}
}

func Test24R2ASearchContentFailClosedWithoutResolver(t *testing.T) {
	res, err := searchContent{d: SearchDeps{}}.Execute(context.Background(),
		r2aParams(t, map[string]any{"path": `C:\x`, "query": "needle"}), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.IsError || !strings.Contains(res.Text, "fail-closed") {
		t.Errorf("no C26 resolver must fail closed, got %+v", res)
	}
	missingQuery, err := searchContent{d: SearchDeps{}}.Execute(context.Background(),
		r2aParams(t, map[string]any{"path": `C:\x`, "query": ""}), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !missingQuery.IsError {
		t.Errorf("an empty query must not turn into a full listing, got %+v", missingQuery)
	}
}

func Test24R2ASearchContentThroughResolverReportsTruncated(t *testing.T) {
	root := r2aCorpus(t)
	long := filepath.Join(root, "long.txt")
	if err := os.WriteFile(long, []byte(strings.Repeat("needle ", 100)+"needle\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	p, canon := r2aCanon(t, []string{root}, root)
	res, err := searchContent{d: SearchDeps{Paths: p, MaxResults: 2}}.Execute(context.Background(),
		r2aParams(t, map[string]any{"path": canon, "query": "needle"}), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if res.IsError {
		t.Fatalf("scan failed: %s", res.Text)
	}
	if !res.Truncated {
		t.Errorf("MaxResults=2 must set Truncated (the spill layer reads that): %q", res.Text)
	}
	if res.Origin != canon {
		t.Errorf("Origin = %q, want the canonical root %q (the C25 taint seam)", res.Origin, canon)
	}
	if !strings.Contains(res.Text, "truncated") {
		t.Errorf("the answer must say truncated in its header: %q", res.Text)
	}
	for _, line := range strings.Split(res.Text, "\n")[1:] {
		if line == "" {
			continue
		}
		i := strings.Index(line, ":")
		j := strings.Index(line[i+1:], ":")
		if i < 0 || j < 0 {
			t.Errorf("a hit must carry path and line: %q", line)
			continue
		}
		snippet := line[i+j+2:]
		if len(snippet) > defaultSearchSnippet+3 {
			t.Errorf("snippet not clipped (%d bytes): %q", len(snippet), snippet)
		}
	}
}

func Test24R2AClipsSnippetWithoutSplittingRunes(t *testing.T) {
	got := clipSnippet(strings.Repeat("汉", 200)) // 600 bytes, over the 200-byte budget
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("clipSnippet must mark the clip, got %q", got)
	}
	body := strings.TrimSuffix(got, "...")
	if len(body) > defaultSearchSnippet {
		t.Errorf("clip is %d bytes, over the %d budget", len(body), defaultSearchSnippet)
	}
	for _, r := range body {
		if r == 0xFFFD {
			t.Errorf("the clip split a UTF-8 rune: %q", body)
			break
		}
	}
	if len(body)%len("汉") != 0 {
		t.Errorf("clip is not a whole number of runes (%d bytes)", len(body))
	}
}

// ---------------------------------------------------------------------------
// registration shape (C1 name, C3 capability, R1 floor, C26 path params)
// ---------------------------------------------------------------------------

func Test24R2AEntriesRegisterWithD34Shape(t *testing.T) {
	r := NewRegistry()
	dir := t.TempDir()
	p := NewPathCanonicalizer([]string{dir}, nil)
	entries := append(
		BuiltinDocumentEntries(DocDeps{Paths: p}),
		BuiltinSearchContentEntries(SearchDeps{Paths: p})...,
	)
	for _, e := range entries {
		if err := r.Register(e); err != nil {
			t.Fatalf("register %s: %v", e.Tool.Name(), err)
		}
	}
	for name := range map[string]bool{"doc.read": true, "search.content": true} {
		got, ok := r.Lookup(name)
		if !ok {
			t.Fatalf("%s not registered", name)
		}
		if got.Decl.Declared != risk.L0 {
			t.Errorf("%s declared risk = %v, want L0 (D34 row)", name, got.Decl.Declared)
		}
		if len(got.Decl.Capabilities) != 1 || got.Decl.Capabilities[0] != CapFSRead {
			t.Errorf("%s capabilities = %v, want fs.read (D34 row)", name, got.Decl.Capabilities)
		}
		if len(got.Decl.PathParams) != 1 || got.Decl.PathParams[0] != "path" {
			t.Errorf("%s PathParams = %v, want [path] so C19 judges it through C26", name, got.Decl.PathParams)
		}
		if got.Decl.Provider != KindBuiltin || !got.Decl.Resident {
			t.Errorf("%s decl = %+v, want builtin + resident", name, got.Decl)
		}
		if len(got.Tool.Parameters()) == 0 {
			t.Errorf("%s has no parameter schema", name)
		}
	}
}
