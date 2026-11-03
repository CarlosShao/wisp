package manifestslot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/plugin"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// Ticket 50 lifecycle tests: Manager's update flow (re-authorization, C16
// refusal, repin + live declaration swap), the restart path over
// plugin_state (AC4 gate per record), uninstall teardown under C11
// DisposalScope (zero files / zero tracked goroutines, ask-about-data
// default keep), and the AC5 net-allowlist verdict reached through an
// actually installed, registered plugin tool.

const (
	lifeID      = "com.example.life"
	lifeInHost  = "open.feishu.cn"
	lifeOutHost = "exfil.evil.net"
)

// lifecycleManifest declares one L0 net tool with a domain allowlist, the
// exact shape SPEC-07 §4 shows.
const lifecycleManifest = `
schema_version = 1
id = "com.example.life"
host_api = "^1.0.0"
name = "Lifecycle Plugin"
version = "0.1.0"

[[tools]]
name = "life.agenda"
description = "Query the agenda."
risk = "L0"
capabilities = ["net"]
net_allowlist = ["open.feishu.cn"]

[tools.parameters]
type = "object"
[tools.parameters.properties]
host.type = "string"
`

// lifecycleManifestV2 is version 0.2.0 with a different allowlist face: the
// old in-list host must start escalating after the update lands.
var lifecycleManifestV2 = strings.NewReplacer(
	`version = "0.1.0"`, `version = "0.2.0"`,
	`net_allowlist = ["open.feishu.cn"]`, `net_allowlist = ["api.example.com"]`,
).Replace(lifecycleManifest)

func echoRun(_ context.Context, _ string, def plugin.ToolDef, params json.RawMessage) (tools.Result, error) {
	return tools.Result{Text: "ok:" + def.Name + ":" + string(params)}, nil
}

// netGateRun is the mock plugin behaviour for AC5 e2e: it judges the call
// target through the real central assessor with the tool's CURRENT declared
// allowlist and reports the verdict in the result text.
func netGateRun(_ context.Context, _ string, def plugin.ToolDef, params json.RawMessage) (tools.Result, error) {
	var p struct {
		Host string `json:"host"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return tools.Result{IsError: true, Text: err.Error()}, nil
	}
	dec, err := plugin.AssessPluginNetCall(def, risk.NetTarget{
		Protocol: "https",
		Host:     p.Host,
		URL:      "https://" + p.Host + "/x",
	}, nil, nil)
	if err != nil {
		return tools.Result{IsError: true, Text: err.Error()}, nil
	}
	return tools.Result{Text: fmt.Sprintf("level=%v rules=%v", dec.Level, dec.RulesHit)}, nil
}

func newLifeManager(root string, st *memStore, auth Authorizer, hostVersion string, run Run) *Manager {
	return &Manager{
		PluginsRoot: root,
		Registry:    tools.NewRegistry(),
		Store:       st,
		Authorizer:  auth,
		HostVersion: hostVersion,
		Exec:        run,
	}
}

func callTool(t *testing.T, reg *tools.Registry, name string, params string) (string, error) {
	t.Helper()
	e, ok := reg.Lookup(name)
	if !ok {
		return "", fmt.Errorf("tool %q not in registry", name)
	}
	res, err := e.Tool.Execute(t.Context(), json.RawMessage(params), nil)
	if err != nil {
		return "", err
	}
	return res.Text, nil
}

// TestInstallManagerUpdateReauthorizesAndRepins is the update half of the
// lifecycle contract: update = re-run the install authorization flow (full
// capability sheet + source + what is being replaced), re-pin
// plugin_state.hash, and swap the live declarations. Identical bytes are a
// no-change success, not a second authorization ceremony.
func TestInstallManagerUpdateReauthorizesAndRepins(t *testing.T) {
	ctx := t.Context()
	st := &memStore{}
	auth := &allowAuth{}
	mgr := newLifeManager(t.TempDir(), st, auth, "1.4.2", echoRun)

	if _, err := mgr.Install(ctx, []byte(lifecycleManifest), "v1.toml"); err != nil {
		t.Fatalf("install: %v", err)
	}
	// No-change: the exact same bytes must not re-authorize or re-pin.
	res, err := mgr.Update(ctx, lifeID, []byte(lifecycleManifest), "v1-same.toml")
	if err != nil {
		t.Fatalf("no-change update: %v", err)
	}
	if !res.NoChange {
		t.Fatalf("no-change update flagged as change: %+v", res)
	}
	if len(auth.saw) != 1 {
		t.Fatalf("authorization calls = %d, want 1 (no-change updates ask nobody)", len(auth.saw))
	}

	if _, err := mgr.Update(ctx, lifeID, []byte(lifecycleManifestV2), "v2.toml"); err != nil {
		t.Fatalf("update: %v", err)
	}
	if len(auth.saw) != 2 {
		t.Fatalf("authorization calls = %d, want 2 (install + update ceremony)", len(auth.saw))
	}
	sheet := auth.saw[1]
	if sheet.PreviousVersion != "0.1.0" {
		t.Errorf("update sheet PreviousVersion = %q, want 0.1.0", sheet.PreviousVersion)
	}
	if sheet.Version != "0.2.0" || sheet.Source != "v2.toml" {
		t.Errorf("update sheet identity/source mismatch: %+v", sheet)
	}
	if len(sheet.Capabilities) != 1 || sheet.Capabilities[0] != "net" {
		t.Errorf("update sheet capabilities = %v, want the full face [net]", sheet.Capabilities)
	}

	wantHash := plugin.ManifestHash([]byte(lifecycleManifestV2))
	if len(st.saved) != 1 || st.saved[0].Hash != wantHash {
		t.Fatalf("state must hold the new pinned hash, got %+v", st.saved)
	}
	// The AC4 gate now vouches for the NEW bytes only: the old hash refuses
	// the installed file (proof the pin moved, not merely appended).
	oldHash := plugin.ManifestHash([]byte(lifecycleManifest))
	if _, err := plugin.LoadInstalled(res.Path, plugin.LoadOptions{PinnedHash: oldHash, PluginID: lifeID}); !errors.Is(err, plugin.ErrTampered) {
		t.Fatalf("load under the old pin err = %v, want ErrTampered", err)
	}
	loaded, err := plugin.LoadInstalled(res.Path, plugin.LoadOptions{PinnedHash: wantHash, PluginID: lifeID})
	if err != nil || loaded.Version != "0.2.0" {
		t.Fatalf("load under the new pin: %v (%+v)", err, loaded)
	}
	// Live declarations swapped: the next executed call receives the NEW
	// ToolDef (allowlist api.example.com), not the snapshot the registry
	// entry was built from.
	seenDef := ""
	mgr.Exec = func(_ context.Context, _ string, def plugin.ToolDef, _ json.RawMessage) (tools.Result, error) {
		seenDef = strings.Join(def.NetAllowlist, ",")
		return tools.Result{Text: "ok"}, nil
	}
	if _, err := callTool(t, mgr.Registry, "life.agenda", `{"host":"x"}`); err != nil {
		t.Fatalf("call after update: %v", err)
	}
	if seenDef != "api.example.com" {
		t.Fatalf("live Run handed def allowlist %q, want the updated face", seenDef)
	}
}

// TestInstallManagerUpdateMajorIncompatibleRefused pins the C16 update rule:
// a new manifest the running host cannot satisfy is refused with the
// explicit C16 error and the OLD install keeps serving — untouched bytes,
// untouched pin, untouched authorization count.
func TestInstallManagerUpdateMajorIncompatibleRefused(t *testing.T) {
	ctx := t.Context()
	st := &memStore{}
	auth := &allowAuth{}
	mgr := newLifeManager(t.TempDir(), st, auth, "1.4.2", echoRun)
	res, err := mgr.Install(ctx, []byte(lifecycleManifest), "v1.toml")
	if err != nil {
		t.Fatalf("install: %v", err)
	}

	bad := strings.Replace(lifecycleManifest, `host_api = "^1.0.0"`, `host_api = "^2.0"`, 1)
	_, err = mgr.Update(ctx, lifeID, []byte(bad), "badhost.toml")
	if err == nil {
		t.Fatal("update to a host-incompatible manifest was accepted")
	}
	if !strings.Contains(err.Error(), "host_api") || !strings.Contains(err.Error(), "major") {
		t.Fatalf("err = %v, want the explicit C16 major-mismatch refusal", err)
	}
	if len(auth.saw) != 1 {
		t.Fatalf("C16-refused update must not reach the authorizer, calls = %d", len(auth.saw))
	}
	disk, err := os.ReadFile(res.Path)
	if err != nil || string(disk) != lifecycleManifest {
		t.Fatalf("old installed bytes must be untouched: err=%v", err)
	}
	if st.saved[0].Hash != res.Hash {
		t.Fatalf("old pin must be untouched: %+v", st.saved)
	}
	if text, err := callTool(t, mgr.Registry, "life.agenda", `{"host":"y"}`); err != nil || text == "" {
		t.Fatalf("the refused update must leave the old tool serving: %v %q", err, text)
	}
}

// TestInstallManagerUpdateDenialKeepsOldInstall: the user says no to the
// update sheet — nothing on disk, in state, or in the live view moves.
func TestInstallManagerUpdateDenialKeepsOldInstall(t *testing.T) {
	ctx := t.Context()
	st := &memStore{}
	auth := &scriptedAuth{answers: []error{nil, errors.New("user declined the update")}}
	mgr := newLifeManager(t.TempDir(), st, auth, "1.4.2", echoRun)
	res, err := mgr.Install(ctx, []byte(lifecycleManifest), "v1.toml")
	if err != nil {
		t.Fatalf("install: %v", err)
	}

	_, err = mgr.Update(ctx, lifeID, []byte(lifecycleManifestV2), "v2.toml")
	if !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("declined update err = %v, want ErrNotAuthorized", err)
	}
	disk, err := os.ReadFile(res.Path)
	if err != nil || string(disk) != lifecycleManifest {
		t.Fatalf("declined update must keep the old bytes: err=%v", err)
	}
	if st.saved[0].Hash != res.Hash {
		t.Fatalf("declined update must keep the old pin")
	}

	// Cross-install guard: an update carrying a DIFFERENT manifest id is
	// refused outright, and updates of never-installed ids are refused.
	other := strings.Replace(lifecycleManifestV2, lifeID, "com.example.other", 1)
	if _, err := mgr.Update(ctx, lifeID, []byte(other), "x.toml"); !errors.Is(err, ErrUpdateIDMismatch) {
		t.Fatalf("cross-id update err = %v, want ErrUpdateIDMismatch", err)
	}
	if _, err := mgr.Update(ctx, "com.example.never", []byte(lifecycleManifestV2), "x.toml"); !errors.Is(err, ErrPluginNotInstalled) {
		t.Fatalf("update of an absent plugin err = %v, want ErrPluginNotInstalled", err)
	}
	if _, err := mgr.Update(ctx, "../escape", []byte(lifecycleManifestV2), "x.toml"); !errors.Is(err, ErrUnsafePluginID) {
		t.Fatalf("unsafe id err = %v, want ErrUnsafePluginID", err)
	}
	if _, err := mgr.Uninstall(ctx, "../escape", UninstallOptions{}); !errors.Is(err, ErrUnsafePluginID) {
		t.Fatalf("uninstall unsafe id err = %v, want ErrUnsafePluginID", err)
	}
}

type scriptedAuth struct {
	answers []error
	calls   int
	saw     []InstallInfo
}

func (a *scriptedAuth) Authorize(_ context.Context, info InstallInfo) error {
	a.saw = append(a.saw, info)
	i := a.calls
	a.calls++
	if i < len(a.answers) {
		return a.answers[i]
	}
	return nil
}

// TestTamperManagerLoadAllReverifiesPinnedHash is the restart path: a fresh
// Manager over the same root + plugin_state re-reads every record through
// the AC4 gate. A tampered plugin is refused by name with an alert, its
// sibling still loads, and the refused tool is NOT callable in the new
// registry.
func TestTamperManagerLoadAllReverifiesPinnedHash(t *testing.T) {
	ctx := t.Context()
	root := t.TempDir()
	st := &memStore{}
	auth := &allowAuth{}
	mgrA := newLifeManager(root, st, auth, "1.4.2", echoRun)
	if _, err := mgrA.Install(ctx, []byte(lifecycleManifest), "v1.toml"); err != nil {
		t.Fatalf("install life: %v", err)
	}
	if _, err := mgrA.Install(ctx, []byte(ac3Manifest), "ac3.toml"); err != nil {
		t.Fatalf("install ac3: %v", err)
	}

	// In-place swap of the installed life manifest (bit-flip grade: any byte
	// change), exactly what plugin_state.hash exists to catch.
	manifestPath := filepath.Join(root, lifeID, "manifest.toml")
	if err := os.WriteFile(manifestPath, []byte(lifecycleManifest+"# swapped\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var alert plugin.CaptureAlerter
	mgrB := newLifeManager(root, st, auth, "1.4.2", echoRun)
	mgrB.Alerter = &alert
	verdicts, err := mgrB.LoadAll(ctx)
	if err != nil {
		t.Fatalf("LoadAll: %v", err)
	}
	byID := map[string]LoadVerdict{}
	for _, v := range verdicts {
		byID[v.PluginID] = v
	}
	if len(verdicts) != 2 {
		t.Fatalf("verdicts = %+v, want one per plugin_state record", verdicts)
	}
	if v := byID[lifeID]; v.Loaded {
		t.Fatalf("tampered plugin loaded: %+v", v)
	} else if !strings.Contains(v.Reason, "tamper") {
		t.Errorf("refusal reason = %q, want it to name the tamper/C14 gate", v.Reason)
	}
	if v := byID["com.example.ac3"]; !v.Loaded {
		t.Fatalf("the healthy sibling must still load: %+v", v)
	}
	if n := len(alert.Alerts()); n != 1 {
		t.Fatalf("tamper alerts = %d, want 1 (visible, never silent)", n)
	}
	if _, ok := mgrB.Registry.Lookup("life.agenda"); ok {
		t.Fatal("refused plugin's tool is registered in the new process")
	}
	if _, err := callTool(t, mgrB.Registry, "ac3.echo", `{}`); err != nil {
		t.Errorf("loaded sibling call: %v", err)
	}
}

// TestHostApiManagerLoadAllRefusesMajorMismatch: a plugin installed under
// host 1.x meets a 2.x host at load time — C16 refuses by name with the
// explicit error, it never silently skips.
func TestHostApiManagerLoadAllRefusesMajorMismatch(t *testing.T) {
	ctx := t.Context()
	root := t.TempDir()
	st := &memStore{}
	auth := &allowAuth{}
	mgrA := newLifeManager(root, st, auth, "1.4.2", echoRun)
	if _, err := mgrA.Install(ctx, []byte(lifecycleManifest), "v1.toml"); err != nil {
		t.Fatalf("install: %v", err)
	}

	mgrB := newLifeManager(root, st, auth, "2.0.0", echoRun)
	verdicts, err := mgrB.LoadAll(ctx)
	if err != nil {
		t.Fatalf("LoadAll: %v", err)
	}
	if len(verdicts) != 1 || verdicts[0].Loaded {
		t.Fatalf("verdicts = %+v, want the plugin refused", verdicts)
	}
	if !strings.Contains(verdicts[0].Reason, "host_api") || !strings.Contains(verdicts[0].Reason, "major") {
		t.Fatalf("reason = %q, want the explicit C16 major-mismatch error", verdicts[0].Reason)
	}
	if _, ok := mgrB.Registry.Lookup("life.agenda"); ok {
		t.Fatal("host-incompatible plugin registered anyway")
	}
}

// TestUninstallManagerTeardownLeavesZeroFilesAndGoroutines is AC3's second
// sentence: uninstall tears down running instances through the C11
// DisposalScope, deterministically. A call that is mid-flight when
// Uninstall starts must end cancelled, its tracked worker must be JOINED
// before Uninstall returns, the disposal result must be clean, and after it
// there must be zero installed files and zero state records.
func TestUninstallManagerTeardownLeavesZeroFilesAndGoroutines(t *testing.T) {
	ctx := t.Context()
	st := &memStore{}
	auth := &allowAuth{}
	var mgr *Manager

	entered := make(chan struct{})
	workerDone := make(chan struct{})
	run := func(callCtx context.Context, pluginID string, _ plugin.ToolDef, _ json.RawMessage) (tools.Result, error) {
		scope, ok := mgr.Scope(pluginID)
		if !ok {
			return tools.Result{IsError: true}, errors.New("scope vanished mid-call")
		}
		// The plugin's long-lived instance, registered with its scope.
		scope.Go("plugin-worker", func(scopeCtx context.Context) {
			select {
			case <-callCtx.Done():
			case <-scopeCtx.Done():
			}
			close(workerDone)
		})
		close(entered)
		<-callCtx.Done()
		return tools.Result{}, callCtx.Err()
	}

	mgr = newLifeManager(t.TempDir(), st, auth, "1.4.2", run)
	res, err := mgr.Install(ctx, []byte(lifecycleManifest), "v1.toml")
	if err != nil {
		t.Fatalf("install: %v", err)
	}

	type outcome struct {
		err error
	}
	calls := make(chan outcome, 1)
	// The mid-flight call must be running while Uninstall starts: an owned
	// worker with recover (same discipline as production code).
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				calls <- outcome{err: fmt.Errorf("mid-flight call panicked: %v", rec)}
			}
		}()
		_, err := callTool(t, mgr.Registry, "life.agenda", `{"host":"x"}`)
		calls <- outcome{err}
	}()
	<-entered

	rep, err := mgr.Uninstall(ctx, lifeID, UninstallOptions{})
	if err != nil {
		t.Fatalf("Uninstall: %v (disposal=%+v)", err, rep.Disposal)
	}

	o := <-calls
	if o.err == nil {
		t.Fatal("mid-flight call survived uninstall; the scope must cancel it")
	}

	// Zero tracked goroutines: the worker had to join before Dispose's
	// step-2 wait returned, so its channel is already closed here.
	select {
	case <-workerDone:
	default:
		t.Fatal("uninstall returned while a plugin worker goroutine was still tracked")
	}
	if rep.Disposal.Incomplete || rep.Disposal.TimedOut {
		t.Fatalf("disposal not clean: %+v", rep.Disposal)
	}
	if len(rep.Disposal.Failed) != 0 {
		t.Fatalf("disposal failures: %+v", rep.Disposal.Failed)
	}
	// Zero files, zero state.
	if _, err := os.Stat(filepath.Dir(res.Path)); !os.IsNotExist(err) {
		t.Fatalf("installed dir must be gone, err=%v", err)
	}
	if len(st.saved) != 0 {
		t.Fatalf("plugin_state records must be gone, got %+v", st.saved)
	}
	// The stale registry name is refused at execution (boundary: names are
	// removed when the host rebuilds the Registry from plugin_state).
	if _, err := callTool(t, mgr.Registry, "life.agenda", `{"host":"x"}`); !errors.Is(err, ErrPluginNotLive) {
		t.Fatalf("post-uninstall call err = %v, want ErrPluginNotLive", err)
	}
}

// TestUninstallManagerAsksAboutDataDefaultKeeps pins the ask-about-data
// default: uninstalling removes the installed files but keeps the plugin's
// data tree, reporting where it lives; only the explicit answer deletes it.
func TestUninstallManagerAsksAboutDataDefaultKeeps(t *testing.T) {
	ctx := t.Context()
	root := t.TempDir()
	dataRoot := t.TempDir()
	st := &memStore{}
	auth := &allowAuth{}
	mgr := newLifeManager(root, st, auth, "1.4.2", echoRun)
	mgr.DataRoot = dataRoot

	if _, err := mgr.Install(ctx, []byte(lifecycleManifest), "v1.toml"); err != nil {
		t.Fatalf("install: %v", err)
	}
	dataDir := filepath.Join(dataRoot, lifeID)
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		t.Fatal(err)
	}
	dataFile := filepath.Join(dataDir, "state.json")
	if err := os.WriteFile(dataFile, []byte(`{"kept":true}`), 0o644); err != nil {
		t.Fatal(err)
	}

	rep, err := mgr.Uninstall(ctx, lifeID, UninstallOptions{})
	if err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	if rep.DataDeleted {
		t.Fatal("default answer must KEEP plugin data")
	}
	if rep.DataDir != dataDir {
		t.Errorf("report DataDir = %q, want %q (the UI needs the address to ask with)", rep.DataDir, dataDir)
	}
	if _, err := os.Stat(dataFile); err != nil {
		t.Fatalf("kept data vanished: %v", err)
	}

	// Second plugin, explicit delete answer.
	if _, err := mgr.Install(ctx, []byte(ac3Manifest), "ac3.toml"); err != nil {
		t.Fatalf("reinstall: %v", err)
	}
	rep, err = mgr.Uninstall(ctx, "com.example.ac3", UninstallOptions{DeletePluginData: true})
	if err != nil {
		t.Fatalf("Uninstall(delete): %v", err)
	}
	if !rep.DataDeleted || rep.DataDir == "" {
		t.Fatalf("delete answer must remove and report the data dir: %+v", rep)
	}
	if _, err := os.Stat(filepath.Join(dataRoot, "com.example.ac3")); !os.IsNotExist(err) {
		t.Fatal("data dir survived the explicit delete answer")
	}
	// Uninstalling what is not installed (and not recorded) is an error,
	// never a silent success.
	if _, err := mgr.Uninstall(ctx, "com.example.never", UninstallOptions{}); err == nil {
		t.Fatal("uninstall of an uninstalled plugin reported success")
	}
}

// TestAllowlistManagerInstalledToolEscalatesE2E is AC5 e2e with a mock
// plugin that actually runs: through install -> register -> registry
// Lookup -> Execute, an in-list host passes on the declared floor (L0) and
// an out-of-list host is auto-escalated to L2 naming R5. After an update
// swaps the allowlist, the previously-in-list host escalates too — the
// verdict always uses the CURRENT declaration, never the registration-time
// snapshot.
func TestAllowlistManagerInstalledToolEscalatesE2E(t *testing.T) {
	ctx := t.Context()
	st := &memStore{}
	auth := &allowAuth{}
	mgr := newLifeManager(t.TempDir(), st, auth, "1.4.2", netGateRun)
	if _, err := mgr.Install(ctx, []byte(lifecycleManifest), "v1.toml"); err != nil {
		t.Fatalf("install: %v", err)
	}

	text, err := callTool(t, mgr.Registry, "life.agenda", `{"host":"open.feishu.cn"}`)
	if err != nil {
		t.Fatalf("in-list call: %v", err)
	}
	if !strings.Contains(text, "level=L0") || strings.Contains(text, "R5") {
		t.Fatalf("in-list verdict = %q, want L0 without an R5 hit", text)
	}

	text, err = callTool(t, mgr.Registry, "life.agenda", `{"host":"exfil.evil.net"}`)
	if err != nil {
		t.Fatalf("out-of-list call: %v", err)
	}
	if !strings.Contains(text, "level=L2") || !strings.Contains(text, "R5") {
		t.Fatalf("out-of-list verdict = %q, want L2 with the R5 rule named", text)
	}

	// Update retires the old domain from the allowlist.
	if _, err := mgr.Update(ctx, lifeID, []byte(lifecycleManifestV2), "v2.toml"); err != nil {
		t.Fatalf("update: %v", err)
	}
	text, err = callTool(t, mgr.Registry, "life.agenda", `{"host":"open.feishu.cn"}`)
	if err != nil {
		t.Fatalf("call after update: %v", err)
	}
	if !strings.Contains(text, "level=L2") || !strings.Contains(text, "R5") {
		t.Fatalf("host retired by the update still passed on the old face: %q", text)
	}
	// And the new face works.
	text, err = callTool(t, mgr.Registry, "life.agenda", `{"host":"api.example.com"}`)
	if err != nil {
		t.Fatalf("new-face call: %v", err)
	}
	if !strings.Contains(text, "level=L0") {
		t.Fatalf("new in-list host = %q, want L0", text)
	}
}
