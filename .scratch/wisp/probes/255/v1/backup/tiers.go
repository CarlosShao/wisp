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
var TierRegistry = map[string]string{
	// Whole-section hot (the hot-apply table in manager.go is derived from
	// exactly these rows).
	"ball":    "hot",
	"session": "hot",
	"audio":   "hot",
	"llm":     "hot",
	"agent":   "hot",
	"privacy": "hot",
	"memory":  "hot",
	"panel":   "hot",
	"cost":    "hot",
	"models":  "hot",
	"observe": "hot",
	"hotkey":  "hot",
	// Locked sections: not a reload tier - any loosening queues an L2 card
	// (planLocked), tightening applies quietly.
	"risk":    "locked",
	"fs":      "locked",
	"net":     "locked",
	"plugins": "locked",
	// [app] is per-key: theme is the one hot key, the rest restart
	// (manager.go planApp; ticket 255 AC#2-ⓑ - registered per key).
	"app.theme":           "hot",
	"app.language":        "restart",
	"app.autostart":       "restart",
	"app.single_instance": "restart",
	// [voice] is per-key: model-pipeline keys reload (apply + event),
	// tuning keys are hot (manager.go planVoice; SPEC-03 sec 3).
	"voice.enabled":              "reload",
	"voice.wake_word.enabled":    "reload",
	"voice.wake_word.keywords":   "reload",
	"voice.asr":                  "reload",
	"voice.asr.provider":         "reload",
	"voice.asr.model":            "reload",
	"voice.tts.provider":         "reload",
	"voice.tts.voice":            "reload",
	"voice.conversation_mode":    "reload",
	"voice.aec":                  "reload",
	"voice.aec.enabled":          "reload",
	"voice.aec.echo_ref":         "reload",
	"voice.realtime":             "reload",
	"voice.realtime.enabled":     "reload",
	"voice.realtime.provider":    "reload",
	"voice.realtime.model":       "reload",
	"voice.realtime.api_key_ref": "reload",
	"voice.realtime.base_url":    "reload",
	"voice.cloud_asr_chain":      "reload",
	"voice.cloud_tts_chain":      "reload",
	"voice.wake_word.thresholds": "hot",
	"voice.wake_word.veto_words": "hot",
	"voice.tts.speed":            "hot",
	"voice.punctuation":          "hot",
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
func TierOf(path string) (tier string, ok bool) {
	t, ok := TierRegistry[path]
	return t, ok
}
