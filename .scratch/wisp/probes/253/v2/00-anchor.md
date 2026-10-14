# 253-v2 对抗验收 · 起手锚＋读数索引

- 腿名：`253-v2`（非实现者验收腿；被裁＝`253-r5` 的 `internal/panel/inbound_roster_253_test.go`）。
- 工作窗口：2026-10-14 13:56–14:18 +08（`date` 现跑，两发见 logs）。
- 起手 HEAD `581dd252`（13:58:59）→ 收笔 HEAD `c2b00de4`（14:03:44，他腿/编排者 commit，与 `internal/panel` 无关）。
- 起手／终态 `internal/panel` porcelain 均空；三枚靶件 sha256（起手＝终态逐字节）：
  - `bridge.go` `fc23919c034d2f0814a36cc0b4bc810eabb9085b94bdccde1aa27efefe736c37`
  - `composer_dispatch.go` `d001fb22ac0c8407f6a90022d302f91b2ed61da415b57cdb1b53eae317b248ad`
  - `inbound_roster_253_test.go` `646125ce4b234717a3d37021cd98f98242447fdc03076faab7fd327fff105e9c`（834 行）
- 起手起跑前置：整腿 `go test ./internal/panel/ -count=1 -v`＝**141 PASS／5 FAIL／0 SKIP**（5 枚逐名见 evidence §7）。
- 窗口内目击他腿：`253-r6` 在 `cmd/wisp` 落 `panel_unbound_253r6_windows_test.go` 并改 `panel_control_windows.go`（本腿零读零碰）。
- 突变纪律：真树种形全部走「`git cat-file blob` 仓外备份 → 种 → 跑 → blob 还原 → sha256 比对」；overlay 只用于**编译面**（测试副本/新文件/守卫行为），因为本尺在**运行时**用 `os.ReadFile` 读真树，overlay 对运行时读取不可见（本腿实测：overlay 把 bridge.go 换成带第 7 名的副本，主测仍绿并自述 `6 names declared`）。
- 本目录：`denoms253v2.go`＋`probe_253v2_test.go`（我自写的两支只读仪器）＋ 13 枚原始 log ＋ 2 枚 diff。
