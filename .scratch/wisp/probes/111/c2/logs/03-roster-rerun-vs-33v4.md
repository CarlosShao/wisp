# 03 · 窗口依赖名册：本腿自己复跑那把尺（第 108 条）＋与 `33-v4` 的件作差

## 那把尺（逐字取自 `33-v4` 的 `verdict.md:240`）

> 尺＝逐枚扫函数体里的开窗调用（`startPanelForTest|webview2.New|bringUp|.Show(|HotShow(|BringUp(`）

本腿复跑（⛔ 未跑任何 `go` 命令，纯 `awk` 体扫；`awk` 不是 Go、不动编译面）：

```
cd cmd/wisp && awk '
FNR==1 { ftag = (/^\/\/go:build/ ? substr($0,12) : "(none)") }
/^func Test[A-Za-z0-9_]*/ { name=$2; sub(/\(.*/,"",name); body=1; n=0; next }
body && /^}/{ if(n>0) printf "%s\t%s\t%s\twindowCalls=%d\n", name, FILENAME, ftag, n; body=0 }
body { c = gsub(/startPanelForTest|webview2\.New|bringUp|\.Show\(|HotShow\(|BringUp\(/,"&"); n += c }
' *_test.go > /d/tmp/wisp111c2/roster-rerun.txt          # rc=0
wc -l < /d/tmp/wisp111c2/roster-rerun.txt                # 15   rc=0
```

本腿复跑的名册（15 枚，锚＝本腿 HEAD `052b393f`；`tag` 取自每文件第 1 行 `//go:build`）：

| 用例 | 文件 | build tag | 体内命中 |
|---|---|---|---|
| `TestAC13BringUpRefusesAThreadWithAQueuedClose` | `panel_resident_windows_test.go` | `windows` | 14 |
| `TestAC13BringUpSurvivesAReusedThreadQuit` | `panel_resident_windows_test.go` | `windows` | 7 |
| `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` | `panel_resident_windows_test.go` | `windows` | 2 |
| `TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever` | `panel_resident_windows_test.go` | `windows` | 2 |
| `TestAC14AwaitedBindingReplyReachesThePage` | `panel_resident_windows_test.go` | `windows` | 1 |
| `TestAC14GoSideEvalPushReachesThePage` | `panel_resident_windows_test.go` | `windows` | 1 |
| `TestAC4FocusReturnToPriorWindowGap33r5` | `panel_host_windows_test.go` | `windows` | 2 |
| `TestAC4PriorFocusSurvivesARefusedPanelSample` | `panel_resident_windows_test.go` | `windows` | 2 |
| `TestBallPanelGesturesReachThePanelThread` | `panel_resident_windows_test.go` | `windows` | 1 |
| `TestPanelHostRealWindowHopAndLifecycle` | `panel_host_windows_test.go` | `windows` | 3 |
| `TestPanelThreadIsSTAAndExitsCleanly` | `panel_resident_windows_test.go` | `windows` | 1 |
| `TestTicket255PanelHostBuildsItsWindowOptions` | `panel_geometry_255_test.go` | `windows` | 7 |
| `TestLive35v2PagePostMessageEnvelopeReachesTheGoDoor` | `panel_transport_live_35v2_windows_test.go` | `windows && winlive` | 3 |
| `TestTicket255RealWindowWidthFollowsTheConfig` | `panel_geometry_255_winlive_test.go` | `windows && winlive` | 5 |
| **`TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive`** | `panel_host_windows_live_test.go` | `windows && winlive` | 2 |

## 与 `33-v4` 那件的作差（`logs/25-window-dependent-roster.txt`，14 枚）

尺：`cut -f1 roster-rerun.txt | sort` 对 `awk '{print $1}' 25-*.txt | grep -E '^Test' | sort`，`comm -23` / `comm -13`，两把都 rc=0。

- **两边同名**：14 枚（`33-v4` 那 14 枚本腿全部复现）。
- **只在本腿这一边**：**1 枚**＝`TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive`（`winlive` 档）。
- **只在 `33-v4` 那一边**：**0 枚**。

差异归因（具名时刻，⛔ 不压成一格）：`33-v4` 取数在 10-08 10:47–10:51、锚 `165b45c4`→`dcbb93d2`；本腿锚 `052b393f`。
`panel_host_windows_live_test.go` 第 1 行是 `//go:build windows && winlive`（本腿现量），它**从来没进过 CI 常规 `cmd/wisp` 那一步**（那一步不带 `-tags winlive`，见 `06-*.md`），也**从来没进过任何本机常规整包读数**（`-- NOT_IN_FILE`，见 `04-*.md`）。
⇒ 这一枚的**存在**不改变"CI 6 枚红"的归因，但**枚数**要按 `052b393f` 报 **15**、不报 14。

⚠ 另：`windowCalls=` 那一栏本腿与 `33-v4` 不同形（如 `TestAC13BringUpRefusesAThreadWithAQueuedClose` 本腿 14、它 1）。**枚数与命中数是两把尺**，本件只用**名字集**作差，命中数只作本腿自证"体扫真跑到了"。

## 名册 × CI 色 × 本机色 的三方对照（核心表，另见 `04-*.md`）

| 用例（窗口依赖名册） | CI run 305 | 本机 10-08 08:18 那发整包 |
|---|---|---|
| `TestAC13BringUpRefusesAThreadWithAQueuedClose` | **PASS** (0.87s) | **PASS** (5.02s) |
| `TestAC13BringUpSurvivesAReusedThreadQuit` | **PASS** (0.70s) | **PASS** (5.03s) |
| `TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever` | **PASS** (0.06s) | **PASS** (0.05s) |
| `TestAC4PriorFocusSurvivesARefusedPanelSample` | **PASS** (0.28s) | **PASS** (5.09s) |
| `TestBallPanelGesturesReachThePanelThread` | **PASS** (0.30s) | **PASS** (5.02s) |
| `TestPanelThreadIsSTAAndExitsCleanly` | **PASS** (0.23s) | **PASS** (5.01s) |
| `TestTicket255PanelHostBuildsItsWindowOptions` | **PASS** (0.00s) | **PASS** (0.00s) |
| `TestAC14AwaitedBindingReplyReachesThePage` | **PASS** (1.21s) | **FAIL** (20.02s) |
| `TestAC4FocusReturnToPriorWindowGap33r5` | **PASS** (2.72s) | **FAIL** (5.15s) |
| `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` | **SKIP** | **FAIL** (20.02s) |
| `TestPanelHostRealWindowHopAndLifecycle` | **FAIL**（`-1` 形） | **FAIL**（`-1` 形） |
| `TestAC14GoSideEvalPushReachesThePage` | **FAIL**（`title=""`） | **FAIL**（`no report ... within 15s`） |
| `TestPanelHostLatencyPercentilesAC2`（不在名册，随 lifecycle 漂） | **SKIP** | **SKIP** |
| `TestLive35v2PagePostMessageEnvelopeReachesTheGoDoor` | — 不编译 | — 不编译 |
| `TestTicket255RealWindowWidthFollowsTheConfig` | — 不编译 | — 不编译 |
| `TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive` | — 不编译 | — 不编译 |

⇒ **名册 15 枚里，12 枚在 CI 那台托管机上是能跑的、8 枚在 CI 是绿的**；红只有 2 枚；另有 2 枚（`AC14Awaited`／`AC4Focus`）**在 CI 绿、在本机红**。
⇒ "这台 runner 没有那能力"这个说法，对上表**只能覆盖 3 枚 winlive 档（今天连编译都没进）**，覆盖不到 CI 那 2 枚红。

rc=0
