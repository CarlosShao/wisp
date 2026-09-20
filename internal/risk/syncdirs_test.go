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
	if !p.IsSyncPath(filepath.Join(syncRoot, "out", "secret.md")).Sync {
		t.Fatal("write under fixture root must be sync")
	}
	if got := p.SyncRoots(); len(got) != 1 || got[0].Provider != "Nutstore" {
		t.Fatalf("fixture roots must be listed, got %+v", got)
	}
	// M-3 (adversarial report): a fixture entry is user-supplied data, not
	// proof that a client syncs there, so it must NOT disarm the
	// under-profile suspect net. Asserting the opposite here is what let a
	// single decoy root switch the safety net off.
	if p.SyncDetectionComplete() {
		t.Fatal("fixture-grade roots must leave the sync-suspect fallback armed")
	}
	if !p.IsSyncPath(filepath.Join(home, "LocalOnly", "a.txt")).Sync {
		t.Fatal("under-profile write must stay sync-suspect while no confirmed root exists")
	}
	// Outside the profile is still not a sync channel.
	if p.IsSyncPath(`D:\plain\data.txt`).Sync {
		t.Fatal("path outside the user profile must not be suspect")
	}
}

// M-3: only env/registry/config-grade evidence may disarm the fallback, and a
// decoy (or a merely-existing default directory) can never do it.
func TestSyncFallbackNotDisarmableByWeakRoot(t *testing.T) {
	home, _, _ := sandbox(t)
	root := filepath.Join(home, "OneDrive")
	if err := os.MkdirAll(root, 0o755); err != nil { // confirmed roots are real dirs
		t.Fatal(err)
	}
	for _, src := range []string{"default", "fixture", "options", ""} {
		p := NewProvenance(ProvOptions{NoProbe: true, HomeDir: home,
			SyncRoots: []SyncRoot{{Provider: "OneDrive", Path: root, Source: src}}})
		if p.SyncDetectionComplete() {
			t.Errorf("source %q must not count as confirmed", src)
		}
		if !p.IsSyncPath(filepath.Join(home, "Documents", "exfil.md")).Sync {
			t.Errorf("source %q: under-profile fallback must stay armed", src)
		}
	}
	for _, src := range []string{"registry", "config", "env"} {
		p := NewProvenance(ProvOptions{NoProbe: true, HomeDir: home,
			SyncRoots: []SyncRoot{{Provider: "OneDrive", Path: root, Source: src}}})
		if !p.SyncDetectionComplete() {
			t.Errorf("source %q must count as confirmed", src)
		}
		if p.IsSyncPath(filepath.Join(home, "Documents", "exfil.md")).Sync {
			t.Errorf("source %q: confirmed root must lift the blanket suspect net", src)
		}
		// ...while the confirmed root itself keeps flagging.
		if !p.IsSyncPath(filepath.Join(root, "notes.md")).Sync {
			t.Errorf("source %q: writes under the confirmed root must be sync", src)
		}
	}
}

// M-3/N-2: a root whose own path C26 cannot verify is "not established", so it
// must not disarm the fallback either.
func TestSyncUnverifiedRootKeepsFallback(t *testing.T) {
	home, _, _ := sandbox(t)
	p := NewProvenance(ProvOptions{NoProbe: true, HomeDir: home, SyncRoots: []SyncRoot{
		{Provider: "OneDrive", Path: filepath.Join(home, "OneDrive", "..", "..", "definitely-missing-9f3a"), Source: "registry"},
	}})
	if p.SyncDetectionComplete() {
		t.Fatal("a registry root that C26 cannot canonicalize is not confirmed")
	}
	if !p.IsSyncPath(filepath.Join(home, "Documents", "x.md")).Sync {
		t.Fatal("suspect fallback must stay armed")
	}
}

// M-2: the OneDrive environment variables are a confirmed-grade location
// signal (this is what the machine actually publishes when HKCU Accounts has
// no UserFolder), and the registry probe logic is now injectable + covered.
func TestSyncEnvConfiguredRoots(t *testing.T) {
	home, _, _ := sandbox(t)
	root := filepath.Join(home, "OneDrive")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	roots := envConfiguredRoots(probeEnv{Home: home, OneDrive: root,
		OneDriveConsumer: root, OneDrivePublic: "  "})
	if len(roots) != 2 {
		t.Fatalf("expected OneDrive + OneDriveConsumer roots, got %+v", roots)
	}
	for _, r := range roots {
		if r.Source != "env" || r.Provider != "OneDrive" {
			t.Errorf("env roots must be graded env/OneDrive, got %+v", r)
		}
	}
	p := NewProvenance(ProvOptions{NoProbe: true, HomeDir: home, SyncRoots: roots})
	if !p.SyncDetectionComplete() {
		t.Fatal("env-grade roots are confirmed evidence")
	}
	if !p.IsSyncPath(filepath.Join(root, "Notes", "new.md")).Sync {
		t.Fatal("env root must flag writes under it")
	}
}

// fakeRegistry drives the portable half of the P12 registry probe.
type fakeRegistry struct {
	sub map[string][]string
	val map[string]string // "<key>|<valueName>" -> data
}

func (f fakeRegistry) subKeys(key string) []string { return f.sub[key] }
func (f fakeRegistry) string(key, name string) string {
	return f.val[key+`|`+name]
}

func TestSyncRegistryProbeTable(t *testing.T) {
	const accounts = `Software\Microsoft\OneDrive\Accounts`
	cases := []struct {
		name string
		fx   fakeRegistry
		want []SyncRoot
	}{
		{"empty hive yields nothing", fakeRegistry{}, nil},
		{"personal + business accounts, decoys ignored", fakeRegistry{
			sub: map[string][]string{accounts: {"Personal", "Business1", "Business2", "ConsumingAccounts", "PriceList"}},
			val: map[string]string{
				accounts + `\Personal|UserFolder`:          `D:\OneDrive`,
				accounts + `\Business1|UserFolder`:         `E:\Contoso`,
				accounts + `\Business2|UserFolder`:         "",         // configured-but-moved-out account
				accounts + `\ConsumingAccounts|UserFolder`: `C:\decoy`, // not an account node
			},
		}, []SyncRoot{
			{Provider: "OneDrive", Path: `D:\OneDrive`, Source: "registry"},
			{Provider: "OneDrive", Path: `E:\Contoso`, Source: "registry"},
		}},
		{"duplicate accounts dedupe", fakeRegistry{
			sub: map[string][]string{accounts: {"Personal", "Business9"}},
			val: map[string]string{
				accounts + `\Personal|UserFolder`:  `D:\OneDrive`,
				accounts + `\Business9|UserFolder`: `d:\onedrive`,
			},
		}, []SyncRoot{
			{Provider: "OneDrive", Path: `D:\OneDrive`, Source: "registry"},
		}},
		{"dropbox installer path", fakeRegistry{
			val: map[string]string{`Software\Dropbox\Dropbox|Path`: `C:\Users\x\Dropbox`},
		}, []SyncRoot{{Provider: "Dropbox", Path: `C:\Users\x\Dropbox`, Source: "registry"}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := registryProbeFor(c.fx)
			if len(got) != len(c.want) {
				t.Fatalf("roots=%+v want %+v", got, c.want)
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Errorf("root %d: got %+v want %+v", i, got[i], c.want[i])
				}
			}
		})
	}
	if got := registryProbeFor(nil); got != nil {
		t.Errorf("nil source (non-Windows) must yield no roots, got %+v", got)
	}
}

// P12 live registry evidence (skipped where the hive has nothing to say):
// asserts the probe's SHAPE, never machine state.
func TestSyncRegistryProbeLive(t *testing.T) {
	roots := registryProbe(probeEnv{Home: userHomeDir()})
	for _, r := range roots {
		if r.Source != "registry" || r.Provider == "" || r.Path == "" {
			t.Fatalf("live registry probe returned a malformed root %+v", r)
		}
		t.Logf("P12 evidence: registry-confirmed %s root %s", r.Provider, r.Path)
	}
	if len(roots) == 0 {
		t.Skip("no registry-grade sync record on this machine (HKCU Accounts without UserFolder is the documented reality here)")
	}
}

// N-5: the suspect net is component-bounded — a sibling directory whose name
// merely starts with the profile name must not be swept in.
func TestSyncSuspectFallbackIsComponentBounded(t *testing.T) {
	base := t.TempDir()
	home := filepath.Join(base, "profile")
	sibling := filepath.Join(base, "profileevil")
	for _, d := range []string{home, sibling} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	p := NewProvenance(ProvOptions{NoProbe: true, HomeDir: home})
	if !p.IsSyncPath(filepath.Join(home, "Documents", "x.txt")).Sync {
		t.Fatal("under-profile must be suspect")
	}
	if p.IsSyncPath(filepath.Join(sibling, "x.txt")).Sync {
		t.Fatal("sibling directory with a name prefix of the profile must not be suspect")
	}
}

// B-1 guard: the fix must not turn every legitimate new-file write into a
// sync channel. A brand-new file in a plain (non-sync) directory, with a
// confirmed root registered, is NOT sync and its content is not scanned.
func TestSyncNormalNewFileWriteNotFlagged(t *testing.T) {
	home, _, _ := sandbox(t)
	work := filepath.Join(home, "work")
	if err := os.MkdirAll(work, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(home, "OneDrive"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := NewProvenance(ProvOptions{NoProbe: true, HomeDir: home, SyncRoots: []SyncRoot{
		{Provider: "OneDrive", Path: filepath.Join(home, "OneDrive"), Source: "registry"}}})
	if !p.SyncDetectionComplete() {
		t.Fatal("sanity: registry-grade root is confirmed")
	}
	target := filepath.Join(work, "brand-new-file.md") // never created
	st := p.IsSyncPath(target)
	if st.Sync {
		t.Fatalf("normal new-file write into a plain dir must NOT be sync: %+v", st)
	}
	p.OpenScope("task-1")
	p.Mark("task-1", SrcFSRead, "/secret", marker)
	if _, ok := p.Inspect("task-1", "fs.write", map[string]any{"path": target, "content": marker}); ok {
		t.Fatal("a plain local write of tainted content must stay out of the R4 channel (SPEC-06 §5)")
	}
}

func TestSyncSuspectFallbackWhenUndetectable(t *testing.T) {
	// Zero confirmed detection -> every path under the user profile is
	// sync-suspect (safe default demanded by the ticket; MINOR N-1: the rule
	// is "anywhere in the profile", which is strictly broader than [fs]
	// allowed_dirs, so the option no longer carries a dead AllowedDirs field).
	home, _, _ := sandbox(t) // no sync dirs inside
	p := NewProvenance(ProvOptions{NoProbe: true, HomeDir: home})
	if p.SyncDetectionComplete() {
		t.Fatal("nothing was detectable")
	}
	st := p.IsSyncPath(filepath.Join(home, "Documents", "report.txt"))
	if !st.Sync || st.Root.Provider != "sync-suspect" {
		t.Fatalf("under-profile path must be sync-suspect, got %+v", st)
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
