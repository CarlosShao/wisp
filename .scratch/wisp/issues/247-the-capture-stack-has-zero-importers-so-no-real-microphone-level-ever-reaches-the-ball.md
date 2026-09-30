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
