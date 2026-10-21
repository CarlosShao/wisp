package main

// Ticket 358 AC#2: `wisp run` must be able to say "该压但压不动" out loud.
//
// WHAT WAS WRONG. Until ticket 358 the report a run carried had one boolean,
// Ran, and the CLI line keyed on it: a pass that folded nothing and a run that
// never came near the budget were the same bytes on the report, so the line
// could only be silent for both. The compressor now carries a
// CompressionState (skipped/attempted/compacted - internal/agent/compress.go,
// shaped by ticket 358 AC#1), which types the two shapes apart. This file pins
// what the CLI does with each of the three, one case per shape, through the
// shipped formatter - never a copy of it.
//
// WHY IT KEYS ON State AND NOT ON Ran. Ran is kept as the compatibility mirror
// of the compacted state (rep.Ran == (rep.State == CompressionCompacted)); its
// meaning on the existing readers is unchanged, which is why the line reads the
// state and leaves Ran alone. The compacted wording is pinned BYTE FOR BYTE -
// the expected string below is the one ticket 355 shipped - because "readable"
// was never a licence to reword a line that already worked.
//
// Every threshold is read back out of agent.BudgetsFor; no number is written by
// hand (ticket 139 AC#1(2)). The reports handed to the formatter are produced by
// the real compressor over real histories, not by struct literals, so a fixture
// that quietly stopped clearing the threshold would redden rather than pass
// tautologically. buildTraceHistory and realCompressionReport are the same
// fixtures ticket 355's file already uses.

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
)

// compactedLine355 is ticket 355's shipped wording, spelled out here rather than
// re-derived from the production format string, so rewording the line (or
// reordering its fields) is red instead of silently accepted. Its arguments are
// the report's own values, so this is not a hand-typed expected string.
func compactedLine355(rep agent.CompressionReport) string {
	return fmt.Sprintf(
		"wisp run: 压缩痕：压前 %d tokens => 压后 %d tokens（history_changed=%t，%d 条消息被折进摘要，保留 %d 轮原文）",
		rep.TokensBefore, rep.TokensAfter, rep.HistoryChanged,
		rep.CompressedMsgs, rep.KeptRawRounds)
}

// mustLine358 is a report's line, or a fatal: several assertions below compare
// two lines and a silent formatter would make them compare equal for the wrong
// reason.
func mustLine358(t *testing.T, rep agent.CompressionReport) string {
	t.Helper()
	line, ok := compressionLine(rep)
	if !ok {
		t.Fatalf("compressionLine answered nothing for state %q", rep.State)
	}
	return line
}

// TestCompressionTraceCompactedLineIsTheOne355Shipped is shape one: a pass that
// really folded history prints exactly the line ticket 355 delivered, with the
// report's own numbers in it.
func TestCompressionTraceCompactedLineIsTheOne355Shipped(t *testing.T) {
	rep := realCompressionReport(t)
	if rep.State != agent.CompressionCompacted {
		t.Fatalf("setup: State = %q for a pass that folded, want %q", rep.State, agent.CompressionCompacted)
	}

	line, ok := compressionLine(rep)
	if !ok {
		t.Fatal("compressionLine answered nothing for a pass that compacted")
	}
	want := compactedLine355(rep)
	if line != want {
		t.Errorf("the compacted line changed (ticket 355's wording is pinned byte for byte):\n got: %s\nwant: %s",
			line, want)
	}

	var out bytes.Buffer
	if printed := printCompressionTrace(&out, rep); !printed {
		t.Fatal("printCompressionTrace reported that it printed nothing for a pass that compacted")
	}
	if got := out.String(); got != want+"\n" {
		t.Errorf("what reached stdout is not the pinned line plus one newline:\n got: %q\nwant: %q",
			got, want+"\n")
	}
}

// attemptedAtKeepRawFloor drives the shape ticket 153 named: the history clears
// the threshold but only one raw round exists, so the fold loop breaks at the
// D15(4) floor immediately. It fails rather than returning a zero report,
// because a fixture that stopped clearing the threshold would turn the
// assertions below into tautologies.
func attemptedAtKeepRawFloor(t *testing.T) agent.CompressionReport {
	t.Helper()
	b := agent.BudgetsFor(4096)
	c := agent.NewCompressor(b, nil)
	hist := buildTraceHistory(1, b.HistoryCompressTokens*4)
	if !c.Need(hist) {
		t.Fatalf("setup: one padded round of %d tokens does not clear %d",
			c.TotalTokens(hist), b.HistoryCompressTokens)
	}
	_, rep, err := c.Compress(context.Background(), hist)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if rep.State != agent.CompressionAttempted {
		t.Fatalf("setup: State = %q for the over-threshold/no-foldable shape, want %q",
			rep.State, agent.CompressionAttempted)
	}
	if rep.TokensBefore <= b.HistoryCompressTokens {
		t.Fatalf("setup: this fixture was meant to be over the fold budget, got %d against %d",
			rep.TokensBefore, b.HistoryCompressTokens)
	}
	return rep
}

// TestCompressionTraceAttemptedSaysItFoldedNothing is shape two, and the point
// of ticket 358: the trigger cleared the threshold and the pass folded nothing,
// so the run must SAY that - with the four readings that make the statement
// checkable - instead of falling silent.
//
// Both shapes that reach this state are driven here, because they are why a line
// of zeros is not enough. In the first the fold loop ran and hit the
// KeepRawRounds floor. In the second the provider's count cleared the threshold
// while the pass's own estimate did not, so the loop never took a round and
// tokens_before == tokens_after == a number BELOW the threshold: a reader seeing
// only those zeros would conclude nothing was attempted, which is exactly the
// mistake this ticket is about.
func TestCompressionTraceAttemptedSaysItFoldedNothing(t *testing.T) {
	floor := attemptedAtKeepRawFloor(t)

	b := agent.BudgetsFor(4096)
	c := agent.NewCompressor(b, nil)
	hist := buildTraceHistory(2, 8)
	observed := b.HistoryCompressTokens*3 + 1
	_, under, err := c.Compress(context.Background(), hist,
		agent.WithTrigger(agent.TriggerReading{
			Need:      true,
			Tokens:    observed,
			Threshold: b.HistoryCompressTokens,
			Source:    agent.TriggerProviderUsage,
			Observed:  observed,
		}))
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if under.State != agent.CompressionAttempted {
		t.Fatalf("setup (ii): State = %q, want %q: the provider count cleared the threshold over a "+
			"history the estimate calls small", under.State, agent.CompressionAttempted)
	}
	if under.TokensBefore > b.HistoryCompressTokens {
		t.Fatalf("setup (ii): this fixture was meant to stay under the fold budget, got %d against %d",
			under.TokensBefore, b.HistoryCompressTokens)
	}
	if under.TriggerSource != agent.TriggerProviderUsage {
		t.Fatalf("setup (ii): TriggerSource = %q, want %q", under.TriggerSource, agent.TriggerProviderUsage)
	}

	for _, tc := range []struct {
		name string
		rep  agent.CompressionReport
	}{
		{"keep_raw_rounds_floor", floor},
		{"provider_triggered_under_fold_budget", under},
	} {
		t.Run(tc.name, func(t *testing.T) {
			line, ok := compressionLine(tc.rep)
			if !ok || line == "" {
				t.Fatalf("an attempt that folded nothing answered %q (ok=%v); state=%q means the run "+
					"DID try, so silence here is the hole ticket 358 exists for",
					line, ok, tc.rep.State)
			}
			if !strings.Contains(line, "超了阈值但一枚都没折") {
				t.Errorf("the attempted line does not say the attempt folded nothing: %s", line)
			}
			for _, never := range []string{"没触发", "未触发", "跳过"} {
				if strings.Contains(line, never) {
					t.Errorf("the attempted line claims %q, which is the skipped shape's answer: %s", never, line)
				}
			}
			for _, want := range []string{
				"tokens_before=" + strconv.Itoa(tc.rep.TokensBefore),
				"tokens_after=" + strconv.Itoa(tc.rep.TokensAfter),
				"trigger=" + strconv.Itoa(tc.rep.TriggerTokens) + " tokens",
				"source=" + string(tc.rep.TriggerSource),
			} {
				if !strings.Contains(line, want) {
					t.Errorf("the attempted line does not carry %q: %s", want, line)
				}
			}
			// Privacy: counts and a source name only, never the folded text.
			if strings.Contains(line, "问题") || strings.Contains(line, strings.Repeat("x", 40)) {
				t.Errorf("the attempted line leaked history content: %s", line)
			}

			var out bytes.Buffer
			if printed := printCompressionTrace(&out, tc.rep); !printed {
				t.Fatal("printCompressionTrace wrote nothing for an attempt that folded nothing")
			}
			if got := out.String(); got != line+"\n" {
				t.Errorf("what reached stdout is not the formatter's line plus one newline:\n got: %q\nwant: %q",
					got, line+"\n")
			}
		})
	}

	// The two shapes must not collapse into one reading: same state, different
	// numbers, so the line has to carry the difference.
	if a, u := mustLine358(t, floor), mustLine358(t, under); a == u {
		t.Errorf("both attempted shapes printed the same line %q, so the numbers on it are not the report's", a)
	}
}

// TestCompressionTraceSkippedStaysSilent is shape three: the trigger never said
// the history was over budget, so nothing was attempted and the run prints
// nothing at all. It pins the FILLED skipped report as well as the zero value,
// so the silence cannot be satisfied by special-casing the empty struct.
func TestCompressionTraceSkippedStaysSilent(t *testing.T) {
	b := agent.BudgetsFor(4096)
	c := agent.NewCompressor(b, nil)
	hist := buildTraceHistory(2, 8)
	if c.Need(hist) {
		t.Fatalf("setup: a %d-token history clears the %d threshold, so this is not the skipped shape",
			c.TotalTokens(hist), b.HistoryCompressTokens)
	}
	_, rep, err := c.Compress(context.Background(), hist)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if rep.State != agent.CompressionSkipped {
		t.Fatalf("setup: State = %q for a history under the threshold, want %q",
			rep.State, agent.CompressionSkipped)
	}

	for _, tc := range []struct {
		name string
		rep  agent.CompressionReport
	}{
		{"filled_skipped_report", rep},
		{"zero_value_no_pass_ran", agent.CompressionReport{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			line, ok := compressionLine(tc.rep)
			if ok || line != "" {
				t.Errorf("a report (state=%q) that attempted no fold printed %q (ok=%v); that shape "+
					"never tried, so a line there is noise", tc.rep.State, line, ok)
			}
			var out bytes.Buffer
			if printed := printCompressionTrace(&out, tc.rep); printed {
				t.Error("printCompressionTrace claimed a print for a report that attempted no fold")
			}
			if out.Len() != 0 {
				t.Errorf("the skipped shape wrote %q, want zero bytes", out.String())
			}
		})
	}
}
