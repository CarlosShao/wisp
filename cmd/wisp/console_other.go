//go:build !windows

package main

// attachParentConsole is a no-op off Windows.
//
// As of 2026-11-02 the Windows shipping build carries -H=windowsgui in
// scripts/build.ps1's $ldflags (:256; the switch landed on the ticket 244
// line, not ticket 07), so the linked binary is a GUI-subsystem executable
// that starts detached from any console. `wisp run` / `wisp doctor` therefore
// rely on attachParentConsole to reach a parent console before printing
// (SPEC-11 §2.2, SPEC-01 §3). Today its early-return branch fires only in the
// redirected-stdout state (`wisp doctor > out.txt`, a usable stdout handle
// already exists); launched from a parent console it does not fire and the
// ATTACH_PARENT_PROCESS path runs. Should the build chain ever drop the flag,
// console builds would take the early return in the normal terminal state
// too — read this as today's fact, not a permanent property.
//
// That problem does not exist on other platforms, where a process started
// from a terminal already has its standard handles.
//
// This file exists only so the untagged callers (main.go, secret.go,
// slo_windows.go) type-check when GOOS is not windows; it changes no Windows
// behaviour.
func attachParentConsole() {}
