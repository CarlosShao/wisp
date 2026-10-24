//go:build windows

package main

// Ticket 375's rulers at the ASSEMBLY level.
//
// These are the cases that go red on the pre-375 tree BY RUNNING (they use no
// API the change introduces): the console assembly whose runSpec carries no
// reply stream is a run where nobody can answer, and before this ticket an L1
// write in that run executed anyway - either silently under auto_approve
// (375-a1's READ-F: 0 cards, file written) or two seconds later under
// ask_every_step (READ-S: 1 card, the window's expiry read as consent).
//
// The two positive controls at the bottom are the other half of the claim: with
// an answer source present, L1 and L2 do exactly what they did before.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/secret"
	"github.com/CarlosShao/wisp/internal/tools"
)

// t375Config writes a loadable config.toml carrying one permission mode. It is
// this file's own fixture rather than run_test.go's newRunFixture because the
// mode is the variable under test and that fixture hard-codes the section file.
func t375Config(t *testing.T, dir, mode string) {
	t.Helper()
	srv := adaptertest.StartMockllm(t)
	body := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/m1"]

[llm.retry]
max = 1
backoff_ms = 1

[llm.providers.acme]
protocol = "openai-chat"
base_url = %q
api_key_ref = "dpapi:acme"

[llm.providers.acme.models.m1]
enabled = true
context_window = 128000

[fs]
allowed_dirs = [%q]

[risk]
l1_window_sec = 1
confirm_timeout_sec = 31
permission_mode = %q
`, srv.Base+"/v1", filepath.ToSlash(dir), mode)
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := secret.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Store("dpapi:acme", fakeStoreKey); err != nil {
		t.Fatal(err)
	}
}

func t375Call(task, corr, path string) agent.ToolRequest {
	return agent.ToolRequest{
		TaskID:        task,
		CorrelationID: corr,
		CallID:        "t375-call",
		Name:          "fs.write",
		Args: json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"375-r1 byte"}`,
			filepath.ToSlash(path))),
	}
}

// Test375UnattendedConsoleRefusesTheL1Write is this ticket's AC#2 in one case,
// under BOTH modes the probe measured: auto_approve (the mode used to silence
// the question, so 0 cards were ever raised) and ask_every_step (the window used
// to open and expire into execution). Neither may write, and neither may show a
// card nobody could answer.
func Test375UnattendedConsoleRefusesTheL1Write(t *testing.T) {
	for _, mode := range []string{risk.ModeAutoApproveName, risk.ModeAskEveryStepName} {
		t.Run(mode, func(t *testing.T) {
			dir := t.TempDir()
			t375Config(t, dir, mode)
			var (
				stdout, stderr bytes.Buffer
				cards          int
				executed       bool
				text           string
				elapsed        time.Duration
			)
			target := filepath.Join(dir, "new-plain.txt")
			code := runTextTask(runSpec{
				argv:    []string{"总结一下 这份笔记"},
				stdout:  &stdout,
				stderr:  &stderr,
				dataDir: dir,
				notify:  swallowPoster,
				onRuntime: func(rt *agentRuntime) {
					// Premise: this run really has no answer source. If a reader
					// ever hands one in, this case stops measuring the ticket.
					if rt.reply != nil {
						t.Fatal("premise broken: this run has a reply surface")
					}
					revoke := rt.gate.AdmitTextTask("t375-unattended")
					defer revoke()
					t0 := time.Now()
					res, err := rt.bridge.Execute(context.Background(),
						t375Call("t375-unattended", "t375-unattended-corr", target))
					elapsed = time.Since(t0)
					if err != nil {
						t.Fatalf("bridge.Execute returned a host fault: %v", err)
					}
					executed = !res.IsError
					text = res.Text
					cards = rt.windowCount()
				},
			})
			if code != 0 {
				t.Fatalf("run exit = %d, stderr:\n%s", code, stderr.String())
			}
			audit := stderr.String()
			if executed {
				t.Fatalf("the L1 write EXECUTED in a run with no answer source (mode %s, %s): a window's "+
					"expiry was read as consent - %q", mode, elapsed.Round(time.Millisecond), text)
			}
			if !strings.Contains(text, "无人可答") {
				t.Errorf("the refusal does not name the cause: %q", text)
			}
			if _, err := os.Stat(target); !os.IsNotExist(err) {
				t.Errorf("the refused write landed on disk anyway: %v", err)
			}
			if cards != 0 {
				t.Errorf("cards = %d, want 0: a card nobody can answer must not be raised", cards)
			}
			if !strings.Contains(audit, "L1-UNANSWERABLE") {
				t.Errorf("the ledger has no L1-UNANSWERABLE line; stderr:\n%s", audit)
			}
			if strings.Contains(stdout.String(), "[确认 L1 ") {
				t.Errorf("the console rendered an L1 card nobody could answer; stdout:\n%s", stdout.String())
			}
			if mode == risk.ModeAutoApproveName {
				if !strings.Contains(audit, "MODE-UNANSWERABLE") {
					t.Errorf("the mode was never reported as downgraded; stderr:\n%s", audit)
				}
				if strings.Contains(audit, "MODE-SILENCE mode=auto_approve") {
					t.Errorf("a call was silenced by auto_approve while nobody could answer it; stderr:\n%s", audit)
				}
			}
		})
	}
}

// Test375AttendedConsoleStillExecutesTheL1Write is the positive control: a run
// WITH a reply surface keeps the pre-375 route exactly - one card, the window's
// expiry means execute, and the file lands.
func Test375AttendedConsoleStillExecutesTheL1Write(t *testing.T) {
	dir := t.TempDir()
	t375Config(t, dir, risk.ModeAskEveryStepName)
	pr, pw := io.Pipe()
	defer pw.Close() // ends the listener's read only after the run has finished

	var (
		stdout, stderr bytes.Buffer
		cards          int
		executed       bool
		text           string
		elapsed        time.Duration
	)
	target := filepath.Join(dir, "attended.txt")
	code := runTextTask(runSpec{
		argv:    []string{"总结一下 这份笔记"},
		stdout:  &stdout,
		stderr:  &stderr,
		dataDir: dir,
		notify:  swallowPoster,
		reply:   pr,
		onRuntime: func(rt *agentRuntime) {
			if rt.reply == nil {
				t.Fatal("premise broken: the reply surface did not attach")
			}
			revoke := rt.gate.AdmitTextTask("t375-attended")
			defer revoke()
			t0 := time.Now()
			res, err := rt.bridge.Execute(context.Background(),
				t375Call("t375-attended", "t375-attended-corr", target))
			elapsed = time.Since(t0)
			if err != nil {
				t.Fatalf("bridge.Execute returned a host fault: %v", err)
			}
			executed = !res.IsError
			text = res.Text
			cards = rt.windowCount()
		},
	})
	if code != 0 {
		t.Fatalf("run exit = %d, stderr:\n%s", code, stderr.String())
	}
	if !executed {
		t.Fatalf("an attended run refused its L1 write: %q", text)
	}
	if _, err := os.Stat(target); err != nil {
		t.Errorf("the allowed write did not land: %v", err)
	}
	if cards != 1 {
		t.Errorf("cards = %d, want 1 (the window is still shown when somebody can answer)", cards)
	}
	if !strings.Contains(stderr.String(), "decision=timeout->execute") {
		t.Errorf("the L1 polarity moved; stderr:\n%s", stderr.String())
	}
	if strings.Contains(stderr.String(), "L1-UNANSWERABLE") {
		t.Errorf("the guard fired in an attended run; stderr:\n%s", stderr.String())
	}
	if elapsed < 2*time.Second {
		t.Errorf("the L1 window did not block for its 2s (%s): the pre-execution window is the point", elapsed)
	}
}

// Test375ResidentGateRefusesTheWindowWithNoBall is assembly D/B's reading at its
// source: the resident leg's own gate, built at boot (85 lines before the ball
// can be bound), asked to run an L1 window in a process whose orb never came up.
//
// Before this ticket the refusal came from the surface (ballCardUI.Prompt's
// errResidentNoBall) - the right outcome by luck of that one surface being
// honest. Now the gate itself refuses, which is what makes the same seat safe
// for a surface that cannot tell (the console prints a card either way).
func Test375ResidentGateRefusesTheWindowWithNoBall(t *testing.T) {
	ra := newResidentApproval()
	if ra.ui.hasBall() {
		t.Fatal("premise broken: an orb is attached")
	}
	revoke := ra.gate.AdmitTextTask("t375-resident")
	defer revoke()

	a, why := ra.gate.PendingWindow(context.Background(), tools.Decision{
		TaskID: "t375-resident", Tool: "fs.write", Level: risk.L1,
		Params: map[string]any{"path": "/work/t375.txt"},
		Args:   json.RawMessage(`{"path":"/work/t375.txt"}`),
		Reason: "R1: 工具声明为下界（L1）",
		Paths:  []string{"/work/t375.txt"},
	})
	if a != tools.AnswerReject {
		t.Fatalf("answer = %q, want reject", a)
	}
	if !strings.Contains(why, "无人可答") {
		t.Errorf("the gate did not name the cause: %q", why)
	}
	if strings.Contains(why, errResidentNoBall.Error()) {
		t.Errorf("the refusal still blames the surface (%q): the guard must not depend on a UI that "+
			"happens to be honest", why)
	}
	if n := ra.ui.displayedCards(); n != 0 {
		t.Errorf("the resident surface counted %d displayed cards with no window attached", n)
	}
}

// Test375InjectedGateMustDeclareItsAnswerSource pins the structural half: the
// mode read and the gate's L1 guard are two declarations of ONE fact, and a host
// that hands the assembly a predicate while its gate has none would get the
// downgrade without the guard - the hole, now with an audit line claiming
// otherwise. (This case names API the change introduces, so on the old tree its
// red is a compile error, not a run.)
func Test375InjectedGateMustDeclareItsAnswerSource(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	ra := newResidentApproval()
	// A gate built WITHOUT Options.Answer: the shape every host had before this
	// ticket, and the shape a new host could still ship by accident.
	bare := approval.New(approval.Options{UI: ra.ui, Logf: func(string, ...any) {}})
	ar, code := assembleRuntime(runSpec{
		stdout:  f.out,
		stderr:  f.err,
		dataDir: f.dir,
		gate:    bare,
		ui:      ra.ui,
		cards:   ra.cards,
		taskCtx: ra.root,
		answer:  func() bool { return true },
		notify:  swallowPoster,
	})
	if ar != nil {
		defer ar.close()
	}
	if code == 0 {
		t.Fatalf("assembleRuntime accepted a gate with no answer source while the host declared one: " +
			"the mode read would drop to ask_every_step and the L1 window would still execute, which is " +
			"the defect ticket 375 is about")
	}
	if !strings.Contains(f.err.String(), "AnswerSource=nil") {
		t.Errorf("the refusal never named the missing half; stderr:\n%s", f.err.String())
	}
}
