//go:build windows

package main

// Ticket 43, autostart: the Run key half, against the real HKCU registry.
//
// The value name is a parameter, so these cases write, read and remove an entry
// of their own ("Wisp43r1Test") through the EXACT function production runs -
// the owner's real "Wisp" entry is never touched. The OFF default the ticket's
// fresh-config cell asks for is read off the compiled default table, and the
// fresh-boot behaviour (nothing to remove, no action) is the no-key branch of
// the same function.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"golang.org/x/sys/windows/registry"
)

const autostartTestValueName = "Wisp43r1Test"

func autostartValue43(t *testing.T) (string, bool) {
	t.Helper()
	k, err := registry.OpenKey(registry.CURRENT_USER, autostartRunKeyPath, registry.QUERY_VALUE)
	if err != nil {
		return "", false
	}
	defer k.Close()
	v, _, err := k.GetStringValue(autostartTestValueName)
	if err != nil {
		return "", false
	}
	return v, true
}

func autostartCleanup43(t *testing.T) {
	t.Helper()
	t.Cleanup(func() { _ = applyAutostartValue(autostartTestValueName, false) })
}

func Test43AutostartRunKeyWriteRemoveRoundTrip(t *testing.T) {
	autostartCleanup43(t)
	if _, exists := autostartValue43(t); exists {
		t.Fatal("a leftover test entry existed before the round trip started")
	}

	// ON: the value appears, pointing at THIS test binary, quoted.
	got := applyAutostartValue(autostartTestValueName, true)
	if !strings.Contains(got, "已写入") {
		t.Fatalf("apply(true) receipt = %q, want the 已写入 line", got)
	}
	v, exists := autostartValue43(t)
	if !exists {
		t.Fatalf("Run value %q missing after apply(true): %s", autostartTestValueName, got)
	}
	exe, _ := os.Executable()
	if want := `"` + exe + `"`; v != want {
		t.Fatalf("Run value = %q, want the quoted executable %q", v, want)
	}

	// ON again: idempotent, and the receipt says nothing was rewritten.
	got = applyAutostartValue(autostartTestValueName, true)
	if !strings.Contains(got, "已存在且一致") {
		t.Fatalf("second apply(true) receipt = %q, want the 已存在且一致 line", got)
	}

	// OFF: the value is gone.
	got = applyAutostartValue(autostartTestValueName, false)
	if !strings.Contains(got, "已删除") {
		t.Fatalf("apply(false) receipt = %q, want the 已删除 line", got)
	}
	if _, exists := autostartValue43(t); exists {
		t.Fatalf("Run value %q survived apply(false): the toggle did not remove it", autostartTestValueName)
	}

	// OFF again on an empty key: honest no-action, not a fake delete.
	got = applyAutostartValue(autostartTestValueName, false)
	if !strings.Contains(got, "无本程序条目") {
		t.Fatalf("second apply(false) receipt = %q, want the 无本程序条目 line", got)
	}
}

func Test43FreshConfigAutostartDefaultsOff(t *testing.T) {
	autostartCleanup43(t)
	// The schema's own zero value is the ticket's OFF default: a fresh config
	// never opt-ins by itself.
	if config.NewDefaults().App.Autostart {
		t.Fatal("the compiled default table ships autostart=true: the OFF default is broken")
	}
	// A boot with an unreadable config falls back to that same table, and the
	// OFF state on an empty Run key is an honest no-action.
	got := autostartSync(filepath.Join(t.TempDir(), "config.toml"), autostartTestValueName)
	if !strings.Contains(got, "autostart=false") || !strings.Contains(got, "无动作") {
		t.Fatalf("fresh-config boot receipt = %q, want the autostart=false no-action line", got)
	}
	if _, exists := autostartValue43(t); exists {
		t.Fatal("a fresh OFF boot wrote a Run entry")
	}
}
