package tools

// Ticket 221 甲形（批准＝台账 A434）：说明书不许对模型许诺一枚没注册的工具。
//
// 票面 AC#1 的判据（编排者在"收 221-c1"一节里改写过的分母，逐字照它写）：
//
//	分母 = 「这次装配真正注册出去的所有工具名的并集」＝包内测试枚举几家
//	       Builtin*Entries 构造器的结果；
//	分子 = 各 Description() / Parameters() 文本里出现的 task. 词根；
//	红句 = 要能指出是哪一句、缺哪一枚。
//
// ⚠ 分母绝不能写成单枚 BuiltinTaskEntries：task.spawn 走的是
// BuiltinSubagentEntries，那样这把尺会把 task.spawn 自己报成"许诺了没注册"
// ＝假红（票面 ④，编排者认账并已改写）。
//
// 这是一把**常驻能力尺**：它扫的是能力（注册名册），不是词面——将来谁在说明书
// 里许诺一枚新工具而没注册，它就响；把 task.cancel 摘掉注册，它也响。

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// taskRootRe pulls a task.* tool root out of text the model reads. It needs a
// letter after the dot, so a parameter named task_id cannot be mistaken for a
// tool root.
var taskRootRe = regexp.MustCompile(`task\.[a-z][a-z0-9_]*`)

// clauseMarks are the separators a reader of these Chinese descriptions actually
// uses. ASCII '.' is deliberately NOT one of them: every tool root contains a
// dot, so cutting on it would report 「可以用 task」 as the offending sentence
// and hide the fact being asserted.
var clauseMarks = []string{"；", "。", "\n"}

// allBuiltinEntriesHere is AC#1's denominator: the union of every builtin entry
// this configuration can register. It is built from the constructors themselves,
// so a re-fork between "what we promise" and "what we register" is what this leg
// measures, and no hand-kept list of names can drift out of sync with it.
func allBuiltinEntriesHere(t *testing.T) []Entry {
	t.Helper()
	paths := NewPathCanonicalizer(nil, nil)
	roster := NewTaskRoster()
	var out []Entry
	for _, batch := range [][]Entry{
		BuiltinFSEntries(FSDeps{Paths: paths, DeleteEnabled: true}),
		BuiltinFSWriteEntries(FSDeps{Paths: paths, DeleteEnabled: true}),
		BuiltinTaskEntries(TaskDeps{Roster: roster, Paths: paths}),
		BuiltinSubagentEntries(SubagentDeps{Roster: roster}),
	} {
		for _, e := range batch {
			if e.Tool == nil {
				t.Fatalf("%s 构造器交出一枚 nil Tool", batchKind(batch))
			}
			// fs.* families overlap by design (BuiltinFSEntries already carries
			// the write half), and a registry would refuse the duplicate - so
			// the union is keyed by name here, first writer wins.
			if !containsName(out, e.Tool.Name()) {
				out = append(out, e)
			}
		}
	}
	if len(out) == 0 {
		t.Fatal("一枚内置工具都没注册到，这把尺就是空转")
	}
	return out
}

func batchKind([]Entry) string { return "Builtin*Entries" }

func containsName(entries []Entry, name string) bool {
	for _, e := range entries {
		if e.Tool.Name() == name {
			return true
		}
	}
	return false
}

// Test221EveryPromisedTaskNameIsRegistered is AC#1's resident ruler.
//
// The un-fixed reading (measured by this leg BEFORE the production change):
// task.spawn's Description() promises 「可以用 task.cancel 单独停它」, while the
// registered union holds task.output and task.spawn from the task family - so
// exactly one root, task.cancel, is promised and not registered. That is the one
// and only red this leg had to see before landing 甲形.
func Test221EveryPromisedTaskNameIsRegistered(t *testing.T) {
	entries := allBuiltinEntriesHere(t)

	registered := map[string]bool{}
	for _, e := range entries {
		registered[e.Tool.Name()] = true
	}

	type promise struct {
		tool     string // whose text promised it
		slot     string // Description() or Parameters()
		sentence string // the clause the model actually reads
	}
	promised := map[string][]promise{}

	for _, e := range entries {
		for _, text := range []struct {
			slot string
			body string
		}{
			{"Description()", e.Tool.Description()},
			{"Parameters()", string(e.Tool.Parameters())},
		} {
			for _, hit := range taskRootRe.FindAllStringIndex(text.body, -1) {
				root := text.body[hit[0]:hit[1]]
				promised[root] = append(promised[root], promise{
					tool: e.Tool.Name(), slot: text.slot,
					sentence: sentenceAround(text.body, hit[0]),
				})
			}
		}
	}

	roots := make([]string, 0, len(promised))
	for root := range promised {
		roots = append(roots, root)
	}
	sort.Strings(roots)

	// The ruler checks its own射程 first: an empty numerator over a non-empty
	// denominator means it is scanning nothing at all, which must not read green.
	if len(roots) == 0 {
		t.Fatal("说明书里一枚 task.* 都没提：这把尺没在扫任何东西（分母非空、分子为空＝判据失效）")
	}

	for _, root := range roots {
		if registered[root] {
			continue
		}
		for _, p := range promised[root] {
			t.Errorf("说明书对模型许诺了一枚不存在的工具：%s 的 %s 写着「%s」，"+
				"而 %q 没有注册进这次装配的并集（并集 %d 枚：%s）",
				p.tool, p.slot, p.sentence, root, len(registered),
				strings.Join(sortedNames(registered), " / "))
		}
	}
}

func sortedNames(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for n := range set {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// sentenceAround returns the clause containing byte offset i, cut on clauseMarks.
// It clamps the result so a red line stays readable instead of dumping a whole
// parameter schema.
func sentenceAround(body string, i int) string {
	start := 0
	for _, mark := range clauseMarks {
		if j := strings.LastIndex(body[:i], mark); j >= 0 && j+len(mark) > start {
			start = j + len(mark)
		}
	}
	end := len(body)
	for _, mark := range clauseMarks {
		if j := strings.Index(body[i:], mark); j >= 0 && i+j < end {
			end = i + j
		}
	}
	clause := strings.TrimSpace(body[start:end])
	if r := []rune(clause); len(r) > 120 {
		clause = fmt.Sprintf("%s…", string(r[:120]))
	}
	return clause
}
