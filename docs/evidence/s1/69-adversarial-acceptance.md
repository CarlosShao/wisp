# 票 69 对抗验收 —— C21 令牌表的双向机器检查（A24-D4）

**验收人**：编排者（**非实现者**；实现方是 `agent-ticket69`）
**验收时间**：2026-09-21 09:36
**被验收的 commit**：`4ca2c6c`（断言落地）、`9c29555`（三次变异检验）
**基线**（用来证明"没弱化"）：`75812ac`
**证据档位标记**：〔独立复现〕= 我亲自敲命令并读真实输出；〔日志＋归档，我抽验〕= 数字来自代理但我核对文件存在且字段对得上；〔仅自述，不背书〕。

## 裁决表（与票面 4 个 AC 框 1:1）

| AC | 裁决 | 我跑的那条判据 | 档位 |
|---|---|---|---|
| **AC#1** 双向机器检查落地 | **PASS** | `go test ./internal/ball/ -run 'Token\|C21\|Table'` → `ok 0.235s`；并在我的独立变异里看到它**自报覆盖面**：`checked 56 constant declarations across 25 geometry rows against 56 exported tokens.go constants`、配色 114 行。表→码与码→表两个方向都在同一份输出里留了字。 | 〔独立复现〕 |
| **AC#2** 变异检验（新断言真会红） | **PASS（我自己重做了一次）** | 我**不采信代理的日志**，在独立树里把 `internal/ball/tokens.go:361 DockTriggerPx = 16` 改成 `17`，先 `grep -n` 证变异落盘（`:361: DockTriggerPx = 17`）再跑：**唯一变红的是 `TestC21GeometryRowsMatchCodeConstants`**，报文 `tokens_table_test.go:838: DockTriggerPx = 17, but this row states no such number (its numbers are: 160ms, 0.42, 16px, 96dpi, 32) - table and code drifted`，**点名表行 155** ⇒ 与代理 M1 的叙述逐字一致。 | 〔独立复现〕 |
| **AC#3** 不许弱化既有 20 条金标准 | **PASS** | `git diff --stat 75812ac HEAD -- internal/ball/tokens.go internal/ball/tokens_test.go` → **输出为空**（两个文件一字未改）。这条我按"最省事做法是先动旧断言"的预期专门查过：没有动。 | 〔独立复现〕 |
| **AC#4** 门禁 | **PASS** | `gofmt -l internal/ball` → 空、rc=0；`go vet ./internal/ball/` → rc=0；`go test ./internal/ball/ -count=2 -v` → `PASS / ok 1.110s`。**计数不变式**：`=== RUN` **108**、顶格 `--- PASS` **96 = 48 个不同测试名 × 2**（正是代理自报的 48）、`--- FAIL\|--- SKIP` **0**、整份日志字符串 `SKIP` 出现 **0** 次（108−96=12 条是缩进的子测试，不是漏跑）。跑之前与跑之后 `git status --porcelain internal/ball` **均为空** ⇒ 这份绿是 HEAD 的绿，不是别人 WIP 的绿。 | 〔独立复现〕 |

## 结论

**PASS，4/4，可以 `-done`。** 本票的价值确实拿到了：在我复现 AC#2 之前，"表写 16px、码是 17"
这类漂移对 CI **完全隐形**（代理的 M1 记录里既有 47 条全绿，含 `dock_test.go:118`——它用
`const trigger = DockTriggerPx` 推导，改值永远不会红；这条我认同其推理，因为它正好解释了我为什么
只看到一条红）。

## 本票**没有**证明的（必须带着走，不许沉底）

代理自己在报告里列了四条没证的，我逐条判归属：

1. **几何行是"行内包含"而不是双射** —— 两行可以互相掩盖同一个数字。
   → 已写进 **票 74 AC#5 段尾**："containment→bijection 只有在能举出一个今天会漏过去的形状时才许收，
   否则记为缺口"。**归属明确，不重开本票。**
2. **`FontFamily` 值无断言** → 同上，票 74 面里点名。
3. **零消费者按字段名匹配，同名字段会偏乐观** → 票 74 的第 3 族（29 条漂移）正是这件事的收口。
4. **表 ↔ `tokens.css` 仍无人机器检查** → 票 74 **AC#4** 硬性要求：要么装跨检，要么写成带票号的
   书面接手，"暂缓"无票号即不合格。

⇒ 四条缺口**全部有带编号的落点**，不构成把本票压成 PARTIAL 的理由。

## 一处需要后来者知道的读日志陷阱（验收时撞到的）

`TestC21GeometryRowsMatchCodeConstants` 变红时会**同时**倒出与本次变异无关的两张清单
（`MATCHED WITHOUT THE PROMISED UNIT (4)`、`NUMBERS NO NAMED CONSTANT CLAIMS (21)`）。
第一次读那份输出很容易把噪声当成变异的后果 ⇒ **只认点名"我改的那个常量"的那一行**。
另：为做这个变异我在仓库根建过一次嵌套 worktree，**这是错的**（并发代理的 `d22scan` 文件遍历会多算
一整份源码，正打在票 70 的"生产文件数下限"守卫上），已 `prune` + 删除，登记在 **registry A38④**。
