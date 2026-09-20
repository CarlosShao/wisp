package openaichat

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/llm/golden"
)

// ---------------------------------------------------------------------------
// THE shared golden harness (ticket 11, AC#1).
//
// This file is the openai-chat leg of the single adaptertest table: the two
// adapters added by ticket 11 (anthropic, openai-responses) run the SAME
// CaseTable rows through the SAME runner, only with their own wire-dialect
// fixtures. Chat-specific wire quirks that have no equivalent in the other
// protocols (missing [DONE] vs finish_reason, tool_call fragments without an
// index, proxy/TLS transport detail) stay in adapter_test.go and
// mockllm_integ_test.go; every assertion that expresses protocol-neutral C6
// behavior was MOVED here and now runs for all three adapters.
// ---------------------------------------------------------------------------

// unit describes this adapter to the shared harness.
func unit() adaptertest.Unit {
	return adaptertest.Unit{
		Name:         "openai-chat",
		EndpointPath: "/chat/completions",
		AuthScheme:   "openai-bearer",
		New: func(t *testing.T, base string, o adaptertest.Options) llm.LlmProvider {
			return New(endpointOptionsFor(base, o.Loose, o.AllowMissingUsage))
		},
		Fixture: fixtureFor,
	}
}

// fixtureFor maps a shared scenario to this dialect's golden file. The ten
// ticket-09 canonical names are unchanged (they are this adapter's canonical
// fixtures); only scenarios added by ticket 11 carry the chat- prefix.
func fixtureFor(sc adaptertest.ScenarioID) string {
	switch sc {
	case adaptertest.ScDisconnectTools:
		return "chat-disconnect-tools"
	case adaptertest.ScThinking:
		return "chat-thinking"
	case adaptertest.ScMidStreamError:
		return "chat-mid-stream-error"
	case adaptertest.ScNoTerminal:
		return "missing-finish" // canonical ticket-09 name
	case adaptertest.ScCancel:
		return "long-text" // canonical ticket-09 name (paced)
	default:
		return string(sc) // canonical ticket-09 names
	}
}

// TestSharedGoldenSuite runs every row of the shared table for openai-chat.
func TestSharedGoldenSuite(t *testing.T) {
	adaptertest.RunAll(t, unit())
}

// TestScenarioCoverage is the anti-shortcut guard: this adapter must have a
// real, parseable fixture for EVERY scenario in the shared table (no skipped
// rows), and every chat-prefixed fixture on disk must be claimed by a
// scenario (no orphan fixtures pretending to be coverage).
func TestScenarioCoverage(t *testing.T) {
	table := adaptertest.CaseTable()
	all := adaptertest.AllScenarios()
	if len(table) != len(all) {
		t.Fatalf("CaseTable has %d rows but AllScenarios lists %d", len(table), len(all))
	}
	seen := map[adaptertest.ScenarioID]bool{}
	claimed := map[string]bool{}
	for _, c := range table {
		if seen[c.ID] {
			t.Fatalf("scenario %s appears twice in the table", c.ID)
		}
		seen[c.ID] = true
		name := fixtureFor(c.ID)
		if name == "" {
			t.Fatalf("scenario %s has no fixture", c.ID)
		}
		claimed[name] = true
		rs, err := golden.LoadFile(filepath.Join(adaptertest.GoldenDir(), name+".sse"))
		if err != nil {
			t.Fatalf("scenario %s: fixture %s: %v", c.ID, name, err)
		}
		if len(rs) == 0 {
			t.Fatalf("scenario %s: fixture %s parses to zero responses", c.ID, name)
		}
	}
	for _, sc := range all {
		if !seen[sc] {
			t.Errorf("scenario %s is not in the executed table", sc)
		}
	}
	assertNoOrphans(t, claimed, "chat-")
}

// assertNoOrphans fails on a prefixed fixture no scenario claims (an orphan
// fixture is fake coverage).
func assertNoOrphans(t *testing.T, claimed map[string]bool, prefix string) {
	t.Helper()
	entries, err := os.ReadDir(adaptertest.GoldenDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		n := e.Name()
		if !strings.HasSuffix(n, ".sse") {
			continue
		}
		base := strings.TrimSuffix(n, ".sse")
		if !strings.HasPrefix(base, prefix) {
			continue
		}
		if !claimed[base] {
			t.Errorf("orphan fixture %s: no scenario maps to it", base)
		}
	}
}

// TestMockllmSharedFaultSuite runs the live-server fault suite (mockllm
// subprocess + /__control injections) for this adapter.
func TestMockllmSharedFaultSuite(t *testing.T) {
	adaptertest.RunFaultSuite(t, unit(), adaptertest.FaultSuite{
		Path: "/chat/completions", Route: "chat",
	})
}
