//go:build windows

package main

// The shell family's tree seam, Windows half (ticket 163's composition-root
// wiring). proc.OpenJobScope is a //go:build windows declaration and run.go is
// an untagged file, so the reference lives here - the cross-vet rule's own
// recipe (build-tag the referencing file, keep the caller untagged). One
// kill-on-close Job per process tree: closing it kills the shell and every
// descendant, which is what turns "cancel => the whole tree is gone" into an
// OS fact instead of a polite Kill on one pid (ticket 163 AC#3).

import (
	"github.com/CarlosShao/wisp/internal/proc"
	"github.com/CarlosShao/wisp/internal/tools"
)

// newShellSessionTree hands out one tree per one-shot shell.exec call and per
// shell.session carrier. *proc.JobScope satisfies tools.SessionTree by method
// set (StartInJob/Close) - the assertion internal/tools' own windows test
// pins - so this closure is the whole production half.
func newShellSessionTree() (tools.SessionTree, error) {
	return proc.OpenJobScope()
}
