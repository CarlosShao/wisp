# 244-a2 — GUI 子系统这一刀的重派普查腿（只读；前腿 `244-a1` 只留空骨架，本腿从零跑尺）

- 腿：`244-a2`（只读普查；⛔ 零产码／零脚本／零测试／零构建／零 `go.mod` 改动；⛔ 没跑过任何 `go build`／`go vet`／`go test`／`scripts/build.ps1`）
- 票面：`.scratch/wisp/issues/244-spec-11-wants-the-gui-subsystem-build.md`（43 行，AC#0–AC#5；**AC#5＝10-01 09:4x 新补，本腿首要问题**）
- 起手锚点：HEAD＝`f1e7a3e2cbfc3f327c9fa64771d2e259351900d8`，同发 `date` 逐字＝`2026-10-01 09:43:51 +0800`（S1）
- 收尾锚点：HEAD＝`7a41db9b5630c8a9d11aee25070fec64f3aea9ed`，`date`＝`2026-10-01 09:50:26 +0800` ⇒ **期间落地一枚提交（`7a41db9b`＝台账 A489＋`internal/panel/composer_dispatch_test.go` 的负向钉反转），我读过的 21 枚文件逐枚重取 blob，全部与起手同号**（表见 §⑥ 末把 S21）⇒ 本件的读数没有踩在 33-r1 的写点上；凡它后续改动，本件一律待复认。
- 前腿残值（具名，不覆盖不删除不搬动）：`.scratch/wisp/probes/244/a1/census.md` 本腿起手现读＝57 行（`wc -l` 我现跑），五节正文逐节"填写中。"，与 `A488` 那句"停在 57 行／1,548 字节"一致 ⇒ **前腿没跑过尺，本腿不受"不许重跑前人读数"约束**。
- ⛔ 本腿不碰 `frontend/**`／`design/**`（零读零写零转述）；未改 `docs/PLAN.md`／`docs/specs/**`／`docs/BUILD.md`／`docs/SLO.md`／`internal/observe/thresholds.go`／golden／`allowlist.txt`／三枚冻结件一字；工作树脏项（起手 `git status --porcelain`＝**204 行**）一律没动。
- ⛔ 票面六格 AC 框一枚没碰（本件是普查表，不是交件）。

> **写作顺序声明（票面纪律）**：本件按硬要求先写满 §⑥§⑦§⑧（尺与自对抗），再写 §①–§⑤ 的结论。

---

## ⑥ 我跑了哪些尺、每条真实读数

**记法**：`〔量〕`＝本腿真在这台机器上跑的命令的读数；`〔推〕`＝读码推出的、本机没量的（下一节⑦里逐条列它的风险）。每条被读过的易变文件都带起手 blob 号（`git rev-parse HEAD:<path>`，取数时刻 `2026-10-01 09:43–09:45 +0800`，HEAD `f1e7a3e2`）。

| # | 尺（命令全文） | 真实读数 |
|---|---|---|
| S1 | `date "+%Y-%m-%d %H:%M:%S %z"`；`git rev-parse HEAD`；`git rev-parse --abbrev-ref HEAD` | `2026-10-01 09:43:51 +0800`／`f1e7a3e2cbfc3f327c9fa64771d2e259351900d8`／`dev` |
| S2 | `git status --porcelain \| wc -l`；`git status --porcelain -- go.mod go.sum cmd internal` | **204 行**；**0 行**（起手这一刻 `cmd/`＋`internal/`＋`go.mod`＋`go.sum` 全是干净的；33-r1 还没落笔） |
| S3 | `objdump -p build/wisp.exe \| grep -i subsystem` | `Subsystem 00000003 (Windows CUI)`（＋`MajorSubsystemVersion 10`）⇒ **票面/A488 那把尺我复认成立，不是新发现** |
| S4 | `ls -l --time-style=long-iso build/wisp.exe` | `30299552` 字节、mtime `2026-09-30 23:43` ⇒ ⚠ **与票面现量表那枚"29,771,385 字节／09-30 11:21"不是同一枚产物**：产物在我这腿之前被重跑过一次（谁跑的不在本腿射程），**子系统维度没变**（仍 CUI）。`build/**` 被 `.gitignore` ⇒ 不可再生，下一位必须自带命令全文 |
| S5 | `grep -rn "submitTask" --include=*.go . \| grep -v .scratch` | 产码只有**一枚文件**：`cmd/wisp/resident_task_source_windows.go`（blob `018152a82829fb72e3a30fdb7b0eeca4b4eb4b2d`）＝定义 `:401` ＋ **两枚调用者 `:322`（注入分支）／`:363`（控制台 `task <文本>` 动词）**。⚠ 同尺**不过滤** `.scratch` 时会多命中 `./.scratch/wisp/probes/246/r2/mut/base-resident_task_source_windows.go` 的同号三行＝**那是别人的突变快照副本，不是产码**；`--include=*.go` 不加 `grep -v` 会把这枚副本算成第三个调用者 ⇒ **调用者枚数的正数＝2，且必须带 `grep -v .scratch`** |
| S6 | `grep -rn "interactiveStdin" --include=*.go . \| grep -v .scratch` | 产码调用者**两枚**：`cmd/wisp/main.go:153`（`cmdRun` 的答复流，blob `e5e98ea89d49810dafb6baa41d5c419ca6cfaaa0`）＋`cmd/wisp/resident_task_source_windows.go:218`（常驻任务源的**唯一**闸门，ruling 2.1 禁止第二把）。定义：`cmd/wisp/approval_reply_stdin_windows.go:41`（blob `a99ac3a93acf0c70ede5ad82eff72f836a1b8982`）＝`GetStdHandle(STD_INPUT_HANDLE)`→`GetConsoleMode` 失败即 nil；POSIX 侧 `cmd/wisp/approval_reply_stdin_other.go:23` 恒返 nil |
| S7 | `grep -rn "attachParentConsole" --include=*.go cmd/ internal/ \| grep -v _test.go` | 调用者**十二枚**：`cmd/wisp/main.go:65`（无参 GUI 路）`:90`(run) `:93`(providers) `:96`(doctor) `:108`(models) `:113`(panel-assets) `:119`(panel-inbound) `:122`(version) `:125`(help) `:128`(default 拒绝路) ＋ **`cmd/wisp/secret.go:183`** ＋ **`cmd/wisp/slo_windows.go:254`**。定义 `cmd/wisp/console_windows.go:33`（blob `fbb33baffbc423515f71875570c95aa0b4270766`）：开头 `if err == nil && out != 0 && out != InvalidHandle { return }` ⇒ 有合法 stdout 就**跳过 attach**；`rebindStdHandle`（`:47`）重绑 `CONOUT$`/`CONIN$`。no-op 对岸 `cmd/wisp/console_other.go:16` |
| S8 | `grep -rln "attachParentConsole\|AttachConsole" --include=*_test.go cmd/ internal/ tools/` | **零命中**（空输出）⇒ 今天**没有任何测试／仪器碰过 attach 通路**，票面 AC#0 那句"真机凭据只有 `docs/BUILD.md:90` 那次临时构建"在本仓仪器面上成立：连一枚读码钉子都没有 |
| S9 | `grep -rn "\.Allow(\|AllowSession(\|PanelAllow(" --include=*.go cmd/ internal/ tools/ \| grep -v _test.go` | 能"允许"的产码入口**只有三枚，且同一枚文件**：`cmd/wisp/approval_reply.go:215`（`s.live.h.Allow`）／`:259`（`AllowSession`）／`:331`（`PanelAllow`，**它不是允许，见 S12**）。blob `33e82acf7069f81aef92795b05cc180e601b5d2b` ⇒ **允许这件事在产码里只有一个漏斗：`replySurface`** |
| S10 | `grep -rn "surface.handle\|\.handle(verb" --include=*.go . \| grep -v .scratch` ＋ `grep -rn "runReplyLoop\|attachReplyListener" --include=*.go . \| grep -v .scratch` | `replySurface.handle` 的产码调用者**两枚**：`cmd/wisp/approval_reply.go:546`（`runReplyLoop`，由 `attachReplyListener` 起，`attachReplyListener` 的唯一产码调用者＝`cmd/wisp/run.go:770`，喂的是 `runSpec.reply`＝`main.go:153` 那枚 `interactiveStdin()`）＋`cmd/wisp/resident_task_source_windows.go:381`（`runConsoleLoop`，同一个 `interactiveStdin()` 闸门下 `:287-290` 才建 surface）。⇒ **控制台键盘是"允许"的唯一产码驱动源，两枚都挂在同一根线上** |
| S11 | `internal/agent/approval/ui.go:139-171`（blob `55bcd1daf251a9ef513ef43b3c63605787252b13`）现读 | `NativeAPI`＝`Allow(ctx,corr,grant)`／`AllowSession`／`Reject`；`PanelAPI`＝`Reject`／`Head`／`View`——**面板侧没有 Allow 方法**（`:164-167` 逐字"It has no Allow method… it is a method that does not exist"）；`ErrPanelAllow`＝`:132-134`「面板来源不得允许（F2 第三层：允许只接受原生侧）」 |
| S12 | `sed -n '726,748p' internal/agent/approval/gate.go`（blob `31e36a849adb6388e8f2bfadd5067582ae72b881`）＋`sed -n '425,455p' internal/agent/approval/replies.go`（blob `190803fdaaa1a0ba846ad6710589665c6b7b83be`） | `DecideFromPanel`（gate.go:730）对 `r.Allow` **按路线直接拒**：打 `PANEL-ALLOW-REJECTED`、`Grant != ""` 时 `q.revokeGrants(corr)` 烧掉那枚卡的新增号，返回 `ErrPanelAllow`；`Replies.PanelAllow`（replies.go:431）**err==nil 才算安全故障**（`:442-444` 逐字"a nil error out of here is a security fault"）。⇒ **`panel-yes` 那枚词（`cmd/wisp/approval_reply.go:574`）是"试一下然后被拒"的仪器形状，不是面板给的允许入口** |
| S13 | `grep -rn "jchv/go-webview2" --include=*.go . \| grep -v .scratch` ＋ `grep -n "webview" go.mod`（blob `6ccb3fd16629e239e03ac36be2caa6169544ad18`） | 产码**零枚**导入：命中只有 `internal/panel/composer_dispatch_test.go:458-459`（测试里的假 go.mod 字符串）＋`scripts/spike/webview2-latency/main.go:31`（独立 module，`scripts/spike/go.mod`）。`go.mod:19`＝`github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808 // indirect` ⇒ **甲（面板路）今天的产码执行者＝0**；`33-r1` 正在把它变成 direct，这一格在我读完之后一定变（标〔锚点期间被写过，待复认〕归它） |
| S14 | `grep -rn -- "-H=" scripts/ tools/ .github/`；`grep -rn "windowsgui" scripts/ tools/ .github/ cmd/ internal/ \| grep -v .scratch` | 第一把：唯一命中＝二进制 blob `scripts/spike/bin/goja-caps.exe`（**不是构建链**）。第二把：全仓 `windowsgui` **只剩 1 行，且是注释**＝`cmd/wisp/console_windows.go:26`；⚠ **票面现量表说"两行注释"，我这一把只抓到一行**，另一行在 `cmd/wisp/console_other.go:7`，它写的是英文「links as a GUI-subsystem binary」，**不含 `windowsgui` 词面** ⇒ 票面那把尺的词面没抓到它，两条都在 §⑤ 的 AC#3 名单里逐枚点名（不是推翻票面，是补它的词面盲区） |
| S15 | `grep -n "go build\|ldflags\|-X " scripts/build.ps1`（blob `3593e87b438a4e2d1673fa0fa7af8285429ed57e`）＋现读 `:103-110`／`:123`／`:162-163` | `$ldflags` 六枚**全是 `-X`**（Version/Commit/BuildDate/DefaultEnv/SherpaOnnxVersion/OnnxRuntimeVersion）、**零枚 `-H`**；唯一构建行 `:123`＝`& $go.Source build -trimpath -ldflags $ldflags -o build\wisp.exe ./cmd/wisp`；`:162` 冒烟＝`& build\wisp.exe doctor`，`:163` 只判 `if ($LASTEXITCODE -ne 0)` ⇒ **冒烟不读 stdout、不读子系统**（§④ 的"谁能看见"里这是个零） |
| S16 | `grep -n "Build wisp.exe\|build.ps1" .github/workflows/ci.yml`＋现读 `:453`／`:550`／`:613`；`grep -n "build.ps1\|wisp.exe" .github/workflows/slo-fresh.yml` | `ci.yml` 三枚步骤调 `scripts/build.ps1 -Env dev`：`:448`「cgo build smoke」→`:453`、`:549`「Build wisp.exe」→`:550`、`:612`「Build wisp.exe (deps cached)」→`:613`（job＝`test-windows:387`／`slo-smoke:532`／`slo-full:590`）。`slo-fresh.yml`＝**零命中**（它不产 `wisp.exe`）⇒ **构建入口名册里 CI 那三枚都只是 build.ps1 的调用者，不是第二处 flag 面** |
| S17 | `grep -rn "wisp\.exe" scripts/*.ps1 scripts/*.sh .github/workflows/*.yml` | 产物消费方：`scripts/build.ps1:146`（`SHA256SUMS` 名单含 `wisp.exe`）／`:162`（doctor 冒烟）；`scripts/slo-check.ps1:69`（默认 `-WispExe build\wisp.exe`）／`:104`（**缺就自己调 build.ps1 -Env dev**）／`:108`／`:155`（进程名清单里认 `wisp.exe`／`wisp-cli.exe`）／`:325`／`:344`／`:365`（三发 `& $WispExe slo … -out <file>`）；`ci.yml:549/612` 见 S16 |
| S18 | `grep -n "go build\|go test\|-ldflags" scripts/wisp-cli-tests.sh scripts/portable-tests.sh scripts/winsec-tests.sh` | 四枚 shell 步骤**全是 `go test` 家族**，**零枚 `-ldflags`／零枚 `-H`** ⇒ 它们产的测试二进制与 `wisp.exe` 的子系统无关；`scripts/wisp-cli-tests.sh` 只是 PATH 装 DLL 的壳（`:111-112` 交回 `portable-tests.sh --scope=cli`） |
| S19 | `go version`；`head -8 go.mod` | `go version go1.27.1 windows/amd64`；`go.mod`＝`module github.com/CarlosShao/wisp`／`go 1.27`／`toolchain go1.27.1`；起手 `git status` 干净（S2） |
| S20 | `for f in $(git ls-files cmd/wisp); do head -1 $f \| grep -o "go:build.*" \|\| echo NONE; done \| sort \| uniq -c` | `cmd/wisp` 71 枚 tracked 文件：**21 枚 `//go:build windows`／7 枚 `//go:build !windows`／3 枚 `//go:build windows && winlive`／40 枚无标签**（含 `main.go`／`run.go`／`approval_reply.go`／`panel_inbound.go`／`doctor.go`／`secret.go`／`models.go`）；三枚 winlive 产码无关文件＝`resident_approval_live_246_windows_test.go`／`resident_ball_live_228_windows_test.go`／`resident_task_source_live_246_windows_test.go`（`internal/ball` 另有 4 枚）⇒ 真机那一族要 `-tags winlive`，而 `grep -rn winlive scripts/ .github/ docs/BUILD.md`＝**零命中**（`A488` K10／`Q-75` 第四次具名，本腿**不重复立案**，只在 §④ 记"切换后这族是否变红"取决于它跑不跑） |
| S21 | 21 枚被我读过的文件逐枚 `git rev-parse HEAD:<path>`（起手 `f1e7a3e2`）＋收尾重取（`7a41db9b`） | **起手＝收尾全部同号**（逐枚：`a99ac3a9`／`fbb33baf`／`7a7f8ea0`／`018152a8`／`af228d67`／`e5e98ea8`／`cff7244f`／`33e82acf`／`d790a7e1`／`357f7f27`／`c8bdccc0`／`d1f1b881`／`4ae16b88`／`545c729e`／`b34ccf13`／`fa59c301`／`31e36a84`／`190803fd`／`55bcd1da`／`3593e87b`／`6ccb3fd1`）⇒ **本腿读过的文件在锚点期间没被写过**，不必标〔待复认〕；但 `7a41db9b` 动了 `internal/panel/composer_dispatch_test.go`（A489 那句"负向钉反转"），**那枚文件我没读、本件也不依赖它** |
| S22 | `grep -rn "\"build\", \"-o\"\|exec.Command(goBin\|go build" --include=*_test.go cmd/ internal/ tools/` | 自己会 `go build` 的仪器**四枚**：`cmd/wisp/secret_argv_windows_test.go:174`＝`exec.Command(goExe,"build","-o",exe,"./cmd/wisp")`（**不带 -ldflags**）／`cmd/wisp/resident_approval_live_246_windows_test.go:417`（build `./cmd/wisp/testdata/esclistener`）／`internal/llm/openaichat/mockllm_integ_test.go:65`／`tools/d22scan/{scan_test.go:2256,selftest_test.go:181}`。⇒ **AC#1 给 build.ps1 加 `-H` 不会穿过这些行**：它们在测试临时目录里重新编一枚**自己的 CUI `wisp.exe`**，与出厂产物的子系统无关——这正是 AC#5 ⓑ 要的"仪器不会因此变红"的第二条独立证据 |
| S23 | `sed -n '45,95p' docs/BUILD.md`（只读，没改一字） | `:52-53`「ldflags 注入版本/commit/构建时间/`DefaultEnv`/两个原生库版本」（无 `-H`）；`:57`「冒烟：运行 `build\wisp.exe doctor`，FAIL 则整个脚本 FAIL」；`:67-71` 那段"无参＝常驻进程"的读法仍写着"悬浮球窗口在票 07""`run` 是 CLI 占位"（**过期文案，禁改，只在 §② 具名登记**）；`:87-88`＝票面点名的"推迟到票 07"；`:90-92`＝临时构建验证 AttachConsole 那三行，逐字「`wisp.exe version` 在调用方控制台正常输出…**且** `wisp.exe version > out.txt` 重定向不受影响（已有合法 stdout 时跳过 attach）」⇒ **这句"跳过 attach"与 S7 读到的 `console_windows.go:35` 那个 early return 是同一条机制**，是本腿 §③ 重定向那一形的唯一旧凭据，且它自己声明是"当时" |
| S24 | `sed -n '60,200p' cmd/wisp/resident_windows.go`（blob `357f7f276f6a2834ea544023c0e429b7b50751b8`）＋`grep -rn "signal.Notify\|os.Interrupt" internal/proc/ cmd/wisp/` | 常驻腿的停机来源**只有一形**：`cmd/wisp/resident_windows.go:107`＝`signal.Notify(bootExit, os.Interrupt, syscall.SIGTERM)`，事件循环 `:200-202` 打印「(Ctrl+C exits cleanly)」并进 `rt.RunEventLoop()`；`internal/proc` 侧同形＝`internal/proc/boot_windows.go:128 signal.NotifyContext(..., os.Interrupt, SIGTERM)`。日志落盘另有 `cmd/wisp/resident_windows.go:63 installLogSink(rt.Layout.DataDir)` ⇒ **slog 走文件、fmt.Printf 走 stdout**，两件事在 GUI 下命运不同（§③） |
| S25 | `sed -n '120,135p;240,258p' cmd/wisp/resident_ball_windows.go`（blob `cff7244f3125f8fd499ef637e5871b6dd10301d1`）＋`sed -n '669,686p' internal/ball/ball_windows.go`（blob `b34ccf1347383e62db67919b0c7837ea859298b5`）＋现读 `internal/ball/tray_windows.go`（blob `fa59c30125899f0ff2af370d3ee099b5b83ce601`） | 托盘菜单构造点＝`internal/ball/tray_windows.go:86-91`（四枚 `appendItem`：打开面板/静音/暂停唤醒/退出，常量 `:20-23`），分流点＝`internal/ball/ball_windows.go:674-684`（`showMenu` 返回值→`b.fire(Events.OnTray*)`，左键 `:671-672` 也走 OnTrayPanel）。**常驻腿的执行者**＝`cmd/wisp/resident_ball_windows.go:128-131`：四枚里三枚是 `recordBallGesture("tray-*")`，第四枚是 `recordTrayExit`；`recordBallGesture`（`:211-214`）只做 `slog.Warn`＋`fmt.Printf`，`recordTrayExit`（`:253-257`）逐字「the tray Exit item has no stop path attached to it」⇒ **托盘四枚按钮在常驻腿里的"真执行者"数＝0**（同尺在 `cmd/balldebug/main.go:199-206` 里有四枚真执行者，那是调试壳，不是出厂路） |
| S26 | `grep -rn "OnTray\|OnSummonHotkey\|OnCancelHotkey\|OnPanelHotkey\|OnMuteHotkey" --include=*.go cmd/wisp internal/ball \| grep -v _test.go` | 热键四枚：`cmd/wisp/resident_ball_windows.go:124-127`＝summon/mute/panel 三枚 `recordBallGesture`，**只有 `:126 OnCancelHotkey: recordCancelHotkey(onCancelEsc)` 有执行者**；`recordCancelHotkey`（`:224-234`）调注入的 `escVetoFunc`，实参＝`cmd/wisp/resident_windows.go:136` 的 `ra.vetoByEsc`，落点 `cmd/wisp/resident_approval_windows.go:171`（`:179 ra.cards.Veto(...)`）⇒ **这是"否决"，不是"允许"**；`grep -n` 同一枚文件里没有 `Allow` 调用（S9 全仓尺已证） |

---

## ⑦ 我可能写错的条目（对抗我自己）

| # | 我差点写下的话 | 判决 | 推翻／复认用的尺 |
|---|---|---|---|
| 7.1 | 「`main.go` 里 `secret` 与 `slo` 两枚分支没调 `attachParentConsole` ⇒ 这两条腿在 GUI 下必盲」 | **我错了，当场推翻**。我第一轮只读了 `cmd/wisp/main.go:100-111` 的分支体就下结论；真身是**这两枚分支在自己文件里调**：`cmd/wisp/secret.go:183`、`cmd/wisp/slo_windows.go:254`（S7）。⇒ 正确说法：**十二枚调用者覆盖全部腿，没有漏的那一枚**，本件 §③ 按这个写 | S7（对 `cmd/`＋`internal/` 全量 grep，不限 `main.go`） |
| 7.2 | 「`build/wisp.exe` 就是票面现量表那枚产物」 | **错**：票面写 29,771,385 字节／09-30 11:21，我读到 30,299,552／09-30 23:43（S4）⇒ **产物已被换过**，我只复认了"子系统＝CUI 这一维没变"，其余（大小、内含 commit）不能与票面同项比。`build/**` 不可再生 ⇒ 我在 §④ 把这条列成"今天的产物读不到今天的构建链"的第二个实例 | S3＋S4 |
| 7.3 | 「切 GUI 之后任务源**绝对**没了」 | **过头，要带条件**。`submitTask` 的 `:322` 注入分支**不看控制台**，看的是 `residentTestTaskText`（`cmd/wisp/resident_task_source_windows.go:178-197`）三把锁：`env==test`＋`WISP_TEST_DATA_DIR` 在场＋数据根逐字相等。⇒ 双击那条路（`WISP_ENV` 默认从构建来，`-X ...DefaultEnv=prod`，见 `scripts/build.ps1:107`）里第一把锁就不过 ⇒ **对 prod 双击确实归零**；但**"永远归零"是谎**，test 路仍能举卡——这正是票面 AC#5 那句"票 246 那族不会因此变红"的机制解释 | S5＋S6＋S15（`-X DefaultEnv`） |
| 7.4 | 「GUI 进程没有父控制台 ⇒ attach 必失败 ⇒ 所有 stdout 消失」 | **前半是平台常识、不是我这台机器的读数**；我只能证到代码形状：`console_windows.go:33-43` 只 `AttachConsole(ATTACH_PARENT_PROCESS)`，且 `:35` 的 early return。⚠ **反例形状我处理不了**：如果 GUI 子系统进程是从 cmd.exe 里手打拉起的，Windows 可能把标准句柄**继承**进来，那时 `out != 0` ⇒ 走 early return、根本不 attach，输出反而看得见（`docs/BUILD.md:92` 那句"重定向不受影响"就是这一形的旧读数）。**这条我本机没量、票面 AC#2 ①② 要的就是它** ⇒ 全件里凡涉及"GUI 下 stdout 可不可见"的句子一律标〔推〕，不写成事实 | S7＋S23（`docs/BUILD.md:90-92`） |
| 7.5 | 「`panel-yes` 是面板给的允许入口」 | **错，且是安全方向上的错**。`handle` 的 `"panel-yes"`（`cmd/wisp/approval_reply.go:574`）走的仍是**控制台**词法，`panelAllow`→`Replies.PanelAllow`→`DecideFromPanel` **按路线拒**并烧新增号（S12），`err==nil` 才算故障。⇒ 它是"试一下就拒"的仪器，**不是一个入口**；本件 §① 把它算进"允许入口"会是错的，我只在"形状名单"里列它并标"不可用（设计上）" | S10＋S12 |
| 7.6 | 「托盘／热键没有执行者」 | **要限定归属**。真执行者在 `cmd/balldebug/main.go:199-206` 是**有的**（四枚托盘＋四枚热键全接），出厂常驻腿 `cmd/wisp/resident_ball_windows.go:124-131` 里**只有 cancel 一枚**（S25/S26）。⇒ 本件写"import 了≠装起来了"那一族的正确版本：**能弹菜单的代码在 `internal/ball`，会办事的代码在 `cmd/wisp` 里只装了一枚否决** | S25＋S26 |
| 7.7 | 「票面说两行注释、我只抓到一行 ⇒ 票面那把尺坏了」 | **不算推翻**：票面 S14 那一把是 `grep -rn windowsgui`，`cmd/wisp/console_other.go:7` 的词面是「links as a GUI-subsystem binary」，**不含 `windowsgui`** ⇒ 词面盲区，两枚文件我都现读到原句（§⑤ 逐枚列行号）。登记它是给票面补盲区，不是说它错 | S14＋`console_other.go:7` 现读 |
| 7.8 | 「winlive 那族仪器切完会红」 | **判不了，别写成结论**。三枚 `cmd/wisp` 的 live 用例（S20）带 `//go:build windows && winlive`，而 `grep -rn winlive scripts/ .github/`＝零命中 ⇒ **整包里它们今天根本不参跑**；它们红不红取决于下一程跑不跑 `-tags winlive`（`A488` 已把这格裁成 `Q-75` 三栏，默认"不做"）。本件只写"看不见"，不写"会红" | S20 |
| 7.9 | 「attach 通路一枚仪器都没有」 | **复认，但尺的形状要说全**：我那把是 `grep --include=*_test.go`（S8 空）。⚠ 同名不带 `_test.go` 后缀的台件（`cmd/wisp/testdata/esclistener/main.go`）不在射程；我另外现读了它 `:60` 只声明 `ShowWindow`、**没有 `AttachConsole`/`FreeConsole`**（`grep -rn "FreeConsole\|AllocConsole\|GetConsoleWindow" cmd/ internal/ tools/`＝**零命中**）⇒ 结论稳，但"零仪器"这句话的射程是"名字里带 attach/console 的测试"，不是"任何测过标准句柄的测试" | S8＋FreeConsole 全仓 grep（零命中） |
| 7.10 | 「我读的文件都是稳的」 | **成立但要带时刻**：起手 HEAD `f1e7a3e2`、收尾 HEAD `7a41db9b`（期间落地 1 枚提交），21 枚文件 blob 逐枚同号（S21）。⇒ 本件读数是稳的；**但 `33-r1` 正在写 `cmd/wisp`＋`internal/panel`，它任何一笔提交都会让 S13（webview2 导入数＝0）作废**，那一格我在 §① 明确标"到 09:50:26 为止" | S21 |
| 7.11 | 「`wisp doctor` 的冒烟在 GUI 下会失效」 | **错**：`scripts/build.ps1:163` 判的是 `$LASTEXITCODE`，不是 stdout ⇒ 冒烟**照样绿**（`cmdDoctor()` 返 bool，`main.go:97-99` 才 `os.Exit(1)`）。真正变的是**人看不看得见那 10 行 PASS**（`cmd/wisp/doctor.go:115-129` 全是 `fmt.Println`/`Printf`＝stdout）。⇒ §④ 写"仪器面＝零"时指的是**没有一枚尺断言过 doctor 的输出可达**，不是说冒烟会红 | S15＋`cmd/wisp/doctor.go` 现读 |
| 7.12 | 「SLO 那条腿会因切 GUI 取不到数」 | **多半不会**：`scripts/slo-check.ps1:325/344/365` 三发都带 `-out <file>`，报告落文件；`wisp slo` 的错误行是 `os.Stderr`（`cmd/wisp/slo_windows.go:221-333`）。⇒ 风险面只剩"stderr 在 GUI 下无处可去"那一形，且 `slo-check.ps1:104` 自己会重跑 build.ps1（**它会跟着产出同一种子系统的 exe**）——这条我只读码，没跑 | S17＋S15 |

---

## ⑧ 判不动的地方（逐条甲／乙／不做＋现量）

> ⛔ 本节**不替编排者做选择**。票面 AC#5 ⓐ 逐字"不许由实现腿自己选"，所以候选只列代价；碰了票面没覆盖的形状就落在这里。

| # | 判不了的事 | 现量（为什么判不了） | 三条出路（甲／乙／不做） |
|---|---|---|---|
| G1 | **GUI 进程的任务源从哪来**（票面 AC#5 ⓐ 原问） | 三条候选的执行者枚数我全量了：甲＝面板——产码 `webview2` 导入 **0**、`go.mod:19` 仍 `// indirect`（S13）；乙＝托盘四枚——常驻腿**真执行者 0**（S25）；丙＝控制台——切完 `interactiveStdin()` 归 nil（S6）。⇒ **"哪条能落地"是排程判断，不是读数判断** | 甲＝按票面既定顺序等 33→248 长出能举任务的面板路；乙＝给托盘加"起一发任务"（要先解决"文本从哪来"，非控制台文本源的合法性票面没写）；**不做＝不许切**（票面 AC#5 ⓑ 的闸门本来就是这个形状）。⛔ 我不选 |
| G2 | **切完双击进程"用什么方式停"**（本腿挖出的新格，票面未覆盖） | 停机来源只有 `cmd/wisp/resident_windows.go:107` 的 `os.Interrupt`/`SIGTERM`；托盘退出项无执行者（`cmd/wisp/resident_ball_windows.go:253-257`）；GUI 子系统＋Explorer 拉起 ⇒ 没有控制台键盘 = **没有 Ctrl+C**。⇒ 若 AC#1 先切、票 228 AC#2 的后端还没落地，**双击起来的常驻进程只能靠任务管理器杀**，D38(e) 十步一步都走不到 | 甲＝把它并进票 228 AC#2／票 248（谁给托盘 Exit 执行者谁负责）；乙＝本票 AC#5 ⓐ 的答案里**并列要求**"任务源＋停机源两条都得有"（现量：切完两条都归零）；**不做＝不行**——这条不是风格问题，是"切了以后关不掉"。⛔ 交编排者裁，票面没写 |
| G3 | **继承句柄那一形**（GUI 子进程从 cmd.exe 拿到不算"无控制台"的输出通路） | `console_windows.go:35` 的 early return＋`docs/BUILD.md:90-92` 的旧读数（它自陈是"当时"、"临时构建"），两者指向同一机制，但**我这台机器上没有任何一把尺能分辨"手打拉起"与 `explorer` 拉起**（那要真跑，本腿禁跑） | 甲＝按票面 AC#2 ①②③ 三形真机各一发（写码腿做）；乙＝把"哪一形走 early return"写成文档待复认；**不做＝把 `BUILD.md:90` 当凭据**（票面 AC#0 已明说那是过期读数，不可继承）⇒ 本件一律标〔推〕 |
| G4 | **"丙＝永远不切"要不要成立**（SPEC-11:50 的字面要求） | `docs/specs/SPEC-11-*.md:50` 逐字（票面引，我没改它一字）"无参 = GUI（`-H=windowsgui`）"——**点名了 flag**；仓内另有第三条去黑框的路我现量过它不存在：`grep -rn "FreeConsole\|AllocConsole" cmd/ internal/ tools/`＝**零命中** ⇒ "保持 CUI、启动后 FreeConsole 掉控制台"是**新产码**，不是现有能力 | 甲＝现在就切（AC#5 ⓐ 未答 ⇒ 票面自己禁止）；乙＝等面板能举任务再切（票面既定序）；丙（票面三形里那枚"永远不切"）＝**只要它含 FreeConsole 之外的形状，就与 SPEC-11:50 的字面 flag 冲突**，本件按票面要求写"与规格冲突"，⛔ 不自创替代形状。是否要给"运行时 FreeConsole"开一格合法性＝交编排者 |
| G5 | **AC#5 ⓑ 那枚"会响的尺"该钉在哪一层** | 现量两形都看不见这一刀：① 全仓零测试碰 attach（S8）；② 自己 `go build ./cmd/wisp` 的四台件都不带 `-ldflags`（S22）⇒ **AC#1 加 `-H` 之后，`go test ./cmd/wisp` 整包会照绿**。票面只给了形状（"能力型正向依赖边＋正控＝把调用者中和成 `var src *residentTaskSource`"），**没给钉该落在 CI 步、测试包还是交件里的读数**；而 `docs/evidence/s1/` 的裁决表归属（谁写哪一格）不在本腿地界 | 甲＝写码腿把钉子做成 `cmd/wisp` 包内的能力型断言（不跑构建）；乙＝做成报表尺（像票 243 AC#3 那族"抽样尺不能当门"）；**不做＝只看构建绿**（票面明令禁止的形状） |
| G6 | **过期文档行要不要在本票处理** | 现量三处过期文案：`docs/BUILD.md:51`（"S5 之前无 frontend"）、`:67-71`（"悬浮球窗口在票 07"、"`run` 是 CLI 占位（真实 agent 循环在票 10）"）、`:87`（票面已点名的"推迟到票 07"）。⇒ `:67-71` 那一段和 `:51` **票面没提**（票面只处理 `:87`）；且 `docs/BUILD.md` 整枚在禁改清单 | 甲＝票 244 交件里逐枚具名登记（票面 AC#4 就是这一形）；乙＝并入票 225 那一族的"文档里写的推迟"扫描尺（AC#4 已交回 225）；**不做＝改 `docs/BUILD.md` 一字**（禁区）⇒ 本件只做登记（§②），没动文件 |
| G7 | **票面 AC#5 标题那句"永久之家"要不要用我的读数复述** | 我没跑票 246 的任何用例（本腿禁跑），`A488` 里 246-v2 那句"AC#7 成立、附两条具名条件"是**编排者的裁定**，不是我量到的 | 甲＝本件只引 AC#5 的三把现量凭据＋我自己 S5/S6 的调用者枚数；乙＝重跑票 246 那族用例（⛔ 越界，本腿零测试）；**不做＝把"任务管线已跑通"写进任何一格**（`A488` 逐字禁止不带栏头写这句）⇒ 按甲写，栏头〔接缝注入〕照带 |

---

## ① GUI 进程里"谁举任务"与"谁能答复卡片"分别还剩什么

填写中。

---

## ② 构建入口名册：今天有几处会产出 `wisp.exe`，切 GUI 要动哪几处

填写中。

---

## ③ 切完之后四条 CLI 腿各自还能不能工作

填写中。

---

## ④ 三形时机（甲＝现在切／乙＝等 33+248／丙＝永远不切）

填写中。

---

## ⑤ AC#3：今天哪几行注释／文案是"切换没发生就是假话"的绝对句

填写中。
