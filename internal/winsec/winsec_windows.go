//go:build windows

package winsec

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// The three principals a private file may name: the current user, and the two
// every Windows backup, AV filter and recovery path assumes. Everything else
// has to go - that removal, not the addition, is what closes ticket 89.
const (
	sidSystem = "S-1-5-18"
	sidAdmins = "S-1-5-32-544"
	// aceCountPrivate is the shape of the DACL this file builds; used only to
	// confirm the descriptor landed (a filesystem that cannot store one - FAT,
	// some network shares - reports success and then reads back something
	// else). Which SIDs are in there is verified by the tests with icacls,
	// because that is the evidence this ticket accepts.
	aceCountPrivate = 3
)

// applyDescriptor is the seam that puts a private DACL on path. Tests replace
// it to inject a failure (AC#5): the point of the seam is that every caller
// above it refuses to place bytes rather than proceeding wide.
var applyDescriptor = applyDescriptorWindows

// narrowNotice is the audit record of one thing a seal did that used to be
// silent: it cleared a principal that stood on this object's DACL *in its own
// right* (not inherited), i.e. somebody had granted it there deliberately.
//
// Ticket 89's acceptance killed the ticket's earlier claim that such a grant
// "fails at the write point": a PROTECTED DACL means a foreign principal on the
// parent never reaches the child the seal writes, so nothing fails there. What
// actually happens is that the next SealDir/SealFile wipes the out-of-band
// grant. That is a *bigger* risk than a refusal - an operator adds a service
// account, sees no error, and the authorization evaporates at the next open of
// the store. Keeping the whitelist closed (acceptance's ruling: the one promise
// here is "only I", and any hole in it un-does the PROTECTED leg) while making
// the removal loud is the branch this takes: the policy stays "winsec owns the
// grants on this tree", and the human who just lost an ACE learns it from a
// warning carrying the path and the principal.
type narrowNotice struct {
	Path string
	// Principals are the SDDL trustee tokens cleared off this object, e.g. "WD"
	// (Everyone) or a bare "S-1-5-6".
	Principals []string
}

// noticeNarrowed is the seam tests swap. The default writes to the default slog
// logger, which is where the app's own log pipeline (internal/observe) already
// lives - winsec must not import it (the graph runs observe -> secret -> winsec,
// so that edge would close a cycle), and a log line is the narrowest channel
// that is still auditable.
var noticeNarrowed = func(n narrowNotice) {
	slog.Warn("winsec: seal cleared principals that were placed on this object explicitly",
		"path", n.Path,
		"cleared", strings.Join(n.Principals, ","),
		"policy", "winsec owns the grants on this tree; out-of-band ACEs are removed at the next seal")
}

// explicitForeignPrincipals lists the ACEs on path that (a) name somebody
// outside the private set and (b) are on this object *itself*. The second half
// is what keeps the notice a signal: the inherited junk this package exists to
// clear is on effectively every file in a profile directory, so reporting those
// would drown the one case anybody needs to hear about, an explicit grant.
// Inheritance is read from the ACE header's INHERITED_ACE bit rather than from
// the "ID" letters in the rendering, and the materialised inherit-only
// companions of our own grants are inside the set anyway.
func explicitForeignPrincipals(path string) ([]string, error) {
	object, err := readDACL(path)
	if err != nil {
		return nil, err
	}
	set, err := privateSetSIDSet()
	if err != nil {
		return nil, err
	}
	var out []string
	for _, ace := range object.aces {
		if ace.inherited {
			continue // inherited from a parent: not this object's own grant
		}
		if ace.grant && set[ace.trustee] {
			continue
		}
		// A deny ACE, an audit ACE, an object ACE, anything that is not one of
		// our three grants: it is about to be gone, and whoever put it there
		// should hear about it - named by the one form of a principal that
		// cannot be confused with somebody else's account.
		out = append(out, ace.trustee+"("+ace.text+")")
	}
	return out, nil
}

// objectDACL is one object's DACL in the form this package judges it by: the
// descriptor as the OS renders it (for the messages) plus every ACE resolved to
// the trustee's SID. Names are never part of a judgment here.
type objectDACL struct {
	sddl      string
	protected bool
	aces      []resolvedACE
}

type resolvedACE struct {
	// trustee is the SID string of the principal the ACE names, read out of the
	// binary ACE. "unreadable" for the ACE shapes whose trustee sits at another
	// offset (object ACEs), which fail closed.
	trustee     string
	grant       bool
	inherited   bool
	inheritOnly bool
	// text is the OS's own SDDL rendering of this ACE, kept for the messages and
	// the audit notice: it is a representation, never an input to a decision.
	text string
}

// readDACL pulls the DACL back off the object by name. Reading it back is the
// point: SetNamedSecurityInfo reporting nil is not evidence that the descriptor
// landed, and the name form also re-checks that the path still resolves to the
// object we meant.
func readDACL(path string) (objectDACL, error) {
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return objectDACL{}, err
	}
	if sd == nil {
		return objectDACL{}, errors.New("read-back returned no descriptor")
	}
	out := objectDACL{sddl: sd.String()}
	if out.sddl == "" {
		return objectDACL{}, errors.New("read-back produced no SDDL")
	}
	out.protected = strings.Contains(out.sddl, "D:P")
	dacl, _, err := sd.DACL()
	if err != nil {
		return objectDACL{}, err
	}
	if dacl == nil {
		// A NULL DACL grants everything to everyone; there is nothing to walk.
		// The missing-DACL leg of verifyPrivate is what refuses this object.
		return out, nil
	}
	// aceGroups renders the same ACEs in the same order, so zipping them gives
	// each judgment its representation for the message next to it.
	texts := aceGroups(out.sddl)
	for i := uint32(0); i < uint32(dacl.AceCount); i++ {
		var ace *windows.ACCESS_ALLOWED_ACE
		if err := windows.GetAce(dacl, i, &ace); err != nil {
			return objectDACL{}, err
		}
		r := resolvedACE{
			grant:       ace.Header.AceType == windows.ACCESS_ALLOWED_ACE_TYPE,
			inherited:   ace.Header.AceFlags&windows.INHERITED_ACE != 0,
			inheritOnly: ace.Header.AceFlags&windows.INHERIT_ONLY_ACE != 0,
		}
		switch ace.Header.AceType {
		case windows.ACCESS_ALLOWED_ACE_TYPE, windows.ACCESS_DENIED_ACE_TYPE:
			r.trustee = (*windows.SID)(unsafe.Pointer(&ace.SidStart)).String()
		default:
			r.trustee = "unreadable"
		}
		if int(i) < len(texts) {
			r.text = texts[i]
		} else {
			r.text = "ace[" + fmt.Sprint(i) + "]"
		}
		out.aces = append(out.aces, r)
	}
	return out, nil
}

// applyDescriptorWindows builds an explicit DACL and writes it into the object's
// security descriptor with SE_DACL_PROTECTED, which is what "inheritance cannot
// widen this" means at the API level: without the flag the parent's grants stay
// in force underneath ours, and the file keeps the foreign ACE it inherited at
// birth.
func applyDescriptorWindows(path string, dir bool) error {
	// Read the before-state first: after the set there is nothing left to
	// compare against. A failure here is not fatal to the seal (the set and the
	// verify below decide that), it only costs the audit line.
	before, beforeErr := explicitForeignPrincipals(path)
	acl, err := privateACL(dir)
	if err != nil {
		return wrapPath(path, err)
	}
	// SetNamedSecurityInfo takes the object by name, so it also validates that
	// the path still resolves to the object we meant (a rename between create
	// and here fails instead of sealing the wrong file).
	if err := windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION,
		nil, nil, acl, nil); err != nil {
		return wrapPath(path, err)
	}
	if err := verifyPrivate(path); err != nil {
		return err
	}
	if beforeErr == nil && len(before) > 0 {
		noticeNarrowed(narrowNotice{Path: path, Principals: before})
	}
	return nil
}

// sealHandle seals a file that is open but still empty. It goes through the
// *name*, not the handle: SetSecurityInfo needs WRITE_DAC on the handle, while
// the handle os.OpenFile hands back only carries the access it asked for
// (GENERIC_WRITE). The owner holds WRITE_DAC over the object implicitly, so the
// named form works where the handle form reports "Access is denied" - which is
// what the first version of this function did.
func sealHandle(f *os.File) error { return applyDescriptor(f.Name(), false) }

// privateACL is the DACL: the current user plus SYSTEM and Administrators, and
// for a directory the inheritance flags that hand exactly that to children.
//
// ACLFromEntries (SetEntriesInAcl) rather than BuildSecurityDescriptor: the
// latter returns a *self-relative* descriptor, whose dacl member is an offset
// and reads back as "no DACL here" - a first attempt at this function failed
// closed on that, which is at least the right direction to fail in.
func privateACL(dir bool) (*windows.ACL, error) {
	entries, sids, err := privateEntries(dir)
	if err != nil {
		return nil, err
	}
	acl, err := windows.ACLFromEntries(entries, nil)
	// The trustee entries hold raw pointers into those SIDs; they have to
	// outlive the conversion, and nothing has to outlive them afterwards
	// because the ACL is a copy on the Go heap.
	runtime.KeepAlive(entries)
	freeSIDs(sids)
	if err != nil {
		return nil, err
	}
	if acl == nil {
		return nil, errors.New("built an empty ACL")
	}
	return acl, nil
}

// privateEntries returns the grant list. GENERIC_ALL rather than a hand-picked
// access mask: an artifact has to be readable by the store, replaceable by the
// retry path and deletable by the quota path, and a partial mask produces
// "access denied" bugs that outvote the safety this is supposed to add.
func privateEntries(dir bool) ([]windows.EXPLICIT_ACCESS, []*windows.SID, error) {
	me, err := currentUserSID()
	if err != nil {
		return nil, nil, err
	}
	others, err := sidStrings(sidSystem, sidAdmins)
	if err != nil {
		_, _ = windows.LocalFree(windows.Handle(unsafe.Pointer(me)))
		return nil, nil, err
	}
	sids := append(others, me)
	var inherit uint32
	if dir {
		// Children of a sealed directory inherit exactly this grant and nothing
		// else, which is what lets one seal cover the files we never create
		// ourselves (SQLite's -wal/-shm, a pre-rename temp).
		inherit = windows.OBJECT_INHERIT_ACE | windows.CONTAINER_INHERIT_ACE
	}
	entries := make([]windows.EXPLICIT_ACCESS, 0, aceCountPrivate)
	for _, sid := range sids {
		entries = append(entries, windows.EXPLICIT_ACCESS{
			AccessPermissions: windows.GENERIC_ALL,
			AccessMode:        windows.GRANT_ACCESS,
			Inheritance:       inherit,
			Trustee: windows.TRUSTEE{
				TrusteeForm:  windows.TRUSTEE_IS_SID,
				TrusteeType:  windows.TRUSTEE_IS_SID,
				TrusteeValue: windows.TrusteeValueFromSID(sid),
			},
		})
	}
	return entries, sids, nil
}

// freeSIDs releases the ConvertStringSidToSid allocations: without it every
// sealed file leaks allocator memory into a long-running assistant process,
// which only shows up as an OOM a few hundred artifacts deep.
func freeSIDs(sids []*windows.SID) {
	for _, sid := range sids {
		if sid != nil {
			_, _ = windows.LocalFree(windows.Handle(unsafe.Pointer(sid)))
		}
	}
}

func sidStrings(s ...string) ([]*windows.SID, error) {
	out := make([]*windows.SID, 0, len(s))
	for _, str := range s {
		sid, err := convertSID(str)
		if err != nil {
			return nil, err
		}
		out = append(out, sid)
	}
	return out, nil
}

func convertSID(str string) (*windows.SID, error) {
	var sid *windows.SID
	if err := windows.ConvertStringSidToSid(windows.StringToUTF16Ptr(str), &sid); err != nil {
		return nil, fmt.Errorf("winsec: sid %q: %w", str, err)
	}
	// Converted SIDs are allocated by the allocator in advapi32 and owned by
	// the caller; BuildSecurityDescriptor copies what it needs out of them.
	return sid, nil
}

// currentUserSID is the SID the os/user package reports for this token (on
// Windows that is a "S-1-..." string, not a uid), so no token query and no
// locale-dependent account name is involved.
func currentUserSID() (*windows.SID, error) {
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("winsec: current user: %w", err)
	}
	if u.Uid == "" || u.Uid[:2] != "S-" {
		return nil, fmt.Errorf("winsec: current user has no SID: %q", u.Uid)
	}
	return convertSID(u.Uid)
}

// verifyPrivate reads the descriptor back and fails if it is not exactly the
// private one that was asked for. SetNamedSecurityInfo returning nil is not
// enough - the entire defect this package exists for is a call that reports
// success while the object keeps the permissions it had - and the check is made
// on the ACEs themselves, resolved to the SIDs they name, because that says
// *who* holds the grant rather than merely how many entries are present. The
// OS's own SDDL rendering rides along in every message, since a number of ACEs
// that fails to be the property is also the thing an operator has to read. So
// "no principal outside the private set appears" is enforced at runtime, not
// only asserted in a test.
func verifyPrivate(path string) error {
	object, err := readDACL(path)
	if err != nil {
		return wrapPath(path, err)
	}
	if !object.protected {
		return wrapPath(path, fmt.Errorf("DACL is missing or unprotected, so inherited grants still apply: %s", object.sddl))
	}
	set, me, err := privateSet()
	if err != nil {
		return wrapPath(path, err)
	}
	var foreign []string
	grantsMe := false
	for _, ace := range object.aces {
		if !ace.grant || !set[ace.trustee] {
			foreign = append(foreign, ace.trustee+"("+ace.text+")")
			continue
		}
		// An inherit-only ACE hands a right to children and holds none over this
		// object, so it must not be counted as "the current user can read this".
		if ace.inheritOnly {
			continue
		}
		if ace.trustee == me {
			grantsMe = true
		}
	}
	// Counting ACEs is not the property: on a container the OS materialises
	// inherit-only companions for the grants we asked for, so a directory's
	// private descriptor legitimately carries six entries. The property is that
	// *nothing outside the private set is named*, and that the owner is not
	// locked out by the repair.
	if len(foreign) > 0 {
		return wrapPath(path, fmt.Errorf("DACL names principals outside the private set (%v): %s", foreign, object.sddl))
	}
	if !grantsMe {
		return wrapPath(path, fmt.Errorf("DACL grants the current user nothing on this object: %s", object.sddl))
	}
	return nil
}

// privateSet is the set of principals a private object may name, in the one
// form the allow side of this package accepts: a resolved SID string. SYSTEM and
// Administrators by their SIDs, and the account this process runs as.
//
// There is deliberately no name in here. SDDL renders a principal by a
// well-known *name* whenever the OS has one for that SID - S-1-5-18 comes back as
// "SY", S-1-5-32-544 as "BA", and the built-in account with RID 500 as "LA" - so
// a whitelist of names has to be maintained against whatever the machine decides
// to print, and the previous version of this function was exactly that: it listed
// "SY" and "BA" and the current user's SID as Go renders it, which passed on a
// box where the job runs as an ordinary account and refused the descriptor this
// package itself wrote on the CI runner, where the job's token *is* the built-in
// Administrator and its own grant came back spelled "LA" (ticket 106). Comparing
// SIDs is not a wider set: "LA" is admitted only when the local Administrator's
// SID is the token's SID, which is the same fact stated without a name.
//
// The reject side keeps the representations: a refusal or a notice carries the
// SDDL text alongside the SID, because names are fine to *report* and unsound to
// decide with.
func privateSet() (allowed map[string]bool, me string, err error) {
	u, err := user.Current()
	if err != nil {
		return nil, "", fmt.Errorf("winsec: current user: %w", err)
	}
	if u.Uid == "" || !strings.HasPrefix(u.Uid, "S-") {
		return nil, "", fmt.Errorf("winsec: current user has no SID: %q", u.Uid)
	}
	// Re-render the token's SID through the system's own converter so the set
	// holds exactly the form (*SID).String() produces for a trustee read out of
	// a binary ACE; the two are then comparable as strings with no name, no
	// case-folding and no alias in between.
	me, err = canonicalSIDString(u.Uid)
	if err != nil {
		return nil, "", err
	}
	system, err := canonicalSIDString(sidSystem)
	if err != nil {
		return nil, "", err
	}
	admins, err := canonicalSIDString(sidAdmins)
	if err != nil {
		return nil, "", err
	}
	return map[string]bool{system: true, admins: true, me: true}, me, nil
}

// privateSetSIDSet is privateSet for callers that only have to ask "is this
// inside the set".
func privateSetSIDSet() (map[string]bool, error) {
	set, _, err := privateSet()
	return set, err
}

// canonicalSIDString resolves a SID string to the form this package compares
// on, and releases the conversion's allocation.
func canonicalSIDString(str string) (string, error) {
	sid, err := convertSID(str)
	if err != nil {
		return "", err
	}
	out := sid.String()
	if _, err := windows.LocalFree(windows.Handle(unsafe.Pointer(sid))); err != nil {
		return "", err
	}
	return out, nil
}

// aceGroups pulls the parenthesised ACEs out of an SDDL descriptor string.
func aceGroups(sddl string) []string {
	var out []string
	for len(sddl) > 0 {
		i := strings.IndexByte(sddl, '(')
		if i < 0 {
			break
		}
		sddl = sddl[i+1:]
		j := strings.IndexByte(sddl, ')')
		if j < 0 {
			break
		}
		out = append(out, sddl[:j])
		sddl = sddl[j+1:]
	}
	return out
}

func sealFile(path string) error { return applyDescriptor(path, false) }

// sealDir narrows a directory and, because inheritance is decided at creation
// time, propagates the same private descriptor to everything already inside it.
// Without the walk, sealing a data root after the fact would leave every
// artifact, blob and sidecar written during the unsealed window exactly as wide
// as it was born. The walk never follows a link (see removeUnlinked's reasoning
// for why that would hand the delete radius to whoever planted it).
func sealDir(path string) error {
	if err := applyDescriptor(path, true); err != nil {
		return err
	}
	return propagatePrivate(path)
}

func propagatePrivate(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil // the subtree vanished under us; there is nothing left to narrow
		}
		// Not-a-permission-error-skip: an entry we cannot enumerate is an entry
		// we cannot prove private, and "could not check, so assume fine" is the
		// wide direction AC#5 forbids.
		return fmt.Errorf("%w: %s: cannot walk the subtree to narrow it: %v", ErrNotSealable, dir, err)
	}
	for _, e := range entries {
		full := filepath.Join(dir, e.Name())
		if isReparsePoint(full) {
			continue // never traverse a link
		}
		if e.IsDir() {
			if err := applyDescriptor(full, true); err != nil {
				return err
			}
			if err := propagatePrivate(full); err != nil {
				return err
			}
			continue
		}
		if err := applyDescriptor(full, false); err != nil {
			return err
		}
	}
	return nil
}

func isReparsePoint(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	attr, ok := info.Sys().(*syscall.Win32FileAttributeData)
	return ok && attr.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0
}

// removeUnlinked deletes the entry at path, never whatever it points at.
//
// A plain os.Remove on Windows cannot unlink a link whose target is a
// non-empty directory (A51②): RemoveDirectoryW resolves the reparse point and
// then reports the *target's* contents as the reason not to delete, so the
// caller is left with a stray subtree it cannot reclaim - and the tempting
// "fix", os.RemoveAll, would instead delete that target tree. Opening the link
// with FILE_FLAG_OPEN_REPARSE_POINT and asking the handle to mark itself
// delete-pending is the form that removes the link and leaves the tree behind.
func removeUnlinked(path string) error {
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	attr, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if ok && attr.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 {
		if err := deleteLink(path); err != nil {
			// The path goes out plainly, not %q: a reclaim loop logs this and a
			// human has to be able to grep it.
			return fmt.Errorf("%w %s: %v: refusing to recurse into whatever it points at",
				ErrIsReparsePoint, path, err)
		}
		return nil
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}

var (
	modkernel32                       = windows.NewLazySystemDLL("kernel32.dll")
	procSetFileInformationByHandle    = modkernel32.NewProc("SetFileInformationByHandle")
	fileDispositionInfoEx             = uint32(windows.FileDispositionInfoEx)
	fileDispositionFlagDelete         = uint32(windows.FILE_DISPOSITION_DELETE)
	fileDispositionFlagPosixSemantics = uint32(0x2) // FILE_DISPOSITION_FLAG_POSIX_SEMANTICS
	fileDispositionFlagIgnoreReadOnly = uint32(0x4)
)

type fileDispositionInformationEx struct{ Flags uint32 }

// deleteLink is the seam the AC#4 "named error" leg injects through: the real
// disposition call is the only way to remove such a link, and when it fails the
// caller must hear about it as ErrIsReparsePoint rather than as a silent skip.
var deleteLink = deleteReparsePoint

func deleteReparsePoint(path string) error {
	h, err := windows.CreateFile(windows.StringToUTF16Ptr(path),
		windows.DELETE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil, windows.OPEN_EXISTING,
		windows.FILE_FLAG_OPEN_REPARSE_POINT|windows.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if err != nil {
		return err
	}
	defer func() { _ = windows.CloseHandle(h) }()
	info := fileDispositionInformationEx{
		Flags: fileDispositionFlagDelete | fileDispositionFlagIgnoreReadOnly | fileDispositionFlagPosixSemantics,
	}
	r1, _, e1 := syscall.SyscallN(procSetFileInformationByHandle.Addr(),
		uintptr(h), uintptr(fileDispositionInfoEx),
		uintptr(unsafe.Pointer(&info)), unsafe.Sizeof(info))
	if r1 == 0 {
		if e1 != 0 {
			return syscall.Errno(e1)
		}
		return errors.New("SetFileInformationBy(FileDispositionInfoEx) failed")
	}
	return nil
}

// ancestorIsLink is RemoveUnlinked's half of ticket 103's AC#2, in the platform's
// own words: the reparse attribute, the same predicate isReparsePoint and
// therefore the same predicate the sealing walk uses, so "is this a link" cannot
// mean two different things inside one package. A junction reports no
// ModeSymlink to os.Lstat on Windows, which is why the portable version in
// winsec_other.go could not have been used here.
func ancestorIsLink(prefix string) bool { return isReparsePoint(prefix) }
