//go:build windows

package winsec

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

// This file is ticket 106's instrument. The defect it pins is that the private
// set used to be judged by the *spelling* a descriptor happens to carry, and
// Windows spells a principal by its account name whenever it has one: the CI
// runner's own seal wrote a grant for the token user, the SDDL renderer printed
// that user's SID as the well-known token "LA", "LA" was not in the hand-listed
// name table, and so a descriptor this package built itself was refused by this
// package's own check. Nothing leaked; the store simply could not open.
//
// The instrument therefore never asserts a fixed verdict for a name. It plants
// the runner's descriptor verbatim, asks this machine what the name resolves to,
// and requires the gate to agree with the *binary* DACL's resolved SIDs - which
// is a judgment that runs identically on a machine where "LA" is the current
// user (the runner) and on one where it is somebody else (this box).

const (
	// runnerDescriptor is the DACL of the directory ticket 106 quotes from run
	// 35586044995, verbatim: our own six-ACE shape, third trustee spelled "LA".
	runnerDescriptor = "D:PAI(A;;FA;;;SY)(A;OICIIO;GA;;;SY)(A;;FA;;;BA)(A;OICIIO;GA;;;BA)" +
		"(A;;FA;;;LA)(A;OICIIO;GA;;;LA)"
	// everyoneSID / administratorsSID / systemSID are the resolved forms of the
	// well-known tokens WD / BA / SY.
	systemSID      = "S-1-5-18"
	administrators = "S-1-5-32-544"
	everyoneSID    = "S-1-1-0"
)

// plantDescriptor writes an SDDL descriptor onto a real directory through the
// OS's own parser, so what lands on disk is the OS's own choice of spelling and
// not the test's.
func plantDescriptor(t *testing.T, path, sddl string) {
	t.Helper()
	var sd *windows.SECURITY_DESCRIPTOR
	p, err := windows.UTF16PtrFromString(sddl)
	if err != nil {
		t.Fatal(err)
	}
	// ConvertStringSecurityDescriptorToSecurityDescriptorW is not exported by
	// x/sys, and it is the only way to ask the OS for its own reading of a name.
	r1, _, e1 := windows.NewLazySystemDLL("advapi32.dll").
		NewProc("ConvertStringSecurityDescriptorToSecurityDescriptorW").
		Call(uintptr(unsafe.Pointer(p)), 1, uintptr(unsafe.Pointer(&sd)), 0)
	if r1 == 0 {
		t.Fatalf("plant %q: %v", sddl, e1)
	}
	acl, _, err := sd.DACL()
	if err != nil {
		t.Fatalf("plant %q: no DACL: %v", sddl, err)
	}
	// PROTECTED, because the descriptor in the CI log is protected: this plants a
	// state to inspect, it does not seal anything, and an unprotected object would
	// be refused for a reason that has nothing to do with the private set.
	if err := windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION,
		nil, nil, acl, nil); err != nil {
		t.Fatalf("plant %q on %s: %v", sddl, path, err)
	}
	if _, err := windows.LocalFree(windows.Handle(unsafe.Pointer(sd))); err != nil {
		t.Fatalf("plant %q: free descriptor: %v", sddl, err)
	}
}

// resolvedTrustees reads the object's DACL back in its binary form and returns
// one entry per ACE: the trustee's SID string, the AceType, and whether the ACE
// is flagged inherited or inherit-only. This is the independent oracle the gate
// is required to agree with; it consults no names at all.
func resolvedTrustees(t *testing.T, path string) []string {
	t.Helper()
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		t.Fatalf("read DACL of %s: %v", path, err)
	}
	acl, _, err := sd.DACL()
	if err != nil {
		t.Fatalf("DACL of %s: %v", path, err)
	}
	out := make([]string, 0, int(acl.AceCount))
	for i := uint32(0); i < uint32(acl.AceCount); i++ {
		var ace *windows.ACCESS_ALLOWED_ACE
		if err := windows.GetAce(acl, i, &ace); err != nil {
			t.Fatalf("GetAce(%s, %d): %v", path, i, err)
		}
		trustee := "unreadable"
		if ace.Header.AceType == windows.ACCESS_ALLOWED_ACE_TYPE ||
			ace.Header.AceType == windows.ACCESS_DENIED_ACE_TYPE {
			trustee = (*windows.SID)(unsafe.Pointer(&ace.SidStart)).String()
		}
		out = append(out, fmt.Sprintf("%d/%#x=%s", ace.Header.AceType, ace.Header.AceFlags, trustee))
	}
	// No LocalFree here: x/sys' GetNamedSecurityInfo already released the win32
	// allocation and handed back a copy on the Go heap.
	return out
}

// privateSetSIDsBySID is the private set computed the way the gate must compute
// it: SIDs only, no name spellings, resolved from the account this process runs
// as.
func inPrivateSet(set map[string]bool, sid string) bool { return set[sid] }

func privateSetSIDsBySID(t *testing.T) map[string]bool {
	t.Helper()
	u, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	return map[string]bool{systemSID: true, administrators: true, u.Uid: true}
}

// TestGateJudgesThePrivateSetByResolvedSID is AC#2 and AC#3 both halves at once:
// the same judgment on a machine where "LA" is the current user and on one where
// it is not.
func TestGateJudgesThePrivateSetByResolvedSID(t *testing.T) {
	set := privateSetSIDsBySID(t)
	la := administratorSID(t)
	named := func(sid string) bool { return inPrivateSet(set, sid) }

	t.Run("runner descriptor is judged by what it resolves to", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "runner")
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		plantDescriptor(t, dir, runnerDescriptor)
		// The OS is free to re-render the ACEs; what must not change is the
		// verdict, which follows the SIDs.
		resolved := resolvedTrustees(t, dir)
		wantInSet := true
		for _, ace := range resolved {
			i := strings.IndexByte(ace, '=')
			if i < 0 {
				t.Fatalf("unparsable oracle entry %q", ace)
			}
			if !inPrivateSet(set, ace[i+1:]) {
				wantInSet = false
			}
		}
		err := verifyPrivate(dir)
		if wantInSet && err != nil {
			t.Fatalf("gate refused a descriptor whose every trustee resolves inside the private set: %v (trustees %v, set %v, LA=%s is current user: %v)",
				err, resolved, setKeys(set), la, named(la))
		}
		if !wantInSet && err == nil {
			t.Fatalf("gate accepted a descriptor naming %v, which is outside the private set %v", resolved, setKeys(set))
		}
	})

	t.Run("the set holds no name in the form it is compared with", func(t *testing.T) {
		// The judgment side of the rule: a spelling is never an input to "is this
		// inside the private set", so nothing that enters the set can be a name.
		// "SY", "BA", "LA" and "ME" are all renderings the OS picks, and "ME" in
		// particular round-trips as a *placeholder* SID that is not the current
		// user - the previous whitelist carried it as if it were.
		set, me, err := privateSet()
		if err != nil {
			t.Fatal(err)
		}
		for entry := range set {
			if !strings.HasPrefix(entry, "S-") {
				t.Fatalf("the private set admits a principal by name, not by SID: %q (set %v)", entry, setKeys(set))
			}
		}
		for _, name := range []string{"SY", "BA", "LA", "ME", "WD", "AU"} {
			if set[name] {
				t.Fatalf("the private still names %q", name)
			}
			if name == "ME" && me == name {
				t.Fatalf("the current user was resolved to the name %q", me)
			}
		}
	})

	t.Run("a name spelling cannot change the verdict", func(t *testing.T) {
		// The runner's descriptor and the numerically identical one describe the
		// same three trustees whenever "LA" resolves to one of them; the gate
		// must not disagree between the two spellings.
		numeric := fmt.Sprintf("D:PAI(A;;FA;;;%s)(A;OICIIO;GA;;;%s)(A;;FA;;;%s)(A;OICIIO;GA;;;%s)(A;;FA;;;%s)(A;OICIIO;GA;;;%s)",
			systemSID, systemSID, administrators, administrators, la, la)
		base := t.TempDir()
		namedDir, numericDir := filepath.Join(base, "named"), filepath.Join(base, "numeric")
		for _, d := range []string{namedDir, numericDir} {
			if err := os.Mkdir(d, 0o700); err != nil {
				t.Fatal(err)
			}
		}
		plantDescriptor(t, namedDir, runnerDescriptor)
		plantDescriptor(t, numericDir, numeric)
		namedErr, numericErr := verifyPrivate(namedDir), verifyPrivate(numericDir)
		if (namedErr == nil) != (numericErr == nil) {
			t.Fatalf("the same three trustees decided differently by spelling alone: named=%v numeric=%v", namedErr, numericErr)
		}
		if a, b := resolvedTrustees(t, namedDir), resolvedTrustees(t, numericDir); strings.Join(a, " ") != strings.Join(b, " ") {
			t.Fatalf("the two spellings did not land on the same DACL: %v vs %v", a, b)
		}
	})
}

// administratorSID is the local account the SDDL renderer spells "LA", taken
// from the system's own well-known-SID table rather than from a name lookup, so
// nothing in these tests depends on how an account happens to be called. On the
// CI runner that account is the one the job runs as, which is the whole defect;
// on an ordinary box it is somebody else. The tests below ask this machine
// instead of assuming either answer.
func administratorSID(t *testing.T) string {
	t.Helper()
	return plantedTrusteeSID(t, "D:(A;;FA;;;LA)")
}

// guestsSID is BUILTIN\Guests (S-1-5-32-546), resolved through the OS's own SDDL
// parser from its numeric form and read back out of the binary ACE. It is the
// principal ticket 112 added to the fixture: an alias can never be the token
// user's own SID, so it is a stranger on a machine that runs the job as the built-in
// Administrator and on one that does not, which is what makes the "cleared and
// named by SID" leg machine-independent. The OS still renders it back by name
// ("BG"), so it exercises the same spelling-vs-SID distinction "LA" did.
func guestsSID(t *testing.T) string {
	t.Helper()
	guests := plantedTrusteeSID(t, "D:(A;;FA;;;S-1-5-32-546)")
	if guests != "S-1-5-32-546" {
		t.Fatalf("asking this machine for BUILTIN\\Guests returned %q, not the alias SID", guests)
	}
	return guests
}

func plantedTrusteeSID(t *testing.T, sddl string) string {
	t.Helper()
	scratch := filepath.Join(t.TempDir(), "la")
	if err := os.Mkdir(scratch, 0o700); err != nil {
		t.Fatal(err)
	}
	plantDescriptor(t, scratch, sddl)
	aces := resolvedTrustees(t, scratch)
	if len(aces) != 1 {
		t.Fatalf("resolving %q left %d ACEs: %v", sddl, len(aces), aces)
	}
	i := strings.IndexByte(aces[0], '=')
	return aces[0][i+1:]
}

// standsOn is the binary DACL read of "this principal is on this object", with no
// name and no rendering involved.
func standsOn(t *testing.T, path, sid string) bool {
	t.Helper()
	for _, ace := range resolvedTrustees(t, path) {
		if i := strings.IndexByte(ace, '='); i >= 0 && ace[i+1:] == sid {
			return true
		}
	}
	return false
}

// TestSealNarrowsAndNamesThePrincipalItRemovedBySID is AC#3 leg (a): a real
// grant left for another account is still not tolerated - the seal narrows it
// and says so, and what it says is the one form that cannot be confused with
// somebody else's account.
func TestSealNarrowsAndNamesThePrincipalItRemovedBySID(t *testing.T) {
	var got []narrowNotice
	orig := noticeNarrowed
	noticeNarrowed = func(n narrowNotice) { got = append(got, n) }
	t.Cleanup(func() { noticeNarrowed = orig })

	set := privateSetSIDsBySID(t)
	root := filepath.Join(t.TempDir(), "store")
	if err := PrivateDirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	la := administratorSID(t)
	if la == "" {
		t.Fatal("no local Administrator account to plant")
	}
	// Ticket 112's red on run 35595651898 was this fixture's assumption, not the
	// notice. The Windows job runs as the built-in Administrator, so "LA" resolves
	// to the token user's SID, which ticket 106 put inside the private set by
	// definition - the seal therefore kept that grant and named only Everyone. The
	// CI's own readings say the instrument was wrong, not the machine:
	//
	//	private_set_sid_windows_test.go:268: ... cleared="S-1-1-0(A;OICI;FA;;;WD)",
	//	  want S-1-5-21-3699639565-2515463329-295617607-500 in it (root DACL now
	//	  [0/0x0=S-1-5-18 0/0xb=S-1-5-18 0/0x0=S-1-5-32-544 0/0xb=S-1-5-32-544
	//	   0/0x0=S-1-5-21-3699639565-2515463329-295617607-500
	//	   0/0xb=S-1-5-21-3699639565-2515463329-295617607-500])
	//
	// So the candidate list is now read off this machine instead of hard-coded: a
	// stranger that can never be the token user (Guests) carries the proof
	// everywhere, and a private-set member is planted on purpose so the membership
	// direction - kept, and never reported as cleared - is asserted too. Which
	// bucket "LA" lands in is measured, named in the failure message, and never
	// skipped.
	guests := guestsSID(t)
	if inPrivateSet(set, guests) {
		t.Fatalf("no stranger to plant on this machine: BUILTIN\\Guests is inside the private set %v", setKeys(set))
	}
	foreign := []string{everyoneSID, guests}
	member := administrators
	if inPrivateSet(set, la) {
		member = la // the runner's shape: the job runs as the built-in Administrator
	} else {
		foreign = append(foreign, la) // this box's shape: LA is somebody else
	}
	t.Logf("planted as strangers: %v; planted as a private-set member: %s (LA=%s, in set: %v, set=%v)",
		foreign, member, la, inPrivateSet(set, la), setKeys(set))

	plant := append(append([]string{}, foreign...), member)
	for _, sid := range plant {
		mustExec(t, "icacls", root, "/grant", "*"+sid+":(OI)(CI)F")
	}
	for _, sid := range plant {
		if !standsOn(t, root, sid) {
			t.Fatalf("the fixture planted nothing to judge for %s: %v", sid, resolvedTrustees(t, root))
		}
	}

	got = nil
	if err := SealDir(root); err != nil {
		t.Fatalf("sealing a tree an operator widened: %v", err)
	}
	var cleared []string
	for _, n := range got {
		if noticeNamesTree(n, root) {
			cleared = append(cleared, n.Principals...)
		}
	}
	joined := strings.Join(cleared, " ")
	for _, want := range foreign {
		if !strings.Contains(joined, want) {
			t.Fatalf("the notice named the cleared principal by spelling only, not by the resolved SID it holds: cleared=%q, want %s in it (root DACL now %v)", joined, want, resolvedTrustees(t, root))
		}
	}
	// The membership direction, which is the one the CI runner was actually
	// exercising when this test went red: a principal inside the private set is
	// kept and is never reported as narrowed away, whichever name the OS chose to
	// render it with.
	if strings.Contains(joined, member) {
		t.Fatalf("the notice named a private-set member as cleared: member=%s cleared=%q (set %v, trustees %v)",
			member, joined, setKeys(set), resolvedTrustees(t, root))
	}
	if !standsOn(t, root, member) {
		t.Fatalf("the seal removed the private-set member %s it was given: %v", member, resolvedTrustees(t, root))
	}
	// And the narrowing itself: nothing outside the private set may stand on the
	// object after the seal, whichever way the OS spells what is left.
	for _, ace := range resolvedTrustees(t, root) {
		i := strings.IndexByte(ace, '=')
		if !inPrivateSet(set, ace[i+1:]) {
			t.Fatalf("seal left a foreign grant on %s: %s (%v)", root, ace, ace)
		}
	}
	if err := verifyPrivate(root); err != nil {
		t.Fatalf("the root is not private after the seal: %v", err)
	}
}

// TestGateRefusesADescriptorThatLeavesARealGrantToAnotherAccount is AC#3 leg (a)
// on the gate itself: a descriptor that hands a stranger a real authorization
// next to ours is not private, whichever spelling the OS chose for the three of
// us. The control leg (the same descriptor without the stranger) keeps this from
// passing for an unrelated reason.
func TestGateRefusesADescriptorThatLeavesARealGrantToAnotherAccount(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	private := fmt.Sprintf("D:PAI(A;;FA;;;%s)(A;;FA;;;%s)(A;;FA;;;%s)", systemSID, administrators, u.Uid)
	widened := private + "(A;;FA;;;" + everyoneSID + ")"

	base := t.TempDir()
	control, widenedDir := filepath.Join(base, "control"), filepath.Join(base, "widened")
	for _, d := range []string{control, widenedDir} {
		if err := os.Mkdir(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	plantDescriptor(t, control, private)
	plantDescriptor(t, widenedDir, widened)

	if err := verifyPrivate(control); err != nil {
		t.Fatalf("the control descriptor - the private set and nobody else - was refused: %v (trustees %v)", err, resolvedTrustees(t, control))
	}
	err = verifyPrivate(widenedDir)
	if err == nil {
		t.Fatalf("a DACL leaving %s a real grant on the object was accepted as private: trustees %v", everyoneSID, resolvedTrustees(t, widenedDir))
	}
	if !strings.Contains(err.Error(), everyoneSID) {
		t.Fatalf("the refusal did not say which principal is outside the set, by SID: %v", err)
	}
	// The same object read through the notice path: a stranger standing on this
	// object's own DACL is exactly what an operator has to be told about, and the
	// notice has to name it by the one form that cannot be another account.
	foreign, err := explicitForeignPrincipals(widenedDir)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, f := range foreign {
		if strings.Contains(f, everyoneSID) {
			found = true
		}
	}
	if !found {
		t.Fatalf("the notice did not report the explicit foreign grant by its SID: %v", foreign)
	}
}

func setKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
