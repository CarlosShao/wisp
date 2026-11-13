package main

// Ticket 219 AC#2 - the 长期 branch leaves a revocable line, and the undo
// round-trips byte for byte.
//
// AC#2's own sentence: 「长期」那一支必须留下可撤销的一行——审计里能查到「谁、
// 什么时候、往 allowed_dirs 加了哪一行、来源是审批卡片」；并有一条判据钉住
// 「加完之后撤销能一字回到加之前」.
//
// WHAT WAS ALREADY PINNED (the two halves AC#2 does NOT re-buy, so no second
// copy is minted here): ticket 201's approval_always_201_test.go pins that the
// widening stores the rule only behind a second, L2-level card, and ticket
// 226's approval_revoke_226_test.go pins that the reply grammar's revoke verb
// removes exactly the stored line and refuses what the process does not hold.
// Neither file pins the ROUND TRIP as AC#2 spells it: the same config.toml,
// byte compared before and after, the four audit fields read off the line, and
// the revocation being REFUSED once the rule is already gone.
//
// THE FOUR AUDIT FIELDS, MEASURED NOT ASSUMED. The production writer of the
// WIDEN-APPLIED line is approval_always.go:140, booked through
// replySurface.record as
//
//	approval: REPLY WIDEN-APPLIED corr="<card id>" tool=config.allow_dir
//	    route="native/always" detail="dir=... section=fs key=allowed_dirs
//	    effect=next-start" when=<RFC3339>
//
// which answers all four on the line itself: 哪一行 = the dir= field, 来源 =
// route="native/always" (the reply grammar's own route label), 谁 = the corr
// field (the card the widening was staged from), 什么时候 = the when= field
// the record writer stamps. The same sentence also reaches the run's durable
// JSONL log (run.go:1428 hands it to the slog sink, whose JSONHandler stamps
// every record); this fixture reads the console copy because that is the copy
// an operator and every 226 tooth already read, and it parses the stamp as a
// real time instead of trusting that the field exists.
//
// THE BYTE-FOR-BYTE HALF AND ITS HONEST BOUNDARY. The store-then-revoke pair
// runs through Manager.AddAllowedDir and Manager.SetAllowedDirs, whose guarded
// writer (internal/config/writeguard.go mergeWrite -> SaveFile) re-lays the
// file in canonical form on every real write. Byte equality of the whole file
// therefore holds only when the file already IS canonical bytes - which is
// exactly the production shape of this fixture: the run's config.toml is
// canonicalized before the task starts (run.go), and this file re-canonicalizes
// it explicitly, so "the file after revoke == the file before always" is a
// byte comparison of MarshalCanonical output against MarshalCanonical output,
// and it is red the moment the revocation forgets a line, adds one, or
// reorders one. On a hand-written file the guarded writer preserves every
// VALUE (that is writeguard AC#1, pinned in internal/config) but re-lays the
// FORM - that boundary is ticket 226's landing cost, documented in
// writeguard.go's own header, and it is reported in the evidence file instead
// of being softened into a weaker comparison here.
//
// THE MUTATION CONTROL. The M1 variant of approval_revoke.go (the
// SetAllowedDirs call handed the ORIGINAL list - a no-op that still books its
// audit lines) is applied at test time through `go test -overlay`, so the
// pristine file on disk is never touched; the round trip must go RED under it
// (the bytes never come back) and the pristine tree stays green.

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
)

// canonicalConfig219 rewrites cfgPath as config.MarshalCanonical bytes, so the
// fixture file IS canonical output before anything touches it - the condition
// under which the guarded writer's own byte-for-byte round trip holds.
func canonicalConfig219(t *testing.T, cfgPath string, c *config.Config) {
	t.Helper()
	body, err := config.MarshalCanonical(c)
	if err != nil {
		t.Fatalf("marshal canonical fixture: %v", err)
	}
	if err := os.WriteFile(cfgPath, body, 0o600); err != nil {
		t.Fatalf("write canonical fixture: %v", err)
	}
}

// canonicalFixture219 loads the host's config through the next start's own
// loader and re-lays it canonically. Returns the bytes now on disk.
func canonicalFixture219(t *testing.T, cfgPath string) []byte {
	t.Helper()
	loaded, _, err := config.LoadFile(cfgPath, nil)
	if err != nil {
		t.Fatalf("the fixture config does not load: %v", err)
	}
	canonicalConfig219(t, cfgPath, loaded)
	body, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read back the canonical fixture: %v", err)
	}
	return body
}

// beforeBytes219 was folded into canonicalFixture219: the bytes it returns are
// the pre-store reading the round trip is measured against.

// driveStoreHoldRevoke drives the 长期 flow to a stored rule, holds the run
// there - both cards answered, the original card retired, the widening worker
// done - until the test has taken its readings, then closes the stream. Each
// revoke line is optional ("" skips); storedOut, when non-nil, receives the
// file bytes between the store and the first revoke: the proof there was a
// line to revoke in the first place.
//
// THE TYPING ORDER IS THE RACE CONTROL. The task's own round completes the
// moment the mock's text-only response lands (the rig's fs.write pend rides
// its own Background context, not the task), and runTextTask's deferred
// rt.close() then cancels the reply root - runReplyLoop refuses any line it
// checks after that. Every line that still needs the listener is therefore
// fed BACK TO BACK once the widening is booked, and FIFO does the sequencing:
// the listener processes revoke1, then revoke2, then the original answer. The
// sentences are waited for AFTER the typing, so a line the teardown overtook
// fails its wait loudly instead of deadlocking the pipe writer.
func driveStoreHoldRevoke(t *testing.T, w *widenRun, pw *io.PipeWriter,
	revoke1, want1, revoke2, want2, answerOriginal string, storedOut *[]byte,
) {
	t.Helper()
	defer func() { _ = pw.Close() }()
	defer close(w.reply)

	waitForSentence187(t, w.h, "卡片编号：", "the first card must show the address the operator answers against")
	if _, err := fmt.Fprintln(pw, "always "+w.corr); err != nil {
		w.setFail(err)
		return
	}
	waitForSentence187(t, w.h, "要存的规则：", "the 长期 answer must print the rule it stores")
	corr := waitWidenCard187(t, w, 15*time.Second)
	if corr == "" {
		w.setFail(errors.New("the second, L2-level re-confirmation card never appeared"))
		return
	}
	if _, err := fmt.Fprintln(pw, "yes "+corr); err != nil {
		w.setFail(err)
		return
	}
	if !waitForAudit187(w.h, "approval: REPLY WIDEN-APPLIED", 25*time.Second) {
		w.setFail(errors.New("the widening never booked WIDEN-APPLIED, so there was no rule to revoke"))
		return
	}
	if storedOut != nil {
		body, err := os.ReadFile(filepath.Join(w.h.dir, configFileName))
		if err != nil {
			w.setFail(err)
			return
		}
		*storedOut = body
	}
	for _, line := range []string{revoke1, revoke2, answerOriginal + " " + w.corr} {
		if line == "" {
			continue
		}
		if _, err := fmt.Fprintln(pw, line); err != nil {
			w.setFail(err)
			return
		}
	}
	for _, step := range []struct{ line, want string }{{revoke1, want1}, {revoke2, want2}} {
		if step.line == "" {
			continue
		}
		waitForSentence187(t, w.h, step.want, "the revoke verb must answer on the console")
	}
}

// Test219AC2WidenAppliedAuditLineCarriesTheFourFields is AC#2's first half:
// the audit line a later auditor acts on answers 谁 corr= / 什么时候 when= /
// 哪一行 dir= / 来源 route="native/always" - all four read off the line, and
// the stamp parsed as a real RFC3339 time instead of trusted.
func Test219AC2WidenAppliedAuditLineCarriesTheFourFields(t *testing.T) {
	w, pw := newWidenRun(t, 40*time.Second, "t219-ac2-audit", "t219-ac2-audit-corr")
	go driveStoreHoldRevoke(t, w, pw, "", "", "", "", "yes", nil)

	if code := w.h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, w.h.out.String(), w.h.err.String())
	}
	waitRoundTripBounded314(t, w)
	if err := w.failure(); err != nil {
		t.Fatal(err)
	}

	audit := w.h.err.String()
	line := ""
	for _, l := range strings.Split(audit, "\n") {
		if strings.Contains(l, "approval: REPLY WIDEN-APPLIED") {
			line = l
			break
		}
	}
	if line == "" {
		t.Fatalf("no WIDEN-APPLIED line was booked; audit:\n%s", audit)
	}
	for _, want := range []string{
		`corr="t219-ac2-audit-corr"`,
		`tool=config.allow_dir`,
		`route="native/always"`,
		"dir=",
		"when=",
	} {
		if !strings.Contains(line, want) {
			t.Errorf("WIDEN-APPLIED line is missing %q; the line is:\n%s", want, line)
		}
	}
	idx := strings.LastIndex(line, " when=")
	if idx < 0 {
		t.Fatalf("no when= stamp to parse; the line is:\n%s", line)
	}
	stamp := strings.TrimSpace(line[idx+len(" when="):])
	when, err := time.Parse(time.RFC3339, stamp)
	if err != nil {
		t.Fatalf("the when= stamp is not RFC3339 (%q): %v\nline:\n%s", stamp, err, line)
	}
	if d := time.Since(when); d < -time.Minute || d > 10*time.Minute {
		t.Errorf("the when= stamp is %v away from this run: %s", d, line)
	}
}

// Test219AC2RevokeRoundTripByteForByte is AC#2's second half, as one round
// trip on the same file: store a rule, revoke it, and the file is byte for
// byte what it was BEFORE the store - not the same length, not "no longer
// contains the directory", the same bytes.
func Test219AC2RevokeRoundTripByteForByte(t *testing.T) {
	w, pw := newWidenRun(t, 40*time.Second, "t219-ac2-roundtrip", "t219-ac2-roundtrip-corr")
	cfgPath := filepath.Join(w.h.dir, configFileName)
	before := canonicalFixture219(t, cfgPath)

	rule := `[fs] allowed_dirs += "` + filepath.ToSlash(w.h.outside) + `"`
	var stored []byte
	w.storedBytes219 = &stored
	go driveStoreHoldRevoke(t, w, pw, "revoke "+rule, "已撤销长期规则", "", "", "yes", &stored)

	if code := w.h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, w.h.out.String(), w.h.err.String())
	}
	waitRoundTripBounded314(t, w)
	if err := w.failure(); err != nil {
		t.Fatal(err)
	}

	after, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(stored), filepath.ToSlash(w.h.outside)) {
		t.Fatalf("the widening never stored the rule this case revokes; config.toml was:\n%s", stored)
	}
	if !bytesEqual219(before, after) {
		t.Fatalf("AC#2 round trip failed: the file after revoke is not byte for byte the file before store.\n"+
			"before (%d bytes):\n%s\nafter (%d bytes):\n%s", len(before), before, len(after), after)
	}
}

// Test219AC2SecondRevokeOfTheSameRuleIsRefused is the trip's closed end: after
// the round trip, the same revoke line is refused out loud - booked as a
// refusal - and the file does not change by one byte.
func Test219AC2SecondRevokeOfTheSameRuleIsRefused(t *testing.T) {
	w, pw := newWidenRun(t, 40*time.Second, "t219-ac2-twice", "t219-ac2-twice-corr")
	cfgPath := filepath.Join(w.h.dir, configFileName)
	before := canonicalFixture219(t, cfgPath)

	rule := `[fs] allowed_dirs += "` + filepath.ToSlash(w.h.outside) + `"`
	var stored []byte
	w.storedBytes219 = &stored
	go driveStoreHoldRevoke(t, w, pw,
		"revoke "+rule, "已撤销长期规则",
		"revoke "+rule, "本进程账上没有",
		"yes", &stored)

	if code := w.h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, w.h.out.String(), w.h.err.String())
	}
	waitRoundTripBounded314(t, w)
	if err := w.failure(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(stored), filepath.ToSlash(w.h.outside)) {
		t.Fatalf("the widening never stored the rule; config.toml was:\n%s", stored)
	}
	after, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytesEqual219(before, after) {
		t.Fatalf("the refused second revocation changed the file.\n"+
			"before (%d bytes):\n%s\nafter (%d bytes):\n%s", len(before), before, len(after), after)
	}
	audit := w.h.err.String()
	if !strings.Contains(audit, `approval: REPLY REFUSED corr="" tool=config.allow_dir route="native/revoke"`) {
		t.Errorf("the refused second revocation was not booked as a refusal:\n%s", audit)
	}
}

// Test219AC2MutationRevokeWriteMustPersistForBytesToComeBack is the tooth's
// mutation control: with the SetAllowedDirs call in approval_revoke.go
// replaced by a write of the ORIGINAL list (the M1 no-op that still books its
// audit lines), the round trip must go RED - the bytes never come back. The
// mutation is applied through `go test -overlay` on a JSON written to
// t.TempDir(), so the pristine file on disk is never touched for one byte
// (AGENTS.md git discipline: no in-place mutation of tracked files).
func Test219AC2MutationRevokeWriteMustPersistForBytesToComeBack(t *testing.T) {
	if testing.Short() {
		t.Skip("the mutation run recompiles the whole package; skip under -short")
	}
	pristine, err := os.ReadFile("approval_revoke.go")
	if err != nil {
		t.Fatalf("read the pristine revoke verb: %v", err)
	}
	const m1Old = "s.rt.mgr.SetAllowedDirs(listWithoutDir(before.FS.AllowedDirs, dir))"
	const m1New = "s.rt.mgr.SetAllowedDirs(before.FS.AllowedDirs)"
	src := string(pristine)
	if !strings.Contains(src, m1Old) {
		t.Fatalf("the mutation anchor moved in approval_revoke.go - rewrite this tooth's M1 against the current line:\n%s", m1Old)
	}
	mutated := strings.Replace(src, m1Old, m1New, 1)

	pkgDir, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	// The test binary's working directory is the PACKAGE directory; the repo
	// root is whatever lies above it, found by walking for go.mod rather than
	// assumed to be one filepath.Dir away.
	repoRoot := pkgDir
	for {
		if _, err := os.Stat(filepath.Join(repoRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(repoRoot)
		if parent == repoRoot {
			t.Fatal("go.mod not found above the package directory; the mutation run has nowhere to point go at")
		}
		repoRoot = parent
	}
	pristinePath := filepath.Join(pkgDir, "approval_revoke.go")
	mutatedPath := filepath.Join(t.TempDir(), "approval_revoke.go")
	if err := os.WriteFile(mutatedPath, []byte(mutated), 0o600); err != nil {
		t.Fatal(err)
	}
	overlay := map[string]map[string]string{
		"Replace": {pristinePath: mutatedPath},
	}
	overlayBody, err := json.Marshal(overlay)
	if err != nil {
		t.Fatal(err)
	}
	overlayPath := filepath.Join(t.TempDir(), "overlay-m1.json")
	if err := os.WriteFile(overlayPath, overlayBody, 0o600); err != nil {
		t.Fatal(err)
	}

	buildDir := filepath.Join(repoRoot, "build")
	sherpaDir := filepath.Join(repoRoot, "third_party", "sherpa-onnx")
	childPath := buildDir + string(os.PathListSeparator) + sherpaDir +
		string(os.PathListSeparator) + os.Getenv("PATH")
	// The child inherits THIS process's environment, minus the temp vars (a
	// set-but-EMPTY TMP makes Windows hand out non-absolute temp roots - the
	// winsec seam refuses them - so absence, not an empty value, is the
	// discipline the dispatch spells for a shell run, spelled here because the
	// child is a second go test invocation inside one test binary).
	env := make([]string, 0, len(os.Environ()))
	for _, kv := range os.Environ() {
		key, _, _ := strings.Cut(kv, "=")
		switch key {
		case "PATH", "GOTMPDIR", "TMPDIR", "TEMP", "TMP":
			continue
		}
		env = append(env, kv)
	}
	env = append(env, "PATH="+childPath)
	cmd := exec.Command("go", "test", "./cmd/wisp",
		"-run", "^Test219AC2RevokeRoundTripByteForByte$",
		"-count=1", "-timeout", "20m",
		"-overlay="+overlayPath)
	cmd.Dir = repoRoot
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("the M1 mutation stayed green (the round trip no longer pins the write):\n%s", out)
	}
	if !strings.Contains(string(out), "FAIL") ||
		!strings.Contains(string(out), "not byte for byte") {
		t.Fatalf("the mutated run did not fail on this tooth's own byte assertion - the harness is broken, not the tooth:\n%s", out)
	}
}

// bytesEqual219 is this file's one byte comparison, spelled locally so no
// caller can quietly swap it for a length or substring check.
func bytesEqual219(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// waitRoundTripBounded314 is ticket 314's bounded replacement for the bare
// <-w.done / <-w.reply pair the three AC#2 cases used to do. w.done closes when
// the fs.write bridge.Execute returns; w.reply closes when the reply driver
// returns. The work behind w.done is bounded by the chain this fixture itself
// lays down: the run's config carries confirm_timeout_sec = 40
// (newReplyHost's l2Wait, approval_reply_201_test.go:154 - the value every
// replyRoute case has passed since ticket 201), and the store is held until
// waitWidenCard187 answers the second card within its own 15s limit
// (driveStoreHoldRevoke, approval_always_219ac2_test.go:135) after the first
// card the driver answers by line. So a hang here cannot be one of those
// clocks quietly failing: it can only mean a step between them never happened,
// which is the silent shape the diagnosis leg measured (8m31s of zero output,
// external kill, no stack, no name). The bound below is therefore NOT a number
// tuned to let this machine pass: it is this fixture's own 40s L2 window - the
// LONGEST clock the awaited work carries - plus a 6s margin for the remaining
// write/audit steps, and it exists only so the failure is named HERE, with the
// readings that decide the case at expiry: whether each of the two awaited
// closes did land inside the bound, and how far the store-hold phase got
// (stored bytes collected = WIDEN-APPLIED booked).
//
// ONE BOUND FOR BOTH CLOSES, NOT A CHAIN: the total wait across done and reply
// can never exceed the bound twice over - done gets whatever the deadline has
// left, and an expired done leaves reply zero, so the case exits by the same
// absolute deadline no matter which close is the wedged one. (A bare <-w.done
// after a satisfied reply was itself an unbounded receive - the exact shape
// this helper exists to end.)
func waitRoundTripBounded314(t *testing.T, w *widenRun) {
	t.Helper()
	const bound = 46 * time.Second
	deadline := time.Now().Add(bound)
	stored := func() bool {
		return w.storedBytes219 != nil && len(*w.storedBytes219) > 0
	}
	clockLine := "the run's own clocks bound the work it awaits at the 40s L2 " +
		"window (confirm_timeout_sec, approval_reply_201_test.go:154) + the 15s " +
		"second-card wait (approval_always_219ac2_test.go:135), so an expiry " +
		"here means a step between them never happened"

	select {
	case <-w.reply:
	case <-time.After(time.Until(deadline)):
		doneClosed := false
		select {
		case <-w.done:
			doneClosed = true
		default:
		}
		t.Fatalf("the reply driver never returned: waited on w.reply %v (the round's whole bound); "+
			"%s. At expiry: the fs.write side %s, the widening reached the store=%v "+
			"- the driver's last typed lines are in the run's stdout/stderr above",
			bound, clockLine,
			map[bool]string{true: "DID return (done closed)", false: "also never returned (done still open)"}[doneClosed],
			stored())
	}
	select {
	case <-w.done:
	case <-time.After(time.Until(deadline)):
		t.Fatalf("the reply driver returned but the fs.write never did: waited on w.done "+
			"until the round's 46s bound expired; %s. At expiry: the widening reached "+
			"the store=%v - the fs.write is parked behind its unanswered or un-expired "+
			"approval card", clockLine, stored())
	}
}
