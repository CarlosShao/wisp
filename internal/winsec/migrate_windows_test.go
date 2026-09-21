//go:build windows

// Ticket 89's returned-for-fix item 2: the worst product in internal/secret is
// not the DPAPI blob, it is `config.toml.bak-plaintext` - the *pre-migration*
// config, which by construction still carries every plaintext api_key the user
// had written down. The migration sealed the file it rewrites and left the
// plaintext copy behind with a decorative 0o600 (AC#1 measured what that means
// on this platform: the parent directory's ACL wins, and on this machine that
// ACL grants two other local accounts MODIFY).
//
// Both readings are taken with icacls, before and after, on the same path, so
// the claim is "this object's foreign principals disappeared", not "a call
// returned nil".
package winsec_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/secret"
)

// plaintextConfigFixture is a config as users wrote it before D33: two
// plaintext keys, one env ref, and unrelated fields.
const plaintextConfigFixture = `schema_version = 1

[app]
language = "zh-CN"

[llm.providers.openai]
protocol = "openai-chat"
api_key = "sk-not-a-real-key-0123456789"

[voice.realtime]
enabled = false
api_key = "rt-not-a-real-key-9876543210"
`

// writeWideConfig places config.toml the way internal/config writes it today:
// a bare mode argument, no descriptor, inside a parent that carries a foreign
// read grant. The parent is not winsec's to seal (PrivateDirAll deliberately
// leaves ancestors alone), which is what keeps this a real baseline instead of a
// strawman.
func writeWideConfig(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(plaintextConfigFixture), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// assertCarriesForeign is the pre-condition leg: a criterion that asserts "no
// foreign principal" on an object that never had one is measuring air, so the
// wide state is checked first and the test refuses to run otherwise.
func assertCarriesForeign(t *testing.T, path string) {
	t.Helper()
	sids, names := aclSIDs(t, path)
	if msg := privateACLError(sids, currentSID(t)); msg == "" {
		t.Fatalf("pre-condition: %s is already private, so the seal below would prove nothing (principals %v)",
			filepath.Base(path), names)
	} else {
		t.Logf("BEFORE seal, %s carries %s", filepath.Base(path), msg)
	}
}

// TestAC2MigrationBackupIsPrivate is the AC#2 coverage claim applied to the
// worst artifact in the package: the migration runs inside a directory that
// hands Everyone read, and the backup of the plaintext config must come out
// named {me, SY, BA} and nothing else.
func TestAC2MigrationBackupIsPrivate(t *testing.T) {
	base := wideParent(t, "root")
	configPath := writeWideConfig(t, base)

	// The config the migration reads is wide - that is the defect class, and
	// pinning it here is what stops the assertion below from passing by
	// construction. (config.toml itself is internal/config's write path, ticket
	// 95's account, not this ticket's.)
	assertCarriesForeign(t, configPath)

	report, err := secret.MigratePlaintext(configPath)
	if err != nil {
		t.Fatalf("MigratePlaintext: %v", err)
	}
	if report.BackupPath == "" {
		t.Fatal("no backup path in the report, so there is nothing to judge")
	}

	// AFTER: the pre-migration plaintext, measured rather than assumed.
	assertNoForeignPrincipalText(t, report.BackupPath)
	assertPrivateACL(t, report.BackupPath)

	// The atomic-rewrite temp file is sealed for the same reason and it is the
	// file that *becomes* config.toml: a rename keeps the descriptor, so the
	// migrated config inherits the tmp's privacy. Without the seal on the tmp,
	// the renamed config keeps the Everyone it inherited at creation.
	if _, err := os.Lstat(configPath + ".migrate-tmp"); err == nil {
		t.Errorf("the migration temp file outlived the rename: %s", configPath+".migrate-tmp")
	}
	assertPrivateACL(t, configPath)

	// Everything the migration created under the data root, recursively.
	sweepPrivate(t, filepath.Join(base, "secrets"))
}

// TestAC2PreExistingMigrationBackupIsRepaired covers the branch that does *not*
// write: an existing backup is never clobbered, so on a machine that migrated
// before this fix landed the plaintext copy stays on disk forever. The repair
// leg has to narrow it, and it has to carry a grant of its own for that leg to
// be load-bearing (an inherited-only grant would be recomputed by the OS when
// the parent changes, and the test would pass with the fix deleted).
func TestAC2PreExistingMigrationBackupIsRepaired(t *testing.T) {
	base := wideParent(t, "root")
	configPath := writeWideConfig(t, base)
	backup := configPath + secret.BackupSuffix
	if err := os.WriteFile(backup, []byte("the original, written before winsec existed"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Somebody's explicit grant on the backup itself - a share fix, a support
	// agent that needed to read it once.
	run(t, "icacls", backup, "/grant", "*"+everyoneSID+":(RX)")

	assertCarriesForeign(t, backup)
	if !rawNamesEveryone(t, backup) {
		t.Fatal("pre-condition: the backup carries no explicit Everyone grant")
	}

	before, err := os.ReadFile(backup)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := secret.MigratePlaintext(configPath); err != nil {
		t.Fatalf("MigratePlaintext over an existing backup: %v", err)
	}

	// Same path, same bytes, narrowed descriptor.
	after, err := os.ReadFile(backup)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Errorf("the repair changed the backup's contents: the first-seen original must win")
	}
	assertNoForeignPrincipalText(t, backup)
	assertPrivateACL(t, backup)
}
