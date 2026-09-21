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
