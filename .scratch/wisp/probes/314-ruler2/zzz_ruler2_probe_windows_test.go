//go:build windows

package main

// Ticket 314, post-audit repair: one COMPILABLE probe file holding every shape
// the 2026-11-14 non-implementer audit found silent (docs/evidence/s1/
// 314-ruler-audit.md sections 1-3). It is copied into two out-of-repo export
// trees - the pre-repair ruler and the post-repair ruler - so the reading is a
// real mutant going red, not a matcher sample going red.
//
// Shape map:
//
//	A2a zzzRuler2D1a   two literal halves inside one Join
//	A2b zzzRuler2D1b   the whole path by + concatenation
//	A2c zzzRuler2D1c   leaf by strings.Join
//	A2d zzzRuler2D1d   Join(...) plus a literal suffix
//	A2e zzzRuler2D1e   concatenated value bound to a local, then to the sink
//	A3a zzzRuler2D2a   root os.Getenv("TEMP")
//	A3a zzzRuler2D2a2  root os.Getenv("TMP")
//	A3b zzzRuler2D2b   root os.UserHomeDir() through its two results
//	A3c zzzRuler2D2c   root is a drive-letter absolute literal
//	A3c zzzRuler2D2c2  drive-letter absolute literal joined with a leaf
//	X1  zzzRuler2PkgSink  package-level DataPath: key
//	X2  zzzRuler2PkgJoin  package-level fixed join handed to a sink in a func
//
// The last declaration is a control in the other direction: a package-level
// fixed join that never reaches a sink must stay silent.

import (
	"os"
	"path/filepath"
	"strings"
)

func zzzRuler2D1a() string {
	return filepath.Join(os.TempDir(), "wisp-audit2a"+"-panel-profile")
}

func zzzRuler2D1b() string {
	return os.TempDir() + "/wisp-audit2b" + "-panel-profile"
}

func zzzRuler2D1c() string {
	return filepath.Join(os.TempDir(), strings.Join([]string{"wisp-audit2c", "-panel-profile"}, ""))
}

func zzzRuler2D1d() string {
	return filepath.Join(os.TempDir(), "wisp-audit2d") + "-panel-profile"
}

func zzzRuler2D1e() *PanelManager {
	p := filepath.Join(os.TempDir(), "wisp-audit2e") + "-panel-profile"
	return NewPanelManager(nil, nil, p)
}

func zzzRuler2D2a() string {
	return filepath.Join(os.Getenv("TEMP"), "wisp-audit3a-panel-profile")
}

func zzzRuler2D2a2() string {
	return filepath.Join(os.Getenv("TMP"), "wisp-audit3a2-panel-profile")
}

func zzzRuler2D2b() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "wisp-audit3b-panel-profile")
}

func zzzRuler2D2c() *PanelManager {
	return NewPanelManager(nil, nil, "C:\\Temp\\wisp-audit3c-panel-profile")
}

func zzzRuler2D2c2() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join("C:\\Temp", "wisp-audit3c2-panel-profile"))
}

type zzzRuler2Sink struct {
	DataPath string
}

var zzzRuler2PkgSink = zzzRuler2Sink{DataPath: filepath.Join(os.TempDir(), "wisp-audit-pkg-sink")}

var zzzRuler2PkgJoin = filepath.Join(os.TempDir(), "wisp-audit-pkg-var")

func zzzRuler2D3b() *PanelManager {
	return NewPanelManager(nil, nil, zzzRuler2PkgJoin)
}

var zzzRuler2PkgPlain = filepath.Join(os.TempDir(), "wisp-audit-pkg-plain")

func zzzRuler2D3c() string {
	return zzzRuler2PkgPlain
}
