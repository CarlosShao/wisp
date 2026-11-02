package plugin

import (
	"strings"
	"testing"
)

// Ticket 50 AC2 — C16 host_api version gate.
//
// Verbatim criteria (issue 50 AC2 / PLAN.md C16 / SPEC-07 §4,§6):
//   "^1.0 plugin vs host 1.x loads; vs 2.x host refuses with clear error",
//   "主版本不兼容 → 拒绝加载并给明确错误，不得静默降级 / 静默跳过".

// gate marks a row whose refusal is a version-gate verdict (as opposed to a
// range/host string the loader cannot even parse); gate verdicts must name the
// plugin, the range, and the host version in one error message.
func TestAC2CaretGate(t *testing.T) {
	tests := []struct {
		name    string
		rng     string
		host    string
		wantOK  bool
		wantErr string // substring; empty when wantOK
		gate    bool
	}{
		// AC2 positive: ^1.0 against a host 1.x loads.
		{"caret1.0/host1.0.0", "^1.0", "1.0.0", true, "", false},
		{"caret1.0/host1.2.3", "^1.0", "1.2.3", true, "", false},
		{"caret1.0/host1.99.99", "^1.0", "1.99.99", true, "", false},
		{"caret1.0/host1.0.0+build", "^1.0", "1.0.0+build.7", true, "", false},
		{"caret1.0/host1.0.1-rc.1", "^1.0", "1.0.1-rc.1", true, "", false},
		{"caret1.2/host1.2.0", "^1.2", "1.2.0", true, "", false},
		{"caret1.2/host1.5.0", "^1.2", "1.5.0", true, "", false},

		// AC2 negative: a 2.x host refuses, with the major-version wording.
		{"caret1.0/host2.0.0", "^1.0", "2.0.0", false, "host major 1", true},
		{"caret1.0/host2.4.1", "^1.0", "2.4.1", false, "refused", true},
		{"caret1.0/host3.0.0", "^1.0", "3.0.0", false, "major-version incompatibility", true},
		{"caret2.1/host1.9.9", "^2.1", "1.9.9", false, "host major 2", true},

		// Boundaries the gate must not get wrong.
		{"caret1.2/host1.1.9 below min", "^1.2", "1.1.9", false, "below the minimum", true},
		{"caret1.0/host0.9.9 major", "^1.0", "0.9.9", false, "host major 1", true},
		{"caret1.2.3/host1.2.2", "^1.2.3", "1.2.2", false, "below the minimum", true},
		{"caret1.2.3/host1.2.3", "^1.2.3", "1.2.3", true, "", false},
		{"caret1.2.3/host1.9.0", "^1.2.3", "1.9.0", true, "", false},

		// Major 0: minor is the compatibility position (standard caret), never a
		// free pass for a host that has not reached the declared API level.
		{"caret0.2/host0.2.5", "^0.2", "0.2.5", true, "", false},
		{"caret0.2/host0.3.0", "^0.2", "0.3.0", false, "refused", true},
		{"caret0.2/host0.1.9", "^0.2", "0.1.9", false, "refused", true},
		{"caret0.0/host0.0.7", "^0.0", "0.0.7", true, "", false},
		{"caret0.0/host0.1.0", "^0.0", "0.1.0", false, "below the minimum", true},

		// "^1" is the same constraint as "^1.0.0" under the mainstream caret
		// reading (node-semver), so it is accepted, not refused.
		{"caret major only ok", "^1", "1.4.0", true, "", false},
		{"caret major only upper", "^1", "2.0.0", false, "host major 1", true},

		// Grammar the contract does not define => explicit refusal, not silence.
		{"bare 1.0", "1.0", "1.4.0", false, "unsupported range syntax", false},
		{"exact =1.0", "=1.0", "1.4.0", false, "unsupported range syntax", false},
		{"comparator >=1", ">=1", "1.4.0", false, "unsupported range syntax", false},
		{"wildcard ^*", "^*", "1.4.0", false, "non-digit", false},
		{"caret alone", "^", "1.4.0", false, "caret needs a version", false},
		{"empty range", "", "1.4.0", false, "host_api is empty", false},
		{"caret leading zero", "^01.0", "1.4.0", false, "leading zero", false},
		{"caret non numeric", "^1.x", "1.4.0", false, "non-digit", false},
		{"caret too many parts", "^1.0.0.0", "1.4.0", false, "at most", false},

		// A host string that cannot be parsed must refuse, not load.
		{"host empty", "^1.0", "", false, "host version", false},
		{"host non numeric", "^1.0", "1.x.0", false, "host version", false},
		{"host two parts", "^1.0", "1.2", false, "host version", false},
		{"host leading zero", "^1.0", "01.2.3", false, "host version", false},
		{"host space", "^1.0", "1.2.3 ", false, "host version", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckHostAPI("com.example.lark", tc.rng, tc.host)
			if tc.wantOK {
				if err != nil {
					t.Fatalf("want load, got refusal: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("want refusal for host_api=%q host=%q, got nil (silent accept)", tc.rng, tc.host)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error %q missing required fragment %q", err, tc.wantErr)
			}
			// Every refusal names the plugin: C16 forbids a silent skip, so the
			// message must say which plugin was refused.
			if !strings.Contains(err.Error(), "com.example.lark") {
				t.Fatalf("refusal %q does not name the plugin", err)
			}
			if !tc.gate {
				return
			}
			// Gate verdicts additionally name both sides of the comparison.
			for _, must := range []string{tc.rng, tc.host} {
				if !strings.Contains(err.Error(), must) {
					t.Fatalf("error %q does not name %q", err, must)
				}
			}
		})
	}
}

// TestAC2ManifestAccessorDrivesGate pins that the loader-facing entry point
// reads the range from the manifest field the validator sealed (host_api), not
// from a re-typed literal.
func TestAC2ManifestAccessorDrivesGate(t *testing.T) {
	m := Manifest{ID: "com.example.lark", HostAPI: "^1.0"}
	if err := m.CompatibleWith("1.8.0"); err != nil {
		t.Fatalf("1.8.0 host must load a ^1.0 plugin: %v", err)
	}
	if err := m.CompatibleWith("2.0.0"); err == nil {
		t.Fatal("2.0.0 host must refuse a ^1.0 plugin")
	}
	m.HostAPI = "^2.0"
	if err := m.CompatibleWith("2.0.0"); err != nil {
		t.Fatalf("range change must take effect: %v", err)
	}
	if err := m.CompatibleWith("1.9.9"); err == nil {
		t.Fatal("^2.0 plugin on 1.9.9 host must refuse")
	}
}

// TestAC2ParseManifestSealsGrammar: a manifest carrying an unparseable
// host_api is rejected by ParseManifest, so the gate never sees a range it
// cannot interpret (the AC1 field check and the AC2 gate are the same seam).
func TestAC2ParseManifestSealsGrammar(t *testing.T) {
	raw := strings.Replace(validManifest, `host_api = "^1.0"`, `host_api = "1.0"`, 1)
	if _, err := ParseManifest([]byte(raw)); err == nil {
		t.Fatal(`host_api = "1.0" (outside the C16 caret grammar) must not parse clean`)
	}
	if _, err := ParseManifest([]byte(validManifest)); err != nil {
		t.Fatalf("the contract example manifest must parse: %v", err)
	}
}

// TestAC2DefaultHostVersionIsReadFromBuildinfo documents the wiring: the gate's
// host side is the stamped product version, and the dev default fails closed.
func TestAC2DefaultHostVersionIsReadFromBuildinfo(t *testing.T) {
	if strings.ContainsAny(DefaultHostVersion, " \t") {
		t.Fatalf("DefaultHostVersion %q is not usable as a host version", DefaultHostVersion)
	}
	// "go run" / non-script builds stamp 0.0.0-dev, which matches no ^1.x range:
	// a plugin must never load against an unidentified host build.
	if err := CheckHostAPI("com.example.lark", "^1.0", "0.0.0-dev"); err == nil {
		t.Fatal("dev default host version must refuse a ^1.0 plugin (fail closed)")
	}
}
