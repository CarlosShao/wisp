package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// Ticket 174, landing leg 174-r1: the ALLOWLIST arm of pointerNotice must name
// every road back that actually exists, and no more.
//
// What was wrong, measured rather than assumed. Until this leg the arm said, in
// one sentence:
//
//	这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；
//	要用户先把所属目录加进 [fs] allowed_dirs 才读得回来
//
// That is a half-truth with the same shape the ticket exists to clear. [fs]
// allowed_dirs is a risk-TIER input, not an execution wall: FSDeps.open
// canonicalizes and never re-checks the allowlist (fs.go's open leg), so an L2
// card answered yes really does fetch the bytes. "会被拒 ... 才读得回来"
// describes only the nobody-approves case and states it as the only case.
//
// The other leg had already been corrected: internal/agent/spill.go's
// pointerNotice has named both roads since 174-c2 measured the asymmetry. So
// the two legs disagreed about one verdict - one said "there is a way, two of
// them", the other said "there is one way and you must take it".
//
// This file is the ruler for the corrected sentence, and it is deliberately
// TWO-ARMED so the wording cannot drift into a slogan:
//
//	text arm      both roads are named, and the only-if claim is gone
//	behaviour arm the SECOND road is real on this very bridge - an approving
//	              gate turns the same fs.read from a refusal into the bytes,
//	              and NoGate still refuses. If road two ever stops existing,
//	              the sentence becomes a new lie and this arm goes red.
//
// The reverse criterion is the silence already pinned elsewhere
// (TestHealthyPointerStaysSilent / ...174r8StaysQuiet): naming more roads must
// not turn an authorized pointer into a warning. This file adds no notice to
// the healthy arm, so those pins keep their meaning.
//
// Injection surface (AGENTS §1.3): real Bridge, real C26, real fs.read, real
// files on disk, real gate seam (GateFuncs through the bridge's own Options).
// No mock stands in for a path decision, and no authorization is widened: the
// arm that must speak authorizes nothing at all.

// twoRoadsBridge174r1 wires fs.* and task.* over ONE canonicalizer and the
// supplied gate - the same shape judgedTaskBridge builds, plus the gate, which
// is what lets the behaviour arm ask the question "is road two real?" of the
// real routing instead of of the wording.
func twoRoadsBridge174r1(t *testing.T, roster *TaskRoster, paths *PathCanonicalizer, g Gate) *Bridge {
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
		Gate:       g, Logf: func(string, ...any) {},
	})
}

// TestPointerOutsideAllowlistNamesEveryRoadBack174r1 is the text half: the
// allowlist arm must name both roads, and must no longer claim there is only
// one. The only-if spelling is banned as a literal because it is the exact
// clause that made the old sentence false - a promise with an "only" in it.
func TestPointerOutsideAllowlistNamesEveryRoadBack174r1(t *testing.T) {
	dir := tempCanonical(t)
	outside := tempCanonical(t)
	body := asciiRun(20000)
	realFile := mustWriteFile(t, filepath.Join(outside, "tool-output-174r1-outside.txt"), body)

	roster := NewTaskRoster()
	roster.Record("bg-174r1", TaskOutput{Text: body, ArtifactPath: mustCanonical(t, realFile)})
	// An EMPTY allowlist is the factory default (SPEC-03 allowed_dirs[]=[]), so
	// this arm authorizes nothing: it is the shape the ticket is about.
	b := twoRoadsBridge174r1(t, roster, NewPathCanonicalizer(nil, nil), NoGate{})

	out := callTaskOutput(t, b, "bg-174r1")
	t.Logf("allowlist arm verbatim: IsError=%v Truncated=%v Text=%q",
		out.IsError, out.Truncated, announceOf(out.Text))

	if out.IsError || !out.Truncated {
		t.Fatalf("a recorded long output must still be answered as a truncation: %+v", out)
	}
	// Both roads must be named. Each marker is one road; requiring both is what
	// makes "the model is told what it can do" a measured claim.
	for _, want := range []string{
		"不在你被授权的目录范围内",      // the verdict itself
		"[fs] allowed_dirs", // road one: the user widens the root
		"L2 卡",              // road two: one approved card still fetches it
	} {
		if !strings.Contains(out.Text, want) {
			t.Errorf("the allowlist notice must name %q, got %q", want, announceOf(out.Text))
		}
	}
	// The only-if clause is gone. This is the load-bearing deletion: "才读得
	// 回来" said the allowlist was the ONLY way back, and that was the lie.
	if strings.Contains(out.Text, "才读得回来") {
		t.Errorf("the notice must not claim the allowlist is the only road back, got %q",
			announceOf(out.Text))
	}
	// The pointer itself stays (PLAN.md:2564; 只截不指＝不合格).
	if pointerRe.FindStringSubmatch(out.Text) == nil {
		t.Fatalf("the pointer must still be given, got %q", announceOf(out.Text))
	}
	// The rectified silence arm is a different arm's job, but the file is the
	// same: this one must not have become a blanket warning. One defect here
	// (authorization), so exactly one notice.
	if n := strings.Count(announceOf(out.Text), noticeLead); n != 1 {
		t.Errorf("this call has exactly one defect (authorization), so one notice; found %d in %q",
			n, announceOf(out.Text))
	}
	_ = dir
}

// TestPointerSecondRoadBackIsRealOnThisBridge174r1 is the behaviour half, and
// the reason the sentence above is allowed to say what it says.
//
// It measures the asymmetry itself on ONE bridge, ONE path, ONE canonicalizer:
//
//	nobody approves (NoGate)          -> refused, L2, and the refusal names the card
//	someone approves (gate says yes)  -> the SAME fs.read returns the whole body
//
// If the second line ever stops holding, "回来的路有两条 ... 要么批下那一张卡"
// becomes a new false promise and this case goes red - which is exactly the
// arrangement the ticket wants: the wording is not permitted to outlive the
// fact it describes.
func TestPointerSecondRoadBackIsRealOnThisBridge174r1(t *testing.T) {
	dir := tempCanonical(t)
	outside := tempCanonical(t)
	body := asciiRun(20000)
	realFile := mustWriteFile(t, filepath.Join(outside, "tool-output-174r1-road2.txt"), body)
	filed := mustCanonical(t, realFile)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	roster := NewTaskRoster()
	roster.Record("bg-174r1-road2", TaskOutput{Text: body, ArtifactPath: filed})

	// Arm one: nobody is there to approve. The notice's "没人批就是直接拒" must
	// be the true reading, so this arm pins the refusal.
	refusing := twoRoadsBridge174r1(t, roster, NewPathCanonicalizer(nil, nil), NoGate{})
	rd := readThroughBridge(t, refusing, realFile)
	t.Logf("road two, nobody approves verbatim: IsError=%v level=%v class=%q text=%q",
		rd.IsError, rd.RiskLevel, rd.ErrorClass, rd.Text)
	if !rd.IsError || rd.RiskLevel != "L2" {
		t.Fatalf("with no gate the out-of-root read must be refused at L2, got %+v", rd)
	}

	// Arm two: a human answers the card. Same path, same canonicalizer, same
	// tool - only the gate answer differs. The bytes must come back whole.
	approving := twoRoadsBridge174r1(t, roster, NewPathCanonicalizer(nil, nil),
		&gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow})
	rd2 := readThroughBridge(t, approving, realFile)
	t.Logf("road two, card approved verbatim: IsError=%v level=%v bytes=%d",
		rd2.IsError, rd2.RiskLevel, len(rd2.Text))
	if rd2.IsError {
		t.Fatalf("road two must be REAL for the notice to name it: an approved card still got %+v", rd2)
	}
	if rd2.Text != body {
		t.Fatalf("road two must return the whole %d bytes, got %d", len(body), len(rd2.Text))
	}
	// And it was the card that changed the outcome, not a widened allowlist:
	// the same bridge still judges the path out of scope.
	if canon, err := NewPathCanonicalizer(nil, nil).Canonicalize(filed); err != nil {
		t.Fatalf("precondition: the filed path must canonicalize, got %v", err)
	} else if NewPathCanonicalizer(nil, nil).InAllowlist(canon) {
		t.Fatal("precondition: this fixture must stay OUTSIDE every authorized root")
	}
}

// TestPointerSecondRoadBackIsBookedAsACard174r1 pins the middle of that claim:
// the read really did route through the L2 queue, so "会要一张 L2 卡" is the
// measured mechanism and not a figure of speech. Without this, the arm above
// could be explained by any widening at all.
func TestPointerSecondRoadBackIsBookedAsACard174r1(t *testing.T) {
	outside := tempCanonical(t)
	body := asciiRun(2000)
	realFile := mustWriteFile(t, filepath.Join(outside, "tool-output-174r1-booked.txt"), body)

	g := &gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow}
	b := twoRoadsBridge174r1(t, NewTaskRoster(), NewPathCanonicalizer(nil, nil), g)

	rd := readThroughBridge(t, b, realFile)
	windows, approvals := g.counts()
	t.Logf("routing verbatim: windows=%d approvals=%d decision=%+v", windows, approvals, g.approvalDecision())
	if rd.IsError {
		t.Fatalf("the approving gate should have let this through, got %+v", rd)
	}
	if approvals != 1 {
		t.Fatalf("the out-of-root read must ask the L2 queue exactly once, got %d (windows=%d)",
			approvals, windows)
	}
	if windows != 0 {
		t.Errorf("an L2 verdict must not be routed to the L1 window, got %d window calls", windows)
	}
	dec := g.approvalDecision()
	if !strings.Contains(dec.Reason, "授权目录之外") {
		t.Errorf("the card shown to the human must name the allowlist verdict, got %q", dec.Reason)
	}
}

// announceRegion174r1 carves the notice out of a stub, so a sweep can judge the
// tool's own words rather than the pointer sentence - which legitimately prints
// the path and is frozen by PLAN.md:2564. Local to this file: the r3/r6 sweeps
// carry their own carve because their bounds differ (they cut at the pointer).
func announceRegion174r1(stub string) string {
	i := strings.Index(stub, noticeLead)
	if i < 0 {
		return ""
	}
	region := stub[i:]
	if j := strings.Index(region, "，全文见 "); j >= 0 {
		region = region[:j]
	}
	return region
}

// TestPointerTwoRoadNoticeLeaksNoPathOrRoots174r1 is the guard that keeps the
// longer sentence from becoming the new leak: naming two roads must not smuggle
// the filed path, the authorized root list, or C26's own error text into the
// notice. Same three classes the r3/r6/r7 sweeps ban, re-asked of the arm this
// leg rewrote - so the rewrite cannot pass its forward arms by widening what it
// says about the host's internals.
func TestPointerTwoRoadNoticeLeaksNoPathOrRoots174r1(t *testing.T) {
	rootA, rootB := tempCanonical(t), tempCanonical(t)
	outside := tempCanonical(t)
	body := asciiRun(20000)
	filed := filepath.Join(outside, "tool-output-174r1-sweep.txt")
	realFile := mustWriteFile(t, filed, body)

	roster := NewTaskRoster()
	roster.Record("bg-174r1-sweep", TaskOutput{Text: body, ArtifactPath: mustCanonical(t, realFile)})
	// Two roots, so the root-list ban is a list and not one directory.
	b := twoRoadsBridge174r1(t, roster, NewPathCanonicalizer([]string{rootA, rootB}, nil), NoGate{})

	out := callTaskOutput(t, b, "bg-174r1-sweep")
	notice := announceRegion174r1(out.Text)
	t.Logf("notice region under the sweep: %q", notice)
	if notice == "" {
		t.Fatalf("this arm must speak before what it says can be swept: %q", announceOf(out.Text))
	}
	for _, banned := range []string{
		mustCanonical(t, realFile), // the filed path literal
		rootA, rootB,               // the authorized root list, entry by entry
		"tools:", "risk:", // C26's own error-text prefixes
	} {
		if strings.Contains(notice, banned) {
			t.Errorf("the two-road notice must not relay %q: notice %q", banned, notice)
		}
	}
}
