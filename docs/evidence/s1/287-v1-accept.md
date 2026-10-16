# 票 287 `AC#2`／`AC#3`／`AC#4` 验收判语表（腿 `287-v1b`，2026-10-16）

**这一枚文件是什么**：`AGENTS.md` §1.5 要求的裁决表落点（`docs/evidence/s1/`）。前任验收腿 `287-v1` 的读数件
`.scratch/wisp/probes/287/v1/00-readings.md:9` 写着"判语在 `docs/evidence/s1/287-v1-accept.md`"——那一枚文件盘上**不存在**
（本腿起手现量：`test -f`＝NO）。按派单规矩，件里指名的文件一律当**待验断言**、交付判据是 `wc -c` 现量 ⇒ 本件就是补上的那一枚。
**本件⛔ 任何 AC 框**（翻框归编排者）、⛔ 改产码或测试件一字（本腿写面只有三处，见 §0-3）。

## 0. 锚、角色、写面、台面起手

1. **裁决锚**：实现笔 `e1973191b839c49473543fca2ce304898a408434`（287-r1 落地，10-15 22:09）＋`e64294a0ed39e7279b262b6e17c472df20b20a99`
   （同腿按⛔ amend 另发的台面登记，10-15 22:11）。两枚都在盘（`git log -1` 现量）。
   实现面现量：`internal/tools/subagent_depth_fixture_287_test.go`＝**260 行／11,650 字节／3 枚用例／`t.Skip` 0 命中／sha256 `02cbaf10…`**；
   `docs/evidence/s1/287-depth-fixture.md`＝**13,631 字节**。
2. **读数锚**：本腿五件门禁跑在 `1a651618d6a5b71e06cb216f1a309304becfd8aa`（＝本腿那笔 `287-v1` 遗产代提），
   突变台件在 `git archive 1a651618` 的仓外导出树上跑。★可比性：`git diff --stat 46007856..1a651618 -- internal/tools`＝**空**、
   夹具那枚 blob 在 `46007856` 与本腿锚点**同哈希**（`a2681ead…`）⇒ 派单第 2 行写的"当前 HEAD＝`46007856`"起手已过期，⛔ 影响这一面。
   写本件时盘上 HEAD＝`da6fda0d37193907465152f3927f7bec20f9a52d`（别腿与编排者在飞），未推现量＝**55 笔**。
3. **本腿写面（只有这三处）**：`.scratch/wisp/probes/287/v1b/**`、本件、票 287 末尾一条 bullet。
   工作树归因：`git status --porcelain -- internal cmd tools` 全程只有**一枚别腿在飞件** `?? cmd/wisp/panel_l1_blocked_220_test.go`
   ⇒ 本腿对 `internal`／`cmd`／`tools` 零字节改动（突变一律在仓外树，跑完逐枚还原＋复跑绿，见 §3 末与 `logs/v1b-12-*`）。

### 判语总表（只有三种：成立／不成立／证据不足＋缺哪一发读数）

| 格 | 子问 | 判据（尺） | 本腿读数 | 判语 |
|---|---|---|---|---|
| **AC#2** | 深度那一支**真走到**了吗（恒真检查） | 仓外树物理摘掉 `:276` 那一支（⛔ `-overlay`），整包跑 | V1：主用例红 7 点＋次序尺红 2 点，红句里第二代真被造出来；整包 **240/2**，红名册只有这两枚 287 用例 | **成立** |
| AC#2 | 一、`Kind` 是夹具自填还是产码给的 | V7 改 `task.go:344` 的 `Kind`；＋静态尺数夹具自造枚数 | 主用例红 `:160`（`Kind="root"`）；夹具内 `TaskOutput{`／`.Record(`／`Kind:` 各 **0** 命中，只调 `PublishSubagent`／`MarkRoot` | **成立** |
| AC#2 | 二、`corr != taskID` 是量的还是注释 | V6 把 `loop.go:676` 的 corr 塌到 taskID | 红点**恰** `:176`／`:229`（`corrOf287` 读持久化 `tool_call` 行），红句逐字见 §3 | **成立** |
| AC#2 | 三、派发经**真 loop** | V8 令 `loop.go:374` 不再调 `AdmitTask` | 三枚用例全红在夹具自己的前置（`:157`／`:205`／`:252`，"AdmitTask 跑了 0 次"） | **成立** |
| AC#2 | 反形/恒真检测器（⛔ 判据换成反形也绿） | V3 把那一支放宽成"有名册行就拒" | **只有**反形正控红（`:213`×2／`:217`／`:221`），主用例与次序尺照绿 | **成立** |
| AC#2 | `:180-182`（两枚 id 同屏那一发）**有没有突变落点** | 本腿补造 V4：判定仍按 taskID、句子把派生者写成 `CorrelationID(ctx)` | 红点**恰只有 `:181`**，整包 **241/1**（唯一红＝主用例） | **成立**（前任那两枚"无人钉"之一由本腿补上） |
| AC#2 | `:197-199`（拒绝先于占位）**有没有突变落点** | 本腿补造 V5：深度块整体搬到 `TryAcquireSubagentSlot` 之后（hunks=2） | 红 `:198`（"还占着 1 枚池位"）＋次序尺 `:255`／`:258`；反形照绿 | **成立**（同上，第二枚） |
| AC#2 | 四件套（种前哈希／复量／红句逐字／还原哈希等值） | 八发逐枚 | 八枚 meta 齐：`PRE-SHA`＝`anchor_hits=1`＝`hunks/改行数`＝红句＝`RESTORED-EQUAL-OK`；还原后整包**复跑绿 242/0/0** | **成立** |
| AC#2 | 〔使用限制〕主用例单独读不出"拒得宽不宽" | V3 下主用例照绿（它只看"有没有拒"） | 见本表第 5 行 | 本格已用**反形**补这一维 ⇒ 判语只在"主用例＋反形＋次序尺"三枚合读时成立 |
| **AC#3** | 甲形那句"代价"说实话没有（⛔ 重裁形、⛔ 动产码） | 三句逐句用尺复量（枚数尺＋现读常量＋V1 的红名册） | 见 §4 问 A／问 B／三句审计 ⇒ 三句**逐句未被任何一发读数证伪** | **成立**（缺口以缺陷具名，见 §6-a） |
| **AC#4** | 写面只 `internal/tools/**`、⛔ 碰 `loop.go` 的 `callCorr`／深度那一支 | `git show --name-only --format=` 对两枚提交，20 枚禁区逐枚求命中 | `e1973191`＝35 枚（32 枚本票 probes＋1 枚证据件＋1 枚票面＋**1 枚新夹具**），剔四处合法落点后剩余**只有夹具**；禁区 20 枚命中**全 0** | **成立** |
| AC#4 | 门禁四件 rc | 本腿自己跑两发（起手一发＋写本件前一发） | `GOFLAGS= go build ./...`＝**0**、`go vet ./internal/tools/`＝**0**、`sh scripts/d22scan.sh`＝**0**（clean）、长度门 `--with-self-test`＝**0**（GREEN，三档正控全 ok） | **成立** |
| AC#4 | 改前改后三数并排＋红名册差集 | 本腿在仓外树**物理摘掉夹具**复量（⛔ 只引腿的自报） | 改前 **239 PASS／0 FAIL／0 SKIP**（RUN 322）→ 改后 **242／0／0**（RUN 325）；红名册差集＝**空**；与腿自报逐枚同数 | **成立** |
| AC#4 | `gofmt -l` 对动过的件空 | 尺＝`gofmt -l`（退码恒 0＝空尺，只数件名） | 夹具＋`subagent_197.go`＝**列出 0 枚**；整目录那 1 枚是 `task_output_pointer_notice_errtext_174r7_test.go`（来路 `67a7ab41`，⛔ 票 287 写面） | **成立** |
| AC#4 | `gofumpt -l` 对动过的件空 | 两把：`gofumpt -l`（走 PATH）＋`/d/work/base/gopath/bin/gofumpt.exe -l`（绝对路径） | PATH 那把 **rc=127 command not found**（287-v1 那发同形）；绝对路径那把**列出 0 枚**（两枚件），整目录列出 **1 枚**＝同一枚 174-r7 存量件 ⇒ 这把尺**非空尺**、正控成立 | **成立**（PATH 缺件属台面形状，⛔ 判据，见 §6-b） |
| AC#4 | `frontend/**`／`design/**` 零读零写零转述 | 本腿自查：只对两枚名取过 `du -sh`／`git ls-files \| wc -l`（枚数与字节） | 零内容读、零写、本件零转述（仪器自己的射程⛔ 算本腿读） | **成立** |
| AC#4 | ⛔ push、commit 必带显式 pathspec | 历史可证的只有**每笔落点名册**；命令形状不在盘上 | 未推现量 55 笔（`git rev-list --count origin/dev..HEAD`）；两枚实现笔的落点全在四处合法面 | **成立**（间接：无一枚越界件可作反证；⛔ 把"命令形状"当已证，见 §6-d） |

## 1. 门禁五件（本腿逐枚复跑，⛔ 引前任与腿的 rc）

尺与读数全在 `.scratch/wisp/probes/287/v1b/10-gates-rerun.md`＋`v1b-01..07` 的 `.log`/`.rc` 成对件＋`logs/v1b-11`／`12`。

- 编译 **rc=0**（0 字节）／包级 **rc=0**（顶层 `^--- PASS`＝**242**、FAIL 0、SKIP 0、`^=== RUN`＝**325**、子测试 `^    --- PASS`＝**83**）／
  静态 **rc=0**／D22 **rc=0**（clean；ban #8 `internal/`=**572**、`cmd/`=138；bans #1-5 `internal/`=239、`cmd/`=41；ban #7 `internal/tools/`=25）／长度门 **rc=0**（GREEN，`tracked paths=13688`／`over-budget=60`／`roster=60`）。
- 与 287-v1 那发**五把包级尺逐枚同数**（242/0/0、325、83）⇒ 前任 §1 这一格由本腿独立复跑撑住。
- 第二发门禁（写本件前，此刻 `cmd/wisp` 面有别腿在飞件）：`build rc=0`／`d22scan rc=0` ⇒ 这一格⛔ 依赖"起手那一刻的洁净"。
- 台面故障归类（⛔ 计入判据）：`gofumpt` 走 PATH＝rc=127（**§6-b 已由绝对路径那把改判，本格判语⛔ 依赖这一发**）；包级日志里 `panic|walk|build failed|no tests to run|0xc0000135|test timed out` 命中 2 行，
  逐枚读＝同一枚既有子测试名 `TestTicket224GrantSourceThatCannotAnswerMeansAsking/panics_mid-lookup`（RUN＋PASS）⇒ 零真红。

## 2. 仓外导出树是"同一块台面"吗（先立的尺，⛔ 拿突变读数当凭据）

- 树＝`git archive 1a651618 | tar -x -C /d/work/tmp/287v1b-mutation-bench/tree-C`；六件套 `go.mod go.sum internal cmd tools scripts` 逐枚验在；
  `.go` 枚数＝总 **1091**／`internal` **572**／`cmd` **139**／`tools` **18**（证非空，且 `internal` 那枚与本腿 d22scan 的 ban #8 分母 **572** 逐字相合）。
- 受种五枚件 CR 计数＝**0**（尺＝`tr -cd '\r' | wc -c`，⛔ 前任踩废过的 `grep -c "$(printf '\r')"`）；`subagent_197.go`＝`86de44c4b79b…`
  ＝仓内工作树＝287-r1/287-v1 两发的还原串（三处同串）。
- **树的基线整包**（`logs/v1b-10-*`）＝**242/0/0、RUN 325、子 PASS 83** ⇒ 与仓内工作树**五把尺逐枚同数**，突变读数可比。
- ⚠ 前两枚树的台面值只作**过程留档**：六件套树＝234/8、加 `docs/PLAN.md` 的树＝239/3，红因逐枚＝`TestTicket252*`（8.3 短名探针对着树自己的路径无从比较）、
  `Test229…`（要 walk `frontend/`）、`Test172…`（要读 `docs/specs/SPEC-07`）⇒ **树缺件＋路径形状**，⛔ 归 287、⛔ 归别腿；tree-C 用整档 archive 消掉。

## 3. `AC#2` 八发突变名册（改法＝逐字节内容锚，⛔ 行号锚；每发还原后哈希等值）

完整逐发读数＝`.scratch/wisp/probes/287/v1b/mutation-roster.txt`（由 `logs/v1b-V*.meta.txt` 现生成）＋`logs/v1b-V*-test.log`／`-test-full.log`。
锚唯一性本腿逐枚现量（`grep -c`，射程＝单枚件）：深度行＝1、`parentID := TaskID(ctx)`＝1、`:279` 那半句＝1、`loop.go:676` 那行＝1、
`task.go` 的 Kind 行＝1、`if l.opt.AdmitTask != nil {`＝1。

| 发 | 落点件 | 改法（逐字节） | hunks/改行 | 红的**指名用例**与红点 | 整包（顶层 PASS/FAIL） | 还原 |
|---|---|---|---|---|---|---|
| V1 | `subagent_197.go:276` | 谓词首项置 `false &&`（＝那一支摘掉） | 1／2 | 主用例 `:169 :173 :179 :185 :189 :192 :195`＋次序尺 `:255 :258`；反形绿 | **240/2**，红名册＝那两枚 287 用例（既有仪器零红） | `86de44c4…` 等值 |
| V2 | `:271` | `TaskID(ctx)`→`CorrelationID(ctx)` | 1／2 | 三枚 287 用例全红（`:169 :173 :179 :185 :195`／`:221`／`:255 :258`） | **238/4**，另红＝既有 `Test283IdentityChainThroughTheRealLoop` | 同上 |
| V3 | `:276` | 收成 `ok {`（有名册行就拒） | 1／2 | **只有**反形 `:213`×2 `:217` `:221`；主用例与次序尺照绿 | ⛔ **作废读数**：跑到既有 `Test197RowExistsBeforeFirstChildModelCall` 挂 9m42s、`panic: test timed out after 10m0s` | 同上 |
| **V4** | `:279` | 句子把 `parentID` 换成 `CorrelationID(ctx)`（判定⛔ 动） | 1／2 | **恰只有主用例 `:181`**；余两枚绿 | **241/1**（唯一红＝主用例） | 同上 |
| **V5** | `:276-280`↔`:282-288` | 深度块整体搬到 `TryAcquireSubagentSlot` 之后（MOVED depth=5／pool=7） | 2／10 | 主用例 `:198`＋次序尺 `:255 :258`；反形绿 | **240/2** | 同上 |
| V6 | `agent/loop.go:676` | `CorrelationID: callCorr(…)`→`CorrelationID: taskID` | 1／2 | 主用例 `:176`＋反形 `:229`（次序尺不调 `corrOf287` ⇒ 照绿） | **237/5**，另红＝`Test283…`／`Test285RosterRowsCarryDistinctCorrPerCall`／`TestLoopPassesDeclaredL1WriteThroughTheGate` | `57b734eb…` 等值 |
| V7 | `tools/task.go:344` | `PublishSubagent` 那行 `Kind`→`TaskKindRoot` | 1／2 | 主用例 `:160`（前置一）＋次序尺 `:255 :258`；反形绿 | **237/5**，另红＝`Test197SubagentPoolCapsAtBridgeCeiling`／`Test197FullPoolRefusesNextSpawnWithReadableReason`／`Test222WaitingParentHoldsNoBridgeSlot` | `91608259…` 等值 |
| V8 | `agent/loop.go:374` | `if l.opt.AdmitTask != nil {`→`if false && …` | 1／2 | 三枚全红在夹具自己的前置：`:157`／`:205`／`:252` | **235/7** | `57b734eb…` 等值 |

**红句逐字（本腿自己复跑出来的原文，含 `文件:行`；全在 `-test.log` 里）**

- V1 `subagent_depth_fixture_287_test.go:189`：`名册里出现第二代（1 枚 [50d6dabc-54db-4bf3-b025-d868011f6aa4]）＝深度上限 1 被走到并破掉；调用者=8bb5067d-925c-4fbe-816e-416f27b36408`
  ⇒ 后果那一维有世界、也被读到（⛔ 只有句子）。
- V3 `:213`（两发同点，`for` 里两枚判据各一发）：`父 kind=root 的调用者被深度判定拒了（那一支已经宽到认不出 kind）：refusal="拒绝派生：深度上限是 " 出现在 text="拒绝派生：深度上限是 1，而派生者 65b150dd-94db-4d61-ad1f-b3cf252931dc 自己就是一枚子代理（父 kind=root）。子代理不许再派子代理。"`
  ⇒ 恒真检测器由这一形打掉；同一发下主用例**照绿**（本格的使用限制）。
- **V4 `:181`**：`深度判定读的身份 = "68724e6d-e265-4bc9-b452-1dc18942e578#call-197-spawn", want task id "68724e6d-e265-4bc9-b452-1dc18942e578"；这一发的 per-call corr = "68724e6d-e265-4bc9-b452-1dc18942e578#call-197-spawn"（corr != taskID）：把 :271 改读 CorrelationID(ctx) 就是这个读数`
  ⇒ `:180-182` 那一发**有落点**（287-v1 §3-④ 判它"无人钉"，本腿补造成功；差因见 §6-e）。
- **V5 `:198`**：`深度拒绝发生之后还占着 1 枚池位：那一支必须在 TryAcquireSubagentSlot 之前返回`
  ⇒ `:197-199` 那一发**有落点**（同上，缺的那一发＝"次序换过来"，本腿落种）。
- V6 `:176`：`task.spawn 的行 corr = "c1c278b6-6b10-4709-a530-e2eab2081af9", want 非空、以 task id "c1c278b6-6b10-4709-a530-e2eab2081af9" 为前缀且不等于它（派发没经真 loop）`
  ⇒ 那一维**是量出来的**（读持久化 `tool_call` 行，⛔ 读 ctx、⛔ 读注释）。
- V7 `:160`：`前置一没落地：调用者 ef6947a5-d5a1-4fd1-a331-4cdea5f8c69c 的行 ok=true Kind="root", want "subagent"`
  ⇒ Kind 由产码 `PublishSubagent`（`task.go:342`）落进名册；夹具自造枚数＝0（静态尺）。
- V8 `:157`：`夹具自己没站住：AdmitTask 跑了 0 次，callerID="", res.TaskID="9033c446-59ee-4a58-9e83-8449378c3620"`
  ⇒ "经真 loop"这一前置由 loop 自己的准入钩子撑住（`subagent_197.go:333` 是产码同形调用点，`grep -n` 现量）。

**还原与复跑**：八发跑完，五枚受种件哈希逐枚＝起手基线，且与**仓内工作树**四枚同串；还原后整包复跑
（`logs/v1b-12-after-restore-green-full.log`）＝**242/0/0、RUN 325、子 PASS 83**、三枚 287 用例逐枚 PASS ⇒ ⛔ 残留。

**★独立复现（本腿自己算的，`23b-leg-mutants-reproduced.md`）**：本腿 V1／V2／V3 三发是**独立落种**（自己的锚、自己的仓外树），落种后字节与 287-r1 交回的
`mutants/M1-disable-branch.go`／`M2-identity-to-corr.go`／`M3-overbroad-any-row.go` **同串**（`69087120ea40…`／`7c7185cc95d9…`／`5aaba840f472…`，尺＝`sha256sum`），
`.bak`＝基线 `86de44c4b79b…`，`mutants/stale-crlf/` 三枚 CR＝**612**、现行五枚 CR＝**0**（＝287-v1 §4 那两把尺的形状，本腿复现）
⇒ 腿那三发的"种"是真的、且本腿逐发的红名册与腿自报的后果逐枚对上；这一条把 `AC#2` 的"凭据"那一问从"读腿的日志"升成"本腿复跑同串"。

## 4. `AC#3`：只裁"甲形那句代价说实话没有"（⛔ 重裁形、⛔ 动产码一字）

编排者已裁＝维持甲（口令「287 改裁乙」在案）。本腿做诚实审计，两问逐问用自己的尺量（`21-ac2-static-rulers.txt`）：

- **问 A：那道"硬拒"可见吗、有人钉吗** ⇒ 可见且已钉。产码 `subagent_197.go:603-612` 的 `subagentToolProvider.Execute` 在 `:604` 那枚名字闸返回
  `Text: "拒绝执行：子代理不能再派生子代理（深度上限 %d），task.spawn 也不在你的工具目录里。"`＋`IsError: true`＝走工具回执通道回给模型；
  既有仪器四点位本腿现读到逐字：`subagent_197_test.go:648`（目录里不许再有 `task.spawn`）／`:652`（目录不许少别的）／
  `:663`（"孩子派孙没有被拒并说明原因"）／`:667`（被拒那枚不许打到父目录）⇒ ⛔ "够不到＝没人知道"那一形。
- **问 B：孩子能不能绕过名字闸间接拿到 `task.spawn`** ⇒ 今天不能，凭**枚数**（⛔ 凭"设计如此"那句话）：
  `"task.spawn"` 字面量在非测试产码里＝**2 处**（`:187` 的 `Name()`、`:583` 的唯一常量）；`newSubagentToolProvider(` 产码＝**2 处**
  （定义 `:579`＋唯一装配点 `:316`，另 1 枚测试直调 `subagent_197_test.go:641`）；名册里 `Kind=subagent` 的发布生产者＝**2 处**
  （`PublishSubagent` 的 `:333`／`:368`）＋收尾那枚 `Record(... Kind: TaskKindSubagent)`＝`subagent_197.go:459`；
  `MarkRoot` 的产码调用者＝**1 处**（`cmd/wisp/run.go:1166`）⇒ 要摸到 `:276` 必须"调用者那行是 subagent"且"它的工具面未被包层"，
  而这两个条件在生产链上互斥（`:316` 先包层）。
- **甲形代价格三句逐句审计**：① "动 0 枚产码"＝§5 越界扫（写面只有那枚新测试件）⇒ 未证伪；② "不会放开深度上限（`:91` 仍 1、`:276` 仍拒）"
  ＝本腿现读 `:85 MaxConcurrentSubagents = 4`、`:91 MaxSubagentDepth = 1`，基线整包三枚 287 用例逐枚 PASS ⇒ 未证伪；
  ③ "第二道的红已由本票夹具给"＝V1（摘掉那一支 ⇒ 指名用例红、既有仪器零红）＋V3（放宽 ⇒ 只有反形红）⇒ 未证伪。
- **判语＝成立**。缺口（⛔ "附条件"，具名进 §6）：甲形那句确实没兑现"最坏后果是什么形状"，问 A／问 B 都要验收方自己去量；
  端到端那一半（`cmd/wisp` 面上"孩子的 `task.spawn` 到不到得了 `:276`"）腿自己标了〔待面空补验〕、本腿也⛔ 有这一发。

## 5. `AC#4` 边界与越界扫（尺＝`git show --name-only --format=`，射程＝两枚提交，⛔ 工作树）

- `e1973191`＝**35 枚**：`​.scratch/wisp/probes/287/r1/**` 32 枚＋`docs/evidence/s1/287-depth-fixture.md` 1 枚＋票 287 1 枚＋
  `internal/tools/subagent_depth_fixture_287_test.go` 1 枚 ⇒ 剔掉四处合法落点后**剩余只有那枚夹具** ⇒ 写面＝票面允许的 `internal/tools/**`，越界 **0**。
- 禁区 20 枚逐枚求命中（两笔合并清单，`grep -F -c`）＝**全 0**：`internal/agent/loop.go`／`internal/approval/`／`internal/risk/`／`internal/panel/`／
  `internal/perm/`／`internal/observe/thresholds.go`／三枚冻结件各 0／`frontend/`／`design/`／`.github/`／`ci.yml`／`docs/PLAN.md`／`docs/specs/`／
  `allowlist.txt`／`cmd/wisp`／`go.mod`／`go.sum`／`.gitignore`。⇒ `callCorr`（`internal/agent/loop.go:603`）零写由"件没进清单"证，⛔ 读它在飞内容当凭据。
- `e64294a0`＝2 枚，全 `.md`（`probes/287/r1/00-readings.md`＋证据件）⇒ ⛔ amend 另发一发的规矩守住、⛔ 碰任何 `.go`。
- 票面三框今天仍逐字 `- [ ]`（`:17`／`:18`／`:19` 原文已随读数落盘）；`e1973191` 对票面 numstat＝**2 增 0 删** ⇒ 腿⛔ 翻框。
  （`AC#1` 那枚 `[x]`＝编排者 10-09 在 `18c9c794` 翻的，⛔ 本腿重开。）
- 三数并排：本腿**物理摘件**复量改前＝239/0/0（RUN 322）⇒ 改后 242/0/0（RUN 325），红名册差集空；
  与腿自报（`287r1-05`／`287r1-10` 两枚日志本腿自己数＝239/0/0 与 242/0/0）逐枚同数。

## 6. 缺陷／未闭合清单（本腿射程内不能证的都写在这一节，⛔ blanks）

- **a｜端到端那一发没有**：`cmd/wisp` 面上"生产链今天仍不可达"的端到端台件＝零发（票 314 死亡区、且此刻 `cmd/wisp` 面有别腿在飞件
  `?? cmd/wisp/panel_l1_blocked_220_test.go`）。§4 问 B 那把枚数尺射程＝`internal/`＋`cmd/` 全部非测试 `.go`，⛔ 顶替端到端台件。
  **这一发⛔ 动摇 `AC#2`／`AC#3` 的判语**（它们只裁"夹具能不能红／代价说实话"），但它是〔待面空补验〕的账，仍挂在票面上。
- **b｜`gofumpt` 那一把：第一发判"证据不足"，第二发改判"成立"**（读数＝`22-gofumpt-absolute.txt`）：
  起手那发 `gofumpt -l` 走 PATH＝**rc=127 command not found**（287-v1 那发同形，本腿 `command -v gofumpt` 正控亦空）⇒ 当时只能由 `gofmt -l` 代交并具名它射程更窄。
  写本件时本腿找到并用了绝对路径那把：`/d/work/base/gopath/bin/gofumpt.exe -l` 对**两枚动过的件列出 0 枚**、对整目录列出 **1 枚**
  （＝同一枚 `task_output_pointer_notice_errtext_174r7_test.go` 存量件）⇒ 这把尺**量得到东西**、两枚件**合规**，`AC#4` 该子句判语＝**成立**。
  残余只是一条台面形状（⛔ 判据）：这台机器的 PATH 里没有 `gofumpt`，后续腿⛔ 再把它当"没装"而交空尺。
- **c｜`internal/tools` 整目录 `gofmt -l` 那 1 枚**（`task_output_pointer_notice_errtext_174r7_test.go`）：来路 `67a7ab41`（票 174-r7）存量债，
  ⛔ 本腿修、⛔ 记到 287-r1 头上。
- **d｜"commit 必带显式 pathspec"的命令形状⛔ 由历史复证**：盘上只有落点名册（可证的＝无一枚越界件可作反证）。⛔ 把它当已证。
- **e｜与 287-v1 那发不同的地方，逐枚具名**（判语表里只允许用本腿复出来的数）：
  1. **`:180-182`／`:197-199` 两枚"无人钉／无突变落点"**：本腿**造出了那两发**（V4 恰只有 `:181` 红、V5 恰红 `:198`＋次序尺）⇒ 与 287-v1 §3-④／§7 结论不同。
     成因＝前任⛔ 自己落种（它 §7 明写"改产码＝禁区"），只数了腿交回的 5 枚名册；本腿按派单把突变放到**仓外导出树**物理落种，禁区一字未动。
     ⚠ 这不是说 287-r1 交回的凭据有假：它的 5 枚名册与还原哈希都在盘上（本腿 §0-1 复量过 `86de44c4…` 三处同串），只是**名册没覆盖这两维**。
  2. **d22scan 的 ban #8 `internal/` 分母**：v1＝571，本腿＝**572**。成因＝两发锚点之间 `322-r1b` 把 `internal/winsec/openleaf_322_other_test.go` 拆成两枚
     （`git diff --stat 46007856..1a651618` 现量＝winsec 三枚＋d22scan 一枚，`internal/tools`＝空）⇒ 别腿上涨，⛔ 归 287、⛔ 当红。
  3. **长度门 `tracked paths`**：v1＝13455，本腿＝**13688**。成因＝同批发出件＋两笔遗产入库。`over-budget=60`／`roster=60` 两发**逐枚同数**。
  4. **长度门正控措辞**：v1 写"控制组 2/3 ok"，本腿日志逐字是三行 `control 1/3 ok`／`control 2/3 ok`／`control 3/3 ok`（＝三档全过）
     ⇒ 本腿日志**否证**"只过了 2/3"那一读，属前任措辞腐坏。
  5. **`bridge.go` 那枚 corr 塌形的点位**：v1 引 `bridge.go:1205`，本腿现量＝`:1206`（`CorrelationID: orDefault(req.CorrelationID, req.TaskID)`）；
     另 `:269-270` 也有一处 `req.CorrelationID = req.TaskID` 的塌形入口。⛔ 影响判语，记为行号腐烂一例。
  6. **本腿补的两枚读数**：V3 的整包档因既有仪器挂死被 10m 超时掐 ⇒ 那一发的整包数**作废**（与 287-r1 自抓 #2、287-v1 §4 末条同形）；
     以及 §2 前两枚树的 8/3 枚红＝树形状台面，⛔ 计入判据。
- **f｜前任那 13 枚遗产的去向**：已按枚入库（笔 `1a651618`，tracked 0→13），字节零改写（sha256 名册＝`00-legacy-integrity-before-add.txt`），⛔ 删、⛔ 改写；
  入库时 git 对 9 枚 `.log`/`.rc` 发了 `LF will be replaced by CRLF` 提示＝仓库 `* text=auto` 的入库归一化（⛔ 本腿编辑它们），工作树字节复核仍逐枚同串。
- **g｜本腿自抓（⛔ 落盘件那一发）**：写本件前那发门禁只跑了 `build`＋`d22scan` 两枚并**只把 rc 写进正文**（⛔ 像起手发那样落成 `.log`/`.rc` 成对件）
  ⇒ 那一发的可复核性弱于其它发；交件后那发（`logs/r2-01..06`）已把**同一格五件**补齐成对落盘，故⛔ 动摇任何判语，只记本腿这一条流程缺陷。
- **h｜本腿⛔ 复跑的东西一律写明**：`cmd/wisp` 整包（票 314 死亡区，且此刻该面有别腿在飞件）、`internal/panel`／`internal/perm`（共享编译面⛔ 本腿射程）、
  287-r1 那五枚交回突变的**原始运行**（本腿只复量它们的字节面：5 枚 `.go`＋`mutants/stale-crlf/` 3 枚＋`.bak` 一枚与基线同串 `86de44c4…`、CR 计数 612/0 与前任 §4 同形）
  ⇒ 这些都⛔ 进本件的判语列，只进"与前任比"的对照列。

## 7. 读数落点名册（哪枚文件是哪把尺；全部在 `.scratch/wisp/probes/287/v1b/`）

- 门禁五件＋格式两把：`gates.sh`／`v1b-01..06` 的 `.log`+`.rc`／`10-gates-rerun.md`；包级 -v 全日志＝`v1b-07-test-tools-full.log`（157,003 字节）＋`07-test-counts.txt`。
- 仓外树建立与基线：`bench-setup.sh`／`.tree-path`／`08-baseline-tree-counts.txt`／`10-baseline-treeC-counts.txt`／`logs/v1b-08|09|10-*`。
- 突变台件：`make-anchors.sh`＋`anchors/`（14 枚逐字节锚）／`mutate.py`／`mutate_v5.py`／`bench-run.sh`／`logs/v1b-V1..V8.meta.txt`＋同名 `-test.log`／`-test-full.log`／`make-roster.sh`＋`mutation-roster.txt`。
- 改前复量：`logs/v1b-11-before-287-absent-full.log`＋`11-before-after-counts.txt`；还原复跑绿：`logs/v1b-12-after-restore-green-full.log`＋`12-after-restore-readings.txt`。
- 只读尺三枚：`20-ac4-boundary-scan.txt`（越界扫＋票面三框原文）／`21-ac2-static-rulers.txt`（夹具自造枚数＋产码 Kind/装配/发布点位＋既有仪器四点位＋常量）／
  `22-gofumpt-absolute.txt`＋`.gofumpt-two.log`＋`.gofumpt-dir.log`（§6-b 那一发改判所用的绝对路径那把尺）。
- 交件后台面复跑：`logs/r2-01..06` 的 `.log`+`.rc`（五件门禁第二落盘发＋`gofumpt` 绝对路径 0 枚列出）。

## 8. 本腿自证三把尺

1. 提交号（⛔ 把笔数写死在本件里——那一行会随本件自己过期；**尺**＝`git log --grep=287-v1b --format=%h` 现量。全部⛔ push）：
   已落的前六笔＝`1a651618`（287-v1 遗产 13 枚按枚入库）／`96f3767d`（门禁五件复跑，18 枚件）／`7dec3a6c`（突变台件八发＋只读尺，56 枚件）／
   `d378635d`（本件＋gofumpt 三枚读数）／`c4ecb2c8`（票尾 bullet）／`f146781c`（交件后台面复跑＋§3 末独立复现）；本节这一行的最后一次改正＝**下一笔**；每笔 `git show --stat` 路径清单逐枚显式（`--pathspec-from-file` 里是一行一枚的具体路径，⛔ 目录形⛔ `-A`）。
2. 判语表＝本件：`D:\work\workspace\projects plans\Wisp\docs\evidence\s1\287-v1-accept.md`；读数根＝
   `D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\287\v1b\`（74 枚件，逐枚 `wc -c` 在名册里）。
3. 门禁逐枚 `rc=`：`GOFLAGS= go build ./...`＝**0**；`go test ./internal/tools/ -count=1 -v`＝**0**；`go vet ./internal/tools/`＝**0**；
   `sh scripts/d22scan.sh`＝**0**；`sh scripts/check-path-length-budget.sh --with-self-test`＝**0**；`gofmt -l`（两枚动过的件）＝**0**；
   `gofumpt -l`（走 PATH）＝**127**；`/d/work/base/gopath/bin/gofumpt.exe -l`（两枚动过的件）＝**0 枚列出**、整目录＝**1 枚**（同一枚 174-r7 存量件＝这把尺的正控）⇒ 见 §6-b。
   门禁一共跑过**三发**：**两发落盘**＝起手发（`v1b-01..07`＋`07-test-counts.txt`）与交件后发（`logs/r2-01..06`，五件 rc 全 0，长度门此时现量 `tracked paths=13967`／`over-budget=60`／`roster=60`、包级 `ok github.com/CarlosShao/wisp/internal/tools 18.240s`）；
   中间那发（写本件前）＝`GOFLAGS= go build ./...` rc=0＋`sh scripts/d22scan.sh` rc=0 两枚，**⛔ 单独落盘件、只在正文记录**＝本腿自抓，具名见 §6-g。
