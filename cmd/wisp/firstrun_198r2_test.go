package main

// Ticket 198 r2 - the four instruments 198-v1 (docs/evidence/s1/
// 198-firstrun-config-v1.md) named as missing, one per open cell of its §6-甲:
//
//  1. AC#2 had no teeth: its V3a (a legal NON-default value planted in the
//     table the file is written from, `theme = "light"`) left all six 198 cases
//     green, because every expectation in this package quoted NewDefaults() -
//     the same function that produced the bytes (198-a1 §5.2 warns about that
//     shape by name). Here the expectation side is rebuilt from the
//     `default:"..."` struct tags alone and compared byte-for-byte with the
//     file the production entry actually wrote, so "copied the tags" and "the
//     writer typed a table / the default table was edited" are different
//     readings. It carries its own positive control: eight planted values
//     (string/int/float/bool/nested/slice/plugins-leaf) must each move the
//     render; a blind comparison fails the case instead of passing it.
//  2. AC#5's failure branch had a live reading and no case: V4 (swallow the
//     sentence at run.go:245-247) stayed green everywhere. The case below
//     drives that branch through the real entry over a data root whose secrets
//     layer is blocked, and pins all three halves 198-v1 named: the sentence
//     gets out, nothing half-written stays on disk, and "tried and could not
//     build" still answers 2 rather than success.
//  3. AC#4 said what is missing but not where to fix it, and the fix is not a
//     line of TOML: api_key_ref names a DPAPI blob. The guidance text now says
//     so (firstrun.go), and the case here checks that every `wisp <command>` it
//     names is a command func main really dispatches (read off the package's
//     AST with the same loader the J1 pin uses) - an invented entry point in
//     operator text is the failure this guards, not a keyword spelling.
//  4. J1's landing point was enforced by other tickets' pins: under V1 the only
//     reds were resident_task_source_246_windows_test.go:390 and this ticket's
//     own AST pin, while 198's three-stage case and the logsink case stayed
//     green. The case below asks the behaviour instead: over two identical
//     fresh data roots, the assembly root the resident leg calls directly
//     (resident_task_source_windows.go:265) must create nothing, and the
//     user-started run entry must create the file. Moving the call into
//     assembleRuntime reddens this ticket's own case, not a neighbour's.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and
// dies at load (0xc0000135, with no --- FAIL line) unless
// third_party/sherpa-onnx is on PATH.

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
)

// ---------------------------------------------------------------------------
// AC#2 - the expectation side, rebuilt from the tags
// ---------------------------------------------------------------------------

// tagDefaults198r2 rebuilds the default table from the `default:"..."` struct
// tags and returns it with the number of leaves it set. It deliberately does
// NOT call config.NewDefaults: that call is the thing under test here, and
// quoting it as the expectation is the vacuity 198-a1 §5.2 names.
func tagDefaults198r2(t *testing.T) (*config.Config, int) {
	t.Helper()
	out := &config.Config{}
	n := setTaggedLeaves198r2(t, reflect.ValueOf(out).Elem(), "", 0)
	if n == 0 {
		t.Fatal("the tag reader set 0 leaves, so every comparison built on it proves nothing")
	}
	return out, n
}

// setTaggedLeaves198r2 walks a struct value and applies every default tag it
// finds, mirroring the schema's own documented conventions: maps stay nil
// (defaults.go:56-58 - go-toml merges into pre-existing maps), `toml:"-"`
// struct fields are not emitted by the head marshal except [plugins], which
// marshalPlugins appends by hand (parse.go:181-191) and is therefore walked
// under that name. Shapes this walker does not model are a hard stop, never a
// silent skip: a silent skip is how a check like this one turns into a formality.
func setTaggedLeaves198r2(t *testing.T, v reflect.Value, prefix string, depth int) int {
	t.Helper()
	if depth > 12 {
		t.Fatalf("config struct nesting past depth 12 at %q: this walk is lost, extend it deliberately", prefix)
	}
	set := 0
	typ := v.Type()
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		fv := v.Field(i)
		tag := f.Tag.Get("toml")
		parts := strings.Split(tag, ",")
		name := parts[0]
		if len(parts) > 1 {
			t.Fatalf("%s.%s carries toml tag options (%q) this instrument does not model; teach it on purpose instead of letting the comparison quietly cover less",
				typ.Name(), f.Name, tag)
		}
		def, hasDef := f.Tag.Lookup("default")
		switch fv.Kind() {
		case reflect.Struct:
			switch {
			case name == "-" && fv.Type() == reflect.TypeOf(config.PluginsSection{}):
				set += setTaggedLeaves198r2(t, fv, "plugins", depth+1)
			case name == "-":
				// Not emitted at all: nothing to compare, nothing to set.
			case name == "":
				t.Fatalf("%s.%s is an untagged struct: this walk cannot name the table it belongs to",
					typ.Name(), f.Name)
			default:
				set += setTaggedLeaves198r2(t, fv, prefix+"."+name, depth+1)
			}
		case reflect.Map:
			if hasDef {
				t.Fatalf("%s.%s is a map carrying a default tag; pre-populating one is exactly the invented entry AC#2 forbids (defaults.go:56-58)",
					typ.Name(), f.Name)
			}
		case reflect.Pointer, reflect.Interface:
			t.Fatalf("%s.%s is a %s: this instrument does not model it, and a default it cannot see is a default it cannot guard",
				typ.Name(), f.Name, fv.Kind())
		default:
			if !hasDef {
				continue
			}
			if err := setLeaf198r2(fv, def); err != nil {
				t.Fatalf("%s.%s: default tag %q cannot be read back: %v", typ.Name(), f.Name, def, err)
			}
			set++
		}
	}
	return set
}

// setLeaf198r2 parses one tag into one leaf. Slice defaults are
// comma-separated by the schema's own convention (defaults.go:91-92), which is
// also why a default containing a comma cannot be expressed - on either side.
func setLeaf198r2(fv reflect.Value, def string) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(def)
	case reflect.Bool:
		b, err := strconv.ParseBool(def)
		if err != nil {
			return err
		}
		fv.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(def, 10, 64)
		if err != nil {
			return err
		}
		fv.SetInt(n)
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(def, 64)
		if err != nil {
			return err
		}
		fv.SetFloat(x)
	case reflect.Slice:
		if def == "" {
			return nil
		}
		items := strings.Split(def, ",")
		out := reflect.MakeSlice(fv.Type(), 0, len(items))
		for _, item := range items {
			ev := reflect.New(fv.Type().Elem()).Elem()
			if err := setLeaf198r2(ev, strings.TrimSpace(item)); err != nil {
				return err
			}
			out = reflect.Append(out, ev)
		}
		fv.Set(out)
	default:
		return errors.New("unsupported default kind " + fv.Kind().String())
	}
	return nil
}

// render198r2 serialises a Config the way SaveFile does: stamp the current
// schema version (loader.go:242-243, the one key in the file that is a
// migration marker rather than a default) and run the canonical writer. Both
// sides of every comparison below go through this same function, so a
// difference here can only be a difference in VALUES, which is precisely what
// AC#2 asks about.
func render198r2(t *testing.T, c *config.Config) []byte {
	t.Helper()
	c.SchemaVersion = config.SchemaVersionCurrent
	data, err := config.MarshalCanonical(c)
	if err != nil {
		t.Fatalf("MarshalCanonical of the expectation: %v", err)
	}
	return data
}

// firstByteDiff198r2 names the first differing line so a red says which key
// drifted instead of dumping two 2.5 KiLobyte files into the log.
func firstByteDiff198r2(want, got []byte) string {
	i := 0
	for i < len(want) && i < len(got) && want[i] == got[i] {
		i++
	}
	wl := strings.Split(string(want), "\n")
	gl := strings.Split(string(got), "\n")
	line := 1 + bytes.Count(want[:min(i, len(want))], []byte("\n"))
	w, g := "<eof>", "<eof>"
	if line-1 < len(wl) {
		w = wl[line-1]
	}
	if line-1 < len(gl) {
		g = gl[line-1]
	}
	return "first difference at line " + strconv.Itoa(line) +
		" (byte " + strconv.Itoa(i) + "): want " + strconv.Quote(w) +
		", got " + strconv.Quote(g) +
		"; total bytes want " + strconv.Itoa(len(want)) + " got " + strconv.Itoa(len(got))
}

// TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags is the
// instrument 198-v1 §6-甲-1 asked for: the file the production entry wrote,
// against the tags rendered independently of NewDefaults.
func TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, configFileName)

	expect, leaves := tagDefaults198r2(t)
	want := render198r2(t, expect)

	// The artifact under test is the one the real entry wrote, not a
	// SaveFile call this test makes itself.
	if code, _, stderrText := run198(t, dir); code != 2 {
		t.Fatalf("exit %d, want 2 (a created config still has no model and no key):\n%s", code, stderrText)
	}
	got, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("first-run wrote no file to compare: %v", err)
	}
	t.Logf("tag-derived default leaves: %d; file: %d bytes", leaves, len(got))
	if !bytes.Equal(want, got) {
		t.Errorf("AC#2 RED: the config.toml first-run wrote is not the schema's `default` tags rendered: %s\n"+
			"every value in that file must be traceable to a default tag (D36 rule 3) - a value the writer typed, "+
			"or a default table edited away from its tags, is exactly what this case exists to catch",
			firstByteDiff198r2(want, got))
	}

	// 正控（种下必红）：每一类 kind 都单独种一枚合法的非默认值。渲染必须因此
	// 与盘上文件不同；哪一种种下去没差别，说明那一种根本进不了比较，这条尺
	// 对它是瞎的 - 那要在这里红，而不是等下一枚变异把整个 AC#2 判成无牙。
	for _, ctrl := range []struct {
		name   string
		seeded func(*config.Config)
	}{
		{"app.theme 改成合法非默认（198-v1 的 V3a 那一形）", func(c *config.Config) { c.App.Theme = config.ThemeLight }},
		{"app.language 改成另一个合法值", func(c *config.Config) { c.App.Language = "en-US" }},
		{"ball.size（Int）", func(c *config.Config) { c.Ball.Size = 60 }},
		{"ball.opacity_idle（Float64）", func(c *config.Config) { c.Ball.OpacityIdle = 0.9 }},
		{"app.autostart（Bool）", func(c *config.Config) { c.App.Autostart = true }},
		{"llm.roles.chat.temperature（嵌套段）", func(c *config.Config) { c.LLM.Roles.Chat.Temperature = 0.2 }},
		{"voice.wake_word.veto_words（Slice 默认）", func(c *config.Config) { c.Voice.WakeWord.VetoWords = []string{"停"} }},
		{"plugins.tier2_enabled（marshalPlugins 那一支）", func(c *config.Config) { c.Plugins.Tier2Enabled = true }},
	} {
		seeded, _ := tagDefaults198r2(t)
		ctrl.seeded(seeded)
		if bytes.Equal(render198r2(t, seeded), got) {
			t.Errorf("AC#2 仪器自证失败：种下「%s」之后渲染仍与盘上文件逐字节相同 ⇒ 这一类值进不了比较", ctrl.name)
		}
	}
}

// TestTicket198R2AC2ExportedDefaultTableStillRendersAsItsTags is the second
// reading: NewDefaults() itself against the same tag-derived expectation. The
// two cases together name WHERE a drift lives - only this one reddens when
// internal/config's own applier changes, only the case above reddens when the
// CLI writes something other than the defaults.
func TestTicket198R2AC2ExportedDefaultTableStillRendersAsItsTags(t *testing.T) {
	expect, _ := tagDefaults198r2(t)
	want := render198r2(t, expect)
	got := render198r2(t, config.NewDefaults())
	if !bytes.Equal(want, got) {
		t.Errorf("AC#2 RED: internal/config's default table is no longer just its `default` tags: %s",
			firstByteDiff198r2(want, got))
	}
}

// ---------------------------------------------------------------------------
// AC#5 - the creation-failure branch
// ---------------------------------------------------------------------------

// TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile guards the branch
// 198-v1 measured live but could not find a case for (its V4, swallowing
// run.go:245-247 into `_ = frErr`, was green everywhere).
//
// The blocked shape is 198-v1 §1-AC#5's "数据根不可建", with the ACL replaced
// by a plain file parked at <dataDir>\secrets: winsec's sealing walk stops
// there with "is not a directory" (winsec.go:186-188), secret.NewStore wraps it
// (store.go:49-51), and ensureFirstRunConfig returns the wrapped failure
// (firstrun.go:79-81). A denied ACL would be the exact live shape but cannot be
// removed by t.TempDir - ticket 101's note at run_mode101_test.go:621-633 books
// that same cost.
func TestTicket198R2AC5CreationFailureIsLoudAndLeavesNoHalfFile(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "secrets"), []byte("standing in the way of the store"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(dir, configFileName)
	if _, err := os.Stat(cfgPath); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("premise broke: %q already exists (err %v)", cfgPath, err)
	}

	code, stdoutText, stderrText := run198(t, dir)

	// 1. the failure gets out, loudly, on the operator's stream.
	if !strings.Contains(stderrText, "首次配置建不出来") {
		t.Errorf("AC#5 RED: the run could not create its first config and never said so (run.go:245-247 is the sentence this case exists for). stderr:\n%s", stderrText)
	}
	if strings.Contains(stdoutText, "首次配置建不出来") {
		t.Errorf("the creation failure was streamed to stdout instead of stderr:\n%s", stdoutText)
	}
	// It names the reason, not just the verdict: the data root that refused.
	for _, mustName := range []string{"数据根不可建", dir} {
		if !strings.Contains(stderrText, mustName) {
			t.Errorf("AC#5: the failure sentence does not name %q, so the user cannot tell which root refused:\n%s", mustName, stderrText)
		}
	}
	// 2. nothing half-written is left behind.
	if _, err := os.Stat(cfgPath); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("AC#5 RED: a creation that failed still left a config.toml at %q (err %v) - a half file reported as success is the shape the ticket forbids", cfgPath, err)
	}
	if strays, _ := filepath.Glob(filepath.Join(dir, ".wisp-config-*.tmp")); len(strays) > 0 {
		t.Errorf("atomicWrite left temp files behind after the failed creation: %v", strays)
	}
	// 3. "tried and could not build" is never a success.
	if code != 2 {
		t.Errorf("AC#5: exit %d, want 2 (a failed first-run must still answer the classified missing-config verdict, not 0)\nstdout:\n%s\nstderr:\n%s",
			code, stdoutText, stderrText)
	}
	// And the run never claims it created the defaults.
	if strings.Contains(stderrText, "新建默认配置") {
		t.Errorf("AC#5 RED: the creation receipt fired on a run that created nothing:\n%s", stderrText)
	}
}

// ---------------------------------------------------------------------------
// AC#4 - "where do I fix this", stated as the path the user can actually walk
// ---------------------------------------------------------------------------

// wispCommand198r2 picks the command word out of operator text like
// "wisp secret set <name>" / "wisp run: ...".
var wispCommand198r2 = regexp.MustCompile(`\bwisp ([a-z][a-z0-9-]*)`)

// TestTicket198R2AC4ReceiptNamesTheRealEntryPoints pins the half 198-v1 found
// missing: the receipt said what is missing but not where to补 it, and the
// answer is not "write another line of TOML" - api_key_ref names a DPAPI blob
// (secret.go:380 prints the reference `wisp secret set` just stored; store.go
// has no plaintext path, and schema.go:384-387 forbids a plaintext field).
//
// The honesty check is a capability check, not a keyword match: every command
// the text tells the user to run must be one func main dispatches, read off the
// package's own AST. Operator text naming an entry point that does not exist is
// the failure mode this guards.
func TestTicket198R2AC4ReceiptNamesTheRealEntryPoints(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, configFileName)
	code, _, stderrText := run198(t, dir)
	if code != 2 {
		t.Fatalf("exit %d, want 2, stderr:\n%s", code, stderrText)
	}
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("no config.toml was created, so there is no receipt to read: %v", err)
	}

	for _, mustName := range []string{
		"wisp secret set", // the only credential entry point that exists
		"api_key_ref",     // the field the reference belongs to
		"dpapi:",          // what the ref actually names (a blob, not a key)
		"env:",            // the one alternative, and where its value lives
		"[llm.providers.", // where the model catalog is filled in
		"text_chain",      // the pair the run resolves against
		"roles.chat",      // ... and the role form of the same answer
	} {
		if !strings.Contains(stderrText, mustName) {
			t.Errorf("AC#4: the first-run receipt never names %q, so it still says what is missing without saying where to fix it:\n%s", mustName, stderrText)
		}
	}

	// Every command the receipt tells the user to run must really dispatch.
	pkg, err := loadMainPackage131(".")
	if err != nil {
		t.Fatalf("loading package main: %v", err)
	}
	var invented []string
	seen := map[string]bool{}
	for _, m := range wispCommand198r2.FindAllStringSubmatch(stderrText, -1) {
		word := m[1]
		if seen[word] {
			continue
		}
		seen[word] = true
		fn := "cmd" + strings.ToUpper(word[:1]) + word[1:]
		if _, ok := pkg.funcs[fn]; !ok {
			invented = append(invented, `wisp `+word+` (no `+fn+` in package main)`)
		}
	}
	if len(invented) > 0 {
		t.Errorf("AC#4 RED: the guidance sends the user to entry points this build does not dispatch: %v", invented)
	}

	// Non-vacuity: the guidance rides the creation branch only. A second run
	// over the same root must repeat neither the receipt nor its fix-up lines -
	// were this assertion satisfied by some other leg's text, the check above
	// would pass on a boot that never created anything.
	_, _, stderr2 := run198(t, dir)
	if strings.Contains(stderr2, "新建默认配置") || strings.Contains(stderr2, "wisp secret set") {
		t.Errorf("the creation guidance fired a second time on a root that already has a config.toml:\n%s", stderr2)
	}
}

// ---------------------------------------------------------------------------
// J1 - this ticket's own landing-point instrument
// ---------------------------------------------------------------------------

// TestTicket198R2J1TheAssemblyRootStillCreatesNothing is 198-v1 §3-V1's other
// half: under that mutation the reds belonged to ticket 246's pin and to this
// ticket's AST reading, while the behavioural cases stayed green. So the
// behaviour is asserted here, over two identical fresh roots: the assembly root
// (which the resident leg reaches directly, resident_task_source_windows.go:265)
// must create nothing, and the user-started run entry must create the file. The
// asymmetry is the ruling; either half flattening into the other reddens a
// case this ticket owns.
func TestTicket198R2J1TheAssemblyRootStillCreatesNothing(t *testing.T) {
	// (a) the shared assembly root, the way the resident leg calls it.
	rootDir := t.TempDir()
	out, errb := &bytes.Buffer{}, &bytes.Buffer{}
	rt, code := assembleRuntime(runSpec{
		stdout:  out,
		stderr:  errb,
		dataDir: rootDir,
		notify:  func(string, string) error { return nil },
	})
	if rt != nil {
		defer rt.close()
	}
	if code != 2 {
		t.Fatalf("assembleRuntime over a fresh root: exit %d, want 2 (Unconfigured)\nstderr:\n%s", code, errb.String())
	}
	if _, err := os.Stat(filepath.Join(rootDir, configFileName)); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("J1 RED: the assembly root created a config.toml nobody asked it for (stat err %v) - the resident leg reaches this function directly, and ticket 246's pin is the only other thing standing between that and a written root",
			err)
	}
	if strings.Contains(errb.String(), "新建默认配置") {
		t.Errorf("J1 RED: the creation receipt fired inside the assembly root:\n%s", errb.String())
	}
	if strays, _ := filepath.Glob(filepath.Join(rootDir, ".wisp-config-*.tmp")); len(strays) > 0 {
		t.Errorf("the assembly root left config temp files in a data dir it must not write: %v", strays)
	}

	// (b) the user-started entry over the same shape of root must still create
	// it - otherwise (a) would be green for the wrong reason (creation broken
	// everywhere) and this case would pin nothing.
	runDir := t.TempDir()
	runCode, _, runStderr := run198(t, runDir)
	if runCode != 2 {
		t.Fatalf("runTextTask over a fresh root: exit %d, want 2, stderr:\n%s", runCode, runStderr)
	}
	if _, err := os.Stat(filepath.Join(runDir, configFileName)); err != nil {
		t.Fatalf("J1 non-vacuity broke: the run entry created no config.toml either, so half (a) proves nothing: %v", err)
	}
	if !strings.Contains(runStderr, "新建默认配置") {
		t.Errorf("the run entry created the file without saying so:\n%s", runStderr)
	}
}
