package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// The D34 write half of the fs family:
//
//	fs.write  new file L1 / overwrite L2 (R8)   D31 temp+atomic-rename writer
//	fs.trash  L1                                the real Shell recycle-bin API
//	fs.move   L1 / cross-volume L2, overwrite L2 (R8)
//	fs.delete NOT REGISTERED unless [fs] delete_enabled=true, then L2 (R8)
//
// Every one of them is a side effect, so every one of them runs through
// sideEffect below, which is the single place that knows what has already
// landed on disk. That ledger is what D31's applied-steps report is built from:
// a cancel is never rendered as a clean undo.

const (
	// defaultMaxWriteBytes caps one fs.write payload. Long content goes through
	// the host-internal artifacts channel (SPEC-07 §3), not a gated tool.
	defaultMaxWriteBytes = 2 * 1024 * 1024
	// defaultWriteChunk is how much content goes out per write step. The size
	// is not performance tuning: it is the granularity of the D31 cancel check,
	// so a multi-megabyte write stays stoppable instead of being one
	// uninterruptible block.
	defaultWriteChunk = 64 * 1024
	// tempPrefix marks a Wisp staging file. It lives in the TARGET's directory
	// so the final rename never crosses a volume - a cross-volume rename is not
	// atomic, which is the whole point of D31.
	tempPrefix = ".wisp-tmp-"
)

// Hooks are the two test seams the write/move/trash paths call at every step
// boundary. Production leaves them nil, so a zero FSDeps is the safe
// configuration rather than a broken one.
type Hooks struct {
	// AtStep observes each boundary (which step is about to run).
	AtStep func(step string)
	// Kill simulates the process dying at that step: returning an error aborts
	// the operation there, with the staging file removed and the target
	// untouched. This is how TestAtomicWriteKillsMidWrite kills a write at an
	// exact byte boundary instead of racing the scheduler.
	Kill func(step string) error
}

// sideEffect is one operation's ledger of what it already did.
type sideEffect struct {
	h     Hooks
	chunk int
	steps []string
	// stopped is set by the first boundary that refused to continue.
	stopped bool
	why     string
}

// ledger opens a side-effect ledger bound to this deps.
func (d FSDeps) ledger() *sideEffect {
	return &sideEffect{h: d.Hooks, chunk: d.chunkSize()}
}

// record appends one FACT to the ledger. The wording is per-fact on purpose:
// the report renders these verbatim under 「已执行」, so a step that did not
// happen must never appear.
func (s *sideEffect) record(format string, args ...any) {
	s.steps = append(s.steps, fmt.Sprintf(format, args...))
}

// at is one stoppable boundary. It returns a non-nil error when the operation
// must not take another step, and it is the ONLY way the fs family aborts: a
// veto, a dead context and a simulated process death all arrive here, so the
// ledger cannot be updated on one path and forgotten on another.
func (s *sideEffect) at(ctx context.Context, step string) error {
	if s.h.AtStep != nil {
		s.h.AtStep(step)
	}
	if s.h.Kill != nil {
		if err := s.h.Kill(step); err != nil {
			return s.stop(step, err.Error())
		}
	}
	if stop, why := Stopped(ctx); stop {
		return s.stop(step, why)
	}
	return nil
}

func (s *sideEffect) stop(step, why string) error {
	if !s.stopped {
		s.stopped = true
		s.why = why
		s.record("在步骤「%s」前停止：%s", step, why)
	}
	return fmt.Errorf("已停止于 %s：%w", step, errors.New(why))
}

// head renders the user-visible headline for a stopped call.
func (s *sideEffect) head(what string) string {
	return fmt.Sprintf("%s已中止（%s）。以下步骤已生效，不会自动回退：\n%s",
		what, s.why, s.render())
}

// stoppedResult closes a ledger that aborted: always an IsError Result, always
// with the ledger attached, never a claim of a clean undo.
func (s *sideEffect) stoppedResult(what string) Result {
	return Result{Text: s.head(what), IsError: true, AppliedSteps: s.snapshot()}
}

func (s *sideEffect) snapshot() []string {
	return append([]string(nil), s.steps...)
}

func (s *sideEffect) render() string {
	if len(s.steps) == 0 {
		return "（无：本次调用在产生任何副作用前就停了）"
	}
	return "- " + strings.Join(s.steps, "\n- ")
}

// writeAll streams a payload in chunk-sized steps, checking every boundary.
func (s *sideEffect) writeAll(ctx context.Context, f *os.File, src io.Reader,
	onUpdate func(string), limit int) (int, error) {
	buf := make([]byte, s.chunk)
	if limit > 0 && len(buf) > limit {
		buf = buf[:limit]
	}
	total := 0
	for {
		n, rerr := src.Read(buf)
		if n > 0 {
			if err := s.at(ctx, fmt.Sprintf("write:%d", total+n)); err != nil {
				return total, err
			}
			w, werr := f.Write(buf[:n])
			total += w
			if werr != nil {
				return total, werr
			}
			if onUpdate != nil {
				onUpdate(fmt.Sprintf("fs.write 已暂存 %d 字节", total))
			}
		}
		if rerr == io.EOF {
			return total, nil
		}
		if rerr != nil {
			return total, rerr
		}
	}
}

// ---------------------------------------------------------------------------
// shared path plumbing
// ---------------------------------------------------------------------------

// canonical resolves one raw path parameter through C26 before any stat or
// open. Every entry point below goes through it: an existence probe on the
// spelled path would ask the OS a question the risk verdict never answered
// (D22's whole objection), and overwrite detection is exactly such a probe.
func (d FSDeps) canonical(raw string) (string, error) {
	if d.Paths == nil {
		return "", errors.New("fs: no C26 resolver configured (fail-closed)")
	}
	return d.Paths.Canonicalize(raw)
}

// dirOf is the parent directory of a C26-canonical path: a string cut at the
// last separator of an ALREADY canonical input. No Clean, no Abs (D22).
func dirOf(canonical string) string {
	c := strings.ReplaceAll(canonical, "/", `\`)
	i := strings.LastIndex(c, `\`)
	if i <= 0 {
		return ""
	}
	return c[:i]
}

// baseOf is the trailing component of a canonical path (string cutting only).
func baseOf(canonical string) string {
	c := strings.ReplaceAll(canonical, "/", `\`)
	if i := strings.LastIndex(c, `\`); i >= 0 {
		return c[i+1:]
	}
	return c
}

// existsViaLstat reports whether one canonical path is there. Lstat, so a
// reparse point is judged as itself and not through its target. A probe error is
// returned rather than folded into "does not exist": an unanswerable existence
// question must land on the expensive side of the gate, not the cheap one.
func existsViaLstat(canonical string) (bool, error) {
	_, err := os.Lstat(canonical)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	default:
		return false, err
	}
}

// ---------------------------------------------------------------------------
// fs.write (D34 row 3, D31)
// ---------------------------------------------------------------------------

type fsWrite struct{ d FSDeps }

type fsWriteArgs struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

var fsWriteSchema = JSONSchema(`{"type":"object","properties":{` +
	`"path":{"type":"string","description":"要写入的文件路径（经 C26 规范化）"},` +
	`"content":{"type":"string","description":"写入的完整内容（UTF-8，覆盖写入）"}}` +
	`,"required":["path","content"],"additionalProperties":false}`)

func (fsWrite) Name() string { return "fs.write" }
func (fsWrite) Description() string {
	return "写入一个文本文件：先写临时文件再原子重命名（D31）。新建文件为 L1，覆盖已有文件为 L2"
}
func (fsWrite) Parameters() JSONSchema { return fsWriteSchema }

// Execute implements Tool.
func (t fsWrite) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a fsWriteArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	if strings.TrimSpace(a.Path) == "" {
		return Result{Text: "缺少 path 参数", IsError: true}, nil
	}
	limit := t.d.writeCap()
	if len(a.Content) > limit {
		return Result{Text: fmt.Sprintf(
			"拒绝写入 %d 字节：超过 fs.write 的单次上限 %d 字节（超长内容属于宿主内部的 artifacts 通道，不是工具参数）",
			len(a.Content), limit), IsError: true}, nil
	}
	target, err := t.d.canonical(a.Path)
	if err != nil {
		return Result{Text: "路径无法解析（按 fail-closed 拒绝）：" + err.Error(), IsError: true}, nil
	}
	res := t.d.stageAndRename(ctx, target, strings.NewReader(a.Content), onUpdate)
	res.Origin = target
	return res, nil
}

// stageAndRename is the D31 writer: stage beside the target, then rename over
// it. The target path is touched by EXACTLY ONE call (os.Rename at the end), so
// any death before that step leaves whatever was there byte-for-byte intact -
// which is what TestAtomicWriteKillsMidWrite asserts against an EXISTING file.
func (d FSDeps) stageAndRename(ctx context.Context, target string, src io.Reader,
	onUpdate func(string)) Result {

	parent := dirOf(target)
	if parent == "" {
		return Result{Text: "拒绝写入卷根目录：" + target, IsError: true}
	}
	se := d.ledger()

	if err := se.at(ctx, "create-temp"); err != nil {
		return se.stoppedResult("写入")
	}
	tmp, err := os.CreateTemp(parent, tempPrefix+"*")
	if err != nil {
		// Nothing was created, so nothing is claimed: no ledger, no undo.
		return Result{Text: "创建临时文件失败：" + err.Error() + "（目标未被改动）", IsError: true}
	}
	tmpName := tmp.Name()
	se.record("在 %s 创建临时文件 %s", parent, baseOf(tmpName))

	// From here on every failure path removes the staging file. The target is
	// still untouched at this point, so that removal is NOT an undo of a user
	// visible change - and the ledger says so in those words, because D31
	// forbids a report that leaves the reader guessing which case they are in.
	discard := func(why error) Result {
		_ = tmp.Close()
		se.record("写入未完成：%v", why)
		if rmErr := os.Remove(tmpName); rmErr == nil {
			se.record("删除临时文件 %s（目标从头到尾未被改动）", baseOf(tmpName))
		} else {
			se.record("临时文件 %s 未能删除，请人工确认：%v", baseOf(tmpName), rmErr)
		}
		return se.stoppedResult("写入")
	}

	written, werr := se.writeAll(ctx, tmp, src, onUpdate, d.writeCap())
	if written > 0 {
		se.record("向临时文件写入 %d 字节", written)
	}
	if werr != nil {
		return discard(werr)
	}
	if err := tmp.Sync(); err != nil {
		return discard(err)
	}
	if err := tmp.Close(); err != nil {
		se.record("临时文件已写出但关闭失败，请人工确认：%v", err)
		return Result{Text: "关闭临时文件失败：" + err.Error(), IsError: true,
			AppliedSteps: se.snapshot()}
	}
	se.record("临时文件已落盘并关闭")

	if err := se.at(ctx, "rename"); err != nil {
		if rmErr := os.Remove(tmpName); rmErr == nil {
			se.record("删除临时文件 %s（目标从头到尾未被改动）", baseOf(tmpName))
		}
		return se.stoppedResult("写入")
	}
	if err := os.Rename(tmpName, target); err != nil {
		se.record("原子重命名失败：%v", err)
		if rmErr := os.Remove(tmpName); rmErr == nil {
			se.record("删除临时文件 %s（目标从头到尾未被改动）", baseOf(tmpName))
		}
		return Result{Text: "原子重命名失败：" + err.Error() + "（目标未被改动）",
			IsError: true, AppliedSteps: se.snapshot()}
	}
	se.record("原子重命名 %s → %s（目标此刻起为新内容）", baseOf(tmpName), target)
	return Result{
		Text: fmt.Sprintf("已写入 %s（%d 字节，临时文件+原子重命名）",
			target, written),
		AppliedSteps: se.snapshot(),
	}
}

// ---------------------------------------------------------------------------
// fs.trash (D34 row 4: L1 BECAUSE the backend is the recycle bin)
// ---------------------------------------------------------------------------

type fsTrash struct{ d FSDeps }

type fsPathArgs struct {
	Path string `json:"path"`
}

var fsTrashSchema = JSONSchema(`{"type":"object","properties":{` +
	`"path":{"type":"string","description":"要放入回收站的路径（经 C26 规范化）"}}` +
	`,"required":["path"],"additionalProperties":false}`)

func (fsTrash) Name() string { return "fs.trash" }
func (fsTrash) Description() string {
	return "把文件/目录放入系统回收站（Shell API，可从回收站还原，因此是 L1 而非 L2）"
}
func (fsTrash) Parameters() JSONSchema { return fsTrashSchema }

// Execute implements Tool.
func (t fsTrash) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a fsPathArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	if strings.TrimSpace(a.Path) == "" {
		return Result{Text: "缺少 path 参数", IsError: true}, nil
	}
	target, err := t.d.canonical(a.Path)
	if err != nil {
		return Result{Text: "路径无法解析（按 fail-closed 拒绝）：" + err.Error(), IsError: true}, nil
	}
	se := t.d.ledger()
	if err := se.at(ctx, "recycle"); err != nil {
		return se.stoppedResult("回收"), nil
	}
	if !recycleBinSupported() {
		// The platform has no bin: refuse. There is no branch below that could
		// turn into an unlink, which is the difference between L1 and a lie.
		return Result{Text: "本平台无 Shell 回收站 API，fs.trash 拒绝执行（绝不退化成删除）：" + target,
			IsError: true}, nil
	}
	detail, err := shellTrash(target)
	if err != nil {
		// The shell said no. Nothing moved, and nothing was deleted either.
		se.record("Shell 回收站调用失败：%v（项目未被删除）", err)
		return Result{Text: "放入回收站失败：" + err.Error() + "（项目未被删除）",
			IsError: true, AppliedSteps: se.snapshot()}, nil
	}
	se.record("通过 %s 将 %s 放入回收站，并已核对还原记录 %s（位于 %s）",
		detail.API, target, detail.Record, detail.Bin)
	if onUpdate != nil {
		onUpdate("fs.trash 已放入回收站")
	}
	return Result{
		Text: fmt.Sprintf("已放入回收站：%s（还原记录 %s，可在回收站还原）",
			target, detail.Record),
		AppliedSteps: se.snapshot(),
	}, nil
}

// ---------------------------------------------------------------------------
// fs.move (D34 row 5: L1, cross-volume L2)
// ---------------------------------------------------------------------------

type fsMove struct{ d FSDeps }

type fsMoveArgs struct {
	From string `json:"from"`
	To   string `json:"to"`
}

var fsMoveSchema = JSONSchema(`{"type":"object","properties":{` +
	`"from":{"type":"string","description":"源路径（经 C26 规范化）"},` +
	`"to":{"type":"string","description":"目标路径（经 C26 规范化）"}}` +
	`,"required":["from","to"],"additionalProperties":false}`)

func (fsMove) Name() string { return "fs.move" }
func (fsMove) Description() string {
	return "移动/重命名文件：同卷是一次可逆重命名（L1），跨卷是复制+永久删除源（L2）"
}
func (fsMove) Parameters() JSONSchema { return fsMoveSchema }

// Execute implements Tool.
func (t fsMove) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a fsMoveArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	if strings.TrimSpace(a.From) == "" || strings.TrimSpace(a.To) == "" {
		return Result{Text: "缺少 from 或 to 参数", IsError: true}, nil
	}
	from, err := t.d.canonical(a.From)
	if err != nil {
		return Result{Text: "源路径无法解析（按 fail-closed 拒绝）：" + err.Error(), IsError: true}, nil
	}
	to, err := t.d.canonical(a.To)
	if err != nil {
		return Result{Text: "目标路径无法解析（按 fail-closed 拒绝）：" + err.Error(), IsError: true}, nil
	}
	se := t.d.ledger()

	if !sameVolume(from, to) {
		return t.crossVolume(ctx, se, from, to, onUpdate), nil
	}
	if err := se.at(ctx, "rename"); err != nil {
		return se.stoppedResult("移动"), nil
	}
	if err := os.Rename(from, to); err != nil {
		se.record("同卷重命名失败：%v（源与目标均保持原状）", err)
		return Result{Text: "移动失败：" + err.Error() + "（源与目标均未被改动）",
			IsError: true, AppliedSteps: se.snapshot()}, nil
	}
	se.record("同卷重命名 %s → %s（再用一次反向重命名即可改回）", from, to)
	if onUpdate != nil {
		onUpdate("fs.move 已完成")
	}
	return Result{
		Text:         fmt.Sprintf("已移动（同卷重命名）：%s → %s", from, to),
		AppliedSteps: se.snapshot(),
	}, nil
}

// crossVolume is the L2 half of fs.move: copy to a staged file beside the
// destination, rename it into place, THEN remove the source. The source removal
// is permanent, which is why D34 puts this branch at L2 and why the ledger
// matters most here: a veto between the copy and the removal leaves TWO copies,
// and the report has to say that out loud.
func (fsMove) crossVolume(ctx context.Context, se *sideEffect, from, to string,
	onUpdate func(string)) Result {
	parent := dirOf(to)
	if parent == "" {
		return Result{Text: "拒绝移动到卷根目录：" + to, IsError: true}
	}
	src, err := os.Open(from)
	if err != nil {
		return Result{Text: "打开源文件失败：" + err.Error() + "（未做任何改动）", IsError: true}
	}
	defer src.Close() //nolint:errcheck // read-only handle

	tmp, err := os.CreateTemp(parent, tempPrefix+"*")
	if err != nil {
		return Result{Text: "创建目标临时文件失败：" + err.Error() + "（源未被改动）", IsError: true}
	}
	tmpName := tmp.Name()
	se.record("在 %s 创建目标临时文件 %s（源仍在 %s）", parent, baseOf(tmpName), from)

	written, werr := se.writeAll(ctx, tmp, src, onUpdate, 0)
	_ = tmp.Close()
	if werr != nil {
		se.record("复制未完成：%v", werr)
		if rmErr := os.Remove(tmpName); rmErr == nil {
			se.record("删除目标临时文件（源未被改动）")
		}
		return se.stoppedResult("跨卷移动")
	}
	se.record("已复制 %d 字节到目标临时文件", written)

	if err := se.at(ctx, "rename-into-place"); err != nil {
		se.record("目标尚未就位：临时文件 %s 仍存在，源未被改动", baseOf(tmpName))
		return se.stoppedResult("跨卷移动")
	}
	if err := os.Rename(tmpName, to); err != nil {
		se.record("重命名失败：%v（源未被改动）", err)
		if rmErr := os.Remove(tmpName); rmErr == nil {
			se.record("删除目标临时文件（源未被改动）")
		}
		return Result{Text: "跨卷移动失败：" + err.Error() + "（源未被改动）",
			IsError: true, AppliedSteps: se.snapshot()}
	}
	se.record("目标已就位：%s", to)

	if err := se.at(ctx, "remove-source"); err != nil {
		// THE honest case: the copy landed and the stop arrived afterwards.
		se.record("源仍保留在 %s：现在有两份内容，需要人工确认", from)
		return Result{
			Text: fmt.Sprintf("跨卷移动部分完成：%s 已就位；%s", to,
				"源未删除（两份内容并存）"),
			IsError:      true,
			AppliedSteps: se.snapshot(),
		}
	}
	if err := os.Remove(from); err != nil {
		se.record("删除源失败：%v（现在有两份内容，请人工确认）", err)
		return Result{Text: "目标已建立，但源删除失败：" + err.Error() + "（两份内容并存）",
			IsError: true, AppliedSteps: se.snapshot()}
	}
	se.record("已永久删除源 %s（跨卷移动的删除不进回收站）", from)
	if onUpdate != nil {
		onUpdate("fs.move 跨卷完成")
	}
	return Result{
		Text:         fmt.Sprintf("已跨卷移动：%s → %s（源已永久删除，不进回收站）", from, to),
		AppliedSteps: se.snapshot(),
	}
}

// ---------------------------------------------------------------------------
// fs.delete (D34 row 6: registered ONLY behind [fs] delete_enabled)
// ---------------------------------------------------------------------------

type fsDelete struct{ d FSDeps }

var fsDeleteSchema = JSONSchema(`{"type":"object","properties":{` +
	`"path":{"type":"string","description":"要永久删除的路径（经 C26 规范化）"}}` +
	`,"required":["path"],"additionalProperties":false}`)

func (fsDelete) Name() string { return "fs.delete" }
func (fsDelete) Description() string {
	return "永久删除一个文件（不进回收站）。仅当配置 [fs] delete_enabled=true 时才会注册"
}
func (fsDelete) Parameters() JSONSchema { return fsDeleteSchema }

// Execute implements Tool. The flag is re-checked HERE as well as at
// registration: a roster is built once, the setting can flip, and a tool that
// outlived its own authorization must still refuse.
func (t fsDelete) Execute(ctx context.Context, params json.RawMessage, _ func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	if !t.d.DeleteEnabled {
		return Result{Text: "fs.delete 未获授权（[fs] delete_enabled=false），已拒绝执行",
			IsError: true}, nil
	}
	var a fsPathArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	if strings.TrimSpace(a.Path) == "" {
		return Result{Text: "缺少 path 参数", IsError: true}, nil
	}
	target, err := t.d.canonical(a.Path)
	if err != nil {
		return Result{Text: "路径无法解析（按 fail-closed 拒绝）：" + err.Error(), IsError: true}, nil
	}
	se := t.d.ledger()
	if err := se.at(ctx, "remove"); err != nil {
		return se.stoppedResult("删除"), nil
	}
	if err := os.Remove(target); err != nil {
		se.record("删除失败：%v（目标未被改动）", err)
		return Result{Text: "删除失败：" + err.Error(), IsError: true,
			AppliedSteps: se.snapshot()}, nil
	}
	se.record("已永久删除 %s（不进入回收站，无法由 Wisp 还原）", target)
	return Result{
		Text:         "已永久删除：" + target + "（未进回收站；找回请用卷快照/备份）",
		AppliedSteps: se.snapshot(),
	}, nil
}

// ---------------------------------------------------------------------------
// risk facts (the R8/R2 input) + registration (D34)
// ---------------------------------------------------------------------------

// writeFacts is the R8 input for the tools that can replace content. Each raw
// path is canonicalized FIRST, because the existence question must be asked of
// the same path R2/R3 judged; "cannot tell" is reported as an overwrite, which
// is the expensive side of the gate.
func (d FSDeps) writeFacts(_ context.Context, _ map[string]any, paths []string) risk.Facts {
	f := risk.Facts{}
	for _, raw := range paths {
		c, err := d.canonical(raw)
		if err != nil {
			f.OverwriteExisting = true
			continue
		}
		ok, err := existsViaLstat(c)
		if err != nil || ok {
			f.OverwriteExisting = true
		}
	}
	return f
}

// moveFacts drives D34's move row: a same-volume rename keeps the L1 floor; an
// existing destination is an overwrite (R8) and a cross-volume move permanently
// destroys the source outside the bin (R8 "delete").
//
// It reads from/to out of the PARAMS, not the judged path list, because the
// bridge hands Facts a sorted set (pathArgs dedups and sorts), and a move is
// directional: which end is the source is exactly the question this answers.
func (d FSDeps) moveFacts(_ context.Context, params map[string]any, paths []string) risk.Facts {
	failClosed := risk.Facts{Irreversible: []string{"delete"}}
	fromRaw, toRaw := paramString(params, "from"), paramString(params, "to")
	if fromRaw == "" || toRaw == "" {
		return failClosed
	}
	from, err := d.canonical(fromRaw)
	if err != nil {
		return failClosed
	}
	to, err := d.canonical(toRaw)
	if err != nil {
		return failClosed
	}
	var f risk.Facts
	if !sameVolume(from, to) {
		f.Irreversible = []string{"delete"}
		return f
	}
	if ok, err := existsViaLstat(to); err != nil || ok {
		f.OverwriteExisting = true
	}
	return f
}

// paramString reads one string parameter without going through the bridge's
// sorted path set.
func paramString(params map[string]any, key string) string {
	if params == nil {
		return ""
	}
	s, _ := params[key].(string)
	return strings.TrimSpace(s)
}

// deleteFacts is R8's permanent-delete class, stated again on top of the L2
// declaration: the operation has no undo at all.
func deleteFacts(context.Context, map[string]any, []string) risk.Facts {
	return risk.Facts{Irreversible: []string{"delete"}}
}

// FSWriteDecl is the host-side declaration for fs.write. The declared level is
// L1 because a NEW file is recoverable by construction (nothing was replaced);
// an overwrite is raised to L2 by R8 through the Facts hook, never by a
// self-serving declaration.
func FSWriteDecl(d FSDeps) Decl {
	return Decl{
		Capabilities: []Capability{CapFSWrite},
		Needs:        []Capability{CapFSWrite},
		Declared:     risk.L1,
		PathParams:   []string{"path"},
		Facts:        d.writeFacts,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// FSTrashDecl is fs.trash: L1 by D34 BECAUSE the backend is the recycle bin. On
// a platform without one the tool refuses at Execute instead of keeping an L1
// promise it cannot honour.
func FSTrashDecl(FSDeps) Decl {
	return Decl{
		Capabilities: []Capability{CapFSWrite},
		Needs:        []Capability{CapFSWrite},
		Declared:     risk.L1,
		PathParams:   []string{"path"},
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// FSMoveDecl is fs.move: L1 same-volume, R8-raised to L2 for a cross-volume
// move or an existing destination.
func FSMoveDecl(d FSDeps) Decl {
	return Decl{
		Capabilities: []Capability{CapFSWrite},
		Needs:        []Capability{CapFSWrite},
		Declared:     risk.L1,
		PathParams:   []string{"from", "to"},
		Facts:        d.moveFacts,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// FSDeleteDecl is fs.delete: D34 pins the declared level at L2 and the R8 delete
// class says so again.
func FSDeleteDecl(FSDeps) Decl {
	return Decl{
		Capabilities: []Capability{CapFSWrite},
		Needs:        []Capability{CapFSWrite},
		Declared:     risk.L2,
		PathParams:   []string{"path"},
		Facts:        deleteFacts,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// BuiltinFSWriteEntries returns the D34 write half on its own, so a host (or a
// test) can assert the delete gating without rebuilding the L0 pair.
//
// fs.delete's ABSENCE when the flag is off is the enforcement: a tool the model
// cannot see cannot be called, which is a stronger guarantee than a call that
// would be refused. D34 asks for exactly that.
func BuiltinFSWriteEntries(d FSDeps) []Entry {
	out := []Entry{
		{Tool: fsWrite{d: d}, Decl: FSWriteDecl(d)},
		{Tool: fsTrash{d: d}, Decl: FSTrashDecl(d)},
		{Tool: fsMove{d: d}, Decl: FSMoveDecl(d)},
	}
	if d.DeleteEnabled {
		out = append(out, Entry{Tool: fsDelete{d: d}, Decl: FSDeleteDecl(d)})
	}
	return out
}
