package winsec

// SetSeamForTest is the test-only way to put a resolver into the sealing seam and
// take it back out. Production code has no such door, and that is the point:
// ticket 108's AC#1 removed SetPathResolver's nil fallback because "free the
// seam, then install the forged resolver" is the order that walked ticket 103's
// guard (probe P1b, measured at icacls SID level), and an install you can undo is
// not a one-time install.
//
// A _test.go file is compiled only into this package's test binary, so the hook
// is not a bypassable guard: code outside internal/winsec still cannot release
// the seam in any order, which is what seam_bypass_108_windows_test.go and
// seam_guard_windows_test.go re-measure rather than assert. Inside the package,
// tests need to try several candidate resolvers against one seam in one binary,
// and leaving the accepted fakes installed would poison every later seal - that
// contamination is itself the first thing this ticket's red run measured.
func SetSeamForTest(r C26Resolver) (restore func()) {
	resolverMu.Lock()
	defer resolverMu.Unlock()
	prev, prevLatched := resolver, seamLatched
	resolver, seamLatched = r, r != nil
	return func() {
		resolverMu.Lock()
		defer resolverMu.Unlock()
		resolver, seamLatched = prev, prevLatched
	}
}

// SeamLatchedForTest reports whether the one-way latch has fired, which is what
// makes "install the floor again" a refused operation rather than a reset.
func SeamLatchedForTest() bool {
	resolverMu.RLock()
	defer resolverMu.RUnlock()
	return seamLatched
}

// TreeOwnershipProbeForTest runs the install-time tree-ownership leg of the seam
// guard against a probe pair the caller supplies, without touching the seam.
//
// It exists because that leg's verdict is the one thing in this package that
// depends on how the machine spells its own temp directory, and the CI runner's
// spelling (C:\Users\RUNNER~1\..., an 8.3 short name for a profile directory
// longer than eight characters) is not reproducible on a laptop by living on it.
// The instrument therefore plants the shape in a temporary directory - a long
// name, asked for in its short form, plus a child that does not exist - and
// judges the real pipeline's answer pair (ticket 112 AC#3). Test-only, like
// SetSeamForTest: nothing outside this package's test binary can steer the guard.
func TreeOwnershipProbeForTest(r C26Resolver, parent, child string) string {
	return treeOwnershipFailureForPair(r, parent, child)
}
