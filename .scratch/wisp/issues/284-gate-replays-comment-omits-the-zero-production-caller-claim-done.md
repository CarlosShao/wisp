# 票 284 — `Gate.Replay` 的注释只写了契约、**没写"零产码调用者／勿删"**：随下一枚真动 `gate.go` 的腿同批补

**立票**：2026-10-09 10:1x 编排者（机主令「及时补票」；来路＝票 278 `AC#2` 裁完剩下的那一格＋台账 `A754`）
**性质**：一行注释的诚实化。**⛔ 不接线、不动行为、不新立仪器**。

## 现量（引用前先重跑；行号与短语同一次取数）

- `HEAD:internal/agent/approval/gate.go:751-754` 的注释块**逐字**写着它是冻结契约的实现半件：`Replay re-displays a refused or expired L2 request under a fresh correlation id (C18 一键重放). It requires the task to be admitted again (D47 still applies to a replay), and it is NOT an answer: the new item must still be approved with a newly minted grant.`
- 定义在 `HEAD:internal/agent/approval/gate.go:755`，唯一下游 `Queue.replay`（`queue.go:598`）与 `replayOf` 字段（`queue.go:69`）。
- **零产码调用者**（编排者自己现跑，整族、剥测试靠路径面）：`git grep -nE "[A-Za-z0-9_]\.Replay\(" HEAD -- ':(exclude).scratch' '*.go' ':!*_test.go'` ＝ **rc=1（空）**；唯一调用点在 `internal/agent/approval/queue_test.go:283`。
- **入口今天结构性不存在**：答案侧两枚接口面都不声明它——`NativeAPI`（`ui.go:143-162`：`Allow`/`AllowSession`/`Reject`）、`PanelAPI`（`ui.go:167-171`：`Reject`/`Head`/`View`）；控制台动词条（`cmd/wisp/approval_reply.go:566-584`）里没有 `replay`。
- 对照件：同族另一枚"建了没接"的 `AskOnTaskRoot` **已经替自己说过话**（`HEAD:cmd/wisp/resident_approval_windows.go:690-693`，逐字 "today its only callers are this package's own cases … rather than letting the existence of the seam read as a running pipeline"）⇒ **本票只求把同一件事在 `Replay` 头上也说出来，形状一致、读者不会再猜。**

## 要建什么

- [x] **AC#1 一句话落地**：在 `gate.go` 的 `Replay` 注释块补一句，内容三要素齐＝〔今天零产码调用者／非死代码、勿删／它替 C18 一键重放留的〕。**⛔ 不许在注释里嵌行号**（本仓刚抓到四处漂：票 278 的 `:697→:698`、`:700→:701`、`:248→:249`、台账 `:749→:755`）——要指就指**具名符号**（`Queue.replay`／`NativeAPI`／`PanelAPI`）。
- [ ] **AC#2 随批同交**：⛔ **不许为这一行单独开一条腿去改 `gate.go`**（那是白占一次 Go 编译面）。本票的落地时机＝**下一枚本来就动 `internal/agent/approval/gate.go` 的腿**，把这一行并进它那一次改动与那一次门禁读数里；在此之前本票只登记。 ⚠实际由单腿 556f0483c 落地＝纪律偏差(A1554)，本框不勾；实质已落 gate.go:872-884
- [x] **AC#3 改完的自证**（由那枚随批腿交，本票不预先勾）：`go vet ./internal/agent/approval/` rc=0＋`sh scripts/d22scan.sh` 纯净树 rc=0；⛔ 注释里**不许写会被自己扫的现在时计数**（本仓定式：把被数的 token 写进注释，"0 次"在那一刻就变成假话）。 ⚠来路条款「由随批腿交」未兑现(A1554)

## 禁区

- ⛔ 不接线（接 `Replay` ＝给"重放"造入口，属功能级，要 owner 一句话）；⛔ 不动 C18 契约文本（`PLAN.md`／`docs/specs/**` 一字节不动）；⛔ 不新立"零调用者"负向钉（票 278 `AC#2` 已裁：负向钉＝真接线那天必须改期望值，正是"为变绿动断言"的温床，且这一族注释腐烂率最高）。
- ⛔ 三枚冻结件一字不动（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）；`frontend/**`／`design/**` 零读零写零转述；⛔ 零翻框（归编排者）；⛔ 零 push。

---
**结案 2026-11-05（A1554）**：非实现者裁决 `docs/evidence/s1/284-verdict.md`＝AC#1 PASS（gate.go:872-884 三要素齐、无行号、具名符号）＋AC#3 PASS（go vet rc0、d22 clean×2）；AC#2 为单腿落地（556f0483c）＝明写过程偏差、不勾。⛔ 不对外称三 AC 全绿。附带更正：实现者证据§2「块内无数字」不实（现跑 4 个数字＝契约号 L2/C18×2/D47）⇒ 口径改「无计数式数字、只有契约号」；交付件 284-gate-replay-zero-caller-comment.md 此前从未提交，本次补提交。
