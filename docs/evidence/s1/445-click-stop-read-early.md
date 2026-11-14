# 票 445 — 球单击的「停」在走落定前读了 D43 快照：AC#1 判定「被调用」＝仪器读早了，修仪器、验牙、普查

- 腿＝票 445 写腿（委派）；2026-11-14 17:0x–17:4x +08；起手 HEAD `c64a4fa22`（票面记的红在 `d2cfa96e8`）。
  本件所有读数都在**本腿工作树**与**仓外导出树** `D:\tmp\wisp\20261114\445-tree` 现取，逐条标了命令与 log 路径。
- 写面：`cmd/wisp/resident_voice_interrupt_390r1_windows_test.go`、`resident_voice_interrupt_381r2_windows_test.go`、
  `resident_voice_wire_r1_windows_test.go`、`resident_voice_acting_395r1_windows_test.go`、
  `voice_e2e_demo_16r1_windows_test.go`（五枚均 `*_test.go`）＋ 本件 ＋ `.scratch/wisp/probes/445/**`（新）＋ 票 445 Progress log 尾一条。
- ⛔ 未动：`cmd/wisp/**` 的任何产码（**0 字节**，§5.5）、阈值／golden／契约／D43 转移表、`scripts/**`、`frontend/**`、其它腿的未提交改动。⛔ push。

---

## 一、AC#1 结论：`TriggerInterrupt` **被调用** —— 红在仪器读早了，不是球与语音腿之间的真漏

### 1.1 两条发布路径，一个可见窗口（机理先写清，否则下面的数字只是噪音）

同一个「走」（D43 的一跳）由两步发布出去，**不是一个原子动作**：

1. `internal/statemachine.Dispatch` 在自己的临界区里就改了状态：`machine.go:144` `m.state = hit.To`，`:160` 解锁，
   `:162-164` 再跑该行的副作用 sink。**`Dispatch` 返回时，新状态已经对所有 `machine.State()` 的读者可见。**
2. 调用方 `cmd/wisp` 的 `rv.dispatch` 在 `Dispatch` **返回之后**才把这一跳追加进 `rv.timeline`
   （`resident_voice_windows.go:1436-1457`，`rv.mu` 在 `:1449` 才拿）。

所以存在一个**真实窗口**：`machine.State() == Listening` 为真，而 `rv.timeline` 里还没有 #27 那一跳。
测试件那枚完成标记等的正是状态：

```go
voice15eWaitFor(t, "the click to end the announcement (aborted or played out)", func() bool {
    return (rv.ttsAborted.Load() >= 1 && rv.machine.State() == statemachine.StateListening) ||
        rv.ttsSpoken.Load() >= 1
})
...
timeline := rv.timelineSnapshot()
last := timeline[len(timeline)-1]
if last.Event != statemachine.EvInterrupt || last.To != statemachine.StateListening {
    t.Fatalf("the click's stop walked %+v, want #27 voice.interrupt -> Listening", last)
}
```

标记被状态满足 ⇒ 立刻读 `timeline` 的最后一枚 ⇒ 可能读到**上一枚**（CI 里读到的就是 `#15 llm.first-token`）。
CI run `43575891409`（head `d2cfa96e8`，job `test-windows` step 10，RUN=846 PASS=603 FAIL=1 SKIP=19）那一行逐字：

```
cmd/wisp/resident_voice_interrupt_390r1_windows_test.go:107: the click's stop walked {From:Thinking To:Speaking Event:llm.first-token AtSample:0 Say:播报：D43#15 首个纯文本 token 让球走 Speaking（TTS 说这一句）}, want #27 voice.interrupt -> Listening
```

### 1.2 「到底被调用了吗」的读数（探针在仓外导出树，产码 0 字节）

探针 `.scratch/wisp/probes/445/probe-445-trigger-window_windows_test.go` 把**同一枪**同时记录三件事：
`attachVoiceInterrupt` 的 hop **真的进了没有**（把 `rv.TriggerInterrupt` 包一层 `triggerCalls atomic.Int64`）、
`ttsAborted`/`ttsSpoken` 的计数、以及**同一时刻**的朴素「读最后一枚」判定（就是 CI 那行断言）。

| 台面 | 分母 | 采样 | 命中（会红的那一读） | `trigger_calls` | `aborted` | `spoken` | 读数文件 |
|---|---|---|---|---|---|---|---|
| 第 1 批（单模式 spin） | 20 | 连续 | **3/20 = 15%** | 20 | 20 | 0 | `.scratch/wisp/probes/445/logs/06-session-readings.txt` §A |
| 第 2 批 | 40（20 spin ＋ 20 票面 5ms 轮询） | 两模式交替 | **1/20 spin ＋ 0/20 5ms** | 40 | 40 | 0 | 同上 |
| 第 3 批（头条） | 200（100 ＋ 100） | 两模式交替 | **4/100 spin ＋ 0/100 5ms** | **200** | **200** | **0** | 同上 |

**合并：连续采样 8/140 = 5.7%；票面那枚 5ms 轮询 0/120。**
三批共 260 发，每一发的 `trigger_calls` 都等于「第几发」（`200/200`），`aborted` 同样 `200/200`，
`spoken` 恒 0 —— 也就是说**球单击每发都真的调到了 `TriggerInterrupt`，每一次都真的把播报打断了**，
被打断的句子一次也没有被算成「念完」。逐发点名（每一行都带 `trigger_calls` 与 `aborted`）：

```
445-REDLINE round 25 (fast sampling): the click's stop walked {From:Thinking To:Speaking Event:llm.first-token ...}, want #27 voice.interrupt -> Listening (trigger_calls=25 aborted=25 state=Listening edges=75)
445-REDLINE round 117 (fast sampling): ... (trigger_calls=117 aborted=117 state=Listening edges=351)
445-REDLINE round 163 (fast sampling): ... (trigger_calls=163 aborted=163 state=Listening edges=489)
445-REDLINE round 165 (fast sampling): ... (trigger_calls=165 aborted=165 state=Listening edges=495)
```

这三批探针还有两个读数值得记下来：

1. **`settle`**（朴素读之后到 #27 真的进 `timeline` 的间隔）：中位数 `0s`（亚微秒），但**最长一发 6.04ms**。
   ⇒ 这个窗口不只是「几条指令的缝」：写者 `rv.dispatch` 可以在 `Dispatch` 返回与追加之间**被抢占**，
   这一发就把窗口撑到 6ms —— 这正是 5ms 轮询的仪器在 CI 里能撞上的量级。
2. 因此「连续采样 5.7% ／ 5ms 轮询 0/120」不是两个矛盾的数字，而是**同一个窗口、两种采样率**；
   CI 那一发红不是「偶然的调度」，是窗口被抢占撑宽时被 5ms 轮询逮住。

### 1.3 逻辑上的旁证（读数之外的独立一条）

要走到那行红，必须 `machine.State()` 从 Speaking 变成 Listening；D43 表里这条边**只有一个入口**：
`resident_voice_engines_windows.go:1228-1236` 的 abort 分支，条件是 `errors.Is(drainErr, errSpeakerAborted)` 且 `moved == true`。
`errSpeakerAborted` 只由探针喇叭的 `Abort()` 返回，而 `Abort()` 只由 `TriggerInterrupt` 调。
⇒ 就算不看探针读数，「到了红行」本身就蕴含「`TriggerInterrupt` 被调用过」。读数与逻辑两条独立路径同一结论。

### 1.4 判定

**AC#1 结论：被调用。** 按票面的岔口，这一支是**仪器读早了**，不是产码缺陷；
⛔ 没有把票 445 当 flake 结掉，也没有把它升格成产码缺陷 —— 它是一枚**真缺陷，缺陷在仪器**。

---

## 二、AC#2 修：等目标那一跳落定再判（带名字的超时行与数字），断言本身不放宽

### 2.1 两个 helper（放在 `resident_voice_interrupt_390r1_windows_test.go` 尾部，全包共用）

```go
func waitHopAfter445(t *testing.T, rv *residentVoice, from int, what string) []voiceTransition   // :421
func waitWalked445(t *testing.T, rv *residentVoice, from statemachine.State, ev statemachine.Event, to statemachine.State, what string) // :444
```

- 上界 **5s**、**1ms** 轮询（不是 sleep 一个拍脑袋的时长）；超时行把**数字**写进红句：
  边数、基线、`state`、`aborted`、`spoken`、最后一枚是什么。两枚 helper 都不吞错：没落定就是红。
- 判据是**断言本身**（「基线之后出现了新的一跳」／「表里有这一行」），不是「等一下应该就有了」。

### 2.2 五处同形位置，各自带读数

| # | 位置（改后行号） | 原来读什么 | 改成 | 说明 |
|---|---|---|---|---|
| 1 | `resident_voice_interrupt_390r1_windows_test.go:99`＋`:117-130` | 标记一满足就读 `timeline` 最后一枚（**CI 红的就是这行**） | 先记基线边界，`waitHopAfter445` 等落定，再判「最后一枚必须是 #27」，并加「只多了 1 枚」断言 | 断言文字未改；**新增更严**的「恰好 1 枚」 |
| 2 | `resident_voice_interrupt_381r2_windows_test.go:330`＋`:350-360` | 同一枚标记之后读最后一枚（同一文件同形） | 同上（`waitHopAfter445` ＋ 恰好 1 枚） | Esc 腿，机理相同 |
| 3 | `resident_voice_wire_r1_windows_test.go:364-377`＋`:389` | `ttsSpoken==1` 之后读「恰好 4 枚」与球的投影 | 先 `waitWalked445`（#26）再 `voice15eWaitFor` 等球的投影出现 Warm，然后才读 4 枚记录 | `ttsSpoken.Add(1)` 在 `:1265`，#26 派发在 `:1282` |
| 4 | `resident_voice_acting_395r1_windows_test.go:153-159` | 同上（`ttsSpoken` 后读记录） | 先 `waitWalked445`（#26）再取记录 | 记录是**日志＋断言**的共同来源，故把快照挪到等待之后 |
| 5 | `voice_e2e_demo_16r1_windows_test.go:329-335` | 同上 | 先 `waitWalked445`（#26）再取记录 | 同上 |

**为什么不是「睡久一点」**：`waitWalked445` 的谓词是「表里有这一行」，永远不会有；睡多久都不会绿。
§3 的变异 m2 就是拿这个当面量的。

### 2.3 修后读数（本腿工作树，主仓）

```
--- PASS: Test395r1VoiceToolTurnWalksActing (12.40s)
--- PASS: Test381r2PlayingPhaseInterruptAbortsAndWalks27 (7.17s)
--- PASS: Test390r1BallClickPlayingInterruptAbortsAndWalks27 (7.12s)
--- PASS: TestWireR1AnnounceEndToEnd (12.78s)
--- PASS: Test16r1VoiceDemoWakeAsrAgentSpeakSettleRealHardware (31.77s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	72.00s
```
（wire_r1 在补了「球的投影也等落定」之后又单独跑过一发：`--- PASS: TestWireR1AnnounceEndToEnd (18.21s)`。）

---

## 三、AC#3 牙：仓外导出树里两枚真变异，修后仪器**真红并逐字点名**，恢复同 sha

导出树 `D:\tmp\wisp\20261114\445-tree`（仓的快照 ＋ `third_party` junction，193.6MB/3608 文件，夹具与钉版 DLL 都是真的）。
脚本 `.scratch/wisp/probes/445/ac3-445.sh`，变异器 `.scratch/wisp/probes/445/mutate-445.py`，
日志 `.scratch/wisp/probes/445/logs/01..05-*.log`。
**仓本身一次都没被变异过**（sha256 见 §3.3；主仓这两个文件与导出树同 sha）。

### 3.1 两枚变异（都是「真缺陷」形状，不是「拔掉断言」）

| 变异 | 改了什么 | 世界变成什么 | 仪器读数 |
|---|---|---|---|
| **m1** `!TriggerInterrupt` | `resident_ball_windows.go` 的 `clickInterruptGesture` 里删掉 `trigger()` 那一行（执行器**仍挂着**、日志仍说 executed=true） | 球单击**真的不再调** `TriggerInterrupt`：播报念完（`spoken=1`、`aborted=0`） | **FAIL**，逐字：`MUTANT 390r1-M1 SHAPE: the ball click did not stop the announcement in flight (spoken=1, aborted=0): ...` |
| **m2** `!edge-recorded` | `resident_voice_windows.go` 的 `rv.dispatch` 里**不追加** `rv.timeline`（状态照常走） | 机器照样走到 Listening，**只有时间线沉默** | **FAIL**，逐字：`the D43 walk never landed within 5s: the click's stop (D43#27 voice.interrupt Speaking -> Listening): timeline has 0 edge(s) (baseline 0), state=Listening, aborted=1 spoken=0, last=<none>` |

读法：
- m1 证明「球→语音腿真漏」这个世界里仪器会**点名**（而且打的是 MUTANT 行，不是超时行）——这正是 AC#1 那一岔口要能分辨的两边。
- m2 证明修后的等待**不是橡皮图章**：状态到了、边永远不来 ⇒ 5s 后带数字红。**「睡久一点」在这里永远绿不了**。

### 3.2 四次运行的顺序与 rc

```
01-baseline-green          rc=0  PASS 11.22s   （导出树与仓同 sha，先证明台面是真的）
02-m1-trigger-never-called rc=1  FAIL  8.68s   （m1）
03-restored-green          rc=0  PASS  8.34s   （恢复后同 sha 再绿）
04-m2-edge-never-recorded  rc=1  FAIL 10.68s   （m2）
```

★ 诚实交代一条弯路：本脚本**第一版**用 `python` 调变异器，而本机 `python` 是 uv 的 trampoline，
在 bash 工具调用里起子进程失败（`uv trampoline failed to spawn Python child process`），
四发跑的都是**没变异的原码**（四发全绿）。这正是「绿的变异体最可怕」：所以第二版把解释器写成
`/c/Users/swq/.dsh/.../python/python.exe`，并让「变异没落地」直接 `exit 4` 结束整条链，
再跑才有上面这四行。

### 3.3 恢复＝同字节（sha256 读数，不是声明）

```
pristine: 1d47eccb205328a744fb1ff4507ef6962fa8cb29826671baf2c5ce9406178425  cmd/wisp/resident_ball_windows.go
pristine: d6d2a659903a8c20858249f5bf520058fcd2b0a28df3d291d9d3351ffeeaafee  cmd/wisp/resident_voice_windows.go
restored: ...同一对 sha256（restore 输出 pristine_same=True），导出树与主仓两份文件 sha256 相同
```
（日志 `logs/00-pristine-sha.txt`、`logs/03-restore-sha.txt`、`logs/05-restore-sha.txt`、`logs/05-restored-sha.txt`。）

---

## 四、AC#4 普查：`git grep` 出全部「读某一跳」的断言，逐条判同形，点名表

命令：`git grep -n 'timelineSnapshot()' -- 'cmd/wisp/*_test.go'`（23 行，含本票新增基线 1 行与 helper 内部 4 行）
与 `... 'voiceStatesSnapshot()' -- ...`（10 行），合计 33 行命中。
判据两条：**(a)** 读的那一刻，那一跳**能不能**还没落定（＝同形）；**(b)** 若不能，是哪条因果把它钉住的。

### 4.1 同形 —— 本票修掉的（5 处，见 §2.2）

| 位置 | 读点 | 判 |
|---|---|---|
| `resident_voice_interrupt_390r1_windows_test.go` 旧 `:107` | `timeline` 最后一枚 | **同形**（CI 红点）→ 修 |
| `resident_voice_interrupt_381r2_windows_test.go` 旧 `:349-353` | `timeline` 最后一枚 | **同形** → 修 |
| `resident_voice_wire_r1_windows_test.go` 旧 `:377-386` | `timeline` 恰好 4 枚 ＋ 球的投影 | **同形**（`ttsSpoken` 在 #26 派发之前）→ 修 |
| `resident_voice_acting_395r1_windows_test.go` 旧 `:159` | 记录快照（供断言） | **同形** → 修 |
| `voice_e2e_demo_16r1_windows_test.go` 旧 `:329` | 记录快照（供断言） | **同形** → 修 |

### 4.2 不同形 —— 逐条给了理由（同一批 `git grep` 命中，不改）

| 位置（改后行号） | 读的是什么 | 判 |
|---|---|---|
| `resident_voice_15e_windows_test.go:405` | 两枚记录 | **不同形**：#11 那一跳由 owner 在把转写交给文本路径**之前**派发，后面的等待（转写出现在 sink）依赖它 |
| `resident_voice_15e_windows_test.go:493-494` | 沉默用例 `len(timeline)==0` | **不同形**：**否定断言**（不该有边），不是「读某一跳」 |
| `resident_voice_15e_windows_test.go:546` | 球投影恰好 `[Listening]` | **不同形**：拒绝用例，其后 `asrLoaded`/`finals` 的等待依赖 #4 已发生 |
| `resident_voice_41a2r1_windows_test.go:162-163` | #5 走（失败信息里读） | **不同形**：读点在 `t.Fatalf` 的实参里（只在已经判败时求值），且其前有 `kwsLoaded` 等待 |
| `resident_voice_41a2r1_windows_test.go:377` | helper `voice41a2r1Walked` | **helper**，是否同形由调用点定：`:162`（上面）、`:421`（在等待谓词里）、`:587`（mute 手势同步走完）——三处都不同形 |
| `resident_voice_41a2r1_windows_test.go:453` | 全表日志 | **不同形**：只打日志 |
| `resident_voice_41a2r1_windows_test.go:432-441` | 球投影含 Armed/Listening | **不同形**：其前有 `asrLoaded`/`handOffs` 等待，ASR 被喂数据发生在 #4/#5 之后 |
| `resident_voice_interrupt_381r2_windows_test.go:212-213` | 空闲用例 `len==0` | **不同形**：否定断言；且手势同步、无在飞播报 |
| `resident_voice_interrupt_381r2_windows_test.go:273`、`:442` | 无 #27／无 #26 边 | **不同形**：否定断言，且该用例没有触发者 |
| `resident_voice_interrupt_381r2_windows_test.go:282` | 全表日志 | **不同形**：只打日志 |
| `resident_voice_interrupt_390r1_windows_test.go:322-323` | 未挂执行器的 `len==0` | **不同形**：否定断言，且无执行器可派发 |
| `resident_voice_wire_r1_windows_test.go:373` | 球的投影（等待谓词里） | **不同形**：谓词本身，本票新增的等待就是它 |
| `resident_voice_acting_395r1_windows_test.go:163` | 球投影日志 | **不同形**：只打日志 |
| `voice_e2e_demo_16r1_windows_test.go:339` | 球投影日志 | **不同形**：只打日志 |
| `voice_e2e_demo_16r1_windows_test.go:452` | helper `demo16Walked` | **helper**，调用点 `:306`（等待谓词里）、`:335/338/341/344/347`（本票修后都在等待之后） |
| `.../resident_voice_15e_windows_test.go:490`、`:497` | 否定断言／日志 | **不同形** |

### 4.3 普查结论

同形的**一共 5 处，全在 `cmd/wisp` 的 `*_test.go` 里**，本票全修（不是只修 CI 点名的这一枚）；
其它读点要么是否定断言、要么只打日志、要么被后面的因果钉住 —— 逐条写在 4.2，没有一条是「看起来像所以放过」。

---

## 五、AC#5 闸门

| 闸门 | 读数 |
|---|---|
| `go vet ./...` | **rc=0，无输出** |
| `gofumpt -l`（本票改的五枚） | **空**（rc=0） |
| `gofumpt -l internal cmd tools scripts frontend` | 命中一枚 `internal\models\acl_sid_121_test.go` —— **不是本票的文件**：工作树副本是 CRLF（268 LF / 268 CRLF），而 `HEAD` 的 blob 是 LF（268 LF / 0 CRLF）且 gofumpt-clean ⇒ Windows 检出伪影，与提交内容无关（§5.2） |
| `sh scripts/d22scan.sh` | **clean**：`no D22 ban violations`；正控 `PASS=45 FAIL=0 SKIP=0`，`examined 482 production Go files` |
| `sh scripts/check-path-length-budget.sh --with-self-test` | **VERDICT GREEN**：tracked=20086，over-budget=69，covered by roster=69，not in roster=0；正控 PASSED |
| 产码 **0 字节** | 本票 commit 只含 5 枚 `*_test.go` ＋ 本件 ＋ `.scratch/wisp/probes/445/**` ＋ 票 445 的 Progress log；`cmd/wisp` 的非测试文件一行未动（§3.3 的 sha256 同值即证） |
| CI `test-windows` **连续两绿** | **本腿做不到也不声称**：推由编排者做（⛔ push），两连绿是编排者的步骤；本腿给的是「同形五处全修 ＋ 仓内五例全绿 ＋ 导出树两枚变异的牙」 |

### 5.2 那枚 gofumpt 命中的完整读数

```
$ gofumpt -l internal cmd tools scripts frontend
internal\models\acl_sid_121_test.go
$ (bytes of the worktree copy)  bytes=12483 LF=268 CRLF=268
$ git show HEAD:internal/models/acl_sid_121_test.go > blob.go   (raw)
  blob bytes=12215 LF=268 CRLF=0 ;  gofumpt -l blob.go  ->  (empty)
$ git status --porcelain -- internal/models/acl_sid_121_test.go  ->  (empty)
```
⇒ 文件在 `HEAD` 的字节是 LF 且 gofumpt 干净（CI 在 Linux 上格式化的是 blob），本机工作树按 `text=auto` 检出成 CRLF。
本腿没有改它，也没有把它算进「我的五枚干净」以外的结论里。

---

## 六、发现了但不属于本腿写面的东西（报，不改）

1. **票面文件 `445-...md` 本身有字节损坏**：第 2978 字节处是一枚 `NUL`，另有若干字符被转义吃掉
   （`test-windows`→TAB+`est-windows`，`voice.interrupt`→`oice.interrupt`，`0xc0000135`→NUL+`xc0000135`，`add -A`→`dd -A`，`reset`→`eset`）。
   本腿按写面只往 Progress log 追加，**没有重写正文**。下一枚读这份票的人会看到那些错字，不是本腿造成的。
2. **2 枚前置事实**（供接手核对，不是缺陷）：`cmd/wisp` 在本波只有本腿一条写腿；本腿所有 `cmd/wisp` 测试
   都经 `scripts/wisp-cli-staged-tests.sh` 跑（同一时刻只跑一发，每次跑完把 `$WISP_CLI_STAGE_DIR/wisp-cli-tests.log` 复制/改名）。

---

## 七、复跑命令（逐条可现跑）

```sh
# 1) 五枚修后的仪器（主仓）
scripts/wisp-cli-staged-tests.sh -test.run '^(Test390r1BallClickPlayingInterruptAbortsAndWalks27|Test381r2PlayingPhaseInterruptAbortsAndWalks27|TestWireR1AnnounceEndToEnd|Test395r1VoiceToolTurnWalksActing|Test16r1VoiceDemoWakeAsrAgentSpeakSettleRealHardware)$'

# 2) 窗口读数（仓外导出树；探针文件在 .scratch/wisp/probes/445/）
cp .scratch/wisp/probes/445/probe-445-trigger-window_windows_test.go /d/tmp/wisp/20261114/445-tree/cmd/wisp/zz445probe_windows_test.go
cd /d/tmp/wisp/20261114/445-tree && WISP_445_ROUNDS=200 scripts/wisp-cli-staged-tests.sh -test.run '^Test445ProbeClickStopTriggerCalledAndWindow$'

# 3) 牙的两次变异（脚本自己 cp 修后文件进树、自己恢复并打 sha256）
sh .scratch/wisp/probes/445/ac3-445.sh /d/tmp/wisp/20261114/445-tree

# 4) 闸门
go vet ./... && gofumpt -l internal cmd tools scripts frontend && sh scripts/d22scan.sh && sh scripts/check-path-length-budget.sh --with-self-test
```

## 八、一句话

**球单击每次都真的调到了 `TriggerInterrupt`（260 发、260 次、被打断的句子一次没算成念完），
红的是那行断言「在走落定之前读了时间线」；修法是把「等标记」换成「等那一跳落定」，
五处同形全修，牙用仓外导出树的两枚真变异验过（m1 点名球→语音腿真漏，m2 让「睡久一点」永远绿不了）。**
