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

未判。

---

## §3 问 3 —— subsystem 那一格：`-H=windowsgui` 在 `scripts/` 的哪一行、PE subsystem 真身

未判。

---

## §4 问 4 —— 无父控制台那一支：`console_windows.go` 的 return 与 `logsink.go` 的 mirror 谁先谁后

未判。

---

## §5 问 5 —— 撞钉预检：新增一枚 exec 级用例会撞哪些既有计数尺

未判。

---

## §6 问 6 —— CI 那一侧：windows job 跑哪些包、什么 tag／env

未判。

---

## §7 问 7 —— 最小落地形状与代价：一发行文＋一枚用例

未判。

---

## §8 落盘尺与「⛔ 这轮没动的东西」自证

未判。
