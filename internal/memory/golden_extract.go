package memory

// Golden extraction contract (ticket 29 AC#2: "Golden extraction set: 10
// scripted conversations -> expected profile slots/values"). This file is the
// CONTRACT; the fixture is testdata/golden_extract.json (embedded below); the
// comparison test is golden_extract_test.go.
//
// Why this is not decoration: the L1 extraction pipeline is NOT wired yet
// (29-v1: the host never assigns Hooks.WarmWindow, and no C5 LLM call feeds
// profile upserts). Without a pinned target, the future wiring leg would pick
// expected slots AFTER seeing what its code produces, which proves nothing.
// The golden set pins the target first: for each of the 10 scripted
// conversations, exactly which slots must land, with exactly which content.
// The wiring leg then gets a real target to shoot at instead of writing the
// test to match its own output.
//
// How the future leg swaps the stub for the real extractor:
//
//  1. Implement Extractor where the C5 chain lives: one small LLM call per
//     task over the conversation turns (via llm.LlmProvider, golden SSE bytes
//     recorded under the C5 golden harness, cf. internal/llm/golden and the
//     internal/agent/testdata/golden/*.sse precedent), run in the Warm window
//     after the response, never on the response path (SPEC-05 section 6).
//     The call returns the ExtractionOutcome for the WHOLE task in one go,
//     which is what batch-commit means: one prompt-cache invalidation per
//     task, not per fact.
//  2. In golden_extract_test.go replace the single construction point
//     newGoldenExtractor() (currently the rule-based stub) with the real
//     chain. Nothing else changes: same fixture, same CompareExtraction,
//     same assertions against the real profile table.
//  3. Run: go test ./internal/memory -run 'TestTicket29_AC2' -count=1.
//     Green = the real chain reproduces the pinned slots/values on all 10
//     scripts; red = the chain disagrees with the pinned contract. If the
//     CONTRACT itself is wrong, amend the golden JSON through a new ticket
//     (it is a reviewed fixture, not a scratch file).
//
// Value formats are part of the contract (the real chain must emit exactly
// these normalized values):
//   - pref.language:    "zh" | "en"  (the language the ASSISTANT replies in,
//     not the language the user typed in -- case-07 types mixed but pins zh)
//   - pref.tone:        "casual" | "concise" | "direct"
//   - habit.work_hours: "weekday HH:MM-HH:MM" (24h, zero-padded)
//   - habit.editor:     "vscode" | "goland"
//   - fact.family:      "mother birthday MM-DD" | "father birthday MM-DD"
//
// Slot names must stay inside the SPEC-02 section 3 enum namespaces
// (pref.* | habit.* | fact.*). habit.editor is a golden-defined member of the
// habit.* namespace (the spec names four members plus an ellipsis; the golden
// set is where further members get pinned).
//
// Case-level Risk uses the codebase's real levels (models.go RiskL0/RiskL1)
// as a sensitivity tier for the outcome: L0 = routine preference/habit slot,
// L1 = personal-fact slot (fact.*) or the undo-intent path, where the
// one-interaction undo affordance (SPEC-05 section 6) is the behavior under
// test. It is fixture metadata for the future pipeline's emphasis, not a new
// runtime semantic.
//
// The fixture lives in internal/memory/testdata/ (package-scoped testdata;
// the repo-root testdata tree is not touched) and is embedded via go:embed so
// tests do not depend on the working directory and out-of-repo -overlay
// mutations of the JSON reach the test binary (that is how the reverse
// criterion mutates one expected value and watches the suite go red).

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// goldenExtractJSON is the embedded golden fixture (testdata/golden_extract.json).
//
//go:embed testdata/golden_extract.json
var goldenExtractJSON []byte

// SlotUpsert is one slot write: the REAL profile row shape minus generated
// columns (Profile.Slot + Profile.Content; Source is fixed to
// ProfileExtracted for extraction outcomes, UpdatedAt is wall clock).
type SlotUpsert struct {
	Slot    string `json:"slot"`
	Content string `json:"content"`
}

// ExtractionOutcome is what one extraction run produces for one conversation:
// the slots to write and the slots the user asked to stop remembering ("别记
// 了" undo intent). The golden fixture's expect block has this same type --
// the contract is literally "the real chain must produce exactly this
// outcome". Semantics are NET per task: a forget cancels any upsert of the
// same slot in the same outcome, and the whole outcome is committed as one
// batch (one prompt-cache invalidation per task, per AC#2).
type ExtractionOutcome struct {
	Upserts []SlotUpsert `json:"upserts"`
	Forgets []string     `json:"forgets"`
}

// GoldenTurn is one scripted utterance. role is "user" or "assistant"; only
// user turns are extraction candidates (what the assistant says about itself
// is not a fact about the user).
type GoldenTurn struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

// GoldenConversation is the minimal transcript shape the extraction seam
// consumes (built from a GoldenCase in the comparison test; not a JSON type).
// The C5 chain implements Extractor by rendering these turns into its
// extraction prompt and parsing the structured reply into an
// ExtractionOutcome.
type GoldenConversation struct {
	ID    string       // fixture case id, carried for logs
	Title string       // fixture case title, carried for logs
	Turns []GoldenTurn // only user turns are extraction candidates
}

// ExtractionRisk tiers (see the file comment). Only these two are valid in
// the golden set: L2 is the tool-gating level (D34), not an extraction tier.
const (
	ExtractionRiskNormal    = RiskL0
	ExtractionRiskSensitive = RiskL1
)

// GoldenCase is one scripted conversation plus its expected outcome.
type GoldenCase struct {
	ID     string            `json:"id"`
	Title  string            `json:"title"`
	Risk   string            `json:"risk"`
	Covers []string          `json:"covers"`
	Seed   *SlotUpsert       `json:"seed,omitempty"`
	Turns  []GoldenTurn      `json:"turns"`
	Expect ExtractionOutcome `json:"expect"`
}

// GoldenExtractSet is the whole fixture: version, label, and the 10 cases.
type GoldenExtractSet struct {
	Version int          `json:"version"`
	Set     string       `json:"set"`
	Cases   []GoldenCase `json:"cases"`
}

// goldenSetCaseCount is the AC#2 pinned case count (10 conversations).
const goldenSetCaseCount = 10

// goldenSetVersion is the fixture schema version the loader accepts.
const goldenSetVersion = 1

// Extractor is the L1 profile extraction seam (ticket 29 AC#2). TODAY the
// only implementation is the deterministic rule-based stub in
// golden_extract_test.go, which exists to prove the golden set is
// self-consistent. The REAL implementation is the C5 LLM chain (small call in
// the Warm window, batch-commit, see the swap instructions above).
type Extractor interface {
	Extract(ctx context.Context, conv GoldenConversation) (ExtractionOutcome, error)
}

// LoadGoldenExtractSet decodes and validates the embedded golden fixture.
// Validation is the contract's teeth: a fixture that drifts (wrong count,
// unknown slot namespace, empty content, bogus risk) fails HERE, so the
// comparison tests below always argue against a well-formed target.
func LoadGoldenExtractSet() (GoldenExtractSet, error) {
	var set GoldenExtractSet
	if err := json.Unmarshal(goldenExtractJSON, &set); err != nil {
		return GoldenExtractSet{}, fmt.Errorf("memory: golden extract fixture: %w", err)
	}
	if set.Version != goldenSetVersion {
		return GoldenExtractSet{}, fmt.Errorf("memory: golden extract fixture version %d, want %d",
			set.Version, goldenSetVersion)
	}
	if len(set.Cases) != goldenSetCaseCount {
		return GoldenExtractSet{}, fmt.Errorf("memory: golden extract set holds %d cases, want %d (AC#2: 10 conversations)",
			len(set.Cases), goldenSetCaseCount)
	}
	seen := make(map[string]bool, len(set.Cases))
	for i := range set.Cases {
		c := &set.Cases[i]
		if c.ID == "" || seen[c.ID] {
			return GoldenExtractSet{}, fmt.Errorf("memory: golden case %d has empty or duplicate id %q", i, c.ID)
		}
		seen[c.ID] = true
		if len(c.Turns) == 0 {
			return GoldenExtractSet{}, fmt.Errorf("memory: golden case %s has no turns", c.ID)
		}
		for _, t := range c.Turns {
			if t.Role != "user" && t.Role != "assistant" {
				return GoldenExtractSet{}, fmt.Errorf("memory: golden case %s has turn role %q (want user|assistant)", c.ID, t.Role)
			}
		}
		switch c.Risk {
		case ExtractionRiskNormal, ExtractionRiskSensitive:
			// ok
		default:
			return GoldenExtractSet{}, fmt.Errorf("memory: golden case %s risk %q (want L0|L1; L2 is the tool-gating level)", c.ID, c.Risk)
		}
		if len(c.Covers) == 0 {
			return GoldenExtractSet{}, fmt.Errorf("memory: golden case %s covers nothing", c.ID)
		}
		if c.Seed != nil {
			if err := validateGoldenSlot(c.Seed.Slot, c.Seed.Content); err != nil {
				return GoldenExtractSet{}, fmt.Errorf("memory: golden case %s seed: %w", c.ID, err)
			}
		}
		for _, u := range c.Expect.Upserts {
			if err := validateGoldenSlot(u.Slot, u.Content); err != nil {
				return GoldenExtractSet{}, fmt.Errorf("memory: golden case %s expect: %w", c.ID, err)
			}
		}
		for _, f := range c.Expect.Forgets {
			if err := validateGoldenSlot(f, "x"); err != nil {
				return GoldenExtractSet{}, fmt.Errorf("memory: golden case %s forget %q: %w", c.ID, f, err)
			}
		}
	}
	return set, nil
}

// validateGoldenSlot checks one slot name against the SPEC-02 section 3 enum
// namespaces and requires non-empty content.
func validateGoldenSlot(slot, content string) error {
	switch {
	case strings.HasPrefix(slot, "pref."),
		strings.HasPrefix(slot, "habit."),
		strings.HasPrefix(slot, "fact."):
	default:
		return fmt.Errorf("slot %q outside the SPEC-02 enum namespaces (pref.*|habit.*|fact.*)", slot)
	}
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("slot %q has empty content", slot)
	}
	return nil
}

// CompareExtraction is the golden comparator: got (what an extractor
// produced) against want (the pinned expectation). Both sides are compared as
// sorted slot sets, duplicates are rejected, and every mismatch names the
// slot and both values. The future C5 chain is judged by exactly this
// function -- swapping the stub for the real chain changes nothing here.
func CompareExtraction(got, want ExtractionOutcome) error {
	gotUpserts := append([]SlotUpsert(nil), got.Upserts...)
	sort.Slice(gotUpserts, func(i, j int) bool { return gotUpserts[i].Slot < gotUpserts[j].Slot })
	for i := 1; i < len(gotUpserts); i++ {
		if gotUpserts[i].Slot == gotUpserts[i-1].Slot {
			return fmt.Errorf("memory: golden extraction: duplicate upsert slot %q", gotUpserts[i].Slot)
		}
	}
	wantUpserts := append([]SlotUpsert(nil), want.Upserts...)
	sort.Slice(wantUpserts, func(i, j int) bool { return wantUpserts[i].Slot < wantUpserts[j].Slot })

	if len(gotUpserts) != len(wantUpserts) {
		return fmt.Errorf("memory: golden extraction: upsert count %d, want %d (got slots %v, want slots %v)",
			len(gotUpserts), len(wantUpserts), slotNames(gotUpserts), slotNames(wantUpserts))
	}
	for i := range wantUpserts {
		if gotUpserts[i].Slot != wantUpserts[i].Slot {
			return fmt.Errorf("memory: golden extraction: upsert slot %q, want %q", gotUpserts[i].Slot, wantUpserts[i].Slot)
		}
		if gotUpserts[i].Content != wantUpserts[i].Content {
			return fmt.Errorf("memory: golden extraction: slot %q content %q, want %q",
				wantUpserts[i].Slot, gotUpserts[i].Content, wantUpserts[i].Content)
		}
	}

	gotForgets := append([]string(nil), got.Forgets...)
	sort.Strings(gotForgets)
	wantForgets := append([]string(nil), want.Forgets...)
	sort.Strings(wantForgets)
	if len(gotForgets) != len(wantForgets) {
		return fmt.Errorf("memory: golden extraction: forget count %d, want %d (got %v, want %v)",
			len(gotForgets), len(wantForgets), gotForgets, wantForgets)
	}
	for i := range wantForgets {
		if gotForgets[i] != wantForgets[i] {
			return fmt.Errorf("memory: golden extraction: forget slot %q, want %q", gotForgets[i], wantForgets[i])
		}
	}
	return nil
}

// slotNames renders an upsert list as bare slot names for error messages.
func slotNames(us []SlotUpsert) []string {
	out := make([]string, len(us))
	for i, u := range us {
		out[i] = u.Slot
	}
	return out
}
