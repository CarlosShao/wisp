//go:build windows

package main

// Ticket 442 AC#2's probe, kept out of the repo's scan set on purpose: it is
// copied into an out-of-repo export tree so the ruler reads it as a normal
// cmd/wisp source, and `go test -c` proves it compiles.
//
// The three shapes below are the two the 2026-11-14 audit called A1b and A1c,
// plus A1c's mirror direction. In each one the temp root and the constant leaf
// are written in DIFFERENT functions, which is exactly what a single-function
// fold cannot see: the value crosses a call.
//
//   - A1b: the helper returns a leaf that avoids the word "profile", so the NAME
//     rule (which needs the spelling) cannot see it; only the SINK rule can, and
//     only if it follows the return value across `plainLeaf442Probe()`.
//   - A1c(root): the root is the helper's parameter and the caller supplies it;
//     the constant leaf is inside the helper.
//   - A1c(leaf): the audit's own reading of the same class - root inside the
//     helper, the constant leaf supplied by the caller.
//
// Every sink call is in a function nothing calls. That is deliberate: the ruler
// judges source shapes, not reachability, and a shape only a dead function holds
// is still a shape the next refactor will move into a live one.

import (
	"os"
	"path/filepath"
)

// A1b - a plain-leaved helper return handed to the sink.
func plainLeaf442Probe() string {
	return filepath.Join(os.TempDir(), "wisp-442-a1b-quicklook")
}

func a1bProbe() {
	_ = NewPanelManager(nil, nil, plainLeaf442Probe())
}

// A1c(root) - ticket 442's reading: the root arrives as a parameter.
func paramRoot442Probe(root string) string {
	return filepath.Join(root, "wisp-442-a1c-rootparam")
}

func a1cRootProbe() {
	_ = NewPanelManager(nil, nil, paramRoot442Probe(os.TempDir()))
}

// A1c(leaf) - the audit's reading: the leaf arrives as a parameter.
func paramLeaf442Probe(leaf string) string {
	return filepath.Join(os.TempDir(), leaf)
}

func a1cLeafProbe() {
	_ = NewPanelManager(nil, nil, paramLeaf442Probe("wisp-442-a1c-leafparam"))
}
