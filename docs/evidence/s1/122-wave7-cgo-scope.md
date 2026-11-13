# 证据件：票 122 波 7 —— 开 cgo 后 pass2 的 24 条真 finding 清零（2026-11-13）

- 写腿：票 122 波 7（派单 `.scratch/wisp/probes/orch/20261113/txt/122-wave7-cgo-scope-dispatch.txt`）
- 口径（派单给定，本机自跑确认）：
  `GOOS=windows CGO_ENABLED=1 D:\work\base\gopath\bin\staticcheck.exe -tags=winlive -checks=all,-SA9009 ./cmd/wisp/ ./internal/audio/aecspike/`
- 本机环境：go1.27.1 windows/amd64；staticcheck 2026.2.2 (0.8.2)；mingw gcc 16.2.0（`E:\work\base\msys64\mingw64\bin` 前置 PATH）

## 读数

**前（改前，2026-11-13 19:0x，24 条逐字）**：

```
cmd\wisp\appearance_wave_windows_test.go:286:6: func floatPtrWave is unused (U1000)
cmd\wisp\panel_control_windows.go:691:24: method SetHtml should be SetHTML (ST1003)
cmd\wisp\panel_control_windows.go:719:6: type ctlRpcMessage should be ctlRPCMessage (ST1003)
cmd\wisp\panel_dispatch_binding_roster_253r1_windows_test.go:103:2: struct field initDoorId should be initDoorID (ST1003)
cmd\wisp\panel_host_windows_test.go:443:6: func pidSetIntersect is unused (U1000)
cmd\wisp\panel_host_windows_test.go:475:2: const mibTcpStateListen should be mibTCPStateListen (ST1003)
cmd\wisp\panel_pageover_33r10_windows_test.go:92:24: method SetHtml should be SetHTML (ST1003)
cmd\wisp\panel_reshow_255r1_windows_test.go:130:24: method SetHtml should be SetHTML (ST1003)
cmd\wisp\panel_transport_35r2_test.go:965:3: the surrounding loop is unconditionally terminated (SA4004)
cmd\wisp\quota_44r5_test.go:32:2: field wroteIt is unused (U1000)
cmd\wisp\resident_approval_live_246_windows_test.go:413:36: runtime.GOROOT has been deprecated since Go 1.24 ... (SA1019)
cmd\wisp\resident_ball_windows.go:674:25: func (*residentBall).currentWakePauseRead is unused (U1000)
cmd\wisp\resident_voice_41a2r1_windows_test.go:485:6: func voice41a2r1FileExists is unused (U1000)
cmd\wisp\resident_voice_bargein_392_windows_test.go:584:3: struct field SwingDb should be SwingDB (ST1003)
cmd\wisp\resident_voice_engines_windows.go:758:26: func (*voiceResidency).waitAcquire is unused (U1000)
cmd\wisp\resident_voice_windows.go:519:2: field warmDiscards is unused (U1000)
internal\audio\aecspike\longrun_windows_test.go:675:4: struct field targetDb should be targetDB (ST1003)
internal\audio\aecspike\longrun_windows_test.go:773:16: should convert e (type gateEpisodeReport) to GateEpisode instead of using struct literal (S1016)
internal\audio\aecspike\rig_windows_test.go:297:2: var subCpu0 should be subCPU0 (ST1003)
internal\audio\aecspike\rig_windows_test.go:300:2: var subCpu should be subCPU (ST1003)
internal\audio\aecspike\rig_windows_test.go:303:2: var vadCpu0 should be vadCPU0 (ST1003)
internal\audio\aecspike\rig_windows_test.go:306:2: var vadCpu should be vadCPU (ST1003)
internal\audio\aecspike\rig_windows_test.go:323:3: struct field SubtractCpuPercent should be SubtractCPUPercent (ST1003)
internal\audio\aecspike\rig_windows_test.go:324:3: struct field VadCpuPercent should be VadCPUPercent (ST1003)
```

全文存本机 `D:\tmp\wisp\20261113\122-wave7\baseline.txt`（仓外临时件）。

**后（改后，同口径复跑）**：`EXIT=0，FINDINGS=0` —— **cmd/wisp 与 internal/audio/aecspike 两包 0 条**。

## 24 条逐条表（符号 / 位置（改前行号）/ 档 / 处置 / 理由）

档义（派单四档）：① ABI 结构体字段 → `//lint:ignore`；② 故意保留的句柄/资源 → ignore；③ 真无主 → 删；④ 只在非默认 tag 有主 → 点名调用者＋tag。
另有一个派单预设的例外形状＝「名字被仓外契约钉死，判不可改」（红线 2 要求的说明见专判①）。三处 `SetHtml` 落在这一形。

### cmd/wisp（16 条）

| # | 符号 | 位置 | 档/处置 | 理由 |
|---|---|---|---|---|
| 1 | `(*panelControl).SetHtml` | panel_control_windows.go:691 | **不可改**＋`//lint:ignore ST1003` | 名字由依赖接口 `webview2.WebView` 钉死（`forks/go-webview2/common.go:62`）；`:329` 有 `var _ webview2.WebView = (*panelControl)(nil)` 显式断言，改名即编译失败。非反射/COM 绑定名（核查见专判①） |
| 2 | `(*docSink33r10).SetHtml` | panel_pageover_33r10_windows_test.go:92 | **不可改**＋ignore | 接口替身：`attachSink33r10` 把它赋给 `PanelManager.w webview2.WebView` 字段（panel_host_windows.go:195） |
| 3 | `(*geoSink255r1).SetHtml` | panel_reshow_255r1_windows_test.go:130 | **不可改**＋ignore | 同文件 `:72` 的 `var _ webview2.WebView = (*geoSink255r1)(nil)` 显式断言钉名 |
| 4 | `ctlRpcMessage` → `ctlRPCMessage` | panel_control_windows.go:719 | ST1003 改名 | 本地类型，3 处引用（:738/:783）全在本文件，同步改 |
| 5 | `initDoorId` → `initDoorID` | panel_dispatch_binding_roster_253r1_windows_test.go:103 | ST1003 改名 | 测试普查结构体字段，6 处引用全在本文件，同步改 |
| 6 | `mibTcpStateListen` → `mibTCPStateListen` | panel_host_windows_test.go:475 | ST1003 改名 | 有主：`:552` 真实调用；全仓仅本文件引用 |
| 7 | `SwingDb` → `SwingDB` | resident_voice_bargein_392_windows_test.go:584 | ST1003 改名 | 匿名 `settingReadings` 字段，`:609` 写 `:610` 读，真用；json tag `swing_db` 不变，序列化形状零变化 |
| 8 | `floatPtrWave` | appearance_wave_windows_test.go:286 | ③ 删 | 全仓 `-w` 尺 0 代码引用（`.scratch`/docs 命中为历史探针与旧证据件的历史读数，非代码引用） |
| 9 | `pidSetIntersect` | panel_host_windows_test.go:443 | ③ 删 | 全仓尺 0 代码引用；同文件 `pidSetSubtract`/`pidSetOfTree` 在用，交集助手无调用者 |
| 10 | `wroteIt`（stub44r5 字段） | quota_44r5_test.go:32 | ③ 删 | 全文件零读零写；非 ABI 字段 |
| 11 | `(*residentBall).currentWakePauseRead` | resident_ball_windows.go:674 | ③ 删 | 唯一字段读者 `mirrorTrayMute`（`:839`）**已持 `muteMux`** 直接读字段——调本 accessor 会自死锁，故字段有主、accessor 无主。孪生 `currentWakePauseToggle` 有测试调用者（resident_voice_15e_windows_test.go:619/:654），保留。注释同步改写并钉「为什么没有配对读回器」 |
| 12 | `voice41a2r1FileExists` | resident_voice_41a2r1_windows_test.go:485 | ③ 删 | 全仓尺 0 代码引用；注释自陈与生产 `voiceFileExists` 平行但从未被用 |
| 13 | `(*voiceResidency).waitAcquire` | resident_voice_engines_windows.go:758 | ③ 删 | 字段 `lastAcquire` 有真读者 `reading()`（`:913`，启动报告行），仅 accessor 无调用者 |
| 14 | `warmDiscards`（字段） | resident_voice_windows.go:519 | ③ 删 | 全仓 token 尺仅声明一条——连写入都没有；邻居计数器全部有读点，唯它无 |
| 15 | SA4004 | panel_transport_35r2_test.go:965 | **真 bug（残留笔误），改** | 见专判②；未机械改 |
| 16 | SA1019 `runtime.GOROOT` | resident_approval_live_246_windows_test.go:413 | **语义保留**：代码不动＋说理注释＋`//lint:ignore SA1019` | 见专判③ |

### internal/audio/aecspike（8 条）

| # | 符号 | 位置 | 档/处置 | 理由 |
|---|---|---|---|---|
| 17 | `targetDb` → `targetDB` | longrun_windows_test.go:675 | ST1003 改名 | 本函数局部匿名结构体字段，3 处引用同步 |
| 18 | S1016 struct literal → conversion | longrun_windows_test.go:773 | 采纳 `GateEpisode(e)` | 两结构体四字段（名/类型/序）逐项一致：`gateEpisodeReport`＝longrun_windows_test.go:526、`GateEpisode`＝longrun.go:275 ⇒ 真同型转换，非类型混用；json tag 不参与 Go 转换恒等 |
| 19-22 | `subCpu0`/`subCpu`/`vadCpu0`/`vadCpu` → `…CPU…` | rig_windows_test.go:297/:300/:303/:306 | ST1003 改名 | 局部变量，全部引用点在 `:297-:339` 本函数内一眼可数 |
| 23-24 | `SubtractCpuPercent`/`VadCpuPercent` → `…CPU…` | rig_windows_test.go:323/:324 | ST1003 改名 | 匿名报告结构体字段，`:330` 构造 `:339` 读取；json tag 不变 |

## 三条专判

### ① SetHtml —— 判「不可改」（红线 2 的答案）

**先做反射/COM 绑定核查**：`git grep -n 'SetHtml'` 全仓——`panel_control_windows.go` 虽然 import `reflect`（:71），但其 `reflect.TypeOf` 只摸 `callbinding` 的 handler 形参，与 `SetHtml` 无关；全仓无字符串形式的调用；COM vtbl 侧全在依赖库内部（forks/go-webview2 的非导出面），本仓不见 `SetHtml` 字样。⇒ **不是反射名，也不是 COM 绑定名**。

**但它同样不可改名**：三处都是依赖接口 `webview2.WebView`（`forks/go-webview2/common.go:26-84`，12 方法：Run/Terminate/Dispatch/Destroy/Window/SetTitle/SetSize/Navigate/**SetHtml**/Init/Eval/Bind）的满足方法：
- `panelControl`：`var _ webview2.WebView = (*panelControl)(nil)`（panel_control_windows.go:329）；
- `geoSink255r1`：`var _ webview2.WebView = (*geoSink255r1)(nil)`（panel_reshow_255r1_windows_test.go:72）；
- `docSink33r10`：`attachSink33r10` 把它赋给 `PanelManager.w webview2.WebView` 字段（panel_host_windows.go:195）。

改名 ⇒ 不再满足接口 ⇒ 三处当场编译失败。派单红线 2 防的是「名字由仓外契约钉死、agent 无权改」这一形——COM 绑定名只是它的一个实例，接口满足名完全同形 ⇒ 同判「不可改」，用 `//lint:ignore ST1003` 登记并写明理由。仓内既有先例（internal/ball/win32_windows.go 等 32 处 `//lint:ignore U1000` 全部带点名理由）说明本仓的 ignore 形状＝点名理由的登记，不是掩盖。
**撤销口令：删三处 ignore 指令行即回到 finding 态。**

### ② SA4004 —— 判「真 bug（残留笔误）」，修而非删

位置：`(*jsParser).equality()`（panel_transport_35r2_test.go:956-967 改前）。

- 外层 `for {` … `return left`（首循环无条件返回）＝死循环壳，即 SA4004 读数。
- 来源判定：`:940 or()` 与 `:969 additive()`（真二元链层）都是 `for { … continue }` 真循环；`equality()` 也是真二元层，本该同形。这个壳是**照抄 or()/additive() 骨架时循环体放错层**的残留：作者把「遍历四算符」塞进内层 `for _, op := range`（一遍扫描）后，忘了外层壳本该是真循环条件。
- **它不只是死代码**：内层一遍扫描在跨算符序时会**意外链式**——`a == b === c` 会把两个算符都挂到同一棵二元树（`==` 处理完后 `===` 继续命中），而真实 JS 语义与 or()/additive() 的形状都是「同一层只吃连续同算符」。`a == b != c` 这类输入当前台件里没有（本文件 137 处比较逐一看过，JS 字面量全是单比较），所以它是**潜伏缺陷**而非已触发缺陷。
- 修法＝剥掉死循环壳、保留内层一遍扫描，并把「为什么不做真循环（不建模 JS 链式比较）」钉进注释：语义上逐字节等价于旧码在壳剥除后的行为，同时把意外链式的容忍一并关掉（旧码的意外容忍恰是那个笔误的另一半）。反控：台件 JS 无链式比较 ⇒ 全部 35r2 用例的期望输出不变。
- 对照组：`or()`/`additive()` 两个真循环未动。

### ③ SA1019 `runtime.GOROOT` —— 判「兜底语义成立，保留＋ignore」，未机械替换

位置：`buildEscListener246`（resident_approval_live_246_windows_test.go:410-415 改前）：
`goExe = "go"` → `LookPath("go")` 成功用之；否则 `filepath.Join(runtime.GOROOT(), "bin", "go.exe")` 存在则用之。

- 派单给的建议面是「换 `go env GOROOT`」。**判定不换**：这个 else 分支只在 `LookPath("go")` 已经失败（go 不在 PATH）时才走到——此时再 `exec.Command("go", "env", "GOROOT")` 是循环依赖，多半同样失败。`go env GOROOT` 替换的不是等价物。
- 弃用理由「构建根在拷机后无意义」在本场景不成立：这是本机测试 rig，构建根就是正确的根。
- 实测探针（`D:\tmp\wisp\20261113\122-wave7\goroot-probe\main.go`，本机跑）：`runtime.GOROOT()`＝`D:\work\base\go`；**PATH 真正移除后** `go env GOROOT` 仍返回同值（Go 工具自身按可执行文件位置推导，不依赖 PATH）——所以两种实现在本机等价，但探针同时暴露了语义面差异：`go env` 依赖能找到 go 可执行文件（本例靠 Windows 对 `os.Executable` 族调用的解析侥幸成功），`runtime.GOROOT()` 不依赖任何外部进程。兜底职责（「主路死了之后的最后一条路」）与后者更匹配。
- 处置：代码零改动；补两段说理注释＋`//lint:ignore SA1019`（点名循环依赖理由）。
**撤销口令：删 ignore＋两段注释即回 finding 态；若要真换，须由编排者裁「接受 LookPath 失败时兜底失效」的语义损失。**

## 名字核查（红线 1 的尺）

改后 `git grep -w`（范围 `cmd` `internal`，全 word）：`floatPtrWave`/`pidSetIntersect`/`mibTcpStateListen`/`ctlRpcMessage`/`initDoorId`/`SwingDb`/`currentWakePauseRead`/`voice41a2r1FileExists`/`waitAcquire`/`warmDiscards`/`wroteIt`/`targetDb`/`subCpu0`/`vadCpu0`/`subCpu`/`vadCpu`/`SubtractCpuPercent`/`VadCpuPercent` ⇒ **全部 0 命中**。
`.scratch/wisp/probes/**` 与 `docs/evidence/s1/122-r1-*.md`、`122-staticcheck-batch*.md` 里的旧名命中是历史读数记录（当时盘面的存照），不属于代码引用，不改写历史件。
名册面：`scripts/portable-tests.sh` GUARD、票 171 名册、`tokens_table_test` 均 `git grep` 不到上述任何旧名（scripts/ 全目录尺 0 命中）。

## 门禁（逐项）

| 门禁 | 结果 |
|---|---|
| `go build ./cmd/wisp/`（CGO=1） | EXIT=0 |
| `go vet ./cmd/wisp/ ./internal/audio/aecspike/`（CGO=1） | EXIT=0 |
| staticcheck（口径见头，两包） | **24 条 → 0 条**（EXIT=0, FINDINGS=0） |
| `gofumpt -l cmd/wisp internal/audio/aecspike` | 干净（一处文件尾空行由本波改动引入，`-w` 当场修） |
| `sh scripts/d22scan.sh` | clean（EXIT=0；`d22scan: clean - no D22 ban violations`，480 产码文件在册） |
| `go test ./cmd/wisp/` 全量 | **未跑**（票 98 豁免，本机跑不了全量）——按派单走 staged 载体定向跑，结果见下节，⛔ 不当绿 |
| `internal/audio/aecspike` 的 go test | **未跑**（sherpa cgo 台件依赖；本波对它的改动＝6 处纯改名＋1 处同型转换，全部由 staticcheck/vet/gofumpt 三把尺覆盖编译与命名面） |

## staged 定向测试（`scripts/wisp-cli-staged-tests.sh`）

跑的用例集＝本波改动文件里不需要真窗/语音模型的那批，**结果：全 PASS，exit 0**（2026-11-13 19:22，staged binary 从 cmd/wisp 目录带 pinned sherpa DLL 运行；log 在 `.scratch/wisp/stage-cli/wisp-cli-tests.log`）：

- wave 6 枚（含跨文件的 `TestWaveAssemblyRootsCarryTheReaders`，`-test.run TestWave` 模式多匹配；删 `floatPtrWave` 后同文件 5 枚全绿）
- 253r1 名册 2 枚（`initDoorID` 改名后普查与 fixture 正反控全绿：census 现读 installPanelTransport@panel_host_windows.go:1171 的 1 Bind＋1 Init、lockstep 判语照出）
- AC3 socket 尺 1 枚（改名后 `mibTCPStateListen` 在真监听正控里 0→1→0 三读数成立；`pidSetIntersect` 删除无影响）
- AC13 pageover 1 枚（最后 act＝nav 到 entry URL，1044 字节 byte-equal）
- 255r1 reshow 6 枚、35r2 transport 9 枚（**SA4004 修改后的 parser 全部用例期望输出不变**——反控坐实「台件无链式比较」）
- 44R5 quota 6 枚（删 `wroteIt` 字段后 stub 行为零变化）
- 392 JSON sentinel 1 枚（`SwingDB` 改名后读数照出）

合计顶层 32 枚 PASS / 0 FAIL / exit 0。
`go test ./cmd/wisp/` 全量与 aecspike 的 go test 未跑（豁免与理由见门禁表）。

## 改动面（本波写入的仓文件）

```
cmd/wisp/appearance_wave_windows_test.go                  # 删 floatPtrWave
cmd/wisp/panel_control_windows.go                         # SetHtml ignore + ctlRpcMessage→ctlRPCMessage
cmd/wisp/panel_dispatch_binding_roster_253r1_windows_test.go  # initDoorId→initDoorID（6 处）
cmd/wisp/panel_host_windows_test.go                       # 删 pidSetIntersect；mibTcpStateListen→mibTCPStateListen
cmd/wisp/panel_pageover_33r10_windows_test.go             # docSink33r10.SetHtml ignore
cmd/wisp/panel_reshow_255r1_windows_test.go               # geoSink255r1.SetHtml ignore
cmd/wisp/panel_transport_35r2_test.go                     # SA4004 修（剥死壳＋注释钉语义）
cmd/wisp/quota_44r5_test.go                               # 删 wroteIt 字段
cmd/wisp/resident_approval_live_246_windows_test.go       # SA1019 说理注释 + ignore（代码零改动）
cmd/wisp/resident_ball_windows.go                         # 删 currentWakePauseRead；注释改写
cmd/wisp/resident_voice_41a2r1_windows_test.go            # 删 voice41a2r1FileExists
cmd/wisp/resident_voice_bargein_392_windows_test.go       # SwingDb→SwingDB
cmd/wisp/resident_voice_engines_windows.go                # 删 waitAcquire
cmd/wisp/resident_voice_windows.go                        # 删 warmDiscards 字段
internal/audio/aecspike/longrun_windows_test.go           # targetDb→targetDB；S1016 转换
internal/audio/aecspike/rig_windows_test.go               # CPU 六处 ST1003 改名
docs/evidence/s1/122-wave7-cgo-scope.md                   # 本件
（commit 时显式 pathspec 逐件点名，⛔ push）
```

红线自检：未动 `build/mut_rvw.go`、未动 `.github/**`、未动阈值/golden/allowlist/契约文字；`git add` 用显式 pathspec；临时件只在 `D:\tmp\wisp\20261113\122-wave7\`。
