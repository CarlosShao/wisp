//go:build probe356

package tools

// Ticket 356 AC#0 probe (read-only). Build-tagged so it never joins the normal
// test count; run with: go test ./internal/tools/ -tags probe356 -run Test356AC0R3DenyPaired -v
//
// Measures the R3 A-class (absolute Deny) defect that pointerNotice cannot name:
// a ~/.ssh A-list file is Deny even when INSIDE an allowed root (R2 satisfied,
// only R3 catches). Uses an APPROVING gateSpy so the counts prove "approval was
// asked 0 times" rather than just "NoGate refuses".

import (
	"os"
	"path/filepath"
	"testing"
)

func Test356AC0R3DenyPaired(t *testing.T) {
	homeRaw := tempRaw(t) // t.TempDir()
	t.Setenv("USERPROFILE", homeRaw)
	t.Setenv("HOME", homeRaw)
	if err := os.MkdirAll(filepath.Join(homeRaw, ".ssh"), 0o700); err != nil {
		t.Fatal(err)
	}
	key := filepath.Join(homeRaw, ".ssh", "id_ed25519")
	if err := os.WriteFile(key, []byte("secret-356-contact"), 0o600); err != nil {
		t.Fatal(err)
	}
	home := mustCanonical(t, homeRaw)
	filed := mustCanonical(t, key)
	body := asciiRun(20000)

	// shape iii: INSIDE allowed root (home) but R3-Deny (the ~/.ssh A rule).
	{
		roster := NewTaskRoster()
		roster.Record("bg-356-iii", TaskOutput{Text: body, ArtifactPath: filed})
		g := &gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow}
		b := twoRoadsBridge174r1(t, roster, NewPathCanonicalizer([]string{home}, nil), g)
		out := callTaskOutput(t, b, "bg-356-iii")
		w, a := g.counts()
		t.Logf("SHAPE_iii notice: IsError=%v Truncated=%v Text=%q | gate window=%d approval=%d",
			out.IsError, out.Truncated, announceOf(out.Text), w, a)
		rd := readThroughBridge(t, b, key)
		t.Logf("SHAPE_iii fs.read: IsError=%v level=%v class=%q text=%q",
			rd.IsError, rd.RiskLevel, rd.ErrorClass, rd.Text)
	}

	// shape i POSITIVE CONTROL: R2-only (out-of-allowlist, NOT R3) + approving
	// gate -> the card is really asked (approval>=1) and the bytes come back.
	// Proves the ruler/route distinguishes R2 from R3 rather than always refusing.
	{
		plain := mustWriteFile(t, filepath.Join(t.TempDir(), "plain-356-outside.txt"), body)
		roster := NewTaskRoster()
		roster.Record("bg-356-i", TaskOutput{Text: body, ArtifactPath: mustCanonical(t, plain)})
		g := &gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow}
		b := twoRoadsBridge174r1(t, roster, NewPathCanonicalizer(nil, nil), g)
		rd := readThroughBridge(t, b, plain)
		w, a := g.counts()
		t.Logf("SHAPE_i fs.read(R2-only): IsError=%v level=%v bytes=%d | gate window=%d approval=%d",
			rd.IsError, rd.RiskLevel, len(rd.Text), w, a)
	}

	// shape ii: OUTSIDE allowed root + R3-Deny (Deny dominates out-of-allowlist).
	{
		roster := NewTaskRoster()
		roster.Record("bg-356-ii", TaskOutput{Text: body, ArtifactPath: filed})
		g := &gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow}
		b := twoRoadsBridge174r1(t, roster, NewPathCanonicalizer(nil, nil), g)
		out := callTaskOutput(t, b, "bg-356-ii")
		w, a := g.counts()
		t.Logf("SHAPE_ii notice: IsError=%v Truncated=%v Text=%q | gate window=%d approval=%d",
			out.IsError, out.Truncated, announceOf(out.Text), w, a)
		rd := readThroughBridge(t, b, key)
		t.Logf("SHAPE_ii fs.read: IsError=%v level=%v class=%q text=%q",
			rd.IsError, rd.RiskLevel, rd.ErrorClass, rd.Text)
	}
}
