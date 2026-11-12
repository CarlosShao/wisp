package main

// The [memory] / [privacy] retention and L1 knobs (the 180-census follow-up's
// wire, ticket 423's second wave). The census measured all of these as
// zero-reader; the production readers this file carries are the two feeds the
// assembly root and the resident leg make:
//
//   - retentionWindowDays180K resolves the L3 window both sections name into
//     one number for memory.StartRetentionJob (run.go's assembleRuntime owns
//     the start; before this wire the job had NO production owner at all, so
//     no row ever aged out - D35 rule 1's exact "silently never run" shape).
//   - memoryL1Enabled180K is the resident gate's fresh read: [memory]
//     l1_enabled=false means the L1 extraction pipeline is not attached, with
//     the boot line saying so instead of the window being silently empty.

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/config"
)

// l3RetentionDefaultDays is the D35 default the two schema keys carry
// (memory.l3_retention_days and privacy.retention_days, both 30; SPEC-02 §4's
// task-log window).
const l3RetentionDefaultDays = 30

// retentionWindowDays180K resolves the L3 retention window in days.
//
// D36 lists the window under BOTH sections - [privacy] retention_days
// ("the L3/forensics retention window", schema comment) and [memory]
// l3_retention_days - and the store has one task_log table, so the two keys
// feed one knob. The rule when they disagree is stated here rather than
// invented silently: the SMALLER window wins, because the window's purpose is
// a privacy promise (forensics rows do not outlive it) and the more
// protective reading is the one neither section contradicts. A key at or
// below zero is treated as "not written" and keeps the D35 default, so a
// hand-typed 0 cannot turn the retention pass into "delete everything now".
func retentionWindowDays180K(c *config.Config) int {
	days := l3RetentionDefaultDays
	if c.Memory.L3RetentionDays > 0 && c.Memory.L3RetentionDays < days {
		days = c.Memory.L3RetentionDays
	}
	if c.Privacy.RetentionDays > 0 && c.Privacy.RetentionDays < days {
		days = c.Privacy.RetentionDays
	}
	return days
}

// memoryL1Enabled180K is the l1_enabled gate's feed: one fresh config.LoadFile
// over the data root (the per-use shape the resident boot uses everywhere it
// reads hot-tier views). Missing/unreadable file = the compiled default table
// (l1_enabled=true), which is the pre-wire shape, and the source string names
// which table decided.
func memoryL1Enabled180K(dataDir string) (enabled bool, source string) {
	cfgPath := filepath.Join(dataDir, configFileName)
	c, _, err := config.LoadFile(cfgPath, nil)
	if err != nil || c == nil {
		c = config.NewDefaults()
		slog.Warn("memory: config unreadable at the l1_enabled gate; the compiled default table decides",
			"path", cfgPath, "err", configErrText(err, c), "l1_enabled", c.Memory.L1Enabled)
		return c.Memory.L1Enabled, "compiled defaults (config.NewDefaults)"
	}
	return c.Memory.L1Enabled, cfgPath
}

// retentionBootLine180K is the one audit sentence the retention start books,
// lifted out so the feed test can assert the same numbers the log carries.
func retentionBootLine180K(days int) string {
	return fmt.Sprintf("wisp run: retention job owned（[memory] l3_retention_days / [privacy] retention_days → task_log 窗口 %d 天）", days)
}
