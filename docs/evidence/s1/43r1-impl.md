# 票 43 落地腿 `43-r1` —— Power & lifecycle events（写腿）

> **尺**：`HEAD=c73bde4b7`（`dev`），现跑 2026-10-31 12:36–13:46 +08。★读数逐条现跑，⛔ 转述。
> **写面（11 枚）**：新 7 枚（`cmd/wisp/lifecycle_windows.go`／`autostart_windows.go`／`crash_recovery.go`／`lifecycle_43_windows_test.go`／`autostart_43_windows_test.go`／`crash_recovery_43_test.go`／`lifecycle_live_43_windows_test.go`）＋改 3 枚（`cmd/wisp/resident_windows.go` 装配挂接／`cmd/wisp/run.go` 存储建立后一行 sweep／`cmd/wisp/config_readers_255.go` 名册按 233 §8 形状只加读者）＋本文件。
> **零字节**：`internal/session/**`（只调公开面：`NewScope/Acquire/Dispose/Stats/Stage`，一行未改）、`internal/risk/**`、`internal/speech/**`、`internal/proc/**`（`RunShutdownSequence`/`ShutdownOptions.Fast` 原样消费）、`internal/memory/**`（只调 `StartTaskLog/FinishTaskLog/ListTaskLogs/TaskLogByID/Open/Close`）、`internal/ball/**`、`internal/audio/**`、`internal/winsec/**`、`docs/PLAN.md`（D43 表一字未动）、`testdata/**`、`frontend/**`、`design/**`、`scripts/**`、`.github/**`、票面、台账。
> **临时件（只建不删）**：`D:\tmp\wisp\20261031\43r1\{bin,logs,cases15e,export-mut,d22scan.exe}`。突变全在仓外 `export-mut` 导出树（`git archive HEAD`＋覆盖），工作树零突变窗口。

## 0. 一句话

常驻进程从此接得住 OS 的两条命：**WM_QUERYENDSESSION**（隐形顶层窗口＋自泵线程 → 真子进程实测 exit **90.8ms**，≤5s 红线余量 55 倍；D38(e) 10 步 trail 全程落盘，4/5/9 照走、只有第 7 步非关键冲刷被跳）与**挂起/恢复**（PowerRegisterSuspendResumeNotification＋WM_POWERBROADCAST 双门归一漏斗；挂起前 L2 待确认全部作废＋计数告知「错过 N 个待确认操作」、采集设备停、SessionScope 在 voice owner 上按 28 的公开面 `Dispose` 释放）；**kill -9 后下一次启动**把开着 的 task_log 行标成 `interrupted`（panel/CLI 列表读同一数据源 ⇒「中断」可见）且结构上零重跑；**autostart** 按 TierRestart 语义在 boot 时把 `HKCU\...\Run` 与 `[app].autostart`（default false）对齐。**PowerRegisterSuspendResumeNotification 实际住在 `powrprof.dll` 而非头文件宣称的 `user32.dll`（本机导出表现勘）**——猜错即 boot 线程 panic，已修并留 Find() 防呆。

## 1. ★★★ 票面 5 格 AC 逐条（原文＋三分类判定＋读数）

| # | 票面原文（逐字） | 三分类 | 读数 |
|---|---|---|---|
| 1 | `Real suspend/resume (psshutdown/manual on runner) → session disposed pre-sleep (no models resident during sleep — sampler), audio restored post-wake, one resume-e2e task works.` | **替代形（如实具名）**：本机做不了真睡眠——真睡眠需无人值守唤醒装置（wake timer＋BIOS 支持），且会话本身会随宿主一起睡死、无人唤醒 ⇒ 测试会话结构上无法自证。**替代形＝两门真事件通道＋真序列**：①`PowerRegisterSuspendResumeNotification` 真注册（live 日志 `电源通知=已注册`，本机实测 powrprof.dll）②顶层窗口收 `WM_POWERBROADCAST`。「真断电」一跳由事件回调直驱/定向注入替代（同一生产体：`deliverPowerEvent`→漏斗→`runSuspendSequence`，296 教训的 hoisted seam）。sampler 面＝`session.Stats` 计数＋引擎 Close 计数，非 RSS 数（RSS 采样器在 `wisp slo`，本腿不重复造） | `Test43SuspendSequenceVoidsDisposesAndStopsCapture`（log 25）：挂起前 `EnginesHot=1`＋引擎 Close=0 → 挂起后 `EnginesHot=0`＋**fake 引擎 Close 恰 1 次**（释放一次、pre-sleep）；`SCOPE-DISPOSED reason=suspend:43-suspend cascade=models-unload,panel-hide,free-os-memory,grants-expire engines=1`（28 的级联原样跑）；采集 `stops=1`（挂起）→ `starts=2`（恢复重开，`runResumeSequence audio_reopened=true`）；**resume-e2e 进程内形状**：清场后下一次 `acquireTurn` 冷装载成功（loaderCalls=2、新 epoch ColdLoads=1/EnginesHot=1）——链完好非仅空。VAD 微型检测器不属 SessionScope、保持驻留（具名：本腿不虚报"全部模型离场"） |
| 2 | `Session-end: broadcast WM_QUERYENDSESSION → process exit ≤5s with steps 4/5/9 done (order audit log), Job children gone.` | **真机可做（等价注入）**：真 `HWND_BROADCAST` 会打进宿主桌面全部应用（外向动作、影响他人窗口）⇒ 用 **SendMessageTimeoutW 定向到子进程自己的 lifecycle 窗口**——消息泵/wndproc/fast 关停/D38(e) 全真，唯「消息从系统来」一跳是注入（真实关机时就是系统发同一条消息） | `TestLive43QueryEndSessionExitsWithin5s`（log 10，真子进程）：**exit in 90.8ms**（单调钟，≤5s ✓）；exit 0（走自身退出路径）；盘上 trail（`logs/*.jsonl`，log 10 的断言）`step 1 < step 4 (stop-audio): ok < step 5 (release-speech-sessions) < step 9 (close-job-object): ok < step 10 (exit)` 顺序在，`fast shutdown: skipping non-critical log/db flushes (D38e fast path)` 在；Job children：step 9 的 `JobScope.Close` 就是「OS 杀光子进程」本体（KILL_ON_JOB_CLOSE），trail 记 ok |
| 3 | `Suspend with pending L2 → voided + counted notice (D43 note updated with final decision).` | **真机可做**。合法行＝**D43 表 #24**（`AwaitingApproval | 队头被拒绝 → Acting（取消该调用）`）——suspend 时自动拒绝就是走这张已冻结行的拒绝漏斗；⛔ 未改 D43 表一字。**「D43 note 更新的最终决定」在此写明**（⇒ 改 PLAN.md 的 note 文字属契约文字改动，归编排者/机主批）：**§16.11#10 的倾向就此定案——挂起遇 `AwaitingApproval` ⇒ 一律作废＋按拒绝处理＋计数告知（「错过 N 个待确认操作」），恢复后不重放；最终措辞在票 54 review 定稿**。L1 窗口（无 reject 动词，SPEC-06 §2 B1）走与 D38(e) 第 3 步同形的 `Forget`＋审计行 | `Test43SuspendSequenceVoidsDisposesAndStopsCapture`：**真 queue 卡**（`approval.Gate.PendingApproval` 导出面直推＋注入 UI 记账，卡真在 C18 队列、ask 真停在 answer 通道）＋一张 L1 账上卡 → `voided=2`；parked ask 收回 `AnswerReject` 且理由含「挂起」；挂起后 `cards.Pending()=0`；console 行逐字 `错过 2 个待确认操作（系统挂起，已全部作废并按拒绝处理；恢复后不会重放）`；审计行 `SUSPEND-L2-VOIDED why="43-suspend" refused_l2=1 voided_l1=1` |
| 4 | `Crash-recovery fixture: kill -9 mid-task → next boot lists task as 中断; no rerun.` | **真机可做（真杀）**：`cmd/wisp/crash_recovery.go`——boot 写 `run-marker.json`（数据根），自身退出链清除；杀进程不跑 defer ⇒ marker 存活＝「上次没走自身退出序列」的唯一判据（不猜 pid、不算墙钟差）。sweep＝`ListTaskLogs` 中 state=running 的行经 `FinishTaskLog("interrupted")` 收口（196 词表＋365 阶梯的现成词；summary/cost 原值保真——COALESCE 语义）。**无重跑是结构性事实**：仓内没有任何 boot 期自动执行任务的路径，本腿未加第二个触发器 | `TestLive43KillNineMidTaskMarksInterruptedOnNextBoot`（log 10，真子进程＋真 `TerminateProcess`＝kill -9 同义）：强杀后 marker 在 ⇒ 第二次 boot stdout 出「已标记为中断」行；DB 复核 `state=interrupted`；`ListTaskLogs` 恰 1 行（**无新行＝零重跑**）；干净退出（CTRL_BREAK 走 D38(e)）后 marker 消失 ⇒ 之后 sweep 0/nil。headless 侧（log 25）：单测钉 marker 往返、running→interrupted 一次、二次 sweep 0（幂等）、无 marker 时 no-op、正常退出不误标 |
| 5 | `Autostart toggle writes/removes Run key; OFF default verified on fresh config.` | **真机可做（真注册表）**：`cmd/wisp/autostart_windows.go`——HKCU（无需提权）、值名参数化（生产 `Wisp`、测试用自己的名，机主真键零触碰）。**TierRestart 语义落地**：`app.autostart` 是 restart 档（233 的 TierRegistry），所以写/删发生在 **boot**（下一个进程开始的时刻），不在改配置瞬间——托盘/配置面板改键 → 233 的 restart-pending 告知 → 重启后 Run key 对齐 | `Test43AutostartRunKeyWriteRemoveRoundTrip`（log 25，真 HKCU）：写 → 注册表读回＝带引号 exe 路径（含空格路径防劈）→ 再写＝「已存在且一致」→ 删 → 读回不存在 → 再删＝「无本程序条目，无动作」（诚实无动作）；`Test43FreshConfigAutostartDefaultsOff`：`config.NewDefaults().App.Autostart==false`（schema 零值即 OFF default）＋fresh 目录 boot receipt `autostart=false…无动作`＋未写任何键 |

## 2. ★★ 关键实现事实（读数背后的三个「为什么」）

1. **两门归一漏斗**：`WM_QUERYENDSESSION` 只进顶层窗口（message-only 收不到广播 ⇒ 建隐形 WS_OVERLAPPED 窗口）；`PBT_APMSUSPEND` 自 Vista 起只发前台应用 ⇒ 补 `PowerRegisterSuspendResumeNotification`（**powrprof.dll**——本机 user32 导出表现勘，猜错=boot 线程 panic，附 `Find()` 防呆降级）。两门都只 post 进一个事件通道，由 `lifecycle-sta` 线程串行执行序列（Spawn 于 `observe.Default`，名册外名＝panel-sta 同款先例，33 的测试钉过「加名册=人工批准」）。
2. **窗口必须建在泵线程上**：首版把 createWindow 放在 boot 线程、pump 放在 Spawn 线程 ⇒ 消息队列属创建线程、无人泵 ⇒ SendMessageTimeout 全超时。修正后窗口在 `lifecycle-sta` 线程内创建＋自泵（panel-sta/ball ui-sta 同款 STA 形状）——live 读数 90.8ms 就是这个修法的直接产物。
3. **fast path 是旗标不是旁路**：`WM_QUERYENDSESSION` 的 wndproc 立即回 TRUE（系统在等回复），随后置 `lifecycleFastExit`＋`rt.RequestStop` ⇒ 事件循环正常返回 ⇒ **boot 自己的 defer 走 `rt.Shutdown(fast=true)`**——D38(e) 十步序列一个字节没动（`internal/proc` 零改动），快只体现在 `ShutdownOptions.Fast`（proc 自己的 D42#7 形状）。

## 3. ★★ 反向判据（三发，全在仓外 `export-mut` 导出树；工作树零接触）

| 发 | 变异（export 树内） | rc | 红句逐字 |
|---|---|---|---|
| M1 | `runSuspendSequence` 的 dispose 挂接摘除（`suspendVoiceEngines` 不再被调） | 1 | `lifecycle_43_windows_test.go:247: EnginesHot after suspend = 1, want 0 (no models resident during sleep)` |
| M2 | `applyAutostartValue` 的 `SetStringValue` 调用拔除 | 1 | `Test43AutostartRunKeyWriteRemoveRoundTrip` FAIL（写路径断言红） |
| M3 | `voidForSuspend` 的 L2 拒绝漏斗拔除（`cards.Reject` 不再被调） | 1 | `lifecycle_43_windows_test.go:224: voided = 1, want 2 (one L2 refused + one L1 window voided)` |
| **还原核**（时机=M3 之后） | ①仓内 `git status -- cmd internal`＝恰本腿写面（2 M＋7 ??），无他人文件；②仓内 `grep -rn MUTANT cmd internal`＝8 命中**全部**是 26r3/381r2 自己的既有失败消息字样，本腿文件 0；③工作树重编重跑 `Test43` **rc=0（log 25，9/9）**；④export-mut 树在 `D:\tmp\…\export-mut`（仓外，只建不删） | 0 | —— |

## 4. ★★ 零回归

| 对象 | 命令 | rc | 读数 |
|---|---|---|---|
| 15e/26R3/41a2r1/381r2 四族合跑 | staging 二进制，cwd=仓根（模型查找面）＋`WISP_15E_TMP/WISP_SILERO_MODEL/WISP_SPIKE_MODELS`（28-r1 §6 同套） | 0 | **19 PASS / 0 FAIL**（log 17）；`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` **逐位同值** |
| `internal/session` | `go test -count=1 ./internal/session` | 0 | `ok … 2.740s`（27 案全过，本腿只调公开面） |
| `internal/config` | `go test -count=1 ./internal/config` | 0 | `ok … 1.724s` |
| 233/223/255 定向 | `Test233HotTier…|Test233ReloadTier…|TestTicket223RunArms…|TestTicket223RestartTier…|TestTicket255Roster…` | 0（修后） | 前 4 枚直接 PASS；255 名册枚红＝本腿新读点未回名册 ⇒ 按 233 §8 形状**只加读者**重裁（`App: {cmd/wisp/autostart_windows.go}`＋`app.theme` 行文字同步「the section does have one production reader since 票 43…answers another key」），修后 `TestTicket255RosterStillMatchesTheActualReadSites` PASS（log 19） |
| `Test43` 原样 | staging `-test.count=2` | 0 | PASS×2 轮（log 23/24），-v 计数 9/9（log 25） |

## 5. ★ 四数＋两尺＋d22scan

- **写面四数**：`git status --porcelain -- cmd internal`＝`M resident_windows.go`、`M run.go`、`M config_readers_255.go`（修名册后入列）＋`?? 7 枚新件`；`git show --stat`（提交后）只含本腿件。
- **甲 `gofmt -l cmd`**：修前 `lifecycle_43_windows_test.go` 1 行（我写的）→ 修后**空**。已知既存假红（`cmd\wisp\models.go` 等 CRLF 件）不在本腿写面、未触碰。
- **乙 `gofumpt -l`**（绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`）：修前命中我 4 枚 → `-w` 修后对**全写面 10 枚单点跑＝空**。
- **d22scan**（`tools/d22scan` 独立 build，扫工作树最终态）：**rc=0 clean**——bans#1-5 internal/=280 cmd/=48、ban#8 internal/=715 cmd/=169（log 20）。**正控**：种一枚真 `≤`（U+2264）进 `internal/config/zz_d22_positive_control_43r1.go` ⇒ **rc=1 且逐字点名** `[emoji] ban #8 glyph in scope internal/ is banned (D23): non-comment text, string literals included`（log 21）⇒ 尺是活的；种后即删、复扫 rc=0（log 22）。
- **`go vet ./cmd/wisp`（默认＋winlive 两档）** rc=0。

## 6. 未做／具名残余（⛔ 藏）

1. **真断电一跳**：三分类判定在 §1 格 1——真睡眠需无人值守唤醒装置且会话随宿主睡死，替代形已如实标注；`psshutdown`/手动睡眠由机主在场时可做（验收腿若做，直接复用 live harness 的 FindWindow＋SendMessage 形状核对 trail）。
2. **托盘「开机自启」菜单项**：票面「toggle in tray menu + config」的托盘半边在 `internal/ball` 的冻结菜单表——非本腿写面（写面=cmd/wisp＋internal/{config,winsec?}）。**config 半边已全通**（任意配置写入方改键 → 重启后 Run key 对齐＋boot receipt 打印）；托盘项归下一张能写 `internal/ball` 的票。
3. **wall-derived schedules 重算**（票面 Key constraints）：本构建没有墙钟派生日程装置（提醒实现未落地），无可重算项——resume 审计行已具名此事实，非静默跳过。
4. **update-pending 消费端**：boot 检测 `update-pending.json` → 登记 `UPDATE-APPLY-HANDOFF` 审计事件（`markUpdatePendingHandoff`，单测在）；更新器本体是票 56，事件文件**不移动不删除**（56 的协议），56 落地后接这条线。
5. **resume 后状态机**：D43 冻结表无挂起边（293 的裁定先例：发明边=改契约）⇒ 状态机原地保持，审计行具名；「重进 Sleeping/Armed」由既有状态自然覆盖（挂起不产生状态迁移）。
6. **D43 note 文字更新**（PLAN.md）：属契约文字，本腿零触碰；最终决定的内容在 §1 格 3 写明，待编排者/机主批。
7. **`internal/panel` 邻接面**：本腿零字节（panel 相关文件一个没碰），233 热档邻接面无需复跑全量——255 名册枚已复跑绿。
