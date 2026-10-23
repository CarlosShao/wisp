# 票 164 终裁腿 `164-v1b` — AC#4「取消要真取消」能不能翻

- 腿：**非实现者终裁腿 `164-v1b`**（裁决者≠实现者）；工作树 `D:\work\workspace\projects plans\Wisp`；分支 `dev`。
- **锚漂移照实记**：起手 `HEAD=e9caddff`，量到一半被别人的 docs-only 提交推到 `7825d31e`，收尾时 `21b528d6`（编排者 A1159）。`git diff --stat 81388d1f HEAD -- <被裁 3 件>` **空** ⇒ 被裁字节没被我以外的人动过；`git diff --name-status 7825d31e 21b528d6 -- internal/tools cmd` **空**。
- 被裁：`821302c9`（仪器 M ＋ 新测试件 A）／`efa950a3`（交付件 A）；件 `docs/evidence/s1/164-r2-impl.md` 134 行。
- 尺的纪律：每条读数下面写明**量的是哪个对象**；跑 go 前一律 `unset GOTMPDIR TMPDIR TEMP TMP; GOFLAGS=`（`go env GOTMPDIR` 现跑＝空）；`git hash-object`＝SHA-1/40 位、`sha256sum`＝SHA-256/64 位，**分开用、不混比**；格式归属判 `[\\/]`；CR 判 `tr -cd "\r" | wc -c`。

## ① 票面前提新鲜度（对 HEAD blob 现量；三条逐条复核）

| # | 它报的腐坏 | 我自己的读数 | 判 |
|---|---|---|---|
| 1 | 定案⑤「`RunAsync` 生产零调用点」不成立 | `grep -rn 'RunAsync' --include=*.go internal cmd tools \| grep -v _test.go \| grep -vE ':[0-9]+:[[:space:]]*//'` ⇒ **2 处真调用**：`cmd/wisp/run.go:1176`、`internal/tools/subagent_197.go:372`；另 1 处定义 `internal/agent/loop.go:328`（原始命中 8 行含 5 行注释，**代码/注释分开数**） | ★成立，但**它报的"1 枚"少数了一枚**：它自己修好的仪器 `[6]` 第 91 行那条 `grep -rn "RunAsync"` 在 HEAD 就会印出 `subagent_197.go:372`；票面 `:43` 也早已写"两处真调用点" |
| 2 | 「名册零写者」不成立 | 三把尺：`.Record(` 非测试**全类型 8 处**；`TaskRoster.Record` 直接调用**4 处**（`subagent_197.go:462`／`task.go:343`／`task.go:356`／`task_backfill.go:141`）；**进名册的生产写点 5 处**（外部对 `Record`/`PublishSubagent`/`MarkRoot` 的调用：`subagent_197.go:338`／`:373`／`:462`、`task_backfill.go:141`、`run.go:1182`） | ★结论成立；**它报的"5 处"只在第三把尺下成立**，它点名的尺（`.Record(`）给 4 或 8 |
| 3 | 「生产注册表只有 `BuiltinFSEntries`」只在锚点成立 | `cmd/wisp/run.go` 现注册 **4 族**：`:528 BuiltinFSEntries`／`:557 BuiltinTaskEntries`／`:799 BuiltinSubagentEntries`／`:838 BuiltinPlanEntries` | ★成立（仪器 `[3]` 量的是锚点 `f206e9f`，不是今天） |

- 行号也有漂：它引「票面 `:67`」＝今天的 **`:75`**（定案⑤）；「`:84`」那句今天在 **`:57`／`:92-93`**，`:84` 不是它；「`:26`」**对**。
- **对 AC#4 判据的影响＝0 条负面的**：三条腐坏只**退掉**定案⑤那句"今天连被测对象都没有"，判据原文一字不动。⇒ AC#4 从"〔今天无法判定〕"变成"可判定"，而 r2 正是按定案⑤要求的顺序做的：**先量原样读数（绿）、再立变异**。**但归因要说清**：对象是票 176/197 落进来的，r2 交的是**判据**，不是守卫。

## ② 核心判点一：这一跳有没有被钉住（我自己造的两发）

**V1（取消只说不做）**：`internal/tools/task.go:462` `cancel()` ⇒ `_ = cancel`。落地：`grep -n 'MUT164v1b-V1T'` ⇒ `462:	_ = cancel // MUT164v1b-V1T`；`go build ./...` **rc=0**；`go test -count=1 -v -run '^Test164AC4CancelLeavesNoWritingTail$' ./internal/tools/` **rc=1**，逐字：
```
    task_164_ac4_cancel_tail_test.go:107: 取消之后那股输出还是被写进名册了：Backfill 没有拒绝——取消之后不许有还在写的尾巴
    task_164_ac4_cancel_tail_test.go:114: 取消后名册那一行被改动了：Text "被停掉前的半句正文" -> "取消之后尾巴还在写的正文"，State "Muted" -> ""（取消之后落到这一行的写就是那条尾巴）
    task_164_ac4_cancel_tail_test.go:125: Cancel 返回成功，但它挂上的句柄没有真停掉子任务：ctx.Err() == nil（取消只是把话说出去了）
--- FAIL: Test164AC4CancelLeavesNoWritingTail (0.00s)
```
**V2（摘掉/绕开写入口判据；我故意换一种写法，⛔ 照抄它的 `&& false`）**：`task_backfill.go:123` `Stopped(ctx)` ⇒ `Stopped(context.Background())`。落地 `123:	if stopped, why := Stopped(context.Background()); stopped {`；build **rc=0**；rc=**1**，逐字：
```
    task_164_ac4_cancel_tail_test.go:107: 取消之后那股输出还是被写进名册了：Backfill 没有拒绝——取消之后不许有还在写的尾巴
    task_164_ac4_cancel_tail_test.go:114: 取消后名册那一行被改动了：Text "被停掉前的半句正文" -> "取消之后尾巴还在写的正文"，State "Muted" -> ""（取消之后落到这一行的写就是那条尾巴）
--- FAIL: Test164AC4CancelLeavesNoWritingTail (0.00s)
```
- **两发可区分**：V1 红 **3 行**（`:107`／`:114`／`:125`）、V2 红 **2 行**（`:125` 仍绿）。区分点唯一＝`:125`（"句柄真跑过"），`:107`／`:114` 是两发共有的"尾巴落地"事实。**判语：可区分，但判据本身只把"取消是假的"与"门是聋的"分开在一条诊断行上**——够用，别把它当成两条独立性质的证。
- **整包尺度要说清**（r2 件里没写它的变异命令）：V1 在 `./internal/tools/` 整包跑**不是红行，是 10 分钟超时**（`panic: test timed out after 10m0s`、卡在 `Test197CancelIsPerRowAndNeverCascades`、`RUN=242/PASS=158/FAIL=0`、new 用例连 `=== RUN` 都没到）；我另跑 `-run '^Test221'` ⇒ `Test221ParentStopsItsOwnChildRowAndStreamSettle` **也红**（命名红），`-run '^TestNoRosterWriteAfterCancel176r1$'` 在 V1 下**仍绿**。V2 整包 ⇒ 红 2 枚：new 用例 ＋ `TestNoRosterWriteAfterCancel176r1`。
- ⇒ **两条缝各自另有钉子**：停法有 221、写入口有 176-r1。**这枚新用例的独有价值＝顺序**（取消 → 经真 `finalize` 落终态 → 再拿子任务自己的 ctx 写尾巴）＋"那一行 `Text`/`State` 一个字不许动"＋`:125` 在尾巴语境里被重说一遍。这与它"拼两条真缝"的说法**一致**。
- **还原**（每发读完立刻单独还原，先问算法）：`sha256sum` 回 `055fad4b…356fbc0`（task.go）／`b0b5a3ab…dd41bec`（task_backfill.go）；`diff` 备份对盘上 `IDENTICAL`；`git status --porcelain -- <file>` 空；`grep -c 'MUT164v1b'`＝**0**；全树 `MUTATION 164-r2`／`MUT164v1b` 命中 **0 文件**。

## ③ 核心判点二：今天绿 ≠ 没有缺陷？——"今天到底还缺什么"

**今天是真的**：`TaskRoster.Cancel` 跑句柄（`task.go:447-465`）⇒ 子任务 ctx 真死；唯一的输出写入口 `TaskBackfill.Backfill` 在碰表之前读一次 `Stopped(ctx)`（`task_backfill.go:123`，`Stopped`＝`cancel.go:86`）⇒ 取消后回填被拒且理由说"取消"。**我找不到一条"今天可达且拦不住"的形**：settle 之后名册的写者今天只剩（a）`finalize` 自己那次落终态（176-r1 的 L-4 明写排除在判据射程外）与（b）被守卫的 `Backfill`；`Record` 本身对取消**故意**无感知（`task_backfill.go:19-24` 写明了），但今天没有第二个生产者在 settle 之后写它（两个 publisher 都在孩子第一次模型调用之前）。

**但 AC#4 今天仍缺三件（都不构成本格退回，具名给下一腿）**：
1. **生产里没有"既能取消、又走这道写入口"的那枚对象**。`run.go:1182` 只 `MarkRoot`，**全仓 `AttachCancel` 只有 `subagent_197.go:374` 一处**（子代理那支），而子代理落终态走 `finalize`、不走 `Backfill`。⇒ `task.cancel` 打 `wisp run` 那枚根行会走 `task.go:459`「没有可停的句柄」——**这一支今天 0 枚测试**。接一行 `AttachCancel` 就补上，而且闸门已经现成（这正是"钉门不钉调用者"买到的）。
2. **调用者侧没有钉子（V3）**：把唯一生产调用点 `run.go:1187` 的 `bg.Root().Ctx` 换成 `context.Background()`，`go build ./...` rc=0、`./internal/tools/` 整包 **rc=0／330/247/0/0 全绿** ⇒ "Backfill 的 ctx 必须是该任务自己的 ctx"今天只是 `task_backfill.go:80-82` 的注释契约，没有测试。`cmd/wisp` 那侧也指望不上（见 ⑥：整包测试 11 分钟被杀、RUN=0，取不到读数）。
3. **产物那一半没有钉子（V4）**：在守卫**上面**插一发 `b.Spills.Prepare(...)` ⇒ 取消后仍往 artifacts 写出 20000 字节的产物，而**两枚钉子全绿**（new 用例 PASS ＋ `TestNoRosterWriteAfterCancel176r1` PASS，rc=0）——因为这枚新用例建的是 `TaskBackfill{Roster: roster}`（**`Spills` 为 nil**），176-r1 那枚只断言"行没进表、理由不是空"，**从不断言盘上没多出字节**。
- ⇒ **判：V3／V4 是"射程边界"，不是"今天的缺陷"**（今天两处行为都对）。要我把 AC#4 判成 ⛔ 的那条形**不存在于今天的树上**；存在的是"判据的两端之外还有两处没人管"，而 AC#4 的原文要求的正是**记录那一行**的尾巴。

## ④ 反向判据（自造、⛔ 只看它的 M1／M2）
V1／V2 见 ②（各先证落地再读数、读完立刻单独还原）。V3／V4 见 ③（各 `grep -n` 证落地、build rc=0、读数后立即还原；`run.go` 回 `58e8d8d6…ab3bfb`、`task_backfill.go` 回 `b0b5a3ab…dd41bec`，`git status` 空、marker 0）。**四发的落点都单独报、没有一次"顺手一起还原"**。

## ⑤ 零产码复核（blob SHA-1 逐枚，⛔ 与 SHA-256 混比）
`for p in $(git ls-tree -r --name-only HEAD -- internal/tools | grep '\.go$'); do git rev-parse HEAD:$p vs git hash-object $p; done` ⇒ **SAME=95／DIFF=0／MISSING=0**；同尺 `cmd/` ⇒ **SAME=149／DIFF=0**。新 `.go` 是 `internal/tools/task_164_ac4_cancel_tail_test.go`（`package tools`、139 行、`_test.go`）⇒ **产码零字节**。四枚相关件 HEAD blob 对工作树：`task.go c0abe997…`／`task_backfill.go 4bfdeb4f…`／`task_164… cac2d7f4…`／`task_output_ac2_before_test.go f409304b…` 全 SAME。

## ⑥ 门禁（四数分开报、⛔ 相减；变异读数与整包串行）
| 门 | 命令 | 基线 `7825d31e` | 收尾 `21b528d6` | 复跑 `21b528d6` |
|---|---|---|---|---|
| build | `GOFLAGS= go build ./...` | **rc=0**（无输出） | **rc=0** | — |
| vet | `go vet ./internal/tools/` | **rc=0**（无输出） | **rc=0** | — |
| 整包 | `go test -count=1 -v ./internal/tools/` | **rc=0**　RUN=330／PASS=247／FAIL=0／SKIP=0 | **rc=1**　RUN=330／PASS=246／FAIL=1／SKIP=0 | **rc=0**　330／247／0／0 |
| 去新用例 | 同上 `-skip '^Test164AC4CancelLeavesNoWritingTail$'` | **rc=0**　RUN=329／PASS=246／FAIL=0／SKIP=0 | **rc=0**　329／246／0／0 | — |
| d22scan | `sh scripts/d22scan.sh` | **rc=0** | **rc=0** | — |
| path-length | `sh scripts/check-path-length-budget.sh` | **rc=0** | **rc=0** | — |
- 逐字末行：d22scan＝`d22scan: clean - no D22 ban violations; live scope work: … ban #8 cmd/ 148 …`；path-length＝`check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`（分母 `tracked paths=17602`→`17604`，over=61、`in the wall interval=0`，两次都是 GREEN）。
- **收尾那枚红必须点名**：`Test222SpawnConclusionArrivesThroughRealBridgeChildren`（`subagent_222_test.go:549`「第 4 枚孩子的工具调用起跑时桥位占用 1 枚, want 4」）。**不是 AC#4 家、不是我变异留下的**：单跑 **5/5 绿**（rc=0×5），整包随即复跑 **rc=0／330/247/0/0**，基线那次也绿 ⇒ 判**负载型 flake**（4 个孩子要同时压在 4 枚桥上，CPU 挤时第 4 枚到点时只数到 1 枚）。⇒ **新增红 0**，但"这一轮整包曾红 1 枚"照实记，不许被四数盖住。
- **格式甲乙两尺都报**：**甲（blob 树、判据）**`git archive 821302c9 | tar -x -C <仓外>` ＋ `gofumpt -l .`（`v0.12.0 (go1.27.1)`）⇒ raw **20 行**、`tr '\\' '/'` 后 `(^|/)\.scratch/` 命中 **20**、**非 `.scratch`＝0**；新测试件**不在**名单（已格式化）。**乙（工作树、只作诊断）**`gofumpt -l .` ⇒ 54 行（`.scratch` 51 ／ 非 `.scratch` **3**：`cmd/wisp/models.go` CR=389、`internal/agent/approval/pending_read.go` CR=131、`internal/agent/tools.go` CR=225，三者 HEAD blob 的 CR 都是 **0**）⇒ 纯工作树 CRLF 假象，rc=2 来自 `.scratch` 里故意写坏的探针件。

## ⑦ 新增红 0（我的核法：名册两向差集，⛔ 只对总数）
- 去新用例名册 **246** 枚 ⊂ 整包名册 **247** 枚：`comm -13`（只在新）⇒ **`Test164AC4CancelLeavesNoWritingTail`**；`comm -23`（被吞）⇒ **0 枚**。基线与复跑两份整包名册 `comm -3` ⇒ **空**（同为 247 枚）。⇒ 新用例就是那个 ＋1，**没有别的名字动**；红名册（主检出）逐名＝**（无）**。
- **`cmd/wisp` 整包四数（⛔ 当判据）**：`go test -count=1 -v ./cmd/wisp/` ⇒ **rc=1**、末行 `*** Test killed: ran too long (11m0s).` / `FAIL	github.com/CarlosShao/wisp/cmd/wisp	33398.976s`、**RUN=0／PASS=0／FAIL=0／SKIP=0**、`0xc0000135` 命中 **0** —— 这是"截断点后整族 RUN=0"那形，**只作诊断、不算门禁**（本程没跑导出树，绝不用它的绝对数横比）。

## ⑧ 写面／越界／正控
- `git show --name-status 821302c9` ⇒ 2 件：`M .scratch/wisp/probes/164/v1/ac2-anchor-reading.sh`、`A internal/tools/task_164_ac4_cancel_tail_test.go`（`.go`＝1、**产码 .go＝0**）；`efa950a3` ⇒ 1 件 `A docs/evidence/s1/164-r2-impl.md`。两笔"禁区/越界"命中 **0**（`docs/specs`／`internal/risk|panel`／`thresholds.go`／golden／`allowlist.txt`／`frontend`／`design`／`scripts`／`.github`）。.scratch 那件是**本票自己的归档仪器**（票面 AC#2b ⑨＋编排者裁定"改仪器 grep"），⛔ 越界。
- **正控（同尺打真提交）**：`bd124b2a` ⇒ `.go` **7 枚／产码 5 枚**；`a8eb0d51` ⇒ **2 枚／产码 1 枚** ⇒ 这把尺看得见生产 `.go`（不然上面那个"产码 0"就是瞎的）。

## ⑨ 仪器 `[6]` 那条要不要在本票里修
- **要修，但不归本腿**。现量：`ac2-anchor-reading.sh:88` 标签写「对 `TaskRoster.Record` 的调用」，`:89` 的尺是 `grep -rn "\.Record(" --include=*.go internal cmd | grep -v _test.go` ⇒ **8 行**，其中**只有 4 行**是名册（另 4 行是 `approval/gate.go:684` 的 `grants.Record`、`approval_reply.go:139` 的 `h.Record`、`resident_approval_windows.go:262/876` 的 `target.Record`/`cards.Record`）。它没有 `[5]` 那种"无条件印结论"的病（`|| echo "-> 0 枚写者"` 只在 grep 失败时印），所以是**射程错、不是假绿**。
- **判**：它是票 164 自己的仪器，和 `[5]` 同族，**该修**；但 `.scratch/wisp/**` 在本腿的禁区里，⇒ **下一枚动 `probes/164/v1/` 的写腿做**，改法两条二选一（我具名）：① 把 `:89` 收窄成 `grep -rnE '(\bRoster|\broster|\br)\.Record\('`（今天 4 行），标签不动；② 标签改成「所有 `.Record(`（含 approval/cards 等别类型）」并**另加一行**名册专用计数。修好之前，`[6]` 那段输出**不许**被读成"名册写者 N 枚"。

## ★终裁：`AC#4` **能翻**（判语〔成立·带条件〕，条件三条一条不减）
1. **判据有牙且两向可区分**：V1 红 `:107/:114/:125`、V2 红 `:107/:114`（`:125` 绿），两发都是我自造、各自先证落地、读完单独还原（hash 逐字节回原值）。AC#4 要的那发"取消后仍写输出 ⇒ 必须红"**成立**。
2. **今天的绿是真的**：取消跑句柄、写入口读取消态；我**没有**找到"今天可达且拦不住"的形。⇒ 这一格今天的读数是**绿的**，⛔ 是"还没对象"。
3. **绿说的是"两条缝"，⛔ 是"整条链"**：V3（调用者换 ctx ⇒ 整包全绿）与 V4（守卫之上先落产物 ⇒ 两枚钉子全绿）是**射程边界**；生产里今天**没有**既挂句柄、又走 `Backfill` 的那枚行（`run.go:1182` 只 `MarkRoot`，`AttachCancel` 全仓只 `subagent_197.go:374`），`task.go:459` 那支 0 测试。⇒ 翻框时**必须把这句写进票面**：AC#4 的〔成立〕覆盖"名册那一行的尾巴"，⛔ 覆盖"盘上产物字节"与"调用者传错 ctx"；接 `AttachCancel` 那一行属下一腿，闸门已现成。
- **留给编排者的一句话**：这格能翻，靠的是**判据**（顺序＋行不许动＋两发可区分），⛔ 靠"新造了守卫"——守卫是票 176/197 落的，r2 交的是钉法，归因别写错。

## 我纠正被裁腿的三处数（都在上面各节留了可重跑的尺）
1. `§0.3 ①`「生产调用点 **1 枚**」⇒ **2 枚**（漏 `subagent_197.go:372`；它修的仪器 `[6]` 自己就会印出来）。
2. `§0.3 ②`「非测试 `.Record(` **5 处**」⇒ 它点名的尺给 **4 处**（名册）／**8 处**（全类型）；"5"只在"进名册的生产写点"这把尺下成立。
3. `§2.3` 的 M1「rc=1 红三行」⇒ **只在靶向 `-run` 下成立**；整包尺度是 **10 分钟超时**（`Test197CancelIsPerRowAndNeverCascades` 卡死、`FAIL=0`）。三条红行本身**逐字复现**，但变异命令该写进件里，否则下一个人照它跑整包会以为"没红"。
