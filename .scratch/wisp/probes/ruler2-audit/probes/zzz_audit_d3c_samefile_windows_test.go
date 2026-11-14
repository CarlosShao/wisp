//go:build windows

package main

// D3 contrast control: the same package-level fixed join and sink, but in the
// SAME file. The repair seeds only same-file package bindings, so this one must
// ring while the cross-file / init() / sync.Once versions stay silent.

import (
	"os"
	"path/filepath"
)

var zzzD3SameFileJoin = filepath.Join(os.TempDir(), "wisp-d3x-quicklook")

func zzzD3SameFileUse() *PanelManager {
	return NewPanelManager(nil, nil, zzzD3SameFileJoin)
}
