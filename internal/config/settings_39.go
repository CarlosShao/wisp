package config

// Ticket 39 AC#1's tier-events leg: SetAppLanguage, the first settings setter
// for a RESTART-tier key ([app].language, tiers.go's per-key TierRestart row).
//
// WHY THIS SETTER IS THE TIER'S FIRST PRODUCTION REACH. Every setter before it
// writes a hot or locked key, so the guarded write path's tier behavior was
// unobservable from the settings route: a save always applied to memory, and
// the receipt's tier word was a hand stamp. This key is the one that proves the
// write path honors the registry (sectionwrite.go's tier-aware branch): the
// save lands in config.toml but NOT in this process's memory - Manager.Config()
// keeps answering the boot value, the receipt says 需要重启 (tierSentence's
// restart sentence), and the next reload books the change in Report.Restart so
// OnRestartPending says it out loud, exactly like a hand edit of the same key
// (manager.go's planApp). The tier judgment is TierOf's ("app.language" ->
// TierRestart), never a rule written here.
//
// NO VOCABULARY OF ITS OWN. language is a free string in the schema (default
// zh-CN): the UI language tag is the frontends' business, and a config-layer
// enum would be a second gate that can lie the way settings_local_audio.go's
// header already ruled for device names. The gate refuses the empty word; the
// guarded path's validate() runs on the base the write would produce.

import "strings"

// keyPathAppLanguage is the one dotted key path SetAppLanguage owns. It is the
// TierRegistry row's own spelling (tiers.go: "app.language": TierRestart).
const keyPathAppLanguage = "app.language"

// SetAppLanguage writes app.language through the guarded key-by-key path every
// settings setter uses. Because the key is registered TierRestart, the write
// lands on disk and the running process keeps its boot value until restart -
// the tier is the registry's, not this method's to choose.
func (m *Manager) SetAppLanguage(language string) ([]string, error) {
	v := strings.TrimSpace(language)
	return m.writeOneKey(keyPathAppLanguage, func(base *Config) bool {
		base.App.Language = v
		return true
	}, nil)
}
