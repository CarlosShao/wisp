package config

import (
	"fmt"
	"os"
	"strconv"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Schema migration (SPEC-03 sec 4.4): a registry of per-version rewrites;
// loading an older file auto-migrates it through the chain, keeping a
// byte-identical backup first. Unmigratable files are left untouched with an
// explicit error (the caller surfaces Unconfigured + guidance) - silently
// resetting would wipe fs/net allowlists, which is a security bug.
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

// applyMigrations walks raw from ver up to SchemaVersionCurrent. ver 0 (file
// without schema_version) is treated as v1, the last unversioned layout.
// On success the original is backed up as config.toml.bak-<ver> and the
// migrated file atomically replaces path.
func applyMigrations(path string, raw []byte, ver int) ([]byte, error) {
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
	// Dry-check: never write a migrated file that would not load.
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
	backup := path + ".bak-" + strconv.Itoa(from)
	if err := os.WriteFile(backup, raw, 0o600); err != nil {
		return nil, observe.Wrap(observe.ClassConfig, err, "config.toml migration backup write")
	}
	if err := atomicWrite(path, data); err != nil {
		return nil, observe.Wrap(observe.ClassConfig, err, "config.toml migration rewrite")
	}
	return data, nil
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
