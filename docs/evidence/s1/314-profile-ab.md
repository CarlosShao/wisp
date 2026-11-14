# 票 314 AC#1/AC#2 — 固定名 profile 夹的归因 ＋ A/B 串行实测（测量腿，非写腿）

> 裁决人：314 测量腿（非实现者，⛔ 未翻任何 AC 框）。时间 `2026-11-14 11:08–11:26 +08`。
> HEAD `15cf0f94b`。机器：本机 Windows（go1.27.1 windows/amd64；跑时 `Get-CimInstance Win32_Processor` 读数 100%，
> 另有两枚常驻产品进程 `wisp` pid 21580 与 `wisp-live` pid 69456（后者由 `.scratch/wisp/stage-cli/wisp-live.exe`
> 于 10:48:40 启动、持 `%LocalAppData%\wisp\orb\69456-c44f3d40\boot-1`），两者都不命名任何测试 profile 夹）。
> 读数根：`.scratch/wisp/probes/314-profile/**`（`results.jsonl`、`results-tight.jsonl`、`results-udf-share.jsonl`、`logs/**`）。

## 0. 一句话结论

**固定名 ⇒ 两发互踩同一个 profile 目录：坐实（本件给了 file:line 与"同夹两个浏览器实例并存"的实测）。但"共享夹状态"不足以解释 139/挂起：串行可及的最强形状下 15 发全绿、零 native 异常。它只在「两个进程同时在 `CreateEnvironmentWithOptions` 里争同一夹」那一瞬才可能致命 —— 那是并发窗口，不是共享夹本身。**

---

## 1. 逐字：那处夹具、那段注释（`68275983e^` 版本，即 de-share 之前）

`cmd/wisp/panel_resident_windows_test.go` 在 `68275983e^` 的行号与内容（现跑 `git show 68275983e^:cmd/wisp/panel_resident_windows_test.go`）：

```
  55: // startPanelForTest brings up the production thread with a production-shaped
  56: // manager (real embedded assets, a real user-data folder under the system temp
  57: // dir, the real fail-closed dispatch chain) and returns it with the cleanups
  58: // already registered.
  59: func startPanelForTest(t *testing.T) (*residentPanel, *PanelManager) {
  60: 	t.Helper()
  61: 	dataPath := filepath.Join(os.TempDir(), "wisp-33r5-panel-profile")
```

`:61` 就是 AC#1 说的那处**固定名**夹具；它上方的注释逐字抄在上面（四行）。

**判（引用层与判断层分开）**

- 引用层：这段注释解释的是三件事 —— 夹具用的是**生产形状**（真 embed assets、真 user-data 夹、真 fail-closed 派发链）、夹在**系统 temp 下**、`t.Cleanup(rp.stop)` 已注册。**它没有一个字解释"为什么是固定名"。**
- 引用层（另一处，唯一写"故意不删"的地方）：`cmd/wisp/panel_host_windows_test.go:604-608` 逐字 ——
  "The WebView2 user-data folder is created here and deliberately NOT deleted (repo rule: throwaways are only ever created).
  t.TempDir's auto-RemoveAll races the browser's own profile-DB teardown and fails on a locked file - that is the teardown,
  not a product assertion." 那是"为什么**不删**"，仍不是"为什么**固定名**"。
- 判断层：固定名是 33-r5 首落（`13acad460` "feat(panel-host): 33-r5 land the panel in the resident process"）顺手写下的名字。
  现跑 `git log -1 --format=%B 13acad460 | Select-String "profile|user-data|temp"` = **零命中**，commit body 里也没有理由。
  ⇒ **AC#1 括号里"票面记着 `:61` 处注释解释原因"不成立：那里没有原因。**固定名是首落的字面量，不是被论证过的设计。
- 现状（HEAD `15cf0f94b`）：`:61` 已换成 `uniqueProfileLeaf`（`68275983e`，2026-11-11 19:22），
  把 `:61`/`:581`/`:658` 三处固定名改成 `<leaf>-<pid>-<nanotime>` 的每发唯一叶，并注册 `t.Cleanup(os.RemoveAll)`；
  其注释逐字（现状 `:55-60`）已经自己写明了理由：
  "only the leaf is made unique with pid + nanotime so runs and processes stop sharing one fixed-name WebView2 user-data folder."

## 2. 因果链（`文件:行`）

**夹怎么变成 WebView2 的 user-data 夹**

| 跳 | 位置 |
|---|---|
| 夹具把 `dataPath` 交给管理器 | `cmd/wisp/panel_host_windows.go:331-332`（`NewPanelManager` 存 `m.dataPath`） |
| `bringUp` 用它建控件 | `cmd/wisp/panel_host_windows.go:625` `newPanelControl(m.windowOptions(), m.dataPath)` |
| 装进 fork 的 Chromium | `cmd/wisp/panel_control_windows.go:355` `browser.DataPath = dataPath`、`:433` `browser.Embed(hwnd)` |
| 落到 WebView2 API 的 `userDataFolder` | `forks/go-webview2/pkg/edge/chromium.go:95` `createCoreWebView2EnvironmentWithOptions(nil, windows.StringToUTF16Ptr(dataPath), 0, e.envCompleted)` |

⇒ 夹具夹名就是 **WebView2 的 user-data 夹**，没有中间改名、没有每进程前缀。

**"同夹"能造成的两种死法（本仓实测记录，票 408，逐字）**

- `internal/ball/orb_windows.go:432-439`：
  "it used to hand EVERY process the same folder (%LocalAppData%\wisp\orb), so `go test ./cmd/wisp` … self-collided on one
  user-data folder: **a contended CreateEnvironment completed with S_OK and a NULL environment out-param, the panic died on
  the recover boundary, the STA thread was gone and ball.New waited forever (whole-package timeout, ticket 408)**."
- `internal/ball/orb_windows.go:439-445`：
  "WebView2's guidance for multiple instances is **one user-data folder per process** … decision record (ticket 408 r1):
  pid alone … still collides when Windows reuses a pid and hands the new process a dead predecessor's folder with its stale locks."
- `internal/ball/orb_windows.go:447-458`（第二种死法，形状是**挂**）：
  "WebView2 **SHARES one browser process across environments that name the same user-data folder** with the same options, and in the
  whole-package run **the previous attach's browser process is still tearing down when the next attach boots - the new environment
  then JOINS a dying process, and the wrapper's Embed waits for a controller-completed callback that never fires** (measured: one
  predecessor ball in the process already cost 7.9s vs 1.3s single; under whole-package load the wait never resolved and **ball.New
  hung 10-17 minutes**)."

两条死法与票面两种形状对得上：`rc=139`（native、无 panic 行、无 verdict）与 90 分钟挂起。**但关键一环是"同时"**：
争用发生在两个 environment **同时**创建、或新的一个**加入正在死的那一个**的瞬间。

**今天的守卫（改变了预期形状）**

- `forks/go-webview2/pkg/edge/chromium.go:188-197`：`EnvironmentCompleted` 遇负 HRESULT 或 `env == nil` ⇒ 具名日志
  `webview2: environment boot failed (hr=…, env nil: …)` ＋ `bootFailed` ＝1，**不再走成功链**（注释逐字点名"measured live:
  res=S_OK with env=0x0 when the WebView2 browser process dies at birth"）。
- 同文件 `:209-217`：controller 同形守；`Embed` 见标志返 `false`（`:108-110`），泵耗尽未 `inited` 也返 `false`（`:123-128`）。
- 上层：`cmd/wisp/panel_control_windows.go:433-436` `if !browser.Embed(hwnd) { … return nil, errors.New("panel host: WebView2 Embed failed …") }`
  ⇒ `bringUp` 返错 ⇒ **测试响亮变红**，不是整进程 139。
- 守卫落地：`da5d559b7`（**2026-11-08 20:13**）；三处 de-share：`68275983e`（**2026-11-11 19:22**）。

## 3. A/B 实测（**每发独占、全程串行**；⛔ 从未同时跑两个测试进程）

### 3.1 方法

- 目标四枚（票面死点 ＋ 同夹具的兄弟）：
  `TestAC14GoSideEvalPushReachesThePage`（`panel_resident_windows_test.go:1002`，死亡点 2 次）、
  `TestAC13BringUpSurvivesAReusedThreadQuit`（`:591`，死亡点 3 次）、
  `TestPanelThreadIsSTAAndExitsCleanly`（`:1034`，10-13 那族死亡点）、
  `TestAC13BringUpRefusesAThreadWithAQueuedClose`（`:667` 叶 `-second`）。
  四枚都真建窗：每发 `MARK299 after CreateWindowExW` ×4、`MARK299 after Embed` ×4（读数在 `logs/*.err.txt`）。
- 二进制三枚（**仓内零触碰**，A 组用 `-overlay` 换掉 `cmd/wisp/panel_resident_windows_test.go` 一个文件）：
  - `bin/wisp-A1-fixed-stale.test.exe`：overlay `overlay/panel_resident_A1_fixed_stale.go` = 把 `uniqueProfileLeaf` 的叶改回
    `filepath.Join(os.TempDir(), prefix)` 并**删掉三处 `t.Cleanup(os.RemoveAll)`** ⇒ 复刻**旧形状**（固定名、夹跨发累积）。
  - `bin/wisp-A2-fixed-fresh.test.exe`：overlay `panel_resident_A2_fixed_fresh.go` = 只把叶改成固定名、**保留**cleanup ⇒
    把"同名"与"旧内容"两个变量拆开。
  - `bin/wisp-b-unique.test.exe`：**HEAD 原样**（每发 pid+nanotime 唯一叶）。
- 每发：`Start-Process <test.exe> -test.run '^(四枚|…)$' -test.v -test.timeout 10m`，读 `ExitCode`（native 码按 `0x{0:X8}` 格式化），
  300s 硬上限（到点 kill 并记 `KILLED-300s`），退出后立刻抓 profile 夹快照与 `msedgewebview2` 命令行里的 `--user-data-dir`。
  每发开跑前 `Get-Process wisp.test` 计数 = **0**（表里 `comp=0` 即这条）。
  两个 harness：`run-one.ps1`（有间隙）、`run-tight.ps1`（**上一发退出即起下一发**，让前发的 `msedgewebview2` 还在退场时就起下一发）。

### 3.2 逐发读数

| 组 | 发 | rc | rcHex | 耗时 | verdicts | panic | bootFailed | Embed failed | 建窗 | 备注 |
|---|---|---|---|---|---|---|---|---|---|---|
| A1 固定/旧夹 | r1 | 0 | 0x00000000 | 4.80s | 4 PASS / 0 FAIL / 0 SKIP | 0 | 0 | 0 | 4 | 夹文件 175->178 |
| A1 固定/旧夹 | r2 | 0 | 0x00000000 | 4.53s | 4 PASS | 0 | 0 | 0 | 4 | 178->179（同一夹继续累积） |
| A1 固定/旧夹 | r3 | 0 | 0x00000000 | 4.60s | 4 PASS | 0 | 0 | 0 | 4 | 179->182 |
| A1 紧循环 | t1 | 0 | 0x00000000 | 6.13s | 4 PASS | 0 | 0 | 0 | 4 | 退出时浏览器 31->48 |
| A1 紧循环 | t2 | 0 | 0x00000000 | 8.22s | 4 PASS | 0 | 0 | 0 | 4 | **前发浏览器仍活（46 起）**，仍绿 |
| A1 紧循环 | t3 | 0 | 0x00000000 | 6.57s | 4 PASS | 0 | 0 | 0 | 4 | 46->47 |
| A2 固定名+清理 | r1 | 0 | 0x00000000 | 4.81s | 4 PASS | 0 | 0 | 0 | 4 | `RemoveAll` 部分成功：175->53 文件 |
| A2 固定名+清理 | r2 | 0 | 0x00000000 | 4.79s | 4 PASS | 0 | 0 | 0 | 4 | 53->20（半删夹仍能起） |
| B 每发独立 | r0(试跑) | 0 | 0x00000000 | 8.30s | 4 PASS | 0 | 0 | 0 | 4 | 该发脚本尚未数 stderr（建窗列后补） |
| B 每发独立 | r1 | 0 | 0x00000000 | 5.49s | 4 PASS | 0 | 0 | 0 | 4 | 新叶 `<prefix>-46884-1794626…` |
| B 每发独立 | r2 | 0 | 0x00000000 | 5.61s | 4 PASS | 0 | 0 | 0 | 4 | 唯一叶数 120->122 |
| B 每发独立 | r3 | 0 | 0x00000000 | 5.85s | 4 PASS | 0 | 0 | 0 | 4 | 122->124 |
| B 紧循环 | t1 | 0 | 0x00000000 | 6.20s | 4 PASS | 0 | 0 | 0 | 4 | 退出时浏览器 31->53 |
| B 紧循环 | t2 | 0 | 0x00000000 | 7.84s | 4 PASS | 0 | 0 | 0 | 4 | 53->76 |
| B 紧循环 | t3 | 0 | 0x00000000 | 8.72s | 4 PASS | 0 | 0 | 0 | 4 | 76->99 |

**15 发：rc 全 0、verdict 全 PASS、`panic` 0、`boot failed` 0、`WebView2 Embed failed` 0、`0xc0000005|access violation` 0。**
A/B 的耗时没有可判别的方向（A1 4.5-4.8s 偏快，是被复用的**热夹**；B 5.5-5.9s），紧循环两组的抬升同形（A1t 6.1-8.2s vs Bt 6.2-8.7s）。

### 3.3 "互踩同一个夹"被真制造出来的直接读数（`results-udf-share.jsonl`）

`udf-share.ps1` 的做法：先跑一发（warm），**它退出后立即**跑第二发并在第二发整个生命周期里每 300ms 采样
`msedgewebview2` 的命令行 `--user-data-dir` 与 `ParentProcessId`。一次只有一个测试进程。

- **A 组（固定名）`A1share2`**：warm pid `49456`、sample pid `344`；11 个采样 tick。
  `coexist` 逐字：`tick=0..3 udf=…\Temp\wisp-33r6-reused-thread-quit\EBWebView warmHosts=1 sampleHosts=1`，
  `tick=6 udf=…\Temp\wisp-33r7-queued-close\EBWebView warmHosts=1 sampleHosts=1`。
  判读：**同一个 user-data 夹上，同一 tick 里同时存在两个顶层浏览器 host** —— 一个的父进程是**已经退出的上一发测试进程**、
  一个是**当前发测试进程**；连续共存约 1.2s（4 个 tick）。这就是"两发互踩同一个 profile 目录"的实测形状。
  而这一发：`warmRc=0 / sampleRc=0`，warm 4.89s、sample **8.41s**（对照 A1 松跑的 4.5-4.8s）。
- **B 组（每发独立）`Bshare`**：`coexist=none`、`maxHostsOnOneUdf=1`、两发 rc=0（7.62s / 8.41s）。
- 诚实边界：`maxHostsOnOneUdf` 这个数会被"父浏览器已死、渲染进程变孤儿"抬高（A 组读到 4），所以本件只把
  **`parent == 测试进程 pid`** 的那一对当作"两个 host 并存"的硬证据，不拿 4 这个数说事。

## 4. 判定

1. **"固定名 ⇒ 两发互踩同一个 profile 目录"：坐实。** 代码链 `panel_host_windows.go:625 -> panel_control_windows.go:355 -> :433 -> chromium.go:95`，
   实测同夹并存两个 host（§3.3 的 `tick=0..3`）。这一半指控成立。
2. **"共享夹状态足以解释 139/挂起"：不充分。** 串行可及的最强形状（同夹两个浏览器实例并存约 1.2s、且夹被上一发写脏/半删）
   下 15 发全绿、零 native 异常。缺的正是 408 记录的那一瞬：**两个进程同时调 `CreateEnvironmentWithOptions` 争同一个夹**
   （`orb_windows.go:436` "a contended CreateEnvironment"），或**新 environment 加入一个正在死的浏览器进程**（`:447-458`）。
   本轮能造出"并存"，造不出"同时创建" —— 而那需要两发真并发，是本腿的边界外。
3. **时间线反证（对 11-13/14 那枚 `rc=139`）**：r3 守卫 11-08 20:13 已落、三处 de-share 11-11 19:22 已落，139 出现在 11-14。
   ⇒ 即便并发窗口里共享夹真的触发，今天的预期形状是**具名红**（`webview2: environment boot failed …` -> `panel host: WebView2 Embed failed`），
   不是整进程 139。要把 139 重新挂回这条链，需要"崩溃发生在守卫射程之外"（浏览器被带走后对已释放 COM 对象的 native AV）——
   本轮**零证据**，不主张。
4. **残余面（de-share 只动了三处，默认档还有建窗用例共享固定夹）**：
   - `cmd/wisp/panel_host_windows_test.go:609` -> `%TEMP%\wisp-33r1-panel-profile`，用例 `TestPanelHostRealWindowHopAndLifecycle`（`//go:build windows`，**真建窗**）。
   - `cmd/wisp/panel_assets_299_windows_test.go:161` -> `%TEMP%\wisp-299r1-panel-profile`，用例 `TestPanelPageFetchesItsOwnAssets299`（默认档，**真建窗**）。
   - `cmd/wisp/panel_holder_319_windows_test.go:53` -> `%TEMP%\wisp-319r1-holder-profile`：`:45-50` 逐字自己说"Nothing here creates a
     window (the headless half never asks bringUp), so the folder is a name and not a profile" ⇒ 不参与争用。
   - `cmd/wisp/panel_unbound_253r6_windows_test.go:73` -> `wisp-253r6-panel-profile`：本腿未判它建不建窗（未做）。
   - winlive 档若被打开跑，另有 `panel_host_windows_test.go:922`、`panel_host_windows_live_test.go:37`、
     `panel_host_windows_budget_winlive_test.go:57`、`panel_death_replay_319_winlive_test.go:75`、`panel_transport_live_35v2_windows_test.go:186`、
     `panel_reshow_255r8_winlive_test.go:125` 六处固定名。
   ⇒ **两发并发时，"互踩"这条通道没有被关死**；它仍是并发窗口里最像的放大器，但本轮读数说它**不是**独立充分因。

## 5. 边界与没做

- ⛔ 未跑整包（四枚定向，最长 8.72s）；⛔ 未做真并发（任何时刻只有一个测试进程；§3.3 的时间重叠来自**上一发已退出、它的浏览器还在退场**，不是两个测试进程同时在跑）。
- ⛔ 零产码/既有测试/阈值/golden/契约/`scripts/**` 改动。A 组是 `-overlay` 换文件，仓内那份 `cmd/wisp/panel_resident_windows_test.go` 一字未动
  （`git status --porcelain cmd/wisp/` 空读数）。
- 未判：`panel_unbound_253r6:73` 建不建窗；`Test219AC2Mutation` 那族嵌套 `go test` 对整包时长的贡献（票 314 别件已记）。
- 本件**不主张** 139 的定因，只判 AC#1 那句"共享夹状态是否充分" = **不充分**，并给出它唯一能致命的那个窗口。

## 6. 复跑（逐字）

```powershell
# 1) 二进制（A1/A2 用 overlay；B 用 HEAD 原样）
go test -c -o .scratch\wisp\probes\314-profile\bin\wisp-b-unique.test.exe ./cmd/wisp/
go test -c "-overlay=$((Resolve-Path .scratch\wisp\probes\314-profile\overlay\overlay-A1-fixed-stale.json).Path)" `
  -o .scratch\wisp\probes\314-profile\bin\wisp-A1-fixed-stale.test.exe ./cmd/wisp/
go test -c "-overlay=$((Resolve-Path .scratch\wisp\probes\314-profile\overlay\overlay-A2-fixed-fresh.json).Path)" `
  -o .scratch\wisp\probes\314-profile\bin\wisp-A2-fixed-fresh.test.exe ./cmd/wisp/
# 2) 逐发（一次一发）
pwsh -NoProfile -File .scratch\wisp\probes\314-profile\run-one.ps1  -Bin <bin> -Tag A1 -Run 1
pwsh -NoProfile -File .scratch\wisp\probes\314-profile\run-tight.ps1 -Bin <bin> -Tag A1t -Runs 3
# 3) 同夹并存取证
pwsh -NoProfile -File .scratch\wisp\probes\314-profile\udf-share.ps1 -Bin <bin> -Tag A1share
```

## 7. 读数根

| 件 | 内容 |
|---|---|
| `.scratch/wisp/probes/314-profile/results.jsonl` | 15 发里的 9 发（A1×3、A2×2、B×4）逐发 rc/耗时/verdict/夹快照/UDF 快照 |
| `.scratch/wisp/probes/314-profile/results-tight.jsonl` | 紧循环 6 发（A1t×3、Bt×3） |
| `.scratch/wisp/probes/314-profile/results-udf-share.jsonl` | 同夹并存取证（A1share×2、Bshare×1），含每 tick host 明细 |
| `.scratch/wisp/probes/314-profile/logs/*.out.txt` / `*.err.txt` | 每发 stdout/stderr 全量（`MARK299 …` 在 stderr） |
| `.scratch/wisp/probes/314-profile/overlay/*` | A 组 overlay 的两个源文件与两个 overlay JSON |
| `.scratch/wisp/probes/314-profile/bin/*` | 三枚测试二进制（各 62.5 MB）—— **只留在盘上、不入 commit**（避免 187 MB 进仓）；要复跑按 §6 现场重编 |
