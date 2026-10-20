//go:build windows

package models

// Ticket 109 rework leg 109-r1, cell B half 2: R-109-3 closed inside AC#4's own
// nail (acceptance table :44 - "把 R-109-2/R-109-3 两条界线的洞补进 AC#4 的用例"; the
// hole measured red at :246-266).
//
// THE HOLE, IN THE TABLE'S WORDS (:119-123): AC#4's second界 says an explicit
// non-(I) foreign ACE on the install directory is a break. Ticket 109's nail
// answered that by examining the lines belonging to ONE principal
// (`aceLinesOwnedBy(t, path, "BUILTIN\Users")`), so a grant standing on the object
// under any other name never entered the loop and all three break lines read
// "allowed" while the object carried one. The table planted exactly that shape and
// measured it (:248-263): `Everyone:(OI)(CI)(M)`, explicit, non-inherited, holding
// write, with "AC#4 as written sees Users=[BUILTIN\Users:(I)(OI)(CI)(M)] (break line
// 2 reads OK=true)".
//
// WHAT TICKET 121 DID AND DID NOT FIX. 121 AC#4 replaced the name comparison with
// SID resolution (acl_sid_121_test.go), which settles WHICH spelling of Users a line
// belongs to. It left the filter itself in place: the nail still examined one
// principal. This leg removes that scoping - the half the table assigns to AC#4 - and
// keeps the SID discipline.
//
// WHY A TEST-SIDE JUDGMENT AND NOT A winsec CALL. The table draws the line itself
// (:321): "判据实现落 internal/models/ 测试；'该用什么口径' 属 internal/winsec/ 地界
// （票 106/108/112 在飞）⇒ 我只登记不跨界改". The judgment lives here; internal/winsec
// is not touched by one byte.
//
// WHY THE OBJECT-WIDE SET IS SOUND AND NOT MERELY BROADER, measured on this host: an
// install directory this package creates carries no explicit ACE at all - every line
// icacls prints on it is inherited (Administrators / SYSTEM / Authenticated Users /
// Users, each with (I); reading in the ticket face). "Somebody widened this one
// object by hand" is therefore exactly "this object has a non-(I) ACE", whoever its
// trustee is. TestAC4R4NailBreakLineTwoStaysQuietOnInheritedWidth is the reverse leg
// that stops that from decaying into "any grant at all is a break", which would undo
// ticket 95's ruling - the same thing AC#2's reverse half forbids.

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// everyoneSID is Well Known SIDs\Everyone, the principal the table planted (:253).
// Spelled as a SID rather than as the text icacls prints for it, because matching
// printed text is the shape this file exists to retire (a name has more than one
// spelling - acl_sid_121_test.go:26-32).
const everyoneSID = "*S-1-1-0"

// ac4BreakLineTwo is AC#4's second界 as ONE judgment, shared by the nail
// (TestAC4InstallDirWidthIsDeliberateAndHasABreakLine) and by the legs below.
//
// It returns every ACE standing on this object that did not arrive by inheritance,
// whoever its trustee is. Break line 2 is written as a sentence about the object
// ("安装目录自己 DACL 上出现非 (I) 的外来 ACE＝破"), and that is how it is computed.
//
// THE SCOPE THAT WAS REMOVED, stated because the acceptance table measured it: ticket
// 109's nail derived this set from the lines owned by ONE principal
// (`aceLinesOwnedBy(t, path, "BUILTIN\Users")`), so a hand-placed grant for anybody
// else never entered the judgment and all three界 read "allowed" -
// docs/evidence/s1/109-adversarial-acceptance.md:119-123, planted and measured red at
// :248-263. The pre-fix form of this function is registered in
// .scratch/wisp/probes/109/r1/cellB-pre-red.log, and mutation M-B (restoring the
// one-principal filter) turns TestAC4R3 and the nail red again.
//
// WHAT THIS DOES NOT CHANGE: identity resolution stays by SID rather than by the text
// icacls prints (ticket 121 AC#4's correction, acl_sid_121_test.go), and AC#4's first
// and third界 still concern BUILTIN\Users specifically, because those two are
// statements about the width ticket 95 sanctioned for that class rather than about
// this object's DACL.
func ac4BreakLineTwo(t *testing.T, path string) []string {
	t.Helper()
	var out []string
	for _, l := range allTrusteesOn(t, path) {
		if !strings.Contains(l, "(I)") {
			out = append(out, l)
		}
	}
	return out
}

// installIntoSeededStore installs the fixture into a store root carrying the ticket
// 95 width (BUILTIN\Users Modify, inheritable) and hands back the install directory.
func installIntoSeededStore(t *testing.T) (string, *Manager) {
	t.Helper()
	store, m := seededStore(t)
	mgr, err := NewManager(Options{DataDir: store, Manifest: m, VerifySignature: true, Attempts: 2, BackoffBase: 0})
	if err != nil {
		t.Fatal(err)
	}
	dir, err := mgr.Ensure(context.Background(), "kws-fixture")
	if err != nil {
		t.Fatalf("install into the seeded store: %v", err)
	}
	return dir, mgr
}

// TestAC4R3NailBreakLineTwoSeesAnExplicitGrantByAnyPrincipal drives AC#4's second界
// against the table's own planted shape.
func TestAC4R3NailBreakLineTwoSeesAnExplicitGrantByAnyPrincipal(t *testing.T) {
	dir, _ := installIntoSeededStore(t)

	// Premise one: as installed, nobody has widened this object by hand. Without this
	// the reading below could be a measurement of the platform rather than of the界.
	if before := ac4BreakLineTwo(t, dir); len(before) != 0 {
		t.Fatalf("the install dir already carries explicit ACEs before anybody touched it, so nothing below can be "+
			"attributed to the plant: %v (full descriptor: %s)", before, icaclsRun(t, dir))
	}

	icaclsRun(t, dir, "/grant:r", everyoneSID+`:(OI)(CI)(M)`)
	carried := allTrusteesOn(t, dir)
	t.Logf("install dir descriptor after planting an EXPLICIT Everyone:(OI)(CI)(M):\n  %v", carried)

	// The premise the table measured: an explicit, non-inherited WRITE grant really is
	// standing on the object now. Checked independently of the界 under test, from the
	// full trustee list, so a broken界 cannot hide behind a broken plant.
	planted := 0
	for _, l := range carried {
		if !strings.Contains(l, "(I)") && (strings.Contains(l, "Everyone") || strings.Contains(l, "S-1-1-0")) {
			planted++
		}
	}
	if planted != 1 {
		t.Fatalf("R-109-3's plant did not land as an explicit non-inherited Everyone grant (found %d in %v)", planted, carried)
	}

	// The界 AC#4 states, applied through the one judgment the nail uses. Break line 2
	// is written as a sentence about the object - "安装目录自己 DACL 上出现非 (I) 的外来
	// ACE＝破" - so the set it reports has to be the object's, not one account's.
	got := ac4BreakLineTwo(t, dir)
	if len(got) == 0 {
		t.Fatalf("AC#4/R-109-3: break line 2 reports nothing while this object carries an explicit non-inherited "+
			"write grant for Everyone (%s). A hand-placed grant the ticket 95 width ruling does NOT sanction is "+
			"invisible to all three界 because the judgment examines one principal's lines; the object has %d ACE(s) "+
			"that did not arrive by inheritance: %v (full descriptor: %v)",
			"Everyone:(OI)(CI)(M)", planted, got, carried)
	}
	joined := strings.Join(got, "|")
	if !strings.Contains(joined, "Everyone") && !strings.Contains(joined, "S-1-1-0") {
		t.Errorf("the break line must name the principal that widened the object, got %v", got)
	}
	for _, l := range got {
		if strings.Contains(l, "(I)") {
			t.Errorf("break line 2 must not report inherited width as a hand-placed grant: %q", l)
		}
	}
	t.Logf("AC#4 break line 2 catches the hand-placed grant whoever holds it: %v", got)
}

// TestAC4R4NailBreakLineTwoStaysQuietOnInheritedWidth is the reverse leg: the width
// ticket 95 sanctioned arrives by inheritance, so the judgment that must catch a
// hand-placed grant must stay silent about the inherited one.
func TestAC4R4NailBreakLineTwoStaysQuietOnInheritedWidth(t *testing.T) {
	dir, mgr := installIntoSeededStore(t)

	if got := ac4BreakLineTwo(t, dir); len(got) != 0 {
		t.Fatalf("AC#4 reverse leg: a freshly installed tree already breaks the principal-agnostic line - %v. This "+
			"class stays inherit-wide on purpose (ticket 95), so an inherited grant must never read as a break.", got)
	}
	// Independently of the界 under test: the full descriptor of a tree this package
	// installs carries nothing that did not arrive from a parent.
	for _, l := range allTrusteesOn(t, dir) {
		if !strings.Contains(l, "(I)") {
			t.Errorf("AC#4 reverse leg: the freshly installed object already holds a non-inherited ACE, so the "+
				"break line below cannot attribute anything to inheritance: %v", l)
		}
	}
	users := aceLinesOwnedBy(t, dir, `BUILTIN\Users`)
	if len(users) == 0 {
		t.Fatalf("AC#4 first界 premise lost: the install dir inherits nothing, so the width ticket 95 bought is gone (icacls %s)", icaclsRun(t, dir))
	}
	for _, l := range users {
		if !strings.Contains(l, "(I)") {
			t.Errorf("AC#4 first界: the Users grant reached this object without the inherited flag: %q", l)
		}
	}
	if _, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(tinyArchiveMembers[0].Path))); err != nil {
		t.Errorf("installed file unreadable through the inherit-wide dir, so the width bought nothing: %v", err)
	}
	if err := mgr.VerifyInstalled("kws-fixture"); err != nil {
		t.Errorf("the hand-off refused the tree the nail just read as within its width: %v", err)
	}
	t.Logf("principal-agnostic break line stays quiet on the sanctioned inherited width: %v", users)
}
