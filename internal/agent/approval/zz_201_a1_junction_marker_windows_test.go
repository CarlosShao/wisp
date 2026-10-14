//go:build windows

package approval_test

// Probe 201-a1 (dispatch 2026-10-14): the junction card-face fragment risk,
// measured live. QUESTION: a junction-refused fs call raises an L2 card whose
// Prompt.Paths carry the marker "<raw path> (无法规范化: …(junction/symlink)…)";
// cmd/wisp's 长期 branch reads that same card and asks it for the widening rule
// (cmd/wisp/approval_always.go:84 → ReplyCard.WideningRule, replies.go:519),
// which derives the directory through dirOf (batch.go:91) — a string split at
// the LAST '/' or '\'. Is there a '/' inside the marker that the split mistakes
// for a path separator, so a "directory" with half the marker in it comes out?
//
// ANSWER, and what changed on 2026-10-14: YES, that is exactly what the probe
// measured (ok=true, "…(junction" as the "dir", the fragment passing both
// AddAllowedDir gates and landing in config.toml — readings, verbatim, in
// .scratch/wisp/probes/201/a1/logs/run3-full.txt). Ticket 318 then fixed it:
// WideningRule now REFUSES the whole card at its entry when any path carries
// the marker (replies.go:521), so leg 3 below flipped its assertions while
// legs 4–5 were removed (their subject cannot exist after the guard).
// Leg 315-r1 set the precedent for this flip (a zz probe becoming the standing
// regression nail for the production shape it reproduced); the pre-fix reading
// is preserved where it was measured, not rewritten here.
//
// THIS FILE IS A PROBE, NOT AN ACCEPTANCE. It fails LOUDLY (never skips) when a
// prerequisite is missing. Every artifact used is real: a real NTFS junction
// (mklink /J, no admin needed), the real C26 canonicalizer, the real
// tools.Bridge, the real approval Gate/Queue/Replies, the real WideningRule.
// On the pre-fix tree legs 3+ are RED by construction; on the fixed tree they
// are GREEN, which is why they run in the package's default run (no env gate —
// an env-gated t.Skip is fatal under tools/d22scan/runtests.sh, ticket 71 AC#3).

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// probeToolsGate is the tools.Gate the bridge routes into. It captures the
// Decision the human's card would be built from and waves BOTH routes through —
// the same adversary shape as internal/tools's bridge_junction_windows_test.go:
// a human who approved the card must still not reach the far side, so nothing
// below depends on a refusal at the gate.
type probeToolsGate struct {
	mu   sync.Mutex
	last tools.Decision
	got  bool
}

func (g *probeToolsGate) set(d tools.Decision) {
	g.mu.Lock()
	g.last, g.got = d, true
	g.mu.Unlock()
}

func (g *probeToolsGate) decision() tools.Decision {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.last
}

func (g *probeToolsGate) PendingWindow(_ context.Context, d tools.Decision) (tools.Answer, string) {
	g.set(d)
	return tools.AnswerAllow, "probe：L1 窗口一路放行（对抗者）"
}

func (g *probeToolsGate) PendingApproval(_ context.Context, d tools.Decision) (tools.Answer, string) {
	g.set(d)
	return tools.AnswerAllow, "probe：L2 卡片被放行（对抗者）"
}

// probeJunction creates a REAL NTFS junction. Same construction as the tools
// corpus's mkRealJunction: mklink /J needs no administrator token.
func probeJunction(t *testing.T, link, target string) {
	t.Helper()
	out, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput()
	if err != nil {
		t.Fatalf("前置条件缺失：造不出真 NTFS junction（mklink /J %s %s: %v: %s）",
			link, target, err, out)
	}
}

// probeBridge wires the real bridge over the real canonicalizer, fs roster and
// assessor, with the capture gate above.
func probeBridge(t *testing.T, allowed string, g *probeToolsGate) *tools.Bridge {
	t.Helper()
	paths := tools.NewPathCanonicalizer([]string{allowed}, nil)
	reg := tools.NewRegistry()
	for _, e := range tools.BuiltinFSEntries(tools.FSDeps{Paths: paths}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	return tools.New(tools.Options{
		Registry: reg, Paths: paths, Gate: g,
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true}),
		Logf:       func(string, ...any) {},
	})
}

// TestZZ201A1JunctionMarkerCardIsRefusedForWidening is the live reproduction of
// the dispatch's chain — junction refused → card face carries the marker → the
// 长期 branch's own derivation reads that card — plus, since ticket 318's fix,
// the negative reading on the same chain: WideningRule must return a refusal,
// not a fragment.
func TestZZ201A1JunctionMarkerCardIsRefusedForWidening(t *testing.T) {
	root := t.TempDir()
	allowed := filepath.Join(root, "allowed")
	outside := filepath.Join(root, "outside")
	for _, d := range []string{allowed, outside} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	secret := "junction-only-bytes-201a1"
	victim := filepath.Join(outside, "victim.txt")
	if err := os.WriteFile(victim, []byte(secret), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(allowed, "jn")
	probeJunction(t, link, outside)

	// Positive control: the junction really reaches the bytes.
	if got, err := os.ReadFile(filepath.Join(link, "victim.txt")); err != nil || string(got) != secret {
		t.Fatalf("junction 没有真的通向目标（%v），后面的读数就什么也没证明", err)
	}

	// ---- leg 1: the junction call, through the real bridge -----------------
	pg := &probeToolsGate{}
	b := probeBridge(t, allowed, pg)
	args, err := json.Marshal(map[string]any{"path": strings.ReplaceAll(filepath.Join(link, "victim.txt"), `\`, "/")})
	if err != nil {
		t.Fatal(err)
	}
	out, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "task-201a1", CorrelationID: "corr-201a1", CallID: "call-201a1",
		Name: "fs.read", Args: args,
	})
	if err != nil {
		t.Fatalf("拒绝必须是判定，不是宿主故障: %v", err)
	}
	if !out.IsError {
		t.Fatalf("桥穿过了 junction：%+v", out)
	}
	if got, _ := os.ReadFile(victim); string(got) != secret {
		t.Fatal("junction 对面的内容被改动了")
	}
	if !pg.got {
		t.Fatal("链断在第一段：bridge 从未把卡片路由给 gate")
	}
	d := pg.decision()
	t.Logf("READING-1 bridge 路由的 Decision: level=%s rules=%v refusal=%q", d.Level, d.RulesHit, d.Refusal)
	if len(d.Paths) != 1 {
		t.Fatalf("卡面 Paths = %q, want 一条", d.Paths)
	}
	marked := d.Paths[0]
	t.Logf("READING-2 卡面逐字 Paths[0] = %q", marked)
	if !strings.Contains(marked, "无法规范化") || !strings.Contains(marked, "(junction/symlink)") {
		t.Fatalf("卡面没有带标记串（链路前提不成立）: %q", marked)
	}

	// Same call in the OTHER separator spelling (backslashes, no fold needed).
	// The pre-fix reading showed the same fragment for this spelling too
	// (READING-2b in probes/201/a1/logs/run3-full.txt, ok=true); the fix must
	// refuse it for the same reason, from the same guard.
	bargs, err := json.Marshal(map[string]any{"path": filepath.Join(link, "victim.txt")})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "task-201a1-bs", CorrelationID: "corr-201a1-bs", CallID: "call-201a1-bs",
		Name: "fs.read", Args: bargs,
	}); err != nil {
		t.Fatalf("反斜杠拼法：拒绝必须是判定: %v", err)
	}
	bsCard := approval.ReplyCard{CorrelationID: "bs", Level: "L2", Tool: "fs.read", Paths: pg.decision().Paths}
	bsDir, _, bsOK := bsCard.WideningRule()
	t.Logf("READING-2b 反斜杠拼法: 卡面=%q dir=%q ok=%v", bsCard.Paths, bsDir, bsOK)
	if bsOK || bsDir != "" {
		t.Fatalf("反斜杠拼法：带标记卡面仍被当作放宽基础（修后两种拼法应同拒）: dir=%q ok=%v", bsDir, bsOK)
	}

	// ---- leg 2: the real gate raises the card, the ledger books it ---------
	ui := newFakeUI()
	g, _, _ := newGate(t, ui)
	ledger := approval.NewReplies()
	ledger.Attach(approval.HostBinding{Gate: g, NativeSource: "probe-native", PanelSource: "probe-panel"})
	revoke := g.AdmitTextTask(d.TaskID)
	defer revoke()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ch := runApproval(t, g, ctx, d) // the bridge's own Decision, verbatim
	p := ui.wait(t)
	ledger.Record(p) // == cmd/wisp/run.go:1386 `u.live.record(p)`
	if err := g.Native().Allow(ctx, p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("Native().Allow: %v", err)
	}
	if a := mustAnswer(t, ch); a.a != tools.AnswerAllow {
		t.Fatalf("第一张卡的答复 = %v", a)
	}
	card, ok := ledger.Look(p.CorrelationID)
	if !ok {
		t.Fatal("账上没有这张卡")
	}
	if len(card.Paths) != 1 || card.Paths[0] != marked {
		t.Fatalf("ledger 里的 Paths 与卡面不一致: %q", card.Paths)
	}
	t.Logf("READING-3 ledger 卡: level=%s tool=%s paths=%q", card.Level, card.Tool, card.Paths)

	// ---- leg 3: THE MEASUREMENT — the 长期 branch's own derivation ---------
	// Ticket 318 flipped this leg on 2026-10-14. The pre-fix reading (verbatim:
	// probes/201/a1/logs/run3-full.txt READING-4/5/6) was ok=true with dir =
	// "…victim.txt (无法规范化: …(junction", the half-marker fragment. Post-fix
	// WideningRule refuses the whole card at its entry, so the chain stops here:
	// cmd/wisp's 长期 branch takes its ok=false branch (approval_always.go:85-91),
	// raises no second card and never calls AddAllowedDir — which is why the
	// pre-fix legs 4-5 (fragment reaching config.toml, second card carrying the
	// same fragment) have no post-fix counterpart left to run: their subject
	// cannot exist. Those readings stay where they were measured, not rewritten.
	dir, rule, ok := card.WideningRule()
	t.Logf("READING-4 WideningRule: ok=%v", ok)
	t.Logf("READING-5 dir  = %q", dir)
	t.Logf("READING-6 rule = %q", rule)
	if ok {
		t.Fatalf("修后带标记卡面仍产出放宽目标（碎片没有死在入口）: dir=%q rule=%q", dir, rule)
	}
	if dir != "" || rule != "" {
		t.Fatalf("拒绝必须同时清空 dir/rule，got (%q,%q)", dir, rule)
	}
	trueParent := strings.ReplaceAll(link, `\`, "/")
	t.Logf("READING-7 真父目录 = %q（修后不再被产出到任何卡/规则文本）", trueParent)

	// ---- negative control: the same derivation on a clean card -------------
	// Ticket 318 kept this control and sharpened it: a real out-of-scope read
	// whose path CAN be canonicalized carries no marker, so it must still derive
	// its true parent - if the new guard ever fired here it would be refusing
	// clean cards, which is the regression direction the fix must not have.
	clean := filepath.Join(outside, "victim.txt")
	pg2 := &probeToolsGate{}
	b2 := probeBridge(t, allowed, pg2)
	cargs, err := json.Marshal(map[string]any{"path": strings.ReplaceAll(clean, `\`, "/")})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := b2.Execute(context.Background(), agent.ToolRequest{
		TaskID: "task-201a1-control", CorrelationID: "corr-201a1-control", CallID: "call-201a1-control",
		Name: "fs.read", Args: cargs,
	}); err != nil {
		t.Fatalf("对照：拒绝必须是判定: %v", err)
	}
	cd := pg2.decision()
	ctrlCard := approval.ReplyCard{CorrelationID: "ctrl", Level: "L2", Tool: cd.Tool, Paths: cd.Paths}
	cdir, crule, cok := ctrlCard.WideningRule()
	t.Logf("READING-11 对照（无标记的越界读）: paths=%q dir=%q rule=%q ok=%v", cd.Paths, cdir, crule, cok)
	if !cok {
		t.Fatalf("对照（无标记）不再推出放宽目标，说明守卫误伤干净卡面: dir=%q ok=%v", cdir, cok)
	}
	if !strings.EqualFold(cdir, filepath.ToSlash(outside)) {
		t.Fatalf("对照推出的不是真父目录: dir=%q want %q", cdir, filepath.ToSlash(outside))
	}
	if strings.Contains(cdir, "无法规范化") {
		t.Fatalf("对照项也带标记串，说明夹具被污染: %q", cdir)
	}
	t.Logf("READING-12 结论：带标记卡面在入口整卡拒（两种拼法 + 空路径标记同拒），干净卡面仍推真父目录 %q", cdir)
}
