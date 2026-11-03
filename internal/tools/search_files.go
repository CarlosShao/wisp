package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// D34 `search.files` (SPEC-07 §3 row:2549 / docs/PLAN.md:2549 - "按名找文件",
// L0, capability fs.read, "白名单内"; ticket 24:22 "name glob within
// allowlist", :26-27 "blacklisted targets excluded even from search RESULTS
// listing", :27 "cancelled mid-scan -> partial results marked truncated").
//
// The match is a NAME glob - the pattern is tested against each entry's
// basename, never against a path, so the only filesystem decision this tool
// makes is "walk inside the canonical root the C26 resolver already approved".
// Glob syntax is the frozen `path/filepath.Match` shape (same syntax the shell
// and `find -name` use, and what the mainstream harness Glob tools expose):
// `*`, `?`, `[chars]`. Case-insensitive, matching the literal semantics of the
// `search.content` sibling. No regex, no boolean query language, no `**`
// path-glob - those would be invented, not landed.
//
// This file imports path/filepath for the pattern MATCHER only. It deliberately
// contains no filepath.Clean/filepath.Abs filesystem decision (D22 ban): the
// root goes through C26 `Canonicalize`, children are joined onto an already
// canonical parent with `joinForListing`, and a pattern carrying a separator is
// refused before anything is touched.
//
// Budgets mirror `search.content` and are enforced here, not by the gate:
// depth (default 8), entries examined (default 2000), results (default 50).
// Hitting any of them - or a cancelled ctx mid-walk - sets `Truncated`, which
// is what the loop's spill layer (D15(3)) reads. Symlink/reparse entries are
// never followed and never printed. Every candidate is vetoed by the A/B
// blacklist BEFORE it is opened and before its name is printed, so a
// blacklisted file cannot leak metadata even as a listing line.

// searchFiles is the D34 `search.files` tool. It reuses SearchDeps because it
// needs exactly the same host face as search.content (resolver, caps,
// blacklist); TextExts and the per-file byte budget are not consulted when the
// match is on a name.
type searchFiles struct{ d SearchDeps }

type searchFilesArgs struct {
	Path       string `json:"path"`
	Pattern    string `json:"pattern"`
	MaxResults int    `json:"max_results"`
}

var searchFilesSchema = JSONSchema(`{"type":"object","properties":{` +
	`"path":{"type":"string","description":"搜索根目录（经 C26 规范化，必须在授权目录内）"},` +
	`"pattern":{"type":"string","description":"文件名片段通配（* ? [字符]，大小写不敏感，不含路径分隔符；例如 *.md 或 *合同*）"},` +
	`"max_results":{"type":"integer","minimum":1,"description":"本次最多返回多少个文件"}}` +
	`,"required":["path","pattern"],"additionalProperties":false}`)

func (searchFiles) Name() string { return "search.files" }

// Description implements Tool.
func (searchFiles) Description() string {
	return "在授权目录内按文件名通配查找文件（只列路径；敏感清单路径不列出；不跟随符号链接）"
}

// Parameters implements C1 Tool.
func (searchFiles) Parameters() JSONSchema { return searchFilesSchema }

// Execute implements C1 Tool.
func (t searchFiles) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a searchFilesArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	if strings.TrimSpace(a.Path) == "" {
		return Result{Text: "缺少 path 参数", IsError: true}, nil
	}
	glob, err := compileNameGlob(a.Pattern)
	if err != nil {
		return Result{Text: "拒绝搜索（fail-closed）：" + err.Error(), IsError: true}, nil
	}
	root, err := t.d.resolveSearch(a.Path)
	if err != nil {
		return Result{Text: "拒绝搜索（fail-closed）：" + err.Error(), IsError: true}, nil
	}
	if cls := t.d.classify(root); cls != risk.ClassNone {
		return Result{Text: "搜索根目录在敏感清单内（" + cls.String() + " 级），按 fail-closed 拒绝", IsError: true}, nil
	}
	st, err := os.Stat(root)
	if err != nil {
		return Result{Text: "打开目录失败：" + err.Error(), IsError: true, Origin: root}, nil
	}
	if !st.IsDir() {
		return Result{Text: "path 不是目录（按名查找需要一个授权内的目录根）", IsError: true, Origin: root}, nil
	}

	limit := t.d.resultsCap()
	if a.MaxResults > 0 && a.MaxResults < limit {
		limit = a.MaxResults
	}
	scan := filesScan{d: t.d, glob: glob, limit: limit}
	out, err := scan.run(ctx, root)
	if err != nil {
		return Result{Text: "搜索失败：" + err.Error(), IsError: true, Origin: root}, nil
	}

	truncated := out.hitResultCap || out.hitEntryCap || out.aborted != ""
	head := fmt.Sprintf("search.files 命中 %d 个文件（扫描 %d 项，上限 %d",
		len(out.matches), out.examined, limit)
	if out.aborted != "" {
		head += fmt.Sprintf("，已中止：%s，truncated", out.aborted)
	} else if out.hitResultCap || out.hitEntryCap {
		head += "，达到上限，truncated"
	}
	if out.excluded > 0 {
		head += fmt.Sprintf("，敏感清单排除 %d 项（不列出）", out.excluded)
	}
	if out.skippedLinks > 0 {
		head += fmt.Sprintf("，符号链接跳过 %d 项（不跟随、不列出）", out.skippedLinks)
	}
	head += ")\n"

	var b strings.Builder
	b.WriteString(head)
	for _, p := range out.matches {
		fmt.Fprintf(&b, "%s\n", p)
	}
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("search.files 命中 %d 个文件", len(out.matches)))
	}
	if len(out.matches) == 0 {
		// An honest zero, same shape as the search.content sibling: a miss is
		// a valid answer (not IsError), but a truncated zero still says so.
		return Result{Text: strings.TrimRight(head, "\n"), Truncated: truncated, Origin: root}, nil
	}
	return Result{Text: b.String(), Truncated: truncated, Origin: root}, nil
}

// compileNameGlob validates the pattern and lowers it for a case-insensitive
// basename match. An empty pattern is refused: "find everything under the root"
// is a listing, not a search, and it would turn one call into a whole-tree dump
// (the same veto search.content applies to an empty query).
func compileNameGlob(pattern string) (string, error) {
	p := strings.TrimSpace(pattern)
	if p == "" {
		return "", errors.New("缺少 pattern 参数（空模式等于全盘列举，拒绝）")
	}
	if strings.ContainsAny(p, `/\`) {
		return "", errors.New("pattern 只匹配文件名，不能含路径分隔符（跨目录请用 * 通配）")
	}
	// filepath.Match reports ErrBadPattern only once it has consumed a chunk
	// longer than the name, so the probe name below exists purely to surface a
	// malformed pattern instead of silently matching nothing.
	if _, err := filepath.Match(p, strings.Repeat("p", 128)); err != nil {
		return "", fmt.Errorf("pattern 语法无效：%w", err)
	}
	return strings.ToLower(p), nil
}

// filesOutcome is what one name scan reports back.
type filesOutcome struct {
	matches      []string
	examined     int
	excluded     int
	skippedLinks int
	hitResultCap bool
	hitEntryCap  bool
	aborted      string
}

type filesScan struct {
	d     SearchDeps
	glob  string
	limit int
}

// run walks the tree with an explicit stack (no recursion, so a deep tree
// cannot grow the goroutine stack) under a depth cap and an entries-examined
// cap. Directories are traversed but never matched: this row finds files.
func (s *filesScan) run(ctx context.Context, root string) (filesOutcome, error) {
	var out filesOutcome
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
			continue // an unreadable directory is not a reason to fail the search
		}
		for _, e := range entries {
			if out.examined >= s.d.filesCap() {
				out.hitEntryCap = true
				return out, nil
			}
			out.examined++
			full := joinForListing(cur.path, e.Name())
			if cls := s.d.classify(full); cls != risk.ClassNone {
				out.excluded++ // dropped before its name is ever printed
				continue
			}
			if e.IsDir() {
				if cur.depth+1 < s.d.depthCap() {
					stack = append(stack, frame{path: full, depth: cur.depth + 1})
				}
				continue
			}
			if e.Type()&fs.ModeSymlink != 0 {
				out.skippedLinks++ // never followed, never printed
				continue
			}
			ok, err := filepath.Match(s.glob, strings.ToLower(e.Name()))
			if err != nil {
				return out, err
			}
			if ok && len(out.matches) < s.limit {
				out.matches = append(out.matches, full)
			}
			if len(out.matches) >= s.limit {
				out.hitResultCap = true
				return out, nil
			}
		}
	}
	return out, nil
}

// SearchFilesDecl is the host-side declaration for search.files: capability
// fs.read (D34 row), R1 floor L0 (C19 owns the verdict), path judged through
// C26. `pattern` is NOT a PathParams entry on purpose - it is matched against
// basenames only, so it never becomes a path the gate would have to judge.
func SearchFilesDecl() Decl {
	return Decl{
		Capabilities: []Capability{CapFSRead},
		Needs:        []Capability{CapFSRead},
		Declared:     risk.L0,
		PathParams:   []string{"path"},
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// BuiltinSearchFilesEntries returns the search.files entry for the Builtin
// slot. Like the document and search.content entries it is not wired into
// bridge.go in this limb - see docs/evidence/s1/24r3-impl.md §5.
func BuiltinSearchFilesEntries(d SearchDeps) []Entry {
	return []Entry{{Tool: searchFiles{d: d}, Decl: SearchFilesDecl()}}
}
