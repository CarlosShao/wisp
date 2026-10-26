# 票 15e 装配与端到端腿 `15e-r1` —— 落地件（写腿）

> **尺**：`HEAD=12be96ad3`（`dev`），现跑 2026-10-26 16:57–17:35 +08。★**十项读数逐条在 §1–§10**，全为现跑，⛔ 转述。
> **写面**：`cmd/wisp/**` —— **改 4 枚**（`models.go` 注释、`resident_audio_windows.go`、`resident_ball_windows.go`、`resident_windows.go`）＋ **新 2 枚**（`resident_voice_windows.go` 954 行／40,678 B、`resident_voice_15e_windows_test.go` 670 行／28,686 B），＋本文件。
> **零字节**：`internal/**`（含 `internal/speech`／`internal/audio`／`internal/models`）、`testdata/**`、`frontend/**`、`design/**`、`scripts/**`、`.github/**`、票面、台账、`docs/specs/**`。
> **临时件（只建不删）**：`D:\tmp\wisp\20261026\15er1\{bin,cases,logs,mut,overlay}`。★**变异走仓外件 ＋ `-overlay`** ⇒ **工作树零突变窗口**。

## 0. 一句话

`15-a1` 查清的**三处 `cmd/wisp` 落点**全部接上，并让 `internal/speech` 第一次有生产 importer：`resident_voice_windows.go` 里 **①** `models.go:25` 预留的 `handOffModel` 调用点（`residentVoiceHandOff`，**首次开句**时才走）**②** 排空 C8 帧通道的那**一个**消费者（`asr-infer`，D38b 临时名册 ＋ C11 作用域）**③** 托盘「暂停唤醒」的执行器与勾选真相源。链路本身用 `15-b` 的判停 ＋ `15-c` 的 `Transcriber`（**未改一行语义**），球走 **D43#4 → D43#11**（`Sleeping→Listening→Thinking`）。

## 1. ★★★ 端到端一发：音频进 ⇒ 判停 ⇒ ASR ⇒ `handOffModel` 被调用

`Test15eSpeechEndToEndHandsTheTranscriptToTheAgent`（staging，见 §8）。夹具＝**32ms 帧 ×181**：1 帧暖机静音 ＋ 250ms 静音 ＋ 3s 合成语音（120Hz 谐波栈×8／3.5Hz 包络，15-b/15-c 的同一支夹具）＋ 2.5s 静音，经**真 C8 缝**（`audio.WavInjector`）进 `audio.NewBoundedFrames()`。

| 读数 | 值 |
|---|---|
| **钩子被调了几次** | **1 次**（`handOffs=1`；⛔ 每句一次、⛔ 启动时一次） |
| **参数是什么** | `id="asr-streaming-paraformer-zh-en"`（签名清单 purpose=`asr-streaming`）＋ `logf != nil`（进度 tick 进日志） |
| **交还走的哪条路** | **生产那条**：`residentVoiceHandOff` → `loadModelConfig` → `openModelStore`（`build/models/manifest.json` + `.minisig` 真验签）→ `store.handOffModel` → `DownloadingBridge.Run` → `Ensure`（缓存命中＋哈希复验）→ `VerifyInstalled` |
| **交还自己的话（原样）** | `wisp models ensure: asr-streaming-paraformer-zh-en 已交还（state=FirstRun，ticks=1，末次进度 100%，2.31s，store=…\models）` ⇒ **exit=0** |
| **字节什么时候打开** | **交还返回 0 之后**（同一枚 `model-downloader` 线程里），路径＝`<data dir>\models\asr-streaming-paraformer-zh-en\{encoder.int8.onnx,decoder.int8.onnx,tokens.txt}`（清单文件名，⛔ `third_party/asr-encoder.int8.onnx` 那套 dev 名） |
| **转写** | `partials=5`、`finals=1`、文本 **`嗯嗯嗯嗯嗯嗯嗯哦哦`**、`closedBy=end`、`speech=3392ms` |
| **进 agent（票 12 文本路）** | `handed=1`，`residentTaskSource.submitTask(text)` 收到与 final **同一串**（test 的 sink 逐字比对） |
| **帧** | `drained=178/181`（丢的 3 帧＝回放解码那 ~0.5s 里最末尾的帧，D38d 计量；见 §11 残余） |
| **problems** | **空**（没有一句降级话） |

★**"函数没报错"不算数**：上表里 `exit=0` 是**清单签名验过 ＋ 三枚文件哈希复验过**的产物；§6 的发 A 证明删掉这次调用后**模型字节根本不存在**。

## 2. ★★ 球走态：`Listening → Thinking` 真的发生了（时序）

`residentBall.setVoiceState` 是本腿唯一的投影跳，它**先记账再投**（`rb.b != nil` 时同一句里 `Ball.SetState`，走的是批准门 UI 早就在用的那条 `PostTask`）。时序全在**音频钟**（已排空样本数），⛔ 墙钟（D42#9）：

| # | 转移 | 事件 | 音频钟位置 | 触发（逐字） |
|---|---|---|---|---|
| 1 | `Sleeping → Listening` | `session.summon`（D43#4） | **13,824（864ms）** | 麦克风开着时的第一次语音（唤醒词＝票 41 未落、召唤热键无执行器 ⇒ D43 里进 Listening 的**唯一合法行**就是 #4） |
| 2 | `Listening → Thinking` | `asr.vad-stop`（D43#11） | **89,600（5.6s）** | 判停（`end`）且语音 **3392ms ≥ 300ms**（守卫 `speech-ge-300ms` 用的是 `Endpoint.Speech`，不是墙钟） |

★**读数**：`ball.voiceStatesSnapshot() = [Listening, Thinking]`（顺序断言）；`timeline` 恰 2 条，`From/Event/To` 逐条对表。★**没改 D43 一行**：表在 `internal/statemachine/table.go`，本片只 `Dispatch` 合法行；不合法就 `noteProblem` 记账，⛔ panic。★**残余**：`Listening` 自己的超时（D43#12）由 `statemachine` 的定时器推动，**没有投影钩子** ⇒ 那一条机器走了球不跟（证据留在 `run`/`dispatch` 的注释里，⛔ 假装没有）。

## 3. ★★ C8 帧通道：`nothing drains it` 那个洞 —— **之前 0 个消费者、现在 1 个**

| 尺 | 命令 | 读数 |
|---|---|---|
| **之前（HEAD）有几个读者** | `git grep -n 'NewBoundedFrames' HEAD -- cmd internal` | 生产面**只有分配点**：`cmd/wisp/resident_audio_windows.go:276 buf := audio.NewBoundedFrames()` |
| **之前有没有人收** | `git grep -n 'ra.frames\|rv.frames\|<-buf\|<-frames' HEAD -- cmd/wisp` | **0 命中** |
| **现在** | 同上（工作树） | 存：`resident_audio_windows.go:297 ra.frames = buf`；收：`resident_voice_windows.go:462`（phase 1）／`:489`（接管后）／`:691`（暂停时空排） ⇒ **1 个消费者** |
| **仪器正控（没有人收时真的会丢）** | `Test15eTheChannelHasNoConsumerWithoutTheLeg`：同一个 `NewBoundedFrames()` 塞 20 帧 | **capacity=6，accept=6，refuse=14，窗口=192ms** ⇒ 缝确实是"没人抽就丢"的形状 |
| **动态证明（有人收）** | 同 §1 | `drained=178/181`；纯静音那一发 `drained=120/120` |

★**注释同步改**：`resident_audio_windows.go` 头注释与 `posture()` 里"帧消费侧（ASR/KWS，`internal/speech` 一块没写）不属于本票"**已改成读接线**（`residentVoice.drainsFrames()`）；`posture()` 现在会说 `C8 通道 1 个消费者（D38b 名册名 asr-infer，C11 作用域）`。

## 4. ★★ 唤醒监听暂停接上（`resident_ball_windows.go:522,596`）

| 面 | 改法 | 读数 |
|---|---|---|
| 手势执行器 | `OnTrayPauseWake` 从 `recordBallGesture("tray-pause-wake")` → `rb.pauseWakeGesture()`；`attachWakePause(rv.togglePause, rv.pausedForTray)` 由 `startResidentVoice` 在语音腿建好后注入 | 有腿 ⇒ 说"已暂停唤醒监听…／已恢复唤醒监听…"；无腿 ⇒ **仍旧**记 unhosted（同一句 `ballGestureWhy`，只加"ticket 15e"一项） |
| 勾选真相源 | `mirrorTrayMute` 的**第二个 bit** 从字面 `false` → `wakeRead()`（有读者才写，没读者仍是 `false`） | `Test15ePauseStopsTheListenerAndTheTrayFollows`：`drained=3 paused_drops=4 fed_before_pause=2 tray_writes=2` ⇒ 暂停时帧**照旧被排空但不再进 VAD/ASR**，勾选跟着变 |
| 状态机 | 暂停且 `Listening` ⇒ D43#13（`EvVeto`，冻结副作用 `audio.discard-buffer`）；恢复且 `Warm` ⇒ D43#29 | 只用表里的行 |
| 293 的四条旧断言 | **仍绿**（没接读者时 `paused=false` 那条路原样） | 见 §9 定向跑 |

★**诚实边界**：暂停不是锁 —— 按下那一瞬间已进 `select` 的一帧还会被喂（test 显式等"观察到暂停"再断言，⛔ 假装瞬时）。

## 5. ★★★ 延迟两段（实测）＋ 球延迟合成

| 段 | 定义（音频钟） | 本次实测 | 15-c 记的 |
|---|---|---|---|
| **① 开句迟**（silero 判"这是语音"的迟） | 本句首样本 − 真语音起点 | **262ms** | 1088ms |
| **② chunk**（Paraformer 出一段文本的粒度） | 首个 partial 的 `Transcript.Latency`（本句首样本 → 出字窗） | **640ms** | 640ms |
| **球延迟合成**（语音起 → 第一段文本） | ①＋② | **902ms ＝ chunk 的 1.41×** —— ★**只算 ② 会低估 41%** | 1728ms（1.7×） |

★**两段都不是常量**：① 是 silero 在这条夹具上的分类迟（本夹具前面有 250ms 数字静音 ＋ 1 帧暖机，判据是从干净静音起步；15-c 那条从第 0 样本就是语音、且是反推出来的），② 是模型给的 chunk 粒度。★★**同一支 `Transcriber` 复现了 15-c 的 640ms，逐位相同**（两次复跑都是 640ms），这是两条腿之间唯一可比的数。
★**另给两个不算进球延迟但必须知道的数**：`seamWaitToFirstText=5.22s`（本次第一个字真正出现在缝上的位置 ⇒ 首用成本：交还 2.31s ＋ 引擎装载 ~2.2s，**每进程一次**）；`finalLatency=4.096s`（判停句全长）。

## 6. ★★ 反向判据：三读数 ＋ 还原核（**时机在最后一发变异之后**）

| 发 | 变异（仓外 `mut\` ＋ `-overlay`） | rc | 红的断言（逐字） |
|---|---|---|---|
| **A** | `resident_voice_windows.go`：把 `exit := hook(voiceAsrModelID, rv.logf)` 换成 `exit := 0`（**交还调用被摘掉**） | **1** | **产品级红**：`voice: ASR 没起来（交还已通过，但字节打不开）：speech: sherpa refused the streaming ASR config (…): a file is unreadable` ＋ `resident_voice_15e_windows_test.go:359: timed out waiting for the utterance to be transcribed and handed over` ⇒ **没有这次调用，`<store>/<id>/` 里就没有字节**（夹具把"装模型"也塞进交还缝：生产里是 `Ensure` 下载，测试里是链接钉住的副本） |
| **B** | **还原**（无 overlay，重编重跑） | **0** | **6 PASS / 0 FAIL**（§9） |
| **C** | **原样绿**（同一发，重跑第二次） | **0** | 6 PASS；`openDelay=262ms／chunk=640ms／ballLatency=902ms` **逐位同值** ⇒ 时序是确定性的，不是抖出来的 |

★**还原核**：①`git hash` 前后同值 —— 工作树 `resident_voice_windows.go` sha256 `60431fdf…c90ca8146e`（变异前）＝ 变异后 `60431fdf…c90ca8146e`；②树内 `MUTANT` 标记 **0 命中**（`cmd/wisp/*.go`）；③变异件只在 `D:\tmp\…\mut\`，`-overlay` 只改编译输入 ⇒ **工作树零突变窗口**。

## 7. ★★ 正控：静音进 ⇒ **不调** `handOffModel`（⛔ 空音频也调一次）

`Test15eSilenceNeverCallsTheHandOff`：**120 帧数字静音**（3.84s）逐帧推入、逐帧等排空。

| 量 | 读数 |
|---|---|
| `handOffs` | **0**（腿里 0 次、缝上 0 次） |
| `asrLoaded` | **false**（ASR 字节没被打开） |
| `finals` / 文本路 | **0 / 空** |
| 球投影 | **`[]`**，机器 `Sleeping`，timeline **空** |
| `drained` | **120/120** ⇒ 绿不是"什么都没跑" |

★**为什么这样才对**：手交挂在 **`EndpointStart`** 上（SPEC-04 §7.1「首次使用某能力」），而静音**开不了句**（15-b 的状态机）⇒ 静音永远走不到能力。⛔ 启动时就调（那就等于"空音频也调一次"）。

## 8. ★★ staging 跑法（完整命令）＋ **为什么 `PATH` 不够**

```
# 1) 编译测试二进制（三枚 sherpa DLL 必须与它同目录）
cd D:\work\workspace\projects plans\Wisp
go test -c -o D:\tmp\wisp\20261026\15er1\bin\wisp.test.exe ./cmd/wisp
copy third_party\sherpa-onnx\*.dll D:\tmp\wisp\20261026\15er1\bin
# 2) 两个模型环境变量（staging 目录里没有 third_party）
set WISP_SILERO_MODEL=<repo>\third_party\spike-models\silero_vad.onnx
set WISP_SPIKE_MODELS=<repo>\third_party\spike-models
set WISP_15E_TMP=D:\tmp\wisp\20261026\15er1\cases
# 3) 就地在仓根跑（二进制在 bin\，映像目录优先，cwd 留在仓里给 build\models\manifest.json）
D:\tmp\wisp\20261026\15er1\bin\wisp.test.exe -test.run=Test15e -test.v=true   =>  6/6, rc=0
```
★**为什么 `PATH` 不够（本片实测，⛔ 引述）**：把 `third_party\sherpa-onnx` 塞进 `PATH`、**不做 staging**，本腿的第一条用例 **0.00s 就红**，逐字：`语音链未构造：VAD 没起来（speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH (on Windows the loader searches System32 before PATH)）`。原因＝Windows 对**依赖 DLL** 的搜索序是「**主模块目录** → System32 → … → `PATH`」：`sherpa-onnx-c-api.dll` 能从 `PATH` 找到，但它依赖的 `onnxruntime.dll` 被 **System32 的 1.17.1** 先接走。★三层口径（本片复述并自测了第 2/3 层）：**什么都不给 ⇒ `0xc0000135`（进程起不来）**；**给 `PATH` ⇒ 起得来但拿 1.17**（本腿：具名拒绝，⛔ 崩）；**staging ⇒ 拿 1.28.2，才全绿**。★`15-b` 的"8 过/8 skip"与 `15-c` 的"10 过/18 skip"口径不同（包长大了），⛔ 矛盾。

## 9. ★ 四数

| 面 | 命令 | rc | 读数 |
|---|---|---|---|
| `cmd/wisp` **定向＋`PATH`**（⛔ 整包：会死在 `MARK299`） | `go test -count=1 -run TestAC3LogSinkLandsInsideTheEnvDataRoot -v ./cmd/wisp` | **0** | `--- PASS`；`ok … 0.070s`（共享编译面没被我改坏） |
| `cmd/wisp` **15e 定向 staging** | §8 | **0** | **6 案 6 过**：端到端／静音正控／交还拒绝／无消费者正控／暂停／生产交还 |
| 同上一发的关键计数 | —— | —— | `handoffs=1 finals=1 handed=1 partials=5 drained=178/181`；`paused_drops=4`；`refuse=14/20`（§3 正控） |
| `internal/speech`（**零字节改动**） | `git status --porcelain -- internal/speech` ＝ **0 行**（HEAD 19 文件）；staging 跑自家全套 | **0** | **28 PASS / 0 FAIL / 0 SKIP**（＝ 15-c 记的 28） |
| `internal/audio`（**零字节改动**） | `git status --porcelain -- internal/audio` ＝ **0 行**；`go test -count=1 ./internal/audio` | **0** | `ok … 17.138s` |

## 10. ★ 格式甲乙两尺 ＋ 门禁

| 尺 | 命令 | 输出 |
|---|---|---|
| 甲 `gofmt` | `gofmt -l cmd/wisp` | 只有 `cmd\wisp\models.go`（**已知 3 条既存假红之一**，CRLF／`autocrlf` 造的；把该文件转 LF 后 `gofmt -l` **空** ⇒ 我的注释改动格式干净） |
| 乙 `gofumpt`（绝对路径） | `D:\work\base\gopath\bin\gofumpt.exe -l cmd/wisp` | 同上，只有 `models.go` |
| **归属判据** | 两尺对本片 **2 枚新文件**单点跑 | **甲空、乙空** |
| `go vet ./cmd/wisp` | —— | **rc 0**（含 `_test.go`） |
| `d22scan` | `go run . -root <repo>`（`tools/d22scan`） | **rc 0** `clean - no D22 ban violations`；`cmd/` 155 枚 Go 全覆盖、bans #1-5 cmd/=43 |
| ★**它抓到了我一条真违规** | `ban #8 cmd/wisp/resident_voice_windows.go:657` | 字符串里的 **`≥`（U+2265）** —— 改成「不小于 300ms」后复跑干净。★**ban#8 的射程是对的**：注释豁免、**字符串不豁免**（`Q-46(c)`）；同一处 `→`（U+2192）虽在仪器的空隙里，也一并改成「到」 |
| ⛔ 错尺 | —— | **没用 ban#8 去扫 `docs/`**（它只扫 `internal/`＋`cmd/` 的 `.go`） |

## 11. 未做／边界（⛔ 越界即作废，逐条自陈）

- **AC5 的"一次重试 → `asr` 错误 →『没听清』"策略未做**：空 final 现在会被记成「这一句没听清」并**不**交给文本路（AC5 键的那条语义在），但**重试与错误类映射**不在本片。
- **交还被拒 ⇒ 本进程的语音链到此为止**（不重试，等下次启动）；重试策略不是本片。
- **D43#12 的 Listening 超时**：机器走、球不跟（`internal/statemachine` 的定时器没有给外部钩子）—— 已在代码注释与 §2 记明。
- **首用成本**：交还 2.31s ＋ 引擎装载 ~2.2s，全部落在**第一句**上（`seamWaitToFirstText=5.22s`）；⛔ 隐藏。回放环 8s：**首次下载**（分钟级）超过它 ⇒ 触发那一句会丢，本片如实说而不假装。
- **切缝本身丢最末尾 3/181 帧**（回放同步解码 ~0.5s，D38d 计量）；出声的句子**不受影响**（环里已有它）。
- **⛔ 加配置键**：模型 id、环长、chunk 预算全是**包内常量**；`[voice].asr.model` 只**读**（不新增键）。
- **未做**：KWS（票 41）、TTS（26）、AC4 的 slot 互斥**实测**（要 TTS 才凑得出）、真语料 CER（15-d）。
- 包内**零裸 `go func(`**（`scope.Go` 两处：`asr-infer`／`model-downloader`，都在 D38b 临时名册里）、**零墙钟超时**（时长一律按样本计）—— 两句都是 `d22scan` rc 0 ＋ 现跑 grep，不是声明。
- **新依赖边一条**：`cmd/wisp → internal/speech`（首次）；`internal/speech` 仍**不 import `cmd/wisp`**，`internal/audio` 不看 `internal/ball`。
