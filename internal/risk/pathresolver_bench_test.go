package risk

import (
	"os"
	"path/filepath"
	"testing"
)

// C26 per-call budget (ticket 18 AC#6): <=1ms per Resolve call on a warm handle
// cache, which is what makes a 50-call-per-task envelope affordable.

func seedResolveTarget(tb testing.TB) string {
	tb.Helper()
	dir := tb.TempDir()
	target := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(target, []byte("wisp c26 bench target"), 0o600); err != nil {
		tb.Fatalf("seed bench target: %v", err)
	}
	return target
}

// warmResolve touches the target enough times that the OS path/handle lookup
// stops paying cold-cache cost, so the measurement reflects the resolver
// pipeline rather than first-touch filesystem work.
func warmResolve(tb testing.TB, target string) {
	tb.Helper()
	for i := 0; i < 64; i++ {
		if _, err := Resolve(target, nil); err != nil {
			tb.Fatalf("warm Resolve(%s): %v", target, err)
		}
	}
}

func BenchmarkResolveWarm(b *testing.B) {
	target := seedResolveTarget(b)
	warmResolve(b, target)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Resolve(target, nil); err != nil {
			b.Fatalf("Resolve(%s): %v", target, err)
		}
	}
}
