//go:build windows

package main

// Adversarial probe for D3 (package-level seeds), part 2: another file uses the
// package-level fixed join declared in zzz_audit_d3a_decl_windows_test.go.

func zzzD3aUse() *PanelManager {
	return NewPanelManager(nil, nil, zzzD3OtherFileJoin)
}
