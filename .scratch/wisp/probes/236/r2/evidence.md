# 236-r2 — AC#1 ＋ AC#1b 的交件（五节）

腿＝`236-r2`；射程＝票 236 的 **AC#1（两格）**与 **AC#1b（三格）**，共五枚拒绝分支，全在 `internal/tools`。
⛔ 本件不答 AC#2／AC#3／AC#4／AC#5／AC#6，判语里也不顺带裁它们。

---

## §0 起手锚

| 尺 | 读数 |
|---|---|
| 撤票口令 `sed -n '3,5p'` 扫 `WITHDRAWN`／`撤`／`作废` | **未触发停手**，尺逐字：`sed -n '3,5p' …236-*.md \| grep -c -e WITHDRAWN -e '撤' -e '作废'` ＝ **`0`（rc=1，零行命中）**。三行内容＝`**Status**：**待派**（编排者 09-29 20:2x 立，起手锚点 c5d88a7f…）`／`来源＝…台账 A451。`／`本票的射程不是"功能没做"…〔契约邻接，要另批〕`。⚠ 口径写明：`撤`／`作废`／`票 269 就地撤` 这些字**在本票别处确实有**（`:54` 作废的是"哪一步红"那一句、`:130` 撤的是**另一枚票 269**），但派单指定的射程是**第 3–5 行**，那一窗零命中 ⇒ 不停手 |
| `git log -1 --format='%h %ad %s' --date=format:'%m-%d %H:%M'` | `3d9b8374 10-06 13:12 probes(232-r2 收尾代提): 写腿撞 150 轮帽但产码与证件全在盘上，编排者只代提未代判` |
| `git status --porcelain -- internal cmd` | **0 行**（空） |
| `date '+%m-%d %H:%M'` | `10-06 13:15` |
| `grep -n '^- \[ \]' .scratch/wisp/issues/236-*.md` | **7 枚**：`:26` AC#1／`:27` AC#2／`:28` AC#3／`:29` AC#4／`:30` AC#5／`:31` AC#6／`:85` AC#1b。**七枚全不归本腿**（派单 §5），本腿不碰框、不改 `Status:`、不改名。 |

> 撤票口令的逐字辨读补尺（我自己跑的，免得只留一句解释）：见 §1 表末行。

---

## §1 现量（票面点名的"零命中"尺，每一把我自己复跑）

尺的口径统一是 **`git --no-pager grep -c "<字面量>" HEAD -- '*_test.go'`**（＝HEAD 提交态的测试文件，不含本机工作树的瞬时脏样，也不含别人在飞的副本）。

| # | 票面断言 | 我复跑的尺 | 原始读数 | 判语 |
|---|---|---|---|---|
| 1 | AC#1：`taskCancel.Execute` 两支 fail-closed 摘掉后**各 0 枚红**（`221-v1` 的 m4／m5，票面"现量"表第 1 行标 ⛔ 未复跑） | 见 §3 两发突变 | §3 逐发记 | 本腿**今天实测**这两发的红名册，票面第 1 行的"未复跑"欠账在此结清 |
| 2 | AC#1 `:707` `caller == ""` 支的字面量在 `_test.go` 零命中 | `任务名册未接线` 见下一行；`不知道是谁要停`／`这条调用没有宿主给的任务 id` | **rc=1 零命中**（两支字面量各自） | 票面断言成立（此支由 §3 m-1a／m-1b 补牙） |
| 3 | AC#1 `:699` `Roster == nil` 支 | `git --no-pager grep -c '任务名册未接线' HEAD -- '*_test.go'` | **rc=1 零命中** | 成立。⚠ 同名册句在产码里 **5 处**（`subagent_197.go:254`／`task.go:515`／`task.go:700`／`task.go:828` 注释／`task_backfill.go:108`），**五处共用同一个前缀 `任务名册未接线`** ⇒ 只断子串 `未接线` 的既有尺（`task_output_leg_test.go:136`）钉的是 **task.output 那一支**，不是 cancel 这一支 |
| 4 | AC#1b `:253-255`（`Roster == nil`） | 同 #3（`任务名册未接线`）＋`拒绝派生子代理` | **两把都 rc=1 零命中** | 成立 ⇒ **需要补钉** |
| 5 | AC#1b `:256-258`（`BaseOptions/ParentTools == nil`） | `宿主没有给出派生用的装配`／`BaseOptions/ParentTools 未接线` | **两把都 rc=1 零命中** | 成立 ⇒ **需要补钉** |
| 6 | AC#1b `:259-263`（`parentID == ""`） | `这条调用没有宿主给的任务 id`／`不知道父任务是谁` | **两把都 rc=1 零命中** | 成立 ⇒ **需要补钉** |
| 7 | 票面 §六（编排者增量）口径「两条字面量 `任务名册未接线`／`没有宿主给的任务 id` 同样 `_test.go` 零命中」 | `git --no-pager grep -c '没有宿主给的任务 id' HEAD -- '*_test.go'` | 见 §5（这把尺我要分清"完整字面量"与"票面那句缩写"，两条读数都记） | — |
| 8 | 基线绿名册 | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/tools/ -count=1 -v` | **rc=0；`^--- PASS`＝198／`^--- FAIL`＝0／`^--- SKIP`＝0；`ok github.com/CarlosShao/wisp/internal/tools 15.403s`**；日志＝`logs/baseline-verbose.txt`（1101 行） | **本包起手全绿**，零枚历史在册红。⇒ §3 的红名册里任何一枚都是**我这一发造成的**，没有"别人的红"可混 |
| 9 | 行号漂移核查（票面 §九："census 里所有 `file:line` 按最新 HEAD 重跑"） | `sed -n` 逐行读 `task.go:687-734`／`subagent_197.go:238-282` | `task.go:699`＝`if t.d.Roster == nil {`、`:707`＝`if caller == "" {`；`subagent_197.go:253`／`:256`／`:259` ＝ 三支的 `if` 行 | **票面五处行号在 HEAD 上逐一对上，本腿无漂移**（不必按"偏 1"修正） |

---

## §2 两格逐格表

| 格 | 分支位置 | 打算答什么（载具） | 命令 | 原始读数 | 判语 |
|---|---|---|---|---|---|
| AC#1 `caller == ""` | `internal/tools/task.go:707-709` | **载具＝票面 §73 现成的那两个**：`build221(t, windowGate221(), provider, false)` ＋ `x.cancel(t, "", target)`；target 用**真派生的孩子**（`x.h.spawn` → `x.h.childRow(t).TaskID`），不是编的 id。钉**完整字面量**。⚠ 为什么必须钉字面量而不是 `IsError`：摘支后调用掉进 `rec.ParentTaskID != caller` 那一支（`:730`），**照样 `IsError==true`** ⇒ 只断 IsError 无牙 | §3 m-1a | §3 | 见 §3 m-1a |
| AC#1 `Roster == nil` | `internal/tools/task.go:699-701` | **载具＝票面 §72 现成的那个**：`mustRegisterTaskEntries(t, TaskDeps{Roster: nil})`（注册的正是 `BuiltinTaskEntries` 全家，含 `task.cancel`）＋ `New(Options{… Gate: windowGate221()})`。⛔ **不能沿用 `task_output_leg_test.go:126` 那份 gate**：它写的是 `Gate: NoGate{}`，而 `task.cancel` 冻在 **L1**（`task.go:662`），NoGate 的 `PendingWindow` 答 `AnswerReject` ⇒ 调用**进不了 Execute**（＝假绿，正是 `task_cancel_221_legs_test.go:13-19` 头部逐字警告过的形状）。本腿**不改**那枚既有函数一字 | §3 m-1b | §3 | 见 §3 m-1b |
| AC#1b `Roster == nil` | `internal/tools/subagent_197.go:253-255` | 载具＝`newSub197Harness` ＋ `buildWith(…, mutate: d.Roster = nil)`（**mutate 钩子已在仓**，`subagent_197_test.go:175`，同形先例 `:720` 用它钉 `Provenance = nil`）⇒ 零新 infra。`task.spawn` 冻在 **L0**（`subagent_197.go:222`）⇒ L0 pass，gate 不参与，NoGate 在这里**不会**造成假绿（判据注释里写清这个不对称） | §3 m-1c | §3 | 见 §3 |
| AC#1b `BaseOptions/ParentTools == nil` | `subagent_197.go:256-258` | 同一 `buildWith` mutate 钩子，各 nil 一枚字段（`\|\|` 的两半都钉） | §3 m-1d | §3 | 见 §3 |
| AC#1b `parentID == ""` | `subagent_197.go:259-263` | ⚠ `spawnResult` 把 `TaskID/CorrelationID` 写死成 `parent197`（`subagent_197_test.go:217`）⇒ **真桥不自造 id**：本腿加一枚**包内未导出**helper，直接 `bridge.Execute` 递 `TaskID:""` ＋ `CorrelationID:""`；`bridge.go:269-271` 只做 `CorrelationID==""→TaskID` 回落，两枚皆空 ⇒ 空（票面 §73 的同一句话） | §3 m-1e | §3 | 见 §3 |

**形状约束（五格共用）**：⛔ 零产码改动、零导出名新增、零放宽既有断言。写面＝一枚新建的包内测试文件 `internal/tools/failclosed_236_teeth_test.go`（helper 全是未导出）。

---

## §3 突变名册（每发：摘哪一支 ＋ overlay 路径 ＋ 哪一枚用例红 ＋ 红句原文 ＋ 还原后绿）

> 载体＝`go test -overlay`（⛔ 不在共享工作树里种针）。合成副本落 `D:/tmp/wisp236r2/`，每发的替换 json 记在下面。
> ⚠ 仪器事实（票面 AC#2 那句，写进判据注释）：**读盘型尺（`os.ReadFile`）对 `-overlay` 结构性不可见**——本腿五枚钉**全是行为尺**（经真桥真发一次调用、断回执文本），没有一枚读盘，所以 overlay 对本腿五支**都量得到**。这一条与 `task_cancel_221_legs_test.go:224` 那把词面尺的射程不是一回事，本腿不替它担保。

（逐发读数见下，m-1a … m-1e。）

---

## §4 门禁读数

（待本腿写满后逐把实跑填入。）

---

## §5 判不动的地方

（本腿自己填，编排者不代填。）
