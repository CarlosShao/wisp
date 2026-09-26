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
