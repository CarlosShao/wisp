package main

// Ticket 255 AC#1 - the receipt may only say 「已立即生效」 about a section this
// host actually acts on.
//
// THREE RULERS IN THIS FILE, each answering a different question:
//
//  1. TestTicket255SplitOnlyClaimsSectionsWithALiveReader - the pure judgement
//     (splitHotTier over every hot name plan() can book). Deterministic, no run.
//  2. TestTicket255HotRowRosterCoversTheRegistry /
//     TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim /
//     TestTicket255RosterStillMatchesTheActualReadSites - the roster's own honesty:
//     every hot TierRegistry row carries a verdict, every verdict cites a
//     `file.go:LINE [token]` whose token uniquely anchors one line of the real
//     source (the LINE is a reading, per 票 255 AC#6 - an insertion above it moves
//     the reading without reddening the ruler), and the set of production files
//     reading each section is exactly the set the roster names. That is the "会响"
//     half: a new hot section with no verdict, a changed or non-unique evidence
//     token, or a reader appearing where the receipt said "no reader" each redden
//     here.
//  3. TestTicket255ReceiptOmitsPanelFromTheImmediateSentence /
//     TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere /
//     TestTicket255ReceiptStillNamesTheLiveReadSection - the stdout readings the
//     ticket asks for, taken through the injection seam AGENTS.md §1.3 names
//     (`wisp run` with a scripted answer stream, the harness ticket 223 already
//     built, so the sentence measured is the production one and not a re-typed
//     copy). The middle one is [panel]'s control re-pointed at [session] once
//     票 255 AC#4 gave [panel] a reader in another process.
//  4. TestTicket324RosterScanRefusesSameNamedNonConfigFieldSite /
//     TestTicket324SectionMembersAreReadOffTheSchema - 票 324's controls for the
//     judgement upgrade in readsFieldAccess. The refusal half is the AC#1反控: a
//     fixture struct that merely OWNS a field named like a section (built under a
//     temp dir, so no production file has to carry a fake field) must not be counted
//     as a reader. The second half is what keeps that refusal from being blindness:
//     a real schema member still counts, and the member set must be a non-empty
//     reading of internal/config's schema.
//
// The positive control is 票 255's own words: 手改 `[panel] width` ⇒ the sentence
// "这些段已立即生效" must not contain `[panel]`. The negative control is a section
// with a real live reader ([llm], read per call by the settings leg) STILL being
// named by that sentence - plus a reading off that live leg proving the panel gets
// the new value while the boot snapshot still holds the old one.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and dies at
// load (0xc0000135) unless third_party/sherpa-onnx is on PATH.

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
)

// rosterScanSkipFiles names the one production file excluded from the reader scan:
// this ticket's own roster, whose verdict strings QUOTE code (the evidence tokens)
// without reading any config. The exclusion is by exact file name so it cannot be
// widened to hide a real reader.
var rosterScanSkipFiles = map[string]bool{
	"config_readers_255.go": true,
}

func repoRootForRoster(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// thisFile = <repo>/cmd/wisp/config_receipt_255_test.go
	return filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
}

func TestTicket255SplitOnlyClaimsSectionsWithALiveReader(t *testing.T) {
	hot := []string{
		"ball", "session", "audio", "llm", "agent", "privacy", "memory",
		"panel", "cost", "models", "observe", "hotkey", "app", "voice",
	}
	split := splitHotTier(hot)

	// Ticket 233 re-adjudicated the [agent] row: execute() reads the live
	// Manager per task (cmd/wisp/run.go:1263 [cfg = rt.mgr.Config()]), so the
	// claimable set is [llm agent] - llm via the settings leg, agent via the
	// composition root's own per-task read. Slice order fixes the join.
	if strings.Join(split.claimable, " ") != "llm agent" {
		t.Errorf("claimable = %v, want exactly [llm agent]: llm is read live per call by the settings leg (cmd/wisp/panel_config_store.go:92), agent by execute() per task since ticket 233 (cmd/wisp/run.go:1263). A section with no live reader must never be in this list.", split.claimable)
	}
	quiet := map[string]bool{}
	for _, name := range split.quiet {
		quiet[name] = true
	}
	for _, name := range hot {
		if name == "llm" || name == "agent" {
			continue
		}
		if !quiet[name] {
			t.Errorf("%q is hot-applied but appears in neither bucket - it would vanish from the receipt (摘尺当修尺, forbidden by AC#1); got claimable=%v quiet=%v", name, split.claimable, split.quiet)
		}
	}
	if len(quiet) != len(hot)-2 {
		t.Errorf("quiet = %v, want the %d non-llm/agent hot names", split.quiet, len(hot)-2)
	}

	// A name plan() booked that the registry cannot explain is neither claimed nor
	// quietly dropped: it is reported as a disagreement.
	odd := splitHotTier([]string{"dumb"})
	if len(odd.disagree) != 1 || odd.disagree[0] != "dumb" {
		t.Errorf("an unregistered hot name must land in disagree, got %+v", odd)
	}
	if !strings.Contains(odd.verdictFor("dumb"), "no hot row") {
		t.Errorf("the disagreement verdict must say why, got %q", odd.verdictFor("dumb"))
	}
	if !strings.Contains(split.verdictFor("app"), "app.theme") {
		t.Errorf("a per-key section's verdict must name the key rows behind it, got %q",
			split.verdictFor("app"))
	}
	if !strings.Contains(split.verdictFor("voice"), "voice.punctuation") {
		t.Errorf("the [voice] verdict must list its hot key rows, got %q", split.verdictFor("voice"))
	}
}

// TestTicket255HotRowRosterCoversTheRegistry is the roster's completeness ruler, in
// the same shape as internal/config's lockedKeyDisposition pin: the registry grows a
// hot row -> the roster owes a verdict; the roster keeps a row the registry no
// longer calls hot -> red.
func TestTicket255HotRowRosterCoversTheRegistry(t *testing.T) {
	hotRows := map[string]bool{}
	for path, tier := range config.TierRegistry {
		if tier == "hot" {
			hotRows[path] = true
		}
	}
	for row := range hotRows {
		verdict, ok := hotRowClaims[row]
		if !ok {
			t.Errorf("hot row %q has no verdict in hotRowClaims (cmd/wisp/config_readers_255.go): ticket 255 AC#1 forbids claiming 已立即生效 without one, and forbids silence too", row)
			continue
		}
		if !hasKnownClaimPrefix(verdict) {
			t.Errorf("hot row %q verdict %q opens with none of the five claim prefixes", row, verdict)
		}
		if !hasRosterEvidence(verdict) {
			t.Errorf("hot row %q verdict %q carries neither a `file.go:LINE [token]` cite nor the 扫描零命中 marker: an AC#1 verdict must say which reading backs it", row, verdict)
		}
		if strings.HasPrefix(verdict, hotClaimConsumed) && len(evidenceCite.FindAllString(verdict, -1)) == 0 {
			t.Errorf("hot row %q is claimed as 已立即生效 on the strength of a live reader, so it must cite that reader's file:line: %q", row, verdict)
		}
		if !strings.HasPrefix(verdict, hotClaimConsumed) && strings.Contains(verdict, "已立即生效") {
			t.Errorf("hot row %q is not consumed yet its verdict says 已立即生效: %q", row, verdict)
		}
	}
	for row := range hotRowClaims {
		if tier, ok := config.TierRegistry[row]; !ok || tier != "hot" {
			t.Errorf("hotRowClaims carries %q but TierRegistry says tier=%q ok=%v: the roster drifted from the tier table", row, tier, ok)
		}
	}
}

func hasKnownClaimPrefix(v string) bool {
	for _, p := range []string{
		hotClaimConsumed, hotClaimSnapshotOnly, hotClaimNoReader,
		hotClaimOtherProcess, hotClaimDebugHostOnly,
	} {
		if strings.HasPrefix(v, p) {
			return true
		}
	}
	return false
}

// zeroScanMarker is the second kind of AC#1 evidence: a verdict that rests on the
// reader scan finding nothing for that field, which
// TestTicket255RosterStillMatchesTheActualReadSites re-measures on every run. A
// reader-less row has no line to cite, so this marker is what keeps its verdict from
// being an assertion of faith.
const zeroScanMarker = "扫描零命中"

func hasRosterEvidence(v string) bool {
	return len(evidenceCite.FindAllString(v, -1)) > 0 || strings.Contains(v, zeroScanMarker)
}

// TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim turns every cite in the
// roster into a checked reading instead of prose: the file must exist, and the
// quoted token must occur EXACTLY ONCE in it - the token's unique hit is the
// criterion (票 255 AC#6, 10-10: 「token 唯一命中当判据，行号只当读数」). The
// recorded :LINE is a reading, so an insertion above it moves the line but never
// reddens this ruler (before AC#6 it did: 313-r2 inserted two lines, the llm row
// moved :96 -> :98, and this ruler called that "the evidence drifted" - punishing
// an edit that never touched the reader). What reddens it is the token: gone ->
// the evidence changed; not unique -> the anchor cannot pin one line and has to
// be sharpened ("SizePx" alone hit 13 lines in internal/ball/ball_windows.go).
// The error names the row, the cite, and what the recorded line says now, so
// re-adjudicating is a reading, not a guess.
func TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim(t *testing.T) {
	root := repoRootForRoster(t)
	citesChecked := 0
	for row, verdict := range hotRowClaims {
		found := evidenceCite.FindAllStringSubmatch(verdict, -1)
		if len(found) == 0 {
			if !strings.Contains(verdict, zeroScanMarker) {
				t.Errorf("row %q carries neither a cite nor the 扫描零命中 marker: %q", row, verdict)
			}
			continue
		}
		for _, c := range found {
			rel, lineNo, token := c[1], c[2], c[3]
			raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
			if err != nil {
				t.Errorf("row %q cites %s which cannot be read: %v", row, rel, err)
				continue
			}
			lines := strings.Split(string(raw), "\n")
			var hits []int
			for i, line := range lines {
				if strings.Contains(line, token) {
					hits = append(hits, i+1)
				}
			}
			switch len(hits) {
			case 1:
				if n, err := strconv.Atoi(lineNo); err == nil && n != hits[0] {
					// The reading drifted - allowed, and worth recording so the
					// next re-adjudication can refresh it (AC#6). The token is
					// what this ruler asserts; the line number never is.
					t.Logf("row %q: recorded reading %s:%s drifted, the token %q now uniquely sits at line %d", row, rel, lineNo, token, hits[0])
				}
			case 0:
				t.Errorf("row %q cites %s:%s for %q and the token is now nowhere in that file (the recorded line reads %q) - the evidence changed, so re-adjudicate this row", row, rel, lineNo, token, recordedLineReading(lines, lineNo))
			default:
				t.Errorf("row %q cites %s:%s for %q, but that token is not unique in the file (it hits lines %v): a content anchor must pin exactly one line, so sharpen the token and re-adjudicate this row", row, rel, lineNo, token, hits)
			}
			citesChecked++
		}
	}
	// Today the roster carries cites on 13 of the 17 rows (the three marker-only
	// rows [session]/[memory]/[cost] and [app.theme] carry none; [panel] carries
	// two and the four [voice] key rows share the section reader's cite). The
	// floor keeps a refactor from quietly turning all the cites into markers.
	const citedRowsFloor = 8
	if citesChecked < citedRowsFloor {
		t.Errorf("the roster carries %d checked cites, want at least %d: verdicts are drifting from readings to prose", citesChecked, citedRowsFloor)
	}
}

// recordedLineReading returns what the cite's recorded :LINE currently says, for
// the error message - the reading is quoted so a re-adjudicator can see the
// change without opening the file. Out-of-range readings are reported as such
// instead of panicking: the number is a reading, not an invariant (AC#6).
func recordedLineReading(lines []string, lineNo string) string {
	n, err := strconv.Atoi(lineNo)
	if err != nil || n < 1 || n > len(lines) {
		return "<the recorded reading is out of range>"
	}
	return strings.TrimSpace(lines[n-1])
}

// TestTicket255RosterStillMatchesTheActualReadSites is the scan half: for every
// Config section the roster adjudicates, the production files that read it must be
// exactly the files sectionReadSites names. A reader appearing where the receipt
// said "no reader" is the AC#1 lie coming back, so it reddens instead of shipping.
func TestTicket255RosterStillMatchesTheActualReadSites(t *testing.T) {
	root := repoRootForRoster(t)
	fields := configSectionsForScan()
	if len(fields) < 14 {
		t.Fatalf("the roster adjudicates %d Config sections, want the 14 hot ones: the walk itself narrowed", len(fields))
	}
	// 票 324's teeth-check for the upgraded judgement, kept next to the scan it gates:
	// every scanned section must resolve to a non-trivial member set off the schema.
	// Without this floor, a scanReadFiles that reads NO member names would refuse every
	// line and this ruler would report "no readers anywhere" as a clean reading - the
	// same 摘尺-as-fix shape AC#1 of 票 255 forbids.
	for _, field := range fields {
		if n := len(configSectionMembers(t, field)); n < 2 {
			t.Errorf("section field %q resolves to %d schema member names, want at least 2: the judgement in readsFieldAccess has nothing to match against, so it can only refuse (摘尺当修尺)", field, n)
		}
	}
	// 票 324 AC#3 wants the sections named one by one, so the loop counts its own rows
	// and logs each one's reading under -v. Reporting only: the judgement below is
	// 票 255's, unchanged.
	sectionLogged := 0
	for _, field := range fields {
		section := sectionNameForField(t, field)
		sectionLogged++
		want, ok := sectionReadSites[field]
		if !ok {
			t.Errorf("section %q (field %s) has roster claims but no sectionReadSites entry - the scan cannot tell a reader-less section from an unmapped one", section, field)
			continue
		}
		got := scanReadFiles(t, root, field)
		t.Logf("SECTION %d/%d section=%q field=%s roster=%v scan=%v match=%v", sectionLogged, len(fields), section, field, want, got, slicesEqualSorted(got, want))
		if !slicesEqualSorted(got, want) {
			t.Errorf("section %q (field %s) is read by %v in production, the roster names %v: something moved and cmd/wisp/config_readers_255.go has not been re-adjudicated",
				section, field, got, want)
		}
	}
	for field := range sectionReadSites {
		found := false
		for _, f := range fields {
			if f == field {
				found = true
			}
		}
		if !found {
			t.Errorf("sectionReadSites names field %q which the roster no longer adjudicates (a hot row left hotRowClaims but its site list stayed)", field)
		}
	}
}

func sectionNameForField(t *testing.T, field string) string {
	t.Helper()
	ft := reflect.TypeOf(config.Config{})
	for i := 0; i < ft.NumField(); i++ {
		if ft.Field(i).Name == field {
			return strings.Split(ft.Field(i).Tag.Get("toml"), ",")[0]
		}
	}
	t.Fatalf("no Config field named %q", field)
	return ""
}

// scanReadFiles returns the repo-relative, forward-slashed production files that read
// this section's fields. Excluded by design: test files (they mock, they prove
// nothing about production), internal/config itself (the schema's own package has to
// compare and copy what it declares - that is not a consumer), this ticket's roster
// file (its verdict strings quote code), testdata trees, and comment lines (票 180's
// census refuses to count a doc comment that names a field).
//
// The section's member names are read off the schema ONCE per field, outside the walk,
// so every line of every file is judged against the same reading (票 324).
func scanReadFiles(t *testing.T, root, field string) []string {
	t.Helper()
	sectionMembers := configSectionMembers(t, field)
	hits := map[string]bool{}
	for _, dir := range []string{"cmd", "internal", "tools", "scripts"} {
		base := filepath.Join(root, dir)
		err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == "testdata" {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)
			if strings.HasPrefix(rel, "internal/config/") || rosterScanSkipFiles[filepath.Base(path)] {
				return nil
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			for _, line := range strings.Split(string(raw), "\n") {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "*") {
					continue
				}
				if readsFieldAccess(trimmed, field, sectionMembers) {
					hits[rel] = true
					break
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("scan %s: %v", dir, err)
		}
	}
	out := make([]string, 0, len(hits))
	for f := range hits {
		out = append(out, f)
	}
	sort.Strings(out)
	return out
}

// readsFieldAccess answers whether this line reaches INTO the section's fields, or
// takes the whole section off a config object: `.Panel.Width`, `Config().Hotkey`,
// `Config.Panel`. A bare `cfg.Panel` handed around as one value is not yet a read of
// anything (and internal/ball's HotkeyConfig has its own Panel field, which is not
// this section), and 票 180's census draws the line the same way. These are exactly
// the three patterns whose readings produced the file lists in sectionReadSites.
//
// 票 324 UPGRADED THE DOTTED PATTERN, AND THE COST IS STATED HERE BECAUSE THE TICKET
// MAKES THE SHIPPED SHAPE BE CHOSEN AND NAMED. It used to stop at "the component
// after `.Field.` starts with a capital", so any struct that happened to own a
// same-named field counted as a reader of that section: internal/panel/composer_dispatch.go's
// `d.Models.HandleModelsListRequest(ctx, req)` reads a dispatch HANDLER hanging off
// *ComposerDispatch (ticket 320 AC#4's optional route), and was counted as a reader of
// [models]. The judgement now requires that component to be a name the schema declares
// INSIDE the section (configSectionMembers): reaching for a member the section does not
// have cannot be a read of that section.
//
// What that buys, and what it does not:
//   - it is a word-level ruler, so it stays fooled at the edges in both directions. A
//     config read spelled through something the schema does not declare (a method on the
//     section, a hop the schema walk does not take) is MISSED; a non-config struct that
//     MIRRORS the schema's member names is COUNTED. The second case is live on purpose:
//     internal/panel/config_handlers.go's `v.Voice.*` lines read panel.SettingsView, a
//     mirror of [voice], and 票 313 adjudicated that render leg into
//     sectionReadSites["Voice"] - this ticket's write face does not include
//     re-adjudicating another section's rows, so the field-name shape is still what the
//     ruler measures there. The receiver-typed alternative (票 324's 甲) would drop that
//     row and owe a roster edit outside this ticket: named as the residual, not quietly
//     paid for with a green.
//   - it is not a name whitelist and not a skip list. Nothing here names a file, so
//     there is nothing for a later wave to add to "make the red go away": the only way
//     to let a new follow-name count is to give internal/config's section struct that
//     field, which is a schema change another ticket owns.
func readsFieldAccess(line, field string, sectionMembers map[string]bool) bool {
	dotted := "." + field + "."
	for i := 0; i+len(dotted) < len(line); i++ {
		if line[i:i+len(dotted)] != dotted {
			continue
		}
		rest := line[i+len(dotted):]
		if rest[0] < 'A' || rest[0] > 'Z' {
			continue
		}
		if sectionMembers[identName(rest)] == false {
			continue
		}
		return true
	}
	for _, prefix := range []string{"Config()." + field, "Config." + field} {
		idx := strings.Index(line, prefix)
		if idx < 0 {
			continue
		}
		end := idx + len(prefix)
		if end >= len(line) || !isIdentByte(line[end]) {
			return true
		}
	}
	return false
}

// identName returns the identifier at the head of s (empty when s does not start with
// one), so readsFieldAccess can ask the schema about the component that follows
// `.Field.` instead of asking it about the first letter.
func identName(s string) string {
	end := 0
	for end < len(s) && isIdentByte(s[end]) {
		end++
	}
	return s[:end]
}

// configSectionMembers reads the section's own member names off internal/config's
// schema instead of typing them here: the exported fields of the struct Config
// declares for this field. It is a READING, so the roster cannot drift from the schema
// without reddening (TestTicket324SectionMembersAreReadOffTheSchema), and a struct of
// some other package that owns a same-named field cannot widen it.
func configSectionMembers(t *testing.T, field string) map[string]bool {
	t.Helper()
	out := map[string]bool{}
	ft := reflect.TypeOf(config.Config{})
	for i := 0; i < ft.NumField(); i++ {
		sf := ft.Field(i)
		if sf.Name != field {
			continue
		}
		st := sf.Type
		for st.Kind() == reflect.Ptr {
			st = st.Elem()
		}
		if st.Kind() != reflect.Struct {
			return out
		}
		for j := 0; j < st.NumField(); j++ {
			if mf := st.Field(j); mf.IsExported() {
				out[mf.Name] = true
			}
		}
	}
	return out
}

func isIdentByte(b byte) bool {
	return b == '_' || (b >= '0' && b <= '9') || (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func slicesEqualSorted(a, b []string) bool {
	x := append([]string(nil), a...)
	y := append([]string(nil), b...)
	sort.Strings(x)
	sort.Strings(y)
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

// --- 票 324: the two controls for the upgraded reader judgement ------------------

// The fake site is the shape 票 324 measured in production: a struct that OWNS a
// field named like a config section, and reaches a member that section's schema does
// not declare (a handler, not a value). It lives in a fixture the scan reads, not in a
// production file - the ticket forbids planting a fake field in the shipped tree.
const rosterScanFakeSite324 = `package panel

type samenameHandler interface{ HandleModelsListRequest() error }

type samenameDispatch struct{ Models samenameHandler }

func (d *samenameDispatch) list() error { return d.Models.HandleModelsListRequest() }
`

// The schema site is the counterpart: the member the [models] section really declares.
// Both fixtures are needed - a ruler that refused EVERYTHING would pass a refusal-only
// control, so the counted half is what makes the refusal half a reading.
const rosterScanSchemaSite324 = `package main

func readMirror(cfg struct{ Models struct{ Mirror []string } }) []string {
	return cfg.Models.Mirror
}
`

// TestTicket324RosterScanRefusesSameNamedNonConfigFieldSite is 票 324 AC#1's 反控,
// taken through the real scan (not through a re-typed predicate): over a fixture tree,
// the handler-shaped same-named field must not be named as a reader of [models] while
// the schema member must be.
func TestTicket324RosterScanRefusesSameNamedNonConfigFieldSite(t *testing.T) {
	root := t.TempDir()
	// The scan walks cmd/internal/tools/scripts and fails the test on a missing tree,
	// so the fixture root carries all four.
	for _, dir := range []string{"cmd/wisp", "internal/panel", "tools", "scripts"} {
		if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(dir)), 0o755); err != nil {
			t.Fatalf("fixture dir %s: %v", dir, err)
		}
	}
	fixtures := map[string]string{
		"internal/panel/samename_fixture.go": rosterScanFakeSite324,
		"cmd/wisp/schembsite_fixture.go":     rosterScanSchemaSite324,
	}
	for rel, src := range fixtures {
		if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(rel)), []byte(src), 0o644); err != nil {
			t.Fatalf("write fixture %s: %v", rel, err)
		}
	}
	got := scanReadFiles(t, root, "Models")
	want := []string{"cmd/wisp/schembsite_fixture.go"}
	// The reading itself, logged so the 反控 is quotable as a measurement and not just
	// as a PASS line (go test -v). fakeSiteRefused is true exactly when the scan named
	// the schema-member fixture and nothing else.
	t.Logf("FIXTURE root=%s scan=%v want=%v fakeSiteRefused=%v", root, got, want, slicesEqualSorted(got, want))
	if !slicesEqualSorted(got, want) {
		t.Errorf("the scan read %v for section field Models over the two fixtures, want %v: d.Models.HandleModelsListRequest in internal/panel/samename_fixture.go is 票 324's same-named NON-config field shape (a dispatch handler, not a value the [models] section declares) and must not be counted as a config reader, while cfg.Models.Mirror must be - a scan that named both has no judgement, and a scan that named neither is 摘尺", got, want)
	}
}

// TestTicket324SectionMembersAreReadOffTheSchema pins the judgement on the two
// production lines 票 324 read, as data rather than as a walk: the [models] schema
// member keeps counting, the handler-shaped same name stops counting, and the
// section-level `Config().Hotkey` shape (no member after it) still counts. It is the
// named case that reddens when the schema-member requirement is degraded to "any
// capitalised component" (AC#1's 正控), and the one that reddens if the member reading
// is broken so hard the ruler refuses everything.
func TestTicket324SectionMembersAreReadOffTheSchema(t *testing.T) {
	models := configSectionMembers(t, "Models")
	if !models["Mirror"] {
		t.Errorf("the [models] schema reading lost its Mirror member (got %v): the real reader line would stop counting, which is 摘尺当修尺", models)
	}
	if models["HandleModelsListRequest"] {
		t.Errorf("the [models] schema reading carries HandleModelsListRequest, which is panel.ModelsListRequestHandler's method and not a config member: the ruler can no longer tell the two apart (got %v)", models)
	}
	schemaSite := "mirrors = cfg.Models.Mirror"
	if !readsFieldAccess(schemaSite, "Models", models) {
		t.Errorf("cmd/wisp/models.go's real read line %q is no longer judged a read of [models]: the ruler went blind instead of going strict", schemaSite)
	}
	fakeSite := "return d.Models.HandleModelsListRequest(ctx, req)"
	if readsFieldAccess(fakeSite, "Models", models) {
		t.Errorf("internal/panel/composer_dispatch.go's handler line %q is judged a read of [models] again: the schema-member evidence was dropped, so a same-named field of any other struct counts as a config reader", fakeSite)
	}
	hotkey := configSectionMembers(t, "Hotkey")
	if !readsFieldAccess("h := mgr.Config().Hotkey", "Hotkey", hotkey) {
		t.Errorf("the Config()-rooted section-level line stopped counting: a whole section taken off the live Manager is a read of it")
	}
	// The upgrade must be a judgement about names the schema declares, not a general
	// refusal: [llm] carries Providers, and the reader 票 324 added to the roster keys
	// on exactly that member.
	llm := configSectionMembers(t, "LLM")
	if !llm["Providers"] || !readsFieldAccess("p, ok := cfg.LLM.Providers[provider]", "LLM", llm) {
		t.Errorf("cmd/wisp/panel_models_320.go's cfg.LLM.Providers line must count as a read of [llm]; schema members = %v", llm)
	}
}

// --- the two stdout readings, through the `wisp run` seam ----------------------

// TestTicket255ReceiptOmitsPanelFromTheImmediateSentence is AC#1's original
// positive control, read off a live process: hand-edit [panel] width while
// `wisp run` is up, and the sentence that says 「已立即生效」 must not name
// [panel].
//
// WHAT 票 255 AC#4 CHANGED HERE, and what it did not. The judgement still holds,
// for a reason that got more specific rather than weaker: [panel] now HAS a
// production reader - the resident panel host's assembly root re-reads it at every
// window creation - but that reader lives in the `wisp` process, and the host
// measured here (`wisp run`) builds no panel window at all. So this seam still may
// not say 已立即生效 about [panel], and the honest line still has to name the edit.
// The strongest form of "no reader anywhere" moved to
// TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere, which samples [session].
func TestTicket255ReceiptOmitsPanelFromTheImmediateSentence(t *testing.T) {
	const plantedWidth = 641 // the schema default is 640 (schema.go's PanelSection.Width tag)
	r := newReloadRun223(t, 40*time.Second)
	r.live(t, func() {
		r.awaitAudit(t, "config: HOT-RELOAD state=armed")
		if got := r.rt.mgr.Config().Panel.Width; got == plantedWidth {
			t.Fatalf("this case's premise moved: the loaded config already carries width %d", plantedWidth)
		}
		r.plant(t, "[fs]", "[panel]\nwidth = 641\n\n[fs]")

		applied := r.awaitAudit(t, "config: HOT-RELOAD state=applied")
		if !strings.Contains(applied, "hot=[panel]") {
			t.Errorf("the reload did not book panel as hot-tier, so this proves nothing: %q", applied)
		}
		out := r.awaitStdout(t, "值已换进本进程内存")

		// THE JUDGEMENT: no line that says 「已立即生效」 may name [panel].
		for _, line := range strings.Split(out, "\n") {
			if strings.Contains(line, "这些段已立即生效") && strings.Contains(line, "[panel]") {
				t.Errorf("ticket 255 AC#1: the receipt claimed a reader-less section took effect immediately: %q", line)
			}
		}
		honest := lastLineWith(out, "值已换进本进程内存")
		if !strings.Contains(honest, "[panel]") {
			t.Errorf("the honest sentence did not name [panel], so the edit vanished from the receipt: %q", honest)
		}
		for _, needle := range []string{"没有会按新值做事的读者", "本次运行不会因此改变行为"} {
			if !strings.Contains(honest, needle) {
				t.Errorf("the honest sentence must say %q: %q", needle, honest)
			}
		}
		// Not a claim that the reload did nothing: the memory did move.
		if got := r.rt.mgr.Config().Panel.Width; got != plantedWidth {
			t.Errorf("live config panel.width = %d after the edit, want %d (the tier did not apply, so the honest line's premise is wrong too)", got, plantedWidth)
		}
		// The ledger carries the per-section verdict, so "who reads it" is a
		// reading off the trail and not an inference from the console.
		line := r.awaitAuditLine(t, "config: HOT-RELOAD-READER section=panel")
		if !strings.Contains(line, "other-process") {
			t.Errorf("the ledger does not name panel's reader as living in another process: %q", line)
		}
		// Re-adjudicated in the 票 255 AC#4 batch. This row used to demand
		// "no-reader" plus a cite to the hard-coded geometry literal
		// (`Width: 420` / `Height: 260` in panel_host_windows.go's create block,
		// at HEAD 67ab595d), and AC#4 removed both: the literal is gone from the
		// create block and the read site moved into the resident assembly root.
		// The cite below is content-anchored and re-read from the file by
		// TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim - per 票 255 AC#6
		// the assertion keys on the file and the token, never on the :LINE reading
		// (a recorded number here would re-import the rot the roster just shed).
		if !strings.Contains(line, "panel_resident_windows.go") || !strings.Contains(line, "[cfg.Panel.Width]") {
			t.Errorf("the ledger's verdict for panel must cite the per-create read site AC#4 installed (file + content anchor): %q", line)
		}
		if n := r.rt.windowCount(); n != 0 {
			t.Errorf("a hot edit with no reader in THIS process displayed %d cards, want 0", n)
		}
	})
}

// TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere is AC#1's positive control
// in the shape the ticket asked for, moved off [panel] by AC#4 and off [session]
// by ticket 233.
//
// WHY IT MOVED HERE. Until 票 255 AC#4 landed, [panel] was the section that proved
// the sentence can lie: its width was hot-applied into memory and read by nobody
// anywhere. That is no longer true of it - the resident panel host re-reads it at
// every window creation - so leaving the strongest form of the control on [panel]
// would measure "no reader in this process" while claiming "no reader at all",
// which is the same word doing two jobs AC#1 exists to separate. [session] took
// the sample over then, and ticket 233 moved it once more: ticket wire-r1's
// resident voice leg reads [session] warm_timeout_sec at its boot LoadFile
// (cmd/wisp/resident_voice_windows.go:1329), so "no reader anywhere" stopped
// being true of [session] too - the roster row was re-adjudicated to
// other-process in the same batch. [ball] takes the sample: internal/config
// tiers.go registers it hot, its size field is hot-applied by plan(), and the
// reader scan re-measures zero production readers for it on every run (the
// roster row's own verdict and TestTicket255RosterStillMatchesTheActualReadSites
// both say so).
func TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere(t *testing.T) {
	r := newReloadRun223(t, 40*time.Second)
	r.live(t, func() {
		r.awaitAudit(t, "config: HOT-RELOAD state=armed")
		r.plant(t, "[fs]", "[ball]\nsize = 66\n\n[fs]")

		applied := r.awaitAudit(t, "config: HOT-RELOAD state=applied")
		if !strings.Contains(applied, "hot=[ball]") {
			t.Errorf("the reload did not book ball as hot-tier, so this proves nothing: %q", applied)
		}
		out := r.awaitStdout(t, "值已换进本进程内存")
		for _, line := range strings.Split(out, "\n") {
			if strings.Contains(line, "这些段已立即生效") && strings.Contains(line, "[ball]") {
				t.Errorf("ticket 255 AC#1: the receipt claimed a section with no reader anywhere took effect immediately: %q", line)
			}
		}
		honest := lastLineWith(out, "值已换进本进程内存")
		if !strings.Contains(honest, "[ball]") {
			t.Errorf("the honest sentence did not name [ball], so the edit vanished from the receipt: %q", honest)
		}
		line := r.awaitAuditLine(t, "config: HOT-RELOAD-READER section=ball")
		if !strings.Contains(line, "no-reader") {
			t.Errorf("[ball] is the reader-less sample now, so its ledger row must still say no-reader: %q", line)
		}
		if strings.Contains(line, "panel") {
			t.Errorf("the [ball] ledger row names panel, so the two verdicts are not separable: %q", line)
		}
	})
}

// TestTicket255ReceiptStillNamesTheLiveReadSection is AC#1's negative control: a
// section this host really does read live must STILL be named by 「已立即生效」, or
// the fix would be 摘尺. [llm] is that section - the settings leg answers the panel
// out of the live Manager (cmd/wisp/panel_config_store.go:92), so the case reads
// that leg back while the run is alive: the new value is what the panel gets, and
// the boot snapshot still holds the old one. That is what makes the claim "有人按新
// 值做事" instead of "内存换了".
func TestTicket255ReceiptStillNamesTheLiveReadSection(t *testing.T) {
	const oldRef = `api_key_ref = "dpapi:acme"`
	const newRef = `api_key_ref = "dpapi:acme-255r2"`
	r := newReloadRun223(t, 40*time.Second)
	r.live(t, func() {
		r.awaitAudit(t, "config: HOT-RELOAD state=armed")
		if got := r.rt.cfg.LLM.Providers["acme"].APIKeyRef; got != "dpapi:acme" {
			t.Fatalf("this case's premise moved: the boot snapshot carries ref %q", got)
		}
		r.plant(t, oldRef, newRef)

		applied := r.awaitAudit(t, "config: HOT-RELOAD state=applied")
		if !strings.Contains(applied, "hot=[llm]") {
			t.Errorf("the reload did not book llm as hot-tier: %q", applied)
		}
		out := r.awaitStdout(t, "这些段已立即生效")
		// [llm] must be claimed, and it must NOT be in the honest bucket. Since
		// ticket 233 [agent] is claimed too (execute() reads the live Manager per
		// task), so the honest bucket must not carry it either.
		line := lastLineWith(out, "这些段已立即生效")
		if !strings.Contains(line, "[llm]") {
			t.Errorf("the live-read section is missing from the immediate sentence (摘尺): %q", line)
		}
		if strings.Contains(line, "[panel]") || strings.Contains(line, "[ball]") {
			t.Errorf("the immediate sentence names a reader-less section: %q", line)
		}
		if h := lastLineWith(out, "值已换进本进程内存"); strings.Contains(h, "[llm]") || strings.Contains(h, "[agent]") {
			t.Errorf("a claimed section was ALSO put in the honest bucket, which contradicts the claim: %q", h)
		}
		// The reader itself, read off the live leg while the process is up.
		view, err := r.rt.settings.ReadSettings(context.Background())
		if err != nil {
			t.Fatalf("the settings leg must answer from the live Manager: %v", err)
		}
		if len(view.Providers) != 1 {
			t.Fatalf("settings view carries %d providers, want the one this run configured", len(view.Providers))
		}
		if view.Providers[0].APIKeyRef != "dpapi:acme-255r2" {
			t.Errorf("the live settings read still answers ref %q, want the planted %q - if this fails, [llm] has no live reader and the receipt must not claim it",
				view.Providers[0].APIKeyRef, "dpapi:acme-255r2")
		}
		if got := r.rt.cfg.LLM.Providers["acme"].APIKeyRef; got == "dpapi:acme-255r2" {
			t.Errorf("the boot snapshot moved too, so nothing here separates a live reader from a memory write: %q", got)
		}
		ledger := r.awaitAuditLine(t, "config: HOT-RELOAD-READER section=llm")
		if !strings.Contains(ledger, "consumed:") {
			t.Errorf("the ledger must record why llm earned the claim: %q", ledger)
		}
		if n := r.rt.windowCount(); n != 0 {
			t.Errorf("a hot-tier llm edit displayed %d cards, want 0", n)
		}
	})
}

// TestTicket255ReceiptSentenceAssemblyIsFiltered pins the sentence-building half
// without the tick's timing: it hands reportReload a Report carrying every hot name
// plan() can book, and demands that the 「已立即生效」 line list only the live-read
// section while the honest line carries the rest. The two cases above measure the
// plant-driven production reading; this one exists so a regression that puts
// rep.Hot back into the claim (the pre-255 shape) reddens on the judgement line
// itself, not on a missing-sentence deadline. The Report is the production type and
// reportReload is the production function - nothing here re-implements the sentence.
func TestTicket255ReceiptSentenceAssemblyIsFiltered(t *testing.T) {
	r := newReloadRun223(t, 40*time.Second)
	r.live(t, func() {
		r.awaitAudit(t, "config: HOT-RELOAD state=armed")
		r.rt.reportReload(&config.Report{
			Hot: []string{
				"ball", "session", "audio", "llm", "agent", "privacy", "memory",
				"panel", "cost", "models", "observe", "hotkey", "app", "voice",
			},
		})
		out := r.awaitStdout(t, "值已换进本进程内存")
		immediate := lastLineWith(out, "这些段已立即生效")
		if immediate == "" {
			t.Fatalf("the immediate-tier sentence disappeared - AC#1 forbids 摘尺当修尺:\n%s", out)
		}
		// Ticket 233: [agent] joined the claimable set (execute() reads the live
		// Manager per task), so the immediate sentence carries it next to [llm].
		if !strings.Contains(immediate, "[llm]") || !strings.Contains(immediate, "[agent]") {
			t.Errorf("the live-read sections are not claimed: %q", immediate)
		}
		for _, forbidden := range []string{"[panel]", "[ball]", "[app]", "[voice]", "[models]", "[hotkey]"} {
			if strings.Contains(immediate, forbidden) {
				t.Errorf("ticket 255 AC#1: the 已立即生效 sentence claims %s, which has no live reader in this host: %q", forbidden, immediate)
			}
		}
		honest := lastLineWith(out, "值已换进本进程内存")
		for _, want := range []string{"[panel]", "[ball]", "[app]", "[voice]", "[models]", "[hotkey]"} {
			if !strings.Contains(honest, want) {
				t.Errorf("the honest sentence omits %s, so the edit vanished from the receipt: %q", want, honest)
			}
		}
		if strings.Contains(honest, "[llm]") || strings.Contains(honest, "[agent]") {
			t.Errorf("a claimed section is in both buckets: %q", honest)
		}
		// Logged for the evidence file (run with -v) so both sentences are quoted
		// off the product, not retyped from the format strings.
		t.Logf("SENTENCES\n  IMMEDIATE %q\n  HONEST    %q", immediate, honest)
		// Every hot name owes a ledger line naming its rows and verdicts.
		for _, name := range []string{"panel", "llm", "app", "voice"} {
			line := r.awaitAuditLine(t, "config: HOT-RELOAD-READER section="+name)
			if !strings.Contains(line, "tier_row=") {
				t.Errorf("the ledger line for %s does not name the tier rows behind the claim: %q", name, line)
			}
			// Logged so the evidence file quotes the production ledger line
			// instead of re-typing it (go test -v).
			t.Logf("LEDGER %s", line)
		}
		perKey := r.awaitAuditLine(t, "config: HOT-RELOAD-READER section=voice")
		if !strings.Contains(perKey, "voice.tts.speed") {
			t.Errorf("[voice] is a per-key section, so its ledger line must list the key rows: %q", perKey)
		}
	})
}

func lastLineWith(haystack, needle string) string {
	lines := strings.Split(haystack, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], needle) {
			return strings.TrimSpace(lines[i])
		}
	}
	return ""
}
