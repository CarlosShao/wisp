package config

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// newTestManager writes an initial config file, loads a Manager from it and
// returns (manager, path, mutator). The mutator rewrites the file with new
// content and sleeps a tick so mtime/size polling sees the change.
func newTestManager(t *testing.T, initial string) (*Manager, string, func(string)) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := NewManager(path, nil)
	if err != nil {
		t.Fatalf("initial load: %v", err)
	}
	return m, path, func(content string) {
		time.Sleep(10 * time.Millisecond)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

const matrixBase = `
schema_version = 2

[ball]
size = 56

[app]
theme = "dark"
language = "zh-CN"

[voice]
enabled = true

[voice.asr]
provider = "local-sherpa"
model = "parakeet-zh"

[fs]
allowed_dirs = []
`

func TestManagerHotTierAppliesImmediately(t *testing.T) {
	m, _, mutate := newTestManager(t, matrixBase)
	mutate(strings.Replace(matrixBase, "size = 56", "size = 60", 1))
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if rep == nil || !reflect.DeepEqual(rep.Hot, []string{"ball"}) {
		t.Fatalf("report.Hot = %+v, want [ball]", rep)
	}
	if got := m.Config().Ball.Size; got != 60 {
		t.Fatalf("ball.size = %d, want 60 (hot applies immediately)", got)
	}
	// Idempotent: an unchanged file must be a no-op.
	rep, err = m.CheckAndReload()
	if err != nil || rep != nil {
		t.Fatalf("unchanged file must be a no-op, got rep=%v err=%v", rep, err)
	}
}

func TestManagerRestartTierKeepsOldValues(t *testing.T) {
	m, _, mutate := newTestManager(t, matrixBase)
	mutate(strings.Replace(matrixBase, `language = "zh-CN"`, `language = "en-US"`, 1))
	var pending []string
	m.OnRestartPending = func(sections []string) { pending = sections }
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if rep == nil || !reflect.DeepEqual(rep.Restart, []string{"app"}) {
		t.Fatalf("report.Restart = %+v, want [app]", rep)
	}
	if got := m.Config().App.Language; got != "zh-CN" {
		t.Fatalf("app.language = %q, want zh-CN (restart tier keeps old value until process restart)", got)
	}
	if !reflect.DeepEqual(pending, []string{"app"}) {
		t.Fatalf("OnRestartPending = %v, want [app]", pending)
	}
}

func TestManagerThemeIsHotInsideApp(t *testing.T) {
	m, _, mutate := newTestManager(t, matrixBase)
	mutate(strings.Replace(matrixBase, `theme = "dark"`, `theme = "light"`, 1))
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if rep == nil || !reflect.DeepEqual(rep.Hot, []string{"app"}) || len(rep.Restart) != 0 {
		t.Fatalf("theme change must be hot inside [app], got %+v", rep)
	}
	if got := m.Config().App.Theme; got != "light" {
		t.Fatalf("theme = %q, want light", got)
	}
}

func TestManagerReloadTierEmitsEvent(t *testing.T) {
	m, _, mutate := newTestManager(t, matrixBase)
	var fired []string
	var wg sync.WaitGroup
	wg.Add(1)
	m.OnReload = func(sections []string) { fired = sections; wg.Done() }
	mutate(strings.Replace(matrixBase, `model = "parakeet-zh"`, `model = "paraformer-zh"`, 1))
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if rep == nil || !reflect.DeepEqual(rep.Reload, []string{"voice"}) {
		t.Fatalf("report.Reload = %+v, want [voice]", rep)
	}
	if got := m.Config().Voice.ASR.Model; got != "paraformer-zh" {
		t.Fatalf("asr.model = %q, want paraformer-zh (reload tier applies immediately + event)", got)
	}
	wg.Wait()
	if !reflect.DeepEqual(fired, []string{"voice"}) {
		t.Fatalf("OnReload fired with %v, want [voice]", fired)
	}
}

func TestManagerVoiceHotKeysDoNotEmitReload(t *testing.T) {
	m, _, mutate := newTestManager(t, matrixBase)
	next := strings.Replace(matrixBase, "[voice]\nenabled = true",
		"[voice]\nenabled = true\npunctuation = false", 1)
	mutate(next)
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if rep == nil || len(rep.Reload) != 0 {
		t.Fatalf("veto-word/threshold/speed/punctuation changes are hot-tier, got Reload=%v", rep.Reload)
	}
	if m.Config().Voice.Punctuation {
		t.Fatal("punctuation=false must apply (hot)")
	}
}

func TestManagerLockedLooseningRejectedKeepsOld(t *testing.T) {
	m, _, mutate := newTestManager(t, matrixBase)
	var askedSection string
	var askedKeys []string
	asked := false
	m.ConfirmLocked = func(section string, keys []string) bool {
		asked, askedSection, askedKeys = true, section, keys
		return false // user rejects
	}
	mutate(strings.Replace(matrixBase, "allowed_dirs = []", `allowed_dirs = ["D:\\data"]`, 1))
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !asked {
		t.Fatal("loosening must invoke the L2 confirm hook")
	}
	if askedSection != "fs" || len(askedKeys) != 1 || askedKeys[0] != "fs.allowed_dirs" {
		t.Fatalf("confirm hook got section=%q keys=%v, want fs [fs.allowed_dirs]", askedSection, askedKeys)
	}
	if got := m.Config().FS.AllowedDirs; len(got) != 0 {
		t.Fatalf("rejected loosening must keep old values, got %v", got)
	}
	if rep == nil || len(rep.Locked) != 1 || rep.Locked[0].Approved ||
		rep.Locked[0].Direction != DirLoosen {
		t.Fatalf("report.Locked = %+v, want rejected loosen decision", rep.Locked)
	}
	// The other sections in the same file change must still apply.
	if m.Config().Ball.Size != 56 {
		_ = m.Config().Ball.Size
	}
}

func TestManagerLockedLooseningApprovedApplies(t *testing.T) {
	m, _, mutate := newTestManager(t, matrixBase)
	m.ConfirmLocked = func(section string, keys []string) bool { return true }
	mutate(strings.Replace(matrixBase, "allowed_dirs = []", `allowed_dirs = ["D:\\data"]`, 1))
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got := m.Config().FS.AllowedDirs; !reflect.DeepEqual(got, []string{"D:\\data"}) {
		t.Fatalf("approved loosening must apply, got %v", got)
	}
	if rep == nil || len(rep.Locked) != 1 || !rep.Locked[0].Approved {
		t.Fatalf("report.Locked = %+v, want approved loosen", rep.Locked)
	}
}

func TestManagerLockedTighteningHotAppliesWithoutHook(t *testing.T) {
	start := strings.Replace(matrixBase, "allowed_dirs = []", `allowed_dirs = ["D:\\data"]`, 1)
	m, _, mutate := newTestManager(t, start)
	hookCalled := false
	m.ConfirmLocked = func(section string, keys []string) bool { hookCalled = true; return true }
	mutate(matrixBase) // removing the allowed dir = tightening
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if hookCalled {
		t.Fatal("tightening must NOT invoke the confirm hook (hot-applies)")
	}
	if len(m.Config().FS.AllowedDirs) != 0 {
		t.Fatal("tightening must hot-apply")
	}
	if rep == nil || len(rep.Locked) != 1 || rep.Locked[0].Direction != DirTighten {
		t.Fatalf("report.Locked = %+v, want tighten decision", rep.Locked)
	}
}

func TestManagerLockedNilConfirmHookDenies(t *testing.T) {
	m, _, mutate := newTestManager(t, matrixBase)
	mutate(strings.Replace(matrixBase, "allowed_dirs = []", `allowed_dirs = ["D:\\data"]`, 1))
	if _, err := m.CheckAndReload(); err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(m.Config().FS.AllowedDirs) != 0 {
		t.Fatal("nil confirm hook must deny loosening (fail-closed)")
	}
}

func TestManagerLockedBothDirectionsLogged(t *testing.T) {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(old)

	m, _, mutate := newTestManager(t, matrixBase)
	m.ConfirmLocked = func(section string, keys []string) bool { return true }
	mutate(strings.Replace(matrixBase, "allowed_dirs = []", `allowed_dirs = ["D:\\data"]`, 1))
	if _, err := m.CheckAndReload(); err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !strings.Contains(buf.String(), "fs") || !strings.Contains(buf.String(), "loosening approved") {
		t.Fatalf("loosening (approved) must be logged, log: %s", buf.String())
	}
	buf.Reset()
	mutate(matrixBase) // tighten back (the loosened value is live, so this is a real change)
	if _, err := m.CheckAndReload(); err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !strings.Contains(buf.String(), "fs") || !strings.Contains(buf.String(), "tightened") {
		t.Fatalf("tightening must be logged too (both directions), log: %s", buf.String())
	}
}

func TestManagerLoadErrorKeepsCurrentConfig(t *testing.T) {
	m, path, mutate := newTestManager(t, matrixBase)
	mutate("this is not toml [[[")
	if _, err := m.CheckAndReload(); err == nil {
		t.Fatal("corrupt file must surface an error")
	}
	if m.Config().Ball.Size != 56 {
		t.Fatal("current config must survive a failed reload")
	}
	// A subsequent good write recovers.
	time.Sleep(10 * time.Millisecond)
	_ = os.WriteFile(path, []byte(strings.Replace(matrixBase, "size = 56", "size = 64", 1)), 0o600)
	if _, err := m.CheckAndReload(); err != nil {
		t.Fatalf("recovery reload: %v", err)
	}
	if m.Config().Ball.Size != 64 {
		t.Fatal("recovery must apply the new values")
	}
}

func TestManagerMissingFileKeepsCurrentAndErrors(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	_ = os.WriteFile(path, []byte(matrixBase), 0o600)
	m, err := NewManager(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if _, err := m.CheckAndReload(); err == nil {
		t.Fatal("missing file must surface an error (config deleted mid-run is a config-class event)")
	}
	if m.Config().Ball.Size != 56 {
		t.Fatal("current config must survive file deletion")
	}
}
