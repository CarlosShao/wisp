//go:build windows

package main

// Ticket 423's feed readings: the census caught five config keys whose slots
// existed and whose production feed did not - the ticket-10 [agent] tail
// (max_rounds / token_budget / loop_guard.repeat_thresholds) and the three
// voice keys ([session] conversation_idle_sec / settling_sec and [voice]
// punctuation) whose owning tickets (28 / 27) had closed without carrying the
// feed away. Each case below pins one feed and states its reverse: the leg is
// organized so deleting the exact feed line reddens exactly the case above.
//
// Staging: NO sherpa DLLs and no model install are needed. The [session]
// cases end at voicePlan/session.ScopeOptions (the plan type, not an engine);
// the punctuation case calls speech.OpenPunctuation, which refuses a bad or
// missing model with an error BEFORE any sherpa call (checkOnnxModelFile);
// the agent case runs the assembly root against the mock LLM. The same
// staging contract the roster and tier cases run under.

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/session"
	"github.com/CarlosShao/wisp/internal/speech"
)

// Test423AgentTailReachesAgentConfig is the ticket-10 tail: execute() hands
// max_rounds / token_budget / loop_guard.repeat_thresholds to agent.Config at
// the loop's construction. The read is execute()'s own published wiring
// (rt.loopOpt - the same production read Test233HotTierAgentReachesTheNextTask
// takes), so the case asserts the loop was configured from a config.toml that
// names non-default values, not from the loop's own fallbacks.
func Test423AgentTailReachesAgentConfig(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	// The runtime capture, the 233 shape: the hook runs inside the assembly
	// root with the runtime it assembled, so the loopOpt this case reads is
	// the one the production root published.
	var live *agentRuntime
	h.rtHook = func(rt *agentRuntime) { live = rt }
	cfgPath := filepath.Join(h.dir, configFileName)
	body, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read config.toml: %v", err)
	}
	// Non-default, chosen against internal/agent's own fallbacks (50 / scaled
	// 200k / 3,5,8) so every slot can only carry the file's number.
	hand := "\n[agent]\nmax_rounds = 9\ntoken_budget = 12345\n\n[agent.loop_guard]\nrepeat_thresholds = [2, 4, 6]\n"
	if err := os.WriteFile(cfgPath, []byte(string(body)+hand), 0o600); err != nil {
		t.Fatalf("plant the [agent] tail: %v", err)
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("wisp run exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	if live == nil {
		t.Fatal("the runtime hook never ran, so nothing here read the assembled host")
	}
	if !live.loopOptSet {
		t.Fatal("execute() never published its options, so the feed cannot be read")
	}
	cfg := live.loopOpt.Config
	if cfg.MaxRounds != 9 {
		t.Errorf("agent.Config.MaxRounds = %d, want the file's 9 (the feed rides execute()'s own wiring; the loop's own fallback is 50)", cfg.MaxRounds)
	}
	if cfg.TokenBudget != 12345 {
		t.Errorf("agent.Config.TokenBudget = %d, want the file's 12345 (the loop's own fallback is the scaled 200k)", cfg.TokenBudget)
	}
	if len(cfg.RepeatThresholds) != 3 || cfg.RepeatThresholds[0] != 2 || cfg.RepeatThresholds[1] != 4 || cfg.RepeatThresholds[2] != 6 {
		t.Errorf("agent.Config.RepeatThresholds = %v, want [2 4 6] from the file (an empty ladder keeps the loop's 3,5,8)", cfg.RepeatThresholds)
	}
	// The premise half: a default file carries no such override, so the values
	// above are provably the feed's and not the fixture's defaults.
	def := config.NewDefaults()
	if def.Agent.MaxRounds == 9 || def.Agent.TokenBudget == 12345 {
		t.Fatalf("this case's premise moved: the shipped defaults already carry the planted values (%d / %d)",
			def.Agent.MaxRounds, def.Agent.TokenBudget)
	}
}

// voice423InstantClock is session.Clock without the waiting (the wire-r1
// shape), with the durations it was asked for kept as the reading: a scope
// built from the feed's options must arm its timers with exactly those
// windows, and the clock is where that number is written down.
type voice423InstantClock struct {
	asked []time.Duration
}

func (c *voice423InstantClock) AfterFunc(d time.Duration, fn func()) session.Timer {
	c.asked = append(c.asked, d)
	return voice423InstantTimer{}
}

type voice423InstantTimer struct{}

func (voice423InstantTimer) Stop() bool { return true }

// errVoice423NoEngines is the loader's answer: this case arms no engine, the
// scope only needs a non-nil Loader to exist.
var errVoice423NoEngines = errors.New("423 case arms no engines")

// Test423SessionKeysReachTheScopeOptions is the boot LoadFile feed: the same
// single pass that has carried [session].warm_timeout_sec since ticket 28 now
// also carries conversation_idle_sec and settling_sec into the plan's duration
// fields, and through the residency's ScopeOptions into the scope at its
// construction. Asserted in two hops - the plan (the boot feed's own receipt)
// and the exact ScopeOptions table the scope opens with - and the scope side
// is closed with a real session.NewScope over a recording clock: the durations
// it arms its three timers with ARE the fed values, which is the "the scope
// really opens on them" half no struct comparison can fake.
func Test423SessionKeysReachTheScopeOptions(t *testing.T) {
	dataDir := voice15eDataDir(t)
	// resolveVoicePlan's one hard requirement is the VAD file; the hatch is the
	// same env var production and the other cases use - no model install.
	t.Setenv("WISP_SILERO_MODEL", voice15eVadModel(t))
	voice15eWriteConfig(t, dataDir, func(c *config.Config) {
		c.Session.WarmTimeoutSec = 41
		c.Session.ConversationIdleSec = 44
		c.Session.SettlingSec = 7
	})
	ra := &residentAudio{frames: make(chan []byte)}
	rv := buildResidentVoice(dataDir, ra, nil, nil, nil)
	if rv.verdict != "语音链已装配" {
		t.Fatalf("the voice leg did not assemble from the bare dataDir: verdict=%q (problems=%v)",
			rv.verdict, rv.problemsSnapshot())
	}
	// Hop 1: the boot feed's own receipt - the plan carries the three windows.
	if got := rv.plan.WarmTimeout; got != 41*time.Second {
		t.Errorf("plan.WarmTimeout = %v, want 41s (the ticket-28 feed this wiring rides)", got)
	}
	if got := rv.plan.ConversationIdle; got != 44*time.Second {
		t.Errorf("plan.ConversationIdle = %v, want the file's 44s (ticket 423's feed; 0 keeps the session package's 30s)", got)
	}
	if got := rv.plan.SettleStatic; got != 7*time.Second {
		t.Errorf("plan.SettleStatic = %v, want the file's 7s (ticket 423's feed; 0 keeps 3s)", got)
	}
	// Hop 2: the ScopeOptions the residency opens one epoch's scope with - the
	// feed survives the plan -> residency hop in the production shape.
	opts := rv.residency.scopeOptions()
	if got := opts.WarmTimeout; got != 41*time.Second {
		t.Errorf("ScopeOptions.WarmTimeout = %v, want 41s", got)
	}
	if got := opts.ConversationIdle; got != 44*time.Second {
		t.Errorf("ScopeOptions.ConversationIdle = %v, want 44s", got)
	}
	if got := opts.SettleStatic; got != 7*time.Second {
		t.Errorf("ScopeOptions.SettleStatic = %v, want 7s", got)
	}
	// Hop 3, the scope side: NewScope over these options arms a timer per
	// window it accepted, with the fed durations verbatim (its 0 = default
	// contracts mean a re-defaulted value would betray the feed here).
	// PrivacyConfirm is granted in-process: the first enable needs an L2 card
	// surface, and this case's subject is the window the card unlocks, not the
	// card (28's own cases did the same).
	clock := &voice423InstantClock{}
	opts.Clock = clock
	opts.Loader = func(kind session.EngineKind) (session.Engine, error) {
		return nil, errVoice423NoEngines
	}
	opts.Hooks.PrivacyConfirm = func() bool { return true }
	sc, err := session.NewScope(opts)
	if err != nil {
		t.Fatalf("NewScope refused the feed's own options: %v", err)
	}
	if ok, err := sc.EnableConversation(); err != nil || !ok {
		t.Fatalf("EnableConversation refused on a scope with a recording clock (ok=%v): %v", ok, err)
	}
	// The conversation timer must carry the file's 44s (NewScope's own default
	// would be 30s). Fire it to prove the scope accepted the value as armed:
	// the Conversation row lands the scope in Warm.
	found := false
	for _, d := range clock.asked {
		if d == 44*time.Second {
			found = true
		}
	}
	if !found {
		t.Fatalf("the scope armed no 44s conversation timer; the durations it asked for: %v", clock.asked)
	}
	if err := sc.Dispose("423 case done"); err != nil {
		t.Fatalf("Dispose: %v", err)
	}
}

// Test423PunctuationSwitchGatesTheStage is the [voice] punctuation switch,
// in the four shapes the ticket's wiring gives it:
//
//	1  the OFF contract      OpenPunctuation(enabled=false) is nil,nil
//	2  the boot feed         the shipped default (true) reaches the plan
//	3  the loader pass       the residency's loader hands the plan field to the
//	                         stage's enabled argument (engines_windows.go:445)
//	4  the acquire gate      a switched-off epoch acquires no punct engine at
//	                         all, in the real "no checkpoint on this machine"
//	                         shape - no spurious problem note, clean raw finals
//
// Halves 3 and 4 are what make the reverse legs real: flipping the loader's
// enabled argument back to the hardcoded `true` it used to be reddens half 3
// (the fake model gets opened and errors), and deleting the acquire gate's
// `r.plan.PunctEnabled &&` reddens half 4 (the empty-path loader error becomes
// a per-turn 标点级没起来 note that the OFF shape must not produce).
func Test423PunctuationSwitchGatesTheStage(t *testing.T) {
	// A file that is NOT an ONNX model (a Go source file passes the path
	// checks and fails the header check): with the feed wired, enabled=false
	// answers nil,nil before any model is touched; with the feed gutted, the
	// same call tries to open it and ERRORS. This works on every machine,
	// including one whose System32 onnxruntime.dll is older than the C API
	// requires (the wire cases' known shadowing posture).
	notAModel := voice15eRepoFile(t, filepath.Join("cmd", "wisp", "run.go"))

	// Half 1, the OFF contract, with the exact arguments the production
	// loader passes.
	p, err := speech.OpenPunctuation(speech.PunctuationConfig{ModelPath: notAModel}, false)
	if err != nil {
		t.Fatalf("OpenPunctuation(enabled=false) errored: %v", err)
	}
	if p != nil {
		t.Fatalf("OpenPunctuation(enabled=false) returned a stage: the switch does not reach the stage's own off contract")
	}

	// Half 2, the boot feed: the schema default of the switch is true and the
	// boot LoadFile must carry it into the plan - a plan built by hand says
	// "on" the way the config default does (the two existing wire cases set
	// it explicitly; this pins the boot side).
	dataDir := voice15eDataDir(t)
	t.Setenv("WISP_SILERO_MODEL", voice15eVadModel(t))
	voice15eWriteConfig(t, dataDir, nil)
	ra := &residentAudio{frames: make(chan []byte)}
	rv := buildResidentVoice(dataDir, ra, nil, nil, nil)
	if rv.verdict != "语音链已装配" {
		t.Fatalf("the voice leg did not assemble from the bare dataDir: verdict=%q (problems=%v)",
			rv.verdict, rv.problemsSnapshot())
	}
	if !rv.plan.PunctEnabled {
		t.Fatal("the boot plan carries PunctEnabled=false off the shipped defaults, want true (the schema default is true and the feed must carry it)")
	}

	// Half 3, the loader pass: the residency hands the plan field to the
	// stage's enabled argument. With the plan switched off, the loader must
	// come back clean with a nil punctuator (the off contract), never an
	// error.
	off := voicePlan{PunctEnabled: false, PunctModelPath: notAModel}
	roff := newVoiceResidency(off, newVoiceOwnerQueue(), nil)
	eng, lerr := roff.loader(session.EnginePunct)
	if lerr != nil {
		t.Fatalf("the residency loader errored for a switched-off plan: %v (the plan field is not riding the enabled argument)", lerr)
	}
	h, ok := eng.(*voiceEngineHandle)
	if !ok || h == nil || h.punct != nil {
		t.Fatalf("the loader answered %T for a switched-off plan, want a handle with a nil punctuator", eng)
	}

	// Half 4, the acquire gate, in the real no-checkpoint shape (switch off +
	// a non-empty but absent model path - exactly what a machine that set
	// punctuation=false without the CT-Transformer file boots with): an ASR
	// turn must acquire no punct engine - the slot sees exactly one load (the
	// ASR itself) - and record no problem. Deleting the gate's PunctEnabled
	// term makes this shape walk into the loader with the missing path
	// (OpenPunctuation(enabled=false) still mints a handle, so the slot counts
	// a second load) - the teeth of this half.
	gated := voicePlan{PunctEnabled: false, AsrDir: "stub", PunctModelPath: "stub/absent-punct-model.onnx"}
	rg := newVoiceResidency(gated, newVoiceOwnerQueue(), nil)
	rg.asrLoader = func() (speech.AsrEngine, error) { return &voiceWireAsr{}, nil }
	if _, _, aerr := rg.acquireTurn(session.EngineASR, "423 off-gate"); aerr != nil {
		t.Fatalf("acquireTurn(ASR) errored with the switch off: %v", aerr)
	}
	if st := rg.stats(); st.ColdLoads != 1 {
		t.Fatalf("a switched-off ASR epoch loaded %d engine(s), want exactly 1 (the ASR): the acquire gate is not reading the switch", st.ColdLoads)
	}
	if probs := rg.problemsSnapshot(); len(probs) != 0 {
		t.Fatalf("a switched-off epoch booked problem notes %v: the acquire gate does not read the switch", probs)
	}

	// The ON control (a real stage opens under enabled=true with the real
	// checkpoint) needs an onnxruntime.dll the C API accepts; on a machine
	// where the OS copy shadows the build's, the stage's constructor refuses
	// with exactly that named error and the case skips honestly rather than
	// pretending the on-state was measured. Staged runs (the wire cases'
	// contract) measure it.
	on, err := speech.OpenPunctuation(speech.PunctuationConfig{ModelPath: voiceWirePunctModel(t)}, true)
	if err != nil {
		if strings.Contains(err.Error(), "onnxruntime") {
			t.Skipf("the on-state control needs a loadable onnxruntime.dll: %v", err)
		}
		t.Fatalf("OpenPunctuation(enabled=true) with the real model errored: %v", err)
	}
	if on == nil {
		t.Fatal("OpenPunctuation(enabled=true) returned nil with a valid model: the control failed")
	}
	if err := on.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}
