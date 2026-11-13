package winsec

// The one icacls parser every ACL test in this repository shares.
//
// WHY THIS IS PRODUCTION CODE (ticket 440, closing ticket 238 AC#4③). The tree
// grew two independent parsers inside test helpers that did not know about
// each other: internal/models/acl_sid_121_test.go's parseACELines, which kept
// the whole ACE line (trustee text plus the "(...)" rights payload), and this
// package's acl_windows_test.go aclSIDs, which returned only SIDs and names
// and dropped the payload. Ticket 238's closing leg named the shape: along
// with the existing helper comes its blind spot - a judgement written on
// aclSIDs cannot ask "is this ACE a deny, is it inherited", because the marker
// for both lives inside the payload it never kept. Measured on a live
// descriptor while landing this file: icacls prints a deny ACE as
// `BUILTIN\Users:(DENY)(DC)` and an inherited one as `NT AUTHORITY\SYSTEM:(I)(F)`
// - both markers sit after the ":(" the parsers split on. Which parser a new
// criterion picked depended on which package the author happened to be sitting
// in, and ticket 238's ledger carried the drift with no owner until this file.
//
// Shared test instrumentation that two packages need is what a small
// production package is for, and winsec is the package every ACL-reading test
// in the tree already reaches for; the two tiers that run these tests (core
// and winsec in scripts/portable-tests.sh) both already pin it. It lands here
// rather than in a new package so that the coverage rosters - which ticket 440
// has no write surface for - do not move at all.
//
// WHAT THIS DELIBERATELY IS NOT. Not an ACL API: the package's own promise
// stays "current-user-only writes" (winsec.go), and nothing in production
// judgement calls in here - every caller today is a test. It does not resolve
// names: the resolver is a parameter, because the real directory service is a
// subprocess that exists only where the descriptor does (the windows test
// files supply one; internal/models' synthetic cases supply a table, which is
// what lets them run on any platform). It does not decide allow versus deny
// for a caller: ParseACELines hands over the rights payload verbatim, and
// policies ("deny beats allow", "integrity labels are not grants") stay with
// the judgement that owns them.
//
// Callers and their ranges, so the next criterion does not re-grow a private
// copy: internal/winsec's tests parse through the aclSIDs wrapper in
// acl_windows_test.go, which keeps the (sids, names) shape its callers know
// and states the payload-shaped blind spot that shape carries; a criterion
// that needs the ACE's type or inheritance uses ParseACELines or
// ACELinesForSID directly, as internal/models' tests do.

import (
	"fmt"
	"strings"
)

// ACELine is one trustee-right pair from an icacls descriptor dump.
type ACELine struct {
	// Trustee is what icacls printed as the holder: either an account name
	// ("BUILTIN\\Users", "DESKTOP-X\\swq") or a bare numeric SID it could not
	// resolve. It is never treated as an identity by itself.
	Trustee string
	// Rights is the "(...)" payload verbatim, e.g. "(I)(OI)(CI)(M)" or
	// "(DENY)(DC)". This is the field the pre-merge helpers dropped, and with
	// it every judgement that needs the ACE's type or inheritance: allow vs
	// deny, inherited vs explicit, modify vs read. Ticket 440's paired
	// criterion turns on exactly this field.
	Rights string
}

// ParseACELines turns an icacls dump for one path into ACE lines.
//
// Two parsing traps this has to survive, both learned from this repository's
// earlier ACL work (acl_windows_test.go's comment records the first one as the
// bug that let a foreign Modify ACE sit in the header line and go uncounted):
// the dump glues the path to the FIRST ACE instead of putting it on its own
// line, and the trailer ("Successfully processed ...") is not an ACE.
func ParseACELines(raw, path string) []ACELine {
	var out []ACELine
	body := strings.ReplaceAll(raw, "\r\n", "\n")
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(strings.ReplaceAll(line, path, ""))
		i := strings.LastIndex(line, ":(")
		if i <= 0 {
			continue
		}
		out = append(out, ACELine{
			Trustee: strings.TrimSpace(line[:i]),
			Rights:  line[i+1:],
		})
	}
	return out
}

// IsNumericSID reports whether icacls printed an unresolved trustee, in which
// case the text already IS the identity and must not be sent for translation.
func IsNumericSID(s string) bool {
	return strings.HasPrefix(strings.ToUpper(strings.TrimSpace(s)), "S-1-")
}

// SIDIdentity is the answer "whose grant is this", as a SID. resolve translates
// an account name the way a live test does (a subprocess); in the synthetic
// cases (internal/models/acl_sid_121_test.go) it is a table, which is what
// lets those cases run on any platform.
func (a ACELine) SIDIdentity(resolve func(string) (string, error)) (string, error) {
	if IsNumericSID(a.Trustee) {
		return strings.ToUpper(strings.TrimSpace(a.Trustee)), nil
	}
	if resolve == nil {
		return "", fmt.Errorf("no resolver for trustee %q", a.Trustee)
	}
	sid, err := resolve(a.Trustee)
	if err != nil {
		return "", err
	}
	if !IsNumericSID(sid) {
		return "", fmt.Errorf("resolver answered %q for %q, which is not a SID", sid, a.Trustee)
	}
	return strings.ToUpper(strings.TrimSpace(sid)), nil
}

// ACELinesForSID returns the lines of an icacls dump that belong to wantSID.
// An entry that cannot be resolved is reported through err rather than dropped:
// "we could not tell whose grant that was" and "that grant is not the one asked
// about" are different answers, and silently equating them is how this hole got
// opened in the first place.
func ACELinesForSID(raw, path, wantSID string, resolve func(string) (string, error)) ([]ACELine, error) {
	want := strings.ToUpper(strings.TrimSpace(wantSID))
	var out []ACELine
	for _, a := range ParseACELines(raw, path) {
		got, err := a.SIDIdentity(resolve)
		if err != nil {
			return nil, fmt.Errorf("%s %s: %w", a.Trustee, a.Rights, err)
		}
		if got == want {
			out = append(out, a)
		}
	}
	return out, nil
}
