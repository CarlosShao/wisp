package panel

// Artifact preview - ticket 192 AC#1 (the preview stack, split out of the
// browser tab) and AC#2 (the injection face: preview content is not input).
//
// WHAT THIS FILE OWNS. One read-only question, asked by the WebView:
// "give me the bytes of artifact <bare name> so I can render them". It owns
// the classification (which four classes are previewable), the byte-vs-path
// decision (bytes, never a local path), the truncation rule for text (a
// window plus a pointer that walks the whole file) and the one door through
// which preview TEXT may become model context, with the C25 mark on it.
//
// WHAT IT DELIBERATELY DOES NOT OWN. No route in the C17 method whitelist
// (ticket AC#4: "C17 不加方法"), no network, no form submit, no download, no
// script execution (AC#3), and no new contract text. The preview is served
// the way the embedded bundle already is: as a response to one intercepted
// resource request under the panel origin, because that is the existing,
// already-gated byte pipe (cmd/wisp/panel_control_windows.go's
// AddWebResourceRequestedFilter arm). Minting a postMessage method for it
// would put a second door on the shared envelope that the panel's own
// composer routes travel, and the only thing the page may do with a preview
// is DRAW it - never send it. AC#2 is the reason that distinction has to be
// structural rather than a promise.
//
// WHY BYTES AND NOT A PATH. The caller is the least trusted part of the
// process (SPEC-06 §1, ban #6). A preview addressed by a filesystem path
// would be a "read any path the page names" capability, which is exactly what
// C26's authorization root and the ticket 175/177/183 exfil family exist to
// stop; ticket 192's own text says "不许把'能预览'做成'模型能读任意路径'".
// So the request carries a BARE ARTIFACT NAME, that name is judged by the
// artifacts store's own guard (memory.ValidArtifactName, injected as
// NameGuard exactly the way attachments.go does it), the file must be a
// regular non-link file inside the store's own artifacts directory, and what
// leaves this function is bytes plus a same-origin URL. No filesystem path is
// ever built from the caller's string beyond joining a name the guard has
// already accepted, and no path is ever echoed into a body or a header: the
// browser never learns where the bytes live (pinned by the criterion that
// asserts the artifacts dir does not appear in the response).
//
// WHY THE HTML CLASS IS SANDBOXED. Previewing host-generated HTML is
// rendering outside content inside the panel window, and that window's origin
// speaks C17: bridge.ParseComposerRequest trusts the envelope's "source"
// field, which is a routing claim, not authority (bridge.go:11-18). A page
// that could run script from a preview would therefore be able to forge
// composer requests. The HTML response is consequently served under
// "Content-Security-Policy: sandbox; default-src 'none'; ..." - scripts do
// not run and the document gets an opaque origin even if it is top-level - so
// the preview can be looked at and cannot speak. Dropping that directive is
// the mutation the AC#2 red proof aims at.
//
// TRUNCATION IS THE D15/164 SHAPE, NOT A NEW ONE. PLAN.md:431 and ticket 164
// AC#3 already define what "too long" means in this product: show the head,
// show the tail, say how long the whole thing is, and give a pointer that can
// fetch the rest. This file reuses that shape at byte granularity for text
// only, with the pointer being a window request on the same route:
//
//	first view  notice (total, head shown, tail shown, elided bytes, resume URL)
//	            + head + "middle N bytes not shown" + tail
//	window      exactly the byte range [offset, offset+limit), cut at a UTF-8
//	            boundary, no envelope at all - so a walker can concatenate
//	            windows and get the file back byte for byte (that is the
//	            criterion, not a comment).
//
// Binary classes are NEVER truncated: a half PNG is a broken image and a half
// PDF is a broken document shown as if it were whole. Over the class cap they
// are refused with both numbers, which is the loud "no" the ticket asks for.
//
// The four answers AC#1 asks for therefore live here as code, not as prose:
// which classes (classifyPreviewArtifact), path or bytes (the guard + the
// artifacts directory + the no-path-in-response criterion), how truncation
// works (textFirstView/textWindow), and what the user sees (the notice lines
// below, pinned by a criterion that parses them).

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

// PreviewPathPrefix is the one URL subspace under the panel origin this file
// answers. It is deliberately outside everything the embedded bundle owns, so
// a preview can never shadow a built asset, and it is a PATH rather than a
// postMessage method so the preview rides the existing resource filter
// (AC#4: no new C17 method).
const PreviewPathPrefix = "/__artifact/"

// The notice lines the truncation view is built from. They are constants
// because the criterion parses them: a notice the page cannot read is the
// "silently truncated" defect ticket 164 AC#3 refuses.
const (
	previewNoticeTotal  = "[wisp 产物预览]"
	previewNoticeResume = "[wisp 续读]"
	previewNoticeBody   = "[wisp 正文]"
	previewNoticeElided = "[wisp 省略]"
	previewHTMLPolicy   = "sandbox; default-src 'none'; img-src data:; style-src 'unsafe-inline'"
	previewImagePolicy  = "default-src 'none'"
)

const (
	// previewHeadBytes / previewTailBytes are what a first view shows: the
	// same head-plus-tail shape D15/164 uses, sized so one screen of a log is
	// legible without a scrollbar.
	previewHeadBytes = 4096
	previewTailBytes = 1024
	// PreviewMaxTextWindow is the largest single window a resume request may
	// ask for, and the default when it asks for nothing.
	PreviewMaxTextWindow = 64 << 10
	// Class caps for the classes that are rendered whole. Text has no cap: it
	// is windowed instead. The numbers are preview budgets, not contract
	// limits - they bound what one resource response may pull into the
	// WebView's process, and nothing here is a security boundary (the security
	// boundary is the name guard and the artifacts directory).
	previewHTMLCap    = 2 << 20
	previewImageCap   = 8 << 20
	previewPDFCap     = 32 << 20
	previewSniffBytes = 512
	// previewQuoteCap bounds the bytes the quote door will scan to prove a
	// quote came from the artifact. Past it the door refuses rather than
	// stamping text it could not check.
	previewQuoteCap = 8 << 20
)

// PreviewKind names one of the four previewable classes (ticket 192 AC#1).
type PreviewKind string

const (
	PreviewKindHTML  PreviewKind = "html"
	PreviewKindImage PreviewKind = "image"
	PreviewKindPDF   PreviewKind = "pdf"
	PreviewKindText  PreviewKind = "text"
)

// PreviewRender says who turns the bytes into pixels, because AC#1 asks
// "which classes can be previewed" and the honest answer differs per class:
// HTML/text/images render in the WebView itself with no plugin, while PDF
// rendering depends on the WebView2 build carrying a usable PDF renderer
// (WebView2's PDF support is runtime- and policy-dependent; some installations
// show a blank frame even on a 200). This file therefore never claims PDF
// pixels: it declares the dependency so the page can fall back to the
// metadata card instead of drawing an empty rectangle and looking broken.
type PreviewRender string

const (
	PreviewRendersInline  PreviewRender = "inline"
	PreviewRendersHostPDF PreviewRender = "host-pdf"
)

// ErrPreviewRefused is the sentinel behind every preview refusal: the guard
// said no, the class is not previewable, the size is over the cap, the name is
// not an artifact, or the quote could not be proven. It exists so a caller (or
// a criterion) can tell "refused with a reason" apart from "the filesystem
// failed" without matching message text - the same shape ErrAttachmentRejected
// established for the composer.
var ErrPreviewRefused = errors.New("panel: artifact preview refused")

// ArtifactsSource is the storage face a preview needs and nothing more: the
// directory the store already owns. *memory.Store satisfies it as-is, so the
// preview can never point itself at a second tree (the same rule
// present.Artifacts and attachments.ArtifactSink follow).
type ArtifactsSource interface {
	ArtifactsDir() string
}

// TaintMarker is the narrow C25 face the quote door needs: Mark one sensitive
// source into one scope. *risk.Provenance satisfies it. Taking the interface
// rather than the engine keeps internal/risk out of this package's imports and
// keeps the door unable to do anything with a scope except record that this
// content came from an artifact.
type TaintMarker interface {
	Mark(scopeID, tool, origin, content string) bool
}

// PreviewSourceTool is the source name a preview quote is marked under. It is
// risk.SrcTaskOutput spelled here as a plain string so this file imports no
// internal/risk; the criterion TestPreview192SourceToolIsTheContractName pins
// it against the contract constant, so the two can never drift apart silently.
//
// Why task.output and not fs.read: the content a preview shows is what a task
// (or the host on a task's behalf) wrote into the artifacts directory, and the
// confirmation card must name the door the content actually came through.
// risk.SrcTaskOutput's own doc says it is "what a background task printed:
// outside content by origin, however domestic the door that reads it back is",
// which is exactly this door. fs.read would claim the model ran a filesystem
// read it never ran.
const PreviewSourceTool = "task.output"

// PreviewLimits are the per-resolver budgets. Zero fields take the defaults
// (a zero limit must never mean "unlimited").
type PreviewLimits struct {
	HeadBytes int64
	TailBytes int64
	MaxWindow int64
	HTMLMax   int64
	ImageMax  int64
	PDFMax    int64
	QuoteMax  int64
}

func defaultPreviewLimits() PreviewLimits {
	return PreviewLimits{
		HeadBytes: previewHeadBytes,
		TailBytes: previewTailBytes,
		MaxWindow: PreviewMaxTextWindow,
		HTMLMax:   previewHTMLCap,
		ImageMax:  previewImageCap,
		PDFMax:    previewPDFCap,
		QuoteMax:  previewQuoteCap,
	}
}

func (l PreviewLimits) normalized() PreviewLimits {
	d := defaultPreviewLimits()
	if l.HeadBytes <= 0 {
		l.HeadBytes = d.HeadBytes
	}
	if l.TailBytes <= 0 {
		l.TailBytes = d.TailBytes
	}
	if l.MaxWindow <= 0 {
		l.MaxWindow = d.MaxWindow
	}
	if l.HTMLMax <= 0 {
		l.HTMLMax = d.HTMLMax
	}
	if l.ImageMax <= 0 {
		l.ImageMax = d.ImageMax
	}
	if l.PDFMax <= 0 {
		l.PDFMax = d.PDFMax
	}
	if l.QuoteMax <= 0 {
		l.QuoteMax = d.QuoteMax
	}
	return l
}

// PreviewResolver answers preview requests for one artifacts directory. It is
// safe for concurrent use: it holds only immutable configuration and opens a
// fresh handle per request.
type PreviewResolver struct {
	dir    ArtifactsSource
	guard  NameGuard
	limits PreviewLimits
}

// NewPreviewResolver wires the resolver to the artifacts face. A nil source or
// a nil guard is a construction error, not a no-op resolver: a preview that
// quietly had no name rule would be a read-any-path door, and one with no
// directory could only lie about what it served.
func NewPreviewResolver(dir ArtifactsSource, guard NameGuard) (*PreviewResolver, error) {
	return NewPreviewResolverWithLimits(dir, guard, PreviewLimits{})
}

// NewPreviewResolverWithLimits is the same construction with explicit budgets
// (the criteria use it to reach the caps without allocating 8 MiB fixtures).
func NewPreviewResolverWithLimits(dir ArtifactsSource, guard NameGuard, limits PreviewLimits) (*PreviewResolver, error) {
	if dir == nil {
		return nil, fmt.Errorf("%w: no artifacts directory attached", ErrPreviewRefused)
	}
	if strings.TrimSpace(dir.ArtifactsDir()) == "" {
		return nil, fmt.Errorf("%w: the artifacts directory has no path", ErrPreviewRefused)
	}
	if guard == nil {
		return nil, fmt.Errorf("%w: no artifact name guard attached", ErrPreviewRefused)
	}
	return &PreviewResolver{dir: dir, guard: guard, limits: limits.normalized()}, nil
}

// PreviewDocument is the answer to one preview request. Exactly one of Body
// (a document) or the window fields (a text range) describes what came back;
// Status says whether it is a document at all.
//
// Headers is the header block the response must carry, already complete - the
// host's response builder prepends nothing of its own for a preview, because
// the CSP that keeps a previewed document from speaking is per-class and
// belongs where the class is known.
type PreviewDocument struct {
	Name        string
	Kind        PreviewKind
	Render      PreviewRender
	Status      int
	ContentType string
	Headers     []string
	Body        []byte
	// TotalBytes is the artifact's size, always reported (also on refusals
	// that were decided by size - the user must see the number).
	TotalBytes int64
	// Offset is the first byte of Body within the artifact (0 for a whole
	// document and for a first view).
	Offset int64
	// NextOffset is the resume pointer: 0 means "nothing more". A first view
	// of a truncated text sets it to the byte after the head it showed.
	NextOffset int64
	Complete   bool
	Reason     string
}

// PreviewURL is the same-origin URL of an artifact's first view. It is the
// only way this package hands out a preview address: escaping happens here so
// a name with a space, a plus or a percent sign round-trips through Lookup.
func PreviewURL(name string) string {
	return PreviewPathPrefix + url.PathEscape(name)
}

// PreviewURLAt is the same URL asking for one window.
func PreviewURLAt(name string, offset, limit int64) string {
	return PreviewURL(name) + "?offset=" + strconv.FormatInt(offset, 10) +
		"&limit=" + strconv.FormatInt(limit, 10)
}

// previewWindow is a parsed resume request.
type previewWindow struct {
	offset  int64
	limit   int64
	present bool
}

// Lookup answers one intercepted request path. requestPath is what is left of
// the URL after the panel origin (the same argument cmd/wisp's
// panelAssetLookup takes), query included.
//
// Every outcome is a PreviewDocument: a 200 with bytes, or a 4xx with Reason.
// There is no error return on purpose - a preview that failed must render as a
// sentence the user can read, not as a Go error the host might swallow into a
// blank pane.
func (r *PreviewResolver) Lookup(ctx context.Context, requestPath string) PreviewDocument {
	if err := ctx.Err(); err != nil {
		return PreviewDocument{Status: http.StatusBadRequest, Reason: "预览请求已被取消"}
	}
	name, window, err := r.parsePreviewPath(requestPath)
	if err != nil {
		return PreviewDocument{Status: http.StatusBadRequest, Reason: err.Error()}
	}
	doc := PreviewDocument{Name: name}
	full := filepath.Join(r.dir.ArtifactsDir(), name)
	info, statErr := os.Lstat(full)
	switch {
	case errors.Is(statErr, os.ErrNotExist):
		doc.Status = http.StatusNotFound
		doc.Reason = fmt.Sprintf("产物 %q 不在 artifacts 目录里（这是『没有这个产物』，不是『内容为空』）", name)
		return doc
	case statErr != nil:
		doc.Status = http.StatusInternalServerError
		doc.Reason = fmt.Sprintf("产物 %q 无法检查：%v", name, statErr)
		return doc
	case info.Mode()&os.ModeSymlink != 0:
		doc.Status = http.StatusForbidden
		doc.Reason = fmt.Sprintf("产物 %q 是一个链接；预览只读 artifacts 目录里的普通文件，不跟随链接", name)
		return doc
	case !info.Mode().IsRegular():
		doc.Status = http.StatusForbidden
		doc.Reason = fmt.Sprintf("产物 %q 不是普通文件（是目录或设备）；预览不读它", name)
		return doc
	}
	f, openErr := os.Open(full)
	if openErr != nil {
		doc.Status = http.StatusInternalServerError
		doc.Reason = fmt.Sprintf("产物 %q 打不开：%v", name, openErr)
		return doc
	}
	defer func() { _ = f.Close() }()

	size := info.Size()
	doc.TotalBytes = size
	if size == 0 {
		// An empty artifact is an empty text file, not a failure: the store
		// can legitimately hold one and "0 字节" is itself the answer.
		doc.Kind = PreviewKindText
		doc.Render = PreviewRendersInline
		doc.Status = http.StatusOK
		doc.ContentType = "text/plain; charset=utf-8"
		doc.Headers = previewHeaders(PreviewKindText, doc.ContentType)
		doc.Complete = true
		return doc
	}
	head, headErr := previewReadAt(f, 0, minInt64(previewSniffBytes, size))
	if headErr != nil {
		doc.Status = http.StatusInternalServerError
		doc.Reason = fmt.Sprintf("产物 %q 无法读取：%v", name, headErr)
		return doc
	}
	kind, contentType, refusal := classifyPreviewArtifact(head)
	if kind == "" {
		doc.Status = http.StatusUnsupportedMediaType
		doc.Reason = refusal
		return doc
	}
	doc.Kind, doc.ContentType = kind, contentType
	doc.Render = PreviewRendersInline
	doc.Headers = previewHeaders(kind, contentType)
	if kind == PreviewKindPDF {
		doc.Render = PreviewRendersHostPDF
	}
	if kind == PreviewKindText {
		if window.present {
			return r.textWindow(f, name, size, window, doc)
		}
		return r.textFirstView(f, name, size, doc)
	}
	limit := r.classCap(kind)
	if size > limit {
		doc.Status = http.StatusRequestEntityTooLarge
		doc.Headers = nil
		doc.Reason = fmt.Sprintf("产物 %q 是 %s，%d 字节，超过本类预览上限 %d 字节；二进制内容不做截断（半个文件不是文件）",
			name, kind, size, limit)
		return doc
	}
	body, readErr := previewReadAt(f, 0, size)
	if readErr != nil {
		doc.Status = http.StatusInternalServerError
		doc.Headers = nil
		doc.Reason = fmt.Sprintf("产物 %q 无法读取：%v", name, readErr)
		return doc
	}
	doc.Status = http.StatusOK
	doc.Body = body
	doc.Complete = true
	return doc
}

func (r *PreviewResolver) classCap(kind PreviewKind) int64 {
	switch kind {
	case PreviewKindHTML:
		return r.limits.HTMLMax
	case PreviewKindImage:
		return r.limits.ImageMax
	case PreviewKindPDF:
		return r.limits.PDFMax
	}
	return 0
}

// parsePreviewPath turns "/__artifact/<name>?offset=&limit=" into a name and an
// optional window request. Unknown query parameters are REFUSED rather than
// ignored: an option nobody honoured would read back as honoured.
func (r *PreviewResolver) parsePreviewPath(requestPath string) (string, previewWindow, error) {
	var window previewWindow
	raw, query := requestPath, ""
	if i := strings.IndexByte(raw, '?'); i >= 0 {
		raw, query = raw[:i], raw[i+1:]
	}
	if !strings.HasPrefix(raw, PreviewPathPrefix) {
		return "", window, fmt.Errorf("%w: %q 不是预览地址（前缀 %q）", ErrPreviewRefused, requestPath, PreviewPathPrefix)
	}
	escaped := strings.TrimPrefix(raw, PreviewPathPrefix)
	if escaped == "" {
		return "", window, fmt.Errorf("%w: 预览地址缺少产物名", ErrPreviewRefused)
	}
	name, err := url.PathUnescape(escaped)
	if err != nil {
		return "", window, fmt.Errorf("%w: 产物名的转义无法解码：%v", ErrPreviewRefused, err)
	}
	// The store's own bare-name rule, injected. A name that spells a path is
	// refused here, before any filesystem call, so the refusal also means
	// "nothing was opened".
	if err := r.guard(name); err != nil {
		return "", window, fmt.Errorf("%w: %v", ErrPreviewRefused, err)
	}
	if strings.TrimSpace(query) == "" {
		return name, window, nil
	}
	values, qErr := url.ParseQuery(query)
	if qErr != nil {
		return "", window, fmt.Errorf("%w: 预览参数无法解析：%v", ErrPreviewRefused, qErr)
	}
	for key := range values {
		if key != "offset" && key != "limit" {
			return "", window, fmt.Errorf("%w: 未知的预览参数 %q（只认 offset 与 limit）", ErrPreviewRefused, key)
		}
		if len(values[key]) != 1 {
			return "", window, fmt.Errorf("%w: 预览参数 %q 出现了 %d 次，只接受一次", ErrPreviewRefused, key, len(values[key]))
		}
	}
	if rawOffset, ok := values["offset"]; ok {
		offset, convErr := strconv.ParseInt(rawOffset[0], 10, 64)
		if convErr != nil || offset < 0 {
			return "", window, fmt.Errorf("%w: offset=%q 不是非负整数", ErrPreviewRefused, rawOffset[0])
		}
		window.offset, window.present = offset, true
	}
	if rawLimit, ok := values["limit"]; ok {
		limit, convErr := strconv.ParseInt(rawLimit[0], 10, 64)
		if convErr != nil || limit < 1 {
			return "", window, fmt.Errorf("%w: limit=%q 不是正整数", ErrPreviewRefused, rawLimit[0])
		}
		if limit > r.limits.MaxWindow {
			return "", window, fmt.Errorf("%w: limit=%d 超过单窗上限 %d 字节", ErrPreviewRefused, limit, r.limits.MaxWindow)
		}
		window.limit, window.present = limit, true
	}
	return name, window, nil
}

// classifyPreviewArtifact decides one artifact's class from its leading BYTES.
// The extension and any declared type are not inputs: the composer already
// learned (attachments.go) that a declared type is a claim, and a preview that
// believed a name would render an .exe exactly like a .png.
//
// The image half reuses the composer's own sniff table on purpose - one list
// for "what is an image in this product" - while HTML/PDF/text ride
// net/http's DetectContentType, which is the stdlib's mainstream answer and
// avoids a hand-rolled HTML heuristic that would call "<div>" a document. A
// text class additionally requires valid UTF-8 in the head: rendering a binary
// as text is mojibake pretending to be a preview.
func classifyPreviewArtifact(head []byte) (PreviewKind, string, string) {
	for _, t := range attachmentTypes {
		if !t.match(head) {
			continue
		}
		if t.kind == "image" {
			return PreviewKindImage, t.mime, ""
		}
		return "", "", fmt.Sprintf("产物是 %s（%s），不在票 192 AC#1 的四类（HTML／图片／PDF／文本）里；预览不做视频播放", t.mime, t.kind)
	}
	raw := http.DetectContentType(head)
	mediaType, params, parseErr := mime.ParseMediaType(raw)
	if parseErr != nil {
		return "", "", fmt.Sprintf("产物类型无法判定（%s）；预览不猜", raw)
	}
	switch mediaType {
	case "application/pdf":
		return PreviewKindPDF, "application/pdf", ""
	case "text/html":
		return PreviewKindHTML, "text/html; charset=utf-8", ""
	case "text/plain", "text/xml", "application/xml":
		if cs := params["charset"]; cs != "" && !strings.EqualFold(cs, "utf-8") {
			return "", "", fmt.Sprintf("产物看起来是 %s 文本，不是 UTF-8（%s）；预览只显示 UTF-8 文本，不猜编码", mediaType, cs)
		}
		if !utf8.Valid(head) {
			return "", "", "产物的开头不是合法 UTF-8；预览只显示 UTF-8 文本，不把二进制当文本画"
		}
		return PreviewKindText, "text/plain; charset=utf-8", ""
	}
	return "", "", fmt.Sprintf("产物类型 %s 不在票 192 AC#1 的四类（HTML／图片／PDF／文本）里；预览不猜它是什么", mediaType)
}

// previewHeaders is the complete header block one response carries. The class
// decides the policy line, and the comments above explain why the HTML line is
// not cosmetic.
func previewHeaders(kind PreviewKind, contentType string) []string {
	headers := []string{
		"Content-Type: " + contentType,
		"X-Content-Type-Options: nosniff",
		"Cache-Control: no-store",
	}
	switch kind {
	case PreviewKindHTML:
		headers = append(headers, "Content-Security-Policy: "+previewHTMLPolicy)
	case PreviewKindImage:
		headers = append(headers, "Content-Security-Policy: "+previewImagePolicy)
	case PreviewKindPDF:
		headers = append(headers, "Content-Disposition: inline")
	}
	return headers
}

// textFirstView is the truncation view: notice, head, elision line, tail.
func (r *PreviewResolver) textFirstView(f *os.File, name string, size int64, doc PreviewDocument) PreviewDocument {
	headBudget, tailBudget := r.limits.HeadBytes, r.limits.TailBytes
	if size <= headBudget+tailBudget {
		body, err := previewReadAt(f, 0, size)
		if err != nil {
			doc.Status = http.StatusInternalServerError
			doc.Headers = nil
			doc.Reason = fmt.Sprintf("产物 %q 无法读取：%v", name, err)
			return doc
		}
		doc.Status = http.StatusOK
		doc.Body = body
		doc.Complete = true
		return doc
	}
	headRead, err := previewReadAt(f, 0, headBudget+utf8.UTFMax-1)
	if err != nil {
		doc.Status = http.StatusInternalServerError
		doc.Headers = nil
		doc.Reason = fmt.Sprintf("产物 %q 无法读取：%v", name, err)
		return doc
	}
	headEnd := previewRuneCut(headRead, int(headBudget))
	head := headRead[:headEnd]
	tailFrom := size - tailBudget
	if tailFrom < int64(headEnd) {
		tailFrom = int64(headEnd)
	}
	tailRead, err := previewReadAt(f, tailFrom, size-tailFrom)
	if err != nil {
		doc.Status = http.StatusInternalServerError
		doc.Headers = nil
		doc.Reason = fmt.Sprintf("产物 %q 无法读取：%v", name, err)
		return doc
	}
	tailSkip := previewRuneStart(tailRead)
	tail := tailRead[tailSkip:]
	elided := size - int64(len(head)) - int64(len(tail))
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s %s · 共 %d 字节 · 显示开头 %d 与结尾 %d 字节，中间 %d 字节未显示\n",
		previewNoticeTotal, name, size, len(head), len(tail), elided)
	fmt.Fprintf(&b, "%s %s · 每窗最多 %d 字节，一直取到结尾；这段指针能把全文一字节不差地取回\n",
		previewNoticeResume, PreviewURLAt(name, int64(len(head)), r.limits.MaxWindow), r.limits.MaxWindow)
	fmt.Fprintf(&b, "%s 开头 %d 字节\n", previewNoticeBody, len(head))
	b.Write(head)
	if len(head) > 0 && head[len(head)-1] != '\n' {
		// Display only: keeps the next notice line from being glued onto the
		// artifact's last line. The byte-exact copy is the window route.
		b.WriteByte('\n')
	}
	fmt.Fprintf(&b, "%s 中间 %d 字节未显示，用上面的续读地址取\n", previewNoticeElided, elided)
	fmt.Fprintf(&b, "%s 结尾 %d 字节\n", previewNoticeBody, len(tail))
	b.Write(tail)
	doc.Status = http.StatusOK
	doc.Body = b.Bytes()
	doc.Offset = 0
	doc.NextOffset = int64(len(head))
	doc.Complete = false
	return doc
}

// textWindow is one byte window with no envelope, so a walker concatenates
// windows into the original file.
func (r *PreviewResolver) textWindow(f *os.File, name string, size int64, window previewWindow, doc PreviewDocument) PreviewDocument {
	if window.offset >= size {
		// Past the end: an empty range that says the walk is over.
		doc.Status = http.StatusOK
		doc.Offset = size
		doc.Complete = true
		return doc
	}
	want := window.limit
	if want <= 0 {
		want = r.limits.MaxWindow
	}
	if window.offset+want > size {
		want = size - window.offset
	}
	buf, err := previewReadAt(f, window.offset, want+utf8.UTFMax-1)
	if err != nil {
		doc.Status = http.StatusInternalServerError
		doc.Headers = nil
		doc.Reason = fmt.Sprintf("产物 %q 无法读取：%v", name, err)
		return doc
	}
	cut := previewRuneCut(buf, int(want))
	doc.Status = http.StatusOK
	doc.Body = buf[:cut]
	doc.Offset = window.offset
	next := window.offset + int64(cut)
	if next >= size {
		doc.Complete = true
		return doc
	}
	doc.NextOffset = next
	return doc
}

// QuoteForModel is the ONE door through which preview text may become model
// context (ticket 192 AC#2). It exists because the alternative shape - the
// page sending whatever it likes and the host stamping whoever asked - would
// let a compromised renderer choose what the C25 index believes: it could
// launder arbitrary text into a task.output mark and either bury a real leak
// in noise or plant fragments that later make innocent calls look tainted.
//
// So the door proves three things before it marks: the text is non-empty, it is
// somewhere inside the named artifact's own bytes (read back here, not taken
// from the page), and the artifact is small enough to have been checked at all.
// Only then does it Mark under PreviewSourceTool with origin artifact:<name>,
// and the mark is what makes the later outgoing call catch (the exfil half of
// the same pipe: prov.CheckText/prov.Inspect see this content from now on).
//
// A Mark that records nothing (content with no matchable characters after
// normalization) is a refusal, not a success: the caller would otherwise send
// text into the model believing it was stamped.
func (r *PreviewResolver) QuoteForModel(ctx context.Context, marker TaintMarker, scopeID, artifact, quoted string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("%w: 引用请求已被取消", ErrPreviewRefused)
	}
	if marker == nil {
		return "", fmt.Errorf("%w: 没有 C25 引擎，引用无法盖戳，因此不许进模型上下文", ErrPreviewRefused)
	}
	if strings.TrimSpace(scopeID) == "" {
		return "", fmt.Errorf("%w: 没有 scope，引用无法盖戳（缺信息按拒绝处理）", ErrPreviewRefused)
	}
	if strings.TrimSpace(quoted) == "" {
		return "", fmt.Errorf("%w: 引用是空的，没有可盖戳的内容", ErrPreviewRefused)
	}
	if err := r.guard(artifact); err != nil {
		return "", fmt.Errorf("%w: 引用的产物名不合格：%v", ErrPreviewRefused, err)
	}
	full := filepath.Join(r.dir.ArtifactsDir(), artifact)
	info, statErr := os.Lstat(full)
	if statErr != nil {
		return "", fmt.Errorf("%w: 引用的产物 %q 读不到：%v", ErrPreviewRefused, artifact, statErr)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("%w: 引用的产物 %q 不是普通文件", ErrPreviewRefused, artifact)
	}
	if info.Size() > r.limits.QuoteMax {
		return "", fmt.Errorf("%w: 产物 %q 有 %d 字节，超过核对上限 %d；无法证明这段文字出自它，不盖戳",
			ErrPreviewRefused, artifact, info.Size(), r.limits.QuoteMax)
	}
	f, openErr := os.Open(full)
	if openErr != nil {
		return "", fmt.Errorf("%w: 引用的产物 %q 打不开：%v", ErrPreviewRefused, artifact, openErr)
	}
	defer func() { _ = f.Close() }()
	body, readErr := previewReadAt(f, 0, info.Size())
	if readErr != nil {
		return "", fmt.Errorf("%w: 引用的产物 %q 无法读取：%v", ErrPreviewRefused, artifact, readErr)
	}
	if !bytes.Contains(body, []byte(quoted)) {
		return "", fmt.Errorf("%w: 这段文字不在产物 %q 里；预览不能拿它的来源替别处的文字盖戳", ErrPreviewRefused, artifact)
	}
	if !marker.Mark(scopeID, PreviewSourceTool, "artifact:"+artifact, quoted) {
		return "", fmt.Errorf("%w: 这段文字没有可匹配的内容，盖戳等于没盖，因此不进模型上下文", ErrPreviewRefused)
	}
	return quoted, nil
}

// --- small byte helpers -------------------------------------------------------

// previewReadAt reads at most n bytes at off, tolerating a short read at EOF.
func previewReadAt(f *os.File, off, n int64) ([]byte, error) {
	if n <= 0 {
		return nil, nil
	}
	buf := make([]byte, n)
	got, err := f.ReadAt(buf, off)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	return buf[:got], nil
}

// previewRuneCut returns the largest index at or below want that does not split
// a UTF-8 rune, given that b carries want plus up to utf8.UTFMax-1 look-ahead
// bytes. When every byte below want is a continuation byte (a caller-chosen
// offset in the middle of a rune) it returns want rather than 0: a window that
// cannot advance is a walk that never ends.
func previewRuneCut(b []byte, want int) int {
	if want >= len(b) {
		return len(b)
	}
	i := want
	for i > 0 && b[i]&0xC0 == 0x80 {
		i--
	}
	if i == 0 {
		return want
	}
	return i
}

// previewRuneStart returns the index of the first rune start in b, i.e. how
// many leading continuation bytes to skip so the tail begins on a character.
func previewRuneStart(b []byte) int {
	i := 0
	for i < len(b) && b[i]&0xC0 == 0x80 {
		i++
	}
	return i
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
