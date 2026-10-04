package agent

import (
	"strings"
	"testing"
)

// Ticket 174, AC#3(ii) on the SECOND leg (internal/agent/spill.go), landed by
// 174-r4. Same forbidden zone, same missing ruler:
//
// spill.go's pointerNotice asks C26 in two steps - Canonicalize(path), then
//	InAllowlist(canon) - and the whole point of that order is that the
//	authorization question is about C26's ANSWER, not about the string the host
//	happened to write down. Every criterion of ticket 174's spill leg answers
//	`return f.inRoot` (spill_pointer_honesty_174_test.go:79-82): the fake never
//	looks at its argument, so "the ask was made with the wrong string" is
//	invisible there - and in production, a real canonicalizer returns a
//	re-spelled path (8.3 expansion, separators, UNC), so asking with the filed
//	string instead of the canonical one flips verdicts on legal pointers.
//
// The judge below therefore answers PER STRING. Two arms, opposite directions,
// over the SAME spilled file, plus a third observation (what string the ask
// actually carried) so the silence cannot be bought by a fake that always says
// yes:
//
//	arm 1  canonical spelling IN, filed spelling OUT  -> must stay silent, and
//	         InAllowlist must have been handed the canonical spelling
//	arm 2  canonical spelling OUT, filed spelling IN  -> must speak, and speak
//	         with BOTH roads back (that is this leg's existing wording)
//
// Nothing here re-words or re-cuts anything the spill leg already produces: the
// pre-fix-sentence pin, the head/tail budget pins and the no-second-path pin of
// ticket 174's own spill criteria stay untouched (that file is pending an
// acceptance flip, and editing it would desync the accept leg's roster diff).
// AGENTS §1.3: this is an authorization-input seam, not a mocked-away
// dependency - PointerJudge is the production interface spill.go declares, and
// what is under test is exactly the composition inside pointerNotice.

// argAware174r4Judge answers the allowlist question by looking at the string it
// is asked about, which is the only way to observe which string the ask used.
type argAware174r4Judge struct {
	canonOut   string
	inFor      map[string]bool
	canonIn    []string
	allowIn    []string
	canonCalls int
}

func (j *argAware174r4Judge) Canonicalize(raw string) (string, error) {
	j.canonCalls++
	j.canonIn = append(j.canonIn, raw)
	return j.canonOut, nil
}

func (j *argAware174r4Judge) InAllowlist(canonical string) bool {
	j.allowIn = append(j.allowIn, canonical)
	return j.inFor[canonical]
}

// stubOf174r4 cuts the bracketed stub out of a spilled reply so a reading that
// gets quoted into an evidence file names the sentence under test instead of
// 68 KB of payload around it. The pointer parse of ticket 174 keeps working on
// this window, which is why it is cut the same way both legs cut it.
func stubOf174r4(text string) string {
	i := strings.Index(text, "[…")
	if i < 0 {
		return text
	}
	if j := strings.Index(text[i:], "…]"); j >= 0 {
		return text[i : i+j+len("…]")]
	}
	return text[i:]
}

// TestSpillPointerAuthorityAsksWithTheCanonicalAnswer174r4 is that pin.
func TestSpillPointerAuthorityAsksWithTheCanonicalAnswer174r4(t *testing.T) {
	// arm 1: only C26's canonical spelling is authorized. The filed path is the
	// same file under a different spelling - the shape production hits when the
	// roster/host hands over an alias, a short name, or any string C26 rewrites.
	filed := "C:\\mock\\filed\\TOOL-O~1.TXT"
	canonical := "C:\\mock\\filed\\tool-output-174-r4.txt"
	j := &argAware174r4Judge{
		canonOut: canonical,
		inFor:    map[string]bool{canonical: true, filed: false},
	}
	out, _, _ := spill174OverThreshold(t, j)
	t.Logf("arm 1 verbatim: Text=%q", stubOf174r4(out.Text))
	if strings.Contains(out.Text, spill174NoticeLead) {
		t.Errorf("C26 answers INSIDE for the string it returned, so a pointer to it must stay quiet;"+
			" got a notice in %q - the authorization question was not asked with C26's answer"+
			" (asked with: %v, canonical: %q)", stubOf174r4(out.Text), j.allowIn, canonical)
	}
	if len(j.canonIn) != 1 || j.canonIn[0] != out.Path {
		t.Errorf("Canonicalize must be asked exactly once with the path the pointer prints:"+
			" got %v, pointer=%q", j.canonIn, out.Path)
	}
	if len(j.allowIn) != 1 || j.allowIn[0] != canonical {
		t.Errorf("InAllowlist must be asked with Canonicalize's answer, not with the filed string:"+
			" got %v, want [%q]", j.allowIn, canonical)
	}

	// arm 2 (reverse control): the SAME judge shape, answers crossed the other
	// way - canonical OUT, filed IN. A reply that stays quiet here is the
	// under-warning half of the exact lie this ticket clears.
	j2 := &argAware174r4Judge{
		canonOut: canonical,
		inFor:    map[string]bool{canonical: false, filed: true},
	}
	out2, _, _ := spill174OverThreshold(t, j2)
	t.Logf("arm 2 verbatim: Text=%q", stubOf174r4(out2.Text))
	if !strings.Contains(out2.Text, spill174NoticeLead) ||
		!strings.Contains(out2.Text, "不在你被授权的目录范围内") {
		t.Fatalf("C26 answered OUT for the string it returned, so the pointer must say so; got %q", stubOf174r4(out2.Text))
	}
	// And it must say BOTH roads back, because [fs] allowed_dirs is a risk-tier
	// input and not an execution wall (174-r2's leg already words it that way;
	// re-pinned here only because this arm is the one that would drift).
	if !strings.Contains(out2.Text, "allowed_dirs") || !strings.Contains(out2.Text, "L2 卡") {
		t.Errorf("the notice owes the reader both roads back, got %q", stubOf174r4(out2.Text))
	}
	if len(j2.allowIn) != 1 || j2.allowIn[0] != canonical {
		t.Errorf("the reverse arm must ask with the same string the forward arm did: got %v", j2.allowIn)
	}
}
