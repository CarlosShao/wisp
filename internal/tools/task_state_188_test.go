package tools

// Ticket 188 AC#2 + AC#3 (dispatch 2026-09-28 §G, ruling A394 (i)).
//
// What these legs own:
//
//	AC#2  TaskOutput now carries a status dimension, and the only legal values
//	      are the 20 names in internal/statemachine/states.go:11-31 - the
//	      verbatim in-repo copy of D43's frozen transition table (its header
//	      line :8 reads "names exactly as in D43 / SPEC-08"; source
//	      docs/PLAN.md:3055 plus :3061-3100). Nothing invented, and the memory
//	      side's legacy words (done / cancelled / running / succeeded, written
//	      by internal/agent/loop.go:983-998 and cmd/wisp/run.go:651-670 into
//	      memory.TaskLog.State, models.go:62 with no CHECK at :57) are refused
//	      here by name. Those four words belong to ticket 196, and per A394
//	      this ticket does not unify the two vocabularies and does not touch
//	      either producer.
//	AC#3  The status dimension has Go-side producers only: the panel may
//	      display and may initiate a request (R20 / ticket 92), never write
//	      this dimension. That is pinned as a source scan with its own positive
//	      controls, and it also pins the read side so ticket 145's snapshot
//	      pump does not get caught by its own gate.
//
// Every negative leg below carries a denominator (how many files it actually
// looked at) so "no findings" can never mean "looked at nothing".

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// d43StateNames is the 20-value set copied out of internal/statemachine/
// states.go:11-31 verbatim, in that file's order. It is written by hand here on
// purpose: the leg that follows fails if this list and the frozen file ever
// drift apart, which is the only way a test can notice a "new" state name being
// added to the vocabulary side.
var d43StateNames = []statemachine.State{
	statemachine.StateFirstRun,
	statemachine.StateSleeping,
	statemachine.StateArmed,
	statemachine.StateMuted,
	statemachine.StateListening,
	statemachine.StateThinking,
	statemachine.StateActing,
	statemachine.StateSpeaking,
	statemachine.StateWarm,
	statemachine.StateConversation,
	statemachine.StateConfirming,
	statemachine.StateAwaitingApproval,
	statemachine.StateSettling,
	statemachine.StateDownloading,
	statemachine.StateUnconfigured,
	statemachine.StateNoNetwork,
	statemachine.StateError,
	statemachine.StateQueued,
	statemachine.StateStuck,
	statemachine.StateWatchdogAlert,
}

// legacyTaskStateWords is ticket 196's vocabulary: what the memory side already
// writes into TaskLog.State. None of them is a D43 name, and none of them may
// enter this dimension.
var legacyTaskStateWords = []string{
	"done", "cancelled", "running", "succeeded", "interrupted", "error",
	"finished", "failed", "stuck", "queued", "pending",
}

func TestTaskState188AC2VocabularyIsD43Only(t *testing.T) {
	if len(d43StateNames) != 20 {
		t.Fatalf("this leg names %d states, not the 20 of D43: %v", len(d43StateNames), d43StateNames)
	}
	for _, s := range d43StateNames {
		if !statemachine.Valid(s) {
			t.Errorf("D43 name %q is not accepted by statemachine.Valid (states.go:39)", string(s))
		}
		got, why := TaskOutput{Text: "x", State: s}.StateAnswer()
		if got != s || why != "" {
			t.Errorf("StateAnswer(%q) = (%q, %q), want the filed name back with no notice", string(s), string(got), why)
		}
	}

	// The legacy family, refused by name. Each one is a string a caller COULD
	// write (statemachine.State is a string type, so the compiler alone does not
	// close this door), and each one is refused at the two doors this ticket
	// owns: the reader and the writer.
	for _, w := range legacyTaskStateWords {
		s := statemachine.State(w)
		if statemachine.Valid(s) {
			t.Fatalf("legacy word %q is now a valid D43 name - either the frozen table moved or this leg is aimed at the wrong file", w)
		}
		got, why := TaskOutput{Text: "x", State: s}.StateAnswer()
		if got != "" || why == "" {
			t.Errorf("StateAnswer(%q) = (%q, %q): a non-D43 word must never be vouched for", w, string(got), why)
		}
		if !strings.Contains(why, "D43") {
			t.Errorf("StateAnswer(%q) notice %q must name the table it was judged against", w, why)
		}
	}

	// Empty is its own answer, and it is not "finished".
	got, why := TaskOutput{Text: "x"}.StateAnswer()
	if got != "" || why == "" || !strings.Contains(why, "没有登记") {
		t.Errorf("StateAnswer(empty) = (%q, %q), want no state and a 未登记 notice", string(got), why)
	}
}

// TestTaskState188AC2ValueComesFromTheRoster is AC#2's "值来自真名册，不许硬编码"
// leg: the status a reader sees must be the one the record in the table holds,
// and it must move when the record moves.
func TestTaskState188AC2ValueComesFromTheRoster(t *testing.T) {
	roster := NewTaskRoster()
	roster.Record("bg-a", TaskOutput{Text: "甲的输出", State: statemachine.StateActing})
	roster.Record("bg-b", TaskOutput{Text: "乙的输出", State: statemachine.StateError})
	roster.Record("bg-c", TaskOutput{Text: "丙的输出"})

	for id, want := range map[string]statemachine.State{
		"bg-a": statemachine.StateActing,
		"bg-b": statemachine.StateError,
	} {
		rec, ok := roster.Look(id)
		if !ok {
			t.Fatalf("Look(%q) said the table has no such task", id)
		}
		got, why := rec.StateAnswer()
		if why != "" || got != want {
			t.Errorf("roster record %q answers (%q, %q), want (%q, \"\")", id, string(got), why, string(want))
		}
	}

	// The third record was filed without a state: same table, different answer.
	// A reader that hardcoded a verdict would get this wrong in exactly one
	// direction, and this leg is what tells the two apart.
	rec, ok := roster.Look("bg-c")
	if !ok {
		t.Fatalf("Look(%q) said the table has no such task", "bg-c")
	}
	if got, why := rec.StateAnswer(); got != "" || why == "" {
		t.Errorf("record with no filed state answered (%q, %q), want the 未登记 notice", string(got), why)
	}

	// Last writer wins is Record's documented semantics (task.go:120-133) and it
	// covers this dimension too: re-file with a different name and the reader
	// follows the table, not a constant.
	roster.Record("bg-a", TaskOutput{Text: "甲的输出", State: statemachine.StateSettling})
	if got, why := mustLook(t, roster, "bg-a").StateAnswer(); got != statemachine.StateSettling || why != "" {
		t.Errorf("after re-filing, bg-a answers (%q, %q), want (Settling, \"\")", string(got), why)
	}
}

func mustLook(t *testing.T, r *TaskRoster, id string) TaskOutput {
	t.Helper()
	rec, ok := r.Look(id)
	if !ok {
		t.Fatalf("Look(%q) said the table has no such task", id)
	}
	return rec
}

// TestTaskState188AC2BackfillIsTheGoSideProducer covers AC#3's other half: the
// producer for this dimension lives on the host/tools side, and it is the only
// place A394 put a writer for this ticket.
func TestTaskState188AC2BackfillIsTheGoSideProducer(t *testing.T) {
	roster := NewTaskRoster()
	ctx := context.Background()

	// (a) Today's production shape (cmd/wisp/run.go:640-647 sets Roster and
	// Spills only): the record lands, and its state dimension stays honestly
	// unfiled rather than being invented by the writer.
	rec, why := TaskBackfill{Roster: roster}.Backfill(ctx, "bg-zero", "零状态正文")
	if why != "" {
		t.Fatalf("backfill with no configured state was refused: %q", why)
	}
	if rec.State != "" {
		t.Errorf("backfill filed %q for a host that filed nothing", string(rec.State))
	}
	if got, notice := mustLook(t, roster, "bg-zero").StateAnswer(); got != "" || notice == "" {
		t.Errorf("unwired producer answered (%q, %q), want the 未登记 notice", string(got), notice)
	}

	// (b) A host that does wire a D43 name gets it into the table, so the field
	// is a carrier and not decoration.
	rec, why = TaskBackfill{Roster: roster, State: statemachine.StateActing}.Backfill(ctx, "bg-acting", "正文")
	if why != "" {
		t.Fatalf("backfill with a D43 state was refused: %q", why)
	}
	if rec.State != statemachine.StateActing {
		t.Errorf("returned record carries %q, want Acting", string(rec.State))
	}
	if got, notice := mustLook(t, roster, "bg-acting").StateAnswer(); got != statemachine.StateActing || notice != "" {
		t.Errorf("roster record answers (%q, %q), want (Acting, \"\")", string(got), notice)
	}

	// (c) The door, with its positive control: wire one of ticket 196's legacy
	// words and the filing is refused before the table is touched - no record,
	// and a reason that names the table.
	before := roster.Count()
	for _, w := range legacyTaskStateWords {
		got, why := TaskBackfill{Roster: roster, State: statemachine.State(w)}.Backfill(ctx, "bg-legacy-"+w, "正文")
		if why == "" {
			t.Fatalf("backfill accepted the legacy word %q: the door is not installed", w)
		}
		if got.State != "" || got.Text != "" {
			t.Errorf("refused backfill still returned a record for %q: %+v", w, got)
		}
		if !strings.Contains(why, w) || !strings.Contains(why, "D43") {
			t.Errorf("refusal for %q says %q: it must name the word and the table", w, why)
		}
		if _, ok := roster.Look("bg-legacy-" + w); ok {
			t.Errorf("a record built from the legacy word %q reached the table", w)
		}
	}
	if after := roster.Count(); after != before {
		t.Errorf("refused filings changed the table: %d records before, %d after", before, after)
	}
}

// ---------------------------------------------------------------------------
// AC#3: no panel-side write leg for this dimension.
// ---------------------------------------------------------------------------

// stateWriteShape is one way source code can write the task status dimension.
type stateWriteShape struct {
	file string
	line int
	why  string
}

// rosterWriteMethods are the two surface names that put a record (and with it
// the status dimension) into the task table. Look and Count are deliberately
// absent: reading the table is what the panel is allowed to do, and ticket 145's
// snapshot pump leg is a reader - a gate that flagged readers would put this
// ticket's own downstream in its own red.
var rosterWriteMethods = map[string]bool{"Record": true, "Backfill": true}

// scanTaskStateWriteLegs walks the non-test Go sources under root and reports
// every shape that writes the task status dimension, plus any panel-side
// whitelisted method constant that even claims the task axis.
//
// It returns the findings and the number of files it actually parsed, because a
// scan whose denominator is zero reports "clean" for a tree it never read.
func scanTaskStateWriteLegs(t *testing.T, root string) ([]stateWriteShape, int) {
	t.Helper()
	var found []stateWriteShape
	examined := 0
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		fset := token.NewFileSet()
		file, perr := parser.ParseFile(fset, path, nil, 0)
		if perr != nil {
			t.Fatalf("parse %s: %v", path, perr)
		}
		examined++
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				return true
			}
			lineno := fset.Position(n.Pos()).Line
			switch v := n.(type) {
			case *ast.CallExpr:
				if sel, ok := v.Fun.(*ast.SelectorExpr); ok && rosterWriteMethods[sel.Sel.Name] {
					found = append(found, stateWriteShape{
						path, lineno,
						"calls the roster/backfill write method " + sel.Sel.Name,
					})
				}
			case *ast.CompositeLit:
				// The record is spelled TaskOutput inside package tools and
				// tools.TaskOutput from any other package, so both shapes count.
				typeName := ""
				switch ty := v.Type.(type) {
				case *ast.Ident:
					typeName = ty.Name
				case *ast.SelectorExpr:
					typeName = ty.Sel.Name
				}
				if typeName != "TaskOutput" {
					return true
				}
				// A record constructed here is a record being written: either it
				// keys the State field, or it is positional and long enough to
				// carry the status.
				if len(v.Elts) >= 3 {
					found = append(found, stateWriteShape{
						path, lineno,
						"constructs a TaskOutput positionally, status included",
					})
					return true
				}
				for _, el := range v.Elts {
					kv, ok := el.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					if k, ok := kv.Key.(*ast.Ident); ok && k.Name == "State" {
						found = append(found, stateWriteShape{
							path, lineno,
							"constructs a TaskOutput with a State key",
						})
					}
				}
			case *ast.AssignStmt:
				for _, lhs := range v.Lhs {
					if sel, ok := lhs.(*ast.SelectorExpr); ok && sel.Sel.Name == "State" {
						found = append(found, stateWriteShape{
							path, lineno,
							"assigns to a State field",
						})
					}
				}
			case *ast.ValueSpec:
				for i, nm := range v.Names {
					if !strings.HasPrefix(nm.Name, "Method") || i >= len(v.Values) {
						continue
					}
					lit, ok := v.Values[i].(*ast.BasicLit)
					if !ok {
						continue
					}
					lower := strings.ToLower(lit.Value)
					if strings.Contains(lower, "task") || strings.Contains(lower, "state") {
						found = append(found, stateWriteShape{
							path, lineno,
							"claims a panel-side whitelisted method on the task axis: " + nm.Name,
						})
					}
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return found, examined
}

func TestTaskState188AC3PanelHasNoWriteLeg(t *testing.T) {
	// The real tree. internal/panel is package-adjacent to the cwd the test runs
	// in (internal/tools), so a relative path is enough and nothing here
	// canonicalizes a path by hand (D22 ban, AGENTS §1.2).
	found, examined := scanTaskStateWriteLegs(t, filepath.Join("..", "panel"))
	if examined < 10 {
		t.Fatalf("the AC#3 scan parsed %d non-test panel sources; the denominator is not credible", examined)
	}
	if len(found) != 0 {
		for _, f := range found {
			t.Errorf("panel-side write leg for the task status dimension: %s:%d %s", f.file, f.line, f.why)
		}
	}
	t.Logf("AC#3 negative leg: %d panel non-test sources examined, 0 write legs", examined)

	// Positive control 1: a planted panel-side write leg, one shape at a time.
	// Without these the negative leg above is a door nobody installed.
	plants := map[string]string{
		"roster-record-call": `package panel

import "github.com/CarlosShao/wisp/internal/tools"

func planted(r *tools.TaskRoster, id string, o tools.TaskOutput) {
	r.Record(id, o)
}
`,
		"roster-backfill-call": `package panel

import "github.com/CarlosShao/wisp/internal/tools"

func planted(b tools.TaskBackfill, id, text string) (tools.TaskOutput, string) {
	return b.Backfill(nil, id, text)
}
`,
		"taskoutput-state-key": `package panel

import (
	"github.com/CarlosShao/wisp/internal/statemachine"
	"github.com/CarlosShao/wisp/internal/tools"
)

func planted() tools.TaskOutput {
	return tools.TaskOutput{Text: "正文", State: statemachine.StateActing}
}
`,
		"state-field-assign": `package panel

import (
	"github.com/CarlosShao/wisp/internal/statemachine"
	"github.com/CarlosShao/wisp/internal/tools"
)

func planted(rec tools.TaskOutput) tools.TaskOutput {
	rec.State = statemachine.StateActing
	return rec
}
`,
		"whitelist-method-claim": `package panel

const MethodTaskStateWrite = "panel.task.state.write"
`,
	}
	for name, src := range plants {
		dir := t.TempDir()
		path := filepath.Join(dir, "planted.go")
		if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
			t.Fatalf("plant %s: %v", name, err)
		}
		got, seen := scanTaskStateWriteLegs(t, dir)
		if seen != 1 {
			t.Fatalf("positive control %s: the scan examined %d files, want 1 - it never read the plant", name, seen)
		}
		if len(got) == 0 {
			t.Errorf("positive control %s went GREEN: the AC#3 gate does not catch a panel-side write leg", name)
			continue
		}
		t.Logf("positive control %s red as required: %s:%d %s", name, filepath.Base(got[0].file), got[0].line, got[0].why)
	}

	// Positive control 2 (the reverse one): the SAME scan must stay quiet about a
	// panel-side READER, because display is what the panel is allowed to do
	// (R20 / ticket 92) and ticket 145's pump leg is exactly that shape. A gate
	// that went red here would be a gate aimed at the wrong thing.
	reader := `package panel

import "github.com/CarlosShao/wisp/internal/tools"

type taskView struct {
	Text  string
	Where string
}

func planted(r *tools.TaskRoster, id string) (taskView, bool) {
	rec, ok := r.Look(id)
	if !ok {
		return taskView{}, false
	}
	return taskView{Text: rec.Text, Where: rec.ArtifactPath}, r.Count() >= 0
}
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "reader.go"), []byte(reader), 0o600); err != nil {
		t.Fatalf("plant reader: %v", err)
	}
	got, seen := scanTaskStateWriteLegs(t, dir)
	if seen != 1 {
		t.Fatalf("reader control: examined %d files, want 1", seen)
	}
	if len(got) != 0 {
		for _, f := range got {
			t.Errorf("the AC#3 gate caught a panel-side READER at %s:%d (%s): it is aimed too wide and would put ticket 145's pump leg in its own red",
				f.file, f.line, f.why)
		}
	}
}
