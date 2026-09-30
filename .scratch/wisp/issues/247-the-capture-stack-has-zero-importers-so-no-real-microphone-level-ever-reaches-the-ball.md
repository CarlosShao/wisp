# 247 — 采集栈写完了但**全仓零 importer**：真麦克风的电平到今天没有任何一处会去读，球上那个数是命令行手填的

Status: OPEN（编排者 09-30 23:0x 立，料全部出自 `240-c1`／票 241 与本编排者 22:5x–23:0x 的现跑尺）
实现者：`247-r1`（写码腿，**按住**，见「排程与串行」） · 前段普查：`247-a1`（只读，在飞） · 裁决者：`247-v1`（必须 ≠ 实现者） · 归口：语音链路（内序＝票 228 球进干活进程 → **真电平（本票）** → 唤醒词引擎）

## 现量（每条都带尺与取数时刻；⚠ 引用前先重跑，别把这几行当常量）

| # | 事实 | 尺与出处 |
|---|---|---|
| 1 | `internal/audio` 的**非测试 importer＝0 枚**：采集栈（真 WASAPI 共享模式 `wasapi_windows.go:210 Open`／`:361 Drain`）与票 241 落的电平生产者半（`level.go:96 LevelOfSamples`、`:118 FrameLevel`、`:137 DecodeFrame`、`wavinjector.go:138 EncodeFrame`）**全在库里、全绿、没有一处被生产代码调用** | 编排者 09-30 22:5x 现跑 `grep -rl "CarlosShao/wisp/internal/audio" --include=*.go . \| grep -v "/internal/audio/" \| grep -vc _test` ⇒ **0** |
| 2 | 球的电平入口 `Ball.SetAudioLevel`（`internal/ball/liquid_windows.go:42`）的**非测试生产者只有 `cmd/balldebug`**（`main.go:418`／`:477`，那枚 flag 的文案逐字 `synthetic audio envelope`）⇒ 今天屏幕上那个"音量在动"是**命令行里填出来的数**，不是声音 | 编排者 09-30 23:0x 现跑 `grep -rn "SetAudioLevel" --include=*.go internal/ cmd/` ⇒ 产码命中只有 `balldebug` |
| 3 | `internal/speech` **只有 `doc.go` 一枚文件**＝ASR/TTS/KWS 那一整块一块都没写 | `ls internal/speech/` 现跑；与项目记忆「240-c1 三分类」一致 |
| 4 | 两枚**方向相反**的开关：`voice.enabled` 默认 **true**、`wake_word.enabled` 默认 **false**（台账 A462 我已按保守默认定：两道都得 owner 手动开） | 票 240-c1／A462；本票的接线不许把它俩的默认值动哪怕一格 |
| 5 | `Sink` 机制与测试都在，但**生产两处建机器都没传** ⇒ 四枚 kws/audio 事件往空里发（`SetLoaded` 同形） | `240-c1` |
| 6 | 「第六环是形状」：`liquidDriven`（`internal/ball/liquid.go:60-68`）**不含 `Armed`** ⇒ 判"接好了"必须看 Armed 收到电平，而不是看 Sleeping 那格 | 票 241 票面 `:39` 那一行遗留的归口；`liquid_windows.go:38-40` 逐字 **no samples and no transcript cross into the ball**（过界的只许有一个 `float32`） |

## 本票要做的（一句话）

把「**真麦克风的 PCM → 一个 `0..1` 的数 → 球的液态**」这一段接成真的，并且**当场回答"接到哪条腿上、谁拥有那个协程、麦克风什么时候真的开、开不了怎么降级"**——这四问里任何一问没答上就不算接好。

## 判据（⛔ 框归编排者，产码腿与验收腿一枚都不许碰）

- [ ] **AC#0 先把四问的代价摆开（不许直接开写）**：① **落点**＝采集协程归常驻那条腿（`cmd/wisp`，与票 246 同一枚进程）还是归 `internal/audio` 自己起？两形各写"要动哪几枚文件＋新增哪几条包级依赖边＋协程 owner 走不走 `observe`／`rt.Registry`"。② **隐私闸门**＝接上之后麦克风在什么条件下真的打开，默认档（`voice.enabled=true` × `wake_word.enabled=false`）下今天这台机器会不会一双击就采音？③ **降级**＝设备被占／无权限／无设备时球的形状与文案说什么（SPEC-05 §3.4：分类＋可见，不许静默）。④ **第六环形状**＝`Armed` 那一格（现量 6）本票做还是留给谁。完成判据＝四问各带现读凭据（文件:行＋尺的读数），由编排者裁后再派 `247-r1`。⛔ **普查腿不许改任何产码**，写点只准落在 `.scratch/wisp/probes/247/a1/census.md`。
- [ ] **AC#1 生产者半有人真调用**：接完之后 `internal/audio` 的非测试 importer **≥1 枚**（尺＝现量 1 那把，接前接后各一发，读数进表），且那枚 importer 就是**跑着的那条腿**（⛔ 不许是又一枚调试用 cmd）。
- [ ] **AC#2 球上的数来自声音，不来自命令行**：一发真机读数＝对麦克风说话与不说话时 `SetAudioLevel` 收到的值**不同**，并给出两形的采样数与出处；⛔ 不许用 `balldebug` 的合成包络冒充（现量 2 那条尾巴要剪掉，不是把它做大）。
- [ ] **AC#3 只过一枚 `float32`**：任何新增代码里，跨过 `internal/ball` 边界的数据**只许是一个标量电平**。尺＝能力型（新增的跨包调用签名里出现 `[]byte`／`[]int16`／samples 即判越界），⛔ 不许做成"扫注释里有没有 samples 这个词"的词面型。
- [ ] **AC#4 隐私默树一格都不许动**：`voice.enabled`／`wake_word.enabled` 的**默认值与读取处**一字不改，且真机读数要证明"默认档下麦克风不会被这条腿打开"（或反过来具名说出它打开了、由哪一行决定）。
- [ ] **AC#5 协程有 owner、有 recover**：⛔ 裸 `go func(`（ban #1，`tools/d22scan/d22scan.exe` 会点红），采集协程必须挂在现成的 owner 上（`observe.Root`／`rt.Registry`），退出序列里被 join——**不许新造 D38(e) 的第 11 步，十步顺序一字不动**。
- [ ] **AC#6 降级照跑**：拔麦／无权限／设备被占三形里任取两形，本进程**仍要把球与任务管线跑起来**并且把损失说响亮（现量：票 128 定的是"没有数据根才拒绝启动"，别把它扩大）。
- [ ] **AC#7 越界检查**：`git diff` 里出现 `internal/speech`（一块没写的包）、`scripts/spike`、`frontend/**`、`design/**`、`PLAN.md`、`docs/specs/**`、`internal/observe/thresholds.go`、golden、`allowlist.txt` 任一路径 ⇒ 直接退回。唤醒词／ASR／TTS **都不属于本票**。
- [ ] **AC#8 门禁四数**：`GOFLAGS= go build ./...`、`gofumpt -l <自己动过的目录>`、`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/audio/ ./internal/ball/ -count=1`、`./tools/d22scan/d22scan.exe`（⚠ 它是**独立模块**，只能跑那枚 exe，`go run ./tools/d22scan` 必失败），逐名照抄终态。

- [ ] **AC#10（09-30 23:2x 编排者追加，来路＝`247-a1` ④ 节；勾要非实现者裁）**：**"接上了"不等于"看得见"**——本票把电平接进生产之后，屏幕上仍然**不会**出现"球随声音呼吸"，因为还有一道独立闸：`prototypeVisuals` 默认**关**（`internal/ball/statevisual.go:102`，产码里只有 `cmd/balldebug:122` 打开它），而 `SetAudioLevel` 在它关着的时候**直接 return**（`internal/ball/liquid_windows.go:52`）。⇒ 本格判据＝票面与本仓文档**不许**把"电平接好"写成"用户能看见球在动"；那一句归口**票 68 AC#2**（`statevisual.go:93` 自己写明"翻默认值"归它），⛔ 本腿不许顺手 `ball.EnablePrototypeVisuals(true)` 来让 AC#2 好看——owner 09-30 明说过「形状算你过关，好不好看以后再说」，翻视觉默认值是**样式决策**，不由接线腿代做。

## 编排者裁定（09-30 23:2x，台账 `A485`）：`247-a1` 的八问逐条判完，落地腿可以派了

`247-a1` 交的是 `.scratch/wisp/probes/247/a1/census.md`（271 行／54,536 字节，27 枚编号尺、三枚 `go list -deps` 名册原文已入盘）。八问裁定：

| 问 | 裁 | 依据与边界（写死，免得落地腿自己猜） |
|---|---|---|
| **P1 默认档谁挡麦** | **甲** | 接线时**必须**把 `c.Audio.MicMutedDefault` 传给 `internal/audio/gate.go:57 WithStartMuted`，并把 `voice.enabled=false` 判为**根本不构造采集器**。⛔ 零默认值改动。现量理由：`voice.enabled`／`wake_word.enabled` 今天**没有任何产码读取处**（只有 `internal/config/manager.go:370`  reload 分层），真能挡麦的是第三枚 `audio.mic_muted_default=true`，而它**只在门被装上时生效**——不装门直接 `mic.Start` 会在 `wasapimic_windows.go:88` 同步等设备开成功＝**一双击就开麦**。乙的"自然行为"取决于接线人是否包了门，不可接受。 |
| **P2 `Armed` 收不收电平** | **甲＋不做** | 本票**不动** `liquidDriven`（`internal/ball/liquid.go:64` 只含六态），"接好了"的判据挂 `Listening`（D43 #4 是今天真能到的态）。`Armed` 归口**唤醒词票**——因为 `EvKwsEnabled` **产码零发射者**（`internal/statemachine/table.go:55`），且 `SPEC-08 §2.1` 逐字写着 Armed **静态**＝改它要先落一枚人工批准的 `A##`，⛔ 不由写码腿顺手做。 |
| **P3 视觉那道闸归谁** | **乙** | AC#2 的取证改用**读数型证据**（进出的两形读数＋`PrototypeVisualsEnabled()`），⛔ 不在常驻腿调 `EnablePrototypeVisuals(true)`。这条新增了一格 **AC#10** 把"看不见"这件事写进票面，防止后续程把"接好了"读成"用户看得见"。 |
| **P4 协程记进哪枚注册表** | **甲，限定形状** | 允许在 `internal/audio` **新增一枚注入位**（把 `wasapimic_windows.go:83`／`wavinjector.go:84` 写死的 `observe.Default.Spawn("audio-capture")` 改吃传入 registry；现成的 `audio.go:180 SpawnCapture(registry, run)` 就是那枚接缝，今天**零调用者**）。⛔ **不许改任何既有导出的语义**，⛔ 不许"两枚注册表并存"（那正是 `247-a1` 的 E6 假绿形状）。⚠ 这一步**动的是票 13/241 已交付的产物**，理由与边界具名记在本条与 `A485`，撤销口令＝「**247 注册表改回**」。 |
| **P5 设备失败要不要推状态** | **甲** | 不推状态：只把 `observe.ClassAudioDevice` 的分类＋`internal/audio/device.go:86/:89` 那两句 guidance **响亮地**打印并落 sink。⛔ **乙一律不要**——`table.go:97` 的 `EvAudioDeviceLost` 起点只有 `Listening`，启动期失败没有合法边，硬造＝改 D43 转移表（契约面）。⛔ 也不许把"拒绝启动"扩大（票 128 只定了没有数据根那一种）。 |
| **P6 中间那枚读帧的人放哪** | **甲** | 电平在**采集腿那一侧**算（`FrameLevel`→一枚 `float32`）再交出去，跨包面只许过那枚标量。乙要多两枚直接依赖，丙会把渲染拖进采集生命周期（`audio→ball` 这条边今天不存在，U6/U7）。 |
| **P7 文件清单要不要重抽** | **甲** | 派 `247-r1` 之前以 `246-r2` 的**终态**重跑那两把尺（它已落了 `cmd/wisp/resident_task_source_windows.go`，`RegisterShutdownHook` 的注册点从 1 处变 3 处——现量见 `A485`）。 |
| **P8 第二枚协程没有名册槽** | **甲** | 电平在**已有的 `audio-capture` 协程内**算完再投递 ⇒ D38(b) 那六名**零膨胀**（`observe/goroutine.go:42-44`/`:270-273` 对名册外名字会 `slog.Warn` 泄漏征兆；`internal/audio/level.go:85-88` 逐字写着这件事）。⛔ 不许用"复用同名"那一支（会让 `Resident` 计数变 2、把 `PLAN.md:2831` 那句"常驻 6 个"稀释掉）。代价要**实测**：投递路径进了 pinned 线程之后，`wasapimic_windows.go:149-152` 那个循环的周期读数要复量一次并进表。 |

**排程不变**：`247-r1` 仍排在 `246-r2` 之后（同占 `cmd/wisp`），并**排在票 33 宿主段之前还是之后由我按 owner 09-30 23:1x 那句功能要求定＝33 先**（他要先看见能点开的界面；本票是"球随声音动"，属他明说过后置的那一类）。

## 禁区（实现腿与验收腿共用）

- ⛔ **不动 `PLAN.md` 与 `docs/specs/**` 一字**；D1–D47／C1–C32／D43 转移表是契约面，改它只能由编排者在台账落一枚 `A##`。
- ⛔ 不动 `internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）；不许为变绿放宽任何断言。
- ⛔ `frontend/**`／`design/**` **两层禁令**：不许读，结论也不许引到它们身上。
- ⛔ **凭据绝不进对话／日志／表**（本票不涉凭据，写在这是防"顺手接云端 ASR"那种越界）。云端那一段属票 61 的归口。
- ⛔ 不许新增词面型仪器；判据一律问能力。
- git：只 commit 不 push；commit 必带显式 pathspec；禁 `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件。

## 排程与串行

- ⛔ **`247-r1`（落地腿）按住，排在 `246-r2` 之后**：本票的落点几乎必然要碰 `cmd/wisp`（常驻那条腿），而 `246-r2` 此刻正在那枚包里改 `run.go`／`resident_windows.go`／`resident_approval_windows.go`。**同包两枚写腿会把红的归因搅浑**（本仓 09-30 实测过两次代价）。`247-a1` 是只读普查，不受此限，但⛔ **它禁跑 `go build`／`go vet`／`go test`**（会吃到 `246-r2` 写到一半的码，读数是脏的）。
- ⚠ 与票 228 的后续片（config.toml 接线、托盘「退出」执行者、球位置持久化）同属地界重叠区：那几枚若与本票都要动 `cmd/wisp`，一律串行。
- ⚠ 跑真机用例前 `tasklist //FI "IMAGENAME eq balldebug.exe"` 与 `wisp.exe` 都必须为 **0**；`internal/audio` 的 `TestLiveWasapiSmoke` 与 `TestPinnedThreadStable10s` 碰真设备，同理。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
- [2026-09-30 23:17:40 +08] agent=247-a1 did=**只读普查腿交件**，写点唯一＝`.scratch/wisp/probes/247/a1/census.md`（271 行／54,536 字节；骨架 60585af9→⑤⑥⑦ 先满 457a824e→①-④ 77a3b442→复量 ⑤bis 5c33ad30）。四问结论各一行（详表＋27 枚尺的原样命令与真实读数在盘上）：**① 落点＝甲形**（`cmd/wisp` 拥有采集），新增包级直接边**恰好 1 枚** `cmd/wisp→internal/audio`，尺 `comm -23 deps-audio deps-cmdwisp` 只出这一名（audio 其余 143 名传递依赖已在 cmd/wisp 的 276 一名册里）；乙-1 直推要多 `audio→ball`（U28：＋`internal/ball`＋`internal/statemachine` 两名）并破 `wasapi_other.go` 的跨 GOOS 承诺，乙-2 注入**与甲同边、并不省 cmd/wisp 改动**；退出 join 用现成接缝 `proc.RegisterShutdownHook(StepStopAudio)`（`boot_windows.go:188`／`shutdown_hooks.go:118-119`／`shutdown.go:163`），`Stop()` 自身已带 2s 有界 join（`wasapimic_windows.go:94-115`）⇒ **十步一字未动、没新造第 11 步**；⚠ 本件新发现一条硬约束：D38(b) 名册里**没有第二枚电平协程的槽**（`observe/goroutine.go:42-44`＋`:81`＋`:270-273`；`internal/audio/level.go:85-88` 的明文早已写着这句话）⇒ 交裁 P8。**② 隐私闸门＝今天不会采音，但挡住它的不是默认值而是"没路径"**：`voice.enabled`=`schema.go:245` true、`wake_word.enabled`=`:197` false 两枚的**产码读取处只有 `internal/config/manager.go:370/:384/:385`**（reload 分层，不是开麦决策，尺 U13）；真挡得住的是派单未点名的第三枚 `audio.mic_muted_default`=`schema.go:277` **true**，但它**只在门被装上时生效**（`gate.go:108` 那个 `if g.effectiveOpen()` × `gate.go:235-237` muted 先返 false ⇒ `openInnerLocked` 不跑、`wasapi_windows.go:210 Open` 不跑）；**不装门**直接 `mic.Start` 则在 `wasapimic_windows.go:88` 同步等设备开成功（送 nil 在 `:182`）⇒ 一双击就开麦。**零默认值改动**，只写"改了会怎样"。**③ 降级＝照跑、响亮、不扩拒绝启动**：四枚用例的**断言本体**逐行读（`hotplug_test.go:378` 五枚断言＝Start 必返 err／class `audio_device`／串含设备名"Busy Mic"／含 guidance／`Stats().LastError` 非空；`:199` 重举**恰一次**且第二枚用新默认设备；`:260` 终态 err＋**buf 从此不再出帧**；`:334` 重举失败要含**上一枚**设备名），文案备好在 `device.go:86/89/103/105/107`；参照系用仓里现成三枚（`resident_ball_windows.go:139-145` 有 err 也 `return rb` 不 `os.Exit`、`resident_windows.go:62-69` "Loud, and it does not stop the app"、246-r2 新落的 `resident_task_source_windows.go:289-293` 注册失败也只 err＋Printf"记为 skipped"）；⛔ 明确判掉一条歧路：D43 #14（`table.go:97`）的 `From` **只有 `Listening`**，启动期失败无合法边，硬造＝改转移表＝契约面。**④ 第六环＝本票不该以"`Armed` 收到电平"为判据**：`liquidDriven` 逐行读数只含 6 态（`liquid.go:64-66`：Listening/Thinking/Speaking/Conversation/Warm/Acting），其余含 Armed/Muted/Sleeping 全 false（20 枚名册自 `states.go:10-30`）；三条实码＝`EvKwsEnabled` **产码零发射者**（U16，Armed 今天不可达）＋`docs/specs/SPEC-08-ui-ball-panel.md:61` 明写 Armed **静态**（同批改它撞 AC#7 的禁改路径）＋`tokens_test.go:155` 只钉 `Opacity` 而 `applyTo` 不写 `Opacity`（所以我只说地界与票 62/74 重叠，不说必然红，见 E8）。**本节最重要产出（牵 AC#2 取证形态）**：即使接上，球今天仍不会"呼吸"——`prototypeVisuals` 默认关（`statevisual.go:102`，产码只有 `balldebug:122` 打开）而 `SetAudioLevel` 在它关时直接 return（`liquid_windows.go:51-54`），且常驻腿可达的 `Sleeping`／`Confirming` 两枚都不在 `liquidDriven` 里 ⇒ 本票能证"电平进了球的门"，证不到"球在动"。⑤bis 复量：`246-r2` 的 `bb37fac2` 在我写 ①-④ 期间落盘，`SetState` 与"audio importer＝0"两枚读数**未过期**，`RegisterShutdownHook` 从 1 处变 3 处（`resident_windows.go:156→157`）。evidence=`.scratch/wisp/probes/247/a1/census.md`（＋名册三枚临时件 `deps-audio.txt`／`deps-ball.txt`／`deps-cmdwisp(-win).txt`，按"只建不删"留着）next=**交编排者裁 8 件**（⑦ 表，逐条甲／乙／不做）：P1 默认档谁挡麦（要不要强制 `WithStartMuted(c.Audio.MicMutedDefault)`＋`voice.enabled=false` 不构造采集器）／P2 Armed 要不要收电平（我判本票不补，归唤醒词票）／P3 `prototypeVisuals` 这道可见性闸归本票还是票 68 AC#2（**决定 AC#2 的取证形态**）／P4 采集协程记进 `observe.Default` 还是 `rt.Registry`（`wasapimic_windows.go:78` 写死前者，假绿风险见 E6）／P5 降级不推状态、只响亮说（乙＝改 D43 表，我建议一律不要）／P6 电平读数放 `cmd/wisp` 还是新包／P7 派单前要不要按 246-r2 终态重抽 ① 的文件清单／P8 名册无槽三选一。⛔ AC 框一枚未碰（8 未勾／0 已勾），票面除本行 Progress log 外一字未改；零产码改动；未跑 `go build`／`vet`／`test`／任何 `./...`；`frontend/**`／`design/**` 未读未引。
