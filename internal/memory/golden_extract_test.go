package memory

// Ticket 29 AC#2 golden comparison tests (see golden_extract.go for the
// contract). The instruments:
//
//   - TestTicket29_AC2_GoldenSetContract: the fixture itself is well-formed
//     (10 cases, unique ids, the coverage axes the ticket demands).
//   - TestTicket29_AC2_GoldenExtractionMatches: per case, the REAL extraction
//     chain runs over the scripted conversation (one C5 call replaying the
//     case's pinned golden-SSE recording, see golden_sse_replay_test.go), the
//     outcome must equal the golden expectation, and applying that outcome
//     through the REAL profile DAO (UpsertProfile / ProfileBySlot /
//     DeleteProfile / ListProfiles) must leave exactly the expected slots with
//     source=extracted. One recorded response per case is also the AC's
//     batch-commit assertion: one C5 call, hence at most one prompt-cache
//     invalidation per task.
//   - TestTicket29_AC2_GoldenBytesThroughPipeline: the same recordings driven
//     through the PRODUCTION exec (L1Pipeline.RunTask: one C5 call, one store
//     transaction) must land the same profile projection, so the golden bytes
//     are tied to the path the host actually runs and not only to a test-only
//     apply loop.
//   - TestTicket29_AC2_ComparatorNegativeControls: a corrupted outcome must
//     be rejected by CompareExtraction, naming the slot. This is the
//     always-on reverse criterion; the file-level one (mutate one expected
//     value in the fixture and watch the suite go red) runs out-of-repo via
//     go test -overlay, because the fixture is embedded.

import (
	"context"
	"strings"
	"testing"
)

// newGoldenExtractor is the single construction point for the extractor under
// test: the REAL chain, an *L1Extractor over the real openai-chat adapter fed
// by each case's pinned golden-SSE recording (golden_sse_replay_test.go). Every
// TestTicket29_AC2 assertion therefore judges the production chain against the
// same fixture, unchanged. History: this returned a rule-based stub until the
// AC#2 leg. The stub is deleted rather than kept alongside, because a double
// that reads only the transcript keeps passing when the real chain's request
// or reply-parsing shape breaks -- which is exactly the failure the golden set
// exists to catch.
func newGoldenExtractor(t *testing.T, set GoldenExtractSet) (Extractor, *goldenSSERouter) {
	t.Helper()
	return newGoldenSSEChain(t, set)
}

// goldenCoverageAxes are the behaviors the ticket demands the 10 scripts
// cover; the fixture's covers fields must touch every one of them.
var goldenCoverageAxes = []string{
	"pref-slots",       // 偏好 -> slots
	"habit-slots",      // 习惯 -> slots
	"multi-turn",       // 多轮增量（同槽位后轮改写，末值生效）
	"noise",            // 噪声句不许误报
	"forget-intent",    // "别记了" 撤销意图
	"mixed-language",   // 中英混排，回复语言锚定
	"english-pref",     // 纯英文偏好（枚举值 en）
	"family-fact",      // fact.* 敏感事实
	"batch-multi-slot", // 一次任务多槽位，整批一次提交
}

func TestTicket29_AC2_GoldenSetContract(t *testing.T) {
	set, err := LoadGoldenExtractSet()
	if err != nil {
		t.Fatalf("LoadGoldenExtractSet: %v", err)
	}
	if len(set.Cases) != 10 {
		t.Fatalf("golden set holds %d cases, want 10 (AC#2: 10 scripted conversations)", len(set.Cases))
	}
	covered := map[string]bool{}
	for _, c := range set.Cases {
		if !strings.HasPrefix(c.ID, "case-") {
			t.Errorf("case id %q does not follow the case-NN-slug shape", c.ID)
		}
		for _, axis := range c.Covers {
			covered[axis] = true
		}
	}
	for _, axis := range goldenCoverageAxes {
		if !covered[axis] {
			t.Errorf("coverage axis %q is not covered by any case; the golden set no longer covers the ticket", axis)
		}
	}
}

func TestTicket29_AC2_GoldenExtractionMatches(t *testing.T) {
	set, err := LoadGoldenExtractSet()
	if err != nil {
		t.Fatalf("LoadGoldenExtractSet: %v", err)
	}
	ex, replay := newGoldenExtractor(t, set)
	for i := range set.Cases {
		c := set.Cases[i]
		t.Run(c.ID, func(t *testing.T) {
			ctx := context.Background()
			s := openTestStore(t)
			if c.Seed != nil {
				// A seed models a slot stored by an EARLIER task, so a
				// forget intent has something real to remove.
				if _, err := s.UpsertProfile(ctx, c.Seed.Slot, c.Seed.Content, ProfileExtracted); err != nil {
					t.Fatalf("seed UpsertProfile: %v", err)
				}
			}
			got, err := ex.Extract(ctx, GoldenConversation{ID: c.ID, Title: c.Title, Turns: c.Turns})
			if err != nil {
				t.Fatalf("Extract: %v", err)
			}
			// Batch-commit: the real chain answered this task from exactly one
			// C5 call (one recorded response), so it costs at most one
			// prompt-cache invalidation - not one per fact.
			if n := replay.callsFor(c.ID); n != 1 {
				t.Fatalf("C5 calls for this task = %d, want exactly 1 (one prompt-cache invalidation per task)", n)
			}
			if err := CompareExtraction(got, c.Expect); err != nil {
				t.Fatalf("golden mismatch: %v", err)
			}
			// Apply through the REAL profile DAO -- this is the batch the
			// Warm window would commit: upserts first, then forgets.
			for _, u := range got.Upserts {
				if _, err := s.UpsertProfile(ctx, u.Slot, u.Content, ProfileExtracted); err != nil {
					t.Fatalf("UpsertProfile %s: %v", u.Slot, err)
				}
			}
			for _, slot := range got.Forgets {
				p, err := s.ProfileBySlot(ctx, slot)
				if err != nil {
					t.Fatalf("forget target %q must exist after applying upserts (seed or upsert): %v", slot, err)
				}
				if err := s.DeleteProfile(ctx, p.ID); err != nil {
					t.Fatalf("DeleteProfile %q: %v", slot, err)
				}
			}
			// Read back through the real table: the final projection must be
			// exactly the golden upserts minus the golden forgets, every row
			// provenance-extracted (the real Profile struct, not a shadow).
			goldenProfileProjection(t, s, c)
		})
	}
	// One call per scripted conversation, and no conversation skipped: this is
	// the seam-level reading of "one prompt-cache invalidation per task".
	if got, want := replay.totalCalls(), len(set.Cases); got != want {
		t.Fatalf("replay served %d C5 calls for %d cases, want exactly one per case", got, want)
	}
	if !replay.exhausted() {
		t.Fatal("not every recording was consumed: a case never reached the C5 seam")
	}
	if n := replay.prefixVariants(); n != 1 {
		t.Fatalf("the extraction system prompt took %d forms across tasks, want 1 (cache-stable prefix)", n)
	}
}

// TestTicket29_AC2_GoldenBytesThroughPipeline drives the pinned recordings
// through L1Pipeline.RunTask, the production exec the Warm window calls: ONE
// C5 call and ONE store transaction per task. The projection read back must be
// the fixture's, exactly as in the comparison test - the difference is that
// nothing here is applied by hand.
func TestTicket29_AC2_GoldenBytesThroughPipeline(t *testing.T) {
	set, err := LoadGoldenExtractSet()
	if err != nil {
		t.Fatalf("LoadGoldenExtractSet: %v", err)
	}
	ex, replay := newGoldenExtractor(t, set)
	ctx := context.Background()
	for i := range set.Cases {
		c := set.Cases[i]
		t.Run(c.ID, func(t *testing.T) {
			s := openTestStore(t)
			if c.Seed != nil {
				if _, err := s.UpsertProfile(ctx, c.Seed.Slot, c.Seed.Content, ProfileExtracted); err != nil {
					t.Fatalf("seed UpsertProfile: %v", err)
				}
			}
			pipe, err := NewL1Pipeline(s, ex)
			if err != nil {
				t.Fatalf("NewL1Pipeline: %v", err)
			}
			got, err := pipe.RunTask(ctx, GoldenConversation{ID: c.ID, Title: c.Title, Turns: c.Turns})
			if err != nil {
				t.Fatalf("RunTask: %v", err)
			}
			if err := CompareExtraction(got, c.Expect); err != nil {
				t.Fatalf("golden mismatch: %v", err)
			}
			if n := replay.callsFor(c.ID); n != 1 {
				t.Fatalf("C5 calls for this task = %d, want exactly 1 (one prompt-cache invalidation per task)", n)
			}
			goldenProfileProjection(t, s, c)
		})
	}
	if got, want := replay.totalCalls(), len(set.Cases); got != want {
		t.Fatalf("replay served %d C5 calls for %d cases, want exactly one per case", got, want)
	}
	if !replay.exhausted() {
		t.Fatal("not every recording was consumed: a case never reached the C5 seam")
	}
}

// goldenProfileProjection requires the real profile table to hold exactly the
// fixture's upserts minus its forgets, every row provenance-extracted. Both
// AC#2 apply paths (DAO loop and production pipeline) read the target back
// through this one instrument.
func goldenProfileProjection(t *testing.T, s *Store, c GoldenCase) {
	t.Helper()
	ctx := context.Background()
	rows, err := s.ListProfiles(ctx)
	if err != nil {
		t.Fatalf("ListProfiles: %v", err)
	}
	want := map[string]string{}
	for _, u := range c.Expect.Upserts {
		want[u.Slot] = u.Content
	}
	for _, f := range c.Expect.Forgets {
		delete(want, f)
	}
	final := map[string]string{}
	for _, r := range rows {
		if r.Source != string(ProfileExtracted) {
			t.Errorf("slot %s stored with source %q, want extracted", r.Slot, r.Source)
		}
		final[r.Slot] = r.Content
	}
	if len(final) != len(want) {
		t.Fatalf("final profile holds %d rows, want %d (got %v, want %v)",
			len(final), len(want), final, want)
	}
	for slot, content := range want {
		if final[slot] != content {
			t.Errorf("slot %s content %q, want %q", slot, final[slot], content)
		}
	}
}

// TestTicket29_AC2_ComparatorNegativeControls is the in-test reverse
// criterion: every kind of corruption must be REJECTED, naming the slot.
func TestTicket29_AC2_ComparatorNegativeControls(t *testing.T) {
	set, err := LoadGoldenExtractSet()
	if err != nil {
		t.Fatalf("LoadGoldenExtractSet: %v", err)
	}

	// V1: an expected slot's content is wrong.
	base := set.Cases[0].Expect
	bad := base
	bad.Upserts = append([]SlotUpsert(nil), base.Upserts...)
	for i := range bad.Upserts {
		if bad.Upserts[i].Slot == "pref.tone" {
			bad.Upserts[i].Content = "formal"
		}
	}
	err = CompareExtraction(base, bad)
	if err == nil || !strings.Contains(err.Error(), "pref.tone") {
		t.Fatalf("corrupted tone content must fail naming pref.tone, got %v", err)
	}

	// V2: an expected upsert is missing entirely.
	dropped := base
	dropped.Upserts = base.Upserts[:1]
	if err := CompareExtraction(dropped, base); err == nil || !strings.Contains(err.Error(), "upsert count") {
		t.Fatalf("dropped upsert must fail on the count, got %v", err)
	}

	// V3: an unrequested extra upsert appears.
	extra := base
	extra.Upserts = append(append([]SlotUpsert(nil), base.Upserts...), SlotUpsert{Slot: "habit.editor", Content: "vscode"})
	if err := CompareExtraction(extra, base); err == nil || !strings.Contains(err.Error(), "upsert count") {
		t.Fatalf("extra upsert must fail on the count, got %v", err)
	}

	// V4: the undo intent is lost (a forget disappears from the outcome).
	var forgetCase *GoldenCase
	for i := range set.Cases {
		if set.Cases[i].ID == "case-06-forget-intent" {
			forgetCase = &set.Cases[i]
		}
	}
	if forgetCase == nil {
		t.Fatal("case-06-forget-intent missing from the golden set")
	}
	noForget := forgetCase.Expect
	noForget.Forgets = nil
	err = CompareExtraction(noForget, forgetCase.Expect)
	if err == nil || !strings.Contains(err.Error(), "forget") {
		t.Fatalf("lost forget intent must fail mentioning forgets, got %v", err)
	}
}
