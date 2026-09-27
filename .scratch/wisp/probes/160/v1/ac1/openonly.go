// 160-v1 acceptance rig, LEG 1/2: my OWN open-without-close sample (not the
// implementer's probes/160/r1/ac1). Statement-form OpenScope compiles on BOTH
// APIs (old: no result; new: discarded result) -> run once pre-change (old
// overlay) and once post-change (working tree) to read both worlds.
package main

import (
	"fmt"
	"os"

	"github.com/CarlosShao/wisp/internal/risk"
)

const marker = "VVaccept160-openonly-marker"

func main() {
	p := risk.NewProvenance(risk.ProvOptions{NoProbe: true, HomeDir: os.TempDir()})
	p.OpenScope("v1-task-leak")
	p.Mark("v1-task-leak", risk.SrcWebFetch, "https://v1/leak", marker)
	hit, ok := p.Inspect("v1-task-leak", "notify", map[string]any{"text": "echo " + marker})
	fmt.Printf("V1-OPENONLY leak-inspect hit=%v src=%q taints=%d guessed-id-taints=%d\n",
		ok, hit.SrcTool, len(p.ScopeTaints("v1-task-leak")), len(p.ScopeTaints("v1-no-such-scope")))
}
