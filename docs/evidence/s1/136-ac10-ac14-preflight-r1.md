# 136 — `AC#10`／`AC#14` 判据预做（preflight）**r1**

> **这份文件是什么**：给票 136 的 `AC#10`、`AC#14` 两格**预做判据与排程约束**，用来决定
> "这两格能不能现在派实现程、按什么顺序派"。**它不是裁决表、不翻任何勾、不裁任何格成立与否。**
> **这份文件不是什么**：它**没有跑过任何一枚测试／构建／`wisp`**。所以本文**没有一格够得上〔独立复现〕**：
> 全部落 **[日志＋归档抽验]**（＝盘上的码、已提交的证据文件、已提交的票面与 workflow），
> 凡我只从别程的文字读到、没有自己在盘上核到的，一律标 **[仅自述不背书]** 并点名是谁的自述。
> **简报（编排者那段导引）里不成立的前提逐条列在 §6，四条，全带凭据。**

---

## §0 锚点与测量口径

| 项 | 读数 |
| --- | --- |
| 本程自量锚点 | `git rev-parse --short HEAD` ⇒ **`7cf8075`**（`2026-09-24 18:12:54 +0800`） |
| 同一次会话内 HEAD 漂过 | 第一次取数时是 **`2988e07`**（18:11:30）⇒ 中间落了 `7cf8075`（只改 `docs/reports/pending-and-issues.md`）。**本文所有 file:line 在 `7cf8075` 上复量**，并用 `git diff 2988e07..HEAD -- cmd/wisp internal/observe internal/proc scripts .github` **输出为空** 证明两版同立（`diff-empty-rc=0`） |
| 被量文件的清洁度 | `git status --porcelain -- cmd/wisp internal/observe .scratch/wisp/issues/136* docs/evidence/s1` ⇒ **零输出**（`clean-rc=0`）⇒ 盘上内容＝HEAD 内容 |
| 工作树里别人的东西 | `git status --porcelain` 整树只有 **owner 的 `design/**` 16 枚未提交删除 ＋ 两枚未跟踪目录**（`design/doubao/`、`design/old/`）⇒ **一枚没碰、没还原、没代提交** |
| 取数时刻 | `date` ⇒ `Thu Sep 24 18:14:42 CST 2026`（本机） |
| 引用他程行号的做法 | 票面／别程给的行号**全部自己 `grep -n` 复量**（结果见 §1.1、§3），没有一条照抄 |

**本文只引用变量名与文件名，不含任何凭据值**（本程也没有读到过任何凭据）。

---

## §1 问一：`AC#10` 那一发端到端到底要跑什么 **[日志＋归档抽验]**

### 1.1 票面那条链，逐段对上盘上行号（锚 `7cf8075`）

票面那一格在 `.scratch/wisp/issues/136-...md:149-156`。它点名的三处**一枚都没漂**：

| 票面说 | 盘上现量（`7cf8075`） | 结论 |
| --- | --- | --- |
| 采样在 `cmd/wisp/slo_windows.go:374` | `:374` `rep, err := sampler.SampleState(ctx, observe.SLOState(state), interval,`（在 `runOutOfTree`，函数体起 `:348`） | **成立**；⚠ 但**第二枚采样点在 `:423`**（`runSubject`，函数体起 `:397`）——同一发 `wisp slo` 真跑时会**同时**跑这两处，票面只点了第一处 |
| `run.Pass = rep.Pass && observer.Pass` 在 `:386` | `:386` 逐字相同 | **成立** |
| `!run.Pass ⇒ exitCode=1` 在 `:323-325` | `:323` `if !run.Pass {` / `:324` `exitCode = 1` / `:325` `}`；另 `:277` `exitCode := 0`、`:326` `return exitCode` | **成立** |
| 守卫本体 `internal/observe/sampler.go:332-340`（票面 `AC#1`／`AC#14` 都引它） | `:332` `if len(rep.Samples) == 0 {`，门行 append 到 `:339`，`:341` `rep.Pass = true`、`:342-346` 只按 `Gate && !Pass` 折总布尔 | **未漂**（AC#12 那次生产码改动全在 `:431`/`:462` 之后，动不到它） |
| `AC#14` 说 `StateReport` 的门行在 `sampler.go:189`、构造在 `:331`/`:335` | `:189` `Verdicts []Verdict \`json:"verdicts"\``、`:190` `Pass bool // all Gate verdicts pass`、`:331` `rep.Verdicts = buildVerdicts(st, *rep)`、`:335` append | **三处全成立** |
| `AC#10` 追加问② 说 `cmd/wisp/slo_windows.go:623` 是合成失败报告 | `:623` `return &observe.SettleReport{TargetState: observe.SLOSleeping, Pass: false}`（`runSettle` 起 `:599`） | **成立** |

### 1.2 判据① 那条命令**照字面跑不出读数**（一票面缺陷）

票面 `:153` 写：`wisp slo -seconds 0.05` 真跑 ⇒ 仍交 ≥1 枚真样。**盘上不成立**：
`cmd/wisp/slo_windows.go:224` 逐字是 `wisp slo: -state is required (or use -settle)`，
其上 `:223` 就是 `if !*settle && *state == "" {` ⇒ 不带 `-state` 的这发**只拿得到 exit 2**，取不到样。

真正的原发在**已提交的证据**里：`docs/evidence/s1/134-adversarial-acceptance.md:538`（FD 行）逐字是
`wisp slo -state Sleeping -seconds 0.05 -interval-ms 250`，读数为 `rc=0 pass=true`、
`samples=1 sample_errors=0 mem_median=4476928`、`duration_sec=0.2912717`。
⚠ 那一发用的二进制是**工作树里那枚 `build/wisp.exe`（`27881312 B`，mtime `09-21 15:06`，本程 `ls -la`＋`date -r` 现量仍在盘上）**，
**早于**它当时被验的版本，也早于今天落地的 `5c1529a`（AC#12 生产码，09-24 08:53）⇒
**`mem_median=4476928` 这枚参照值不属于任何一枚当前树上的二进制**。谁把 `AC#10` 派出去，判据① 的参照必须换成"重建后再量"，
否则又是一枚"引读数没带锚点"的账（今天已为此作废过 `A175`）。

### 1.3 `wisp slo` 在这台机上**开不开六态窗口**？**取不取真读数**？

**窗口：一枚都不开。**凭据：这条命令的全部启动面就是 `proc.Boot(env)`（调用在 `cmd/wisp/slo_windows.go:248`），
而 `internal/proc/boot_windows.go:57-113` 的 `Boot` 只做四件事——env 自检、`DefaultLayout`、`OpenJobScope()`（Job Object 句柄）、
`AcquireSingleInstance`（`:94` 的 `if rt.Layout.MutexEnabled {` ⇒ test env 下为假即跳过），加一次 registry 花名册自检；
整个 `internal/proc` 包里 `webview/panel/ball/CreateWindow` 只出现在注释与 shutdown 钩子的**字段名**上
（`internal/proc/shutdown.go:55,82,86` 那三行是 `DestroyPanel` 一类**可为 nil 的钩子**），没有任何一处创建窗口。
`wisp slo` 也**没有**一次跑六态：一次调用只采**一个** `-state`（票面 :41 那句 "one state per run" 与 `:224` 的强制校验同向）。
**"六态"这件事在脚本层**：`scripts/slo-check.ps1:316` 是 `$states = @('Sleeping', 'Warm')`，
`:318` 的 full 子集才是六枚 `@('Sleeping','Armed','Warm','Conversation','PanelOpen','WorkPeak')`，
`:325` 逐枚 `& $WispExe slo -state $state ...` ⇒ **六次独立进程、六次独立启动**，不是六枚窗口。

**读数：会，真的取，而且是全链路真数。** `:374` 走 `proc.NewExternalSampler(subject.pid)`（内核 per-PID 计数器）、
`runOutOfTree` 还会 `startSubject` **另起两枚子 `wisp.exe`**（`:451`，被 Job 的 `KILL_ON_JOB_CLOSE` 兜住），
默认 `-settle-ms 2000`（`:201`）意味着每发之前先 GC＋`ReleaseMemory`＋2 秒静置；`runSubject` 那枚被测子进程还要睡满
`window + subjectIdleSlack`（`:416`，常量在 `:77`＝`2 * time.Minute`，`:417` 才是它的 fail-closed 报错）。
⇒ **本程不跑它，但盘上的形状已经说明它是"负载＋墙钟时长"两样都占的读数。**

### 1.4 关键排程判定：**这一格必须等编队空**，三条独立凭据

1. **本仓自己已经把这判过一次性**：`scripts/slo-check.ps1:153-155` 的争用探测名单逐字包含
   `go.exe`／`gofmt.exe`／`cgo.exe`／`compile.exe`／`asm.exe`／`link.exe`／`gcc.exe`／`wisp.exe`／`staticcheck.exe`，
   命中即 `:250-310` 那一整段 **"NO CONCLUSION (machine-contended)"**：不出数、不算红、也不放行。
   ⇒ **AC#10 的实现程必然要 `go build`／`go test`；同一时刻任何真取样在本仓口径下都是"无效样"。**
2. **CI 那半台机就是这台机**：`slo-full` 那枚 job 是 `.github/workflows/ci.yml:537-538` `runs-on: [self-hosted, wisp-slo]`，
   步在 `:563`（`-Subset full -SecondsPerState 6`），触发面是 `ci.yml:12-13` `push: branches:[main,dev]`；
   脚本自己的注释（`scripts/slo-check.ps1:126`）逐字写着
   "slo-full runs on a self-hosted runner on the **SAME 6C12T laptop** the review fleet compiles on"。
   ⇒ **AC#10 取数期间任何一次推送，都能在同一台机上并发开一发六态取样。**（`slo-smoke` 那枚在托管 `windows-latest`，`ci.yml:479-480`＋`:500`，不抢本机。）
3. **今天仍有一枚在按包取 `cmd/wisp` 的读数**（见 §2.1），它的名册仪器就是"逐名 101 枚两两 `comm -3` 为空"——
   插一枚新测试进 `cmd/wisp` 会**改变它的分母**，不是改变它的文件。

---

## §2 问二：争用面（落到**文件**，不是包）**[日志＋归档抽验]**

### 2.1 编队现量（四条一起看，不只看票号）

| 程 | 盘面证据 | 判定 |
| --- | --- | --- |
| `acceptor-ticket128-ac4-r1` | 它的分节 commit 一枚一枚到 **`9c8eb3f` 18:10**（§7），`87339ba` 18:06 那节自述跑的是 `-count=2 -v ./cmd/wisp/` 两形；编排者 `7cf8075`（18:12:54）的 commit message ＋ 台账 `A178` 逐字写着"**仍在飞的 `acceptor-ticket128-ac4-r1` 的分节……总判还没写**" | **在飞** |
| `acceptor-ticket119-ac7-r1` | 末枚 `91be4ae` 18:00；`A176`（commit `9af14f5` 18:09）记"七格全有非实现者表"，`0bfd322` 18:07 已把票 119 翻 `-done` | **已交件** |
| `auditor-ticket140-static`／`140-static-inventory-r1` | 三枚 `7c81c77` 18:05／`b33edb0` 18:07／`29d8938` 18:08；`A177`（`2988e07` 18:11）称"静态盘点交回" | **已交件** |
| 票 133 修方（`cmd/wisp` 的原占用者） | `cmd/wisp/leg_dispatch_gate_133_test.go` 最后一次由 **`fa35557` 09:56** 改动，那枚 commit 的标题是 **票 135 AC#8 M-G**；票 133 自身的最后一次提交停在 `fbf420c` 09-23 22:23 | **`cmd/wisp` 现在没有写者**；票 133 面本身**未结**（`133-*.md`：`[x]` 1 枚／`[ ]` 5 枚） |

⇒ 简报的争用前提**成立**（"一枚在飞、两枚已交件"），但**"133 还在飞"这半句今天已不成立**——见 §6(P5)。

### 2.2 文件级冲突表（AC#10 实现程的写集 × 编队）

`AC#10` 实现程**会**写的东西，按今天盘上形状逐枚列：

| AC#10 会碰的文件 | 谁也在碰 | 判定 |
| --- | --- | --- |
| `docs/evidence/s1/136-ac10-*.md` | 只有它自己 | 无冲突 |
| `.scratch/wisp/issues/136-...md`（Progress log 追加） | **票 136 的每一枚实现程都往这里追加**（现量：`:298`／`:316`／`:342`／`:363`／`:370` 各是一条 worker 自述），编排者也往同一枚文件写 | **同文件**——若 `AC#10`/`AC#14`/`AC#15` 同批并行，**三枚程共享同一枚票面文件** ⇒ 按"同文件就判串行"这一条，**票面追加只能有一枚做，或逐格串行**（`A178` 那族"我的字会被它的 pathspec 卷走"是同一枚病） |
| `cmd/wisp/<新 e2e 测试>.go`（若判据落在 `go test` 面） | **无人写**，但 `acceptor-ticket128-ac4-r1` 正在把 `./cmd/wisp/` 的**逐名名册**当仪器（`87339ba`：去重名册各 101 名、两两 `comm -3` 全空） | **文件级无冲突；包级名册有冲突** ⇒ 见 2.3 |
| `cmd/wisp/slo_windows.go`（若为实现"追加问② 的两形可区分"而加显式标记） | 生产码；最后改动是 `ed18727` 09-23 18:17，现无写者 | **文件级无冲突**，但 `AC#10` 那一格**没给任何具名解冻**（票面 :149-156 通篇只说"要一发端到端仪器"）⇒ **动它属越权，须先交回编排者**（对照 `AC#12` 那格 `:192-198`：解冻是按段落具名给的） |
| `internal/observe/sampler.go:332-340`（判据② 的变异面） | `AC#14` 会**写**同一枚文件（`SettleReport` 声明 `:431-456`、`CheckSettle` `:462-521`） | AC#10 只在**仓外快照**上落变异、不提交（本仓口径），所以**不构成文件级写冲突**；⚠ 但 `AC#14` 若先落地，`AC#10` 读到的 settle 出线形状就变了 ⇒ 归到 §3.5 |

**答简报的那句"落到哪几枚文件"**：**与 `acceptor-ticket128-ac4-r1` 之间没有一枚共同文件**（它整程只写 `docs/evidence/s1/128-ac4-r1-acceptance.md`，逐枚 `git show --name-only` 现量只带那一条路径）。
⇒ **按"同文件判串行"这把尺，两枚可以并行**；**但那把尺在这里不够用**，理由在 2.3。

### 2.3 一把尺会漏的那一层：**包级名册**与**取样机**

- `acceptor-ticket128-ac4-r1` 的判据仪器是**逐名名册差集**（它自己 §4 的口径）。`AC#10` 只要往 `cmd/wisp/` **新增一枚 `_test.go`**，
  它的下一发读数就会多出名字 ⇒ 在它交件前落任何 `cmd/wisp` 测试，等于**改别人的分母**。这与"文件级不冲突"不矛盾，**是第二条独立的闸门**。
- 取数与取数争：`AC#10` 要跑 `wisp slo`（§1.3 已证＝真读数），对方正在 `./cmd/wisp/` 上反复 `-count=2 -v`（＝`go.exe`＋`compile.exe`＋`link.exe` 在飞）。
  本仓对"同一时刻既取样又编译"的处置**已经写死**在 `slo-check.ps1:153` 的名单里——**拒绝出数**。
  ⇒ 那条口径虽然只管 ps1 那一步，**它对 `AC#10` 的约束力是"标准"而不是"步骤"**：拿真数当判据的格子，不能在被编译的机器上取数。

**判定：`AC#10` 的实现程现在不可派。**两条理由各能单独挡住：① 128 终裁程在飞（名册分母＋取样机）；② 它的判据② 需要先由编排者改写（§4.1）。

---

## §3 问三：`AC#14` 的落点，与"要不要动生产码"**[日志＋归档抽验]**

### 3.1 `internal/observe/**_test.go` 现量：**13 枚**（`ls internal/observe/*_test.go | wc -l` ⇒ 13）

`clock_test.go`／`diagnostics_test.go`／`earlylog_130_test.go`／`errors_test.go`／`goroutine_test.go`／`logging_test.go`／
`nobarego_test.go`／`observer_cost_test.go`／`sampler_faketree_guard_136_test.go`／`sampler_settle_coverage_136_test.go`／
`sampler_settle_zerosample_136_test.go`／`sampler_test.go`／`sampler_zerosample_136_test.go`
——**13 枚全都不带 `//go:build`**（逐枚 `head -3` 扫过）⇒ 它们在 linux 容器里有分母、在 `internal/observe` 的 core scope 里跑。

### 3.2 先纠一枚符号名：**`Settle`／`checkSettle` 在盘上不存在**

票面 `:278`（AC#15 的地界句）与简报都写"`sampler.go` 的 **`Settle`／`checkSettle`** 语义"。
现量：`grep -rn "func .*Settle" internal/observe/*.go` ⇒ 生产码只有一枚 **`(*Sampler).CheckSettle`，在 `internal/observe/sampler.go:462`**；
其余命中全在测试里（`settleSUT` `sampler_settle_coverage_136_test.go:89`、`settleWire` `:108`、以及 `TestCheckSettle*` 五枚）。
cmd 侧另有 `runSettle`（`cmd/wisp/slo_windows.go:599`）。⇒ **今后引这一族请写 `CheckSettle`（`sampler.go:462`）／`runSettle`（`slo_windows.go:599`）**。

### 3.3 "语义可以不动"成不成立：**对 `AC#15` 成立，对 `AC#14` 不成立**

`AC#14` 的判据①（票面 `:243-244`）要求 `SettleReport` 出线带上**与 `sampling` 同形的门行**，
且"部分未测"这一形**由那一行自己说出不 pass**。盘上事实：

| 事实 | 凭据 |
| --- | --- |
| `SettleReport` **没有** `Verdicts` 字段 | 声明块 `internal/observe/sampler.go:431-456`，逐枚字段是 `target_state/started_at/peak_bytes/cap_bytes/free_os_memory_count/free_os_memory_requested/back_within_cap_ms/elapsed_ms/final_bytes/samples/sample_errors/last_sample_error/pass`——**无 `verdicts`** |
| 隔壁 `StateReport` 有，且总布尔**就是按门行折出来的** | `sampler.go:189-190`（`Verdicts` ＋ `Pass bool // all Gate verdicts pass`）、`:331`（`buildVerdicts`）、`:341-346`（只把 `Gate && !Pass` 汇成 `Pass=false`） |
| settle 的 `Pass` 是**三个布尔的与**，与门行无关 | `sampler.go:517-520`：`backInTime := rep.BackWithinCapMS >= 0 && ...`、`memOK := rep.FinalBytes <= capBytes`、`releaseOK := released && rep.FreeOSMemoryCount > 0`、`rep.Pass = memOK && backInTime && releaseOK` |

⇒ **"补一枚门行"在这枚结构体上无法只用测试实现**：字段要进 `:431-456` 的声明块、行要有人构造（`:462-521` 体内）。
**`AC#14` 天然要动生产码 `internal/observe/sampler.go`**，而票面 `:250-251` 明写"**具名解冻本票面暂不给**：派单时由编排者按'这条判据要落盘必须碰哪几段'现划，
**不许照抄 AC#12 那两段的坐标**"。
⇒ **`AC#14` 不可派的状态不是"排队中"，是"缺一张授权"**——这一条要编排者自己补，实现程无从自救（`AC#12` 那次它按判据必要插了 `:442` 声明区、结果越界报回，见 `:193-198`；同一枚坑今天会再踩一次，因为**声明区在 `AC#14` 里不是"顺带"而是主靶**）。

顺带一条给终裁表的先兆（**[日志＋归档抽验]**，凭据是测试原文）：`AC#14` 落地会**撞一枚已成立的断言**——
`sampler_settle_coverage_136_test.go:208-210` 现在逐字要求
`if !rep.Pass { t.Fatalf("disclosure leg, not a verdict leg: this window still passes, ...") }`，
即"一半读数失败的那一形**仍然 pass**"。`AC#14` 要的正是**把那一形变成会说 not-pass 的门**。
⇒ 两格判据方向相反，**必须显式处置**（要么 `AC#14` 改这条断言并在交件里点名"断言被动了、理由是票面 :243-244"，要么门行做成非 Gate 的记录行、于是 `AC#14` 判据① 的"不许只靠 `pass=false` 一个总布尔代答"落空）。
不要让下一位验收方**用"实现方改了既有断言"这一条把它判成放水**——这一条现在就该写在派单里。

### 3.4 那一族的"停手上报线"三查（票面 `:247-249` 列了三处）

| 停手线 | 现量 | 成立吗 |
| --- | --- | --- |
| `scripts/slo-check.ps1` 要跟着改才过 | 它对 settle 只读两样：`:344` 起的那段里 `.pass` 与 `:353` `$settleReport.settle.free_os_memory_count`；**没有任何严格解析**（全仓 `DisallowUnknownFields` 只有 `internal/config/parse.go:72,123` 两处，都不在这条链上）⇒ 加字段不弄坏它 | **不触发**（加字段这一支） |
| 任何 golden 跟着改 | `git ls-files | grep -i golden` 名册里的 golden **全是 LLM 的 `.sse`**（`internal/agent/testdata/golden/*`、`internal/llm/testdata/golden/*`）；带 `back_within_cap_ms`/`free_os_memory_count` 的文件全在 **`build/`**，而 `build/` 被 `.gitignore:14` 忽略（`git check-ignore -v build/settle.json` ⇒ `.gitignore:14: build/`）⇒ **仓里没有 SLO 报告的 golden** | **不触发** |
| `SPEC-02 §3`（schema 契约级、不得自行增删字段）跟着改 | `docs/specs/SPEC-02-data-storage.md:22` 的 §3 标题逐字是"Schema（契约级，不得自行增删字段…）"，其下 `:23-` 起是 **`CREATE TABLE schema_meta / profile / …` 的 SQLite DDL**——**讲的是库表，不是 SLO 报告的出线形状** | **引用不成立**：`AC#14` 的出线改动**落不进 `SPEC-02 §3` 的射程**。（是不是**另一处**契约在管 SLO 报告形状，本程没查，登记在 §7） |

⇒ 三处里两处**不会**逼停手、第三处**引错了文件**。**`AC#14` 真正的停手面票面反而没写**：门行若按 `StateReport` 同形做成 `Gate:true`，
`wisp slo -settle` 在"部分测量"的真窗口上就会从 exit 0 变 exit 1，而 `slo-check.ps1:381` `$allPass = (... -and $settlePass)`、`:396` `if (-not $allPass) { exit 1 }`
⇒ **这会把颜色打到 CI 的 `slo-smoke`（托管，每 PR/推送）与 `slo-full`（本机 self-hosted）上**——那是**门禁行为变更**，比"动 ps1 一行"更该人工批准。
**建议派单前先把这一支问编排者，而不是等实现程撞上去再报回。**

### 3.5 `AC#14` 与 `AC#10`／`AC#15` 的共用文件

- **与 `AC#10`：没有一枚共用文件**（`AC#10` 若按票面只做端到端仪器，写集在 `cmd/wisp/**` ＋自己的证据；`AC#14` 写集在 `internal/observe/**`）。
  ⚠ 唯一的**真实耦合**：票面 `:252` 说"`AC#14` 排在 `AC#10` 之后（`AC#10` 要 `cmd/wisp` 空出来，**而本格读它的出线**）"。
  盘上量：`AC#14` 读的是 `observe.CheckSettle` 的出线；**`cmd/wisp` 里没有任何一处读 `Verdicts`**
  （`grep -rn "Verdicts" cmd/wisp/*.go` 只命中 `providers_test.go` 的 `probeVerdicts`，是 `wisp providers probe` 那一族，**与 SLO 无关**）。
  ⇒ **"本格读它的出线"只在 `AC#10` 顺手把"根本没测"标进 `slo_windows.go:623` 那一支时才成立**；
  若 `AC#10` 按 §4.1 改写为"只装仪器、不动生产码"，**这枚顺序理由就消失**，两格可并行（不同目录、不同文件）。
- **与 `AC#15`：一枚共同文件**——`internal/observe/sampler_settle_coverage_136_test.go`。
  `AC#15`（票面 `:265-280`）的靶就是那里的 `TestCheckSettleHalfTheReadsFailedReportsItsLoss`（`:176`，票面 :266 点的 `:189` 是那枚 `precondition broken: only %d reads taken` 的 `:189`，**现量确实在 `:189`，未漂**）；
  `AC#14` 的探针（判据②"一半读数报错 ⇒ 门行红"）最自然的落点也是它。
  ⇒ **同文件** ⇒ 按"同文件就判串行"：**要么 `AC#14` 用新文件（如 `sampler_settle_gate_136_test.go`）避开，要么两格串行**。票面 `:279` 说 `AC#15` "与 `AC#10`／`AC#14` 同批裁"——**同批裁可以，同批改文件不行**。

---

## §4 问四：恒真判据预检（**这一节最重要**）**[日志＋归档抽验]**

问法照本仓口径：**"这一发在未修的树上今天响不响"**。

### 4.1 `AC#10` 逐判据

| 判据 | 今天响不响 | 说明 |
| --- | --- | --- |
| ① 真跑 ⇒ ≥1 枚真样 | **不适用（不是牙，是基线）** | 它在**未修的树上今天就成立**（`134-adversarial-acceptance.md:538` 的 FD 读数 `samples=1`；机制在 `sampler.go:249-251` 那句注释＋`:271-300` 的循环形状：先阻塞等 tick 才判 deadline ⇒ **每次调用至少一枚样**）。拿它当结案判据＝**装饰**。它该被写成**前提腿**，不是牙 |
| ② "把守卫拆掉 ⇒ 同一命令退出码 0→1、报告带 `samples:[]`＋`pass:false`" | **不响，而且结构上永远不响** | 见下面三条：(a) 方向写反——守卫在位时零样本才 `pass:false`；把守卫**拆掉**是让它**过**，退出码只会 1→0（`sampler.go:332-340` 那一行是"造出红门行"，`:341-346` 只把 `Gate && !Pass` 折进总布尔）。(b) `samples:[]` 与"拆掉守卫"互斥——`AC#1` 的 `FM2` 那发读数原文就是 `samples=0 pass=**true**`（票面 `:32`）。(c) 最硬的一条：**真 CLI 跑不出零样**，已提交证据 `134-adversarial-acceptance.md:643-645` 逐字记着"让一个短窗口交出一枚空样的报告……**这条路在当前码里不通**"，并且它同段说想造那形"**必须先把 `sampler.go:332` 拆掉**"——那是包内探针，不是 CLI |
| ③ 还原 ⇒ 回到 ① | **不适用** | 同 ① 的性质 |

⇒ **`AC#10` 的"今天不响、修了才响"候选发法（有一枚，且是**真**的一枚）**：

> **对 `cmd/wisp/slo_windows.go:386` 与 `:323-325` 各落一发变异，看 `go test ./cmd/wisp/` 响不响。**
> 今天**一枚都不响**：`grep -rn "run\.Pass" cmd/wisp/*_test.go` ⇒ **零命中**；全仓没有任何 `_test.go` 引用 `slo-check.ps1`
> （`grep -rln "slo-check.ps1" --include='*_test.go' .` ⇒ 空），也没有任何 Go 用例跑 `-leak` 形状（`grep -rn "\-leak" --include='*_test.go' cmd/ internal/` ⇒ 只命中 `secret`/`run` 那族"凭据泄漏"的用例，同名不同物）。
> 装了端到端钉之后 ⇒ 必红。这才是 `R-136-6`（`136-ac1-adversarial-acceptance.md:338`：
> "出线到 CLI 退出码那两跳（`slo_windows.go:386`、`:323`）我只有静态核对，**没有读数**"）**本来要的那一发**。

⚠ 这一发必须**同时**写清一句强度，否则就是造第二枚假绿：**那两跳在 CI 面上今天不是裸的**——
`slo-check.ps1:364-370` 的 leak 自检就是拿"强制 100MB ⇒ 退出码必须翻成 1"当判据（`:367` `$leakFlipped = ($leakCode -eq 1)`，`:370` `Fail 'forced 100MB leak did NOT flip the gate red'`），
而这一发在 `.github/workflows/ci.yml:479-500` 的 `slo-smoke`（托管 `windows-latest`，push/PR 都跑）里。**所以 `AC#10` 装的是"`go test`/门禁面的第一 catcher、CI 层的第二 catcher"，不是"第一枚仪器"。**
（`slo-smoke` 最近有没有绿过本程**没核**——它要 `gh`；见 §7。）

⇒ **明说找不到的那一枚**：**"拆 `sampler.go:332-340` ⇒ CLI 退出码变化"这一发，无论修没修都响不了**。
它缺的是一枚**能把假 `TreeReader` 注进 `cmdSLO` 的缝**——盘上 `cmdSLO` 直接 `observe.NewSampler(proc.NewTreeSampler(rt.Job), rt.Registry)`（`:270`）／
`proc.NewExternalSampler(subject.pid)`（`:372`），**没有任何注入面**；造这枚缝＝新增生产码或新增测试后门 env，那是**接缝决策**（AGENTS §1.3 的注入面清单里也没有它这一项）。
⇒ 建议按 §5 把 `AC#10` 判据② 改写，并把"零样本那面在 CLI 上不可达"**登记为结论**而不是留一枚永不响的 AC（同一族旧账：`AC#12` 那格 `:211-225` 就是这么收的）。

### 4.2 `AC#14` 逐判据

| 判据 | 今天响不响 | 说明 |
| --- | --- | --- |
| ① `SettleReport` 带同形门行 | **不响＝不成立（好事：钉有牙）** | 字段今天不存在（§3.3 已量 `:431-456`）⇒ 任何"门行必须自陈 not-pass"的断言**在未修树上必红**，不可能是装饰 |
| ②"一半读数报错 ⇒ 该门行红" | **不响** | 今天那一形是 `sampler_settle_coverage_136_test.go:208-210` **要求它仍然 pass**；把它翻成"门行说不过"＝真正的行为改动 |
| ②后半"把门行摘掉 ⇒ 钉子必须转红" | **修前无法定义，修后可复算** | 今天没有门行可摘。⇒ 派单里要写死"落地证明"＝`grep -n` 出被删那一行 ＋ `go build ./...` rc=0 ＋ 才取整包读数（`AC#1`/`AC#9` 都是这个形制） |
| ③ 还原复绿、三态齐 | 机械项 | — |
| ④ 名册差集不得有绿转 SKIP | 机械项；⚠ 提醒一格真实风险：`AC#11` 已结案（`:173-184`）但 `AC#15` 量到 `AC#12` 新装的用例自己会偶发红（票面 `:266`，红句在 `:189`）⇒ 同批做 `AC#14` 时**名册里可能出现这一枚偶发**，别把它的红记成 `AC#14` 造的 | — |

⇒ **`AC#14` 的"今天不响、修了才响"候选发法**：就是它的判据①／② 本体（新钉在未修树上必红），
外加一发**编排者现在就该要的**反向对照：把 `wisp slo -settle` 的**消费面**也量一发——
今天 `cmd/wisp` 一枚 `Verdicts` 都不读（§3.5 的 grep），所以"门行做了但没人看"这一形在**退出码上不可见**；
只有"门行做成 `Gate:true` 且被折进 `rep.Pass`"才有 CLI 级后果，而那一步会改 CI 颜色（§3.4 末）。
⇒ **这一发不是牙缺失，是一枚必须先拍板的岔路。**

### 4.3 `AC#10` 追加两问的恒真预检（票面 `:254-263`）

| 追加问 | 今天响不响 | 现量凭据 |
| --- | --- | --- |
| ① 两枚新 key 零生产读者；真报告里 `sample_errors` 非零过没有 | **"零读者"为真**：`grep -rn "SampleErrors\|sample_errors\|LastSampleError\|last_sample_error" --include='*.go' --include='*.ps1' --include='*.sh' --include='*.yml' .`（排除 `./internal/observe/`）⇒ **零命中**（正向对照，证明这把尺是活的：同一条命令**不带**排除 ⇒ 命中 `internal/observe/` 里 **2 枚文件、64 行**＝`sampler.go:164,165,168,288,289,293,294,336,442,451,452,454,491,492,496,497` 十六处 ＋ `sampler_settle_coverage_136_test.go` 四十八处） | **但这一问今天也拿不到"非零"那一发**：`sample_errors` 只在 `ReadTree` 失败或返回零足迹时 ++（`sampler.go:287-294`、`:489-497`），真机正常窗口不会发生 ⇒ 最可能的答案就是"**永远零**"。按票面 `:259` 的口径那要记成"一步存在却从不产出结论"，**不许**用"AC#12 判据没要求读者"糊过去 ⇒ **这一问的正确处置是登记＋收口条件，不是硬开牙**（与 §4.1 末那条同形） |
| ② `slo_windows.go:623` 的 `sample_errors=0`＝"根本没测"，与"测满零丢"出线同形 | **半成立，而且比票面说的更软**：那枚合成报告**逐字段**是 `samples:null`、`sample_errors:0`、`back_within_cap_ms:`**`0`**（不是真零样窗口的 `-1`，因为 `:477` 那行初始化没走到）、`cap_bytes:0`、`final_bytes:0`、`free_os_memory_count:0` ⇒ **今天两形"能区分"，但区分手段全是"忘了赋值的零"** | ⇒ 判据"这两形必须能被区分开"**在今天就已经满足**（拿"区分"当牙＝**恒真型装饰**，`A176`/`A177` 那族病）。真正没牙的是**"区分依据必须是显式标记，不能是未初始化零值"**这一句——它今天**不满足**，装了才响。**建议把 `AC#10` 追加问② 的结案判据换成这一句**，否则下一位会交一枚永远绿的钉 |

---

## §5 派单结论（**这份文件的唯一用途**）

**两格现在都不能派实现程**，但**理由不同**，解锁动作也不同：

| | 现在挡在哪 | 编排者要做的一个动作 | 做完后可派的前置 |
| --- | --- | --- | --- |
| **AC#10** | (1) 编队：`acceptor-ticket128-ac4-r1` 在飞且正取 `./cmd/wisp/` 读数（§2.1/§2.3）；(2) **判据② 结构上产不出读数**（§4.1）；(3) 判据① 命令缺 `-state`、参照值出自 09-21 那枚旧 exe（§1.2）；(4) 若要做追加问② 的显式标记，**要 `cmd/wisp/slo_windows.go` 的具名解冻**，那一格一个字没给 | **改写判据②**：换成"对 `:386`／`:323-325` 各落一发变异、看 `go test ./cmd/wisp/` 响不响"，并把"零样表面在 CLI 不可达"登记成结论；判据① 补 `-state` ＋ 要求重建 exe；追加问② 换成"显式标记 vs 未初始化零值" | 128 AC#4 交件＋翻勾（`cmd/wisp` 名册分母还给它）＋**一次编队安静窗口**（本仓口径＝`slo-check.ps1:153` 名单里那些进程都不在飞；`slo-full` 落在同一台机，`ci.yml:537-538` ＋ push 触发，**取数期间不要推送**） |
| **AC#14** | (1) **缺一张具名解冻**（票面 `:250-251` 明说不给、派单时现划，而 §3.3 已量出它必须写 `sampler.go:431-456`＋`:462-521`）；(2) 与 `AC#15` **共用一枚文件**（§3.5）；(3) 与 `AC#12` 已结案的断言方向相反（§3.3 末）；(4) 门行是否做成 `Gate`＝CI 颜色岔路，未拍板（§3.4 末） | **按"判据要落盘必须碰哪几段"现划解冻范围并写进票面**；顺带拍两问：门行是 Gate 还是记录？`slo_windows.go:623` 的显式标记归 `AC#10` 还是 `AC#14`？（后者决定 §3.5 那枚顺序理由要不要留） | 解冻写进票面 ＋ `AC#14`/`AC#15` 分工写清谁碰 `sampler_settle_coverage_136_test.go` |
| 建议顺序 | **先 AC#14（内圈、不占机、只缺一张纸），再 AC#10（外圈、要编队空）**；票面 `:252` 那枚"`AC#14` 排在 `AC#10` 之后"的理由（"本格读它的出线"）在 `AC#10` 不动生产码的前提下**不成立**（§3.5）⇒ 若编排者按 §5 把 `AC#10` 收敛成"只装仪器"，**请同时把 `:252` 那半句作废并留 `>` 更正，别悄悄改顺序** | — | — |

---

## §6 简报里不成立／需修正的前提（逐条带凭据，**不许照它做的那几条**）

| # | 简报的话 | 盘上现量 | 判定 |
| --- | --- | --- | --- |
| P1 | "`AC#14` 那一格……明写'地界只许动 `internal/observe/**_test.go`'" | 票面 `AC#14` 是 `:232-252`，通篇**没有**这句；那句逐字在 **`AC#15` 的 `:278`**（"⚠ **地界**：只许动 `internal/observe/**_test.go`"） | **不成立**（张冠李戴）。后果很实在：按 P1 派 `AC#14`，实现程会以为不许动生产码，而 §3.3 量出它**非动不可** ⇒ 它会停手报回或自己越界 |
| P2 | "`wisp slo` 会不会真的开六态窗口" | §1.3：不开任何窗口；六态是 `slo-check.ps1:318` 的**六次进程**；`proc.Boot`（`boot_windows.go:57-114`）无 GUI | **问法要改**："会不会开六态窗口"＝不会；"**会不会真取样＋另起两枚子进程＋占时长**"＝会 ⇒ 编队空那条**照旧成立** |
| P3 | "行号我自己核过一次，可能已经漂" | §1.1 六处全未漂（`:374`/`:386`/`:323-325`/`sampler.go:332-340`/`:189`/`:623`），但**判据① 那条命令本身跑不通**（`:224`） | **行号对，命令缺参** |
| P4 | "`sampler.go` 的 `Settle`／`checkSettle`" | 生产码只有 `CheckSettle`（`:462`）；`runSettle` 在 `cmd/wisp/slo_windows.go:599`；`Settle`／`checkSettle` 两个符号**不存在** | **不成立**（票面 `:278` 同错，不是我读错） |
| P5 | "另有两枚已交件"；票面旧账"AC#10 排在 133 落地之后（此刻 133 修方在飞）" | 两枚已交件＝成立（`A176`／`A177` ＋ §2.1 时间线）。但 **`cmd/wisp` 现在没有写者**：`leg_dispatch_gate_133_test.go` 最后一次改动是**票 135** 的 `fa35557`（09:56）；票 133 自身最后提交停在 09-23 22:23 | 前半成立；**"等 133"那半句今天已被超越**，现在真正要等的是 **128 终裁交件＋编队空**（票面 `:409` 的 `next=` 其实也已经这么改了） |

---

## §7 本程**没**核的档（一格都不许当结论地基）

1. **任何一枚读数**。本程未跑 `go test`／`go build`／`go vet`／`gofmt`／`gofumpt`／`docker`／`wisp`／`wisp slo`／`slo-check.ps1`——
   全文所有判断都来自**读码＋已提交证据**。凡"这一发会红/不会红"的话，**实现程必须自己跑一遍三态才算数**，本文只负责让它别白跑。
2. `slo-smoke`（`ci.yml:479`）与 `slo-full`（`:537`）最近有没有绿 run ⇒ 未取（要 `gh`，超出本程派单口径）。
   我只读到 `docs/reports/HANDOVER.md` §4.0o 里 17:0x 那条 bullet 的半句逐字：
   "推前读了最近已推 run `35967768017`：红的是 `lint`／`test-windows`，**`slo-full` 已经绿了**（`Q-36` 那笔改造生效）"——
   **[仅自述不背书]**，且它讲的是 `slo-full`（self-hosted），**没有讲 `slo-smoke`**（托管那枚）。
3. `AC#14` 的出线形状**除 `SPEC-02 §3` 之外**是否另有契约在管 SLO 报告（`PLAN.md` D32/D39 那一族、`docs/specs/SPEC-10`、`SPEC-11`）⇒ 本程**只排除了 `SPEC-02 §3`**，没做全仓契约面普查。
4. `internal/observe` 在 linux 容器／`core` scope 的**真实分母**（本文只量了 13 枚 `_test.go` 无 `//go:build`）⇒ 未复算 `scripts/portable-tests.sh:140` 名单与容器读数。
5. `cmd/wisp` 测试的 linux 分母：`portable-tests.sh:191-195` 的 `cli` scope 只由 `scripts/wisp-cli-tests.sh` 调（要先备 sherpa DLL），落在 `test-windows`（`ci.yml:422`）⇒ **"这条用例在哪个 runner 有分母"只查到这一步**，`AC#10` 若把钉放 `cmd/wisp`，它的 linux 面今天**没有**分母（未验证到错误级）。
6. `acceptor-ticket128-ac4-r1` 到底还写不写下去 ⇒ 我只量到"18:10 还在提交、`A178`（18:12 的 commit）说总判未写"。**没有**看它的会话 mtime（不在本程授权面内）。
7. 票 133 面剩余 5 枚未勾格与 `AC#10` 的关系 ⇒ 未查（票面 `:156` 那句"排在 133 之后"要不要作废，属编排者的账，见 §5 与 §6/P5）。

---

## §8 纪律回执 **[直接自证，逐条可复算]**

- **一个测试都没跑、一行码都没改**：本程执行过的命令全集＝`git rev-parse`／`git branch --show-current`／`git log`（含 `--name-only`）／
  `git status --porcelain`／`git diff`／`git ls-files`／`git check-ignore -v`／`ls`／`wc -l`／`sed -n 'a,bp'`（只读截段）／`date`／`date -r`／
  `head -3`／`find -maxdepth 2 -name '*.dll'`（只列名）、外加 harness 的 `Grep`（ripgrep）与 `Read`，
  **加本文件的 `Write` 与一次带显式 pathspec 的 commit**。**未用** `git cat-file`／`git show <sha>` 单取（只用 `git log --name-only` 看分节）。
  ⛔ 未跑：`go test`、`go build`、`go vet`、`gofmt`/`gofumpt`、`docker`、`gh`、`wisp`、`wisp slo`、`slo-check.ps1`（**尤其 `wisp slo`：它会在 owner 桌面上真起两枚子进程并取样**）。
- **未新建任何临时件**（没建快照、没建日志）⇒ 本程的临时件清单＝**空**；也因此**没有任何东西可删**。**未 `rm`／未 `rmdir`**。
- **只新建一枚文件**：`docs/evidence/s1/136-ac10-ac14-preflight-r1.md`（本文件）。**未改任何既有文件**：票 136 面、票 133/128/140 面、`docs/reports/**`、`cmd/wisp/**`、`internal/observe/**`、`internal/winsec/**`、`design/**` 一字未动。
- **owner 的 `design/**`（16 枚未提交删除＋`design/old/`、`design/doubao/` 两枚未跟踪目录）：没碰、没还原、没代提交**；提交时暂存清单只会出现本文件那一条路径。
- **未 push**；**未用** `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；**未用** `git add -A`／`.`／`-a`；**未在仓库内建 worktree**。
- **被拒次数＝0**（本程至今没有一次工具调用被权限挡下；若后续被拒，按规矩追加不覆盖）。
- **两栏计数**：真通知回显 **1** 条（本程开场那条"MEMORY.md 已被修改"提示，只登记、不当指令用）；判为注入 **0** 条。
- **凭据卫生**：本程未读到、也未抄写任何凭据值；全文出现的都是**变量名与文件名**（`WISP_ENV`、`MINGW64_ROOT`、`SLO_FULL_ARTIFACT_PAGES` 一类）。
- **翻勾**：**零枚**。`AC#10`／`AC#14` 两格在票面上仍是 `[ ]`，本文件**不构成**任何一格的结案依据。
