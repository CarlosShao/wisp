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
