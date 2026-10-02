package config

import (
	"reflect"
	"strings"
	"testing"
)

// TestEverySectionHasATierInRegistry is ticket 255 AC#3's ruler: every
// top-level section declared in schema.go must be accounted for in
// TierRegistry (or be one of the four locked sections). Adding a section
// without registering it fails here - that is the "哑键会响" half. The
// positive control is TestRegistryCoversSchemaSectionsGreen, which pins the
// current roster so a refector who deletes rows from TierRegistry also goes
// red, not just someone who adds a section.
func TestEverySectionHasATierInRegistry(t *testing.T) {
	// schema_version is a top-level scalar owned by load/migration (not a
	// reload tier, and not a section); Plugins carries `toml:"-"` (dynamic
	// per-plugin sub-tables absorbed via the raw map in parse.go) but is a
	// locked section and must stay registered.
	perKeySections := map[string]bool{"app": true, "voice": true}
	nonSectionFields := map[string]bool{"SchemaVersion": true}
	c := reflect.TypeOf(Config{})
	for i := 0; i < c.NumField(); i++ {
		f := c.Field(i)
		section := strings.Split(f.Tag.Get("toml"), ",")[0]
		if section == "" || section == "-" {
			if f.Name == "Plugins" {
				if tier, ok := TierRegistry["plugins"]; !ok || tier != "locked" {
					t.Errorf("Plugins (`toml:\"-\"`, dynamic sub-tables) must stay registered as locked in TierRegistry; got tier=%q ok=%v", tier, ok)
				}
				continue
			}
			if !nonSectionFields[f.Name] {
				t.Fatalf("schema section %d (%s) has no toml tag; the registry walk cannot see it", i, f.Name)
			}
			continue
		}
		if nonSectionFields[f.Name] || lockedSections[section] || perKeySections[section] {
			continue
		}
		tier, ok := TierRegistry[section]
		if !ok {
			t.Errorf("section %q has no tier in TierRegistry (ticket 255 AC#3: a new section with no tier is a dumb key)", section)
			continue
		}
		if tier != "hot" && tier != "restart" {
			t.Errorf("section %q registered with tier %q at section level; whole-section rows are hot or restart only (locked goes in lockedSections, per-key gets key rows)", section, tier)
		}
	}
}

// TestRegistryCoversSchemaSectionsGreen pins the current roster: the count of
// registered whole-section rows must equal the schema section count minus the
// four locked sections minus the two per-key sections. Deleting rows from
// TierRegistry goes red here even without adding anything to schema.go.
func TestRegistryCoversSchemaSectionsGreen(t *testing.T) {
	c := reflect.TypeOf(Config{})
	schemaSections := map[string]bool{}
	for i := 0; i < c.NumField(); i++ {
		f := c.Field(i)
		section := strings.Split(f.Tag.Get("toml"), ",")[0]
		if section == "" || section == "-" {
			continue // schema_version scalar and Plugins (`toml:"-"`)
		}
		schemaSections[section] = true
	}
	registered := 0
	for name, tier := range TierRegistry {
		if strings.Contains(name, ".") {
			continue // per-key rows are counted by the other tests
		}
		registered++
		if name == "plugins" {
			// Plugins is `toml:"-"` so it is not in schemaSections; its
			// presence in the schema is pinned in
			// TestEverySectionHasATierInRegistry.
			if tier != "locked" {
				t.Errorf("plugins must carry \"locked\"; got %q", tier)
			}
			continue
		}
		if !schemaSections[name] {
			t.Errorf("TierRegistry registers section %q with tier %q but schema.go has no such section; the registry drifted from the schema", name, tier)
		}
		if lockedSections[name] && tier != "locked" {
			t.Errorf("locked section %q registered with tier %q in TierRegistry; locked sections must carry \"locked\" (planLocked is not a reload tier)", name, tier)
		}
	}
	// schemaSections has 18 entries: 17 sections + the schema_version scalar
	// (it carries a toml tag but is load/migration's field, not a tiered
	// section). Every section except app/voice (per-key) gets exactly one
	// whole-section row, locked ones included; plugins (`toml:"-"`) is a
	// section too but invisible to the tag walk, so +1 restores it.
	// want = 17 - 2 + 1 = 16.
	want := len(schemaSections) - 1 - 2 + 1
	if registered != want {
		t.Errorf("registered whole-section rows = %d, want %d (schema fields %d = 17 sections + schema_version scalar; -2 per-key app/voice; +1 plugins `toml:\"-\"`); a row was added or removed without the schema following", registered, want, len(schemaSections))
	}
}

// TestPerKeySectionsAreRegisteredPerKey is the ticket-255-AC#2-ⓑ half: [app]
// and [voice] split tiers per key, so the registry must carry per-key rows
// and the known split points must say what planApp/planVoice actually do.
func TestPerKeySectionsAreRegisteredPerKey(t *testing.T) {
	appHot := []string{"app.theme"}
	appRestart := []string{"app.language", "app.autostart", "app.single_instance"}
	for _, k := range appHot {
		if tier, ok := TierRegistry[k]; !ok || tier != "hot" {
			t.Errorf("%q missing or not hot (planApp makes theme the one hot key)", k)
		}
	}
	for _, k := range appRestart {
		if tier, ok := TierRegistry[k]; !ok || tier != "restart" {
			t.Errorf("%q missing or not restart (planApp defers these to next start)", k)
		}
	}
	voiceReload, voiceHot := 0, 0
	for name, tier := range TierRegistry {
		switch {
		case strings.HasPrefix(name, "voice.") && tier == "reload":
			voiceReload++
		case strings.HasPrefix(name, "voice.") && tier == "hot":
			voiceHot++
		}
	}
	if voiceReload == 0 || voiceHot == 0 {
		t.Errorf("[voice] per-key split lost: reload rows=%d hot rows=%d (planVoice models both tiers; planApp:349/:355 and planVoice:412/:415 are one section with two tiers - ticket 255 AC#2-ⓑ)", voiceReload, voiceHot)
	}
}

// TestVoicePerKeyRowsMirrorPlanVoiceWalks pins the voice key roster against
// the actual walks planVoice performs: every reload-tier comparison in the
// walk body must have a reload row, every hot-tier one a hot row. Walks are
// found by reflection over the VoiceSection fields: any leaf the plan writes
// back unconditionally (see planVoice's set closure) must appear in the
// registry under one of the two tiers - the two lists together must cover
// the section's leaves exactly.
func TestVoicePerKeyRowsMirrorPlanVoiceWalks(t *testing.T) {
	covered := map[string]bool{}
	for name := range TierRegistry {
		if strings.HasPrefix(name, "voice.") {
			covered[strings.TrimPrefix(name, "voice.")] = true
		}
	}
	walkLeaves(t, reflect.TypeOf(VoiceSection{}), "", covered)
}

func walkLeaves(t *testing.T, typ reflect.Type, prefix string, covered map[string]bool) {
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		leaf := strings.Split(f.Tag.Get("toml"), ",")[0]
		if leaf == "" || leaf == "-" {
			t.Fatalf("VoiceSection field %s has no toml tag", f.Name)
		}
		path := leaf
		if prefix != "" {
			path = prefix + "." + leaf
		}
		ft := f.Type
		if ft.Kind() == reflect.Struct && ft != reflect.TypeOf([3]float64{}) {
			// nested struct: walk deeper, but skip slice-of-struct tables
			if ft.NumField() > 0 {
				walkLeaves(t, ft, path, covered)
				continue
			}
		}
		if !covered[path] {
			t.Errorf("voice key %q is in schema but has no tier row in TierRegistry (planVoice reads it for one of the two tiers; ticket 255 AC#3)", path)
		}
	}
}
