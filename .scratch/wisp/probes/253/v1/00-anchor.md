# 253-v1 对抗验收 · 起手锚

- 腿名：`253-v1`（非实现者裁决腿；对面实现腿为 `253-r1`）
- 取量时刻：2026-10-08 18:19 CST（`date` 实跑）
- HEAD：`8dd239f14d8338752882b16a1fdfa94719ad0454`（2026-10-08 18:01:41 +0800，A725 落账）
- 分支：`dev`
- 被验收对象：`cmd/wisp/panel_dispatch_binding_roster_253r1_windows_test.go`（579 行，commit `65f4c968`），两枚顶层用例 `TestTransportDoorBindingMatchesRoster253r1` / `TestBindingRosterBitesItsOwnFixtures253r1`（4 子夹具 good／drifted-bind／empty-bind／untied-init）
- 编排者已复跑：`go test . -run '253r1' -v` rc=0 PASS

## 工作树在飞（只登记，不碰）

`git status --porcelain | head -20` 显示别人的在飞改动（≥20 项，含 `.gitignore`、`.scratch/wisp/probes/152|161|242|268` 若干、`design/**` 若干 D/M）。本腿不碰其中任何一项；本腿写面仅 `.scratch/wisp/probes/253/v1/*.md`。

## 靶子

- 生产码靶：`cmd/wisp/panel_host_windows.go` 的 `installPanelTransport`（约 `:800-810`）——突变只种在这块，种完必还原。
- 测试注入点：`cmd/wisp` 包内测试，注 PATH 后跑；不跑整包。
