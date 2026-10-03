package main

// Ticket 255, the item the orchestrator ruled onto this leg in ledger `A549`
// (「判不动五条我裁两条＝… restartTierKeys 幻影注释归 255-r2 清」).
//
// THE DEFECT. cmd/wisp/config_reload.go's own comment on `restartTierKeys` claims
// the list "can be compared by a reader - and by a test - instead of drifting".
// 票 255's 现量 (§1.1 T2 row) measured that second half as a phantom:
//
//	grep -rn restartTierKeys --include=*_test.go cmd internal   =>  0 hits
//
// so the sentence the restart tier prints to the operator was backed by a comment
// naming a test that did not exist - the same species of claim-without-a-reader this
// ticket was立 for on the hot tier (AC#1). The comment also named a symbol that has
// never existed (`applyApp`; the real one is `planApp`).
//
// WHAT MAKES IT TRUE. Each key the sentence names is driven ONE AT A TIME through
// the real config.Manager (SaveFile -> NewManager -> CheckAndReload), and must land
// in Report.Restart and never in Report.Hot; `app.theme`, the one key planApp
// hot-applies, is the positive control on the other side. The list itself is pinned
// name-for-name, and every named key must be registered "restart" in
// config.TierRegistry - so the sentence's word list (T2) and ⓑ's registry cannot
// drift apart, which is the drift 票 255's §1.1 table exists to count.
//
// No product behaviour changes here: this file only makes an existing claim
// checkable. PATH warning (ticket 98): this package's test binary links sherpa-onnx
// and dies at load unless third_party/sherpa-onnx is on PATH.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
)

// restartTier255Cases are the [app] keys the restart sentence names, with the
// default text the canonical file carries and a replacement that is BOTH different
// and a different length (the poll's fingerprint is mtime+size; equal-length plants
// would make the case depend on mtime granularity alone).
var restartTier255Cases = []struct {
	key      string
	from, to string
}{
	{"app.language", `language = 'zh-CN'`, `language = 'en-US-for-255-r2'`},
	{"app.autostart", "autostart = false", "autostart = true"},
	{"app.single_instance", "single_instance = true", "single_instance = false"},
}

func TestTicket255RestartTierKeysAreBackedByATest(t *testing.T) {
	// (1) name-for-name: the list the operator's sentence prints.
	if got := strings.Join(restartTierKeys, " "); got != "app.language app.autostart app.single_instance" {
		t.Errorf("restartTierKeys = %q, want the three keys planApp defers - the sentence in reportRestartPending prints this list verbatim", got)
	}
	// (2) the same keys must be registered restart in ⓑ's registry (one source).
	for _, c := range restartTier255Cases {
		tier, ok := config.TierOf(c.key)
		if !ok || tier != "restart" {
			t.Errorf("restartTierKeys names %q but config.TierRegistry says tier=%q ok=%v (tiers.go)", c.key, tier, ok)
		}
	}
	// (3) behaviour: one key at a time, through the real Manager.
	for _, c := range restartTier255Cases {
		c := c
		t.Run(c.key, func(t *testing.T) {
			rep := reloadOneAppKey255(t, c.key, c.from, c.to)
			if !contains255(rep.Restart, "app") {
				t.Errorf("%s changed but rep.Restart=%v names no app entry - the restart tier is not what planApp produced", c.key, rep.Restart)
			}
			if contains255(rep.Hot, "app") {
				t.Errorf("%s is a restart key yet rep.Hot=%v carries app - the receipt would print it under 「已立即生效」", c.key, rep.Hot)
			}
		})
	}
	// (4) the positive control on the other side: theme is the one hot [app] key.
	rep := reloadOneAppKey255(t, "app.theme", `theme = 'dark'`, `theme = 'light'`)
	if !contains255(rep.Hot, "app") {
		t.Errorf("app.theme changed but rep.Hot=%v does not carry app - planApp's hot branch moved", rep.Hot)
	}
	if contains255(rep.Restart, "app") {
		t.Errorf("a theme-only edit was booked as restart-tier: rep.Restart=%v", rep.Restart)
	}
	if tier, ok := config.TierOf("app.theme"); !ok || tier != "hot" {
		t.Errorf("TierOf(app.theme) = %q/%v, want hot", tier, ok)
	}
}

// reloadOneAppKey255 writes the canonical config, plants one [app] key, and returns
// the Report the production CheckAndReload answers with. Nothing here mocks the
// tier engine: Manager.plan/planApp decide, this only reads what they booked.
func reloadOneAppKey255(t *testing.T, key, from, to string) *config.Report {
	t.Helper()
	path := filepath.Join(t.TempDir(), configFileName)
	if err := config.SaveFile(path, config.NewDefaults()); err != nil {
		t.Fatalf("SaveFile canonical config: %v", err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read it back: %v", err)
	}
	text := string(raw)
	if !strings.Contains(text, from) {
		t.Fatalf("the canonical file does not carry %q for %s, so this plant would prove nothing:\n%s", from, key, text)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("NewManager on a config the app itself wrote: %v", err)
	}
	time.Sleep(20 * time.Millisecond)
	next := strings.Replace(text, from, to, 1)
	if next == text {
		t.Fatalf("plant of %q changed nothing", from)
	}
	if err := os.WriteFile(path, []byte(next), 0o600); err != nil {
		t.Fatalf("plant write: %v", err)
	}
	rep, err := mgr.CheckAndReload()
	if err != nil {
		t.Fatalf("CheckAndReload after planting %s: %v", key, err)
	}
	if rep == nil {
		t.Fatalf("CheckAndReload answered nil for a planted %s change - the poll never saw it", key)
	}
	return rep
}

func contains255(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}
