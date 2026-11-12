//go:build windows

package main

import (
	"fmt"
	"log/slog"
)

// attachMemoryL1IfEnabled180K is the resident boot's l1_enabled gate, the
// addressable seam the wiring test drives: enabled runs the ticket-29 attach,
// disabled leaves the WarmWindow hook body-less and says the key decided
// that. It runs where the assembly root used to call attachMemoryL1
// unconditionally (resident_windows.go).
//
// Windows-only by build tag: it is the one retention wire whose parameter
// list names *residentVoice, a type that only exists behind the windows
// build (resident_voice_windows.go). Everything else in the 180-census
// retention wire is cross-platform and stays in memory_retention_180k.go.
func attachMemoryL1IfEnabled180K(rv *residentVoice, run *agentRuntime, dataDir string) {
	l1, src := memoryL1Enabled180K(dataDir)
	if !l1 {
		slog.Info("memory: L1 extraction not attached",
			"reason", "[memory] l1_enabled=false", "config_source", src)
		fmt.Printf("wisp: 记忆抽取未接（[memory] l1_enabled=false，来源 %s）：Warm 窗口不挂抽取管线，画像不学习\n", src)
		return
	}
	rv.attachMemoryL1(run)
}
