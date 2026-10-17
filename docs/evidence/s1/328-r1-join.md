# 328-r1 落地件①：丙-带边界 的形状、四形实测、去牙成对两发、剩下的零读数

腿＝`328-r1`（落地腿）。起手锚＝**`648eda04`**（`git log -1` 现认：Sat Oct 17 12:34:32 2026 +0800，⛔ 抄派单里那句号）。
交件侧锚＝**`f705b082`**。写面＝`internal/panel/**`＋`internal/agent/**`（四枚提交逐枚 `--name-only` 见件②）。
裁形出处＝票 328 末尾 `编排者裁形＝丙-带边界` 那条 bullet；料＝`docs/evidence/s1/328-a1-census.md`（锚 `7535d83b`）。
本腿⛔ 选形、⛔ 翻任何 AC 框、⛔ 改票名／`Status:`／判据文字／加 `-done`；交的是料，不是判语。

## 一、落下来的形状（硬框逐枚对账）

名册侧（`internal/panel/subagent_roster_197.go`，产码一笔＝`b2d55342`）：

- `:276 rosteredTaskIDs(state)`（新增，非导出）＝成员只有"本次渲染自己那批行的 `TaskID`"，空号不入集。
- `:227 known := rosteredTaskIDs(state)` → 三处 `indexWaitingKey` 都带上它（卡片一枚 ⇄ L1 窗口两枚，不对称还在原处）。
- `:303 indexWaitingKey(set, known, key)` 第三臂，`ⓐ` 的判据本体：
  `:311 for task := range known` → `:312 if strings.HasPrefix(key, task+"#")` → `:313 set[task] = true`。
  两半都是判据：段必须**逐字节等于一枚在册任务号**，且原串里紧跟其后的是 `#`。⇒ 裁形的第一句话，⛔ 我另起一形。
- 第一臂（`set[key]=true`）与第二臂（`splitClashSuffix` 那枚数字改写）一字未动。
- `⛔` 新 JSON 键／`⛔` 新导出 `panel.*` 名／`⛔` 新 C17 方法名：本腿新增标识符只有 `rosteredTaskIDs`（非导出）
  与 `indexWaitingKey` 的一枚入参；没撞到需要停手上报的那格。尺＝交件侧 `git grep -n 'func [A-Z]' internal/panel/subagent_roster_197.go`。
- `ⓑ` `subagent_blocked_220_test.go:141`（`child-1#retry`）翻 **false→true**；`:142`（`child-1x`）原地不动；`:140`（`other-1#9`）仍 false。
- `⓭` 两处自陈改为说实话：`TaskRosterSectionFrom` 头补上 328 那一族（生产 corr 今天真落在名册上），
  `splitClashSuffix` 文档里那句 `a wider match than any code in this tree produces` 就地过期、改写为
  "变的不是这枚数字测试，而是 join 还接受哪一形"。数字臂的判据文字一字未动。
- `⓮` 见件②那一节：四面 diff 行数逐枚报 **0**。

`ⓐ` 正控那一发（`internal/agent/corr_roster_join_328_test.go`，一笔＝`f705b082`）：corr **不是本文件写的**——
真 golden SSE（`tool-then-text`）→ 真 C5 适配器 → `run`→`executeCalls`→`loop.go:676 callCorr(...)`，
corr 从 `h.echo().Calls()[0].Req.CorrelationID` 读出来，先四把尺钉住它就是 `<taskID>#<callID>` 那一形，
再走真搬运 `panel.NativeVerdict{...}.CardView()`（`pump.go:98`→`approval.go:72`，宿主 `liveVerdicts` 同一枚纯函数）
→ 真名册渲染（`TaskRosterSectionFrom`→`taskRosterSectionFrom`＝`pump.go:333` 调的那枚），断言 `:123` 那一行为 true。
★本票 ⛔ 的正是"测试里硬写 `task#3`／`#2` 那种撞名形"：这一发里没有一枚硬写的 corr。
跨包可达性（派单坑①）：`internal/panel` 非导出面只 import `projctx/risk/streamkey`（现量，见件②尺表），
`agent` 侧测试文件 import `panel` 无回环 ⇒ **没撞语言层、没开新导出名**。

## 二、必答①：四形逐枚点不亮＝实测读数

尺＝导出台树 `wisp-328r1-after-tree`（锚 `f705b082`）`go test -v ./internal/panel/ -run '…' -count=1 -timeout 30m`；
断言行＝`internal/panel/subagent_corr_task_prefix_328_test.go:51`（`assertBlocked(t, data, "child-1", c.want)`，读的是 `Publish()` 出来的线内 JSON）。
逐枚 `--- PASS` 子用例名（全绿那发＝`logs/after-panel-targeted.log`，RUN=19／PASS_TOP=5／FAIL_TOP=0／SKIP=0）：

| 形 | 用例（子测试名，逐字） | 期望 | 读数 |
|---|---|---|---|
| `T#call_e1`（支A，`loop.go:607`） | `arm A provider call id lights its task` | true | **亮**（PASS） |
| `T#call-0`（支B，`loop.go:605`） | `arm B id-less call lights its task` | true | **亮** |
| `T#call-0#7`（队列叠改，`queue.go:168`） | `arm B with the queue clash rewrite stacked lights its task` | true | **亮** |
| `T#9`（撞名形） | `the queue rewrite of a bare task id lights its task` | true | **亮** |
| 加测 `T#call_e1#7` | `arm A rewrite of a bare task id lights its task` | true | 亮 |
| `approval-<seq>`（空 corr 回退，`queue.go:165`） | `the empty-corr fallback names no row` | false | **仍不亮** |
| `host:*`（宿主发起，`gate.go:281`） | `a host-initiated correlation names no row` | false | **仍不亮** |
| 建在没人登记过的号上 | `a correlation built on a task nobody filed names no row` | false | 不亮（在册那半的牙） |

`T`＝该台名册自己登记的 `child-1`（夹具 `rows197()`，与票 220 那族同一枚台面）；每一枚同场还断言另一行 `root-1`＝false。
支B/叠改这两形的**字节出处**另有一发实测：`TestTheLoopMintsTheShapesTheRosterAccepts`（`:165`）把
`callCorr(T,"",0)`／`callCorr(T,"call_e1",0)` 的输出逐枚钉成 `T#call-0`／`T#call_e1` 并各自点亮自己那一行 ⇒ 名册那张表喂的字节＝铸造点产的字节，⛔ 按谓词文字推。

## 三、去牙：成对两发（落点不同、逐字红句、当场还原）

两发都打在**仓外导出树**的物理字节上，内容锚＝函数体内那枚 token（⛔ 绝对行号）；脚本＝`logs/mutate.sh`，逐字日志＝`logs/mut-M1.log`／`logs/mut-M2.log`。
两发 pristine 与还原后同一串：`sha256=450e1d1c47ee715fb74b56116a827988b0d9bfbc261fc06748c4704b4bd87028`，
`git hash-object=2da6c596369c1ee38672ae2656a7d5a3ecbaabfb`＝锚 `f705b082` 那枚 blob（`pristine_equals_blob_check=YES`，`restored_same_as_pristine=YES`）。

**发一（硬框 ⓒ 那一发）**：`strings.HasPrefix(key, task+"#")` → `strings.HasPrefix(key, task)`（抹掉分隔符那半边＝纯前缀）。
落点用例＝`:142` 那一枚，红（`go_test_rc=1`）：
`subagent_blocked_220_test.go:147: row "child-1" blockedOnApproval=true，期望 false； packets={"pending":[{"correlationId":"child-1x",…`（packet 全文逐字在日志里）
同场 `child-1#retry` 那枚仍 PASS ⇒ 这一发只对"分隔符"敏感，正是隔板该有的样子。

**发二（正控的牙）**：`set[task] = true` → `_ = task`（第三臂不再登记任何东西＝改回今天的形状）。
落点用例＝`ⓐ` 那一发，红（`go_test_agent_rc=1`）：
`corr_roster_join_328_test.go:124: the row for task "55b29c12-e4f2-47b1-9581-2c8e12e5b0ad" reads not-blocked although a card of its own call ("55b29c12-e4f2-47b1-9581-2c8e12e5b0ad#call_e1") is pending: this is ticket 328's defect, measured on a correlation the loop minted`
同发还红（`go_test_panel_rc=1`）：`:165` 两支各一句、`subagent_corr_task_prefix_328_test.go:51` 四形各一句（`child-1#call_e1`／`#call-0`／`#call-0#7`／`#call_e1#7`，句式 `row "child-1" blockedOnApproval=false，期望 true`），
以及 `subagent_blocked_220_test.go:147` 的 `child-1#retry`（⇒ `ⓑ` 那一枚翻向是**承重**的，不是装饰）。
★同一发里 `child-1#9`／`child-1#147` **仍绿** ⇒ 数字臂没被我碰、两臂各自承重，作差可分（⛔ 拿"红了一堆"当强度）。

## 四、`ⓑ` 那一枚翻向的账（逐枚具名，买真相）

- `subagent_blocked_220_test.go:141`：`{"a non-numeric suffix is not the queue's form", "child-1#retry", false}`
  → `{"a non-numeric tail on a filed task id lights (328)", "child-1#retry", true}`。
  为什么是买真相：`child-1#retry` 正是 `callCorr` 给一枚在册任务拼出来的形状（`<taskID>#<非数字尾>`），
  那一行本来就在等这张卡；旧期望 false 就是把票 328 的缺陷本体钉成了判据。用例**没删**、⛔ `t.Skip`、判据文字未动。
- `:142`（`child-1x`＝有前缀无分隔符）**原地不动**，且发一证明它是活的：抹掉分隔符 ⇒ 它红。
- `:140`（`other-1#9`）仍 false；`:143`（空 corr）仍 false。表头散文（`:119-131`，用例函数 `:132` 之前那段）改为点名"两族形＋分隔符那一半是牙"。

## 五、必答②：拆完之后还剩哪几格零读数（具名给尺）

★**照抄那一行，不美化**：屏幕上那一行今天**仍不会亮**——`frontend` 里 `blockedOnApproval` 命中文件数 **0**
（尺＝`git grep -l 'blockedOnApproval' f705b082 -- frontend | wc -l`，量的，交件侧现跑＝0；`-- frontend/src` 亦 0；
大小写不敏感同族 `git grep -il 'blockedonapproval'` 亦 0）。本票买的是"名册 join 得上了"那一跳，不是"人看见了"。
- 正控那一枚我要**更正票面/普查件的一个数**：同尺族 `correlationId` 在 `frontend` 命中 **11** 枚文件（`frontend/src` 内 **7** 枚），
  ⛔ 复现出票上写的"4 枚"；两枚数都列在这里、以我这把尺的 cmd 行为准（cmd 逐字在 `logs/end-rulers.log` 尺1）。
- 页面侧接收器自陈（现量行号，⛔ 抄普查件的 `:27-33`）：`cmd/wisp/panel_snapshot_deliver_329.go:28`
  `// window.wispOnPanelSnapshot does not exist yet …`＋`:57 const panelSnapshotGlobal = "window.wispOnPanelSnapshot"`。
- 生产包那句自陈（现量）：`cmd/wisp/subagent_carrier_197_test.go:598`
  `// 1. blockedOnApproval is never asserted on a production packet.` ⇒ 真 packet 上那一格仍零仪器（归票 329／前端会话）。
- `cmd/wisp/**` 整包本腿没跑（⛔ 派单：`228-r3b` 在取证，重包同一时刻只跑一个）⇒ 那包里的名册行
  （`subagent_blocked_197_test.go:190/:223`、`subagent_carrier_197_test.go:311` 三枚 false 断言）**没有改后读数**。
  静态界（量的）：尺 `git grep -nE '"[A-Za-z0-9_.-]+#[0-9]+"' f705b082 -- cmd/wisp`＝2 命中，逐枚看过都是注释里的
  `"AC#0"`／`"AC#3"`（`resident_audio_windows.go:15`、`slo_report_144_windows_test.go:873`），卡片 corr 里带 `#` 者＝**0 枚**
  ⇒ 第三臂要求 `#`，那一包不可能被我的放宽点亮新行；这一格是**推算的界**，不是颜色读数。
- `internal/tools/**` 里断言 `BlockedOnApproval` 的用例枚数＝**0**（`git grep -n 'BlockedOnApproval' f705b082 -- internal/tools | wc -l`，量的）⇒ 那包无本票读数需求。
- ★新找到的一格零读数（本腿现量，⛔ 抄件）：队列撞名改写 `<corr>#<seq>`（`queue.go:168`）**全仓没有任何用例真产过它**——
  尺 `git grep -nE '"[^"]*#[0-9]+"' f705b082 -- internal/agent/approval | wc -l`＝**0**；正控＝同尺族
  `git grep -nE '"[A-Za-z0-9_.-]+#[0-9]+"' f705b082 -- internal cmd | wc -l`＝**11 行**（逐枚看过：真值 7 枚＝
  `subagent_blocked_220_test.go` 的 `:108/:112/:138/:139/:140/:194` 与本腿新件的 `:42`；注释／假阳 4 枚＝
  同两枚文件的 `:15`／`:20`，加 `cmd/wisp/resident_audio_windows.go:15` 的 `"AC#0"` 与
  `cmd/wisp/slo_report_144_windows_test.go:873` 的 `"AC#3"`）。
  ⇒ `T#call-0#7` 那一形的"叠改"半边在本仓仍是**夹具形**，⛔ 经真队列。为什么没接上：`callCorr` 非导出且在 `internal/agent`，
  `Queue.push` 非导出且在 `internal/agent/approval`，两枚铸造点在不同包、都要新导出名才链得起来（⛔ 撞线即停手上报），
  所以本腿把这两半边分别实测（第二节两支的字节出处＋第三节发二），并把链不起来这件事具名登记在这里。
- `splitClashSuffix` 放宽后仍**没有任何测试直呼它**（量的，交件侧）：`git grep -n 'splitClashSuffix' f705b082 -- internal cmd`
  ＝ def `:333`＋doc `:318`＋唯一读者 `:308`，三命中；正控＝同尺族 `indexWaitingKey` 命中 8 行。
  ⇒ 普查件 §二(丙)(c) 那句"放宽后这一格没有任何仪器"对本腿**部分**成立：`child-1x`／`other-1#9`／四形表是**行为级**仪器，
  直呼这枚函数的仪器仍为零。
