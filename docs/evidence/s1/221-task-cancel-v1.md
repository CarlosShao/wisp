# 221 — 非实现者对抗验收表（腿 `221-v1`）

- **锚点**：本腿自取 `git rev-parse HEAD` = **`29215a93696a2127b4558c2417554c6f04d14af6`**（短 `29215a93`，与派单给的线索一致）。
- **被验对象**：`a7993a9b`（实现腿起手锚点）..`29215a93`，三发＝`453eebab`（草案＋预检，现核＝**零产码**，只有票面 31 行＋probe plan 82 行）／`c44b30c4`（甲形＋乙形同一发）／`29215a93`（AC#5 读数记录）。
- **实现腿**：`221-r1`。它把票面五枚 AC 全部自行勾成 `[x]` ⇒ **本表把这五枚勾一律读成〔实现方自述〕**，逐格重裁。
  ⚠ 这件事本身是一处治理缺陷（`AGENTS.md` §0.3「裁决者≠实现者」），记在最后那节。
- **批准边界**：台账 `A434`（`docs/reports/pending-and-issues.md:9317-9331`），七条写死；撤销口令「撤 221 甲」。
- **可迁移性核过**：`git diff c44b30c4..HEAD -- internal cmd` ＝**空**；`internal/tools/task.go` 与 `subagent_197.go` 工作树 md5 ＝ HEAD md5
  （`4138177e29ffff427776a73eb0d61a9b` / `06caf8f5465ff1c47c310da27fe3ffab`）⇒ 归档读数可直接用。
- **本腿起跑前置**：每次 `go test` 前 `tasklist //FI "IMAGENAME eq go.exe" | grep -c go.exe` ＝ **0**；突变全部走 `-overlay`，
  收尾 `git status --porcelain -- internal cmd` ＝**空**（13 发全核，见下表末行）。
- **判语口径**：成立／不成立／附条件成立，每格带本腿自己跑出来的读数。

## 判语汇总（九格，正文在各格小节）

| # | 格子 | 判语（本腿现跑裁出） |
|---|---|---|
| AC#1 | 说明书与现实一致（常驻能力尺） | **成立**（m1 摘注册那行 ⇒ 这把尺响；分母＝四家构造器的并集，不是单枚 `BuiltinTaskEntries`；扫能力不扫词面 ⇒ 注释污染不成立） |
| AC#2 | 生产调用者 ≥1 且名册那行不消失（`A434` ③④⑤） | **成立**（非测试调用者现量 **1** 枚＝`task.go:740`，12 枚命中里零注释命中；③借名 `Muted` 未新造、m12 有牙；④名册**根本没有删行 API**；⑤只进现成审计 sink、上屏字段零新增） |
| AC#3 | 两枚正控＋父取消不级联 | **成立**（m2/m3 各**恰 1 枚红**＝拒绝来自权限那一步不是门；m6 证实"被拒不动行"不是装饰；m10 三枚红＝非级联双钉；m7 反向证实 NoGate 拦的是门）⚠ **附一条缺格**：两支 fail-closed（`caller==""`／`Roster==nil`）m4/m5 **各 0 枚红**＝今天无尺 |
| AC#4 | DEFERRED 登记（只裁本票那枚标记） | **成立**（摘标记与接线同在 `c44b30c4` 一发，`git show` 两条计数各＝1；`task.list` 那行是 diff 的**上下文行**＝逐字未动）。尺本腿用"overlay 换读路径＋正控"测出：**(a) 支有牙**（m14 放回标记 ⇒ FAIL），**(b) 支无牙**（m13 摘掉 `task.list` 标记 ⇒ PASS，叙述行 `:279` 把计数撑住）＝下一格要修的形状缺陷，不是本票的账 |
| AC#5 | 整包终态读数 | **成立**（本腿自发 rc=1：23 ok／3 FAIL／5 无测试；名级红 6 枚＝在册 panel 4＋ball 1＋`cmd/wisp` 一枚计时红（隔离复量 3/3 绿）⇒ **本票新造红 0 枚**）⚠ 更正：在册名册不是稳定集合，`risk/R-116-1` 与 `cmd/wisp/TestTicket223*` 会换位 |
| 乙形 | 三处删句 | **成立**（三处逐条现读；m8/m9 各恰 1 枚红；既有钉 `Test197CancelIsPerRowAndNeverCascades` asis 绿＝"停掉父任务不会级联"没被顺手删） |
| 额外-A | `run.go` 是否只改注释 | **成立＝自称可复核为真**（非注释删除 **0** 枚，非注释新增也 **0** 枚） |
| 额外-B | 票 175 那枚普查尺 | **成立（不是削弱）**（−7 全是同一枚 map 键的 gofumpt 重排、值文本逐字未变，零断言被删；m11 删新答话 ⇒ 普查尺恰 1 枚红＝载荷非装饰） |
| 前提核验 | 可见但拒的三条前提 | **(i) 真**（孩子目录里**有** `task.cancel`，过滤器只摘 `task.spawn` 一枚）；**(ii) 真且挡驾无钉**（拒绝句回显 `ParentTaskID`＋存在/不存在两形，靠 uuidv4 不可猜，但仓里无尺钉"任务 id 必须宿主铸造"，熵源失败那支还退化成时钟派生） |

## 名册（本腿现跑，`-run 'Test221|Test197|TestEveryRegisteredToolIsClassifiedForMarking175r2' -v`）

改前（asis，`.scratch/wisp/probes/221/v1/logs/mut-asis.txt`）：**rc=0／两层结果行 20 枚／`--- FAIL` 0 枚**
＝8 枚 `Test221*` ＋ 11 枚 `Test197*` ＋ 1 枚票 175 普查尺，全绿。

| 突变（`-overlay`，工作树零写入） | 拆掉的那一形 | 红名册（本腿读数） |
|---|---|---|
| `m1-cancel-unregistered` | `BuiltinTaskEntries` 里 `task.cancel` 那行注册 | **6 枚红**，含 `Test221EveryPromisedTaskNameIsRegistered`（AC#1 的尺有牙）＋ `…TaskCancelIsRegisteredAtItsFrozenLevel` ＋ `…ParentStopsItsOwnChildRowAndStreamSettle` ＋ `…SubagentCannotStopSiblingOrItself` ＋ `…ParentCancellationStillDoesNotCascade` ＋ `…UnderNoGateStopsAtTheWindow` |
| `m2-no-parent-check` | `rec.ParentTaskID != caller` 那一支拒绝（改成 `if false &&`） | **恰好 1 枚红**＝`Test221SubagentCannotStopSiblingOrItself` |
| `m3-no-selfstop-check` | `target == caller` 自停那一支 | **恰好 1 枚红**＝同一枚 AC#3 用例 |
| `m4-no-empty-caller-guard` | `caller == ""` 那枚 fail-closed | **0 枚红（rc=0）** ⇒ 那一支没牙（见 AC#3 判语的缺格） |
| `m5-no-nil-roster-guard` | `t.d.Roster == nil` 那枚 fail-closed | **0 枚红（rc=0）** ⇒ 同上（实现腿自己把 M7 留在桌上，本腿证实） |
| `m6-cancel-before-refusal` | 把 `Roster.Cancel(target)` 注到拒绝判定**之前** | **1 枚红**＝`Test221SubagentCannotStopSiblingOrItself` ⇒ "被拒的调用没动行"那句断言不是装饰 |
| `m7-level-l0` | `Declared: risk.L1` 改成 `risk.L0` | **2 枚红**＝`…TaskCancelIsRegisteredAtItsFrozenLevel` ＋ `…UnderNoGateStopsAtTheWindow` ⇒ 级别有钉，且 NoGate 那枚拦下的确实是**门** |
| `m8-old-promise-returned` | 乙形①②的旧裸许诺放回 `Description()` | **1 枚红**＝`Test221SpawnDescriptionPromisesOnlyWhatIsTrue` |
| `m9-old-receipt-returned` | 乙形③旧回执放回 `task.spawn` 不等了那段 | **1 枚红**＝`Test221ParentCancellationStillDoesNotCascade` |
| `m10-cascade-returns` | `context.WithCancel(context.WithoutCancel(ctx))` → `WithCancel(ctx)` | **3 枚红**＝`Test221ParentCancellationStillDoesNotCascade` ＋ 既有 `Test197CancelIsPerRowAndNeverCascades` ＋ `Test197FullPoolRefusesNextSpawnWithReadableReason` ⇒ 双钉 |
| `m12-finalize-wrong-state` | `finalize` 把停掉那一行写成 `Settled` 而不是 `Stopped`(＝借来的 `Muted`) | **1 枚红**＝`Test221ParentStopsItsOwnChildRowAndStreamSettle` |
| `m11-census-answer-deleted` | 删掉票 175 普查表里新加的 `task.cancel` 那行 | **1 枚红**＝`TestEveryRegisteredToolIsClassifiedForMarking175r2` ⇒ 那行是**答问**不是放宽 |

---

## AC#1 说明书与现实一致（常驻能力尺）

**判据原句（票面「收 `221-c1`」改写后的口径，逐字）**：
> 尺的分母＝"这次装配真正注册出去的所有工具名的并集"（现跑＝`BuiltinTaskEntries`＋`BuiltinSubagentEntries`＋其余 `Builtin*Entries` 的构造结果，落点＝包内测试枚举这几家构造器），分子＝各 `Description()`／`Parameters()` 文本里出现的 `task.` 词根。⇒ 今天这一发应恰好红一枚＝`task.cancel`。⚠ 不许把分母写成单枚 `BuiltinTaskEntries`。

**我跑了什么**
- 尺＝`internal/tools/task_cancel_221_test.go:89 Test221EveryPromisedTaskNameIsRegistered`；现跑绿：`--- PASS`（`.scratch/.../logs/t221-baseline.txt`，rc=0）。
- 分母来源＝`allBuiltinEntriesHere()`（`:42`）逐家调 4 枚构造器；本腿现跑 `grep -rn '^func Builtin.*Entries' internal/tools/*.go | grep -v _test` ＝
  **恰好 4 枚**（`fs.go:325`／`fs_write.go:772`／`subagent_197.go:233`／`task.go:600`）⇒ 并集把四家全枚举了，**不是**票面④纠正掉的单枚 `BuiltinTaskEntries` 那一形。
- 有牙证明（`m1`）：摘掉注册那行 ⇒ 这把尺**响**，红句逐字（本腿读数）＝
  「说明书对模型许诺了一枚不存在的工具：task.spawn 的 Description() 写着「task.cancel 只有派生它的那枚父任务能用它单独停孩子——子代理停兄弟、停自己都一律被拒」，而 "task.cancel" 没有注册进这次装配的并集（并集 9 枚：fs.delete / fs.edit / fs.list / fs.move / fs.read / fs.trash / fs.write / task.output / task.spawn）」
  ⇒ **指出是哪一句、哪个槽位、缺哪一枚、并集几枚**，四件齐；且没把 `task.spawn` 报成假红。
- 词面 vs 能力：分子取 `e.Tool.Description()`／`e.Tool.Parameters()` 的**运行时返回值**，分母取**构造器结果**⇒ 注释、日志、票面文字都进不了这把尺（与 AC#2 那把 grep 的 self-pollutable 问题是两回事）。`taskRootRe`＝`task\.[a-z][a-z0-9_]*`，`task_id` 这类参数名不会被当词根。
- 尺自扫射程：`len(roots)==0` 直接 `t.Fatal`（`:130`）⇒ 分子空转不会读成绿。

**判语**：**成立**。
⚠ 两处射程注（不算缺陷，记给下一位）：① 分母是**手列四家构造器**，将来第 5 枚 `Builtin*Entries` 出现时这把尺不会自己知道（Go 里也没法枚举包内函数，只能是纪律）；
② 它只扫内置工具，插件目录里的说明文字不在射程内（与 D34 那张冻表的射程一致，本票不欠）。

---

## AC#2 生产调用者 ≥1 且名册那一行不消失（对 `A434` ③④⑤）

**判据原句（票面）**：
> `TaskRoster.Cancel` 有**生产调用者**（现跑点数：改前 0 枚 ⇒ 改后 ≥1 枚，且调用者在装配路径上、不在测试里）；停一名孩子后**它自己的流要有终态**（不许让那行凭空消失，要留下"被谁停的"），名册那一行状态要落到 D43 已有的名（**不许新造态名**）。

**我跑了什么**
- 点数（本腿现跑，排除注释的写法）：`grep -rn '\.Cancel(' --include='*.go' internal/ cmd/ | grep -v _test` ＝ 12 枚命中，逐枚读后**名册那枚只有 1 处**＝`internal/tools/task.go:740`（`taskCancel.Execute` 内），其余 11 处全是 `context.CancelFunc`／`RunningTask.Cancel`／`replyRoot`／`reloadRoot`。
  ⇒ **改前 0 ⇒ 改后 1**，实现腿自称的数对了。装配路径核：`cmd/wisp/run.go:429` 的 `for _, e := range tools.BuiltinTaskEntries(...)` 注册。
  ⚠ 本腿另核了"注释冒充调用者"这一形：那 12 枚命中里**没有一行是注释**；`subagent_197.go:39/:197`、`task.go:30/:684` 写的是 `TaskRoster.Cancel`（无左点无括号），进不了这把尺——实现腿把 run.go 注释改成不含调用语法的写法这件事**成立**。
- `A434` ③（借名、不新造态名）：`subagent_197.go:109` `subagentStateStopped = statemachine.StateMuted` 是**改前就在**的借名（`:100` 注释自己写着「D43 没有 "cancelled" 这个名字」）；
  `git diff a7993a9b..HEAD -- internal/tools/subagent_197.go` 里那一族常量**一行未动**⇒ 本票没造新名。用例侧断言：`task_cancel_221_legs_test.go:325-331` 既比 `final.State == subagentStateStopped` 又跑 `statemachine.Valid(final.State)`。
  有牙证明＝`m12`：让 `finalize` 把停掉那一行写成 `Settled` ⇒ 唯一红的正是 `Test221ParentStopsItsOwnChildRowAndStreamSettle`。
  ⚠ 派单让我"摘 `Roster.Cancel` 里改那一行的分支"这一攻法**在代码里没有对象**：`Cancel`（`task.go:447-464`）只做"查 id→查行在不在→查句柄→`cancel()`→`return true, ""`"，**它从不写那一行**；那一行的终态由派生侧 `finalize`（`subagent_197.go:413-433`）落。本腿因此改攻 `finalize`（＝m12），并把这一处口径差异具名记在这里。
- `A434` ④（那一行不许消失）：本腿读了名册全部方法（`grep -n 'func (r \*TaskRoster)' internal/tools/task.go` ＝ Record / WatchRow / Look / Count / Descendants / PublishSubagent / MarkRoot / TryAcquireSubagentSlot / InFlightSubagents / RunningSubagentIDs / AttachCancel / DetachCancel / Cancel）
  ⇒ **名册根本没有"删行"这枚 API**，"消失"这一形在结构上就不可达；用例还逐枚比了 `ParentTaskID`／`Kind`／`Label`（`:333-336`），并另钉了一枚隔壁行不受牵连（`:356-363`）。
- `A434` ⑤（"被谁停"只进审计、不加上屏字段）：
  产码侧 `Cancel` **不写谁停的**（返回 `(true, "")`），`taskCancel.Execute` 也**没有**新写审计——审计来自桥里那两条**现成 sink**：`bridge.go:992` 的 `tools: call kind=%s task=%s corr=%s tool=%s …`（`book()`，每发调用都记）与 `bridge.go:1030` 落库的 `memory.ToolCall{TaskID, Tool, ArgsJSON, …}`。
  生产接线本腿现读：`cmd/wisp/run.go:584` `Logf: rt.auditf`、`:571` `Journal: mem` ⇒ 这两条 sink 在真路径上是接上的，不是只活测试。
  用例把这一维钉住（`:369-388`）：审计文本须含 `tool=task.cancel` 与 `task=<调用者>`，落库那行须 `Tool=="task.cancel"` 且 `ArgsJSON` 含目标 id 且 `TaskID==调用者`。
  上屏字段零新增：`git diff a7993a9b..HEAD --stat` 里**没有** `internal/panel/**`、没有 `TaskOutput` 结构体改动 ⇒ 没有借这枚票加契约面字段。
- **判这格时要说清的边界**：审计记的是"谁调了 `task.cancel`、参数是哪枚 id"，**不是**"谁停了谁"这一语义事件；被窗口／被权限拒的调用同样留这一行（`book()` 对所有结局都记）。
  按 `A434` ⑤ 把票面 AC#2 那句"要留下'被谁停的'"**降级成〔审计可查〕**的原文口径，这个形状是**合规的**；实现腿还把自己先前写在回执里那句"这一停记在宿主审计里"**收回**了（票面 log 续 3 末段），收窄成当下为真的两句话——这一步是**减谎**，不是缺账。

**判语**：**成立**（③④⑤三条各自有读数；⑤ 的实现方式是"复用现成 sink＋用例证明"，不是新造载体）。

---

## AC#3 两枚正控＋父取消不级联

**判据原句（票面）**：
> 造两枚假腿——① 子代理停**兄弟**② 子代理停**自己**——**两枚都必须被拒**；并一条判据钉住"父取消仍不级联"（`subagent_197.go:327` 那形不许被本票改坏）。

**我跑了什么**（先证实派单点名的假绿机制，再攻两枚正控）
- 机制自证：`NoGate.PendingWindow` 回 `AnswerReject`（`internal/tools/gate.go`），而 `task.cancel` 声明 L1 ⇒ 走 `NoGate` 的桥**到不了 `Execute`**。
  本腿用 `m7`（L1→L0）反向证实：改级别后 `Test221TaskCancelUnderNoGateStopsAtTheWindow` **变红**（本腿读数 2 枚红）⇒ 那枚用例观察到的确实是**门**给的拒绝，而不是权限给的。
  两枚正控用的桥是 `windowGate221()`（把 L1 窗口如实答复成 `AnswerTimeout`＝到点即执行，与 `bridge.go` 的 L1 分支同向），所以拒绝只能来自 `Execute` 里那几支判定。
- 正控①（停兄弟）：`m2` 摘掉 `rec.ParentTaskID != caller` 那一支 ⇒ **恰好 1 枚红**＝`Test221SubagentCannotStopSiblingOrItself`。
  ⇒ 那枚"拒绝"**确实**是权限那一步给的（派单 (a) 要的答案＝成立）。
- 正控②（停自己）：`m3` 摘掉 `target == caller` 那一支 ⇒ **恰好 1 枚红**＝同一枚用例 ⇒ 有牙（实现腿 log 里那条"M3 第一次是假的、断言只写了『含自己』"的自翻，本腿复核其收窄后的断言：现在只认「不许停自己」那一支的文案，`if false` 之后没有第二枚文案能替代它）。
- (e) 被拒的调用没动任何行：`m6` 把 `Roster.Cancel(target)` 注到拒绝判定之前 ⇒ **1 枚红**＝同一枚用例（断言在 `:426-431`，读的是 `state` 仍为 `Thinking`）⇒ 那条断言不是装饰。
- (d) 父取消不级联：`m10` 把 `context.WithoutCancel(ctx)` 摘掉 ⇒ **3 枚红**（本票新钉 ＋ 既有 `Test197CancelIsPerRowAndNeverCascades` ＋ `Test197FullPoolRefusesNextSpawnWithReadableReason`）⇒ 双钉，本票改坏不了。
  "没有被级联取消"那枚字面量在乙形删句后**逐字还在**：`subagent_197.go:397`（本腿现读）与回执用例 `:490`；`Test221ParentCancellationStillDoesNotCascade` asis 绿。
- 两枚正控都用**真桥＋宿主 minted 的调用者身份**（`x.cancel` 走 `bridge.Execute`，`CorrelationID` 由宿主盖章；`caller` 取 `CorrelationID(ctx)`，不是参数），本腿读了 `taskCancel.Execute` 全支：参数里没有"父任务 id"这一维可填 ⇒ 模型无法靠改参数冒充别人家的父。

**判语**：**成立**，但**附一条缺格**（不是本票的账，写在下面以便翻勾）：
两枚 fail-closed 支（`caller == ""`、`t.d.Roster == nil`）**今天没有任何用例能响**——`m4`/`m5` 各摘一支，**rc=0、0 枚红**（本腿读数）。
实现腿自己在票面 log 里把 M7 留在桌上没自批，本腿证实它说的属实且**范围比 M7 更宽一枚**（`caller == ""` 那一支同样无尺）。
⇒ AC#3 的**两枚正控＋非级联**这三件票面明写的事成立；两支 fail-closed 无钉是本票之后该补的一格（建议见最后一节）。

---

## AC#4 DEFERRED 登记（按改写后口径，只裁本票那一枚标记）

**判据原句（票面「收 `221-c1`」改写后）**：
> 本票**不承诺**"与 `SPEC-12 §5` 双向 1:1 对得上"。改成本票自己这枚标记的两条：(a) 若走甲形，`internal/tools/task.go:20-27` 头部把 `task.cancel` 标为 DEFERRED 的那段要摘掉；(b) 摘标记与接线同一次核销。1:1 那枚仪器本身归票 225。

**我跑了什么**
- (a) 现读 `internal/tools/task.go:18-32` 头部：DEFERRED 那一行只剩 `task.list`；`task.cancel` 改写成 `L1 … (this file, 221)`。
  尺＝`Test221DeferredMarkerForCancelLiftedButListStillMarked`（读包内 `task.go` 源文本）asis **绿**。
- (b) 同一次核销（本腿自己跑 `git show`）：`git show c44b30c4 -- internal/tools/task.go | grep -c '^-.*task.cancel.*DEFERRED'` ＝ **1**（摘标记），
  `… | grep -c '^+.*{Tool: taskCancel{d: d}'` ＝ **1**（接线），**同一发**；前一发 `453eebab` 现读**零产码**（只有票面＋probe plan）。
  ⇒ 没有"先摘标记再补实现"的中间态那一发。
- `task.list` 那枚标记一字不动：`git diff a7993a9b..HEAD -- internal/tools/task.go | grep -c 'task.list.*DEFERRED'` ＝ **1**，且那一行在 diff 里是**上下文行**（非 `+`/`-`）⇒ 逐字未改。
- 没有顺手动别的登记（本腿现跑）：`git diff a7993a9b..HEAD --name-only -- internal cmd` ＝ **6 枚**（`run.go`／`subagent_197.go`／`task.go`＋两枚新用例＋票 175 那枚普查件），`git diff --name-only -- docs/specs docs/PLAN.md` ＝ **0 行**；
  `grep -n 'DEFERRED' internal/tools/task.go` 现量 **3 行**（`:23` 是 `task.list` 那枚标记，`:32`／`:279` 是叙述），**没有一行含 `task.cancel`**；
  `DEFERRED(D-xx)` 那枚与 `SPEC-12 §5` 对账用的标记形态在本票 6 枚写面上**全部为 0 命中**（＝本票没动过它，对账仍归票 225）。
- ⚠ **本腿一把尺的牙测到了意外结果，具名记在这里**：`Test221DeferredMarkerForCancelLiftedButListStillMarked` 的实现是 `os.ReadFile("task.go")`（`task_cancel_221_legs_test.go:224`）
  ⇒ **直接对 `task.go` 做 `-overlay` 突变对这把尺不可见**（overlay 只改编译期视图，不改运行时读盘）：m13/m14 两发直接跑都是 rc=0／0 枚红，**那两发不作证据**。
  本腿因此换一种攻法（仍零写工作树）：overlay 只替换**测试文件里的读路径**，让那把尺的判定逻辑逐字跑在突变体文本上，并配一枚正控。三发读数（`logs/teeth-*.txt`）：
  | 发 | 尺读到的是 | 读数 |
  |---|---|---|
  | `teeth-control` | `git cat-file blob HEAD:internal/tools/task.go` 的原样副本 | **PASS**（⇒ 这套装法本身是通的，不是空转） |
  | `teeth-m14-cancel-marked-again` | 把 `task.cancel` 的 DEFERRED 行**放回**而注册仍在 | **FAIL**，红句逐字＝「task.go 里还有 1 行把 task.cancel 标成 DEFERRED：它已经注册并有实现，『在册＋无实现＋无人认领』正是票 164 AC#1 要杀的那一形」（`:241`）⇒ **AC#4(a) 那支有牙** |
  | `teeth-m13-list-marker-lifted` | **摘掉** `task.list` 那行 DEFERRED 标记而不接线 | **PASS（不响）** ⇒ **"task.list 的标记不许被摘"这一支今天没有牙** |
  m13 不响的原因本腿现读定位：那把尺数的是"**含 `DEFERRED` 的行里含 `task.list`**"，而 `task.go` 里除了标记行还有**一行叙述**（HEAD `:279`「…task.list is DEFERRED (§7 :1531)…」）同样命中；
  摘掉标记行后计数＝1（本腿对拉：HEAD 2 行 → m13 1 行）⇒ `listMarked == 0` 那一支永远等不到。
  ⚠ 这支的不对称要说清：`cancelMarked != 0` 是**越有越响**（任何 DEFERRED＋task.cancel 的行都算，包括叙述行，偏严但保守＝安全方向），
  `listMarked == 0` 是**越无越响**（靠叙述行就能被撑住＝偏松）。⇒ 下一格该补的是把 (b) 那支从"扫词面"改成"扫能力"（看 `BuiltinTaskEntries` 里到底有没有 `task.list`），不是再加一句词面。

**判语**：**成立**——(a)(b) 两句的**事实**由 `git show` 逐字核实（摘标记与接线同在 `c44b30c4`、`task.list` 那行是 diff 的上下文行＝逐字未动），
不依赖那把尺；而那把尺本腿测出：**(a) 支有牙（m14 响）、(b) 支无牙（m13 不响，原因如上）**。
⚠ 按派单口径重申：**不要**因为"和 `SPEC-12 §5` 对不上"判本票——那是票 225 的账（`grep -n 'task\.' docs/specs/SPEC-12-roadmap-governance.md` ＝ 0 命中这件事，本腿复核仍为 0）。

---

## AC#5 整包终态读数

**判据原句（票面）**：
> `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/...` 到终态＋逐名比红名集合（`-run` 单跑不算）；gofumpt 用 `"$GOPATH/bin/gofumpt.exe"`。

**我跑了什么**（本腿自己那一发；起跑前 `tasklist //FI "IMAGENAME eq go.exe" | grep -c go.exe` ＝ **0**，一次只跑一发）
- 命令＝`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 ./cmd/wisp ./internal/...`，日志 `.scratch/wisp/probes/221/v1/logs/ac5-full.txt`，**跑到终态 rc=1**。
- 包级（正确尺现跑）：`grep -cP '^ok[[:space:]]'` ＝ **23**／`grep -cP '^FAIL\t'` ＝ **3**（`cmd/wisp` 118.938s／`internal/ball` 0.241s／`internal/panel` 2.982s）／无测试文件 **5** ⇒ 31 枚包，与实现腿续 3 的自洽式（23+3+5）对得上。
- 名级红（两层同一把尺，`grep -cE '^[[:space:]]*--- FAIL'`）＝ **6 枚**：
  `cmd/wisp/TestTicket223RefusedLooseningKeepsOldValues`（5.15s，`config_reload_223_test.go:379`「the operator is not told the loosening did not take effect」）＋
  `internal/ball/TestC21TableColourRowsMatchTokensCSS` ＋ `internal/panel` 4 枚
  （`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`）。
- **在册常红 5 枚对上**：panel 4＋ball 1，失败正文本腿现读逐字＝`read design/assets/tokens.css: open …design\assets\tokens.css: The system cannot find the path specified`
  ⇒ 落在**别人在工作树里删掉的那棵 `design/**`** 上；三枚冻结件本腿零改动（`git status --porcelain -- internal/panel internal/perm` 空）。按派单口径**不报成新缺陷、一字不许动**。
- ⚠ **本发新出现的一枚（不在实现腿终态那发、也不在派单点名的在册名册里）＝`cmd/wisp/TestTicket223RefusedLooseningKeepsOldValues`**：
  归因做了——单包隔离复量 **3/3 全绿**（`5.055s`／`5.019s`／`5.034s`，日志 `recheck223-1..3.txt`）；这一族属票 223 的写面，本票三发对 `cmd/wisp` 只有注释（额外-A 的 0 枚非注释删除读数）
  ⇒ 判 **〔待复量〕已过：整包并发／本机 CI 同跑造成的计时红，不是本票造成的红**。
- ⚠ 与实现腿终态那发（续 3 第 3 发）的差异具名：那发是 `internal/risk/TestResolvePerCallBudget`（`R-116-1`）红而 `cmd/wisp` 绿，本发正好相反（`risk` 绿、`cmd/wisp` 红）
  ⇒ **"改前在册红＝一枚不多一枚不少的稳定名册"这句话本身不成立**，那两枚计时红（`risk/R-116-1`、`cmd/wisp/TestTicket223*`）会换位；
  实现腿"本票新增红 0 枚"的**结论仍成立**，但它的比名口径要补一半：逐名比必须配隔离复量。
- 本票唯一动过的包 `internal/tools` 本发 **`ok 17.435s`**。

**判语**：**成立**——整包到终态、红名集合逐名可归因、**本票三发新造的红＝零枚**；
附带一条给编排者的口径更正（不是退回理由）：在册红名册**不是稳定集合**，今后"逐名比红名集合"这句要自带"计时类红单包隔离复量"那一半。

---

## 乙形三处（`subagent_197.go`，与甲形同一发）

**判据原句（票面＋`A434` ⑥）**：
> 乙形那三处（`subagent_197.go:196`、`:197` 的尾巴、`:389`）与甲形同一发落地——⚠ 改的时候必须原样保住 `"停掉父任务不会级联"` 那句。

**我跑了什么**
- 三处逐条现读（行号是本腿 `git diff` 现量，不是票面那组已漂移的号）：
  ① `Description()` 那句裸许诺"可以用 `task.cancel` 单独停它"已换成带边界的说法；② 同句尾巴那个括号"（要停它得单独停）"已删；
  ③ `task.spawn` 不等了那段回执（`subagent_197.go:395-400`）由"可以单独停它：它的流键是…"换成"它的父任务 %s 还能用 task.cancel 单独停它（只有派生它的那一枚能停）；它自己的流键是 %s"。
- 有牙证明：`m8`（①②旧句还原）⇒ 红 `Test221SpawnDescriptionPromisesOnlyWhatIsTrue`；`m9`（③旧回执还原）⇒ 红 `Test221ParentCancellationStillDoesNotCascade`。两枚都是**恰 1 枚红**＝钉在该钉的地方。
- "停掉父任务不会级联"没被顺手删：既有钉子 `Test197CancelIsPerRowAndNeverCascades`（`subagent_197_test.go:796-799` 逐字断那枚字面量）asis **绿**（本腿读数）；新用例 `:263` 也把它列进 `keep` 清单。
- 换成"当下为真的话"这件事本腿逐句读了新文本：它写的是"只有派生它的那枚父任务能用它单独停孩子——子代理停兄弟、停自己都一律被拒"，
  与代码里的判定支一一对得上（AC#3 那几支）；并明写"由用户停某一枚子代理…今天没有落点"⇒ 与 `A434` ⑤ 的"不许算完成"同向，没有藏缺格。

**判语**：**成立**。

---

## 额外-A　`cmd/wisp/run.go` 的 +19/−10 是否每一行删除都是注释

**判据原句（派单）**：逐行看删除行，任何一条非注释删除都要具名报出来。

**我跑了什么**
- `git diff a7993a9b..HEAD -- cmd/wisp/run.go | grep '^-' | grep -v '^---' | grep -vcP '^-\s*//'` ＝ **0** ⇒ 10 枚删除行**全是注释**。
- 反向也核（派单没要求但更严）：`grep '^+' | grep -v '^+++' | grep -vP '^\+\s*//'` ＝ **空** ⇒ 19 枚新增行也全是注释，**没有一行产码**。
- 逐枚看了那 10 行内容：删的是"这一族只有 `task.output`"那一段旧说明（接完线就成了新的谎），换成语法化的两行族表＋生产落点说明。
- 顺带核了 AC#2 那把 grep 尺的污染问题：新注释写的是 `TaskRoster's Cancel had zero production callers` 与"count it with a grep for calls of that method outside _test"，
  **不含** `.Cancel(` 这个调用语法 ⇒ 本腿点尺那发读数（12 枚命中零注释命中）证实这句话现在是**真**的。

**判语**：**成立**（自称"只改注释"＝可复核为真）。

---

## 额外-B　`ticket175r2_stamp_live_test.go` 为什么被本票动 +12/−7

**判据原句（派单）**：如果它是为了迁就 221 而削弱了票 175 的某枚断言，直接判退回并写清被删断言的原句与它的存在理由。

**我跑了什么**
- 逐行读 diff：**−7／+7 是同一枚 map 的键被 gofumpt 重新对齐**（新增了一枚更长的键 `"task.cancel"`），**+5 是新增的 `task.cancel` 那一行＋四行说明注释**。
  **零枚断言被删**：这个 hunk 里没有 `if`/`t.Error`/`t.Fatal` 被删；被"删"的 7 行全是 `"fs.read": …` 这类表项的原样重排（值文本逐字未变，本腿对拉过）。
- 存在理由（本腿读 `:327-374` 的注释与断言原文）：票 175 AC#3 那把普查尺要求"本包注册的每一枚内置工具都必须在分类表里就地答一句"，
  未答者按 `risk.IsSensitiveSource` 分两支报 `t.Errorf`。⇒ **注册 `task.cancel` 必然要求这行答案**，不写它包就红。
- 有牙证明：`m11`（删掉那行答案）⇒ 红 `TestEveryRegisteredToolIsClassifiedForMarking175r2`（恰 1 枚）⇒ 那行是**载荷**，不是装饰。
- ⚠ 一句要说清的射程：这把普查尺只查**"有没有就地答"**＋答的方向与 `IsSensitiveSource` 是否自相矛盾，**不查答得对不对**。
  "unstamped（回执是本进程名册的事，没有外部正文回流）"这一判**本腿读了 `taskCancel.Execute` 的全部回执文本**（含目标 id、父 id、状态名，全是宿主自己记的东西，不含任何外部正文）⇒ 分类与现读一致，我认这句答。

**判语**：**成立**（不是削弱；是普查尺要求的就地答话）。

---

## 编排者已裁项（"可见但拒"）的三条前提核验

**我不重裁这一格，只验前提**（裁定＝保持可见但拒，理由三条见派单 §6）。

- (i) **孩子目录里到底有没有 `task.cancel`**：有。现读 `subagent_197.go:549` `const subagentHiddenTool = "task.spawn"` 与 `:552-567` 的目录过滤循环——
  过滤器**只按这一枚名字**摘（`if info.Name == subagentHiddenTool { continue }`），别家一律整份继承（`:293` `opt.Tools = newSubagentToolProvider(t.d.ParentTools)`）。
  ⇒ 派单那句"要过滤就等于新造一维可见性"的前提**成立**（今天的形状只有"摘 `task.spawn` 一枚"这一处特例，不是通用可见性机制）。
- (ii) **回执是否回显任意 target 的 `ParentTaskID` 与"存在/不存在"**：是，两处都是。
  ① `task.go` 归属拒绝支文本 `"拒绝停止：%s 是任务 %s 派生的孩子…"` 直接回显 `rec.ParentTaskID`；
  ② `known == false` 那支回答"查不到这个任务 %s"，与"存在但无权"是**两句不同的话** ⇒ 对任意 id 有一个 exists/not-exists 的oracle。
  今天不构成实际泄漏，靠的是**id 不可猜**这一条：`internal/agent/loop.go:1123 newTaskID()` 用 `crypto/rand` 铸 36 字符 uuidv4 形状，
  `Loop.Run`/`RunAsync`（`:321`/`:332`）**不接受**调用方给的 id ⇒ 模型侧无法自选一个短 id 去枚举。
  ⚠ **但这条挡驾在仓里没有保护它的尺**：`grep -rn newTaskID --include='*_test.go'` 只命中 `internal/agent/compress_trace_test.go:524-527`（那是**trace 属性形状**检查：`len==36 && dashes==4`）与 `:557`（同一次铸造不相等），
  **没有任何用例钉住"任务 id 必须由宿主铸造／`Loop` 不许开收外部 id 的入口"**。而 `internal/tools/task_backfill.go:38` 明写着 **call id 是 model-supplied**——
  也就是说"模型能自己命名一个 id"这件事在本仓**已经有一处是真的**（工具调用 id），只是任务 id 还没有。谁哪天把 task id 也接到模型给的值上，上面那枚 exists/oracle 就变成跨行名册内视。
  ⚠ 另一处形状：`newTaskID` 在 `rand.Read` 失败时**退化成 `time.Now().UnixNano()` 派生**（`loop.go:1126-1130`）——那种状态下 id 可猜，而这枚退化支同样没有尺。
  ⇒ **不判本票**（这不在 A434 的射程里，本票也没扩大它），按派单要求写成"下一格该有什么尺"，见最后一节第 1 条。

**判语**：三条前提中 (i)(ii) 均**核实为真**，其中 (ii) 的"不可猜"依赖的是一条**无钉**的现状。

---

## 门禁读数（本腿自己跑）

本腿自己跑的三件，都在终态树（HEAD `29215a93`）上：

- `"$GOPATH/bin/gofumpt.exe" -l internal/tools/ cmd/wisp/` ＝ **空**（rc=0）⇒ 本票那 6 枚写面没有格式债（裸 `gofumpt` 会 exit 127 出假读数，本腿按派单口径用了 `$GOPATH/bin` 那枚）。
- `go vet ./internal/tools/ ./cmd/wisp/` ＝ **rc=0**。
- `bash scripts/d22scan.sh` ＝ **rc=0**，尾行逐字（日志 `.scratch/wisp/probes/221/v1/logs/d22scan.txt`）＝
  「d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=219, bans #1-5 cmd/=29, ban #6 frontend/=85, ban #7 internal/tools/=22, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=462, ban #8 cmd/=63」。
  ⚠ 本腿只转述这枚**计数与结论**，没读 `frontend/**`／`design/**` 的任何内容（那两棵树在工作树里正被别人改动，见 AC#5 那 5 枚在册红的成因）。
- 突变卫生：13 发（12 枚 AC 攻法＋1 发 asis 对照）跑完后 `git status --porcelain -- internal cmd` ＝ **空**，逐次核过起跑前 `go.exe` 计数＝0，全程一次只跑一发。

---

## 没做完／留给编排者

1. **AC#5 与门禁已补完（本腿自发）**，剩余没做的只有下面几条。
   ⚠ 但有一条**要编排者亲手做**的收尾：整包那一发里出现过一枚**不在派单在册名册**的红（`cmd/wisp/TestTicket223RefusedLooseningKeepsOldValues`），
   本腿隔离复量 3/3 绿、判定为计时红，**但这枚名册今晚已经和实现腿那发不一致**（那发是 `risk/R-116-1` 红、本发是 `cmd/wisp` 红）
   ⇒ 在册常红名册本身需要有人重钉一次（本腿没有那枚名册的写面，也不该顺手改别人的红账）。
2. **AC#4 那把尺测出一支真的没牙，修法要另批**：`Test221DeferredMarkerForCancelLiftedButListStillMarked` 的 (a) 支（`task.cancel` 不许再被标 DEFERRED）**有牙**（`teeth-m14` FAIL），
   (b) 支（`task.list` 的标记不许被摘）**无牙**（`teeth-m13` 摘掉真标记行仍 PASS，因为 `task.go:279` 那行叙述同时含 `task.list` 与 `DEFERRED`，把 `listMarked` 撑到 1）。
   ⇒ 下一格＝把 (b) 从"扫词面"改成"扫能力"（看 `BuiltinTaskEntries` 里到底有没有 `task.list` 这枚注册），本腿零产码、没有动它。
   ⚠ 另记一条仪器事实，供今后台词面型尺的子程用：**读盘的尺（`os.ReadFile`）对 `-overlay` 突变结构性不可见**，
   要测这种尺的牙就得 overlay **替换测试文件本身的读路径**并配一枚正控（本腿 `logs/teeth-control.txt` 走的就是这条路）。
3. **两支 fail-closed 无尺（`caller == ""`、`t.d.Roster == nil`）**：`m4`/`m5` 读数都是 **0 枚红（rc=0）**。票面 AC#3 明写的那三件（两枚正控＋非级联）成立，
   但这两支"名字写在那儿、没有任何用例能响"的形状正是本票要杀的同一形（票 164 AC#1 族）。⇒ 建议下一格：一枚定向尺"未接线名册／无宿主 id 的调用必须在 `Execute` 那一步被拒且不动任何行"，
   并且**不要**把它塞回 221 的账里（本票不欠，但账要有人接）。
4. **跨行名册内视的钉子缺失**（见 §(ii)）：建议立一格＝"任务 id 只能由宿主铸造"这条不变式配一把尺——现仓里只有 trace 属性那枚形状检查（`compress_trace_test.go:524`），
   它挡不住"`Loop` 新增一枚收外部 id 的入口"。同时 `task_cancel.Execute` 的 exists/not-exists 两形（"查不到这个任务"vs"是别人家的孩子"）也值得一句"这组回答只对本调用者的子树以外使用"的裁定。
5. **`A434` ⑤ 的"被谁停"仍然是通用 per-call 审计，不是"停止"这一语义事件**：如果编排者认为审计里要能**单独查到"这一行被谁停过"**（而不是"谁调过 task.cancel 且参数是谁"），
   那是一个**新载体**，本腿判它不在 A434 已批范围内（要另批），本票没做也没声称做了——这一点实现腿写得比票面诚实（它把回执里那句承诺收回过）。
5. **治理缺陷（要记在编排者那本账上，不是代码账）**：`221-r1` 自行把票面 5 枚 AC 勾成 `[x]`，违反 `AGENTS.md` §0.3；本表把它全部读成〔实现方自述〕并逐格重裁。
   另有一处派单口径与实际形状不符，本腿就地更正而不是照抄：AC#2 那格派单给的攻法"摘 `Roster.Cancel` 里改那一行的分支"在代码里**没有对应物**（`Cancel` 不写那一行），本腿改攻 `finalize` 的落账支（m12）。
