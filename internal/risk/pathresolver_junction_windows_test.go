//go:build windows

package risk

import (
	"testing"
)

// TestPathResolverJunctionWindows is the ticket-08 PLACEHOLDER for the C26
// PathResolver junction cases (SPEC-10 §6 red-team case 8): real
// `mklink /J` junction, 8.3 short name, UNC and \\?\ spellings of an A-tier
// file must all be rejected by the resolver once it lands.
//
// The PathResolver itself is implemented by ticket 18 (fs tool surface) with
// the red-team matrix in ticket 20 - until then this test pins the SLOT in
// the CI matrix (test-windows job) so the gate cannot silently lose the
// case. It asserts the package exists and reports honestly that the real
// assertions arrive with 18/20; the placeholder skip note is the one
// sanctioned by the ticket-08 dispatch, not a CI-level skip (the test-windows
// job itself always runs this test to this assertion).
func TestPathResolverJunctionWindows(t *testing.T) {
	t.Skip("PLACEHOLDER (ticket 08): C26 PathResolver junction/8.3/UNC/\\\\?\\ cases land with tickets 18 (resolver) and 20 (red-team matrix); this slot pins the test-windows job until then")
}
