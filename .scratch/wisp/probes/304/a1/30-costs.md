# 304-a1 — 表③：三形代价表＋必答两问

射程＝`20-roster.md` §1 那 **18 枚"父不等它"** 的 spawn 点（尤其票面 `:15` 点名的两枚 `mockllm` 源＝表①的 #1 与 N2）。本腿⛔ 跑 Go／⛔ 跑门禁 ⇒ **凡"CI 色／rc 色"一律标〔读码〕或〔需真机〕**，并落进 `90-final.md` 的 `Q-304-a1-*`。

## 1. 三形代价（逐形：改哪几处／现成仪器／blast radius／平台覆盖／风险）

| 形 | 要改的落点（现量行号） | 现成仪器与先例（量来的） | 代价／风险 | 平台覆盖 |
|---|---|---|---|---|
| **甲＝Windows Job Object＋`JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`**（测试二进制自己建 job，把 mockllm 塞进去） | 两枚 spawn：`internal/llm/openaichat/mockllm_integ_test.go:73`/`Start :80` ＋ `internal/llm/adaptertest/mockllm.go:72`/`Start :79`（★票面 `:15` 硬约束＝**两枚必同批**）；job 生命周期各一枚 `t.Cleanup(job.Close)` | **本仓已有整枚仪器**：`internal/proc/jobscope_windows.go` —`OpenJobScope :71-86`（`CreateJobObject(nil,nil)` `:72` ＋ `KILL_ON_JOB_CLOSE` `:77`）、`StartInJob :110-123`（assign 失败即当场 Kill＋Wait，注释逐字 "a child can never escape the Job by accident" `:107-109`）、`Close :170-181`。**测试侧先例**＝`internal/proc/jobscope_windows_test.go:41`（真把子代塞 job）＋`t.Cleanup(job.Close)` `:111`/`:161`/`:204`＋**判据 `waitGone :92-104`（逐字 "helper survived for %v after Job close"）**。**产码同形**＝`cmd/wisp/slo_windows.go:500 rt.Job.StartInJob`。**import 先例**＝尺 `grep -rn "wisp/internal/proc" internal cmd tools` 现量 **23 枚文件**，其中 `internal/llm/tempdir_resolved_124_test.go:15`＝同族已在 import | ⛔ 新 import 边：`adaptertest` 与 `openaichat` 现量⛔ import `internal/proc`（尺＝`grep -rn "wisp/internal/proc" internal/llm/` 只命中上面那一枚 `:15`）；job 归测试二进制所有 ⇒ **一次跑内⛔ 复用别家的 job**；job 只端**显式 assign 的**成员 ⇒ `go build` 那两枚（`openaichat:65`/`adaptertest:65`）⛔ 进 job（它们本就是同步等的）；句柄继承会把"最后一枚句柄关"往后拖 ⇒ 现量 `CreateJobObject(nil, nil)`＝属性参数为 nil ⇒ 不可继承（**读码级**，那半句"Go runtime ⛔ 把它传给子代"需真机） | **⛔ windows 专属**：`jobscope_windows.go:1` ＝`//go:build windows` ⇒ ubuntu 那一支（`ci.yml:420`/`:487 --scope=core`，而 `internal/llm/…` 在 core pin `portable-tests.sh:172-176`）**⛔ 被甲形保护**，`mockllm`（无 `.exe`）在 ubuntu 上同样能漏 |
| **乙＝子进程自己盯父 PID**（`mockllm` 加 `--parent-pid`） | ⛔ 一处变三处：① `tools/mockllm/main.go`（flag 名册现量 **4 枚**＝`:43 addr`／`:44 golden-dir`／`:45 models`／`:46 print-addr` ⇒ 新增第 5 枚＋watch 循环＋退出路径），②③ 两枚 spawn 点各多传一个参数（票面 `:15` 同批那条照旧生效） | **⛔ 任何现成仪器**：尺＝`grep -rEn "parent-pid\|ParentPid\|ParentPID\|GetParentProcess\|PPID" --include=*.go internal cmd tools third_party` 现量 **0 命中**；mockllm 今天**没有自退路**＝`main.go:49-66` 现读＝`listen → print MOCKLLM_ADDR → Serve(ln)`，`control.go` 端点现量 **8 枚** handler（`grep -c "func (s \*Server) handle"`＝8：fail-next/latency/truncate/capability/reset/last-request/state/models），**⛔ shutdown 端点** | ⚠ **跨模块**：`tools/mockllm/go.mod` 现量＝`module github.com/CarlosShao/wisp/tools/mockllm` ＋ `go 1.27`，**⛔ require／⛔ replace**，根目录 **⛔ `go.work`**（尺＝`ls go.work*`＝No such file）⇒ 它**用⛔ `internal/observe.NewTimeout`**（`internal/observe/clock.go:32-49`）＝仓里那枚单调钟 helper；`main.go:1-5` 注释逐字 "an **stdlib-only** module" ⇒ 加 require/replace＝**改模块图**（本票 `AC#1` 硬约束⑤之外的一条新风险）。stdlib 内可行的读码形＝`os.FindProcess(ppid)` 报错即父不在（⛔ 需 x/sys）；⚠ 该形对"pid 复用"的判据需真机（本腿在 `10-residual.md` §1 已量到一枚**真 pid 复用现例**：`11724` 现为 `RuntimeBroker.exe`） | **⛔ 平台绑定**（linux 那一支也能盯父）＝比甲多覆盖 ubuntu CI 那一支——但 linux 侧"怎么读父在不在"是**另一份实现**（读码级） |
| **丙＝不动＋登记具名已知残留＋把 `mockllm.exe` 写进起手闸名册** | ⛔ 动产品码；只动**模板/名册**文本（`AC#4` 那一格） | 编队实践**今天已经在数它**（现量，尺＝`grep -n "IMAGENAME eq" docs/reports/HANDOVER.md`）：`HANDOVER.md:2431`、`:2455`、`:2466`、`:2474`、`:2490` 五处逐字列 `mockllm.exe`＝**0 枚** ⇒ ⚠ **票面 `:5` ① 那句"`mockllm.exe` 不在这两名册里"对模板成立、对最新实践已过期**（现量五行＝证据） | 名册**模板**那几处确实还只列两名（具名行，尺＝`grep -rn "IMAGENAME eq" .scratch/wisp/probes/orch/`）：`2026-10-09-live-window-runsheet.md:45`（只 balldebug）、`296-ac3-runsheet.md:7`（wisp＋balldebug）、`a793-append.md:7`（三枚含 `wisp-296.exe`，⛔ mockllm）、`handover-40as.md:4`／`handover-40at.md:4`（wisp／balldebug）。★**CI 那侧⛔ 任何进程计数**＝尺 `grep -cE "tasklist\|IMAGENAME\|Get-CimInstance" .github/workflows/ci.yml` 现量 **0** ⇒ 托管 runner 上漏了⛔ 有人数、也⛔ 有人清 | 代价＝票面 `:45` 自己写的那句："机主台面上会继续累积挂着⛔ 睡的假大脑，而我们的起手闸门一直看不见它"；本腿给它补两个数：**每台机器一枚假大脑＝一枚 Listen 回环套接字**（票面 `:10` 那两枚各持 `127.0.0.1:54887`/`58918`）＋**同症状目录现量 21 枚**（`10-residual.md` §3）。⚠ 丙**⛔ 解决累积**，只把它变成可见的 |

## 2. 必答 (I)：甲形在**托管 runner** 上会⛔ 把 CI 的别的进程一起端掉？（job 的射程边界）

结论形状＝**"⛔ 会"这一半停在读码级，"会不会"那一半的最后一层需真机/需 CI 现量**。分层（每层标级）：

1. **〔读码〕射程＝显式成员**：本仓的 job 只有两个入口——`Assign(p *os.Process)`（`jobscope_windows.go:89`）与 `StartInJob(cmd)`（`:110`）；**没有**"把整机/整棵树塞进 job"那种调用。尺＝`grep -rEn "\bAssign\(|StartInJob\(" --include=*.go internal cmd tools` 现量 **7 行**（件＝`logs/grep-m9-job-entrypoints.txt`）＝**2 枚方法定义**（`:89`/`:110`）＋**1 枚方法内部自调**（`:117`）＋**3 枚真调用点**（`internal/proc/jobscope_windows_test.go:41`、`:207`、`cmd/wisp/slo_windows.go:500`；⚠ `jobscope_windows_test.go:42` 那行是 `t.Fatalf` 里的字样，⛔ 调用点）。⇒ 未被 assign 的进程（runner 代理、别的 step、机主那 18 枚 `msedgewebview2`）**⛔ 在 job 里**，`KILL_ON_JOB_CLOSE` 碰⛔ 到它们。
2. **〔读码〕job 关闭的时机**：`Close` 只 `CloseHandle`（`:178`），而 `OpenJobScope` 用 `CreateJobObject(nil, nil)`（`:72`）＝⛔ 具名、属性 nil ⇒ 句柄⛔ 可继承。外力 `TerminateProcess` 杀测试二进制时 OS 销毁其句柄表 ⇒ job 句柄关 ⇒ **OS 端子内成员**。这条正是本票要的 (c) 形机制，但它**⛔ 等于实测**＝`Q-304-a1-A`（派单 §0.3 与票面 `:24` 都写着"⛔ 读文档当读数"）。
3. **〔读码〕嵌套 job**：若 runner 已把测试二进制放进它自己的 job——Win8+ 允许子 job，关闭子 job ⛔ 动父 job 成员；Win8 之前那枚 `AssignProcessToJobObject` 直接返回失败，而本仓 helper 在失败时**先 Kill 那一枚子代再报错**（`:117-121`）⇒ 两种历史形状下都⛔ 误端别人，最坏是那一枚子代没跑成＋测试红。**托管 runner 到底是哪一种**＝需真机/需 CI 现量。
4. **〔需真机／需 CI 现量〕本仓今天已经在这两台机器上跑 job 用例**：`TestJobScopeKillsChildOnClose` 一类（`jobscope_windows_test.go:107`/`:157` 的 `OpenJobScope`）在 **hosted `windows-latest`** 的射程里——尺＝`scripts/portable-tests.sh:249`（`windows)` 分支）＋`:250-254` 的 scope 数组**首项就写着 `./internal/proc/`（`:251`）**＋`win_pin:200`＋`ci.yml:534 runs-on: windows-latest`＋`ci.yml:797 bash scripts/portable-tests.sh --scope=windows`。⇒ **如果甲形会端掉 runner 上的别的进程，本仓早就在托管色上见过它**。这一问本腿⛔ 能交"CI 色"读数（⛔ 跑 Go、⛔ push、⛔ 读 run 日志）＝**`Q-304-a1-B`**。
5. **〔需机主〕自托管那一支**：`ci.yml:880 runs-on: [self-hosted, wisp-slo]`＝**机主自己那台笔记本**（台账同口径现量记在 `docs/reports/pending-and-issues.md:946`、`:3429`，本腿⛔ 动台账）。甲形在那台机器上的 blast radius 落在**真台面**上：本腿现量那 18 枚 `msedgewebview2` 属 `SearchHost`/`wetype_update`/`clipsync-desktop` 三棵别家树、父链全活（`10-residual.md` §2）⇒ **⛔ assign 就⛔ 碰**（读码级）。⚠ 另外那台机器上**产码本来就在用 job**（`cmd/wisp/slo_windows.go:500`）＝甲形⛔ 引入一枚新边界，而是**把已经在用的那枚边界搬到测试侧**。

## 3. 必答 (II)：乙形"轮询父 PID"撞⛔ 撞 D22 那条"⛔ 用墙钟时间差实现超时"

★按派单令：**先看仪器实际扫什么**，再看规格怎么写，两栏分开。

**仪器那一栏（现读字面，`tools/d22scan/main.go`）**

- 正则字面：`:166 wallclockRe = regexp.MustCompile(\`\.Sub\(time\.Now\(\)\)\`)`——**只认 `.Sub(time.Now())` 这一枚拼写**（与派单 §2(C)(II) 记忆里那句一致，本腿逐字核过）。
- 第二支：`:167 unixTimeRe = time\.Now\(\)\.(Unix|UnixNano|UnixMilli)\(` ＋ `:168 timeoutWordRe = (?i)timeout|deadline|expire|\bttl\b|budget|until`，**同行两者都命中**才报（应用点 `:789-791`）。
- 应用范围（`:781-785`，逐字读）：`for i, line := range lines { code := strings.TrimSpace(line); if code == "" \|\| strings.HasPrefix(code, "//") { continue } … }` ⇒ **注释与空行豁免、代码⛔ 豁免**。
- 逐形判断（尺＝上面两条正则字面，结论＝词面级，⛔ 跑过门禁）：
  | 候选写法 | 在⛔ 仪器射程 | 仓内今天的同形先例（现量行号） |
  |---|---|---|
  | `for { if !parentAlive { os.Exit(0) }; time.Sleep(500*time.Millisecond) }` | **⛔ 命中**（⛔ 一枚 `.Sub(time.Now())`，也⛔ 有 Unix＋timeout 字样同行） | `internal/memory/concurrent_test.go:256-260`（`for`＋`time.Sleep(10ms)` 轮询） |
  | `case <-time.After(interval):` | **⛔ 命中** | `internal/llm/openaichat/mockllm_integ_test.go:102`、`internal/proc/jobscope_windows_test.go:92-104`（`waitGone`） |
  | `end := time.Now().Add(d)` ＋ `time.Now().After(end)` | **⛔ 命中** | `internal/tools/fs_staging_windows_test.go:391-402` |
  | `observe.NewTimeout(...)`＋`deadline.Expired()` | **⛔ 命中**（其内部 `internal/observe/clock.go:40` 是 `time.Since(t.start)`＝单调钟拼写，⛔ 匹配 `\.Sub\(time\.Now\(\)\)`；`:33` 的 `start: time.Now()` 亦⛔ 匹配） | `internal/memory/concurrent_test.go:255`、`internal/llm/adaptertest/mockllm.go:105-107` |
  | `time.Now().Sub(start) > timeout`／`end.Sub(time.Now())` | **★命中第一支** ⇒ 报 `wallclock-timeout` | —（本腿⛔ 写候选码，只报形状） |
  | `time.Now().Unix()` 与 `deadline`/`until` 字样同行 | **★命中第二支** | — |

  ⇒ **形状结论**：乙形**只要⛔ 把"父在不在"的存在性轮询写成上面前三行那一类，就落在仪器射程之外**；把它写成"墙钟差比大小"才撞。⚠ **但 mockllm 够⛔ 着 `observe`**（独立 module，`tools/mockllm/go.mod` 现量⛔ require/⛔ replace）⇒ 它只剩手写 `time.Until`／`time.After`／`time.Since` 这三枚单调安全拼写可用，或改模块图（＝甲/乙之外的一笔代价）。
  ⚠ **本腿⛔ 跑 `sh scripts/d22scan.sh`**（派单 §0.3 硬边界）⇒ "以上判断在真仪器下 rc=0"是**欠读数**＝`Q-304-a1-C`。

**规格那一栏（另写，⛔ 与仪器混为一谈）**

- `AGENTS.md §1.2` 那条写作"**用'墙钟时间差'实现超时**"（该文件自述它是薄索引、⛔ 权威来源；权威＝`PLAN.md` 的 D22 正文与 `tools/d22scan`）。**类别表述比正则宽**：轮询"父进程还在不在"本身⛔ 是超时（存在性判断，⛔ 时长判断）；但乙形若含"**等超过 X 还没退就放弃**"那一半，那一半**落进类别**，即便它⛔ 落进正则。
- 这条解释**⛔ 由本腿裁**：票面 `:25` 写的是"先答这两问再摆给我裁"＝编排者的动作。本腿交的是**两栏射程图**＋一句具名待定案＝`Q-304-a1-D`（规格侧"存在性轮询算⛔ 算超时实现"未定案，⛔ 权威文本里没有这一条的判例）。
