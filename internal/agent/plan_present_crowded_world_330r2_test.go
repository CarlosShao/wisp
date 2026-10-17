package agent

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/memory"
)

// Ticket 330, leg 330-r2: the "multi-candidate world" that gap G2 of
// docs/evidence/s1/165-v2-verdict.md says is missing. NOTHING here changes a
// production byte: the two cases below only read internal/agent/inject.go as it
// already behaves.
//
// WHAT WAS BROKEN (G2, quoted shape). internal/agent/plan_present_doorbell_165r3_test.go
// mirrors the frozen D34 row of plan.present and hands that mirror to the loop.
// The mirror carries a residency switch (Resident: true at :69 of that file).
// 165-v2 flipped it to false and the whole package stayed green, rc=0. The
// mechanism, re-verified at this anchor rather than taken on faith: Inject
// (internal/agent/inject.go:44-52) splits the directory into the resident arm
// and the retrieval arm, and bm25Select (internal/agent/inject.go:172-245)
// keeps EVERY candidate whose score is above zero (:219). The 165 case's
// directory holds exactly one candidate, and the query "帮我把这个项目收拾一下"
// shares the unigrams 我/一/下 with the mirrored description "先交一份「我要干
// 什么」的清单并停下等点头" - so the brake is retrieved anyway, lands in the
// round-1 body anyway, and the assertion at :219 of that file never sees the
// tier. The switch was decoration because the world had nothing to compete with.
//
// THE WORLD THIS FILE ADDS. One more entry in the same directory: a long-tail
// NON-RESIDENT candidate whose index line is deliberately big. bm25Select
// clips the retrieval arm at the D15(2) index budget (internal/agent/
// inject.go:232-243, SectionToolIndex = 300 tokens at the 128k reference
// window): the highest-scoring candidate is admitted unconditionally, and any
// later hit is dropped once used+cost exceeds the budget. The decoy outscores
// the mirrored brake by a wide margin (it shares 帮/把/这/个/项/目/收/拾 plus
// the bigrams 帮我/这个/个项/项目/收拾 with the query, where the brake shares
// only 我/一/下), so in this world the two non-resident candidates cannot both
// occupy the capped slot. That is what makes the residency switch decisive:
//
//	resident arm          -> the brake is injected no matter what the slot says
//	retrieval arm (capped) -> the brake loses the slot to the decoy and is NOT
//	                          offered to the model at all
//
// THE TWO CASES.
//
// TestCrowdedWorldStillOffersTheBrake330r2 drives the genuine loop over this
// directory and reads the round-1 request body: it asserts only that the brake
// is offered and that the doorbell row is still booked. It never names the
// decoy, so deleting the decoy from the fixture returns this case to the 165
// single-candidate world byte-for-byte - which is the positive control that the
// tooth below hangs on the world, not on a hard-coded expectation.
//
// TestResidencyTierDecidesTheInjectionArm330r2 reads the arms of Inject
// directly. No case in this package called Inject directly at this anchor
// (grep 'Inject(' over *_test.go returns nothing; the only production caller is
// internal/agent/loop.go:573), so this is the first instrument that pins WHICH
// arm a mirrored entry lands in. Its gates are ordered so that the three
// mutation recipes in the dispatch book each trip a DIFFERENT line:
// premise, then residency of the brake, then the decoy's arm, then the
// decoy's retrieval.
//
// WHAT THIS FILE DOES NOT CLAIM. The decoy is a retrieval-only fixture entry.
// It is not a mirror of a real roster row and asserts nothing about
// internal/tools: the mirror-to-real-roster reconciliation is gap G1, owned by
// ticket 330 AC#1/AC#2 and by leg 330-r1 in cmd/wisp/**, which package agent
// cannot reach (internal/tools/bridge.go:13 already imports internal/agent).
// Declared is NOT re-teethed here either: 165-v2 measured that switch as
// already armed (M-b2 and M-g), so a second nail on it would be a duplicate.

// crowdedWorldQuery330r2 is the same user text the 165 doorbell case runs; the
// literal is restated so this file's mutations can be applied to it alone.
const crowdedWorldQuery330r2 = "帮我把这个项目收拾一下"

// longTailName330r2 is the fixture's second, NON-resident candidate. The name
// is Latin on purpose: Tokenize (internal/agent/inject.go:127-163) yields no
// CJK terms from it, so every shared term it scores on comes from the
// description below, which keeps the score arithmetic readable.
const longTailName330r2 = "ledger.sweep330r2"

// longTailDesc330r2 is sized, not incidental: its index line costs 291 of the
// 300 retrieval-budget tokens, and the mirrored brake's line costs 20, so the
// pair (311) cannot both fit. Premise gate 2 of the arm case pins that
// arithmetic so a later edit of this description cannot silently disarm it.
const longTailDesc330r2 = "扫描这个项目的根目录与全部子目录，帮我把散落的构建产物、过期脚本和临时压缩包逐一登记，" +
	"这个项目里每一份文件的完整路径、字节大小与最后一次修改时间都要列进长清单，" +
	"清单按目录分组再按大小排序，名字带着旧日期前缀的单独收拢到一处，" +
	"牵连到的 README 段落要在旁边标注，方便你逐条核对，" +
	"我不写入也不删除任何文件，只把这个项目当前收拾到一半的样子汇报给你，" +
	"等你点头之后这一轮才动第一份文件，动之前会把计划贴在这条清单上，" +
	"每一条待办都要留下可复跑的检核方式与预计影响范围，" +
	"跑完再把结果按同样的顺序回报一次。" +
	"凡是体积超过一兆的产物都要单独点名，包括缓存目录里那份没人认领的旧索引，" +
	"每一条都要说清它在项目里的角色，是脚本、是产物、还是收拾时顺手挪走的临时件，" +
	"这份清单也要把我一开始提过的那句要求复述一遍，免得核对时要回去翻记录，" +
	"全部登记完毕再把总数、总大小与最占地方的三项一起报回来。"

// crowdedDoorbellWorld330r2 is the 165 mirror PLUS the long-tail candidate. It
// composes planPresentDirectory165r3() instead of copying it, so the frozen row
// stays mirrored in exactly one place (a second copy would be a second mirror
// for gap G1 to reconcile).
func crowdedDoorbellWorld330r2() []ToolInfo {
	return append(planPresentDirectory165r3(), ToolInfo{
		Name:        longTailName330r2,
		Description: longTailDesc330r2,
		Resident:    false,
		RiskLevel:   RiskL1,
		Parameters:  json.RawMessage(`{"type":"object","properties":{"root":{"type":"string"}},"required":["root"],"additionalProperties":false}`),
	})
}

// indexLineCost330r2 repeats the cost formula bm25Select uses to clip the
// retrieval index (internal/agent/inject.go:237) so the premise of the crowded
// world is stated as a number instead of a hope.
func indexLineCost330r2(t ToolInfo) int {
	return ApproxTokens("- "+t.Name+": "+strings.TrimSpace(t.Description)) + 1
}

func armHas330r2(tools []ToolInfo, name string) bool {
	for _, t := range tools {
		if t.Name == name {
			return true
		}
	}
	return false
}

func defsHave330r2(defs []llm.ToolDef, name string) bool {
	for _, d := range defs {
		if d.Name == name {
			return true
		}
	}
	return false
}

func defNamesOf330r2(defs []llm.ToolDef) []string {
	out := make([]string, 0, len(defs))
	for _, d := range defs {
		out = append(out, d.Name)
	}
	return out
}

func namesOf330r2(tools []ToolInfo) []string {
	out := make([]string, 0, len(tools))
	for _, t := range tools {
		out = append(out, t.Name)
	}
	return out
}

// TestCrowdedWorldStillOffersTheBrake330r2 is the world-side tooth: the genuine
// Loop, the genuine openai-chat adapter and the genuine store, over a directory
// with two candidates. The brake is resident, so the capped retrieval slot is
// not its slot to win, and it must be offered in round 1.
//
// THE TOOTH: replay gap G2's mutation in this tree (the mirror entry's
// Resident: true -> false at internal/agent/plan_present_doorbell_165r3_test.go:69).
// The brake then competes for the retrieval slot, loses it to the decoy, and is
// never put on the wire: this case goes red at the mustContain below. In the
// 165 single-candidate world the same mutation is invisible (the lone candidate
// is admitted unconditionally), which is exactly the gap this file closes.
func TestCrowdedWorldStillOffersTheBrake330r2(t *testing.T) {
	dir := sealableTempDir124(t)
	store, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	defer store.Close()

	h := newHarness(t, planPresentFixture165r3,
		withTools(NewEchoProvider(crowdedDoorbellWorld330r2()...)),
		withJournal(store),
		withConfig(func(c *Config) { c.ArtifactsDir = filepath.Join(dir, "artifacts") }))
	res := h.run(crowdedWorldQuery330r2)

	bodies := h.requestBodies()
	if len(bodies) != 2 {
		t.Fatalf("provider requests = %d, want 2 (the recorded call must open a second round)", len(bodies))
	}
	// THE JUDGMENT of this case: the model was offered this door. This needle is
	// the wire form of internal/agent/plan_present_doorbell_165r3_test.go:219,
	// which is the line the 165-v2 verdict named as blind to the tier.
	mustContain(t, "330r2 crowded world: round 1 body must offer the brake",
		string(bodies[0]), `"name":"`+planPresentName165r3+`"`)

	// The doorbell itself still rings over the crowded directory: the sink books
	// the row before the risk switch, and the ungated host still refuses it.
	rows, err := store.ListToolCallsByTask(context.Background(), res.TaskID)
	if err != nil {
		t.Fatalf("ListToolCallsByTask: %v", err)
	}
	if n := sinkCount165r3(rows, planPresentName165r3); n != 1 {
		t.Errorf("sink reading for %q = %d, want 1 over %d rows", planPresentName165r3, n, len(rows))
	}
	row := rowOfTool165r3(t, rows, planPresentName165r3)
	if row.Decision != DecisionReject {
		t.Errorf("row.decision = %q, want %q (an ungated host must not run an L2 write)",
			row.Decision, DecisionReject)
	}
	if got := h.echo().CallCount(); got != 0 {
		t.Errorf("executed tool calls = %d, want 0", got)
	}
}

// TestResidencyTierDecidesTheInjectionArm330r2 reads the two arms of Inject for
// the same world. Unlike the case above it names the decoy, so its gates are
// the arm-bookkeeping half of the instrument, not the world half.
func TestResidencyTierDecidesTheInjectionArm330r2(t *testing.T) {
	b := BudgetsFor(ReferenceContextWindow)
	world := crowdedDoorbellWorld330r2()

	// Gate 1 (premise): the mirror plus exactly one long-tail candidate.
	if len(world) != 2 {
		t.Fatalf("330r2 premise: the crowded world must be the frozen mirror plus exactly one retrieval decoy, got %d entries %+v",
			len(world), namesOf330r2(world))
	}
	// Gate 2 (premise): the slot really is scarce - both lines fit the budget on
	// their own, the pair does not.
	costBrake, costDecoy := indexLineCost330r2(world[0]), indexLineCost330r2(world[1])
	if costDecoy > b.SectionToolIndex {
		t.Fatalf("330r2 premise: the decoy index line alone costs %d, over the SectionToolIndex budget %d, so the retrieval arm would be empty even in the world with no mutation at all",
			costDecoy, b.SectionToolIndex)
	}
	if costBrake+costDecoy <= b.SectionToolIndex {
		t.Fatalf("330r2 premise: the two index lines cost %d+%d = %d, within the SectionToolIndex budget %d, so nothing competes for the slot and the tier cannot decide anything",
			costBrake, costDecoy, costBrake+costDecoy, b.SectionToolIndex)
	}

	inj := Inject(crowdedWorldQuery330r2, world, b)

	// Gate 3 (THE G2 SWITCH): the mirror declares the brake resident, so the
	// resident arm must carry it. Replay of gap G2's mutation (Resident true ->
	// false in plan_present_doorbell_165r3_test.go:69) trips exactly this line.
	if !armHas330r2(inj.Resident, planPresentName165r3) {
		t.Fatalf("residency switch ignored: resident arm = %+v, want %q in it - the mirror flags it Resident: true and Inject put it somewhere else",
			namesOf330r2(inj.Resident), planPresentName165r3)
	}
	// Gate 4 (de-fang shot): the decoy is declared NON-resident, so it must not
	// be handed the unconditional arm. Flipping Resident: false -> true on the
	// decoy above trips exactly this line.
	if armHas330r2(inj.Resident, longTailName330r2) {
		t.Fatalf("tier leak: resident arm = %+v, want %q kept out of it - the fixture declares it Resident: false and only the tier decides this arm",
			namesOf330r2(inj.Resident), longTailName330r2)
	}
	// Gate 5 (control shot): the decoy is in the world because retrieval put it
	// there. Replacing crowdedWorldQuery330r2 with a query that shares no term
	// with it trips exactly this line, and only this line.
	if !armHas330r2(inj.Retrieved, longTailName330r2) {
		t.Fatalf("retrieval arm = %+v, want %q in it: it scores above zero against %q, so the capped slot is the only thing that could have dropped it",
			namesOf330r2(inj.Retrieved), longTailName330r2, crowdedWorldQuery330r2)
	}

	// Downstream of both arms, non-fatally: the wire set is the two arms, and a
	// retrieval hit means no list_tools fallback (internal/agent/inject.go:57-72).
	if !defsHave330r2(inj.ToolDefs, planPresentName165r3) {
		t.Errorf("tool defs = %+v, want %q on the wire", defNamesOf330r2(inj.ToolDefs), planPresentName165r3)
	}
	if !defsHave330r2(inj.ToolDefs, longTailName330r2) {
		t.Errorf("tool defs = %+v, want %q on the wire", defNamesOf330r2(inj.ToolDefs), longTailName330r2)
	}
	if inj.AdvertiseListTools {
		t.Errorf("AdvertiseListTools = true, want false: the retrieval hit means the D15(2) fallback is not offered")
	}
	if len(inj.IndexText) == 0 {
		t.Errorf("IndexText is empty, want the retrieved decoy rendered into the (2) suffix index")
	}
}
