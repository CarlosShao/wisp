# 152 r1 — 真 subject 中途死掉那一条出口，第一次被**量**到：它在生产接线上今天产不出自己那句红

- 交付程角色：**实现程（票 152）**，写码＋取证
- 派单：`.scratch/wisp/dispatches/2026-09-26-091x-r152-r153-r151.md` 派单 B
- 工单：`.scratch/wisp/issues/152-the-real-subject-death-shape-has-never-been-measured-only-faked-to-be-reachable-and-the-offset-fallback-is-a-leg-that-changes-no-reading.md`
- 来路：`docs/evidence/s1/149-three-unpinned-outlets-r1-accept-r1.md` 第 11 格第 1、2 条（同一枚洞第三次挂账）
- 地界（本程**只**动这四类路径）：`cmd/wisp/slo_windows.go` · `cmd/wisp/slo_report_144_windows_test.go` ·
  `.scratch/wisp/probes/152/**` · 本件

---

## 第 0 格　锚点·工作树·互斥·两枚仪器坑的绕法

### 0.1 开工现读

```
$ git rev-parse --short HEAD
5365cb2                                    <- 本程开工锚点（派单共同锚点是 86b0161）

$ git log --format='%h %ad %s' --date=format:'%H:%M' 86b0161..HEAD
5365cb2 09:10 docs(派单存档): 票 152／153／151 三枚写码程派出 …   <- 编排者自己那一枚，只含 .scratch/wisp/dispatches/
ff550f3 09:24 test(153 AC#1): …                                    <- 票 153 的程，internal/agent/**
cdf2471 09:24 feat(frontend 动画开关机器): …                       <- 前端会话
ad9b29f 09:30 feat(frontend 设置屏动画开关): …                     <- 前端会话
b23c7f7 09:30 feat(153 AC#2): …                                    <- 票 153 的程

$ git log --oneline 86b0161..HEAD -- cmd/wisp/slo_windows.go cmd/wisp/slo_report_144_windows_test.go
（空）
```

⇒ **锚点 ≠ 派单的 `86b0161`，按派单那句"不算漂移"处理**，点名是谁的：漂移全部来自
①编排者自己的派单存档（`5365cb2`，一枚路径 `.scratch/wisp/dispatches/2026-09-26-091x-…md`）、
②票 153 的程（`internal/agent/**`）、③前端会话（`frontend/**`＋`docs/reports/frontend-session-log-zcode.md`）。
**本票两枚被测文件在 `86b0161..HEAD` 区间里零枚 commit 碰过**（上面那条 pathspec 现量为空），
且工作树字节与 `86b0161` **逐枚 md5 相同**：

```
$ for f in cmd/wisp/slo_windows.go cmd/wisp/slo_report_144_windows_test.go; do
    a=$(git show 86b0161:$f|md5sum|cut -d' ' -f1); b=$(md5sum < $f|cut -d' ' -f1); echo "$f same=$([ "$a" = "$b" ] && echo YES)"
  done
cmd/wisp/slo_windows.go                  same=YES   md5=0403d5196f4bc0c20993dca7ceeea2b8
cmd/wisp/slo_report_144_windows_test.go   same=YES   md5=4d07365e12e9e156771020f3b638dcfa
```

⚠ 顺带一枚对上别家读数的交叉核验：这两枚 md5 与票 149 验收件 §2.1 在 `c3f7224` 上量到的
`0403d519…` / `4d07365e…` **逐字符相同** ⇒ 从 149 结线到本程开工，这条洞所在的两枚文件**一个字节都没动过**，
"第三次挂账"期间没有任何修法混进来。

### 0.2 工作树与互斥（文件级）

```
$ git status --porcelain | grep -E 'cmd/wisp|internal/'
 M cmd/wisp/run.go                          <- 票 151 的程在写（mtime 09:25:30）
?? cmd/wisp/task_scope_close_151_test.go    <- 同上
 M internal/tools/bridge.go                 <- 同上
```

⇒ 与派单 C（票 151）同目录 `cmd/wisp/`，**按文件级判互斥成立**：它动 `run.go` ＋新增
`task_scope_close_151_test.go`，本程动 `slo_windows.go` ＋ `slo_report_144_windows_test.go`，交集为空。
⇒ **本程的"零命中／干净"类宣称一律不含 `cmd/wisp/run.go`、`internal/tools/bridge.go`、`frontend/**`、`design/**`**：
它们要么被别人正写着，要么派单明令不计进任何宣称。
⇒ 由此得一条方法约束：**门禁四数与静态三把尺一律在 `git archive` 快照里跑**（§0.4），
不在那棵会被人同时改的活树上跑，否则本程报的数里会混进 151 的中间态。

### 0.3 工具版本（现读，未背数）

```
$ go version
go version go1.27.1 windows/amd64
$ D:/work/base/gopath/bin/gofumpt.exe --version
v0.12.0 (go1.27.1)
$ ls third_party/sherpa-onnx/*.dll
onnxruntime.dll  sherpa-onnx-c-api.dll  sherpa-onnx-cxx-api.dll        (3 枚)
```

### 0.4 两枚仪器坑的**绕法**与**绕没绕成的判据**

坑①（`exit status 0xc0000135` ＋ 0 条 `=== RUN`，因至少两枚：dll 不在位 / `PATH` 写成 `D:/…` 盘符正斜杠形）：
⇒ 本程所有 `go test` 的 `PATH` 一律用 **MSYS 形** `/d/work/…`，且**每一枚日志都现数 `=== RUN` 枚数**当"跑没跑到"的唯一凭据；
"dll 在位"只当必要条件、不当充分条件。本程实际取到的枚数：探针 #1 `RUN=8`、探针 #2 `RUN=2`、
门禁快照 `RUN=136`（见各日志），**没有一枚是 0**。

坑②（`-overlay` 与 `-cover*` 合用时 overlay 被静默忽略）：
⇒ 本程的**变异全部走 `-overlay`、且命令里不出现任何 `-cover*` 旗标**；
本程**没有跑覆盖剖面**（本票不需要"走到哪一行"这一类读数，需要的是"印出哪一句"，那是 stdout 上的文本读数）。
⇒ 反向自证：凡带 `-overlay` 的命令行原文都落盘在日志首行，复算的人一眼能看见它没有 `-cover`。

### 0.5 仓外临时件（只建不删）

`D:\work\tmp\wisp152\` 下：`scratch/probe/`（两枚探针源＋三枚 overlay.json）、
`scratch/e2e/`＋`e2e2/`＋`e2e3/`（三版端到端日志与 stderr 原文）、`snap-pre/`（`git archive 5365cb2` 的快照＋dll）、
`wisp.exe` / `wisp-anchor.exe` 两枚构建产物、`roster-pre.txt`。
**仓库目录内没有建过 worktree 或 checkout**；探针源另抄一份进 `.scratch/wisp/probes/152/` 供复算（149 同形）。

**放水两问自答**：① 本格零判据改动、零断言方向改动（只读＋建仓外目录＋一枚快照）；② 本格不使用任何 helper。

**第 0 格判定**：**成立**（锚点现读、漂移逐枚点名到人、两枚被测文件与工作树三方 md5 相符、互斥按文件级成立）。

**本程没测什么（本格）**：没核 `86b0161..HEAD` 之外更早区间里这俩文件的历史（那是 144/147/149 的账）；
没查那 3 枚 dll 与 CI runner 用的是不是同一批字节（派单没要求，149 验收件同样登记为未查）。

---

## 第 1 格　AC#1 主读数：**真起子进程**，量"subject 中途死掉"那条出口今天到底印不印

### 1.1 判据（本格自定，写死在跑之前）

工单 AC#1 要求的是"把'真 subject 中途死掉'**量成一次读数**，而不是再造一个假对象"。本格把"读数"定死成三枚可观测量：

1. 循环**印出的那一句原文**（`err.Error()`，不解析、不美化）；
2. 那句**点不点名停在第几字节**（句里有没有 `… at offset N`）；
3. **`s.exited()` 在被判定的那一瞬返回值** ＋ 同一瞬**操作系统对那枚 pid 的答复**
   （`OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION)` ＋ `GetExitCodeProcess`，
   与 `internal/tools/staging_live_windows.go:34-46` 同形，是**外部证人**不是 Go 自己数）。

要"真"，本格的下述四条一条不退：真 `exec.Command` 起真进程（真 pid）、真退出码 7、
真落在盘上的真文件（子进程自己 `os.ReadFile` 之外的写）、**真 `os.ReadFile` 读**（`readReportFile` 留 nil＝生产接线）。
被动的变量只有 case 13 当年**手工钉死的那两半**：循环开跑前有没有 `Wait()`、读的是不是脚本。

### 1.2 现量命令原文

```
$ cd "D:\work\workspace\projects plans\Wisp"
$ export PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:$PATH"     # MSYS 形，坑①
$ go test -count=1 -overlay "D:/work/tmp/wisp152/scratch/probe/overlay-probe.json" -v -run TestP152SubjectDeathProbe ./cmd/wisp/
# 日志原文：.scratch/wisp/probes/152/probe-exited-reachability.log
# RUN=8 TOPPASS=1 SUBPASS=8 FAIL=0 SKIP=0 P152lines=64
```

探针文件**没有落进仓库**：`overlay-probe.json` 把仓外源映射成 `cmd/wisp/zz152probe_windows_test.go`，
现量 `ls cmd/wisp/zz152probe_windows_test.go` → `No such file or directory`，`git status --porcelain cmd/wisp` 里
**没有本程的任何文件**。

### 1.3 读数（7 格矩阵，每格都是真子进程；红字是本程加的强调，原文在日志里）

| # | 格名 | 子进程 | 循环前 Wait？ | 读文件 | OS 答复 | `s.exited()` 循环前 | **印出的句子（原文）** | elapsed |
|---|---|---|---|---|---|---|---|---|
| 1 | `prod-prefix-unreaped` | 写 989B 前缀后 exit 7 | **否（＝生产）** | 真 `os.ReadFile` | terminated code=7 | **false**（`exitCode=-1`） | `wisp slo: subject 2572 never wrote a complete report within 4s (last read: 989 bytes read, document still open at offset 989: the tail had not arrived)` | **4.02s（＝预算烧尽）** |
| 2 | `prod-nofile-unreaped` | 没建文件就 exit 7 | **否（＝生产）** | 真 `os.ReadFile` | terminated code=7 | **false** | `wisp slo: subject 3056 never wrote a complete report within 4s (last read: 0 bytes read, no report file yet)` | 4.02s |
| 3 | `prod-prefix-unreaped-count` | 同 1 | 否 | **计数壳包的 `os.ReadFile`** | terminated code=7 | false | 同 1 那一句（timeout 腿） | 4.02s，**真读 186 次** |
| 4 | `prod-prefix-REAPED` | 同 1 | **是**（`cmd.Wait()`） | 真 `os.ReadFile` | — | **true**（code 7） | `wisp slo: subject 64584 exited (code 7) without writing its report (989 bytes read, document still open at offset 989: the tail had not arrived)` | **0s** |
| 5 | `prod-nofile-REAPED` | 同 2 | 是 | 真 `os.ReadFile` | — | true | `wisp slo: subject 25924 exited (code 7) without writing its report (0 bytes read, no report file yet)` | 0s |
| 6 | `case13-shape-scripted-REAPED` | 同 1 | 是 | 脚本（case 13 原样） | — | true | 与 4 同形（exited 腿），脚本被读 **1** 次 | 0s |
| 7 | `unreaped-scripted` | 同 1 | 否 | 脚本（case 13 原样） | terminated code=7 | false | 与 1 同形（timeout 腿），脚本被读 **189** 次 | 4s |

### 1.4 读数判读（四句，每句都只靠上表的列）

1. **`exited()` 那一条出口在生产接线上一次都没印过自己的句子。**
   三格"没 Wait"的格（1／2／3／7）印的全是 **timeout 腿**，四格印 exited 腿的（4／5／6）全都 Wait 过。
   ⇒ 翻转它的**唯一变量是那一行 `cmd.Wait()`**，与读的是真文件还是脚本无关（对比 3↔7、4↔6）。
2. **外部证人同瞬指认"这枚 pid 已经死了"**：`os_says=alive=false code=7 (GetExitCodeProcess: terminated)`，
   而同一瞬 `s.exited()=false`、`s.exitCode()=-1`。⇒ 这不是"死得不够快"（那一格是**等到 OS 说死**才开跑的），
   是 `exec.Cmd.ProcessState` 在没有 `Wait()` 之前**永远是 nil**。
3. **"到点红"确实印了，也确实点名的停在第几字节**——但是由 **timeout 腿**印的：
   `989 bytes read, document still open at offset 989`。⇒ 工单问的"印的句子点名不点名停在第几字节"这一半，
   答案是**点名的**；不成立的是**归因**：一具尸体被报成"写得慢"，而且**为此把预算整段烧光**（4.02s／186 次真读；生产那一档是 33s）。
4. **case 13 的读数本身不是假的**（第 6 格复算到同一句），它缺的是**前提**：
   它把"循环开跑前已经有人 `Wait()` 过"当成给定，而生产接线里没有这个人。
   ⇒ 工单第 1 条那句"只被 case 13 **造**到可达，没被**量**到"——**本程量到了，而且量到的是反向**：
   那条支在生产接线上**今天走不到**，不是"走到了但不说话"。

### 1.5 机制的字据（现量名册，不是推理）

`s.exited()` 读 `s.cmd.ProcessState`（`slo_windows.go:794`），而 `ProcessState` 只由 `(*exec.Cmd).Wait` 填。
本程把整枚文件里唯一那几个相关调用点名：

```
$ grep -n '\.Wait()\|\.Run()\|StartInJob\|ProcessState' cmd/wisp/slo_windows.go
498:	p, err := rt.Job.StartInJob(cmd)
794:	return s.cmd != nil && s.cmd.ProcessState != nil && s.cmd.ProcessState.Exited()
801:	return s.cmd.ProcessState.ExitCode()
816:		_ = s.cmd.Wait()                       <- 全文件唯一一处 Wait，在 stop() 里面

$ grep -n 'go func\|^	go ' cmd/wisp/slo_windows.go
（空）                                          <- 没有哪条 goroutine 会去收尸

$ grep -n 'collectReport\|defer .*stop()' cmd/wisp/slo_windows.go
371:	defer subject.stop()
380:	defer instrumented.stop()
392:	observer, err := instrumented.collectReport()     <- 收尸在 392 之后（defer），判死在 392 之内

$ sed -n '110,123p' internal/proc/jobscope_windows.go    # StartInJob 的真身
114:	if err := cmd.Start(); err != nil { … }
117:	if err := j.Assign(cmd.Process); err != nil {
118:		_ = cmd.Process.Kill()
119:		_ = cmd.Wait()                                  <- 只有指派失败那一支才 Wait
```

⇒ `collectReport()` 从 `startSubject()`（只 `Start`）拿到 `s`，一路到 `collectReportWithin` 的循环，
**没有任何一处调用过 `Wait`**；`stop()` 是 `defer`，排在 392 之后。§1.4 第 1 条那句"翻转它的唯一变量是 `Wait()`"
在这张名册上是同一件事的两种说法：读数是现象，名册是它为什么必然如此。

**放水两问自答**：① 本格**零生产码改动**（`cmd/wisp` 在跑完本格之后仍与 `86b0161` md5 相同），
也没有为让哪一格"通过"改过任何判据——工单预填的那三种可能里，本程量到的是它写的那句"如果量出来那条路今天**产不出**红句"
的**加强版**（产不出**它自己那句**，产得出**邻居那句**）；② 本格的 helper 全部是被测包自己的
（`slo144Report`／`scriptedReader`／`collectReportWithin`／`s.exited`），探针自己只新增 OS 证人一枚（`OpenProcess`＋`GetExitCodeProcess`），
**没有 mock 掉任何一件被测物**。

**第 1 格判定**：**成立**——但成立的是**取证**那一半：`exited()` 支在生产接线上不可达，
到点红由 timeout 腿出，**归因错（尸体报成慢 writer）＋代价（整段预算）**。
⇒ AC#2 按工单那句"不许预填结论"走 ⓑ/ⓒ，见第 3 格。

**本程没测什么（本格）**：
1. **没量"真 wisp subject 在 `os.WriteFile` 那一笔里断掉"**（真前缀由真 subject 写出的形状）：
   那要求把刀落在一次单 write 的中途，144/147 已量过"本机撞得出部分读、但唯一形状是 `len=0`"，
   本程**沿用**那条已知不可复现的结论，用"真子进程写下的真前缀"代替（§1.3 第 1 行）。
   ⇒ 谁先被骗：以为"`exited()`＋真前缀"这形在真 subject 上永远到不了的人——本程证的是**接线到不了**，不是**世界到不了**。
2. 没量 `exited()` 支在 `reportCorrupt` 之后的先后（第 778 行的 corrupt 早返回把 781 挡在前面）——
   那一支是 149 的地界，本票明令不重开。
3. 每格只跑一遍（`-count=1`），没做重复性标定；`4.02s`／`186 次` 这类数是**本机读数**，
   能引的是"烧没烧尽预算"这个量级，不是那枚 186。

---

## 第 2 格　同一条根因的第二枚出口：`waitReady()` 那句"exited early (code N) before reporting ready"

**判据**：§1.5 的名册对 `slo_windows.go:520` 那一支同样成立（`waitReady` 也在任何 `Wait` 之前问 `s.exited()`）。
这一格不接受"名册推出来所以不量"——本票的整个来历就是"三张表都写没测"，所以本程把它量了：
真子进程、**不写 ready 标记**、exit 7，然后跑**真的 `s.waitReady()`**（它等的是常量
`subjectReadyBudget = 60s`，本程一字节没动，所以这一发要烧满 60 秒）。

### 2.1 现量命令原文

```
$ cd /d/work/tmp/wisp152/snap-pre                      # git archive 5365cb2 ＋ dll 3 枚
$ export PATH="/d/work/tmp/wisp152/snap-pre/third_party/sherpa-onnx:$PATH"
$ go test -count=1 -overlay D:/work/tmp/wisp152/scratch/probe/overlay-probe2-samppre.json -v -run TestP152BWaitReadyProbe ./cmd/wisp/
# 日志原文：.scratch/wisp/probes/152/probe-waitready-unreaped.log
```

### 2.2 读数（五行全引）

```
P152B|child|os_alive=false os_code=7 how=GetExitCodeProcess: terminated pid=48704
P152B|predicate|exited_before_waitReady=false exitCode=-1 ready_file_exists=false
P152B|waitReady|elapsed=1m0.02s err=wisp slo: subject 48704 never reported ready within 1m0s
P152B|after|exited=false exitCode=-1
P152B|after_reap|exited=true exitCode=7
```

⇒ 与第 1 格同形：**OS 说它死了（code 7）、`s.exited()` 说没死、循环烧满整段预算、印的是预算句**
（`never reported ready within 1m0s`），而 `exited early (code 7) before reporting ready` 那句一次没印。
最后一行是本程在测量之后自己 `Wait()` 了一下的反向证明：同一枚 `s`、同一枚子进程，
`Wait()` 之后 `exited()=true / exitCode=7`——**变的只有收尸这一件事**。

⇒ **所以这不是一枚"出口写歪了"的洞，是一枚"根因在生命周期上"的洞**：`sloSubject` 的
`exited()` 语义（"有人收过尸"）与它的两处用法（"进程死没死"）之间缺一个收尸点，两处调用都缺同一个东西。

**放水两问自答**：① 本格零码改动、零预算改动（`subjectReadyBudget` 未被引用也未被改，那一发烧的正是它自己）；
② helper 只有包内真函数＋OS 证人，无 mock。

**第 2 格判定**：**成立**（第二枚出口同因、同形，且已量成读数）。

**本程没测什么（本格）**：没量 ready 标记**已写出之后**才死的那一格（那是第 1 格的接线，结论同一条）；
没量 `startSubject` 指派失败那一支（`jobscope_windows.go:119` 那一枚 `Wait` 是**唯一**在 `collectReport` 之前
可能出现的收尸点，而它只在 Assign 失败时走、且走完就直接返回错误、进不到循环）——本程**没有构造过那一发**。

---

## 第 3 格　AC#1 的端到端补量：**真 `wisp slo`**＋**真 subject 子进程被真 TerminateProcess**＝操作员看到的那一句原文

第 1 格量的是包内的真子进程（真 pid／真退出码／真文件／真 `os.ReadFile`），但"subject"那一枚是测试二进制扮的。
本格把最后半段距离也量掉：跑**真的 `wisp.exe slo`**（观察器），让**真的 `wisp.exe slo -subject -self-sample`** 子进程
在半路被 `Stop-Process`（＝`TerminateProcess`）打死，然后抄操作员 stderr 上的原文。

### 3.1 现量命令原文（三版脚本只建不删，逐版点名）

```
$ bash /d/work/tmp/wisp152/scratch/e2e-subject-kill.sh      # v1
$ bash /d/work/tmp/wisp152/scratch/e2e-subject-kill2.sh     # v2
$ bash /d/work/tmp/wisp152/scratch/e2e-subject-kill3.sh     # v3 ＝本格引用的那一版
```

- **v1 第一次跑出的是"空名册"**：`Get-CimInstance` 那行写在 bash 双引号里，`-Filter "Name='wisp.exe'"` 的引号
  到了 PowerShell 已经变形 ⇒ 枚举个数＝0，脚本却照样往下印 `victim pid=`（空）与 `observer rc=0`。
  日志原文：`.scratch/wisp/probes/152/e2e-v1-broken-instrument.log`（`NO VICTIM FOUND - … this run measures nothing`）。
- **v1 重跑（同一条命令）出的第二发更坏**：枚举这次成功了，但脚本用 `cut -d'|' -f1` 取 pid，取回来的是
  `62424\`（一枚带反斜杠的脏 token），`Stop-Process` 当场 `Cannot bind parameter 'Id'`、`stop rc=1`，
  而脚本继续把"kill 过了"的叙述写完，最后印 `observer rc=0`。
  日志原文：`.scratch/wisp/probes/152/e2e-v1-second-run-broken-pid.log`。
  ⇒ **这一发如果没人看 rc，会被读成"杀了子进程、观察器还是绿的"——一条彻底的假读数。** 本程据它写了 v2 的三道 guard。
- **v2 的读数是真的但接线不纯**：它用的 `wisp.exe` 建自 09:26 的**活工作树**，而那棵树在 09:25:30 已经带上
  票 151 未提交的 `cmd/wisp/run.go`（§0.2 点过名）⇒ 那枚 exe 里混着别人的中间态。v2 的句子与 v3 逐字符相同
  （`e2e-subject-kill-real-cli.log`），但**本件引的读数一律按 v3**，v2 降级为"同形的第二发、接线不纯"。
- **v3＝本格凭据**：先证字节再跑——

```
$ for f in cmd/wisp/slo_windows.go cmd/wisp/slo_report_144_windows_test.go cmd/wisp/run.go; do
    a=$(git show 5365cb2:$f|md5sum|cut -d' ' -f1); b=$(md5sum < /d/work/tmp/wisp152/snap-pre/$f|cut -d' ' -f1); echo "$f same=$(…)"
  done
cmd/wisp/slo_windows.go                 anchor=0403d5196f4bc0c20993dca7ceeea2b8 snap=… same=YES
cmd/wisp/slo_report_144_windows_test.go anchor=4d07365e12e9e156771020f3b638dcfa snap=… same=YES
cmd/wisp/run.go                         anchor=862af1373b868e7fa51f276363de015f snap=… same=YES   <- 151 的中间态不在里面
$ (cd /d/work/tmp/wisp152/snap-pre && go build -o /d/work/tmp/wisp152/scratch/wisp-anchor.exe ./cmd/wisp)
build rc=0 exe=/d/work/tmp/wisp152/scratch/wisp-anchor.exe size=28482471
```

### 3.2 对照：同一条命令、同一台机、不杀任何东西（阴性对照）

```
$ … wisp.exe slo -state Sleeping -seconds 3 -settle-ms 0 -interval-ms 250 -out base-report.json
baseline rc=0 elapsed=7s  baseline report bytes=17083
```
⇒ 这台机器、这套旗标能把整条 `wisp slo` 跑绿（17083 字节报告），所以 3.3 那发红**不是**环境红。

### 3.3 读数（v3 全引）

```
observer shell pid=40155 at 09:33:28
--- attempt 1 sees: ---
56208 …\wisp-anchor.exe slo -state Sleeping -seconds 6 -settle-ms 0 -interval-ms 250 -out …/kill-report.json
12736 …\wisp-anchor.exe slo -subject -state Sleeping -seconds 6 … -ready-file C:\Temp\wisp-slo-subject-2583625739\ready -out …\subject-report.json
28528 …\wisp-anchor.exe slo -subject -state Sleeping -seconds 6 … -ready-file C:\Temp\wisp-slo-subject-571289499\ready -out …\subject-report.json -self-sample
victim pid=28528
Stop-Process rc=0
victim GONE (no process object)
victim subject report file at kill time: ls: cannot access 'C:\Temp\wisp-slo-subject-571289499\subject-report.json': No such file or directory
observer rc=2 elapsed_since_kill=36s
--- observer stderr verbatim ---
wisp slo: in-tree record unavailable (fail-closed, no silent downgrade to the tree basis): wisp slo: subject 28528 never wrote a complete report within 33s (last read: 0 bytes read, no report file yet)
kill report bytes=NO FILE (failed closed)
```

四枚点名读数：
1. **受害者是按它自己的命令行认出来的**（命令行里带 `-self-sample` 的那一枚，pid 28528），
   不是按"第一个 wisp.exe"猜的——同瞬 OS 名册里有三枚 wisp-anchor.exe（观察器＋两枚 subject），
   认错就整发作废，所以那三行原文留在日志里。
2. **杀是真的**：`Stop-Process rc=0` ＋ `victim GONE (no process object)` ＋ 它的报告文件此刻还不存在。
3. **观察器 rc=2（fail-closed 没有变成绿）**，等场 **36s**（≈ 常量 `subjectReportBudget+subjectGrace = 33s` ＋窗口尾巴），
   然后印的是——
4. **操作员看到的那一句原文**：
   `wisp slo: in-tree record unavailable (fail-closed, no silent downgrade to the tree basis): wisp slo: subject 28528 never wrote a complete report within 33s (last read: 0 bytes read, no report file yet)`

⇒ **这一句里同时有对的半句和错的半句**：
对的——它点名了 pid、点名了停在第几字节（`0 bytes read`）、并且真的红（rc=2，不降级）；
错的——它把这具尸体说成"33 秒内没写完"，而 OS 在 36 秒前就已经报告它死了、`Stop-Process` 的退出码本来就在手上。
⇒ 工单第 1 条"谁先被骗：排『`wisp slo` 到点红该看哪一句』的人（D37：错误要映射到可行动的原因）"——
**本程量到的正是这一发**：可行动的原因（subject 死了，code 几）在这一句里被换成了不可行动的原因（等得不够久）。
⇒ 而本票要它说的那一句 `exited (code N) without writing its report` ——**在这条端到端路上一次都没出现过**；
它与第 1 格第 4／5／6 格的差异只有"有没有人在循环开跑前 `Wait()` 过"。

**放水两问自答**：① 本格零码改动（exe 建自快照字节，md5 三枚现量相符）；
② 没有 mock：观察器、subject、`Stop-Process`、stderr 全是真的；本程自造的两枚坏仪器（v1 两发）
**按"仪器坏了"点名而不是当读数**，与票 149 验收件 §N-1 第 7 行同一条规矩。

**第 3 格判定**：**成立**（AC#1 的"真 subject 中途死掉"已有端到端读数；红句印、字节点名、归因错、代价 33s，四味全在原文里）。

**本程没测什么（本格）**：
1. **端到端那一发是"死在写报告之前"**，不是"死在 `os.WriteFile` 那一笔中间"——真 subject 的唯一一笔写无法用外力
   稳定切进中途（§1.4 第 1 条同一条已知形状）。**带前缀的端到端读数**由第 1 格第 1 行给（真子进程写的真 989 字节前缀），
   那一发的 subject 是测试二进制扮的。⇒ 谁先被骗：把这两格读成"真 subject＋真前缀已经端到端量过"的人。
2. 没量"被杀的是**被测** subject（不带 `-self-sample`）"那一支——那一支归 `sampler`，不经过 `collectReportWithin`，
   本票地界外，本程**没有**顺手测。
3. 没量 `wisp slo -settle` 那条腿；没跑真 CI（无 push 权限）。

---

## 第 4 格　AC#2：`exited()` 那一支该印什么——按 149 的三档答，答案是 **ⓒ：说不好 ⇒ 停手上报**

**判据**（工单 AC#2 原文三档）：ⓐ 该点名（补断言）／ⓑ 设计上不该走到 `summary()`（改成不会静默的形状＋一发"仍走到即红"）／
ⓒ 说不好 ⇒ 停手上报。工单同时写了："不许预填结论……如果量出来'那条路今天**产不出**红句'，那正是要的答案，按 ⓑ/ⓒ 处理并停手上报"。

### 4.1 为什么 ⓐ 与 ⓑ 都不合本程量到的形

- **ⓐ 不成立**（不是"错"，是**已被 149 收过、且治不到本程量到的那一发**）：149 按 ⓐ 补的就是 case 13，
  而 case 13 钉的是"**那条支被走到之后**句子带不带得上最后一次读数"。本程第 1／2／3 格量到的是**走到之前**就断了：
  生产接线上没有任何人收尸 ⇒ `s.exited()` 恒假 ⇒ 那一支一次都没被执行过。
  再补一枚断言只能钉"手工接线造出来的可达"，那正是本票要收掉的东西。
- **ⓑ 不合形**：ⓑ 的措辞是"设计上不该走到 `summary()`"，其判据是"改成不会静默的形状＋一发**仍走到即红**"。
  本程量到的不是"走到了但静默"，是**整条支走不到**（第 1 格四发 unreaped 全印邻居那句、第 2 格 `waitReady` 同形、
  第 3 格端到端同形）。照 ⓑ 硬改，造出来的是**一枚永远不会响的格**——工单明令不许。
- **真正缺的那一味不在本票地界里**：缺的是 subject 生命周期上的**收尸点**。本程把候选列出来，一条都不自己动手：
  ①在 `collectReportWithin` 的循环里补一次 `Wait`/轮询；②把 `exited()` 改成问 OS（`GetExitCodeProcess`，
  本程探针用的就是这一形）；③给 subject 配一枚 reaper goroutine——第三形还要过 D22 裸 `go func` 的 owner/recover 与
  D38 并发与线程模型。三形都改到 `slo_windows.go` 之外的接线（`startSubject`／`stop()`／`internal/proc` 的 Job），
  而"到点红该由谁在几秒内定性成什么"是 D37 错误映射的形状问题——**方案没覆盖，按 D22 闸门③停手，不自行假设**。

⇒ **本程 AC#2 那一半一字节生产码都没改**（第 5 格改的是 AC#3 那条退路，与这一支无关），也没有为它开"要它响"的格。

### 4.2 承重两问——两句都答，尺是本程自己的（`probes/152/my152.py`，原始读数在 `mut-anchor/` 与 `mut-shipped/`）

| 变异 | 跑在哪版字节 | RUN | 顶层 PASS | 顶层 FAIL | 红名 |
|---|---|---|---|---|---|
| `e1-exited-drop-summary`（`last.summary()` 换成写死的 `"nothing read"`） | 锚点 5365cb2 | 20 | 12 | **1** | `TestSLO149ExitedGiveUpSentenceCarriesTheLastReading` |
| `e2-exited-drop-summary-case13off`（同一发＋case 13 摘掉） | 锚点 | **19** | 12 | **0** | —（逃逸） |
| `e3-case13-not-a-test-only`（只摘 case 13，码不动） | 锚点 | 19 | 12 | 0 | — |
| `h1-exited-arm-deleted`（**整条 `if s.exited()` 支删掉**） | 交付字节 | 21 | 13 | **1**（case 13 跑满 60.05s） | 同上 |
| `h2-exited-arm-deleted-case13off`（删支＋摘 case 13） | 交付字节 | **20** | 13 | **0** | —（逃逸） |

`h1` 的红句原文（`probes/152/mut-shipped/h1-exited-arm-deleted.log`）：

```
slo_report_144_windows_test.go:837: exited give-up "wisp slo: subject 63728 never wrote a complete report within 30s (last read: 8 bytes read, document still open at offset 8: the tail had not arrived)" does not name the exit code
slo_report_144_windows_test.go:840: an exited subject was reported as a spent budget, not as a dead child: "wisp slo: subject 63728 never wrote a complete report within 30s (last read: 8 bytes read, document still open at offset 8: the tail had not arrived)"
slo_report_144_windows_test.go:843: exited subject was read 18485 times, want exactly 1 (the loop sees ProcessState and gives up on the spot)
```

- **问一："摘掉任意一味，是否存在一发变异从此打不红？"——存在，而且比 149 报的那发更凶。**
  ① 摘掉 case 13 ⇒ `e1`（摘 `last.summary()` 那半）从此 0 红；
  ② 摘掉 case 13 ⇒ 连**把整条 exited 支删掉**都 0 红（`h2`：FAIL 1→0、RUN 21→20）；
  ⇒ 全仓对那条支的覆盖**只有 case 13 一枚持有者**，而它靠的是"测试替循环把尸收了"这一手生产里不存在的接线。
  把它摘掉等价于"那条支不存在"——**没有任何读数会知道**。这一味因此是**承重的**（不是装饰），
  但它承重的那条路今天不通：`h1` 印出来的句子与第 3 格端到端那句**同形**
  （`never wrote a complete report within …（last read: …）`）——
  **今天的生产输出，就等于"那条支被删掉"那一发的输出。**
- **问二："摘掉它有没有任何外部可见读数变过？"——有。**
  摘掉 case 13：`=== RUN` 20→19（交付字节 21→20）、顶层 PASS 枚数动、名册少一枚；
  摘掉 `last.summary()` 那半：红的**枚数**与**原文**都动（`e1`/`h1` 三行红句在案）。
  ⇒ 与本票第 5 格那条退路**不同类**：那条是"摘了零读数变化"（装饰腿），这一味是"摘了必响"（承重腿）——
  只是它响的前提（走到那条支）在生产接线上今天不发生。

**放水两问自答**：① 本程没有为了让哪一格"通过"改过判据：`e*`/`h*` 是**造变异去量它**，不是去放宽断言；
② 尺是本程自己写的（`my152.py`），overlay 落盘后回读证明替换进了编译字节，且全程不带 `-cover*`（坑②）。

**第 4 格判定**：**ⓑ/ⓒ 停手上报（取 ⓒ 那一支：三档都没覆盖"支不可达"这一形）**；
承重两问**均已答**，凭据五发变异读数。**本程交付的是这次取证与"该由谁在哪收尸"这句待定案，不是一个修法。**

**本程没测什么（本格）**：
1. **没量"补了收尸点之后"的行为**——本程没实现任何候选修法，所以"改了就会印 exited 句"只是第 1 格
   第 4／5 那两发 `REAPED` 读数的**外推**（那两发是真的，但外推是本程说的，不是量到的）。
2. **没量 `waitReady` 之外还有没有第三处**同样在 `Wait` 之前问 `exited()`：本程的名册只覆盖
   `cmd/wisp/slo_windows.go` 一枚文件（`exited()` 的调用者见 §1.5），
   `internal/proc` 里若有同形调用**未被本程点名**。
3. `h1` 那发的 18485 次真读与 60.05s 是**本机单次读数**，能引的是"烧尽预算＋反复重读"这个量级。

---

## 第 5 格　AC#3：那条退路选了"明说自己不知道"那一支——两问都有读数，"在未修码上响不响"＝**响**

**判据**（工单 AC#3）：删掉，或者换成"会明说自己不知道"的形状，两选一、写清选了哪个为什么；
加注释／说成"纵深防御"／"留着反正不害"三件不算收；不许为它开"要它响"的格，但判据要能答"这一发在未修码上响不响"。

### 5.1 现量：那条腿"摘了零读数变化"，本程自己复算（不复用 149 的数）

尺＝`probes/152/my152.py`（同文件多枚替换**串联**；overlay 落盘后**回读**断言"新文本在／旧文本不在"，
否则 FATAL 不出数；命令行按 shlex 引号落盘——裸 `|` 的 `-run` 粘回去会变成管道，那不是"命令原文"）。
被验树＝`git archive 5365cb2` 快照＋dll 三枚；选择＝`TestSLO144|TestSLO147|TestSLO149`；全程零 `-cover*`（坑②）。

| 发 | 改了什么 | RUN | 顶层 PASS | FAIL | SKIP | 与 asis 比 |
|---|---|---|---|---|---|---|
| `asis` | 无 overlay | 20 | 13 | 0 | 0 | — |
| `f1-fallback-to-zero` | 锚点那行 `return inputOffset` → `return 0` | 20 | 13 | 0 | 0 | **四数逐枚相同 ⇒ 零枚外部读数变化** |
| `f2-fallback-to-minus1` | 同一行 → `return -1` | 20 | 13 | 0 | 0 | **同样零变化** |

⇒ 149 第 5 格那枚"装饰腿"读数在本锚点**复算成立**。`f2` 还额外量到一件事：
**只把值换成 -1、句子不认**，外部读数仍然一动不动——所以本票要收的不是那枚数，是**那句承诺**
（`g3` 的红句印的正是 `-1` 冒充位置的样子）。

### 5.2 选了哪一支，为什么

**选"换成会明说自己不知道"，不选"删掉"**：字面删不掉——149 验收件 §5.1 量到"函数还有返回值，
删掉默认臂 Go 直接编不过"；本程按"别家结论当未验证断言"的规矩自己撞了一遍（把 `return offsetUnknown` 摘掉＝
`missing return`）。所以**能删的只有那句承诺，不是那行代码**。落地形状（`cmd/wisp/slo_windows.go`，commit `10e3585`）：

- `contradictionOffset(err, inputOffset)` → **`contradictionOffset(err) int64`**：两家错误都不叫位置那一臂交回
  `offsetUnknown`（新具名常量，值 -1），调用点**不再把 `dec.InputOffset()` 递进去**——
  那枚数正是本函数存在的全部理由要说谎的数（149：18 发 corrupt 里 8 发印 0）。
- 新增 `contradictSubjectReportErr(size, offset, err)`：`offset >= 0` 那臂的字符串与改前**逐字节相同**
  （case 11／2／5／12 都靠它，改一个字符就会红——这就是"命名臂没被动过"的负向凭据：交付字节上整条选择
  21 枚 RUN／14 枚顶层 PASS／0 红，见 `probes/152/fixed-tree-selection.log`）；
  `offset < 0` 那臂印 `at an offset the decoder did not name`。
- 新增 case 14 `TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne`：
  两枚"叫得出位置"的臂用的是**真解码器错误**（`readSubjectReport` 造、`errors.As` 从 `%w` 包装里取回，不是手搓夹具），
  只有"叫不出"那一臂 handed 一枚非 JSON 错误——那正是这条臂过去区分不了的东西。

### 5.3 承重两问＋"这一发在未修码上响不响"

| 发 | 改了什么 | RUN | PASS | FAIL | 红名／红句 |
|---|---|---|---|---|---|
| `asis`（交付字节） | 无 overlay | 21 | 14 | 0 | — |
| `g1-restored-borrowed-value` | 把借来的值放回去（`return offsetUnknown` → `return 0`）＝未修码的行为 | 21 | 13 | **1** | case 14 |
| `g1-restored-borrowed-value-case14off` | 同一发＋摘掉 case 14 | **20** | 13 | **0** | —（逃逸） |
| `g2-case14-not-a-test-only` | 只摘 case 14，码不动 | **20** | 13 | 0 | — |
| `g3-renderer-always-prints-position` | 句子那一臂不认 `offset < 0` | 21 | 13 | **1** | case 14 |
| `g4-renderer-never-prints-position` | 句子那一臂永远不印位置 | 21 | 12 | **2** | case 11 **＋** case 14 |

`g1` 的红句原文（`probes/152/mut-shipped/g1-restored-borrowed-value.log`）：

```
slo_report_144_windows_test.go:748: contradictionOffset for an error that names no position is 0, want -1 (offsetUnknown); any other number is a byte position this error never named - that is the value ticket 149 killed
```

`g3` 的红句原文（＝"只改值、不改句子"会印出来的样子）：

```
slo_report_144_windows_test.go:754: no-position corrupt sentence "39 bytes contradict a subject report at offset -1: invalid character: the runner wrote something that is not a json error type" is not the admission shape
```

⚠ **一把不够用的尺要点名**：把交付版用例**原文**叠到锚点的 `slo_windows.go` 上跑，是**编译期就断**的——
`probes/152/n1-newtest-on-anchor-code.log` 原文四行：`not enough arguments in call to contradictionOffset`、
`undefined: offsetUnknown`（两处）、`undefined: contradictSubjectReportErr`，`=== RUN` 枚数 **0**。
⇒ "在未修码上响不响"这一问，本程**不以**那发为凭据（编不过是构建状态，不是读数），
凭据落在 `g1`：它在交付字节上把未修码的行为**逐字放回去**，case 14 当场红，红句里带的就是那枚被借的数。

- **问一："摘掉任意一味，是否存在一发变异从此打不红？"——存在。** 摘掉 case 14 ⇒ `g1` 从 1 红变 **0 红**
  （`g1-restored-borrowed-value-case14off.log`），同一发里 `g3` 也没人接
  ⇒ **case 14 是"没有位置就不许印位置"这条承诺的唯一持有者。**
  ⚠ 诚实的一半：`g4`（命名臂被消音）由 **case 11 与 case 14 各抓**，那一臂**不是** case 14 独占的覆盖面。
- **问二："摘掉它有没有任何外部可见读数变过？"——有。** `RUN` 21→20、顶层 `PASS` 14→13
  （`g2-case14-not-a-test-only.log`），名册差集少 `TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne` 一枚
  ⇒ 与本票第 4 格那枚承重的 exited 腿同形、与它自己开局那条退路"摘了零读数变化"的装饰形状**相反**。

**放水两问自答**：① 本程**没有**为这条腿开"要它响"的格——它响的是"把借来的值放回去"这一发，
判据方向是"不许印自己没有的位置"，与 147／149 已钉的三档同向，没有任何断言被放宽、没有阈值被动；
② 新增的 helper（`contradictSubjectReportErr`）在本票地界内，只被生产调用点与本票 case 14 用，
**没有替换掉任何一件被测物**：`readSubjectReport` 仍是真解码器、真 `%w` 包装、真句子。

**第 5 格判定**：**成立**（AC#3 选了第二支并有可复算凭据；两问都答；"未修码上响不响"＝响，凭据 `g1`）。

**本程没测什么（本格）**：
1. **没有量到那条臂在真实文档上被走到**——本票全部读数的共同前提仍是"今天走不到"。
   本程**没有**证明它不可达（那要枚举 `encoding/json` 的错误全集，含未来版本）。
   ⇒ 谁先被骗：把"明说自己不知道"读成"这一发已经有人见证过"的人——见证的是接缝，不是世界。
2. 没做 `offsetUnknown` 与 `note` 两腿的语义合并验证（都写 -1，同族不同来源）：本程只在注释里写清，
   没造一发"把 -1 换成正数会怎样"的读数。
3. 没验 case 14 那三枚真错误在**别的 Go 版本**上仍归 Syntax／UnmarshalType 两家（与 149 第 4 格同一条洞）。
