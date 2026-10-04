# 证据件 265-r1 — 票 265 形 ⓐ-Ⅰ 落地（晚绑定 holder，`GrantRecorder` 由 `cmd/wisp` 自己实现）

**腿**：`265-r1`（产码落地腿）
**对象票**：`.scratch/wisp/issues/265-resident-gate-has-no-session-grant-writer.md`
**裁定出处**：票面「编排者裁定（形）＝ⓐ-Ⅰ」＋台账 `A601` §4／§5、`A602`（具名解冻那枚 AST 钉）
**本件九节状态**：骨架先落盘（每格初值＝**未判**），此后每做完一格 Edit 一节并 commit 一次。

---

## 0. 本腿的运行约束与开工闸门

（本节答：开工闸门两条的现跑读数＋时刻；`git status --porcelain -- cmd internal`；在飞腿名单与本腿的归因边界；sherpa PATH 这一格怎么解决。）

- 闸门① `git status --porcelain -- cmd internal`：**未判**
- 闸门② `go test ./cmd/wisp/` 需要 sherpa PATH：**未判**
- 闸门③ `winlive` 一律不跑：**未判**

## 1. 形状复认（ⓐ-Ⅰ 三枚写面的开工时刻现量）

（本节答：三枚文件在 `git cat-file blob HEAD:` 那儿的原始行号与原文；写面边界——`internal/agent/approval/**` 与 `cmd/wisp/run.go` 零改动怎么自证；有没有越界。）

- 1.1 `cmd/wisp/resident_approval_windows.go` 建门字面量与 `:185-192` 那段注释：**未判**
- 1.2 `cmd/wisp/resident_task_source_windows.go` 的绑定位时序（`:281` 之后、`submitTask` 之前）：**未判**
- 1.3 `cmd/wisp/resident_approval_risk_256_windows_test.go` 的 AST 钉（`:445`／`:452-455`）：**未判**
- 1.4 写面零越界自证：**未判**

## 2. 产码：holder 的形状与晚绑定接线

（本节答：holder 类型逐字设计；未绑定窗口为什么必须走 error 路径（硬约束第 1 条）；`run.session == nil`（mint 失败）落在哪一支；为什么只能晚绑定（硬约束第 2 条）；账本唯一构造点没被动（硬约束第 3 条）。）

- 2.1 holder 类型与 `Record` 的未绑定支：**未判**
- 2.2 建门字面量 `Grants:` 的实传（5 枚 → 6 枚）：**未判**
- 2.3 绑定点与它的前后邻行：**未判**
- 2.4 `:185-192` 那段注释改写后的原文（不许留反话）：**未判**
- 2.5 顺带同批改释两枚（`resident_approval_windows.go:32-33`／`resident_task_source_windows.go:52-53`）：**未判**

## 3. 判据：新建的钉与解冻后的钉

（本节答：普查件 §5.2 具名的"常驻腿那一格零枚用例守"补了哪几枚；每枚断什么、正控形是什么；256 那枚 AST 钉解冻后的期望集与反向断言搬到哪儿去了。）

- 3.1 新建判据件与用例名册：**未判**
- 3.2 256 AST 钉的解冻落地（A602 边界内）：**未判**
- 3.3 硬约束第 1 条的钉（holder 未绑定 ⇒ 指名用例红）：**未判**
- 3.4 AC#3 ⓐ 支的正控（摘掉接线 ⇒ 指名用例红）：**未判**

## 4. 突变自证（每一发：种下必红 ＋ 还原 ＋ md5 三读数）

（本节答：每一发突变的原样命令、红句**逐字**、还原用的副本出处（必须是 `git cat-file blob HEAD:<path>`）、起手＝还原后＝HEAD blob 三枚 md5。）

- M1 holder 未绑定 ⇒ 让 `Record` 静默 `return 0, nil`：**未判**
- M2 摘掉 `Grants:` 那一枚实传：**未判**
- M3 摘掉绑定位（`:281` 之后那一发）：**未判**
- M4 种"第二枚 mint"形（证本腿没走 ⓐ-Ⅲ）：**未判**

## 5. AC#1 那两句假话 → 真话的兑现

（本节答：`approval_reply.go:279-281`（给用户）与 `:277-278`（给审计）两句**文字一字未动**的自证；写侧接上之后两句为什么为真；判据怎么指得出"写侧没接上"这一情形。）

- 5.1 两句的文字零改动：**未判**
- 5.2 兑现链（holder → `session.Ledger.Record` → `GRANT-RECORDED`）：**未判**
- 5.3 判据对"没接上"的敏感度：**未判**

## 6. 门禁四数读数（带取数时刻）

（本节答：`sh scripts/d22scan.sh`、`sh scripts/check-path-length-budget.sh --with-self-test`、`gofumpt -l <本腿文件>`、`go vet ./cmd/wisp/` 四门读数＋时刻；定向 `go test` 的红名集合逐名比对，哪几枚是既有红。）

- 6.1 d22scan：**未判**
- 6.2 路径长度门：**未判**
- 6.3 gofumpt：**未判**
- 6.4 go vet（含 `-tags winlive` 只编译不跑）：**未判**
- 6.5 定向 go test 红名集合逐名比对：**未判**

## 7. 我判不动／量不到的地方（具名＋归口）

（本节答：⛔ 不许用"应该没问题"填空。每一格给：判到哪一步就停了、归口哪张票／哪枚编排者面。）

- 7.1 第三枚过期注释 `internal/agent/approval/approval_always.go:165`（不在本腿写面）：**未判**
- 7.2 其余：**未判**

## 8. 污染面自证与提交名册

（本节答：本腿动过的每一枚文件（包级路径＋numstat）；临时件只建在 `.scratch/wisp/probes/265/r1/`；commit 名册与 pathspec 逐枚；票面 AC 框零改动的自证。）

- 8.1 写面与 numstat：**未判**
- 8.2 提交名册：**未判**
- 8.3 票面 AC 框零改动：**未判**
