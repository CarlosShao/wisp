# 16x-c1 · 六枚契约级票（160／162／163／164／165／168）的冻结文本射程普查

- 程名＝`readonly-16x-c1`（只读普查）｜派单＝`.scratch/wisp/dispatches/2026-09-26-225x-readonly-16x-c1-contract-lines-census.md`
- 被派时刻 `2026-09-26 22:5x`｜**本程锚点（step 0 现量）＝`0d1764731973a998c5d6da5839070151fb9b04a9`**（`dev`）
- **本程只写两枚面**：本文件＋`.scratch/wisp/probes/16x/c1/**`。**不改任何票面、不改 `docs/PLAN.md`／`docs/specs/**`／任何 `.go`／任何 `.yml`。**
- 本文件**只提文字、不落文字**：下面所有"建议逐字改成"都是给编排者执行的草稿，**本程一个字都没往冻结件里写**。
- ⚠ 与前端那支会话**零往来**：§4.7 只答"我们自己的票面有没有要求写 `frontend/**`"，判据＝票面原文＋两把尺，**没有读过 `frontend/**` 的任何内容**（只出现过一次按文件名列举，见 §5 尺C）。

---

## 1. step-0 四件原文

```
$ date
Sat Sep 26 22:51:17 CST 2026

$ git rev-parse HEAD
0d1764731973a998c5d6da5839070151fb9b04a9

$ git rev-parse --abbrev-ref HEAD
dev

$ git status --porcelain | head -20        # 全量枚数见下条
 M .scratch/wisp/issues/169-a-single-decorative-glyph-in-frontend-makes-the-repo-wide-gate-red-and-kills-every-ci-step-after-it.md
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/158/accept-r1/mut/guard.no1.go
 M .scratch/wisp/probes/158/accept-r1/mut/guard.no2.go
 M .scratch/wisp/probes/158/accept-r1/mut/guard.no3.go
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
 D design/index.html
 D design/screens/approval.html
 D design/screens/ball.html
 D design/screens/chat.html
 D design/screens/config.html
 D design/screens/cost.html
 D design/screens/firstrun.html

$ git status --porcelain | wc -l
49
```

**六枚票面都在（枚数＝6，与派单 §0 的第 4 件一致）**：

```
$ ls .scratch/wisp/issues/ | grep -E '^(160|162|163|164|165|168)-'
160-scope-open-returns-a-handle-carrying-its-own-closer.md
162-add-the-patch-one-region-tool-instead-of-rewriting-whole-files.md
163-persistent-shell-session-that-still-looks-one-shot.md
164-background-jobs-cannot-be-read-fix-the-roster-then-add-the-output-leg.md
165-plan-first-then-act-a-brake-the-voice-entry-needs.md
168-built-in-browser-as-a-disable-able-plugin-two-rules-to-set-first.md

$ ls .scratch/wisp/issues/ | grep -cE '^(160|162|163|164|165|168)-'
6
```

**版本口径（本程所有行号的前置声明）**：本程引的每一行都按上面的锚点取（`git show <锚>:<文件>`），
落盘件在 `.scratch/wisp/probes/16x/c1/`：

```
$ git show 0d17647:docs/PLAN.md > .scratch/wisp/probes/16x/c1/PLAN.md.0d17647
$ wc -l .scratch/wisp/probes/16x/c1/PLAN.md.0d17647
3595
$ cmp .scratch/wisp/probes/16x/c1/PLAN.md.0d17647 docs/PLAN.md && echo IDENTICAL
IDENTICAL                      # ⇒ 本锚点上 docs/PLAN.md 工作树＝blob，行号两把尺同一读数
$ git diff --stat HEAD -- docs/PLAN.md     # 空输出＋rc=0＝未改动
$ cmp .scratch/wisp/probes/16x/c1/provenance.go.0d17647 internal/risk/provenance.go
$ cmp .scratch/wisp/probes/16x/c1/bridge.go.0d17647 internal/tools/bridge.go
```

---

## 2. 本程**没**测什么（每一格＝只给了推断、没给凭据，逐条列）

| # | 格子 | 本程做到哪一步 | 没做到的那半 | 档位 |
|---|---|---|---|---|
| 2.1 | 六枚票的**AC 可达成性** | 只读票面＋读被票面点名的那几处码/文字 | **没有跑任何一发变异/用例**（本程是普查，未跑 `go test`、未跑门） | 〔我本轮现跑过的只有 grep/cmp/git show〕 |
| 2.2 | "改后不减少确认次数"这一判断（票 165 的代价栏） | 只核到 `internal/agent/approval/batch.go` 的**建造者注释与阈值常量**（D45-1 的聚合发生在**单次调用内**、且 L2 永不聚合） | **没有实跑一发**证"批完方案后仍逐次弹卡"；这一句是**从注释与常量推的** | 〔我未验，只是推断〕 |
| 2.3 | 票 160 的**跨包误关**那一发（AC#3 反向判据） | 只读了 `OpenScope`/`CloseScope`/`Detector` 的声明与两处生产调用点 | **没有构造那发变异**，也没量"摘掉认身份会怎样" | 〔我未验，只是推断〕 |
| 2.4 | 票 162 的 L1／L2 之争 | 量到 **R8 的 `overwrite` 是生产码里的结构化事实类**（`internal/risk/rules_irreversible.go:10,25`）＋`fs.write` 覆盖行现写 L2（`PLAN.md:2534`） | **没有量到"若真要 L1 需要什么样的可撤销机制"**——我给的"必须先有备份/撤销"是推理，不是读数 | 读数硬、推论那半〔我未验〕 |
| 2.5 | 票 163 的会话形状可行性（PTY／stdin／退出码标记） | 只量了我们侧的约束（ban #4 墙钟、`D33` argv 强制在 `PLAN.md:2507`、`internal/proc/**` 存在） | **外部三仓的原文（dsh 的 `tool-bash-persistent`、`node-pty` 那三条）我一字未读**——它们不在本仓，票面自标的档位我也没复算 | 〔仅指外部读数未取〕；本仓侧有凭据 |
| 2.6 | 票 164 的"名存实无" | **两把尺量过**：生产 `Name() string` 只有 6 枚 `fs.*`；`"task.` 字面量在生产码 0 命中（详见 §4.4 第 5 问） | 没有量"注册表运行时实际条目"（要跑起来才知道）——我引的是 A308 的读数＋我自己的静态两把尺 | 静态读数〔我本轮现跑过〕／运行时那半〔A308 读数，我未复算〕 |
| 2.7 | 票 168 的两张白名单**该定成什么内容** | 只量到"它们尚未定稿"这件事的**六处原文**（§4.6 第 2 问） | **没有提议任何一个具体方法名之外的内容**；我给的两枚方法名（`browser.view.open`/`browser.view.navigate`）是我造的名字，**仓里从未出现过**，属提案不属读数 | 提案（不是读数） |
| 2.8 | `D34` 声明枚数 | 我用两把尺都**没能复现出 A308 那句"36 枚"**（我按行切分得 34 个去重串，且其中 `list`/`cancel` 是被 `reminder.` 前缀折出来的碎片 ⇒ 我的尺子不合格） | **所以本文件不用"36"这个数**；我用的是可复现的**枚举行数 33／33**（命令在 §6） | 〔A308 的 36＝编排者读数，我未复算，且我这把尺证不了〕 |
| 2.9 | `frontend/**` 与 `design/**` 的实际内容 | **零读取**（只在 §5 尺C 里按**文件名**列了 8 枚命中 `shell.exec` 字面的件，未打开任何一枚） | 因此"面板/界面上会显示什么字"这一问，我只从**冻结件已写死的文案**回答（如 `PLAN.md:2025` 那句「点击悬浮球以批准」），**没有看任何现有界面** | 〔按派单 §4 令：故意不测〕 |
| 2.10 | 兄弟程 `161-r3` 的产出 | **零采信**：`tools/d22scan/**`、`ci.yml`、`probes/158/accept-r1/mut/**` 一律未进任何结论 | 无 | 〔按派单 §2 令：故意不测〕 |
| 2.11 | 同源拷贝普查里 `docs/evidence/**` 与 `docs/reports/**` | **主动排除在语料外**（它们是裁决史／历史读数，不是"自称逐字抄自"的契约拷贝） | 所以"某句在证据件里被引过 N 次"本程不答；那里头很多是**带版本的旧行号**，数进来只会污染这张表 | 〔故意不测，理由已核：尺A-1 语料三条〕 |

---

## 4. 六枚各一节（派单 §3 那七问，逐枚）

### 4.1 票 160 — `OpenScope` 改成自带关闭动作的句柄

**这枚的批准粒度不是"改哪一行字"，是"开不开 `internal/risk/**` 这块地"**（派单 §3 第 6 问明令单独写清）。
下面照七问走，但**第 2／3 问在这枚上是次要的**：文字改动只在 `provenance.go` 头注释那一处，
真正贵的是**写面解冻**＋**这枚票自己的契约轴自相矛盾**（第 6 问）。

**① 票面现读**（`.scratch/wisp/issues/160-scope-open-returns-a-handle-carrying-its-own-closer.md`，本程 `Read` 全文，53 行）：

- Status（`:3`）：`立而不派（要动 internal/risk/**＝冻结面，且要兑现 DEFERRED(C25-loop-wiring) ⇒ 改契约＝人工批准，owner 点头之前不派）`
- AC（`:26`–`:30`）：AC#1 未修码读数／AC#2 换形状（硬约束：不新增 goroutine、不用墙钟差、`RiskLevel` 判定语义一字节不许变）／AC#3 反向判据／AC#4 明写治不到什么／**AC#5 契约轴**
- AC#5（`:30`）逐字：**`本票只许动 internal/risk/provenance.go 及其测试＋自己的证据件；docs/PLAN.md、docs/specs/**、thresholds.go、golden、allowlist.txt、tools/d22scan/**、frontend/**、design/** 零字节`**
- 本票不解决（`:34`–`:37`）：不裁 `Q-56`／不做 `OpenTask/CloseTask` 同改造／不把门接进 CI／**只收 C25-loop-wiring 的第 (1) 条**
- 前置条件（`:41`）：①owner 批准动 `internal/risk`（一句话即可，撤销口令写进票面进度）②票 158 结案 → **`:52` 已把②改判**为"票 158 五格成立＋AC#5 残余有票号"

**② 要不要动冻结文本**：**要动 `docs/PLAN.md` 的哪一行？＝不动。**（票面 `:30` 自己钉死 `docs/PLAN.md`、`docs/specs/**` 零字节，本程核这一条与冻结面无冲突。）
真正的冻结面是**包**，且**它的身份在仓规里写着**：

- `.scratch/wisp/issues/README.md:180`（逐字）：`Never modify D1–D47 / C1–C32 / R1–R9 / D43 transition table (contract change = human approval).`
- `.scratch/wisp/issues/README.md:141`（逐字）：`⚠ 本票要动 internal/risk/＝冻结区 ⇒ 改法先报编排者判；` ← **`internal/risk` 的"冻结区"身份在这里，不在 `docs/specs/**` 里**
- `docs/specs/SPEC-12-roadmap-governance.md:38`（逐字）：`改 C1–C32 或 D1–D47 = 人工批准；同步更新 PLAN.md、docs/DECISIONS.md、受影响切片卡。`

被这条批准放开的**具体文字**只有一处（它在冻结包内，是契约注释）：

```
$ git show 0d17647:internal/risk/provenance.go | awk 'NR>=55 && NR<=61 {printf "%d\t%s\n", NR, $0}'
55	// DEFERRED(C25-loop-wiring): the agent loop / tool providers (tickets 10, 20,
56	// 21, 22, 26) must (1) OpenScope at task start and Defer(CloseScope) on the
57	// task's DisposalScope, (2) call Mark(...) on every SPEC-06 §5 sensitive
58	// source output, (3) gate outgoing calls with Inspect(scope, tool, params)
59	// (or CheckText for the TTS/HTTP-body channels), and (4) record the R4 source
60	// name in the tool_call forensics row via the existing DAO — Hit.Fragment
61	// itself must NEVER be persisted (it is sensitive content; native card only).
```

**③ 建议的替换文本**（只改 `:55-57` 那三行的第 (1) 条，第 (2)(3)(4) 条一字不动——票面 `:37` 就写了"只收第 (1) 条"）：

```go
// DEFERRED(C25-loop-wiring): the agent loop / tool providers (tickets 10, 20,
// 21, 22, 26) must (1) LANDED as of ticket 160: opening a scope hands back a
// handle that owns its own close, (2) call Mark(...) on every SPEC-06 §5 sensitive
```

（第三行续接原 `:57` 后半句 `source output, (3) gate outgoing calls ...`，即**整段只把 "(1) OpenScope at task start and Defer(CloseScope) on the task's DisposalScope," 换成上面那半句**，其余字节不动。）

**④ 新增还是改写**：**改写既有文字**（一枚 `DEFERRED` 标记句 → 就地记"第 (1) 条已兑现"）。
比 162/163/164 那三枚"加一行"贵，因为**这枚标记正对着双向交叉核对**：
`PLAN.md:1546` 逐字 `**强制要求**：代码内每个 ` + '`// DEFERRED(D-xx): ... → docs/DEFERRED.md#锚点`' + ` 标记必须在本表有对应条目，**反向亦查**`，
`docs/specs/SPEC-12-roadmap-governance.md:94` 逐字 `**双向交叉核对**：代码内每个 ` + '`// DEFERRED(D-xx): … → docs/DEFERRED.md#锚点`' + ` 必须在表有对应条目，`。
⚠ 顺带一枚**本程量到的既存缝**（不是 160 造的，但 160 一动就会撞上）：

```
$ grep -c "loop-wiring" .scratch/wisp/probes/16x/c1/PLAN.md.0d17647   # 尺1：内容串
0        （grep rc=1＝零命中）
$ grep -c "C25-loop-wiring" docs/specs/SPEC-12-roadmap-governance.md  # 尺2：另一枚文件里的同一名字
0        （grep rc=1＝零命中）
```
⇒ `DEFERRED(C25-loop-wiring)` 在 `PLAN.md` 与 `SPEC-12` 的登记表里**都没有条目**（两把尺都零命中）。
它**严格说不在 `D-xx` 那族里**（`PLAN.md:1546` 的规矩字面只覆盖 `DEFERRED(D-xx)`），所以这不算违令；
但 160 把它"就地记为已兑现"之后，表里**依然**没有这一条 ⇒ **登记表的 1:1 双向核对从此永远数不平它**。
⚠ 本程**不**建议 160 顺手补这条（那要动 `PLAN.md` §7，正踩票面 `:30` 的零字节），只把它的**归属写清**：
补条目＝另一次批准，或归 §7 登记表的既有 owner 待定项。

**⑤ 同源拷贝有几份**（这一问在这枚上＝"这段注释所钉的成对规矩在别处还有几处自称同源"）：

- **尺A（词形）**：`grep -rnE "逐字|照抄|verbatim|唯一来源|权威表" docs/PLAN.md docs/specs AGENTS.md | wc -l` → **28** 行命中
  （原始逐行输出在 `.scratch/wisp/probes/16x/c1/copy-census.log` 尺A-1 段，那次跑的是**更宽**的八支样式所以得 33；
  两支的差别＝多出来的 `一一对应`／`与 D3x 表一致`／`与 D43 转移表` 三支）。
  ⚠ **这把尺的已知缺陷（写下来免得下一位当它是精确计数）**：`逐字` 无词边界，会把 `逐字节稳定`（`PLAN.md:3003`、`SPEC-02:131`、`SPEC-05:106`）也算进来
  ⇒ **它只能当"有没有人声称逐字拷贝"的上界用**；本格的结论（0 份）靠的是逐枚读那 28 行，不是那个数。
  **28 行里没有任何一行声称"逐字抄自 `provenance.go` 的头注释"** ⇒ 文档侧的同源拷贝：**0 份**。
- **尺B（内容串）**：

```
$ grep -rn "OpenScope at task start" docs/ internal/ cmd/ | grep -v docs/evidence
internal/risk/provenance.go:56:// 21, 22, 26) must (1) OpenScope at task start and Defer(CloseScope) on the
internal/risk/provenance.go:453:			logf("risk/C25: Inspect/CheckText on unregistered scope %q while other scopes hold taints: fail-closed R4 (ensure OpenScope at task start)", scopeID)
```
  ⇒ 那半句原文在**文档里 0 处**、在**码里 2 处且都在同一枚文件内**（`:453` 是同包的一句日志文案，不是拷贝）。
- **尺B-补（被改那条标记的名字在别处被复述了几次）**：

```
$ git grep -n "C25-loop-wiring" 0d17647 -- ':!.scratch' ':!docs/evidence' ':!docs/reports'
internal/risk/provenance.go:55:// DEFERRED(C25-loop-wiring): the agent loop / tool providers (tickets 10, 20,
internal/tools/bridge.go:627:// DEFERRED(C25-loop-wiring) item (1). Idempotent; Execute opens lazily so a
cmd/wisp/panel_assets.go:164:// DEFERRED(C25-loop-wiring) item 1); a CLI probe has no task, so it names one
```
  ⇒ **除本体之外还有 2 处注释在复述"第 (1) 条"**（`bridge.go:627`、`panel_assets.go:164`）。
  它们**不是"逐字拷贝"，是"跟着这条标记说话的注释"**——把 `:56` 改成 LANDED 而不改它们，同一枚标记就会出现**三种状态**（本仓"名册与注释各说各话"那一族的形状）。
  ⚠ 那 2 处**本来就落在本程第 ⑥ 问必须放开的两枚文件里** ⇒ 不缺批准，只缺"一起改"这一句。
- **受签名改动直接影响的调用者（两把尺各数一次，都是 2 枚）**：

```
$ git grep -l "OpenScope" 0d17647 -- '*.go' ':!*_test.go' ':!.scratch' | grep -v internal/risk/provenance.go
0d17647:cmd/wisp/panel_assets.go
0d17647:internal/tools/bridge.go
$ git grep -n "prov.OpenScope" 0d17647 -- '*.go' ':!*_test.go' ':!.scratch'
0d17647:cmd/wisp/panel_assets.go:232:	prov.OpenScope(taintSourceScopeID)
0d17647:internal/tools/bridge.go:642:		b.prov.OpenScope(taskID)
```
  逐枚现读（本锚点）：
  - `internal/tools/bridge.go:633-644` `func (b *Bridge) OpenTask(...)`，其中 `:638-642` 就是票 160 进度段说的那两套并行记账：
    `open := b.scopes[taskID]` / `b.scopes[taskID] = true` 与 `:642 b.prov.OpenScope(taskID)` 分列两行（**本程读码自取，与票面 `:46-47` 同读**）；
  - `cmd/wisp/panel_assets.go:232`，且 `:162-165` 的注释自陈它是**故意不关**的（逐字，跨两行折行）：
    `a CLI probe has no task, so it names one` / `and closes nothing - the engine lives and dies inside this process.`
    ⇒ **这一枚不是遗漏，是与新形状直接冲突的一形**：句柄自带关闭动作之后，"开一个、不关、随进程死"要么显式写成"故意丢弃句柄"，要么被类型系统挡住。
- ⚠ **三把尺都数过，才允许写下面这句**：**文档侧同源拷贝＝0 份**；**同文件内复述＝1 处**（`:453`）；
  **跨文件复述（会随这条标记变陈的注释）＝2 处**；**受签名改动影响的真实调用者＝2 枚**。
  这**是三件不同的事**，别并成一句。另：`docs/evidence/s1/{154-host-id-never-closed-r1.md:336, 151-task-scope-never-closed-r1.md:200,299}` 也复述过（3 处），
  但那是裁决史、不是契约拷贝（§2 第 2.11 格的排除理由），**不计进"要同批改的份数"**。

**⑥ 不动它能不能落地**：**不能**，且**理由与文字无关**——
签名从 `OpenScope(scopeID string)` 变成"交回句柄"，**编译器会当场要 2 枚调用者一起改**，
而票面 `:30` 的契约轴**只许动 `provenance.go` 及其测试**。
⇒ 这一格本程必须说白：**票 160 现在的 AC#5 是自相矛盾的**，实现程无论怎么走都会越界（或是"用 mock 代替真的"那类假绿）。
**两种解法，我推荐乙**：
甲＝owner 只批"开地"，编排者随后**改票面**把 AC#5 扩到那 2 枚调用者（票面归编排者，`.scratch/wisp/issues/README.md:198` 已写明"那张票面归我"）；
乙＝**批语里一次划清射程**（照 09-26 `thaw-panel-for-145` 那份特批件的写法：owner 给一句话、编排者把射程写死到文件名），即下面第 ⑦ 问那句。
⚠ 撤销口令照既有制度必须随批语一起给（`.scratch/wisp/dispatches/README.md` 每份记录固定五字段：程名·被派时刻·锚点 sha·授权范围（含"不含哪些"）·撤销口令）。

**⑦ 要 owner 点的那一句（一句话，可直接回"照这个改"）**：

> **「批我动 `internal/risk/provenance.go`（含其测试），并连带放开 `internal/tools/bridge.go`、`cmd/wisp/panel_assets.go` 这两枚唯一生产调用点为配合签名改动所必需的最小改动；`RiskLevel` 判定语义、`D43` 转移表、`docs/PLAN.md`、`docs/specs/**` 一字不动；撤销口令＝"160 别动"。」**

推荐错了的代价（照 `thaw-panel-for-145` 的写法把"代价"随批语一起摆出来）：
放开那 2 枚调用者＝同时打开 `internal/tools`（票 158/161 的地界）与 `cmd/wisp`（装配根）两片写面 ⇒
**代价是同片写码程要串行**（票面 `:41` 前置条件③"编队空"这条仍在），最坏情况是这枚把 161 的尾巴顶到后面；
**不放开**的代价＝这枚永远只能停在"普查尺事后数"，而票面 `:46-50` 已经量到"门装上了还是能从门旁边走过去"那一发。

---

### 4.2 票 162 — 补"只改这一小块"那枚工具（`fs.edit`）

**① 票面现读**（`.scratch/wisp/issues/162-add-the-patch-one-region-tool-instead-of-rewriting-whole-files.md`，`Read` 全文，52 行）：

- Status（`:3`）：`立而不派（往内置工具表加东西＝改 D34＝人工批准；owner 已答"要的"，但那是批能力要，路径他还没选 ⇒ 见"我的推荐"）`
- AC（`:32`–`:37`）：AC#1 未修码读数／AC#2 四枚拒绝各一枚用例／AC#3 换行 BOM 正反／AC#4 原子性／AC#5 风险档与门控／**AC#6 契约轴**
- AC#5（`:36`）逐字：**`fs.edit 判定必须与 fs.write 同族（L1）`**，且**越界路径必须被 `PathResolver` 拒**；不许新增豁免、不许动 `allowlist.txt`
- AC#6（`:37`）逐字：**`本票只许动 docs/PLAN.md 的 D34 那一行（需 owner 单独批准，且批准原文要落台账）＋ internal/tools/** 实现与测试＋自己证据件。docs/specs/**、internal/risk/**、thresholds.go、golden、frontend/**、design/** 零字节`**
- 本票不解决（`:41`–`:43`）：不做"人逐块批准改动"那个界面（`:41`）／不做机器自审工具／不补 `fs.read` 分页
- **owner 已答的那半**（`:51`）：`owner 已答"要的"＝批能力要；路径（核心/插件/内部不露）他未选`

**② 要动的具体行（file:line ＋ 现在逐字写着什么）**：**两枚**，一枚是必动的权威表，一枚是被它牵着动的拷贝。

```
$ awk 'NR>=2529 && NR<=2535 {printf "%d\t%s\n", NR, $0}' .scratch/wisp/probes/16x/c1/PLAN.md.0d17647
2529	| 工具 | 动作 | RiskLevel | capability | 落切片 | 备注 |
2530	|---|---|---|---|---|---|
2531	| `fs.read` | 读文本/二进制摘要 | **L0**（授权目录内）/ L2（越界） | `fs.read` | S3 | 经 C26；长输出走宿主内部 spill（**不门控**，见注②）；结果打 taint |
2532	| `fs.list` | 列目录 | L0 / L2（越界） | `fs.read` | S3 | — |
2533	| `fs.write` | **新建**文件 | **L1** | `fs.write` | S3 | **temp + 原子 rename**（D31 已定，因取消不原子） |
2534	| `fs.write` | **覆盖已存在** | **L2** | `fs.write` | S3 | C19/R8 不可逆 |
2535	| `fs.move` | 移动/重命名 | **L1** | `fs.write` | S3 | **跨盘 = 复制+删 → L2** |
```

⇒ **插入点＝`docs/PLAN.md:2534` 之后**（紧跟"覆盖已存在"那一行，同族相邻）。
它**同时牵住**这三行，本程建议它们**一字不动**：

```
$ awk 'NR==2524 || NR==2525 || NR==2981 || NR==1176 {printf "%d\t%s\n", NR, $0}' .scratch/wisp/probes/16x/c1/PLAN.md.0d17647
1176	- 工程要求：`fs.write` 必须走**临时文件 + 原子 rename**，否则取消会留损坏文件；
2524	RiskLevel 列是 **C19 融合后的默认结论**；实际判定仍由 C19 的 R1–R9 规则集在调用时算出
2525	（例如越界路径会把 L0 升到 L2）。
2981	| R8 | 不可逆性：永久删除 / 覆盖已存在 / 关机重启 / 关闭窗口 / 发送类 | → **L2** |
```

**③ 建议的替换文本（逐字可贴，一行；照现有表格排版）**：

```
| **`fs.edit`** | **字面定位替换已存在文件里的一小段** | **L2** | `fs.write` | S3 | 定位＝字面子串且**必须唯一命中**；匹不上/多处命中＝硬错且**绝不落盘**；写走 temp＋原子 rename（`:1176` 同规）；**R8 覆盖已存在 → L2**，与 `fs.write` 覆盖行同级 |
```

⚠⚠ **本程在这里改了推荐，且必须让 owner 看见我改了什么**：票面 `:36` 现在写的是 **L1**，
我给的是 **L2**。凭据三条（都是本程现量，不是推理）：

```
$ git show 0d17647:internal/risk/rules_irreversible.go | awk 'NR==10 || NR==25 {printf "%d\t%s\n", NR, $0}'
10		"overwrite":    "覆盖已有内容",
25			classes = append(classes, "overwrite")
$ git show 0d17647:internal/tools/fs_write.go | awk 'NR==17 {printf "%d\t%s\n", NR, $0}'
17	//	fs.write  new file L1 / overwrite L2 (R8)   D31 temp+atomic-rename writer
$ awk 'NR==2534 {print}' .scratch/wisp/probes/16x/c1/PLAN.md.0d17647
| `fs.write` | **覆盖已存在** | **L2** | `fs.write` | S3 | C19/R8 不可逆 |
```
⇒ 在**现有语义**下，一枚"替换已存在文件的旧内容"的调用，事实钩子必然带 `overwrite`，R8 必然判 **L2**；
把 D34 那行写成 L1 ＝**表与判定器互相矛盾**，而 `R1–R9` 与 `D34` 都在禁改名单里（`README:180`）。
⇒ 想真拿 L1，必须**先有可撤销那味**（写前把原件留一份、面板可撤销，类比 `fs.trash` 的 L1 论证 `PLAN.md:2536`）——
**那是新机器，不是表格加一行**，且要另一次批准。〔这一句是推断，见 §2 第 2.4 格〕
⇒ **验收面也跟着**，两处文字都要求"判定必须等于表"：

```
$ grep -n "与 D34 表一致" .scratch/wisp/probes/16x/c1/PLAN.md.0d17647 docs/specs/*.md
.scratch/wisp/probes/16x/c1/PLAN.md.0d17647:1668:    → **每条须成功，且风险级判定与 D34 表一致**（`fs.trash` 必须是 L1 不是 L2）
docs/specs/SPEC-00-product-overview.md:119:| ① 系统操作与文件治理 | … | S3 | 10 条真实指令逐条跑，风险级判定与 D34 表一致（SPEC-10 §7.1） |
```
⇒ 表里写 L1、判定器判 L2 ⇒ **场景①那 10 条的验收格自己打红**（这两行就是那格的判据原文）。

**④ 新增还是改写**：**新增**（在 `:2534` 后插一行）。
**不**改写 `:2533`/`:2534`/`:2524`/`:2981`/`:1176`/`SPEC-06:41` 任何一枚既有文字——
这是这张表里**最便宜的一种改动**，也正好是票面 `:37` 已经写明"批准原文要落台账"的那一档。
（对照：若走"把 `fs.write` 覆盖行的措辞改成能涵盖 edit"那条路＝**改写既有文字**，本程**不推荐**，见 §2 与 §3 的"改写"列。）

**⑤ 同源拷贝有几份（两把尺都数过）**：

- **尺A（词形，谁自称是 D34 的拷贝）**：

```
$ grep -n "内置工具权威表" docs/specs/SPEC-07-tools-and-plugins.md docs/PLAN.md AGENTS.md
docs/specs/SPEC-07-tools-and-plugins.md:1:# SPEC-07 · 内置工具权威表与插件系统
docs/specs/SPEC-07-tools-and-plugins.md:32:## 3. D34 内置工具权威表（唯一来源；落 `docs/TOOLS.md`）
docs/PLAN.md:351:> ⚠ **本节的「六件套」清单已被 D34 内置工具权威表取代。**
docs/PLAN.md:1747:| **`docs/TOOLS.md`** | **第四轮新增**：D34 内置工具权威表（工具名/动作/RiskLevel/capability/落哪个切片）。**必须单独成文** …
docs/PLAN.md:2514:## 16.5 D34 — 内置工具权威表（**取代 D14 的六件套；`docs/TOOLS.md` 的唯一来源**）
AGENTS.md:152:| D34 | 内置工具权威表（取代 D14 的六件套；`docs/TOOLS.md` 的唯一来源） | `PLAN.md:2514` |
```
  逐枚判读：`PLAN.md:2514` 是**本体**；`PLAN.md:351`（D14 的作废声明）与 `:1747`（交付物清单）是**指回它**；
  `AGENTS.md:152` 是**薄索引一行**（只抄标题＋给行号，**不抄表格**）；
  **真正"另抄一份表格"的＝1 枚文件＝`docs/specs/SPEC-07-tools-and-plugins.md`，那张表在它里面（标题在 `:32`，行体在 `:37-71`）**。
- **尺B（内容，把被抄的具体内容串当尺）**：`fs.*` 那一族行在拷贝里的落点（**行号＝该文件绝对行号**）：

```
$ awk 'NR>=37 && NR<=71 && /\| `?fs\./ {printf "%d\t%s\n", NR, $0}' docs/specs/SPEC-07-tools-and-plugins.md
39	| `fs.read` | 读文本/二进制摘要 | L0（授权内）/ L2（越界） | fs.read | S3 | 经 C26；长输出宿主内部 spill；结果打 taint |
40	| `fs.list` | 列目录 | L0 / L2（越界） | fs.read | S3 | — |
41	| `fs.write` | 新建文件 | L1 | fs.write | S3 | temp + 原子 rename（D31） |
42	| `fs.write` | 覆盖已存在 | L2 | fs.write | S3 | C19/R8 |
43	| `fs.move` | 移动/重命名 | L1 | fs.write | S3 | 跨盘 = 复制+删 → L2 |
44	| `fs.trash` | 移入回收站 | **L1** | fs.write | S3 | ⭐ 可用性收益最大：场景①大部分操作 L2→L1，缓解 B2 |
45	| `fs.delete` | 永久删除 | L2 | fs.write | S3 | **默认不注册**，`[fs] delete_enabled=true` 才有 |
48	| `search.files` | 按名找文件 | L0 | fs.read | S3 | 白名单内 |
49	| `search.content` | 全文检索 | L0 | fs.read | S3 | 结果打 taint |
54	| `doc.read` | PDF/docx/pptx → 文本 | L0 | fs.read | S3 | xlsx 与 OCR DEFERRED |
```
  **两张表的枚数现在相等**（这把尺最容易漏，所以两条命令分开）：

```
$ awk 'NR>=2531 && NR<=2563' .scratch/wisp/probes/16x/c1/PLAN.md.0d17647 | grep -c "^| "
33
$ awk 'NR>=39 && NR<=71' docs/specs/SPEC-07-tools-and-plugins.md | grep -c "^| "
33
```
⇒ **D34 表格的逐行同源拷贝＝1 份**（`SPEC-07:37-71`），它的 `fs.write` 覆盖行在 **`SPEC-07:42`** ⇒ 我的新行要**同批镜像到 `SPEC-07:42` 之后**。

```
$ awk 'NR>=2531 && NR<=2563' .scratch/wisp/probes/16x/c1/PLAN.md.0d17647 | grep -c "^| "
33
$ awk 'NR>=39 && NR<=71' docs/specs/SPEC-07-tools-and-plugins.md | grep -c "^| "
33
```
⇒ **D34 表格的逐行同源拷贝＝1 份**（`SPEC-07:37-71`），它的 `fs.write` 覆盖行在 **`SPEC-07:42`** ⇒ 我的新行要**同批镜像到 `SPEC-07:42` 之后**。
- ⚠ **第三枚会跟着变谎的地方**（词形尺抓不到、内容尺也抓不到，只有读码看得到）：

```
$ git show 0d17647:internal/tools/fs.go | awk 'NR>=318 && NR<=323 {printf "%d\t%s\n", NR, $0}'
318	// BuiltinFSEntries returns the whole D34 fs family the configuration allows:
319	// the L0 pair plus the write half (fs.write / fs.trash / fs.move), with
320	// fs.delete present ONLY when [fs] delete_enabled is true (FSDeps.
321	// DeleteEnabled). Callers that want the write half on its own use
322	// BuiltinFSWriteEntries.
323	func BuiltinFSEntries(d FSDeps) []Entry {
```
  ⇒ `fs.go:318` 那句 **"the whole D34 fs family"** 是一个**自称等于全表的断言**：
  只往表里加 `fs.edit`、不往这里加，这行注释就从那一刻起变成假话（而票面 `:37` 恰好**允许**动 `internal/tools/**`，所以这条不缺批准、只缺提醒）。
- ⚠ **没有机器尺在盯这张表的对齐**（两把尺，各自量的是不同形状的东西）：

```
$ git grep -l "docs/PLAN" 0d17647 -- '*.go' '*.sh' '*.yml' '*.ps1' ':!.scratch' ; echo rc=$?
rc=1                       # 尺1（零命中）：没有任何码/脚本/CI 读 docs/PLAN.md
$ git grep -n "SPEC-07" 0d17647 -- '*.go' '*.sh' '*.yml' '*.ps1' ':!.scratch' | grep -E "table|§3|roster|名册"
cmd/wisp/run.go:331:	// shortcut). SPEC-07 §3's frozen S1 roster names the first two builtins
internal/tools/bridge.go:100:// through (SPEC-01 §3 / SPEC-07 §2). It implements agent.ToolProvider, the
internal/tools/doc.go:2:// funnels every capability through enforcement (SPEC-01 §3, SPEC-07 §2).
internal/tools/fs_write.go:29:	// the host-internal artifacts channel (SPEC-07 §3), not a gated tool.
```
  尺2 的 4 枚命中**逐枚判读**：`bridge.go:100` 与 `doc.go:2` 指的是 SPEC-07 **§2**（接口面），与表格无关；
  **真点名 §3 那张表的只有 2 枚**（`run.go:331` 与 `fs_write.go:29`），它们**不抄表格文字、只把那张表当权威引用** ⇒
  词形尺（"自称逐字抄自"）给 **0**、内容尺（表格行）给 **1 份**（`SPEC-07:37-71`）、"拿它当权威"的注释给 **2 枚**。
  ⚠ 三个数不是一回事，**别并成一句**；**但"没有机器比对"这件事两把尺都成立**（尺1 零命中＋尺2 四枚全是注释、无解析）。
⇒ **表与码、`PLAN.md` 与 `SPEC-07` 三者之间的漂移是静默的**（本仓"改一把尺之前先数它有几份拷贝"那条教训在这里就长这样）。

**⑥ 不动 `D34` 能不能落地**：**不能（附带条件：只能"先只写码、名册那行由 owner 一次批完"）**。最短理由三条：
(1) 工具名必须由 C1 注册器过校验，`PLAN.md` 表里没有这个名字＝**实现比权威表多一枚**，正是 `Q-58` 那一族（`cmd/wisp/run.go:330-340` 那段注释就是"两个冻结名打架、没人裁定"停在码里的现物）；
(2) 票面 `:36` 自己要求"判定必须与 `fs.write` 同族"，而同族结论在表里；
(3) `docs/specs/**` 若不同批镜像（`SPEC-07:42` 之后），**SPEC-07 那张自称权威表的拷贝当场落后一行** —— 而票面 `:37` **明令 `docs/specs/**` 零字节** ⇒
**本程量到票 162 的 AC#6 与它自己要动的东西冲突**：批语必须把 `SPEC-07` 那一枚镜像行**一起放开**，否则这枚票交出来的就是**一张已知会漂的表**。

**⑦ 要 owner 点的那一句**：

> **「批 `docs/PLAN.md` 的 D34 在 `fs.write` 覆盖行（`:2534`）之后插一行 `fs.edit`，**风险级按 R8 写 L2**（票面原写 L1，我改判，理由在证据件 §4.2③），同批把这一行镜像进 `docs/specs/SPEC-07-tools-and-plugins.md:42` 之后；除此之外 `docs/specs/**` 其余零字节，批准原文我落台账（现最大号 `A316`）。」**

不按推荐来的代价：坚持 **L1** ⇒ 要么验收格自己打红（`PLAN.md:1668`），要么**回去改 R8／改判定器**（`R1–R9` 禁改名单＋`internal/risk/**` 冻结面，比加一行贵一个数量级）；
不镜像 `SPEC-07` ⇒ 留一行静默漂移，下一枚读 `SPEC-07` 的程会照旧表实现。

---

### 4.3 票 163 — 常驻 shell 会话（对外仍像一次性）

**① 票面现读**（`.scratch/wisp/issues/163-persistent-shell-session-that-still-looks-one-shot.md`，`Read` 全文，44 行）：
- Status（`:3`）：`立而不派（要扩 shell.exec 的语义＝动 D34/D46 那条线，人工批准；且它和票 161 抢同一批审批判据）`
- 硬约束（`:18`–`:21`）：不许墙钟差／不许裸 `go func(`／会话寿命必须绑任务作用域／"别照抄 Windows 那一半"清单
- AC（`:25`–`:29`）：AC#1 未修码读数／AC#2 认退出码不靠时间／AC#3 会话绑任务（取消必收）／**AC#4 审批形状不许松，且 `:28` 逐字写着 `⚠ 不许动 internal/agent/approval/**，要动先回来报`**／AC#5 契约轴（`:29`）
- AC#5（`:29`）逐字：**`只许动 docs/PLAN.md 里 shell.exec 那一行（需 owner 单独批准）＋ internal/tools/**／internal/proc/** 实现与测试＋证据件。frontend/**、design/**、thresholds.go、golden、allowlist.txt、tools/d22scan/** 零字节`**
- 本票不解决（`:33`–`:35`）：`:33` 不做"给人用的集成终端"那个界面；`:35` **不引入 PTY 依赖**（要引依赖是另一次批准）

**② 要动的具体行**：票面自己指的是 `docs/PLAN.md:2563`（`shell.exec` 那一行）。逐字现读：

```
$ awk 'NR==2563 {printf "%d\t%s\n", NR, $0}' .scratch/wisp/probes/16x/c1/PLAN.md.0d17647
2563	| `shell.exec` | 任意命令 | **L2**（白名单命中 → L1） | `shell` | S3 | **默认禁用**（D14 结论保留）；**argv 向量强制**（D33） |
```

**但本程量到：动这一行是最贵的那条路。** 因为 `shell.exec` 这条线在定稿里被**三处**文字钉住，其中一处明确写着"不改"：

```
$ awk 'NR==2507 || NR==2640 || NR==2662 {printf "%d\t%s\n", NR, $0}' .scratch/wisp/probes/16x/c1/PLAN.md.0d17647
2507	- **`shell.exec` 必须强制 argv 向量** → 参数类型是 `argv: string[]`，**不接受单一命令字符串**。
2640	> 而 D14 恰恰把 `shell.exec` 默认禁用了，且 D33 强制 argv 向量、拒绝命令字符串。**
2662	- `shell.exec` = **通用任意命令**，默认禁用，开启后每次 L2，argv 向量强制。**保持不变。**
```
- `:2507` 属 **D33 安全必修五项**（`PLAN.md:2357` 的标题逐字：`## 16.4 D33 — 安全必修五项（**实现前必须全部落地，不得推迟到「以后加固」**）`）；
- `:2662` 属 **D46**，而且**它已经预告了今天这件事**：票 163 要动 `shell.exec` ＝ 直接撞一句写着"**保持不变**"的定稿文字。
⇒ **走"扩 `shell.exec` 语义"这条路＝改写三枚已定稿的句子（其中两枚在 D33/D46），而不是加一行。**

**③ 建议的替换文本**：**不改 `:2507`／`:2563`／`:2640`／`:2662` 任何一字**，改为**在 `:2563` 之后插一行新工具**（形状与 `fs.edit` 同策）：

```
| **`shell.session`** | **在一条活着的会话里跑下一条命令（工作目录与环境延续）** | **L2** | `shell` | S3 | 一次性外观、常驻内核；**会话寿命绑任务作用域**；超时＝真 deadline（ban #4）；每发命令仍按 R6 判；**不引 PTY**（引依赖另批） |
```

**④ 新增还是改写**：**新增**（推荐路，一行）／若走票面原路（扩 `shell.exec` 语义）＝**改写 3 枚既有句**（`:2507`、`:2563`、`:2662`）＋连带 `docs/specs/SPEC-06:139-140`。
**这条差价就是本格要给 owner 看的东西。**

**⑤ 同源拷贝有几份（两把尺）**：
- **尺A（词形）**：`grep -n "内置工具权威表" ...` 同 §4.2 第 ⑤ 问 ⇒ 表格拷贝仍只有 **1 份**（`SPEC-07`），`shell.exec` 那行在拷贝里的落点＝**`SPEC-07:71`**（现读：`| `shell.exec` | 任意命令 | L2（白名单命中 → L1） | shell | S3 | **默认禁用**；argv 向量强制 |`）。
- **尺B（内容串）**：`shell.exec`＋（默认禁用／argv 向量）的组合，本锚点全语料 **7 处**，逐枚判读不同性质的"拷贝"：

```
$ grep -rnE '`shell\.exec`.*(默认禁用|argv 向量)' docs/PLAN.md docs/specs AGENTS.md 2>/dev/null
docs/PLAN.md:2029:| 10 | ~~**邮件 / 日历 / IM 明确 REJECTED**~~ → ✅ **已定案（2026-09-18）** | …而 **D14 把 `shell.exec` 默认禁用了，不能为接一个 CLI 就打开整个 shell**…（整行 1.7 KB，全行见 copy-census.log 尺B-1 段）
docs/PLAN.md:2507:- **`shell.exec` 必须强制 argv 向量** → 参数类型是 `argv: string[]`，**不接受单一命令字符串**。
docs/PLAN.md:2563:| `shell.exec` | 任意命令 | **L2**（白名单命中 → L1） | `shell` | S3 | **默认禁用**（D14 结论保留）；**argv 向量强制**（D33） |
docs/PLAN.md:2640:> 而 D14 恰恰把 `shell.exec` 默认禁用了，且 D33 强制 argv 向量、拒绝命令字符串。**
docs/PLAN.md:2662:- `shell.exec` = **通用任意命令**，默认禁用，开启后每次 L2，argv 向量强制。**保持不变。**
docs/specs/SPEC-06-security-gatekeeping.md:139:- `shell.exec` 强制 argv 向量（`argv: string[]`，不接受命令字符串）；确需字符串 →
docs/specs/SPEC-07-tools-and-plugins.md:71:| `shell.exec` | 任意命令 | L2（白名单命中 → L1） | shell | S3 | **默认禁用**；argv 向量强制 |

$ 同上 | wc -l
7
```
  （⚠ `:2029` 那一枚的**整行有 1.7 KB**，本块里我用 `…` 折了一次中段，**枚数与判定按整行读**；整行原文在 `.scratch/wisp/probes/16x/c1/copy-census.log` 尺B-1 段第 1 条。
  ⚠ 这一支样式**漏掉了 `AGENTS.md`**（它不含 `shell.exec`＋默认禁用的组合，实测 0 命中），所以语料里写了 `AGENTS.md` 也只是把它当被扫件，不是当命中来源。）
  分类（**别把 7 当成"7 份表格拷贝"**）：表格行拷贝 **1**（`SPEC-07:71`）＋本体 **1**（`:2563`）；
  **规矩断言 4**（`:2507` D33／`:2662` D46／`SPEC-06:139-140`／`:2640` 引言）；**论证性复述 1**（`:2029` §15 第 10 项定案正文）。
  ⚠ 另有 **1 处码侧注释**引用同一条规矩：`internal/config/unwired.go:70` 与 `:76` 两句 `lands: "it lands with the shell.exec tool (SPEC-06 sec 10…)"` —— 见下第 ⑥ 问，那才是这枚票真正的隐形第三写面。

**⑥ 不动它能不能落地**：**不能**，且这枚有**两枚连带的硬事实**必须让 owner 知道：

1. **`shell.exec` 本身今天也是名存实无**（与 164 同族，两把尺）：

```
$ git grep -n '"shell\.exec"' 0d17647 -- '*.go' ':!*_test.go' ':!.scratch'
（零命中，rc=1）
$ git grep -nE 'Name\(\) string[[:space:]]*\{ return "shell' 0d17647 -- '*.go' ':!.scratch'
（零命中，rc=1）
```
   ⇒ 票面 `:3` 写的"扩 `shell.exec` 的语义"在**代码侧没有那个被扩的东西**：这不是"改一个已有工具"，是**同时新造两枚工具**（`shell.exec` 本体＋会话腿）。
2. **仓里有一枚"诚实钥匙"登记簿会跟着变谎**（本程读码自取，票面未提）：

```
$ git show 0d17647:internal/config/unwired.go | awk 'NR>=57 && NR<=66 {printf "%d\t%s\n", NR, $0}'
57	// unwiredKeys is the guard table. Deleting a row is part of landing the
58	// capability it names: a load-time error must not outlive the consumer that
59	// justifies it.
60	var unwiredKeys = []unwiredKey{
61		{
62			path:    "risk.shell_enabled",
63			missing: "no shell.exec tool is registered in internal/tools, so nothing reads this flag",
64			lands:   "it lands with the shell.exec tool itself (SPEC-07 sec 3, S3, default-disabled per D14)",
65			fires:   func(c *Config) bool { return c.Risk.ShellEnabled },
66		},
```
   ⇒ `[risk] shell_enabled`／`allow_shell_string`／`shell_allowlist` 三枚键今天被**加载期响亮拒绝**（`validateUnwired`，`unwired.go:99-111`），
   理由逐字就是"no shell.exec tool is registered"。**会话腿一旦落地，这三句 `missing` 里至少第一句变成假话**，
   而 `internal/config/**` **既不在票面 `:29` 的允许写面里，也不在它的零字节名单里**（这一格是空白，不是遗漏 owner 就能看到）。
   ⚠ 还有钉子在旁边：`TestEveryLockedSectionKeyIsAccountedFor`（`internal/config/unwired_test.go:311`）——
   新增 🔒 `[risk]` 键不登记就红。⇒ **要么复用 `shell_enabled` 当会话的总开关（则 `unwired.go` 必须改），要么新增一键（则测试＋D36 的 🔒 语义一起跟上）**；两条路都不是"只动 `internal/tools`＋`internal/proc`"。
3. **审批腿**：票面 `:28` 自己禁动 `internal/agent/approval/**`；会话的"每发命令各判一次"能不能在不碰它的情况下成立，本程**未量**（§2 第 2.5 格）。

**⑦ 要 owner 点的那一句**：

> **「批 D34 在 `shell.exec` 行（`PLAN.md:2563`）之后加一行 `shell.session`（L2、capability `shell`、S3、寿命绑任务、不引 PTY），同批镜像进 `SPEC-07:71` 之后；`shell.exec` 那一行与 D33 的 `:2507`、D46 的 `:2662`「保持不变」一字不动；并放开 `internal/config/unwired.go` 里那三行 `shell_*` 的 `missing` 文案（只改文案、不删守卫、不动阈值）。」**

不这么批的代价：走票面原路（扩 `shell.exec` 语义）＝**改写三枚定稿句子**（其中一枚是 D33 安全必修项、一枚写着"保持不变"），
本程按"改写比新增贵"的规矩判它**不推荐**；不批＝两条腿都停在纸面，且 `unwired.go:63` 那句"今天没人读它"会一直是对的。

---

### 4.4 票 164 — 后台任务读不到输出：先定名册，再补输出腿

**① 票面现读**（`.scratch/wisp/issues/164-background-jobs-cannot-be-read-fix-the-roster-then-add-the-output-leg.md`，`Read` 全文，28 行）：
- Status（`:3`）：`立而不派（D34 表里那两枚名存实无，先要一次名册更正＝人工批准，再谈实现）`
- AC#1（`:14`）逐字：**`名册先行：现量 D34 里所有 task.* 名字 ↔ 生产注册表里的实现，列出差集；差集非空就不许进实现格，先出一份"要么补名、要么删名、要么标 DEFERRED 且五字段齐"的裁定表。⚠ 不许悄悄改名（Q-58 那两枚就是"名字与规矩打架、没人裁定"的前例）`**
- AC#3（`:16`）：截断必须带**可找回的指针**，并指名按 `PLAN.md:431` 的既有规矩做（**本程核该行号＝对**，见下）
- AC#5（`:18`）逐字：**`docs/specs/**、internal/risk/**、internal/panel/**、thresholds.go、golden、allowlist.txt、frontend/**、design/** 零字节；D34 那一行若要动＝owner 单独批准、原文落台账`**
- 本票不解决（`:22`–`:23`）：`:23` `不做面板上"看得见后台任务"的界面（那是 frontend/**，owner 委托会话的写面）`
- 前置（`:27`）：`AC#1 那份名册裁定先归我判（编排者），判完再派实现；票 163 之后`

**② 要动的具体行**：`docs/PLAN.md:2561`（唯一权威落点），逐字：

```
$ awk 'NR>=2559 && NR<=2562 {printf "%d\t%s\n", NR, $0}' .scratch/wisp/probes/16x/c1/PLAN.md.0d17647
2559	| **`reminder.set/list/cancel`** | 一次性提醒（**只发通知**） | L1 / L0 / L1 | `notify` | **S4** | ⭐ **对 D12 的部分推翻**，见 16.5.4 |
2560	| `memory.save` / `memory.recall` | 显式记忆 | L1 / L0 | `memory` | S4 | D20 |
2561	| `task.list` / `task.cancel` | 查看/取消任务 | L0 / L1 | — | **S7** | **D31「被阻塞的东西必须可见」的落地工具**（原文只有 UI 要求，没给 Agent 侧的查询手段） |
2562	| `list_tools` | 元工具（D15 两级注入） | L0 | — | **S1** | — |
```
被 AC#3 点名的那枚既有规矩（**票面行号在本锚点成立**）：

```
$ awk 'NR==431 {printf "%d\t%s\n", NR, $0}' .scratch/wisp/probes/16x/c1/PLAN.md.0d17647
431	| **③ 单个工具结果** | **> 4000 token**（约 6000 汉字 / 16KB 文本） | 全文落 `artifacts	ool-output-<id>.txt`；上下文里只留 **头 500 token + 尾 200 token + 总长度 + 文件路径**；模型可用 `fs.read` 按需再读（复用 D10 机制） |
```
（⚠ 那一行里 `` `artifacts	ool-output-<id>.txt` `` **不是我在转述时打错的**：本程按码点现量，`docs/PLAN.md:431` 的原文里就是一个 **U+0009 制表符**顶在 `artifacts` 与 `ool-output` 之间（写的人本意是 `artifacts\tool-output-<id>.txt`，反斜杠被吞成了一次转义）：

```
$ awk 'NR==431' docs/PLAN.md | python -c "import sys;s=sys.stdin.read().rstrip('\n');print([hex(ord(c)) for c in s[66:76]], s.count(chr(92)), s.count(chr(9)))"
['0x66', '0x61', '0x63', '0x74', '0x73', '0x9', '0x6f', '0x6f', '0x6c', '0x2d'] 0 1
                                                             ^^^^=U+0009 TAB   ^反斜杠枚数=0  ^制表符枚数=1
```
⇒ **定稿里被票 164 当判据引用的那枚路径字符串是坏的**。本程**未修**（冻结件、且改它＝人工批准），
**只登记**：它同时是 `PLAN.md` 与 §10 验收读数的一个隐性坑（任何程序照字面读这行会拿到一个带 TAB 的路径）。
工作树与该行的锚点版逐字节相同（`git diff --stat HEAD -- docs/PLAN.md` 空），所以这不是 checkout 改出来的。

**③ 建议的替换文本**：**`:2561` 一字不动**（那是"改写"），改为在它**之后插一行**，另在 §7 登记表加一条 DEFERRED：

```
| **`task.output`** | **读一个后台任务吐了什么**（含"太长怎么续读"） | **L0** | — | **S7** | 164 的输出腿；截断按 `:431` 的既有规矩（头 500＋尾 200＋总长＋路径），**只截不指＝不合格** |
```
§7 登记表（`PLAN.md:1530` 那一族所在）新增一条（**照 `:1530` 的五字段排版**，五字段齐）：

```
| **DEFERRED** | **`task.list` / `task.cancel` 的实现（D34:2561 名存实无）** | 名册在册但生产注册表零实现（两把尺见 16x-c1 证据件 §4.4⑤）；先补名册裁定再谈实现 | 注册表里 `task.list`/`task.cancel` 各有一枚真实现并过契约测试 | 票 163（同一批命令执行面）＋票 164 AC#1 裁定表 | Agent 侧查不到"谁被阻塞"，D31 那条"必须可见"只剩 UI 半条腿 |
```

**④ 新增还是改写**：**两处都是新增**（表加一行＋登记表加一行）。
本程**明确不建议**去改写 `:2561` 那句 备注（"D31『被阻塞的东西必须可见』的落地工具"）——它是那行的**理由**，不是状态声明；
状态更正的正确落点就是 §7 登记表（票面 `:14` 自己给的第三选项"标 DEFERRED 且五字段齐"）。

**⑤ 同源拷贝有几份（两把尺＋第三把）**：
- **尺A（词形）**：同 §4.2⑤ 的 `内置工具权威表` 一查，表格拷贝 **1 份**（`SPEC-07`），本行落点 **`SPEC-07:69`**：
  `$ awk 'NR==69' docs/specs/SPEC-07-tools-and-plugins.md` → `| `task.list` / `task.cancel` | 查看/取消任务 | L0 / L1 | — | **S7** | 「被阻塞的东西必须可见」的 Agent 侧手段 |`
- **尺B（内容串 `task.list`）**：全语料 **4 处**（命令与原始输出在 `.scratch/wisp/probes/16x/c1/copy-census.log` 尺B-2 段）：
  `docs/PLAN.md:2561`（本体）· `docs/PLAN.md:3109`（**D44 的 S7 验收清单里也点这俩名**）· `docs/specs/SPEC-06-security-gatekeeping.md:108` · `docs/specs/SPEC-07-tools-and-plugins.md:69`
  ⇒ **改名或加注必须同批扫到这 4 处**，否则"名册"这件事会以另一种形式回来（`Q-58` 那一族的形状）。
- **尺C（生产侧到底有没有，两把尺）**：

```
$ git grep -nE 'func \([a-zA-Z0-9_ *]+\) Name\(\) string[[:space:]]*\{ return "' 0d17647 -- '*.go' ':!*_test.go' ':!.scratch' | wc -l
6
$ git grep -n '"task\.' 0d17647 -- 'internal/' 'cmd/' ':!*_test.go' ':!.scratch'
0d17647:internal/statemachine/table.go:253:		SideEffects: []string{"task.ctx-keep"},   ← 唯一命中，且不是工具名
```
  ⇒ 生产注册表**只有 6 枚 `fs.*`**（`fs.read/fs.list/fs.write/fs.trash/fs.move/fs.delete`），**`task.*` 一枚都没有** ⇒ 票面 `:4` 那句"`task.list`／`task.cancel` 作为工具在生产注册表里并不存在"**成立**。
  ⚠ 本程第一次跑这把尺时用了**不带 `[[:space:]]*`** 的正则，得 **5**（漏了 `fs.go:188` 那枚被 gofumpt 对齐过的 `fsList`）；**两把尺都跑第二遍才对上 6**——记在这里，免得下一位拿 5 说事。
- **另记（不是拷贝、但会被名册更正牵到的两处注释）**：`docs/reports/pending-and-issues.md:7642`（A308 里那句"名册级更正"已经写过了）· `.scratch/wisp/issues/README.md` 未列 `task.*`。
  按 §2 第 2.11 格的排除理由，**不计进"要同批改的份数"**。

**⑥ 不动 `D34` 能不能落地**：**AC#1（名册裁定）能**，**AC#3（输出腿）不能**。
- 能的那半：裁定表是**编排者手里**的东西（票面 `:27` 就写了"先归我判"），零冻结面文字。
- 不能的那半：**新工具必须有名字**；名字不在 D34 ⇒ 与 `Q-58` 同形（`cmd/wisp/run.go:330-340` 那段注释就是"两个冻结名打架、停在码里等 owner"的现物，本程现读在位）。
- ⚠ **附带条件一枚（本程量到的票面自相矛盾）**：票面 `:18` 令 `docs/specs/**` 零字节，
  而它要动的 `:2561` 有 **1 份 `SPEC-07:69` 逐行拷贝**＋**1 份 SPEC-06:108 内容引用** ⇒
  **不改 `SPEC-07` 就交出一张已知会漂的表**；要改就违背自己的零字节 ⇒ 批语里必须一起放开（同 §4.2⑥）。
- ⚠ 还有第二枚连带：`:3109` 那枚 **D44 的 S7 验收清单**也点了 `task.list`/`task.cancel`。它不是一张拷贝，它是**验收判据**；
  名册若被改（删名/换名），**D44 那一行会跟着失去指代对象** ⇒ 本程把它列为"**必须同批核、不必同批改**"的那一类（现状：不改名 ⇒ 它不受影响）。

**⑦ 要 owner 点的那一句**：

> **「批两枚新增、零改写：D34 在 `:2561` 之后加一行 `task.output`（L0、S7，截断按 `:431`），§7 登记表加一条 `task.list`/`task.cancel` 的 DEFERRED（五字段齐）；`:2561` 原文一字不动；同批把 `SPEC-07:69` 之后镜像那一行；`docs/PLAN.md:3109` 的 S7 验收清单只核不改。」**

不这么批的代价：把 `:2561` 直接**删掉**＝让 `PLAN.md:3109` 的验收清单指到一个不存在的名上（`SPEC-12 §3` 的"任一侧落空＝缺口回报"那一形）；
**不动它**＝那两枚继续"名册在册、注册表零实现"，而票 164 的 AC#1 永远交不出差集为空的读数。

---

### 4.5 票 165 — "先出方案、你点头再动手"那一档

**① 票面现读**（`.scratch/wisp/issues/165-plan-first-then-act-a-brake-the-voice-entry-needs.md`，`Read` 全文，39 行）：
- Status（`:3`）：`立而不派（新增一档＝动 D31/D43/C12 那条状态与审批线，人工批准；且 owner 已定过"改状态机要谨慎"）`
- 硬约束（`:19`–`:21`）：`:19` **`不许新造一套状态机。D43/C12 那张表是冻结的；本票优先复用现有的 AwaitingApproval 那一族，若量到"必须新增状态才成立"，停手上报`**
- AC（`:25`–`:29`）：AC#1 未修码读数／**AC#2 复用形状（不新增状态；必须新增则本票退回转 `Q##`）**／AC#3 批准只批这一批（两发变异）／AC#4 面板不许代答／**AC#5 契约轴**
- AC#5（`:29`）逐字：**`docs/PLAN.md、docs/specs/**、internal/risk/**、internal/panel/**、internal/agent/approval/**、thresholds.go、golden、allowlist.txt、frontend/**、design/** 零字节`**
- 本票不解决（`:33`）：`不做界面（谁显示方案、点哪儿算点头＝frontend/**，owner 委托会话的写面）`

**② 要动哪一行**：**本程把三条候选路各自的落点都量了一遍，结论是"这三条路的落点全在 `:29` 那句零字节名单里"**：

| 候选路 | 必动的东西 | 逐字现读（本锚点） | 在 `:29` 的禁改名单里吗 |
|---|---|---|---|
| 甲：方案是一枚**工具**（走 C18 队列） | `D34` 加一行 | `docs/PLAN.md:2562` 就是那张表的元工具行，新行只能挨着它插 | **在**（`docs/PLAN.md`） |
| 乙：方案是一枚**权限档**（第四档） | `internal/risk/mode.go` 的词表 | `65	modeVocabularyPattern = "ask_every_step\|ask_high_risk\|auto_approve"` | **在**（`internal/risk/**`） |
| 丙：只改**审批层的放行范围**（批完方案就免问后续） | `internal/agent/approval/**` | `internal/agent/approval/batch.go:14` `// BatchMinOps is D45-1's "N >= 3 同质 L1 操作 -> 合并为一个确认".` ＋ `:15` `BatchMinOps = 3`（`:11-12` 逐字：`They are numbers from the contract, not tunables: changing one is a contract change.`） | **在**（`internal/agent/approval/**`） |

⇒ **本票的"哪一行必须改"这一问，正确答案是先改票面**（它的 `:29` 把这枚票唯一可能的三条路全禁了）。

**③ 建议的替换文本**（走甲路，代价与差别见第 ⑥ 问）：

```
| **`plan.present`** | **交回一份"我要干什么"的清单并停下等点** | **L2** | — | S3 | 复用 D43 #17/#23/#24/#25 那一族，**不新增状态**；**批准只批清单里那一批**（清单外的写仍各自弹卡）；**本票因此不减少确认次数**，它买的是"先看意图再看手" |
```

**④ 新增还是改写**：**新增**（表加一行）。
⚠ **但票面自己许诺的收益是"改写"级的**——"一批改动只问一次"这句话若要做实，
必动的就是 **D45 那三条不可逾越的限制**（`PLAN.md:2148`–`:2153`，现读第一条逐字 `**L2 永不进入任何持久授权**（含会话级）`），
外加 `docs/specs/SPEC-06-security-gatekeeping.md:115-120` 那份 D45 拷贝，外加 §7 里那两枚 RESERVED 行（`PLAN.md:1532`、`:1533`）——
**那全是改写已定稿的文字，且 D45 那两枚机制在语料里落在 22 行**（命令与逐枚见本节第 ⑤ 问）。

**⑤ 同源拷贝有几份（两把尺）**：
- **`plan.present` 这个名字**：**0 份**（它还不存在）。两把尺：

```
$ grep -niE "plan.?mode|计划模式|先方案|propose" .scratch/wisp/probes/16x/c1/PLAN.md.0d17647 | wc -l
0        （尺1：内容串，定稿整篇零命中）
$ grep -n "steering" .scratch/wisp/probes/16x/c1/PLAN.md.0d17647
2731	| `[agent]` | `max_rounds`(50 兜底) … `steering_enabled` | `hot` |
         （尺2：唯一那枚"看起来像刹车"的键，全篇只出现 1 次、且从未被定义过它做什么）
```
  ⚠ 第二把尺顺带量到一枚**既存空洞**：`[agent] steering_enabled` 在定稿里**只有名字、没有语义** ⇒ 任何程把它当"就是这一档"来实现＝自行假设（D22 闸门③／`AGENTS.md §2`）。
- **`AwaitingApproval` 那一族**（甲路复用的东西）：**词形尺 5 枚文件、内容尺 18 处**，其中**表格拷贝 1 份**：
  `$ git grep -c "AwaitingApproval" 0d17647 -- docs/PLAN.md 'docs/specs/*' AGENTS.md` → `AGENTS.md:1 · docs/PLAN.md:10 · SPEC-06:1 · SPEC-08:5 · SPEC-12:1`
  转移表本体的拷贝＝`docs/specs/SPEC-08-ui-ball-panel.md:85-130`（它的 #23/#24/#25 在 `:111`/`:112`/`:113`，与 `PLAN.md:3077/3078/3079` 逐行对得上）。
  ⚠ **本程在这枚票上量到一张"两张表已经不一致"的活证**（这一条与 165 直接相关，因为复用的就是它）：

```
$ awk 'NR==3051' .scratch/wisp/probes/16x/c1/PLAN.md.0d17647
3051	20 态（§2）· 40 条转移。**未列出的转移一律非法**（未定义即停，D22 闸门③）。
$ awk 'NR==85' docs/specs/SPEC-08-ui-ball-panel.md
85	## 3. 状态机权威转移表（C12/D43——20 态 42 条转移；未列出的一律非法）
$ git show 0d17647:internal/statemachine/doc.go | awk 'NR>=6 && NR<=7 {printf "%d\t%s\n", NR, $0}'
6	//   - the authoritative D43 transition table (table.go: the 40 frozen rows
7	//     plus SPEC-08 §3's appended #41/#42) and the per-state timeout table
```
  ⇒ 定稿说 **40**、`SPEC-08` 说 **42**（并把 #41/#42 插在 `:119-120`，顺序也乱了）、生产码的 `doc.go` 自己写明"40 frozen rows **plus** SPEC-08's appended #41/#42"。
  第三把尺（枚数）：`$ awk 'NR>=3055 && NR<=3094' 定稿件 | grep -c "^| "` → **40 行**；`$ awk 'NR>=89 && NR<=130' SPEC-08 | grep -c "^| "` → **42 行**。
  ⚠ **这不是本程造的**（D47 那批追加没回写定稿 §16.10 的那张表），但它正是"任何新转移该往第几号插"这个问题的现物 ⇒ **票 165 若要加任何一条转移，它会掉进这 40/42 缝里**。
- **D45 那两机制在语料里的枚数**（丙路要改的东西，**按文件分**）：

```
$ grep -rc "批量聚合\|作用域会话授权" docs/PLAN.md docs/specs/*.md AGENTS.md | grep -v ":0$"
docs/PLAN.md:15
docs/specs/SPEC-01-architecture.md:1
docs/specs/SPEC-02-data-storage.md:1
docs/specs/SPEC-06-security-gatekeeping.md:3
docs/specs/SPEC-10-testing-acceptance.md:1
docs/specs/SPEC-12-roadmap-governance.md:1
$ grep -rn "批量聚合\|作用域会话授权" docs/PLAN.md docs/specs/*.md AGENTS.md | wc -l
22
```
  ⇒ 逐枚清单已在本程取数记录里（`PLAN.md` 的 15 枚行号：`1532 1590 1639 1795 2133 2135 2138 2140 2144 2698 2844 2985 3105 3109 3144`；`SPEC-06` 的 3 枚：`115 117 152`）；
  **丙路若走，这 22 行里"哪几行是真要改"本程未逐枚判**（那要读每一行的语境，已越出普查射程，见 §2 第 2.2 格）。
  ⚠ `AGENTS.md` 那把尺 **0 命中**（`grep -c` 直接给了 `AGENTS.md:0` 后被 `grep -v ":0$"` 滤掉，所以看起来像"漏了一枚文件"，不是漏）。

**⑥ 不动它能不能落地**：**不能落地**（三档里的"不能"，且是最硬那种）——
不是"要改一行才能落地"，是**票面 `:29` 自己把三条路全堵了**：
甲堵在 `docs/PLAN.md`、乙堵在 `internal/risk/**`、丙堵在 `internal/agent/approval/**`。
⇒ **必须先有一句批准（放开其中一条路的具体落点），这枚票才有可派形状。**
附带两条本程量到的实话：
1. **甲路不做 D45 改写 ⇒ 不减少确认次数**：`AGENTS.md` 禁的"面板侧给出允许"这条红线（`PLAN.md:2420` 逐字 `**`approval.decide` 的「允许」决策不接受来自面板的调用。**`）意味着卡片照弹、点还在那儿点；
   ⇒ **票面 `:15` 那句"有了先方案这一档，一批改动只问一次"在甲路下不成立**（本程未实跑，见 §2 第 2.2 格）。
2. **甲路复用的 `AwaitingApproval` 有一格行为至今未定义**，且它是 `AGENTS.md §2` 具名的"未定义即停"项：

```
$ awk 'NR==1530' .scratch/wisp/probes/16x/c1/PLAN.md.0d17647 | cut -c1-150
1530	| **DEFERRED** | **`AwaitingApproval` 遇系统挂起的行为定案**（§16.11 第 10 条） | 本轮识别但未定案。倾向**作废并判拒绝**…| S7 切片卡明确定义，并写进 D43 转移表 | S7 | 挂起恢复后待批准请求的行为**未定义** → 按 D22 闸门③「未定义即停」处理 |
```
   同一条在 `PLAN.md:3228-3231` 与 `docs/specs/SPEC-06-security-gatekeeping.md:110-111` 各写了一次 ⇒ **"方案被批准"这一发如果正好赶上系统挂起，行为没有定义**。
   ⚠ 这一条**不该由 165 顺带解决**（它自己 `:19` 就说了新造状态要停手），但 owner 批甲路时应知道：他批的那张卡挂在一枚已知未定义的行为上。

**⑦ 要 owner 点的那一句（含本程的判语）**：

> **「票 165 我按甲路批：放开两枚落点——`docs/PLAN.md` 的 D34 加一行 `plan.present`（L2、S3、复用 D43 现有 #17/#23/#24/#25，不新增状态）＋同批镜像 `SPEC-07 §3`；`internal/agent/approval/**` 与 `internal/risk/**` 仍然零字节。⇒ 明示代价：**这一档今天只买到"先看意图再看手"，不买到"一批只问一次"**；要后者必须另批 D45 的改写（语料 22 行落点，含 §7 两枚 RESERVED 行）。」**

⚠ 三条本程要更正的票面读数（都带凭据，逐条）：
1. `:15` 引 `PLAN.md:2118` 算的是"**≈33 分钟**"（`500 次阻止窗口 × 每次约 4s`），票面写成"半小时"——**同一笔账，措辞漂移**（凭据：`awk 'NR==2118' 定稿件`）。
2. `:3` 的 Status 说"新增一档＝动 D31/D43/C12" —— **本程量到：走甲路可以完全不动 D43/C12**（复用现成四行）；"必须动状态机"这件事**不成立**，除非要走乙路（第四枚权限档，那才动 `internal/risk/mode.go:65` 的词表）。
3. `:29` 里 `internal/panel/**` 被列零字节 —— **仓规其实明令它不冻结**：`.scratch/wisp/issues/README.md:44` 逐字 `**不冻结**（别扩大解释）：`internal/panel/` 与 `cmd/wisp/` 的 **Go 侧接线**`（`docs/reports/HANDOVER.md:545` 同句）。
   ⇒ 这一条不是缺陷（票面自愿收窄射程是合法的），但**owner 批的时候别以为"解冻 panel"需要他开口**——需要他开口的是 `docs/PLAN.md` 那一行。

---

### 4.6 票 168 — 内置浏览器：两件必须先立

**① 票面现读**（`.scratch/wisp/issues/168-built-in-browser-as-a-disable-able-plugin-two-rules-to-set-first.md`，`Read` 全文，44 行）：
- Status（`:3`）：`立而不派（前置＝票 50（Tier-1 清单插件）＋票 51（Tier-2 goja）都没做完，插件地基不存在；且它要动 C17/C24 那两张未定案的白名单＝人工批准）`
- 推荐（`:12`）：`A 做成可关的插件；B 只在宿主里留一个受控的第二视图，不许插件自己塞页面进来`
- 两件必须先立（`:21`–`:22`）：逛到的内容必须打"外部内容"标记（与抓取同枚标记）／一次浏览的会话不得跨任务常驻
- AC（`:26`–`:32`）：**AC#1 地基前置证明（票 50/51 未结案 ⇒ 不许开工，缺任何一枚停手上报）**／AC#2 注册位第二名当场失败／AC#3 两引擎同一接口／AC#4 外部内容标记正反两向／AC#5 面板边界零越界（`:30` 逐字含 `只加"显示/导航"这一类，不加"决定"这一类`）／AC#6 关掉要真关（对着 D32 那两个数）／**AC#7 契约轴（`:32`）**
- AC#7（`:32`）逐字：**`docs/PLAN.md、docs/specs/**、internal/risk/**、internal/panel/**、internal/agent/approval/**、thresholds.go、golden、allowlist.txt、frontend/**、design/** 零字节`**
- 本票不解决（`:36`–`:39`）：`web.search` 走哪条路不答／`:37` `不做"给人看的浏览器界面"（界面腿 owner 委托）`／不接 MCP／`:39` **不许把浏览器做成受门控的内置工具**

**② 要动哪一行**：**这枚票的批准不是"改哪一行字"，是"两张尚未定案的白名单要不要现在就定稿"**（与 160 同型）。它自己的 AC#7 也把 `docs/specs/**` 禁死了，而那两张表**正好只在 specs 里**：

```
$ grep -n "方法白名单" docs/specs/*.md
docs/specs/SPEC-08-ui-ball-panel.md:156:### 5.2 C17 PanelBridge 方法白名单【SPEC 提案，S5 定稿走契约批准】
docs/specs/SPEC-08-ui-ball-panel.md:235:- C17：correlationId 路由测试（并发 10 个请求乱序回复）、方法白名单拒绝、resync 后前端状态与
docs/specs/SPEC-06-security-gatekeeping.md:131:3. **服务端二次授权**：…；PanelBridge 方法白名单 +
docs/specs/SPEC-12-roadmap-governance.md:48:- C24 GojaHostAPI 初始集与 C17 方法白名单定稿（S7/S5 切片卡批准）
```
⇒ **方法清单本体只在 `SPEC-08:161-174` 一处**（14 行方法表，`panel.resync` 起、事件推送止）；
⇒ **C24 的初始入口集只在 `SPEC-07:117-120` 一处**，它自己那行就写着"【SPEC 提案，S7 定稿走契约批准】"：

```
$ awk 'NR>=115 && NR<=120 {printf "%d\t%s\n", NR, $0}' docs/specs/SPEC-07-tools-and-plugins.md
115	- 所有能力经 host 注入（goja 无 `fetch`/定时器/IO）；**C24 `GojaHostAPI` 是 JS 侧唯一宿主入口
116	  全集**，未列出的入口即契约违规（字符串扫描可验）。
117	  【SPEC 提案，S7 定稿走契约批准】初始入口集：
118	  `wisp.http.request`（cap: net）· `wisp.fs.read|write`（cap: fs.read|fs.write）·
119	  `wisp.clipboard.read|write` · `wisp.notify.send` · `wisp.sysinfo.get` ·
120	  `wisp.state.get|set`（插件自身 KV，无 capability——插件自己的数据）。
```
（这两处未定案的身份，`AGENTS.md §2` 未定义即停清单里点得名：`- ` + '`C24 GojaHostAPI`' + ` 初始集与 ` + '`C17`' + ` 方法白名单定稿（S7/S5 切片卡批准）`。）

**③ 建议的替换文本**（**只在 owner 真要现在定 B 路时才贴**；本程的推荐是**现在不定、只做第一问那两枚标记**，见第 ⑥ 问）：
- 若定 B 路：`docs/specs/SPEC-08:156` 那枚标题去掉未定案方括号 →
  `### 5.2 C17 PanelBridge 方法白名单（2026-09-26 owner 定稿；浏览器视图那两枚另批）`，
  并在方法表（`:174` 之前）**加两行**：

```
| `browser.view.open` | Go→前端推送 | —（只开一个受控第二视图；**不给出任何允许**） |
| `browser.view.navigate` | Go→前端推送 | —（导航目标必须是任务内已批准的 URL；**永不扩"决定"这一类**） |
```
- 若只立"两件规矩"（本程推荐）：**零文字改动**，两枚标记都是**行为要求**，落点在实现里：
  外部内容标记复用 `internal/risk/provenance.go` 那族源标记（`SrcWebFetch = "web.fetch"`，现读 `:79`），
  会话归属复用票 160 的句柄形状（票面 `:22` 自己就是这么指的）。

**④ 新增还是改写**：走 B 路＝**两处都是改写＋两处新增**（标题那句"提案/未定稿"是要被**改掉**的已定稿措辞，方法表加两行是新增）；
走"只立两枚标记"路＝**零改写零新增**（不动冻结件）。⚠ **这个差价是本枚最该让 owner 看见的东西**：
168 唯一"必须现在批"的那一格（第 ⑥ 问）其实**不需要改任何文字**，需要的是**排期判断**。

**⑤ 同源拷贝有几份（两把尺）**：
- **C17 方法清单**：**1 份**（就在 `SPEC-08:161-174`，它就是本体）。词形尺给出 6 处**只谈"要有白名单"而不列方法**的指涉：
  `PLAN.md:2432`、`PLAN.md:2963`、`PLAN.md:3107`、`SPEC-06:131`、`SPEC-12:48`、`AGENTS.md:74`
  ⇒ **改方法表不会牵连那 6 处**（它们不抄内容）。
- **C24 初始入口集**：**1 份**（`SPEC-07:117-120` 本体，`grep -rn "wisp\.http\.request" docs/` → 仅此一处）；
  指涉它的：`PLAN.md:1374`（C24 契约行，只说"全集/未列出即违规"，不列名字）· `PLAN.md:3109`（S7 验收）· `SPEC-12:48` · `AGENTS.md:74`。
- ⚠ **别把上面这两组数并成"拷贝很多"**：**内容尺给 1＋1**，**词形尺给 6＋4 枚"只指名不抄内容"的引用**。

**⑥ 不动它能不能落地**：**不能，且理由与冻结文字无关（本程两把尺现量）**：

```
$ ls .scratch/wisp/issues/ | grep -E "^0?5[01]-"
50-tier1-manifest-plugins.md
51-tier2-goja.md
$ ls .scratch/wisp/issues/ | grep -cE '^(50|51)-.*-done$'
0        （grep rc=1＝零命中）
$ grep -n "Status" .scratch/wisp/issues/50-tier1-manifest-plugins.md .scratch/wisp/issues/51-tier2-goja.md
.scratch/wisp/issues/50-tier1-manifest-plugins.md:3:**Status:** ready-for-agent
.scratch/wisp/issues/51-tier2-goja.md:3:**Status:** ready-for-agent
```
⇒ 两枚都不带 `-done` 后缀（`AGENTS.md §1.5`：那后缀是防重领的唯一键），Status 也都不是 done ⇒ **票 168 的 AC#1（`:26`）判"不许开工"直接成立**。
⇒ 它还排在票 160 之后（`:43` 前置条件写死"票 160 结案"），而 160 正是 §4.1 里那枚"要 owner 批开地"的票。
⇒ **本程对 owner 的建议：168 这一枚不要现在批文字。** 它今天唯一能批的是"**先立那两枚规矩**"（票面 `:21-22`），
而那两枚**不需要动任何冻结件**；方法白名单（C17/C24）定稿是 **S5/S7 切片卡**该做的事（`:48` 与 `AGENTS.md §2` 都是这么写的），
现在为了一枚地基还没有的浏览器提前定稿，正是本仓 `Q-52` 那格"为还会变的东西签冻结契约"的形状。

**⑦ 要 owner 点的那一句**：

> **「票 168 我这轮只批那一半不需要动冻结件的：把票面 `:21-22` 那两件（逛到的内容必须打与抓取同一枚外部内容标记／一次浏览的会话不得跨任务常驻）钉成实现要求，`docs/specs/**` 与 `docs/PLAN.md` 一字不动；C17 方法白名单与 C24 初始集仍按 S5/S7 切片卡定稿，**不提前为浏览器开洞**；本票继续立而不派，等票 50／51 双双 `-done` ＋票 160 结案。」**

不这么批的代价：若现在就定 C17/B 路那两枚方法名 ⇒ **在一枚依赖的插件地基（票 50/51）尚不存在、且界面腿归别人**的情况下把白名单钉死，
将来 S5/S7 切片卡真要定稿时只能改写它；若连"两枚标记"都不批 ⇒ `:21` 那条正面漏口留在（dsh 那形"抓的打了标记、逛的什么都没打"，票面 `:4` 已把它记成反面教材）。

---

（下节：§4.7 owner 那句"这六张里有没有要写 `frontend/**` 的"。）
