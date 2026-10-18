package main

// Ticket 223 AC#1/AC#2/AC#3/AC#4/AC#6/AC#7 - the production readings for the
// hot-reload wiring, taken through the injection seam AGENTS.md §1.3 names for
// this CLI (`wisp run` with a scripted answer stream).
//
// WHAT MAKES THESE READINGS PRODUCTION AND NOT A MOCK. No case below assigns
// Manager.ConfirmLocked, Manager.OnRestartPending or Manager.OnReload. Each one
// drives runTextTask, so the tick, the hook and the restart notice are the ones
// cmd/wisp installed (config_reload.go), and the L2 card is answered by the
// reply listener the ticket 201 host attaches - the same route a human typing at
// this terminal travels. Deleting rt.startConfigReload() reddens every case in
// this file; that is the mutation reading the ticket asks for, and it is why the
// assertions read the card and the live memory rather than a field somebody set.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and
// dies at load (0xc0000135) unless third_party/sherpa-onnx is on PATH.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// reloadCaseBudget bounds one wait on the 1s tick plus a card answer. It is a
// test-side ceiling, not a product deadline: nothing on the reload path measures
// elapsed time to decide anything (config_reload.go's header carries that
// argument, and 票 223 AC#1 asks for it by name).
const reloadCaseBudget = 40 * time.Second

// restartSentence272 is the causal half of the operator's restart notice: the
// words that only a planted edit can be evidence for. 票 272 AC#2's two arms
// (zero copies before the plant, exactly one copy inside the window) name it so
// they search the same string the case already searches for at :622 (the await)
// and :647 (the operator's own copy) - those two literals are left verbatim
// because 票 272 AC#3 forbids touching an existing assertion.
const restartSentence272 = "本次运行不会生效"

type reloadRun223 struct {
	h  *replyHost
	pw *io.PipeWriter
	rt *agentRuntime
}

func newReloadRun223(t *testing.T, l2Wait time.Duration) *reloadRun223 {
	t.Helper()
	r := &reloadRun223{h: newReplyHost(t, l2Wait)}
	pr, pw := io.Pipe()
	r.pw = pw
	r.h.reply = pr
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })
	return r
}

// live drives the composition root and parks inside it: body runs after
// assembleRuntime returned - so the reload tick, the D36 hook and the reply
// listener are all up - and before the agent task starts.
func (r *reloadRun223) live(t *testing.T, body func()) {
	t.Helper()
	r.h.rtHook = func(rt *agentRuntime) {
		r.rt = rt
		body()
	}
	code := r.h.run("总结一下 这份笔记")
	if code != 0 {
		t.Fatalf("wisp run exit %d\nstdout:\n%s\nstderr:\n%s",
			code, r.h.out.String(), r.h.err.String())
	}
}

func (r *reloadRun223) cfgPath() string { return filepath.Join(r.h.dir, configFileName) }

// plant rewrites one targeted piece of the live config.toml and proves the
// rewrite landed: a plant that silently misses is how a case would pass while
// proving nothing. The mtime+size fingerprint is what the poll reads, so the
// pause plus the length change is what makes exactly one tick see it.
func (r *reloadRun223) plant(t *testing.T, old, new string) {
	t.Helper()
	path := r.cfgPath()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config.toml: %v", err)
	}
	text := string(raw)
	if !strings.Contains(text, old) {
		t.Fatalf("plant target %q is not in config.toml:\n%s", old, text)
	}
	out := strings.Replace(text, old, new, 1)
	if out == text {
		t.Fatalf("plant of %q changed nothing", old)
	}
	time.Sleep(20 * time.Millisecond)
	if err := os.WriteFile(path, []byte(out), 0o600); err != nil {
		t.Fatalf("plant write: %v", err)
	}
}

// writeOver replaces the whole file - for the AC#4 cases, where the shape of the
// failure is the subject.
func (r *reloadRun223) writeOver(t *testing.T, body string) {
	t.Helper()
	time.Sleep(20 * time.Millisecond)
	if err := os.WriteFile(r.cfgPath(), []byte(body), 0o600); err != nil {
		t.Fatalf("write config.toml: %v", err)
	}
}

// awaitAudit waits for one sentence to reach this run's audit trail (stderr plus
// the persistent JSONL sink, both written by agentRuntime.auditf) and returns
// the whole trail, so a caller can also assert what is NOT in it.
func (r *reloadRun223) awaitAudit(t *testing.T, needle string) string {
	t.Helper()
	deadline := time.NewTimer(reloadCaseBudget)
	defer deadline.Stop()
	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()
	for {
		got := r.h.err.String()
		if strings.Contains(got, needle) {
			return got
		}
		select {
		case <-deadline.C:
			t.Fatalf("audit never carried %q within %v; full stderr:\n%s",
				needle, reloadCaseBudget, got)
		case <-tick.C:
		}
	}
}

// awaitStdout waits for one sentence to reach the operator's stdout and returns
// the whole stream. It exists because the reload path writes the audit trail
// BEFORE the operator line (reportRestartPending emits both auditf calls first,
// then the Fprintf - cmd/wisp/config_reload.go:282/:284 -> :288), so a case
// that awaited the audit and then read stdout ONCE has a real window and
// reddens intermittently (ticket 223 r2 AC#5: measured 1/25 here, 1/5 on the
// v1 leg, 1/15 by the orchestrator). The wait is bounded and monotonic - the
// same time.Timer/time.Ticker shape awaitAudit uses, not a wall-clock
// subtraction (AGENTS.md 1.2's banned shape). It is NOT an assertion and not a
// loosening: a sentence that never arrives still fails the case, at this
// helper's own deadline, with the full stdout and stderr dumped.
func (r *reloadRun223) awaitStdout(t *testing.T, needle string) string {
	t.Helper()
	deadline := time.NewTimer(reloadCaseBudget)
	defer deadline.Stop()
	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()
	for {
		got := r.h.out.String()
		if strings.Contains(got, needle) {
			return got
		}
		select {
		case <-deadline.C:
			t.Fatalf("stdout never carried %q within %v; full stdout:\n%s\nstderr:\n%s",
				needle, reloadCaseBudget, got, r.h.err.String())
		case <-tick.C:
		}
	}
}

// windowSince232 is the windowed read ticket 232 AC#2 asks for: one of this
// host's streams FROM a mark taken earlier, never the whole stream. It exists
// because a needle searched over the full stream is satisfied by bytes the run
// printed BEFORE the edit was planted - a start-up banner that happens to carry
// the same words would make the restart-tier assertions true without the
// restart notice ever saying anything (232's "前向风险" row, and AC#4's mutation
// H, which is measured in .scratch/wisp/probes/232/r2/logs/).
//
// The mark is a snapshot string, not an index: both of this host's streams are
// syncWriters whose bytes only ever get appended (reply listener, tick notices
// and the reload path all Write, nothing Rewinds), so the current stream must
// contain the snapshot. If it somehow does not, the case fails loudly instead of
// silently searching the wrong window. Test-only helper on the test's own host:
// it adds no exported name and touches no production file.
func windowSince232(t *testing.T, w *syncWriter, mark, stream string) string {
	t.Helper()
	now := w.String()
	if !strings.HasPrefix(now, mark) {
		t.Fatalf("232 window: the %s stream no longer starts with the mark taken before the plant;\nmark:\n%s\ncurrent:\n%s",
			stream, mark, now)
	}
	return strings.TrimPrefix(now, mark)
}

// stdoutSince232 is windowSince232 for the operator's stream.
func (r *reloadRun223) stdoutSince232(t *testing.T, mark string) string {
	t.Helper()
	return windowSince232(t, r.h.out, mark, "stdout")
}

// awaitStdoutSince232 waits for a sentence to arrive in the WINDOWED part of the
// operator's stdout (the bytes written after `mark`) and returns that window
// grown to include it. It is awaitStdout with a start marker instead of a
// whole-stream search: a copy of the sentence that existed BEFORE the edit was
// planted cannot satisfy it, which is exactly the shape 232 AC#4 pins.
//
// The await is the same monotonic shape the helpers above use - time.NewTimer
// plus a 20ms time.NewTicker, no wall-clock subtraction (d22scan ban #4), no
// bare goroutine (ban #1), no Skip: a sentence that never arrives reddens at
// reloadCaseBudget with both streams and the window dumped.
func (r *reloadRun223) awaitStdoutSince232(t *testing.T, mark, needle string) string {
	t.Helper()
	deadline := time.NewTimer(reloadCaseBudget)
	defer deadline.Stop()
	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()
	for {
		win := r.stdoutSince232(t, mark)
		if strings.Contains(win, needle) {
			return win
		}
		select {
		case <-deadline.C:
			t.Fatalf("stdout never carried %q AFTER the plant within %v; window since the mark:\n%s\nfull stdout:\n%s\nstderr:\n%s",
				needle, reloadCaseBudget, win, r.h.out.String(), r.h.err.String())
		case <-tick.C:
		}
	}
}

// awaitAuditSince232 is awaitStdoutSince232's twin for the audit trail: it waits
// for a line to reach stderr AFTER the mark was taken, so a pre-plant audit line
// that happens to carry the same words cannot open the window.
func (r *reloadRun223) awaitAuditSince232(t *testing.T, mark, needle string) string {
	t.Helper()
	deadline := time.NewTimer(reloadCaseBudget)
	defer deadline.Stop()
	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()
	for {
		win := windowSince232(t, r.h.err, mark, "audit")
		if strings.Contains(win, needle) {
			return win
		}
		select {
		case <-deadline.C:
			t.Fatalf("the audit trail never carried %q AFTER the plant within %v; window since the mark:\n%s\nfull stderr:\n%s",
				needle, reloadCaseBudget, win, r.h.err.String())
		case <-tick.C:
		}
	}
}

// awaitCard waits for the reload path to display a card and returns it.
func (r *reloadRun223) awaitCard(t *testing.T, tool string) liveCard {
	t.Helper()
	deadline := time.NewTimer(reloadCaseBudget)
	defer deadline.Stop()
	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()
	for {
		for _, c := range r.rt.liveCards.pending() {
			if c.Tool == tool {
				return c
			}
		}
		select {
		case <-deadline.C:
			t.Fatalf("no %q card was displayed within %v; stdout:\n%s\nstderr:\n%s",
				tool, reloadCaseBudget, r.h.out.String(), r.h.err.String())
		case <-tick.C:
		}
	}
}

// answer types one line on the operator's stream, exactly as the CLI's stdin
// would.
func (r *reloadRun223) answer(t *testing.T, line string) {
	t.Helper()
	if _, err := fmt.Fprintln(r.pw, line); err != nil {
		t.Fatalf("reply stream: %v", err)
	}
}

// awaitLive polls cond while the run is alive, so a change that lands on the
// tick's clock is read off the product instead of assumed.
func (r *reloadRun223) awaitLive(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.NewTimer(reloadCaseBudget)
	defer deadline.Stop()
	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()
	for {
		if cond() {
			return
		}
		select {
		case <-deadline.C:
			t.Fatalf("%s never happened within %v\nstdout:\n%s\nstderr:\n%s",
				what, reloadCaseBudget, r.h.out.String(), r.h.err.String())
		case <-tick.C:
		}
	}
}

// awaitAuditLine is awaitAudit narrowed to the one line that matched. Cases that
// assert what a sentence does NOT say must read one line: the whole trail
// mentions "L2" in the armed line's own detail, and a substring test over it
// would flag a tightening as if it had cost a card.
func (r *reloadRun223) awaitAuditLine(t *testing.T, needle string) string {
	t.Helper()
	trail := r.awaitAudit(t, needle)
	lines := strings.Split(trail, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], needle) {
			return strings.TrimSpace(lines[i])
		}
	}
	t.Fatalf("the trail returned no line for %q:\n%s", needle, trail)
	return ""
}

// allowedLine renders the [fs] allowed_dirs line the way newReplyHost seeds it.
func allowedLine(dirs ...string) string {
	quoted := make([]string, 0, len(dirs))
	for _, d := range dirs {
		quoted = append(quoted, fmt.Sprintf("%q", filepath.ToSlash(d)))
	}
	return "allowed_dirs = [" + strings.Join(quoted, ", ") + "]"
}

// TestTicket223RunArmsTheReloadTick is AC#1 and AC#2's hot tier: `wisp run` now
// has a production poller, and a hand edit lands in memory without a restart.
// The boot snapshot (rt.cfg) is the "before" and the live Manager is the
// "after", so this reads a change rather than a field somebody set.
func TestTicket223RunArmsTheReloadTick(t *testing.T) {
	// 64 is inside [ball]'s validated range (validate.go: 44-72) and different
	// from the schema default 56 (schema.go's `default:"56"` tag), so the plant
	// is a real change and a legal one.
	const plantedSize = 64
	r := newReloadRun223(t, 40*time.Second)
	r.live(t, func() {
		armed := r.awaitAudit(t, "config: HOT-RELOAD state=armed")
		if !strings.Contains(armed, "goroutine=watchdog") || !strings.Contains(armed, "owner=config") {
			t.Errorf("the armed line does not name its trigger shape: %q", armed)
		}
		if !strings.Contains(r.h.out.String(), "配置热加载已接管") {
			t.Errorf("the operator is never told the tick is running; stdout:\n%s", r.h.out.String())
		}
		if r.rt.cfg.Ball.Size == plantedSize {
			t.Fatalf("this case's premise moved: the boot snapshot already carries size %d",
				plantedSize)
		}

		// A hot-tier edit, planted while the process is alive.
		r.plant(t, "[fs]", fmt.Sprintf("[ball]\nsize = %d\n\n[fs]", plantedSize))
		applied := r.awaitAudit(t, "config: HOT-RELOAD state=applied")
		if !strings.Contains(applied, "hot=[ball]") {
			t.Errorf("the reload did not book ball as hot-tier: %q", applied)
		}
		if got := r.rt.mgr.Config().Ball.Size; got != plantedSize {
			t.Fatalf("live config ball.size = %d after the edit, want %d (the tier did not apply)",
				got, plantedSize)
		}
		if got := r.rt.cfg.Ball.Size; got == plantedSize {
			t.Fatalf("the boot snapshot moved too, so nothing here proves a mid-run change: %d", got)
		}
		// Ticket 255 AC#1 rewrote these three lines. They used to demand
		// 「这些段已立即生效」 for a [ball] edit - and [ball] has no production
		// reader (cmd/wisp/config_readers_255.go's `ball` row, 票 180's census
		// class "D"), so that sentence was precisely the lie 票 255 was立 for. The
		// rewrite is two-sided, not a loosening: the operator must still get a
		// user-visible line naming this edit, AND the forbidden claim must be absent.
		out := r.awaitStdout(t, "值已换进本进程内存")
		if !strings.Contains(out, "[ball]") {
			t.Errorf("no user-visible line names the hot edit at all; stdout:\n%s", out)
		}
		if strings.Contains(out, "这些段已立即生效") {
			t.Errorf("ticket 255 AC#1: a reader-less hot edit was reported as 已立即生效; stdout:\n%s", out)
		}
		// AC#3's other half, in one reading: a hot edit costs no card.
		if n := r.rt.windowCount(); n != 0 {
			t.Errorf("a hot-tier edit displayed %d confirmation cards, want 0", n)
		}
	})
}

// TestTicket223HandEditedFsLooseningCostsAnL2Card is AC#3's positive control and
// AC#6's production half: the re-confirmation travels through the hook cmd/wisp
// assigned, an unanswered card leaves the old values in place, and while the
// card is parked Config() is still readable - which is the whole reason
// manager.go had to be split.
func TestTicket223HandEditedFsLooseningCostsAnL2Card(t *testing.T) {
	r := newReloadRun223(t, 90*time.Second)
	r.live(t, func() {
		r.awaitAudit(t, "config: HOT-RELOAD state=armed")
		r.plant(t, allowedLine(r.h.dir), allowedLine(r.h.dir, r.h.outside))
		card := r.awaitCard(t, configReloadTool)
		if card.Level != "L2" {
			t.Errorf("the re-confirmation card is %s, want L2 (D33)", card.Level)
		}
		// The card text is what the operator reads; the ledger's ReplyCard has no
		// reason field, so the naming is asserted off the console.
		shown := r.h.out.String()
		if !strings.Contains(shown, "[确认 L2 "+configReloadTool+"]") {
			t.Errorf("the console did not render the reload card: %s", shown)
		}
		if !strings.Contains(shown, "fs.allowed_dirs") {
			t.Errorf("the card does not name the key it is asking about: %s", shown)
		}

		// AC#6's planted reading: the card is open, so this process is parked on
		// a human. Config() must still answer, and it must still answer with the
		// OLD values (D36 rule 1: 不得静默生效).
		type read struct{ dirs []string }
		ch := make(chan read, 1)
		go func() {
			// A bare `go` in a *_test.go: tools/d22scan's bans #1-5 skip test
			// files by design (票 223 census E.3 records that "测试里绿" proves
			// nothing about production shape), and this goroutine IS the probe -
			// routing it through the registry would hide a freeze behind a leak
			// report.
			ch <- read{dirs: r.rt.mgr.Config().FS.AllowedDirs}
		}()
		select {
		case got := <-ch:
			for _, d := range got.dirs {
				if filepath.ToSlash(d) == filepath.ToSlash(r.h.outside) {
					t.Fatalf("the pending loosening was live BEFORE it was confirmed: %v", got.dirs)
				}
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Config() froze while the D36 confirmation was parked - the hook is still under the Manager lock")
		}
		if n := r.rt.windowCount(); n != 1 {
			t.Errorf("the reload displayed %d cards, want exactly 1 for one loosening", n)
		}

		r.answer(t, "yes "+card.CorrelationID)
		answered := r.awaitAudit(t, "config: D36-CONFIRM state=answered section=fs")
		if !strings.Contains(answered, "result=allow") {
			t.Fatalf("the allow did not reach the ledger: %q", answered)
		}
		section := r.awaitAudit(t, "config: D36-SECTION section=fs direction=loosen")
		if !strings.Contains(section, "effect=applied-after-L2") {
			t.Errorf("an approved loosening is not booked as applied after the L2: %q", section)
		}
		r.awaitLive(t, "the confirmed loosening reaching memory", func() bool {
			return len(r.rt.mgr.Config().FS.AllowedDirs) == 2
		})
		if !strings.Contains(r.h.out.String(), "放宽已经过 L2 重新确认") {
			t.Errorf("no user-visible line says the loosening was confirmed and applied; stdout:\n%s",
				r.h.out.String())
		}
		// The honesty line about what did NOT move (the C26 canonicalizer).
		if !strings.Contains(r.h.out.String(), "仍按启动时建好的 C26 名单") {
			t.Errorf("an approved [fs] loosening did not say that this run's path verdicts are unchanged")
		}
	})
}

// TestTicket223RefusedLooseningKeepsOldValues is AC#3's fail-closed half on the
// same production route: the answer is "no", so the file keeps the operator's
// line and the process keeps the strict value - and the operator is told which
// one happened.
func TestTicket223RefusedLooseningKeepsOldValues(t *testing.T) {
	r := newReloadRun223(t, 90*time.Second)
	r.live(t, func() {
		r.plant(t, allowedLine(r.h.dir), allowedLine(r.h.dir, r.h.outside))
		card := r.awaitCard(t, configReloadTool)
		r.answer(t, "no "+card.CorrelationID+" 这次不要放宽")

		denied := r.awaitAudit(t, "config: D36-CONFIRM state=answered section=fs")
		if !strings.Contains(denied, "result=deny") {
			t.Fatalf("a refused card did not resolve as a deny: %q", denied)
		}
		section := r.awaitAudit(t, "config: D36-SECTION section=fs direction=loosen")
		if !strings.Contains(section, "effect=kept-old-values") {
			t.Errorf("a refused loosening is not booked as keeping the old values: %q", section)
		}
		if !strings.Contains(r.h.out.String(), "放宽本次没有生效") {
			t.Errorf("the operator is not told the loosening did not take effect; stdout:\n%s",
				r.h.out.String())
		}
		// Give the tick another cycle: the denial must be stable, not a race
		// before a late apply.
		time.Sleep(3 * time.Second)
		dirs := r.rt.mgr.Config().FS.AllowedDirs
		if len(dirs) != 1 || filepath.ToSlash(dirs[0]) != filepath.ToSlash(r.h.dir) {
			t.Fatalf("a refused loosening moved memory: %v", dirs)
		}
		// The file keeps what the operator wrote: refusing is not rewriting.
		raw, err := os.ReadFile(r.cfgPath())
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(raw), filepath.ToSlash(r.h.outside)) {
			t.Errorf("the refused hand edit was erased from config.toml instead of left alone:\n%s", raw)
		}
	})
}

// TestTicket223TighteningRaisesNoCard is AC#3's mandated reverse assertion: the
// tightening direction must hot-apply WITHOUT a card, or the security rail
// becomes harassment.
func TestTicket223TighteningRaisesNoCard(t *testing.T) {
	r := newReloadRun223(t, 90*time.Second)
	r.live(t, func() {
		r.plant(t, allowedLine(r.h.dir), allowedLine())
		line := r.awaitAuditLine(t, "config: D36-SECTION section=fs direction=tighten")
		if strings.Contains(line, "L2") {
			t.Errorf("a tightening was routed through the L2 confirmation: %q", line)
		}
		// The stronger form of the same claim: the D36 hook was never asked.
		if strings.Contains(r.h.err.String(), "config: D36-CONFIRM") {
			t.Errorf("the confirmation hook ran for a tightening; trail:\n%s", r.h.err.String())
		}
		if got := r.rt.mgr.Config().FS.AllowedDirs; len(got) != 0 {
			t.Fatalf("a tightening must hot-apply: %v", got)
		}
		if n := r.rt.windowCount(); n != 0 {
			t.Errorf("tightening [fs] displayed %d cards, want 0", n)
		}
		if strings.Contains(r.h.out.String(), configReloadTool) {
			t.Errorf("the reload card surface appeared for a tightening; stdout:\n%s", r.h.out.String())
		}
		if strings.Contains(r.h.err.String(), "state=denied") {
			t.Errorf("a tightening was booked as a denial: %s", r.h.err.String())
		}
	})
}

// TestTicket223ModeLooseningChangesTheRunningModeAfterAllow is AC#2's first tier
// read as BEHAVIOUR, not as a field: the permission mode is the live consumer of
// [risk] in this host (perm.Store answers the bridge from Manager.Config() once
// per tool call), so an approved loosening changes how the next call is gated
// without restarting the process.
func TestTicket223ModeLooseningChangesTheRunningModeAfterAllow(t *testing.T) {
	r := newReloadRun223(t, 90*time.Second)
	r.live(t, func() {
		if before := r.rt.modes.PermissionMode(); before != risk.ModeAskEveryStep {
			t.Fatalf("this case's premise moved: the run came up in %v, want the strictest档", before)
		}
		r.plant(t, "[risk]\n", "[risk]\npermission_mode = \"auto_approve\"\n")
		card := r.awaitCard(t, configReloadTool)
		shown := r.h.out.String()
		if !strings.Contains(shown, "risk.permission_mode") {
			t.Errorf("the card does not name risk.permission_mode: %s", shown)
		}
		if got := r.rt.modes.PermissionMode(); got != risk.ModeAskEveryStep {
			t.Fatalf("the running档 widened before anyone confirmed it: %v", got)
		}
		r.answer(t, "yes "+card.CorrelationID)
		r.awaitAudit(t, "config: D36-SECTION section=risk direction=loosen approved=true")
		r.awaitLive(t, "the confirmed [risk] loosening reaching the running档", func() bool {
			return r.rt.modes.PermissionMode() == risk.ModeAutoApprove
		})
	})
}

// TestTicket223RestartTierSaysItWillNotApply is AC#7: the restart tier has an
// exit now, and it names the three things silence hid - which sections, that
// this run will not use them, and why. Deliberately its own case: 票 223 forbids
// folding this tier into the immediately-effective one.
func TestTicket223RestartTierSaysItWillNotApply(t *testing.T) {
	r := newReloadRun223(t, 40*time.Second)
	r.live(t, func() {
		if r.rt.cfg.App.Autostart {
			t.Fatal("this case's premise moved: the boot snapshot already autostarts")
		}
		// 票 232 AC#2: the two marks are the operator stream and the audit trail
		// as they stand BEFORE the edit is planted. Everything asserted from a
		// window below is therefore bytes this run printed AFTER the plant, and
		// a start-up banner that happens to carry the same words (232's
		// forward-risk row, AC#4's mutation H) cannot stand in for the restart
		// notice.
		mark := r.h.out.String()
		errMark := r.h.err.String()
		// 票 272 AC#2, arm 1 - the one that reddens this ticket's M shape:
		// no byte printed BEFORE the plant may carry the operator's restart
		// sentence. Why this is its own arm instead of "move the window start
		// later" (232/272's shape 乙): the start-up banner is written
		// synchronously by startConfigReload (run.go:813) into a syncWriter
		// with no pipe and no copier in between
		// (approval_reply_201_test.go:64), so the mark taken just above already
		// excludes every start-up byte. Measured, not assumed - the
		// shape "banner carries the words, operator sentence deleted" reddens
		// on the UNMODIFIED ruler at its :592 (the same await sits at :622
		// below) with an EMPTY window after 40s in
		// .scratch/wisp/probes/272/r2/logs/AC1-oldruler-MDEL.txt, while the
		// same banner next to the real sentence passes (AC1-oldruler-M.txt),
		// which is exactly the ambiguity this arm removes: it bans the copy
		// itself, wherever in the pre-plant stream it sits. A copy that
		// arrives between this snapshot and the plant is inside the window
		// and is caught by arm 2's count. The normal run prints zero copies,
		// so this arm costs no green (AC1-oldruler-cur.txt PASS).
		prePlant := r.h.out.String()
		if n := strings.Count(prePlant, restartSentence272); n != 0 {
			t.Errorf("the operator stream carried %q %d time(s) BEFORE this run planted the edit; a copy nobody planted cannot be evidence that this run will not use the new values. stdout before the plant:\n%s",
				restartSentence272, n, prePlant)
		}
		r.plant(t, "[fs]", "[app]\nautostart = true\n\n[fs]")
		pending := r.awaitAuditSince232(t, errMark, "config: HOT-RELOAD state=restart-pending sections=[app]")
		if !strings.Contains(pending, "effect=next-process-start") {
			t.Errorf("the restart notice does not say when it takes effect: %q", pending)
		}
		why := r.awaitAuditSince232(t, errMark, "config: RESTART-PENDING detail=")
		// Ticket 223 r2 AC#5: the operator sentence lands on stdout AFTER the
		// two audit lines just awaited (config_reload.go:282/:284 -> :288),
		// so the single read that stood here caught the stream mid-write and
		// reddened this case intermittently. Poll for the sentence with
		// awaitStdout instead; every assertion below is unchanged verbatim,
		// and the positive control (deleting the product sentence) still
		// reddens - awaitStdout itself fails at its deadline.
		//
		// Ticket 232 AC#2 changed WHERE the sentence is awaited and WHAT the
		// content needles may read: awaitStdoutSince232 finds it in the window
		// that starts at `mark`, so bytes printed before the plant cannot
		// satisfy it, and every content needle below is pinned to ONE NAMED
		// stream instead of to a `why+out` concatenation - which is what let
		// 232-v2's mutation B delete all three explanation clauses from the
		// operator's sentence and stay green off the audit trail alone. The
		// awaits are unchanged in shape and deadline (reloadCaseBudget,
		// time.NewTimer + 20ms time.NewTicker), nothing was removed and nothing
		// became an either/or (AC#5); the window reads are all STRICTER.
		out := r.awaitStdoutSince232(t, mark, "本次运行不会生效")
		// Which stream actually carries each needle was MEASURED, not assumed:
		// .scratch/wisp/probes/232/r2/logs/P0-needle-stream-probe.txt (HEAD
		// product, read-only probe injected by overlay) says "app.autostart" and
		// "开机自启" are on the operator's stdout and so are the three clause
		// markers below, while "重启进程后生效" lives ONLY on the audit line -
		// reportRestartPending's Fprintf (config_reload.go:321-326) never says
		// it. Pinning that one to stdout would mean writing it into the product
		// sentence, which 票 232's 禁区 forbids, so it is pinned to the named
		// audit stream instead - read from the post-plant window, so a banner
		// cannot satisfy it either. "需要重启进程" is the stdout-side equivalent
		// of the same promise and is pinned here so the operator's half carries
		// it too.
		for _, needle := range []string{"app.autostart", "开机自启", "需要重启进程", "原因：", "交给平台层", "没有被丢掉"} {
			if !strings.Contains(out, needle) {
				t.Errorf("the operator's restart sentence omits %q; stdout since the plant:\n%s", needle, out)
			}
		}
		if !strings.Contains(why, "重启进程后生效") {
			t.Errorf("the restart notice never says %q on the audit trail either; audit since the plant:\n%s\nstdout since the plant:\n%s",
				"重启进程后生效", why, out)
		}
		if got := r.rt.mgr.Config().App.Autostart; got {
			t.Fatal("the restart tier applied mid-run, which is exactly what D36 forbids")
		}
		if !strings.Contains(out, "本次运行不会生效") {
			t.Errorf("the operator is not told the edit will not land this run; stdout:\n%s", out)
		}
		// 票 272 AC#2, arm 2 - the count the ticket names as shape 甲: the
		// window since the mark must carry EXACTLY ONE copy of the sentence.
		// Zero means the notice never arrived, two means something else said
		// it too, and the ruler before this ticket asked neither: it only ever
		// asked "contains", so a second verbatim copy of the operator's own
		// Fprintf passed it in 2.22s (shape REPEAT,
		// .scratch/wisp/probes/272/r2/logs/AC1-oldruler-REPEAT.txt).
		// The count is taken off a re-read made AFTER the audit line
		// state=applied, because reportReload books that line only once
		// CheckAndReload has returned (config_reload.go:179) while the
		// operator's sentence is written inside that very call
		// (config_reload.go:321-326) - so by then every stdout byte this plant
		// causes is already in the stream and the window cannot be counted
		// mid-sentence. This is not a window-start move (票 272 AC#3 forbids
		// loosening the start): `mark` still opens the window, the await above
		// still ends it, and this line only re-reads what is already there.
		applied := r.awaitAuditSince232(t, errMark, "config: HOT-RELOAD state=applied")
		if win := r.stdoutSince232(t, mark); strings.Count(win, restartSentence272) != 1 {
			t.Errorf("the window since the mark carries %q %d time(s), want exactly one copy caused by this plant; stdout since the plant:\n%s\naudit since the plant:\n%s",
				restartSentence272, strings.Count(win, restartSentence272), win, applied)
		}
		// Two tiers, two sentences: the restart notice must not be dressed up as
		// an immediate one. Read over the WHOLE stream on purpose - a windowed
		// read here would be narrower than the one this assertion has always
		// had, and 票 232 AC#5 forbids loosening.
		if whole := r.h.out.String(); strings.Contains(whole, "这些段已立即生效：[app]") {
			t.Errorf("the restart tier was reported as immediately effective; stdout:\n%s", whole)
		}
		if n := r.rt.windowCount(); n != 0 {
			t.Errorf("a restart-tier edit displayed %d cards, want 0 (nothing loosened)", n)
		}
	})
}

// TestTicket223FailureSentencesAreDistinct is AC#4: 缺失 / 语法错 / schema 拒绝
// each get their own sentence and must not borrow each other's. The fourth
// required sentence (热加载被禁用) is a host state rather than a read failure and
// lives in TestTicket223PanelInboundSaysHotReloadIsDisabled; the third one
// (权限不够) needs a Windows ACL and lives in
// TestTicket223PermissionDeniedSitsInItsOwnSentence.
func TestTicket223FailureSentencesAreDistinct(t *testing.T) {
	cases := []struct {
		name      string
		plant     func(t *testing.T, r *reloadRun223)
		wantCause string
	}{
		{
			name: "缺失",
			plant: func(t *testing.T, r *reloadRun223) {
				if err := os.Remove(r.cfgPath()); err != nil {
					t.Fatalf("remove config.toml: %v", err)
				}
			},
			wantCause: "cause=missing",
		},
		{
			name: "语法错",
			plant: func(t *testing.T, r *reloadRun223) {
				// Deliberately WITHOUT a schema_version line: this is
				// loader.go's branch 1 - bytes from which not even a version
				// can be picked. Which sentence a file WITH a declared
				// version gets is decided by which pipeline owns its error,
				// and 票 223 r2 narrowed that boundary: declared OLDER
				// (below current) belongs to the migration pipeline - that
				// is what migrate_test.go:123 pins - while declared CURRENT
				// (or newer) with an unparseable body belongs to 语法错,
				// because a file already carrying this build's version has
				// no migration to run. Both halves are pinned by
				// TestTicket223R2FailureSentenceRouting. "Not TOML at all,
				// not even a version" is what this sub-case means.
				r.writeOver(t, "this is not toml [[[\n")
			},
			wantCause: "cause=syntax",
		},
		{
			name: "声明了版本但坏在后面",
			plant: func(t *testing.T, r *reloadRun223) {
				r.writeOver(t, "schema_version = 1\n[llm\nbroken ===\n")
			},
			wantCause: "cause=migration",
		},
		{
			name: "schema未知键",
			plant: func(t *testing.T, r *reloadRun223) {
				r.writeOver(t, "schema_version = 2\n\nthis_key_does_not_exist = 1\n")
			},
			wantCause: "cause=unknown-key",
		},
	}
	// The whole family, so each case can assert that the OTHER sentences are
	// absent instead of trusting that one shared line means one shared cause.
	// 票 231 AC#2's cause joins the roster: a fifth "read refused" shape that is
	// not in here is a shape nothing protects from swallowing a neighbour.
	all := []string{
		"cause=missing", "cause=syntax", "cause=unknown-key", "cause=invalid",
		"cause=permission", "cause=unclassified", "cause=migration",
		"cause=newer-build", "state=disabled",
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newReloadRun223(t, 40*time.Second)
			r.live(t, func() {
				r.awaitAudit(t, "config: HOT-RELOAD state=armed")
				before := r.rt.mgr.Config()
				tc.plant(t, r)
				trail := r.awaitAudit(t, "config: HOT-RELOAD state=not-applied "+tc.wantCause)
				for _, other := range all {
					if other == tc.wantCause {
						continue
					}
					if strings.Contains(trail, other) {
						t.Errorf("the %s reading also carries %q, which is a different cause:\n%s",
							tc.name, other, trail)
					}
				}
				// "读不到" and "生效了" cannot both be true: the failed reload
				// must have left the running config where it was.
				now := r.rt.mgr.Config()
				if now.Ball.Size != before.Ball.Size ||
					len(now.FS.AllowedDirs) != len(before.FS.AllowedDirs) {
					t.Errorf("a failed reload moved the live config: ball %d vs %d, dirs %v vs %v",
						now.Ball.Size, before.Ball.Size, now.FS.AllowedDirs, before.FS.AllowedDirs)
				}
			})
		})
	}
}

// TestTicket223PanelInboundSaysHotReloadIsDisabled is AC#4's fourth sentence on
// a real production host that does not own a tick: `wisp panel-inbound` opens
// the same config.toml, has no approval gate to confirm a loosening on, and says
// so instead of leaving an edit to die silently.
func TestTicket223PanelInboundSaysHotReloadIsDisabled(t *testing.T) {
	dir := newInboundDataDir33(t, risk.ModeAskEveryStepName)
	code, out, errLog := runInboundLeg33(t, "", dir)
	if code != 0 {
		t.Fatalf("panel-inbound exit %d\nstdout:\n%s\nstderr:\n%s", code, out, errLog)
	}
	if !strings.Contains(errLog, "config: HOT-RELOAD state=disabled host=panel-inbound") {
		t.Fatalf("the host that owns no tick did not say so; stderr:\n%s", errLog)
	}
	for _, other := range []string{"cause=missing", "cause=syntax", "cause=permission", "state=armed"} {
		if strings.Contains(errLog, other) {
			t.Errorf("the disabled sentence carries %q, which belongs to another cause:\n%s", other, errLog)
		}
	}
}
