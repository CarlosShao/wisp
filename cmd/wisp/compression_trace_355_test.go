package main

// Ticket 355 AC#2 and AC#3: the compression trace must be answerable from the
// production side, not only by a human opening <dataDir>\logs\wisp-*.jsonl.
//
// What is pinned here, and why it is pinned THIS way:
//
//   - The three flavours (tokens_before, tokens_after, history_changed) come
//     from ONE report. history_changed is the compressor's own reading
//     (internal/agent/compress.go: rep.TokensAfter != rep.TokensBefore ||
//     the message count moved), and it can contradict Ran. Recomputing it in
//     this package from other counters would let the closing line and the
//     on-disk record disagree, so the reading is carried on the report and
//     this file feeds a REAL one through the REAL formatter.
//
//   - The report handed to compressionLine is produced by agent.NewCompressor
//     over a history that genuinely clears the threshold - not by a struct
//     literal. That is what makes the reverse judgement bite: delete the
//     `rep.HistoryChanged = ...` assignment in compress.go and this test goes
//     red, because the pass really did move the history while the report would
//     then claim it did not.
//
//   - The "most common shape" is answered out loud rather than papered over,
//     and since ticket 358 the answer differs per state: compressionLine is
//     silent for the shapes where nothing was attempted (a run that never
//     cleared the threshold, and the never-reached-compressor report), and the
//     second test below is the fang against a line that prints zeros for a
//     measurement nobody took. An attempt that folded nothing is NOT silent any
//     more - it prints its own line - see that test's own header for the ruling
//     that re-cut it. The budget contract that decides when a fold is possible
//     is Q-312's and is deliberately untouched.
//
// Every threshold is read back out of agent.BudgetsFor; none is a written
// number (ticket 139 AC#1(2) forbids citing the reference value as a trigger).

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm"
)

// buildTraceHistory makes `rounds` conversational rounds that each clear a
// known token weight, so the compressor's own budget decides whether a fold
// happens instead of this fixture guessing.
func buildTraceHistory(rounds, padBytes int) []llm.Message {
	pad := strings.Repeat("x", padBytes)
	hist := make([]llm.Message, 0, rounds*3)
	for i := 0; i < rounds; i++ {
		id := fmt.Sprintf("c%02d", i)
		hist = append(hist,
			llm.Message{
				Role:    llm.RoleUser,
				Content: []llm.Content{llm.TextPart{Text: fmt.Sprintf("问题 %d %s", i, pad)}},
			},
			llm.Message{
				Role: llm.RoleAssistant,
				Content: []llm.Content{llm.ToolUsePart{
					ID: id, Name: "echo", Input: json.RawMessage(`{"text":"hi"}`),
				}},
			},
			llm.Message{
				Role: llm.RoleTool,
				Content: []llm.Content{llm.ToolResultPart{
					ID:      id,
					Content: []llm.Content{llm.TextPart{Text: "结果 " + pad}},
				}},
			},
		)
	}
	return hist
}

// realCompressionReport drives the shipped compressor over a history that
// really is over its threshold and returns the report it produced. It fails the
// test rather than returning a zero report if no fold happened: a fixture that
// silently stopped clearing the threshold would turn every assertion below
// into a tautology.
func realCompressionReport(t *testing.T) agent.CompressionReport {
	t.Helper()
	b := agent.BudgetsFor(4096)
	c := agent.NewCompressor(b, nil)
	hist := buildTraceHistory(b.KeepRawRounds+3, 400)
	if !c.Need(hist) {
		t.Fatalf("setup: history of %d tokens does not clear the threshold %d, so this pass folds nothing",
			c.TotalTokens(hist), b.HistoryCompressTokens)
	}
	_, rep, err := c.Compress(context.Background(), hist)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if !rep.Ran {
		t.Fatal("setup: the compressor reported no fold over a history that clears the threshold and has more raw rounds than KeepRawRounds")
	}
	return rep
}

// TestCompressionTraceThreeFlavoursReachTheClosingLine is AC#2's positive half:
// a pass that really folded the history yields a line carrying all three
// values, and the values on the line are the report's values.
func TestCompressionTraceThreeFlavoursReachTheClosingLine(t *testing.T) {
	rep := realCompressionReport(t)

	// The report must be able to answer the third flavour at all. This is the
	// assertion the reverse judgement turns red: without the assignment in
	// compress.go the field stays false while the history demonstrably moved.
	if !rep.HistoryChanged {
		t.Errorf("HistoryChanged = false for a pass that replaced the history: "+
			"%d -> %d tokens, %d messages folded into the summary",
			rep.TokensBefore, rep.TokensAfter, rep.CompressedMsgs)
	}
	if rep.TokensAfter >= rep.TokensBefore {
		t.Errorf("tokens did not drop: %d -> %d, so this fixture no longer proves a real fold",
			rep.TokensBefore, rep.TokensAfter)
	}

	line, ok := compressionLine(rep)
	if !ok {
		t.Fatal("compressionLine answered nothing for a pass that ran")
	}
	// Each flavour by name, read off the line the operator actually sees.
	for _, want := range []string{
		fmt.Sprintf("压前 %d tokens", rep.TokensBefore),
		fmt.Sprintf("压后 %d tokens", rep.TokensAfter),
		"history_changed=true",
	} {
		if !strings.Contains(line, want) {
			t.Errorf("the closing line does not carry %q: %s", want, line)
		}
	}
	// Privacy: counts only. The fixture's prose sentinel must not travel.
	if strings.Contains(line, "问题") || strings.Contains(line, strings.Repeat("x", 40)) {
		t.Errorf("the closing line leaked history content: %s", line)
	}
}

// TestCompressionTraceSilentWhenNothingCompressed is the fang against a line
// that prints when nothing was attempted. It pins BOTH halves of the shape
// ticket 358 re-cut:
//
//	(a) a run whose history never cleared the trigger, i.e. the report's state is
//	    skipped, still answers NOTHING - a line there would be noise, and a line
//	    of zeros would claim a measurement nobody took;
//	(b) a pass that DID clear the threshold and folded nothing (state
//	    attempted) must say so out loud, with the four readings that make the
//	    claim checkable. It is NOT silent any more.
//
// WHY THIS TEST CHANGED, IN THE RULING'S OWN WORDS (decision of 2026-10-21 17:3x
// +08, archived verbatim at
// .scratch/wisp/probes/orch/20261021/txt/358-ruling-nail-355.txt, taken up by
// ticket 358 AC#2). Until this run the second half asserted silence, and the
// argument written at the time was:
//
//	"ok=false covers both 'this run never compressed' and 'the run is not the
//	 one that compressed'... That is stated rather than papered over - a line
//	 printing zeros would claim a measurement nobody took."
//
// THAT ARGUMENT WAS TRUE THEN AND ITS PREMISE IS GONE. It rested on the two
// shapes being indistinguishable on the report - both Ran=false, no field able
// to tell them apart - so any line emitted for one was also a claim about the
// other. Ticket 358 gave CompressionReport a CompressionState
// (skipped/attempted/compacted, internal/agent/compress.go), which makes "never
// triggered" and "triggered but could not fold" different values of one field.
// Keeping the old silence would therefore not be caution; it would be an
// outdated nail holding back a reading that now exists, which is the very
// "spec and instrument drift apart" failure the ruling names.
//
// The cut is stronger than what it replaces, not looser: it still forbids a
// line for the skipped shape AND now requires the attempted shape to spell out
// that nothing was folded (rather than only forbidding a zeros line). No
// assertion was deleted to turn anything green.
func TestCompressionTraceSilentWhenNothingCompressed(t *testing.T) {
	// (a) The ordinary shape: a run whose history never cleared the threshold, so
	// the loop never called the compressor and Result.Compression is zero. Its
	// State is the zero string, which is "no pass ran", not "we looked and said
	// no" - and it prints nothing.
	if line, ok := compressionLine(agent.CompressionReport{}); ok || line != "" {
		t.Errorf("a run that never reached the compressor printed %q (ok=%v); a line of zeros would "+
			"claim a measurement nobody took", line, ok)
	}

	// (b) The shape ticket 153 named and ticket 358 made readable: the threshold
	// is cleared but fewer than KeepRawRounds+1 raw rounds exist, so the fold loop
	// breaks immediately. Ran is false - and the line must now SPEAK.
	// The pad is sized from the budget itself so this stays over the threshold
	// even after scaling: one round of 4000 bytes is ~1000 tokens against the
	// 384 this window derives, with no number written by hand.
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
	if rep.Ran {
		t.Fatal("setup: a single raw round folded, so this is not the shape ticket 153 described")
	}
	if rep.State != agent.CompressionAttempted {
		t.Fatalf("setup: State = %q, want %q for a pass that cleared the threshold and folded nothing",
			rep.State, agent.CompressionAttempted)
	}

	line, ok := compressionLine(rep)
	if !ok || line == "" {
		t.Fatalf("a pass that cleared the threshold but folded nothing answered %q (ok=%v): the state "+
			"is %q, so this run DID try and the reading must say it folded nothing - silence here is "+
			"the hole ticket 358 exists for", line, ok, rep.State)
	}
	// It must say, in words, that the attempt folded nothing - not merely print
	// zeros, and not claim the trigger never fired.
	if !strings.Contains(line, "超了阈值但一枚都没折") {
		t.Errorf("the attempted line does not say the attempt folded nothing: %s", line)
	}
	for _, never := range []string{"没触发", "未触发", "跳过"} {
		if strings.Contains(line, never) {
			t.Errorf("the attempted line claims %q, which is the skipped shape's answer: %s", never, line)
		}
	}
	// The four readings, each by name and each carrying the report's own value.
	for _, want := range []string{
		fmt.Sprintf("tokens_before=%d", rep.TokensBefore),
		fmt.Sprintf("tokens_after=%d", rep.TokensAfter),
		fmt.Sprintf("trigger=%d tokens", rep.TriggerTokens),
		fmt.Sprintf("source=%s", string(rep.TriggerSource)),
	} {
		if !strings.Contains(line, want) {
			t.Errorf("the attempted line does not carry %q: %s", want, line)
		}
	}
	// The reading is taken from the report, so an attempted pass whose numbers
	// really are equal must still be readable as an attempt (this is exactly what
	// zeros alone could not say).
	if rep.TokensBefore != rep.TokensAfter {
		t.Errorf("setup: this fixture was meant to fold nothing, yet %d => %d",
			rep.TokensBefore, rep.TokensAfter)
	}

	// (c) The third half cannot go untested either: the skipped STATE must be
	// silent even when it is a filled report rather than the zero value, so
	// nobody can make the silence above pass by only special-casing {}.
	skipped := agent.CompressionReport{State: agent.CompressionSkipped}
	if line, ok := compressionLine(skipped); ok || line != "" {
		t.Errorf("a report that names the skipped state printed %q (ok=%v); that shape never "+
			"attempted a fold, so a line there is noise", line, ok)
	}
}

// TestCompressionTraceReachesStdout is AC#2's exit half: the line is emitted on
// the run's own stdout, after the closing summary, by the shipped helper the
// composition root calls. The bytes go through printCompressionTrace, which is
// what execute() invokes - not through a formatter this test re-implements.
func TestCompressionTraceReachesStdout(t *testing.T) {
	rep := realCompressionReport(t)
	var out bytes.Buffer

	if printed := printCompressionTrace(&out, rep); !printed {
		t.Fatal("printCompressionTrace reported that it printed nothing for a pass that ran")
	}
	got := out.String()
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("the compression line is not newline-terminated, so it would run into the next one: %q", got)
	}
	if n := strings.Count(got, "\n"); n != 1 {
		t.Errorf("the compression reading should be exactly one line, got %d: %q", n, got)
	}
	for _, want := range []string{"压缩痕", "压前", "压后", "history_changed"} {
		if !strings.Contains(got, want) {
			t.Errorf("stdout does not carry %q: %s", want, got)
		}
	}

	// The silent shape writes zero bytes - not a blank line, which would be a
	// paragraph the operator has to interpret.
	var quiet bytes.Buffer
	if printed := printCompressionTrace(&quiet, agent.CompressionReport{}); printed {
		t.Error("printCompressionTrace claimed a print for a run that never compressed")
	}
	if quiet.Len() != 0 {
		t.Errorf("a run that never compressed wrote %q, want zero bytes", quiet.String())
	}
}
