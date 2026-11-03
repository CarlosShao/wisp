package manifestslot

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/plugin"
	"github.com/CarlosShao/wisp/internal/tools"
)

func okRun(context.Context, string, plugin.ToolDef, json.RawMessage) (tools.Result, error) {
	return tools.Result{Text: "ok"}, nil
}

type memStore struct{ saved []InstalledRecord }

func (m *memStore) SaveInstalled(_ context.Context, r InstalledRecord) error {
	// Upsert by id, mirroring memory.UpsertPluginState's ON CONFLICT so the
	// update flow re-pins in place instead of stacking records.
	for i, s := range m.saved {
		if s.ID == r.ID {
			m.saved[i] = r
			return nil
		}
	}
	m.saved = append(m.saved, r)
	return nil
}

func (m *memStore) DeleteInstalled(_ context.Context, id string) error {
	for i, r := range m.saved {
		if r.ID == id {
			m.saved = append(m.saved[:i], m.saved[i+1:]...)
			return nil
		}
	}
	return errors.New("not recorded: " + id)
}

func (m *memStore) ListInstalled(_ context.Context) ([]InstalledRecord, error) {
	out := make([]InstalledRecord, 0, len(m.saved))
	for _, r := range m.saved {
		out = append(out, r)
	}
	return out, nil
}

type allowAuth struct{ saw []InstallInfo }

func (a *allowAuth) Authorize(_ context.Context, info InstallInfo) error {
	a.saw = append(a.saw, info)
	return nil
}

type denyAuth struct{ saw *[]InstallInfo }

func (d denyAuth) Authorize(_ context.Context, info InstallInfo) error {
	*d.saw = append(*d.saw, info)
	return errors.New("user declined")
}

func testInstaller(t *testing.T) (Installer, *memStore, *allowAuth, *tools.Registry) {
	t.Helper()
	root := t.TempDir()
	reg := tools.NewRegistry()
	st := &memStore{}
	auth := &allowAuth{}
	return Installer{
		PluginsRoot: root,
		Registry:    reg,
		Store:       st,
		Authorizer:  auth,
		HostVersion: "1.4.2",
		Runner:      okRun,
	}, st, auth, reg
}

const ac3Manifest = `
schema_version = 1
id = "com.example.ac3"
host_api = "^1.0.0"
name = "AC3 Plugin"
version = "0.1.0"

[[tools]]
name = "ac3.echo"
description = "Echo the input."
risk = "L0"
capabilities = ["fs.read"]

[tools.parameters]
type = "object"
[tools.parameters.properties]
text.type = "string"
`

func TestAC3_InstallE2E_RegistersVisibleTool(t *testing.T) {
	ctx := t.Context()
	in, st, auth, reg := testInstaller(t)
	res, err := in.Install(ctx, []byte(ac3Manifest), "/tmp/ac3.toml")
	if err != nil {
		t.Fatalf("Install: %v", err)
	}
	if res.Manifest.ID != "com.example.ac3" {
		t.Fatalf("id = %q, want com.example.ac3", res.Manifest.ID)
	}
	if res.Hash == "" || !strings.HasPrefix(res.Hash, "sha256:") {
		t.Fatalf("hash shape = %q, want sha256:<hex>", res.Hash)
	}

	wantDir := filepath.Join(in.PluginsRoot, "com.example.ac3")
	wantPath := filepath.Join(wantDir, "manifest.toml")
	if res.Path != wantPath {
		t.Fatalf("Path = %q, want %q (SPEC-07 §6 layout)", res.Path, wantPath)
	}
	disk, err := os.ReadFile(wantPath)
	if err != nil || string(disk) != ac3Manifest {
		t.Fatalf("disk mismatch: err=%v bytes=%d", err, len(disk))
	}

	if len(st.saved) != 1 || st.saved[0].Hash != res.Hash {
		t.Fatalf("Store.Saved: %+v want one pinned record", st.saved)
	}
	if st.saved[0].Version != "0.1.0" {
		t.Errorf("state version = %q, want 0.1.0", st.saved[0].Version)
	}

	e, ok := reg.Lookup("ac3.echo")
	if !ok {
		t.Fatalf("Lookup after install: missing (list=%v)", toolNames(reg))
	}
	if e.Decl.Provider != tools.KindManifest {
		t.Errorf("Decl.Provider = %q, want manifest", e.Decl.Provider)
	}
	if e.Decl.PluginID != "com.example.ac3" {
		t.Errorf("Decl.PluginID = %q, want com.example.ac3 (AC6 identity)", e.Decl.PluginID)
	}
	if len(e.Decl.Capabilities) != 1 || e.Decl.Capabilities[0] != tools.CapFSRead {
		t.Errorf("Decl.Capabilities = %v, want [fs.read]", e.Decl.Capabilities)
	}
	if len(e.Decl.Needs) != 1 || e.Decl.Needs[0] != tools.CapFSRead {
		t.Errorf("Decl.Needs = %v, want the declared face and nothing wider", e.Decl.Needs)
	}
	if e.Decl.Declared.String() != "L0" {
		t.Errorf("Decl.Declared = %v, want L0 floor from the manifest", e.Decl.Declared)
	}

	res2, err := e.Tool.Execute(ctx, json.RawMessage(`{"text":"hi"}`), nil)
	if err != nil || res2.IsError {
		t.Fatalf("Execute through installed slot: err=%v res=%+v", err, res2)
	}

	if len(auth.saw) != 1 {
		t.Fatalf("Authorizer calls = %d, want exactly 1", len(auth.saw))
	}
	info := auth.saw[0]
	if info.Source != "/tmp/ac3.toml" {
		t.Errorf("review sheet Source = %q", info.Source)
	}
	if len(info.Capabilities) != 1 || info.Capabilities[0] != "fs.read" {
		t.Errorf("review sheet capabilities = %v, want [fs.read]", info.Capabilities)
	}
	if info.PluginID != "com.example.ac3" || info.Version != "0.1.0" {
		t.Errorf("review sheet identity mismatch: %+v", info)
	}

	loaded, err := plugin.LoadInstalled(wantPath, plugin.LoadOptions{
		PinnedHash: st.saved[0].Hash,
		PluginID:   "com.example.ac3",
	})
	if err != nil {
		t.Fatalf("post-install LoadInstalled (AC4 gate): %v", err)
	}
	if loaded.ID != "com.example.ac3" || len(loaded.Tools) != 1 {
		t.Fatalf("loaded manifest shape = %+v", loaded)
	}
}

func TestAC3_ReinstallSameID(t *testing.T) {
	ctx := t.Context()
	in, _, _, _ := testInstaller(t)
	if _, err := in.Install(ctx, []byte(ac3Manifest), "ac3.toml"); err != nil {
		t.Fatalf("first install: %v", err)
	}
	_, err := in.Install(ctx, []byte(ac3Manifest), "ac3.toml")
	if !errors.Is(err, ErrAlreadyInstall) {
		t.Fatalf("second install err = %v, want ErrAlreadyInstall", err)
	}
	if !strings.Contains(err.Error(), "com.example.ac3") {
		t.Errorf("error should name the plugin id: %v", err)
	}
}

func TestAC3_AuthorizerDenies(t *testing.T) {
	ctx := t.Context()
	root := t.TempDir()
	reg := tools.NewRegistry()
	st := &memStore{}
	var saw []InstallInfo
	in := Installer{
		PluginsRoot: root, Registry: reg, Store: st,
		Authorizer: denyAuth{&saw}, HostVersion: "1.4.2", Runner: okRun,
	}

	_, err := in.Install(ctx, []byte(ac3Manifest), "denied.toml")
	if !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("denied install err = %v, want ErrNotAuthorized", err)
	}
	if len(saw) != 1 {
		t.Fatalf("review sheet calls = %d, want 1 before refusal", len(saw))
	}
	if len(st.saved) != 0 {
		t.Fatalf("no state record expected, got %+v", st.saved)
	}
	if names := toolNames(reg); len(names) != 0 {
		t.Fatalf("registry must be empty after denial, got %v", names)
	}
	if _, err := os.Stat(filepath.Join(root, "com.example.ac3")); !os.IsNotExist(err) {
		t.Fatalf("no plugin dir expected; err=%v", err)
	}
}

func TestAC3_SchemaInvalid(t *testing.T) {
	ctx := t.Context()
	in, st, auth, reg := testInstaller(t)
	bad := strings.Replace(ac3Manifest, `name = "ac3.echo"`, `name = "Echo!"`, 1)
	_, err := in.Install(ctx, []byte(bad), "bad.toml")
	if err == nil || !strings.Contains(err.Error(), "C14") {
		t.Fatalf("schema-bad err = %v, want C14 refusal", err)
	}
	if len(st.saved) != 0 || len(auth.saw) != 0 || len(toolNames(reg)) != 0 {
		t.Fatalf("schema refusal must leave zero trace")
	}
}

func TestAC3_HostAPIRefuses(t *testing.T) {
	ctx := t.Context()
	in, st, auth, reg := testInstaller(t)
	in.HostVersion = "2.0.0"
	_, err := in.Install(ctx, []byte(ac3Manifest), "host2.toml")
	if err == nil || !strings.Contains(err.Error(), "host_api") {
		t.Fatalf("host2 err = %v, want host_api C16 refusal", err)
	}
	if len(st.saved) != 0 || len(auth.saw) != 0 || len(toolNames(reg)) != 0 {
		t.Fatalf("C16 refusal must leave zero trace")
	}
}

// TestAC3_DuplicateToolNameRefused is the C1 collision reverse: a second
// plugin declaring a tool name that is already registered is refused whole —
// before the user is even asked to authorize it, and before anything lands.
// tools.Registry silently skips names it already has on provider registration,
// so the refusal must happen here, not there.
func TestAC3_DuplicateToolNameRefused(t *testing.T) {
	ctx := t.Context()
	in, st, auth, reg := testInstaller(t)
	if _, err := in.Install(ctx, []byte(ac3Manifest), "first.toml"); err != nil {
		t.Fatalf("first install: %v", err)
	}

	// Same tool name, different plugin id.
	dup := strings.Replace(ac3Manifest, "com.example.ac3", "com.example.other", 1)
	_, err := in.Install(ctx, []byte(dup), "dup.toml")
	if !errors.Is(err, ErrToolNameTaken) {
		t.Fatalf("dup-tool install err = %v, want ErrToolNameTaken", err)
	}
	if len(st.saved) != 1 || st.saved[0].ID != "com.example.ac3" {
		t.Fatalf("state must hold only the surviving install, got %+v", st.saved)
	}
	if _, statErr := os.Stat(filepath.Join(in.PluginsRoot, "com.example.other")); !os.IsNotExist(statErr) {
		t.Fatalf("refused dir must not exist; err=%v", statErr)
	}
	if names := toolNames(reg); len(names) != 1 {
		t.Fatalf("registry = %v, want only the first install's tools", names)
	}
	if len(auth.saw) != 1 {
		t.Fatalf("authorizer must not be consulted for a manifest that cannot install, saw %d", len(auth.saw))
	}
}

func TestAC3_RestartLoadUsesPinnedHash(t *testing.T) {
	ctx := t.Context()
	in, st, _, _ := testInstaller(t)
	res, err := in.Install(ctx, []byte(ac3Manifest), "boot.toml")
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	rec := st.saved[0]

	m, err := plugin.LoadInstalled(res.Path, plugin.LoadOptions{
		PinnedHash: rec.Hash, PluginID: rec.ID,
	})
	if err != nil {
		t.Fatalf("LoadInstalled: %v", err)
	}
	if m.ID != rec.ID {
		t.Fatalf("reload mismatch: %+v vs %+v", m, rec)
	}

	if err := os.WriteFile(res.Path, []byte(ac3Manifest+"# flip\n"), 0o644); err != nil {
		t.Fatalf("flip: %v", err)
	}
	var alert plugin.CaptureAlerter
	_, err = plugin.LoadInstalled(res.Path, plugin.LoadOptions{
		PinnedHash: rec.Hash, PluginID: rec.ID, Alerter: &alert,
	})
	if !errors.Is(err, plugin.ErrTampered) {
		t.Fatalf("tampered load err = %v, want ErrTampered", err)
	}
	if len(alert.Alerts()) != 1 {
		t.Fatalf("tamper alerts = %d, want 1", len(alert.Alerts()))
	}
}

func TestAC3_Uninstall(t *testing.T) {
	ctx := t.Context()
	in, st, _, reg := testInstaller(t)
	res, err := in.Install(ctx, []byte(ac3Manifest), "un.toml")
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if len(toolNames(reg)) != 1 {
		t.Fatalf("tool should be visible after install")
	}

	if err := in.Uninstall(ctx, "com.example.ac3"); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	if len(st.saved) != 0 {
		t.Fatalf("state record must be gone, got %+v", st.saved)
	}
	if _, err := os.Stat(filepath.Dir(res.Path)); !os.IsNotExist(err) {
		t.Fatalf("plugin dir must be gone; err=%v", err)
	}
}

func TestAC3_InstallerUsesMemoryDAO(t *testing.T) {
	ctx := t.Context()
	root := t.TempDir()
	dbDir := t.TempDir()
	st, err := memory.Open(dbDir)
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	dao := memoryDAO{store: st}
	reg := tools.NewRegistry()
	auth := &allowAuth{}
	in := Installer{
		PluginsRoot: root, Registry: reg, Store: dao,
		Authorizer: auth, HostVersion: "1.4.2", Runner: okRun,
	}

	res, err := in.Install(ctx, []byte(ac3Manifest), "dao.toml")
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	rec, err := st.PluginStateByID(ctx, "com.example.ac3")
	if err != nil {
		t.Fatalf("PluginStateByID: %v", err)
	}
	if rec.ID != "com.example.ac3" || rec.Version != "0.1.0" || rec.Hash != res.Hash {
		t.Fatalf("DAO record mismatch: %+v res.Hash=%q", rec, res.Hash)
	}
	if !rec.Enabled {
		t.Errorf("Enabled = false after install")
	}

	if err := in.Uninstall(ctx, "com.example.ac3"); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	if _, err := st.PluginStateByID(ctx, "com.example.ac3"); !errors.Is(err, memory.ErrNotFound) {
		t.Fatalf("DAO err = %v, want ErrNotFound", err)
	}
}

type memoryDAO struct{ store *memory.Store }

func (m memoryDAO) SaveInstalled(ctx context.Context, r InstalledRecord) error {
	caps := []string{"fs.read"}
	cj, err := json.Marshal(caps)
	if err != nil {
		return err
	}
	return m.store.UpsertPluginState(ctx, memory.PluginState{
		ID: r.ID, Version: r.Version, Enabled: true,
		CapabilitiesJSON: string(cj), Hash: r.Hash,
	})
}

func (m memoryDAO) DeleteInstalled(ctx context.Context, id string) error {
	return m.store.DeletePluginState(ctx, id)
}

func (m memoryDAO) ListInstalled(ctx context.Context) ([]InstalledRecord, error) {
	all, err := m.store.ListPluginStates(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]InstalledRecord, 0, len(all))
	for _, p := range all {
		out = append(out, InstalledRecord{ID: p.ID, Version: p.Version, Hash: p.Hash, ExeHash: p.ExeHash})
	}
	return out, nil
}

func toolNames(r *tools.Registry) []string {
	var out []string
	for _, e := range r.List() {
		out = append(out, e.Tool.Name())
	}
	return out
}
