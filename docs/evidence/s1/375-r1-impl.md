# 票 375 落地腿 `375-r1`（写腿）：〔乙〕动作层 —— 「能问」≠「非 nil」，无人可答时 L1 不许「到点即执行」

- 腿 `375-r1`（写腿）。开工 HEAD `ad9785b28`（开工前核过：`git diff --cached` 空、`cmd/wisp`／`internal/agent/approval` 干净；工作树里别人的脏＝`.scratch/**` 与 `??` 怪名，未碰未清）。
- 纪律：未翻框、未改票名／`Status:`／`-done`、未 push、未碰台账。读数一律现跑；探针全在仓外 `D:\tmp\wisp\20261024\375r1\**`，经 `-overlay` 注入（仓内零探针字节）。
- ★**没碰 `cmd/wisp/run.go` 的 `consoleApprovalUI.Prompt`**：`git diff -- cmd/wisp/run.go | grep -c '^[+-].*Prompt'` ＝ **0**（与票 376 的 5/1 落点不重叠）。

## ① 为什么是〔乙〕（引普查 `375-a1`，⛔ 重做它的量）

★普查 §② 的 D 行原话：腿在 `resident_windows.go:224` 交出、球在 `:298/:309` 才绑 ⇒ **早 85 行**、`READ-H3 orbAtGateConstruction=false`、`panel_inbound.go:247-254` 自述「腿不能事后填」⇒ **〔甲〕boot 降级对 D 判不了**；而普查 §④ 判「**`AC#2` 的验收判据只有〔乙〕能满足；用〔甲〕做，`AC#2` 必红**」——本笔的 READ-S 复跑把这句话逐字证实了：最严档下那一发是 **1 卡、`2.027s`、仍然写入成功**（`ANSWER-EXPIRED … timeout->execute after=2s`）。
⇒ ★★**判据是「最有价值＋主流 harness」**（铁律）：主流 harness 的动作层在**调用那一发**判「有没有人答」，无则 `unavailable`／fail-closed，⛔ 在 boot 猜；〔乙〕买到的是**四个装配全覆盖**，〔甲〕买到的只有 A/B 的姿态可见。★**故本笔只落〔乙〕**。

## ② 动作层在哪、为什么必须在那层

| 落点 | 是什么 | 为什么必须在这里 |
|---|---|---|
| `internal/agent/approval/reach.go`（新，162 行）| `AnswerSource`（活谓词）＋ `Gate.answerable()` ＋ `UnsilencedMode`（档位读的包装）| 「有没有人能答这张卡」是**审批层自己的事实**；`options` 在 `New` 里拷贝 ⇒ 只有这一层能同时管住 **L1 窗口**与**档位读**这两个消费者，且 `Gate` 已有 `Window()/Queue()` 这类只读面 |
| `internal/agent/approval/gate.go`（+52/−3）| `Options.Answer`／`Gate.answer`／`PendingWindow` 的守卫 | L1 的极性（到点执行）写在 `gate.go:325-344`，⛔ 在桥或装配根改得动 |
| `cmd/wisp/run.go`（+73/−1）| `runSpec.answer`；`assembleRuntime` 里**一个**谓词喂两半；`tools.Options.Modes: approval.UnsilencedMode(...)` | 装配根是唯一知道「这条腿此刻有没有答复源」的地方（A：`rt.reply != nil`；B/D：`ra.ui.hasBall`），且谓词必须**活**（D 在 boot 判不了） |
| `cmd/wisp/resident_{approval,task_source}_windows.go`（+9／+6）| 两处都注入 `ra.ui.hasBall` | 与 `ballCardUI.Prompt` 自己的判据**同一个谓词**，守卫与界面不可能各说一套 |

★**⛔ 没落在**：`internal/perm`（372 刚改完、且 `PermissionMode()` 的降级是 372 的合同）／`internal/risk`（`Mode.Screen` 是 C19 之外的第二个判官，本笔只在**读**侧动手）／`internal/tools`（`b.permissionMode()` 的 panic 兜底在那里，包装器在它之内照样被 `recover` 覆盖）／`run.go` 的 `Prompt`（票 376）。

## ③ 改了什么（逐行）

1. `reach.go`：`type AnswerSource func() bool`；`Gate.answerable()`（**nil＝未声明 ⇒ 保持今天的行为**；`false` ⇒ 拒）；`Gate.AnswerSource()` 只读面（装配尺）；`UnsilencedMode(src, answer, logf)` ⇒ `PermissionMode()`：`answer()==false && 配置档 ≠ ask_every_step` 时返回 `ask_every_step`（＝372 的降级目标），并**每 episode 落一行** `approval: MODE-UNANSWERABLE configured=… effective=ask_every_step because=no-answer-source`（读侧每调用一次，⛔ 每调用落一行）。
2. `gate.go`：`Options.Answer`（L1 专用，注释写明 L2 为什么不碰）；`PendingWindow` 在 **R7 升级之后、开窗前**：`if !g.answerable()` ⇒ 落 `approval: L1-UNANSWERABLE … decision=reject because=no-answer-source` ＋ 返回 `AnswerReject`（**不打卡、不计时**）。到点执行那行的注释补一句：375 没有翻极性，只是在它**上游**拒掉一个不是问题的窗口。
3. `run.go`：`runSpec.answer`；`answer := s.answer`，无注入闸门时＝`func() bool { return rt.reply != nil }`（**两腿各自的唯一事实，⛔ 互相 OR**）；注入分支多两道**结构性**拒绝（缺谓词／闸门没有谓词 ⇒ 印话 ＋ 退 2，理由：档位降了而闸门照跑窗口＝洞还在）；console 闸门 `Answer: answer`；`Modes: approval.UnsilencedMode(rt.modes, answer, rt.auditf)`（**`rt.modes` 对象本身一字未动**：`perm.Store` 的 372 降级与行照旧）。
4. `resident_approval_windows.go`：闸门 `Answer: ra.ui.hasBall`（闭包，构造时球还没绑）；`resident_task_source_windows.go`：`answer: ra.ui.hasBall`。

## ④ `READ-F`：写入成功 → **拒**（唯一验收判据）

★改前（HEAD `ad9785b28`，同一探针同一发 `fs.write` 新文件、`configured=auto_approve`、`replyNil=true`）：
`READ-F write: cards=0 executed=true assessedLevel=L1 elapsed=26ms`；`MODE-SILENCE mode=auto_approve assessed=L1 effective=L0`；`call … decision=allow outcome=success`。
★★改后（同一发）：
`READ-F write: cards=0 executed=false assessedLevel=L1 elapsed=26ms`；`approval: MODE-UNANSWERABLE configured=auto_approve effective=ask_every_step because=no-answer-source`；`approval: L1-UNANSWERABLE corr=p375-f-corr tool=fs.write task=p375-f decision=reject because=no-answer-source`；`tools: call kind=refused … decision=reject outcome=error`；拒因逐字＝「本进程此刻没有答复源（没有可交互的控制台输入流，也没有悬浮球窗口），L1 确认窗口无人可答：「到点即执行」的前提是有人能看见并否决这张卡，无人可答时已 fail-closed 拒绝执行」。
⇒ ★**⛔「2s 后成功」、⛔「只出现降级行」**：`executed=false`、文件不在盘上（`Test375UnattendedConsoleRefusesTheL1Write` 同一发断 `os.Stat` 不存在）。

## ⑤ `READ-S`（最严档 ＋ 无答复源）：也必须变「拒」

改前：`cards=1 executed=true elapsed=2.027s`、`ANSWER-EXPIRED … decision=timeout->execute after=2s`、`[started] 确认窗口结束，未收到否决，开始执行`。
★改后：`cards=0 executed=false elapsed=21ms`（还原核那发 27ms）、`L1-UNANSWERABLE`、**零卡、零 `ANSWER-EXPIRED`**。
⇒ 到点执行那一支**没被删**（正控 §⑧ 仍是 `timeout->execute`），被删的是**不是问题的那种窗口**。

## ⑥ `READ-W`（L2／R1+R8）：仍「超时拒」

改前：`cards=1 executed=false assessedLevel=L2 elapsed=31.016s`、`ANSWER-EXPIRED … decision=timeout->reject after=31s`、`MODE-REDLINE mode=auto_approve still asks`。
★改后：`cards=1 executed=false elapsed=31.051s`（还原核 `31.027s`）、**同一条 `ANSWER-EXPIRED … timeout->reject after=31s`**、卡面同形（`窗口 31.0s`、`[warning] 审批将在 30 秒后自动拒绝`）。
★**诚实登记一处非零漂移**：`MODE-REDLINE … still asks` 这次**不出现**——因为有效档已是 `ask_every_step`，`Screen` 在「最严档什么都问」那一支就返回了，不需要红线豁免来解释「为什么还问」。`configured=auto_approve` 这个事实由上一行 `MODE-UNANSWERABLE` 携带。★L2 的**问与拒**一字未动（守卫只挂 L1 分支）。

## ⑦ 四个装配逐条（改前／改后）

| # | 装配 | 谓词（改后）| 改前 | 改后 |
|---|---|---|---|---|
| **A** | `wisp run` 控制台（`run.go:672` 一带）| `rt.reply != nil` | `auto_approve`⇒0 卡写入成功（READ-F）；`ask_every_step`⇒1 卡 2s 后写入（READ-S）| **两档都拒**（§④⑤），0 卡、`L1-UNANSWERABLE` |
| **B** | 常驻任务 runtime（`resident_task_source_windows.go:271-284` 注入 `ra.gate`）| `ra.ui.hasBall` | 无球时 L1 被 `ballCardUI.Prompt` 的 `errResidentNoBall` 拒（**这个面恰好诚实**）；但 `auto_approve` 下**根本不到闸门** ⇒ 同形的洞 | 无球 ⇒ 档位降到 `ask_every_step` 且闸门**自己**拒（`Test375ResidentGateRefusesTheWindowWithNoBall`：改前红「确认界面不可达（常驻进程没有悬浮球窗口，卡片无处呈现）」／改后绿「…无人可答…」）|
| **C** | `wisp panel-inbound` | 无（`confirm=nil`）| 372 的 nil-leg 降级已是正解 | **零改动**（无闸门、无 bridge；本笔不动 `internal/perm`）|
| **D** | 常驻面板档位链（`resident_windows.go:224`）| 同上 `ra.ui.hasBall` | ★**boot 判不了**（腿上得比球早 85 行）⇒〔甲〕做不到 | ★**活谓词每次调用现读** ⇒ 覆盖；球在＝照旧问、球不在＝拒 |

★D 的**强制面**（bridge 读的那一档）改后＝`ask_every_step`；★**展示面**（`perm.Store.PermissionMode()`／快照）仍是配置值 ⇒ 见 §⑬ 残留①。

## ⑧ 正控：有答复源时 L1 与 L2 **一字不变**

同一探针、同一发，只用 `io.Pipe`（活的、无人打字）当答复源：

| 读数 | 改前 | 改后 | 还原核 |
|---|---|---|---|
| `READ-PL1`（L1 窗口）| cards=1 executed=true **2.033s** `timeout->execute` | cards=1 executed=true **2.041s** 同句 | **2.038s** |
| `READ-PL2`（L2 未答）| cards=1 executed=false **31.023s** `timeout->reject` | cards=1 executed=false **31.024s** 同句 | **31.022s** |
| `READ-PL2A`（L2 答 yes）| cards=1 executed=true **268ms** `ANSWER-ALLOW route=native` | cards=1 executed=true **180ms** 同句 | **294ms** |

## ⑨ 新用例逐名 ＋ 为什么今天响

- `cmd/wisp/approval_unanswerable_375_test.go`：`Test375UnattendedConsoleRefusesTheL1Write`（2 子案 auto_approve／ask_every_step）· `Test375AttendedConsoleStillExecutesTheL1Write`（正控）· `Test375ResidentGateRefusesTheWindowWithNoBall` · `Test375InjectedGateMustDeclareItsAnswerSource`。
- `internal/agent/approval/reach_375_test.go`：`Test375L1RefusesAWindowNobodyCouldSee` · `Test375L1AnswerableKeepsTimeoutMeansExecute` · `Test375UndeclaredHostIsUntouched`（nil≠false，防「修成清洗」）· `Test375L2IgnoresTheAnswerSource`（边界 2）· `Test375UnsilencedModeOnlyDropsAModeNobodyCouldUse`（5 子案）· `Test375UnsilencedModeReArmsAfterItClears`。
- ★**今天响的实证**（把三条**不引用新 API**的用例放到 HEAD 上跑，用 `-overlay` 把本文件换成去掉新缝那一条的副本）：`Test375UnattendedConsoleRefusesTheL1Write` **FAIL ×2**（`the L1 write EXECUTED … auto_approve, 200ms`／`ask_every_step, 2.048s`）、`Test375ResidentGateRefusesTheWindowWithNoBall` **FAIL**（旧句 `确认界面不可达（常驻进程没有悬浮球窗口…）`）、`Test375AttendedConsoleStillExecutesTheL1Write` **PASS**（正控绿前绿后）。
- ★**诚实**：`reach_375_test.go` 与 `Test375InjectedGateMustDeclareItsAnswerSource` 引用新缝 ⇒ 在 HEAD 上是**编译红**（缝不存在），⛔ 充作「运行红」；运行红由上面三发承担。

## ⑩ 反向判据：三读数 ＋ 还原核（★核在最后一发变异之后量）

| 变异（仓外导出树＋`-overlay`，从**最终树**导出）| `READ-F` | `READ-S` |
|---|---|---|
| ★摘掉**模式读**那一跳（`Modes: rt.modes`）| **cards=0 executed=true 20ms**（回到写入成功）| cards=0 executed=false 28ms（守卫单独够用）|
| ★摘掉**闸门**那一跳（删 L1 守卫）| **cards=1 executed=true 2.025s**（回到「2s 后写入成功」）| **cards=1 executed=true 2.029s** |
| 两跳都摘 | **cards=0 executed=true 29ms** | **cards=1 executed=true 2.047s** |
| **不动（还原核，最终树）** | **cards=0 executed=false 48ms** | **cards=0 executed=false 27ms** |

★`READ-W` 在全部变异与还原核里都是 `cards=1 executed=false`＋`timeout->reject`（还原核 31.027s）⇒ 变异只动物理上属于 375 的那两跳。★还原核时机：**最后一发变异（mut12）之后**才跑（`restore-final.txt`），同一刻 `sha256sum` 五枚改动文件（`restore-sha.txt`，`run.go a69ef890…`／`gate.go cab0caa7…`／`reach.go be110e95…`／`resident_task_source e8c6eb83…`／`resident_approval 9d246ed1…`）⇒ 上面四行全出自**同一个最终树状态**。

## ⑪ 四数（改前／改后分开，⛔ 相减）

| 尺 | 改前 | 改后 |
|---|---|---|
| `cmd/wisp` **定向**（69 名，`PATH="$PWD/third_party/sherpa-onnx:$PATH"`）| 68 pass / **1 fail**（`TestAC246DevLegIgnoresTheTestTaskInjection`＝桌面态红，HEAD 同样红）| 68 pass / **1 fail**（同一条，同因）|
| ★超集（83 名＝69 ＋ 所有 `approval.New(` 文件里的用例）| —（未跑超集）| **83 pass / 0 fail** |
| `internal/agent/approval` | **78** top ＋ **19** sub ＝ 97 案 / 0 fail | **84** top ＋ **24** sub ＝ 108 案 / 0 fail |
| 动过的包之邻（`internal/tools`，零字节改动）| — | **250** top ＋ **87** sub ＝ 337 案 / 0 fail |

## ⑫ 门禁 rc ＋ 格式（甲乙两尺都报）

- `go build ./...` ⇒ **rc=0**；`go vet ./cmd/wisp ./internal/agent/approval ./internal/tools` ⇒ **rc=0**；`d22scan -root .` ⇒ **rc=0**（`clean - no D22 ban violations`）。
- 格式**甲尺**：`/d/work/base/gopath/bin/gofumpt.exe -l <12 枚改动 .go>` ⇒ **rc=0、输出空**；**乙尺**：`gofmt -l` 同 12 枚 ⇒ **rc=0、输出空**。⛔ 未用 `gofumpt -l .`（会扫 `.scratch` 故意写坏的探针）。
- 本笔 **12 枚 `.go`**：9 枚改动（+209/−12）＋ 3 枚新增（162／274／308 行）。
- ★已知非本笔红：`internal/risk` 的 `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（`internal/panel` 的两个读者）——HEAD 用 `-overlay` 复跑**同样红**；`internal/perm` 绿。

## ⑬ 诚实登记（⛔ 藏）

1. ★**展示面/记账面仍报配置值**：`perm.Store.PermissionMode()` 与面板快照是 `auto_approve`，而 bridge 读的是 `ask_every_step`（372 的降级在 `internal/perm`，本笔写面不含它；`374-a1` ⑤⑥ 已把「L2 卡印真 deadline」「记账面只许用有效档」分给票 376／377）。
2. ★`MODE-REDLINE` 行在「无人可答」场景消失（§⑥）——有效档已是最严档，红线豁免不再是它仍要问的理由。
3. ★**A 上 L1 窗口的否决通道仍是 0 枚**（`run.go` 给闸门空 `Channels`）——普查③-6；本笔的判据是**答复源**（`reply` 流），⛔ 否决通道；正控要求它一字不变，故未动。
4. ★**测试适配 9 处**（都是**给出答复源**，⛔ 放宽断言）：`run_test.go`（`runFixture.reply` 字段 ＋ `TestComposedGateBlocksAWriteForTwoSeconds`）· `panel_l1_blocked_220_test.go` · `run_mode101_test.go`（`t101host.start`，覆盖 3 枚）· `ticket224_assembly_test.go` · `resident_task_source_246_windows_test.go`（注入支的两枚正例加 `answer: ra.ui.hasBall`，球禁面那句改成 375 的句子）· `resident_approval_246_windows_test.go`（把「为什么拒」的句子换成 375 的，**结果断言一字未动**）。
