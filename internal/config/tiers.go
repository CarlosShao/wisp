package config

// TierRegistry is the one machine-readable answer to "which section (or which
// key, inside the per-key sections) lands on which reload tier" (ticket 255
// AC#2-ⓑ). It mirrors, and is the source of, the hot-apply table in
// manager.go's plan(): a section named here with tier "hot" is exactly the
// set plan() walks. Adding a new top-level section to schema.go without
// adding it here fails TestEverySectionHasATierInRegistry, which is what
// keeps this table from going stale the way the [risk] four did (unwired.go
// has the same shape for locked-section keys).
//
// Values:
//
//	"hot"     - plan() hot-applies the whole section when it differs.
//	"restart" - changes are only picked up at the next process start.
//	"locked"  - planLocked(): loosening needs an L2 confirmation (risk/fs/net/
//	            plugins); not a reload tier, registered so the completeness
//	            pin can tell "accounted for" from "forgotten".
//
// Per-key sections ([app], [voice]) are NOT whole-section entries: each key
// gets its own row (see below), because a section can be hot and restart at
// the same time (planApp/planVoice) - registering them whole would be the
// "粒度只到段" lie ticket 255 stands on.
//
// TICKET 233 - THE TIER WORDS ARE THE THREE CONSTANTS, NOT WORDS. The map's
// values are typed through config.Tier (TierHot/TierReload/TierRestart) and
// the "locked" marker string. Before this retag the three Tier* constants had
// zero non-test reference in the repository outside their own declaration
// lines in schema.go - the D36 three-tier enum existed while the table that
// actually tiered every key spelled its values as bare string literals, so a
// typo in one of them ("het"/"relaod") would compile, pass every typed check
// and silently re-tier a key. TierOf still returns a string so its callers
// (cmd/wisp's receipt roster, restartTierKeys' pin) keep their shape.
var TierRegistry = map[string]Tier{
	// Whole-section hot (the hot-apply table in manager.go is derived from
	// exactly these rows).
	"ball":    TierHot,
	"session": TierHot,
	"audio":   TierHot,
	"llm":     TierHot,
	"agent":   TierHot,
	"privacy": TierHot,
	"memory":  TierHot,
	"panel":   TierHot,
	"cost":    TierHot,
	"models":  TierHot,
	"observe": TierHot,
	"hotkey":  TierHot,
	// Locked sections: not a reload tier - any loosening queues an L2 card
	// (planLocked), tightening applies quietly. TierLocked is not one of D36's
	// three reload tiers; it is the marker planLocked's sections carry so the
	// completeness pin can tell "accounted for" from "forgotten".
	"risk":    TierLocked,
	"fs":      TierLocked,
	"net":     TierLocked,
	"plugins": TierLocked,
	// [app] is per-key: theme is the one hot key, the rest restart
	// (manager.go planApp; ticket 255 AC#2-ⓑ - registered per key).
	"app.theme":           TierHot,
	"app.language":        TierRestart,
	"app.autostart":       TierRestart,
	"app.single_instance": TierRestart,
	// [voice] is per-key: model-pipeline keys reload (apply + event),
	// tuning keys are hot (manager.go planVoice; SPEC-03 sec 3).
	"voice.enabled":              TierReload,
	"voice.wake_word.enabled":    TierReload,
	"voice.wake_word.keywords":   TierReload,
	"voice.asr":                  TierReload,
	"voice.asr.provider":         TierReload,
	"voice.asr.model":            TierReload,
	"voice.tts.provider":         TierReload,
	"voice.tts.voice":            TierReload,
	"voice.conversation_mode":    TierReload,
	"voice.aec":                  TierReload,
	"voice.aec.enabled":          TierReload,
	"voice.aec.echo_ref":         TierReload,
	"voice.realtime":             TierReload,
	"voice.realtime.enabled":     TierReload,
	"voice.realtime.provider":    TierReload,
	"voice.realtime.model":       TierReload,
	"voice.realtime.api_key_ref": TierReload,
	"voice.realtime.base_url":    TierReload,
	"voice.cloud_asr_chain":      TierReload,
	"voice.cloud_tts_chain":      TierReload,
	"voice.wake_word.thresholds": TierHot,
	"voice.wake_word.veto_words": TierHot,
	"voice.tts.speed":            TierHot,
	"voice.punctuation":          TierHot,
}

// lockedSections lists the four planLocked sections; the completeness pin
// walks every schema section and requires it to be either registered here or
// present in TierRegistry with a reload tier.
var lockedSections = map[string]bool{
	"risk":    true,
	"fs":      true,
	"net":     true,
	"plugins": true,
}

// TierOf returns the registered tier of a section or key path, and whether it
// is registered at all. The panel write path (cmd/wisp) is the intended
// caller once 255-r2 lands: the receipt's "when does this take effect"
// sentence must be produced from this table, not hand-written per field.
//
// The returned tier is a string view of the Tier value the registry holds
// (ticket 233: the values ARE the three D36 constants plus TierLocked), so
// existing string-based callers keep compiling unchanged. A caller that wants
// the typed value can compare against the Tier constants directly.
func TierOf(path string) (tier string, ok bool) {
	t, ok := TierRegistry[path]
	return string(t), ok
}
