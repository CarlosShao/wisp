# 244-c1 — 票 244「切 GUI 子系统」前置普查（只读第一趟；零产码／零构建／零 go test）

- 腿：`244-c1`（派单存档 `.scratch/wisp/probes/orch/20261102/txt/244c1-dispatch.txt`，2026-11-02 10:4x）
- 票面：`.scratch/wisp/issues/244-spec-11-wants-the-gui-subsystem-build.md`（blob `9b80604d`，109 行，全文已读；**七枚 AC 框一枚未碰**，收尾复尺见 §6）
- 起手三核（10:4x 现跑）：`git rev-parse HEAD`＝**`c9ef940c853d417d528e0c47e1b0887fcac6cf14`**（≥ 派单要求 `c9ef940c8`，相等）；`git diff --cached --stat`＝**空**（暂存区干净）；工作树脏项（`M cmd/wisp/run.go`＋数百枚未跟踪）全部属 337r2 写腿与编排者暂存件，**本腿一枚没碰**。
- 收尾锚点（11:0x）：HEAD 仍 `c9ef940c8`；`git diff --name-only HEAD -- scripts/build.ps1 scripts/slo-check.ps1 .github/workflows/ci.yml cmd/wisp/console_other.go docs/BUILD.md .scratch/wisp/probes/244/a2/census.md`＝**空**（本腿读过的文件在取数期间与 HEAD 同号，按 a2 腿 §⑦7.14 的口径 blob＋diff 双尺）。取数时刻 `2026-11-02 10:46–11:02 +0800`。
- 记法：`〔量〕`＝本腿本机现跑；`〔推〕`＝读码/引用推断；`〔263量〕`＝票 263 证据件在该机器上实测过的同形读数（本腿未重跑）。

---

## 1. AC#0① 现跑：产物子系统读数——**派单预测"应=00000003"与真读数相反，切换已经落地**

尺（本腿现跑，binutils objdump＝`E:\work\base\msys64\mingw64\bin\objdump.exe`）：

```text
objdump -p build/wisp.exe | Select-String Subsystem      # 即票面 grep -i subsystem 的等价尺
```

真实读数〔量〕：

```text
MajorSubsystemVersion	10
MinorSubsystemVersion	0
Subsystem		00000002	(Windows GUI)      ← 不是 00000003 (Windows CUI)
Characteristics 0x26
DllCharacteristics 00008160
```

产物身份具名（`build/**` 被 `.gitignore`＝不可再生，票面 AC#1 要求自带产生命令；本腿⛔动构建链，只能报"读到了什么"）：

- `build\wisp.exe`＝**31,555,449 字节／mtime 2026-10-15 18:18:31**〔量〕。
- 对比票面时间线：票面现量表（09-30）＝29,771,385 字节／CUI(3)；票 §9 第 2 条（10-03）＝"mtime 09-30 23:43、比 244-r1 早、GUI 产物落在临时目录没覆盖 build/"。⇒ **今天盘上这枚既不是票面那枚也不是 §9 那枚——10-15 有人真产出过一枚 GUI 子系统产物**。谁、用哪条命令产出，不在本腿射程（无命令史可考）〔量不到，具名〕。
- objdump 全文落 `D:\tmp\wisp\20261102\244c1\objdump-wisp-exe.txt`〔量〕。

**结论（读数，不是建议）**：派单/票面 AC#0 的"CUI(3) 预期"与票面 AC#1 的"改前读数 CUI(3)"在今天都**已过期**；票 §9（10-03 收 244-c2）那句"出厂链 3 枚已带／继承 `-H=windowsgui`"与本读数一致。票面 AC#1"改后读数 GUI(2)"那一发——**就这枚产物而言已成立**，但票面 §9 第 2 条同时警告过"旗标已落≠台上 exe 没黑窗"；今天两半都落齐了（链带旗＋台上产物＝带旗产物），唯缺一发**由 scripts/build.ps1 当场产出并同场读子系统**的记录（`build/**` 不可再生 ⇒ 见 §4 名册第 1-4 行的连带面）。

## 2. AC#0② 构建链读数：派单说的 `:103-110` 是过期行号，现状在 `:244-257`

现读 `scripts/build.ps1`（blob `29df41f6…`，313 行）〔量〕：

- `$ldflags` 数组现居 **`:249-257`**（`:103-110` 今天被前端步骤的注释与 `Fail` 文案占据——票面行号自 `3593e87b` 时代腐烂，与 §9 第 1 条"`:115`"一样是旧尺读数）。
- 数组**共七枚**：六枚 `-X` **一字未动**（`:250` Version／`:251` Commit／`:252` BuildDate／`:253` DefaultEnv／`:254` SherpaOnnxVersion／`:255` OnnxRuntimeVersion）＋ **`:256` `"-H=windowsgui"`**。
- `:244-248` 是落地时写的 why 注释（SPEC-11 §2.2／attachParentConsole 早退语义），与票 §9"上面 5 行 why 注释从 `:103` 起、六枚 `-X` 一字未动"逐字对得上，只是行号又随前端步骤（票 274）整体下移。
- 构建行 `:270`＝`& $go build -trimpath -ldflags $ldflags -o build\wisp.exe ./cmd/wisp`；冒烟 `:307-310` 见 §4 名册第 1 行。

**全仓 `-H` 面复尺**〔量〕：`grep -rn -- '-H=' scripts/ tools/ .github/`——除 `build.ps1:256/:244/:157-158(注释)`、`slo-check.ps1:58/:158(注释)` 外唯一命中仍是二进制 blob `scripts/spike/bin/goja-caps.exe`（票面/a2 S14 同读数）。

旁读（边界内）：`cmd/wisp/console_other.go:7`〔量〕那句绝对句「The Windows build links as a GUI-subsystem binary (ticket 07)」——**旗落之后这句从假话变成了真话**（归口仍是票面 AC#3 的"带条件的事实句"改写，本腿不动它；它举的出处"ticket 07"依旧错，切换实际由票 244 线落地）。

## 3. AC#5ⓐ GUI 进程的任务源从哪来——引用 `244-a2` 普查件（在），并记录它的时效边界

`.scratch/wisp/probes/244/a2/census.md` **存在**〔量：glob＋全文读〕（211 行／71,898 字节，锚点 10-01 `f1e7a3e2→697b4fae`；其读数按票 §5 编排者收件已被复认过一轮）。本腿依派单走引用路线，不重跑 cmd/wisp 尺（边界＝cmd/wisp 只许读 `console_other.go`，见 §5）。

a2 现量（转引，标 a2 时刻）：

- `submitTask` 产码调用者总数＝**2**（`cmd/wisp/resident_task_source_windows.go:322` 注入分支／`:363` 控制台 `task <文本>` 动词）；两枚可达性分别挂在「`env==test`＋三把锁」与「真控制台缓冲（`interactiveStdin()`）」上 ⇒ **prod＋双击 GUI 一形，两条同时归零**（a2 §①.1 A 表 + §⑦7.3）。
- 候选甲（面板＝票 33＋248）到 a2 收钟态＝**宿主文件已进 HEAD（`697b4fae`）、`dispatchRaw` 未装 `Message`、常驻进程零 `PanelManager` 装配、"面板能举一发任务"＝0**（a2 S30/S31）。
- 候选乙（托盘）＝常驻腿真执行者 0（a2 S25）；热键切后仍活但只有 cancel 一枚、且是否决不是举任务（a2 S26/§①.1 丙行）。

**本格答案（=票 §5 编排者已裁的口径，本腿不另选）**：GUI 进程任务源＝**裁甲，面板路（票 33＋248）**；停机源（AC#5 ⓐ′ 第二腿，§5 G2 补格）＝**乙-1＝`internal/proc` stop-request 位，归口票 228 新 AC#11**。"允许"那一半与任务源正交：面板侧结构上永无 `Allow`（`ui.go:164-171`，票 §5 修法禁区原文），原生侧归票 228 AC#2（托盘第五枚，§5 甲格裁甲-4）。
**落地状态读数〔量不到，本腿具名〕**：10-01（a2）之后到今天（11-02），面板路是否已长出"能举一发任务"的装配、228 AC#11 是否落地——**超出本腿只读边界**（cmd/wisp 仅一枚文件许读），交下一枚持 cmd/wisp 读权的腿或编排者复量；本腿不写"已/未落地"任何一端。

## 4. AC#6 名册：谁在消费 `wisp.exe` 的退出码／stdout（逐条 file:line＋"翻成 GUI 后还成立吗"）

尺〔量〕：`grep wisp\.exe → scripts/**`（27 命中，全列于 build.ps1／slo-check.ps1 两文件，`.sh` 零命中）；`grep 'wisp|build.ps1|slo-check.ps1' → .github/workflows/`；`grep 'LASTEXITCODE|Start-Process|WaitForExit|-Wait' → scripts/**`（20 命中逐枚定性）；docs 面 `grep 'wisp\.exe.*>' → docs/**`。

| # | 消费者 | file:line（现行号，blob 锚） | 消费什么／怎么起 | 翻成 GUI 后还成立吗 |
|---|---|---|---|---|
| 1 | **`scripts/build.ps1` 冒烟** | `:309-310`〔量〕`& (build\wisp.exe) doctor` → `if ($LASTEXITCODE -ne 0) { Fail }` | `&` 起产物＋读自动变量 | **不成立⚠本票头号残量**。〔263量〕同形实测（`slo-check.ps1:157-164` 逐字记录：CUI 探针 2049ms 带 rc＝3 返回，GUI 探针 **70ms 放空、下一次取值在严格模式下死**）⇒ GUI 子系统下 `&` 不等待；本文件 `:40-41` 有 `$ErrorActionPreference='Stop'`＋`Set-StrictMode 2.0`，但 `$LASTEXITCODE` 已被 `:270` 的 `& go build`（rc=0）赋过值 ⇒ `:310` 读到的是**前一道命令的陈值 0 ⇒ 冒烟自 -H 落地起结构性空转（doctor 真 FAIL 也不会红）**。〔推，本腿⛔动构建链未现跑；机制凭据＝263 实测〕。**票 263 只修了 `slo-check.ps1`，没扫这枚同形姊妹文件**——slo-check 自己的形状钉也自陈射程只有本文件（`:233-247`，`$path = $PSCommandPath`）。⇒ 具名交编排者裁归口（本腿⛔改，AC#6 判据="只建名册不改任何消费者"）。 |
| 2 | **`scripts/slo-check.ps1`** | `:533` 逐态 `slo -state… -out`／`:557` settle／`:583` leak——全走 `Invoke-ExternalProgram`（`:143-200`，`Process.Start`＋**`WaitForExit()`**＋真 `ExitCode`）；`:307` 起 `build.ps1` 同形；`:292` 跑前能力钉（cmd /c exit 7）＋形状钉（本文件字面禁自动变量） | 退出码＝D32 门的判据主体 | **成立**（263 已修，本腿现读复认〔量〕：三发＋build 路零处 `&`＋自动变量直读）。钉的自陈残口〔量：`:240-246` 原文〕＝运行期构名取值／别文件的直读它不抓，外层= `scripts/slo-freshness.sh` P3。stdout 面：三发都带 `-out <file>`，报告落盘〔244票面§⑦7.12 同读数〕。 |
| 3 | **CI build 步骤 ×3** | `ci.yml:847-852`「cgo build smoke」／`:1115-1117`「Build wisp.exe」／`:1192-1193`「Build wisp.exe (deps cached)」，全＝`powershell -File scripts/build.ps1 -Env dev` | 只消费 **build.ps1 的退出码**（步骤级），不直接起 `wisp.exe` | **成立（作业起法）＋继承第 1 行的盲**〔量〕：旗标全在 build.ps1 内，CI 无第二处 flag 面（a2 S16 复认）；但三步的 `doctor` 冒烟绿就是第 1 行那种绿。 |
| 4 | **CI SLO 步骤 ×2** | `ci.yml:1119-1121`「SLO smoke gate」／`:1195-1196`「SLO full gate」＝`powershell -File scripts/slo-check.ps1 …` | 消费 slo-check.ps1 退出码（其内部再逐发消费 wisp.exe 码） | **成立**〔量〕：继承第 2 行的已修形；`:1230` `if-no-files-found: warn` 是 134 AC#6 的装饰变更，与本面无关。 |
| 5 | `slo-fresh.yml` | 全文件 `wisp` 仅 `:8` 注释一枚命中〔量〕 | 不起 `wisp.exe` | **非消费者**（a2 S16 复认成立）。 |
| 6 | **人手文档行** | `docs/BUILD.md:90-92`（`:92`＝`wisp.exe version > out.txt`）〔量：现读禁改件，一字未动〕 | 人跑 CLI 腿的 stdout／重定向 | **成立，且有真机凭据**〔量外引用〕：`docs/evidence/s1/244-black-console-r1.md:95-96`——244-r1 用**带 -H 的自建 GUI exe** 实测 `run --help > out 2>&1` rc=2／1782 字节、`version` rc=0／267 字节、`-h` rc=0／2464 字节 ⇒ 重定向与父控制台两态下 stdout 都接得住（机制＝`console_windows.go` 有合法 stdout 即早退不 attach，与 build.ps1 `:246-248` 注释一致）。禁改件该段仍写"临时构建验证"＝过期自述，归口票面 AC#3/G6 一族，本腿只登记。 |
| 7 | 编排者清理令 | `docs/reports/orchestrator-handoff-20261015.md:113` `taskkill /F /PID <wisp.exe>`〔量：grep 命中〕 | 消费进程存在性，**不消费退出码** | 成立／与子系统无关（taskkill 不等不收码）。登记它因它点名杀常驻 `wisp.exe`——GUI 化不改变可杀性。 |
| 8 | **确认非消费者**（防"名册漏计＝我漏的"） | `scripts/dev/ball-cycle.ps1:31/:108`（起 `balldebug`，且已用 `WaitForExit(15000)`）〔量〕；`scripts/spike/run.ps1:27/:33/:51/:71`（起 go/gcc/spike bins，非 wisp.exe）〔量〕；`scripts/{fetch-deps,sign-models}.ps1`（无 wisp）〔量〕；`scripts/*.sh` 对 `wisp.exe` **零命中**（wisp-cli/portable/winsec 三 tests.sh 全 `go test` 族——a2 S18 转引＋本腿 grep 复认〔量〕）；Go 台件自 build 五处（`secret_argv_windows_test.go:174` 等，**无 `-ldflags` ⇒ 台件二进制永远 CUI，读不到出厂产物这一维**——a2 S22＋票 §9 第 1 条复认〔转引〕，另 §9 补计 `internal/llm/openaichat/mockllm_integ_test.go:65`） | — | 均不消费出厂 `wisp.exe`；其中台件一族＝AC#5ⓑ"仪器看不见这一维"的名册外证据，转引不重跑。 |

**名册口径与量不到（具名）**：① 射程＝`scripts/**`＋`.github/**`＋`docs/**` 的字面尺与票面/§9/a2 的转引复认；**人手会话里临时 `& build\wisp.exe …` 读 `$LASTEXITCODE` 的形状枚举不了**（写"量不到"）；② `build/wisp.exe` 10-15 那发的产生命令史——不可再生、无记录，写"量不到"；③ 第 1 行的"空转已被触发过几次"要看 CI run log，本腿未取（可取，非判据必需）；④ `wisp-cli.exe` 只在 `slo-check.ps1:363` 进程名清单里被**认名**，全仓无产它的步骤——不列消费者，登记于此防误删。

## 5. 边界合规与暂存件

- 只读面：`scripts/**`、`docs/**`、`.github/workflows/**`（AC#6 判据点名"CI 步"）、`cmd/wisp/console_other.go`（唯一一枚，按派单点名）、票面与 a2 普查件。**未读未写** `frontend/**`／`design/**`（零转述——本件对前端的唯一提及＝build.ps1 注释里的事实句引用其行号，不含内容转述）。
- ⛔ 零动码／构建链／票面／台账：未跑任何 `go build`/`go test`/`build.ps1`；`M cmd/wisp/run.go`（337r2 在写）与全部未跟踪项原样未碰；无 push。
- 临时件目录＝`D:\tmp\wisp\20261102\244c1\`（`objdump-wisp-exe.txt`，只建不删）。

## 6. 收尾复尺（交件前最后一发）

- `git rev-parse HEAD`＝`c9ef940c8…`〔量，与起手同号：期间编排者未落新提交〕；本件读过的六枚文件 `git diff --name-only HEAD`＝空（见头部收尾锚点段）。
- 票面 AC 框〔量〕：`grep '^- \[ \]'`＝**7**／`grep '^- \[x\]'`＝**0**——本腿前后同数，一枚未碰。
- 本腿唯一提交＝本文件（显式 pathspec，不 push）。

## 7. 读数一句话汇总（给编排者收件用）

1. **AC#0①**：`build/wisp.exe` 子系统今天＝**00000002 (Windows GUI)**，派单预期"=3"过期；产物 mtime 2026-10-15 18:18:31／31,555,449 字节，链上 `scripts/build.ps1:256` 现居 `-H=windowsgui`（六枚 `-X` 在 `:250-255` 一字未动，a2/票 §9 的 `:103-110`/`:115` 都是腐烂行号）。
2. **AC#5ⓐ**：依 a2＝任务源归票 33＋248 面板路（裁甲）、停机源归票 228 AC#11（乙-1）；"允许"永远走不到面板侧（结构禁令原文在票 §5）。11-02 的落地态本腿量不到（边界具名）。
3. **AC#6**：名册 8 行；**唯一新增危险读数＝`build.ps1:309-310` 的 `&`＋`$LASTEXITCODE` 冒烟在 GUI 产物上是结构性空转**（263 修了 slo-check、没修这枚姊妹文件，slo-check 的形状钉射程自陈只有本文件）——只登记不改，归口交编排者裁。其余具名消费者（slo-check 三发／CI 五步／BUILD.md 人手行／taskkill 行）翻 GUI 后成立，凭据逐行带锚。
