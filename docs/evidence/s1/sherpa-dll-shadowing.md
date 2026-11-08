# 票 409 裁决表 — System32 的 onnxruntime.dll 遮挡 PATH 上的新版（语音族 30 红）

- 日期：2026-11-08
- 腿：Wisp 编队写腿（票 409）
- 范围：`scripts/stage-sherpa-dlls.sh`（新）· `scripts/speech-tests.sh`（新）· 本文件。零 Go 产品码、零断言改动、零 thresholds/golden 改动。

## 1. 根因（现测，非转述）

| 件 | 位置 | 版本 | 大小 |
|---|---|---|---|
| 旧 | `C:\Windows\System32\onnxruntime.dll` | 1.17.260311-1434.1.os-germanium.6bcd20d | 10,572,960 B |
| 新 | `third_party/sherpa-onnx/onnxruntime.dll` | 1.28.2 | 17,799,168 B |

Windows 加载器标准搜索序：**exe 所在目录 → System32 → … → PATH（最后）**。`internal/speech` 以 cgo 链接 `sherpa-onnx-c-api.dll`，其二级依赖 `onnxruntime.dll` 在**进程载入时**解析——System32 的 1.17.1 赢，PATH 上的 1.28.2 输。`checkOnnxruntime()`（`internal/speech/vad.go`）读版本拒收旧版，把每个模型背书测试转成 skip：

```
silero VAD unavailable in this process: speech: onnxruntime 1.17.1 is older than
the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is
shadowing the one on PATH (on Windows the loader searches System32 before PATH)
```

实测（2026-11-08，HEAD `eec09a028`，裸 `go test ./internal/speech/ -count=1`）：
**RUN=106／PASS=42／FAIL=0／SKIP=47**，rc=0——假绿形状。整包口径下 `tools/d22scan/runtests.sh` 把 SKIP 记 fatal（票 71 AC#3：「SKIP is not a pass」），即整包语音族红的机制：**不是断言红，是遮蔽把真跑挡成 skip、runner 把 skip 记红**。

## 2. 被否的修法（各一句为什么）

1. **`TestMain`/`init()` 显式 `LoadLibraryEx` 绝对路径**：太晚。cgo 静态导入树的解析发生在进程入口之前、任何 Go 代码之前，此时 onnxruntime.dll 已按搜索序绑定到 System32 旧版。
2. **PATH 前置 `third_party/sherpa-onnx`**（票 98 在 `scripts/wisp-cli-tests.sh` 的既有配方）：对一级导入（`sherpa-onnx-c-api.dll` 本身）有效，对二级依赖无效——System32 在 PATH 之前。本票基线（PATH 已前置）仍 47 skip，实锤。
3. **改 System32**：系统目录，不可逆动作，禁。

## 3. 落地修法：测试二进制与三枚 DLL 同目录（staging）

加载器搜索序第一步就是 exe 所在目录——唯一能在载入时、在 System32 之前生效的杠杆。这也是 Windows 应用的标准部署形状（DLL 与 exe 同目录），与本仓产线 `scripts/build.ps1` 的形状同构（`build\wisp.exe` + 三枚 DLL 同目录）。`internal/speech/fixtures_test.go` 的 `sileroModelPath` 逐字预留了这一形状（「run from a staging directory that holds the sherpa DLLs」时向上找 repo root）。

两个新件：

- `scripts/stage-sherpa-dlls.sh DEST`：从 `deps.toml` 的 `[sherpa-onnx.dll."…"]` 钉值读 DLL 名（与 `wisp-cli-tests.sh` 同一合同，不复制名单），GUARD 校验后 `cp` 到 DEST。幂等，setup 失败响红不静默。
- `scripts/speech-tests.sh`：carrier。`go test -c -o <stage>/speech.test.exe ./internal/speech/` → stage 三枚 DLL 到 exe 旁 → 在包目录（`go test` 同款 cwd 合同）运行二进制 → 打 runtests 同形四数。`WISP_SPEECH_STAGE_DIR` 可改 stage 根，默认 `.scratch/wisp/stage-speech/`。

复跑命令：

```bash
unset GOTMPDIR TMPDIR TEMP TMP
export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH
bash scripts/speech-tests.sh
```

## 4. 验收读数（实测）

- 修前（裸 `go test ./internal/speech/ -count=1 -timeout 20m -v`）：**RUN=106／PASS=42／FAIL=0／SKIP=47**，rc=0（skip 假绿；skip 理由 100% 为 onnxruntime 1.17.1 遮蔽文案）。
- 修后（carrier 真跑，2026-11-08 20:57–20:59 本腿实测，HEAD `75ff9c4eb`）：**RUN=106／PASS=83／FAIL=0／SKIP=6，rc=0，全程 147 秒**（`bash scripts/speech-tests.sh`，log 存 `.scratch/wisp/stage-speech/speech-tests.log`）。
- 作差（修前 → 修后）：**PASS 42→83（+41）／SKIP 47→6（−41）／FAIL 0→0／RUN 106→106**。41 枚遮蔽 skip 全部转为真跑 `--- PASS`（含 ASR/VAD/TTS/transcriber/cgo 边界全族）；log 中 `1.17.1` 提及次数＝**0**——遮蔽在本机口径下归零。
- 残留 6 枚 skip（逐名，均与遮蔽无关）：`TestPunctuationOnThe300CharDictationFixture`／`TestPunctuationIsRestoredPerExample`／`TestPunctuatorPositiveControls`／`TestPunctuatorLatencyWithinBudget`／`TestFinalTranscriptCrossesTheSeamPunctuated`／`TestPunctuationOnTheRealAsrFinal`——理由全部为**标点模型缺件**（需 `WISP_PUNC_MODEL` 指向 CT-Transformer `model.int8.onnx`，或放 `third_party\spike-models\model.int8.onnx`）。修前它们同样在 47 之列，但当时在更早的遮蔽检查处就被挡下；遮蔽治好后它们推进到标点模型查找处、因模型文件不在盘上而 skip。修前口径「skip 理由 100% 为 1.17.1」对这 6 枚应精化为「被 1.17.1 遮蔽挡在标点检查之前」。标点模型取件（`third_party/model-fetch/punc-int8.tar.bz2`）是另一件事，不在本票射程。
- 载入验证：`cgo_boundary_windows_test.go` 用 `EnumProcessModules` 枚举进程载入镜像并断言来源路径——staging 后这些断言以真跑（非 skip）通过，即进程内载入的 onnxruntime 来自 stage 目录的新版。

## 5. 边界与如实

- 裸 `go test ./internal/speech/`（不经 carrier）在这台装有 System32 旧版的机器上仍会 skip——这是 Windows 加载器的物理决定，任何 Go 侧代码都改变不了载入时解析。carrier 是本仓对该物理事实的应对，与票 98 的 `wisp-cli-tests.sh` 同为「step 载体」形状。
- **遮蔽的射程＝本机独有，CI 上不存在**（现核 `.github/workflows/ci.yml`）：CI 对 `internal/speech` 只挂了纯 Go 子包 `internal/speech/cer`（ubuntu，零 cgo，`ci.yml:754-764`），cgo 父包根本不在 CI 射程；GitHub runner 的 System32 也不携带 onnxruntime.dll。故本票修的是**本机门禁天天红的口径**，不是 CI 红。
- 本票不翻任何票面 AC 勾/Status，只动 Progress log；不动 `frontend/**`、`internal/ball/**`、`cmd/wisp/**`、`PLAN.md`、specs、thresholds、golden；未触碰 System32。
