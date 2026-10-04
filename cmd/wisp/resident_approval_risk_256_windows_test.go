//go:build windows

package main

// Ticket 256 AC#1/AC#2, the [risk] half: the four resident-leg rulers.
//
// WHAT THE LEG CHANGED. The resident process built its approval gate from three
// approval.Options fields (UI / Channels / Logf), so [risk]'s confirm_timeout_sec
// and l1_window_sec did nothing here - the gate ran the compiled 300s / 3s
// whatever config.toml said. Two fields are now passed, from a construction-time
// read of the host's own config.toml (newResidentApprovalWithConfig, called by
// runResident with rt.Layout.DataDir).
//
// THE FOUR RULERS, and which fact each one is the only one that can hold:
//
//   ① TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants
//      no host view => the compiled constants, and the provenance says so.
//   ② TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow
//      a seeded legal value moves the resident leg's own gate.
//      Its positive control is Queue().Timeout(): it is the ONLY valid positive
//      control, because Window() is clamped into [2s,3s] by gate.go's New plus
//      queue.go's MaxL1Window, so a seeded l1_window_sec of 99 or of 1 reads back
//      green either way - Window is used here as the clamp-still-holds reverse
//      control and never as proof the value was taken.
//   ③ TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly
//      the limitation, run rather than asserted in prose: approval.Gate and
//      Queue copy these two in New/NewQueue and nothing in the repository
//      re-applies them, so what this leg closes is "the gate THIS launch is
//      built with follows [risk]" and NOT "edit the file and the running process
//      follows". No ruler in this ticket may be written as the second sentence,
//      and this one is what makes that unmissable.
//   ④ TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims
//      the landing site, measured off the syntax tree rather than off words: the
//      field set this file really passes to approval.New, and the denominator
//      read off approval's own Options declaration. Comments are invisible to it
//      by construction, so writing the word "[risk]" anywhere cannot turn it
//      green. It reddens when the shipped count changes in either direction.
//   ⑤ TestTicket256ResidentBootPassesTheDataDirToTheGate
//      the other half of "did it actually land": the production call site really
//      hands over rt.Layout.DataDir instead of the no-view signature.
//
// WHAT IS NOT HERE, NAMED SO THE ABSENCE IS NOT MISREAD AS A PASS:
// Options.Grants stays unset. The session ledger is constructed only inside
// cmd/wisp/run.go's assembleRuntime, roughly 129 lines AFTER this gate is built
// and behind a conditional early return (256-a2 census §0/§2), and g.grants has
// one writer and no late-binding entry. Closing that half means a second minted
// ledger or a new holder type - new seams, not a field. It is filed as pending
// ticket 265, and ④ deliberately asserts Grants is ABSENT so nobody can quietly
// claim AC#1 as whole from this leg.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/config"
)

// risk256Dir is this package's own directory, resolved off the caller's file so
// the two syntax-tree rulers read the same files whether `go test` runs them
// from the package dir or a prebuilt test binary runs them from elsewhere. It
// returns the PACKAGE dir; the one ruler that reaches outside it (the Options
// denominator, which lives in internal/agent/approval/gate.go) takes two more
// Dir() steps, the same way repoRootForRoster takes three from runtime.Caller's
// file path. Getting that depth wrong is how a ruler ends up parsing
// cmd/internal/... and reporting a missing file as a code defect.
func risk256Dir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed for the 256 ruler")
	}
	return filepath.Dir(thisFile)
}

// writeRiskConfig256 lays down a minimal readable config.toml carrying the
// passed [risk] body. The base is the canonical short form
// resident_hotkey_258_windows_test.go:193 already uses: schema_version alone
// loads, so what is planted is exactly what is being measured.
func writeRiskConfig256(t *testing.T, dir, riskBody string) string {
	t.Helper()
	body := "schema_version = " + strconv.Itoa(config.SchemaVersionCurrent) + "\n"
	if riskBody != "" {
		body += "[risk]\n" + riskBody
	}
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("plant config.toml in %s: %v", dir, err)
	}
	return path
}

// ---------------------------------------------------------------- ① default tier

// TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants is ruler
// ①: with no host config view - either no dataDir at all (the signature the 246
// roster keeps) or a dataDir whose config.toml will not read - the gate this leg
// builds must run on approval's compiled constants, and the provenance field
// must name which of the two fallback shapes it is.
//
// The read face is Queue().Timeout() (gate.go's Queue() plus queue.go's
// Timeout()); that is the one face the orchestrator's §8.3 ruled valid. The
// paired positive control at the end is what makes this a ruler and not a
// tautology: the SAME directory, once it holds a planted timeout, must read
// differently, or "the fallback produced the default" would be true of a gate
// that ignored config entirely.
func TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants(t *testing.T) {
	shapes := []struct {
		name string
		// dir is resolved per case: "no-view" is the empty dataDir, "unreadable"
		// is a real directory with no config.toml in it.
		kind string
		prov string
	}{
		{"no host view at all", "no-view", riskProvenanceNoView},
		{"config.toml missing", "unreadable", riskProvenanceUnreadable},
	}
	for _, s := range shapes {
		t.Run(s.name, func(t *testing.T) {
			var ra *residentApproval
			dir := ""
			switch s.kind {
			case "no-view":
				ra = newResidentApproval()
			case "unreadable":
				dir = t.TempDir()
				ra = newResidentApprovalWithConfig(dir)
			}
			if ra.gate == nil {
				t.Fatal("no gate built, so nothing below measures anything")
			}
			if got := ra.gate.Queue().Timeout(); got != approval.DefaultApprovalTimeout {
				t.Errorf("Queue().Timeout() with no usable [risk] view = %v, want the compiled %v "+
					"(300s is the contract default; a fourth number here means the fallback grew a value of its own)",
					got, approval.DefaultApprovalTimeout)
			}
			if got := ra.gate.Window(); got != approval.DefaultL1Window {
				t.Errorf("Window() with no usable [risk] view = %v, want the compiled %v", got, approval.DefaultL1Window)
			}
			if ra.riskProvenance != s.prov {
				t.Errorf("riskProvenance = %q, want %q: the fallback has to be SAYABLE, not just present "+
					"(ticket 258's hotkeyProvenance* trio is the precedent for that word)", ra.riskProvenance, s.prov)
			}

			// Positive control, paired to this negative ruler in the same case:
			// plant a timeout the constants cannot produce and require it to show
			// up. If the read face were hard-wired, or if newResidentApproval were
			// silently the only path, this sub-assertion is the one that reddens.
			if s.kind == "unreadable" {
				writeRiskConfig256(t, dir, "confirm_timeout_sec = 45\n")
				seeded := newResidentApprovalWithConfig(dir)
				if got := seeded.gate.Queue().Timeout(); got != 45*time.Second {
					t.Errorf("positive control failed: seeded confirm_timeout_sec = 45 but Queue().Timeout() = %v, "+
						"so ruler ① above was measuring a face that ignores config at all", got)
				}
				if seeded.riskProvenance != riskProvenanceRead {
					t.Errorf("positive control: provenance = %q, want %q once the file reads",
						seeded.riskProvenance, riskProvenanceRead)
				}
			}
		})
	}
	// The old signature and the new one with an empty dataDir must be the same
	// shape - that is what "delegation, not a second implementation" means.
	if a, b := newResidentApproval(), newResidentApprovalWithConfig(""); a.riskProvenance != b.riskProvenance {
		t.Errorf("newResidentApproval() and newResidentApprovalWithConfig(\"\") disagree on provenance: %q vs %q",
			a.riskProvenance, b.riskProvenance)
	}
}

// ------------------------------------------------------------------ ② positive control

// TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow is ruler ②: the
// resident leg's own gate follows the planted [risk] value.
//
// Timeout is the positive control. Window is the CLAMP-STILL-HOLDS reverse
// control, stated per case, because gate.go's New clamps it into
// [MinL1Window, MaxL1Window] = [2s, 3s] and queue.go:122 pins MaxL1Window at 3s:
// a seeded 99 and a seeded 1 both read back inside the band, so treating
// Window() as proof the config was taken would be a permanently-green ruler.
// The one seeded window value that IS informative is 2s, because the compiled
// default this leg fell back to before ticket 256 was 3s.
func TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow(t *testing.T) {
	cases := []struct {
		name      string
		seed      string
		wantTime  time.Duration
		wantWin   time.Duration
		winIsWhat string // what the Window() reading is allowed to prove
	}{
		{
			name:     "timeout only - the window then comes from the schema tag, not from the gate's compiled constant",
			seed:     "confirm_timeout_sec = 45\n",
			wantTime: 45 * time.Second,
			wantWin:  2 * time.Second,
			winIsWhat: "a fact this leg had to learn by running: once a readable config.toml is consulted, an " +
				"UNseeded l1_window_sec arrives as the schema tag's 2 (internal/config/schema.go's " +
				`default:"2"` + ") and NOT as approval's compiled 3s, so the window moves to 2s even though " +
				"nobody typed a window. Before ticket 256 this process always held 3s. Whatever the 255/265 " +
				"wording for the ⓑ sentence ends up saying, it has to say 2s here and not 3s",
		},
		{
			name:     "window seeded to 2s - differs from the pre-256 compiled 3s, so it is a real reading",
			seed:     "l1_window_sec = 2\n",
			wantTime: approval.DefaultApprovalTimeout,
			wantWin:  2 * time.Second,
			winIsWhat: "a genuine move off the compiled 3s default, AND it is inside the band, so it is " +
				"the only window case here that counts as a positive control",
		},
		{
			name:      "both seeded",
			seed:      "confirm_timeout_sec = 45\nl1_window_sec = 2\n",
			wantTime:  45 * time.Second,
			wantWin:   2 * time.Second,
			winIsWhat: "same as above: 2s is inside [2s,3s], so the band let the file through",
		},
		{
			name:     "window way too large - reverse control, the clamp must survive",
			seed:     "l1_window_sec = 99\n",
			wantTime: approval.DefaultApprovalTimeout,
			wantWin:  approval.MaxL1Window,
			winIsWhat: "the CLAMP, not the config: 99s is refused down to MaxL1Window, and that is exactly " +
				"why Window() must never be used as the positive control (it would read green for 99 and for 3)",
		},
		{
			name:      "window below the floor - reverse control, the clamp must survive",
			seed:      "l1_window_sec = 1\n",
			wantTime:  approval.DefaultApprovalTimeout,
			wantWin:   approval.MinL1Window,
			winIsWhat: "the CLAMP again: 1s is raised to MinL1Window",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			writeRiskConfig256(t, dir, c.seed)
			ra := newResidentApprovalWithConfig(dir)
			if ra.riskProvenance != riskProvenanceRead {
				t.Fatalf("provenance = %q, want %q: the file reads, so the gate must not claim a fallback",
					ra.riskProvenance, riskProvenanceRead)
			}
			if got := ra.gate.Queue().Timeout(); got != c.wantTime {
				t.Errorf("Queue().Timeout() = %v, want %v for seed %q - this is AC#2's positive control, "+
					"a miss means the resident leg still ignores confirm_timeout_sec", got, c.wantTime, c.seed)
			}
			if got := ra.gate.Window(); got != c.wantWin {
				t.Errorf("Window() = %v, want %v for seed %q; what this reading proves: %s",
					got, c.wantWin, c.seed, c.winIsWhat)
			}
			// The band is a judgement, not a preference: it must never be opened
			// by feeding it a value. If it moved, this whole ticket's Window()
			// readings would have to be re-read, so it is pinned here rather than
			// trusted from queue.go.
			if w := ra.gate.Window(); w < approval.MinL1Window || w > approval.MaxL1Window {
				t.Errorf("Window() = %v is outside [%v,%v]: gate.go's clamp is gone, and every Window() "+
					"reading in this ticket has to be re-adjudicated", w, approval.MinL1Window, approval.MaxL1Window)
			}
		})
	}
}

// ----------------------------------------------------------------- ③ the limitation

// TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly is ruler ③: it RUNS
// the limitation instead of describing it. approval.Gate copies the window into
// g.window and NewQueue copies the deadline into q.timeout, both inside
// construction, and the repository has no re-apply path for either. So the
// sentence this ticket is allowed to say is "the gate this launch is built with
// follows [risk]" - never "change the file and the live process changes".
//
// The paired positive control (the file really did become 90s, read back through
// config.LoadFile after the rewrite) is what keeps this from being a stale-read
// artifact: without it, "still 45" would be equally consistent with my rewrite
// having never landed.
func TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly(t *testing.T) {
	dir := t.TempDir()
	writeRiskConfig256(t, dir, "confirm_timeout_sec = 45\n")
	ra := newResidentApprovalWithConfig(dir)
	if got := ra.gate.Queue().Timeout(); got != 45*time.Second {
		t.Fatalf("setup: the gate was not built with the seeded 45s, got %v; ruler ② owns that failure, "+
			"but this case cannot measure a limitation on a value it never had", got)
	}

	writeRiskConfig256(t, dir, "confirm_timeout_sec = 90\n")

	// Positive control on the rewrite itself: the file on disk really says 90.
	c, _, err := config.LoadFile(filepath.Join(dir, configFileName), nil)
	if err != nil || c == nil {
		t.Fatalf("re-reading the rewritten config.toml: err=%v c==nil=%v - the rewrite is unverifiable", err, c == nil)
	}
	if c.Risk.ConfirmTimeoutSec != 90 {
		t.Fatalf("the rewrite did not land: config.toml now answers confirm_timeout_sec = %d, want 90; "+
			"without this the assertion below would be measuring a file that never changed",
			c.Risk.ConfirmTimeoutSec)
	}

	// The limitation, run: the already-built gate keeps the value it was built
	// with. If a future leg adds a re-apply path for [risk], THIS case reddens and
	// says out loud that AC#2's wording may now be upgraded - that is the point of
	// pinning a limitation rather than writing it in prose.
	if got := ra.gate.Queue().Timeout(); got != 45*time.Second {
		t.Errorf("the live gate moved from its construction value (%v) to %v after config.toml changed. "+
			"Either a re-apply path was added for [risk] - in which case ticket 256's stated limitation is "+
			"obsolete and has to be re-adjudicated, and the 255/265 wording for the ⓑ sentence with it - or "+
			"the timeout is now being read from somewhere other than construction, which is worse.",
			45*time.Second, got)
	}
	if ra.riskProvenance != riskProvenanceRead {
		t.Errorf("provenance drifted after the file changed: %q. Provenance is a construction-time receipt, "+
			"not a live probe.", ra.riskProvenance)
	}
	// And a freshly built leg over the same directory does take the new number:
	// the pair "old gate keeps 45 / new gate takes 90" is the whole shape of the
	// limitation in one case.
	fresh := newResidentApprovalWithConfig(dir)
	if got := fresh.gate.Queue().Timeout(); got != 90*time.Second {
		t.Errorf("a gate built AFTER the rewrite reads %v, want 90s: without this half the assertion above "+
			"would be indistinguishable from 'the leg never reads config'", got)
	}
}

// ------------------------------------------------------------ ④ the landing site

// optionsFieldSet256 parses one Go file and returns, for every call of the shape
// <pkg>.New(<pkg>.Options{...}), the field names the composite literal really
// passes, plus how many such literals the file holds. Comments are not in the
// tree, so this cannot be turned green by writing a word anywhere.
func optionsFieldSet256(t *testing.T, path string) (sets [][]string, literals int) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("parse %s (the instrument, not the code): %v", path, err)
	}
	ast.Inspect(f, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		sel, ok := lit.Type.(*ast.SelectorExpr)
		if !ok || sel.Sel == nil || sel.Sel.Name != "Options" {
			return true
		}
		pkg, ok := sel.X.(*ast.Ident)
		if !ok || pkg.Name != "approval" {
			return true
		}
		literals++
		var names []string
		for _, e := range lit.Elts {
			kv, ok := e.(*ast.KeyValueExpr)
			if !ok {
				t.Errorf("%s: an approval.Options element is not key: value form, so the field count "+
					"would silently under-measure; re-read this ruler", path)
				continue
			}
			id, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}
			names = append(names, id.Name)
		}
		sets = append(sets, names)
		return true
	})
	return sets, literals
}

// declaredOptionsFields256 reads the denominator off approval's own declaration
// instead of hard-coding "10": a field added to Options shows up here, which is
// what makes ④'s "5 of N" a reading rather than a remembered number.
func declaredOptionsFields256(t *testing.T) []string {
	t.Helper()
	root := filepath.Dir(filepath.Dir(risk256Dir(t)))
	path := filepath.Join(root, "internal", "agent", "approval", "gate.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("parse %s for Options' field count: %v", path, err)
	}
	for _, d := range f.Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, s := range gd.Specs {
			ts, ok := s.(*ast.TypeSpec)
			if !ok || ts.Name.Name != "Options" {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				t.Fatalf("Options in %s is not a struct any more", path)
			}
			var out []string
			for _, fld := range st.Fields.List {
				for _, n := range fld.Names {
					out = append(out, n.Name)
				}
			}
			return out
		}
	}
	t.Fatalf("no Options type declaration in %s: the denominator of ruler ④ cannot be read", path)
	return nil
}

// TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims is ruler ④: the
// resident leg's real Options field set, counted off the syntax tree. Before
// ticket 256 it was 3 (UI / Channels / Logf); it is now 5. Either direction of
// change reddens - deleting the two new fields, or adding a third one nobody
// adjudicated.
//
// Grants is asserted ABSENT on purpose. AC#1's other half is not this leg's work
// (256-a2 census §0 judged it a different proposition, and it is filed as
// pending ticket 265), so this ruler refuses to let the shipped shape be read as
// "Options fully fed".
func TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims(t *testing.T) {
	path := filepath.Join(risk256Dir(t), "resident_approval_windows.go")
	sets, literals := optionsFieldSet256(t, path)
	if literals != 1 {
		t.Fatalf("found %d approval.Options literals in %s, want exactly 1. Two would mean the resident leg "+
			"builds a second gate, which TestAC246ResidentPipelineAsksThroughTheOneGate (the ar.gate != ra.gate "+
			"pointer identity pin) already refuses; this ruler is the earlier, cheaper version of that.",
			literals, path)
	}
	if len(sets) != 1 {
		t.Fatalf("collected %d field sets from %d literals - the instrument and the count disagree", len(sets), literals)
	}

	got := map[string]bool{}
	for _, n := range sets[0] {
		if got[n] {
			t.Errorf("Options field %q passed twice in the same literal, so the count below double-counts it", n)
		}
		got[n] = true
	}
	want := []string{"UI", "Channels", "Window", "ApprovalTimeout", "Logf"}
	for _, w := range want {
		if !got[w] {
			t.Errorf("the resident leg does NOT pass Options.%s to approval.New. That is the whole of ticket 256's "+
				"[risk] half: Window and ApprovalTimeout are the two fields this leg exists to feed.", w)
		}
	}
	if got["Grants"] {
		t.Errorf("the resident leg now passes Options.Grants. That is NOT ticket 256-r1's scope: the session " +
			"ledger has no value to hand over at this moment (256-a2 census §0), and the half was filed as " +
			"pending ticket 265. Either a seam was minted outside that ticket, or this leg grew.")
	}
	for n := range got {
		found := false
		for _, w := range want {
			if w == n {
				found = true
			}
		}
		if !found {
			t.Errorf("the resident leg passes an unexpected Options field %q: the shipped field set is no longer "+
				"the one ticket 256 adjudicated, so re-read what this gate is now being fed", n)
		}
	}
	if len(got) != len(want) {
		t.Errorf("resident leg passes %d of Options' fields, want %d", len(got), len(want))
	}

	declared := declaredOptionsFields256(t)
	t.Logf("landing site reading: resident leg passes %d of %d declared approval.Options fields (was 3 of %d before ticket 256)",
		len(got), len(declared), len(declared))
	if len(declared) < 10 {
		t.Errorf("approval.Options now declares %d fields, fewer than the 10 the 256-a2 census §1 rostered: "+
			"the denominator moved, so that census table has to be re-adjudicated", len(declared))
	}
	for n := range got {
		found := false
		for _, d := range declared {
			if d == n {
				found = true
			}
		}
		if !found {
			t.Errorf("Options.%s is passed but not declared in gate.go - the ruler read a stale tree", n)
		}
	}
}

// TestTicket256ResidentBootPassesTheDataDirToTheGate is ruler ⑤, the other half
// of "did it land": the shipped resident boot must hand its own data dir to the
// config-fed constructor. Measured off the tree of resident_windows.go, so a
// reverted call site reddens even though every other ruler here would still pass
// (they construct the gate directly and would never notice the boot's choice).
func TestTicket256ResidentBootPassesTheDataDirToTheGate(t *testing.T) {
	path := filepath.Join(risk256Dir(t), "resident_windows.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("parse %s (the instrument, not the code): %v", path, err)
	}
	var fn *ast.FuncDecl
	for _, d := range f.Decls {
		if FD, ok := d.(*ast.FuncDecl); ok && FD.Name != nil && FD.Name.Name == "runResident" {
			fn = FD
		}
	}
	if fn == nil {
		t.Fatalf("no runResident in %s: the boot function this ruler measures has moved or been renamed", path)
	}

	var fed, bare int
	var argOK bool
	ast.Inspect(fn, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		id, ok := call.Fun.(*ast.Ident)
		if !ok {
			return true
		}
		switch id.Name {
		case "newResidentApprovalWithConfig":
			fed++
			if len(call.Args) == 1 {
				if sel, ok := call.Args[0].(*ast.SelectorExpr); ok && sel.Sel != nil && sel.Sel.Name == "DataDir" {
					argOK = true
				}
			}
		case "newResidentApproval":
			// Exact name, not prefix: the two constructors are distinct calls.
			bare++
		}
		return true
	})
	if fed != 1 {
		t.Errorf("runResident calls newResidentApprovalWithConfig %d times, want exactly 1", fed)
	}
	if bare != 0 {
		t.Errorf("runResident still calls the no-host-view newResidentApproval() %d time(s): the shipped "+
			"resident leg would be back to the compiled 300s / 3s, which is the exact defect ticket 256 was "+
			"filed for. The bare signature exists for the 246 roster only (orchestrator ruling §8.2).", bare)
	}
	if !argOK {
		t.Errorf("runResident's call does not pass a *.DataDir value, so the gate is not being handed the " +
			"host's own directory: the read would go to the wrong config.toml or to none")
	}
}
