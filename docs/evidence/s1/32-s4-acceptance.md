# 32 — S4 切片验收（腿 A ＝ 票 32 · 非实现者 · 只读＋重跑）

**日期**：2026-11-15 12:0x +08
**锚**：HEAD `43702bac7eb4b290434baa8ed394d666be35deca`（`dev`）
**被验票**：`.scratch/wisp/issues/32-s4-acceptance.md`（5 格 AC，逐字见 §2／§3）
**前置四票**：`.scratch/wisp/issues/{27,29,30,31}-*-done.md`（盘上都在）
**写面**：本文件 ＋ `.scratch/wisp/probes/32/**`。⛔ 未动产码、票面、台账、阈值、golden、契约、`docs/specs/**`；⛔ 碰 `cmd/wisp`（只读）；⛔ `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`；⛔ push。
**方法**：五格逐格**自己现读**产码与测试，再自己跑能跑的门；票面与前置票的判语一律**当未验证**，只在"它们声称的那件事"上被当作线索，读数一律落回 `file:line` 原文。凡本腿现跑的数标【现跑】，转引的数标【读自】并给出处。
**口径**：判语只用三种——**满足／部分满足／不满足**；环境上根本量不到的格写**不可验＋缺什么**，⛔ 记满足、⛔ 记失败。

---

## 1. 先说清三件事（不然后面五格会被误读）

1. **本机跑不了真语音引擎**（不是产品的红，是这台机器的环境）。`go test ./internal/speech` 对真引擎族整族 SKIP，逐字原因：

   > `punctuation stage unavailable on this host: speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH (on Windows the loader searches System32 before PATH)`

   根因现量：`C:\Windows\System32\onnxruntime.dll` FileVersion `1.17.260311-1434.1.os-germanium.6bcd20d`，而仓内 `third_party\sherpa-onnx\onnxruntime.dll` 是 `1.28.2` ⇒ 加载器先命中 System32。⇒ **任何需要真 ASR/TTS/标点模型的现场复跑，在本机跑不出来**；本报告凡涉及真引擎的读数，都是读别人在**别的机器/别的 staging 目录**上跑出来的证据件，条目里逐条标了出处。
2. **`docs/evidence/s4/` 这个目录不存在**（`Test-Path docs\evidence\s4` => `False`；`docs/evidence` 下只有 `precheck, s0, s1, s2, s3`）。票面 Key constraints 那句「screen recordings + logs + tool traces committed under docs/evidence/s4/」**今天零字节**，而 ③④ 两场真机验收的"证据"正是这一条。
3. **③④ 两场一次都没跑过**。这不是推断：全仓（`cmd`+`internal`，含测试）搜 `五轮`／`5-round`／`five-round` 命中 **0**；`docs` 下 `mp4/webm/mkv` **0 枚**；没有任何一次把「口述→标点→agent→`input.type` 落到目标窗」串起来的 run。**票 32 的 5 格里，没有一格是被"一次完整场景 run"证明的**——下面每格的"部分满足"都只覆盖它自己的分句，不是场景级的 PASS。

---

## 2. 逐格判语（一行五列）

| AC 原文逐字（`.scratch/wisp/issues/32-s4-acceptance.md:36-44`） | 判语 | 依据（`file:line` ＋ 读到的原文） | 可核命令 | 若部分满足／不可验，缺什么 |
|---|---|---|---|---|
| **AC1**（`:36-37`）`Scenario ③ PASS: video + transcript + input.type trace; elevated-window failure message captured; score sheet ≥8/10.` | **部分满足** | ① `input.type` 真机 300 字成立：`docs/evidence/s1/23-input-type-real-window.md:12`「runA（ASCII）：注入 300 字 → 窗内 **300 字**、`wm_char=300`、读取 sha256 与期望**逐字相同**」＋`:59` sha256 `1f47881e…fbf` 同；② 提权拒绝**修后**真机成立：`docs/evidence/s1/23-uipi-end-to-end-after-fix.md:23-25`「`elevated_is_error=true` / `elevated_text=目标程序以管理员身份运行，Wisp 无法向其输入（…本工具不会尝试提权绕过）…` / `elevated_fixture_chars_unchanged=true`」；产码闸在打第一字节之前：`internal/tools/system_backend_windows.go:1438-1440`「`if elevated { return 0, "", ErrElevatedTarget }`」；③ 8/10 有表但是**单样本自评**：`docs/evidence/s1/27-r1-impl.md:19`「`302 字 ⇒ 26 个标点`，逐字对照后我判 **8/10**」 | `go test ./internal/tools -run Test23 -count=1`【现跑 ok 0.152s】<br>`Test-Path docs\evidence\s4`【现跑 False】<br>`go test ./internal/speech -run TestPunctuationOnThe300CharDictationFixture -count=1 -v`【现跑 SKIP，原因见 §1】 | ①**没有"口述→光标"这条产品路**：全仓无听写模式、`internal/agent` 的 system prompt 对 `input.type` **零指引**（`Get-ChildItem internal\agent\*.go | Select-String 'input\.type|听写|打字'` => 0），所以 300 字落在窗内是**探针直接调 tool**（`23-input-type-real-window.md:33`「调真正的 `window.focus` / `system.get` / `input.type` 三个 tool 的 `Execute`」），不是 ASR 转写喂进去的；②**无录像／无转写／无 trace 入库**（§1-2、§1-3）；③ 8/10 是**文本夹具自评**（`internal/speech/punctuate_test.go:259-294` 的 `dictation300` 字符串直接进 `AddPunct`，不是音频过 ASR），D29「human acceptance is the bar」要的是**人**签，本格没人签 |
| **AC2**（`:38-40`）`Scenario ④ PASS: 5-round log; latency table (round1 vs round2 delta >1s proving Warm); TTS ≥7/10; fact-recall demonstrated. **D47: barge-in ≤400ms gate + no-self-transcription PASS; external-audio false-trigger evidence; Conversation memory row backfilled.**` | **部分满足** | ① 延迟表成立（组件级、真引擎）：`docs/evidence/s1/28-r1-impl.md:16`「冷 `acquire=3.2586296s`／…**第一请求 3.4086s**；Warm-1 `acquire=0ns`／**171.35ms** ⇒ **19.9x**；`Stats().LoadsAfterWarm=0`」——**delta＝3.24s > 1s**，且暖的是"不用装"（`:34`「暖的不是"算得快"，是"不用装"」）；② 零自激转写成立：`docs/evidence/s1/26pr1-impl.md:80-83`「`triggers=0、finals Δ=0、handed Δ=0`」＋`docs/evidence/s1/26r3-impl.md:35`「`transcripts_after=1`（相对播报前**零增量**）」；③ 外部音源有 P15 口径读数：`docs/evidence/s1/59-v3-verdict.md:60`「对照票面外部音源目标「≤1 次/10min」：**−30 组合门 0/10min 达标**」；④ barge-in 门在、凡火 ≤400ms：`docs/evidence/s1/26pr1-impl.md:65-69`（99/56/67/65/40ms），但 `docs/evidence/s1/26p-v2-verdict.md:18-22` 复跑 6 窗**只 1 火（24ms）**并「终判＝有条件 PASS：(a) 验收语义按本件 §3 裁定重 key（5 窗至少 1 火＋凡火必 400ms 内）」 | `go test ./internal/session -run TestTicket28Conversation -count=1 -v`【现跑 4 PASS / ok 0.048s】<br>`go test ./internal/audio -run 'TestHotplug\|TestGateConversation' -count=1`【现跑 ok 2.169s】<br>`Select-String docs\SLO.md -Pattern 'Conversation'`【现跑：`:51` 仍 `未测`】<br>`Select-String -Path (生产 .go) -Pattern 'Profiles:'`【现跑：`cmd/wisp` 0 命中】 | ①**五轮日志零字节**（§1-3），本格第一条就没兑现；②**TTS ≥7/10 无人评表**：`docs/evidence/s1/26-r1-impl.md:75`「`P7 表不存在（R26-6 不挡结案）`」，台账只有替代分支 `docs/reports/pending-and-issues.md:51`「[R5] 默认 TTS 走云 TTS，供应商 StepFun」——**票 32 没把 R5 写成 7/10 的替代**，所以这一分句**不可验（要机主双盲人评）**；③**fact-recall 在生产上不成立**：`internal/agent/loop.go:154-156`「Profiles supplies the (3) L1 section; the caller reads them from memory」——而**没有任何生产调用者**给它赋值（`cmd/wisp` 的 `agent.Config` 在 `run.go:2091-2116` 构造，`Profiles` 不在其中；全仓 `ListProfiles` 只被 `internal/memory/privacy.go:50,202` 与 `undo.go:27` 读）⇒ 抽取链在（`internal/memory/l1_warm.go:199 func (w *L1WarmWindow) Fire()`＋`cmd/wisp/resident_voice_engines_windows.go:516 WarmWindow: r.warmWindow`），但**画像永远进不了 prompt**；④ barge-in **不在默认构建里**：`scripts/build.ps1:270`「`go build -trimpath -ldflags … -o … ./cmd/wisp`」**不带 `-tags`** ⇒ `internal/audio/aec/aec_stub.go` 生效（`NewApm` 恒 `ErrUnavailable`），monitor 不武装；≤400ms 那些数是 `-tags webrtcap` 专用构型的读数；⑤ Conversation 内存行**没回填**：`docs/SLO.md:51`「`Conversation`（开麦多轮） | 同 Warm + 采集缓冲 ≤20MB | **未测** | 待 S4（依赖 Warm 态口径裁定）」；`:532` 那行 4.26MB 是 `posture=skeleton`（`:536-538`「六个 state 的 `posture` **全是 `skeleton`**——今天的采样器不加载 KWS/ASR/TTS」），不是 ASR+TTS 共存的数；⑥ 外部音源只在 **P15 rig** 上量过（59-v3），不是 ④ 的产品 run |
| **AC3**（`:41`）`Conversation privacy-confirm-once + red-ring visuals verified by user.` | **不可验（人眼那一半）** | 代码侧本腿自己读通并跑绿：① confirm-once `internal/session/scope.go:527` `func (s *Scope) EnableConversation() (bool, error)`，`:537` `first := !s.confirmed`，`:546` `if !confirm() { … return false, nil }`，首次才 `s.prompts++`；② 红环 `internal/ball/statevisual.go:205-212` `case statemachine.StateConversation:` ⇒ `v.RingColor = p.Danger`、`v.RingWidth = ConvRingStrokePx`、`IconMic`；色值 `internal/ball/tokens.go:141` `Danger: hex(0xE07A70, 1)`（暗色 `:190 hex(0xBC544C,1)`）、线宽 `:577 ConvRingStrokePx = 2.0` | `go test ./internal/session -run TestTicket28ConversationPrivacyConfirmOnce -count=1 -v`【现跑 PASS 0.00s】<br>`Select-String internal\ball\statevisual.go -Pattern 'StateConversation' -Context 0,7` | **票面这一格的动作就是"user 验"**，本腿不能代签：缺**机主一次 Conversation 首开会话**——看两张东西（①一张 L2 隐私卡且第二次不再弹；②球的红环比 Confirming 不弱），签在一句话上即可。代码面已把"该弹一次"与"环是 danger 2px"钉住，但**"屏幕上真的是红的"只有眼睛能判**。（另注：本机真语音链被 §1-1 挡住，不是这格的原因——这格本来不需要引擎，只要起进程开 Conversation。） |
| **AC4**（`:42`）`D10 tier hits during ④ consistent with 30's matrix (no regressions).` | **部分满足** | 矩阵与装配都在、测试我自己跑绿：① 分档 `internal/present/tier.go:81 func Classify(text string) (Tier, string)`，阈值冻结（`:90-97` 用 `longHanziLimit`/`longCharLimit`/`shortLimit`）；② 四档过真管线 `internal/present/e2e_present_30_test.go:70 func Test30FourResultsThroughTheRealPipeline`；③ 生产是**同一条** Dispatcher：`cmd/wisp/run.go:2401-2402`「`if rt.presenter != nil && res.Status == agent.StatusCompleted { rt.presenter.Dispatch(ctx, present.Input{TaskID: res.TaskID, Text: res.Text}) }`」；④ 语音结果的 Speak 通道由同一装配交出：`cmd/wisp/resident_present_windows.go:216-220`「`ch.Speak = presenterSpeaker(run.voiceAnnounce)`」 | `go test ./internal/present -count=1`【现跑 ok 3.980s】<br>`Select-String cmd\wisp\run.go -Pattern 'presenter.Dispatch'` | 「**during ④**」这半句**没有观测**（④ 没跑，§1-3）⇒ 只能说"矩阵与装配一致、无回归可验的码面无回归"，**不能说 ④ 期间分档正确**。缺一次 ④ 里**短结果与长结果各至少一次**的混跑（看 `present` 的 receipt/notice/ball 三通道落点是否与 30 的矩阵逐档对上）。 |
| **AC5**（`:43-44`）`Pre-mortem: TTS device yanked mid-speak; mic device yanked mid-dictation — correct error classes + no stuck states.` | **部分满足** | ① **麦克风侧**：错误类＋具名＋文案＋"确实报出来"都有测试并跑绿——`internal/audio/hotplug_test.go:288-297`「`got != observe.ClassAudioDevice` ⇒ fatal」「`error must name the device`」「`error must carry the in-use guidance`」；文案在 `internal/audio/device.go:159/162`；② **TTS 侧**码路在：无设备 `cmd/wisp/resident_voice_engines_windows.go:1168-1175`「`speaker, serr := rv.newSpeaker(…); if serr != nil { rv.ttsNoDevice.Add(1); rv.noteProblem("播报没有输出设备：" …)`」＋`if moved { … EvInterrupt … }`；播不完 `:1269-1272`「`if drainErr != nil { … noteProblem("播报没播完：" …) }`」；③ **不卡死**（码面）：每个失败分支都从 Speaking 出去（`:1172-1175` → #27；`:1264-1267` → #27；`:1288-1298` → #26/#28） | `go test ./internal/audio -run 'TestHotplug\|TestOpenOccupied' -count=1 -v`【现跑 4 PASS】<br>`Select-String -Path (生产 .go) -Pattern 'EvAudioDeviceLost'`【现跑 4 处，**零 dispatch**】<br>`Select-String cmd\wisp\*_test.go -Pattern 'drain.*err\|播报没播完'`【现跑 0 命中】 | ①**麦克风"口述中被拔"在产品里不出状态**：`EvAudioDeviceLost` 只在 `internal/statemachine/events.go:30` 定义、`table.go:97` 登记（`D43: 14, From: StateListening … To: StateError`），**生产零派发点**（另两处 `cmd/wisp/resident_audio_windows.go:39,358` 是注释与打印）；语音消费侧从不读 `mic.Err()`（`cmd/wisp/resident_voice_*.go` 搜 `.Err()` 零命中）⇒ 设备死了，球可能**停在 Listening/Warm 不动**，这正是票面说的"stuck state"；②**TTS"播放中拔设备"没有夹具**：只有"打开时没设备"与"drain 失败"两条码路，`cmd/wisp/*_test.go` 里没有一枚注入 `Play/Drain` 中途失败的用例；③ **TTS 失败没有 D37 类**：`noteProblem`（`cmd/wisp/resident_voice_windows.go:1563-1572`）只追加散文＋slog，不带 `observe.ClassAudioDevice`（该类在 `internal/observe/errors.go:25`，麦克风侧用了，TTS 侧没用）；④ 真机热拔一次没做过（要硬件＋管理员） |

---

## 3. 逐格详证（表里塞不下的原文与读数）

### AC1 · 场景③

**票面要求（逐字，`:16-18`）**：「continuous 300-char dictation → text appears AT THE TARGET APP's cursor via input.type (saving to a file does NOT count as pass); punctuation subjective ≥8/10; plus one elevated-window attempt → explicit UIPI message」。

本腿**自己跑到的**：

```
go test ./internal/tools -run Test23 -count=1 -v
--- PASS: Test23UIPIPredicateRefusesElevatedTargetBeforeAnyByte (0.00s)
--- PASS: Test23UIPIPredicateLetsLiveNonElevatedTargetThrough (0.00s)
--- PASS: Test23IntegrityOfProcessAnchorsOnThisProcess (0.00s)
--- PASS: Test23ProcessNameIsTheD45BindingKey (0.00s)
ok  github.com/CarlosShao/wisp/internal/tools  0.152s
```

（`Test23UIPIPredicateLetsLiveNonElevatedTargetThrough` **不是 SKIP**：谓词对真普通权限进程这一半在本机也是活的，只有提权那一半需要真 High 窗。）

**读到的真机读数（⛔ 不是我跑的，是 23 那两件的底片）**：

- `docs/evidence/s1/23-input-type-real-window.md:13` **修复前**逐字：「对**真 High 完整性窗口**（Task Manager，IL 12288）调用 `input.type`：**没有** `ErrElevatedTarget`，回的是 `IsError=false` ＋「已向 taskmgr.exe 注入 15 个字符」」——判「**不成立（产品缺陷，已定位到 `file:line` 并可复现）**」。
- `docs/evidence/s1/23-uipi-end-to-end-after-fix.md:28` **修复后**逐字：「**同一个目标窗（`il=0x3000`）：修复前 `IsError=false`＋"已注入 15 个字符"，修复后 `IsError=true`＋冻结拒绝句＋fixture 零变化。**」
- 产码侧的闸：`internal/tools/system_backend_windows.go:1411-1413` 注释＋`:1438-1440` 代码「`if elevated { return 0, "", ErrElevatedTarget }`」，错误值 `internal/tools/system_bridge.go:163-165` 逐字「目标程序以管理员身份运行，Wisp 无法向其输入（Windows UIPI 平台硬限制：普通权限进程不能向更高完整性级别的窗口发送输入，本工具不会尝试提权绕过）。请把目标应用以普通权限运行后重试」。

**为什么仍判"部分满足"**（三条，都是本格自己的分句）：

1. **"video + transcript"零字节**：`docs/evidence/s4/` 不存在，`docs` 下无任何 mp4/webm/mkv。
2. **"300 字口述到光标"这条产品路不存在**。300 字那一发的形状是**探针直接调 tool**（`23-input-type-real-window.md:33`），中间**没有 ASR、没有标点、没有 agent**；而产品侧没有任何"听写模式"把转写接到 `input.type`（`internal/agent` 的 prompt 对 `input.type` 零指引、`cmd/wisp` 无听写入口）。`input.type` 只作为 **L2 工具**存在（`internal/tools/system.go:409` `func InputTypeDecl() Decl { return systemDecl(risk.L2, CapInput, nil) }`），要模型自己决定调、且要过 L2 卡。
3. **8/10 是自评单样本**：`27-r1-impl.md:19` 自己写「我判 **8/10**」，同一件 `:43` 给了扣分理由；而 S4 票面 `:29-30` 明确「human (user) sign-off required for both subjective scores (D29 rule: human acceptance is the bar)」。⇒ **这一分句不可验，缺机主主观评分**。

### AC2 · 场景④

**票面要求（逐字，`:19-25`）**：5 轮、第 2 轮起首字 ≤1.5s（用 round1-vs-round2 delta >1s 证明 Warm）、TTS ≥7/10（P7 表）、记住第 1 轮事实（L1 profile live）、D47 三条（≤400ms 停播／不自激转写／外部音源误触发）、Conversation 内存行回填。

**延迟那一分句是真的**（`28-r1-impl.md:29-32`）：

| 轮 | acquire | 首文本 | 第一请求 |
|---|---|---|---|
| 冷启 | 3.2586296s | 150.0ms | **3.4086319s** |
| Warm-1 | 0ns | 171.35ms | **171.3518ms** |
| Warm-2/3 | 0ns | 189.39/161.19ms | 同 |

delta ＝ 3.24s，> 1s，且 `Stats().LoadsAfterWarm=0` ⇒「Warm 活着」这个断言**在组件级成立**。

**"记住第 1 轮事实"这一分句，生产上是断的**。这是本报告最重要的一条独立发现：

- 注入点存在：`internal/agent/loop.go:154-156`「`// Profiles supplies the (3) L1 section; the caller reads them from memory (ticket 29 owns extraction, this ticket owns only the injection point).`」；`prompt.go:164` 把它拼进 `(3)` 段。
- 但**没有调用者**：`cmd/wisp` 全目录 `Profiles:` **0 命中**；`run.go:2091-2116` 构造的 `agent.Config` 里没有这个字段；全仓生产侧 `ListProfiles(` 只出现在 `internal/memory/privacy.go:50`（`ListPrivacy` 的 `PrivacyProfile` 分支）／`:202`（`ExportPrivacy` 的 `PrivacyProfile` 分支）／`undo.go:27`（撤销）——**没有一处把它喂给 loop**。
- ⇒ 抽取链（`internal/memory/l1_extract.go` ＋ `l1_warm.go:199 Fire()` ＋ `cmd/wisp/resident_voice_engines_windows.go:516 WarmWindow: r.warmWindow`）会把画像**写进库**，但 `(3)` 段**永远是空的**。票 29 自己的 5 格 AC（`29-memory-l1l2-done.md:37-48`：延迟隔离／golden 抽取／LRU／undo／save-recall＋隐私 API）**没有一格钉"注入要有生产者"**，而「L1 injection: full ~400-token profile into prompt suffix (③ section)」写在同一票的 Key constraints（`:22`）⇒ 那一句今天**只是注释里的意图**。本票这一分句因此**不成立**，缺的是把 `ListProfiles` 接到 `agent.Config.Profiles` 的那一根线（以及"每轮刷新还是每任务刷新"的口径）。

**barge-in 三分句的实况**：

- 门与判据在：`internal/audio/gate.go:166-194` `SetConversation`（Path C 开麦窗）、`:273-281`「`if g.speaking && g.path == PathT && !g.conversation`」；组合门 `internal/audio/bargein.go:82 BargeinTightenedEnergyDbfs = -30`、`:420 func NewBargeinGate`；生产挂载 `cmd/wisp/resident_voice_engines_windows.go:1125 rv.bargein.arm()`。
- 读数：`26pr1-impl.md:65-69` 五轮 **99/56/67/65/40ms**（全 ≤400ms）；复跑 `26p-v2-verdict.md:18-22` **6 窗 1 火（24ms）**，其余 5 窗被压到 −39.35..−43.44 dBFS ⇒ 终判「**有条件 PASS**」，语义重 key 成「5 窗至少 1 火＋凡火必 400ms 内」。⇒ **"每轮必火"不成立**，本格不能拿 26pr1 的 5/5 当"门稳定"。
- **默认构建里没有 barge-in**：`scripts/build.ps1:270` 的 `go build` 不带 `-tags webrtcap`；`internal/audio/aec/aec_stub.go:1` 的 build tag 是 `!windows || !cgo || !webrtcap`，stub 的 `NewApm` 恒 `ErrUnavailable`；`cmd/wisp/resident_voice_bargein_windows.go:42`「ErrUnavailable and the monitor stays inert」。⇒ 出厂 `wisp.exe` 里 Path C 插话**不会发生**，那些 ≤400ms 读数都来自 `-tags webrtcap` 的专用构型。
- 不自激转写：`26pr1-impl.md:80-83`「`triggers=0、finals Δ=0、handed Δ=0`」＋`26r3-impl.md:35`「`ECHO-26R3 … transcripts_after=1`（相对播报前**零增量**）」＋夹具 `cmd/wisp/resident_voice_selfecho_26r3_windows_test.go:129 Test26R3AnnounceClosesTheCaptureGate`。这一分句在本机**不可复跑**（真引擎被 §1-1 挡住），但它有独立夹具与两件证据，按"读到的"判成立。
- 外部音源：`59-v3-verdict.md:60`「对照票面外部音源目标「≤1 次/10min」：**−30 组合门 0/10min 达标**、−45 不达标」；同件 `:52-54` 如实记了 music w24 那枚纯回声事件。⇒ **P15 rig 上达标**，但 ④ 的产品 run 没做过、也没有"背景放音乐时产品误触发几次"的读数。

**TTS ≥7/10**：`26-r1-impl.md:75`「`P7 表不存在（R26-6 不挡结案）`」；替代分支是 `pending-and-issues.md:51` R5（默认走 StepFun 云 TTS）。**票 32 没有把 R5 认成 7/10 的替代**，所以这一分句**不可验，缺机主双盲人评表（≥7/10）**。

**Conversation 内存行**：`docs/SLO.md:51` 仍 `未测`；`:127-130` 写着「`Conversation` 预算 = Warm + TTS 增量，**由 S4 实测回填**」；`:532` 那行 4.26MB 的 `posture` 是 `skeleton`，`:536-538` 自己声明不加载 KWS/ASR/TTS ⇒ **不能拿它冒充**。⇒ 本分句**不满足**（既没实测数，也没有"超预算→降级模式"的验证）。

### AC3 · 隐私一次确认 ＋ 红环

代码侧（本腿现读＋现跑）已在上表钉住。补两条原文：

- `internal/session/scope.go:527-560` 的形状：`first := !s.confirmed` → `if first { if confirm == nil { …fail-closed… }; if !confirm() { … return false, nil } }` → `s.confirmed = true` → 只有 `first` 才 `s.prompts++`。`session/scope.go:124-125` 定了 30s/15s 的 idle。
- 红环不是"随便一个色"：`internal/ball/tokens.go:577` 的注释逐字「Conversation ring "不得弱于 Confirming"」，`statevisual.go:205-212` 给 `Danger`＋`ConvRingStrokePx`＋`IconMic`。

**判不可验的原因只有一个**：票面写的是"verified by user"，而本腿是 agent。**缺机主一次开麦会话的人眼签收**（看红环＋第一张隐私卡只弹一次）。

### AC4 · D10 分档

`internal/present/tier.go:76-99` 的阶梯（first match wins：空→none；结构化→structured；>400 汉字或 >800 字符→long；>60 汉字→medium；否则 short）。测试我自己跑绿（`ok 3.980s`）。生产同一条路：`run.go:2401-2402` 只在 `StatusCompleted` 时分发（失败线走 `else` 那条老路：notify ＋ `voiceAnnounce`，`run.go:2387-2399` 的注释逐字「Only a COMPLETED result has a spoken reply; a failure line keeps the legacy pure-text announce it always had」），通道装配 `resident_present_windows.go:204-239`。

照 30 的矩阵逐档对得上（60/61、400/401、800、structured-beats-length），**"无回归"可验的那一半是绿的**；缺的是 ④ 这一次 run 里的实际命中记录。

### AC5 · 拔设备预演

麦克风侧（真测试、本机跑绿）：

```
--- PASS: TestGateConversationWindowOnPathT (1.38s)
--- PASS: TestHotplugReenumerateOnce (0.28s)
--- PASS: TestHotplugStaleHandleFailsLoudly (0.39s)
--- PASS: TestHotplugReenumerateFailureNamesDevice (0.07s)
--- PASS: TestOpenOccupiedAndPermissionDenied (0.00s)
ok  github.com/CarlosShao/wisp/internal/audio  2.169s
```

`hotplug_test.go:257-259` 的自我声明（注释跨三行，此处合写）是「acceptance criterion 4 second half - … the supervisor surfaces Error(audio_device) with the device name instead of silently keeping the stale handle」。

但**产品状态机那一半没接**：`EvAudioDeviceLost` 生产零派发（§2 表里那条尺），`resident_audio_windows.go:358` 自己写着「这一条不推状态、也不拒绝启动」，而语音消费侧从不读 `mic.Err()`。⇒ 口述中途拔麦：**错误类与文案会浮出来（在采集腿的 Stats/Err 上），但球不会走 D43#14，可能停在 Listening**。

TTS 侧：两条失败码路都离开 Speaking（不会卡在 Speaking），但 ① 没有"播放中拔设备"的夹具；② `noteProblem` 只有散文、不带 D37 类；③ `drainErr` 那条**仍然计 `ttsSpoken.Add(1)`**（`:1273`），即"没播完也记成播过一句"——这不是卡死，是账不准，本腿如实记下。

---

## 4. 两段收尾

### ① 满足 N／总 M

**满足 0／总 5**；其中 **部分满足 4 格（AC1、AC2、AC4、AC5）**，**不可验 1 格（AC3，人眼签收那一半）**，**不满足 0 格**（没有哪一格被证伪到"产品做不到"，AC1 的提权拒绝那半曾不成立、已修并复验）。

⇒ 读法：**S4 切片今天拿不到验收**，不是因为某一格判红，而是 5 格**全部**缺"人/现场"那一半：③④ 两场没跑过、`docs/evidence/s4/` 零字节、两个人评（8/10、7/10）都没有、五轮日志不存在。码面能给的已经给了（UIPI 拒绝、Warm 3.24s delta、零自激转写、D10 矩阵、麦热拔错误类），**这些够了"实现"，不够"验收"**。

### ② 四张前置票各承担哪一格（逐票一行）

- **27-punctuation → AC1**（承担 `score sheet ≥8/10` 那一分句）。它自己交的 8/10 是**单样本自评**（`27-r1-impl.md:19,43`），本票要的是人评 ⇒ 它把"标点阶段能做对"证到了，把"人给 8/10"没证。
- **29-memory-l1l2 → AC2**（承担 `remembers one fact from round 1 (L1 profile live)`）。抽取与存储落地（`l1_warm.go:199`、`WarmWindow: r.warmWindow`），但**画像到 prompt 没有生产者**（`agent.Config.Profiles` 全仓零赋值；29 的 Key constraints `:22` 写了"L1 injection"，5 格 AC 里没有一格钉它）⇒ 这一票承担了格子，却承担不了"live"两个字。
- **30-result-routing-d10 → AC4**（承担 D10 矩阵与"无回归"）。矩阵、四档真管线、生产同一条 Dispatcher 都在，本腿跑绿；它不承担"④ 期间"这个观测。
- **31-reminders → 本票 5 格零格承担**。票面 5 条 AC 里没有任何 `reminder.*` 判据（逐字核过 `:36-44`），而 31 又确实列在 `Blocked by` 里 ⇒ **这是本票的结构缺口**：31 的落地只能由它自己票面的 5 格证明，本票（S4 切片验收）**无法替它签收**。要么把 31 从 `Blocked by` 摘掉，要么在 5 格之外补一格 reminder 判据——两条都要改票面文字，本腿不改，抬头报给编排者。

---

## 5. 复跑清单（本腿现跑过的，逐条可重放）

```powershell
cd "D:\work\workspace\projects plans\Wisp"
git rev-parse HEAD                                  # 43702bac7eb4b290434baa8ed394d666be35deca
go test ./internal/present -count=1                 # ok 3.980s
go test ./internal/tools -run Test23 -count=1       # ok 0.152s
go test ./internal/session -run TestTicket28Conversation -count=1 -v   # 4 PASS / ok 0.048s
go test ./internal/audio -run 'TestHotplug|TestOpenOccupied|TestGateConversation' -count=1 -v  # 5 PASS / ok 2.169s
$env:PATH = "$PWD\third_party\sherpa-onnx;$PWD\build;$env:PATH"
go test ./internal/speech -run TestPunctuationOnThe300CharDictationFixture -count=1 -v  # SKIP（onnxruntime 1.17.1 遮蔽，见 §1-1）
Test-Path docs\evidence\s4                          # False
```

空白底的完整输出（含每条命令的原文与 rc）在 `.scratch/wisp/probes/32/readings.txt`。

---

## 6. 本腿的写面与撤销

- 写面：本文件 ＋ `.scratch/wisp/probes/32/readings.txt`（新建）。
- ⛔ 未翻任何 AC 框、⛔ 未改票面/台账/阈值/golden/契约/`docs/specs/**`、⛔ 未动产码、⛔ 未 push。
- 撤销：`git rm --cached docs/evidence/s1/32-s4-acceptance.md && git commit -m "撤销 32 腿 A"`（或直接在下一笔里把这枚文件删掉）；探针目录可整目录留着，它只含读数。
