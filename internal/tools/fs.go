package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// The D34 fs family lands slice by slice. This file holds the L0 pair only:
//
//	fs.read  L0 inside [fs] allowed_dirs, L2 out of scope via R2
//	fs.list  L0 / L2, same rule, capability fs.read
//
// fs.write / fs.trash / fs.move (and the delete_enabled-gated fs.delete) are
// ticket 20 segment 2: they need the D31 temp+atomic-rename writer, the
// applied-steps report and the [fs] allowed_dirs first-use ask flow, none of
// which the L0 pair exercises. Result.AppliedSteps and Decl.Provider are the
// seams they attach to, so landing them does not reshape anything here.

// FSDeps is what an fs tool needs from the host.
type FSDeps struct {
	// Paths is the C26 canonicalizer. An fs tool MUST reach the disk through
	// it: filepath.Clean/Abs outside risk's resolver is a CI-failing pattern
	// (D22), and skipping it would let a spelled path open something the risk
	// verdict never looked at.
	Paths *PathCanonicalizer
	// MaxReadBytes caps one fs.read (default 256 KiB; the loop's spill layer
	// (D15(3)) owns what happens above the context budget, not this tool).
	MaxReadBytes int
	// MaxListEntries caps one fs.list (default 500).
	MaxListEntries int
}

const (
	defaultMaxReadBytes   = 256 * 1024
	defaultMaxListEntries = 500
)

func (d FSDeps) readCap() int {
	if d.MaxReadBytes > 0 {
		return d.MaxReadBytes
	}
	return defaultMaxReadBytes
}

func (d FSDeps) listCap() int {
	if d.MaxListEntries > 0 {
		return d.MaxListEntries
	}
	return defaultMaxListEntries
}

// open resolves a model-supplied path through C26 and hands back the canonical
// spelling the OS then opens. The resolution is the same call R2/R3 judged, so
// the bytes read are the bytes the gate approved.
func (d FSDeps) open(raw string) (string, error) {
	if d.Paths == nil {
		return "", errors.New("fs: no C26 resolver configured (fail-closed)")
	}
	c, err := d.Paths.Canonicalize(raw)
	if err != nil {
		return "", err
	}
	if strings.Contains(c, "??") {
		return "", errors.New("fs: unresolved path component")
	}
	return c, nil
}

// ---------------------------------------------------------------------------
// fs.read (C1 Tool, D34 row 1)
// ---------------------------------------------------------------------------

type fsRead struct{ d FSDeps }

type fsReadArgs struct {
	Path     string `json:"path"`
	MaxBytes int    `json:"max_bytes"`
}

// fsReadSchema is the C1 Parameters document.
var fsReadSchema = JSONSchema(`{"type":"object","properties":{` +
	`"path":{"type":"string","description":"文件路径（可含 ~ 与环境变量，经 C26 规范化）"},` +
	`"max_bytes":{"type":"integer","minimum":1,"description":"本次最多读多少字节"}}` +
	`,"required":["path"],"additionalProperties":false}`)

func (fsRead) Name() string { return "fs.read" }
func (fsRead) Description() string {
	return "读取一个 UTF-8 文本文件的内容（越界路径会被升为 L2 审批）"
}
func (fsRead) Parameters() JSONSchema { return fsReadSchema }

// Execute implements Tool.
func (t fsRead) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a fsReadArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	if strings.TrimSpace(a.Path) == "" {
		return Result{Text: "缺少 path 参数", IsError: true}, nil
	}
	canon, err := t.d.open(a.Path)
	if err != nil {
		return Result{Text: "路径无法解析（按 fail-closed 拒绝）：" + err.Error(), IsError: true}, nil
	}
	f, err := os.Open(canon)
	if err != nil {
		return Result{Text: "打开失败：" + err.Error(), IsError: true, Origin: canon}, nil
	}
	defer f.Close() //nolint:errcheck // read-only handle, Close's error is not actionable

	limit := t.d.readCap()
	if a.MaxBytes > 0 && a.MaxBytes < limit {
		limit = a.MaxBytes
	}
	buf, err := io.ReadAll(io.LimitReader(f, int64(limit)+1))
	if err != nil {
		return Result{Text: "读取失败：" + err.Error(), IsError: true, Origin: canon}, nil
	}
	truncated := len(buf) > limit
	if truncated {
		buf = buf[:limit]
	}
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("fs.read 读取 %d 字节", len(buf)))
	}
	return Result{
		Text:      string(buf),
		Truncated: truncated,
		Origin:    canon,
	}, nil
}

// ---------------------------------------------------------------------------
// fs.list (C1 Tool, D34 row 2)
// ---------------------------------------------------------------------------

type fsList struct{ d FSDeps }

type fsListArgs struct {
	Path  string `json:"path"`
	Limit int    `json:"limit"`
}

var fsListSchema = JSONSchema(`{"type":"object","properties":{` +
	`"path":{"type":"string","description":"目录路径（经 C26 规范化）"},` +
	`"limit":{"type":"integer","minimum":1,"description":"最多列多少条目"}}` +
	`,"required":["path"],"additionalProperties":false}`)

func (fsList) Name() string           { return "fs.list" }
func (fsList) Description() string    { return "列出一个目录里的条目（名称/类型/大小）" }
func (fsList) Parameters() JSONSchema { return fsListSchema }

// Execute implements Tool.
func (t fsList) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a fsListArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	if strings.TrimSpace(a.Path) == "" {
		return Result{Text: "缺少 path 参数", IsError: true}, nil
	}
	canon, err := t.d.open(a.Path)
	if err != nil {
		return Result{Text: "路径无法解析（按 fail-closed 拒绝）：" + err.Error(), IsError: true}, nil
	}
	dir, err := os.Open(canon)
	if err != nil {
		return Result{Text: "打开目录失败：" + err.Error(), IsError: true, Origin: canon}, nil
	}
	defer dir.Close() //nolint:errcheck // read-only handle

	names, err := dir.Readdirnames(-1)
	if err != nil && len(names) == 0 {
		return Result{Text: "读取目录失败：" + err.Error(), IsError: true, Origin: canon}, nil
	}
	sort.Strings(names)
	limit := t.d.listCap()
	if a.Limit > 0 && a.Limit < limit {
		limit = a.Limit
	}
	truncated := len(names) > limit
	if truncated {
		names = names[:limit]
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s (%d 条目%s)\n", canon, len(names), sizeSuffix(truncated, limit))
	for _, n := range fi(names, canon, ctx) {
		fmt.Fprintln(&b, n)
	}
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("fs.list 列出 %d 条目", len(names)))
	}
	return Result{Text: b.String(), Truncated: truncated, Origin: canon}, nil
}

// sizeSuffix annotates a capped listing.
func sizeSuffix(truncated bool, limit int) string {
	if truncated {
		return fmt.Sprintf("，已达上限 %d，truncated", limit)
	}
	return ""
}

// fi renders one line per entry with its stat, stopping early if the call's
// ctx is gone (a cancelled listing must not keep hitting the disk).
func fi(names []string, dir string, ctx context.Context) []string {
	out := make([]string, 0, len(names))
	for _, n := range names {
		if ctx.Err() != nil {
			return append(out, "(已中止："+ctx.Err().Error()+")")
		}
		line := n
		if st, err := os.Lstat(joinForListing(dir, n)); err == nil {
			kind := "f"
			switch {
			case st.IsDir():
				kind = "d"
			case st.Mode()&os.ModeSymlink != 0:
				kind = "l"
			}
			line = fmt.Sprintf("%s %s %d %s", kind, n, st.Size(),
				st.ModTime().Format(time.RFC3339))
		} else {
			line = fmt.Sprintf("? %s stat-error", n)
		}
		out = append(out, line)
	}
	return out
}

// joinForListing appends one directory entry to a C26-canonical directory. It
// is string concatenation onto an ALREADY canonical parent (entries come from
// Readdirnames, so they carry no separators), which is why it needs no Clean
// step and must not grow into a general path joiner.
func joinForListing(dir, name string) string {
	if dir == "" {
		return name
	}
	if strings.HasSuffix(dir, `\`) || strings.HasSuffix(dir, "/") {
		return dir + name
	}
	return dir + `\` + name
}

// ---------------------------------------------------------------------------
// registration (D34 capability + R1 floor per row)
// ---------------------------------------------------------------------------

// FSReadDecl is the host-side declaration for fs.read: capability fs.read,
// R1 floor L0 (C19 may and does raise it), path judged through C26.
func FSReadDecl() Decl {
	return Decl{
		Capabilities: []Capability{CapFSRead},
		Needs:        []Capability{CapFSRead},
		Declared:     risk.L0,
		PathParams:   []string{"path"},
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// FSListDecl is the host-side declaration for fs.list (same capability as
// fs.read per the D34 table).
func FSListDecl() Decl {
	return Decl{
		Capabilities: []Capability{CapFSRead},
		Needs:        []Capability{CapFSRead},
		Declared:     risk.L0,
		PathParams:   []string{"path"},
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// BuiltinFSEntries returns the segment-1 fs entries. The composition root
// registers them on a Registry (or hands them to NewBuiltinProvider).
func BuiltinFSEntries(d FSDeps) []Entry {
	return []Entry{
		{Tool: fsRead{d: d}, Decl: FSReadDecl()},
		{Tool: fsList{d: d}, Decl: FSListDecl()},
	}
}
