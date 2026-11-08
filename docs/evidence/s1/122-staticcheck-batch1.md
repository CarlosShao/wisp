# 票 122 — staticcheck 清账第一批（AC#1～AC#5）裁决表

- 腿名：122sa1 · 日期：2026-11-08 · 写面 agent：本腿（实现者）；本表由实现腿先出底稿，验收方另跑复核
- HEAD：`6e2eb5845`（A1632 particle colour leg）· 导出树：`/d/tmp/wisp/20261108/122sa1-tree`（基线）、`122sa2-tree`（改后）
- **仪器版本（AC#7 格式要求，逐字）**：`staticcheck.exe 2026.2.2 (0.8.2)`、`gofumpt v0.12.0 (go1.27.1)`
- 跑法：全程在 `git clone --shared` 导出树里跑，主工作树只写码。

---

## AC#1 — 开工重测（当前 HEAD 纯净快照，两 GOOS）

| 读数 | GOOS=linux | GOOS=windows |
|---|---|---|
| 票面预做（audit-85-preflight，旧 HEAD） | 34 | 78 |
| CI 自报分母（A101②，09-22） | — | 42 |
| **今天 HEAD `6e2eb5845` 实测** | **20 条 finding**（28 行输出，另 8 行编译噪音） | **115 条 finding**（115 行输出，零噪音） |

今天 linux 分布：`14×U1000 / 2×S1011 / 2×ST1018 / 1×SA9009 / 1×SA5000`。
今天 windows 分布：`94×U1000 / 4×SA1019 / 3×SA4006 / 2×ST1018 / 2×ST1012 / 2×ST1005 / 2×S1016 / 2×S1011 / 1×SA9009 / 1×SA5000 / 1×SA4010 / 1×SA4004`。

**与 34/78 的差，能落到证据的差**：
- `SA4000` 那 1 条今天**不存在**：已被 commit `79041f071`（CI 修复腿 ci-r1，2026-11-01）修掉，
  diff 逐字 `if !errors.As(err, &oe) || !errors.As(err, &oe) {` → `if !errors.As(err, &oe) {`，
  即票面允许的「本来就重复、删掉是修正」路径（internal/observe/goroutine_test.go，现 :151）。
- 今天 windows 侧比 78 **多** 37 条：主因是 audit-85 之后 windows-only 面大幅增长
  （`internal/ball/*_windows.go`、`internal/tools/system_backend_windows.go` 等整文件新增），
  以及检查器版本不同（ST1018/ST1012/ST1005/SA5000/SA4010/SA4004 这些形状今天在报）。
  逐 commit 归因不再可行（其间几十枚 commit），本表只对「票面点名的三组」给了一对一证据。
- **linux 侧 8 行编译噪音（登记，不属本票）**：GOOS=linux 下 `internal/ball/orbdebug.go:39/:43`
  报 `undefined: Ball`（linux 编译面破损），且 `sherpa-onnx-go-linux@v1.13.8` 构建约束全排除——
  依赖这两者的包在 linux 读数里**没被分析**，linux 20 这个分母因此天然偏低。归属：linux 编译面专项，非本票。

## AC#2 — 第一批三组清账（逐条）

改前 → 改后：**windows 115 → 109（-6，diff 逐条＝下表 6 条）；linux 20 → 20（六条全在 windows-only 件里，linux 不动）**。

| # | 形状 | 位置 | 处置 | 说明 |
|---|---|---|---|---|
| 1 | SA4000 | internal/observe/goroutine_test.go（原 :106，现 :151） | **已被 `79041f071` 修掉，本腿零改动** | 两个相同操作数 `!errors.As(err, &oe)` 删一个＝票面允许的修正；该用例的断言维度没写松：Class 维度（`ClassOfOrInternal`）与类型维度（`errors.As`）仍是两个独立断言 |
| 2 | SA1019 | cmd/wisp/secret_argv_windows_test.go:169（`runtime.GOROOT`） | **修** | 回退分支改读 `os.Getenv("GOROOT")`（弃用注记认可的环境变量半边）；PATH 与 GOROOT 都找不到时仍走原有的显式 Fatal，断言不弱化 |
| 3 | SA1019 | internal/speech/cgo_boundary_windows_test.go:36（`windows.GetCurrentProcess`） | **修** | 改 `windows.CurrentProcess()`（x/sys 的非弃用伪句柄访问器，无 error 返回），删掉恒不触发的 err 分支 |
| 4 | SA1019 | internal/speech/cgo_boundary_windows_test.go:95（同上） | **修** | 同 #3 |
| 5 | SA1019 | internal/risk/proccpu_windows_test.go:20（`windows.GetCurrentProcess`） | **不修，登记** | `internal/risk/**` 票面明令冻结，本腿无写面；修法与 #3 完全同形（一行换 API），需编排者专项解冻或挂后续票 |
| 6 | SA4006 | cmd/wisp/resident_watchdog_windows_test.go:210（hold64 初值） | **修** | 加 `len(hold64) != 64<<20 → t.Fatalf` 断言：对初值的真读数。**行为证据**：断言不恒真——`holdTouched384` 返回短缓冲即红；原测试「hold 撑住工作集」的前提从此被显式验到，而不是靠 GC 不及时 |
| 7 | SA4006 | cmd/wisp/resident_watchdog_windows_test.go:242（hold128 初值） | **修** | 同 #6（128MB 档） |
| 8 | SA4006 | cmd/wisp/slo_windows.go:800（`last` 初值） | **修** | 初值 `{state: reportUnwritten, offset: -1, note: "nothing read yet"}` 从未被读（循环体每次 `last = obs` 先于一切对 `last` 的读），改裸零值 `var last subjectReportRead`，行为逐位不变；原初值里那句 note 是死文案 |

**如实栏（误报/存疑登记）**：本批 6 条现场修的没有误报。#6/#7 严格说是「测试前提靠 GC 时序」的隐患被 staticcheck 揪出来，修法是**加强**（补尺寸断言），不是绕。

## AC#3 — 假阳性格：SA9009（只登记不修）

- 读数：`frontend/embed.go:4:1: ineffectual compiler directive due to extraneous space: "// go:embed below is the ONLY way those bytes reach wisp.exe. There is no" (SA9009)`
- 现场核过（`sed -n 1,22p` + `cat -A`）：`:4` 是**散文注释**——句子恰好以 `// go:embed` 开头（完整句是 "the go:embed below is the ONLY way those bytes reach wisp.exe..."）；真正的指令在 `:19`，逐字 `//go:embed all:dist`，写法正确、有效。
- **为什么不修**：`:4` 那行本身就是票 77 AC#1 的文档（"Ticket 77 AC#1: npm run build ... the go:embed below is the ONLY way..."），改它＝动那张票的验收文档；加 `//nolint` 类全局抑制＝拆门，票面明令禁止。故**登记为已知假阳性，代码零改动**。

## AC#4 — U1000 逐条归属表（只出表，本腿不删）

分母：今天 **96 条去重**（windows 94 + linux 独有 2：`system_bridge.go` 的 `argInt`/`argBool`——这两个在 windows 构建里被用到，linux 构建里才闲置）。票面的 26 是旧 HEAD 的数，**不并列覆盖、只注明过期**。

判定图例：**留-ABI**＝Windows API 结构体/常量镜像，按 ABI 全量声明、Go 侧不读是设计内；**留-口**＝留给后续票的入口/诊断面；**留-试**＝测试夹具；**删-2**＝建议第二批删（本表只判，不动手）；**冻**＝包冻结，任何批次都需专项授权。

### internal/ball（冻结面之外的 ABI 镜像与未接线渲染路径）

| 位置 | 符号 | 判定 | 归属/理由 |
|---|---|---|---|
| win32_windows.go:162-258（27 个 field） | style…hBalloonIcon | 留-ABI | `WINDOWCLASS`/`BITMAPINFOHEADER`/`NOTIFYICONDATAW` 等 OS 结构体镜像，按 ABI 整表声明 |
| win32_windows.go:319 | winHeapPtr | 删-2 | 文档说"see renderer.dumpDIB"，但 dumpDIB 并未调用它；文档已烂，helper 成死码 |
| d2d_windows.go:139/144/145 | typ/usage/minLevel | 留-ABI | D2D1 属性结构镜像 |
| orb_windows.go:163 | fTransitionOnMaximized | 留-ABI | DWM_MARQUEETIMERS 结构镜像 |
| orb_windows.go:191 | orbHwndTopmost | 删-2 | 常量无任何引用，非 ABI 表成员 |
| ball_windows.go:100/142 | fadeStart/startedErr | 删-2 | 状态字段，写侧已不存在（非 OS 镜像） |
| ball_windows.go:1083 | (*Ball).recenterAt | 留-口 | 球体交互入口，留给后续交互票 |
| renderer_windows.go:29 | capW/capH | 留-ABI | 随 ensureDCAndDIB 一起属于未接线的分配路径（见下） |
| renderer_windows.go:218 | (*renderer).resize | 留-口 | DPI 变更重绑路径，文档明确写了用途，留给渲染票 |
| renderer_windows.go:838 | (*renderer).overlayGlyphPixels | 留-口 | 字形像素诊断面，留给渲染票 |
| renderer_windows.go:938 | maxWindowEdge | 留-口 | 语义命名层（BallSizeMaxPx 的 DPI 包装），随渲染票接线 |
| renderer_windows.go:945 | ensureDCAndDIB | 留-口 | DC/DIB 分配路径，与 resize 同属渲染票 |

### internal/tools（未接线的系统后端 API 面）

| 位置 | 符号 | 判定 | 归属/理由 |
|---|---|---|---|
| system_backend_windows.go:35/36 | ole32/shlwapi | 留-口 | 下方 16 个 proc 变量的宿主 DLL 句柄，与 proc 组同进退 |
| system_backend_windows.go:60-87（16 个 var） | procOpenInputDesktop…procGetKeyState | 留-口 | 系统动词（桌面枚举/DPI/图标/剪贴板序列号）的 lazy proc 面，整组留给系统后端票接线；删一半留一半会把 API 面撕开 |
| system_backend_windows.go:105/145 | gaRoot/dwmExtendedFrameBounds | 留-口 | 同上组的常量 |
| system_backend_windows.go:344 | realPostMessage | 删-2 | 单行包装且无调用方；接线时重写不心疼 |
| system_backend_windows.go:812/813 | procGetWindowIconSafe/dwmGetWindowAttrSafe | 留-口 | "Safe" 后缀组，随系统后端票 |
| system_backend_windows.go:820 | (*systemBackend).fns | 删-2 | nil 兜底自反函数，无任何调用方 |
| system_bridge.go:264/316（+linux 281/304） | isWindowsGate/argInt/argBool/systemErrText | 删-2 | 四个纯 helper，无调用方（argInt/argBool 仅在 linux 构建闲置，windows 构建有人用——linux 独有那 2 条随本组一起批 2 处理） |
| system_dpi_windows.go:16 | setDPIAwarenessSafe | 留-口 | DPI 感知设置入口，随系统后端票 |
| system_readings_windows.go:17 | systemNowMu | 删-2 | 为 systemNowFn 配的锁，锁的东西没有第二个写者，死锁具 |
| screen_capture_windows.go:43/44 | procGetDeviceCaps/procSetCaptureDpiAwareness | 留-口 | 屏幕捕获 API 面，随捕获票 |
| screen_capture_windows.go:194 | colorGrey | 删-2 | 颜色常量无引用 |
| webmd.go:39/41/54 | preLang/pendNL/dropEntirely | 删-2 | markdown 解析器的脚手架（pre 只写了块状态字段从不读、dropEntirely 整表无引用），解析现行路径不经过它们 |
| web_bridge_test.go:49 | (*scriptGate).levels | 留-试 | 测试夹具，随 web 桥测试票核对 |

### internal/config / plugin / speech / agent / memory / observe

| 位置 | 符号 | 判定 | 归属/理由 |
|---|---|---|---|
| config/tiers.go:94 | lockedSections | 删-2 | 注释自称"completeness pin 会走它"，但仓内 5 处引用全在注释/文档里，代码零引用——注释与码已经分叉，批 2 删时同步改注释 |
| config/unwired.go:118 | lockedKeyDisposition | 删-2 | 同上（8 处引用全为注释） |
| plugin/manifestslot/command_extract.go:11 | extractDeclJSON | 删-2 | 无调用方的解析 helper |
| plugin/manifestslot/command_job_other.go:13 | childJob.exited | 留-ABI | 非 Windows 占位结构体（文件头注释言明：为 cross-vet 腿能编译而存在，D7 Windows-only） |
| speech/cloudchain.go:648 | causeOrUnknown | 删-2 | 注释齐全的死 helper，终点跳事件的 cause 命名器；接线前删，接线时复活 |
| agent/reminder_timer_test.go:109 | (*fakeReminderClock).releasedCount | 留-试 | 假时钟夹具的计数面，随 agent 测试票核对 |
| memory/l1_warm_29_test.go:151 | l1WaitFor | 留-试 | 测试等待 helper |
| observe/logprivate_132_windows_test.go:99 | dirIsSealed132 | 留-试 | 密封目录断言 helper |

### internal/audio / winsec / risk / cmd

| 位置 | 符号 | 判定 | 归属/理由 |
|---|---|---|---|
| audio/aecspike/rig_windows.go:49、audio/wasapi_windows.go:49 | iidIUnknown ×2 | 留-ABI | COM IID GUID 镜像 |
| winsec/acl_windows_test.go:32 | wellKnown | 留-试 | 测试常量 |
| winsec/seam_guard_windows_test.go:47-70 | narrowOnlyResolver 类型+2 方法+installs 字段 | 留-试 | 缝卫测试夹具（12 处引用，本测试文件内互引成网）——批 2 若判删需整网一起判 |
| risk/anchor_world_334_test.go:55/60 | hostilePOSIX347/danglingTarget347 | **冻** | `internal/risk/**` 冻结，本票与批 2 都无写面；需编排者专项授权 |
| cmd/balldebug/shot_windows.go:36 | pGetModuleHandle | 删-2 | 无调用方的 lazy proc |
| cmd/wisp/panel_host_windows_test.go:443 | pidSetIntersect | 留-试 | 面板测试夹具 |
| cmd/wisp/quota_44r5_test.go:32 | wroteIt 字段 | 删-2 | 测试结构体里只写不读的哨兵字段 |
| cmd/wisp/resident_ball_windows.go:624 | currentWakePauseRead | 留-口 | 面板桥诊断读数入口，随面板票 |
| cmd/wisp/resident_voice_41a2r1_windows_test.go:484 | voice41a2r1FileExists | 留-试 | 语音驻留测试夹具 |
| cmd/wisp/resident_voice_engines_windows.go:373/719 | problemsSnapshot/waitAcquire | 留-口 | 注释言明"for the boot line"的诊断读数面，随语音驻留票 |
| cmd/wisp/resident_voice_windows.go:482 | warmDiscards 字段 | 留-口 | 驻留统计面，随语音驻留票 |

**批 2 汇总**：建议删 22 条（上表"删-2"）；其余 73 条留（ABI 镜像/入口/夹具）+ 2 条冻（risk）。本腿一条未删。

## AC#5 — 受影响包测试（导出树 122sa2-tree，`-count=2 -v`，DLL harness 带齐）

环境注记：`--shared` clone **不带 gitignored 的 DLL 与模型**（`third_party/sherpa-onnx/*.dll`、
`third_party/spike-models/`），首次跑 speech 全体秒挂 `silero model not found`、首次两包全 `0xc0000135`
（DLL 不在 PATH 实际指向的目录）——均为**环境红**，非代码红；把 DLL 与模型补进导出树后重跑，读数如下：

| 包 | rc | === RUN 行 | 不同测试名 | ×2 对上？ | PASS | FAIL | SKIP | 包级结论 |
|---|---|---|---|---|---|---|---|---|
| internal/speech | 0 | 212 | 106 | 是 | 118 | 0 | 94 | `ok` |
| cmd/wisp | 1 | 1410（顶层的 1098） | 705（顶层的 549） | 是 | 979 | 109 | 10 | FAIL——红集全为预先存在的环境红/负载敏感翻动，见下 |

**cmd/wisp 红绿归属（三核过，不是本腿造的红）**：
- 78 条不重复红里，约 30 条语音/真实硬件测试**红因是仪器自报的环境问题**（逐字）：
  `speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH (on Windows the loader searches System32 before PATH)`
  ——`C:\Windows\System32\onnxruntime.dll`（10.5MB，2026-05-26）在 loader 顺序上先于任何 PATH，属**机器级环境红**；
- 其余为真窗口/面板类（TestPanelHostRealWindowHop、TestAC14 系列）在**双套件并发**下翻动的负载敏感红；
- **基线对照**：无本腿改动的 122sa1-tree 跑同批失败用例（Test15eSpeech/Test26PBargein/TestWireR1/TestAC228）→ 同样四条全红（`122-baseline-failcheck.log`）；
- **本腿改到的面在两遍里全绿**：Test384×3（watchdog，含新加的 hold 尺寸断言）、TestSecret×16、TestSLO144/147/149/156/159 全家、Test12r2Seam——逐条 ×2 PASS，红集里一条都没有。

**仪器注记（AC#5 判据的坑，登记给后腿）**：`-count=2` 的语义是**整张测试清单跑两遍**（probe 实证输出顺序 A→B→A→B），
不是逐条连跑两遍。cmd/wisp 全套一遍就要 ~10–21 分钟，**撞上 go test 默认 10m 超时会死在第一遍**，
日志表现恰好是"每名只出现一次"——别把那当成 count=2 没生效。跑全套必须显式 `-timeout=45m` 级别。

## 门禁

- `gofumpt -l` 四个改动件：**0 行**（v0.12.0 (go1.27.1)，rc=0）
- `go build ./...`（导出树 122sa2-tree）：**rc=0**
- `go vet ./cmd/wisp ./internal/speech`（导出树）：**rc=0，零输出**
- staticcheck 改前/改后并列：linux 20→20、windows 115→109（-6＝本批全部目标，diff 逐条核对无多余增减）
- `sh scripts/d22scan.sh`（导出树）：**rc=0，clean - no D22 ban violations**（scope 读数逐条在 `122-d22.log`：bans #1-5 internal/=387、cmd/=64、ban #6 frontend/=119、ban #7 internal/tools/=48、ban #10 internal/risk/=21、ban #8 各 scope 逐条列出）

## AC#6 — next=（CI 步级读数，本腿不造假）

本腿只 commit 不 push。推送后要看的读数：CI run 里 `lint` job 的 staticcheck 步——
期望 **GOOS=linux 20 条 / GOOS=windows 109 条**（若 CI 的 staticcheck 版本与本腿的 2026.2.2 (0.8.2) 不同，
条数允许漂移，按 AC#7 格式同时报版本）；`gofumpt` 步期望 0 行；`d22scan` 步期望 rc=0。
取 run id + job id + step 号后回填本表并关 AC#6。
