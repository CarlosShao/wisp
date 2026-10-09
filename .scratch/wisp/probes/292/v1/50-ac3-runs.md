# 292-v1 50 — AC#3 整包两发（辅证：名册作差）

环境尺（每发起手前必跑，逐字 stdout）：

- 第 1 发前：`tasklist //FI "IMAGENAME eq wisp.exe"` ＝ `INFO: No tasks are running which match the specified criteria.`；`balldebug.exe` 同。
- 第 1 发后、第 2 发前：两把同逐字（同上那句）。
- 第 2 发后：两把同逐字 ⇒ **零残留进程**，因此不需要 `Get-CimInstance Win32_Process` 判父链，也没有别人的进程要具名报回。

尺（逐字，两发同一把）＝
`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -v > <raw>.md 2>&1`
三数尺＝`grep -cE '^[[:space:]]*--- PASS'`／`'--- FAIL'`／`'--- SKIP'`（原始日志已入库：`30-run1-raw.md` 311,023 字节、`31-run2-raw.md`）。

| 发 | rc | PASS | FAIL | SKIP | 顶层红名册 |
|---|---|---|---|---|---|
| v1-run1 | 1 | **382** | **5** | **2** | `TestPanelHostRealWindowHopAndLifecycle`／`TestAC4FocusReturnToPriorWindowGap33r5`／`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`／`TestAC14AwaitedBindingReplyReachesThePage`／`TestAC14GoSideEvalPushReachesThePage` |
| v1-run2 | 1 | **383** | **4** | **2** | 同上去掉 `TestAC4FocusReturnToPriorWindowGap33r5` |

- 两发包尾逐字：第 1 发 `FAIL	github.com/CarlosShao/wisp/cmd/wisp	477.894s`、第 2 发 `FAIL	github.com/CarlosShao/wisp/cmd/wisp	488.127s`（制表符分隔，`rc=1`）。
- SKIP 两发同两枚，逐字：`--- SKIP: TestPanelHostLatencyPercentilesAC2 (0.00s)`、`--- SKIP: TestAC247LiveMicrophoneLevelsReachTheBallSeam (0.00s)`。
- **±1 漂移实测**：`TestAC4FocusReturnToPriorWindowGap33r5` 第 1 发 `--- FAIL ... (5.10s)`、第 2 发 `--- PASS ... (5.24s)` ⇒ 票面 AC#3 那条"同树两发自己会漂"的警告**我这把尺亲眼复现**，且漂的是**基线里的既有红**，不是新增红。
- 两发红名册**并集**＝那 5 枚，**交集**＝4 枚 ⊂ 基线 5 枚 ⇒ **新增红 0 枚**（判据口径＝并集，从严）。

## 对拉基线（编排者 18:1x 两发，无探针）

尺＝`grep -hE '^\s*--- FAIL' .scratch/wisp/probes/orch/logs/c1-post-{1,2}.txt | sort | uniq -c` ⇒ 5 枚各出现 **2 次**（逐字时长两发不同：PanelHost `5.06s/5.29s`、AC13 `20.01s/20.07s`、AC14Awaited `20.01s/20.02s`、AC14GoSide `20.01s`×2、AC4Focus `5.11s/5.15s`）⇒ 编排者两发都是**同一套 5 枚**。
`grep -cE '^\s*--- PASS' c1-post-1.txt` ＝ **0** ⇒ 坐实"编排者那两发不带 `-v`、PASS 恒 0 是尺的口径不是全红"（裁定节 C1 的说法成立）。

## 腿那两枚外因红的复核

尺＝`grep -E '^[[:space:]]*--- FAIL: (TestAC246DevLegIgnoresTheTestTaskInjection|TestTicket223ModeLooseningChangesTheRunningModeAfterAllow)' 30-run1-raw.md 31-run2-raw.md` ⇒ **零命中（rc=1）**；
反向尺＝同两名配 `--- PASS` ⇒ 四行逐字（每发各两枚，`TestAC246 … (4.11s)`／`(3.72s)`、`TestTicket223 … (2.27s)`／`(2.36s)`）。
⇒ **无探针态下这两枚两发皆绿**，与裁定节 C1 一致；腿那一发把它们读成红，归因"活着的开发 `wisp.exe`（PID 9084）"**与我的读数不矛盾**（我那两发起手/收尾盘上都没有 `wisp.exe`）。

## AC#3 判语：【成立】（新增红 0 枚；主证是 40 件的语义面零，本格为辅证）
