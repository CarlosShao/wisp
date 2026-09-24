package risk

import (
	"reflect"
	"strings"
	"testing"
)

// --- R9 fail-closed ----------------------------------------------------------

func TestR9PanicInSubAssessorFailsClosed(t *testing.T) {
	a := NewRiskAssessor()
	a.RegisterSubAssessor(panicSubAssessor{})

	d := a.Assess("fs.read", nil, Facts{})
	if d.Level != L2 {
		t.Fatalf("panic sub-assessor must fail-closed to L2, got %v", d.Level)
	}
	if len(d.RulesHit) != 1 || d.RulesHit[0] != R9 {
		t.Fatalf("want [R9], got %v", d.RulesHit)
	}
	if !strings.Contains(d.Reason, "panic") {
		t.Fatalf("reason must name the panic, got %q", d.Reason)
	}
}

func TestR9PanicInWiredDependencyFailsClosed(t *testing.T) {
	a := NewRiskAssessor().WithCanonicalizer(&stubCanon{panicOn: "D:/proj/boom.txt"})

	d := a.Assess("fs.write", nil, Facts{Paths: []string{"D:/proj/boom.txt"}})
	if d.Level != L2 || len(d.RulesHit) != 1 || d.RulesHit[0] != R9 {
		t.Fatalf("panicking dependency: want L2/[R9], got %v/%v (%s)", d.Level, d.RulesHit, d.Reason)
	}
}

// --- fusion -------------------------------------------------------------------

func TestFusionTakesMaxSeverity(t *testing.T) {
	a := NewRiskAssessor()

	// Declared L0 + R8 send class: computed rule overrides a declared L0.
	d := a.Assess("mail.send", nil, Facts{Irreversible: []string{"send"}})
	if d.Level != L2 || len(d.RulesHit) != 1 || d.RulesHit[0] != R8 {
		t.Fatalf("declared L0 + send: want L2/[R8], got %v/%v", d.Level, d.RulesHit)
	}

	// Declared L2 + R1 floor + R5 hit: fused at L2, rules sorted R1 R5.
	d = a.Assess("web.open", nil, Facts{
		Declared: L2,
		Network:  &NetTarget{Protocol: "https", Host: "10.0.0.1"},
	})
	if d.Level != L2 || !reflect.DeepEqual(d.RulesHit, []RuleID{R1, R5}) {
		t.Fatalf("declared L2 + private range: want L2/[R1 R5], got %v/%v", d.Level, d.RulesHit)
	}

	// L1 floor + L2 rule: max wins.
	d = a.Assess("shell.exec", nil, Facts{Declared: L1, ShellArgv: []string{"git", "status"}, ShellAllowlist: []string{"git"}})
	if d.Level != L1 || !reflect.DeepEqual(d.RulesHit, []RuleID{R1, R6}) {
		t.Fatalf("declared L1 + allowlisted shell: want L1/[R1 R6], got %v/%v", d.Level, d.RulesHit)
	}
}

func TestFusionSingleDenyBeatsAnyConfirm(t *testing.T) {
	a := NewRiskAssessor().
		WithCanonicalizer(&stubCanon{allow: map[string]bool{"~/.git-credentials": true}}).
		WithSensitiveClassifier(&stubClassifier{tiers: map[string]Tier{"~/.git-credentials": TierA}})

	// Tier-A Deny + batch scale L2: Deny is terminal.
	d := a.Assess("fs.delete", nil, Facts{
		Paths:      []string{"~/.git-credentials"},
		BatchCount: 5000,
	})
	if d.Level != Deny || !reflect.DeepEqual(d.RulesHit, []RuleID{R3, R7}) {
		t.Fatalf("Deny must beat L2: want Deny/[R3 R7], got %v/%v", d.Level, d.RulesHit)
	}
}

func TestFusionR4BlocksSessionOverride(t *testing.T) {
	a := NewRiskAssessor().
		WithCanonicalizer(&stubCanon{}). // allowlist empty: every path is out of bounds
		WithTaintDetector(&stubTaint{src: "web.fetch", hitOn: "leak"})

	d := a.Assess("notify", map[string]any{"text": "leak"}, Facts{Paths: []string{"D:/proj/x"}})
	if d.Level != L2 || !reflect.DeepEqual(d.RulesHit, []RuleID{R2, R4}) {
		t.Fatalf("want L2/[R2 R4], got %v/%v", d.Level, d.RulesHit)
	}
	if !d.SessionOverrideBlocked {
		t.Fatal("R4 hit must block D45 session override")
	}
}

// --- golden snapshots ---------------------------------------------------------

// TestDecisionGoldenSnapshots pins the exact Decision objects (level, rules
// order, card-renderable reason) for representative calls. The confirmation
// card (21/37) renders rules_hit and reason verbatim; any diff here is a
// contract change and needs human review.
func TestDecisionGoldenSnapshots(t *testing.T) {
	wired := NewRiskAssessor().
		WithCanonicalizer(&stubCanon{allow: map[string]bool{
			"D:/proj/ok.txt": true,
			"D:/proj/.env":   true,
		}}).
		WithSensitiveClassifier(&stubClassifier{tiers: map[string]Tier{"D:/proj/.env": TierB}}).
		WithTaintDetector(&stubTaint{src: "web.fetch"})

	cases := []struct {
		name  string
		a     *RiskAssessor
		tool  string
		facts Facts
		want  Decision
	}{
		{
			name: "clean L0 call",
			a:    wired, tool: "fs.read",
			facts: Facts{},
			want:  Decision{Level: L0, RulesHit: []RuleID{}, Reason: "无规则命中（L0 直接执行）"},
		},
		{
			name: "R1 declared floor",
			a:    wired, tool: "fs.write",
			facts: Facts{Declared: L1},
			want:  Decision{Level: L1, RulesHit: []RuleID{R1}, Reason: "R1: 工具声明为下界（L1）"},
		},
		{
			name: "R2 out of allowlist",
			a:    wired, tool: "fs.write",
			facts: Facts{Paths: []string{"D:/outside/x.txt"}},
			want:  Decision{Level: L2, RulesHit: []RuleID{R2}, Reason: "R2: 目标路径在授权目录之外: D:/outside/x.txt"},
		},
		{
			name: "R3 tier B",
			a:    wired, tool: "fs.read",
			facts: Facts{Paths: []string{"D:/proj/.env"}},
			want:  Decision{Level: L2, RulesHit: []RuleID{R3}, Reason: "R3: 命中 B 档敏感路径（可单文件豁免）: D:/proj/.env"},
		},
		{
			name: "R5 deny-listed protocol",
			a:    wired, tool: "web.open",
			facts: Facts{Network: &NetTarget{Protocol: "vbscript"}},
			want:  Decision{Level: Deny, RulesHit: []RuleID{R5}, Reason: "R5: 协议被明确拒绝: vbscript"},
		},
		{
			name: "R6 metacharacters",
			a:    wired, tool: "shell.exec",
			facts: Facts{ShellArgv: []string{"cmd", "/c", "a & b"}},
			want:  Decision{Level: L2, RulesHit: []RuleID{R6}, Reason: "R6: argv 含元字符（&）"},
		},
		{
			name: "R7 at threshold",
			a:    wired, tool: "fs.delete",
			facts: Facts{BatchCount: 50},
			want:  Decision{Level: L2, RulesHit: []RuleID{R7}, Reason: "R7: 单次调用影响 50 个文件（>=50）"},
		},
		{
			name: "R8 send class",
			a:    wired, tool: "mail.send",
			facts: Facts{Irreversible: []string{"send"}},
			want:  Decision{Level: L2, RulesHit: []RuleID{R8}, Reason: "R8: 不可逆操作（发送类外发操作）"},
		},
		{
			name: "R9 panicking sub-assessor",
			a:    func() *RiskAssessor { x := NewRiskAssessor(); x.RegisterSubAssessor(panicSubAssessor{}); return x }(),
			tool: "fs.read",
			want: Decision{Level: L2, RulesHit: []RuleID{R9}, Reason: "R9: 判定器 panic，fail-closed 升 L2（插件判定器 panic: judge exploded）"},
		},
		{
			name: "R2 fail-closed on resolver error",
			a:    NewRiskAssessor().WithCanonicalizer(&stubCanon{errOn: "C:/broken"}), tool: "fs.write",
			facts: Facts{Paths: []string{"C:/broken"}},
			want:  Decision{Level: L2, RulesHit: []RuleID{R2}, Reason: `R2: 路径无法规范化，按越界处理（fail-closed）: "C:/broken"`},
		},
	}

	for _, tc := range cases {
		got := tc.a.Assess(tc.tool, nil, tc.facts)
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("%s:\n  want %+v\n  got  %+v", tc.name, tc.want, got)
		}
	}
}
