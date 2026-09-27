// 160-v1 acceptance rig, foreign-by-id close: my OWN shot (the implementer's
// probes/160/r1/ac1-closebyid is positive-control context only). Compiles ONLY
// against the pre-change tree (p.CloseScope disappeared in f576cf08).
// SHOT A = silent shape claimed by 160-r1 §3(a); SHOT B = control where another
// scope still holds taints, where the fail-closed log MUST show. Two-direction
// test of the [inference -> reading] upgrade claim about scopeMarks' last branch.
package main

import (
	"fmt"
	"os"

	"github.com/CarlosShao/wisp/internal/risk"
)

const marker = "VVaccept160-foreign-close-marker"

func foreignClose(p *risk.Provenance, id string) { p.CloseScope(id) }

func main() {
	p := risk.NewProvenance(risk.ProvOptions{NoProbe: true, HomeDir: os.TempDir()})
	p.OpenScope("v1-victim-a")
	p.Mark("v1-victim-a", risk.SrcWebFetch, "https://v1/a", marker)
	beforeA := len(p.ScopeTaints("v1-victim-a"))
	foreignClose(p, "v1-victim-a")
	_, hitA := p.Inspect("v1-victim-a", "notify", map[string]any{"text": marker})
	fmt.Printf("V1-FOREIGN-A no-other-taints before=%d after=%d inspectHit=%v\n",
		beforeA, len(p.ScopeTaints("v1-victim-a")), hitA)

	q := risk.NewProvenance(risk.ProvOptions{NoProbe: true, HomeDir: os.TempDir()})
	q.OpenScope("v1-victim-b")
	q.Mark("v1-victim-b", risk.SrcWebFetch, "https://v1/b", marker)
	q.OpenScope("v1-witness")
	q.Mark("v1-witness", risk.SrcFSRead, "/witness", marker)
	foreignClose(q, "v1-victim-b")
	hitB, okB := q.Inspect("v1-victim-b", "notify", map[string]any{"text": "zz"})
	fmt.Printf("V1-FOREIGN-B witness-tainted inspectHit=%v src=%q\n", okB, hitB.SrcTool)
}
