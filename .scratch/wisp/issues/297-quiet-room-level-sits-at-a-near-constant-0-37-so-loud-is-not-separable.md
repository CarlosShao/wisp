# 票 297 — 安静房间的 delivered level 就有 **mean≈0.3677（−8.6 dBFS）**、六窗 max/mean 全挤在 1.08–1.11：今天"响得起来"那一格判不了，因为**不知道喂给电平尺的那串数是什么**（真机现量，缺的是仪器不是人）

**立票**：2026-10-09 18:1x 编排者（来路＝**票 247 `AC#2` 那两发真机读数**：机主在场的真机窗口，我亲跑；件 `.scratch/wisp/probes/orch/2026-10-09-live-window-findings.md` §1，逐字原始输出在 `…-a1-run1-quiet.raw.md` / `…-a1-run2-voice.raw.md`）
**性质**：★**这一格今天既不能判绿也不能判红**。票 247 `AC#2` 的红句是"the loud phase is not distinguishable from silence"，而它给出的六个 mean 全挤在 0.367–0.383 —— **"分不开"这件事本身已经测到了**，但**为什么分不开**有三个形状、代价与归口各不相同，本票只裁"是什么"。⛔ 本票**不许**被读成"电平链坏了"，也不许被读成"链路通了、只差声音不够大"——两条我都没有证据。
**为什么要紧**：票 291 那把尺的定义（"电平该随声音变"）与票 247 `AC#2`/`AC#4` 的收口**都卡在这一问上**；而球的液态（票 68）将来吃的就是这枚数——如果它吃到的是一条 0.37 的近常数底，球会**一直在呼吸、且对声音无感**，那是用户第一眼就看得到的一等缺陷。

## 现量（每条都带尺；⛔ 引用前先重跑，行号与读数都是快照）

- **六窗读数**（尺＝`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -run 'TestAC247LiveMicrophoneLevelsReachTheBallSeam' -count=1 -v`；打印点逐字 `cmd/wisp/resident_audio_247_live_windows_test.go:184`）：
  - run1（完全安静）：`quiet-A: samples=93 mean=0.367691 max=0.409285 levels_per_s=31.00` / `sound: samples=189 mean=0.382482 max=0.413813 levels_per_s=31.50` / `quiet-B: samples=94 mean=0.383328 max=0.408784 levels_per_s=31.33`
  - run2（机主开口 6 秒）：`quiet-A: … mean=0.369753 max=0.408560` / `sound: … mean=0.382104 max=0.402393` / `quiet-B: … mean=0.381650 max=0.405371`
  - ⇒ **sound/quiet-A＝1.0402**（判据要 2 倍，`:194`）；⚠⚠ **`quiet-B` 比 `sound` 还高**（run1：0.383328 > 0.382482）——"最吵的那一窗"不是最响的；⚠ run2 的 `sound`（外放＋人声）比 run1 的 `sound`（只有外放）**低 0.000378** ⇒ 人声在那把尺下**没有可测贡献**。
- **计数器（同一发，逐字 `:181`）**：run1 `AC#2 counters at the seam: levels=376 frames_sent=6 frames_dropped=370 last_err=""`；run2 `levels=375 frames_sent=6 frames_dropped=369`；两发退出行逐字（`:26`）`wisp: audio capture leg stopped: levels_delivered=376 frames_sent=6 frames_dropped=370 reopens=0`。⇒ **31 枚/秒、32 ms 一窗、`last_err=""`、零 reopen**——链路在跑，不是拿不到数。
- **电平尺在"零"上是有钉的**（⇒ 0.37 不是尺子的偏置）：`internal/audio/level.go` 的 `LevelOfSamples` 是纯 RMS/`LevelFullScale`（`:96-107`，逐字 `	rms := math.Sqrt(float64(sumSquares) / float64(len(samples)))` / `	return rms / LevelFullScale`，**无任何加性底**），而"全零帧读 0"这条有常驻判据（尺＝`grep -n "silent\|MinLevel" internal/audio/level_test.go` ⇒ `:46` 那节 `// AC#1 first of the three: …` 起，`:52-55` 断言全零帧读 `MinLevel`）。⇒ **0.3677 只能是"喂进去的字节不是零"**：`0.3677×32768 ≈ 12,050` counts。
- **设备路径上有两条"本该读到零"的路**（⇒ 现象与二者之一有关）：`internal/audio/wasapi_windows.go:394-403` 逐字 `	if silent || data == nil {` / `		return make([]int16, frames)`（WASAPI 打 `AUDCLNT_BUFFERFLAGS_SILENT` 时给全零）与 `:363 floating := s.format.tag == waveFormatFloat` ⇒ **要么设备真的在流 ≈12,000 counts 的非零样本，要么 `convertPacket` 那一支把字节切错了**。
- ★**`convertPacket` 的射程＝今天没有任何仪器跑过它**（尺＝`grep -rn "convertPacket" --include=*.go internal cmd` ⇒ **3 命中全在 `wasapi_windows.go`**：调用点 `:385`、定义 `:393-394`，⛔ 零 `_test.go`）。⇒ **它唯一的生产调用点在真设备上**，所以"切错字节"这一形**在仓内不可能被测到**。⚠ 但它依赖的两枚函数**都有单元测试**（`internal/audio/resample_test.go:158 TestMonoDownmixAndFloatConvert`，用到 `:160 MonoDownmix(stereo, 2)` 与 `:169 FloatToPCM16(f)`）⇒ 所以本票要问的是**枚/通道的拼装那一层**（`frames*channels` 的切片、`ch` 从哪来、`bits==32` 与非浮点 tag 同时出现时怎么办），⛔ 不是"把 `FloatToPCM16` 重测一遍"。
- **本机那次窗口的旁证（不是结论）**：常驻腿那发（票 290 `AC#3` 形②）里 `dropped_frames_total` 在两扇开门窗口内分别涨 `1→704`（23.3 s）与 `722→1073`（11.8 s），与 32 ms/帧对得上 ⇒ **采集线程确实在跑、帧确实在产出**；那一发没有电平读数可看（常驻路径零打印器，只有退出行）。
- ⚠ **不能排除"设备真相"那一支**：本机麦克风增益/增强（Boost、AGC、降噪）或一枚虚拟音频设备，都可能让安静房间真有 −8.6 dBFS 的本底；**这一支要靠"同一枚设备、同一把尺、已知幅度"的夹具来排除，不能靠读码排除**；而"已知幅度/理论 RMS"那把尺今天不存在（票 291 的 B③/B④ 两格，件 §0 表已具名）。
- ⛔ **不许把机主的扬声器没响当作已测事实**：红句自己给了这一候选（逐字 `:195` 末段 `(speakers off, or routed to an endpoint this microphone cannot hear)`），但本次窗口**没有单独问他"外放你听见了吗"**⇒ 这是一个**未答的事实问句**，不是结论；若第 ① 形（声音根本没到麦）成立，那 0.37 的底仍需第 ②/③ 形解释，**两支互相独立**。

## 要建什么（⛔ 先判别，再决定谁修；三格都不许动票 291 那把尺的定义）

- [ ] **AC#0 把"喂进电平尺的到底是哪串字节"打出来（只产**测试侧**仪器，⛔ 不动产码语义）**：新增一枚**真设备**在程用例（可 `t.Skip`，但⛔ 不许 skip 掉本票要的那一发——它归编排者本机跑），要求它**同时**打印：ⓐ 解析后的 mix format 四元组（tag／channels／rate／bits，逐字来自 `internal/audio/wasapi_windows.go:175-180` 那枚 `waveFormat`），ⓑ **第一枚真包的前 16 个样本值**（`int16` 逐枚），ⓒ 同一段字节按"当作 float32 头 8 个"重读一遍的 8 枚值，ⓓ 该窗的 `FrameLevel`。判据＝四行都有；**并自带一条判别句**：若 ⓑ 呈"大/小交替且 |值|≈12,000" ⇒ 指向切错字节；若 ⓑ 看起来像平滑小信号而 ⓒ 像垃圾 ⇒ 形状反过来；若 ⓑⓒ 都平滑且 RMS 真≈0.37 ⇒ **设备本底**那一支。**完成判据＝这枚件在改任何产码前先跑一发，读数逐字进 `probes/297/a1/`。**
- [ ] **AC#1 用已知幅度钉一次"能不能分开"**（⛔ 这一格是票 291 B③/B④ 的**同族**，⛔ 本票不重造尺）：走**注入接缝**（`C8 AudioSource`＝wav，仓里现成：`internal/audio/wavinjector.go:52`、`internal/audio/capturelevel_windows_test.go:19 TestAC247RealCaptureLoopEmitsLevels` 那枚"真采集环＋假 opener＋真 wav 字节"的台件）灌两窗**已知幅度**（如 0.05 与 0.4 的同长窗），要求 `mean` 之比≈幅度之比。⇒ 判据＝**同一把尺在注入面能分开、在真麦面分不开** 这条成对读数拿到手；拿到之后才允许说"问题在设备/转换那一侧，不在 `LevelOfSamples`"。⚠ 这一格会**同时**补掉票 291 的一条具名欠（B③ 的"理论 RMS"那一半），归口写在件里、⛔ 不许顺手勾票 291 的框。
- [ ] **AC#2 三形择一处置（⛔ 由 AC#0/#1 的读数裁，不由写腿挑）**：
  - **①声音根本没到麦**（外放在另一枚端点／扬声器 mute）⇒ 这是**仪器与流程**问题，落点＝把票 247 那发 `sound` 窗改成"点名端点并把端点名打出来"，⛔ 不改电平尺定义；
  - **②字节切错**（`convertPacket` 那一层的 frames/channels/bits 拼装）⇒ 这是**产码缺陷**，落点＝`internal/audio/wasapi_windows.go:385-403` 同处＋**AC#0 那枚件升级为常驻判据**（今天它零仪器，见现量第 5 条）；
  - **③设备本底就是这样**⇒ 那"响得起来"那一格**不能靠 mean 绝对值判**，必须靠**基线相减或比值**——这一支会改到票 291 那把尺的**判据形状**（⛔ 不是默认值、⛔ 不是 2x 那个数本身），属**功能/仪器决策**，要摆机主一句话（含"不做"栏：不做的形状＝球拿到一条 0.37 的常数底、永远不对声音起反应）。
  判据＝三形各带"读数支持它的那一行逐字"＋"落点文件:行"。
- [ ] **AC#3 默认档一字不动**：`mic_muted_default`（`internal/config/schema.go:277`）与本机 `config.toml`（现量：该节第 5 行逐字 `mic_muted_default = true`）**⛔ 一个字都不改**；⛔ 不许为了让电平好看而动增益/重归一/加压缩（`internal/audio/level.go:40-48` 那段具名写死了"曲线的决定权在消费侧、且这枚数仍是 raw scale 的定义"）。
- [ ] **AC#4 门禁＋越界**：`GOFLAGS= go build ./...` rc=0；`sh scripts/d22scan.sh` rc=0；`gofmt -l` 空；`$(go env GOPATH)/bin/gofumpt.exe -l` 空（⛔ 裸 `gofumpt` rc=127 不算跳过）；`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ ./internal/audio/ -count=1` 改前改后各**≥2 发取交集**、逐名作差＝新增红 0 枚（⚠⛔ **改前那两发起手前必跑 `tasklist //FI "IMAGENAME eq wisp.exe"`＝0**：本仓今日实测一枚活着的常驻探针会把整包名册搅出假红，而 AC#0/#1 恰好要起真进程）；每把门禁件自落一行 `rc=N`（⛔ 0 字节＝那格没交）；`git show --stat` 名册只含本票自认的件＋`probes/297/**`；`frontend/**`／`design/**`／三枚冻结件／golden／`thresholds.go`／`allowlist.txt`／D43 表零字节；⛔ 零 push、commit 必带显式 pathspec、临时件只建不删（README 规则 8）。

## 禁区

- ⛔ **不许放宽 `resident_audio_247_live_windows_test.go:194` 那枚判据**，也不许把 `:195` 那句红句改软——它是本票唯一的手柄；那两行逐字＝`:194` `	if loud.mean() <= quietA.mean()*2 || loud.max() <= quietA.max() {`（⚠ **两半**：mean 要 2 倍、max 要更大）＋ `:195` `		t.Fatalf("the loud phase is not distinguishable from silence: %s | %s | %s"+`；要改判据形状**必须**先有 AC#2 ③ 的裁定＋机主一句话。
- ⛔ 不许顺手勾票 247（`AC#2`/`AC#4`/`AC#6`）、票 290（`AC#3`）、票 291（`AC#2`）任何框——三张票的读数归非实现者验收腿裁，本票只交付料。
- ⛔ 不许新造协程（D38(b) 名册零膨胀）、⛔ 裸 `go func(`（ban #1）、⛔ 在 `internal/audio` 里读配置或碰球。
- git：只 commit 不 push；⛔ `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件、不建 worktree。

## 排程

`297-a1`（＝AC#0 那枚仪器＋AC#1 的注入面成对读数，⛔ 零产码语义改动；写面＝新增 `*_test.go`＋`probes/297/**`）→ 编排者按 AC#2 三形裁 → 落点腿（①／②任一）或机主清单（③）。
⛔ 与票 296 **同属地界但不同问题**（一枚是热键的键值合并、一枚是电平字节），⛔ 不许混批；`297-a1` 要真设备 ⇒ **必须独占 `cmd/wisp` 与音频设备面**（同一次窗口只许一枚腿碰麦），排在 `292` 的验收与复跑之后。

**Status:** **未开工**，但**现象已定实**（六窗读数在盘上、`LevelOfSamples` 零形有钉、`convertPacket` 零仪器三条都现跑核过）。⛔ 零翻框、零 push。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
- [2026-10-09 18:07:16 +08] agent=编排者 did=立票 297（真机六窗读数＋三形候选，LevelOfSamples 零形有钉、convertPacket 零仪器两条现跑核过）next=派 297-a1（AC#0 仪器＋AC#1 注入面对拉，独占音频设备面）
