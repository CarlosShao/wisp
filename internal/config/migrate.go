package config

import (
	"fmt"
	"os"
	"strconv"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/winsec"
)

// Schema migration (SPEC-03 sec 4.4): a registry of per-version rewrites.
//
// Ticket 345 (2026-11-02, owner-approved shape A) split the pipeline in two:
//
//   - READ (migrateRaw): loading an older file migrates it through the chain
//     IN MEMORY and hands the migrated bytes to the decoder. The read
//     pipeline writes no bytes - the operator's file, its comments and its
//     hand-chosen key order, stay on disk exactly as left.
//   - WRITE (backupOldSchema, called by SaveFile): the migration lands on an
//     explicit write, which was always going to replace the file anyway
//     (D36 rule 3: SaveFile stamps the current version). Before it does, the
//     pre-write bytes are kept verbatim as config.toml.bak-<ver>.
//
// Before this split, "read config.toml once" was itself a disk write whenever
// the file was behind on schema_version (backup + canonical rewrite inside
// readConfigFile) - ticket 345's named root of the lost-comments reports.
//
// Unmigratable files are left untouched with an explicit error (the caller
// surfaces Unconfigured + guidance) - silently resetting would wipe fs/net
// allowlists, which is a security bug. The error sentences below say "the
// file was left untouched"; since ticket 345 the read path means it by
// construction, and the write path refuses to run backupOldSchema over a
// file whose migrated bytes would not load (SaveFile only ever writes an
// already-validated in-memory config, and mergeWrite re-reads first).
//
// Safety invariant: the migrated bytes are decoded and validated BEFORE the
// file is replaced, so a migration can never write a config that would not
// load.

// migrationStep rewrites raw from version From to version To.
type migrationStep struct {
	To int
	Fn func(raw []byte) ([]byte, error)
}

// migrations is the registry: migrations[v] upgrades a v-file one step.
// v1 - D36 draft: [llm] default_provider/fallback_provider/timeout/
//
//	temperature + flat providers.<n>.{base_url, model, api_key_ref}.
//
// v2 - catalog supplement (SPEC-03 sec 3.1); see SchemaVersionCurrent.
var migrations = map[int]migrationStep{
	1: {To: 2, Fn: migrateV1toV2},
}

// migrateRaw walks raw from ver up to SchemaVersionCurrent and returns the
// migrated bytes. ver 0 (file without schema_version) is treated as v1, the
// last unversioned layout. It is PURE (ticket 345): nothing is read from or
// written to the disk here - the backup and the rewrite that used to ride
// along moved to backupOldSchema, which SaveFile runs on an explicit write.
func migrateRaw(raw []byte, ver int) ([]byte, error) {
	from := ver
	if from == 0 {
		from = 1
	}
	data := raw
	cur := from
	for cur != SchemaVersionCurrent {
		step, ok := migrations[cur]
		if !ok {
			return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
				"config.toml: no migration registered from schema version %d (this build understands up to %d); "+
					"the file was left untouched - fix or restore it manually, it will never be silently reset",
				cur, SchemaVersionCurrent))
		}
		next, err := step.Fn(data)
		if err != nil {
			return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
				"config.toml: cannot migrate from schema version %d: %v; "+
					"the file was left untouched - fix or restore it manually, it will never be silently reset",
				cur, err))
		}
		cur = step.To
		data = next
	}
	// Dry-check: a migrated config that would not load is refused here, on
	// the read side, before any write could carry it anywhere.
	check := NewDefaults()
	if err := decodeStrict(data, check); err != nil {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: migration to schema version %d produced an invalid config (%v); the file was left untouched",
			SchemaVersionCurrent, err))
	}
	normalizeConfig(check)
	applyPresets(check)
	if err := validate(check); err != nil {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: migration to schema version %d produced an invalid config (%v); the file was left untouched",
			SchemaVersionCurrent, err))
	}
	return data, nil
}

// backupOldSchema keeps the pre-write bytes of an older-schema config.toml
// verbatim as <path>.bak-<ver>, immediately before SaveFile replaces them -
// the write-side half of what applyMigrations used to run on every read
// (ticket 345). Pure reads never reach it.
//
// Three cases deliberately back up nothing:
//   - no file at path (first-run): there are no bytes to keep;
//   - no version can be picked out of the bytes: the read pipeline owns that
//     report (ticket 223 AC#4 routing), and writing a .bak-<ver> would invent
//     a version number for a file that never declared one;
//   - the file already declares the current or a newer version: no migration
//     is being taken over by this write.
//
// The backup can still carry plaintext api_key values (secret.MigratePlaintext
// exists precisely because config.toml may hold them). A mode argument does
// not land on Windows (ticket 89, A51①), so 0o600 here was decoration: the
// bytes landed with whatever the parent directory's DACL hands down, which is
// readable by other local accounts. winsec.PrivateFile seals the handle before
// the first content byte and refuses the write if the seal fails; because this
// call rewrites an existing backup rather than skipping it, PrivateFile is
// also the repair path for a wide backup left by an older build.
func backupOldSchema(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil // missing or unreadable: nothing verbatim to keep; the atomicWrite that follows answers for the path
	}
	ver, peekErr := peekSchemaVersion(raw)
	if peekErr != nil {
		return nil
	}
	from := ver
	if from == 0 {
		from = 1 // the last unversioned layout, the same reading migrateRaw uses
	}
	if from >= SchemaVersionCurrent {
		return nil
	}
	backup := path + ".bak-" + strconv.Itoa(from)
	if err := winsec.PrivateFile(backup, raw, 0o600); err != nil {
		return observe.Wrap(observe.ClassConfig, err, "config.toml migration backup write")
	}
	return nil
}

// migrateV1toV2 maps the D36 draft model onto the catalog schema (SPEC-03
// sec 3.1):
//
//	llm.timeout                    -> llm.timeout_ms
//	llm.temperature                -> llm.roles.chat.temperature
//	llm.default_provider (+model)  -> llm.text_chain[0], llm.roles.chat.{provider,model}
//	llm.fallback_provider (+model) -> llm.text_chain[1]
//	providers.<n>.model            -> providers.<n>.models.<model>{enabled=true}
//
// Everything else (all other sections, all other keys) passes through
// unchanged.
func migrateV1toV2(raw []byte) ([]byte, error) {
	var root map[string]any
	if err := toml.Unmarshal(raw, &root); err != nil {
		return nil, fmt.Errorf("not valid TOML: %w", err)
	}
	llm, _ := root["llm"].(map[string]any)
	if llm != nil {
		if v, ok := llm["timeout"]; ok {
			llm["timeout_ms"] = v
			delete(llm, "timeout")
		}
		temp, hasTemp := llm["temperature"]
		dp, _ := llm["default_provider"].(string)
		fp, _ := llm["fallback_provider"].(string)
		providers, _ := llm["providers"].(map[string]any)

		// Pass 1: hoist each flat provider model into the catalog.
		modelOf := map[string]string{}
		for name, pv := range providers {
			pm, ok := pv.(map[string]any)
			if !ok {
				continue
			}
			model, _ := pm["model"].(string)
			if model == "" {
				continue
			}
			modelOf[name] = model
			models, _ := pm["models"].(map[string]any)
			if models == nil {
				models = map[string]any{}
			}
			if _, exists := models[model]; !exists {
				models[model] = map[string]any{"enabled": true}
			}
			pm["models"] = models
			delete(pm, "model")
		}

		// Pass 2: chain and chat role.
		var chain []string
		if dp != "" && modelOf[dp] != "" {
			chain = append(chain, dp+"/"+modelOf[dp])
		}
		if fp != "" && fp != dp && modelOf[fp] != "" {
			chain = append(chain, fp+"/"+modelOf[fp])
		}
		if len(chain) > 0 {
			llm["text_chain"] = chain
		}
		if hasTemp || (dp != "" && modelOf[dp] != "") {
			roles, _ := llm["roles"].(map[string]any)
			if roles == nil {
				roles = map[string]any{}
				llm["roles"] = roles
			}
			chat, _ := roles["chat"].(map[string]any)
			if chat == nil {
				chat = map[string]any{}
				roles["chat"] = chat
			}
			if hasTemp {
				chat["temperature"] = temp
			}
			if dp != "" && modelOf[dp] != "" {
				chat["provider"] = dp
				chat["model"] = modelOf[dp]
			}
		}
		delete(llm, "default_provider")
		delete(llm, "fallback_provider")
		delete(llm, "temperature")
	}
	root["schema_version"] = SchemaVersionCurrent
	out, err := toml.Marshal(root)
	if err != nil {
		return nil, fmt.Errorf("re-marshal: %w", err)
	}
	return out, nil
}
