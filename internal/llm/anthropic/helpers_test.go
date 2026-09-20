package anthropic

import (
	"github.com/CarlosShao/wisp/internal/llm"
	"net/http"
)

// endpointOptionsFor builds adapter options pointing at base. The transport is
// pinned to a proxy-less client so machine-level proxy settings cannot leak
// into adapter tests (same discipline as the openai-chat tests). Fake key
// only: no real credential is ever read by this package's tests.
func endpointOptionsFor(base string, loose, allowMissingUsage bool) llm.EndpointOptions {
	var o llm.EndpointOptions
	o.Provider = "mock"
	o.Model = "mock-small"
	o.Protocol = Protocol
	o.BaseURL = base
	o.APIKey = "sk-ant-test-not-a-real-key"
	o.ContextWindow = 8192
	o.Compat.Loose = loose
	o.Compat.AllowMissingUsage = allowMissingUsage
	o.HTTPClient = &http.Client{Transport: &http.Transport{Proxy: nil}}
	return o
}
