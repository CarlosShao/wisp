// Package main is 160-r1's AC#1 rig: the MINIMAL PRODUCTION SHAPE that leaks,
// measured against the REAL risk.Provenance (not a stand-in like 160-c1's
// shape160.go, which re-implemented both candidate shapes instead of importing
// the engine).
//
// The shape is copied from cmd/wisp/panel_assets.go's detector(): build a
// Provenance, OpenScope(...), Mark(...), hand the detector back — and never
// close. That is the whole program. The question this rig answers is what the
// three layers do about it:
//
//	(a) go build / go vet — do they refuse to compile an open-without-close?
//	(b) the doorbells (G5/G6/G7) — see ../gate-pre-change-HEAD.txt
//	(c) post-hoc — is "which scope is still open" observable at all from
//	    outside the package? LEG 2 tries the only exported routes.
package main

import (
	"fmt"
	"os"

	"github.com/CarlosShao/wisp/internal/risk"
)

const marker = "ZZsecret-token-abcdefgh"

func main() {
	// LEG 1: the production shape, verbatim in kind — open, mark, never close.
	p := risk.NewProvenance(risk.ProvOptions{NoProbe: true, HomeDir: os.TempDir()})
	p.OpenScope("task-leaked") // <-- return value (today: none) discarded outright
	p.Mark("task-leaked", risk.SrcWebFetch, "https://example/a", marker)

	hit, ok := p.Inspect("task-leaked", "notify", map[string]any{"text": marker})
	fmt.Printf("LEG1-open-and-leak   inspectHit=%v srcTool=%q marks=%d\n",
		ok, hit.SrcTool, len(p.ScopeTaints("task-leaked")))

	// LEG 2: can an outsider answer "which scope is open and never closed"?
	// ScopeTaints is the only exported per-scope view, and it needs the id —
	// so a leak is only findable by someone who already knows the name.
	// There is NO exported "list open scopes" at all (see LEG 2b's compile
	// note in the evidence file): the census is the only post-hoc detector.
	fmt.Printf("LEG2-outsight        taints(id known)=%d taints(id guessed)=%d\n",
		len(p.ScopeTaints("task-leaked")), len(p.ScopeTaints("task-other")))

	// LEG 3 lives in ../ac1-closebyid/ as its OWN module, on purpose: after
	// AC#2 that program must stop compiling, and a compile failure there must
	// not blind the two readings above.
}
