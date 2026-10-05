# 244-a4 普查 —— 首启回执文案的「stderr 真到达」通道（只读普查，不产码）

> 腿号 `244-a4`。工作目录 `D:\work\workspace\projects plans\Wisp`。
> 起手锚点 `cfd97636`（`git log -1`，2026-10-05 12:58 +08:00）；本文件起手时刻 2026-10-05 13:01 +0800。
> 硬纪律：**一枚 `go` 命令都没跑**（`go test`／`go build`／`go vet`／`go env` 全部为 0 次调用）；
> 写点只有本文件；`grep`／`find` 的根一律显式限定为 `cmd internal tools scripts docs .scratch`。
> 本文件是**骨架 + 逐节读数**：下面是 §0–§7 的骨架，每节的正文由后续 Edit 逐节补齐并逐节 commit。

---

## §0 起手复认（派单给的每一条 `file:line` 断言的现读结果）

判语：**七条锚里四条例句成立、三条要更正**——`console_windows.go:38-51` 那句"没有父控制台就 return"**与本文件的真身方向相反**（详见 §4）；`A619` 那句"AC#1 不勾＝缺 AC#2 真机两发"**不在 A619**（`AC#1 不勾` 四字全台账 0 命中）；`objdump` 本机**有**，所以 §3 那一问不欠"量不到"。

| 派单断言 | 现读（本腿自己跑／读，13:0x–13:1x） | 判语 |
|---|---|---|
| `cmd/wisp/main.go:159-163` 附近把 `os.Stdout`／`os.Stderr` 交给 `runTextTask` | `cmd/wisp/main.go:159-164` 是那一枚 `runSpec` 字面量：`:161 stdout: os.Stdout`、`:162 stderr: os.Stderr`、`:164 })` | **成立、未漂**（引用块比派单多一枚闭合行；行号逐字对得上） |
| `cmd/wisp/console_windows.go:38-51`＝"没有父控制台就 return"那段 | `:33-43` 是 `attachParentConsole`：**提前 return 在 `:34-37`**，条件是 `GetStdHandle(STD_OUTPUT_HANDLE)` 拿到**有效句柄**（注释 `:31-32` 自己写明"no-op when a usable stdout already exists"）；`:39` 才是 `AttachConsole(ATTACH_PARENT_PROCESS)` 且**失败不 return**；`:47-51` 的 `rebindStdHandle` 在 `CreateFile("CONOUT$")` 失败时 `return`＝什么都不重绑 | **段号对、句子反**：本文件里"return"那一支守的是"已有可用 stdout"，**不是**"没有父控制台"；"没有父控制台"落在 `:39`＋`:49-51` 的**失败不重绑**上。⇒ §4 按真身写，不照派单那句抄 |
| `cmd/wisp/logsink.go:160` 那台 tee 的 mirror 那一支 | `cmd/wisp/logsink.go:160` 逐字 `mirror: slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})` | **成立、一字未漂** |
| `-H=windowsgui` 在 `scripts/build.ps1` | `scripts/build.ps1:115`（数组枚 `:108-116` 的第七枚；why 注释 `:103-107`），消费在 `:129` 的 `go build -trimpath -ldflags $ldflags -o build/wisp.exe ./cmd/wisp`；**无任何条件分支**（`$ldflags` 一条路拼死）；落档 commit `cc6eaa65`（10-02 19:39 `244-r1: build.ps1 追加 -H=windowsgui…`） | **成立**，且"什么条件下才传"的答案＝**永远传**（详见 §3） |
| `TestAC1ResidentLegInstallsItsLogListenerOnDisk` | 真在 `cmd/wisp/resident_sink_nail_127_windows_test.go:415`（尺：`grep -rn "func TestAC1ResidentLegInstallsItsLogListenerOnDisk" cmd/wisp`＝1 命中） | **不是幻影名**，且它正是 §5 要点名的那枚带载先例 |
| 票 244 的 AC 面 | `.scratch/wisp/issues/244-spec-11-wants-the-gui-subsystem-build.md`：`grep -c "^- \[ \]"`＝**7**、`grep -c "^- \[x\]"`＝**0**（AC#0/1/2/3/4/5/6 全未勾） | **整票未动**；本腿一枚框都没碰 |
| `A619` 那句"AC#1 不勾＝缺 AC#2 真机两发" | `grep -rn "AC#1 不勾" docs/reports/pending-and-issues.md`＝**0 命中**；`"244 AC#2" 真机两发` 的真身在 **`A557` §5（`:11063`）**："只读备料＝…`244 AC#2` 真机两发（要构建＝等写面空窗）"；`A619` 节内（`:12118-12158`）现读**零枚** `244` 命中 | **归错账**：这句话的盘上真身是 `A557` 不是 `A619`。同一族语义在**票面 AC#0**（`244-…md:21`）逐字写着"AC#1 单独做完会留下一枚『看着对、其实没人验过』的产物"＋"AC#1 与 AC#2 不许拆成两批改"，**引用这一格请按票面 `:21`＋`A557` `:11063`，别按 `A619`** |

另：本腿起手时 `git log -1`＝`cfd97636`（10-05 12:58），与 `A621` §3 那句"残余＝没有任何 exec 级用例钉住通道绑定 ⇒ 归票 244，随即派 `244-a4`"同一条链——裁决表凭据句现读在 `docs/evidence/s1/257-clean-machine-provider-registry-v1.md:150`（判语）与 **`:161`**（残余那句逐字："**没有任何用例从真子进程口径测过这枚通道**（两枚测试都在同进程里传 buffer）…要长期钉住需一发行文＋一枚 `exec` 级用例"）。

---

## §1 问 1 —— 现成先例：`cmd/wisp` 里「编译一次 wisp.exe 再用 `exec.Command` 起真子进程」那族

判语：**这一族今天有 10 枚文件／21 枚调用点，但没有一枚是"编译一次"**——`buildWispForTest` 是**每次调用都重建一枚新 exe**的无缓存 helper（`t.TempDir()`＋`go build`＋拷 DLL），"编译一次"这个前提在盘上不存在。⚠ **判过"某句话到了真子进程 stderr"的＝有，但只有 4 枚断言、且全在常驻腿（无参）的 slog-mirror 那一支；`runTextTask` 那枚注入 writer 的话，21 枚调用点里零枚从子进程口径判过。**

现尺与读数（逐条现跑，13:0x）：

```
$ grep -rn "exec.Command(" cmd/wisp | wc -l          # 派单指定的起手尺
23
$ grep -rln "exec.Command(" cmd/wisp | wc -l
13
$ grep -rn "exec.CommandContext(" cmd/wisp | wc -l    # 派单的尺漏了这一形
2
$ grep -rn "buildWispForTest(t)" cmd/wisp | wc -l
21
$ grep -rln "buildWispForTest(t)" cmd/wisp | wc -l
10
```

**两把尺的集合差**（现跑 `comm -23 <(grep -rl "buildWispForTest(t)" cmd/wisp|sort) <(grep -rl "exec.Command(" cmd/wisp|sort)`）：**5 枚文件调 helper 却没有自己的 `exec.Command(`**——`dataroot_128_windows_test.go`（走 `exec.CommandContext`，`:128`）、`resident_approval_246_windows_test.go`、`resident_ball_228_windows_test.go`、`resident_ball_live_228_windows_test.go`、`resident_hotkey_live_258_windows_test.go`（后四枚复用 127 那枚 `bootResidentLeg`，同包可见：现读调用点在 `resident_approval_246_windows_test.go:187`、`resident_ball_228_windows_test.go:56/:187`、`resident_ball_live_228_windows_test.go:101`、`resident_hotkey_live_258_windows_test.go:43/:85/:131/:207`）。⇒ **派单那把尺（`exec.Command(`＝23）会漏计 5 枚真起子进程的文件**，下一位取数请按 `buildWispForTest(t)`＝21 枚调用点这把数。

**helper 真身**（`cmd/wisp/secret_argv_windows_test.go:161`，文件带 `//go:build windows`＝该文件 `:1`）：

- `:172 dir := t.TempDir()` → `:173 exe := filepath.Join(dir, "wisp.exe")` → `:174 cmd := exec.Command(goExe, "build", "-o", exe, "./cmd/wisp")` → `:184-197` 把 `third_party/sherpa-onnx/*.dll` 逐枚拷进同一目录。
- **没有缓存、没有 `sync.Once`**（尺：`grep -rn "sync.Once" cmd/wisp --include=*.go | wc -l`＝2，两枚都在产码里：`notify_windows.go:74`、`panel_resident_windows.go:125`），**没有跨用例共享的 exe 名册**。
- **不带 `-ldflags`** ⇒ 台件二进制是 **CUI 子系统**，与出厂件（`scripts/build.ps1:115` 带 `-H=windowsgui`）**不同形态**——这正是票 244 §5 G5 那句"包内断言天生看不见 subsystem 这一维"的落点（本腿现读复认：`cmd/wisp/secret_argv_windows_test.go:174` 那一发 `go build` 的实参只有 `-o` 与包名两枚）。
- 前置门槛：`runtime.GOARCH != "amd64"` 直接 `t.Fatalf`（`:163-165`）；`third_party/sherpa-onnx` 无 DLL 时 `t.Fatalf` 要你先跑 `scripts/fetch-deps.ps1`（`:184-188`）。

**逐枚点名句柄**（10 枚 `buildWispForTest` 用户，各判什么／怎么读子进程输出）：

| 文件（全在 `cmd/wisp/`） | 调用点 | 子进程 argv | 读 stdout／stderr 的方式 | 它判的是什么 |
|---|---|---|---|---|
| `secret_argv_windows_test.go` | `:271`、`:367` | `secret set … --from-stdin`／`secret set`／`secret set never --value=…` | `:238`、`:339`、`:373` 两枚 `bytes.Buffer` 分开接 `cmd.Stdout`／`cmd.Stderr` | argv 里不许有明文（从子进程 **PEB** 读命令行，不是读 buffer）；`:354`／`:385` 只是**负向**（stdout／stderr 里不许出现那枚 key）。**没有一枚正向断言过任何一句话** |
| `early_log_nail_130_windows_test.go` | `:186`、`:280` | `secret list`（`:121`） | `:124-125` `out`＋`errb` 分开；`:130` 把 `errb.String()`／`out.String()` 一起交回 | 判的是**盘上 JSONL 里早期记录的次序**；stderr 只进 `t.Fatalf` 的诊断文本（`:128`），**不作为断言对象** |
| `dataroot_128_windows_test.go` | `:50` | 4 条腿的表驱动：`run "…"`／`secret list`／`doctor`／**无参常驻**（`:59-64`） | `:131-133` **stdout 与 stderr 合进同一枚 `bytes.Buffer`** | 逐腿判退码≠0 ＋ `markers` 那几句话在**合并输出**里出现（`:74-78`）⇒ ⚠ **这一枚是今天最接近"真子进程口径判一句话"的先例，也是唯一一枚把 `run` 腿真跑起来判句子的**，但它**分不清那句话落在 stdout 还是 stderr**（通道绑定对它失明） |
| `resident_sink_nail_127_windows_test.go` | `:416`、`:521`、`:565` | **无参**（`:219 cmd := exec.Command(exe)`） | `:229 cmd.Stdout, cmd.Stderr, cmd.Stdin = leg.stdout, leg.stderr, nil`，两枚 `lockedBuf`（`:174-193`，带 `has(sub)`） | ★ **这族里唯一真判过"某句话到了真子进程 stderr"的**：`:487 leg.stderr.has(residentInstallMsg)`、`:493 strings.Count(leg.stderr.String(), residentEarlyResolverMsg) != 1`、`:529 pollUntil127(… leg.stderr.has(residentRefusalPrefix))`、`:541 leg.stderr.has(residentRefusalPromise)`。四枚判的全是 **slog mirror（`logsink.go:160`）** 那一条通道，**不是 `runTextTask` 的注入 writer** |
| `resident_task_source_246_windows_test.go` | `:297`、`:359`、`:415` | **无参**（`:457`，`bootResidentLegWithEnv`） | 复用上面那枚 `residentLeg`＋`lockedBuf`（`:462`） | 判三句话到**真子进程 stdout**（`:302 :303 :368-370 :428-430` 全用 `leg.stdout.has(…)`）＋盘上无 `wisp.db`。⚠ 这批句子是常驻腿自己的启动汇报（`fmt.Printf`），不是 `runTextTask` 的回执 |
| `resident_task_source_live_246_windows_test.go` | `:123`、`:261` | wisp 腿＝**无参**（`bootResidentLegWithEnv`，用例在 `:124`／`:261-262` 取 exe）；另一枚 exec 是 **Esc 台件** `esclistener -watch`（`:493`，rig 由 `buildEscListener246` 现编，见 `resident_approval_live_246_windows_test.go:417`，源在 `cmd/wisp/testdata/esclistener`） | `:498 cmd.StdoutPipe()` 逐行读台件输出；wisp 腿走 `lockedBuf`；`:506 cmd.Stderr = a.errBuf` | 真窗／真卡那两发的活体读数全在 **`leg.stdout.has(…)`**（`:134 :149 :239 :270`）；台件 stderr 只作诊断 |
| `resident_approval_246_windows_test.go` | `:185` | 无参 | 经 `bootResidentLeg`（同 127 那枚 helper） | "出货进程真持有那一步"（cancel step 挂钩），盘上记录为主 |
| `resident_ball_228_windows_test.go` | `:54`、`:185` | 无参 | 同上 | 球进常驻的报与账（盘上＋stdout 句子） |
| `resident_ball_live_228_windows_test.go` | `:90` | 无参 | 同上 | 真机桌面上球窗存在性 |
| `resident_hotkey_live_258_windows_test.go` | `:41`、`:78`、`:121`、`:194`（**四枚＝本族最密**） | 无参 | 同上 | 启动汇报里点名来源与默认、配置热键、免重启重绑、占用组合报新值 |

**不属于这一族、但被派单那把尺一并捞进来的 3 枚**（逐枚具名，免得下一位误当先例）：
- `cmd/wisp/panel_host_gate_test.go:310/:322/:340/:354/:394/:407`＝起的是 `git`（rev-parse／archive）与一枚 `go build ./...`，**不 exec `wisp.exe`**。
- `cmd/wisp/resident_approval_live_246_windows_test.go:417`＝`go build -o … ./cmd/wisp/testdata/esclistener`，`:446` exec 的是那枚**台件**，不是 wisp.exe。
- `cmd/wisp/slo_report_144_windows_test.go:881`＝`cmd.exe /c exit 7`；`config_reload_perm_223_windows_test.go:52`／`logsink_windows_test.go:112`＝`icacls`。
- 另有两枚 **re-exec 自家测试二进制**（不是 wisp.exe）：`secret_argv_windows_test.go:450`（正控：往子 argv 里种明文再要求同一把尺抓到）与 `slo_exit_os_156_windows_test.go:322`；加 `leg_dispatch_gate_133_test.go:1580 exec.CommandContext(ctx, exe, "-test.list", ".*")`（`exe`＝`os.Executable()`＝测试二进制本身，见 §5）。
- 产码侧两处 `exec.Command`（非用例）：`cmd/wisp/doctor.go:316`（`cc --version`）、`cmd/wisp/slo_windows.go:492`。

**子进程收尾**：这一族**没有"起了不收"的形状**——`bootResidentLeg` 与 `bootResidentLegWithEnv` 都在 `Start()` 之后立刻 `go func(){ leg.done <- cmd.Wait() }()` ＋ `t.Cleanup(leg.stop)`（`resident_sink_nail_127_windows_test.go:230-239`、`resident_task_source_246_windows_test.go:463-469`，注释 `:235-238` 逐字写明理由："no child may outlive it: an orphaned wisp.exe would hold its temp dir (and, in the dev/prod envs, the session mutex) open past the test"）；`slo156Spawn`（`slo_exit_os_156_windows_test.go:318-328`）`Start()` 后 `t.Cleanup(func(){ _ = cmd.Wait() })`；`dataroot_128` 用 `exec.CommandContext`＋`cmd.Run()`。**唯一"输出不设防"的先例＝`slo156Spawn` 的 `cmd.Stdout, cmd.Stderr = nil, nil`（`:324`）**，那是刻意的（子进程 PRINTS 不该被误当观察者控制台），照抄进本问那枚新用例就是盲的。

---

## §2 问 2 —— 通道绑定的真身：`os.Stdout`／`os.Stderr` 交给了谁

判语：**绑定发生在 `cmd/wisp/main.go:159-164`，那三行今天既没有被任何用例驱动、也没有被任何真子进程判过**——`cmdRun`（唯一持有那枚 `runSpec` 的函数）全仓**只有 2 处命中**（尺：`grep -rn "cmdRun(" cmd/wisp --include=*.go` → `main.go:91 os.Exit(cmdRun(args[1:]))`、`main.go:151 func cmdRun(args []string) int`），**零枚测试调它**。`runTextTask` 另有 8 枚测试调用点，**8 枚全部在同一进程里传 `bytes.Buffer`／`io.Discard`**。⇒ 257-v1 `:161` 那句残余**成立**，且本腿把它量窄了一格：**连"同进程但走真文件描述符"那一形都没用在 `runTextTask` 上**（`captureLeg131` 那台管道机只喂给了 `cmdModels`／`cmdSecret`，见下）。

**绑定逐字**（现读 `cmd/wisp/main.go`）：

- `:89-91` `case "run": attachParentConsole(); os.Exit(cmdRun(args[1:]))`——**attach 早于绑定**（与票 244 §9 第 5 条同一条读数，本腿复认：`attachParentConsole` 在 `:90`，`os.Stdout` 捕获在 `:161`）。
- `:151-165 cmdRun`：`:152 printVersions("")`（两行版本走 `fmt.Printf`＝**stdout**）→ `:153 reply := interactiveStdin()` → `:154-158` 无控制台时向 `os.Stderr` 打那一句"L2 卡会等到超时后按拒绝处理…" → `:159-164 return runTextTask(runSpec{argv: args, stdout: os.Stdout, stderr: os.Stderr, reply: reply})`。
- 下游：`cmd/wisp/run.go:104-105` 是 `runSpec` 的 `stdout io.Writer`／`stderr io.Writer` 两枚字段；`:182-187` 的 nil 兜底把默认再绑一次 `os.Stdout`／`os.Stderr`（**运行时读全局**，不是构造时绑定）；`:245 ensureFirstRunConfig(s.dataDir, s.stderr)` 把 `s.stderr` 交给回执生产者（`cmd/wisp/firstrun.go:72`，句面 `:92-95` 与 `:111-…` 那两枚 `fmt.Fprintf(stderr, …)`）。
- ⚠ **`stderr` 那一枚 writer 同时还是错误与归因句的出口**（`run.go:196` 用法句、`:208` 起的数据根拒绝句……），所以"首启回执真到达"这句话在产码侧本来就与退码 2 那一族共用同一条通道。

**`runTextTask` 的 8 枚测试调用点（尺＝`grep -rn "runTextTask(" cmd/wisp | wc -l`＝9，扣 `main.go:159` 与 `run.go:181` 定义）**：

| 调用点 | 传的 writer | 口径 |
|---|---|---|
| `cmd/wisp/firstrun_198_test.go:32`（经 `run198` `:29-40`） | `stdout: out, stderr: errb`（`:31` 两枚 `bytes.Buffer`） | 同进程内存。★ 那枚"新建默认配置＋路径"的正向断言就在 **`:98`**，`t.Logf` 逐字回显 stderr 在 `:101` |
| `cmd/wisp/firstrun_257_test.go:79` | 复用 `run198` | 同进程内存（尺＝`grep -rn "run198(" cmd/wisp | wc -l`＝**15 处命中**：定义 1 ＋调用 14，分布在 `firstrun_198_test.go`(7)、`firstrun_198r2_test.go`(5)、`firstrun_257_test.go`(1)、`firstrun_257_nonpreset_test.go`(1)、`firstrun_acl_198_windows_test.go:24`(1) ⇒ **这一族的 5 枚文件共用同一枚同进程 helper，没有一枚走子进程**） |
| `cmd/wisp/firstrun_257_nonpreset_test.go:55` | 复用 `run198` | 同上 |
| `cmd/wisp/run_test.go:122` | `runFixture` 的 buffer | 同进程内存 |
| `cmd/wisp/run_mode101_test.go:148` | 同上族 | 同进程内存 |
| `cmd/wisp/approval_reply_201_test.go:172` | 同上族 | 同进程内存 |
| `cmd/wisp/dataroot_128_test.go:117` | `stdout: io.Discard, stderr: &out` | 同进程内存 |
| `cmd/wisp/logsink_windows_test.go:152`、`:319` | `out, errb := &bytes.Buffer{}, &bytes.Buffer{}` | 同进程内存 |

**"真到达"这一格的三枚最接近先例，各自差在哪一格**（全部现读复认，⛔ 都不算通道钉）：

1. `cmd/wisp/dataroot_128_windows_test.go:131-133`＋`:74-78`——**真子进程**、**真 argv**（含 `run` 腿 `:60`），但 stdout 与 stderr **合进同一枚 `bytes.Buffer`** ⇒ 它能说"那句话到了某处"，**说不出"那句话到了 stderr"**。这一枚是"差半枚口径"的先例。
2. `cmd/wisp/leg_sink_nail_131_windows_test.go:161-207 captureLeg131`——把 `os.Stdout`/`os.Stderr` 换成**真 `os.Pipe()`**（`:164-172`，收尾 `:190-197` 还原），再在 `:351-357` 用 `cmdModels(…, modelsIO{stdout: os.Stdout, stderr: os.Stderr, …})` 判 `:393 !strings.Contains(stderr, "wisp models ensure: 交还被拒绝")`。**这是仓里最硬的一枚"句子真写进 stderr 文件描述符"的钉**，但它（i）是**同进程**（`go test` 的进程自己换柄，不是 `wisp.exe` 起来自己绑），（ii）钉的是 `cmdModels` 那枚 `modelsIO`，**不是 `runTextTask` 的 `runSpec`**。同族还有一枚更软的：`cmd/wisp/panel_assets_143_test.go:85-86` 把 `os.Stdout/os.Stderr` 换成**临时文件**再直接调 `cmdPanelAssets`。
3. `cmd/wisp/resident_sink_nail_127_windows_test.go:487/:493/:529/:541`——**真子进程＋真 stderr**，四枚正向 `has(...)` 都在，但**那条通道是 slog 的 mirror（`logsink.go:160`），不是 `runSpec.stderr`**：无参常驻腿根本不进 `runTextTask` 的 CLI 分支（它走 `resident_windows.go:33`），而首启回执只在 `run` 腿产生（`run.go:245`）。⇒ **两枚先例各占一半：一半有真进程没有真 stderr，另一半有真 stderr 没有那枚 writer。**

**量不到的那一格具名**：`run198` 那族用例"stderr 里有没有回执"与真进程"stderr 句柄到底是不是那根管道"之间的这一步，**今天整个包没有一枚用例跨过去**（尺：本包 22 处 `os.Stdout`/`os.Stderr` 命中全在测试自己的换柄或子进程 argv 里，`grep -rn "os\.Stdout\|os\.Stderr" cmd/wisp --include=*_test.go | wc -l`＝22；**零枚**用例从 `exec.Command(exe, "run", …)` 之后读 `cmd.Stderr` 里的首启回执串）。

---

## §3 问 3 —— subsystem 那一格：`-H=windowsgui` 在 `scripts/` 的哪一行、PE subsystem 真身

判语：**旗标在 `scripts/build.ps1:115`、无条件、今天就在出厂链里；但盘上那枚 `build/wisp.exe` 比这面旗早 2 天，`objdump` 现读仍是 `Subsystem 00000003 (Windows CUI)`**——"默认构建出来的 exe 的 subsystem"这一问有**两枚答案**，必须分开交：**"链今天会产出什么"＝GUI(2)〔由 `.scratch/probes-244-r1/out/wisp.exe` 的真读数佐证，那枚是 244-r1 用同一份 `$ldflags` 现编的〕；"仓里现在躺着什么"＝CUI(3)**。⛔ 本腿不跑构建（`build.ps1:129` 就是 `go build`），所以"今天再跑一次 build.ps1 会得到 GUI(2)"这一发**量不到，归编排者**——它正是票 244 **AC#1** 那句"改前 CUI(3) 与改后读数都要写在交件里"欠的那一发。

**旗标那一行与条件**（现读，⛔ 没有 if）：

```
scripts/build.ps1:103-107   六行 why 注释（点 SPEC-11 §2.2 ＋点名 attachParentConsole 的早退条件）
scripts/build.ps1:108-116   $ldflags = ( 六枚 "-X ..." , "-H=windowsgui" ) -join ' '
                     :115       ← 旗标本体，数组第七枚，**无守卫、无条件分支**
scripts/build.ps1:121-124  环境变量：CGO_ENABLED=1 / CC / GOOS=windows / GOARCH=amd64
scripts/build.ps1:129      & $go.Source build -trimpath -ldflags $ldflags -o build\wisp.exe ./cmd/wisp
scripts/build.ps1:130      if ($LASTEXITCODE -ne 0) { Fail ... }
scripts/build.ps1:166-169  冒烟：& build\wisp.exe doctor ＋ if ($LASTEXITCODE -ne 0) { Fail }
```

现尺（13:1x 全部现跑）：

```
$ grep -rn "windowsgui" scripts | wc -l
4        # build.ps1:103（注释）/ :115（旗标）/ slo-check.ps1:58（注释）/ :158（注释）
$ git log --format="%h %cI %s" -S'-H=windowsgui' -- scripts/build.ps1
cc6eaa65 2026-10-02T19:39:20+08:00 244-r1: build.ps1 追加 -H=windowsgui，双击不再带黑控制台
$ git log --format="%h %cI %s" -3 -- scripts/build.ps1     # 旗标之后没有第二笔
cc6eaa65 …   ← 此后再无人动过 build.ps1
$ git ls-files "*.exe" | wc -l
0        # 库里零枚 tracked exe ⇒ 任何"读 PE subsystem"的尺都必须先构建（复认 244-c2 §3 那句）
```

**PE 头真身**（⛔ 没读注释判形态；尺＝仓里既有定式 `objdump -p <exe> | grep -i subsystem`，本机 `command -v objdump`＝`/e/work/base/msys64/mingw64/bin/objdump` **有**）：

| 产物 | mtime／字节（`ls -l --time-style=long-iso`） | `objdump -p` 现读 |
|---|---|---|
| `build/wisp.exe` | **2026-09-30 23:43**，30,299,552 | `Subsystem 00000003 (Windows CUI)`（`MajorSubsystemVersion 10`） |
| `build/wisp228.exe` | 2026-09-30 16:14，30,184,669 | `Subsystem 00000003 (Windows CUI)` |
| `build/balldebug.exe` | 2026-09-30 11:21，7,295,488 | `Subsystem 00000003 (Windows CUI)` |
| `.scratch/probes-244-r1/out/wisp.exe`（244-r1 的临时出厂件，**没覆盖 `build/`**；它的构建脚本 `.scratch/probes-244-r1/build-probe.ps1` 现读 `:103/:108/:115/:129` 与出厂链逐行同形） | 2026-10-02 19:37，30,973,922 | **`Subsystem 00000002 (Windows GUI)`** |
| `cmd/wisp` 台件（`buildWispForTest` 现编，`cmd/wisp/secret_argv_windows_test.go:174`） | 每次 `t.TempDir()`，跑完即弃 | **量不到**（不许跑 `go build`）；⇒ 只能由那一发的实参**只有 `-o`＋包名、零枚 `-ldflags`** 判它**必然是 CUI**〔读码，非读数〕 |

**新旧判定**（现尺）：`ls -l --time-style=long-iso cmd/wisp/main.go cmd/wisp/run.go cmd/wisp/firstrun.go`＝`main.go 2026-09-30 23:16`／`run.go 2026-10-02 16:23`／`firstrun.go **2026-10-05 12:36**`。⇒ **`build/wisp.exe`（09-30 23:43）比 `run.go` 与 `firstrun.go` 都旧，且比那面旗（`cc6eaa65`，10-02 19:39）早 2 天**——票 244 §9 第 2 条那句"r1 的 GUI 产物落在 `.scratch/probes-244-r1/`、没覆盖 `build/`"**本腿现读复认成立**，而且它比 `A558` 当时记的"09-30 23:43"没有动过（同 mtime、同字节数）。⇒ **今天任何"双击已经没有黑框了"的句子都是假话**（票面 §9 第 2 条逐字禁止这一句）。

**顺带量到一枚本问没问、但直接打在 §6/§7 上的形状**（⚠ 标〔读码推断＋盘上尺〕，本腿不许跑 `build.ps1` 所以**没有读数**）：`scripts/build.ps1:166-169` 的冒烟是 `& build\wisp.exe doctor` ＋ `if ($LASTEXITCODE -ne 0) { Fail }`，同一文件 `:28` 是 `Set-StrictMode -Version 2.0`；旗标落地后 `wisp.exe` 是 GUI 子系统件，**`&` 不等它退出**——这条因果链在本仓已有**两枚具名实证**：`scripts/slo-check.ps1:50-64`（"since ticket 244 scripts/build.ps1 passes -H=windowsgui … a bare call to one is NOT waited for"）与 `:155-161`（实测对照："the console build returned 2049ms later with exit=3, the GUI build fell through in 70ms"，探针盘上现读在 `.scratch/wisp/probes/263/r1/`，含 `bin/`）。`build.ps1` 里 `:130` 那发 `go build` 已经把 `$LASTEXITCODE` 赋成 0，`:169` 读到的**很可能是上一发的 0**，于是冒烟步在"doctor 到底跑没跑成"这一格**恒绿**＝**step 层的绿等于什么都没测**的现成形状。⇒ **具名归口**：这是票 244 **AC#6**（票面 `:31-33`"下游消费者名册"）射程内的一枚消费者，本腿只登记、⛔ 不判它今天是否真的恒绿（判它要跑构建）。

---

## §4 问 4 —— 无父控制台那一支：`console_windows.go` 的 return 与 `logsink.go` 的 mirror 谁先谁后

判语（三句，每句都指得到行）：
1. **顺序上没有错位**：`attachParentConsole` 永远**早于**"writer 绑定"与"mirror 绑定"这两件事，所以 `logsink.go:160` 的 mirror 与 `runSpec.stderr` 取的是**同一批柄**——不存在"mirror 抓到重绑之前的旧柄"这种形状（今天没有那条路径）。
2. **但"双击起法"那一支压根不产生那句话**：首启回执只在 `runTextTask` 里印（`cmd/wisp/run.go:245` → `cmd/wisp/firstrun.go:92-95`），而双击＝**无参**＝`main.go:60-67` 的 resident 分支，**从不进 `cmdRun`**。⇒ ⚠ 具名推翻派单的一处隐含前提：**票 257 那格"stderr 真到达"与票 244 的"双击那一支"不是同一条通道**，一枚 exec 级用例能钉的是"从终端起／被重定向"那一支，**买不到 AC#2 ③ 那一发双击**。
3. **`console_windows.go` 里那句"没有父控制台就 return"在本文件不存在**（见 §0 第 2 行）：早退的条件是**"已经有可用 stdout"**（`:34-37`），"没有父控制台"落在 `:39` 的 `AttachConsole` 失败＋`:47-51` 的 `CreateFile("CONOUT$")` 失败**静默不重绑**上。这一处方向之差直接决定 stdout 那一行的去向，所以下面逐态写。

**谁先谁后（现读逐条，全带行号）**：

| 腿 | 顺序 | 绑定发生在 |
|---|---|---|
| `wisp run`（终端起／重定向起） | `main.go:90 attachParentConsole()` → `:91 os.Exit(cmdRun(…))` → `main.go:152 printVersions`（`fmt.Printf`＝stdout）→ `:154-157` 无控制台时向 `os.Stderr` 喊那一句 → **`main.go:161-162` 把 `os.Stdout`／`os.Stderr` 交进 `runSpec`** → `run.go:181 runTextTask` → `run.go:222 installLogSink`（mirror 在 `logsink.go:160` 读**那一刻的** `os.Stderr`）→ `run.go:245 ensureFirstRunConfig(s.dataDir, s.stderr)` | attach **早于**绑定（票 244 §9 第 5 条同读数，复认） |
| `wisp`（无参＝双击那一支） | `main.go:65 attachParentConsole()` → `:66 runResident()` → `resident_windows.go:34 printVersions`（stdout）→ `:66 installLogSink` → `:71` 装不上时向 `os.Stderr` 喊"持久日志未启用" → `:88` 启动汇报 `fmt.Printf`（stdout） | 同上：attach 早于 mirror |
| 常驻腿的内部任务源（**今天唯一能印首启回执的非 CLI 入口**） | `resident_task_source_windows.go:266-267 stdout: os.Stdout, stderr: os.Stderr`（`runSpec` 字面量），且这条腿**只在 `console != nil` 时才受理任务**（`:346-348` 那个 `if console != nil` 才起 `runConsoleLoop`） | 无 attach 之外的额外绑定 |

**双击起法下"stdout 那一行"到底去了哪**——两态分开写，⛔ 不许合成一句：

- **态甲＝今天盘上那枚件（CUI(3)，§3 现读）**：无参双击 ⇒ 系统为这个 CUI 进程**新建一枚控制台**（＝票 244 票面 `:5` 与 `:15` 记的那只黑框），`attachParentConsole` 在 `:34-37` 因为 `GetStdHandle(STD_OUTPUT_HANDLE)` 有效**直接 return**，于是一切 `fmt.Printf` 都写进那只黑框。⇒ **看得见，但看见的方式就是那只框**；`resident_windows.go:55-60` 那段注释（"double click the icon, no terminal attached, stderr going nowhere"）**说的不是态甲、是态乙**——今天盘上这枚件双击起来是有框、有 stderr 去向的，这正是票面 AC#3 要点名改口的两行注释之一。
- **态乙＝旗标落地之后（GUI(2)）**：无父控制台 ⇒ `:39` 的 `AttachConsole(ATTACH_PARENT_PROCESS)` 失败，`:40-42` 三枚 `rebindStdHandle` 在 `:48-51` 因为 `CreateFile("CONOUT$")` 失败而**逐枚静默 return**，`os.Stdout`／`os.Stderr` **保持进程启动时的柄**。那一枚"启动时的柄"到底是什么（Explorer 给 GUI 进程的是 NUL 还是无效柄）**本腿量不到**——判它要一发真机 GUI 双击（＝票面 AC#2 ③，归编排者；`docs/BUILD.md:90` 那枚"当时的临时构建"读数**不可继承**，票面 AC#0 与 §9 第 2 条两处都写了这句）。**两支共同点只有一条：写进一枚没有读者的柄**，所以态乙下"那句话消失"是**无害的哑**，不是丢数据的崩（票面 `:5` 第 ③ 行"最坏后果是什么形状"就是这么裁的）。
- **态丙＝本腿唯一能从盘上反证的那一支（测试起法）**：21 枚台件都是**把子进程 stdout/stderr 重定向到管道**起的（`resident_sink_nail_127_windows_test.go:229`、`resident_task_source_246_windows_test.go:462`、`secret_argv_windows_test.go:238`……），重定向本身使 `GetStdHandle` 拿到**有效管道柄** ⇒ `console_windows.go:35-37` 早退 ⇒ 句子的去向＝**那根管道**。⇒ 盘上佐证：127 那四枚正向 `leg.stderr.has(...)`（`:487/:493/:529/:541`）今天在默认档名册里是绿的（`docs/evidence/s1/228-resident-ball-v1.md:13` 记的默认档整包 `rc=0／零枚 --- FAIL／148.168s` 就含这一族）。〔读码推断＋盘上第二手读数，本腿没跑任何用例〕
  ⚠ **这一态的推论要给写腿**：新用例**在测试进程里能可靠看到那枚首启回执**，恰恰是因为"重定向使 attach 早退"这一条——**它不是"双击也看得见"的证据**，两个态别混进同一句判语。

**顺手量到的一格（本问没问，归票 244 AC#6 的名册）**：`cmd/wisp/main.go` 的 switch 里**两枚分支不调 `attachParentConsole`**——`case "secret"`（`:100-101`）与 `case "slo"`（`:110-111`），而 `secret.go:200-202` 照样绑 `os.Stdout`／`os.Stderr`、`:226` 照样 `installLogSink`。⇒ 旗标落地后这两条腿在 GUI 子系统下的可见输出**只靠继承柄**：从终端起＝继承的是控制台柄（有效，attach 那段本来也不需要跑）；管道起＝有效；**双击／无终端起＝与态乙同形**。`docs/BUILD.md` 不在本腿射程，且这两枚是不是"该 attach"属设计判断——⛔ 本腿只登记，判它归编排者。

---

## §5 问 5 —— 撞钉预检：新增一枚 exec 级用例会撞哪些既有计数尺

判语：**会打红的只有 A 组两枚 AST 门的"仪器形"红，且三枚禁行写死；B/C/D 三组共 12 枚计数形钉全部钉在"自己那枚子进程／自己那枚 temp 目录／产码文件形状"上，加一枚文件、加一枚用例都打不红它们**。另有一枚**反向陷阱**必须写进派单：`firstrun_198_test.go:91` 那句"stdout 必须为空"在**真子进程口径下天然不成立**（版本行三行走 stdout），照抄它会把新用例生下来就判红。

⚠ **本节的引文口径**（免得下一位把折叠当篡改）：所有"逐字"都是**同一行里的原文**，只把连续空白（缩进、注释块的对齐空格）折成单个空格，⛔ 没有改过任何字符、没有省略中段；带 `…` 的那几处才是有省略，且省略号前后仍是原文。

**A 组｜会读到你新加那枚文件的两扇门**（⚠ 行为型钉不写常量名只写产物，这两枚就是；光 grep 新符号名不算数）

| 尺／行 | 断言原句（逐字） | 加一文件／加一用例会不会红 |
|---|---|---|
| `cmd/wisp/leg_dispatch_gate_133_test.go:179` `pkg, err := loadPackage133(".")` → `:365 entries, err := os.ReadDir(dir)` → `:372 if e.IsDir() \|\| !strings.HasSuffix(name, ".go") { continue }` → `:438 isTest := strings.HasSuffix(name, "_test.go")` | 见左：**`_test.go` 一样进 AST**，`isRunnableCase133`（定义在 `:516`）把真用例分进 `p.tests`、其余分进 `p.helpers`（`:468-472`，注释逐字："Helpers, methods, benchmarks: still walkable in loose mode … They are not cases, and they can no longer cover a leg by being named after one."） | **不因"多一枚文件"红**；红的是下面四枚仪器形 |
| 同文件 `:213-216` `for _, b := range pkg.blind { t.Errorf("AC#1/#2 RED (the instrument, not the code): %s\n"…` | blind 只有 5 个追加点：`:758`／`:765`（"calls %s, a package-level function value this directory declares but never defines"）、`:788`（"calls through field %s of package-level var %s…"）、`:1753`（"cannot read %s to scan it for coverage rulings"）、`:1768`（"%s:%d carries %s with no leg name after it"） | ⛔ **禁行 1**：新文件里不许新增包级 `var x = someFunc`／结构体字段携带函数；⛔ **禁行 2**：不许出现 `WISP-LEG-COVERAGE-RULING` 词面却不带真腿名（`main.go:74-87` 那三枚是唯一合法形） |
| 同文件 `:187-190`（下界常量 `minLegs133 = 11` 定义在 `:148`） | `if len(legs) < minLegs133 { t.Errorf("AC#1/#2 RED: the leg census enumerated %d legs (%d is the floor, measured when this gate landed), so a dispatch rewrite is being read as \"no legs\" rather than reported. Census: %s", …) }` | 只数 `func main` 的分支 ⇒ **加用例不红**；**加命令**才红（还要同步 `usage` 块，见 `:218 censusVsUsage133`） |
| 同文件 `:195-206`（orphan sink callers） | `t.Errorf("AC#1 RED: %s reaches %s (or is a production entry of this package's dispatch), and no chain from func main reaches it any more.…")` | `sinkCallers133` 只走 `allDecls133()`＝`p.decls`＋`p.methods`（非测试声明）⇒ **测试用例与测试辅助函数都不进这张图**，新用例就算直接调 `installLogSink` 也不会被记成 orphan（现例：`logsink_test.go:128` 今天就在调它且门是绿的） |
| 同文件 `:227-230`＋`:1632 runRosterReds135` | `t.Errorf("AC#1/#2 RED: leg %q (%s) is booked in the ledger below as covered by the case %q, and this round's test binary has no such case…")` | 只有**登记名册 `legCovers133`（`:169-174`）里的每一枚名**被逐字核对；枚数 `n` 只进 disclosure 句（`:1687`）不断言 ⇒ **加用例不红**；⛔ **禁行 3**：不许把新用例名塞进 `legCovers133`／`legNails131` 那两张登记册，却让它落在本编译拿不到的 tag 后面（那正是 R-133-9 的形状） |
| 同文件 `:1543 const rosterReadTimeout135 = 2 * time.Minute`（注释 `:1540` 记实测 0.098 s） | 起子进程跑 `exe -test.list '.*'`（`:1580`） | 这枚门**每次跑都要列全包名册** ⇒ ⛔ 新用例**绝不允许在 `init()` 或包级变量初始化里跑构建**，否则 2 min 预算被吃穿 |
| `cmd/wisp/leg_sink_gate_131_test.go:559-586 loadMainPackage131`（注释 `:548-558` 逐字："parses every .go file in dir, **whatever its build constraints say**"） | `:382 if len(legNails131) == 0 {`＋同族 blindness check | 加用例不红；⚠ 具名：`:587-592` 的 ruling 扫描**跳过 `_test.go`**（`if !isTest { …bytes.Contains(src, []byte(legSinkRulingMarker)) }`）⇒ 在测试文件里写 marker **不会被算作一条 ruling**，别指望它当豁免用 |

**B 组｜stdout／stderr 形状与枚数钉**（逐枚原句＋判语）

| `file:line` | 断言原句（逐字，取判据那一行） | 判"加一文件／加一用例" |
|---|---|---|
| `cmd/wisp/firstrun_198_test.go:91-93` | `if stdoutText != "" { t.Errorf("stage 1: stdout = %q, want empty (a refused run streams no task output)", stdoutText) }` | 打不红（同进程 buffer）。⚠ **反向陷阱**：真子进程口径下 stdout＝三行版本（`main.go:152 printVersions`→`fmt.Printf`），`257-v1:159` 现量 119 字节 ⇒ **新用例不许照抄这句判据** |
| `cmd/wisp/firstrun_198_test.go:98-99` | `if !strings.Contains(stderrText, "新建默认配置") \|\| !strings.Contains(stderrText, cfgPath) { t.Errorf("stage 1: the run created the file but never said so with its path …") }` | 打不红。★这一枚就是**要复制到子进程口径的那句判据**，串与路径两枚都在 |
| `cmd/wisp/firstrun_257_test.go:88`／`firstrun_257_nonpreset_test.go:59` | `if !strings.Contains(stderrText, "新建默认配置") {` | 打不红（同一枚 helper） |
| `cmd/wisp/resident_sink_nail_127_windows_test.go:493` | `if n := strings.Count(leg.stderr.String(), residentEarlyResolverMsg); n != 1 {` | 数的是**它自己那枚子进程**的 stderr ⇒ 打不红。⚠ 唯一会连坐的形＝新用例复用 `bootResidentLeg` 又改动那台 tee |
| `cmd/wisp/resident_sink_nail_127_windows_test.go:487/:529/:541` | `if !leg.stderr.has(residentInstallMsg) {`／`pollUntil127(200, func() bool { return leg.stderr.has(residentRefusalPrefix) })` | 同上，打不红 |
| `cmd/wisp/early_log_nail_130_windows_test.go:257/:266` | `if n := strings.Count(stderr, early130ResolverMsg); n != 1 {`／`… early130InstallMsg); n != 1 {` | 打不红（`runSecretListLeg` 每发一枚新进程） |

**C 组｜文件枚数／名册钉**

| `file:line` | 原句 | 判语 |
|---|---|---|
| `cmd/wisp/firstrun_257_test.go:435` | `t.Errorf("AC#3 RED: firstrun.go now declares %d functions, want 1 (ensureFirstRunConfig). A new helper here is a new surface…")` ＋ `:437` 逐字签名钉 `"func ensureFirstRunConfig(dataDir string, stderr io.Writer) (bool, error) {"` | ⇒ ⛔ 新用例**不许往 `firstrun.go` 加任何 helper、不许改那枚签名**；"一发行文"只能动 `firstrun.go` 里 `:92-95`／`:111-…` 那两枚 `Fprintf` 的**文案串**，签名与枚数一动不动 |
| `cmd/wisp/logsink_test.go:81-101` | `before, err := os.ReadDir(cwd)` … `if len(after) != len(before) { t.Errorf("the refused install changed the working directory: %d entries before, %d after", len(before), len(after)) }`（cwd＝**包源码树**，该用例注释 `:72-75` 逐字："which for a test run is the package source tree"） | ⇒ ⛔ 新用例**绝不许把编出来的 exe、拷出来的 DLL、任何日志落进 `cmd/wisp/`**；沿用 `buildWispForTest` 的 `t.TempDir()`（`secret_argv_windows_test.go:172`）就是安全形 |
| `cmd/wisp/dataroot_128_windows_test.go:79` | `assertDirEmpty128(t, cwd, "real leg "+leg.name)` | 打不红（cwd 是它自己 `:66` 的 `t.TempDir()`）；⚠ 但它**只判合并输出**，别拿它当通道证据 |
| 六枚跳过 `_test.go` 的 AST 尺：`cmd/wisp/config_receipt_255_test.go:295`、`leg_dispatch_gate_133_test.go:1748`、`resident_ball_228_test.go:74`、`resident_grant_writer_265_windows_test.go:547`、`panel_inbound_33_test.go:213`＋`:289`、`panel_host_gate_test.go:218` | 例：`if e.IsDir() \|\| !strings.HasSuffix(name, ".go") \|\| strings.HasSuffix(name, "_test.go") { continue }` | 加测试文件**零影响**（已逐枚读到位） |

**D 组｜构建次数／耗时预算**（现尺：`grep -rn "buildWispForTest" cmd/wisp --include=*_test.go | grep -c "len("`＝**0**；`grep -c "timeout-minutes" .github/workflows/ci.yml`＝**0**）
⇒ **没有任何一枚仪器钉构建枚数或包时长**。代价落在时长上：默认档整包现成读数 **148.168 s**（`docs/evidence/s1/228-resident-ball-v1.md:13`），台账 `A621` §7 给 `268-r1` 的排程口径是"整包预期 ~8.5 min 不许压种子"。每多一枚 `buildWispForTest` 调用＝**多一发完整 cgo＋sherpa 链接**（无缓存，§1 已量），⚠ 这就是"复用现成 exe"那一格真正的价：想要它快，就得先把 helper 改成跨用例缓存——那是**改 21 枚调用点共用的台件底座**，属另一格活，不在"一发行文＋一枚用例"的射程内。

**E 组｜真窗毒源那一族＋"起了子进程却没收尾"那一形**
- 派单点名的先例真身＝`cmd/wisp/resident_sink_nail_127_windows_test.go:415 TestAC1ResidentLegInstallsItsLogListenerOnDisk`。它的毒源形状量得出三样：**无参常驻腿**＝真球窗＋四枚全局热键＋**每会话单实例互斥**（`cmd/wisp/resident_windows.go:42 rt, err := proc.Boot(env)`，全仓只有 `resident_windows.go:42` 与 `slo_windows.go:261` 两枚 `proc.Boot` 调用者，现尺）＋`CREATE_NEW_PROCESS_GROUP`/`GenerateConsoleCtrlEvent`（`resident_sink_nail_127_windows_test.go:222/:253`）。
- 现尺：`grep -rn "t.Parallel()" cmd/wisp --include=*_test.go | wc -l`＝**0** ⇒ 同包内串行；但**两发并发 `go test ./cmd/wisp` 一定互洗**——`resident_task_source_246_windows_test.go:432-436` 逐字把这形写成红：`"AC#7 RED (desktop state, not our code): a dev Wisp is already running in this session, so this leg handed its activation over and exited. Stop it and re-run."`。
- ⇒ **具名结论（这格最值钱的一句）**：新用例走 **`run` 腿**就不继承这枚毒源（`cmd/wisp/run.go` 零枚 `proc.Boot`），这是它今天能被放进默认档门禁当尺用的前提；反之**复用 `bootResidentLeg` 就连带继承热键、真窗与互斥**。
- "起了子进程却没收尾"：**名册里零枚**（§1 末逐枚列的 `t.Cleanup(leg.stop)`／`go func(){…cmd.Wait()}`／`CommandContext` 三形覆盖全部 21＋2 枚调用点）。最接近的是 `cmd/wisp/slo_exit_os_156_windows_test.go:318-328 slo156Spawn`：`Start()` 后**不由自己 Wait**、由 witness 从 PID 重开柄，但它自带 `t.Cleanup(func() { _ = cmd.Wait() }) // never leave a child behind`（逐字 `:328`），且 `cmd.Stdout, cmd.Stderr = nil, nil`（`:324`）⇒ **照抄这一形会得到一枚看不见任何输出的用例**，本问要的恰是反形。
- ⛔ 顺带一枚 CI 侧硬约束（细节见 §6）：`scripts/portable-tests.sh` 的严格跑法里**任何 top-level `--- SKIP` 都是红的**（`tools/d22scan/runtests.sh:24`"any `--- SKIP` -> SKIP is NOT a pass"，`scripts/portable-tests.sh:18`）⇒ 新用例**不许"没有桌面就 skip"**，要么判红要么按 `ledger=(` 名册（`scripts/portable-tests.sh:494-509`）具名登记。

**量格结论**：`grep -rn "objdump" scripts tools cmd internal .github 2>/dev/null | wc -l`＝**0**，`grep -rln "debug/pe" cmd internal tools --include=*.go | wc -l`＝**0** ⇒ **全仓今天没有一把读 subsystem 的仪器**（复认票 244 §9 第 3 条"会响的尺今天是真空"）。

---

## §6 问 6 —— CI 那一侧：windows job 跑哪些包、什么 tag／env

判语：**新用例会进 CI，且指得到步名与那一行**——job `test-windows`（`ci.yml:419`，`runs-on: windows-latest` `:420`，job 级 `env: WISP_ENV: test` `:421-422`）里那一步 **"cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)"＝`ci.yml:487`**，执行行＝**`ci.yml:507 run: bash scripts/wisp-cli-tests.sh`**（`shell: bash` `:505`、`if: ${{ !cancelled() }}` `:506`）。⇒ 落点链逐跳现读：`scripts/wisp-cli-tests.sh:112 bash "$portable" --scope=cli` → `scripts/portable-tests.sh:226-231 cli) … scope=(./cmd/wisp/) ; pinned=$cli_pin` → `:83 strict="$root/tools/d22scan/runtests.sh"` → `:598 sh "$strict" "${scope[@]}" -count=1 -skip "$skip_pattern"`。**只要新用例（i）在 `./cmd/wisp/` 包里、（ii）只带 `//go:build windows`、不带额外 tag、（iii）不 SKIP，它就在那一步里被跑到。**

**"step 层的绿可以等于什么都没测"这一坑：这一族今天堵住了，堵它的三行都在**：
- `tools/d22scan/runtests.sh:22`「zero top-level `--- PASS`/`--- FAIL` lines -> the pattern matched nothing」、`:24`「any `--- SKIP` -> SKIP is NOT a pass」、`:86-99` 现算 `passed/failed/skipped` 并在 `skipped>0` 时非零退出（`:99`）。
- `scripts/portable-tests.sh:18-19` 同两条规则（外加 `:19` "zero top-level PASS *and* zero FAIL is fatal (the green no-op)"）。
- ⚠ **但它不钉枚数**：`grep -rn "PASS=33" .github/workflows/ci.yml scripts/*.sh` 只命中**注释**（`ci.yml:491`"its 33 top-level cases go PASS=33 FAIL=0 SKIP=0"、`scripts/wisp-cli-tests.sh:20` 同句）⇒ **名册枚数是散文、不是仪器**；新用例进了 CI 这件事，靠的是"包在 scope 里"那一跳，不靠那枚 33。
- ⚠ 另一枚"绿但没测"的现存形状在**同一 job 的上一步**：`ci.yml:485 run: powershell … scripts/build.ps1 -Env dev`，而 `build.ps1:166-169` 的 doctor 冒烟在 GUI 子系统旗标之后**可能读到上一发 `go build` 留下的 0**（§3 末那格，⛔ 本腿无读数，判它要跑构建）。

**`winlive` 是什么门槛、默认跑不跑**：
- **默认不跑，CI 里零命中**：尺＝`grep -c "winlive" .github/workflows/ci.yml`＝**0**；`scripts/*.sh` 里也无 `--tags`／`-tags winlive`（现读 `portable-tests.sh:598` 那一发的实参只有 `-count=1` 与 `-skip`）。⇒ **今天没有任何一枚 CI 步跑过 winlive 档**，与 `docs/evidence/s1/228-resident-ball-v1.md:315` 那句"`-tags winlive` 在 `ci.yml` 里零命中 ⇒ live 档只在这台机上跑过"现读复认一致。
- 门槛本体＝构建标签：`grep -rl "^//go:build windows && winlive" cmd internal --include=*_test.go | wc -l`＝**11 枚文件**，其中 **`cmd/wisp/` 6 枚**（`panel_geometry_255_winlive_test.go`、`panel_host_windows_live_test.go`、`resident_approval_live_246_windows_test.go`、`resident_ball_live_228_windows_test.go`、`resident_hotkey_live_258_windows_test.go`、`resident_task_source_live_246_windows_test.go`）＋`internal/ball` 等 5 枚。跑法只有手敲：`go test -tags winlive ./cmd/wisp -run <名>`（现成读数见 `docs/evidence/s1/245-esc-not-a-standby-global-hotkey-r1.md:12`：每次跑 winlive 前必量 `tasklist` 的 `balldebug.exe`／`wisp.exe`／`go.exe` 三把尺＝计数 0）。
- ⇒ **这一格对最小落地形状的直接影响**：新用例**不该**带 `winlive`——它判的是 stderr 通道，不需要桌面；带上就永远进不了 CI，等于把那格残余重新交回手工腿。

**其余相关步（逐枚指到行，⛔ 说不出当它不存在）**：
- lint job（`ci.yml:65-66`，ubuntu）：`ci.yml:108` "D22 seven-ban + emoji scan (tools/d22scan)" → `:134 run: sh scripts/d22scan.sh`；`ci.yml:136` "Tracked path-length budget (ticket 262)" → `:166 run: sh scripts/check-path-length-budget.sh --with-self-test`；`ci.yml:168` "gofmt (gofumpt)" → `:171 gofumpt -l . tools/d22scan tools/mockllm`；`ci.yml:206` "go vet (module)" → `:207 run: go vet ./...`；`ci.yml:213` "staticcheck"。⇒ **新用例会过这四把静态尺**：emoji 那一把对 `_test.go` 与注释都算射程，但**注释豁免、字符串不豁免**（`tools/d22scan/main.go:1150-1162`，Q-46(c)：`commentRangesFor` 先把注释抹掉再 `emojiRe.MatchString`，判语原文"non-comment text, string literals included; comments are exempt per Q-46(c)"）；盘上尺复认：`grep -rhoP '[\x{1F000}-…]' cmd internal --include=*.go` 有 `⚠`／`⛔`／`∩` 三枚字符、20 枚文件命中，而门是绿的 ⇒ **它们全在注释里**；⛔ **新用例不许把这类字形放进任何字符串字面量**（那正是断言用的用户文案）。
- 路径长度那一把（票 262）：`scripts/check-path-length-budget.sh:149 HAT_NAME_LIMIT=100`、`:161 ISSUES_PREFIX_LEN=21`、`:156 WORST_PREFIX=44`、`:253-254 HAT_RELATIVE=121 / FULL_PATH_BUDGET=165`，名册外多一枚超预算 tracked path 即红（形 A）。现尺：`git ls-files | awk '{print length}' | sort -n | tail -1`＝**180**（＝已登记在名册里的最枚）。⇒ 新测试文件名建议按 `cmd/wisp/…_windows_test.go` 一族短名走，别去碰那枚 165 的预算线。
- ubuntu 侧**不跑** `cmd/wisp`：`test-core`（`ci.yml:309-310`）与 `portable-tests.sh` 的 core（`:206-217`）／windows（`:219-224`）两 scope 的名册里都没有 `./cmd/wisp/`（windows 那 8 枚是 `internal/proc`／`secret`／`config`／`risk`／`ball`／`perm`／`plugin`／`cmd/llmrecord`），`scripts/wisp-cli-tests.sh:65` 的 GUARD 逐字"this is the windows leg of the cmd/wisp gate"。⇒ **新用例的 CI 覆盖面只有一枚 windows runner**，这是"通道绑定"这一格未来唯一的门禁处，写派单时要具名。
- slo 两 job（`ci.yml:564-565` slo-smoke、`:622-623` slo-full self-hosted `[self-hosted, wisp-slo]`）**不跑 `cmd/wisp` 用例**，只跑 `scripts/slo-check.ps1` 消费 `wisp.exe`（`:584`/`:647`）；slo-full 那枚 runner 与票 244 §排程那格"取数期间不改构建链"直接相关。

---

## §7 问 7 —— 最小落地形状与代价：一发行文＋一枚用例

判语：**最小形成立、代价可接受，但它今天买到的是"把 257-v1 那发不可重放的手工读数换成默认档门禁里的一枚用例"，⛔ 它不满足票 244 AC#2 的字面、也碰不到 subsystem 那一维**。派单要的那句"对照 AC 与 `A619`"在这里要给准：**盘上没有 `A619` 那句**（§0 已量：`AC#1 不勾` 全台账 0 命中），等价判据的真身是**票面 AC#0 `:21`**"AC#1 单独做完会留下一枚『看着对、其实没人验过』的产物"＋"AC#1 与 AC#2 不许拆成两批改"，与 `A557:11063`"244 AC#2 真机两发（要构建＝等写面空窗）"。⇒ **按这两句读：这一枚用例既不是 AC#1、也不是 AC#2，它是 257 残余③自己的一格**，别把它记成 244 的任何一枚 AC 的兑现。

**最小形（写腿照抄即可，全部指得到现读行号）**

| 格 | 取哪一枚现成形 | 出处（现读） |
|---|---|---|
| 落点文件 | **新建** `cmd/wisp/firstrun_channel_exec_244_windows_test.go`（短名，⛔ 别去碰路径预算：§6 已量 `FULL_PATH_BUDGET=165`、tracked 最长 180 已在册） | `scripts/check-path-length-budget.sh:149/:156/:161/:253-254` |
| build tag | `//go:build windows` **单标签**，⛔ 不带 `winlive` | §6：CI 里 `winlive` 零命中＝带上就永远进不了门禁；同形先例 21 枚 `//go:build windows` 文件 |
| 编 exe | 复用 `buildWispForTest(t)`，**不要**新造 helper（`firstrun.go` 那枚"函数枚数＝1"钉不禁测试 helper，但 §5 A 组禁行 1 管包级 var） | `cmd/wisp/secret_argv_windows_test.go:161-199` |
| 起子进程 | `exec.Command(exe, "run", "<任务文本>")` ＋ `cmd.Dir = filepath.Dir(exe)` ＋ `cmd.Env = append(os.Environ(), "WISP_ENV=test", "WISP_TEST_DATA_DIR="+dataDir)` ＋ `cmd.Stdin = nil` | 三形逐枚在：argv/env `secret_argv_windows_test.go:234-236`、`cmd.Stdin = nil`（＝"无可交互控制台"那一态）`resident_sink_nail_127_windows_test.go:229`、`early_log_nail_130_windows_test.go:125` |
| 收两柄 | `lockedBuf` ×2（活进程可读）＋ `pollUntil127` 等句子，⛔ 不许 `time.Sleep` 猜时长 | `resident_sink_nail_127_windows_test.go:174-193`（`lockedBuf`/`has`）、`:317 pollUntil127(attempts int, done func() bool) bool` |
| 收尾 | `t.Cleanup` 里 `cmd.Wait()`；**不许留孤儿子进程** | 定式原文 `resident_sink_nail_127_windows_test.go:235-238`、`slo_exit_os_156_windows_test.go:328` |
| 判据（正向） | 真子进程 **stderr** 里同时有 `"新建默认配置"` 与那枚 `cfgPath`；再加一句前提：跑之前 `os.Stat(cfgPath)` 必须 `fs.ErrNotExist`、跑之后必须存在（⛔ 前提不成立就等于判了个恒真） | 要复制的判据逐字在 `cmd/wisp/firstrun_198_test.go:98-99`（同进程口径），前提形在 `:78-80`＋`257-v1:158` 那句"先 `t.Fatalf` 断文件此前不存在＋退码仍 2＋stderr 里有 `新建默认配置`" |
| 判据（退码） | `rc == 2` 不许放宽（首建之后仍按未配置失败） | `firstrun_198_test.go:86-90` 逐字："The exit code is asserted at 2 on purpose" |
| ⛔ 不许抄的判据 | `stdout == ""`（真进程必有版本行三发）；`assertDirEmpty128` 那枚合并输出形 | §5 B 组第一行；`dataroot_128_windows_test.go:131-133` |

**"一发行文"落哪一格**（三枚候选，各带归属，⛔ 本腿一枚都不写）：
1. **票 244 票面新格**（＝最合规那一枚：把"通道绑定"归进 244 名下并写明它与 AC#1/AC#2 的关系）⇒ **归编排者落笔**，且票面现在 7 枚框全未勾、本腿一枚没碰（尺：`grep -c "^- \[ \]"`＝7、`grep -c "^- \[x\]"`＝0）。
2. **台账 `A##` 一行**（真相源：`docs/reports/pending-and-issues.md`，只追加不删）⇒ 归编排者。
3. **产码注释一行**（`cmd/wisp/firstrun.go:88-95` 那段"说实话的回执"注释里加一句"这句话今天由 `cmd/wisp/firstrun_channel_exec_244_windows_test.go` 在真子进程 stderr 上钉着"）⇒ ⚠ 这一枚**有仪器后果**：`tools/d22scan` ban #9 phantom-citation 专扫 `internal/`＋`cmd/` 的**产码注释**里指向不存在路径的引用（`tools/d22scan/main.go:50-59`，且 `_test.go` 与 `tools/**` 在它射程外）⇒ **只有在新测试文件已经真落盘之后**才许写这行注释，否则写腿自己造出一枚幻影引用（本仓今天已抓到第五例幻影用例名）。

**代价**（能量到的都给了尺，量不到的具名）
- 时长：**每多一枚 `buildWispForTest` 调用＝多一发完整 cgo＋sherpa 链接**（无缓存，§1 量死：21 枚调用点＝21 发构建）。基线现成读数＝默认档整包 **148.168 s**（`docs/evidence/s1/228-resident-ball-v1.md:13`），台账 `A621` §7 的排程口径"整包预期 ~8.5 min"。⇒ **精确增额本腿量不到**（要跑构建），能说的是它落在"与 127/246 那族同一枚价"这一档，不是新增量级。
- 毒源继承：走 `run` 腿 ⇒ **零** `proc.Boot`、不起球窗、不抢全局热键、不吃每会话互斥（§5 E 组那把尺）；⛔ 别顺手复用 `bootResidentLeg`，那会把 127 那族带载偶发一并接进来。
- 排程冲突（写腿前必查）：`268-r1` 此刻正在 `cmd/wisp` 写（`A621` §7 逐字），同包并发＝互洗读数；票 244 自己那格"排程与串行"也写过同一条形（⛔ 与 `228-r1` 串行）。⇒ **新用例排在 `268-r1` 交回之后**。
- CI 面：只在 `test-windows` 那一步生效（§6），ubuntu 两 job 永远不跑它 ⇒ 这一格的回归网只有一枚 runner 宽。

**买到什么／买不到什么**（逐格对照票面 AC）
- ✅ **买到**：257-v1 `:150`/`:161` 那句"生产把它绑到 `os.Stderr` 这一步没有任何用例钉着"从〔手工一发、⛔ 不可重放为门禁〕升级成〔默认档门禁里的用例〕；顺带让"`wisp run` 的 stderr 通道"在 `run` 腿上有了一枚能被打红的钉（ Mutation：把 `firstrun.go:92` 的目标换成 `io.Discard`，新用例必红——这正是 `257-v1` §8 M1 那一发的自动形）。
- ⛔ **买不到 1**：**AC#2 的字面三发真机读数**（票面 `:23`）——①"在父控制台里跑、回复文本真出现在那个控制台"要的是**人眼看那台控制台**，用例给的是"管道里看得见"；②"`wisp doctor > out.txt` 重定向仍有效"是另一枚 argv；③"双击／explorer 拉起不再出现黑框＋常驻腿照常起"§4 判语 2 已具名：**那一支根本不产生这句话**。⇒ **⛔ 不许用这枚用例去勾 AC#2**（票面 `:23` 原文那句"⛔ 不许用『单测里 mock stdout』代替真机"针对的是 mock；真子进程用例不是 mock，但它也**不是** AC#2 那一发的替身）。
- ⛔ **买不到 2**：**subsystem 这一维**。`buildWispForTest` 零 `-ldflags` ⇒ 台件二进制永远 CUI（§1/§3），而票 244 §5 G5 那句"包内断言天生看不见 subsystem 这一维"就是为这一格写的。要钉 subsystem 得另一枚尺（读 `debug/pe` 或 `objdump`，且**必须先有构建**，因为库里 tracked exe＝0 枚，§3 现尺）——归 AC#1/AC#5 ⓑ，不归本问。
- ⛔ **买不到 3**：**态乙那一发**（GUI 件无父控制台时那句话的真去向：NUL 还是无效柄）——量不到，归编排者真机（票面 AC#0 与 §9 第 2 条两处都禁止继承 `docs/BUILD.md:90` 那枚过期读数）。

---

## §8 落盘尺与「⛔ 这轮没动的东西」自证

**七问判语一句话版**（正文在各自那节，每节都带现尺读数）：

| 问 | 判语 | 主要那把尺（读数） |
|---|---|---|
| 1 现成先例 | 这一族 10 枚文件／21 枚调用点，⛔ 没有一枚是"编译一次"（`buildWispForTest` 零缓存、零 `-ldflags`）；判过"某句话到真子进程 stderr"的只有 4 枚、且全在 slog-mirror 那一支 | `grep -rn "buildWispForTest(t)" cmd/wisp \| wc -l`＝21；`grep -rn "exec.Command(" cmd/wisp \| wc -l`＝23（**漏计 5 枚真起子进程的文件**，含 `dataroot_128` 那枚 `CommandContext`） |
| 2 通道绑定 | 绑定在 `main.go:159-164`；`cmdRun` 零枚测试调用者；`runTextTask` 的 8 枚测试调用点**八枚同进程传 buffer**；`run198` 那枚共享 helper 15 处命中、零枚走子进程 ⇒ **派单那格残余成立** | `grep -rn "cmdRun(" cmd/wisp --include=*.go`＝2（定义＋dispatch）；`grep -rn "run198(" cmd/wisp \| wc -l`＝15 |
| 3 subsystem | 旗标无条件在 `build.ps1:115`；盘上出厂件仍 **CUI(3)**（比旗标早 2 天、比 `run.go`/`firstrun.go` 旧）；244-r1 那枚临时出厂件真读 **GUI(2)**；"今天重跑 build.ps1"那一发**量不到，归编排者**＝AC#1 欠的那一发 | `objdump -p build/wisp.exe \| grep -i subsystem`＝`00000003 (Windows CUI)`；`git ls-files "*.exe" \| wc -l`＝0 |
| 4 无父控制台 | attach 永远早于绑定（无错位）；★**双击那一支根本不产生这句话**（无参＝resident 分支，从不进 `cmdRun`）；`console_windows.go` 里"没有父控制台就 return"这句**不存在**（早退条件是"已有可用 stdout"`:34-37`） | 现读 `cmd/wisp/console_windows.go:33-43`、`cmd/wisp/logsink.go:160`、`cmd/wisp/run.go:222/:245` |
| 5 撞钉预检 | 会红的只有两枚 AST 门的**仪器形**（三枚禁行写死）；12 枚计数形钉都打不红；★反向陷阱＝`firstrun_198_test.go:91` 的"stdout 必须为空"在真进程口径天然不成立；★新用例走 `run` 腿**不继承真窗毒源**；"起了子进程却没收尾"名册里**零枚** | `grep -rn "t.Parallel()" cmd/wisp --include=*_test.go \| wc -l`＝0；`grep -rn "objdump" scripts tools cmd internal .github \| wc -l`＝0；`grep -rln "debug/pe" cmd internal tools --include=*.go \| wc -l`＝0 |
| 6 CI | 会进 CI，指到步名＋那一行：`ci.yml:487`＋`:507` → `wisp-cli-tests.sh:112` → `portable-tests.sh:226-231` → `runtests.sh`（零 PASS 即 fatal、任何 SKIP 即 fatal）；⚠"33 枚"只在注释里；**`winlive` 在 ci.yml 现量 0 命中＝默认不跑**，新用例不许带这枚标签 | `grep -c "winlive" .github/workflows/ci.yml`＝0；`grep -c "timeout-minutes" .github/workflows/ci.yml`＝0；`grep -rn "PASS=33" .github scripts/*.sh`＝3 处、**全 `#` 注释** |
| 7 最小形 | 形状齐（tag／helper／argv／两柄／收尾／判据逐格都指得到现读行）；买到＝把 257-v1 那发手工读数升级成默认档门禁；⛔ 买不到 AC#2 三发真机字面、买不到 subsystem、买不到态乙 | 基线两把现成读数：整包 148.168 s（`228-v1:13`）／排程口径 ~8.5 min（`A621` §7）；精确构建增额**量不到** |

**推翻／更正派单的五处**（逐枚具名，⛔ 没照抄进表）：① "编译一次 wisp.exe 再用 exec.Command 起真子进程"——**不存在缓存形**；② `console_windows.go:38-51` 那句"没有父控制台就 return"——**方向相反**；③ `A619` 那句"AC#1 不勾＝缺 AC#2 真机两发"——**归错账**（真身＝票面 AC#0 `:21`＋`A557:11063`，`AC#1 不勾` 全台账 0 命中）；④ 起手尺 `grep -rn "exec.Command(" cmd/wisp`——**漏计 5 枚**（`exec.CommandContext` 与复用 127 helper 的四枚）；⑤ 隐含前提"exec 级用例能把 AC#2 双击那一发钉住"——**买不到**（§4 判语 2）。

**落盘尺**（收尾现跑，`85c7e51e` 之后复量）
- `wc -l .scratch/wisp/probes/244/a4/census.md`＝**327**／`wc -c`＝**62,270**。
- 占位尺 `grep -cE '待[填]|填写[中]|未判'`＝**1**，而那唯一一行是本节自指（`grep -nE` 现读＝本行，落在"占位尺"那一枚句子里）；把它自己剔掉再量 `grep -E '待[填]|填写[中]|未判' \| grep -vc "占位尺"`＝**0** ⇒ **真实占位 0 枚**。对照骨架那一笔（`git show 05ac9fc0:… \| grep -cE …`＝**10**）⇒ 七问＋§0＋§8 十一格里，本腿把 10 枚占位全部答完。
- ⚠ 自指这一行永远会被那把尺捞到，**下一位复量时请一并复跑那条排除式**，别把它读成一格没答。
- commit 链（逐枚 `git log --format="%h %cI %s" -- .scratch/wisp/probes/244/a4/census.md` 复跑）：`05ac9fc0` 13:01:41 骨架 → `625717d9` 13:08:43 §0＋§1 → `7b5a6d13` 13:11:12 §2 → `273cb2d5` 13:17:01 §3 → `8734a4ba` 13:18:24 §4 → `e00cbd6d` 14:08:41 §5＋§6 → `b70cf069` 14:30:17 §7 → `85c7e51e` 14:35:48 §8 →（末一笔＝本行终尺更正）。

**具名"量不到"清单**（⛔ 全部用推测顶过；每格都写了归谁）

1. **"今天再跑一次 `scripts/build.ps1` 会得到 GUI(2) 的那一发读数"**——量不到（判它要跑 `go build`，本腿禁）；归编排者＝票 244 **AC#1** 欠的那一发（§3）。
2. **"台件 `wisp.exe` 的 subsystem 实测"**——量不到（同上，且 `git ls-files "*.exe"`＝0 枚 tracked exe 可偷读）；只给"零 `-ldflags` ⇒ 必然 CUI"的〔读码〕判语（§3 表末行）。
3. **"GUI 件无父控制台时那句 `fmt.Printf` 的真去向（NUL 还是无效柄）"**——量不到（要真机双击一发 GUI 件）；归编排者＝票面 **AC#2 ③**，且 `docs/BUILD.md:90` 那枚过期读数**不可继承**（§4 态乙）。
4. **"`build.ps1:166-169` 冒烟是否已因 GUI 旗标恒绿"**——量不到（判它要跑构建）；只登记为票面 **AC#6** 名册里一枚待答消费者（§3 末）。
5. **"新用例的精确构建增额（秒）"**——量不到（要跑构建）；只给两把现成基线读数（整包 148.168 s／排程口径 ~8.5 min）（§7 代价）。
6. **"CI 那台 windows runner 有没有可建窗口的桌面"**——本腿量不到（`gh run` 不在授权内、也无读数）；本问不依赖它（新用例不该带 `winlive`），⛔ 不许据此句判任何事（§6）。

**⛔ 这轮没动的东西**（自证，全部现跑）

- 逐笔 `git show --name-only` 复跑：上面 **8 枚 commit 每一枚的文件清单只有 `.scratch/wisp/probes/244/a4/census.md` 一行**，没有第二枚文件。
- `git status --porcelain -- cmd internal tools scripts docs`＝**1 行**：` M cmd/wisp/resident_approval_windows.go`——**那是 `268-r1` 的活，不是我的**（本腿对它零写入、只读定位过它的存在，⛔ 未改、未 stash、未 add）。⇒ 我这轮在该四界里贡献 **0 行**。
- **一枚 `go` 命令都没跑**：无 `go test`／`go build`／`go vet`／`go env`／`go list`（唯一执行过的二进制是 `objdump`、`git`、`grep`、`ls`、`wc`、`find`、`command -v`）。`268-r1` 在 `cmd/wisp` 取整包读数的名册因此没有被本腿洗过。
- 票面／台账／证据件／`docs/**` **一字未改**；票 244 的 7 枚 `- [ ]` 框一枚没碰（收尾复尺 `grep -c "^- \[ \]" .scratch/wisp/issues/244-spec-11-wants-the-gui-subsystem-build.md`＝7、`grep -c "^- \[x\]"`＝0）。
- `frontend/**`／`design/**` **零读零写零转述**（本腿所有 grep／find 的根都显式限定为 `cmd internal tools scripts docs .github .scratch`；尺：`find` 只从 `.scratch` 起、`grep -r` 全部带显式目录参数）。
- temp 件只建不删，全在 `.scratch/wisp/probes/244/a4/`（本文件的 8 枚 commit message txt 都落在同目录，⛔ 没进任何一次 pathspec）。
