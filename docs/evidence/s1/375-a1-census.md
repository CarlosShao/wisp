# 375-a1 普查：「能问」⛔＝「非 `nil`」——问不出的全部条件 ＋ boot 可判性（只读腿）

- 腿 `375-a1`（**只读**）。开工 HEAD `66456e479`；**全部读数归属 `8ced487bb`**（见"脏树"一条）。本件只读：未落 `.go`、未改票面/台账/`Status:`、未 push、未翻框。
- ★**脏树（开工后发现，已按纪律绕开、⭐上报）**：12:24 开工时 `cmd/wisp`／`internal/perm` **干净**；12:26 编排者提交 `8ced487bb`；此后 `cmd/wisp/run.go` 出现**别人的未提交脏**（票 373 件一 `taskLifecycleJournal`，`git diff --numstat`＝`98 4`，新块在 `noteTask` 之后、并加 `sync/atomic`）。12:32 一发 `go test` 撞上其半写态（`run.go:1153: undefined: taskLifecycleJournal`）。⇒ **不碰、不清**；把 `git show 8ced487bb:cmd/wisp/run.go` 导到 `D:\tmp\wisp\20261024\375a1\export\run.go`，`-overlay` 同盖 `run.go` ＋ 探针（仓内零 `.go` 字节）⇒ 读数归属 HEAD、免疫并发写。
- 探针：`D:\tmp\wisp\20261024\375a1\zz375a1_probe_test.go`（5 发）。原文：`head-Test375a1*.txt`（脏 `run.go` 那两发另存 `Test375a1ReadF.txt`／`Test375a1StrictModeSameWrite.txt`，只作旁证，**本件读数一律引 `head-*`**）。
- ★**尺的射程（落笔前五问）**：①对象＝**被 `perm.New` 收到的那枚 `Confirm`** 在**四个装配根**上"**此刻**答得出吗"，以及**同一发工具调用**在不同有效档下的路由；⛔ 不是"`Confirm` 是否 `nil`"（那是 372 的尺）。②形状＝`READ-x key=value` 逐字 ＋ `[audit]` 行逐字。③数的是**代码调用**：生产 `perm.New(` **2** 枚、生产 `Confirm:` 赋值 **4** 枚（`.scratch/**` 的 probe 备份⛔计）。④总数核对：生产 `Set(` 调用者 **2** 枚（皆 `internal/panel/composer_handlers.go:147,163`）⇒ 没有第五条"能改档"的路；生产 `perm.New` 只有 `store.go:140` 一个定义。⑤超时数由 `Gate.Queue().Timeout()` 现读，⛔ 不抄。

## ① 三发复跑（READ-F／G0／G2）＋ 两发我加的对照

| 读数 | 装配 | 条件 | 复跑结果（逐字） |
|---|---|---|---|
| **READ-F′** | `wisp run` 控制台（`run.go:668-674`，`confirm` 非 `nil`） | 配置 `auto_approve` ＋ 无答复源（`reply=nil`）＋ `fs.write` **新文件** | `configured=auto_approve effective=auto_approve snapshot=auto_approve replyNil=true replyListenerNil=true queueTimeout=31s l1Window=2s`；`cards=0 executed=true assessedLevel=L1 elapsed=38ms`；日志 `tools: MODE-SILENCE mode=auto_approve assessed=L1 effective=L0 tool=fs.write rules=[R1]` ⇒ ★**0 张卡、写入成功**（与原读数同形） |
| **READ-S**（我加） | 同上，**最严档** | 配置 `ask_every_step` ＋ 无答复源，同一发写 | `cards=1 executed=true elapsed=2.04s`；`approval: ANSWER-EXPIRED … decision=timeout->execute after=2s` ⇒ ★★**最严档也照写**，L1 窗口"到点即执行" |
| **READ-W**（我加） | 同上，`auto_approve` | 同档，`fs.write` **覆盖已有文件**（L2/R1+R8） | `cards=1 executed=false assessedLevel=L2 elapsed=31.024s`；`MODE-REDLINE mode=auto_approve still asks … rules=[R1 R8]`；`ANSWER-EXPIRED … decision=timeout->reject after=31s` ⇒ ★**31s 超时→拒**；卡面 `窗口 0.0s`、真 deadline 只在 `[warning] 审批将在 30 秒后自动拒绝` |
| **READ-G0′** | 常驻面板档位链（`resident_windows.go:224` → `panel_resident_windows.go:233` → `panel_inbound.go:272`） | 球窗**不在** ＋ `auto_approve` | `legNonNil=true orbPresent=false effective=auto_approve snapshot=auto_approve configured=auto_approve downgradeLines=0` ⇒ ★**0 行 `MODE-DOWNGRADE`**（与原读数同形） |
| **READ-G2′** | 同上（写侧） | 同上，页发"切 `auto_approve`" | `err=perm: 切到 auto_approve 被拒绝: L2 强确认未通过（reject）：审批界面不可达，已 fail-closed 拒绝执行 elapsed=1ms effective=ask_every_step switchLines=1` ⇒ ★**1ms 当场拒**、档位不动 |

★**这三＋两发一起说清了**：安全网确实在，但它守的是 **L2**（READ-W/G2）；**L1 的到点极性是"执行"**（READ-S），而 `auto_approve` 把 L1 连问都不问（READ-F′ 的 `MODE-SILENCE`）。

## ② 「boot 时能不能判出这条腿问不出」——逐条装配 ＋ 判语

| # | `perm.New` 被调的那一行 | 那一刻"问得出吗"这个信息**存不存在** | 判语 |
|---|---|---|---|
| **A** | `run.go:672`（`Confirm` 在 `:674`），在 `assembleRuntime(s runSpec)`（`:400`）里 | ★**存在**：`s.reply` 是**同一函数作用域的字段**，且由调用方 `main.go:153 interactiveStdin()` 在进 `runTextTask` 之前算好（`:154-158` 还为此打印了响亮句） | ★〔甲〕**可行**（读 `s.reply`）。但 boot 只能说"有没有这条流"，⛔ 说不了"有没有人在打字"——后者只有 `5m0s` 的 deadline 能回答 |
| **B** | 同上 `run.go:672`，`spec.gate=ra.gate`（`resident_task_source_windows.go:275`） | ★**存在**：B 在 `resident_windows.go:367` 才装配，而球在 `:309 bindBallHost` 已绑 ⇒ `ra.ui.hasBall()` 此刻可读、且在全进程生命期**不再变**（`detachBall` 只在 `:806-807` 的退出序列） | ★〔甲〕**可行**（读 `hasBall()`） |
| **C** | `panel_inbound.go:272`（`Confirm` 在 `:277`） | ★**无对象**：`confirm=nil`（`:236` 传 nil、`:346-353` nil 进 nil 出）⇒ 372 的读侧降级**已经**生效 | ★**无需判**；现状即正解 |
| **D** | `panel_inbound.go:272`（同上），由 `resident_windows.go:224` 把腿交进去 | ★★**不存在**：`:224` 早于 `:298 startResidentBall`／`:309 bindBallHost` **85 行**；READ-H3 现测：`orbAtGateConstruction=false`。且 `panel_inbound.go:247-254` 自己写明"腿不能事后填"（Store 在 `New` 里拷贝 `Confirm`、handler 每请求只读自己的 nil-ness） | ★★〔甲〕**在 D 这条装配点上做不到**（除非把 `:224` 挪到 `:309` 之后，或把腿改成**活谓词**——那就是〔乙〕） |

★★**判语**：〔甲〕**可做但不是全覆盖**——A/B 能判、**C 不用判、D 判不了**；〔乙〕**四个全可行**，且对 B/D 而言**动作层那道闸已经存在**（`ballCardUI.Prompt` 无球即 `errResidentNoBall`，`resident_approval_windows.go:890-896` ⇒ `gate.go:536-540` 立刻 fail-closed）；缺的只是 **(i) 这个信号没接进 `PermissionMode()`**（`store.go:201` 只判 `confirm==nil`）、**(ii) 没落一行 `MODE-DOWNGRADE`**、**(iii) A 的 `Prompt` 打得出卡所以永远不报错**（`run.go:1533-1578`），A 需要一个**新**谓词（`reply != nil`）——而那正好是〔甲〕的数据。
★★**所以**：〔甲〕对 A/B 与〔乙〕等价（事实静态、读一次＝读每次）；〔甲〕治不了 D（事实在交接点还不存在）；〔甲〕**更治不动 READ-F**（见 ④）。⇒ **承重的一半必须是〔乙〕；〔丙〕＝可行**，但要把 D 的腿改成活谓词、A/B 的读数落成 boot 那行明话。
★★**停手上报的两条不触发**：①READ-F 与"超时成拒"**不矛盾**（④ 具名：两个对象）；②boot 判不出的**只有 D**，而**〔乙〕做得到** ⇒ 本票**不需要换形**，只需把 AC#2 的承重放在〔乙〕。

## ③ 「问不出」的全部条件（逐条，⛔ 三条）

1. **控制台（A）**：`interactiveStdin()`（`approval_reply_stdin_windows.go:41-52`）判＝`GetStdHandle(STD_INPUT_HANDLE)` 非 0/非 Invalid **且** `GetConsoleMode(h,&mode)` 成功；⇒ **管道／重定向文件／Explorer 无控制台／CI 全部 nil**（`:47-50`），非 Windows 恒 nil（`approval_reply_stdin_other.go:23`）。本进程现测 `interactiveStdinNil=true`（READ-H1）。nil ⇒ `attachReplyListener` 不装配（`run.go:857-859`）⇒ 卡照打、无人答。★**nil 时仍非 `nil` 的腿**＝本票的病。
2. **球窗（B/D）**：判＝`ballCardUI.currentBall() != nil`，**判的那一行**在 `resident_approval_windows.go:891-896`（`b := u.currentBall(); if b == nil ⇒ errResidentNoBall`）；绑定在 `:849-853 attachBall` + `resident_windows.go:309`，解绑在 `:806-807 detachBall`。READ-H3：交接点 `false`、`bindBallHost(nil)=false`（`:571-578` 的 false 分支才是"本进程没有悬浮球窗口"的权威句，`:309` **丢掉了它的返回值**）。
3. **ctx 先死**：`gate.go:576-580` ⇒ `abandon` ⇒ 拒（⛔ 不取消任务 root）。READ-H4 现测：控制台腿 **1ms** 拒（`任务上下文已结束，审批请求已作废并按拒绝处理`）、常驻腿 **1ms** 拒（`审批界面不可达…`）⇒ 两条腿、两个理由、同一 1ms。
4. **无人答**：`confirm_timeout_sec` 默认 **300**（`config/schema.go:460`）⇒ `Queue().Timeout()=5m0s`（READ-H2 现读）；带底 **31**、顶 3600（`config/validate.go:128-129,145`）；提示提前 **30s**（`queue.go:146-147`，READ-W 逐字）。L2 到点 ⇒ `expire` ⇒ **拒**（`gate.go:571-574`）；★**L1 到点 ⇒ 执行**（`gate.go:344` `AnswerTimeout`、`bridge.go:466-469`）——**"无人答"在 L1 上不是拒**。★另一枚会咬人的读数：L1 窗口配置 `1s` 被夹到 **2s**（`MinL1Window`，`queue.go:138`；READ-H2 `config.l1_window_sec=1`／`l1Window=2s`）。
5. ★**第五个条件（结构性的；⛔ 不是第五条装配——装配仍只有 A–D 四条）**：**A/B 这两条装配的 `Confirm` 在生产里没有任何调用者**。`rt.modeWrites`（`run.go:692-701`）**写明"生产无读者"**（`:331-345` 自述），生产 `Set` 只有面板 handler 两处；⇒ 这两条腿不是"问了没人答"，而是**根本不会被问**——档位只能从 `config.toml` 进场，而读侧降级因 `confirm != nil` 不触发。这就是 READ-F′ 的机制。
6. ★**L1 的窗口在控制台上也否决不了**：`run.go:632` 给 gate 的是空 `NewChannels()`、`main.go:159-164` 不设 `replyVeto`，`main.go:157` 自己就说"L1 窗口没有人能否决"；READ-S 实测 2s 到点直接执行。

## ④ 「超时成拒」 vs `READ-F`「写入成功」——**不矛盾，但票面那半句要改**

- ★**两件事量的是两个对象**：READ-F′ 那一发**根本没到 gate**（`cards=0`、无 `[确认…]` 行、`MODE-SILENCE … effective=L0` ⇒ 问题在 `bridge.go:319-322` 就被档位删掉了）；"超时成拒"量的是**L2 卡**（READ-W 同一 seam **31.024s ⇒ 拒**、READ-G2 1ms ⇒ 拒）。⇒ **不矛盾**。
- ★★**但票面把 A 的"问不出的形状"写成"超时成拒／ctx 先死当场拒"（票面 :33）是**漏了 L1**：在 A 上，"问不出"对 **L1** 的形状是**到点执行**（READ-S 2.04s ⇒ 写入成功）。所以"READ-F 写入成功"**⛔ 不是**"超时路径坏了"，而是**这条调用压根不是一等问题**。
- ★★★**这条直接改 `AC#2` 的可行性**：〔甲〕把有效档降到 `ask_every_step` 之后，READ-F 那一发**变成 READ-S**——**仍然写入成功**（2s 后）。⇒★**`AC#2` 的验收判据（"READ-F 从写入成功变成拒"）只有〔乙〕能满足**；用〔甲〕去做，`AC#2` 会红。

## ⑤ `374-a1` 的 `AC#1` 建议 ①–⑥ 逐条判（该不该做／属哪张票）

| # | 建议 | 该不该做 | 属哪张票 |
|---|---|---|---|
| ① | 可问性≠非 nil | ★**该做，但形式必须是〔乙〕**：boot 降级治不动 READ-F（④）；B/D 的信号已在 `Prompt` 里、只需接进 `PermissionMode()` 并落行；A 需新谓词 `reply != nil` | **375**（`AC#2` 承重项） |
| ② | 四个装配点逐条声明 | ★该做（**注释/装配处声明**），但 D 的声明必须写"活谓词"而不是 boot 快照，否则是把 `:224` 的假话写成合同 | **375** |
| ③ | 降级"当时"落行 | ★该做，**但只作为〔乙〕拒那一发的伴随行**（现无 tick、也没有"跨到降级"的那条边；374-a1 的 READ-H 已量到该跨点上 `MODE-DOWNGRADE` 仍 **0 行**） | **375**（随 ①） |
| ④ | 无答复源不许静默宽档 | ★该做：现在 `main.go:154-158` 只说了 L2 会超时，**没说 L1 窗口会无人否决并执行**（READ-S）⇒ 先把这句说全（文件层） | **375**（与 ① 同笔） |
| ⑤ | L2 卡印真实 deadline | ★该做；**我复现了 `窗口 0.0s`**（READ-W 逐字，`gate.go:534` 传 `win=0` × `run.go:1563`）⇒ 376 的立案前提**新鲜** | **376**（已立） |
| ⑥ | 记账面只许用有效档 | ★该做，但它与"能问≠非 nil"无因果（是 `approval_always.go:151-156` 把**配置档**写进 `Decision.Mode`） | ★**独立小票**（⛔ 塞进 375） |

## ⑥ 票面前提新鲜度（对 HEAD blob `8ced487bb` 复量）

| 票面断言 | 现量 |
|---|---|
| `store.go:201` 只判 `confirm == nil`（:14） | ★**仍真**（`s.confirm == nil && configured != ModeAskEveryStep`） |
| READ-F＝0 卡＋写入成功（:22） | ★**复现**（`cards=0 executed=true 38ms`） |
| READ-G0＝`MODE-DOWNGRADE` 0 行（:23） | ★**复现**（`downgradeLines=0`） |
| READ-G2＝1ms 当场拒（:24） | ★**复现**（`elapsed=1ms`） |
| 四条装配 A–D 与行号（:31-36） | ★**全部对上**：`run.go:672/674`、`resident_task_source_windows.go:271-283`、`panel_inbound.go:236/272/277`、`resident_windows.go:224`；生产 `Confirm:` 赋值 4 枚、`perm.New(` 2 枚 ⇒ **无第五条装配** |
| 「无第五条生产装配」（:38） | ★**独立复核为真**（另加：`Set(` 生产调用者 2 枚，都在同一文件） |
| `confirm_timeout_sec` 默认 300 ⇒ `5m0s`、带底 31、实测 `31.001s`（:55） | ★默认/带底**对**（READ-H2 `5m0s`；`31.024s` 是本发实测，deadline 就是 31s，+24ms 是开销）；★**但 `31.001s` 是 374-a1 那一发的数**，⛔ 当常量 |
| A 的"问不出的形状＝超时成拒／ctx 先死当场拒"（:33） | ★★**要改**：这两形只对 **L2** 卡成立；A 上 **L1** 的形状是**到点执行**（READ-S），且 A 的腿**无调用者**（③-5） |

## ⑦ 门禁 rc

- `go vet -overlay <探针> ./cmd/wisp` ⇒ **rc=0**；`go vet -overlay <探针+HEAD run.go> ./cmd/wisp` ⇒ **rc=0**。
- `go test -overlay <探针+HEAD run.go> ./cmd/wisp -run '^Test375a1…$' -count=1 -v`：`ReadF`／`StrictModeSameWrite`／`PageChainNoBall`／`BootFacts`／`OverwriteUnderAutoApprove` **5 发全 rc=0**（`PATH="$PWD/third_party/sherpa-onnx:$PATH"`、`GOFLAGS=`、只 `unset GOTMPDIR`）。
- 旁证那一发（脏树窗口内、未导出 HEAD）：`ReadF` rc=0、`StrictModeSameWrite` rc=0；`PageChainNoBall`／`BootFacts` 撞上半写态 rc=1（build failed）⇒ **故全部重跑到 HEAD blob**。

## ⑧ 提交与 `.go` 计数

- ★**本件两枚提交（⭐第一枚主题写错，已按"⛔ amend"改走追加）**：正文 `cd774a511`（**主题缺了"票 375 普查腿 375-a1"前缀**，内容无缺陷）；本枚主题修正 `git log -1 --format='%h %s' -- docs/evidence/s1/375-a1-census.md` 读回，主题逐字以 `票 375 普查腿 375-a1` 开头。⛔ amend，故正文不写死本枚哈希；读数基线 `8ced487bb`。
- **本笔 `.go` 改动＝0 枚**：`git show --name-only --format= <本件提交>` 只有本文件（已核 `cd774a511`：1 个 `A`、`.go` 0 枚）；5 发探针全在仓外、经 `-overlay` 注入。
