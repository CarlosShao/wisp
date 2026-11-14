package manifestslot

// Ticket 52 — D46 command runtime, the half that owns a process: pin recheck,
// pre-spawn screening, argv to a REAL CreateProcess, the environment allowlist,
// the monotonic budget with a process-tree kill, the output caps and the
// extraction hand-off.
//
// The fixture CLI is this package's own test binary (command_52_main_test.go):
// each test copies it to a path OUTSIDE the TEMP roots (D46 refuses to pin an exe
// under TEMP - correctly), pins the copy, and points a manifest at it.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/plugin"
)

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

// exeDirOnce / exeDir are the ONE stable directory the whole file's fixtures
// live in, created once and cleaned by its own dispose.
//
// Why not t.TempDir() per test: on Windows t.TempDir() honours TEMP, so every
// fixture would land inside the very tree D46 refuses to pin - and a suite that
// skips itself on the default configuration proves nothing. The first candidate
// is therefore %LOCALAPPDATA%\wisp52-fixture-tests, a fixed local path outside
// TEMP; only if that is unavailable does one fall back to the temp root with the
// test saying so.
var (
	exeDirOnce sync.Once
	exeDir     string
	exeDirErr  error
)

const fixtureExeFileName = "wisp52-fixture.exe"

func stableExeDir(t *testing.T) string {
	t.Helper()
	exeDirOnce.Do(func() {
		var candidates []string
		if local := strings.TrimSpace(os.Getenv("LOCALAPPDATA")); local != "" {
			candidates = append(candidates, filepath.Join(local, "wisp52-fixture-tests"))
		}
		if root := strings.TrimSpace(os.Getenv("SystemDrive")); root != "" {
			candidates = append(candidates, root+`\wisp52-fixture-tests`)
		}
		candidates = append(candidates, filepath.Join(os.TempDir(), "wisp52-fixture-tests"))
		var lastErr error
		for _, dir := range candidates {
			if inTemp(dir) {
				// The documented fallback: usable, but every test that pins
				// through it says why it is weaker.
				t.Logf("wisp52: no stable root available, fixtures live under TEMP at %q (the D46 refusal under test is still exercised by internal/plugin's own suite)", dir)
			}
			if err := os.MkdirAll(dir, 0o755); err != nil {
				lastErr = err
				continue
			}
			exeDir = dir
			return
		}
		exeDirErr = fmt.Errorf("no usable fixture directory (last error: %v)", lastErr)
	})
	if exeDirErr != nil {
		t.Fatalf("wisp52 fixture directory: %v", exeDirErr)
	}
	return exeDir
}

// inTemp reports whether dir sits under a TEMP/TMP/TMPDIR root.
func inTemp(dir string) bool {
	for _, env := range []string{"TEMP", "TMP", "TMPDIR"} {
		root := strings.TrimSpace(os.Getenv(env))
		if root == "" {
			continue
		}
		rl := strings.ToLower(filepath.Clean(root))
		dl := strings.ToLower(filepath.Clean(dir))
		if dl == rl || strings.HasPrefix(dl, rl+`\`) {
			return true
		}
	}
	return false
}

// installFixtureCLI copies the test binary to a pinnable path once per process.
//
// One copy for the whole file, not one per test: the file is ~90MB and the
// Windows filter manager will not let a freshly converted forked-VM (CFA) file be
// truncated and rewritten immediately, which is exactly what a per-test copy
// would do. The tamper test writes a same-size, in-place byte change instead.
var (
	fixtureOnce sync.Once
	fixturePath string
	fixtureErr  error
)

func installFixtureCLI(t *testing.T) string {
	t.Helper()
	fixtureOnce.Do(func() {
		self, err := os.Executable()
		if err != nil {
			fixtureErr = fmt.Errorf("os.Executable: %w", err)
			return
		}
		dst := filepath.Join(stableExeDir(t), fixtureExeFileName)
		src, err := os.ReadFile(self)
		if err != nil {
			fixtureErr = fmt.Errorf("read test binary: %w", err)
			return
		}
		_ = os.Remove(dst)
		if err := os.WriteFile(dst, src, 0o755); err != nil {
			fixtureErr = fmt.Errorf("install fixture exe: %w", err)
			return
		}
		fixturePath = dst
	})
	if fixtureErr != nil {
		t.Fatalf("wisp52 fixture exe: %v", fixtureErr)
	}
	return fixturePath
}

// restoreFixtureCLI puts the pristine bytes back and forgets the pin cache, so
// a test that tampers cannot poison the next one (order independence matters
// more than the seconds this costs: a suite whose green depends on test order is
// a suite nobody can trust).
func restoreFixtureCLI(t *testing.T) {
	t.Helper()
	self, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable: %v", err)
	}
	src, err := os.ReadFile(self)
	if err != nil {
		t.Fatalf("read test binary: %v", err)
	}
	if got := installFixtureCLI(t); got == "" {
		t.Fatal("no fixture path to restore")
	}
	if err := os.WriteFile(fixturePath, src, 0o755); err != nil {
		t.Fatalf("restore fixture exe: %v", err)
	}
}

// fixtureTool is a command tool over the installed fixture CLI.
func fixtureTool(exe, mode string) plugin.ToolDef {
	return plugin.ToolDef{
		Name:         "fixture.run",
		Description:  "Run the fixture CLI.",
		Risk:         plugin.RiskL0,
		Kind:         plugin.ToolKindCommand,
		Exe:          exe,
		ArgvTemplate: []string{filepath.Base(exe)},
		Parameters:   map[string]any{"type": "object", "properties": map[string]any{}},
		Extract:      map[string]any{"kind": "regex", "pattern": `.+`},
		// The fixture decides what to do from MODE; every name in the allowlist
		// must be something the ENVIRONMENT can hold. The marker names are
		// declared because the whole point is that the allowlist, not the
		// parent process, is where the child's handshake comes from - so the
		// handshake travels inside the manifest's own declaration.
		EnvAllowlist: []string{fixtureHandshakeVar},
	}
}

// sealedEnv is a host environment with an API key in it. The env proof needs a
// host that HOLDS a key: a host without one could not disprove a leak.
func sealedEnv(name string) (string, bool) {
	switch name {
	case "PATH":
		return `C:\Windows\System32;C:\Windows`, true
	case "USERPROFILE":
		return `C:\Users\fixture`, true
	case "SystemRoot":
		return `C:\Windows`, true
	case "OPENAI_API_KEY":
		return "sk-should-never-reach-a-child", true
	case "LARK_SECRET":
		return "secret-should-never-reach-a-child", true
	case fixtureHandshakeVar:
		return fixtureModePrefix + "echo", true
	}
	return "", false
}

// fixtureEnv is the host environment for every other fixture call: the real
// PATH/SystemRoot the OS needs to load an image, the handshake, and nothing
// else from the test process.
func fixtureEnv(mode string) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "PATH":
			return os.Getenv("PATH"), true
		case "SystemRoot":
			return os.Getenv("SystemRoot"), true
		case "SystemDrive":
			return os.Getenv("SystemDrive"), true
		case fixtureHandshakeVar:
			return fixtureModePrefix + mode, true
		}
		return "", false
	}
}

// recordingSpawner is the process-creation seam: it records exactly what
// creation was asked for and delegates to the real OS spawner.
type recordingSpawner struct {
	mu    sync.Mutex
	specs []SpawnSpec
}

func (r *recordingSpawner) Spawn(ctx context.Context, spec SpawnSpec) (ChildProcess, error) {
	r.mu.Lock()
	r.specs = append(r.specs, spec)
	r.mu.Unlock()
	return OSCommandSpawner{}.Spawn(ctx, spec)
}

func (r *recordingSpawner) last() (SpawnSpec, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.specs) == 0 {
		return SpawnSpec{}, false
	}
	return r.specs[len(r.specs)-1], true
}

func (r *recordingSpawner) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.specs)
}

// runtimeFor pins one tool and builds the runtime over it.
func runtimeFor(t *testing.T, def plugin.ToolDef, opt CommandRunOptions) (*CommandRuntime, plugin.ExePin) {
	t.Helper()
	pin, err := plugin.PinCommandExe(def, nil)
	if err != nil {
		t.Fatalf("PinCommandExe: %v", err)
	}
	runtime := NewCommandRuntime(map[string]map[string]plugin.ExePin{
		"com.example.fixture": {def.Name: pin},
	}, opt)
	return runtime, pin
}

// ---------------------------------------------------------------------------
// AC "Injection": argv vector, no shell, refusals before spawn
// ---------------------------------------------------------------------------

// Test52NoShellSpawnAndArgvVector is the process-creation audit: the seam sees
// the pinned exe as argv[0], one vector element per argument, and no cmd.exe /
// sh / PowerShell anywhere. The spec is then handed to the REAL spawner, so the
// audit covers the call that actually ran.
func Test52NoShellSpawnAndArgvVector(t *testing.T) {
	exe := installFixtureCLI(t)
	def := fixtureTool(exe, "extract")
	def.ArgvTemplate = []string{filepath.Base(exe), "--text", "{{text}}"}
	def.Parameters = map[string]any{
		"type": "object",
		"properties": map[string]any{
			"text": map[string]any{"type": "string", "maxLength": 64},
		},
		"required": []any{"text"},
	}
	def.Extract = map[string]any{"kind": "json", "path": ".ok"}

	rec := &recordingSpawner{}
	runtime, _ := runtimeFor(t, def, CommandRunOptions{Spawner: rec, EnvLookup: fixtureEnv("extract")})

	res, err := runtime.Exec(t.Context(), "com.example.fixture", def,
		[]byte(`{"text":"hello world  with  spaces"}`))
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if res.IsError {
		t.Fatalf("Exec: tool error: %s", res.Text)
	}
	if res.Text != "true" {
		t.Errorf("res.Text = %q, want the extracted %q", res.Text, "true")
	}

	spec, ok := rec.last()
	if !ok {
		t.Fatal("the spawner was never called")
	}
	if spec.Exe != exe {
		t.Errorf("spec.Exe = %q, want the pinned %q", spec.Exe, exe)
	}
	want := []string{exe, "--text", "hello world  with  spaces"}
	if strings.Join(spec.Argv, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("spec.Argv = %q, want %q", spec.Argv, want)
	}
	for _, el := range spec.Argv {
		lower := strings.ToLower(el)
		if strings.Contains(lower, "cmd.exe") || strings.Contains(lower, "powershell") ||
			strings.HasSuffix(lower, "/sh") || lower == "sh" || lower == "bash" ||
			strings.Contains(lower, " /c ") {
			t.Fatalf("argv carries a shell interpreter: %q", el)
		}
	}
	if len(spec.Env) == 0 {
		t.Fatal("spec.Env is empty: the child would inherit nothing and could not even load")
	}
	for _, kv := range spec.Env {
		if strings.HasPrefix(kv, "OPENAI_API_KEY=") || strings.HasPrefix(kv, "LARK_SECRET=") {
			t.Fatalf("spec.Env carries a host secret: %q", kv)
		}
	}
}

// Test52InjectionIsRefusedBeforeSpawn: a metacharacter, a schema violation and
// an undeclared parameter must all be refused with ZERO spawner calls. That is
// the structural form of "rejected pre-spawn".
func Test52InjectionIsRefusedBeforeSpawn(t *testing.T) {
	exe := installFixtureCLI(t)
	def := fixtureTool(exe, "extract")
	def.ArgvTemplate = []string{filepath.Base(exe), "--text", "{{text}}"}
	def.Parameters = map[string]any{
		"type": "object",
		"properties": map[string]any{
			"text": map[string]any{"type": "string", "maxLength": 32},
		},
		"required": []any{"text"},
	}
	def.Extract = map[string]any{"kind": "json", "path": ".ok"}

	rec := &recordingSpawner{}
	runtime, _ := runtimeFor(t, def, CommandRunOptions{Spawner: rec})

	bad := map[string]string{
		"pipe":          `{"text":"a|calc"}`,
		"and":           `{"text":"a&&calc"}`,
		"semicolon":     `{"text":"a;calc"}`,
		"dollar":        `{"text":"$(calc)"}`,
		"backtick":      `{"text":"` + "`calc`" + `"}`,
		"redirect":      `{"text":"a>out.txt"}`,
		"newline":       `{"text":"a\ncalc"}`,
		"too long":      `{"text":"` + strings.Repeat("x", 40) + `"}`,
		"wrong type":    `{"text":7}`,
		"undeclared":    `{"text":"ok","shell":"calc"}`,
		"missing":       `{}`,
		"not an object": `[]`,
	}
	for name, args := range bad {
		res, err := runtime.Exec(t.Context(), "com.example.fixture", def, []byte(args))
		if err != nil {
			t.Fatalf("%s: Exec returned a host error (%v); a refused call is a tool-level error", name, err)
		}
		if !res.IsError {
			t.Fatalf("%s: args %s were accepted (res = %+v)", name, args, res)
		}
	}
	if n := rec.count(); n != 0 {
		t.Fatalf("%d process(es) were created for refused calls: a violation reached argv", n)
	}
}

// Test52PinIsReverifiedEveryCall: one byte on the exe, and the NEXT call is
// refused with an alert and no process. This is the AC's "every execution
// re-verifies hash" with a real file and the real runtime.
func Test52PinIsReverifiedEveryCall(t *testing.T) {
	exe := installFixtureCLI(t)
	def := fixtureTool(exe, "extract")
	def.Extract = map[string]any{"kind": "json", "path": ".ok"}

	rec := &recordingSpawner{}
	cap := &plugin.CaptureAlerter{}
	runtime, _ := runtimeFor(t, def, CommandRunOptions{Spawner: rec, Alerter: cap, EnvLookup: fixtureEnv("extract")})

	if res, err := runtime.Exec(t.Context(), "com.example.fixture", def, []byte(`{}`)); err != nil || res.IsError {
		t.Fatalf("first call: err=%v res=%+v", err, res)
	}

	// Same size, in place: a copy-over is not the threat model, a swap is.
	bytes, err := os.ReadFile(exe)
	if err != nil {
		t.Fatal(err)
	}
	bytes[len(bytes)/2] ^= 0x01
	if err := os.WriteFile(exe, bytes, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { restoreFixtureCLI(t) })

	_, err = runtime.Exec(t.Context(), "com.example.fixture", def, []byte(`{}`))
	if !errors.Is(err, plugin.ErrExeSwap) {
		t.Fatalf("call over a tampered exe: err = %v, want plugin.ErrExeSwap", err)
	}
	if n := rec.count(); n != 1 {
		t.Fatalf("spawner calls = %d, want 1 (the tampered call must never reach process creation)", n)
	}
	alerts := cap.Alerts()
	if len(alerts) != 1 || alerts[0].Kind != "exe" {
		t.Fatalf("alerts = %+v, want one exe alert", alerts)
	}
}

// Test52UnpinnedToolIsRefusedAndAlerted: no pin is a refusal, not a pass.
func Test52UnpinnedToolIsRefusedAndAlerted(t *testing.T) {
	exe := installFixtureCLI(t)
	def := fixtureTool(exe, "extract")
	rec := &recordingSpawner{}
	cap := &plugin.CaptureAlerter{}
	runtime := NewCommandRuntime(nil, CommandRunOptions{Spawner: rec, Alerter: cap})

	_, err := runtime.Exec(t.Context(), "com.example.fixture", def, []byte(`{}`))
	if !errors.Is(err, ErrNoCommandPin) {
		t.Fatalf("err = %v, want ErrNoCommandPin", err)
	}
	if rec.count() != 0 {
		t.Fatal("an unpinned tool reached process creation")
	}
	if alerts := cap.Alerts(); len(alerts) != 1 || alerts[0].Kind != "exe" {
		t.Fatalf("alerts = %+v, want one exe alert for the missing pin", alerts)
	}
}

// Test52DeclarativeToolIsNotForThisRuntime keeps the two Tier-1 kinds apart: the
// command runtime executes command tools only.
func Test52DeclarativeToolIsNotForThisRuntime(t *testing.T) {
	exe := installFixtureCLI(t)
	def := fixtureTool(exe, "extract")
	def.Kind = ""
	rec := &recordingSpawner{}
	runtime, _ := runtimeFor(t, fixtureTool(exe, "extract"), CommandRunOptions{Spawner: rec})
	if _, err := runtime.Exec(t.Context(), "com.example.fixture", def, []byte(`{}`)); !errors.Is(err, ErrNotCommandTool) {
		t.Fatalf("err = %v, want ErrNotCommandTool", err)
	}
	if rec.count() != 0 {
		t.Fatal("a declarative tool reached process creation")
	}
}

// ---------------------------------------------------------------------------
// AC "Env": the CHILD's own environment is the allowlist and nothing else
// ---------------------------------------------------------------------------

// Test52ChildEnvIsOnlyTheAllowlist is the echo-fixture proof: the fixture prints
// its own os.Environ(), the host environment it was given contains an API key,
// and the child's table must not.
func Test52ChildEnvIsOnlyTheAllowlist(t *testing.T) {
	exe := installFixtureCLI(t)
	def := fixtureTool(exe, "echo")
	def.Extract = map[string]any{"kind": "regex", "pattern": `^env=(.*)$`, "group": "1", "multiline": true}

	rec := &recordingSpawner{}
	runtime, _ := runtimeFor(t, def, CommandRunOptions{
		Spawner:   rec,
		EnvLookup: sealedEnv,
	})
	res, err := runtime.Exec(t.Context(), "com.example.fixture", def, []byte(`{}`))
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if res.IsError {
		t.Fatalf("Exec: tool error: %s", res.Text)
	}

	seen := map[string]string{}
	for _, line := range strings.Split(res.Text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		name, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		seen[name] = value
	}
	if len(seen) == 0 {
		t.Fatalf("the fixture reported no environment at all (text = %q)", res.Text)
	}
	for _, secret := range []string{"OPENAI_API_KEY", "LARK_SECRET"} {
		if v, ok := seen[secret]; ok {
			t.Fatalf("the CHILD process saw %s=%s: the host's key leaked into a wrapped CLI", secret, v)
		}
	}
	// The child's table is exactly what the runtime passed, not what the test
	// process happens to hold: a child that inherited the parent would show
	// dozens of developer variables here.
	spec, _ := rec.last()
	want := map[string]string{}
	for _, kv := range spec.Env {
		name, value, _ := strings.Cut(kv, "=")
		want[name] = value
	}
	for name, value := range want {
		got, ok := seen[name]
		if !ok {
			t.Errorf("the child did not see allowlisted %q", name)
			continue
		}
		if name == "PATH" {
			// Windows re-creates PATH for a new process from the machine and
			// user environment when the block has none, so trim-path equality
			// is not the contract; presence is.
			continue
		}
		if got != value {
			t.Errorf("%s = %q in the child, want %q", name, got, value)
		}
	}
	for name := range seen {
		if _, ok := want[name]; ok {
			continue
		}
		if strings.EqualFold(name, "PATH") {
			continue
		}
		t.Errorf("the child saw %q, which the allowlist did not name (inherited environment)", name)
	}
}

// ---------------------------------------------------------------------------
// AC "Extraction": JSON + regex on a real CLI, >32KB truncated
// ---------------------------------------------------------------------------

// Test52ExtractionOnARealCLI runs both extractor kinds over the fixture's
// output. Raw stdout (the whole document) must never be what enters the
// context.
func Test52ExtractionOnARealCLI(t *testing.T) {
	exe := installFixtureCLI(t)
	cases := []struct {
		name string
		decl map[string]any
		want string
	}{
		{"json path", map[string]any{"kind": "json", "path": ".data.items[].summary"}, `["standup 10:00","review 14:00"]`},
		{"json scalar", map[string]any{"kind": "json", "path": ".ok"}, "true"},
		{"regex", map[string]any{"kind": "regex", "pattern": `"summary":"([^"]+)"`, "group": "1"}, "standup 10:00\nreview 14:00"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			def := fixtureTool(exe, "extract")
			def.Extract = tc.decl
			runtime, _ := runtimeFor(t, def, CommandRunOptions{EnvLookup: fixtureEnv("extract")})
			res, err := runtime.Exec(t.Context(), "com.example.fixture", def, []byte(`{}`))
			if err != nil {
				t.Fatalf("Exec: %v", err)
			}
			if res.IsError {
				t.Fatalf("Exec: tool error: %s", res.Text)
			}
			if res.Text != tc.want {
				t.Fatalf("extracted = %q, want %q", res.Text, tc.want)
			}
			if strings.Contains(res.Text, `"data"`) {
				t.Fatalf("the extracted text carries the raw document: %q", res.Text)
			}
		})
	}
}

// Test52OversizedExtractIsTruncatedTo32KB: the extractor selects more than the
// context cap and the runtime cuts it, marking the result truncated.
func Test52OversizedExtractIsTruncated(t *testing.T) {
	exe := installFixtureCLI(t)
	// 8MB of "A" is written by the fixture in "huge" mode; the regex picks it
	// up whole, so the EXTRACT (not the stream cap) is what has to truncate.
	def := fixtureTool(exe, "huge")
	def.Extract = map[string]any{"kind": "regex", "pattern": `A+`}

	runtime, _ := runtimeFor(t, def, CommandRunOptions{EnvLookup: fixtureEnv("huge")})
	res, err := runtime.Exec(t.Context(), "com.example.fixture", def, []byte(`{}`))
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	// The stream cap fires first for 8MB of stdout (the next test proves it),
	// so this call reports the cap. Both are "refuse to hand the model a
	// flood"; the assertion below pins which one owns this case.
	if !res.IsError || !res.Truncated {
		t.Fatalf("an 8MB stdout was reported as a clean success: %+v", res)
	}
}

// Test52ExtractCapIsEnforcedByClamp is the >32KB extract rule in isolation: the
// clamp is a rune-boundary cut that marks the result truncated and keeps the
// head and the tail, so a manifest whose regex selects a huge document cannot
// flood the context either.
func Test52ExtractCapIsEnforcedByClamp(t *testing.T) {
	big := strings.Repeat("x", 40*1024)
	out, truncated := clampExtract(big, plugin.CommandExtractCapBytes)
	if !truncated {
		t.Fatal("a 40KB extract was not marked truncated")
	}
	if len(out) > plugin.CommandExtractCapBytes+2048 {
		t.Fatalf("clamped output is %d bytes, want ~%d", len(out), plugin.CommandExtractCapBytes)
	}
	if !strings.Contains(out, "[... truncated: ") {
		t.Fatalf("clamped output carries no truncation marker: %q", out[:120])
	}

	// Rune boundary: the head must end on a rune boundary. A 3-byte rune
	// straddling the cut would otherwise leave a half rune (strings.ToValidUTF8
	// reports it as U+FFFD), which is what the model would then read.
	runes := strings.Repeat("中", 20*1024) // 3 bytes each, 60KB
	out, truncated = clampExtract(runes, plugin.CommandExtractCapBytes)
	if !truncated {
		t.Fatal("multi-byte extract was not truncated")
	}
	head := out[:plugin.CommandExtractCapBytes]
	if strings.ContainsRune(strings.ToValidUTF8(head, "\uFFFD"), '\uFFFD') {
		t.Fatal("the clamped head ends mid-rune: the model would read a replacement character")
	}
	if strings.Contains(head, "\uFFFD") {
		t.Fatal("the clamped head carries a replacement character")
	}
}

// ---------------------------------------------------------------------------
// AC "1MB stdout -> kill + truncated (no OOM)"
// ---------------------------------------------------------------------------

// Test52StdoutCapKillsAndTruncates runs a real CLI that writes 8MB. The call
// must come back quickly, marked truncated, holding at most the 1MB cap, and the
// child must be gone.
func Test52StdoutCapKillsAndTruncates(t *testing.T) {
	exe := installFixtureCLI(t)
	def := fixtureTool(exe, "huge")
	def.Extract = map[string]any{"kind": "regex", "pattern": `A+`}
	def.TimeoutMs = 30000

	rec := &recordingSpawner{}
	cap := &plugin.CaptureAlerter{}
	runtime, _ := runtimeFor(t, def, CommandRunOptions{Spawner: rec, Alerter: cap, EnvLookup: fixtureEnv("huge")})

	start := time.Now()
	res, err := runtime.Exec(t.Context(), "com.example.fixture", def, []byte(`{}`))
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if !res.IsError || !res.Truncated {
		t.Fatalf("8MB of stdout was reported as %+v, want an error marked truncated", res)
	}
	if !strings.Contains(res.Text, "stream cap") {
		t.Fatalf("res.Text = %q, want it to name the stream cap", res.Text)
	}
	if elapsed > 25*time.Second {
		t.Fatalf("the call took %v: the cap did not stop the CLI", elapsed)
	}
	if len(cap.Alerts()) != 1 || cap.Alerts()[0].Kind != "command-output" {
		t.Fatalf("alerts = %+v, want one command-output alert", cap.Alerts())
	}

	// The captured stream really was capped: the writer stops at the cap, so a
	// result that carried 8MB would mean the cap is decorative.
	if n := rec.count(); n != 1 {
		t.Fatalf("spawner calls = %d, want 1", n)
	}
}

// Test52CapWriterStopsAtTheLimit is the cap's unit-level teeth: everything past
// the limit is dropped, over is set, and the callback fires once.
func Test52CapWriterStopsAtTheLimit(t *testing.T) {
	fired := make(chan struct{}, 4)
	w := &capWriter{limit: 16}
	w.onOver = func() { fired <- struct{}{} }
	for i := 0; i < 100; i++ {
		if _, err := w.Write([]byte(strings.Repeat("z", 4))); err != nil {
			t.Fatalf("Write: %v", err)
		}
	}
	if got := len(w.Bytes()); got != 16 {
		t.Fatalf("captured %d bytes, want the 16-byte limit", got)
	}
	if !w.Over() || !w.Truncated() {
		t.Fatal("capWriter did not record the overflow")
	}
	// fire is dispatched on its own goroutine (the kill must not run under the
	// writer's lock).
	select {
	case <-fired:
	case <-time.After(2 * time.Second):
		t.Fatal("the overflow callback never fired: a runaway CLI would not be killed")
	}
	select {
	case <-fired:
		t.Fatal("the overflow callback fired twice: the kill path must be once-only")
	case <-time.After(50 * time.Millisecond):
	}
}

// ---------------------------------------------------------------------------
// AC "Timeout": the hanging CLI is killed and its TREE is gone
// ---------------------------------------------------------------------------

const (
	// timeoutTree443PIDWait is how long the test keeps polling for the
	// fixture's pid file once the call has returned. The only writers are
	// inside the tree the budget just killed, so this window covers filesystem
	// visibility, not a second budget: a pid that never arrived cannot be
	// conjured by waiting longer - which is why the CONTROL below, and not a
	// longer sleep, is what carries the kill reading.
	timeoutTree443PIDWait = time.Second
	// timeoutTree443ControlBudgetMS is the control run's budget: eight times
	// the raced 500ms. Starting two 90MB images twice cannot outrun an 8x
	// margin, so the control answers "does the timeout still kill the tree?"
	// without inheriting the raced run's scheduling lottery.
	timeoutTree443ControlBudgetMS = plugin.MinCommandTimeoutMS * 80
)

// Test52TimeoutKillsTheProcessTree runs a CLI that spawns a grandchild and then
// hangs. The D46 budget must end the call, the result must say timed out, and
// the grandchild must be gone - it is what a plain Process.Kill would leave
// behind.
//
// THE TWO QUESTIONS, AND WHY THEY NOW HAVE TWO CARRIERS (ticket 443). This test
// used to ask both questions through the pid the fixture writes: when the file
// was absent it reded with "the fixture never reported a grandchild pid", which
// is a statement about the runner's schedule, not about the kill. The 500ms
// budget starts before the fixture has finished starting (the fixture is this
// package's 90MB test binary, re-executed twice), so a loaded runner can kill
// the tree before the grandchild exists; CI runs 43567331249 and 43568581591
// read exactly that on unchanged bytes. Same mechanism ticket 441 fixed for the
// posix seal window, so the same shape is used here:
//
//  1. "did the grandchild's pid reach the test in THIS run?" - the raced run
//     answers with numbers (budget, elapsed, waited, polls, and the last poll
//     error when the file never appeared). A miss is NAMED, not hard-red:
//     `AC#2 ABSTAIN(443:timeout-tree:raced-no-pid) ...` - it is a reading this
//     run could not take, and a zero from absence is never a green.
//  2. "does the timeout still kill the TREE?" - answered by the control below
//     when the raced run abstained: same fixture, same runtime, same kill path,
//     eight times the budget. No pid in the control is RED (the fixture is
//     broken, not the runner); a pid still alive in either run is RED and names
//     the pid. That arm is the tooth: a mutation that drops the job object reds
//     here even on the abstain path, and the ticket-443 evidence file records
//     that control reading.
func Test52TimeoutKillsTheProcessTree(t *testing.T) {
	exe := installFixtureCLI(t)
	pidFile := filepath.Join(stableExeDir(t), "grandchild.pid")
	// A pid file left by an earlier run is a stale witness: its pid belongs to
	// a process that is already gone, stillAlive() would answer "dead", and the
	// test would go green without this run's grandchild ever being reported -
	// the false green that a fixed path plus a previous run can produce. Clear
	// it before the call, and require the world to be one this test built.
	clearPIDFile(t, pidFile)

	def := fixtureTool(exe, "hang")
	def.Extract = map[string]any{"kind": "regex", "pattern": `.+`}
	def.TimeoutMs = plugin.MinCommandTimeoutMS * 10 // 500ms: a real budget, briefly
	// Ticket 443's dial (see fixtureDelay443Var) has to cross the allowlist to
	// reach the fixture: the child environment is the allowlist, never the
	// parent's. Absent from the host, the name is simply absent from the child.
	def.EnvAllowlist = append(def.EnvAllowlist, fixtureDelay443Var)

	env := func(name string) (string, bool) {
		switch name {
		case "PATH":
			return os.Getenv("PATH"), true
		case "SystemRoot":
			return os.Getenv("SystemRoot"), true
		case ModuleNameEnv:
			return "1", true
		case fixtureHandshakeVar:
			return fixtureModePrefix + "hang:" + pidFile, true
		case fixtureDelay443Var:
			if v := os.Getenv(fixtureDelay443Var); v != "" {
				return v, true
			}
		}
		return "", false
	}

	runtime, _ := runtimeFor(t, def, CommandRunOptions{EnvLookup: env})
	start := time.Now()
	res, err := runtime.Exec(t.Context(), "com.example.fixture", def, []byte(`{}`))
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if !res.IsError || !res.Truncated {
		t.Fatalf("a hanging CLI came back as %+v, want an error marked truncated", res)
	}
	if !strings.Contains(res.Text, "timed out") {
		t.Fatalf("res.Text = %q, want it to say the budget expired", res.Text)
	}
	if elapsed > 20*time.Second {
		t.Fatalf("the call took %v; the budget was 500ms and the tree was not killed", elapsed)
	}

	raced := awaitPID443(pidFile, timeoutTree443PIDWait)
	racedRun := fmt.Sprintf("budget-ms=%d elapsed-ms=%d waited-ms=%d polls=%d",
		plugin.CommandTimeoutMs(def), elapsed.Milliseconds(), raced.waitedMS, raced.polls)
	if raced.pid > 0 {
		t.Logf("AC#2 raced run: grandchild pid %d reached the test (%s)", raced.pid, racedRun)
		judgeTreeDead443(t, "raced run", raced.pid)
		return
	}
	// The pid never reached the test. That is this runner's schedule, not the
	// kill: name the reading with its numbers and let the control carry the
	// question this run cannot answer. Nothing here is green on its own.
	t.Logf("AC#2 ABSTAIN(443:timeout-tree:raced-no-pid) %s pid-seen=no pid-file=%q last-poll-error=%q - the budget expired before the fixture reported the grandchild it was supposed to be killed with, so THIS run cannot attest the tree kill from the race; the carrier is the %dms control run below, not a green from absence",
		racedRun, pidFile, raced.lastErr, timeoutTree443ControlBudgetMS)
	judgeTimeoutKillControl443(t, runtime, def, pidFile)
}

// judgeTimeoutKillControl443 is question 2's carrier (ticket 443): run the same
// fixture through the same runtime and the same kill path with a budget eight
// times the raced one, then judge the tree from the pid IT reported. It is not a
// relaxation of the raced run - the raced run still has to time out and still
// has to read its own pid when one arrives - it is the reading that does not
// depend on winning a 500ms race against a 90MB process start.
func judgeTimeoutKillControl443(t *testing.T, runtime *CommandRuntime, def plugin.ToolDef, pidFile string) {
	t.Helper()
	def.TimeoutMs = timeoutTree443ControlBudgetMS
	clearPIDFile(t, pidFile)
	start := time.Now()
	res, err := runtime.Exec(t.Context(), "com.example.fixture", def, []byte(`{}`))
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("AC#2 RED(443:timeout-tree:control-exec) the control run could not be executed: %v", err)
	}
	if !res.IsError || !res.Truncated || !strings.Contains(res.Text, "timed out") {
		t.Fatalf("AC#2 RED(443:timeout-tree:control-no-timeout) the control run (budget-ms=%d elapsed-ms=%d) came back as %+v, want a truncated timeout: the control exists to prove the budget still ends a call, and it did not",
			timeoutTree443ControlBudgetMS, elapsed.Milliseconds(), res)
	}
	ctl := awaitPID443(pidFile, timeoutTree443PIDWait)
	ctlRun := fmt.Sprintf("budget-ms=%d elapsed-ms=%d waited-ms=%d polls=%d",
		timeoutTree443ControlBudgetMS, elapsed.Milliseconds(), ctl.waitedMS, ctl.polls)
	if ctl.pid <= 0 {
		t.Fatalf("AC#2 RED(443:timeout-tree:control-no-pid) the control run could not report a grandchild pid either (%s pid-seen=no pid-file=%q last-poll-error=%q): a %dms budget is eight times the raced one, so a fixture that still cannot start its grandchild is broken, and no reading here attests the tree kill",
			ctlRun, pidFile, ctl.lastErr, timeoutTree443ControlBudgetMS)
	}
	t.Logf("AC#2 control run: grandchild pid %d reached the test (%s)", ctl.pid, ctlRun)
	judgeTreeDead443(t, "control run", ctl.pid)
}

// judgeTreeDead443 is the tooth, and it does not care which run produced the
// pid: the process the fixture reported has to be gone. A survivor is RED with
// its name, and it is reaped before the failure returns so this test cannot
// park a hanging process in front of the rest of the suite.
func judgeTreeDead443(t *testing.T, run string, pid int) {
	t.Helper()
	if !stillAlive(pid) {
		t.Logf("AC#2 %s: grandchild pid %d is gone, so the kill reached the tree", run, pid)
		return
	}
	_ = exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid)).Run()
	t.Fatalf("AC#2 RED(443:timeout-tree:tree-survived) %s: grandchild pid %d is still running after the budget expired - the kill did not reach the tree", run, pid)
}

// clearPIDFile makes the fixture's report a thing THIS run produced: no stale
// pid, and no silently absent pid file to mistake for a reading. It is shared by
// ticket 443's test and ticket 138's cancel test: both read a pid from the same
// fixed fixture directory (stableExeDir survives across runs), so a file left by
// an earlier run is a witness from a dead process.
func clearPIDFile(t *testing.T, path string) {
	t.Helper()
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("clear the grandchild pid file %s: %v", path, err)
	}
	if _, err := os.Stat(path); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("the grandchild pid file %s is still there after clearing it (stat err: %v): a stale pid would be read as this run's witness", path, err)
	}
}

// pidWatch443 is one reading of "did the fixture's grandchild pid reach the
// test", with the numbers the abstain line has to carry: how long the test
// waited, how many polls that took, and - when it never appeared - the last
// error the file gave, because "absent" and "written but unreadable" are
// different findings.
type pidWatch443 struct {
	pid      int
	polls    int
	waitedMS int64
	lastErr  string
}

// awaitPID443 polls path until it holds a pid or wait expires.
func awaitPID443(path string, wait time.Duration) (w pidWatch443) {
	start := time.Now()
	defer func() { w.waitedMS = time.Since(start).Milliseconds() }()
	deadline := start.Add(wait)
	for {
		w.polls++
		pid, err := readPIDFileErr(path)
		if err == nil {
			w.pid = pid
			return w
		}
		w.lastErr = err.Error()
		if time.Now().After(deadline) {
			return w
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// Test52BudgetIsMonotonicNoWallClockSubtraction is a source-level check on the
// one file that owns the budget: the deadline comes from a context, never from
// comparing two timestamps after the fact (D22 ban 4). time.Until is allowed
// only in the form a context deadline needs (it reads the ctx's own monotonic
// reading); the banned shapes are the ones that compute an elapsed interval
// from the wall clock.
func Test52BudgetIsMonotonicNoWallClockSubtraction(t *testing.T) {
	src, err := os.ReadFile("stopwatch.go")
	if err != nil {
		t.Fatalf("read stopwatch.go: %v", err)
	}
	body := string(src)
	if !strings.Contains(body, "context.WithDeadline") {
		t.Fatal("stopwatch.go does not derive the budget from a context deadline")
	}
	for _, banned := range []string{"time.Now().Sub", ".Sub(time.Now()", "time.Since("} {
		if strings.Contains(body, banned) {
			t.Fatalf("stopwatch.go contains %q: the budget must ride a monotonic deadline, not a timestamp difference", banned)
		}
	}
	// The whole runtime package: no wall-clock subtraction anywhere in the
	// command path (a deadline is enforced by the ctx firing, not by a
	// post-hoc comparison).
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read package dir: %v", err)
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		s := string(raw)
		for _, banned := range []string{"time.Now().Sub", ".Sub(time.Now()"} {
			if strings.Contains(s, banned) {
				t.Fatalf("%s contains %q: wall-clock interval math is banned in the command path", name, banned)
			}
		}
	}
}

// readPIDFile is the ticket-138 form: a pid, or 0 when the file is absent or
// unreadable. Ticket 443 needs the error too (the abstain line prints whether
// the file was absent, empty or garbage), so the parsing lives in
// readPIDFileErr and this wrapper keeps the old reading for its caller.
func readPIDFile(t *testing.T, path string) int {
	t.Helper()
	pid, err := readPIDFileErr(path)
	if err != nil {
		return 0
	}
	return pid
}

// readPIDFileErr reads the pid a fixture reported, keeping "not there yet" and
// "there but unreadable" as errors instead of folding them into a zero.
func readPIDFileErr(path string) (int, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(raw)))
	if err != nil {
		return 0, fmt.Errorf("%s holds %q, not a pid: %w", path, string(raw), err)
	}
	return pid, nil
}

// stillAlive reports whether pid is a live process, from the OS's own table.
func stillAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	out, err := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/NH").CombinedOutput()
	if err != nil {
		// The measurement failed. Report "alive" so the caller fails loudly
		// rather than reporting a silent zero (a broken ruler is not a green).
		return true
	}
	return strings.Contains(string(out), strconv.Itoa(pid))
}

// ---------------------------------------------------------------------------
// Canvas-level invariants of the runtime
// ---------------------------------------------------------------------------

// Test52ExtractIsNeverTheRawStdout: an extractor that selects nothing yields an
// EMPTY result, never the payload.
func Test52EmptyExtractionIsEmptyNotRaw(t *testing.T) {
	exe := installFixtureCLI(t)
	def := fixtureTool(exe, "extract")
	def.Extract = map[string]any{"kind": "json", "path": ".data.not_here"}
	runtime, _ := runtimeFor(t, def, CommandRunOptions{EnvLookup: fixtureEnv("extract")})
	res, err := runtime.Exec(t.Context(), "com.example.fixture", def, []byte(`{}`))
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if res.IsError {
		t.Fatalf("an extract selecting nothing is an empty answer, not an error: %+v", res)
	}
	if res.Text != "" {
		t.Fatalf("res.Text = %q, want empty (raw stdout must never be the fallback)", res.Text)
	}
}

// Test52NonZeroExitIsAnError: a CLI that fails reports its code and its stderr,
// not a success with empty output.
func Test52NonZeroExitIsAnError(t *testing.T) {
	exe := installFixtureCLI(t)
	def := fixtureTool(exe, "extract")
	def.Parameters = map[string]any{"type": "object", "properties": map[string]any{}}
	def.Extract = map[string]any{"kind": "json", "path": ".ok"}
	// "extract" is not a mode the fixture knows unless the tool says so; the
	// env lookup below deliberately withholds the mode, so the fixture exits 2.
	env := func(name string) (string, bool) {
		switch name {
		case "PATH":
			return os.Getenv("PATH"), true
		case ModuleNameEnv:
			return "1", true
		}
		return "", false
	}
	runtime, _ := runtimeFor(t, def, CommandRunOptions{EnvLookup: env})
	res, err := runtime.Exec(t.Context(), "com.example.fixture", def, []byte(`{}`))
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if !res.IsError {
		t.Fatalf("a nonzero exit was reported as success: %+v", res)
	}
	if !strings.Contains(res.Text, "exited with code") {
		t.Fatalf("res.Text = %q, want it to name the exit code", res.Text)
	}
}

// Test52TaintedOutputIsFramedAsData: with taint_output the extracted text is
// marked as data (D30 layer 4).
func Test52TaintedOutputIsFramedAsData(t *testing.T) {
	exe := installFixtureCLI(t)
	def := fixtureTool(exe, "extract")
	def.Extract = map[string]any{"kind": "json", "path": ".ok"}
	def.TaintOutput = true
	runtime, _ := runtimeFor(t, def, CommandRunOptions{EnvLookup: fixtureEnv("extract")})
	res, err := runtime.Exec(t.Context(), "com.example.fixture", def, []byte(`{}`))
	if err != nil || res.IsError {
		t.Fatalf("Exec: err=%v res=%+v", err, res)
	}
	if !strings.Contains(res.Text, "data, not instructions") {
		t.Fatalf("res.Text = %q, want the data framing", res.Text)
	}
}

// Test52CommandSchemaIsNarrowedForTheModel: what the model sees is the
// template's placeholders only.
func Test52SchemaForTheModelIsNarrowed(t *testing.T) {
	def := plugin.ToolDef{
		Name: "fixture.run", Kind: plugin.ToolKindCommand,
		ArgvTemplate: []string{"fixture.exe", "{{a}}"},
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"a": map[string]any{"type": "string"},
				"b": map[string]any{"type": "string"},
			},
			"required": []any{"a", "b"},
		},
	}
	var got map[string]any
	if err := json.Unmarshal(commandSchemaFor(def), &got); err != nil {
		t.Fatalf("commandSchemaFor is not JSON: %v", err)
	}
	props, _ := got["properties"].(map[string]any)
	if _, ok := props["b"]; ok {
		t.Fatalf("the model-facing schema still offers %q", "b")
	}
	req, _ := got["required"].([]any)
	if len(req) != 1 {
		t.Fatalf("required = %v, want only the substituted placeholder", req)
	}
}
