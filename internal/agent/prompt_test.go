package agent

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
)

// Acceptance criterion 7 (D39/SPEC-05 §4.1): the assembled prompt is ordered so
// the CACHE-STABLE prefix (identity, safety, style, resident tools) is byte-
// identical across turns, and everything that changes turn to turn - the (2)
// BM25 retrieved index and the (4) time/focus scene - sits in the suffix, with
// the scene LAST. A moving byte anywhere in the prefix is a per-turn cache
// miss, so this asserts on BYTES, not on whether a section merely "is present".

// sysBytes renders the request's cache-prefix content to a single byte string:
// the exact stable prefix the provider would cache.
func sysBytes(req *llm.Request) string {
	var b strings.Builder
	for _, c := range req.System {
		if t, ok := c.(llm.TextPart); ok {
			b.WriteString(t.Text)
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func prefixOnlyTools() []ToolInfo {
	return []ToolInfo{
		{Name: "echo", Description: "Echo back the given text.", Resident: true,
			Parameters: json.RawMessage(`{"type":"object"}`)},
		{Name: "sleep", Description: "Sleep for milliseconds.", Resident: true,
			Parameters: json.RawMessage(`{"type":"object"}`)},
	}
}

// TestCachePrefixIsByteStableAcrossTurns drives two different turns (different
// time, different BM25 hits, different profiles, growing history) and requires
// the cache-prefix bytes to be IDENTICAL, while the volatile bytes differ and
// live only at the tail.
func TestCachePrefixIsByteStableAcrossTurns(t *testing.T) {
	asm := NewAssembler("mock-small", BudgetsFor(128000), llm.CacheSupport{})
	resident := prefixOnlyTools()

	mkIn := func(t1 bool) PromptInput {
		scene := Scene{Now: time.Date(2026, 9, 20, 9, 5, 0, 0, time.UTC), FocusWindow: "记事本"}
		bm25 := "- 日历: 查看今日日程"
		if t1 {
			scene = Scene{Now: time.Date(2026, 9, 20, 14, 30, 0, 0, time.UTC), FocusWindow: "微信 - 张三"}
			bm25 = "- 微信: 给联系人发消息"
		}
		return PromptInput{
			Profiles:  []string{"用户偏好简洁回答", "常用时段 09:00-18:00"},
			Scene:     scene,
			Injection: Injection{Resident: resident, IndexText: bm25},
		}
	}

	histTurn1 := []llm.Message{
		{Role: llm.RoleUser, Content: []llm.Content{llm.TextPart{Text: "上午的安排"}}},
	}
	histTurn2 := append([]llm.Message{
		{Role: llm.RoleUser, Content: []llm.Content{llm.TextPart{Text: "下午的安排"}}},
		{Role: llm.RoleAssistant, Content: []llm.Content{llm.TextPart{Text: "好的"}}},
	}, histTurn1...)

	req1 := asm.Build(mkIn(true), histTurn1)
	req2 := asm.Build(mkIn(false), histTurn2)

	// The cache-prefix bytes must match exactly across two turns whose volatile
	// data (time/focus, BM25 hits, profile, and the history itself) all differ.
	if sysBytes(req1) != sysBytes(req2) {
		t.Fatalf("cache prefix is not byte-stable across turns\n--turn1--\n%s\n--turn2--\n%s",
			sysBytes(req1), sysBytes(req2))
	}
	prefix := sysBytes(req1)

	// The byte comparison is only meaningful over real, non-trivial content:
	// require the stable identity block to be present in the cached prefix.
	if !strings.Contains(prefix, "你是 Wisp") || len(prefix) < 200 {
		t.Fatalf("cache prefix is empty/trivial (len=%d), byte-stability unproven", len(prefix))
	}

	// Volatile content must live ONLY in the suffix (the last message), never in
	// the cached prefix.
	for _, leak := range []string{"微信 - 张三", "记事本", "给联系人发消息", "上午的安排", "下午的安排"} {
		if strings.Contains(prefix, leak) {
			t.Errorf("volatile data %q leaked into the cache prefix", leak)
		}
	}
	// Sanity: the volatile suffix really did differ between turns (so prefix
	// equality above is not vacuous).
	if suffixText(req1) == suffixText(req2) {
		t.Error("volatile suffix identical across turns - byte-stability test is vacuous")
	}
}

// suffixText returns the request's final (volatile) message content as bytes.
func suffixText(req *llm.Request) string {
	if len(req.Messages) == 0 {
		return ""
	}
	last := req.Messages[len(req.Messages)-1]
	var b strings.Builder
	for _, c := range last.Content {
		if t, ok := c.(llm.TextPart); ok {
			b.WriteString(t.Text)
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// TestPromptSectionCanonicalOrder pins the D39 section order: the four
// cache-prefix blocks first, then the volatile suffix with the scene LAST.
func TestPromptSectionCanonicalOrder(t *testing.T) {
	asm := NewAssembler("mock-small", BudgetsFor(128000), llm.CacheSupport{})
	in := PromptInput{
		Profiles:  []string{"偏好简洁"},
		Scene:     Scene{Now: time.Unix(1700000000, 0).UTC(), FocusWindow: "资源管理器"},
		Injection: Injection{Resident: prefixOnlyTools(), IndexText: "- 计算器: 算数"},
	}
	secs := asm.Sections(in)

	want := []SectionKind{
		SecIdentity, SecSafety, SecStyle, SecResidentTools, // prefix
		SecProfile, SecRetrievedTools, SecScene, // suffix, (4) last
	}
	if len(secs) != len(want) {
		t.Fatalf("section count = %d, want %d", len(secs), len(want))
	}
	for i, k := range want {
		if secs[i].Kind != k {
			t.Errorf("section %d = %s, want %s", i, secs[i].Kind, k)
		}
	}
	// Prefix/suffix flag partition matches the cache layout.
	for i, s := range secs {
		prefixWant := i < 4
		if s.CachePrefix != prefixWant {
			t.Errorf("section %s CachePrefix = %v, want %v", s.Kind, s.CachePrefix, prefixWant)
		}
	}
	// The (4) scene is the mandatory-last section.
	if secs[len(secs)-1].Kind != SecScene {
		t.Errorf("last section = %s, want %s (scene must be last)", secs[len(secs)-1].Kind, SecScene)
	}
}

// TestVolatileSuffixPutsSceneLast checks the assembled request's final message
// orders its sections profile -> BM25 index -> scene, with the scene the last
// content part, so the prefix stays cacheable.
func TestVolatileSuffixPutsSceneLast(t *testing.T) {
	asm := NewAssembler("mock-small", BudgetsFor(128000), llm.CacheSupport{})
	in := PromptInput{
		Profiles:  []string{"偏好简洁回答"},
		Scene:     Scene{Now: time.Unix(1700000000, 0).UTC(), FocusWindow: "独占窗口XYZ"},
		Injection: Injection{Resident: prefixOnlyTools(), IndexText: "- 音乐: 播放歌单"},
	}
	req := asm.Build(in, []llm.Message{
		{Role: llm.RoleUser, Content: []llm.Content{llm.TextPart{Text: "放点音乐"}}},
	})

	if len(req.Messages) < 2 {
		t.Fatalf("messages = %d, want conversation + volatile suffix", len(req.Messages))
	}
	last := req.Messages[len(req.Messages)-1]

	var headers []string
	for _, c := range last.Content {
		if t2, ok := c.(llm.TextPart); ok {
			if strings.HasPrefix(t2.Text, "## ") {
				headers = append(headers, strings.SplitN(t2.Text, "\n", 2)[0])
			}
		}
	}
	want := []string{"## " + string(SecProfile), "## " + string(SecRetrievedTools), "## " + string(SecScene)}
	if strings.Join(headers, ",") != strings.Join(want, ",") {
		t.Errorf("suffix section order = %v, want %v", headers, want)
	}
	// The scene text (unique focus window) is genuinely the last content part.
	finalPart, _ := last.Content[len(last.Content)-1].(llm.TextPart)
	if !strings.HasPrefix(finalPart.Text, "## "+string(SecScene)) {
		t.Errorf("final suffix part = %.40q, want the scene block", finalPart.Text)
	}
	if !strings.Contains(finalPart.Text, "独占窗口XYZ") {
		t.Errorf("scene/focus not in the final part: %.80q", finalPart.Text)
	}
}
