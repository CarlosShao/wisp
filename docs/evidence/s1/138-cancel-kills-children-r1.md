# 138 实现件 r1 — 取消一轮任务时那一轮自己起的子进程：先把"无归属"变成读数，再裁修法

- 工单：`.scratch/wisp/issues/138-cancelling-a-task-kills-nothing-the-task-itself-spawned-single-jobscope-no-terminatejobobject.md`
- 交件程：**实现程**（本件作者＝写码那枚，**本件不是验收表、本程不给自己打勾**）
- 交件时刻：2026-09-25 22:5x–23:0x +08
- 落点（编排者简报）：`internal/agent/**`
- 本件每一条读数都是**本程现跑**；票面与简报里的读数只用作"它说了什么"的对照。

---

## 0　锚点、环境、工作树（进场第一读）

**判据**：HEAD 必须是简报给的那枚 sha；脏工作树里别人的东西一律不碰、不 add、不还原。

### 0.1　命令原文与读数

```
$ git rev-parse HEAD
64858d6838ced46fbe7bcce38dc4f7bb163d2f9c

$ git branch --show-current
dev

$ git log --oneline -6
64858d6 docs(台账 A265): 票 147 结线五格全成立；验收程造出四发，其中一发我拒绝开格、写成"不许开格"的限定语
de1ee7b ticket(147 结线＋144/147 改名＋149 新立): 五框按非实现者表翻勾；同族三处零牙出口立成下一票
e9412c9 evidence(147 对抗验收 r1 第 9 格): 落盘名册七枚、仓外临时件 35 枚只建不删、复跑最小路径与变异尺字面量自锁
1642886 evidence(147 对抗验收 r1 第 5-8 格): 本程没测什么八条、结论修正记录十一条（含推翻本程自己三句）、伪授权两数分栏、五格总裁决
c107761 evidence(147 对抗验收 r1 第 4 格): 门禁五发独立重跑——CI 同形两边四数、名册 comm、假红 8 枚的来源、naive/anchored 归属差、d22scan 分母随树
c62b884 evidence(147 对抗验收 r1 第 3 格): 四件报回逐件重走——corrupt 支零牙(B14)/complete 支有牙(B15)、append-only 与 31/2 相符、offset:-1 实为四处、锚点三枚零代码

$ git status --porcelain
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
 D design/index.html
 D design/screens/approval.html
 D design/screens/ball.html
 D design/screens/chat.html
 D design/screens/config.html
 D design/screens/cost.html
 D design/screens/firstrun.html
 D design/screens/palette.html
 D design/screens/privacy.html
 D design/screens/security.html
 D design/screens/states.html
 D design/screens/tasks.html
 M docs/reports/frontend-session-log-zcode.md
 M frontend/scripts/gen-tokens.mjs
 M frontend/src/App.tsx
 M frontend/src/components/nav-rail.tsx
 M frontend/src/components/panel-skeleton.tsx
 M frontend/src/main.tsx
 M frontend/src/styles/theme.css
?? .zcodeignore
?? design/doubao/01-ball-states.jpg
?? design/doubao/demo/lib/
?? design/doubao/demo/rb-files.js
?? design/doubao/demo/rb-plugins.js
?? design/doubao/demo/rb-review.js
?? design/doubao/demo/rb-terminal.js
?? design/doubao/demo/rightbar.js
?? design/doubao/demo/screens/home.js
?? design/doubao/demo/screenshots/
?? design/doubao/demo/sidebar.js
?? design/old/
```

**读数**：HEAD **就是**简报给的 `64858d6`，分支 `dev`，没有别人的 commit 插进来。
工作树脏在 `design/**`（owner 自己挪动的 16 枚删除）、`frontend/**`＋`docs/reports/frontend-session-log-zcode.md`
（另一枚前端会话在写）、`.zcodeignore`（未跟踪）——**本程一枚未碰**；`cmd/wisp/**` 归同时在跑的票 149，本程未写一枚。

### 0.2　仪器版本（现跑，不背读数）

```
$ go version
go version go1.27.1 windows/amd64

$ D:/work/base/gopath/bin/gofumpt.exe --version
v0.12.0 (go1.27.1)
```

⇒ `gofumpt` 与票面 AC#4 那句"宿主已装 v0.12.0"相符；**本程未执行任何 `go install`**。

---

## 1　AC#1① —— 那一枚 `JobScope` 的创建点数与调用者（现量）

**判据**（票面 AC#1① 原文）：*"那一枚 `JobScope` 的创建点数与调用者（生产码里几枚、测试里几枚；是不是整机共用一枚）"*。

### 1.1　创建点：全仓只有 `OpenJobScope` 一处生产者

```
$ git grep -n 'OpenJobScope' -- '*.go'
internal/proc/boot_windows.go:87:	job, err := OpenJobScope()
internal/proc/jobscope_windows.go:60:// JobScope owns the process-wide Job Object. Create with OpenJobScope; the
internal/proc/jobscope_windows.go:68:// OpenJobScope creates the Job Object with KILL_ON_JOB_CLOSE. The caller
internal/proc/jobscope_windows.go:71:func OpenJobScope() (*JobScope, error) {
internal/proc/jobscope_windows.go:85:	return &JobScope{handle: h}, nil
internal/proc/jobscope_windows_test.go:107:	job, err := OpenJobScope()
internal/proc/jobscope_windows_test.go:109:		t.Fatalf("OpenJobScope: %v", err)
internal/proc/jobscope_windows_test.go:157:	job, err := OpenJobScope()
internal/proc/jobscope_windows_test.go:159:		t.Fatalf("OpenJobScope: %v", err)
internal/proc/jobscope_windows_test.go:200:	job, err := OpenJobScope()
internal/proc/jobscope_windows_test.go:202:		t.Fatalf("OpenJobScope: %v", err)

$ git grep -n '&JobScope{\|JobScope{}' -- '*.go'
internal/proc/jobscope_windows.go:85:	return &JobScope{handle: h}, nil

$ git grep -n 'CreateJobObject' -- '*.go'
internal/proc/jobscope_windows.go:72:	h, err := windows.CreateJobObject(nil, nil)
internal/proc/jobscope_windows.go:74:		return nil, fmt.Errorf("proc: CreateJobObject: %w", err)
scripts/spike/common/winshell.go:50:	procCreateJobObjectW        = modKernel32.NewProc("CreateJobObjectW")
scripts/spike/common/winshell.go:454:	job, _, err := procCreateJobObjectW.Call(0, 0)
```

**读数**：

| 形状 | 枚数 | 逐名 |
|---|---|---|
| `JobScope` 结构体字面量 | **1** | `internal/proc/jobscope_windows.go:85`（在 `OpenJobScope` 内部）⇒ 没有任何绕过 `OpenJobScope` 造 JobScope 的后门 |
| `OpenJobScope()` **调用点：生产码** | **1** | `internal/proc/boot_windows.go:87`（`proc.Boot` 第 3 步） |
| `OpenJobScope()` **调用点：测试码** | **3** | `internal/proc/jobscope_windows_test.go:107 / :157 / :200` |
| 产品码里的 `CreateJobObject` | 1 | 同上 `:72`；`scripts/spike/**` 那两枚是 spike 夹具，不在任何构建产物里（`scripts/` 非 Go 包路径的一部分） |

⇒ 生产侧**每个进程恰好一枚** Job Object：`Boot` 是唯一创建者，没有任何"每任务一枚"的第二创建点。

### 1.2　调用者：这枚整机 JobScope 今天被谁用

```
$ git grep -n 'NewTreeSampler' -- '*.go'
cmd/wisp/slo_other.go:16:// process tree by proc.NewTreeSampler / proc.NewExternalSampler for a specific
cmd/wisp/slo_windows.go:281:	sampler := observe.NewSampler(proc.NewTreeSampler(rt.Job), rt.Registry)
cmd/wisp/slo_windows.go:297:		run.Settle = runSettle(ctx, sampler, proc.NewTreeSampler(rt.Job))
internal/proc/treemetrics_windows.go:36:// NewTreeSampler binds the sampler to a JobScope.
internal/proc/treemetrics_windows.go:37:func NewTreeSampler(job *JobScope) *TreeSampler {

$ git grep -n 'rt.Job' -- '*.go'
cmd/wisp/slo_windows.go:281:	sampler := observe.NewSampler(proc.NewTreeSampler(rt.Job), rt.Registry)
cmd/wisp/slo_windows.go:297:		run.Settle = runSettle(ctx, sampler, proc.NewTreeSampler(rt.Job))
cmd/wisp/slo_windows.go:498:	p, err := rt.Job.StartInJob(cmd)
internal/proc/boot_windows.go:91:	rt.Job = job
internal/proc/boot_windows.go:146:	if rt.Job != nil {
internal/proc/boot_windows.go:147:		hooks.CloseJob = func(context.Context) error { return rt.Job.Close() }
internal/proc/boot_windows_test.go:43:	if rt.Job == nil || rt.Job.Closed() {
internal/proc/boot_windows_test.go:58:	if !rt.Job.Closed() {

$ git grep -n 'CloseJob' -- '*.go'
internal/proc/boot_windows.go:147:		hooks.CloseJob = func(context.Context) error { return rt.Job.Close() }
internal/proc/boot_windows_test.go:54:	jobRec := records[StepCloseJob-1]
internal/proc/boot_windows_test.go:55:	if jobRec.Step != StepCloseJob || jobRec.Skipped || jobRec.Err != nil {
internal/proc/shutdown.go:45:	StepCloseJob
internal/proc/shutdown.go:88:	CloseJob              func(ctx context.Context) error // step 9 (proc.JobScope.Close)
internal/proc/shutdown.go:183:	run(StepCloseJob, hooks.CloseJob, 0)
internal/proc/shutdown_test.go:27:	hooks.CloseJob = mk(StepCloseJob)
internal/proc/shutdown_test.go:39:		StepFreeOSMemory, StepCloseJob,
internal/proc/shutdown_test.go:89:		CloseJob:              mk(StepCloseJob),
internal/proc/shutdown_test.go:107:	for _, mustRun := []ShutdownStep{StepStopAudio, StepReleaseSpeechSessions, StepFreeOSMemory, StepCloseJob} {
scripts/spike/common/winshell.go:451:// AttachToKillOnCloseJob assigns the current process to a new Job Object
scripts/spike/common/winshell.go:453:func AttachToKillOnCloseJob() (windows.Handle, error) {
scripts/spike/shell-baseline/main.go:138:		j, err := common.AttachToKillOnCloseJob()
scripts/spike/xy-verdict/main.go:95:	job, err := common.AttachToKillOnCloseJob()
```

**读数**：生产侧对那枚整机 JobScope 的调用者共 **三处、两个用途**——

| # | 调用者 | 用途 | 是不是"任务级" |
|---|---|---|---|
| 1 | `cmd/wisp/slo_windows.go:281`、`:297` | D32 SLO 采样（`TreePrivateBytes` / 进程树读数） | 否，进程级度量 |
| 2 | `cmd/wisp/slo_windows.go:498` | `wisp slo` 把**被测子进程**塞进整机 Job（`StartInJob`） | 否，SLO 测量夹具 |
| 3 | `internal/proc/boot_windows.go:147` | 关停序列第 9 步 `Close()`（`KILL_ON_JOB_CLOSE` 连带清理） | 否，**进程退出路径** |

⇒ 审计代理那半截（"全进程共用一枚、有生产调用者"）**成立**：整机一枚，创建点 1 枚，调用者 3 处。
票面 §0 只敢断到"声明层"的那半，本格把它量成了读数。

⚠ 归属说明：`CloseJob` 那发还命中 `scripts/spike/**` 三枚，那是另一枚符号
`common.AttachToKillOnCloseJob` 的**子串**命中（spike 夹具，不在构建产物里），本表按符号名把它排除，
**不是**把它们算成"零命中"。

### 1.3　放水两问自答

- **断言方向动没动**：没动——本格不写断言，只做普查；普查命令是票面 AC#1① 给的三个形状的原文。
- **helper 是不是原有的那枚**：本格无 helper、无测试改动。
- 另答一条本仓的坑：**"零命中"宣称的分母**。`git grep ... -- '*.go'` 的射程是全仓跟踪的 `.go`（含 `scripts/spike/**`、`tools/**`），上表把落在射程里的**每一枚**命中都逐名贴出，包括两枚 spike 命中——没有把它们悄悄算成"无"。

### 1.4　本格没测什么（按"漏了它谁会先被骗"排序）

1. **没在运行时数过 Job 里真有几枚进程**。读数全部来自静态普查；`TreeProcessCount()` 在真机上会返回几，本程没量（量它要 `wisp slo`，那会另起两枚子 `wisp.exe` 并抢 CPU，且 `cmd/wisp` 此刻是票 149 的写者）。
2. **没穷尽"绕开 `OpenJobScope` 拿到第二枚 Job Object"的写法**。本格的 `CreateJobObject` 普查能抓到 `windows.CreateJobObject`，也能抓到 lazy-proc 那族（`scripts/spike/common/winshell.go:50` 的 `NewProc("CreateJobObjectW")` 就是被它抓到的——符号名是**字面串**，grep 打得着）。打不着的是**把符号名拼出来**的写法（`"Create"+"JobObjectW"`）与不经这些名字的第三方封装；本程没排那条，因为仓内没有任何一枚这样的现存代码，排它要先造一台仪器。

---

## 2　AC#1② —— 仓内有没有 `AssignProcessToJobObject`（及等价形状）；有 ⇒ 谁被分配进去过

**判据**（票面 AC#1② 原文）：*"仓内有没有 `AssignProcessToJobObject`（及等价形状）；有 ⇒ 谁被分配进去过"*。

### 2.1　命令原文与读数

```
$ git grep -n 'AssignProcessToJobObject' -- '*.go'
internal/proc/jobscope_windows.go:100:		if err := windows.AssignProcessToJobObject(j.handle, windows.Handle(handle)); err != nil {
internal/proc/jobscope_windows.go:101:			assignErr = fmt.Errorf("proc: AssignProcessToJobObject(pid %d): %w", p.Pid, err)
scripts/spike/common/winshell.go:490:	if err := windows.AssignProcessToJobObject(jobHandle, windows.CurrentProcess()); err != nil {

$ git grep -nE '\.Assign\(|\.StartInJob\(' -- '*.go'
cmd/wisp/slo_windows.go:498:	p, err := rt.Job.StartInJob(cmd)
internal/proc/jobscope_windows.go:117:	if err := j.Assign(cmd.Process); err != nil {
internal/proc/jobscope_windows_test.go:41:		if _, err := job.StartInJob(cmd); err != nil {
internal/proc/jobscope_windows_test.go:207:	if _, err := job.StartInJob(cmd); err == nil {

$ git grep -nE 'CREATE_SUSPENDED|ResumeThread|CreationFlags' -- '*.go'
cmd/balldebug/diff_windows.go:412:	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x00000200} // NEW_PROCESS_GROUP
cmd/wisp/resident_sink_nail_127_windows_test.go:200:	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.CREATE_NEW_PROCESS_GROUP}
```

### 2.2　读数

**有。** 产品码里唯一一处 `AssignProcessToJobObject` 在 `internal/proc/jobscope_windows.go:100`，包在 `(*JobScope).Assign` 里；
`Assign` 的调用者**只有一枚**（`StartInJob` 内部那处，`:117`），**外部直调零枚**
（`\.Assign\(` 的普查里除 `:117` 再无别的 `Assign` 命中）。
⇒ 进 Job 的唯一门 = `StartInJob`。

| 谁被分配进去过 | 位置 | 性质 |
|---|---|---|
| `wisp slo` 的被测子进程（`wisp.exe slo -subject …`） | `cmd/wisp/slo_windows.go:498` | **唯一的产品侧分配点**，落在 SLO 测量夹具里 |
| 测试夹具：重执行的测试二进制自身（`-test.run=^TestHelperProcess$`） | `internal/proc/jobscope_windows_test.go:41`（经 `startHelper`） | 测试，非产品 |
| spike：把 **spike 进程自己**塞进 kill-on-close Job | `scripts/spike/common/winshell.go:490` | 夹具，不在构建产物 |
| **一轮 agent 任务起的子进程** | —— | **零枚**（本格不宣称"不可能"，只报"现存代码里没有调用者"） |

**"等价形状"这一问的答案是否。** 把进程塞进 Job 的等价写法是 `CREATE_SUSPENDED` 起、`AssignProcessToJobObject`、再 `ResumeThread`
（先挂起再归属，才没有归属前就跑起来的窗口）。`-E 'CREATE_SUSPENDED|ResumeThread|CreationFlags'` 那发普查
**零枚** `CREATE_SUSPENDED`、**零枚** `ResumeThread`，两枚 `CreationFlags` 命中都是 `CREATE_NEW_PROCESS_GROUP`（进程组，与 Job 无关）。

### 2.3　顺带量到的一处现存形状（不是本票的修法，是给 AC#2 判料）

`StartInJob` 的门是"**先 `cmd.Start()` 再 `Assign`**"（`jobscope_windows.go:114`→`:117`）：

```
$ sed -n '114,121p' internal/proc/jobscope_windows.go
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("proc: StartInJob: %w", err)
	}
	if err := j.Assign(cmd.Process); err != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		return nil, err
	}
```

⇒ 在 `Start` 与 `Assign` 之间，那枚子进程**已经在无 Job 状态下运行**；它若在这条窗口里自己再 fork，
生下来的孙进程就**永远不在**这枚 Job 里（Job 归属是"起我的时候我在哪，我就在哪"）。
`Assign` 失败时它会被 `Kill`+`Wait` 收掉，所以票 03 那句"防逃逸"就"**分配失败不留孤儿子进程**"这半是属实的；
就"**没有任何逃逸窗口**"这半，本格的读数是**它没覆盖这一层**。那两枚出处的原文（现跑）：

```
$ grep -n '防逃逸' docs/evidence/s0/03-adversarial-acceptance.md
97:  `StartInJob` = Start 后立即 Assign，**Assign 失败先 Kill+Wait 再返回错误**（防逃逸）；

$ grep -n 'escape-proof' .scratch/wisp/issues/03-skeleton-runtime-rules-done.md
60:- [2026-09-19T07:18:37Z] agent=T03-impl did=JobScope-C30-proc (OpenJobScope KILL_ON_JOB_CLOSE, Assign/StartInJob escape-proof, TreePrivateBytes+TreeProcessCount via pid-list+psapi PrivateUsage; tests: READY-handshaked child killed on job close, multi-child accounting, empty-job zero, failed-start leaves no child; race-clean) next=single-instance+env-fork
```

⚠ 本程**没有**为这条造用例（造它要真起进程并卡那条窗口），也没有把它算成本票的缺口——它记在 §7"没测什么"里，
并把话留在 AC#2 那族修法真被批准时一起裁。

### 2.4　正面控制：这条门今天真的通着电

普查只说"有没有人调用"，不说"调用有没有用"。所以本格真跑一次那三枚 Job 用例（宿主原生 Windows）：

```
$ go test -count=1 -run 'TestJobScope|TestStartInJob' -v ./internal/proc/
=== RUN   TestJobScopeKillsChildOnClose
--- PASS: TestJobScopeKillsChildOnClose (0.04s)
=== RUN   TestJobScopeTreeAccounting
--- PASS: TestJobScopeTreeAccounting (0.07s)
=== RUN   TestStartInJobRejectsBadCommand
--- PASS: TestStartInJobRejectsBadCommand (0.01s)
PASS
ok  	github.com/CarlosShao/wisp/internal/proc	0.162s
rc=0
```

⇒ 三枚全过、`rc=0`、无 `--- FAIL:`。整机 Job 的"关了就杀"这条**在真机上是活的**；
本票要找的"任务级级联"**没有对应的用例**，因为它没有对应的实现——这条判语在 AC#3 那格落字。

### 2.5　放水两问自答

- **断言方向动没动**：没动。本格跑的是**原有**的三枚 `TestJobScope*` 用例，一字未改（`git diff` 在 §5 零改动自证里出）。
- **helper 是不是原有的那枚**：是——`internal/proc/jobscope_windows_test.go:27 startHelper` 是票 03 交付的那枚，本程未改它、也未新造注入面。
- 补一条：`-run` 的正则是"这三枚"，不是"这一包"；本格的 `rc=0` 只代表这三枚，**不代表 `internal/proc` 整包**（整包读数在 §6 的门禁格里另出）。

### 2.6　本格没测什么

1. **没证明"归属窗口"真能被利用**（§2.3 那条）：没有用例，也没有测量，只是读码得到的形状。若有人要按它开票，得先有一把尺。
2. **没验 `wisp slo` 那枚被测子进程在真实运行里确实落在 Job 内**（`TreeProcessCount()>0`）——那要跑 `wisp slo`，会另起子进程并抢 CPU，且 `cmd/wisp` 是票 149 的写者。
3. **没在 Linux 侧看任何等价性**：`internal/proc` 的 Job 那半边整个是 `_windows.go`，非 Windows 侧今天**没有对应文件**（这条是 AC#1③ 与 AC#2③ 的交点，话留在 §4）。

---

## 3　AC#1③ —— 今天哪一条生产路径真的会 exec 出子进程（逐名列候选）

**判据**（票面 AC#1③ 原文）：*"今天哪一条生产路径真的会 exec 出子进程：逐枚列出候选（`fs.*` 之外还有哪些），
并注明 `shell.exec` / D46 命令插件今天注册了没有（`internal/tools/unwired.go` 那条
"no shell.exec tool is registered" 的注释是不是仍然属实）"*。

⚠ **票面这一格里程打偏了一处**（不改票面一字，另起登记）：那条注释**不在 `internal/tools/unwired.go`**，
该文件不存在；原文在 **`internal/config/unwired.go:63`**。见 §3.4。

### 3.1　生产码里 `os/exec` 的持有者（非测试文件，逐名）

```
$ git grep -l '"os/exec"' -- '*.go' | grep -v '_test.go'
cmd/balldebug/diff_windows.go
cmd/wisp/doctor.go
cmd/wisp/slo_windows.go
internal/llm/adaptertest/mockllm.go
internal/proc/jobscope_windows.go
scripts/spike/webview2-latency/main.go
tools/d22scan/gitignore.go

$ git grep -l '"os/exec"' -- '*.go' | wc -l
38
```

⇒ 全仓 38 个文件 import `os/exec`，**其中 31 枚是 `_test.go`**；剩下 7 枚按"是不是产品运行时、
是不是在某一轮任务里"逐枚拆：

| # | 位置 | exec 的是什么 | 在不在一轮 agent 任务的路径上 |
|---|---|---|---|
| 1 | `cmd/wisp/doctor.go:316`（`gccVersion`，被 `doctor.go:49` 调） | `$CC --version`（默认 `gcc`） | **否**：`wisp doctor` 自检子命令，一次跑完就退出，没有 `RunningTask`、没有取消语义 |
| 2 | `cmd/wisp/slo_windows.go:490`（`startSubject`） | 自身 `wisp.exe slo -subject …` | **否**：SLO 测量夹具；而且它是**唯一**把子进程塞进整机 Job 的产品调用点（§2.2 表第 1 行） |
| 3 | `cmd/balldebug/diff_windows.go:409` | 自身子进程（截图比对） | **否**：`cmd/balldebug` 是开发工具，不是交付产物 |
| 4 | `internal/llm/adaptertest/mockllm.go:65,72` | `go build` ＋ `mockllm.exe` | **否**：包名不带 `_test.go` 后缀，但**只被测试 import**（`git grep -n 'llm/adaptertest' -- '*.go' \| wc -l` ＝ **13**，再 `grep -c '_test.go'` ＝ **13**） |
| 5 | `internal/proc/jobscope_windows.go` | **不 exec**：只为 `StartInJob(cmd *exec.Cmd)` 这个参数类型 import | —— |
| 6 | `scripts/spike/webview2-latency/main.go:283,294` | spike 自身二进制 | **否**：spike，不在构建产物 |
| 7 | `tools/d22scan/gitignore.go:373` | `git` | **否**：CI 仪器（票面点名零字节的那棵），不在产品运行时 |

**产品运行时（`wisp` 交付物）在一轮任务里 exec 出的子进程数：0。**

### 3.2　这一问不能只用 import 图回答（本程差点用错的一把尺）

一开始本程想拿"`os/exec` 在不在 `internal/agent` 的依赖闭包里"当尺，量完发现**这把尺会给出反的结论**：

```
$ go list -deps ./internal/agent | grep -c '^os/exec$'
1

$ go list -deps -f '{{.ImportPath}}|{{join .Imports " "}}' ./internal/agent | awk -F'|' '$2 ~ /(^| )os\/exec( |$)/ {print $1}'
modernc.org/libc

$ go list -deps ./internal/agent | grep -c 'CarlosShao/wisp/internal/proc'
0
```

⇒ `os/exec` **确实在** `internal/agent` 的传递依赖里，但它是 `modernc.org/libc`（SQLite 驱动链进来的）import 的，
不是 agent 自己；而 `internal/agent` 的依赖闭包里**根本没有 `internal/proc`**（第三发 `rc=1`、计数 0）——
也就是说，**今天的 agent 包连那枚整机 Job Object 都拿不到**，更谈不上"每轮一枚"。
本格因此改以"调用点普查"（§3.1）为准，import 图只作对照。

### 3.3　生产任务路径实际能调到的工具名册（`fs.*` 之外还有没有）

`wisp run` 的组合根只注册这一族：

```
$ grep -n 'for _, e := range tools.BuiltinFSEntries' cmd/wisp/run.go
341:	for _, e := range tools.BuiltinFSEntries(tools.FSDeps{

$ sed -n '341,347p' cmd/wisp/run.go
	for _, e := range tools.BuiltinFSEntries(tools.FSDeps{
		Paths:         rt.paths,
		DeleteEnabled: cfg.FS.DeleteEnabled,
	}) {
		if err := reg.Register(e); err != nil {
			fmt.Fprintf(s.stderr, "wisp run: 工具注册失败（%s）：%v\n", e.Tool.Name(), err)
			return rt, 2

$ sed -n '323,329p' internal/tools/fs.go
func BuiltinFSEntries(d FSDeps) []Entry {
	out := []Entry{
		{Tool: fsRead{d: d}, Decl: FSReadDecl()},
		{Tool: fsList{d: d}, Decl: FSListDecl()},
	}
	return append(out, BuiltinFSWriteEntries(d)...)
}
```

名字逐枚（`git grep -n 'func.*Name() string' -- internal/tools/fs.go internal/tools/fs_write.go`）：
`fs.read` / `fs.list` / `fs.write` / `fs.trash` / `fs.move` / `fs.delete`（最后一枚仅 `[fs] delete_enabled=true` 时注册）。
这六枚里唯一与"外部世界"打交道的是 `fs.trash`，它走的是 **`shell32.dll!SHFileOperationW` 进程内 COM 调用**
（`internal/tools/recycle_windows.go:58`），**不产生子进程**。
⇒ `fs.*` 之外的候选：**零枚**。

### 3.4　`shell.exec` 与 D46 命令插件今天注册了没有

```
$ git grep -n '"shell.exec"' -- '*.go' | wc -l
6
$ git grep -n '"shell.exec"' -- '*.go' | grep -c '_test.go'
6

$ git grep -n 'ShellEnabled' -- '*.go' | grep -v '_test.go'
internal/config/manager.go:341:	if !old.ShellEnabled && new.ShellEnabled {
internal/config/manager.go:343:	} else if old.ShellEnabled && !new.ShellEnabled {
internal/config/schema.go:451:	ShellEnabled bool `toml:"shell_enabled" default:"false"`
internal/config/unwired.go:65:		fires:   func(c *Config) bool { return c.Risk.ShellEnabled },
```

- **`shell.exec` 注册了没有：没有。** `"shell.exec"` 这个字面串全仓 6 枚命中，**6 枚全在 `_test.go`**
  （`internal/risk/assessor_test.go`、`internal/risk/rules_test.go` 测 R6 规则；`internal/config/unwired_test.go` 测那条守卫）。
  产品侧 `risk/rules_shell.go` 那条 R6 规则**写好了但没有生产者**：`ShellEnabled` 的非测试命中只有
  schema 定义、加载期守卫、和"热加载时把它记成 loosen/tighten 的一行差异报告"（`manager.go:341`，不是消费者）。
- **那条注释还属不属于实：属实。** 原文在 **`internal/config/unwired.go:60-67`**：

```
$ sed -n '60,67p' internal/config/unwired.go
var unwiredKeys = []unwiredKey{
	{
		path:    "risk.shell_enabled",
		missing: "no shell.exec tool is registered in internal/tools, so nothing reads this flag",
		lands:   "it lands with the shell.exec tool itself (SPEC-07 sec 3, S3, default-disabled per D14)",
		fires:   func(c *Config) bool { return c.Risk.ShellEnabled },
	},
```

  ⚠ 票面把它的位置写成 `internal/tools/unwired.go`——**那枚文件不存在**（`git ls-files internal/tools/unwired.go` 空）。
  文字内容属实，路径不属实；本程不改票面一字，只在此登记。
- **D46 Tier-1 命令插件注册了没有：没有实现，只留了槽。**

```
$ grep -n 'var providerSlots' internal/tools/registry.go
52:var providerSlots = []SlotErr{

$ sed -n '52,56p' internal/tools/registry.go
var providerSlots = []SlotErr{
	{Kind: KindManifest, Ticket: "ticket 50 lands Tier-1 manifests"},
	{Kind: KindGoja, Ticket: "ticket 51 lands the Tier-2 goja runtime"},
	{Kind: KindMCP, Ticket: "REJECTED: D13/16.9#7, interface slot only"},
}

$ ls internal/plugin/
disposal.go
disposal_test.go
doc.go
```

  manifest 槽一旦被组合根问到就返回 `ErrSlotNotLanded`（`registry.go:29-31`），`internal/plugin/` 整包只有 C11 的
  disposal 与 doc；D46 那张"外部命令插件"表**没有一处实现**。

### 3.5　读数小结（AC#1③ 的一句话答案）

**今天没有任何一条生产路径会在"一轮任务"里 exec 出子进程。** 会 exec 的四枚产品侧位置
（doctor / slo 被测子进程 / balldebug / 测试夹具 mockllm）都不在 `RunningTask` 的生命周期里，
其中 slo 那枚已经归了整机 Job。任务侧之所以起不了子进程，是因为**任务能调的工具只有六枚 `fs.*`**，
而 `shell.exec`（SPEC-07 的 S3）与 D46 命令插件（票 50）都还没落地。

### 3.6　放水两问自答

- **断言方向动没动**：没动，本格零断言。
- **helper 是不是原有的那枚**：本格没有 helper；所有读数来自 `git grep` / `go list` / `sed -n` 三把现成的尺。
- 补一条自答：**本格没有把"grep 零命中"当成"不可能"**。§3.2 就是本程自己先拿错尺、再换对的记录，
  它同时给出一条对本仓有用的判语：**import 图证明不了"不会 spawn"，只有调用点能**。

### 3.7　本格没测什么

1. **没有跑一次真任务并数它的子进程数**（例如 `wisp run` 全程用 `Get-CimInstance Win32_Process` 盯子节点）。
   那要构 `cmd/wisp`（票 149 的写者正在那棵里改），本程不构、不跑。⇒ §3.5 那句"0"是**静态调用点**读数，
   不是运行时观测读数。谁要把"0"升格成运行时结论，得等 `cmd/wisp` 空出来。
2. **没有排第三方库里自发子进程的可能**。`os/exec` 经 `modernc.org/libc` 在闭包里（§3.2），
   本格只证明了"agent 自己的码不调用它"，没证明"闭包里没有任何一枚第三方码会调它"。
   （诚实补一句：真要排这一条，得审计 `modernc.org/libc` 与 sherpa 那族 cgo 绑定的实现，本程没做。）
3. **没有验 `wisp run` 之外的组合根**（常驻球 / 面板那条）注册了哪些工具——因为那条组合根本程没找到
   （`internal/panel`、`internal/speech` 至今只有 `doc.go`）。若它已在别的未跟踪分支里长出来了，本格的"六枚名册"要重算。

---

## 4　AC#1 总裁决 —— "已实现功能里的缺口"还是"给未来功能预留的洞"

**判据**（票面 AC#1 末段原文）：*"这一格允许得出的结论是'今天无入口踩得到'——那也是有效交付，但要**同时**答：
那它是'已实现功能里的缺口'还是'给未来功能预留的洞'？若是后者 ⇒ **必须**去 `SPEC-12 §5` 按五字段
登记为 RESERVED／DEFERRED，并把本票降级为 `ready-for-human`，**不许**为了结案去造一条假的现实危害。"*

### 4.1　裁决：**预留的洞**，不是已实现功能里的缺口

三发读数支撑（全部现量于 §1/§2/§3，逐条可复算）：

1. **任务侧起不了子进程**：一轮任务能调的工具只有六枚 `fs.*`（§3.3），唯一碰外部世界的 `fs.trash` 走进程内 COM（`SHFileOperationW`）。
   会 exec 的四枚产品侧位置（doctor / slo / balldebug / 测试夹具）**没有一枚在 `RunningTask` 的生命周期里**（§3.1）。
2. **任务侧连整机 Job 都拿不到**：`go list -deps ./internal/agent | grep -c 'internal/proc'` ＝ **0**（§3.2 第三发）。
   "每轮任务一枚自己的 Job"这个形状，**今天缺的不只是那枚 Job，还缺一条从组合根把 `internal/proc` 递进 `agent.Options` 的通路**。
3. **连"任务 scope"这个容器都还没有生产构造点**：SPEC-01 §6 写死了层级
   （`docs/specs/SPEC-01-architecture.md:196`：*"层级：会话 scope（C31）⊃ 任务 scope ⊃ 工具 scope ⊃ 插件 scope"*），
   C11 `plugin.DisposalScope` 也确实实现了（`internal/plugin/disposal.go:155 NewDisposalScope`），但——

```
$ git grep -n 'NewDisposalScope' -- '*.go'
internal/memory/retention_test.go:193:	scope := plugin.NewDisposalScope("test-retention", nil,
internal/plugin/disposal.go:128:// not usable; use NewDisposalScope.
internal/plugin/disposal.go:152:// NewDisposalScope creates a scope derived from parent (nil = detached). The
internal/plugin/disposal.go:155:func NewDisposalScope(name string, parent context.Context, opts ...Option) *DisposalScope {
internal/plugin/disposal_test.go:17:	s := NewDisposalScope(name, context.Background(), opts...)
internal/plugin/disposal_test.go:168:	s := NewDisposalScope("task-tail", context.Background(),
internal/risk/provenance_test.go:307:	scope := plugin.NewDisposalScope("session-1", context.Background())
```

⇒ 生产码里 `NewDisposalScope` 的调用者 **0 枚**（七枚命中：1 枚函数声明、2 枚注释、4 枚测试）。
"任务级归属"这件事在这仓里**目前只有契约文本，没有任何生产载体**。

### 4.2　为什么本程**不**把它做成代码（AC#2 那族的修法未触发的三条理由）

票面 AC#2 的引导句是条件句：*"**如果** AC#1 判定确有可走路径 ⇒ 修法只许落在'每轮任务一枚自己的 Job Object、取消时级联'这一族"*。
§4.1 的裁决把这个条件判**假**，于是本程按票面后半走"登记＋上报"，不写码。三条独立理由：

1. **没有生产者**：级联要杀的那些子进程今天不存在（§3.5）。实现了也没人调用它，
   正好落成本仓抓过的那族"一步存在却从不产出结论"（票面 AC#2③ 的括号里点名的就是这族）。
2. **AC#3 的牙造不出来**（不是不想造，是没有可摘的东西）：AC#3 要求"必须有一枚用例在'级联'被摘掉时转红"。
   今天没有级联，"摘掉级联"这发变异**在物理上无法落笔**。若为了让 AC#3 响而先造一级联再摘它，
   那就是票面禁止的"造一条假的现实危害来结案"。
3. **落点与在飞的票 149 硬撞**：任何"每轮一枚 Job"的接线都要经过 `cmd/wisp` 的组合根——现量：
   `agent.New` 在产品码里**只有一枚调用者**，`agent.Options` 也只在那里填一次。

```
$ git grep -n 'agent.New(' -- '*.go' | grep -v '_test.go'
cmd/wisp/run.go:547:	loop, err := agent.New(agent.Options{
```

   而 `cmd/wisp/**` 此刻是票 149 的写者。这不是判据不足，是**同一时刻同一文件两枚写者**——本仓的规矩是切开落点，不是叠上去。

⚠ **本程没有做的两件事，以及为什么没做**（这两件都在票面 AC#1 的字面要求里，但都落在本程被冻的面上）：

| 票面要求 | 本程动作 | 理由 |
|---|---|---|
| "去 `SPEC-12 §5` 按五字段登记" | **未落笔**，只在 §4.3 备好可直接粘贴的那一行 | `docs/specs/**` 是简报点名的零字节面；且 `AGENTS.md §1.1` 把"登记表与代码标记 1:1 双向"列成受审项，本程单侧改表会造出一枚**没有代码标记对应、也没有验收表**的登记 |
| "把本票降级为 `ready-for-human`" | **未改 Status 一行**，只在 §9 与票面进度日志里把判据与请求报回编排者 | 简报第 0 条："票面框由编排者按非实现者验收表来定"；Status 属同一族的票面框 |

### 4.3　可直接粘贴的登记行（五字段齐全，类型＝DEFERRED）

`SPEC-12 §5` 的类型语义（逐字，`docs/specs/SPEC-12-roadmap-governance.md:61`）：
**DEFERRED** ＝ 有计划不做；**RESERVED** ＝ 只留接口位、无实现计划、不得误读为待办。

⇒ 按这语义，本条该记 **DEFERRED**，不是 RESERVED：这件事**有计划**（`shell.exec` 排在 SPEC-07 的 S3、
D46 命令插件是**已采纳**项并有票 50），"只留接口位、无实现计划"那句套不上。票面写的是"RESERVED／DEFERRED"两可，
本程按表内语义取 DEFERRED，理由如上，**若编排者要取 RESERVED，改的是这一格类型列，不是内容**。

表头是六列（`类型｜项｜为什么现在不做（依据）｜完成判据｜前置｜当前残缺表现`），本行**逐列填满**，
"五字段缺一不可"在任何一种数法下都成立：

```markdown
| DEFERRED | 任务级子进程归属（每轮任务一枚 Job Object、取消时级联） | 票 138 AC#1 现量：今天没有任何一条生产路径在一轮任务里 exec 子进程——任务可调的工具只有六枚 `fs.*`（`fs.trash` 走进程内 `SHFileOperationW`）、`shell.exec` 未注册（`internal/config/unwired.go:63` 守卫仍在拦）、D46 Tier-1 manifest 槽返回 `ErrSlotNotLanded`、`internal/agent` 依赖闭包里没有 `internal/proc`、连 C11"任务 scope"都没有生产构造点 ⇒ 属预留洞不是实装缺口。另：AC#1 硬约束不许顺手改 `cancelled` 文案（那半是 D37 17 类表＝契约级） | 一轮任务自己起的每枚子进程，在 `RunningTask.Cancel()` 之后随该轮那枚 Job 一起没了；且有一枚用例在"级联"被摘掉时转红（票 138 AC#3 那族的牙） | `shell.exec`（SPEC-07 §3，S3）或 D46 命令插件（票 50）任一先落地，并给 `agent.Options` 一条从组合根拿 `internal/proc` 的通路 | 按"停"只掐 context，那一轮若起了子进程没人负责关；**今天无入口踩得到**，落地当日才成真伤 |
```

**并且**：这条登记一旦落下，票 138 本身的结案形状就是"AC#1 成立、AC#2/AC#3 按条件句未触发、票转 `ready-for-human`"。
本程把这行原样交给编排者，**不代他按批准键**。

### 4.4　放水两问自答

- **断言方向动没动**：没动。本格是"判有没有"，不是"判好坏"；三发支撑都是 `git grep` / `go list` 的原始输出。
- **helper 是不是原有的那枚**：本格无 helper、无测试、无码改。
- 反向自答（这格最容易被判"放水"的地方）：**本格有没有为了不做工而把'缺口'说成'预留'？**
  判据是外部的、可复算的：`internal/agent` 闭包里 `internal/proc`＝0、生产 `NewDisposalScope`＝0、
  `"shell.exec"` 生产注册＝0。三发都是"数出来的"，不是"选出来的"。
  若这三发里任何一枚被复算推翻，本格的裁决就要翻（复算命令逐字贴在 §3.2 与 §4.1③）。

### 4.5　本格没测什么

1. **没排除"别的分支上已经有 shell.exec"**：本格读的是 `dev`@`64858d6` 的跟踪文件。
   若某条未合并分支落了 `shell.exec`，"无入口"那条当场失效——这也是本程把裁决交回人拍板而不是自己结案的原因之一。
2. **没量"未来落地时该挂在哪"**：SPEC-01 §6 那句层级文本与 C11 的实现都在盘上，
   但"每轮任务的 Job 该由 `DisposalScope.Defer` 收、还是由 `RunningTask` 自己收"是**设计选择**，
   票面 AC#2 只给了族名没给归属 ⇒ 本程不替人拍这一枚。
3. **没验 `cancelled` 文案那半与本格有没有牵连**：本程一字未动 `D37` 那 17 类（见 §5 零改动自证），
   所以那半截契约变更**没有被触发**；但本程也没去证明"它不可能被牵连"——那是验收表那侧的活。


