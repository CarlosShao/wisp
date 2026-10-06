package main

// 231-r1 · AC#1 现复现台件（overlay 注入件，物理件在 .scratch/wisp/probes/231/r1/）。
//
// 这一枚只读数、不判等：它走的是票 223 那套真宿主（newReloadRun223 -> runTextTask ->
// cmd/wisp 装的 tick），种一份"schema_version 比本程序高、正文合法且解析得开"的
// config.toml，然后把操作员实际看到的那一条审计行原文打进日志。它不放宽任何断言，
// 因为它压根没有断言 —— 判决靠读原文，不靠颜色。
//
// 种子形状（.scratch/wisp/probes/231/a1/census.md §2.1 末段）：必须"解析得开"，
// 否则 internal/config/loader.go:104 抢在 :120 之前返回 config.toml parse，
// 那一形归 cause=syntax（＝票面 AC#3 要保留的 (a) 形）。
//
// PATH warning (ticket 98)：cmd/wisp 测试二进制链 sherpa-onnx，PATH 不带它会 0xc0000135。

import (
	"strings"
	"testing"
	"time"
)

func TestProbe231R1NewerBuildOperatorLine(t *testing.T) {
	const body = "schema_version = 99\n\n[ball]\nsize = 64\n"
	r := newReloadRun223(t, 40*time.Second)
	r.live(t, func() {
		r.awaitAudit(t, "config: HOT-RELOAD state=armed")
		r.writeOver(t, body)
		trail := r.awaitAudit(t, "config: HOT-RELOAD state=not-applied cause=")
		for _, line := range strings.Split(trail, "\n") {
			if strings.Contains(line, "state=not-applied") {
				t.Logf("AC#1 OPERATOR LINE (verbatim) =\n%s", line)
			}
		}
		t.Logf("AC#1 FULL TRAIL (stderr) =\n%s", trail)
	})
}
