//go:build windows

package main

// Ticket 33 AC#10 - [panel]'s fields must have PRODUCERS, and the judgement has
// to be able to tell "read the value and really act on it" from "read it in a
// test".
//
// WHAT AC#10 ASKS FOR, verbatim from the ticket (33-panel-host-c27.md:245): the
// five [panel] fields must have a producer, width/height/scale must really set
// the window bounds, and hot must really take effect. The ticket's own warning
// is the ruler this file is built around: 判据要能区分"读值并真设边界"与"只在测试里
// 读一下", and ⛔ "把 default:\"640\" 改掉" is not a delivery.
//
// WHAT WAS ALREADY LANDED BEFORE THIS LEG, re-measured rather than assumed:
// width and height have had a production producer since 票 255 AC#4 -
// panelGeometrySource re-reads config.toml per create and windowOptions hands the
// pair to the ONE create site (cmd/wisp/panel_geometry_255_test.go is that
// ticket's ruler, and this file does not re-litigate it). So AC#10's remaining
// debt is the OTHER three fields, and this leg pays two of them:
//
//	enabled              - was read by nobody; now gated at the top of bringUp
//	keep_alive_in_session - was read by nobody; now decides whether Hide destroys
//	scale                - STILL has no producer, and the reason is a hard block,
//	                       not a preference: see TestAC10ScaleHasNoProducerAndTheReasonIsNamed
//
// EACH RULER HERE HAS A NAMED RED SHAPE (a test that can only ever be green is
// the failure this repository calls 恒真判据):
//
//	(a) TestAC10DisabledPanelRefusesBeforeTheRuntimeProbe - the red shape is
//	    deleting the enabled gate OR moving it below the probe. Both leave the
//	    spy on panelRuntimeVersion called, which this case demands stay at zero.
//	    That is the whole point of the ORDER: a policy refusal must not ask the
//	    machine a question.
//	(b) TestAC10DisabledRefusalIsAPolicyAnswerNotAMachineOne - the red shape is
//	    writing the refusal into runtimeMissing. That would make ticket 37's
//	    native fallback card fire on a user who merely switched the panel off.
//	(c) TestAC10EnabledSourceReadsTheRealConfigFile - the red shape is a
//	    snapshot-at-assembly source: the second read after a rewrite answers the
//	    first value, which contradicts [panel] being registered hot.
//	(d) TestAC10KeepAliveFalseTearsTheWindowDownOnHide - the red shape is Hide
//	    ignoring the source (leaving the window parked) or destroying when the
//	    source says keep.
//
// NO REAL WINDOW IS CREATED ANYWHERE HERE, and that is not a shortcut: (a), (b)
// and (c) stop before the create by construction (that IS the claim), and (d)
// drives Hide against a host whose hwnd is 0, which is Hide's documented silent
// return - so what (d) measures is the DECISION (did teardown run, and was
// Destroy reached) and never a Win32 effect. The on-screen half of
// keep_alive_in_session is a winlive reading and this file does not pretend to
// have one.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and dies
// at load (0xc0000135) unless third_party/sherpa-onnx is on PATH.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
)

// --- (a) the policy gate runs BEFORE the machine probe -------------------------

// TestAC10DisabledPanelRefusesBeforeTheRuntimeProbe is AC#10's enabled producer,
// measured as an ORDER.
//
// The spy on panelRuntimeVersion is the whole ruler. A host that refuses on
// policy without asking the machine leaves it at 0; a host that probes first (or
// ignores [panel] enabled and probes at all) leaves it at 1. Those are the two
// red shapes, and neither can be produced by a test that only reads the value.
func TestAC10DisabledPanelRefusesBeforeTheRuntimeProbe(t *testing.T) {
	var probes atomic.Int32
	restore := panelRuntimeVersion
	panelRuntimeVersion = func() (string, error) {
		probes.Add(1)
		return "120.0.0.0", nil
	}
	t.Cleanup(func() { panelRuntimeVersion = restore })

	m := NewPanelManager(panelDispatchForGeometry255(), nil, t.TempDir(),
		withPanelEnabled(func() bool { return false }))

	err := m.bringUp(context.Background())
	if err == nil {
		t.Fatal("a host told [panel] enabled = false brought a window up anyway: AC#10's enabled producer is not wired to the gate")
	}
	if !errors.Is(err, errPanelDisabledByConfig) {
		t.Errorf("the refusal is not the named one: %v (want errors.Is(errPanelDisabledByConfig), because a caller has to be able to tell a policy refusal from a machine one)", err)
	}
	if got := probes.Load(); got != 0 {
		t.Errorf("the runtime probe ran %d time(s) for a panel disabled in config.toml: the policy gate must sit ABOVE the machine probe, otherwise a user who switched the panel off still pays a Runtime query and ticket 37's fallback card cannot tell the two refusals apart", got)
	}
	if created := m.IsCreated(); created {
		t.Error("a disabled panel recorded a created window")
	}
}

// TestAC10DisabledRefusalIsAPolicyAnswerNotAMachineOne keeps the two refusal
// classes apart. The red shape is writing the policy refusal into the runtime
// flag: ticket 37's native fallback card reads Unavailable() to decide whether a
// native surface may stand in, and firing it for a user who merely wrote
// enabled = false would be the card claiming a machine defect that does not
// exist.
func TestAC10DisabledRefusalIsAPolicyAnswerNotAMachineOne(t *testing.T) {
	m := NewPanelManager(panelDispatchForGeometry255(), nil, t.TempDir(),
		withPanelEnabled(func() bool { return false }))

	if err := m.bringUp(context.Background()); err == nil {
		t.Fatal("bringUp succeeded for a disabled panel")
	}
	if missing, _ := m.Unavailable(); missing {
		t.Error("the policy refusal set the RUNTIME-missing flag: a panel switched off in config.toml is not a machine without WebView2, and ticket 37's fallback card reads that flag")
	}
	disabled, reason := m.Disabled()
	if !disabled {
		t.Error("Disabled() does not report the refusal that just happened, so the state has no reader")
	}
	if reason == "" {
		t.Error("the refusal carries no sentence for the operator")
	}
}

// TestAC10EnabledHostStillProbesTheRuntime is the reverse control for (a): the
// gate must be a GATE and not a wall. A host with the default (no source) and a
// host explicitly told true both reach the runtime probe exactly once, which is
// the pre-ticket behaviour every existing construction site keeps.
func TestAC10EnabledHostStillProbesTheRuntime(t *testing.T) {
	cases := []struct {
		name string
		opts []panelHostOption
	}{
		{name: "no source at all keeps the host enabled", opts: nil},
		{name: "an explicit true source is enabled", opts: []panelHostOption{withPanelEnabled(func() bool { return true })}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var probes atomic.Int32
			restore := panelRuntimeVersion
			panelRuntimeVersion = func() (string, error) {
				probes.Add(1)
				// Answer the missing shape so bringUp stops here: this case is
				// about the probe being REACHED, not about a window opening.
				return "", nil
			}
			t.Cleanup(func() { panelRuntimeVersion = restore })

			m := NewPanelManager(panelDispatchForGeometry255(), nil, t.TempDir(), tc.opts...)
			err := m.bringUp(context.Background())
			if !errors.Is(err, errPanelRuntimeMissing) {
				t.Fatalf("an ENABLED host did not stop at the runtime probe: %v", err)
			}
			if got := probes.Load(); got != 1 {
				t.Errorf("the runtime probe ran %d time(s) for an enabled host, want exactly 1: the enabled gate must not become a wall", got)
			}
		})
	}
}

// --- (c) the source is a live read, not a snapshot ----------------------------

// TestAC10EnabledSourceReadsTheRealConfigFile drives the PRODUCTION closure over a
// real config.toml, and rewrites the file in between. The red shape is a source
// that answers once: [panel] is registered hot, so a saved value has to be able to
// take effect without a restart, and a snapshot would answer the first value
// twice.
func TestAC10EnabledSourceReadsTheRealConfigFile(t *testing.T) {
	dir := t.TempDir()
	writePanelConfigAC10(t, dir, map[string]string{"enabled": "false"})

	src := panelEnabledSource(dir)
	if src() {
		t.Fatal("a config.toml carrying [panel] enabled = false was read as enabled")
	}

	writePanelConfigAC10(t, dir, map[string]string{"enabled": "true"})
	if !src() {
		t.Error("after [panel] enabled was rewritten to true the source still answers false: the value was snapshotted at assembly, which contradicts [panel] being registered hot")
	}

	// The unreadable branch: no file means the schema default, never a silent
	// disable. This is the direction panelEnabledSource documents.
	missing := panelEnabledSource(filepath.Join(dir, "no-such-dir"))
	if !missing() {
		t.Error("an unreadable config.toml disabled the panel: a parse error must not switch a feature off, it must fall back to the schema default (true)")
	}
}

// TestAC10KeepAliveSourceReadsTheRealConfigFile is (c) for the second field.
func TestAC10KeepAliveSourceReadsTheRealConfigFile(t *testing.T) {
	dir := t.TempDir()
	writePanelConfigAC10(t, dir, map[string]string{"keep_alive_in_session": "false"})

	src := panelKeepAliveSource(dir)
	if src() {
		t.Fatal("a config.toml carrying [panel] keep_alive_in_session = false was read as true")
	}

	writePanelConfigAC10(t, dir, map[string]string{"keep_alive_in_session": "true"})
	if !src() {
		t.Error("after keep_alive_in_session was rewritten to true the source still answers false: the value was snapshotted, not re-read")
	}

	missing := panelKeepAliveSource(filepath.Join(dir, "no-such-dir"))
	if !missing() {
		t.Error("an unreadable config.toml tore the window down: the fallback direction must be C27's keep-warm posture")
	}
}

// --- (d) keep_alive_in_session=false really destroys on hide -------------------

// TestAC10KeepAliveFalseTearsTheWindowDownOnHide is the second field's producer,
// measured as a DECISION plus a real teardown.
//
// Hide against a host with no window returns before any Win32 call (its
// documented silent return), so this case can run headlessly - but it would then
// prove nothing about teardown. So the case gives the host a window HANDLE it
// owns: a plain HWND is enough because teardown is Destroy's bookkeeping
// (created/hwnd/w) and Destroy's only Win32 effect is on the library's control,
// which is nil here. What is asserted is exactly the branch AC#10 owes: with the
// source answering false the host no longer reports a created window, and with it
// answering true the window is still there.
func TestAC10KeepAliveFalseTearsTheWindowDownOnHide(t *testing.T) {
	cases := []struct {
		name          string
		keepAlive     func() bool
		wantDestroyed bool
	}{
		{
			name:          "keep_alive_in_session = false destroys after hiding",
			keepAlive:     func() bool { return false },
			wantDestroyed: true,
		},
		{
			name:          "keep_alive_in_session = true keeps the window warm",
			keepAlive:     func() bool { return true },
			wantDestroyed: false,
		},
		{
			name:          "no source at all keeps C27's hide-don't-destroy posture",
			keepAlive:     nil,
			wantDestroyed: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var opts []panelHostOption
			if tc.keepAlive != nil {
				opts = append(opts, withPanelKeepAlive(tc.keepAlive))
			}
			m := NewPanelManager(panelDispatchForGeometry255(), nil, t.TempDir(), opts...)
			// A window that exists but owns no library control: enough for Hide's
			// bookkeeping path and for Destroy's own clearing, and no Win32 call
			// this file would then be pretending to have measured.
			m.mu.Lock()
			m.hwnd = 0x1234
			m.created = true
			m.shown = true
			m.mu.Unlock()

			m.Hide()

			if got := m.IsCreated(); got == tc.wantDestroyed {
				t.Errorf("after Hide, IsCreated()=%v with keep_alive_in_session answering %v (want created=%v): the field's only observable meaning is whether a hide leaves the WebView2 warm",
					got, !tc.wantDestroyed, !tc.wantDestroyed)
			}
			if tc.wantDestroyed {
				if m.KeepAliveOnHide() {
					t.Error("KeepAliveOnHide() answered true for a source that says false, so the accessor and the behaviour disagree")
				}
			}
		})
	}
}

// TestAC10AssemblyRootWiresBothSources is the WIRING ruler, and it exists because
// a source nobody hands to a host is the exact shape this repository calls
// "文案在、控件不在" moved to the wiring layer: the closure would compile, the
// per-field tests above would pass, and the production window would never ask it
// anything.
//
// It reads the assembly root's own text rather than a hand-copied struct, so the
// red shape is deleting either option from the NewPanelManager call at
// newResidentPanelManager.
func TestAC10AssemblyRootWiresBothSources(t *testing.T) {
	root := filepath.Join("panel_resident_windows.go")
	data, err := os.ReadFile(root)
	if err != nil {
		t.Fatalf("read the assembly root: %v", err)
	}
	src := string(data)
	for _, want := range []struct{ opt, src string }{
		{"withPanelEnabled", "panelEnabledSource"},
		{"withPanelKeepAlive", "panelKeepAliveSource"},
	} {
		if !strings.Contains(src, want.opt+"("+want.src+"(dataDir))") {
			t.Errorf("the assembly root does not hand %s(%s(dataDir)) to NewPanelManager: the source exists and no host ever asks it, so [panel]'s field has no producer in the production window", want.opt, want.src)
		}
	}
}

// --- the field that still has no producer, with the reason named --------------

// TestAC10ScaleHasNoProducerAndTheReasonIsNamed is the honest half of AC#10.
//
// [panel] scale ("UI zoom factor; 0 = follow system DPI", internal/config/schema.go)
// is the one field this leg does NOT land, and the block is structural rather than
// a scheduling choice:
//
//	the WebView2 zoom lives on ICoreWebView2Controller2.PutZoomFactor, and this
//	repository's fork declares that method ONLY as a raw ComProc vtbl slot
//	(forks/go-webview2/pkg/edge/ICoreWebView2Controller2.go:16). There is no Go
//	wrapper, the vtbl field is unexported, and Chromium.controller is unexported
//	too (forks/go-webview2/pkg/edge/chromium.go), so cmd/wisp cannot reach the
//	zoom factor without editing forks/** - which this leg may not touch.
//
// This case pins that reasoning so it cannot rot silently: if someone adds the
// wrapper, the first assertion fails and whoever did it is told to land the
// producer. It is deliberately a claim about the FORK and not about AC#10, so it
// stays true regardless of how the ticket is finally closed.
func TestAC10ScaleHasNoProducerAndTheReasonIsNamed(t *testing.T) {
	fork := filepath.Join("..", "..", "forks", "go-webview2", "pkg", "edge", "ICoreWebView2Controller2.go")
	data, err := os.ReadFile(fork)
	if err != nil {
		t.Skipf("the fork's controller2 file is not readable from here (%v), so the scale claim cannot be checked", err)
	}
	src := string(data)
	if strings.Contains(src, "func (i *ICoreWebView2Controller2) PutZoomFactor(") {
		t.Fatal("the fork now exports a Go PutZoomFactor wrapper: [panel] scale can be landed, and TestAC10ScaleHasNoProducerAndTheReasonIsNamed's excuse is no longer true - wire it to the create path and delete this case")
	}
	if !strings.Contains(src, "PutZoomFactor") {
		t.Fatal("the fork no longer even declares PutZoomFactor: the reasoning this case records has rotted, re-measure before repeating it")
	}
}

// writePanelConfigAC10 puts a loadable config.toml with the given [panel] rows in
// dir. schema_version is current so readConfigFile takes no migration branch
// (that branch rewrites the file, and a rewrite from under a re-read test would
// be a second mover in a measurement about who read what when).
func writePanelConfigAC10(t *testing.T, dir string, rows map[string]string) {
	t.Helper()
	body := "schema_version = 2\n\n[panel]\n"
	for _, k := range []string{"enabled", "width", "height", "keep_alive_in_session", "scale"} {
		v, ok := rows[k]
		if !ok {
			continue
		}
		body += k + " = " + v + "\n"
	}
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
