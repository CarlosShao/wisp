# 220-v1 判语表（非实现者终裁·票 220 五格）

- 腿名 `220-v1`｜角色＝裁决者（可写测试件⛔ 动产码；本波依派单更正 (1)⛔ 新建 `cmd/wisp/**`、`internal/panel/**` 任何件）
- 起手 `git rev-parse HEAD`＝`ca0e30fce32425e31b9882c5d947e56a3b17d851`｜`date '+%F %T %z'`＝`2026-10-16 11:44:33 +0800`（起手现跑）
- 票框尺（现跑）＝`grep -c '^- \[ \]' <票面>`＝**5**，`grep -c '^- \[x\]'`＝**0**（rc=1）⇒ 本腿一枚框没翻
- 被裁实现方＝`220-r1`（`ce693ce4`＋`480b970d2`）＋`220-r2`（`2d10bd0f`＋`fc4f704b`，撞 150 轮帽死在最后一步）
- 逐枚 `git show --name-only` 写面核证（⛔ 区间尺）：
  `ce693ce4`＝`internal/agent/approval/{gate.go,ticket220_l1_window_read_test.go,window_read.go}`
  `480b970d2`＝`internal/panel/{pump.go,subagent_blocked_220_test.go,subagent_roster_197.go}`＋probes 14 枚
  `2d10bd0f`＝`cmd/wisp/{panel_pump.go,run.go}`
  `fc4f704b`＝`cmd/wisp/panel_l1_blocked_220_test.go`＋probes 3 枚
  ⇒ 四笔写面全在票 220 声明的射程内（`internal/agent/approval`＋`internal/panel`＋`cmd/wisp`＋probes），**零越界**

## 台面（本腿自己的，⛔ 引他人读数）

- 三枚仓外 `git archive` 导出树（七件套 `go.mod go.sum internal cmd tools scripts frontend`）：
  `bench`（HEAD，突变台）／`full`（HEAD，整包台）／`pre`（`fd01e7962`＝接线那笔的父）
  件数尺＝`find <tree> -name '*.go' | wc -l`＝**743**（HEAD 两枚）／**741**（pre）⇒ 非空
- 手工补的台面缺件（派单更正 (2)(3)＋放腿时刻那节）：`frontend/dist` 整目录（**4 枚文件**）、
  `build/*.dll`（3）、`third_party/sherpa-onnx/*.dll`（3）、另拷 3 枚 `.dll` 进 `cmd/wisp/` 旁；
  `full` 树再加 `docs`＋`design`（从 HEAD `git archive` 取，⛔ 动仓内）＝消掉两枚"读 docs/PLAN.md／design/assets/tokens.css 的台面缺件假红"
- 三枚树 `go build ./cmd/wisp/` 退码**全 0** 后才开始量（`logs/00-tree.txt`）
- `GOTMPDIR`＝`cygpath -w` 出来的**原生反斜杠形**且先 `mkdir -p`；跑 go 前置 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`
- ⛔ `-overlay`：五发突变全部打在导出树的物理字节上，每发先备份⇒跑⇒还原⇒`sha256` 对回（`logs/10-mut-check.txt`），
  还原后台面自证＝三枚文件树内 sha ⇄ `git cat-file HEAD:<f>` sha **逐枚 match=YES**（`logs/30-restore-verify.txt`）

## AC#1 join 读对（corr 与 taskID 双查⛔ 新增导出方法）＝**成立**

① 双查那两行的当前行号（内容锚，现跑 `git show HEAD:… | grep -n`）：
- `internal/panel/subagent_roster_197.go:214`＝`indexWaitingKey(waiting, card.CorrelationID)`（索引侧＝卡片的 corr，含去 `#十进制序号`）
- `internal/panel/subagent_roster_197.go:236`＝`BlockedOnApproval: row.TaskID != "" && waiting[row.TaskID],`（查询侧＝行的 taskID）
- L1 那一半的两枚索引＝`:219`（`w.CorrelationID`）／`:220`（`w.TaskID`）；去 `#序号` 的识别器＝`:262-270 indexWaitingKey`＋`:277-288 splitClashSuffix`

② 去牙反控（MU-A：把 `:214` 退成 220 之前的单键形 `waiting[card.CorrelationID] = true`，即不再认队列改写）⇒
`./internal/panel/ -count=1` 红名集合从基线 **8 枚**变 **11 枚**，新增三枚**必须红**＝
`TestASecondCardForTheSameTaskStillLightsItsRow`／`TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt`（红掉的子测试逐字＝`the queue's own rewrite lights` 与 `a multi-digit seq lights`）／`TestABlockedRowKeepsItsRunningDimension`。
逐字红句（`logs/20-panel-MU-A.txt`，⛔ 转述）：
`subagent_blocked_220_test.go:109: row "child-1" blockedOnApproval=false，期望 true； packets={"pending":[{"correlationId":"child-1#9",…`
`subagent_blocked_220_test.go:139: row "child-1" blockedOnApproval=false，期望 true； …`（同发两枚子测试：`child-1#9` 与 `child-1#147`）
`subagent_blocked_220_test.go:191: the second card is pending and the row does not say so`

③ 反形两发都跑了（**不是只跑一形**）：
- MU-B（L1 那一半只按 corr 索引，丢 `w.TaskID`）⇒ 红名 8→**10**，新增 `TestAL1WindowInWaitingLightsItsRow`（`:152` 逐字：`row "child-1" blockedOnApproval=false，期望 true`）＋`TestABlockedRowKeepsItsRunningDimension`（`:200` 逐字：`root row = {…BlockedOnApproval:false…}, want blocked AND still carrying its filed state`）
- MU-C（查询侧改成只按 corr 维＝`row.StreamKey`）⇒ 红名 8→**13**，除上面两族外另红 `TestTheRosterReaderPutsSubagentsOnTheWire`（`subagent_roster_197_test.go:165` 逐字：`the child has an L2 card on the queue and the packet does not say so: …`）
⇒ **两形都拦得住＝这把尺对这一格敏感**（⛔"两形都放行"那一判不成立）。基线 8 枚红的具名归因见文末〔台面/他票噪声〕。

④ 导出方法新增枚数尺（带分母，现跑）：
- 尺＝`git grep -h '^func (g \*Gate)' <ref> -- internal/agent/approval/{gate,report,window_read}.go | wc -l`（⛔ 测试件）：
  `ce693ce4^`＝**23** 枚 ⇄ `ce693ce4`＝**24** ⇄ `HEAD`＝**24** ⇒ 本票射程净增**1 枚**＝`Gate.LiveL1Windows()`（`window_read.go:63`）＝**AC#2 甲形那枚只读枚举口**，具名批准＝台账 `A562` §1（`pending-and-issues.md:11111` 起，本腿现读到）＋`A1002`（`:18680`）沿用
- 尺＝逐笔导出名差集（`git show <c>:<f>` 名册 diff）：`480b970d2` 在 `internal/panel/pump.go` 净增**1 枚导出类型**`L1WindowWait`、在 `subagent_roster_197.go` **0 枚**；`2d10bd0f` 两文件**0 枚**；`ce693ce4` 净增 `type L1Window`＋`Gate.LiveL1Windows`
⇒ **AC#1 那一格自身＝新增导出方法 0 枚**（`TaskRosterSectionFrom` 的对外签名逐字未动＝`:187-189` 仍转调内部 `taskRosterSectionFrom`）；本票全射程新增的导出面＝1 枚方法＋2 枚类型，全部归 AC#2 甲形与其数据形状，⛔ 记在 AC#1 名下。

## AC#2 让 L1 可见（落点已裁＝甲）＝**成立**

① "改前 0 枚生产者 ⇄ 改后 2 枚"＝**本腿自己现跑**，两发同一把尺、两个台面：
尺＝`grep -rn 'L1Windows' --include='*.go' <tree>/cmd | grep -v '_test.go' | grep -vE ':[0-9]+:\s*//' | wc -l`（射程＝该导出树 `cmd/`；**量的**）
- 改前锚树（`fd01e7962`）＝**0 枚**（同尺全部命中＝0，含测试与注释）
- HEAD 树（`ca0e30fce`）＝**2 枚**＝`cmd/wisp/panel_pump.go:108`（reader 本体 `wins := rt.gate.LiveL1Windows()`）＋`cmd/wisp/run.go:731`（装配处 `L1Windows: rt.liveL1WindowWaits,`）
⇒ 我自己数的计数＝**0 ⇄ 2**（与实现件那枚数偶然同值，但尺、树、时刻都是本腿现跑的；`logs/01-count.txt`＋`01b-count.txt`）

② 真腿用例确实**经真 run**（⛔ fixture 直喂）——`cmd/wisp/panel_l1_blocked_220_test.go` 现读＋我自己那一发：
开窗走 `rt.bridge.Execute(ctx, agent.ToolRequest{TaskID: rootID, CorrelationID: corr, Name: "fs.write", …})`（`:130-133`），
窗口活着与否是 `rt.gate.LiveL1Windows()` 现读（`:142` 轮询、`:174` 轮询到关），
根行 id 由 `waitRootRow197(rt)` 从这枚 run 的名册取（`:98`），点亮那一包字节只认 `rt.lastPanelSnapshot()`（`:159-163`）且带 `snapshotCount` 钉（`:220-223`）
我自己那一发＝**绿**：`--- PASS: TestRunPacketLightsTheRosterRowOfItsLiveL1Window (20.39s)`，`logs/04-cmdwisp-base.txt` 里那一行同屏读数＝
`L1 window task=06e1ed6f-… corr=l1-220-corr-1792122967550428000 | pre blocked=false | lit at publish #2 (was #1; pending 0, queue 0) | post blocked=false | lit bytes: {"pending":[],…"blockedOnApproval":true,…}`

③ 去牙＝摘掉 `run.go` 那一行供数（MU-D）⇒ 该用例**必须红且红在装配那一格**，实测红句逐字（`logs/21-cmdwisp-MU-D.txt`）：
`panel_l1_blocked_220_test.go:216: PRE-CHANGE / WIRING READING: an L1 window was live on task "ad5ba047-e1ec-4f8c-b944-403396ad94ae" (corr "l1-220-corr-1792123027139025600") and its roster row still read blockedOnApproval=false - the pump's reader is only as good as the producer the composition root points at it (cmd/wisp/run.go's PumpSources literal)`
⇒ 同一发里 `--- FAIL` 只有 1 枚、`setupErr` 为空（窗口真开出来）、`:249` 那行照打 `pre blocked=false | … | pending 0, queue 0 | post blocked=false`
⇒ **红在装配那一格，⛔ 红在建窗失败／台面故障**（`RUN=1 TOPPASS=0 FAIL=1 SKIP=0`，rc=1）

④ 反形＝供一枚常量非空列表（MU-E：`liveL1WindowWaits` 不看窗口、按名册行恒供 `[]L1WindowWait{{CorrelationID:"const-220v1", TaskID: 行 id}}`）⇒ **被抓住**，实测红句逐字（`logs/21-cmdwisp-MU-E.txt`）：
`panel_l1_blocked_220_test.go:212: the pre-window packet carried no row for "84a3350b-f974-4754-8e24-f0543fa7177c", or already lit it: the readings below would describe a constant, not a state`
⇒ 抓住它的判据字节＝**窗口还没开那一包**名册行的 `blockedOnApproval`（用例在 `:211-214` 就把这一发当"常量"拦下，`t.Fatalf`）
④′ "名册那行随窗口结束回到 `false`"这一发**本腿自己复跑过**＝② 那一发绿里的 `post blocked=false` ＋ 断言面 `:229-232`；
    且 `:211` 的 `pre found && !blocked` 与 `:215` 的 lit 双向都响 ⇒ 这一维随状态动，不是常量、也不是建包时的固定字面

## AC#3 两形各有正控＝**成立**（附一枚具名射程缺口，见"残余"）

- 形一〔L1 等待中〕：正控两枚——渲染级 `TestAL1WindowInWaitingLightsItsRow`（MU-B 掏掉 taskID 那一半即红）＋
  生产级 `TestRunPacketLightsTheRosterRowOfItsLiveL1Window`（MU-D 摘供数即红）。
  **坏世界真被造出来的同屏证据**在我自己那一发的行内：`pending 0, queue 0`＋`"pending":[]`＋`"blockedOnApproval":true`
  ＋同一发还带 `[确认 L1 fs.write]`／`ANSWER-EXPIRED` 两句界面自述（用例 `:240-245` 断的就是这两句），⛔ 把"没复现"读成"修法成立"那一支在此不成立
- 形二〔同任务第二张卡〕：正控＝`TestASecondCardForTheSameTaskStillLightsItsRow`＋`TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt`
  （六枚子测试含"非队列形不点亮"那一支：`a non-numeric suffix is not the queue's form`／`a task id that merely prefixes does not light`）；
  MU-A 一掏改写识别即三枚红⇒**种 X 必响**成立
- 缺口具名（⛔ 翻框、⛔ 我补）：形二的"坏世界"今天只在 fixture 层被造出来（手写串 `child-1#9`），
  全仓零枚用例经真队列把改写形跑出来（尺＝`grep -rn 'clash' internal/agent/approval/*_test.go`＝**0 命中**＋
  `grep -rnE '"[^"]*#[0-9]+"' --include='*_test.go' internal cmd` 只命中 `internal/panel/subagent_blocked_220_test.go` 那六枚手写样本）

## AC#4 ⛔ 把布尔改成态／⛔ 新造 D43 名＝**成立**

- `blockedOnApproval` 今天仍是布尔＝`internal/panel/subagent_roster_197.go:132`（`BlockedOnApproval bool` json 键 `blockedOnApproval`），
  且名册那一行未被替换：我自己那一发的字节里同一行同时带 `"status":""`/`"statusKnown":false`/`"statusReason":"宿主没有登记这一维…"` 与 `"blockedOnApproval":true`⇒附加的一维，不是替换（票面"外部对照"那条约束）
- 越界尺（现量）：本票射程（`ce693ce4^..HEAD`）内 `internal/statemachine/` 的 `git diff` **0 行**⇒状态名册与转移表零改动；
  名册里 `StateConfirming`/`StateAwaitingApproval` 的当前行号＝`internal/statemachine/states.go:21`/`:22`（现跑，与票面 `:16` 写的一致）
- 引 D43 那几枚编号的当前行号自证没动转移表（`internal/statemachine/table.go`，现跑 `grep -n 'D43: 17, …'`）：
  **#17**＝`:118`（→Confirming）与 `:123`（→AwaitingApproval）｜**#21**＝`:150`｜**#22**＝`:154`｜**#23**＝`:158`/`:163`｜**#24**＝`:168`/`:172`｜**#25**＝`:176`/`:180`
  （派单"四枚编号"若指回复那两枚 D43 名所出入的四条边，则＝`#17/#21/#22/#23`；本表把六行都给出来，判语不依赖是哪四枚）
- C17 白名单／三枚冻结件命中数：`git diff ce693ce4^..HEAD -- internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go`＝**0 行**；
  白名单载体件（`internal/panel/bridge.go`、`composer_dispatch.go`、`cmd/wisp/panel_inbound.go`）**不在本票四笔的写面文件名册里**（逐枚 `git show --name-only` 见抬头）；
  逐笔新增 `"panel.*.*"` 字面量＝**0 枚**（四笔全 0）；`internal/panel/**`＋`cmd/wisp/**` 内新增 `State*` 常量＝**0 枚**

## AC#5 整包终态读数＝**不可判**〔缺那一发"整包跑到终态"的 ./cmd/wisp 名册〕

- 我带的超时＝**显式 `-timeout 30m`**（./cmd/wisp 三发与 ./internal/... 两发全带上）。三发都没到超时——**没跑到终态的原因是测试进程被拦腰打断，⛔ 是 `panic: test timed out`**
- `./cmd/wisp -count=1 -v` 三发四数（同一枚仓外 `full` 树＝`ca0e30fce` 导出树，⛔ 共享工作树；分母取自该导出树）：
  发1（`logs/40-full-cmdwisp.txt`）＝`RUN=203 PASS_TOP=145 FAIL=3 SKIP=0 rc=1`（包级 161.777s）
  发2 同台面复跑（`logs/42-rerun-cmdwisp2.txt`）＝`RUN=203 PASS_TOP=145 FAIL=3 SKIP=0 rc=1`（168.634s，逐数同发1）
  发3 排除发1/2 拦腰那枚后（`logs/43-cmdwisp-skipAC13.txt`，`-skip '^TestAC13BringUpSurvivesAReusedThreadQuit$'`）＝`RUN=206 PASS_TOP=150 FAIL=1 SKIP=0 rc=1`（165.630s）
- 三发**全部未达终态**，断点逐名（原始 `-v` 日志末行都是 `=== RUN <名>` 紧跟 `FAIL github.com/CarlosShao/wisp/cmd/wisp …s`，零 `panic:` 句、零该用例自己的 `--- FAIL`＝进程硬终止）：
  发1／发2 断在 `TestAC13BringUpSurvivesAReusedThreadQuit`（`cmd/wisp/panel_resident_windows_test.go:523`，票 33 的 poisoned-thread 形）
  发3 断在 `TestAC14GoSideEvalPushReachesThePage`
  ⇒ 归因⛔ 本腿下（派单第二节那六枚的归因已由编排者成对台件闭掉；我这三发**到不了**那六枚里的五枚，具名在下面）
- 断点之前跑到的红名逐名归因（**全⛔ 本票射程、全⛔ 代码红**）：
  `TestAC2AC3DegradedLegsStillDeliverTheirVerdict/models` 与 `TestAC2ModelsLegBooksItsHandOffVerdictOnDisk`＝**台面缺件**，红句逐字 `leg_sink_nail_131_windows_test.go:494/349: committed models/manifest.json not found from the test working directory`
  ⇒ 这把尺自证：我在发3 前把 `models/`（HEAD 跟踪件，`git ls-files models`＝2）补进树，那两枚转绿，发3 的 `FAIL` 只剩 1 枚
  `TestCleanCheckoutBuilds_AC11`＝**导出树没有 `.git`**，红句逐字 `panel_host_gate_test.go:404: git rev-parse --show-toplevel: exit status 128`
- 那六枚已知红在我自己三发里的颜色（`logs/41-cmdwisp-sixred-verbatim.txt`）：
  `TestAC4FocusReturnToPriorWindowGap33r5` **三发全 PASS**（0.67s／1.48s／1.68s）⇒ 它不是恒定红，本机前台状态随时序变；
  其余五枚（`TestAC14GoSideEvalPushReachesThePage`／三枚 `Test258*`／`TestAC246DevLegIgnoresTheTestTaskInjection`）**未跑到**
  ⇒ 所以本腿⛔ 报"整包含 6 枚红"、⛔ 报"本波新增红"，只能报"三发未达终态＋到达过的全绿/台面红"
- 本票那枚真腿用例在**整包里**（⛔ `-run` 单跑）也绿：发3 名册逐字 `--- PASS: TestRunPacketLightsTheRosterRowOfItsLiveL1Window (12.16s)`
- `./internal/... -count=1 -timeout 30m` **达终态**，两发四数（分母取自 HEAD 导出树）：
  发1（`logs/40-full-internal.txt`）＝`RUN=0 PASS_TOP=0 FAIL=11 SKIP=0 rc=1`（非 `-v` ⇒ 那两枚恒 0，名册以 `^--- FAIL` 计；包级 `ok`／`FAIL` 逐行在件内）
  发2（`logs/44-internal-with-models.txt`，补 `models/` 后）＝`FAIL=10 rc=1`，包级 `ok=23／no test files=4／FAIL=3`（`internal/agent/approval`、`internal/panel`、`internal/risk`）
  逐名归因（发2 那 10 枚）：`internal/panel` 8 枚＝其中 7 枚同发1 的台面／他票／前端地界（两枚跨语对照键差 `Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`＝HEAD 跟踪的前端件落后于机主未提交件，本腿⛔ 按其调色、⛔ 读 `frontend/**`；`TestC21DesignTokensFourWayAgree`＝四处令牌不一致；`TestComposerRenderFixtureTellsTheTruth`；`TestReadGitOnThisRepositoryIsSelfConsistent`＝无 `.git`；两枚 `TestStreamLog*35r8`＝票 35 堆上限）＋第 8 枚 `TestPanelHostIsAttachedAndNamesTheWindowHops`＝**负载／临时目录形**，红句逐字 `composer_dispatch_test.go:436: scan … for host capability: open ….gotmp\TestSealNarrowsAndNamesThePrincipalItRemovedBySID2895323098\002\la: Access is denied.`，低负载单发复跑＝**绿**（`logs/45-flake-confirm.txt`）；
  `internal/agent/approval` 1 枚 `TestTenOpsInOneToolCallGetOneConfirm`＝**墙钟计时形**，红句逐字 `batch_test.go:164: 确认界面从未收到提示（等待 2s 后放弃）`＋`fakes_test.go:193: 派生的测试协程未在 3s 内退出（闸门存在无限等待）`，同尺发1 该包为 `ok 3.456s`，低负载单发复跑＝**绿**（同件）；
  `internal/risk` 1 枚 `TestResolvePerCallBudget`＝本机 C26 每调用预算读数（红句自述 `1758386 ns/op = 1.758 ms/op (budget 1.000 ms, 1196 samples)`），⛔ 动 `internal/observe/thresholds.go` 一字、⛔ 本票射程
  ⇒ 票 220 自己的判据面在两发里都不是红名；panel 那四枚 220 用例我另发逐名复跑全绿（`logs/22-panel-postrestore-green.txt`：`--- PASS` ×4，`ok … 0.133s`，且还原后整包红名回到基线 8 枚）
- 缺的那一发读数（具名）＝〔**无活实例／前台可控世界**下的 `./cmd/wisp -v -count=1 -timeout 30m` 整包终态名册〕。
  本腿⛔ 停 PID 21892（机主 2026-10-15 19:03 按台账 `A988` 自起的 `build\wisp.exe`），⛔ 为让整包跑完而排除任何别的用例；〔缺世界〕交编排者裁。

## 能力类必答：这一格在产品界面上今天看得见的有几维

- **0 维**看得见。尺（现跑，射程＝HEAD 树全部非测试 `.go`）：`grep -rn 'BlockedOnApproval|blockedOnApproval' --include='*.go' internal cmd | grep -v _test.go`
  ＝**8 行**，全部是**生产者与注释**（`subagent_roster_197.go:127/132/236`、`pump.go:199/327`、`panel_pump.go:78`、`run.go:728`），
  **零枚 Go 侧读者把这一维交给页面**
- 出口的实读：产码里 `Publish()` 只有 `cmd/wisp/panel_pump.go:446` 一处，其 `Out`＝`run.go:732 rt.bookPanelSnapshot`（`panel_pump.go:360` 定义＝写这枚 run 的台账），
  加上 `run.go:742 rt.ui.publish`（控制台面）。`internal/panel/pump.go:15-20` 逐字自述"WHAT THIS FILE IS NOT: the transport…no Go -> page channel in this tree today"，
  `pump.go:366` 的错误文案自述最后一公里在票 33/35 名下⇒ **实现件 §3 那句"传输那一跳仍属票 33／35"＝本腿复认（带上面这两把尺）**
- 包上（页面若接得到）此刻并存的是**两维**：`blockedOnApproval`（本票那一格）＋`status/statusKnown/statusReason`（宿主登记的运行态那一维），二者互不替换；
  名册渲染那枚函数＝`internal/panel/subagent_roster_197.go:198 taskRosterSectionFrom`（行装配 `:226-252`，排序 `:253`），
  沿路：`internal/panel/pump.go:323-327`（读 `L1Windows`⇒`taskRosterSectionFrom(p.src.Tasks(), cards, waits)`）⇒`:327`⇒上面那枚 `:198`

## 门禁四数（本腿自己跑的，`logs/50-gates.txt`／`51-d22scan-detail.txt`／`52-gates-worktree.txt`）

| # | 尺 | rc | 读数与归因 |
|---|---|---|---|
| 1 | `go build ./...`（仓外 `full` 导出树） | **0** | 全树编译面净 |
| 2 | `sh scripts/d22scan.sh` | 导出树 **1** ⇄ 共享工作树 **0** | 导出树那一发全是 `[phantom-citation]`（ban #9／票 212）＝**台面形状**：`git archive` 那几枚目录里没有 `.scratch/**`，产码注释引用的探针件在树上不存在（逐名 15＋ 枚在 `51-d22scan-detail.txt`，红句自述 `repo HEAD violates: …`）；同一把尺在共享工作树逐枚 clean＝`d22scan: clean - no D22 ban violations`，射程＝bans #1-5 `internal/`240／`cmd/`41、ban #6 `frontend/`116、ban #7 `internal/tools/`25、ban #8 `design/`39＋`frontend/`116＋`internal/`577＋`cmd/`139 Go 件 ⇒ ⛔ 计入判据、⛔ 当代码红 |
| 3 | `sh scripts/check-path-length-budget.sh --with-self-test` | 导出树 **2** ⇄ 工作树 **0** | 导出树那一发逐字自述 `check-path-length-budget.sh: not inside a git work tree - refusing to guess a denominator`＝那把尺**拒猜分母**（⛔ 冒充"门过了"）；工作树那一发 `VERDICT GREEN`：over-budget 60／名册覆盖 60／名册外 **0**、最长相对名 185 字符、字节视角与字符视角选中同一批 60 枚（现跑一致） |
| 4 | `go vet ./cmd/wisp/` 与 `go vet ./internal/panel/`（导出树） | **0**／**0** | 无输出＝只证编译面，⛔ 用作任何一格的判据 |
| 5 | `gofmt -l cmd/wisp internal/panel` | rc=0，列表非空 | 只点名存量件 `cmd\wisp\models.go`（与 `220-r2` 那发同名单）；本腿写面＝**零枚 Go 文件**，⛔ 代修存量债 |


## 残余（具名，⛔ 记进任何一格 AC）

1. **卡片那一半在生产里今天点不亮行**（本腿现量，⛔ 本票判据、属能力缺口）：`internal/agent/loop.go:603-608 callCorr` 逐字
   `if callID == "" { return fmt.Sprintf("%s#call-%d", taskID, index) }`／`return taskID + "#" + callID`，并在 `:676` 作为 `CorrelationID` 交给桥；
   ⇒ 生产侧卡的 corr 恒为 `<taskID>#<callID>` 形，`#` 后那一段是 `call-<i>`（含字母）或 provider 的 call id，
   而名册去牙识别器 `splitClashSuffix`（`subagent_roster_197.go:277-288`）只认**末 `#` 后全 ASCII 数字**；
   队列改写那形（`queue.go:163-168`）要求两枚**同 corr** 的卡才发生，按上面那形今天的循环造不出碰撞
   ⇒ **净效果**：`blockedOnApproval` 今天靠 L1 窗口那一半点亮（本票 AC#2 那一跳，真 run 已证），
   L2 卡片那一半只在"provider 给纯数字 callID"那一形才亮＝**票 220 射程外的既有缺口**，与编排者 10-15 那条读数同族（`docs/reports/pending-and-issues.md:14937`：真 provider 若给纯数字 `callID` 卡片那一半才点亮、"树内不可复现"），归口＝编排者裁（建议另立一枚：真 run 上种两枚 L2 卡的正控）
2. **没裁完的格＝AC#5 一格**，缺的那一发读数具名＝〔无活实例／前台可控世界下的 `./cmd/wisp -v -count=1 -timeout 30m` 整包终态名册〕（三发全部未达终态：两枚拦腰用例逐名＝`TestAC13BringUpSurvivesAReusedThreadQuit`、`TestAC14GoSideEvalPushReachesThePage`）。
   AC#1／AC#2／AC#3／AC#4 四格读数齐、已裁死；`./internal/...` 那发**达终态**、已并入 AC#5 的读数面
3. **台面定式再加一枚（写给编排者的七件套清单）**：`models/`（HEAD 跟踪 2 枚）不在清单里 ⇒ 未补时 `internal/models` 三枚红＋`cmd/wisp` 两枚红全是"committed models/manifest.json not found"那一形（本腿两棵台面各自实测：补进树后逐枚转绿）
4. 台面/他票噪声（HEAD 导出树 `./internal/panel/ -count=1` 基线 8 枚红，逐名归因）：
   `TestGitDimensionHasNoModelCallableTool`／`TestC21DesignTokensFourWayAgree`＝台面缺件形（红句自述 `read docs/PLAN.md … cannot find the path`／`read design/assets/tokens.css …`，`bench` 树没这两枚目录；`full` 树补齐后应消，见 AC#5 节）；
   `TestReadGitOnThisRepositoryIsSelfConsistent`＝导出树没有 `.git`（红句自述"这棵工作区目录及其以上都没有 .git"）；
   `TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`＝跨语对照在 **HEAD 跟踪件**上为红（红句自述 `Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`），前端未提交件归机主外部 agent，本腿⛔ 按其调色、⛔ 动其字节；
   `TestComposerRenderFixtureTellsTheTruth`＝渲染夹具自述 `render fixture verified across 0 painted states`；
   `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`＝票 35 的堆上限用例，非本票射程
4. 工作树 `go.mod` 那枚未提交改动（mtime `2026-10-10 22:28`）＝**还脏**，本腿⛔ 碰／⛔ 提交／⛔ revert，归编排者处置
