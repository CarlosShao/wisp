# 票 297 非实现者验收裁决表 v1（297v1 · 只读裁决 · 2026-11-11 18:3x +08）

- 裁决者：非实现者验收腿 `297v1`（⛔ 零产码、零票面改动、零 git 写操作；被验件＝`297-l1`／`297-l2`／`297-roster` 三腿交付）。
- 裁决基线：票面 `.scratch/wisp/issues/297-quiet-level-sits-at-a-near-constant-0-37-so-loud-is-not-separable.md` AC#0–AC#4 ＋ 编排者裁定节 ⓐ–ⓔ ＋ Progress log（297-l1 / 297-l2 / 297-roster 三条）。
- 证据件全部现读（read/glob），未跑任何 `go test`／`go build`（派单禁令）。

## 0. probes/297 盘点（现量 `Get-ChildItem -Recurse`，字节为文件 Length）

| 组 | 关键件 | 字节 |
|---|---|---|
| a1（只读普查，2026-10-10） | `00-anchor-and-rulers.md` / `10-convert-packet-reachability.md` / `20-instrument-shape-options.md` / `30-known-amplitude-vs-ticket-291.md` / `40-existing-nails-and-ci.md` | 8,964 / 15,977 / 14,160 / 12,778 / 12,975 |
| l1（AC#0+#1 落地腿） | `00-readings-297-l1.md`(8,426)＋9 枚 verbatim log（baseline/after audio 56B×4、baseline/after cmdwisp ~303–345KB×4、`live-instrument-run{1,2}.log` 1,285×2、`live-247-sixwindow-after.log` 5,048、`gate-build.log` 602、`gate-build-scoped.log` 19、`gate-d22scan.log` 27,522、`gate-path.log` 798） | 全在 |
| l2（AC#2 形①＋AC#5＋AC#3 腿） | `00-readings-297-l2.md`(8,225)＋`formone-gated-run{1,2,3}.log`(424/5,957/2,703)＋`gates.log` 27,761＋`ac5-lists/`（8 件，含两枚 linux 测试二进制 6.28MB/6.50MB）等 | 全在 |
| l3（CI 名册收口腿） | `00-roster-readings.md`(4,971)＋`l2-roster-verify-fullwindows.log`(544,077)＋`no-script-13pkg-probe.log`(526,997)＋`l2-line-{head,verbatim}.txt`(2,662×2) | 全在 |

⇒ **l1/l2/l3 读数件全部真在、非 0 字节**（AC#4"0 字节＝那格没交"条款对三腿均成立）。

## 1. 逐格裁决

| AC | 裁决 | 判据（读数支持行，全部现读） |
|---|---|---|
| **AC#0** 把喂进电平尺的字节打出来（测试侧仪器＋真机一发，读数落盘） | **PASS** | 三枚仪器真在：`internal/audio/convertpacket_297_windows_test.go`（S1–S5 五面＋parseWaveFormat 四元组）、`livepacket_297_windows_test.go`（`WISP_LIVE_MIC=1` 门，打印 ⓐⓑⓒⓓ）。真机读数 verbatim 落盘：`probes/297/l1/live-instrument-run1.log` 逐字 `[A] mix format quad: tag=3 channels=2 rate=48000 bits=32` ＋ `[B] … first 16: [0 0 … 0]` ＋ `[C] … first 8: [0 0 … 0]` ＋ `[D] 64 windows … mean=0.003513 max=0.078400` ＋ `[verdict] no signature matched`（run2 同形 mean=0.002175）。四行判据全齐。⚠ 一处偏差（不影响判）：票面原文写读数进 `probes/297/a1/`，实际落 `probes/297/l1/`——a1 已被 2026-10-10 只读普查腿占用、落地腿改名 l1 并在读数档第 1 条具名说明，属具名改名非丢料。 |
| **AC#1** 已知幅度钉"能不能分开"（注入面成对读数） | **PASS** | 仪器真在：`amplitudepair_297_windows_test.go`（真采集环＋假 opener＋真 wav，0.05/0.4 两窗，判据 mean 比＝幅度比 ±0.01）。读数 verbatim（`probes/297/l1/00-readings-297-l1.md` §4）：`window A (amp 1638 counts) mean=0.049988, window B (amp 13107 counts) mean=0.399994` ＋ `mean ratio = 8.0018, amplitude ratio = 8.0018 -- the ruler SEPARATES at the injection seam`。成对另一半（真麦面）＝票 247 修复前 1.0402（票面现量）＋修复后真机底 0.0013–0.0035（`live-247-sixwindow-after.log` verbatim）同在盘上。⇒ 同一把尺注入面分得开、真麦面（修复前）分不开的成对读数到手；`LevelOfSamples` 被赦免。⛔ 未勾票 291 任何框（l1/l2 Progress log 与读数档均零翻框）。 |
| **AC#2** 三形择一处置（由 AC#0/#1 读数裁） | **PASS** | 三形各有逐字读数与落点：**形②**（字节切错＝旧 0.37 底成因，已被票 300 @24 修复治掉）——合成切错形状 `level = 0.383045 (RMS 12552 counts)`（convertpacket 用例 D2 行）↔ 旧底 0.367691（12,050 counts）↔ 修复后真机 `mean=0.003513/0.002175`（l1 run1/run2 log），三读数同页对齐；**形③**（设备本底 0.37）＝**排除**——同一 Realtek 麦同一把尺修复后底 0.0013–0.0035（−50 dBFS 量级），`[verdict] no signature matched`；**形①**（声音没到麦）＝**SUPPORTED**（l2 两次独立真机）：`formone-gated-run3.log` 逐字 `[render 1..5] endpoint tap: rms=640.055832 … captured=48000`（端点喂声铁证 5/5）＋ `[volume] session volume=1.000 muted=false` ＋ `endpoint volume=0.200 muted=false`（软件可读轴全开）＋ 麦渲染覆盖窗 mean 0.001288–0.001959 全 ≤ 安静底 0.002108（best 比 0.93x）＋ `[verdict] … form ① SUPPORTED this run`；旁证＝麦听得见房间突发声（quiet max=0.121012）却听不见渲染（覆盖窗 max ≤0.005050）⇒ 麦路径本身活的。落点：①仪器/流程支——端点点名已落（`livepacket` A0 行打印 capture/render 端点名、247 用例 `:156` 同），不改电平尺定义；剩余嫌疑＝物理层（扬声器物理静音/没接/路由），**机主耳朵定**——这正是形①自身措辞内含的收口，不构成"形①未定案"。与派单预置结论（l2 形① SUPPORTED）一致，读数支撑成立。 |
| **AC#3** 默认档一字不动 | **PASS** | `internal/config/schema.go` `MicMutedDefault … default:"true"` 真读在盘（现量 :299，票面引 :277 为行号漂移、值逐字未动）；本机 `%APPDATA%\wisp\` 无 config.toml（默认档即生效档；票面"config.toml 第 5 行"是立票时点状态，l2 如实记录、非被腿改动）；`level.go` raw-scale 块（无加性底、无增益/压缩）与 `internal/ball/liquid.go:30` `SilenceLevelGate = 0.06` 原样；l2 `git diff --stat` 四文件＝空。⇒ 零改动坐实。 |
| **AC#4** 门禁＋越界 | **PARTIAL** | 已达：audio/aecspike 改前改后各 2 发全绿（l1 18.08–19.08s、l2 audio 18.0–18.8s＋aecspike 0.19–0.32s，新增红 0）；cmd/wisp 门 **l1 达标**（改前改后各 2 发红名册同集＝TestTicket198…/Test313…/TestLockedSuffix…33r11 三枚在册历史红，交集作差新增红 0）；d22scan rc0 clean（458/462 文件）、path budget GREEN、gofmt/gofumpt 空、vet rc0、wisp.exe 起手＝0、每把门禁件自落 rc 行（spot-check `gate-build-scoped.log`＝`SCOPED_BUILD_RC=0`）且无 0 字节门禁件；l3 名册腿整档实跑两发（SKIP 根已除、Test52 红坐实为 D46 fixture 预存 flake）。**两处未全达（如实记）**：Ⓐ 票面要求 `GOFLAGS= go build ./...` rc=0，两腿整树 build 均 rc=1——归因 gitignored 未跟踪遗留分片 `build/mut_rvw.go`（预存本机状态，scoped build `./cmd/... ./internal/... ./tools/...` rc=0，两读数均落档未抹）；Ⓑ l2 那一波 cmd/wisp 门不可判（423 槽主树在变＋两发撞 600s 包超时，l2 自记漏 `-timeout 30m`）——l1 已判达标、l2 写面（aecspike）经 `go list -deps -test ./cmd/wisp` 证明不在 cmd/wisp 导入图上，结构上零影响。⇒ 按票面字面 rc=0 与"改前改后各≥2 发"全格标准，这两处各欠一角，判 PARTIAL 不判 FAIL（欠角均预存/外因、如实落档、无放宽无隐瞒）。 |

## 2. 越界检查（现读坐实）

- `resident_audio_247_live_windows_test.go:194/:195` 判据与红句一字未动：`live-247-sixwindow-after.log` 里红句逐字照打（`the loud phase is not distinguishable from silence: … (speakers off, or routed to an endpoint this microphone cannot hear)`），判据仍生效（该发 --- FAIL 即判据咬合的证据）。
- 三腿 Progress log 自报与读数档交叉核对：AC#0..AC#4 框全未勾（勾框权归本裁决，本表不代勾票 291/247/290 任何框）。
- l1 自报纪律偏差（--no-verify、调用数超配）已由编排者在票面裁定节 ⚠ 记账；l2 自报 cmd/wisp 门漏 `-timeout 30m` 记过——均为自曝非隐瞒。
- 本裁决零产码、零票面、零 Status 改动；唯一写面＝本证据件。

## 3. 裁决汇总

**AC#0 PASS · AC#1 PASS · AC#2 PASS（形① SUPPORTED，物理层待机主耳朵）· AC#3 PASS · AC#4 PARTIAL**（整树 build rc=1＝gitignored 预存遗留＋l2 波 cmd/wisp 门不可判，欠角如实落档）。

票 297 的"今天判不了那一格"已可判：旧 0.37 底＝形②（float32 按 int16 切，票 300 修复已治），形③排除，现红句成因＝形①（声音没到麦，软件轴全开，机主耳朵定）。
