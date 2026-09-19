package secret

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Migration vocabulary (D33, D36#5, SPEC-03 §5.3): config.toml must never
// carry a plaintext api_key; the first-run migration rewrites it into
// api_key_ref = "dpapi:<blob-id>" with a DPAPI blob under <data>\secrets\ and
// a backup of the original file.
const (
	// PlaintextKey is the config key migrated away from.
	PlaintextKey = "api_key"

	// RefKey replaces it.
	RefKey = "api_key_ref"

	// BackupSuffix turns "config.toml" into "config.toml.bak-plaintext"
	// (ticket 06): the pre-migration original, kept forever.
	BackupSuffix = ".bak-plaintext"

	// migrateTmpSuffix is the temp file used for the atomic rewrite.
	migrateTmpSuffix = ".migrate-tmp"
)

// MigrationReport describes one MigratePlaintext run. It is the user-visible
// material for the D33 notice (no secret values, only paths and field names).
type MigrationReport struct {
	ConfigPath string
	// BackupPath is the pre-migration backup, "" when nothing was migrated.
	BackupPath string
	// Migrated lists the dotted TOML paths that were rewritten.
	Migrated []string
	// Refs maps each migrated path to the dpapi ref that now holds the value.
	Refs map[string]string
}

// Notice renders the user-visible migration notice (D33). Empty when nothing
// was migrated. It contains no secret material by construction.
func (r *MigrationReport) Notice() string {
	if r == nil || len(r.Migrated) == 0 {
		return ""
	}
	return fmt.Sprintf("config security (D33): %d plaintext api_key field(s) in %s were migrated to DPAPI-protected references and removed from the file; the original was backed up to %s",
		len(r.Migrated), r.ConfigPath, r.BackupPath)
}

// plaintextField is one discovered plaintext api_key: the table it lives in,
// and its dotted path (for the report and the deterministic blob id).
type plaintextField struct {
	node   map[string]any
	value  string
	dotted string
}

// MigratePlaintext rewrites plaintext api_key fields in the config file at
// configPath into DPAPI-protected api_key_ref references (D33, SPEC-03 §5.3):
//
//  1. every string api_key != "" anywhere in the tree (nested tables and
//     arrays of tables) is stored as a DPAPI blob under
//     <dir(configPath)>\secrets\<blob-id> - SPEC-02 §6 places config.toml in
//     the data root, so the secrets dir is its sibling;
//  2. the field becomes api_key_ref = "dpapi:<blob-id>" (deterministic blob
//     id derived from the TOML path, so re-runs address the same blob);
//  3. the original file is backed up as config.toml.bak-plaintext (an
//     existing backup is never clobbered - the FIRST original wins);
//  4. a user-visible notice is logged (slog) and returned in the report.
//
// Idempotent: a config without plaintext api_key fields is a no-op (no
// backup, no rewrite, no notice). Empty-string api_key values are not
// secrets and are left alone. If storing any blob fails, the file is left
// untouched. The rewrite re-serializes the TOML tree (comments/formatting are
// preserved only in the backup).
func MigratePlaintext(configPath string) (*MigrationReport, error) {
	raw, err := os.ReadFile(configPath)
	if errors.Is(err, fs.ErrNotExist) {
		// First run without a config yet: nothing to migrate, not an error.
		return &MigrationReport{ConfigPath: configPath, Refs: map[string]string{}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("secret: migrate: read %s: %w", configPath, err)
	}
	var tree map[string]any
	if err := toml.Unmarshal(raw, &tree); err != nil {
		return nil, fmt.Errorf("secret: migrate: parse %s: %w", configPath, err)
	}

	var found []plaintextField
	var walk func(node map[string]any, prefix []string)
	walk = func(node map[string]any, prefix []string) {
		for k, v := range node {
			switch t := v.(type) {
			case map[string]any:
				walk(t, append(append([]string{}, prefix...), k))
			case []any:
				for i, e := range t {
					if m, ok := e.(map[string]any); ok {
						walk(m, append(append([]string{}, prefix...), fmt.Sprintf("%s[%d]", k, i)))
					}
				}
			case string:
				if k == PlaintextKey && t != "" {
					found = append(found, plaintextField{
						node:   node,
						value:  t,
						dotted: joinPath(prefix, k),
					})
				}
			}
		}
	}
	walk(tree, nil)

	report := &MigrationReport{ConfigPath: configPath, Refs: map[string]string{}}
	if len(found) == 0 {
		return report, nil
	}

	// SPEC-02 §6: config.toml sits in the data root; secrets\ is its sibling.
	st, err := NewStore(filepath.Dir(configPath))
	if err != nil {
		return nil, fmt.Errorf("secret: migrate: %w", err)
	}

	// Phase 1: persist every blob first - any failure here leaves the config
	// file untouched.
	for i := range found {
		f := &found[i]
		ref := RefPrefixDPAPI + sanitizeBlobID(f.dotted)
		if err := st.Store(ref, f.value); err != nil {
			return nil, fmt.Errorf("secret: migrate: %s: %w", f.dotted, err)
		}
		f.node[RefKey] = ref
		report.Refs[f.dotted] = ref
		report.Migrated = append(report.Migrated, f.dotted)
	}
	// Phase 2: strip the plaintext fields (all blobs already persisted).
	for i := range found {
		delete(found[i].node, PlaintextKey)
	}

	// Backup the original. An existing backup is never overwritten: the
	// first-seen original is the one worth keeping.
	backupPath := configPath + BackupSuffix
	if _, statErr := os.Stat(backupPath); errors.Is(statErr, fs.ErrNotExist) {
		if werr := os.WriteFile(backupPath, raw, 0o600); werr != nil {
			return nil, fmt.Errorf("secret: migrate: write backup %s: %w", backupPath, werr)
		}
	} else if statErr != nil {
		return nil, fmt.Errorf("secret: migrate: stat backup %s: %w", backupPath, statErr)
	}
	report.BackupPath = backupPath

	// Atomic rewrite: temp file in the same directory, then rename.
	out, err := toml.Marshal(tree)
	if err != nil {
		return nil, fmt.Errorf("secret: migrate: serialize: %w", err)
	}
	tmpPath := configPath + migrateTmpSuffix
	if werr := os.WriteFile(tmpPath, out, 0o600); werr != nil {
		return nil, fmt.Errorf("secret: migrate: write %s: %w", tmpPath, werr)
	}
	if rerr := os.Rename(tmpPath, configPath); rerr != nil {
		_ = os.Remove(tmpPath)
		return nil, fmt.Errorf("secret: migrate: replace %s: %w", configPath, rerr)
	}

	// User-visible notice (D33). The fields logged are paths and counts only.
	slog.Warn("plaintext api_key migrated to DPAPI-protected references (D33)",
		"config", configPath,
		"fields", len(report.Migrated),
		"paths", strings.Join(report.Migrated, ","),
		"backup", report.BackupPath)
	return report, nil
}

// joinPath renders prefix.key in dotted form.
func joinPath(prefix []string, key string) string {
	if len(prefix) == 0 {
		return key
	}
	return strings.Join(prefix, ".") + "." + key
}

// sanitizeBlobID maps a dotted TOML path to a safe blob file name: anything
// outside [A-Za-z0-9._-] (e.g. array indexes) becomes '_', capped at 128
// chars. The mapping is deterministic so re-runs address the same blob.
func sanitizeBlobID(dotted string) string {
	if len(dotted) > 128 {
		dotted = dotted[:128]
	}
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			return r
		case r == '.' || r == '_' || r == '-':
			return r
		default:
			return '_'
		}
	}, dotted)
}
