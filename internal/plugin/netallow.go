package plugin

import (
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// AC5 — plugin net allowlist (ticket 50, SPEC-07 §4/§6, R5).
//
// The manifest grammar already fails closed at C14 level: a tool declaring
// "net" must ship a non-empty net_allowlist (validator.go). What this file
// owns is the runtime verdict: the target of one outbound call is judged by
// the central C19 assessor with the tool's declared allowlist as the R5
// domain gate.
//
// Verdicts (SPEC-07 §4 "白名单外 → 升 L2", SPEC-06 R5):
//   - host inside the declared allowlist, clean protocol    -> declared floor passes (an L0 tool stays L0)
//   - host outside it                                        -> L2 with rule R5 (auto-escalation, never silent)
//   - deny-listed / unknown protocol                        -> R5's own verdict (Deny / L2), allowlist cannot launder it
//
// The matcher mirrors risk's private hostAllowlisted exactly (exact host or
// subdomain of an entry, case-insensitive, optional :port on the host): two
// allowlist semantics in one system is how "evil-example.com" starts passing
// for "example.com".

// NetDomainAllowed reports whether host matches one of the allowlist entries:
// exact spelling, or a subdomain of an entry (entry "example.com" covers
// "api.example.com"). Matching is case-insensitive; a trailing ":port" on the
// host is ignored. Entries with schemes, paths or wildcards never match —
// C14 already refuses such an allowlist at install; a hand-built one must not
// gain wildcard behaviour here.
func NetDomainAllowed(host string, allowlist []string) bool {
	host = normalizeNetHost(host)
	if host == "" {
		return false
	}
	for _, entry := range allowlist {
		entry = normalizeNetHost(entry)
		if entry == "" {
			continue
		}
		if host == entry || strings.HasSuffix(host, "."+entry) {
			return true
		}
	}
	return false
}

// NetAllowed is the per-tool form: is host inside this tool's declared
// net_allowlist?
func (t ToolDef) NetAllowed(host string) bool {
	return NetDomainAllowed(host, t.NetAllowlist)
}

// normalizeNetHost lowercases, trims, and drops a single trailing :port.
func normalizeNetHost(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if i := strings.LastIndex(s, ":"); i >= 0 && !strings.Contains(s[:i], ":") {
		// host:port spelling; a bare IPv6 keeps its colons (colon in the head
		// means this is not host:port).
		s = s[:i]
	}
	return s
}

// MergeNetAllowlists unions the plugin-declared and host-configured feeds
// (config [net] allowlist may tighten or widen the global gate; the union
// keeps the plugin's own list load-bearing either way). Deduped, lowercased,
// in plugin-then-config order.
func MergeNetAllowlists(feeds ...[]string) []string {
	seen := map[string]bool{}
	var out []string
	for _, feed := range feeds {
		for _, entry := range feed {
			e := normalizeNetHost(entry)
			if e == "" || seen[e] {
				continue
			}
			seen[e] = true
			out = append(out, e)
		}
	}
	return out
}

// LevelForRisk maps a manifest risk spelling to the host's risk.Level.
// An unknown spelling is an error, never a silent default to the safest-
// looking level (C2: the declaration is only the R1 lower bound).
func LevelForRisk(s string) (risk.Level, error) {
	switch s {
	case RiskL0:
		return risk.L0, nil
	case RiskL1:
		return risk.L1, nil
	case RiskL2:
		return risk.L2, nil
	}
	return risk.L0, fmt.Errorf("plugin: risk %q is not one of %s/%s/%s (C2)", s, RiskL0, RiskL1, RiskL2)
}

// AssessPluginNetCall judges one outbound call made by a plugin tool through
// the central C19 assessor, with the tool's declared net_allowlist as the R5
// domain gate. cfgAllowlist carries the host-side [net] allowlist feed (nil
// in tests and for plugins whose gate is purely their own declaration).
// a == nil uses a bare built-in-rule assessor; production hosts pass the one
// wired with C26/C25 dependencies so R2/R3/R4 still fuse on the same call.
func AssessPluginNetCall(def ToolDef, target risk.NetTarget, cfgAllowlist []string, a *risk.RiskAssessor) (risk.Decision, error) {
	level, err := LevelForRisk(def.Risk)
	if err != nil {
		return risk.Decision{}, err
	}
	if a == nil {
		a = risk.NewRiskAssessor()
	}
	target.Allowlist = MergeNetAllowlists(def.NetAllowlist, cfgAllowlist, target.Allowlist)
	dec := a.Assess(def.Name, nil, risk.Facts{
		Declared: level,
		Network:  &target,
	})
	return dec, nil
}
