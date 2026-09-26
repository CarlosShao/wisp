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

（下节：§4.3 票 163。）
