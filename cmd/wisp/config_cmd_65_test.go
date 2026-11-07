package main

// Ticket 65's five teeth for the `wisp config` CLI leg (the owner's today
// usable path for the ball style switch): get answers the schema default
// before any file exists, set jarvis lands and reads back, set glass switches
// back, an out-of-vocabulary value is refused by name with config.toml
// byte-identical, and the closed key roster refuses unknown keys. The mutation
// half of tooth five (lifting SetBallStyle's enum guard must turn tooth four
// red) is a leg drill recorded in docs/evidence/s1/65-ball-style-cli.md §3,
// not a resident test - the guard is internal/config's, and this file's job is
// to prove the CLI surface passes refusals through instead of papering over
// them.
//
// Every case drives the REAL entry point cmdConfig over the REAL layout
// resolution, isolated the way the command surface isolates tests: WISP_ENV=
// test + WISP_TEST_DATA_DIR (the identity contract internal/proc documents on
// TestDataDirEnv) pointing at a fresh temp dir. No test touches the operator's
// own %APPDATA%.

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/proc"
)

// withConfigTestEnv65 points the whole command surface at a fresh temp data
// root for one test and returns that root.
func withConfigTestEnv65(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("WISP_ENV", "test")
	t.Setenv(proc.TestDataDirEnv, dir)
	return dir
}

// runConfigCmd65 drives cmdConfig with captured streams and returns
// (exit code, stdout, stderr).
func runConfigCmd65(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, errb bytes.Buffer
	rc := cmdConfig(args, configIO{stdout: &out, stderr: &errb})
	return rc, out.String(), errb.String()
}

// TestConfigLegGetSetGuarded65 is teeth one to three AND the legCovers133 nail
// for the `config` dispatch leg: get on a missing file answers the schema
// default (not a local literal), set jarvis first-run-creates the file and
// reads back jarvis, set glass switches back. The receipt's tier sentence is
// asserted here too, because "立即生效" would be the one lie this leg must
// never print.
func TestConfigLegGetSetGuarded65(t *testing.T) {
	dir := withConfigTestEnv65(t)

	// 牙① get，文件不存在：schema 默认值（default 标签），rc 0，且 stderr 明说
	// 这是默认值不是文件里的值。
	rc, out, errb := runConfigCmd65(t, "get", "ball.style")
	if rc != 0 {
		t.Fatalf("get on missing file: rc=%d stderr=%s", rc, errb)
	}
	if got := strings.TrimSpace(out); got != "glass" {
		t.Fatalf("get on missing file: want schema default glass, got %q (stdout=%q)", got, out)
	}
	if !strings.Contains(errb, "schema 默认值") {
		t.Fatalf("get on missing file must say the value is the schema default on stderr, got: %s", errb)
	}

	// 牙② set jarvis：首建 config.toml + 守卫写，收据点名文件与键值，档位句
	// 说的是重启档；随后 get 读回 jarvis，文件真在数据根。
	rc, out, errb = runConfigCmd65(t, "set", "ball.style", "jarvis")
	if rc != 0 {
		t.Fatalf("set jarvis: rc=%d stderr=%s", rc, errb)
	}
	if !strings.Contains(out, "ball.style = jarvis") {
		t.Fatalf("set receipt must name key and value, got stdout: %s", out)
	}
	if !strings.Contains(out, "重启") {
		t.Fatalf("set receipt must say the CLI write lands at the restart tier, got stdout: %s", out)
	}
	cfgPath := filepath.Join(dir, "config.toml")
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("set on missing file must first-run-create config.toml at %s: %v", cfgPath, err)
	}
	rc, out, errb = runConfigCmd65(t, "get", "ball.style")
	if rc != 0 {
		t.Fatalf("get after set jarvis: rc=%d stderr=%s", rc, errb)
	}
	if got := strings.TrimSpace(out); got != "jarvis" {
		t.Fatalf("get after set jarvis: want jarvis, got %q", got)
	}

	// 牙③ set glass 切回：写回后 get 读回 glass。
	rc, out, errb = runConfigCmd65(t, "set", "ball.style", "glass")
	if rc != 0 {
		t.Fatalf("set glass: rc=%d stderr=%s", rc, errb)
	}
	if !strings.Contains(out, "ball.style = glass") {
		t.Fatalf("set glass receipt must name key and value, got stdout: %s", out)
	}
	rc, out, errb = runConfigCmd65(t, "get", "ball.style")
	if rc != 0 {
		t.Fatalf("get after set glass: rc=%d stderr=%s", rc, errb)
	}
	if got := strings.TrimSpace(out); got != "glass" {
		t.Fatalf("get after set glass: want glass, got %q", got)
	}
}

// TestConfigSetBallStyleInvalidValueRefused65 is tooth four: an
// out-of-vocabulary value is refused by name with the schema vocabulary in the
// refusal, the command exits 1, and config.toml is byte-identical afterwards -
// the guard refuses BEFORE a byte moves, and this test holds it to that.
//
// WHICH GUARD SPEAKS is pinned on purpose (tooth five's mutation half): the
// refusal the CLI passes through must be the setter's own pre-write sentence
// (SetBallStyle's 不在枚举内 + the schema vocabulary), not merely any refusal.
// A second layer would still refuse the write without it - WriteSection runs
// validate(candidate) over the base the write would produce, and validate.go
// judges ball.style against the same ballStyles list - so lifting the
// pre-write guard alone keeps the file safe but changes WHO answers. This
// test goes red exactly there, which is what the mutation drill in
// docs/evidence/s1/65-ball-style-cli.md §3 records.
func TestConfigSetBallStyleInvalidValueRefused65(t *testing.T) {
	dir := withConfigTestEnv65(t)
	if rc, _, errb := runConfigCmd65(t, "set", "ball.style", "glass"); rc != 0 {
		t.Fatalf("seeding set glass: rc=%d stderr=%s", rc, errb)
	}
	cfgPath := filepath.Join(dir, "config.toml")
	before, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	rc, out, errb := runConfigCmd65(t, "set", "ball.style", "紫色")
	if rc != 1 {
		t.Fatalf("invalid value must be refused with rc=1, got rc=%d stdout=%q stderr=%s", rc, out, errb)
	}
	for _, want := range []string{"ball.style", "不在枚举内", "glass|jarvis"} {
		if !strings.Contains(errb, want) {
			t.Fatalf("refusal must be the setter's own named sentence (field, enum wording, schema vocabulary); missing %q in stderr: %s", want, errb)
		}
	}
	after, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatalf("refused write must leave config.toml byte-identical: %d bytes before, %d after", len(before), len(after))
	}
	rc, out, errb = runConfigCmd65(t, "get", "ball.style")
	if rc != 0 {
		t.Fatalf("get after refused write: rc=%d stderr=%s", rc, errb)
	}
	if got := strings.TrimSpace(out); got != "glass" {
		t.Fatalf("value after refused write: want glass, got %q", got)
	}
}

// TestConfigKeyRosterRefused65 holds the CLI to its closed roster: an unknown
// key on get and on set, an unknown subcommand, and a set without a value are
// all usage refusals (rc=2) that name the one key this surface owns, and none
// of them creates or touches a config.toml.
func TestConfigKeyRosterRefused65(t *testing.T) {
	dir := withConfigTestEnv65(t)
	cfgPath := filepath.Join(dir, "config.toml")

	rc, _, errb := runConfigCmd65(t, "get", "llm.providers")
	if rc != 2 || !strings.Contains(errb, "ball.style") {
		t.Fatalf("get of unknown key: want rc=2 with the roster named, got rc=%d stderr=%s", rc, errb)
	}
	rc, _, errb = runConfigCmd65(t, "set", "ball.size", "56")
	if rc != 2 || !strings.Contains(errb, "ball.style") {
		t.Fatalf("set of unknown key: want rc=2 with the roster named, got rc=%d stderr=%s", rc, errb)
	}
	rc, _, _ = runConfigCmd65(t, "set", "ball.style")
	if rc != 2 {
		t.Fatalf("set without a value: want rc=2, got rc=%d", rc)
	}
	rc, _, _ = runConfigCmd65(t, "frobnicate")
	if rc != 2 {
		t.Fatalf("unknown subcommand: want rc=2, got rc=%d", rc)
	}
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		t.Fatalf("refused commands must not create config.toml at %s (stat err=%v)", cfgPath, err)
	}
}

// TestConfigListShowsRoster65: `wisp config list` answers with the owned key,
// the current effective value and the schema's own vocabulary, with and
// without a config.toml on disk.
func TestConfigListShowsRoster65(t *testing.T) {
	withConfigTestEnv65(t)

	rc, out, errb := runConfigCmd65(t, "list")
	if rc != 0 {
		t.Fatalf("list on missing file: rc=%d stderr=%s", rc, errb)
	}
	if !strings.Contains(out, "ball.style = glass") || !strings.Contains(out, "glass|jarvis") {
		t.Fatalf("list must show key, default value and schema vocabulary, got stdout: %s", out)
	}

	if rc, _, errb := runConfigCmd65(t, "set", "ball.style", "jarvis"); rc != 0 {
		t.Fatalf("set jarvis: rc=%d stderr=%s", rc, errb)
	}
	rc, out, errb = runConfigCmd65(t, "list")
	if rc != 0 {
		t.Fatalf("list after set: rc=%d stderr=%s", rc, errb)
	}
	if !strings.Contains(out, "ball.style = jarvis") {
		t.Fatalf("list must show the current file value, got stdout: %s", out)
	}
}
