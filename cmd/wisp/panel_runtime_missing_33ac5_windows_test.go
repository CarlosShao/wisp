//go:build windows

// Ticket 33 AC#5's fixture: "Runtime-missing fixture (rename/mask loader) ->
// fallback event, app alive, native L2 flag on."
//
// WHY THE MASK LANDS ON THE PROBE SEAM AND NOT ON A DLL FILE. The fork embeds
// WebView2Loader.dll into the binary (webviewloader/module_amd64.go) and falls
// back to loading it from memory when the disk copy is missing, so renaming or
// hiding the DLL cannot produce the missing-runtime world on a machine that
// HAS the Runtime installed. The dimension that actually distinguishes the
// world is the loader's own version query
// (GetAvailableCoreWebView2BrowserVersionString): it answers ("", nil) - no
// error, no install - when the Runtime is absent. So the fixture swaps
// panelRuntimeVersion (cmd/wisp/panel_host_windows.go), the named seam whose
// production value IS that real query, and the production path is exercised
// for real: bringUp probes BEFORE newPanelControl, refuses by name, emits the
// panel.unavailable event and sets the L2 native fallback flag without ever
// reaching the library's create path (whose log.Fatal call sites are why the
// probe exists at all - A500's correction of this ticket's J2 answer).
//
// The flag's named consumer is ticket 37's masked-runtime fallback test
// ("with panel runtime masked, L2 still approvable via native card (33
// fixture)"); until that card lands, this file is the flag's reader, which is
// what keeps it from being a write-only field.
package main

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
)

// withMaskedRuntime swaps the probe seam for one that answers the
// missing-runtime world the way the real loader does: ("", nil) - no error,
// no install. The real value is restored by the returned cleanup.
func withMaskedRuntime(t *testing.T) {
	t.Helper()
	orig := panelRuntimeVersion
	panelRuntimeVersion = func() (string, error) { return "", nil }
	t.Cleanup(func() { panelRuntimeVersion = orig })
}

// withCapturedSlog returns the buffer the default slog handler writes to while
// the test runs. Package tests are sequential (zero t.Parallel in cmd/wisp),
// so the swap cannot interleave with another test's logging.
func withCapturedSlog(t *testing.T) *bytes.Buffer {
	t.Helper()
	buf := &bytes.Buffer{}
	orig := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(buf, nil)))
	t.Cleanup(func() { slog.SetDefault(orig) })
	return buf
}

// TestAC5RuntimeMissingRefusesNamedSetsTheFallbackFlag33ac5 is the fixture
// itself: mask the Runtime, ask for the panel, and demand the three AC#5
// products - a named (never fatal) refusal, the panel.unavailable event, and
// the native L2 fallback flag ON, with the window never created. Reaching the
// assertions at all is the "app alive" half: the pre-probe shape found a
// missing runtime by walking into the library's create path, whose failure
// modes on this world include log.Fatal.
func TestAC5RuntimeMissingRefusesNamedSetsTheFallbackFlag33ac5(t *testing.T) {
	withMaskedRuntime(t)
	logs := withCapturedSlog(t)

	mgr := NewPanelManager(&panel.ComposerDispatch{}, nil, filepath.Join(t.TempDir(), "panel-webview2"))

	err := mgr.Show(context.Background())
	if !errors.Is(err, errPanelRuntimeMissing) {
		t.Fatalf("a masked runtime must refuse by name (errPanelRuntimeMissing), got: %v", err)
	}
	if !strings.Contains(err.Error(), "Evergreen") {
		t.Errorf("the refusal must carry the Evergreen guidance, never an auto-download: %v", err)
	}
	if mgr.IsCreated() {
		t.Errorf("the probe must refuse BEFORE the create: the window exists in the missing-runtime world")
	}
	on, reason := mgr.Unavailable()
	if !on {
		t.Errorf("AC#5's native L2 fallback flag must be ON after the refusal (reason recorded: %q)", reason)
	}
	if strings.TrimSpace(reason) == "" {
		t.Errorf("the flag's reason must be recorded, not an empty string")
	}
	if !strings.Contains(logs.String(), "panel.unavailable") {
		t.Errorf("the panel.unavailable event was not emitted, log buffer: %q", logs.String())
	}
}

// TestAC5FreshManagerFlagStartsOff33ac5 is the negative control: the flag is a
// state read, not a promise - a manager that never refused must read
// (false, ""), so a consumer cannot mistake "nothing happened yet" for
// "runtime missing".
func TestAC5FreshManagerFlagStartsOff33ac5(t *testing.T) {
	mgr := NewPanelManager(&panel.ComposerDispatch{}, nil, filepath.Join(t.TempDir(), "panel-webview2"))
	if on, reason := mgr.Unavailable(); on || reason != "" {
		t.Fatalf("a fresh manager must read (false, \"\"), got (%t, %q)", on, reason)
	}
}

// TestAC5ProbeMapsEveryMissingAnswerToTheNamedRefusal33ac5 pins the probe's
// own contract at the seam: a loader error and an empty version both refuse by
// name (fail-closed - the create path they guard has log.Fatal call sites),
// while a real version passes. It is also the seam-wiring nail: the fixture
// and the production bring-up both reach the create path only through
// panelRuntimeVersion, so dropping that default breaks this file's compile
// instead of silently probing a different world.
func TestAC5ProbeMapsEveryMissingAnswerToTheNamedRefusal33ac5(t *testing.T) {
	if panelRuntimeVersion == nil {
		t.Fatal("the probe seam must be wired (panelRuntimeVersion nil)")
	}
	orig := panelRuntimeVersion
	t.Cleanup(func() { panelRuntimeVersion = orig })

	panelRuntimeVersion = func() (string, error) { return "", nil }
	if err := probePanelRuntime(); !errors.Is(err, errPanelRuntimeMissing) {
		t.Errorf("an empty version (no install) must refuse by name, got: %v", err)
	}
	panelRuntimeVersion = func() (string, error) { return "", errors.New("loader gone") }
	err := probePanelRuntime()
	if !errors.Is(err, errPanelRuntimeMissing) {
		t.Errorf("a loader error must refuse by name (fail-closed), got: %v", err)
	}
	if !strings.Contains(err.Error(), "loader gone") {
		t.Errorf("the refusal must carry the loader's own failure verbatim, got: %v", err)
	}
	panelRuntimeVersion = func() (string, error) { return "138.0.3.7", nil }
	if err := probePanelRuntime(); err != nil {
		t.Errorf("a present runtime must pass the probe, got: %v", err)
	}
}
