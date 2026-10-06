package tools

// Ticket 236 AC#2: the (b) branch of the DEFERRED ruler has to scan CAPABILITY,
// not word counts - and measured today it has no teeth.
//
// What this file puts on trial is one thing only: whether the instrument can
// tell a true claim from a false one. The claim itself (task.list is DEFERRED,
// PLAN.md §7 :1531 / D34; ticket 236 must not register it) is NOT re-decided
// here, and the marker <-> SPEC-12 §5 two-way reconciliation is ticket 225's,
// not this file's. 票面 AC#2 逐字：「225 管『标记与 SPEC-12 §5 双向对账』，
// 本格只管这把尺本身有没有牙」。
//
// ---------------------------------------------------------------------------
// The ruler on trial, and the two reasons its (b) branch is toothless
// ---------------------------------------------------------------------------
//
// Carrier being measured: task_cancel_221_legs_test.go:223-247,
// Test221DeferredMarkerForCancelLiftedButListStillMarked. It does
// os.ReadFile("task.go") (:224) and counts lines containing "DEFERRED",
// branching on whether the same line also names task.cancel (a) or task.list (b).
// The (b) branch is :244-246 (listMarked == 0 -> Errorf).
//
//  1. The count is propped by a narrative sentence, so it cannot fall to 0 for
//     the reason it claims. Measured at 交件 §1: `grep -n DEFERRED task.go` = 3
//     lines (:23 the roster header, :32 a sentence about that header, :279 the
//     Count() doc comment); 2 of them also carry "task.list" (:23 and :279).
//     Delete :23 - the very line the branch is supposed to guard - and
//     listMarked is still 1: the leg would report "标记还在" over a file whose
//     marker line is gone. 交件 §3 measures this shape twice: M-2a (the teeth-m13
//     copy, arithmetic in logs/r3c/mut/copies-and-landing-proofs.txt) and the
//     read-path positive control, where the archived leg is shown green over the
//     deleted-marker text and red only over the all-markers-gone copy.
//  2. A read-disk ruler is structurally blind to `go test -overlay`, which is
//     the instrument this ticket must mutate with. The overlay replaces the
//     bytes the COMPILER sees; os.ReadFile opens the PHYSICAL file at run time,
//     which no overlay ever touched. So a mutation that removes every
//     DEFERRED+task.list line is invisible to the leg guarding it: with this
//     carrier you cannot even MEASURE whether the ruler has teeth. 票面 §四
//     states this as 「今天连『测这把尺的牙』都做不到」. This is an instrument
//     fact, pinned here because AC#2 asks for it verbatim, and it is not a code
//     defect: the answer is not to make a disk read overlay-visible, it is to
//     stop resting a claim on a comment. Test236R3InstrumentFactReadDiskIs
//     BlindToOverlayCarriesIt below is the resident carrier that makes the
//     blindness measurable (it goes red on its OWN the moment task.go is
//     overlaid, while the read-disk leg stays green over the same mutation).
//
// The same blindness covers that leg's (a) branch, since both branches share one
// os.ReadFile carrier. This file does not adjudicate (a) (see 交件 §5 item 1).
//
// ---------------------------------------------------------------------------
// Where the readings in this file come from (attribution, AC#2's 交件 rule)
// ---------------------------------------------------------------------------
//
// No judgement below is copied from the earlier legs. 236-r3 and 236-r3b both
// died of a service fault before writing any adjudication (their evidence file
// was 36 lines of "打算答" intent, commit d9aff5fd; the orchestrator re-submitted
// it as 未验证半成品). This file's own bytes are leg 236-r3c: every rc, every
// FAIL line and every roster quoted in 交件 §2/§3 was produced by r3c on
// 10-06 16:2x-16:4x from the copies under D:/tmp/wisp236r3c/, and the archived
// probe files under .scratch/wisp/probes/236/r3/logs/** were used as a SHAPE
// reference only - never as a credential.
//
// ---------------------------------------------------------------------------
// The AC#2 nail: ask the roster, not the comment
// ---------------------------------------------------------------------------
//
// "task.list is DEFERRED" is a claim about capability: the family this host
// registers must not contain that row. So the nail reads BuiltinTaskEntries'
// registered names - the same denominator 票面 §四 points at, because
// task_cancel_221_test.go's allBuiltinEntriesHere and
// ticket175r2_stamp_live_test.go's classification table are already
// capability rulers (⛔ nothing new is invented, nothing existing is deleted or
// loosened: the lexical leg keeps both branches and the count-shaped assertion at
// task_cancel_221_legs_test.go:183 keeps counting).
//
// A capability nail is reachable by the mutation instrument because the roster is
// COMPILED. 交件 §3 records which shapes leg 236-r3c measured:
//
//	M-2a (teeth-m13)  delete the marker row at :23 only -> this file's roster nail
//	      stays GREEN, which is correct: the capability did not change. The leg
//	      that cannot tell this apart from a real code change is the archived
//	      lexical one, and it stays green for the wrong reason (see 1 above).
//	M-2b  register a task.list row -> the nail goes red BY NAME, while the
//	      archived lexical leg stays green (:23 is still on disk): the exact
//	      divergence the old shape cannot see. Measured, 交件 §3.
//	M-2e  unregister task.cancel -> the nail goes red on the second half (the
//	      roster must contain the two rows that have implementations). Measured.
//	M-2c  register an unrelated new task.* row -> this nail is designed to stay
//	      green (it is a name set, not a count) while the count-shaped leg at
//	      task_cancel_221_legs_test.go:183 goes red. NOT measured by r3c: the
//	      dispatch capped this leg at four mutations, so the 假红 prediction of
//	      票面 §四 stays an argument from source, not a reading. 交件 §5 says so.

import (
	_ "embed"
	"os"
	"strings"
	"testing"
)

// taskSrcPath236r3 is the path a read-disk ruler opens. Kept as a constant so
// the overlay positive control (交件 §3 M-2p) can rewrite just this line.
const taskSrcPath236r3 = "task.go"

// compiled236r3TaskSrc is task.go as the COMPILER saw it. //go:embed is resolved
// at build time, so `-overlay` DOES reach it - the paired half of the instrument
// fact above: this reading and the disk reading can only diverge under overlay.
//
//go:embed task.go
var compiled236r3TaskSrc string

// deferredListMarkerLines236r3 counts the marker shape the archived ruler
// guards: a line carrying both tokens. Same two tokens, same predicate, kept as
// a pure function so 交件 §3's readings are auditable line by line.
func deferredListMarkerLines236r3(text string) int {
	var n int
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "DEFERRED") && strings.Contains(line, "task.list") {
			n++
		}
	}
	return n
}

// registeredTaskNames236r3 is the AC#2 denominator: the names the task family
// actually registers. Refuses to read green over an empty denominator.
func registeredTaskNames236r3(t *testing.T) []string {
	t.Helper()
	names := make([]string, 0, 4)
	for _, e := range BuiltinTaskEntries(TaskDeps{Roster: NewTaskRoster()}) {
		if e.Tool == nil {
			t.Fatal("BuiltinTaskEntries 交出一枚 nil Tool：名册读不出来，这把尺就是空转")
		}
		if strings.TrimSpace(e.Tool.Name()) == "" {
			t.Fatal("BuiltinTaskEntries 交出一枚无名工具：名册读不出来，这把尺就是空转")
		}
		names = append(names, e.Tool.Name())
	}
	if len(names) == 0 {
		t.Fatal("一枚 task 家族工具都没注册到：这把尺没在扫任何东西（分母为空＝判据失效）")
	}
	return names
}

// Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment is the AC#2 nail.
// Two name-set claims, both about capability, both reachable by -overlay:
// task.list must NOT be in the roster, and the two rows that have
// implementations must be. Deliberately NOT a count, so a legitimately
// approved third task tool cannot false-red (交件 §3 M-2c measures that pair).
func Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment(t *testing.T) {
	names := registeredTaskNames236r3(t)

	for _, n := range names {
		if n == "task.list" {
			t.Errorf("DEFERRED 那一支和名册分叉了：BuiltinTaskEntries 的注册名册里出现了 task.list 这一行"+
				"（现名册：%s）。§7 :1531 与 D34 都还把它记成 DEFERRED，「在册＋无实现＋无人认领」正是票 164 "+
				"AC#1 要杀的那一形；真要解冻得由人工批准、标记与接线同批核销（A434 批准范围 item 7），"+
				"而不是留下一行注释替一枚已注册的工具说它还不可用。"+
				"这一枚就是 AC#2 要的牙：读注释那把尺在同一发突变下照旧全绿（交件 §3 M-2b）。",
				strings.Join(names, " / "))
		}
	}

	var sawOutput, sawCancel bool
	for _, n := range names {
		switch n {
		case "task.output":
			sawOutput = true
		case "task.cancel":
			sawCancel = true
		}
	}
	if !sawOutput || !sawCancel {
		t.Errorf("task 家族名册 = %s, want 含 task.output 与 task.cancel（task.output 是票 164 的，"+
			"task.cancel 是票 221 甲形接的线）——名册少一枚就是「说明书与能力分叉」的另一形",
			strings.Join(names, " / "))
	}
	t.Logf("现名册：%s（DEFERRED 的判据读的是这一行，不是注释里的词面）", strings.Join(names, " / "))
}

// Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt is the carrier for the
// instrument fact AC#2 asks to be written into the judgement, and it is NOT a
// nail about task.list: it owns one invariant that is true at HEAD - the bytes
// the compiler built this package from and the bytes os.ReadFile opens are the
// same bytes. Two readings of the marker count, compared:
//
//	diskCount     = deferredListMarkerLines236r3(os.ReadFile(taskSrcPath236r3))
//	compiledCount = deferredListMarkerLines236r3(compiled236r3TaskSrc)
//
// They can only separate under `-overlay`, which is what makes this leg the
// overlay-visible witness for mutations of task.go: 交件 §3 M-2a runs as a PAIR -
// this leg red + the archived read-disk leg green - and that pair is the proof
// the blindness is structural rather than "the overlay failed to land".
// The positive control M-2p (overlay rewriting just the read path inside the
// archived test file) is what turns the archived leg red at all.
//
// ⚠ Consequence for whoever runs the next mutation on task.go: this leg goes red
// under ANY -overlay of that file, benign ones included. That is the reading it
// exists to produce, not a defect and not a new always-red - at HEAD, and in CI,
// the two readings are equal and the leg passes. Do not "fix" it by comparing
// less: comparing less is what made the archived leg unmeasurable.
func Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt(t *testing.T) {
	onDisk, err := os.ReadFile(taskSrcPath236r3)
	if err != nil {
		t.Fatalf("读不到 %s（读盘型尺的落点就是包内源文件）：%v", taskSrcPath236r3, err)
	}
	diskText := string(onDisk)
	diskCount := deferredListMarkerLines236r3(diskText)
	compiledCount := deferredListMarkerLines236r3(compiled236r3TaskSrc)

	if compiled236r3TaskSrc != diskText {
		t.Errorf("编译期与运行期读到的 task.go 字节不一致（编译 %d 字节 / 盘上 %d 字节，标记行数 %d 对 %d）："+
			"只有 -overlay 能让它们分开。这一发就是那条仪器事实的落地证明：盘上那把尺看不见突变，"+
			"所以它绿的时候不能当作「标记还在」的证据。标记账号＝票 236 AC#2，未修码读数见交件 §3 M-2a",
			len(compiled236r3TaskSrc), len(diskText), compiledCount, diskCount)
	}
	if diskCount != compiledCount {
		t.Errorf("同一枚 DEFERRED＋task.list 标记行数，编译器看到 %d 行、os.ReadFile 看到 %d 行："+
			"两个数只有在 -overlay 下才会分开（交件 §3 M-2a 的配对读数）",
			compiledCount, diskCount)
	}
	t.Logf("成对读数：盘上 %d 行 / 编译期 %d 行（HEAD 上两数相等＝这把尺没在空转）", diskCount, compiledCount)
}
