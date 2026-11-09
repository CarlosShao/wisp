package llm

import (
	"sort"
	"sync"
)

// Ticket 145 AC#2b②'s exported tier vocabulary: the thinking-intensity words
// each registered protocol's adapter ACTUALLY accepts in
// Request.ThinkingIntensity - the fact census
// 145-snapshot-field-census-r1.md §2.14/§4.3 found locked inside three
// package-private tables with no exported face, so a picker could not answer
// "which tiers does THIS adapter map" without importing adapter internals.
//
// THE TRUTH STAYS WITH THE ADAPTERS. The vocabulary is not restated here:
// each adapter derives its registration from its own mapping table
// (openairesponses' effortLevels, anthropic's thinkingBudgets), so the words
// cannot drift from what the adapter maps at request time. openai-chat maps
// no tier at all (its adapter doc says the intensity is intentionally not
// sent) and registers nothing - an empty answer for a registered protocol is
// the measured "this adapter accepts no tier", not a failure.
//
// "off" is deliberately NOT a vocabulary member: it is not a tier, it is the
// request asking for no thinking, and every adapter already handles it
// outside its tier table. Sending the config's four-value enum out verbatim
// as "what everyone supports" is exactly the shape AC#2b forbids.

var (
	thinkingVocabMu sync.RWMutex
	// thinkingVocab maps a registered protocol enum to the tier words its
	// adapter maps, sorted. Adapters self-register from their own package
	// init, the same procedure RegisterProtocol uses.
	thinkingVocab = map[string][]string{}
)

// RegisterThinkingLevels installs the tier words the named protocol's adapter
// accepts in Request.ThinkingIntensity. Call it from an adapter's package
// init, derived from that adapter's own mapping table; registering the same
// protocol twice panics (programming error, caught by tests), matching
// RegisterProtocol's discipline.
func RegisterThinkingLevels(protocol string, levels []string) {
	thinkingVocabMu.Lock()
	defer thinkingVocabMu.Unlock()
	if _, dup := thinkingVocab[protocol]; dup {
		panic("llm: thinking levels for protocol " + protocol + " registered twice")
	}
	sorted := append([]string(nil), levels...)
	sort.Strings(sorted)
	thinkingVocab[protocol] = sorted
}

// ThinkingLevels returns the tier words the named protocol's adapter actually
// accepts, sorted, as a copy the caller cannot mutate the registry through.
// An empty result means the adapter maps no tier (openai-chat) or the
// protocol is not registered - in both shapes the honest answer to the picker
// is "offer no tier choices", never the config's own four-value enum.
func ThinkingLevels(protocol string) []string {
	thinkingVocabMu.RLock()
	defer thinkingVocabMu.RUnlock()
	return append([]string(nil), thinkingVocab[protocol]...)
}
