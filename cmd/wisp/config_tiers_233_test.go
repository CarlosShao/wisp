package main

// Ticket 233 - the readings the tier table had and no resident case carried:
// the hot tier reaching execute()'s own inputs (改了 → 不重启 → 行为真变了), the
// reload tier's event received by a hook cmd/wisp assigned (the tier had an
// emitter and no consumer in `wisp run`), and the tier's sentence on the
// operator's own stdout inside the post-plant window.
//
// WHAT MAKES THESE READINGS PRODUCTION AND NOT A MOCK. No case below assigns
// Manager.OnReload. Each one drives the assembled host (newReloadRun223 or
// assembleRuntime), so the hook assignment is the one cmd/wisp's
// startConfigReload installs and the event travels the route a real hand edit
// takes: plant -> 1s tick -> CheckAndReload -> planVoice -> commit -> OnReload.
// Deleting the assignment line reddens the reload case (its sentence never
// arrives) - that is this ticket's mutation M1, measured with the line
// commented out and restored byte-identical afterwards.
//
// The hot case reads execute()'s OWN wiring: after the tick hot-applies an
// [agent] edit, the NEXT task the same process runs must be built from the new
// values. The read is rt.loopOpt - the agent.Options execute() itself published
// for the task it just ran (run.go's `rt.loopOpt, rt.loopOptSet = opts, true`)
// - so the case asserts the loop was configured from the live Manager and not
// from the boot snapshot, without re-implementing any of execute().
//
// The needle discipline is 票 232/272's: every stdout assertion is waited for
// (time.NewTimer + 20ms time.NewTicker - no sleep-shortening, no wall-clock
// subtraction, d22scan ban #4) inside the window that starts at a mark taken
// BEFORE the plant, so a pre-plant copy of any sentence cannot satisfy anything.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and
// dies at load (0xc0000135) unless third_party/sherpa-onnx is on PATH.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
)

// reloadTierSentence233 is the operator sentence reportReloadTier prints. It is
// the reload tier's whole receipt: which sections, that the event was received,
// and what THIS host does not reload.
const reloadTierSentence233 = "本宿主没有可重载的子系统"

// auditReloadTier233 is the audit line the same hook writes, so the ledger and
// the operator's stream can be asserted separately (票 223 AC#5's rule: the
// sentence must live where its reader is, and each stream carries its own half).
const auditReloadTier233 = "config: RELOAD-TIER state=received"

// Test233ReloadTierEventReachesAProductionConsumer drives a real [voice] reload
// through the assembled host's own tick and requires the hook cmd/wisp assigned
// to fire: the operator's stdout carries the reload-tier receipt inside the
// post-plant window, the audit trail carries the same event with its effect
// verdict, and the memory itself moved (the tier applies immediately - what the
// receipt says did NOT happen is the subsystem reload, not the memory write).
func Test233ReloadTierEventReachesAProductionConsumer(t *testing.T) {
	r := newReloadRun223(t, 40*time.Second)
	r.live(t, func() {
		r.awaitAudit(t, "config: HOT-RELOAD state=armed")
		// The boot snapshot must not already hold the value we are planting -
		// otherwise nothing here proves a mid-run change arrived. The default
		// file carries the sherpa defaults, so the model line is absent.
		if r.rt.cfg.Voice.ASR.Model == "paraformer-zh-233" {
			t.Fatalf("this case's premise moved: the boot snapshot already carries model %q",
				"paraformer-zh-233")
		}
		mark := r.h.out.String()
		errMark := r.h.err.String()
		// The reload-tier edit, in the 223 plant shape: the host's fixture file
		// carries no [voice.asr] section, so one is APPENDED ahead of [fs]
		// (the same anchor the 223/227/255 plants use). The model line is a
		// different-length addition, which keeps the mtime+size fingerprint
		// honest.
		r.plant(t, "[fs]", "[voice.asr]\nmodel = \"paraformer-zh-233\"\n\n[fs]")
		// The tick's own tier verdict first: a reload-tier booking, not hot.
		applied := r.awaitAuditSince232(t, errMark, "config: HOT-RELOAD state=applied")
		if !strings.Contains(applied, "reload=[voice]") {
			t.Fatalf("the [voice] edit was not booked reload-tier: %q", applied)
		}
		// The tier table says voice.asr.model is reload - the 为什么 side of
		// the receipt.
		if tier, ok := config.TierOf("voice.asr.model"); !ok || tier != "reload" {
			t.Fatalf("TierOf(voice.asr.model) = %q/%v, want reload (the tier table the receipt stands on)", tier, ok)
		}
		// The event reached the production consumer: the audit line only the
		// assigned hook writes, waited inside the post-plant window.
		auditWin := r.awaitAuditSince232(t, errMark, auditReloadTier233)
		if !strings.Contains(auditWin, "sections=[voice]") {
			t.Fatalf("the reload-tier audit line does not name [voice]: %q", auditWin)
		}
		if !strings.Contains(auditWin, "effect=this-host-reloads-nothing") {
			t.Fatalf("the reload-tier audit line does not say what this host does not do: %q", auditWin)
		}
		// The operator's own stream, same window discipline: the sentence must
		// arrive AFTER the plant (no pre-plant copy can satisfy this).
		out := r.awaitStdoutSince232(t, mark, reloadTierSentence233)
		for _, needle := range []string{"[voice]", "重载档", "已生效到本进程内存", "不会触发模型重载"} {
			if !strings.Contains(out, needle) {
				t.Errorf("the reload-tier sentence omits %q; stdout since the plant:\n%s", needle, out)
			}
		}
		// The value itself IS live in the Manager (the tier applies
		// immediately); the receipt names what did not happen - the model swap.
		if got := r.rt.mgr.Config().Voice.ASR.Model; got != "paraformer-zh-233" {
			t.Errorf("live config voice.asr.model = %q, want the planted value (reload tier applies immediately)", got)
		}
	})
}

// Test233HotTierAgentReachesTheNextTask pins the hot tier's real-behaviour
// half: execute() reads the LIVE Manager per task (ticket 233's re-pointing of
// run.go's cfg line), so a hot-applied [agent] edit changes what the NEXT task
// the same process runs is built with. 改了 → 不重启 → 行为真变了 - read through
// the production composition (the second execute's own agent.Options), not
// through a memory probe on Manager.
func Test233HotTierAgentReachesTheNextTask(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	cfgPath := filepath.Join(h.dir, configFileName)
	var live *agentRuntime
	h.rtHook = func(rt *agentRuntime) {
		live = rt
		// The tick must be armed before the plant, so the window is the
		// production one and not a race against assembly.
		deadline := time.NewTimer(reloadCaseBudget)
		defer deadline.Stop()
		tick := time.NewTicker(20 * time.Millisecond)
		defer tick.Stop()
		for {
			if strings.Contains(h.err.String(), "config: HOT-RELOAD state=armed") {
				break
			}
			select {
			case <-deadline.C:
				t.Fatalf("the reload tick never armed; stderr:\n%s", h.err.String())
			case <-tick.C:
			}
		}
		// Premise: the boot snapshot holds no per-tool override (the schema
		// default is 0 = each tool's own class default).
		if rt.cfg.Agent.PerToolTimeoutMS != 0 {
			t.Fatalf("this case's premise moved: the boot snapshot already carries per_tool_timeout_ms %d",
				rt.cfg.Agent.PerToolTimeoutMS)
		}
		// The hot-tier edit, appended to the file this run is watching.
		hand := "\n[agent]\nper_tool_timeout_ms = 2500\n"
		data, err := os.ReadFile(cfgPath)
		if err != nil {
			t.Fatalf("read config.toml: %v", err)
		}
		time.Sleep(20 * time.Millisecond)
		if err := os.WriteFile(cfgPath, []byte(string(data)+hand), 0o600); err != nil {
			t.Fatalf("hand edit config.toml: %v", err)
		}
		// Wait for the tick's own booking of the change - the same shape the
		// 223/227 cases use (wait on the product's statement, never outrun it).
		applied := false
		deadline2 := time.NewTimer(reloadCaseBudget)
		defer deadline2.Stop()
		tick2 := time.NewTicker(20 * time.Millisecond)
		defer tick2.Stop()
		for !applied {
			if strings.Contains(h.err.String(), "hot=[agent]") {
				applied = true
				break
			}
			select {
			case <-deadline2.C:
				t.Fatalf("the tick never booked the [agent] edit as hot within the budget; stderr:\n%s",
					h.err.String())
			case <-tick2.C:
			}
		}
	}
	code := h.run("总结一下 这份笔记")
	if code != 0 {
		t.Fatalf("wisp run exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	if live == nil {
		t.Fatal("the runtime hook never ran, so nothing here read the assembled host")
	}
	// THE BEHAVIOUR READING: the task this same process ran was configured from
	// the live Manager - its loop options carry the planted per-tool budget and
	// the boot snapshot does not. rt.loopOpt is execute()'s own published
	// wiring for that task; if execute() still read the boot snapshot, this
	// would answer 0 and the case would fail.
	if live.loopOptSet != true {
		t.Fatal("execute() never published its options, so the per-task wiring cannot be read")
	}
	if got := live.loopOpt.Config.PerToolTimeout; got != 2500*time.Millisecond {
		t.Fatalf("the task that ran in this process was built with PerToolTimeout=%v, want 2.5s "+
			"(the hot-applied [agent] edit did not reach execute()'s own inputs - 改了 → 不重启 → 行为没变)", got)
	}
	// The frozen half stays frozen: the boot snapshot still holds the old
	// value, which is what separates "the tick applied it" from "the boot
	// already had it".
	if got := live.cfg.Agent.PerToolTimeoutMS; got == 2500 {
		t.Fatalf("the boot snapshot moved too, so nothing here proves a mid-run change: %d", got)
	}
}

// Test233ReloadTierReceiptOnlyForReloadTier asserts the reload-tier sentence is
// not reachable by the other two tiers' traffic: a hot edit produces no
// reload-tier receipt (its own honest receipt still prints), and the sentence
// count stays at zero until a reload-tier edit actually lands.
func Test233ReloadTierReceiptOnlyForReloadTier(t *testing.T) {
	r := newReloadRun223(t, 40*time.Second)
	r.live(t, func() {
		r.awaitAudit(t, "config: HOT-RELOAD state=armed")
		mark := r.h.out.String()
		// A HOT edit: [ball] size. Its own honest receipt is expected; the
		// reload-tier sentence must not appear.
		r.plant(t, "[fs]", "[ball]\nsize = 67\n\n[fs]")
		_ = r.awaitAuditSince232(t, r.h.err.String(), "hot=[ball]")
		hotOut := r.awaitStdoutSince232(t, mark, "值已换进本进程内存")
		if !strings.Contains(hotOut, "[ball]") {
			t.Errorf("the hot edit's own honest receipt never named [ball]; stdout:\n%s", hotOut)
		}
		if n := strings.Count(r.h.out.String(), reloadTierSentence233); n != 0 {
			t.Errorf("a hot-tier edit produced the reload-tier receipt %d time(s): the tier sentences are not separable; stdout:\n%s",
				n, r.h.out.String())
		}
	})
}
