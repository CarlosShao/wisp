package anthropic

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
)

// The anthropic leg of THE shared golden harness (ticket 11 AC#1): same
// runner, same CaseTable rows, same expectations as openai-chat - only the
// wire-dialect fixture bytes differ.

func unit() adaptertest.Unit {
	return adaptertest.Unit{
		Name:         Protocol,
		Prefix:       "anthropic-",
		EndpointPath: "/messages",
		AuthScheme:   "anthropic-header",
		New: func(t *testing.T, base string, o adaptertest.Options) llm.LlmProvider {
			return New(endpointOptionsFor(base, o.Loose, o.AllowMissingUsage))
		},
		Fixture: fixtureFor,
	}
}

// fixtureFor maps a shared scenario onto this dialect's golden file.
func fixtureFor(sc adaptertest.ScenarioID) string {
	if sc == adaptertest.ScCancel {
		return "anthropic-long-text" // paced 30-chunk fixture (chat's long-text)
	}
	return "anthropic-" + string(sc)
}

// TestSharedGoldenSuite runs every row of the shared table for anthropic.
func TestSharedGoldenSuite(t *testing.T) { adaptertest.RunAll(t, unit()) }

// TestScenarioCoverage proves no row was skipped and no fixture is orphaned.
func TestScenarioCoverage(t *testing.T) { adaptertest.AssertCoverage(t, unit()) }

// TestMockllmSharedFaultSuite runs the live mockllm fault suite.
func TestMockllmSharedFaultSuite(t *testing.T) {
	adaptertest.RunFaultSuite(t, unit(), adaptertest.FaultSuite{
		Path: "/messages", Route: "messages",
	})
}
