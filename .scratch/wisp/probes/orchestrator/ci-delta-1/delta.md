# ci-delta-1：8ae4c23e vs 0589fd9c 的 CI 红名册差（只读取证）

## 0. 起手锚
- `date -Iseconds` = `2026-10-02T08:25:56+08:00`（收尾复量 `2026-10-02T08:5x+08:00`，全程未跑任何吃 CPU 的东西）
- `git log -1` = `a35f7f5e`（`dev`，工作树脏：`.gitignore`/`design/**`/`.scratch/**` 为他人未提交件，本腿未动）
- 新发 run `36889094435`（push 事件，workflow `ci`，headSha `8ae4c23e2ab4…`，16:03:34Z，failure）
- 基线 run `36789659174`（push 事件，workflow `ci`，headSha `0589fd9cc3c4…`，09-30 23:10:15Z，failure）
- 推送区间 `git rev-list --count 0589fd9c..8ae4c23e` = 231
- ⛔ 未把 `36940536372` 计入本差集。但票面前提是错的：**它不是"另一支 workflow"**——`gh run view 36940536372`
  回 `"workflowName":"ci","event":"schedule","headSha":"8ae4c23e…","status":"completed","conclusion":"failure"`。
  同一支 `ci`，同 headSha，已跑完且也红。它是"同码二次采样"的现成负载/抖动对照，本腿按票面没有并进来。
- 取数方法（避开三条已知漏计）：`gh run view <run> --log` 全文下载（各 1.03MB / 1.78MB，5 个作业有日志），
  自建 `slice3.awk` 按 `##[group]Run <cmd>` 把每一行归到它自己的步骤，再 grep `--- FAIL:`／`--- SKIP:`／`##[error]`／
  `^FAIL <pkg>`；**没有用 `--log-failed`**（它只给失败作业、不给失败步骤里的 PASS/SKIP，会少算）。
  未跑步骤单列在 §2，不当绿。日志原文与切片留在本目录（`logs_*.txt`、`raw_*.txt`、`d_*.tsv`、`fails_*.tsv`）。
  本机读数与 CI 读数全程分开：本差集的分母**只有 CI 自己的 `-v` 计数**。

## 1. 两发名册与作差

### 1.1 作业级（6 枚，两发同一批作业名）
| 作业 | 基线 0589fd9c | 新发 8ae4c23e | 差 |
|---|---|---|---|
| lint-frontend | success | success | — |
| slo-smoke | success | success | — |
| test-windows | failure | failure | — |
| lint | failure | failure | — |
| test-core | failure | failure | 作业名不变，**但红的原因整枚换掉了**（§1.3） |
| slo-full | **success** | **failure** | **新增 1 枚作业级红**，且原因不是测试（§1.4） |

### 1.2 步骤级红名册（逐枚，两发各一份）
基线红步骤 5 枚：`lint#7 gofmt(gofumpt)`、`lint#11 staticcheck`、`test-windows#7 cmd/wisp CLI tests`、
`test-windows#8 Portable windows tests`、`test-core#7 Portable package tests(core)`。
新发红步骤 5+2 枚：同前 5 枚，外加 `slo-full#5 SLO full gate (six states + settle + leak)`（结论 `null`）
与 `slo-full#6 Upload SLO report`（`null`）。
静态件的红句两发逐字相同：`##[error].scratch/wisp/probes/185/c1/mut/fs_broken.go:4:1: imports must appear before other declarations`
（`git ls-files --error-unmatch` 确认这枚被跟踪＝票 185 的坏变体进了 gofumpt 的分母，两发同形，非新增）。

### 1.3 测试级名册作差（本次唯一有意义的尺）
runtests.sh 自己印的四个数（逐字，`packages=` 已截短）：

| 步骤 | 基线 | 新发 |
|---|---|---|
| `./cmd/wisp/`（test-windows#7） | `PASS=110 FAIL=7 SKIP=0, RUN=191` | `PASS=147 FAIL=11 SKIP=2, RUN=241` |
| `./internal/proc/ … ./cmd/llmrecord/`（test-windows#8） | `PASS=311 FAIL=12 SKIP=1, RUN=454` | `PASS=320 FAIL=13 SKIP=1, RUN=468` |
| core scope（test-core#7） | `PASS=921 FAIL=4 SKIP=0, RUN=1365` | **没跑到 go test**（见 §2） |

包级 `FAIL` 行：基线 3 行（`cmd/wisp 424.719s`、`internal/risk 6.244s`、`internal/panel 0.784s`），
新发 2 行（`cmd/wisp 579.297s`、`internal/risk 6.247s`）。`internal/panel` 那行消失＝分母没了，不是修好。

**新增红 6 枚测试名 + 2 枚新 SKIP**（逐名，全在 test-windows）：
1. `TestTicket223HandEditedFsLooseningCostsAnL2Card`（cmd/wisp）
2. `TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking`（cmd/wisp）
3. `TestAC246ResidentPipelineAsksThroughTheOneGate`（cmd/wisp）
4. `TestPanelHostRealWindowHopAndLifecycle`（cmd/wisp）
5. `TestResolvePerCallBudget`（internal/risk）
6. staticcheck 新增 3 条发现（`lint#11`，同一枚新测试文件）
   新 SKIP：`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`、`TestPanelHostLatencyPercentilesAC2`

**消失的基线红 4 枚**（test-core，`--- FAIL` 名）：
`TestC21DesignTokensFourWayAgree`、`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`、
`TestComposerContractTypesMatchFrontend`、`TestApprovalCardViewJSONKeysMatchFrontendTypes`。
⇒ **判定：4 枚全部是"未测到"，不是"修好了"**（步骤在 `go test` 之前就被 GUARD C 打死）。

净数：新增红名 **6**（测试名口径；含 SKIP 口径则 **8**），消失 **4**（全为未采到），作业级新增 **1**（slo-full）。

### 1.4 slo-full：这一发的红没有日志
`gh run view 36889094435 --log` 逐作业回 `log not found: 110459845223`；
`gh api repos/CarlosShao/wisp/check-runs/110459845223/annotations` 唯一一条：
> `failure .github:1 The self-hosted runner lost communication with the server. Verify the machine is running and has a healthy network connection. Anything in your workflow that terminates the runner process, starves it for CPU/Memory, or blocks its network access can cause this error.`

`output.title/summary/text` 全为 `null`，作业 `run_attempt=1`，跑 12 分 03 秒（16:03:39→16:15:42）后掉线。
基线那发 slo-full 在同一台 `wisp-selfhosted-01` 上排队 2h18m 后 3 分 32 秒跑完并 success（日志里 `wisp doctor: PASS`）。
⇒ **这一枚不是断言红，是采集红**：D32 六态/settle/leak 的任何读数这一发根本没产生。runner 现在 `status=online`。

## 2. "一步红吃掉后续步骤"清单（这些读数永久采不到，不许当绿）
1. `lint#8 gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)` — skipped（**两发同**）
   ⇒ "受跟踪集合作分母"那一版读数**一次都没拿到过**；#7 用的是全量分母，才把 `.scratch/…/fs_broken.go` 算进来。
2. `lint#9 go vet (module)` — skipped（两发同）⇒ 整个 push 区间 `go vet` **无读数**。
3. `lint#10 go vet (tools/d22scan module)` — skipped（两发同）⇒ 同上。
4. `test-core#7 Portable package tests (core scope)` — 新发这枚**跑起来了但 3 秒内死在守卫**，`go test` 一行没执行
   ⇒ 基线那 1365 个 `=== RUN` 的读数在新发**整体缺失**（含 §1.3 那 4 枚消失的红）。红句逐字：
   > `portable-tests.sh: GUARD C - scope mode=core resolved to a DIFFERENT package`
   > `portable-tests.sh:   set than the one pinned next to it. Pinned: 25, resolved: 35.`
   > `portable-tests.sh:   25a26,35` 后面 10 行逐字是 `> go: downloading github.com/dustin/go-humanize v1.0.1`、
   > `… google/uuid v1.6.0`、`… pelletier/go-toml/v2 v2.2.4`、`… remyoudompheng/bigfft …`、`… golang.org/x/crypto v0.57.0`、
   > `… x/sys v0.48.0`、`… modernc.org/libc v1.75.7`、`… mathutil v1.7.1`、`… memory v1.12.1`、`… sqlite v1.59.0`

   凭据＋尺子出处：`scripts/portable-tests.sh:261` 是 `go list "${scope[@]}" >"$resolved" 2>&1`——**stderr 被并进了分母**。
   35 = 25（真包）+ 10（`go: downloading` 噪声）。`git diff --stat 0589fd9c..8ae4c23e -- scripts/portable-tests.sh` 为空，
   `git diff --stat … -- internal/risk/` 也为空 ⇒ 包集合没变，变的是这台机器的模块缓存是冷的。
   这是**仪器自己的口径 bug**，不是产品断言失败，也不是可以放过的绿。
5. `slo-full#5`、`#6`、`#11`、`#12` = `null`（新发）＋整个作业日志缺失（§1.4）。
6. `test-windows#16/#17`、`test-core#15` = Post 步骤 skipped（无读数含义，仅登记）。

## 3. 逐枚归因（新增红 6 + SKIP 2）
| # | 红句（逐字，截断处标 …） | 成因判定 | 凭据 |
|---|---|---|---|
| 1 | `config_reload_223_test.go:303: the console did not render the reload card: wisp run: 答复监听已接入（卡片上给了编号）…` ／ `:306: the card does not name the key it is asking about: …` | **代码**（真回归） | 该测试文件 `ec7a034d` 加入、区间内 `git diff` 为空；基线逐字 `--- PASS: TestTicket223HandEditedFsLooseningCostsAnL2Card (2.13s)`，新发 `--- FAIL: … (2.43s)`。区间内 `internal/config internal/panel cmd/wisp internal/agent` 共 `32 files changed, 10598 insertions(+), 100 deletions(-)` |
| 2 | `ticket224_assembly_test.go:316: covered call booked decision="timeout", want "allow_session_grant"` ／ `:319: … grant_id=<nil>, want the covering row 1: tool_call.grant_id is the column D45-2 promised for forensics` ／ `:324: audit is missing the GRANT-USE line naming row 1` | **代码**（D45-2 装配侧没接：会话授权未被查询，卡片等到超时） | 新文件 `1c601fab`（本次 push 带入，`git merge-base --is-ancestor` 判为 NEW-in-push），基线里不存在此名；耗时 41.72s |
| 3 | `resident_task_source_246_windows_test.go:126: the refusal did not come from the resident surface: text = "审批界面不可达，已 fail-closed 拒绝执行", want it to carry "常驻进程没有悬浮球窗口，卡片无处呈现"` | **代码**（常驻面走的是通用 fail-closed 文案，常驻专用那句在另一条分支上先返回了）；与面板正确性无关，是文案出口接线 | 新文件 `2071f59e`（NEW-in-push，+486 行）；同发日志上方逐字有 `approval: 卡片无处呈现（本进程没有悬浮球窗口），已 fail-closed 拒绝` |
| 4 | `panel_host_windows_test.go:544: cold bring-up measured on this box: -1.000 ms (HEAD 8ae4c23 at read time, 2026-10-01T16:09:20Z)` ／ `:546: cold bring-up did not produce a browser round trip (got -1.000); the message channel did not come up on the real window` | **口径**为主：文件头自己写着 `It runs on a private locked OS thread (hostThreadHarness) and only on this desktop: winlive has no CI job.`，而 `wisp-cli-tests.sh` 的 `-skip` 名单两发逐字相同（`TestDefaultDeadlineWallClockMeasurement|TestSubprocessCrashWriter|TestHelperProcess|TestLiveWasapiSmoke|TestRealDownloadVadThroughPipeline|TestRealDownloadPuncArchiveThroughPipeline|TestSyncRegistryProbeLive`），没人给它挂豁免 ⇒ 一枚"只在真桌面"的用例被推进了 CI 分母。次要可能＝负载（同发 cmd/wisp 包墙钟 424.7s -> 579.3s，runner 机器 1000001150 -> 1000001157） | 新文件 `697b4fae`（NEW-in-push，+912 行）；同文件 `TestAC4FocusReturnToPriorWindowGap33r5` 新发 `--- PASS: … (5.27s)`，`:753` 量到真实前台句柄迁移（`after Show 0x20204 | panel hwnd 0x20204`）⇒ 桌面是真的，缺的是 coldMs 那一站 |
| 5 | `pathresolver_budget_norace_test.go:34: C26 Resolve: 1241491 ns/op = 1.241 ms/op (budget 1.000 ms, 1902 samples)` ／ `:37: C26 budget breach: Resolve averages 1.241491ms per call, budget 1ms` | **负载/环境**（D32 毫秒预算，撞票面点名的那一族） | 基线同尺读数逐字 `C26 Resolve: 961784 ns/op = 0.962 ms/op (budget 1.000 ms, 1677 samples)` 且 `--- PASS: TestResolvePerCallBudget (1.81s)`；`git diff --stat 0589fd9c..8ae4c23e -- internal/risk/` **为空**，测试文件也未变（`ec2c8799`，区间无 delta）⇒ 被测代码一字未动，+29% 每调用耗时只能来自机器 |
| 6 | staticcheck 新 3 条（步骤 `lint#11`，逐字）：`cmd/wisp/subagent_selfapproval_197_test.go:580:2: should replace loop with verbs = append(verbs, planted...) (S1011)`、`:598:34: field holder is unused (U1000)`、`:600:6: type carrierWithDoor197 is unused (U1000)` | **代码**（新测试文件带的死形状；两发静态发现数 44 -> 47，其余 44 条逐字相同） | `subagent_selfapproval_197_test.go` = `a818df46`，NEW-in-push，+710 行；差集用 `comm` 对 `sc_new.txt/sc_base.txt` 算得，"只在基线有"一侧为空 |
| S1 | `panel_resident_windows_test.go:315: Resolve(entry): panel: embedded assets are not built (run npm run build in frontend/)` ＋ `:317: AC#13 has no subject in this tree: the embed resolves no entry …` ⇒ `--- SKIP: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (0.00s)` | **口径**（dist 未产） | `git ls-files frontend/dist` 只有 1 行 `frontend/dist/.gitkeep`；`tools/d22scan/runtests.sh:98-102` 逐字：`if [ "$skipped" -ne 0 ]; then … SKIP is not a pass (ticket 71 AC#3) … exit 1` |
| S2 | `panel_host_windows_test.go:794: no cold/hot sample recorded in this process: the lifecycle test did not run in this binary … Named skip` ⇒ `--- SKIP: TestPanelHostLatencyPercentilesAC2 (0.00s)` | **代码/口径连锁**：它是 #4 的下游（#4 死在 546 行 ⇒ 样本从没记上） | `portable-tests.sh: four numbers (all from -v output): === RUN=241 --- PASS=147 --- FAIL=11 --- SKIP=2` |

### 3.1 票面点名的"与面板正确性无关的既有形状"——逐一裁决
- **sherpa DLL 没摆好**：**两发均不成立**。真 `wisp doctor: PASS`，逐字含
  `[PASS] sherpa-onnx C API runtime 1.13.8 matches build pin (exe dir D:\a\wisp\wisp\build)`、
  `[PASS] onnxruntime.dll colocated file version 1.28.2.0 matches build pin 1.28.2`、两条 `[PASS] DLL colocated: sherpa-onnx-c{,cxx}-api.dll`。
  日志里另一处 `wisp doctor: FAIL`（含 `[FAIL] … does not match build pin unknown`）**两发逐字相同**，
  且它是 `TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/doctor` 这条 **PASS 的子用例**故意造的场景（`commit=unknown built=unknown`），不是环境缺件。
- **`TestHelperProcess`**：两发各出现 9 次，全部在 `-skip` 名单里，**两发都没有它的 `--- FAIL`** ⇒ 本次差集里它是 0 枚。
- **`TestSyncRedTeamRealOneDrive`（OneDrive 根为空）**：**两发都在**、逐字相同、且是 **SKIP 不是 FAIL**：
  `syncdirs_redteam_windows_test.go:220: no live sync root on this machine (detected roots: [])`。它是 windows scope 那 `SKIP=1` 的唯一成员。
- **`wisp slo` 的 300s 审批超时族**：**两发同形，非新增**——`--- FAIL: TestComposedGateBlocksAWriteForTwoSeconds (301.70s)`（新）
  vs `(301.12s)`（基线）。名字说 2 秒、实跑 301 秒＝它等的是审批超时。既有账，不该记在本差集头上。
- **D32 毫秒预算**：只有 #5 一枚，且归到负载（§3 第 5 行）。
- **core scope 的 4 枚前端契约红**：`internal/panel` 的 `frontend_hygiene_test.go`／`tokens_fourway_test.go`／`composer_test.go`／`approval_test.go`
  读的是 `frontend/**` 与 `design/**` 的漂移——**既有形状，本腿按票面没有打开那些文件**，只抄 CI 日志行。

## 4. 预测核对：错了，而且错得不少
预测＝"红名集合只多 1 枚"＝`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`，机制＝SKIP 被 `runtests.sh:98` 判红。
- 该枚在两发的状态：基线 **不存在**（`grep -c AC13 logs_base_all.txt` = 0）；新发 **SKIP**（不是 FAIL）。
- 机制对了一半、**这一发没有兑现**：那一步的退出走的是 `if [ "$rc" -ne 0 ]` 这条更早的分支——
  逐字 `runtests.sh: go test exited 1 - packages=[./cmd/wisp/ …] top-level: PASS=147 FAIL=11 SKIP=2`。
  11 枚 FAIL 先把步打死，`:98` 的 SKIP 规则**根本没轮到**。⇒ 它是"埋在 11 枚 FAIL 后面的第 12 枚"，FAIL 清零那天才会自己站起来。
- 事实层面（`frontend/dist` 只有 `.gitkeep` ⇒ embed 解析不到 entry ⇒ 具名 skip）预测**是对的**，日志逐字坐实。
- 行号差一点：预测 `panel_resident_windows_test.go:314`，CI 实际报 `:315` 与 `:317`。
- **错了 6 枚**：新增红是 6 枚测试名（#1–#6），不是 1 枚；其中只有 #4 与 AC13 同源于"面板宿主"这条线。
  另有 1 枚**预测外的第二个 SKIP**（`TestPanelHostLatencyPercentilesAC2`）和 **4 枚消失但其实是未采到**（test-core）。
- **作业级**恰好也是 +1（slo-full），但那是 runner 掉线（§1.4），与预测的机制毫无关系——这个"+1 对上了"是巧合，不能当验证。

## 5. WebView2 面板正确性那一格：读到了（独立读数，非旁证）
出处＝新发 run `36889094435` 作业 `test-windows` 步骤 #7（`bash scripts/wisp-cli-tests.sh`，`./cmd/wisp/`）。基线同发 **0 命中**（`grep -c 'AC#14' raw_tw_base.txt` = 0），所以这是本 push 新带进来的路径产生的**新读数**。逐字：
- `panel_resident_windows_test.go:825: AC#14 nail 1 (reply hop), page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)` ⇒ `--- PASS: TestAC14AwaitedBindingReplyReachesThePage (1.75s)`
- `panel_resident_windows_test.go:867: AC#14 nail 2 (Eval push hop), page's own words: title="PUSHED-33R5-OK"` ⇒ `--- PASS: TestAC14GoSideEvalPushReachesThePage (0.56s)`
- `wisp: panel window is up (test-harness, cold 1333.9 ms, hot path 0.0 ms)`（新发 5 次，基线 0 次）
为什么这算"运行库在"而不是"桩在跑"（读源码核过，非按常识填空）：这两枚走的是生产构造器
`NewPanelManager(disp, builtinAssetsOrTestNil(t), dataPath)`（`cmd/wisp/panel_resident_windows_test.go:59-67`），
执行 JS 的入口是 `cmd/wisp/panel_resident_windows_test.go:166-180` 的 `w := rp.mgr.currentWindow(); w.Eval(js)`——
真窗口的 `Eval`；页面再用 `window.wispDispatch(...)` 把话**回灌进 Go**，`awaitReport` 等的是回灌。
`"REPLIED"` 只有在页面里 `await` 的那个 Promise 被 Go 侧兑现后才可能出现，`title="PUSHED-33R5-OK"` 只有真脚本能改。
⇒ **这台云端 Windows 机器上 WebView2 Runtime 存在且能起真页面、能双向通信**。前一腿那句"装了"（出处＝issue 自述，旁证级）
到这里被一条独立读数顶上了。
**这一发没给的读数**：日志里没有任何版本串／`msedgewebview2.exe`／`WebView2Loader.dll`／`CreateCoreWebView2Environment`／注册表行
（`grep -icE 'NeverStealsFocus|CoreWebView2|msedgewebview2|WebView2 Runtime'` = 0 命中，全作业 `webview2|msedge` 只有 10 行，
逐字是 `go: downloading github.com/jchv/go-webview2 …` 与产品自己的中文错误文案）。⇒ **行为级证据＝在；版本号＝没有。**
**反向的一条、必须并列摆着**：同一作业同一发里 `TestPanelHostRealWindowHopAndLifecycle` 量到 `cold bring-up … -1.000 ms`
＝那一条路径上"浏览器来回"没成。所以"Runtime 在"不等于"宿主冷启动那一站在 CI 上也通"（§3 第 4 行）。

## 6. 下一步：三堆分明
**A. 要修（代码级，动产品码）**
1. `TestTicket223HandEditedFsLooseningCostsAnL2Card` — 唯一一枚"基线绿 -> 新发红"且测试文件未变 ⇒ 区间内 10598 行插入里找回归；先修它，别去改尺子。
2. `TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking` — 装配侧接 D45-2 会话授权：命中授权要出 `allow_session_grant`、要写 `tool_call.grant_id`、要出 GRANT-USE 审计行。
3. `TestAC246ResidentPipelineAsksThroughTheOneGate` — 常驻面要把自己那句"没有悬浮球窗口，卡片无处呈现"送到出口，而不是被通用 fail-closed 文案抢先。
4. staticcheck 新 3 条（`subagent_selfapproval_197_test.go:580/598/600`）— 新测试文件里的死字段/死类型/可展平的 append，删死形状即可（**不是**关 staticcheck）。

**B. 要登记（具名挂账级，本腿没动任何账文件）**
5. `scripts/portable-tests.sh:261` 的 `2>&1` 把 `go: downloading` 灌进 GUARD C 的分母（25 -> 35）——仪器口径 bug；
   它一次吃掉了 core scope 1365 个 `=== RUN` 的读数，并把 4 枚基线红伪装成"消失"。挂账名建议 `GUARD C stderr 并进分母（冷缓存即炸）`。
6. `TestPanelHostRealWindowHopAndLifecycle` 的自述 `only on this desktop: winlive has no CI job` 与 CI 分母不一致——
   要么给它一个真桌面作业，要么把这句改成能执行的判据。**具名登记，不是塞进 `-skip` 让它变绿。**
7. `TestPanelHostLatencyPercentilesAC2` 的 SKIP（`panel_host_windows_test.go:794`）＝上面那枚的下游连锁，跟着一起挂账。
8. `TestResolvePerCallBudget` 0.962 -> 1.241 ms/op、`internal/risk` 码未变（§3 第 5 行）＝CI 那台机器的负载账；
   两发的 runner 机器号不同（1000001150 vs 1000001157）是现成凭据。**阈值一字节都不许动。**
9. `lint#8/#9/#10` 三枚步骤被 #7 吃掉（含 `go vet` 两枚）；而 #7 的红句逐字是 `.scratch/wisp/probes/185/c1/mut/fs_broken.go`
   进了分母——被跟踪的坏变体和 gofumpt 的射程冲突，登记为"让 #8 那版受跟踪集合分母能真正跑起来"的前置。
10. `TestComposedGateBlocksAWriteForTwoSeconds` 301.7s/301.12s（两发同形）——300s 审批超时账，既有的，别记在本次推送头上。
11. 票面把 `36940536372` 说成"另一支 workflow"，实测是同一支 `ci` 的 schedule 触发、同 `headSha 8ae4c23e`、现已 completed/failure。登记这条前提更正。

**C. 要重跑才能判（本腿不重跑、不 cancel、不 rerun）**
12. `slo-full` 整作业：D32 六态/settle/leak 全部读数缺失（§1.4），要一次 runner 在线、日志可取的完整作业才能判。
13. core scope 全部 1365 枚（含 §1.3 那 4 枚"消失"的红）：要一次**模块缓存温热**的跑，GUARD C 才不吃掉分母。
14. `TestResolvePerCallBudget`：同一 SHA 背靠背跑两次（现成的 `36940536372` 就是同码样本），才能把"负载"与"码"分开钉死。
15. `go vet` 两枚与 `lint#8` 那版分母：#7 绿了才第一次有读数。

## 7. 我可能判错的条目（附"如果错了后果"）
- **#4 判成"口径为主"**：凭据是文件头那句自述＋`-skip` 名单未变＋AC#4 焦点跳数量到了真句柄。若其实 #6（33-r5 把面板搬进常驻进程）
  改坏了 coldMs 的记账路径，那"口径"就是"产品"的马甲——后果：把一个真宿主回归记成挂账，下次没人再看 coldMs。
- **#5 判成纯负载**：`internal/risk` 区间 diff 为空，但 `Resolve` 的开销可能间接吃 `internal/config` 的默认值（区间内改了 10598 行）。
  若如此，后果＝D32 预算被产品码悄悄推过线，而我用"机器慢"替它开了脱。
- **#1 判成"真回归"**：只证明"测试文件未变 + 基线 PASS + 新发 FAIL"，没定位到哪一枚码改的。若其实是该用例自身依赖的
  时序在 CI 上抖（2.13s -> 2.43s），后果＝把一枚 flake 记成回归，去查不存在的那条线。
- **§5 判成"Runtime 在"**：只验到"真窗口 + 真 Eval + 回灌"。若 `go-webview2` 在无 Runtime 时能退化成别的实现，
  后果＝给"装了"这条旁证盖了一个不该盖的章；版本号那一格我明写了没有。
- **§1.3 的 6/4/8 三个数**：尺子是 `--- FAIL:`/`--- SKIP:` 行 + runtests.sh 自印的四个数。若 gh 的作业日志拼接有丢行（本次 slo-full 整作业就没有），
  后果＝差集偏小。已用两个独立尺子互校（我 grep 的名字数 与 runtests.sh 的 `FAIL=11`/`SKIP=2` 逐枚对得上），但 slo-full 无法这样互校。

## 8. 判不动的地方（甲＝补得上，乙＝补不上，不做）
- 甲：`TestTicket223HandEditedFsLooseningCostsAnL2Card` 到底哪一枚码改的——`git log -L` 逐文件比对区间内 32 枚改动文件，谁都能做（本腿没做，因为只要给方向不要定罪）。
- 甲：GUARD C 到底是"10 行噪声"还是"真多了 10 枚包"——把 `go list ./...` 的 25 条基线名单与本次日志里的名单并排点一遍即可（我已用 `internal/risk` diff 为空＋基线 `resolved=25` 侧面坐实）。
- 乙：`slo-full` 这一发到底差在哪个阈值——作业日志 `log not found`、check-run output 三个字段全 `null`，只有 runner 掉线一条注解；**不重跑就没有任何读数**。不做（⛔ 票面禁重跑）。
- 乙：WebView2 的具体版本号——日志里没有，`msedge --version` 不在任何步骤里。不做（⛔ 不许靠常识填空，也不许在取证腿里跑任何东西）。
- 乙：36940536372 的红名册——**票面明令不混入本次比对**，所以即便它是同码现成对照，本腿也没展开。
