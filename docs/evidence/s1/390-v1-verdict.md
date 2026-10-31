# 票 390 终裁腿 `390-v1` —— dev `e12200317`（390-r1）上终裁（非实现者）

> **尺**：被裁实现＝`e1220031741`（3 枚 `cmd/wisp` +438/−8）＋证据笔 `867606544`（两枚均 dev 祖先，`merge-base --is-ancestor` rc 0；开工 tip=`41295d4d5`＝编排者 A1329 收件笔）。★**全部实验出自仓外导出树** `D:\tmp\wisp\20261031\390v1\export`（`git archive e1220031741`→Git Bash tar 解出 **18179/18179 文件零缺**，三件实现 `git hash-object` 与提交 blob **逐位同**）；纯净对照树＝`e1220031741^`＝`f62b0fc3b`（**零 390 字节**，同法 18178/18178）。★staging＝三 DLL 贴身（`third_party/sherpa-onnx` 全家：onnxruntime／sherpa-onnx-c-api／sherpa-onnx-cxx-api）＋`WISP_SILERO_MODEL/SPIKE_MODELS/TTS_DIR/KWS_DIR/15E_TMP` 五 env 指主仓资产＋junction `build`/`third_party`（两目录 **0 枚 `.go`**，分母不受影响）＋CWD=导出树根＋`-test.timeout 30m`＋GOTMPDIR 每跑置空＋Git Bash runner。★共享树纪律：`26p-r1`（cmd/wisp）/`387-r1`（internal/speech）/`386-r1`（tools/mockllm）在飞，本腿零碰零归因；⛔ 未改任何仓库文件（除本件）；⛔ push/amend/reset/rebase/stash/checkout。

## 0 一句话判语

**五格全翻、票 390 可带残余结案**：381-r1 的"需 internal/ball API 扩展"申报**被本腿坐实推翻**（OnClickBall 本就是既有事件面且 pre-390 就在消费它）；端到端 ≤400ms 本腿两跑 10 轮 28–50ms 全过（专用发 33/38/28/37/29、合跑发 35–50）；M1 反向（本腿**独立造**突变件）主案红句逐字＋Esc 侧 3/3 照常绿＋还原核绿；零回归红名集合零新增（唯一红 Test384 经**纯净对照树同句同点红**＋单跑 ×2 绿归因）；d22scan 分母四数逐枚命中且正控两发可逆。实现件读数本腿逐项复核**全真**，一处 staging 口径差（"四 DLL"实为三枚全家）具名 §7，不动摇实质。

## 1 ★★★ 裁点①「ball API 无需扩展」推翻本身复核（最关键裁点）

| 证据 | 读数（现跑） |
|---|---|
| 既有事件面字段 | `internal/ball/ball_windows.go:50`＝`OnClickBall func()`（`Events` 结构第 1 枚回调；`:48` 注释"Nil entries are skipped"） |
| 既有 fire 点 | `:632`＝`b.fire(b.opts.Events.OnClickBall)`——`wmLButtonUp` 非拖拽非 dock 分支（`:615-634`），逐字命中 |
| pre-390 就在消费它 | `e1220031741^:cmd/wisp/resident_ball_windows.go`＝`OnClickBall: func() { recordBallGesture("click") }`——**旧闭包本来就挂在这一事件面上**（把单击记为 unhosted） |
| 390 只换闭包体 | `git diff f62b0fc3b e1220031741 -- resident_ball_windows.go`：闭包体 `recordBallGesture("click")`→`rb.clickInterruptGesture()`；`internal/ball` **零字节**（该文件最后改动＝`b1e59d639` 10-04 票 260-r1） |
| 381-r1 原申报 | `docs/evidence/s1/381r1-impl.md:49` ②逐字："球单击→预留钩子但未装配执行器（**需 internal/ball API 扩展**，不在本片范围）" |

**判：推翻成立**。单击事件面在冻结包里既有、既有消费者在、装配只需 cmd/wisp 侧换闭包＋晚挂执行器——381-r1 的申报不成立，且属"API 已够"而非"停手交回"情形（A1329 定性维持）。AC#1 的"ball API 扩展须具名交回"分支**未触发**。

## 2 ★★★ 裁点② 端到端 ≤400ms 复刻（staged、真 waveOut 真 matcha、球手势入口）

`Test390r1BallClickRealDeviceWithin400ms`（触发＝`ball.clickInterruptGesture()`，非裸 TriggerInterrupt；click 落在出声 300ms 处）：

| 发 | unstick/reset（click→#27 走完／→waveOutReset 落地） | 控制轮 | 门 |
|---|---|---|---|
| 专用发 rc 0（34.9s） | **33/38/28/37/29 ms**（abort_returned 全 0） | `parked_to_done=2.108s` 照播完 | `gate 6/6`、`spoken=1 aborted=5` |
| 六族合跑内 rc=1 发 | **50/49/35/37/50 ms**（负载下同带） | `parked_to_done=2.116s` | 同上 |

**≤400ms：MET**——两跑 10 轮分布 28–50ms，与实现件五轮（33/39/22/31/39）及 381-v3 的 Esc 带（35–48）同量级。`CLICK-27 aborted=1 spoken=0 next_spoken=1 probe1_aborts=1 probe2_aborts=0 gate 2/2`（被打断句不计 spoken、下一句照常、stale abort 不及第二句）在专用发、合跑发、还原核三发逐字同读数。⚠ `announce_to_done=9.944s` vs `parked_to_done=2.108s`＝381-v3 §10③ 具名的 Drain 早回残余（归 26 侧），本腿读数与它一致、非 390 新账。

## 3 ★★★ 裁点③ 反向 M1 复刻（本腿独立造突变件）＋定向自证＋还原核

- **突变**：`clickInterruptGesture` 摘 `trigger()` 调用（attach／记录分支保留，`_ = trigger`＋`slog.Warn("…unhooked (390v1-M1)")`）。仓外 `D:\tmp\wisp\20261031\390v1\mut\`＋`-overlay`（映射**导出树**路径，工作树零突变窗口）。SHA256：本腿突变件 `BD37DD5F…` ≠ 原件 `E229F9DA…`（＝实现件申报的原件值，逐位同）≠ 实现腿突变件 `D6C8C44B…`——三件互异、独立复刻。
- **主案红**：`Test390r1BallClickPlayingInterruptAbortsAndWalks27` @ M1 二进制 **rc=1**，红句逐字＝`MUTANT 390r1-M1 SHAPE: the ball click did not stop the announcement in flight (spoken=1, aborted=0): the click gesture did not fire the attached TriggerInterrupt, so the sentence played out - the second trigger source is not wired (ticket 390: 球单击→TriggerInterrupt 真挂)`（测试 `:92` 自带断言句）。
- **定向自证**：**同一枚** M1 二进制上 `Test381r2PlayingPhaseInterruptAbortsAndWalks27`（`ABORT-27 aborted=1 spoken=0 … gate 2/2`）＋`Test381r2IdleInterruptIsInert`＋`Test390r1UnattachedClickDoesNotInterrupt`（`UNATTACHED spoken=1 aborted=0`）**3/3 PASS rc=0**——摘的是球单击这一跳，interrupt 本体与 Esc 半边分毫未动。
- **还原核**：无 overlay 原二进制重跑主案 rc=0、`CLICK-27` 同读数；导出树文件哈希复验 `E229F9DA…` 未动。

## 4 ★★ 裁点④ 正控（「无事发生」的那一半，原二进制 rc 0 3/3）

- **idle 球单击惰性**：`Test390r1IdleBallClickIsInert` PASS（0.00s）——Sleeping 不动、timeline 空、problems 空、`ttsAborted/ttsSpoken=0`；nil host／未挂接 host 单击不 panic（`:293-296`）。
- **UNATTACHED 反形**：`Test390r1UnattachedClickDoesNotInterrupt` PASS——`spoken=1 aborted=0 problems=0`，播报照播完；与主案**唯一 delta＝attach 那一行**（`:327` 注释逐字），挂接敏感性的测试内证明。

## 5 ★★ 裁点⑤ 零回归（红名集合比）＋延迟三数

| 发 | 四数 | 红名 |
|---|---|---|
| 六族合跑 @ 390 树 | **42 RUN／41 PASS／1 FAIL／0 SKIP** | 唯一红 `Test384ResidentWatchdogEndToEnd`（3.41s） |
| 六族合跑 @ 纯净树 `f62b0fc3b` | **38 RUN／37 PASS／1 FAIL／0 SKIP** | 唯一红**同一枚**（3.38s）**同句同点**：`resident_watchdog_windows_test.go:226: 384: never observed the tree back under the Armed cap (a real within-cap reading) within 3s` |
| `Test384` 单跑 ×2 @ 390 二进制 | 2 发全绿 rc 0 | ——合跑红＝3s 时窗的负载形状，单跑即绿 |

两树红名差＝恰 4 枚 `Test390r1*`（全绿）；实现腿报的 41a2 合跑时窗红本腿两树都未见（与"负载漂"定性相容）。**红名集合零新增成立**，Test384 红＝HEAD 既有、非 390-r1 回归。`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms`（合跑 `Test15e`）**逐位同**基线。

## 6 ★★ 裁点⑥ 线程形状核（381-v2「abort 从触发方线程出手」同形满足）

单击与 Esc **同一枚 `ball.Events` 字面量的两枚字段**（`resident_ball_windows.go:405` `OnClickBall: func(){ rb.clickInterruptGesture() }`／`:408` `OnCancelHotkey: …recordCancelHotkey(onCancelEsc)…`），由**同一个窗体过程**经**同一个 `fire()`** 内联调出：`ball_windows.go:632`（WM_LBUTTONUP）与 `:663`（wmHotkey/hkCancel）；`fire()`＝`:728-732` 直接 `fn()`，契约注释 `:724-727`「runs a consumer callback on the STA thread…No goroutine is spawned per gesture」。回调体非阻塞：`clickInterruptGesture`＝锁下读回＋一次调用（`:692-694` 注释）；`TriggerInterrupt`（`:548-562`）＝buffered send（default 臂）＋`speakAbort` 直接 abort——ui-sta 线程零等待零 Win32。Esc 侧链（`resident_audio_windows.go:395` `ra.approval.voiceInterrupt`→`TriggerInterrupt`；`vetoByEsc` `resident_approval_windows.go:645`）同线程同形状。**判：同在 ball ui-sta 线程、381-v2 的线程要求两触发源同形满足**（代码形状核；381-v3 §3 同法先例）。

## 7 ★ 裁点⑦ d22scan（导出树）＋正控两发

- `go run . -root <export>`（树内 `tools/d22scan` 自身）**rc=0 clean**；分母 **bans#1-5 internal/=286 cmd/=51、ban#8 internal/=726 cmd/=177**（另 design 14／frontend 116／tools 25）——**派单四数 286/51/726/177 逐枚命中**，与实现件申报逐数同。⚠ 本导出无 `.git` ⇒ 扫描器显式打出"gitignore NOT APPLIED（全扫方向）"；archive 只含 tracked 文件，分母与实现腿 clone 法逐数一致，非隐账。
- **正控两发＋可逆**：`✓`(U+2713) 串 ⇒ rc=1 `[emoji] cmd/wisp/zz_390v1_poscontrol_emoji.go:5`（string literals included）；裸 `go func(` ⇒ rc=1 `[bare-goroutine] :7`（use observe.Registry.Spawn）；两探针撤除复跑 **rc=0 分母逐数回落**。
- staging 口径差具名：A1326 卡写"四 DLL"，`third_party/sherpa-onnx` 实存**三枚全家**，本腿按三枚贴身跑，15e 真设备／41a2r1 KWS／381r2／390r1 全绿——三枚即足，"四"疑计口误，不影响任何读数。`WISP_PUNC_MODEL` 未设（`punct_attached=false`，`[voice] punctuation` 归独立票，无案需要它）。

## 8 ★★★ 五格逐条判 ＋ 收口判定

| 格 | 判 | 依据（本腿自跑，⛔ 采信实现件字面） |
|---|---|---|
| **AC#1** 球单击→TriggerInterrupt 真挂（对照 Esc 形状；ball 扩展须具名交回） | **✅ 可翻** | attach 行＝`resident_voice_windows.go:1717`（`startResidentVoice` 的 `if rb != nil` 块、紧挨 `attachWakePause`，同一所有权形状）；手势入口经 `clickInterruptGesture`（全部读数走球入口）；**ball API 无需扩展推翻坐实**（§1），冻结面零字节 |
| **AC#2** 端到端 ≤400ms＋门重开 | **✅ 可翻** | §2 两跑 10 轮 28–50ms 全 ≤400、控制轮照播完两发、门 6/6＋CLICK-27 2/2 三发同读数 |
| **AC#3** 摘挂接⇒红（三读数＋还原核） | **✅ 可翻** | §3 本腿独立突变件红句逐字＋同突变 Esc 侧 3/3 绿（定向自证）＋还原核绿＋三 SHA 互异 |
| **AC#4** 零回归（六族＋Esc 通路不动） | **✅ 可翻** | §5 两树红名集合零新增＋Test384 归因三发＋`262/640/902` 逐位同；Esc 通路案含突变下全绿 |
| **AC#5** 纪律格 | **✅** | `git show --stat`：恰 3 枚 cmd/wisp＋证据笔 1 枚；`internal/{ball,risk,session,speech}` 零字节（ball 最后改动 10-04）；d22scan clean＋正控可逆；票面/台账/push 零动作 |

- **可否带残余结案：可**。残余全非 AC 项，逐条定性：①**单击 summon 半仍无执行者**——`ballGestureWhy`/boot 行/未寄宿 bullet 三处文案已随真（diff 在证）；归口＝A1329 收件账＋码面具名（GUI 手势执行者家族，228 AC#2 面），**本腿未找到给"click→summon"单独立票的票号**，请编排者落一个具名归口或在台账明示按 228 面挂账；②**Test384 合跑 3s 时窗**——纯净树同红＋单跑 ×2 绿＝HEAD 既有负载形状，归口 384 侧，390 零触碰其代码；③（随行核）**Drain 早回**（`announce_to_done` 9.9s vs `parked_to_done` 2.1s）＝381-v3 §10③ 既有残余、账仍归 26 侧，本腿读数与之一致。
- 未定义即停三条全部未触发（①事件面既有＝§1；②线程同形＝§6；③381-r2 本体零字节＝`resident_voice_windows.go` 全 diff 恰 +7 行 attach＋注释）。

## 9 纪律核

写面＝本件＋`D:\tmp\wisp\20261031\390v1\**`（export/headtree/mut/bin/logs/cases 全仓外；两突变件、两正控探针已撤并复跑可逆）。共享树零字节改动、零归因；rc 一律重定向取值；GOTMPDIR 每跑置空；变异走仓外 `-overlay`、还原核在末发之后；音频全程单进程串行（错峰）；⛔ push。★踩坑一笔：pwsh 直调测试二进制把 `-test.v` 拆成裸 `-test`（381-v3 §1 已具名同款），改 Git Bash runner 后 rc 正确——复跑腿请直接用 bash。

---

**交回**：编排者——票 390 五格全翻、按 §8 带残余结案；残余①请落具名归口（summon 半）；381-r1 申报推翻请随本件在台账记一句（"需 internal/ball API 扩展"不成立＝OnClickBall 既有面）。
