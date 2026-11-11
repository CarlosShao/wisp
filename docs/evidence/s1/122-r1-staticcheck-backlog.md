# 票 122 — staticcheck 清账腿 122-r1 裁决表（AC#1／#2／#4 底稿，AC#3／#5 顺带全交）

- 腿名：122-r1 · 日期：2026-11-11 · 写面：本腿（实现者）；验收方另派
- HEAD：`e31937a90`（A1709+A1710）；分支 `dev`
- **仪器版本（AC#7 格式要求，逐字）**：`staticcheck.exe 2026.2.2 (0.8.2)`、`gofumpt v0.12.0 (go1.27.1)`
  （本机 GOPATH bin：`D:\work\base\gopath\bin\`；staticcheck mtime 2026-11-01，与 CI 钉的 `@2026.2.2` 同版，读数可比）
- 跑法：主工作树 tracked `.go` 零未提交改动（仅 `.scratch` 未跟踪件与 1 枚 html 件），
  为躲开 gitignore 残渣 `build\mut_rvw.go`（编译错、非 HEAD）与 `cmd\wisp` 两枚未跟踪 323 测试件，
  本腿读数全部跑在 `git archive HEAD` 导出的**仓外纯净快照**（`%TEMP%\wisp122r1\snap`）：
  导出后 tracked `.go` 计数 1722 == `git ls-files '*.go'` 1722、文件集合 diff=0（逐名比对）。
  纯净性对 Go 读数等价于 HEAD；主工作树零写入（本腿只写本表与 `.scratch/wisp/probes/122/r1/`）。

---

## AC#1 — 开工重测（当前 HEAD 纯净快照，两 GOOS，root module）

| 读数 | GOOS=linux | GOOS=windows |
|---|---|---|
| 票面预做（audit-85-preflight，旧 HEAD；标题里的 34/78） | 34 | 78 |
| CI 自报分母（A101②，09-22） | — | 42 |
| 批 1 收口（122sa1，HEAD `6e2eb5845`，11-08） | 20 | 115 |
| 批 2 收口（122sb2，HEAD `31482f271`，11-08） | 14 | 98 |
| **今天 HEAD `e31937a90` 实测（本腿现跑）** | **57 条 finding** | **103 条 finding** |

原始仪器输出：`.scratch/wisp/probes/122/r1/ac1-linux-head.txt`、`ac1-windows-head.txt`（逐行，可复跑）。

今天分布（按 checker）：

- **windows 103**：`88×U1000 / 2×S1011 / 2×S1016 / 2×ST1005 / 2×ST1012 / 2×ST1018 /
  1×SA1019 / 1×SA4004 / 1×SA4010 / 1×SA5000 / 1×SA9009`
- **linux 57**：`50×U1000 / 2×S1011 / 2×ST1018 / 1×S1016 / 1×SA5000 / 1×SA9009`
  （另有 1 行依赖噪音：`sherpa-onnx-go-linux@v1.13.8` 构建约束排除——本机无 cgo linux 工具链时
  `go build` 同样不可达，属环境形状，非树的问题；linux 分母因此天然偏低）

**与批 2 收口（14/98）的差，逐块归因（都落到证据）**：

1. **SA4000（1 条）已绝迹**：批 1 表 (#1) 记录、本腿现核——commit `79041f071`（CI 修复腿 ci-r1，
   2026-11-01）把 `if !errors.As(err, &oe) || !errors.As(err, &oe) {` 改成
   `if !errors.As(err, &oe) {`（internal/observe/goroutine_test.go，现 :151）。今天两 GOOS 均无 SA4000。
2. **SA4006 今天 0 条**（批 1 时 3 条）：批 1 #6/#7/#8 全部修掉（cmd/wisp watchdog 两处补尺寸断言、
   slo_windows.go 初值改裸零值）。今天两 GOOS 均无 SA4006。
3. **SA1019 4→1**：批 1 修 3（cmd/wisp GOROOT 1 条＋internal/speech GetCurrentProcess 2 条），
   剩 1 条在 `internal/risk/proccpu_windows_test.go:20`（冻结，见 AC#2 表 #4）。
4. **linux 14→50 的主因＝分母修复，不是新债**：批 1/批 2 时 `internal/ball` 在 linux 有编译噪音
   （`orbdebug.go` undefined: Ball，批 1 表 AC#1 注记「linux 分母天然偏低」）；票 418 已把
   orbdebug 拆成平台中立＋windows-only 两件，本腿现跑 `GOOS=linux go build ./internal/ball/`＝**rc=0**。
   噪音消失后 linux 把整棵树真正分析完，暴露出 **44 条 linux-only U1000**（几乎全是
   `internal/ball/jarvis.go` 平台中立渲染数学＝只被 windows-only 代码引用的形状）。这批「新增」
   大部分是**一直存在、此前被编译噪音挡住没被数到**的面。
5. **windows 98→103 的净变化**：批 2 后新落代码带来的新 finding 6 枚（逐枚有 commit 归属，见 AC#4 表），
   减去批 2 后消失的若干枚（例：`problemsSnapshot` 批 1 表记「留-口」，今天 cmd/wisp 已有 79 处
   真调用——测试与产码接线把它养活了，不再报）。
6. **U1000 跨 GOOS 去重并集＝132 条**（windows 88 ∪ linux 50；两 GOOS 共同命中 6、
   linux 独有 44、windows 独有 82）。票面标题的 26 与批 1 的 96 都是旧账，**不并列覆盖、只注明过期**。

另两个 module（CI staticcheck 步同样遍历）：`tools/d22scan` windows 0 / linux 0；
`tools/mockllm` windows 0 / linux 0（全 rc=0）。

## AC#2 — 第一批三组五条（SA4000×1＋SA1019×3＋SA4006×1）逐条处置

**如实结论：本腿零行代码改动。五条里 2 条已在 HEAD 之前被前批修掉、3 条落在本腿无写面的
冻结/越界地界**。逐条：

| # | 形状 | 位置（票面原锚） | 今天实测 | 处置与归属 |
|---|---|---|---|---|
| 1 | SA4000 | internal/observe/goroutine_test.go:106 | **不存在（两 GOOS 均无）** | **已被 `79041f071`（ci-r1，2026-11-01）修掉，本腿零改动**。真缺陷证据见下节 |
| 2 | SA1019 | cmd/wisp/**（runtime.GOROOT，批 1 #2） | **不存在** | 已被批 1（122sa1）修掉：改读 `os.Getenv("GOROOT")`，断言不弱化 |
| 3 | SA1019 | internal/speech/cgo_boundary_windows_test.go:36/:95（批 1 #3/#4） | **不存在** | 已被批 1 修掉：改 `windows.CurrentProcess()`，删恒不触发的 err 分支 |
| 4 | SA1019 | internal/risk/proccpu_windows_test.go:20（批 1 #5） | **仍在**（windows 侧唯一 SA1019） | **不修，归口登记**：`internal/risk/**` 票面明令冻结，本腿零写面（⛔ 命中即跳过）。修法与 #3 完全同形（一行换 `windows.CurrentProcess()`＋删 err 分支）。取证腿已交：`docs/evidence/s1/122-sa1019-risk-verdict.md`（API＝x/sys v0.49.0 `windows.GetCurrentProcess` 逐字 Deprecated，error 恒 nil）。**归下一波/待编排者具名解冻**；解冻前 windows 分母恒含这条 |
| 5 | SA4006 | cmd/wisp/**（批 1 #6/#7/#8 共 3 条） | **不存在（两 GOOS 均无）** | 已被批 1 修掉（hold64/hold128 补尺寸断言＝加强；slo_windows.go 初值改裸零值＝行为逐位不变） |

### SA4000 为什么是真缺陷（票面要求的说明＋行为证据）

- 原形状（`git show 79041f071` diff 逐字）：
  `if !errors.As(err, &oe) || !errors.As(err, &oe) {` —— `||` 的两个操作数**逐字相同**。
- 错在哪：这枚用例（TestPanicInWorkerSurvivesAndCancelsRoot）要验的是 panic 错误的两个维度——
  类别维度（`ClassOfOrInternal(err) == ClassInternal`，:147）与类型维度（`errors.As(err, &oe)`，:151）。
  重复操作数让那行**看起来在验两件事、实际只验一件**（`A || A` 恒等于 `A`）；
  staticcheck 的 SA4000 报的正是这种「二元运算两侧同表达式」的失效断言。
- 为什么删是修正不是写松：票面给的两条路里这属于「本来就重复，删掉是修正」——不存在
  「第二个操作数本该是另一种错误类型」的读法（若本该不同，diff 里会出现第二个 `var`/target，
  实际没有；且 `errors.As` 对同一 target 两次求值是纯重复）。
- 行为证据（AC#2 要求的「仍红于该红的形状」）：现测试 `:151-153` 保留完整——
  `!errors.As(err, &oe)` 成立时 `t.Fatalf("panic error not an observe *Error: %v", err)`，
  断言没有变恒真；类别维度与类型维度仍是两个独立断言（:147 与 :151）。修后该包测试今天
  `-count=2` 全绿（见 AC#5）。

## AC#3 — 假阳性格：SA9009（只登记不修）

- 今天读数（两 GOOS 各 1 条，同一条）：`frontend\embed.go:4:1: ineffectual compiler directive due to
  extraneous space: "// go:embed below is the ONLY way those bytes reach wisp.exe. There is no" (SA9009)`
- 现场核过（本腿 read 逐字）：`:4` 是**散文注释**——票 77 AC#1 的文档句
  「the go:embed below is the ONLY way those bytes reach wisp.exe...」恰好以 `// go:embed` 开头；
  真正的指令在 `:19`，逐字 `//go:embed all:dist`，写法正确、有效。
- **为什么不修**：改 `:4`＝动票 77 AC#1 的验收文档；`frontend/**` 归机主外部 agent，本腿零写面；
  加 `//nolint` 类抑制＝拆门，票面明令禁止。⇒ **登记为已知假阳性，代码零改动**。
  这条在有人重排那句话的第一个词之前会一直在，CI lint 步会因它恒红——这不是本票能修的格。

## AC#4 — U1000 逐条归属表（只出表，本腿一条未删）

- 分母：**132 条**（两 GOOS 去重并集；windows 88 / linux 50 / 共同命中 6）。
- 判定图例：**留-ABI**＝Windows/COM 结构体·常量镜像，按 ABI 全量声明、Go 侧不读是设计内；
  **留-口**＝留给后续票的入口/诊断面/未接线产码；**留-试**＝测试夹具；
  **留-平**＝平台不对称（符号被 windows-only 代码消费，linux 构建里才闲置）；
  **删-c**＝建议下一波删（零引用已现核；**本腿零删**）；**冻**＝包冻结；
  **越界-cmd**＝落在 `cmd/**`（本腿写面外，⛔ 硬改，归下一波）。
- 判例注记（批 2 实证，本腿沿用）：这版 staticcheck 的 U1000 **不把 `_test.go` 引用计为使用**，
  「工具说 unused」不等于「删了不红」——删前必须全仓 grep 含测试。本表每条「删-c」都已做过这一步。

### A. 两 GOOS 共同命中（6 条）

| 位置 | 符号 | 判定 | 理由或归属 |
|---|---|---|---|
| internal/tools/shell_exec.go:68 | field out（waitOutcome） | 删-c | 票 163（`75f9479fa`，11-10）新落产码的死字段：捕获实际走 :116 的 `outBuf/errBuf`（guardedBytes），waitOutcome 里这两个 builder 全仓零写零读，唯一命中＝定义行 |
| internal/tools/shell_exec.go:69 | field errs（waitOutcome） | 删-c | 同上 |
| internal/tools/shell_session_test.go:31 | field closeErr | 删-c | 测试夹具死字段：全文件零引用（bridge.go:1058 的同名局部变量是另一物），唯一命中＝定义行 |
| internal/agent/reminder_timer_test.go:109 | (*fakeReminderClock).releasedCount | 留-试 | 假时钟夹具的计数面，随 agent 测试票核对（批 1 同判） |
| internal/memory/l1_warm_29_test.go:151 | l1WaitFor | 留-试 | 测试等待 helper（批 1 同判） |
| internal/tools/web_bridge_test.go:49 | (*scriptGate).levels | 留-试 | web 桥测试夹具（批 1 同判） |

### B. linux 独有（44 条）——三块平台不对称形状

**B1. internal/ball/jarvis.go（38 条，:17-:48、:55、:57、:59、:71、:82、:98、:141、:155、:164、:172、:176、:228）**
逐条列出（判定全部相同）：

| 行 | 符号 | 判定 | 理由或归属 |
|---|---|---|---|
| 17-20 | jarvisRingInner / jarvisRingMain / jarvisRingCore / jarvisRingField | 留-平 | jarvis 皮肤的平台中立渲染数学：整文件只被 `jarvis_windows.go`（产码 :32/:65/:83）与 `jarvis_windows_test.go`（:93-:131 等）消费——windows 面真在用，linux 构建里才闲置；删任何一枚＝windows 构建或 winlive 测试红 |
| 27-34 | jarvisRingInnerFrac / jarvisRingMainIdleFrac / jarvisRingMainThinkFrac / jarvisRingMainAlertFrac / jarvisRingCoreFrac / jarvisCoreDotFrac / jarvisRingTiltFrac / jarvisFieldMinFrac | 留-平 | 同上 |
| 35-48 | jarvisFieldMaxFrac / jarvisGlowRadiusFrac / jarvisCoreGlowFrac / jarvisDotRadiusPx / jarvisDotBrightPx / jarvisBrightShare / jarvisRingJitterMain / jarvisRingJitterInner / jarvisParticleCountBase / jarvisParticleCountMin / jarvisParticleCountMax / jarvisListenPulseFrac / jarvisScatterPerLevel / jarvisParticleSeed | 留-平 | 同上 |
| 55/57/59/71 | jarvisRng / newJarvisRng / (*jarvisRng).f64 / jarvisDot | 留-平 | 同上（jarvis_windows.go 的粒子场直接持有 []jarvisDot） |
| 82/98 | jarvisParticleCount / jarvisLayout | 留-平 | 同上（windows 产码 :32 的每帧路径） |
| 141/155/164/172/176 | jarvisMainRingFrac / jarvisRingSpin / jarvisListenPulse / jarvisScatterAmp / jarvisClamp01 | 留-平 | 同上 |
| 228 | jarvisFamilyColors | 留-平 | 同上（windows 产码 :65/:83 取粒子三色） |

**归属票：jarvis 渲染面（internal/ball）后续票**；本票与批 2 均无 ball 写面（派单 ⛔）。

**B2. internal/audio/device.go（3 条）**

| 行 | 符号 | 判定 | 理由或归属 |
|---|---|---|---|
| 60 | captureDeviceLister | 留-平 | `[audio] input_device` 解析面：唯一消费者 `wasapimic_windows.go:293`（windows-only）＋180k 测试 :45/:64/:76/:95；windows 真在用（票 180 腿 K 落的） |
| 66 | inputDeviceDefault | 留-平 | 同上（:79 的 "default" 短路值） |
| 77 | selectCaptureDevice | 留-平 | 同上（windows 每次开麦与热插拔重开都走它） |

**B3. 其余（3 条）**

| 位置 | 符号 | 判定 | 理由或归属 |
|---|---|---|---|
| internal/plugin/manifestslot/command_job_other.go:13 | field exited | 留-ABI | 非 Windows 占位结构体（文件头言明为 cross-vet 能编译而存在，D7 Windows-only；批 1 同判） |
| internal/tools/system_bridge.go:274 | argInt | 留-平 | windows 构建有 15 个真调用点（system_backend_windows.go:900 等）；批 2 现核推翻批 1 的判例，沿用 |
| internal/tools/system_bridge.go:297 | argBool | 留-平 | 同上（1 个真调用点） |

### C. windows 独有（82 条）

**C1. internal/ball/win32_windows.go（24 条，:175-:271、:332）**

| 行 | 符号 | 判定 | 理由或归属 |
|---|---|---|---|
| 175/177/178/180/182/183/185 | style / clsExtra / winExtra / icon / bkgnd / menuName / iconSm | 留-ABI | WINDOWCLASS 镜像，按 ABI 整表声明（批 1 同判） |
| 220-225 | compression / sizeImage / xPelsPerMeter / yPelsPerMeter / clrUsed / clrImportant | 留-ABI | BITMAPINFOHEADER 镜像（批 1 同判） |
| 230 | bmiColors | 留-ABI | BITMAPINFO 镜像（批 1 同判） |
| 247 | hoverTime | 留-ABI | TRACKMOUSEEVENT 类镜像 |
| 262-271 | dwState / dwStateMask / szInfo / uVersion / szInfoTitle / dwInfoFlags / guidItem / hBalloonIcon | 留-ABI | NOTIFYICONDATAW 镜像（批 1 同判） |
| 332 | winHeapPtr | 删-c（归 ball 写面腿） | 现核零引用（批 2 判例：文档指向的 dumpDIB 并不用它）；⛔ ball 非本腿写面 |

**C2. internal/tools/system_backend_windows.go（21 条）**

| 行 | 符号 | 判定 | 理由或归属 |
|---|---|---|---|
| 35/36 | ole32 / shlwapi | 留-口 | 下方 15 个 lazy proc 的宿主 DLL 句柄，整组同进退（批 1 同判）；归属系统后端接线票 |
| 60/63/72/76-87 | procOpenInputDesktop / procMonitorFromWindow / procGetDpiForWindow / procGetWindowIcon / procOleInitialize / procOleUninitialize / procDwmGetWindowAttr / procMapVirtualKeyExW / procGetAsyncKeyState / procIsWindowEnabled / procGetWindowLongPtrW / procGetAncestor / procRegisterClipFmt / procGetClipboardSeqNum / procGetKeyState（15 枚） | 留-口 | 系统动词（桌面枚举/DPI/图标/OLE/键态/剪贴板序列号）的 lazy proc 面，删一半留一半会把 API 面撕开（批 1 同判） |
| 105/145 | gaRoot / dwmExtendedFrameBounds | 留-口 | 同组的常量 |
| 808/809 | procGetWindowIconSafe / dwmGetWindowAttrSafe | 留-口 | "Safe" 后缀组，随系统后端票 |

**C3. internal/ball 其余（15 条）**

| 位置 | 符号 | 判定 | 理由或归属 |
|---|---|---|---|
| d2d_windows.go:139/144/145 | typ / usage / minLevel | 留-ABI | D2D1 属性结构镜像（批 1 同判） |
| orb_windows.go:165 | fTransitionOnMaximized | 留-ABI | DWM_MARQUEETIMERS 镜像（批 1 同判） |
| orb_windows.go:193 | orbHwndTopmost | 留-试 | 批 2 现核推翻批 1：winlive 测试 4 处真引用（`orb_glass62_live_windows_test.go`/`orb_live_windows_test.go`），删必红 winlive 编译 |
| orb_windows.go:466 | orbDataLock | 留-口 | 票 408r3（`da5d559b7`，11-08）落的**句柄寄存处**：:508 写入、句柄开着本身就是 sweep 的活性信号（注释言明 held OPEN until process exit）——「只写不读」是设计，不是死变量；归属 orb 数据目录/sweep 票。⚠ 附带发现：**若后续票确认无人读它，建议改为 `_ = orbDataLock` 注册行或具名注释消音，而不是删赋值**（删赋值＝句柄被 GC 语义收走的风险面） |
| ball_windows.go:118/160 | fadeStart / startedErr | 留-口（归 ball 写面腿） | 批 1 判「删-2」、批 2 沿用「⛔ ball 非其写面」；本腿同界，判例不变 |
| ball_windows.go:1113 | (*Ball).recenterAt | 留-口 | 球体交互入口，随渲染/交互票（批 1 同判） |
| renderer_windows.go:29 | capW / capH（2 条） | 留-ABI | 随 ensureDCAndDIB 的分配路径（批 1 同判） |
| renderer_windows.go:218/838/938/945 | (*renderer).resize / overlayGlyphPixels / maxWindowEdge / ensureDCAndDIB | 留-口 | DPI 重绑/字形诊断/语义命名/DC+DIB 分配，批 1 逐条判留，随渲染票 |

**C4. internal/winsec（6 条）**

| 位置 | 符号 | 判定 | 理由或归属 |
|---|---|---|---|
| acl_windows_test.go:32 | wellKnown | 留-试 | 测试常量（批 1 同判） |
| seam_guard_windows_test.go:47/49/53/70 | narrowOnlyResolver + Resolve + ResolveAccounted + installs | 留-试 | 缝卫夹具，本文件内互引成网（12 处），删需整网一起判（批 1 同判） |
| leftstanding_354_windows_test.go:68 | leftStandingInjected354r1 | 留-试 | 票 354r1（`393430dc9`）测试常量表成员：注释 declaring「the injected refusal this leg names as its own」；同表其余 5 枚常量都在用，独它闲置＝夹具冗余而非死逻辑；随 354 族测试票裁（若裁删，一行的事） |

**C5. 其余 windows-only（12 条）**

| 位置 | 符号 | 判定 | 理由或归属 |
|---|---|---|---|
| internal/tools/screen_capture_windows.go:42/43 | procGetDeviceCaps / procSetCaptureDpiAwareness | 留-口 | 屏幕捕获 API 面，随捕获票（批 1 同判） |
| internal/tools/system_dpi_windows.go:16 | setDPIAwarenessSafe | 留-口 | DPI 感知设置入口，随系统后端票（批 1 同判） |
| internal/risk/anchor_world_334_test.go:55/60 | hostilePOSIX347 / danglingTarget347 | **冻** | `internal/risk/**` 冻结（台账 A35 首裁），本票与后续任何批次都需专项授权（批 1 同判） |
| internal/audio/aecspike/rig_windows.go:49 | iidIUnknown | 留-ABI | COM IID GUID 镜像（批 1 同判） |
| internal/audio/wasapi_windows.go:49 | iidIUnknown | 留-ABI | 同上 |
| internal/observe/logprivate_132_windows_test.go:99 | dirIsSealed132 | 留-试 | 密封目录断言 helper（批 1 同判） |
| cmd/balldebug/shot_windows.go:36 | pGetModuleHandle | 越界-cmd（倾向删） | 现核零引用（批 2 判例）；`cmd/**` 非本腿写面，**归下一波** |
| cmd/wisp/appearance_wave_windows_test.go:286 | floatPtrWave | 越界-cmd（倾向删） | appearance wave（`7f3177f93`，11-11）落下的 helper，全仓唯一命中＝定义行；**归下一波** |
| cmd/wisp/panel_host_windows_test.go:443 | pidSetIntersect | 越界-cmd（留-试） | 面板测试夹具（批 1 同判）；**归下一波**核对 |
| cmd/wisp/quota_44r5_test.go:32 | wroteIt | 越界-cmd（倾向删） | 批 2 现核「只写不读……全文件零引用」；**归下一波** |
| cmd/wisp/resident_ball_windows.go:674 | currentWakePauseRead | 越界-cmd（留-口） | 面板桥诊断读数入口，随面板票（批 1 同判）；**归下一波** |
| cmd/wisp/resident_voice_41a2r1_windows_test.go:485 | voice41a2r1FileExists | 越界-cmd（留-试） | 语音驻留测试夹具（批 1 同判）；**归下一波** |
| cmd/wisp/resident_voice_engines_windows.go:758 | (*voiceResidency).waitAcquire | 越界-cmd（留-口） | 批 1 表曾把同文件 :373 problemsSnapshot 与 :719 waitAcquire 并记「留-口」；今天 problemsSnapshot 已有 79 处真调用不再报，waitAcquire 仍闲置——判例沿用「留-口」，**归下一波**核对是否被语音驻留票接走 |
| cmd/wisp/resident_voice_windows.go:519 | warmDiscards | 越界-cmd（留-口） | 驻留统计面，随语音驻留票（批 1 同判）；**归下一波** |

**AC#4 汇总**：132 条 ⇒ 判删候选 **3 条**（shell_exec.go out/errs＋shell_session_test.go closeErr，
全部零引用已现核）；冻 **2 条**（risk）；越界-cmd **8 条**（cmd/wisp 7＋cmd/balldebug 1，归下一波）；
其余 **119 条**留（ABI 镜像／入口面／夹具／平台不对称面）。**本腿一条未删（派单 ⛔）**。
批 1/批 2 已兑现的删除（SA4000 修＋11 条死码）在本表里只以「已消失」的形态出现，不重复占行。

## AC#5 — 门禁（导出快照，rc 现跑）

本腿**零行 Go 代码改动**，无受影响包；门禁按「基线全绿」口径交，AC#2 五条的修后行为证据由
前批证据＋本腿 `internal/observe -count=2` 复跑兜底：

| 门禁 | 读数 |
|---|---|
| `go build ./...`（GOOS=windows，快照） | **rc=0** |
| `go vet ./...`（GOOS=windows，快照） | **rc=0，零输出** |
| `gofumpt -l . tools/d22scan tools/mockllm`（快照） | 只标 `.scratch\wisp\probes\**` 下 **20 枚 0927 时代的探针件**（tracked，e13196bd4 落的；非本腿产物、非 internal/cmd 面）；**internal/cmd/tools/scripts/frontend 五树＝0 行**。CI 的 gofumpt 步（ci.yml:308）显式只扫这五树 ⇒ 与 CI 同形读数＝0 行（另跑一遍 `gofumpt -l internal cmd tools scripts frontend`＝**0 行 rc=0**） |
| `sh scripts/d22scan.sh`（GOOS=windows，快照） | **rc=0，clean**；scope：bans#1-5 internal/=402、cmd/=74、ban#6 frontend/=122、ban#7 internal/tools/=51、ban#10 internal/risk/=21、ban#8 internal/=1003、cmd/=270、frontend/=122、design/=14（台账各 scope 不降：较批 2 的 387/64/119/48/21 全为持平或随代码自然增长） |
| `go test ./internal/observe -count=2 -v`（SA4000 所在包，主工作树） | **rc=0，`ok`**；`=== RUN` **238 行**＝不同测试名 **119 × 2**（整包：goroutine_test 的 TestPanicInWorkerSurvivesAndCancelsRoot／TestNoopTaskReturnsToBaseline 全在内，另含 quota/log/slo/diagnostics 诸件；名单逐行在日志）；PASS 238 / FAIL 0 / SKIP 0（逐行见 `.scratch/wisp/probes/122/r1/ac5-observe-count2.log`） |

（注：本腿未改任何包，`-count=2` 四数只对 SA4000 所在的 observe 包现跑兜底；cmd/wisp 全套
属越界写面，本腿不碰也不跑全套——批 1 已登记其环境红基线。）

## 写面越界检查名册（交付要求）

- 本腿改动文件全集（拟 commit）：`docs/evidence/s1/122-r1-staticcheck-backlog.md`（本表）、
  `.scratch/wisp/probes/122/r1/ac1-windows-head.txt`、`ac1-linux-head.txt`、`ac5-observe-count2.log`。
  **零 .go 文件改动，零删除。**
- 命中即跳过（具名）：SA1019 `internal/risk/proccpu_windows_test.go:20`（⛔ risk 冻结）；
  AC#4 表里全部 `internal/ball/**` 与 `cmd/**` 条目（⛔ 越界写面）。**越界＝0**（未改任何越界文件）。
- ⛔ 三禁自查：未加任何 `//lint:ignore`／`//nolint`；未改任何 checker 开关；未动 `.github/workflows/ci.yml`
  与任何 lint 步。清账判据一字未动。

## AC#6 — next=（本腿只 commit 不 push；CI 读数由编排者取）

推送后看 CI run 里 **`lint` job 的 staticcheck 步**（.github/workflows/ci.yml:383 起，root＋d22scan＋mockllm
三 module；CI 钉 `@2026.2.2` 与本腿同版）：

- 期望：**GOOS=windows 103 条 / GOOS=linux 57 条**（root module；d22scan 与 mockllm 两 module 期望各 0 条）；
  步自报行 `staticcheck self-report: ... findings=...` 应与上数一致（CI 侧还含两 module 的 0，其
  findings 口径以步日志为准）。
- ⚠ **CI lint 步当前应恒红，红因＝SA9009 `frontend/embed.go:4` 这一条已知假阳性**（frontend 归
  外部 agent，无票能修它，除非授权重排那句话的第一个词）。编排者取读数时请把「红」按
  findings=103/57 核，别把假阳性那 1 条当新债。
- `gofumpt` 步（:307，v0.12.0 钉版）期望 0 行；`d22scan` 步期望 rc=0。
- 取 run id ＋ job id ＋ step 号回填后关 AC#6。**不要用本表 rc=0 冒充 CI 绿。**

## 给编排者的三行账

1. **票面 AC#2 的五条在 HEAD 已全部有着落**：2 条前批修掉、3 条冻结/越界不修——本腿的活
   实际落在「证明这件事」＋AC#4 的 132 行表＋linux 分母被编译噪音吃掉半边的更正。
2. **判删候选只剩 3 条**（shell_exec.go ×2＋shell_session_test.go ×1），全在 `internal/tools`——
   若下波给 tools 写面，一腿可清；其余 119 留／2 冻／8 越-cmd。
3. **linux 读数从此可信**：ball linux 编译面已修（票 418），57 这个分母不再是「噪音压低的数」；
   后续票以 103/57 为新基线（同版本 staticcheck 2026.2.2 下）。
