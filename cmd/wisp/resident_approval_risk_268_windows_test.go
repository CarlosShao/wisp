//go:build windows

package main

// Ticket 268 AC#1/AC#2: the resident leg has to answer with TWO DIFFERENT NAMES
// for "there is no config.toml" and "config.toml is there and the loader refused
// what it says", and the refused one has to reach a human eye.
//
// WHY THE DEFECT EXISTED AT ALL. Before ticket 267 any number in
// [risk].confirm_timeout_sec flowed into the gate this process builds. Ticket 267
// added the band [31, 3600], so an out-of-band number now dies in the loader -
// and residentRiskGateValues treated that refusal as "the file would not read",
// answering with riskProvenanceUnreadable and logging a sentence blaming a file
// that is sitting right there. The safety outcome is unchanged (the fallback is
// still approval's compiled 300s, and 300 > 30 keeps C18's warning armed); what
// was lost is the receipt.
//
// THREE RULERS, and the one fact each is the only one that can hold:
//
//   ① TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig
//      the two shapes carry different provenance words, the four words are four
//      distinct literals, the refused shape still runs on the compiled constants
//      (AC#2), and the planted out-of-band number is NOT live (Q-77 stays shut).
//      Its premise control is the config.LoadFile call at the top: it proves the
//      planted file exists and is refused for a reason that is not "missing", so
//      the ruler below cannot be measuring two absent files.
//   ② TestTicket268RefusedRiskConfigReachesStdoutOnce
//      AC#1's "at least once to the user's eyes", run: from the same construction
//      the terminal face carries exactly one line naming this shape, and it
//      carries the loader's own answer in it. Its control is the missing shape on
//      the same face, which must stay as silent as it has been since ticket 256.
//   ③ TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords
//      the shape of the branch, read off the syntax tree: the split is made by
//      errors.Is against fs.ErrNotExist, the branch prints, and no strings.* call
//      and no err.Error() sits in that function. Comments are invisible to it, so
//      writing the word "sentinel" anywhere cannot turn it green - and the
//      rejected alternative (config_reload.go's strings.HasPrefix on the detail
//      string) would redden it.
//
// WHAT THESE RULERS DELIBERATELY DO NOT CLAIM:
//   - They do not make the refused value take effect. AC#2 keeps the compiled
//     300s, and whether the band's floor should move is Q-77 (owner's call).
//   - They do not name WHICH loader rule refused. "refused at load" is a class,
//     not a diagnosis: ticket 267's band, a syntax error, an unknown key and a
//     migration refusal all land in this one branch, and separating them would
//     mean reading error prose (see ruler ③'s prohibition). The verbatim loader
//     answer travels in the log attribute and in the stdout line instead.
//   - They do not add a second [risk] read point. Everything here goes through
//     the constructor ticket 256 already ships, whose AST ruler in
//     resident_approval_risk_256_windows_test.go counts the constructors
//     runResident calls - that ruler, and the 300s and construction-time pins in
//     the same file, are untouched by this ticket.
//   - They do not touch the console leg. wisp run still exits 2 with its own
//     loud sentence (ticket 267, AC#3 here).
//
// The helpers reused from resident_approval_risk_256_windows_test.go are that
// file's own risk256Dir (this package's directory) and writeRiskConfig256 (a
// minimal readable config.toml carrying a passed [risk] body): same package, same
// purpose, and duplicating them here would mint a second place to be wrong about
// what "a planted config.toml" means.

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/config"
)

// risk268RefusedSeed is the number from ticket 268's own story: a user who thinks
// 300s is too long writes 20. It is below the band's floor (confirmTimeoutSecMin
// = 31 in internal/config/validate.go), which is exactly why the loader refuses
// the whole file. The band itself is not this ticket's to move.
const risk268RefusedSeed = "confirm_timeout_sec = 20\n"

// risk268StdoutMarker is the prefix of the line ticket 268 AC#1 adds, and the
// ruler below counts THAT, not any stdout traffic.
const risk268StdoutMarker = "wisp: resident [risk]:"

// ---------------------------------------------------------------- AC#1 + AC#2

func TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig(t *testing.T) {
	missingDir := t.TempDir()
	missing := newResidentApprovalWithConfig(missingDir)

	refusedDir := t.TempDir()
	refusedPath := writeRiskConfig256(t, refusedDir, risk268RefusedSeed)

	// Premise first, because every assertion below is about a file that EXISTS.
	if _, err := os.Stat(refusedPath); err != nil {
		t.Fatalf("the planted config.toml is not there (%v), so this case would be measuring the missing shape twice", err)
	}
	_, _, loadErr := config.LoadFile(refusedPath, nil)
	if loadErr == nil {
		t.Fatalf("setup: %s was accepted by the loader, so there is no refused branch to name; the seed is %q", refusedPath, risk268RefusedSeed)
	}
	if errors.Is(loadErr, fs.ErrNotExist) {
		t.Fatalf("setup: the loader answered with the not-exist sentinel (%v) on a file that stats clean - the split this ticket lands is built on exactly that sentinel, so the premise has to hold", loadErr)
	}
	// A test may match prose (production may not; ruler ③ owns that half): this
	// is the check that the refusal really is ticket 267's band.
	if !strings.Contains(loadErr.Error(), "out of range [31, 3600]") {
		t.Fatalf("setup: the planted file was refused, but not by the band - the loader said %q. The ruler below would then be naming some other shape refused", loadErr)
	}

	refused := newResidentApprovalWithConfig(refusedDir)

	if missing.riskProvenance != riskProvenanceUnreadable {
		t.Errorf("missing-file provenance = %q, want %q: ticket 256's ruler ① pins this word to the absent file and this ticket does not move it",
			missing.riskProvenance, riskProvenanceUnreadable)
	}
	if refused.riskProvenance != riskProvenanceRefusedAtLoad {
		t.Errorf("refused-file provenance = %q, want %q: AC#1 asks for a name of its own for the shape ticket 267's band created",
			refused.riskProvenance, riskProvenanceRefusedAtLoad)
	}
	if missing.riskProvenance == refused.riskProvenance {
		t.Errorf("the two shapes are folded back into one word %q - that fold IS ticket 268: the user wrote a number, the gate runs the compiled 300s, and the receipt blames a file that is sitting in that directory",
			missing.riskProvenance)
	}
	if !strings.Contains(refused.riskProvenance, "present") {
		t.Errorf("refused-file provenance %q does not say out loud that the file is present, which is the whole difference from the missing shape",
			refused.riskProvenance)
	}

	// Four words, four literals: folding any pair of them reddens here even if
	// the case names above were rewritten to match.
	words := []struct {
		name  string
		value string
	}{
		{"riskProvenanceRead", riskProvenanceRead},
		{"riskProvenanceUnreadable", riskProvenanceUnreadable},
		{"riskProvenanceRefusedAtLoad", riskProvenanceRefusedAtLoad},
		{"riskProvenanceNoView", riskProvenanceNoView},
	}
	seen := map[string]string{}
	for _, w := range words {
		if prev, dup := seen[w.value]; dup {
			t.Errorf("provenance constants %s and %s share the literal %q - the four readings have to be four names", prev, w.name, w.value)
		}
		seen[w.value] = w.name
	}

	// AC#2: both fallback shapes still run on the compiled constants, and the
	// refused number is NOT what the gate holds.
	for _, c := range []struct {
		name string
		ra   *residentApproval
	}{
		{"missing config.toml", missing},
		{"refused config.toml", refused},
	} {
		if got := c.ra.gate.Queue().Timeout(); got != approval.DefaultApprovalTimeout {
			t.Errorf("Queue().Timeout() with a %s = %v, want the compiled %v. AC#2: the fallback must not grow a value of its own, and the refused number must not sneak in through the door ticket 268 just opened (Q-77 is the owner's call, not this leg's)",
				c.name, got, approval.DefaultApprovalTimeout)
		}
		if got := c.ra.gate.Queue().Timeout(); got == 20*time.Second {
			t.Errorf("Queue().Timeout() = 20s: the value the band refused became live with %s", c.name)
		}
		if got := c.ra.gate.Window(); got != approval.DefaultL1Window {
			t.Errorf("Window() with a %s = %v, want the compiled %v", c.name, got, approval.DefaultL1Window)
		}
	}

	// Paired positive control on the SAME directory: once the number is inside
	// the band the file reads, so "refused" above is the file's content talking
	// and not this test's choice of directory.
	writeRiskConfig256(t, refusedDir, "confirm_timeout_sec = 45\n")
	fixed := newResidentApprovalWithConfig(refusedDir)
	if fixed.riskProvenance != riskProvenanceRead {
		t.Errorf("after repairing the seed in the same directory, provenance = %q, want %q - without this half the refused verdict above could be a property of the directory, not of the value",
			fixed.riskProvenance, riskProvenanceRead)
	}
	if got := fixed.gate.Queue().Timeout(); got != 45*time.Second {
		t.Errorf("after repairing the seed, Queue().Timeout() = %v, want 45s", got)
	}
}

// ------------------------------------------------------------------- AC#1 face 2

func TestTicket268RefusedRiskConfigReachesStdoutOnce(t *testing.T) {
	dir := t.TempDir()
	writeRiskConfig256(t, dir, risk268RefusedSeed)

	var ra *residentApproval
	out := captureStdout128(t, func() { ra = newResidentApprovalWithConfig(dir) })

	if ra.riskProvenance != riskProvenanceRefusedAtLoad {
		t.Fatalf("this case is about the refused branch, but the gate was built with provenance %q - nothing below measures what it claims to",
			ra.riskProvenance)
	}
	if n := strings.Count(out, risk268StdoutMarker); n != 1 {
		t.Errorf("the refused branch wrote %d stdout lines carrying %q, want exactly 1. AC#1 asks that this shape reach a user's eye at least once from a terminal; one line is the shape resident_windows.go's [hotkey] fallback has carried since ticket 258",
			n, risk268StdoutMarker)
	}
	if !strings.Contains(out, "out of range [31, 3600]") {
		t.Errorf("the stdout line has to carry the loader's own answer, not just this file's class name; captured stdout was:\n%s", out)
	}
	if !strings.Contains(out, "DefaultApprovalTimeout=300s") {
		t.Errorf("the stdout line has to say what the gate runs on instead; captured stdout was:\n%s", out)
	}

	// Control on the same face: the missing shape keeps its silence. Ticket 268
	// opens the refused branch only - AC#3 forbids smoothing the legs' different
	// postures into one, and the same logic covers the two fallback shapes.
	missingDir := t.TempDir()
	missingOut := captureStdout128(t, func() { newResidentApprovalWithConfig(missingDir) })
	if strings.Contains(missingOut, risk268StdoutMarker) {
		t.Errorf("the missing-file branch wrote a %q line: the two shapes are sharing this face again, which is the fold this ticket exists to close",
			risk268StdoutMarker)
	}
}

// ----------------------------------------------------------- AC#1 shape of branch

func TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords(t *testing.T) {
	path := filepath.Join(risk256Dir(t), "resident_approval_windows.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("parse %s (the instrument, not the code): %v", path, err)
	}
	var fn *ast.FuncDecl
	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Name != nil && fd.Name.Name == "residentRiskGateValues" {
			fn = fd
		}
	}
	if fn == nil {
		t.Fatalf("no residentRiskGateValues in %s: the function this ruler measures has moved or been renamed", path)
	}

	sentinelIs, proseCalls, printfs := 0, []string{}, 0
	returned := map[string]bool{}
	ast.Inspect(fn, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			sel, ok := node.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pkg, ok := sel.X.(*ast.Ident)
			if !ok {
				return true
			}
			switch {
			case pkg.Name == "errors" && sel.Sel.Name == "Is" && len(node.Args) == 2:
				if arg, ok := node.Args[1].(*ast.SelectorExpr); ok && arg.Sel != nil && arg.Sel.Name == "ErrNotExist" {
					sentinelIs++
				}
			case pkg.Name == "strings":
				proseCalls = append(proseCalls, "strings."+sel.Sel.Name)
			case pkg.Name == "err" && sel.Sel.Name == "Error":
				proseCalls = append(proseCalls, "err.Error()")
			case pkg.Name == "fmt" && sel.Sel.Name == "Printf":
				printfs++
			}
		case *ast.ReturnStmt:
			for _, r := range node.Results {
				if id, ok := r.(*ast.Ident); ok {
					returned[id.Name] = true
				}
			}
		}
		return true
	})

	if sentinelIs != 1 {
		t.Errorf("residentRiskGateValues makes %d errors.Is(err, fs.ErrNotExist) calls, want exactly 1: the missing/refused split has to rest on the sentinel, which is the one machine-readable fact the loader's chain carries (internal/config/loader.go keeps the original error in it, and config_reload.go's cause=missing branch already reads the same sentinel)",
			sentinelIs)
	}
	if len(proseCalls) != 0 {
		t.Errorf("residentRiskGateValues classifies by prose through %v: matching an error message is the shape config_reload.go:396 already carries and ticket 268 was told not to copy - a reworded loader sentence would silently re-fold the two shapes",
			proseCalls)
	}
	if printfs != 1 {
		t.Errorf("residentRiskGateValues holds %d fmt.Printf calls, want exactly 1: AC#1's user-visible face is the other half of this branch, and a branch that only logs to disk does not reach the eye of anyone who launched the process from a terminal",
			printfs)
	}
	for _, name := range []string{"riskProvenanceRead", "riskProvenanceUnreadable", "riskProvenanceRefusedAtLoad", "riskProvenanceNoView"} {
		if !returned[name] {
			t.Errorf("residentRiskGateValues never returns %s: one of the four readings has no exit of its own in the tree", name)
		}
	}
}
