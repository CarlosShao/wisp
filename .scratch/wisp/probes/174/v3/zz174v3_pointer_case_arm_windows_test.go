// Probe for acceptance leg 174-v3 (ticket 174 AC#3(ii)): a spelling family
// NEITHER of 174-r4's two arms uses. r4's arms discriminate by 8.3 short names
// (over-report) and by workspace narrowing (under-report); this one flips CASE
// only - the filed string is strings.ToUpper of a real file's canonical path.
//
// C26 folds case (paths.go InAllowlist -> foldPath; SPEC-06 §4 treats canonical
// paths as case-insensitive), so the answer is INSIDE and the pointer must stay
// silent. Any containment computed on the FILED string by a case-sensitive
// prefix comparison answers OUTSIDE and speaks: the two answers disagree, which
// is what makes this arm able to tell them apart.
//
// It compiles into package tools through `go test -overlay` (see overlay json
// next to it), so internal/tools/ stays byte-clean. Injection surface is the
// real one (AGENTS §1.3): real bridge, real C26, real fs.read, real file.
package tools

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestV3PointerFollowsC26OnCaseFlippedSpelling174(t *testing.T) {
	dir := tempCanonical(t)
	body := asciiRun(20000)
	good := mustWriteFile(t, filepath.Join(dir, "tool-output-174v3-case.txt"), body)
	filed := strings.ToUpper(good)
	t.Logf("case-flip fixture: dir=%q real=%q filed=%q", dir, good, filed)
	if filed == good {
		t.Fatalf("fixture: %q has no letters to flip", good)
	}
	if strings.HasPrefix(filed, dir) {
		t.Fatalf("fixture broken: %q must NOT be a case-sensitive prefix of the canonical root %q,"+
			" or a lexical decision and C26's decision agree and this arm measures nothing", filed, dir)
	}

	roster := NewTaskRoster()
	roster.Record("bg-case", TaskOutput{Text: body, ArtifactPath: filed})
	b := judgedTaskBridge(t, roster, NewPathCanonicalizer([]string{dir}, nil))
	out := callTaskOutput(t, b, "bg-case")
	t.Logf("case arm verbatim: IsError=%v Truncated=%v Text=%q", out.IsError, out.Truncated, announceOf(out.Text))
	if out.IsError || !out.Truncated {
		t.Fatalf("a recorded long output must still be answered as a truncation: %+v", out)
	}
	for _, banned := range []string{noticeLead, "读不到", "不在你被授权的目录范围内", "未接线", "不可找回"} {
		if strings.Contains(out.Text, banned) {
			t.Errorf("C26 authorizes the real file whatever the case of the filed string; found %q in %q"+
				" - that is a decision made on the string, not on C26's answer", banned, announceOf(out.Text))
		}
	}
	if !strings.Contains(out.Text, filed) {
		t.Fatalf("the pointer must still name the filed path, got %q", announceOf(out.Text))
	}
	// The silence must be TRUE: following the pointer as printed works.
	rd := readThroughBridge(t, b, filed)
	if rd.IsError || rd.RiskLevel != "L0" || rd.Text != body {
		t.Fatalf("control leg: fs.read on the case-flipped spelling must return the whole %d bytes at L0,"+
			" got isError=%v level=%v len=%d", len(body), rd.IsError, rd.RiskLevel, len(rd.Text))
	}

	// Reverse control: same filed string, authority removed (the factory default
	// shape) - the same input must flip to speaking, or the silence above means
	// nothing. fs.read through the same bridge must agree it is refused.
	roster2 := NewTaskRoster()
	roster2.Record("bg-case-noauth", TaskOutput{Text: body, ArtifactPath: filed})
	b2 := judgedTaskBridge(t, roster2, NewPathCanonicalizer(nil, nil))
	out2 := callTaskOutput(t, b2, "bg-case-noauth")
	t.Logf("case arm, no root verbatim: IsError=%v Truncated=%v Text=%q",
		out2.IsError, out2.Truncated, announceOf(out2.Text))
	if !strings.Contains(out2.Text, noticeLead) || !strings.Contains(out2.Text, "不在你被授权的目录范围内") {
		t.Fatalf("with nothing authorized the case arm must say the pointer cannot be followed, got %q",
			announceOf(out2.Text))
	}
	if rd2 := readThroughBridge(t, b2, filed); !rd2.IsError || rd2.RiskLevel != "L2" {
		t.Fatalf("the notice says fs.read is refused; same bridge says %+v", rd2)
	}
}
