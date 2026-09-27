package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

// Ticket 174 AC#2, shape 丙: the pointer in a task.output stub stays, and the
// reply says out loud - in the text the model actually reads - when that
// pointer cannot be followed right now.
//
// Two shapes, judged by two different facts, plus the arm that has neither:
//
//	(a) authorization: the copy file is outside every authorized root, so
//	    fs.read on it lands at R2/L2. Judged ONLY through C26 (the same
//	    *PathCanonicalizer fs.read is judged by). Nothing here cleans or
//	    absolutizes a path by hand - AGENTS §1.2 bans that, and a hand-rolled
//	    check would drift from what fs.read really does.
//	(b) existence: the copy file is not there at all, or is a directory. The
//	    pre-fix branch tested only `ArtifactPath != ""`, so both spelled
//	    "全文见 …" and said nothing else (probe 174-c1 §5, re-measured here).
//	(+ ) no judge wired: fail CLOSED and say so - never assume readable.
//
// The reverse criterion this file exists to protect is the SILENCE arm:
// TestHealthyPointerStaysSilent. A reply that warned about every pointer would
// not be honest, it would be noise, and the ticket's fix would then be free to
// pass its own forward legs while quietly making the tool untrustworthy.
//
// Injection surface (AGENTS §1.3): the real bridge, the real C26
// canonicalizer, the real fs.read, real files on disk. No mock stands in for a
// real component here, and no authorization is widened to make an arm pass -
// [fs] allowed_dirs is configured only in the arms that model a host that was
// already configured.
//
// This file references no test in task_output_leg_test.go /
// task_output_ac2_before_test.go: those carry ticket 164's already-accepted
// criteria and are not this ticket's to edit.

// judgedTaskBridge wires fs.* AND task.* over ONE canonicalizer - the wiring
// this ticket adds. TaskDeps.Paths nil is a leg of its own (the fail-closed
// arm), not a default anything relies on.
func judgedTaskBridge(t *testing.T, roster *TaskRoster, paths *PathCanonicalizer) *Bridge {
	t.Helper()
	reg := NewRegistry()
	for _, e := range append(BuiltinFSEntries(FSDeps{Paths: paths}),
		BuiltinTaskEntries(TaskDeps{Roster: roster, Paths: paths})...) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	return New(Options{
		Registry: reg, Paths: paths,
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil}),
		Gate:       NoGate{}, Logf: func(string, ...any) {},
	})
}

func callTaskOutput(t *testing.T, b *Bridge, id string) agent.ToolOutcome {
	t.Helper()
	out, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "174-r1", CallID: "c-" + id, Name: "task.output",
		Args: mustArgs(t, map[string]any{"task_id": id}),
	})
	if err != nil {
		t.Fatalf("task.output must answer with data, not a Go error: %v", err)
	}
	return out
}

// readThroughBridge is the fact behind the wording: what fs.read actually does
// with the path the stub just named.
func readThroughBridge(t *testing.T, b *Bridge, path string) agent.ToolOutcome {
	t.Helper()
	out, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "174-r1", CallID: "c-read", Name: "fs.read",
		Args: mustArgs(t, map[string]any{"path": filepath.ToSlash(path)}),
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func mustWriteFile(t *testing.T, p, body string) string {
	t.Helper()
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

const noticeLead = "注意：" // what a speaking reply is made of

// TestPointerOutsideAuthorizedRootSpeaks is shape (a): a REAL copy file, on
// disk, readable by the OS - and unreachable anyway, because no root authorizes
// it. That distinction is the whole point: only C26 can answer it, so the fix
// may not be a stat.
func TestPointerOutsideAuthorizedRootSpeaks(t *testing.T) {
	dir := tempCanonical(t)
	body := asciiRun(20000)
	realFile := mustWriteFile(t, filepath.Join(dir, "tool-output-174-outside.txt"), body)

	roster := NewTaskRoster()
	roster.Record("bg-a", TaskOutput{Text: body, ArtifactPath: mustCanonical(t, realFile)})

	// An EMPTY allowlist is the default config (SPEC-03 allowed_dirs[]=[]), not
	// a widened one: this arm authorizes nothing at all.
	b := judgedTaskBridge(t, roster, NewPathCanonicalizer(nil, nil))
	out := callTaskOutput(t, b, "bg-a")
	t.Logf("shape (a) verbatim: IsError=%v Truncated=%v Text=%q", out.IsError, out.Truncated, announceOf(out.Text))

	if out.IsError || !out.Truncated {
		t.Fatalf("a recorded long output must still be answered as a truncation: %+v", out)
	}
	if !strings.Contains(out.Text, noticeLead) || !strings.Contains(out.Text, "读不到") {
		t.Fatalf("the reply must say the pointer cannot be followed, got %q", announceOf(out.Text))
	}
	if !strings.Contains(out.Text, "不在你被授权的目录范围内") {
		t.Errorf("shape (a) must name authorization as the reason, got %q", announceOf(out.Text))
	}
	// And the pointer itself stays (PLAN.md:2564; 只截不指＝不合格).
	if pointerRe.FindStringSubmatch(out.Text) == nil {
		t.Fatalf("the pointer must still be given, got %q", announceOf(out.Text))
	}
	// The wording must match the verdict fs.read really returns: same bridge,
	// same canonicalizer, real call.
	rd := readThroughBridge(t, b, realFile)
	if !rd.IsError || rd.RiskLevel != "L2" {
		t.Fatalf("expected the refusal the notice describes, got %+v", rd)
	}
	if rd.ErrorClass != "user_rejected" {
		t.Errorf("error_class = %q, want user_rejected (R2 outside the authorized dirs)", rd.ErrorClass)
	}
}

// TestPointerToMissingCopyFileSpeaks is shape (b) first arm: the host filed a
// path that is not there. Pre-fix this took the "有指针" branch and stayed
// silent, so "后半段不可找回" never fired either - the stub pointed at air.
func TestPointerToMissingCopyFileSpeaks(t *testing.T) {
	dir := tempCanonical(t)
	body := asciiRun(20000)
	ghost := filepath.Join(dir, "tool-output-does-not-exist.txt")
	if _, err := os.Stat(ghost); !os.IsNotExist(err) {
		t.Fatalf("the fixture must be absent for this leg to mean anything: %v", err)
	}

	roster := NewTaskRoster()
	roster.Record("bg-ghost", TaskOutput{Text: body, ArtifactPath: mustCanonical(t, ghost)})
	// Authorization is fine here (dir IS the root): only existence is broken,
	// which is what keeps shape (b) from being shape (a) in disguise.
	b := judgedTaskBridge(t, roster, NewPathCanonicalizer([]string{dir}, nil))

	out := callTaskOutput(t, b, "bg-ghost")
	t.Logf("shape (b) missing file verbatim: IsError=%v Truncated=%v Text=%q",
		out.IsError, out.Truncated, announceOf(out.Text))
	if !strings.Contains(out.Text, "并不存在") || !strings.Contains(out.Text, "读不到") {
		t.Fatalf("a pointer at a missing copy file must say so, got %q", announceOf(out.Text))
	}
	if pointerRe.FindStringSubmatch(out.Text) == nil {
		t.Fatalf("the pointer is still owed to the model, got %q", announceOf(out.Text))
	}
	if strings.Contains(out.Text, "不在你被授权的目录范围内") {
		t.Errorf("this path IS authorized; the notice must not claim otherwise: %q", announceOf(out.Text))
	}
}

// TestPointerToNonRegularPathSpeaks is shape (b) second arm: the name exists,
// but it is a directory. This is the shape ticket 164's own
// TestTruncationShapeIsTheD15Triple sits on, so the silence was never an
// accident of a bad fixture - it is the branch.
func TestPointerToNonRegularPathSpeaks(t *testing.T) {
	dir := tempCanonical(t)
	body := asciiRun(20000)
	roster := NewTaskRoster()
	roster.Record("bg-dir", TaskOutput{Text: body, ArtifactPath: dir})
	b := judgedTaskBridge(t, roster, NewPathCanonicalizer([]string{dir}, nil))

	out := callTaskOutput(t, b, "bg-dir")
	t.Logf("shape (b) directory verbatim: IsError=%v Truncated=%v Text=%q",
		out.IsError, out.Truncated, announceOf(out.Text))
	if out.IsError {
		t.Fatalf("the tool still answers, it does not refuse the whole read: %+v", out)
	}
	if !strings.Contains(out.Text, "不是一般文件") || !strings.Contains(out.Text, "读不到") {
		t.Fatalf("a pointer at a directory must say the full text is not there, got %q", announceOf(out.Text))
	}
	if pointerRe.FindStringSubmatch(out.Text) == nil {
		t.Fatalf("the pointer must survive the honesty fix, got %q", announceOf(out.Text))
	}
}

// TestHealthyPointerStaysSilent is the reverse criterion, and the one that
// keeps this fix from being a slogan: in the root, on disk, a real regular
// file. Pre-fix silence here was accidental; post-fix silence is the contract,
// and it must coincide with fs.read actually returning the whole text.
func TestHealthyPointerStaysSilent(t *testing.T) {
	dir := tempCanonical(t)
	body := asciiRun(20000)
	good := mustWriteFile(t, filepath.Join(dir, "tool-output-174-healthy.txt"), body)

	roster := NewTaskRoster()
	roster.Record("bg-good", TaskOutput{Text: body, ArtifactPath: mustCanonical(t, good)})
	b := judgedTaskBridge(t, roster, NewPathCanonicalizer([]string{dir}, nil))

	out := callTaskOutput(t, b, "bg-good")
	t.Logf("silent control verbatim: IsError=%v Truncated=%v Text=%q",
		out.IsError, out.Truncated, announceOf(out.Text))
	if out.IsError || !out.Truncated {
		t.Fatalf("a healthy pointer must be a plain truncation: %+v", out)
	}
	for _, banned := range []string{noticeLead, "读不到", "不存在", "不是一般文件", "未接线", "不可找回"} {
		if strings.Contains(out.Text, banned) {
			t.Errorf("a pointer that works must not be flagged: found %q in %q", banned, announceOf(out.Text))
		}
	}
	// The silence has to be TRUE, not merely quiet.
	rd := readThroughBridge(t, b, good)
	if rd.IsError || rd.RiskLevel != "L0" || rd.Text != body {
		t.Fatalf("control leg: fs.read must return the whole %d bytes at L0, got isError=%v level=%v len=%d",
			len(body), rd.IsError, rd.RiskLevel, len(rd.Text))
	}
}

// TestUnwiredJudgeFailsClosedInReply is the arm the dispatch names: no
// canonicalizer wired (what cmd/wisp/run.go:362 builds today) must NOT be read
// as "readable". The answer says it cannot verify, in the reply.
func TestUnwiredJudgeFailsClosedInReply(t *testing.T) {
	dir := tempCanonical(t)
	body := asciiRun(20000)
	realFile := mustWriteFile(t, filepath.Join(dir, "tool-output-174-nojudge.txt"), body)

	roster := NewTaskRoster()
	roster.Record("bg-nojudge", TaskOutput{Text: body, ArtifactPath: mustCanonical(t, realFile)})
	reg := NewRegistry()
	if err := reg.Register(Entry{Tool: taskOutput{d: TaskDeps{Roster: roster}}, Decl: TaskOutputDecl()}); err != nil {
		t.Fatal(err)
	}
	b := New(Options{
		Registry: reg, Paths: NewPathCanonicalizer([]string{dir}, nil),
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil}),
		Gate:       NoGate{}, Logf: func(string, ...any) {},
	})

	out := callTaskOutput(t, b, "bg-nojudge")
	t.Logf("fail-closed verbatim: IsError=%v Truncated=%v Text=%q", out.IsError, out.Truncated, announceOf(out.Text))
	if !strings.Contains(out.Text, "未接线") || !strings.Contains(out.Text, "按读不到处理") {
		t.Fatalf("no judge must fail closed out loud, got %q", announceOf(out.Text))
	}
	if strings.Contains(out.Text, "不在你被授权的目录范围内") || strings.Contains(out.Text, "并不存在") {
		t.Errorf("without a judge the reply may not assert either verdict, only the missing tool: %q", announceOf(out.Text))
	}
	if pointerRe.FindStringSubmatch(out.Text) == nil {
		t.Fatalf("failing closed is not the same as withholding the pointer, got %q", announceOf(out.Text))
	}
}

// TestBothPointerDefectsAreReportedSeparately proves the two shapes are two
// facts, not one else-branch: a path outside every root AND absent must say
// both, so neither check can silently shadow the other.
func TestBothPointerDefectsAreReportedSeparately(t *testing.T) {
	dir := tempCanonical(t)
	body := asciiRun(20000)
	// OUTSIDE the only authorized root AND absent: two independent facts. The
	// first run of this leg put the ghost inside dir instead, so only the
	// existence note fired - that was the fixture being wrong, not the code.
	elsewhere := tempCanonical(t)
	roster := NewTaskRoster()
	roster.Record("bg-both", TaskOutput{
		Text:         body,
		ArtifactPath: filepath.Join(elsewhere, "nowhere", "gone.txt"),
	})
	b := judgedTaskBridge(t, roster, NewPathCanonicalizer([]string{dir}, nil))

	out := callTaskOutput(t, b, "bg-both")
	t.Logf("both shapes verbatim: IsError=%v Text=%q", out.IsError, announceOf(out.Text))
	if !strings.Contains(out.Text, "不在你被授权的目录范围内") {
		t.Errorf("the authorization leg is missing: %q", announceOf(out.Text))
	}
	if !strings.Contains(out.Text, "并不存在") {
		t.Errorf("the existence leg is missing: %q", announceOf(out.Text))
	}
	if n := strings.Count(announceOf(out.Text), noticeLead); n != 2 {
		t.Errorf("both notices expected, found %d in %q", n, announceOf(out.Text))
	}
}

// TestPointerNoticeKeepsTheD153StubShape is the collateral guard: honesty is
// added to the announcement, not substituted for it. The head/tail windows, the
// omitted/total numbers and the pointer sentence must read as ticket 164
// accepted them, and the notice must live in Result.Text - not in the progress
// callback.
func TestPointerNoticeKeepsTheD153StubShape(t *testing.T) {
	dir := tempCanonical(t)
	body := asciiRun(20000)
	roster := NewTaskRoster()
	roster.Record("bg-shape", TaskOutput{Text: body, ArtifactPath: dir})
	paths := NewPathCanonicalizer([]string{dir}, nil)

	// The progress callback is the tool's own - Bridge.Execute takes none - so
	// this leg calls the tool directly and reads the same Text the bridge would
	// hand the model.
	var progress []string
	out, err := taskOutput{d: TaskDeps{Roster: roster, Paths: paths}}.Execute(
		context.Background(), mustArgs(t, map[string]any{"task_id": "bg-shape"}),
		func(s string) { progress = append(progress, s) })
	if err != nil {
		t.Fatal(err)
	}
	if got := len(headOf(out.Text)); got != 500*4 {
		t.Errorf("head = %d bytes, want the D15 2000 the accepted tests pin", got)
	}
	if got := len(tailOf(out.Text)); got != 200*4 {
		t.Errorf("tail = %d bytes, want the D15 800 the accepted tests pin", got)
	}
	if headOf(out.Text) != body[:2000] || tailOf(out.Text) != body[len(body)-800:] {
		t.Error("the notice must not disturb which bytes are kept")
	}
	if !strings.Contains(out.Text, "总长 "+itoa(len(body))+" 字节") {
		t.Errorf("total length must still be announced, got %q", announceOf(out.Text))
	}
	if n := omittedAnnounced(out.Text); n != len(body)-2800 {
		t.Errorf("省略 = %d, want %d", n, len(body)-2800)
	}
	if !strings.Contains(out.Text, "全文见 ") {
		t.Errorf("the pointer sentence must stay verbatim, got %q", announceOf(out.Text))
	}
	if strings.Contains(out.Text, "不可找回") {
		t.Errorf("a pointer is given, so the no-copy announcement must not fire: %q", announceOf(out.Text))
	}
	for _, p := range progress {
		if strings.Contains(p, noticeLead) {
			t.Errorf("the notice must be in the reply, not smuggled into the progress callback: %q", p)
		}
	}
}

// TestHealthyReplyStillMatchesThePreFixTemplate is the anti-noise pin in its
// strongest form: for a pointer that works, the reply must be byte-identical to
// the sentence this ticket measured on the UNFIXED code - same template, same
// numbers, nothing appended. Honesty was added where a fact demanded it, not
// sprayed over every stub.
func TestHealthyReplyStillMatchesThePreFixTemplate(t *testing.T) {
	dir := tempCanonical(t)
	body := asciiRun(20000)
	good := mustWriteFile(t, filepath.Join(dir, "tool-output-174-template.txt"), body)
	roster := NewTaskRoster()
	roster.Record("bg-tpl", TaskOutput{Text: body, ArtifactPath: mustCanonical(t, good)})
	b := judgedTaskBridge(t, roster, NewPathCanonicalizer([]string{dir}, nil))

	out := callTaskOutput(t, b, "bg-tpl")
	head, tail := headOf(out.Text), tailOf(out.Text)
	want := fmt.Sprintf(
		"%s\n[…输出已落文件：省略 %d 字符，总长 %d 字节 / 约 %d token，全文见 %s…]\n%s",
		head, len(body)-len(head)-len(tail), len(body), agent.ApproxTokens(body),
		mustCanonical(t, good), tail)
	if out.Text != want {
		t.Errorf("a healthy pointer drifted from the pre-fix sentence:\n got %q\nwant %q", out.Text, want)
	}
}

// announceOf carves the bracketed announcement sentence out of a stub, i.e. the
// part of the reply that is about the pointer.
func announceOf(stub string) string {
	i := strings.Index(stub, "[…")
	if i < 0 {
		return stub
	}
	j := strings.Index(stub[i:], "…]")
	if j < 0 {
		return stub[i:]
	}
	return stub[i : i+j+len("…]")]
}
