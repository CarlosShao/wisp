package observe

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// TestNoBareGoFuncInProductionCode enforces the D22/D38b ban on bare
// `go func(` in production code: every goroutine must go through
// Registry.Spawn so it has a name, owner, exit condition and the recover
// boundary. The real AST scan lands with ticket 08; this grep-level check
// keeps the rule from regressing meanwhile.
func TestNoBareGoFuncInProductionCode(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// thisFile = <repo>/internal/observe/nobarego_test.go
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))

	bare := regexp.MustCompile(`\bgo func\(`)
	checked := 0
	for _, dir := range []string{filepath.Join(repoRoot, "internal"), filepath.Join(repoRoot, "cmd")} {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == "testdata" {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			for i, line := range strings.Split(string(data), "\n") {
				if bare.MatchString(line) {
					t.Errorf("%s:%d: bare `go func(` is banned (D22/D38b): use observe.Registry.Spawn", path, i+1)
				}
			}
			checked++
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", dir, err)
		}
	}
	if checked < 10 {
		t.Fatalf("scan checked only %d files; repo layout changed?", checked)
	}
}
