package tools

import (
	"sort"
	"strings"
)

// C3 capability face (SPEC-07 §2). The eleven tokens below are the frozen set:
// a token outside it is invalid input, not a capability, and no caller may
// extend the list at runtime (extending C3 is a contract change).
//
// The enforcement rule is asymmetric on purpose: a tool that needs a capability
// it never declared is REJECTED (see hardReject in bridge.go), it is not an
// error return. See Bridge.Execute for why the distinction is load-bearing.
type Capability string

// The frozen C3 set (SPEC-07 §2, verbatim order).
const (
	CapFSRead    Capability = "fs.read"
	CapFSWrite   Capability = "fs.write"
	CapNet       Capability = "net"
	CapClipboard Capability = "clipboard"
	CapShell     Capability = "shell"
	CapSysinfo   Capability = "sysinfo"
	CapScreen    Capability = "screen"
	CapInput     Capability = "input"
	CapWindow    Capability = "window"
	CapNotify    Capability = "notify"
	CapMemory    Capability = "memory"
)

// AllCapabilities is the C3 set in declaration order (11 entries; the count is
// asserted by TestCapabilitySetIsFrozen so a future edit cannot silently
// widen or narrow the contract face).
var AllCapabilities = []Capability{
	CapFSRead, CapFSWrite, CapNet, CapClipboard, CapShell, CapSysinfo,
	CapScreen, CapInput, CapWindow, CapNotify, CapMemory,
}

// Valid reports whether c is one of the eleven frozen tokens.
func (c Capability) Valid() bool {
	for _, a := range AllCapabilities {
		if a == c {
			return true
		}
	}
	return false
}

// String implements fmt.Stringer.
func (c Capability) String() string { return string(c) }

// capSet is a capability set with membership tests; the zero value is empty
// and usable.
type capSet map[Capability]bool

func newCapSet(caps ...Capability) capSet {
	s := make(capSet, len(caps))
	for _, c := range caps {
		s[c] = true
	}
	return s
}

func (s capSet) has(c Capability) bool { return s != nil && s[c] }

// missing returns the members of need that s does not carry (sorted, so the
// rejection text a user sees is stable).
func (s capSet) missing(need []Capability) []Capability {
	var out []Capability
	for _, c := range need {
		if !s.has(c) {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// dedupe returns caps in first-seen order without duplicates.
func dedupe(caps []Capability) []Capability {
	seen := make(capSet, len(caps))
	out := make([]Capability, 0, len(caps))
	for _, c := range caps {
		if seen.has(c) {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	return out
}

// joinCaps renders a capability list for a user-visible reason string.
func joinCaps(caps []Capability) string {
	strs := make([]string, len(caps))
	for i, c := range caps {
		strs[i] = string(c)
	}
	return strings.Join(strs, ", ")
}
