// Package main is 160-r1's AC#1 LEG 3: the by-id close route, called from a
// package that never opened the scope. This is the exact shape
// internal/tools/bridge.go:695 uses today.
//
// PRE-CHANGE reading: this compiles, vets clean, and RUNS — the foreign close
// succeeds and the victim's taints vanish with no error and no log. That is the
// hole AC#2 criterion (iii) has to shut.
//
// POST-CHANGE reading: this must not compile. The compile error is the
// deliverable; `go build` output is captured verbatim into the evidence file.
package main

import (
	"fmt"
	"os"

	"github.com/CarlosShao/wisp/internal/risk"
)

const marker = "ZZsecret-token-abcdefgh"

func main() {
	p := risk.NewProvenance(risk.ProvOptions{NoProbe: true, HomeDir: os.TempDir()})
	// The victim opens and taints its own scope.
	p.OpenScope("task-victim")
	p.Mark("task-victim", risk.SrcWebFetch, "https://example/victim", marker)
	before := len(p.ScopeTaints("task-victim"))

	// An unrelated owner — it never opened "task-victim" — names the id and
	// closes it. Nothing about this call is distinguishable from the victim
	// closing its own scope.
	p.CloseScope("task-victim")

	_, hit := p.Inspect("task-victim", "notify", map[string]any{"text": marker})
	fmt.Printf("foreign-by-id-close taintsBefore=%d taintsAfter=%d inspectHitAfter=%v\n",
		before, len(p.ScopeTaints("task-victim")), hit)
}
