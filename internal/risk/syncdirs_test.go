package risk

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
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

// syncEnv is the P12 fixture of ticket 82: one injected sync root, one plain
// directory, and a profile home that contains NEITHER of them.
type syncEnv struct {
	p           *Provenance
	home        string // fake profile, empty on purpose (see membershipEngine)
	cloud       string // the sibling tree that holds root and work
	root        string // existing dir, injected as `grade` evidence
	work        string // existing plain dir, no client anywhere near it
	syncTarget  string // brand-new file under root   -> sync BY MEMBERSHIP
	plainTarget string // brand-new file under work   -> not sync, same route
}

// membershipEngine builds the fixture whose two verdicts are decidable on EVERY
// platform, which is what AC#1 of ticket 82 needed before any of this family
// could be called Windows-only.
//
// The older fixtures (m7Engine's pre-ticket-82 shape, syncSandbox in
// syncdirs_redteam_windows_test.go) put the root and the plain directory INSIDE
// the profile, so "not a sync write" there could only mean "the under-profile
// suspect net was disarmed", and disarming it takes a root that is both
// confirmed-GRADE and C26-CANONICAL (syncSet.finalize). Canonicalization is the
// half POSIX does not have: syncSet.add only sets canonical=true when Resolve
// reports Resolved, and pathresolver_other.go's resolveHandle is the
// DEFERRED(macOS/Linux) stub that ticket 55 owns. So that precondition has no
// object on POSIX — it is not an assertion that happens to fail, there is
// nothing there to assert on.
//
// Moving both directories out of the profile removes the net from the
// decision without removing anything from the subject: match() compares the
// roots BEFORE the fallback and the fallback only covers paths under the
// profile, so membership is the only thing that can flip either verdict here,
// on either platform. The preconditions below assert that by ATTRIBUTION
// (Root.Source), so today's green can never be read as "the fallback decided
// this" — the mistake ticket 75's report made illegal.
//
// grade is the Source the root is injected with, because which grades count as
// confirmed is itself tiered: the portable cases pass "registry", the POSIX
// tier sweeps every grade name to pin that none of them is actionable there.
func membershipEngine(t *testing.T, grade string) syncEnv {
	t.Helper()
	base := t.TempDir()
	e := syncEnv{
		home:  filepath.Join(base, "profile"),
		cloud: filepath.Join(base, "cloud"),
	}
	e.root = filepath.Join(e.cloud, "OneDrive")
	e.work = filepath.Join(e.cloud, "work")
	for _, d := range []string{e.home, e.root, e.work} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	e.syncTarget = filepath.Join(e.root, "Notes", "out.md") // never created
	e.plainTarget = filepath.Join(e.work, "brand-new.md")   // never created
	e.p = NewProvenance(ProvOptions{NoProbe: true, HomeDir: e.home, SyncRoots: []SyncRoot{
		{Provider: "OneDrive", Path: e.root, Source: grade},
	}})
	if st := e.p.IsSyncPath(e.syncTarget); !st.Sync || st.Root.Source == "suspect-fallback" {
		t.Fatalf("precondition: %s must be sync by ROOT MEMBERSHIP, not by the fallback: %+v", e.syncTarget, st)
	}
	if st := e.p.IsSyncPath(e.plainTarget); st.Sync {
		t.Fatalf("precondition: %s must be judged a plain non-sync target: %+v", e.plainTarget, st)
	}
	return e
}

// outOfProfile builds a real, absolute path that is NOT under the fake profile,
// with the platform's own separator. Two tests in this file used to write
// `D:\plain\data.txt` for that job: on Windows it is a genuine out-of-profile
// path, and on POSIX it is a backslash-named file that happens to be RELATIVE,
// so the assertion passed without ever reaching the "outside the profile"
// concept the case is about. (Spellings that genuinely need a Windows-shaped
// input belong in a Windows-tier file — syncdirs_redteam_windows_test.go is
// where the `\\?\` and 8.3 forms live — and TestSyncRegistryProbeTable's
// "nil source (non-Windows) must yield no roots" case is the portable, honest
// statement of the same fact about the probe itself.)
func outOfProfile(t *testing.T, home string) string {
	t.Helper()
	other := filepath.Join(filepath.Dir(home), "outside-profile")
	if err := os.MkdirAll(other, 0o755); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(other, "plain", "data.txt") // never created, like a real fs.write target
	if isUnder(normPath(p), normPath(home)) {
		t.Fatalf("test premise broken: %q is under the profile %q", p, home)
	}
	return p
}

// TestSyncDefaultLocationProbe asserts the weakest grade stays a directory
// probe: it is only about a directory that exists under the profile.
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
	if p.IsSyncPath(outOfProfile(t, home)).Sync {
		t.Fatal("path outside the user profile must not be suspect")
	}
}

// M-3, portable half: a weak-grade root (a directory that merely exists at a
// default location, a fixture line, an injected option) is user-supplied data,
// not proof that a client syncs there, so it must NEVER disarm the
// under-profile suspect net. That half has an object wherever the net is armed,
// i.e. on both platforms.
//
// The MIRROR half — which grades DO disarm it — is the part that is not
// portable today, and ticket 82 AC#1 refuses to pretend otherwise: reaching
// "confirmed" needs a C26-canonical root, and on POSIX nothing is ever
// canonical until ticket 55 lands. That half now lives in
// syncdirs_windows_test.go (TestSyncConfirmedGradeDisarmsFallbackWindows), and
// syncdirs_other_test.go asserts the POSIX reality (no grade name at all
// disarms it here, fail-closed) instead of keeping a Windows-only expectation
// as a permanent red.

func TestSyncFallbackNotDisarmableByWeakRoot(t *testing.T) {
	home, _, _ := sandbox(t)
	root := filepath.Join(home, "OneDrive")
	if err := os.MkdirAll(root, 0o755); err != nil { // confirmed roots are real dirs
		t.Fatal(err)
	}
	for _, src := range []string{"default", "fixture", "options", ""} {
		p := NewProvenance(ProvOptions{
			NoProbe: true, HomeDir: home,
			SyncRoots: []SyncRoot{{Provider: "OneDrive", Path: root, Source: src}},
		})
		if p.SyncDetectionComplete() {
			t.Errorf("source %q must not count as confirmed", src)
		}
		if !p.IsSyncPath(filepath.Join(home, "Documents", "exfil.md")).Sync {
			t.Errorf("source %q: under-profile fallback must stay armed", src)
		}
		// The injected root itself flags whichever grade it came with.
		if !p.IsSyncPath(filepath.Join(root, "notes.md")).Sync {
			t.Errorf("source %q: writes under the injected root must be sync", src)
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

// M-2: the OneDrive environment variables are a confirmed-grade LOCATION signal
// (this is what the machine actually publishes when HKCU Accounts has no
// UserFolder), and the registry probe logic is now injectable + covered.
//
// Ticket 82 AC#1 keeps the portable half here: what envConfiguredRoots returns
// is a pure function of the environment it is handed (grade, provider, the
// blank-value trim), and the verdict for a write under that root is decided by
// ROOT MEMBERSHIP, which match() checks before the fallback, so it is assertable
// on POSIX too. What is NOT assertable there is the grade's consequence — that
// env evidence switches the under-profile net off — because on POSIX no root is
// ever C26-canonical: that consequence is pinned per platform by
// TestSyncConfirmedGradeDisarmsFallbackWindows (windows) and
// TestSyncNoGradeIsConfirmedOnPosix (!windows).
func TestSyncEnvConfiguredRoots(t *testing.T) {
	home, _, _ := sandbox(t)
	root := filepath.Join(home, "OneDrive")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	roots := envConfiguredRoots(probeEnv{
		Home: home, OneDrive: root,
		OneDriveConsumer: root, OneDrivePublic: "  ",
	})
	if len(roots) != 2 {
		t.Fatalf("expected OneDrive + OneDriveConsumer roots, got %+v", roots)
	}
	for _, r := range roots {
		if r.Source != "env" || r.Provider != "OneDrive" {
			t.Errorf("env roots must be graded env/OneDrive, got %+v", r)
		}
	}
	p := NewProvenance(ProvOptions{NoProbe: true, HomeDir: home, SyncRoots: roots})
	if st := p.IsSyncPath(filepath.Join(root, "Notes", "new.md")); !st.Sync || st.Root.Source != "env" {
		t.Fatalf("env root must flag writes under it BY MEMBERSHIP (not by the suspect net), got %+v", st)
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

// TestSyncRegistryProbeLive (the P12 live registry evidence case) lives in
// syncdirs_windows_test.go as of ticket 93 AC#2: the registry hive it probes
// does not exist on POSIX, so keeping it here meant a `--- SKIP` that bare
// `go test` booked as `ok` on every CI run. Body moved verbatim; nothing about
// its assertions changed.

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
// registry-grade root registered, is NOT sync and its content is not scanned.
//
// Ticket 82 AC#1: this ran on Windows only until now because the fixture kept
// the plain directory inside the profile and required SyncDetectionComplete()
// before "not sync" could mean anything; on POSIX that precondition has no
// object, so the case died at its own sanity check. With membershipEngine the
// subject (a new file that is not under any root stays out of channel R4) is
// asserted on both platforms, and the positive control below keeps the
// exemption from being "nothing is ever sync".
func TestSyncNormalNewFileWriteNotFlagged(t *testing.T) {
	e := membershipEngine(t, "registry")
	if st := e.p.IsSyncPath(e.plainTarget); st.Sync {
		t.Fatalf("normal new-file write into a plain dir must NOT be sync: %+v", st)
	}
	e.p.OpenScope("task-1")
	e.p.Mark("task-1", SrcFSRead, filepath.Join(e.home, "secret.txt"), marker)
	if _, ok := e.p.Inspect("task-1", "fs.write", map[string]any{"path": e.plainTarget, "content": marker}); ok {
		t.Fatal("a plain local write of tainted content must stay out of the R4 channel (SPEC-06 §5)")
	}
	// Positive control on the same engine: the only thing that flips the verdict
	// is which directory the bytes land in.
	if _, ok := e.p.Inspect("task-1", "fs.write", map[string]any{"path": e.syncTarget, "content": marker}); !ok {
		t.Fatal("ESCAPIABLE: the same tainted bytes aimed at the injected root must be an R4 channel")
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
	if p.IsSyncPath(outOfProfile(t, home)).Sync {
		t.Fatal("path outside the user profile must not be suspect")
	}
}

func TestSyncUnresolvablePathFailClosed(t *testing.T) {
	home, _, _ := sandbox(t)
	p := NewProvenance(ProvOptions{
		NoProbe: true, HomeDir: home,
		SyncRoots: []SyncRoot{{Provider: "OneDrive", Path: filepath.Join(home, "OneDrive")}},
	})
	// Empty/absent path spelling: unverifiable -> sync (strict side).
	if !p.IsSyncPath("").Sync {
		t.Fatal("empty path must fail-closed as sync")
	}
}

// TestSyncDotDotTailFailsClosed is the N-10 hardening from the re-verification
// (R2/R5): `link\..\x` is folded by lexical cleaning BEFORE reparseComponents
// audits the chain, so the junction that stood where `link\..` used to be is
// invisible to the classifier. That was only safe because Windows folds `.`
// and `..` the same way before it opens the file — an OS property, not an
// invariant of this code. It is now one: a spelling carrying a folded `..`
// fails closed as sync-suspect instead of being classified from a string that
// no longer describes the write.
//
// Ticket 82 AC#1: the hardening itself is lexical (hasFoldedDotDot runs on the
// raw spelling before any OS is consulted), so it has an object on POSIX as
// well — but the fixture had to change for this case to prove it. The old one
// put `work` inside the profile and required SyncDetectionComplete() so that the
// "plain write stays clean" control could not be satisfied by luck; on POSIX
// that precondition has no object, so the case died before its first assertion.
// membershipEngine puts both directories outside the profile, which makes the
// negative control mean non-membership on either platform.
func TestSyncDotDotTailFailsClosed(t *testing.T) {
	e := membershipEngine(t, "registry")
	sep := string(filepath.Separator)
	folded := filepath.Join(e.work, "x.md") // what the raw spelling below cleans to
	raw := strings.Join([]string{e.root, "Notes", "..", "..", "work", "x.md"}, sep)
	if filepath.Clean(raw) != folded {
		t.Fatalf("test premise broken: %q does not clean to %q", raw, folded)
	}
	if st := e.p.IsSyncPath(folded); st.Sync {
		t.Fatalf("precondition broken: the folded target must be a plain non-sync write: %+v", st)
	}
	st := e.p.IsSyncPath(raw)
	if !st.Sync {
		t.Fatalf("ESCAPIABLE (N-10): a spelling with a folded `..` was classified anyway: %+v", st)
	}
	if !strings.Contains(st.Why, "..") {
		t.Errorf("Why must name the '..' hardening, got %q", st.Why)
	}
	// The exfil gate follows the verdict, and the plain spelling is unaffected
	// (this is the false-positive side: a normal new-file write stays clean).
	e.p.OpenScope("task-1")
	if !e.p.Mark("task-1", SrcFSRead, filepath.Join(e.home, "secret.txt"), marker) {
		t.Fatal("mark rejected")
	}
	if _, ok := e.p.Inspect("task-1", "fs.write", map[string]any{"path": raw, "content": marker}); !ok {
		t.Error("ESCAPIABLE (N-10): fs.write through a folded `..` escaped the sync gate")
	}
	if _, ok := e.p.Inspect("task-1", "fs.write", map[string]any{"path": folded, "content": marker}); ok {
		t.Error("false positive: plain new-file write flagged after the N-10 hardening")
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
