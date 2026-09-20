package openairesponses

import (
	"net/http"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
)

// The openai-responses leg of THE shared golden harness (ticket 11 AC#1).

func unit() adaptertest.Unit {
	return adaptertest.Unit{
		Name:         Protocol,
		Prefix:       "responses-",
		EndpointPath: "/responses",
		AuthScheme:   "openai-bearer",
		New: func(t *testing.T, base string, o adaptertest.Options) llm.LlmProvider {
			return New(endpointOptionsFor(base, o.Loose, o.AllowMissingUsage))
		},
		Fixture: fixtureFor,
	}
}

// fixtureFor maps a shared scenario onto this dialect's golden file.
func fixtureFor(sc adaptertest.ScenarioID) string {
	if sc == adaptertest.ScCancel {
		return "responses-long-text" // paced 30-chunk fixture
	}
	return "responses-" + string(sc)
}

// TestSharedGoldenSuite runs every row of the shared table for
// openai-responses.
func TestSharedGoldenSuite(t *testing.T) { adaptertest.RunAll(t, unit()) }

// TestScenarioCoverage proves no row was skipped and no fixture is orphaned.
func TestScenarioCoverage(t *testing.T) { adaptertest.AssertCoverage(t, unit()) }

// TestMockllmSharedFaultSuite runs the live mockllm fault suite.
func TestMockllmSharedFaultSuite(t *testing.T) {
	adaptertest.RunFaultSuite(t, unit(), adaptertest.FaultSuite{
		Path: "/responses", Route: "responses",
	})
}

// endpointOptionsFor builds adapter options pointing at base with a fake key
// and a proxy-less transport (no real credential is ever read here).
func endpointOptionsFor(base string, loose, allowMissingUsage bool) llm.EndpointOptions {
	var o llm.EndpointOptions
	o.Provider = "mock"
	o.Model = "mock-small"
	o.Protocol = Protocol
	o.BaseURL = base
	o.APIKey = "sk-test-not-a-real-key"
	o.ContextWindow = 8192
	o.Compat.Loose = loose
	o.Compat.AllowMissingUsage = allowMissingUsage
	o.HTTPClient = &http.Client{Transport: &http.Transport{Proxy: nil}}
	return o
}
