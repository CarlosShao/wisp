package main

// Ticket 223 r2 AC#4 - the permanent routing case for the failure sentences.
//
// WHAT THIS PINS. Which sentence an unreadable config.toml gets is decided by
// which pipeline owns the error, and the split - measured by 223-v1's 17-shape
// overlay, then fixed in internal/config/loader.go (r2 branch), then pinned
// here - is:
//
//   - bytes from which not even a schema_version can be picked   -> 语法错
//     (loader.go branch 1);
//   - a declared CURRENT (or newer) version whose body does not parse ->
//     语法错 (loader.go r2 branch; before r2 these were booked cause=invalid
//     and told "语法没问题", which is the lie this case forbids);
//   - a declared OLDER version with a broken body -> the migration pipeline
//     speaks (migrate_test.go:123 pins that shape at the config layer);
//   - a file that PARSES but the schema rejects -> its own sentence, never
//     语法错. This half is why the table asserts the wrong cause too: a fix
//     that routed everything to 语法错 would pass a one-sided test.
//
// The reading is the production classifier itself: config.Manager.
// CheckAndReload feeds describeReloadFailure exactly what reloadOnce does
// (cmd/wisp/config_reload.go:153-155), only synchronously and without the
// tick, so each line logged here is the line an operator would see.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and
// dies at load (0xc0000135) unless third_party/sherpa-onnx is on PATH.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
)

func TestTicket223R2FailureSentenceRouting(t *testing.T) {
	cases := []struct {
		name        string
		body        string
		wantCause   string // the cause this shape must be booked as
		notCause    string // a cause it must never borrow
		notFragment string // a sentence half that contradicts wantCause ("" = none)
	}{
		// The four shapes 223-v1 measured as swapped (overlay C1/E1/G1/L1),
		// plus the newer-version body-broken shape (overlay J1) and the
		// version-free baseline (overlay A1) that branch 1 already covered.
		{
			"声明当前版_注释以方括号开头_C1", "# [fs] 这不是表头\nschema_version = 2\nbroken [[[\n",
			"cause=syntax", "cause=invalid", "语法没问题",
		},
		{
			"声明当前版_CRLF_E1", "schema_version = 2\r\nbroken [[[\r\n",
			"cause=syntax", "cause=invalid", "语法没问题",
		},
		{
			"声明当前版_无空格_G1", "schema_version=2\nbroken [[[\n",
			"cause=syntax", "cause=invalid", "语法没问题",
		},
		{
			"声明当前版_缩进版本行_L1", "   schema_version = 2\nbroken [[[\n",
			"cause=syntax", "cause=invalid", "语法没问题",
		},
		{
			"声明未来版_正文语法坏_J1", "schema_version = 99\nbroken [[[\n",
			"cause=syntax", "cause=invalid", "语法没问题",
		},
		{
			"读不出版本_A1_基线不动", "this is not toml [[[\n",
			"cause=syntax", "cause=invalid", "语法没问题",
		},
		// The migration pipeline keeps exactly its own shape: a declared
		// version BELOW current. This is migrate_test.go:123's file, seen
		// from the sentence side - r2 must not steal it.
		{
			"声明旧版_坏表头_A2_仍归迁移", "schema_version = 1\n[llm\nbroken ===\n",
			"cause=migration", "cause=syntax", "",
		},
		// And a file that parses is NOT syntax: routing is decided by
		// whether the document parsed, not by which sentence is shorter.
		{
			"解析得开_未知键_不抢语法错", "schema_version = 2\n\nthis_key_does_not_exist = 1\n",
			"cause=unknown-key", "cause=syntax", "",
		},
		// Ticket 231 AC#3: a file that DECLARES A NEWER VERSION *and parses* is its
		// own shape - it never reaches decodeStrict or validate (loader.go:120
		// returns first), so booking it cause=invalid told the operator "内容被校验
		// 拒绝（值不合法或引用解不开）" about a file whose only property is that
		// another build wrote it. J1 above is the OTHER half of the same line and
		// stays untouched: declared-newer with a body that does NOT parse is
		// 语法错 (loader.go:104 short-circuits first). This row is the one that
		// makes the newer-build branch reachable from the permanent table - before
		// it no case in this file planted that shape at all, so deleting the branch
		// reddened nothing (probes/231/a1/census.md §3.3/§6.2 #2).
		{
			"声明未来版_正文解析得开_归更高版本自己那句", "schema_version = 99\n\n[ball]\nsize = 64\n",
			"cause=newer-build", "cause=invalid", "语法没问题",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.toml")
			if err := config.SaveFile(path, config.NewDefaults()); err != nil {
				t.Fatalf("seed canonical config: %v", err)
			}
			mgr, err := config.NewManager(path, nil)
			if err != nil {
				t.Fatalf("NewManager on a canonical config must work: %v", err)
			}
			time.Sleep(30 * time.Millisecond)
			if err := os.WriteFile(path, []byte(tc.body), 0o600); err != nil {
				t.Fatalf("plant: %v", err)
			}
			time.Sleep(30 * time.Millisecond)
			_, err = mgr.CheckAndReload()
			if err == nil {
				t.Fatalf("the planted file must fail the reload, got a clean one")
			}
			line := describeReloadFailure(err)
			t.Logf("ROUTING %q\n  raw err = %v\n  PRODUCTION LINE = %s", tc.name, err, line)
			if !strings.HasPrefix(line, tc.wantCause) {
				t.Errorf("shape %s is booked %q, want it to open with %q - the two sentences swapped again",
					tc.name, line, tc.wantCause)
			}
			if strings.Contains(line, tc.notCause) {
				t.Errorf("shape %s borrows %q, which is a different pipeline's answer:\n%s",
					tc.name, tc.notCause, line)
			}
			if tc.notFragment != "" && strings.Contains(line, tc.notFragment) {
				t.Errorf("shape %s says %q, which contradicts its own cause:\n%s",
					tc.name, tc.notFragment, line)
			}
		})
	}
}
