package config

import (
	"fmt"
	"os"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Load pipeline (SPEC-03 sec 4.1): read file -> schema_version check ->
// (migrate with backup when older) -> strict decode onto struct-tag defaults
// -> preset inheritance -> semantic validation -> api_key_ref resolution.
// Plaintext secrets never enter the Config struct; they live in Resolved.

// SecretResolver resolves api_key_ref values to plaintext secrets (SPEC-03
// sec 4.1: DPAPI decrypt or env read). *secret.Store satisfies it; tests
// inject fakes.
type SecretResolver interface {
	Resolve(ref string) (string, error)
}

// Resolved carries the plaintext secrets of one loaded config, keyed for the
// consumers (llm adapters per provider name, voice realtime). The Config
// keeps only refs (storage boundary, SPEC-03 sec 3.1).
type Resolved struct {
	// ProviderKeys maps provider name -> api key (only providers with a
	// non-empty api_key_ref appear).
	ProviderKeys map[string]string
	// RealtimeKey is the voice.realtime api key ("" when unset).
	RealtimeKey string
}

// LoadFile loads and validates path per the sec 4.1 pipeline. res may be nil
// (no resolution; refs stay untouched in the Config). A resolution failure
// returns an error - the caller maps it to the Unconfigured state (sec 4.1:
// never fall back to a half-configured runtime silently).
func LoadFile(path string, res SecretResolver) (*Config, *Resolved, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, observe.Wrap(observe.ClassConfig, err, "config.toml read")
	}
	ver := peekSchemaVersion(raw)
	if ver > SchemaVersionCurrent {
		return nil, nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: schema_version %d was written by a newer build (this build understands %d); upgrade Wisp or restore a backup",
			ver, SchemaVersionCurrent))
	}
	if ver != SchemaVersionCurrent {
		raw, err = applyMigrations(path, raw, ver)
		if err != nil {
			return nil, nil, err
		}
	}
	cfg := NewDefaults()
	if err := decodeStrict(raw, cfg); err != nil {
		return nil, nil, err
	}
	normalizeConfig(cfg)
	applyPresets(cfg)
	if err := validate(cfg); err != nil {
		return nil, nil, err
	}
	resolved, err := resolveRefs(cfg, res)
	if err != nil {
		return nil, nil, err
	}
	return cfg, resolved, nil
}

// peekSchemaVersion reads only the schema_version key (non-strict; the full
// decode happens after migration may have rewritten the file).
func peekSchemaVersion(raw []byte) int {
	var probe struct {
		SchemaVersion int `toml:"schema_version"`
	}
	_ = toml.Unmarshal(raw, &probe)
	return probe.SchemaVersion
}

// applyPresets fills protocol/base_url from the built-in preset for any
// provider leaving them empty (SPEC-03 sec 3: presets supply protocol and
// base_url defaults only). Explicit values are never overwritten.
func applyPresets(c *Config) {
	for name, p := range c.LLM.Providers {
		pre, ok := LookupPreset(name)
		if !ok {
			continue
		}
		if p.Protocol == "" || p.BaseURL == "" {
			if p.Protocol == "" {
				p.Protocol = pre.Protocol
			}
			if p.BaseURL == "" {
				p.BaseURL = pre.BaseURL
			}
			c.LLM.Providers[name] = p
		}
	}
}

// resolveRefs resolves every non-empty api_key_ref through res. A failure
// names the failing ref and provider; nothing partial is returned.
func resolveRefs(c *Config, res SecretResolver) (*Resolved, error) {
	out := &Resolved{ProviderKeys: map[string]string{}}
	if res == nil {
		return out, nil
	}
	for name, p := range c.LLM.Providers {
		if p.APIKeyRef == "" {
			continue
		}
		key, err := res.Resolve(p.APIKeyRef)
		if err != nil {
			return nil, observe.Wrap(observe.ClassConfig, err, fmt.Sprintf(
				"config.toml: llm.providers.%s.api_key_ref (%s)", name, p.APIKeyRef))
		}
		out.ProviderKeys[name] = key
	}
	if ref := c.Voice.Realtime.APIKeyRef; ref != "" {
		key, err := res.Resolve(ref)
		if err != nil {
			return nil, observe.Wrap(observe.ClassConfig, err, fmt.Sprintf(
				"config.toml: voice.realtime.api_key_ref (%s)", ref))
		}
		out.RealtimeKey = key
	}
	return out, nil
}

// SaveFile serializes c in canonical form and atomically replaces path. The
// written file always carries the current schema version (D36 rule 3: the
// struct is the truth; the file is its serialization).
func SaveFile(path string, c *Config) error {
	if c == nil {
		return observe.New(observe.ClassConfig, "save: nil config")
	}
	cp := deepCopyConfig(c)
	cp.SchemaVersion = SchemaVersionCurrent
	data, err := MarshalCanonical(cp)
	if err != nil {
		return err
	}
	return atomicWrite(path, data)
}
