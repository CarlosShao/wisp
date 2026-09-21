package risk

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/winsec"
)

// Ticket 102 AC#1: the reproduction of acceptor-ticket94's PROBE B2
// (docs/evidence/s1/94-adversarial-acceptance.md:262-291, residual R-a).
//
// expandInput (pathresolver.go) runs env / ~ expansion as the FIRST step of the
// C26 pipeline — that step is contract-mandated (SPEC-06 §4 line 50:
// 「展开(env / ~) → 绝对化 → Clean → ...」, mirrored by PLAN.md:2375), so this test
// does NOT ask for it to be removed. What it asks for is the thing that is
// missing today: when expansion moves the path onto a different tree, no
// production decision leg may act on that other tree and hand its caller a
// plain success. The invariant asserted below is exactly the ticket's wording:
//
//	"结果树与请求树必须同一棵，否则必须报错"
//
// "请求树" (the requested tree) is what the caller itself computes: the lexical
// absolute form of the spelling it was holding, without any expansion. That is
// literally what internal/memory/open.go:181-183 does after winsec returns nil,
// and what every sealing/judging leg compares against afterwards. lexCanonical
// is the pipeline's own lexical step, so calling it here adds no second
// normalizer (and no D22 surface: the Abs/Clean stay inside pathresolver.go).
//
// Each leg below is a real production entry point, not a probe:
//  1. c26Pipeline.Resolve  - the seam internal/risk installs into winsec
//     (winsec_c26.go:33), i.e. the exact function PROBE B2 measured.
//  2. winsec.ResolvePath   - what the sealing caller gets back (nil = "sealed").
//  3. syncSet.resolveTarget - the fs.write "is this a sync dir" judgement.

// requireSameTreeOrRefusal is the AC#1 assertion. A leg passes when it either
// refuses (non-nil error) or hands back a spelling of the very tree the caller
// named. Acting on a different tree with err == nil is the fail-open shape.
func requireSameTreeOrRefusal(t *testing.T, leg, input string, got string, err error) {
	t.Helper()
	if err != nil {
		t.Logf("%s: refused as required: %v", leg, err)
		return
	}
	requested := lexCanonical(input)
	if normPath(got) == normPath(requested) {
		t.Logf("%s: same tree as requested (pct input=%d got=%d requested=%d) got=%s",
			leg, strings.Count(input, "%"), strings.Count(got, "%"), strings.Count(requested, "%"), got)
		return
	}
	t.Errorf("FAIL-OPEN: %s acted on a tree the caller did not name and reported success (err = nil)\n"+
		"  caller's spelling : %s\n"+
		"  requested tree    : %s\n"+
		"  tree acted upon   : %s\n"+
		"  invariant         : 结果树与请求树必须同一棵，否则必须报错",
		leg, input, requested, got)
}

func TestC26ExpansionMustNotRewriteOntoAnotherTree(t *testing.T) {
	t.Run("percent_env_var_spelling", func(t *testing.T) {
		const name = "WISP102PROBE_TARGET"
		t.Setenv(name, "elsewhere")
		root := t.TempDir()
		// A directory whose literal name contains %VAR% is legal on NTFS: this is
		// PROBE B2's shape, where the caller means THIS tree and C26 answers for
		// another one.
		callerTree := filepath.Join(root, "a%"+name+"%b")
		if err := os.MkdirAll(callerTree, 0o700); err != nil {
			t.Fatalf("mkdir caller tree: %v", err)
		}
		input := callerTree + sepStr + "artifacts"
		if _, err := os.Stat(input); err == nil {
			t.Fatalf("probe tree already exists: %s", input)
		}

		c26 := c26Pipeline{}
		got, err := c26.Resolve(input)
		requireSameTreeOrRefusal(t, "risk.c26Pipeline.Resolve (seam installed into winsec)", input, got, err)

		got2, err2 := winsec.ResolvePath(input)
		requireSameTreeOrRefusal(t, "winsec.ResolvePath (sealing caller's view)", input, got2.String(), err2)

		got3, err3 := (&syncSet{}).resolveTarget(input)
		requireSameTreeOrRefusal(t, "risk.syncSet.resolveTarget (fs.write sync judgement)", input, got3, err3)
	})

	t.Run("leading_tilde_spelling", func(t *testing.T) {
		// Nothing is created under the real home directory: every leg below is
		// read-only, the assertion is on which tree each one names.
		for _, input := range []string{`~\wisp102-tilde-probe\artifacts`, `~/wisp102-tilde-probe/artifacts`} {
			c26 := c26Pipeline{}
			got, err := c26.Resolve(input)
			requireSameTreeOrRefusal(t, "risk.c26Pipeline.Resolve ~ "+input, input, got, err)

			got2, err2 := winsec.ResolvePath(input)
			requireSameTreeOrRefusal(t, "winsec.ResolvePath ~ "+input, input, got2.String(), err2)

			got3, err3 := (&syncSet{}).resolveTarget(input)
			requireSameTreeOrRefusal(t, "risk.syncSet.resolveTarget ~ "+input, input, got3, err3)
		}
	})
}
