# 票 314 独立验收审计 v1（review-and-audit 路由组 · 非实现者）

- **谁／何时**：314-audit-v1 验收审计腿，2026-11-11 19:0x–19:2x +08，树 HEAD `8cfde490e`（`git rev-parse --short HEAD` 现读）。
- **射程**：只读核对 ＋ 本件。未改任何 `.go`／票面／台账／Status；未跑 `go build`／`go test`；未 `git add/commit/push`；未删任何文件；未在仓库内建临时件。
- **审计对象**：`docs/evidence/s1/314-diagnosis-v1.md`（47 行，下称"诊断"）。
- **角色**：本腿非该诊断的实现者，独立复核；结论按 实／虚／未证实 三档记。

---

## ① 结论

**部分成立。**

- **实**：诊断的「事实核」成立——`cmd/wisp` 面板测试件确实有多处**固定名共享**的 WebView2 UDF（`os.TempDir()` 下的固定 leaf），三枚死点测试用的就是其中两枚；408 的三枚提交（`e7f420a6c`／`eec09a028`／`da5d559b7`）**一枚都没碰**这些 profile 命名。
- **虚／未证实**：诊断的「因果核」（314 的整进程无 panic 死亡**由**固定名共享剖面引起）**证据不足**——诊断自己写了"结构上可能""今天是否真复现：量不到"；而票 314 自己的 q5 读数（"夹先挪走、跑中重生成，照样死"）指向相反方向。
- **不完整**：诊断的「最小修」只盖到三枚死点里的**两枚**，漏掉**最高频**死点 `TestAC13BringUpSurvivesAReusedThreadQuit`（q1–q5 占 3/5；CI run `38690995413` 死点）自己的固定名 profile `:581`；把"改一处即可"写成了近似全称。

> 一句话：**诊断说的"有固定名共享 profile"是真的；"那就是 314 的死因"没有被证明；"改 `startPanelForTest` 一处就够"不成立。**

---

## ② 逐条核对表

| # | 诊断的 claim | 复核到哪（文件:行／命令） | 判定 |
|---|---|---|---|
| 1 | `startPanelForTest` 定义在 `panel_resident_windows_test.go:59-67`，UDF＝`os.TempDir()/wisp-33r5-panel-profile`（:61，固定名） | `cmd/wisp/panel_resident_windows_test.go:59-67`；`:61` 逐字 `filepath.Join(os.TempDir(), "wisp-33r5-panel-profile")`；`:63` `NewPanelManager(disp, builtinAssetsOrTestNil(t), dataPath)` | **实** |
| 2 | `wisp-33r5-panel-profile` 全仓命中 6 处，代码只有 1 处（`_test.go`），其余是文档转述 | `grep` 6 命中＝`cmd/wisp/panel_resident_windows_test.go:61` ＋ `docs/reports/HANDOVER.md:2735` ＋ `docs/evidence/s1/314-diagnosis-v1.md:20/22/33` ＋ `docs/evidence/s1/33-panel-host-c27-v2.md:91` | **实**（代码 1 处、在测试夹具、固定名） |
| 3 | `startPanelForTest` 是**三个**死亡点测试共用的**唯一**起 panel 入口（"改一处即可"单点） | 调用点 `:455 :781 :947 :990 :1022 :1073 :1125` ＋ `panel_death_replay_319_winlive_test.go:88`；三死点里只有 `:990`（AC14GoSideEvalPush）与 `:1022`（PanelThreadIsSTA）走它；`TestAC13BringUpSurvivesAReusedThreadQuit`（`:580`）**不走**，自己 `NewPanelManager`（`:583`） | **虚**（是 2/3 的入口，不是 3/3；且被 6+ 枚非死点测试共用，不是死点专属单点） |
| 4 | 三枚死点测试的 UDF 至今仍是固定名共享 profile | `:61`（33r5，2 枚用）＋ `:581`（33r6，AC13BringUp 用）＋ `:658`（33r7，AC13BringUpRefuses 用）；`%TEMP%` 顶层现量残留 `wisp-33r5-panel-profile`／`wisp-33r6-reused-thread-quit`／`wisp-33r7-queued-close`／`wisp-33r1-panel-profile`／`wisp-299r1-panel-profile` | **实** |
| 5 | 面板 UDF 设定点链 `NewPanelManager`(:331)→`newPanelControl`(:625)→`browser.DataPath`(:355)→`browser.Embed`(:433) | `panel_host_windows.go:331`／`:625`；`panel_control_windows.go:338`／`:355`／`:433` | **实**（行号逐枚对上） |
| 6 | 408 三枚提交没碰 profile 命名（写面只有 `internal/ball`＋fork＋几枚球/热键测试件） | `git show --stat`：`e7f420a6c`＝`internal/ball/orb_windows.go`＋`orb_datapath_windows_test.go`；`eec09a028`＝`internal/ball/orb_windows.go`＋6 枚 `cmd/wisp/resident_*_windows_test.go`（approval_246／ball_boot_bounded／hotkey_258／hotkey_296／hotkey_v1probe／mute_294）＋票面/证据；`da5d559b7`＝`forks/go-webview2/pkg/edge/chromium.go`＋其牙＋`internal/ball/orb_windows.go`＋`orb_jarvis_live_windows_test.go`＋票面/证据。**零枚 `panel_*_windows_test.go`** | **实** |
| 7 | "408 的修复不解决 314" | 同上；408 票面 `:4` 逐字裁定 314「同族不同根」，`:39` 把面板固定名 profile 列为**条件性**"下一个具名嫌疑"（原文"若整包仍超时，这是下一个具名嫌疑"） | **实**（在"408 没碰面板夹具"这个窄义上） |
| 8 | "314 的整进程无 panic 死亡**由**固定名共享 WebView2 剖面引起" | 无任何现跑读数支持；诊断自己写"结构上可能""今天是否真复现：量不到"；票 314 `:19` 逐字"q5 干掉'共享 user-data 夹状态'主嫌（夹先挪走、跑中重生成，照样死）" | **未证实／证据不足**（且与票面 q5 相抵触，见 ③.1） |
| 9 | "q5 的读数已把「共享 user-data 夹」从 **orb 侧**主嫌摘掉，面板侧固定名 profile 仍是活的共享面" | q5 日志根 `C:\Users\swq\AppData\Local\wisp-test-gotmp\orch311-q-20261013\logs` 现读 **EXISTS=False**（父根 `wisp-test-gotmp` True）⇒ 无法复核 q5 到底挪了哪枚夹；票 314 `:31`（立票时写的 AC#1）逐字指向 `:61`（面板 profile），不是 orb | **虚**（retrospective 重构；无盘上凭据支持"q5 只动了 orb 侧"） |
| 10 | 最小修＝只改 `startPanelForTest`（`:61`）一处，"死点三枚里两枚立刻脱共享" | 同 #3；`:581`（最高频死点）与 `:658` **不在**诊断的刀口清单里（诊断只补了 `wisp-33r1-panel-profile` 三处为第二刀） | **半实**（"两枚"这句诚实；但"最小修"清单不完整） |

---

## ③ 风险与遗漏

1. **因果核未被证明，且票面有一条反向读数。** 票 314 `:19` 逐字把"共享 user-data 夹状态"列为**已被 q5 排除**的主嫌。诊断用"q5 只动了 orb 侧"化解，但：(a) q5 日志根现读不存在，无法复核挪的是哪枚夹；(b) 票 314 `:31` 的 AC#1 在立票时就指向 `:61`（面板 profile），所以 q5 的"user-data 夹"更可能就是面板侧。**另一支同样读不通**：q5 的死点是 `TestAC13BringUpSurvivesAReusedThreadQuit`（用 `:581` 的 33r6），若当时只挪了 `:61` 的 33r5，则 q5 也**没有**证伪 33r6 那条。⇒ 两种读法都成立、都无盘上凭据 ⇒ 这条只能记「**未定性**」，⛔ 不能记"诊断已证明"。
2. **"最小修"漏掉最高频死点。** `TestAC13BringUpSurvivesAReusedThreadQuit` 在 q1–q5 占 3/5、在 CI run `38690995413` 是死点，它用 `:581` 的 `wisp-33r6-reused-thread-quit`（固定名、`%TEMP%` 顶层现量在）。只改 `:61` 对它**零效果**。诊断的刀口清单（`:61` 第一刀、`wisp-33r1-panel-profile` 三处第二刀）**没有 `:581`／`:658`**。
3. **死点集比"三枚"宽。** 票 314 `:52` 记 CI run `38687919508` 死在 `TestLegacySubShapeOneHookDiesInAReentryLoop`（`panel_transport_35r2` 族）之后；该族用 `panel_transport_live_35v2_windows_test.go:186` 的 `wisp-35v2-panel-profile`（固定名）。诊断只说"2026-10-15 CI 两发早死同签名"，没点名这第四枚死点，也没把它计入刀口。
4. **固定名 profile 面比诊断列的全。** 现量 `grep os.TempDir()` 在 `cmd/wisp` 命中 16 处，其中**固定名 profile 12 处（10 枚文件）**：
   `panel_resident_windows_test.go:61`（33r5）／`:581`（33r6）／`:658`（33r7）；
   `panel_host_windows_test.go:620,919` ＋ `panel_host_windows_live_test.go:37`（33r1）；
   `panel_assets_299_windows_test.go:161`（299r1）；`panel_death_replay_319_winlive_test.go:75`（319r2）；
   **`panel_holder_319_windows_test.go:53`（319r1-holder，诊断漏列）**；
   `panel_reshow_255r8_winlive_test.go:125`（255r8）；`panel_transport_live_35v2_windows_test.go:186`（35v2）；
   `panel_unbound_253r6_windows_test.go:73`（253r6）。
   （另 2 处 `dataroot_128_test.go:257,259` 是 `EvalSymlinks(os.TempDir())` 读数，不是 profile。）
5. **"408 后未复发"这条间接对照与诊断的假设相抵触，诊断没做 reconciliation。** 诊断 `:29` 自己记：408 收口腿在导出树三笔整包跑到终态（931s／579s，无早死）。若面板固定名 profile 是 314 的**充分**死因，408 那三笔（面板 profile 原样共享）也该死。诊断把它归为"导出树、非 prescribed 形，不能当清白证明"——合理但弱。**可 reconciliation 的一支诊断没写**：314 的死可能需要"orb 共享 UDF 把进程搞不稳"＋"面板固定名 profile"两因子叠加，408 拿掉 orb 那半后单因子不足致死。
6. **副作用（若采纳最小修）**：
   - `t.TempDir()` 形：每枚测试独立、收尾自动删（比 `<pid>-<8hex>` 干净）；但会**换根**——profile 从 `%TEMP%` 搬到 `GOTMPDIR` 根，而本仓已登记 `GOTMPDIR` 换根／换卷会让别的用例变色（票 310／311／317），且 `t.TempDir()` 在封存根上收尾可能删不掉（`docs/reports/orchestrator-handoff-20261015.md:361`）。
   - `<pid>-<8hex>` 形：不留清理（垃圾界要另设），且 408 的锁文件清扫形状不能直接搬（那是 `internal/ball` 产码面，测试夹具里搬＝另一套）。
   - 并发：`grep t.Parallel()` 在 `cmd/wisp` **0 命中** ⇒ 包内顺序执行，"并发互挤"不是包内机制；共享的真实面是**跨跑次／跨进程残留**（`%TEMP%` 顶层残留已在盘）。
   - 隔离性：唯一化只会改善、不会变差；真实代价是磁盘残留 ＋ 换根带来的读数漂移风险。
7. **诊断没漏的自相矛盾**：`:33` 的"死点三枚里两枚立刻脱共享"与 `:27` 的"死点三枚测试的 UDF 至今仍是固定名共享 profile"并列时，容易被读成"改一处即可"；建议收口时把"两枚"写进标题句。

---

## ④ 复核过的命令与读数

| 命令（只读） | 读数 |
|---|---|
| `git rev-parse --short HEAD` | `8cfde490e` |
| `git log --oneline -8` | 8 枚，HEAD＝`8cfde490e`（297 收口）；408 三枚在更早 |
| `git show --stat e7f420a6c` | 2 件：`internal/ball/orb_datapath_windows_test.go` ＋ `internal/ball/orb_windows.go`（+296/−5） |
| `git show --stat eec09a028` | 9 件：`internal/ball/orb_windows.go` ＋ 6 枚 `cmd/wisp/resident_*_windows_test.go`（球／热键／静音）＋新有界等待件 ＋ 票面／证据；**0 枚 `panel_*`** |
| `git show --stat da5d559b7` | 6 件：`forks/go-webview2/pkg/edge/chromium.go` ＋ 其牙 ＋ `internal/ball/orb_windows.go` ＋ `orb_jarvis_live_windows_test.go` ＋ 票面／证据；**0 枚 `panel_*`** |
| `grep wisp-33r5-panel-profile`（全仓） | 6 命中；代码 1 处（`cmd/wisp/panel_resident_windows_test.go:61`），文档 5 处 |
| `grep startPanelForTest`（全仓） | 16 命中；定义 `panel_resident_windows_test.go:59`；调用 8 处（本文件 7 ＋ `panel_death_replay_319_winlive_test.go:88`） |
| `grep '^func Test' panel_resident_windows_test.go` | 11 枚；三死点在 `:580`／`:989`／`:1021` |
| `grep os.TempDir()` in `cmd/wisp` | 16 命中，其中固定名 profile **12 处（10 枚文件）**，见 ③.4 |
| `grep t.Parallel()` in `cmd/wisp` | **0 命中** |
| `Test-Path C:\Users\swq\AppData\Local\wisp-test-gotmp\orch311-q-20261013\logs` | **False**（父根 `wisp-test-gotmp`＝True）⇒ 诊断"LOG ROOT MISSING"**复认** |
| `[System.IO.Path]::GetTempPath()` ＋ `Get-ChildItem -Directory -Filter 'wisp-*'` | `C:\Users\swq\AppData\Local\Temp\`；残留含 `wisp-33r5-panel-profile`／`wisp-33r6-reused-thread-quit`／`wisp-33r7-queued-close`／`wisp-33r1-panel-profile`／`wisp-299r1-panel-profile`（另有 60+ 枚 `wisp-test-<pid>`） |
| 读 `panel_resident_windows_test.go:50-124` | `startPanelForTest` `:59-67`；`:61` 固定名；`:63` `NewPanelManager(..., dataPath)` |
| 读 `panel_resident_windows_test.go:570-609` | `TestAC13BringUpSurvivesAReusedThreadQuit` `:580`；`:581` `wisp-33r6-reused-thread-quit`；`:583` 自建 manager |
| 读 `panel_resident_windows_test.go:978-1062` | AC14GoSideEvalPush `:989`→`:990` 走 `startPanelForTest`；PanelThreadIsSTA `:1021`→`:1022` 走 `startPanelForTest` |
| 读 `.scratch/wisp/issues/314-cmdwisp-suite-process-death-in-window-zone.md`（58 行） | `:19` q5 排除共享夹；`:31` AC#1 指向 `:61`；`:48` 294-r1b 5/5 全死；`:52` CI 第二发死点 `TestLegacySubShapeOneHookDiesInAReentryLoop` |
| 读 `.scratch/wisp/issues/408-cmdwisp-whole-package-times-out-shared-webview2-udf-done.md`（58 行） | `:4`"同族不同根"；`:39` 残余共享态具名面板固定名 profile 为"下一个嫌疑" |
| 读 `docs/evidence/s1/314-diagnosis-v1.md`（47 行） | 审计对象全文 |

---

## ⑤ 给编排者的一句

诊断的**事实核可用**（固定名共享 profile 面 ＋ 408 没碰＝实），但因果核**未定性**、最小修**漏最高频死点**。建议把处置拆成两问：㈠先跑 AC#0（直跑测试二进制取退出码）＋ q5 的**可复核重跑**（把 `%TEMP%\wisp-33r5-panel-profile` 与 `%TEMP%\wisp-33r6-reused-thread-quit` **两枚都**挪走、各自 ≥2 发，判"共享夹状态"是否充分）；㈡修夹具时把 `:61`／`:581`／`:658`（必要时全 12 处）一起唯一化，⛔ 只改 `:61`。
