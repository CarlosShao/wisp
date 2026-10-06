package main

// 231-v1 (adversarial acceptance leg) · AC#1 re-implementation of the read-only
// probe. Physical file lives in .scratch/wisp/probes/231/v1/ and is mapped into
// cmd/wisp only through `go test -overlay`; the worktree never contains it.
//
// Same host as ticket 223's permanent tests (newReloadRun223 -> live -> writeOver
// -> the installed 1s tick), same planted shape (schema_version above this build
// with a body that DOES parse, the only shape that reaches loader.go's
// newer-build return instead of being short-circuited by the parse check).
//
// It asserts nothing on purpose: the verdict comes from reading the operator's
// own line off the audit trail, not from a colour.
//
// PATH (ticket 98): cmd/wisp's test binary links sherpa-onnx; without
// third_party/sherpa-onnx and build on PATH it dies with 0xc0000135 and zero
// --- FAIL, which reads as "green" if you do not look.

import (
	"strings"
	"testing"
	"time"
)

func TestProbe231V1NewerBuildOperatorLine(t *testing.T) {
	const body = "schema_version = 99\n\n[ball]\nsize = 64\n"
	r := newReloadRun223(t, 40*time.Second)
	r.live(t, func() {
		r.awaitAudit(t, "config: HOT-RELOAD state=armed")
		r.writeOver(t, body)
		trail := r.awaitAudit(t, "config: HOT-RELOAD state=not-applied cause=")
		for _, line := range strings.Split(trail, "\n") {
			if strings.Contains(line, "state=not-applied") {
				t.Logf("AC#1 OPERATOR LINE (verbatim) =\n%s", line)
			}
		}
		t.Logf("AC#1 FULL TRAIL (stderr) =\n%s", trail)
	})
}
