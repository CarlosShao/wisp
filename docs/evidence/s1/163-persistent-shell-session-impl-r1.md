# 163 实现腿裁定表 — shell.exec ＋ shell.session 落地（AC#2/AC#3/AC#4 主体）＋ AC#5 残格登记

- 时刻：2026-11-10 22:2x +08｜腿名 `163-impl-r1`（写腿：internal/tools/** ＋ internal/proc/** 零字节（proc 只被引用未被改）＋ 证据件）
- 锚：开工现 HEAD ＝ `eec9fe05422dbfb4964ddbbdcc109be322509a36`（branch `dev`）
- 隔离：共享主树直写本腿射程（internal/tools/** 为本腿独占；`cmd/wisp/run.go`、`internal/audio/**` 等 9 枚脏文件是**别的在飞腿**的，本腿零触碰，收工读数见 §5）
- 环境令照办：`export LC_ALL=C`；`unset GOTMPDIR TMPDIR TEMP TMP`；`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`；**CGO_ENABLED=0**（⛔ 默认 1；本腿全部 go 命令在该环境下跑）；gofumpt v0.12.0（`/d/work/base/gopath/bin/gofumpt.exe`）。

## 0. 开工三核（现跑）

| 核 | 读数 | 命令 |
|---|---|---|
| HEAD | `eec9fe0542…` | `git rev-parse HEAD` |
| 分支 | `dev` | `git rev-parse --abbrev-ref HEAD` |
| 本腿射程脏栏 | **0 行**（开工时 `git status --porcelain -- internal/tools/ internal/proc/` 空；全树 1575 行脏属别的在飞腿） | 同左 |
| AC#1b 三数对账 | 冻结文本 4 行在册（`PLAN.md:2574/:2575`＋`SPEC-07:74/:75`，本腿亲读）；真注册表两枚 Lookup=false（`163-ac1b-recount.md` §1，本腿 grep 零实现命中复核）；代码普查 `shell.session`=0／`shell.exec` 非测试 6 枚全非实现（本腿开工 grep 复核同数） | `grep -n 'shell\.exec\|shell\.session' docs/PLAN.md docs/specs/SPEC-07-tools-and-plugins.md`；`grep -rn 'shell\.session' --include=*.go internal/ cmd/`（开工时 0） |

## 1. 交付了什么（7 枚产码/判据件 ＋ 2 枚文书件）

| 件 | 性质 | 内容 |
|---|---|---|
| `internal/tools/shell.go` | 新产码 | 家族面：ShellDeps／SessionTree／SessionShell 两个缝、Decl×2（L2、能力 token `shell`、PathParams `cwd`、C22 预算 30m）、`BuiltinShellEntries`（**Enabled=false ⇒ 返回 nil，家族不存在**＝D14 默认禁用，FSDeps.DeleteEnabled 同形）、Facts 钩子（只喂 ShellArgv＋ShellAllowlist；ShellString 恒 false＝不实现字符串模式）、`joinSessionArgv`（会话行拼装：引号/换行/NUL/空元素 fail-closed 拒绝、8000 字符帽）、`newSessionMarker`（12 字节 crypto/rand hex，每发新铸） |
| `internal/tools/shell_exec.go` | 新产码 | `shell.exec` 本体：argv 强制、cwd 走 C26、每次调用一枚独立 kill-on-close 树（`NewTree` 缝）、Wait 收尾协程（`observe.Default.Spawn`，名 `tool-exec-shell-<task id>`）、**退出码＝Wait 状态字段**（`exitCodeOf`：ExitError→ExitCode()，其余 -1 不claim）、C22 到点⇒关树＋**不发退出码** |
| `internal/tools/shell_session.go` | 新产码 | `shell.session` 本体：任务作用域名册 `ShellSessions`（per-task 一行、`TerminateForTask`=唯一收口）、载体 `cmdSessionShell`（cmd.exe /q /d 走管道、不引 PTY、自有树）、标记协议（`<行> 2>&1` + `echo <marker> %errorlevel%`）、等待循环三源 select（`notify` wake-up／`streamOver`／call ctx）、超时⇒判失去同步＋回收＋**零退出码**、会话进程退出⇒Wait 字段取码、cwd 只在建立那发生效（活会话传 cwd＝响亮拒绝） |
| `internal/tools/bridge.go` | 小改 | `Options.TaskCleanup func(taskID string) int`（任务结束清道夫缝，nil＝行为逐字节不变）；`CloseTask` 先跑清道夫再走既有 C25 close（**同一边界、无第二族收尾动词**；收落数上审计行 `TASK-JANITOR`） |
| `internal/tools/spawn_roster_census_229_test.go` | 判据登记 | 229 普查的复审册补 2 行（本腿两枚 Spawn 站点，`tool-exec-` per_task 前缀族）——按该测试自己的补救路径，非松判据 |
| `internal/tools/shell_session_test.go` | 新判据 | 单元 7 枚：标记携码（7/0 两向）／**摘标记变异**（baitResponder：有退码形文本、无边界 ⇒ 拒判＋回收＋下一发新会话）／shell 退出走 Wait 字段／名册收口两向（自家收、外人收不动、空 id 收零）／五形响亮拒绝／cwd 只建会话时生效／join 行拼装拒绝表 |
| `internal/tools/shell_session_teeth_test.go` | 新判据 | 桥级牙口 4 枚：**AC#4**（第二发仍过门 visits==2、拒收不毁会话、万能力 GrantSource 下 L2 仍到门）／**AC#3** 正（CloseTask 收会话＋他任务收不到＋B 会话存活）／**AC#3 变异负控**（拆掉 TaskCleanup 接线 ⇒ CloseTask 后会话还在——这就是"取消后会话还在"那一发的活样本，证明上面正测有牙） |
| `internal/tools/shell_live_windows_test.go` | 新判据（windows） | 真 cmd.exe＋真 `proc.OpenJobScope` 6 枚：cwd 延续（AC#1 缺口的闭合读数）／标记携码 7／`exit 5` 走 Wait 字段／超时重置两向／shell.exec 的 Wait 码＋**a1 活样本反猜**（打印 `exit 1` 退码 0 ⇒ 判 0）／shell.exec 超时关树零码 |
| `docs/evidence/s1/163-persistent-shell-session-impl-r1.md` | 文书 | 本件 |
| `.scratch/wisp/issues/163-…-one-shot.md` | 文书 | Progress log 追加（只追加不抹） |

**proc 侧零字节**：`internal/proc/**` 一行未改——本腿只**引用**它（`proc.OpenJobScope` 经 `SessionTree` 缝，编译期断言 `var _ SessionTree = (*proc.JobScope)(nil)` 钉住方法集）。判定依据：`163-ac1b-recount.md` §3 R4/R5 指名 `proc.JobScope.StartInJob` 就是仓内正解，缝已在。

## 2. 每 AC 读数（全部现跑，命令可逐字复跑）

### AC#2 埋标记取回退出码，不许靠时间；并答"摘掉标记是否存在一发超时从此判错"

- **正向**：`go test -run '^Test163SessionMarkerCarriesExitCode$' -count=1 ./internal/tools/` ⇒ **ok**（标记边界携码 7→IsError、0→成功，标记串零泄漏进结果）。
- **活壳**：`go test -run '^Test163LiveSessionMarkerExitCode$|^Test163LiveShellExecExitCodeFromWaitField$' -count=1 ./internal/tools/` ⇒ **ok**；其中 a1 活样本（打印 `exit 1` 退码 0 的进程）判 **0**——退出码出自 Wait 状态字段，文本猜不出来（`Test163LiveShellExecExitCodeFromWaitField` 反猜支）。
- **摘标记的答案＝"存在，且那一发从此必判错"**：摘掉标记 ⇒ 流里没有边界 ⇒ 超时那一发的输出与**下一发**的输出不可区分——这就是误判本体。代码按"那世界不可达"建造：超时路径（`run` 的 `ctx.Done` 支）判**失去同步**、关 kill-on-close 树、收尸、弃掉全部已积输出、**不报退出码**，下一发开全新会话。
- **变异两向**：负向 `Test163SessionMarkerIsNotAGuess`（bait：有 `exit 1` 形文本、无标记 ⇒ 拒判、拒绝文本含「失去同步」、结果不含「退出码：0/1」、树已关、下一发换新载体 made==2）⇒ **ok**；正向同上标记支 ⇒ **ok**。
- **不许靠时间**：全家族零 `time.Now()` 差、零 sleep 轮询（等待=三源 select；预算=call ctx deadline）。仪器读数：`sh scripts/d22scan.sh` ⇒ **clean**（ban #4 之形零命中）。

### AC#3 会话寿命绑任务，取消/结束⇒会话跟着收，带变异证明

- **正**：`Test163AC3TaskEndCollectsSessionTeeth` ⇒ **ok**：`b.CloseTask("t")` 后 `Sessions.Count()` 1→0、载体 close 计数 1；他任务 `CloseTask("other")` 收不动本会话（Count 仍 1、B 的下一发仍 `退出码：0`）；两任务各收各的。
- **变异负控**：`Test163AC3UnwiredJanitorIsTheDefect` ⇒ **ok**——把 `Options.TaskCleanup` 接线拆掉，`CloseTask` 后会话**还在**（Count==1、零 close）——"任务结束（取消走同一 revoke→CloseTask 边界）后会话还在"那一发的活样本；它FAIL 的世界＝正测 PASS 的世界，牙口成立。
- **运行中取消**：会话命令在跑时任务取消 ⇒ call ctx（任务根派生）`Done` ⇒ 同一 `ctx.Done` 支回收（与超时同路径，`Test163LiveSessionTimeoutResets` ⇒ **ok**：0.7s 预算内 `ping -n 30` 被收，下一发新会话，收尾名册 Count==1）。
- **唯一收口**：`TerminateForTask` 是唯一收集器（无 by-name/by-age 扫）；清道夫挂 `CloseTask`——任务结束边界本体，无第二族收尾动词（G6 词汇零触碰：`grep -n 'OpenTask\|CloseTask' internal/tools/shell*.go` ⇒ 0）。

### AC#4 审批不许松：不许"授权一次就一直放行"

- **同会话第二发仍过 C19**：`Test163AC4SessionNeverCarriesAuthority` ⇒ **ok**：第一发 visits==1（L2 卡），第二发 visits==2（**会话存在不携带任何授权**，C19+门逐发重跑）；第三发拒收（门自己的理由穿透回结果），会话不被拒收毁掉，第四发放行仍正常——拒收是授权面的答案，不是会话的判决。
- **D45 授权不覆盖会话路线**：`Test163AC4GrantsNeverCoverTheSessionRoute` ⇒ **ok**：万能力 GrantSource（什么都答 covered）在，L2 判决仍到门（visits==1）——授权缝只服务 L1（bridge 既有分支），会话路线结构性摸不到它。
- **红线**：`internal/agent/approval/**` 零字节（读数 §5）；判级完全走既有 R6（`rules_shell.go` 一字未动，risk 侧零改——ac1b §3 R5"判级器已在册且已接线"照办）。

### AC#5 契约轴（主体做完，残格登记）

- **做了**：`docs/PLAN.md`/`docs/specs/**`/`thresholds.go`/golden/`allowlist.txt`/`tools/d22scan/**`/`frontend/**`/`design/**` **零字节**（读数 §5）；实现＋测试全在 `internal/tools/**`（允许写面）；`internal/proc/**` 零字节。
- **没做（诚实登记）**：`internal/config/unwired.go` 三行 `shell_*` missing 文案在工具落地后**开始撒谎**（"no shell.exec tool is registered"将成假话）。它既不在票面允许写面、也不在零字节名单，裁定表 §5 AC#5 行指名**待 owner 一句话**（三选一：批文案放开／批复用 `shell_enabled` 当会话总开关／维持现状）。**本腿选择了"不动它"**：〔风险如实报〕工具在 `[risk] shell_enabled=true` 时注册，而那枚 config 键今天被加载期响亮拒绝 ⇒ **组合根未接线前（run.go 归别的腿/编排者），生产路径上没人能把 shell_enabled 打开**，unwired 守卫与工具注册在接线日之前不共存，谎言不会在今日生产路径上被读到——但接线那天必须先有 owner 那句话，否则 TestEveryLockedSectionKeyIsAccountedFor 语义冲突。**残格＝AC#5 的 unwired 三行文案，等 owner，不硬做。**

## 3. 门禁读数（全部现跑）

| 门 | 读数 |
|---|---|
| gofumpt -l internal/tools/ internal/proc/ | **空**（rc 0） |
| go vet ./internal/tools ./internal/proc | **rc 0** |
| go build ./internal/tools ./internal/proc | **rc 0** |
| go test -count=1 ./internal/tools/ | **ok 32.3s**（含本腿 17 枚新测试） |
| go test -count=1 ./internal/proc/ | **ok 3.0s** |
| sh scripts/d22scan.sh | **clean**（bans #1-5 internal/=396 cmd/=66；ban #8 internal/=991 cmd/=255；frontend/=121 design/=39） |
| sh scripts/check-path-length-budget.sh | **rc 0 VERDICT GREEN**（over-budget=62 全部 in roster、roster==tree） |

## 4. AC#1 的反向读数（票面勾过的那一格，落地后复量）

- 一发会话内 `echo up`（cwd=临时目录）＋ 第二发 `cd`：**第二发打印第一发的目录**（`Test163LiveSessionCwdCarriesOver` ⇒ ok，大小写不敏感比对命中 C26 canonical 拼写）——a1 量的"两发独立 spawn 目录丢"缺口，在本腿落地后**闭合**。

## 5. 反向两读数（"禁区零字节"与"未启用即不存在"）

**① 禁区零字节（本腿提交的文件清单即凭据，`git show --stat <本提交>` 可复跑）**：
本提交 pathspec 只含 7 枚产码/判据件＋2 枚文书件（§1 表）。禁区复核（提交后现跑，输出即读数）：
`git show --stat HEAD -- frontend/ design/ docs/PLAN.md docs/specs/ internal/config/ internal/agent/ internal/risk/ internal/agent/approval/ thresholds.go tools/d22scan/ '*.golden*' '*allowlist*'` ⇒ **空**（0 枚文件）。

**② 未启用即不存在（D14 默认禁用的结构读数）**：
`go test -run '^Test163' -count=1 ./internal/tools/` 全绿，其中家族面由 `BuiltinShellEntries(ShellDeps{})` 的 nil 分支保证：未启用 ⇒ 返回 nil ⇒ 注册表里没有 shell.* ⇒ 模型侧敲门即「未知工具」。另 229 普查与 tier 一致性测试在全包绿里同跑（新增站点已按普查自己的登记路径入册，非松判据）。

## 6. 本表没测什么（逐格照实）

- `bridge.Execute` 审批链路的**端到端真 UI**（要接 approval UI）未构造——AC#4 的读数出自 GateFuncs 形测试门＋桥真链（C19→route→门），卡面渲染那半归票 21 的面。
- 组合根接线（`cmd/wisp/run.go` 注册 BuiltinShellEntries＋TaskCleanup 接线）**未做**：run.go 正被另一条在飞腿占用（工作树 9 枚脏文件），且接线会踩 AC#5 的 unwired 三行文案（等 owner）。本腿交付的是"接线即用"的家族：一条 `BuiltinShellEntries(ShellDeps{Enabled:…, Sessions: NewShellSessions(), Paths:…, NewTree: procTreeFactory 闭包})` ＋ `Options.TaskCleanup: roster.TerminateForTask` 两行即完成接线。
- `task.spawn` 的后台形与本会话形的组合（163 的 spawn 半边，a1 §11 ②点的 `run.go:355-356` 注释）未动——那是 176/164 族的既有格。
- 非 Windows 平台（`shell_live_windows_test.go` 带 windows 构建 tag）零覆盖——Wisp 是 Windows-only（D7），家族本体跨平台（fake 载体测试全绿），生产载体只在 Windows 量过。

## 7. 撤销口令

**「撤 163-impl-r1」**＝revert 本提交即可：产码 5 枚（shell.go/shell_exec.go/shell_session.go/bridge.go 小改/spawn_census 2 行登记）＋判据 3 枚＋文书 2 枚，全部随这一笔；`internal/proc/**` 与全部禁区零字节，无第二笔要追。

## 8. 自纠（提交后追加，不改上文）

- **CGO_ENABLED 自纠**：提交前的一轮测试在未显式 export CGO_ENABLED 的 shell 里跑（默认 1）；提交后复跑（显式 `export CGO_ENABLED=0`）：build rc0／vet rc0／`go test -count=1 ./internal/tools/` ok 33.0s／`./internal/proc/` ok 1.4s——**CGO=0 下全绿，结论不变**（本腿零 cgo 依赖，两环境读数一致）。读数时刻 2026-11-10 22:3x，HEAD＝提交 `75f9479fa`。
- 提交式自查：`git show --stat 75f9479fa -- <禁区 pathspec>` ⇒ 0 行（§5 读数①已在提交后现跑）；本腿两处 pathspec（add/commit）逐字同一组 10 枚路径；`git status --porcelain -- internal/tools/ internal/proc/` ⇒ 0 行。
