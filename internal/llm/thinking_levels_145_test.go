package llm_test

// Ticket 145 AC#2b②'s value teeth for the exported tier vocabulary. The
// registration sites derive each adapter's words from ITS OWN mapping table
// (openairesponses' effortLevels, anthropic's thinkingBudgets), so the
// assertions here are hardcoded from those tables - a table that silently
// loses or gains a tier reds here, and a second copy of the vocabulary that
// drifts from the request-time table reds here too.

import (
	"reflect"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
	_ "github.com/CarlosShao/wisp/internal/llm/anthropic"
	_ "github.com/CarlosShao/wisp/internal/llm/openaichat"
	_ "github.com/CarlosShao/wisp/internal/llm/openairesponses"
)

func Test145ThinkingLevelsComeFromTheAdaptersOwnTables(t *testing.T) {
	// The three words each tiered adapter maps, sorted - "off" is not a
	// member anywhere: it is the no-thinking request, not a tier.
	for _, c := range []struct{ protocol string }{
		{"openai-responses"}, {"anthropic"},
	} {
		got := llm.ThinkingLevels(c.protocol)
		want := []string{"high", "low", "medium"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("ThinkingLevels(%q) = %v, want %v (the adapter's own tier words, sorted)", c.protocol, got, want)
		}
	}
}

func Test145OpenaiChatRegistersKnownButZeroTiers(t *testing.T) {
	// openaichat's own doc: the intensity is intentionally not sent, so its
	// vocabulary is empty. The empty answer must not be distinguishable from
	// a failure to register - both come back empty, and both mean "offer no
	// tier choices" to the picker.
	if got := llm.ThinkingLevels("openai-chat"); len(got) != 0 {
		t.Errorf("ThinkingLevels(openai-chat) = %v, want empty (this adapter maps no tier)", got)
	}
	if got := llm.ThinkingLevels("no-such-protocol"); len(got) != 0 {
		t.Errorf("ThinkingLevels(unregistered) = %v, want empty", got)
	}
}

func Test145ThinkingLevelsReturnIsACopy(t *testing.T) {
	got := llm.ThinkingLevels("openai-responses")
	if len(got) == 0 {
		t.Fatal("expected the registered vocabulary, got empty")
	}
	got[0] = "mutated"
	if again := llm.ThinkingLevels("openai-responses"); again[0] == "mutated" {
		t.Errorf("the caller mutated the registry through the returned slice: %v", again)
	}
}

func Test145ThinkingLevelsDoubleRegistrationPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("registering the same protocol's vocabulary twice must panic (programming error), matching RegisterProtocol's discipline")
		}
	}()
	llm.RegisterThinkingLevels("openai-responses", []string{"low"})
}
