package plugin

// Ticket 52 — D46 command plugins, the half that lives in internal/plugin:
// install-time pinning, the every-execution recheck, parameter screening, argv
// vector building, the environment allowlist and the declarative extractors.
//
// Every test here is written against the ACCEPTANCE CRITERIA on
// .scratch/wisp/issues/52-d46-command-plugins.md, and the comment above each one
// names which criterion it closes. The runtime half (spawn, timeout, caps) is in
// internal/plugin/manifestslot/command_52_test.go.
//
// The fixtures are REAL files: a byte-tamper test that hashes a string cannot
// prove that a one-byte change on disk is noticed, which is the whole point of
// the pin (SPEC-07 §7 红队: 篡改 exe 一字节拒执行).

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

// fixtureExeDir is a directory OUTSIDE the TEMP roots, because D46 (rightly)
// refuses to pin an exe under TEMP. The fake CLI is a file rather than a real
// executable for the pin tests: nothing here spawns it, and the runtime tests
// (manifestslot) do spawn a real process.
//
// On POSIX t.TempDir() answers the OS scratch root (TMPDIR, defaulting to
// /tmp) - the very shape the windows guard below skips away from, and on a CI
// container that root is recycled per run. Refusing to measure there would
// leave the whole pin family unexecuted on ubuntu (that is exactly what run
// 43362781437's SKIP rows were), so the POSIX branch plants the fixtures in a
// STABLE tree instead: $GITHUB_WORKSPACE on CI (the checkout the job itself
// owns, writable by the runner user), /var/tmp otherwise - neither is named by
// TEMP/TMP/TMPDIR, so underTempRoot accepts it and the D46 rule under test is
// exercised for real. The pin only has to outlive this one test, which it does
// by construction.
func fixtureExeDir(t *testing.T) string {
	t.Helper()
	if runtime.GOOS != "windows" {
		base := strings.TrimSpace(os.Getenv("GITHUB_WORKSPACE"))
		posixRoot := "GITHUB_WORKSPACE"
		if base == "" {
			base = "/var/tmp"
			posixRoot = "/var/tmp"
		}
		dir, err := os.MkdirTemp(base, "wisp-d46-fixture-")
		if err != nil {
			// No stable POSIX tree on this host: say so rather than fail
			// inside the code under test. Named so a ledger row can quote it.
			t.Skipf("no stable POSIX fixture root available (%s: %v); the D46 stability fixtures need a tree outside the TEMP/TMPDIR roots", posixRoot, err)
		}
		if env, root, inTemp := underTempRoot(dir); inTemp {
			// This host's TEMP/TMPDIR answer covers the stable tree too, so
			// the pin rule under test would refuse the fixture by the very
			// rule being measured - say so instead of red-herring.
			os.RemoveAll(dir)
			t.Skipf("every POSIX fixture root lives under %s (%q) on this host; the D46 stability fixtures need a tree outside the TEMP/TMPDIR roots", env, root)
		}
		t.Cleanup(func() { os.RemoveAll(dir) })
		return dir
	}
	base := t.TempDir()
	if _, _, inTemp := underTempRoot(base); inTemp {
		// t.TempDir() honours TMPDIR/TMP/TEMP, so a developer shell that
		// points those at a non-default root would land the fixture inside
		// the very tree D46 refuses. Say so instead of failing inside the
		// code under test.
		t.Skipf("t.TempDir() lives under a TEMP root (%s); the D46 stability fixtures need a stable tree", base)
	}
	return base
}

// writeFixtureExe writes a stand-in CLI whose bytes are known.
func writeFixtureExe(t *testing.T, dir, name, body string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o755); err != nil {
		t.Fatalf("write fixture exe: %v", err)
	}
	return p
}

// commandTool is a schema-valid D46 declaration over a fixture exe.
func commandTool(exe string) ToolDef {
	return ToolDef{
		Name:        "lark.send_message",
		Description: "Send one message.",
		Risk:        RiskL0,
		Kind:        ToolKindCommand,
		Exe:         exe,
		ArgvTemplate: []string{
			filepath.Base(exe), "im", "send", "--to", "{{to}}", "--text", "{{text}}",
		},
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"to":   map[string]any{"type": "string", "maxLength": 32},
				"text": map[string]any{"type": "string", "maxLength": 512},
			},
			"required": []any{"to", "text"},
		},
		Extract: map[string]any{"kind": "json", "path": ".data.message_id"},
	}
}

func mustScreen(t *testing.T, def ToolDef, args string) CommandCall {
	t.Helper()
	call, err := ScreenCommandParams(def, json.RawMessage(args))
	if err != nil {
		t.Fatalf("ScreenCommandParams(%s): %v", args, err)
	}
	return call
}

// ---------------------------------------------------------------------------
// AC "Pinning": install computes the hash, one tampered byte refuses + alerts
// ---------------------------------------------------------------------------

// Test52PinComputesAndVerifiesHash is the AC's first half: installing a command
// tool computes exe_hash from the bytes on disk (the manifest need not declare
// one), and a second verification over unchanged bytes passes.
func Test52PinComputesAndVerifiesHash(t *testing.T) {
	dir := fixtureExeDir(t)
	exe := writeFixtureExe(t, dir, "lark-cli.exe", "v1:this-is-the-real-cli")

	pin, err := PinCommandExe(commandTool(exe), nil)
	if err != nil {
		t.Fatalf("PinCommandExe: %v", err)
	}
	if !strings.HasPrefix(pin.Hash, HashAlgorithm) {
		t.Errorf("pin.Hash = %q, want the %q prefix (same shape as plugin_state.hash)", pin.Hash, HashAlgorithm)
	}
	if len(pin.Hash) != len(HashAlgorithm)+64 {
		t.Errorf("pin.Hash = %q, want sha256 + 64 hex digits", pin.Hash)
	}
	wantSize := int64(len("v1:this-is-the-real-cli"))
	if pin.Size != wantSize {
		t.Errorf("pin.Size = %d, want %d (the byte count the hash was taken over)", pin.Size, wantSize)
	}
	if pin.Path == "" {
		t.Fatal("pin.Path is empty: a pin that does not name the file it hashed cannot be re-verified")
	}
	if err := pin.Verify(); err != nil {
		t.Fatalf("Verify over unchanged bytes: %v", err)
	}
}

// Test52OneByteTamperRefusesAndAlerts is the AC's second half, with a REAL file
// and a ONE-BYTE change: the exact red-team case in SPEC-07 §7.
//
// A `Verify` that hashed a cached string, or a refusal that stayed silent, would
// fail here.
func Test52OneByteTamperRefusesAndAlerts(t *testing.T) {
	dir := fixtureExeDir(t)
	exe := writeFixtureExe(t, dir, "lark-cli.exe", "v1:this-is-the-real-cli")
	pin, err := PinCommandExe(commandTool(exe), nil)
	if err != nil {
		t.Fatalf("PinCommandExe: %v", err)
	}

	// One byte, in place: same size, same mtime second, different content.
	orig, err := os.ReadFile(exe)
	if err != nil {
		t.Fatal(err)
	}
	tampered := make([]byte, len(orig))
	copy(tampered, orig)
	tampered[len(tampered)/2] ^= 0x01
	if err := os.WriteFile(exe, tampered, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := pin.Verify(); !errors.Is(err, ErrExeSwap) {
		t.Fatalf("Verify after a one-byte tamper: err = %v, want ErrExeSwap", err)
	}

	cap := &CaptureAlerter{}
	err = VerifyExeAlert(pin, "com.example.lark", "lark.send_message", cap)
	if !errors.Is(err, ErrExeSwap) {
		t.Fatalf("VerifyExeAlert: err = %v, want ErrExeSwap", err)
	}
	alerts := cap.Alerts()
	if len(alerts) != 1 {
		t.Fatalf("alerts = %d, want exactly 1 (a refusal without an alert is the silent swap this AC forbids)", len(alerts))
	}
	a := alerts[0]
	if a.Kind != "exe" {
		t.Errorf("alert.Kind = %q, want \"exe\" so the panel can render it apart from a manifest tamper", a.Kind)
	}
	if a.PluginID != "com.example.lark" {
		t.Errorf("alert.PluginID = %q", a.PluginID)
	}
	if a.Expected != pin.Hash || a.Actual == "" || a.Actual == a.Expected {
		t.Errorf("alert expected/actual = %q/%q, want the pinned hash and the recomputed one", a.Expected, a.Actual)
	}
}

// Test52ManifestPinMakesItUnpinnable: a pin that names bytes the file does not
// have must ALSO be computed from disk at install. A manifest that ships an
// exe_hash is a manifest asking the host to check it against the file.
func Test52ManifestPinMismatchIsRefusedAsASwap(t *testing.T) {
	dir := fixtureExeDir(t)
	exe := writeFixtureExe(t, dir, "lark-cli.exe", "v1:this-is-the-real-cli")
	def := commandTool(exe)
	def.ExeHash = "sha256:" + strings.Repeat("00", 32) // not the bytes on disk

	if _, err := PinCommandExe(def, nil); !errors.Is(err, ErrExeSwap) {
		t.Fatalf("PinCommandExe with a wrong declared exe_hash: err = %v, want ErrExeSwap", err)
	}

	right, _, err := HashFileBytes(exe)
	if err != nil {
		t.Fatal(err)
	}
	def.ExeHash = right
	pin, err := PinCommandExe(def, nil)
	if err != nil {
		t.Fatalf("PinCommandExe with the correct declared hash: %v", err)
	}
	if !strings.EqualFold(pin.Hash, right) {
		t.Errorf("pin.Hash = %q, want %q", pin.Hash, right)
	}
}

// Test52UnpinnedToolIsRefused: "no pin" is a refusal, never "nothing to verify".
func Test52UnpinnedToolIsRefused(t *testing.T) {
	if err := (ExePin{}).Verify(); !errors.Is(err, ErrNoPin) {
		t.Fatalf("Verify on an empty pin: err = %v, want ErrNoPin", err)
	}
	pins, err := ParseExeHashField("")
	if err != nil || len(pins) != 0 {
		t.Fatalf("ParseExeHashField(\"\") = %v, %v; want no pins and no error", pins, err)
	}
}

// Test52InstallRefusesUnstableExe is the plan's adversarial case verbatim
// (PLAN.md:2678: `exe` 指向一个可变路径 → 必须拒绝安装), split into the spellings
// that make a path movable. Both gates that can see the manifest are checked:
// the C14 validator (loader) and PinCommandExe (installer).
func Test52InstallRefusesUnstableExe(t *testing.T) {
	tmp := t.TempDir()
	cases := []struct {
		name string
		exe  string
		want string
	}{
		{"temp variable", `%TEMP%\x.exe`, "TEMP"},
		{"scratch percent var", `%USERPROFILE%\AppData\Local\Temp\x.exe`, "movable"},
		{"home relative", `~\bin\lark-cli.exe`, "movable"},
		{"relative", `bin/lark-cli.exe`, "absolute"},
		{"bare name", `lark-cli.exe`, "absolute"},
		{"unc share", `\\server\share\lark-cli.exe`, "UNC"},
		{"script needing a shell", filepath.Join(tmp, "lark-cli.cmd"), "executable image"},
		{"batch file", filepath.Join(tmp, "lark-cli.bat"), "executable image"},
		{"powershell", filepath.Join(tmp, "lark-cli.ps1"), "executable image"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			def := commandTool(tc.exe)
			errs := Validate(&Manifest{
				SchemaVersion: CurrentSchemaVersion,
				ID:            "com.example.lark",
				HostAPI:       "^1.0",
				Name:          "Lark",
				Version:       "0.1.0",
				Tools:         []ToolDef{def},
			})
			var found bool
			for _, e := range errs {
				if strings.HasSuffix(e.Field, ".exe") && strings.Contains(e.Reason, tc.want) {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("C14 accepted exe %q (want a .exe refusal mentioning %q); errors: %v", tc.exe, tc.want, errs)
			}
			if _, perr := PinCommandExe(def, nil); !errors.Is(perr, ErrExeUnstable) {
				t.Fatalf("install accepted exe %q: err = %v, want ErrExeUnstable", tc.exe, perr)
			}
		})
	}
}

// Test52AbsoluteTempPathIsRefusedForTheTempReason: the TEMP refusal must be the
// TEMP rule, not an accident of "not absolute". A path under the real TEMP root
// IS absolute, so a test that only asserted "refused" would pass even if the
// stability check were deleted.
//
// Why a probe and not "os.Getenv(TEMP) == \"\"": a hosted windows runner DOES
// answer TEMP (as an 8.3 short spelling whose `~` used to bury this rule - see
// validateExe's comment), so the old env check both skipped here on ubuntu and
// mis-measured there. When the host answers a TEMP/TMP/TMPDIR root, every arm
// runs against that real root; when it answers nothing, the installer arm is
// still measured with a REAL file planted in the recycled tree this platform
// does offer (t.TempDir()) - the same ErrExeUnstable, still naming TEMP - so
// the stability rule is asserted here on every host, never skipped.
func Test52AbsoluteTempPathIsRefusedForTheTempReason(t *testing.T) {
	root := strings.TrimSpace(os.Getenv("TEMP"))
	if root == "" {
		root = strings.TrimSpace(os.Getenv("TMP"))
	}
	if root == "" {
		root = strings.TrimSpace(os.Getenv("TMPDIR"))
	}

	if root != "" {
		exe := filepath.Join(root, "d46-fixture", "lark-cli.exe")
		def := commandTool(exe)

		errs := Validate(&Manifest{
			SchemaVersion: CurrentSchemaVersion,
			ID:            "com.example.lark",
			HostAPI:       "^1.0",
			Name:          "Lark",
			Version:       "0.1.0",
			Tools:         []ToolDef{def},
		})
		var reasons []string
		for _, e := range errs {
			if strings.HasSuffix(e.Field, ".exe") {
				reasons = append(reasons, e.Reason)
			}
		}
		if len(reasons) == 0 {
			t.Fatalf("C14 accepted %q, which lives under TEMP (%s)", exe, root)
		}
		if !strings.Contains(reasons[0], "TEMP") {
			t.Fatalf("first exe refusal = %q, want it to name TEMP (the stability rule), not just absoluteness", reasons[0])
		}
		_, perr := PinCommandExe(def, nil)
		if !errors.Is(perr, ErrExeUnstable) || !strings.Contains(perr.Error(), "TEMP") {
			t.Fatalf("install refusal = %v, want ErrExeUnstable naming TEMP", perr)
		}
		return
	}

	// No TEMP/TMP/TMPDIR answer on this host. On POSIX the recycled tree this
	// platform does offer is still a subject: os.TempDir()'s default (/tmp)
	// is the recycled scratch root and underTempRoot refuses it by name. On
	// WINDOWS the env-driven rule has no subject when no variable answers
	// (os.TempDir would fall through to a user-profile tree the rule does not
	// name), so the case says so instead of measuring nothing.
	if runtime.GOOS == "windows" {
		t.Skip("no TEMP/TMP/TMPDIR answer on this windows host: the env-driven TEMP rule has no subject to measure here (every CI windows runner answers TEMP; measured on run 43362781437)")
	}
	// Plant a REAL file in the platform's recycled default tree - the
	// installer hashes real bytes - and the refusal is the same ErrExeUnstable
	// naming TEMP. A pinned result here would mean the TEMP rule never fired.
	exe := writeFixtureExe(t, t.TempDir(), "lark-cli.exe", "v1:d46-temp-reason")
	pin, err := PinCommandExe(commandTool(exe), nil)
	if err == nil {
		t.Fatalf("PinCommandExe pinned %q (%s): a TEMP-root exe must be refused with ErrExeUnstable naming TEMP", pin.Hash, exe)
	}
	if !errors.Is(err, ErrExeUnstable) {
		t.Fatalf("install refusal = %v, want ErrExeUnstable", err)
	}
	if !strings.Contains(err.Error(), "TEMP") {
		t.Fatalf("install refusal = %v, want it to name TEMP (the stability rule)", err)
	}
}

// Test52ForwardSlashWindowsPathIsPinnable is the SPEC-07 §4 example itself:
// `exe = "C:/Users/me/bin/lark-cli.exe"`. A Windows path spelled with forward
// slashes is absolute, and refusing it would refuse the contract's own sample.
//
// On POSIX the same criterion is measured inside its own platform: one file,
// pinned twice, first through the native spelling, then through a
// slash-rewritten spelling of the same file - the "spelling must not change
// what is pinned" property, which is what the SPEC example pins on Windows.
// The pre-triage shape skipped unexecuted on ubuntu; it now asserts there.
func Test52ForwardSlashWindowsPathIsPinnable(t *testing.T) {
	dir := fixtureExeDir(t)
	exe := writeFixtureExe(t, dir, "lark-cli.exe", "v1:forward-slash")
	if runtime.GOOS != "windows" {
		slashSpelling := filepath.ToSlash(exe)
		native, err := PinCommandExe(commandTool(exe), nil)
		if err != nil {
			t.Fatalf("PinCommandExe(native spelling %q): %v", exe, err)
		}
		slash, err := PinCommandExe(commandTool(slashSpelling), nil)
		if err != nil {
			t.Fatalf("PinCommandExe(%q): %v (the rewritten spelling of the same file must pin the same bytes)", slashSpelling, err)
		}
		if native.Hash != slash.Hash || native.Path != slash.Path {
			t.Errorf("the two spellings of one file pinned differently: %s@%s vs %s@%s", native.Hash, native.Path, slash.Hash, slash.Path)
		}
		return
	}
	forward := filepath.ToSlash(exe)
	back, err := PinCommandExe(commandTool(exe), nil)
	if err != nil {
		t.Fatalf("PinCommandExe(backslash spelling): %v", err)
	}
	fwd, err := PinCommandExe(commandTool(forward), nil)
	if err != nil {
		t.Fatalf("PinCommandExe(%q): %v (this is the spelling SPEC-07 §4 uses)", forward, err)
	}
	if !strings.EqualFold(back.Hash, fwd.Hash) {
		t.Errorf("the two spellings of one file pinned different bytes: %s vs %s", back.Hash, fwd.Hash)
	}
}

// Test52ManifestPinsEveryCommandTool: exe_hash is per TOOL, not per plugin, so a
// plugin wrapping two CLIs pins both (pinning only the first would leave the
// second unverified).
func Test52ManifestPinsEveryCommandTool(t *testing.T) {
	dir := fixtureExeDir(t)
	first := writeFixtureExe(t, dir, "one.exe", "one")
	second := writeFixtureExe(t, dir, "two.exe", "two")

	a := commandTool(first)
	a.Name = "lark.one"
	b := commandTool(second)
	b.Name = "lark.two"
	b.ArgvTemplate = []string{filepath.Base(second), "run"}

	pins, err := PinManifestExe(Manifest{Tools: []ToolDef{a, b}}, nil)
	if err != nil {
		t.Fatalf("PinManifestExe: %v", err)
	}
	if len(pins) != 2 {
		t.Fatalf("pins = %d, want 2 (one per command tool)", len(pins))
	}
	if pins["lark.one"].Hash == pins["lark.two"].Hash {
		t.Error("both tools pinned the same hash: the per-tool keying is not in effect")
	}

	// The plugin_state round trip must preserve both, path included.
	encoded := EncodeExePins(pins)
	back, err := ParseExeHashField(encoded)
	if err != nil {
		t.Fatalf("ParseExeHashField(%q): %v", encoded, err)
	}
	if len(back) != 2 {
		t.Fatalf("round trip kept %d pins, want 2", len(back))
	}
	for name, want := range pins {
		got, ok := back[name]
		if !ok {
			t.Fatalf("round trip lost tool %q", name)
		}
		if got.Hash != want.Hash || !strings.EqualFold(got.Path, want.Path) {
			t.Errorf("round trip of %q = %+v, want %+v", name, got, want)
		}
	}
}

// Test52SingleToolPinCollapsesToTheBareHash keeps ticket 50's plugin_state shape
// for the one-tool case (a bare "sha256:<hex>"), which is what the existing
// installer tests assert.
func Test52SingleToolPinCollapsesToTheBareHash(t *testing.T) {
	pin := ExePin{Path: `C:\bin\cli.exe`, Hash: "sha256:" + strings.Repeat("ab", 32)}
	enc := EncodeExePins(map[string]ExePin{"lark.one": pin})
	if enc != pin.Hash {
		t.Fatalf("EncodeExePins(one tool) = %q, want the bare hash %q", enc, pin.Hash)
	}
	back, err := ParseExeHashField(enc)
	if err != nil {
		t.Fatalf("ParseExeHashField: %v", err)
	}
	if got := back[""]; got.Hash != pin.Hash {
		t.Fatalf("legacy bare-hash form parsed as %+v", got)
	}
}

// ---------------------------------------------------------------------------
// AC "Injection": meta-character and schema-violating params never reach argv
// ---------------------------------------------------------------------------

// Test52MetacharactersAreRejected is the ban list of PLAN.md:2667 verbatim.
// Each character is tested alone, so a screen that dropped one from the set
// fails on that character instead of on a compound string.
func Test52MetacharactersAreRejected(t *testing.T) {
	dir := fixtureExeDir(t)
	def := commandTool(writeFixtureExe(t, dir, "lark-cli.exe", "v1"))
	banned := []string{"|", "&", ";", "$", "`", "<", ">", "\n", "\r"}
	for _, c := range banned {
		t.Run("text/"+map[string]string{"\n": "LF", "\r": "CR"}[c]+c, func(t *testing.T) {
			call, err := ScreenCommandParams(def, json.RawMessage(`{"to":"bob","text":"a`+escapeJSON(c)+`rm -rf /"}`))
			if err == nil {
				t.Fatalf("parameter carrying %q was accepted (call = %+v)", c, call)
			}
			if !errors.Is(err, ErrParamMetachar) {
				t.Fatalf("err = %v, want ErrParamMetachar", err)
			}
		})
		// The same check on the OTHER placeholder, so the screen is per value
		// rather than "the biggest field".
		t.Run("to"+map[string]string{"\n": "LF", "\r": "CR"}[c]+c, func(t *testing.T) {
			_, err := ScreenCommandParams(def, json.RawMessage(`{"to":"a`+escapeJSON(c)+`b","text":"hi"}`))
			if !errors.Is(err, ErrParamMetachar) {
				t.Fatalf("parameter %q carrying %q: err = %v, want ErrParamMetachar", "to", c, err)
			}
		})
	}
	if _, err := ScreenCommandParams(def, json.RawMessage(`{"to":"bob","text":"hello, world"}`)); err != nil {
		t.Fatalf("a benign value was refused: %v", err)
	}
}

func escapeJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b[1 : len(b)-1])
}

// Test52SchemaViolationsAreRejected covers type, maxLength, enum, range, the
// closed-object rule and the required list. Every case must be refused BEFORE a
// process exists: the check is a pure function, so "before" is structural here
// (the runtime test proves the same ordering end to end).
func Test52SchemaViolationsAreRejected(t *testing.T) {
	dir := fixtureExeDir(t)
	base := commandTool(writeFixtureExe(t, dir, "lark-cli.exe", "v1"))

	withProps := func(props map[string]any, required ...string) ToolDef {
		def := base
		req := make([]any, 0, len(required))
		for _, r := range required {
			req = append(req, r)
		}
		def.Parameters = map[string]any{"type": "object", "properties": props, "required": req}
		return def
	}
	// The screen inspects only the parameters the argv template binds
	// (command_argv.go substitutes the placeholders, then screens those
	// values). The base template binds {{to}}/{{text}} alone, so a schema
	// constraint on `days` would never be consulted; the range cases below
	// carry a template that actually binds {{days}} so the constraint is live.
	withArgv := func(def ToolDef, tail ...string) ToolDef {
		def.ArgvTemplate = append(append([]string{}, def.ArgvTemplate...), tail...)
		return def
	}

	cases := []struct {
		name string
		def  ToolDef
		args string
		want error
	}{
		{"wrong type", withProps(map[string]any{"text": map[string]any{"type": "string"}}, "text"), `{"text":42}`, ErrParamSchema},
		{"over maxLength", withProps(map[string]any{"text": map[string]any{"type": "string", "maxLength": 4}}, "text"), `{"text":"12345"}`, ErrParamSchema},
		{"under minLength", withProps(map[string]any{"text": map[string]any{"type": "string", "minLength": 3}}, "text"), `{"text":"ab"}`, ErrParamSchema},
		{"outside enum", withProps(map[string]any{"text": map[string]any{"type": "string", "enum": []any{"a", "b"}}}, "text"), `{"text":"c"}`, ErrParamSchema},
		{"below minimum", withArgv(withProps(map[string]any{"days": map[string]any{"type": "integer", "minimum": 1, "maximum": 30}}, "days"), "--days", "{{days}}"), `{"days":0}`, ErrParamSchema},
		{"above maximum", withArgv(withProps(map[string]any{"days": map[string]any{"type": "integer", "minimum": 1, "maximum": 30}}, "days"), "--days", "{{days}}"), `{"days":31}`, ErrParamSchema},
		{"not an integer", withArgv(withProps(map[string]any{"days": map[string]any{"type": "integer"}}, "days"), "--days", "{{days}}"), `{"days":1.5}`, ErrParamSchema},
		{"missing required", withProps(map[string]any{"text": map[string]any{"type": "string"}}, "text"), `{}`, ErrParamSchema},
		{"undeclared parameter", withProps(map[string]any{"text": map[string]any{"type": "string"}}, "text"), `{"text":"ok","extra":"x"}`, ErrParamSchema},
		{"no type declared", withProps(map[string]any{"text": map[string]any{"maxLength": 4}}, "text"), `{"text":"ok"}`, ErrParamSchema},
		{"unsupported keyword", withProps(map[string]any{"text": map[string]any{"type": "string", "minProperties": 1}}, "text"), `{"text":"ok"}`, ErrParamUnsupported},
		{"arguments not an object", base, `["text"]`, ErrParamSchema},
		{"arguments null", base, `null`, ErrParamSchema},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			call, err := ScreenCommandParams(tc.def, json.RawMessage(tc.args))
			if err == nil {
				t.Fatalf("accepted %s (call = %+v)", tc.args, call)
			}
			if !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
		})
	}
}

// Test52ArgvIsAVectorWithOneArgumentPerSlot is the "argv vector asserted" half
// of the injection AC in the pure layer: a value full of spaces and quotes stays
// ONE argument, and nothing is ever joined or shell-quoted.
func Test52ArgvIsAVectorWithOneArgumentPerSlot(t *testing.T) {
	dir := fixtureExeDir(t)
	exe := writeFixtureExe(t, dir, "lark-cli.exe", "v1")
	def := commandTool(exe)
	args := `{"to":"bob smith","text":"hello world  with  spaces"}`
	argv, err := BuildCommandArgv(def, mustScreen(t, def, args))
	if err != nil {
		t.Fatalf("BuildCommandArgv: %v", err)
	}
	want := []string{exe, "im", "send", "--to", "bob smith", "--text", "hello world  with  spaces"}
	if len(argv) != len(want) {
		t.Fatalf("argv = %q, want %q", argv, want)
	}
	for i := range want {
		if argv[i] != want[i] {
			t.Fatalf("argv[%d] = %q, want %q (argv = %q)", i, argv[i], want[i], argv)
		}
	}
	// The contract's own shape: no element may be a joined command line.
	if strings.Contains(argv[0], " ") {
		t.Errorf("argv[0] = %q carries a space: element 0 is the exe, never text", argv[0])
	}
}

// Test52ArrayPlaceholderExpandsToSeveralArguments and the embedded form's refusal
// of containers: the two shapes a template may use, each with its bound.
func Test52ArrayPlaceholderExpandsToSeveralArguments(t *testing.T) {
	dir := fixtureExeDir(t)
	exe := writeFixtureExe(t, dir, "ffmpeg.exe", "v1")
	def := ToolDef{
		Name:         "media.transcode",
		Description:  "Transcode.",
		Risk:         RiskL0,
		Kind:         ToolKindCommand,
		Exe:          exe,
		ArgvTemplate: []string{filepath.Base(exe), "-i", "{{inputs}}", "-y", "out.mp4"},
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"inputs": map[string]any{
					"type": "array", "maxItems": 3,
					"items": map[string]any{"type": "string", "maxLength": 16},
				},
			},
			"required": []any{"inputs"},
		},
		Extract: map[string]any{"kind": "regex", "pattern": `out\.mp4`},
	}
	argv, err := BuildCommandArgv(def, mustScreen(t, def, `{"inputs":["a.mp4","b.mp4"]}`))
	if err != nil {
		t.Fatalf("BuildCommandArgv: %v", err)
	}
	want := []string{exe, "-i", "a.mp4", "b.mp4", "-y", "out.mp4"}
	if strings.Join(argv, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("argv = %q, want %q", argv, want)
	}
	if _, err := ScreenCommandParams(def, json.RawMessage(`{"inputs":["a","b","c","d"]}`)); err == nil {
		t.Fatal("maxItems 3 accepted 4 items")
	}

	embedded := def
	embedded.ArgvTemplate = []string{filepath.Base(exe), "-i={{inputs}}"}
	call := mustScreen(t, embedded, `{"inputs":["a.mp4"]}`)
	if _, err := BuildCommandArgv(embedded, call); err == nil {
		t.Fatal("an array placeholder was embedded in one argument instead of being refused")
	}
}

// Test52UndeclaredPlaceholderIsNeverArgv: a template placeholder without a
// schema is a C14 error, and a call that omits a needed placeholder is an error
// rather than an empty argv slot (an empty slot silently becomes a different
// command).
func Test52UndeclaredPlaceholderIsNeverArgv(t *testing.T) {
	dir := fixtureExeDir(t)
	def := commandTool(writeFixtureExe(t, dir, "lark-cli.exe", "v1"))
	def.Parameters = map[string]any{
		"type":       "object",
		"properties": map[string]any{"to": map[string]any{"type": "string"}},
	}
	errs := Validate(&Manifest{
		SchemaVersion: CurrentSchemaVersion, ID: "com.example.lark", HostAPI: "^1.0",
		Name: "Lark", Version: "0.1.0", Tools: []ToolDef{def},
	})
	var named bool
	for _, e := range errs {
		if strings.Contains(e.Reason, "no schema is declared") && strings.Contains(e.Field, "text") {
			named = true
		}
	}
	if !named {
		t.Fatalf("C14 accepted {{text}} with no schema; errors: %v", errs)
	}

	// Runtime half: the placeholder exists in the schema but the call omits it.
	def.Parameters["properties"] = map[string]any{
		"to":   map[string]any{"type": "string"},
		"text": map[string]any{"type": "string"},
	}
	call, err := ScreenCommandParams(def, json.RawMessage(`{"to":"bob"}`))
	if err != nil {
		t.Fatalf("ScreenCommandParams: %v", err)
	}
	if _, err := BuildCommandArgv(def, call); err == nil {
		t.Fatal("BuildCommandArgv filled a missing placeholder with an empty string")
	}
}

// Test52PlaceholderInArgv0IsRefused: element 0 is the pinned exe and never text.
func Test52PlaceholderInArgv0IsRefused(t *testing.T) {
	dir := fixtureExeDir(t)
	def := commandTool(writeFixtureExe(t, dir, "lark-cli.exe", "v1"))
	def.ArgvTemplate = []string{"{{exe}}", "run"}
	errs := Validate(&Manifest{
		SchemaVersion: CurrentSchemaVersion, ID: "com.example.lark", HostAPI: "^1.0",
		Name: "Lark", Version: "0.1.0", Tools: []ToolDef{def},
	})
	var found bool
	for _, e := range errs {
		if strings.HasSuffix(e.Field, "argv_template[0]") {
			found = true
		}
	}
	if !found {
		t.Fatalf("C14 accepted a placeholder in argv[0]; errors: %v", errs)
	}
}

// Test52ParameterSchemaIsNarrowedToTheTemplate: the model-facing schema offers
// only what the template substitutes, so the model cannot believe it changed a
// call by filling a property the template ignores.
func Test52ParameterSchemaIsNarrowedToTheTemplate(t *testing.T) {
	dir := fixtureExeDir(t)
	def := commandTool(writeFixtureExe(t, dir, "lark-cli.exe", "v1"))
	def.Parameters["properties"].(map[string]any)["ignored"] = map[string]any{"type": "string"}
	def.Parameters["required"] = []any{"to", "text", "ignored"}

	schema := ParameterSchema(def)
	props, _ := schema["properties"].(map[string]any)
	if _, ok := props["ignored"]; ok {
		t.Fatalf("ParameterSchema still offers %q, which the argv template never substitutes", "ignored")
	}
	for _, want := range []string{"to", "text"} {
		if _, ok := props[want]; !ok {
			t.Fatalf("ParameterSchema dropped %q, which the template does use", want)
		}
	}
	req, _ := schema["required"].([]any)
	if len(req) != 2 {
		t.Fatalf("required = %v, want the two placeholders only", req)
	}
}

// ---------------------------------------------------------------------------
// AC "Env": allowlist only, never the parent environment
// ---------------------------------------------------------------------------

// Test52ChildEnvIsOnlyTheAllowlist proves the child environment is built from
// the allowlist alone. The lookup is sealed: it answers for ANY name it is
// asked, so a build that called os.Environ (or passed anything through) would
// show up as an extra variable here.
func Test52ChildEnvIsOnlyTheAllowlist(t *testing.T) {
	dir := fixtureExeDir(t)
	def := commandTool(writeFixtureExe(t, dir, "lark-cli.exe", "v1"))
	def.EnvAllowlist = []string{"LARK_TOKEN"}

	sealed := func(name string) (string, bool) { return "v:" + name, true }
	env := CommandEnv(def, sealed)

	got := map[string]bool{}
	for _, kv := range env {
		name, _, _ := strings.Cut(kv, "=")
		got[name] = true
	}
	for _, name := range baseEnvNames {
		if !got[name] {
			t.Errorf("child env is missing the contract's base name %q", name)
		}
	}
	if !got["LARK_TOKEN"] {
		t.Error("child env is missing the manifest-declared LARK_TOKEN")
	}
	if len(got) != len(baseEnvNames)+1 {
		t.Errorf("child env carries %d names (%v), want %d base + LARK_TOKEN", len(got), got, len(baseEnvNames))
	}
	for _, secret := range []string{"OPENAI_API_KEY", "ANTHROPIC_API_KEY", "DEEPSEEK_API_KEY", "AWS_SECRET_ACCESS_KEY", "LARK_SECRET"} {
		if got[secret] {
			t.Errorf("child env carries %q; the host's keys must never reach a wrapped CLI", secret)
		}
	}

	// An allowlisted name the host does not have is simply absent.
	absent := CommandEnv(def, func(string) (string, bool) { return "", false })
	if len(absent) != 0 {
		t.Errorf("CommandEnv with an empty host = %v, want no variables at all", absent)
	}
}

// Test52ChildEnvNamesMatchesWhatTheChildGets keeps the audit view and the
// runtime view from drifting.
func Test52ChildEnvNamesMatchesWhatTheChildGets(t *testing.T) {
	def := ToolDef{Kind: ToolKindCommand, EnvAllowlist: []string{"FFMPEG_LOGLEVEL"}}
	names := ChildEnvNames(def)
	seen := map[string]bool{}
	for _, n := range names {
		seen[n] = true
	}
	for _, want := range append(append([]string{}, baseEnvNames...), "FFMPEG_LOGLEVEL") {
		if !seen[want] {
			t.Errorf("ChildEnvNames = %v, want %q in it", names, want)
		}
	}
}

// Test52PathIsTrimmed: PATH travels, but without empty components.
func Test52PathIsTrimmed(t *testing.T) {
	sep := string(filepath.ListSeparator)
	v := "C:\\a" + sep + sep + "C:\\b" + sep
	got := trimPathValue(v)
	if !strings.Contains(got, "C:\\a") || !strings.Contains(got, "C:\\b") {
		t.Fatalf("trimPathValue(%q) = %q, want both real components kept", v, got)
	}
	if strings.Contains(got, sep+sep) {
		t.Fatalf("trimPathValue(%q) = %q, still carries an empty component", v, got)
	}
}

// ---------------------------------------------------------------------------
// Timeout resolution
// ---------------------------------------------------------------------------

// Test52TimeoutResolution: default 15000ms, manifest-overridable, clamped into
// the range the validator guarantees.
func Test52TimeoutResolution(t *testing.T) {
	cases := []struct {
		ms   int
		want int
	}{
		{0, DefaultCommandTimeoutMS},
		{DefaultCommandTimeoutMS, DefaultCommandTimeoutMS},
		{5000, 5000},
		{MinCommandTimeoutMS - 1, MinCommandTimeoutMS},
		{MaxCommandTimeoutMS + 1, MaxCommandTimeoutMS},
	}
	for _, tc := range cases {
		def := ToolDef{Kind: ToolKindCommand, TimeoutMs: tc.ms}
		if got := CommandTimeoutMs(def); got != tc.want {
			t.Errorf("CommandTimeoutMs(%d) = %d, want %d", tc.ms, got, tc.want)
		}
	}
	if DefaultCommandTimeoutMS != 15000 {
		t.Errorf("DefaultCommandTimeoutMS = %d, want 15000 (D46)", DefaultCommandTimeoutMS)
	}
	if CommandStdoutCapBytes != 1<<20 || CommandStderrCapBytes != 1<<20 {
		t.Errorf("stream caps = %d/%d, want 1MB each", CommandStdoutCapBytes, CommandStderrCapBytes)
	}
	if CommandExtractCapBytes != 32*1024 {
		t.Errorf("CommandExtractCapBytes = %d, want 32KB", CommandExtractCapBytes)
	}
}

// Test52TimeoutFieldIsValidated: C14 refuses a budget the runtime would clamp,
// so a manifest cannot silently ask for a number it does not get.
func Test52TimeoutFieldIsValidated(t *testing.T) {
	dir := fixtureExeDir(t)
	exe := writeFixtureExe(t, dir, "lark-cli.exe", "v1")
	for _, ms := range []int{-1, MinCommandTimeoutMS - 1, MaxCommandTimeoutMS + 1} {
		def := commandTool(exe)
		def.TimeoutMs = ms
		errs := Validate(&Manifest{
			SchemaVersion: CurrentSchemaVersion, ID: "com.example.lark", HostAPI: "^1.0",
			Name: "Lark", Version: "0.1.0", Tools: []ToolDef{def},
		})
		var found bool
		for _, e := range errs {
			if strings.HasSuffix(e.Field, ".timeout_ms") {
				found = true
			}
		}
		if !found {
			t.Errorf("C14 accepted timeout_ms = %d; errors: %v", ms, errs)
		}
	}
}

// ---------------------------------------------------------------------------
// Declarative extraction
// ---------------------------------------------------------------------------

// Test52JSONExtractor covers the grammar's shapes on real CLI-shaped payloads.
func Test52JSONExtractor(t *testing.T) {
	const doc = `{"data":{"items":[{"summary":"a","n":1},{"summary":"b","n":2}],"message_id":"om_1"},"has_more":false}`
	cases := []struct {
		path string
		want string
	}{
		{".data.message_id", "om_1"},
		{"$.data.message_id", "om_1"},
		{".data.items[].summary", "[\"a\",\"b\"]"},
		{".data.items[0].n", "1"},
		{`.data.items[]["summary"]`, "[\"a\",\"b\"]"},
		{".has_more", "false"},
		{"..message_id", "om_1"},
		{".data.missing", ""},
	}
	for _, tc := range cases {
		ext, err := NewExtractor(map[string]any{"kind": "json", "path": tc.path})
		if err != nil {
			t.Fatalf("NewExtractor(%q): %v", tc.path, err)
		}
		got, err := ext.Extract([]byte(doc))
		if err != nil {
			t.Fatalf("Extract(%q): %v", tc.path, err)
		}
		if got != tc.want {
			t.Errorf("Extract(%q) = %q, want %q", tc.path, got, tc.want)
		}
	}
}

// Test52RegexExtractor covers whole-match, named group, multiline and the
// no-match-is-empty rule.
func Test52RegexExtractor(t *testing.T) {
	cases := []struct {
		name string
		decl map[string]any
		in   string
		want string
	}{
		{"whole match", map[string]any{"kind": "regex", "pattern": `sha:[0-9a-f]{8}`}, "rev sha:deadbeef ok", "sha:deadbeef"},
		{"named group", map[string]any{"kind": "regex", "pattern": `sha:(?P<hex>[0-9a-f]{8})`, "group": "hex"}, "rev sha:deadbeef", "deadbeef"},
		{"numeric group", map[string]any{"kind": "regex", "pattern": `sha:([0-9a-f]{8})`, "group": "1"}, "rev sha:deadbeef", "deadbeef"},
		{"multiline", map[string]any{"kind": "regex", "pattern": `^id:(\w+)`, "group": "1", "multiline": true}, "x\nid:abc\ny", "abc"},
		{"no match", map[string]any{"kind": "regex", "pattern": `nope`}, "anything", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ext, err := NewExtractor(tc.decl)
			if err != nil {
				t.Fatalf("NewExtractor: %v", err)
			}
			got, err := ext.Extract([]byte(tc.in))
			if err != nil {
				t.Fatalf("Extract: %v", err)
			}
			if got != tc.want {
				t.Errorf("Extract = %q, want %q", got, tc.want)
			}
		})
	}
}

// Test52ExtractorGrammarIsClosed: an unknown kind, a missing required field, a
// keyword the kind does not read, and an invalid pattern are all refusals. A
// typo must never leave a tool with "no extractor" and raw output on the path.
func Test52ExtractorGrammarIsClosed(t *testing.T) {
	bad := []map[string]any{
		{},
		{"kind": "xml"},
		{"path": ".a"},
		{"kind": "json"},
		{"kind": "json", "path": ".a", "pattern": "x"},
		{"kind": "json", "path": ".a["},
		{"kind": "json", "path": "a"},
		{"kind": "regex"},
		{"kind": "regex", "pattern": "("},
		{"kind": "regex", "pattern": "x", "group": "missing"},
		{"kind": "regex", "pattern": "x", "path": ".a"},
	}
	for _, decl := range bad {
		if _, err := NewExtractor(decl); err == nil {
			t.Errorf("NewExtractor(%v) accepted a declaration the grammar does not define", decl)
		}
	}
}

// Test52ExtractorNeverReturnsRawOutputOnFailure: a payload that is not JSON is a
// tool-level error whose text names the extractor, never the bytes.
func Test52ExtractorNeverReturnsRawOutputOnFailure(t *testing.T) {
	ext, err := NewExtractor(map[string]any{"kind": "json", "path": ".data.id"})
	if err != nil {
		t.Fatal(err)
	}
	const payload = "SECRET-TOKEN-abcdef not json at all"
	got, err := ext.Extract([]byte(payload))
	if err == nil {
		t.Fatalf("Extract accepted a non-JSON payload and returned %q", got)
	}
	if strings.Contains(err.Error(), "SECRET-TOKEN") {
		t.Fatalf("the error text carries the payload: %v", err)
	}
	if got != "" {
		t.Fatalf("Extract returned %q alongside an error; the extractor must never hand back raw output", got)
	}
}

// Test52ExtractIsRequired: a command tool without an extractor is refused at
// C14, because raw stdout must never be the fallback.
func Test52ExtractIsRequired(t *testing.T) {
	dir := fixtureExeDir(t)
	def := commandTool(writeFixtureExe(t, dir, "lark-cli.exe", "v1"))
	def.Extract = nil
	errs := Validate(&Manifest{
		SchemaVersion: CurrentSchemaVersion, ID: "com.example.lark", HostAPI: "^1.0",
		Name: "Lark", Version: "0.1.0", Tools: []ToolDef{def},
	})
	var found bool
	for _, e := range errs {
		if strings.HasSuffix(e.Field, ".extract") {
			found = true
		}
	}
	if !found {
		t.Fatalf("C14 accepted a command tool with no extractor; errors: %v", errs)
	}
}

// ---------------------------------------------------------------------------
// C14 surface of the command kind
// ---------------------------------------------------------------------------

// Test52CommandFieldsAreLegalOnlyOnACommandTool: on a declarative tool the D46
// fields keep naming their owner instead of being read and ignored.
func Test52CommandFieldsAreLegalOnlyOnACommandTool(t *testing.T) {
	dir := fixtureExeDir(t)
	exe := writeFixtureExe(t, dir, "lark-cli.exe", "v1")
	raw := `
schema_version = 1
id = "com.example.lark"
host_api = "^1.0"
name = "Lark"
version = "0.1.0"

[[tools]]
name = "lark.search"
description = "d"
risk = "L0"
capabilities = ["net"]
net_allowlist = ["api.example.com"]
exe = "` + strings.ReplaceAll(exe, `\`, `\\`) + `"

[tools.parameters]
type = "object"
`
	if _, err := ParseManifest([]byte(raw)); err == nil {
		t.Fatal("the loader accepted `exe` on a declarative tool")
	} else if !strings.Contains(err.Error(), "ticket 52") {
		t.Fatalf("err = %v, want it to name ticket 52 (the field's owner)", err)
	}
}

// Test52CommandKindValidates is the positive C14 case for the whole D46
// surface, and the negative cases for each field of it.
func Test52CommandKindValidates(t *testing.T) {
	dir := fixtureExeDir(t)
	exe := strings.ReplaceAll(writeFixtureExe(t, dir, "lark-cli.exe", "v1"), `\`, `\\`)
	tmpl := `
schema_version = 1
id = "com.example.lark"
host_api = "^1.0"
name = "Lark"
version = "0.1.0"

[[tools]]
name = "lark.send_message"
description = "Send one message."
risk = "L1"
capabilities = ["net"]
net_allowlist = ["open.feishu.cn"]
kind = "command"
exe = "` + exe + `"
argv_template = ["lark-cli.exe", "im", "send", "--to", "{{to}}"]
timeout_ms = 5000
env_allowlist = ["LARK_TOKEN"]
irreversible = ["send"]

[tools.parameters]
type = "object"
required = ["to"]
[tools.parameters.properties]
to = { type = "string", maxLength = 32 }

[tools.extract]
kind = "json"
path = ".data.message_id"
`
	m, err := ParseManifest([]byte(tmpl))
	if err != nil {
		t.Fatalf("a complete D46 declaration was refused: %v", err)
	}
	def := m.Tools[0]
	if !def.IsCommand() {
		t.Fatal("IsCommand() = false for kind = \"command\"")
	}
	if def.TimeoutMs != 5000 || len(def.EnvAllowlist) != 1 || len(def.Irreversible) != 1 {
		t.Fatalf("decoded D46 fields = %+v", def)
	}

	cases := []struct {
		name   string
		mutate func(string) string
		field  string
	}{
		// "unknown kind" names the AUDIT finding, not the .kind branch of
		// Validate: kind = "shell" does not select the command key table, so
		// the entry's D46 fields are refused as unknown keys (tools[0].exe …)
		// and ParseManifest returns at the audit gate before field validation
		// runs. The manifest is still refused and the error still names the
		// offending tools[0] entry; only the field specificity differs.
		{"unknown kind", func(s string) string { return strings.Replace(s, `kind = "command"`, `kind = "shell"`, 1) }, "tools[0]"},
		{"unknown irreversible class", func(s string) string {
			return strings.Replace(s, `irreversible = ["send"]`, `irreversible = ["reformat-disk"]`, 1)
		}, ".irreversible[0]"},
		{"duplicate irreversible", func(s string) string {
			return strings.Replace(s, `irreversible = ["send"]`, `irreversible = ["send", "send"]`, 1)
		}, ".irreversible[1]"},
		{"bad env name", func(s string) string {
			return strings.Replace(s, `env_allowlist = ["LARK_TOKEN"]`, `env_allowlist = ["LARK-TOKEN"]`, 1)
		}, ".env_allowlist[0]"},
		{"empty argv", func(s string) string {
			return strings.Replace(s, `argv_template = ["lark-cli.exe", "im", "send", "--to", "{{to}}"]`, `argv_template = []`, 1)
		}, ".argv_template"},
		{"empty argv element", func(s string) string {
			return strings.Replace(s, `"--to", "{{to}}"`, `"--to", ""`, 1)
		}, ".argv_template[4]"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseManifest([]byte(tc.mutate(tmpl)))
			if err == nil {
				t.Fatalf("C14 accepted %s", tc.name)
			}
			if !strings.Contains(err.Error(), tc.field) {
				t.Fatalf("err = %v, want it to name %s", err, tc.field)
			}
		})
	}
}

// Test52DuplicateToolNamesRefused: two [[tools]] entries with one name would
// collapse in the slot's byName map, and the loser is the entry nobody notices.
func Test52DuplicateToolNamesRefused(t *testing.T) {
	errs := validateDuplicateToolNames([]ToolDef{{Name: "a.b"}, {Name: "a.b"}})
	if len(errs) != 1 || !strings.Contains(errs[0].Reason, "duplicate") {
		t.Fatalf("validateDuplicateToolNames = %v, want one duplicate finding", errs)
	}
}
