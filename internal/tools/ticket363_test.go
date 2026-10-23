package tools_test

// Ticket 363: the audit line must not claim a question that was never asked.
//
// The chain the ticket measured, in one call:
//
//	bridge.go:311  sil := mode.Screen(verdict)             L0 -> "unidentified level"
//	bridge.go:323  case sil.Kept != "":                    Kept is non-empty for L0
//	bridge.go:324  "MODE-REDLINE ... still asks ..."        THE FALSE STATEMENT
//	bridge.go:432  case risk.L0: decision = allow          no gate, no card
//
// So the log said "still asks" on a call that raised zero cards and had already
// returned the file's contents. The same sentence was false in the opposite
// direction for a Deny verdict: refused outright, zero cards, "still asks".
//
// Two fixes are under test here, and each has its own failing assertion:
//
//	Screen recognises L0 (internal/risk/mode.go)   -> Kept stays empty, so the
//	                                                  bridge has nothing to claim
//	the kept-question line names the real outcome  -> a Deny says "refuses"
//
// The cards are counted where the cards are: the Gate interface is exactly two
// callbacks (gate.go:110-117, PendingWindow = the L1 window, PendingApproval =
// the C18 queue), so this file counts CALLS TO THOSE TWO and never a UI tick.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// --- fixtures ---------------------------------------------------------------

// t363Gate counts the two cards and answers "the user let it through", so an
// assertion about a card being raised is never also an assertion about refusal.
type t363Gate struct {
	mu        sync.Mutex
	windows   int
	approvals int
	seen      []tools.Decision
}

func (g *t363Gate) PendingWindow(_ context.Context, d tools.Decision) (tools.Answer, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.windows++
	g.seen = append(g.seen, d)
	return tools.AnswerAllow, ""
}

func (g *t363Gate) PendingApproval(_ context.Context, d tools.Decision) (tools.Answer, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.approvals++
	g.seen = append(g.seen, d)
	return tools.AnswerAllow, ""
}

func (g *t363Gate) counts() (int, int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.windows, g.approvals
}

// t363Modes is the read-only ModeSource the bridge reads once per call.
type t363Modes struct{ m risk.Mode }

func (s t363Modes) PermissionMode() risk.Mode { return s.m }

// t363Log is the audit sink (production: rt.auditf), verbatim.
type t363Log struct {
	mu    sync.Mutex
	lines []string
}

func (l *t363Log) logf(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lines = append(l.lines, fmt.Sprintf(format, args...))
}

func (l *t363Log) all() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.lines...)
}

type t363Harness struct {
	b    *tools.Bridge
	gate *t363Gate
	log  *t363Log

	mu   sync.Mutex
	decs []tools.Decision
}

func (h *t363Harness) record(d tools.Decision) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.decs = append(h.decs, d)
}

func (h *t363Harness) last() tools.Decision {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.decs) == 0 {
		return tools.Decision{}
	}
	return h.decs[len(h.decs)-1]
}

// t363New composes the SAME shape cmd/wisp composes (run.go:764-788): the real
// D34 fs roster, the real C26 canonicalizer, the real C19 assessor built by
// tools.New from those two, the injected gate and the injected mode source. The
// only substitution is the audit sink.
func t363New(t *testing.T, mode risk.Mode, allowed string) *t363Harness {
	t.Helper()
	paths := tools.NewPathCanonicalizer([]string{allowed}, nil)
	reg := tools.NewRegistry()
	for _, e := range tools.BuiltinFSEntries(tools.FSDeps{Paths: paths}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	h := &t363Harness{gate: &t363Gate{}, log: &t363Log{}}
	h.b = tools.New(tools.Options{
		Registry:   reg,
		Paths:      paths,
		Gate:       h.gate,
		Modes:      t363Modes{m: mode},
		Logf:       h.log.logf,
		OnDecision: h.record,
	})
	return h
}

func t363Exec(t *testing.T, h *t363Harness, tool, args string) (agent.ToolOutcome, error) {
	t.Helper()
	return h.b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "task-363", CorrelationID: "corr-363", CallID: "call-363",
		Name: tool, Args: json.RawMessage(args),
	})
}

func t363Canon(t *testing.T, p string) string {
	t.Helper()
	c, err := tools.NewPathCanonicalizer(nil, nil).Canonicalize(p)
	if err != nil {
		t.Fatalf("Canonicalize(%q): %v", p, err)
	}
	return c
}

// t363AskClaimsNeedACard is the invariant behind both fixes: "still asks" is a
// claim about a CARD, so it is only true when the routing really raised one.
func t363AskClaimsNeedACard(t *testing.T, label string, lines []string, windows, approvals int) {
	t.Helper()
	for _, l := range lines {
		if !strings.Contains(l, "still asks") {
			continue
		}
		if windows+approvals == 0 {
			t.Errorf("%s: the audit line claims a question while the call raised 0 cards "+
				"(L1 windows=%d, L2 approvals=%d):\n  %s", label, windows, approvals, l)
		}
	}
}

func t363LineWith(lines []string, needle string) (string, bool) {
	for _, l := range lines {
		if strings.Contains(l, needle) {
			return l, true
		}
	}
	return "", false
}

// --- (1) L0: the direction the ticket names first ---------------------------

// TestTicket363L0ReadNeverClaimsAQuestion is the ticket's core reading turned
// into a nail. Same call, three modes: an ordinary read-only fs.read inside the
// allowed directory. Today the bridge stamps "still asks" plus "无法识别的风险级"
// on it, sends ModeKept with the same sentence to the tool_call row, and raises
// zero cards while the tool returns the file's contents.
func TestTicket363L0ReadNeverClaimsAQuestion(t *testing.T) {
	root := t.TempDir()
	allowed := t363Canon(t, root)
	target := filepath.Join(root, "note.txt")
	if err := os.WriteFile(target, []byte("hello"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	for _, m := range []risk.Mode{risk.ModeAskEveryStep, risk.ModeAskHighRisk, risk.ModeAutoApprove} {
		h := t363New(t, m, allowed)
		out, err := t363Exec(t, h, "fs.read", fmt.Sprintf(`{"path":%q}`, filepath.ToSlash(target)))
		if err != nil {
			t.Fatalf("mode=%s: Execute: %v", m, err)
		}
		w, a := h.gate.counts()

		// The SAME CALL's facts: the tool ran and no card was raised. These two
		// are what makes the old log line a lie rather than a wording quibble.
		if out.IsError || out.Text != "hello" {
			t.Errorf("mode=%s: the L0 read did not return the file: isError=%v text=%q",
				m, out.IsError, out.Text)
		}
		if w != 0 || a != 0 {
			t.Errorf("mode=%s: an L0 read raised L1 windows=%d L2 approvals=%d, want 0/0 "+
				"(the L0 route never reaches a gate)", m, w, a)
		}

		// The record the tool_call row is built from.
		d := h.last()
		if d.Level != risk.L0 {
			t.Errorf("mode=%s: booked level = %v, want the assessed L0", m, d.Level)
		}
		if d.Mode != m {
			t.Errorf("mode=%s: the decision records mode %v", m, d.Mode)
		}
		if d.ModeKept != "" {
			t.Errorf("mode=%s: ModeKept = %q on an L0 read; Kept names the red line that "+
				"kept a question (risk/mode.go:143-146) and L0 had no question to keep - "+
				"every read-only row would carry a red line that does not exist", m, d.ModeKept)
		}
		if d.ModeSilenced {
			t.Errorf("mode=%s: ModeSilenced = true on an L0 read; nothing was asked, so the "+
				"mode removed nothing", m)
		}

		// The operator-visible line, three named ways of being wrong.
		lines := h.log.all()
		for _, l := range lines {
			if strings.Contains(l, "still asks") {
				t.Errorf("mode=%s: the audit says 'still asks' on a call that raised 0 cards "+
					"and already returned %q:\n  %s", m, out.Text, l)
			}
			if strings.Contains(l, "无法识别的风险级") {
				t.Errorf("mode=%s: the audit calls L0 an unrecognised level; L0 is the judge's "+
					"own zero value (risk/assessor.go:46) and Screen itself emits it:\n  %s", m, l)
			}
			if strings.Contains(l, "MODE-SILENCE") || strings.Contains(l, "MODE-REDLINE") {
				t.Errorf("mode=%s: an L0 verdict is neither a silenced question nor a kept "+
					"one, so neither mode line may be printed for it:\n  %s", m, l)
			}
		}
		t363AskClaimsNeedACard(t, "mode="+m.String(), lines, w, a)

		// ...and the line that does tell the truth is still there.
		if _, ok := t363LineWith(lines, "tool=fs.read risk=L0 decision=allow outcome=success"); !ok {
			t.Errorf("mode=%s: no audit line books the call as decision=allow with the "+
				"assessed risk=L0; lines:\n  %s", m, strings.Join(lines, "\n  "))
		}
	}
}

// --- (2) Deny: the reverse direction ---------------------------------------

// TestTicket363DenyIsRefusedNotAsked is the half Screen cannot fix: a tier-A
// path is refused by route's `case risk.Deny` without touching the gate
// (bridge.go:437-439), so "still asks" is false the other way round. The fix is
// in the log line: it must name the refusal, not a question.
func TestTicket363DenyIsRefusedNotAsked(t *testing.T) {
	// A profile this test owns, so the frozen %APPDATA%\wisp\config.toml anchor
	// names a file here and not the machine's real config (same trick as
	// ticket90's tier-A test).
	home := t.TempDir()
	appdata := filepath.Join(home, "AppData", "Roaming")
	local := filepath.Join(home, "AppData", "Local")
	cfgDir := filepath.Join(appdata, "wisp")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cfg := filepath.Join(cfgDir, "config.toml")
	if err := os.WriteFile(cfg, []byte("[llm]\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("APPDATA", appdata)
	t.Setenv("LOCALAPPDATA", local)

	allowed := t363Canon(t, t.TempDir())
	for _, m := range []risk.Mode{risk.ModeAskEveryStep, risk.ModeAutoApprove} {
		h := t363New(t, m, allowed)
		out, err := t363Exec(t, h, "fs.read", fmt.Sprintf(`{"path":%q}`, filepath.ToSlash(cfg)))
		if err != nil {
			t.Fatalf("mode=%s: Execute: %v", m, err)
		}
		w, a := h.gate.counts()

		d := h.last()
		if d.Level != risk.Deny {
			t.Fatalf("mode=%s: assessed level = %v, want Deny (tier A is absolute)", m, d.Level)
		}
		if w != 0 || a != 0 {
			t.Errorf("mode=%s: a Deny raised L1 windows=%d L2 approvals=%d, want 0/0 (a Deny "+
				"must not become a question)", m, w, a)
		}
		if !out.IsError || !strings.Contains(out.Text, "该目标被绝对禁止访问") {
			t.Errorf("mode=%s: the Deny came back as isError=%v text=%q, want the absolute-deny "+
				"refusal", m, out.IsError, out.Text)
		}

		lines := h.log.all()
		claim, ok := t363LineWith(lines, "MODE-REDLINE")
		if !ok {
			t.Fatalf("mode=%s: no MODE-REDLINE line naming the red line that kept this verdict; "+
				"lines:\n  %s", m, strings.Join(lines, "\n  "))
		}
		if strings.Contains(claim, "still asks") {
			t.Errorf("mode=%s: the audit says 'still asks' on a Deny; the call was refused "+
				"without a card (0 windows, 0 approvals):\n  %s", m, claim)
		}
		if !strings.Contains(claim, "refuses") {
			t.Errorf("mode=%s: the kept-verdict line does not name the refusal; want it to "+
				"say the call is refused with no card raised:\n  %s", m, claim)
		}
		t363AskClaimsNeedACard(t, "mode="+m.String()+" deny", lines, w, a)

		if _, ok := t363LineWith(lines, "decision=reject"); !ok {
			t.Errorf("mode=%s: no audit line books decision=reject; lines:\n  %s",
				m, strings.Join(lines, "\n  "))
		}
	}
}

// --- (3) positive controls: the levels that already told the truth ----------

// TestTicket363OtherLevelsKeepTellingTheTruth is the guard against "fixed the
// log by making everything quiet". It exercises the two levels Screen has
// always recognised and the one it refuses, and pins BOTH halves: L1 still
// opens its window under the default mode and prints no kept-verdict line at
// all, L2 still opens its approval card and the "still asks" line is TRUE
// there (one card exists), and the silenced L1 still goes through MODE-SILENCE.
func TestTicket363OtherLevelsKeepTellingTheTruth(t *testing.T) {
	t.Run("L1 window under ask_every_step", func(t *testing.T) {
		root := t.TempDir()
		allowed := t363Canon(t, root)
		h := t363New(t, risk.ModeAskEveryStep, allowed)
		target := filepath.Join(root, "fresh.txt")
		out, err := t363Exec(t, h, "fs.write", fmt.Sprintf(`{"path":%q,"content":"x"}`, filepath.ToSlash(target)))
		if err != nil {
			t.Fatalf("Execute: %v", err)
		}
		w, a := h.gate.counts()
		if w != 1 || a != 0 {
			t.Errorf("L1 write under ask_every_step raised windows=%d approvals=%d, want 1/0 "+
				"(the reversible-write window must still open)", w, a)
		}
		if out.IsError || !strings.Contains(out.Text, "已写入") {
			t.Errorf("the approved L1 write did not run: isError=%v text=%q", out.IsError, out.Text)
		}
		lines := h.log.all()
		if l, ok := t363LineWith(lines, "MODE-REDLINE"); ok {
			t.Errorf("an unsilenced L1 verdict printed a kept-verdict line; the card was "+
				"raised, so there is nothing to explain:\n  %s", l)
		}
		if l, ok := t363LineWith(lines, "MODE-SILENCE"); ok {
			t.Errorf("ask_every_step silenced an L1 verdict:\n  %s", l)
		}
	})

	t.Run("L1 silenced by ask_high_risk", func(t *testing.T) {
		root := t.TempDir()
		allowed := t363Canon(t, root)
		h := t363New(t, risk.ModeAskHighRisk, allowed)
		target := filepath.Join(root, "fresh-silent.txt")
		if _, err := t363Exec(t, h, "fs.write", fmt.Sprintf(`{"path":%q,"content":"x"}`, filepath.ToSlash(target))); err != nil {
			t.Fatalf("Execute: %v", err)
		}
		w, a := h.gate.counts()
		if w != 0 || a != 0 {
			t.Errorf("ask_high_risk opened windows=%d approvals=%d on an L1 write, want 0/0", w, a)
		}
		lines := h.log.all()
		if l, ok := t363LineWith(lines, "MODE-SILENCE"); !ok {
			t.Errorf("a silenced L1 verdict left no MODE-SILENCE line; lines:\n  %s", strings.Join(lines, "\n  "))
		} else if !strings.Contains(l, "not an allow") {
			t.Errorf("the MODE-SILENCE line must keep saying a silenced question is not an "+
				"allow:\n  %s", l)
		}
		if l, ok := t363LineWith(lines, "MODE-REDLINE"); ok {
			t.Errorf("a silenced verdict also printed the kept-verdict line:\n  %s", l)
		}
	})

	// L2 twice: kept by ask_high_risk's own rule, and kept by the R1/R8 red line
	// under auto_approve. In both the card exists, so "still asks" is true.
	for _, m := range []risk.Mode{risk.ModeAskHighRisk, risk.ModeAutoApprove} {
		t.Run("L2 overwrite under "+m.String(), func(t *testing.T) {
			root := t.TempDir()
			allowed := t363Canon(t, root)
			over := filepath.Join(root, "over.txt")
			if err := os.WriteFile(over, []byte("ORIGINAL"), 0o600); err != nil {
				t.Fatalf("write: %v", err)
			}
			h := t363New(t, m, allowed)
			out, err := t363Exec(t, h, "fs.write", fmt.Sprintf(`{"path":%q,"content":"NEW"}`, filepath.ToSlash(over)))
			if err != nil {
				t.Fatalf("Execute: %v", err)
			}
			w, a := h.gate.counts()
			if w != 0 || a != 1 {
				t.Errorf("an irreversible overwrite under %s raised windows=%d approvals=%d, "+
					"want 0/1 (any mode must still ask about R8)", m, w, a)
			}
			if out.IsError || !strings.Contains(out.Text, "已写入") {
				t.Errorf("the approved L2 write did not run: isError=%v text=%q", out.IsError, out.Text)
			}
			lines := h.log.all()
			claim, ok := t363LineWith(lines, "still asks")
			if !ok {
				t.Fatalf("the kept L2 verdict lost its audit line; lines:\n  %s",
					strings.Join(lines, "\n  "))
			}
			if !strings.Contains(claim, "because=") {
				t.Errorf("the kept-verdict line no longer names WHY the mode did not apply:\n  %s", claim)
			}
			// The claim is true here: exactly one card was raised for it.
			t363AskClaimsNeedACard(t, "L2/"+m.String(), lines, w, a)
		})
	}
}
