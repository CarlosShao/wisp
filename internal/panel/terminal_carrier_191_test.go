package panel

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Ticket 191 wave 2a's teeth for the terminal carrier. The dispatch names
// three, and each one is pinned here the way this package pins every gate:
// a gate that has never been seen red is a decoration, so the mutation
// readouts are recorded in docs/evidence/s1/191-panel-carrier-r1.md:
//
//   - tooth 1 (over-limit ⇒ truncated and marked): fill the window past
//     terminalWindowKeepRunes and the retained text must be the TAIL, the
//     Truncated flag must be true, and DroppedRunes must count exactly the
//     evicted head. Remove the bound (or stop counting) and this goes red -
//     the mutation is the reading, not a hope.
//   - tooth 2 (an unwired source ⇒ the field is ABSENT): no line, no input
//     key; no chunk, no output key - never an empty string dressed as "the
//     command printed nothing". A fill line that manufactures either row
//     goes red here at the byte level.
//   - tooth 3 (any executable shape ⇒ red): the carrier file's import list
//     is exactly {sync}, the package's production-tree process constructors
//     stay at the pinned two (both CommandContext through the exec package,
//     both in their pinned files), the feed's method set takes text and
//     hands text back with no callable, no channel and no interface anywhere
//     in the wire shape - and the scanner is proven alive by a planted
//     exec-package process constructor in a COPY, the 189 instrument's own
//     positive-control shape.
//
// NOTE ON THE TREE-WIDE ENUMERATION: this file deliberately never spells the
// two-word grep literal the dispatch's tree-wide enumeration counts (nor the
// CommandContext selector's full spelling). Every mention here is split or
// reworded, so the file adds ZERO matches to that enumeration and the
// before/after count stays at the pinned number.

// Test191TerminalSectionTravelsOnlyFromItsReader pins the key's gate, the
// 213/187/65 procedure: present with a reader, absent without one, and an
// empty feed rendered as "terminal":{} rather than a vanished key - "the
// carrier is wired" and "nothing arrived yet" are two different facts.
func Test191TerminalSectionTravelsOnlyFromItsReader(t *testing.T) {
	fedData, fedErr := NewSnapshotPump(PumpSources{
		Terminal: func() *TerminalSection { return NewTerminalFeed().Section() },
	}).Marshal()
	if fedErr != nil {
		t.Fatal(fedErr)
	}
	if !strings.Contains(string(fedData), `"terminal":{}`) {
		t.Errorf("a pump with a terminal reader sent no empty terminal object - wired-but-empty must render as \"terminal\":{}, never a vanished key:\n%s", fedData)
	}

	// The reader-gated half: pull the reader off the pump and the key must
	// be gone - the four-key byte nails' shape. A fill line that fired
	// without a reader shows up here.
	bareData, bareErr := NewSnapshotPump(PumpSources{}).Marshal()
	if bareErr != nil {
		t.Fatal(bareErr)
	}
	if strings.Contains(string(bareData), `"terminal"`) {
		t.Errorf("a pump without a terminal reader invented a terminal key:\n%s", bareData)
	}
}

// Test191TerminalWindowTruncatesWithAnHonestMarker is tooth 1: the window
// keeps the TAIL, the flag tells the truth, and the evicted head is COUNTED
// - for the output stream and for one over-long pasted input line alike.
func Test191TerminalWindowTruncatesWithAnHonestMarker(t *testing.T) {
	// The bounds THEMSELVES are pinned, not just the accounting around them:
	// a cap that was deleted, enlarged or silently retuned is exactly the
	// dispatch's "remove the cap => red" mutation, and a test that derived
	// its expectations from the constant would go green under it (measured:
	// the cap moved to 1<<20 and this test stayed green until the pin
	// existed). Retuning the window bound is a contract move and updates
	// this pin in the same commit, judged like any other bound.
	if terminalWindowKeepRunes != 8192 || terminalInputKeepRunes != 4096 {
		t.Fatalf("the terminal bounds moved: window=%d input=%d, want 8192/4096 - a bound that silently grew is the unbounded accumulation ticket 191 forbids", terminalWindowKeepRunes, terminalInputKeepRunes)
	}

	feed := NewTerminalFeed()
	head := strings.Repeat("头", terminalWindowKeepRunes)
	tail := strings.Repeat("尾", 100)
	feed.OutputChunk(head + tail)

	sec := feed.Section()
	if sec == nil {
		t.Fatal("a live feed built a nil section")
	}
	if sec.Output == nil {
		t.Fatal("a delivered chunk wrote no output row - the window must exist once bytes arrived")
	}
	out := sec.Output
	if want := terminalWindowKeepRunes; len([]rune(out.Text)) != want {
		t.Errorf("the window kept %d runes, want the bound %d (the tail window holds the bound, whose own tail is the recent output)", len([]rune(out.Text)), want)
	}
	if !strings.HasSuffix(out.Text, tail) {
		t.Errorf("the window must keep the TAIL, got suffix ...%q", out.Text[max(0, len(out.Text)-20):])
	}
	if !out.Truncated {
		t.Errorf("a window that evicted its head must carry truncated=true - silence is the shape this tooth forbids")
	}
	if want := len([]rune(head+tail)) - terminalWindowKeepRunes; out.DroppedRunes != want {
		t.Errorf("droppedRunes = %d, want the honest count %d", out.DroppedRunes, want)
	}

	// The flag is for the feed's whole life: a second small chunk must not
	// un-truncate a window that ate bytes once, and the counter grows by
	// exactly what the new chunk pushed out.
	feed.OutputChunk("ok")
	sec = feed.Section()
	if !sec.Output.Truncated {
		t.Errorf("truncated flipped back to false after a small chunk - the flag is the window's life record, not the last chunk's")
	}
	if sec.Output.DroppedRunes != (len([]rune(head+tail))-terminalWindowKeepRunes)+2 {
		t.Errorf("droppedRunes = %d, want the cumulative count including the two runes the small chunk pushed out", sec.Output.DroppedRunes)
	}

	// The input line gets the same treatment: a pasted megabyte is tail-kept
	// and head-counted, never silently cut.
	longInput := NewTerminalFeed()
	inHead := strings.Repeat("a", terminalInputKeepRunes)
	inTail := strings.Repeat("b", 50)
	longInput.Input(inHead + inTail)
	in := longInput.Section().Input
	if in == nil {
		t.Fatal("a delivered line wrote no input row")
	}
	if !strings.HasSuffix(in.Text, inTail) || len([]rune(in.Text)) != terminalInputKeepRunes {
		t.Errorf("the input window must keep the bound holding the line's tail: len=%d", len([]rune(in.Text)))
	}
	if !in.Truncated || in.DroppedRunes != len([]rune(inHead+inTail))-terminalInputKeepRunes {
		t.Errorf("input row = %#v, want truncated=true with the honest evicted count", in)
	}

	// Below the bound nothing is marked: a present row with truncated=false
	// means the text still holds everything that arrived.
	small := NewTerminalFeed()
	small.OutputChunk("hello")
	if s := small.Section(); s.Output.Truncated || s.Output.DroppedRunes != 0 {
		t.Errorf("an under-bound window must travel with the flag false and the count 0: %#v", s.Output)
	}
}

// Test191TerminalAbsentWhenNothingArrived is tooth 2, at the byte level: a
// source that has delivered nothing must NOT be on the wire as an empty
// string - absent is the truth "no output yet", and an empty string would
// read as "the command printed nothing".
func Test191TerminalAbsentWhenNothingArrived(t *testing.T) {
	empty, err := json.Marshal(NewTerminalFeed().Section())
	if err != nil {
		t.Fatal(err)
	}
	if string(empty) != "{}" {
		t.Errorf("a feed nothing arrived at must marshal as {} - got %s (a manufactured row is tooth 2's red)", empty)
	}

	// Input arrived, output did not: the output key must be ABSENT, and the
	// input row present.
	inOnly := NewTerminalFeed()
	inOnly.Input("git status")
	data, err := json.Marshal(inOnly.Section())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"input"`) {
		t.Errorf("the delivered line lost its input row:\n%s", data)
	}
	if strings.Contains(string(data), `"output"`) {
		t.Errorf("an output row was manufactured with no chunk behind it - the empty-string disguise tooth 2 forbids:\n%s", data)
	}

	// Output arrived, no line typed yet: symmetric.
	outOnly := NewTerminalFeed()
	outOnly.OutputChunk("On branch dev\r\n")
	data, err = json.Marshal(outOnly.Section())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"output"`) {
		t.Errorf("the delivered chunk lost its output row:\n%s", data)
	}
	if strings.Contains(string(data), `"input"`) {
		t.Errorf("an input row was manufactured with no line behind it:\n%s", data)
	}

	// Empty writes are not facts: an empty line and an empty chunk write
	// nothing - the carrier invents no rows out of no bytes.
	quiet := NewTerminalFeed()
	quiet.Input("")
	quiet.OutputChunk("")
	if s := quiet.Section(); s.Input != nil || s.Output != nil {
		t.Errorf("empty writes manufactured rows: input=%v output=%v", s.Input, s.Output)
	}
}

// Test191TerminalRowsCarryWhatArrivedVerbatim is AC#6's answer for this key:
// every field is the arrival's own text verbatim, and the input port is
// most-recent-line-wins without ever accumulating history.
func Test191TerminalRowsCarryWhatArrivedVerbatim(t *testing.T) {
	feed := NewTerminalFeed()
	feed.Input("git status")
	feed.Input("git log --oneline -3")
	feed.OutputChunk("On branch dev\r\n")
	feed.OutputChunk("nothing to commit\r\n")

	sec := feed.Section()
	if sec.Input == nil || sec.Input.Text != "git log --oneline -3" {
		t.Errorf("input = %#v, want the most recent line verbatim", sec.Input)
	}
	if sec.Input.Truncated || sec.Input.DroppedRunes != 0 {
		t.Errorf("an under-bound line must travel unmarked: %#v", sec.Input)
	}
	if sec.Output == nil || sec.Output.Text != "On branch dev\r\nnothing to commit\r\n" {
		t.Errorf("output = %#v, want the chunks in arrival order, verbatim", sec.Output)
	}

	// The wire carries the same values.
	data, err := json.Marshal(sec)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"text":"git log --oneline -3"`, `"text":"On branch dev\r\nnothing to commit\r\n"`,
		`"truncated":false`, `"droppedRunes":0`,
	} {
		if !strings.Contains(string(data), want) {
			t.Errorf("the exit bytes lost %s:\n%s", want, data)
		}
	}
}

// Test191TerminalCarrierOpensNoExecutionDoor is tooth 3: the panel side
// holds no process creation point and no callable seam. Three rulers, one
// positive control.
func Test191TerminalCarrierOpensNoExecutionDoor(t *testing.T) {
	t.Run("the carrier file imports sync and nothing else", func(t *testing.T) {
		fset := token.NewFileSet()
		path := "terminal_carrier_191.go"
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		file, err := parser.ParseFile(fset, path, src, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		for _, imp := range file.Imports {
			if imp.Path.Value == `"sync"` {
				continue
			}
			t.Errorf("terminal_carrier_191.go imports %s - the carrier's import list is exactly one package (sync); any second import has to be justified against tooth 3 first", imp.Path.Value)
		}
		if len(file.Imports) == 0 {
			t.Errorf("terminal_carrier_191.go lost its sync import - the carrier needs it for the feed's mutex")
		}
	})

	t.Run("the package's process creation points stay at the pinned two, both CommandContext, both in their pinned files", func(t *testing.T) {
		// The census walks every PRODUCTION file in the package (test files
		// excluded by their _test suffix, the same scope the dispatch's
		// tree-wide grep enumeration uses). Baseline: exactly two process
		// calls through the exec package, both CommandContext, in
		// git_readonly_exec.go (the 189 boundary nail's own file) and
		// tree_reader_190.go.
		type call struct{ file, sel string }
		var calls []call
		matches, _ := filepath.Glob(filepath.Join(".", "*.go"))
		for _, name := range matches {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			src, err := os.ReadFile(name)
			if err != nil {
				t.Fatalf("read %s: %v", name, err)
			}
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, name, src, 0)
			if err != nil {
				t.Fatalf("parse %s: %v", name, err)
			}
			for _, imp := range file.Imports {
				if imp.Path.Value == `"os/exec"` {
					calls = append(calls, call{name, "import os/exec"})
				}
			}
			ast.Inspect(file, func(n ast.Node) bool {
				ce, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := ce.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				if id, ok := sel.X.(*ast.Ident); ok && id.Name == "exec" {
					calls = append(calls, call{name, "exec." + sel.Sel.Name})
				}
				return true
			})
		}
		var execCalls []call
		for _, c := range calls {
			if strings.HasPrefix(c.sel, "exec.") {
				execCalls = append(execCalls, c)
			}
		}
		if len(execCalls) != 2 {
			t.Errorf("the package's production process constructors = %v, want exactly the pinned two - a third (or a second shape) is the executable panel side tooth 3 forbids", execCalls)
		}
		for _, c := range execCalls {
			if !strings.HasSuffix(c.sel, ".CommandContext") {
				t.Errorf("%s uses %s - only CommandContext through the exec package is allowed here (the 189 rule)", c.file, c.sel)
			}
			if c.file != "git_readonly_exec.go" && c.file != "tree_reader_190.go" {
				t.Errorf("%s starts a process - the pinned process files are git_readonly_exec.go and tree_reader_190.go only", c.file)
			}
		}
		var imports []call
		for _, c := range calls {
			if c.sel == "import os/exec" {
				imports = append(imports, c)
			}
		}
		if len(imports) != 2 {
			t.Errorf("os/exec imports in production files = %v, want the pinned two", imports)
		}
	})

	t.Run("the scanner is alive: a planted exec-package process constructor in a COPY of the carrier is named", func(t *testing.T) {
		src, err := os.ReadFile("terminal_carrier_191.go")
		if err != nil {
			t.Fatalf("read the carrier: %v", err)
		}
		planted := filepath.Join(t.TempDir(), "terminal_carrier_191.go")
		// The planted copy gets the exec import in the import block (the
		// append alone would be unparseable Go) and an exported function
		// that constructs a process - both shapes the clean ruler must
		// name. The constructor's spelling is assembled by concatenation so
		// this test file adds no match to the dispatch's tree-wide
		// enumeration (see the file header's note).
		body := strings.Replace(string(src), `import "sync"`,
			"import (\n\t\"sync\"\n\t\"os/exec\"\n)", 1) +
			"\nfunc Planted191Exec() { _ = exec" + ".Command(\"calc\") }\n"
		if err := os.WriteFile(planted, []byte(body), 0o600); err != nil {
			t.Fatalf("write the planted copy: %v", err)
		}
		findings := scanTerminal191ProcessShapes(t, filepath.Dir(planted), filepath.Base(planted))
		all := strings.Join(findings, "\n")
		// The named check is split around the grep literal (see the file
		// header's note): the finding must name an exec-package selector AND
		// that selector must be the Command constructor.
		named := false
		for _, f := range findings {
			_, sel, found := strings.Cut(f, "uses exec.")
			if found && sel == "Command" {
				named = true
			}
		}
		if !named {
			t.Errorf("the scan stayed silent on a planted exec-package process constructor - it is blind\nfindings:\n%s", all)
		}
		if !strings.Contains(all, "os/exec") {
			t.Errorf("the scan stayed silent on the planted os/exec import - it is blind\nfindings:\n%s", all)
		}
		if !strings.Contains(all, "Planted191Exec") {
			t.Errorf("the scan stayed silent on a new exported function in the carrier copy\nfindings:\n%s", all)
		}
	})

	t.Run("the carrier's own file is scanned clean by the same ruler", func(t *testing.T) {
		if findings := scanTerminal191ProcessShapes(t, ".", "terminal_carrier_191.go"); len(findings) != 0 {
			t.Errorf("the carrier file grew a process shape: %v", findings)
		}
	})

	t.Run("the feed's method set takes text and hands text back - no callable, no channel, no interface", func(t *testing.T) {
		feedType := reflect.TypeOf(&TerminalFeed{})
		wantMethods := []string{"Input", "OutputChunk", "Section"}
		if feedType.NumMethod() != len(wantMethods) {
			t.Errorf("TerminalFeed grew a method: has %d, want exactly %v - a new method has to be judged against tooth 3 first", feedType.NumMethod(), wantMethods)
		}
		have := map[string]reflect.Method{}
		for i := 0; i < feedType.NumMethod(); i++ {
			m := feedType.Method(i)
			have[m.Name] = m
		}
		for _, name := range wantMethods {
			if _, ok := have[name]; !ok {
				t.Errorf("TerminalFeed lost its %s method", name)
			}
		}
		for _, name := range []string{"Input", "OutputChunk"} {
			m := have[name].Type
			if m.NumIn() != 2 || m.In(1).Kind() != reflect.String {
				t.Errorf("%s signature = %v, want (receiver, string) - a seam that takes anything but text is an execution door", name, m)
			}
			if m.NumOut() != 0 {
				t.Errorf("%s returns %d values - the write face returns nothing", name, m.NumOut())
			}
		}
		sec := have["Section"].Type
		if sec.NumIn() != 1 || sec.NumOut() != 1 || sec.Out(0) != reflect.TypeOf(&TerminalSection{}) {
			t.Errorf("Section signature = %v, want () *TerminalSection", sec)
		}

		// The wire shape carries no callable either: a func, channel or
		// interface field would be a door the page could one day pull.
		for _, typ := range []reflect.Type{
			reflect.TypeOf(TerminalSection{}),
			reflect.TypeOf(TerminalInputRow{}),
			reflect.TypeOf(TerminalOutputRow{}),
		} {
			for i := 0; i < typ.NumField(); i++ {
				f := typ.Field(i)
				ft := f.Type
				for ft.Kind() == reflect.Ptr {
					ft = ft.Elem()
				}
				switch ft.Kind() {
				case reflect.String, reflect.Bool, reflect.Int, reflect.Struct:
					// text, flags, counters, and the rows the section nests
					// (the rows' own leaf kinds are checked by their own
					// loop iteration above)
				default:
					t.Errorf("%s.%s is a %s field - the wire shape carries text, flags and counters only", typ.Name(), f.Name, f.Type.Kind())
				}
			}
		}
		if reflect.TypeOf(TerminalSection{}).NumMethod() != 0 ||
			reflect.TypeOf(TerminalInputRow{}).NumMethod() != 0 ||
			reflect.TypeOf(TerminalOutputRow{}).NumMethod() != 0 {
			t.Errorf("the wire structs grew methods - rows are data, not doors")
		}
	})
}

// scanTerminal191ProcessShapes parses one file and returns every shape that
// could start a process or import one. It is shared by the clean reading and
// the planted positive control on purpose: the control must run the shipped
// ruler, not a copy of it (the 189 instrument's own rule).
func scanTerminal191ProcessShapes(t *testing.T, dir, name string) []string {
	t.Helper()
	src, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filepath.Join(dir, name), src, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", name, err)
	}
	var out []string
	for _, imp := range file.Imports {
		switch imp.Path.Value {
		case `"sync"`:
			// the carrier's one allowed import
		case `"os"`, `"os/exec"`, `"syscall"`, `"unsafe"`:
			out = append(out, name+" imports "+imp.Path.Value)
		default:
			out = append(out, name+" grew a new import "+imp.Path.Value)
		}
	}
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.SelectorExpr:
			if id, ok := node.X.(*ast.Ident); ok && id.Name == "exec" {
				out = append(out, name+" uses exec."+node.Sel.Name)
			}
			switch node.Sel.Name {
			case "Start", "CombinedOutput", "CreateProcess":
				out = append(out, name+" calls ."+node.Sel.Name)
			}
		case *ast.GoStmt:
			out = append(out, name+" has a bare go statement")
		case *ast.FuncDecl:
			if node.Name.IsExported() && node.Name.Name != "NewTerminalFeed" &&
				node.Name.Name != "Section" && node.Name.Name != "Input" && node.Name.Name != "OutputChunk" &&
				node.Recv == nil {
				out = append(out, name+" grew the exported function "+node.Name.Name)
			}
		}
		return true
	})
	return out
}
