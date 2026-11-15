//go:build windows

package main

import (
	"fmt"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 46's post-audit leg, item 1: DiskHeldStats gets its production reader.
//
// The audit (docs/evidence/s1/46b-audit-20261115.md §3) judged this reading
// "不成立" with a named reason: internal/observe/logging.go's
// LogPipeline.DiskHeldStats had ZERO non-test callers, so the only numbers that
// can tell "the file is empty because bytes are HELD" apart from "the file is
// empty because nothing was written yet" existed inside test processes only.
// The audit's own ruling (echoed by 46-s6-acceptance.md 审计收件 §裁定 1) is to
// give it a reader, not to delete it.
//
// WHY THE BOOT ROW AND NOT `wisp doctor` / `wisp diagnostics`. Both of those
// are SEPARATE processes: a CLI doctor cannot read the resident process's
// in-memory gate, and the bundle's disk-guard entry (BundleOptions.DiskStatus)
// exists for a caller that owns a *DiskGuard at build time - the CLI legs
// install no pipeline at all (main.go:173's own note: the diagnostics leg
// "reaches no installLogSink"). The one process that owns the gated writer is
// the resident (installLogSink, resident_windows.go), so the one honest
// production reader of the held/lost/replayed counters is this boot: print the
// reading at the moment the pipeline exists, every start.
//
// WHY ALWAYS-ONE-LINE. The dispatch forbids the shape "output only when
// abnormal": a line that appears only under pressure teaches the operator
// nothing about the healthy baseline and cannot be grep-diffed to prove the
// gate exists. Normal is therefore its own sentence, and a nil pipeline (the
// installLogSink refusal branch) says the gate is absent rather than pretending
// a reading.
//
// The counts come from LogPipeline.DiskHeldStats, which returns ok=false when
// the sink is not disk-gated (guard absent: the probe failed at install, or no
// data dir) - that posture gets its own sentence, because "ungated" and "held
// 0" must not read alike (the same distinction the install record's
// disk_gated/disk_guard attributes make).
func diskHeldStatsLine46(p *observe.LogPipeline) string {
	if p == nil {
		return "日志磁盘门控读数：无（持久日志管线未装配，本次启动没有可读的门控计数）"
	}
	held, lost, replayed, gated := p.DiskHeldStats()
	if !gated {
		return "日志磁盘门控读数：未装门控（安装时未装配磁盘守卫，日志写入不经任何扣留/回放门）"
	}
	if held == 0 && lost == 0 && replayed == 0 {
		return "日志磁盘门控读数：未被扣住（门控健康：扣住 0 条／丢弃 0 条／回放 0 条）"
	}
	return fmt.Sprintf("日志磁盘门控读数：日志被扣住 %d 条／丢弃 %d 条／回放 %d 条", held, lost, replayed)
}

// pipeline46 is the nil-tolerant read of the boot sink's pipeline: the
// installLogSink refusal path leaves the sink nil, and the boot row must print
// the "no pipeline" sentence instead of panicking on the way there.
func (s *logSink) pipeline46() *observe.LogPipeline {
	if s == nil {
		return nil
	}
	return s.pipeline
}
