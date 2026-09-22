//go:build !windows

// Ticket 119 rework (R-119-1): the `wisp secret` command surface is the second
// production reader of the same OS answer ticket 119's option 2 taught
// internal/proc and cmd/wisp's resolveDataDir to resolve.
//
// `wisp secret` never goes through proc.DefaultLayout: resolveSecretLayout
// reads os.UserConfigDir() itself and hands the value straight to
// proc.LayoutFor. So the resolution ticket 119 installed in the proc layer did
// not reach it, and the acceptor measured the consequence on a real binary:
// `wisp doctor` printed the resolved tree while `wisp secret list` still died
// with "winsec: refusing to seal .../varlink/home/.config/wisp-dev/secrets"
// (rc=1, dev, HOME / $HOME/.config / XDG_CONFIG_HOME behind a symlink).
//
// These cases drive resolveSecretLayout, i.e. the route itself, and they
// observe the outcome through secret.NewStore - the exact call openStore makes,
// which is where the refusal was produced. The Windows-only part of the store
// is the protector, reached by Store/Resolve, so NewStore is a legitimate call
// on this platform; a test that only looked at the returned string could not
// see the failure this ticket is about.
//
// Every expected value here is computed from the filesystem (os.Lstat,
// os.SameFile and the platform's own filepath.EvalSymlinks), never from
// proc.SealableRoot: the rule written in internal/proc/envfork_test.go ("an
// expectation computed by the function under test cannot fail") is the reason
// the acceptor could resolve declared roots as well and still read 56/56 green.
// The premise assertions below therefore fail the case outright when the two
// spellings happen to name one string, because a case that cannot tell them
// apart would prove nothing.
//
// Ruling R-119-3, pinned here and by
// TestAC2POSIXInjectedTestDataDirStandsAsDeclared119 in internal/winsec: a value
// this process obtains by *asking the OS* (os.UserConfigDir, os.TempDir) is
// resolved, whatever environment variable wrote it, because a data root is a
// location contract about a tree; a value this process was *handed by name*
// (WISP_TEST_DATA_DIR) is returned verbatim, because the harness that set it
// compares strings with it.
package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/secret"
	"github.com/CarlosShao/wisp/internal/winsec"
)

// plantedLink119b is a directory the caller spells one way and the kernel reads
// another: the shape this whole ticket family is about.
type plantedLink119b struct {
	link  string // <base>/<name>link119, a symlink to real
	real  string // <base>/<name>real119, the tree the link names
	clean string // the kernel's own spelling of that same tree
}

// plantLink119b creates "<base>/<name>real119" and points
// "<base>/<name>link119" at it. It returns once the filesystem has confirmed
// both that the link is a link and that the two spellings name one object.
func plantLink119b(t *testing.T, base, name string) plantedLink119b {
	t.Helper()
	p := plantedLink119b{
		link: filepath.Join(base, name+"link119"),
		real: filepath.Join(base, name+"real119"),
	}
	if err := os.Mkdir(p.real, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(p.real, p.link); err != nil {
		t.Skipf("cannot build the symlink this case measures (%v): without it the shape does not exist", err)
	}
	info, err := os.Lstat(p.link)
	if err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("Lstat(%s) after os.Symlink: info=%v err=%v", p.link, info, err)
	}
	p.clean, err = filepath.EvalSymlinks(p.link)
	if err != nil {
		t.Fatalf("EvalSymlinks(%s): %v", p.link, err)
	}
	// Premise, read from the filesystem: the two spellings must name one object
	// and must not be the same string, or nothing below can tell resolution from
	// identity - the false green this repository has been bounced for twice.
	if p.clean == p.link {
		t.Fatalf("premise broke: %q is spelled the same resolved and unresolved, so this case would assert nothing", p.link)
	}
	a, err := os.Stat(p.link)
	if err != nil {
		t.Fatalf("stat through the link %s: %v", p.link, err)
	}
	b, err := os.Stat(p.real)
	if err != nil {
		t.Fatalf("stat the named tree %s: %v", p.real, err)
	}
	if !os.SameFile(a, b) {
		t.Fatalf("premise broke: %s and %s are not the same object", p.link, p.real)
	}
	return p
}

// assertSecretRouteSeals runs the production composition (resolveSecretLayout
// then secret.NewStore) against the config root the OS answered, and pins both
// halves of the fix: the layout names the tree the kernel reads, and the store's
// first sealing call is accepted.
func assertSecretRouteSeals119b(t *testing.T, declared, wantCfg string) {
	t.Helper()
	l, err := resolveSecretLayout(buildinfo.EnvDev)
	if err != nil {
		t.Fatalf("resolveSecretLayout(dev): %v", err)
	}
	if l.Env != buildinfo.EnvDev {
		t.Fatalf("layout env = %q, want dev", l.Env)
	}
	if l.Portable {
		t.Fatalf("layout came back portable=true: a portable.txt next to the test binary would move this case's data root")
	}
	wantRoot := filepath.Join(wantCfg, "wisp-dev")
	t.Logf("os.UserConfigDir answered %s -> resolveSecretLayout(dev).DataDir = %s", declared, l.DataDir)
	if l.DataDir != wantRoot {
		t.Errorf("AC#1 RED: the secret route's data root is %q, want the tree the OS reads (%q); the other spelling is reached through a link, which is what winsec's placement floor (ticket 113) refuses",
			l.DataDir, wantRoot)
	}
	st, err := secret.NewStore(l.DataDir)
	if err != nil {
		t.Fatalf("AC#1 RED: secret.NewStore refused the dev data root the secret route announced: %v", err)
	}
	if st.Dir() != filepath.Join(wantRoot, "secrets") {
		t.Errorf("AC#1: the store opened %q, want %q", st.Dir(), filepath.Join(wantRoot, "secrets"))
	}
	// The bytes landed in the tree the OS named, and the link-side spelling is a
	// name for that same object rather than the address of a second store.
	opened, err := os.Stat(st.Dir())
	if err != nil {
		t.Fatalf("stat %s: %v", st.Dir(), err)
	}
	seen, err := os.Stat(filepath.Join(declared, "wisp-dev", "secrets"))
	if err != nil {
		t.Fatalf("AC#1 RED: %s was never created, so the seal did not land in the tree this call names: %v", filepath.Join(declared, "wisp-dev", "secrets"), err)
	}
	if !os.SameFile(opened, seen) {
		t.Errorf("AC#1 RED: %s and %s are not one object - the route wrote a second tree", st.Dir(), declared)
	}
	if got := opened.Mode().Perm(); got != 0o700 {
		t.Errorf("AC#1 RED: %s holds mode %o, want 0700 - winsec verifies that mode and refuses to promise a seal without it", st.Dir(), got)
	}
}

// TestAC1POSIXSecretRouteSymlinkedHomeBecomesSealable119 is the dev shape the
// acceptor measured as rc=1 before and after ticket 119: $HOME spelled through a
// symlink, so os.UserConfigDir() answers through it ($HOME/.config).
func TestAC1POSIXSecretRouteSymlinkedHomeBecomesSealable119(t *testing.T) {
	base := t.TempDir()
	p := plantLink119b(t, base, "home")
	viaCfg := filepath.Join(p.link, ".config")
	cleanCfg := filepath.Join(p.clean, ".config")
	if err := os.Mkdir(viaCfg, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", p.link)
	t.Setenv("XDG_CONFIG_HOME", "") // empty, not unset: this is how UserConfigDir picks $HOME/.config

	cfg, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("os.UserConfigDir with HOME=%s: %v", p.link, err)
	}
	if cfg != viaCfg {
		t.Fatalf("premise broke: this case measures the $HOME/.config branch but os.UserConfigDir() answered %q", cfg)
	}
	assertSecretRouteSeals119b(t, cfg, cleanCfg)
}

// TestAC1POSIXSecretRouteSymlinkedXDGConfigHomeBecomesSealable119 is the other
// branch of the same OS read. XDG_CONFIG_HOME is a value somebody else declared,
// and ruling R-119-3 says an OS answer is resolved whichever variable carried it,
// because the layout is a location contract. Without this leg a fix that
// special-cased $HOME would still read green.
func TestAC1POSIXSecretRouteSymlinkedXDGConfigHomeBecomesSealable119(t *testing.T) {
	base := t.TempDir()
	p := plantLink119b(t, base, "xdg")
	t.Setenv("HOME", filepath.Join(base, "home-not-asked"))
	t.Setenv("XDG_CONFIG_HOME", p.link)

	cfg, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("os.UserConfigDir with XDG_CONFIG_HOME=%s: %v", p.link, err)
	}
	if cfg != p.link {
		t.Fatalf("premise broke: os.UserConfigDir() answered %q, want the declared %q", cfg, p.link)
	}
	assertSecretRouteSeals119b(t, cfg, p.clean)
}

// TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119 is the AC#3 reverse
// half on the route this rework touches. Resolving the config root takes the
// system's link out of the ancestor chain; a link planted inside the data root
// the route announced still carries a seal out of the tree the call names, and
// must still be refused - otherwise "make `wisp secret` work" would have bought
// itself by softening the floor, the move ticket 107 was returned for.
func TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119(t *testing.T) {
	base := t.TempDir()
	p := plantLink119b(t, base, "home")
	if err := os.Mkdir(filepath.Join(p.link, ".config"), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", p.link)
	t.Setenv("XDG_CONFIG_HOME", "")

	l, err := resolveSecretLayout(buildinfo.EnvDev)
	if err != nil {
		t.Fatalf("resolveSecretLayout(dev): %v", err)
	}
	if err := winsec.PrivateDirAll(filepath.Join(l.DataDir, "secrets"), 0o700); err != nil {
		t.Fatalf("AC#3: the route refused to seal its own tree: %v", err)
	}
	foreign := filepath.Join(base, "someone-elses-tree")
	if err := os.MkdirAll(foreign, 0o700); err != nil {
		t.Fatal(err)
	}
	keep := filepath.Join(foreign, "keep-me.txt")
	if err := os.WriteFile(keep, []byte("not this tree's data"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(keep, 0o666); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(keep)
	if err != nil {
		t.Fatal(err)
	}
	planted := filepath.Join(l.DataDir, "secrets2")
	if err := os.Symlink(foreign, planted); err != nil {
		t.Skipf("no symlink inside the data root: %v", err)
	}

	err = winsec.PrivateDirAll(planted, 0o700)
	t.Logf("AC#3 PrivateDirAll(%q) -> %v", planted, err)
	if err == nil {
		t.Errorf("AC#3 RED: sealing through a link planted inside the resolved data root returned nil")
	} else if !errors.Is(err, winsec.ErrUnresolvedPath) {
		t.Errorf("AC#3: refusal did not name ErrUnresolvedPath: %v", err)
	}
	after, err := os.Stat(keep)
	if err != nil {
		t.Fatalf("the foreign file is gone: %v", err)
	}
	if !os.SameFile(before, after) || before.Mode() != after.Mode() {
		t.Errorf("AC#3 RED: %s was %s and is now %s - the seal left the named tree", keep, before.Mode(), after.Mode())
	}
	// And the route's own tree is still sealable next to the planted link, or the
	// leg now over-refuses on exactly the shape option 2 exists for.
	if err := winsec.PrivateFile(filepath.Join(l.DataDir, "artifact.txt"), []byte("tool output"), 0o600); err != nil {
		t.Errorf("AC#3 RED: PrivateFile refused a plain path inside the resolved root: %v", err)
	}
	if _, statErr := os.Lstat(filepath.Join(foreign, "artifact.txt")); !errors.Is(statErr, fs.ErrNotExist) {
		t.Errorf("AC#3 RED: the artifact appeared in the foreign tree instead (%v)", statErr)
	}
}
