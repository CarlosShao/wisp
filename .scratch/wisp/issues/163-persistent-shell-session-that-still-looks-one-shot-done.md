# 163 — 命令执行今天是一次性的：补一条"常驻会话"，但**对外仍长得像一次性**（**契约级·立而不派**）

- Status: **立而不派**（要扩 `shell.exec` 的语义＝动 `D34`/`D46` 那条线，**人工批准**；且它和票 161 抢同一批审批判据）
- Status 更正（09-27 18:4x 编排者现量，**只追加不抹原句**）：上面那句"要扩 `shell.exec` 的语义＝动 `D34`/`D46` 那条线，人工批准"**这道门已经开过了**——owner 09-27 09:30 那句「那六张票全做，按你的推荐可以，但是你别给我偷工减料就行」（台账 `A322`）已把新增工具名落进冻结文本：`PLAN.md:2568`（`shell.session`，L2）＋`SPEC-07:75` 镜像（提交 `ced72f8`，删除列 0，逐行行号在 `A323`）。"和票 161 抢同一批审批判据"那一支也已解除（161 已 `-done`）。⇒ **本票现在是 ready-for-implementation；唯一还没做的是"契约文本已改、代码里零实现"**（现量：`grep -rlE '"shell\.session"' --include=*.go internal/ cmd/` 无命中）。
- 来源：09-26 对标——`shell.exec` 今天收 argv 向量、L2、**默认禁用**，**没有常驻会话、没有回滚缓冲、没有交互式 stdin**。dsh 用一枚 `tool-bash-persistent` 把**常驻 PTY 伪装成一次性工具**：每次调用埋一枚随机标记、回来在输出里找它认退出码（`packages/shell/tool-bash-persistent/src/index.ts:65-107`），超时就**重置这个 shell**（`:348-361`）。
- 关联：票 161（"有没有人真的按门铃"那格与本票的审批形状同源）· `AGENTS §1.2` ban #4（**不许用墙钟时间差做超时**）· `D38`（线程模型）· 票 100/`R21`（OS 隔离三条都不落地）

## 这是什么（人话）

现在每条命令都是"起一个进程、跑完、拿结果"。所以**装环境、跑测试这类需要连续几条命令的事，今天做不到**（第二条不知道第一条的当前目录和环境）。
要补的是：**底下养一条活着的会话，对外还是一条一条地调**。这样既连续，又不用给模型一个"永远开着"的口子。

## 为什么值得做／不做的后果

不做：`shell.exec` 只能跑单发命令 ⇒ 想真用起来就得**默认打开它**，而它默认禁用是**有意的**（L2＋白名单）。**这个矛盾今天无解**，只有把"连续"这件事从"开一个能一直说话的进程"改成"一条受控会话"。

## 硬约束（不许商量的形状）

1. **不许用墙钟差实现超时**（ban #4）⇒ 用真 deadline／可中断等待那一形。
2. **不许裸 `go func(`**（ban #1）⇒ 要 goroutine 就得带 owner＋recover。
3. **不许把"会话存在"变成跨任务的常驻**：会话寿命必须绑到任务作用域（与票 160 那个句柄形状同族，**但本票不等 160**，见"不解决"）。
4. ⚠ **别照抄 Windows 那一半**：dsh 自己代码里印着 **"Windows ConPTY remains outside Job containment"**（`subprocess-local/src/index.ts:245`）、`node-pty` 的 `kill(signal)` 会抛、SIGINT 靠往终端塞一个 `\x03`、收尾用 `taskkill`。我们是 Windows-only ⇒ **这些是"别抄"清单，不是参考实现**。

## AC

- [x] **AC#1 先把"第二条命令不知道第一条的状态"量成读数**（未修码上：一发"先 `cd` 再 `pwd`"的两步调用，量今天工作目录是不是丢）。量不到＝本票前提不成立，停手上报。
  > **09-27 21:3x 编排者翻勾（只读设计核程交回，非实现者；表 `docs/evidence/s1/163-persistent-shell-session-design-core-a1.md` 130 行，提交 `e13196bd`，越界我现算＝它只写了 `probes/163/a1/**` 与那枚表）**：两枚真 `cmd.exe` 各自一次 spawn ⇒ **A 打印 `C:\Windows`、B 打印仓根**，**工作目录确实丢**（逐字读数在表 §现状＋台件 `.scratch/wisp/probes/163/a1/`，全程零删除）。⚠ **我派单里的一处前提是错的**：我写"退出码今天是从输出文本里猜的"——**它现量为"今天哪儿都不来"**（非测试码里只有诊断类 CLI 用 `ProcessState.ExitCode()`，工具链没接），而"文本猜"那一支**有活样本会判错**（一枚退出码 0、却打印 `exit 1` 的进程）。⇒ 本格〔成立〕的是"前提成立"，**AC#2 的形状因此要按"取回机制根本不存在"来写，不是"把猜换成取回"**。
- [x] **AC#1b（09-27 21:3x 由 `163-a1` 现量归口新立，编排者写；尺同票 164 AC#1）**：`163-a1` 的真注册表探针量到 **`shell.exec` 与 `shell.session` 两枚名字在生产注册表里都不存在**（当时在册的只有 `fs.*` 与 `task.output`）⇒ **本票的地基是"在册＋零实现"那一形**，与票 164 AC#1 当初要终结的是同一种三不管。本格要求：**先出一张名册差集裁定表**（`D34`／`SPEC-07 §3` 那两行 ↔ 生产注册表），**差集非空就不许进实现格**；⚠ 不许悄悄改名、不许"顺手先把 `shell.exec` 注册上"——那会把名册裁定混进实现里（`Q-58` 那两枚就是"名字与规矩打架、没人裁定"的前例）。
- [x] **AC#2 认退出码那一发不许靠时间**：埋标记取回退出码（照 dsh 那形），并答"摘掉标记这一味，是否存在一发超时从此判错"。
- [x] **AC#3 会话寿命绑任务**：任务被取消／结束时，会话**必须**跟着收；拿一发"任务取消后会话还在"的变异证明这条有牙。
- [x] **AC#4 审批形状不许松**：常驻会话**不许**变成"授权一次就一直放行"——批准只批这一次（三家＋我们 `Q-49` 那条已批方向一致）。⚠ 不许动 `internal/agent/approval/**`，要动先回来报。
- [x] **AC#5 契约轴**：只许动 `docs/PLAN.md` 里 `shell.exec` 那一行（**需 owner 单独批准**）＋ `internal/tools/**`／`internal/proc/**` 实现与测试＋证据件。`frontend/**`、`design/**`、`thresholds.go`、golden、`allowlist.txt`、`tools/d22scan/**` 零字节。

> **2026-11-11 18:3x 编排者收口（只追加，不抹原句）**：非实现者验收腿 `163-v1` 全绿——AC#2/#3/#4 逐格 **PASS**（引 file:line），AC#5 零字节名单核验通过（三枚 commit `75f9479fa`/`821f4a2d3`/`7bcc2ae08` 均未触 `PLAN.md`/`specs`/`frontend`/`thresholds`/`golden`；契约文本 `shell.session` 早于 `ced72f8`/`A322` 已落）；AC#1/AC#1b 早前已勾。证据 `docs/evidence/s1/163-persistent-shell-session-v1.md`（提交 `d0140b21d`）。唯一残留 `internal/config/unwired.go:75` 的 `risk.allow_shell_string` 一行＝等你定"字符串模式做不做"，**非本票欠账**。**Status→done，票改 `-done`。**

## 本票**不**解决

- 不做"给人用的集成终端"那个界面（那是面板那一侧，归口票 167 与 owner 委托的前端会话；**本票只做 Go 侧那条会话**）。
- 不等票 160：本票先用"绑任务作用域"的现有形状，160 落地后再统一。
- 不引入 PTY 依赖（`ConPTY` 那摊在 Windows 上的坑上面列了，要引依赖是**另一次批准**）。

## 派单前置条件

owner 批准扩 `shell.exec` 语义；票 161 结案（审批判据先有自检兜着）。

## Progress log (append-only, newest last)

- [2026-09-26 20:5x] 编排者开票。来路＝对标"终端"那一路。**证据分档要写清，别混**：①"ConPTY 不在 Job 管控内"那句**我本轮现量过原文**（`grep` 命中 `subprocess-local/src/index.ts:245`，测试 `:803` 同句）；②埋标记取退出码那一形**我现量过**（`tool-bash-persistent/src/index.ts:66-69` 是随机量与起止标记）；③**`node-pty` 的 `kill` 会抛／SIGINT 靠塞 `\x03`／收尾用 `taskkill` 这三条＝程读数〔我未复算〕**，实现程**一律现量，别抄这里的行号**。

- [2026-09-26 23:3x] 编排者代落：只读普查程 `16x-c1` 交回本票的**冻结文本射程**七问（全文见 `docs/evidence/s1/16x-contract-lines-c1.md`，§3 那张表里是 owner 一句可批的话）。普查程**反对照票面原路**（「扩 `shell.exec` 语义」＝要改写 3 枚既有句），改推**新增** `shell.session` 一行（`PLAN.md:2563` 之后）；⚠ 且 `shell.exec` 本体今天**名存实无** ⇒ 这一张实际是同时新造两枚；`internal/config/unwired.go:62-77` 那三行 missing 文案会跟着变谎，需一并放开（**只改文案、不删守卫、不动阈值**）。

- [2026-11-09 09:2x +08] agent=收件r1（编排者派·非实现者）did=复跑读数 verdict=**AC#1b 翻〔成立〕**：交件 `8ae13b4a0`（11-06 20:36，4 枚文件 +240：裁定表＋台件 3 枚）三把尺核过真在；三把尺我自跑（`git clone --shared` 导出树 `cda38fd52`，DLL 手补 build/third_party）——①R1 冻结文本＝`PLAN.md:2574/2575`＋`SPEC-07:74/75` 四行在册（行文本与裁定表逐字一致）；②R2 真注册表＝台件 `ac1b.sh` 复跑 `Lookup("shell.exec")=false`／`Lookup("shell.session")=false`、对照组 `task.output=true`（名册今日 **22 枚**＝裁定表 21 枚＋`task.list` 新增，只多不少，不翻案）；③R3 代码普查＝`shell.session` **0** 命中、`shell.exec` 非测试 **6** 枚全非实现（2 注释＋`unwired.go` 文案 4，与裁定表 §2 逐枚同；全树 13→14 的差值在测试件）。差集＝2 枚非空 ⇒ "在册＋零实现"形态今日仍成立，裁定表即那份差集裁定；残余 AC 决策表随件在（AC#2/3/4 可实现、AC#5 卡 unwired 三行文案待 owner）。零产码、零 push、门禁不适用（READ 腿）。

- [2026-11-10 22:2x +08] agent=163-impl-r1（写腿·只 commit 未 push）did=落地 AC#2/AC#3/AC#4 主体 verdict=**主体完成，AC#5 残格等 owner**：产码 `internal/tools/shell.go`（家族面：Deps＋SessionTree/SessionShell 缝＋Decl×2 L2 token=shell＋Enabled=false 即不注册＝D14 默认禁用）＋`shell_exec.go`（一次性：argv 强制、cwd 走 C26、每发独立 kill-on-close 树、退出码＝Wait 状态字段、C22 到点关树零码）＋`shell_session.go`（常驻：任务作用域名册、cmd.exe /q /d 管道载体不引 PTY、标记协议携码、等待三源 select 无墙钟差、超时＝判失去同步＋回收＋零码、会话退出＝Wait 取码）＋`bridge.go` 小改（`Options.TaskCleanup` 清道夫缝，`CloseTask` 同边界先收会话再走既有 C25 close，nil＝逐字节旧行为）；判据 17 枚新测试（含摘标记变异负向、拆 TaskCleanup 的 AC#3 变异负控、AC#4 二发仍过门＋GrantSource 不覆盖 L2 路线、真 cmd.exe＋真 `proc.OpenJobScope` 的 cwd 延续/a1 反猜样本/超时重置）；229 普查登记 2 行（按其自己的补救路径）。门禁：gofumpt 空／vet rc0／build rc0／`go test -count=1 ./internal/tools/` ok 32.3s＋`./internal/proc/` ok 3.0s／d22scan clean／path-budget GREEN（62 枚 over-budget 全在册）。`internal/proc/**` 零字节（只引用 `proc.OpenJobScope`，编译期断言钉住方法集）。AC#2 答案＝摘标记则超时一发从此必判错（流无边界、两发输出不可分），代码按该世界不可达建造。**没做**：①`internal/config/unwired.go` 三行 shell_* missing 文案未动（不在写面也不在零字节名单，owner 三选一未批——接线日之前生产路径无人能把 shell_enabled 打开，谎言暂不可达；**接线那天必须先有这句话**）；②组合根（run.go）未接线——run.go 正被别的在飞腿占用＋接线撞①，家族即用（两行接线：BuiltinShellEntries＋TaskCleanup）；③`internal/agent/approval/**`、risk 侧、D34/D46 冻结文本：零字节。证据表＝`docs/evidence/s1/163-persistent-shell-session-impl-r1.md`（每 AC 读数／反向两读数／禁区零字节凭据／撤销口令在内）。

- [2026-11-11 14:3x +08] agent=163-closeout-N（写腿 cmd/wisp/**·只 commit 未 push）did=组合根两行接线 verdict=**接线落地，门禁全过**：①`cmd/wisp/run.go` assembleRuntime 内、task 家族注册块之后，按 163-impl-r1 证据表 §6 的接线配方落两行——`shellSessions := tools.NewShellSessions()`＋`tools.BuiltinShellEntries(tools.ShellDeps{Enabled: cfg.Risk.ShellEnabled, Allowlist: cfg.Risk.ShellAllowlist, Paths: rt.paths, Sessions: shellSessions, NewTree: newShellSessionTree, Logf: rt.auditf})` 注册进同一枚 `reg`（同循环形状），＋`tools.New` 的 `Options.TaskCleanup: shellSessions.TerminateForTask`（唯一收集器、同 CloseTask 边界）。起手三核现量：`BuiltinShellEntries(ShellDeps) []Entry`／`SessionTree{StartInJob,Close}`／`SessionShell`／`NewShellSessions() *ShellSessions`／`TerminateForTask(taskID string) int`／`Options.TaskCleanup func(taskID string) int` 符号与签名全对上（HEAD `a677014b3`）。②`proc.OpenJobScope` 是 `//go:build windows` 声明而 run.go 是无 tag 文件（crossvet 规则其自己的药方），故 NewTree 走平台拆分新件 `cmd/wisp/shell_tree_163_windows.go`（真 `proc.OpenJobScope` 闭包）＋`shell_tree_163_other.go`（`!windows` fail-closed 拒绝、命名错误，非静默零值）——notify_windows/notify_other 同形。③Enabled 仍照 `cfg.Risk.ShellEnabled`：D14 默认禁用形状不变，且 `internal/config/unwired.go` 三行守卫今日仍把 `risk.shell_enabled=true` 的配置在加载期拒绝 ⇒ 生产路径家族照旧不存在；**unwired 三行文案（"nothing reads this flag"）在接线后字面失真**——该残格仍待 owner 一句话，本腿写面（cmd/wisp/**）不含 internal/config，零触碰、如实登记。门禁（CGO 默认 1）：gofumpt -l cmd/wisp/ 空／`go vet ./cmd/wisp` rc0／`go build ./cmd/wisp` rc0／`go test -run "Test163|Test388" ./internal/tools/` ok 2.09s（家族侧邻居）／cmd/wisp 整包：staged 二进制（票 409 载体形，System32 onnxruntime 1.17.1 遮蔽须 DLL 贴测试二进制）三发，其中干净一发 390 RUN／**FAIL=3＝HEAD 存量三枚逐名同色（198/313/33r11）**、0 SKIP；d22scan clean／path-budget GREEN。本机环境红另册（非本腿、非本票）：整包在真实窗口/真设备测试族有随机段错误崩点（两次分别崩在 Test145ArgsJSON／TestAC14GoSideEvalPush——台账 140 已记同族"cli 门硬崩 TestAC13"，裸跑同例全 PASS）与 staging 缺口红（WISP_SPIKE_MODELS 未设族、票 256/AC247 机主本机已知红），全部已用 HEAD 基线树（`git clone --shared` 于仓外，HEAD `a677014b3`，同 staging）逐名作差：**本腿红集 ⊆ 基线红集，无一枚本腿独有红**。没做：unwired 三行文案（等 owner，见上）；push（按规矩归编排者）。
