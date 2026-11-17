package config

// Ticket 39 AC#1's full section-tree read face (readface wave 1, 2026-11-17):
// one walk over the D36 Config tree that answers `sections[]{key,value,tier}`
// for EVERY section, not only the model/voice half ticket 313 read and the
// ball/panel six the snapshot pushes. The config editor page renders from this
// and from nothing else - the GUI never owns defaults or a second tree (D36
// rule 4), so this walk is the one machine-readable "what is in the file".
//
// Where the tier comes from: TierRegistry (tiers.go) is the source of the tier
// words, consulted per key path first (the [app]/[voice] per-key rows) and per
// section second. A locked section (risk/fs/net/plugins) carries the registry's
// own "locked" marker - the page renders its lock affordance from exactly this
// word, and a made-up fourth word here would be a second vocabulary.
//
// What a Value may contain: config.toml facts only. The schema keeps plaintext
// credentials OUT of Config (D36 rule 5 - only api_key_ref references ride the
// struct; the plaintext lives in Resolved), so a walk over Config cannot leak a
// secret no matter how the walk is written. Composite values (slices, the
// [plugins] dynamic tables) are serialized as compact JSON so the row stays
// one string on the wire.

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// SectionRow is one leaf of the D36 tree: a dotted key path ("ball.size",
// "voice.asr.provider", "plugins.<id>.enabled"), the value as it sits in the
// struct, and the tier word the registry gives that path.
type SectionRow struct {
	Key   string
	Value string
	Tier  string
}

// SectionRows walks the whole Config tree in schema field order (maps sorted
// by key so two walks of one config marshal identically). A nil config
// answers nil - the caller reports "the leg did not read", never an empty tree
// dressed as "the file has nothing in it".
func SectionRows(c *Config) []SectionRow {
	if c == nil {
		return nil
	}
	var out []SectionRow
	walkSectionValue(&out, "", reflect.ValueOf(*c))
	return out
}

// tierForPath answers the tier word for one dotted key path: the per-key
// registry row first, the section row second, "" when neither is registered
// (a caller that receives "" treats it as "unknown tier", never as a tier).
func tierForPath(path string) string {
	if t, ok := TierOf(path); ok {
		return t
	}
	section := path
	if i := strings.IndexByte(path, '.'); i > 0 {
		section = path[:i]
	}
	if t, ok := TierOf(section); ok {
		return t
	}
	return ""
}

// walkSectionValue appends one row per leaf reached from v under prefix.
// Structs descend through their toml tags (`toml:"-"` fields are skipped, the
// same rule the marshal side uses), maps descend per key (sorted), slices and
// scalars are leaves.
func walkSectionValue(out *[]SectionRow, prefix string, v reflect.Value) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			emitSectionRow(out, prefix, "null")
			return
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			tag := f.Tag.Get("toml")
			if tag == "-" && f.Name != "Plugins" {
				// `toml:"-"` is the marshal side's own skip (App.Portable).
				// Config.Plugins carries the same tag (parse.go keeps it out
				// of the plain head marshal) but its canonical segment is
				// "plugins" and the dynamic sub-tables MUST walk - the one
				// deliberate exception, named.
				continue
			}
			name := tag
			if f.Name == "Plugins" {
				name = "plugins"
			}
			if name == "" {
				// Untagged reachable fields: PluginsSection.Entries only.
				// The entries hang directly under the "plugins" prefix the
				// walk already carries - the same shape the diff walk treats
				// them as (a transparent map under the section).
				child := prefix
				walkSectionValue(out, child, v.Field(i))
				continue
			}
			walkSectionValue(out, joinSectionKey(prefix, name), v.Field(i))
		}
	case reflect.Map:
		if prefix == "" {
			return
		}
		keys := make([]string, 0, v.Len())
		for _, k := range v.MapKeys() {
			keys = append(keys, fmt.Sprint(k.Interface()))
		}
		sort.Strings(keys)
		for _, k := range keys {
			walkSectionValue(out, joinSectionKey(prefix, k), v.MapIndex(reflect.ValueOf(k)))
		}
	case reflect.Slice, reflect.Array:
		emitSectionRow(out, prefix, compactJSON(v))
	default:
		emitSectionRow(out, prefix, scalarSectionValue(v))
	}
}

// emitSectionRow appends one row when the prefix names a key. A walk that
// somehow reached a leaf with no name (the empty-prefix root) must not emit a
// nameless row - the page cannot render a keyless slot.
func emitSectionRow(out *[]SectionRow, key, value string) {
	if key == "" {
		return
	}
	*out = append(*out, SectionRow{Key: key, Value: value, Tier: tierForPath(key)})
}

func joinSectionKey(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}

// scalarSectionValue renders one scalar the way the file spells it: strings
// verbatim (a ref travels as a ref), bools as true/false, numbers bare.
func scalarSectionValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'g', -1, 64)
	default:
		return fmt.Sprint(v.Interface())
	}
}

// compactJSON renders one composite leaf as one line of JSON. A slice of
// strings ("a","b") becomes ["a","b"] - the same array the TOML file carries.
func compactJSON(v reflect.Value) string {
	data, err := json.Marshal(v.Interface())
	if err != nil {
		// Slices and arrays of scalars cannot fail to marshal; if a later
		// schema addition ever makes that untrue, the row says so instead of
		// quietly rendering an empty value.
		return "<unrenderable>"
	}
	return string(data)
}
