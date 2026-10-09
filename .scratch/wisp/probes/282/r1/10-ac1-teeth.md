# 票 282 `AC#1` 行为尺 — 种刀牙口读数（腿 282-r1，2026-10-09）

票面判据句（`sed -n '14p'` 逐字，未改一字）：

```
- [ ] **AC#1 行为尺**：把"任务身份改读访问器"那一跳**种坏**（例：让 `TaskID(ctx)` 返回空或回退到 corr），**指名用例必须红**——若全绿 ⇒ 具名写"该支今天没有仪器"。四件套（种前 hash／改的那行 `sed` 复量／红句 file:line／还原 hash 逐字回＋porcelain 空）。
```

## 0. 起手锚与形状

- 起手 HEAD：`6041d7f9369c7b0963e23d9587728ead63bc3871`（本腿第 1 笔 commit 后＝`c2b422a490ca30f1a973c1635680b0c4b606d789`）
- 基线（⛔ 已带 `-count=1`，无缓存重放）：`go test ./internal/tools/ ./internal/agent/ -count=1 -v`
  - `internal/tools`：**PASS=209 / FAIL=0 / SKIP=0**，`ok github.com/CarlosShao/wisp/internal/tools 13.803s`
  - `internal/agent`：**PASS=84 / FAIL=1 / SKIP=0**，`FAIL github.com/CarlosShao/wisp/internal/agent 1.861s`
  - 那枚既有红（不归本腿、不修不 Skip、归票 286）：`--- FAIL: TestGoldenSingleToolCall (0.00s)`
- 访问器真身（`git rev-parse` 面）：`internal/tools/cancel.go:67 func TaskID(ctx context.Context) string`，体为 `h, _ := cancelOf(ctx); return h.taskID`。
- 种刀形状选择：两发都走票面举例的**「回退到 corr」**那一支（把调用点读的那枚身份换回 `CorrelationID(ctx)`），不走「返回空」。
  理由（现量）：「返回空」两发都会塌进同一处 fail-closed 支，而那两个空支已被 `failclosed_236_teeth_test.go` 逐字面量钉着，
  种它是测一枚已经有人看的门；「回退到 corr」才是票 242 落地前那个真实缺陷形状，也是本票要复核的那一跳。

## 1. 发一 · 种在 `internal/tools/subagent_197.go:262`

四件套：

- 种前 hash：`git hash-object` = `8c9266a7dfaf62fbf9418cf8055a5905eb98b8b1`；
  `sha256sum` = `f2248443bc7001ad9892ffe434b1d902f81bf1c46556eca101743129059c10ec`；
  该行 `sed -n '262p'` 原文 = `	parentID := TaskID(ctx)`（该文件 CR=0，工作树 hash 与 `HEAD:internal/tools/subagent_197.go` 逐字相同）。
- 种后那一行 `sed -n '262p'` 复量 = `	parentID := CorrelationID(ctx)`
  LANDS 自证：种后该文件 `CorrelationID(ctx)` 命中数 = `1`（种前该文件此形状命中 0，`grep -n` 空）。
  种后 hash = `aef8339652277f75c8043db055fc9a42fd00a835`（与种前不等 ⇒ 突变确实落盘）。
- 指名用例红（同一命令形状 `go test ./internal/tools/ ./internal/agent/ -count=1 -v`，`testrc=1`）：

  ```
  --- FAIL: Test283IdentityChainThroughTheRealLoop (0.06s)
      internal/tools/ticket283_corr_identity_rulers_test.go:194: 子代理行的 ParentTaskID = "caf08939-ea19-4564-bbc5-deb69566e3b2#call-283-spawn", want 宿主 task id "caf08939-ea19-4564-bbc5-deb69566e3b2"（task.spawn 读的是 TaskID(ctx)，不是 per-call corr）
  ```

  同发连带红（不是另一发，是同一条用例的第二条断言）：

  ```
      internal/tools/ticket283_corr_identity_rulers_test.go:206: task.cancel 的回执没有走到权威判定通过那一支: "本用例不带 L2 队列"
  ```

  该发整包名册：`--- FAIL: Test283IdentityChainThroughTheRealLoop` ＋既有红 `--- FAIL: TestGoldenSingleToolCall`，两枚，无第三枚。
  计数：`internal/tools` **PASS=208 / FAIL=1 / SKIP=0**（`FAIL github.com/CarlosShao/wisp/internal/tools 13.798s`）、
  `internal/agent` **PASS=84 / FAIL=1 / SKIP=0**（`FAIL ... 1.866s`）。
- 还原：`sed -n '262p'` = `	parentID := TaskID(ctx)`；
  `git hash-object` = `8c9266a7dfaf62fbf9418cf8055a5905eb98b8b1`（与种前**逐字等值**）；
  `sha256sum` = `f2248443bc7001ad9892ffe434b1d902f81bf1c46556eca101743129059c10ec`（与种前逐字等值）。
  `git status --porcelain -- internal cmd` = **空**。

判语：**这一支有仪器，且指名用例确实红。**

## 2. 发二 · 种在 `internal/tools/task.go:708`

四件套：

- 种前 hash：`git hash-object` = `3ec8d498f1613f0b7eeadb0b840d4c77cc04ed84`；
  `sha256sum` = `e74197192745af05b3677e8ed38315a68a74c10621dbd46d0795f6f7e59d7bf0`；
  该行 `sed -n '708p'` 原文 = `	caller := TaskID(ctx)`（该文件 CR=0，工作树 hash 与 `HEAD:internal/tools/task.go` 逐字相同）。
- 种后那一行 `sed -n '708p'` 复量 = `	caller := CorrelationID(ctx)`
  LANDS 自证：种后该文件 `grep -n 'CorrelationID(ctx)'` = 唯一命中 `708:	caller := CorrelationID(ctx)`（种前该文件此形状命中 **0 枚**，`grep -n` 输出为空）。
  种后 hash = `cedb4762a1edd577b6072a0ae5738abb44300fea`（与种前不等 ⇒ 突变确实落盘）。
- 指名用例红（同一命令形状，`testrc=1`）：

  ```
  --- FAIL: Test283IdentityChainThroughTheRealLoop (0.07s)
      internal/tools/ticket283_corr_identity_rulers_test.go:206: task.cancel 的回执没有走到权威判定通过那一支: "拒绝停止：f4f7abb1-9a1c-456e-a959-5fa17112b75f 是任务 cbe4af20-760f-464e-9ba7-cc48f6bd25d3 派生的孩子，不是调用者 cbe4af20-760f-464e-9ba7-cc48f6bd25d3#call-283-cancel 的孩子——只有父任务能停自己的孩子，兄弟之间、旁支之间都停不了（票 221 AC#3 的正控）。"
      internal/tools/ticket283_corr_identity_rulers_test.go:209: 父任务停自己的孩子被身份判定拒了（调用者读到的不是 task id）: "拒绝停止：f4f7abb1-9a1c-456e-a959-5fa17112b75f 是任务 cbe4af20-760f-464e-9ba7-cc48f6bd25d3 派生的孩子，不是调用者 cbe4af20-760f-464e-9ba7-cc48f6bd25d3#call-283-cancel 的孩子——只有父任务能停自己的孩子，兄弟之间、旁支之间都停不了（票 221 AC#3 的正控）。"
  ```

  该发整包名册：同样只有 `Test283IdentityChainThroughTheRealLoop` ＋既有红 `TestGoldenSingleToolCall` 两枚。
  计数：`internal/tools` **PASS=208 / FAIL=1 / SKIP=0**（`FAIL ... 14.059s`）、`internal/agent` **PASS=84 / FAIL=1 / SKIP=0**（`FAIL ... 1.925s`）。
  同形状未被这发打红的邻尺（逐枚点名，用来界定射程）：
  `--- PASS: Test221ParentStopsItsOwnChildRowAndStreamSettle (0.06s)`、
  `--- PASS: Test221SubagentCannotStopSiblingOrItself (0.00s)`、
  `--- PASS: Test236R2TaskCancelRefusesWhenHostGaveNoCallerID (0.00s)`、
  `--- PASS: Test285RosterRowsCarryDistinctCorrPerCall (0.07s)`、
  `--- PASS: Test197ChildCannotDeriveSubagent (0.00s)`。
  ⇒ 父停子在**直接派发**夹具（`corr == taskID`，桥只做 `CorrelationID → TaskID` 回落）下分不出两枚身份，只有真 loop 那条链分得出——这正是票 283 头注写的盲区二，本发实测复现。
- 还原：`sed -n '708p'` = `	caller := TaskID(ctx)`；
  `git hash-object` = `3ec8d498f1613f0b7eeadb0b840d4c77cc04ed84`（与种前**逐字等值**）；
  `sha256sum` = `e74197192745af05b3677e8ed38315a68a74c10621dbd46d0795f6f7e59d7bf0`（与种前逐字等值）。
  `git status --porcelain -- internal cmd` = **空**。

判语：**这一支也有仪器，指名用例确实红。**

## 3. 方法学前置的现量确认（票面 AC#3 那句，本尺同样吃它）

「夹具必须 `TaskID != CorrelationID`」——本腿不引注释、不引别人结论，两枚现量：

1. 发二的红句里两枚身份**同屏出现且不相等**：调用者 `cbe4af20-760f-464e-9ba7-cc48f6bd25d3#call-283-cancel`（corr）
   对 父任务 `cbe4af20-760f-464e-9ba7-cc48f6bd25d3`（task id）。发一的 :194 同形（`...#call-283-spawn` 对 task id）。
   ⇒ 真 loop 链上两枚身份确实是两个值，前置成立。
2. 干净树基线那一跑（同一命令 `go test ./internal/tools/ ./internal/agent/ -count=1 -v`，未加 `-run`）里逐枚现量到绿：
   `--- PASS: Test283TaskIDAccessorIsTheTaskIDNotTheCorr (0.00s)`（它自带 handle `corr:"corr-283"` 对 `taskID:"task-283"`）、
   `--- PASS: Test283IdentityChainThroughTheRealLoop (0.06s)`、
   `--- PASS: Test285RosterRowsCarryDistinctCorrPerCall (0.06s)`、`--- PASS: Test285CallCorrIsDistinctPerCall`、
   `--- PASS: Test285LoopDispatchesOneCorrPerCall`。
   尺件归属已核：`git log --oneline -1 039ec93c` = 「283-r1 补尺…（新文件，零碰既有件）」、`git show --stat` =
   `internal/tools/ticket283_corr_identity_rulers_test.go | 227 +++++`（227 行新建零删除，与票面 §AC#3 那句一致）；
   `git log --oneline -1 9a00a890` = 「票 285 AC#4 第五形补尺…」；两枚均 `git ls-tree HEAD` 在树内。

## 4. 两发的门禁与自检

- `gofmt -l internal/tools/subagent_197.go internal/tools/task.go` = 零输出，`rc=0`
- `gofumpt -l`（`D:\work\base\gopath\bin\gofumpt.exe`，只读档）同两枚文件 = 零输出，`rc=0`
- `GOFLAGS= go build ./...` `rc=0`
- `sh scripts/d22scan.sh`（纯净树）`rc=0`，末行 `d22scan: clean - no D22 ban violations`

## 5. 本腿自己的一次失败尝试（照实带，不藏）

发一第一次尝试用了 `sed -i '262s|.*-|\t...|'`，样式里的 `.*-` 要求那一行含 `-` 字符，而 `parentID := TaskID(ctx)` 不含 ⇒
**突变根本没落盘**：那次跑完后 `sed -n '262p'` 仍是原句、`git hash-object` 仍等于种前 `8c9266a7…`、名册只有既有红 `TestGoldenSingleToolCall` 一枚。
⇒ 那一发**不算读数**（它只是又复跑了一遍干净树），本件§1 的发一是重跑的那一发（样式锚换成 `TaskID(ctx)`→`CorrelationID(ctx)`，并加了 LANDS 自证）。
教训与既有台件腐烂同形：**行号锚的突变必须自带落地自证**，本腿此后每发都打「种后 `grep -n` 命中」这一行。

## 6. 收工三数（还原后的纯净树，与基线逐字同形）

命令同上（`go test ./internal/tools/ ./internal/agent/ -count=1 -v`），`finalrc=1`（只因既有红）：

- `internal/tools`：**PASS=209 / FAIL=0 / SKIP=0**，`ok github.com/CarlosShao/wisp/internal/tools 13.727s`
- `internal/agent`：**PASS=84 / FAIL=1 / SKIP=0**，`FAIL github.com/CarlosShao/wisp/internal/agent 1.850s`
- 整包 FAIL 名册只有一枚，且就是起手那枚既有红：`--- FAIL: TestGoldenSingleToolCall (0.00s)`
- 与基线（tools 209/0/0 · agent 84/1/0）**逐格相等** ⇒ 两发突变没有留下任何字节。
- 门禁：`GOFLAGS= go build ./...` **rc=0**；`sh scripts/d22scan.sh`（纯净树）**rc=0**（末行 `d22scan: clean - no D22 ban violations`）。
- 禁区零改动自证：`git status --porcelain -- internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go internal/agent/approval/gate.go internal/agent/loop.go internal/tools/loop_approval_test.go docs internal/observe/thresholds.go` = **空**；
  `git status --porcelain -- internal cmd | wc -l` = **0**。

## 7. 本格结论

`AC#1` **两发各有指名红** ⇒ 不需要写"该支今天没有仪器"；两枚调用点（`subagent_197.go:262`／`task.go:708`）的
"任务身份改读访问器"那一跳**都有牙**，且红全部落在票 283 尺二 `internal/tools/ticket283_corr_identity_rulers_test.go:194/206/209` 上。
零断言被改、零 `t.Skip`、零产码字节残留。本格不翻框（归编排者）。
