//go:build windows

package risk

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/sys/windows"
)

// B-1 red team (adversarial report §4 B-1, §6 recipes P10/P12/P13): fs.write
// targets are usually brand-new files, where C26's handle query cannot run, so
// the pre-fix code compared the *lexical* fallback spelling against the sync
// roots — and `\\?\`-prefixed / 8.3-short-name spellings of a real sync root
// slipped through while the bytes landed inside that root.
//
// Every spelling below is produced by the OS (GetShortPathNameW), never typed
// by hand, and each case is executed twice: once asserting the gate verdict and
// once letting the gate decide whether the write may happen at all, then
// verifying on disk that no byte landed.

// tryShortPath returns the OS-assigned 8.3 spelling of p and whether the
// volume generated an alias at all (NtfsDisable8dot3NameCreation, or a name
// that is already 8.3-clean, means there is none). Spellings are always read
// back from the OS — never typed by hand.
func tryShortPath(p string) (string, bool) {
	p16, err := windows.UTF16PtrFromString(p)
	if err != nil {
		return "", false
	}
	n, err := windows.GetShortPathName(p16, nil, 0)
	if err != nil {
		return "", false
	}
	buf := make([]uint16, n)
	if _, err := windows.GetShortPathName(p16, &buf[0], n); err != nil {
		return "", false
	}
	short := windows.UTF16ToString(buf)
	if strings.EqualFold(short, p) {
		return short, false
	}
	return short, true
}

// shortPathOf is tryShortPath for the cases where an alias is the premise.
func shortPathOf(t *testing.T, p string) string {
	t.Helper()
	short, ok := tryShortPath(p)
	if !ok {
		t.Skipf("no 8.3 alias generated for %s on this volume", p)
	}
	return short
}

// syncSandbox builds <tmp>\profile with a real (existing) OneDrive root holding
// an existing Notes directory, registered as registry-grade evidence so the
// sync verdict can only come from ROOT MEMBERSHIP, never from the blanket
// under-profile fallback.
func syncSandbox(t *testing.T) (p *Provenance, home, root string) {
	t.Helper()
	base := t.TempDir()
	home = filepath.Join(base, "profile")
	root = filepath.Join(home, "OneDrive Copilot Sync") // long name => real 8.3 alias exists
	for _, d := range []string{root, filepath.Join(root, "Notes"), filepath.Join(home, "work")} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	p = NewProvenance(ProvOptions{NoProbe: true, HomeDir: home, SyncRoots: []SyncRoot{
		{Provider: "OneDrive", Path: root, Source: "registry"}}})
	if !p.SyncDetectionComplete() {
		t.Fatal("sanity: the injected registry-grade root must be confirmed (fallback must NOT be what flags these writes)")
	}
	return
}

func TestSyncRedTeamNewFileSpellingsDeniedAndLandNoBytes(t *testing.T) {
	p, home, root := syncSandbox(t)
	shortRoot := shortPathOf(t, root)
	t.Logf("sync root long  = %s", root)
	t.Logf("sync root 8.3   = %s", shortRoot)
	rel := filepath.Join(`Notes`, `out2.md`)

	p.OpenScope("task-1")
	if !p.Mark("task-1", SrcFSRead, filepath.Join(home, "secrets.txt"), "token "+marker) {
		t.Fatal("mark rejected")
	}

	cases := []struct{ name, target string }{
		{"plain long spelling", filepath.Join(root, rel)},
		{"extended-length \\\\?\\ spelling", `\\?\` + filepath.Join(root, rel)},
		{"OS 8.3 short-name spelling", filepath.Join(shortRoot, rel)},
		{"8.3 + \\\\?\\ combined", `\\?\` + filepath.Join(shortRoot, rel)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// The target must not exist yet — that is the whole point (a
			// pre-existing target was already caught before this fix).
			if _, err := os.Lstat(c.target); err == nil {
				t.Fatalf("test precondition broken: %s already exists", c.target)
			}
			// Precondition of the bypass: plain C26 on this spelling reports
			// "not resolved", i.e. the returned form is lexical and keeps the
			// \\?\ prefix / 8.3 alias that the sync-root comparison must see
			// expanded. This is what the ancestor anchoring now repairs.
			if res, err := Resolve(c.target, nil); err != nil || res.Resolved {
				t.Fatalf("test premise broken: Resolve(%q) resolved=%v err=%v", c.target, res.Resolved, err)
			}
			st := p.IsSyncPath(c.target)
			if !st.Sync {
				t.Fatalf("EXPLOITABLE: spelling judged non-sync: %+v", st)
			}
			if st.Root.Source == "suspect-fallback" {
				t.Fatalf("verdict must come from root membership, not the fallback: %+v", st)
			}
			if st.Root.Provider != "OneDrive" {
				t.Fatalf("wrong root attribution: %+v", st)
			}
			hit, ok := p.Inspect("task-1", "fs.write", map[string]any{"path": c.target, "content": "leak " + marker})
			if !ok {
				t.Fatal("EXPLOITABLE: fs.write of tainted content into a sync root did not hit")
			}
			if hit.Channel != ChSyncWrite {
				t.Errorf("channel: got %q want %q", hit.Channel, ChSyncWrite)
			}
			// Gate-controlled write: the caller writes only when the gate
			// allows. Prove nothing landed.
			if !fileExists(c.target) {
				if err := gatedWrite(p, c.target, "leak "+marker); err != nil {
					t.Fatalf("control write path broken: %v", err)
				}
			}
			if fileExists(c.target) {
				t.Fatalf("EXPLOITABLE: %d bytes landed inside the sync root", mustSize(t, c.target))
			}
		})
	}
}

// TestSyncRedTeamGuardPlainNewFileWriteAllowed is the other half of B-1: the
// fix must not flag every legitimate new-file write. The same spellings that
// defeat the sync test, aimed at a plain non-sync directory, must stay
// L0/L1 (no R4) and the bytes must land.
func TestSyncRedTeamGuardPlainNewFileWriteAllowed(t *testing.T) {
	p, home, root := syncSandbox(t)
	shortWork := shortPathOf(t, filepath.Join(home, "work"))
	t.Logf("non-sync dir 8.3 = %s", shortWork)

	p.OpenScope("task-1")
	p.Mark("task-1", SrcFSRead, filepath.Join(home, "secrets.txt"), "token "+marker)

	for _, target := range []string{
		filepath.Join(home, "work", "plain-new.md"),
		`\\?\` + filepath.Join(home, "work", "plain-ext.md"),
		filepath.Join(shortWork, "plain-83.md"),
	} {
		st := p.IsSyncPath(target)
		if st.Sync {
			t.Fatalf("false positive: normal new-file write flagged as sync: %+v", st)
		}
		if _, ok := p.Inspect("task-1", "fs.write", map[string]any{"path": target, "content": "fine " + marker}); ok {
			t.Fatalf("false positive: R4 on a plain local write (%s)", target)
		}
		// The gate really does let this one through (proof the deny above is
		// a verdict, not an I/O accident).
		if err := gatedWrite(p, target, "local bytes"); err != nil {
			t.Fatalf("gate blocked a legitimate local write: %v", err)
		}
		if !fileExists(target) {
			t.Fatalf("control write did not land: %s", target)
		}
	}
	// And the same engine still flags the sync root in the same run.
	if !p.IsSyncPath(filepath.Join(root, "Notes", "x.md")).Sync {
		t.Fatal("regression: plain sync-root write no longer flagged")
	}
}

// TestSyncRedTeamUnverifiableChainFailClosed: no existing ancestor at all
// (dead volume) cannot be anchored -> strict sync-suspect, never "unknown =
// allowed".
func TestSyncRedTeamUnverifiableChainFailClosed(t *testing.T) {
	p, _, _ := syncSandbox(t)
	for _, target := range []string{
		`Q:\definitely\not\here\out.md`,
		`\\host-that-does-not-exist\share\out.md`,
		`\\?\Q:\weird\out.md`,
	} {
		st := p.IsSyncPath(target)
		if !st.Sync || st.Root.Source != "suspect-fallback" {
			t.Errorf("unverifiable chain must fail closed as sync-suspect: %+v", st)
		}
	}
}

// TestSyncRedTeamRealOneDrive is the machine-reality leg of AC#4: the SAME
// bypass spellings aimed at the operator's actual, configured sync root.
// It never writes into the cloud folder — the assertion is that the gate
// denies, and that nothing appeared on disk afterwards.
func TestSyncRedTeamRealOneDrive(t *testing.T) {
	home := userHomeDir()
	if home == "" {
		t.Skip("no profile home")
	}
	p := NewProvenance(ProvOptions{}) // live env: env/registry/config/default probes
	var root string
	for _, r := range p.SyncRoots() {
		if dirExists(r.Path) {
			root = r.Path
			t.Logf("P12 evidence: live %s root via %s: %s", r.Provider, r.Source, r.Path)
			break
		}
	}
	if root == "" {
		t.Skipf("no live sync root on this machine (detected roots: %+v)", p.SyncRoots())
	}
	p.OpenScope("task-1")
	if !p.Mark("task-1", SrcFSRead, filepath.Join(home, "contract.pdf"), "clause "+marker) {
		t.Fatal("mark rejected")
	}
	rel := filepath.Join("Wisp-T19 Redteam Probe Dir", "out2.md")
	targets := []string{
		filepath.Join(root, rel),
		`\\?\` + filepath.Join(root, rel),
	}
	if short, ok := tryShortPath(root); ok {
		targets = append(targets, filepath.Join(short, rel))
	} else {
		t.Logf("no 8.3 alias exists for the live root %s (name is already 8.3-clean); the \\\\?\\ legs below still cross the boundary", root)
	}
	if len(targets) < 2 {
		t.Fatalf("expected at least the long + \\\\?\\ spellings, got %d", len(targets))
	}
	for _, target := range targets {
		if fileExists(target) {
			t.Fatalf("precondition: probe path already exists: %s", target)
		}
		st := p.IsSyncPath(target)
		if !st.Sync {
			t.Fatalf("EXPLOITABLE on the real sync root: %+v judged non-sync", st)
		}
		t.Logf("real root spelling judged sync=%v root=%s/%s why=%q", st.Sync, st.Root.Provider, st.Root.Source, st.Why)
		hit, ok := p.Inspect("task-1", "fs.write", map[string]any{"path": target, "content": "leak " + marker})
		if !ok {
			t.Fatalf("EXPLOITABLE: real sync root write not an R4 channel for %q", target)
		}
		t.Logf("R4 hit channel=%s source=%s", hit.Channel, hit.Source())
		if err := gatedWrite(p, target, "leak "+marker); err != nil {
			t.Fatalf("gate-controlled write failed: %v", err)
		}
		if fileExists(target) {
			t.Fatalf("EXPLOITABLE: bytes landed in the real sync root at %s", target)
		}
	}
}

// gatedWrite performs the write only if the C25 gate allows it — mirroring the
// call site the loop wiring (tickets 20/21) will have. It returns an error if
// the write was attempted and failed, so a "no bytes landed" assertion cannot
// be satisfied by an unrelated I/O error.
func gatedWrite(p *Provenance, target, content string) error {
	if _, hit := p.Inspect("task-1", "fs.write", map[string]any{"path": target, "content": content}); hit {
		return nil // denied by R4: nothing is written, which is the point
	}
	return os.WriteFile(target, []byte(content), 0o644)
}

func fileExists(p string) bool {
	_, err := os.Lstat(p)
	return err == nil
}

func mustSize(t *testing.T, p string) int64 {
	t.Helper()
	st, err := os.Lstat(p)
	if err != nil {
		return -1
	}
	return st.Size()
}
