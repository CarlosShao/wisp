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
//     every hot TierRegistry row carries a verdict, every verdict cites
//     `file.go:LINE [token]` and that line still contains that token in the real
//     source, and the set of production files reading each section is exactly the
//     set the roster names. That is the "会响" half: a new hot section with no
//     verdict, a drifted evidence line, or a reader appearing where the receipt
//     said "no reader" each redden here.
//  3. TestTicket255ReceiptOmitsPanelFromTheImmediateSentence /
//     TestTicket255ReceiptStillNamesTheLiveReadSection - the two stdout readings the
//     ticket asks for, taken through the injection seam AGENTS.md §1.3 names
//     (`wisp run` with a scripted answer stream, the harness ticket 223 already
//     built, so the sentence measured is the production one and not a re-typed copy).
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

	if strings.Join(split.claimable, " ") != "llm" {
		t.Errorf("claimable = %v, want exactly [llm]: it is the only hot section this host reads live (cmd/wisp/panel_config_store.go:92 takes s.mgr.Config() per call). A section with no live reader must never be in this list.", split.claimable)
	}
	quiet := map[string]bool{}
	for _, name := range split.quiet {
		quiet[name] = true
	}
	for _, name := range hot {
		if name == "llm" {
			continue
		}
		if !quiet[name] {
			t.Errorf("%q is hot-applied but appears in neither bucket - it would vanish from the receipt (摘尺当修尺, forbidden by AC#1); got claimable=%v quiet=%v", name, split.claimable, split.quiet)
		}
	}
	if len(quiet) != len(hot)-1 {
		t.Errorf("quiet = %v, want the %d non-llm hot names", split.quiet, len(hot)-1)
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

// TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim turns every file:line in the
// roster into a checked reading instead of prose: the file must exist, the line
// number must be in range, and the line must still carry the quoted token. This is
// what "逐名给出谁读它的 file:line 证据" has to mean to be auditable.
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
			n, err := strconv.Atoi(lineNo)
			if err != nil {
				t.Errorf("row %q has an unparseable line number %q", row, lineNo)
				continue
			}
			raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
			if err != nil {
				t.Errorf("row %q cites %s which cannot be read: %v", row, rel, err)
				continue
			}
			lines := strings.Split(string(raw), "\n")
			if n < 1 || n > len(lines) {
				t.Errorf("row %q cites %s:%s but the file has %d lines", row, rel, lineNo, len(lines))
				continue
			}
			got := strings.TrimSpace(lines[n-1])
			if !strings.Contains(got, token) {
				t.Errorf("row %q cites %s:%s for %q, that line now reads %q - the evidence drifted, so re-adjudicate this row (and if the reader really moved, the receipt's claim may have moved with it)",
					row, rel, lineNo, token, got)
			}
			citesChecked++
		}
	}
	// Today eight of the seventeen rows back their verdict with a real line (the
	// consumed [llm] row, the three near-miss readers [agent]/[models]/[hotkey], and
	// the three "the value is decided elsewhere" rows [ball]/[privacy]/[observe] plus
	// [panel], whose row cites the hard-coded geometry). The floor keeps a refactor
	// from quietly turning all the cites into markers.
	const citedRowsFloor = 8
	if citesChecked < citedRowsFloor {
		t.Errorf("the roster carries %d checked cites, want at least %d: verdicts are drifting from readings to prose", citesChecked, citedRowsFloor)
	}
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
	for _, field := range fields {
		section := sectionNameForField(t, field)
		want, ok := sectionReadSites[field]
		if !ok {
			t.Errorf("section %q (field %s) has roster claims but no sectionReadSites entry - the scan cannot tell a reader-less section from an unmapped one", section, field)
			continue
		}
		got := scanReadFiles(t, root, field)
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
func scanReadFiles(t *testing.T, root, field string) []string {
	t.Helper()
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
				if readsFieldAccess(trimmed, field) {
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
func readsFieldAccess(line, field string) bool {
	dotted := "." + field + "."
	for i := 0; i+len(dotted) < len(line); i++ {
		if line[i:i+len(dotted)] != dotted {
			continue
		}
		if next := line[i+len(dotted)]; next >= 'A' && next <= 'Z' {
			return true
		}
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

// --- the two stdout readings, through the `wisp run` seam ----------------------

// TestTicket255ReceiptOmitsPanelFromTheImmediateSentence is AC#1's positive control,
// read off a live process: hand-edit [panel] width while `wisp run` is up, and the
// sentence that says 「已立即生效」 must not name [panel]. The edit is not nothing -
// plan() really overwrote the memory, and the receipt says that in its own honest
// line - but nobody acts on the new width, and 票 255 AC#4 (not this ticket, not
// this leg) is what would change that.
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
		if !strings.Contains(line, "no-reader") {
			t.Errorf("the ledger does not name panel as reader-less: %q", line)
		}
		if !strings.Contains(line, "panel_host_windows.go:304") {
			t.Errorf("the ledger's verdict for panel must cite the hard-coded geometry: %q", line)
		}
		if n := r.rt.windowCount(); n != 0 {
			t.Errorf("a reader-less hot edit displayed %d cards, want 0", n)
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
		line := lastLineWith(out, "这些段已立即生效")
		if !strings.Contains(line, "[llm]") {
			t.Errorf("the live-read section is missing from the immediate sentence (摘尺): %q", line)
		}
		if strings.Contains(line, "[panel]") || strings.Contains(line, "[ball]") {
			t.Errorf("the immediate sentence names a reader-less section: %q", line)
		}
		if h := lastLineWith(out, "值已换进本进程内存"); strings.Contains(h, "[llm]") {
			t.Errorf("llm was ALSO put in the honest bucket, which contradicts the claim: %q", h)
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
			Hot: []string{"ball", "session", "audio", "llm", "agent", "privacy", "memory",
				"panel", "cost", "models", "observe", "hotkey", "app", "voice"},
		})
		out := r.awaitStdout(t, "值已换进本进程内存")
		immediate := lastLineWith(out, "这些段已立即生效")
		if immediate == "" {
			t.Fatalf("the immediate-tier sentence disappeared - AC#1 forbids 摘尺当修尺:\n%s", out)
		}
		if !strings.Contains(immediate, "[llm]") {
			t.Errorf("the live-read section is not claimed: %q", immediate)
		}
		for _, forbidden := range []string{"[panel]", "[ball]", "[app]", "[voice]", "[agent]", "[models]", "[hotkey]"} {
			if strings.Contains(immediate, forbidden) {
				t.Errorf("ticket 255 AC#1: the 已立即生效 sentence claims %s, which has no live reader in this host: %q", forbidden, immediate)
			}
		}
		honest := lastLineWith(out, "值已换进本进程内存")
		for _, want := range []string{"[panel]", "[ball]", "[app]", "[voice]", "[agent]", "[models]", "[hotkey]"} {
			if !strings.Contains(honest, want) {
				t.Errorf("the honest sentence omits %s, so the edit vanished from the receipt: %q", want, honest)
			}
		}
		if strings.Contains(honest, "[llm]") {
			t.Errorf("[llm] is in both buckets: %q", honest)
		}
		// Every hot name owes a ledger line naming its rows and verdicts.
		for _, name := range []string{"panel", "llm", "app", "voice"} {
			line := r.awaitAuditLine(t, "config: HOT-RELOAD-READER section="+name)
			if !strings.Contains(line, "tier_row=") {
				t.Errorf("the ledger line for %s does not name the tier rows behind the claim: %q", name, line)
			}
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
