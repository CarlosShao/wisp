//go:build !windows

// Ticket 119 AC#1: the collateral the ticket 113 leg produced, cut as cases that
// can be re-run instead of re-argued.
//
// What was measured (Linux container, `ln -s /realpriv /varlink` +
// TMPDIR=/varlink/w119tmp): 81 failing test lines across internal/winsec,
// internal/config, internal/agent and internal/memory, every one of them
// "winsec: refusing to seal ... reaches it through the link at /varlink", and
// four of them belonging to ticket 113's own AC#3 reverse half. Same packages,
// same container, TMPDIR in a real directory: zero failures. That is the shape
// this file pins.
//
// The semantics ticket 119 picked is option 2: winsec's rule is not narrowed.
// The layer that decides which tree to seal - internal/proc's data-root
// resolution and cmd/wisp's resolveDataDir, which is what `wisp run`,
// `wisp providers` and `wisp doctor` open the store from - resolves what the OS
// answered (proc.SealableRoot) and hands the floor a spelling that names a real
// tree. So the cases below drive those routes rather than re-implementing them:
// a test that only exercised winsec would stay green whether the caller learned
// anything, which is why this file imports internal/proc from inside
// internal/winsec's test binary (an external test package, so no cycle: the
// production edge stays winsec <- proc).
//
// The reverse halves matter as much here as they did in ticket 113: the root
// nobody resolved must still be refused (option 2 is not "the floor got
// softer"), and a link planted inside an already-resolved data root must still
// be refused, or the fix would have traded the leg away on exactly the route it
// touched.
//
// Three instrument notes, because all three have bitten this repository:
//   - these containers run as uid 0, and root is not stopped by a mode bit. The
//     0700 assertions below are about the mode the filesystem agreed to store on
//     the directory the seal created, which is what winsec verifies; they are not
//     a claim that another account was locked out, and nothing here should be
//     read as one;
//   - an expected value is written against the *real* tree (asked of the
//     filesystem with filepath.EvalSymlinks, see cleanSpelling119) or against the
//     link's own name - never against a spelling produced by the function the
//     case is testing. Ticket 119's acceptance measured what that costs: the two
//     fixtures that asked proc.SealableRoot for their expectations stayed green
//     while production was mutated into resolving the declared roots it is
//     supposed to leave alone, because SealableRoot is idempotent.
//   - which base a case stands on is part of its judgement, not a default. A case
//     whose subject is a link it plants itself must stand on a resolved base: on
//     a host whose temp dir is behind a link (macOS' real shape, and the
//     container's TMPDIR=/ac7link shape) the floor refuses the whole spelling at
//     that first component, the walk never reaches the planted link, and the
//     refusal the case reads belongs to somebody else. Ticket 137 AC#4 resolved
//     the roots of the cases that were not about the unresolved-root shape; the
//     two here that are (TestAC1POSIXUnresolvedSymlinkedRootStillRefused119 and
//     TestAC2POSIXInjectedTestDataDirStandsAsDeclared119) got the same treatment
//     from ticket 119 AC#7, together with the attribution check
//     refusalCreditsLink137, because root and reading only mean anything as a
//     pair.
package winsec_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/proc"
	"github.com/CarlosShao/wisp/internal/winsec"
)

// linkShape119 builds "<base>/<linkName> -> <base>/<realName>" under base and
// returns both spellings. The link is the OS's own shape (a directory the caller
// spells one way and the kernel reads another), which is the whole subject.
type linkShape119 struct {
	link string // spelled through the link, e.g. /tmp/.../varlink119
	real string // the tree the link names, /tmp/.../realpriv119
}

func newLinkShape119(t *testing.T, base, name string) linkShape119 {
	t.Helper()
	s := linkShape119{
		link: filepath.Join(base, name+"link119"),
		real: filepath.Join(base, name+"real119"),
	}
	if err := os.Mkdir(s.real, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(s.real, s.link); err != nil {
		t.Skipf("cannot build the symlink this case measures (%v): without it the shape does not exist", err)
	}
	if info, err := os.Lstat(s.link); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("Lstat(%s) after os.Symlink: info=%v err=%v", s.link, info, err)
	}
	return s
}

// spelledThrough returns the link-side spelling of a sub path.
func (s linkShape119) spelledThrough(sub string) string { return filepath.Join(s.link, sub) }

// inRealTree returns the real-tree spelling of the same sub path.
func (s linkShape119) inRealTree(sub string) string { return filepath.Join(s.real, sub) }

// assertSealedDir is the "the seal happened, on that directory" half: it checks
// through the link spelling so a reader can see the two names describing one
// object, and asserts the mode winsec verified.
func assertSealedDir119(t *testing.T, viaLink, real string) {
	t.Helper()
	info, err := os.Stat(viaLink)
	if err != nil {
		t.Fatalf("AC#1: %s (the tree the link names, seen through %s) is not there: %v", real, viaLink, err)
	}
	if !info.IsDir() {
		t.Fatalf("AC#1: %s is not a directory", real)
	}
	if got := info.Mode().Perm(); got != 0o700 {
		t.Errorf("AC#1: %s landed %o, not 0700 - the seal did not run on the resolved tree", real, got)
	}
}

// cleanSpelling119 asks the filesystem how the kernel spells a path that already
// exists. It is deliberately filepath.EvalSymlinks and not proc.SealableRoot:
// the latter is what several of these cases exist to test, and an expectation
// computed by the function under test cannot fail.
func cleanSpelling119(t *testing.T, path string) string {
	t.Helper()
	clean, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("EvalSymlinks(%s): %v", path, err)
	}
	return clean
}

// TestAC1POSIXSymlinkedTempDirRouteBecomesSealable119 is the WISP_ENV=test route
// with TMPDIR behind a symlink - the container shape, and the macOS shape (its
// per-user TMPDIR lives under /var, which is a symlink). Before ticket 119 the
// data root proc announced reached itself through the link, so the first
// sealing call refused and `wisp run`/`wisp providers` died in
// secret.NewStore / memory.Open.
func TestAC1POSIXSymlinkedTempDirRouteBecomesSealable119(t *testing.T) {
	base := t.TempDir()
	shape := newLinkShape119(t, base, "var")
	sub := shape.spelledThrough("w119tmp")
	if err := os.Mkdir(sub, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TMPDIR", sub)
	t.Setenv(proc.TestDataDirEnv, "") // take the os.TempDir() branch, not the injection

	root := proc.TestDataDir()
	t.Logf("AC#1 TMPDIR=%s (a link) -> proc.TestDataDir()=%s", sub, root)
	if strings.Contains(root, "varlink119") {
		t.Errorf("AC#1 RED: the data root is still spelled through the link (%s), so the caller declared an unresolved root and winsec has to refuse it", root)
	}
	if err := winsec.PrivateDirAll(root, 0o700); err != nil {
		t.Errorf("AC#1 RED: the test-env data root behind a symlinked TMPDIR could not be sealed: %v", err)
	}
	assertSealedDir119(t, filepath.Join(sub, filepath.Base(root)), root)
}

// TestAC1POSIXSymlinkedConfigDirRouteBecomesSealable119 is the route R-113-B
// registered only as a claim: dotfiles managers link $HOME/.config, and the
// dev/prod data root hangs off os.UserConfigDir. Same instrument as above, on
// the other OS read, so a half-fix that resolves TMPDIR and not the config root
// goes red here (ticket 119 AC#5).
func TestAC1POSIXSymlinkedConfigDirRouteBecomesSealable119(t *testing.T) {
	base := t.TempDir()
	shape := newLinkShape119(t, base, "home")
	homeViaLink := shape.spelledThrough("user")
	if err := os.Mkdir(homeViaLink, 0o700); err != nil {
		t.Fatal(err)
	}
	// XDG_CONFIG_HOME set but empty is how os.UserConfigDir picks $HOME/.config;
	// the value has to be cleared or a container that exports it wins, and then
	// this case would measure the wrong read.
	t.Setenv("HOME", homeViaLink)
	t.Setenv("XDG_CONFIG_HOME", "")

	cfg, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("AC#1: os.UserConfigDir with HOME=%s: %v", homeViaLink, err)
	}
	if cfg != filepath.Join(homeViaLink, ".config") {
		t.Fatalf("AC#1: the case is measuring the wrong config root: os.UserConfigDir()=%q", cfg)
	}
	sum, err := proc.Summarize(buildinfo.EnvProd)
	if err != nil {
		t.Fatalf("AC#1: proc.Summarize(prod): %v", err)
	}
	t.Logf("AC#1 HOME=%s (a link) -> prod data root=%s", homeViaLink, sum.DataDir)
	if strings.Contains(sum.DataDir, "homelink119") {
		t.Errorf("AC#1 RED: the prod data root is still spelled through the link (%s)", sum.DataDir)
	}
	if err := winsec.PrivateDirAll(sum.DataDir, 0o700); err != nil {
		t.Errorf("AC#1 RED: the config-dir data root behind a symlinked $HOME/.config could not be sealed: %v", err)
	}
	assertSealedDir119(t, filepath.Join(cfg, "wisp"), sum.DataDir)
}

// TestAC1POSIXUnresolvedSymlinkedRootStillRefused119 is today's other half, kept
// on purpose: a caller that hands the floor an unresolved root is refused, and
// the refusal names the link. That is the residual this ticket registers instead
// of swallowing - our own suites pass t.TempDir() straight to the floor, so on a
// host whose temp dir is behind a symlink they are red until their roots are
// resolved the same way (ticket 118/111's harness land, not option 2's).
//
// Ticket 119 AC#7 tightened the *reading* of that refusal, and it is a different
// claim than the one above: this case is now one of only two places in the
// package that still answer for the unresolved-root shape (ticket 137 AC#4 moved
// the other POSIX roots onto resolved spellings), so a PASS here has to be about
// the link this case planted. On a host whose temp dir is behind a link the floor
// refuses the whole spelling at that first component and never walks down to the
// planted link, which is why the two assertions below come as a pair:
//   - the base is resolved (cleanSpelling119, the same precedent as the AC#3 case
//     in this file), so the only link left in the spelling is the one planted
//     here, and the case can go green *and* go red;
//   - the refusal is read for attribution with refusalCreditsLink137, so it
//     cannot be answered by an ambient link above the tree.
//
// Either half on its own is worthless, and ticket 137 AC#1's MUT-D is the
// measurement that says so: with only the sentinel check this case stayed green
// through that mutation in the symlinked shape (the host link refunded the
// refusal), and with only the named assertion it would be red there whatever the
// implementation does.
func TestAC1POSIXUnresolvedSymlinkedRootStillRefused119(t *testing.T) {
	base := cleanSpelling119(t, t.TempDir())
	shape := newLinkShape119(t, base, "var")
	foreign := newLinkShape119(t, base, "fgn")
	unresolved := shape.spelledThrough("data")
	// Premise, not expectation: if anything above this tree were still a link,
	// every refusal below it would belong to that link and the attribution check
	// would be reading a sentence this case did not earn. Asked of the floor, and
	// the AC#3 case below uses the same gate for the same reason.
	if err := winsec.PrivateDirAll(filepath.Join(base, "premise"), 0o700); err != nil {
		t.Fatalf("premise broke: the resolved base %s is itself refused by the floor (%v), so the link planted here is not the only link in these spellings", base, err)
	}

	err := winsec.PrivateDirAll(unresolved, 0o700)
	t.Logf("AC#1 PrivateDirAll(%q) -> %v", unresolved, err)
	if err == nil {
		t.Errorf("AC#1 RED: PrivateDirAll(%q) was accepted, so the floor no longer refuses an unresolved root: the link this case planted at %s is invisible to it", unresolved, shape.link)
	} else if !errors.Is(err, winsec.ErrUnresolvedPath) {
		t.Errorf("AC#1: refusal did not name ErrUnresolvedPath: %v", err)
	} else if !refusalCreditsLink137(err, shape.link) {
		t.Errorf("AC#1 RED: the refusal of %q named ErrUnresolvedPath but did not credit the link this case planted at %q (%v): an ambient link above the tree answered for it, so this says nothing about the leg under test", unresolved, shape.link, err)
	}
	if _, statErr := os.Lstat(shape.inRealTree("data")); !errors.Is(statErr, fs.ErrNotExist) {
		t.Errorf("AC#1 RED: the refusal created %s anyway (%v)", shape.inRealTree("data"), statErr)
	}
	if _, statErr := os.Lstat(foreign.inRealTree("data")); !errors.Is(statErr, fs.ErrNotExist) {
		t.Errorf("AC#1 RED: the refusal wrote into a tree nobody named (%v)", statErr)
	}
}

// TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119 is ticket 119 AC#3: the
// reverse half has to stay pinned on the route the fix touches. Resolving the
// root removes the *system's* link from the ancestor chain; a link somebody
// plants inside the data root still carries the seal out of the tree this call
// names, and still gets refused with the foreign tree untouched.
func TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119(t *testing.T) {
	base := t.TempDir()
	// Declared the way option 2 says a caller declares it: the OS-side links in
	// the ambient temp dir are resolved away, so the case measures the link
	// inside the data root and not wherever the test binary happens to live. The
	// clean spelling is asked of the filesystem, not of proc.SealableRoot (see
	// cleanSpelling119).
	root := filepath.Join(cleanSpelling119(t, base), "data")
	if err := winsec.PrivateDirAll(root, 0o700); err != nil {
		t.Fatalf("AC#3: PrivateDirAll refused its own tree: %v", err)
	}
	victim := filepath.Join(base, "victim")
	if err := os.Mkdir(victim, 0o700); err != nil {
		t.Fatal(err)
	}
	keep := filepath.Join(victim, "keep-me.txt")
	if err := os.WriteFile(keep, []byte("not this tree's data"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(keep, 0o666); err != nil {
		t.Fatal(err)
	}
	before := statFact113(keep)
	if err := os.Symlink(victim, filepath.Join(root, "out")); err != nil {
		t.Skipf("no symlink here: %v", err)
	}

	spelled := filepath.Join(root, "out", "keep-me.txt")
	err := winsec.SealFile(spelled)
	t.Logf("AC#3 SealFile(%q) -> %v", spelled, err)
	if err == nil {
		t.Errorf("AC#3 RED: sealing through a link inside a resolved data root returned nil")
	} else if !errors.Is(err, winsec.ErrUnresolvedPath) {
		t.Errorf("AC#3: refusal did not name ErrUnresolvedPath: %v", err)
	}
	if after := statFact113(keep); after != before {
		t.Errorf("AC#3 RED: the foreign file was %s and is now %s", before, after)
	}
	// And the resolved tree's own file is still sealable, or the leg now over-refuses
	// on the very route option 2 is supposed to make work.
	mine := filepath.Join(root, "artifact.txt")
	if err := winsec.PrivateFile(mine, []byte("tool output"), 0o600); err != nil {
		t.Errorf("AC#3 RED: PrivateFile refused a plain path inside the resolved root: %v", err)
	}
}

// TestAC2POSIXInjectedTestDataDirStandsAsDeclared119 pins the line option 2 draws
// between an OS answer and a caller declaration (ticket 119's ruling on
// R-119-3): WISP_TEST_DATA_DIR is returned verbatim, so the harness that set it
// owns the spelling it declared and winsec answers it the way it answers any
// other caller - resolved spellings seal, links get refused. Rewriting an
// injected root would move the store out from under the harness that named it.
//
// Both directions are asserted, and the expectations come from the filesystem:
//   - a root declared *through* the link must come back exactly as declared (an
//     implementation that resolves injections, which is the discipline this case
//     exists to keep, returns the other spelling instead), and it must then be
//     refused by the floor and create nothing;
//   - the same tree declared by its clean spelling must come back exactly as
//     declared too, and must seal - so the rule is "as declared", not "refused".
//
// The first leg is why the fixture is spelled through the link. Ticket 119's
// acceptance ran production through the mutation "resolve the injected root as
// well" and this case read green, because its expectation used to be computed by
// proc.SealableRoot, which is idempotent: as declared and resolved-and-rejoined
// were the same string, so no implementation could contradict it.
//
// Ticket 119 AC#7 is the same disease on the other half of this case: the
// placement leg of leg one was only checked for its sentinel, so on a host whose
// temp dir is behind a link the floor refused leg one at that ambient component,
// the planted link was never reached, and the case stayed green through ticket
// 137 AC#1's MUT-D. Both medicines are in, for the reason spelled out on
// TestAC1POSIXUnresolvedSymlinkedRootStillRefused119: a resolved base, and a
// refusal read for which link it credits. Leg two is this case's premise gate -
// it requires the floor to seal a plain path under the same base, so an ambient
// link above the tree reddens the case rather than letting leg one pass for the
// wrong reason.
func TestAC2POSIXInjectedTestDataDirStandsAsDeclared119(t *testing.T) {
	base := cleanSpelling119(t, t.TempDir())
	shape := newLinkShape119(t, base, "inj")

	// Leg one: declared through the link, so "verbatim" and "resolved" are
	// different strings and only one of them can satisfy the assertion.
	declared := shape.spelledThrough(filepath.Join("harness", "picked"))
	asTheKernelSpellsIt := filepath.Join(cleanSpelling119(t, shape.link), "harness", "picked")
	if declared == asTheKernelSpellsIt {
		t.Fatalf("premise broke: %q is already the kernel's own spelling, so this leg could not tell an untouched root from a resolved one", declared)
	}
	t.Setenv(proc.TestDataDirEnv, declared)

	if got := proc.TestDataDir(); got != declared {
		t.Errorf("AC#2 RED: an injected data root was rewritten: got %q, want the declared %q. Rewriting it resolves the link the caller chose, which is exactly what option 2 refuses to do to a declared root (the OS-side reads, os.TempDir and os.UserConfigDir, are resolved elsewhere and have their own cases)", got, declared)
	}
	// The price of "as declared", pinned rather than implied: the floor answers a
	// root that reaches itself through a link with a refusal, and refuses without
	// creating anything in the tree the link names. Since ticket 119 AC#7 the
	// refusal is also read for *which* link it credits: leg one planted it, so an
	// answer bought by a link above the base is not this case's answer.
	err := winsec.PrivateDirAll(declared, 0o700)
	t.Logf("AC#2 PrivateDirAll(%q) -> %v", declared, err)
	if err == nil {
		t.Errorf("AC#2 RED: the floor accepted a declared root that reaches itself through the link at %s, so the placement leg is gone on the route option 2 leaves untouched", shape.link)
	} else if !errors.Is(err, winsec.ErrUnresolvedPath) {
		t.Errorf("AC#2: refusal did not name ErrUnresolvedPath: %v", err)
	} else if !refusalCreditsLink137(err, shape.link) {
		t.Errorf("AC#2 RED: the refusal of the declared root %q named ErrUnresolvedPath but did not credit the link this case planted at %q (%v): an ambient link above the base answered for it, so the walk under test never ran", declared, shape.link, err)
	}
	if _, statErr := os.Lstat(asTheKernelSpellsIt); !errors.Is(statErr, fs.ErrNotExist) {
		t.Errorf("AC#2 RED: the refusal created %s anyway (%v)", asTheKernelSpellsIt, statErr)
	}

	// Leg two: the same tree declared the way a harness on a normal system
	// declares it. Declared-and-clean still comes back untouched and still seals,
	// which is what keeps leg one from being read as "injections are refused".
	injected := filepath.Join(cleanSpelling119(t, shape.real), "harness", "picked")
	t.Setenv(proc.TestDataDirEnv, injected)
	if got := proc.TestDataDir(); got != injected {
		t.Errorf("AC#2 RED: a clean injected data root was rewritten: got %q, want the declared %q", got, injected)
	}
	if err := winsec.PrivateDirAll(injected, 0o700); err != nil {
		t.Errorf("AC#2 RED: a clean injected data root was refused: %v", err)
	}
	assertSealedDir119(t, injected, injected)
}
