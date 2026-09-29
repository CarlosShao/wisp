//go:build windows

package main

// The console's answer stream, and when it is really there (ticket 201).
//
// interactiveStdin is the gate on the whole reply listener, so its rule is worth
// stating exactly: it hands back os.Stdin only when the standard input handle is
// a console screen-buffer input, i.e. when a human is sitting at a terminal that
// this process owns. Everything else - a pipe, a redirected file, Explorer's
// detached start with no console at all, CI - gets nil, which means "no answer
// source" and therefore the pre-201 posture: the card is shown, nobody answers
// it, and each route resolves on its own clock.
//
// Why the distinction is load-bearing rather than tidy:
//
//   - A redirected stdin is somebody else's bytes. An approval answer is the
//     one input in this process that can authorise a write to disk, so a channel
//     a script can pre-fill is not a "user confirmed it" - it is a way to make
//     the confirmation arrive without a user. Reading a piped stdin for an
//     approval reply would be the caller-controlled-selector failure this
//     repository keeps naming (M-7/C-3) wearing a keyboard.
//   - A console handle that GetConsoleMode accepts is, conversely, evidence of a
//     real interactive desktop input buffer: the same call the ball's own console
//     code leans on, and it fails for every non-console file type.
//
// The trade-off, written down because it is a cost and not a bug: a run whose
// stdin is a pipe cannot be answered, and says so at startup rather than going
// quiet. That is the same shape as every other degraded boot in this file set
// (SPEC-05 §3.4 - the downgrade has to be visible).

import (
	"io"
	"os"

	"golang.org/x/sys/windows"
)

// interactiveStdin returns the operator's console input stream, or nil when this
// process has no interactive console to read answers from.
func interactiveStdin() io.Reader {
	h, err := windows.GetStdHandle(windows.STD_INPUT_HANDLE)
	if err != nil || h == 0 || h == windows.InvalidHandle {
		return nil
	}
	var mode uint32
	if err := windows.GetConsoleMode(h, &mode); err != nil {
		// Not a console input buffer: a pipe, a file, or no handle at all.
		return nil
	}
	return os.Stdin
}
