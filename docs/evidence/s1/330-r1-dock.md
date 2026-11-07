# 票 330-r1 交件读数（G1 对接钉 · cmd/wisp 面）

- 日期：2026-11-07 11:2x +08；腿＝330-r1 写腿；锚＝HEAD `766bf71f4`。
- 交件枚：`cmd/wisp/plan_present_roster_dock_330r1_test.go`（一枚用例 `Test330DockPinsTheWireNameToTheAssembledRoster`）＋ `.scratch/wisp/probes/330/r1/**`＋本件＋票 330 票尾 Progress log 一条。`internal/**` 零字节（被突变件当场还原，sha 逐位回基线）。

## 判据与读数

- 判据（票面甲形）：`plan-present.sse` 线上那枚工具名 ⇄ 装配出来的真名册必须逐字相等。两半都现读、零枚新名字字面量（AC#2 尺：带引号的 `"plan.present` 字面量＝0 枚）：
  - 线上名读自 `internal/agent/testdata/golden/plan-present.sse:12`（结构化抽 `"function"."name"`）；
  - 真名册读自 `cmd/wisp/run.go:1126` 注册进 `rt.bridge.Registry()` 的那一册（名字最终出自 `internal/tools/plan_present.go:72`）。
- 正向发：PASS，rc=0；线上名 `"plan.present"` ＝ 真名册 22 枚里的一行，`Resident=true`（读数行 `plan_present_roster_dock_330r1_test.go:134`）。
- G1 重放发：夹具（golden `:12`）与镜像（`planPresentName165r3`，`plan_present_doorbell_165r3_test.go:38`）一起改成 `plan.presentX`、真名册不动 ⇒ FAIL，rc=1，红在 `:128`，红句具名 22 枚全名单。165-v2 报的那发绿已变红＝本票唯一判据成立。还原后 sha256 双双回基线（`7d26fa88…`／`6542773d…`），`git status internal/agent` 空。

## 门禁读数表

| 门 | 读数 |
|---|---|
| `gofumpt -l cmd/wisp/` | 列空 rc=0 |
| `go vet ./cmd/wisp` | rc=0 |
| `go build ./cmd/wisp` | rc=0 |
| `go test ./cmd/wisp -run 'Test330' -count=1 -timeout 30m` | PASS rc=0 |
| `sh scripts/d22scan.sh` | clean rc=0 |

- 变异有牙的两发读数、还原核验、中途 d22scan 抓到本腿自己一枚 `⛔`(U+26D4) 字符串字面量后改口的记录，全在 `.scratch/wisp/probes/330/r1/330r1-firings.md`。

## 未做到／留给编排者（宁缺毋造）

- AC#5 的整包四数台面本腿未跑（派单口径＝只点名跑 `Test330`，整包 rc=1 是既有真窗族，⛔ 当判据）；AC#5 那格归编排者裁决。
- AC#7：票 329 未落地（近 15 枚提交零 329 命中），本腿锚＝`766bf71f4`；同面串行由派单保证，AC 框判定权归编排者。
- G2（`Resident` 档位牙）已由 330-r2 在 `internal/agent` 面收口（票尾 330-r2 条），本票不重复。
