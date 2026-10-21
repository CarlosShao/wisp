package tools

// 票 172 机器侧装配件：尺的机器侧（桥 + 探针名册）与文本侧"可指副本"的演示件。
// 尺本身（解析 / 对账 / 读数 / 红句）在 tier_consistency_172_ruler_test.go，
// 用例在 tier_consistency_172_cases_test.go。

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// tier172Bridge wires the C4 Builtin slot the way cmd/wisp/run.go does: the same
// three entry builders, delete_enabled on (D34 has a row for fs.delete).
func tier172Bridge(t *testing.T, allowed string) *Bridge {
	t.Helper()
	paths := NewPathCanonicalizer([]string{allowed}, nil)
	reg := NewRegistry()
	roster := BuiltinFSEntries(FSDeps{Paths: paths, DeleteEnabled: true})
	roster = append(roster, BuiltinTaskEntries(TaskDeps{Roster: NewTaskRoster(), Paths: paths})...)
	roster = append(roster, BuiltinSubagentEntries(SubagentDeps{})...)
	for _, e := range roster {
		if err := reg.Register(e); err != nil {
			t.Fatalf("register %s: %v", e.Tool.Name(), err)
		}
	}
	return New(Options{Registry: reg, Paths: paths, Logf: func(string, ...any) {}})
}

// tier172Calls is one canonical call per (tool, mode) the frozen table distinguishes,
// over a machine whose allowlist is one temp dir alone. The tier values are NOT
// written here: this is only the glue telling the machine which situation to judge.
func tier172Calls(t *testing.T) ([]tier172Call, *Bridge) {
	t.Helper()
	in, out := sealableTempDir124(t), sealableTempDir124(t)
	existing := writeUnder(t, in, "there.txt", "172\n")
	outside := writeUnder(t, out, "outside.txt", "172\n")
	if err := os.Mkdir(filepath.Join(in, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	slash := filepath.ToSlash
	return []tier172Call{
		{"fs.read", map[string]any{"path": slash(existing)}, "授权目录内的既有文件"},
		{"fs.read", map[string]any{"path": slash(outside)}, "授权目录之外的既有文件"},
		{"fs.list", map[string]any{"path": slash(filepath.Join(in, "sub"))}, "授权目录内的目录"},
		{"fs.list", map[string]any{"path": slash(out)}, "授权目录之外的目录"},
		{"fs.write", map[string]any{"path": slash(filepath.Join(in, "brand-new.txt"))}, "在授权目录里新建"},
		{"fs.write", map[string]any{"path": slash(existing)}, "覆盖授权目录里的既有文件"},
		{"fs.edit", map[string]any{"path": slash(existing)}, "替换授权目录里既有文件的一小段"},
		{"fs.move", map[string]any{"from": slash(existing), "to": slash(filepath.Join(in, "renamed.txt"))}, "同卷内重命名"},
		{"fs.trash", map[string]any{"path": slash(existing)}, "移入回收站"},
		{"fs.delete", map[string]any{"path": slash(existing)}, "永久删除"},
		{"task.output", map[string]any{"task_id": "172-no-such-row"}, "读一枚不存在的任务行"},
		{"task.cancel", map[string]any{"task_id": "172-no-such-row"}, "停一枚不存在的任务行"},
		{"task.spawn", map[string]any{"task": "172 probe"}, "派生一枚子代理"},
	}, tier172Bridge(t, in)
}

// tier172RepoRoot walks up to the tree that holds docs/PLAN.md; tests run with the
// package directory as CWD, so the frozen text is two levels up.
func tier172RepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "docs", "PLAN.md")); err == nil {
			return dir
		}
		if parent := filepath.Dir(dir); parent != dir {
			dir = parent
			continue
		}
		break
	}
	t.Fatalf("no docs/PLAN.md above %s", dir)
	return ""
}

// tier172TextAt is where the text side comes from: the tracked frozen text under
// `root`, or a copy when the demonstration seam is set. The copy is what makes the
// red shape demonstrable without touching the tracked file (票面 AC#1/AC#4).
func tier172TextAt(t *testing.T, root string) (string, string) {
	t.Helper()
	path := filepath.Join(root, "docs", "PLAN.md")
	if p := os.Getenv(tier172DocEnv); p != "" {
		path = p
	}
	return tier172Read(t, path), path
}

func tier172Read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

// flipTier172Row rewrites one row's RiskLevel cell in a COPY of the text. It refuses
// to return unless exactly one row changed and the parse agrees with the change, so a
// "red" further down cannot come from an edit that missed.
func flipTier172Row(t *testing.T, text, tool, from, to string) string {
	t.Helper()
	lines, hit := strings.Split(text, "\n"), 0
	for i, ln := range lines {
		if !strings.HasPrefix(ln, "|") {
			continue
		}
		cells := strings.Split(ln, "|")
		if len(cells) < 4 {
			continue
		}
		if names := tier172Names(cells[1]); len(names) != 1 || names[0] != tool {
			continue
		}
		if strings.TrimSpace(cells[3]) != "**"+from+"**" {
			continue
		}
		cells[3] = " **" + to + "** "
		lines[i] = strings.Join(cells, "|")
		hit++
	}
	if hit != 1 {
		t.Fatalf("flip %s %s->%s touched %d rows, want exactly 1", tool, from, to, hit)
	}
	out := strings.Join(lines, "\n")
	rows, err := parseTier172Table(out)
	if err != nil {
		t.Fatalf("flipped copy no longer parses: %v", err)
	}
	for _, r := range rows {
		if r.tool == tool && len(r.tiers) == 1 && r.tiers[0] == to {
			return out
		}
	}
	t.Fatalf("flipped copy does not read %s as [%s]", tool, to)
	return ""
}
