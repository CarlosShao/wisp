package memory

// Ticket 29 AC#2 golden comparison tests (see golden_extract.go for the
// contract). Three instruments:
//
//   - TestTicket29_AC2_GoldenSetContract: the fixture itself is well-formed
//     (10 cases, unique ids, the coverage axes the ticket demands).
//   - TestTicket29_AC2_GoldenExtractionMatches: per case, the extractor runs
//     over the scripted conversation, the outcome must equal the golden
//     expectation, and applying that outcome through the REAL profile DAO
//     (UpsertProfile / ProfileBySlot / DeleteProfile / ListProfiles) must
//     leave exactly the expected slots with source=extracted.
//   - TestTicket29_AC2_ComparatorNegativeControls: a corrupted outcome must
//     be rejected by CompareExtraction, naming the slot. This is the
//     always-on reverse criterion; the file-level one (mutate one expected
//     value in the fixture and watch the suite go red) runs out-of-repo via
//     go test -overlay, because the fixture is embedded.

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// newGoldenExtractor is the single construction point for the extractor under
// test. SWAP POINT for the future C5 chain: return the real Extractor here
// and every TestTicket29_AC2 test judges the real chain against the same
// golden fixture, unchanged (full procedure in golden_extract.go).
func newGoldenExtractor() Extractor {
	return goldenStubExtractor{}
}

// goldenStubExtractor is a deterministic, rule-based STUB Extractor. It reads
// ONLY the conversation text (never the case id), so the golden expectations
// it passes against are provably derivable from the scripts. It is NOT the
// production extractor: the real implementation is the C5 LLM chain (small
// call in the Warm window, one batched outcome per task), delivered by the
// future wiring leg -- swap it in at newGoldenExtractor().
//
// Rules (each fires on user turns only; later turns override earlier values
// for the same slot; a forget turn cancels any pending upsert of that slot):
//   - forget intent ("别记了" family): the slot is named by a keyword in the
//     same turn (说话/语气/风格 -> pref.tone, 语言/中文/英文 -> pref.language,
//     上班/工作时间/作息 -> habit.work_hours, 编辑器/ide -> habit.editor,
//     我妈/我爸/家人 -> fact.family); no keyword, no forget (noise).
//   - reply language: "都用中文|用中文回|中文回复" -> zh, "reply in english"
//     -> en (the language the assistant replies in, not the user's).
//   - tone: 简短/简洁/说重点/要点 -> concise, 别太正式/随意/口语 -> casual,
//     直接说结论/直接一点/别绕 -> direct.
//   - work hours: needs 工作日 plus 早上|上午 N点 and 下午|晚上 N点; content
//     "weekday HH:MM-HH:MM", afternoon/evening hours shifted +12.
//   - editor: goland | vscode (lowercase match).
//   - family: 我妈/我母亲 (mother) or 我爸/我父亲 (father) plus a M月D号 date;
//     content "mother birthday MM-DD".
type goldenStubExtractor struct{}

var (
	goldenStartRe = regexp.MustCompile(`(?:早上|上午)\s*(\d{1,2})\s*点`)
	goldenEndRe   = regexp.MustCompile(`(?:下午|晚上)\s*(\d{1,2})\s*点`)
	goldenDateRe  = regexp.MustCompile(`(\d{1,2})\s*月\s*(\d{1,2})\s*[号日]`)
)

func (goldenStubExtractor) Extract(_ context.Context, conv GoldenConversation) (ExtractionOutcome, error) {
	out := ExtractionOutcome{}
	pending := map[string]string{}
	for _, t := range conv.Turns {
		if t.Role != "user" {
			continue // assistant turns are never facts about the user
		}
		// One turn, one action: a forget turn skips upsert scanning.
		if slot, ok := goldenForgetTarget(t.Text); ok {
			if !slices.Contains(out.Forgets, slot) {
				out.Forgets = append(out.Forgets, slot)
			}
			delete(pending, slot) // net semantics: forget cancels a same-slot upsert
			continue
		}
		if v := goldenReplyLanguage(t.Text); v != "" {
			pending["pref.language"] = v
		}
		if v := goldenTone(t.Text); v != "" {
			pending["pref.tone"] = v
		}
		if v := goldenWorkHours(t.Text); v != "" {
			pending["habit.work_hours"] = v
		}
		if v := goldenEditor(t.Text); v != "" {
			pending["habit.editor"] = v
		}
		if v := goldenFamily(t.Text); v != "" {
			pending["fact.family"] = v
		}
	}
	for slot, content := range pending {
		out.Upserts = append(out.Upserts, SlotUpsert{Slot: slot, Content: content})
	}
	sort.Slice(out.Upserts, func(i, j int) bool { return out.Upserts[i].Slot < out.Upserts[j].Slot })
	sort.Strings(out.Forgets)
	return out, nil
}

// goldenForgetTarget resolves the "别记了" undo intent: the trigger phrase
// plus a slot keyword in the same turn. No keyword match means the turn is
// noise, not an undo intent.
func goldenForgetTarget(text string) (string, bool) {
	trigger := strings.Contains(text, "别记了") ||
		strings.Contains(text, "别记这条") ||
		strings.Contains(text, "忘掉这条") ||
		strings.Contains(strings.ToLower(text), "forget this")
	if !trigger {
		return "", false
	}
	lower := strings.ToLower(text)
	for _, kw := range []struct {
		slot  string
		words []string
	}{
		{"pref.tone", []string{"说话", "语气", "风格"}},
		{"pref.language", []string{"语言", "中文", "英文"}},
		{"habit.work_hours", []string{"上班", "工作时间", "作息"}},
		{"habit.editor", []string{"编辑器", "ide", "代码工具"}},
		{"fact.family", []string{"我妈", "我爸", "家人"}},
	} {
		for _, w := range kw.words {
			if strings.Contains(lower, w) {
				return kw.slot, true
			}
		}
	}
	return "", false
}

func goldenReplyLanguage(text string) string {
	lower := strings.ToLower(text)
	switch {
	case strings.Contains(lower, "都用中文"),
		strings.Contains(lower, "用中文回"),
		strings.Contains(lower, "中文回复"):
		return "zh"
	case strings.Contains(lower, "reply in english"),
		strings.Contains(lower, "in english please"),
		strings.Contains(lower, "用英文回"):
		return "en"
	}
	return ""
}

func goldenTone(text string) string {
	switch {
	case goldenContainsAny(text, "简短", "简洁", "说重点", "要点"):
		return "concise"
	case goldenContainsAny(text, "别太正式", "随意", "口语"):
		return "casual"
	case goldenContainsAny(text, "直接说结论", "直接一点", "别绕"):
		return "direct"
	}
	return ""
}

func goldenWorkHours(text string) string {
	if !strings.Contains(text, "工作日") {
		return "" // one-off meeting times are not a habit
	}
	start := goldenStartRe.FindStringSubmatch(text)
	end := goldenEndRe.FindStringSubmatch(text)
	if start == nil || end == nil {
		return ""
	}
	s, _ := strconv.Atoi(start[1])
	e, _ := strconv.Atoi(end[1])
	if e < 12 {
		e += 12
	}
	if s >= e {
		return ""
	}
	return fmt.Sprintf("weekday %02d:00-%02d:00", s, e)
}

func goldenEditor(text string) string {
	lower := strings.ToLower(text)
	switch {
	case strings.Contains(lower, "goland"):
		return "goland"
	case strings.Contains(lower, "vscode"),
		strings.Contains(lower, "vs code"),
		strings.Contains(lower, "visual studio code"):
		return "vscode"
	}
	return ""
}

func goldenFamily(text string) string {
	who := ""
	switch {
	case strings.Contains(text, "我妈"), strings.Contains(text, "我母亲"):
		who = "mother"
	case strings.Contains(text, "我爸"), strings.Contains(text, "我父亲"):
		who = "father"
	default:
		return ""
	}
	m := goldenDateRe.FindStringSubmatch(text)
	if m == nil {
		return ""
	}
	month, _ := strconv.Atoi(m[1])
	day, _ := strconv.Atoi(m[2])
	return fmt.Sprintf("%s birthday %02d-%02d", who, month, day)
}

func goldenContainsAny(text string, words ...string) bool {
	for _, w := range words {
		if strings.Contains(text, w) {
			return true
		}
	}
	return false
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
	ex := newGoldenExtractor()
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
		})
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
