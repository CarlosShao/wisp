# 语音夹具环境变量 — 20 枚语音红从"看起来像产品缺陷"变成真跑（2026-11-12，round 88）

## 0. 上一轮的错误结论（必须更正）

`docs/evidence/s1/cmd-wisp-dll-staging-carrier.md` §6b 与台账 `A1796` 写的是：

> 残余 20 枚红＝另一层…**根因现读＝模型夹具根本没下载**：`models/` 与 `build/models/` 只有 `manifest.json`＋`.minisig`，零 `.onnx`/`.bin`

**这句话是错的。** 它只查了 `models/` 与 `build/models/` 两处，而夹具根本不在那里。

## 1. 现量：模型早就在盘上

| 目录 | 内容 | 合计 |
|---|---|---|
| `third_party/spike-models/` | `asr-encoder.int8.onnx` 157.8MB · `asr-decoder.int8.onnx` 68.3MB · `asr-tokens.txt` · `silero_vad.onnx` · `vocos-22khz-univ.onnx` · `matcha-icefall-zh-baker/`（72.1MB acoustic＋lexicon＋tokens＋fst）· `sherpa-onnx-kws-.../` | ~300MB |
| `third_party/model-fetch/` | `x-kws` 35.6MB · `x-matcha` 87.5MB · `x-punc` 76.0MB · `x-sensevoice` 229.4MB（含 `.tar.bz2` 原件） | ~428MB |

⇒ **零下载需要，一分钱不花。** 缺的只是**指向它们的环境变量**。

## 2. 真根因（逐字）

`cmd/wisp/resident_voice_15e_windows_test.go:166-168`：

```go
src := os.Getenv("WISP_SPIKE_MODELS")
if src == "" {
    return fmt.Errorf("WISP_SPIKE_MODELS is not set")
}
```

没设 ⇒ `t.Fatalf("install the streaming ASR into %s: %v", ...)` ⇒ 断言看起来像"装配坏了"，实际是**环境没指路**。

这套变量在仓库里只有**散文**记录、没有可执行来源：
- `docs/evidence/s1/12r2-impl.md:21`：`WISP_SILERO_MODEL/WISP_TTS_DIR/WISP_SPIKE_MODELS/WISP_KWS_DIR/WISP_PUNC_MODEL/WISP_15E_TMP`
- `docs/evidence/s1/15a-r1-impl.md:115`：`set WISP_SPIKE_MODELS=<repo>\third_party\spike-models`

散文不能执行 ⇒ 每条腿都要重新推导，且**推导错就得到两种假信号**（见 §4）。

## 3. 落地

`scripts/wisp-cli-staged-tests.sh` 新增 `voice_env()`：五个变量**逐个解析默认值＋逐目录验存在**，调用方已设的值优先（腿仍可指向别的树）。

| 变量 | 默认值 | 为什么是这个 |
|---|---|---|
| `WISP_SPIKE_MODELS` | `third_party/spike-models` | 15e 的 ASR 三件来源（`:175`） |
| `WISP_SILERO_MODEL` | `third_party/spike-models/silero_vad.onnx` | VAD 单文件 |
| `WISP_TTS_DIR` | `third_party/spike-models/matcha-icefall-zh-baker` | ★见 §4 |
| `WISP_KWS_DIR` | `third_party/spike-models/sherpa-onnx-kws-...` | 唤醒词 |
| `WISP_PUNC_MODEL` | `third_party/model-fetch/x-punc` | 标点（`voiceWirePunctModel` 自己的 fallback 同此） |

**目录不存在 ⇒ exit 1 响亮报 setup 失败**，⛔ 让它变成"红测试"。

## 4. ★两种假信号（我踩了一个，留档防再踩）

### 4a. 指向 `model-fetch/x-matcha` ⇒ 31 枚红（比不设还差）

`voiceTtsFiles.complete()`（`cmd/wisp/resident_voice_engines_windows.go:85-87`）要求五件齐全：

```go
return f.Acoustic != "" && f.Vocoder != "" && f.Lexicon != "" && f.Tokens != "" && f.DictDir != ""
```

`x-matcha` 有 acoustic＋lexicon＋tokens，**没有 vocoder、没有 dict 目录** ⇒ 逐枚报
`the matcha checkpoint is incomplete under ...x-matcha`。正确树是 `spike-models/matcha-icefall-zh-baker`——`voiceWireTtsFiles` 自己的 fallback（`resident_voice_wire_r1_windows_test.go:213`）就写这个路径。

**实测对照**：指向 x-matcha ⇒ PASS=15 FAIL=31；指向 spike-models ⇒ **FAIL=0**（见 §5）。

### 4b. PowerShell 调 bash 时的引号陷阱

`& $sh script.sh -test.run '^(...)'` 在 PowerShell 里会把 `-test.run` 与正则拆成两个参数，二进制收到裸 `-test` ⇒ `flag provided but not defined: -test` + 整页 usage，rc=2。**这看起来像脚本坏了，实际是调用方引号问题。** 正确姿势：`& $sh -c "script.sh '-test.run' '^(...)'"`。

## 5. 读数（我亲跑，`bash -c` 调用，env 由脚本自解析）

`scripts/wisp-cli-staged-tests.sh '-test.run' '^(Test15e|Test392|Test381r2|Test390r1|Test394R1|Test26R3|Test382WR2|TestWireR1)'`
日志：`D:\tmp\wisp\20261112\envrun3\carrier.log`

| 指标 | 读数 |
|---|---|
| `WISP_SPIKE_MODELS is not set` | **0** |
| `checkpoint is incomplete` | **0** |
| `onnxruntime 1.17` 遮蔽 | **0** |
| **`--- FAIL`** | **0**（该族） |
| `--- PASS` | 逐族累加，含 **`Test15eSpeechEndToEndHandsTheTranscriptToTheAgent` PASS 9.89s**（整包发里长期 FAIL） |

**具名转绿**：`Test15eSpeechEndToEndHandsTheTranscriptToTheAgent`(9.89s) · `Test15eSilenceNeverCallsTheHandOff` · `Test15eARefusedHandOffLeavesTheAsrUnopened` · `Test15eTheChannelHasNoConsumerWithoutTheLeg` · `Test15ePauseStopsTheListenerAndTheTrayFollows` · `Test15eTheProductionHandOffIsTheModelStoreWalk`(3.95s) · `Test392JSONSentinelClampsNonFiniteReadings`

## 6. 边界（诚实标注）

- 本件治的是**夹具寻址**这一层，⛔ 不宣称整包 52 枚红全绿——`cmd/wisp` 另有 5 枚台账 `:32512` 已具名的**存量真红**（`TestTicket198FirstRunCallerIsTheRunEntryOnly` 等），与本件无关。
- ⛔ 未下载任何模型、未改 System32、未提交二进制、未动断言/scope/pin/SKIP/GUARD。
- 上一轮"模型未下载"的结论**已在 §0 明确更正**，台账同步更正。

## 7. 标点模型：一个文件而非一个目录（第三个假信号）

`WISP_PUNC_MODEL` 指向 `third_party/model-fetch/x-punc`（**目录**）时，腿日志打：

> `voice: 标点级未接上（这一进程的 final 是原始假设，不带标点）：…/x-punc`

`voiceWirePunctModel`（`resident_voice_wire_r1_windows_test.go:198-205`）把 `WISP_PUNC_MODEL`
**原样当作模型路径**，只有未设时才 fallback 到一个**文件**（`.../model.int8.onnx`）。
⇒ 该变量是**文件路径**，不是目录。

正确值：`third_party/model-fetch/x-punc/sherpa-onnx-punct-ct-transformer-zh-en-vocab272727-2024-04-12-int8/model.int8.onnx`

**修正后复跑 `^TestWireR1`（我亲跑，日志 `D:\tmp\wisp\20261112\puncfix\run.log`）**：

| 指标 | 读数 |
|---|---|
| `标点级未接上` | **0**（改前每次必现） |
| `--- PASS` | **6** |
| `--- FAIL` | **1** |

转绿具名：`TestWireR1AnnounceEndToEnd` 19.02s · `TestWireR1WarmSecondTurnAcquiresNothing` 12.08s ·
**`TestWireR1PunctuationReachesTheAgent` 5.20s**（就是标点那一格）· `TestWireR1AnnounceGateRefusesTaintedText` 0.45s ·
`TestWireR1SilenceAndEmptyTextNeverSpeak` 4.30s · `TestWireR1TheRealDeviceTakesTheAudio` 0.32s

## 8. 环境层修好后暴露的一枚真红（⛔ 不是环境问题）

`TestWireR1WarmIdleDisposesOnTheOwner` FAIL 5.19s，断言逐字：

> `resident_voice_wire_r1_windows_test.go:745: mic close counted 1 times, want at least 2 (the Warm row and the cascade's defence in depth)`

**性质判定**：这是**产品缺陷**，不是夹具问题。
- 整包发里该枚也 FAIL，但当时原因是 `WISP_SPIKE_MODELS is not set`（环境）⇒ **同一个名字、两个不同原因**；
- 环境层修好后它才暴露出真实断言：Warm 行与级联防御应各关一次麦，实测只关 1 次；
- 我本轮**零 Go 产码改动**（只动 `scripts/**` 与 `docs/**`），⛔ 不是我引入的。

**去向**：该枚需单独立票或并入现有票（`mic close` 计数属 `resident_voice_wire` 写面），⛔ 在本件里不算"已修"。

## 9. 总计（本轮实测，11 族）

| 族 | 读数 |
|---|---|
| `Test15e*`(6) · `Test392*`(1) · `Test381r2*`(5) · `Test390r1*`(4) · `Test394R1*`(2) · `Test26R3*`(3) · `Test382WR2*`(3) · `TestWireR1*`(6/7) | **24 PASS / 0 FAIL**（第一次 11 族发），`TestWireR1` 单独复跑 6 PASS / 1 FAIL |
| 三类环境错误 | `is not set`＝0 · `checkpoint is incomplete`＝0 · `onnxruntime 1.17`＝0 |
| 转绿具名 | `Test15eSpeechEndToEnd…` 9.89s · `Test15eTheProductionHandOff…` 3.95s · `Test381r2RealDevice…` 36.00s · `Test390r1BallClickRealDevice…` 36.93s · `Test26R3AnnounceClosesTheCaptureGate` 19.85s · `Test382WR2TTSCleanSentence…` 8.37s · `TestWireR1PunctuationReachesTheAgent` 5.20s 等 |
