# 票 164 · `task.output` 落点普查＋AC#2–#5 判据表（只读设计核 164-c1）

- 派单＝`.scratch/wisp/dispatches/2026-09-27-185x-readonly-164-c1-task-output-landing-survey-and-criteria.md`
- 票面＝`.scratch/wisp/issues/164-background-jobs-cannot-be-read-fix-the-roster-then-add-the-output-leg.md`（含 09-27 18:4x 追加的 Status 更正）
- 性质＝**只读**：本程**零字节改动**任何既有文件，**未写一行码**，**唯一写入＝本文件（新建）**。
- 证据分档标记：〔现跑〕＝本程本轮在同一棵树上跑过的命令；〔第二见证〕＝引别处的读数＋它的锚点；〔仅自述〕／〔推断〕＝没有读数支撑。
- ⚠ 本程**没有跑** `.scratch/wisp/probes/154/gate-clauses.sh` 与 `.scratch/wisp/probes/161/r6/flip-declaration.sh`（此刻有验收程在拿它们做名册差集）；关于它们的一切都是**只读文件**得来的，逐条标了行号。

---

## 1. step-0 三件〔现跑〕

```
date '+%Y-%m-%d %H:%M %z'                       -> 2026-09-27 18:32 +0800
git rev-parse --abbrev-ref HEAD                 -> dev
git rev-parse HEAD                              -> a651fe80464d17c5a8134af69bd4f2446bde53cb
git status --porcelain                          -> 见下（工作树里躺着别人的东西，一枚未 add）
```

现场（`git status --porcelain`，摘录）：` M .gitignore`、` M .scratch/wisp/probes/152/my152.py`、` M .scratch/wisp/probes/161/r6/logs/flip-*.txt`（9 枚）、` D design/**`（owner 未提交的删除）、` M design/doubao/*`。⇒ 本程**一枚都没碰、一枚都没 add**。

---

## 2. 本程**没**查什么（照实记，别当查过）

1. **没跑**票 154／161 那两把尺（明令禁跑，见文件头）。§5 那一节的每一句都是"读脚本 + 引行号"，不是读数。
2. **没跑** `go test ./...`、`go vet`、`tools/d22scan`：全篇没有任何"绿／红"的测试读数。
3. **没查**注册表是否用别的拼法（变量拼接、常量表、`"task." + name`）注册这三枚 ⇒ 这一条与票面 AC#1 裁定段 `:24` 自己承认的洞是**同一个洞**，本程**没有把它补上**：我只跑了字面量普查（§3.2 问② 那两条命令）。
4. **没做字节级核对** `PLAN.md:431` 那段路径字面量（见 §3.1 末的坑）。
5. **没查** `agent.Config.ArtifactsDir` 在组合根里由谁填（命令在 §3.4 问④ 的待办里给），也**没查** `fs.allowed_dirs` 默认是否包含 artifacts 目录。
6. **没核** SPEC-12 §5 那"五字段"的**权威命名**是否与 `PLAN.md:1515` 表头逐字同源（我只抽了表头与那一行，见 §3.1）。
7. **没动** `internal/tools/**` 一字节（派单硬规矩）。

---

## 3. 四问逐答

### 3.1 问① —— 冻结文本到底要求什么

**`docs/PLAN.md:2564` 整行**〔现跑：`git show HEAD:docs/PLAN.md | awk 'NR>=2555 && NR<=2570'`〕逐字：

> `| **`task.output`** | **读一个后台任务吐了什么**（含"太长怎么续读"） | **L0** | — | **S7** | 票 164（2026-09-27 owner 批准新增）的输出腿。截断按 **D15「单个工具结果」那一行的既有规矩**（头 500 token＋尾 200 token＋总长＋文件路径），**只截不指＝不合格**（必须给可续读的路径） |

**`docs/specs/SPEC-07-tools-and-plugins.md:71` 整行**〔现跑：同一手法，`NR>=60 && NR<=80`〕逐字：

> `| **`task.output`** | 读一个后台任务吐了什么（含续读） | **L0** | — | **S7** | 票 164 新增，镜像 `PLAN.md` D34；截断按 D15 规矩（头 500＋尾 200＋总长＋路径），**只截不指＝不合格** |

逐字段：

| 字段 | 读数 | 出处 |
|---|---|---|
| 档位 | **L0**（派单那句"是 L0 不是 L2"**核对为真**） | `PLAN.md:2564`、`SPEC-07:71` |
| 上游／capability 归属 | **—（空）** | 同一行第 4 列；`task.list`/`task.cancel`（`:2563`）第 4 列同样是 `—` |
| 切片 | **S7** | 两行一致 |
| 续读机制 | **文本只规定到"给一个可续读的路径"这一层**，规定的是**截断形状**（头 500 token＋尾 200 token＋总长＋路径），**没有**规定游标／偏移／分页／`offset` 参数／"读第 N 段"这类续读 API | `PLAN.md:2564` 括注＋`PLAN.md:431`（D15 那行原文，见下） |

D15 被指向的那一行 `docs/PLAN.md:431`〔现跑同上〕：

> `| **③ 单个工具结果** | **> 4000 token**（约 6000 汉字 / 16KB 文本） | 全文落 … 上下文里只留 **头 500 token + 尾 200 token + 总长度 + 文件路径**；模型可用 `fs.read` 按需再读（复用 D10 机制） |`

紧邻 `:432` 还有硬上限行：`> 1MB 原始输出 ⇒ 直接截断到 1MB 再走上一行，并标 truncated=true`。

⚠ **一处引文坑（上报，不修）**：awk 打印 `:431` 时那段路径显示成 `artifacts<TAB>ool-output-<id>.txt`，而 `internal/agent/spill.go:20` 的注释写的是 `<artifacts>\tool-output-<encoded id>.txt` ⇒ 怀疑源文件里是**反斜杠**、被打印/管道层当 `\t` 吃掉了。**本程没做字节级核对**（核法：`git show HEAD:docs/PLAN.md | sed -n '431p' | od -c | head`）。**引用 `:431` 时别把它当逐字凭据。**

§7 那条 DEFERRED〔现跑：`git show HEAD:docs/PLAN.md | awk 'NR>=1520 && NR<=1545'`；表头在 `:1515`〕：

- 表头（`:1515`）：`| 类型 | 项 | 为什么现在不做 | 完成判据（可验证） | 前置依赖 | 当前残缺表现 |` ⇒ 类型之外的**五字段**。
- `:1531` 逐字（**五字段齐，派单那句核对为真**）：
  > `| **DEFERRED** | **`task.list` / `task.cancel` 的实现**（D34 名册在册、生产注册表零枚） | 票 164 AC#1 现量（2026-09-26）：名册在册但生产注册表**零实现**，两把尺的读数见 `docs/evidence/s1/16x-contract-lines-c1.md` §4.4⑤；**先补名册裁定再谈实现** | 注册表里 `task.list`／`task.cancel` **各有一枚真实现**并过契约测试 | 票 163（同一批命令执行面）＋票 164 AC#1 裁定表 | Agent 侧查不到"谁被阻塞"，D31 那条"必须可见"**只剩 UI 半条腿** |`
- ⚠ **这条 DEFERRED 不覆盖 `task.output`**（项里只写 `task.list`/`task.cancel`）。⇒ 票面 AC#1 那句"AC#3 落地那天若仍零实现，必须变成五字段 DEFERRED 行或经 owner 批准从 D34 摘掉"（票面 `:23`）**今天没有落点**，它是一笔**尚未到期**的账，不是一条已存在的登记。

**射程结论（设计层，不是码）**：文本对 `task.output` 的**强制**只有三条 —— ①L0（不许弹卡）② 输出>4000 token 时**必须**在返回文本里带上"可续读的路径＋总长＋头尾"，只截不指＝不合格 ③ 归属 S7。**"怎么续读"的接口形状（游标？偏移？分页参数？还是就把路径交给 `fs.read`？）规格没有规定** ⇒ 属未定义即停的范围（AGENTS §2 精神／D22 闸门③），实现程要的是**一次定案**，不是自己填。

---

### 3.2 问② —— 盘上现状（"后台任务"这个概念今天有没有载体）

普查命令（两条都跑了）〔现跑〕：

```
grep -rlE '"task\.(output|list|cancel)"' --include=*.go internal/ cmd/ tools/   -> rc=1（零命中）
grep -rhoE 'task\.[a-z]+' internal/ cmd/ tools/ --include=*.go | sort | uniq -c -> 3  task.ctx
git grep -rnE '^type [A-Za-z_]*(Task|Job|Bg|Background)[A-Za-z_]* ' -- '*.go'   -> 4 枚（下表）
git grep -rn 'RunAsync' -- '*.go' | grep -v '^\.scratch'                         -> 仅 internal/agent/*_test.go
```

⇒ **派单那句"Go 代码里这个工具零实现"核对为真**；`task.ctx` 那 3 处是 `internal/statemachine/table.go:253` 的副作用串 `task.ctx-keep`（D43 转移表那一族），与 D34 名册无关，票面 AC#1 也是这么判的。

今天真实存在的**四个**"看起来像载体"的东西（逐个带 file:line，全部 `git grep` 证过）〔现跑〕：

| 符号／文件 | 它是什么 | 能不能当 `task.output` 的载体 |
|---|---|---|
| `type RunningTask struct`＝`internal/agent/loop.go:287`；`func (l *Loop) RunAsync(...) *RunningTask`＝`internal/agent/loop.go:321`；方法 `Root()`/`Cancel()`/`Wait()`/`Pending()`＝`:297/:300/:308/:318` | Agent 环路里的**内存态**后台任务句柄：ID、observe root/handle、`done chan`、`result Result` | **半条**：它是唯一"起一个东西、之后能读它"的形状，`Wait()` 给 `Result`（`internal/agent/loop.go:82`：`TaskID/Status/Text/Message/…/ToolLog`）。⚠ **但生产码里零枚 spawn 点**：`RunAsync` 的调用点**只有测试**（`internal/agent/control_test.go:81,131`、`forensics_test.go:99`、`loop_golden_test.go:172,251,275`）与 `.scratch/wisp/probes/158/**` 的桥副本。⇒ 今天没有"活着的后台任务名册"可查 |
| `type TaskLog struct`＝`internal/memory/models.go:58`；DAO＝`internal/memory/dao_tasklog.go:25/49/77/93`（`StartTaskLog`/`FinishTaskLog`/`TaskLogByID`/`ListTaskLogs`） | SQLite 的 L3 任务行（D20）：`ID/StartedAt/EndedAt*/State/QueryText/SummaryText*/Cost*/ErrorClass` | **不是**：字段里**没有任何输出正文**（`SummaryText` 是收尾摘要）。写它的只有 `cmd/wisp/run.go:614/619/626` 经 `internal/agent/journal.go:19/20/129/144` ⇒ **一次 `wisp run` 一行**，不是"一个后台任务一行" |
| `type ToolCall struct`＝`internal/memory/models.go:74`；DAO＝`internal/memory/dao_toolcall.go:17/98/111/121`（`ListToolCallsByTask` 等） | 取证行（D35）：`Tool/ArgsJSON/RiskLevel/Decision/Outcome/ErrorClass/CorrelationID` | **不是**：存**入参**与结论，**不存返回正文**。现量：`git grep -niE 'result_text|output|stdout' -- internal/memory/*.go | grep -v _test` 只命中 `artifacts.go:18`、`open.go:28` 两条**注释** ⇒ **DB 里没有"工具输出了什么"这一列** |
| `observe.Registry`＝`internal/observe/goroutine.go:262/344/383/405`（`Spawn`/`Count`/`Snapshot`/`RosterReport`） | goroutine 花名册；`RunningTask` 正是以 `agent-task-<id>` 名义 `Spawn`（`loop.go:324`） | **可当 `task.list` 的形状**（它是"谁在跑"的现成册子），但 `GoroutineInfo` 是诊断形状、**没有输出内容**，也不是给模型读的工具面 |
| `Spiller`／`Spill`＝`internal/agent/spill.go:36/44`，`Prepare`＝`:84`，桩文本＝`:140`，`artifactName`＝`:184`；接线路径＝`internal/agent/loop.go:238`（`NewSpiller(opt.Config.ArtifactsDir, b)`） | D15(3) 落盘机制：**已经在盘上、就是它** | **这是"输出正文"今天唯一的落点**：`<artifacts>\tool-output-<encoded id>.txt`，桩文本逐字 `%s\n[…输出已落文件：省略 %d 字符，总长 %d 字节 / 约 %d token，全文见 %s…]\n%s`（`spill.go:140-143`）。⚠ **命名按 tool-call id 编码（`artifactName(callID, seq)`，`:184`），不按 task id** ⇒ **"某任务的输出在哪"这个映射今天不存在** |
| `Store.ArtifactsDir()`＝`internal/memory/open.go:249`；`ListArtifacts`/`PutArtifact`/`ValidArtifactName`/`DeleteArtifact`＝`internal/memory/artifacts.go:96/71/91/280` | artifacts 目录与读写口（500MB LRU 配额） | **是读回全文的现成口子**（不经门控），但生产侧消费者只有 `internal/panel/attachments.go:61` 的接口声明（面板那一侧）⇒ **Agent 侧今天没有任何工具能读它** |

**零命中／占位（明写，别顺着规格编）**〔现跑〕：

- `git grep -rnE 'type [A-Za-z_]*(Task|Job)[A-Za-z_]* ' -- '*.go'` 只 4 枚：`internal/agent/compress.go:148 type traceTaskKey struct{}`（ctx key）、`internal/agent/loop.go:287 RunningTask`、`internal/memory/models.go:58 TaskLog`、`internal/proc/jobscope_windows.go:62 JobScope`（**Windows 专用**）。
- `internal/session/`＝**只有 `doc.go`（18 行）**，末行写着 `DEFERRED(SessionScope): implemented by ticket 28`；`internal/watchdog/`＝**只有 `doc.go`**，写着 `DEFERRED(watchdog loop/thresholds): implemented by ticket 42`。⇒ 两包今天**零实现**，`session > task > tool` 的 C11 层级只在文档里。
- `internal/tools/` 生产注册表里 `task.*` **零枚**；生产注册点只有 `cmd/wisp/run.go:341` 的 `tools.BuiltinFSEntries(...)`（`internal/tools/fs.go:330` ⇒ `fs.read`/`fs.list` + `BuiltinFSWriteEntries`（`fs_write.go:772`）＝票面说的六枚）。
- `cmd/**`／`internal/**` 生产码里**没有** "background／后台" 语义的起跑口：`git -c core.quotePath=false grep -rniE 'background|后台' -- cmd internal --include=*.go | grep -v _test` 命中的全是 `context.Background()` 与窗口绘制/文案注释。

**⇒ 本问的结论（一句话）**：**"后台任务"今天没有可服务的载体。** 最接近的两块拼不出这一发：`RunningTask` 有身份、有 `Cancel`、有 `Result`，但**生产里没人 spawn 它、也没有"跑完的输出留在哪"的概念**；`spill.go` 有全文与路径，但**只按 tool-call id 寻址、且没有任何工具能读它**。⇒ 这一票真正要造的**不是"读文件"那一腿，而是"任务级输出定位符"这一层**（谁产它、按什么 key 存、活结束/进程重启后还在不在）。

---

### 3.3 问③ —— AC#2–#5 判据预做（一格一张小表）

形状照票面既有判据：**未修码上响不响／修完必须怎么响／要造它需要哪一枚台件**。

#### AC#2「未修码读数：今天起一个后台东西再读它的输出能不能做到」

| | |
|---|---|
| 未修码上响不响 | **该响，而且我已经有支持它的读数**〔现跑〕：`RunAsync` 生产调用点＝0（只 `internal/agent/*_test.go`）；`task.*` 工具名在 Go 源码字面量＝0；`cmd/wisp/run.go:341` 只注册 fs 六枚；`ToolCall`/`TaskLog` 无输出列。**⇒ 做不到＝本票成立，派单的判断没被推翻。** |
| 修完必须怎么响 | 判据要**读一个真发生过的调用**，不是读名字：`wisp run` 走真桥（`internal/tools/bridge.go:247` 的 `b.reg.Lookup(req.Name)`），同一发 `task.output` 调用**修前**必须落到"查无此工具"那一支、**修后**必须返回真 `Result.Text`。⚠ **我今天说不出修前那一发的逐字读数**（没跑过桥，见 §2 第 2 条）⇒ 台件跑出来才算，别把 `Lookup` 的行为当已知。 |
| 台件 | `internal/agent/harness_test.go`（239 行）＋`internal/agent/testtools_test.go`（138 行）＝golden provider 起环路的现成套件；`internal/llm/adaptertest/mockllm.go`＝C5 注入口（AGENTS §1.3 许可的接缝）。注入面**只用这两个**，不要新造 mock 代替真桥。 |
| ⚠ 不可满足风险 | 若判据写成"注册表里出现 `task.output` 这个名字"⇒ 它与 AC#1 那把"名字在不在 Go 源码里"的尺同形，**测不出"能不能做到这件事"**；那是本票 §3.2 已经量过的零，会**在未修码上就"绿"**（恒真判据的形状）。判据必须落在**调用级**。 |

#### AC#3「输出太长那一格：截断必须带可找回的指针」

| | |
|---|---|
| 未修码上响不响 | **机制今天半响：截断有、指针有，但"续读"那一腿没人能走。**〔现跑〕`spill.go:140-143` 的桩文本已经带 `全文见 <path>`（形状正是 D15/`PLAN.md:2564` 要的"头＋尾＋总长＋路径"），但**没有任何工具能读那个路径**（§3.2），且文件名按 tool-call id 编码（`spill.go:184`）。⇒ 造一发"拿到路径、再读回全文"，今天**读不回来**＝本票那一腿。 |
| 修完必须怎么响 | 必须**测到"读回来"**，不是测"给出了路径"：同一发输出 >4000 token 时，桩里给的那个指针**必须能被 `task.output` 自己（或 `fs.read`）二次读出且字节对得上**；总长字段必须与实际全文长度一致。只断言"字符串里含 `全文见`"＝**只截不指的判据版**。 |
| 台件 | `internal/memory/artifacts.go:96 ListArtifacts`／`:71 PutArtifact`／`:91 ValidArtifactName`（现成的目录级读写与名字校验）；`internal/agent/truncation_test.go`（153 行）＋`spill_test.go`（310 行）＝已有的截断/落盘断言形状可照。 |
| ⚠ 两条**必须定案才能落**（不是实现程可以自己填的） | ① **续读接口形状**规格没写（§3.1）：游标／偏移／分页／还是把路径交给 `fs.read`。② **`fs.read` 那条"按需再读"今天是否真走得通我**没查**（§2 第 5 条）：`rt.paths` 只用 `cfg.FS.AllowedDirs` 构造（`cmd/wisp/run.go:327`），artifacts 目录是否在里面**无读数**；`PLAN.md:431` 说的"模型可用 `fs.read` 再读"若实际被 C26 allowlist 挡在外面，那**是一笔既有账**（不属本票、别顺手修），但它会决定 `task.output` 只能自带读腿。命令：`git show HEAD:cmd/wisp/run.go | grep -n ArtifactsDir`、`grep -n 'ArtifactsDir' internal/config/*.go` |

#### AC#4「取消要真取消：取消后不许有还在写的尾巴（变异：取消后仍写输出 ⇒ 必须红）」

| | |
|---|---|
| 未修码上响不响 | **今天无法判定〔原因〕：没有"被取消的后台写者"这个对象可让变异作用。**〔现跑〕生产侧零 spawn 点（§3.2）；`RunningTask.Cancel()`（`loop.go:299`）只是 `root.Cancel()`（ctx 取消），而 `ctx` 取消之后**没有任何"任务输出"在写**。⇒ 现在写这发判据，它要么打空、要么变成一发永不响的装饰。**本程不硬造。** |
| 修完必须怎么响 | 变异必须作用在**真写字的那个 fd/文件**上：取消之后那发写**必须失败或被观察到没发生**，红句要点名"取消后仍落了几字节"。可借的既有形状：`internal/tools/fs_edit_ac4b_kill_windows_test.go:176 TestFSEditRealTaskkillLeavesTheTargetWholeAndNamesTheStagingResidue`，同文件 `:36` 明写"There is no t.Skip in this file: a missing taskkill.exe fails LOUDLY" ⇒ **不许 skip 的形制**照抄。 |
| 台件 | `internal/proc/jobscope_windows.go`：`OpenJobScope():71`（带 `KILL_ON_JOB_CLOSE`，`:78-83`）、`StartInJob(cmd):110`、`Close():170`、`TreePIDs():157`。⚠ 这枚**只在 Windows 存在** ⇒ 直接落到问④ 的分母问题。另外 `internal/tools/cancel.go:19` 的 `CancelBus`（`Vetoed/Started/Report/Complete`）是 D31 现成的"取消后回报已做步骤"契约，`tools.Result.AppliedSteps`（`tool.go:54-59`）是它的载体。 |

#### AC#5「契约轴零字节：`docs/specs/**`、`internal/risk/**`、`internal/panel/**`、`thresholds.go`、golden、`allowlist.txt`、`frontend/**`、`design/**`」

| | |
|---|---|
| 未修码上响不响 | **这一格今天就能量，且形状是"改完之后必须仍为 0 hunks"，不是"现在响不响"。** 基线〔现跑〕：`git -c core.quotePath=false diff --name-only a651fe8 -- docs/specs internal/risk internal/panel tools/d22scan/allowlist.txt internal/observe/thresholds.go frontend design` ⇒ 本程自己**一条没出**（本程没动任何既有文件）。 |
| 修完必须怎么响 | 同上，`git diff --numstat <实现前锚点>..<锚后> -- <那族路径>` 的**变更文件数列**必须为 0。⚠ `internal/memory/**` **不在** AC#5 的禁改名单里，但它是票 103／106／119 那族证据件里反复自证的"未碰"面 ⇒ 若实现要动 DAO，**得在本票里明写，别静默扩面**。 |
| 台件 | 无新台件；`git diff --numstat`＋按路径筛（名册筛法：先 `git -c core.quotePath=false`）。 |
| ⚠ **提前发现的三处"零字节 vs 落地"张力（上报，不替它决定）** | ① **能力面是闭集**：`internal/tools/capability.go:17` 写着 "The frozen C3 set"（常量块 `:19-29`，11 枚：fs.read/fs.write/net/clipboard/shell/sysinfo/screen/input/window/notify/memory），`AllCapabilities` 的枚数由 `TestCapabilitySetIsFrozen` 钉住。`Registry.validate`（`registry.go:116-150`）在注册时就校 `validCaps` ⇒ **`task.output` 必须复用其中一枚能力，或者 C3 契约面被加宽**（第 12 枚＝C1–C32 契约变更＝**owner 批准**，不是实现程的选择）。候选只有 `fs.read` 或 `memory`，两者语义都不完全对得上"读任务输出"。名字形状本身没问题：`validToolName`（`registry.go:247-267`）只要求 `namespace.action`（`task.output` 合格）。<br>② **`internal/risk/provenance.go:84-92` 有一份按工具名的冻结标记表**（`sensitiveSourceTools`＝SPEC-06 §5 那 8 枚）。若设计要求"`task.output` 的返回内容打 taint/标记"，那**要动 `internal/risk/**` ＝ AC#5 零字节冲突**。缓解读数（同一文件注释 `:94-95`）："a `Mark()` for any other tool is still honored fail-closed, just logged" ⇒ 不加进名单**是**fail-closed，不是裸奔；但"要不要加"属定案。⚠ 顺带：`provenance.go:133-136` 的 `channelTable` 也是**按 D34 工具名**建的外泄通道表，`task.output` 不在册（`web.search`/`notify`/`clipboard.write`/`fs.write` 四枚在册）。<br>③ `internal/risk/rules_network.go:15-24` 那个 `"shell": true` 是**协议**deny 表（`deniedProtocols`），**不是工具命名空间表** ⇒ 别把它当"新命名空间要在 risk 里登记"的依据（我一开始差点读错，核了原文才排除）。 |

---

### 3.4 问④ —— 最小改动面＋跨平台分母

**落点清单（我的读法：这是"正解落哪"的普查结论，不是施工图；实现程可推翻，推翻时以它的读数为準）**：

| # | 文件 | 新增／改既有 | 为什么是它 |
|---|---|---|---|
| 1 | `internal/tools/task.go`（**新建**） | 新增 | 现成样板＝`internal/tools/fs.go`：一文件一族 `Tool` 实现 + `XxxDecl()`（`fs.go:296/309`）+ `BuiltinXxxEntries(deps)`（`fs.go:330`、`fs_write.go:772`）。`task.output` 走同一形状，`Decl.Declared: risk.L0`（形状见 `fs.go:300`） |
| 2 | `cmd/wisp/run.go:341` 那段注册循环 | 改既有 | 今天唯一的组合根注册点；`task.output` 要进生产注册表只能在这里 append（注意它上面 `:330-340` 有一整段"notify 为什么还不是工具"的注释——那是**同一枚命名形状冲突**的先例，别绕过去静默处理） |
| 3 | 任务级输出定位那一层的**新载体** | 待定案后新增 | §3.2 结论：今天没有"task → 输出"的映射。两条候选形状：**a)** 在 `internal/agent/spill.go` 的命名里加 task 维度（⚠ 会撞 `internal/agent/spill_name_injectivity_test.go`(368 行)、`spill_path_invariant_test.go`(706 行)、`spill_test.go`、`spill_acl_windows_test.go` 四枚既有断言）；**b)** 一个进程内的任务表（`RunningTask`＋`Result.Text`＋`Wait()` 现成）由工具读取（⚠ 零持久 ⇒ 重启后读不到，是否可接受属定案）。**本程不选** |
| 4 | `internal/tools/task_*_test.go`（新增，可多枚） | 新增 | AC#2/#3/#4 的判据都落在这里；跨平台那一族见下 |
| — | `internal/tools/capability.go`、`internal/risk/**`、`internal/panel/**` | **零字节**（AC#5） | 张力见 §3.3 AC#5 那三处；要动就必须先要批准，不许实现时顺手 |

**跨平台分母（真量过）**〔现跑〕：

```
docker --version                                        -> Docker version 29.6.2, build dfc9... (本机可读)
docker images --format '{{.Repository}}:{{.Tag}}' | grep golang -> golang:1.27 / golang:1.27-alpine
docker run --rm golang:1.27 go version                  -> go version go1.27.1 linux/amd64
```

⇒ 派单那句"本机 Docker 可跑 `golang:1.27`"**核对为真**（不是"镜像在"，是**真起得来**）。

CI runner 名册〔现跑：`grep -nE 'runs-on:' .github/workflows/ci.yml`〕：`ubuntu-latest` ×3（`:66`/`:278`/`:665`）、`windows-latest` ×2（`:388`/`:533`）、`[self-hosted, wisp-slo]` ×1（`:591`）。

`!windows` 用例的分母在哪〔现跑：`git grep -rl '//go:build !windows'`〕：仓里这一族有 20 个文件（`internal/tools/platform_other.go`、`internal/tools/staging_live_other.go`、`cmd/wisp/console_other.go`、`cmd/wisp/notify_other.go`、`cmd/wisp/resident_other.go`、`cmd/wisp/slo_other.go`、`internal/secret/protect_other.go`、`internal/winsec/*_other_test.go` 等）；`ls internal/*/*_windows*.go | wc -l` = **72**。

⇒ **本票的分母结论**：
1. AC#4 的"取消后不许还在写"若借用 `internal/proc.JobScope`（`KILL_ON_JOB_CLOSE`，`jobscope_windows.go:78-83`），那是 **Windows-only**：它的红**只能在 `windows-latest`（或本机 Windows）量**，`_windows_test.go` 在 linux 上根本不参与编译。
2. 反过来，`*_other_test.go`／`//go:build !windows` 那族**只能在 linux 容器量**（`docker run --rm -v <repo>:/src golang:1.27 go test ...`）。
3. ⇒ **判据里不许出现"本机绿＝跨平台绿"**：`task.output` 的跨平台腿要么明写"linux 容器里跑，命令＝上面那条 `docker run`"，要么这票就**只承诺 Windows 分母**并把 linux 那一腿登记成 DEFERRED（五字段齐）。这一条属定案，别由实现程自己填。

---

## 4. 提前排雷那一问（G5/G6/G7，以及本程新发现的一枚）

**做法声明**：只读 `.scratch/wisp/probes/154/gate-clauses.sh`，**未执行**（派单硬令；此刻另一程正在拿它做名册差集）。下面引的每行都是该文件在锚 `a651fe8` 的**文件内容**，不是它的读数。

名册（`:69`）：`LEGS_EXPECT='G1 G1b G2 G3 G4 G5 G5pos G5neg G6 G6pos G6neg G7 G7pos G7neg'`
形制（`:2-24` 那段注释＋`:314`）：每腿开跑前先 `want <ring|quiet>`＋`want_n <基线枚数>`；pair 腿"响＝**实测 > 基线**"（`:314` 逐字："响＝实测 > 基线；want_n 登记的那一发就是这一行的右边"），**拿不到 `want_n` 直接 rc=4**（`:133`），实测 < 基线＝`STALE`（`:156-157`）。

| 腿 | 现册声明（文件内容，逐行号） | `task.output` 落地会不会被点名 |
|---|---|---|
| **G5** OpenScope↔成对关（`want G5 ring` / `want_n 1`，`:384-386`；合方 `'CloseScope|Close\(\)'`，开方 `'OpenScope'`，`handle` 模式，排 `*_test.go` 与 `internal/risk/*`，`:386-387`） | 基线那 1 枚＝`cmd/wisp/panel_assets.go:243`（`:373-376` 逐字点名） | **会被点名，且是真漏**。新落点（`internal/tools/task.go` 或组合根）若 `prov.OpenScope(...)` 却不在**同一文件**里出现 `Close()`/`CloseScope`，实测变 2 > 基线 1 ⇒ 那枚 BAD 是**真漏**（开了 C25 scope 没关），**不是"名册要更新"**；名册该更新的只有"合法成对的新文件"（它同文件里既有 `Close()` 就不会被算进来） |
| **G6** OpenTask↔CloseTask（`pair … 'OpenTask' 'CloseTask' - … :!internal/tools/bridge.go'`，`:452-453`；`:51` 声明 G6 今天按设计就该响；反向 1 枚＝`cmd/wisp/run.go`，`:467`） | `run.go` 那一发 `CloseTask` 同时在 **G2**（`want G2 ring`/`want_n 2`，`:359-361`）在册 | 若 `task.output` 的实现**派发工具调用**（经桥 `OpenTask`）却没在同文件 `CloseTask` ⇒ G6 主尺**新增未成对**＝真漏。⚠ 注意 `:447` 那句：正向 0 枚**不是**"都关好了"，纯粹因为开方唯一生产调用点在 bridge.go 内部被排掉了 ⇒ 新增第二个生产开方点会**同时**动 G6 与 G2 的读数（G2 基线 2 → 3），**这一枚要提前打招呼**，否则下一程会把它当成别人跑歪 |
| **G7** Defer↔(New)DisposalScope（`':476-481'`；开方写 `Defer(Named)?`、合方写 `(New)?DisposalScope`，理由在 `:479-481` 逐字：`-Ew 'Defer'` 对 `DeferNamed` 是 rc=1） | `:51` 今天按设计该响 | 若新落点用 `plugin.DisposalScope`/`Defer` 挂收尾（任务级资源寿命最自然的挂法）⇒ 只要**同一文件内**成对就不响；跨文件挂收尾会被反向那一味点名。定义本体 `internal/plugin/disposal.go`（`:477` 给了行号：`type DisposalScope :129`、`Defer :186`、`DeferNamed :191`）被排，与 G5 排 `internal/risk` 同理 |
| **G3**（**本程新发现的一枚，派单没问但会咬人**）`want G3 quiet` / `want_n 0`，pattern `^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID`，射程 `internal/agent`（`:366-368`） | 今天 0 枚，且**声明为 quiet** | **只要实现程为了拿输出给 `internal/agent.Loop` 加一个收 `taskID` 的导出方法**（例如 `Loop.Output(ctx, taskID string)` 之类，§3.4 候选 b 的形状），G3 立刻实测 1 > 基线 0 ⇒ **BAD**。而这条腿在册的理由（`:366` 那句"Q-56 那一支落地"）是**另一码事**：它是用来盯 Q-56 的，不是给票 164 用的。⇒ **要么把"读任务输出"的入口放在 tools 侧而不是 Loop 的导出方法上，要么这一枚必须走一次名册更新（＝改尺，属另一程地界，且 `probes/**` 在本程禁改名单里）。本程把冲突写在这里，一字节没动。** |

**另附**（同一文件、与"新增生产码"直接相关的两条既有形状，实现程会撞到）〔现跑读文件〕：
- `tools/d22scan` 的 ban `pathresolver-bypass`（`tools/d22scan/main.go:15`）与 `bare-goroutine`（allowlist 里逐字：`bare-goroutine internal/observe/goroutine.go … Exempted BY FILE PATH only (never by call shape), so every other 'go <anything>' in the tree is a finding`，见 `tools/d22scan/allowlist.txt` 头 8 行）。⇒ `task.output` **不许**裸 `go func(`（要走 `observe.Registry.Spawn`，`goroutine.go:262`），也**不许**在 `risk.PathResolver` 之外用 `filepath.Clean|Abs` 做文件系统决策（AGENTS §1.2 逐字，`PLAN.md:1288`）——"按 taskID 去 artifacts 目录拼路径"正是最容易踩的那一脚。

---

## 5. 我这张表里，**今天无法判定**的格（不许当判据用）

1. **AC#4 的"未修码上响不响"〔今天无法判定〕**——原因：没有"被取消的后台写者"这个对象（`RunAsync` 生产零调用点、没有任何任务输出在被写）。**本程没有为它硬造判据**（见 §3.3）。
2. **AC#2 修前那一发的桥级逐字读数〔今天无法判定〕**——原因：本程禁跑测试与桥；我只证明了"注册表里没有这个名字"，**没证明**"发一个 `task.output` 调用会得到什么错误形状"（`bridge.go:247` 的 `Lookup` 失败分支行为未读）。
3. **`fs.read` 能否读回 artifacts 指针〔今天未查，不是"判定不了"，是没查〕**——见 §3.3 AC#3 那行末的两条命令。这一条会直接决定 `task.output` 该不该自带读腿，属**实现前一问**。
4. **G5/G6/G7/G3 的真读数〔本程没量〕**——派单禁跑；我只给了声明与行号。任何"它今天响几枚"的话都得由跑尺那一程出。
5. **linux 分母上"哪些用例真的会跑"〔未逐条核〕**——我只数了 `//go:build !windows` 文件（20 枚）与 `_windows*` 文件（72 枚），**没有**跑容器内 `go test`。

---

## 6. `next=`（给下一格的建议，不替它写码）

1. **先要一次定案，再开实现**（AGENTS §0 那句"未定义即停"）：三条都是**规格没说**的 —— ① 续读接口形状（游标／偏移／分页／交给 `fs.read`）；② 任务级输出的**持久性**（进程重启后要不要还读得到 → 决定 §3.4 候选 a 还是 b）；③ `task.output` 用**哪一枚既有 C3 能力**声明自己（第 12 枚＝契约变更，必须 owner 批）。
2. **AC#2 先做成"调用级"判据**（注册表里出现名字＝不算），注入面只用 `C5 LlmProvider` golden + `wisp run`（AGENTS §1.3）。
3. **AC#3 的判据必须断言"读回来"**：桩里的路径/指针与实际全文**逐字节对得上**，总长与实际一致；"只截不指"要能红。
4. **动 `internal/agent` 之前先读 G3 那一腿**（`:366-368`）：给 `Loop` 加收 `taskID` 的导出方法会把它从 quiet 打成 BAD。这条**优先于**任何"顺手加个方法"的冲动。
5. **动 `internal/tools/**` 生产文件之前，先确认那两把尺没在被验收程跑**（票面 `:4` 末句）。
6. 若最终决定"这票这一格先不落" ⇒ `task.output` 必须按票面 `:23` 那条硬账变成 §7 的五字段 DEFERRED 行（与 `:1531` 同形），**不许**继续"在册＋无实现＋无人认领期限"。
7. **`PLAN.md:431` 那处引文坑**（§3.1）与 **"`fs.read` 能不能读回 artifacts"**（§5.3）各是独立的小账，别混进本票的实现格。

---

### 附：本程写面自证

本程唯一写入＝本文件（新建）。既有文件零改动、零 `git add`（除本文件）、零 push。发现但**没有动手**的"一行改法"：无 —— 本程没发现任何"改一行就能过"的形状；§3.3 那三处张力全部是**要人拍板**的形状，不是笔误。
