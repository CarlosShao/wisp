# 374-a1 普查：`Confirm` 的「有腿」vs「问得出来」（只读腿）

- 腿：`374-a1`（只读）。基线 HEAD `48d193298`（2026-10-24 12:0x +08）。本件**只读**：未落任何 `.go`、未改票面/台账、未 push。
- 探针全部在仓外：`D:\tmp\wisp\20261024\374a1\`（`-overlay` 注入 `cmd/wisp` 的**测试**文件，仓内零字节）。读数原文件：`readA/readB/readDEC/readEF/readG2/readH.txt`。
- ★**尺的射程（落笔前五问的自答）**：①对象＝`perm.New(perm.Options{Confirm:…})` 的**装配值**与它落到 `approval.Gate` 的那张卡；②形状＝`grep 'perm.New('` 的**代码调用**（⛔ 不算字符串 `t.Fatalf("perm.New: %v")`，该处 1 枚）；③我数的是代码，不是注释；④总数核对：`.go` 里 `perm.New(` 代码调用 **4 枚**（生产 2＋测试 2）——第 5 枚是**字符串**；⑤超时数由 `Gate.Queue().Timeout()` 现读（`config.confirm_timeout_sec`，默认 300、带 [31,3600] ）。

## ① 四条装配：`Confirm` 交给 `perm.New` 的逐条清单

生产只有 **2 个 `perm.New(` 调用点**，但被 **4 个装配根**各自喂出不同的 `Confirm` 值（这就是"四条"）：

| # | 装配根（谁建 runtime/chain） | `perm.New` 行 | 那一枚 `Confirm` 是什么 |
|---|---|---|---|
| A | `wisp run` 控制台：`main.go:159` → `run.go:249` | `run.go:672`（`Confirm:` 在 `:674`） | `run.go:668-671`：`s.modeConfirm` 为 nil ⇒ `rt.confirmModeSwitch`（`run.go:893-895`）⇒ `confirmModeSwitchVia(rt.gate,…)`（`:905`）⇒ `gate.PendingApproval`。**永远非 nil**。 |
| B | 常驻任务的 runtime：`resident_task_source_windows.go:271-284`（注入 `ra.gate/ra.ui/ra.cards`，`reply` 留 nil） | 同上 `run.go:672` | 同一个方法值，但 `rt.gate` ＝ **悬浮球的 gate**（`:275`）⇒ 卡落在球上。**非 nil**。 |
| C | `wisp panel-inbound`：`main.go:120` → `panel_inbound.go:236` | `panel_inbound.go:272`（`Confirm:` 在 `:277`） | `storeConfirmFromPanel(nil)` ⇒ **nil**（`panel_inbound.go:346-349`：nil 进 nil 出）。 |
| D | 常驻面板档位链：`resident_windows.go:224` → `panel_resident_windows.go:233` → 同一 `panel_inbound.go:272` | 同上 | `panelModeConfirm(ra.confirmModeSwitch)` ⇒ 非 nil 的转发闭包（`run.go:938-945`）。**非 nil**。 |

测试里另有 2 个 `perm.New(` 点（`internal/panel/composer_gate_agreement_367_test.go:72`、`cmd/wisp/panel_mode_confirm_114_windows_test.go:202`，均 `Confirm` 缺省＝nil）；`cmd/wisp/run_mode101_test.go:153` 走 `runSpec.modeConfirm` 注入。**未发现第五条生产装配。**

## ② 逐条判「问不出的形状」

| # | 非 nil？ | 问不出的形状 | 判据 |
|---|---|---|---|
| A | 是 | ★**超时成拒**（有控制台、无人答）；★**「永远没被问过」**（生产无调用者）；ctx 先结束＝**当场拒** | 卡由 `consoleApprovalUI.Prompt` 打出（`run.go:1533-1578`，返回 nil≠允许），只有 `AnswerAllow` 才放过（`gate.go:545-555`）；否则 `queue.expire` ⇒ `AnswerTimeout` ⇒ 拒（`queue.go:533-548`）。READ-A：`elapsed=31.001s`、`result=refused-confirm secs=31`、`after=ask_every_step`。READ-E：ctx 已死 ⇒ **2ms** 拒（`任务上下文已结束…按拒绝处理`）。★**`rt.modeWrites` 在生产无读者**（`run.go:333-345` 自述；`grep modeWrites` 生产仅构造处）⇒ A 的这条腿在 `wisp run` 里**没有任何调用者**。 |
| B | 是 | ★**当场拒**（球窗不在）／★**超时成拒**（球在、没人点） | `ballCardUI.Prompt` 无球 ⇒ `errResidentNoBall`，`PendingApproval` 立刻 `q.abandon` ⇒ 拒（`resident_approval_windows.go:890-896`＋`gate.go:536-540`）。READ-D：**1ms** 拒（`审批界面不可达，已 fail-closed 拒绝执行`）。 |
| C | **否** | ★**当场拒**（无卡可弹、写侧直接拒）＋★读侧降级 | `Set` 的目标门 `RequiresCard(to) && s.confirm==nil` ⇒ `ResultRefused`（`store.go:268-274`）；面板先拒成 `ErrNoL2Confirm`（`composer_handlers.go:153-158`）。READ-C1：`err=…确认腿未接入…`、1 行 `MODE-REFUSED`；★**未到 `Set`**（`composer_handlers.go:153-158` 在 `Set` 之前返回）⇒ 连卡都不会入队。 |
| D | 是 | 同 B（同一个 `confirmModeSwitchVia`＋同一个 `ballCardUI`） | READ-G2：非 nil 腿＋无球面 ⇒ 卡已入队（`approval: queued … tool=permission.mode`）后 **1ms** 拒；`effective` 不变。 |

★全表**没有**「返回允许」形，也**没有**「永不返回」形：`PendingApproval` 的五个出口（admit 失败／queue 失败／nonce 失败／Prompt 失败／answer·deadline·ctx）只有 `AnswerAllow` 返回 nil（`gate.go:502-583`），而 `AnswerAllow` 必须花掉 `grantNonce` 签发的单次令牌（`queue.go:359-378,394+`），令牌只由显示过那张卡的 UI 收下。**本票最坏的两形在"会被调用的腿上"都不存在。**

## ③ `wisp run` 那条 seam（`run.go:668-671`）

- ★**能问到人的条件（三条同时）**：①stdin 真是本进程的控制台输入缓冲（`interactiveStdin()`，`approval_reply_stdin_windows.go:41-52`）⇒ `attachReplyListener`（`run.go:858`）；②卡上打印了编号（`run.go:1569-1571`）；③有人按语法输入 `yes <编号>`（`approval_reply.go:461-463`）。缺①时 `cmdRun` **响亮说明**"L2 卡会等到超时后按拒绝处理"（`main.go:154-158`）。
- ★**问不到的条件**：管道/重定向/Explorer 无控制台/CI（同上）；有控制台但无人答；ctx 先结束；**以及最重要的：生产里没人调 `Set`**（②A 表末行）。
- ★**控制台跑（无 WebView2 面板）的读数**：`queueTimeout=31s`（把 `confirm_timeout_sec` 设到带底 31）⇒ **`elapsed=31.001s`、超时成拒、`result=refused-confirm`、档位不动**；默认配置的读数是 **`confirm_timeout_sec=300` ⇒ `Queue().Timeout()=5m0s`**（READ-B）。★**拒，⛔ 不是允许**——本票最坏的那一格在这一条 seam 上未被触发。
- ⚠**同位发现（卡片文案）**：L2 卡的 `窗口` 打的是 `Prompt.Window`，而 `PendingApproval` 给 L2 传的 `win=0`（`gate.go:534` 与 `:597,609`），故控制台卡印 **`窗口 0.0s`**（READ-A 逐字），真实 deadline 在 `Prompt.Deadline` 里、控制台面**没读**（`run.go:1563`）。真正倒计时只剩"超时前 30s"那句 warning。

## ④ 有没有「非 nil 但永不拒」的腿？——没有；以及我怎么排除的

**判据（四步，逐条可重跑）**：
1. 生产 `Confirm` 的产生点只有两个方法值：`run.go:893-895`（A/B）与 `resident_approval_windows.go:111-113`（D），**两者都调同一个** `confirmModeSwitchVia`（`run.go:905`）⇒ 「腿的行为」只有一个定义，不可能分叉出第二种。
2. 该函数的 nil 出口唯一：`ans == tools.AnswerAllow`（`run.go:920-926`）。
3. `AnswerAllow` 的唯一入口要**花令牌**：`PendingApproval → grantNonce（gate.go:512）→ promptFor(…, nonce) → ui.Prompt`，`allow/allowScoped` 必 `spend`；`PanelAPI` 根本没有 `Allow`（`gate.go:626` 注释＋`:630`）。⇒ 无"人"就无令牌，无令牌就无 allow。
4. 两个转发闭包都不吞：`storeConfirmFromPanel`、`panelModeConfirm` 都是 **nil 进 nil 出**，非 nil 时原样转发（`panel_inbound.go:346-353`、`run.go:938-945`）；生产 `Confirm:` 赋值共 5 处＝2 处 `perm.Options`（`panel_inbound.go:277`、`run.go:674`）＋2 处面板 handler（`panel_inbound.go:288`、`run.go:699`）＋1 处结构体拷贝（`store.go:150`），**没有一处是字面量闭包**。
- ★**但本票真正的坏格不是这个形状**：是「**非 nil 但问不出**」本身。因为读侧只问 `confirm == nil`（`store.go:201`），非 nil 就**不降级**。READ-G0 实测：面板链（会出页面的那条）腿非 nil、球窗**不在**、`config.toml` 写 `auto_approve` ⇒ `effective=auto_approve`、`Snapshot().Mode=auto_approve`、`MODE-DOWNGRADE=0 行`。READ-F 实测 `wisp run` 无答复源＋配置 `auto_approve` ⇒ 一次 `fs.write` **0 张卡、写入成功**（`writeCards=0 writeExecuted=true`）。⇒ ★**"静默通过"在"非 nil 但问不出"这一格上已经成立**，只是它靠的是**配置值**进场，不是靠一次成功的切档。

## ⑤ `Snapshot().Mode` 与 `run.go:724` 报的是有效档还是配置档

- 两者都是**有效档**：`Snapshot()` 直接返回 `s.PermissionMode()`（`store.go:326-332`）；泵的读数是 `Mode: rt.modes.PermissionMode`（`run.go:724`）；boot 的 `MODE-READ … mode=` 也是它（`run.go:709-710`）。页面上是 `ModeView.Current = mode.String()`（`composer.go:156-179,288-305`，泵源 `pump.go:141-142`）。⇒ ★**页面显示的就是"真在跑"的档**（372-v1 那句成立，我复核为真）。
- ★**"用户以为拨成功了、其实被降级"看得见吗＝半可见**：看得见**结果**（页面 `current` 是降级后的档），看不见**原因与事实**——`ModeView` 只有三格（`composer.go:156-179`：`current`/`names`/`l2ConfirmNames`），**没有降级位**，唯一记号是构造时那一行 `MODE-DOWNGRADE`（日志，不在页面）；而请求本身的回执是 **`nil`（applied）**（READ-C2/READ-H1）⇒ 只读回执的人会以为成功。★另有一处**反向**读数：`rt.permissionMode()`（`approval_always.go:151-156`）读的是 `mgr.Config().PermissionMode()`＝**配置档**，被 `:122` 用作 `Decision.Mode` 写进卡/审计 ⇒ 那是全仓**唯一**把配置档当"在跑的档"报出去的地方（纯展示/记账，不参与执行）。

## ⑥ `MODE-DOWNGRADE` 的小口子（只在构造时一行）

- 唯一落行点：`store.go:167-175`（`New` 内），`grep '"MODE-DOWNGRADE"'` 生产**仅此 1 处**（其余 3 处是测试断言）。
- ★**进程活着把盘改宽，会不会当场落降级行？——不会，0 行**。READ-H：构造时 `configured=ask_every_step` ⇒ `downgradeLines=0`；随后把 `ask_high_risk` 落盘（`err=nil`）⇒ `effective=ask_every_step`、`Snapshot=ask_every_step`、`configuredOnDisk=ask_high_risk`、**`downgradeLines=0`**、`switchLines=1`（`perm: MODE-SWITCH from=ask_every_step to=ask_high_risk … result=applied`）。⇒ 降级**开始生效**的那一刻，日志里**一条降级行都没有**。
- ★**当时的证据只有 `MODE-SWITCH` 一条，而且它说的是反话（applied）**⇒ **判：除它之外没有任何"当时"的证据**。旁证都不"当时"：`config: D36-SECTION … effect=…` 只在 `wisp run` 的 reload 路径（`config_reload.go:182-245`；且 `[risk]` 是 locked 段 `tiers.go:41`，手改放宽要过 `confirmLockedLoosening` 的 L2 卡，nil 钩子＝拒）；`wisp panel-inbound` **根本没有 reload tick**（`config_reload.go:101-103` 具名）；`MODE-READ` 只在 boot 一次。真正那一行要等**下一次启动**、且文件仍更宽时由 `New` 补写。
- 边界：READ-C3 实测，legless 进程硬跑 `CheckAndReload()` ⇒ `locked=1`、放宽**不生效**（`configuredNow` 不动）⇒ 该口子当前只能由 `Set` 到**中间档**（M4 不弹卡）打开，⛔ 不能由手改直接打开。

## ⑦ 票 374 该写哪几条 `AC`（只建议，⛔ 未改任何票面）

1. **AC#1 可问性 ≠ 非 nil**：`perm.Options` 必须能区分"有腿"与"这张卡真能到人"；判据：腿非 nil 但无面（球窗不在/无控制台）时，boot 的 `PermissionMode()` 必须降级，且落 1 行 `MODE-DOWNGRADE`。现状反例＝READ-G0/READ-F。
2. **AC#2 每个装配点自己声明**：4 个装配根（①表 A–D）逐条在装配处写下"这张卡在什么条件下能到人"；静态尺：任何 `Confirm` 非 nil 的 `perm.New` 调用点，同函数内必须有该判定的读法。⛔ 让 `perm` 猜。
3. **AC#3 降级要"当时"落行**：读侧从"未降级"跨到"降级"的那一刻必须新增 1 行（⛔ 只靠下次启动）；判据＝READ-H 的剧本，`downgradeLines` 从 0 变 1，且行内 `result=downgraded-no-channel`。
4. **AC#4 `wisp run` 无答复源时不许静默宽档**：`reply==nil` 且配置档宽于 `ask_every_step` ⇒ 必须响亮说明"本进程按 X 档跑且无人可答"（现状只说了"L2 卡会超时"）；判据＝READ-F 的剧本，卡数必须非 0 或启动即降级。
5. **AC#5 L2 卡要印真实 deadline**：控制台卡 `窗口 0.0s`（`gate.go:534`×`run.go:1563`）必须改成真实秒数。
6. **AC#6 动作层读数只许用有效档**：`Decision.Mode` 之类记账面（`approval_always.go:151-156`）不得报配置档，或必须同时带"有效档"。

## ⑧ 门禁 rc

- `go vet -overlay <probe> ./cmd/wisp` ⇒ **rc=0**；`go test -overlay <probe> ./cmd/wisp -run 'Test374a1…'`（A/B/C/D/E/F/G/H 各跑一次，`PATH="$PWD/third_party/sherpa-onnx:$PATH"`、`GOFLAGS=`、未 unset `TMPDIR/TEMP/TMP`）⇒ **全 rc=0**。
- 探针末次全绿：A `PASS(35.45s)` · B/C/D/E/F/G/H 均 `PASS`。输出：`D:\tmp\wisp\20261024\374a1\read*.txt`。

## ⑨ 提交与 `.go` 计数

- 本件提交：`git log -1 --format=%h -- docs/evidence/s1/374-a1-census.md` 读回（⛔ amend，故不在自身正文写死哈希）；基线 `48d193298`。
- **本笔 `.go` 改动＝0 枚**（`git show --name-only --format= <本件提交>` 应只有本文件）；探针 8 个测试全在仓外、经 `-overlay` 注入。
