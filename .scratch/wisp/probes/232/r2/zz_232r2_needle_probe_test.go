package main

// 232-r2 read-only probe (mapped into cmd/wisp by an overlay; the shared work
// tree never holds this file). It measures WHICH stream carries each of the
// three restart-tier needles, so AC#2's "stdout only" rewrite is decided on a
// reading and not on a guess. It asserts nothing.

import (
	"strings"
	"testing"
	"time"
)

func TestProbe232R2WhichStreamCarriesEachNeedle(t *testing.T) {
	r := newProbe232Run(t)
	r.live(t, func() {
		if r.rt.cfg.App.Autostart {
			t.Fatal("probe premise moved: the boot snapshot already autostarts")
		}
		mark := r.h.out.String()
		r.plant(t, "[fs]", "[app]\nautostart = true\n\n[fs]")
		why := r.awaitAudit(t, "config: RESTART-PENDING detail=")
		out := r.awaitStdout(t, "本次运行不会生效")
		win := strings.TrimPrefix(r.h.out.String(), mark)
		for _, needle := range []string{"app.autostart", "开机自启", "重启进程后生效"} {
			t.Logf("PROBE needle=%q stdout-full=%v stdout-window=%v audit=%v why+out=%v",
				needle,
				strings.Contains(out, needle),
				strings.Contains(win, needle),
				strings.Contains(why, needle),
				strings.Contains(why+out, needle))
		}
		t.Logf("PROBE stdout-window-since-mark follows:\n%s", win)
	})
}

func newProbe232Run(t *testing.T) *reloadRun223 {
	t.Helper()
	return newReloadRun223(t, 40*time.Second)
}
