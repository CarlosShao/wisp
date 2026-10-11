# 309-r2 · 40 `b.prov == nil` 与 `b.injected` 两肢：注释一行，⛔ 造仪器（job 5）

裁语＝编排者 §4 第 4 件：⛔ 新造夹具，照既有格式在注释里另记一句（具名"这一肢今天⛔ 仪器＋谁该为此响＋归哪一族"）。
落点＝`internal/tools/bridge.go` 的 `assessorFor` 文档块（内容锚 `assembly-is-wrong family, and NEITHER has an instrument`，
现量 `:966-977`，⛔ 行号锚——本票现量段早写了行号会腐烂）。

## 1. 本腿跑的两发（⛔ 沿用 309-v1 的数，自己重跑）

尺＝`logs/prov-rigs.txt`，读数＝`logs/prov-rigs.stdout`＋`logs/provnil-mutation.log`＋`logs/injected-mutation.log`；
台面＝工作树（`internal/tools` 干净⇄本腿注释面），HEAD 面 `6f389389`，14:4x +08；每发 `LANDS` 自证（旧形 0 命中／新形 1 命中）＋当场还原（`RESTORE cmp: byte-identical OK`，sha `069cb275…`）。

| 发 | 摘掉的肢 | 摘后守卫逐字 | `go test ./internal/tools/ -count=1`（整包） |
|---|---|---|---|
| MUT-E | `b.prov == nil` | `if b.injected || taskID == "" {` | **rc=1**，顶层 16 枚红（`TestToolConcurrencyCeilingIsFour`／`TestPerToolTimeoutHonoredViaContext`／`TestRiskDecisionRoutesEachLevelToItsBranch`／`Test236R2TaskCancelRefusesWhenRosterIsUnwired`／`TestTicket224LiveGrantStopsTheL1Question`…） |
| MUT-F | `b.injected` | `if b.prov == nil \|\| taskID == "" {` | **rc=0**（全绿） |
| MUT-0 | ⛔ 摘（基线） | 原形 | rc=0，顶层 PASS 211／FAIL 0／SKIP 0 |

⚠ MUT-F 那条 `grep -c` 读数＝**3** 枚同形，⛔ 是我打出来的：`internal/tools/bridge.go` 里 `OpenTask`(`:848`) 与 `CloseTask`(`:915`) 自己就带
`if b.prov == nil || taskID == "" {` 这枚守卫；LANDS 的判据是**旧形 0 命中**＋`:984` 那一枚新形。别把那个 3 读成"我动了三处"。

## 2. 与 309-v1 那发的差别（★这⛔ 是"编排者裁定里⛔ 成立的一处"，具名顶回）

派单 §5 与票面 §2 写的是："'摘掉 `b.prov == nil` 或 `b.injected`' 两发**各 rc=0**＝仪器对那两肢瞎"。
读了 v1 的脚本 `logs/mut-rig2.txt` 才对上原因：**那一发是 `-run Test309AssessorForEmptyTaskIDDropsTaintDetector`（只跑那枚仪器）**，⛔ 整包。
两发读数⛔ 矛盾，是**两把尺**：

- 只跑 309 仪器（v1 面）：两肢各 rc=0 ⇒ "那枚仪器对这两肢瞎"＝**真**。
- 整包（本腿面）：摘 `b.injected` rc=0＝**整包也瞎**；摘 `b.prov == nil` **rc=1＝整包⛔ 瞎**，
  但响的是**管道噪声**：`prov == nil` 的桥一旦摘掉那一肢就会走到 `b.prov.Detector(taskID)`，
  空指针进去 ⇒ 调用普遍被判成 `L2 审批通道尚未接入…已拒绝执行`（`logs/provnil-mutation.log` 里逐枚可见），
  ⛔ 一枚断言讲的是"这一肢该不该在"。⇒ "**没有具名这一肢的仪器**"这一句对两肢都成立；
  "两肢同色（都绿）"⛔ 成立。⇒ 注释按这个写，⛔ 抄 v1 的"两肢各 rc=0"。

## 3. 交付的那行注释到底说了哪三件事

1. **归哪一族**：`assembly-is-wrong`，⛔ policy（与票面现量早已把 `injected` 归装配族另裁同一口径）。
2. **今天⛔ 仪器**：`b.injected` 整包摘掉也⛔ 响（本腿 rc=0）；`b.prov == nil` 整包响＝管道噪声、⛔ 具名这一肢（本腿 rc=1＋上面那串名册）。
   并写明 309-v1 那对 rc=0 是"scoped 到 309 仪器"那把尺，两把尺⛔ 混用。
3. **谁该为此响**：`whoever first wires that field - or relies on prov being non-nil - owns the test that has to ring`；
   `injected` 今天恒 false 的根据＝**R359-4 写点普查**（`grep -rn 'Assessor:' internal cmd tools --include=*.go`＝2 行＝
   1 枚字段自己的行尾注释＋1 枚 `_test.go` 里的写点 ⇒ 生产写点 0 枚，读数 `logs/assessor-census.log`）。

⛔ 新造夹具／⛔ 动布尔式／⛔ 删肢／⛔ 动 `internal/risk`／⛔ 改 `lit236*`／⛔ 给 `agent.ToolRequest` 加字段——全部守住，
证据见 `50-gates-behaviour.md` 的 B1–B3（剥注释后 blob 与 HEAD `cmp` 同一）。
