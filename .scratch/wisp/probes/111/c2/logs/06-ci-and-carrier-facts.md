# 06 · CI 载体与运具（carrier）的盘上事实（本腿全部现跑，只读）

## 1 `runs-on` 名册（尺＝`grep -n 'runs-on' .github/workflows/ci.yml`，rc=0；本腿锚 `052b393f`）

```
66:    runs-on: ubuntu-latest            → lint          (:65)
396:   runs-on: ubuntu-latest            → test-core     (:395)
506:   runs-on: windows-latest           → test-windows  (:505)
791:   runs-on: windows-latest           → slo-smoke     (:790)
849:   runs-on: [self-hosted, wisp-slo]  → slo-full      (:848)
923:   runs-on: ubuntu-latest            → lint-frontend (:922)
```
（作业名与行的对应＝`grep -nE '^  [a-z][a-z0-9-]*:$'`，rc=0。）
同一把尺跑在**已推 tip** `cc31526` 上（`git show cc31526:.github/workflows/ci.yml | grep -nE '^  [a-z][a-z0-9-]*:|^    runs-on'`，rc=0）：
`lint` 66 ubuntu／`test-core` 310 ubuntu／**`test-windows` 420 `windows-latest`**／**`slo-smoke` 654 `windows-latest`**／`slo-full` 712 `[self-hosted, wisp-slo]`／`lint-frontend` 786 ubuntu ⇒ **两棵树同形**。
⇒ 派单那句"只有 `slo-full` 是 self-hosted"**成立**（本腿自量，与 `111-c1` 一致，不是复述）。

## 2 workflow 里有没有"装 WebView2 / 等 WebView2"的步骤

尺＝`grep -in 'webview\|edge\|evergreen\|runtime' .github/workflows/ci.yml`（rc=0）⇒ **只有 2 行、都是散文注释**（`:445`、`:755`，讲"declared out of scope"与"compiled WINDOWS test binary"），**没有任何一步安装或等待运行库**。
`find`/`grep` 全仓无 WebView2 安装件（`111-c1` 与本腿都没量到安装步骤）。
⇒ **CI 那台机器上有没有运行库，不是 workflow 决定的**（见 `07-*.md` 的直接读数）。

## 3 `test-windows` 那条腿里前端 bundle 到底建没建

尺＝`grep -in 'npm\|node' .github/workflows/ci.yml`（rc=0）⇒ `npm ci`/`npm run build`/`npm run render:*` 全在 **`:930-991` 那一段，主人是 `lint-frontend`（ubuntu，`:923`）**；`test-windows`（`:506-644`）里 **0 枚 npm 步**。
CI 侧的正面读数（run 305 同发，`clean-run305.txt`）：`panel_host_gate_test.go:109 ... built=false entry-bytes=0 entry-err=panel: embedded assets are not built (run npm run build in frontend/)`。
⇒ **`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` 的 SKIP 是这条结构造成的**（没有 entry 就没有主体），⛔ 与 WebView2 无关。产码侧对应：`cmd/wisp/panel_host_windows.go:472-488 serveEntry` 在无 bundle 时返错、`:450-452` 改走 `serveNotBuiltNoticeLocked`。

## 4 步级 `if:` 与那道"从未执行"的形状（本腿现量，⚠ 具名尺子）

| 尺（逐字命令） | HEAD `052b393f` |
|---|---|
| `grep -cE '^\s+if:'` | **13** |
| `grep -cE '^[[:space:]]+if:[[:space:]]*\$\{\{[[:space:]]*!cancelled\(\)[[:space:]]*\}\}$'` | **12** |
| `grep -cE '^    if:'`（作业级） | **0**（rc=1＝零命中，⛔ 不是"命令失败"） |
| `grep -cE '^\s+continue-on-error:'` | **0** |
| `grep -c 'timeout-minutes'` | **0** |
| `grep -c 'winlive' .github/workflows/ci.yml` | HEAD **7** ／ 已推 tip `cc31526` **0**（`git show cc31526:... \| grep -c winlive`，rc=1） |

⇒ 台账 `A692`/票 111 `:375` 那串"12 处"是**加步之前**的数，本腿这发是**加步之后**的数（13 = 12 `!cancelled()` ＋ 1 `always()`）；两数分属两时刻＋两把尺，⛔ 压成一格就是错（第 130 条）。
⇒ `winlive compile gate` 那一步**只存在于未推 SHA**：`07`/`02` 两格里它在真实 run 里的颜色是**零读数**（连 `skipped` 都没有），这一条与 `111-c1` 的三枚硬读数一致，本腿复跑了第三枚（`git show ... | grep -c winlive`＝0）。

## 5 运具对 SKIP 的规矩（★这一格决定 ⓑ 那支能不能成立）

`tools/d22scan/runtests.sh` 头部逐字（`:22-26`）：
```
#   2. any `--- SKIP`                             -> SKIP is NOT a pass. A test
#      that legitimately cannot run on this platform must be moved out of the
#      step's scope, not silenced here (that would be lowering a threshold).
```
实现逐字（`:98-102`）：
```
if [ "$skipped" -ne 0 ]; then
    echo "runtests.sh: $skipped test(s) SKIPPED and SKIP is not a pass (ticket 71 AC#3) - $summary" >&2
    echo "runtests.sh: either make the step's scope exclude it or run it where it can execute;" >&2
    echo "             do not relax this script." >&2
    exit 1
fi
```
（逐字对照尺＝`sed -n '96,103p' tools/d22scan/runtests.sh | cat -A`，rc=0。）
⇒ **在这一步里给用例加 `t.Skip` 兜底，换不来绿**：`FAIL=0` 之后仍会因 `SKIP!=0` 走 `exit 1`。今天 run 305 那 2 枚 SKIP 就是这句话射程内的两枚活样本。
⇒ 本仓既有的**合法出口**是"把 scope 摘出去"，不是"在测试里静默"：
- `cmd/wisp` 那一步现在就带着一份**具名 `-skip` 名单**（CI 日志逐字）：`-skip ^(TestDefaultDeadlineWallClockMeasurement|TestSubprocessCrashWriter|TestHelperProcess|TestLiveWasapiSmoke|TestRealDownloadVadThroughPipeline|TestRealDownloadPuncArchiveThroughPipeline|TestSyncRegistryProbeLive)$`；
- 它背后的台账在 `scripts/portable-tests.sh`：`grep -oE '"Test[A-Za-z0-9_]+' scripts/portable-tests.sh | tr -d '"' | sort -u | wc -l` 现量＝**11 枚**（尺 rc=0；行头形状尺 `grep -cE '^\s+"Test'` 同为 11），逐枚：`TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`／`TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`／`TestD34WriteMatrix`／`TestDefaultDeadlineWallClockMeasurement`／`TestHelperProcess`／`TestLiveWasapiSmoke`／`TestRealDownloadPuncArchiveThroughPipeline`／`TestRealDownloadVadThroughPipeline`／`TestSubprocessCrashWriter`／`TestSyncRegistryProbeLive`／`TestWorkspaceSwitchRefusesAJunctionToOutside`。每行带 scope／平台／理由／`next=` 归属，且脚本自己写着 `An unaccounted SKIP also still goes red, so a new skip cannot be parked here by accident`（`:32-33`，逐字尺＝`grep -n 'skip' scripts/portable-tests.sh`，rc=0）；
- `scripts/wisp-cli-tests.sh:60-67` 另一形：**非 windows 就直接 `exit 2`**（"Not a skip and not a pass"），⛔ 不静默。
⇒ **本腿现量：那 6 枚红名没有一个进过这份台账**（`grep -c 'TestPanelHostRealWindowHopAndLifecycle\|TestAC14GoSideEvalPushReachesThePage\|TestTicket223PermissionDenied\|TestRunPacketCarries\|TestAC1AlwaysBranch' scripts/portable-tests.sh` ＝ **0**，rc=1＝零命中）。

## 6 "t.Skip 不是选项"这句在本仓的原文位置

`internal/ball/live_guard_windows_test.go:170`（本腿现量，尺＝`grep -n 'Skip' internal/ball/live_guard_windows_test.go`，rc=0）逐字：
`"Kill pid %d (a leftover balldebug/wisp) and re-run. t.Skip is not an option here.",`
同文件 `:13`：`(registry A5). Skipping past that is what hid it, so nothing here skips`。
⇒ 票 35 `:272` 引的那句"同族判据逐字写着 `t.Skip` is not an option here"**成立**，本腿复跑到了原文与行号。

## 7 `cmd/wisp` 那一步自己声称的历史读数（⚠ 与 CI 实况不同形，登记不裁）

`.github/workflows/ci.yml:574-582` 注释逐字：`on windows its 33 top-level cases go PASS=33 FAIL=0 SKIP=0 rc=0 once the pinned third_party/sherpa-onnx DLLs are on PATH`。
CI 实况（`02-*.md`）：同一步 8 发里 `=== RUN` 是 **261→335**、`FAIL` 是 **6→12**、`SKIP` 是 1→2，**12/12 发 failure**。
⇒ 那句注释说的是**当初接入时的那一发**（33 枚规模），今天这一步 335 枚；本腿**不判**这句话现在还算不算真话，只把两边数字并排放着（⛔ 未动 `.github` 一字）。

rc=0
