// Probe fixture for acceptance leg 174-v1 (non-implementer). It is NOT part of
// the tree: it compiles into package tools through `go test -overlay` (see
// overlay.json next to it), so internal/tools/ stays byte-clean.
//
// Purpose: 174-v1's OWN positive case for ticket 174 AC#2's reverse arm -
// a pointer that is healthy (inside an authorized root, a real regular file on
// disk, output long enough to spill) must produce ZERO occurrences of
// "注意：" in the WHOLE reply, and the very same bridge must read the whole
// file back at L0. Two legs differ only in how the host filed the path:
// canonicalized, or raw as a spill would hand it over.
package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

func TestProbe174V1HealthyPointerStaysSilentAndReadable(t *testing.T) {
	root := tempCanonical(t)
	nested := filepath.Join(root, "spill", "day1")
	if err := os.MkdirAll(nested, 0o700); err != nil {
		t.Fatal(err)
	}
	body := asciiRun(30000)
	f := mustWriteFile(t, filepath.Join(nested, "tool-output-174-v1-healthy.txt"), body)

	for _, leg := range []struct {
		name string
		path string
	}{
		{"canonical-recorded-path", mustCanonical(t, f)},
		{"raw-recorded-path", f},
	} {
		t.Run(leg.name, func(t *testing.T) {
			roster := NewTaskRoster()
			roster.Record("bg-v1", TaskOutput{Text: body, ArtifactPath: leg.path})
			b := judgedTaskBridge(t, roster, NewPathCanonicalizer([]string{root}, nil))

			out := callTaskOutput(t, b, "bg-v1")
			if out.IsError || !out.Truncated {
				t.Fatalf("a healthy long output must answer as a plain truncation: %+v", out)
			}
			if n := strings.Count(out.Text, noticeLead); n != 0 {
				t.Errorf("noticeLead count = %d, want 0 across the WHOLE reply", n)
			}
			for _, banned := range []string{"读不到", "不存在", "不是一般文件", "未接线", "不可找回"} {
				if strings.Contains(out.Text, banned) {
					t.Errorf("healthy reply carries %q", banned)
				}
			}
			if !strings.Contains(out.Text, "全文见 ") {
				t.Errorf("the pointer must be given: %q", announceOf(out.Text))
			}
			if !strings.HasPrefix(out.Text, body[:2000]) || !strings.HasSuffix(out.Text, body[len(body)-800:]) {
				t.Error("the head/tail windows must be untouched on the healthy arm")
			}
			rd := readThroughBridge(t, b, f)
			if rd.IsError || rd.RiskLevel != "L0" || rd.Text != body {
				t.Fatalf("same bridge must read back all %d bytes at L0: isError=%v level=%s len=%d",
					len(body), rd.IsError, rd.RiskLevel, len(rd.Text))
			}
			t.Logf("174-v1 positive leg %s: bytes=%d noticeCount=0 fs.read level=%s", leg.name, len(body), rd.RiskLevel)
		})
	}
}

// v1AllowGate is 174-v1's own gate: it answers the L2 card the way a user who
// says "允许" would. NoGate (what the tree ships today) always rejects, so the
// measured shape of "outside the allowlist" is a refusal - this leg asks
// whether that is the ONLY shape.
type v1AllowGate struct{}

func (v1AllowGate) PendingWindow(context.Context, Decision) (Answer, string) {
	return AnswerAllow, "174-v1 probe: L1 allowed"
}

func (v1AllowGate) PendingApproval(context.Context, Decision) (Answer, string) {
	return AnswerAllow, "174-v1 probe: L2 allowed by the user"
}

// TestProbe174V1OutsideRootIsReadableOnceTheCardIsAnswered probes the direction
// of the notice's wording: "fs.read 会被拒；要用户先把所属目录加进 [fs]
// allowed_dirs 才读得回来". If an answered L2 card already lets fs.read return
// the bytes, that sentence is stricter than the machine (the road is gated, not
// closed), and the pointer is over-promising in the safe direction.
func TestProbe174V1OutsideRootIsReadableOnceTheCardIsAnswered(t *testing.T) {
	root := tempCanonical(t)
	elsewhere := tempCanonical(t)
	body := asciiRun(4096)
	f := mustWriteFile(t, filepath.Join(elsewhere, "outside-root.txt"), body)

	reg := NewRegistry()
	paths := NewPathCanonicalizer([]string{root}, nil)
	for _, e := range BuiltinFSEntries(FSDeps{Paths: paths}) {
		if err := reg.Register(e); err != nil {
			t.Fatal(err)
		}
	}
	b := New(Options{
		Registry: reg, Paths: paths, Gate: v1AllowGate{},
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil}),
		Logf:       func(string, ...any) {},
	})
	out, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "174-v1", CallID: "c-v1-card", Name: "fs.read",
		Args: mustArgs(t, map[string]any{"path": filepath.ToSlash(f)}),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("174-v1 answered-card leg: out-of-root fs.read isError=%v level=%s bytes=%d",
		out.IsError, out.RiskLevel, len(out.Text))
	if !out.IsError && out.RiskLevel == "L2" && out.Text == body {
		t.Log("CONCLUSION: with an answered L2 card the path IS readable, so the notice's " +
			"\"会被拒/要用户先加 allowed_dirs\" is stricter than the machine (a gated road, not a closed one)")
	}
}
