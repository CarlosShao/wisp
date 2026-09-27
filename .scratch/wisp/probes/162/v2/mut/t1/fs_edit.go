package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// The D34 row ticket 162 adds (docs/PLAN.md 的 D34 表，2026-09-27 owner 批准):
//
//	fs.edit  字面定位替换已存在文件里的一小段   L2 (R8 覆盖已存在)
//
// It exists because until now the only way to change one line was fs.write,
// which takes the WHOLE file back and overwrites with it: a model that re-copies
// and drops a block loses that block silently. Ticket 162 AC#1 measured the
// claim on unmodified code (internal/tools/fswrite_silentloss_ac1_test.go) -
// 6 lines, 222 bytes gone, and neither the Result, the tool_call row, the audit
// line nor a later fs.read said anything at all.
//
// fs.edit's contract is the opposite shape: the caller names the bytes it
// replaces, so there is nothing to re-copy and nothing to drop. The price is
// that locating is now a question we can answer with a hard error, and every
// refusal in this file happens BEFORE the first byte is written.
//
// What IS deliberately NOT here: the Unicode repair layer Pi runs as a second
// pass after an exact match fails (判据形状 1 的第二半) - an edit that only
// lands after a fuzzy repair is precisely the kind of silent content change this
// tool exists to prevent, so it needs its own decision, not a bolt-on. The
// risk-routing assertions (AC#5) are C19's and internal/risk's to own.
//
// AC#3b (r3) added the two refusals this file used to be silent about - an old
// that swallows the file's BOM, and a new that breaks the file's own line-ending
// convention. Those are NOT the normalization layer above: they refuse a shape,
// they never widen what matches.

// fsEdit is the D34 locate-and-replace tool.
type fsEdit struct{ d FSDeps }

// fsEditArgs is the argument object: one path and a batch of literal edits.
//
// OLD is a *string, not a string, so "the model left the field out" is told
// apart from "the model asked for an empty old" - the second one is the
// file-wide insertion this tool must refuse (判据形状 4), the first is a
// malformed call. Both refuse; only the wording differs, and the wording is the
// remedial action the model reads.
type fsEditArgs struct {
	Path  string     `json:"path"`
	Edits []fsEditOp `json:"edits"`
}

type fsEditOp struct {
	Old *string `json:"old"`
	New *string `json:"new"`
}

var fsEditSchema = JSONSchema(`{"type":"object","properties":{` +
	`"path":{"type":"string","description":"要改写的已存在文件（经 C26 规范化）"},` +
	`"edits":{"type":"array","minItems":1,"description":"一批字面替换，按顺序定位；任何一枚不成立则整批不落盘","items":{` +
	`"type":"object","properties":{` +
	`"old":{"type":"string","description":"要被替换的那一段原文，字面子串、非空、必须在文件里唯一命中"},` +
	`"new":{"type":"string","description":"替换成什么（可以比 old 长或短）"}}` +
	`,"required":["old","new"],"additionalProperties":false}}}` +
	`,"required":["path","edits"],"additionalProperties":false}`)

func (fsEdit) Name() string { return "fs.edit" }
func (fsEdit) Description() string {
	return "在已存在的文件里按字面子串定位并替换一小段（票 162）。old 必须逐字一致且唯一命中；" +
		"匹不上、命中多处、old 为空、old 与 new 相同或同批内区间重叠都是硬错，且整批一个字也不落盘。" +
		"写入走临时文件＋原子重命名（D31）。改一个文件里的一行请用本工具，不要用 fs.write 重抄整个文件"
}
func (fsEdit) Parameters() JSONSchema { return fsEditSchema }

// refusal is one refusal's shape: the finding, plus the action that fixes it.
//
// 判据形状 2 demands the remedial text be IN the error, not implied by it. The
// English sentences below are the upstream wording this ticket copies (票面
// 判据形状 2/3 quote them); they are not decoration - the model reads them and
// self-corrects (D37), so a paraphrase that drops "all whitespace and newlines"
// would lose the exact thing the model got wrong.
const (
	refusalExactWording = "old 必须与文件内容逐字一致，包含全部空白与换行" +
		"（the old string must match exactly including all whitespace and newlines）"
	refusalUniqueWording = "请在 old 里带上更多上下文，让它只命中一处" +
		"（provide more surrounding context to make it unique）"
)

func refusal(why, remedial string) Result {
	return Result{
		Text:    "fs.edit 拒绝执行：" + why + "。" + remedial + "。本次调用一个字也没有写，目标文件保持原状。",
		IsError: true,
	}
}

// Execute implements Tool. Every refusal below returns IsError with the ledger
// empty - the D31 writer is reached only after the whole batch has resolved, so
// "all or nothing" is a property of this order of operations, not a promise.
func (t fsEdit) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a fsEditArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	if strings.TrimSpace(a.Path) == "" {
		return Result{Text: "缺少 path 参数", IsError: true}, nil
	}
	if len(a.Edits) == 0 {
		return refusal("edits 是空的，没有任何要改的段落",
			"要新建文件请用 fs.write；要改这一段就把它放进 edits 里"), nil
	}
	target, err := t.d.canonical(a.Path)
	if err != nil {
		return Result{Text: "路径无法解析（按 fail-closed 拒绝）：" + err.Error(), IsError: true}, nil
	}
	limit := t.d.writeCap()

	// fs.edit only rewrites something that is already there. If the file is not
	// there, there is nothing to locate in it, and "create it empty and then
	// fail to find the text" is a worse answer than refusing.
	exists, err := existsViaLstat(target)
	if err != nil {
		return refusal("无法确认目标是否存在（存在性问题答不出来时按贵的那一侧走）："+err.Error(),
			"请确认路径可读，或改用 fs.read 先看一眼"), nil
	}
	if !exists {
		return refusal("目标文件不存在，没有可定位的内容："+target,
			"新建文件请用 fs.write；本工具只改已存在的文件"), nil
	}

	raw, err := os.ReadFile(target)
	if err != nil {
		return Result{Text: "读取目标失败：" + err.Error() + "（目标未被改动）", IsError: true, Origin: target}, nil
	}
	content := string(raw)
	if len(content) > limit {
		return refusal(fmt.Sprintf("目标文件 %d 字节，超过本工具单次处理上限 %d 字节", len(content), limit),
			"请缩小到单个文件，或改用宿主内部的 artifacts 通道"), nil
	}

	// 判据形状 1: locating is a LITERAL byte-string match, at exact precedence,
	// and nothing else. No line numbers, no diff format, no fuzzy pass.
	//
	// 判据形状 4's all-or-nothing needs every position resolved before anything
	// is applied, and a later old must be counted against the ORIGINAL bytes -
	// otherwise an earlier edit would quietly move the second one's match count
	// and the refusal would depend on the order the model happened to list them.
	type located struct {
		old, new    string
		start, stop int
	}
	matches := make([]located, 0, len(a.Edits))
	for i, op := range a.Edits {
		if op.Old == nil {
			return refusal(fmt.Sprintf("第 %d 枚编辑没有给 old 字段", i+1), refusalExactWording), nil
		}
		if op.New == nil {
			return refusal(fmt.Sprintf("第 %d 枚编辑没有给 new 字段（删除段落也要显式写成 new=\"\"）", i+1),
				"把要替换成什么写出来，即使答案是空字符串"), nil
		}
		old, updated := *op.Old, *op.New
		if old == "" {
			return refusal(fmt.Sprintf("第 %d 枚编辑的 old 是空字符串：那等于往文件里插一段而不改任何原文", i+1),
				"请把 old 写成要被替换掉的那一段原文（非空），整段插入请改成锚定到具体上下文"), nil
		}
		if old == updated {
			return refusal(fmt.Sprintf("第 %d 枚编辑的 old 与 new 相同：这一枚不改变任何东西", i+1),
				"请去掉这一枚，或把 new 改成真正要换成的内容"), nil
		}
		count := strings.Count(content, old)
		switch {
		case count == 0:
			return refusal(
				fmt.Sprintf("第 %d 枚编辑的 old 在文件里找不到（0 命中，文件 %d 字节）", i+1, len(content)),
				refusalExactWording+"；先用 fs.read 读回原文再照抄那一段"), nil
		case count > 1:
			return refusal(
				fmt.Sprintf("第 %d 枚编辑的 old 在文件里命中 %d 处（多于 1 处即为歧义，不取第一处）", i+1, count),
				refusalUniqueWording), nil
		}
		start := strings.Index(content, old)
		// AC#3b shape ①: an old that starts inside the file's UTF-8 BOM swallows
		// those bytes into the replaced region, and the file loses its encoding
		// marker without the caller ever naming it. Refused HERE, before any
		// other edit in the batch is applied and before the write, so the whole
		// call stays all-or-nothing.
		if n := leadingBOMLen(content); (n > 0 && start < n) || start == 0 { // TAMPER-V2-t1
			return refusal(
				fmt.Sprintf("第 %d 枚编辑的 old 从文件第 %d 字节起，圈住了文件头那 %d 字节 EF BB BF"+
					"（UTF-8 BOM，文件 %d 字节）", i+1, start, n, len(content)),
				"BOM 是文件自己的编码标记，不是内容：请把 old 和 new 都从 BOM 之后的第一个字节开始抄；"+
					refusalExactWording), nil
		}
		matches = append(matches, located{old: old, new: updated, start: start, stop: start + len(old)})
	}

	// 判据形状 4's other half: two edits landing on overlapping bytes cannot both
	// be honoured, and applying one while silently losing the other is the exact
	// failure mode this ticket exists to remove. The ranges are counted here, on
	// the resolved batch, and the whole call is refused - never "first one wins".
	for i := range matches {
		for j := i + 1; j < len(matches); j++ {
			if matches[i].start < matches[j].stop && matches[j].start < matches[i].stop {
				return refusal(fmt.Sprintf("第 %d 与第 %d 枚编辑定位到重叠的区间（[%d,%d) 与 [%d,%d)）",
					i+1, j+1, matches[i].start, matches[i].stop, matches[j].start, matches[j].stop),
					"请把重叠的那两枚合并成一枚，或改成互不相交的段落"), nil
			}
		}
	}

	// Apply in ascending byte order. The positions above all refer to the
	// ORIGINAL content, so each shift is the accumulated length change of the
	// edits already applied; the model's listing order must not change the bytes
	// that come out, which is why the sort happens before the splice.
	sort.Slice(matches, func(i, j int) bool { return matches[i].start < matches[j].start })
	var b strings.Builder
	next, shift := 0, 0
	for _, m := range matches {
		start, stop := m.start+shift, m.stop+shift
		b.WriteString(content[next:start])
		b.WriteString(m.new)
		next = stop
		shift += len(m.new) - len(m.old)
	}
	b.WriteString(content[next:])
	updated := b.String()

	// AC#3b shape ②: a new that spells its breaks in a style the file does not
	// use turns a pure-CRLF (or pure-LF, or pure-CR) file into a mixed-ending one
	// - today silently, at 3 CRLF pairs plus 1 bare LF.
	//
	// Why this refusal cannot weld the door: every old is matched BYTE-EXACTLY
	// against content, so when content is pure <one convention> no old taken from
	// it can contain a foreign break. Every foreign break in `updated` therefore
	// came out of a new the caller spelled, and only out of that new - an edit
	// that copies the file's own convention through keeps landing. Files that
	// have no convention (no breaks at all) or already mix are NOT touched by
	// this check: there is nothing there to preserve, and guessing one is how a
	// tool starts rewriting bytes nobody asked about.
	if conv := lineEndingConvention(content); conv != "" && foreignEndings(conv, updated) > 0 {
		who := "本批 new"
		crlf, cr, lf := countEndings(updated)
		for i, m := range matches {
			if foreignEndings(conv, m.new) > 0 {
				who = fmt.Sprintf("第 %d 枚编辑的 new", i+1)
				break
			}
		}
		return refusal(
			fmt.Sprintf("%s 会改写这份文件的行尾约定：文件原本是纯 %s，改完会变成 CRLF %d 组 / 裸 CR %d 条 / 裸 LF %d 条",
				who, endingName(conv), crlf, cr, lf),
			fmt.Sprintf("请把 new 里的换行按文件自己的行尾拼写（纯 %s 的文件：%s）；"+
				"确实要把整份文件换成别的行尾，请改用 fs.write 整文件写回，本工具不猜行尾。%s",
				endingName(conv), endingSpelling(conv), refusalExactWording)), nil
	}

	if len(updated) > limit {
		return refusal(fmt.Sprintf("改完的内容 %d 字节，超过单次上限 %d 字节", len(updated), limit),
			"请缩小改动范围，或改用宿主内部的 artifacts 通道"), nil
	}

	// The ONLY write in this function, and it is behind every check above:
	// D31's staged writer, temp file in the target's own directory then one
	// atomic rename (PLAN.md 的那条工程要求；票面 判据形状 6 明令不许照抄上游
	// 的"整文件交出去写")。
	res := t.d.stageAndRename(ctx, target, strings.NewReader(updated), onUpdate)
	res.Origin = target
	if !res.IsError {
		res.Text = fmt.Sprintf(
			"fs.edit 已改写 %s：%d 枚编辑全部生效，%d 字节 / %d 行 → %d 字节 / %d 行（临时文件+原子重命名）",
			target, len(matches), len(content), countNL(content), len(updated), countNL(updated))
	}
	return res, nil
}

// countNL is a display-only line count for the success line: this tool's own
// answer says how long the file is now, which is the number fs.write's answer
// has no counterpart for.
func countNL(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n")
}

// ---------------------------------------------------------------------------
// AC#3b's two shape checks
// ---------------------------------------------------------------------------

// leadingBOMLen is how many of a file's first bytes are a UTF-8 byte-order mark
// (EF BB BF), else 0. Spelled as bytes because the three bytes ARE the thing: a
// literal BOM in this source line would be invisible to whoever reads it.
func leadingBOMLen(s string) int {
	if strings.HasPrefix(s, "\xef\xbb\xbf") {
		return 3
	}
	return 0
}

// countEndings returns how one blob spells its breaks: (crlf pairs, bare CR,
// bare LF). A pure-CRLF blob has bareLF 0; a pure-LF blob has bareCR 0.
func countEndings(s string) (crlf, bareCR, bareLF int) {
	crlf = strings.Count(s, "\r\n")
	bareCR = strings.Count(s, "\r") - crlf
	bareLF = strings.Count(s, "\n") - crlf
	return crlf, bareCR, bareLF
}

// lineEndingConvention names the ONE kind of break a blob uses exclusively:
// "crlf", "lf" or "cr". It answers "" for a blob with no break at all and for
// one that already mixes - in both cases there is no convention to preserve, so
// the caller must not refuse anything on its behalf.
func lineEndingConvention(s string) string {
	crlf, cr, lf := countEndings(s)
	switch {
	case crlf > 0 && cr == 0 && lf == 0:
		return "crlf"
	case lf > 0 && cr == 0 && crlf == 0:
		return "lf"
	case cr > 0 && crlf == 0 && lf == 0:
		return "cr"
	}
	return ""
}

// foreignEndings counts the breaks in s that are NOT how conv's file spells
// them - the bytes this tool refuses to be the ones that introduced.
func foreignEndings(conv, s string) int {
	crlf, cr, lf := countEndings(s)
	switch conv {
	case "crlf":
		return cr + lf
	case "lf":
		return crlf + cr
	case "cr":
		return crlf + lf
	}
	return 0
}

// endingName / endingSpelling render a convention for the two sentences the
// model reads: what the file is, and how to spell a break in that file.
func endingName(conv string) string {
	switch conv {
	case "crlf":
		return "CRLF"
	case "lf":
		return "LF"
	case "cr":
		return "CR"
	}
	return conv
}

func endingSpelling(conv string) string {
	switch conv {
	case "crlf":
		return `换行写成 \r\n，回车加换行两字节`
	case "lf":
		return `换行写成 \n，单字节换行`
	case "cr":
		return `换行写成 \r，单字节回车`
	}
	return "按文件原有行尾拼写"
}

// ---------------------------------------------------------------------------
// registration (D34)
// ---------------------------------------------------------------------------

// FSEditDecl is the host-side declaration for fs.edit.
//
// The declared floor is L2, not L1: fs.edit can only act on a file that already
// exists (there is nothing to locate in an absent one), so "覆盖已存在内容" is
// not a fact some call happens to have - it is the tool's only mode. That is
// R8's conclusion stated at the place D34 froze it
// (docs/PLAN.md 的 D34 表 fs.edit 那一行)，and it is why no Facts hook is
// needed here to raise anything. Risk ROUTING (how the card reads, and the
// out-of-scope-path half) stays C19's: this ticket declares, it does not judge
// (AC#5 is r2's, and internal/risk is out of its approved scope).
func FSEditDecl(d FSDeps) Decl {
	return Decl{
		Capabilities: []Capability{CapFSWrite},
		Needs:        []Capability{CapFSWrite},
		Declared:     risk.L2,
		PathParams:   []string{"path"},
		Resident:     true,
		Provider:     KindBuiltin,
	}
}
