//go:build !windows

package main

// interactiveStdin has no answer stream to hand out off Windows.
//
// The product's platform scope is Windows (D7), and the console test this file's
// Windows counterpart leans on - GetConsoleMode on the standard input handle,
// which distinguishes a terminal's input buffer from a pipe - has no portable
// equivalent that this leg is willing to assert. Returning nil keeps the
// fail-closed reading: no answer source, so a card on a non-Windows build stays
// exactly as unanswered as it was before ticket 201, and the CLI tests that DO
// answer cards pass their own reader through runSpec.reply (the `wisp run` seam
// AGENTS.md §1.3 names) instead of relying on this function.
//
// This file exists so the untagged callers (main.go, approval_reply.go) type-check
// under GOOS=linux, the same job console_other.go does for attachParentConsole.
// It changes no Windows behaviour.

import "io"

// interactiveStdin always reports "no interactive console".
func interactiveStdin() io.Reader { return nil }
