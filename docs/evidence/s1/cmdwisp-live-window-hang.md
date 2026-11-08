# 票 408 r2 —— cmd/wisp 整包「真窗/热键活体测试逐个无限挂起」第二根因裁定与修复

- 时刻：2026-11-08 下午（钟点以本文件各读数的日志时间戳为准）
- 腿：Wisp 写腿（票 408 第二根因）
- 前序：e7f420a6c / c656c0ed8（r1，UDF 按进程唯一化 + 锁文件清扫 + 四牙）已落；另一腿在 r1 树上三发整包，`cmd/wisp` 仍未到终态，两个死点：`Test258BridgeRebindsLiveKeysFromConfigEdit` 挂 17m35s、`TestAC294MuteHotkeyRegistersLiveOnARealBallWindow` 挂 10 分钟。

## 1. 起手三核读数

- ①派单点名的串/行都在：`orbDataPath`/`acquireOrbDataDir`/`orbSweepStaleDirs` 全文在 `internal/ball/orb_windows.go`（r1 形状原样）；两枚死点测试全文在 `cmd/wisp/resident_hotkey_258_windows_test.go:266` 与 `cmd/wisp/resident_mute_294_windows_test.go:95`。
- ②引文出处：r1 的「S_OK + NULL env ⇒ panic ⇒ recover 吞 ⇒ STA 死 ⇒ ball.New 无限等」机制串逐字对上 `forks/go-webview2/pkg/edge/chromium.go`（`Embed` 的 `for { if inited… GetMessageW }` 无界泵，chromium.go:92-113；`EnvironmentCompleted` 只在 res<0 时 Fatal，S_OK+NULL env 直接 `env.vtbl.AddRef` nil 解引用，chromium.go:171-181）。
- ③飞件核：跑测前后 `tasklist` 无 wisp/go 测试进程；工作树 tracked 文件零改动（`git status --short | grep -v "^??"` 空读数）；跑完后桌面上的 48 枚 msedgewebview2 经命令行归因**全部是 Windows SearchHost 的**（`--webview-exe-name=SearchHost.exe`，开机起、父进程全活），**不是测试残留**——跑后零孤儿。

## 2. 复现实验（本腿实测）

| 实验 | 结果 |
|---|---|
| 单跑 `TestAC294MuteHotkeyRegistersLiveOnARealBallWindow` | **PASS 1.48s**（进程内第一枚球） |
| 单跑 `Test258BridgeRebindsLiveKeysFromConfigEdit` | **PASS 1.33s** |
| 三连 `TestAC246ChannelNeedsBothWindowAndExecutor` + 258Bridge + AC294 | 全 PASS，但**第二枚球的 `ball.New` 花 7.89s**（单跑 1.33s 的 6 倍） |
| 整包一发（r1 树，本腿复跑，`-timeout 40m`） | **393s 无声硬死**（无 Go panic、无 signal 行）：死于 `TestBookedSnapshotReachesHostBitForBit329`，RUN=345 / 出榜 277，两枚死点测试根本没跑到；该测试单跑 PASS 1.92s |

结论：两枚死点测试**单跑皆绿**——挂死不是测试件结构问题，是**同进程顺序多枚 WebView2 启动的语境退化**；失败形态在「无限挂起」与「无声硬死」之间随负载漂移，四发整包（他腿×3+本腿×1）没有一发到过终态。

## 3. 第二根因（本腿裁定）

**同一 UDF + 同参数 ⇒ WebView2 让新 environment JOIN 既有浏览器进程（官方多实例共享行为）。整包里前一枚球的浏览器树还在拆除时，下一枚 attach 在同一 per-pid UDF 上启动 ⇒ 新 env JOIN 一个濒死进程 ⇒ controller 完成回调永不触发 ⇒ `Embed` 的 `GetMessageW` 泵永不退出（chromium.go:92-113 无任何上限）⇒ `createOnSTA` 不返回 ⇒ `close(s.started)` 永不执行 ⇒ `ball.New`/`waitStarted` 无限等。**

支撑读数：

1. r1 只把 UDF 唯一到**进程**；`go test ./cmd/wisp` 整包是**一个进程**里顺序 ~7 枚 orb boot 共享同一枚 per-pid UDF——r1 消掉了跨进程撞夹，没消掉进程内顺序 boot 的 join 门。
2. 量测：进程内仅**一枚**前驱球就令下一枚 `ball.New` 从 1.33s → 7.89s（6×，join 濒死进程的等待成本）。
3. r1 留档的姊妹症状同门：同夹争用下 `CreateEnvironment` 返回 S_OK + NULL env（nil 解引用 panic），panic 死在 `close(s.started)` 之上（recover 边界在 observe 包装层）⇒ 同样无限等。两条门殊途同归：**boot 失败在产线等待点上没有失败路径**。
4. `internal/ball` 默认档**零枚**真 orb（唯一开 orb 的 `orb_live_windows_test.go` 是 `windows && winlive` 档）⇒ 「同进程顺序多枚 orb」这条路径从未被任何绿套件覆盖，`cmd/wisp` 整包是唯一跑它的地方。

## 4. 修复（两层，均为本腿写面）

### ① 产线（`internal/ball/orb_windows.go`，派单预授权「仅当必须调 UDF 策略」）

**r2 裁定：per-process 夹只当 ROOT，每枚 orb attach 在其下领独立叶 `<per-process>/boot-<n>`**（`orbBootDataDir`，进程内原子计数；site 夹随叶）。任何两枚 environment 从此**不可能**指到同一文件夹 ⇒ join 濒死进程的门物理消失。垃圾界形状不变：叶挂在持锁的 per-process 夹下，宿主死后被下一进程的首枚 attach 连夹带叶全扫。每叶首启建新 profile 的成本实测 ~1s（单跑整枚测试含全新 profile 共 1.48s），整包 ~7 叶 ≈ +7s，可忽略。

四牙（`orb_datapath_windows_test.go`）**零改动全绿**：tooth 1（进程内幂等）测的是 per-process 层，仍成立；跨进程 pid 叶、LOCALAPPDATA 回退、清扫三态全不受影响。fallback：叶建夹失败退回 per-process 夹（r1 形状）并 WARN，attach 不因 mkdir 失败而失败。

### ② 测试件（`cmd/wisp` 测试层，派单核心诉求「等不到要失败而不是挂死」）

新增 `cmd/wisp/resident_ball_boot_bounded_windows_test.go`：`startResidentBallBounded` = `startResidentBall` 的有界形（`select` + `time.After`，预算 2m ≈ 最坏健康 boot 7.9s 的 15 倍），到点 `t.Fatalf` 并**打全量 goroutine 栈**（挂死状态欠证据的那份诊断）。接线全部 **8 处**默认档 `startResidentBall` 测试调用点（approval_246×1、hotkey_258×3、hotkey_296×2、v1probe×1、mute_294×1）。断言零改动、skip 口径零改动；winlive 档不在整包射程未动。

已知代价（文件头具名）：到点失败**不**拆除已挂死的 boot——该 goroutine 连同其 STA 线程/窗/热键留在进程里，已红的 run 后续可能被污染出连带红；这是挂死本来的代价，使其可见而非掩盖。

## 5. 本腿范围外、如实上报的发现（零写入）

1. **存量红两枚**：`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` / `TestTicket255RosterStillMatchesTheActualReadSites`——`[ball]` 节实际读点 `cmd/wisp/ball_style_65_windows.go`+`config_cmd_65.go`，roster（`config_readers_255.go`）仍空。票 65 落地后的 roster 欠账，归该票收尾腿。
2. **`go build ./...` 整树红**：`build/mut_rvw.go`（package main、符号未定义）——突变审查存档拷贝，存量非本腿；按域 `go build ./internal/... ./cmd/...` rc0。
3. **整包负载语境的无声硬死**：`TestBookedSnapshotReachesHostBitForBit329`（面板区，单跑 PASS 1.92s）——无 panic 的进程级死亡，形态属同一退化类；面板产线写面不在本腿，上报不修。
4. **面板侧 WebView2 controller 从不 Close**（`panel_control_windows.go` WM_DESTROY 只 `ctlDropWindow`+PostQuitMessage，对照球侧票 65 的 controller.Close 教训）：健康路径靠父窗销毁让浏览器树异步退出（winlive 实测 ≤2s），负载下退出延迟即 33-r5 量到的 1/6 整包红——面板区退化风险的结构性源头，建议归面板收尾腿裁决（不属本腿写面，零改动）。

## 6. 验收读数（整包 `-timeout 40m -count=1`，r2 树，2026-11-08 18:30 发）

**到终态：931.402s（15.5 分钟）跑完，全程零无限挂起。RUN=708 / PASS=497 / FAIL=54 / SKIP=2。**

红名册分类（54）：

- **存量红 6**（与 r2 前那发 393s 硬死跑出的 6 枚逐一同名）：255Roster×2（`config_readers_255.go` 欠账）、198FirstRunCaller、313WritableUnion、AC4FocusReturn（对焦测量，桌面态敏感）、LockedSuffix。
- **两枚死点 = 有界响红 120.01s ×2**（`Test258BridgeRebindsLiveKeysFromConfigEdit`、`TestAC294MuteHotkeyRegistersLiveOnARealBallWindow`）：watchdog 到点 `t.Fatalf` + 全量栈，挂死变红。现场栈在本跑抓到活体：`EnvironmentCompleted(res=S_OK, env=0x0)` **同步**回调 → `chromium.go:175` nil 解引用 panic（observe recover 接住）→ `close(s.started)` 永不执行 → `ball.New` 无限等（r1 姊妹门在 virgin per-boot 夹上复现 ⇒ NULL-env 快速失败与 UDF 争用无关，是 WebView2 loader 在浏览器进程出生即死时的行为）。全跑 panic-recovered 恰好 2 次，全在两枚死点上；**中间五枚球 boot（258Mutation/258Occupied/296×2/v1probe）全部 <1.2s 绿** ⇒ 楔死是间歇性的、非位置性的。
- **语音族 ~30（0.00-0.08s 极速红）**：`onnxruntime 1.17.1`（System32 先于 PATH 被搜到，遮挡 1.28）——本机环境类，与代码无关。
- **他腿落树欠账 2**：256Options（产线多喂了 "Answer" 字段，7≠6）、AC247HandingTheLevel（`prototypeVisuals` 装配里被打开）——均在 r2 前那发没跑到的区域。
- **重执行族 ~6（5-10s）**：AC1ResidentLeg×3、AC246DevLeg、12r2、16r1Premortem×2——子进程「未到事件循环」，单跑语境绿（同 329 的负载语境类）；与本次 diff 无机制关联（不经过 `startResidentBall`/orb 路径）。

残余（写面外，报编排者）：**fork 级 nil-env 守卫**（`forks/go-webview2/pkg/edge/chromium.go`：`EnvironmentCompleted` 对 `env==nil` 记失败 + `Embed` 泵加失败退出臂返回 false）——落了它，今天的两枚 120s 红会按设计退化成「D2D 回退、boot 绿」；以及可选的 `waitStarted` 有界化（产线语义，须机主/编排者裁）。

## 7. CI 可跑性判断

默认档在无桌面 CI 上安全：无桌面 ⇒ `CreateWindowExW` 失败 ⇒ `createOnSTA` 返回 err ⇒ `rb.b == nil` ⇒ 两枚测试走既有 SKIP-LOUD 口径，不进 WebView2 路径，无挂死类。挂死类只在「有桌面 + 同进程顺序多枚 WebView2 boot」的整包语境出现；r2 的 per-boot 叶把该语境的 join 门拔掉，测试层预算再把残余挂死变响红。
