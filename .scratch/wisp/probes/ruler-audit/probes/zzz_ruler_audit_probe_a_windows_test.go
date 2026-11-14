//go:build windows

package main

// Ticket 314 ruler audit, PROBE SET A. This file exists ONLY in the export
// tree C:\Temp\wisp-ruler-audit-20261114 (from git archive HEAD). It is never
// added to the repo; it exists so the ruler's own scan set sees the three
// escape shapes with real source files, real compilation, and a real red.

import (
	"os"
	"path/filepath"
	"strings"
)

// A1a - the audit's attack 1 verbatim: the fixed leaf is built inside a helper
// and the helper's return value is handed to the DataPath sink.
func auditMkProfile() string {
	return filepath.Join(os.TempDir(), "wisp-audit1a-panel-profile")
}

func auditAttack1a() {
	_ = NewPanelManager(nil, nil, auditMkProfile())
}

// A1b - the same helper-return shape, but the leaf avoids the word "profile":
// the SINK rule is the only rule that could see it.
func auditMkPlain() string {
	return filepath.Join(os.TempDir(), "wisp-audit1b-quicklook")
}

func auditAttack1b() {
	_ = NewPanelManager(nil, nil, auditMkPlain())
}

// A1c - the header's own escape surface 1: the root lives inside the helper,
// the literal is spelled at the call site, so they never meet in one expression.
func auditMkParam(leaf string) string {
	return filepath.Join(os.TempDir(), leaf)
}

func auditAttack1c() {
	_ = NewPanelManager(nil, nil, auditMkParam("wisp-audit1c-panel-profile"))
}

// A2a - the leaf assembled from two half literals inside the same Join.
func auditAttack2a() {
	_ = NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-33r1"+"-panel-profile"))
}

// A2b - the whole path assembled by concatenation, three pieces.
func auditAttack2b() {
	_ = NewPanelManager(nil, nil, os.TempDir()+"/wisp-33r1"+"-panel-profile")
}

// A2c - the leaf assembled by strings.Join.
func auditAttack2c() {
	_ = NewPanelManager(nil, nil, filepath.Join(os.TempDir(),
		strings.Join([]string{"wisp-33r1", "-panel-profile"}, "")))
}

// A3a - %TEMP% through os.Getenv instead of os.TempDir().
func auditAttack3a() {
	_ = NewPanelManager(nil, nil, filepath.Join(os.Getenv("TEMP"), "wisp-audit3a-panel-profile"))
}

// A3b - the root is the home directory.
func auditAttack3b() {
	home, _ := os.UserHomeDir()
	_ = NewPanelManager(nil, nil, filepath.Join(home, "wisp-audit3b-panel-profile"))
}

// A3c - an absolute literal root, no root expression at all.
func auditAttack3c() {
	_ = NewPanelManager(nil, nil, `C:\Temp\wisp-audit3c-panel-profile`)
}

// X1 - bonus, NOT in the header's escape list: a package-level struct literal
// whose DataPath key carries a non-profile fixed leaf. The SINK rule only walks
// function bodies, so the DataPath-key arm cannot see this.
var auditPackageSink = auditSink{DataPath: filepath.Join(os.TempDir(), "wisp-auditpkg-sink")}

type auditSink struct{ DataPath string }
