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

- [x] **AC#1 两支 fail-closed 要有常驻尺**：给 `taskCancel.Execute` 的 `caller == ""` 与 `Roster == nil` 两支各立一枚判据，**形状＝摘掉那一支必然有具名用例变红**（照 `221-v1` 的 m2／m3 那两发的做法：各恰 1 枚红才算有牙）。完成判据＝两发突变读数（红→还原绿）＋工作树 `git status --porcelain -- internal cmd` 为空。残缺表现＝今天有人把任一支删掉，全仓零枚测试会响。
- [x] **AC#2 那把 DEFERRED 尺的 (b) 支改扫能力**：从"数词面"改成"看 `BuiltinTaskEntries` 的注册名册里有没有 `task.list` 这一行"，并**先证它今天无牙**（`teeth-m13` 那形：摘标记行仍 PASS＝未修码读数）。⚠ 与票 225 分开算账：225 管"标记与 `SPEC-12 §5` 双向对账"，本格只管**这把尺本身有没有牙**。⚠ 附一条仪器事实要写进判据注释：**读盘型尺（`os.ReadFile`）对 `-overlay` 结构性不可见**，测它的牙必须 overlay 替换测试文件里的读路径并配正控。
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

- [x] **AC#1b `task.spawn` 侧同族三支同样零尺**：`internal/tools/subagent_197.go:253-255`（`Roster == nil`）、`:256-258`（`BaseOptions/ParentTools == nil`）、`:259-263`（`parentID == ""`）——三条字面量在 `_test.go` 零命中。判据形状与 AC#1 逐字相同（**钉完整字面量**＋摘一支必有具名用例红）。裁定理由＝同包（`internal/tools`）、同 harness、同"拒绝来自哪一道门"的坑，拆开会造成两腿抢同一枚文件＋两套载具。

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

### 2026-09-29 23:5x 编排者增量（二）— 收 `236-a2`：⛔ **AC#6a 已交付**，而 **`Q-66` 的两个关键数由我自己复跑确认＝甲法（分母排除台件）确实能让门变绿，乙法是个追不完的移动靶**

- 台件＝`.scratch/wisp/probes/236/a2/census.md`（**35,093 字节／231 行**，我 `wc -c`／`wc -l`；六节齐全、"没做完"9 条非空，另 13 份读数码件）。起手锚点 `98df640a`，交回 `a207d129`。零编译、零 commit、`git diff --cached` 空。
- **AC#6a（"先逐枚列名＋归属＋引用面"）＝已满足**：8 枚逐枚给了脏因（**7 枚是纯格式**、只有 `185/c1/mut/fs_broken.go` 是解析错误 `:4:1`，尺 rc=2 而其余 rc=1）、归属票、以及逐枚引用面判定（谁引的是**行列/错误文案**、谁只是路径提及、**零枚被引字节数或 md5**）。⇒ 本票"现量"表第 4 行那句"改字节会洗掉别人引用的读数"至此**彻底改口径**：真实风险不是字节数，是**"这枚文件是不是别人的逐字节可复核副本"**。
- ✅✅ **我自己复跑的两个数（这两个数决定裁定，不引腿的）**：
  - `git ls-files '*.go'`＝**731**；`git ls-files '.scratch/*.go'`＝**176**；⇒ 非台件那 **555 枚**里 `gofumpt -l` **零输出（rc=0）**。**⇒ 甲法（格式门分母排除 `.scratch/**`）一步就能让 step 7 变绿，并把 step 9／10 两道 `go vet` 从 `skipped` 放出来**（票面丙与"甲做不到"那句作废——腿在 census 的"甲"一节推的"还剩 7 枚所以照红"，是把它自己那次 `git ls-files -z '*.go'` 的**全量**名册当成了甲的效果；甲排的是台件目录，产品码那 555 枚今天本来就干净）。
  - 枚 5/6 的可复核性我逐字节验了：`git cat-file blob d1c440ad:internal/tools/subagent_197.go | cmp - .scratch/wisp/probes/197/r1c/pre/subagent_197.go` ＝ **IDENTICAL**（两边都 **19,314** 字节）⇒ **乙法格式化这两枚＝当场销毁一条"这是某 commit 那份产品码的逐字节改前快照"的凭据**（该承诺在 `docs/evidence/s1/197-subagent-entity-r1b.md:93` 具名写着）。
- ⚠ **乙法是移动靶（这条我原本没料到，腿也没料到，我自己复跑发现数还在长）**：现跑 `git ls-files -z -o '*.go' | xargs -0 gofumpt -l`＝**11 枚脏**（腿在 `a207d129` 量到 7；我跑时同段另有腿正在建突变副本 ⇒ **这是个会随每一枚写腿落盘而增长的数**）。⇒ 乙法"逐枚洗干净"之后**每轮都要重追**，甲法一次到位。⚠ 引用这个 11 时必须带口径（未跟踪、且含瞬时突变副本）。
- **一处腿读数我未复跑、标〔转述〕**：step 11 `staticcheck` 今晚同为 `failure`〔转述 `236-a1`，我 `gh api` 那份步级名册里确有 `11|staticcheck|failure`⇒ **这条我现在算复认了**〕。**但"门绿≠lint 绿"这句留着**：甲法买到的是两道 `go vet` 的读数，`staticcheck` 那一格的红是**另一件事**（在册 36 条只能在那版尺上量，本机版解不开 go1.27 的 export data＝已知坑，别把它并进来算）。
- **裁定状态**：`Q-66` 仍待 owner 一句（甲＝改 `ci.yml:136` 那步的分母，属动 CI 正文＝要人工批准；乙＝销毁逐字节副本＋追不完，**我不推荐，且现量给了销毁证据的具体形状**）。**不答＝维持现状，本票其余五格照走。**

---

## 2026-10-05 12:3x 编排者增量 — AC#6 的**两个关键数由只读腿 `269-a1` 复跑推翻**（原话不抹，按这一节读）＋**票 269 撤、并案回本票**

- ★**并案（记我）**：我 10-05 11:51 把"lint 的 gofmt 步被仓里自己的入库夹具顶成常红、并把三步无守卫的步吃回采不到"**另立成了票 269**（commit `05680e44`）——查下来这**就是本票 AC#6 同一枚物理缺陷的第二次立案**（`A575` 那条裁定逐字写着"同一物理缺陷链上只记一次"）。⇒ **票 269 就地撤（口令「269 撤」已用）**，普查件保留在原路径不动、不改写归属：`.scratch/wisp/probes/269/a1/census.md`（**350 行／68,999 字节／占位 0**，两笔 `fdffddd5`→`0f690a1c`）此后**当本票 AC#6 的凭据件引用**。⚠ **字母冲突必须重命名**：本票的"甲／乙"（甲＝分母剔除台件、乙＝逐枚整理脏样）与那件普查的"甲／乙／丙"（乙＝样本搬家、丙＝补 `if:` 守卫）**不是同一支**，今后引用一律改写成 **形 I＝分母剔除台件／形 II＝样本搬家／形 III＝补 `!cancelled()` 守卫**。
- ⛔⛔ **`A457` 那两个数已过期**（我那节逐字写"甲法（分母排除台件）确实能让门变绿，乙法是个追不完的移动靶"）：**形 I 单独做今天不让门变绿**——现跑名册（HEAD 整树归档 × gofumpt v0.12.0）＝**tracked 分母脏 19 枚（起手锚）→ 20 枚（普查期间 `50e0299c` 12:09 在飞写腿入库 `probes/259/r1/gofumpt-negctl/probe.go`）**，其中**一枚根本不是台件**＝`tools/d22scan/selftest.go`（首入库 `5e8748b3`，10-03，属**真违规**、归 161 那批我代提的码）；且形 I 之后链头红只是从 `ci.yml:168` 挪到 `:184`（那一步喂的是 `attrib.sh:341` 的 **tracked 全集、不剔样本**）⇒ **"吞三步"照旧**。⇒ **`Q-66` 要重新摆的两个数＝①剔完台件还剩哪几枚（今天＝那枚真违规）②修完它之后链头红落在哪一步**，⛔ 不许再拿 09-29 那句"甲法能让门变绿"去要那句话。
- ⚠ **一枚仪器坑具名（我自己在派单里差点用错）**：`gofumpt -l` 跑**本机工作树**在这台 Windows 机器上会被 CRLF 污染——我现跑 `git ls-files -z '*.go' | xargs -0 gofumpt -l` 得 **32 枚**（其中 26 枚在 `.scratch`、看起来"6 枚不在"＝`cmd/wisp/models.go`／`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`internal/risk/provenance.go`／`internal/tools/bridge.go`／`tools/d22scan/selftest.go`）；**同一条尺在 HEAD 归档上＝19–20 枚、非台件只有 1 枚**。⇒ **判 CI 那道门的分母只许用归档形状跑**（`git archive HEAD` 或 `git cat-file blob HEAD:<path>`），本机那 5 枚多出来的是**行尾不是格式**（`A615` 已记过一次预存 `cmd/wisp/models.go` CRLF）。这条写进 `Q-66` 的续料，⛔ 谁拿我本机那 32 枚去裁都是裁错。
- **本轮落地的那一小格（不需 owner 批准）**＝`236-r1`：只把 `tools/d22scan/selftest.go` **那一枚非台件真违规**格式化到 gofumpt 干净（它属 161 那批、不是任何人的证据件），并**报修完之后 tracked 分母还剩几枚、是不是全在 `.scratch`**（那正是 `Q-66` 缺的那枚数）。⛔ 不许碰 `.scratch/**` 任何样本（形 II 未被裁）、⛔ 不许改 `ci.yml` 正文（形 I 要 `Q-66`）、⛔ 不许给任何步加 `if:`（形 III 未被裁）。

### 2026-10-06 14:2x 写腿 `236-r2` 增量 — **AC#1＋AC#1b 同腿交回：五枚拒绝分支全部装上常驻尺，八发突变读数全在件**

- 交件＝`.scratch/wisp/probes/236/r2/evidence.md`（五节全实、占位 0；另 `logs/overlay/*.json`＋`logs/diff/*.diff`＋`logs/mut/*.txt` 八发原始日志）。起手锚 `3d9b8374 10-06 13:12`，写面＝新建 `internal/tools/failclosed_236_teeth_test.go`（六枚用例，**零产码改动、零导出名、零放宽既有断言**）；⛔ 七枚 AC 框／`Status:`／改名本腿一枚没碰（归编排者）。
- 票面"现量"表第 1 行那笔 ⛔ 未复跑＝**今天结清**：载体＝`go test -overlay`（⛔ 不在共享工作树种针），m-1a／m-1b／m-1c／m-1e **各恰 1 枚具名红**、m-1d 整支摘＝2 枚（`||` 两半一起没）、m-1d1／m-1d2 两半各摘＝各恰 1 枚，八发全部**还原后绿**（rc=0，PASS=204／FAIL=0）；基线 198 全绿⇒名册里没有别人的红可混。票面 §72 预言的"摘支后掉进『查不到』那一支照样绿"**被实测命中**，这是"只断 `IsError` 无牙"的直接证据；五枚钉因此都断**完整字面量**（`!=` 逐字节相等，不是 `Contains`）。
- ★★★ **一条本票没预料到的仪器读数（请编排者裁，本腿不修也不判）**：`subagent_197.go:256-258` 那道 guard 今天**同时充当 `:278` 的 nil 防护**——摘整支且 `BaseOptions == nil` 时是一次**空函数值调用 panic**，第一版读数＝`rc=1, PASS=25, FAIL=2`，**其后约 179 枚用例零读数**（票面 §九 说"并发＝互相洗读数"，这次一枚未受管的 panic 自己就洗掉了整包）。修法落在测试面（panic 收成该用例自己的 1 枚红，只可能把"崩"变红、不可能把红变绿）；⛔ 产码那个洞一字未动，性质与归属见件 §5 第 2 条。
- 另两条要请编排者裁的（本腿只报读数）：§5 第 3 条——`||` 复合 guard 上"各恰 1 枚红"有三种读法，三发证据都摆齐了，**判语不预填**；§5 第 4 条——AC#1 `caller == ""` 那一枚的 **`ErrorClass` 这一维今天量不到**，缺的具体一行＝现成载具 `x.cancel`（`task_cancel_221_legs_test.go:116-127`）把 `ErrorClass` 丢了，修它＝改别人共用的 helper 形状，超出本格射程（⛔ 这与"读盘尺对 overlay 不可见"不是同一条，本腿五枚都是行为尺，overlay 全部量得到）。
- 门禁：`sh scripts/d22scan.sh` rc=0（正控先绿：35/0/0；正文 `clean - no D22 ban violations`，ban #8 覆盖 `internal/` 513 枚含 `_test.go`）；`gofmt -l` 与 `gofumpt v0.12.0 -l` 对本腿新文件**都零输出**；`go vet ./internal/tools/` rc=0；`git status --porcelain -- internal cmd` 起手 0 行、commit 后回到 0 行。
- 2026-10-06 17:0x 写腿 `236-r3c`（票 236 AC#2，一腿一格）交回＝`git commit a3a7d535`：四发突变**本腿自己重跑**（正控／m-2b／m-2e／teeth-m13）＋三发仪器配对，判语写满 `.scratch/wisp/probes/236/r3/evidence.md` §0–§6（原始输出 `logs/r3c/**`＝mut 23／overlay 7／gates 14；`236-r3` 那 36 行骨架逐字留档在 `evidence-skeleton-236r3.md`、未删）；**读数三条**＝①teeth-m13（只摘 `task.go:23`）红名册恰 1 枚、归档词面尺与能力钉都 PASS，副本上共现行数 2→1 ⇒ 票面要的"今天无牙"**机制性坐实**；②m-2b（注册 `task.list`）＝旧尺全绿、新尺具名红（红句自报名册三名）⇒ 换扫能力买到的牙；③m-2e（注销 `cancel`）＝10 枚具名红，反向也有牙；**仪器事实**（读盘尺对 `-overlay` 结构性不可见）已按票面要求配"替换测试文件读路径"的正控变成可量的：改读路径后归档尺在"两枚标记全摘"的字节上会红（`task_cancel_221_legs_test.go:245` 红句在件）⇒ "无牙"的精确含义＝对 overlay 无牙、不是谓词写坏；§5 七条判不动不空白（(a) 支与 m-2c 今天**未实测**、本格故意拦不住"注释说谎"那一形并写明为什么不造那枚词面尺）；门禁＝d22scan rc=0（正控 35/0/0，ban #8 internal/=514）、gofmt／gofumpt v0.12.0 对本腿文件零输出、go vet rc=0、未突变态三枚 PASS；`git status --porcelain -- internal cmd` 起手与四发后各 0 行；⛔ 七枚 AC 框一枚未勾、`Status:` 未动、零产码改动（对载体那枚测试文件只改注释三处）。；★本腿自记一条写错的引用（原句不抹、另见件 §5 第 8 条）：第一版把注释改后字节里的行号**猜**成 `:144`／`:163`／`:201`／`:207`，实测＝`:155`／`:174`／`:212`／`:218`／`:222`（对照表已落盘 `logs/r3c/gates/line-number-map.txt`）⇒ 更正提交＝第二条 commit，⛔ 不影响任何 rc／红名册／红句正文（那些全取自日志原文）。；★同一轮复量再抓到本腿自己第二处过头措辞（§1.6 那把 grep 尺原写「四行全不是比注册名」，复跑发现其中一行就是本腿那枚钉 `:154`）⇒ 改准为「除本腿那枚之外三行都不是比注册名；本腿落笔之前＝零枚」，判语方向不变、口径变严。两处更正同在这笔里。
- 2026-10-06 17:2x 终裁腿 `236-v1`（票 236 AC#1＋AC#1b＋AC#2 三格独立复认，非实现者）起手＝锚 `f87a696c`：本腿唯一 go 腿、五把尺全过（`sed -n 3,5p` 零撤票口令／5 前缀 porcelain 0 行／写腿七枚提交 `c06d569a`→`f87a696c` 逐枚验存在＋证件 237 行与 448 行 47,553 字节复量一致／未勾 7 枚已勾 0 枚／`internal/tools/task.go` 盘上＝HEAD＝`4138177e…`），基线 `go test ./internal/tools/ -count=1` rc=0 顶层 206/0/0 已抄名，骨架＋§0 落 `.scratch/wisp/probes/236/v1/verdict.md`，突变只在 `go test -overlay`（副本 `D:/tmp/wisp236v1/`），⛔ 七枚框一枚不翻、`docs/evidence/s1/` 零写。

## 编排者翻勾记录（2026-10-06 18:3x +08）

凭据＝**非实现者终裁腿 `236-v1`** 的裁决表 `.scratch/wisp/probes/236/v1/verdict.md`（289 行／63,098 字节，五节全实、占位符尺 0 命中；三笔提交 `ea5b2dc9`／`55324446`／`1e8200b7`，编排者逐枚 `git log` 验存在）。⛔ 本编队没有一把框是实现腿自己勾的（`236-r2`／`236-r3c` 交件时七枚全空，§4 第 10 行复量未勾 7／已勾 0 与起手逐字一致）。

- **AC#1 → 翻勾**。判语＝成立。裁决腿自己造了 6 发突变（不是复引用腿日志）：`v-m1a`／`v-m1b` 两枚删除发各拿**恰 1 枚具名红**（`Test236R2TaskCancelRefusesWhenHostGaveNoCallerID`／`…WhenRosterIsUnwired`），另加 `A-1`／`A-2` 永假条件、`A-6` 翻 `IsError`、`A-7` 两句字面量对调 —— **四形全红**。⇒ 牙咬在行为上，不在字面量的存在性上。⚠ 本票 §72 我原先预言的"落到相邻分支照样绿"**被推翻（往好的方向）**，推翻它的读数在 §2 第 1 行。
- **AC#1b → 翻勾**。判语＝成立，三支都有牙：`v-m1c`（`Roster == nil`）／`v-m1d1`＋`v-m1d2`（`\|\|` 两半各恰 1 枚红）／`v-m1e`（`parentID == ""`，同一用例内三处断言一起红，含"名册行数＝2 want 1"那枚**后果钉**）。★`236-r2` §5 第 3 条留给验收腿的那一裁（"一枚语句管两枚字段"时"各恰 1 枚"怎么读）**已由裁决腿明确判定**：半边摘⇒恰 1 枚红是兑现形；整支摘⇒2 枚红是**正确反应、不降级**（判它破会反向激励拆语句，而拆语句正是产码改动）。
- **AC#2 → 翻勾**。判语＝成立，而且**派单点名的两枚欠量由裁决腿补齐**（我此前在 `A641` 第 4 节把这两格记成欠账，现在按新读数升回来——⚠ 我的降级只能由新读数升回来，这条就是那次）：①(a) 支＝`B-2a` 注入一行喂给旧尺 ⇒ 恰 1 枚具名红，正控 `B-2b`（逐字节相同副本）rc=0／206／0 证明红不来自"改读路径"这个动作，`B-2c` 证明纯 overlay 时它读不到 ⇒ **(a) 支谓词有牙、瞎在载体**；②`m-2c`＝`B-3` 注册无关的 `task.note` ⇒ 枚数形旧尺**误红**（`:184`）而新钉按名字集合不误红 ⇒ "保留计数、另加名集"这个落点**实测选对**；③决定性那发是 `B-6`：把标记行已删的字节**绕过 overlay 从读路径直接喂给旧尺，它照样 206／0** ⇒ 无牙是**谓词本性**，比写腿那份（只证到 `B-1`＝overlay 没落地）硬一层。

**跟着这三格一起入账、但不在本票射程的三件（⛔ 不许被读成"本票已闭合"）**
1. ★**产码洞 `internal/tools/subagent_197.go:278`**：那道 guard 今天**同时充当 `:278` 空函数值调用的 nil 防护**（`A-4-without-recover` ⇒ PASS=25／FAIL=2／其后约 181 枚零读数；正控 `ctrl-no-recover-only` 只换测试文件 ⇒ rc=0／206／0，所以失明只能归因产码那一次调用）。裁决腿明确不判该不该修（超射程），⛔ 但留下一条**后程硬约束：任何把 `:256` 拆成两支语句的修法，必须先处理 `:278`**。⇒ 具名立案＝**票 270**。
2. **附带读数：票 221 那把枚数尺会在合法新增第三枚 task 工具时假红**（`B-3`，`task_cancel_221_legs_test.go:184` 数的是 `len(names) != 2`）。⛔ 裁决腿不判它算缺陷、也不许任何人为变绿去放宽它——那是**别人的尺的地界**。⇒ 具名归口＝**票 271**（票 221 已 `-done`，按治理裁定"勾交付物＋把洞归口另一张票＋原话不改"）。
3. **`ErrorClass` 那一维今天没有可钉的牙**：裁决腿实测 `P-1`／`P-2` ⇒ 那一格**手上就握着 `agent.ToolOutcome`、`ErrorClass="tool"` 一跑就读得到**，所以 `236-r2` §5 第 4 条那句"无法断"在措辞上**说过头了**（缺的是"没读"不是"拿不到"）；而 `P-3`／`P-4` 证明分支摘掉后三元组**逐字不变**（`bridge.go:583-587` 对所有 `IsError` 一律书 `ClassTool`，`tools.Result` 没有字段供工具自填）⇒ 就算补断言也不会因突变而红。AC#1 因此**不降级**；要让 D37 那一维能区分"是哪道门拒的"＝**产码改动＋〔契约邻接〕**，本票不动，登记在台账 `A643` 等排程。

**另记两条给后程引用本件的人**：①那枚仪器尺 `Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt` 在**任何**把不同 `task.go` 字节 overlay 进来的全包发上都红（25 发里红 11 发／绿 3 发），看到它红应读作"overlay 换进了不同字节"，⛔ 不是新常红、更不许"修"成少比几样；②裁决腿自记一次越界（临时件建在写面外＋`rm -f` 删过一枚，后果量零），编排者裁见台账 `A643`，**原句不抹**。

**剩余四枚（AC#3／AC#4／AC#5／AC#6）保持未勾**——裁决腿 §5 第 9 行逐字声明"本件不答"，缺的是各自的读数而不是判语：AC#3 欠"task id 换成模型给的值时几枚尺会红"的现量；AC#4 欠 `ci.yml` junction 那步的"改前两发现读两发"；AC#5 欠与 `portable-tests.sh` four numbers 互咬的名册口径；AC#6 欠 tracked 分母（〔契约邻接〕，三支未预选）。本票**不加 `-done`**。
- 2026-10-06 20:0x 只读普查腿 `236-a4`（票 236 **AC#3 第一半**＝"先量清今天有几枚尺会红"，⛔ 零 `go` 命令／零产码／本腿一枚框未碰，起手与交件 `git status --porcelain -- cmd internal` 各 0 行）交回＝件 `.scratch/wisp/probes/236/a4/census.md`（385 行／35,757 字节，四节全实、占位符尺 0 命中；原始读数 77 份落 `a4/logs/**`，只建不删）。**读数三条**：①★**编排者那句『0 枚』不成立，真值＝1 枚，尺是 `internal/agent/compress_trace_test.go:526`**（`if len(task) != 36 || strings.Count(task, "-") != 4`，用例 `TestCompressionTraceCarriesTheOwningTaskID`／函数 `:478`；机制＝`task` 从真 `h.run()` 的 `Result.TaskID` 读回，换成 golden 那批模型 id（现量最短 1 字符、`msg_mock1`/`call_e1`/`r`）则 `len != 36` 成立 ⇒ `t.Errorf` ⇒ 红）；★**漏因＝AC#3 让筛的两目录是 `internal/tools`＋`cmd/wisp`，这一枚射程在 `internal/agent` ⇒ 按包开尺看不见跨包那一枚**（与"读盘尺对 `-overlay` 结构性瞎"同族）。②**但"0 枚"在票面两目录的口径内成立**，且★**这 1 枚红买不到 AC#3 要的东西**：类别 B 来源钉**四把负向尺全 0 命中**（TaskID↔CallID 互比／熵与 guessab 词面／`UnixNano` 4 命中全在 `internal/proc` mutex／`rand.` 在 `internal/agent` 测试面＝退化支无注入接缝、不可达），类别 C 字面量钉 27 行／48 串／19 枚文件**判会红 0 枚**（唯二真从生产铸造流读回来的 `cmd/wisp/subagent_carrier_197_test.go:203` 与 `cmd/wisp/run_test.go:138` 都只做自洽比对；`subagent_stream_key_197_test.go:55` 那枚 `"00000000-0000-4000-8000-000000000000"` 方向相反、今天**在给"模型给一枚 uuid 形 id"开绿灯）；分母尺逐字在件：116 枚测试文件／36 枚含 task-id 词（票面那句对上）／**203 行命中里只有 18 行带断言动词**。③**铸造点只有两枚、没有第三枚**（`git grep -n "newTaskID(" -- internal cmd` 锚带括号形状＝tracked 全量 6 行：`:329`／`:340` 两枚调用＋测试 `:557` 一行两次＋**两行注释**（`pointer_183_cli_seam_test.go:59`、`task_backfill.go:38`）；`loop.go:333` 是**传**不是**铸**），熵支退化形逐字＝`loop.go:1132-1137`（★格式化在 `if` 之外 ⇒ **退化形与正常形同 36 字符形状**，形状尺分不出；"唯一"站得住、"不可猜"站不住——⛔ 读码推的，没真跑）。**另两格现量**：D34 表（`docs/PLAN.md:2516` 起、表本体 `:2529`、task 两行 `:2564`／`:2565`）里**没有任何一枚工具的动作是铸造 id** ⇒ 钉这一格**不需要新工具面、不触票面禁区**（★顺带登记一个事实：表里**没有 `task.spawn` 行**，`grep -nE "task\.spawn" docs/PLAN.md` ＝ 0 命中，本腿不判它算不算缺陷）；`taskOutput.Execute`（`internal/tools/task.go:503-560` 全文读过）**不比对 caller** ⇒ "同进程内点名读走别的任务的输出"静态可达且无尺（⛔ 不写成漏洞结论；`pointer_185_cli_seam_test.go:327` T5 钉的是 C25 作用域借道、不是这道门）。**§4 写满**：静态尺结构性看不见五类／要真跑才知道的三件（含 A-1 定向突变最小形＝只换 `:329`＋`:340`）／`TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3` 的牙**分两档答**（本机 Windows 跑得到，`:404` 真 `mklink /J` 不 skip；CI ubuntu 进得了 core tier 作用域但 `:408 t.Skipf`，且★HEAD 现量 `git show HEAD:scripts/portable-tests.sh | grep -c "NotRelayTheFiledPath"` ＝ **0** ⇒ 不在 ledger 名册 ⇒ archived `run-37166458550-failed.log` 里 `--- SKIP`×2 ＋ `portable-tests.sh: unaccounted SKIP lines` ＋ `strict runner exited 1`，`runtests.sh` 现量 `top-level: PASS=1010 FAIL=4 SKIP=1`；windows tier `:249-253` **不含 `internal/tools`**）⇒ 今天 CI 跑不到它的执行路径、只跑得到"它 skip 了"这一声；⛔ 本腿不裁该进 ledger 还是改平台 tag（票 111 AC#10／票 236 AC#4 的地界）。**"钉该落在哪一层"本腿不裁**，四条理由写满（落点判据必须是真读数／三层候选天花板互不覆盖／票面 `:28` 自己定的顺序是"先量清…再定落点"／仓里已有同族成熟形状可抄＝`internal/session` minted-shape 锁三段双向，`grants_test.go:152`＋`session.go:67`/`:94`＋`ticket224_assembly_test.go:433`，且它不靠新工具面）⇒ **AC#3 的框要由另一个 agent 凭真读数翻**；本腿零产码改动 ⇒ d22scan／gofmt／gofumpt／go vet **一律没跑**（跑它们要 `go`），件与 logs 全是 Markdown/txt、不在 `.scratch/` 根下落文件。
