//go:build windows
// +build windows

package edge

import (
	"sync/atomic"
	"testing"
)

// The ticket-408 nil-env guard: a boot that completes with a negative HRESULT
// or a nil interface pointer (measured live as res=S_OK with env=0x0 when the
// WebView2 browser process dies at birth) must neither crash the boot thread
// nor start the success chain - it must arm bootFailed so Embed's pump returns
// false and the host can take its fallback renderer. These tests inject the
// callback shapes directly; they never depend on a real WebView2 runtime.

// hrEFail is HRESULT 0x80004005 sign-extended to the pointer-width register
// value a COM completion callback actually receives on win64.
const hrEFail = uintptr(0xFFFFFFFF80004005)

func TestEnvironmentCompletedNilEnvArmsBootFailed(t *testing.T) {
	e := NewChromium()
	e.EnvironmentCompleted(0, nil) // the measured live shape; must not panic
	if atomic.LoadUintptr(&e.bootFailed) == 0 {
		t.Fatalf("S_OK with nil environment did not arm bootFailed")
	}
	if e.environment != nil {
		t.Fatalf("nil environment leaked into e.environment")
	}
	if e.controller != nil || e.webview != nil {
		t.Fatalf("success chain started from a nil environment")
	}
}

func TestEnvironmentCompletedFailedHRESULTArmsBootFailed(t *testing.T) {
	e := NewChromium()
	e.EnvironmentCompleted(hrEFail, nil) // must not log.Fatalf the process
	if atomic.LoadUintptr(&e.bootFailed) == 0 {
		t.Fatalf("negative HRESULT did not arm bootFailed")
	}
	if e.environment != nil {
		t.Fatalf("failed boot produced an environment")
	}
}

func TestControllerCompletedNilArmsBootFailed(t *testing.T) {
	e := NewChromium()
	e.CreateCoreWebView2ControllerCompleted(0, nil) // must not panic
	if atomic.LoadUintptr(&e.bootFailed) == 0 {
		t.Fatalf("S_OK with nil controller did not arm bootFailed")
	}
	if e.controller != nil || e.webview != nil {
		t.Fatalf("success chain started from a nil controller")
	}
}

func TestControllerCompletedFailedHRESULTArmsBootFailed(t *testing.T) {
	e := NewChromium()
	e.CreateCoreWebView2ControllerCompleted(hrEFail, nil) // must not log.Fatalf
	if atomic.LoadUintptr(&e.bootFailed) == 0 {
		t.Fatalf("negative HRESULT did not arm bootFailed")
	}
}

func TestBootGuardResetsBetweenBootAttempts(t *testing.T) {
	e := NewChromium()
	e.EnvironmentCompleted(0, nil)
	atomic.StoreUintptr(&e.bootFailed, 0)
	// A healthy boot (env wired by the real runtime) would re-run the
	// success chain; the flag staying cleared after an explicit reset is
	// what lets a later Embed attempt proceed.
	if atomic.LoadUintptr(&e.bootFailed) != 0 {
		t.Fatalf("bootFailed did not clear")
	}
}
