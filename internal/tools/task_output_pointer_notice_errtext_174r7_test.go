package tools

import (
	"fmt"
	"strings"
	"testing"
)

// Ticket 174 AC#2d, third delivery (174-r7): a ruler over the one thing the
// canonicalize-error notice was still relaying - the upstream error text.
//
// The history of this arm, so the next reader does not re-derive it: 174-v1
// found that this `case err != nil` arm is the ONLY place in pointerNotice that
// concatenates err.Error() into model-visible text, and that it had zero
// criteria (MUT-V5: removing the whole arm left the package green). 174-r2
// pinned that the arm speaks and is not the unwired sentence; 174-r3 and
// 174-r6 pinned what its sentence must not relay (the filed path literal, the
// root list, road-back promises) - and 174-r3 measured, without committing a
// red assertion (A375), that the arm DOES relay the upstream text verbatim:
// `（tools: empty path）` on the whitespace arm and the ErrReparseDenied
// sentence on the junction arm. The relay point is unbounded by construction
// (174-a3 §3 form 丙: any future Canonicalize error text would ride it into the
// context), and A587 struck the "frozen sentence, needs owner" gate as one
// line, zero criteria moved, reversible. Form 乙 strips the parenthetical, and
// the target sentence is verbatim the one internal/agent/spill.go already
// speaks for the same failure - the two legs now say the same thing. This
// criterion is the ruler that keeps it stripped.
//
// The ban is derived INDEPENDENTLY: the forbidden string is whatever the SAME
// PathCanonicalizer constructor hands back for the SAME input the bridge is
// wired with - not one copied literal - so "relay whatever C26 says" is what
// goes red, not one spelling of it.
//
// One arm only, and deliberately portable: the whitespace shape is the one
// non-empty literal Canonicalize refuses on every platform (174-r2/174-r6's
// fixture; paths.go's TrimSpace guard runs ahead of risk.Resolve). The
// path-literal half of the same ban needs a real junction and is already
// pinned by TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3 (windows,
// ledged in scripts/portable-tests.sh) plus its portable form at ...174r6.
// This file adds the half that had no ruler at all, and re-asks the other two
// classes in the same sweep so the criterion stands alone: the branch must
// speak (its own claim on the record) AND nothing about the failure may be
// relayed.

// TestCanonicalizeErrorNoticeRelaysNoUpstreamErrorText174r7 is the AC#2d
// second half: the notice names the step and stops.
func TestCanonicalizeErrorNoticeRelaysNoUpstreamErrorText174r7(t *testing.T) {
	rootA, rootB := tempCanonical(t), tempCanonical(t)
	body := asciiRun(20000)

	paths := NewPathCanonicalizer([]string{rootA, rootB}, nil)
	// The upstream text, derived through the same constructor on the same
	// input. Non-nil is a precondition: if Canonicalize ever stops failing
	// here, this criterion must go red loudly, not turn vacuously green.
	_, cerr := paths.Canonicalize(canonicalizeFailsHere)
	if cerr == nil {
		t.Fatalf("fixture precondition broken: Canonicalize(%q) no longer fails", canonicalizeFailsHere)
	}
	upstream := cerr.Error()
	t.Logf("upstream text this criterion bans from the notice: %q", upstream)

	roster := NewTaskRoster()
	roster.Record("bg-errtxt-174r7", TaskOutput{Text: body, ArtifactPath: canonicalizeFailsHere})
	b := judgedTaskBridge(t, roster, paths)

	out := callTaskOutput(t, b, "bg-errtxt-174r7")
	t.Logf("err-text criterion verbatim: IsError=%v Truncated=%v announce=%q",
		out.IsError, out.Truncated, announceOf(out.Text))

	// Forward half. The branch speaks, by its own claim - the same sentence
	// 174-r2 pinned - and the pointer is still owed (PLAN.md:2564).
	if out.IsError || !out.Truncated {
		t.Fatalf("a recorded long output must still be answered as a truncation: %+v", out)
	}
	if !strings.Contains(out.Text, noticeLead) || !strings.Contains(out.Text, "读不到") {
		t.Fatalf("a path C26 cannot normalize must be handled as unreadable, got %q", announceOf(out.Text))
	}
	if !strings.Contains(out.Text, "连规范化都没通过") {
		t.Fatalf("the reply must name normalization as the step that failed, got %q", announceOf(out.Text))
	}
	if !strings.Contains(out.Text, "全文见 "+canonicalizeFailsHere) {
		t.Fatalf("failing closed is not withholding the pointer, got %q", out.Text)
	}
	// The notice region: from the first 注意： up to the frozen pointer
	// sentence - the same carve 174-r3/174-r6 use, for the same reason (the
	// pointer sentence legitimately prints the filed literal).
	i := strings.Index(out.Text, noticeLead)
	if i < 0 {
		t.Fatalf("the branch speaks above, so a notice region must exist: %q", out.Text)
	}
	notice := out.Text[i:]
	if j := strings.Index(out.Text[i:], "，全文见 "); j >= 0 {
		notice = out.Text[i : i+j]
	}
	t.Logf("notice region under the err-text ban: %q", notice)

	// The relay ban. Three forbidden classes, asked on the branch's own output:
	// the filed path literal, the authorized root list, and C26's own state -
	// the upstream text itself plus the packages' error prefixes, so a future
	// wrapper (risk.PathResolver, an fs wrapper) is caught by class, not by one
	// copy of today's string.
	for _, banned := range []string{
		canonicalizeFailsHere,                    // the filed literal, as filed
		fmt.Sprintf("%q", canonicalizeFailsHere), // ...or in its quoted spelling
		rootA, rootB,                             // the authorized root list, entry by entry
		strings.ToLower(rootA), strings.ToLower(rootB), // ...as the comparison fold spells it
		upstream,          // C26's own error text, verbatim
		"tools:", "risk:", // ...and the packages' error prefixes
	} {
		if strings.Contains(notice, banned) {
			t.Errorf("the canonicalize-error notice must not relay %q: notice %q", banned, notice)
		}
	}
}
