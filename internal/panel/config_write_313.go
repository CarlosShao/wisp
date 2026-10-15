package panel

// Ticket 313's writable-field additions (A877 rulings 1 and 2).
//
// WHY THIS IS A SECOND TABLE AND NOT SIX MORE ROWS IN THE FIRST ONE. Ticket
// 248's table in config_handlers.go is pinned by a frozen contract test
// (config_route_248_test.go: TestAC1FieldTableCarriesNoApprovalVocabulary) whose
// `want` list is the seven ticket-248 names, exactly. Ticket 313 was sent with
// "只许追加 ... ⛔ 改任何既有断言" (313r1a-send §2), so the additions live in their
// own roster and the 248 pin keeps meaning what it meant: that list is the
// ticket-248 surface, not a running total. The write gate asks BOTH rosters
// (knownWritableField), so a name is accepted by the door only when one of the
// two tables carries it, and the locked-family guard above both is unchanged.
//
// WHAT EACH FIELD OWNS.
//
//	voice_cloud_asr_chain / voice_cloud_tts_chain  the cloud fallback chains,
//	  one "provider/model" element each (comma-separated in the value string).
//	  A877 ruling 1 opened these two names; every element passes the same
//	  checkChainElement ruler validateCatalog runs at load.
//	voice_realtime_provider / voice_realtime_model  the realtime head, as two
//	  single-key fields and never one pair field (A877 ruling 1's β shape: the
//	  one-key-per-call discipline is not bent where the per-key shape can work).
//	provider_row_create  A877 ruling 2's explicit row-creation step: provider
//	  in configProvider, an optional model id in configModel. It is the only
//	  field that touches a provider section that does not exist yet, and it is
//	  explicit by name - see internal/config/settings_313.go's header for the
//	  implicit/explicit line and why the older setters still refuse.
//	role_chat_pair  llm.roles.chat's provider and model in one write, because
//	  validateRoles refuses both half states, so from an unset role the two
//	  single-key setters can never reach a legal state.
//
// NO FIELD NAME HERE REACHES A LOCKED FAMILY (risk. / fs. / net. / plugins. /
// privacy. / models. / audio.), and none carries a word of the approval
// vocabulary the frozen grant-boundary instrument keys on - both are asserted in
// config_313_test.go rather than promised in a comment.
//
// TICKET 320's FIRST ADDITION lives in the same roster and follows the same two
// rules the 313 names follow. It is a FIELD, not a method: no new socket name, no
// ConfigStore method, so the 248 pin and its legs stay untouched.
//
//	provider_billing  llm.providers.<name>.billing, the pay-per-token|plan
//	  switch. The enum is enforced in internal/config before the write, the
//	  same way the three older leaf setters refuse their out-of-band values.

import "sort"

const (
	// FieldVoiceCloudASRChain / FieldVoiceCloudTTSChain are the two cloud voice
	// chains (A877 ruling 1). Value = the chain, elements comma-separated.
	FieldVoiceCloudASRChain = "voice_cloud_asr_chain"
	FieldVoiceCloudTTSChain = "voice_cloud_tts_chain"
	// FieldVoiceRealtimeProvider / FieldVoiceRealtimeModel are the realtime
	// head's two leaves (A877 ruling 1's β shape). Value = the name / the id.
	FieldVoiceRealtimeProvider = "voice_realtime_provider"
	FieldVoiceRealtimeModel    = "voice_realtime_model"
	// FieldProviderRowCreate is the explicit provider-row (and optional one-model)
	// creation step (A877 ruling 2). Provider in the selector; model optional.
	FieldProviderRowCreate = "provider_row_create"
	// FieldProviderBilling is llm.providers.<name>.billing (ticket 320 AC#1).
	// Value = plan or pay-per-token, checked against the schema's own enum.
	FieldProviderBilling = "provider_billing"
	// FieldRoleChatPair is llm.roles.chat's provider+model, written together.
	FieldRoleChatPair = "role_chat_pair"
)

// fields313 is the roster this file owns, sorted at the accessor. It is kept
// apart from config_handlers.go's allWritableFields on purpose (see the header).
var fields313 = []string{
	FieldProviderBilling,
	FieldProviderRowCreate,
	FieldRoleChatPair,
	FieldVoiceCloudASRChain,
	FieldVoiceCloudTTSChain,
	FieldVoiceRealtimeModel,
	FieldVoiceRealtimeProvider,
}

// The selector columns of the 313 roster, same shape as the 248 table's three
// maps so the gate can ask one question per column across both tables.
var (
	fields313NeedsProvider = map[string]bool{
		FieldProviderBilling: true, FieldProviderRowCreate: true, FieldRoleChatPair: true,
	}
	fields313NeedsModel = map[string]bool{
		FieldRoleChatPair: true,
	}
	fields313TakesValue = map[string]bool{
		FieldProviderBilling:    true,
		FieldVoiceCloudASRChain: true, FieldVoiceCloudTTSChain: true,
		FieldVoiceRealtimeProvider: true, FieldVoiceRealtimeModel: true,
	}
)

// AdditionalWritableFields returns ticket 313's writable fields, sorted. It is
// the 313 twin of WritableFields, and a test reads the roster off here rather
// than copying it.
func AdditionalWritableFields() []string {
	out := append([]string(nil), fields313...)
	sort.Strings(out)
	return out
}

// additionalWritableField asks this file's roster, not a copy of it.
func additionalWritableField(field string) bool {
	for _, f := range fields313 {
		if f == field {
			return true
		}
	}
	return false
}

// acceptedWritableFields is the union of both rosters, sorted, for the refusal
// text: a page told "可写的只有 ..." must be told about every name the door
// would actually accept, or the sentence sends it chasing the wrong list.
func acceptedWritableFields() []string {
	out := append(WritableFields(), AdditionalWritableFields()...)
	sort.Strings(out)
	return out
}

// needsProviderFor / needsModelFor / takesValueFor answer one selector question
// across both rosters, so the gate's three checks stay one loop-shaped question
// each instead of growing a second branch per column.
func needsProviderFor(field string) bool {
	return fieldNeedsProvider[field] || fields313NeedsProvider[field]
}

func needsModelFor(field string) bool {
	return fieldNeedsModel[field] || fields313NeedsModel[field]
}

func takesValueFor(field string) bool {
	return fieldTakesValue[field] || fields313TakesValue[field]
}
