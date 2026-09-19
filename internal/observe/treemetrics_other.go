//go:build !windows

package observe

// Placeholder tree reader for non-Windows builds (ticket 08). The SLO metric
// set is Windows Job Object semantics (C30); on other platforms the package
// must compile (ubuntu test-core) and the sampling logic stays testable via
// injected fake readers. The product sampler never runs here.

// NullTreeReader always fails with ErrTreeUnsupported.
type NullTreeReader struct{}

// ReadTree implements TreeReader.
func (NullTreeReader) ReadTree() (TreeMetrics, error) {
	return TreeMetrics{}, ErrTreeUnsupported
}
