package tools

// Ticket 271, shape 甲 (A648 reopen): the count tripwire at
// task_cancel_221_legs_test.go:183-185 is upgraded from an undocumented
// alarm into a DOCUMENTED one, and the false-red shot it was accused of
// becomes a resident measurement instead of an argument from source.
//
// ---------------------------------------------------------------------------
// The shape pick (AC#0), compressed
// ---------------------------------------------------------------------------
//
// The ticket laid out three shapes:
//
//	甲  keep the count assertion, document it as a deliberate tripwire,
//	    and pin a positive control that proves it fires (this file).
//	乙  replace the count with a name-set judgment at :184.
//	丙  do nothing; register the false red as a known risk.
//
// Pick = 甲. The reason that decides it: A648 ruled 丙 on 10-06 only because
// every red-count in the three cost tables was transcribed, never measured
// ("表里几枚红那一栏全部是转述"), and left an explicit reopen condition --
// any leg about to wire a further task tool re-opens this ticket, which then
// auto-promotes to 甲 with the landing leg supplying the un-run numbers.
// The reopen ran twice and both readings are now resident: ticket 47's
// task.list took the family from 2 to 3 (this file's mutant expectations
// moved with that landing, same commit), and this positive control keeps
// converting the one load-bearing un-run reading -- 236r3's archived M-2c
// gap, recorded in the header of tasklist_deferred_236r3_teeth_test.go as
// "NOT measured by r3c ..." -- into something the suite itself re-proves on
// every run:
//
//	B-3 shot (register an UNRELATED task.note):
//	  - the count tripwire Test221TaskCancelIsRegisteredAtItsFrozenLevel
//	    goes red, and the red names the tripwire message (4 registered) --
//	    the alarm is wired, not dead;
//	  - the capability nail Test236R3TaskListDeferralIsPinnedByTheRoster
//	    NotByAComment stays GREEN on the same mutant -- the red is FALSE
//	    in the capability dimension (task.note is not task.list, the one
//	    name the nail demands), exactly the false-red the ticket names,
//	    now measured instead of transcribed;
//	  - the shared tree is never edited: the mutant exists only under
//	    -overlay against a temp copy, and the restore proof below shows
//	    disk bytes == compiled bytes and the real roster back at baseline.
//
// Why not 乙: the name-set ruler ALREADY exists (236r3's nail is one, and it
// is deliberately "NOT a count"). Converting :184 would duplicate that
// coverage while deleting the one dimension no other ruler carries -- the
// census (271-a1 §1) counted the family-size sensors in the whole repo and
// found exactly this one; every capability ruler only knows task.list by
// name, so a wrong registration under ANY other name is invisible to them.
// 乙 trades a unique sensor for a duplicate, and it edits a done ticket's
// predicate without a written unfreeze. Why not 丙: 丙 keeps the false red
// undocumented and the M-2c reading un-measured forever, which is the exact
// "checked the box but never measured" the orchestrator refused to sign.
//
// What the next leg sees on legal-third day (the cost 甲 accepts, in
// writing, per AC#2 -- this is NOT permission to widen the assertion):
//
//	1. the tripwire at :184 goes red on your legal third tool too -- that
//	   red is BY DESIGN; read this header and the comment above the test
//	   in task_cancel_221_legs_test.go;
//	2. the procedure: confirm your third tool's ticket is real, then take
//	   the shape decision back to ticket 271 / the ledger (A648 records
//	   the standing procedure) -- do NOT widen `!= 2` in place to go green;
//	3. the capability nail stays green for you either way -- it is the
//	   false-red-free half of the pair, and this file proves both halves
//	   on one mutant.
//
// Known blind spot, unchanged by any shape (census 271-a1 §4): no ruler in
// this repo can check APPROVAL -- legality of a third tool is process, not
// code. That dimension stays with the ledger.
//
// ---------------------------------------------------------------------------
// Mechanics
// ---------------------------------------------------------------------------
//
// The mutant is byte surgery on the //go:embed copy of task.go: the roster
// block gains one entry whose Tool only carries a Name ("task.note") -- the
// same shape 236-v1's B-3 probe used. The rulers under trial read
// BuiltinTaskEntries' names and never Execute anything, so a stub body is
// honest here; the appended type compiles inside the overlaid task.go
// without new imports (context and encoding/json are already imported
// there). The inner run targets exactly the two tests on trial; everything
// else in the suite is out of scope for this leg.

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// taskGo271 is task.go as the COMPILER sees it -- the same bytes the rulers
// under trial were built from. The mutant is derived from this copy, never
// from a re-read of disk, so the surgery is anchored to what actually
// compiles at HEAD.
//
//go:embed task.go
var taskGo271 string

// task271RosterAnchor is the exact line the mutant appends its third entry
// after. Count-checked before surgery: if the roster block moves or is
// reworded, this leg refuses to run instead of silently overlaying nothing
// (a silent no-op overlay would fake the positive control green).
const task271RosterAnchor = "{Tool: taskCancel{d: d}, Decl: taskCancelDecl()},"

// task271NoteStub is appended to the mutant task.go. It mirrors the shape
// 236-v1's B-3 used for the false-red reading: Name() only, never executed.
const task271NoteStub = `
// task271Note is ticket 271's M-2c mutant: an unrelated third registration.
// It exists only under -overlay produced by taskfamily_roster_271_tripwire_test.go
// and is never executed by any ruler under trial.
type task271Note struct{ d TaskDeps }

func (task271Note) Name() string { return "task.note" }
func (task271Note) Description() string { return "ticket 271 M-2c mutant: unrelated third registration, never executed" }
func (task271Note) Parameters() JSONSchema { return JSONSchema("{}") }
func (task271Note) Execute(ctx context.Context, params json.RawMessage, onUpdate func(delta string)) (Result, error) {
	return Result{Text: "task.note is a registration mutant and must never be called"}, nil
}
`

// write271UnrelatedThirdMutant overlays task.go with a copy whose roster
// carries the unrelated third entry, and returns the overlay JSON path.
// Everything lands in a temp dir outside the repo, so the shared working
// tree is never edited in place (AC#1's overlay rule).
func write271UnrelatedThirdMutant(t *testing.T) string {
	t.Helper()
	if strings.Count(taskGo271, task271RosterAnchor) != 1 {
		t.Fatalf("task.go 的名册锚点命中 %d 次（want 恰 1 次）：task.go 的 BuiltinTaskEntries 已漂移，"+
			"这枚正控拒绝在未知锚点上开刀（盲目 overlay 一份没变的拷贝会把正控假绿）",
			strings.Count(taskGo271, task271RosterAnchor))
	}
	mutated := strings.Replace(taskGo271, task271RosterAnchor,
		task271RosterAnchor+"\n\t\t{Tool: task271Note{d: d}, Decl: TaskOutputDecl()},", 1)
	mutated += task271NoteStub

	dir := t.TempDir()
	mutantPath := filepath.Join(dir, "task.go")
	if err := os.WriteFile(mutantPath, []byte(mutated), 0o644); err != nil {
		t.Fatalf("写不出 mutant task.go：%v", err)
	}
	// The overlay key must be the ABSOLUTE path go compiles task.go under.
	// os.Getwd + filepath.Join is the pattern every bridge/fs test here uses
	// for read-only path naming; no filepath.Abs (d22scan ban #2), and no
	// filesystem decision is being made -- this only names the file to
	// overlay inside a temp dir go test reads.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("取不到包目录（overlay 键要绝对路径）：%v", err)
	}
	abs := filepath.Join(wd, "task.go")
	overlay := struct {
		Replace map[string]string `json:"Replace"`
	}{Replace: map[string]string{
		filepath.ToSlash(abs): filepath.ToSlash(mutantPath),
	}}
	overlayPath := filepath.Join(dir, "overlay.json")
	raw, err := json.Marshal(overlay)
	if err != nil {
		t.Fatalf("写不出 overlay JSON：%v", err)
	}
	if err := os.WriteFile(overlayPath, raw, 0o644); err != nil {
		t.Fatalf("写不出 overlay JSON 文件：%v", err)
	}
	t.Logf("271 overlay JSON：%s", raw)
	return overlayPath
}

// run271InnerGoTest runs one filtered go test against the overlay and
// returns the combined output. Skipped only when no go toolchain is on PATH
// -- the skip names itself rather than faking a pass.
func run271InnerGoTest(t *testing.T, overlayPath, runPattern string) (string, bool) {
	t.Helper()
	if _, err := exec.LookPath("go"); err != nil {
		t.Skipf("PATH 上找不到 go 工具链，-overlay 内层跑不了（跳过原因具名，不算通过也不算失败）：%v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "test", ".", "-run", runPattern, "-count=1", "-v", "-overlay", overlayPath)
	out, err := cmd.CombinedOutput()
	text := string(out)
	if ctx.Err() != nil {
		t.Fatalf("内层 go test 超时（3 分钟）：%s", text)
	}
	return text, err == nil
}

// Test271B3UnrelatedThirdToolTripsOnlyTheCountTripwire is the M-2c false-red
// measurement, resident. One mutant, two rulers, both halves of the claim:
//
//	half 1 (alarm wired):  the count tripwire FAILS, and the failure text is
//	                       the tripwire's own (naming the count and the new
//	                       name) -- so the red is attributable to :184, not
//	                       to collateral anywhere else;
//	half 2 (red is false): the capability nail PASSES on the same mutant --
//	                       a name-set ruler correctly treats an unrelated
//	                       third tool as no capability change.
//
// If half 1 ever comes back green, the tripwire died (someone widened the
// count or gated it) and this leg says so by name. If half 2 ever comes back
// red on this mutant, the capability nail grew a count-shaped tooth and has
// started false-reding itself.
func Test271B3UnrelatedThirdToolTripsOnlyTheCountTripwire(t *testing.T) {
	overlayPath := write271UnrelatedThirdMutant(t)

	// half 1: the tripwire fires, and the red is ITS text.
	out, ok := run271InnerGoTest(t, overlayPath, "^Test221TaskCancelIsRegisteredAtItsFrozenLevel$")
	if ok {
		t.Fatalf("B-3 台件（注册无关 task.note）下枚数哨兵居然绿了：task_cancel_221_legs_test.go:183 的"+
			"枚数断言已经死掉（被放宽或被门掉），这枚正控就是来抓这件事的。内层输出：\n%s", out)
	}
	if !strings.Contains(out, "FAIL: Test221TaskCancelIsRegisteredAtItsFrozenLevel") {
		t.Fatalf("内层红了，但红的不是枚数哨兵那一枚（红错了地方，读数不可归属）。内层输出：\n%s", out)
	}
	if !strings.Contains(out, "task 家族注册了 4 枚") || !strings.Contains(out, "task.note") {
		t.Fatalf("枚数哨兵红了，但红句没有报出 4 枚与 task.note（哨兵文案被改，红得不可读）。内层输出：\n%s", out)
	}

	// half 2: the same mutant leaves the capability nail green.
	out2, ok2 := run271InnerGoTest(t, overlayPath, "^Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment$")
	if !ok2 {
		t.Fatalf("同一 B-3 台件下 236r3 的名集钉竟然红了：名集尺长出了枚数牙、开始自己也假红，"+
			"这正是本票不许发生的另一半。内层输出：\n%s", out2)
	}
	if !strings.Contains(out2, "task.note") {
		t.Fatalf("名集钉绿了，但它的 Logf 没报出 task.note（现名册没真的变到 3 枚，读数存疑）。内层输出：\n%s", out2)
	}
	t.Logf("B-3 台件 resident 读数：枚数哨兵红（假红，红句具名）＋名集钉绿——两把尺的分工在同一发突变下实测成立")
}

// Test271OverlayOnlyDisciplineRosterBackAtBaseline is AC#1's restore proof,
// resident: after the mutant runs, the shared tree must show zero edits and
// the roster the COMPILED binary sees must be back at its baseline shape.
//
//	(a) disk task.go bytes == the bytes this binary was built from -- the
//	    mutant only ever existed as an -overlay copy in a temp dir, so any
//	    difference here means someone edited task.go in place;
//	(b) the real capability nail passes on the un-overlaid roster -- invoked
//	    directly rather than re-stated, so this proof cannot drift away from
//	    the ruler it restores;
//	(c) the baseline count is LOGGED, not asserted: this leg deliberately
//	    creates no new count-shaped red of its own (AC#2).
func Test271OverlayOnlyDisciplineRosterBackAtBaseline(t *testing.T) {
	onDisk, err := os.ReadFile("task.go")
	if err != nil {
		t.Fatalf("读不到 task.go：%v", err)
	}
	if string(onDisk) != taskGo271 {
		t.Fatalf("盘上 task.go 与编译进本测试的字节不一致（差 %d 字节）：restore 的前提是 mutant 只活在"+
			"-overlay 里；出现差异说明有人对共享树原地动刀，还原自证即失败",
			len([]byte(taskGo271))-len(onDisk))
	}

	// (b) run the real ruler, not a replica of it.
	Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment(t)

	names := registeredTaskNames236r3(t)
	t.Logf("起手读数：task 家族名册 %d 枚（%s）——此行只记录、不断言；枚数警报的全部辖区在"+
		" task_cancel_221_legs_test.go:183 的哨兵与本文件的 B-3 正控，本测试不新造第二把枚数红",
		len(names), strings.Join(names, " / "))
}
