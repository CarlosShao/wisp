// Command mkv2 (166-v1 bench, runs in the PRE-166 exported tree) writes a real
// schema-v2 database with the binary that existed before ticket 166 landed, and
// seeds one row per table that the v3 upgrade must not disturb.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: mkv2 <data-dir>")
		os.Exit(2)
	}
	dir := os.Args[1]
	s, err := memory.Open(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mkv2: Open:", err)
		os.Exit(1)
	}
	defer s.Close()
	ctx := context.Background()
	if _, err := s.UpsertProfile(ctx, "pref.language", "zh-CN", memory.ProfileManual); err != nil {
		fmt.Fprintln(os.Stderr, "mkv2: UpsertProfile:", err)
		os.Exit(1)
	}
	if _, err := s.AddMemory(ctx, "记得带伞", "带", "伞"); err != nil {
		fmt.Fprintln(os.Stderr, "mkv2: AddMemory:", err)
		os.Exit(1)
	}
	db, err := sql.Open(memory.DriverName, "file:"+dir+"/wisp.db")
	if err != nil {
		fmt.Fprintln(os.Stderr, "mkv2: raw open:", err)
		os.Exit(1)
	}
	defer db.Close()
	var wm string
	if err := db.QueryRow(`SELECT value FROM schema_meta WHERE key='schema_version'`).Scan(&wm); err != nil {
		fmt.Fprintln(os.Stderr, "mkv2: watermark:", err)
		os.Exit(1)
	}
	var n int
	if err := db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`).Scan(&n); err != nil {
		fmt.Fprintln(os.Stderr, "mkv2: table count:", err)
		os.Exit(1)
	}
	fmt.Printf("MKV2 schema_version=%s tables=%d at=%s\n", wm, n, time.Now().Format(time.RFC3339))
	for _, want := range []string{"session", "session_turn"} {
		var has int
		if err := db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?`, want).Scan(&has); err != nil {
			fmt.Fprintln(os.Stderr, "mkv2: hasTable:", err)
			os.Exit(1)
		}
		fmt.Printf("MKV2 hasTable(%s)=%d\n", want, has)
	}
}
