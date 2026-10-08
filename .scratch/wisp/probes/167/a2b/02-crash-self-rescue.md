# 167-a2b 格二：面板「崩溃自救」——Go 侧有没有"重启后状态"的真源

- 现量时间 2026-10-08 18:2x +0800；根锚 HEAD `8dd239f1`；全部 `git show HEAD:` / `git grep HEAD` 现量，零 Go 命令。
- 主尺（派单指定）：`git grep -niE 'recover|resume|restart|crash' HEAD -- 'internal/panel/*.go' 'cmd/wisp/*.go' ':(exclude)*_test.go'`。
- 判据（⭐）：面板**出向读面**（`Snapshot`／`ApprovalCardView`／视图那族结构体）里有没有一枚字段承载"上次崩溃/上次异常退出"。

## 0. 判语（先说结论）

- **没造**：面板出向读面今天**没有**一枚字段承载"上次崩溃/上次异常退出"（三族结构 27 枚键逐枚核过＋json 键面全扫 0 命中＋宽键名扫描 rc=1）。
- 主尺 47 行命中按用途剥完：**注释 32 行**（三族）＋**产码 15 行**（全属"配置 restart 档＋其审计"一脉，**只写审计行/枚举值，不供重启后状态**）。
- 与"重启后状态"最直接的 Go 侧文本是 `cmd/wisp/run.go:284-285` 的反向陈述（**不保留**，见 §3 ④）。

## 1. 主尺输出结构（rc=0；口径＝行；射程＝`internal/panel/*.go`＋`cmd/wisp/*.go` 排除 `*_test.go`）

| 分面 | 行数 | 命令 |
|---|---|---|
| 命中总行 | **47** | `git grep -niE ... \| wc -l` |
| 其中注释行（匹配在 `//` 注释内） | **32** | `grep -cE ':[0-9]+:[[:space:]]*//'`（rc=0） |
| 其中产码行 | **15** | `grep -vE ':[0-9]+:[[:space:]]*//'`（逐行见 §2） |

## 2. 逐枚剥离：按类型报枚数（行口径，全部具名可复核）

### ① 配置 restart 档（D36 三档生效级别 hot/reload/restart）＋其审计句——产码 **15 枚**，注释 **17 枚**

- 产码 15 枚（逐行）：`cmd/wisp/config_reload.go` **:115**（`rt.mgr.OnRestartPending = rt.reportRestartPending` 赋值）/:121/:179/:180（审计句）/:311（`func reportRestartPending`）/:315/:317/:320/:326（审计句）/:338（`var restartTierKeys`）共 10 枚；`cmd/wisp/panel_config_store.go` **:228**/**:**275（`res.Tier = panel.EffectiveRestart`）共 2 枚；`internal/panel/config_handlers.go` **:201**（`EffectiveRestart EffectiveTier = "restart"`）/**:**434（`case EffectiveRestart:`）共 2 枚；`cmd/wisp/resident_ball_windows.go` **:331**（printf"热键手改无需重启"句）1 枚。
- 性质：供的是"**某些配置改动要下次进程启动才生效**"这一提醒＋审计行，**不是**"重启后的运行状态"；重启之后它自己也无处被读（是写给日志/人的一句话＋枚举值）。
- 注释 17 枚具名：`config_readers_255.go:97`/`:315`、`config_reload.go:5`/`:18`/`:19`/`:302`/`:303`/`:329`/`:334`/`:336`/`:337`、`panel_config_store.go:29`、`panel_inbound.go:177`、`panel_resident_windows.go:176`、`resident_approval_windows.go:351`、`config_handlers.go:426`/`:427`。

### ② recover 边界纪律注释（D22 ban #1：`Registry.Spawn`/`run` 持有 recover 与 panic sink）——**8 枚注释，0 枚产码调用**

- 具名：`cmd/wisp/approval_reply.go:453`、`config_reload.go:117`/`:118`、`panel_host_windows.go:376`、`panel_resident_windows.go:55`/`:256`、`resident_task_source_windows.go:453`、`internal/panel/composer.go:104`。
- 射程内**零枚产码 `recover()` 调用**；8 枚全是"别在这里写裸 recover / recover 归 Registry"的纪律说明。

### ③ crash 字面注释——**3 枚**

- `cmd/wisp/doctor.go:103`（data dir 解析失败"lands here as a FAIL line ... rather than as a crash"——走 FAIL 不走 crash）；
- `cmd/wisp/panel_resident_windows.go:17`（"worse to ship than a crash"——面板线程嵌套泵风险说明）；
- `cmd/wisp/resident_windows.go:91`（"An exit request that arrives DURING boot is now a request, not a crash"——进程**活着时**启动期 Ctrl+C 登记，不是崩后恢复）。
- 三枚都是进程活着时的设计说明，不供任何状态。

### ④ "重启不保留"与装配注释——**4 枚**

- `cmd/wisp/run.go:284-285` 逐字："v1 keeps no roster across restarts (票 164 定案②), so a restart answers "查不到这个任务" rather than empty." ⇒ **反向陈述：重启后 Go 侧不记得任务表**；
- `run.go:686`/`:811`（审计/装配注释，含 `OnRestartPending` 归属说明）。

⇒ 剥完全部：**没有一枚供"重启后的状态"**；唯一产码族（①）供的是"改动待重启"的审计与枚举。

## 3. ⭐判据：出向读面字段逐枚核（没造）

- `Snapshot`（`internal/panel/composer.go:57`）6 枚键：`pending`/`results`/`composer`/`generatedAt`/`instructions`/`tasks`——无崩溃/异常退出位。
- `ApprovalCardView`（`internal/panel/approval.go:39`）10 枚键：`correlationId`/`tool`/`args`/`level`/`rulesHit`/`reason`/`reasonKnown`/`sessionOverrideBlocked`/`callChain`/`decidedBy`——无。
- `ComposerState`（`composer.go:235`）11 枚键（名册同格一件 §2）——无。
- 面板 json 键面全扫：`git grep -n 'json:"' HEAD -- 'internal/panel/*.go' ':(exclude)*_test.go' | grep -ciE 'crash|recover|resume|exit|restart'` ⇒ **0**（grep rc=1 即零命中）。
- 宽键名扫描：`git grep -niE 'lastExit|previousExit|lastCrash|unclean|abnormal|lastRun|dirtyExit|lastBoot|priorExit' HEAD -- 'internal/panel/*.go' ':(exclude)*_test.go'` ⇒ **rc=1 零命中**。
- ⇒ **判语＝没造**（"有报名"不成立，如实写"没造"）。

## 4. 排除项与旁证（为什么不能拿它们充数）

- `internal/observe` 的 panic 收集器**不许**当"面板知道自己崩过"（派单点名；本仓有过一枚已撤回的假话就长这样）：现量 `defaultPanicSink` 在 `internal/panel/*.go`＋`cmd/wisp/*.go`（排除测试）**rc=1 零命中**——面板与 cmd 侧根本不引用它，定义只在 `internal/observe/goroutine.go:167`/`:236`/`:245`。"崩溃栈被接住并写日志"那组读数归 167-a2 census／票 264 线，与"面板读面有没有崩溃状态"无关。
- 唯一与"死/恢复"沾边的面板侧注释是 `composer.go:101-104`（`NewSnapshot` 上方）逐字要点："The composer section is MANDATORY ... a panel that closes and reopens - or a WebView process that dies - recovers exactly what Go holds rather than an empty input row ..."——它说的是**重连即按 Go 现持有状态整发快照**（产码侧对应 `:112-115` mode 缺省 `"unknown"` 兜底），**不含任何"崩过"维度**、且本身是注释不是字段。

## 5. 判不动（缺哪把尺）

- "面板该不该有'上次异常退出'位"＝产品/契约判断（页面对应声明归 `Q-51` 面）：读数不裁。
- WebView2 进程真崩后**行为面**（重连延迟、快照重推是否覆盖页面残留、用户实际看到什么）需要跑面证据；本程零 Go 命令（253-v1 在飞禁并发 Go 面）⇒ 行为面判不动。
- ⛔ 未提议动 `frontend/**` 任何一寸。
