package risk

import (
	"fmt"
	"strings"
	"testing"
)

// --- test doubles -----------------------------------------------------------

// stubCanon is a fake C26 resolver: identity canonicalization, allowlist and
// tiers keyed on the raw spelling. errOn/panicOn inject failure modes.
type stubCanon struct {
	allow   map[string]bool
	errOn   string
	panicOn string
}

func (s *stubCanon) Canonicalize(raw string) (string, error) {
	if s.panicOn == raw {
		panic("canonicalize panic")
	}
	if s.errOn == raw {
		return "", fmt.Errorf("unresolvable: %s", raw)
	}
	return raw, nil
}

func (s *stubCanon) InAllowlist(canonical string) bool { return s.allow[canonical] }

type stubClassifier struct{ tiers map[string]Tier }

func (s *stubClassifier) Classify(canonical string) Tier { return s.tiers[canonical] }

type stubTaint struct {
	src   string
	hitOn string // hits when params["text"] equals this value
	seen  map[string]any
	asked int
}

func (s *stubTaint) TaintHit(params map[string]any) (string, bool) {
	s.seen = params
	s.asked++
	text, _ := params["text"].(string)
	return s.src, s.hitOn != "" && text == s.hitOn
}

type panicSubAssessor struct{}

func (panicSubAssessor) Assess(string, map[string]any, Facts) Decision { panic("judge exploded") }

// --- per-rule matrix: one positive + one negative + edge per rule ------------

func TestRuleDeclaredR1(t *testing.T) {
	a := NewRiskAssessor()

	d := a.Assess("any.tool", nil, Facts{Declared: L2}) // positive: floor raised
	if d.Level != L2 || len(d.RulesHit) != 1 || d.RulesHit[0] != R1 {
		t.Fatalf("declared L2: want L2/[R1], got %v/%v", d.Level, d.RulesHit)
	}

	d = a.Assess("any.tool", nil, Facts{}) // negative: L0 declaration never hits R1
	if d.Level != L0 || len(d.RulesHit) != 0 {
		t.Fatalf("declared L0: want L0/no hits, got %v/%v", d.Level, d.RulesHit)
	}

	d = a.Assess("any.tool", nil, Facts{Declared: Level(9)}) // edge: invalid input fail-closed
	if d.Level != L2 || len(d.RulesHit) != 1 || d.RulesHit[0] != R9 {
		t.Fatalf("invalid declared: want L2/[R9], got %v/%v", d.Level, d.RulesHit)
	}
}

func TestRulePathAllowlistR2(t *testing.T) {
	a := NewRiskAssessor().WithCanonicalizer(&stubCanon{allow: map[string]bool{"D:/proj/a.txt": true}})

	d := a.Assess("fs.write", nil, Facts{Paths: []string{"D:/other/evil.txt"}}) // positive
	if d.Level != L2 || len(d.RulesHit) != 1 || d.RulesHit[0] != R2 {
		t.Fatalf("outside allowlist: want L2/[R2], got %v/%v", d.Level, d.RulesHit)
	}

	d = a.Assess("fs.write", nil, Facts{Paths: []string{"D:/proj/a.txt"}}) // negative
	if d.Level != L0 || len(d.RulesHit) != 0 {
		t.Fatalf("inside allowlist: want L0/no hits, got %v/%v", d.Level, d.RulesHit)
	}

	a2 := NewRiskAssessor().WithCanonicalizer(&stubCanon{errOn: "C:/broken/link"}) // edge: resolver error fail-closed
	d = a2.Assess("fs.write", nil, Facts{Paths: []string{"C:/broken/link"}})
	if d.Level != L2 || len(d.RulesHit) != 1 || d.RulesHit[0] != R2 {
		t.Fatalf("unresolvable path: want L2/[R2], got %v/%v", d.Level, d.RulesHit)
	}

	d = a.Assess("fs.write", nil, Facts{}) // dormant without paths
	if len(d.RulesHit) != 0 {
		t.Fatalf("no paths: want no hits, got %v", d.RulesHit)
	}
}

func TestRuleSensitivePathR3(t *testing.T) {
	a := NewRiskAssessor().
		WithCanonicalizer(&stubCanon{allow: map[string]bool{
			"~/.git-credentials": true, // stub: allowlist isolation for R3 only
			"D:/proj/.env":       true,
			"D:/proj/notes.txt":  true,
		}}).
		WithSensitiveClassifier(&stubClassifier{tiers: map[string]Tier{
			"~/.git-credentials": TierA,
			"D:/proj/.env":       TierB,
		}})

	d := a.Assess("fs.read", nil, Facts{Paths: []string{"~/.git-credentials"}}) // positive A
	if d.Level != Deny || len(d.RulesHit) != 1 || d.RulesHit[0] != R3 {
		t.Fatalf("tier A: want Deny/[R3], got %v/%v", d.Level, d.RulesHit)
	}

	d = a.Assess("fs.read", nil, Facts{Paths: []string{"D:/proj/notes.txt"}}) // negative
	if d.Level != L0 || len(d.RulesHit) != 0 {
		t.Fatalf("clean path: want L0/no hits, got %v/%v", d.Level, d.RulesHit)
	}

	d = a.Assess("fs.read", nil, Facts{Paths: []string{"D:/proj/.env"}}) // edge B
	if d.Level != L2 || len(d.RulesHit) != 1 || d.RulesHit[0] != R3 {
		t.Fatalf("tier B: want L2/[R3], got %v/%v", d.Level, d.RulesHit)
	}
}

func TestRuleTaintR4(t *testing.T) {
	det := &stubTaint{src: "web.fetch", hitOn: "leak"}
	a := NewRiskAssessor().WithTaintDetector(det)
	params := map[string]any{"text": "leak"}

	d := a.Assess("notify", params, Facts{}) // positive
	if d.Level != L2 || len(d.RulesHit) != 1 || d.RulesHit[0] != R4 || !d.SessionOverrideBlocked {
		t.Fatalf("taint hit: want L2/[R4]/blocked, got %v/%v/%v", d.Level, d.RulesHit, d.SessionOverrideBlocked)
	}
	if det.seen["text"] != "leak" || det.asked != 1 {
		t.Fatalf("taint detector must receive params, got %v (asked=%d)", det.seen, det.asked)
	}

	d = a.Assess("notify", map[string]any{"text": "hello"}, Facts{}) // negative
	if d.Level != L0 || len(d.RulesHit) != 0 || d.SessionOverrideBlocked {
		t.Fatalf("no taint: want L0/no hits/unblocked, got %v/%v/%v", d.Level, d.RulesHit, d.SessionOverrideBlocked)
	}
	if det.asked != 2 {
		t.Fatalf("taint detector must run on every call, asked=%d", det.asked)
	}
}

func TestRuleNetworkR5(t *testing.T) {
	a := NewRiskAssessor()

	cases := []struct {
		name      string
		target    NetTarget
		wantLevel Level
		wantHit   bool
	}{
		{"deny-listed protocol -> Deny", NetTarget{Protocol: "javascript", Host: "x"}, Deny, true},
		{"private range -> L2", NetTarget{Protocol: "https", Host: "10.0.0.5"}, L2, true},
		{"metadata IP is link-local -> L2", NetTarget{Protocol: "https", Host: "169.254.169.254"}, L2, true},
		{"ipv6 loopback -> L2", NetTarget{Protocol: "https", Host: "::1"}, L2, true},
		{"host:port spelling -> L2", NetTarget{Protocol: "https", Host: "192.168.1.1:8080"}, L2, true},
		{"unknown protocol -> L2", NetTarget{Protocol: "gopher", Host: "gopher.example.com"}, L2, true},
		{"domain outside allowlist -> L2", NetTarget{Protocol: "https", Host: "evil.example.net", Allowlist: []string{"example.com"}}, L2, true},
		{"URL over cap -> L2", NetTarget{Protocol: "https", Host: "example.com", URL: "https://example.com?q=" + strings.Repeat("a", 2049)}, L2, true},
		{"allowlisted https -> clean", NetTarget{Protocol: "https", Host: "api.example.com", Allowlist: []string{"example.com"}}, L0, false},
	}
	for _, tc := range cases {
		d := a.Assess("web.open", nil, Facts{Network: &tc.target})
		if tc.wantHit {
			if d.Level != tc.wantLevel || len(d.RulesHit) != 1 || d.RulesHit[0] != R5 {
				t.Fatalf("%s: want %v/[R5], got %v/%v (%s)", tc.name, tc.wantLevel, d.Level, d.RulesHit, d.Reason)
			}
			continue
		}
		if d.Level != L0 || len(d.RulesHit) != 0 {
			t.Fatalf("%s: want clean L0, got %v/%v (%s)", tc.name, d.Level, d.RulesHit, d.Reason)
		}
	}

	d := a.Assess("web.open", nil, Facts{}) // negative: not a network call
	if d.Level != L0 || len(d.RulesHit) != 0 {
		t.Fatalf("no network target: want L0/no hits, got %v/%v", d.Level, d.RulesHit)
	}
}

func TestRuleShellR6(t *testing.T) {
	a := NewRiskAssessor()

	cases := []struct {
		name      string
		facts     Facts
		wantLevel Level
	}{
		{"metachar pipe -> L2", Facts{ShellArgv: []string{"cmd", "/c", "dir | findstr x"}}, L2},
		{"metachar redirect+subshell -> L2", Facts{ShellArgv: []string{"sh", "-c", "cat f > g $(x)"}}, L2},
		{"clean allowlisted -> L1", Facts{ShellArgv: []string{"git", "status"}, ShellAllowlist: []string{"git"}}, L1},
		{"clean allowlisted full path -> L1", Facts{ShellArgv: []string{`C:\Program Files\Git\bin\git.exe`, "status"}, ShellAllowlist: []string{"git.exe"}}, L1},
		{"clean not allowlisted -> L2", Facts{ShellArgv: []string{"powershell", "-Command", "Get-Date"}, ShellAllowlist: []string{"git"}}, L2},
		{"string mode -> L2 regardless", Facts{ShellString: true, ShellArgv: []string{"git", "status"}, ShellAllowlist: []string{"git"}}, L2},
	}
	for _, tc := range cases {
		d := a.Assess("shell.exec", nil, tc.facts)
		if d.Level != tc.wantLevel || len(d.RulesHit) != 1 || d.RulesHit[0] != R6 {
			t.Fatalf("%s: want %v/[R6], got %v/%v (%s)", tc.name, tc.wantLevel, d.Level, d.RulesHit, d.Reason)
		}
	}

	d := a.Assess("fs.read", nil, Facts{}) // negative: not a shell call
	if d.Level != L0 || len(d.RulesHit) != 0 {
		t.Fatalf("no argv: want L0/no hits, got %v/%v", d.Level, d.RulesHit)
	}
}

func TestRuleBatchScaleR7(t *testing.T) {
	a := NewRiskAssessor()

	cases := []struct {
		name      string
		facts     Facts
		wantLevel Level
	}{
		{"exactly 50 -> L2 (edge)", Facts{BatchCount: 50}, L2},
		{"49 -> clean (edge)", Facts{BatchCount: 49}, L0},
		{"5000 -> L2", Facts{BatchCount: 5000}, L2},
		{"fallback to len(Paths) at 60 -> L2 (edge)", Facts{Paths: make([]string, 60)}, L2},
		{"fallback to len(Paths) at 3 -> clean", Facts{Paths: []string{"a", "b", "c"}}, L0},
	}
	for _, tc := range cases {
		d := a.Assess("fs.delete", nil, tc.facts)
		if d.Level != tc.wantLevel {
			t.Fatalf("%s: want %v, got %v/%v", tc.name, tc.wantLevel, d.Level, d.RulesHit)
		}
		if tc.wantLevel == L2 && (len(d.RulesHit) != 1 || d.RulesHit[0] != R7) {
			t.Fatalf("%s: want [R7], got %v", tc.name, d.RulesHit)
		}
	}
}

func TestRuleIrreversibleR8(t *testing.T) {
	a := NewRiskAssessor()

	cases := []struct {
		name      string
		facts     Facts
		wantLevel Level
	}{
		{"send class -> L2", Facts{Irreversible: []string{"send"}}, L2},
		{"power class -> L2", Facts{Irreversible: []string{"power"}}, L2},
		{"no classes -> clean", Facts{}, L0},
		{"overwrite existing -> L2 (edge)", Facts{OverwriteExisting: true}, L2},
		{"unknown class fails closed -> L2 (edge)", Facts{Irreversible: []string{"reformat-disk"}}, L2},
		{"close-window -> L2", Facts{Irreversible: []string{"close-window"}}, L2},
	}
	for _, tc := range cases {
		d := a.Assess("sys.power", nil, tc.facts)
		if d.Level != tc.wantLevel {
			t.Fatalf("%s: want %v, got %v/%v (%s)", tc.name, tc.wantLevel, d.Level, d.RulesHit, d.Reason)
		}
		if tc.wantLevel == L2 && (len(d.RulesHit) != 1 || d.RulesHit[0] != R8) {
			t.Fatalf("%s: want [R8], got %v", tc.name, d.RulesHit)
		}
	}
}
