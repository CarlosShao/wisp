//go:build windows

package main

// Adversarial probe for the repaired ticket-314 shape ruler, audit leg 2.
// Every shape here is a fixed-name WebView2 data-path folder spelled in a way
// the repair claims to fold or to recognise. Copied into an out-of-repo export
// tree; never compiled into the repo.

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// --- D1: leaves that are fully constant but not written as one literal -----

// Same-file const: the whole leaf is constant, just named.
const zzzD1SuffixSameFile = "-pan" + "el-profile"

func zzzD1a() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-d1a"+zzzD1SuffixSameFile))
}

func zzzD1b() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-d1b"+zzzD1SuffixOtherFile))
}

func zzzD1c() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(),
		fmt.Sprintf("%s%s", "wisp-d1c", "-panel-profile")))
}

func zzzD1d() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(),
		strings.Replace("wisp-d1d_PANEL_profile", "_", "-", -1)))
}

func zzzD1e() *PanelManager {
	parts := []string{"wisp-d1e", "-panel-profile"}
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), strings.Join(parts, "")))
}

func zzzD1f() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(append([]string{os.TempDir()}, "wisp-d1f", "-panel-profile")...))
}

func zzzD1g() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-d1g", "panel-profile"))
}

func zzzD1h() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(),
		filepath.Join("sub", "wisp-d1h-panel-profile")))
}

func zzzD1i() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-d1i")+filepath.Join("", "-panel-profile"))
}

// --- D2: roots the repair does not name -----------------------------------

func zzzD2a() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.Getenv("TMPDIR"), "wisp-d2a-panel-profile"))
}

func zzzD2b() *PanelManager {
	const envName = "TEMP"
	return NewPanelManager(nil, nil, filepath.Join(os.Getenv(envName), "wisp-d2b-panel-profile"))
}

func zzzD2c() *PanelManager {
	exe, _ := os.Executable()
	return NewPanelManager(nil, nil, filepath.Join(filepath.Dir(exe), "wisp-d2c-panel-profile"))
}

type zzzD2Dirs struct{}

func (zzzD2Dirs) TempDir() string { return "" }

func zzzD2d() *PanelManager {
	d := zzzD2Dirs{}
	return NewPanelManager(nil, nil, filepath.Join(d.TempDir(), "wisp-d2d-panel-profile"))
}

func zzzD2e() *PanelManager {
	dir, _ := os.MkdirTemp("", "wisp-d2e")
	return NewPanelManager(nil, nil, filepath.Join(dir, "wisp-d2e-panel-profile"))
}

func zzzD2f() *PanelManager {
	return NewPanelManager(nil, nil, path.Join(os.TempDir(), "wisp-d2f-panel-profile"))
}

func zzzD2g() *PanelManager {
	root := os.Getenv("TEMP")
	return NewPanelManager(nil, nil, filepath.Join(root, "wisp-d2g-panel-profile"))
}

func zzzD2h() *PanelManager {
	root := os.Getenv("TEMP") + string(filepath.Separator) + "wisp"
	return NewPanelManager(nil, nil, filepath.Join(root, "wisp-d2h-panel-profile"))
}
