# `cmd/wisp` 装配腿 `wire-r1` —— 落地件（写腿）

> **尺**：`HEAD=cbb43c948`（`dev`），现跑 2026-10-26 19:0x–19:5x +08。★**十项读数逐条在 §1–§9**，全为现跑，⛔ 转述。
> **写面**：`cmd/wisp/**` —— **新 3 枚**（`resident_voice_engines_windows.go` 951 行／31,930 B、`voice_out_windows.go` 417 行／13,869 B、`resident_voice_wire_r1_windows_test.go` 794 行）＋ **改 3 枚**（`resident_voice_windows.go` 1022→1403 行、`run.go` +21、`resident_windows.go` +22/−?）＋本文件。
> **零字节**：`internal/**`（★`git status --porcelain -- internal/` ＝ **0 行**；含 `internal/speech`／`internal/session`／`internal/ball`／`internal/risk`）、`testdata/**`、`frontend/**`、`design/**`、`scripts/**`、`.github/**`、票面、台账、`docs/specs/**`。
> **临时件（只建不删）**：`D:\tmp\wisp\20261026\wirer1\{bin,speechbin,cases,logs,mut,overlay,fmt}`。★**三发变异全走仓外件 ＋ `-overlay`** ⇒ **工作树零突变窗口**。

## 0 一句话

三张票（26 TTS／27 标点／28 warm）卡的是**同一处**：`residentVoice` 旁边没有装配点。本腿把三件一次接完——**播报**（`risk.ChTTS` 闸 → `TtsEngine.Speak` → waveOut 出声 → 球 `Speaking` → Warm）、**标点**（`Transcriber.SetPunctuator` 的唯一生产调用点）、**warm 池**（`session.Scope` epoch ＋ `speech.Slot` 单族驻留；回合结束落 Warm，第二次 acquire **0s**）。

## 1 ★★★ 读数一：TTS 播报端到端一发（逐跳）

`TestWireR1AnnounceEndToEnd`（★**真 ASR ＋ 真 matcha TTS ＋ 真标点 ＋ 录音声卡**，一次现跑）：

| 跳 | 量 | 读数（逐字） |
|---|---|---|
| 触发 | 生产触发器 | `run.go` 的 D10 结果呈现处：`rt.voiceAnnounce(notifyBody(res))`（与通知**同一句**）→ 队列入 owner |
| ① 闸 | `risk.ChTTS` | `gate_ok=1 gate_refused=0`（★拦下时 `ttsGated=1` 且**不合成**，见 §4-M1） |
| ② 引擎 | 串行槽 | `slot=tts loads=2 closes=1` ⇒ 生命周期 **load:asr → close:asr → load:tts**（⛔ 与 ASR 共驻） |
| ③ Speak | C9 缝 | `chunks=1 samples=89611 rate=22050` |
| ④ 出声 | 渲染 | `peak=0.4686 rms=0.0702`（★若只断 `samples>0`，静音也能过）；**真设备同机现跑**：`DEVICE rate=22050 samples=4410 ms=200 abort=0`（440Hz 200ms 播完） |
| ⑤ 球 | D43 | 4 条边：`Sleeping+summon→Listening`(#4)、`Listening+vad-stop→Thinking`(#11)、`Thinking+first-token→Speaking`(#15)、`Speaking+speak-done→Warm`(#26) |
| ⑥ 落 Warm | session | `stage=Warm engines_hot=2 loads_after_warm=0 refs=0` |

★**acquire 读数**：`cold(turn1)=3.1097033s`／`announce(tts)=4.2876017s` ⇒ **换族是有价的**（D32 16.3.3），见 §9-③。

## 2 ★★★ 读数二：标点端到端一发（改前 vs 改后）

`TestWireR1PunctuationReachesTheAgent`（真 `Transcriber` ＋ 真 CT-Transformer，已知转写）：

| 面 | 文本（逐字） | `Punctuated` |
|---|---|---|
| **改前**（ASR 原始假设，`SetPunctuator` 未调用） | `今天下午三点提醒我开会` | `false` |
| **改后**（`adopt` 里 `tr.SetPunctuator(p)`） | **`今天下午三点提醒我开会。`** | **`true`**（`punctFinals=1`） |
| 进 agent 的那串（`voiceTranscriptSink`） | 与 final **同一串**（逐字比对） | —— |
| 词一个没动 | 剥掉 `。，？、` 后与原始串逐字相等 | —— |

★**唯一的调用点**：`resident_voice_windows.go` 的 `adopt`（本仓里唯一造 `Transcriber` 的地方）；★**partial 结构性不带标点**（`partial()` 一行未改，27-r1 的牙仍在）。
★**15e 的夹具**（`嗯嗯嗯嗯嗯嗯哦哦`）不加标点不是缺陷：那串是填充音，标点模型无事可做 ⇒ 本读数的对照用已知句。

## 3 ★★★ 读数三：warm 端到端一发（回合间）

`TestWireR1WarmSecondTurnAcquiresNothing`（真 ASR，两回合）：

| 量 | 读数（逐字） |
|---|---|
| 第一回合 acquire（冷装载） | **2.9790732s** |
| 回合结束 | `warm_turn_end=1` ⇒ `SCOPE-WARM reason=turn-end engines_hot=2`（★**⛔ 全关**） |
| 麦克风半边 | `mic_closes=2 stages_at_mic_close=[Awake Awake]` ⇒ **关在转移之前**（"closed FIRST" 是序） |
| 第二回合 acquire | **`0s`** |
| 装载次数（强判据） | `loader_calls=1`（两回合全程只装一次） |
| 会话自己的计数 | `loads_after_warm=0 warm_resumes=1`（`SCOPE-WARM-RESUME resumes=1 engines_hot=1`） |
| 两回合产出 | `finals=2 partials=10 handed=2 drained=325` |
| 90s→3s→dispose | `WARM-IDLE stage=Disposed engines_hot=0 engine_closes=1 mic_closes=2 stages=[Awake Disposed]` ⇒ **级联落在 owner goroutine**（见 §9-④） |

## 4 ★★ 反向判据：三读数 ＋ 还原核（时机＝**最后一发变异之后**）

| 发 | 变异（仓外 `mut\` ＋ `-overlay`，各只差一坨） | rc | 逐字红句 |
|---|---|---|---|
| **M1** | 摘掉 `ChTTS` 检查（`CheckText` 整块） | **1** | `timed out waiting for the gate to refuse the sentence` ⇒ 无闸 ⇒ 带毒句照播 |
| **M2** | 摘掉 `tr.SetPunctuator(p)` | **1** | `the final reached the agent UNPUNCTUATED ("今天下午三点提醒我开会") with a stage attached: SetPunctuator is not on the path` |
| **M3** | 回合结束改成 `closeEpoch`（＝"回合结束就全关"，28-r1 停手前的形状） | **1** | ①`stage after the turn = "Disposed", want Warm`；②`the ASR factory ran 2 times over two turns, want 1`；③**变慢**：`warm wake took its ASR borrow back acquire=2.4803697s`（绿 0s／冷 2.3–3.0s） |
| **还原** | 无 overlay 重编重跑 | **0** | 13 案全绿（§5） |
| **原样** | 再跑一发（同一二进制） | **0** | 13 案全绿，`262/640/902` 与 `cold_acquire` 同量级逐发复现 |

★★**还原核（19:49，三发变异之后）**：①树内 `MUTANT M[123] wire-r1` **0 命中**；②三枚变异件对工作树 `git diff --no-index --stat` ＝ **3／2／5 insertions**（只差变异块本身）；③`git status --porcelain -- cmd/wisp` 只列我 6 枚；④三个改件的 sha256 已记（`DDB47A50…`／`63D7479B…`／`4DFD982C…`）。

## 5 ★★ `15-e` 零回归（一人一发现跑）

| 尺 | 读数（逐字） |
|---|---|
| staging 6 案 | **6 PASS／0 FAIL**（端到端／静音正控／交还拒绝／无消费者正控／暂停／生产交还） |
| 延迟三数 | `LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` ⇒ **262/640/902 逐位同值** |
| 关键计数 | `drained=181/181 partials=5 finals=1 handed=1 handoffs=1 asrLoaded=true problems=[]` |
| 时间线 | 仍**恰 2 条**（`#4`／`#11`）⇒ 本腿没往 15e 的回合里塞边 |
| 另一发（19:45） | `drained=179/181`（15e-r1 记 178，28-r1 记 181）⇒ **末帧丢弃是本夹具的固有抖动**，非本腿回归 |
| 全量 | `TestWireR1|Test15e` ＝ **13 PASS／0 FAIL／0 SKIP**，rc=0 |

## 6 ★★ 正控：静音／空串／纯标点 ⇒ 都不出声

`TestWireR1SilenceAndEmptyTextNeverSpeak`：`POSITIVE-CONTROL finals=0 handed=0 announce_asked=4 empty=3 silent=1 spoken=0 samples=0`
- **120 帧数字静音**进 ⇒ `finals=0`、**没有**播报（`announceAsked=0` 在该段）⇒ "静音进 ⇒ 不播报、不调 Speak"。
- `""`／`"   "`／`"\t\n"` ⇒ 队前短路（`empty=3`），**零引擎装载**（`slot_loads=0`）。
- **纯标点** `。，、！？` ⇒ **`Speak` 被调、引擎自己拒绝**（`ErrNothingToSpeak`，`silent=1 samples=0`）⇒ ★**⛔ 空音频也播一句**（26-r1 的 M2 牙仍在）。

## 7 ★★ staging 跑法（完整命令）＋ 为什么 `PATH` 不够

```
cd D:\work\workspace\projects plans\Wisp
go test -c -o D:\tmp\wisp\20261026\wirer1\bin\wisp.test.exe ./cmd/wisp
copy third_party\sherpa-onnx\*.dll D:\tmp\wisp\20261026\wirer1\bin
set WISP_SILERO_MODEL=<repo>\third_party\spike-models\silero_vad.onnx
set WISP_SPIKE_MODELS=<repo>\third_party\spike-models
set WISP_PUNC_MODEL=<repo>\third_party\model-fetch\x-punc\<punc id>\model.int8.onnx
set WISP_TTS_DIR=<repo>\third_party\spike-models\matcha-icefall-zh-baker
set WISP_15E_TMP=D:\tmp\wisp\20261026\wirer1\cases
D:\tmp\wisp\20261026\wirer1\bin\wisp.test.exe "-test.run=TestWireR1|Test15e" "-test.v=true" "-test.timeout=30m"   =>  13/13, rc=0
```
★**为什么 `PATH` 不够（本机现跑，⛔ 引述）**：Windows 对**依赖 DLL** 的搜索序是「**主模块所在目录** → System32 → … → `PATH`」。`sherpa-onnx-c-api.dll` 从 `PATH` 找得到，但它依赖的 `onnxruntime.dll` 被 **System32 的 1.17.1** 先接走，而 linked C API 要 1.28。本腿**实测了第一层**：`go test ./internal/speech`（不 staging）＝ **rc 1，进程寄存器转储后 FAIL**（§8-②）；staging 后同一套 **55/0/0**。
★**`-test.run` 只跑定向**：整包会死在 `MARK299`；★**必须 `-timeout 30m`**（6 案里最长一条 ~12s，但模型装载在冷缓存上更久）。★**引号**：PowerShell 里 `-test.run=…` 不加引号会被拆成 `-test` ＋ 剩余 ⇒ 必须写成 `"-test.run=…"`。

## 8 ★ 四数

| 面 | 命令 | rc | 读数 |
|---|---|---|---|
| `cmd/wisp` **定向**（★带 `PATH`，⛔ 整包） | `go test -count=1 -run=TestAC3LogSinkLandsInsideTheEnvDataRoot -v ./cmd/wisp` | **0** | `--- PASS … 0.00s`；`ok … 0.067s` |
| `cmd/wisp` **本腿全量定向 staging** | §7 | **0** | **13 PASS／0 FAIL／0 SKIP** |
| `internal/speech`（**零字节**，独立跑） | `git status --porcelain -- internal/speech`＝**0 行**；staging：`go test -c` ＋ 三枚 DLL ＋ `WISP_*` | **0**／**0** | 未 staging ⇒ **rc 1**（缺 DLL，见 §7）；**staged ⇒ 55 PASS／0 FAIL／0 SKIP** |
| `internal/session`（**零字节**） | `git status --porcelain -- internal/session`＝**0 行**；`go test -count=1 ./internal/session` | **0** | `ok … 2.292s` |
| `go vet ./cmd/wisp`（含 `_test.go`） | —— | **0** | 空 |

## 9 ★ 格式甲乙两尺 ＋ 门禁 rc

| 尺 | 输出 |
|---|---|
| 甲 `gofmt -l internal cmd tools scripts frontend` | **恰 3 行**：`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`／`cmd\wisp\models.go` ⇒ **已知 3 条既存假红**（CRLF／`autocrlf` 造），**⛔ 本腿责任** |
| 乙 `D:\work\base\gopath\bin\gofumpt.exe -l`（同树） | **逐字同上 3 行** |
| **归属判据** | 两尺对我 **3 枚新件**单点跑 ＝ **甲空、乙空**；改件里只有 `models.go` 红（`crlf=389 lf=389`，其余 `.go` 全 LF ⇒ 我的改动不是 CRLF 造的） |
| 门禁 `d22scan`（`cd tools/d22scan && go run . -root <repo>`） | **先 rc 1／后 rc 0**：第一发 **2 finding**，★都是**我的真违规** —— `resident_voice_engines_windows.go:794,924` 的字符串里 **`⛔`(U+26D4)**（在 U+2600–U+27BF 带内）。改成「不念出来／不播一段近静音」后复跑 **rc 0 `clean - no D22 ban violations`** |
| ★**分母涨了** | `ban #8 cmd/` **155 → 158**（＋3 我的新件）、`bans #1-5 cmd/` **43 → 45**（＋2 生产件）；`internal/` 一面 **679／263 不变**（零字节） |
| ⛔ 错尺 | **没用 ban#8 去扫 `docs/`**（它只扫 `.go`／`design/`／`frontend/` 文本） |

## 10 未做／边界（逐条自陈）

1. ★★**播报闸的 taint 寿命**（最重要的一条）：`risk.Scope.Close` **会丢掉**该作用域的标记，而任务作用域在任务结束就关 ⇒ 结果句出生时"任务的毒"已经没了。本腿的闸读的是自己的进程级作用域 `voice-announce`（**已注册**，因此不会走 `scopeMarks` 的 unbound 分支）。★**它不是装饰**：M1 摘掉它当场红（§4），`TestWireR1AnnounceGateRefusesTaintedText` 里 `refused=1 passed=0 spoken=0 samples=0 slot_loads=0`（**合成之前**就拒）。★**要把任务的毒带到播报点**＝改 `internal/risk` 的作用域寿命（本腿禁写面）⇒ 具名残余。
2. ★**"关麦"关的是喂食门，不是设备**：C8 帧仍被排空、VAD 仍判定（否则 Warm 里永远听不到下一句）；把设备真关掉会让 warm 唤醒不可达，而唯一的另一唤醒源（KWS，票 41）**还没实现**。⇒ 证据里写的是 `mic_closes=2 stages=[Awake …]` 这个**序**，⛔ 冒充设备关闭。
3. ★**换族有价，逐条给数**：播报会把 ASR 的字节释放（`slot loads=2 closes=1`）⇒ 播报 acquire **3.8–4.3s**、下一回合 ASR 冷装载 **2.3–3.1s**。这是 D32 16.3.3 的 on-demand 档，⛔ 假装免费。另加一条**包内常量**闸：播报上限 `voiceAnnounceMaxRunes=120` 字（长文本走 D10 通知）。
4. ★**D43 的两处既存空档，⛔ 本腿拿表去填**：①`Warm + llm.first-token` 无行 ⇒ 从 Warm 播报时**球不动**（仍出声，行由 §1 的 Thinking 起点保证）；②`Thinking + asr.vad-stop` 无行 ⇒ 语音单回合后球停在 Thinking 而**会话已 Warm**（表没有 `Thinking→Warm`）。两条都是**表冻结**的直接后果，日志逐字留证。
5. ★**`internal/audio` 的 `DEFERRED(playback)` 由本腿在 `cmd/wisp` 侧回答**：`voice_out_windows.go` 用 **winmm `waveOut`**（非 WASAPI：不要 COM/STA、不要 Go→C 函数指针；22.05k mono int16 在它射程内）。★设备完成等待的预算是**按样本折算的音频时长 ＋ 固定宽限**，⛔ 墙钟差值。
6. ★**本腿不新增配置键**：TTS/标点路径**只读**模型库布局（`<store>/<manifest id>/<archive 内文件名>`），staging 用 `WISP_PUNC_MODEL`／`WISP_TTS_DIR`（与 15e 的 `WISP_SILERO_MODEL` 同一条既有先例）；`[session].warm_timeout_sec` **只读**。`[voice] punctuation` 开关仍需单开票（动 D36）。
7. ★**包内零裸 `go func(`**：owner 队列的排空在既有 `asr-infer` 线程的 select 里（D38b）；session 的定时器回调经 `voiceOwnerClock` **改投 owner**（级联关引擎绝不在 runtime 定时器线程上跑）。**零墙钟超时**：时长一律按样本计。两句都是 `d22scan` rc 0 ＋ 现跑 grep，不是声明。
8. ★**新依赖边 1 条**：`cmd/wisp → internal/session`（首次，C31 的 keep-alive 所有者）；`internal/speech`／`internal/risk` 的**既有语义零改**，`internal/ball` **零改**（只调用它已有的 `Ball.SetState`）⇒ 派单硬边界 1 因此不触发"停手上报"。
