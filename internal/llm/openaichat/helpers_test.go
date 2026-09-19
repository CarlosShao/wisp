package openaichat

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Test helpers shared by the adapter tests.

// classOf extracts the D37 class from an error chain (observe.ClassOf).
func classOf(err error) observe.ErrorClass {
	c, _ := observe.ClassOf(err)
	return c
}

// errorsAs finds the *observe.Error in an error chain.
func errorsAs(err error, target **observe.Error) bool {
	return errors.As(err, target)
}

// Test helpers shared by the adapter tests.

// wrapper lets httptestServer carry *httptest.Server's URL without
// embedding an external type in the public test surface.
type wrapper = httptest.Server

// URL returns the server base URL.
func (s *httptestServer) URL() string { return s.wrapper.URL }

// endpointOptionsFor builds adapter options pointing at base (the model id
// matches the fixtures/mockllm vocabulary). The transport is pinned to a
// proxy-less client so machine-level system proxy settings cannot leak into
// adapter tests; tests that exercise proxying override the client.
func endpointOptionsFor(base string, loose, allowMissingUsage bool) llm.EndpointOptions {
	o := endpointOptionsWithClient(base, loose, allowMissingUsage, noProxyClient())
	return o
}

// endpointOptionsWithClient is the explicit-transport variant.
func endpointOptionsWithClient(base string, loose, allowMissingUsage bool, c *http.Client) llm.EndpointOptions {
	var o llm.EndpointOptions
	o.Provider = "mock"
	o.Model = "mock-small"
	o.Protocol = "openai-chat"
	o.BaseURL = base
	o.APIKey = "sk-test"
	o.ContextWindow = 8192
	o.Compat.Loose = loose
	o.Compat.AllowMissingUsage = allowMissingUsage
	o.HTTPClient = c
	return o
}

// noProxyClient is a deterministic transport (no env, no system proxy).
func noProxyClient() *http.Client {
	return &http.Client{Transport: &http.Transport{Proxy: nil}}
}

// envProxyClient honors HTTP(S)_PROXY exactly like the production default.
func envProxyClient() *http.Client {
	return &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment}}
}

// newTestServerHandler starts a plain httptest server around h.
func newTestServerHandler(t *testing.T, h http.HandlerFunc) *httptestServer {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return &httptestServer{srv}
}

// newTLSServer starts an HTTPS server with a self-signed certificate that is
// NOT in any trust store (the corporate TLS-interception simulation).
func newTLSServer(t *testing.T, h http.HandlerFunc) *httptestServer {
	t.Helper()
	srv := httptest.NewTLSServer(h)
	t.Cleanup(srv.Close)
	_ = tls.VersionTLS12 // referenced to keep the crypto/tls import honest
	return &httptestServer{srv}
}
