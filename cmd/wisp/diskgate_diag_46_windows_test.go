//go:build windows

package main

// Ticket 46's post-audit teeth for items 1 and 2 (the audit件
// docs/evidence/s1/46b-audit-20261115.md, "必须回退的" #1 and "另记 2 条" #1).
//
// Both are BEHAVIOUR teeth (A1880): each drives the production function and
// asserts a reading, and each names in its red sentence what is missing - never
// "a string in a file changed". Deleting the production wiring turns each one
// red for the reason the wiring exists.
//
// Item 1's tooth: the disk gate's boot reading. The production reader is the
// line built from the pipeline the boot owns (diskHeldStatsLine46 over
// sink.pipeline46(), printed unconditionally at resident boot). Three
// behaviours are pinned:
//
//   - a healthy gated pipeline answers the 未被扣住 sentence (the
//     always-one-line contract the dispatch demands: normal has its own words);
//   - obs.Disk = nil (the guard never assembled) answers the 未装门控 sentence
//     - the reading CHANGES when the gate is absent, which is exactly the
//     mutation the audit's ruling names;
//   - a tripped gate answers the held/lost/replayed sentence with real counts.
//
// Item 2's tooth: the boot row's three states. An assembled-but-unmounted
// gate (no voice residency: first install without `wisp models ensure` is the
// common shape) must say so with the named reason instead of the assembled
// sentence - the row the audit measured as a false claim.

import (
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// diskHeldSentenceHeld46 is the held/lost/replayed sentence's stable head: the
// reading an operator greps under pressure.
const diskHeldSentenceHeld46 = "日志磁盘门控读数：日志被扣住"

// Test46DiskHeldStatsLineChangesWhenTheGateIsAbsent drives the production
// reader (diskHeldStatsLine46) over the pipeline the production boot path
// builds (installLogSink with LogConfig.Disk fed by diskGuardForDataDir). The
// mutation the audit's ruling names - obs.Disk = nil - is realized here by the
// fixture seam that DOES exist: diskGuardForDataDir answers no guard, which is
// exactly what a nil obs.Disk produces in InitLog (logging.go:149's arm never
// runs, so the sink carries no gated writer and DiskHeldStats answers
// ok=false).
func Test46DiskHeldStatsLineChangesWhenTheGateIsAbsent(t *testing.T) {
	prevDefault := slog.Default()
	t.Cleanup(func() { slog.SetDefault(prevDefault) })

	// (a) the gate assembled (production seam, real guard, healthy volume):
	// the row must say 未被扣住 - the normal posture has its own sentence.
	probe := &diskProbe46{}
	probe.free.Store(2 << 30)
	guard, err := observe.NewDiskGuard(observe.DiskGuardOptions{
		Volume: t.TempDir(),
		Probe:  probe.probe,
	})
	if err != nil {
		t.Fatalf("46: fixture guard: %v", err)
	}
	restore := swapDiskGuard46(func(string) (*observe.DiskGuard, string) {
		return guard, "46 fixture probe (2GiB free)"
	})
	t.Cleanup(restore)

	healthyDir := t.TempDir()
	healthy, err := installLogSink(healthyDir)
	if err != nil {
		t.Fatalf("46: installLogSink (healthy): %v", err)
	}
	t.Cleanup(func() { healthy.close() })
	healthyLine := diskHeldStatsLine46(healthy.pipeline46())
	if !strings.Contains(healthyLine, "未被扣住") {
		t.Fatalf("46: the healthy gated pipeline's row = %q, want the 未被扣住 sentence: a line that only speaks under pressure teaches no baseline and the gate's existence is unprovable from the boot output", healthyLine)
	}
	if strings.Contains(healthyLine, diskHeldSentenceHeld46) {
		t.Fatalf("46: the healthy gated pipeline's row = %q, want no held/lost/replayed counts: a healthy gate is not a tripped gate", healthyLine)
	}

	// (b) obs.Disk = nil, by the assembly's own seam: no guard, so the sink
	// carries no gated writer and the row must CHANGE - naming 未装门控 instead
	// of pretending a gate exists. This is the tooth's load-bearing half: with
	// the old shape (ok=false reading rendered as zeros, or output only when
	// abnormal) this branch cannot be told apart from (a).
	restore2 := swapDiskGuard46(func(string) (*observe.DiskGuard, string) {
		return nil, "46 fixture: probe failed at install, the log writer is ungated"
	})
	t.Cleanup(func() { restore2() })

	ungatedDir := t.TempDir()
	ungated, err := installLogSink(ungatedDir)
	if err != nil {
		t.Fatalf("46: installLogSink (ungated): %v", err)
	}
	t.Cleanup(func() { ungated.close() })
	ungatedLine := diskHeldStatsLine46(ungated.pipeline46())
	if !strings.Contains(ungatedLine, "未装门控") {
		t.Fatalf("46: the ungated pipeline's row = %q, want the 未装门控 sentence: with the disk guard absent (obs.Disk = nil) the boot output must change, otherwise an absent gate and a healthy gate read alike", ungatedLine)
	}

	// (c) the tripped posture: the held/lost/replayed counts read off the
	// running gate (the operator's "is the log WAITING or LOST" question).
	trippedProbe := &diskProbe46{}
	trippedProbe.free.Store(1 << 20)
	trippedGuard, err := observe.NewDiskGuard(observe.DiskGuardOptions{
		Volume: t.TempDir(),
		Probe:  trippedProbe.probe,
	})
	if err != nil {
		t.Fatalf("46: fixture guard (tripped): %v", err)
	}
	restore3 := swapDiskGuard46(func(string) (*observe.DiskGuard, string) {
		return trippedGuard, "46 fixture probe (1MB free)"
	})
	t.Cleanup(func() { restore3() })

	trippedDir := t.TempDir()
	tripped, err := installLogSink(trippedDir)
	if err != nil {
		t.Fatalf("46: installLogSink (tripped): %v", err)
	}
	t.Cleanup(func() { tripped.close() })
	trippedLine := diskHeldStatsLine46(tripped.pipeline46())
	if !strings.Contains(trippedLine, diskHeldSentenceHeld46) {
		t.Fatalf("46: the tripped pipeline's row = %q, want the held/lost/replayed sentence: the reading nobody but the gate can give is exactly what the boot row must carry under pressure", trippedLine)
	}
	if !strings.Contains(trippedLine, "丢弃 0 条") {
		t.Fatalf("46: the tripped pipeline's row = %q, want the 丢弃 count in it: a held tail without the lost count hides the bytes the bound dropped", trippedLine)
	}

	// (d) no pipeline at all (the install refusal path): the row names the
	// absence instead of panicking or pretending.
	nilLine := diskHeldStatsLine46(nil)
	if !strings.Contains(nilLine, "无") {
		t.Fatalf("46: the no-pipeline row = %q, want the absence sentence: the installLogSink refusal branch must not read like a healthy gate", nilLine)
	}
}

// Test46BootRowSaysAssembledButUnmountedWithoutAResidency is item 2's tooth:
// the three-state boot row. A gate handed to a boot with no voice residency
// (first install without `wisp models ensure`) must print
// 已装配但未挂载 with the named reason; the mounted gate keeps the assembled
// sentence; a gate over no watchdog loop keeps the 未装配 sentence.
//
// Red when the mount half is dropped: with the old two-state line() an
// unmounted gate answered the assembled sentence (the audit's false claim),
// and this case's first assertion goes red on exactly that.
func Test46BootRowSaysAssembledButUnmountedWithoutAResidency(t *testing.T) {
	// (a) assembled but unmounted: the production attach half runs and finds
	// nothing to mount on - the voice leg never built a residency. The verdict
	// is the named reason the row carries.
	rv := &residentVoice{verdict: "语音链未构造：[voice] enabled=false（配置来源 config.toml）"}
	tree := &settleTree46{bytes: 30 << 20}
	reg := observe.NewRegistry()
	root := observe.NewRoot("46-settle-line-unmounted")
	t.Cleanup(root.Cancel)
	gate := newResidentSettleGate(settleLoop46(t, tree, reg, 50*time.Millisecond), tree, reg, root, nil)
	if rv.attachSettleGate(gate) {
		t.Fatal("46: the production attach half answered mounted on a boot with no residency: the dispose cascade could never reach the gate, so the boot row would have a mount reading that is a lie")
	}
	line := gate.line()
	if !strings.Contains(line, "已装配但未挂载") {
		t.Fatalf("46: an assembled-but-unmounted gate's row = %q, want the 已装配但未挂载 sentence: with the old two-state line this boot claimed 装配 while the dispose cascade could never reach the gate (the audit's false claim), so the first install without `wisp models ensure` lied at every start", line)
	}
	if !strings.Contains(line, rv.verdict) {
		t.Fatalf("46: the unmounted row = %q, want the boot's own named reason (%q): an unnamed reason is a second false claim - it does not point at the actual cause", line, rv.verdict)
	}

	// (b) the mount lands: the row says 装配 again, through the same
	// production attach half (the residency object is the one the cascade
	// reads, settleGateNow).
	rv2 := &residentVoice{residency: newVoiceResidency(voicePlan{}, newVoiceOwnerQueue(), nil)}
	if !rv2.attachSettleGate(gate) {
		t.Fatal("46: the production attach half refused to mount on a real residency: the dispose cascade would verify nothing")
	}
	mountedLine := gate.line()
	if !strings.Contains(mountedLine, "收口闸：装配") {
		t.Fatalf("46: the mounted gate's row = %q, want the 装配 sentence: a real residency must not be described as unmounted or the operator stops trusting the row", mountedLine)
	}

	// (c) the no-watchdog state is unchanged: 未装配, verbatim shape.
	if got := (*residentSettleGate)(nil).line(); !strings.Contains(got, "未装配") {
		t.Fatalf("46: the no-watchdog row = %q, want the 未装配 sentence: a process with no loop has no verifier and the row must keep saying so", got)
	}
}
