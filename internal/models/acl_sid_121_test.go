package models

// Ticket 121 AC#4, second half: R-109-3 from acceptor-ticket109b, "按名字而不是
// 按 SID 过滤 ACL".
//
// Ticket 109's Windows leg measures who actually holds the write in the hand-off
// window by reading what icacls prints. The version of that measurement shipped
// in ticket 109 filtered those lines by comparing the trustee TEXT against a
// hand-typed account name (aceLinesNaming / namesPrincipal: strings.EqualFold on
// whatever icacls chose to print). internal/winsec's own standing rule is the
// opposite, and it is not a style preference: an account name is a label some
// other authority resolves, so "does this ACE belong to the principal I am
// asking about" has to be answered with a SID (ticket 118's
// TestAC118NamesEveryoneIsJudgedBySIDNotByTwoBytes is one instance of the same
// correction landing there). This file moves internal/models' two ACL-reading
// helpers onto that footing.
//
// THE PARSER LIVES IN WINSEC NOW (ticket 440). The tree used to carry this
// parser twice - here and in internal/winsec's acl_windows_test.go - as two
// test helpers that did not know about each other and answered different
// questions: this one kept the whole ACE line (trustee plus the "(...)" rights
// payload, which is what let the synthetic cases below assert on Rights), the
// other returned only SIDs and names. A criterion written on the second could
// not ask "is this ACE a deny / inherited", because the markers live in the
// payload it dropped. Ticket 440 merged them into winsec.ParseACELines /
// ACELine / SIDIdentity / ACELinesForSID (internal/winsec/aclparse.go) - the
// package every ACL-reading test already reaches for, and one both CI tiers
// that run these tests already pin. What stays here is the synthetic-dump
// denominator: the identity logic is a text function over an icacls dump, so
// judging it on both CI legs instead of only the Windows one needs a resolver
// that is a table, not a subprocess.
//
// THE SHAPE BEING DEFENDED AGAINST, once, concretely: icacls prints the SID
// instead of a name whenever it cannot resolve the trustee, and real profiles on
// this machine's hosts do carry such ACEs (observed verbatim while writing this:
// `S-1-5-21-3623186960-731165060-4091685855-1717338598:(I)(OI)(CI)(M,DC)` on a
// freshly created temp directory). A name filter reads that line as "not
// BUILTIN\Users" and reports nothing, which is a green that says "nobody foreign
// can write here" about a Modify grant nobody can even name.

import (
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/winsec"
)

// aceLine is the parser's record: one trustee-right pair from an icacls
// descriptor dump. It is an alias, not a new shape - the fields (and the
// rights payload the discarded copy used to drop) are defined once, next to
// the parser.
type aceLine = winsec.ACELine

// parseACELines is the shared icacls parser. Its two parsing traps (the path
// glued to the FIRST ACE, and the "Successfully processed" trailer that is not
// an ACE) are documented once, at the definition.
var parseACELines = winsec.ParseACELines

// isNumericSID reports whether icacls printed an unresolved trustee, in which
// case the text already IS the identity and must not be sent for translation.
var isNumericSID = winsec.IsNumericSID

// aceLinesForSID returns the lines of an icacls dump that belong to wantSID.
// It is the shared filter (internal/winsec/aclparse.go); its one rule this
// ticket exists over - an unresolvable trustee is reported through err, never
// silently dropped, because "cannot tell whose grant that was" is not the same
// answer as "that grant is not the one asked about" - is documented there.
var aceLinesForSID = winsec.ACELinesForSID

// tableResolver is the synthetic dump's directory service.
func tableResolver(m map[string]string) func(string) (string, error) {
	return func(name string) (string, error) {
		if sid, ok := m[name]; ok {
			return sid, nil
		}
		return "", fmt.Errorf("unknown trustee %q in this fixture", name)
	}
}

const usersWanted = "S-1-5-32-545"

// TestAC47TrusteesAreJudgedBySIDNotByText is the portable half of R-109-3: the
// four ways icacls can spell one principal, and what the filter must do about
// each. The old name comparison passes the first case and silently misses the
// other three, which is the whole claim.
func TestAC47TrusteesAreJudgedBySIDNotByText(t *testing.T) {
	const path = `C:\store`
	cases := []struct {
		name  string
		dump  string
		want  int
		match string // trustee text that must be picked
		miss  string // trustee text that must NOT be picked
	}{
		{
			name:  "canonical name resolves to the wanted SID",
			dump:  path + " BUILTIN\\Users:(OI)(CI)(M)\n" + "  NT AUTHORITY\\SYSTEM:(I)(OI)(CI)(F)\n",
			want:  1,
			match: `BUILTIN\Users`,
			miss:  `NT AUTHORITY\SYSTEM`,
		},
		{
			name: "short spelling of the same account is the same principal",
			// The alias shape: icacls answers with whatever the LSA calls the
			// account on this host, and ticket 83's whole lesson in this
			// repository is that a name has more than one spelling.
			dump:  path + " Users:(OI)(CI)(M)\n" + "  BUILTIN\\Administrators:(I)(OI)(CI)(F)\n",
			want:  1,
			match: "Users",
			miss:  `BUILTIN\Administrators`,
		},
		{
			name: "a trustee icacls could not resolve is still an identity",
			// No translation is possible for this line and no name exists for
			// it, so a text filter can only ever answer "not the one I asked
			// about" - a green about a grant it never looked at.
			dump:  path + " S-1-5-21-3623186960-731165060-4091685855-1717338598:(I)(OI)(CI)(M,DC)\n" + "  BUILTIN\\Users:(I)(OI)(CI)(M)\n",
			want:  1,
			match: `BUILTIN\Users`,
			miss:  "S-1-5-21-3623186960-731165060-4091685855-1717338598",
		},
		{
			name:  "the wanted principal arriving as a bare SID is found, not skipped",
			dump:  path + " S-1-5-32-545:(I)(OI)(CI)(M)\n" + "  DESKTOP-H\\swq:(I)(OI)(CI)(F)\n",
			want:  1,
			match: "S-1-5-32-545",
			miss:  `DESKTOP-H\swq`,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resolve := tableResolver(map[string]string{
				`BUILTIN\Users`:          usersWanted,
				"Users":                  usersWanted,
				`NT AUTHORITY\SYSTEM`:    "S-1-5-18",
				`BUILTIN\Administrators`: "S-1-5-32-544",
				`DESKTOP-H\swq`:          "S-1-5-21-1-1-1-1001",
			})
			got, err := aceLinesForSID(tc.dump, path, usersWanted, resolve)
			if err != nil {
				t.Fatalf("the filter refused to answer: %v", err)
			}
			if len(got) != tc.want {
				t.Fatalf("wanted %d ACE line(s) owned by %s, got %d from:\n%s",
					tc.want, usersWanted, len(got), tc.dump)
			}
			if !strings.Contains(got[0].Trustee, tc.match) {
				t.Errorf("the line found is not the one expected (%q), got %q", tc.match, got[0].Trustee)
			}
			if strings.Contains(got[0].Trustee, tc.miss) {
				t.Errorf("the line found belongs to the principal that must be excluded: %q", got[0].Trustee)
			}
			if !strings.Contains(got[0].Rights, "(") {
				t.Errorf("rights payload lost while parsing: %q", got[0].Rights)
			}
		})
	}
}

// TestAC48WrongSIDFindsNothingIsTheNegativeControl that stops the filter from
// being "count every ACE line". A name-based filter cannot express this check at
// all, which is why R-109-3 is not cosmetic: without the identity comparison
// there is no way to ask "and does anybody ELSE hold a grant here".
func TestAC48WrongSIDFindsNothing(t *testing.T) {
	const path = `C:\store`
	dump := path + " BUILTIN\\Users:(OI)(CI)(M)\n" + "  BUILTIN\\Users:(I)(OI)(CI)(RX)\n" + "\nSuccessfully processed 1 files; Failed processing 0 files\n"
	resolve := tableResolver(map[string]string{`BUILTIN\Users`: usersWanted})

	owned, err := aceLinesForSID(dump, path, usersWanted, resolve)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if len(owned) != 2 {
		t.Fatalf("both Users ACEs must be found (inheritable and explicit), got %d: %v", len(owned), owned)
	}
	for _, other := range []string{"S-1-5-32-544", "S-1-5-18", "S-1-5-21-3623186960-731165060-4091685855-1717338598"} {
		misses, err := aceLinesForSID(dump, path, other, resolve)
		if err != nil {
			t.Fatalf("resolve for %s: %v", other, err)
		}
		if len(misses) != 0 {
			t.Errorf("the filter matched %d ACE line(s) to %s, which holds nothing in this descriptor: %v",
				len(misses), other, misses)
		}
	}
	if all := parseACELines(dump, path); len(all) != 2 {
		t.Errorf("the parser must skip the trailer and keep both ACEs, got %d", len(all))
	}
}

// TestAC49UnresolvableTrusteeIsReportedNotSkipped pins the difference between
// "not this principal" and "cannot tell": an ACE whose name the resolver does not
// know must come back as an error, because quietly dropping it is exactly how a
// foreign grant disappears from a security measurement.
func TestAC49UnresolvableTrusteeIsReportedNotSkipped(t *testing.T) {
	const path = `C:\store`
	dump := path + " DOMAIN\\nobody-here:(OI)(CI)(M)\n"
	_, err := aceLinesForSID(dump, path, usersWanted, tableResolver(map[string]string{}))
	if err == nil {
		t.Fatal("an unresolvable trustee was skipped, which turns an unknown grant into a green")
	}
	if !strings.Contains(err.Error(), "nobody-here") {
		t.Errorf("the error must name the trustee it could not resolve: %v", err)
	}
}

// TestAC440TheSharedParserKeepsTheRightsPayload is ticket 440's paired
// criterion over the merge: it is RED against the discarded copy's shape and
// GREEN against the merged one. The dump is a live-printed descriptor shape
// (the same two ACE forms a real sealed dir shows): a deny ACE, and an
// inherited one - the two markers a "this ACE is allowed, not denied, and not
// inherited" seal judgement asks about, both printed AFTER the ":(" the
// parsers split on.
//
// Why the red half is a shape assertion instead of a call to a dead symbol:
// the identity-only shape (per-line SID+name strings, no payload) is what
// internal/winsec's aclSIDs returned before the merge, and the merge deleted
// that shape on purpose. A test may not resurrect a deleted symbol, so the
// red half pins the shape itself and the two questions it cannot answer; the
// mutation-shaped reading of the same pair (strip the Rights field from
// ACELine and every assertion here fails) is in
// docs/evidence/s1/440-acl-parser-drift.md.
func TestAC440TheSharedParserKeepsTheRightsPayload(t *testing.T) {
	const path = `C:\store`
	dump := path + ` BUILTIN\Users:(DENY)(DC)` + "\n" +
		`  NT AUTHORITY\SYSTEM:(I)(F)` + "\n" +
		"\nSuccessfully processed 1 files; Failed processing 0 files\n"

	// GREEN, the merged parser: the payload survives the parse, so a judgement
	// on it can tell a deny from a grant and an inherited ACE from an explicit
	// one.
	lines := parseACELines(dump, path)
	if len(lines) != 2 {
		t.Fatalf("the parser must keep both ACEs and skip the trailer, got %d: %v", len(lines), lines)
	}
	if lines[0].Trustee != `BUILTIN\Users` || lines[0].Rights != "(DENY)(DC)" {
		t.Errorf("the deny marker did not survive the parse: %+v", lines[0])
	}
	if lines[1].Rights != "(I)(F)" {
		t.Errorf("the inherited marker did not survive the parse: %+v", lines[1])
	}

	// RED, the discarded copy's shape: identity-only rows - the (sids, names)
	// return shape aclSIDs carried before ticket 440. It still answers WHOSE
	// grant each line is; for both type questions the answer is simply absent,
	// so the two "matched" checks below can never fire - that emptiness is the
	// blind spot, pinned. The live wrapper additionally resolved each name
	// through a subprocess; the resolution is not what was lost, the payload
	// was.
	type identityOnlyRow struct{ SID, Name string }
	table := tableResolver(map[string]string{
		`BUILTIN\Users`:       "S-1-5-32-545",
		`NT AUTHORITY\SYSTEM`: "S-1-5-18",
	})
	var identityOnly []identityOnlyRow
	for _, a := range lines {
		sid, err := a.SIDIdentity(table)
		if err != nil {
			t.Fatalf("the identity half of the merged parser: %v", err)
		}
		identityOnly = append(identityOnly, identityOnlyRow{SID: sid, Name: a.Trustee})
	}
	if len(identityOnly) != 2 {
		t.Fatalf("the identity half must still see both rows, got %d", len(identityOnly))
	}
	for _, r := range identityOnly {
		if strings.Contains(r.SID, "DENY") || strings.Contains(r.Name, "DENY") ||
			strings.Contains(r.SID, "(I)") || strings.Contains(r.Name, "(I)") {
			t.Errorf("the identity-only shape unexpectedly carries payload markers: %+v", r)
		}
	}
}
