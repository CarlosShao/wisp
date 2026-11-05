package risk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 230 AC#4 (the risk-side seam): the B-tier override grew a SECOND,
// separately keyed source - the operator's standing list from config.toml's
// risk.blacklist_overrides - without changing what Gate means.
//
// These are the four readings the AC asks for, one subtest each:
//  1. Gate(c, m) and GateStanding(c, m, nil) are the same verdict, so every
//     caller that predates the seam is behaviorally unchanged;
//  2. a persisted B entry allows the file, and its absence keeps the
//     fail-closed "one L2 confirm" answer;
//  3. an A-tier path placed in the persisted set (and in the L2 set, for good
//     measure) is still ClassA with no Allow and no NeedL2;
//  4. exactly ONE security line per allowed use, and it names the source
//     (source=l2-confirm / source=config:risk.blacklist_overrides); the config
//     line and the config Reason never claim an L2 confirm.
//
// The fixture files live in a temp dir: the B tier is decided by file NAME, and
// the comparison key is OverrideKey output, so no profile anchor is involved.

// logLines230 swaps the package security sink for one test and returns the
// slice the lines are appended to.
func logLines230(t *testing.T) *[]string {
	t.Helper()
	lines := &[]string{}
	prev := Logf
	Logf = func(format string, args ...any) {
		*lines = append(*lines, fmt.Sprintf(format, args...))
	}
	t.Cleanup(func() { Logf = prev })
	return lines
}

// writeFile230 creates one file under a fresh temp dir and returns both its
// spelling and its C26 canonical form.
func writeFile230(t *testing.T, parts ...string) (raw, canonical string) {
	t.Helper()
	dir := t.TempDir()
	full := filepath.Join(append([]string{dir}, parts...)...)
	if err := os.MkdirAll(filepath.Dir(full), 0o700); err != nil {
		t.Fatalf("mkdir for %q: %v", full, err)
	}
	if err := os.WriteFile(full, []byte("fixture"), 0o600); err != nil {
		t.Fatalf("write %q: %v", full, err)
	}
	res, err := Resolve(full, nil)
	if err != nil {
		t.Fatalf("Resolve(%q): %v", full, err)
	}
	return full, res.Canonical
}

// TestGateStandingIsGateWhenThePersistedSourceIsNil proves reading 1 by
// comparing the two entry points' whole PathDecision, over every branch that
// has an override map today: default B, allowed B, A with a key that would
// cover it, A without one, and a path outside both tiers.
func TestGateStandingIsGateWhenThePersistedSourceIsNil(t *testing.T) {
	_, bPath := writeFile230(t, "proj", ".env")
	_, aPath := writeFile230(t, "repo", ".git", "config")
	_, plain := writeFile230(t, "notes.txt")

	bKey := OverrideKey(bPath)
	aKey := OverrideKey(aPath)

	cases := []struct {
		name string
		path string
		l2   map[string]bool
	}{
		{"b-default", bPath, nil},
		{"b-l2-confirmed", bPath, map[string]bool{bKey: true}},
		{"a-default", aPath, nil},
		{"a-keyed-anyway", aPath, map[string]bool{aKey: true}},
		{"outside-both-tiers", plain, nil},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			logs := logLines230(t)
			got := Gate(tc.path, tc.l2)
			gateLines := append([]string(nil), *logs...)
			*logs = nil
			want := GateStanding(tc.path, tc.l2, nil)
			standingLines := append([]string(nil), *logs...)

			if got != want {
				t.Fatalf("Gate and GateStanding(.., nil) disagree:\n\tGate           %+v\n\tGateStanding   %+v", got, want)
			}
			if strings.Join(gateLines, "\n") != strings.Join(standingLines, "\n") {
				t.Fatalf("the two entry points logged different lines:\n\tGate         %q\n\tGateStanding %q",
					gateLines, standingLines)
			}
		})
	}
}

// TestGateStandingPersistedExemptionAllowsButSaysSo covers readings 2 and 4 for
// the config source: absent it is "ask first"; present it allows; the single
// line and the returned Reason both name the config list and neither claims an
// L2 confirm.
func TestGateStandingPersistedExemptionAllowsButSaysSo(t *testing.T) {
	_, bPath := writeFile230(t, "proj", ".env")
	key := OverrideKey(bPath)
	if key != normPath(bPath) {
		t.Fatalf("OverrideKey(%q) = %q, want the comparison fold %q: the exported key must BE "+
			"the fold overrideApplies uses, not a second implementation of it",
			bPath, key, normPath(bPath))
	}

	logs := logLines230(t)

	d := GateStanding(bPath, nil, nil)
	if d.Class != ClassB || d.Allow || !d.NeedL2 {
		t.Fatalf("without any source the B-tier verdict changed (fail-closed default broken): %+v", d)
	}
	if len(*logs) != 1 {
		t.Fatalf("default deny emitted %d security lines, want exactly 1: %q", len(*logs), *logs)
	}
	*logs = nil

	d = GateStanding(bPath, nil, map[string]bool{key: true})
	if d.Class != ClassB || !d.Allow || d.NeedL2 {
		t.Fatalf("a persisted B-tier exemption did not allow the file: %+v", d)
	}
	if len(*logs) != 1 {
		t.Fatalf("a config-sourced allow emitted %d security lines, want exactly 1: %q", len(*logs), *logs)
	}
	line := (*logs)[0]
	if !strings.Contains(line, "source=config:risk.blacklist_overrides") {
		t.Errorf("the allow line does not name the config source: %q", line)
	}
	if strings.Contains(line, "l2-confirm") || strings.Contains(strings.ToLower(line), "l2 confirmed") {
		t.Errorf("a config-sourced allow claims an L2 confirm: %q", line)
	}
	if !strings.Contains(line, "rule=id_*") && !strings.Contains(line, "rule=.env") {
		t.Errorf("the allow line does not name the rule that matched: %q", line)
	}
	if !strings.Contains(d.Reason, "risk.blacklist_overrides") {
		t.Errorf("the config-sourced Reason does not name its source: %q", d.Reason)
	}
	if strings.Contains(d.Reason, "one L2 confirm + logged") {
		t.Errorf("the config-sourced Reason claims the L2-confirm sentence: %q", d.Reason)
	}
	*logs = nil

	// The L2 source keeps its own sentence, and its line says l2-confirm.
	d = GateStanding(bPath, map[string]bool{key: true}, nil)
	if !d.Allow {
		t.Fatalf("the L2 confirmation stopped unlocking its own file: %+v", d)
	}
	if len(*logs) != 1 || !strings.Contains((*logs)[0], "source=l2-confirm") {
		t.Fatalf("the L2-sourced allow must emit exactly one line naming source=l2-confirm, got %q", *logs)
	}
	if strings.Contains((*logs)[0], "config:") {
		t.Errorf("the L2-sourced allow line names the config source: %q", (*logs)[0])
	}
	*logs = nil

	// Both sources covering one file: still exactly ONE line, and it is the
	// human-answered one (the stronger claim).
	d = GateStanding(bPath, map[string]bool{key: true}, map[string]bool{key: true})
	if !d.Allow || len(*logs) != 1 || !strings.Contains((*logs)[0], "source=l2-confirm") {
		t.Fatalf("both sources must yield one allow and one line naming l2-confirm, got allow=%v lines=%q",
			d.Allow, *logs)
	}
}

// TestGateStandingPersistedExemptionDoesNotUnlockA is reading 3: the tier check
// runs before either map, so an A path keyed into BOTH sources is still ClassA,
// with no Allow and no L2 route, and it emits exactly one A-tier deny line.
func TestGateStandingPersistedExemptionDoesNotUnlockA(t *testing.T) {
	_, aPath := writeFile230(t, "repo", ".git", "config")
	if got := Classify(aPath); got != ClassA {
		t.Fatalf("fixture guard: %q classifies %s, want A", aPath, got)
	}
	key := OverrideKey(aPath)
	both := map[string]bool{key: true}

	logs := logLines230(t)
	d := GateStanding(aPath, nil, both)
	if d.Class != ClassA || d.Allow || d.NeedL2 {
		t.Fatalf("A tier was unlocked by a persisted exemption: %+v", d)
	}
	if len(*logs) != 1 {
		t.Fatalf("the A deny emitted %d lines, want exactly 1: %q", len(*logs), *logs)
	}
	if strings.Contains((*logs)[0], "source=") {
		t.Errorf("the A-tier deny line names an override source, which would mean a map was read: %q", (*logs)[0])
	}
	*logs = nil

	d = GateStanding(aPath, both, both)
	if d.Class != ClassA || d.Allow || d.NeedL2 {
		t.Fatalf("A tier was unlocked with both sources keyed: %+v", d)
	}
	if len(*logs) != 1 {
		t.Fatalf("the A deny emitted %d lines with both sources keyed, want exactly 1: %q", len(*logs), *logs)
	}
}
