package tools

import (
	"archive/zip"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/ledongthuc/pdf"
)

// Ticket 24 (D34 `doc.read` row, SPEC-07 §3 "PDF/docx/pptx -> 文本", capability
// fs.read, L0 inside the allowlist, results taint-marked via Result.Origin and
// the bridge's C25 mark).
//
// The docx reader is the zero-dependency half of the ticket: OOXML is a zip
// container whose body is `word/document.xml`, so stdlib `archive/zip` plus
// `encoding/xml` is enough to pull the text out (ticket:16-17 "docx XML text
// extraction with pure-Go"). search.content lives in search_content.go.
//
// The PDF reader (24-r5) rides the whitelisted pure-Go parser
// `github.com/ledongthuc/pdf` (BSD-3) rather than an invented one: ToUnicode
// CMaps, embedded font encodings and the xref table are exactly what the ticket
// forbids reimplementing here. Page text is walked per page and concatenated
// under the caller's byte budget; the spill layer (D15(3)) owns anything above
// the context budget, so this reader only ever reports Truncated.
//
// What is NOT here, and why (every one is a registered deferral, never a silent
// success):
//
//	DEFERRED(xlsx, D34)  SPEC-12 §5 row "`doc.read` xlsx/OCR（D34）". No
//	   extractor is registered for .xlsx, so doc.read refuses it rather than
//	   returning a blob of shared-strings XML dressed up as a document.
//	DEFERRED(OCR, D34)   same §5 row. A no-text-layer PDF (a scan) must produce
//	   the clear error pinned in scannedPDFReason below, never an empty result
//	   (ticket:20-21). The wording lives here so the PDF limb inherits it with
//	   the test that hangs on it.

// ErrUnsupported is the seam's "this format has no extractor here" signal. It
// is errors.Is-checkable so a caller can tell "not implemented" apart from
// "this file is broken" - the second must never be laundered into the first.
var ErrUnsupported = errors.New("tools: no text extractor for this document format")

// Extracted is one document's text plus the honest flags.
type Extracted struct {
	// Text is the extracted body (UTF-8, one '\n' per paragraph).
	Text string
	// Truncated says the extraction stopped at the caller's byte budget while
	// more text was available.
	Truncated bool
	// Runs counts the text runs it walked. Zero runs with no error means the
	// document really has no extractable text, which doc.read treats as an
	// error case, not an empty success.
	Runs int
}

// TextExtractor is the seam doc.read dispatches on: one implementation per
// document format, reached only through a C26-canonical path.
type TextExtractor interface {
	// Format is the lowercase format name ("docx", "pdf").
	Format() string
	// ExtractText reads canonicalPath (already C26-canonical, already gated)
	// and returns at most limit bytes of text. limit <= 0 means the caller
	// wants the deps' default.
	ExtractText(ctx context.Context, canonicalPath string, limit int64) (Extracted, error)
}

// ExtractorSet maps a file extension to the extractor that owns it. An
// extension with no entry is refused; an extension whose entry is a deferral
// stub is refused with that stub's own reason.
type ExtractorSet struct {
	byExt map[string]TextExtractor
	exts  []string
}

// NewExtractorSet builds the table. Two extractors for one format is a wiring
// bug and is refused here, so a composition root cannot silently shadow one
// format's reader with another.
func NewExtractorSet(extractors ...TextExtractor) (*ExtractorSet, error) {
	s := &ExtractorSet{byExt: map[string]TextExtractor{}}
	for _, e := range extractors {
		if e == nil {
			return nil, errors.New("tools: extractor set: nil extractor")
		}
		ext := strings.ToLower(strings.TrimSpace(e.Format()))
		if ext == "" {
			return nil, errors.New("tools: extractor set: empty format name")
		}
		if _, dup := s.byExt[ext]; dup {
			return nil, fmt.Errorf("tools: extractor set: duplicate extractor for %q", ext)
		}
		s.byExt[ext] = e
		s.exts = append(s.exts, ext)
	}
	return s, nil
}

// Lookup resolves one canonical path to its extractor by extension.
func (s *ExtractorSet) Lookup(canonicalPath string) (TextExtractor, string, bool) {
	if s == nil || len(s.byExt) == 0 {
		return nil, "", false
	}
	ext := strings.ToLower(strings.TrimPrefix(extOfPath(canonicalPath), "."))
	if ext == "" {
		return nil, "", false
	}
	e, ok := s.byExt[ext]
	return e, ext, ok
}

// Formats lists the registered formats (for descriptions and errors).
func (s *ExtractorSet) Formats() []string {
	if s == nil {
		return nil
	}
	out := append([]string(nil), s.exts...)
	sort.Strings(out)
	return out
}

// extOfPath returns the trailing ".ext" of a canonical path. It splits on the
// separators C26 emits and never cleans or absolutizes anything: path
// decisions belong to the resolver, not here (D22).
func extOfPath(p string) string {
	base := p
	if i := strings.LastIndexAny(p, `/\`); i >= 0 {
		base = p[i+1:]
	}
	if i := strings.LastIndexByte(base, '.'); i > 0 {
		return base[i:]
	}
	return ""
}

// ---------------------------------------------------------------------------
// docx: the zero-dependency reader (archive/zip + encoding/xml)
// ---------------------------------------------------------------------------

// docxBodyPart is the single zip entry that holds a WordprocessingML body. The
// name is matched exactly, which is also why there is no zip-slip surface here:
// nothing is ever written to disk from the archive.
const docxBodyPart = "word/document.xml"

// docxMainNS is the WordprocessingML main namespace as a prefix. Go's decoder
// does not resolve xmlns URIs, so a prefixed document reports Space "w"; a
// default-namespaced document reports "". Both are accepted, anything else is
// somebody else's markup and is skipped.
const docxMainNS = "w"

// docxWMLURI is the full namespace URI, accepted for writers that emit it as
// the literal prefix (rare, but harmless to allow).
const docxWMLURI = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"

type docxText struct{}

func (docxText) Format() string { return "docx" }

// ExtractText implements TextExtractor. A malformed, non-zip or body-less
// package is an error out of this method, so doc.read reports it as a parse
// failure:坏 docx 绝不静默成功.
func (docxText) ExtractText(ctx context.Context, canonicalPath string, limit int64) (Extracted, error) {
	if limit <= 0 {
		limit = defaultMaxTextBytes
	}
	zr, err := zip.OpenReader(canonicalPath)
	if err != nil {
		return Extracted{}, fmt.Errorf("docx: 不是有效的 OOXML 包（zip 打不开）：%w", err)
	}
	defer zr.Close() //nolint:errcheck // read-only archive handle

	var body *zip.File
	for _, f := range zr.File {
		if f.Name == docxBodyPart {
			body = f
			break
		}
	}
	if body == nil {
		return Extracted{}, fmt.Errorf("docx: 包里缺少 %s（不是 Word 正文包）", docxBodyPart)
	}
	rc, err := body.Open()
	if err != nil {
		return Extracted{}, fmt.Errorf("docx: 打开 %s 失败：%w", docxBodyPart, err)
	}
	defer rc.Close() //nolint:errcheck // read-only entry handle

	// Decompression budget. The caller's limit is on extracted TEXT; a zip
	// that inflates far past it still must not be allowed to expand without
	// bound, so the inflating stream is capped at a multiple of the budget and
	// hitting that ceiling is reported as Truncated, never as a clean read.
	budget := limit*docxInflateFactor + docxInflateSlack
	cr := &countingReader{r: io.LimitReader(rc, budget)}
	got, err := extractDocxXML(ctx, cr, limit)
	if err != nil {
		return got, err
	}
	if cr.n >= budget && !got.Truncated {
		// The stream ended exactly at the inflate ceiling: the body may
		// continue past it, so the text cannot be claimed complete.
		got.Truncated = true
	}
	return got, nil
}

const (
	docxInflateFactor = 8
	docxInflateSlack  = 4096
)

// countingReader records how many bytes were actually consumed.
type countingReader struct {
	r io.Reader
	n int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n += int64(n)
	return n, err
}

// extractDocxXML walks the WordprocessingML token stream and keeps the text.
// Paragraph ends become newlines, `w:br` a newline and `w:tab` a tab; anything
// that is not a text run is ignored (formatting, relationships, numbering
// references). A decode error is returned, never swallowed.
func extractDocxXML(ctx context.Context, r io.Reader, limit int64) (Extracted, error) {
	dec := xml.NewDecoder(r)
	var (
		b         strings.Builder
		runs      int
		truncated bool
		inText    bool
	)
walk:
	for {
		if err := ctx.Err(); err != nil {
			return Extracted{}, err
		}
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return Extracted{}, fmt.Errorf("docx: %s 解析失败：%w", docxBodyPart, err)
		}
		switch tv := tok.(type) {
		case xml.StartElement:
			if !docxTokenNS(tv.Name) {
				continue
			}
			switch tv.Name.Local {
			case "t":
				inText = true
			case "tab":
				inText = false
				if !appendText(&b, "\t", limit) {
					truncated = true
					break walk
				}
			case "br":
				inText = false
				if !appendText(&b, "\n", limit) {
					truncated = true
					break walk
				}
			}
		case xml.EndElement:
			if !docxTokenNS(tv.Name) {
				continue
			}
			switch tv.Name.Local {
			case "t":
				inText = false
			case "p":
				if !appendText(&b, "\n", limit) {
					truncated = true
					break walk
				}
			}
		case xml.CharData:
			if !inText {
				continue
			}
			if len(tv) == 0 {
				continue
			}
			runs++
			if !appendText(&b, string(tv), limit) {
				truncated = true
				break walk
			}
		}
	}
	return Extracted{Text: b.String(), Truncated: truncated, Runs: runs}, nil
}

// docxTokenNS accepts the main-wordprocessingml prefix (or an unprefixed
// element) and rejects foreign namespaces.
func docxTokenNS(n xml.Name) bool {
	return n.Space == "" || n.Space == docxMainNS || n.Space == docxWMLURI
}

// appendText writes s unless the budget is already spent; false means the
// caller must stop and report the result as truncated.
func appendText(b *strings.Builder, s string, limit int64) bool {
	if int64(b.Len())+int64(len(s)) > limit {
		return false
	}
	b.WriteString(s)
	return true
}

// ---------------------------------------------------------------------------
// PDF text extractor: ledongthuc/pdf (BSD-3 pure Go, whitelisted in SPEC-01 §96)
// ---------------------------------------------------------------------------

// pdfText extracts text from PDFs using the ledongthuc library. It handles
// multi-page concatenation, scanned PDFs (no text layer), and truncation.
type pdfText struct{}

func (pdfText) Format() string { return "pdf" }

// ExtractText implements TextExtractor. It walks pages, extracts plain text with
// Unicode support, and reports scanned or broken PDFs as errors. Multi-page
// extraction respects the caller's byte limit (C15 context budget spill layer).
func (pdfText) ExtractText(ctx context.Context, canonicalPath string, limit int64) (Extracted, error) {
	if limit <= 0 {
		limit = defaultMaxTextBytes
	}
	// Fail closed: a bad PDF is an error, never a silent empty success.
	f, err := os.Open(canonicalPath)
	if err != nil {
		return Extracted{}, fmt.Errorf("pdf: 无法打开文件：%w", err)
	}
	defer f.Close() //nolint:errcheck // read-only

	fi, err := f.Stat()
	if err != nil {
		return Extracted{}, fmt.Errorf("pdf: 无法读取文件信息：%w", err)
	}
	size := fi.Size()

	reader, err := pdf.NewReader(f, size)
	if err != nil {
		return Extracted{}, fmt.Errorf("pdf: 不是有效的 PDF（解析失败）：%w", err)
	}

	numPages := reader.NumPage()
	if numPages == 0 {
		return Extracted{}, fmt.Errorf("pdf: %s", scannedPDFReason)
	}

	var (
		b         strings.Builder
		runs      int
		truncated bool
		pageErr   error
	)
	for i := 1; i <= numPages; i++ {
		if err := ctx.Err(); err != nil {
			return Extracted{}, err
		}
		page := reader.Page(i)
		if page.V.Kind() == pdf.Null {
			// A missing page index is not an error; it is just no content here.
			continue
		}
		// ledongthuc reports a page-level panic as an error here (its GetPlainText
		// recovers internally). One bad content stream must not void a document
		// whose other pages read fine - but it is remembered, so a document that
		// yields nothing is never guessed at as a scan when it actually broke.
		text, err := page.GetPlainText(nil)
		if err != nil {
			if pageErr == nil {
				pageErr = err
			}
			continue
		}
		// Page text from ledongthuc carries layout newlines (one per BT block);
		// the body is what matters here, so the shell is trimmed and pages are
		// joined with a single newline.
		body := strings.TrimSpace(text)
		if body == "" {
			continue
		}
		if b.Len() > 0 {
			if !appendText(&b, "\n", limit) {
				truncated = true
				break
			}
		}
		runs++
		if !appendTextWithin(&b, body, limit) {
			truncated = true
			break
		}
	}
	// Scanned PDF (or every page empty): zero text and no content-stream failure.
	// This is the ticket's "no silent empty result" case, and the marker says what
	// it observed instead of guessing at why (OCR is DEFERRED(OCR, D34)).
	if runs == 0 {
		if pageErr != nil {
			return Extracted{}, fmt.Errorf("pdf: 内容流解析失败：%w", pageErr)
		}
		return Extracted{}, fmt.Errorf("pdf: %s（疑似扫描件，无文本层）", scannedPDFReason)
	}
	return Extracted{Text: b.String(), Truncated: truncated, Runs: runs}, nil
}

// appendTextWithin writes as much of s as the remaining budget holds and returns
// false when s did not fit whole. Unlike appendText it never throws a run away:
// a truncated read still owes the caller the head of the text, cut on a rune
// boundary so the result stays valid UTF-8.
func appendTextWithin(b *strings.Builder, s string, limit int64) bool {
	room := limit - int64(b.Len())
	if room <= 0 {
		return false
	}
	if int64(len(s)) <= room {
		b.WriteString(s)
		return true
	}
	cut := s[:room]
	for !utf8.ValidString(cut) && len(cut) > 0 {
		cut = cut[:len(cut)-1]
	}
	b.WriteString(cut)
	return false
}

// scannedPDFReason is what a PDF with no text layer owes its reader (ticket:
// 20-21, "no silent empty result"). It is a const, not a dead string, so the
// limb that lands the parser and this limb's test both pin the same wording.
const scannedPDFReason = "无可提取文本（扫描件暂不支持 OCR）"

// ---------------------------------------------------------------------------
// doc.read (C1 Tool, D34 `doc.read` row)
// ---------------------------------------------------------------------------

// DocDeps is what doc.read needs from the host: the same C26 resolver every fs
// tool uses, the extractor table, one byte budget and the A/B classifier.
type DocDeps struct {
	// Paths is the C26 canonicalizer. nil fails closed, exactly as FSDeps does.
	Paths *PathCanonicalizer
	// Extractors is the format table. nil selects DefaultExtractors().
	Extractors *ExtractorSet
	// MaxTextBytes caps one extraction (default defaultMaxTextBytes). The
	// loop's spill layer (D15(3)) owns what happens above the context budget,
	// so this tool reports Truncated and does not invent a spill of its own.
	MaxTextBytes int
	// Classify is the A/B blacklist seam; nil means risk.Classify. doc.read and
	// search.content both refuse a classified path BEFORE reading a byte, so a
	// blacklisted document never leaks even its text (ticket:25-26).
	Classify func(string) risk.Class
}

// defaultMaxTextBytes is the doc.read byte budget when the host sets none.
const defaultMaxTextBytes = 256 * 1024

func (d DocDeps) textCap() int64 {
	if d.MaxTextBytes > 0 {
		return int64(d.MaxTextBytes)
	}
	return defaultMaxTextBytes
}

func (d DocDeps) set() *ExtractorSet {
	if d.Extractors != nil {
		return d.Extractors
	}
	s, err := NewExtractorSet(DefaultExtractors()...)
	if err != nil {
		return nil
	}
	return s
}

func (d DocDeps) classify(canonical string) risk.Class {
	if d.Classify != nil {
		return d.Classify(canonical)
	}
	return risk.Classify(canonical)
}

// DefaultExtractors is the D34 doc.read roster as far as it lands today: docx
// and PDF both real (24-r5). xlsx is absent on purpose - DEFERRED(xlsx, D34),
// as is OCR for scanned PDFs, which pdfText refuses with scannedPDFReason.
func DefaultExtractors() []TextExtractor {
	return []TextExtractor{docxText{}, pdfText{}}
}

// resolveDoc canonicalizes a model-supplied path through C26, with the same
// fail-closed rules as FSDeps.open: no resolver, an unresolvable spelling, or
// an unresolved component all refuse the call.
func (d DocDeps) resolveDoc(raw string) (string, error) {
	if d.Paths == nil {
		return "", errors.New("doc: no C26 resolver configured (fail-closed)")
	}
	c, err := d.Paths.Canonicalize(raw)
	if err != nil {
		return "", err
	}
	if strings.Contains(c, "??") {
		return "", errors.New("doc: unresolved path component")
	}
	return c, nil
}

type docRead struct{ d DocDeps }

type docReadArgs struct {
	Path    string `json:"path"`
	MaxChar int    `json:"max_chars"`
}

var docReadSchema = JSONSchema(`{"type":"object","properties":{` +
	`"path":{"type":"string","description":"文档路径（经 C26 规范化）"},` +
	`"max_chars":{"type":"integer","minimum":1,"description":"本次最多返回多少字节文本"}}` +
	`,"required":["path"],"additionalProperties":false}`)

func (docRead) Name() string { return "doc.read" }

// Description implements Tool. It names the deferred formats up front instead of
// letting the model discover them by trial.
func (docRead) Description() string {
	return "提取文档正文文本（docx 与 PDF 已支持；PDF 无文本层（扫描件）与 xlsx 推迟，OCR 见 DEFERRED）"
}

func (docRead) Parameters() JSONSchema { return docReadSchema }

// Execute implements C1 Tool.
func (t docRead) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a docReadArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	if strings.TrimSpace(a.Path) == "" {
		return Result{Text: "缺少 path 参数", IsError: true}, nil
	}
	canon, err := t.d.resolveDoc(a.Path)
	if err != nil {
		return Result{Text: "路径无法解析（按 fail-closed 拒绝）：" + err.Error(), IsError: true}, nil
	}
	// The blacklist check runs before any open: a classified document is not
	// read and then filtered, it is not read at all.
	if cls := t.d.classify(canon); cls != risk.ClassNone {
		return Result{Text: "路径在敏感清单内（" + cls.String() + " 级），按 fail-closed 不读取", IsError: true}, nil
	}
	ex, ext, ok := t.d.set().Lookup(canon)
	if !ok {
		return Result{Text: fmt.Sprintf(
			"不支持的文档格式 %s（已注册：%s；xlsx/OCR 属 DEFERRED，登记见 SPEC-12 §5）",
			ext, strings.Join(t.d.set().Formats(), "/")), IsError: true, Origin: canon}, nil
	}
	limit := t.d.textCap()
	if a.MaxChar > 0 && int64(a.MaxChar) < limit {
		limit = int64(a.MaxChar)
	}
	got, err := ex.ExtractText(ctx, canon, limit)
	if err != nil {
		if errors.Is(err, ErrUnsupported) {
			return Result{Text: "该格式暂未实现：" + err.Error(), IsError: true, Origin: canon}, nil
		}
		return Result{Text: "解析失败：" + err.Error(), IsError: true, Origin: canon}, nil
	}
	// "no silent empty result" (ticket:20-21): zero extracted text is an error.
	if strings.TrimSpace(got.Text) == "" {
		return Result{Text: fmt.Sprintf("%s：%s 未提取到任何文本（空文档或图片型文档）",
			scannedPDFReason, ext), IsError: true, Origin: canon}, nil
	}
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("doc.read 提取 %d 字节（%d 个文本段）", len(got.Text), got.Runs))
	}
	return Result{
		Text:      got.Text,
		Truncated: got.Truncated,
		Origin:    canon,
	}, nil
}

// DocReadDecl is the host-side declaration for doc.read: capability fs.read
// (D34 row), R1 floor L0 (C19 still owns the verdict), path judged via C26.
func DocReadDecl() Decl {
	return Decl{
		Capabilities: []Capability{CapFSRead},
		Needs:        []Capability{CapFSRead},
		Declared:     risk.L0,
		PathParams:   []string{"path"},
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// BuiltinDocumentEntries returns the D34 document-read entry for the host's
// Builtin slot. Wiring it into the bridge is a shared file (bridge.go), outside
// this limb's write surface (ticket 24 / 24-r2a: new document and search files
// only) - see docs/evidence/s1/24r2a-impl.md §4.
func BuiltinDocumentEntries(d DocDeps) []Entry {
	return []Entry{{Tool: docRead{d: d}, Decl: DocReadDecl()}}
}
