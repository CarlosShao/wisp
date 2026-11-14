//go:build windows

package main

// Adversarial probe for D3 (package-level seeds), part 1: declarations.
// File A holds a package-level fixed join, an init()-assigned fixed join and a
// sync.Once-assigned fixed join. Only the plain one is used from another file.

import (
	"os"
	"path/filepath"
	"sync"
)

var zzzD3OtherFileJoin = filepath.Join(os.TempDir(), "wisp-d3a-quicklook")

var zzzD3InitValue string

func init() {
	zzzD3InitValue = filepath.Join(os.TempDir(), "wisp-d3b-quicklook")
}

var zzzD3Once sync.Once

var zzzD3OncePath string

func zzzD3InitOnce() {
	zzzD3Once.Do(func() {
		zzzD3OncePath = filepath.Join(os.TempDir(), "wisp-d3c-quicklook")
	})
}

func zzzD3bUse() {
	zzzD3InitOnce()
	_ = NewPanelManager(nil, nil, zzzD3InitValue)
}

func zzzD3cUse() {
	zzzD3InitOnce()
	_ = NewPanelManager(nil, nil, zzzD3OncePath)
}
