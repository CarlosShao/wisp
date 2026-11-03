# 24r4-stop: PDF text extraction API not available in pdfcpu v0.16.4

**Ticket:** 24 (`doc.read`) — leg r4 (PDF text layer implementation)  
**Date:** 2026-11-03 10:5x +08  
**Status:** STOPPED (hard gate: no suitable API in whitelisted dependency)

## Gate rule

Per task directive: if pdfcpu text extraction API is unusable/substandard ⇒ stop and write this file; do NOT implement your own parser ("invented parser path" forbidden by ticket:16-17).

## Evidence

### A. Build block (repo-level dependency graph)

```bash
$ go build github.com/pdfcpu/pdfcpu/pkg/api
D:\work\base\gopath\pkg\mod\github.com\pdfcpu\pdfcpu@v0.16.4\pkg\filter\ccittDecode.go:25:2: missing go.sum entry for module providing package golang.org/x/image/ccitt; to add:
	go get github.com/pdfcpu/pdfcpu/pkg/filter@v0.16.4
D:\work\base\gopath\pkg\mod\github.com\pdfcpu\pdfcpu@v0.16.4\pkg\pdfcpu\types\string.go:28:2: missing go.sum entry for module providing package golang.org/x/text/unicode/norm; to add:
	go get github.com/pdfcpu/pdfcpu/pkg/pdfcpu@v0.16.4
D:\work\base\gopath\pkg\mod\github.com\pdfcpu\pdfcpu@v0.16.4\pkg\pdfcpu\writeImage.go:31:2: missing go.sum entry for module providing package github.com/hhrutter/tiff; to add:
	go get github.com/pdfcpu/pdfcpu/pkg/pdfcpu@v0.16.4
... (more missing transitive deps)
```

`go.mod` line shows `// indirect`; `go.sum` contains zero entries for the six transitive modules that pdfcpu actually requires:

- `golang.org/x/image/{ccitt,webp}`
- `golang.org/x/text/{unicode/norm,encoding/charmap}`
- `go.yaml.in/yaml/v3`
- `github.com/hhrutter/tiff`
- `github.com/xdg-go/stringprep`
- `github.com/mattn/go-runewidth`

Fixing this would require `go mod tidy` or `go get -t`, which modifies `go.mod` and `go.sum`—outside the declared write surface for leg r4 (only `internal/tools/document_text.go` + companion test file permitted). Without go.sum repair, any change importing `api` cannot compile.

### B. Capability gap (pdfcpu v0.16.4 lacks text extraction API)

#### B1 ExtractContent returns raw content streams

Signature:

```go
func ExtractContent(c context.Context, rs io.ReadSeeker, selectedPages []string, digestContent func(io.Reader, int) error, conf *model.Configuration) (err error)
func ExtractPageContent(c context.Context, ctx *model.Context, pageNr int) (io.Reader, error)
```

Impl of `ExtractPageContent` (extract.go:866-890):

```go
bb, err := ctx.PageContent(d, pageNr)   // inflates filters, returns raw operator bytes
return bytes.NewReader(bb), nil
```

The returned reader yields literal PDF content-stream bytes, e.g.:

```
BT /F1 12 Tf 100 500 Td (Hello) Tj ET q Q cm 10 0 0 10 0 0 cm
```

This is **not decoded Unicode text**. The string `(Hello)` may be in a non-Unicode character encoding (e.g., WinAnsi, MacRoman, custom CIDKey maps); ToUnicode CMaps exist but are **not applied** by these APIs.

#### B2 No helper to decode content to Unicode

Search of v0.16.4 sources for exported text decoders found none:

```bash
grep -r "func .*(Text|String).*\(.*Content" pkg/ | grep -v "_test.go" # none
```

ToUnicode appears only in font-related paths (`font/install.go`, `font/metrics.go`, `font/fontDict.go`): used for watermarks & metrics, **never** for turning `Tj/TJ` operators into Unicode.

Thus, proper text extraction requires implementing:

- Content-stream operator tokenizer (scan for `BT..ET`, handle `q/Q`, `cm`, etc.)
- Resource dictionary traversal (Font resources, ToUnicode CMaps as XObjects)
- Encoding detection (WinAnsi/UTF16BE/BMCID/CMap lookup)
- Hex-string vs literal-string handling + CMap decoding (RFC 3294)
- Multi-byte CJK via CIDKey maps

That is precisely the "reimplement ToUnicode CMaps + embedded font encodings + xref cross table" the stub comment flags as "invent your own parser path," which the ticket explicitly forbids.

#### B3 Upstream status confirmation (citations from discussion threads)

External discussions confirm that **text extraction with CMap support has never been merged** into pdfcpu v0.x:

- Discussion #580 "text extraction, page rendering": ongoing feature request thread
- Issue #122 "Text extraction": long-standing open issue (no implementation in current major version)

(Exact URLs not fetched due to network restrictions; see search results in session logs.)

## Hard-gate decision

Both conditions fail:

1. ❌ API does not perform Unicode text extraction from content streams
2. ❌ Dependency chain not built-in-repo; building would break scope boundaries

Per task instructions: **"停手（写 24r4-stop.md 说明，⛔硬造）"**

## Consequences

The PDF branch in `internal/tools/document_text.go` must remain a stub. All items referencing PDF text extraction stay marked:

- `DEFERRED(D34-1)` remains in code comments
- Registration entry stays absent (xlsx/OCR already deferred)

Acceptance criteria shift: `doc.read` supports `"docx"` only until a future revision where a human-approved parser dependency provides text extraction. Scanned-PDFs return the fixed marker:

```text
无可提取文本（扫描件暂不支持 OCR）
```

## Required follow-ups (out of scope for leg r4)

- Human approval to either:
  - Add a different text-extraction library to SPEC-01 whitelist, or
  - Accept that PDF text extraction is out-of-scope for this phase
- If approved, coordinate `go mod tidy` to repair transitive dependencies

---

*Leg r4 delivered on time: probe confirmed stop condition. Next step awaits human decision on pivot path or scope adjustment.*
