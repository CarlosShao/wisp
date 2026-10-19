//go:build windows

// Tickets 342 AC#1 and 348 AC#2/AC#3, leg 342-r1: the Windows world, racer and liveness proof for
// the two ANCESTOR-leg instruments. One world, THREE FORMS, two separate gates - each test file
// judges its own numerator, and nothing here is borrowed from the unix readings. Both tickets' ACs
// require per-platform readings because the ancestor shapes are different mechanisms: a link and a
// lexical containment rule on unix, a reparse point and OBJ_DONT_REPARSE here.
//
// WHY THREE FORMS AND NOT ONE. This is the坑 the dispatch named for the Windows half - "os.Root 钉
// 祖先⛔ 钉叶子" (342-a1's W-ROOTSYM row, followed=true, three runs) - and the three forms are what
// makes it measurable instead of merely asserted:
//
//   - winAncFormJunction: `mklink /J` at the ancestor, into somebody else's tree. A junction's stored
//     substitute name is ABSOLUTE, so the pinned chain's descent hits it with OBJ_DONT_REPARSE, Go
//     reads the link back, doInRoot restarts with an absolute spelling and refuses - "path escapes
//     from parent". That is an ACT-TIME refusal (342-a1 measured this exact row: WYI damage=0,
//     root-refuse 1312/1500), and it is form 乙's answer for the shape Windows users actually meet.
//   - winAncFormSymDir: a directory SYMLINK at the ancestor, with a RELATIVE target. That stays inside
//     the root the chain descends from, so containment does NOT refuse it and the walk would follow it
//   - which is why the chain re-verifies every component against the name it was reached through
//     (sameObjectPinnedWindows). This is the form that separates "pinned the ancestor" from "trusted
//     the ancestor". Needs developer mode: unavailable construction is a t.Skipf, never a quiet pass.
//   - winAncFormLeafSym: the ancestor stays REAL and the LEAF name is flipped between a real file and
//     a relative symlink at the foreign object. os.Root follows a leaf link inside the root (that IS
//     W-ROOTSYM), so the chain alone would hand back a handle on somebody else's object - this form
//     exists to prove the leaf identity re-check inside openPinnedLeaf runs and refuses, and to name
//     the residual it cannot undo: an O_TRUNC the followed open already carried out.
//
// LIVENESS, WINDOWS FACE (票 342 AC#1 names the key: reparseAt(ancestor) 复证). A cycle counts only
// when reparseAt reports TRUE immediately after the flip's rename. That is this package's own probe
// (reparse_windows_test.go:200), the same rule the 322-r1 Windows route instrument applies.
//
// WHY THE WITNESS IS READ BY ITS OWN ABSOLUTE SPELLING. The question is "what happened to somebody
// else's object". Once the ancestor or the leaf is flipped, the caller's spelling also names that
// object, so reading through it would report on whatever the attacker is currently pointing at.
//
// TWO DAMAGE COUNTERS, KEPT APART (票 348 AC#3: the deletion pin may not reuse the truncation pin's
// green): `truncated` = the outside object lost its body to the route's open; `deleted` = the
// refusal-path unlink removed it. Both are reachable on the pre-landing bytes, which is exactly why
// these are two tickets.
package winsec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	// winAncAttempts342r1 is lower than the unix denominator, and the reason is measured rather than
	// taste: each cycle here costs a cmd-spawned mklink (os.Symlink can only make a symlink, which
	// needs developer mode), so 322-r1's own Windows route instrument set the family precedent at 400
	// attempts. What has to reproduce across legs is the SHAPE (damage>0 before, 0 after, liveness in
	// the same reading); every run prints its own denominator.
	winAncAttempts342r1 = 400
	winAncPace342r1     = 2 * time.Millisecond
	winAncLeaf342r1     = "leaf-342r1w.bin"
	winAncOutside342r1  = "someone-elses-tree-w"
	// Different lengths, same rule as the unix world: an O_TRUNC that lands and writes nothing leaves a
	// short body, and equal-length bodies would let that read as "untouched".
	winAncOwnBody342r1 = "the tree's own object, under the caller's name"
	winAncWitness342r1 = "foreign bytes are not yours to truncate, and this body is longer on purpose"
	winAncPayload342r1 = "raced private payload 342r1 windows ancestor"
)

// The three flip forms.
const (
	winAncFormJunction = "junction"
	winAncFormSymDir   = "symlink-dir"
	winAncFormLeafSym  = "leaf-symlink"
)

// winAncLegs342r1 are the answers this route can produce on this platform, in a declared order so two
// runs diff line by line. The first three are form 乙's Windows face; the rest were already on the
// route before this landing and keep their own names.
var winAncLegs342r1 = []struct{ why, name string }{
	{"path escapes from parent", "pinned chain, act-time refusal of a reparse target (form B's answer: os.Root + OBJ_DONT_REPARSE)"},
	{"the name carries a link, so the pinned object", "pinned chain, sameObjectPinnedWindows (a relative ancestor the walk could have followed)"},
	{"the name carries a link, so this handle", "leaf identity re-check inside openPinnedLeaf (os.Root pins ancestors, NOT the leaf)"},
	{"this handle is not the object the name carries now", "leaf identity re-check, SameFile face"},
	{"the reparse-point probe could not read the name", "preOpenLeafGuard's unreadable-name refusal (ticket 322's Windows leaf leg)"},
	{"the name carries a reparse point", "preOpenLeafGuard's leaf refusal (ticket 322's Windows leaf leg)"},
	{"reaches it through the link at", "placement floor before the act (platformVerifyPlacement)"},
	{"path traverses a reparse point", "the installed C26 resolver's own refusal (resolveString -> risk), i.e. steady state BEFORE any act"},
	{"cannot apply a private security descriptor", "a seal refusal on this route"},
	{"not provably resolved", "C26 resolveString refused the spelling"},
}

const (
	winAncOther342r1 = "some other failure"
	winAncNil342r1   = "none: the call returned nil"
)

func spokeWinAnc342r1(err error) string {
	if err == nil {
		return winAncNil342r1
	}
	msg := err.Error()
	for _, leg := range winAncLegs342r1 {
		if strings.Contains(msg, leg.why) {
			return leg.name
		}
	}
	return winAncOther342r1
}

// winAncIsFormLeg342r1 reports whether an answer came from one of form 乙's own legs (the pinned chain
// or the leaf re-check), i.e. whether the refusal is attributable to this landing rather than to a leg
// that was already there.
func winAncIsFormLeg342r1(leg string) bool {
	return leg == winAncLegs342r1[0].name || leg == winAncLegs342r1[1].name ||
		leg == winAncLegs342r1[2].name || leg == winAncLegs342r1[3].name
}

// winAncStats342r1 is one run's reading. Every field is printed; each test file gates its own.
type winAncStats342r1 struct {
	attempts      int
	form          string
	cycles        int
	provedReparse int32 // reparseAt(ancestor or leaf) proved TRUE right after the flip's rename
	truncated     int   // ticket 342's numerator
	deleted       int   // ticket 348's numerator
	formRefusals  int   // attempts answered by one of form 乙's legs
	lands         int   // attempts that returned nil: the pre-landing shape, where the route "succeeds"
	legs          map[string]int
	firstProof    string
	firstFlip     string
	flipErrs      int
}

type winAncWorld342r1 struct {
	t         *testing.T
	form      string
	base      string
	data      string
	sub       string
	name      string
	foreign   string
	witness   string
	parkReal  string
	staged    string
	cycles    int32
	reparseOK int32
	flipErrs  int32
	stop      chan struct{}
	done      chan struct{}
	stopOnce  sync.Once
	flipMu    sync.Mutex
	firstFlip string
}

func newWinAncWorld342r1(t *testing.T, form string) *winAncWorld342r1 {
	t.Helper()
	if form != winAncFormJunction && form != winAncFormSymDir && form != winAncFormLeafSym {
		t.Fatalf("unknown Windows flip form %q", form)
	}
	w := &winAncWorld342r1{t: t, form: form, base: SealableTempDirForTest124(t)}
	w.data = filepath.Join(w.base, "data")
	w.sub = filepath.Join(w.data, "sub")
	w.name = filepath.Join(w.sub, winAncLeaf342r1)
	w.foreign = filepath.Join(w.base, winAncOutside342r1)
	w.witness = filepath.Join(w.foreign, winAncLeaf342r1)
	w.parkReal = filepath.Join(w.data, "parked-real-sub")
	// The staged name lives IN THE DIRECTORY THE FLIP LANDS IN, because a relative reparse target is
	// resolved against the directory that CONTAINS the link: staged at `<base>` and renamed into `data`
	// would leave the target pointing at something the link's own directory cannot reach. (The junction
	// form carries an ABSOLUTE substitute name, so it does not care - which is precisely why containment
	// refuses that one and not this one.)
	if w.form == winAncFormLeafSym {
		w.staged = filepath.Join(w.sub, "staged-leaf-flip")
	} else {
		w.staged = filepath.Join(w.data, "staged-anc-flip")
	}
	for _, d := range []string{w.data, w.sub, w.foreign} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatalf("stage %s: %v", d, err)
		}
	}
	if err := os.WriteFile(w.name, []byte(winAncOwnBody342r1), 0o600); err != nil {
		t.Fatalf("give the caller a name to stand on: %v", err)
	}
	if err := os.WriteFile(w.witness, []byte(winAncWitness342r1), 0o600); err != nil {
		t.Fatalf("stage the witness: %v", err)
	}
	// Privilege honesty in this package's own shape: a junction needs no privilege, so a failure here
	// is a real bench fault; a directory symlink needs developer mode, so THAT unavailable construction
	// is a t.Skipf that acceptance has to count (mkDirSymlink's rule), never a quiet pass.
	if err := w.stageFlip(); err != nil {
		if w.form == winAncFormSymDir {
			t.Skipf("cannot construct a directory symlink here: %v (needs developer mode or "+
				"SeCreateSymbolicLinkPrivilege; the junction form is the no-privilege one)", err)
		}
		t.Fatalf("stage the flip: %v", err)
	}
	t.Cleanup(w.settleAll)
	return w
}

// settleAll leaves nothing that a later os.RemoveAll could traverse: a junction standing inside a temp
// root whose target has already been deleted is exactly the shape A51② booked for RemoveAll, so every
// staged reparse point is unlinked (never followed) through this package's own entry point.
func (w *winAncWorld342r1) settleAll() {
	_ = removeWinAnc(w.staged)
	_ = removeWinAnc(w.name)
	if err := w.toReal(); err != nil {
		w.t.Logf("AC#1'/AC#2' cleanup note: settle at exit: %v", err)
	}
	_ = os.RemoveAll(w.foreign)
}

// removeWinAnc deletes an entry without following it, tolerating absence. For a reparse point that
// means unlinking the POINT, which is RemoveUnlinked's promise; os.Remove on a directory reparse point
// with a non-empty target is what A51② measured as impossible.
func removeWinAnc(name string) error {
	if _, err := os.Lstat(name); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return RemoveUnlinked(name)
}

// execNoCheck runs a helper and hands its output back as TEXT plus an error, so a failure can become a
// Skipf instead of mustExec's Fatalf (that helper is the family's, but it decides the run is over).
func execNoCheck(name string, args ...string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = &out, &out
	err := cmd.Run()
	return out.String(), err
}

// stageFlip arms the staged name for this form's next cycle.
func (w *winAncWorld342r1) stageFlip() error {
	if err := removeWinAnc(w.staged); err != nil {
		return err
	}
	switch w.form {
	case winAncFormJunction:
		out, err := execNoCheck("cmd", "/c", "mklink", "/J", w.staged, w.foreign)
		if err != nil {
			return fmt.Errorf("mklink /J: %v (%s)", err, strings.TrimSpace(out))
		}
	case winAncFormSymDir:
		// A RELATIVE target, resolved against the directory the link ends up in (`data`), so it lands on
		// the foreign tree INSIDE the volume root the chain descends from: containment does not refuse
		// it, so the answer has to come from the identity gate.
		if err := os.Symlink(filepath.Join("..", winAncOutside342r1), w.staged); err != nil {
			return err
		}
	case winAncFormLeafSym:
		// ★The target is the foreign FILE, not the foreign directory. The first version of this case
		// pointed at the directory, so the flipped leaf resolved to a directory and the open answered
		// EISDIR - the run never drove the shape its gate claims to test, and "no leg answered" was a
		// bench fault and NOT a reading about the landing (本腿自抓，⛔ 计入判据；that case's earlier
		// numbers must be re-measured before anyone cites them).
		// Same resolution rule as above: the link ends up inside `sub`, so the target walks up from
		// there to `base`, then into the foreign tree and onto the one name in it.
		if err := os.Symlink(filepath.Join("..", "..", winAncOutside342r1, winAncLeaf342r1), w.staged); err != nil {
			return err
		}
	}
	if !reparseAt(w.staged) {
		return fmt.Errorf("the staged name %s is not a reparse point, so this fixture races nothing", w.staged)
	}
	return nil
}

// toReal lands the ordinary state for THIS form. It is form-aware on purpose: the ancestor forms flip
// `sub` itself, while the leaf form keeps `sub` real and flips only the leaf name - so asking the leaf
// form to remove `sub` would fail ENOTEMPTY (a directory holding a link) and the first settle would
// end the run with a bench fault instead of a reading.
func (w *winAncWorld342r1) toReal() error {
	if w.form == winAncFormLeafSym {
		if err := removeWinAnc(w.name); err != nil {
			return err
		}
		if err := os.MkdirAll(w.sub, 0o700); err != nil {
			return err
		}
		return os.WriteFile(w.name, []byte(winAncOwnBody342r1), 0o600)
	}
	if err := removeWinAnc(w.sub); err != nil {
		return err
	}
	if _, err := os.Lstat(w.parkReal); err != nil {
		// Nothing parked (a first settle, or a previous settle that already rebuilt the real state):
		// re-create it rather than rename a ghost.
		if err := os.MkdirAll(w.sub, 0o700); err != nil {
			return err
		}
		return os.WriteFile(w.name, []byte(winAncOwnBody342r1), 0o600)
	}
	if err := os.Rename(w.parkReal, w.sub); err != nil {
		return err
	}
	return os.WriteFile(w.name, []byte(winAncOwnBody342r1), 0o600)
}

// toFlip lands the swapped state and PROVES it: the cycle counts only when reparseAt reports TRUE at
// the name the flip just moved in. That is the liveness key, and for the leaf form it is the leaf.
func (w *winAncWorld342r1) toFlip() error {
	switch w.form {
	case winAncFormLeafSym:
		if err := removeWinAnc(w.name); err != nil {
			return err
		}
		if err := os.Rename(w.staged, w.name); err != nil {
			return err
		}
		if !reparseAt(w.name) {
			return fmt.Errorf("leaf %s is not a reparse point after the flip, so the state was never entered", w.name)
		}
	default:
		if err := os.Rename(w.sub, w.parkReal); err != nil {
			return err
		}
		if err := os.Rename(w.staged, w.sub); err != nil {
			return err
		}
		if !reparseAt(w.sub) {
			return fmt.Errorf("ancestor %s is not a reparse point after the flip, so the state was never entered", w.sub)
		}
	}
	atomic.AddInt32(&w.reparseOK, 1)
	return w.stageFlip()
}

func (w *winAncWorld342r1) startRacer() {
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
			if err := w.toFlip(); err != nil {
				w.noteFlip(err)
				continue
			}
			atomic.AddInt32(&w.cycles, 1)
			time.Sleep(winAncPace342r1)
			select {
			case <-w.stop:
				return
			default:
			}
			if err := w.toReal(); err != nil {
				w.noteFlip(err)
				continue
			}
			time.Sleep(winAncPace342r1)
		}
	}()
}

func (w *winAncWorld342r1) noteFlip(err error) {
	atomic.AddInt32(&w.flipErrs, 1)
	w.flipMu.Lock()
	defer w.flipMu.Unlock()
	if w.firstFlip == "" {
		w.firstFlip = err.Error()
	}
}

func (w *winAncWorld342r1) stopRacer() {
	w.stopOnce.Do(func() { close(w.stop) })
	<-w.done
	if err := w.toReal(); err != nil {
		w.t.Fatalf("settle the world back to the real state: %v", err)
	}
}

func (w *winAncWorld342r1) run(attempts int, call func() error) winAncStats342r1 {
	w.startRacer()
	st := winAncStats342r1{attempts: attempts, form: w.form, legs: map[string]int{}}
	for i := 0; i < attempts; i++ {
		err := call()
		leg := spokeWinAnc342r1(err)
		st.legs[leg]++
		switch {
		case winAncIsFormLeg342r1(leg):
			st.formRefusals++
			if st.firstProof == "" {
				st.firstProof = "answered by the landing's own leg: " + err.Error()
			}
		case leg == winAncNil342r1:
			st.lands++
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
		}
		w.restockWitness()
	}
	st.cycles = int(atomic.LoadInt32(&w.cycles))
	w.stopRacer()
	st.provedReparse = atomic.LoadInt32(&w.reparseOK)
	st.flipErrs = int(atomic.LoadInt32(&w.flipErrs))
	w.flipMu.Lock()
	st.firstFlip = w.firstFlip
	w.flipMu.Unlock()
	return st
}

func (w *winAncWorld342r1) restockWitness() {
	w.t.Helper()
	if err := os.WriteFile(w.witness, []byte(winAncWitness342r1), 0o600); err != nil {
		w.t.Fatalf("restock the witness: %v", err)
	}
}

func (w *winAncWorld342r1) readWitness() (kind, detail string) {
	w.t.Helper()
	info, lErr := os.Lstat(w.witness)
	if lErr != nil {
		if os.IsNotExist(lErr) {
			return "deleted", lErr.Error()
		}
		return "untouched", "lstat: " + lErr.Error()
	}
	body, rErr := os.ReadFile(w.witness)
	if rErr != nil {
		return "untouched", "read: " + rErr.Error()
	}
	if string(body) != winAncWitness342r1 {
		return "truncated", fmt.Sprintf("size %d of %d, content %q", len(body), info.Size(), body)
	}
	return "untouched", ""
}

func (s winAncStats342r1) print(t *testing.T, tag string) {
	t.Helper()
	t.Logf("%s: form=%s attempts=%d cycles=%d provedReparse=%d formBLegRefusals=%d truncated=%d deleted=%d "+
		"returnedNil=%d flipErrs=%d",
		tag, s.form, s.attempts, s.cycles, s.provedReparse, s.formRefusals, s.truncated, s.deleted,
		s.lands, s.flipErrs)
	for _, leg := range winAncLegs342r1 {
		if n := s.legs[leg.name]; n > 0 {
			t.Logf("%s': leg that answered = %s : %d attempts", tag, leg.name, n)
		}
	}
	for _, extra := range []string{winAncOther342r1, winAncNil342r1} {
		if n := s.legs[extra]; n > 0 {
			t.Logf("%s': leg that answered = %s : %d attempts", tag, extra, n)
		}
	}
	if s.firstProof != "" {
		t.Logf("%s': first proof verbatim: %s", tag, s.firstProof)
	}
	if s.firstFlip != "" {
		t.Logf("%s': first racer error verbatim: %s", tag, s.firstFlip)
	}
}
