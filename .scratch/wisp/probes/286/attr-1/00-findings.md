# 票 286 — 只读归因腿 golden-attr-1 交件（AC#1 归因 + AC#3 名册指针）

取数时刻：2026-10-09 11:4x +0800（`date '+%F %H:%M %z'` 现量 11:41；正文另处时间记号亦来自同批现量）。
HEAD 现量 = `98482adc`（11:40，另一枚腿落的台账件；起手时是 `6041d7f9`）。`internal/agent/loop.go` 自 `bd124b2a` 起、`loop_golden_test.go` 自 `62dda11d` 起都未被再改 ⇒ 下面两态对拉的行号对 `bd124b2a..98482adc` 全段稳定。
零改动：只新建本 `.md` 与 `10-registry.md`；⛔ 没 checkout／stash／worktree／reset，没读工作树字节，没跑任何 go/脚本/gh/网络。AC#2（权威）归你已裁；⛔ 我没动 `loop_golden_test.go` 的任何期望。

---

## AC#1 — 请求侧 corr 从哪一跳开始带 `#call-` 后缀

### 边界那一枚 commit
**`bd124b2a`** — subject 起手 `242-corrland-1 落地：loop 每次工具调用铸独立 corr（callCorr=taskID#callID，任务级溯源仍走 taskID）…`，author date **2026-10-08 19:56:27 +0800**。它 `merge-base --is-ancestor bd124b2a HEAD` = **YES**（在 dev 主线上， reachable）。

### 两态对拉（对象层，逐字原文）
候选名册来源：`git log --oneline -- internal/agent/loop.go internal/agent/loop_golden_test.go`（从今天往前推到 10-08，起点覆盖到 ticket 10 的 `ae23da67`＝2026-09-20）。逐候选 `git show <ref>:<path>` 对拉，只有 `bd124b2a` 处发生状态翻转：

| ref | `loop.go` 请求侧 corr 那行 | 产码 corr 形状 | `loop_golden_test.go:69` 期望 | 请求侧那一枚绿/红 |
|---|---|---|---|---|
| parent `52501e22`（10-08 19:52） | `:654  TaskID: taskID, CorrelationID: taskID, CallID: p.call.ID,` | 逐字 = taskID（无 `callCorr`，`grep -c func callCorr` = 0） | `CorrelationID != res.TaskID` 在（:69） | **绿**（corr 就是 taskID） |
| **`bd124b2a`**（10-08 19:56） | `:676  TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i),` | `callCorr` 返回 `taskID + "#" + callID`（callID 空则 `%s#call-%d`） | `CorrelationID != res.TaskID` 仍在 :69（**未跟着搬**） | **红**（corr 永不等 taskID） |
| `98482adc`（HEAD，今日 11:40） | 同 `bd124b2a`（loop.go 未再改） | 同上 | 同 :69（golden test 未再改） | **仍红** |

⇒ 请求侧 corr 带 `#call-`/`#callID` 后缀的**起始落地跳 = `bd124b2a` 一枚**；红自它诞生、一直躺到今天（`loop.go`/`loop_golden_test.go` 二者自 `bd124b2a`/`62dda11d` 后无人再动 ⇒ 中间没有任何 commit 修过，也没人重新放大红）。

### 那枚期望自它诞生起钉的是什么（读代码，不按票号猜）
- `git log -S 'CorrelationID != res.TaskID' -- internal/agent/loop_golden_test.go` 的唯一命中 = **`ae23da67`**（`feat(agent): ReAct loop core … (ticket 10)`，2026-09-20 11:20）。⇒ 这枚等值期望**生于 ticket 10**。
- 同一 ref 读产码 `ae23da67:internal/agent/loop.go:599`：`			TaskID: taskID, CorrelationID: taskID, CallID: p.call.ID,`（且 `grep -c func callCorr` = 0，无 callCorr）。⇒ **它诞生时钉的就是"请求侧 corr 逐字等于 task id"，而且当时为真、为绿**（因为那阵子产码就是把 `CorrelationID` 直接写成 `taskID`）。
- 结论句：`:69` 不是一枚"错钉"，它当年钉的是**票 242 之前产码的真实形状**（corr=taskID）；242 把产码换成 per-call corr 后，这枚钉就变成"钉住旧缺陷"那一族——与你 AC#2 引的票 282 `AC#2` 先例（"那枚等值钉当年钉的是这个 bug 本身"，`A758`）同形，只是 242 那批只搬了 `loop_approval_test.go:213`、漏搬了同形的 `loop_golden_test.go:69`。**注意 `bd124b2a` 的 subject 自称"…loop_approval_test.go:213 唯一等值断言具名重判"** —— 这句"唯一"就是它漏网的根因（同仓请求侧还有 `golden:69` 这第二枚，只是那一跳没扫到）。

### 你转述的锚——三处对不上盘上原文，具名报回（以盘上为准）
1. ⚠ 你说断言在 `loop_golden_test.go:70`。盘上 HEAD：`:69` 才是那枚 `if call.Req.TaskID != res.TaskID || call.Req.CorrelationID != res.TaskID {`；`:70` 是 `t.Errorf(...)`。FAIL 输出里的 `:70` 是 Go 报 `t.Errorf` 的行，不是 `if` 的行。票面"现量"与 AC#2 也写成 `:70`——同一处，写腿动手时锚应按 `:69` 找那条 `if`。（`:341` 你对了。）
2. ✓ 其余锚全部对上原文：`loop.go:603 callCorr` 定义、`:676` 唯一调用点（形状 `taskID + "#" + callID`）、`:369 newTaskJournal(l.opt.Journal, taskID, taskID)`、`journal.go:59` 签名 `newTaskJournal(j Journal, taskID, corrID string)`、`journal.go:81 CorrelationID: t.corrID` —— 逐字属实。
3. 我**没**独立核 `docs/PLAN.md:1368` 的字面（读数留给权威格）；归因不依赖它。我引的两处对象层旁证是 `loop.go:594-601`（callCorr doc）与 `loop.go:363-368`（journal 注释），二者都自指 `docs/PLAN.md:1368` 的 C18 句"corr 坐在 task id 旁边、不要求等于它"，与 AC#2 一致。

---

## AC#3 — 名册（详见 `10-registry.md`，此处只给头条）

尺覆盖了等值断言／`==` 正形与反形（两 corr 比相等）／`CorrelationID: task-token` 构造／`CorrelationID = task-token` 赋值／位置参数 `taskID, taskID`／`orDefault(Corr,Task)` 回落／注释里的 `TaskID == CorrelationID` 句——不是只扫一种字样。真树名册 **25 枚**：

- **该搬 = 1**：`HEAD:internal/agent/loop_golden_test.go:69` `	if call.Req.TaskID != res.TaskID || call.Req.CorrelationID != res.TaskID {`
- **该留 = 19**：`:341`(账本)、`panel_pump_test.go:323`(结果流)、`ticket285_corr_distinct_rulers_test.go:119`／`loop_approval_test.go:219`／`ticket283_corr_identity_rulers_test.go:222`／`ticket285_corr_rows_rulers_test.go:72`(4 枚 242 后正确尺)、`bridge.go:270`(空才回落)、`loop.go:369`(账本残余)，+ 9 枚宿主 fixture（`subagent_blocked_197_test.go:121`/`subagent_carrier_197_test.go:173`/`task_scope_close_151_test.go:39`/`bridge_scope_open_ticket158_test.go:39`/`pointer_183_cli_seam_test.go:127`/`pointer_185_cli_seam_test.go:120`/`task_cancel_221_legs_test.go:119`/`ticket175r2_stamp_live_test.go:95`/`failclosed_236_teeth_test.go:145`）+ `internal/panel/pump_test.go:57` + `declared_l0_risk_179_test.go:204`（账本侧）。逐条带原文见 10-registry 三/二节。
- **待人裁 = 4**：过期注释句（把"loop 派发 corr==task"当事实写，242 后为假）——`internal/panel/subagent_roster_197.go:57`、`internal/panel/subagent_roster_197_test.go:468`、`internal/panel/subagent_blocked_220_test.go:155`、`cmd/wisp/subagent_carrier_197_test.go:20`。都不是红、不在请求侧、跨 `internal/panel`/`cmd/wisp` 包 ⇒ 不建议写腿在 AC#4 那一笔里顺手改。
- **冻结 out-of-scope = 1**：`internal/perm/ticket90_persist_test.go:123`（三枚冻结件之一，git grep 顺带命中，未整读、不动）。
- 探针桶（不算动作）：`.scratch/**` 的 corr==task 家族 **162 行 / ~29 件**死突变/探针副本。

**"第二枚同形会不会漏"这一格的答案：不会。** 全仓请求侧"要求 corr==taskID"的断言只有 `golden:69` 一枚；其余同族要么是账本/结果流侧（`:341`、`panel_pump:323`，今天都靠 task-keyed 产码为绿、且**该留**），要么已是 242 之后"必须 ≠taskID"的正确尺（4 枚），要么是宿主 fixture 的合法输入（绿、不动）。⇒ 名册给写腿的净动作面只有 `golden:69` 一枚（＋你那 4 句注释要不要一起订正，等你裁）。

---

## `:341` 那一枚的单独判语
`HEAD:internal/agent/loop_golden_test.go:341` 逐字 `if row.Tool != "echo" || row.CorrelationID != res.TaskID {`。
它绿是**对的**：`:333` 取 `store.ListToolCallsByTask`，那行的 `CorrelationID` 由 `loop.go:369 newTaskJournal(l.opt.Journal, taskID, taskID)` 的第二实参经 `journal.go:59→:81 CorrelationID: t.corrID` 落成 taskID ⇒ 账本侧今天就是 taskID。这笔残余已由票 282 `AC#4(a)` 裁"留"（`A758`）。⛔ 不许为了和请求侧"看起来一致"去动 `:341`／`loop.go:369`／`journal.go:81`——动它＝悄悄改账本侧形状。名册里与之同侧的账本形（`declared_l0_risk_179_test.go:204/227/248` 三发 `newTaskJournal(...,task,task)`）同样按"账本侧、该留"判。

---

## 量不到的格子 + 缺什么读数（要真跑的分开写）
1. ⚠ **"今天整树只 `internal/agent` 一枚红"我没能在只读层证实**：你 11:2x 只跑了 `./internal/agent/ ./internal/tools/`（票面现量）。`cmd/wisp`（`panel_pump_test.go:323` 住这儿）、`internal/panel` 从没在这轮跑过。⇒ 我对 `:323` 判"绿、该留"的凭据是**对象层产码链**（`run.go:1266 Append(e.TaskID)`→`pump.go:526 CorrelationID:key`），**不是实跑 rc**。要坐实"cmd/wisp 无第二枚红"，缺的读数＝一次 `go test ./cmd/wisp/ ./internal/panel/ -count=1`（⛔ 归写腿/在飞 `282-r1` 独占 Go 面，我不跑）。
2. ⚠ 名册里 9 枚宿主 fixture（`#10-#18`）与 `#19/#20`：我判"绿、该留"是基于"它们把 corr 直接喂 `bridge.Execute`、不经 loop，且 prod 允许 corr==task"。**没逐枚跑过**；若你要求逐枚 rc，缺同一次 `./internal/tools/ ./cmd/wisp/` 的 -v 名册。
3. ✓ 能静态定的我都定了：边界 commit（两态对拉）、`:69` 诞生形状（`ae23da67` + `loop.go:599`）、`:341`/`:323` 的 task-keying 产码链——这些不需要跑。
4. ⚠ 我给不了"搬完 rc=0"那格：那是 AC#4（写腿＋反形自证＋`go test`），按规矩在名册之后、由在飞 Go 独占腿做。

## 我顶回你的话（直说）
- **顶回一处行号**：`:70` 应为 `:69`（见上）。票面 AC#2 正文也带着这个 `:70`，写腿别按 `:70` 去搬。
- **不顶回你的"怀疑起点"**：这次你怀疑的起点（票 242 的 correlation 落地那批）**按尺判是对的**——边界确为 `bd124b2a`（242-corrland-1）。我没有为反对而盖章，是两态对拉独立落到同一枚；顺带把你 AC#2 的权威判（搬期望、不回退产码）在归因层印证了一遍。
- **补你一条**：`bd124b2a` subject 自称"`loop_approval_test.go:213` 是**唯一**等值断言"——这个"唯一"就是它漏掉 `golden:69` 的根因；建议把这句也记进 242 的未收口账，避免下一位再信"唯一"。
- 名册里那 4 枚**过期注释**（`subagent_roster_197.go:57` 等）＝我判"待人裁"、不并进 AC#4 那一笔；要不要同批订正由你裁（改它们跨 `internal/panel`/`cmd/wisp` 包，动别处包注释不属"搬 golden:69"）。
