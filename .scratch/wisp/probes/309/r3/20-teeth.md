# `309-r3` 新射程有没有牙：逐发原始读数

台件＝`logs/rig-shapes.txt`（可逐字重跑），汇总＝`logs/rig-shapes.log`，逐发＝`logs/rig-{pristine,new-add,new-del,make-add,make-del,conv-add,conv-del}.log`（每枚件末行自落 `rc=N`）。
台面＝**仓外导出树**（`git archive a1f84d39 | tar -x`；⛔ 从共享工作树取变异面、⛔ `cp` 未提交的件 ⇒ ⛔ 复现 `309-r2` 自报的 `Q-309-r2-tree-rigs-mixed-face`）。
指名钉＝`Test309ToolRequestProductionConstructionNail`。共享工作树全程 `git status --porcelain -- internal cmd tools scripts`＝**0 行**（读数＝`rig-shapes.log` 末节，"empty above == 0 rows == clean"）。

每发的"红"都引到用例名＋失败文案本身（⛔ 只报颜色）。三发红句**同一段**，只有名册与分母两处不同，所以逐字引一次、其余两发引名册差值。

---

## 0. 正控（原树，⛔ 任何变异）

```
--- PASS: Test309ToolRequestProductionConstructionNail (0.21s)
rc=0
AC#2 枚数钉：… 生产构造点=1（==1，那枚＝internal/agent/loop.go），分母=295 枚非测试源（每 root=cmd=40 internal=231 scripts=11 tools=13）
```
⚠ 分母＝**295**，⛔ 等于 `309-r2` 报的 293：增量两枚在 `internal/`（`internal=231` ⇄ 那次的 229），来自本波另一枚腿 `304-r1` 已提交的 `internal/llm/adaptertest/**`（尺＝`git diff --name-only 3974ab4d a1f84d39 -- '*.go'` 名册，读数＝`logs/roster-control.log`）。⛔ 归因给我＝那两枚⛔ 是本腿写的；本腿只报"分母今天 295"。

---

## 1. RIG-7 `new(agent.ToolRequest)`

加了什么：`internal/tools/zz309r3new.go`（139 字节，`sha256[:16]=078bace59e060fde`），`package tools`，逐字函数体 `func rig309r3New() *agent.ToolRequest { return new(agent.ToolRequest) }`。
★这一枚就是 `309-v2` 那发打穿 "every/wherever" 的 **B14 同形同位**（v2 记为 `Q-309-v2-3-blind-shape-rings-green`：加进 `internal/tools/` 而钉绿）。

ADD ⇒ 红（`logs/rig-new-add.log`，`rc=1`）：
```
--- FAIL: Test309ToolRequestProductionConstructionNail (0.23s)
toolrequest_nail_309_test.go:390: 枚数钉响：roots=["..\..\cmd" "..\..\internal" "..\..\tools" "..\..\scripts"] 里数到 2 枚 agent.ToolRequest 生产构造点，
票 309 钉的是 1 枚且那枚落在 internal/agent/loop.go（the one newTaskID() dispatch, id always 36 chars, never empty）。
名册=internal/agent/loop.go:675(same-package-bare-literal) | internal/tools/zz309r3new.go:5(new-of-tool-request)。
分母=296 枚非测试源（每 root=cmd=40 internal=232 scripts=11 tools=13）。
```
⇒ 名册真带上新那枚的 **shape 标签 `new-of-tool-request`**（⛔ 只多一个数）；分母 295→296（＝台件自己那枚文件），红键的是分子。

REMOVE ⇒ 同一枚钉绿（`logs/rig-new-del.log`，`rc=0`）：
```
--- PASS: Test309ToolRequestProductionConstructionNail (0.23s)   rc=0
生产构造点=1 … 分母=295（每 root=cmd=40 internal=231 scripts=11 tools=13）
```
还原证明两行（同一棵树里对 `$REV` 的两个 blob `cmp`＋`sha256sum`）：
```
cmp internal/tools/toolrequest_nail_309_test.go: IDENTICAL   sha256[:16] tree=c4aa3d46a6daf3d8  commit=c4aa3d46a6daf3d8
cmp internal/tools/bridge.go: IDENTICAL                      sha256[:16] tree=e6c685da6bfbcfdf  commit=e6c685da6bfbcfdf
```
（另加一行"那枚变异件真被拿掉"＝`removed? still present count: NO`，与 `tree non-test .go count now: 295`。）

---

## 2. RIG-8 `make([]ToolRequest, 0)`（裸名＋长度写法 `0`，落在 `internal/agent`）

加了什么：`internal/agent/zz309r3make.go`（83 字节，`sha256[:16]=099425cc8a58b9e0`），`package agent`，逐字 `func rig309r3Make() []ToolRequest { return make([]ToolRequest, 0) }`。
★一发同时钉两件事：`make` 那枚 built-in 形被计 ∧ **同包裸名**那支在 `internal/agent` 里也管用（`309-r2` 之前整包被跳过那一族）。

ADD ⇒ 红（`logs/rig-make-add.log`，`rc=1`，`--- FAIL: Test309ToolRequestProductionConstructionNail (0.21s)`）；名册差值（红句其余部分与 RIG-7 逐字同一段）：
```
名册=internal/agent/loop.go:675(same-package-bare-literal) | internal/agent/zz309r3make.go:3(make-slice-of-tool-request)
数到 2 枚 … 分母=296 枚非测试源（每 root=cmd=40 internal=232 scripts=11 tools=13）
```
REMOVE ⇒ 绿（`logs/rig-make-del.log`，`--- PASS: … (0.19s)`，`rc=0`，构造点=1、分母回 295）。
还原证明两行：同 RIG-7 那两枚 `cmp … IDENTICAL` ＋ `tree=c4aa3d46a6daf3d8/commit=…`、`tree=e6c685da6bfbcfdf/commit=…`（逐字在 `logs/rig-shapes.log`）。

---

## 3. RIG-9 定义型 `type X agent.ToolRequest` ＋ 转换 `agent.ToolRequest(x)`

加了什么：`internal/tools/zz309r3conv.go`（182 字节，`sha256[:16]=4cd38fe571f279d7`），`package tools`，逐字两行
`type rig309r3X agent.ToolRequest` ＋ `func rig309r3Conv(x rig309r3X) agent.ToolRequest { return agent.ToolRequest(x) }`。

ADD ⇒ 红（`logs/rig-conv-add.log`，`rc=1`，`--- FAIL: Test309ToolRequestProductionConstructionNail (0.21s)`）；名册差值：
```
名册=internal/agent/loop.go:675(same-package-bare-literal) | internal/tools/zz309r3conv.go:7(conversion-to-tool-request)
数到 2 枚 … 分母=296 枚非测试源（每 root=cmd=40 internal=232 scripts=11 tools=13）
```
REMOVE ⇒ 绿（`logs/rig-conv-del.log`，`--- PASS: … (0.15s)`，`rc=0`，构造点=1、分母回 295）。
还原证明两行：同前两节那两枚 `cmp`/`sha256[:16]` 对。

---

## 4. 六发汇总表（每发一枚件、件末自落 `rc=`）

| 发 | 加/拿 | 位置 | 具名用例 | 原始 rc | shape 标签（名册里真出现的那串） |
|---|---|---|---|---|---|
| 正控 | ⛔ 变异 | `git archive a1f84d39` 原树 | `Test309ToolRequestProductionConstructionNail` | **rc=0** PASS | 构造点=1 |
| RIG-7 | 加 | `internal/tools/` | 同一枚 | **rc=1** FAIL | `new-of-tool-request` |
| RIG-7 | 拿 | 同上 | 同一枚 | **rc=0** PASS | 回到 1 枚 |
| RIG-8 | 加 | `internal/agent/` | 同一枚 | **rc=1** FAIL | `make-slice-of-tool-request` |
| RIG-8 | 拿 | 同上 | 同一枚 | **rc=0** PASS | 回到 1 枚 |
| RIG-9 | 加 | `internal/tools/` | 同一枚 | **rc=1** FAIL | `conversion-to-tool-request` |
| RIG-9 | 拿 | 同上 | 同一枚 | **rc=0** PASS | 回到 1 枚 |

---

## 5. 整包名册改前⇄改后（同面正控，⛔ 一把正则推整族）

台件＝`logs/roster-control.txt`，读数＝`logs/roster-control.log`；两棵树都出自 `git archive`（同一枚台面）：
`BEFORE=3974ab4d5b49…`（本腿起手锚）、`AFTER=a1f84d39a929…`（本腿 (B) 那笔）。
名册尺＝`grep -E '^(    )*--- (PASS|FAIL|SKIP): '` → 剥耗时列 `sed -E 's/ \([0-9.]+s\)$//'` → `sort` → `comm`。

| 尺 | before | after |
|---|---|---|
| 名册行数（含子测试） | 279 | 279 |
| 顶层 `--- ` 枚数 | 211 | 211 |
| 子测试 `    --- ` 枚数 | 68 | 68 |
| `FAIL:` 枚数 | 1 | 1 |
| `SKIP:` 枚数 | 0（empty above＝0 枚） | 0（empty above＝0 枚） |
| 整包 `go test` rc | 1 | 1 |
| **newly RED**（after 有、before 无的 FAIL） | — | **`newly_red_count=0`**（`logs/fail-new.txt` 空＝0 枚） |
| added / removed 名册差 | — | added 1／removed 1，两枚**同一枚子测试**：`TestTicket252R1BothSpellingsAnswerTheSame/<树名>`——用例名里嵌了 `mktemp` 目录名，⛔ 真新增⛔ 真删除（具名报，⛔ 洗成 0/0） |

唯一那枚红（两发同一枚）＝`--- FAIL: TestTicket252R1AlignmentAddsNoAuthorization`＝导出树缺件那一族，`309-r2`／`309-v2` 各自在两棵树里量到同一枚，归票 298 `AC#5` 已承接的账，⛔ 本票、⛔ 本腿造成（改前那发就红）。

⚠ 区间纯度具名：`3974ab4d..a1f84d39` 里除了本腿三笔还有 `2a6f11a5`（编排者落账）与 `304-r1` 两笔（`b9f37892`、`e875f00f`，写面 `internal/llm/**`），
所以"改前⇄改后"两棵树的 `.go` 差值＝本腿两件 ＋ 别腿四件（`git diff --name-only … -- '*.go'` 名册在 `logs/roster-control.log`）。⇒ 名册⛔ 新增红这一条仍然只说明"那一整段区间⛔ 让 `internal/tools` 多红"，⛔ 是本腿单独的因果证明；本腿单独的因果证明＝上面那六发＋`30-gates.md` 的 `go vet`／`go build`。
