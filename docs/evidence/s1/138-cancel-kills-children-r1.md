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
2. **没验 `Boot` 之外有没有别的路径能拿到第二枚 Job Object**——例如"某人自己 `windows.CreateJobObject` 再塞进 `JobScope{handle:…}`"。本格的 `CreateJobObject` 普查覆盖了这个形状，但 `x/sys/windows` 的 `NewLazySystemDLL` 动态取符号这类**语法上抓不到**的绕法没排（仓内 `modpsapi` 就是这种形状，所以它不是假设）。
