package agent

import (
	"context"
	"strings"
	"testing"
	"time"
)

// Acceptance criterion 6: the D11(1) control layer short-circuits with ZERO
// provider round-trips. The assertion is on the mockllm/replayer request count,
// so a regression that quietly routed "停" through the model would fail loudly.

func TestControlWordsNeverReachTheProvider(t *testing.T) {
	cases := []struct {
		utterance string
		verb      ControlVerb
	}{
		{"停", ControlStop},
		{"取消", ControlCancel},
		{"重说", ControlRepeat},
		{"大声点", ControlLouder},
		{"确认", ControlConfirm},
		{"  停 \n", ControlStop}, // transcript punctuation/whitespace tolerance
		{"。确认。", ControlConfirm},
	}
	for _, tc := range cases {
		var seen []ControlVerb
		h := newHarness(t, "text-reply", withControl(func(v ControlVerb, _ string) ControlOutcome {
			seen = append(seen, v)
			return ControlOutcome{Handled: true, Text: "ok"}
		}))

		res := h.run(tc.utterance)

		if res.Status != StatusControl {
			t.Errorf("%q: status = %s, want control", tc.utterance, res.Status)
		}
		if got := h.requests(); got != 0 {
			t.Errorf("%q: provider requests = %d, want 0 (D11(1) is local regex)", tc.utterance, got)
		}
		if len(seen) != 1 || seen[0] != tc.verb {
			t.Errorf("%q: handler saw %v, want [%s]", tc.utterance, seen, tc.verb)
		}
		if e := h.sink.LastOf(EvControl); e == nil || e.Verb != tc.verb {
			t.Errorf("%q: control event = %+v", tc.utterance, e)
		}
		// No task_log row is opened for a control word: it is not a task.
		if len(h.loop.History()) != 0 {
			t.Errorf("%q: control word polluted the conversation: %s",
				tc.utterance, dumpHistory(h.loop.History()))
		}
	}
}

// Near-misses are business intent, not control words: the two-layer split of
// D11 must not swallow real requests.
func TestNearMissUtterancesGoToTheLoop(t *testing.T) {
	for _, u := range []string{"取消这个任务吧", "停止播放音乐", "重说一下今天的安排", "确认一下会议室", "大声点说话好吗"} {
		if _, ok := MatchControl(u); ok {
			t.Errorf("MatchControl(%q) = true, want false (bare-word rule)", u)
		}
	}
	h := newHarness(t, "text-reply")
	res := h.run("取消这个任务吧")
	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s), want completed", res.Status, res.Message)
	}
	if h.requests() != 1 {
		t.Errorf("requests = %d, want 1 (it is business intent)", h.requests())
	}
}

// The loop's own default handler can stop the task it is running (D11(1) is
// wired to cancellation without any host code in between).
func TestDefaultControlCancelsRunningTask(t *testing.T) {
	h := newHarness(t, "slow-tool", withConfig(func(c *Config) {
		c.PerToolTimeout = 10 * time.Second
	}))
	task := h.loop.RunAsync(context.Background(), "慢一点也行")
	if !h.waitForFirstRequest(3 * time.Second) {
		t.Fatal("provider never received a request")
	}

	res := h.loop.Run(context.Background(), "停")
	if res.Status != StatusControl {
		t.Fatalf("control status = %s", res.Status)
	}
	if !strings.Contains(res.Message, "已停止") {
		t.Errorf("ack = %q", res.Message)
	}
	runner := task.Wait()
	if runner.Status != StatusCancelled {
		t.Errorf("running task status = %s (%s), want cancelled", runner.Status, runner.Message)
	}
	if runner.RootPending != 0 {
		t.Errorf("root pending = %d, want 0", runner.RootPending)
	}
	// The control word itself cost no round-trip; only the cancelled task did.
	if h.requests() != 1 {
		t.Errorf("requests = %d, want 1", h.requests())
	}
}

// An unhandled verb (no host wired) must be reported, not silently dropped.
func TestControlWithoutAConsumerIsVisible(t *testing.T) {
	h := newHarness(t, "text-reply")
	res := h.run("大声点")
	if res.Status != StatusControl {
		t.Fatalf("status = %s", res.Status)
	}
	if res.Message == "" {
		t.Error("an unhandled control word must still say something")
	}
	if h.requests() != 0 {
		t.Errorf("requests = %d, want 0", h.requests())
	}
}

// Steering: a runtime insert joins the conversation before the next model call
// and costs no extra round-trip of its own (Pi inner queue, D21#4).
//
// The slow-tool fixture is deliberate: round 1's 400ms sleep keeps the loop
// inside its first round while the test goroutine calls Steer, so the insert
// deterministically lands in round 2's drainSteering() rather than racing an
// instant echo provider. The assertion (round 2's wire body carries the steered
// text, no extra round-trip) is unchanged.
func TestSteeringInsertsIntoNextRound(t *testing.T) {
	h := newHarness(t, "slow-tool")
	task := h.loop.RunAsync(context.Background(), "现在天气怎么样")
	if !h.waitForFirstRequest(3 * time.Second) {
		t.Fatal("no request")
	}
	if !h.loop.Steer("顺便说下湿度") {
		t.Fatal("steering refused (agent.steering_enabled?)")
	}
	res := task.Wait()
	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s)", res.Status, res.Message)
	}
	if !strings.Contains(dumpHistory(h.loop.History()), "顺便说下湿度") {
		t.Errorf("steering text never joined the conversation:\n%s", dumpHistory(h.loop.History()))
	}
	// It reached round 2's wire body.
	bodies := h.requestBodies()
	if len(bodies) != 2 || !strings.Contains(string(bodies[1]), "顺便说下湿度") {
		t.Errorf("round 2 body lacks the steered instruction (bodies=%d)", len(bodies))
	}
	// Steering cost no round-trip of its own: one task still spent exactly the
	// provider turns the golden sequence defines.
	if h.requests() != 2 {
		t.Errorf("requests = %d, want 2 (steering adds no LLM call)", h.requests())
	}
}

func TestSteeringDisabled(t *testing.T) {
	h := newHarness(t, "text-reply", withConfig(func(c *Config) { c.SteeringEnabled = false }))
	if h.loop.Steer("插一句") {
		t.Error("Steer must report false when agent.steering_enabled is off")
	}
}
