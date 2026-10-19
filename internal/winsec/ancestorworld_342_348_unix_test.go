//go:build unix

// Tickets 342 AC#1 and 348 AC#2/AC#3, leg 342-r1: the shared world, racer and liveness proof for
// the two ANCESTOR-leg instruments. One world, two separate gates - each test file judges its own
// numerator, because the two tickets book two different damages, and the shape of the flip decides
// which of the two answers it is testing (see the two target forms below).
//
// WHY THESE INSTRUMENTS HAD TO BE NEW. 342-a1 and 342-v1 both measured that nothing on the shelf
// drives PrivateFile through a swapped ANCESTOR: TestPOSIX120SealWindowLeafSwapIsClosed exercises
// SealFile only, and TestAC1POSIX…322r1's liveness key is the LEAF's ELOOP. A green suite therefore
// says nothing about this leg, which is the reason the two tickets exist.
//
// THE WORLD. The caller names `<base>/data/sub/leaf-342r1.bin` and owns a real file there. `sub`
// flips, by rename only, between that real directory and a SYMLINK into somebody else's tree, and
// the outside tree carries one witness under the SAME leaf name - so in both states the spelling
// resolves to a file and the only question is WHICH file. The witness is sampled after every
// attempt and then restocked, so each attempt gets its own reading and one witness serves 1500
// attempts.
//
// TWO TARGET FORMS, BECAUSE THEY TEST DIFFERENT ANSWERS. Go 1.27's os.Root (which is what this
// package's pinPlacement is built from) refuses a link whose target is absolute - the Root is rooted
// at "/", so the target escapes - and follows a link whose target stays inside that root (winsec
// _other.go states this: "os.Root does not enforce this package's floor rule - it follows a symbolic
// link whose target stays inside the root"). So:
//
//   - ancTargetRel: a BARE sibling name, `someone-elses-tree`, with no ".." in it. This is the shape
//     the stdlib follows and only sameObjectPinned can refuse: doInRoot's escape rule fires on a ".."
//     element (measured on this leg's first run - a `../x` target answered "path escapes from parent"
//     and never reached the identity gate), while a bare in-root target is resolved to another object
//     inside the same root and handed back as if it were the one pinned. The OLD code walks through it
//     with nothing to refuse: open, seal and write all reach the foreign object.
//   - ancTargetAbs: the outside tree's absolute spelling. The pinned walk refuses it at Root's escape
//     face BEFORE any open, so nothing in the foreign tree is ever created; the OLD code opens
//     through it, then the seal's own walk refuses, and only then does the refusal-path unlink fire -
//     which is ticket 348's whole shape, deterministically (342-a1's A3: unlink err=<nil>, witness
//     gone). This is the form that gives the deleted numerator its teeth.
//
// WHY THE FLIP IS CLEAR-THEN-RENAME. rename(2) cannot replace an existing DIRECTORY with a symlink
// (ENOTDIR/EISDIR), so the outgoing state is removed first and the incoming one renamed in from a
// staged name. That leaves a real gap during which the caller's path resolves to nothing; attempts
// landing there are counted as ENOENT gaps, not hidden, and they are not liveness.
//
// LIVENESS IS PROVED, NOT INFERRED. A cycle counts only when Lstat immediately after the rename
// reports a symlink at the ancestor name - the same rule the Windows side of this family applies with
// reparseAt (openleaf_guard_322_windows_test.go). Without it, a zero damage numerator would be
// indistinguishable from a window nobody visited, which is the failure both tickets name in their ACs.
//
// THE THREE DAMAGE COUNTERS ARE KEPT APART. truncated is ticket 342's shape (the open followed the
// flipped ancestor and its O_TRUNC landed there); deleted is ticket 348's (the refusal-path unlink
// removed a name in somebody else's tree); narrowed is ticket 120's (a seal that wandered). One
// instrument's green does not discharge another ticket's - ticket 348 AC#3 verbatim.
package winsec

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	// ancAttempts342r1 is the census's and the verdict's own denominator, so magnitudes stay
	// comparable with 342-a1's and 342-v1's readings.
	ancAttempts342r1 = 1500
	// ancPace342r1 is how long each state stands per cycle: long enough for several production calls
	// to complete inside the link state, short enough that a red is a window the product walks rather
	// than one the fixture manufactured by standing still.
	ancPace342r1 = 200 * time.Microsecond
	// ancTargetRel and ancTargetAbs are the two link spellings, for the two reasons above.
	ancTargetRel = "rel"
	ancTargetAbs = "abs"

	ancLeaf342r1       = "leaf-342r1.bin"
	ancOutsideDir342r1 = "someone-elses-tree"
	// ancOwnBody342r1 and ancWitnessBody342r1 differ in LENGTH on purpose: an O_TRUNC that lands and
	// writes nothing leaves a short body, and equal-length bodies would let that read as "untouched".
	ancOwnBody342r1     = "the tree's own object, under the caller's name"
	ancWitnessBody342r1 = "foreign bytes are not yours to truncate, and this body is longer on purpose"
	ancPayload342r1     = "raced private payload 342r1 through a swapped ancestor"
	// ancChainELOOP342r1 is the kernel's wording, ambiguous between legs - which is why the attribution
	// below keeps 322-r1's "open of the name, not a deeper openat" exclusion.
	ancChainELOOP342r1 = "too many levels of symbolic links"
)

// The three answers form 乙's pinned walk can produce. Which of them named the CALLER'S LEAF decides
// where on the route the refusal came from: the parent-side walk (pinnedParent, BEFORE the open)
// cannot mention the leaf because it never looks at it; the seal-side walk (applyDescriptor on the
// name, AFTER the open) always does. That distinction is what makes the paired readings sayable
// instead of just quieter, so it is kept in the attribution rather than flattened into one counter.
const (
	whyPinnedLink342r1   = "the name carries a link, so the pinned object"
	whyPinnedSwap342r1   = "the name carries a different object than the one that was pinned"
	whyPinnedEscape342r1 = "path escapes from parent"
)

// ancLegs342r1 lists the answers this route can produce, in a declared order so two runs diff line by
// line.
var ancLegs342r1 = []struct{ why, name string }{
	{whyPinnedLink342r1, "pinned walk, sameObjectPinned link face (winsec_other.go through pinPlacement)"},
	{whyPinnedSwap342r1, "pinned walk, sameObjectPinned SameFile face"},
	{whyPinnedEscape342r1, "pinned walk, Root's own escape refusal (absolute-target link)"},
	{"reaches it through the link at", "placement floor before the act (platformVerifyPlacement)"},
	{"the name carries a link, so this descriptor", "sealHandle's inline link face (winsec_other.go)"},
	{"the bytes would not land in what was just sealed", "sealHandle's inline SameFile face"},
	{ancChainELOOP342r1, "ticket 322's leaf leg: the open of a name carrying a link"},
	{"no such file or directory", "the racer's own clear-then-rename gap (ENOENT): absent, not followed"},
	{"is a directory", "the open landed on the flipped ancestor itself"},
	{"not provably resolved", "C26 resolveString refused the spelling"},
}

const (
	ancOtherLeg342r1 = "some other failure"
	ancNilLeg342r1   = "none: the call returned nil"
)

// spokeAnc342r1 attributes one answer to a leg. ELOOP is credited to ticket 322's leaf leg only when
// the OPEN OF THE NAME reported it: the pinned walk also produces chain ELOOP, and counting those as
// the leaf leg would fake evidence for a leg that did not answer (322-r1's exclusion, same shape).
func spokeAnc342r1(err error) string {
	if err == nil {
		return ancNilLeg342r1
	}
	msg := err.Error()
	if strings.Contains(msg, ancChainELOOP342r1) && strings.Contains(msg, "open ") && !strings.Contains(msg, "openat ") {
		return ancLegs342r1[6].name
	}
	for _, leg := range ancLegs342r1 {
		if leg.why == ancChainELOOP342r1 {
			continue
		}
		if strings.Contains(msg, leg.why) {
			return leg.name
		}
	}
	return ancOtherLeg342r1
}

// ancIsPinnedWalk342r1 reports whether an answer came from the package's own pinning walk (either
// face of sameObjectPinned, or Root's escape face). Both the parent-side call (form 乙's new answer)
// and the seal-side call (already on the route) format through these, so the leaf test below tells
// them apart.
func ancIsPinnedWalk342r1(msg string) bool {
	return strings.Contains(msg, whyPinnedLink342r1) ||
		strings.Contains(msg, whyPinnedSwap342r1) ||
		strings.Contains(msg, whyPinnedEscape342r1)
}

// ancStats342r1 is one run's whole reading. Every field is printed; each test file gates its own.
type ancStats342r1 struct {
	attempts         int
	target           string
	cycles           int
	provedLinkCycles int32 // the instrument's own liveness floor: states proved to be links
	truncated        int   // ticket 342's numerator
	deleted          int   // ticket 348's numerator
	narrowed         int   // ticket 120's numerator
	untouched        int
	gaps             int // attempts that answered ENOENT (the racer's gap, counted not hidden)
	liveness         int // attempts the pinned walk refused
	pinnedParentLeg  int // ...BEFORE the open: the text names the parent, not the leaf (form 乙's answer)
	sealSideLeg      int // ...AFTER the open: the text names the caller's leaf (the leg 120 booked)
	legs             map[string]int
	firstProof       string
	firstGap         string
	firstFlip        string
	flipErrs         int
}

// ancWorld342r1 holds both trees, the staged names and the racer's counters.
type ancWorld342r1 struct {
	t          *testing.T
	target     string
	base       string
	foreign    string
	witness    string // <foreign>/<leaf>: somebody else's object, same leaf name
	data       string // <base>/data
	anc        string // <data>/sub: the ancestor that gets flipped
	name       string // <anc>/<leaf>: the spelling PrivateFile is handed
	stageReal  string // a prepared real directory carrying the caller's leaf, renamed into anc
	stageLink  string // a prepared symlink at the outside tree, renamed into anc
	cycles     int32
	linkCycles int32
	flipErrs   int32
	stop       chan struct{}
	done       chan struct{}
	stopOnce   sync.Once
	flipMu     sync.Mutex
	firstFlip  string
}

func newAncWorld342r1(t *testing.T, target string) *ancWorld342r1 {
	t.Helper()
	if target != ancTargetRel && target != ancTargetAbs {
		t.Fatalf("unknown ancestor link form %q", target)
	}
	w := &ancWorld342r1{t: t, target: target, base: SealableTempDirForTest124(t)}
	w.data = filepath.Join(w.base, "data")
	// The foreign tree is a SIBLING of the flipped ancestor's own directory, inside `data`. That is
	// what makes the link spelling that reaches it a BARE name: a symlink at `data/sub` resolves its
	// target relative to `data/`, so `someone-elses-tree` lands on `data/someone-elses-tree` with no
	// ".." element anywhere - the one spelling Root follows (see the header).
	w.foreign = filepath.Join(w.data, ancOutsideDir342r1)
	w.witness = filepath.Join(w.foreign, ancLeaf342r1)
	w.anc = filepath.Join(w.data, "sub")
	w.name = filepath.Join(w.anc, ancLeaf342r1)
	w.stageReal = filepath.Join(w.base, "stage-real")
	w.stageLink = filepath.Join(w.base, "stage-link")
	for _, d := range []string{w.foreign, w.data, w.stageReal} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatalf("stage the world (%s): %v", d, err)
		}
	}
	if err := os.WriteFile(filepath.Join(w.stageReal, ancLeaf342r1), []byte(ancOwnBody342r1), 0o600); err != nil {
		t.Fatalf("stage the real directory: %v", err)
	}
	if err := os.Rename(w.stageReal, w.anc); err != nil {
		t.Fatalf("give the caller a real ancestor: %v", err)
	}
	// rename takes the staged directory AWAY, so it has to be re-armed here exactly as toReal re-arms
	// it after each cycle: without this the very first settle back to the real state fails with ENOENT
	// and the whole run races an ancestor that is a link or absent. (本腿自抓的台损，⛔ 计入判据。)
	if err := os.Mkdir(w.stageReal, 0o700); err != nil {
		t.Fatalf("re-arm the staged directory: %v", err)
	}
	// The foreign witness: wide, with its own bytes, named by nobody in the caller's tree. Create then
	// chmod - os.WriteFile masks the mode with ^umask and the container's umask is 0022, so a 0666
	// argument here silently lands 0644 (the same bench bug that made "narrowed" fire on every attempt
	// of this leg's first pair run; ⛔ 计入判据).
	if err := os.WriteFile(w.witness, []byte(ancWitnessBody342r1), 0o600); err != nil {
		t.Fatalf("stage the witness: %v", err)
	}
	if err := os.Chmod(w.witness, 0o666); err != nil {
		t.Fatalf("widen the witness: %v", err)
	}
	if err := w.stageLinkSpelling(); err != nil {
		t.Fatalf("stage the link: %v", err)
	}
	info, err := os.Lstat(w.anc)
	if err != nil || info.Mode()&os.ModeSymlink != 0 {
		t.Fatalf("the ancestor did not stage as a real directory: %v mode=%v", err, info.Mode())
	}
	return w
}

// linkTarget is the spelling the staged symlink carries for this run's form.
func (w *ancWorld342r1) linkTarget() string {
	if w.target == ancTargetAbs {
		return w.foreign
	}
	// A bare sibling name: the link's target resolves against the directory the link lives in, which
	// is `data`, the same directory that holds the foreign tree. No ".." anywhere - with a ".." element
	// Root's own escape rule answers first and the identity gate under test never runs (measured on
	// this leg's first container run: every refusal in the rel form came back "path escapes from
	// parent", i.e. the instrument was driving Root's containment rule, not sameObjectPinned).
	return ancOutsideDir342r1
}

// stageLinkSpelling arms the staged symlink for the next cycle.
func (w *ancWorld342r1) stageLinkSpelling() error {
	if err := os.Remove(w.stageLink); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return os.Symlink(w.linkTarget(), w.stageLink)
}

// ancClear removes whatever stands at the ancestor name WITHOUT following it: a symlink is unlinked,
// a directory is emptied and then removed. Following here would reproduce the bug inside the fixture.
func ancClear(name string) error {
	info, err := os.Lstat(name)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return os.Remove(name)
	}
	if info.IsDir() {
		_ = os.Remove(filepath.Join(name, ancLeaf342r1))
		return os.Remove(name)
	}
	return os.Remove(name)
}

// toReal lands the ordinary state and re-arms the staged directory for the next cycle. The staged
// directory carries the caller's own leaf, so the real state is a real file at the spelling - the
// object a by-name O_TRUNC would have taken - not an absent name.
func (w *ancWorld342r1) toReal() error {
	if err := ancClear(w.anc); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(w.stageReal, ancLeaf342r1), []byte(ancOwnBody342r1), 0o600); err != nil {
		return err
	}
	if err := os.Rename(w.stageReal, w.anc); err != nil {
		return err
	}
	return os.Mkdir(w.stageReal, 0o700)
}

// toLink lands the swapped state and PROVES it: a cycle counts only when Lstat immediately after the
// rename reports a symlink at the ancestor name. This is the liveness key both tickets demand.
func (w *ancWorld342r1) toLink() error {
	if err := ancClear(w.anc); err != nil {
		return err
	}
	if err := os.Rename(w.stageLink, w.anc); err != nil {
		return err
	}
	info, err := os.Lstat(w.anc)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("ancestor %s is not a symlink after the flip, so the state was never entered", w.anc)
	}
	atomic.AddInt32(&w.linkCycles, 1)
	return w.stageLinkSpelling()
}

// startRacer alternates the two states until stopped. A plain goroutine, not an owner-managed one:
// d22scan's bare-goroutine ban is about production code, and this package's instrument family is the
// precedent - the racers in openleaf_322_other_test.go and posix_seal_window_120_other_test.go are
// declared the same way.
func (w *ancWorld342r1) startRacer() {
	w.stop = make(chan struct{})
	w.done = make(chan struct{})
	go func() {
		defer close(w.done)
		defer func() {
			if p := recover(); p != nil {
				w.noteFlip(fmt.Errorf("racer panic: %v", p))
			}
		}()
		for {
			select {
			case <-w.stop:
				return
			default:
			}
			if err := w.toLink(); err != nil {
				w.noteFlip(err)
				continue
			}
			atomic.AddInt32(&w.cycles, 1)
			time.Sleep(ancPace342r1)
			select {
			case <-w.stop:
				return
			default:
			}
			if err := w.toReal(); err != nil {
				w.noteFlip(err)
				continue
			}
			time.Sleep(ancPace342r1)
		}
	}()
}

func (w *ancWorld342r1) noteFlip(err error) {
	atomic.AddInt32(&w.flipErrs, 1)
	w.flipMu.Lock()
	defer w.flipMu.Unlock()
	if w.firstFlip == "" {
		w.firstFlip = err.Error()
	}
}

// stopRacer settles the world in the REAL state before any reading below it is taken, so a post-run
// check is about the witness and not about which half-cycle the racer happened to be in.
func (w *ancWorld342r1) stopRacer() {
	w.stopOnce.Do(func() { close(w.stop) })
	<-w.done
	if err := w.toReal(); err != nil {
		w.t.Fatalf("settle the world back to the real ancestor: %v", err)
	}
}

// run drives the production entry point over the flipping ancestor and samples the foreign witness
// after every attempt. Both tickets call PrivateFile, because both damages live on that one route.
func (w *ancWorld342r1) run(attempts int, call func() error) ancStats342r1 {
	w.startRacer()
	st := ancStats342r1{attempts: attempts, target: w.target, legs: map[string]int{}}
	for i := 0; i < attempts; i++ {
		err := call()
		leg := spokeAnc342r1(err)
		st.legs[leg]++
		if err != nil {
			msg := err.Error()
			switch leg {
			case ancLegs342r1[7].name:
				st.gaps++
				if st.firstGap == "" {
					st.firstGap = msg
				}
			}
			// Where the pinned walk answered decides what this attempt PROVES. A parent-side answer
			// (the text names the parent only) means the route stopped BEFORE opening anything: form
			// 乙's new leg, and the only one that can be claimed as "the open never followed". A
			// seal-side answer carries the caller's leaf name, which means an open already happened.
			if ancIsPinnedWalk342r1(msg) {
				st.liveness++
				if strings.Contains(msg, ancLeaf342r1) {
					st.sealSideLeg++
				} else {
					st.pinnedParentLeg++
					if st.firstProof == "" {
						st.firstProof = "refused by the pinned parent before the open: " + msg
					}
				}
			}
		}
		switch kind, detail := w.readWitness(); kind {
		case "deleted":
			st.deleted++
			if st.firstProof == "" {
				st.firstProof = fmt.Sprintf("witness ABSENT after attempt %d (%s)", i, detail)
			}
		case "truncated":
			st.truncated++
			if st.firstProof == "" {
				st.firstProof = fmt.Sprintf("witness body replaced after attempt %d (%s)", i, detail)
			}
		case "narrowed":
			st.narrowed++
		default:
			st.untouched++
		}
		w.restockWitness()
	}
	st.cycles = int(atomic.LoadInt32(&w.cycles))
	w.stopRacer()
	st.provedLinkCycles = atomic.LoadInt32(&w.linkCycles)
	st.flipErrs = int(atomic.LoadInt32(&w.flipErrs))
	w.flipMu.Lock()
	st.firstFlip = w.firstFlip
	w.flipMu.Unlock()
	return st
}

// restockWitness writes the foreign object back exactly as it was, widened to 0666 so a seal that
// wandered shows up as a mode change too (ticket 113's instrument).
//
// create-then-chmod, NOT WriteFile(…, 0666): os.WriteFile masks the mode argument with ^umask and the
// container runs with umask 0022, so a 0666 argument lands 0644 silently. That made "narrowed" fire on
// every attempt on BOTH trees - a bench artifact wearing a red's clothes (本腿自抓，⛔ 计入判据; the
// first pair run's narrowed=1423/1500 on the pre-landing tree was this, not a seal).
func (w *ancWorld342r1) restockWitness() {
	w.t.Helper()
	if err := os.WriteFile(w.witness, []byte(ancWitnessBody342r1), 0o600); err != nil {
		w.t.Fatalf("restock the witness: %v", err)
	}
	if err := os.Chmod(w.witness, 0o666); err != nil {
		w.t.Fatalf("widen the witness back to 0666: %v", err)
	}
}

// readWitness classifies the foreign object after one attempt. "deleted" is ticket 348's shape,
// "truncated" is ticket 342's, "narrowed" is ticket 120's.
func (w *ancWorld342r1) readWitness() (kind, detail string) {
	w.t.Helper()
	info, lErr := os.Lstat(w.witness)
	if lErr != nil {
		if errors.Is(lErr, fs.ErrNotExist) {
			return "deleted", lErr.Error()
		}
		return "untouched", "lstat: " + lErr.Error()
	}
	body, rErr := os.ReadFile(w.witness)
	if rErr != nil {
		return "untouched", "read: " + rErr.Error()
	}
	if string(body) != ancWitnessBody342r1 {
		return "truncated", fmt.Sprintf("size %d, content %q", len(body), body)
	}
	if info.Mode().Perm() != 0o666 {
		return "narrowed", info.Mode().String()
	}
	return "untouched", ""
}

// print writes the whole reading on one screen, legs in the declared order. Both test files call it,
// so the two tickets' evidence sits side by side without either borrowing the other's numerator.
func (s ancStats342r1) print(t *testing.T, tag string) {
	t.Helper()
	t.Logf("%s: target=%s attempts=%d cycles=%d provedLinkCycles=%d pinnedWalkRefusals=%d "+
		"pinnedParentBeforeOpen=%d refusedAfterOpen=%d truncated=%d deleted=%d narrowed=%d untouched=%d "+
		"enoentGaps=%d flipErrs=%d",
		tag, s.target, s.attempts, s.cycles, s.provedLinkCycles, s.liveness, s.pinnedParentLeg,
		s.sealSideLeg, s.truncated, s.deleted, s.narrowed, s.untouched, s.gaps, s.flipErrs)
	for _, leg := range ancLegs342r1 {
		if n := s.legs[leg.name]; n > 0 {
			t.Logf("%s': leg that answered = %s : %d attempts", tag, leg.name, n)
		}
	}
	for _, extra := range []string{ancOtherLeg342r1, ancNilLeg342r1} {
		if n := s.legs[extra]; n > 0 {
			t.Logf("%s': leg that answered = %s : %d attempts", tag, extra, n)
		}
	}
	if s.firstProof != "" {
		t.Logf("%s': first proof verbatim: %s", tag, s.firstProof)
	}
	if s.firstGap != "" {
		t.Logf("%s': first gap verbatim: %s", tag, s.firstGap)
	}
	if s.firstFlip != "" {
		t.Logf("%s': first racer error verbatim: %s", tag, s.firstFlip)
	}
}
