//go:build !windows

package speech

// workingSetBytes has no portable spelling; the residency readings are a Windows
// measurement because the product is a Windows product (D7) and because the SLO
// gate's unit is the private working set (SLO.md section 7).
func workingSetBytes() (int64, bool) { return 0, false }
