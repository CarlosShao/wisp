# `304-r1` 10 — 落地形＝甲（Job Object，`KILL_ON_JOB_CLOSE`）接到两枚 mockllm spawn 点

改了什么（`git diff --stat b9f37892..worktree -- internal/llm/` 现量＝**4 files changed, 93 insertions(+), 2 deletions(-)**，尺＝`git diff --stat`，射程＝`internal/llm/`，比的是我这枚锚 commit 的树 ⇄ 工作树）：

| 文件 | 新/改 | 逐枚行（现量 `grep -n`） | 内容 |
|---|---|---|---|
| `internal/llm/adaptertest/childjob_windows.go` | 新 | `:1 //go:build windows`、`:31 func StartChildGuarded`、`:33 proc.OpenJobScope()`、`:39 job.StartInJob(cmd)`、`:40 _ = job.Close()`（失败支）、`:48-50 t.Cleanup(job.Close)` | windows 那一支：开一枚 C30 `proc.JobScope`（`KILL_ON_JOB_CLOSE`），**只把这一枚子进程**用 `StartInJob` 塞进去，成员一存在就挂 `t.Cleanup(job.Close)` |
| `internal/llm/adaptertest/childjob_other.go` | 新 | `:1 //go:build !windows`、`:21 func StartChildGuarded`、`:23 return cmd.Start()` | **非 windows 那一支＝明确 no-op**（见下"跨平台那支长什么样"） |
| `internal/llm/adaptertest/mockllm.go` | 改 | `:85 if err := StartChildGuarded(t, proc.cmd); err != nil {`（原 `:79 proc.cmd.Start()`），注释在 `:79-84` | spawn 点 **#1＝票面 `:15` 的 helper**（`adaptertest/mockllm.go:72` 起、`:65` 建）；回收仍留在同件原有 `t.Cleanup`（Kill＋Wait＋`spawnRoot.Wait`），⛔ 动其语义 |
| `internal/llm/openaichat/mockllm_integ_test.go` | 改 | `:87 if err := adaptertest.StartChildGuarded(t, proc.cmd); err != nil {`（原 `:80 proc.cmd.Start()`），注释在 `:81-86`；新增 import `adaptertest` | spawn 点 **#2**（票面 `:15` 的第二枚，原有 `t.Cleanup` 在 `:113`，⛔ 动） |

**两枚同批**＝票面 `:15`／派单 §1 的硬条件：只改一枚就留下第二个泄漏源。两枚现在走**同一段代码**（`StartChildGuarded`），⛔ 两份各写各的 job。

## 用的工具＝在库的那枚，⛔ 新造封装
- `internal/proc/jobscope_windows.go`：`OpenJobScope`／`StartInJob`／`Close`（派单 §1 给的行号 `:71/:110/:170` 是**快照**；我这轮现量的定义行位置⛔ 重述，本件只声明"调的就是这三枚导出符号"，产码同形先例＝`cmd/wisp/slo_windows.go` 的 `rt.Job.StartInJob(cmd)`）。
- 没有新 job 封装、⛔ 抄 `CreateJobObject`：`childjob_windows.go` 里唯一与 job 有关的动作就是 `OpenJobScope`/`StartInJob`/`Close` 三枚调用。
- `StartInJob` 自带"assign 失败就当场 Kill＋Wait 那一枚子代"的语义（`jobscope_windows.go:117-121`，⛔ 我重写），所以子代⛔ 可能"因为 assign 失败而逃出去"。

## 跨平台那支长什么样（边界①，逐字对着派单 §2.5 写）
`//go:build windows` 那支走 `StartInJob`；`//go:build !windows` 那支＝`return cmd.Start()`，**文件头注释里具名写着**：这不是 guard、⛔ 把它读成受保护——`internal/proc.JobScope` 在 windows 构建约束后面，非 windows 无句柄可关，所以 **ubuntu/linux 那一支今天仍会漏**（外力杀／timeout panic 之后子进程照旧活着），这枚修法⛔ 覆盖它，实测欠账＝`Q-304-r1-1`。⛔ 任何"看起来两边都管住了"的写法：两支的函数名相同、语义⛔ 相同，且差异写在**代码注释**里而不只写在我的件里。

## 边界②（⛔ 整棵树、⛔ 顺手收别人）
- 每枚 job 只有**一个显式成员**＝`StartInJob` 起的那一枚 `mockllm.exe`；⛔ `Assign` 任何别的进程，⛔ 把 `go build` 之类塞进同一枚 job（`go build` 那枚子进程的残留见 `Q-304-r1-2`）。
- ⛔ 新增任何"扫机器上的 mockllm 然后停掉"的机制。停进程只发生在两处：①`proc.JobScope` 的 OS 语义（关句柄只端**该 job 的成员**），②我这轮测量用完之后的处置＝复用编排者的 `mockllm-procs.ps1 -CleanIds`，它内建三条守卫（名字必须是 `mockllm.exe` ＋ 命令行含本仓 golden 目录 ＋ 父 PID 现量不存在），任一条不满足就打 `skip pid=… reason=…` 而⛔ 停（逐字在 `logs/paired-*.txt` 的 `clean:` 行）。

## ⛔ 动过的东西（越界清单，供 `304-v1` 反查）
SLO 阈值／golden（`internal/llm/testdata/golden/**` 零字节）／`internal/observe/**`（`mockllmRegistry`、`observe.NewTimeout`、D38 名目全部原样）／`tools/mockllm/**`（乙形落选，⛔ 加 `--parent-pid`）／`tools/d22scan/**`／`internal/proc/**`（只调用，⛔ 改一字）／三枚冻结件／`ci.yml`·`.github/**`／`docs/**`／票面任何人写的原句／别人的 `.go`。断言⛔ 放宽、⛔ 造 `--- SKIP` 换绿（本轮读数用的用例⛔ 改过判定条件，只改了"怎么起子进程"）。
- ⚠ 一处**必然的副作用**要具名：`internal/llm/adaptertest`（非测试件）新增一条到 `internal/proc` 的 import 边＝`A850` §3 甲形代价表里写着的那一枚；`internal/llm` 下原有 `internal/proc` 入边只有 `internal/llm/tempdir_resolved_124_test.go:15`（〔读码级〕，抄自 `304-a1` 表③，我⛔ 另跑尺）。
- ⚠ 反向窗口要具名：为了交 `AC#2` 的反形，`internal/llm/adaptertest/mockllm.go` 与 `internal/llm/openaichat/mockllm_integ_test.go` 在 **17:50–18:05** 之间被临时改回"无 job"形状（⛔ `stash`／⛔ `reset`，用 `git show b9f37892:<path>` 重定向到我自己的两枚件），随后逐字节还原（`cmp`＋`sha256sum` 双证在 `20-paired-readings.md` ＋ `logs/restore-verify.txt`）。这段窗口里在飞的其他腿看到的是**无 job 的那两枚文件**。
