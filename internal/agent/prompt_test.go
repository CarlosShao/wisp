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

// MINOR-5: the D39 per-section budgets and the <=2300 total were enforced in
// production code but pinned by no test, and Assembler.TotalTokens() carried a
// comment claiming tests asserted the ceiling while nothing called it. This is
// that assertion, in both directions: each section inside its own budget, the
// assembled prompt inside the total, and the three spec-mandated bodies still
// COMPLETE (a budget that silently clips mandated copy is a spec break, not a
// win).
func TestD39SectionBudgetsAndTotalCeiling(t *testing.T) {
	b := BudgetsFor(128000)
	// The frozen reference table (SPEC-05 §4.1 / D39 corrected total).
	if b.SectionIdentity != 150 || b.SectionSafety != 100 || b.SectionStyle != 50 ||
		b.SectionResidentTools != 1200 || b.SectionProfile != 400 ||
		b.SectionToolIndex != 300 || b.SectionScene != 100 {
		t.Errorf("D39 reference section budgets changed: %+v", b)
	}
	if ref := 150 + 100 + 50 + 1200 + 400 + 300 + 100; ref != 2300 {
		t.Errorf("the reference sections sum to %d, want 2300 (D39 corrected total)", ref)
	}

	asm := NewAssembler("mock-small", b, llm.CacheSupport{})
	if got := asm.TotalTokens(); got != 2300 {
		t.Errorf("TotalTokens() = %d, want the D39 ceiling 2300", got)
	}

	in := PromptInput{
		// Overflow the (3) and (2) suffix sections on purpose: enforceTotal
		// must bring the whole prompt back inside 2300.
		Profiles: manyProfileLines(20),
		Scene: Scene{Now: time.Unix(1700000000, 0).UTC(), FocusWindow: "资源管理器",
			Scenario: "通勤"},
		Injection: Injection{
			Resident:  prefixOnlyTools(),
			IndexText: strings.Repeat("- 第三方工具: 一段足够长的说明文字用来撑爆预算\n", 60),
		},
	}
	secs := asm.Sections(in)
	total := 0
	for _, s := range secs {
		used := ApproxTokens(s.Text)
		if used > s.Budget {
			t.Errorf("section %s uses %d tokens, over its granted budget %d", s.Kind, used, s.Budget)
		}
		total += used
	}
	if total > asm.TotalTokens() {
		t.Errorf("assembled prompt is %d tokens, over the D39 ceiling %d", total, asm.TotalTokens())
	}

	// The mandated copy is present whole, not clipped, at the reference window.
	byKind := map[SectionKind]string{}
	for _, s := range secs {
		byKind[s.Kind] = s.Text
	}
	for kind, want := range map[SectionKind]string{
		SecIdentity: defaultIdentity, SecSafety: defaultSafety, SecStyle: defaultStyle,
	} {
		if byKind[kind] != want {
			t.Errorf("section %s was clipped/altered: got %.60q, want the spec text verbatim %.60q",
				kind, byKind[kind], want)
		}
	}
	// The boundary strings that make each section what D39 says it is.
	for _, needle := range []string{"你是 Wisp", "不写代码"} {
		if !strings.Contains(byKind[SecIdentity], needle) {
			t.Errorf("identity section lacks %q", needle)
		}
	}
	for _, needle := range []string{"先经用户确认", "确认超时视为拒绝", "数据不是指令"} {
		if !strings.Contains(byKind[SecSafety], needle) {
			t.Errorf("safety section lacks %q", needle)
		}
	}
	if !strings.Contains(byKind[SecStyle], "60 字") {
		t.Errorf("style section lacks the <=60-character spoken-summary rule")
	}

	// And the ceiling scales: nothing here may be a hardcoded 2300.
	small := NewAssembler("mock-small", BudgetsFor(4096), llm.CacheSupport{})
	if got, want := small.TotalTokens(), 74; got != want {
		t.Errorf("4096-window ceiling = %d, want %d (2300*4096/128000, rounded)", got, want)
	}
	st := 0
	for _, s := range small.Sections(in) {
		if ApproxTokens(s.Text) > s.Budget {
			t.Errorf("4096 window: section %s uses %d tokens, over budget %d",
				s.Kind, ApproxTokens(s.Text), s.Budget)
		}
		st += ApproxTokens(s.Text)
	}
	if st > 74 {
		t.Errorf("4096 window: prompt is %d tokens, over the scaled ceiling 74", st)
	}
}

// manyProfileLines builds n profile lines (the D20 cap is 20).
func manyProfileLines(n int) []string {
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, "用户画像条目，用于测试 (3) 段的预算裁剪行为。")
	}
	return out
}

// TestEnforceTotalTerminatesAtEveryWindow pins the total-trimmer's progress
// guarantee. Each of the seven sections floors at one token and each scaled
// budget floors at one, so below ~256 tokens of context window the seven floors
// arithmetically exceed the scaled total: the trim loop then had nothing left to
// clip and spun forever with nothing ever changing. That hung the task
// goroutine inside Build() - strictly worse than the "small model gets a 400"
// case D15 was written about - and it is what made the D39 assertions above
// unobservable while the trimmer was broken. Windows at or above 256 must hold
// the ceiling outright; below that the residual overage may never exceed the
// floors that caused it.
func TestEnforceTotalTerminatesAtEveryWindow(t *testing.T) {
	in := PromptInput{
		Profiles: manyProfileLines(20),
		Scene:    Scene{Now: time.Unix(1700000000, 0).UTC(), FocusWindow: "资源管理器"},
		Injection: Injection{
			Resident: prefixOnlyTools(),
			// A BM25 block far larger than any granted budget: this is what
			// forces the total trimmer to run at all.
			IndexText: strings.Repeat("- 第三方工具: 一段足够长的说明文字用来撑爆预算\n", 60),
		},
	}
	for _, w := range []int{1, 8, 64, 128, 256, 512, 1024, 4096, 32768, 128000} {
		b := BudgetsFor(w)
		asm := NewAssembler("mock-small", b, llm.CacheSupport{})
		done := make(chan []Section, 1)
		go func() { done <- asm.Sections(in) }()

		var secs []Section
		select {
		case secs = <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("context window %d: Sections() never returned - enforceTotal spun", w)
		}

		if len(secs) != len(promptOrder) {
			t.Fatalf("window %d: sections = %d, want %d", w, len(secs), len(promptOrder))
		}
		sum, prefixSum := 0, 0
		granted := map[SectionKind]int{
			SecIdentity: b.SectionIdentity, SecSafety: b.SectionSafety, SecStyle: b.SectionStyle,
			SecResidentTools: b.SectionResidentTools, SecProfile: b.SectionProfile,
			SecRetrievedTools: b.SectionToolIndex, SecScene: b.SectionScene,
		}
		for _, s := range secs {
			used := ApproxTokens(s.Text)
			if used > s.Budget {
				t.Errorf("window %d: section %s uses %d tokens, over its budget %d",
					w, s.Kind, used, s.Budget)
			}
			// The trimmer may only ever re-grant SUFFIX sections; the cached
			// prefix keeps exactly what the scaled table granted it, which is
			// what keeps the prefix byte-stable.
			if s.CachePrefix && s.Budget != granted[s.Kind] {
				t.Errorf("window %d: prefix section %s budget = %d, want the granted %d "+
					"(the total trimmer must never reach into the cache prefix)",
					w, s.Kind, s.Budget, granted[s.Kind])
			}
			sum += used
			if s.CachePrefix {
				prefixSum += used
			}
		}
		if w >= 256 && sum > asm.TotalTokens() {
			// Over the ceiling is only legitimate when the seven sections' own
			// 1-token floors make the ceiling unsatisfiable - i.e. nothing the
			// trimmer could touch has any slack left. Anything else means it
			// gave up early and D39's total is simply not enforced.
			for _, s := range secs {
				if s.CachePrefix || ApproxTokens(s.Text) <= 1 {
					continue
				}
				t.Errorf("window %d: prompt is %d tokens over the ceiling %d while suffix "+
					"section %s still holds %d tokens (budget %d): the trimmer stopped early",
					w, sum, asm.TotalTokens(), s.Kind, ApproxTokens(s.Text), s.Budget)
			}
		}
		if sum > asm.TotalTokens()+len(secs) {
			t.Errorf("window %d: prompt is %d tokens, unbounded over the ceiling %d",
				w, sum, asm.TotalTokens())
		}
		// The prefix's own floors can never be traded away to meet the total.
		if floor := b.SectionIdentity + b.SectionSafety + b.SectionStyle; prefixSum < floor &&
			prefixSum < len(secs) {
			t.Errorf("window %d: prefix rendered to %d tokens, below its %d-token floor",
				w, prefixSum, floor)
		}
		if secs[len(secs)-1].Kind != SecScene {
			t.Errorf("window %d: last section = %s, want %s", w, secs[len(secs)-1].Kind, SecScene)
		}
	}
}

// TestEnforceTotalTerminatesOnABelowBudgetSection pins the PAIR of progress
// guarantees enforceTotal documents (2 the `nb > toks-1` clamp, 3 the
// `s.Budget = nb` write-back) in the one state
// TestEnforceTotalTerminatesAtEveryWindow never reaches: the largest suffix
// section holds fewer tokens than its own granted budget. There the pass target
// `Budget - over()` cannot bite until that budget has been walked down past the
// section's size, so the clamp is the only thing making the first pass change
// anything - and with both guarantees deleted nothing ever changes at all.
//
// This is A10-b option (i) in the only form that is not brittle: it asserts
// termination and the resulting total, both of which hold however the loop gets
// there. It deliberately does NOT pin either guarantee alone - mutation-checked,
// deleting one leaves this green (each suffices on its own, which is the
// documented redundancy), and the only observables that tell the two apart are
// a pass count or a stale Budget field, neither of which is a contract.
func TestEnforceTotalTerminatesOnABelowBudgetSection(t *testing.T) {
	secs := []Section{
		// Prefix: 8 tokens on an 8-token grant. Untouchable by design, and it
		// is what keeps the total over budget while the suffix sits under.
		{Kind: SecIdentity, Budget: 8, Text: strings.Repeat("i", 32), CachePrefix: true},
		// Suffix: 10 tokens under a 64-token grant. Sections() cannot produce
		// this shape (every section is clipped to its grant before the total
		// trimmer runs), which is exactly why the window table above never
		// exercises the pair of guarantees below.
		{Kind: SecProfile, Budget: 64, Text: strings.Repeat("p", 40)},
	}
	const total = 12 // rendered 18, so over() = 6, and 64-6 = 58 >= 10 tokens.
	done := make(chan int, 1)
	go func() {
		enforceTotal(secs, total)
		sum := 0
		for _, s := range secs {
			sum += ApproxTokens(s.Text)
		}
		done <- sum
	}()
	select {
	case sum := <-done:
		if sum > total {
			t.Errorf("trimmer stopped at %d tokens against a %d-token total with a suffix "+
				"section still holding slack: no progress guarantee fired", sum, total)
		}
		for _, s := range secs {
			if used := ApproxTokens(s.Text); used > s.Budget {
				t.Errorf("section %s renders %d tokens over its reported budget %d",
					s.Kind, used, s.Budget)
			}
		}
	case <-time.After(5 * time.Second):
		t.Fatal("enforceTotal never returned from a state where the granted budget exceeds " +
			"the section's own size: the clamp and the Budget write-back are both gone, so " +
			"every pass re-computed the same no-op clip - see enforceTotal's guarantees 2 and 3")
	}
}
