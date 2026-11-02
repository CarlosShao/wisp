package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// D34 `search.content` (SPEC-07 §3 row: 全文检索, L0, capability fs.read,
// "结果打 taint"; ticket:22-27 "full-text within allowlist; results
// taint-marked; huge result sets -> spill", "blacklisted targets excluded even
// from search RESULTS listing", "cancelled mid-scan -> partial results marked
// truncated").
//
// On-demand only: no index, no daemon (ticket Out of scope:29-31). The scan
// walks the C26-canonical root, reads text-ish files, and stops at its own caps
// or at ctx - whichever comes first.
//
// Match semantics are a literal, case-insensitive substring. SPEC-07 §3 gives
// this row only its name, level, capability and taint note (24-r1 §4 said the
// parameter/scan面 is thin), so this file deliberately keeps the surface at the
// minimum the ticket states: one root, one query, one result cap, plus the
// allowlist and blacklist vetoes. No regex, no boolean query language - those
// would be invented, not landed.
//
// The D34 sibling rows `search.files` (name glob) and `search.apps` (installed
// apps, capability sysinfo) are NOT here: A1435 scoped this limb to the
// zero-dependency doc+search pair, and the ticket's own parallel-slot note
// splits them. What is still owed is listed in docs/evidence/s1/24r2a-impl.md §4.

// SearchDeps is what search.content needs from the host.
type SearchDeps struct {
	// Paths is the C26 canonicalizer; nil fails closed.
	Paths *PathCanonicalizer
	// MaxResults caps the matches one call returns (default
	// defaultMaxSearchResults). Hitting it sets Truncated, which is the signal
	// the loop's spill layer (D15(3)) reads; this tool never spills to disk on
	// its own.
	MaxResults int
	// MaxFiles caps how many text files are searched (default
	// defaultMaxSearchFiles).
	MaxFiles int
	// MaxDepth caps how far below the root the walk goes (default 8).
	MaxDepth int
	// MaxFileBytes is the budget per file (default 256 KiB, the same figure as
	// one fs.read). A larger file is skipped whole and never searched
	// half-way.
	MaxFileBytes int
	// TextExts are the lowercase extensions considered text. Empty selects
	// DefaultSearchExts().
	TextExts []string
	// Classify is the A/B blacklist seam; nil means risk.Classify. A
	// classified path is dropped before it is opened AND before its name is
	// printed, so a blacklisted file never appears in a listing (ticket:25-26).
	Classify func(string) risk.Class
}

const (
	defaultMaxSearchResults = 50
	defaultMaxSearchFiles   = 2000
	defaultMaxSearchDepth   = 8
	defaultSearchSnippet    = 200
)

// DefaultSearchExts is the text-ish extension set the on-demand scan reads.
// Anything outside it (binaries, documents, media) is skipped by extension, and
// a NUL byte in what it does read is a second, content-based veto.
func DefaultSearchExts() []string {
	return []string{
		"txt", "md", "markdown", "mdx", "log", "csv", "tsv", "json", "jsonl",
		"yaml", "yml", "toml", "ini", "cfg", "conf", "xml", "html",
		"htm", "css", "go", "py", "js", "jsx", "ts", "tsx", "rs", "java",
		"c", "h", "cc", "cpp", "hpp", "cs", "rb", "php", "sh", "ps1", "bat",
		"sql", "proto", "srt", "vtt",
	}
}

func (d SearchDeps) resultsCap() int {
	if d.MaxResults > 0 {
		return d.MaxResults
	}
	return defaultMaxSearchResults
}

func (d SearchDeps) filesCap() int {
	if d.MaxFiles > 0 {
		return d.MaxFiles
	}
	return defaultMaxSearchFiles
}

func (d SearchDeps) depthCap() int {
	if d.MaxDepth > 0 {
		return d.MaxDepth
	}
	return defaultMaxSearchDepth
}

func (d SearchDeps) fileCap() int64 {
	if d.MaxFileBytes > 0 {
		return int64(d.MaxFileBytes)
	}
	return defaultMaxReadBytes
}

func (d SearchDeps) exts() map[string]bool {
	src := d.TextExts
	if len(src) == 0 {
		src = DefaultSearchExts()
	}
	m := make(map[string]bool, len(src))
	for _, e := range src {
		e = strings.ToLower(strings.TrimPrefix(strings.TrimSpace(e), "."))
		if e != "" {
			m[e] = true
		}
	}
	return m
}

func (d SearchDeps) classify(canonical string) risk.Class {
	if d.Classify != nil {
		return d.Classify(canonical)
	}
	return risk.Classify(canonical)
}

// resolveSearch canonicalizes the scan root through C26 (the same fail-closed
// rules as the fs family) and then enforces the ticket's "within allowlist"
// clause inside the tool too: a root outside the authorized dirs is refused
// rather than turned into a whole-disk scan. The bridge still judges the same
// canonical spelling through Decl.PathParams for R2/R3, so a search that is
// legitimately out of scope goes through the gate, not around it.
func (d SearchDeps) resolveSearch(raw string) (string, error) {
	if d.Paths == nil {
		return "", errors.New("search: no C26 resolver configured (fail-closed)")
	}
	c, err := d.Paths.Canonicalize(raw)
	if err != nil {
		return "", err
	}
	if strings.Contains(c, "??") {
		return "", errors.New("search: unresolved path component")
	}
	if !d.Paths.InAllowlist(c) {
		return "", errors.New("search: 根目录不在授权目录内（fail-closed，不做全盘扫描）")
	}
	return c, nil
}

type searchContent struct{ d SearchDeps }

type searchContentArgs struct {
	Path       string `json:"path"`
	Query      string `json:"query"`
	MaxResults int    `json:"max_results"`
}

var searchContentSchema = JSONSchema(`{"type":"object","properties":{` +
	`"path":{"type":"string","description":"检索根目录（经 C26 规范化，必须在授权目录内）"},` +
	`"query":{"type":"string","description":"要查找的字面文本（大小写不敏感，非正则）"},` +
	`"max_results":{"type":"integer","minimum":1,"description":"本次最多返回多少条命中"}}` +
	`,"required":["path","query"],"additionalProperties":false}`)

func (searchContent) Name() string { return "search.content" }

// Description implements Tool.
func (searchContent) Description() string {
	return "在授权目录内做按需全文检索（字面匹配，结果含路径与行号；敏感清单路径不列出）"
}

func (searchContent) Parameters() JSONSchema { return searchContentSchema }

// Execute implements C1 Tool.
func (t searchContent) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a searchContentArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	if strings.TrimSpace(a.Path) == "" {
		return Result{Text: "缺少 path 参数", IsError: true}, nil
	}
	if a.Query == "" {
		return Result{Text: "缺少 query 参数（空查询不做全盘列举）", IsError: true}, nil
	}
	root, err := t.d.resolveSearch(a.Path)
	if err != nil {
		return Result{Text: "拒绝检索（fail-closed）：" + err.Error(), IsError: true}, nil
	}
	if cls := t.d.classify(root); cls != risk.ClassNone {
		return Result{Text: "检索根目录在敏感清单内（" + cls.String() + " 级），按 fail-closed 拒绝", IsError: true}, nil
	}
	st, err := os.Stat(root)
	if err != nil {
		return Result{Text: "打开目录失败：" + err.Error(), IsError: true, Origin: root}, nil
	}
	if !st.IsDir() {
		return Result{Text: "path 不是目录（全文检索需要一个授权内的目录根）", IsError: true, Origin: root}, nil
	}

	limit := t.d.resultsCap()
	if a.MaxResults > 0 && a.MaxResults < limit {
		limit = a.MaxResults
	}
	scan := searchScan{d: t.d, needle: strings.ToLower(a.Query), limit: limit, textExt: t.d.exts()}
	out, err := scan.run(ctx, root)
	if err != nil {
		return Result{Text: "检索失败：" + err.Error(), IsError: true, Origin: root}, nil
	}

	truncated := out.hitResultCap || out.hitFileCap || out.aborted != ""
	head := fmt.Sprintf("search.content 命中 %d 处（扫描 %d 文件，上限 %d",
		len(out.matches), out.scanned, limit)
	if out.aborted != "" {
		head += fmt.Sprintf("，已中止：%s，truncated", out.aborted)
	} else if out.hitResultCap || out.hitFileCap {
		head += "，达到上限，truncated"
	}
	if out.excluded > 0 {
		head += fmt.Sprintf("，敏感清单排除 %d 项（不列出）", out.excluded)
	}
	head += ")\n"

	var b strings.Builder
	b.WriteString(head)
	for _, m := range out.matches {
		fmt.Fprintf(&b, "%s:%d:%s\n", m.path, m.line, m.snippet)
	}
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("search.content 命中 %d 处", len(out.matches)))
	}
	if len(out.matches) == 0 {
		// An honest zero: the scan ran and found nothing. A miss is a valid
		// answer, so it is not IsError - but it is never dressed up as a
		// listing, and a truncated zero still carries Truncated.
		return Result{Text: strings.TrimRight(head, "\n"), Truncated: truncated, Origin: root}, nil
	}
	return Result{Text: b.String(), Truncated: truncated, Origin: root}, nil
}

// searchHit is one printed match: a canonical path, a 1-based line, a clipped
// snippet.
type searchHit struct {
	path    string
	line    int
	snippet string
}

// searchOutcome is what one scan reports back, including the caps it hit.
type searchOutcome struct {
	matches      []searchHit
	scanned      int
	excluded     int
	hitResultCap bool
	hitFileCap   bool
	aborted      string
}

type searchScan struct {
	d       SearchDeps
	needle  string
	limit   int
	textExt map[string]bool
}

// run walks the tree with an explicit stack (no recursion, so a deep tree
// cannot grow the goroutine stack) and a depth cap. Directories that are
// symlinks/reparse points are never followed.
func (s *searchScan) run(ctx context.Context, root string) (searchOutcome, error) {
	var out searchOutcome
	type frame struct {
		path  string
		depth int
	}
	stack := []frame{{path: root, depth: 0}}
	for len(stack) > 0 {
		if err := ctx.Err(); err != nil {
			out.aborted = err.Error()
			return out, nil
		}
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		entries, err := os.ReadDir(cur.path)
		if err != nil {
			continue // an unreadable directory is not a reason to fail the scan
		}
		for _, e := range entries {
			full := joinForListing(cur.path, e.Name())
			if cls := s.d.classify(full); cls != risk.ClassNone {
				out.excluded++
				continue
			}
			info, err := e.Info()
			if err != nil {
				continue
			}
			switch {
			case e.IsDir():
				if cur.depth+1 >= s.d.depthCap() {
					continue
				}
				stack = append(stack, frame{path: full, depth: cur.depth + 1})
			case info.Mode()&fs.ModeSymlink != 0:
				// Never followed, never read, never printed.
			default:
				if out.scanned >= s.d.filesCap() {
					out.hitFileCap = true
					return out, nil
				}
				if !s.textExt[strings.ToLower(strings.TrimPrefix(extOfPath(full), "."))] {
					continue
				}
				if info.Size() > s.d.fileCap() {
					continue // skipped whole, never searched half-way
				}
				out.scanned++
				stop, err := s.scanFile(ctx, full, &out)
				if err != nil {
					return out, err
				}
				if stop {
					return out, nil
				}
			}
		}
	}
	return out, nil
}

// scanFile reads one candidate file and appends its matches. The bool return
// says "stop the whole scan" (result cap reached or ctx gone).
func (s *searchScan) scanFile(ctx context.Context, path string, out *searchOutcome) (bool, error) {
	if err := ctx.Err(); err != nil {
		out.aborted = err.Error()
		return true, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return false, nil
	}
	if hasNUL(raw) {
		return false, nil // binary inside the text-extension veto: skip
	}
	text := strings.ReplaceAll(string(raw), "\r\n", "\n")
	for i, line := range strings.Split(text, "\n") {
		if !strings.Contains(strings.ToLower(line), s.needle) {
			continue
		}
		out.matches = append(out.matches, searchHit{path: path, line: i + 1, snippet: clipSnippet(line)})
		if len(out.matches) >= s.limit {
			out.hitResultCap = true
			return true, nil
		}
	}
	return false, nil
}

// hasNUL reports whether the bytes look like a binary payload.
func hasNUL(raw []byte) bool {
	for _, b := range raw {
		if b == 0 {
			return true
		}
	}
	return false
}

// clipSnippet keeps a printable single line within defaultSearchSnippet bytes
// without splitting a UTF-8 rune.
func clipSnippet(line string) string {
	line = strings.TrimSpace(strings.TrimRight(line, "\r"))
	if len(line) <= defaultSearchSnippet {
		return line
	}
	cut := defaultSearchSnippet
	for cut > 0 && !isRuneStart(line, cut) {
		cut--
	}
	return line[:cut] + "..."
}

// isRuneStart reports whether i begins a UTF-8 rune (not a continuation byte).
func isRuneStart(s string, i int) bool {
	return i < len(s) && s[i]&0xC0 != 0x80
}

// SearchContentDecl is the host-side declaration for search.content: capability
// fs.read, R1 floor L0 (C19 owns the verdict), path judged through C26
// (ticket:26 - every search is path-resolved via C26).
func SearchContentDecl() Decl {
	return Decl{
		Capabilities: []Capability{CapFSRead},
		Needs:        []Capability{CapFSRead},
		Declared:     risk.L0,
		PathParams:   []string{"path"},
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// BuiltinSearchContentEntries returns the search.content entry for the Builtin
// slot. Like the document entry it is not wired into bridge.go in this limb -
// see docs/evidence/s1/24r2a-impl.md §4.
func BuiltinSearchContentEntries(d SearchDeps) []Entry {
	return []Entry{{Tool: searchContent{d: d}, Decl: SearchContentDecl()}}
}
