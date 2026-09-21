//go:build windows

package winsec

// Ticket 118's own cases: the two holes acceptor-ticket104 booked as R-104-1 and
// R-104-6, and (further down) the two legs agent-ticket115b measured to have no
// teeth. Nothing in this file judges a production decision - it pins, with cases,
// two judgments the production code already makes and the suite did not hold.
//
// R-104-1 (AC#1, AC#2). The default notifier renders a three-valued `kind=` field
// (explicit / inherited / explicit+inherited) out of which bucket each cleared
// grant landed in. acceptor-ticket104 measured that nothing pinned it: its M5
// (delete the whole switch, so every notice says "explicit") left `go build` rc=0
// and all four AC#1 cases green, and its M4 (swap the two buckets) reddened only
// the structural case while TestAC1DefaultLogSaysInherited stayed green - because
// that case asked `strings.Contains(out, "inherited")`, and the attribute NAME
// `cleared_inherited=` carries the word "inherited" in every rendering, whether
// or not the kind did. So the pins below compare the VALUE of the kind attribute
// as one token, and they compare each bucket's payload by the SID of the
// principal it must hold. A swapped pair of buckets cannot satisfy both.
//
// R-104-6 (AC#3). `namesEveryone` used to answer "does this text contain the two
// bytes WD (or S-1-1-0, or Everyone)?" - which any output can satisfy for the
// wrong reason: a directory named WD-team-share whose only foreign grant is
// NT AUTHORITY\SERVICE reads as Everyone to it. The instrument below reads a
// trustee out of the positions a trustee can actually occupy (icacls aligns every
// record under its trustee column, and SDDL puts the trustee last in an ACE),
// asks this machine which SID that token names, and compares SIDs. That is the
// same rule ticket 106 wrote for the private set and ticket 115 for attribution:
// a representation is not an identity, so decide with the identity.
//
// R-115-b/AC#7 is at the bottom of this file, together with what changed in
// inherited_narrow_notice_104_windows_test.go and why.

import (
	"bytes"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/sys/windows"
)

// serviceSID is NT AUTHORITY\SERVICE: a principal outside the private set whose
// rendering shares no bytes with Everyone's, which is what makes it usable as the
// second bucket's trustee in a fixture that has to tell the two apart.
const serviceSID = "S-1-5-6"

// ---------------------------------------------------------------------------
// AC#1 / AC#2: the rendered kind= field and the two bucket payloads
// ---------------------------------------------------------------------------

// warnLines118 reads the WARN records out of captured log text as attribute maps.
func warnLines118(t *testing.T, text string) []map[string]string {
	t.Helper()
	var out []map[string]string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "level=WARN") {
			continue
		}
		out = append(out, parseAttrs118(t, line))
	}
	return out
}

// renderNoticeText118 is renderNotices118 for a case that has to quote the raw
// lines it judged as well as read them.
func renderNoticeText118(t *testing.T, body func()) (string, []map[string]string) {
	t.Helper()
	var buf bytes.Buffer
	origDefault := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(origDefault) })
	body()
	text := buf.String()
	return text, warnLines118(t, text)
}

// renderNotices118 runs body with the DEFAULT notifier in place - no seam is
// installed, the renderer is the thing under test - and returns every WARN line
// it wrote as its attribute map. A test that reached for `noticeNarrowed` here
// would be reading its own builder instead of the one that ships.
func renderNotices118(t *testing.T, body func()) []map[string]string {
	t.Helper()
	_, lines := renderNoticeText118(t, body)
	return lines
}

// parseAttrs118 splits one slog text line into its attributes. It is a parser and
// not a search, because every assertion below has to be able to say "the value of
// kind is this token and nothing else": a substring test on this line cannot tell
// kind=inherited from the attribute name cleared_inherited=, which is exactly how
// R-104-1 stayed green through M4.
func parseAttrs118(t *testing.T, line string) map[string]string {
	t.Helper()
	out := map[string]string{}
	for i := 0; i < len(line); {
		eq := strings.IndexByte(line[i:], '=')
		if eq < 0 {
			break
		}
		key := strings.TrimSpace(line[i : i+eq])
		rest := i + eq + 1
		var val string
		if rest < len(line) && line[rest] == '"' {
			end := strings.IndexByte(line[rest+1:], '"')
			if end < 0 {
				t.Fatalf("unterminated quoted value in %q", line)
			}
			val = line[rest+1 : rest+1+end]
			i = rest + end + 2
		} else {
			end := strings.IndexByte(line[rest:], ' ')
			if end < 0 {
				val = line[rest:]
				i = len(line)
			} else {
				val = line[rest : rest+end]
				i = rest + end + 1
			}
		}
		if key == "" {
			continue
		}
		if _, dup := out[key]; dup {
			t.Fatalf("attribute %q appears twice in %q, so the line cannot be read as a pair list", key, line)
		}
		out[key] = val
	}
	if _, ok := out["level"]; !ok {
		t.Fatalf("this line carries no level attribute, so it is not a slog record: %q", line)
	}
	return out
}

// bucketSIDs118 turns one rendered bucket ("", or a comma-joined list of
// `SID(ace-text)` tokens as foreignPrincipals builds them) into the SIDs of the
// principals it names. The trustee is read from the field that holds it; nothing
// here searches the text for a name.
func bucketSIDs118(t *testing.T, value string) []string {
	t.Helper()
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	var out []string
	for _, tok := range strings.Split(value, ",") {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		i := strings.IndexByte(tok, '(')
		if i <= 0 {
			t.Fatalf("AC#1: the rendered bucket holds %q, which is not the trustee(ace) shape foreignPrincipals builds", tok)
		}
		sid, ok := resolveTrusteeToken118(t, tok[:i])
		if !ok {
			t.Fatalf("AC#1: the rendered bucket names a trustee this machine cannot resolve: %q", tok)
		}
		// The ace text carries its own trustee in the last field; the two have to
		// agree, or the token was assembled from two different objects.
		if ace := strings.TrimSuffix(strings.TrimSuffix(tok[i:], ")"), ""); strings.Contains(ace, ";") {
			fields := strings.Split(ace, ";")
			if inner, ok := resolveTrusteeToken118(t, strings.TrimSpace(fields[len(fields)-1])); ok {
				if !strings.EqualFold(inner, sid) {
					t.Fatalf("AC#1: token %q names %s as its trustee and %s inside the ace text", tok, sid, inner)
				}
			}
		}
		out = append(out, sid)
	}
	return out
}

// assertRenderedNotice118 is the whole pin for one rendering: kind must be exactly
// this token, and each bucket must name exactly the principals the seal took off
// that bucket. wantExplicit/wantInherited are SID strings.
func assertRenderedNotice118(t *testing.T, lines []map[string]string, wantKind string, wantExplicit, wantInherited []string) map[string]string {
	t.Helper()
	if len(lines) != 1 {
		t.Fatalf("AC#1: the seal rendered %d WARN line(s), want exactly 1: %v", len(lines), lines)
	}
	attrs := lines[0]
	if got := attrs["kind"]; got != wantKind {
		t.Errorf("AC#1/AC#2 RED: the notice rendered kind=%q, want %q. This field is the only thing that tells an operator \"an out-of-band grant on this file is gone\" from \"this file stopped being as reachable as its directory\", and the two buckets carrying it are what the rest of this assertion checks. line=%v", got, wantKind, attrs)
	}
	for _, e := range []struct {
		attr string
		want []string
		got  []string
	}{
		{"cleared", wantExplicit, bucketSIDs118(t, attrs["cleared"])},
		{"cleared_inherited", wantInherited, bucketSIDs118(t, attrs["cleared_inherited"])},
	} {
		if strings.Join(e.got, " ") != strings.Join(e.want, " ") {
			t.Errorf("AC#1/AC#2 RED: the %s attribute names %v, want exactly %v (kind=%q). A bucket whose payload moved to the other attribute is the M4 shape acceptor-ticket104 measured as half-unguarded.",
				e.attr, e.got, e.want, attrs["kind"])
		}
	}
	if attrs["path"] == "" {
		t.Errorf("AC#1: the notice carries no path for an operator to grep: %v", attrs)
	}
	return attrs
}

// grantStandsOn118 is the binary-DACL answer to "does this principal hold a grant
// on this object", read through this package's own DACL reader so the answer is a
// trustee SID and an inherited bit rather than a rendering. wantInherited decides
// which of the two buckets the grant has to stand in.
func grantStandsOn118(t *testing.T, path, sid string, wantInherited bool) bool {
	t.Helper()
	object, err := readDACL(path)
	if err != nil {
		t.Fatalf("readDACL(%s): %v", path, err)
	}
	canonical, err := canonicalSIDString(sid)
	if err != nil {
		t.Fatalf("canonicalize %s: %v", sid, err)
	}
	for _, ace := range object.aces {
		if strings.EqualFold(ace.trustee, canonical) && ace.grant && ace.inherited == wantInherited {
			return true
		}
	}
	return false
}

// requireGrantStandsOn118 is a fixture proof: a case that seals a grant which is
// not there measures nothing, so the shape is asserted from the binary DACL
// before any seal runs.
func requireGrantStandsOn118(t *testing.T, path, sid string, wantInherited bool) {
	t.Helper()
	if grantStandsOn118(t, path, sid, wantInherited) {
		return
	}
	object, err := readDACL(path)
	which := "explicit"
	if wantInherited {
		which = "inherited"
	}
	t.Fatalf("the fixture did not plant an %s %s grant on %s: sddl=%q err=%v", which, sid, path, object.sddl, err)
}

// everyoneStandsOn118 is the binary-DACL answer to "does Everyone still hold a
// grant on this object", in either bucket. It replaces the two places ticket 104
// asked that question of icacls *text*; the text answer cannot say which bucket a
// grant stands in, and the object it is read from is not a rendering.
func everyoneStandsOn118(t *testing.T, path string) bool {
	t.Helper()
	return grantStandsOn118(t, path, everyoneSID, false) || grantStandsOn118(t, path, everyoneSID, true)
}

// noticeRendersKindOnTree118 answers "did the WARN line this seal wrote for THIS
// tree carry kind=<want>?", comparing the VALUE of the attribute as one token.
// R-104-1 is a note about this exact distinction: `strings.Contains(out,
// "inherited")` is true of every rendering the default notifier produces, because
// the attribute is NAMED cleared_inherited=, so it stayed green when M4 swapped the
// buckets and the kind value moved from "inherited" to "explicit".
//
// The tree scoping is not a convenience: a case that installs the capture over a
// whole fixture sees every narrowing in it, and on a machine whose temp root hands
// down foreign ACEs (measured on this box: two S-1-5-21-* trustees arrive with the
// directory) those lines are about other objects. Reading some other object's line
// as this one's is the comparison defect ticket 115 ruled on, so the same rule
// answers here.
func noticeRendersKindOnTree118(t *testing.T, text, spelling, want string) bool {
	t.Helper()
	lines := noticeLinesAboutTree118(t, warnLines118(t, text), spelling)
	if len(lines) != 1 {
		t.Errorf("the capture holds %d WARN line(s) about %s, want exactly 1: %q", len(lines), spelling, text)
		return false
	}
	return lines[0]["kind"] == want
}

// noticeClearsOnTreeExactly118 answers "are the principals this rendered notice
// cleared for THIS tree exactly these, named by SID and read out of the bucket
// fields they stand in?" The SDDL alias of a principal is two bytes long, so
// searching a whole log line for it is the same instrument defect R-104-6 booked
// for Everyone, one principal away: any other field of the same line can satisfy
// it, and so can any other line.
func noticeClearsOnTreeExactly118(t *testing.T, text, spelling string, want ...string) bool {
	t.Helper()
	lines := noticeLinesAboutTree118(t, warnLines118(t, text), spelling)
	if len(lines) != 1 {
		t.Errorf("the capture holds %d WARN line(s) about %s, want exactly 1: %q", len(lines), spelling, text)
		return false
	}
	var got []string
	for _, attr := range []string{"cleared", "cleared_inherited"} {
		got = append(got, bucketSIDs118(t, lines[0][attr])...)
	}
	wantCanonical := make([]string, 0, len(want))
	for _, sid := range want {
		canonical, err := canonicalSIDString(sid)
		if err != nil {
			t.Fatalf("canonicalize %s: %v", sid, err)
		}
		wantCanonical = append(wantCanonical, canonical)
	}
	if strings.Join(got, " ") != strings.Join(wantCanonical, " ") {
		t.Errorf("the notice about %s cleared %v, want exactly %v (line: %v)", spelling, got, wantCanonical, lines[0])
		return false
	}
	return true
}

// noticeLinesAboutTree118 is the tree-scoped filter every assertion above uses:
// the path attribute of a rendered line goes through this package's own
// attribution rule, which is the same rule the captured-struct cases use.
func noticeLinesAboutTree118(t *testing.T, lines []map[string]string, spelling string) []map[string]string {
	t.Helper()
	var out []map[string]string
	for _, l := range lines {
		if l["path"] == "" {
			t.Fatalf("a rendered notice carries no path attribute: %v", l)
		}
		if noticeNamesTree(narrowNotice{Path: l["path"]}, spelling) {
			out = append(out, l)
		}
	}
	return out
}

func privateRoot118(t *testing.T, name string) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), name)
	if err := PrivateDirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	return root
}

func writeChild118(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("artifact"), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// TestAC118KindFieldSaysExplicitForAGrantThatStoodOnTheObjectItself is kind= value
// one of three (AC#1): the ticket 89 bucket, a grant somebody placed on this
// object with their own hand.
func TestAC118KindFieldSaysExplicitForAGrantThatStoodOnTheObjectItself(t *testing.T) {
	root := privateRoot118(t, "store")
	child := writeChild118(t, root, "explicit-only.txt")
	mustExec(t, "icacls", child, "/grant", "*"+serviceSID+":(RX)")
	requireGrantStandsOn118(t, child, serviceSID, false)

	lines := renderNotices118(t, func() {
		if err := SealFile(child); err != nil {
			t.Errorf("SealFile: %v", err)
		}
	})
	assertRenderedNotice118(t, lines, "explicit", []string{serviceSID}, nil)
}

// TestAC118KindFieldSaysInheritedForAGrantTheParentHandedDown is kind= value two of
// three, and the ticket 104 leg: the child never held the grant in its own DACL.
// M5 (the switch gone) reddens this case and the next one, which is the whole of
// AC#1's criterion.
func TestAC118KindFieldSaysInheritedForAGrantTheParentHandedDown(t *testing.T) {
	root := privateRoot118(t, "store")
	mustExec(t, "icacls", root, "/grant", "*"+everyoneSID+":(OI)(CI)(RX)")
	child := writeChild118(t, root, "inherited-only.txt")
	requireGrantStandsOn118(t, child, everyoneSID, true)

	lines := renderNotices118(t, func() {
		if err := SealFile(child); err != nil {
			t.Errorf("SealFile: %v", err)
		}
	})
	assertRenderedNotice118(t, lines, "inherited", nil, []string{everyoneSID})
}

// TestAC118KindFieldSaysBothWhenTheObjectCarriedItsOwnGrantToo is kind= value three
// of three, and the case AC#2 needs on its own: both buckets are non-empty, so
// swapping them cannot change the kind token - only the two payloads, which is
// precisely what "the buckets swapped but the log field did not" looked like when
// it passed unnoticed. This case must therefore go red on M4 while still being
// green on M5-free code, and it does not depend on any other case to say so.
func TestAC118KindFieldSaysBothWhenTheObjectCarriedItsOwnGrantToo(t *testing.T) {
	root := privateRoot118(t, "store")
	mustExec(t, "icacls", root, "/grant", "*"+everyoneSID+":(OI)(CI)(RX)")
	child := writeChild118(t, root, "both-buckets.txt")
	mustExec(t, "icacls", child, "/grant", "*"+serviceSID+":(RX)")
	requireGrantStandsOn118(t, child, everyoneSID, true)
	requireGrantStandsOn118(t, child, serviceSID, false)

	lines := renderNotices118(t, func() {
		if err := SealFile(child); err != nil {
			t.Errorf("SealFile: %v", err)
		}
	})
	attrs := assertRenderedNotice118(t, lines, "explicit+inherited", []string{serviceSID}, []string{everyoneSID})
	// The one thing that made M4 invisible: the word "inherited" is in this line
	// twice over, once as a value and once as an attribute name, so only the
	// value comparison above can tell the two renderings apart. Recorded, so a
	// later reader does not "simplify" this back into a Contains.
	if !strings.Contains(attrs["cleared_inherited"], "WD") && !strings.Contains(attrs["cleared_inherited"], everyoneSID) {
		t.Errorf("AC#2: the inherited bucket does not name Everyone at all: %q", attrs["cleared_inherited"])
	}
}

// ---------------------------------------------------------------------------
// AC#3: namesEveryone, judged by SID
// ---------------------------------------------------------------------------

// everyoneSIDCanonical118 is the value every comparison below ends at.
func everyoneSIDCanonical118(t *testing.T) string {
	t.Helper()
	sid, err := canonicalSIDString(everyoneSID)
	if err != nil {
		t.Fatalf("canonicalize %s: %v", everyoneSID, err)
	}
	return sid
}

var trusteeSIDCache118 = map[string]string{}

// resolveTrusteeToken118 asks this machine which principal a trustee token names
// and answers with that SID. Three forms exist in the wild and all three go
// through the system rather than through a needle: the SDDL alias ("WD"), the
// numeric SID ("S-1-1-0"), and whatever display name the renderer chose
// ("Everyone", `NT AUTHORITY\SERVICE`), which is the form that is localized on a
// non-English box and therefore the one a name whitelist can never carry.
func resolveTrusteeToken118(t *testing.T, token string) (string, bool) {
	t.Helper()
	token = strings.TrimSpace(token)
	if token == "" || len(token) > 261 {
		return "", false
	}
	if cached, ok := trusteeSIDCache118[token]; ok {
		return cached, cached != ""
	}
	sid, ok := "", false
	if converted, err := convertSID(token); err == nil {
		// The same allocator discipline privateEntries documents: a converted SID
		// is the caller's to release.
		sid, ok = converted.String(), true
		freeSIDs([]*windows.SID{converted})
	} else if u, uerr := user.Lookup(token); uerr == nil {
		sid, ok = u.Uid, true
	} else if g, gerr := user.LookupGroup(token); gerr == nil {
		sid, ok = g.Gid, true
	}
	if !ok {
		trusteeSIDCache118[token] = ""
		return "", false
	}
	canonical, err := canonicalSIDString(sid)
	if err != nil {
		t.Fatalf("canonicalize the SID this machine answered for %q (%s): %v", token, sid, err)
	}
	trusteeSIDCache118[token] = canonical
	return canonical, true
}

// trusteeColumns118 returns every token in s that occupies a position a trustee
// can hold. Two renderings are read, both by structure:
//
//   - icacls prints "<subject> <trustee>:(flags)(perms)" for the first record and
//     aligns every further record under the trustee column, so the indentation of
//     the second record is where the first record's subject column ends. Reading
//     the column instead of the line is what keeps a directory named "Everyone",
//     or one named "WD-team-share", out of the answer.
//   - SDDL puts the trustee last inside the parentheses of an ACE, so every ACE
//     group's final field is a trustee - which is also the shape of a
//     narrowNotice bucket token's inner text.
//
// A token is returned whether or not it resolves; the caller decides. Nothing here
// searches for a name.
func trusteeColumns118(s string) []string {
	lines := strings.Split(s, "\n")
	subjectPad := 0
	for _, l := range lines[1:] {
		body := strings.TrimLeft(l, " \t")
		if body == "" || !strings.ContainsAny(body, "(:") {
			continue
		}
		if i := strings.Index(l, body); i > 0 {
			subjectPad = i
			break
		}
	}
	var out []string
	record := func(rec string) {
		rec = strings.TrimSpace(rec)
		if rec == "" {
			return
		}
		if i := strings.IndexByte(rec, '('); i > 0 {
			out = append(out, strings.TrimRight(rec[:i], " \t:"))
			return
		}
		if i := strings.IndexByte(rec, ':'); i > 0 {
			out = append(out, strings.TrimSpace(rec[:i]))
		}
	}
	for i, l := range lines {
		if i == 0 && subjectPad > 0 {
			if len(l) > subjectPad {
				l = l[subjectPad:]
			} else {
				l = ""
			}
		}
		record(l)
	}
	for _, ace := range aceGroups(s) {
		if n := strings.LastIndexByte(ace, ';'); n >= 0 {
			if tok := strings.TrimSpace(ace[n+1:]); tok != "" {
				out = append(out, tok)
			}
		}
	}
	return out
}

// namesEveryone answers "does this text name the Everyone principal?" by reading
// every trustee position it can find, asking this machine which SID each one
// names, and comparing SIDs.
//
// It replaces the two-byte probe acceptor-ticket104 booked as R-104-6, which
// answered that question by searching the whole text for "WD" (or "S-1-1-0", or
// "Everyone"). The old rule could not tell a trustee from any other byte of an
// output: measured on this box, `icacls D:\...\WD-team-share` over a directory
// whose only foreign grant is NT AUTHORITY\SERVICE contains "WD" and names
// Everybody's-but-nobody, and the old answer was "Everyone" - which is a fixture
// proof that passes for a reason unrelated to the fixture. TestAC118NamesEveryoneIsJudgedBySIDNotByTwoBytes
// plants exactly that object and requires the opposite answer, and requires the
// real Everyone grant to still be named, so neither a constant false nor a
// substring search can satisfy it.
func namesEveryone(t *testing.T, s string) bool {
	t.Helper()
	everyone := everyoneSIDCanonical118(t)
	for _, tok := range trusteeColumns118(s) {
		if sid, ok := resolveTrusteeToken118(t, tok); ok && strings.EqualFold(sid, everyone) {
			return true
		}
	}
	return false
}

// TestAC118NamesEveryoneIsJudgedBySIDNotByTwoBytes is AC#3's both halves, on the
// three renderings the instrument actually receives: icacls text, SDDL ACE text,
// and a narrowNotice bucket list.
func TestAC118NamesEveryoneIsJudgedBySIDNotByTwoBytes(t *testing.T) {
	everyone := everyoneSIDCanonical118(t)

	// Half one, forward: every legal representation of a real Everyone grant, on a
	// real object, read out of the OS's own rendering. Without this the new rule
	// could be satisfied by a function that answers false to everything.
	wide := privateRoot118(t, "wide")
	wideChild := writeChild118(t, wide, "everyone-stands-here.txt")
	mustExec(t, "icacls", wideChild, "/grant", "*"+everyoneSID+":(RX)")
	requireGrantStandsOn118(t, wideChild, everyone, false)
	wideText := icaclsRaw(t, wideChild)
	if !namesEveryone(t, wideText) {
		t.Fatalf("AC#3 RED: namesEveryone does not recognise a real explicit Everyone grant read back from the object: %s", wideText)
	}
	object, err := readDACL(wideChild)
	if err != nil {
		t.Fatalf("readDACL: %v", err)
	}
	if !namesEveryone(t, object.sddl) {
		t.Errorf("AC#3 RED: namesEveryone does not recognise the SDDL form of the same grant: %s", object.sddl)
	}

	// Half two, the substitute principal: the same object shape, foreign grant
	// coming from NT AUTHORITY\SERVICE instead. Neither of the three needles the
	// old rule searched for appears in its rendering.
	svc := privateRoot118(t, "svc")
	svcChild := writeChild118(t, svc, "service-stands-here.txt")
	mustExec(t, "icacls", svcChild, "/grant", "*"+serviceSID+":(RX)")
	requireGrantStandsOn118(t, svcChild, serviceSID, false)
	svcText := icaclsRaw(t, svcChild)
	if namesEveryone(t, svcText) {
		t.Fatalf("AC#3 RED: namesEveryone admitted a SERVICE grant as Everyone: %s", svcText)
	}
	if grantStandsOn118(t, svcChild, everyone, false) || grantStandsOn118(t, svcChild, everyone, true) {
		t.Fatalf("AC#3: the fixture is not what this leg claims: Everyone really does stand on %s", svcChild)
	}

	// Half three, R-104-6's actual failure mode: a principal with no relation to
	// Everyone, in an output whose other bytes spell the two letters the old rule
	// searched for. The directory's own name is where they come from, which is the
	// "换一种渲染 / 别的字段名" case the report named.
	impostor := filepath.Join(privateRoot118(t, "WD-team-share"), "artifact.txt")
	if err := os.WriteFile(impostor, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	mustExec(t, "icacls", impostor, "/grant", "*"+serviceSID+":(RX)")
	impostorText := icaclsRaw(t, impostor)
	if !strings.Contains(impostorText, "WD") {
		t.Fatalf("AC#3: the fixture planted no WD bytes outside a trustee position, so this leg would prove nothing: %s", impostorText)
	}
	if namesEveryone(t, impostorText) {
		t.Errorf("AC#3 RED: namesEveryone still reads the two bytes WD anywhere in the text as Everyone; this object's only foreign grant is %s: %s", serviceSID, impostorText)
	}
	if !grantStandsOn118(t, impostor, serviceSID, false) {
		t.Fatalf("AC#3: the substitute principal is not standing on its own object, so the leg above measured nothing: %s", impostorText)
	}

	// Half four, the call site that made this instrument load-bearing: a
	// narrowNotice bucket list. The cleared buckets are what
	// TestSealReportsThePrincipalsItCleared asks this question about.
	got := captureNotices115(t)
	if err := SealFile(impostor); err != nil {
		t.Fatalf("SealFile: %v", err)
	}
	if len(*got) != 1 {
		t.Fatalf("AC#3: the seal emitted %d notice(s), want 1: %+v", len(*got), *got)
	}
	buckets := strings.Join(append(append([]string{}, (*got)[0].Principals...), (*got)[0].Inherited...), ",")
	if namesEveryone(t, buckets) {
		t.Errorf("AC#3 RED: a notice that cleared only %s was admitted as naming Everyone: %q", serviceSID, buckets)
	}
	if err := SealFile(wideChild); err != nil {
		t.Fatalf("SealFile: %v", err)
	}
	if len(*got) != 2 {
		t.Fatalf("AC#3: two seals, %d notices: %+v", len(*got), *got)
	}
	buckets = strings.Join(append(append([]string{}, (*got)[1].Principals...), (*got)[1].Inherited...), ",")
	if !namesEveryone(t, buckets) {
		t.Errorf("AC#3 RED: a notice that cleared Everyone was not recognised: %q", buckets)
	}
	t.Logf("AC#3 buckets read: SERVICE-only %q, Everyone %q", strings.Join((*got)[0].Principals, ","), buckets)
}
