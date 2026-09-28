package risk

// Ticket 183 dispatch section 6, the mandatory frequency question: does a REAL
// body of text ever carry one of the artifact path's own >=8-rune windows, or
// was "-output-" a fixture artifact of 179-r2's filler line? The repo had no
// reading for this dimension (183-a1 section 2 named it as unmeasured).
//
// This is a PROBE, mapped into the risk package with
// .scratch/wisp/probes/183/r1/overlay-freq-head.json, which pins
// provenance.go + taintmatch.go to their HEAD blobs (git show HEAD:<path>) so
// every reading below is the UNFIXED production rule answering for itself -
// not a hand-derived predicate.
//
// Run:
//
//	go test -count=1 -overlay=.scratch/wisp/probes/183/r1/overlay-freq-head.json -run TestFreq183 -v ./internal/risk/
//
// Shape per corpus file: the host stub (省略/总长/全文见 <artifacts path>) with
// the real file's text as the tail body, marked through MarkWithHostPath with
// that path declared, then the model's reread Inspected as fs.read{path}. A hit
// means: on unfixed code this real body would have cost the model its own
// reread of the pointer the host wrote for it.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func freqCorpus183(t *testing.T) []string {
	t.Helper()
	root := repoRoot183(t)
	rel := []string{
		"AGENTS.md",
		"docs/PLAN.md",
		"docs/reports/HANDOVER.md",
		"docs/specs/SPEC-06-security-gating.md",
		"docs/specs/SPEC-08-ball-and-panel.md",
		"internal/agent/loop.go",
		"internal/tools/bridge.go",
		"cmd/wisp/run.go",
		"scripts/d22scan.sh",
		"go.mod",
	}
	var out []string
	for _, r := range rel {
		p := filepath.Join(root, filepath.FromSlash(r))
		b, err := os.ReadFile(p)
		if err != nil {
			t.Logf("corpus %s unreadable: %v", r, err)
			continue
		}
		out = append(out, r+"|"+sanitizeBody183(string(b)))
	}
	// A couple of bodies that are not repo text at all: what a real task output
	// of the two most path-shaped kinds looks like.
	out = append(out,
		"synthetic:dir-listing|"+strings.Join([]string{
			"listing of C:\\Users\\swq\\AppData\\Roaming\\wisp\\artifacts follows",
			"tool-output-agent-task-00000000.txt  20040 bytes",
			"C:\\Users\\swq\\AppData\\Roaming\\wisp\\config.toml",
		}, "\n"),
		"synthetic:plain-prose|Wisp is a lightweight desktop assistant. It should feel calm, quick and honest, and never shout at the user.",
	)
	return out
}

func sanitizeBody183(s string) string {
	r := []rune(s)
	if len(r) > 6000 {
		r = r[:6000]
	}
	return string(r)
}

func repoRoot183(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			t.Fatal("no go.mod above the test working directory")
		}
		wd = parent
	}
}

func TestFreq183RealBodiesAgainstTheUnfixedRule(t *testing.T) {
	path := `C:\Users\swq\AppData\Roaming\wisp\artifacts\tool-output-agent-task-c66c0634.txt`
	np := []rune(normalizeTaint(path))
	nset := map[string]bool{}
	for i := 0; i+contractMinFragmentChars <= len(np); i++ {
		nset[string(np[i:i+contractMinFragmentChars])] = true
	}
	hits, total := 0, 0
	var lines []string
	for _, c := range freqCorpus183(t) {
		name, body := c, ""
		if i := strings.Index(c, "|"); i >= 0 {
			name, body = c[:i], c[i+1:]
		}
		total++
		stub := "报告头部 中段 […输出已落文件：省略 17200 字符，总长 20000 字节，全文见 " + path +
			"…]\n尾部段落（真实正文如下）\n" + body
		p := testProv(t, baseOptions(t))
		if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", stub, path) {
			t.Fatalf("corpus %s: mark refused", name)
		}
		h, ok := p.Inspect("task-1", "fs.read", map[string]any{"path": path})
		frag := ""
		if ok {
			hits++
			// Which of the path's own windows the real body spells outside the
			// span - reported, not assumed.
			nb := []rune(normalizeTaint(stub))
			spanLo := runeIndexOf(nb, np)
			seen := map[string]bool{}
			for i := 0; i+contractMinFragmentChars <= len(nb); i++ {
				if i < spanLo+contractMinFragmentChars && i+contractMinFragmentChars > spanLo {
					continue // window touches the declared span: excluded at HEAD too
				}
				if w := string(nb[i : i+contractMinFragmentChars]); nset[w] && !seen[w] {
					seen[w] = true
					if len(frag) < 120 {
						frag += " " + w
					}
				}
			}
		}
		lines = append(lines, fmt.Sprintf("%-34s hit=%-5v %s", name, ok, strings.TrimPrefix(frag, " ")))
		_ = h
	}
	t.Logf("REAL-BODY FREQUENCY on UNFIXED code: %d of %d bodies make the host's own pointer unreadable", hits, total)
	for _, l := range lines {
		t.Logf("  %s", l)
	}
}
