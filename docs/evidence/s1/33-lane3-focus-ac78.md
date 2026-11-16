# 33 lane-3（2026-11-16）— `:76` 焦点往返 ／ AC#7 ／ AC#8：读数与牙

写腿＝33 lane-3（本波唯一 `cmd/wisp` 写腿）。基线 HEAD＝`d28826829`，开工时 `git status --porcelain -- cmd internal`＝0 行。

## 0. 载具与纪律（自证）

- `go test ./cmd/wisp` 全部经载具 `scripts/wisp-cli-staged-tests.sh`（`D:\work\soft\Git\bin\bash.exe`），一次一发；每发跑前 `Get-Process wisp.test`＝空（三发共查三次）。
- 夹具在盘（`third_party/spike-models`＝True）⇒ 载具走**全分母模式**（"voice fixture environment resolved"），无名册跳过。
- 仪器坑一枚（给下一程）：PowerShell 原生调用 bash 时把 `-test.run` **在点号处拆开**（实测参量到达 bash 时为 `-test`／`.run`／模式三段），正是载具头注写的调用方错误形状（"flag provided but not defined: -test"）。解法＝整条命令作为一个字符串走 `bash -c '<整条>'`。第一发就这样白跑了一次（0.00s FAIL，零测试被执行），载具没有误报。
- 变异三发（A/B/C）全在 `cmd/wisp` 产码上，逐发：改 → 载具单跑 → 还原 → `git diff HEAD`＝0 行核验。零枚突变体进提交。
- ⚠ 复原方式披露：M-E（internal/panel）那发用了 `git checkout -- internal/panel/pump.go`（**单文件 pathspec**，非被禁的 `checkout .` 全树形；该文件当时唯一改动就是我自己那发突变，还原后 `git diff HEAD`＝0 已核）。

## 1. `:76` 焦点往返（automated 半）＝**环境红，具名跳过＋数字**

尺＝`TestAC4FocusReturnToPriorWindowGap33r5`（`cmd/wisp/panel_host_windows_test.go:1042`，33-r4 建尺、33-r5 落产码修），**一字未改、未放宽**。载具全分母模式实跑，**FAIL 5.26s**，红因＝断言 1。

**逐字读数（`panel_host_windows_test.go:1120` 的日志行，HEAD d28826829）**：

```
foreground before any panel 0x94e0122 | the ruler's own editor window 0x8aa0cf6 |
foreground while hidden (prior) 0x94e0122 | after Show 0x94e0122 | panel hwnd 0x115705d8 |
prevFocus recorded at Show 0x94e0122 | after Hide 0x94e0122 |
Hide attempted restore to 0x94e0122 (SetForegroundWindow 1, SetFocus 0)
```

**桌面事实（独立探针，13:1x 同窗取数）**：前台＝`0x94E0122`，属 pid 44296＝**Weixin（别的进程）**；Wisp 球 pid 49348 在跑但**没占前台**。

**判读（句柄事实，非时延）**：
- 断言 1 红（"panel did not take the foreground on Show"）：根因是 **Windows 前台锁**，不是产码焦点簿记——尺自己的 editor（`0x8aa0cf6`）从未拿到前台（它的 `SetForegroundWindow` 被拒，`waitForegroundIs` 预算耗尽时前台仍是 `0x94e0122`）⇒ 往返的第一环"editor 持焦点"在这台桌面上**没有 subject**。
- 断言 2 **PASS**：宿主记下的 prevFocus＝`0x94e0122` ≠ 面板 hwnd `0x115705d8`（没把自己记成来处）。
- 断言 3 **PASS**：Hide 之后前台回到 prior `0x94e0122`，`SetForegroundWindow` 答案＝1（跨进程回还真成功，还给的是 Weixin 的窗）。

**具名跳过**：按派单——`:76` 在"前台被别的进程占着"时取不准，本轮读数如上，**两个方向都不当产品判定**；不杀机主进程、不动机主窗口。与 `:74` 同一条队：**等安静桌面复量**。

**manual 半（留机主，具名）**：需机主看一眼——随便开一个有文字的编辑器（如记事本），按面板热键（面板弹出），再按热键收起面板，看**闪烁的光标是否回到编辑器里、接着打字落在编辑器里**。一眼即可，无需数字。

## 2. AC#7 ＝ 能力型判据：全量包字节在生产链上真被读

`Marshal()` 直接非测试调用者：仍 **0**（grep 只命中定义/测试），按派单判据**不数调用者、数能力**。

**全量包字节走的生产链（file:line，本轮逐一现读）**：
1. 触发：状态移动 → `rt.publishPanelSnapshot()`——装配点 `cmd/wisp/run.go:1347`（`rt.ui.publish`）、`run.go:2071`（consoleSink.publish）、`cmd/wisp/resident_windows.go:666`（`holder.setSnapshotPublish(src.run.publishPanelSnapshot)`）。
2. `cmd/wisp/panel_pump.go:511` `rt.pump.Publish()` → `internal/panel/pump.go:583` Publish → `:572` Marshal：全量包字节在此成型。
3. `Out = rt.bookPanelSnapshotBlocking`（接线 `cmd/wisp/run.go:1337`）→ `cmd/wisp/panel_pump.go:409-412` 留存 `lastSnap/lastSnapBytes/snapSeen`。
4. `cmd/wisp/panel_pump.go:523` `rt.pushPanelSnapshotToHost()` → **`cmd/wisp/panel_snapshot_deliver_329.go:95` `rt.lastPanelSnapshot()`＝全量包字节的生产读者** → `:105` `h.evalPanelJS(panelSnapshotEvalJS(data))`，字节**逐字嵌入、零再编码**（`:65-67`）。
5. 宿主：`cmd/wisp/panel_host_windows.go:661` `evalPanelJS` → 队列 → `:728` `deliverPanelPush` → `:772` `w.Eval(js)` → 页面全局 `window.wispOnPanelSnapshot`（`panel_snapshot_deliver_329.go:57`）。
6. 常驻接线：`cmd/wisp/resident_windows.go:618-620` `src.run.setPanelSnapshotHost(rp)`——本进程唯一的 PanelManager 成为推送目的地。

**能力尺（绿）**：`TestBookedSnapshotReachesHostBitForBit329`——三条尺（marshal 侧＝留下的字节 ≡ 测试自己的 json.Marshal；push 侧＝宿主收到的 JS ≡ `panelSnapshotEvalJS(那串字节)`；state-pin＝包里真带着队列那一张卡）＋收尾包＋台账 digest 方程。**PASS 2.98s**（13:20 发）与 **PASS 2.24s**（13:33 复认发）。

**牙（MUTANT-A）**：`pushPanelSnapshotToHost` 开头加 `return`（生产读者被摘）⇒ 该尺 **RED**：`:141` push 侧红（宿主收到 `""`，want＝完整逐字包）＋`:164` 收尾包红；marshal 侧与 state-pin **保持绿**＝红被归因到被摘的那一跳、别无连带。还原后 diff＝0。

**F-PUMP-4 撕裂形（`internal/panel/ac7_tear_experiment_test.go`，三发全 PASS）**：`TestAC7TornPacketIsObservable`（T1 的卡＋T2 的 result 同包，撕裂可观测）／`TestAC7PumpLockDoesNotCloseTheTear`（发现成立：泵锁整程持有包**仍撕裂**——票面开的"持锁一次性快照"药方**量出来不成立**）／`TestAC7PumpLockGuardsOnlyTheCounters`。

## 3. AC#8 ＝ 两枚已知分支＋本轮"量一次"普查

**ⓐ 超 440 丢 ids 分支**（`cmd/wisp/panel_pump.go:460-466`，`summaryClamp=:473`）：绿尺 4/4 PASS（`BooksThroughThePumpEntry`／`DropsIdsAndKeepsTheDigest`／`ShortPacketKeepsItsIds`／`TryEntry`+`Blocking` 对）。**负控分得开**：短包保留 ids、超包 `pending=-` 且 depth+sha256 仍钉住包。
**牙（MUTANT-C，本轮独立复量；r14 的 M1 早已记过）**：`if false && len(line) > summaryClamp` ⇒ `DropsIdsAndKeepsTheDigest` **RED 双牙**（`:105` 2347 字节包 ids 仍逐字列出；`:114` 行 496 runes > 440），`ShortPacketKeepsItsIds` **同发保持 GREEN**＝"丢"与"不丢"两种状态分得开，负控不是恒真。还原后 diff＝0。

**ⓑ EvToolStart：已不再不可达，复核读数**：
- **首个生产发射者**＝`internal/agent/loop.go:651-662`：`stream()` 把 `llm.EvToolCallStart` 原样转成 `agent.EvToolStart`（TaskID＋适配器自己的 CallID/ToolName）。
- **消费者**＝`cmd/wisp/run.go:2747` `case agent.EvToolStart:` → 打印＋`c.runFeed.ToolStarted`（`:2753`）＋`changed = true`（`:2754`）→ publish（`:2812-2814`）。
- **红尺**：`Test145RunningUsageAndDoneTravelTheWire` PASS 0.00s；**MUTANT-B**（整个 case 块从 consoleSink.Publish 摘掉）⇒ **RED** `panel_run_145_test.go:174` `running = [], want the start event's own call in flight`。还原后 diff＝0。

**本轮一次性普查（AC#8"接出站通道量一次"）**：`internal/panel` 整包 `-coverprofile`，`ok 25.3s`、**86.6% statements**。`pump.go` 零执行块分类：
- **真行为分支、当时零执行、本轮各补一枚会响的用例**（新件 `internal/panel/ac8_zero_branch_33lane3_test.go`，两发 PASS）：
  - `pump.go:830-836` `StreamLog.Close` 未知键分支（"stream ended and nobody had noticed it start"）——普查前**所有**现存调用者都只 Close 自己 Append 过的键（`pump_test.go:176/:226`、`subagent_stream_197_test.go:162`），分支零跑。新尺＝`TestAC8CloseUnknownKeyStillRecordsADoneChunk33lane3`（幽灵键留 done 空 text 行＋已 Append 键不动的内控）。**MUTANT-D**（`if false && !ok`）⇒ RED `:47` "left no row at all (1 chunks)"。还原后 diff＝0。
  - `pump.go:414-419` credential 读者空答→`CredentialUnknown`——普查前没有一发测试用空答驱动它。新尺＝`TestAC8EmptyCredentialReaderAnswersUnknownNotSilent33lane3`（三臂：空答→unknown＋known=true；真答→逐字；无读者→known=false＝票 216 AC#4 那句两句话规则）。**MUTANT-E** ⇒ RED `:77` 空答逐字穿成 `""`。还原后 diff＝0。
- **守卫块（nil 接收者/空键），按形状分类留置**：`:351`（Snapshot nil）／`:575`（Marshal nil）／`:585`、`:590-592`（Publish nil＋marshal-error 防御支）／`:610`（Publishes nil）／`:696`、`:708`（absorb/render nil）／`:780`（NewStreamLog maxKeys≤0 默认）／`:794`（Append nil/空 text）／`:826`（Close nil/空键）／`:844`（Chunks nil）／`:937`/`:949`/`:965`/`:981`（Truncated/Elided/DroppedKeys/DroppedCount nil）／`:1014`/`:1029`（TruncationFor nil/未知键）。
- `cmd/wisp/panel_pump.go` 的分支可达性：未做 covered-binary 跑（⛔ cmd/wisp 裸形状禁），可达性由上面的具名尺＋突变钉（clamp 支、两个 booking 入口、blocking/try 对）。

## 4. 门禁终态

- `gofmt -l` 新件：净（首写未格式化被 flag，提交前已 `-w`）。⚠ `gofmt -l internal/panel/ cmd/wisp/` 另列 `cmd/wisp/approval_always.go`＝**HEAD 上既存**、本腿未碰（diff＝0）、未顺手改。
- `go vet ./internal/panel/` rc=0；`d22scan` rc=0 clean（bans #1-5 internal/=405·cmd/=89；ban #8 internal/=1023·cmd/=313）。
- `internal/panel` 整包（含两枚新尺）：`ok 24.5s`。
- `cmd/wisp` 载具复认发（7 枚具名尺）：**7/7 PASS**、`ok 3.00s`。
- 写面：`cmd/wisp` **净 0 行**（三枚突变体逐字节还原）；`internal/panel` ＋1 新测试件；`docs/evidence/s1/33-*`＋票面进度行。
- 禁区零触碰：`thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43／`design/**`／`frontend/**` diff＝0；未 push；未放宽任何断言；票面 AC 框一枚未碰。

## 5. 没打完的

- `:76` automated 半：本轮**无决定性绿**——唯一桌面采样时前台属外进程（Weixin pid 44296），读数已录（§1），待安静桌面复量。
- `:76`／AC#7／AC#8 三格**均未勾**（勾归非实现者）；本腿只交读数＋牙。
- `:74` 延迟与 AC#10-scale：按派单一字未动。
