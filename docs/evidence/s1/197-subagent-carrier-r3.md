# 197-r3b 载体层证据件 — 名册与每枚子代理的流真进快照（续前一腿 372 行，非重写）

> 票 197 第 3 层（载体层）的**续腿**：前一腿 `197-r3` 撞到 150 轮上限，产码全在工作树里未提交、证据件没写。
> 本腿（`197-r3b`）的活是**接着做完**：判据落在生产路径、补上"拔掉泵必红"的反向正控、提交并写这份证据件。
> 派单：`.scratch/wisp/dispatches/2026-09-28-212x-impl-197r3b-finish-the-uncommitted-carrier.md`（工单正文＝同目录 `2026-09-28-210x-impl-197r3-*.md`）。
> 生成时刻 `2026-09-28 21:5x +08`。**本文件只写盘上可核的事实**；所有读数出自本腿亲手跑的命令，探针输出全在
> `.scratch/wisp/probes/197/r3b/`。

---

## ① 起手锚点与名册（本腿现取，未采信派单里的任何行号）

| 项 | 读数 | 尺 |
|---|---|---|
| 起手 HEAD | `2f18768d`（编排者的 A414 台账＋派 197-r3b 那一枚） | `git log --oneline -1` |
| 起手工作树里的本票产码 | 6 枚已跟踪文件 `M`（372 插入／33 删除）＋**4 枚未跟踪文件**（不是派单说的 1 枚，见 §⑥-1） | `git diff --numstat`／`git status --short` |
| 起手包读数：`internal/tools` | **`ok 17.240s`** | `go test -count=1`（带 sherpa/build PATH） |
| 起手包读数：`internal/panel` | **`FAIL 2.508s`**＝在册那 4 枚（逐名见下） | 同上 |
| 起手包读数：`cmd/wisp` | **`FAIL 92.418s`**＝红 **1** 枚，**不是在册红**，是前一腿产码自带的一枚跨瞬断言 | 同上，`.scratch/wisp/probes/197/r3b/start-gate-3pkgs.txt` |

起手那枚 `cmd/wisp` 红的逐字读数（本腿复现，非引用）：

```
--- FAIL: TestRunPacketCarriesTheSubagentItsRosterRowFed
    subagent_carrier_197_test.go:362: the run-driven packet already reported 0 slots and a settled row while the spawn call was still open: 1
    subagent_carrier_197_test.go:370: packet tasks section: rows=2 poolCap=4 child=f1ee91f7-... runStatus=Thinking afterStatus=Settling key=subagent:f1ee91f7-... bytes=2789
```

⇒ **本腿开写面之前先有的结论**：载体层的**产码不缺跳**（快照里有那枚子代理、带它的 taskID／D43 状态／父任务／它自己的流键，
正文也确实进了它那一页），**缺的是那一枚断言本身写错了**：它拿"run 自己发布的包里读到的 `inFlightSlots`"（一个**中途**的瞬时）
去比"spawn 调用返回之后再读一次的状态"（**另一个**瞬时）。整包跑（负载高、run 最后一次发布落中途）就红，
`-run` 单跑（前一腿 21:42 那一发）就绿——它自己那份 `cmd-wisp-carrier-run1.txt` 与 `final-v-four-pkgs.txt` 就是同一棵树上的两枚相反读数。
修法见 §②(末)，**没有放宽任何断言**：改成"同一枚包内的两半不许互相矛盾"，那才是这条不变量真正说的东西
（`internal/tools/subagent_197.go:352-355` 的注释写明池位**先还**、行**后**定稿，所以"包里写着 Settling 却还占着位"才是矛盾）。

**在册红名册（起手逐名，本腿一未修、二未 Skip、三未放宽）**：`internal/panel` 4 枚——
`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／
`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourwayAgree`。
逐包绿点数见 §④（本腿在终态用 `-v` 现数，前一腿 21:12 那份 `start-roster-counts.txt`（panel 90／tools 158／wisp 93）**不算本腿的读数**，只作对照）。

---

## ② 三格各自的改前 → 改后读数＋落点逐跳行号

⚠ 行号一律**本腿现读**（`git show HEAD:<file> | grep -n`，锚 `335b8d2b`），不抄派单、不抄前一腿证据件。
⚠ 本文件里 `cmd/wisp/subagent_carrier_197_test.go` 的行号**锚 `335b8d2b`**；第二枚 commit `7e3e7663` 在第 116 行处加了 3 行读数日志，
其后每一枚行号在 `7e3e7663` 上各 **＋3**（`:275→:278`、`:320→:323`、`:339→:342`、`:412→:415`、`:545→:548`）。
§③ 引的红句是**变异那一发自己打出来的行号**（跑在 `335b8d2b` 上），照原样抄，不改。
"改前"＝`2f18768d` 那棵树上载体层不存在的形；本腿用 §③ 的单点变异把那一形**在同一棵树上复现并打出字节**，
所以改前读数不是引用，是本腿亲手跑出来的。

### (a) 名册进快照：面板能列出"现在有哪几枚子代理、每枚什么状态、是谁派出来的"

| | 读数 |
|---|---|
| **改前（字节级）** | 泵**没有名册读者**时，线上只有 5 枚顶层键：`pending,results,composer,generatedAt,instructions`——**没有 `tasks`**。逐字读数＝`.scratch/wisp/probes/197/r3b/m1-red-cmdwisp.txt` 里 `packetTasks197` 打出的整包（那一发就是把 `cmd/wisp/run.go:491` 的 `Tasks: rt.taskRosterState,` 换成 `Tasks: nil,`，即前一腿产码入库之前生产路径的真形） |
| **改后（生产路径）** | 同一枚 run 自己发布的包里有那枚孩子：`rows=2 poolCap=4 child=<uuid> runStatus=Thinking afterStatus=Settling key=subagent:<uuid> bytes=2789`（`cmd/wisp/subagent_carrier_197_test.go` 的 `t.Logf`，逐名读数见 §③/§④），且 case 逐条断言 `taskId/kind/parentTaskId/label/statusKnown/status∈D43/streamKey`＋根的同一行 |
| **落点** | `internal/panel/composer.go:91` `Tasks *TaskRosterSection json:"tasks,omitempty"` → `internal/panel/pump.go:157` `Tasks func() TaskRosterState`（读者位）→ `internal/panel/pump.go:256-263` `if p.src.Tasks != nil { snap.Tasks = TaskRosterSectionFrom(p.src.Tasks(), cards) }` → `internal/panel/subagent_roster_197.go:176` `TaskRosterSectionFrom`（渲染＋排序＋fail-closed）→ 装配根 `cmd/wisp/run.go:491` `Tasks: rt.taskRosterState,` → 读者本体 `cmd/wisp/panel_pump.go:146` `taskRosterState()`（枚举行＝`admitTask` 见过的 id ∪ `Roster.Descendants`，状态＝`TaskOutput.StateAnswer()`，池位＝`InFlightSubagents()`，天花板＝`tools.MaxConcurrentSubagents`，截断三读数＝本 run 正在写的 `StreamLog`）→ 枚举集合的生产者 `cmd/wisp/run.go:266 seenTasks`＋`:658 rt.noteTask(taskID)`（挂在 `admitTask` 上，根环与孩子环都过这一枚钩子） |

**键的选择（派单 §2(a) 要写清的那一条）**：走的是**嵌进现有 `Snapshot` 的一枚 `omitempty` 指针**，**不是**无条件顶层键。
⇒ `internal/panel/pump_test.go:123` 与 `:291` 那两枚"四枚顶层键"的钉**一字未动、今天仍绿**（终态读数见 §④）。
⚠ 但那一绿**不能算载体层的功劳**：派单 §2(a) 自己点破了——那两枚钉读的是没填 `omitempty` 段的 fixture，**加了顶层键也看不见**。
本腿因此把**能看见的那一半**钉成常驻判据，而不是去动那两枚钉：
`internal/panel/subagent_roster_197_test.go:194 TestAPumpWithoutARosterReaderSendsFourKeys` 两向都钉
（**没有**读者 ⇒ 恰好四枚、且**不许**有 `tasks`；**有**读者 ⇒ 一枚孩子都没派也必须发这枚键，
"这一程还没派过"与"这枚宿主看不见子代理"在线上必须是两句话）。

**契约面这一跳本腿做了什么**：新增 `tasks` 键**确实**把 `TestApprovalCardViewJSONKeysMatchFrontendTypes` 与
`TestComposerContractTypesMatchFrontend` 两枚在册红**的红因加长了一枚**（现读逐字见 §④：`Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`）。
**没改尺、没加豁免、没为过尺少发一个键**；界面侧要补的是 `frontend/src/lib/panel.ts` 里 `PanelSnapshot` 的 `tasks` 一枚与
`TaskRosterSection`／`TaskRowView` 两张表的字段（逐字段形状见本文件末"要交给界面那支的键集"一节）。

### (b) 每枚子代理的流进快照：点进去那一页有东西可渲染

| | 读数 |
|---|---|
| **改前** | 孩子环路**没有 sink**（`opt.Sink = nil` 没人覆盖）⇒ 它流出的正文被丢掉，那一页只剩生命周期两行。逐字红句（本腿 §③-M3 亲手跑出来的）：`the child's own streamed reply never reached its own row: "已派生：载体层正控：把一句话原样说出来"`。另外 `Close` 那一半当时不存在 ⇒ 收口的孩子仍 `done=false`，页面上永远在流 |
| **改后** | 同一枚 case 断言三件事：孩子那一页有 `echo:`（它自己流出来的正文真到了它自己的行）、**没有**根的 prompt 那个记号（两枚任务没被并成一行）、收口之后 `done=true`；截断那半＝case 2 的 `bound crossed: rows=34 streams=34 truncated=true elided=0 dropped=[]` |
| **落点** | `internal/tools/subagent_197.go:153` `Stream SubagentStreamSink`＋`:160` 那两枚方法（`Append`/`Close`）→ `:287` `opt.Sink = subagentTextSink{d: t.d}`→ `:481` `subagentTextSink.Publish`（只把 `EvTextDelta` 追加、`EvDone` 收口，推理增量不进结果道，理由沿用 `consoleSink`）→ `:443/:460` `feed`/`feedDone`（空 task id ⇒ 键是 `""` ⇒ **一个字都不写**）→ 装配根 `cmd/wisp/run.go:541` `Stream: rt.stream,`（把 `*panel.StreamLog` **本体**交出去，不再是一层闭包，因为 `Close` 是活的一半）→ 读侧 `internal/panel/pump.go:628 TruncationFor`（单条流的代价）＋ `cmd/wisp/panel_pump.go:198-204`（节级三读数 `InFlightSubagents/Truncated/ElidedRunes/DroppedKeys`） |

**流键的单一铸造点**（派单 §2(b) 要的那一收敛）：`internal/streamkey/streamkey.go:29/:39` 现在是全仓唯一写 `"subagent:"`
字面量的地方，`internal/tools`（`subagent_197.go:53`）与 `internal/panel`（`pump.go:430`）两枚常量都**改成了别名**；
前一腿留下的两枚"两包相等"钉（`cmd/wisp/subagent_stream_key_197_test.go`）本腿没拆牙，另加一枚常驻的
`TestSubagentStreamKeyHasOneMintSite`（`cmd/wisp/subagent_carrier_197_test.go:412`，走 `go/ast` 数**字符串字面量**、注释不豁免）。
⇒ **顺手关掉一形前一腿记在 §⑤ 的洞**：空白 task id 时 panel 侧返 `""`、tools 侧返 `"subagent:"`，两边不一致；
现在两边都调 `streamkey.Subagent`，同一枚函数（`internal/streamkey/streamkey.go:39`）。

**"少了多少"必须是数字**：节级 `streamTruncated/streamElidedRunes/droppedStreamKeys` ＋ 行级
`streamTruncated/streamElidedRunes/streamDropped` 都上了字节。**没测到什么**在这一格是硬事实，本腿复核过（不是抄）：
mockllm 把答案截在 400 **字节**、环路又把 scene 块排在前面 ⇒ 单枚孩子的整行约 150 rune，**过不了 512-rune 的保留窗**，
所以真装配那一发只能是 `truncated=true elided=0`（真值，不是缺读）；"一行真少了中间"的算术在
`internal/panel/subagent_roster_197_test.go:331 TestTruncationFactsRideThePacket` 里钉着（真 `StreamLog`、真泵、
真字节，逐字断言字节里出现 `"streamElidedRunes":88`）。

### (c) 状态词只许 D43 那 20 枚；取消不级联

- 载体层**不判状态**：`internal/panel` 不 import `internal/statemachine`，`subagent_roster_197.go` 只把宿主
  `TaskOutput.StateAnswer()` 的三半带过去。这是**被迫的**、且判据在场：`TestTaskState188AC3PanelHasNoWriteLeg`
  会遍历本包每一枚非测试源件找"面板侧写状态"的那根线（前一腿把它写进注释，本腿复跑确认它仍绿）。
- 载体层唯一**强制**的一件事＝被 judging 拒绝的名字不许当成状态上线（`internal/panel/subagent_roster_197.go:215`
  那个 `if !view.StatusKnown`：把值**丢掉**而不是翻译，翻译就是票 196 防的第二套词表）。这一条本腿装了牙并跑红（§③-M2）。
- 生产路径上的状态维：case 1 直接拿 `statemachine.Valid(statemachine.State(child.Status))` 判线上那枚名字；
  根行那一维今天**没人填**（裁定 `A394` 把那根线留给票 196），所以断言的是"空**并且带一句原因**"，
  而不是替它编一个——将来谁把根的状态填上，这一条就是会**先红**的那一条，不会悄悄变绿。
- **取消不级联**：本腿没动这条语义、没重写它绕开它，载体层也不新增任何取消路径；常驻判据仍是
  `internal/tools/subagent_197_test.go:716 Test197CancelIsPerRowAndNeverCascades`（终态读数 §④：tools 整包 158 PASS／0 FAIL）。

### (末) 本腿唯一改动的那一枚判据（前一腿 `cmd/wisp/subagent_carrier_197_test.go:361-364`）

改前那三行的条件写成 `sect.InFlightSlots != 0 && joined.Status == Settling`：左边来自 **run 自己发布的那一枚包**（中途瞬时），
右边来自 **spawn 调用返回之后再读一次**（另一个瞬时）——两枚瞬时之间孩子可以随便收口，所以它不是一条不变量，
而是一枚掷硬币（整包红、单跑绿，§① 那两枚相反读数就是它）。红句自己还写着"already reported 0 slots"，
而它的条件要求的是 `!= 0`，**报的不是它看见的东西**。
改后＝**同一枚包内的两半不许互相矛盾**（`child.Status == Settling` 与 `sect.InFlightSlots != 0` 同时为真才算红），
这正是 `internal/tools/subagent_197.go:349-350` 注释许诺的那条序（**池位先还、行后定稿**）。
⇒ 没有少断言任何东西：孩子收口后读到的必须是 `Settling`（`:351`）、池位必须清空（`:356`）两枚原样留着。
**这一改动是本腿在判据上的唯一改动，且它让整包从红变绿**（`gate-cmdwisp-after-assert-fix.txt`：`ok 103.731s`）。

---

## ③ 四发变异／正控（本腿亲手跑，读数逐字）＋"哪一枚 commit 才算未修码"

台件：`.scratch/wisp/probes/197/r3b/mutate-197r3b.sh`（单点、串行、每发跑完立刻
`git cat-file blob HEAD:<path> > <path>` 还原，再用 `git status --porcelain` 自证还原干净）。
逐发输出＝同目录 `m0-*`／`m1-*`／`m2-*`／`m3-*`／`m4-*.txt`，汇总＝`mutation-log.txt`。
**基线（未变异）先跑绿**：`m0-baseline-cmdwisp rc=0 ok 3.776s`、`m0-baseline-panel rc=0 ok 0.054s`。

| # | 单点改动（落点） | 结果 | 逐字读数 |
|---|---|---|---|
| **M1** 把泵那根线拔掉 | `cmd/wisp/run.go:491` `Tasks: rt.taskRosterState,` → `Tasks: nil,` | `cmd/wisp` **红 2 枚**（`rc=1`）；`internal/panel` 五枚载体判据**全绿**（`rc=0`）⇒ **包内判据对装配回归是瞎的，只有生产路径那一头发牙** | `--- FAIL: TestRunPacketCarriesTheSubagentItsRosterRowFed (1.91s)`／`--- FAIL: TestRunPacketReportsTheStreamLogPastItsBound (1.67s)`；两枚都停在 `subagent_carrier_197_test.go:275/545: the packet carries no tasks key at all: {"pending":[],"results":[...]...,"instructions":{...}}`——打出来的整包**确实只有 5 枚顶层键、没有 `tasks`** |
| **M2** 让被拒绝的状态名照样上线 | `internal/panel/subagent_roster_197.go:215` `if !view.StatusKnown {` → `if false {` | `internal/panel` **红 1 枚** | `--- FAIL: TestAStatusOutsideD43NeverReachesTheWire`：`subagent_roster_197_test.go:298: status = "running": a name the host's judge refused still reached the wire, which is the second vocabulary ticket 197 §2(c) forbids`＋`:322: a blank status with no reason is the shape StateAnswer exists to prevent` |
| **M3** 孩子不留自己的 sink | `internal/tools/subagent_197.go:287` `opt.Sink = subagentTextSink{d: t.d}` → `opt.Sink = nil` | `cmd/wisp` **红 1 枚**，且**只红在流那一跳**（名册那半照旧绿）⇒ 名册与流是两枚独立的牙 | `--- FAIL: TestRunPacketCarriesTheSubagentItsRosterRowFed (1.85s)`：`subagent_carrier_197_test.go:339: the child's own streamed reply never reached its own row: "已派生：载体层正控：把一句话原样说出来" - this is the hop the carrier leg exists to close`（那一页只剩生命周期一行） |
| **M4** 把流键字面量再铸一枚 | `internal/panel/pump.go:430` 别名 → `= "subagent:"` | `cmd/wisp` **红 1 枚**（铸造点扫描） | `--- FAIL: TestSubagentStreamKeyHasOneMintSite (0.12s)`：`the stream key literal "subagent:" is minted in 2 non-test sources [internal/panel/pump.go internal/streamkey/streamkey.go], want exactly one: internal/streamkey/streamkey.go` |

每发还原后的复跑（同一台件里紧跟着跑）：`m1-restored-cmdwisp rc=0`／`m2-restored-panel rc=0`／
`m3-restored-cmdwisp rc=0`／`m4-restored-cmdwisp rc=0`，四枚 `RESTORED_CLEAN`，脚本末
`git status --porcelain -- cmd internal` **输出为空**。

**"哪一枚 commit 才算未修码"**（派单点名要写清的）：
- 载体层**不存在**的那棵树＝ **`2f18768d`**（本腿起手 HEAD）。未修码的**线上形状**本腿没有用 checkout 去打
  （共享工作树里禁 `reset`/`checkout .`），而是用 **M1 在同一棵树上复现**：那一发打出来的整包就是"泵没有名册读者"时
  生产路径真会发的字节——`tasks` 一枚都没有。
- 载体层**入库**的第一枚 commit＝ **`335b8d2b`**（本腿，10 枚文件、1748 插入／33 删除）。它已经把前一腿那枚跨瞬断言写对，
  所以 `335b8d2b` 上**整包是绿的**（`gate-cmdwisp-after-assert-fix.txt ok 103.731s`）。
- 本腿稍后的**第二枚** commit＝**判据补牙**（`internal/tools` 的 `Close` 断言把前一腿留下的那枚**没人调用**的
  `streamClosed` helper 变成活的判据；`cmd/wisp` 加一行把 `tasks` 那节的**原始字节**打进 `-v` 日志，
  这样 §② 引用的改后读数是从 run 发布的字节里取的，不是本文件自己再 marshal 一遍）。
  ⇒ 在这一枚之前，"孩子收口后那一页必须 `done=true`"在 tools 层**没人钉**（只有 cmd/wisp 那一发钉着）；
  它算未修码的是**判据面**、不是产码面：产码（`feedDone`）在 `335b8d2b` 已经在，红不到。

正向读数（改后、生产路径，本腿自己的 `-v` 跑）：

```
packet tasks section: rows=2 poolCap=4 child=bcfd9825-... runStatus=Thinking afterStatus=Settling key=subagent:bcfd9825-... bytes=2789
tasks wire bytes: {"rows":[{"taskId":"4dc60786-...","label":"总结一下 rootprompt197r3 aaa…,"kind":"root",...
```

## ⑤ 我没测什么（具名，不留空）

1. **`blockedOnApproval` 在真装配上没钉。** 它只在 `internal/panel` 用一枚真 `ApprovalCardView` 钉住那个 join
   （`subagent_roster_197_test.go:103`），生产路径那发只测到"没有卡的行不许报 blocked"这一向。
   ⚠ 前一腿 §未测 1 那句"**这棵树里没有任何办法从外面造出那枚卡**"**说重了**，本腿现读更正：宿主自己就有开卡的线——
   `cmd/wisp/run.go:577-579` 为模式切换走 `rt.gate.AdmitTextTask(taskID)` ＋ `rt.gate.PendingApproval(ctx, tools.Decision{...})`
   签一枚真 L2 卡；把 `CorrelationID` 换成**孩子的** task id 就是 AC#4 要的那一发（难点只剩"孩子的 id 要在它活着时读到"，
   `taskRosterState()` 每 5ms 轮询即可）。**本腿没做**：派单 §1 写死"三件，按顺序，别扩范围"，而这一格不在那三件里。
   **要不要那一发＝编排者一句话**，代价＝cmd/wisp 里 20-30 行测试、一次 4 秒的 `-run` 复跑，产码零改动。
2. **"一行真少了中间"在真装配上不可达**（本腿复核过、不是抄）：mockllm 把答案截在 400 **字节**、环路把 scene 块排在前面
   ⇒ 单枚孩子的整行约 150 rune，**过不了 512-rune 保留窗**，所以 `bound crossed: rows=34 … elided=0` 是**真值**。
   算术那半在 `internal/panel` 的 `TestTruncationFactsRideThePacket` 钉（真 `StreamLog`、真泵、字节里出现 `"streamElidedRunes":88`）。
3. **`droppedStreamKeys` 在真装配上没点过名**：硬天花板＝`DefaultStreamKeys × 2`＝64 枚**真孩子环路**，33 枚已经是本套件最慢一发。
   那一半同样在 panel 层钉（`TestDroppedStreamsAreNamedOnTheWire`）。
4. **点击那一跳没测**：所有读数停在**包字节**。传输仍归票 33（`Snapshot` 出向那一条以外没有接收器、没有 router、没有 WebView2 宿主——
   包里的 `composer.git.switchBlocked` 那句宿主原话就是这条边界的自述）。没有任何一枚判据证明"页面渲染过一行"。
5. **名册的字节规模上限没量**：饱和名册（4 枚孩子各留 512-rune 窗）× 包大小 vs 账本那枚上限，全仓无尺；
   前一腿 §未测 4 记的这条本腿没动。⇒ 若界面那支要一次拉全部历史行，这条得先量。
6. **根任务那一行的状态维今天没人填**（裁定 `A394` 把那根线留给票 196），所以"根行也带一枚 D43 名"这条在生产上没法钉；
   本腿钉的是**空并且带原因**（`subagent_carrier_197_test.go:320`）。将来谁填上，那一枚断言就是**先红**的那一枚。
7. **"孩子正在跑"的那一瞬没在真装配上抓到**：mockllm 下没有让一枚孩子稳定比根慢的仪器（前一腿 §未测 4 的推理本腿认同）。
   本腿只多拿到一枚副产品：M1 那发打出来的包里孩子的行是 `"done":false`（**中途瞬时的字节**确实存在），
   但那不是"孩子比根慢"的可控读数。
8. **空白 task id 那一路两侧同答，本腿没有专门钉。** `streamkey.Subagent` 现在是两侧唯一入口（M4 钉住"只有一个铸造点"），
   但"空白 id ⇒ 两侧都返空串"没有自己的常驻用例（前一腿 §⑤2 记的那枚不一致，是靠**收进同一枚函数**关掉的，不是靠新判据）。

---

## ⑥ 对派单的不服（逐条；能给读数就给读数）

1. **未跟踪件是 4 枚，不是派单说的 1 枚**（派单 §0／台账 `A414`）。盘上除 `cmd/wisp/subagent_carrier_197_test.go` 之外还有
   `internal/panel/subagent_roster_197.go`（229 行）、`internal/panel/subagent_roster_197_test.go`（477 行）、
   `internal/streamkey/streamkey.go`（44 行），mtime 全在 21:11–21:26 那一窗（前一腿在飞时段）。
   ⇒ **这三枚是承重的**：`composer.go:91` 引用的 `TaskRosterSection`、`pump.go:430` 引用的 `streamkey.SubagentPrefix` 都在里面。
   **照派单 §1.3 那张 pathspec 清单（"那 6 枚＋新测试件"）提交，得到的是一棵编不过的树。**
   本腿按 10 枚逐名点名提交（`335b8d2b`：1748 插入／33 删除）。台账那句"372 行产码全在工作树未提交"只数了 6 枚已跟踪文件的 diff。
2. **前一腿的 `cmd/wisp` 不是"读数还没回来"——回来的那一发是红的**（派单 §0 括号里那句）。盘上就有：
   `.scratch/wisp/probes/197/r3/final-v-four-pkgs.txt`（21:43–21:44）里
   `FAIL github.com/CarlosShao/wisp/cmd/wisp 104.387s` ＋ `--- FAIL: TestRunPacketCarriesTheSubagentItsRosterRowFed`，
   红句逐字见 §①。同目录 `cmd-wisp-carrier-run1.txt`（21:42）里同一枚用例是 `--- PASS (1.72s)`——**同一棵树、两枚相反读数**，
   差别只在整包 vs `-run` 单跑。前一腿把后者当成了终态。⇒ 这条不是"没来得及写判据"，是**判据本身写错**（§② 末）。
3. **派单 §2(a) 只点了 `TestApprovalCardViewJSONKeysMatchFrontendTypes` 一枚，实际是两枚**：
   新增顶层键同样打红 `TestComposerContractTypesMatchFrontend`（两枚都从真 marshal 的字节里数键）。
   本腿终态逐字＝`approval_test.go:129: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`。
   两枚都在起手那 4 枚在册红里 ⇒ "名册未变"成立，但**红因变了**，界面那支要补的不止一张表（键集见 §④.3 末）。
4. **"入向那一跳今天到底要不要"——不要。** 派单 §1 让本腿"觉得必须有入向方法就停手具名报"。本腿现读的结论是**不需要**：
   名册行自带 `streamKey`，同一枚包的 `results` 节里 `correlationId === streamKey` 那一条就是那一页的正文——
   M1 那发打出来的整包正是"两半都在、只有名册那半没来"的字节，反向证明了这一点。
   ⚠ 顺手更正 `internal/panel/subagent_roster_197.go:36` 那句"would be a **fifth composer method**"：
   本腿现读 `grep -rn "PanelBridge" internal/ cmd/`（非测试码）**只有两枚注释命中**
   （`internal/panel/doc.go:2/:9`、`internal/agent/approval/gate.go:595`），**今天没有 Go 侧 `PanelBridge` 接口**，
   所以"第几枚方法"这个数**无法在盘上核**；准确的写法是"那是 C17 白名单里新添一枚入向方法名＝契约变更"。
   界面那支若坚持要宿主记住"当前格装的是哪枚会话"，才需要那一枚——**这一跳具名交回编排者落 `A##`**，本腿没有加、也没有改名。
5. **派单 §0 那句"internal/panel 红＝恰好起手那 4 枚"本腿复跑复核＝真**（终态仍 4 枚、逐名同一批）。
   但同一句里的 `internal/tools ok 16.471s` 本腿复跑是 `ok 17.240s`（起手）／终态 `PASS=158 FAIL=0`——
   秒数不该被当读数引用，**枚数才是**（这也是为什么 §④ 全用 `-v` 点数）。
6. **`pump_test.go:123`／`:291` 派单给的行号本腿现读为**准**（逐字仍是那串 `"composer,generatedAt,pending,results"`），
   本腿没动它们。但派单 §2(a)② 那条"要无条件顶层键就得破这两枚钉"的取舍，实际落点是**第三条路**：
   键仍是 `omitempty`（不破那两枚钉），而"有读者必发"由**新加的常驻判据** `TestAPumpWithoutARosterReaderSendsFourKeys` 钉住。
   派单没写这条落点，前一腿代码里写了——本腿认同并把它写进 §②(a)。
7. **本腿自己的一枚小账**：`335b8d2b` 的 commit message 末尾多了一串字面 `$''`（本腿用 `printf` 传 message 时打错，
   禁 `--amend` 所以留着）。产码与文件集不受影响（`git show --stat` 10 枚文件、1748/33）。**这是本腿的错，不是派单的。**
8. **前一腿那 43 行 `internal/tools` 测试件里有一枚没人调用的 helper**（`subagent_197_test.go:266 streamClosed`）——
   本腿把它变成活的判据（第三枚 commit，§③ 末），因为它原本是那一腿"没来得及写完"的那半截。

---

## ⑦ 入库件清单（本腿，逐枚点名）

| commit | 内容 |
|---|---|
| `335b8d2b` | 10 枚文件：`cmd/wisp/panel_pump.go`、`cmd/wisp/run.go`、`cmd/wisp/subagent_carrier_197_test.go`、`internal/panel/composer.go`、`internal/panel/pump.go`、`internal/panel/subagent_roster_197.go`、`internal/panel/subagent_roster_197_test.go`、`internal/streamkey/streamkey.go`、`internal/tools/subagent_197.go`、`internal/tools/subagent_197_test.go`（1748 插入／33 删除） |
| 第二枚＝**`7e3e7663`**（判据补牙） | `internal/tools/subagent_197_test.go`（把 `streamClosed` 变成活的判据，3/6 行、**删除列＝0**）＋`cmd/wisp/subagent_carrier_197_test.go`（把 `tasks` 那节的原始字节打进 `-v` 读数） |
| 第三枚＝本枚（提交时自指 `git log -1` 可取） | 本文件＋`.scratch/wisp/probes/197/r3b/**`（起手读数、门禁逐包、台件与四发变异逐发输出、`pre/` 里变异前快照的 5 枚源件） |

**只 commit、未 push**；写面之外零触碰（`docs/evidence/s1/152-*.md`、`.gitignore`、`.scratch/wisp/probes/152/**`、
`.scratch/wisp/probes/161/r6/logs/**`、`design/**` 那 16 枚删除、`.scratch/wisp/.scratch/**` 一枚都没进本腿的 commit）。

## ④ 门禁逐包读数 ＋ 在册红名册逐名比对 ＋ 格式门

### ④.1 起手那一发（本腿亲手跑，`.scratch/wisp/probes/197/r3b/start-gate-3pkgs.txt`）

`internal/tools` **`ok 17.240s`**／`internal/panel` **`FAIL 2.508s`**（红 4 枚＝在册那 4 枚）／
`cmd/wisp` **`FAIL 92.418s`**（红 1 枚＝前一腿产码自带的跨瞬断言，逐字见 §①）。

### ④.2 终态逐包读数（`-v` 现数，同一棵树上跑两遍：第一遍 `final-*`＝`335b8d2b`，第二遍 `final2-*`＝补牙之后）

| 门禁 | 读数（本腿亲手跑） |
|---|---|
| `go test -count=1 -v ./internal/panel/` | `rc=1` **PASS=95 FAIL=4 SKIP=0**（4 枚＝在册；95＝前一腿 21:12 那发的 90 ＋ 本票载体层新增 5 枚） |
| `go test -count=1 -v ./internal/tools/` | `rc=0` **PASS=158 FAIL=0 SKIP=0**（前一腿那 43 行是**改既有用例**、没加新用例，所以枚数照旧） |
| `go test -count=1 -v ./internal/agent/` | `rc=0` **PASS=73 FAIL=0 SKIP=0**（池的契约拷贝在这包，照 r1b 的门禁口径一起跑） |
| `go test -count=1 -v ./internal/streamkey/` | `rc=0` **PASS=0 FAIL=0**＝`[no test files]`（一枚常量＋一枚函数，判据在两包与铸造点扫描里） |
| `go test -count=1 -v ./cmd/wisp/`（整包，**不是 `-run`**） | `rc=0` **PASS=96 FAIL=0 SKIP=0**（96＝起手 93 ＋ 本票载体层 3 枚；⚠ 长跑约 90–105 秒） |
| `go build ./...` | **rc=0**，输出空 |
| `go vet ./internal/panel/ ./internal/tools/ ./internal/agent/ ./internal/streamkey/ ./cmd/wisp/` | **rc=0**，输出空 |
| `sh scripts/d22scan.sh` | **rc=0 clean**：`no D22 ban violations`；ban#1-5 `internal/`=216、`cmd/`=24 枚产码；ban#7 `internal/tools/`=22；ban#8 `internal/`=**453**、`cmd/`=**50**（注释与 `_test.go` 在射程内，比 r1b 那发的 450／49 各多，多的就是本腿入库的源件与测试件）、`design/`=39、`frontend/`=85 |
| gofumpt（`"$GOPATH/bin/gofumpt.exe"` v0.12.0） | `-l internal cmd` **输出为空＝0 枚未净**。⚠ 前一腿交下来的两枚未跟踪测试件（`internal/panel/subagent_roster_197_test.go`、`cmd/wisp/subagent_carrier_197_test.go`）**起手是脏的**，本腿跑过 `-w` 之后才为 0 |

**补牙之后的第二遍复跑（终态＝本文件入库的那一棵树，`final2-summary.txt`）**：
`internal/tools rc=0 **PASS=158 FAIL=0 SKIP=0**`、`cmd/wisp rc=0 **PASS=96 FAIL=0 SKIP=0**`。
`internal/panel`／`internal/agent`／`internal/streamkey` 第二遍**没重跑**——第二枚 commit 只动两枚测试件
（`internal/tools/subagent_197_test.go`、`cmd/wisp/subagent_carrier_197_test.go`），碰不到那三包；
若要三包也按终态数，那是验收程的活，本腿在此具名而不是替它报数。
两枚逐名读数在第二遍里照旧：`packet tasks section: rows=2 poolCap=4 child=d0d84534-… runStatus=Thinking afterStatus=Settling key=subagent:d0d84534-… bytes=2789`、
`bound crossed: rows=34 streams=34 truncated=true elided=0 dropped=[]`。

### ④.3 在册红名册逐名比对（起手 4 枚 vs 终态 4 枚，`internal/panel`）

| 起手在册红（派单给名） | 终态现量（`grep '^--- FAIL' final-v-panel.txt`） | 结论 |
|---|---|---|
| `TestApprovalCardViewJSONKeysMatchFrontendTypes` | `--- FAIL` 在 | **名册未变**（本腿一未修、二未 Skip、三未放宽） |
| `TestComposerContractTypesMatchFrontend` | `--- FAIL` 在 | **名册未变** |
| `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` | `--- FAIL` 在 | **名册未变** |
| `TestC21DesignTokensFourwayAgree` | `--- FAIL` 在 | **名册未变** |

**没有第 5 枚、也没有任何一枚消失**；`cmd/wisp` 与 `internal/tools` 终态各 **0** 枚红。
两枚名册钉的红因**加长了一枚**（这就是派单 §2(a) 说的"这是它该做的工"）：

```
approval_test.go:129: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare
```

`[instructions]` → `[instructions tasks]`：**`tasks` 是载体层新加的那一枚**，界面侧 `PanelSnapshot` 还没声明它。
要交给界面那支的键集（逐字段，本腿只出 Go 侧现读，`frontend/**` 零写面、内容不转述）：

- `PanelSnapshot` 新键 `tasks`，可缺（Go 侧 `omitempty` 指针）：`tasks?: TaskRosterSection`
- `TaskRosterSection` ＝ `rows: TaskRowView[]`、`inFlightSlots: number`、`poolCap: number`、
  `streamTruncated: boolean`、`streamElidedRunes: number`、`droppedStreamKeys: string[]`（**空数组＝没有丢，不是 null**）
- `TaskRowView` ＝ `taskId: string`、`label: string`、`kind: string`（宿主词表：`root`／`subagent`）、
  `parentTaskId: string`（**空＝根行**）、`status: string`（D43 那 20 枚名字之一，`statusKnown` 为假时**恒空**）、
  `statusKnown: boolean`、`statusReason?: string`（缺省键；空状态时的那句原因，宿主原话带过来）、
  `streamKey: string`（**点进去那一页的地址**：同一枚包 `results` 节里 `correlationId === streamKey` 那一条就是正文）、
  `blockedOnApproval: boolean`、`streamTruncated: boolean`、`streamElidedRunes: number`、`streamDropped: boolean`
- `ResultChunk` 那三枚键（`correlationId`/`text`/`done`）**一枚没加**——"这一行是谁的"由名册行的 `streamKey` 回答，
  所以 `approval_test.go:134` 那句 `ResultChunk <-> ResultChunkView: 3 JSON keys reconciled` 终态照旧。

---
