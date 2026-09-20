package agent

import (
	"context"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Acceptance criterion 5: the C22 gradient ladder [3,5,8] injects graduated
// reminders, the top rung drives Stuck with an explicit user-visible message
// (never a silent stop), and the per-tool timeout fires.

func TestLoopGuardLadderRemindersThenStuck(t *testing.T) {
	h := newHarness(t, "repeat-echo")
	res := h.run("把目录再列一遍")

	if res.Status != StatusStuck {
		t.Fatalf("status = %s (%s), want stuck", res.Status, res.Message)
	}
	if res.Brake != BrakeRepeat {
		t.Errorf("brake = %s, want %s", res.Brake, BrakeRepeat)
	}
	// Rungs 3 and 5 inject reminders; rung 8 stops. The ninth response (a text
	// answer) must never be requested.
	want := []int{3, 5, 8}
	if len(res.ReminderLevels) != len(want) {
		t.Fatalf("reminder levels = %v, want %v", res.ReminderLevels, want)
	}
	for i, lv := range want {
		if res.ReminderLevels[i] != lv {
			t.Errorf("level %d = %d, want %d", i, res.ReminderLevels[i], lv)
		}
	}
	if h.requests() != 8 {
		t.Errorf("requests = %d, want 8 (stopped at the top rung)", h.requests())
	}
	// Rounds 1..7 executed; the 8th was refused before execution.
	if n := h.echo().CallCount(); n != 7 {
		t.Errorf("executions = %d, want 7", n)
	}

	// The reminders really were injected into the conversation (D43 row #19
	// side effect "agent.inject-reminder"), twice (rungs 3 and 5).
	var injected int
	for _, m := range h.loop.History() {
		for _, p := range m.Content {
			if t, ok := p.(llmText); ok && strings.HasPrefix(t.Text, "[系统提醒]") {
				injected++
			}
		}
	}
	if injected != 2 {
		t.Errorf("injected reminders = %d, want 2:\n%s", injected, dumpHistory(h.loop.History()))
	}

	// Stuck is visible AND names what repeated (C22 forbids a silent stop).
	st := h.sink.LastOf(EvStuck)
	if st == nil {
		t.Fatal("no stuck event published")
	}
	for _, needle := range []string{"echo", "连续调用 8 次", "参数相同", "token"} {
		if !strings.Contains(st.Text, needle) && !strings.Contains(res.Message, needle) {
			t.Errorf("stuck message must contain %q; got event=%q result=%q", needle, st.Text, res.Message)
		}
	}
	recorded := h.sink.Events()
	if len(recorded) == 0 {
		t.Fatal("sink recorded nothing")
	}
	reminders := 0
	for _, e := range recorded {
		if e.Kind == EvReminder {
			reminders++
			if e.RepeatLevel == 0 || e.ToolName != "echo" {
				t.Errorf("reminder event = %+v", e)
			}
		}
	}
	if reminders != 3 {
		t.Errorf("reminder events = %d, want 3 (one per rung)", reminders)
	}
}

// The ladder is config data, not a constant: a smaller configured ladder stops
// earlier, which is what proves the thresholds are read from config.
func TestLoopGuardLadderIsConfigurable(t *testing.T) {
	h := newHarness(t, "repeat-echo", withConfig(func(c *Config) {
		c.RepeatThresholds = []int{2}
	}))
	res := h.run("把目录再列一遍")
	if res.Status != StatusStuck {
		t.Fatalf("status = %s (%s), want stuck", res.Status, res.Message)
	}
	if h.requests() != 2 {
		t.Errorf("requests = %d, want 2 (single-rung ladder fires on the 2nd repeat)", h.requests())
	}
	if !strings.Contains(res.Message, "连续调用 2 次") {
		t.Errorf("message = %q", res.Message)
	}
}

// A changed call sequence resets the counter (the brake targets *consecutive*
// identical calls, D43 row #19 wording).
func TestLoopGuardResetsOnDifferentCall(t *testing.T) {
	g := NewGuard(GuardConfig{Budgets: BudgetsFor(128000)})
	one := toolCalls("call_a", "echo", `{"text":"x"}`)
	for i := 1; i <= 2; i++ {
		if _, rep := g.ObserveTurn(one); rep && i < 3 {
			t.Fatalf("rung fired early at %d", i)
		}
	}
	other := toolCalls("call_b", "echo", `{"text":"y"}`)
	if _, rep := g.ObserveTurn(other); rep {
		t.Fatal("a different argument set must reset the repeat counter")
	}
	if rem, _ := g.ObserveTurn(other); rem.Level != 0 {
		t.Fatalf("first repeat of the new call must not fire, got %+v", rem)
	}
	if rem, _ := g.ObserveTurn(one); rem.Level != 0 {
		t.Fatalf("back to the old call must restart at 1, got %+v", rem)
	}
}

func TestLoopGuardTokenBudgetScalesWithWindow(t *testing.T) {
	big := BudgetsFor(128000)
	if big.TokenBudget != 200000 {
		t.Errorf("128k budget = %d, want the frozen 200k", big.TokenBudget)
	}
	small := BudgetsFor(4096)
	if small.TokenBudget >= big.TokenBudget {
		t.Errorf("4k-window budget = %d, must shrink below %d", small.TokenBudget, big.TokenBudget)
	}
	if got := (4096 * 200000) / 128000; small.TokenBudget != got {
		t.Errorf("4k-window budget = %d, want the proportional %d", small.TokenBudget, got)
	}
	g := NewGuard(GuardConfig{Budgets: small})
	if g.TokenBudget() != small.TokenBudget {
		t.Errorf("guard budget = %d, want %d", g.TokenBudget(), small.TokenBudget)
	}
}

// Per-tool timeout: the call is aborted cooperatively at timeoutMs, the row is
// booked with a tool class, and the loop keeps going (the model is told).
func TestPerToolTimeoutFires(t *testing.T) {
	dir := t.TempDir()
	store, err := memory.Open(dir, memory.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	defer store.Close()

	h := newHarness(t, "slow-tool", withJournal(store),
		withConfig(func(c *Config) {
			c.PerToolTimeout = 60 * time.Millisecond
			c.ArtifactsDir = filepath.Join(dir, "artifacts")
		}))
	started := time.Now()
	res := h.run("慢一点也行")
	elapsed := time.Since(started)

	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s), want completed", res.Status, res.Message)
	}
	if elapsed > 2*time.Second {
		t.Errorf("loop waited %v: the tool was not aborted at 60ms", elapsed)
	}
	if len(res.ToolLog) != 1 {
		t.Fatalf("tool log = %+v", res.ToolLog)
	}
	l := res.ToolLog[0]
	if l.Outcome != OutcomeError || !strings.Contains(l.Text, "超时") {
		t.Errorf("timeout not reported to the model: %+v", l)
	}
	if l.ErrorClass != string(observe.ClassTool) {
		t.Errorf("error_class = %q, want tool (self-correction class, D37)", l.ErrorClass)
	}
	// The result that fed round 2 carries the timeout message.
	if !historyHasResultContaining(h.loop.History(), "超时") {
		t.Errorf("round-2 history lacks the timeout result:\n%s", dumpHistory(h.loop.History()))
	}
	rows, err := store.ListToolCallsByTask(context.Background(), res.TaskID)
	if err != nil || len(rows) != 1 {
		t.Fatalf("tool_call rows = %+v (%v)", rows, err)
	}
	if rows[0].Outcome != OutcomeError {
		t.Errorf("row outcome = %q, want error", rows[0].Outcome)
	}
	// The second turn answered with text, proving the loop continued.
	if res.Text != "那一步超时了，我先给结论。" {
		t.Errorf("final text = %q", res.Text)
	}
}

// ---------------------------------------------------------------------------

// MAJOR-2 / MINOR-4: the timeout above is only proven for a provider that
// BREAKS its own contract. The ToolProvider contract (tools.go:82-85) says a
// host failure comes back as an error, and a real host bridge (ticket 20+) will
// therefore return ctx.Err() on a timed-out call. dispatch used to pass that
// error alongside the structured TOOL_TIMEOUT, and executeCalls prefers a
// non-nil error - so the model saw "context deadline exceeded" and the row was
// booked error_class="internal". This is the assertion on the EMITTED
// error_class for the contract-honest branch.
func TestPerToolTimeoutOfContractHonestToolIsToolClass(t *testing.T) {
	dir := t.TempDir()
	store, err := memory.Open(dir, memory.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	defer store.Close()

	tools := newBlockingProvider() // aborts by RETURNING ctx.Err(), as contracted
	t.Cleanup(tools.release)
	h := newHarness(t, "slow-tool", withJournal(store), withTools(tools),
		withConfig(func(c *Config) {
			c.PerToolTimeout = 60 * time.Millisecond
			c.ArtifactsDir = filepath.Join(dir, "artifacts")
		}))
	res := h.run("慢一点也行")

	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s), want completed (a timeout is not a task failure)",
			res.Status, res.Message)
	}
	if len(res.ToolLog) != 1 {
		t.Fatalf("tool log = %+v", res.ToolLog)
	}
	l := res.ToolLog[0]
	if l.Outcome != OutcomeError {
		t.Errorf("outcome = %q, want %s", l.Outcome, OutcomeError)
	}
	// THE decisive assertion: the tool's own error must not steal the C22 class.
	if l.ErrorClass != string(observe.ClassTool) {
		t.Errorf("error_class = %q, want %s (the host error must not downgrade a "+
			"timeout to an unself-correctable class, D37)", l.ErrorClass, observe.ClassTool)
	}
	if !strings.Contains(l.Text, "超时") || !strings.Contains(l.Text, "sleep") {
		t.Errorf("model-visible text = %q, want the structured TOOL_TIMEOUT naming "+
			"the tool", l.Text)
	}
	if strings.Contains(l.Text, "context deadline exceeded") {
		t.Errorf("the raw host error leaked into the model-visible text: %q", l.Text)
	}
	// What actually fed round 2 is the timeout line, not the ctx error.
	if !historyHasResultContaining(h.loop.History(), "超时") {
		t.Errorf("round-2 history lacks the timeout result:\n%s", dumpHistory(h.loop.History()))
	}
	if strings.Contains(dumpHistory(h.loop.History()), "context deadline exceeded") {
		t.Errorf("conversation carries the raw host error instead of TOOL_TIMEOUT:\n%s",
			dumpHistory(h.loop.History()))
	}
	rows, err := store.ListToolCallsByTask(context.Background(), res.TaskID)
	if err != nil || len(rows) != 1 {
		t.Fatalf("tool_call rows = %+v (%v)", rows, err)
	}
	if rows[0].ErrorClass != string(observe.ClassTool) {
		t.Errorf("persisted error_class = %q, want %s", rows[0].ErrorClass, observe.ClassTool)
	}
	if rows[0].Outcome != OutcomeError {
		t.Errorf("persisted outcome = %q, want %s", rows[0].Outcome, OutcomeError)
	}
	// The loop kept going and the second round answered.
	if res.Text != "那一步超时了，我先给结论。" {
		t.Errorf("final text = %q", res.Text)
	}
}

// MINOR-3: turnSignature's comment promised encoding/json canonicalization while
// the code only trimmed, so equivalent calls written with a different key order
// or spacing escaped the brake entirely (a comment reading stronger than the
// defence). The code now matches the comment; these are both halves of that.
func TestLoopGuardCanonicalizesEquivalentArgs(t *testing.T) {
	g := NewGuard(GuardConfig{Budgets: BudgetsFor(128000)})
	// Four spellings of ONE object, including a whitespace/newline variant.
	variants := []string{
		`{"a":1,"b":2}`,
		`{ "b":2, "a":1 }`,
		"{\n  \"b\": 2,\n  \"a\": 1\n}",
		`{"b":2,"a":1}`,
	}
	var fired []int
	for i, v := range variants {
		rem, rep := g.ObserveTurn(toolCalls("call_sig", "echo", v))
		if !rep {
			continue
		}
		fired = append(fired, rem.Level)
		if rem.Repeat != i+1 {
			t.Errorf("turn %d (%s): repeat counter = %d, want %d - the variants must "+
				"fold onto one signature", i, v, rem.Repeat, i+1)
		}
	}
	if len(fired) != 1 || fired[0] != 3 {
		t.Errorf("rung fired %v across four spellings of one object, want exactly [3] "+
			"(key order/whitespace variants are the same call)", fired)
	}
}

// The canonicalization must not collapse DIFFERENT calls into one signature:
// that would stop a legitimately retrying model. Values, key sets and number
// literals all have to stay distinguishable (json.Number is what keeps a big id
// from being rounded into its neighbour).
func TestLoopGuardKeepsDistinctArgsDistinct(t *testing.T) {
	cases := [][2]string{
		{`{"a":1}`, `{"a":2}`},
		{`{"a":1,"b":2}`, `{"a":1}`},
		{`{"id":10000000000000000000001}`, `{"id":10000000000000000000002}`},
		{`not-json`, `not json either`},
	}
	for n, pair := range cases {
		g := NewGuard(GuardConfig{Budgets: BudgetsFor(128000)})
		for i := 0; i < 8; i++ {
			v := pair[0]
			if i%2 == 1 {
				v = pair[1]
			}
			if rem, rep := g.ObserveTurn(toolCalls("call_d", "echo", v)); rep {
				t.Errorf("case %d (%q vs %q) fired rung %d at turn %d: distinct calls "+
					"were treated as repeats", n, pair[0], pair[1], rem.Level, i)
				break
			}
		}
	}
}

// Non-JSON arguments have no structure to canonicalize, so exact trimmed text is
// the only available comparison - and it must still drive the ladder.
func TestLoopGuardRepeatsNonJSONArgs(t *testing.T) {
	g := NewGuard(GuardConfig{Budgets: BudgetsFor(128000)})
	for i := 0; i < 2; i++ {
		if _, rep := g.ObserveTurn(toolCalls("call_n", "echo", `raw-argument-text`)); rep {
			t.Fatalf("rung fired early at %d", i)
		}
	}
	rem, rep := g.ObserveTurn(toolCalls("call_n", "echo", `raw-argument-text`))
	if !rep || rem.Level != 3 {
		t.Errorf("identical non-JSON args not counted as repeats: rep=%v rem=%+v", rep, rem)
	}
}

// ---------------------------------------------------------------------------

func historyHasResultContaining(hist []llmMessage, needle string) bool {
	for _, m := range hist {
		if m.Role != llmRoleTool {
			continue
		}
		for _, p := range m.Content {
			tr, ok := p.(llmToolResult)
			if !ok {
				continue
			}
			for _, ip := range tr.Content {
				if t, ok := ip.(llmText); ok && strings.Contains(t.Text, needle) {
					return true
				}
			}
		}
	}
	return false
}
