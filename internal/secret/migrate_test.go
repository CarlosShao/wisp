package secret

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

// plaintextFixture is the seeded "first run" config: two plaintext api_key
// fields, one pre-existing env ref, an empty api_key, and unrelated fields
// that must survive the rewrite verbatim.
const plaintextFixture = `# user's original config (comments survive in the backup)
schema_version = 1

[app]
language = "zh-CN"
theme = "dark"

[llm]

[llm.providers.openai]
protocol = "openai-chat"
base_url = "https://api.openai.com/v1"
api_key = "sk-openai-plain-1122334455"

[llm.providers.deepseek]
protocol = "openai-chat"
base_url = "https://api.deepseek.com"
api_key_ref = "env:WISP_TEST_LLM_KEY"

[voice.realtime]
enabled = false
api_key = "rt-plain-abcdef9988"

[cost]
alert_threshold = 0.8
`

const openaiSecret = "sk-openai-plain-1122334455"
const realtimeSecret = "rt-plain-abcdef9988"

func writeFixture(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

// TestMigratePlaintextMigratesAndBacksUp is acceptance criterion 2: seeded
// plaintext config -> migrated to dpapi refs, fields removed, backup exists,
// user-visible notice logged, secrets resolvable again via the new refs.
func TestMigratePlaintextMigratesAndBacksUp(t *testing.T) {
	logs := captureLog(t)
	path := writeFixture(t, plaintextFixture)
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	report, err := MigratePlaintext(path)
	if err != nil {
		t.Fatalf("MigratePlaintext: %v", err)
	}

	// Both plaintext fields migrated, under their deterministic dotted paths.
	if len(report.Migrated) != 2 {
		t.Fatalf("Migrated = %v, want the 2 plaintext fields", report.Migrated)
	}
	openaiRef := report.Refs["llm.providers.openai.api_key"]
	realtimeRef := report.Refs["voice.realtime.api_key"]
	if openaiRef == "" || realtimeRef == "" {
		t.Fatalf("Refs incomplete: %v", report.Refs)
	}
	if !strings.HasPrefix(openaiRef, RefPrefixDPAPI) || !strings.HasPrefix(realtimeRef, RefPrefixDPAPI) {
		t.Fatalf("refs must be dpapi refs: %v", report.Refs)
	}

	// Backup exists and is byte-identical to the original (comments included).
	if report.BackupPath == "" {
		t.Fatal("BackupPath empty after a real migration")
	}
	if !strings.HasSuffix(report.BackupPath, "config.toml"+BackupSuffix) {
		t.Errorf("backup name = %q, want config.toml.bak-plaintext", report.BackupPath)
	}
	backup, err := os.ReadFile(report.BackupPath)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if !bytes.Equal(backup, original) {
		t.Error("backup is not byte-identical to the original config")
	}

	// The rewritten config: no plaintext api_key anywhere, refs in place,
	// unrelated values untouched.
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read migrated config: %v", err)
	}
	if bytes.Contains(raw, []byte("api_key = \"")) {
		t.Errorf("migrated config still contains a plaintext api_key:\n%s", raw)
	}
	var tree map[string]any
	if err := toml.Unmarshal(raw, &tree); err != nil {
		t.Fatalf("parse migrated config: %v", err)
	}
	llm := tree["llm"].(map[string]any)["providers"].(map[string]any)
	openai := llm["openai"].(map[string]any)
	if openai[RefKey] != openaiRef {
		t.Errorf("openai api_key_ref = %v, want %q", openai[RefKey], openaiRef)
	}
	if _, has := openai[PlaintextKey]; has {
		t.Error("openai still carries a plaintext api_key key")
	}
	if openai["base_url"] != "https://api.openai.com/v1" {
		t.Errorf("unrelated field lost: openai.base_url = %v", openai["base_url"])
	}
	if tree["app"].(map[string]any)["language"] != "zh-CN" {
		t.Error("unrelated section lost: [app].language")
	}
	realtime := tree["voice"].(map[string]any)["realtime"].(map[string]any)
	if realtime[RefKey] != realtimeRef {
		t.Errorf("realtime api_key_ref = %v, want %q", realtime[RefKey], realtimeRef)
	}
	if llm["deepseek"].(map[string]any)[RefKey] != "env:WISP_TEST_LLM_KEY" {
		t.Error("pre-existing env ref was touched by the migration")
	}
	if _, has := realtime[PlaintextKey]; has {
		t.Error("voice.realtime still carries a plaintext api_key key")
	}

	// The blobs decrypt back to the original secrets (round-trip through the
	// migration), one file per ref under <dir>\secrets\.
	st, err := NewStore(filepath.Dir(path))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	for ref, want := range map[string]string{openaiRef: openaiSecret, realtimeRef: realtimeSecret} {
		got, err := st.Resolve(ref)
		if err != nil {
			t.Fatalf("Resolve(%q): %v", ref, err)
		}
		if got != want {
			t.Errorf("migrated secret for %q does not round-trip", ref)
		}
	}

	// Notice logged, user-visible, and free of secret material.
	out := logs.String()
	if !strings.Contains(out, "migrated") || !strings.Contains(out, "D33") {
		t.Errorf("migration notice missing from the log:\n%s", out)
	}
	if !strings.Contains(out, filepath.Base(report.BackupPath)) {
		t.Errorf("notice should name the backup file:\n%s", out)
	}
	if strings.Contains(out, openaiSecret) || strings.Contains(out, realtimeSecret) {
		t.Errorf("log path contains plaintext secrets:\n%s", out)
	}
	if r := report.Notice(); r == "" || !strings.Contains(r, "D33") || strings.Contains(r, openaiSecret) {
		t.Errorf("report.Notice() = %q, want a D33 notice without secrets", r)
	}
}

// TestMigratePlaintextIdempotent is acceptance criterion 2b: the second run
// is a no-op - no fields, no backup rewrite, no file change, no notice.
func TestMigratePlaintextIdempotent(t *testing.T) {
	path := writeFixture(t, plaintextFixture)
	if _, err := MigratePlaintext(path); err != nil {
		t.Fatalf("first MigratePlaintext: %v", err)
	}
	afterFirst, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read migrated config: %v", err)
	}
	backupAfterFirst, err := os.ReadFile(path + BackupSuffix)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}

	logs := captureLog(t)
	report, err := MigratePlaintext(path)
	if err != nil {
		t.Fatalf("second MigratePlaintext: %v", err)
	}
	if len(report.Migrated) != 0 {
		t.Errorf("second run migrated %v, want nothing", report.Migrated)
	}
	if report.BackupPath != "" {
		t.Errorf("second run wrote a backup: %q", report.BackupPath)
	}
	if report.Notice() != "" {
		t.Errorf("second run produced a notice: %q", report.Notice())
	}
	afterSecond, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config after second run: %v", err)
	}
	if !bytes.Equal(afterFirst, afterSecond) {
		t.Error("second run modified the config file")
	}
	backupAfterSecond, err := os.ReadFile(path + BackupSuffix)
	if err != nil {
		t.Fatalf("read backup after second run: %v", err)
	}
	if !bytes.Equal(backupAfterFirst, backupAfterSecond) {
		t.Error("second run modified the backup")
	}
	if out := logs.String(); strings.Contains(out, "migrated") {
		t.Errorf("second run logged a migration notice:\n%s", out)
	}
}

// TestMigratePlaintextNoConfigIsNoOp: no config yet -> clean no-op.
func TestMigratePlaintextNoConfigIsNoOp(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	report, err := MigratePlaintext(path)
	if err != nil {
		t.Fatalf("MigratePlaintext on missing config: %v", err)
	}
	if len(report.Migrated) != 0 || report.BackupPath != "" {
		t.Errorf("missing config must be a no-op, got %+v", report)
	}
}

// TestMigratePlaintextKeepsFirstBackup: an existing backup is never
// clobbered - the first-seen original wins.
func TestMigratePlaintextKeepsFirstBackup(t *testing.T) {
	path := writeFixture(t, plaintextFixture)
	sentinel := []byte("first original - do not lose")
	if err := os.WriteFile(path+BackupSuffix, sentinel, 0o600); err != nil {
		t.Fatalf("seed backup: %v", err)
	}
	if _, err := MigratePlaintext(path); err != nil {
		t.Fatalf("MigratePlaintext: %v", err)
	}
	got, err := os.ReadFile(path + BackupSuffix)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if !bytes.Equal(got, sentinel) {
		t.Error("existing backup was overwritten")
	}
}

// TestMigratePlaintextLeavesNonSecretsAlone: empty-string api_key values are
// not secrets and unparseable files are explicit errors, not silent rewrites.
func TestMigratePlaintextLeavesNonSecretsAlone(t *testing.T) {
	t.Run("empty api_key is not migrated", func(t *testing.T) {
		path := writeFixture(t, "[llm.providers.x]\napi_key = \"\"\n")
		before, _ := os.ReadFile(path)
		report, err := MigratePlaintext(path)
		if err != nil {
			t.Fatalf("MigratePlaintext: %v", err)
		}
		if len(report.Migrated) != 0 {
			t.Errorf("empty api_key migrated: %v", report.Migrated)
		}
		after, _ := os.ReadFile(path)
		if !bytes.Equal(before, after) {
			t.Error("no-op run modified the file")
		}
	})

	t.Run("unparseable config is an explicit error", func(t *testing.T) {
		path := writeFixture(t, "not [valid toml ===")
		if _, err := MigratePlaintext(path); err == nil {
			t.Fatal("unparseable config must fail explicitly")
		}
	})
}
