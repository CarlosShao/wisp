package main

// Ticket 198 AC#1/AC#3 in the shape 普查件 198-a1 §5.1 demanded, plus the two
// pins that keep ruling J1 (账 A514) from rotting: the receipt sentence stays
// distinct from 票 223's cause=missing line, and the creation caller set is
// exactly the user-started run entry.
//
// The three-stage AC#1 case is deliberate about what it does NOT prove: the
// byte comparison against MarshalCanonical(NewDefaults()) is an AC#1
// determinism check (second run touches nothing, delete-and-rerun reproduces
// the same bytes). As an AC#2 verdict it would be vacuous - implementation
// and expectation would quote the same function (198-a1 §5.2 warns this by
// name), and the reflection-vs-file reading AC#2 needs belongs to 198-r2.

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
)

// run198 drives the user-started entry once over dir and returns the exit code
// plus both streams. notify is a recorder that never posts: nothing here may
// reach a provider, and a real toast would make the case depend on the shell.
func run198(t *testing.T, dir string) (int, string, string) {
	t.Helper()
	out, errb := &bytes.Buffer{}, &bytes.Buffer{}
	code := runTextTask(runSpec{
		argv:    []string{"总结一下 这份笔记"},
		stdout:  out,
		stderr:  errb,
		dataDir: dir,
		notify:  func(string, string) error { return nil },
	})
	return code, out.String(), errb.String()
}

// staticSectionHeads198 reads the `toml:"..."` tags of the exported Config
// struct - the 18 section fields minus the `toml:"-"` one (Plugins, marshalled
// separately by MarshalCanonical's plugins tail) and minus schema_version,
// which is a leaf, not a table. Deriving the names from the struct rather than
// typing them is what keeps this check from becoming a second roster the
// schema is free to outgrow.
func staticSectionHeads198(t *testing.T) []string {
	t.Helper()
	var heads []string
	ct := reflect.TypeOf(config.Config{})
	for i := 0; i < ct.NumField(); i++ {
		f := ct.Field(i)
		tag := f.Tag.Get("toml")
		if tag == "" || tag == "-" {
			continue
		}
		name := strings.Split(tag, ",")[0]
		if f.Type.Kind() == reflect.Struct {
			heads = append(heads, name)
		}
	}
	if len(heads) == 0 {
		t.Fatal("the reader found no section tags, so every head assertion below would pass for nothing")
	}
	return heads
}

// TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone is the three-stage
// positive control from 198-a1 §5.1: create on the first run, touch nothing on
// the second, re-create after a delete - each stage killing one lazy
// implementation (an in-memory flag, a non-atomic rewrite, defaults smuggled
// in from a previous read).
func TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, configFileName)

	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		t.Fatalf("premise broke: a fresh temp data dir already holds %q (err %v)", cfgPath, err)
	}

	// Stage 1: zero config -> the run creates the file AND still fails loudly.
	// The exit code is asserted at 2 on purpose: first-run must not read as
	// "configured now, all clear" - with no model and no key the run still
	// refuses (票面 §3 "不许顺手放宽", 198-a1 R13).
	code, stdoutText, stderrText := run198(t, dir)
	if code != 2 {
		t.Fatalf("stage 1: exit %d, want 2 (created config still lacks model/key; stdout %q stderr %q)",
			code, stdoutText, stderrText)
	}
	if stdoutText != "" {
		t.Errorf("stage 1: stdout = %q, want empty (a refused run streams no task output)", stdoutText)
	}
	first, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("stage 1: after `wisp run` over a fresh machine, config.toml is not there: %v\nstderr:\n%s", err, stderrText)
	}
	if !strings.Contains(stderrText, "新建默认配置") || !strings.Contains(stderrText, cfgPath) {
		t.Errorf("stage 1: the run created the file but never said so with its path (AC: 不许静默建完就当没事发生). stderr:\n%s", stderrText)
	}
	t.Logf("stage 1 receipt (verbatim stderr):\n%s", stderrText)
	t.Logf("stage 1 config.toml: %d bytes, first line %q", len(first), strings.SplitN(string(first), "\n", 2)[0])

	// Stage 2: second run must not rewrite. Bytes AND mtime identical, and no
	// temp file from atomicWrite's rename left behind in the data dir.
	before, err := os.Stat(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	code, _, stderr2 := run198(t, dir)
	if code != 2 {
		t.Fatalf("stage 2: exit %d, want 2, stderr:\n%s", code, stderr2)
	}
	after, err := os.Stat(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if !after.ModTime().Equal(before.ModTime()) {
		t.Errorf("stage 2: mtime moved (%v -> %v): the second run rewrote a config that was already there",
			before.ModTime(), after.ModTime())
	}
	second, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Errorf("stage 2: bytes changed between run one and run two; a creation branch still fired")
	}
	if strings.Contains(stderr2, "新建默认配置") {
		t.Errorf("stage 2: the creation receipt fired again:\n%s", stderr2)
	}
	if strays, _ := filepath.Glob(filepath.Join(dir, ".wisp-config-*.tmp")); len(strays) > 0 {
		t.Errorf("stage 2: atomicWrite temp files left in the data dir: %v", strays)
	}

	// Stage 3: delete and rerun - the branch reads DISK, not a process flag.
	if err := os.Remove(cfgPath); err != nil {
		t.Fatal(err)
	}
	code, _, stderr3 := run198(t, dir)
	if code != 2 {
		t.Fatalf("stage 3: exit %d, want 2, stderr:\n%s", code, stderr3)
	}
	third, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("stage 3: deleting config.toml did not make the next run re-create it: %v", err)
	}
	if !bytes.Equal(first, third) {
		t.Errorf("stage 3: the re-created file differs from the first one; defaults are being sourced from somewhere other than NewDefaults")
	}
}

// TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten is the 种下必响
// control for 198-a1 §4-N5: ticket 101's 权限读不到 stand-in parks a DIRECTORY
// under the config name. "stat failed -> create" would rename straight over it;
// the fs.ErrNotExist test must treat a directory as present and stay out.
func TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, configFileName)
	if err := os.Mkdir(cfgPath, 0o700); err != nil {
		t.Fatal(err)
	}
	code, _, stderrText := run198(t, dir)
	if code != 2 {
		t.Fatalf("exit %d, want 2 (an unreadable config must stay a classified failure, 票 101 AC#3):\n%s", code, stderrText)
	}
	if strings.Contains(stderrText, "新建默认配置") {
		t.Errorf("the first-run branch claimed to create over a directory standing in the file's place:\n%s", stderrText)
	}
	if !strings.Contains(stderrText, "配置未就绪") {
		t.Errorf("the existing fail-closed sentence disappeared; the directory shape no longer reads as the config-class failure it is:\n%s", stderrText)
	}
	info, err := os.Stat(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Fatalf("the directory under the config name is gone - creation overwrote what it was told to leave alone")
	}
	if entries, _ := os.ReadDir(cfgPath); len(entries) != 0 {
		t.Errorf("the untouched directory grew entries: %v", entries)
	}
}

// TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables
// answers AC#1's "真产出一份带全部 section 的默认 config.toml" at the file
// level, and records the J2 reading as a pin: the 17 tagged static section
// tables each appear exactly once, while the dynamic tables (llm.providers,
// models.local_override, plugins entries) appear ZERO times - NewDefaults
// leaves maps nil (defaults.go:77-78), and filling one "so the file looks
// complete" is exactly the invented-default the ticket forbids.
func TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, configFileName)
	if code, _, stderrText := run198(t, dir); code != 2 {
		t.Fatalf("exit %d, want 2, stderr:\n%s", code, stderrText)
	}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	text := "\n" + string(data) // so the first head matches the counted form
	for _, head := range staticSectionHeads198(t) {
		if n := strings.Count(text, "\n["+head+"]\n"); n != 1 {
			t.Errorf("static section [%s] appears %d times, want exactly 1", head, n)
		}
	}
	// The 18th section, Plugins (`toml:"-"`, schema.go:127-131), is emitted by
	// marshalPlugins' own tail with its single static key tier2_enabled - the
	// "18 static sections, dynamic tables empty" reading of J2, pinned here so
	// neither half can silently disappear.
	if n := strings.Count(text, "\n[plugins]\n"); n != 1 {
		t.Errorf("[plugins] head appears %d times, want exactly 1 (marshalPlugins emits the section, parse.go:181-191)", n)
	}
	for _, forbidden := range []string{"[llm.providers", "[models.local_override", "[plugins."} {
		if strings.Contains(text, "\n"+forbidden) {
			t.Errorf("the created file carries the dynamic table %q anyway - that would be an invented entry, not a default", forbidden)
		}
	}
	if n := strings.Count(text, "\nschema_version = "); n != 1 {
		t.Errorf("schema_version lines = %d, want 1 (SaveFile stamps it, loader.go:243)", n)
	}

	// The file the app itself wrote must load (same promise
	// internal/config/unwired_test.go:86-95 pins inside the config package;
	// re-checking it here closes the loop on the CLI-side path).
	if _, _, err := config.LoadFile(cfgPath, nil); err != nil {
		t.Fatalf("a config first-run wrote cannot be loaded: %v", err)
	}
}

// TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences is the 198-a1
// §4-N9 control: 票 223's routing cases assert that each cause sentence
// excludes its siblings, and cause=missing's own text ("本次运行继续用内存里的
// 旧配置") is a lie for a boot that has no previous config. The creation
// receipt therefore shares none of that sentence's markers.
func TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences(t *testing.T) {
	dir := t.TempDir()
	_, _, stderrText := run198(t, dir)
	if !strings.Contains(stderrText, "新建默认配置") {
		t.Fatalf("no creation receipt to test against; stderr:\n%s", stderrText)
	}
	for _, borrowed := range []string{
		"cause=missing",
		"config.toml 读不到：文件不存在",
		"本次运行继续用内存里的旧配置",
	} {
		if strings.Contains(stderrText, borrowed) {
			t.Errorf("the first-run receipt reuses hot-reload wording %q (票 223 keeps the four causes four sentences):\n%s",
				borrowed, stderrText)
		}
	}
}

// TestTicket198FirstRunCallerIsTheRunEntryOnly is J1 as a standing instrument,
// not a comment: the AST call set of ensureFirstRunConfig must be exactly
// {runTextTask}. Adding the call to assembleRuntime - which the resident leg
// shares (resident_task_source_windows.go:265), and whose two pins
// (resident_task_source_246_windows_test.go:389-391,
// logsink_windows_test.go:159-164) this placement keeps green - turns this
// case red before any of them can be "loosened to land".
func TestTicket198FirstRunCallerIsTheRunEntryOnly(t *testing.T) {
	pkg, err := loadMainPackage131(".")
	if err != nil {
		t.Fatalf("loading package main: %v", err)
	}
	// Non-vacuity: the loader must see both the creation function and the two
	// entries the ruling names, or the call-set reading below proves nothing.
	for _, mustSee := range []string{"ensureFirstRunConfig", "runTextTask", "assembleRuntime"} {
		if _, ok := pkg.funcs[mustSee]; !ok {
			t.Fatalf("the call-graph reader does not see %q; this instrument is blind", mustSee)
		}
	}
	var callers []string
	for name, fi := range pkg.funcs {
		if fi.calls["ensureFirstRunConfig"] {
			callers = append(callers, name)
		}
	}
	if len(callers) != 1 || callers[0] != "runTextTask" {
		t.Errorf("AC J1 RED: ensureFirstRunConfig is called by %v, want exactly [runTextTask] - the resident leg reaches the assembly root directly, and first-run creation shared there writes a config.toml its pins forbid", callers)
	}
}
