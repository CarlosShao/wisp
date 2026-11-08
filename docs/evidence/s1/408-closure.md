# 证据：票 408 收口 — 三笔合验的终态读数与红名册作差（复核腿，零产码）

- 时刻：2026-11-08 20:15–20:50 +0800；腿＝408 收口复核腿；只 commit 不 push。
- 起点：主树 HEAD `b7a2b6721`（含三笔 eec09a028 / da5d559b7 / 062e5fe70）。
- 方法：★`git -c core.autocrlf=false clone --shared` 导出树 `/d/tmp/wisp/20261108/408-closure`
  （检出 `b7a2b6721`；盘面 3 枚 embed 资产曾带 CRLF，已就地归 LF，`git diff` 空读数；从主树拷入
  untracked DLL：`build/`×3 ＋ `third_party/sherpa-onnx/`×3）。所有门禁与测试都在导出树跑，
  主树 tracked 面零写入。

## 1. 起手三核读数（现跑，非转述）

1. 票面 `.scratch/wisp/issues/408-cmdwisp-whole-package-times-out-shared-webview2-udf.md` 全文（56 行，
   含 AC#1–#4 与 Progress log 3 条）现读；AC 框与 Status 本腿零触碰。
2. 证据件现读：`docs/evidence/s1/cmdwisp-live-window-hang.md`（r2，72 行）、
   `docs/evidence/s1/webview2-nil-env-guard.md`（r3，94 行）。
3. 三笔 diff 现读（`git show --stat`）：见 §2 名册；飞件核＝主树 tracked 面零改动
   （`git status --porcelain | grep -v "^??"` 空），唯一相关 untracked 件见 §5 缺口①。

## 2. 三笔名册（`git show --name-only`）

| 笔 | 内容 | 牙 | 牙读数（本腿导出树现跑） |
|---|---|---|---|
| `eec09a028`（r2） | `internal/ball/orb_windows.go`（orb per-boot `boot-<n>` 叶）＋`cmd/wisp/resident_ball_boot_bounded_windows_test.go`（`startResidentBallBounded` 2m 有界）＋5 处测试接线 | 新测试件自身＋`orb_datapath_windows_test.go` 四牙零改动全绿 | `go test ./internal/ball`＝ok 0.482s |
| `da5d559b7`（r3） | `forks/go-webview2/pkg/edge/chromium.go`（nil-env/负 HRESULT 守卫＋`bootFailed` 标志＋Embed 泵失败返 false）＋`internal/ball/orb_windows.go`（接住半＋`DebugOrbUserDataDir()`）＋`orb_jarvis_live_windows_test.go:380` 一行 | `chromium_bootguard_windows_test.go` 五枚注入牙 | `go test ./pkg/edge -run "TestEnvironmentCompleted\|TestControllerCompleted\|TestBootGuard"`＝ok 0.068s（5/5） |
| `062e5fe70`（r3b） | `internal/ball/sta_windows.go`（`waitStarted` 裸 `<-s.started` → `select`＋`time.NewTimer(bootWaitBudget=30s)` 到点具名错误） | ⚠**牙 untracked 未入库**（`internal/ball/sta_wait_started_bounded_windows_test.go`，主树 `??`）——本腿拷入导出树现跑：**PASS 0.20s**；突变验证＝把 `waitStarted` 复原成裸 receive ⇒ 该牙 **FAIL 5.045s**（牙自己的 5s 护栏触发，不挂死），突变后 `git checkout` 复原 | 牙绿＋突变红（但未入库，见 §5①） |

## 3. ★整包终态（AC#1 读数；导出树 `b7a2b6721`，2026-11-08 20:36:41 发，45m 预算）

**到终态：579.515s（9.7 分钟）跑完，零无限挂起、零真 panic、零超时。rc=1。
RUN=707 / PASS=497 / FAIL=51 / SKIP=4。**

基线（`eec09a028` r2 那发）：931.402s，RUN=708 / PASS=497 / FAIL=54 / SKIP=2。

### 作差

| 项 | 基线 r2 | 本次 | 差 |
|---|---|---|---|
| 时长 | 931.402s | 579.515s | **−351.9s**（主体＝两枚 120s 有界红消失） |
| RUN | 708 | 707 | −1（他腿落树 roster 变化，非本三笔） |
| PASS | 497 | 497 | 0 |
| FAIL | 54 | 51 | **−3**（死点 2 绿＋roster 净变 1） |
| SKIP | 2 | 4 | +2（环境 skip，非本三笔射程） |

### 逐名红名册（51，按家族；基线对照）

- **两枚死点：红→绿（本收口的决定性读数）**
  `Test258BridgeRebindsLiveKeysFromConfigEdit` **PASS 0.06s**、
  `TestAC294MuteHotkeyRegistersLiveOnARealBallWindow` **PASS 0.65s**。
  活体日志（20:43:48.447）：`webview2: environment boot failed (hr=80010106, env nil: true)` →
  `ball: webview orb not attached; D2D fallback stays active (err="... boot guard: nil environment")`
  ——nil-env 守卫在整包语境**真实触发一次**并按设计流进 D2D 降级，等待方拿到具名错误，不再无限等。
- **存量红 6（与基线逐名同名）**：TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim（0.01s）、
  TestTicket255RosterStillMatchesTheActualReadSites（1.24s）、TestTicket198FirstRunCallerIsTheRunEntryOnly（0.25s）、
  Test313WritableUnionEqualsTheTwoRosters（0.00s）、TestAC4FocusReturnToPriorWindowGap33r5（4.69s）、
  TestLockedSuffixNeverTakesTheLockItself33r11（0.03s）。
- **他腿落树欠账 2（与基线同名）**：TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims（0.00s）、
  TestAC247HandingTheLevelToTheSeamIsNotVisibility（0.00s）。
- **重执行/负载语境族 4（基线 6）**：TestAC1ResidentLegInstallsItsLogListenerOnDisk（6.01s）、
  TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink（6.51s）、TestAC246DevLegIgnoresTheTestTaskInjection（7.96s）、
  Test12r2SettleWalkFiveSamplesEndToEnd（1.54s）。（基线的 16r1Premortem×2 本跑 0.01s 极速红，归入下面语音族。）
- **语音/onnx 族 39（基线 ~38）**：15e×5、16r1×2、26P×5、26R3×2、381r2×3、382WR2×3、390r1×3、392×3、
  394R1×2、395r1×1、41a2r1×4、WireR1×6——全部 ≤0.11s 极速红，`onnxruntime 1.17.1`（System32 先于 PATH）
  遮蔽类，**另一腿在治**；+1 来自 eec09a028 之后他腿落树的新语音票测试（381r2/382/390/392/394/395/41a2/WireR1 波次），非本三笔回退。

## 4. winlive 终态（AC#2 读数；导出树，同尺 `-run 'TestBallLiveOrb' -count=1 -timeout 30m`）

| 发 | 时刻 | rc | 读数 | 判定 |
|---|---|---|---|---|
| 1 | 20:23:50 | 1 | 9 FAIL 全 0.00s | **作废**：撞别腿在主树在飞的 winlive（`ball.test.exe` pid 4936，其 go 父进程 20:21:26 起自主树）；测试的 foreign-window 护栏拒量 |
| 2 | 20:27:22 | 1 | 53s 无声进程死亡，无 panic/无孤儿/无 WER/无事件日志 | 环境 GPU 抖动族（r3 证据 §6 已立档同族），非代码红 |
| 3 | 20:29:48 | **0** | **10/10 全 PASS，244.689s** | ★终态读数 |

## 5. 门禁（AC#3；导出树，`/d/work/base/gopath/bin/gofumpt.exe`）

- gofumpt `-l internal/ball cmd/wisp forks/go-webview2/pkg/edge`：三笔改动件**全不在名单**＝空；
  名单上 6 枚（`ICoreWebView2AcceleratorKeyPressedEventHandler.go`／`chromium_386.go`／`comproc_go17.go`／
  `comproc_go18.go`／`corewebview2.go`／`guid.go`）为 fork 上游存量、三笔零触碰（如实具名，不粉饰）。
- `go vet ./internal/ball ./cmd/wisp`：rc0。`go build ./internal/... ./cmd/...`：rc0
  （整树 `go build ./...` 的 `build/mut_rvw.go` 存量红系 r2 已留档的突变存档拷贝，非本三笔）。
- `go test ./internal/ball -count=1`：ok 0.482s（不含 §5① 那枚 untracked 牙）；拷入后 ok 1.264s。
- `scripts/d22scan.sh`：rc0，clean。
- LF：导出树 3 枚 embed 资产曾带 CRLF（clone 态），已归 LF；本腿写面两件（证据件＋票面 log 行）LF。

## 6. ★收口判定

- **AC#1 ✓**：整包 579.515s 到终态，⛔ 超时未再发生；两枚死点绿。
- **AC#2 ✓**：winlive 10/10 PASS 244.689s（rc0）。
- **AC#3 ✓**：门禁五件全过（gofumpt 按改动件口径空；6 枚 fork 上游存量格式件如实记账）。
- **AC#4（条件项）**：整包未再超时 ⇒ 无新死点可报，条件未触发＝空集满足。

**结论：票 408 三笔合起来在终态验收下成立，AC 全覆盖，可以结**——唯一下挂点见 §7①。

## 7. 如实缺口（逐格点名）

1. **`062e5fe70` 的牙未入库**：`internal/ball/sta_wait_started_bounded_windows_test.go` 在主树是
   untracked（`??`），`062e5fe70` 只含产码 `sta_windows.go`。本腿已替它补齐读数（绿 0.20s＋突变红
   5.045s，见 §2），但**牙不在仓库里＝产码有牙未落档**。复核腿写面禁碰 `.go`，须编排者裁定补一枚
   commit 收编该牙（撤销口令：不收编则该牙保持 untracked，突变验证读数以本文件为凭）。
2. 语音族 39 枚红为 `onnxruntime` System32 遮蔽环境类，**另一腿在治**，不阻塞本票（票面 r2 log 已记账）。
3. SKIP 2→4、RUN 708→707 来自 eec09a028 之后他腿落树（145 wave4 等），与本三笔无关，逐名见 §3。
