# 287-depth-fixture — 票 287 `AC#2`∥`AC#3`∥`AC#4` 落地腿交件读数

⚠ **这份是实现腿自己的交件读数，不是裁决表**：`AC#2` 的缺口审计与对抗验收必须由**另一枚非实现者腿**做
（D22 双角色／`SPEC-12 §4.3`），本腿⛔ 给自己盖章。逐字红句、四件套、门禁日志一律在
`.scratch/wisp/probes/287/r1/`（本件只给结论与点位，⛔ 把结论当读数）。

腿：`287-r1`（写腿，`internal/tools/**` 唯一写腿）。起手锚 `400927c4`（2026-10-15 21:16 +08）、
交件复量锚 `afcec913`（22:04 +08），期间 `dev` 漂了 7 发 ⇒ 本件每一枚行号都是**本腿现跑**的号。
跑的是票 287 的 `AC#2`∥`AC#3`∥`AC#4`；`AC#1` 已闭合（只读腿 `287-a1`），本腿未碰、⛔ 未翻任何 `- [ ]` 框、
⛔ 未改票名/Status/台账、⛔ 未改产码一字（三发种坏全部物理改写＋sha256 还原，见第 3 节）。

写面（逐枚）：`internal/tools/subagent_depth_fixture_287_test.go`（新建，260 行）＋本件＋
`.scratch/wisp/probes/287/r1/**`＋票 287 末尾 1 条 bullet。
越界：零。⛔ `cmd/wisp/**`、⛔ `internal/panel/**`、⛔ `internal/winsec/**`（120-v1 在飞，跑期间两次出现在
`git status` 里，本腿一字未碰、具名登记）、⛔ 三枚冻结件、`frontend/**`／`design/**` 零读零写零转述。

## 1. 现量锚（`grep -n` 于 `internal/tools/subagent_197.go`，起手与收工各跑一次、同值）

| 问的事 | 现量号 | 票面旧号 |
|---|---|---|
| 调用者身份那一读 `parentID := TaskID(ctx)` | **271** | 262 |
| 深度判定入口 `if rec, ok := ...; ok && rec.Kind == TaskKindSubagent {` | **276** | 267 |
| 拒语两行（`拒绝派生：深度上限是 %d…` ＋ `子代理不许再派子代理。`） | **278-279** | 270 |
| 池位判定 `TryAcquireSubagentSlot` | 282 | — |
| 孩子的 `opt.Tools` 被换成包层 provider | **316** | 296 |
| 孩子那一行的真发布点（`opt.AdmitTask` 里的 `PublishSubagent`） | 332-333 | — |
| 目录过滤那一半 `if info.Name == subagentHiddenTool {` | 595-597 | — |
| 名字闸那一半（`拒绝执行：子代理不能再派生子代理…`） | **604-610** | 573-580 |
| 深度常量 `MaxSubagentDepth = 1` | 91（注释 87-90） | — |
| 生产装配把桥交给孩子的地方 `cmd/wisp/run.go:785` | 785 | 782 |

⇒ 派单预告的"锚已腐烂 `+9`"逐枚对上；本腿只用自己的号。

## 2. `AC#2`：造出来的那枚夹具

**形状**：调用者那一程是**真 `agent.Loop`**，工具面是**未包层的真桥**（`Tools: x.h.bridge`）。
产码今天走不到 `:276`——`:316` 把孩子的工具目录换成 `newSubagentToolProvider`，`:604` 的名字闸在桥之前就把
`task.spawn` 拒了（`287-a1` 尺③）。所以唯一能让 `:276` 被走到的形状，就是"结构那一半没装上"的形状；
夹具造的正是它，⛔ 把包层拆掉（那是产码/契约面，本腿一字未动）。

**AC#2 点名缺一个就作不出读数的三件事，逐枚落成真断言**（点位是用例文件里的号）：

| 前置 | 怎么造 | 怎么钉 |
|---|---|---|
| ① 调用者那一行 `Kind == TaskKindSubagent` | 在 `agent.Loop` 的 `AdmitTask` 钩子里调**产码自己的** `Roster.PublishSubagent`（`:76`），这正是产码给孩子落同一枚行的地方（`:333`） | `Look(callerID)` 断 `Kind`，不对即 `t.Fatalf`（`:158-161`）；⛔ 用"字段非空"充当判据 |
| ② 派发经真 loop ⇒ `corr != taskID` | 请求由 C5 接缝的假 provider 在第一程发出，loop 自己铸 `callCorr(taskID, callID, i)` | `corrOf287`（`:124-143`）读**持久化的** `tool_call` 行：非空＋以 task id 为前缀＋**不等于** task id；红句把 `who/taskID/corr` 三枚同屏（`:180-182`）＝票 282/283 那条方法学前置 |
| ③ 请求 `task.spawn` | 同上那一程 | `spawnLog287`（`:103-121`）要 `res.ToolCalls == 1` 且回执 `Name == "task.spawn"` |

**三枚用例，各自守哪一行**：

1. `Test287SubagentCallerRefusedAtTheRosterDepthCheckThroughTheRealLoop` ⇒ 守 `subagent_197.go:276`
   那一支：拒语两句（`:278` 的「拒绝派生：深度上限是 1」＋`:279` 的「子代理不许再派子代理。」——
   后者就是 `AC#1` 尺①数出来**全部 `*_test.go` 零命中**的那句，今天由这一枚钉住）、身份那一维
   （`派生者 <task id>` 必须等于 loop 的 task id、`父 kind` 必须是 `subagent`）、**后果**（名册里不许出现第二代、
   `parent197` 的孩子行仍是 1 枚）、以及 `Rejected == false`（证明这一发读到的是名册判定数、⛔ 闸门数）。
2. `Test287RootCallerWithTheSameShapeStillDerives` ⇒ 反形正控：同一形状、调用者那一行是 `root` 时那一支**不许**拒，
   且必须真派出一枚 `Kind=subagent` 的孩子、回执里的孩子 id 与名册行同一枚。这一枚打掉恒真检测器。
3. `Test287DepthRefusalWinsOverThePoolRefusal` ⇒ 守 `:276` 在 `:282` **之前**返回：把池位占满（4 枚）后，
   子代理调用者读到的必须是深度那句、⛔ 池满那句。

注入只在 C5 接缝（AGENTS §1.3）：唯一的假是 provider；roster／bridge／spawner／loop／journal 全是真的，
⛔ 用 mock 代替真的。⛔ `t.Skip`：全包 SKIP 计数 0（见第 5 节）。

## 3. `AC#2` 判据：三发种坏（四件套齐；逐字红句在 `287r1-M1-*.log`／`M2`／`M3`）

统一尺＝`go test ./internal/tools/ -count=1 -run 'Test287' -v`；种前 hash＝种后还原 hash＝
`86de44c4b79b35172c5ca4bd42862f3659e1d9268ab992f975c282eaa6b172e3`（与 `4d90dfcd` 那发的基线同串）；
每发的 `sed -n '271p;276p'` 复量与内容锚命中（唯一性）由 `run-mutations.sh` 自带打印。

| 发 | 那一行怎么种（落地尺寸 hunks=1／改行数=1） | 红 | 还原 |
|---|---|---|---|
| M1 摘掉深度那一支 | `ok && rec.Kind == TaskKindSubagent` 前插 `false &&` | 用例 1＋用例 3（用例 1 七发点位：`:169`/`:173`/`:179`/`:185`/`:189`/`:192`/`:195`） | `RESTORED-EQUAL-OK` |
| M2 身份改读 corr | `TaskID(ctx)` → `CorrelationID(ctx)`（`:271`） | 用例 1＋2＋3（反形也红：孙代理那行挂到 corr 上，`:221`） | `RESTORED-EQUAL-OK` |
| M3 那一支放宽成"有名册行就拒" | `ok && rec.Kind == TaskKindSubagent` → `ok` | **只有**用例 2（`:213`/`:217`/`:221`） | `RESTORED-EQUAL-OK` |

M1 的红句里**第二代真被造出来**（`名册里出现第二代（1 枚 [df68984f-…]）＝深度上限 1 被走到并破掉`），
⛔ 只报句子就当读数。用例 1 的池位那一发（`:197-199`，钉"深度拒绝先于占位"）在 M1/M2 下不红——
它的牙要"把 `:276` 搬到 `:282` 之后"那种改动才会响，本腿⛔ 为凑读数造第四发，具名登记为**未证**。

## 4. `AC#3` 代价表（⛔ 本票不许落地、⛔ 读成已裁）

先给"设计如此"那一半的产码自述（本腿只引，不裁）：`subagent_197.go:87-90` 逐字写着深度 1
「Enforced structurally (the child's tool directory has no task.spawn in it) **and again by the roster check in
Execute**, because "reachable but refused" is a weaker guarantee than "not offered"」，
`cmd/wisp/run.go:778` 写着「depth 1 by structure, not by counting」。⇒ **:276 是产码自己声明的第二道**，
今天不可达是**设计的后果**，不是遗漏。

再给"它挡住了什么"的现量：`:276` 下游的四支判定在**深度 0 的 root 调用者**身上全部可达、且各有仪器
（池位上限＝`subagent_197_test.go` 的 `Test197SubagentPoolCapsAtBridgeCeiling`；无 task id 的 fail-closed＝
`failclosed_236_teeth_test.go`；`BaseOptions`/`ParentTools` 未接线＝`nilbase_270_test.go`；
`AdmitTask` 未接线＝同一族那枚"子代理永不自批"）。**被结构闸挡住的只有深度那一支自己**——
而它的触发条件（调用者行＝subagent）今天没有任何装配能造出来。

| 形 | 要动的文件（逐枚具名＋本腿现量号） | 会不会把深度上限放开 | 代价（本腿现量，非推理） |
|---|---|---|---|
| **甲：维持现状**（结构闸在前，`:276` 是纵深防御） | **零枚**产码。要动只有注释：`internal/tools/subagent_197.go:87-90`（补一句"第二道今天由票 287 的夹具钉住"）与 `cmd/wisp/run.go:778` 同步 | **不会**：`MaxSubagentDepth = 1`（`:91`）与两处判定都不动 | 深度那一支在生产链上永不可达 ⇒ 它的红**只能由夹具给**（本票已交，M1/M2/M3 三发齐）；读码层证据在 `287-a1`，端到端档在 `cmd/wisp` 面〔待面空补验〕 |
| **乙：摘掉名字闸**（孩子的 `task.spawn` 打到桥，`:276` 成为实际拦下者） | `internal/tools/subagent_197.go:603-612`（`subagentToolProvider.Execute` 那一支）。若连目录过滤一起摘＝再多一枚点位：同文件 `:595-597`。文本面要跟着改：`:544-546` 注释、`:207` 给模型读的说明、`cmd/wisp/run.go:778` 那句"depth 1 by structure" | **不会自动放开**：`:91` 仍是 1、`:276` 仍拒 ⇒ 第二代仍出不来；但保证从"不 offering"降级为"reachable but refused"（产码注释明写这是 weaker），且**只剩一枚运行期判定**守着——任何忘记包层的装配（D46 Tier-1 外部命令那一族、插件面）把风险面交给 `:276` 一处 | 现量（`run-form-b-census.sh`，整包档 `-timeout 180s`）：只摘名字闸＝**PASS=241／FAIL=1**，唯一红的既有仪器是 `Test197ChildCannotDeriveSubagent`（两发点位 `subagent_197_test.go:663`「孩子派孙没有被拒并说明原因」＋`:667`「被拒的那一枚还是打到了父目录上」）；再摘目录过滤＝**PASS=241／FAIL=1**，同一枚函数内四发点位（`+ :648`「孩子的工具目录里还有 task.spawn」＋`:652`）。两形下本票三枚用例**照绿**（它们读的是名册那一支，包层与否不影响）。两趟都没有挂死、都没有新增红以外的账 |

⇒ **交编排者裁**。本腿的读码倾向（⛔ 裁）：甲形的"代价"已经被本票的夹具付清了——第二道现在有仪器、
有后果读数；乙形买到的只是"让第二道成为唯一那道"，卖的是产码注释里那句 stronger guarantee。

## 5. `AC#4`：边界与门禁（逐条带 `rc=`，日志同名在 probes 目录）

| 尺 | 命令 | rc | 读数 |
|---|---|---|---|
| 编译 | `GOFLAGS= go build ./...` | **0** | 全树（含别人的在飞面） |
| 格式 | `gofmt -l internal/tools/subagent_depth_fixture_287_test.go internal/tools/subagent_197.go` | **0** | 输出 0 行 |
| 格式 | `/d/work/base/gopath/bin/gofumpt -l`（同上两枚） | **0** | 输出 0 行 |
| 静态 | `go vet ./internal/tools/` | **0** | 干净 |
| 包级 | `go test ./internal/tools/ -count=1 -v` 改前 | **0** | **239 PASS／0 FAIL／0 SKIP**（起手锚 `400927c4`） |
| 包级 | 同命令改后 | **0** | **242 PASS／0 FAIL／0 SKIP**（`287r1-10-test-post-full.log`） |
| 红名册差集 | `grep '^--- FAIL'` 改前后各一次 | — | 改前 0 枚／改后 0 枚 ⇒ **差集空**；新增 3 枚全 PASS |
| D22 门 | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`（`internal/ 568 Go files`、`cmd/ 136`；含本腿新件） |
| 长度门 | `sh scripts/check-path-length-budget.sh --with-self-test` | **0** | `VERDICT GREEN`，`tracked paths=13361`（终局档现量；同一把尺 21:51 那发是 13332，差 29 枚＝别人在飞的件）、`over-budget=60`、`roster entries=60` |

既有红不算新债、也⛔ 由本腿跑：`cmd/wisp` 的 `TestTicket255RosterStillMatchesTheActualReadSites`（归口票 324，
该编译面 324-r1 独占）与 `internal/panel` 5 枚——本腿**未跑这两面**（派单 ⛔），只引派单数字并具名标"未复跑"。

台面自抓（细节在 probes 件 §5）：① 首版把 `⊂`（U+2282，字符串不豁免）写进 `t.Errorf` 文案 ⇒ d22scan 真红过一次，
改文案后 rc=0；② M3 的整包档跑到 `Test197RowExistsBeforeFirstChildModelCall` 时挂住（M3 把派生一律拒了），
本腿掐掉该趟并具名登记先红的 5 枚 ⇒ ⛔ 当本票读数；③ 突变体第一趟被 Python 文本模式写成 CRLF，
"落地尺寸"账本失真 ⇒ 改逐字节替换后第二趟重跑，hunks=1／改行数=1 逐发对上；④ ⚠ 复跑台件的人先看：本机
`git config core.autocrlf = true`（现量）、`git add` 已对本目录 `.log`/`.bak`/`mutants/*.go` 打过
"LF will be replaced by CRLF" 的警告 ⇒ 用 `subagent_197.go.pre-mutation.bak` 当字节基线之前先跑
`grep -c $'\r'` 证它＝0（本腿交件时现量 0，新用例文件 disk sha＝blob sha `02cbaf10da0cedd4…`），
不是 0 就改用 `git show HEAD:internal/tools/subagent_197.go` 重落一份（probes 件 §5.7）。

## 6. 没做完／待人拍板（具名，⛔ blanks）

- 〔待面空补验〕夹具的"生产链今天仍不可达"在 `cmd/wisp` 端没有端到端台件（本腿无该编译面权限）。
- ⛔ 未裁 `AC#3`（只交代价表）、⛔ 未翻任何 `AC` 框、⛔ 未改票名/Status/台账、⛔ 未动深度这一维任何数值。
- 用例 1 的池位次序那一发（`:197-199`）只对"搬运判定顺序"那种改动敏感，本腿**没有**为它造种坏 ⇒ 具名未证。
- 顺手量到、不归本票：`MaxConcurrentSubagents = 4`（`:85`）与票面 `:11` 的"并发 8"仍分家（`287-a1` §5.5 已登记，
  归口票 211/324 那一族）⇒ 本腿一字未动。
- 只 commit、⛔ push（推送权在编排者，`A813`）。
