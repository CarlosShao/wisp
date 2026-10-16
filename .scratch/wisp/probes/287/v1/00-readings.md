# 287-v1 读数件（验收腿·非实现者·只读产码）

时钟：2026-10-15 23:2x +08 起手｜派单＝`.scratch/wisp/probes/orch/20261015/txt/287v1-dispatch.txt`＋同目录 `-send.txt`（冲突以后者为准）
HEAD 现量＝`dce9eb90d56608340f72ac8a6820bb854db2f8dc`（23:16:22，＝派单与释放记录那一发本身）
⚠ 具名一笔过期：派单第 2 行写的"当前 HEAD＝`988ce413`"在我起手时已经不是 HEAD（编队随后落了 `dce9eb90`）⇒ 本件所有读数都在 `dce9eb90` 上量。
裁决锚两枚都在盘：`git log -1 e1973191`＝10-15 22:09「287-r1 落地腿交件」、`e64294a0`＝22:11「追加一笔台面登记（⛔ amend）」。
面空：`git status --porcelain -- internal/tools`＝**0 行**（量，射程＝该目录，工作树）。

本件只给读数与尺。判语在 `docs/evidence/s1/287-v1-accept.md`。

---

## 1. 门禁五件：我自己跑的（⛔ 引腿的 rc）

| 尺 | 命令（全部仓根） | rc | 读数（写全集＋单位） |
|---|---|---|---|
| 编译 | `GOFLAGS= go build ./...` | **0** | 全树干净。**⛔ 本腿这一发没有 `internal/winsec` 的红可另节**＝在飞腿 `322-r1` 在 23:2x 这两枚件（`winsec.go`／`winsec_other.go`／`winsec_windows.go` 三枚 M）此刻仍能编译；我按铁律没读它们的内容当凭据 |
| 包级 | `go test ./internal/tools/ -count=1 -v` | **0** | **顶层用例 242 PASS／0 FAIL／0 SKIP**；第二把尺同日志 `^=== RUN`＝**325 枚**、`^    --- PASS`（子测试级）＝**83 枚** ⇒ 三把尺同日志、三个数，单位不同都自洽 |
| 静态 | `go vet ./internal/tools/` | **0** | 空输出 |
| D22 | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`；分母（ban #8 射程）`internal/` **571 枚 Go files**、`cmd/` **138**；bans #1-5 那支 `internal/`=239、`cmd/`=41（⚠ 两支射程不同，见 §2） |
| 长度门 | `sh scripts/check-path-length-budget.sh --with-self-test` | **0** | `VERDICT GREEN`、`tracked paths=13455`、`over-budget=60`、`roster=60`、控制组 2/3 ok（⚠ 必须在仓根跑，否则 rc=2 是射程拒绝） |
| 格式 | `gofmt -l internal/tools/subagent_depth_fixture_287_test.go internal/tools/subagent_197.go` | 0 | **列出 0 枚件名**（退码恒 0＝空尺，只能这样写） |
| 格式（扩射程对照） | `gofmt -l internal/tools`（整目录，量） | 0 | **列出 1 枚**：`internal/tools/task_output_pointer_notice_errtext_174r7_test.go` ⇒ 注释对齐差（`gofmt -d` 前 12 行＝同一段 `[]string` 的字面量对齐）。来路＝`67a7ab41`（2026-10-14，票 174-r7）**存量债**，⛔ 不在票 287 写面（`e1973191` 的父项清单没有它，§7 那把尺证） |
| 格式（门禁名义那把） | `gofumpt -l` 同上两枚 | **127** | `command not found`＝台面没装 ⇒ 本腿**只能**用 `gofmt -l` 交这一格，且具名它比门禁那把窄；⛔ 冒充"门禁那把绿" |

台面故障归类（⛔ 计入判据）：`gofumpt` rc=127 一枚；`grep -c "$(printf '\r')"` 这把尺**本腿自己踩废过一次**（见 §4 末）。
红句里 `walk`／`build failed`／`no tests to run`／`connection refused`／`0xc0000135` 命中数：我的包级日志 **0**（量，`grep -ciE`，射程＝`v1-01-test-head-full.log`）。

---

## 2. A996 欠着的那一把（d22scan 一度 rc=1 归属）

- 我现跑：`sh scripts/d22scan.sh` **rc=0**、`clean`；`internal/` 分母 **571**（腿交件时 568 ⇒ +3 枚是别腿上涨，⛔ 当红、⛔ 归我）。
- 射程具名（派单点的那一条，本腿复量）：bans #1-5 那一支 **不含** `_test.go`（`internal/`=239），ban #8 那一支**才含**（`internal/`=571，日志逐字写着 `comments and _test.go included`）⇒ 那枚 `⊂`（U+2282）属 **#8 那一支**，与 #1-5 不是同一条禁。
- **盘上痕迹（不是推测）**：`.scratch/wisp/probes/287/r1/stale-pre-trim-287r1-M3.log:21` 与 `:22` 两行逐字带着
  `subagent_depth_fixture_287_test.go:249: ... "拒绝派生：深度上限是 " ⊂ text=...`。
  ⇒ 那一发的 `⊂` 出现在**测试件 `t.Errorf` 文案里**（`⊂` 是印出来的判据文案，不是注释），且点位是 **:249**（旧长版夹具）；
  当前夹具（260 行）里 `⊂` 命中 **0**（量，`grep -c '⊂'`；正控＝同一把尺对同文件找 `。` 命中 **11** ⇒ 不是空尺），
  那一枚反形判据今天落在 **:211-213**，文案已改成 `出现在 text=%q`。
- **判语（就"归属"这一问）＝腿中途自己改了它自己的测试件**：既不是存量债（`⊂` 从没进过任何产码/别人件），也不是编排者修的（`e64294a0` 只动两枚 `.md`，§7 那把尺证）。腿的 `287r1-11-d22scan.log` 尾部逐字＝`clean ... internal/=568` ⇒ 它交件那一发确实绿；两发的差别只可能是盘上那枚测试件被它自己改过。

---

## 3. `AC#2` 三问（逐问给读数，⛔ 概括成"夹具造好了"）

夹具＝`internal/tools/subagent_depth_fixture_287_test.go`（260 行、3 枚用例；`wc -l` 量，工作树）。三枚用例在我自己的 HEAD 整包跑里逐枚 `--- PASS`（§1）。

### ① 真的走到深度判定那一支了吗？（攻恒真）

**走到了**，凭三发读数，不是凭句子：
- 主用例绿着的时候，回执正文是产码 `subagent_197.go:278-279` 那句**带真 uuid 的**成品（`拒绝派生：深度上限是 1，而派生者 <uuid> 自己就是一枚子代理（父 kind=subagent）。子代理不许再派子代理。`），这句今天在**任何产码路径**上说不出（§6 那把计数尺）。
- M1（`false && ok && rec.Kind == TaskKindSubagent`）下：主用例 **7 发**红（`:169/:173/:179/:185/:189/:192/:195`）＋次序尺 **2 发**红（`:255/:258`），红句里**第二代真被造出来**（`:189` 逐字 `名册里出现第二代（1 枚 [df68984f-542b-46f4-85ea-e3f0537adb3b]）＝深度上限 1 被走到并破掉`）⇒ 后果那一维有世界、也被读到。
- M3（那一支放宽成 `ok`，即"有名册行就拒"）下：**只有**反形正控用例红（`:213`×2、`:217`、`:221`）。⇒ **"判据换成反形必须红"这一问由 M3 答**，⛔ 主用例单独读：主用例在 M3 下**照绿**（它只看"有没有拒"，看不见"拒得宽不宽"）。这条是本格最重要的使用限制，写进验收表。

### ② `corr != taskID` 这一维是量出来的还是注释里写的？

**量出来的**，两枚独立来源：
- 测试件侧：`corrOf287`（夹具 `:124-143`）读 `x.mem.ListToolCallsByTask(ctx, taskID)` 的**持久化 `tool_call` 行**的 `CorrelationID` 字段，判三件：非空、`strings.HasPrefix(corr, taskID)`、`corr != taskID`（`:140`）；并要求该 tool 的行数 **恰 1 枚**（`:136-138`）。⛔ 读 ctx、⛔ 读注释。
- 产码侧（我现读，射程＝两枚件）：`internal/agent/loop.go:603 callCorr` 铸 `taskID + "#" + callID`（`callID==""` 时 `taskID#call-<i>`）⇒ 前缀成立且**结构上不可能等于** taskID；`internal/tools/bridge.go:1205` 写 `CorrelationID: orDefault(req.CorrelationID, req.TaskID)` ⇒ **直调桥那一形会把 corr 塌到 taskID 上**，正是 `:140` 那一支要拦的世界。
- ⚠ 诚实记账：那个"直调桥"的坏世界**没有台件把它造出来过**（5 枚突变全在产码侧，没有一枚动夹具）。所以本格的说法只能是"判据对着一个产码里真实存在的塌形"，⛔ 写成"已被突变证过"。

### ③ 那枚 `Kind` 是夹具自己塞的字段还是真装配给的？

**产码自己的 API、产码自己的发布点形状**，不是夹具自造的字段：
- `TaskRoster.PublishSubagent` 定义在 `internal/tools/task.go:342`，体内 `Record(taskID, TaskOutput{ParentTaskID:…, Kind: TaskKindSubagent, State:…})` ⇒ `Kind` 由**产码**写进名册，夹具只是调用方。
- 产码自己的两处调用＝`internal/tools/subagent_197.go:333`（孩子的 `opt.AdmitTask` 钩子里）与 `:368`（派生 goroutine 首落）。夹具在 `agent.Loop` 的 `AdmitTask` 钩子里调它（夹具 `:69-79`）＝**同一枚钩子形状、同一枚 API**，⛔ 自造结构体、⛔ 手改字段。
- 反形那一支用 `MarkRoot`（`task.go:355`）⇒ 两形都经产码的落库路径。
- 名册同一性：夹具往 `x.h.roster` 落行、判定读的是桥 `t.d.Roster`——若是两枚对象，判定就不会拒、主用例必红；今天主用例绿 ⇒ 同一枚（这是构造性反证，⛔ 我另造台件）。
- 借用的件全部**预先存在**（我逐枚 `grep -rl` 证，射程＝`internal/tools`、剔掉夹具自己）：`build221`/`harness221`/`windowGate221`＝`task_cancel_221_legs_test.go`；`fake197Provider`/`askSpawnFirstTurn`/`parent197`＝`subagent_197_test.go` ⇒ 夹具没顺手改别的测试件。

### ④ 同形判据点名（派单点的 `winsec_other.go` 那一坑）

夹具内**没有**"两对同形判据互冒充"的形状：同形判据是**跨用例**摆的，各自的世界不同——
- `HasPrefix(log.Text, refuse287HeadFmt)` 出现在 `:168`（主用例，池空）与 `:254`（次序尺，池满 4 枚）⇒ 同一把尺、两个世界，⛔ 互冒充（M1 下两枚都红，M3 下两枚都绿，行为一致＝它们判的是同一件事）。
- `InFlightSubagents()` 出现在 `:197`（主用例，want 0）与 `:250`（次序尺，want 4，且是 `t.Fatalf` 正控分母）⇒ 一枚读"拒绝不许占位"，一枚读"分母是我声称的那个"。
- ⇒ 真正的**两枚无突变落点的判据**（具名，⛔ 追认）：
  - `:197-199`（深度拒绝先于占位）：腿**自己已具名未证**（证据件 §3 末行）⇒ 本腿复核一致：M1/M2 下它不红，5 枚突变没有一枚动 `:276`↔`:282` 的次序。
  - `:180-182`（`who[1] != res.TaskID`，即"两枚 id 同屏"那一发）：**腿没具名**。它要红，需要"分支仍按 taskID 判、但 Sprintf 里把 `parentID` 换成 corr"那种改法——**5 枚突变里没有这一枚**（M2 动的是判定的读、不是句子的写），所以主用例在 M2 下是从 `:169/:173/:179/:185/:195` 红的，⛔ 从 `:181`。⇒ 这一发今天**无人钉**；本格其余部分不受影响（身份那一维由 M2 从"整句消失"那一侧钉住了）。

---

## 4. 突变台件算不算凭据（派单交我判的那一格）

`e1973191` 交回的突变目录现量＝**5 枚 `.go`**（`mutants/`：`B1-name-gate-off.go`／`B2-name-gate-and-listing-off.go`／`M1-disable-branch.go`／`M2-identity-to-corr.go`／`M3-overbroad-any-row.go`）**＋ 子目录 `mutants/stale-crlf/` 内 3 枚同名 `.go`**＋顶层 **6 枚 `stale-pre-trim-*.log`**＋**1 枚 `subagent_197.go.pre-mutation.bak`**（全部 `ls` 量，非算）。

- **"5 枚有没有把 CRLF 那一套另算一次"这一问的答案＝没有另算，且腿在 `e1973191` 正文里逐字点了它**：交回正文末段写着「留档⛔删：`stale-pre-trim-*.log`（夹具压到 260 行前的读数，行号已过期）与 `mutants/stale-crlf/*.go`」⇒ 5 枚＝**凭据名册**，`stale-crlf/` 与 `stale-pre-trim-*`＝**自报的过期留档**，两本账没混。本腿据这两条判：**5 枚算凭据、stale 两套不算**。
- 字节同一性（我现量，`sha256sum` 三处形态）：`git show HEAD:internal/tools/subagent_197.go`＝工作树＝`.bak` 三者**同串** `86de44c4b79b35172c5ca4bd42862f3659e1d9268ab992f975c282eaa6b172e3` ⇒ `.bak` 是今天的合法基线，腿的"三发还原后与基线同串"这一句**成立**。
- 落地尺寸（我现量，尺＝`diff -u <BAK> <mutant>` 后 `grep -c '^@@'` / `grep -c '^+[^+]'`）：M1/M2/M3 **1 hunk／1 改行**；B1 1/1；B2 2/2（设计如此＝闸＋目录两处）。⇒ 与腿自报逐枚对上。
- 锚落点（派单问①：⛔ 行号锚，认内容短语）：台件 `run-mutations.sh:36-47` 三枚锚都是**逐字节内容短语**，且 `:50-53` 要求 `anchor_hits==1` 否则 `refuse to seed`、`:55-57` 另判 `LANDS`。我在**今天的 HEAD 件上**独立复量唯一性（尺＝`grep -c`，射程＝单枚件）：`if rec, ok := t.d.Roster.Look(parentID); ok && rec.Kind == TaskKindSubagent {`＝**1**、`parentID := TaskID(ctx)`＝**1**、`if req.Name == subagentHiddenTool {`＝**1** ⇒ 锚没腐烂。
  三枚锚**全部落在函数体内**：`subagentSpawn.Execute`（M1/M2/M3）与 `subagentToolProvider.Execute`／`.Tools`（B1/B2），⛔ 包级 var／⛔ 注释。
- 还原凭据：`287r1-20-mutations-run-all.txt` 逐发带 `BASELINE-SHA`／`WROTE … anchor_hits=1 lands=…`／`hunks=1 改行数=1`／`RESTORED-EQUAL-OK`（3 发齐）／每发后 `git status --porcelain -- internal` 空行 ⇒ 派单要的"四件套"在盘上齐。
- **`stale-crlf/` 的归属判语（这一格本腿自己量出来的，⛔ 抄腿）**：尺＝`tr -cd '\r' | wc -c`（射程＝逐枚件字节）：三枚 stale 各 **CR=612**、三枚现行各 **CR=0**、`.bak` CR=0；把它们 `tr -d '\r'` 后与现行孪生件 `cmp` ⇒ **逐枚全同**（内容只差行尾，⛔ 是第二套突变）。对基线 `diff` 的假账也被我复现：**hunks 仍是 1**（整枚文件行行不同＝挤成一坨），但 `改行数` 从 1 涨到 **570** ⇒ 腿自抓 #3 那句"落地尺寸失真"**成立**，且比它自报的更具体：**`hunks` 这把尺对 CRLF 不敏感、只有 `改行数` 那把看得见**（后续复跑的人别只报 hunks）。
- ⚠ **本腿踩废的一把尺（具名＋正控）**：我先前用 `grep -c "$(printf '\r')"` 对同样 11 枚件量，**逐枚得 0**（＝空尺，命令替换把 CR 吞掉的形状）；换 `tr -cd '\r' | wc -c` 才量出 612。⇒ 本件所有 CRLF 数一律以 `tr` 那把为准；`grep` 那一发的"0 命中"**作废**，⛔ 任何人拿它当"stale-crlf 目录已无 CR"的读数。
- 6 枚 `stale-pre-trim-*`：逐枚点位与现夹具行号**不同**（例：M3 旧版红在 `:249`，现版红在 `:213`）⇒ 判为**过期留档**，不算凭据；只有 `stale-pre-trim-287r1-M3.log` 在本腿账上另有用途＝它是 §2 那枚 `⊂` 的**盘上唯一物证**。
- `stale-pre-trim-287r1-M3-full.log`（987 行）尾部逐字停在 `=== RUN Test197RowExistsBeforeFirstChildModelCall` 且无包级收尾行 ⇒ 腿自抓 #2（跑到既有仪器挂死被掐）属实：归类面残留，⛔ 计入任何判据。

---

## 5. `AC#3`：甲形那句"代价"说实话没有（⛔ 重裁这一格）

编排者已裁＝维持甲（口令「287 改裁乙」在案）。本腿只做"诚实审计"，两问逐问量：

**问 A：拒句可不可见？** ⇒ **可见，且今天已有仪器钉着**（这条腿的甲形代价格没写）：
- 产码 `internal/tools/subagent_197.go:604-610`（`subagentToolProvider.Execute`）返回 `agent.ToolOutcome{Text: "拒绝执行：子代理不能再派生子代理（深度上限 1），task.spawn 也不在你的工具目录里。", IsError: true}` ⇒ 走工具回执通道回给模型。
- 既有钉＝`internal/tools/subagent_197_test.go:663`（`if !out.IsError || !strings.Contains(out.Text, "深度")`）＋`:648`（孩子目录里不许有 `task.spawn`）＋`:652`（目录不许少别的）＋`:667`（被拒那枚不许打到父目录）。⇒ "硬拒"这一维**不是哑的**，⛔ "够不到＝没人知道"那一形。

**问 B：孩子能不能间接拿到 `task.spawn`？** ⇒ **今天不能**，且这把尺我给得出来（腿没给）：
- `subagentHiddenTool`＝`"task.spawn"`（`subagent_197.go:583`，唯一常量），全仓非测试产码里点名 `"task.spawn"` 的只有 `:187`（`subagentSpawn.Name()`）与 `:583` 两处（量，`grep -rn '"task\.spawn"' internal/ cmd/ --include=*.go | grep -v _test`）⇒ 没有第二枚工具顶着这个名字。
- 包层只有一处装配：`newSubagentToolProvider(` 产码命中 **1 枚**＝`:316`（另 1 枚＝`subagent_197_test.go:641` 测试直调）。
- 产码里起环路的装配只有 **2 处**：`internal/tools/subagent_197.go:342`（孩子，`opt.Tools` 已在 `:316` 被包）与 `cmd/wisp/run.go:1084`（宿主根，它的调用者行不在名册里当 subagent）。
- 名册里 `Kind=subagent` 的行只有 **2 个生产者**：`:333` 与 `:368`，两者都由 `subagentSpawn.Execute` 这一路发出，而这一路自己进 `:276` 之前必须先过 `:316` 的装配。
⇒ **"生产链今天够不到 `:276`"这一句是真的**，且它靠的是**枚数**（1 处包层／2 处装配／2 处发布／1 个名字常量），⛔ 靠"设计如此"那句话。

**审计结论（写成判语）**：甲形代价格里那三句（"动 0 枚产码"、"不会放开深度上限（`:91` 仍 1、`:276` 仍拒）"、"第二道的红已由夹具给"）**逐句未被任何一发读数证伪**——`:91 = MaxSubagentDepth = 1` 我现读到；"第二道现在有仪器"由 §3 的 M1/M3 两发撑住。
**但它确实没兑现"最坏后果是什么形状"这一维**（A/B 两问都要我自己去量），并且把端到端那一半挂成〔待面空补验〕。⛔ 我造出"那句写得虚"的反证（§5 的 A/B 反而把它补成了实话）⇒ 判语走"成立"，缺口以缺陷具名，⛔ 写"附条件"。

乙形那一栏我复核过的部分：`run-form-b-census.sh` 两发的整包包级日志（`287r1-B1-*.log`／`B2-*.log`）我现量＝**各 241 顶层用例 PASS／1 FAIL**，唯一红的既有函数是 `Test197ChildCannotDeriveSubagent`，**票 287 三枚用例在两形下逐枚照绿**（我在两份 160 KB 日志里逐枚 `^--- PASS: Test287` 量到 6 行）；两枚日志里 `panic` 字样只命中 1 枚**子测试名**（`TestTicket224…/panics_mid-lookup`，PASS），无 `build failed`／无 `test timed out` ⇒ 不是台面故障件。

---

## 6. `AC#4` 越界扫（我独立跑，⛔ 拿腿的自述当尺）

尺＝`git show --name-only --format=` 对 `e1973191`＋`e64294a0` 两枚父项清单求差集/求命中（射程＝HEAD 的 commit 树，⛔ 工作树）。

- `e1973191` 父项 **35 枚**，剔掉三处合法落点（`.scratch/wisp/probes/287/r1/**` 27 枚、`docs/evidence/s1/287-depth-fixture.md` 1 枚、票 287 1 枚）后**只剩 1 枚**：`internal/tools/subagent_depth_fixture_287_test.go` ⇒ 写面＝票面允许的 `internal/tools/**`，**越界 0 枚**。
- 禁区逐枚求命中（两枚提交合并清单，`grep -F -c`）：`internal/agent/loop.go` **0**／`internal/approval/` **0**／`internal/risk/` **0**／`internal/panel/` **0**／`internal/perm/` **0**／`internal/observe/thresholds.go` **0**／三枚冻结件（`tokens_fourway_test.go`、`l2_grant_boundary_test.go`、`ticket90_persist_test.go`）**各 0**／`frontend/` **0**／`design/` **0**／`.github/` **0**／`ci.yml` **0**／`docs/PLAN.md` **0**／`docs/specs/` **0**／`allowlist.txt` **0**／`cmd/wisp` **0**。⇒ `callCorr`（`internal/agent/loop.go:603`）零写＝由"件没进清单"这一条证，⛔ 我再读它在飞内容。
- `e64294a0` 父项 **2 枚**＝`probes/287/r1/00-readings.md`＋`docs/evidence/s1/287-depth-fixture.md` ⇒ ⛔ amend 另发一发的规矩守住了，两枚都是 `.md`，⛔ 碰任何 `.go`。
- 票面三框状态：`AC#2`／`AC#3`／`AC#4` 今天仍逐字 `- [ ]`（我读票 17-19 行原文）；`e1973191` 对票面的 numstat＝**2 增 0 删**，两行增量＝1 枚空行＋1 条 bullet ⇒ **腿没翻任何框**（`AC#1` 那枚 `[x]` 是编排者 10-09 翻的，在 `18c9c794`，⛔ 本腿重开）。
- `git status --porcelain -- internal/tools`＝**0 行**（起手 §0 复量）⇒ 287-r1 已收，写面空。

---

## 7. 没裁完／本腿不能裁（具名，⛔ blanks）

1. **`cmd/wisp` 端的端到端"不可达"台件**——本腿**没有**这一发读数，且派单 ⛔ 我跑（票 314 死亡区）。⇒ 这一格本腿留"证据不足"，缺的那一发＝`cmd/wisp` 面上"孩子的 `task.spawn` 在真装配里到底走不到 `:276`"的端到端台件（腿自己也标了〔待面空补验〕）。它⛔ 动摇 §5-B 那把枚数尺（枚数尺是全集计数，射程＝`internal/`＋`cmd/` 全部非测试 `.go`），⛔ 拿它当"已证"。
2. **`:180-182` 那一发无人钉**（§3-④）：缺的那一发＝"分支按 taskID 判、句子写 corr"那枚第 6 突变。本腿⛔ 自己造（改产码＝禁区）。
3. **`:197-199` 那一发无突变落点**：腿已自报未证；本腿复核一致，⛔ 追认。缺的那一发＝"`:276` 与 `:282` 换次序"那枚突变。
4. **`gofumpt` 那一把**：台面 rc=127 未装 ⇒ 本格由 `gofmt -l`（更窄）代交，已具名；要按门禁名义那把绿，缺的那一发＝一台装了 `gofumpt` 的机器。
5. `internal/tools` 整目录 `gofmt -l` 那 1 枚存量件（§1）：⛔ 本腿修（不在我写面、不在票 287 写面），⛔ 记到 287-r1 头上（来路＝`67a7ab41`）。
