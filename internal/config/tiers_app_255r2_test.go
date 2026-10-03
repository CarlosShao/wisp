package config

// Ticket 255 AC#3's residual half, assigned to 255-r2 by the orchestrator in the
// ticket's own 翻勾 note: 「范围披露：[app] 新增键不在键级尺射程（固定名单非反射走查）
// ——我裁＝记范围不算欠账，补尺（一行循环）归 255-r2 顺手做」.
//
// WHAT WAS ALREADY MEASURED. tiers_255_test.go's TestPerKeySectionsAreRegisteredPerKey
// names [app]'s four keys from a fixed list (appHot/appRestart). It answers "are
// these four registered" - it CANNOT hear a fifth key: add
// `PortraitMode bool` to AppSection in schema.go and every ruler in this package
// stays silent, which is exactly the shape that made the [risk] four keys lie for a
// whole ticket (unwired.go's header tells that story).
//
// WHAT THIS FILE ADDS. The missing reflection walk: every leaf key AppSection
// declares must carry a tier row in TierRegistry, and every `app.*` row in
// TierRegistry must correspond to a leaf that still exists. Both directions are
// needed - the first is 票 255's "新增一枚哑键必红" capability shape, the second is
// the pin that keeps someone from "quieting" the ruler by deleting rows (the
// 定式 TestRegistryCoversSchemaSectionsGreen already applies to the section rows).
//
// THE MUTATION READING (impl.md §4) plants `portrait_mode` in AppSection without a
// registry row, names it here, then registers it and re-measures green.

import (
	"reflect"
	"strings"
	"testing"
)

// TestEveryAppKeyIsRegisteredInTierRegistry is the AC#3 residual ruler.
func TestEveryAppKeyIsRegisteredInTierRegistry(t *testing.T) {
	typ := reflect.TypeOf(AppSection{})
	leaves := map[string]bool{}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		leaf := strings.Split(f.Tag.Get("toml"), ",")[0]
		if leaf == "" || leaf == "-" {
			// AppSection.Portable carries `toml:"-"` on purpose (schema.go's own
			// comment: it is decided by portable.txt next to the exe, "never by
			// config.toml", and a literal `portable` key is rejected as unknown).
			// It is not a tiered key, so it must NOT gain a registry row either -
			// the reverse walk below checks that too.
			continue
		}
		leaves[leaf] = true
		path := "app." + leaf
		tier, ok := TierRegistry[path]
		if !ok {
			t.Errorf("app key %q (AppSection.%s) has no tier row in TierRegistry: ticket 255 AC#3's residual - a NEW [app] key is invisible to the fixed-name list, so it must be caught by the walk (register it in tiers.go, or delete the key)",
				path, f.Name)
			continue
		}
		switch tier {
		case "hot", "reload", "restart":
		default:
			t.Errorf("app key %q registered with tier %q, which is not a reload tier (hot/reload/restart; a locked section's key never lands here)",
				path, tier)
		}
	}

	// The other direction: a row that names a key the schema no longer declares
	// is a registry that drifted from the truth it claims to describe.
	for path, tier := range TierRegistry {
		if !strings.HasPrefix(path, "app.") {
			continue
		}
		if !leaves[strings.TrimPrefix(path, "app.")] {
			t.Errorf("TierRegistry registers %q (tier %q) but AppSection declares no such key - the row is stale or the key was renamed without touching tiers.go",
				path, tier)
		}
	}
}

// TestAppKeyWalkSeesTheRosterItClaims guards the ruler itself: if the walk ever
// looks at fewer leaves than TierRegistry registers for [app], the ruler has been
// narrowed (a refactor, or a new key smuggled in behind `toml:"-"` to dodge
// registration) and would go green while proving nothing. The count is a floor, not
// an equality, on purpose: the honest way to grow [app] is to add the key AND its
// tier row, and that must not redden anything (票 255 AC#3's 正控 is literally
// "补登记后不响").
func TestAppKeyWalkSeesTheRosterItClaims(t *testing.T) {
	typ := reflect.TypeOf(AppSection{})
	walked := 0
	for i := 0; i < typ.NumField(); i++ {
		leaf := strings.Split(typ.Field(i).Tag.Get("toml"), ",")[0]
		if leaf == "" || leaf == "-" {
			continue
		}
		walked++
	}
	registered := 0
	for path := range TierRegistry {
		if strings.HasPrefix(path, "app.") {
			registered++
		}
	}
	if walked < registered || registered < 4 {
		t.Errorf("[app] walk covers %d leaves, TierRegistry registers %d app keys (want walked>=registered>=4: language/theme/autostart/single_instance): one side of the ruler was narrowed", walked, registered)
	}
}
