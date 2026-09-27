package v1

import "path/filepath"

func abs(p string) string {
	x, _ := filepath.Abs(p)
	return filepath.Clean(x)
}
