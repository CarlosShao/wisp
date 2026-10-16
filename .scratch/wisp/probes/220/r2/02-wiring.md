# 票 220 r2 — 接线读数（改后）与快照真色

- 腿名 `220-r2`｜落笔时刻 `2026-10-16 11:04 +08`（`date` 现跑）
- 接线 commit：`2d10bd0f`（`cmd/wisp/panel_pump.go` + `cmd/wisp/run.go`）
- 起手 HEAD（本腿第 0 步）：`e396301d`；改前基线见同目录 `01-baseline.md`。

## 1 生产调用点尺（两发，同一把尺）

尺＝`grep -rn 'L1Windows' cmd/wisp/*.go`（射程＝`cmd/wisp` 顶层 133 枚 `.go`；工作树；**量的**）。

| 时刻 | 命中 | 生产调用点枚数 |
|---|---|---|
| 改前（HEAD `e396301d` 与工作树一致） | 零行、`git grep -c` rc=1 | **0 枚** |
| 改后（接线 commit `2d10bd0f`） | `cmd/wisp/panel_pump.go:108`（reader 本体 `rt.gate.LiveL1Windows()`）<br>`cmd/wisp/run.go:731`（装配处 `L1Windows: rt.liveL1WindowWaits,`）<br>`cmd/wisp/panel_l1_blocked_220_test.go:5,141,173`（测试件，不算生产） | **2 枚**（两枚文件各 1） |

沿到消费端（`internal/panel/pump.go`，本腿一字未动）：`:200` 字段 → `:323-327` 读它 →
`taskRosterSectionFrom(p.src.Tasks(), cards, waits)`（`:327`）→ 那一格的 join 在
`internal/panel/subagent_roster_197.go`（220-r1 的 corr＋taskID 双查）。

⇒ 派单 (a)① 成立：改前 0 枚 ⇄ 改后 2 枚，且这一条链在现树逐枚行号可查（本腿落笔时现读）。
⛔ 拿 `go build`／`go vet` 的 rc=0 当"接上了"——本腿两处 rc 也都记在 `05-gates.md`，
但它们只证明编译面，不证明这一格；证明这一格的是上面那枚数与下面那发颜色。

## 2 快照真色（同一枚用例的两发颜色，逐字）

用例＝`cmd/wisp/panel_l1_blocked_220_test.go` 的
`TestRunPacketLightsTheRosterRowOfItsLiveL1Window`（真腿：装配那一形，⛔ fixture 直喂；
窗口经 `rt.bridge.Execute` 的真 L1 路开出来，`rt.gate.LiveL1Windows()` 现读为活）。

### 2.1 改前那一发（红，逐字；日志＝`prechange-test.log`）

```
    panel_l1_blocked_220_test.go:215: PRE-CHANGE / WIRING READING: an L1 window was live on task "e87919dd-aef8-4c9a-aba2-25b304202a09" (corr "l1-220-corr-1792119597202068100") and its roster row still read blockedOnApproval=false - the pump's reader is only as good as the producer the composition root points at it (cmd/wisp/run.go's PumpSources literal)
--- FAIL: TestRunPacketLightsTheRosterRowOfItsLiveL1Window (13.01s)
```

⇒ 这一发红**不是**"应该会红"：窗口真开出来了（同一次跑里 `setupErr` 为空、包外那几发
`pre blocked=false / post blocked=false / pending 0 / queue 0` 全绿、写盘真落了），
缺的只有生产者。这一发就是票面 AC#2 要的"改前读数当场跑出来"。

### 2.2 改后那一发（绿，逐字；日志＝`postchange-test.log`）

```
    panel_l1_blocked_220_test.go:249: L1 window task=d2d9866b-5840-4201-a413-3aee41b2df22 corr=l1-220-corr-1792119892843526500 | pre blocked=false | lit at publish #2 (was #1; pending 0, queue 0) | post blocked=false | lit bytes: {...}
--- PASS: TestRunPacketLightsTheRosterRowOfItsLiveL1Window (12.87s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	13.080s
```

那一发 `lit bytes` 里与本判据有关的两枚字节（逐字摘自日志，⛔ 转述）：

```
"pending":[],
"tasks":{"rows":[{"taskId":"d2d9866b-5840-4201-a413-3aee41b2df22", ... "blockedOnApproval":true, ...
```

⇒ 真色三件同时成立：① 窗口活着时那一行为 `true`（wire 上的字节，不是 Go 侧解码）；
② `pending` 仍是空数组 ⇒ 窗口没被塞进可答复段（220-r1 拒绝的第二扇门没开，票 146 的
`LiveApprovals` 判据射程不变，本腿同发读回 `queue 0`）；
③ 同一枚任务的行在窗口结束后回到 `false` ⇒ 这一维随状态动，不是常量。
④ `lit at publish #2 (was #1)`：本探针只在改前那一格自己发布过一次（`#1`），
点亮那一发是 `#2` ⇒ 点亮那包字节是**这枚 run 自己的发布**（`consoleApprovalUI.Prompt` 里
`u.publish()`，接线后在 run.go:1451／`PendingWindow` 在窗口内 Prompt，gate.go:283-294），
不是测试自己喂的一包。

## 3 本腿在这一格上⛔ 不主张的东西

- ⛔ 界面已经看得见：传输那一跳仍属票 33／35 的地界，本腿止于"页面能读到的那包含这一格"。
- ⛔ 名册渲染语义有任何改动（corr＋taskID 双查是 `480b970d2` 的格子），
  ⛔ 把 `internal/panel/subagent_blocked_220_test.go:151-152` 那枚 fixture 读成生产可达。
- ⛔ 新导出方法名／新 `panel.*` 名（reader 是非导出的 `liveL1WindowWaits`，
  `L1WindowWait` 用的是 220-r1 已落的那一枚，C17 白名单一字未动）。
- 点亮的是**附加的一维**：那包字节的 `status` 仍是宿主填的（本例 root 行的诚实空＋理由），
  `blockedOnApproval:true` 没有把它吃掉（票 220「外部对照」那条约束）。
