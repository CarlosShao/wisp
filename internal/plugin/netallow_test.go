package plugin

import (
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// TestAllowlistMatcherMirrorsR5 pins AC5's matching semantics against the
// central rule's own behaviour: the plugin gate must not invent a second,
// subtly wider notion of "in the allowlist".
func TestAllowlistMatcherMirrorsR5(t *testing.T) {
	allow := []string{"open.feishu.cn", "Api.Example.com"}
	cases := []struct {
		host string
		want bool
	}{
		{"open.feishu.cn", true},
		{"OPEN.FEISHU.CN", true},
		{"search.api.example.com", true}, // subdomain of an entry (case-insensitive entry)
		{"api.example.com:443", true},    // :port on the host is ignored
		{"api.example.com.", false},      // trailing dot is not normalization here
		{"feishu.cn", false},             // parent of an entry never matches
		{"evil-example.com", false},      // suffix without the dot boundary
		{"evil.open.feishu.cn", true},    // genuine subdomain
		{"", false},
		{"https://open.feishu.cn", false}, // a scheme is not a host
	}
	for _, c := range cases {
		if got := NetDomainAllowed(c.host, allow); got != c.want {
			t.Errorf("NetDomainAllowed(%q) = %v, want %v", c.host, got, c.want)
		}
	}
}

func TestAllowlistToolFormUsesOwnDeclaration(t *testing.T) {
	def := ToolDef{Name: "lark.search", NetAllowlist: []string{"api.example.com"}}
	if !def.NetAllowed("search.api.example.com") {
		t.Error("subdomain of the declared entry must pass")
	}
	if def.NetAllowed("api.evil.com") {
		t.Error("undeclared host must not pass")
	}
}

func TestAllowlistMergeDedupesAndLowercases(t *testing.T) {
	got := MergeNetAllowlists(
		[]string{"A.Example.com", "open.feishu.cn"},
		[]string{"a.example.com", "  ", "cdn.example.net"},
		nil,
		[]string{"a.example.com:443", "A.EXAMPLE.COM"},
	)
	want := "a.example.com|open.feishu.cn|cdn.example.net"
	if j := strings.Join(got, "|"); j != want {
		t.Fatalf("merged = %q, want %q", j, want)
	}
}

func TestAllowlistLevelForRiskFailsClosed(t *testing.T) {
	for _, s := range []string{"L0", "L1", "L2"} {
		if _, err := LevelForRisk(s); err != nil {
			t.Errorf("LevelForRisk(%q): %v", s, err)
		}
	}
	for _, s := range []string{"", "l0", "L3", "LOW"} {
		if _, err := LevelForRisk(s); err == nil {
			t.Errorf("LevelForRisk(%q) accepted an unknown C2 spelling", s)
		}
	}
}

// TestAllowlistR5InListPassesL0 is AC5's first half, e2e through the REAL
// central assessor: an L0-declared plugin tool calling a host inside its own
// declared net_allowlist keeps its declared floor (no R5 hit).
func TestAllowlistR5InListPassesL0(t *testing.T) {
	def := ToolDef{
		Name: "lark.calendar_agenda", Risk: RiskL0,
		Capabilities: []string{"net"}, NetAllowlist: []string{"open.feishu.cn"},
	}
	dec, err := AssessPluginNetCall(def, risk.NetTarget{
		Protocol: "https", Host: "open.feishu.cn", URL: "https://open.feishu.cn/agenda",
	}, nil, nil)
	if err != nil {
		t.Fatalf("AssessPluginNetCall: %v", err)
	}
	if dec.Level != risk.L0 {
		t.Fatalf("in-list verdict = %v (%s), want L0 pass", dec.Level, dec.Reason)
	}
	for _, r := range dec.RulesHit {
		if r == risk.R5 {
			t.Fatalf("in-list call must not hit R5, got %v", dec.RulesHit)
		}
	}
}

// TestAllowlistR5OutOfListEscalatesL2 is AC5's second half: the same tool
// reaching a host outside the declared list is auto-escalated to L2 with the
// R5 rule named (SPEC-07 §4 "白名单外 → 升 L2"; never a silent pass).
func TestAllowlistR5OutOfListEscalatesL2(t *testing.T) {
	def := ToolDef{
		Name: "lark.calendar_agenda", Risk: RiskL0,
		Capabilities: []string{"net"}, NetAllowlist: []string{"open.feishu.cn"},
	}
	dec, err := AssessPluginNetCall(def, risk.NetTarget{
		Protocol: "https", Host: "exfil.evil.net", URL: "https://exfil.evil.net/stolen",
	}, nil, nil)
	if err != nil {
		t.Fatalf("AssessPluginNetCall: %v", err)
	}
	if dec.Level != risk.L2 {
		t.Fatalf("out-of-list verdict = %v, want L2", dec.Level)
	}
	var hitR5 bool
	for _, r := range dec.RulesHit {
		if r == risk.R5 {
			hitR5 = true
		}
	}
	if !hitR5 {
		t.Fatalf("out-of-list verdict must hit R5, got %v (%s)", dec.RulesHit, dec.Reason)
	}
	if !strings.Contains(dec.Reason, "exfil.evil.net") {
		t.Errorf("reason must name the offending host, got %q", dec.Reason)
	}
}

// TestAllowlistR5ProtocolDeniesAreNotLaundered proves the plugin allowlist
// cannot weaken R5's own protocol rules: an exact-listed "evil.example" is
// still refused/Denied for a denied protocol, and a private-range IP stays
// L2 even when someone puts it in the allowlist.
func TestAllowlistR5ProtocolDeniesAreNotLaundered(t *testing.T) {
	def := ToolDef{
		Name: "x.y", Risk: RiskL0,
		Capabilities: []string{"net"}, NetAllowlist: []string{"evil.example", "127.0.0.1"},
	}
	dec, err := AssessPluginNetCall(def, risk.NetTarget{
		Protocol: "javascript", Host: "evil.example", URL: "javascript:evil.example",
	}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Level != risk.Deny {
		t.Fatalf("denied protocol inside the allowlist = %v, want Deny (allowlist must not launder R5)", dec.Level)
	}
	dec, err = AssessPluginNetCall(def, risk.NetTarget{
		Protocol: "http", Host: "127.0.0.1", URL: "http://127.0.0.1/x",
	}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Level != risk.L2 {
		t.Fatalf("loopback inside the allowlist = %v, want L2 (SSRF defence outranks the plugin list)", dec.Level)
	}
}
