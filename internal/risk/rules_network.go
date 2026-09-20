package risk

import (
	"fmt"
	"net/netip"
	"strings"
)

// MaxURLLength is the frozen outbound URL cap (SPEC-06 §6.3): longer URLs can
// smuggle file content into the query string.
const MaxURLLength = 2048

// deniedProtocols is the frozen protocol deny list for network targets
// (SPEC-06 §6.5): explicit denies, verdict Deny. file:/mailto:/ms-settings:
// belong to the separate file.open allowlist, not to R5 network targets.
var deniedProtocols = map[string]bool{
	"javascript": true,
	"vbscript":   true,
	"mshta":      true,
	"shell":      true,
	"search-ms":  true,
	"tel":        true,
}

// webProtocols is the only protocol set a network target may speak
// (SPEC-06 §6.3: web.open 仅 http/https).
var webProtocols = map[string]bool{
	"http":  true,
	"https": true,
}

// ruleNetwork implements R5. Verdicts, highest severity wins:
//   - deny-listed protocol -> Deny
//   - host in a private/loopback/link-local range (SSRF defense,
//     incl. 169.254.169.254) -> L2
//   - protocol outside http/https, or unknown -> L2
//   - host outside the configured domain allowlist -> L2 (no allowlist
//     configured = domain gate off; the [net] config owns that choice)
//   - URL longer than MaxURLLength -> L2
func ruleNetwork(ctx *assessCtx) *contribution {
	n := ctx.facts.Network
	if n == nil {
		return nil
	}
	scheme := strings.ToLower(strings.TrimSpace(n.Protocol))
	if scheme == "" && n.URL != "" {
		if i := strings.Index(n.URL, "://"); i > 0 {
			scheme = strings.ToLower(n.URL[:i])
		}
	}
	host := strings.Trim(strings.TrimSpace(n.Host), "[]")

	type verdict struct {
		level  Level
		reason string
	}
	var candidates []verdict
	add := func(l Level, format string, args ...any) {
		candidates = append(candidates, verdict{l, fmt.Sprintf(format, args...)})
	}

	if deniedProtocols[scheme] {
		add(Deny, "R5: 协议被明确拒绝: %s", scheme)
	}
	if isPrivateHost(host) {
		add(L2, "R5: 私有网段目标（SSRF 防御）: %s", host)
	}
	switch {
	case scheme == "":
		add(L2, "R5: 协议未知")
	case !webProtocols[scheme]:
		add(L2, "R5: 协议不在 web 允许列表: %s", scheme)
	}
	if host != "" && len(n.Allowlist) > 0 && !hostAllowlisted(host, n.Allowlist) {
		add(L2, "R5: 域名在白名单之外: %s", host)
	}
	if len(n.URL) > MaxURLLength {
		add(L2, "R5: URL 超长（>%d 字符）", MaxURLLength)
	}
	if len(candidates) == 0 {
		return nil
	}
	worst := candidates[0]
	for _, c := range candidates[1:] {
		if c.level > worst.level {
			worst = c
		}
	}
	return &contribution{rules: []RuleID{R5}, level: worst.level, reason: worst.reason}
}

// isPrivateHost reports whether the host is an IP literal inside the SSRF
// ranges (SPEC-06 §6.3): 10/8, 172.16/12, 192.168/16, 127/8, 169.254/16,
// ::1, fc00::/7, fe80::/10, plus unspecified addresses. DNS names are NOT
// resolved here — pure logic only; resolution-time checks are the net layer's
// job (documented integration boundary).
func isPrivateHost(host string) bool {
	if host == "" {
		return false
	}
	// Strip a trailing :port for plain host:port spellings; bare IPv6 keeps
	// its colons (more than one colon means it is not host:port).
	if i := strings.LastIndex(host, ":"); i >= 0 && strings.Count(host, ":") == 1 {
		host = host[:i]
	}
	ip, err := netip.ParseAddr(host)
	if err != nil {
		return false
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
		ip.IsUnspecified() || ip.IsLinkLocalMulticast()
}

// hostAllowlisted matches the host exactly or as a subdomain of an allowlist
// entry (entry "example.com" covers "api.example.com"), case-insensitively.
func hostAllowlisted(host string, allowlist []string) bool {
	host = strings.ToLower(host)
	for _, entry := range allowlist {
		entry = strings.ToLower(strings.TrimSpace(entry))
		if entry == "" {
			continue
		}
		if host == entry || strings.HasSuffix(host, "."+entry) {
			return true
		}
	}
	return false
}
