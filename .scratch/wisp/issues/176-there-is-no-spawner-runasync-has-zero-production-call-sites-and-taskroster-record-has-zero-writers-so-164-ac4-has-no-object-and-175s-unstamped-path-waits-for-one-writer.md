# 176 — **没有一个起跑口**：`RunAsync` 生产零调用点、`TaskRoster.Record` 零写者 ⇒ 票 164 的 AC#4 没有对象、票 175 的"洗掉来源"那条路等着自动通、`task.output` 永远只能回答"查不到这个任务"——**而这件"谁去把后台任务真跑起来"的事，票池里没有一张票认领**

- Status: **立而不派**（编排者排程：**排在票 175 的落地腿之后**——那枚写者一动手就必须与本票同批，见 AC#2；且它与票 163／165 抢同一批 `internal/tools/**` 文件，串行）。
- 来源：**编排者本轮排程时自己撞出来的**：我为了执行 `A344` 里那句"票 175 必须与写者同批落地"去现量"那位写者归哪张票"，**结果发现归不到任何一张**——我此前一直把这件事记在票 163 名下，**那是错的**（票 163 是"常驻终端会话"，另一码事，见 `A345` 的自纠）。台账 `A345`。
- 关联：票 164（AC#4 的被测对象在这里诞生）· 票 175（**洗戳那条路今天不通只因为没人写**）· 票 154（已 `-done`，RESERVED＋触发门：宿主自带 task id 那一形"至今没人关"）· 票 163（**不是**本票，别再混）· `PLAN.md:1531`（`task.list`/`task.cancel` 五字段 DEFERRED 行）· D31（并发与审批模型）· D42

## 这是什么（人话）

昨天到今天我们把"读后台任务的输出"这枚工具做出来了，验收也过了，但它带着一条我抄进票面的原话：**"响亮地做不到，仍然是做不到。"** 原因是缺了最前面那一环：**今天没有任何一段生产代码会把一个任务放到后台去跑**。我把这句话量成了两个具体的数（都不是猜的，`164-v1` 现量、我复核过）：

- `RunAsync`（环路里那个"起一个后台东西"的入口）在**生产码里零调用点**；
- 新名册的写入方法 `TaskRoster.Record` **零写者**（只有测试在写）。

⇒ 所以真机上每一次调用 `task.output` 都只会得到那句"查不到这个任务"。**功能在册、有实现、没人喂**——这正是票 164 AC#1 当初要终结的那种"三不管"，只是前进了一格。

⚠ **更要紧的是：本票不落地，票 175 就只是一句登记**。那枚"外来的内容换个门读回来就不打戳"的路，**今天就差一个写者**；写者一出现，路自动通，而裁决程已现量确认**没有任何一枚既有用例、也没有任何一条 CI 步骤会因此变红**。

## 为什么值得做／不做的后果

不做：票 164 停在"有实现、不可达"；票 175 停在"已知洞、没人接线就没事"；票 154 那枚触发门永远等不到人来撞；`PLAN.md:1531` 那行 DEFERRED 的完成判据（"注册表里各有一枚真实现"）也永远不会满足。**最坏的形状不是"没做完"，是"三张票都以为自己不归口"。**

## 我的推荐（**owner 不用选，除非他要说"别做"**）

**先派票 175 的落地腿（小、独立、不依赖任何规格），再派本票**；本票第一腿**只做"起跑口该长在哪一侧"的落点裁定**（只读设计核），不直接写码。理由：本票一动就同时碰 `internal/agent`（环路那一侧）与 `cmd/wisp`（组合根），而票面定案④＋G3 那条腿明令**不许**在 `Loop` 上加"收 taskID 的导出方法"⇒ 落点是真决策，不该由写手临场挑。**不按推荐来的坏处**：一上来就派写手，最容易的那条路就是"给 Loop 加一个方法"——那会把名册上声明为安静的那条腿打红，而且是被当成"尺子误报"糊过去。

## AC（每格都要答"这一发在**未修码**上响不响"）

- [x] **AC#1 先把"没有起跑口"量成可复制的两个数**（不是新判断，是**把本票的存在钉在读数上**）：`RunAsync` 生产调用点枚数、`TaskRoster.Record` 生产写者枚数（两枚都要**排测试**的口径＋命令原文）。⚠ 若任一枚变成"≥1"，**本票作废、回来改票 164/175 的条件**，不许留着。
- [ ] **AC#2 与票 175 同批（硬约束，不是建议）**：本票任何落地腿，**必须与"外来的读回来要盖戳"那枚判据在同一次交付里存在**。⚠ 写者落了、戳还没有 ⇒ 那次交付**判不通过**（判据本体在票 175 AC#3，落地形状＝常驻正控用例，非尺子腿——这一条由 `175-c1` 现量裁出：`gate-clauses.sh` 那一族**没有任何 CI 步骤会跑它**）。
- [ ] **AC#3 取消语义在本票诞生**（票 164 AC#4 等的就是这个对象）：起跑口一落地，"任务被取消之后不许还有人在写它的输出"必须**同批**做成红／绿两向判据——⚠ 不许再写成"等 164 那格以后再说"（那是票 164 定案⑤登记过的同一枚形状）。
- [ ] **AC#4 落点形状禁区（逐条自证）**：(i) **不许**给 `internal/agent.Loop` 加"收 `taskID` 的导出方法"（票 164 定案④＋`gate-clauses.sh` 的 G3 腿 pattern，现量在 `:363-366` 一带，行号自己量）；(ii) **不许**把宿主内部 artifacts 写入做成受门控的 Tool（`AGENTS §1.2` 逐字禁止项）；(iii) 裸 `go func(` 一律走 `observe` 那套 owner/recover（`tools/d22scan` 的 `bare-goroutine`）；(iv) 非 `risk.PathResolver` 之外不许用 `filepath.Clean|Abs` 做文件系统决策（同一条）。
- [ ] **AC#5 契约轴**：`docs/PLAN.md`（含 `:1531` 那行 DEFERRED 的**完成判据列**若要改）、`docs/specs/**` 一字节不许动，动＝owner 单独批准＋批准原文落台账；`internal/risk/**`、`internal/panel/**`、`thresholds.go`、golden、`allowlist.txt`、`frontend/**`、`design/**` 零字节。⚠ 若实现程认为 D31/S7 的哪一句与本票冲突 ⇒ **停手上报**，不自填。

## 本票**不**解决

- 不做常驻终端会话（那是票 163）。
- 不做 `task.list`／`task.cancel` 的实现本体（`PLAN.md:1531` DEFERRED 在册；本票只把"起跑口"这一前置条件立起来，做完之后那两枚仍需各自的腿）。
- 不裁"跨进程重启的任务名册"（票 164 定案②已裁 v1 不做）。

## Progress log

- 09-28 00:0x｜**176-a1（只读设计核，零产码，未勾任何框）**：AC#1 两枚数现量完成，**票不作废**——`RunAsync` **存在**（`internal/agent/loop.go:321`）但生产调用点 **0 枚**（排测试；含测试 6 枚全在 `_test.go`），`TaskRoster.Record`（`internal/tools/task.go:123`）生产写者 **0 枚**（含测试 16 枚）。尺形原文与可复算件＝`docs/evidence/s1/176-background-start-port-census-a1.md` §1 ＋ `.scratch/wisp/probes/176/a1/census.sh`（读数 `readings2.txt`）。⚠ 停手条件（任一枚 ≥1）**未触发**；⚠ 上一枚假交件里"`RunAsync` 不存在"那半句地基读数与本程现量不符（表 §1 末），建议更正 `A356`。
  - 顺带多量一枚：`ArtifactPath:` 生产赋值点 **0 枚**（`grep -rn 'ArtifactPath:' --include='*.go' internal/ cmd/ | grep -v _test.go`）⇒ 票 175 那枚精确豁免的载具（`internal/tools/task.go:244`/`:253`）在生产里**今天恒不可达**——本条由 R3＋R7 自推，未引 `177-c1` 的表当证据。
  - **AC#3 归口本票（一句理由：判据必须与它的被测对象同票；164 定案⑤已把 AC#4 记成"连被测对象都没有"）**，但**今天不可能响**：生产无起跑口 ⇒ 无可被取消的后台写者 ⇒ 本格属"待起跑口同批诞生"，**本程不翻勾**。盘上最接近的 `internal/agent/control_test.go:77` 只断言环路 `StatusCancelled`、一字节不碰名册，**不是它的替身**。现成陷阱一枚带给写手：`internal/agent/loop.go:956-958` 的 `terminalWriteCtx`（`context.WithoutCancel`）是仓里今天唯一合法的"取消之后仍然写"（调用点 `:666-672`，注释逐字交代理由）⇒ 红／绿判据射程**必须显式排除**那一支。
  - **AC#4 逐条自证**：(i) 我的落点用现成导出的 `RunningTask.ID`（`loop.go:288`）＋`Wait()`/`Cancel()`/`Root()`，**不新增任何 Loop 方法**；🔴 **标红回禀**：G3 腿（`probes/154/gate-clauses.sh:363-366`，pattern `'^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID'`）是按**字面 `taskID`** 抓的，改参数名或包进 struct 则读数纹丝不动 ⇒ **"腿安静 ≠ 合规"，不许写手拿它当通行证**。(ii) 落点保持 `Spiller` 的宿主内部写入（`internal/agent/spill.go:29-31` 逐字 "deliberately not a gated tool"）；🔴 **标红回禀**：若落地腿为了让 C25 认下后台产物而把 artifacts 写入暴露成受门控 Tool，即正面撞 `AGENTS §1.2` 逐字禁止项。
  - **落点结论一句**：最不外溢的接点在 **`cmd/wisp/run.go` 的 `execute()`（`:586`、环路装配 `:588-617`、`loop.Run` `:622`）**——形状＝"调现成 `RunAsync` -> 拿 `RunningTask` -> `Wait()` 之后 `Record`"；⚠ id **必须仍由环路铸**，宿主自造 id 会掉进 `run.go:573-575` **已登记为没收口的开放端**（无 admission、无边界）。`internal/tools` 对本票是 consume-only、`internal/agent` 撞禁区 (i) ⇒ **三条路只剩 `cmd/**` 活的**。本程**未动 `cmd/**`**（派单 §5 禁改、未具名放开）⇒ 留给落地腿，且**落地腿的派单必须具名放开 `cmd/wisp/run.go`**，否则只能去撞 (i)。面板侧那条管子（`run.go:438-446` 的 `StreamLog`/`SnapshotPump`）**不同型、搭不上**，且 `internal/panel/**` 是 AC#5 零字节区。
  - 起手交代：`go test -count=1 ./internal/risk/ ./internal/tools/` 在锚点 `fbcdf441` 上 **risk 红**（`TestResolvePerCallBudget` 1.098ms 对预算 1.000ms，机器负载敏感那形；`thresholds.go`／判据一字节未动、也不许为变绿动），tools ok。门禁三枚全 rc=0：`scripts/d22scan.sh`｜`probes/154/gate-clauses.sh`（逐字"腿数＝14 声明与实测不符＝0"、G3 quiet 基线 0 实测 0）｜`tools/d22scan/runtests.sh -C tools/d22scan ./...`。**没跑第 4 枚**（`probes/161/r6/flip-declaration.sh`）、没跑全仓 `./...`。
  - ⚠ 同格另有一枚 **`176-a1b`** 在飞（写面 `-a1b.md`＋`probes/176/a1b/**`，与本条不撞名）；本条只登记 176-a1 自己的读数，**未读、未采纳它的任何结论**。

### 09-28 00:0x +08｜176-a1b（只读设计核·**重派枚**，零产码、未勾任何框）

- 锚点（本枚 §0 现跑）：`f6a4eebd38e8ce9d53878e8443ad063645e33b2b` / 短号 `f6a4eebd` / 分支 `dev` / `git status --porcelain -- internal/ cmd/` **空**。表＝`docs/evidence/s1/176-background-start-port-census-a1b.md`，台件＝`.scratch/wisp/probes/176/a1b/logs/`。
- **AC#1 两枚数（都排测试，命令在表 §1）**：`RunAsync` 生产调用点 **0**、`TaskRoster.Record` 生产写者 **0** ⇒ **两枚都没有 ≥1，本票不作废**，175 排程前提不变。
- ⚠ **要更正台账/派单的一处**：`RunAsync` **不是「不存在」**，它**声明在** `internal/agent/loop.go:321`（`func (l *Loop) RunAsync(ctx, input string) *RunningTask`），测试侧 6 枚调用点、生产侧 0 枚。**「不存在」与「存在但零调用点」是两种事实**：后者意味着落地腿**只需要有人调它**，不需要先造函数。票面正文第 11 行的口径本来就是对的，错的是转述它的那句台账（`1ed0bc14`）。
- **落点（表 §2，四条形状禁区之下）**：
  - **先例有**，且就是同一枚 `RunAsync` 自己：`loop.go:324` 走 `observe.Registry.Spawn("agent-task-"+id, "agent", ...)`（owner＋recover 全在 `internal/observe/goroutine.go:262/:281/:296`）；本仓生产码裸 `go func(` 实测 **0 枚**。同族先例另有 `internal/memory/writer.go:79`（队列→落库的回填）。
  - **D15(3) `Spiller` 可复用**（`NewSpiller` 是导出的、不附属于 `Loop`），`Spill.Text/Path` 正对 `TaskOutput.Text/ArtifactPath`；但**别裸送 taskID** —— `newTaskID()`（`loop.go:1107`）是小写十六进制 UUID 形，编码后就是它自己，会和**模型侧供给**的 tool-call id 落进同一目录同一前缀，撞上就是 `spill.go:117-131` 那条 last-writer-wins 静默换字节。⇒ 照 `loop.go:324` 的命名法加逻辑前缀（`agent-task-<id>`），零新增导出面。
  - **最不外溢的接点＝`cmd/wisp/run.go` 的 `execute()`**（`:622` 那一发 `loop.Run` 换 `RunAsync` ＋ 收尾 `rt.tasks.Record(res.TaskID, ...)`，约三行；名册今天与 `internal/` 零耦合，全仓只有 `:216`/`:361`/`:362` 三处引用）。**否决**了两条替代：挂 `EvDone`（`loop.go:922-924` 在 `res.Text==""` 时用状态文案顶掉 summary ⇒ 有损）、扩 `admitTask`（闭包拿不到 `Result`）。
  - ⚠ **本枚未动 `cmd/**`，留给落地腿** —— 三个候选全在 `cmd/wisp/run.go` 里，而只读腿没有这道具名放开。**不给这句，写手只能越权或去撞禁区 (i)。**
  - **面板侧没有可搭车的异步管**：`internal/panel/pump.go` 的出口是同步函数指针（`PumpSources.Out`，`:126-130`），文件头 `:15-23` 自陈「不是 transport、停在字节边界」。
- **AC#3 取消判据（表 §3）**：挂点裁在 `tools.TaskRoster` 那一层配 `tools.Stopped(ctx)`（`cancel.go:74-82`），**不挂桥**（`CancelBus` 的键是 correlation id、语义是单次工具调用）、**不挂 Loop**（任何 `func (l *Loop) X(ctx, taskID)` 直撞禁区 (i)）。红/绿两向形状已写清（绿＝取消后名册零变化、红＝故意留一处取消后仍写），并**明写这一发在未修码上「不可能响」**（今天没有起跑口）⇒ 属「待起跑口同批诞生」，**不是已通过**。`AppliedSteps` 那套管的是「已经落了什么」，**不能**顶替「别再写了」。
- **归口**：判据归 **本票 176**，不回 164 AC#4（表 §3.4，一句理由：被测对象由 176 造出来，164 只是先指出它缺位）。
- **门禁三枚全 rc=0**：`scripts/d22scan.sh`（clean）｜`probes/154/gate-clauses.sh` → 逐字「**# 腿数＝14 声明与实测不符＝0**」「名册=14 声明=14 记账=14 缺腿=0 空头声明=0」｜`tools/d22scan/runtests.sh` → `PASS=34 FAIL=0 SKIP=0`。**未跑** `flip-declaration.sh`、**未跑** 全仓 `./...`。
- ⚠ **起手那发逐包测试是红的，本枚没去放宽它**：`go test -count=1 ./internal/risk/ ./internal/tools/` → `internal/risk` **FAIL**（`TestResolvePerCallBudget`，C26 Resolve 实测 `1198746 ns/op = 1.199 ms/op`，预算 `1.000 ms`，981 样本，`pathresolver_budget_norace_test.go:34/:37`），`internal/tools` ok。本枚零产码、且 `internal/ cmd/` 起手即净 ⇒ 非本枚引入；按 `AGENTS §1.1` 阈值一字节未动、未重跑求绿。
- ⚠ **给排程的一条实测陷阱**：`gate-clauses.sh` 的 **G2** 腿 `want G2 ring` / `want_n 2`，而那 2 行里**有一行只是注释**（`cmd/wisp/run.go:564` 提到 `CloseTask`）。起跑口落地腿在 `cmd/wisp/run.go` 新写任何含 `OpenTask|CloseTask` 字面的说明 ⇒ 2 变 ≥3 ⇒ 该腿按「声明与实测不符」响。两条出路（避字面 / 走人工批准重登基线）都要编排者选，写手不许临场定。
- **禁区自证**：本枚结论不撞 (i)(ii)(iii)(iv)（表 §5，含尺形与行号）。另钉一条：**G3 安静 ≠ 合规** —— 那枚是文本尺，参数改名成 `id`/`tid` 就躲得过而实质照违。
- ⚠ **别让门禁绿被读成能用**：票 177 那枚豁免在生产里**今天仍是惰性的**（没有任何人填名册、`task.output` 也还没进 C25 名册），本票的起跑口**今天仍不存在**。
- **另有一枚在飞（本枚实测撞上，已按派单 §0.5#2 处置）**：本枚 23:54 取锚时 `probes/176` 还不在 `git status` 未跟踪清单里；00:01 现查已存在 `.scratch/wisp/probes/176/a1/{census.sh,readings.txt,readings2.txt}`，**mtime 23:58:59–23:59:35 全落在本枚工作时段内**且 `git check-ignore` 说没被忽略 ⇒ **先前那枚 `176-a1` 此刻真的还在跑**。本枚**没有读**那三枚文件的任何一个字节（只取名字/大小/mtime），它的任何读数与结论一律不继承、不比对；本枚写面已改名 `a1b`，盘上零重叠。截至 00:04 `git status --porcelain -- <本票面>` 仍为空 ⇒ 它**还没往票面上追加过任何东西**，故无原文可引；它若随后追加，以它自己那段为准。
- ⚠ 登记一条盘上事实（不是本枚结论）：本枚写表期间 `docs/evidence/s1/176-background-start-port-census-a1.md` **已出现且未跟踪** ⇒ 台账 `1ed0bc14` 那句「表不存在」到这里已过期；「提交 `150f798c` 不存在」仍成立（`git cat-file -t 150f798c` → `fatal: Not a valid object name`）。**两枚表的取舍不由本枚做**，本枚未读它的内容。
- 预算：**未超**（顶 ≤45）—— 含事后更正、提交与提交后校验在内不超过 38 次，确切数在最终回禀里给。零删除命令、零 push、`internal/**`/`cmd/**`/`tools/**`/`docs/PLAN.md`/`docs/specs/**`/`docs/reports/**` 零字节、别人的票面与 `probes/**` 既有台件零改动。

### 09-28 00:1x 编排者追加（两枚同格普查都真交件⇒ AC#1 翻勾；⚠ 台账 `A356` 里那句"`RunAsync` 不存在"是**我编的读数**，在此具名作废；三件裁定落定）

1. **AC#1 翻勾的凭据是我自己重跑的三条**（不是引两枚普查的数）：`grep -rn 'func .*RunAsync' --include='*.go' internal/` → **`internal/agent/loop.go:321` 在场**；`grep -rn '\.RunAsync(' --include='*.go' internal/ cmd/ | grep -v _test.go` → **空（0）**；`grep -rn '\.Record(' … | grep -v _test.go | grep -v slog` → **空（0）**。⇒ 两枚都＝0，**票 176 不作废、票 175 的排程前提不变**（票面 AC#1 那句"若任一枚 ≥1 本票作废"未触发）。
2. ⚠ **我这一轮最该记下的一条是我的错**：台账 `A356`（提交 `1ed0bc14`）里我写了"`RunAsync` 不存在＝`grep …` → 0 命中"并标了**〔我本轮现跑〕**——**我根本没跑成**：那条命令里前面一枚 `grep` 无命中返回 1，`&&` 链当场断掉，后面的命令一条都没执行，而我把"票面标题的说法"当成我自己的读数写了进去。⇒ **本行就是作废记录，`A356` 原文不抹**（更正走追加，见 `A357`）。**真事实是：`RunAsync` 存在、有实现、有 6 枚测试调用点、生产零调用者**——这把落地腿的活从"新造一枚函数"改成"调用现成那一枚"，是不同的活。
3. **两枚独立普查交叉核过同一批数**（`8eb2f320` 与 `b52de78b`），我逐枚现量后**采纳四条**：先例＝`RunAsync` 自己（`loop.go:324` 用 `reg.Spawn("agent-task-"+id,"agent",…)`，裸 `go func(` 生产 0 枚）；落点＝`cmd/wisp/run.go` 的 `execute()`（现读 `:622` 一带）；`Spiller` 可复用**但 taskID 不得裸送**（与模型侧 call id 同目录同前缀，撞上就是 `spill.go:117` 那处静默换字节）；取消判据挂 `tools.TaskRoster`＋现成 `tools.Stopped(ctx)`，**不挂桥、不挂 `Loop`**。
4. **我的三件裁定**（写手腿不用再问）：**(a) 具名放开 `cmd/wisp/run.go`**——边界＝只在组装点接那几行，其余一字节不动；撤销口令**「run.go 别动」**；**(b) 否决"改 `TaskRoster.Record` 为取消即拒写"那一支**（它会改掉 `internal/tools/task.go:121-122` 逐字的 **last writer wins**，那属契约轴邻近）——只用"`Wait()` 之后才 Record"＋"Record 之前问 `ctx.Err()`／`Root().Pending()`"这两条；**(c) 不为后台记录员新登记 D38 名**（那要动 `PLAN.md`＝契约变更）——复用 `agent-task-` 前缀；⚠ 若因此落不进档，`CategoryUnknown` 会 `slog.Warn`，**那是要看的症状，不许静音**。
5. **两枚在飞的陷阱，抄进派单**：**(i) G3 安静≠合规**——那条腿按字面 `taskID` 抓，改名或包进 struct 就纹丝不动，**不许拿它当通行证**；**(ii) G2 的 `want_n 2` 里有 1 行是注释**（`cmd/wisp/run.go:564` 提到 `CloseTask`），**在组装点新写含 `OpenTask|CloseTask` 字面的说明就会把那条腿顶红**。
6. **起手那发 `internal/risk` 红我判成争用、不是回退**（三枚读数）：两枚普查在 `23:51`／`23:54` 各自报 `TestResolvePerCallBudget` `1.098`／`1.199 ms/call` 对 `1ms` 预算；**编队空出来之后我自己跑**＝`-v` 单腿 `--- PASS (1.20s)`（**真执行、不是被跳过**）、整包 `ok 6.058s`。`thresholds.go` 与那枚判据文件自 `fbcdf441` 起 `git diff --numstat` **为空**。⇒ 判：**争用读数**，落地那天照普查 `next=` 第 7 条重跑一次留档。
7. **AC 现状**：本票 **1 勾／4 未勾**（AC#2 是同批硬约束、AC#3 取消语义、AC#4 形状禁区、AC#5 契约轴都要落地那天裁）。**排程**：`177 已落地` ⇒ 先派 **`175-r2`**（把 `task.output` 送进 C25 名册＋桥级正/反两判据；票 177 AC#2 的桥级那一发与 AC#3 都等它）⇒ 再派 **`176-r1`**（起跑口）。
