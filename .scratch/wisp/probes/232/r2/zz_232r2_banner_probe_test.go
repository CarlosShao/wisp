package main

// 232-r2 read-only probe #2 (in cmd/wisp only through an overlay; the shared
// work tree never holds this file). AC#4 asks WHICH mechanism makes a
// start-up banner satisfy the restart-tier needles. This measures it directly:
// does the operator stream ALREADY carry the words before the edit is planted?
// It asserts nothing and waits on nothing but the audit line the plant causes.

import (
	"strings"
	"testing"
	"time"
)

func TestProbe232R2BannerAvailableBeforePlant(t *testing.T) {
	r := newReloadRun223(t, 40*time.Second)
	r.live(t, func() {
		if r.rt.cfg.App.Autostart {
			t.Fatal("probe premise moved: the boot snapshot already autostarts")
		}
		mark := r.h.out.String()
		t.Logf("PROBE2 BEFORE plant: stdout contains %q", "本次运行不会生效")
		for _, n := range []string{"本次运行不会生效", "app.autostart", "开机自启", "原因：", "交给平台层", "没有被丢掉", "重启进程后生效"} {
			t.Logf("PROBE2 needle=%q present-in-mark(pre-plant)=%v", n, strings.Contains(mark, n))
		}
		r.plant(t, "[fs]", "[app]\nautostart = true\n\n[fs]")
		r.awaitAudit(t, "config: RESTART-PENDING detail=")
		win := strings.TrimPrefix(r.h.out.String(), mark)
		for _, n := range []string{"本次运行不会生效", "app.autostart", "开机自启", "原因：", "交给平台层", "没有被丢掉"} {
			t.Logf("PROBE2 needle=%q present-in-window(post-plant)=%v", n, strings.Contains(win, n))
		}
	})
}

var _ = time.Second
