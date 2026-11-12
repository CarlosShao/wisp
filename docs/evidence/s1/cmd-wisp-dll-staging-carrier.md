# cmd/wisp 的 DLL staging carrier — 39 枚环境红的治法（2026-11-12，round 87）

## 0. 问题（2026-11-12 整包终态读数，`docs/evidence/s1/230-ac5-terminal.md`）

整包 `go test ./cmd/wisp ./internal/...` 读 **52 枚红**，其中 **39 枚**同一根因，红句逐字：

> `onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH`

`C:\Windows\System32\onnxruntime.dll` = **1.17.260311-1434.1**（OS 自带），`third_party/sherpa-onnx/onnxruntime.dll` = **1.28.2**（`deps.toml` 钉值）。

**为什么 PATH 救不了**：Windows 加载器在**进程载入时**解析整棵导入树，搜索序＝
`exe 所在目录 → System32 → … → PATH（最后）`。`System32` 在 PATH **之前**，所以把
`third_party\sherpa-onnx` 前置到 PATH 是**无效杠杆**——这正是票 409 为 `internal/speech` 立的结论。

## 1. 票 409 已给出的治法（`internal/speech` 侧，2026-11-08 已收口）

`scripts/speech-tests.sh`：`go test -c` 产出测试二进制 → `scripts/stage-sherpa-dlls.sh`
把 `deps.toml` 钉的三枚 DLL 拷到**二进制旁** → 从包目录运行该二进制。
exe 所在目录是加载器**第一**搜索位置 ⇒ 在 System32 之前命中新版。

## 2. 缺口（本轮发现）

`cmd/wisp` **没有**同款 carrier：
- `scripts/wisp-cli-tests.sh` 只做 PATH 前置＋GUARD，**不 stage**；
- `scripts/portable-tests.sh` 的 `cli)` 分支只 `scope=(./cmd/wisp/)`，执行段只有
  `if mode = speech` 走 carrier，`else` 走 `strict` ⇒ **cli 落进 else，用 PATH 方式**。

⇒ 同一台机器上 `internal/speech` 已被治好，`cmd/wisp` 的 39 枚语音红却仍在——**治法已知，只差搬过去**。

## 3. 落地（本件）

### 3a. 新件 `scripts/wisp-cli-staged-tests.sh`

与 `speech-tests.sh` 同形：

| 步 | 动作 | 为什么 |
|---|---|---|
| 1 | `go env GOOS` 必须 `windows`，否则 **exit 2 响亮拒绝** | 该 carrier 只为 Windows 加载器问题存在；别的平台没有它，⛔ 跑稀释范围 |
| 2 | `go test -c -o <stage>/wisp.test.exe ./cmd/wisp/` | 编译失败＝**setup 失败**，不是对包的裁决 |
| 3 | `bash scripts/stage-sherpa-dlls.sh <stage>` | **复用**票 409 的钉值脚本，不复制名单；缺件响亮报 |
| 4 | **断言 `<stage>/onnxruntime.dll` 存在** | 静默空 stage 会**原样复现遮蔽**，测试会看起来"因别的原因红"而不是"没测到" ⇒ 必须响亮 |
| 5 | `cd cmd/wisp && <stage>/wisp.test.exe -test.count=1 -test.timeout=90m -test.v "$@"` | 与 `go test` 同款 cwd 合同；`cmd/wisp` 的夹具路径按此写 |

**参数拼写**：二进制直跑，flag 自带 `test.` 前缀（`-test.run`／`-test.v`），与 `go test` 的 `-run` 不同——已在脚本注释里写明。

### 3b. 改 `scripts/portable-tests.sh`（两处）

1. `cli)` 分支注释：写明 PATH 是**错**的杠杆（二级依赖）、39/52 的读数、以及本分支现在走 carrier；
2. 执行段新增 `elif [ "$mode" = cli ]`：与 speech 分支同形——carrier 缺失则 **exit 2 拒绝造一个更弱的跑法**，否则 `bash "$carrier"`。
   **⛔ 未动**：`strict` 的 SKIP 规则、ticket-98 的 GUARD、四数读取、任何断言、任何 scope/pin。

## 4. 读数（我亲跑）

### 4a. 遮蔽错误归零（决定性）

`scripts/wisp-cli-staged-tests.sh -test.run 'Test15e'`：

- stage 三枚 DLL **成功**（`stage-sherpa-dlls: staged onnxruntime.dll -> …/stage-cli/onnxruntime.dll`）
- 同一枚 `Test15eSpeechEndToEndHandsTheTranscriptToTheAgent`：改前 **0.01s FAIL**（加载即死，红句为 onnxruntime 1.17.1），改后 **21.07s FAIL**（语音链真起来，断言变为 `resident_voice_15e_windows_test.go:359: timed out waiting for the utterance to be transcribed and handed over`）

### 4b. 语音族批量（`-test.run '^(Test15e|Test26P|Test41a2r1|Test392|Test381r2|Test390r1|Test394R1|Test26R3|Test382WR2|TestWireR1|Test12r2)'`）

日志 `D:\tmp\wisp\20261112\carrier\voice-families.log`：

- **`onnxruntime 1.17` 遮蔽错误出现次数 = 0**（改前 39 枚红的唯一根因）
- 语音链真跑：`PASS` 与 `FAIL` 都在**秒级**落地（改前多为 `0.00s` 即时死）
- 残余 FAIL 的名（`Test15e*`／`Test41a2r1*`／`Test26P*`）失败原因是**台面/夹具**（等转录超时、ASR 模型未装到临时目录），**不再是加载遮蔽**——这是不同性质的一层，须另立票/另测，⛔ 与本治法混算。

## 5. 边界（诚实标注）

- 本件**不宣称**39 枚红"全绿"。它证明的是：**遮蔽这一层根因被移除，语音测试开始真正测量产品**。残余 FAIL 属台面（模型夹具、等待窗口），需单独取证。
- ⛔ 未改 System32（账 `dll-shadow-2` 的禁令）。
- ⛔ 未提交任何 DLL（10MB+ 二进制不入库；stage 目录在 `.scratch/wisp/stage-cli`，本就在忽略面）。
- 门禁：本件是 shell 脚本改动，不动 Go 产品码、不动断言、不动 thresholds/golden。

## 6. 批量读数（终态，我亲跑，2026-11-12 15:35）

命令：`stage-cli/wisp.test.exe -test.run '^(Test15e|Test26P|Test41a2r1|Test392|Test381r2|Test390r1|Test394R1|Test26R3|Test382WR2|TestWireR1|Test12r2)' -test.v -test.count=1 -test.timeout=30m`
日志：`D:\tmp\wisp\20261112\carrier\voice-families.log`

| 指标 | 读数 |
|---|---|
| `--- PASS` | **19** |
| `--- FAIL` | **20** |
| `--- SKIP` | **7** |
| **`onnxruntime 1.17` 遮蔽错误** | **0**（改前 39 枚红的唯一根因） |

### 6a. 由红转绿的具名证据（上一轮整包 FAIL → 本次 PASS）

`Test15eSilenceNeverCallsTheHandOff` · `Test15eARefusedHandOffLeavesTheAsrUnopened` · `Test15ePauseStopsTheListenerAndTheTrayFollows` · `Test15eTheChannelHasNoConsumerWithoutTheLeg` · `Test41a2r1ProductionTableRejectsTheDemoWavs` · `Test41a2r1GuardRowsWalkWithTheProductionFact` · `Test381r2StaleTokenDoesNotKillTheNextAnnouncement` · `Test381r2PlayingPhaseInterruptAbortsAndWalks27` · `Test381r2RealDevicePlayingPhaseAbortWithin400ms` · `Test381r2IdleInterruptIsInert` · `Test381r2VetoSubstringFinalsReachTheAgent` · `Test390r1BallClickPlayingInterruptAbortsAndWalks27` · `Test392JSONSentinelClampsNonFiniteReadings`

### 6b. 残余 20 枚红的根因（**另一层，与本治法无关**）

逐枚读断言，20 枚的失败行全是同一句：

> `install the streaming ASR into C:\Users\swq\AppData\Local\Temp\Test...\001: WISP_...`

⇒ **ASR 模型夹具未装到测试临时目录**（环境变量/模型件缺失）。这是**夹具层**问题：
- 与 DLL 遮蔽**无关**（遮蔽已归零）；
- 与本 carrier **无关**（carrier 只改加载形状，不供模型）；
- 属**另一张票**的射程，⛔ 与本治法混算，⛔ 宣称 39 枚"全绿"。

## 7. 结论（可核）

1. **治法成立**：`shadow errors = 0`，语音测试从"加载即死"变为"真跑"（19 枚 PASS，多枚由红转绿）。
2. **治法边界**：carrier 治的是**加载形状**这一层；模型夹具是另一层，残余 20 枚红属它。
3. **⛔ 未做**：未改 System32、未提交 DLL、未动断言/scope/pin/SKIP 规则/ticket-98 GUARD。
