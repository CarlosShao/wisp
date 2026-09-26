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

- **尺A（词形）**：`grep -rnE "逐字|照抄|verbatim|唯一来源|权威表" docs/PLAN.md docs/specs AGENTS.md` → 全语料 33 行命中（原始输出在 `.scratch/wisp/probes/16x/c1/copy-census.log` 尺A-1 段），
  其中**没有一行**声称"逐字抄自 `provenance.go` 的头注释" ⇒ 词形这把尺给出：**0 份自称同源的拷贝**。
- **尺B（内容串）**：`grep -rn "OpenScope at task start" docs/ internal/ cmd/` 与 `git grep -n "OpenScope" <锚> -- '*.go' ':!*_test.go' ':!.scratch'` →
  文档侧 0 命中；**代码侧 2 枚生产调用点**（它们是**受这次签名改动影响的调用者**，不是"自称逐字抄"的拷贝）：

```
$ git grep -l "OpenScope" 0d17647 -- '*.go' ':!*_test.go' ':!.scratch' | grep -v internal/risk/provenance.go
0d17647...:cmd/wisp/panel_assets.go
0d17647...:internal/tools/bridge.go
```

  逐枚现读（本锚点）：
  - `internal/tools/bridge.go:633-644` `func (b *Bridge) OpenTask(...)`，其中 `:638-642` 就是票 160 进度段说的那两套并行记账：
    `open := b.scopes[taskID]` / `b.scopes[taskID] = true` 与 `:642 b.prov.OpenScope(taskID)` 分列两行（**本程读码自取，与票面 `:46-47` 同读**）；
  - `cmd/wisp/panel_assets.go:232 prov.OpenScope(taintSourceScopeID)`，且 `:162-166` 的注释自陈它是**故意不关**的：
    `a CLI probe has no task, so it names one / and closes nothing - the engine lives and dies inside this process.`（逐字）
    ⇒ **这一枚不是遗漏，是与新形状直接冲突的一形**：句柄自带关闭动作之后，"开一个、不关、随进程死"要么显式写成"故意丢弃句柄"，要么被类型系统挡住。
- ⚠ **两把尺都数过，才允许写这句**：`provenance.go` 头注释的**同源拷贝＝0 份**；**受影响的真实调用者＝2 枚**（这是两件不同的事，别并成一句）。

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

（§4.2 起逐枚续写；§3 六行表在六枚全部落地后一次给出。）
