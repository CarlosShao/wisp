# 35-g3 背压格（票 35 `:49`）写码腿 — 推包队列有界 ＋ `failedPost` 读点 ＋ 四枚新尺 ＋ `heapDelta` 下溢

票 35 `.scratch/wisp/issues/35-panel-bridge-c17.md` 的 `:49` 格（「Backpressure」），含格内 10-08 11:3x 就地裁定。
派单＝`.scratch/wisp/probes/orch/20261115/txt/35-g3-backpressure-dispatch.txt`（逐字读过，四件事＋边界）。

本腿写面：`cmd/wisp/**`（本波唯一写腿）、本件。`internal/panel/**` **只读**（另一条腿在用；本腿一个字节没改，§⑤ 给出该改的人该怎么改）。

---

## ⓪ 起手锚点

- 时刻：`2026-11-15 18:4x +08`
- HEAD：`d9cc64ccc7db1703013a50f3286b5b1498fc7446`（本腿未 commit 前）
- 台件目录：`.scratch/35g3/`（突变源＋overlay JSON，只建不删）
- 读数日志：`.scratch/35g3-targeted-{1,2,3}.log`、`.scratch/35g3-stability.log`、`.scratch/35g3/mut-m*.log`、`.scratch/35g3-wholepkg-after.log`

---

## ⓪′ 裁定遵从（先声明，再说做了什么）

**我按 10-08 那枚就地裁定做，不按格里的旧三句做。** 格 `:49` 原文三句是「merges，no unbounded memory（heap cap asserted），content integrity kept」；10-08 11:3x 裁定（票 `:51` 逐字）把 **「merges」和「content integrity kept」两句判给票 197 leg B 的重裁取代**（`internal/panel/pump.go:632`：「Overflow TRUNCATES and never merges」），**只剩「no unbounded memory（heap cap asserted）」这半句还算分**。

因此本腿：

- **只实现并只认领**「推包队列无界内存」这半句（§①）；
- **⛔ 没有把 pump 改成会合并**，也 ⛔ 没有把「merges」这句旧话以任何形状搬回来（pump 一个字节没动；`internal/panel/pump.go:632` 的规则原样，逐字 `Overflow TRUNCATES and never merges (ticket 197 leg B re-cut this rule, ...)`，它的具名尺是 `internal/panel/pump_test.go:165 TestTheStreamLogTruncatesInsteadOfMerging`）。派单引的 `pump.go:397` 在本树已长到 `:632`（文件变长，行号漂移；条文逐字一致）；
- **不认领「content integrity kept」**：推队列按设计**会丢中间包**（整状态快照，新的盖旧的），丢了就计数（`pushCoalesced`）——这与 197 leg B 的「截断、不合并」同向，而与旧句「事件一个不丢」相反。本腿新尺只认一条更窄、可查的话：**交给页面的是一个完整的包，字节不动，且是最新那个**（§③ 第二枚）；
- 我做的「合并」只在**推的频率轴**上（coalesce：只留最新那一句），⛔ 不在内容轴上（不拼接/不改写/不切片），报告 §⑤ 有逐字边界。

---

## ① 队列形状：上限值＋依据＋落点

### 形状

`PanelManager` 上一条**有界、合并式（coalescing）**队列（`cmd/wisp/panel_host_windows.go:334-342` 字段与长注释）：

- `pushPending []string` —— 待送句子，**上限 `panelPushPendingCap` 句**；
- 写入方 `evalPanelJS`（`:659`）在 `pushMu` 下 append，**超界丢最旧**（`m.pushPending[0] = ""` 后 `append(m.pushPending[:0], m.pushPending[over:]...)`，`:674-682`）——**⛔ 永不丢最新**，并且**把被丢那一格清空**，否则 backing array 会替被丢的包续命，「纸面有界」正是本格要结束的东西；
- 只有**第一条**推文安排唤醒：`schedule := !m.pushScheduled`，然后 `w.Dispatch(func(){ m.deliverPanelPush(w) })`（`:710-712`）；
- 消费方 `deliverPanelPush`（`:726`）在面板线程取走整批（上限即 cap 句），清 `pushScheduled`，**逐句整包 `w.Eval(js)`**（`:769-771`）——⛔ 不把两句塞进一次 Eval；
- 计数：`pushDelivered/pushCoalesced/pushSkipped`，读口 `PushStats()`（`:786`）。

### 上限值：`panelPushPendingCap = 1`（`cmd/wisp/panel_host_windows.go:352`）

**为什么是 1（依据，不是口味）**：这个 hop 送的**是整状态快照**（`panelSnapshotEvalJS` → `window.wispOnPanelSnapshot(<整包>)`），页面收到后**整包替换**；任何尚未送达的旧包，被新包盖掉后**对屏幕零影响**。所以「留 N 句」中 N>1 的每一句都**只买内存、不买语义**。界＝「1 句等待 ＋ 1 句在飞」。

**为什么不是「永不触发的大数」**：值在代码里是**载荷**——裁剪循环读它（`:674-682`），改大即改行为；突变 m5（1→64）当场把四枚尺打红（§④），pinned 从 4145 B 涨到 265216 B。

**为什么不是 0**：0＝连最新也不留＝丢最新，是本格 ⛔ 明令禁止的那一侧。

**内存依据（只读，一字节没改）**：`internal/observe/thresholds.go:23` 的 `memCapPanel int64 = 600 << 20`（600 MiB，panel 态的整行内存上限）。它说明的是「**一份 panel 状态**的预算量级」；推队列的风险不在这 600 MiB 之内，而在「**每事件一份拷贝 × 无界 N**」——600 MiB 除 4 KiB 是 15 万句，洪水一秒就能到。所以界落在 **N=1**，队列内存从此与洪水速率无关。

### 盘上读口（两条，都只每 burst 一次，日志本身也有界）

- `panel.push.superseded`（`panel_host_windows.go:703`）——**队列第一次装满时**立即写一条：消费者**永不醒来**也留痕（带 `queue_cap`、`coalesced_total`、`delivered_total`、`pending_bytes`）；
- `panel.push.burst`（`:754`）——**burst 结束时**一条，带 `burst_superseded`（本 burst 到底盖掉几句）、`coalesced_total`、`delivered_sentences`、`skipped_total`。

> 为什么两条：第一条能证明「堵住了」，第二条才能给出**总数**——第一次装满时洪水还没结束，总数那时不存在。只留第一条，盘上永远只写着 `coalesced_total=1`。

---

## ② `failedPost` 的真实读点（第二件事）

改前：`failedPost` 只有测试读（35-recon ⑤ 层 4 量到 non-test readers = 0）。改后**两个非测试读点**：

1. **每条拒绝一条记录**（`cmd/wisp/panel_resident_windows.go:517`）：
   `slog.Warn("panel.post.refused", "detail", ..., "failed_posts", refused, "queue_cap", cap(rp.tasks), "window_live", false)` —— 拒绝分支可达时 `wRef` 必为 nil，所以 `window_live` 写字面量 `false` 而不是解锁后再读 `rp.wRef`（那会是一处数据竞争）；
2. **boot/状态行**（`panel_resident_windows.go:807`）：`"... window %s, %d show request(s), %d toggle(s), %d refused post(s))"`，末位 `rp.failedPost.Load()`。

**从盘上读得到**：`TestPanelPostRefusalHasAReadPointOnDisk35r9` 走**真 `post()`**（`tasks` 是无缓冲、无读者的 real channel ⇒ 必拒），拒两次后：`failedPost==2`、状态行含 `2 refused post(s)`、`installLogSink` 的 JSONL 里有**两条** `panel.post.refused`，末条 `failed_posts=2`、`queue_cap=0`。三条断言同发成立，读点不是「代码里有个字段」。

---

## ③ 四枚 35r8 尺 → 四枚**新**尺（第三件事）

四枚旧尺在 `internal/panel/backpressure_bound_35r8_test.go`（本波只读，**一字未改**）。它们**没毛病**，毛病是**射程**：**四枚量全是 StreamLog 与一只内存 sink，没有一枚碰过 `Dispatch` / `dispatchq` / 推队列**（35-recon ⑤：`35r8 四枚现跑全绿 … 量的是 StreamLog；⛔ 没有一枚尺打在推队列上`）。

| 旧尺（行） | 它量什么 | 为什么对本格不够 | 新尺（`cmd/wisp/panel_push_backpressure_35r9_windows_test.go`）| 它量什么 |
|---|---|---|---|---|
| `:254 TestStreamLogFanOutFloodKeepsRetentionUnderCap35r8` | StreamLog 扇出保留量次线性 | 量的是 pump 的行账本；推队列根本不在那条链上 | `:460 TestPanelPushQueueFloodUnderBlockedConsumerStaysBounded35r9` | 400/1600 推同 4 KiB 包，**各只钉 1 句**、只 1 次唤醒、`coalesced=attempts-1`、cap 是出货值 |
| `:321 TestStreamLogFloodBelowKeyBoundIsBounded35r8` | 低于 key 界时 StreamLog 有界 | 它的「堵」是 `Out` 通道；页面队列的堵是 `Dispatch` 闭包，两回事 | 同上（1x/4x 两读，**两个都断言**，含 heap 维） | 同上；本机 1x heap 真的下溢，读作 underflow（§②′） |
| `:385 TestStreamLogDroppedNamingLedgerIsBounded35r8` | 丢了什么、账本有界 | 它数的是被丢的**行**；推队列丢的是**整包**，且必须「整包」而不是「包的一部分」 | `:552 TestPanelPushQueueDeliversTheNewestWholePacket35r9` | 送出去的**是最新那一整句**（逐字节）、全局只调 1 次、`seq=attempts-1`、中间包**计入** coalesced |
| `:421 TestFloodRulerFiresOnAnUnboundedSink35r8` | 正控：无界 sink 上尺会响 | 正控只对**它自己那只 sink** 成立；推队列的无界形是 HEAD 的「每推一次 Dispatch」 | `:670 TestPanelPushRulerFiresOnAnUnboundedQueue35r9` | 同一台架、同一谓词，把界拿掉（每推一次 Dispatch）：400/1600 都判 OVER；**同形洪水走出货队列必须不响** |

**另加两枚（不是替换，是本格缺的另外两面）**：

- `:600 TestPanelPushQueueSupersessionsAreReadableOnDisk35r9` —— 从 `<dataDir>/logs/*.jsonl` 读回 `panel.push.superseded`（首装）与 `panel.push.burst`（总数 299），即「损失是盘上的事实」；
- `:765 TestPanelPostRefusalHasAReadPointOnDisk35r9` —— §② 的读点（`failedPost` 计数＋状态行＋盘上两条记录）。

**为什么新建四枚而不是改旧的**：① `internal/panel/**` 本波只读（另一条腿在用）；② 旧尺的主体（StreamLog）与推队列**不是同一条队列**，改它们会把两件事搅在一起，等于用「StreamLog 绿」冒充「推队列有界」——正是 recon ⑤ 点名的空洞。旧四枚仍全绿、仍守它们自己的主体。

---

## ②′ `heapDelta` 的 uint64 下溢（第四件事）

### 洞（在只读区，本腿不改，只精确报给能改的人）

`internal/panel/backpressure_bound_35r8_test.go`：

- `:64` `heapDelta uint64`（字段）；
- `:207` `heapDelta: after - before`（无符号相减）；
- `:237` `if m.heapDelta > heapCap`（把 ~2^64 判成「远超上限」）；
- `:303` / `:371` / `:411` 三处 `m1.heapDelta` 只 Print，**1x 的 heap 维从不断言**（所以这个洞至今没红）。

`after < before`（GC 还得比洪水分配的多）时 `after-before` 回绕到约 `2^64`，谓词立刻读成「远超 cap」⇒ **尺自己造的假红**。

### 精确补丁（给有写权的一条腿；本腿没动）

```go
// :64
heapDelta int64
// 新增
heapUnderflow bool
// :207 前后
delta, under := heapDeltaSigned(after, before) // after>=before ? after-before : (0,true)
// :237
if !m.heapUnderflow && m.heapDelta > heapCap { ... }
// :303/:371/:411 的 m1 行：把 heap 维从 Print 改成断言（与 m4 同谓词）
```

### 本腿的替代实现（`cmd/wisp` 侧，不等那条腿）

`panel_push_backpressure_35r9_windows_test.go`：`heapDelta35r9(before, after uint64) (int64, bool)` —— 回绕即 `(0, true)`；谓词 `overPushCap35r9` **拒绝把 underflow 读成 over cap**；**1x 与 4x 两读都断言**（heap 维不再是 Print-only）；`TestPanelPushHeapDeltaUnderflowIsNotReadAsOverCap35r9`（`:720`）钉住三件事：缩水读成 `(0,true)`、真增长仍判 OVER（修完没被阉）、旧无符号值在签名字段里塌成 `-600`。

**这台机现跑就撞上了这个洞**：1x 洪水的 heap 读数 `underflow=true, delta=0`（`35g3-targeted-3.log` 逐字）——正是旧形状会读成 `2^64-0` 的那种场景；新尺把它记为 underflow，其他四维照常判界。

---

## ④ 牙（A1880：删行为要红，不是断言源码串）

全部用 `go test -overlay`（`overlay-m*.json`），**不碰共享工作树**，突变源在 `.scratch/35g3/`。

| 突变 | 改法 | 结果（`mut-m*.log`）|
|---|---|---|
| m1 无界回退 | `evalPanelJS` 回退成 HEAD 形：每次 `w.Dispatch(func(){ w.Eval(js) })` | **4 红**：Flood / DeliversNewest / OnDisk / PositiveControl；`wakeups=400`、pinned 1657490 B |
| m3 丢最新 | 队列满时**不收新包**（保留最旧） | **2 红**：Flood（pinned 4142 ≠ 4145）＋DeliversNewest（收到的 `seq=0`） |
| m5 改大上限 | `panelPushPendingCap` 1→64 | **4 红**：pinned 265216 B ≫ 4145、`Delivered:64` |
| m2 删盘上读点 | `post()` 拒绝分支删掉 `slog.Warn`（只留计数） | **1 红**：PostRefusalOnDisk（0 条记录，want 2） |
| m4 删状态行读点 | `statusLine()` 删掉 `%d refused post(s)` 与实参 | **1 红**：PostRefusalOnDisk（状态行不再带数） |

m2 第一次生成时因残留未用变量编译失败（`declared and not used: refused`），**重生成后**才是上表那发；⛔ 没拿「编译失败」当红。

---

## ⑤ 运行读数

### 定向（本腿六枚尺）

```
go test ./cmd/wisp -run 'TestPanelPush|TestPanelPostRefusalHasAReadPointOnDisk35r9' -count=1 -v -timeout 300s
```

`35g3-targeted-3.log`：**6/6 PASS**，`ok ... 0.188s`。`-count=3`（`35g3-stability.log`）：**rc=0，3 连绿**。

逐字读数（1x / 4x，4 KiB 包）：

```
1x waiting=1 bytes=4144 wakeups=1 coalesced=399 heap+0B(underflow=true)
4x waiting=1 bytes=4145 wakeups=1 coalesced=1599 heap+13304B
bound waiting=1 pinned=4145 heap=1065156
```

### 整包（A1882：锚点／最后 `=== RUN`／本波新尺在截断点哪一侧）

命令（车＝派单里那条，`-timeout` 放宽到 1800s，防 24 分钟级整包被测试超时误杀）：

```
$env:PATH = "...\third_party\sherpa-onnx;...\build;$env:PATH"
go test ./cmd/wisp -count=1 -timeout 1800s -v
```

**发 1（原样车，`35g3-wholepkg-after.log`）**

- 锚点 `TestZZ315A1`：**未命中（0 次出现）** —— 进程在 `438.462s` 消失。
- 最后一条 `=== RUN`：`TestAC4PriorFocusSurvivesARefusedPanelSample`（它自己的 teardown 日志是全文件末行）。
- 本波新尺：**在截断点之前**，六枚在日志第 `3387-3418` 行、全部 `--- PASS`（Flood 0.03s／DeliversNewest 0.01s／OnDisk 0.02s／PositiveControl 0.03s／HeapUnderflow 0.00s／PostRefusal 0.01s）；截断点在 3952 行。
- 逐数 `PASS=357／FAIL=0／SKIP=1`：**没有一条 `--- FAIL`**，包行仍 `FAIL` —— 整进程消失（票 314 的硬形态）。

**发 2（把发 1 的死亡点那一枚测试 `-skip` 掉，`35g3-wholepkg-skip-death.log`）—— 整包跑到底**

- `-skip 'TestAC4PriorFocusSurvivesARefusedPanelSample'`（只跳这一枚，其余全跑；已知红一枚没跳）。
- **锚点命中：3/3 PASS** —— `TestZZ315A1RootLoopFirstBridgeCall`（0.21s）、`TestZZ315A1SiblingChildFirstBridgeCall`（0.29s）、`TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`（0.21s），三枚就是最后三条 `=== RUN`。
- 最后一条 `=== RUN`：`TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`（即锚点第三枚，跑到尾了）。
- 逐数 `PASS=618／FAIL=41／SKIP=10`，`752.948s`。
- **本波相关面红册＝空**：41 枚红**全部**是语音链（`speech: onnxruntime 1.17.1 is older than the 1.28 … a second onnxruntime.dll is shadowing the one on PATH`，System32 里那份 1.17.1 抢先）＋派单点名的已知红 `TestAC246DevLegIgnoresTheTestTaskInjection`（37.74s）；`TestPanel*`／`TestPanelPush*`／`TestPanelPost*`／`TestAC4*` 无红。**HEAD overlay 复算**（`.scratch/35g3-voice-head.log`）：`Test15eSilenceNeverCallsTheHandOff`／`Test26R3AnnounceClosesTheCaptureGate`／`TestWireR1SilenceAndEmptyTextNeverSpeak` 三枚在 **HEAD 下同样红、同一句 onnxruntime 1.17.1** ⇒ 这 41 枚与本腿无关，是本机 DLL 环境。
- 另两枚派单点名的已知红这一发是**绿**的：`TestAC4FocusReturnToPriorWindowGap33r5` PASS（1.33s）、`Test320FetchLegWritesTheKeyNowhereOnDisk` PASS（0.11s）。

### 整包进程死亡（本机既有现象，非本腿）

发 2 已经跑到底（锚点三枚是最后三条 RUN，rc=1 由 41 枚语音红造成，不是死亡）；发 1 和基线则死在真窗区：

- **改前基线**（`.scratch/35g3-baseline-cmdwisp.log`，HEAD＋无本腿新文件）：435.9s 后**整进程消失** —— 逐数 `--- PASS=351`、`--- FAIL=0`、`--- SKIP=1`（**一条 `--- FAIL` 都没有**），最后一条 `=== RUN` 是 `TestAC4PriorFocusSurvivesARefusedPanelSample`（它的 teardown 日志就是全文件末行），**锚点未到**，末行 `FAIL ... 435.945s`。
- **A/B**：同一条面板批（`TestPanel|TestStream|TestAC14|TestAC13|TestAC4`）在 **HEAD overlay** 下也死（22 PASS 后死于 `TestAC4PriorFocus…`，`35g3-panelbunch-head.log`）；**本腿代码**下同样死（29 PASS 后死于 `TestAC452PanelVoice…`，`35g3-panelbunch.log`）。⇒ 死亡与本次改动无关，是票 314 那族「真窗区整进程死亡」。

---

## ⑥ 门禁四数

- `go build ./cmd/wisp`：rc=0
- `go vet ./cmd/wisp/`：rc=0
- `gofumpt -l`（三枚改过的文件）：无输出
- `D:\work\soft\Git\bin\bash.exe scripts/d22scan.sh`：**rc=0，clean**（`35g3-d22scan.log` 末行 `d22scan.sh: OK - control rc=0, scan rc=0`；首跑两条已消：`cmd/wisp/panel_host_windows.go:310` 引了当时还没落盘的**本件**，`:349` 把测试文件名折行成 `panel_push_backpressure_35r9_` 被判 phantom-citation ⇒ 本件落盘＋改写叙述后复跑干净）
- `D:\work\soft\Git\bin\bash.exe scripts/check-path-length-budget.sh`：**rc=0，GREEN**（denominator 20703、over-budget 68、roster 68、not-in-roster 0、wall interval 0）

---

## ⑦ 未做／边界（诚实清单）

1. **⛔ 没 push**（本腿只 commit）。
2. **`internal/panel/**` 一个字节没改**：那边要改的正是 §②′ 的 uint64 洞；补丁已写清，留给有写权的腿。
3. **页面侧接收器仍缺**：`window.wispOnPanelSnapshot` 在那边的页面契约还没落地（票 329 probes），本腿只保证「Go 侧把最新整包送进 Eval」这半程。
4. **`PushStats` 还没接进 boot 状态行**：本腿按要求把 `failedPost` 接了两处读点；推队列的读口是 `PushStats()`＋两条日志记录。要不要再进状态行，属下一腿的排版决定（本腿没顺手改状态行，避免和 `resident_cancel_key_wording_260r3` 那类逐字断言打架）。
5. **整包颜色**（票 314）：原样车那一发**死**了（锚点未到），`-skip` 掉死亡点那一枚后**整包跑到底、锚点 3/3 命中**；但全包仍是 `rc=1`，因为 41 枚语音红（System32 的 onnxruntime 1.17.1 抢先，HEAD 下同样红）。**我没有为了变绿碰任何一枚语音测试或 DLL**，也没把「41 红」说成绿。

---

## ⑧ 交件判语

格 `:49` 现在还剩的那半句「flood events under a blocked consumer → no unbounded memory（heap cap asserted）」在**推队列**上**有界**且**被四枚新尺逐维断言**（句数、字节、唤醒数、heap），正控证明谓词会响，盘上读得到损失，`failedPost` 有了两个非测试读点。旧三句里的「merges／content integrity kept」按 10-08 裁定**不实现、不认领**。
