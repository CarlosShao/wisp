// Command dumpdb (166-v1 bench, runs in the HEAD exported tree) prints the whole
// schema roster of a wisp.db WITHOUT migrating it (mode "pre"), or after the
// production Open has run the chain (mode "post"), or after the baseline binary
// (which cannot reach v3) opened it (mode "base").
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/CarlosShao/wisp/internal/memory"
)

func dumpRaw(dir, label string) {
	db, err := sql.Open(memory.DriverName, "file:"+strings.ReplaceAll(dir, "\\", "/")+"/wisp.db?_pragma=busy_timeout(3000)")
	if err != nil {
		fmt.Fprintln(os.Stderr, label+": raw open:", err)
		os.Exit(1)
	}
	defer db.Close()
	var wm string
	if err := db.QueryRow(`SELECT value FROM schema_meta WHERE key='schema_version'`).Scan(&wm); err != nil {
		fmt.Fprintln(os.Stderr, label+": watermark:", err)
		os.Exit(1)
	}
	fmt.Printf("[%s] schema_version=%s\n", label, wm)
	rows, err := db.Query(`SELECT type, name FROM sqlite_master WHERE name NOT LIKE 'sqlite_%' ORDER BY type DESC, name`)
	if err != nil {
		fmt.Fprintln(os.Stderr, label+": roster:", err)
		os.Exit(1)
	}
	defer rows.Close()
	for rows.Next() {
		var typ, name string
		if err := rows.Scan(&typ, &name); err != nil {
			fmt.Fprintln(os.Stderr, label+": scan:", err)
			os.Exit(1)
		}
		extra := ""
		if typ == "table" {
			crows, err := db.Query(`SELECT name FROM pragma_table_info(?)`, name)
			if err != nil {
				fmt.Fprintln(os.Stderr, label+": table_info:", err)
				os.Exit(1)
			}
			var cols []string
			for crows.Next() {
				var c string
				if err := crows.Scan(&c); err != nil {
					fmt.Fprintln(os.Stderr, label+": col scan:", err)
					os.Exit(1)
				}
				cols = append(cols, c)
			}
			crows.Close()
			extra = " (" + strings.Join(cols, ",") + ")"
		}
		fmt.Printf("  %s %s%s\n", typ, name, extra)
	}
	if err := rows.Err(); err != nil {
		fmt.Fprintln(os.Stderr, label+": rows:", err)
		os.Exit(1)
	}
	var ic string
	if err := db.QueryRow(`PRAGMA integrity_check`).Scan(&ic); err != nil {
		fmt.Fprintln(os.Stderr, label+": integrity:", err)
		os.Exit(1)
	}
	fmt.Printf("[%s] integrity_check=%s\n", label, ic)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: dumpdb <pre|post> <data-dir>")
		os.Exit(2)
	}
	mode, dir := os.Args[1], os.Args[2]
	if mode == "pre" {
		dumpRaw(dir, "PRE-no-migration")
		return
	}
	s, err := memory.Open(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "post: Open:", err)
		os.Exit(1)
	}
	defer s.Close()
	dumpRaw(dir, "POST-after-production-Open")
	ctx := context.Background()
	var content string
	rdb, err := sql.Open(memory.DriverName, "file:"+strings.ReplaceAll(dir, "\\", "/")+"/wisp.db?_pragma=busy_timeout(3000)")
	if err != nil {
		fmt.Fprintln(os.Stderr, "post: raw open:", err)
		os.Exit(1)
	}
	defer rdb.Close()
	if err := rdb.QueryRow(`SELECT content FROM profile WHERE slot='pref.language'`).Scan(&content); err != nil {
		fmt.Fprintln(os.Stderr, "post: seeded profile lost:", err)
		os.Exit(1)
	}
	fmt.Printf("POST seeded-profile=%q\n", content)
	// Use the ticket-166 outlets on a database that came from a pre-166 v2 file.
	if err := s.InsertSession(ctx, memory.Session{ID: "v166-upgraded", Title: "upgraded-db", Workspace: "wk-up"}); err != nil {
		fmt.Fprintln(os.Stderr, "post: InsertSession:", err)
		os.Exit(1)
	}
	got, err := s.ListSessions(ctx, memory.SessionListOptions{Workspace: "wk-up"})
	if err != nil {
		fmt.Fprintln(os.Stderr, "post: ListSessions:", err)
		os.Exit(1)
	}
	fmt.Printf("POST ListSessions(wk-up) rows=%d first=%q\n", len(got), func() string {
		if len(got) == 0 {
			return ""
		}
		return got[0].Title
	}())
}
