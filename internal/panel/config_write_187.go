package panel

// Ticket 187's writable-field addition: role_chat_thinking, the thinking tier
// of llm.roles.chat, riding the settings route ticket 248 opened.
//
// WHY A THIRD ROSTER AND NOT A ROW IN AN OLDER ONE. Ticket 248's table is
// pinned by a frozen contract test whose `want` list is the seven 248 names,
// exactly; ticket 313's roster is pinned the same way in config_313_test.go
// and by the disjointness assertion between the two. The additions live in
// their own roster, the door asks ALL THREE (knownWritableField's union), and
// every older pin keeps meaning exactly what it meant. Same discipline the
// 313 header records; the model half of ticket 187 needs NO new field at all,
// because role_chat_model (248) and role_chat_pair (313) already write it.
//
// THE TIER IS IN THIS TABLE, NOT IN THE LEG'S HANDS (ticket 187 AC#1/AC#8).
// cmd/wisp's ApplySetting stamps EffectiveRestart on every [llm] write - true
// for the model fields (the provider chain is built once at assembly and
// never rebuilt) and FALSE for the thinking tier, whose reader is the per-task
// re-resolution over the LIVE config (run.go's per-task ResolveRole over
// rt.mgr.Config(), fed by commitSectionWrite's in-memory apply). A receipt
// that told the user to restart a process for a tier the next task already
// picks up would be the wrong-tier lie ticket 248's AC#8 wording rule exists
// to prevent, so this roster carries the measured tier per field and the
// settings door corrects the stamp for these fields after the leg answers.
// The leg's own stamp stays untouched for every field it already owns.
//
// NO APPROVAL VOCABULARY, NO LOCKED FAMILY. The field name carries neither
// (asserted in config_write_187_test.go), the locked-family guard above both
// older rosters is unchanged and runs BEFORE this roster is consulted, and
// the write cannot touch a permission档: the config-side setter owns the
// vocabulary check, nothing here imports the permission or risk packages,
// and the same test reads this file's bytes to keep it that way (AC#3's
// nail).

import "sort"

// FieldRoleChatThinking is llm.roles.chat.thinking_intensity. Value = the
// tier word itself; the enum (and the current model's own declared list)
// is enforced in internal/config before anything is written, the same way
// provider_billing's enum is.
const FieldRoleChatThinking = "role_chat_thinking"

// fields187 is the roster this file owns, sorted at the accessor. Kept apart
// from the 248 and 313 rosters on purpose (see the header).
var fields187 = []string{
	FieldRoleChatThinking,
}

// fields187TakesValue: the one field needs the envelope's plain value half.
var fields187TakesValue = map[string]bool{
	FieldRoleChatThinking: true,
}

// fieldTier187 is the measured effective tier per field in THIS roster, from
// the ticket 187 AC#1 readings (docs/evidence/s1/187-model-switchable.md §2):
// the thinking tier's reader re-resolves per task from the live config, so
// the next task picks a tier write up. The model fields live in the older
// rosters and keep the leg's restart stamp, which is correct for them.
var fieldTier187 = map[string]EffectiveTier{
	FieldRoleChatThinking: EffectiveNextTask,
}

// WritableFields187 returns ticket 187's writable fields, sorted. It is the
// 187 twin of WritableFields / AdditionalWritableFields, and a test reads the
// roster off here rather than copying it.
func WritableFields187() []string {
	out := append([]string(nil), fields187...)
	sort.Strings(out)
	return out
}

// writableField187 asks this file's roster, not a copy of it.
func writableField187(field string) bool {
	for _, f := range fields187 {
		if f == field {
			return true
		}
	}
	return false
}

// writableFieldTier187 reports the measured tier for a field THIS roster owns.
// The bool is false for every name outside the roster, so the door corrects
// only what this table actually carries and never guesses a tier for an
// older field.
func writableFieldTier187(field string) (EffectiveTier, bool) {
	tier, ok := fieldTier187[field]
	return tier, ok
}
