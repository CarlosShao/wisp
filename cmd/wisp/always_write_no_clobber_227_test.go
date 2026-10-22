package main

// Ticket 227 AC#6 - the one hop ticket 226's only host-driven case left
// unpinned: under the REAL production tick, a hand edit is still seen.
//
// WHY THIS CASE EXISTS. TestAC1AlwaysBranchDoesNotRevertAHandEditedKey drives
// the whole assembled run and ends by asserting the OPPOSITE direction of D36's
// hot tier: it requires the running process to still hold the theme this boot
// loaded ("the guarded write applied the hand-edited theme to the running
// process ... reading a hand edit back is ticket 223"). That assertion was
// written while the comment above it was true - "while the process is running and
// nothing polls" - and it stopped being true at 248095d1, when ticket 223 wired
// Manager.CheckAndReload to a 1s tick that `wisp run` arms on every start
// (cmd/wisp/config_reload.go's Spawn -> configReloadTick -> reloadOnce).
//
// NOTE ON LINE NUMBERS: the two quotes above are identified by their words, not
// by a `file:LINE`. This leg appended a correction to that comment, which moved
// every line below it, so a number here would have been stale the moment it was
// written - the rot this repository has paid for before.
//
// So the 226 case is green only while its own hand-edit-to-read window stays
// under one tick; it is a witness with a stopwatch in its mouth. Measured on this
// box (docs/evidence/s1/227-r1-impl-part1.md sec 1): the same scenario held for 2.5s
// after the hand edit reddens that assertion with
// `got "light", want the dark this run loaded`, while the unmodified case reads
// it 58ms after the edit and never reaches a tick.
//
// This case pins the direction that is actually true, on the same assembled host
// and with the same hand-edit shape, and it does it by WAITING for the tick's own
// statement instead of by outrunning it: no sleep is shortened, no interval is
// widened, and nothing here measures how long anything took.
//
// WHAT MAKES THIS A RESIDENT CASE AND NOT A Manager UNIT CASE. Nothing below
// constructs a config.Manager, assigns a hook, or calls CheckAndReload. The poll
// is the one cmd/wisp installs: assembleRuntime -> startConfigReload ->
// observe.Default.Spawn("watchdog", "config", ...) -> configReloadTick, and the
// only thing this case does is put a file on disk and wait for that goroutine's
// own audit line. Deleting the Spawn reddens this case; that reading is in
// docs/evidence/s1/227-r1-impl-part1.md sec 4.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and dies
// at load (0xc0000135) unless third_party/sherpa-onnx is on PATH.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// ac6Applied227 is the tick's own statement that it read a change and applied it.
// `hot=[app]` is the tier verdict, so a reload that booked any other section
// cannot satisfy this wait.
const ac6Applied227 = "config: HOT-RELOAD state=applied hot=[app]"

// waitForWindow227 waits until the operator's stream, read FROM mark, carries
// every needle, and returns that window. A copy that existed before the plant
// cannot satisfy it. The wait is bounded and monotonic - a time.Timer plus a
// 20ms time.Ticker, the shape cmd/wisp/config_reload_223_test.go's awaits use -
// never a wall-clock subtraction (AGENTS.md 1.2's banned shape).
//
// It exists because the reload path writes its audit lines BEFORE the operator's
// own sentence (config_reload.go:183-198), so a single read at the moment the
// audit line arrives has a real window; that is the shape ticket 351 AC#1
// classified as instrumentation rather than a product defect. A receipt that
// never arrives still fails, at this helper's own deadline, with the window
// dumped.
func waitForWindow227(h *replyHost, mark string, needles []string, limit time.Duration) (string, bool) {
	deadline := time.NewTimer(limit)
	defer deadline.Stop()
	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()
	var win string
	for {
		now := h.out.String()
		if !strings.HasPrefix(now, mark) {
			// The stream only ever grows (syncWriter appends), so this can only
			// mean the mark is not this host's. Say so instead of searching the
			// wrong window.
			return "227 window: the operator stream no longer starts with the mark taken before " +
				"the plant\nmark:\n" + mark + "\ncurrent:\n" + now, false
		}
		win = strings.TrimPrefix(now, mark)
		all := true
		for _, n := range needles {
			if !strings.Contains(win, n) {
				all = false
				break
			}
		}
		if all {
			return win, true
		}
		select {
		case <-deadline.C:
			return win, false
		case <-tick.C:
		}
	}
}

// TestAC6AHandEditIsStillSeenUnderTheRealPoll is ticket 227 AC#6: the hop from
// "the operator edited config.toml" to "the running process holds that value"
// under the production tick.
//
// The scenario is ticket 226's, unchanged - hand edit [app] theme while the run
// is alive, then answer one 「以后这类都别问」 so the guarded write runs in the
// same process - and the only addition is that this case then waits for the tick
// instead of ending before it. What it pins:
//
//	(a) the tick's own audit line says it applied a hot [app] change,
//	(b) the live Manager now holds the hand-edited theme,
//	(c) the boot snapshot does NOT (so this is a mid-run change, not something
//	    the run started with),
//	(d) the operator got a user-visible receipt for it,
//	(e) the file still carries the hand edit, and ticket 226's guarded write
//	    still ran and still stored the rule it was answering for.
//
// (b) is the direction that 226 assertion forbids and D36 requires. (c) is what
// keeps (b) from being satisfied by the boot value. (e) is what keeps this case
// from being a quieter neighbour of the 226 scenario: the same widening answer is
// answered here, through the same gate, and the hand edit survives it.
func TestAC6AHandEditIsStillSeenUnderTheRealPoll(t *testing.T) {
	w, pw := newWidenRun(t, 40*time.Second, "t227-ac6", "t227-ac6-corr")
	cfgPath := filepath.Join(w.h.dir, configFileName)

	inner := w.h.rtHook
	w.h.rtHook = func(rt *agentRuntime) {
		// The hand edit, in ticket 226's own shape: one [app] key appended to the
		// file this run is watching, with nothing whatever to do with [fs].
		data, err := os.ReadFile(cfgPath)
		if err != nil {
			w.setFail(fmt.Errorf("read config.toml: %w", err))
			return
		}
		hand := string(data) + "\n[app]\ntheme = \"light\"\n"
		if err := os.WriteFile(cfgPath, []byte(hand), 0o600); err != nil {
			w.setFail(fmt.Errorf("hand edit config.toml: %w", err))
			return
		}
		// The 232 shape: everything asserted from `mark` on is bytes this run
		// printed AFTER the plant, so the start-up banner cannot stand in for the
		// reload receipt.
		mark := w.h.out.String()

		// Ticket 226's widening, so the guarded write runs in this same process.
		// inner also parks until the fs.write card reaches the native ledger.
		inner(rt)

		// The one hop this ticket is about: wait for the PRODUCTION tick to say
		// it applied the edit. This is a wait on the product's own statement - it
		// either arrives or this case fails with the whole trail dumped - not a
		// sleep, and not a shortened interval.
		if !waitForAudit187(w.h, ac6Applied227, 40*time.Second) {
			w.setFail(fmt.Errorf("the production reload tick never booked %q within 40s: the "+
				"hand edit was written while the run was alive, so either the tick is not "+
				"running or it did not see the change\nstderr:\n%s", ac6Applied227, w.h.err.String()))
			return
		}
		// (b) the running process holds the hand edit.
		if got := rt.mgr.Config().App.Theme; got != "light" {
			w.setFail(fmt.Errorf("the running process did not pick up the hand-edited theme "+
				"(got %q, want \"light\"): D36 makes app.theme a hot-tier key, so the tick that "+
				"just booked hot=[app] must have applied it", got))
			return
		}
		// (c) the boot snapshot did not move: this was a mid-run change.
		if got := rt.cfg.App.Theme; got == "light" {
			w.setFail(fmt.Errorf("the boot snapshot already held the hand-edited theme (%q), "+
				"so nothing here proves the tick did anything", got))
			return
		}
		// (d) the operator's own stream got a receipt naming this section. It is
		// WAITED for, not read once: reportReload writes its audit lines before
		// the operator's Fprintf (config_reload.go:183-198), so a single read here
		// is the exact window ticket 351 classified as instrumentation shape. The
		// wait is the same monotonic timer+ticker shape the rest of this package
		// uses - no wall-clock subtraction, no sleep.
		win, ok := waitForWindow227(w.h, mark,
			[]string{"配置热加载：", "[app]"}, 40*time.Second)
		if !ok {
			w.setFail(fmt.Errorf("no user-visible reload receipt naming [app] reached the "+
				"operator after the plant within 40s; stdout since the plant:\n%s", win))
			return
		}
	}

	go driveAlways(t, w, pw, true)

	if code := w.h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, w.h.out.String(), w.h.err.String())
	}
	<-w.done
	<-w.reply
	if err := w.failure(); err != nil {
		t.Fatalf("%v\nstdout:\n%s\nstderr:\n%s", err, w.h.out.String(), w.h.err.String())
	}

	// (e) the file still carries the hand edit, and the guarded write ran.
	text, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(text), "theme = 'light'") {
		t.Errorf("the hand-edited theme is not in config.toml any more, so the tick that "+
			"applied it was reading a file this run had already reverted:\n%s", text)
	}
	if !strings.Contains(string(text), filepath.ToSlash(w.h.outside)) {
		t.Errorf("the run did not store the rule it was answering for, so this is not the "+
			"ticket-226 scenario:\n%s", text)
	}
	audit := w.h.err.String()
	if !strings.Contains(audit, "approval: REPLY WIDEN-APPLIED") {
		t.Errorf("the widening was not booked as applied, so the guarded write never ran in "+
			"this process and (e) proves nothing:\n%s", audit)
	}
	// The two hops stay distinguishable: this line is the TICK's, and only the
	// tick writes it (config_reload.go's reportReload).
	if !strings.Contains(audit, "config: HOT-RELOAD-READER section=app") {
		t.Errorf("the tick did not book a reader verdict for [app], so the reload receipt this "+
			"case rests on came from somewhere else:\n%s", audit)
	}
}
