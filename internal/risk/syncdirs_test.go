package risk

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

// P12 detection is fixture-driven so it runs anywhere; the Windows registry
// probe itself is exercised by the machine-reality test at the bottom (and
// its conclusions are recorded in docs/PRECHECK.md §P12).

func sandbox(t *testing.T) (home, lapp, app string) {
	t.Helper()
	base := t.TempDir()
	home = filepath.Join(base, "profile")
	lapp = filepath.Join(base, "local")
	app = filepath.Join(base, "roaming")
	for _, d := range []string{home, lapp, app} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return
}

func TestSyncDefaultLocationProbe(t *testing.T) {
	home, _, _ := sandbox(t)
	oned := filepath.Join(home, "OneDrive")
	if err := os.MkdirAll(oned, 0o755); err != nil {
		t.Fatal(err)
	}
	roots := defaultLocationProbe(probeEnv{Home: home})
	if len(roots) != 1 || roots[0].Provider != "OneDrive" || roots[0].Source != "default" {
		t.Fatalf("expected the OneDrive default dir to be detected, got %+v", roots)
	}
}

func TestSyncDropboxHostDBConfig(t *testing.T) {
	home, lapp, _ := sandbox(t)
	db := filepath.Join(lapp, "Dropbox", "host.db")
	if err := os.MkdirAll(filepath.Dir(db), 0o755); err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(home, "My Dropbox Moved") // relocated client: string matching alone would miss it
	if err := os.WriteFile(db, []byte("C:\\Users\\x\\Dropbox\\\\.dropbox\n"+
		base64.StdEncoding.EncodeToString([]byte(root))+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	roots := clientConfigProbe(probeEnv{Home: home, LocalAppData: lapp})
	for _, r := range roots {
		if r.Provider == "Dropbox" && r.Path == root {
			return
		}
	}
	t.Fatalf("relocated Dropbox root from host.db must be detected, got %+v", roots)
}

func TestSyncFixtureFallbackAndMatch(t *testing.T) {
	home, _, _ := sandbox(t)
	fx := filepath.Join(home, "syncfixture.json")
	syncRoot := filepath.Join(home, "Cloud", "NutstoreSync")
	body := `[{"provider":"Nutstore","path":"` + jsonEsc(syncRoot) + `"}]`
	if err := os.WriteFile(fx, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	p := NewProvenance(ProvOptions{NoProbe: true, HomeDir: home, SyncFixturePath: fx})
	if !p.SyncDetectionComplete() {
		t.Fatal("fixture roots count as detected")
	}
	st := p.IsSyncPath(filepath.Join(syncRoot, "out", "secret.md"))
	if !st.Sync || st.Root.Provider != "Nutstore" {
		t.Fatalf("write under fixture root must be sync, got %+v", st)
	}
	if p.IsSyncPath(filepath.Join(home, "LocalOnly", "a.txt")).Sync {
		t.Fatal("non-sync path must not flag")
	}
}

func TestSyncSuspectFallbackWhenUndetectable(t *testing.T) {
	// Zero detection success -> [fs] allowed_dirs under the user profile are
	// sync-suspect (safe default demanded by the ticket).
	home, _, _ := sandbox(t) // no sync dirs inside
	p := NewProvenance(ProvOptions{
		NoProbe:     true,
		HomeDir:     home,
		AllowedDirs: []string{filepath.Join(home, "Documents")},
	})
	if p.SyncDetectionComplete() {
		t.Fatal("nothing was detectable")
	}
	st := p.IsSyncPath(filepath.Join(home, "Documents", "report.txt"))
	if !st.Sync || st.Root.Provider != "sync-suspect" {
		t.Fatalf("under-profile allowed dir must be sync-suspect, got %+v", st)
	}
	// Outside the profile is not suspect.
	if p.IsSyncPath(`D:\plain\data.txt`).Sync {
		t.Fatal("path outside the user profile must not be suspect")
	}
}

func TestSyncUnresolvablePathFailClosed(t *testing.T) {
	home, _, _ := sandbox(t)
	p := NewProvenance(ProvOptions{NoProbe: true, HomeDir: home,
		SyncRoots: []SyncRoot{{Provider: "OneDrive", Path: filepath.Join(home, "OneDrive")}}})
	// Empty/absent path spelling: unverifiable -> sync (strict side).
	if !p.IsSyncPath("").Sync {
		t.Fatal("empty path must fail-closed as sync")
	}
}

func TestSyncRootsIntrospection(t *testing.T) {
	home, _, _ := sandbox(t)
	p := NewProvenance(ProvOptions{NoProbe: true, HomeDir: home, SyncRoots: []SyncRoot{
		{Provider: "GoogleDrive", Path: filepath.Join(home, "Google Drive")},
	}})
	roots := p.SyncRoots()
	if len(roots) != 1 || roots[0].Provider != "GoogleDrive" {
		t.Fatalf("roots view: %+v", roots)
	}
}

// TestSyncDetectionOnThisMachine is the self-hosted evidence run (AC: "Sync-dir
// detection on a machine with OneDrive configured"): it exercises the REAL
// registry + config probes on the dev machine and records what was found. It
// never asserts machine-specific state (CI-safe); it asserts the fail-closed
// invariant: either roots were detected, or suspects are active.
func TestSyncDetectionOnThisMachine(t *testing.T) {
	p := NewProvenance(ProvOptions{}) // live env
	if p.SyncDetectionComplete() {
		for _, r := range p.SyncRoots() {
			t.Logf("P12 evidence: detected %s root via %s: %s", r.Provider, r.Source, r.Path)
		}
	} else {
		t.Log("P12 evidence: no sync location detected on this machine -> under-profile allowed dirs fall back to sync-suspect")
	}
	home := userHomeDir()
	if home == "" {
		t.Skip("no profile home")
	}
	if !p.SyncDetectionComplete() {
		if !p.IsSyncPath(filepath.Join(home, "Documents", "x.txt")).Sync {
			t.Fatal("fail-closed invariant broken: undetectable -> suspect must flag")
		}
	}
	if od := filepath.Join(home, "OneDrive"); dirExists(od) {
		if !p.IsSyncPath(filepath.Join(od, "anything.txt")).Sync {
			t.Fatalf("OneDrive present at %s but writes there were not flagged", od)
		}
		t.Logf("P12 evidence: OneDrive at %s flagged as sync channel", od)
	}
}

func dirExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.IsDir()
}

func jsonEsc(s string) string {
	return filepath.ToSlash(s)
}
