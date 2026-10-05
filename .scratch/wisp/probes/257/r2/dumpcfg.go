//go:build ignore

// Probe for ticket 257-r2 §3: print the exact bytes the production first-run
// path writes, so the hand-add patch the receipt teaches can be verified
// against the real section layout instead of an assumed one.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/config"
)

func main() {
	dir, err := os.MkdirTemp("", "wisp257dump")
	if err != nil {
		panic(err)
	}
	path := filepath.Join(dir, "config.toml")
	if err := config.SaveFile(path, config.NewDefaults()); err != nil {
		panic(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("=== %s (%d bytes) ===\n%s\n", path, len(data), string(data))
	fmt.Println("=== section heads ===")
	for _, line := range splitLines(string(data)) {
		if len(line) > 0 && line[0] == '[' {
			fmt.Println(line)
		}
	}
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return out
}
