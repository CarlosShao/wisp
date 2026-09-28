// Package streamkey owns the one spelling of the stream keys this harness
// writes and reads (ticket 197 §0, ruling A406).
//
// WHY A PACKAGE FOR ONE CONSTANT. A subagent's work page is addressed by a
// stream key, and two packages have to agree on it: the writer is
// internal/tools (task.spawn feeds the child's lifecycle lines and its sink
// feeds its deltas) and the reader is internal/panel (StreamLog keys its rows
// by it and the snapshot's results section carries them). Those two must not
// import each other - the tool layer has no business knowing the panel's view
// model, and the panel has no business knowing the bridge - so before this
// package each of them carried its OWN copy of the literal and only the
// composition root could compare them (cmd/wisp/subagent_stream_key_197_test.go).
//
// Two copies of one spelling is how a display starts lying about provenance:
// change it on the writing side and every roster row keeps its state while the
// page for it goes empty, and no assertion in either package can see that,
// because each one is self-consistent. So the spelling now lives here and both
// sides alias it; internal/streamkey is the only place in this repository where
// the literal itself may be written, and TestSubagentStreamKeyHasOneMintSite
// (cmd/wisp) is the resident nail that goes red when somebody re-forks it.
package streamkey

import "strings"

// SubagentPrefix is ticket 197 §0's frozen shape: lower-case, colon separated,
// and the task id is passed through verbatim after it. Nothing about a task id
// is transformed here - a key that renamed its own agent would be a second
// identity for the same row.
const SubagentPrefix = "subagent:"

// Subagent returns the stream key for one subagent task id.
//
// An empty or whitespace-only id yields "": a stream belonging to no task is
// not a stream, and appending under "" would open a row the panel cannot
// attribute to anybody. That guard used to live only on the reading side
// (panel.SubagentStreamKey) while the writing side had no such branch, which
// is the divergence ledger A406 points at; aliasing this function closes it for
// both sides at once.
func Subagent(taskID string) string {
	if strings.TrimSpace(taskID) == "" {
		return ""
	}
	return SubagentPrefix + taskID
}
