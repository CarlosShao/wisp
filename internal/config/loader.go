package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	cfg, err := readConfigFile(path)
	if err != nil {
		return nil, nil, err
	}
	resolved, err := resolveRefs(cfg, res)
	if err != nil {
		return nil, nil, err
	}
	return cfg, resolved, nil
}

// readConfigFile runs the sec 4.1 pipeline up to (and including) semantic
// validation, and stops before api_key_ref resolution: it returns the values a
// file declares, not the secrets behind them.
//
// Ticket 226 split it out for one reason. A programmatic write must re-read the
// file it is about to replace so it can merge instead of overwrite - and that
// re-read needs NO secret resolution (it is a read of values, and the refs
// travel back into the file untouched). Routing that re-read through LoadFile
// would make "the DPAPI blob is unreadable right now" refuse a config write that
// has nothing to do with secrets, while NewManager/CheckAndReload keep using
// LoadFile unchanged. One pipeline, one order of failures, two entry points.
func readConfigFile(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, observe.Wrap(observe.ClassConfig, err, "config.toml read")
	}
	ver, peekErr := peekSchemaVersion(raw)
	if ver == 0 && peekErr != nil {
		// Ticket 223 AC#4: when not even a schema version can be picked out of
		// these bytes, the file's shape is unreadable and that is a SYNTAX
		// failure - it must be named as one, not as a migration problem. It used
		// to be unnameable, because this probe discarded its error: the version
		// stayed 0, the pipeline fell into applyMigrations, and a broken line
		// came back reported as "cannot migrate from schema version 1".
		// Measured, not theorized - 票 223's own case planted
		// "schema_version = 2\n\nthis is not toml [[[" and read that sentence;
		// since 票 223 r2 the branch below names that shape correctly too.
		//
		// The ver != 0 half is deliberate and it is not a loosening: a file that
		// declares a version OLDER than this build belongs to the migration
		// pipeline, and TestMigrateCorruptFileUntouched (migrate_test.go:123)
		// requires that path's guidance ("the file was left untouched") for a
		// declared-1 file with a broken table header. Which pipeline owns the
		// error is decided by the version the file declares, not by which
		// message is easier.
		return nil, observe.Wrap(observe.ClassConfig, peekErr, "config.toml parse")
	}
	if peekErr != nil && ver >= SchemaVersionCurrent {
		// Ticket 223 AC#4, second half (r2): the bytes DID declare a version,
		// so branch 1 above does not fire, and the document still does not
		// parse - the failure is the syntax of the body and the operator must
		// hear the 语法错 sentence for it. 223-v1 measured the missing half:
		// four declared-current shapes whose only fault was a broken line
		// (comment-led table bracket, CRLF, tight "schema_version=2", indented
		// version line) fell through to decodeStrict, whose DecodeError text
		// describeReloadFailure books as cause=invalid - "语法没问题" said to a
		// file whose one and only problem is syntax. A file already carrying
		// this build's version (or a newer one) has no migration to run, so a
		// parse failure inside it can only be what it is. Pinned, with the
		// older-version half left to the migration branch above, by
		// TestTicket223R2FailureSentenceRouting in cmd/wisp.
		return nil, observe.Wrap(observe.ClassConfig, peekErr, "config.toml parse")
	}
	if ver > SchemaVersionCurrent {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config.toml: schema_version %d was written by a newer build (this build understands %d); upgrade Wisp or restore a backup",
			ver, SchemaVersionCurrent))
	}
	if ver != SchemaVersionCurrent {
		raw, err = applyMigrations(path, raw, ver)
		if err != nil {
			return nil, err
		}
	}
	cfg := NewDefaults()
	if err := decodeStrict(raw, cfg); err != nil {
		return nil, err
	}
	normalizeConfig(cfg)
	applyPresets(cfg)
	if err := validate(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// peekSchemaVersion reads only the schema_version key (non-strict about the
// rest of the document; the full decode happens after migration may have
// rewritten the file). It reports the probe's own failure: the caller must not
// read "no version could be picked out of these bytes" as "version 0", because
// that fall-through sent a syntactically broken file into the migration path and
// made 语法错 unnameable (ticket 223 AC#4).
func peekSchemaVersion(raw []byte) (int, error) {
	var probe struct {
		SchemaVersion int `toml:"schema_version"`
	}
	err := toml.Unmarshal(raw, &probe)
	if err == nil {
		return probe.SchemaVersion, nil
	}
	// The document does not parse. It may still have DECLARED a version, and
	// that distinction is what decides which pipeline owns the error: a file
	// that says "schema_version = 1" and is broken further down belongs to the
	// migration path, whose guidance names the version and promises the file was
	// left untouched (migrate_test.go:123 pins that sentence for exactly that
	// shape); a file from which no version can be picked at all has a syntax
	// problem, and 票 223 AC#4 requires it to be called that.
	//
	// This is a fallback for the failure path only - a file that parses never
	// reaches it - and it stops at the first table header, because below one the
	// key would be a nested table's, not the document's version.
	if declared, ok := declaredSchemaVersion(raw); ok {
		return declared, err
	}
	return 0, err
}

// declaredSchemaVersion reads a top-level schema_version assignment line by
// line. It is deliberately dumb: it is only ever asked what a file that will not
// parse claimed, and its job is to route the error message, not to validate it.
func declaredSchemaVersion(raw []byte) (int, bool) {
	for _, line := range strings.Split(string(raw), "\n") {
		trimmed := strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(trimmed, "[") {
			return 0, false
		}
		key, rest, found := strings.Cut(trimmed, "=")
		if !found || strings.TrimSpace(key) != "schema_version" {
			continue
		}
		v, err := strconv.Atoi(strings.TrimSpace(rest))
		if err != nil {
			return 0, false
		}
		return v, true
	}
	return 0, false
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
