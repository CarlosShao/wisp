// File: internal/risk/markseam_394_test.go — 票 394: mark 观察缝测试（第三形）
//
// Status: 实施中（394-r1 写腿 2026-11-01）
// Environment: 本地全套（纯 Go，无 staging 依赖）；gofumpt 绝对路径；突变走仓外 -overlay
//
// 三条硬边界遵守：
// 1. ⛔ 改 CheckText/Scope.Close/RecordTaintInheritance 既有契约 —— 本文件只测新增缝
// 2. ⛔ internal/tools/** —— 4 处生产 mark 点的汇流口（Mark/MarkWithHostPath）在此直接驱动
// 3. 契约面只在 risk —— 观察者接口未导出，白盒断言读 gate.inherited

package risk

import (
	"strings"
	"testing"
)

// seam394Raw mixes every normalization trigger the contract normalizer folds:
// ASCII whitespace, the ideographic space, a zero-width space, mixed case and
// full-width Latin/digits/punctuation. The announcement quotes the HALF-WIDTH
// spelling, so raw-vs-raw the two share nothing — only the NORMALIZED forms
// meet. A seam that recorded raw content could never produce the refusal the
// tests below demand, which is what pins "normalizeTaint 履行在缝内".
const seam394Raw = "  SECRET　Ｔｏｋｅｎ－３９４Ｒ１ＳＥＡＭ\u200b（十分钟内有效）  "

// seam394Quote is the announcement that quotes the taint's normalized form.
const seam394Quote = "令牌是 token394r1seam，请确认"

// Test394SeamRecordsNormalizedContentAtMarkTime is AC#1's core reading: a plain
// Mark with RAW content lands in the chain already normalized, tool kept
// parallel, and the closed task scope's taint refuses the announcement through
// form B with the source naming the tool — no RecordTaintInheritance call
// anywhere in the test.
func Test394SeamRecordsNormalizedContentAtMarkTime(t *testing.T) {
	prov := NewProvenance(ProvOptions{NoProbe: true})
	gate := NewTTSGate(prov, DefaultTTSGateConfig())

	const taskID = "task-394r1-unit"
	scope := prov.OpenScope(taskID)
	if !prov.Mark(taskID, SrcSearchContent, "search", seam394Raw) {
		t.Fatal("Mark refused the raw taint")
	}

	if got := gate.CountInheritanceLinks(); got != 1 {
		t.Fatalf("seam recorded %d links, want 1 (one link per scope, engine-side posture)", got)
	}
	link := gate.inherited[0]
	if link.TaskScopeID != taskID {
		t.Errorf("link scope = %q, want %q", link.TaskScopeID, taskID)
	}
	if want := normalizeTaint(seam394Raw); len(link.TaintContent) != 1 || link.TaintContent[0] != want {
		t.Errorf("chain content = %q, want the normalized form %q (raw content must never be stored)",
			link.TaintContent, want)
	}
	if len(link.ToolSources) != 1 || link.ToolSources[0] != SrcSearchContent {
		t.Errorf("tool sources = %v, want [%s]", link.ToolSources, SrcSearchContent)
	}
	if link.MarkedAt == 0 {
		t.Error("MarkedAt was not stamped")
	}

	if err := scope.Close(); err != nil {
		t.Fatalf("close the task scope: %v", err)
	}
	hit, refused, source := gate.CheckTextWithInheritance(taskID, ChTTS, seam394Quote)
	if !hit || !refused {
		t.Fatalf("seam-fed form B did not refuse the quoting sentence (hit=%v refused=%v)", hit, refused)
	}
	if !strings.Contains(source, SrcSearchContent) || !strings.Contains(source, "from closed task scope") {
		t.Errorf("refusal source = %q, want tool name + closed-scope wording", source)
	}
}

// Test394SeamUpsertsOneLinkPerScope pins the retention posture: repeated marks
// in one scope fold into ONE link with parallel content/tool slices, a second
// scope gets its own link, and the exported RecordTaintInheritance still
// APPENDS (its contract untouched — never upserted into the seam's link).
func Test394SeamUpsertsOneLinkPerScope(t *testing.T) {
	prov := NewProvenance(ProvOptions{NoProbe: true})
	gate := NewTTSGate(prov, DefaultTTSGateConfig())

	prov.OpenScope("task-394a")
	prov.OpenScope("task-394b")
	prov.Mark("task-394a", SrcFSRead, "o1", "AAAAAAAAAA1")
	prov.Mark("task-394a", SrcWebFetch, "o2", "BBBBBBBBBB2")
	prov.Mark("task-394b", SrcDocRead, "o3", "CCCCCCCCCC3")

	if got := gate.CountInheritanceLinks(); got != 2 {
		t.Fatalf("seam recorded %d links, want 2 (one per scope)", got)
	}
	var linkA *taintInheritanceLink
	for _, l := range gate.inherited {
		if l.TaskScopeID == "task-394a" {
			linkA = l
		}
	}
	if linkA == nil {
		t.Fatal("no link for task-394a")
	}
	if len(linkA.TaintContent) != 2 || len(linkA.ToolSources) != 2 {
		t.Fatalf("task-394a link holds %d taints / %d tools, want 2/2 (parallel slices)",
			len(linkA.TaintContent), len(linkA.ToolSources))
	}
	if linkA.TaintContent[0] != normalizeTaint("AAAAAAAAAA1") || linkA.TaintContent[1] != normalizeTaint("BBBBBBBBBB2") {
		t.Errorf("taint order drifted: %v", linkA.TaintContent)
	}
	if linkA.ToolSources[0] != SrcFSRead || linkA.ToolSources[1] != SrcWebFetch {
		t.Errorf("tool order drifted: %v", linkA.ToolSources)
	}

	gate.RecordTaintInheritance("task-394a", [][]rune{[]rune("manual394record")}, []string{SrcWebFetch})
	if got := gate.CountInheritanceLinks(); got != 3 {
		t.Fatalf("after the explicit record: %d links, want 3 — RecordTaintInheritance appends, it does not upsert", got)
	}
}

// Test394SeamCoversMarkWithHostPathSite drives the exact funnel production's
// four mark sites use (bridge.go:687 goes through MarkWithHostPath). The seam
// records the full normalized body, and the ticket 177/183 exemption machinery
// is untouched: windows spelling the declared host path stay non-evidence from
// this mark, body windows stay evidence.
func Test394SeamCoversMarkWithHostPathSite(t *testing.T) {
	prov := NewProvenance(ProvOptions{NoProbe: true})
	gate := NewTTSGate(prov, DefaultTTSGateConfig())

	const taskID = "task-394r1-hostpath"
	prov.OpenScope(taskID)
	const hostPath = `D:\tmp\wisp\20261101\394r1\artifact.txt`
	// The 175 shape the exemption is real for: the host embeds its own minted
	// path INTO the content (a task.output stub), so the declared span actually
	// locates inside the body — a body without the path exempts nothing.
	const body = "reread stub: " + hostPath + " (0 bytes) token-394R1HPQ"
	if !prov.MarkWithHostPath(taskID, SrcTaskOutput, "stub", body, hostPath) {
		t.Fatal("MarkWithHostPath refused the body")
	}

	if got := gate.CountInheritanceLinks(); got != 1 {
		t.Fatalf("seam recorded %d links, want 1", got)
	}
	if got := gate.inherited[0].TaintContent[0]; got != normalizeTaint(body) {
		t.Errorf("chain content = %q, want normalized body %q", got, normalizeTaint(body))
	}

	// The index half is unchanged: the declared path alone is no evidence from
	// this mark (177/183), while a foreign body window still hits (R4).
	if _, hit := prov.CheckText(taskID, ChTTS, hostPath); hit {
		t.Error("the declared host path became evidence again — the seam disturbed the 177/183 exemption")
	}
	if _, hit := prov.CheckText(taskID, ChTTS, "token-394R1HPQ"); !hit {
		t.Error("a body window stopped matching — the seam disturbed the index")
	}
}

// Test394SeamNotRecordedWhenMarkRecordsNothing: a mark that normalizes to
// nothing returns false and records nothing, so the seam must not conjure a
// link out of it.
func Test394SeamNotRecordedWhenMarkRecordsNothing(t *testing.T) {
	prov := NewProvenance(ProvOptions{NoProbe: true})
	gate := NewTTSGate(prov, DefaultTTSGateConfig())

	prov.OpenScope("task-394empty")
	if prov.Mark("task-394empty", SrcFSRead, "o", "   \u3000\u200b  ") {
		t.Fatal("Mark accepted whitespace-only content")
	}
	if got := gate.CountInheritanceLinks(); got != 0 {
		t.Fatalf("seam recorded %d links from a refused mark, want 0", got)
	}
}

// Test394SeamDisabledConfigRegistersNothing: a gate built with inheritance off
// never registers (the chain would never be scanned — registration would be
// pure memory waste), and the direct-check half keeps working.
func Test394SeamDisabledConfigRegistersNothing(t *testing.T) {
	prov := NewProvenance(ProvOptions{NoProbe: true})
	gate := NewTTSGate(prov, TTSGateConfig{
		EnableTaintInheritance: false,
		MaxAncestorDepth:       16,
		EnableL2Confirmation:   true,
	})

	const taskID = "task-394off"
	_ = prov.OpenScope(taskID)
	prov.Mark(taskID, SrcFSRead, "o", seam394Raw)
	if got := gate.CountInheritanceLinks(); got != 0 {
		t.Fatalf("disabled gate recorded %d links, want 0", got)
	}
	if _, refused, _ := gate.CheckTextWithInheritance(taskID, ChTTS, seam394Quote); !refused {
		t.Error("L2 direct check regressed while inheritance is off")
	}
}

// Test394SeamConcurrentMarksStayConsistent hammers the notify→upsert path from
// many goroutines: one shared scope folds EVERY mark into one link (no lost or
// duplicated taint), each goroutine's own scope gets its own link. Same
// owner/recover shape as TestTTSGate_ConcurrentAccess (d22scan ban on bare
// go func).
func Test394SeamConcurrentMarksStayConsistent(t *testing.T) {
	prov := NewProvenance(ProvOptions{NoProbe: true})
	gate := NewTTSGate(prov, DefaultTTSGateConfig())

	const workers = 8
	const perWorker = 25
	prov.OpenScope("task-394shared")

	done := make(chan bool, workers)
	for w := 0; w < workers; w++ {
		go func(w int) {
			defer func() {
				if rec := recover(); rec != nil {
					t.Errorf("worker %d panicked: %v", w, rec)
				}
				done <- true
			}()
			own := "task-394w" + string(rune('a'+w))
			prov.OpenScope(own)
			for i := 0; i < perWorker; i++ {
				prov.Mark("task-394shared", SrcWebFetch, "o", "shared-taint-"+string(rune('a'+w))+"-"+string(rune('A'+i))+"-fragment")
				prov.Mark(own, SrcWebFetch, "o", "own-taint-fragment-"+string(rune('a'+w)))
			}
		}(w)
	}
	for w := 0; w < workers; w++ {
		<-done
	}

	var shared *taintInheritanceLink
	for _, l := range gate.inherited {
		if l.TaskScopeID == "task-394shared" {
			shared = l
		}
	}
	if shared == nil {
		t.Fatal("the shared scope got no link")
	}
	if got := len(shared.TaintContent); got != workers*perWorker {
		t.Errorf("shared link holds %d taints, want %d — marks were lost or duplicated", got, workers*perWorker)
	}
	if len(shared.ToolSources) != len(shared.TaintContent) {
		t.Errorf("parallel slices diverged: %d tools vs %d taints", len(shared.ToolSources), len(shared.TaintContent))
	}
	if got := gate.CountInheritanceLinks(); got != workers+1 {
		t.Errorf("chain holds %d links, want %d (one per own scope + shared)", got, workers+1)
	}
}
