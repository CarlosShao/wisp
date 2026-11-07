package main

// 票 330-r1（G1 对接钉）：plan-present.sse 线上那枚工具名 ⇄ 装配出来的真名册必须逐字相等。
//
// 缺口（docs/evidence/s1/165-v2-verdict.md G1，台账 A1009）：internal/agent 的门铃用例
// 听的是自己抄下来的名册（planPresentName165r3 一枚常量喂四处），与 golden 夹具同步改名
// 即自洽——「镜像⇄真名册」之间零枚用例对接。机理＝导入环（internal/agent 只能包内形状、
// ⛔ import internal/tools，bridge.go:13），真名册那一半只在 cmd/wisp 够得着，所以钉落在这面。
//
// 判据的两半都⛔不落第二份硬编码字面量（AC#2）：
//   - 线上名：现读 internal/agent/testdata/golden/plan-present.sse 的字节，从
//     "function":{"name":...} 结构里抽出来（wireNamesFromSSE330r1）；
//   - 真名册：newReplyHost 走 runTextTask -> assembleRuntime 的真装配，读
//     rt.bridge.Registry()（cmd/wisp/run.go 的 BuiltinPlanEntries 注册进的那一个）。
//
// 两半在测试进程里现读现比，本文件不含任何工具名字面量。夹具名与镜像名一起被改成
// plan.presentX 而真名册不动的那一发必须红（G1 重放，读数在 docs/evidence/s1/330-r1-dock.md）。
//
// PATH warning（票 98）：本包测试二进制链 sherpa-onnx，PATH 缺 third_party/sherpa-onnx
// 时加载即死（0xc0000135），不是缺陷。

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"
)

// wireNamesFromSSE330r1 reads the golden fixture and returns every tool-call
// name carried on the wire, in order of first appearance. The parse is
// structural on purpose: the "name" string of any "function" object counts and
// nothing else does, so a renamed fixture yields the renamed name verbatim
// (the G1 replay depends on exactly that) and no tool-name literal lives in
// this file.
func wireNamesFromSSE330r1(t *testing.T, path string) []string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("读不到 golden 夹具 %s：%v", path, err)
	}
	var names []string
	var walk func(node any)
	walk = func(node any) {
		switch v := node.(type) {
		case map[string]any:
			if fn, ok := v["function"].(map[string]any); ok {
				if name, ok := fn["name"].(string); ok && name != "" && !slices.Contains(names, name) {
					names = append(names, name)
				}
			}
			for _, child := range v {
				walk(child)
			}
		case []any:
			for _, child := range v {
				walk(child)
			}
		}
	}
	for _, line := range strings.Split(string(raw), "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(trimmed, "data:"))
		if payload == "" || payload == "[DONE]" {
			continue
		}
		var chunk any
		if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
			t.Fatalf("golden 夹具 %s 有一行不是 JSON：%v（%q）", path, err, payload)
		}
		walk(chunk)
	}
	return names
}

// Test330DockPinsTheWireNameToTheAssembledRoster is the G1 dock pin: the one
// name the plan-present wire carries must be a row of the roster the real
// assembly registered. The positive firing (fixture and roster untouched)
// answers green because the two sides name the same tool; the G1 replay
// (fixture AND mirror renamed to plan.presentX, real roster untouched) must
// answer red on the judgment below, because the renamed wire name is no longer
// a row of anything the assembly registered.
func Test330DockPinsTheWireNameToTheAssembledRoster(t *testing.T) {
	// 线上名：从仓内 golden 夹具现读。⛔ 第二份字面量：名字本身只存在于夹具字节里。
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller 拿不到本文件路径，golden 夹具定位不了")
	}
	wirePath := filepath.Join(filepath.Dir(thisFile), "..", "..", "internal", "agent", "testdata", "golden", "plan-present.sse")
	wire := wireNamesFromSSE330r1(t, wirePath)
	if len(wire) == 0 {
		t.Fatal("golden 夹具里一枚工具名都抽不出来：仪器空转，对接钉无从谈起")
	}
	if len(wire) != 1 {
		t.Fatalf("golden 夹具线上名 %d 枚（%v），该夹具恰好载一枚 call；多枚时对接钉须回票面改形，不在这里顺手放宽", len(wire), wire)
	}

	// 真名册：走真装配，读 run.go 注册进 bridge 的那一个 Registry，不重建第二份名册。
	h := newReplyHost(t, 40*time.Second)
	var (
		roster   []string
		entryOK  bool
		resident bool
	)
	h.rtHook = func(rt *agentRuntime) {
		for _, e := range rt.bridge.Registry().List() {
			roster = append(roster, e.Tool.Name())
		}
		entry, ok := rt.bridge.Registry().Lookup(wire[0])
		entryOK = ok
		resident = ok && entry.Decl.Resident
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("装配跑不起来 exit=%d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	if len(roster) == 0 {
		t.Fatal("装配出来的真名册是空的：run.go 一枚工具都没注册上，对接钉必红（此时先修装配，不要放宽判据）")
	}

	// THE JUDGMENT：线上名必须逐字等于真名册里的某一行。
	if !entryOK || !slices.Contains(roster, wire[0]) {
		t.Fatalf("G1 对接钉：plan-present.sse 线上名 %q 不在装配出来的真名册（%d 枚）里：%v"+
			"——镜像⇄真名册失配，正是 165-v2 G1 报的那发绿要变红的形状", wire[0], len(roster), roster)
	}
	if !resident {
		t.Errorf("真名册里那枚 %q 的 Resident=false，冻结行（PLAN.md D34 / internal/tools/plan_present.go:129）写的是常驻", wire[0])
	}
	t.Logf("对接读数：线上名 %q ＝ 真名册（%d 枚）里的一行；Resident=%v", wire[0], len(roster), resident)
}
