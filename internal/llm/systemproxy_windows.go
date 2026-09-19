//go:build windows

package llm

import (
	"net/http"
	"net/url"
	"strings"

	winreg "golang.org/x/sys/windows/registry"
)

// Windows system proxy support (D42#4 / SPEC-05 sec 3.4: "respect the system
// proxy, WinHTTP default"). The WinINET per-user settings under
// HKCU\Internet Settings are what WinHTTP-configuration default resolves to
// on user sessions; we read them live on every request (cheap registry read,
// honours changes without restart) and parse the WinINET server string
// grammar:
//
//	ProxyEnable = 1 && ProxyServer =
//	  "host:port"                          (single proxy for all protocols)
//	  "http=host:port;https=host:port"     (per-protocol overrides)
//	  "<local>" inside ProxyOverride        bypass loopback/intranet hosts
func systemProxyFunc() func(*http.Request) (*url.URL, error) {
	return func(req *http.Request) (*url.URL, error) {
		if req == nil || req.URL == nil {
			return nil, nil
		}
		host := req.URL.Hostname()
		if host == "" {
			return nil, nil
		}
		enable, server, override := readWinINETProxy()
		if !enable || server == "" {
			return nil, nil
		}
		if overrideBypass(override, host) {
			return nil, nil
		}
		u := parseWinINETServer(server, req.URL.Scheme)
		if u == nil {
			return nil, nil
		}
		return u, nil
	}
}

// readWinINETProxy fetches (ProxyEnable, ProxyServer, ProxyOverride).
// Missing key/values read as disabled; a failed registry open (service
// session, stripped ACL) means "no system proxy" - the env fallback in
// ProxyFunc still applies.
func readWinINETProxy() (enabled bool, server, override string) {
	k, err := winreg.OpenKey(winreg.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Internet Settings`, winreg.QUERY_VALUE)
	if err != nil {
		return false, "", ""
	}
	defer k.Close()
	if v, _, err := k.GetIntegerValue("ProxyEnable"); err == nil {
		enabled = v != 0
	}
	if v, _, err := k.GetStringValue("ProxyServer"); err == nil {
		server = v
	}
	if v, _, err := k.GetStringValue("ProxyOverride"); err == nil {
		override = v
	}
	return enabled, server, override
}

// parseWinINETServer handles both the flat "host:port" form and the
// per-protocol "scheme=host:port;..." form. Returns nil when the string
// carries no usable proxy for the wanted scheme.
func parseWinINETServer(server, scheme string) *url.URL {
	server = strings.TrimSpace(server)
	if server == "" {
		return nil
	}
	if !strings.Contains(server, "=") && !strings.Contains(server, ";") {
		return proxyURL("http", server)
	}
	var chosen string
	for _, part := range strings.Split(server, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		scheme2, addr, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		scheme2 = strings.ToLower(strings.TrimSpace(scheme2))
		if scheme2 == strings.ToLower(scheme) || (scheme2 == "https" && scheme == "wss") {
			chosen = strings.TrimSpace(addr)
			break
		}
	}
	if chosen == "" {
		return nil
	}
	return proxyURL(scheme, chosen)
}

// proxyURL builds the proxy URL. WinINET addresses without a scheme are HTTP
// proxies (CONNECT for https targets).
func proxyURL(scheme, addr string) *url.URL {
	if !strings.Contains(addr, "://") {
		addr = "http://" + addr
	}
	u, err := url.Parse(addr)
	if err != nil {
		return nil
	}
	return u
}

// overrideBypass mirrors the WinINET ProxyOverride bypass grammar: the
// "<local>" token bypasses plain-hostname (intranet) targets, entries match
// by host suffix or full host; a leading/mildcard "*" matches everything.
func overrideBypass(override, host string) bool {
	for _, ent := range strings.Split(override, ";") {
		ent = strings.ToLower(strings.TrimSpace(ent))
		if ent == "" {
			continue
		}
		if ent == "<local>" {
			if !strings.Contains(host, ".") {
				return true
			}
			continue
		}
		ent = strings.TrimPrefix(ent, ".")
		if ent == "*" {
			return true
		}
		h := strings.ToLower(host)
		if h == ent || strings.HasSuffix(h, "."+ent) {
			return true
		}
	}
	return false
}
