# 票 174 · 验收腿 `174-v3` 锚与门禁（只读；判语交回编排者）

- 腿＝`174-v3`（面＝只读）：裁 **AC#3 三条禁区**，按 `174-r4` 交付件 `.scratch/wisp/probes/174/r4/evidence.md`（§2 判据表／§8「AC#3 三半的状态」；它自己写「我不声称 AC#3 结题」）。
- 窗口：起手 HEAD `9065c659`（15:26）→ 交卷 HEAD `037bfb7f`（15:29；他腿提交，与我无重面）。判语相关文件：`internal/tools/task.go` md5 `ded82068…`（＝HEAD blob；`2549a18a→HEAD` 只差 `taskCancel.Execute` 一处，pointerNotice 区零差）、`internal/agent/spill.go` md5 `f31147b4…`＝r4 起手 pristine、r4 两枚判据文件今日 blob `36705837…`／`83dc4955…`（后者含 `16901acb` 的一次签名跟随改动：`:196 SetWorkspaceRoot(res.Canonical, res)`）。
- 台件：`logs/01-zone-i.txt`｜`logs/02-zone-ii.txt`｜`logs/03-zone-iii-and-gates.txt`（`04-gates.txt` 为同内容留盘未提交，清单见下）。`mut/`（pristine 两枚＋三发突变）与 `overlay-*.json` 五枚留盘未提交。

## 门禁（全部现跑，锚 `037bfb7f`）

| 门 | 命令 | 读数 |
|---|---|---|
| 定向 -run（本案主选器） | `go test -count=1 -run '174r4' ./internal/tools/ ./internal/agent/` | **rc=0**（15:33:4x 与 15:35:4x 两跑） |
| 整包 | `go test -count=1 ./internal/tools/ ./internal/agent/ ./internal/config/` | **rc=0** 安静窗（tools ok 14.325s／agent ok 2.304s／config ok 1.115s） |
| d22scan | `sh scripts/d22scan.sh`（GOTMPDIR 未设，复验 `76b9e214`） | **clean rc=0**（`ban #7 internal/tools/ examined 23`、`ban #8 internal/ 550`、`cmd/ 127`） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | **rc=0 VERDICT GREEN**：`tracked paths=12778 over-budget=60 covered by roster=60 not in roster=0`（**60=60**） |

## 窗口逐报（他腿在飞；两处，均不归我）

1. `internal/tools/paths.go` 在 **15:33:50 现读为 ` M`**（252-r4 的在树突变窗）：同窗我跑整包 ⇒ 四枚 FAIL 全是它的（`TestTicket252R1AlignmentAddsNoAuthorization`／`TestTicket107bProbeCLinkInsideAllowedRootStaysOutside`／`TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses`／`TestTicket252R2BothContainmentsReachable`，红句逐字在 logs/04 段的 15:34 块里），**其中两条明写 `InAllowlist` 只剩 1 次包含**＝那时 C26 行为被改过。**安静窗重跑（15:35:38–15:35:51，porcelain 只剩 253-v3 的两枚 `cmd/wisp` 文件、`internal/` 空）全部 rc=0**，四枚 FAIL 消失 ⇒ 归因坐实、不记在本腿头上。15:35:09 复读 `paths.go` 三处包含（`:208/:221/:227`）齐、`git diff --stat -- internal/tools/paths.go` 空。
2. `cmd/wisp/panel_control_windows.go`（` M`）＋`cmd/wisp/panel_unbound_253v3_windows_test.go`（未跟踪）＝253-v3 的在飞窗；本腿未读未跑未归因该包。

## 不动产码自证

三条突变全走 `go test -overlay`（五枚 overlay json 在台件目录），**tracked 零字节**：起手/终态 `git status --porcelain -- internal/ cmd/` 除上面他腿两枚外空；`task.go`／`spill.go` md5 全程 `ded82068…`／`f31147b4…`（与起手抽件逐字节同值）；`internal/config/`、`internal/risk/`、`docs/specs/**`、`scripts/**` 零字节。⛔ 本腿未翻任何框、未动 `Status:`／票名（起手/终态 `grep -c "^- \[ \]"`＝**5**、`^- \[x\]`＝**4** 逐字不变）。

## 留盘未提交（规则 8 只建不删）

`logs/04-gates.txt`（与 03 合并前的同内容版）、`overlay-arm.json`／`overlay-mutA.json`／`overlay-mutB.json`／`overlay-mutC.json`／`overlay-defaults.json`、`mut/`（`task-pristine.go`／`spill-pristine.go`／`task-MUT-V3A.go`／`task-MUT-V3B.go`／`spill-MUT-V3C.go`）。
