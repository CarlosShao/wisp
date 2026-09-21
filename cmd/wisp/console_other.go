//go:build !windows

package main

// attachParentConsole is a no-op off Windows.
//
// The Windows build links as a GUI-subsystem binary (ticket 07), which starts
// with no console at all, so `wisp run` / `wisp doctor` must attach to the
// parent console before printing (SPEC-11 §2.2, SPEC-01 §3). That problem does
// not exist on other platforms, where a process started from a terminal
// already has its standard handles.
//
// This file exists only so the untagged callers (main.go, secret.go,
// slo_windows.go) type-check when GOOS is not windows; it changes no Windows
// behaviour.
func attachParentConsole() {}
