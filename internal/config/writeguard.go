package config

// The guarded write path (ticket 226).
//
// WHY THIS FILE EXISTS. Ticket 201's 「以后这类都别问」 branch, and ticket 90's
// permission-mode switch before it, both persisted their one key by handing the
// WHOLE in-memory snapshot to SaveFile. SaveFile is a serializer, not a merge:
// deepCopyConfig -> MarshalCanonical -> atomicWrite, with no step that ever reads
// the file back (loader.go, unchanged by this file). So every programmatic write
// rewrote config.toml from the values this process loaded at startup - and with
// nothing polling the file in production (ticket 223 measures that gap), the
// snapshot is by construction older than the file the moment the operator edits
// it. The result was the shape ticket 226 was filed for: a hand edit of ANY key
// (a model name, a volume, a hotkey) is silently reverted by one chat answer,
// and the write then adopted the file's new mtime+size as "my own", so even the
// reload path could not see what it had just destroyed.
//
// WHICH CANDIDATE THIS IS. Ticket 226 lists 甲 (re-read and replace only that
// section), 乙 (compare the stat, refuse on divergence) and 丙 (three-way merge).
// This is 甲, chosen over 乙 on two measured costs, both read off this repo
// rather than theorized:
//
//   - 乙 still has to name the conflicting keys (its own AC#2 half), and nothing
//     in this package can diff a Config today, so 乙 pays the same differ this
//     file pays. It buys no extra protection for that money.
//   - 乙 cannot recover inside one run. With no poll wired, refusing leaves the
//     divergence in place forever: every later 「一直」 answer would be refused
//     until the process restarts. The gate the ticket wants closed is
//     "do not overwrite", and 甲 closes it without switching a working answer
//     into a permanently unusable one.
//
// 丙 (merge the two sides key by key) was NOT taken: it would have to decide, per
// key, whether memory or disk wins, and every one of those decisions is a silent
// behaviour this ticket is not authorized to invent. What 丙 does have that 甲
// also delivers here is AC#2's report, so the one precondition the ticket put on
// 丙 is met by this file anyway - the write names every key it actually changed.
//
// THE CONTRACT THIS FILE KEEPS FROM THE OLD PATH. The two safety behaviours
// allowdirs.go carried are untouched and are not this file's to soften: an entry
// with a relative hop never reaches the writer at all (fail-closed), and a failed
// write rolls the in-memory value back so the value the runtime acts on and the
// value the next start reads cannot diverge. What changed is only the COVERAGE of
// the persist - which key of the file this write may touch.
//
// WHAT 甲 COSTS, STATED PLAINLY (the ticket asks for one sentence, not a
// silence): MarshalCanonical is a struct serialization, so it never carried
// comments and never will - a guarded write re-lays the file in canonical form
// and any comment or hand-chosen key order in it is gone. VALUES are not lost any
// more (that is AC#1); formatting is. The other loss is narrower and deliberate:
// because SaveFile always writes the current schema version, a file the loader
// had to migrate is re-laid out at that version by this write - the same thing
// the old path did, and the reason readConfigFile below runs migrations first.

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// The key paths a guarded caller owns. Spelled the way manager.go's direction
// auditor spells them, so a log line here and a D36 verdict there name one thing
// (the same convention unwired.go's table documents).
const (
	keyFSAllowedDirs      = "fs.allowed_dirs"
	keyRiskPermissionMode = "risk.permission_mode"
)

// maxReportedKeys bounds how many key paths one log line may carry. A diverged
// config.toml can differ in hundreds of leaves (a file written by an older build
// is the ordinary case); the bound keeps the report readable and says so rather
// than truncating quietly.
const maxReportedKeys = 24

// mergeWrite persists one programmatic change to the watched config.toml by
// re-reading the file, replacing ONLY the key this caller owns on what the file
// actually holds, and writing that back. It never hands the in-memory snapshot to
// SaveFile as the base again. Caller holds m.mu.
//
// ownedKey is the dotted key path this caller owns; setOn stamps that key onto
// the base it is given. Both are the caller's business - mergeWrite cannot tell
// one key from another and does not try to.
//
// Three outcomes, none of them silent:
//
//  1. The file cannot be re-read (missing is fine, unreadable or invalid is not).
//     A file this process cannot parse is a file it must not overwrite: the edit
//     that made it unparseable is exactly the thing AC#1 forbids destroying, and
//     guessing a base from the snapshot is what this ticket filed. The error is
//     returned and nothing is written.
//  2. The file already carries this key's new value. Nothing is written at all -
//     "we rewrote the file to say what it already said" would be a clobber with
//     no benefit (it re-lays the whole file for nothing) - and the caller's
//     in-memory change stands, because the value it holds is the one on disk.
//  3. The write happens. Its report names every key path that differs between the
//     old file and the new one (AC#2's "which keys did you actually write"), and
//     every key path the file now holds that this process does NOT hold in
//     memory - the hand edits it preserved rather than applied.
//
// AC#3's narrowing of statOwnWrite lives at the end of case 3: the file's new
// stat is adopted only when the bytes just written are exactly what memory holds.
// When the merge carried foreign values into the file, adoption would launder
// them - the reload would never look at content this process never saw - so the
// watch state is left alone and the next CheckAndReload reads the file back and
// makes the D36 verdict on it, which for a loosening is a denial (fail-closed).
func (m *Manager) mergeWrite(ownedKey string, setOn func(base *Config)) error {
	// (1) What does the file hold right now? This is the whole point of the
	// ticket: the answer must be read, not assumed.
	var disk *Config
	_, statErr := os.Stat(m.path)
	switch {
	case statErr == nil:
		loaded, err := readConfigFile(m.path)
		if err != nil {
			return observe.Wrap(observe.ClassConfig, err,
				"config: refused to write config.toml; the file it holds cannot be read back, "+
					"and overwriting an unreadable file is the silent revert ticket 226 exists to stop")
		}
		disk = loaded
	case fileMissing(statErr):
		// No file yet: there is nothing here to clobber, so the snapshot is a
		// legitimate base and writing it creates what first-run did not.
	default:
		return observe.Wrap(observe.ClassConfig, statErr,
			"config: refused to write config.toml; it cannot be stat'ed, so its contents are unknown")
	}

	base := deepCopyConfig(m.cur)
	if base == nil {
		// A Manager built without a loaded config cannot exist in production
		// (NewManager loads or fails); this is the hand-built case, and defaults
		// are what the old path would have serialized anyway.
		base = NewDefaults()
	}
	if disk != nil {
		base = deepCopyConfig(disk)
	}
	setOn(base)

	// (2) The keys this write would actually change in the file. Against the
	// file's own content, never against memory - "what did I write" and "what do
	// I hold" are two different questions and only the first one is this report.
	written := []string{ownedKey}
	if disk != nil {
		written = diffKeyPaths(disk, base)
		if len(written) == 0 {
			slog.Info("config: config.toml already carries this key; nothing written",
				"key", ownedKey)
			return nil
		}
	}

	// (3) Persist. The caller keeps the in-memory half and rolls it back if this
	// fails, so memory and file cannot diverge on a failed write.
	if err := SaveFile(m.path, base); err != nil {
		return err
	}
	// What the file now says that this process does not: the hand edits 甲
	// preserved. Computed against memory AFTER the caller's own change, so the
	// key this write owns is not reported as somebody else's edit.
	diverged := diffKeyPaths(m.cur, base)

	if len(diverged) == 0 {
		slog.Info("config: wrote merged change into config.toml (ticket 226: only the keys listed "+
			"changed value, every other key keeps the value the file already held; the file is re-laid "+
			"in canonical form, so comments and hand-chosen key order are not preserved)",
			"key", ownedKey, "wrote", truncateKeys(written))
		// AC#3: the bytes on disk are exactly the bytes this process holds, so
		// adopting their stat claims only what this write produced.
		m.statOwnWrite()
		return nil
	}
	// AC#3's other half: do NOT adopt. The file changed, and the change contains
	// content this process has never read, so the next poll must see it.
	slog.Warn("config: wrote merged change into config.toml and kept hand edits this process "+
		"does not hold (ticket 226); the file was NOT claimed as our own write, so the next "+
		"reload reads it back and a loosening there is denied until it is confirmed (D36 rule 1). "+
		"Only the keys listed under wrote changed value; the file is re-laid in canonical form, "+
		"so comments and hand-chosen key order are not preserved",
		"key", ownedKey, "wrote", truncateKeys(written), "kept_in_file_not_in_memory", truncateKeys(diverged))
	return nil
}

// diffKeyPaths lists the dotted config.toml key paths whose values differ
// between a and b, sorted, both non-nil. It walks the SAME struct tags the
// serializer uses, so a reported path is a real key in the file, and it follows
// the two exceptions the serializer has: schema_version (never a user setting -
// SaveFile forces it) and non-serialized `toml:"-"` fields, which are skipped;
// Config.Plugins, which is `toml:"-"` yet serialized by hand in canonical form,
// is walked under "plugins".
//
// It is a leaf walk, not a patch builder: slices are compared whole (an added
// allowed_dirs entry is a changed key, not a diff inside one), and map elements
// are compared per key so a plugin or a provider is named as [plugins.<id>] /
// [llm.providers.<name>].
func diffKeyPaths(a, b *Config) []string {
	if a == nil || b == nil {
		// One side absent is the caller's business, not a "no difference" answer.
		// mergeWrite handles the missing file before it gets here; reaching this
		// branch means a new caller forgot to, and saying so beats returning [].
		slog.Error("config: diffKeyPaths called with a nil side; the caller must handle an absent file")
		return []string{"(nil-side)"}
	}
	out := make([]string, 0, 8)
	diffStruct("", reflect.ValueOf(*a), reflect.ValueOf(*b), &out)
	slices.Sort(out)
	return out
}

// diffStruct walks two structs of one type, one key path per field.
func diffStruct(prefix string, a, b reflect.Value, out *[]string) {
	t := a.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Name == "SchemaVersion" {
			// SaveFile forces it to the current version; it is not a value the
			// operator wrote, and reporting it would bury the keys that were.
			continue
		}
		name, transparent := configFieldName(f)
		if !transparent && name == "" {
			continue // never serialized: nothing in the file to differ
		}
		path := prefix
		if name != "" {
			path = joinKey(prefix, name)
		}
		diffValue(path, a.Field(i), b.Field(i), out)
	}
}

// configFieldName resolves the key a field is serialized under. It returns
// ("", true) for a field that carries no key of its own but is still written out
// (PluginsSection.Entries: the map hangs transparently under [plugins]), and
// ("", false) for a field the file never holds.
func configFieldName(f reflect.StructField) (name string, transparent bool) {
	tag, present := f.Tag.Lookup("toml")
	if tag == "-" {
		if f.Name == "Plugins" {
			return "plugins", false
		}
		return "", false
	}
	if !present || tag == "" {
		if f.Type.Kind() == reflect.Map {
			return "", true
		}
		return "", false
	}
	return strings.Split(tag, ",")[0], false
}

// diffValue walks one field or map element: struct and map are descended into,
// everything else is a leaf compared as a value.
func diffValue(path string, a, b reflect.Value, out *[]string) {
	a, b = unwrapInterface(a), unwrapInterface(b)
	if !a.IsValid() || !b.IsValid() {
		if a.IsValid() != b.IsValid() {
			appendKey(out, path)
		}
		return
	}
	if a.Type() != b.Type() {
		appendKey(out, path)
		return
	}
	switch a.Kind() {
	case reflect.Struct:
		diffStruct(path, a, b, out)
	case reflect.Map:
		diffMap(path, a, b, out)
	default:
		if !reflect.DeepEqual(a.Interface(), b.Interface()) {
			appendKey(out, path)
		}
	}
}

// diffMap compares two maps key by key so one plugin or one provider is named by
// its own id. A map whose keys are not strings is compared as a leaf: no
// config.toml path names it anyway.
func diffMap(path string, a, b reflect.Value, out *[]string) {
	if a.Type().Key().Kind() != reflect.String || b.Type().Key().Kind() != reflect.String {
		if !reflect.DeepEqual(a.Interface(), b.Interface()) {
			appendKey(out, path)
		}
		return
	}
	keys := make([]string, 0, a.Len()+b.Len())
	seen := map[string]bool{}
	for _, v := range append(a.MapKeys(), b.MapKeys()...) {
		k := v.String()
		if !seen[k] {
			seen[k] = true
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		child := joinKey(path, k)
		av, bv := a.MapIndex(reflect.ValueOf(k).Convert(a.Type().Key())),
			b.MapIndex(reflect.ValueOf(k).Convert(b.Type().Key()))
		if !av.IsValid() || !bv.IsValid() {
			appendKey(out, child) // present on one side only: the key appeared or went away
			continue
		}
		diffValue(child, av, bv, out)
	}
}

// joinKey appends one segment to a dotted key path.
func joinKey(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}

// appendKey adds one path unless the report is already long enough to read.
func appendKey(out *[]string, path string) {
	*out = append(*out, path)
}

// truncateKeys bounds a report line and says what it bounded.
func truncateKeys(keys []string) []string {
	if len(keys) <= maxReportedKeys {
		return keys
	}
	return append(slices.Clone(keys[:maxReportedKeys]),
		fmt.Sprintf("(+%d more key paths)", len(keys)-maxReportedKeys))
}
