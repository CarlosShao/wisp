package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"strings"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/winsec"
)

// Parsing and canonical serialization.
//
// [plugins] is the one section go-toml cannot decode into a typed struct:
// its per-plugin sub-tables ([plugins.<id>]) are dynamic keys, and the
// decoder has no catch-all map field for struct members. The decode target
// is therefore Config's own field layout with the Plugins field swapped for
// a raw map absorber (built once via reflect.StructOf, so it can never drift
// from the schema), and the typed PluginsSection is reassembled afterwards.
// Marshaling is the inverse: Config.Plugins carries `toml:"-"` and is
// appended manually in canonical form. Both directions are pinned by the
// round-trip tests.

var (
	pluginsAbsorberType = reflect.TypeOf(map[string]any(nil))
	pluginsFieldIdx     = -1
	decodeType          = buildDecodeType()
)

// buildDecodeType mirrors Config with Plugins replaced by a raw map
// absorber tagged `toml:"plugins"`. All other fields are reused verbatim
// from the Config type (same order, same tags), keeping the two in lockstep.
func buildDecodeType() reflect.Type {
	ct := reflect.TypeOf(Config{})
	fields := make([]reflect.StructField, ct.NumField())
	for i := 0; i < ct.NumField(); i++ {
		f := ct.Field(i)
		if f.Name == "Plugins" {
			pluginsFieldIdx = i
			fields[i] = reflect.StructField{
				Name: "Plugins",
				Type: pluginsAbsorberType,
				Tag:  reflect.StructTag(`toml:"plugins"`),
			}
			continue
		}
		fields[i] = f
	}
	return reflect.StructOf(fields)
}

// decodeStrict overlays data onto cfg (which must already carry defaults).
// Unknown keys anywhere outside [plugins.<id>] produce errors naming the key
// and its 1-based line; unknown keys inside a plugin entry are caught during
// the typed entry rebuild and name the key plus the line within the entry.
func decodeStrict(data []byte, cfg *Config) error {
	dv := reflect.New(decodeType).Elem()
	cv := reflect.ValueOf(cfg).Elem()
	for i := 0; i < decodeType.NumField(); i++ {
		if i == pluginsFieldIdx {
			continue // absorber starts nil
		}
		dv.Field(i).Set(cv.Field(i))
	}
	dec := toml.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dv.Addr().Interface()); err != nil {
		return formatDecodeError(err)
	}
	for i := 0; i < decodeType.NumField(); i++ {
		if i == pluginsFieldIdx {
			ps, err := buildPlugins(dv.Field(i).Interface().(map[string]any))
			if err != nil {
				return err
			}
			cv.Field(i).Set(reflect.ValueOf(ps))
			continue
		}
		cv.Field(i).Set(dv.Field(i))
	}
	return nil
}

// buildPlugins reassembles the typed [plugins] section from the raw
// absorber: the fixed tier2_enabled key plus one strict entry decode per
// plugin sub-table.
func buildPlugins(raw map[string]any) (PluginsSection, error) {
	ps := PluginsSection{}
	for k, v := range raw {
		if k == "tier2_enabled" {
			b, ok := v.(bool)
			if !ok {
				return ps, observe.New(observe.ClassConfig, fmt.Sprintf(
					"config.toml: plugins.tier2_enabled must be a boolean, got %T", v))
			}
			ps.Tier2Enabled = b
			continue
		}
		if k == "" {
			return ps, observe.New(observe.ClassConfig, "config.toml: plugin id must not be empty")
		}
		if strings.ContainsAny(k, "/") {
			return ps, observe.New(observe.ClassConfig, fmt.Sprintf(
				"config.toml: plugin id %q must not contain '/'", k))
		}
		m, ok := v.(map[string]any)
		if !ok {
			return ps, observe.New(observe.ClassConfig, fmt.Sprintf(
				"config.toml: plugins.%s must be a table ([plugins.%s]), got %T", k, k, v))
		}
		var entry PluginEntry
		frag, err := toml.Marshal(m)
		if err != nil {
			return ps, observe.Wrap(observe.ClassConfig, err, fmt.Sprintf("config.toml: plugins.%s", k))
		}
		fdec := toml.NewDecoder(bytes.NewReader(frag))
		fdec.DisallowUnknownFields()
		if err := fdec.Decode(&entry); err != nil {
			return ps, observe.Wrap(observe.ClassConfig, err,
				fmt.Sprintf("config.toml: [plugins.%s] entry (line numbers relative to the entry)", k))
		}
		if ps.Entries == nil {
			ps.Entries = make(map[string]PluginEntry, len(raw))
		}
		ps.Entries[k] = entry
	}
	return ps, nil
}

// formatDecodeError turns go-toml errors into D37 config-class errors.
// Strict (unknown-key) errors carry one DecodeError per offending key with
// the full key path and its 1-based line; type errors carry position too.
func formatDecodeError(err error) error {
	var sme *toml.StrictMissingError
	if errors.As(err, &sme) {
		parts := make([]string, 0, len(sme.Errors))
		for _, de := range sme.Errors {
			row, _ := de.Position()
			parts = append(parts, fmt.Sprintf("unknown key %q at line %d",
				strings.Join(de.Key(), "."), row))
		}
		return observe.New(observe.ClassConfig, "config.toml: "+strings.Join(parts, "; "))
	}
	var de *toml.DecodeError
	if errors.As(err, &de) {
		row, col := de.Position()
		msg := strings.TrimPrefix(de.Error(), "toml: ")
		return observe.New(observe.ClassConfig,
			fmt.Sprintf("config.toml: line %d, col %d: %s", row, col, msg))
	}
	return observe.Wrap(observe.ClassConfig, err, "config.toml parse")
}

// MarshalCanonical serializes cfg into the canonical config.toml layout:
// every static field emitted explicitly (no omitempty - the file must always
// round-trip the full struct, D36 rule 4), then the [plugins] section with
// its dynamic per-plugin sub-tables.
func MarshalCanonical(c *Config) ([]byte, error) {
	if c == nil {
		return nil, observe.New(observe.ClassConfig, "marshal: nil config")
	}
	head, err := toml.Marshal(c)
	if err != nil {
		return nil, observe.Wrap(observe.ClassConfig, err, "config.toml marshal")
	}
	tail, err := marshalPlugins(c.Plugins)
	if err != nil {
		return nil, err
	}
	return append(head, tail...), nil
}

// marshalPlugins emits the whole [plugins] section: tier2_enabled plus one
// quoted, correctly ordered sub-table per plugin entry.
func marshalPlugins(ps PluginsSection) ([]byte, error) {
	inner := map[string]any{"tier2_enabled": ps.Tier2Enabled}
	for id, e := range ps.Entries {
		inner[id] = e
	}
	out, err := toml.Marshal(map[string]any{"plugins": inner})
	if err != nil {
		return nil, observe.Wrap(observe.ClassConfig, err, "config.toml [plugins] marshal")
	}
	return out, nil
}

// atomicWrite writes data to path via temp-file + rename so readers never
// observe a partially written config (GUI writes and migration both funnel
// through this).
func atomicWrite(path string, data []byte) error {
	parent := "."
	if i := strings.LastIndexAny(path, `/\`); i >= 0 {
		parent = path[:i]
	}
	tmp, err := os.CreateTemp(parent, ".wisp-config-*.tmp")
	if err != nil {
		return observe.Wrap(observe.ClassConfig, err, "config.toml temp file")
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // no-op once renamed
	// Seal before a single content byte, not after: os.Chmod(0o600) at the end
	// was decoration on Windows (ticket 89, A51① - the mode argument never
	// lands, the parent directory's DACL decides), and rename carries the
	// descriptor over, so the temp file's ACL *is* config.toml's ACL afterwards
	// (same fact secret/migrate.go states for its own pre-rename temp). A wide
	// empty file leaks nothing; a wide file holding the serialized config -
	// endpoints, model paths, and on an un-migrated machine a plaintext
	// api_key - leaks to every local account.
	if err := winsec.SealFile(tmpName); err != nil {
		_ = tmp.Close()
		return observe.Wrap(observe.ClassConfig, err, "config.toml temp file seal")
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return observe.Wrap(observe.ClassConfig, err, "config.toml write")
	}
	if err := tmp.Close(); err != nil {
		return observe.Wrap(observe.ClassConfig, err, "config.toml write")
	}
	if err := os.Rename(tmpName, path); err != nil {
		return observe.Wrap(observe.ClassConfig, err, "config.toml replace")
	}
	return nil
}

// fileMissing reports whether err is a not-exist condition.
func fileMissing(err error) bool {
	return errors.Is(err, fs.ErrNotExist)
}
