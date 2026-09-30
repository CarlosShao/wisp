package session

// 224-v2 probe (overlay-injected; no tracked file is touched).
//
// Question it answers: the dialect decision says rows and paths "are both slash
// form", but production's canonical form comes out of internal/risk's
// lexCanonical = filepath.Abs + filepath.Clean, i.e. BACKSLASH on Windows, and the
// repository's own wildcard idiom (the frozen guards plant filepath.Join(dir,"*"))
// is backslash too. This probe reads what patternCovers / storeablePattern /
// Record+Covering actually do for those forms, and prints readings only.

import (
	"context"
	"path/filepath"
	"testing"
)

func TestProbe224v2DialectForms(t *testing.T) {
	ctx := context.Background()
	dir := `C:\work\notes`

	cases := []struct {
		name             string
		pattern, concrete string
	}{
		{"A backslash-wildcard row vs a backslash child (Windows idiom)", dir + `\*`, dir + `\todo.txt`},
		{"B slash-wildcard row vs a backslash child (mixed forms)", dir + `/*`, dir + `\todo.txt`},
		{"C slash-wildcard row vs a slash child (what the pins feed)", `/work/notes/*`, `/work/notes/todo.txt`},
		{"D backslash-wildcard row vs a slash child", dir + `\*`, `/work/notes/todo.txt`},
		{"E filepath.Join(dir,\"*\") row vs backslash child", filepath.Join(`C:\work`, `notes`, `*`), dir + `\todo.txt`},
		{"F exact backslash row vs same backslash child", dir + `\todo.txt`, dir + `\todo.txt`},
		{"G exact row vs case-differing child (drive+name)", `C:\Work\new.txt`, `c:\work\new.txt`},
		{"H exact row vs separator-folded child", `C:/work/new.txt`, `C:\work\new.txt`},
	}
	for _, tc := range cases {
		covered := patternCovers(tc.pattern, tc.concrete)
		err := storeablePattern(tc.pattern)
		t.Logf("%-62s patternCovers=%-5v storeablePattern_err=%v", tc.name, covered, err)
	}

	store := openStore(t, filepath.Join(t.TempDir(), "data"))
	l := newTestLedger(t, store)
	for _, p := range []string{dir + `\*`, dir + `/*`, filepath.Join(`C:\work`, `notes`, `*`)} {
		id, err := l.Record(ctx, "fs.write", p)
		gotID, ok := l.Covering(ctx, "fs.write", []string{dir + `\todo.txt`})
		t.Logf("LEDGER row %q accepted(err=%v id=%d) -> Covering(%q)=%v grant_id=%d",
			p, err, id, dir+`\todo.txt`, ok, gotID)
	}

	// Control: the same ledger does honour a slash-form row against a slash child,
	// so anything above is a form mismatch and not "the matcher is dead".
	if _, err := l.Record(ctx, "fs.write", `/ctrl/notes/*`); err != nil {
		t.Fatalf("control Record: %v", err)
	}
	if _, ok := l.Covering(ctx, "fs.write", []string{"/ctrl/notes/todo.txt"}); !ok {
		t.Fatalf("control slash row does not cover a slash child: probe premise broken")
	}
	t.Logf("CONTROL slash-form row covers slash-form child = true")
}
