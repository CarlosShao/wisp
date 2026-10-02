# 244-c1 只读普查 —— 构建入口名册 / PE 头现量 / 修法形状与影响面 / 两条入口对控制台的需求差异

## §0 起手锚

- 时刻：**2026-10-02 19:09 +0800**（起手 `date` 现取）。
- 锚点：进场 `git log -1` 现读 **`f17b1165ec4d6828891c68a7d1cf561fb3597b4e`**（2026-10-02 19:05:26 +0800），分支 **`dev`**。本文全部行号与读数对这枚锚负责。
- 票面：`.scratch/wisp/issues/244-spec-11-wants-the-gui-subsystem-build.md`（AC 框一枚未碰）。
- 规格要求（票面凭据，本腿复读原文）：`docs/specs/SPEC-11-build-deploy-containerization.md:50` 逐字「CLI 与 GUI 同一二进制：无参 = GUI（`-H=windowsgui`）；`wisp run` 子命令 `AttachConsole(ATTACH_PARENT_PROCESS)` 输出。【SPEC】」。
- 纪律自证：**全程零 Go 命令**（含 `go build`／`go version`）；PE 头用现存 MSYS2 objdump（`/e/work/base/msys64/mingw64/bin/objdump`，GNU objdump 2.47.20260726）现量；`docs/BUILD.md`（S0 冻结件）只读一字未动；`frontend/**`／`design/**` 零读零写。
- 派单范围说明：grep 根按 hard rules 用 `cmd internal tools docs scripts .scratch`；`.github/**` 里的构建入口（ci.yml）以只读 grep 补覆盖，未改一字。

## §1 构建入口名册与旗标现状

**结论先给：全仓所有产生 `wisp.exe` 的入口没有一处传 `-H=windowsgui`；`-H` 字样只活在 `docs/BUILD.md` 的"推迟记录＋临时构建追述"两行里。**

| # | 入口 | 位置 | 构建命令形状 | `-ldflags -H=windowsgui`？ |
|---|---|---|---|---|
| 1 | **`scripts/build.ps1`（唯一权威主入口）** | `scripts/build.ps1:103-110`、`:123` | `go build -trimpath -ldflags $ldflags -o build\wisp.exe ./cmd/wisp`；`$ldflags` 六枚**全是 `-X`**（`Version`／`Commit`／`BuildDate`／`DefaultEnv`／`SherpaOnnxVersion`／`OnnxRuntimeVersion`，`:104-109`），**零枚 `-H`** | **否** |
| 2 | `.github/workflows/ci.yml`（三处 step） | `ci.yml:453`（test-windows cgo build smoke）、`:550`（slo-smoke）、`:613`（slo-full） | 三处全部 `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/build.ps1 -Env dev`——CI 不自带旗标，全跟 1 号走 | **否（继承 #1）** |
| 3 | `.github/workflows/slo-fresh.yml` | 全文 | 只跑 `scripts/slo-freshness.sh` freshness 探针（查 GitHub API），**零构建** | 无此维度 |
| 4 | `scripts/slo-check.ps1` | `:69`、`:104-108` | 不直接构建；`-WispExe` 默认 `build\wisp.exe`，缺失时转手调 build.ps1（`:104` 逐字 "wisp.exe missing; building (scripts/build.ps1 -Env dev)"） | **否（继承 #1）** |
| 5 | `scripts/dev/ball-cycle.ps1` | `:29-31` | `go build -o $env:TEMP\wisp-balldebug.exe cmd\balldebug`，零旗标 | **否**（balldebug 台件，非出厂件） |
| 6 | `scripts/spike/run.ps1` | `:40-53` | `go build -o <bin> [-tags <t>] ./$Pkg`，零 ldflags；产物落 `scripts/spike/bin/*.exe` | **否**（spike 台件） |
| 7 | `scripts/sign-models.ps1` | `:44,52,64` | `go run ./tools/signmodels …`，不留 PE 产物 | 无此维度 |
| 8 | **测试台件（现做 wisp.exe 的）** | `cmd/wisp/secret_argv_windows_test.go:161-174`（`buildWispForTest`：`go build -o <tmp>/wisp.exe ./cmd/wisp`，无 `-ldflags`）；同族 `cmd/wisp/resident_approval_live_246_windows_test.go:417`（esclistener）、`cmd/wisp/panel_host_gate_test.go:407`（`go build ./...`）、`tools/d22scan/selftest_test.go:181` 与 `scan_test.go:2256` | 全部零旗标 | **否** ⇒ 包内台件二进制**天生 CUI**（与票面 §5 G5、`244-a2` 读数一致，本腿复读 `:174` 复认成立） |
| 9 | `docker/mockllm.Dockerfile` | `:10` | `CGO_ENABLED=0 go build -trimpath -o /out/mockllm .`（Linux 容器，无 PE 维度）；`docker/builder.Dockerfile` 无 go build | 无此维度 |
| 10 | Makefile 族 | 根目录与 `scripts/**` | **不存在**（`Makefile`/`justfile`/`Taskfile`/`*.mk` 零命中） | — |
| 11 | `docs/BUILD.md`（**冻结件，只读**） | `:52-53`、`:87`、`:90` | `:87` 逐字「windowsgui 子系统切换（`-H=windowsgui`）**推迟到票 07**」；`:90`「用 `go build -ldflags "-H=windowsgui"` **临时构建**后…」——全仓文档里 `-H` 的唯一出现处，都是**追述不是构建链** | **否** |

**旗标 grep 复核**（显式根，ripgrep）：`windowsgui` 产码/脚本命中＝**两行注释**（`cmd/wisp/console_windows.go:26`、`cmd/wisp/console_other.go:7`）＋文档三处（SPEC-11:50、BUILD.md:87/90）＋台账/证据件转述；`-H=`/`-H ` 在 `scripts/` 与 `.github/` 里**零命中**。与票面"构建链里根本没有这一维"一致。（票面另记一枚二进制 blob 命中 `scripts/spike/bin/goja-caps.exe`——ripgrep 默认跳二进制，本腿未复现该命中；已另行 objdump 量过该文件 subsystem＝CUI，见 §2。）

## §2 PE 头现量（objdump -p，逐字引）

量法：`objdump -p <exe> | grep -i '^[[:space:]]*Subsystem'`。**台上 19 枚 exe，subsystem 全部是 `Subsystem 00000003 (Windows CUI)`，零枚 GUI(2)**。逐枚：

| exe（盘上路径） | mtime | 字节数 | objdump 逐字行 |
|---|---|---|---|
| `build/wisp.exe` | 2026-09-30 23:43 | 30,299,552 | `Subsystem		00000003	(Windows CUI)` |
| `build/balldebug.exe` | 2026-09-30 11:21 | 7,295,488 | `Subsystem		00000003	(Windows CUI)` |
| `build/wisp228.exe` | 2026-09-30 16:14 | 30,184,669 | `Subsystem		00000003	(Windows CUI)` |
| `build/wisp77.exe` | 2026-09-21 13:44 | 27,820,531 | `Subsystem		00000003	(Windows CUI)` |
| `build/wisp77_partial.exe` | 2026-09-21 13:44 | 27,555,827 | `Subsystem		00000003	(Windows CUI)` |
| `build/sloprobe.exe` | 2026-09-20 07:30 | 2,906,112 | `Subsystem		00000003	(Windows CUI)` |
| `build/33p1-receipt.exe` | 2026-10-01 12:48 | 6,009,856 | `Subsystem		00000003	(Windows CUI)` |
| `build/33p1-reentry.exe` | 2026-10-01 12:45 | 6,044,672 | `Subsystem		00000003	(Windows CUI)` |
| `wisp.exe`（仓根，`.gitignore:2 *.exe`） | 2026-10-01 14:59 | 30,930,354 | `Subsystem		00000003	(Windows CUI)` |
| `balldebug.exe`（仓根） | 2026-09-21 08:40 | 7,069,696 | `Subsystem		00000003	(Windows CUI)` |
| `signmodels.exe`（仓根） | 2026-09-25 23:16 | 7,838,720 | `Subsystem		00000003	(Windows CUI)` |
| `tmp-none/x211check.exe` | — | — | `Subsystem		00000003	(Windows CUI)` |
| `scripts/spike/bin/` 7 枚（goja-caps／model-residency／shell-baseline／speech-baseline／webview2-latency／xy-verdict／xy-verdict-nocgo） | — | — | 逐枚 `Subsystem		00000003	(Windows CUI)` |

三枚要点：

1. **票面立票读数与本腿读数对不上（产物流转痕迹，不是矛盾）**：票面记 `build/wisp.exe`＝09-30 **11:21**／29,771,385 B；本腿现量＝09-30 **23:43**／30,299,552 B。即**立票（16:1x）之后这枚产物又被重建过一次**，但 subsystem 读数未变（仍 CUI(3)）——票面结论未被推翻，反而多一层"台面产物会被后来者覆盖、出生命令不可自证"的证据（票面 AC#1 已点名 `build/**` 被 `.gitignore:14` 掉、交件必须带命令全文）。
2. `build/balldebug.exe`（09-30 11:21）与票面引 243-c1 表第 23 行一致，同为 CUI。
3. "台上无产物"不成立——是"**台上产物全 CUI**"；`build/**` 族属出厂链产物位（被 gitignore），仓根三枚与 `tmp-none/`、`spike/bin/` 均为手工/台件残留。

## §3 修法形状与影响面（只列形状，不选形）

### 3.1 修法形状清单

- **形状甲（改权威出口）**：`scripts/build.ps1:103-110` 的 `$ldflags` 数组**加一枚元素** `"-H=windowsgui"`（该数组 `-join ' '` 后整体喂给 `:123` 的 `-ldflags`，追加一枚 `-H` 与六枚 `-X` 并存无冲突）。影响面：唯一权威出口，CI 三处 step（§1 #2）与 slo-check 缺件自建路径（§1 #4）**自动跟随**，一次改对全线。GOOS/GOARCH 已被脚本钉成 windows/amd64（`:117-118`），`-H` 无交叉编译语义问题。
- **形状乙（改别的入口）**：不动 build.ps1，在 slo-check 自建分支或测试台件加旗标——会造出**第二真相源**（slo 腿量的是出厂 exe，改台件等于没改出厂；CI 走 build.ps1 也吃不到）。仅登记，不展开。
- **形状丙（保持 CUI＋运行时 FreeConsole）**：票面 §5 G4 已裁"与 SPEC-11:50 字面 flag 要求冲突"；`244-a2` 现量 `FreeConsole|AllocConsole` 在 `cmd/ internal/ tools/` 零命中，属新产码。本腿不复跑，只转抄裁决。

### 3.2 对控制台输出一族的影响面（静态读码，未真机）

- **attach 的判据链**（`cmd/wisp/console_windows.go:33-43`）：`attachParentConsole` 先 `GetStdHandle(STD_OUTPUT_HANDLE)`，**已有可用 stdout 就直接 return**（`:35`）；否则 `AttachConsole(ATTACH_PARENT_PROCESS)`＋重绑三枚 std handle（`CONOUT$`/`CONIN$`，`:40-42`）。三种场景推演：
  - **父控制台里跑**：切换后进程自身无控制台 ⇒ stdout 无效 ⇒ attach 父控制台 ⇒ 输出落回调用方终端（`docs/BUILD.md:90-92` 当年临时构建追述过，但那是**过期读数**，票面 AC#0 明说不可继承）。
  - **重定向**（`wisp doctor > out.txt`）：stdout 合法 ⇒ attach 直接 return ⇒ 不受影响。
  - **explorer 双击**：无父控制台 ⇒ `AttachConsole` 失败（函数注释"Fails harmlessly"）⇒ 黑框消失＝切换目的；代价见 3.3。
- **attach 覆盖清点（dispatch 面）**：`cmd/wisp/main.go` 无参与八枚子命令分支全在 dispatch 处调 `attachParentConsole()`（`:65` 无参、`:90` run、`:93` providers、`:96` doctor、`:108` models、`:113` panel-assets、`:119` panel-inbound、`:122` version、`:125` help、`:128` default）；**两枚例外在各自函数体内调**——`secret`（`cmd/wisp/secret.go:183`）与 `slo`（`cmd/wisp/slo_windows.go:254`）。静态覆盖是**全的**；CI 的 slo 采样进程是 PowerShell 子进程（切换后为 GUI 子系统的控制台父），`slo_windows.go` 的 attach 分支与重定向 JSON 输出是否照吃，本腿**没有真机读数**，归票面 AC#2 的三条真跑。
- **已知会被切换改变的行为面（转抄票面＋源码复读）**：AC#5 三把凭据里的 ①②——`interactiveStdin()`（`cmd/wisp/approval_reply_stdin_windows.go:41-52`）以 `GetConsoleMode` 成功为判据；Explorer 拉起的 GUI 进程 `GetStdHandle(STD_INPUT_HANDLE)` 无控制台 ⇒ 返回 nil ⇒ 常驻腿任务入口"控制台键盘"一支关闭（`resident_task_source_windows.go:218-226` 打印 `taskEntryDisabledClaim` 后照常带球常驻）。**停机面**：`resident_windows.go:107` 的 `os.Interrupt/SIGTERM` 注册在先（注释 ：93-99），但"双击进程收不到 Ctrl+C"的形状归票 228 AC#11 与票 43，本腿不展开。

## §4 两条入口对控制台的需求差异（run 腿 vs 常驻腿）

**核心差异一句话：`wisp run` 的 stdout 审计行族是**它的交付物本身**（结果、审计、卡片全靠打印）；常驻腿对控制台是**可选增值**（有则白得一条任务入口＋答复路，无则打印降级声明、球照常）。**

| 维度 | `wisp run`（`cmd/wisp/main.go:89-91` → `run.go`） | 常驻腿（`main.go:60-67` → `resident_windows.go`） |
|---|---|---|
| attach 时机 | `main.go:90` dispatch 处先 attach，`cmdRun` 之后才跑 | `main.go:65` attach 后进 `runResident()` |
| **stdout 依赖度＝高** | ① `run.go:931-935` `auditf`：`[audit] …` 行 **Fprint 到 stderr**＋落盘 sink（"stderr is the operator standing at a terminal"；注释自陈双通道设计）；② `run.go:1260-1300` `consoleSink.Publish`：流式回复、`[工具 …]`、`[错误 …]` 全 `Fprint(c.out)`；③ `run.go:1143-1152` 收尾三行（`任务 … 结束`／`回复已完成`）打 stdout；④ `consoleApprovalUI.Prompt`（`:1344-1389`）把确认卡正文打 `u.out`；⑤ 启动即 `printVersions`（`main.go:152`）。**run 腿不打印＝这腿没交付** | `resident_windows.go` 的 Printf 面（`:45` 二次启动提示、`:78/81` 关停序列、`:85` boot 报告、`:146` 面板不可用、`:216/227/229-233` 事件循环状态）**全部是给人看的旁证**，落盘真相走 `installLogSink`（`:63`，票 117 注释 ：52-57 逐字："double click the icon, no terminal attached, stderr going nowhere… needs a listener on disk"） |
| **stdin 依赖度** | `run.go:153` `interactiveStdin()` 非 nil 才有答复路；nil 时 `main.go:154-158` **打印响亮声明**（"L2 卡会等到超时后按拒绝处理…"）后照跑 | 同一枚 `interactiveStdin()`（`resident_task_source_windows.go:218`）作**任务入口开关**：nil ⇒ 打 `taskEntryDisabledClaim`（`:85-88`）⇒ 常驻腿**继续常驻**（文件头 ：29-33 逐字"A MISSING CONSOLE IS NEVER A REASON TO STOP"） |
| 控制台丢失时的形态 | **腿本身失效一半**（输出不可达＋无人答复卡片），退出码与任务仍会走完 | **腿不受损**：球、托盘、审批门照装；只是"任务来源"与"答复通道"两格降级（`src.surface` 仅 `console != nil` 时建，`:287-290`） |
| run 腿要打印审计行吗 | **要**。`auditf` 的 stderr 半边就是 run 腿的现场审计（票 105 注释 ：918-925 明写两半各有其主）；"标准输出＋审计行"族＝派单点名的那族，全在 attach 覆盖之后才写 | 常驻腿**不需要** stdout 审计行（真相在 sink 落盘），但它的 Printf 面在**有控制台**时是操作员唯一的即时反馈，切换后从终端启动这条仍靠 attach 通 |

**差异对切换的含义（只陈述，不裁）**：`wisp run`／`doctor` 的交付对 attach 通路是**硬依赖**——切换后这条通路一旦不成立，run 腿从终端发起的输出与答复同时归零且无落盘兜底（auditf 落盘那半在，但流式回复与收尾行不落盘）；常驻腿是**软依赖**，降级路径已在码里打印出来。这正是票面 AC#0「只加 flag 不验 CLI＝把 CLI 那条腿盲切」与 AC#2 三条真跑的形状依据。

## §5 我可能判错的条目＋量不到的格子

**可能判错**：

1. **§1 #8 的枚举可能不全**：`_test.go` 里现做 exe 的位我只 grep 了 `exec.Command(… "build"` 形状；若有经 `go test -c`、`go run` 或脚本间接现做 wisp.exe 的测试路径，我可能漏。已覆盖的四枚：`secret_argv_windows_test.go:174`、`resident_approval_live_246_windows_test.go:417`、`panel_host_gate_test.go:407`、`tools/d22scan/{selftest,scan}_test.go:181/2256`。
2. **`ci.yml` 行号会漂**：`:453/:550/:613` 对锚 `f17b1165` 有效；该文件注释极长，任何新 step 都会推移行号——引用时以"三处 `scripts/build.ps1` 调用"为准而非行号。
3. **attach 场景推演是静态的**：§3.2 三场景与"slo 采样子进程"一格全部是读码推演，**没有任何真机读数**；若 Windows 对 GUI 子系统进程的 std handle 语义与推演不符（如重定向句柄继承细节），判错处在此。
4. **`build/wisp.exe` 的归属**：我按 mtime 与大小判它出自 build.ps1 链（含 `-X` 六枚），但**没有**读它的版本资源段核对 commit——不排除它是某腿手工 `go build` 的产物；这对"subsystem=CUI"的结论无影响，对"出厂链产物"的标签可能错。
5. **spike/bin 七枚的构建者**：我按 `scripts/spike/run.ps1` 判它们出自该脚本，但盘上文件 mtime 未逐枚核对与脚本运行的对应关系。

**量不到的格子**：

1. **切换后的真机三形态读数**（父控制台跑／重定向／explorer 拉起）——票面 AC#2 的三条真跑，本腿零 Go 命令＋不跑产码，量不到。
2. **`wisp doctor > out.txt` 在切换后的实际行为**——同上，AC#2 ②。
3. **包内台件对 subsystem 维度的失明程度**——"台件天生 CUI 所以包内断言看不见这一维"是票面 G5 的裁决，本腿未种变异复跑验证。
4. **`.scratch/wisp/issues/` 全票池的 `windowsgui` 现查**——票面已现跑过（两枚 -done 票），本腿用 ripgrep 复核了产码/脚本/文档根，未对 issues/ 全量重跑同一把尺（票面读数未受挑战）。
5. **`docs/BUILD.md` 的冻结状态由谁守**——本腿确认它一字未动，但"冻结件清单"的权威文本在 `SPEC-12`／票 244 禁区节，我只转抄票面禁区，未独立读 SPEC-12 原文核对。
