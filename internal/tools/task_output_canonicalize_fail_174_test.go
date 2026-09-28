package tools

import (
	"context"
	"strings"
	"testing"
)

// Ticket 174 AC#2d, first resident criterion for the branch nothing pinned.
//
// pointerNotice (internal/tools/task.go:318-340) answers in five shapes, and
// the eight criteria shipped by 174-r1 pin four of them: unwired judge (:319),
// outside every authorized root (:327), copy file missing (:332), not a regular
// file (:334). The fifth - `case err != nil` at task.go:325-326, the arm where
// C26 cannot even normalize the string the host filed - had ZERO coverage:
//   grep -c '规范化' internal/tools/task_output_pointer_notice_test.go -> 0
//   grep -rn "连规范化都没通过" --include=*_test.go internal/         -> no hit
// so 174-v1's MUT-V5 removed its whole sentence and the package stayed green
// (rc=0, 124 PASS). This file is the criterion that makes that impossible.
//
// Which arm production actually walks today is the unwired one
// (cmd/wisp/run.go builds TaskDeps without Paths - 174-c2 §1); wiring is NOT
// this leg's to do, so both arms here are constructed at the seam, which is
// exactly where the shape was measured in the first place (AGENTS §1.3: real
// bridge, real C26 canonicalizer, real roster records; no mock stands in for a
// real component).
//
// Three forms, all load-bearing:
//
//	forward   Canonicalize RETURNS AN ERROR (not "normalized fine but outside
//	          the roots" - that arm already has TestPointerOutsideAuthorizedRootSpeaks)
//	          => the reply takes the fail-closed branch and is never an empty
//	          success: truncation answered, pointer still given, notice present.
//	reverse   the branch is named by its OWN claim ("连规范化都没通过"), so
//	          muting it, replacing it with "" (assume readable), rerouting it to
//	          the allowlist verdict, or rerouting it to the unwired sentence each
//	          turn a leg red here. Measured with -overlay, never by editing a
//	          tracked file; see docs/evidence/s1/174-ac2d-criterion-r2.md §3.
//	boundary  the reply must not say this path can be read ("读得回来"), must not
//	          hand the model the authorized root list (dir is checked absent),
//	          and must not fall through to the no-pointer announcement. This is
//	          deliberately NOT a byte-for-byte text pin: :299 and :349 of
//	          task_output_pointer_notice_test.go freeze the D15(3) stub shape and
//	          the healthy-reply template, and this leg may not disturb them.

// canonicalizeFailsHere is the input that reaches task.go:325 on this host
// without needing a reparse point on disk: a whitespace-only artifact path is
// non-empty, so the pointer branch at task.go:242 takes it, and
// PathCanonicalizer.Canonicalize (paths.go:120-122) refuses it with
// "tools: empty path" before risk.Resolve is ever called.
//
// The other reachable error is risk.ErrReparseDenied, which needs a real
// junction (mkRealJunction lives in the windows-only
// bridge_junction_windows_test.go); that shape is registered as NOT measured by
// this leg rather than faked.
const canonicalizeFailsHere = "   "

// TestCanonicalizeFailureFailsClosedInReply174 is the forward form: C26 returns
// an error, and the reply must say the pointer cannot be followed instead of
// passing it off as healthy.
func TestCanonicalizeFailureFailsClosedInReply174(t *testing.T) {
	dir := tempCanonical(t)
	body := asciiRun(20000)

	roster := NewTaskRoster()
	roster.Record("bg-c174", TaskOutput{Text: body, ArtifactPath: canonicalizeFailsHere})
	// A judge IS wired here, and its root is a real tree: this arm is the one
	// that gets past `d.Paths == nil` and into `Canonicalize`, which is the whole
	// point of the leg (the unwired arm is pinned by 174-r1 already).
	b := judgedTaskBridge(t, roster, NewPathCanonicalizer([]string{dir}, nil))

	out := callTaskOutput(t, b, "bg-c174")
	t.Logf("canonicalize-error verbatim: IsError=%v Truncated=%v announce=%q",
		out.IsError, out.Truncated, announceOf(out.Text))

	if out.IsError || !out.Truncated {
		t.Fatalf("a recorded long output must still be answered as a truncation, not a refusal: %+v", out)
	}
	// The branch speaks: fail closed, in the text the model reads.
	if !strings.Contains(out.Text, noticeLead) || !strings.Contains(out.Text, "读不到") {
		t.Fatalf("a path C26 cannot normalize must be handled as unreadable, got %q", announceOf(out.Text))
	}
	// This branch specifically - by its own claim, so muting it cannot hide
	// behind another branch's sentence.
	if !strings.Contains(out.Text, "连规范化都没通过") {
		t.Fatalf("the reply must name normalization as the step that failed, got %q", announceOf(out.Text))
	}
	// Not an empty success: the pointer itself is still owed (PLAN.md:2564).
	// pointerRe cannot serve here - it needs a non-space path, and this fixture's
	// filed string IS spaces - so the leg pins the printed pointer against the
	// exact string the roster holds instead. (Measured on the first run of this
	// leg: the reply ends "…，全文见    …]" - pointer given, no \S to match.)
	if !strings.Contains(out.Text, "全文见 "+canonicalizeFailsHere) {
		t.Fatalf("failing closed is not withholding the pointer, got %q", announceOf(out.Text))
	}

	// Boundary form: nothing in the reply may promise a road back, and the
	// notice may not recite the authorization verdict or the root list.
	for _, banned := range []string{
		"读得回来", "随时可读", "可以读回", "不在你被授权的目录范围内", "不可找回", "未接线",
	} {
		if strings.Contains(out.Text, banned) {
			t.Errorf("the canonicalize-error reply must not contain %q: %q", banned, announceOf(out.Text))
		}
	}
	if strings.Contains(out.Text, dir) {
		t.Errorf("the notice must not relay the authorized root list to the model, found %q in %q",
			dir, announceOf(out.Text))
	}
}

// TestCanonicalizeFailureIsNotTheUnwiredArm174 pins that the two fail-closed
// sentences stay two different facts on the SAME fixture. Without this leg a
// "fix" could route the normalization error into the unwired-judge text (or the
// other way round) and keep every forward assertion green, which is the shape
// 174-v1 counted as a criterion with no teeth.
func TestCanonicalizeFailureIsNotTheUnwiredArm174(t *testing.T) {
	body := asciiRun(20000)
	dir := tempCanonical(t)

	wiredRoster := NewTaskRoster()
	wiredRoster.Record("bg-wired-174", TaskOutput{Text: body, ArtifactPath: canonicalizeFailsHere})
	wired, err := taskOutput{d: TaskDeps{
		Roster: wiredRoster,
		Paths:  NewPathCanonicalizer([]string{dir}, nil),
	}}.Execute(context.Background(), mustArgs(t, map[string]any{"task_id": "bg-wired-174"}), nil)
	if err != nil {
		t.Fatalf("task.output answers with data, not a Go error: %v", err)
	}
	if !strings.Contains(wired.Text, "连规范化都没通过") {
		t.Fatalf("wired judge must report the normalization failure, got %q", announceOf(wired.Text))
	}
	if strings.Contains(wired.Text, "未接线") {
		t.Errorf("the judge IS wired here; the reply may not blame the missing wiring: %q", announceOf(wired.Text))
	}

	nilRoster := NewTaskRoster()
	nilRoster.Record("bg-nil-174", TaskOutput{Text: body, ArtifactPath: canonicalizeFailsHere})
	unwired, err := taskOutput{d: TaskDeps{Roster: nilRoster}}.Execute(
		context.Background(), mustArgs(t, map[string]any{"task_id": "bg-nil-174"}), nil)
	if err != nil {
		t.Fatalf("task.output answers with data, not a Go error: %v", err)
	}
	t.Logf("unwired control verbatim: announce=%q", announceOf(unwired.Text))
	if !strings.Contains(unwired.Text, "未接线") || !strings.Contains(unwired.Text, "按读不到处理") {
		t.Fatalf("no judge must stay the fail-closed arm 174-r1 pinned, got %q", announceOf(unwired.Text))
	}
	if strings.Contains(unwired.Text, "连规范化都没通过") {
		t.Errorf("without a judge nothing was ever canonicalized; the reply may not claim it was: %q",
			announceOf(unwired.Text))
	}
}
