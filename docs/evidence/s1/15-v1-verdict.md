# 票 15 终裁腿 `15-v1` —— 五片终裁（非实现者）
> **尺**：`HEAD=de02d5053`（写作时父提交；`dev`，2026-10-26 18:2x +08；五笔代码提交全为祖先。读的过程中 HEAD 从 `a148fda16` 走到 `de02d5053`，后十一笔是基线集/CI/证据/台账；`git diff 2c0be38dd..de02d5053 -- internal cmd go.mod` **空**）。
> **写面**：本文件 ＋ `D:\tmp\wisp\20261026\15v1\**`。**落点前 `ls`：`docs/evidence/s1/15-v1-verdict.md` 不存在**（18:14）⇒ 用本号、不拆 `-partN`。
> **方法**：五笔提交逐个 `git show --name-status`；测试全在 staging（三枚 DLL 同目录＋`WISP_SILERO_MODEL`/`WISP_SPIKE_MODELS`）跑；突变全在仓外 `mut\` ＋ `-overlay`；还原核在**最后一发变异之后**（§3、§5、§6 各一次）。
> **⛔ 未改**：`internal/**`、`cmd/**`、`testdata/**`、`.scratch/wisp/**`、`docs/reports/**`、`docs/specs/**`、`docs/PLAN.md`、`docs/audit/**`；**⛔ push**。工作树里别人的脏（`testdata/asr-baseline/**` 基线集腿、`.scratch/**`）**零碰零清**。
## 0 一句话判语
五片**都够格落地**（提交边界干净、零回归双向成立、三条核心证明都能自造复现）；但票 15 **不是六条全翻**：**AC1／AC3／AC6 可翻**（各带一条证据边界），**AC2／AC4／AC5 是具名残余**。`15-a` 的参数类型放宽**我判对**。唯一曾红的现行面是基线集腿先交的 **10 列 `sources.tsv`**（17:52 实测 `cer` 23 PASS／1 FAIL／0 SKIP），**已被 15d3/15d4 闭环**（18:0x 起 24 PASS／0 SKIP／0 FAIL）；AC2 仍缺“真 ASR 跑分”（CI 接线已由 `cerci-r1` 接，但只证 gate、不证 6%/15%）。
## 1 票面前提新鲜度（6 条 AC 全未勾）
| AC | 票面要点（逐字节选） | 现在状态 | 覆盖它的片 |
|---|---|---|---|
| AC1 | wav→VAD 判停→ transcript 进 agent（12 的 text path）端到端；球 `Listening→Thinking` | **可翻** | `15-b`＋`15-c`＋`15-e`；`15-e` §1/§2 是唯一整链证据 |
| AC2 | `near_clean`≤6% 且 `noisy_far`≤15%（提交集；CI-emitted report） | **⛔ 翻** | `15-d`（计算器/契约/gate）＋基线集 `15d2-r1`→`15d3-r1`→`15d4-r1`；**无真 ASR 跑分**；CI 已接但只证 gate |
| AC3 | 说-停-说切分；<300ms 忽略；60s 自动判停 | **可翻** | `15-b` §1/§2（状态机＋真模型）；`15-c` 的 speak-pause 两 final |
| AC4 | slot 互斥：ASR+TTS 串行驻留、卸载 RSS delta、切换延迟 vs 1.6s | **⛔ 翻** | 只有 `15-a` 的 `TtsEngine` 接口位；**无 slot、无 TTS（票 26 0/7）** |
| AC5 | 纯静音→一次重试→`asr` 错误＋「没听清」 | **⛔ 翻** | `15-c` 只给空 final 信号、`15-e` 只记「没听清」；**无重试/错误类** |
| AC6 | ASR 推理单核钉死、无线程池扇出、>3s 单核只在 work-peak 规则内 | **可翻**（边界见 §12） | `15-b` 的 `IntraOpThreads=1`＋CPU 形状；`15-c` 的 ASR config literal 1 |

★票面本身：`Status: ready-for-agent`、`Claimed by: —`、`Last update: 2026-09-19`（**37 天**）、`Progress log` 标题下 **0 行**、6 条 `- [ ]` **一条未勾**。前置 `02/13/14` 均 `-done`、0 未勾；**未申报前置**票 12 仍 `review`（6 勾/3 未勾）、票 26 仍 `ready-for-agent`（0/7）。⇒ **五片证据在 `docs/evidence/`，票面未同步**；翻框/改 `Status` 是编排者收尾动作，⛔ 本腿写面。
## 2 核心判点一：五片零回归的双向（我自己造）
| 尺 | 改前（`a736ad0ea`/`ca8834d72` 状态） | 改后（`2c0be38dd`，HEAD 同代码） | 判 |
|---|---|---|---|
| `internal/speech` staging 全套 | **28 PASS／0 FAIL／0 SKIP**，rc 0 | **31 PASS／0 FAIL／0 SKIP**，rc 0 | ✅ |
| 原 28 条**逐名** | — | 当前 31 名集合差集 **missing=0**；新增恰 3 条：`TestC9EngineRostersArePinned`／`TestC9ReservedSlotsAreImplementable`／`TestTheTranscriberRunsOverAnyAsrEngine` | ✅ |
| `cmd/wisp` 15e staging | **6 案 6 过**，rc 0 | **6 案 6 过**，rc 0 | ✅ |
| 15e 延迟三数 | `openDelay=262ms chunk=640ms ballLatency=902ms` | **逐位同值**：262／640／902 | ✅ |
| 15e 关键计数 | `handoffs=1 finals=1 handed=1 partials=5` | 逐位同值；final `"嗯嗯嗯嗯嗯嗯嗯哦哦"`；`exit=0` | ✅ |
| 静音正控 | `handoffs=0 finals=0 ball=[] state=Sleeping` | 逐位同值；`drained=120/120` | ✅ |
| `cmd/wisp` 293 旧断言 | — | **7/7 PASS**（不是 4 条；文件里 7 个 `TestAC293*`） | ✅ |
| `internal/audio` | `ok 16.7s` | `ok 16.711s` | ✅ |

★★**不当作回归的两个漂移**：①`drained` 我这一发 **181/181**、15e 记 **178/181**、15a 复跑 **166/171/181** ⇒ 192ms 有界通道在回放解码窗口的丢帧，**⛔ 当常量**；②15e 球走态第二边我这一发 `AtSample=92672（5.792s）`、15e 记 `89600（5.6s）` ⇒ 同一 D43#11 边、位置随丢帧漂移，**⛔ 当矛盾**（首边两边都是 864ms）。
## 3 核心判点二：`15-a` 接口的牙（五发，全我自己造）
突变件在 `D:\tmp\wisp\20261026\15v1\mut\`，经 `-overlay` 换入；工作树零突变窗口；最后一发变异后重编无 overlay 的 31 案 ⇒ rc 0、31/31。

| 发 | 变异 | rc | 逐字读数（节选） |
|---|---|---|---|
| A | `asr.go`：`Fed`→`FedSamples` | **1（编不过）** | `*StreamingAsr does not implement AsrEngine (missing method Fed)`；`engines.go:78` 的编译期断言先响 |
| B | 整段合成 TTS | **1（编不过）** | `wholeClipTts does not implement TtsEngine (wrong type for method Speak)`／`have Speak(string) ([]float32, error)`／`want Speak(string, func([]float32) error) error` |
| C | 只会报词的 KWS | **1（编不过）** | `wordsOnlyKws does not implement WakeWordEngine (wrong type for method Accept)`／`have Accept([]float32) ([]string, bool)`／`want Accept([]float32, []KeywordHit) []KeywordHit` |
| D | `RealtimeEngine` 调工具 | **1（编不过）** | **`e.CallTool undefined (type RealtimeEngine has no field or method CallTool)`** ⇒ 零工具权限是**结构性**的 |
| E | 给 `WakeWordEngine` **加** `Armed() bool`（stub 一并补，编译保持绿） | **编 0／跑 1** | `WakeWordEngine roster = [Accept Armed Close], want exactly [Accept Close]`（`engines_test.go:53`，`reflect` 名册钉） |

★★★**E 是这套牙里唯一的“加法”钉子**：Go 的结构化类型对**加方法**是瞎的（编译 rc 0、所有实现仍满足、所有调用点仍过），**只有 `c9Rosters` 的 `reflect` 名册会响**。⇒ C32 的“零工具权限”不是评审语，是**加一个方法就得改一行数据并被拒**。四发“减法/错签名”由编译器负责，E 由数据名册负责，两边互补。
## 4 核心判点三：`15-e` 端到端（我自己跑的那一发）
`Test15eSpeechEndToEndHandsTheTranscriptToTheAgent`（staging，`-test.run=Test15e`，6/6、rc 0）：

| 读数 | 我的值 | 备注 |
|---|---|---|
| 帧 | **181 帧**；`drained=181/181` | 夹具 1 暖机帧＋250ms 静音＋3s 合成语音＋2.5s 静音 |
| `handOffModel` | **被调 1 次**；`handoffs=1`；参数 `id="asr-streaming-paraformer-zh-en"` | 测试的 `seen` 也只 1 条同 id；**⛔ 每句一次、⛔ 启动时一次** |
| 交还 | **`exit=0`**；生产路 `residentVoiceHandOff`→`handOffModel`（`build/models/manifest.json`＋`.minisig` 真验签→`Ensure`→`VerifyInstalled`） | `wisp models ensure: … 已交还（state=FirstRun，ticks=1，末次进度 100%，3.38s）` |
| 转写 | `partials=5 finals=1 handed=1`；final **`"嗯嗯嗯嗯嗯嗯嗯哦哦"`** | `speech=3392ms closedBy=end`；`problems=[]` |
| 进文本路 | `submitTask(text)` 收到**同一串**；`handed=1` | seam 是 `voiceTranscriptSink`（生产装配传 `residentTaskSource`）；边界见 §12-AC1 |
| 球走态 | `[Listening, Thinking]`；timeline 2 条：`Sleeping+session.summon→Listening`（**13824＝864ms**）、`Listening+asr.vad-stop→Thinking`（**92672＝5.792s**） | D43 表**零 diff**；`AtSample` 位置随丢帧漂移，见 §2 |
| 延迟两段 | **开句迟 262ms ＋ chunk 640ms ＝ 球延迟 902ms**（只算 chunk **低估 41%**） | `seamWaitToFirstText=5.51s`；`finalLatency=4.096s` |
| 静音正控 | **120 帧 ⇒ `handoffs=0 finals=0 ball=[] state=Sleeping`**，`drained=120/120` | 绿不是“什么都没跑” |
| C8 0→1 | 15e 前的 `12be96ad3`：`NewBoundedFrames` 只有分配点、**接收点 0**；现在 `resident_voice_windows.go:462/489/691` 三个接收点**同属一个 `asr-infer` 消费者** | 无消费者正控：`capacity=6 accepted=6 refused=14 window=192ms` |

★**15-c 的同一支 `Transcriber` 复现 640ms 逐位同值**；开句迟 15c 记 1088ms 是另一夹具（第 0 样本即语音、且是反推）⇒ **两数如实并列，⛔ 当矛盾**。
## 5 核心判点四：`15-c` 的 `B` 那发（流式→整段，我自己造）
变异：`asr.go` 的 `Feed` 只攒不喂、`Finish` 才一次性灌入（仓外 `Bseg_asr.go` ＋ `-overlay`）。同一夹具（1.2s 语音＋1500ms 静音；3s 语音那发用于 partial 时序）：

| 量 | 流式（原样） | 整段（B 变异） | 判 |
|---|---|---|---|
| 句内 partial | **3 条**（frame 44/63/82，首条 `utterance t=640ms`） | **0 条**（`0 partial(s) over 3s of speech, want at least 2`） | ★**流式真在流的唯一证明** |
| 首字延迟（句起点→首文本，音频钟） | **640ms**（预算 800ms，余 160ms） | **1.792s（句闭）＋ 252.99ms（flush）＝ 2.045s**；对 800ms 预算 **2.56×** | ❌ |
| 判停→final 的 flush 腿（墙钟） | **529.8µs**（低于本机 ~0.57ms tick） | **252.9935ms**；对 0.53ms 基线 **477×** | ❌ |
| speak-pause-speak | 3 partial／2 final | **0 partial／2 final**，`no partials were produced while the utterances were open` | ❌ |

★我这一发的 2.56×／477× 与 15c 记的 **2.5×／205×** 同向不同值：**墙钟腿不是常量**（15c 基线 1.05ms、变异 214.6ms；我 0.53ms、253ms）；**确定性的是 0 partial ＋ final 仍报 28608 样本（1.792s）**。⇒ 结论不靠倍数，靠“句内 0 partial、首文本只在句闭后出现”。`TestFinalArrivesWithinTheBudgetOfTheVadStop` 在 B 下 **PASS**（253ms<800ms），恰好说明真正的罪是“首字被推迟到句尾”，不是那个 flush。
## 6 核心判点五：`15-d` 的「⛔ 永远绿」（我自己造）
| 发 | 树 | rc | 读数 |
|---|---|---|---|
| 缺失 | 仓外 `certree`（有 `go.mod`＋`cer` 8 文件、**无 `testdata`**） | **0** | **23 PASS／1 SKIP／0 FAIL**；SKIP 逐字：`baseline corpus absent (…sources.tsv): … SKIP, not pass: a green over 0 utterances would be a lie.` |
| **skip 前先断言** | 同上 | 0 | `baseline_test.go:137` 先 `LoadSet` 并断言 `errors.Is(err, ErrManifestMissing)`，**然后**才 `t.Skipf`（`:138` 是那条断言） |
| 变异 C：缺失⇒空集 | 同上＋`sources.go` 改 `return []Source{}, nil`＋删 `LoadSet` 零行守卫 | **1** | **20 PASS／4 FAIL**：`baseline_test.go:138 absent baseline set: LoadSet returned <nil>, want ErrManifestMissing`（★**连 skip 那条也变红**）／`:176`／`sources_test.go:129`／`:236` |
| 今天（基线集已落） | 仓内 `testdata/asr-baseline`（5 列、60+40） | **0** | **24 PASS／0 SKIP／0 FAIL**；`checkCommittedSet` 真跑 60+40、逐 wav 哈希、perfect=0、empty=100% fail；**skip 按其设计消失** |

★三重证明成立：①缺失是**具名错误**（`ErrManifestMissing`＋`fs.ErrNotExist`），不是空集；②skip **之前**先断言具名错误，变异 C 让这条变红；③`Set{}.Score`／`Report{}.Verdict`／`Aggregate{}.Judge` 无分母全部具名报错。`runtests.sh:224-235` 把 `--- SKIP` 判红是**故意的**：基线集必须先于 CI 接线落地——**现在基线集先落、CI 后接（`8d6e150f8`），我现跑 24/0/0，SKIP 不存在**；⛔ 把 skip 删掉或把红线当成缺陷。
## 7 核心判点六：`sources.tsv` 契约 vs 在飞的基线集腿（我的判语）
- **17:52 实测（碰撞真实发生）**：基线集腿先交的 `sources.tsv` 是 **10 列**（`file⇥sha256⇥source_utt_id⇥duration_s⇥speaker⇥condition⇥snr_db⇥noise_source⇥text_raw⇥text_norm`），15d 冻的是 **5 列**（`file⇥sha256⇥speaker⇥condition⇥origin`）⇒ 现行 `go test ./internal/speech/cer` **rc 1、23 PASS／1 FAIL／0 SKIP**，逐字 `baseline manifest header is not the frozen one: line 1 is "file\tsha256\tsource_utt_id…", want "file\tsha256\tspeaker\tcondition\torigin"`。**不是 15d 的错**：基线集腿先派、15d 后冻契约（时序风险，`15-a1` 预测过）。
- **18:0x 已闭环**：`c6ea9f66e`（15d3-r1）把 `sources.tsv` 改成 **5 列**、把多出列另存 `sources-extra.tsv`（6 列：`file/source_utt_id/duration_s/snr_db/noise_source/text_raw`）；`3ccb8c198`（15d4-r1）改生成器，**重跑仍产 5 列＋6 列**，`sources.tsv` 首行与 `sources.go:26` 的 `manifestHeader` **我现跑逐字节 `-ceq` = yes**。现行 `cer` **24 PASS／0 SKIP／0 FAIL**。
- **怎么记**：★**记为“时序残余——已闭环”**，不记为 AC2 的现行阻塞；证据＝两次 live 读数（23/1 → 24/24）＋两笔提交＋`sources-extra.tsv` 对读。★**仍要盯**：任何未来生成器重跑必须保持 5 列首行（15d4 已做，重跑逐字节同）；`sources-extra.tsv` **不被 15d 读**（无第二真相源）。★**AC2 的新阻塞是“真 ASR 跑分”**（CI 已由 `8d6e150f8` 接成 `runtests.sh` step，24/0/0；但只证 gate），见 §12/§13。
## 8 核心判点七：编排者准的那处（`NewTranscriber` 参数放宽）
**判：对，我复核通过。** `git show --numstat 2c0be38dd -- internal/speech/transcriber.go` **8 增／6 删**：两处类型（字段＋参数）＋两处注释，**零语义**；`git diff` 逐字只有 `*StreamingAsr`→`AsrEngine` 与注释。证据：`internal/speech` **28→31、missing=0**；15e **6/6**；`TestTheTranscriberRunsOverAnyAsrEngine` 用**与 sherpa 零共享代码**的 `cloudStubAsr` 走**真判停器**，逐字 `19 partial(s), 1 final(s), final "far-side transcript" closed by end over 28608 samples (1.792s)`，且 `last.Latency == samplesDuration(eng.Fed())` **逐位成立**。
★**理由**：C9 的价值就是“可互换”；不改参数类型 ⇒ 接口**零消费者**、纯装饰——正是 `15-a1` 判的“建了没人调”。**零行为改动、代价极小**；而“不许动参数类型”会让这条证明消失，**为守纪律丢价值**，正是铁律第 1 条禁的。★**代价**：`RealtimeEngine` 的方法集在票 60 前被冻结一版（零实现零 importer；§3-E 的名册钉让改进来必须显式过门）——可接受。
## 9 核心判点八：`d22scan` 分母涨了没有
**涨了，不是空跑绿。** 现跑 `cd tools/d22scan && go run . -root <repo>`：**rc 0** `clean - no D22 ban violations`；`ban #8 internal/ = 666`（`15-a` 前 `ca8834d72` 导出树同工具 **664**）⇒ **+2**；`bans #1-5 internal/ = 259`（前 **258**）；`cmd/ = 155` 不变；`ban #8 cmd/ = 155`。两枚新 `.go`（`engines.go`／`engines_test.go`）被真走到。★父提交导出树那次 rc 1 是**导出缺 `docs/`/`.scratch/` 的 phantom-citation**，⛔ 代码红；分母 664／258 是同一工具现跑的。
## 10 核心判点九：`cmd/wisp` 的边界（4 个既有文件）
`ca8834d72` 改 4 个既有文件＋新 2 枚；我逐项核：

| 面 | 读数 | 判 |
|---|---|---|
| `models.go` | `git diff 12be96ad3 ca8834d72` 全在注释：旧“doc.go-only”段改成 15e 已接线；**无代码行** | ✅ |
| `resident_audio_windows.go` | 加 `frames`/`voice` 字段、`attachVoice`、`stop()` 在 `gate.Stop()` 前 `voice.stop()`（D38(e) 第 4 步的同一 body，不是第 11 步）、`posture()` 改读 `drainsFrames()/absentReason()` | ✅ 语义变更**有意**，nil-safe；无旧断言红 |
| `resident_ball_windows.go` | `OnTrayPauseWake` 从 `recordBallGesture` 改为 `rb.pauseWakeGesture()`；`mirrorTrayMute` 第二 bit 从 literal false 改为 `wakeRead()`（无读者仍 false）；`setVoiceState` 先记账再投 | ✅ 293 的旧语义保留 |
| `resident_windows.go` | 组装根在 `startResidentAudio` 后 `attachVoice(startResidentVoice(...))`；`voiceTasks` 只在 `src.run != nil` 时非 nil（防 nil deref） | ✅ |
| D43 表 | `git diff 12be96ad3 ca8834d72 -- internal/statemachine` **空**；工作树 `git status --porcelain -- internal/statemachine` **空** | ✅ 一行未改，只 `Dispatch` 合法行 |
| 旧断言 | `TestAC293*` **7/7 PASS**；`TestAC3LogSink…` PASS（0.102s）；`go vet ./cmd/wisp` rc 0；15e 6/6 | ✅ |

★唯一需要记的语义变更＝**关机顺序**（语音链先停、设备后停）；理由写在 `resident_audio.stop` 注释里（先释放 sherpa 会话再停线程），测试全绿，**未发现破别的东西**。
## 11 四数复跑 ＋ 格式甲乙两尺 ＋ 门禁 rc
| 面 | 命令 | rc | 读数 |
|---|---|---|---|
| `internal/speech` staging | `speech.test.exe -test.v`（三 DLL 同目录） | **0** | **31 PASS／0 FAIL／0 SKIP** |
| `cmd/wisp` 15e staging | `wisp.test.exe -test.run=Test15e -test.v` | **0** | **6/6 PASS**；`262/640/902`、`handoffs=1 finals=1 handed=1`、`drained=181/181` |
| `internal/audio` | `go test -count=1 ./internal/audio` | **0** | `ok … 16.711s` |
| `cmd/wisp` 定向 | `-run TestAC3LogSinkLandsInsideTheEnvDataRoot` | **0** | `--- PASS … 0.102s` |
| `internal/speech/cer` 现行 | `go test -count=1 -v ./internal/speech/cer` | **0** | **24 PASS／0 SKIP／0 FAIL**（60+40 已载入） |
| CI CER gate | `bash tools/d22scan/runtests.sh ./internal/speech/cer/` | **0** | `PASS=24 FAIL=0 SKIP=0`（我现跑；即 `8d6e150f8` 的新 step） |
| 293 定向 | `-test.run=TestAC293` | **0** | **7/7 PASS** |
| 门禁 `d22scan` | `cd tools/d22scan && go run . -root <repo>` | **0** | `clean`；`ban#8 internal/=666`、`bans#1-5 internal/=259` |
| `go vet` | `./internal/speech`；`./cmd/wisp` | **0／0** | 空 |
| 甲 `gofmt -l internal cmd tools scripts frontend` | — | **0** | 恰 3 行：`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`／`cmd\wisp\models.go` |
| 乙 `gofumpt -l`（`/d/work/base/gopath/bin/gofumpt.exe`） | — | **0** | **同上 3 行、逐字相同**；`internal/speech` 与 `cmd/wisp` 新文件**两尺皆空** |
| 三层 DLL | ①无 staging 无 PATH：**rc 1 `0xc0000135`**（进程起不来，⛔ 测试红）②只给 PATH：**rc 0、12 PASS／19 SKIP**，skip 原文 `onnxruntime 1.17.1 is older than the 1.28…a second onnxruntime.dll is shadowing the one on PATH`③staging：**31/31** | — | ★`PATH` 不够，**staging 才全绿**；15b 的 8/8 与 15c 的 10/18 是**包长大**，⛔ 当矛盾 |

★**3 条格式红是假红（终证）**：三枚 `blob CR=0`、`worktree CR=131/225/389`、`git diff HEAD` 空；把 blob 拷出来跑 `gofmt`／`gofumpt` **各 0 条** ⇒ `autocrlf` 显示面，**⛔ 缺陷**（尺清单第 17 处）。
## 12 票 15 的 6 条 AC 逐条：能翻／⛔ 翻＋缺什么
| AC | 判 | 缺什么／边界（具名） |
|---|---|---|
| **AC1** | **能翻** | 边界：测试用 `voice15eSink`（`submitTask` 的测试实现）；生产装配在 `resident_windows.go` 传 `residentTaskSource`。真 agent 回合归票 12（仍 `review` 6/3），**不是 15 的欠账**；wav→VAD→ASR→`handOffModel`→文本路 seam→球走态均已证。 |
| **AC2** | **⛔ 翻** | 缺：①把 `15-c` 的流式 ASR 跑在 100 句提交集上，产出**逐件＋聚合 CER 报告**（当前 `cer` 只有计算器/`checkCommittedSet` 的 perfect/empty，**没有 ASR 假设产生器**，`grep StreamingAsr internal/speech/cer` = 0）；②CI step 已接（`8d6e150f8`；`runtests.sh`；我现跑 24/0/0），但 cerci 自陈「6%/15% are EXERCISED, not MET」＋自托管 `asr-cer` job 未做；③`noisy_far` 15% 实测；若失败按 D44 记 `DEFERRED`＋UI-hint 票。基线集/契约/CI 已不是阻塞（§7）。 |
| **AC3** | **能翻** | 边界：300ms 门**只在状态机层承重**（299→0、300/301→1 转移）；真 silero 对 150/250ms 合成突发根本不报语音，所以端到端那条“没有端点”不能单独证明 300ms 门。60s 真模型 `[start max-duration start end-of-stream]`、`Duration/Samples` 截到 60000/960000。 |
| **AC4** | **⛔ 翻** | 缺：①slot mutex（`internal/speech`/`cmd/wisp` 无 slot 实现，只有 `TtsEngine` 接口位）；②TTS（票 26 `ready-for-agent` 0/7）；③ASR+TTS 争用/串行驻留；④卸载 RSS delta（采样器）；⑤切换延迟 vs 1.6s。**15 自己凑不出实测**。 |
| **AC5** | **⛔ 翻** | 缺：①一次重试；②`asr` 错误类映射；③「没听清」的用户可见语义。现状：空 final 只记 `这一句没听清…（重试与错误类映射不在本片）`、**不交文本路**；15c/15e 都明写未做。 |
| **AC6** | **能翻**（边界） | 边界：`asrRecognizerConfig` 的 `NumThreads` **literal 1**、`EnableEndpoint=0`、16k/80 维由 `TestAsrRecognizerConfigPinsTheShape` 钉；**实测 CPU 形状来自 VAD 会话**（`intra_op=1: threads+0 cpu/wall=1.07`；`intra_op=4: threads+3 cpu/wall=3.95`；32s 负载）。ASR 自己的 cpu/wall 没有单独采样；若机主要“ASR 专属 >3s 采样”，那是一条小残余。 |

★**建议翻框顺序**：AC1/AC3/AC6 可立即翻；AC2/AC4/AC5 保持 `- [ ]` 并在票面写上 §13 的具名残余。**我不改票面**。
## 13 残余清单（具名，⛔ 模糊）
| # | 残余 | 证据／归属 |
|---|---|---|
| R1 | **AC2 的真 ASR 跑分**：无 `StreamingAsr`→`cer.Set.Score` 的 runner；CI step 已接（`8d6e150f8`；`runtests.sh`；24/0/0），但 cerci 自陈「6%/15% are EXERCISED, not MET」＋自托管 `asr-cer` job 未做，只证“集合可载入＋gate 有牙” | `internal/speech/cer` 无 ASR import；`cerci-r1-impl.md:42/:138`；`cerci-r1` 不属五片 |
| R2 | **AC4 slot mutex＋TTS**：无 slot、无 TTS 实现、无 RSS delta、无 1.6s 实测 | 票 26；`internal/speech/engines.go` 只 `TtsEngine` 接口 |
| R3 | **AC5 重试/错误类**：只有空 final 记录 | `cmd/wisp/resident_voice_windows.go:653`；15c/15e §未做 |
| R4 | **D43#12 `Listening` 超时**：机器定时器走、球不跟（无投影钩子） | 15e §2/§11；`internal/statemachine` 超时无外部钩子 |
| R5 | **首用成本**：交还 3.38s（我这一发）＋引擎装载 ~2.2s；`seamWaitToFirstText=5.51s`；首次下载分钟级 > 8s 回放环 ⇒ 首句会丢 | 15e §1/§5/§11；⛔ 隐藏 |
| R6 | **回放窗口丢帧**：192ms 有界通道；15e 178/181、15a 166/171/181、我 181/181 ⇒ **0–15 帧**，丢的是句尾，不影响出声句 | `resident_audio` 通道 cap 6；⛔ 当常量 |
| R7 | **KWS（41）／TTS（26）／AEC（59）**：接口位 reserved，未实现；`doc.go` 的 `DEFERRED(engines)` 仍写“ticket 15 实现”——**半过期**（15 已实现 ASR/VAD 半边） | `15-a` §9；摘 marker 属 `225-deferred-registry-audit` |
| R8 | **注释残债**：`internal/models/assembly_reachability_121_test.go:50` 仍说 `internal/speech` 是 `doc.go`；**测试不会红（名单手维护）**，但“注释说假话” | `15-a1` §3；收尾片/121 的账 |
| R9 | **C9 三个 reserved slot**（TTS/KWS/Realtime）零实现；`RealtimeEngine` 方法集已被 15-a 冻一版 | 票 26/41/60/61；§3-E 名册钉 |
| R10 | **`sources.tsv` 时序残余已闭环**（不是现行阻塞）：15d2 交 10 列 → 15d3 改 5 列＋extra → 15d4 改生成器；现行 24/24 | `c6ea9f66e`／`3ccb8c198`；§7 |
## 14 提交号 ＋ 逐笔越界（五片）
| 片 | commit | 短统计 | 路径／越界 |
|---|---|---|---|
| `15-b` VAD | `d5d25ddad` | **7 files, +1529** | `internal/speech` 6 枚＋本片证据；**无 `cmd/`、无 `testdata/`** |
| `15-c` 流式 ASR | `a736ad0ea` | **5 files, +1274／-0** | 4 枚 `.go`＋证据；**0 deletions** ⇒ 15-b 7 枚一字节未动 |
| `15-d` CER | `12be96ad3` | **8 files, +1776** | 全在 `internal/speech/cer/`；**无 `cmd/`、无 `testdata/`**；证据件另笔 `edec2597e` |
| `15-e` 装配 | `ca8834d72` | **7 files, +2129／-39** | `cmd/wisp` 4 改＋2 新＋证据；D43 表 0 diff |
| `15-a` 基座 | `2c0be38dd` | **4 files, +655／-6** | `engines.go`／`engines_test.go`／`transcriber.go`（8/6，§8）＋证据 |

★**逐笔越界 0**：15b/c/d 纯 `A`；15a **1 `M`＋3 `A`**；15e **4 `M`＋3 `A`**、**无 `D`**；五笔均无 `testdata/`、`.scratch/wisp/issues/`、`docs/specs/`、`docs/PLAN.md`、`docs/reports/`、`docs/audit/`。★五笔全为 HEAD 祖先；`git diff 2c0be38dd..HEAD -- internal cmd go.mod` **空**（后十一笔只动基线集/CI/证据/台账；五片之外的 `8d6e150f8` 改 `.github/workflows/ci.yml`＋`scripts/portable-tests.sh`，是 AC2 的 CI 半边）。
## 15 交回父代理
**未定义即停：无。** 唯一限制：`cerci-r1`（CI 接线）在**写件末落地**（`8d6e150f8`），我现跑其 `runtests.sh` **24/0/0**；其注释自陈不证真 ASR，AC2 仍 ⛔。HEAD 仍在动，复跑请以“五笔代码提交＋`git diff 2c0be38dd..HEAD -- internal cmd go.mod` 空”为快照；本文件落点时的 `ls`／HEAD／时间已写在件头。
