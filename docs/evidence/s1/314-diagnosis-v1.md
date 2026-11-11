# 票 314 诊断 v1 — 只读诊断腿（glm·禁跑整包）

- **谁／何时**：314-diag-v1 只读诊断腿，2026-11-11 18:45:55 +08，树 HEAD `8cfde490e`。
- **射程**：只读＋本件；未改任何产码／票面／Status；未跑任何测试（只 `go vet`＋`go test -c` 编译）；未动 `.github/**`；未在仓库内建临时件（编译产物落 `D:/tmp/wisp/20261111/314diag/wisp.test.exe`，**从未被执行**）。
- **门禁读数**：`go vet ./cmd/wisp` rc=0；`go test -c -o D:/tmp/wisp/20261111/314diag/wisp.test.exe ./cmd/wisp` rc=0（62,455,223 字节）。

## 三核（读到什么写什么）

1. **票 314 全文现读**（`.scratch/wisp/issues/314-cmdwisp-suite-process-death-in-window-zone.md`，末笔 2026-11-08 408 分界立票 `5bbf23762`）：死法签名＝死点全落在 `CreateWindowExW` 返回到首几条 wndproc 之间、无 Go panic、无 `--- FAIL`、无 `exit status`，整进程消失；q1–q5 五发 4 枚死点测试里只落两枚（`TestAC13BringUpSurvivesAReusedThreadQuit`×3＋`TestAC14GoSideEvalPushReachesThePage`×2）；后补 294-r1b 5/5 全死把第四枚 `TestPanelThreadIsSTAAndExitsCleanly` 也圈进来；2026-10-15 CI 两发早死同签名。**原始日志根 `C:\Users\swq\...\orch311-q-20261013\logs\` 本刻已不在盘上**（实测 LOG ROOT MISSING）——取证面只剩票面转述。
2. **票 408 全文现读**（`408-...-done.md`）：根因＝orb 侧 `orbDataPath()` 固定共享 `%LocalAppData%\wisp\orb` → 多进程/多 boot 自挤 → `EnvironmentCompleted` `S_OK`+NULL env → panic/recover 吞/STA 死/无限等 → 整包超时。**408 自己逐字裁定 314「同族不同根」**，且 §残余共享态具名列出面板测试件固定名 profile 目录为下一嫌疑。
3. **修法落地核验（git log 现读，均已在当前树）**：
   - r1 `e7f420a6c`：`internal/ball/orb_windows.go` `orbDataPath()` 改 `<base>/wisp/orb/<pid>-<8hex>`＋锁文件清扫；牙 `internal/ball/orb_datapath_windows_test.go` 已入库（git ls-files 在）。
   - r2 `eec09a028`（＋编译修正 `c656c0ed8`）：orb 侧再降级为 per-boot 叶 `boot-<n>`；`startResidentBallBounded` 有界启动等待（internal/ball，不罩 cmd/wisp 面板路径）；牙 `internal/ball/sta_wait_started_bounded_windows_test.go` 已入库（408 验牙腿补提交）。
   - r3 `da5d559b7`：fork 级 nil-env 守卫 `forks/go-webview2/pkg/edge/chromium.go`——`EnvironmentCompleted`（:188，nil env → `bootFailed` :195）/`CreateCoreWebView2ControllerCompleted`（:209，:215）/`Embed` 泵断出（:108）→ Embed 返 false，替掉 log.Fatalf 杀进程。

## 死点三枚测试与 UDF 的覆盖关系（本诊断的核心事实）

- 三枚死点测试全部经真实 WebView2 bring-up：
  - `TestAC13BringUpSurvivesAReusedThreadQuit`（`cmd/wisp/panel_resident_windows_test.go:580`）UDF＝`os.TempDir()/wisp-33r6-reused-thread-quit`（:581）。
  - `TestAC14GoSideEvalPushReachesThePage`（:989）与 `TestPanelThreadIsSTAAndExitsCleanly`（:1021）都走 `startPanelForTest`（:59-67），UDF＝`os.TempDir()/wisp-33r5-panel-profile`（:61，**固定名**）。
- 面板 UDF 设定点：`NewPanelManager(disp, assets, dataPath)`（`panel_host_windows.go:331`）→ `bringUp` → `newPanelControl(m.windowOptions(), m.dataPath)`（:625）→ `browser.DataPath = dataPath`（`panel_control_windows.go:355`）→ `browser.Embed(hwnd)`（:433，fork 守卫罩住这条）。
- 同类固定名共享 UDF 还有：`wisp-33r1-panel-profile`（`panel_host_windows_test.go:620,919`＋`panel_host_windows_live_test.go:37`）、`wisp-299r1-panel-profile`、`wisp-319r2-live-profile`、`wisp-35v2-panel-profile`、`wisp-33r5-panel-profile`、`wisp-33r6/33r7`、`wisp-253r6`、`wisp-255r8` 等——**408 的修一枚都没盖到这些**（408 写面只有 `internal/ball`＋fork）。
- **面板路径拿到的 408 收益只有 r3 nil-env 守卫**（fork 共享代码，Embed 返 false → `newPanelControl` 返错 → bringUp 报具名错），以及由此把「NULL env panic → recover 吞」这一族从死面里摘掉。

## 判定

**一句话**：票 408 修掉了 314 的同族近亲（orb 侧共享 UDF 超时族）并给面板路径上了 nil-env 守卫，但 **314 的「整进程无 panic 死亡」没有被 408 解释或修掉**——死点三枚测试的 UDF 至今仍是固定名共享 profile（408 自己具名的下一嫌疑），且死法签名（CreateWindowExW 后、无 verdict、无输出、整进程消失）与 408 的两族死法（可 recover 的 panic／有界等待响红）形状都对不上。

**是否仍可能复现**：结构上可能（共享 UDF 碰撞类原样在树里；r3 守卫只接得住「回调拿到干净 HRESULT＋nil 指针」这一种形状，接不住 WebView2 loader/原生侧硬死——那正是 314 无 Go 输出的形状）。**今天是否真复现：量不到**——本腿禁跑整包，AC#0（直跑测试二进制取退出码）与 AC#1/#2 的现场读数仍待能占 `cmd/wisp` 色的腿；2026-11-08 408 收口腿在导出树三笔整包均跑到终态（931s/579s，rc=1 有红但**无早死**）是「408 后该签名未复发」的仅有的间接对照（导出树、非本票 prescribed 形，不能当清白证明）。

## 最小修／收口建议

1. **最小修（结构对齐 408 已裁定的形状）**：把 cmd/wisp 面板测试件的固定名 profile 目录按 408 的「锁文件＋按进程唯一」或至少 per-boot 叶形状唯一化。第一刀最小面＝只改 `startPanelForTest`（panel_resident_windows_test.go:61）一处：`wisp-33r5-panel-profile` → `wisp-33r5-panel-profile-<pid>-<8hex>`（或 t.TempDir()），死点三枚里两枚立刻脱共享；`wisp-33r1-panel-profile` 三处为第二刀。改动全在 `*_test.go` 夹具，零产码面。
2. **收口顺序**：先走 AC#0（直跑本腿已编好的 `wisp.test.exe` 取退出码——**须由获准跑窗口测试的腿执行**，本腿禁跑），拿到退出码定性死法（native 异常码 vs os.Exit vs timeout panic），再决定修是「夹具唯一化」还是「须机主真机静桌面窗口」（AC#1/#2 的 A/B 与安静桌面读数本来就要求真机窗区，headless/忙桌面下量不到就具名写量不到，⛔ 假装能自动修）。
3. ⛔ 不建议在定性前动 `bringUp`／fork 产码——r3 守卫已罩住可 Go 感知的那半；剩下那半（原生硬死）没有 Go 侧可打的补丁，只能靠复现退出码定位。

## 受 314 牵连的「整包终态 AC」票清单（一行）

- 票 294-r1b（prescribed-form 整包 5/5 全死，等 314 定/修后按 AC#5 重跑 2+2）；票 311（q1–q5 死亡出其射程、另一立票的来历）；票 310（死的后果＝临时目录残留，同因不同面）；票 329／220 的超时族已由 408 认领（不同根）；票 303 已 `-done` 不受牵连。另：任何以 `cmd/wisp` 四数为分母的终态 AC 都受 `A976` 口径约束（先证跑到 `TestZZ315A1*` 终点锚，否则四数不作分母）。

## AC 逐条现状（票 314 AC#0–#4）

- **AC#0 死法取证**：未做（须真跑测试二进制；本腿只编译好弹药 `D:/tmp/wisp/20261111/314diag/wisp.test.exe`，rc=0，未运行）。仍待认领。
- **AC#1 最小重现＋归因 (a)（共享夹状态）**：未做；q5 的读数已把「共享 user-data 夹」从 orb 侧主嫌摘掉，但**面板侧固定名 profile 本刻仍是活的共享面**（本诊断核心事实），A/B 实验仍有效且更值得跑。
- **AC#2 归因 (b)/(d)（安静桌面／main 对照）**：未做；须机主静桌面窗口，⛔ 本腿自约，量不到。
- **AC#3 CI 暴露与推送前门**：2026-10-15 已实证 CI 会看见早死（两发 `38690995413`/`38687919508`）；推送前按本票形状重跑的门要求仍在，⛔ 未动 `.github/**`。
- **AC#4 处置**：本诊断即其前置——判定＝408 修了同族不同根，314 未被修掉；最小修建议见上节；处置（谁修/登记多长）待编排者拍。
