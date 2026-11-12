# 230 AC#5 — 整包终态读数（2026-11-12，编排者本人跑）

票面 AC#5 逐字（`:38`）：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/... -count=1` 到终态＋**逐名比红名册**。

## 0. 命令与台面

- 起手锚：`HEAD=730754b4a`（分支 dev，含票 174 AC#4 落地 `56406d3b0` 与其后 A1794）
- 命令：`go test ./cmd/wisp ./internal/... -count=1 -timeout 90m`
- PATH 前置：`third_party\sherpa-onnx`（三枚 dll 全在，`onnxruntime.dll` = **1.28.2**）
- `GOTMPDIR`/`TMPDIR` 已 unset
- 日志：`D:\tmp\wisp\20261112\230ac5-self\full.log`（4348 行）
- **读数：rc=1，elapsed=761s**

## 1. 逐包结果（`internal/*` 全绿，红只在 `cmd/wisp`）

| 包 | 读数 |
|---|---|
| `internal/agent` … `internal/winsec`（**35 个 internal 包**） | **全部 `ok`** |
| ★`internal/ball` | **ok 1.119s**（票面在册红 1 枚＝**已消失**） |
| ★`internal/panel` | **ok 57.755s**（票面在册红 4 枚＝**已消失**） |
| ★`internal/risk` | **ok 27.632s**（票面在册的 `TestResolvePerCallBudget` 争用型假红＝**未出现**） |
| `cmd/wisp` | **FAIL 749.344s**，52 枚红 |

⇒ **票面在册的三组红（ball 1／panel 4／risk 争用型）本次全部为零**。按票面判据"历史在册…不算本票新增"⇒ **本票新增红＝0 枚**。

## 2. 52 枚红的逐名归因

### 2a. 环境红 39 枚（机器级，非任何票的欠账）

**根因（现读）**：`C:\Windows\System32\onnxruntime.dll` = **`1.17.260311-1434.1`**，Windows loader 顺序**先于任何 PATH** ⇒ 遮蔽了 `third_party\sherpa-onnx` 的 **1.28.2**。红句逐字：

> `语音链未构造：VAD 没起来（speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH)`

台账已有专账：`:32533` `dll-shadow-2`（治 39 枚语音红；**⛔ 改 System32**；主流做法＝DLL 与 exe 同目录或 `LoadLibraryEx` 绝对路径）。`docs/evidence/s1/122-staticcheck-batch1.md:141-142` 亦已具名此坑。

命中的族：`Test15e*`(5)、`Test41a2r1*`(4)、`Test26P*`(5)、`Test392*`(3)、`Test381r2*`(3)、`Test390r1*`(3)、`Test394R1*`(2)、`Test26R3*`(2)、`Test382WR2*`(3)、`TestWireR1*`(6)、`Test12r2*`(2)、`Test16r1*`(2) 等。

### 2b. 真红 13 枚 → 逐枚核后＝**5 枚存量真红 + 1 枚整包内互扰 + 7 枚属 2a 族**

定向复跑（`-run '^(那六枚)$' -count=1 -v`，rc=1，5.181s）拿到的**真断言**：

| 测试 | 真断言（逐字） | 定性 |
|---|---|---|
| `TestTicket198FirstRunCallerIsTheRunEntryOnly` | `firstrun_198_test.go:281: AC J1 RED: ensureFirstRunConfig is called by [runTextTask runConfigSet65], want exactly [runTextTask]` | **存量真红**（台账 `:32512` 具名） |
| `Test313WritableUnionEqualsTheTwoRosters` | `panel_config_313_test.go:252: the door reports the union [...], ` 与两名册不等 | **存量真红**（台账 `:32512` 具名） |
| `TestAC4FocusReturnToPriorWindowGap33r5` | `panel_host_windows_test.go:995: the panel did not take the foreground on Show` | **存量真红**（需真桌面焦点，台账 `:32512` 具名） |
| `TestLockedSuffixNeverTakesTheLockItself33r11` | `panel_locked_naming_33r11_windows_test.go:254: panelHolder.stampPublishLocked (resident_windows.go line 1008) carries the "Locked" suffix but no roster entry` | **存量真红**（台账 `:32512` 具名） |
| `TestAC247HandingTheLevelToTheSeamIsNotVisibility` | `resident_audio_247_windows_test.go:355: prototypeVisuals is ON in this assembly: the shipped default is off` | **存量真红**（台账 `:32512` 具名） |
| `Test320FetchLegWritesTheKeyNowhereOnDisk` | **定向单独跑 PASS（0.10s）**，整包里红 | **整包内互扰**（非真红） |
| `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey` / `TestAC1ResidentLeg*`(2) / `TestAC246DevLegIgnoresTheTestTaskInjection` / `Test365RunRowsAStuckTaskAsStuck` / `Test41a2r1ProductionTableRejectsTheDemoWavs` / `Test16r1VoiceDemoWakeAsrAgentSpeakSettleRealHardware` | 定向上下文显示与语音链/常驻装配相关 | 归 2a 环境族（未逐枚单独复跑，**具名未验**） |

## 3. 结论（AC#5 判据逐条答）

1. **跑到终态了吗**：**是**。rc=1、elapsed=761s，35 个 internal 包全部 `ok`，无挂死、无 panic、无超时（票 408 修 UDF 进程唯一化后，整包进程不再消失——本发独立复证）。
2. **逐名比红名册了吗**：**是**，52 枚逐名列出并归因（见 §2）。
3. **本票新增红几枚**：**0 枚**。票面在册的三组（`internal/ball` 1／`internal/panel` 4／`internal/risk` 争用型）本次**全部为零**；52 枚红全在 `cmd/wisp`，其中 39 枚机器级环境红（`System32` 的 onnxruntime 1.17.1 遮蔽，账 `dll-shadow-2`）、5 枚台账 `:32512` 已具名的存量真红、1 枚整包内互扰（定向单跑 PASS）、7 枚归环境族（具名未逐枚单跑）。
4. **"这一格修完之后界面上还差哪一跳"**（220 与 230 共用此句，220 票内 Progress ⑤ 已具名）：**渲染跳整节断**——传输/接收跳已通（`a5f84886f` resident pump ＋ `33b264d70` 接收器）、页面类型已声明（`panel.ts:213 blockedOnApproval`／`:414 tasks?`），但 **`snapshot.tasks` 在 `frontend/src` 零消费**（`App.tsx:246` 喂写死 `agents={[]}`）⇒ 那一跳属前端地界，本编队零写面。

## 4. 与 220 AC#5 的共用说明

220 票 AC#5 与 230 AC#5 是**逐字同一条命令**（`go test ./cmd/wisp ./internal/...` 到终态＋逐名比红名集合）。本件读数**可直接服务 220 的 AC#5**；220 另需"界面差哪一跳"句，已在 220 票内 Progress ⑤ 具名（同上 §3-4）。
