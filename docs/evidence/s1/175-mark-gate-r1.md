# 175-r1 — 把 C25 的盖戳从「按工具名放行」改成「按内容来源放行」，并落一枚能红的常驻正控用例

- 派单：`.scratch/wisp/dispatches/2026-09-27-203x-impl-175-r1-make-the-stamp-not-name-gated-and-land-a-canary.md`（`7687ae0f`）
- 票面：`.scratch/wisp/issues/175-…-name-is-not-on-the-roster.md`（AC#2／AC#3 两格逐字见下 §4／§6）
- 前置裁定表：`docs/evidence/s1/175-c25-marking-roster-scope-c1.md`（A 支：名册＝举例定义）——**当前置结论用，不当判据**
- 本程写面：`internal/tools/bridge.go`（唯一一枚既有产码）＋`internal/tools/bridge_mark_provenance_ticket175_test.go`（新）＋`.scratch/wisp/probes/175/r1/**`＋本文件。**`internal/risk/**` 一字节未动**（`git status --porcelain -- internal/risk/` 见 §8）。
- 起手写码锚点 `7687ae0fbd01c089f01a5d13c0a58b072f828064`，分支 `dev`；`date` = `Sun Sep 27 20:35:53 CST 2026`。

## 1. step-0 五件＋改前基线

| 件 | 命令 | 读数 |
|---|---|---|
| 时间 | `date` | `Sun Sep 27 20:35:53 CST 2026` |
| 分支 | `git rev-parse --abbrev-ref HEAD` | `dev`（逐次复核：改前／改后各一次，均为 `dev`） |
| 锚点 | `git rev-parse HEAD` | `7687ae0f…` |
| 干净 | `git status --porcelain -- internal/tools/ internal/risk/ cmd/wisp/` | **空** |
| 改前基线 | `go test -count=1 ./internal/tools/ ./internal/risk/` | `ok tools 15.441s`／`ok risk 4.235s`（逐包，未跑全仓 `./...`） |

四数，两向（**顶层口径与子测试口径分开**；命令：`go test -count=1 -v <pkg>`＋`grep -c`）：

| 包 | 口径 | 改前（HEAD 快照树 `/tmp/wisp175pre`，`git archive HEAD` 解出） | 改后（工作树） |
|---|---|---|---|
| `internal/tools/` | 顶层 | `ok`，PASS=1 FAIL=0 | **`FAIL`**，PASS=1 FAIL=1 ← §7 那两枚腿 |
| `internal/tools/` | 子测试 | RUN=165 PASS=165 FAIL=0 SKIP=0 | RUN=168（+3＝本程新增三腿）PASS=166 **FAIL=2** SKIP=0 |
| `internal/risk/` | 顶层 | `ok` | `ok` |
| `internal/risk/` | 子测试 | RUN=172 PASS=171 FAIL=0 SKIP=1 | RUN=172 PASS=171 FAIL=0 SKIP=1（一字未动，逐字相同） |

名册口径现量（不是引裁定的）：`risk.IsSensitiveSource("task.output") = false`、`IsSensitiveSource("web.fetch") = true`（台件第一行）。

## 2. 本程没测／没做什么

- 没跑全仓 `go test ./...`（本仓规矩：逐包）；只跑了 `internal/tools/`、`internal/risk/` 两包＋`tools/d22scan` 的正控族。
- 没读、没引 `frontend/**`、`design/**`（他人领地）。d22scan 的 `design/` 计数 30→39 是工作树里他人未提交件造成的，与本程无关，也没去核。
- 没量票 176 的起跑口（`RunAsync`/`Record` 生产写者）——本程不接写者，只补「读回来要盖戳」这一侧。
- 没动 `[fs]`／`allowlist.txt`／档位（AC#4），`task.output` 仍是 L0（§7 那发读数里 `RiskLevel:L2` 是 **R4 事后升级**，不是声明档位，两回事）。
- 没测「非 builtin 的 C4 槽位（manifest/goja/mcp）结果要不要盖戳」：`Registry.RegisterProvider` 今天对这三槽回 `ErrSlotNotLanded`，生产路径一枚都注册不进来，所以那一支判据会是死码，本程不落（D46 命令插件那一路由票 163/176 与裁语 §4 的 `PLAN.md:2661` 管，届时判据本体已就位＝§6 的分类表那一腿）。

## 3. AC#2 未修码那一发（台件 `.scratch/wisp/probes/175/r1/`，零改动产码）

台件是一枚独立 Go 模块（`go.mod` 里 `replace github.com/CarlosShao/wisp => ../../../../../`），驱动**真 `risk.Provenance`＋真 `tools.Bridge`**（AGENTS §1.3 允许的接缝），只用两枚日志汇（`tools.Options.Logf` 与导出的 `risk.Logf`）做仪表——**这正是把「Mark 有没有被调用」与「戳有没有落地」分开的办法**。读数落本程自身目录（`ZZ175R1_DIR` 未设时程序直接死，不写继承来的 CWD）。

命令：`cd .scratch/wisp/probes/175/r1 && ZZ175R1_DIR=$PWD ZZ175R1_OUT=rig-reading-before-fix.txt go run .`

`rig-reading-before-fix.txt`（未修码，逐字）：

```
A1 execute            err=<nil> iserror=false risk=L0 text_prefix="外来内容探针 ZZ175r1-forei"
A2 mark_return_early  bridge saw no marking path: bridge lines mentioning "task.output" = []
A3 risk_log C25 lines ["risk/C25: sync-dir detection has no confirmed … (P12 safe default)"]
A4.1 scope_taints      n=0 []
A4.2 r4_inspect_hit    false src="" channel=""
A4.3 close_audit       ["tools: C25 scope closed task=175r1-A-task-output was_open=false dropped=0 open_scopes=0 close_err=<nil>"]
```

三问逐答（派单 §1 要的就是这三件）：

- **(i) `prov.Mark` 到底有没有被调用＝没有。** 分两半证：A3 里 `risk` 侧**没有** `Mark tool "task.output" is not a SPEC-06 §5 source` 这一行——而同一枚台件的 C 腿（越过桥、直接 `p.Mark(taskC, "task.output", …)`）**就是靠这一行**证明名册外名字走得进引擎：`C1 direct Mark() returned=true`＋`C2` 里出现了那一行。同一枚引擎、同一个名字，桥这一侧连日志都没有 ⇒ 早退发生在 `bridge.go` 的调用方，不在引擎。另一半：A4.3 `was_open=false dropped=0`——`mark()` 里 `b.OpenTask` 在早退之后，所以 scope 根本没开。
- **(ii) 返回的 `Result` 上有没有来源信息＝没有。** `E1 tool_result_origin "" err=<nil>`：`task.output` 不填 `Result.Origin`，而它的 `Decl.PathParams` 是空的 ⇒ 桥那侧 `origin` 的两条来路（Origin／`dec.Paths[0]`）都是空。修好之后仍然是空（见 §5 的残留）。
- **(iii) `Inspect`／R4 那一侧看不看得见＝看不见。** `A4.2 r4_inspect_hit false`。对照组 B 腿（同一枚桥、同一段内容、名册内的名字 `web.fetch`）：`B3.1 scope_taints n=1 [{… web.fetch https://example/foreign-body 45 …}]`、`B3.2 r4_inspect_hit true src="web.fetch https://example/foreign-body"`、`B3.3 … was_open=true dropped=1`。**同一份外部内容，换个门回来就没戳**＝本格成立（未修码上这一发不响）。

## 4. 选哪一支：乙。另两支为什么不选（都是现量，不是推理）

派单 §2 三形的代价各自量在 `.scratch/wisp/probes/175/r1/mut-shapeA/`（甲）与本表 §2／§4（丙）。两形都用 `go test -overlay=` 跑变异体，**产码一字节没退回去过**（`internal/tools/bridge.go` 全程只有乙那一处改动）。

**甲（凡成功结果都走 Mark）＝拒，代价已量出来。** 变异体 `bridge.shapeA-unconditional-mark.go`（overlay 换 `internal/tools/bridge.go`），落地证明：

```
internal/tools/bridge.go:552:	if b.prov == nil || dec.TaskID == "" || !risk.IsSensitiveSource(dec.Tool) {
.scratch/wisp/probes/175/r1/mut-shapeA/bridge.shapeA-unconditional-mark.go:552:	if b.prov == nil || dec.TaskID == "" {
```

`go test -count=1 -overlay=… ./internal/tools/` 的读数（`readings-shapeA-mutant-full.txt`；同一枚用例不带 overlay 时 `ok 0.151s`）：

```
--- FAIL: TestSensitiveReadAcrossTheBridgeOpensItsTaskScope
    bridge_scope_open_ticket158_test.go:98: 非敏感源 fs.list 也开了 scope … 开账必须挂在敏感源上，不能挂在每一次 Execute 上
--- FAIL: TestFSListSummarizesADirectory        fs_test.go:111: fs.list marked 1 tainted sources, want 0
--- FAIL: TestAtomicWriteKillsMidWrite/a_clean_write_lands_and_round_trips  … L2 审批通道尚未接入（票 21），已拒绝执行
--- FAIL: TestLongOutputPointerRecoversEveryByte …
--- FAIL: TestPointerPast256KiBIsNotFullyReadable  …
```

⇒ 甲不是一枚「更宽的盖戳」，它把**本进程自己写出来的确认文本**也当外部内容盖戳，随后 R4 拿这些戳去顶后续调用（`L2 审批通道尚未接入…已拒绝执行`），当场红掉五枚既有常驻用例，其中 `bridge_scope_open_ticket158_test.go:98` 那枚是**反向对照腿本体**（票 158 明写「开账必须挂在敏感源上」）。改法要活下来只能顺手把那一枚既有测试与 `fs_test.go:35` 那族一起改掉——那既超出本程写面，也不是修 bug。

**丙（产码自声明 `Decl` 加一位）＝拒，两条独立理由。** ① 写面不够：`Decl` 的本体在 `internal/tools/tool.go`，而本程写面只有 `internal/tools/bridge.go` 一枚既有产码文件（派单 §5），要落丙就得先扩射程——按派单「若这一支动到 … 停手上报」的精神本程不自己扩。② 判据方向不对：`Decl` 由产码方填，「这枚结果是不是外部内容」交给被门控的那一方声明，正是 `internal/risk/provenance.go` 头部 M-7 那一条拒的形状（豁免不许随产码方给的名字旅行）。裁语 §4 的 A 支要求的是「凡外部内容都要盖戳」，不是「谁承认谁盖」。

**乙（tools 侧立本地判据）＝取。** 落点 `internal/tools/bridge.go`：`marksProvenance(tool) = risk.IsSensitiveSource(tool) || outsideContentTools[tool]`，`mark` 的早退改吃这一枚。名册一字节没动，`internal/risk/**` 没动（A 支裁定如此）。乙自带的弱点＝「新工具要靠有人加名字」，本程用 §6 的第三枚判据（分类表全覆盖）把它从「静默漏」变成「红用例」。

⚠ **派单 §2 那句「甲会同时动 G6/G7 那两族腿的读数」现量为假**（对我这一程、也对甲那一形）：`.scratch/wisp/probes/154/gate-clauses.sh` 的 G6 主尺是 `git grep` 文本普查且显式排掉定义本体 `internal/tools/bridge.go`（脚本 `:452-453`），盖戳形状怎么改都不增不减那一族的行数。改前改后各一次逐字对比见 §7 表。这一条不按派单那句话去凑数。

## 5. 改后两向（不响→响）＋残留

`ZZ175R1_OUT=rig-reading-after-fix.txt go run .`（同一台件、同一码形，只多乙那一处改动）：

```
A3 risk_log C25 lines [ … "risk/C25: Mark tool \"task.output\" is not a SPEC-06 §5 source; recorded fail-closed as sensitive anyway (origin=\"\")" ]
A4.1 scope_taints      n=1 [{175r1-A-task-output task.output  45 1790513416}]
A4.2 r4_inspect_hit    true src="task.output" channel="notify"
A4.3 close_audit       ["tools: C25 scope closed task=175r1-A-task-output was_open=true dropped=1 …"]
```

两向并读：`scope_taints n=0 → n=1`、`r4_inspect_hit false → true`、`was_open=false dropped=0 → was_open=true dropped=1`、引擎那条名册外日志「没有 → 有」。

**残留（就地登记，不当已解决）**：`origin` 仍为空（`task.output` 不填 `Result.Origin`，填它要动 `internal/tools/task.go`，超出本程写面）。后果＝确认卡上念得出的是「来自 `task.output`」而念不出是哪个后台任务（`Hit.Source()` 在 `Origin` 为空时只回工具名）。修法＝`task.output` 把 roster key 填进 `Result.Origin`，属票 163/176 那一批。

## 6. AC#3＝常驻正控用例（本单真正买的东西）

落点 `internal/tools/bridge_mark_provenance_ticket175_test.go`（跑在 `scripts/portable-tests.sh --scope=core` 已含的 `./internal/tools/...` 里，不挂在任何 CI 不跑的尺上——裁语 §6 现量 `gate-clauses.sh` 在 `ci.yml` 零引用，所以判据没长在那族腿上）。三枚判据：

1. `TestTaskOutputMarkAndGate`：`task.output` 读回来的内容，来源标记必须非空，且 `Inspect` 那一侧必须命中并报出 `src="task.output"`（钉「盖了戳且 R4 看得见」，不是一颗空的索引条目）。
2. `TestNonContentResultStaysUnmarked`：`fs.list` 不得盖戳、不得开账（钉住乙与甲的边界，也钉住 §4 里甲红掉的那味）。
3. `TestEveryRegisteredToolIsClassifiedForMarking`：对本包注册的每一枚 builtin 工具问一句「盖不盖戳、是不是被判定过」，**没被分类＝红**。这一腿才是「不许依赖有人记得改名册」的落点：票 163/176 接上写者、或下一程新增一枚读回外部内容的工具而没登记，红的是 CI 核心域里的一枚用例，不是靠谁记得。

**能红（两向读数＋变异落地证明）**。变异体 `.scratch/wisp/probes/175/r1/mut-ac3/bridge.name-gated.go`（overlay，把早退改回名册闸门＝退回修前形状）：

```
internal/tools/bridge.go:578:	return risk.IsSensitiveSource(tool) || outsideContentTools[tool]
.scratch/wisp/probes/175/r1/mut-ac3/bridge.name-gated.go:578:	return risk.IsSensitiveSource(tool) // MUTANT 175-r1: the name gate is back (pre-fix shape)
```

`go test -count=1 -overlay=.scratch/wisp/probes/175/r1/mut-ac3/overlay.json ./internal/tools/ -run 'TestTaskOutputMarkAndGate|…'` ⇒

```
--- FAIL: TestTaskOutputMarkAndGate (0.00s)
    bridge_mark_provenance_ticket175_test.go:119: task.output 读回来的外部内容没有盖 C25 来源标记：scope "task-175-mark" 的污点表=[]。…
    bridge_mark_provenance_ticket175_test.go:130: 外来的任务输出被读回来之后，R4 一侧不再看得见它 …
FAIL
```

不带 overlay 的同一枚命令 ⇒ `ok github.com/CarlosShao/wisp/internal/tools 0.061s`（`readings-shipped.txt`）。**明写**：本用例钉的是「读回来要盖戳」这件事，不钉「有人真在写」——`TaskRoster.Record` 生产零写者（票 176），所以每一腿都自己注入 record。

## 7. 门禁五枚（改前／改后各一次）＋本程撞出来的那一格

| 门禁 | 命令 | 改前 | 改后 |
|---|---|---|---|
| 逐包 go test | `go test -count=1 ./internal/tools/ ./internal/risk/` | `ok tools`／`ok risk` | **`FAIL tools`（§7.1 两枚）**／`ok risk` |
| d22scan | `sh scripts/d22scan.sh` | rc=0 `clean - no D22 ban violations` | rc=0 `clean` |
| 正控族 | `sh tools/d22scan/runtests.sh -C tools/d22scan ./...` | `ok 22.369s`，RUN=76 | rc=0，`PASS=34 FAIL=0 SKIP=0, === RUN=76, [no tests to run]=0` |
| 名册差集 | 两份输出 `grep -oE '^[a-z]+/[^ ]+\.go:[0-9]+ .*' \| sort \| comm -3` | 两侧皆空（clean，零枚违规行）⇒ **差集为空**；唯一差异＝扫描面 `internal/ 425 → 426 Go files`（本程新增那一枚 `_test.go` 本体） | 同 |
| gate-clauses | `sh .scratch/wisp/probes/154/gate-clauses.sh` | rc=0，`腿数＝14 声明与实测不符＝0` | rc=0，`腿数＝14 声明与实测不符＝0` |
| **G6/G7** | 同上，逐腿 | `G6 1/1`、`G6pos 0/0`、`G6neg 1/1`、`G7 3/3`、`G7pos 0/0`、`G7neg 4/4` | **逐字相同**（声明＝基线＝实测，一枚没动） |
| flip | `sh .scratch/wisp/probes/161/r6/flip-declaration.sh` | 未跑（起手已由他人跑脏，见下） | rc=0 `GREEN …（baseline rc=0, 14 legs booked）` |
| gofumpt | `gofumpt --version`＋`gofumpt -l -w <两枚 .go>` | — | `v0.12.0 (go1.27.1)`，`-l` 两枚皆空（无重排） |

被 flip 弄脏的跟踪路径（跑完现列，不还原、不提交）：`.scratch/wisp/probes/161/r6/logs/flip-{1,2,3,4,5,6,baseline,restored}.txt`（` M`）＋`flip-7.txt`（`??`）。**这一组在起手 step-0 的 `git status` 里就已经是同一批脏件**，本程只是又跑了一次同一枚脚本，没有新增路径。

### 7.1 ⚠ 本程撞出来的一格：补戳会让票 164 的两枚常驻腿红（形状无关，不是乙独有的代价）

`go test -count=1 -v ./internal/tools/` 改后红的那两枚（改前 0 枚红）：

```
--- FAIL: TestLongOutputPointerRecoversEveryByte
    task_output_leg_test.go:220: the pointer must be re-readable with fs.read, got: {Text:L2 审批通道尚未接入（票 21），已拒绝执行 IsError:true RiskLevel:L2 …}
--- FAIL: TestPointerPast256KiBIsNotFullyReadable
    task_output_leg_test.go:381: the first segment must read back, got {Text:L2 审批通道尚未接入（票 21），已拒绝执行 …}
```

机制（现量，不是设想）：D15(3) 的续读桩文本里**逐字嵌着产物路径**（`…全文见 %s…`，`internal/tools/task.go:229`），盖戳之后那段路径就进了污染索引；同一 task scope 里下一次 `fs.read` 的参数只有 `{"path": <那个路径>}`，`Inspect` 对 path 形参数**不豁免**（`provenance.go:662` 那一支：`!gateOpen && !isPathKey(k)` 才跳，path key 永远扫）⇒ 命中 ⇒ R4 升 L2 ⇒ 测试用 `NoGate{}` ⇒ 「已拒绝执行」。两腿都用同一个 `TaskID`（`task_output_leg_test.go:181`／`:213`）。

⇒ **甲／乙／丙三形都会撞同一格**（甲那一形撞得更多，见 §4）。这不是选形问题，是裁定表 §4「补上 task.output 的盖戳＝修 bug，不需 owner 批」漏量的一格：盖戳一旦补上，「桩里那个宿主自己报出来的路径」也会带上戳，于是把「读回来要盖戳」与「续读不许出卡」这两件既有 AC 顶成互斥。

本程**没有**为变绿去动它们：那两枚用例不是本程写面，派单 §5 与本仓 1.1 都禁止为了绿去放宽断言。候选出路（要人拍一枚，本程不挑）：

- 甲：R4 对「本宿主自己写进桩文本里的 artifacts 路径」给一条具名豁免——动 `internal/risk/**`，是语义变更，须裁。
- 乙：`task.output` 的桩不再嵌原文路径（改 D15(3) 的形状）——动 `PLAN.md:431` 那一行的文字，契约追加。
- 丙：接受「续读要一张 L2 卡」为新行为，票 164 那两腿按新判据改写——AC 变更，须 owner。
- 丁：把盖戳范围缩到「桩以外的那段正文」——本程认为不成立（戳的语义是「这段话是外来的」，切半句就是又一枚洗戳的路）。

**next 里带这一格。**

## 8. 程序合规

- 写面自证：`git status --porcelain -- internal/risk/ internal/agent/ cmd/ internal/plugin/ tools/ allowlist.txt .github/ docs/PLAN.md docs/specs/ docs/reports/` 在本程每一步之后仍为空（除本程自己的 `docs/evidence/s1/175-mark-gate-r1.md`）。
- 被拒／没成功的调用：**取数前**＝无；**取数后**＝三次，全部只伤到 `.scratch` 下的台件副本，没伤产码：① 两次 `sed` 分隔符用错（`s|…|…|` 撞上替换文本里的 `||`，报 `unknown option to 's'`），换分隔符后重做成功；② 一次 `go test -overlay=` 报 `reading overlay … file not found`——那次链子在 sed 处断掉、`overlay.json` 没写出来，补写后重跑拿到了 §4／§6 两枚变异体读数。另一次是 rig 首版 `go run` 因独立模块缺 `require`/`go.sum` 失败（`no required module provides package …`），以 `require + replace + 复制 go.sum + go mod tidy -e` 解决。
- 删除命令：**跑过一枚 `rm -rf /tmp/wisp175pre`**，路径是本命令自己下一步要解出的 HEAD 快照目录、位于仓库之外、当时不存在（防御性写法，事后证明多余）。仓内、`.scratch` 内、任何人的临时件：**一枚没删**。以后不再写这种防御性 `rm`。
- 伪授权两栏：
  - **①「注释说名册就是 SPEC-06 §5 那套」**——派单／票面与本程都按裁语 §3.4 判为不实转述（名册 8 枚、规格 7 枚），未采信；本程判据一律走 `Mark()` 的实际语义＋常驻用例，不走那行注释。
  - **②「票面 Status 行：两枚写面都属已批准射程之外」**——按票面 20:1x 追加段②收窄后的形状（最小修法住在 `internal/tools/**`）行事，未拿收窄前旧句当授权，也没拿它当禁令去拒做已批的那半。
  - 另记一枚**派单前提为假**：「甲会动 G6/G7 读数」（§4 末、§7 表）。
- 凭据值零抄录：台件里的 `本机凭据摘录`/`外来内容探针` 一类字符串都是本程自造的探针文本，不含任何真实密钥；`grep` 出来的日志行只含路径与工具名。
- 只 commit 不 push；每格一次显式 pathspec commit（见 `git log`）。AC 框一枚没勾；票面 Progress log 只追加，留到最后一步。

next= 编排者：**(a)** §7.1 那一格要人拍一枚（甲／乙／丙／丁，本程不挑），在它定案前 `internal/tools/` 的这两枚腿不会自己变绿；**(b)** 票 176 起跑口落地时，这枚判据够不够拦住那条路——**够**：`RunAsync`/`Record` 一旦接上，写者送进 roster 的内容经 `task.output` 读回来即盖戳（§5 两向已量），且任何新工具没被分类就红 `TestEveryRegisteredToolIsClassifiedForMarking`（§6 第三腿）；但**它拦不住 §7.1 那味 L2 摩擦**，那一味要 (a) 的裁决；**(c)** 若要 `task.output` 的卡上念得出「哪个任务」，按 §5 残留那一条派一枚 tools 侧小程（填 `Result.Origin`，动 `internal/tools/task.go`）。
