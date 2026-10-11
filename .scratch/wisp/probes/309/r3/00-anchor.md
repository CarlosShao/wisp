# `309-r3` 起手锚（写腿，⛔ 裁决者）

## 0. 起手读数（本轮第 2 次调用现跑，⛔ 抄派单）

| 读数 | 值 | 尺（可逐字重跑） | 射程 | 面 |
|---|---|---|---|---|
| HEAD | `3974ab4d` | `git rev-parse --short HEAD` | 整仓 | 工作树所属仓 |
| 分支 | `dev` | `git rev-parse --abbrev-ref HEAD` | 整仓 | 同上 |
| 工作树脏行数 | **0** | `git status --porcelain -- internal cmd tools scripts \| wc -l`＝件 `logs/st.txt`，`wc -l`＝0 | `internal` `cmd` `tools` `scripts` 四根 | 工作树 ⇄ HEAD |
| 全索引暂存枚数 | **0** | `git -c core.quotePath=false diff --cached --name-only \| grep -c ''`（`grep -c` 命中 0 时 rc＝1 ⇒ 该发按 `N=$(… \|\| true)` 取数，**empty above＝0 枚**） | 全索引（⛔ 归因，只报数） | 索引 |

## 1. 别腿在飞（具名，⛔ 我动它）

- 派单 §0.2 具名 `304-r1`（写面 `internal/llm/**`）此刻可能在同一工作树/索引里 commit。
- ★我这一刻的读数：暂存枚数 **0**、四根工作树 **0 行** ⇒ **我没看到任何别家路径**，也**没清空/没归位过任何东西**（⛔ `reset`／`stash`／`checkout .`／`clean`；允许的归位动作只有 `git restore --staged -- <我自己的路径>`，本轮未用到）。
- 我的名册尺一律**限定在自己的 pathspec 内**：`internal/tools/**` ＋ `.scratch/wisp/probes/309/r3/**` ＋ 票 309 一枚文件。

## 2. 写面（派单 §2.3，⛔ 扩一格）

1. `.scratch/wisp/probes/309/r3/**`（本腿自家件与 `logs/`）
2. `internal/tools/toolrequest_nail_309_test.go`（补射程，任务 (A)）
3. `internal/tools/bridge.go` 的**注释行**（改正 "every/wherever" 那句，任务 (B)；⛔ 动语义字节）
4. `.scratch/wisp/issues/309-*.md` 末尾**一条** append-only bullet（1 加／0 删）

## 3. 本笔 commit 的范围

只交 `00-anchor.md`（＋本件的消息件不进 commit，落 `logs/`）。`git add` 与 `git commit` 两处都带显式 pathspec。
