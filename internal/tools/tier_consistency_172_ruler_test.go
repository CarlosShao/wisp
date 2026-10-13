package tools

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// 票 172 尺件：冻结文本里写的档位 与 这台机器实际给的档位 之间的一致性尺。
//
// 洞（票 162-v3 §4② 现量，本票 2026-10-13 23:1x 自复跑）：把 PLAN.md D34 表里 fs.edit
// 的档位 L2 改成 L1 => 全仓零枚用例红（没有任何 Go 码在运行时读契约文本，测试硬编的是
// 码侧的 L2）。这份文件补上那半句，分母＝整张表，不是一枚工具（票 171 AC#2：只钉一枚＝
// 私货尺）。
//
// 两边取数：文本侧从**给定路径**读 D34 表（所以能指向副本，真件一字不动），RiskLevel 格
// 里每一枚档位记号收成该工具档位集合的一条；机器侧每枚落地工具走**桥自己的判定路径**
// （pathArgs -> factsFor -> assessorFor(...).Assess），即"这次调用机器实际给哪一档"：声明、
// 路径抽取、Facts 钩子全是生产路径那一套，本尺不另列期望表。判据：每枚工具，机器探针给出的
// 档位集合必须与文本档位集合**相等**（不是"包含"）。
//
// 摘掉哪一处、哪条因此不响（AC#2）：FSEditDecl 的 Declared 改 L1 => 机器剩 {L1}、文本
// {L2} => Test172FrozenTierMatchesTheTierTheGateGives 红并点名 fs.edit；把 fs.edit 从注册
// 表摘掉（或改名）=> 落地名册钉红；改真文本那一行 => 同一发红（副本演示见用例件）。
//
// 覆盖不到、如实登记（不假绿）：备注列里没写进 RiskLevel 格的升级条款（如 fs.move 的"跨盘
// -> L2"）、未落地工具的越界/升级腿 —— 只对账两边都拿得到的读数。
const tier172DocEnv = "WISP172_TIER_DOC"

// 现量（2026-10-13 本腿，两处冻结文本逐字一致）：表里 37 行 markdown，其中 3 行各含多枚
// 工具（reminder.* 三枚 / memory.* 两枚 / task.list+task.cancel）=> (工具, 档位) 条 41 条、
// 去重工具名 40 枚；机器侧两边都有读数的 9 枚。这几个数是分母（AC#3）：任何一处变了都是
// 契约变更，钉住它们是为了不让分母悄悄缩水。
const (
	tier172TierEntries = 41
	tier172DocNames    = 40
	tier172LandedNumb  = 9
)

var (
	// 两边都有读数、真被对账的 9 枚（fs family + task.output/task.cancel）。
	tier172Landed = []string{
		"fs.delete", "fs.edit", "fs.list", "fs.move", "fs.read",
		"fs.trash", "fs.write", "task.cancel", "task.output",
	}
	// 机器名册里有、两处冻结文本都没有它的行：今天只有 task.spawn（票 197）。
	tier172Undocumented = []string{"task.spawn"}
	tier172TierRe       = regexp.MustCompile(`L[012]`)
)

// tier172Row is one parsed table row: one tool name and the tiers its cell gave.
type tier172Row struct {
	tool  string
	tiers []string
}

// tier172Reading is one machine probe: what the gate gave for one canonical call.
type tier172Reading struct {
	level risk.Level
	why   string
}

type tier172Verdict struct {
	tool     string
	doc      []string
	machine  []tier172Reading
	unlanded bool
}

// tier172Call is one canonical call to judge, and why that call is canonical.
type tier172Call struct {
	tool   string
	params map[string]any
	why    string
}

// parseTier172Table pulls the D34 table out of frozen text: found by its header row,
// ended by the first line that is not a markdown row, so a sibling table (16.5.5's)
// cannot leak in.
func parseTier172Table(text string) ([]tier172Row, error) {
	var rows []tier172Row
	inside := false
	for _, ln := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
		if !inside {
			inside = strings.HasPrefix(ln, "| 工具 | 动作 | RiskLevel")
			continue
		}
		if !strings.HasPrefix(ln, "|") {
			break
		}
		if strings.HasPrefix(ln, "|---") {
			continue
		}
		cells := strings.Split(strings.Trim(ln, "|"), "|")
		if len(cells) < 3 {
			return nil, fmt.Errorf("row has %d cells: %q", len(cells), ln)
		}
		names, tiers := tier172Names(cells[0]), tier172TierRe.FindAllString(cells[2], -1)
		if len(names) == 0 || len(tiers) == 0 {
			return nil, fmt.Errorf("row carries no name or no tier: %q", ln)
		}
		if len(names) > 1 && len(names) != len(tiers) {
			return nil, fmt.Errorf("%d names vs %d tiers: %q", len(names), len(tiers), ln)
		}
		for i, n := range names {
			t := tiers // one tool, N modes ("L0（授权内）/ L2（越界）")
			if len(names) > 1 {
				t = tiers[i : i+1] // positional ("task.list / task.cancel" -> L0 / L1)
			}
			rows = append(rows, tier172Row{tool: n, tiers: t})
		}
	}
	switch {
	case !inside:
		return nil, errors.New("no D34 header row")
	case len(rows) == 0:
		return nil, errors.New("D34 table parsed to zero rows")
	}
	return rows, nil
}

// tier172Names splits one first cell into tool names. Markdown shorthand inside the
// cell ("reminder.set/list/cancel") is completed from the shared namespace, so the
// names come out spelled the way the registry spells them.
func tier172Names(cell string) []string {
	cell = strings.NewReplacer("`", "", "*", "", " ", "").Replace(cell)
	var out []string
	for _, part := range strings.Split(cell, "/") {
		if part == "" {
			continue
		}
		if !strings.Contains(part, ".") && len(out) > 0 {
			if ns, _, ok := strings.Cut(out[0], "."); ok {
				part = ns + "." + part
			}
		}
		out = append(out, part)
	}
	return out
}

// tier172Docs folds parsed rows into one tier set per tool.
func tier172Docs(rows []tier172Row) map[string]map[string]bool {
	out := map[string]map[string]bool{}
	for _, r := range rows {
		if out[r.tool] == nil {
			out[r.tool] = map[string]bool{}
		}
		for _, tier := range r.tiers {
			out[r.tool][tier] = true
		}
	}
	return out
}

// tier172Reconcile is the ruler proper: frozen text in, per-tool verdicts out. The
// machine side is the caller's bridge (机器侧装配件), so one text can be judged twice.
func tier172Reconcile(text string, b *Bridge, calls []tier172Call) ([]tier172Verdict, []string, error) {
	rows, err := parseTier172Table(text)
	if err != nil {
		return nil, nil, err
	}
	doc, machine := tier172Docs(rows), map[string][]tier172Reading{}
	for _, c := range calls {
		entry, ok := b.reg.Lookup(c.tool)
		if !ok {
			return nil, nil, fmt.Errorf("probe for %s but it is not registered", c.tool)
		}
		facts := b.factsFor(context.Background(), entry, c.params, pathArgs(c.params, entry.Decl.PathParams))
		level := b.assessorFor("").Assess(c.tool, c.params, facts).Level
		machine[c.tool] = append(machine[c.tool], tier172Reading{level, c.why})
	}
	var verdicts []tier172Verdict
	for _, tool := range slices.Sorted(maps.Keys(doc)) {
		v := tier172Verdict{tool: tool, doc: slices.Sorted(maps.Keys(doc[tool])), machine: machine[tool]}
		v.unlanded = len(v.machine) == 0
		verdicts = append(verdicts, v)
	}
	var undocumented []string
	for _, e := range b.reg.List() {
		if _, ok := doc[e.Tool.Name()]; !ok {
			undocumented = append(undocumented, e.Tool.Name())
		}
	}
	slices.Sort(undocumented)
	return verdicts, undocumented, nil
}
