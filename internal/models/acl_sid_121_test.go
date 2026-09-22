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
// WHY THE PARSER IS HERE AND UNTAGGED. The identity logic is a text function
// over an icacls dump, so it can be judged on both CI legs instead of only on
// the Windows one: the synthetic dumps below are the mutation denominator for
// the shape name-matching cannot express at all (a trustee printed as a bare
// numeric SID, and a short-spelled name that resolves to the wanted SID), while
// the live-object seeding stays where it belongs, in the _windows_test.go files
// that have icacls. Nothing in this file shells out.
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
)

// aceLine is one trustee-right pair from an icacls descriptor dump.
type aceLine struct {
	// Trustee is what icacls printed as the holder: either an account name
	// ("BUILTIN\\Users", "DESKTOP-X\\swq") or a bare numeric SID it could not
	// resolve. It is never treated as an identity by itself.
	Trustee string
	// Rights is the "(...)" payload verbatim, e.g. "(I)(OI)(CI)(M)".
	Rights string
}

// parseACELines turns an icacls dump for one path into ACE lines.
//
// Two parsing traps this has to survive, both learned from this repository's
// earlier ACL work (winsec acl_windows_test.go's comment records the first one
// as the bug that let a foreign Modify ACE sit in the header line and go
// uncounted): the dump glues the path to the FIRST ACE instead of putting it on
// its own line, and the trailer ("Successfully processed ...") is not an ACE.
func parseACELines(raw, path string) []aceLine {
	var out []aceLine
	body := strings.ReplaceAll(raw, "\r\n", "\n")
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(strings.ReplaceAll(line, path, ""))
		i := strings.LastIndex(line, ":(")
		if i <= 0 {
			continue
		}
		out = append(out, aceLine{
			Trustee: strings.TrimSpace(line[:i]),
			Rights:  line[i+1:],
		})
	}
	return out
}

// isNumericSID reports whether icacls printed an unresolved trustee, in which
// case the text already IS the identity and must not be sent for translation.
func isNumericSID(s string) bool {
	return strings.HasPrefix(strings.ToUpper(strings.TrimSpace(s)), "S-1-")
}

// sidIdentity is the answer "whose grant is this", as a SID. resolve translates
// an account name the way a live test does (a subprocess); in the synthetic legs
// below it is a table, which is what lets these cases run on any platform.
func (a aceLine) sidIdentity(resolve func(string) (string, error)) (string, error) {
	if isNumericSID(a.Trustee) {
		return strings.ToUpper(strings.TrimSpace(a.Trustee)), nil
	}
	if resolve == nil {
		return "", fmt.Errorf("no resolver for trustee %q", a.Trustee)
	}
	sid, err := resolve(a.Trustee)
	if err != nil {
		return "", err
	}
	if !isNumericSID(sid) {
		return "", fmt.Errorf("resolver answered %q for %q, which is not a SID", sid, a.Trustee)
	}
	return strings.ToUpper(strings.TrimSpace(sid)), nil
}

// aceLinesForSID returns the lines of an icacls dump that belong to wantSID.
// An entry that cannot be resolved is reported through err rather than dropped:
// "we could not tell whose grant that was" and "that grant is not the one asked
// about" are different answers, and silently equating them is how this hole got
// opened in the first place.
func aceLinesForSID(raw, path, wantSID string, resolve func(string) (string, error)) ([]aceLine, error) {
	want := strings.ToUpper(strings.TrimSpace(wantSID))
	var out []aceLine
	for _, a := range parseACELines(raw, path) {
		got, err := a.sidIdentity(resolve)
		if err != nil {
			return nil, fmt.Errorf("%s %s: %w", a.Trustee, a.Rights, err)
		}
		if got == want {
			out = append(out, a)
		}
	}
	return out, nil
}

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
