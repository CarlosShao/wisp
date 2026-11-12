# 428 — `Test26PWithoutAECTheEchoTranscribesVerbatim` 恒红：不是"判定器没开 episode"，是**判定器从未 armed**（台面问题）

- 票面：`.scratch/wisp/issues/428-the-26p-no-aec-echo-test-sees-zero-episodes-so-the-judge-never-received-the-announcement-audio.md`
- 写面：**仅** `cmd/wisp/resident_voice_bargein_26p_windows_test.go`（台面侧；产码零字节）
- 实现者：写腿 428-26p（裁决者≠实现者，翻框归编排者）

---

## §1 AC#0 — 只读取证（本格零写；探针走仓外 `-overlay`）

### 1.1 稳定性：**稳定红，3/3**

```
$ D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^Test26PWithoutAECTheEchoTranscribesVerbatim$'"
run1: --- FAIL: Test26PWithoutAECTheEchoTranscribesVerbatim (22.86s)
run2: --- FAIL: Test26PWithoutAECTheEchoTranscribesVerbatim (24.42s)
overlay-probe run: --- FAIL: Test26PWithoutAECTheEchoTranscribesVerbatim (23.81s)
```

⇒ **稳定红，不是偶发**。三发都落在同一句、同一读数（见 §1.3），耗时 22.86s / 24.42s / 23.81s
（量级一致，方差来自 TTS 合成与 waveOut 播放的墙钟）。与票面 round-88 的 16.78s 同族同句。

### 1.2 `:216` 完整格式串（票面说日志里被截断，这里逐字抄源码）

```go
215: 	if len(read.Episodes) == 0 {
216: 		t.Fatalf("MUTANT 26p-M2 SHAPE: the raw echo opened zero episodes - the judge did not see the announcement audio the mic surface carried (the 59-r1 rig measured 5/5 verbatim on this shape)")
217: 	}
```

注意这句**本身就把结论写错了一半**：它断言"judge did not see the announcement audio"，
但本程量出来的事实是**judge 根本没起来**（不是"起来了但没看见"）。这是 §3 修形的直接依据。

### 1.3 ★岔路判定：**"音频没送到判定器"（台面问题）**，不是"送到了但判定器没开 episode"

**判据＝仓外 overlay 探针**（拷到 `D:\tmp\wisp\20261112\428-26p\probe\`，`go test -overlay` 编译；
本枚是运行时行为检查 ⇒ overlay 有效）。探针插在 `:214` 之前，读 `readings()` 的**全部**字段：

```
428-PROBE armed=false reason="[voice] aec.enabled=false：AEC 插话检测不装配（无 AEC 的插话必自触发，P15/A1307；D47 Path C 由配置关闭）"
          frames_drained=0 processed_samples=0 episodes=0
          aec_bypass=true aec_available=false ready_samples=0 lag=0 lag_strength=0.000
          triggers=0 stream_capped=false far_samples=0
428-PROBE tts spoken=2 aborted=0 samples=281536
```

**逐字段读**：

| 读数 | 值 | 含义 |
|---|---|---|
| `armed` | **false** | 判定器**从未 arm** —— 这是根因，不是症状 |
| `reason` | `[voice] aec.enabled=false…` | arm 被**配置开关**拒绝，拒绝原因自报 |
| `frames_drained` | **0** | drain goroutine 一帧都没读（它压根没被 spawn） |
| `processed_samples` | **0** | 没有任何样本进过判定器 |
| `episodes` | **0** | 因此必然零 episode |
| `tts samples` | 281536 | **播报确实出声了**（2 句都播完，`spoken=2`） |

⇒ **判定器实际收到 0 帧**（不是"很多帧但 Episodes 空"）。按票面判据：
**收到 0 帧 ⇒ 前者＝台面/注入问题**。票面的第二个可能（判定器行为问题）被**排除**。

**为什么不是"注入没送到"**：注入侧是健康的 —— `tts samples=281536` 证明播报链路完整，
且 §1.4 的修复后读数（`frames_drained=200`）证明**同一条注入路径在判定器 armed 后立刻收到 200 帧**。
⇒ 断点在 **arm**，不在 `rig.inject`。

### 1.4 ★同一探针的对照发：把 `AECEnabled` 声明成 true ⇒ 判定器立刻活

同一个探针，唯一差别是给 rig 的 `voicePlan` 补上 `AECEnabled: true`（§3 的修形）：

```
428-PROBE armed=false reason="" frames_drained=200 processed=102080 episodes=1
          aec_bypass=true aec_available=false ready=9728 lag=501
M2-VERBATIM episodes=1 passed=0 triggers=0 transcript="我已经完成明天下午三点的会议改到四点请准时参加谢谢" distinctive=true
```

`voice: barge-in monitor armed for the announcement window roster=tts-infer aec=false aec_bypass=true dynamic_gating=true`
这行日志在修复前**一次都没出现**，修复后出现两次（两段播报各一次）。

⇒ **200 帧进判定器、开出 1 个 episode、逐字转写命中**（`distinctive=true`）。
**这一发把"是 arm 而不是注入"钉成了读数**，不是推断。

### 1.5 音频从注入到读数的完整路径（逐段读源码，⛔ 非猜）

```
rig.inject(echoWav)                          resident_voice_bargein_26p_windows_test.go:183-194
  → audio.NewWavInjector(wav) / inj.Start(ctx, r.frames)        （注入 C8 seam 帧通道）
  → rv.bargein.drain() 读 m.frames                              resident_voice_bargein_windows.go:741
      m.framesDrained.Add(1)                                    :745      ← 探针的 frames_drained 就是这里
      audio.DecodeFrame → d.pending                             :746-768
      d.armJudge()  需要 far≥300ms 且 pending≥600ms             :784-861
          → audio.AlignByEnvelopeCorrelation 量 lag             :813
          → audio.NewBargeinGate(...)                            :837/:843
      d.process() 每 160 样本（10ms）喂一次                      :866-928
          res = nearBlock（aecBypass：raw capture）              :900
          m.det.Accept(...) → EndpointStart → d.judge(ep)        :916-921
              → d.gate.Evaluate(...)                            :990
              → m.episodes = append(...)                        :1010-1011  ← 探针的 episodes 读这里
  → rig.rv.bargein.readings()                                  resident_voice_bargein_windows.go:638-677
```

**唯一的开关在入口**：`m.arm()` 的第一个 case 就是 `case !m.aecEnabled:` ⇒ 直接 return false
（`resident_voice_bargein_windows.go:392-393`、`:405-406`）。**arm 返回 false 之后 `drain` 从未被 spawn**
（`:482-484` 在 arm 内部），所以 `frames_drained` 恒 0、`episodes` 恒空。

### 1.6 三个岔路的另外两条

**(2) 这枚测试读哪些 env？有没有 `WISP_26P_*` 一类别名变量没设？**

读 `new26pRig`（`:137-168`）与 `voiceWireStart`（`resident_voice_wire_r1_windows_test.go:226-243`）：

- `new26pRig` 的 `voicePlan` 只填三格：`VadModelPath` / `Source` / `Tts`（`:139-143`）。
- env 面走的是 carrier 已解析的五个变量（`WISP_SPIKE_MODELS` / `WISP_SILERO_MODEL` /
  `WISP_TTS_DIR` / `WISP_KWS_DIR` / `WISP_PUNC_MODEL`），由 `scripts/wisp-cli-staged-tests.sh:136-173` 解析。
- **没有 `WISP_26P_*` 一类别名变量**（`git grep 'WISP_26P'` 零命中）⇒ 这条岔路**排除**。
- **真正缺的不是 env，是 `voicePlan` 的一个字段**：`AECEnabled`。它是**编译期零值 `false`**，
  不是环境变量 —— 所以"检查 env 是否齐"这个方向查不到它。

**(3) 与票 392 的 AEC 调参有关吗？**

**无关，但同根**。票面说"本轮 `Test392*` 两枚 PASS" —— 本程现读 round-88 日志，
那两枚其实是 **SKIP**（不是 PASS）：

```
=== RUN   Test392KeepalivePostureFindingsAndColdDefault
--- SKIP: Test392KeepalivePostureFindingsAndColdDefault (3.89s)
=== RUN   Test392AecTuningDoubleTalkSwing
--- SKIP: Test392AecTuningDoubleTalkSwing (3.69s)
```

skip 原因是**同一句** `the AEC3 binding is not in this build (-tags webrtcap …)`
（`resident_voice_bargein_392_windows_test.go:404`/`:554`/`:647` 的 `rig.requireAEC3()`）。

- 票 392 的三枚**也用 `newVoice26pRig`**（`:403`/`:553`/`:646`）⇒ **同一个 `AECEnabled` 零值把它们也一起废了**。
  它们因为 `requireAEC3()` 先 Skip（本机无 webrtc 源码树，`third_party/webrtc-audio-processing` 不存在），
  **Skip 把同一个台面缺陷掩盖了** —— 所以"判定器在别的用例里是活的"这个前提**在本机不成立**：
  本机没有任何一枚 26p/392 集成用例真正跑到过判定器。
- 本枚 `Test26PWithoutAEC…` 是**唯一不 requireAEC3** 的（它的 bypass 半不建 APM），
  所以它是**唯一把这个缺陷暴露出来的那枚** —— 这也解释了为什么"只有它红"。
- 结论：**与 392 的 AEC 调参语义无关**（那是 `aec.Tuning`），但**与 392 的台面同根**（同一个 rig）。

---

## §2 根因（一句话 ＋ `文件:行`）

**commit `e83494cf5`（2026-11-11，180-census leg K）给 `newResidentVoice` 加了
`rv.bargein.setAECEnabled(plan.AECEnabled)`（`cmd/wisp/resident_voice_windows.go:591`），
而 `newVoice26pRig` 手搓的 `voicePlan` 从不设这一格（`cmd/wisp/resident_voice_bargein_26p_windows_test.go:139-143`），
零值 `false` 让 `arm()` 的第一个 case 直接拒绝（`cmd/wisp/resident_voice_bargein_windows.go:392-393`），
判定器整族再也没 armed 过。**

这条 commit 的 `git show` 只动了产码、**没动 rig**（`git show e83494cf5 --name-only | grep 26p` 零命中），
所以是一次**产码新增一个开关、台面没跟上**的回归。

**为什么 `AECEnabled` 的零值恰好是"关"**：schema 里 `[voice] aec.enabled` 的默认是
`default:"true"`（`internal/config/schema.go` 的 `AECConfig.Enabled`），
所以**生产路径永远填 true**，零值 false 只在**手搓 plan** 的测试里出现 —— 这正是它藏得住的原因。

---

## §3 AC#1 — 选形：**修台面（测试侧）**，产码零字节

### 3.1 改了什么

```go
	plan := voicePlan{
		VadModelPath: voice15eVadModel(t),
		Source:       "26p test",
		// Ticket 428: a plan built by hand must say "on" the way the boot plan
		// does for the [voice] aec.enabled schema default (true). Since
		// e83494cf5 newResidentVoice feeds plan.AECEnabled into the monitor's
		// setAECEnabled, and setAECEnabled(false) refuses EVERY arm - so the
		// zero value here silently disabled the whole 26p/392 integration
		// family (the M2 reading saw frames_drained=0 and read it as "the
		// judge saw nothing"). Same shape ticket 423 ruled for PunctEnabled
		// (resident_voice_wire_r1_windows_test.go:286-288).
		AECEnabled: true,
		Tts:        voiceWireTtsFiles(t),
	}
```

**一行语义改动 ＋ 六行注释**，只碰 `newVoice26pRig` 的 plan 字面量。断言**一个字节都没动**
（`:215-217` 的红句原样保留）。

### 3.2 为什么修台面而不是修产码（选形理由）

1. **判据已定量**：§1.3 探针读到 `frames_drained=0` 且 `reason` **自报**是配置开关 ⇒ 断点在台面输入，
   不在产码逻辑。产码的 `setAECEnabled` 语义是对的（`aec.enabled=false` 就该拒绝装配，P15/A1307 的既定裁定）。
2. **产码侧没有可改的东西而不违反既定裁定**：要让这枚测试过，要么把 `setAECEnabled(false)` 的拒绝去掉
   （＝推翻 P15/A1307 与 `Test180KAECSwitchGatesTheMonitor`），要么给测试开后门 —— 两者都是改产品语义。
3. **票 423 已有同形先例且已裁定**：`PunctEnabled` 那格同样是"手搓 plan 必须说出生产默认值"，
   当时的结论是**台面补齐**（`resident_voice_wire_r1_windows_test.go:286-288` 的注释逐字：
   "a plan built by hand must say `on` the way the boot plan does for the config default"）。
   本程**沿用同一裁定**，不新造规则。
4. **产码与台面二选一**（票面硬约束）：既然判"台面"，产码就零字节 ⇒ 突变控制才能证明"这一行承重"。

### 3.3 修复后读数（carrier 真跑，两发）

```
$ D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^Test26PWithoutAECTheEchoTranscribesVerbatim$'"
postfix-run1: M2-VERBATIM episodes=1 passed=0 triggers=0 transcript="不已经完成明天下午三点的会议改到四点请准时参加谢谢" distinctive=true
              --- SKIP: Test26PWithoutAECTheEchoTranscribesVerbatim (26.94s)
postfix-run2: M2-VERBATIM episodes=1 passed=0 triggers=0 transcript="不已经完成明天下午三点的会议改到四点请准时参加谢谢" distinctive=true
              --- SKIP: Test26PWithoutAECTheEchoTranscribesVerbatim (24.37s)
```

**★重要且必须说清**：修复后这枚**不是 PASS，是 SKIP** —— 红句消失了（`episodes=1`，
`M2-VERBATIM` 那行打出来了），但测试继续往下走，在**还回 AEC 的那半**（`:249-250` 的
`rig2.requireAEC3()`）Skip 掉，原因逐字：

```
resident_voice_bargein_26p_echo_windows_test.go:250: the AEC3 binding is not in this build
(-tags webrtcap + the deps.toml CGO flags): the integration readings need it:
aec: webrtc-audio-processing binding not built (requires cgo + -tags webrtcap; see deps.toml [webrtc-audio-processing])
```

⇒ **本机能证明的是"判定器活了、200 帧进、开 1 episode、逐字转写复现"**（M2 的前半，也就是票面关心的那半）；
**不能证明的是"还回 AEC 后零 episode"**（restore 半）—— 那半需要 `-tags webrtcap` 的构建，
而本机 `third_party/webrtc-audio-processing` **源码树不存在**（`Test-Path` 实测 MISSING），
`go test -tags webrtcap` 直接编译失败（`api/scoped_refptr.h: No such file or directory`）。
**这是本机 staging 的缺口，不是本票的缺陷**，具名记在 §7。

---

## §4 AC#1 突变控制：两发**真跑**

**仪器**：`go test -overlay`，拷贝全在仓外 `D:\tmp\wisp\20261112\428-26p\mut\`。
本枚是运行时行为检查 ⇒ overlay 有效（对照票 427 那枚读盘源码的静态检查，overlay 是假绿仪器）。

### 4.1 发①：把修好的那行删掉（`AECEnabled: true` ⇒ `false`）⇒ **必红，实测红**

```
resident_voice_bargein_26p_echo_windows_test.go:216: MUTANT 26p-M2 SHAPE: the raw echo opened zero
episodes - the judge did not see the announcement audio the mic surface carried (the 59-r1 rig
measured 5/5 verbatim on this shape)
--- FAIL: Test26PWithoutAECTheEchoTranscribesVerbatim (21.75s)
rc=1
```

⇒ **逐字复现原红句**，耗时 21.75s（与改前 22.86s/24.42s 同量级）。**这一行承重，坐实。**

### 4.2 发②：恢复 ⇒ **绿（红句消失）**

```
resident_voice_bargein_26p_echo_windows_test.go:245: M2-VERBATIM episodes=1 passed=0 triggers=0
transcript="不已经完成明天下午三点的会议改到四点请准时参加谢谢" distinctive=true
--- SKIP: Test26PWithoutAECTheEchoTranscribesVerbatim (27.26s)
rc=0
```

⇒ 红句消失、`episodes=1`、逐字转写命中。**恢复即绿（本机口径下"绿"＝越过 `:216`，见 §3.3 的 SKIP 说明）。**

### 4.3 两发读数对照表

| 发 | `AECEnabled` | `armed` | `frames_drained` | `episodes` | 判决 | 耗时 |
|---|---|---|---|---|---|---|
| 改前（基线） | 未设（false） | false | **0** | **0** | **FAIL** `:216` | 22.86s / 24.42s |
| 发①（删修复） | false | false | **0** | **0** | **FAIL** `:216` | 21.75s |
| 发②（恢复） | true | — | **200** | **1** | 越过 `:216`，后 SKIP | 27.26s |
| 修复后 carrier 两发 | true | — | **200** | **1** | 越过 `:216`，后 SKIP | 26.94s / 24.37s |

---

## §5 AC#2 — 越界检查（`git diff --name-only` 逐行）

```
$ git diff --name-only
.scratch/wisp/issues/426-warm-idle-mic-close-counted-once-against-two-because-the-second-defence-never-runs.md
.scratch/wisp/issues/427-the-12r2-ast-check-still-looks-for-marksettling-inside-openscope-but-it-moved-into-scopeoptions-hooks-logf.md
.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go
.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt
cmd/wisp/resident_voice_bargein_26p_windows_test.go
frontend/fixtures/composer-states.html
```

**逐行归属**（⛔ 藏他腿的活）：

| 行 | 是谁的 | 判据 |
|---|---|---|
| `cmd/wisp/resident_voice_bargein_26p_windows_test.go` | **本腿（唯一写面）** | 本程唯一编辑（§3.1 那一行） |
| `.scratch/wisp/issues/426-…`、`.scratch/wisp/issues/427-…` | 他腿（426/427 收件改名） | 本程起手（20:26）首次 `git status` 即在，本程零触碰 |
| `.scratch/wisp/probes/291/…`、`.scratch/wisp/probes/orch/…`、`frontend/fixtures/composer-states.html` | 起手前既存脏改 | 同上，本程零触碰 |
| `cmd/wisp/resident_watchdog_windows_test.go` | **他腿 429（在飞）** | 本程 20:47 撞到它编译失败（`"os" imported and not used`），等其落定后重跑；本程 ⛔ 未触碰、⛔ 未 commit |

- 本腿**产码零字节**：`resident_voice_bargein_windows.go` / `resident_voice_windows.go` / `internal/**` 未改。
- 禁面全部未动：`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／
  `tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`。
- `build/mut_rvw.go` 未删（未跟踪残留件，非本腿写面）。
- commit 用**显式 pathspec**（只 `cmd/wisp/resident_voice_bargein_26p_windows_test.go` ＋ 本证据件），⛔ push。

---

## §6 AC#3 — 门禁四读数（本程现跑）

| 尺 | 命令 | 读数 |
|---|---|---|
| 构建 | `go build ./cmd/wisp/` | **rc=0**（按票面用 `./cmd/wisp/`，⛔ 不用 `./...`，避开 `build/mut_rvw.go` 残留） |
| 静态检查 | `go vet ./cmd/wisp/` | **rc=0** |
| 格式 | `gofumpt -l cmd/wisp/resident_voice_bargein_26p_windows_test.go`<br>`gofumpt -l cmd/wisp/` | **两尺均零输出**（rc=0） |
| D22 禁令 | `scripts/d22scan.sh`（⚠ ⛔ 用 `go run ./tools/d22scan`，独立 module） | **clean，rc=0**：`d22scan: clean - no D22 ban violations`；`runtests.sh: OK - packages=[./...] top-level: PASS=45 FAIL=0 SKIP=0` |

**载波口径**：测试一律经 `D:\work\soft\Git\bin\bash.exe` ＋ `scripts/wisp-cli-staged-tests.sh`；
本机 PATH 上的 `bash` 是 WSL 转发器（一枚测试都不跑），本程未用。

---

## §7 AC#4 — 与票 26P／392 的关系（⛔ 不翻任何框）

1. **本票只处理这一枚**：`Test26PWithoutAECTheEchoTranscribesVerbatim` 的 `:216` 红句从红转绿。
2. **⛔ 不翻票 26P／392 任何框**；⛔ 不宣称 26P/392 全绿。
3. **★必须回填的事实**（本程只记，⛔ 自行改票）：修好的是**共享 rig** `newVoice26pRig`，
   所以 `Test26P*` 另外四枚与 `Test392*` 三枚**也一并被解封**了 —— 但它们在本机**仍然 SKIP**
   （`requireAEC3()`，本机无 webrtc 源码树）。⇒ **它们的真实状态在本机仍然未测**，
   需要一台带 `-tags webrtcap` 的 staging 才能判。**这条是票 26P/392 的欠账，不是本票能关的。**
4. **票面"本轮 `Test392*` 两枚 PASS"这条前提，本程实测是 SKIP**（§1.6(3) 逐字日志）。
   本程 ⛔ 改票面，只在此具名：**该前提与 round-88 日志不符**，请编排者核。
5. **AEC 调参语义（`aec.Tuning`）本程零触碰**；本票与它无关（§1.6(3)）。

### 7.1 具名残余（⛔ 藏）

1. **restore 半在本机不可测**：`-tags webrtcap` 需要 `third_party/webrtc-audio-processing` 源码树，
   本机 `Test-Path` 实测 **MISSING**，`go test -tags webrtcap` 编译失败
   （`api/scoped_refptr.h: No such file or directory`）。⇒ 修复后本枚是 **SKIP**，不是 PASS。
   **"还回 AEC3 零 episode"这一对在本机未验证**，需要带 webrtc 的 staging 复跑。
2. **`Test26P*` 另外四枚 + `Test392*` 三枚本程未跑**（票面明示不碰，且本机必 SKIP）。
   它们是否因同一修复转绿，本机无法回答（同上，缺 webrtc）。
3. **本程未跑整包 `cmd/wisp`**：本机整包有他腿在飞（429 的 `resident_watchdog_windows_test.go` 20:47 还在改，
   本程 20:47 撞到过一次编译失败）。全包读数不归本腿声明。

---

## §8 停手条件自核

1. **未定义即停？** 未触发：票面已给定"台面/产码二选一、⛔ 两边都改"，本程按 §1.3 的探针读数判"台面"
   并写明理由（§3.2），无自行发明判据。
2. **要不要改产码？** 不要：`setAECEnabled` 的拒绝语义是 P15/A1307 的既定裁定（且有
   `Test180KAECSwitchGatesTheMonitor` 钉住），改它＝推翻裁定。
3. **要不要翻框？** 不要：翻框归编排者，本程只追加 Progress log。
