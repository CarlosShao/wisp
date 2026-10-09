# 282-v1 · AC#2／AC#4／AC#5 三格现量（非实现者，只读，零 `.go` 改动）

取样一律 `git show <rev>:<path>`／`git grep … HEAD`；工作树只在 AC#4(b) 的**第二列**用一次（那一格本来就要工作树字节）。
尺子里行尾一律 `tr -cd '\r' | wc -c`（⛔ 不用 `grep -c $'\r'`，它在本机给假 0）。

---

## AC#5 越界检查——名册结论：**越界 0 枚**

### 名册找回（不采信转述的三枚 sha，先自取）

- `git log --oneline --grep='242-corrland'` ⇒ 命中 5 行，其中**属于这条腿的产码/探针笔**＝`133b1bfa`／`bd124b2a`／`0abe217c`；另两行（`f5c63b7c`／`bbbc18fa`）是编排者的裁定与落账笔。
- `git log --oneline 133b1bfa~1..0abe217c` ＝ 6 枚，其中 `133b1bfa`／`bd124b2a`／`0abe217c` 是本腿，夹在中间的 `88bc61b5`／`0699f49f`／`52501e22` 属 `done-class-b-1` 与编排者补票，⛔ 不计入本腿射程。
- `git merge-base --is-ancestor 133b1bfa 0abe217c` ⇒ 真（区间方向没写反）。
- ⚠ **转述的三枚里只有 `bd124b2a` 碰产码**：`133b1bfa` 落 `.scratch/wisp/probes/242/corrland1/00-anchor.md`（1 件）、`0abe217c` 落 `.scratch/wisp/probes/242/corrland1/10-landing.md`（1 件）。⇒ 对这两枚而言"没碰禁区"是**平凡成立**（它们一件产码都没落），我在下面仍然按名册扫过 numstat，不省略。

### 逐枚 numstat（`git show --numstat`，`bd124b2a` 为唯一有产码笔的一枚）

| 枚 | 落地面（numstat 全量，含 `.scratch`） |
|---|---|
| `133b1bfa` | `.scratch/wisp/probes/242/corrland1/00-anchor.md`（仅此 1 件） |
| `bd124b2a` | `internal/agent/corr_percall_242_test.go` 145/0（新建）· `internal/agent/loop.go` 26/3 · `internal/tools/bridge.go` 7/3 · `internal/tools/cancel.go` 17/5 · `internal/tools/loop_approval_test.go` 10/2 · `internal/tools/subagent_197.go` 4/1 · `internal/tools/task.go` 7/5（共 7 件） |
| `0abe217c` | `.scratch/wisp/probes/242/corrland1/10-landing.md`（仅此 1 件） |

### 禁区扫描（三枚各扫一遍，pathspec 逐名列出）

命令形：`git show --numstat --format= <sha> -- docs/PLAN.md docs/specs internal/observe/thresholds.go frontend design internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go`
再一把：`git show --name-only --format= <sha> -- '*golden*' internal/agent/testdata/golden internal/llm/testdata/golden tools/mockllm/testdata/golden internal/llm/golden`

| 禁区名 | `133b1bfa` | `bd124b2a` | `0abe217c` |
|---|---|---|---|
| `docs/PLAN.md` | 0 | 0 | 0 |
| `docs/specs/**` | 0 | 0 | 0 |
| `internal/observe/thresholds.go` | 0 | 0 | 0 |
| `frontend/**` | 0 | 0 | 0 |
| `design/**` | 0 | 0 | 0 |
| `internal/panel/tokens_fourway_test.go` | 0 | 0 | 0 |
| `internal/panel/l2_grant_boundary_test.go` | 0 | 0 | 0 |
| `internal/perm/ticket90_persist_test.go` | 0 | 0 | 0 |
| golden 面（`*golden*` 两把我自取的名册：`internal/agent/testdata/golden` 11 枚 `.sse`／`internal/llm/testdata/golden` 40 枚 `.sse`／`tools/mockllm/testdata/golden` 2 枚／`internal/llm/golden/` 3 枚／`internal/agent/loop_golden_test.go`／`internal/llm/openaichat/harness_golden_test.go`／`tools/mockllm/goldenfmt.go`） | 0 | 0 | 0 |

⇒ **碰了的：无。没碰的：上面 9 行 × 3 枚＝27 次 numstat／name-only 扫描全部零命中**（不是"看起来没越界"，是对着具名 pathspec 扫的）。
⇒ 顺带一条对编排者有用的事实：`bbbc18fa` 的裁定边界①（"不动契约文本"）在盘上同样成立——`bd124b2a` 的 diff 里引用了 `docs/PLAN.md:1368`，但**只引用、未改**（`git show --numstat bd124b2a -- docs/PLAN.md docs/specs` 输出空）。

---

## AC#2 有没有顺手放宽——判定：**为变绿而放宽 0 处；"三形皆红"成立；但格式轴有 1 处真实的松面（具名，见 S-4）**

### 这一轮它一共只碰了两枚断言载体

- 改既存断言句：**1 枚**＝`internal/tools/loop_approval_test.go`。
- 新增断言句：**1 枚**＝`internal/agent/corr_percall_242_test.go`（145 行全新文件）。
- ⛔ 其余 5 件（`loop.go`／`bridge.go`／`cancel.go`／`subagent_197.go`／`task.go`）是产码，不含断言；它们的"改变别的用例通过条件"的面向在 S-5／S-6 单独判。

### 行号与短语同一次取数（票面 `:215`、编排者转述 `:213`–`:215`，HEAD 实测已漂）

`git show HEAD:internal/tools/loop_approval_test.go | grep -n` ⇒

```
219:	if r.CorrelationID == "" || r.CorrelationID == res.TaskID ||
220:		!strings.HasPrefix(r.CorrelationID, res.TaskID) {
222:			r.CorrelationID, res.TaskID)
```
`git show bd124b2a:…` 同一把尺 ⇒ 行号 **219／220／222 一致**（这枚文件自 `bd124b2a` 起没人再动过，`git log --oneline -- internal/tools/loop_approval_test.go` 里最后两笔是别票早于它）。⇒ 漂的是**票面与转述的行号**，不是盘。

### 两半原文（`git show bd124b2a -- internal/tools/loop_approval_test.go`，−/＋逐字）

- 旧句（−，2 行）：
  - `if r.CorrelationID == "" || r.CorrelationID != res.TaskID {`
  - `t.Errorf("correlation_id = %q, want the task id %q (C18)", r.CorrelationID, res.TaskID)`
- 新句（＋，10 行，含 6 行注释）：
  - `if r.CorrelationID == "" || r.CorrelationID == res.TaskID ||` / `!strings.HasPrefix(r.CorrelationID, res.TaskID) {`
  - `t.Errorf("correlation_id = %q, want a non-empty per-call id routing back to task %q (C18)", …)`
- 取数顺带核到的编译面：`strings` 这枚 import **不是新加的**——`git show bd124b2a~1:internal/tools/loop_approval_test.go` 第 13 行已有 `"strings"` ⇒ 新句没有藏"为了让 diff 好看而动 import"的侧账。

### 逐形判定（把"应被判红的形状"一枚一枚喂给新旧两句，纯布尔，不跑测试）

| 形 | 旧句 | 新句 | 判定 | 凭据（新句里能指的那一枚子句） |
|---|---|---|---|---|
| S-1 空串 `""` | 红 | 红 | **等价** | `r.CorrelationID == ""`（`:219` 第一支） |
| S-2 旧形等值 `corr == taskID` | **绿**（这正是旧句要求的值） | **红** | **更强** | `r.CorrelationID == res.TaskID`（`:219` 第二支）＝旧句的**通过集**被新句判红 |
| S-3 无前缀形（与本任务无关的 corr） | 红 | 红 | **等价** | `!strings.HasPrefix(r.CorrelationID, res.TaskID)`（`:220`） |
| S-4 前缀＋任意后缀（`taskID+"#anything"`） | 红 | **绿** | **更松（唯一一处松面）** | 新句只钉"前缀回本任务"，⛔ **不钉分隔符与后缀形状**；旧句的通过集只有 `{taskID}` 一枚，新句的通过集是无穷集 |
| S-5 `taskID` 前缀但整体相等（等价于 S-2） | — | 红 | 更强 | 同 S-2 |

⇒ 票面 AC#2 要的**三形皆红（空／旧形等值／无前缀）＝成立**（S-1／S-2／S-3 三行）。
⇒ 编排者转述里那句"比旧句更严"**只对了一半**：对 S-2 是新增拦截（更严），对 S-4 是把"钉死唯一值"换成"钉一族前缀匹配"（更松）。**这一处松面不是"为变绿而放宽"**——旧句在 242 之后本来就**必须**死（loop 现在铸的是 `taskID#callID`，旧句会把正确行为判红），裁的方向就是把等值钉换成路由钉；但**松的那一根轴是真实存在的，且今天没有任何别处补上它**：
  - 尺面现量：`git grep -n 'callCorr|"#"|#call-' HEAD -- '*.go' ':(exclude).scratch'` ⇒ 全仓对 `callCorr` 的**形状**没有任何断言；`internal/agent/corr_percall_242_test.go:139` 用的也是 `strings.HasPrefix` 同一族（不是等值），`:142` 只钉 `c1 != taskID`。
  - 于是 S-4 这一族里能藏住的具名形状：**id-less 调用回落支** `fmt.Sprintf("%s#call-%d", taskID, index)`（`internal/agent/loop.go:603-607`）——轮内 index 在**跨轮**重复时（第 2 轮的第 0 枚＝`taskID#call-0`，与第 1 轮第 0 枚同值）两枚调用会**共用一枚路由键**，而 `loop_approval_test.go:219-220` 与新用例（其 provider 给了 `call_p1`/`call_p2`，走的是非空支）**都不红**。⇒ 这条不属于"它放宽了既存钉子"，属于"新钉子的射程比旧钉子的等值窄在形状轴上"；我把它的**读数欠账挂在 AC#1**（见 `20-held-cells.md`），本腿不判红绿、不种坏。

### 反向面（它把东西变强了 ⇒ 谁的通过条件被改了）——具名

| 位 | 变化 | 会被它变红的形状 | 今天盘上有没有这种形状 |
|---|---|---|---|
| S-5 `internal/tools/subagent_197.go:262` `parentID := CorrelationID(ctx)` → `TaskID(ctx)`（票面/转述写作 `internal/agent/subagent_197.go`，**那枚路径在 HEAD 不存在**，见文末顶回①） | 身份来源从"corr（`bridge.go:269-271` 还带 `corr←taskID` 回落）"换成"只认 `req.TaskID`，⛔ **没有反向回落**（`bridge.go:528` `taskID: req.TaskID` 原样带）" | 只填 `CorrelationID` 不填 `TaskID` 的 `task.spawn` 派发 ⇒ 今天会走 `:263` 那句 fail-closed 文案 | **无**：`git grep -n -A3 '"task.spawn"\|"task.cancel"' HEAD -- '*_test.go' ':(exclude).scratch'` ⇒ 命中的构造点全部同时给 `TaskID`（`subagent_197_test.go:217`／`subagent_carrier_197_test.go:173`／`subagent_selfapproval_197_test.go:116`／`task_cancel_221_legs_test.go:119`／`failclosed_236_teeth_test.go:145` 与 `:288`（两者皆空＝本就是 fail-closed 用例）／`subagent_197_test.go:655` 只给 `TaskID`，改前后塌成同一值（`orDefault`）⇒ 不翻） |
| S-6 `internal/tools/task.go:708` `caller := CorrelationID(ctx)` → `TaskID(ctx)`（注释同批改写） | 同上 | 同上形状的 `task.cancel` 派发 | **无**（同一把尺扫过；`internal/tools/task_cancel_221_legs_test.go:380-381` 那枚 `r.TaskID != parent197` 钉的是桥记的 tool_call 行，而行里的 `TaskID` 取自 **`bridge.go:1116 TaskID: req.TaskID`**，本腿未碰它；同处 **`:1124 CorrelationID: orDefault(req.CorrelationID, req.TaskID)`** ⇒ 生产组合下 tool_call 行的 correlation 就是 per-call corr，这条正好支撑 AC#4(a) 的"生产侧没问题"） |
| S-7 产码 `ToolRequest` 构造点全仓唯枚 | `git grep -n 'ToolRequest{' HEAD -- '*.go' ':(exclude)*_test.go' ':(exclude).scratch'` ⇒ **只命中 `internal/agent/loop.go:675`** | — | ⇒ S-5/S-6 的收紧面**碰不到任何生产派发**（唯一生产构造点同时给两枚 id），只能碰测试夹具，而夹具名册见上 |
| S-8 `internal/agent/loop.go:363→369` 删掉 `// C18: correlation == task id` 换成 6 行"C18 从不要求同值" | 注释，非断言 | — | **核实为实话**：`git show HEAD:docs/PLAN.md` 第 1368 行 C18 原文＝「全局 FIFO；每项含 `correlationId` / 所属任务标识 / 工具名 / …」——两枚字段**并列列出、未写相等**；路由义务在 C17（`:1367`「回复必须按 correlationId 路由」）与 C27（`:1377`「按 correlationId 分区渲染」）。⇒ 旧注释是一句 C18 没做过的许诺，新注释更贴原文；契约文件本身零碰（AC#5 已扫） |

### 新用例自身有没有"看着像尺其实是软尺"

`git show HEAD:internal/agent/corr_percall_242_test.go`（145 行）：
- `grep -n 't.Skip\|testing.Short\|//go:build'` ⇒ **0 命中**（无跳、无 tag）。
- 断言全是 `t.Fatal`/`t.Fatalf`/`t.Errorf`，**没有一枚降级成 `t.Logf`**：`:96` `t.Fatal(...)`／`:101` `t.Fatalf`／`:104` `t.Fatalf`（同键塌缩）／`:113` `:116` `:120` `:129` `:132` `:136` `t.Fatalf`／`:140` `:143` `t.Errorf`。
- `:120` 那枚 `select … case <-done: t.Fatal("the task completed although the second ask was never answered…")` 是**反向钉**（不许提前完成），属加强不属放宽。
- ⚠ 唯一"软"的地方已在 S-4 记过：`:139` 用前缀不用等值 ⇒ 整条新用例也不钉 corr 的**形状**。

---

## AC#4 两笔残余逐笔现量（三形之一＋尺）

### (a)「journal 非生产组合行仍写 taskID」⇒ **真存在**（HEAD 可指），且"非生产组合"这一定性也成立

尺（全部 `git show HEAD:<path>`）：
1. `git show HEAD:internal/agent/loop.go | grep -n newTaskJournal` ⇒ **`:369: j := newTaskJournal(l.opt.Journal, taskID, taskID)`**（第二枚实参仍是 `taskID`）。
2. 那枚第二实参真的会落进 `tool_call` 行的 correlation 列：`git show HEAD:internal/agent/journal.go` ⇒ `:59 func newTaskJournal(j Journal, taskID, corrID string)`，`startCall` 里 **`:81 CorrelationID: t.corrID`**。⇒ "行里记任务级"不是推测，是写路径可指。
3. "生产组合"这一侧：`git grep -n 'Journal:' HEAD -- '*.go' ':(exclude).scratch'` ⇒ 生产只有两枚且**分属两层**——`cmd/wisp/run.go:750 Journal: mem` 在 `assembleRuntime`（382–836）里、属于 **`tools.New(tools.Options{…})`＝桥的 journal**；`cmd/wisp/run.go:1003 Journal: nil` 在 `(rt *agentRuntime).execute`（990–）里、才是 **`agent.Options`＝loop 的 journal**（同处注释逐字：「The loop's own tool_call rows are deliberately NOT booked: the bridge already writes the authoritative one」）。⇒ **残余只在"有人给 loop 直接接 Journal"的非生产组合**（盘上名册：`internal/agent/declared_l0_risk_179_test.go:204/227/248` 三枚 `newTaskJournal(store, task, task)` 是这条组合的现役消费者）。
4. 生产那一侧的行里到底是什么：`git show HEAD:internal/tools/bridge.go` ⇒ **`:1124 CorrelationID: orDefault(req.CorrelationID, req.TaskID)`**（桥记的 tool_call 行吃 `req.CorrelationID`＝per-call corr），而 `req.CorrelationID` 由 `loop.go:676 callCorr(taskID, p.call.ID, i)` 供给 ⇒ "生产组合下行里是 per-call corr"这半句**在盘上可指**，不是它的自述。
5. 登记面：`git grep -n 'DEFERRED(' HEAD -- 'internal/agent/*.go' 'internal/tools/*.go' ':(exclude).scratch'` ⇒ 10 枚命中里**没有一枚**是这条残余（`loop.go:407` 那枚是 `D28-1` Warm 窗钩子，与 journal 无关）。⇒ 它只活在探针 md（`.scratch/wisp/probes/242/corrland1/10-landing.md` §6 第 2 条）里，**没走 AGENTS §1.1 的 DEFERRED↔`SPEC-12 §5` 1:1 登记**。⚠ 这一条请编排者裁：不是"它放宽了产码"，是"它的推迟没登记"。
6. 归票：票 283（`283-corr-landing-named-two-new-blind-spots…-done.md`）只裁两枚身份盲区补尺，不涉 journal；`284`／`285` 名册不含 journal。**没有现役票认领它** ⇒ 我的建议形：判**留**，账挂到"下一枚真动 `internal/agent/journal.go` 或给 loop 接非 nil Journal 的票"；若编排者要它变成显式债，就得补 `DEFERRED(...)` 标记＋`SPEC-12 §5` 一行（那是产码改动，⛔ 本腿不做）。

### (b)「`internal/tools/bridge.go` 既有 CRLF 没洗」⇒ **从来就不是那样**（作为该文件的债）；工作树那一侧确实带 CRLF，但那是检出态

尺（**两列并排**，行尾只按 `tr -cd '\r' | wc -c`）：

| 取数 | CR 字节数 |
|---|---|
| 工作树 `internal/tools/bridge.go`：`tr -cd '\r' < internal/tools/bridge.go \| wc -c` | **1284** |
| `git show HEAD:internal/tools/bridge.go \| tr -cd '\r' \| wc -c` | **0** |
| `git show bd124b2a:internal/tools/bridge.go \| tr -cd '\r' \| wc -c` | **0** |
| `git show bd124b2a~1:internal/tools/bridge.go \| tr -cd '\r' \| wc -c` | **0**（⇒ "既有"在落地侧**从来没有过**，不是它没洗） |

旁证（同一把尺的四枚）：
- `wc -l < internal/tools/bridge.go` ＝ **1284** ＝ CR 字节数 ⇒ 全文件每一行都带 CR，是**整文件检出态**，不是局部脏行。
- `git ls-files --eol -- internal/tools/bridge.go` ＝ **`i/lf    w/crlf  attr/text eol=lf`**（index 侧 LF／工作树侧 CRLF／属性写着 eol=lf）。同尺扫另 6 件本腿产物：`corr_percall_242_test.go`／`loop.go`／`cancel.go`／`loop_approval_test.go`／`subagent_197.go`／`task.go` **全部 `i/lf w/lf`** ⇒ 它自报的"其余 6 枚均 LF"在 index 侧与检出侧都对。
- `git status --porcelain -- internal/tools/bridge.go` ⇒ **空**（工作树与 index 一致 ⇒ 那 1284 枚 CR 不是谁在飞的改动，是 checkout 产物）。
- `git config --get core.autocrlf` ＝ **`true`**；`git show HEAD:.gitattributes` 第 **4** 行＝`*.go text eol=lf`。
- **归口早已有票且已结**：票 280 名册第 9 行逐字把 `bridge.go` 列在"其余 5 枚行尾 artifact"里，给的现量与工作树 CR 枚数与我**完全对上**（`334/131/225/1115/1284` 的最后一枚＝1284），判语＝「**是尺的射程错，不是文件的债**」，AC#3 处置＝「5 枚行尾 artifact＝**不改**」，`Status: done`（2026-10-09 10:3x）。⇒ 这一笔**不是悬账**，本腿判"从来就不是那样"并把它指向票 280。

---

## 顶回编排者／票面的前提（以盘上原文为准，⛔ 不迁就转述）

1. **`internal/agent/subagent_197.go` 不存在**。`git cat-file -e HEAD:internal/agent/subagent_197.go` ⇒ `does not exist in 'HEAD'`；`git ls-tree -r --name-only HEAD | grep subagent_197` 的产码枚＝**`internal/tools/subagent_197.go`**（numstat 也是这枚）。票面 `:9` 与本腿派单转述都写作 `internal/agent/`。台账 `A749` 已具名同一处出入（并注"hash 是真的"），本腿独立复现。⚠ 票面一字未动，改名归编排者。
2. **转述的三枚 sha 不是三枚产码笔**：只有 `bd124b2a` 落产码（7 件），`133b1bfa`／`0abe217c` 各落 1 枚 `.scratch` 探针 md。AC#5 对后两枚是平凡成立，我在表里仍逐枚给了扫描名册。
3. **行号已漂**：票面 `:215`、派单转述 `:213`–`:215` ⇒ HEAD 实测 **`:219`–`:222`**（行号与短语同一次 `git show` 取，见上）。
4. **派单给 AC#4(b) 的"既有 CRLF 没洗"这一前提，在落地侧不成立**（三枚 rev CR＝0、票 280 已判"非债"）。⇒ 那条腿不是"漏洗了债"，是**自报过头**（它把检出态当成文件债写进残余）。这条与"它没放宽任何东西"不矛盾，但它属"自报账目与盘不符"一族，值得进台账。
5. **"比旧句更严"这句判语只对一半**（S-2 更严／**S-4 更松**，见 AC#2 表）。⇒ 编排者转述若写成"全轴更严"，那是一处需要缩的表述。
