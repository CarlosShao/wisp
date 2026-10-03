package agent

// Ticket 174 AC#2c: 「同一句"全文见"在两条腿上、只有一说实话」.
//
// internal/tools/task.go's stub got its honesty clause in 174-r1 (the extra %s
// slot between "约 %d token" and "，全文见"). internal/agent/spill.go's D15(3)
// stub - the one the working path actually puts into the model's context - was
// still printing a path it never asked anybody about. These criteria pin the
// second leg to the same verdict, on four questions:
//
//	speaks    each of the three shapes that make a pointer un-readable says so,
//	          in the reply text the model reads (not a log line, not a Go error)
//	quiet     a healthy pointer (C26 says in-root) is byte-identical to the
//	          pre-fix sentence, with 注意： appearing zero times - otherwise
//	          "tell the truth" degenerates into spray
//	bounded   the canonicalize-failure shape names the fact and stops: no
//	          upstream error text, no allowlist roots, no C26 internal state
//	one       the notice never adds a second path-shaped token, so the pointer
//	          regex ticket 177's exemption window is keyed to still matches
//	          exactly once, and this leg cannot change that text's shape
//
// Plus the wiring pair, because a field nobody reads passes all of the above
// and still tells the model nothing on a real run: Config.PointerJudge must
// reach the very Spiller whose Prepare produced the text that reached history.

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// spill174NoticeLead is what "this leg is speaking" looks like.
const spill174NoticeLead = "注意："

// spill174OldTemplate is the sentence as it stood BEFORE this ticket, spelled
// out here on purpose: an expectation re-derived from the value under test is
// self-echo, not a ruler (174-v1's finding about the tools leg's pin #8).
const spill174OldTemplate = "%s\n[…输出已落文件：省略 %d 字符，总长 %d 字节 / 约 %d token，全文见 %s…]\n%s"

// spill174PointerRe is the pointer the host wrote. Local on purpose: this
// package does not import internal/tools, and sharing that regex would be the
// cross-package coupling ticket 174's ruling tells me not to create.
var spill174PointerRe = regexp.MustCompile(`全文见 (\S+)…`)

// fake174Judge is a C26 stand-in that only ever answers the two questions
// PointerJudge asks. roots and internalState are strings a real judge could
// leak; the criteria below demand they never reach the reply.
type fake174Judge struct {
	canonCalls  int
	allowCalls  int
	lastCanonIn string
	canonOut    string
	canonErr    error
	inRoot      bool

	// quoted back only inside this struct - the reverse ruler asserts neither
	// shows up in the model-visible text.
	roots         []string
	internalState string
}

func (f *fake174Judge) Canonicalize(raw string) (string, error) {
	f.canonCalls++
	f.lastCanonIn = raw
	if f.canonErr != nil {
		return "", fmt.Errorf("%s: %w", f.internalState, errors.New(strings.Join(f.roots, "|")))
	}
	if f.canonOut != "" {
		return f.canonOut, nil
	}
	return raw, nil
}

func (f *fake174Judge) InAllowlist(canonical string) bool {
	f.allowCalls++
	return f.inRoot
}

// spill174OverThreshold spills one payload with an optionally wired judge and
// returns the stub plus the pieces the caller needs to check numbers from the
// outside.
func spill174OverThreshold(t *testing.T, judge PointerJudge) (Spill, string, Budgets) {
	t.Helper()
	b := BudgetsFor(128000)
	dir := sealableTempDir124(t)
	sp := NewSpiller(dir, b)
	if judge != nil {
		sp.WithPointerJudge(judge)
	}
	head := "HEADSTART-"
	tail := "-TAILFINISH"
	full := head + strings.Repeat("0123456789abcdef ", b.SpillTokens) + tail
	out, err := sp.Prepare("call_honesty", full)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if !out.Spilled {
		t.Fatalf("precondition: payload under the spill threshold (tokens %d vs %d)",
			ApproxTokens(full), b.SpillTokens)
	}
	return out, full, b
}

// 1. speaks: unwired judge must NOT be read as "readable".

func TestSpillPointerUnwiredJudgeFailsClosed174(t *testing.T) {
	out, _, _ := spill174OverThreshold(t, nil)

	if !strings.Contains(out.Text, "授权判定者未接线") {
		t.Errorf("unwired judge must say so in the reply, got %q", out.Text)
	}
	if !strings.Contains(out.Text, "无法核实") {
		t.Errorf("fail-closed wording must not claim the path is fine: %q", out.Text)
	}
	if !strings.Contains(out.Text, "按读不到处理") {
		t.Errorf("fail-closed must land on the strict side, not on silence: %q", out.Text)
	}
	// The pointer itself stays (PLAN.md:431's shape, and 164's accepted AC#3):
	// honesty is an added clause, never a withdrawn promise.
	if m := spill174PointerRe.FindStringSubmatch(out.Text); m == nil || m[1] != out.Path {
		t.Errorf("pointer must survive the notice, got %q", out.Text)
	}
}

// 2. speaks: out-of-root names the gap and both roads back.

func TestSpillPointerOutsideAllowlistNamesBothRoadsBack174(t *testing.T) {
	j := &fake174Judge{inRoot: false, roots: []string{`C:\authorized\only`}, internalState: "c26-state"}
	out, _, _ := spill174OverThreshold(t, j)

	for _, want := range []string{
		"这条路径现在读不到",
		"不在你被授权的目录范围内",
		"[fs] allowed_dirs", // road one
		"L2 卡",              // road two: 174-c2 measured allowlist is a tier input, not a wall
	} {
		if !strings.Contains(out.Text, want) {
			t.Errorf("out-of-root notice must contain %q, got %q", want, out.Text)
		}
	}
	// Judged through C26, not by hand: both methods, on the printed path.
	if j.canonCalls != 1 || j.lastCanonIn != out.Path {
		t.Errorf("verdict must canonicalize the printed path once, got %d calls, last=%q",
			j.canonCalls, j.lastCanonIn)
	}
	if j.allowCalls != 1 {
		t.Errorf("verdict must ask InAllowlist once, got %d", j.allowCalls)
	}
	// And it must not quote the answer's innards back at the model.
	for _, banned := range []string{`C:\authorized\only`, "c26-state"} {
		if strings.Contains(out.Text, banned) {
			t.Errorf("notice leaked judge internals %q into model-visible text", banned)
		}
	}
}

// 3. speaks + bounded: canonicalize failure says the fact, quotes no upstream text.

func TestSpillPointerCanonicalizeFailureQuotesNoUpstreamText174(t *testing.T) {
	j := &fake174Judge{
		roots:         []string{`D:\leaked\root`, out174Sentinel},
		internalState: "REPARSE_DENIED_INTERNAL",
	}
	j.canonErr = errors.New("boom")
	out, _, _ := spill174OverThreshold(t, j)

	if !strings.Contains(out.Text, "连规范化都没通过") {
		t.Errorf("a path C26 cannot canonicalize must be announced as un-readable: %q", out.Text)
	}
	if j.allowCalls != 0 {
		t.Errorf("InAllowlist must not run on an uncanonicalized path, got %d calls", j.allowCalls)
	}
	// The half AC#2d says the tools leg is still missing, armed here from day
	// one: neither the error's own words nor anything the judge could otherwise
	// route into the sentence may reach the model.
	for _, banned := range []string{"boom", out174Sentinel, `D:\leaked\root`, "REPARSE_DENIED_INTERNAL"} {
		if strings.Contains(out.Text, banned) {
			t.Errorf("notice spliced %q out of the judge into model-visible text", banned)
		}
	}
	if strings.Count(out.Text, spill174NoticeLead) != 1 {
		t.Errorf("exactly one notice expected, got %q", out.Text)
	}
}

// out174Sentinel is a stand-in for a path-shaped upstream detail.
const out174Sentinel = `\\vserver\share\secret`

// 4. quiet: the healthy positive stays silent, to the byte.

func TestSpillHealthyPointerStillMatchesThePreFixSentence174(t *testing.T) {
	j := &fake174Judge{inRoot: true}
	out, full, b := spill174OverThreshold(t, j)

	if strings.Contains(out.Text, spill174NoticeLead) {
		t.Errorf("a readable pointer must not be decorated, got %q", out.Text)
	}
	// Independent reconstruction of the pre-fix sentence: head/tail re-cut from
	// the ORIGINAL payload with the same budget, not sliced back out of the
	// reply under test.
	head := takeTokens(full, b.SpillHeadTokens)
	tail := takeTokensLast(full, b.SpillTailTokens)
	want := fmt.Sprintf(spill174OldTemplate,
		head, len(full)-len(head)-len(tail), len(full), ApproxTokens(full), out.Path, tail)
	if out.Text != want {
		t.Errorf("healthy pointer drifted from the pre-fix sentence:\n got %q\nwant %q",
			out.Text, want)
	}
}

// 5. one: the notice never grows a second path, so ticket 177's exemption
// window (keyed to the single host-written pointer) keeps the same text shape.

func TestSpillPointerNoticeAddsNoSecondPath174(t *testing.T) {
	cases := map[string]PointerJudge{
		"unwired":  nil,
		"outRoot":  &fake174Judge{inRoot: false},
		"canonErr": &fake174Judge{canonErr: errors.New("x"), internalState: "s", roots: []string{"r"}},
	}
	for name, j := range cases {
		t.Run(name, func(t *testing.T) {
			out, _, _ := spill174OverThreshold(t, j)
			if got := strings.Count(out.Text, "全文见"); got != 1 {
				t.Fatalf("%s: 全文见 appears %d times, want exactly 1 (177 keys its exemption window to this)",
					name, got)
			}
			m := spill174PointerRe.FindStringSubmatch(out.Text)
			if m == nil || m[1] != out.Path {
				t.Fatalf("%s: pointer unparsable or wrong: %q", name, out.Text)
			}
			notice := out.Text[strings.Index(out.Text, " token")+len(" token") : strings.Index(out.Text, "，全文见")]
			if strings.ContainsAny(notice, `/\`) {
				t.Errorf("%s: notice carries path characters, which is a second pointer in all but name: %q",
					name, notice)
			}
		})
	}
}

// 6. the clause sits inside the bracket: the head stays the head, the tail
// stays the tail, and the stub is still smaller than what it replaced.

func TestSpillPointerNoticeKeepsHeadTailAndBudget174(t *testing.T) {
	j := &fake174Judge{inRoot: false}
	out, full, _ := spill174OverThreshold(t, j)

	if !strings.HasPrefix(out.Text, "HEADSTART-") {
		t.Errorf("notice displaced the head: %q", out.Text[:min(60, len(out.Text))])
	}
	if !strings.HasSuffix(out.Text, "-TAILFINISH") {
		t.Errorf("notice must sit before the pointer, not after the tail: %q",
			out.Text[max(0, len(out.Text)-60):])
	}
	if ApproxTokens(out.Text) >= ApproxTokens(full) {
		t.Errorf("stub (%d tokens) is not smaller than the original (%d)",
			ApproxTokens(out.Text), ApproxTokens(full))
	}
	// The truncation marker and the notice coexist: 174's clause must not
	// crowd out the raw-cap announcement spill_test.go already pins.
	sp := NewSpiller(sealableTempDir124(t), Budgets{
		RawOutputCapBytes: 200, SpillTokens: 40, SpillHeadTokens: 10, SpillTailTokens: 10,
	}).WithPointerJudge(j)
	s2, err := sp.Prepare("call_both", strings.Repeat("x", 800))
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	for _, want := range []string{"truncated=true", "硬上限", "不在你被授权的目录范围内"} {
		if !strings.Contains(s2.Text, want) {
			t.Errorf("stub lost %q when both announcements apply:\n%s", want, s2.Text)
		}
	}
}

// 7/8. wiring pair: Config.PointerJudge must reach the Spiller whose output
// reaches the model's history. Both arms, so a dead field cannot pass.

func TestLoopSpillNoticeReachesHistoryUnwired174(t *testing.T) {
	dir := sealableTempDir124(t)
	h := newHarness(t, "spill-tool", withConfig(func(c *Config) {
		c.ContextWindow = 4096
		c.ArtifactsDir = filepath.Join(dir, "artifacts")
		c.PerToolTimeout = 2 * time.Second
	}))
	res := h.run("把长结果整理一下")
	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s)", res.Status, res.Message)
	}
	if len(res.ToolLog) != 1 || !res.ToolLog[0].Spilled {
		t.Fatalf("precondition: tiny window must spill, got %+v", res.ToolLog)
	}
	// Landing: the text handed to the next round is the stub WITH the clause.
	if !strings.Contains(res.ToolLog[0].Text, "授权判定者未接线") {
		t.Errorf("tool log carries a pointer with no verdict attached: %q", res.ToolLog[0].Text)
	}
	// Delivery: the clause is in what the provider sees next round, not just in
	// the returned struct.
	if !historyHasResultContaining(h.loop.History(), "授权判定者未接线") {
		t.Errorf("round-2 history lost the honesty clause:\n%s", dumpHistory(h.loop.History()))
	}
	// The artifact really is on disk and really is the thing pointed at, so the
	// clause is about reachability, never about existence (who lands it).
	m := spill174PointerRe.FindStringSubmatch(res.ToolLog[0].Text)
	if m == nil {
		t.Fatalf("no pointer in %q", res.ToolLog[0].Text)
	}
	if body, err := os.ReadFile(m[1]); err != nil || len(body) == 0 {
		t.Errorf("pointed artifact missing/empty (%v)", err)
	}
}

func TestLoopHonorsConfigPointerJudgeQuietly174(t *testing.T) {
	dir := sealableTempDir124(t)
	j := &fake174Judge{inRoot: true}
	h := newHarness(t, "spill-tool", withConfig(func(c *Config) {
		c.ContextWindow = 4096
		c.ArtifactsDir = filepath.Join(dir, "artifacts")
		c.PerToolTimeout = 2 * time.Second
		c.PointerJudge = j
	}))
	res := h.run("把长结果整理一下")
	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s)", res.Status, res.Message)
	}
	if len(res.ToolLog) != 1 || !res.ToolLog[0].Spilled {
		t.Fatalf("precondition: tiny window must spill, got %+v", res.ToolLog)
	}
	if j.canonCalls == 0 {
		t.Fatalf("Config.PointerJudge never consulted: the field is dead wiring")
	}
	if strings.Contains(res.ToolLog[0].Text, spill174NoticeLead) {
		t.Errorf("wired+in-root must stay quiet, got %q", res.ToolLog[0].Text)
	}
	if !historyHasResultContaining(h.loop.History(), "输出已落文件") {
		t.Errorf("round-2 history lost the stub:\n%s", dumpHistory(h.loop.History()))
	}
}
