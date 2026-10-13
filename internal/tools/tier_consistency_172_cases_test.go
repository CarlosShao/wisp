package tools

// 票 172 用例件：三发用例。尺在 tier_consistency_172_ruler_test.go，机器侧装配与
// 文本侧"可指副本"演示件在 tier_consistency_172_machine_test.go。
//
// 每发答"未修码上响不响"：Test172DocCopyThatLies* 的对照段是**未修**读数（本票 23:1x
// 在真文本上自复跑过：改真文本那行 => 零枚红），另两发在未修码上绿、在变异上红。

import (
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// tier172Mismatched names the tools whose two sides disagree; unlanded rows have no
// machine side to compare, so they are not counted here.
func tier172Mismatched(verdicts []tier172Verdict) []string {
	var out []string
	for _, v := range verdicts {
		if !v.unlanded && !slices.Equal(v.doc, tier172LevelsOf(v.machine)) {
			out = append(out, v.tool)
		}
	}
	return out
}

// tier172Report prints every row's reading and reddens on a mismatch, naming the tool
// and both sides the way AC#1/AC#2 ask.
func tier172Report(t *testing.T, verdicts []tier172Verdict, undocumented []string) {
	t.Helper()
	unlanded := 0
	for _, v := range verdicts {
		got := tier172LevelsOf(v.machine)
		switch {
		case v.unlanded:
			unlanded++
			t.Logf("[未落地] %-18s 文本 %v 机器 无名册（未注册为 C4 工具）", v.tool, v.doc)
		case slices.Equal(v.doc, got):
			t.Logf("[对得上] %-18s 文本 %v 机器 %v", v.tool, v.doc, got)
		default:
			var probes []string
			for _, r := range v.machine {
				probes = append(probes, r.why+" -> "+r.level.String())
			}
			t.Errorf("%s：文档写 %s、机器给 %s（探针：%s）—— 冻结文本与门给的档位不一致，"+
				"两边必须由同一个 ticket 改齐", v.tool, strings.Join(v.doc, "/"),
				strings.Join(got, "/"), strings.Join(probes, "；"))
		}
	}
	bad := tier172Mismatched(verdicts)
	t.Logf("分母 %d 枚工具：对得上 %d，对不上 %d，未落地 %d；机器有而契约无 %v",
		len(verdicts), len(verdicts)-unlanded-len(bad), len(bad), unlanded, undocumented)
}

func tier172LevelsOf(readings []tier172Reading) []string {
	seen := map[string]bool{}
	var out []string
	for _, r := range readings {
		if s := r.level.String(); !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	slices.Sort(out)
	return out
}

// tier172TiersIn parses one frozen text and pins its (tool, tier) entry count as the
// denominator.
func tier172TiersIn(t *testing.T, what, text string) map[string]map[string]bool {
	t.Helper()
	rows, err := parseTier172Table(text)
	if err != nil {
		t.Fatalf("%s: %v", what, err)
	}
	if len(rows) != tier172TierEntries {
		t.Errorf("%s 的 (工具, 档位) 条数 %d，现量钉住的是 %d —— 这张表是冻结契约，条数变了"+
			"一枚就该有人知道", what, len(rows), tier172TierEntries)
	}
	return tier172Docs(rows)
}

// Test172RulerReadsBothFrozenTablesAndTheyMirror is the ruler's own parachute: the
// parser must really find the table (an empty parse would make every later check
// green), and the two frozen texts must still tell the same tier per tool.
func Test172RulerReadsBothFrozenTablesAndTheyMirror(t *testing.T) {
	root := tier172RepoRoot(t)
	plan := tier172TiersIn(t, "PLAN.md D34",
		tier172Read(t, filepath.Join(root, "docs", "PLAN.md")))
	spec := tier172TiersIn(t, "SPEC-07 §3",
		tier172Read(t, filepath.Join(root, "docs", "specs", "SPEC-07-tools-and-plugins.md")))
	if got := slices.Sorted(maps.Keys(plan["fs.edit"])); !slices.Equal(got, []string{"L2"}) {
		t.Fatalf("parser positive control: PLAN.md 的 fs.edit 行读出 %v，要的是 [L2]", got)
	}
	if len(plan) != tier172DocNames {
		t.Errorf("D34 去重工具名 %d 枚，现量钉住的是 %d 枚", len(plan), tier172DocNames)
	}
	for _, tool := range slices.Sorted(maps.Keys(plan)) {
		if _, ok := spec[tool]; !ok {
			t.Errorf("%s 在 PLAN.md D34 有行，SPEC-07 §3 镜像没有", tool)
			continue
		}
		want, got := slices.Sorted(maps.Keys(plan[tool])), slices.Sorted(maps.Keys(spec[tool]))
		if !slices.Equal(want, got) {
			t.Errorf("%s：PLAN.md D34 写 %v，SPEC-07 §3 镜像写 %v —— 两处冻结文本自己先不一致",
				tool, want, got)
		}
	}
}

// Test172FrozenTierMatchesTheTierTheGateGives is the live ruler over the real frozen
// text (or over the copy the demonstration seam points it at).
func Test172FrozenTierMatchesTheTierTheGateGives(t *testing.T) {
	root := tier172RepoRoot(t)
	text, docPath := tier172TextAt(t, root)
	if os.Getenv(tier172DocEnv) == "" && docPath != filepath.Join(root, "docs", "PLAN.md") {
		t.Fatalf("default text side must be the tracked docs/PLAN.md, got %s", docPath)
	}
	calls, b := tier172Calls(t)
	verdicts, undocumented, err := tier172Reconcile(text, b, calls)
	if err != nil {
		t.Fatalf("ruler: %v", err)
	}
	tier172Report(t, verdicts, undocumented)
	if !slices.Equal(tier172Undocumented, undocumented) {
		t.Errorf("机器名册里出现了没有契约行的工具 %v，今天是 %v：新工具要么带着它的契约行"+
			"（与机器同一票落），要么在这里被显式登记", undocumented, tier172Undocumented)
	}
	var landed []string
	for _, v := range verdicts {
		if !v.unlanded {
			landed = append(landed, v.tool)
		}
	}
	if len(landed) != tier172LandedNumb || !slices.Equal(landed, tier172Landed) {
		t.Errorf("落地（两边都有读数）的工具从 %v 变成 %v：少一枚就等于有一行契约文本不再"+
			"被对账，多一枚就得多一行分母", tier172Landed, landed)
	}
}

// Test172DocCopyThatLiesReddensTheRuler is AC#1's tooth: the SAME ruler, pointed at a
// copy whose fs.edit row lies, must red and name fs.edit, while the real frozen text
// is never written to.
func Test172DocCopyThatLiesReddensTheRuler(t *testing.T) {
	root := tier172RepoRoot(t)
	real := tier172Read(t, filepath.Join(root, "docs", "PLAN.md"))
	flipped := flipTier172Row(t, real, "fs.edit", "L2", "L1")
	if err := os.WriteFile(filepath.Join(t.TempDir(), "PLAN-172-copy.md"),
		[]byte(flipped), 0o600); err != nil {
		t.Fatal(err)
	}
	calls, b := tier172Calls(t)
	control, _, err := tier172Reconcile(real, b, calls)
	if err != nil {
		t.Fatalf("ruler: %v", err)
	}
	if bad := tier172Mismatched(control); len(bad) != 0 {
		t.Fatalf("对照组：尺在没改过的真文本上就红了（%v），红形演示无从归属", bad)
	}
	verdicts, _, err := tier172Reconcile(flipped, b, calls)
	if err != nil {
		t.Fatalf("ruler over the copy: %v", err)
	}
	if bad := tier172Mismatched(verdicts); !slices.Equal(bad, []string{"fs.edit"}) {
		t.Errorf("副本里只把 fs.edit 那行改成 L1，尺红在这些工具上 %v，要的正好是 [fs.edit]", bad)
	}
	for _, v := range verdicts {
		if v.tool != "fs.edit" {
			continue
		}
		got := strings.Join(v.doc, "/") + "|" + strings.Join(tier172LevelsOf(v.machine), "/")
		if got != "L1|L2" {
			t.Errorf("副本里 fs.edit 那行改成 L1 之后，尺给的读数是 %q，要的是文本 L1 / 机器 L2"+
				" —— 尺必须喊出'文档写 L1、机器给 L2'", got)
		}
	}
}
