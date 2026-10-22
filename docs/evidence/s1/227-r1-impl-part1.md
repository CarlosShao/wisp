# 227-r1 落地腿（写代理）·票 227 AC#6 · part1/2

起手锚 `527b7377`。只做票面最后一格 AC#6：给「真轮询之下，手改仍看得见」补一枚常驻用例，
并把 `always_write_no_clobber_226_test.go` 那句已不成立的注释就地追加更正（原句一字未删）。
写面：`cmd/wisp/always_write_no_clobber_227_test.go`（新）、
`cmd/wisp/always_write_no_clobber_226_test.go`（只加注释）、本文件对（part1 读数，part2 门禁）。

## sec 1 改前读数：真轮询下那一跳今天是什么

票面两条前提现跑逐条复认（`logs/08-premises.txt`）。前提 1：`nothing polls` 仍逐字在场，且是
跨行形（单行 `grep 'nothing polls'` 为 0，`227-v1` 栽的就是那把尺）：

```
:43  "	// The edit lands after assembly, while the process is running and nothing"
:44  "	// polls: [app] theme, a key with nothing whatever to do with [fs]."
:43 ends with 'nothing'  YES      :44 carries 'polls:'  YES
```

前提 2：`:102` 方向与 D36 热档相反（HEAD 读数；本腿没动这条断言）：

```
102:	if got := rt.mgr.Config().App.Theme; got != "dark" {
104:			"(got %q, want the dark this run loaded): reading a hand edit back is ticket 223", got)
```

前提 3：`config_reload.go:123` `observe.Default.Spawn("watchdog", "config", ...)`、`:92`
`configReloadPollInterval = time.Second`、`run.go:861 rt.startConfigReload()` 三处在位，
`248095d1` 是 HEAD 祖先。故两条前提都成立，本腿无需走「前提不成立则停手上报」那支。

跑满 1 秒后旧用例会怎样（决定性读数，`logs/07-probe-226-ordering-3x.txt`）：台件是票 226 那枚
用例的同一顺序（手改、一直、受守卫写、收口读），唯一差别是受守卫写落账之后再活 `hold` 时长。
两臂同码、连打三发：

```
hold=0s   (原样)  : got "dark"  ; hand-edit to read span 84ms   ; applied_lines=0  => 绿
hold=2.5s (跨 tick): got "light" ; hand-edit to read span 2.592s ; applied_lines=1  => 会在 :102 转红
hold=0s   (原样)  : got "dark"  ; hand-edit to read span 95ms   ; applied_lines=0  => 绿
hold=2.5s (跨 tick): got "light" ; hand-edit to read span 2.567s ; applied_lines=1  => 会在 :102 转红
hold=0s   (原样)  : got "dark"  ; hand-edit to read span 99ms   ; applied_lines=0  => 绿
hold=2.5s (跨 tick): got "light" ; hand-edit to read span 2.573s ; applied_lines=1  => 会在 :102 转红
```

三发同向、`applied_lines` 逐发 0/1：旧用例今天的绿靠它自己 84-99ms 的收口窗口，一旦活过一枚
tick 就红在 `:102`。两臂跑的是同一个 1s tick，只差「写完还活多久」，故不是靠调间隔或缩短 sleep
得到的。另一发（无受守卫写、只有手改）把「谁读回来」钉到 tick 头上
（`logs/02-probe-real-tick.txt`）：`t=0.1s theme="dark" applied=false` /
`t=0.6s theme="dark" applied=false` / `t=1.1s theme="light" applied=true`。
旧用例自身颜色（`logs/01-before-226-timing.txt`）三发全 PASS（3.59s / 7.49s / 4.40s），
但这不构成反证：整发时长由 mockllm 冷启与门控决定，与「手改到收口读」那段无关（`223-v2` 早把
「整发跑不满 1 秒」这个解释推翻过），真正的量是上面那个 84-99ms 窗口。

## sec 2 新用例：逐名，以及为什么今天响

`TestAC6AHandEditIsStillSeenUnderTheRealPoll`，落在
`cmd/wisp/always_write_no_clobber_227_test.go`。它断的五件事每件都要真轮询跑起来才有：

| 臂 | 断什么 | 为什么今天只有真 tick 能满足 |
|---|---|---|
| (a) | 审计里出现 `config: HOT-RELOAD state=applied hot=[app]` | 这行只有 `config_reload.go` 的 `reportReload` 写，只由 tick 调 `reloadOnce` 触达 |
| (b) | `rt.mgr.Config().App.Theme == "light"` | D36 把 `app.theme` 列为 hot，写进内存的是 tick 的 plan/commit |
| (c) | `rt.cfg.App.Theme != "light"` | 开机快照不许动，否则 (b) 是启动值冒充的 |
| (d) | 操作员流里（从手改那一刻起的窗口内）出现 `配置热加载：` 且带 `[app]` | 232 那形：手改前印的同样字样不许顶替 |
| (e) | 文件仍带 `theme = 'light'`、仍带存下的规则、仍 `WIDEN-APPLIED` | 226 那一跳（受守卫写）在同一进程里真跑了，手改没被它回滚 |

不许裸 Manager、不许手动 `CheckAndReload`：用例体里没有 `config.NewManager`、没有任何
`CheckAndReload()` 调用、没有给 `Manager.ConfirmLocked` / `OnRestartPending` 赋值。它只做两件事：
往盘上写文件、等产品自己的审计行（`waitForAudit187` 加 `waitForWindow227`）。场景复用票 201 的
`newWidenRun` / `driveAlways`（已是「装配好的 `wisp run` 加脚本化答复流」）。没有缩短任何 sleep、
没有调大 tick：`configReloadPollInterval` 零改动；两个等待都是 `time.NewTimer` 加 20ms
`time.NewTicker` 的 40s 上界，是失败上界而非成功条件，超时即 `t.Fatalf` 并 dump 全部流。

## sec 3 它走的是真轮询路径的证据

调用链（生产码，本腿未改）：`runTextTask`、`assembleRuntime`、`run.go:861
rt.startConfigReload()`、`config_reload.go:113 observe.NewRoot("config-reload")`、`:123
observe.Default.Spawn("watchdog", "config", rt.reloadRoot, rt.configReloadTick)`、`:136
time.NewTicker(configReloadPollInterval)`、`:148 rt.reloadOnce()`、`:157 rt.mgr.CheckAndReload()`。
用例拿到同一个 `*agentRuntime`：`replyHost.run` 的 `onRuntime` 钩子
（`approval_reply_201_test.go:182-186`）在 `assembleRuntime` 返回后、任务开跑前被调。
反向判据把它钉死（sec 4）：摘掉 `:123` 的 `Spawn`、或摘掉 `:148` 的 `reloadOnce()`，(a) 都不再
出现、用例即红，说明它等的是那条真 goroutine。未突变时红句 dump 里逐字带
`[audit] config: HOT-RELOAD state=armed tick=1s ... goroutine=watchdog owner=config`，
随后 `[audit] config: HOT-RELOAD state=applied hot=[app]` 与
`[audit] config: HOT-RELOAD-READER section=app tier_row=[app.theme] ...`。

## sec 4 反向判据：三读数逐字，加 sha256

两发突变，落点不同（M1 摘「谁在跑」，M2 摘「跑的时候干什么」），都只动
`cmd/wisp/config_reload.go`、都只换一行、都当场还原；每发落刀前后各验一次命中数。
M1 落点 `:123`，原文
`rt.reloadHandle = observe.Default.Spawn("watchdog", "config", rt.reloadRoot, rt.configReloadTick)`，
改后 `rt.reloadHandle = nil // MUTATION M1 (227-r1): the production poll is not spawned`。
命中数落刀前 1、落刀后 `old: 0  new: 1`（`logs/41-mut-m1-apply.txt`），`git diff` 为 1 增 1 删、CR 0。
变异红（`logs/42-mut-m1-red.txt`）：

```
    always_write_no_clobber_227_test.go:192: the production reload tick never booked "config: HOT-RELOAD state=applied hot=[app]" within 40s: the hand edit was written while the run was alive, so either the tick is not running or it did not see the change
--- FAIL: TestAC6AHandEditIsStillSeenUnderTheRealPoll (43.56s)
FAIL	github.com/CarlosShao/wisp/cmd/wisp	43.723s
```

还原绿（`logs/43-mut-m1-restored-green.txt`）：

```
--- PASS: TestAC6AHandEditIsStillSeenUnderTheRealPoll (4.99s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	5.258s
```

M2 落点 `:148`，原文 `rt.reloadOnce()`，改后
`_ = 0 // MUTATION M2 (227-r1): the tick wakes but never polls`；命中数落刀前 1、落刀后
`old: 0  new: 1`（`logs/44-mut-m2-apply.txt`），`state=armed` 那行仍在（`:124`，故意不碰）。
变异红（`logs/45-mut-m2-red.txt`）：

```
    always_write_no_clobber_227_test.go:192: the production reload tick never booked "config: HOT-RELOAD state=applied hot=[app]" within 40s: the hand edit was written while the run was alive, so either the tick is not running or it did not see the change
--- FAIL: TestAC6AHandEditIsStillSeenUnderTheRealPoll (43.48s)
FAIL	github.com/CarlosShao/wisp/cmd/wisp	43.667s
```

M2 比 M1 更硬：`armed` 横幅照旧印、用例仍然红，故它不是被启动横幅喂饱的。原样绿
（`logs/60-selector-precommit.summary.txt`，还原后同一选择器）：

```
--- PASS: TestAC1AlwaysBranchDoesNotRevertAHandEditedKey (3.05s)
--- PASS: TestAC6AHandEditIsStillSeenUnderTheRealPoll (3.78s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	7.002s
```

还原核 sha256 三处逐字相同（`logs/87-sha-restore.txt`）：

```
config_reload.go  worktree : 72762bd5992479c7768d011b7126967556685952d777cc3fe2a5ff33896c602c
config_reload.go  restore  : 72762bd5992479c7768d011b7126967556685952d777cc3fe2a5ff33896c602c
config_reload.go  HEAD blob: 72762bd5992479c7768d011b7126967556685952d777cc3fe2a5ff33896c602c
git diff --quiet -- cmd/wisp/config_reload.go  =>  yes: no diff at all
```

读数为 raw 件，全部在 `D:\tmp\wisp\20261022\227r1\logs\`（本文件引用的 `logs/xx` 即其文件名）。
