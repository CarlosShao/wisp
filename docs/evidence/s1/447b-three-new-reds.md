# 票 447 第二轮 — 三枚新 CI 红的定性（2026-11-14 21:0x–21:5x +08，写腿）

**起手 HEAD** `2512c900b`（`git rev-parse HEAD`，现跑）。写面＝`cmd/wisp/**` 的 `*_test.go` ＋
`internal/agent/**` 的 `*_test.go` ＋ 本文件 ＋ `.scratch/wisp/probes/447b/**` ＋ 本票 Progress log。
**产码改动＝0 字节**（`git show --name-only` 见 §7）。

> 口径：三条红句、三枚落点、四数都从 `gh run view --job <id> --log` 的**原文**里取，不是转述。
> 用到的窗口逐字剪进 `.scratch/wisp/probes/447b/logs/*.redwindow.txt`；完整 job 日志（3–4 MB/发）
> 留在工作盘 `.scratch/wisp/probes/447/logs/`，**不进货**（同上一腿口径，见末节）。

## 判决表

| # | 靶子 | CI job | 判 | 一句话 |
|---|---|---|---|---|
| 3 | `Test285LoopDispatchesOneCorrPerCall` | test-core（linux）step7 | **(c) 同形，但落点在尺不在产码** | 尺的 `recordCorrProvider` 在真并发下无锁 append，丢一次写；产码 0 字节；**已修尺** |
| 2 | `Test26R3AnnounceClosesTheCaptureGate` | test-windows step10 | **(a) 读早了** | 等的是 `stage==Warm`，判的是 gate，而 `SetSpeaking(false)` 在 stage 之后；CI 那一发正好落在两道语句之间 |
| 1 | `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` | test-windows step10 | **(d) 墙钟上界（⚠ 带 (c) 风险标记）** | 15 s 预算内 8 次探针**一次回话都没有**，尺把「空回话」当成「一枚 id 都不在」下了判；已把两个读数拆开 |

**有没有 (c)？没有产码 (c。** 唯一 0.00 s 就红的那枚（#3）查到底：产码的并发是对的，错的是尺。

---

## 三枚红的原文（逐字，来自 job 日志）

### #3 `Test285LoopDispatchesOneCorrPerCall` — `43593532737` / job `132443116610` / step `Portable package tests`

```
2026-11-14T12:13:01.3722448Z === RUN   Test285LoopDispatchesOneCorrPerCall
2026-11-14T12:13:01.3724014Z     ticket285_corr_distinct_rulers_test.go:110: dispatched requests = 1, want 2 (one per tool call):
    [{TaskID:3aa5e7eb-0740-4cb9-be97-94fa6c08ae5d CorrelationID:3aa5e7eb-0740-4cb9-be97-94fa6c08ae5d#call_p2
      CallID:call_p2 Name:echo Args:{"text":"second"} Timeout:2s}]
2026-11-14T12:13:01.3725272Z --- FAIL: Test285LoopDispatchesOneCorrPerCall (0.00s)
2026-11-14T12:13:01.3759504Z FAIL	github.com/CarlosShao/wisp/internal/agent	2.132s
```

四数 `RUN=3083 PASS=2090 FAIL=1 SKIP=0`，**这个包只红这一枚**。
读数形状本身就在指路：**少了一条，活下来的是 `call_p2`**（后写的那条），`CorrelationID` 本身是好的
（`<task>#call_p2`，per-call corr 对）。少一条、且活的是"最后写的那条"＝典型的**丢写**，不是"产码少派了一次"。

### #2 `Test26R3AnnounceClosesTheCaptureGate` — `43591105754` / job `132435815655` / step10

```
11:54:56.406Z voice: capture gate closed for the announcement reason=speaking path=T
11:54:56.406Z voice: D43 walked from=Thinking event=llm.first-token to=Speaking
11:54:56.826Z voice: D43 walked from=Speaking  event=tts.speak-done to=Warm
11:54:56.826Z voice[scope]: session: SCOPE-WARM reason=announce-end engines_hot=1
11:54:56.826Z voice: capture gate reopened (speech ended) reason=announce-done path=T
             resident_voice_selfecho_26r3_windows_test.go:159: MUTANT 26r3-M2 SHAPE: the capture gate was left
             in Speaking after the sentence finished (closed=1 reopened=0 speaking=true) ...
11:54:56.876Z wisp: voice leg stopped: ... gate_closed=1 gate_reopened=1 ...
--- FAIL: Test26R3AnnounceClosesTheCaptureGate (9.87s)
```

**同一毫秒**里：`to=Warm` 先写、`capture gate reopened` 后写；尺读到的是中间那一刻。
而且**同一发**的 stop 行是 `gate_closed=1 gate_reopened=1` —— 门最后是开的，只是被判早了。

### #1 `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` — `43591208796` / job `132436112929` / step10

```
11:50:31.902Z === RUN   TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe
11:50:31.902Z     panel_resident_windows_test.go:462: AC#13 probes from the resolved entry (1068 bytes): 1 id(s) [root]
11:50:32.421Z wisp: panel window is up (test-harness, cold 511.5 ms, hot path 0.0 ms)
11:50:47.494Z     panel_resident_windows_test.go:470: AC#13 no probe attempt over 8 attempt(s) saw a present id; last report: ""
11:50:47.494Z     panel_resident_windows_test.go:472: AC#13 page answer (head ff0edb4b): "" - 0 of 1 probe id(s) present in the live document
11:50:47.494Z     panel_resident_windows_test.go:475: after a real cold start the live document contains NONE of the
                 1 element ids the embedded entry declares (the page itself answered "")
--- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (15.64s)
```

四数 `RUN=847 PASS=605 FAIL=1 SKIP=18`，**同一发里 `TestPanelPageFetchesItsOwnAssets299` 是绿的**：
`299 AC#2 page report ... href="https://wisp.panel/index.html" ready="complete" sheets=1 rules=110 root="yes"
roottext=116 resjs=1 rescss=1` —— 同一个进程、同一个 entry URL、同一条资源过滤器、同一扇 door，
**另一个窗口能把整页读出来**。

---

## #3 `Test285LoopDispatchesOneCorrPerCall` — 判「真缺陷，但缺陷在尺」

### 凭什么 0.00 s 就红

1. **产码确实真并发**：`internal/agent/loop.go:791` 为每一枚调用 `observe.Spawn` 一个 goroutine，
   `:882/:884` 在里面调 `dispatch` → `(*recordCorrProvider).Execute`（D38d 并发上限 4）。
   ⇒ 两条 `Execute` 同时跑，**不是**顺序派发。
2. **尺的写入没有锁**：`ticket285_corr_distinct_rulers_test.go:37` 原文
   `p.reqs = append(p.reqs, req)`。两条 goroutine 同时读 `len=0`、同时写 `index 0`，
   后写的覆盖先写的 ⇒ `len==1`，留下的正是后写的那条（CI 里是 `call_p2`）。
3. **可判事实（不是推理）**：同一份 HEAD、同一枚用例，`-race` 直接点名：

```
$ go test -race -count=1 -run '^Test285LoopDispatchesOneCorrPerCall$' ./internal/agent
WARNING: DATA RACE
Previous write at 0x00c00000a390 by goroutine 18:
  github.com/CarlosShao/wisp/internal/agent.(*recordCorrProvider).Execute()
      .../ticket285_corr_distinct_rulers_test.go:37
  github.com/CarlosShao/wisp/internal/agent.(*Loop).dispatch()  loop.go:884 / :882
  github.com/CarlosShao/wisp/internal/agent.(*Loop).executeCalls.func1()  loop.go:794
Previous write at ... (第二条同样的告警)
--- FAIL: Test285LoopDispatchesOneCorrPerCall (0.02s)
```
   日志：`probes/447b/logs/02-pristine-race-red.txt`（在仓外导出树上跑 HEAD 的原始尺）。
4. 所以 0.00 s：**不需要等任何落定**，两条 goroutine 一撞就丢写，撞上就红。

### 改了什么（只动尺，断言一字节没动）

`ticket285_corr_distinct_rulers_test.go`
- `recordCorrProvider` 加 `mu sync.Mutex`；`Execute` 的 append 加锁；
- 新增 `requests() []ToolRequest`（同锁下取快照），用例读的是它，**不再是裸 `p.reqs`**；
- 断言原文（`len(reqs)!=2`、`task id`、`collapse onto task`、`share one routing key`、前缀、后缀自路由）
  **一行没删没放宽**，只把三处 `p.reqs` 换成 `reqs`。

### 复跑读数（本机，低位）

| 命令 | 读数 |
|---|---|
| `go test -count=5 -run '^Test285(LoopDispatchesOneCorrPerCall|CallCorrIsDistinctPerCall)$' ./internal/agent` | 10/10 PASS（`ok … 0.140s`） |
| `go test -race -count=3 -run '^Test285LoopDispatchesOneCorrPerCall$' ./internal/agent` | `ok … 1.159s`，**零 DATA RACE** |

### 牙（仓外导出树，`%TEMP%\wisp447b`，从 `git archive HEAD` 出）

| 发 | 突变 | 读数 |
|---|---|---|
| M-285 | `loop.go` 的 `return taskID + "#" + callID` → `return taskID + "#call"`（call id 不再进路由键） | **红且点名**：`ticket285_corr_distinct_rulers_test.go:150: both calls of one task dispatched with corr "<task>#call": the two asks share one routing key`，`--- FAIL (0.01s)`；SHA256 还原后与仓同值、再绿（`ok … 0.059s`） |

日志：`03-m285-collapse-red.txt`、`04-m285-restored-green.txt`。

---

## #2 `Test26R3AnnounceClosesTheCaptureGate` — 判 (a) 读早了

### 判据（两条，都可核）

1. **时间顺序**（上面原文）：`to=Warm` 与 `capture gate reopened` 在**同一毫秒**，且 reopen 在后。
2. **代码顺序**：`resident_voice_windows.go` 的 `finish()` 先走 D43 `Speaking→Warm`（`tts.speak-done`），
   然后才 `SetSpeaking(false)`；而尺等的是 `rv.residency.stage()==session.StageWarm`（旧 `:151`），
   拿到的却是 gate 的读数（旧 `:154-160`）。**等待谓词 ≠ 断言谓词**，中间那段就是红窗。
3. **同一发的反证**：`gate_closed=1 gate_reopened=1`（stop 行）＋ `speakGateOff==1`（尺自己 `:161` 的等待随后就成立）
   ⇒ 门确实开了，只是判早了。

⇒ 与票 447 上一轮 `WireR1` 的 `waitStage447` 同形（(a)）。

### 改了什么

`resident_voice_selfecho_26r3_windows_test.go`
- 新增 `waitGate26r3(t, g, within)`：10 ms 轮询、单调时钟、上界＝本包既有的 `panelThreadWait`，
  谓词＝**gate 自己的落定态**（`closed>=1 && reopened>=1 && !speaking`）；
- 旧 `:154-160` 的「M1/M2 命名」整段搬到**超时路径**：`closed==0` ⇒ `MUTANT 26r3-M1 SHAPE`，
  否则 ⇒ `MUTANT 26r3-M2 SHAPE`（**两句原文保留**，尾上加了可 grep token `447-26r3-gate-not-settled`
  与本发耗时）；
- stage 的等待（`voice15eWaitFor … StageWarm`）保留在原位；
- `:165-180` 的精确计数（`closed==1`、`reopened==1`、事件序列、leg 自己计数器）**一行未动**。

### 复跑读数（载具，一次一发，`-test.count=3`）

```
447…  resident_voice_selfecho_26r3_windows_test.go:233: GATE-26R3 closed=1 reopened=1 events=2
      mid_speech_transcripts=0 spoken_samples=89490 gate_on=1 gate_off=1      ×3
--- PASS: Test26R3AnnounceClosesTheCaptureGate (24.91s / 25.08s / 20.93s)
```
日志：`05-carrier-both-count3.txt`。

**牙欠账（照记）**：26r3-M2 的突变（去掉 `SetSpeaking(false)`）要在 `cmd/wisp` 的产码上改，
而 voice 夹具（`third_party/`，git-ignored）不在仓外导出树里 ⇒ 本轮机位的牙**没法在导出树上跑**。
本轮的牙只给了 #3 与 #1 各一发（见上/下）。⇒ 若要 26r3 的牙，需要一枚**允许临时改产码**的载具
或一台带夹具的导出树；我不在本仓工作树里改产码。

---

## #1 `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` — 判 (d)，⚠ 带 (c) 风险标记

### 判 (d) 的依据（含量到的反证）

1. **15.64 s ＝ 15 s 预算用光**：`panelThreadWait=15s`（`panel_resident_windows_test.go:53`），
   `probeAttemptSettle=2s` ⇒ 8 次尝试 × 2 s，最后一次被 deadline 截短；`attempts=8` 与读数自洽。
2. **不是「读早了一次」**：尺有界重问 8 次、跨 15 s，**一次回话都没有**（`last report: ""`），
   而同一进程的 299 用例在同一发里读出了整页（`ready=complete`）。⇒ (a) 的"等一等就好"形状不成立
   （重问已经在等）；(b) 无硬件/夹具前提缺失（本机同码 1.66–1.96 s 绿）。
3. **不是「页面在但内容不对」**：AC#13 的**已知缺陷形状**（把 `serveEntry` 排在 `firstRoundTrip` 之前，
   即 33 号票修掉的那个顺序）在同一预算内的读数是 **`718 attempt(s) … last report: "0"`**
   （`probes/299/r1b` 的 m1 导出树读数；本轮我在导出树复现，见下 §牙）。
   ⇒ 「空回话」与「0」是**两个不同的世界**：前者页面根本没说话，后者页面说话了说没有。
4. **不是"机器整体慢 4×"**：同发冷启动 `cold 511.5 ms`，绿发 `445.8 ms`（`43591051106` 同用例 2.36 s、attempt 2 拿到 `"1"`）
   ＝ **1.15×**，不是 30×。⇒ 一条 1068 字节、走进程内资源过滤器、无网络的导航，在 15 s 里没落定，
   按"负载折算"解释不了。
5. 因此我把它判成 **(d) 墙钟上界**：这一发的判据**只**建在「15 s 内必须落定」这个墙钟上，而在那 15 s 里
   **尺连一个读数都没拿到**。

### ⚠ (c) 风险（必须写清，不许含糊）

同一份产码、同一套测试，6 发里 3 红 3 绿（六个 head `ff0edb4b0` / `1984618a4` / `6aa244b67` /
`445ea5b98` / `468cd2894` / `35faf8064` **全是 dispatch/台账/票面的文档提交**
—— `git show --stat` 逐枚核过，**产码零改动**）。⇒ 这只证明"偶发"，**不证明产码无罪**。
「导航 15 s 没落定」也可能是产码的真缺陷（`serveEntry` 只发 `Navigate(panelEntryURL)`、
没有任何 `NavigationCompleted` 读面、没有重试、没有失败读数：`panel_host_windows.go:789` 一带）。
**我这一轮没有把它判成 (c)**，因为手上没有一个"产码错了"的正面事实；反过来也没有一个"产码没错"的事实。
⇒ 下一发如果**在新的尺下**再红成 `447-ac13-handover-never-answered`，那次读数就是判 (c) 的起点
（新尺把「页面没说话」单独立成 failed measurement，不再冒充内容判决）。

### 改了什么（把两个读数拆开，断言不动）

`panel_resident_windows_test.go`
- 新增 `ac13HandoverProbeID = "447-ac13-handover"` 与 `awaitHandoverEntryDocument(t, rp, h)`：
  **先问"你是哪个文档、完了没有"**（`location.href|document.readyState|document.title`），
  谓词＝`readyState=="complete"`（**只有这一条**，故意不要求 href 是 entry），有界 `panelThreadWait`；
  超时 ⇒ `Fatalf` 打 **`447-ac13-handover-never-answered`**，带 `attempts`、`reportCount`、`last`、
  `describeRequests(h)`（"门到底收到了什么"）；
- `TestAC13ColdStart…` 在 `awaitEntryReport` 之前插入这一次 settle；
- **AC#13 的内容断言（`:535` 那段 `present == 0` ⇒ 红，读出 `"none of the ids"`）一字未改**。

**为什么谓词只要 `complete` 而不是"href 是 entry"**：若要求 href 是 entry，那"顺序反了"的世界
（页面停在探针页，`about:blank|complete|`）会落到我这条"没落定"的红上，**AC#13 的原句子就再也不点名了**
—— 那才是钝化。现在的分工是：
- 顺序反 / 内容不在 ⇒ settle 立刻过（探针页也是 `complete`）⇒ **原句红**（牙在，见 §牙）；
- 一个文档都不回话 ⇒ 新的 failed-measurement 红（带数），**不再冒充**内容判决。

### 复跑读数

| 场景 | 读数 |
|---|---|
| 载具，`-test.count=3`（与 26R3 同发） | `447-ac13-handover-settled on probe attempt 5 / 7 / 8: "https://wisp.panel/index.html|complete|Wisp"`，随后 `entry-id report arrived on probe attempt 1: "1"`；**PASS 1.46 / 1.75 / 1.74 s** |
| 仓外导出树 baseline（同一份改过的尺） | `settled on probe attempt 7` ⇒ `page answer "1"`；**PASS 1.75 s** |
| 导出树 + 顺序突变（M-AC13） | `settled on probe attempt 1: "about:blank|complete|"` ⇒ `page answer ""`→`"0"` ⇒ **原 AC#13 句子红**，`--- FAIL (16.44 s)` |
| 还原后（SHA256 与仓同值） | `settled on probe attempt 6` ⇒ `page answer "1"` ⇒ **PASS 1.96 s** |

▲ 注意本地 settle 需要 5–8 次尝试才 `complete`：这 5–8 次**都不是 2 s**，
每次都是拿到一个新回话（`readyState=loading`）就立刻重问 —— 也就是说**导航进行中的文档是会回话的**。
⇒ CI 那一发"一次都不回话"比"慢"更像"那时窗口里没有能跑脚本的文档"。这条读数已写进上面 (c) 风险段。

---

## 窄化/放宽台账（AC#4）

| 项 | 变化 | 判 |
|---|---|---|
| AC#13 内容断言 | **未改**（句、阈值、`present==0` 判据全同） | 无放宽 |
| AC#13 新增 settle 红 | 新增一条**失败的测量**出口（不是把红改绿） | 只增读数 |
| 26R3 M1/M2 命名句 | 两句原文保留，搬到超时路径 | 无放宽（红仍红，且仍点名） |
| 26R3 精确计数 / 事件序列 / leg 计数器 | 未改 | 无放宽 |
| 285Loop 断言 | 未改（只换读取通道 + 给记录器加锁） | 无放宽 |
| 阈值 / golden / `thresholds.go` / 契约 | **零改动** | — |

## 没做 / 没翻的

- **AC#5（CI 连续 5 发 `FAIL` 恒 0）**：要推送后由编排者用 `collect-ci-fails.ps1` 复读，本腿不 push。
- **#1 的 (c) 最终定性**：需要下一发 CI 在新尺下再红一次（`447-ac13-handover-never-answered`）才有正面事实。
- **26R3 的牙**：导出树无 voice 夹具，见上。

## 整包四数（载具，本机，2026-11-14 21:24–21:48，一次一发）

```
个  RUN=853  PASS=611  FAIL=2  SKIP=17      （ok/FAIL 行：FAIL github.com/CarlosShao/wisp/cmd/wisp 1817.00s）
```

两枚红**都不是本腿的**，且都是票面点过名的本机已知红：
- `TestAC4FocusReturnToPriorWindowGap33r5`（5.42s）— 前台被占；
- `TestAC246DevLegIgnoresTheTestTaskInjection`（37.42s）— 「a dev Wisp is already running in this session」（编排者起的驻留球）。

★ 三枚靶子里两枚在同一次整包里**绿**：`447-ac13-handover-settled on probe attempt 8: "https://wisp.panel/index.html|complete|Wisp"`、
`GATE-26R3 closed=1 reopened=1 events=2 … gate_on=1 gate_off=1`。日志：`11-carrier-whole-package.txt`。

## 本 commit 动了什么（`git show --name-only`）

```
cmd/wisp/panel_resident_windows_test.go                      （settle 闸＋新 helper）
cmd/wisp/resident_voice_selfecho_26r3_windows_test.go        （gate 落定谓词＋M1/M2 搬到超时路径）
internal/agent/ticket285_corr_distinct_rulers_test.go        （记录器加锁＋requests() 快照）
docs/evidence/s1/447b-three-new-reds.md                      （本文件）
.scratch/wisp/issues/447-…md                                 （Progress log，只追加）
.scratch/wisp/probes/447b/logs/**                            （台件）
```
**产码 0 字节**；`thresholds.go` 0 枚；`internal/panel/**` / `internal/tools/**` / `internal/llm/**` /
`tools/d22scan/**` / `frontend/**` / `design/**` 0 枚。

## 复跑命令（逐字）

```powershell
# 载具（cmd/wisp 只有这一条路；⛔ 同时跑两发；⚠ PowerShell 里每个 -test.* 必须单独引号，
# 不然 PS 会把 -test.count=3 拆成 -test + .count=3，二进制报 "flag provided but not defined: -test"）
$env:WISP_CLI_STAGE_DIR="$PWD\.scratch\wisp\stage-447b"
& "D:\work\soft\Git\bin\bash.exe" scripts/wisp-cli-staged-tests.sh `
  '-test.count=3' '-test.run' '^(TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe|Test26R3AnnounceClosesTheCaptureGate)$'

# internal/agent（可直接 go test，一次一发）
go test -count=5 -run '^Test285(LoopDispatchesOneCorrPerCall|CallCorrIsDistinctPerCall)$' ./internal/agent
go test -race -count=3 -run '^Test285LoopDispatchesOneCorrPerCall$' ./internal/agent
```

## 台件清单（`.scratch/wisp/probes/447b/logs/`）

`43591208796.ac13.redwindow.txt`（#1 红窗原文）· `43591105754.26r3.redwindow.txt`（#2 红窗原文）·
`43593532737.285loop.redwindow.txt`（#3 红窗原文）· `43591051106.ac13.greenwindow.txt`（对照绿发窗口）·
`02-pristine-race-red.txt`（-race 点名）· `03-m285-collapse-red.txt` / `04-m285-restored-green.txt`（#3 牙）·
`05-carrier-both-count3.txt`（载具 ×3）· `06-export-build.txt` / `07-m0-export-ac13-baseline.txt` /
`09-mut-ac13-order-swap-red.txt` / `10-mut-restored-green.txt`（#1 牙）· `11-carrier-whole-package.txt`（整包四数）·
`285-test-HEAD-pristine.go.txt` / `loop-go-pristine.txt` / `panel_host-go-pristine.txt`（突变前的原字节，供 SHA256 对账）。

⚠ **三份完整 job 日志（3.7 MB ×2 ＋ 3.1 MB）不进货**，按上一腿的同一口径：它们留在
`.scratch/wisp/probes/447/logs/` 的工作盘上（untracked），本文件 §红窗 已经把用到的行逐字摘出来，
任何一行都可用 `gh run view <id> --job <jid> --log` 重取。
