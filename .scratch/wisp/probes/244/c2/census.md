# 244-c2 只读普查 —— 构建入口名册复量 / CUI 残余路径 / subsystem 尺有无 / GUI 档下 stdout 可见性

> 本腿代号 `244-c2`，只读普查腿。今天 2026-10-03。⛔ 全程零 Go 命令、零构建命令、零写产码。
> 前件 A＝`c1`（`.scratch/wisp/probes/244/c1/census.md`，commit `0a88d7a8`，锚 `f17b1165`）：量到 11 枚构建入口全零 `-H`、台上 19 枚 exe 全 CUI。
> 前件 B＝落地件 `scripts/build.ps1`（commit `cc6eaa65`，244-r1）：第 115 行已追加 `-H=windowsgui`。本腿**对着现树复量**，⛔ 不照抄 c1 读数。

## §0 锚

- 时刻（起手 `date -Iseconds`）：**2026-10-03T09:13:13+08:00**。
- 起手 `git log -1 --format=%h`：**`3fe377e3`**；同法取数时 `git status --porcelain -- cmd internal tools docs scripts .scratch | wc -l`＝**355**（共享工作树，多腿在飞，dirty 属正常态）。
- ⚠ 落 §0 骨架这发命令时 HEAD 已漂到 **`72cce7a0`**（HANDOVER 4.0z 停车 commit）——本仓是**活体共享工作树**，锚点会在腿飞行期间被别人推进；本腿所有行号/读数对"取数时刻的盘上内容"负责，取数时 census 相关 8 枚文件（`scripts/build.ps1`／`cmd/wisp/console_windows.go`／`console_other.go`／`main.go`／`config_reload.go`／`run.go`／`.github/workflows/ci.yml`／`docs/BUILD.md`）`git status` **全部 clean＝读的是已 commit 内容，不是别腿的在飞脏写**。
- 票面：`.scratch/wisp/issues/244-spec-11-wants-the-gui-subsystem-build.md`（AC 框一枚未碰；末尾「越序」裁定在台账 `A548`，本腿已读）。
- 规格要求：`docs/specs/SPEC-11-build-deploy-containerization.md:50`（票面凭据）。
- 纪律自证：全程零 Go 命令／零构建；PE subsystem 本腿**量不到**（库内零枚 tracked exe，`*.exe`＋`build/` 均在 `.gitignore:2/14`；sanctioned 法＝objdump 只对 tracked exe，无枚可量——详见 §6）；`docs/BUILD.md`／`ci.yml` 等冻结件只读一字未动；`frontend/**`／`design/**` 零读零写零转述（grep 曾被动命中 frontend 路径，本腿未读其内容、本件不引用）。

## §1 入口名册复量（已带 / 仍不带 `-H`，逐枚）

**一句话结论（对现树复量，非照抄 c1）**：c1 那句"11 枚全零 `-H`"是**落地前**状态；244-r1（`cc6eaa65`）之后，**唯一权威出厂入口 `scripts/build.ps1` 直接带上 `-H=windowsgui`，两台继承者（CI 三 step／slo-check 缺件自建）随动转 GUI ⇒ 出厂链全线带 `-H`**。**仍零 `-H` 而今天还会产 exe 的只剩 3 枚入口，且全是非出厂件**：`ball-cycle`（balldebug）、`spike/run`（spike bins）、测试台件（包内 wisp.exe／esclistener／d22scan／mockllm）。#3/#7/#9/#10 本就没有 PE subsystem 这一维，#11 是文档、其教的那条产线命令继承 #1。

| # | 入口（会产出 exe 的每一枚命令行都算） | 现树位置（复量行号） | `-H=windowsgui`？ | 与 c1 差异 |
|---|---|---|---|---|
| 1 | **`scripts/build.ps1`（唯一权威主入口）** | `$ldflags` 数组＝`:108-116`，**`"-H=windowsgui"` 在 `:115`**，上面挂 5 行 why 注释 `:103-107`（逐字引 `cmd/wisp/console_windows.go:33-43`）；六枚 `-X`（`:109-114`）逐字未动；`go build -trimpath -ldflags $ldflags -o build\wisp.exe ./cmd/wisp` 在 **`:129`**（c1 记 `:103-110`/`:123`，因 r1 加注释行号整体下移 6 行） | **已带** | ★ c1＝不带 → 现带（这是本腿与 c1 的核心差异枚） |
| 2 | `.github/workflows/ci.yml`（三处 step，**冻结件只读**） | `:453`（cgo build smoke）、`:550`（slo-smoke）、`:613`（slo-full）全是 `powershell … -File scripts/build.ps1 -Env dev`，CI 不自带旗标 | **已带（继承 #1）**——〔射程判断，非内容引用：这三 step 产出的 `build\wisp.exe` subsystem 由 build.ps1 的旗标决定，本腿未跑、不改此件〕 | c1＝不带 → 随 #1 转 |
| 3 | `.github/workflows/slo-fresh.yml` | `:70/:78/:83` 只跑 `scripts/slo-freshness.sh`（GitHub-API freshness 探针） | 无此维度（零构建） | 同 c1 |
| 4 | `scripts/slo-check.ps1` | `:69` `-WispExe` 默认 `build\wisp.exe`；`:104` 缺件时转手调 `scripts/build.ps1 -Env dev` | **已带（继承 #1）** | c1＝不带 → 随 #1 转 |
| 5 | `scripts/dev/ball-cycle.ps1` | `:30` `& $go build -o $exe cmd\balldebug`，零旗标 | **仍不带**（balldebug dev 台件，非出厂件） | 同 c1（仍零） |
| 6 | `scripts/spike/run.ps1` | `:46` `$args = @("build", "-o", $exe)`，零 ldflags；产物落 `scripts/spike/bin/*.exe` | **仍不带**（spike 台件） | 同 c1（仍零；c1 记 `:40-53`，本腿复量构建行在 `:46`） |
| 7 | `scripts/sign-models.ps1` | `:44/:52/:64` `& $go.Source run ./tools/signmodels …`，`go run` 不留 PE | 无此维度 | 同 c1 |
| 8 | **测试台件（自己现做 exe 的）** | `cmd/wisp/secret_argv_windows_test.go:174`（`buildWispForTest`：`go build -o …/wisp.exe ./cmd/wisp`，**无 `-ldflags`**，`:158-160` 注释自陈"the ticket's claim is about the shipped command line, not a test harness's … colocates the DLLs the way scripts/build.ps1 does"——**只搬 DLL、不走 build.ps1 旗标**）；`cmd/wisp/resident_approval_live_246_windows_test.go:417`（esclistener testdata exe）；`cmd/wisp/panel_host_gate_test.go:407`（`go build ./...`，多包 smoke、不落单件 exe）；**`internal/llm/openaichat/mockllm_integ_test.go:65`（`go build -o mockllm .`，c1 未列＝本腿补上的一枚）**；`tools/d22scan/selftest_test.go:181` 与 `scan_test.go:2256`（d22scan.exe 工具件） | **仍不带** ⇒ 包内现做的 `wisp.exe` 台件二进制**天生 CUI**（与票面 §5 G5 一致，本腿逐枚复读复认成立） | 同 c1（仍零）；枚举个数比 c1 多 1（mockllm_integ，c1 §5 第 1 条自认枚举个数可能不全，本腿兑现） |
| 9 | `docker/*.Dockerfile` | `docker/mockllm.Dockerfile:10` `CGO_ENABLED=0 go build -trimpath -o /out/mockllm .`（Linux 容器，无 PE subsystem）；`docker/builder.Dockerfile` **只有交叉编译 env（`:9`），全文无 `go build` 命令行** | 无此维度 | 同 c1（builder 无 go build 复认） |
| 10 | Makefile 族 | `Makefile`/`justfile`/`Taskfile`/`*.mk` find 零命中 | — 不存在 | 同 c1 |
| 11 | `docs/BUILD.md`（**冻结件，只读**） | `:45` 教 `powershell … scripts\build.ps1 -Env dev`（→ 继承 #1，现产 GUI）；`:90` 教 `go build -ldflags "-H=windowsgui"` 是**临时构建**（自带 -H，非产线）；`:52-53` 步骤描述只列版本/commit 那族 `-X`、未提 `-H`；`:87` 仍逐字写"推迟到票 07"＝**过期文本**（票面 §5 G6 已裁甲、本票不改它一字） | 文档教的**产线命令继承 #1**（:45→build.ps1） | 同 c1（:87/:90 内容一字未变） |

**旗标 grep 复量（显式根 `cmd internal tools docs scripts .github`，ripgrep）**：产码/脚本里 `-H`／`windowsgui` 命中＝`scripts/build.ps1:103/115`（旗标本体＋why 注释）、`cmd/wisp/console_windows.go:26`、`cmd/wisp/console_other.go:7`（两行注释）；docs 里＝`SPEC-11:50`＋`BUILD.md:87/90`＋台账/证据件转述。**结论：构建链里 `-H=windowsgui` 现在有一枚真旗标（build.ps1:115），不再是 c1 时的"只活在文档两行里"。**

## §2 今天仍会产 CUI exe 的路径名册

**分两层答"哪几条路径今天仍会产出带黑窗口的 exe"（对应派单三点①②③）：**

**① 直接 `go build ./cmd/wisp`（不带 ldflags）——仍产 CUI。** 244-r1 只改了 `scripts/build.ps1` 这一枚权威出口；`-H=windowsgui` 是靠 build.ps1 的 `$ldflags` 注入的，**任何绕开该脚本、手打 `go build ./cmd/wisp`（或 `go build -o 某名 ./cmd/wisp`）的命令行，今天产出的仍是 CUI(3)**。台账 `A548` 第 2 条"下次任何人构建就切了"的正解就是这一层——切的是走 build.ps1 那条，不是裸 `go build` 那条。仓里把这条"裸构建"编码进代码的地方＝测试台件 `cmd/wisp/secret_argv_windows_test.go:174`（`buildWispForTest`，无 `-ldflags`），它今天产出的 `wisp.exe` 恒为 CUI（⇒ 票面 §5 G5 那句"包内断言天生看不见 subsystem 这一维"的形状来源，本腿复认）。

**② 文档教的命令——今天教的事两条都不产黑窗口的出厂件；但有一处过期正文是"手动重建即漏 `-H`"的地雷：**
- `docs/BUILD.md:45`（**冻结件，只读**）逐字教 `powershell -ExecutionPolicy Bypass -File scripts\build.ps1 -Env dev` ⇒ 这条**继承 #1 的旗标，今天产 GUI**，不产黑窗口。
- `docs/BUILD.md:90` 教的 `go build -ldflags "-H=windowsgui"` **临时构建**自带 `-H` ⇒ GUI（且明写是"临时构建"、不是产线）。
- ⚠ 地雷具名：`docs/BUILD.md:52-53` 把构建步骤描述成"CGO_ENABLED=1、CC、-trimpath，ldflags 注入版本/commit/构建时间/DefaultEnv/两个原生库版本"——**只列六枚 `-X` 族、正文没提 `-H`**；`:87` 仍逐字写"windowsgui 子系统切换推迟到票 07"＝**过期文本**（旗标已落、票 07 早结案、这一切换由 244-r1 补做）。⇒ 一个照 `:52-53` 正文**手敲 go build 命令**而不跑 `:45` 那条脚本的人，会**漏掉 `-H` 产出 CUI**。⛔ 本腿不改 `BUILD.md` 一字（禁区）；过期文本的更正责任在票面 §5 G6（已裁甲、登记不改文件），本腿只具名指认行号。
- `.scratch/wisp/issues/244-*.md:15-16` 票面"现量表"记的是 `objdump -p build/wisp.exe` 与 `$ldflags` 六枚全 `-X`——那是**立票时读数**（落地前），本腿已在 §1 复量为"现带"。
- `docs/reports/HANDOVER.md` 里凡 `go build ./...` 字样（`:1190/:1387/:1447` 等）＝历轮"编译门 rc=0"的**叙述**，不是教人怎么产 wisp.exe 的命令行；HANDOVER 不新教构建入口，产线命令一律指回 BUILD.md/build.ps1。⇒ HANDOVER **无一条教人产出 CUI 出厂件**。

**③ CI 构建步——今天产 GUI（继承 #1），⛔ ci.yml 只读只报行号与形状。** `.github/workflows/ci.yml` 三处构建 step（`:453` cgo build smoke／`:550` slo-smoke Build wisp.exe／`:613` slo-full Build wisp.exe）逐字都是 `powershell … -File scripts/build.ps1 -Env dev`，**CI 自身不带旗标、不改 ldflags**（本腿只读该冻结件，未见任何 `-H`／`ldflags` 注入行）。〔射程判断，非内容引用：这三 step 产物 subsystem 由 build.ps1 的 `$ldflags` 决定，244-r1 后＝GUI；本腿零构建、未对 CI 产物做 objdump，此判仅由"CI→build.ps1→带 -H"的调用链静态得出。〕

**净结论**：
- **出厂链（走 build.ps1 的每一条：`:45` 文档命令／CI 三 step／slo-check 缺件自建）今天全线产 GUI，不再产黑窗口出厂件。**
- **今天仍产 CUI exe 的只剩三类**：(a) 裸 `go build ./cmd/wisp`（人工／`buildWispForTest` 台件 `:174`）；(b) `scripts/dev/ball-cycle.ps1:30`→`balldebug.exe`（dev 台件）；(c) `scripts/spike/run.ps1:46`→`scripts/spike/bin/*.exe` 及 d22scan／mockllm 台件。**三类都不是出厂 `wisp.exe`。**
- ⚠ 本腿未做真机 objdump（无 tracked exe，见 §6）；上面"GUI/CUI"判语＝**旗标文本层＋构建脚本调用链的静态读码**，出厂产物的现时 PE subsystem 读数仍缺、由票面 AC#1/AC#2 的落地腿带命令全文交。

## §3 subsystem 尺的有无与最小落点

**有没有一把"会响的尺"？—— 没有。仓里今天没有任何测试／脚本检查产物的 PE subsystem（GUI(2) vs CUI(3)）。**

尺（显式根 `cmd internal tools docs scripts .github docker`，ripgrep，`-i`）：`objdump|Subsystem|IMAGE_SUBSYSTEM|debug/pe|OptionalHeader|IMAGE_FILE_HEADER|00000002|windowsgui` 命中的全部落点＝
- `scripts/build.ps1:103/115`——**旗标本体＋ why 注释**（它让产物变成 GUI，但**不检验**产物是不是 GUI）；
- `cmd/wisp/console_windows.go:26`、`cmd/wisp/console_other.go:7`、`cmd/wisp/main.go:144`——**注释**（描述 windowsgui 语义，非断言）；
- `docs/specs/SPEC-11:50`、`docs/BUILD.md:87/90`、`docs/reports/**`、`docs/evidence/s1/244-black-console-r1.md`、票面——**规格/文档/证据件里的文字与手工 objdump 追述**，不是自动化尺；
- `tools/**`、`internal/**` 里 `Subsystem` 命中的都是**无关义的同名英文词**（"owner names the subsystem"、"isolated subsystems" 等，`config_reload.go:116`／`observe/goroutine.go:232/254` 一类），**零枚**读 PE 头。

现成构建校验面逐枚核过、**都不看 subsystem**：
- `scripts/build.ps1` step 6（`:166-169`）冒烟＝跑 `wisp.exe doctor` 看 rc，**不 objdump、不读 PE**；step 5（`:151-163`）只算 SHA256SUMS。
- `tools/d22scan`（AGENTS.md §1.2 那台代码形状仪器）的 ban 是"裸 go／filepath 越权／明文密钥／墙钟超时／镜像站哈希／面板侧 L2 允许／emoji"——**没有一条扫 `-H` 或 subsystem**。
- `scripts/slo-check.ps1`、`scripts/wisp-cli-tests.sh`、`scripts/portable-tests.sh`、`scripts/winsec-tests.sh`：grep subsystem/objdump **零命中**。
- ⇒ 与票面 §5 G5 一致："包内台件二进制天生 CUI、包内断言看不见 subsystem 这一维"——**且不止包内看不见，全仓根本没有一把读 subsystem 的尺**。

**要装一把这样的尺，最小落点在哪：** subsystem 是"链接后 PE OptionalHeader"的属性，尺只能对**已产出的 exe** 读。最小落点＝**在构建产出之后加一发 subsystem 读数**，两形任选（本腿只列形状、不选、不写）：
- **形甲**：`scripts/build.ps1` 的 step 3 之后／step 6 冒烟区（现 `:129` 产 exe、`:166-169` 跑 doctor）之间，读 `build\wisp.exe` 的 subsystem，非 GUI 则 `Fail`。
- **形乙**：一枚**外部构建后台件**（不在 `cmd/wisp` 包内，否则踩 G5 的"台件天生 CUI 看不见"），跑 build.ps1 产 exe → 读该 exe subsystem → 断言 == 2。
- 形状上应与票面 AC#5 ⓑ 那把"会响的尺"（`A548` 第 3 条已归口 `255-r2`/`258-r1` 落地腿）并列登记：**它是"零任务生产者会响"的尺，不是"subsystem 会响"的尺**——两码事，别并成一把。

**能不能在不构建的前提下测？—— 不能。** 三点具名：
1. 库内**零枚 tracked exe**（`git ls-files '*.exe'`＝0；`*.exe`＋`build/` 在 `.gitignore:2/14`）⇒ 没有可离线读的 PE 样本，`debug/pe`/objdump 都无米下锅。
2. subsystem 只在**链接后**才存在，源码/脚本层"看着带了 `-H`"≠产物真是 GUI（c1 §5 第 3 条与票面 §1 第 22 行都点过这个词面/产物差别）。
3. 本腿**一枚构建/Go 命令都不许跑**（第一硬规）⇒ 连造一个样本都做不到。

⛔ 不提议放宽任何既有断言；此处不写"做不到"式否证入账（本仓铁律：否证必须自己跑过作用面才许入账，本腿一枚都没跑）。⇒ **需要谁来量**＝能在本机跑 build.ps1＋对**新产出** exe 做 objdump 的**落地／验收腿**：票面 AC#1（改前 CUI(3)／改后 GUI(2) 两枚读数都要带命令全文入账）与 AC#2（父控制台／重定向／双击三形态真跑），且 subsystem 这把新尺的落地按 `A548` 归口那族（`255-r2`/`258-r1` 同写面）。

## §4 stdout 可见性〔读码推断，未跑〕

**问题**：exe 切成 GUI subsystem 后，`wisp run` 那两条往 `rt.stdout` 写的 stdout 句子还看得见吗？

**先定位这两条**（`cmd/wisp/config_reload.go`，全 `Fprintf(rt.stdout, …)`）：
- `:124-127` `startConfigReload` 的"配置热加载已接管（每 1s 检查一次…）"（`rt.auditf` 在 `:120` 先落 stderr/sink，stdout 这条是给人看的那半）。
- `:171`（立即生效段）、`:190`（放宽没生效）、`:195`（放宽已过 L2 写进内存）、`:201-203`（fs 要重启）、`:288-293`（`reportRestartPending` 重启档那句）。

**读码链（关键在时序，非旗标本身）：**
1. `main.go:90` 对 `run` 分支**先调 `attachParentConsole()`，再 `os.Exit(cmdRun(…))`**。
2. `attachParentConsole`（`console_windows.go:33-43`）：若 `GetStdHandle(STD_OUTPUT_HANDLE)` 已是可用句柄就 `:35-36` **直接 return**（重定向／父控制台已给的合法 stdout 走这枝）；否则 `:39` `AttachConsole(ATTACH_PARENT_PROCESS)`＋`:40-42` 三发 `rebindStdHandle`，把 `os.Stdout`/`os.Stderr`/`os.Stdin` **重新赋值成新开的 `CONOUT$`/`CONIN$`**（`:55-60`）。
3. `cmdRun`（`main.go:161`）把 **此刻的 `os.Stdout`** 读进 `runSpec.stdout`，`run.go:390` 再赋给 `rt.stdout`。**因为 attach 在 `:90` 早于 `cmdRun` 里这次读取，`rt.stdout` 拿到的是"重绑后/早返回后的那个可用 stdout"，不是切 GUI 之前的死句柄。**

**结论〔读码推断，未跑〕：在"有父控制台可 attach"与"stdout 被重定向"两种形态下，这两条 stdout 句子应当仍看得见**——判据是上面第 2、3 步的**赋值时序**（attach 先重绑 `os.Stdout`，`rt.stdout` 后捕获），与 c1 §3.2/§4 的判据链同向、本腿独立复认。唯一"看不见"的分支＝**attach 失败且无重定向**（如 `explorer` 无父控制台拉起 `wisp run`）：`AttachConsole` 失败"无害"（`:38` 注释），`rebindStdHandle` 的 `CreateFile("CONOUT$")` 在无控制台时也失败并 `:50-51` return ⇒ `os.Stdout` 停在 GUI 进程原始的无效句柄 ⇒ `Fprintf(rt.stdout,…)` 写进死柄、**无输出但无害**（退出码与任务照走完，c1 §4）。

**只有哪一发真机读数能定它**：票面 **AC#2 ①**（在父控制台里 `wisp run "…"`，看这两条 `wisp run: 配置…` 句子真出现在那个终端）与 **②**（`wisp doctor > out.txt` 重定向仍有效）。⚠ 本腿零构建、不能造 GUI 样本，也无法验证 Win32 对"GUI 子系统进程被 cmd 拉起时 `GetStdHandle` 返回可用还是无效"的语义（这正是 c1 §5 第 3 条自陈的静态盲点）——**该语义一旦与推演不符，判错处就在上面第 2 步的 return/attach 分支选择**。

## §5 我可能写错的条目（自我对抗）

（取数中）

## §6 量不到的地方（具名）

（取数中）

## §7 交件判语

（取数中）
