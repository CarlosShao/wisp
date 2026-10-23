# 票 196 契约文字腿 `196-c1` · 裁决表（写腿：契约文字）

**这一腿做了什么**：把 `task_log.state` 的 **12 枚值域名册**写进两处契约文字——`docs/specs/SPEC-02-data-storage.md` §3 与
`docs/PLAN.md` D35 的 `task_log` 块——让"**代码里有 `CHECK (state IN (…12 枚…))`、契约文字一个字都没有**"这个落差消失。
**只加这一件事**，⛔ 改任何既有句子、⛔ 动 D 号/C 号、⛔ 动任何 `.go`。

## 1. 机主批准（逐字，这是本腿唯一的依据）

> 「**1. 可以**」

上文＝编排者问"票 196 那 12 枚名册要不要写进 `SPEC-02` §3 与 `D35` 的 `task_log` 块"。
⇒ 本腿获准改的**只有这两处契约文字**、且**只准加"`task_log.state` 的值域名册/CHECK"这一件事**。

## 2. 12 枚从代码提取（⛔ 手抄、⛔ 从票面抄）

**唯一来源**＝`internal/memory/taskstate.go`：`const` 块 7 枚（旧词）＋ `taskStateRoster` 里 5 枚（D43 名）。
CHECK 的值表由 `schema.go` 的 `buildDDLV4()` 用 `taskStateRoster` 生成（`schema.go:294-324`，值表 = 该切片顺序）。

**逐字命令**：

```sh
sed -n '/^const (/,/^)/p' internal/memory/taskstate.go \
  | sed -n 's/^[[:space:]]*TaskState[A-Za-z]*[[:space:]]*TaskState = "\([a-z]*\)".*/\1/p'
sed -n '/^var taskStateRoster = /,/^}/p' internal/memory/taskstate.go \
  | grep -o 'statemachine\.State[A-Za-z]*' | sed 's/statemachine\.State//'
```

**输出（逐字）**：

```
legacy 7 = running done succeeded cancelled interrupted stuck error
d43 5    = Thinking Settling Muted Stuck Error
```

**12 枚（＝`taskStateRoster` 顺序＝生成 CHECK 的顺序）**：

```
 1 running    2 done        3 succeeded   4 cancelled
 5 interrupted 6 stuck      7 error       8 Thinking
 9 Settling  10 Muted      11 Stuck      12 Error
```

★**两处易错，我当场量到并纠正**（写在这里，让下一个人不重犯）：
1. `grep -o '"[a-z]*"'` 会把 `succeeded` 那行**行尾注释**里的 `"done"` 也算成一枚 ⇒ 量出 **13 枚**。
   第一次我就这么错了；改成 `sed` 只取赋值字面量后＝**12 枚**。（＝尺清单第 1 例"把注释算进证据"）
2. `stuck`/`Stuck`、`error`/`Error` **是两对大小写不同的独立取值**（`state TEXT` 无 `COLLATE`，SQLite 默认 BINARY 比较大小写敏感）
   ⇒ 12 枚是真 12 枚，⛔ 有人把大写那两枚当"笔误"合并掉。

## 3. 改前：契约里确实没有它（读数）

尺＝**在 `HEAD` 的同一个 `task_log` 定义块里**逐枚 `grep -o -w`（⛔ 全文——全文里 `error`/`done`/`Thinking` 会因别的语境命中，那是假读数）：

```sh
git show HEAD:docs/specs/SPEC-02-data-storage.md | sed -n '/^CREATE TABLE task_log (/,/^);/p'
git show HEAD:docs/PLAN.md | grep '^| `task_log` |'
```

| 枚 | SPEC-02 的 `task_log` 块 | PLAN D35 的 `task_log` 行 |
|---|---|---|
| running / done / succeeded / cancelled | **0** / **0** / **0** / **0** | **0** / **0** / **0** / **0** |
| interrupted / stuck / error | **1** / **0** / **0** | **0** / **0** / **0** |
| Thinking / Settling / Muted / Stuck / Error | **0** ×5 | **0** ×5 |

★**唯一的 1 处命中，逐字如下**（在**列注释**里，⛔ 值域声明）：

```
  state          TEXT NOT NULL,          -- 结束态（含 'interrupted'）
```

⇒ 严格讲是 **11/12 命中 0、1/12 只在注释里出现一次**。⛔ 写成"12 枚全 0"——那是把注释当证据。

**D35 那一行改前原文**（只列列名、只字不提值域）：

```
| `task_log` | `id` PK(TEXT), `started_at`, `ended_at`, `state`, `query_text`, `summary_text`, `cost_tokens_in`, `cost_tokens_out`, `cost_amount`, `error_class` | `started_at` | **30 天** | **L3 任务日志**（D20）。`query_text` 须按 §14.4 脱敏 |
```

## 4. 两处插入的逐字内容（两份文件**逐字节相同**，机检 `IDENTICAL`）

```
**`task_log.state` 的值域（票 196；本段与代码同源，⛔ 是第二份手工名册）**：合法取值＝**12 枚**，顺序＝`internal/memory/taskstate.go` 的 `taskStateRoster`（`schema.go` 的 `buildDDLV4` 按该切片顺序生成 CHECK；改一处须同改另一处）：
`running` · `done` · `succeeded` · `cancelled` · `interrupted` · `stuck` · `error`（既有词 7 枚；`running`＝行已打开、尚无终态；`succeeded` 已无生产者，只为旧行保留）＋ `Thinking` · `Settling` · `Muted` · `Stuck` · `Error`（D43 名 5 枚，任务行可合法持有；⛔ 全 20 枚——`Sleeping`/`Listening` 是球的状态）。
⚠ `stuck`/`Stuck` 与 `error`/`Error` 是**两对大小写不同的独立取值**（TEXT 默认 BINARY 比较），⛔ 当成笔误合并成一枚。**名册之外的写入被 CHECK 拒绝**（`CHECK constraint failed`）⇒ 遇到"库拒绝"先照这张名册核，⛔ 当 bug；该 CHECK 由 schema v4 的 `task_log` 重建迁移加上（SPEC-02 §3 的 v1 DDL 块本身不含 CHECK）。
```

**落点**：`SPEC-02` 插在 §3 主 ```sql 围栏**闭合行之后**（`:124` 之后、`provider_health` 那条之前）；`PLAN.md` 插在 D35 八行表**之后**、`**三条必须写死的规则**：`之前。
两处都是**纯插入**：各 `+4` 行（3 行正文 ＋ 1 行空行分隔）。

★**为什么插在围栏外、⛔ 围栏内**：`internal/memory/schema.go:10-18` 声明 `ddlV1` 必须与 `SPEC-02 §3` 的 ```sql 块**逐字节一致（含注释）**。
我是手量过的：两者今天**已有 2 行**不同，且差异全是 ban #8 逼出来的 ASCII 化（Go 里 `<=` / `->`，文档里 `≤` / `→`；`d22scan` 的 `ban #8` 只扫 `internal/`＋`cmd/`）。
⇒ 往围栏**内**加一行 `--` 注释会把这份差异从 2 行推成 3 行，而我又⛔ 改 `.go` 去配平；插在围栏**外**则围栏一个字节不动（§6 有复量）。

## 5. 交叉核：契约的 12 枚 vs 代码的 12 枚

```
docs/specs/SPEC-02-data-storage.md    contract=12   per-name same order = YES   set equal = YES
docs/PLAN.md                          contract=12   per-name same order = YES   set equal = YES
diff code-roster-12.txt contract-spec.txt  -> rc=0（无差异行）
diff code-roster-12.txt contract-plan.txt  -> rc=0（无差异行）
```

契约那侧是**机检**取的（先剥掉两处括号注释，再按出现顺序抓反引号词），⛔ 由我口报：
`running done succeeded cancelled interrupted stuck error Thinking Settling Muted Stuck Error`——与代码逐枚相同、顺序相同。

**改后全文命中**（每枚 ≥1，两文件都有）：`running 2/2 · done 1/21 · succeeded 2/2 · cancelled 2/3 · interrupted 2/1 · stuck 2/2 · error 3/4 · Thinking 1/11 · Settling 1/25 · Muted 1/5 · Stuck 2/14 · Error 2/25`（格式＝SPEC-02/PLAN）。

## 6. 原句一个字未动 ＋ 删除 0 行

```
$ git diff --numstat -- docs/PLAN.md docs/specs/SPEC-02-data-storage.md
4	0	docs/PLAN.md
4	0	docs/specs/SPEC-02-data-storage.md
```

- **删除列＝0 / 0**。
- `git diff -U0 | grep '^-' | grep -vc '^---'` ＝ **0** ⇒ diff 里**一条 `-` 行都没有**（唯一的 `-` 行是 `--- a/…` 文件头）⇒ 既有字节一个没动。
- **复量围栏未动**：重跑 `ddlV1` vs §3 围栏正文 → 差异**仍是那 2 行**（4 个 diff 行＝2 处增删对），与改前**同值**。
- `git status --porcelain -- docs` ＝ 只有 ` M docs/PLAN.md` ＋ ` M docs/specs/SPEC-02-data-storage.md`（无第三枚、无 D/A）。

## 7. 写面／未碰

- 本腿只动：`docs/specs/SPEC-02-data-storage.md`、`docs/PLAN.md`、本件。
- ⛔ 动：`internal/**`、`cmd/**`、`scripts/**`、`.github/**`、`docs/reports/**`、`frontend/**`、`design/**`、票面、台账、D 号/C 号。
- ⚠ 工作树里有**别人的未提交改动**：`internal/observe/logging.go`、`internal/proc/envfork.go`、`internal/proc/envfork_test.go`
  **不是我的**，我一个字节没碰、也**没有**进我的提交（我的 `add`/`commit` 是逐个文件字面点名）。

## 8. 门禁 rc（`unset GOTMPDIR` 后跑，⛔ 走管道取 rc）

| 门 | rc | 读数 |
|---|---|---|
| `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`；`ban #8` 扫了 internal/ 630 ＋ cmd/ 148 ＋ frontend/ 116 ＋ design/ 39，**⛔ 扫 docs/** |
| `sh scripts/check-path-length-budget.sh` | **1（预存，⛔ 本腿引起）** | `VERDICT RED`：过预算 **62** 条 vs 名册 **61** 条；唯一不在名册的那枚＝`.scratch/wisp/issues/365-two-taskLogState-functions-disagree-so-statusstuck-writes-stuck-on-one-path-and-error-on-the-other.md`（name 105 字符） |

★**pathlen 的 RED 归因（⛔ 我把它记成绿、也⛔ 记成"我的锅"）**：
1. 那把尺量的是**路径长度**，而本腿**一个路径都没新增**——`git diff --name-status` 只有两行 `M`，零 `A`/`D`；
2. 那枚犯规文件**在 `8beec66b1` 就进树了**（`git log --oneline -1 -- <那枚>` 逐字＝`8beec66b1 编排者A1166：…立票365…`），**早于我开工**；
3. 我这两枚文件相对路径 12 / 34 字符（预算 165），且内容改动⛔ 影响路径长度。
⇒ 结论：**RED 与本腿无关，本腿没有让它变红、也没有让它变绿**。⛔ 把 rc=1 说成 rc=0。

## 9. 新发现（⛔ 我自己调和，交编排者/机主）

1. ★**`state` 的 v1 列注释比名册窄**：`-- 结束态（含 'interrupted'）` 说这列是**结束态**，而名册里 `running` 是**行刚打开、尚无终态**的取值
   （`taskstate.go:59` `// row opened, no terminal state yet`；`journal.go:131` 确实插 `running`）。
   ★那行注释在**逐字节一致**的围栏里、我⛔ 改 ⇒ 我只写了值域，**没有去动它、也没有在这里替它圆场**。这是一处**待机主定**的措辞落差。
2. ★**票面括注与 `196-r1` 的读数不一致**：本腿派单写"D43 里**任务行真会取的** 5 枚"，而 `196-r1` 的 §9 欠账第 3 条量的是"**今天没有生产者写它们**、为迁移预留"。
   ⇒ 我按**代码的原话**写（`taskstate.go:35-42`："the five a task record **can legitimately hold**"）⇒ 契约里写"**可合法持有**"，**⛔ 写"今天确实写"**。
3. `schema.go:10-18` 说 `ddlV1` 与 SPEC-02 §3 ```sql 块"逐字节一致（含注释）"——**实测今天差 2 行**（`<=`/`≤`、`->`/`→`，ban #8 逼出的 ASCII 化）。
   ⇒ 那句注释本身**过期**；本腿没改它（⛔ 我的写面），只把它当证据用。

## 10. 采信前请注意（本腿自纠 3 次）

1. 第 1 次量代码：`grep -o '"[a-z]*"'` → **13 枚**（多一个 `done`，来自注释）⇒ 改 `sed` 取赋值字面量 → **12 枚**。
2. 第 1 次量契约：`split("（")[0]` 在第一个全角括号处切断 → 只剩 **7 枚** ⇒ 改成"先剥括号再抓词" → **12 枚**。
3. 第 1 次取 rc：把 `'-'` 写进了 PowerShell 单引号串 ⇒ 命令被截断、`git diff` 根本没跑（⛔ 那次读数）⇒ 改成脚本文件落盘再跑。
★三次都是**尺的问题、不是结论的问题**——所以本件每个数都给了**可逐字重跑**的命令。
