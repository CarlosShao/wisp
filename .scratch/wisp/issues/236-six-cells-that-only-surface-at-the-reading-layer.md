# 236 — **六枚只在"读数层"暴露的仪器缺口**（`221-v1` 与 `ci-red-1` 两枚非实现者腿交回）：两支 fail-closed 无尺／那把 DEFERRED 尺的 (b) 支扫词面无牙／"任务 id 必须宿主铸造"零钉／`ci.yml` 那句"今天会红"已过期／**在册红名册不是稳定集合**／**一枚刻意坏的夹具被入库之后污染了格式门的 tracked 分母、并吃掉两道 `go vet`**

- **Status**：**待派**（编排者 09-29 20:2x 立，起手锚点 `c5d88a7f`＝本票现读 `git rev-parse --short HEAD`）。
  来源＝`docs/evidence/s1/221-task-cancel-v1.md`（**37,097 字节**，我本人 `wc -c`；三发 `98007cc0`→`f6186158`→`c5d88a7f`）与 `.scratch/wisp/probes/ci-red/ci-red-1.md`（**35,452 字节**，我本人 `wc -c`，未提交）。台账 `A451`。
- **本票的射程不是"功能没做"，是"没人能证明它没坏"**。⛔ **零枚 AC 允许放宽任何现有断言**；凡"改门／改分母"那一支一律标〔契约邻接，要另批〕。

## 现量（起手逐条复算，别信这里的行号）

| # | 事实 | 谁的读数 | 我复跑了吗 |
|---|---|---|---|
| 1 | `taskCancel.Execute` 里 `caller == ""` 与 `t.d.Roster == nil` 两支 fail-closed，摘掉之后**各 0 枚红、rc=0** | `221-v1` 的 m4／m5 | ⛔ 未复跑（腿自述） |
| 2 | 票 221 新写那把 DEFERRED 尺：**(a) 支有牙、(b) 支无牙**——`teeth-m13` 摘掉 `task.list` 的标记行仍 PASS | `221-v1` | ✅ **我现读机制**：`internal/tools/task.go:279` 那行注释**同时含** `task.list` 与 `DEFERRED` ⇒ 词面计数被叙述句撑住（HEAD 命中 2 行：`:23` 与 `:279`） |
| 3 | 全仓**没有任何尺钉"任务 id 必须由宿主铸造"**；只有 `internal/agent/compress_trace_test.go:524` 一枚**形状**检查（36 字符＋4 枚连字符），射程是 trace 那一族 | `221-v1` | ✅ 我 `grep -rn newTaskID --include=*_test.go`＝5 命中，逐条读：`:524` 形状、`:557` 两枚不同、`pointer_183_cli_seam_test.go:59` 是指针注释不是断言 ⇒ **无"谁铸造"这一维** |
| 4 | `.scratch/wisp/probes/185/c1/mut/fs_broken.go` **已被入库**（`4813567e`，09-28 14:42）⇒ 进了 `lint::gofmt (gofumpt) - the tracked set is the denominator` 那一步的分母 | `ci-red-1` | ✅ **两半都我自己跑**：`git ls-files --error-unmatch` **成功**＝tracked；`gofumpt.exe -l` 该文件 **rc=2**、原文 4:1 "imports must appear before other declarations" |
| 5 | 今晚两发 push CI 红名册**逐字相同＝顶层 23 枚（含缩进子测试 25 枚）**，⛔ **不是我 A450 里写的 11 枚** | `ci-red-1` | ✅ 我认：我那 11 是用 `gh run view --log-failed \| grep -oE` 抽的，**任何一条过滤器都复现不出 23**（见 AC#5 那条尺） |
| 6 | `ci.yml:519-524` 那句"这一步今天在 windows-latest **会失败**（`RUNNER~1` 8.3 短名让 A 档判成 B）" | `ci-red-1` 称今晚实读 `--- PASS: TestPathResolverJunctionWindows` | ⛔ 未复跑（腿自述；且需要托管 runner 的环境才知道过期是否长期成立） |

## 为什么必须先修（不是文案洁癖）

- 第 1、3 条是**同一族**：两枚"拒绝"今天靠**代码里的那一行**存在，而**没有任何仪器**能在有人删掉它时报警。票 221 的整条安全结论（子代理不许停兄弟／不许停自己）就建在这两行之上。
- 第 4 条已经在**吃读数**：`lint` 那一步一红，后面的 `go vet (module)` 与 `go vet (tools/d22scan module)` **今晚两发零读数**（项目已知坑"一步红会吃掉后续步，skipped 不产日志＝读数永久采不到"，这次不是推测、是实测到了）。
- 第 5 条影响**所有"逐名比红名集合"的判据**：我这一整天用它做过若干次"零新增红"的结论，而名册本身在两发之间会换位 ⇒ 那条尺**缺了一半**（缺隔离复量那一半）。

## 落点（每格都要现跑读数；⛔ 不许用 mock 代替真的）

- [ ] **AC#1 两支 fail-closed 要有常驻尺**：给 `taskCancel.Execute` 的 `caller == ""` 与 `Roster == nil` 两支各立一枚判据，**形状＝摘掉那一支必然有具名用例变红**（照 `221-v1` 的 m2／m3 那两发的做法：各恰 1 枚红才算有牙）。完成判据＝两发突变读数（红→还原绿）＋工作树 `git status --porcelain -- internal cmd` 为空。残缺表现＝今天有人把任一支删掉，全仓零枚测试会响。
- [ ] **AC#2 那把 DEFERRED 尺的 (b) 支改扫能力**：从"数词面"改成"看 `BuiltinTaskEntries` 的注册名册里有没有 `task.list` 这一行"，并**先证它今天无牙**（`teeth-m13` 那形：摘标记行仍 PASS＝未修码读数）。⚠ 与票 225 分开算账：225 管"标记与 `SPEC-12 §5` 双向对账"，本格只管**这把尺本身有没有牙**。⚠ 附一条仪器事实要写进判据注释：**读盘型尺（`os.ReadFile`）对 `-overlay` 结构性不可见**，测它的牙必须 overlay 替换测试文件里的读路径并配正控。
- [ ] **AC#3 "任务 id 由谁铸造"要有一枚钉**：先量清今天有几枚尺会在"task id 换成模型给的值"时变红（我现读＝**0 枚**，只有形状检查），再定钉的落点。⚠ 背景两条硬事实：`internal/tools/task_backfill.go:38` 明写**普通 spill 路径的 call id 是 model-supplied**、且"裸 task id 是完全可以被够到的 call id"；`agent.newTaskID()` 熵源失败那一支**退化成时钟派生**。⛔ 本票**不许**顺手把 `allowed_dirs` 之类做成硬边界、也不许新造工具面（D34 无对应行）。
- [ ] **AC#4 `ci.yml` 那句过期注释要不要动**：只改注释、零行为改动；改前先**现读两发**确认 junction 那步真的 PASS（`ci-red-1` 给的是今晚两发，别拿它当长期事实——托管 runner 的环境会变）。⚠ 低利害、可逆，**不该拿去找 owner 拍**，由编排者或下一位读 CI 的人顺手定。
- [ ] **AC#5 "逐名比红名集合"这条尺要自带隔离复量那一半**：把规矩写成两发——**整包名册红 ＋ 单包隔离复量**（`221-v1` 实测：`cmd/wisp/TestTicket223RefusedLooseningKeepsOldValues` 整包 5.15s 红、隔离 **3/3 绿**（5.055／5.019／5.034s）＝计时红，不记账）；并把"在册常红名册今晚与实现腿那发不一致（两枚计时红换位）"重钉一次，落点写进哪份仪器说明由执行者定。⚠ 顺带纠正一条**我自己的坏尺**：`gh run view --log-failed` 抽名册会**漏计**（我抽到 11、真数 23），且**步名归属在部分 run 上全是 `UNKNOWN STEP`** ⇒ 数名册必须与 `portable-tests.sh` 自报的 four numbers **互咬**（今晚 4／7／12 三发各自相加＝23 对得上），并写清"名级 23 ≠ 包级 3"两个口径。
- [ ] **AC#6 tracked 分母被入库的刻意夹具污染**：`.scratch/wisp/probes/185/c1/mut/fs_broken.go` 是**故意编译不过**的正控夹具（`undefinedSymbol185c1`），它入库之后合法地进了格式门的 tracked 分母 ⇒ `lint::gofmt (tracked set)` 常红，并**吃掉后面两道 `go vet`**。三支摆清楚（⛔ 本票不预选，因为涉及门的分母＝〔契约邻接〕）：
  **甲**＝只把该文件的 **import 顺序排好**、保留"未定义符号"那一枚坏点（它的用途是符号坏、不是格式坏）——最小、零门改动；⚠ 代价＝它是票 185 已归档读数的对象，**改字节会洗掉别人引用的读数**，动之前要先核 `185` 名下有无 md5／字节数引用它。
  **乙**＝门的分母排除 `mut/` 这类刻意夹具目录——**动 CI 形状**，照票 134 的 C+B 先例要人工批准。
  **丙**＝不动，只把"这一步今天为真红、且它会吃掉两道 go vet"写进 CI 拓扑文档与停车点，让后续程别再当"CI 没测"是环境噪声。
  完成判据＝任选一支之后**必须真拿到一次 `go vet` 读数**（不是"应该能拿到"）。

## 禁区

- ⛔ 不放宽／不删除任何现有断言；不许为了变绿把 `lint` 那一步改成 `continue-on-error`。
- ⛔ **临时件只建不删**（`issues/README` 规则 8）：不许用"删掉 `fs_broken.go`"当作甲的省事做法。
- ⛔ 不碰 `frontend/**`／`design/**`、`PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）。
- ⛔ 工作树里别人的 46 枚脏改动（含 `docs/evidence/s1/152-*.md`、`probes/161/r6/logs/**`）一枚不许顺手提交。
- 派单必写：**产码腿的 AC 框一枚都不许碰**（票 221 刚为这一漏撞过一次，`A450`）；凡写进派单的 grep 尺先在本机跑一遍并把真实读数抄进派单。

## Progress log (append-only, newest last)

### 2026-09-29 22:5x 编排者增量 — 收 `236-a1`（只读普查腿），**顶正本票三处**＋新增一格＋重切一格判据

- 交付件＝`.scratch/wisp/probes/236/a1/census.md`（**43,613 字节／243 行**，我本人 `wc -c`／`wc -l`；九节齐全、"没做完"一节非空）。起手锚点 `82110540`，交回时 HEAD `5759d7fb`（同段推进 4 枚，全归 `235-r1` 与编排者）。
- **本腿合规**：零 `go test`／零编译／零 commit（`235-r1` 独占突变窗口）；跟踪文件零改动；只在 `.scratch/wisp/probes/236/a1/` 下建件。它并报了一处**我没写清的东西**：起手段里"编排者已暂存 5 条"——我现跑 `git diff --cached --name-status` ＝ **空**，那 5 条已随 `5759d7fb` 提交，索引现在干净，**不是遗留**。

#### 一、顶正本票"现量"表第 4 行（**我这行写错了，A451 同源的那半句一起作废**）

- 票面原句（第 4 行，⛔ 不抹）：「`fs_broken.go` 已被入库 ⇒ **进了 `lint::gofmt (gofumpt) - the tracked set is the denominator` 那一步的分母**」。
- **我亲跑的步级读数**（`gh api repos/CarlosShao/wisp/actions/jobs/109376664767 --jq '.steps[] | "\(.number)|\(.name)|\(.conclusion)"'`，逐字）：
  - `7|gofmt (gofumpt)|failure` ← **红的是这一步**（`ci.yml:136-144`，正文＝`gofumpt -l . tools/d22scan tools/mockllm`，走的是**工作树**，不是 tracked 名册）
  - `8|gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)|skipped`
  - `9|go vet (module)|skipped`／`10|go vet (tools/d22scan module)|skipped`
  - `11|staticcheck|failure`／`12|mockllm module vet|success`
- ⇒ **作废的是"哪一步红"这一句，不是"文件已入库"这一句**（后者我两半都自己跑过，仍成立）。"吃掉后面两道 `go vet`"**属实且更硬**：现在有 step conclusion 级实据（两枚 `skipped`）。
- 台账与停车点的同一处也要标过期（⛔ 不删原句，按台账规矩追加更正，落 `A455`）：`docs/reports/HANDOVER.md:397` 与 `docs/reports/pending-and-issues.md:8983` 写的「`-l` 名单**确实只有 4 枚**＋第 5 枚错误行」——**现量＝8 枚脏样**（7 枚在 `-l` 名单 ＋ 1 枚让尺退出 2）。多出的 3 枚是 09-29 之后新入库的（`197/r1c/pre/` 两枚、`222/v1/mutations/` 一枚）〔腿读数，名册见 census `### 全仓现量`〕。⚠ 谁按"4 枚"派活都会低估分母。

#### 二、AC#3 收窄：不是"零钉"，是**钉错了维**

- 我原句「先量清今天有几枚尺会在'task id 换成模型给的值'时变红（我现读＝**0 枚**，只有形状检查）」——**"谁铸造"这一维 0 枚，成立**；但票面背景把"裸 task id 会静默换字节"那一风险也一起算成了零覆盖，**错了**：`internal/tools/ticket176r1_start_port_test.go:128-168`（`TestTaskArtifactKeyStaysOutOfTheCallIDNamespace176r1`）**已在钉 `agent-task-` 前缀**，两半都钉（`:144-146` 前缀存在性＋`:152-167` 用"模型可重复的裸 task id 当 call id"真撞一次并逐字节比）。⇒ 摘掉 `task_backfill.go` 的前缀拼接，那枚用例**必红**。**前缀生效性不是零载具。**
- 收窄后的 AC#3 射程＝**只钉"键的来源"**：`Record`/`Look`/`Cancel` 对键的内容零假设（`task.go:215-218` 只挡空串），`cmd/wisp/run.go:892` 用 `res.TaskID`、`internal/agent/loop.go:323` 用 `newTaskID()` 填 `RunningTask.ID`——**"两者相等"这件事今天没有尺**（census 的修法 3，⚠ 需要一发真跑组合根的 CLI 接缝用例，候选载具＝`pointer_183/185_cli_seam_test.go` 那两枚形状）。
- 熵源退化支（`internal/agent/loop.go:1125-1131`，直引 `crypto/rand`、**无注入缝**）：**今天无载具** ⇒ 修法 2（加 `var randRead = rand.Read`）属〔形状改动，非本腿可批〕，且与 AGENTS §1.3"只在既有接缝注入"直接冲突 ⇒ **本格不许顺手为测试开后门**；性质钉（修法 1）若采用，交回时必须写明"覆盖不到熵失败"，否则是拿假绿换一条勾。
- ⚠ **口径更正（我的尺与腿的尺不同字，两条都要留）**：腿写「`grep -rn "没有宿主铸的 id" --include=*.go .`＝**1 命中**」；我单独复跑＝**5 命中**（`internal/tools/task_backfill.go:110` ＋ `.scratch/wisp/probes/176/r1/mut/{M1,M2,M3,M4}-*.go` 各 1 行拷贝），腿的口径未声明。**结论一致**：`--include=*_test.go` 现跑（范围＝`internal cmd tools`）＝**零命中（rc=1）**，两条字面量（`任务名册未接线`／`没有宿主给的任务 id`）同样 **`_test.go` 零命中**。

#### 三、AC#1 有现成载具，本格**零新 infra**（这是普查最要紧的一条）

- `Roster == nil` 支：`mustRegisterTaskEntries(t, TaskDeps{Roster: nil})`（**`internal/tools/task_output_leg_test.go:141-150`**，我现读到函数体）注册的正是 `BuiltinTaskEntries` 全家（含 `task.cancel`）；**同形先例已在仓**——同文件 `:113` `TestMissingArgsAndUnwiredRoster` 的 `:126-138` 用完全一样的形状钉了 `task.output` 那一支，只是它断的是子串 `未接线`。**⇒ 本格＝把那份形状抄到 cancel 的调用名＋**完整字面量**上**（只断 `IsError` 的尺没牙：摘支后掉进"查不到"或"不是调用者的孩子"那一支，照样绿）。
- `caller == ""` 支：`build221(t, windowGate221(), provider, false)` ＋ 现成 `x.cancel(t, "", target)`（`task_cancel_221_legs_test.go:74`／`:116-127`）；真桥不自造 id（`bridge.go:246-247` 只做 `CorrelationID==""→TaskID` 回落，两枚皆空⇒空）。
- ⛔ **`h222` 不能当载具**：`newH222` 的桥是 `Gate: NoGate{}` 且 `wireSpawn` 根本不注册 task 家族——NoGate 的 `PendingWindow` 答 `AnswerReject`，调用**进不了 Execute**，观察到的"拒绝"是那道门不是权限判定（＝假绿，正是 `task_cancel_221_legs_test.go:13-19` 头部注释逐字警告过的形状）。这条我已收进我自己的验收纪律。
- ⛔ `build221` 的 roster 写死在 `:101` ⇒ nil-roster 走不通，别给写腿留"改 `build221` 签名"这条路（动既有 8 处调用者，代价大于收益）。

#### 四、AC#2 补一层比"词面被叙述句撑住"更硬的机制（我现读复认）

- 尺本体 `task_cancel_221_legs_test.go:223-247`：`:224` 是 **`os.ReadFile("task.go")`**——**读盘型尺对 `go test -overlay` 结构性不可见**（overlay 只替换编译器看到的字节，测试二进制运行时打的是物理盘）。⇒ 今天**连"测这把尺的牙"都做不到**，不是没测、是测不出。这条必须写进判据注释。
- 词面被撑住的机制我也复认：`grep -n "DEFERRED" internal/tools/task.go` ＝ **3 行**（`:23`／`:32`／`:279`），其中同时含 `task.list` 的是 **2 行**（`:23` 名册抬头、`:279` 叙述注释）⇒ 摘 `:23` 后 `listMarked` 仍＝1，照绿。
- 改扫能力后的落点（现成尺，⛔ 不新造）：`task_cancel_221_legs_test.go:171` 那枚的 `len(names) != 2` 是**计数形**（假红风险：合法新增第三枚 task 工具会被误判）；`task_cancel_221_test.go:42-69` `allBuiltinEntriesHere` 与 `ticket175r2_stamp_live_test.go:349` 的分类表**已经是能力尺**。建议**保留计数、另加名集**（⛔ 删既有断言本票禁止）。"task.list 不在名册里"的名集形今天**零载具**（`grep -rn "task\.list" --include=*_test.go .`＝4 命中，没有一枚写 `registered["task.list"]`）〔腿读数，名集尺我未复跑〕。

#### 五、**新增一格 AC#1b**（普查顺带顶出、我裁并入本票；⛔ 不让它 silently 变成"以后加固"）

- [ ] **AC#1b `task.spawn` 侧同族三支同样零尺**：`internal/tools/subagent_197.go:253-255`（`Roster == nil`）、`:256-258`（`BaseOptions/ParentTools == nil`）、`:259-263`（`parentID == ""`）——三条字面量在 `_test.go` 零命中。判据形状与 AC#1 逐字相同（**钉完整字面量**＋摘一支必有具名用例红）。裁定理由＝同包（`internal/tools`）、同 harness、同"拒绝来自哪一道门"的坑，拆开会造成两腿抢同一枚文件＋两套载具。

#### 六、AC#5 的口径必须自带"仪器自己也数不全"这一半

- 现量（腿跑，机制我复认）：`scripts/portable-tests.sh:426-430` 那四把计数全是**行首锚定**（`^--- PASS`／`^--- FAIL`／`^--- SKIP`）⇒ **结构上看不到缩进子测试**；`scripts/winsec-tests.sh:134` 同形。⇒ "four numbers"与"名级 25（顶层 23）"**只能对拉顶层口径**，拿它当"子测试层也数过了"＝假绿（这正是我自己少算 37% 那一坑的镜像）。
- 规矩本体今天**只活在停车点散文里**：`grep -cE "红名|隔离复量|逐名" .scratch/wisp/issues/README.md`＝**0**，`scripts/` 与 `.github/` 亦 0〔腿读数〕。census 摊了四个落点各自的"会误导谁"，⛔ 它不代裁——**我裁：走它列的第 4 条**（把 `docs/reports/HANDOVER.md:331` 那枚正尺升级成"整包名册红＋单包隔离复量＋口径声明"三条，并追加一条 `A##` 指向它），不新开文档（本仓有"索引写了、文件不存在"的前科，AGENTS §4 末段 ⛔ 就是为此写的）。
- ⚠ 名册口径两件事写死进判据：今晚两发红名册**逐字相同**（25 行名，含 2 行缩进子测试）；与本机基线 `roster-baseline-f7478d37.txt`（**17 枚**）差 8 枚的原因是**跑范围不同**（本机包集 vs CI 全 job `--log-failed`）⇒ 判据里不声明口径就还会再撞一次"11 vs 23"。步名归因：`gh api .../jobs/<id>` 给得出准确步名（我上面亲跑复认），`--log-failed` 在部分 run 上是 `UNKNOWN STEP` ⇒ 给不出就用正文自报行（`runtests.sh: OK - packages=…`／`portable-tests.sh: four numbers`）认步。

#### 七、AC#4 改判：**不是"过期"，是"指错对象"**；本格今晚**不做**

- 腿给的正文读数：`ci.yml:519-524` 那句「the step **FAILS** on windows-latest today」对**那一步**为假——今晚两发该步 `--- PASS: TestPathResolverJunctionWindows (0.03s)` ＋ `runtests.sh: OK … PASS=1 FAIL=0`。而 **8.3 短名那一因没消失，只是不砸在这一步上**：同包里 `TestPathResolverShortNameAListDenied`／`…UNCAListDenied`／`…ExtendedLengthPrefixAListDenied` 三枚 FAIL，失败正文正是 `want "C:\\Users\\RUNNER~1\\…"`〔腿读数＋日志行号，我未逐行复跑；`gh api` 那一份我只能证步级结论〕。
- ⛔ 注释写在门的正文里，改它会被读成"动门的形状"（先例逐字在 `ci.yml:146-151`）。⇒ 我裁：**AC#4 与 AC#6 错开、由非实现者写腿各做各的**，我不自己写（我写了就同时是这格的实现者，后面没人裁）。判据补一句：改后的注释要写"读数＋日期＋run 号"，不写"今天为真"。

#### 八、AC#6 的**完成判据本身站不住**，重切＋摆一枚 Q

- 现量三支**各自都不能**满足我原先写的"任选一支之后必须真拿到一次 `go vet` 读数"：甲（只修 `fs_broken.go` 的 import 顺序）之后 step 7 仍红于**剩下 7 枚清单**（`bash -e` 的 `exit 1` 照样吃掉两道 vet）；乙（改 `attrib.sh:341` 的分母）**根本碰不到红的那一步**（红在 `ci.yml:136`）；丙（不动）零改动＝lint 一直红、vet 一直被吃。⇒ 原句「完成判据＝任选一支之后必须真拿到一次 `go vet` 读数」**作废**，按下面两支重切（⛔ 原句不抹，靠本节指回）。
- **AC#6a（口径，任何修法之前先做）**：把 tracked 分母里那 8 枚脏样逐枚列名＋归属票＋"为什么是坏的"（census 已给名册：7 枚 `-l`＋1 枚解析错误，全在 `.scratch/wisp/probes/**`）。
- **AC#6b（修法，二选一，我推荐甲）**：
  - **甲（推荐）**＝格式门的分母**剔除台件目录**（`.scratch/**`）——**这确实是动 `ci.yml` 正文**，属〔契约邻接〕，要 owner 一句话，立 **`Q-66`**。理由：台件是**证据件**，它的价值在**字节保真**，让 lint 去判它形状本来就错配；产品码那 555 枚仍被扫，覆盖率不减反稳。⚠ 代价要写清：tracked 分母从 **726 → 555**（171 枚免检），其中 `197/r1c/pre/subagent_197.go` 是**产品码改前快照**——免检不等于藏东西，真产品码那枚照扫。
  - **乙（不推荐）**＝逐枚把 8 枚脏样整理干净。⛔ 这会把**别人证据件的字节**改掉，而且 census 已证明票 185 的"防恒绿自证"正是靠 `fs_broken.go:4:1 imports must appear before other declarations` **那句错误文本**（`docs/evidence/s1/185-reread-owner-census-c1.md:48` 具名引用它）；甲形修完之后错误行变成 `undefined: undefinedSymbol185c1`，**rc=1 仍在但理由换句** ⇒ 后程复算 185 会判"凭据对不上"。票面原先写"改字节会洗掉字节数/md5 引用"——**现查是 0 枚字节数或 md5 引用**，真正的引用面是**行列与错误文案**（腿尺：`grep -rn "fs_broken" docs .scratch/wisp/issues`＝13 命中，其中 3 处以上引用 `4:1` 那句）。这条也顶正本票 AC#6 甲段的措辞。
  - ⛔ 两支都不许把 step 7 改成 `continue-on-error`、不许加 `if:`、不许删件（`issues/README` 规则 8）。
- **不答 `Q-66` 的默认动作**＝维持现状（lint 常红、两道 vet 一直被吃），**不阻塞本票其余五格**；代价＝后续程继续把"CI 没测"当环境噪声，而"球常驻资源达标"那类结论已经因此缺过一次证据。

#### 九、派单口径（写腿排队时照这个）

- ⛔ **本票所有写腿必须在 `235-r1` 之后**：AC#1／AC#1b／AC#2／AC#3 全在 `internal/tools`（同一枚包、同一批 harness、`subagent_222_test.go` 正在被改），并发＝互相洗读数。
- 一腿一格（或 AC#1＋AC#1b 同腿），AC#6 单独一腿且**不动 `internal/**`**；AC#4/AC#5 属文档面，可与 AC#6 之后的腿并。
- 派单必写：**产码腿的 AC 框一枚都不许碰**；凡写进派单的 grep 尺先在本机跑一遍并把真实读数抄进派单；**census 里所有 `file:line` 按最新 HEAD 重跑一遍再用**（本轮已实测"行号偏 1"两处：`task_backfill.go` 引 `:34-42`／`:38` 与 `open.go` 引 `:505`；漂移是共享工作树的常态，不是腿的读数错误）。
- census 的"没做完"九条里，**对本票有直接后果的是第 1 条**：AC#1／AC#2／AC#3 的**突变读数一枚都没跑**（⛔ 235-r1 在飞）。⇒ 这三格今天只有"落点＋载具＋形状推理＋静态零命中"，**没有实测温差**，验收腿必须补两发（红→还原绿）＋`git status --porcelain -- internal cmd` 为空。
