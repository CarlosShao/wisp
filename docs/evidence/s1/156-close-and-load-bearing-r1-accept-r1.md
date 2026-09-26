# 票 156 — 非实现者验收表（accept-r1）：七格逐格裁 ＋ 承重按操作定义攻 ＋ 收口纪律

- 裁的角色：**非实现者验收程（accept-r1）**。本程**不是**任何一格的实现者：只裁、只记，**一个字的生产码／测试文件／注释都没动**。
- **被验版本＝`7ca1130`（钉死）**。派单＝`.scratch/wisp/dispatches/2026-09-26-165x-accept-156-seven-cells-pinned.md`（`ea59d28`）。
- 本程 step 0 的 **HEAD＝`ea59d28` ≠ 钉住的 `7ca1130`**（差一枚＝派单自身那枚 commit）⇒ **按派单 §"HEAD 会往前走"那一句，照常裁 `7ca1130`**，锚不动。
  取数一律 `git show 7ca1130:<路径>`（本程写面之外的东西一律不许从工作树读）。
- 票面七枚框**现量＝7 枚未勾、0 枚已勾**（`git show 7ca1130:…156-…md | grep -c '^- \[ \] '` ＝ **7**；`'^- \[x\] '` ＝ **0**），
  与派单 §1 声称的形状一致；本程**不勾、不改票面**。
  ⚠ 已入库的票面在 `7ca1130` 之后**长了 6 行**（编排者 17:0x 的进度条目，`git show HEAD:票面`＝94 行）——
  本程逐字读过那 6 行：**里面没有替本程预写的判语**，它写的是"七框仍由它定，我不勾"＋把 §4 两枚以"可以顶回我"的形交给本程。
  本程引任何票面文字一律用 `7ca1130` 那版（88 行），不用这版。

## 凭据档位（本表每一枚数都带档，无档的数不许出现）

| 档 | 意思 |
|---|---|
| **〔自跑〕** | accept-r1 本程当轮自己跑出的读数（命令原文在表里，落地件在 `.scratch/wisp/probes/156-accept/**`） |
| **〔日志＋归档，本程抽验〕** | 别的程入库的日志，**本程自己重新数过四数／名册**，不是转述 |
| **〔某程读数，本程未复算〕** | 只引用，没重算——**本表凡是这一档都点名是谁的读数，且不用它下判** |

本程自己的尺全部落在 `.scratch/wisp/probes/156-accept/`：
`acc156.py`（**import r1 那把 `probes/156/my156.py`**，本程一枚变异字面都没重打）·
`ac6-census.sh` · `ac6/golden-fix.sh` · `gate/recompute-landed.sh` · `gate/gate156-accept.sh`。
日志在 `logs/`、名册在 `logs/roster/`、门禁在 `gate/`。

---

## 第 0 格　三条前提独立重走（派单 §2，不许继承）＋本程的取数面

### 0.1 P1 码真的在 `7ca1130` 里，且 `exited()` 走的是 (b) 支问 OS

```
$ git show 7ca1130:cmd/wisp/slo_windows.go            # 1031 行，本程落在 probes/156-accept/anchor/slo_windows.go
$ grep -n 'exited()\|cmd.Wait()\|ProcessState\|queryProcess\|exitStatus\|WithHandle\|WaitForSingleObject\|GetExitCodeProcess'
854:func (s *sloSubject) exited() bool          -> 855: gone, _ := s.exitStatus()
867:func (s *sloSubject) exitStatus() (gone bool, code int)
878:	if ps := s.cmd.ProcessState; ps != nil {          # (a) 我们的记录，且只在它存在时才问
881:	// (b) the OS, for the window this type actually lives in.
882:	return s.queryProcess()                            # ← 本票选定的那一味
904:func (s *sloSubject) queryProcess() (gone bool, code int)
909:	err := s.cmd.Process.WithHandle(func(h uintptr) {
910:		if waited, _ := windows.WaitForSingleObject(windows.Handle(h), 0); waited != windows.WAIT_OBJECT_0 {
913:		if cerr := windows.GetExitCodeProcess(windows.Handle(h), &osCode); cerr != nil {
936:		_ = s.cmd.Wait()                                  # 全文件唯一一枚 cmd.Wait()
$ grep -c 's\.exited()'                               3   （:522 waitReady ／ :822 collectReportWithin ／ :932 stop 的守卫）
```

**判〔自跑〕：P1 成立**。行号**本程自己量**：(b) 支确实在 **`:881-882`**（与 r2 证据件 §5.0 写的同号，
但**这不是继承**——本程上面那发 grep 是自己跑的；r2 同号纯属两版之间 `slo_windows.go` 没被动过）。
`cmd.Wait()` 全文件仍**唯一一枚**（`:936`，即票面 `:15` 说的"只在 `stop()` 收尾"那一条路径，位移＝上方插入行）。

⚠ **一处本程量到、与票面文字不同而 r1 已点名的形状**：票面 `:17` 写"`exited()` 的调用者全仓只有三枚：`:520`／`:837`／`:851`"，
其中 `:837` 是"`exitCode()` 的守卫"。**在 `7ca1130` 上这一句对不上**：`exitCode()` 现在问的是 `exitStatus()`（`:860`），
**不再是 `exited()` 的调用者**；`exited()` 的调用者是 `:522`／`:822`／`:932` 三枚，其中 `:822`（`collectReportWithin`）
正是票面漏记、r1 §0.2 前提③ 判"不成立"并点名要编排者认的那一枚。
⇒ 本程判：**r1 对前提③的处置正确且更重要的一层是——落地后那句"三枚"重新变真了，但指的是另一组三枚（`:522`/`:822`/`:932`）**。
`exitStatus()` 的注释 `:865` 写的"the three exited() call sites"在本锚点上**与码一致**（本程 `grep -c` ＝3），不是过期文字。
这一处不改判任何格，只登记：**引用"哪三枚调用者"必须带版本**，票面那组和 HEAD 那组同名不同址。〔自跑〕

### 0.2 P2 预算常量一字未动

```
$ git diff --numstat 9835d81 7ca1130 -- cmd/wisp/slo_windows.go        85	4 cmd/wisp/slo_windows.go
$ git diff 9835d81 7ca1130 -- cmd/wisp/slo_windows.go | grep -cE '^[-+].*(subjectReportBudget|subjectGrace|subjectReadyBudget|subjectPollInterval)'
0        （rc=1＝零命中；本程按坑单把 rc 打出来，没放进 && 链）
$ git show 9835d81:… | sed -n '75,88p'  vs  7ca1130 版 sed -n '77,90p'      diff 空 ⇒ const 区逐字节相同
77→79  subjectGrace = 3 * time.Second      83→85  subjectReadyBudget  = 60 * time.Second
84→86  subjectReportBudget = 30 * time.Second   85→87  subjectPollInterval = 20 * time.Millisecond
```

**判〔自跑：P2 成立〕**——四枚常量**值一字未动**，只有行号随上方插入漂移（77→79 等）。
`internal/observe/thresholds.go` 在 `7ca1130` 存量＝1 枚，且本票 10 枚实现 commit 上命中 0（见第 6 格）。

### 0.3 P3 没新增 go.mod 依赖

```
$ git rev-parse 9835d81:go.mod 7ca1130:go.mod      6ccb3fd1…  6ccb3fd1…   ← 同一枚 blob＝逐字节未动
$ git diff --numstat 9835d81 7ca1130 -- go.mod go.sum          （空，rc=0）
$ git log --format='%h %ad %s' --date=… -1 -- go.mod           5eb0f6b  09-21 09:46  ← 本票 era 之前 5 天
$ 逐枚（10 枚实现 commit，git diff-tree）                      gomodgomsum=[] 全部
$ 正控（同一把尺打在真动过 go.mod 的 commit）                  5eb0f6b -> go.mod hits=1   ← 尺会响
$ git show 7ca1130:go.mod | sed -n '11p'                       golang.org/x/sys v0.48.0
```

**判〔自跑：P3 成立〕**，且**这条是本程把自己的第一发仪器错换来的**：
本程最初用 `git show --numstat -1 --format='' <c> -- go.mod` 数命中，**十枚全报 `1	1	go.mod`**——
看着像"每枚 commit 都动了 go.mod"。反查后定性：

⚠ **新坑（本程现量，给坑单补档）**：**`git show --numstat -1 --format='' <c> -- <路径>` 不是逐枚尺**。
空 `--format=` 会让它对**不同的 `<c>` 吐出同一批数**，且**为本票没碰过的路径伪造命中**：
`7ca1130`／`3fbb1ce` 两枚都报 `slo_windows.go=85/4`（＝整票 range 数），`cc58445` 报 `56/17`；
`go.mod=1/1`、`go.sum=2/0`、`docs/PLAN.md=4/4` 对**三枚不同 commit 一字不变**。
⇒ 危害方向＝**假命中**（会把零命中判成越界）；正控（`5eb0f6b`→go.mod＝1）在**这把坏尺上照样"响"**，所以正控救不了它。
本程改用 `git diff-tree -r --no-commit-id --name-only <c> -- <路径>`，并用**真负控**（`7ca1130` 只动证据件 → 全空）
＋**真正控**（`5eb0f6b`→go.mod；`1218192`→internal/risk；`191f0d6`→docs/specs；`5866c6f`→docs/PLAN.md）双向钉住。〔自跑〕

### 0.4 本程取数面与互斥（工作树一个字没读成被验版本）

- 本程 step 0 现量脏工作树＝**42 条**（`git status --porcelain | wc -l`）：`design/**` 25 枚、
  `docs/evidence/s1/152-…-accept-r1.md`、`docs/reports/pending-and-issues.md`、`probes/152/my152.py`、
  `?? .zcodeignore`、`?? probes/156/__pycache__/`、`?? probes/156/mut-156-r2/asis.log` 等 ⇒ **一律不碰、不还原、不提交、不算进任何零命中宣称**。
- 但"跑测试"必须落在**字节等于锚点的工作树**上，本程没有含糊：
  `git diff --name-only 7ca1130 -- cmd/wisp internal/tools go.mod go.sum` ＝ **空**〔自跑〕
  ⇒ 工作树那两包与 `go.mod`/`go.sum` 与 `7ca1130` **逐字节相同**（跑在被验版本上，不是跑在别人的脏东西上）。
  本程门禁的 `meta.txt` 把这行钉在文件里。
- 争用（坑单"开测前查一次"）：`tasklist` 现量 **Runner.Listener.exe＝1、Runner.Worker.exe＝0、在跑的 go test＝0**〔自跑〕，
  与 r3/r4 的 `meta.txt`（listener=1／worker=0）**同形** ⇒ 本程的 ms 数与他们可比，但**不是"无争用"**。
- 仪器版本现读：`go version go1.27.1 windows/amd64`；`third_party/sherpa-onnx/*.dll` 3 枚在位；`gofumpt.exe` 在 `$(go env GOPATH)/bin`。〔自跑〕

### 0.5 伪授权反查（本程进场就做一次，第 7 格给两栏计数）

形状①"伪造我正在读的那枚文件里的段落"——本程对**四枚引用最多的文件**做 `7ca1130`／`HEAD`／工作树三向字节比：

```
156-exited-asks-os-r1.md            674/674/674 行  md5 553fbcd6 三向相同  diff(anchor,worktree)=0 行
156-r4-gates-and-cleanup-r1.md      449/449/449 行  md5 7085cf2d 三向相同  diff=0 行
AGENTS.md                            165/165/165 行  md5 008fc391 三向相同  diff=0 行
票面 156-…md                          88/94/94 行  ← anchor 88，HEAD 94：差的那 6 行＝编排者 17:0x 追加的进度条目
152-…-accept-r1.md（别程半件）        709/709/709 行  锚点 md5＝HEAD md5 ee39d32a，工作树 16f45130：diff＝6 行（3 加 3 删），
                                                  git status 仍 ` M`、7ca1130..HEAD 碰它的 commit 0 枚 ⇒ **至今未提交**
执行过的两把尺（额外补的一发，见 8.6 第 10 条）：
probes/156/my156.py                 锚点 md5 2168c128c157 ＝ 工作树 2168c128c157
probes/156/zz156probe_windows_test.go  锚点 d8987891201b ＝ 工作树 d8987891201b
`git diff --name-only 7ca1130 -- .scratch/wisp/probes/156/` 输出**空** ⇒ 已跟踪件相对锚点**零漂移**；
`git status --porcelain -- .scratch/wisp/probes/156/` 只列四枚 `??`
（`__pycache__/`、`mut-156-r2/asis.log`、`zero156-r4-head.sh`、`zero156-r4-work/`）——未跟踪件不进 `git diff`，两者是两把尺
```

⇒ 两枚实现件**没有被插段落**（工作树＝锚点＝HEAD，逐字节）；票面那 6 行**是已入库的真追加**，
本程读了全文，**不含替本程预写的判语**（原文写"七框仍由它定，我不勾"）。
形状②"伪造 commit／整节台账、替我要裁的那格预写判语"——本程对本表引用的**每一枚号**打了 `git cat-file -t`：
`7ca1130 ea59d28 9835d81 cc58445 0e95353 2ce77e1 2262c2b 3fbb1ce 1137781 1d53526 a4ec1f4 d661e8c 7ad2d98 10e3585 0f18652 5eb0f6b 1218192 191f0d6 5866c6f cbbbdf8`
**全部＝`commit`，零枚 fatal**〔自跑〕。派单 §"伪授权已发展到第 20 代"那两形在本程的取数面里各扫了一遍，读数在第 7 格。

**第 0 格判定**：**成立**——三条前提本程独立重走并各自给出会响的正控；
一处与票面不同的调用者口径已点名；两枚实现件三向字节反查无插入；十六枚被引号全部 `cat-file` 可解析。
**本程没测什么（本格）**：没验 `third_party/sherpa-onnx` 那 3 枚 dll 与 CI runner 用的是不是同一批字节（与 r1 §0.4 同一笔欠账）。

---

## 第 1 格　AC#1 —— "改前那一发"复现成**本程自己的读数**

**票面判据（`:33-35`）**：真起子进程、让它在写报告写到一半时死掉，量"OS 已知道／码还不知道"那个窗口有多大
（**具名次数与 ms**）；不许引票 152 的数当凭据；复现不出来就停手上报。

### 1.1 本程的尺与"改前"这一发是靠什么成立的

本程**没有**去 checkout 改前那棵树（派单禁读脏工作树、也禁在仓内建 worktree），
用的是**单点回退**：在 `7ca1130` 的字节上，用 r1 已入库的变异 `OS_ARM_REMOVED`
把 `:881-882` 那两行 (b) 支换成 `return false, exitCodeUnknown`，其余一字不改（`-overlay`，从不与 `-cover` 同用）。

**为什么这一发可以当"改前那一发"读**（本程自己量的等价性，不是推断）：

```
$ git show 9835d81:cmd/wisp/slo_windows.go | sed -n '832,841p'
832 func (s *sloSubject) exited() bool {
833	return s.cmd != nil && s.cmd.ProcessState != nil && s.cmd.ProcessState.Exited()
834 }
836 func (s *sloSubject) exitCode() int {
837	if !s.exited() { return -1 }
840	return s.cmd.ProcessState.ExitCode()
```
改前那枚谓词＝**只看 `ProcessState`**；`exitStatus()＋m1` 的可达状态里：
`s.cmd==nil`／`ProcessState==nil` 两支都答 `(false,-1)`，`ProcessState!=nil` 那支走 (a) 返回 `ps.Exited(), ps.ExitCode()`
⇒ **与 `:833` 逐状态等价**；唯一本程查过的边界是"前置守卫 `s.cmd.Process==nil` 会不会挡掉 (a)"——
`exec.Cmd.Process` 在 `Start()` 之后不再被置 nil（`Wait()` 只填 `ProcessState`），所以不存在
`Process==nil 且 ProcessState!=nil` 的可达状态〔自跑：读的是 `7ca1130` 的 `:867-882` 与改前 `:832-841`，非跑出来的〕。

### 1.2 本程读数（`logs/accept-probe-on-m1.log`，选择器 `^TestP156SubjectDeathProbe$`，`RUN=5 PASS=1 FAIL=0 SKIP=0 LOADED=YES`）

| cell | OS 何时知道 | 那一刻我们的记录 | 循环烧掉什么（**具名次数**） | 印出来的是哪句 |
|---|---|---|---|---|
| `prefix-unreaped` | `os_knew_at_ms=38`（`polls=20`、code=7 terminated） | `ProcessState_nil=true exited=false exit_code=-1` | `elapsed_ms=2009` **reads=187**、`agree=os_dead_to_code_dead_ms=-1 asks=130` | **BUDGET**：`subject 54356 never wrote a complete report within 2s (last read: 989 bytes read, document still open at offset 989…)` |
| `nofile-unreaped` | `os_knew_at_ms=35`（`polls=19`、code=7） | 同上 | `elapsed_ms=2008` **reads=190**、`agree=-1 asks=114` | **BUDGET**：`… (last read: 0 bytes read, no report file yet)` |
| `live-child`（阴性对照） | `alive=true code=259 polls=1` | `exited=false` | `elapsed_ms=2007` **reads=192** | **BUDGET**（正确：没死的孩子就该说预算） |
| `prefix-REAPED-case13shape` | `os_knew_at_ms=165` | `ProcessState_nil=false exited=true code=7` | `elapsed_ms=0` **reads=1** | **DEAD-CHILD**（(a) 支没被 m1 摘到） |

**窗口（AC#1 要的那枚数）**：`window=os_knew_at_ms=38_code_still_said_alive_after_loop_ms=2211`（`prefix`）／
`35 / 2210`（`nofile`）——**内核 38ms 就知道，那发里我们到死都没知道**（`agree=-1`＝循环全程 130 次问 `exited()` 零次为真）。

### 1.3 同一条尺在**发货字节**上的对照（`logs/accept-probe-on-ship.log`，同选择器，`RUN=5 PASS=1 FAIL=0 SKIP=0`）

`prefix-unreaped`：`os_knew_at_ms=48` → `agree=os_dead_to_code_dead_ms=7 asks=6`、`loop=elapsed_ms=0 **reads=1** verdict_kind=DEAD-CHILD`、
`sentence=subject 23104 **exited (code 7) without writing its report** (989 bytes read, document still open at offset 989…)`；
`nofile-unreaped`：`36` / `agree=6 asks=5` / `reads=1` / DEAD-CHILD；
`live-child`：`elapsed_ms=2005 reads=188` **仍是 BUDGET**；`prefix-REAPED-case13shape`：`reads=1` DEAD-CHILD **一字未变**。

### 1.4 与 r1 入库那发对照（本程自己重数，不是转述）

〔日志＋归档，本程抽验〕`probes/156/probe-pre/probe-156.log`：`44/2209`、`reads=186`、`asks=145`；
`probe-post/probe-156.log`：`50/6`、`reads=1`。与本程 `38/2211`、`187`、`130` 同形、同方向，
差值属本机差（r1 自己在 §5.3 也拒做跨机宣称）。**票 152 的 `18485 次／60.05s` 本程一枚都没当凭据。**

**放水两问自答**：① 断言方向——本程零判据改动（本格只跑探针与读 git 对象，探针不断言生产行为）；
② helper——本程用的尺是 **r1 已入库那把**（`my156.py` 的 `PROBE_ADD`/`OS_ARM_REMOVED`/`run`），
本程的 `acc156.py` 只加 cell 命名与打印，**一枚变异字面都没重打**，也就没有"造一发打不红的变异"这扇门。

**第 1 格判定**：**〔成立〕**〔全部关键数为本程 1.2／1.3 自跑；1.4 那两枚为抽验〕。
票面三要件逐条对上：**真子进程**（`exec.Command(os.Args[0], …)` 重入，pid／退出码 7 在原文里）、
**写到一半死掉**（子进程真写 `doc[:989]` 共 989／1978 字节后 `exit 7`；`nofile` 那发连文件都不存在）、
**窗口有 ms 与具名次数**（38ms→2211ms、187 次真读、130 次问）——**且本程是自己复现出来的，不是复述 r1 的。**

**本程没测什么（本格）**：
1. **没量"子进程死在 `os.WriteFile` 系统调用内部"** 那一形（本探针的半截是孩子主动写一半再退，属真前缀、不属被打断）——
   与 r1 §1「没测什么」第 2 条同一笔欠账，本程没能力在这一票里还。
2. **没量真 `wisp slo` 端到端**（真 Job Object 里打死真 subject 后操作员看到的原文）。
3. **本格的窗口数来自探针自挑的 `2s/10ms` 预算**，不是生产的 `30s/20ms`；生产常量那一形在第 5 格的 `m1` 发里量到
   （本程现跑 **`dead subject was read 1470 times, want exactly 1`／30.05s**，见 5.3）。本格不做外推宣称。

---

## 第 2 格　AC#2 —— 落地乙：`exited()` 从此问 OS ＋ 三件硬约束 ＋ 两向举证

**票面判据（`:36-41`）**：`exited()`（或等价的口径）改为**问 OS**；①不许动预算常量 ②不许新增 go.mod 依赖
③不许用墙钟时间差实现超时、不许裸 `go func(`（要 goroutine 得带 owner/recover）、选丙那一形要先回来批；
并且**举证"这一发在改前不响、改后会响"**。

### 2.1 落地的形状（本程自己读 `7ca1130`，不引 r2 的行号）

`exited()`（`:854-857`）与 `exitCode()`（`:859-862`）现在都是 `exitStatus()`（`:867-883`）的薄壳；
`exitStatus()` 两支：**(a)** `:878-880` 我们的记录（只在 `ProcessState != nil` 时才问，注释 `:871-877` 写明理由——
`Wait()` 之后 os/exec 不再出借句柄、pid 还可能被发给别人，所以收过尸之后"我们的记录"不是偷懒而是唯一安全的一支），
**(b)** `:881-882` → `queryProcess()`（`:904-922`）：`Process.WithHandle` → `WaitForSingleObject(h,0)` 判**等态**、
`GetExitCodeProcess` 读**码**，两问拆开；任何未决（`WAIT_FAILED`、`GetExitCodeProcess` 失败、句柄借不到）一律答"没死"。

本程核过这不是"顺手 `GetExitCodeProcess` 一发完事"那一形（派单 `:20` 警告的正是它）：
`queryProcess` **先等态后码**，所以一个以 259 退出的 subject 能被如实报 259 而不是被读成活着——
这一条本程**没有**量到发生过（见 2.5 第 2 条），但**反向的一发本程量到了**：探针 `live-child` 那发
`witness=alive=true code=259`，发货字节上 `exited()=false`、烧满预算、印 BUDGET〔自跑，1.3〕
⇒ **等态门没有把活孩子认成死的**。

### 2.2 三件硬约束——本程自己的尺，逐枚带正控

| 约束 | 本程的尺 | 读数 | 判 |
|---|---|---|---|
| ① 预算常量一字不许动 | `git diff 9835d81 7ca1130 -- cmd/wisp/slo_windows.go \| grep -cE '^[-+].*(subjectReportBudget\|subjectGrace\|subjectReadyBudget\|subjectPollInterval)'` | **0**（rc=1＝零命中）；const 区 `sed` 双向 `diff` **空** | **成立**〔自跑，见 0.2〕 |
| ② 不许新增 go.mod 依赖 | `git rev-parse 9835d81:go.mod 7ca1130:go.mod` ＋ `git diff-tree` 逐枚 10 枚 | 同一枚 blob `6ccb3fd1…`；10 枚全 `[]`；正控 `5eb0f6b`→`go.mod` 会响 | **成立**〔自跑，见 0.3〕 |
| ③-a 不许墙钟超时（ban #5） | `grep -acE 'time\.Now\|time\.Since' <(git diff … \| grep '^+')`＝**0**；`sed -n '904,922p' \| grep -ac 'time\.'`＝**0**；`observe.NewTimeout` 锚点 **2** 枚＝改前 **2** 枚（一枚没少） | **成立**〔自跑〕 |
| ③-b 不许裸 `go func(`（ban #1）／不许 reaper（丙未批） | `grep -ac 'go func'` 锚点 `slo_windows.go`＝**0**＝改前 **0**；新测试文件 `grep -cn 'go func'`＝**0**（rc=1）；`grep -acE 'reaper\|watchdog'`＝**0** | **成立**〔自跑〕 |
| ④（票面 `:65`）`stop()` 的 `Kill→Wait→RemoveAll` 不许动 | `git diff … \| grep -aE '^[-+]' \| grep -avE '^(\+\+\+\|---)' \| grep -aiE 'stop\|Kill\|RemoveAll\|Wait\(\)'` 命中 **3** 行，其中**非注释行＝0**（再 pipe `grep -avE '^[+-][[:space:]]*//'` 零命中，rc=1）；`cmd.Wait()` 全文件仍**唯一一枚**（`:936`） | **成立**〔自跑，见 0.1〕 |

`AGENTS.md §1.2` 那条"`filepath.Clean\|Abs` 只能在 `risk.PathResolver` 之外不出现"本程也顺手扫了新测试文件：
`grep -nE 'filepath\.(Clean\|Abs)'` ＝ **0 命中**〔自跑〕。

### 2.3 举证"改前不响、改后会响"——两向都是本程自己跑的

| 同一发、同一 cell、只换那一味 | **摘掉 (b)（=m1，行为等价于改前谓词，等价性论证见 1.1）** | **发货字节（`7ca1130`）** |
|---|---|---|
| `prefix-unreaped` | `exited()` 130 次全 false、`reads=187`、`elapsed=2009ms`、印 **BUDGET** | `exited()` 第 6 次亮（`agree=7ms`）、**`reads=1`**、`elapsed=0ms`、印 **DEAD-CHILD `exited (code 7) without writing its report`** |
| `nofile-unreaped` | 114 次全 false、`reads=190`、BUDGET | `reads=1`、DEAD-CHILD |
| `live-child`（阴性对照） | BUDGET、`reads=192` | BUDGET、`reads=188` —— **一字未变** |
| `prefix-REAPED-case13shape` | DEAD-CHILD、`reads=1` | DEAD-CHILD、`reads=1` —— **(a) 支没被摘到** |

〔全部 自跑，`logs/accept-probe-on-m1.log` 与 `logs/accept-probe-on-ship.log`〕
**用发货用例而不是探针**的第二向在第 5 格：`m1` 之下 case 17 量到 **`read 1470 times, want exactly 1`／30.05s**（生产常量），
同一 cell 零变异那一发 6/6 全绿〔自跑，`logs/accept-teeth-m1-jia.log`／`accept-teeth-asis-jia.log`〕。

反向也核了：**本票的改动面里没有一处以"让它响"为目的地动过既有断言**——
`slo_report_144_windows_test.go` 在整票上 `67 加／1 删`，而**非注释的改动行＝0**
（`git diff … \| grep -aE '^[+-]' \| grep -avE '^(\+\+\+\|---)' \| grep -avE '^[+-]//'` 零命中，rc=1）〔自跑〕；
整票生产码面**只有 `cmd/wisp/slo_windows.go` 一枚文件**，`internal/` 零枚〔自跑〕。

### 2.4 放水两问自答（本格）

① **断言方向动没动**——没动：本程零判据、零阈值、零 golden 改动（第 0 格那把 `git diff-tree` 尺在 10 枚 commit 上
对 `thresholds.go`／golden／`allowlist.txt`／`slo-check.ps1`／`tools/d22scan` 全 0 命中，见第 6 格）；
② **helper 是不是原有的那枚**——本程这一格用的两把尺全是**已入库的**：`probes/156/my156.py`（变异字面与 `run`）＋
`probes/156/zz156probe_windows_test.go`（探针），本程只加命名与打印；**没有新造变异、没有新造判据**。

**第 2 格判定**：**〔成立〕**。乙落地且经本程自己两向举证；四件（票面三件硬约束＋`stop()` 语义）各有本程的零命中读数；
**owner 批的是乙，本程量到的就是乙**（没有换丙：全文件 `go func`＝0、无 reaper）。

### 2.5 本格**没测什么**（三条缺口是本程自己核过"确实没钉"，不是转述）

1. **`sloSubject.stop()` 今天零枚用例**——本程自己数：`grep -rc '\.stop()' cmd/wisp/*_test.go` 非零的只有
   `resident_sink_nail_127_windows_test.go`（6 枚，是 `leg.stop()`，另一类型）与 `leg_dispatch_gate_133_test.go`（1 枚），
   `slo_exit_os_156_windows_test.go`／`slo_report_144_windows_test.go` **各 0 枚** ⇒
   r1 §2.2 末段那条**口径后果**（`stop()` 里 `!s.exited()` 现在会跳过对已死未收尸 subject 的 `Kill()`）**没有用例钉**。
   **这一条不改判 AC#2**（票面 `:65` 只要求"不许动收尾语义"，本程证实未动），但它是一笔**真实的未钉后果**，登记在第 8 格。
2. **subject 以 259 退出**那一形零枚用例（新测试文件里 `259` 只出现在 `slo156StillActive` 这枚**见证侧**常量，
   `:65-68`／`:374`）⇒ "先等态后码所以能如实报 259"仍是**接口性质**的论证，本程没把它升成读数。
3. **`WAIT_FAILED` 兜底支没有用例走到**（`:899-903` 注释写明兜底＝改前行为）⇒ 本程同意 r1 "硬造＝放水"的判断，
   **不为此开一格永不响的判据**；本程的读数侧证据是 `live-child` 那发两版逐字同形。

---

---

## 第 3 格　AC#3 —— 会响的检：新用例钉住"OS 说死 ⇒ 我们立刻算死"，且没把既有断言弄松

**票面判据（`:42-44`）**：新增用例必须钉住那一形，**并且**不能把票 152 已钉住的既有断言弄松（放水两问自答）；
不许为"让它响"而 `t.Skip`、不许改 golden／阈值。

### 3.1 五枚新用例在本程自己那发门禁里全绿（`gate/accept-post/cmdwisp.log`，144／84／0／0 rc=0）〔自跑〕

```
$ grep -aE '^--- PASS: TestSLO156' accept-post/cmdwisp.log | wc -l      5
$ grep -aE '^--- (PASS|FAIL): TestSLO14[4-9]|…TestSLO152' 同一枚日志    13 枚全 PASS
   含票 149 的 case 13 `TestSLO149ExitedGiveUpSentenceCarriesTheLastReading (0.04s)` 与
   票 152 的 case 14 `TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne (0.00s)`
$ grep -ac '^--- SKIP' 两包各                                        0
```

⇒ 在本程的锚点上，**新增五枚与票 144/147/149/152 那十三枚同时为绿**——"没弄松既有断言"这一句
本程不是靠名册数的，是靠**逐枚名字**在自家日志里读出来的。

### 3.2 钉的是不是那一形——本程逐条读断言（不是读注释）

`cmd/wisp/slo_exit_os_156_windows_test.go`（397 行、5 枚 `func Test`）里与本判据直接相关的四行：

| 行 | 断言（原文缩写） | 为什么它就是"钉这一形"而不是"钉这个答案" |
|---|---|---|
| `:159` | OS 报 terminated code 7 而 `ProcessState` 仍 nil ⇒ `exited()` 必须 true | 这正是 Q-57 拍的那一形 |
| **`:165`** | 问完 OS **不许**把 `ProcessState` 填出来（"asking the OS must not turn into reaping the child"） | **这一枚挡住"顺手在 `exited()` 里 `Wait()` 一下"那形**——那一形满足 `:159`/`:162` 却在此红。票面 AC#3 要的就是这一层 |
| `:170` | 第二次问仍 true | 挡住"只认一次"的实现 |
| `:191`／`:202`／`:212` | 活着 ⇒ false；我们 `Kill` 且未收尸 ⇒ true（(b)）；`Wait()` 之后 ⇒ true（(a)） | 双向界＋两支的分界都在 |
| `:257` | **`dead subject was read %d times, want exactly 1`** | "立刻"两字在这里落地：生产接缝上烧预算就红 |
| `:285`／`:288` | `waitReady` 那一支也要指名尸体与码，不许印"never reported ready within 1m0s" | 第二个出口也有主人（票 152 第 2 格量过这句"今天根本产不出"） |

### 3.3 本程自己的**逐味回退名册**——谁红谁不红，八发全自跑（`logs/accept-mut-*.log`，选择器 `TestSLO149Exited|TestSLO156`）

| 摘掉哪一处（本程自跑） | 四数 | 红的名字 |
|---|---|---|
| 零变异 asis | 6／6／**0**／0 | — |
| **m1＝只摘 (b) 支**（`:881-882`） | 6／2／**4**／0 | 15、16、17、18 |
| m2＝等态门永不判死（`:910` 改 `!= windows.WAIT_OBJECT_0`→`!= 1`） | 6／2／**4**／0 | 15、16、17、18 |
| m3＝除了 `WAIT_FAILED` 全判死（把活孩子认成死的） | 6／5／**1**／0 | **只有 16** |
| m4＝摘 (a) 支（`ProcessState` 那三行） | 6／4／**2**／0 | 16 ＋ **票 149 的 case 13** |
| m5＝等态留着、码不问了（`osCode = 1`） | 6／3／**3**／0 | 15、17、18 |
| m1＋改名掉 case 15 | 5／2／**3**／0 | 16、17、18 |
| m1＋改名掉 case 17 | 5／2／**3**／0 | 15、16、18 |
| m4＋改名掉 case 13 | 5／4／**1**／0 | 16 |

**名册两向**（防"一枚 panic 吞掉同包其余读数"）：以上每一发对 asis 做 `comm` 双向，
`asis-only`／`mut-only` **全为空**；三枚 `caseNN-off` 各只少**它自己改名那一枚**
（`TestSLO156ExitedAsks…`／`TestSLO156ReportLoop…`／`TestSLO149ExitedGiveUp…`）〔自跑，`logs/roster/`〕。
⚠ `logs/roster/*.log.runs.txt` 那批 0 行件＝**本程自己把 `$f.log` 又拼了一次 `.log` 造出来的作废名册**（第 8 格仪器错第 4 条），
**不是任何一发的读数**；权威名册是不带 `.log` 的那批（`<cell>.runs.txt`／`<cell>.verdicts.txt`）。按"临时件只建不删"留在盘上并同批入库，免得下一位按名字猜。

⇒ **本格的读数结论**：①"会响"不是宣称——(b) 支（m1/m2）、(a) 支（m4）、码那读（m5）、双向界（m3）**四种回退各有主人**；
②**case 16 是全包唯一一枚"把活孩子认成死的"的主人**（m3 只红它）；
③ case 13（票 149 的既有用例）**仍是 (a) 支的主人之一**，本程的 m4 直接量到它红——
这一条把 r1 §4 那句"case 13 走 (a) 那条支"从说明升成读数。〔自跑〕

### 3.4 helper 与"没弄松"的两把尺（本程自己重打，纠正一次自己的假零）

本程第一次问"判定用的 helper 是不是原有的那枚"时打的是 `grep -c 'func scriptedReader'` ⇒ **假 0**：
`scriptedReader` 在 144 文件里是 `type`（`:263`）不是 `func`。按名字形状重打并**对改前那版**比：

```
$ 9835d81 版 slo_report_144_windows_test.go：scriptedReader 2 处、scriptMissing 1、scriptBytes 1、slo144Report 1  ← 全部票 156 之前就存在
$ 7ca1130 版 slo_exit_os_156_windows_test.go 里重新定义这四枚？                       0／0／0／0（没有影子 helper）
$ 本票新写的 slo156* 只有五枚：Spawn／Witness／OsDead／OsAlive／PlayRole，
  它们唯一的 t.Fatalf 全在"夹具没照名字死／活"（:327、:363、:367、:372、:375、:388），无一条评产品
```

⇒ 判据用的 helper **全是原有的那四枚**，新增五枚只是取版器〔自跑〕。
"没弄松既有断言"的另一向：整票对该文件的 diff `67 加／1 删`而**非注释改动行＝0**（2.3 末那把尺），
新用例**没有**把任何一枚 `TestSLO14x/15x` 改名、跳过或删掉（3.1 的逐枚名字 ＋ 3.3 的名册两向）。

### 3.5 本格的两处**如实登记**（不是缺陷申报，是本程拒绝把没量到的写成量到的）

1. **色觉级"每枚用例唯一持有"只在 case 16 上成立**（m3 只红它）。15/17/18 的互不替代本程量到的是
   **断言行级**（`:165`／`:257`／`:285` 各仅一枚）——本程**没有**为它们各造一发"只红这一枚"的变异，
   因此**不采纳** r1 第 3 格表头那句"每枚在第 4 格都有它**唯一持有**的那一发"的全称形式。〔自跑〕
2. **case 19 不是"会响的检"**：它在 m1/m2/m4/m5 下**从不红**（本程上面八发里 case 19 全绿），这是**设计如此**
   （它钉的是夹具角色名实相符）。但 r1 第 3 格判定句把 19 与 15/17/18 并列为"四枚会响的检"——
   本程按自己的读数把它归到**夹具守卫**，与 r1 自己那张表里 19 那一行（"夹具自己"）一致，与判定句不一致。
   ⇒ 记为**证据件内部一处口径不齐**（第 8 格），**不改判 AC#3**：判据要的"这一形"由 15/16/17/18 四枚钉住，本程量到了。

**放水两问自答（本格）**：① 断言方向——本程与四程都没动任何既有判据（`t.Skip` 整票 0 命中、golden／阈值见第 6 格 0 命中、
`SKIP=0` 在本程自家门禁两包都是 0）；② helper——判定用的四枚 helper **全部早于本票**（3.4 现量），
本程复用的尺是 r1 入库的 `my156.py`，**没有新造变异字面、没有新造判据**。

**第 3 格判定**：**〔成立·带条件〕**。
**〔这一档是本程写完第 8 格之后回改的：原判"成立"，与 8.2 那张档位表不一致，本程取严的那一枚＝带条件；
C3（本格 `:283`／`:285` 那两行把 `144/84`＋"5 枚"写成 r1 自己那发 `gate-post`，实为 `143/83`＋4 枚）与
C5（同处判定句把 case 19 列为"会响的检"，本程八发变异里它从不红）都落在**本格自己的证据段**里，
所以本程不把"条件"只记进第 8 格而让这里留着"成立"——两处必须同档。〕
四枚新用例（15/16/17/18）钉住"OS 说死⇒我们立刻算死"，
其中 `:165` 单独挡住"改成去收尸"那一形、`:257` 单独挡住"烧完预算才认"那一形；
既有十三枚（含 case 13／case 14）在本程自家门禁里逐枚名字为绿，名册两向无任何一枚被吞或被改名。
**两枚条件＝C3＋C5**，全部为记录级、最小闭合集合＝**只追加更正、不动 r1 正文一字**（见 8.2）。

**本程没测什么（本格）**：`exitStatus()` 在 `s == nil`／`s.cmd == nil` 两支的行为（本程没造用例，r1 §3 同笔）；
`WAIT_FAILED` 兜底支（见 2.5 第 3 条）；`-race` 一枚没跑 ⇒ 新增每轮询 1–3 枚内核调用与 `stop()` 的竞态**无判据**。

---

## 第 4 格　AC#4 —— 票 152 §8.3 那三枚**注释级**动作（①②③逐枚裁 ＋ "分开 commit"这一条）

**票面判据（`:45-46`）**：把"现量第 4 条"那三枚动作做掉；它们**不影响 AC#2 成立与否**，
**分开 commit、分开报数**；②那枚改动前**必须**按自己的锚点现量那三行注释里的数。

### 4.1 动作①（补回 case 13 被删的裁定头，并写明"ⓐ 是 149 的裁定、152 的 AC#2 走 ⓒ"）——**成立**

本程不引 r2 的字节核，自己重打这把尺：

```
$ git show 10e3585 -- cmd/wisp/slo_report_144_windows_test.go | grep -a '^-//' | sed 's/^-//'   → 1 行
$ md5sum <那一枚删除行>                                          a756e044366a5c957bef470027a86bc0
$ md5sum <7ca1130 版 :835>                                       a756e044366a5c957bef470027a86bc0
$ diff <(那一枚) <(:835)                                        空 ⇒ 逐字节补回，票面 ① 前半成立
$ awk 'NR>=849&&NR<=857' 7ca1130 版                             "// Which ⓐ this is, and which one it is not (ticket 156 AC#4(i)…"
                                                                "The letter belongs to TICKET 149's AC#3 … ⓐ/ⓑ/ⓒ …
                                                                 It is NOT the letter ticket 152's AC#2 drew: 152 took ⓒ"
```

⇒ 票面 ① 要的两半（**补回**＋**按 ⓒ 写明那两枚字母各是谁的裁定**）**都在发货面上**，
且后半写在 `:849-857` 那段注释里、带真名路径（`docs/evidence/s1/152-…-accept-r1.md §2.2`，r2 §4.1 补的"指路"）。〔自跑〕

### 4.2 动作②（`:737-740` 那三行里的读数）——**成立**，且"改动前按自己的锚点现量"这一条本程复现了

本程的三发配对（尺＝`acc156.py ac4`，它 **import r1 的 `my156.py`** 的 `spec`/`case_off`，一枚字面没重打；
选择器＝票 152 原选择器 `TestSLO144|TestSLO147|TestSLO149|TestSLO152`，**不含**本票新增五枚）：

| 发 | 本程四数 | 红名 |
|---|---|---|
| `accept-ac4-g1-case14off`（借回值放回去 **＋** case 14 改名） | **RUN=20 PASS=13 FAIL=0 SKIP=0** rc=0 | — |
| `accept-ac4-g1-case14live`（同一发编辑、case 14 活着） | RUN=21 PASS=13 **FAIL=1** SKIP=0 | `TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne` |
| `accept-ac4-asis-152-surface`（不放任何变异） | RUN=21 **PASS=14** FAIL=0 SKIP=0 | — |

⇒ **票面 ② 给的 `20/13/0` 对**，且本程同时量出**错因**（不是"抄错一位"）：`21/14/0` 正是这选择器的**零变异基线**，
被贴进了一句承诺**变异读数**的话——三发同尺同形，这一句是本程的读数而不是转述。〔自跑〕

发货面核到字节：

```
$ awk 'NR==739' 7ca1130 版   "// g1-restored-borrowed-value-case14off.log reports === RUN=20, PASS=13, FAIL=0);"
$ grep -a -c 'reports === RUN=21, PASS=14, FAIL=0' 7ca1130 版     0   （rc=1＝零命中）
$ grep -a -c 'RUN=21, PASS=14, FAIL=0' 7ca1130 版                 1   ← 只剩 :749 那段"讲历史"的更正块里
$ 本程自己数那枚被 :739 点名的 tracked 日志：git ls-tree … mut-post/g1-…case14off.log = 1 枚（tracked）
  四数 RUN=20 PASS=13 FAIL=0 SKIP=0，首两行的 -run 图案与本程同形
```

⇒ `:739` 那句**今天两个断言都为真**：那一形的读数对，它点名的日志真的也印 20/13/0。〔自跑〕

### 4.3 动作③（`mut-shipped/` 那批入库 **或** 把票 152 实现件 §4.2 档位降级，两选一、写清选了哪个）——**成立，但闭的人是编排者**

```
$ git ls-tree -r --name-only 7ca1130 .scratch/wisp/probes/152/mut-shipped/ | wc -l    9
$ git ls-tree -r --name-only 7ca1130 .scratch/wisp/probes/152/mut-post/    | wc -l    6
$ git diff-tree -r --no-commit-id --name-only 7ad2d98 | wc -l                          9   ← 选的是"入库"那一支
```

票面 progress log（`7ca1130` 版 `:74-75`）写明选了哪一支、以及为什么不必再降级 ⇒ **两选一＋写清"选了哪个"两半都齐**。
**但这一支不是四个实现程做的**：r2 §4.3 判两支都越界、只交读数；**编排者 16:1x 自己闭的**。
本程如实记：**AC#4③ 成立，归属＝编排者**，不算任何一程的交件。〔自跑〕

### 4.4 票面"分开 commit"这一条——**②③做到，①没做到**（本程逐枚数文件清单）

| commit | 程 | 碰的文件 | 与 AC#2 的码同枚？ |
|---|---|---|---|
| `cc58445` | r1 | 11 枚，全 `probes/156/**` ＋ 证据件 | 否（那枚只有凭据） |
| **`0e95353`** | r1† | 25 枚：`slo_windows.go`＋新测试文件＋**`slo_report_144_windows_test.go`** | **是——AC#4① 与 AC#2 的码在同一枚 commit 里** |
| `2ce77e1` | r2 | 7 枚（AC#4 注释面＋r2 日志＋证据件） | 否 |
| `1d53526` | r4 | **1 枚**（`slo_report_144…_test.go`，7 加／7 删） | 否 |
| `7ad2d98` | 编排者 | 9 枚 `probes/152/mut-shipped/**` | 否 |

†＝编排者代提（r1 停在轮次上限）。⇒ 票面 `:46` 那句"**分开 commit**"对 **①** 未兑现：
它被装在 AC#2 那枚 checkpoint 里一起落地。本程量到这一枚的**成因**而不是猜：
`cc58445`（r1 交完第 0-1 格）之后 r1 就没再单独提过码——第 2/3 格的码与 AC#4① 的注释**同时**出现在 `0e95353`。
**这不是 AC#2 的缺陷**（本格票面自己写了"它们不影响 AC#2 的成立与否"），是 **AC#4 自己的报数纪律缺一半**。〔自跑〕

### 4.5 派单 §4 那笔账（①"历史化失效"的三处指涉）——本程的裁：**记录级，不是 AC#4② 的扣分项**

**编排者的裁定**：不改别人已入库正文、只许追加更正；把定性留给本程。
**本程不顶回来，并交出自己的读数说明为什么**：

```
$ 票面 ② 的射程（`:46`＋现量第 4 条 `:23`）写的坐标是 "同文件 :737-740 那三行"
$ awk 'NR>=737&&NR<=740' 7ca1130 版      → 这四行今天读到的就是被改真的那句（:739＝20/13/0）
$ 那三处过期指涉在 :742 / :749 / :743＋:752 —— 全在 :740 之外，属 r2 的更正块
$ grep -a -c 'RUN=21, PASS=14, FAIL=0'   = 1 枚，位置＝:749，且它的上下 10 行逐字在解释"这是什么的数"
```

⇒ 判据是**射程**：AC#4② 这条框只裁 `:737-740`，那里今天为真；`:742-752` 的问题是**时态**（"is neither reading"、
"at this anchor" 指的是 r2 的锚点），不是**数值**——今天按 `:739` 拿走错数的人不存在（本程 `grep -c` ＝0 命中那句句式）。
所以：**不扣 AC#4②**，也**不另开一格永不响的判据**；这三处的正确落点＝**另立一张记录级收口票**（像票 157 那样），
用**追加**的形把三处写成过去时，**不改正文一字**。最小闭合集合本程写在第 8 格。〔自跑〕

**放水两问自答（本格）**：① 断言方向——本程与四程都没动任何判据：整票 `slo_report_144_windows_test.go`
非注释改动行＝**0**（见 2.3 末尾那把尺），本程另证 `t.Skip` 在整票 diff 里 **0 命中**〔自跑〕；
② helper——本格不使用 helper，三发的尺是 **r1 已入库那把**（`my156.py` 的 `spec`/`case_off`/`run`）。

**第 4 格判定**：**〔成立·带条件〕**。
①成立（逐字节补回＋字母指路）／②成立（本程自己现量 `20/13/0` 并把那句改真，`reports === RUN=21…` 今天 0 命中）／
③成立但闭的人是编排者。**两枚条件**：
**C1** 票面"分开 commit"在 ① 上未兑现（与 AC#2 的码同枚 `0e95353`）——**最小闭合集合：无需动码，只需在台账记这一笔形，
或在下一张注释级票里把 ① 的注释面单独提一次**；
**C2** `:742`–`:752` 三处时态过期——**最小闭合集合：一张记录级收口票，用追加更正（不改正文）闭，见 4.5 与第 8 格**。



## 第 5 格　AC#5 —— 承重那一格**按操作定义攻**（不是"读一遍说成立"）

**定义（派单 §3 逐字）**：**摘掉任意一味，是否存在一发变异从此打不红？** 反面：**别要求"同一发变异里新加那支先响"**
（它挂在 `else if` 上、前提正是被抹掉的东西）。再问一句：**摘掉它有没有任何外部可见读数变过？**
**票面 `:48` 另加一条**：`摘一味导致编不过`**不算有牙**——要重做一发改得过、行为不变的对照。

本味＝`cmd/wisp/slo_windows.go` @ `7ca1130` 的 **`:881-882`**（`// (b) the OS…` ＋ `return s.queryProcess()`），
位置与内容本程自己在第 0 格钉过（**不是**照 r2 的行号去找）。摘法＝`my156.py` 的 `OS_ARM_REMOVED`
把两行换成 `// MUTATION m1…` ＋ `return false, exitCodeUnknown`：**(a) 一字未动、签名一字未动、编得过**
——本程下面每一发都真的编译并跑出了 `=== RUN`，这一条不是声称而是读数〔自跑〕。

### 5.1 问①——**三形本程全部自己重跑**（派单点名"至少两形，尤其丙"），答案是**分层的**

| 形 | 本程的发（`logs/`） | 选择器 | 本程四数 | 判 |
|---|---|---|---|---|
| **甲 单摘这一味** | `accept-teeth-m1-jia` | `TestSLO149Exited\|TestSLO156` | RUN=6 PASS=2 **FAIL=4** SKIP=0 | **有牙**：红的正是 15/16/17/18；仍绿的两枚＝case 19（夹具）＋**票 149 的 case 13**（它走 (a)，m1 拿它无可奈何——这正是 m1 只摘了 (b) 的证据） |
| **乙 摘这一味＋本票四枚持有者逐枚改名** | `accept-teeth-m1-plus-newcases-off-yi` | 全 SLO 五族 | RUN=22 PASS=15 **FAIL=0** SKIP=0 | **存在打不红的那一发**，形状＝"摘一味 **并且** 把持有它的四枚用例一起摘掉" |
| **丙 摘这一味，只问本票之前的旧名册** | `accept-teeth-m1-old-surface-bing` | `TestSLO144\|147\|149\|152` | RUN=21 PASS=14 **FAIL=0** SKIP=0 | 与**同选择器零变异基线** `accept-teeth-asis-old-surface-bing-asis`＝21/14/0 **四个数一字不差**，且**名册两向 `comm` 双向皆空**（21↔21，见 3.3 末与 `logs/roster/`）⇒ 改前那份测试面对这一味**零持有**。**这才是"打不红"的严格形**，也正是 AC#3 必须新增用例的理由 |

⇒ **本格按定义的答案：问①"是"——存在打不红的那一发**（乙、丙两形），
但**不是**"摘了没人能觉"那一形：两发都要靠**摘掉用例本身**才成立，而单摘一味（甲）红四枚。
本程**没有**把乙/丙读成"这一味没牙"，也**没有**只报甲。〔全部自跑；与 r2 §5.1 的 6/2/4、22/15/0、21/14/0 **同数**，
那是〔日志＋归档，本程抽验〕的对照，不是本程的凭据〕

### 5.2 问②——**变了，两个口径都变**，且变的正是操作员看到的那一句

**操作员口径**（本票的靶；本程自己两发探针，同一 cell 只换那一味）：

| cell | 零变异（`accept-probe-on-ship`） | 摘掉这一味（`accept-probe-on-m1`） |
|---|---|---|
| `prefix-unreaped` | `sentence=… **exited (code 7) without writing its report** (989 bytes read…)`、`loop=elapsed_ms=0 **reads=1** DEAD-CHILD`、`window=48/7`、`agree=7ms asks=6` | `sentence=… **never wrote a complete report within 2s** (last read: 989 bytes read…)`、`loop=elapsed_ms=2009 **reads=187** BUDGET`、`window=38/**2211**`、`agree=**-1** asks=130` |
| `nofile-unreaped` | DEAD-CHILD、`reads=1`、`window=36/6` | BUDGET、`reads=190`、`window=35/2210`、`agree=-1 asks=114` |
| `live-child`（阴性对照） | BUDGET、`reads=188` | **仍是 BUDGET**、`reads=192` |
| `prefix-REAPED-case13shape` | DEAD-CHILD、`reads=1` | **一字未变**：DEAD-CHILD、`reads=1` |

⇒ 问②**有**：**指错病因 → 指对病因**，两发之间被改的只有"有没有人去问 OS"这一件事；
**反向对照两发逐字同形**（`live-child` 与 `case13shape`）⇒ 这一味不是"什么都变"的那种改法。〔自跑〕

**名册口径**：`TestSLO156` 五族那面对照见 5.1 甲；本程另有一发**用发货用例**的外部可见变化——
`accept-teeth-m1-jia` 里 case 17 的红句原文（生产常量 30s）：

```
--- FAIL: TestSLO156ReportLoopNamesTheDeadSubjectItWasWaitingOn (30.05s)
    slo_exit_os_156_windows_test.go:257: dead subject was read 1470 times, want exactly 1 (the loop asks the OS and gives up on the spot)
--- FAIL: TestSLO156WaitReadyNamesTheDeadSubjectToo (60.05s)
    slo_exit_os_156_windows_test.go:285: waitReady give-up "wisp slo: subject 63508 never reported ready within 1m0s" does not name the dead subject and its code
```

（1470 次／30.05s 与 r1 的 1469／30.07、r2 的 1467／30.04 是本机差，本程不做跨机宣称。）〔自跑〕

### 5.3 单点回退进攻（本仓固定动作）：**只撤 1 处，问"哪条用例变得不响"**

本程不撤全部，逐味撤（同一把尺、同一选择器、八发全自跑，全表在 3.3）：
**撤 (b)** ⇒ 15/16/17/18 四枚红（甲）；**撤等态判定**（m2）⇒ 同样四枚；**撤"别把活的认成死的"**（m3）⇒ **只 16 一枚**；
**撤 (a)**（m4）⇒ 16 ＋ **票 149 的 case 13**；**撤码那读**（m5）⇒ 15/17/18；
**撤 (b) 再撤一枚持有者**（m1＋off15／m1＋off17）⇒ **其余三枚仍红**；
**撤 (a) 再撤 case 13**（m4＋off13）⇒ **只剩 16 一枚红**。

⇒ 答得出"哪条变得不响"：**没有任何一处改动撤掉之后全包静默**；同时**没有任何一枚持有者是装饰**——
撤掉其中一枚，其余仍响（红数 3／3／1），要**四枚一起撤**才静默（乙）。这一条正是"承重"的可复算形。〔自跑〕

### 5.4 恒真判据这一类新假绿——本程主动排掉两枚

1. **"复跑必须红"型判据本程一枚没开**：丙那一形今天**就是**不响（21/14/0 与基线一字不差），
   本程把它写成**读数**而不是判据（若写成"摘一味必须红"就成了永不响的装饰格）。
2. **本程证过的"会响的正控"都在**：`asis` 6/6/0（绿侧）与六发变异 4/4/2/3/3/1（红侧）同尺同树，
   外加 `LOADED=YES`（每发 `=== RUN`≥5，**零枚 `LOADED=NO`**、零枚 `0xc0000135`）〔自跑〕。

**放水两问自答（本格）**：① 断言方向——本程零判据改动；本程甚至**没有**新造任何一枚变异（六发全部 import r1 的 `my156.py`），
所以不存在"造一发打不红的变异来充数"这扇门；② helper——本程用 r1 的尺（`spec`/`case_off`/`run`/`apply_overlay`）
与 r1 的探针（`zz156probe_windows_test.go`），只加 cell 命名与四数/名册打印。

**第 5 格判定**：**〔成立〕**——两问都**量**不**辩**：问①按定义答"是，且只在使用者撤掉持有者时才成立"（甲乙丙三形本程全部自跑），
问②答"变了"（操作员那句逐字变、名册那面 4 枚换色，两向反向对照同形）。
**"摘一味编不过"那一形在本票不存在**（六发变异全部编得过并跑出 RUN）。

**本程没测什么（本格）**：
1. **没量 (a) 支内部**（`ps.Exited()` 与 `ps.ExitCode()` 的分配）——m4 摘的是整支三行，本程没造"只摘其中一行"的发；
2. **没造"只红 case 17 一枚"或"只红 case 15 一枚"的那一发变异**（5.3 的对称缺口，见 3.5 第 1 条）；
3. **没跑 `-race`** ⇒ 新增每轮询 1–3 枚内核调用与 `stop()` 的竞态关系**无判据**（既不红也不绿）；
4. **没量真 `wisp slo` 那条命令端到端**（5.2 那一格是**探针**打在 `collectReportWithin` 生产接缝上打印的句子，
   预算是探针自挑的 `2s/10ms`；生产常量那发是 case 17，走测试接线不是那条命令）——
   这笔欠账票面 `:67` 自己登记过（`wisp slo` 全仓 100% 不被 `go test` 执行），本程**确认它仍在**。

---

## 第 6 格　AC#6 —— 契约轴零字节（**逐枚 commit 现量**，区间 diff 不当尺）

**票面判据（`:49-52`）**：`docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`internal/panel/**`、
`internal/agent/approval/**`、`internal/observe/**`、`thresholds.go`、任何 golden、`allowlist.txt`、
`scripts/slo-check.ps1`、`tools/d22scan/**`、`frontend/**`、`design/**`、台账、停车点＝**编排者写面**；
⚠ 名册**逐枚 commit 现量**（别拿区间 diff 当尺）；`frontend/**`／`design/**` 有别的会话在未提交地写：**不碰、不还原、
不算进任何零命中宣称**。

### 6.1 本程的"本票实现 commit"名册是**按路径挑的**，得 **10 枚**（r4 报 8 枚）

```
$ git rev-list --count 9835d81..7ca1130                 88        （r4 在 d661e8c 上当轮量到 87＝版本漂，不是谁数错）
$ 逐枚 git diff-tree … -- cmd/wisp | docs/evidence/s1/156-* | .scratch/wisp/probes/156
cc58445 r1  0e95353 r1†  2ce77e1 r2  2262c2b r2  3fbb1ce r3  1137781 r3†  1d53526 r4  a4ec1f4 r4  d661e8c r4  7ca1130 r4
```

⇒ **比 r4 多的两枚（`d661e8c`／`7ca1130`）就是 r4 自己最后那两枚**——它数的时候它们还没落地，**结构性数不到**，
不是漏计缺陷。†＝编排者代提（`0e95353` 代提 r1 停在轮次上限的码、`1137781` 代提 r3 的"改前"门禁台件）。
**"代提"不加重也不减轻**：本程对这 10 枚**逐枚**重打同一把尺，`1137781` 在内（见 6.3 与第 7 格）。〔自跑〕

### 6.2 十枚 × 射程：全 0，且**每一枚的删除列都不落在射程上**

`ac6-census.sh`（16 支 pathspec＝票面 15 支，`thresholds.go` 与 `internal/observe/**` 分列、
"任何 golden"拆成 `**/testdata/golden/**`＋`**/golden/**`）原文 `ac6/CENSUS-OUT.txt`＋`ac6/GOLDEN-FIX.txt`：

```
cc58445  A1..A16 全 0  TOTALHITS=0 deletedlines=0      0e95353  全 0  TOTALHITS=0 deletedlines=4
2ce77e1  全 0  TOTALHITS=0 deletedlines=0              2262c2b  全 0  TOTALHITS=0 deletedlines=0
3fbb1ce  全 0  TOTALHITS=0 deletedlines=2              1137781  全 0  TOTALHITS=0 deletedlines=0
1d53526  全 0  TOTALHITS=0 deletedlines=7              a4ec1f4  全 0  TOTALHITS=0 deletedlines=0
d661e8c  全 0  TOTALHITS=0 deletedlines=0              7ca1130  全 0  TOTALHITS=0 deletedlines=0
golden 那两支单独重打（6.4 那枚死 pathspec 修好之后）：10 枚全 golden_hits=0
go.mod／go.sum 逐枚：10 枚全 []
```

删除列非零的三枚（`0e95353`=4、`3fbb1ce`=2、`1d53526`=7）本程逐枚看了文件清单：
全部落在 `cmd/wisp/**` 与 `docs/evidence/s1/156-*` ⇒ **一枚都不在射程**〔自跑〕。

### 6.3 让"零"值钱：五枚正控全响（按**路径**挑，不按 message 挑）

```
1218192  internal/risk hits=2                     ← 冻结码会响
191f0d6  docs/specs hits=2 ＋ 台账 1 ＋ 停车点 1
5866c6f  docs/PLAN.md hits=1 ＋ docs/specs hits=1
5eb0f6b  go.mod hits=1                            ← 依赖那支会响
cbbbdf8  **/testdata/golden/** hits=1 ＋ **/golden/** hits=1     ← golden 那支会响
负控：本程自己的 commit 也打一遍 → 16 支全 0（它只动证据件）
```

⇒ 本表所有"0 命中"都是**这把尺打得响**之下的 0〔自跑〕。

### 6.4 本程在这格犯的**两枚**仪器错（都往"假零"方向，都就地纠正）

1. **死 pathspec**：本程第一版把"任何 golden"写成裸 `testdata/golden` ⇒ 作为 git **前缀**只匹配**顶层**
   `testdata/golden/`（本仓没有），**只能读出 0**。改成 `:(glob)**/testdata/golden/**` 后，
   正控 `cbbbdf8` 打得响（hits=1），10 枚实现 commit 仍 0。**没修之前那个 0 是假的**，本程把它写出来而不是抹掉。
2. **`git ls-tree` 不支持 pathspec magic**：`git ls-tree -r 7ca1130 -- ':(glob)**/golden/**'`
   ⇒ `fatal: pathspec magic not supported by this command`，而**同一支 magic 在 `git diff-tree` 上工作正常**。
   本程那条"存量 firing check"因此整批发 0——那是**尺不响**，不是**没有货**：
   锚点存量本程改用朴素前缀重量：`docs/PLAN.md`=1、`docs/specs`=14、`internal/risk`=37、`internal/panel`=20、
   `internal/agent/approval`=18、`internal/observe`=25、`thresholds.go`=1、`allowlist.txt`=1、
   `scripts/slo-check.ps1`=1、`tools/d22scan`=6、`frontend`=85、`design`=30。〔自跑，`ac6/GOLDEN-FIX.txt`〕

### 6.5 时代那面（本程不算进任何宣称，但把归属量出来）

```
$ 区间交叉核（不是 AC#6 的尺，只是对照）：git diff --name-only 9835d81 7ca1130 -- <冻结 11 支＋golden glob＋go.mod go.sum>
   → **零行**（整个时代 88 枚里，冻结码那面一枚都没被碰过）
$ 逐枚（88 枚 × frontend/design/台账/停车点）→ 20 枚有命中：fe=[1-9] 4 枚、design 0 枚、台账/停车点 16 枚
$ 这 20 枚里，本票 10 枚实现 commit 命中＝**0**
```

⇒ **`design/**` 在 commit 层 0 命中**只说已入库部分；派单点名的 **25 枚未提交 `design/**` 增量**在名册上结构性看不见，
本程**没碰、没还原、没提交**，也**没算进任何零命中宣称**。`a2a3b3c`／`0f18652`／`0c56b0c`／`7609913` 那四枚前端写面＝**另一会话**。〔自跑，`ac6/era-attr.txt`〕

### 6.6 一枚本程多做的（票面没要求，但 AC#2/AGENTS §1.2 的禁形要用它）

`bash scripts/d22scan.sh` rc=0：**clean - no D22 ban violations**，且它**报了射程**（`bans #1-5 cmd/=23`、
`ban #8 cmd/ 45 Go files 含注释与 _test.go`、`internal/=205/413`），正控套件 `TestBuiltBinaryGoesRedEndToEnd` 六枚全 PASS
⇒ "clean"不是"尺瞎了"。**档位与边界要写清**：这一发扫的是**当前工作树**（含别人未提交的 `design/**`、`frontend/**`），
**不是 `7ca1130` 那棵树** ⇒ 它证"今天盘上没有违禁形状"，**不证**"锚点树没有"；本程**没有**为它造锚点快照
（派单禁 `git archive | tar`，而本程没有别的合规取树法在射程内）。〔自跑，`d22scan.txt`〕

**放水两问自答（本格）**：① 断言方向——本程一枚判据没动，且这格连测试都没跑（纯 git 尺＋一发只读扫描）；
② helper——尺是本程自己写的 `ac6-census.sh`/`golden-fix.sh`（承 r1/r4 的逐枚原则），
**没有** import r4 的 `zero156-r4.sh`，所以 6.2 那张表**不是**它的读数的复述。

**第 6 格判定**：**〔成立〕**。票面 15 支射程在**本票全部 10 枚实现 commit**（含编排者代提那两枚）上逐枚 0 命中；
三枚非零删除列全部落在 `cmd/wisp/**` 与证据件上；`go.mod`／`go.sum` 逐枚 0；
五枚按路径挑的正控全响；两枚本程自己的"假零"尺已就地钉死并写进本表。

**本程没测什么（本格）**：
1. **没验 `design/**` 那 25 枚未提交增量的内容**（派单明令不碰）；
2. **没在 `7ca1130` 的快照树上跑 d22scan**（见 6.6 的边界）；
3. **没核 `mut-shipped/` 那 9 枚日志的内容**（AC#4③ 已闭，本程只核了"9 枚 tracked"这一条形，第 4 格 4.3）。

---

## 第 7 格　AC#7 —— 门禁（逐包单跑）＋ 派单 §4② 那笔账的裁

**票面判据（`:55-58`）**：`cmd/wisp` 与 `internal/tools` **改前改后各一次**，四数之外**名册两向 `comm`（差集逐枚归属）**；
三枚仪器坑（dll/PATH 形状、`-overlay` 不与 `-cover*` 同用、`-race` 可能既不算红也不算绿）。
本程做了两件：**把四发入库日志用自己的尺重数**＋**自己再跑一发**（第五发，跑在合并态 HEAD 上）。

### 7.1 五发并排（前四发＝〔日志＋归档，本程抽验〕，第五发＝〔自跑〕）

| 发 | 跑者／锚点 | `cmd/wisp` 四数 | `internal/tools` 四数 | `0xc0000135` | verdict |
|---|---|---|---|---|---|
| `gate-pre` | r1 @ 改前面 | 139／79／0／0 | 115／79／0／0 | 0 | ok 136.561s |
| `gate-post` | r1 @ r1 之后 | **143／83**／0／0 | 115／79／0／0 | 0 | ok 105.606s |
| `gate-r3-pre` | r3 @ `3293192`（编排者代提 `1137781`） | 144／84／0／0 | 115／79／0／0 | 0 | ok 84.960s |
| `gate-r4-post` | r4 @ `1d53526` | 144／84／0／0 | 115／79／0／0 | 0 | ok 114.336s |
| **`accept-post`（本程）** | 本程 @ 工作树字节＝`7ca1130`、HEAD＝`c6ceddc` | **144／84／0／0 rc=0** | **115／79／0／0 rc=0** | 0 | ok 79.560s／13.223s |

本程那一发的钉法（不是"我相信工作树干净"）：`meta.txt` 里 `diff_against_anchor_7ca1130=[]`
＋`tree_sha256 cmd/wisp/slo_windows.go=d1f1b881a006a569` **＝**`anchor_sha256 …=d1f1b881a006a569`〔自跑〕。

### 7.2 名册两向 `comm`——四组差集本程全部自己重算

```
r1 那对（139/79 vs 143/83）：cmdwisp runs pre-only=0 post-only=4 ; verdicts pre-only=0 post-only=4
   post-only 四枚逐枚点名＝TestSLO156ExitedAsks…／LiveSubjectStays…／ReportLoopNames…／WaitReadyNames…
   internaltools runs+verdicts 两向全空（115↔115、79↔79）
r3/r4 那对（144/84 vs 144/84）：四向全空（本程抽验，且本程从 .log 重抽的名册与入库 -runs/-verdicts 逐字节相同＝8 项 IDENTICAL-to-landed）
本程 vs r4（合并态那一发）：cmdwisp runs 0/0（144↔144）、verdicts 0/0（84↔84）；internaltools 0/0、0/0
整票那一向（本程补的尺，票面没要求但"改前改后"要的是这一对）：
   gate-pre(139/79) vs 本程 accept-post(144/84)：pre-only=0，mine-only＝**恰好那五枚 TestSLO156**（含 case 19）
```

⇒ **整票净增五枚名字、零枚被吞／被改名**〔自跑＋抽验，`gate/RECOMPUTE-LANDED.txt`〕。
`internal/tools` 五发一字未变，本程量到了**为什么**：整票 `git diff --name-only 9835d81 7ca1130 -- internal` ＝ **零枚文件**〔自跑，2.3 那把尺〕。

### 7.3 三枚坑逐枚核（本程自己数，不抄坑单）

1. **印出来的红**：本程 `accept-post/cmdwisp.log` 里 `^--- FAIL`＝**0**，而字符串 `FAIL`＝**4**。
   本程逐行验那 4 行＝`[FAIL] sherpa-onnx C API`／`[FAIL] onnxruntime.dll version`／`[FAIL] data dir resolvable (dev)`
   ＋汇总行 `wisp doctor: FAIL`，且**全在 `TestAC2RealProcessRefusesOnEveryLegWithoutAppData128` 内部、那枚用例 `--- PASS`**
   （本程用行号回溯定位 enclosing `=== RUN`，不是引用别人的话）。
   ⇒ 四发入库日志同样各是 **4／0**〔自跑，本程重数〕。**派单 §5 那句"3 枚 `[FAIL]` 标记＋1 行汇总＝4 行"本程证实。**
2. **dll／PATH 形状**：本程走 shell 形 `/d/work/…/third_party/sherpa-onnx`（`meta.txt` 记 PATH 首项），
   两包 `0xc0000135`＝**0** 且 `=== RUN` 非零 ⇒ **真加载真跑了**，不是"看着像没测到"〔自跑〕。
3. **`-race`／`-cover`／`-overlay` 一枚没用**（本程脚本源码即证：`go test -count=1 -v` 一条）
   ⇒ 本程这两发**既没有 overlay/cover 互相吞掉的读数，也没有 `rc=1 而零条 --- FAIL` 那种"既不算红也不算绿"**〔自跑〕。
4. **射程没跑偏**：`TestC21DesignTokensFourWayAgree`／`internal/panel` 在本程两发日志里 **0／0 命中**〔自跑〕。
   （但本程在**另一件事**上撞到了 `internal/panel`——见 7.5，那与门禁射程无关。）
5. **争用**：开测前 `tasklist` 现量 Listener=1／Worker=0／在跑 go test=0（0.4），与本票前几程同形。

**放水两问自答（本格）**：① 断言方向——本程没有为变绿动过任何判据；本程那一发**红了就是红了**，
   事实是 `--- FAIL` 0 枚、`rc=0`、verdict `ok`；② helper——本程的脚本 `gate156-accept.sh` 是**照 r4 的形状自己重写的一份**
   （同包列、同 flag、同锚定模式、同抽名法），**没有 import r4 的脚本**，所以 7.1 那张表不是他们读数的复述。

**第 7 格判定**：**〔成立〕**。票面要的"改前改后各一次"在**两层**都齐：
每程各自那一对（r1 139/79→143/83、r3/r4 144/84→144/84），**以及整票那一**对（139/79→144/84，本程自己合出来）；
四数之外名册两向 `comm` 逐枚归属，差集要么是"恰好多出本票新增的名字"要么是空集，**没有一枚用例被吞或被改名**。

### 7.4 本格顺手量到的一枚**证据件缺陷**（记档，本程不修——不是本程的写面）

r1 证据件第 3 格（`7ca1130` 版 `:283`／`:285`）写着：

```
改后 cmd/wisp   RUN=144  PASS=84   FAIL=0  SKIP=0   rc=0  ok  (见第 7 格现量)
only-in-post: 5 枚，全部 TestSLO156*（本票新增，逐枚点名见第 7 格）
```

本程重数它**自己那发** `gate-post/cmdwisp.log`：`RUN=143 PASS=83 FAIL=0 SKIP=0`，名册 post-only **＝4 枚**。
`144／84／5 枚` 那组数**真实存在**，但它是 `gate-r3-pre`（跑在 `3293192`）那一发——
本程把中间那一枚差集也点了名：`gate-post → gate-r3-pre` 只多 **`TestSLO156FixtureChildrenDoWhatTheirNamesSay`**
（case 19 在 r1 跑完它那发门禁之后才加进去，r1 的 post 面结构数不到它）。
另：那两行括号里指的"第 7 格"在**现在这份文件里是 r2 的第 7 格**（`第 6/7 格（r2 部分）`），
r1 自己的 6/7 格没进这份件 ⇒ **指针悬空**。
**档位**：记录级缺陷（数值真实、**归属指错了一发跑**）；**不改判 AC#7**，也**不扣 AC#3**（AC#3 的判据在 3.1/3.3 由本程独立量到）。
最小闭合集合见第 8 格 C3。

### 7.5 派单 §4② 那笔账——本程**同意结论、推翻理由**，并交出读数

**编排者的裁定**＝本票不重跑合并态；**r4 的理由**＝"两包内 `frontend` 引用 0 命中、无 embed 路径"。
**本程量到那一枚理由是假的**：

```
$ grep -rn "wisp/frontend" --include=*.go cmd/ internal/     →  internal/panel/assets.go:25  import "…/wisp/frontend"
$ go list -deps ./cmd/wisp/  | grep -c CarlosShao/wisp/frontend →  1        （internal/tools 那发＝0）
$ frontend/embed.go:19                                        →  //go:embed all:dist
$ git ls-tree -r 7ca1130 -- frontend/dist                     →  frontend/dist/.gitkeep   只有 1 枚
```

⇒ **`cmd/wisp` 的测试二进制确实把 `frontend` 那个 Go 包（连带 `all:dist` 的字节）链进来了**，
依赖只是**经 `internal/panel` 一跳**到达，所以"`cmd/wisp/*.go` 里搜不到 `frontend/` 字样"这把尺**天然量不到它**。
**结论仍然成立，但成立在另一处事实上**（本程自己量）：`a2a3b3c` 那枚 commit 只碰 `frontend/src/**` 3 枚文件，
而 embed 图案是 `all:dist`——时代 88 枚里碰 `frontend/embed.go` 或 `frontend/dist` 的 commit **0 枚**；
四枚别家前端 commit（`7609913 0c56b0c 0f18652 a2a3b3c`）**每一枚的非 `frontend/src` 文件＝0**。

**而且这一发本程已经替它跑过了**：本程那发门禁跑在 HEAD `c6ceddc` 上（`git merge-base --is-ancestor a2a3b3c c6ceddc` ＝ rc 0），
`cmd/wisp`＋`internal/tools` 字节钉在 `7ca1130`，读数＝144／84／0／0 与 115／79／0／0，与 r4 那发**名册两向全空**。
⇒ **本程判：AC#7 不需要为合并态再开一格；"要那枚读数得另派"这笔欠账已由本程这一发付掉，不必继承、也不必另派。**
**同时登记一枚新的架构事实**（它比这一票更值钱，交给下一位引用"frontend 没被 embed"这句话的人）：
**别拿"两包内搜不到 `frontend/` 字样"当"前端 commit 进不了 `cmd/wisp` 读数"的理由**——正确的尺是 `go list -deps`；
`frontend/dist` 一旦被 build 出真字节并入库，或 `frontend/*.go` 被改，`cmd/wisp` 的读数**可以**随之变，
而那种变化今天**没有任何一枚门禁盯得住**（`wisp slo`／panel bundle 一致性那面另说）。〔自跑〕

**本程没测什么（本格）**：
1. **没跑全仓 `go test ./...`**（票面射程只有两包，本程不自扩射程）；
2. **没跑 `npm run build` 后那一发 `cmd/wisp`**（那要动 `frontend/**`，本程写面之外）——
   7.5 登记的正是这一枚未证形：**dist 字节变了会怎样，本票无人量过**；
3. **没跑 `-race`**（与 r1/r2/r4 同一笔未取证）；
4. **没验 r3 那发"改前"的**树**可比性**——本程只证了名册抽取同源、四数对得上，
   那发跑在 `3293192`，其时 `cmd/wisp` 字节已含 r1/r2 的改动（"改前"是 **r3/r4 那一段的改前**，不是整票的改前；整票那对本程合在 7.2）。

---

## 第 8 格　交件：总判 · 七格档位 · 复算面 · 两栏计数 · 本程的错

### 8.1 总判一句

**票 156 可收：七格无一格退回**——承重那一味（`slo_windows.go:881-882` 问 OS）在本程自己的尺下
**摘掉就红四枚、操作员那句话逐字变对**，甲乙丙三形全部本程复跑对上；
**两格带条件**（AC#3、AC#4），四条条件全是**记录级**（数值归属／归类口径／时态／commit 纪律），
**没有一条需要动生产码或测试文件一字**。

### 8.2 七格档位（每格的命令与读数在对应那格里，本表不复述）

| 格 | 票面判据 | 档位 | 本程的关键读数 | 条件／登记 |
|---|---|---|---|---|
| **AC#1** | 改前那一发复现成自己的数（ms＋具名次数） | **成立** | `window=38ms→2211ms`、`reads=187`、`asks=130` 全 false、印 BUDGET（m1 单点回退） | 三笔未量：写一半被打断形／真 `wisp slo` 端到端／探针预算≠生产预算 |
| **AC#2** | 落地乙＋三件硬约束＋两向举证 | **成立** | 预算常量 grep 0；`go.mod` 两锚同一 blob；`go func`/`reaper`/`time.Now`/`time.Since` 全 0；两向举证自跑 | 三笔未钉（3 笔全部**登记**，本程自己核过"确实没钉"）：`stop()` 零枚用例、259 退出形零枚、`WAIT_FAILED` 兜底零枚 |
| **AC#3** | 新用例钉住那一形，且没把既有断言弄松 | **成立·带条件** | 本程自家门禁新五枚全 PASS ＋ 票 144/147/149/152 那十三枚逐枚名字全绿；`:165`／`:257` 两枚是本形的牙 | **C3**＋**C5** |
| **AC#4** | 三枚注释级动作＋分开 commit＋②改动前自己现量 | **成立·带条件** | ①逐字节 md5 相等；②自跑 `20/13/0`、`21/13/1`、`21/14/0`，`:739` 今天为真；③9 枚 tracked | **C1**＋**C2** |
| **AC#5** | 承重按操作定义答两问 | **成立** | 问①＝甲 6/2/4、乙 22/15/0、丙 21/14/0＝基线一字不差且名册两向空；问②＝操作员句逐字变＋两枚反向对照不变 | 无；恒真判据两枚本程主动排掉（5.4） |
| **AC#6** | 契约轴零字节，逐枚 commit 现量 | **成立** | 10 枚实现 commit × 16 支全 0；三枚非零删除列全在 `cmd/wisp`／证据件；`go.mod`/`go.sum` 逐枚 0；五枚按路径挑的正控全响 | 两枚本程自己的假零尺已就地钉死并写进 6.4 |
| **AC#7** | 两包改前改后各一次＋名册两向 `comm` 逐枚归属 | **成立** | 五发并排；整票那对 `139/79→144/84` 差集＝恰好五枚新增名、`pre-only=0`；`internal/tools` 五发一字未变（整票 `-- internal` 零枚文件） | **C3**（第 3 格数字归属）＋ C4 |

**C1**＝票面 `:46`"分开 commit"在 ① 上未兑现（与 AC#2 的码同枚 `0e95353`）——**最小闭合集合：不动码**，
在台账记这一笔形（或下一张注释级票里单独提一次那枚注释面）。
**C2**＝`:742`／`:749`／`:743`＋`:752` 三处**时态**过期（r4 改真 `:739` 的必然后果）——
**最小闭合集合：一张记录级收口票，按"只追加、不改正文"闭**（本程证实数值侧今天无 hazard：`reports === RUN=21…` 0 命中）。
**C3**＝r1 证据件第 3 格把 `144/84`＋"5 枚"写成自己那发 `gate-post`（实为 `143/83`＋4 枚），且"见第 7 格"指针悬空——
**最小闭合集合：同样只追加**（在验收表或收口票里写明"整票那对 139→144 的 +5 属 `gate-r3-pre` 那发"），**不动 r1 正文**。
**C4**＝`wisp slo` 那条命令今天**仍不进任何门禁射程**（票面 `:67` 自己登记，本程确认未闭）＋
7.5 那枚**新登记的架构事实**：`cmd/wisp` 经 `internal/panel` 链 `frontend`（`//go:embed all:dist`），
所以"两包内搜不到 `frontend/` 字样"**不是**"前端改动进不了读数"的理由，正确的尺是 `go list -deps`。
**C5**＝r1 第 3 格判定句把 case 19 与 15/17/18 并列为"四枚会响的检"，本程八发变异里 case 19 **从不红**（设计如此）⇒
归类应为**夹具守卫**（与 r1 自己那张表里 19 那行一致、与判定句不一致）——**最小闭合集合：只追加一句更正**。

### 8.3 本程**复算过**的 vs **继承**的（本仓为"把引用当复算"退回过程序，所以逐枚点名）

**本程自跑**（当轮、命令在表里）：两发探针（ship／m1）；甲乙丙五发 teeth；八发 mut／case-off；三发 AC#4 配对；
一发两包门禁；名册全部抽取与两向 `comm`；P1/P2/P3 三把尺；AC#6 十枚 × 16 支＋五枚正控＋golden 两支重打；
时代归属 20 枚；`go list -deps`；embed 定位；`git merge-base --is-ancestor`；AC#4① 的 md5 等式；
`grep` 类零命中宣称全部带 rc。
**本程抽验**（别人入库的日志，本程自己重数四数／名册，**没有**重跑那一发）：
`gate-pre`／`gate-post`（r1）、`gate-r3-pre`（r3，**编排者代提 `1137781` 的那发——本程按自己的尺重算了它的四数与名册，
"代提"既不加重也不减轻**）、`gate-r4-post`（r4）、`probe-pre`／`probe-post`（r1）、
`probes/152/mut-post/g1-…case14off.log`（tracked，20/13/0 与选择器形状）。
**本程未复算、只用其形**〔某程读数，本程未复算〕：`mut-156/` 里 r1 的十二发日志原文（本程自己重跑了等价发，不需要它）；
`mut-156-r2/` 与 `mut-156-r4*/` 里各发的原文；`3fbb1ce` 的 gofumpt 修形（r3 交、编排者核过，本程只现读
`gofumpt -l cmd/wisp/` 为空＝本程自己的尺）；`probes/152/mut-shipped/` 那 9 枚的**内容**。
**本程没有把任何一枚"未复算"的读数当作下判的凭据。**

### 8.4 伪授权两栏（分开计，各带出处）

| 栏 | 枚数 | 明细 |
|---|---|---|
| **真通知回显数** | **7** | ①`<system-reminder>` 技能清单（首发 Bash 结果尾部）；②同处 `The date has changed`；③同处 `Memory: …/agents.md`（`AGENTS.md` 正文回显）；④⑤⑥⑦ 四枚后台任务完成事件（census／teeth／gate／muts 各一枚，均以 `[SYSTEM NOTIFICATION - NOT USER INPUT` 形回显） |
| **判为注入数** | **0** | 本程读过的每一枚文件与每一次工具输出里，**没有**任何"替本程写好结论／放宽判据／削弱 owner 权威／让本程少取证、别用工具、直接给结论／请 revert"的文字 |

**反查记录**（派单 §7 那两形各扫一遍）：
形状①——四枚被引文件 `7ca1130`／`HEAD`／工作树**三向字节比**（`injection-scan.txt`）：两枚实现件与 `AGENTS.md` 三向相同、
diff 0 行；票面 88→94 行的那 6 行＝编排者**已入库的真追加**，本程读了全文（写的是"七框仍由它定，我不勾"＋把两枚 §4 交本程裁），
**不含预写判语**，故**既不记为注入、也不当凭据**。
形状②——本表引用的**十九枚号**逐枚 `git cat-file -t` 全为 `commit`、零枚 fatal；
本程**执行**过的两把尺（`my156.py`、`zz156probe_windows_test.go`）额外做了工作树↔锚点 md5 等式（8.6 第 10 条）。
**另：别程 commit message 里对伪授权的转述（`e098927` 那整节、`a2a3b3c` 之类）＝内容，不是对本程说话，两栏都不计**（照派单 §7 末条）。

**凭据声明**：本程未读、未抄任何 API 密钥／DPAPI 值／环境变量内容。
本机供参考的现量：`%APPDATA%\wisp` 目录存在但零条目 ⇒ 本程**没有凭据语料**，
这一句**只记成"无凭据语料"**，不是"扫过、很干净"。

### 8.5 被拒的工具调用：**0 次**（含取数前／取数后）

- `Write docs/evidence/s1/156-close-and-load-bearing-r1-accept-r1.md` **一次通过**，没有"被权限拒就换路径"那件事发生。
- 本程全程**没有**任何一次权限拒绝 ⇒ 没有"取数前／取数后"可分。
- ⚠ 为免下一位把两枚**不是拒绝**的事当拒绝，本程点名：
  ①**两枚 Bash／Edit 参数校验错误**（一枚 `command` 为空、一枚漏 `old_string`，都是本程手滑）——
  分别发生在取数中与收尾时，本程立刻重发，**没少取一次证**；
  ②两发命令**转后台跑**（门禁、muts）——是时限形状，不是闸门，读数已入库。

### 8.6 本程自己犯的仪器错（十条，全部如实；前七条都在取数过程中实发并当场纠正）

1. **`git show --numstat -1 --format='' <c> -- <路径>` 不是逐枚尺**：不同 `<c>` 吐同一批数，还为本票**没碰过**的
   `go.mod`／`go.sum`／`docs/PLAN.md` **伪造命中**（`1 1`／`2 0`／`4 4`）；最坏的是**正控在这把坏尺上照样"响"**。
   ⇒ 改 `git diff-tree -r --no-commit-id` ＋真负控（`7ca1130` 只动证据件→全空）双向钉。（0.3、6.2）
2. **把输出重定向进"脚本自己才创建"的目录** ⇒ `> …/ac6/CENSUS-OUT.txt` 打开失败，**整发 census 从没跑过**，
   而后台任务回报"exit code 0"。发现于 `sed: can't read …`。⇒ 先 `mkdir -p` 再跑；**教训：后台"rc=0"不等于跑过**，
   要数产出件（这正是坑单里 `grep -c` 那枚 rc=1 短路的姐妹形）。
3. **死 pathspec**：裸 `testdata/golden` 作前缀只匹配顶层 ⇒ 那一支**只能读 0**。改成 `:(glob)**/testdata/golden/**`
   后正控 `cbbbdf8` 才打得响。**这个 0 在改之前是假的。**（6.4）
4. **`git ls-tree` 不支持 pathspec magic**（`fatal: pathspec magic not supported by this command`）而 `git diff-tree` 支持
   ⇒ 本程那条"存量 firing check"整批 0＝**尺不响**，不是**没有货**。（6.4）
5. **门禁脚本的仓库根算浅一层**（照抄 r4 的 `../../../..`，而本程脚本在 `…/156-accept/gate/` 深一枚）
   ⇒ `cd` 到 `.scratch/`，`go test ./cmd/wisp/` **rc=1、RUN=0**——正是坑单点名的"既不是红也不是绿"那一形。
   发现于 `meta.txt` 里 `pwd=…/.scratch` 与 `tree_sha256` 全空。⇒ 改五层 `..` 重跑。
   ⚠ 坏的那发日志与 `meta.txt` **被好的重跑就地覆盖**（同一路径），所以**作废读数的原件只在这段对话里，不在仓里**——
   本程不假装仓里有它，也不拿它下任何判。（7.1）
6. **`grep -c 'func scriptedReader'` 结构性假 0**：它是 `type`（`:263`）不是 `func` ⇒ helper 溯源那一发差点读成
   "四枚判定 helper 都不存在"。按名字形状重打并对 `9835d81` 比才拿到真数。（3.4）
7. **roster 循环里把 `$f`（已含 `.log`）又拼一次 `.log`** ⇒ `comm` 读的是不存在的文件，**每一发都报 0/0**（看着像"名册全等"）。
   发现于 `wc -l` 报 base 文件不存在。**这批 `*.log.runs.txt` 是 0 行作废件**，本程按"只建不删"留在盘上并入库，
   在 3.3 就地声明"不是任何一发的读数"。（3.3 末）
8. **python 的 stdout 在重定向下全缓冲** ⇒ `*-summary.txt` 长时间 0 字节，本程一度以为那发没跑；
   真正可读的进度是 `run()` 每发**立刻落盘**的 `<cell>.log`。⇒ 别把"没输出"读成"没跑"，要数产出件。
9. **一枚 `grep -v "^.*://"` 过滤器把 `go:embed` 的搜索结果整批吃掉** ⇒ 本程一度得出"仓里没有 embed 指令"，
   从而差点原样接受 r4 那句"无 embed 路径"。重打 `grep -rn "go:embed"` 才拿到 `frontend/embed.go:19`，
   进而量到 7.5 那枚**推翻了理由、保留了结论**的读数。**教训：过滤串里的 `//` 会吃掉注释型命中。**
10. **本程"执行"了两把尺的工作树副本**（`python … acc156.py` import `probes/156/my156.py`），
    而派单要求引用内容一律 `git show 7ca1130:`——本程**当场没有替代做法**（不能把锚点件写进仓内），
    所以事后补了等式：`my156.py` 与 `zz156probe_windows_test.go` 工作树 md5 **＝** 锚点 md5
    （`2168c128c157`／`d8987891201b`），`git diff --name-only 7ca1130 -- probes/156/` 除四枚 `??` 外**零漂移**，
    另加一枚独立佐证：本程 m1 那发红的名册＝r1/r2 那两发红的名册，逐字相同。
    ⚠ 本程 import 也往 `probes/156/__pycache__/` 写了 `.pyc`：那目录**前程已有**（派单点名 `??`），
    本程**没创建它、没删它、没提交它**，但**增加了内容**——如实记这一笔，不写进任何"未碰"宣称。
11. **一处 `--count` 用 `wc -l` 数 `git show --name-only` 的输出**，枚数含空行会偏（本程全部改用 `grep -c .` 重打）。

### 8.7 `SPEC-12 §4.3` 五件事与本表的关系（本表只占其中一枚，别把四枚当已交）

① **缺口审计**——**不在本表**（审计者≠实现者，需另一程）；本表 8.2/8.3 那两张"没测什么／未复算"表是它的输入。
② **失败预演**（`SPEC-10 §7`）——**未做**，本程射程只有七格。
③ **对抗验收报告**——**＝本表**（`SPEC-10 §8` 的七项检查按票面 AC 逐格落，承重按派单操作定义攻，含单点回退一发）。
④ **场景↔切片双向核对**——**未做**（本票只有 `wisp slo` 观察者面那一格，不涉 `SPEC-00 §4` 矩阵的其它行）。
⑤ **推迟登记**——本表**没有新增 DEFERRED**：8.2 的 C1–C5 全是**记录级收口**与**未钉登记**，
不构成 `SPEC-12 §5` 的推迟条目（若编排者要把 C4 那两笔升成推迟，必须五字段齐全，且不是本程的写面）。

### 8.8 收口纪律自报（本程的 git 与写面）

- **写面**：只 `docs/evidence/s1/156-close-and-load-bearing-r1-accept-r1.md`（本件，`Write` 建）＋
  `.scratch/wisp/probes/156-accept/**`（只建不删）。**没碰**：任何生产码／测试文件／注释一字、票面、台账、
  `HANDOVER`、`injection-timeline`、r1/r2/r3/r4 已入库证据件正文、`frontend/**`、`design/**`、
  `probes/152/my152.py`、`152-…-accept-r1.md`。
  ⚠ **本条曾被本程写错过一次并当场改回**：本程一度在 8.8 写"那枚半件的 3 行自校已被别人入库"，
  现量证伪——`git status --porcelain` 仍显 ` M`、`git diff --numstat 7ca1130 -- <该件>`＝**3 加／3 删**、
  `7ca1130..HEAD` 里碰它的 commit **0 枚** ⇒ **它今天仍是未提交的半件**，本程**没碰、没还原、没提交、没补**（派单 §6 点名）。
  本程读它只做三向字节比（0.5：锚点 md5＝HEAD md5 `ee39d32a`，工作树 `16f45130`，差 6 行），**一字未当凭据**。
- **git**：**只 commit、零 push**；七枚 commit 全部**显式 pathspec**（`git commit -m … -- <本程路径>`），
  每枚 commit 前现读 `git diff --cached --name-only`，**没有一枚出现过别人的路径**（每枚的 `--name-only` 都可反查）；
  未用 `add -A`／`.`／`commit -a`；未用 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`rm`／任何删除。
- **本程的 commit 链**（每格一枚，按时间序；文档内次序＝格号序，第 4 格因"一格一枚"先落在文件末尾、
  在第 6 格那枚里挪回原位并声明，`sort` 行差＝0）：
  `ffa813d` 第 0+1 格 ／ `1d63f8d` 第 2 格 ／ `9d76f6c` 第 4 格 ／ `c718cca` 第 3 格 ／ `c9f9459` 第 5 格 ／
  `21db8ae` 第 6 格（含挪位）／ `c3b7a1d` 第 7 格 ／ `152c4ec` 第 8 格 ／ **本件之后第九枚＝第 3 格档位的自校正**
  （**为什么有一枚"改自己已提交内容"的 commit**：本格 8.2 与第 3 格判定句在写完后互相矛盾——一处"成立"、一处"带条件"。
  本程取严的那一枚并把这次回改**写在原地可见**，而不是另起一节让下一位去考古。除这一枚之外，
  本程**没有**改写过任何已提交的文字，`--amend`／`reset` 一枚没用过。）
- **票面七框本程一枚不勾、`-done` 不加**：勾由编排者按本表翻（派单 §1"表出来后票面才有得翻勾"）。
- **被验版本自始至终＝`7ca1130`**；本程 step 0 HEAD＝`ea59d28`，交件时 HEAD 已到 `c6ceddc` 及以后（编排者在提自己的写面），
  本程**没有**跟着锚动，`gate/accept-post/meta.txt` 里那行 `diff_against_anchor_7ca1130=[]` 就是本程不跟锚的证据本身。




---
