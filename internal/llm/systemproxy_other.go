//go:build !windows

package llm

import (
	"net/http"
	"net/url"
)

// Non-Windows platforms have no WinINET settings to consult: the system proxy
// component is a no-op and ProxyFunc falls through to the HTTP(S)_PROXY env
// rules (http.ProxyFromEnvironment).
func systemProxyFunc() func(*http.Request) (*url.URL, error) { return nil }
