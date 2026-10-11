# 304-a1 — 表①：本仓"测试里起子进程"的全名册

## 0. 尺名册（每把都写射程目录＋枚数来源）

| 号 | 尺（逐字） | 射程 | 枚数 | 量来／算来 | 件 |
|---|---|---|---|---|---|
| R1 | `grep -rEn "exec\.(Command\|CommandContext)\(" --include=*_test.go internal cmd tools` | 三枚顶层 Go 目录的 `_test.go` | **66 枚 spawn 语句** | 量来的（`wc -l`） | `logs/grep-m1-testonly.txt` |
| R2 | 同 R1 但 `--include=*.go`（⛔ 限 `_test.go`） | 同上三目录全部 `.go` | **72 枚** ⇒ 比 R1 **多 6 枚＝非测试文件里的 spawn** | 量来的 | `logs/grep-m2-allgo.txt` |
| R3 | `grep -rEn "\.Start\(\)" --include=*_test.go internal cmd tools` | 同 R1 | **17 枚** ＝"父进程不等它"那一档 | 量来的 | `logs/grep-m4-start-in-tests.txt` |
| R4 | `grep -rEn "StartInJob\(" --include=*_test.go …` | 同 R1 | **3 行＝2 枚调用点**（`:41` 真起／`:207` 预期起不来）；`:42` 是 `t.Fatalf` 里的字样，⛔ 调用 | 量来＋算来 | `logs/grep-m8-startinjob-in-tests.txt` |
| R5 | `grep -rEn "\.(Output\|CombinedOutput\|Run)\(\)" --include=*_test.go …` | 同 R1 | **51 行**（＝同步形的语句行；比"R1−R3"多两行，因为有的件里同一把尺还命中了别家的调用形 ⇒ **只当"同步形有哪些行"用，⛔ 当分母**） | 量来的 | `logs/grep-m7-sync-forms.txt` |
| R6 | `grep -rEn "os\.StartProcess\|SysProcAttr\|go build -o" --include=*.go internal cmd tools` | 同 R2 | **3 枚**（全是 `SysProcAttr`，`os.StartProcess` **0 枚**） | 量来的 | `logs/grep-m3-other-shapes.txt` |
| R7 | `find . -type d -name testdata` ＋在那些目录里打 `exec.Command` | 全仓 testdata | testdata 目录 **6 枚**；testdata 里 spawn **0 枚（rc=1）** | 量来的 | 本件 §3"看不见什么" |
| R8 | 全仓 `exec.Command` 减 `internal cmd tools`（`grep -rEn "exec\.Command" --include=*.go . \| grep -v …`） | `.`（含 `.scratch`、`scripts`） | **17 枚**：`.scratch/wisp/probes/**` 15 枚＋`scripts/spike/webview2-latency/main.go` 2 枚 | 量来的 | `logs/grep-m5-outside-scope.txt` |
| R9 | `grep -rEln "exec.Command" --include=*_test.go internal cmd tools`（去重到文件） | 同 R1 | **41 枚测试文件** | 量来的 | `logs/grep-m6-distinct-testfiles.txt` |

**枚数是整族还是抽样＝整族**：R1/R2 的每一枚都逐枚进了本件（R2 的 72＝R1 的 66＋非测试的 6，作差＝量来的两条尺相减，另逐枚点清那 6 枚见 §2）。分档按"父进程等不等它"：`66＝17（R3 的 Start 形）＋1（R4 的 `:207`，注定起不来）＋48（同步 inline 形）`——**这 48 枚的同步形逐枚读过**（每件都翻了 `R5` 命中的那一行与 spawn 行的相邻关系，⛔ 一把正则推整族）。

provenance＝**工作树**（起手 `git status --porcelain -- internal cmd tools`＝**0 行** ⇒ 工作树＝`d75db9ee` 的 blob，两态同值这一条是量来的）。

## 1. 会留下活子代那一档（18 枚，逐枚：文件:行｜起什么｜回收挂在哪｜外力 `TerminateProcess` 时子代命运）

★**这一档就是本票的缺陷面**：`t.Cleanup` 只在测试二进制**自己走到收尾**时跑；`defer` 只在 `panic`／`Goexit`／正常返回时跑——**外力 `TerminateProcess` 三种都不跑**（＝票面 `:5` ③ 与派单 §1 特别点名那一类）。最后一列是**读码级**判断（依形状＋票面机制段），⛔ 实测；实测归表②（本腿⛔ 做）。

| # | 文件:行（spawn｜Start） | 起的是什么 | 回收挂在 | 外力杀父时（读码级） |
|---|---|---|---|---|
| 1 | `internal/llm/openaichat/mockllm_integ_test.go:73`｜`:80` | **`mockllm.exe` 假大脑（HTTP 服务器，持 Listen 套接字）** | `t.Cleanup` Kill＋Wait `:106`；就绪失败两枚分支体内 Kill `:99`/`:103` | **留着不退**（服务器⛔ 自退）＝票面那两枚孤儿的形；⚠ 回收注册在就绪等待**之后**（对比 #17 那一枚注册在之前） |
| 2 | `internal/memory/concurrent_test.go:230`｜`:236` | 自 re-exec 测试二进制（`-test.run=^TestSubprocessCrashWriter$`，名册已核＝decl 在 `:184`） | ⛔ `t.Cleanup`／⛔ `defer`：**体内内联** Kill `:262`＋Wait `:265` | 留着，直到它自带的 `-test.timeout=60s` 自己炸（＝**有限**残留） |
| 3 | `internal/proc/jobscope_windows_test.go:33`｜`:41`（有 job）或 `:45`（无 job） | 自 re-exec `TestHelperProcess`（decl `:59`，名册已核） | **Job Object**：`t.Cleanup(job.Close)` `:111`/`:161` ＋ OS 的 `KILL_ON_JOB_CLOSE` | **这一枚是例外**：OS 关句柄即端子代＝甲形在本仓已在跑的活样本（见 `30-costs.md`） |
| 4 | `internal/proc/singleinstance_windows_test.go:30`｜`:40` | 同上 helper（`single-instance-holder`） | `t.Cleanup` Kill `:43`（⛔ Wait） | 留着，直到 helper 自退（⛔ 事件时 30s，`:80-82` 那段） |
| 5 | `internal/tools/bridge_a18_kill_windows_test.go:107`｜`:114` | 自 re-exec（`-test.run=^TestA18RealTaskkill…$`，decl `:176` 已核） | **`defer reap()` `:129`**（真 `taskkill /F` `:121`＋Wait `:126`） | `defer` ⛔ 跑 ⇒ 留着，直到 `-test.timeout=90s` |
| 6 | `internal/tools/fs_edit_ac4b_kill_windows_test.go:125`｜`:131` | 自 re-exec（`^TestFSEditRealTaskkill…$`，decl `:176` 已核） | **`defer func(){…}` `:136`**（taskkill `:138`＋Wait `:143`） | 同上：留着到 `-test.timeout=90s` |
| 7 | `internal/tools/fs_staging_windows_test.go:51`｜`:52` | `cmd.exe /c exit 0`（秒退夹具） | 体内内联 Wait `:57` | 已经退了，⛔ 残留面 |
| 8 | `internal/tools/fs_staging_windows_test.go:383`｜`:388` | 自 re-exec 持句柄者（`^TestSweepSparesATempFileHeldByAnotherProcess$`，decl `:310` 已核） | **调用方 `defer` `:330-333`**（taskkill＋Wait） | 留着到 `-test.timeout=120s` |
| 9 | `cmd/wisp/resident_approval_live_246_windows_test.go:446`｜`:457` | `esclistener.exe`／wisp 子进程 | `t.Cleanup` Kill＋Wait `:460-463` | 留着（⚠ `//go:build windows && winlive` ⇒ CI 只编译⛔ 跑，`ci.yml:623` 那把 winlive `go vet` 尺） |
| 10 | `cmd/wisp/resident_audio_247_live_windows_test.go:239`｜`:240` | `powershell.exe` 声源 | `t.Cleanup(stop)` `:251`（`sync.Once` 包 Kill＋Wait `:244-248`） | 留着；tag＝`windows`（⛔ winlive）⇒ **CI 形状要读 skip 逻辑，本腿⛔ 跑 ⇒ `Q-304-a1-E`** |
| 11 | `cmd/wisp/resident_sink_nail_127_windows_test.go:219`｜`:230` | `wisp.exe` 形子进程，**`CREATE_NEW_PROCESS_GROUP` `:222`** | `t.Cleanup(leg.stop)` `:238`（stop＝Kill `:287`；Wait 在协程 `:234`） | 留着；★新进程组＝控制台 `Ctrl+C`/`Ctrl+Break` **⛔ 端得到它**（R6 量到 3 枚 `SysProcAttr` 全这一形） |
| 12 | `cmd/wisp/resident_task_source_246_windows_test.go:457`｜`:463` | 同上，`CREATE_NEW_PROCESS_GROUP :460` | `t.Cleanup(leg.stop)` `:468` | 同 #11 |
| 13 | `cmd/wisp/resident_task_source_live_246_windows_test.go:493`｜`:507` | `wisp.exe -watch`（常驻） | `t.Cleanup` Kill＋Wait `:510-514` | 留着（`windows && winlive` ⇒ CI ⛔ 跑） |
| 14 | `cmd/wisp/secret_argv_windows_test.go:234`｜`:246` | `wisp.exe`（卡在 stdin 读，为读它自己的 PEB 命令行） | 体内内联：`stdin.Close()`＋Wait `:256`/`:263` | 父死⇒写端句柄关⇒子读到底（**读码推**，⛔ 实测）；若那推不成立＝一枚常驻 `wisp.exe` |
| 15 | `cmd/wisp/secret_argv_windows_test.go:335`｜`:344` | 同一枚二进制，另一形 stdin | 体内内联 Start `:344`＋Wait `:351` | 同 #14 |
| 16 | `cmd/wisp/secret_argv_windows_test.go:450`｜`:460` | 自 re-exec 探针子（`^TestProcessCommandLineProbeHelperProcess$`，decl `:420` 已核） | 体内内联 Wait `:466`/`:470`（⛔ Kill） | 留着到 `-test.timeout=120s` |
| 17 | `cmd/wisp/slo_exit_os_156_windows_test.go:322`｜`:326` | 自 re-exec 角色子（`^TestSLO156FixtureChildrenDoWhatTheirNamesSay$`，decl `:100` 已核） | **`t.Cleanup` 只 Wait、⛔ Kill** `:329`（注释逐字"never leave a child behind"）；helper 文件 `:313-315` 自述"starts a role child WITHOUT reaping it" | 孩子按角色自退：`exit` 形立刻退、`sleep` 形睡 `slo156SleepFor=30s`（`:73`）⇒ **最多 30s**；兜底 `-test.timeout=300s` |
| 18 | `internal/proc/jobscope_windows_test.go:206`｜`:207` | 不存在的 exe（正向控制，`TestStartInJobRejectsBadCommand` decl `:199` 已核） | — | ⛔ 子代（`Start` 必失败）；⚠ 它**不在 R3 的 17 里**＝`grep '.Start()'` 这把尺看不见它（`Start` 藏在 `jobscope_windows.go:114`） |

**票面 `:15` 那两枚 spawn 点的重量（现跑）**：尺＝`grep -rEn "print-addr" --include=*.go internal cmd tools` 命中 `adaptertest/mockllm.go:72` ＋ `openaichat/mockllm_integ_test.go:73` ＝**两枚，整族**（＋ `tools/mockllm/main.go:46` 的 flag 定义本身）⇒ 与票面"命中恰两枚"同值（票面那把尺的字面是 `'print-addr'`，本腿的是同一枚字样）。★票面 `:15` 明写**两枚必同批改**——本表把它们记为 **#1 与 #16（下表 §2 的 N2）**，任何只改一枚的落地形都留第二个泄漏源。

## 2. 非测试文件里的 6 枚（R2−R1 的作差，逐枚）

| # | 文件:行 | 是什么 | 回收 |
|---|---|---|---|
| N1 | `internal/llm/adaptertest/mockllm.go:65` | `go build -o <t.TempDir()>/mockllm.exe .`（建夹具） | 同步 `CombinedOutput` `:67`＝父等它 |
| **N2** | **`internal/llm/adaptertest/mockllm.go:72`**｜`Start :79` | **`mockllm.exe` 假大脑＝票面孤儿的另一枚源** | **`t.Cleanup` Kill＋Wait＋协程 join `:98-104`**；注册**早于**就绪等待（`:82-83` 注释逐字"Cleanup is registered before the readiness wait so every failure path below still reaps the child"），失败分支体内再 Kill `:126` | 
| N3 | `cmd/balldebug/diff_windows.go:409` | 产码（调试工具）起子进程，`CREATE_NEW_PROCESS_GROUP :412` | 产码侧，⛔ 本票射程（票面 `:18`／`:35` 已把产码那侧归 `Q-57`／票 156） |
| N4 | `cmd/wisp/doctor.go:316` | `cc --version` | 同步 `Output()` |
| N5 | `cmd/wisp/slo_windows.go:492` | 被测主体子进程 | **已经进 Job**：`rt.Job.StartInJob(cmd) :500` ⇒ 产码侧的甲形先例 |
| N6 | `tools/d22scan/gitignore.go:373` | 工具自己的 `git`（独立 module，`ci.yml:246` 那把 `go vet (tools/d22scan module)` 尺） | `exec.CommandContext`＝ctx 约束 |

★**R1（打在 `_test.go` 上的那把尺）整档看不见 N1/N2**——而 N2 恰是本票立票的那枚孤儿的头号来源。⇒ 只跑 `grep --include=*_test.go` 的普查会**报出一枚假的"泄漏源只有两枚里的一枚"**；本票要求两把并排（R1＋R2），这就是理由。

## 3. 同步 inline 那一档（**48 枚＝集合差算出来的，⛔ 手数**）

推导尺（可逐字重跑，件＝`logs/partition-66-18-48.txt`＋两份名册 `logs/partition-start18.txt`／`logs/partition-sync48.txt`）：

`cut -d: -f1,2 logs/grep-m1-testonly.txt | sort`（66） 对 `18 枚 Start 档名册`（§1 全表）作差 `comm -23` ⇒ **48**；两向作差皆空（⇒ §1 那 18 枚**一枚都⛔ 漂**、都真在 R1 的 66 里）。

按**实参形**分组（枚数＝6＋2＋8＋7＋7＋8＋1＋9＝**48**，分组靠逐枚读过的调用形，⛔ 一把正则推整族）：

| 形（跑的是什么） | 枚数 | 行名册（`文件:行`） |
|---|---|---|
| `icacls` 直写形（权限夹具） | **6** | `internal/agent/spill_acl_windows_test.go:79`、`internal/config/private_acl_windows_test.go:204`、`internal/models/no_seal_ruling_windows_test.go:85`、`internal/winsec/inherited_narrow_notice_104_windows_test.go:49`、`cmd/wisp/config_reload_perm_223_windows_test.go:52`、`cmd/wisp/logsink_windows_test.go:112` |
| 通用 helper `run`／`mustExec`（`exec.Command(name, args…)`，实参由调用方给＝`icacls` 与 `cmd /c mklink` 两形） | **2** | `internal/winsec/acl_windows_test.go:38`、`internal/winsec/reparse_windows_test.go:25` |
| `cmd /c mklink /J`（junction 夹具，直写形） | **8** | `internal/memory/artifacts_junction_tripwire_windows_test.go:39`、`internal/risk/pathresolver_junction_windows_test.go:50`、`internal/tools/bridge_junction_windows_test.go:52`、`internal/tools/paths_ticket107b_probes_test.go:31`、`internal/tools/paths_workspace_test.go:203`、`internal/tools/task_output_pointer_notice_test.go:406`、`internal/winsec/resolve_windows_test.go:190`、`internal/winsec/seam_guard_windows_test.go:104` |
| `taskkill /F`（**杀手本身也是子进程**，秒退） | **7** | `internal/tools/bridge_a18_kill_windows_test.go:121`、`:141`、`internal/tools/fs_edit_ac4b_kill_windows_test.go:138`、`:157`、`internal/tools/fs_staging_windows_test.go:331`、`:349`、`:398` |
| `git` | **7** | `cmd/wisp/panel_host_gate_test.go:310`、`:322`、`:340`、`:354`、`:394`、`tools/d22scan/scan_test.go:1386`、`:2146` |
| `go` 工具链（build／vet／list） | **8** | `internal/llm/openaichat/mockllm_integ_test.go:65`、`internal/models/assembly_reachability_121_test.go:265`（`go list -e -deps ./cmd/wisp`）、`internal/proc/crossvet_test.go:67`（`go vet`）、`cmd/wisp/panel_host_gate_test.go:407`、`cmd/wisp/resident_approval_live_246_windows_test.go:417`、`cmd/wisp/secret_argv_windows_test.go:174`、`tools/d22scan/scan_test.go:2256`、`tools/d22scan/selftest_test.go:181` |
| 自 re-exec 测试二进制（同步等它） | **1** | `internal/config/c26_seam_posix_125_test.go:174`（`//go:build !windows`＝ubuntu 那一支） |
| 跑自家／别的二进制并等它退出 | **9** | `cmd/wisp/dataroot_128_windows_test.go:128`（`CommandContext`＋`defer cancel`＋`Run`，120s ctx＝**守卫⛔ 测量**，注释逐字在 `:119-122`）、`cmd/wisp/early_log_nail_130_windows_test.go:121`、`cmd/wisp/leg_dispatch_gate_133_test.go:1580`（`-test.list`）、`cmd/wisp/resident_audio_247_live_windows_test.go:217`（`tasklist`）、`cmd/wisp/secret_argv_windows_test.go:369`、`cmd/wisp/slo_report_144_windows_test.go:881`（`cmd.exe /c exit 7`＝非零退出是夹具本身）、`internal/models/acl_sid_121_windows_test.go:65`（`powershell` 读 SID）、`tools/d22scan/scan_test.go:2351`、`tools/d22scan/selftest_test.go:190` |

⇒ **甲/乙/丙的射程只该打 §1 那 18 枚**：同步形那一档父进程等它退出，外力杀父时它们**要么已经退了、要么跟着管道/控制台一起没**（读码级）——但"跟着一起没"**⛔ 是本仓证过的**（票面 `:5` ③ 只说 Windows ⛔ 父死子随），⇒ 归 `Q-304-a1-A`（表②）。

## 4. 这几把尺看不见什么（派单 §4：命中 0 ⛔ 等于没问题）

1. **R7＝testdata 里 spawn 0 枚（rc=1）**——但 `cmd/wisp/testdata/esclistener` 确实是**被测试起起来的夹具**：spawn 点在 testdata **外面**（`resident_approval_live_246_windows_test.go:417` 建它、`:446` 跑它）。⇒ "testdata 里 0 命中"⛔ 等于"没有夹具子进程"；这把尺只能证"夹具⛔ 在自己目录里藏 spawn 码"。
2. **R1（`_test.go`）整档看不见 `adaptertest/*.go`**＝N2，本票头号残留源（§2 具名）。
3. **R3（`.Start()`）看不见 `StartInJob` 那一支**（`Start` 在 `jobscope_windows.go:114` 里）＝#18 那一枚；也⛔ 看得见"同一个 spawn 变量既可能被 `Start` 也可能被 `StartInJob`"这种**分支形**（#3）。
4. **R1/R2 都看不见⛔ 走 `os/exec` 的起法**：`os.StartProcess` 全 scope **0 枚**（R6）、`syscall` 直接 `CreateProcess` 的形本腿**没打尺**＝⛔ 排除，只登记（`Q-304-a1-I`）。
5. **词面尺⛔ 能证的**：`mockllm` 之外还有没有**别的假大脑/别的长驻夹具**（例：`docker compose` 起的那枚 mock-llm，`ci.yml:430` 的 `Start compose test services (mock-llm on 18080)`＝**编排器起的，⛔ 测试二进制起的**，本名册按票面"测试里起子进程"的射程把它记为⛔ 入册、但它是同一族残留的**另一种来源**）⇒ 见 `30-costs.md` 丙形一节。
6. **R8 的 17 枚全在编译射程之外**：`.scratch/**` 以点开头 ⇒ Go 工具链⛔ 进这些目录（`go/build` 的.walk 规则，**读码级**，本腿⛔ 跑 `go list`）；`scripts/spike/webview2-latency/main.go` 那 2 枚＝**手动 spike 程序**（⛔ 测试起的），落在编译射程内但⛔ 测试名册里。
7. **R9＝41 枚测试文件** ⇄ 票面**没有快照值**可比 ⇒ 只当"文件枚数"报，⛔ 当"用例枚数"用。
