//go:build !windows

package main

// dllFileVersion has no implementation off Windows: it reads a PE
// VS_FIXEDFILEINFO through version.dll, which does not exist here.
//
// Reporting ("", false) rather than a synthesized version is deliberate -
// doctor.go maps !ok to fail("onnxruntime.dll colocated", ...), and a non-
// Windows tree genuinely has no colocated onnxruntime.dll. A stub that
// returned a version would turn "cannot check" into a false PASS.
//
// This file exists only so the untagged doctor.go type-checks when GOOS is
// not windows; it changes no Windows behaviour.
func dllFileVersion(string) (string, bool) { return "", false }
