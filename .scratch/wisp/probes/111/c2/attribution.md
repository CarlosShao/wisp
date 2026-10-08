# 111-c2 · 取数＋归因件（票 111「CI 上这 6 枚红」）

- 腿：`111-c2`（非实现者、只读取数＋归因；⛔ 未碰产码/测试/`.github`，⛔ 未开真窗，⛔ 未 push）
- 起手锚：`.scratch/wisp/probes/111/c2/00-anchor.md`（**已单独一笔进仓＝`aba879c6`**）；锚上现量 HEAD `052b393f`（＝派单那枚，"不早于"成立）、`origin/dev..HEAD` **312**、`cnb/dev..HEAD` **492**
- 落笔时刻：`2026-10-08 11:4x +0800`
- 大输出全落**仓外** `D:/tmp/wisp111c2/`（8 发 run 日志 1.8–2.2 MB／发），仓内只落 `.md`；⛔ 未把整发日志读进上下文
- 写面自报：本腿**只新建 `.md`**（本件＋`logs/*.md` 7 枚＋锚），⛔ 零 `.sh`/`.ps1`/`.txt`/`.out`
- Go 面自报：⛔ 未跑 `go build`/`go vet`/`go test`；**跑过 `go env GOMODCACHE` 一发**（rc=0，只为定位依赖库源码路径，读的是模块缓存里的 `webview.go`，⛔ 不涉编译）；未跑 `go list`

证据件（逐枚自带尺与 rc）：`logs/01-run305-red-sentences.md`｜`logs/02-cross-run-colour-matrix.md`｜`logs/03-roster-rerun-vs-33v4.md`｜`logs/04-local-colours.md`｜`logs/05-minus-one-exits.md`｜`logs/06-ci-and-carrier-facts.md`｜`logs/07-webview2-runtime-evidence.md`｜`logs/08-dedupe-scan.md`（八枚，字节数见交付时的 `ls -la`）。

---

## 问 1 · 逐名归因（六枚＋`TestPanelHostLatencyPercentilesAC2`）

红句原文全量在 `logs/01`；本机色全量在 `logs/04`；三方对照表在 `logs/03`。逐枚一行：

| 用例 | CI run 305 红句要点（原文在 `logs/01`） | 它要的那件能力 | 本机色（10-08 08:18 那发整包，件 `probes/35/r4/logs/orch/full-package-after-r4.txt`，锚 `bcfb459b`） |
|---|---|---|---|
| `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey` | `always_write_no_clobber_226_test.go:91` 盘上存 `.../runneradmin/.../003`、期望 `.../RUNNER~1/.../003` | **不要窗、不要 DLL、不要运行库**：要的是"同一目录只有一种拼法"——`%TEMP%` 在这台机上是 **8.3 短名**、产品存的是**长名** | **PASS** (1.80s)；本机整包件里 `RUNNER~1` 命中 **0**（`grep -c` rc=1）、`C:\Users\swq` 263 次 ⇒ 本机根本没有短名这一形 |
| `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` | `approval_always_201_test.go:141` `config.toml did not gain the stored rule "C:/Users/RUNNER~1/..."` | 同上（短名 vs 长名） | **PASS** (1.32s) |
| `TestTicket223PermissionDeniedSitsInItsOwnSentence` | `config_reload_perm_223_windows_test.go:95` `attempt 1: icacls denied nothing, so this case cannot show the permission sentence` | 要的是**这台机的 ACL 真能拒绝一个 `Everyone:(R)` deny**（用例的正控：`Stat` 成功、`ReadFile` 必须失败）；⛔ 不要窗 | **PASS** (2.54s) |
| `TestRunPacketCarriesTheLoadedInstructionFiles` | `instructions_200r2_test.go:167` 同一行里并排两种拼法（`RUNNER~1` vs `runneradmin`，Bytes 都是 269） | 短名 vs 长名（同第 1、2 枚一族） | **PASS** (1.39s) |
| `TestPanelHostRealWindowHopAndLifecycle` | `panel_host_windows_test.go:660 → :662`，`got -1.000`（**未测量哨兵**，见问 3） | 真窗＋WebView2 通道＋`wispProbeRT` 绑定那一跳在 **5 秒单调截止内**完成 | **FAIL** 同名同句（`-1.000`，5.06s）⇒ **这枚不是"只有 CI 红"**；本机 10-07 起连发四件都是 `-1`（`logs/04`） |
| `TestAC14GoSideEvalPushReachesThePage` | `panel_resident_windows_test.go:869` `the page reports its title as "", want "PUSHED-33R5-OK"`——**页报回了、报的是空标题** ⇒ 要的是 **Go→页的 `Eval` 落地**（不是页→Go 的门，门是通的） | 同上（真窗＋Eval push） | **FAIL** 但**红在另一出口** `:866`：`no report "ac14-push" from the page within 15s (what DID arrive at the door: nothing at all)` ⇒ 本机连页→Go 都没回（20.01s） |
| `TestPanelHostLatencyPercentilesAC2`（派单追加） | run 305 是 **SKIP**（`:985` `no cold/hot sample recorded in this process`）；**run 300 不是"换成"一个新因**，它红在 `:1005 cold P95 3743.192 ms over 1 runs exceeds the D32 panel cold budget 1500 ms`（hot 侧 25.996 ms 在 200 ms 预算内） | 要的是"同一次运行里 lifecycle 真把样本记上"＋**那台机器的冷起能进 1500 ms** | **SKIP**（本机同发，同因：lifecycle 终止在 `:662`） |

### 名册作差（本机窗口依赖名册 vs CI 红名册）

本腿按第 108 条**自己复跑了那把尺**（`logs/03`）：**15 枚**（`33-v4` 那件是 14 枚；差的恰一枚＝`TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive`，`winlive` 档，两边名册都只在"我这多"那侧）。
- **两边同名（红）**：`TestPanelHostRealWindowHopAndLifecycle`、`TestAC14GoSideEvalPushReachesThePage` ⇒ 2 枚。
- **只在 CI 红那侧**：另外 4 枚（短名族 3 ＋ `icacls` 1）**根本不在窗口依赖名册里**（它们体扫零命中，本腿复跑确认它们不在 15 枚内）⇒ **"6 枚红都属窗口族"这个读法不成立**。
- **只在窗口名册那侧**：13 枚；其中 **8 枚在 CI 那台机器上是 PASS**（`logs/03` 的表），**2 枚本机红而 CI 绿**（`TestAC14AwaitedBindingReplyReachesThePage`、`TestAC4FocusReturnToPriorWindowGap33r5`），3 枚 winlive 档两侧都**不编译**。

rc=0

---

## 问 2 · 托管那台 `windows-latest` 上 WebView2 Runtime 有没有

**答：有，而且不止一档证据；但"有运行库"与"看得见的窗"是两格。** 全量在 `logs/07`。

- workflow 侧：`grep -in 'webview\|edge\|evergreen\|runtime' .github/workflows/ci.yml` ⇒ 只有 2 行散文注释（`:445`、`:755`），**没有任何安装/等待运行库的步骤**（`logs/06` §2）⇒ CI 上有没有它**不由 workflow 决定**。
- 红句本身算哪一支：**"有运行时但通道没起"**。判据是产码给"没运行时"那支留了名字：`cmd/wisp/panel_host_windows.go:394-396` `runtime missing or blocked` ⇒ 那支会红在 `panel_host_windows_test.go:641-642`；**CI 实际红在 `:662`**，且同发 `:647 IsCreated()`、`:650 hwnd != 0` 都没响。
- 仓里已量过（先例都在）：`33-n1`（`probes/33/n1/verdict.md` §4）结论逐字"托管 `windows-latest` **有** Evergreen WebView2 运行库"，⛔ 不许用甲推乙、R-a…R-e 五枚读数列出但未采；台账 `A526` 已从**一发真推送的 CI 日志**把甲升为【已证】（`panel window is up ... cold 231.2 ms` ×5 ＋真 hwnd ＋ `TestAC4Focus...` PASS）。
- 本腿新取（⛔ 不是二手）：run 294–304 **七发**逐字带 `our tree webview=7 (machine-wide 7)`／`machine-wide msedgewebview2=7`、`same HWND ... across hide->re-show=true`、`hot re-show 38.649 ms`；run 305 里 `TestAC4FocusReturnToPriorWindowGap33r5` **PASS** 且 `after Show 0x10204 == panel hwnd 0x10204`。⇒ **子进程真起、窗真建、前台真归它**。

**这一格今天取不到的（具名，欠哪一发读数写清）**：
1. **R-a**（进程内 `GetAvailableCoreWebView2BrowserVersionString` ＋HKLM/HKCU 两枚 `pv` 键的**当场版本串**）——**零读数**；欠的是一发带采集步骤的真实 run（⛔ 本腿无权推、无权改 `.github`）。
2. **R-c**（对窗 DC `BitBlt`/`PrintWindow` 落 PNG ＋非纯色字节数）——**零读数**；"画没画"到今天仍判不了。
3. **R-e**（同一 job 连跑 10 次的冷起分布 min/med/p95 ＋本机对照）——**零读数**；今天只有单发 `-count=1` 的 3 126–3 941 ms。
4. **run 305 那一发为什么从正读数漂到 `-1`**——同 SHA 的 run 304 给出 3414.339 ms 正读数 ⇒ 本腿**没有**能区分这两种的读数；欠同一枚 SHA 上多发连跑。
5. `winlive` 那 3 枚（`TestLive35v2PagePostMessageEnvelopeReachesTheGoDoor`／`TestTicket255RealWindowWidthFollowsTheConfig`／`TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive`）在任何一侧都**零色**——它们今天连编译都没被 CI 查过（`logs/06` §4：winlive 步只存在于未推 SHA，`git show cc31526:...|grep -c winlive`＝0，rc=1）。

rc=0

---

## 问 3 · `got -1.000` 这个哨兵值

全量在 `logs/05`。锚 `cmd/wisp/panel_host_windows.go`（本腿 HEAD `052b393f` 现量，`grep -n 'return -1'` 恰 4 行）：

| 出口 | 行号 | 触发条件 |
|---|---|---|
| E1 | **845** | `w == nil`（`:842` 从 `m.w` 取） |
| E2 | **860** | `w.Bind("wispProbeRT", func() string {...})` 返错 |
| E3 | **878** | `ctx.Err() != nil` |
| E4 | **882** | `time.Now().After(deadline)`，`deadline := t0.Add(5 * time.Second)`（`:868`；`t0` 来自 `bringUp:319`，**开窗之前**起算） |

值只带一个 `float64`：`coldStartPageHandover:448` → `bringUp:421-425` 写进 `m.lastColdMs` → 访问器 `:285` → 测试 `panel_host_windows_test.go:659` → `:662 t.Fatalf(... got %.3f)`。**没有 error 位、没有原因栏。**

**CI 那发红句能不能区分是哪一条出口？——红句本身不能；缺的那一栏是"哪条出口（E1/E2/E3/E4）"。**
但把红句＋同发别的读数＋码摆在一起能**排掉三条**（⚠ 这是推理、不是仪器读数，逐条依据在 `logs/05`）：E1 被同发未响的 `:647/:650` 与 `:413-417` 的赋值顺序排除；E3 被测试的 ctx 形状排除（`:625` Background＋defer cancel，`hostThreadHarness.bringUp:110-112` 无 `WithTimeout`）；E2 **静态不可达**（`go-webview2` 的 `Bind` 只在"不是 func"或"返回值 >2"时返错，`webview.go:450-457`，而实参是 `func() string`）⇒ **只剩 E4＝5 秒单调截止到期**。

口径自守：`-1` 是**未测量哨兵**，⛔ 不许读成"时延 -1 毫秒"，⛔ 不许当"环境没有运行时"的证据；`probes/35/v6/live-wave.md:44` 那个 `-1` 本件**一格都没引**。
另具名一条形状差：**8 发里只有 run 305 是 `-1`**，其余 7 发（含 run 300）是"量到正数、超 1500 ms 预算"那一支（`logs/02`）——`111-c1` 把 run 300 也写成 `-1`，本腿顶回（见文末顶回第 1 条）。

rc=0

---

## 问 4 · 要什么事实才能定"这 6 枚在 CI 的期望色"（⛔ 不选形）

### ⓐ 视为真缺陷要修 —— 它需要的事实前提
1. 每一步红都来自**产品答错**而不是环境形状。本腿读数把这条**先切掉一半**：6 枚里 4 枚的红句是环境形状（3 枚短名拼法＋1 枚 `icacls` 拒不了），它们的**判据本身在两台机器上表现不同**（`logs/03`/`logs/04`）。
2. **即使 6 枚 FAIL 全修，那一步也不会绿**：运具 `tools/d22scan/runtests.sh:98-102` 逐字 `SKIP is not a pass ... exit 1`，而 run 305 同发带 **2 枚 SKIP**（`TestPanelHostLatencyPercentilesAC2`、`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`）。⇒ ⓐ 的完整前提必须**连这两枚 SKIP 的成因一起解**（其中一枚的成因是结构性的：`test-windows` 里没有 `npm run build` 步，`logs/06` §3）。
3. 修短名那一族要有**能复现它的机器**：本机 263 次长名、`RUNNER~1` 零命中（`logs/04`）⇒ 本机今天复现不了；能复现的只有托管那台 ⇒ 每修一发都要推一次才有色（⛔ 本腿未推）。
4. 归口已存在（⛔ 不新立）：短名机制＝票 252（未 `-done`）；分诊总框＝票 70（`blocked-on-owner`）。

### ⓑ 视为环境不支持、该带 skip 兜底 —— 它需要的事实前提
1. **本仓"带 skip 兜"的规矩是反的**：`runtests.sh` 头部规则 2（`:22-26`）逐字——"legitimately cannot run on this platform must be **moved out of the step's scope, not silenced here** (that would be lowering a threshold)"；`:98-102` 实现。⇒ **在 `cmd/wisp` 这一步里加 `t.Skip` 换不来绿**，只会把 FAIL 换成 SKIP 红。
2. 既有先例（同一机制链，已落地）：`cmd/wisp` 那一步**现在就带着一份具名 `-skip` 名单**（7 枚，CI 日志逐字在 `logs/06` §5），它背后的台账是 `scripts/portable-tests.sh` 里**11 枚具名行**（每行带 scope／平台／理由／`next=` 归属），脚本自写 `An unaccounted SKIP also still goes red, so a new skip cannot be parked here by accident`（`:32-33`）；规矩的票面出处＝**票 71 AC#3**（`[x]`）。⇒ ⓑ 的真实形状**不是"给测试加兜"，是"往那 11 枚台账里加具名行"**，且这属运具层（⛔ 本腿没动）。
3. ⛔ 本腿现量：**这 6 枚红名一枚都不在那份台账里**（`grep -c ... scripts/portable-tests.sh`＝0，rc=1）⇒ 今天它们没被任何一档"具名环境豁免"覆盖。
4. 同族判据原文（⛔ 与 ⓑ 直接冲突的那一支）：`internal/ball/live_guard_windows_test.go:170` 逐字 `t.Skip is not an option here.`；`scripts/wisp-cli-tests.sh:60-67` 非 windows 直接 `exit 2`（"Not a skip and not a pass"）。⇒ 走 ⓑ 必须先裁"这两处原文要不要收窄"。
5. 台账里那枚纪律（派单提到的"winlive 无 skip 兜 ⇒ 红是环境红"）原文＝票 35 `:272`：`cmd/wisp` 那 7 枚 winlive 件里**没有一枚用 `t.Skip` 兜**（腿现量 `t.Skip`＝0、只有 `t.Fatalf`）⇒ **那一格说的是 winlive 那 7 枚，⛔ 射程不覆盖今天这 6 枚**（这 6 枚在默认 `windows` 档，`logs/03` 的 tag 列逐枚可查）。引用它之前要把射程写明。

### ⓒ 视为该从 CI 摘出、归本机量 —— 它需要的事实前提
1. **摘出的手段有哪两把**：改 workflow 那一步的 scope（＝动 `.github`，⛔ 本腿无权、且按 `A523` 属契约级），或往 `portable-tests.sh` 的具名台账加行（运具层）。作业级 `if:` 被 D22 mode-6 禁（票 35 `:270-271` 逐字；本腿复量 `^    if:`＝0，rc=1，与那句一致）。
2. **"归本机量"这条路今天的真实缺口**：本件**本机色全部是引用**，⛔ 本腿没复跑（派单⛔ `go test`；`logs/04` 已按第 108 条具名写成"它取过、不是我为真"）。⇒ ⓒ 的前提是**有一枚持有测试面的腿在 `052b393f` 上重跑整包**，逐名对上 `logs/03` 那张表。
3. **ⓑ/ⓒ 与"红就是环境红"在本仓的先例**：`A451:9564`（icacls／短名＝环境，逐字）＋`A617:12088`（"runner 机器耦合…⛔ 不新立案"）＋`A526:10718`（`cold=3499.615 ms` 是**那台机器自己的时序**，⛔ 不许当我们的性能读数）。⇒ 这三处是**已经落成文的处置**，不是要新造的规矩。
4. **⚠ 两枚必须先摆清的事实，⛔ 本腿不裁**：(i) 窗口族在 CI **不是"跑不了"**——8 枚窗口依赖用例同发绿、还有 `msedgewebview2` 子进程计数（`logs/07`），所以"摘出"的理由只能是**时序/像素档不可复现**（`33-n1` §6.2④"CI 档里任何毫秒级阈值都不可复现"），不能是"没有运行库"；(ii) 同一枚 `TestPanelHostRealWindowHopAndLifecycle` **本机也红**（10-07 起四件 `-1`），摘进本机档**不会让它绿**。
5. 步级 `if: !cancelled()` 那条纪律（派单点名要写的先例）：本腿复量 HEAD＝**13 枚步级 `if:`（12 枚 `!cancelled()` ＋ 1 枚 `always()`）、作业级 0**（`logs/06` §4）。"在 yaml 里"≠"跑过"的实证：`go vet (module)`／`go vet (tools/d22scan module)` 12 发全 `[skipped]`＝从未执行（`111-c1` 的普查，本腿未重跑 gh 普查、只复用了它的结论并具名）。⇒ 任何新增步**不带步级 `if:` 就是一个死格**。

rc=0

---

## 问 5 · 红名册的稳定性（⛔ 未把大日志读进上下文）

取数：本腿自己 `gh run view --log-failed` 拉 **6 发**（304/303/302/301/299/294，逐发 rc=0、err 文件 0 B），加 `111-c1` 落在仓外的 305/300 两发原文**由本腿重切**⇒ 共 **8 发**。逐发字节、切片行数、命令全在 `logs/02`。

- **稳定集（8 发全红）＝2 枚**：`TestRunPacketCarriesTheLoadedInstructionFiles`、`TestPanelHostRealWindowHopAndLifecycle`（⚠ 后者**同名两种坏法**：305 是 `-1`／294–304 是"正数超 1500 ms 预算"，`logs/02` 逐发列出）。
- **近 7 发稳定、最早那发绿＝2 枚**：`TestAC1AlwaysBranchDoesNotRevertAHandEditedKey`、`TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card`（run 294 **PASS**，299–305 FAIL）⇒ 这一族是在 `941805d`→`c6cf66e` 之间**进入**名册的。
- **`TestTicket223PermissionDeniedSitsInItsOwnSentence`＝8 发全红**（唯一一枚 8/8 都红且只有一种红句的）。
- **漂移枚**：
  - `TestAC14GoSideEvalPushReachesThePage`：FAIL(305,304,303)／**PASS(302,301,300,299)**／FAIL(294) ⇒ **漂**。
  - `TestPanelHostLatencyPercentilesAC2`：**FAIL→SKIP 互换**（294–304 FAIL、305 SKIP），机理不是随机：它由 lifecycle **怎么红**决定（红在 `:665` 是 `Errorf`→样本记上→它跑分位并红；红在 `:662` 是 `Fatalf`→样本没记上→它 SKIP）。⇒ **这两枚是同一个因的两个果**。
  - `TestTicket223HandEditedFsLooseningCostsAnL2Card`：只在 303、299 两发红（本机同族 `-count=3` 曾 3/3 绿，`A451:9626`）。
  - 名册**规模**也在漂：`=== RUN` 261→323→335，`FAIL` 12→6→8→7→6。
- **`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`＝8 发全 SKIP**（同因：`built=false entry-bytes=0`，`logs/06` §3）⇒ 它是**结构性的、不是漂的**。

**答：不是稳定集。** 6 枚里**只有 3 枚**在"8 发全红"这一档；派单里"run 300 的后两枚换成…"那句**部分成立**——换上去的 `TestPanelHostLatencyPercentilesAC2` 不是新因（见上），而 run 300 的 lifecycle 与 run 305 的 lifecycle **红在不同行、不同坏法**。
归口（⛔ 不新立）：**票 236 标题第 5 格**逐字"**在册红名册不是稳定集合**"＋台账 `A451:9626`"**'在册红名册是稳定集合'这句话今天第三次被推翻**"，本腿是第四次实证。

rc=0

---

## 问 6 · 有没有既有票/框已经裁过这件事

**有，四条机制链全部有名字；⛔ 不新立任何票。** 扫描逐变体（含"CI 上没有 WebView2／真窗不进 CI／窗口依赖名册／环境红／skip 兜底／仅本机可量／machine-local"七组字面 ＋ 命中数与 rc）在 `logs/08-dedupe-scan.md`。归口：

| 机制链 | 已在哪一格 |
|---|---|
| 8.3 短名 ↔ 长名分家 | **票 252**（未 `-done`，10-02 立）；先例件票 115（`-done`）、票 106/112（`-done`，winsec 在 CI runner 首跑红） |
| 环境红／runner 机器耦合 | 台账 **`A451:9564`**（`icacls denied nothing`＝托管 runner 环境、`RUNNER~1`＝环境）＋**`A617:12088`**（"其余 18 枚 → runner 机器耦合…⛔ 不新立案"，并逐枚点名 6 枚短名／2 枚面板冷启超预算／1 枚 icacls）＋票 260 `:20` 与 `:11677` 那枚"既有环境红"先例 |
| 真窗／仅本机可量 | **票 35 `:267-272`**（10-08 08:5x 裁定 `:75(c)` 走〔仅本机可量〕、⛔ 不搬真窗进 CI）＋**票 111 `AC#11` ⓒ**＋产码注释逐字 `winlive has no CI job`（`panel_host_windows_test.go:605-609`、`:649`；`panel_host_windows_live_test.go:36`） |
| 托管机有没有运行库（甲）与"能不能真开可见窗"（乙） | **`33-n1`** §4／§5／§6（甲=装了、乙不升档、R-a..R-e 待采）＋台账 **`A526`**（甲升【已证】、`cold=3499.615 ms` 不许当我们的读数） |
| 红名册不稳定 | **票 236 第 5 格**＋`A451:9626`（第三次推翻）＋票 70（CI 分诊总框，`blocked-on-owner`） |
| SKIP 不是绿（"skip 兜"的规矩） | **票 71 AC#3**（`[x]`）＋`tools/d22scan/runtests.sh:22-26/:98-102`＋`scripts/portable-tests.sh` 11 枚具名台账 |

**唯一今天找不到既有格子的一格**（⛔ 本腿不立、只具名）：`TestAC14GoSideEvalPushReachesThePage` **CI 红在 `:869`（页报回空标题）而本机红在 `:866`（页零报回）**——"同一枚用例在两侧坏在不同出口"这个具体形状，账上只有尺（`A617`／票 35 那格"红名册不区分坏法要连文案引"）、没有归口票。
另一处具名缺口：`icacls`＋`runner` 同现的票有 9 枚，但没有一枚名字里带 `TestTicket223PermissionDeniedSitsInItsOwnSentence`（逐枚 `grep -c` 全 0，rc=1）⇒ 它今天只有台账 `A451:9564` 的归因、没有票面格子。

rc=0

---

## 具名顶回（⚠ 派单原文与 `111-c1` 的件，凡与盘上/CI 原文冲突一律以原文为准）

1. **`111-c1` 与本腿派单把 run 300 的 lifecycle 红写成"`got -1.000`；the message channel did not come up"** ⇒ **只有 run 305 是那一形。** 盘上原文：run 300 `panel_host_windows_test.go:660: cold bring-up measured on this box: **3743.192 ms**` → `:665 cold bring-up 3743.2 ms exceeds D32 panel cold budget 1500 ms`；299/301/302/303/304/294 **同形**。件 `ci-colours.md` Q5 第 177 行与票 111 追加段第 392 行都写成"两发都是 `-1`"⇒ **两处需更正**（本腿逐发红句在 `logs/02`）。
2. **派单："6 枚红里有两枚是面板那族"⇒ 成立；但剩下的四枚不是一族产品缺陷。** 逐字读数：3 枚是 `RUNNER~1`↔`runneradmin` 同目录两种拼法、1 枚是 `icacls denied nothing`（用例自己的正控失效）。⇒ "一枚待修的缺陷"这个单数前提**不成立**，至少是 3–4 条互不相干的机制链。
3. **票 35 `:270` 那枚裁定理由（"每次 push 都会在机主屏幕上开真窗"、依据是 `gh api` "只有这台机在线"）与 `A701/A703` 里我承认作废的"唯一 runner"是同一族** ⇒ 本腿 `runs-on` 现量（两棵树同形，`logs/06` §1）确认**只有 `slo-full` 自托管**、`test-windows`/`slo-smoke` 是托管。⚠ 但那枚裁定的**结论**不因此自动翻：作业级 `if:` 被禁 ⇒ 挂到 `slo-full` 那台确实每次 push 会开机主屏幕；**挂到托管那台则是另一道题**，而"托管那台能不能真开可见窗"是 `33-n1` §5.2 明写"⛔ 不许用甲推乙"的**乙**，今天仍欠 R-b/R-c。⇒ 本腿只报"第 1 条理由的前提是错的、第 2/3 条不依赖它"，⛔ 不翻裁定。
4. **`33-v4` 的窗口依赖名册 14 枚** ⇒ 本腿在同仓现量（HEAD `052b393f`）**15 枚**，多出的恰是 `TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive`（`winlive` 档）。⛔ 不是它错，是**两枚锚之间落过新件**；引用必须带 ref（第 130 条）。
5. **派单："台账里有一枚纪律 `winlive` 无 skip 兜 ⇒ 红是环境红"** ⇒ 原文（票 35 `:272`）射程只覆盖 **winlive 那 7 枚文件**；今天这 6 枚里**没有一枚是 winlive 档**（`logs/03` 逐枚 tag 现量：12 枚 `windows`、3 枚 `windows && winlive`，6 枚红全在 `windows` 档里跑）。⛔ 别把它当覆盖这 6 枚的判据。
6. **`A703`/`A692` 那串"步级 `if:` 12 处"** ⇒ 本腿 HEAD 现量 **13**（12 `!cancelled()` ＋ 1 `always()`），`^    if:`（作业级）**0**。两数分属"加步前／加步后"两个时刻（票 111 `:375` 已把这族坑写过），⛔ 压成一格就是错。
7. 与派单**一致、本腿复量成立**的三句（照实报，不冒充顶回）：`runs-on` 三行那句、`cmd/wisp` 那步真跑且红（`RUN=335 PASS=230 FAIL=6 SKIP=2`）、run 300 前 4 枚同名＋后两枚换成 `TestPanelHostRealWindowHopAndLifecycle`／`TestPanelHostLatencyPercentilesAC2`。

rc=0

---

## 本腿自报（纪律与射程）

- ⛔ 未 `go build`/`go vet`/`go test`；**跑过 `go env GOMODCACHE` 一发**（rc=0，定位依赖库源码只读用）；未跑 `go list`。
- ⛔ 未改产码/测试/CI/台账；`.github/**` 只 `grep`＋`git show`。票面 111 **只在文末追加一节**；追加前后未勾框枚数尺＝`grep -cE '^[[:space:]]*- \[ \]'` **追加前 5／追加后 5**，`- \[x\]` **6／6**（写进票末那一节的末尾）。⛔ 未翻、未增、未删任何勾。
- ⛔ 未 push、未开真窗、未在仓库目录内建 worktree/checkout。Git：`git add -- <显式 path>` ＋ `git commit -- <显式 path>`，⛔ 无 `-A`/`.`/`--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`。
- `gh` 全部一次通过（`run list` 一发 rc=0；6 发 `--log-failed` 逐发 rc=0、err 文件 0 B）⇒ **本腿没有需要重试的失败**，也没有据失败下"取不到"的结论。
- 大输出：`D:/tmp/wisp111c2/`（`run*.log` 8 枚 1.8–2.2 MB／`cli-run*.txt`／`clean-run*.txt`／`roster-rerun.txt`／`file-tags.txt`／`ci20.json`），⛔ 未入库、⛔ 未删（临时件只建不删）。
- 仓外 `D:/tmp/wisp111c1/`（`111-c1` 的原料）本腿**只复用了两发原始 run 日志并自己重切**，⛔ 没有引用它的切片结论当本腿读数。
- 一件诚实的射程限制：**本机色全部是引用、本腿没能复跑**（派单⛔ `go test`）——这一格具名记在问 1／问 4 ⓑⓒ／`logs/04` 三处，欠的那一发＝持有测试面的腿在 `052b393f` 上重跑整包。

rc=0
