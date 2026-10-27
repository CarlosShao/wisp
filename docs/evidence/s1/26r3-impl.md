# `cmd/wisp` 承接腿 `26-r3` —— `SetSpeaking` 一针接线 ＋ self-echo 夹具（落地件，写腿）

> **尺**：终 `HEAD=06609feaf`（`dev`），现跑 2026-10-27 17:5x–19:2x +08。★开工核：`git diff --cached --name-status`＝0 行；`cmd/wisp`／`internal/speech`／`internal/session`／`internal/audio` 四面 `git status --porcelain` 全 0 行。
> **写面**：`cmd/wisp/**` 4 枚（`resident_voice_windows.go`／`resident_voice_engines_windows.go`／`resident_windows.go` 改 ＋ `resident_voice_selfecho_26r3_windows_test.go` 新）＋本文件。**⛔ 零字节**：`internal/**`（含 `internal/audio`／`internal/risk`；`slot.go`/`kws.go` 是 `41-a2` 的在飞件，⛔ 当我的脏、零碰）、`testdata/**`、`frontend/**`、`design/**`、`scripts/**`、`.github/**`、票面、台账。
> **临时件（只建不删）**：`D:\tmp\wisp\20261026\26r3\{bin,speechbin,cases,logs,mut,overlay,fmt,fmt2,export}`。★三发变异全走**仓外 `git clone --shared` 导出树 ＋ `-overlay=<json>`** ⇒ 工作树零突变窗口（全程 `git status --porcelain -- cmd/wisp` 只在我的 4 枚）。

## 0 一句话 ＋ 一条路上抓的生产缺陷

一针接上：`announce()` 路径在换族之后、首个样本产出之前 `SetSpeaking(true)` 关采集门，`gateRelease`（defer，全函数唯一出口）保证**六条出口**（出声／闸拦／无引擎／无设备／合成失败／引擎拒读）都 `SetSpeaking(false)` 恰一次。★**顺带修复一件让整条播报链在生产上死着的前置缺陷**：`buildResidentVoice`（生产构造器）从不建 owner 队列 ⇒ `rv.queue=nil` ⇒ `announce()` 首行即拒、owner select 的 queue case 永不触发、session 定时器级联直接落在 runtime timer 线程——所有测试走 `newResidentVoice` 所以 staging 全绿而生产形状是死的。队列已在生产构造器补建（`resident_voice_windows.go` buildResidentVoice 注释具名）。

## 1 ★ 票面 7 格逐条（原文 ＋ 本腿覆盖判定）

| 格 | 票面原文（逐字摘） | 判 |
|---|---|---|
| AC1 | `Speak path: short result → TTS audible, ball Speaking with success breathing; interrupt mid-playback → instant stop + Listening.` | **⛔ 本腿不覆盖**：出声＋球 Speaking 已由 wire-r1 落地；缺的是**生产打断通道**（热键／球单击／否决词→#27）＋≤400ms 停播读数 ＝ tri-v1 `R26-1`，承接票。 |
| AC2 | `Half-duplex (Path T): playback injects no ASR transcripts (self-echo fixture); capture gate closed during playback (13 counter asserts zero frames).` | **★本腿覆盖主半**：`SetSpeaking` 生产接线（§2）＋self-echo 夹具（§3）：关门窗口 transcript 增量＝0、门态读数＝PathT speaking ⇒ 采集**拒**。tri-v1 `R26-2` 的"采集门未驱动"收口；"ring 重放含播放期回声"机理随门闭**结构性消除**（播放期 inner 源不再被交接采集，新帧根本不产生）。 |
| AC3 | `Path C barge-in (D47, blocked by 59)…` | **⛔ 不覆盖**：前置票 59 未清（tri-v1 `R26-3`）。 |
| AC4 | `Latency: first-audio ≤800ms warm / ≤1600ms cold-load measured (10 runs, P50/P95 recorded).` | **⛔ 本腿不测**：tri-v1 已带读数判"能翻（半格，冷档 NOT MET 必注）"；`R26-4` 归换族／驻留策略票。 |
| AC5 | `Serial slot: ASR and TTS never co-resident (sampler peak = max, asserted); unload frees.` | **⛔ 本腿不测**：tri-v1 已带读数判"能翻"（比值 0.41<0.75）。 |
| AC6 | `P7: score sheet committed ≥7/10 or blocked-decision logged for cloud-TTS fallback.` | **⛔ 本腿不碰**：第二支（`R5` 云 TTS blocked-decision）已入账，tri-v1 判"能翻"。 |
| AC7 | `C25 wiring… every announced string goes through risk.Provenance.CheckText… a hit means the announce is replaced by the L2 confirmation` | **⛔ 本腿不碰**：调用点已由 wire-r1 落地（M1 牙仍在）；"任务毒到不了 voice-announce"＋"命中无 L2 确认"＝tri-v1 `R26-5`，承接票。 |

## 2 ★★★ `SetSpeaking` 接线一发（含中止路径）

- **接在哪**：`speak()` HOP 2b（`resident_voice_engines_windows.go`）——HOP 1（ChTTS 闸）→ HOP 2（`acquireTurn(EngineTTS)` 换族，`asrGone` 置位）→ **HOP 2b `SetSpeaking(true)`** → HOP 3/4（合成＋waveOut）→ 出口全走 `gateRelease()` ⇒ `SetSpeaking(false)`。闸拦／无引擎两条在针**之前**返回，门从未闭 ⇒ 无需放（计数也不动，见正控）；针后五条出口（无 TTS 句柄／无设备／合成失败含 `ErrNothingToSpeak`／drain 失败／成功）全被 defer 覆盖。
- **为什么在换族后**：门闭之前的帧都已在 192ms 有界通道里（`internal/audio/gate.go` 注释原句：闭门时通道里的帧都是 TTS 开始**前**采的）——在首个样本产出前闭合才保得住该不变式。
- **staged 读数**（`Test26R3AnnounceClosesTheCaptureGate`，真 ASR＋真 TTS）：
  `GATE-26R3 closed=1 reopened=1 events=2 mid_speech_transcripts=0 spoken_samples=89407 gate_on=1 gate_off=1`；
  日志三行：`voice: capture gate wired to the announce path seam=audio.HalfDuplexGate.SetSpeaking path=T` → `voice: capture gate closed for the announcement reason=speaking path=T` → `voice: capture gate reopened (speech ended) reason=announce-done path=T`。
  探测在 `SetSpeaking` 调用瞬间读链态：**关门时 `asrGone=true asrLoaded=true`**（换族已发生、phase1 重挂未跑＝真 mid-speech）；句终落 `Warm`；REOPEN 后第二句真转写（`REOPEN-26R3 finals=2 warm_wake=1 stage=Warm`）。
- **观测面**：`announceReading()` 加 `门闭=N 放开=N`；`stopLine()` 加 `gate_closed=/gate_reopened=`；`setSpeakingCount()` 读数助手（nil 安全，`ZEROVALUE closed=0 reopened=0 stop_line_ok=1`）。

## 3 ★★★ self-echo 夹具一发（量的是门，⛔ 不说设备已关）

`Test26R3AnnounceClosesTheCaptureGate` 的 echo 段：**播报窗口内**向 C8 缝直推 6 帧（恰有界窗容量 192ms）⇒ 读数：`ECHO-26R3 gate_state_during_feed=speaking(PathT capture refused) frames_fed=6 drained_total=187 transcripts_after=1`（transcripts 相对播报前**零增量**）。★**门在哪个状态拒/放**：`PathT speaking ⇒ effectiveOpen=false（gate.go:234-242）⇒ 门**拒**采集交接（`Open()=false`）；本台架无设备，读数给的就是门态——生产上 `SetSpeaking(true)` 还会把 inner 源 `Stop()`（门不再把采集交接给内源，D16 语义），帧继续被消费侧排空而非堆积。★tri-v1 警句照守：全程写"门"与"采集交接"，⛔ 写"设备已关"。

## 4 ★★★ 反向判据两发 ＋ 还原核（时机＝最后一发变异之后）

| 发 | 变异（导出树 `D:\tmp\wisp\20261026\26r3\export` @`6c7733a42`，`-overlay` JSON 只换 `resident_voice_engines_windows.go`） | rc | 逐字红句 |
|---|---|---|---|
| **M1** | 摘掉 `SetSpeaking(true)` 块（`gateRelease` 与 false 保留） | **1** | `MUTANT 26r3-M1 SHAPE: the announce path never closed the capture gate (closed=0 reopened=0): SetSpeaking(true) is not on the path, so the D16 half-duplex close never happens and the microphone would hear the announcement` |
| **M2** | 摘掉 `gateRelease` 里的 `SetSpeaking(false)`（true 块保留） | **1** | `MUTANT 26r3-M2 SHAPE: the capture gate was left in Speaking after the sentence finished (closed=1 reopened=0 speaking=true): SetSpeaking(false) is not on the path; capture stays closed - stuck in Speaking - and no later utterance can ever be captured again`（★"卡在 Speaking"的可观测后果写进红句） |
| **还原核**（M2 之后） | 无 overlay 重编重跑 `Test26R3` | **0** | **2/2 PASS**：`GATE-26R3 closed=1 reopened=1 … spoken_samples=89172`／`ECHO…`／`REOPEN…` 全回绿；树内 `MUTANT 26r3-M[12] SHAPE` 字符串只存在于测试文件 4 行（具名红句锚点）；工作树 `git status --porcelain -- cmd/wisp`＝0 行 |
| **原样** | 终字节二进制再跑（与 15e 合跑，§6） | **0** | 16/16 PASS |

★三读数齐：变异红（M1/M2）／还原绿／原样绿。★变异与整包门禁**串行**（同一 shell 会话逐发跑，无并行）。

## 5 ★★ 正控：静音 ⇒ 不播报、`SetSpeaking` 不被调

`Test26R3SilenceNeverTouchesTheGateNeedle`：120 帧数字静音 ⇒ `SILENCE-26R3 finals=0 announce_asked=0 spoken=0 samples=0 gate_closed=0 gate_reopened=0`——门针（两套计数：recorder 与原子计数）零触发。⇒ "静音不播报"补上了"本该关门的那根针没动"这一半；与 wire-r1 既有 `TestWireR1SilenceAndEmptyTextNeverSpeak`（13/13 里复跑仍绿）互为正反。

## 6 ★★ `15e` 零回归（终字节二进制，staging）

| 尺 | 读数 |
|---|---|
| `TestWireR1|Test15e|Test26R3` 合跑 | **16 PASS／0 FAIL／0 SKIP**，rc 0（15e 6＋wire-r1 7＋本腿 3） |
| 延迟三数 | `LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` ⇒ **262/640/902 逐位同值** |
| 关键计数 | `drained=180/181`（另一发 179/181；**178/179/181 都是本夹具固有抖动，⛔ 当判据**）、`partials=5 finals=1 handed=1 handoffs=1 asrLoaded=true problems=[]` |
| 时间线 | 仍恰 2 条（#4／#11）⇒ 本腿没往 15e 回合里塞边 |
| 静音正控（15e 侧） | `SILENCE drained=120 handoffs=0 finals=0 ball=[] state=Sleeping` |
| 针零回归 | wire-r1 七案 `gate_closed=0 gate_reopened=0`（未接门缝的台架针不动）＋ `TestWireR1AnnounceEndToEnd` 播报读数同型（samples≈89.8k、slot loads=2 closes=1、落 Warm） |

## 7 ★★ staging 跑法（完整命令）＋ 为什么 `PATH` 不够（三层口径）

```
cd "D:\work\workspace\projects plans\Wisp"
$env:GOFLAGS=''
go test -c -o D:\tmp\wisp\20261026\26r3\bin\wisp.test.exe ./cmd/wisp
copy third_party\sherpa-onnx\*.dll D:\tmp\wisp\20261026\26r3\bin\
$env:WISP_SILERO_MODEL="<repo>\third_party\spike-models\silero_vad.onnx"
$env:WISP_SPIKE_MODELS="<repo>\third_party\spike-models"
$env:WISP_PUNC_MODEL="<repo>\third_party\model-fetch\x-punc\<punc 目录>\model.int8.onnx"
$env:WISP_TTS_DIR="<repo>\third_party\spike-models\matcha-icefall-zh-baker"
$env:WISP_15E_TMP="D:\tmp\wisp\20261026\26r3\cases"
D:\tmp\wisp\20261026\26r3\bin\wisp.test.exe "-test.run=Test26R3" "-test.v=true" "-test.timeout=30m"   ⇒ 3/3, rc=0
D:\tmp\wisp\20261026\26r3\bin\wisp.test.exe "-test.run=TestWireR1|Test15e|Test26R3" "-test.v=true" "-test.timeout=30m"  ⇒ 16/16, rc=0
```
①Windows 对依赖 DLL 的搜索序是「主模块所在目录 → System32 → … → `PATH`」：`sherpa-onnx-c-api.dll` 从 `PATH` 找得到，但它依赖的 `onnxruntime.dll` 被 **System32 的 1.17.1** 先接走而 C API 要 1.28 ⇒ 不 staging 就是 `0xc0000135` 进程起不来（tri-v1 §7 同口径现证，⛔ 当测试红）。②所以 DLL 必须**拷到测试二进制旁边**（`copy` 那行），`PATH` 替代不了这一层。③模型文件走 `WISP_*` 环境变量（15e 既有先例，非配置键）。★`-test.run` 只跑定向（整包会死在 `MARK299`）、必抬 `-timeout 30m`、PowerShell 里 `-test.run=…` 必须带引号。

## 8 ★ 四数

| 面 | 命令 | rc | 读数 |
|---|---|---|---|
| `cmd/wisp` 定向（带 PATH、staged、只 `-run`） | `wisp.test.exe "-test.run=TestAC3LogSinkLandsInsideTheEnvDataRoot" -v -timeout=30m` | **0** | `--- PASS … 0.00s` |
| `internal/speech`（staging，`-c`＋三 DLL＋`WISP_*`） | `speech.test.exe -test.v=true -test.timeout=30m` | **0** | **67 PASS／0 FAIL／0 SKIP**（含 `41-a2` 合入后的 KWS 面；开工时该目录有其在飞脏，收工时 `git status` 0 行） |
| `internal/session`（零字节） | `git status --porcelain -- internal/session`＝0 行；`go test -count=1 ./internal/session` | **0** | `ok … 2.417s`（终 HEAD 复跑 `2.815s`） |
| `go vet ./cmd/wisp` | —— | **0** | 空（途中一次 rc1 是 `41-a2` 在飞窗口的 `kws.go:142`，重试即 0，⛔ 归它） |

## 9 ★ 格式甲乙两尺 ＋ `d22scan`（导出树）＋ 门禁 rc

| 尺 | 工作树 | 导出树（同 `HEAD`） |
|---|---|---|
| 甲 `gofmt -l internal cmd tools scripts frontend` | **恰 3 行**：`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`／`cmd\wisp\models.go` | **空、rc 0** |
| 乙 `gofumpt -l`（`D:\work\base\gopath\bin\gofumpt.exe`） | **逐字同上 3 行** | **空、rc 0** |
| 逐件终证 | `git show HEAD:<file>` 原字节 blob `CR=0/0/0`、worktree `CR=131/225/389`；blob 拷出跑 `gofmt`/`gofumpt` ⇒ **不点名** | ⇒ 3 条是 `autocrlf` 显示面假红，**⛔ 本腿责任**（tri-v1 §8 同法复证） |
| 本腿 4 枚单点 | 两尺皆空（`gofumpt -l cmd\wisp` 只剩已知 `models.go`） | —— |

`d22scan`（`git clone --shared` 导出树 @`06609feaf`，`git status` 0 行；`cd tools/d22scan && go run . -root <export>`）：干净树 **rc 0 `clean`**，分母 `bans#1-5 internal/=264、cmd/=45、ban#8 internal/=684、cmd/=159`。★途中真抓我一条：`resident_voice_windows.go:311` 注释 `cmd/wisp-side` 被 **ban #9 幻影引证**点名（rc 1）⇒ 改 "on the cmd/wisp side of the seam" 追加笔 `06609feaf` 后 rc 0。**正控两发**（注入导出树、验后即移）：A 字符串放 `U+2713` ⇒ rc 1 `[emoji] ban #8 …string literals included`；B 裸 `go func(` ⇒ rc 1 `[bare-goroutine] …use observe.Registry.Spawn`；删两枚复跑 ⇒ rc 0 `clean`（分母回落）＝正控可逆（第 21 处定式）。★⛔ 未拿 ban#8 扫 `docs/`。

## 10 提交号 ＋ 交回

| 笔 | commit | 面 |
|---|---|---|
| 接线＋夹具 | `ad309e0ac` | `M` 3 `.go` ＋ `A` 1 `.go`（4 文件 ＋384/−5） |
| 反向判据红句具名化 | `6c7733a42` | `M` 1（测试件 ＋16/−6） |
| 门针正控案 | `003203ab5` | `M` 1（测试件 ＋54） |
| 幻影引证修正 | `06609feaf` | `M` 1（注释单行 ＋2/−1） |

`add`/`commit` 两处**显式逐文件 pathspec**；`-F` 真临时文件、写在 `--` 之前；主题逐字以 `票 26 承接腿 26-r3` 开头；⛔ push。四笔全在 `cmd/wisp/**`，与 `41-a2`（`d8dcf4c31`）及编排者收件笔（`75900bf9b`）**追加交错**、零改写。★**未定义即停：无停手项**——两件都做完；唯一边界：AC2 的"13 counter zero frames"半句以门态读数等价收口（本台架无设备，设备级零帧归生产装配态）。
