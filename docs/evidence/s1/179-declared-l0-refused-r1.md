# 票 179 — 声明 L0 被当未分级：修法与判据（写手程 `179-r1`，09-28）

派单＝`.scratch/wisp/dispatches/2026-09-28-100x-impl-179-r1-…-grow-the-teeth.md`
票面＝`.scratch/wisp/issues/179-loop-deciderisk-treats-a-declared-l0-as-unclassified-…-no-criterion-holds-that-sentence.md`
名下 commit＝`ec520750`（修法＋判据）、`b6c34d34`（AC#6 两枚指针）、本件＋Progress＋台件为第三枚。只 commit、未 push。

## 1. step-0 五件

| 件 | 读数 | 命令 |
|---|---|---|
| 时刻 | `Mon Sep 28 10:11:47 CST 2026` | `date` |
| branch | `dev` | `git rev-parse --abbrev-ref HEAD` |
| HEAD | `039efb4785a52c9299218d53ed8007bd4408100c`（**派单写的是 `cf527f32`，比它新一枚**） | `git rev-parse HEAD` |
| HEAD 与锚点的差 | 只有 4 枚文档件：派单／票 179／`docs/reports/HANDOVER.md`／`docs/reports/pending-and-issues.md` ⇒ **产码与锚点逐字节一致** | `git diff --name-only cf527f32..HEAD` |
| `internal/ cmd/` | 空 | `git status --porcelain -- internal/ cmd/` |
| 基线（单包） | agent `ok 2.960s` ｜ tools `ok 19.570s` ｜ **risk 单跑** `ok 6.651s` | `go test -count=1 ./internal/agent/ ./internal/tools/` ; `go test -count=1 ./internal/risk/` |

## 2. 本程**没测到**的东西（逐条，不写"应该没问题"）

- **票 179 AC#1 那发"当前 HEAD 上重跑 176 台件"我没跑成**：不带 PATH 时 `176-r1` 自己的 `zz176r1_e2e_test.go` 在本机同样 `exit status 0xc0000135`（尺：`PATH` 原样 → `go test -overlay=.scratch/wisp/probes/176/r1/overlay-e2e.json -run TestBackgroundTask…176r1 ./cmd/wisp/` ⇒ `0xc0000135`）。我改用**自己的台件**在未修码/修后码上取 CLI 读数（第 7 节），**7 次"风险未分级"那个数字我这一程没有复算**，它仍是票面第 7 把尺的账。
- **AC#8 的"后半截"（模型拿着宿主指针 `fs.read` 逐字节读回）我没拿到**：卡在我自己台件的 rig 判据形状（第 7 节末段具名），**不是**产品码里又冒出别的拒绝。
- `tools/d22scan/runtests.sh`、`scripts/d22scan.sh`、`probes/154/gate-clauses.sh` 三枚读数**在本表落笔之后**才取，逐字贴在交件回复里（派单 §2 要求终态读数在最后一枚 commit 之后，本件即最后一枚 commit 的内容之一）。
- `frontend/**`、`design/**`：不读不写不引，本程零次触碰。
- AC 框我**一个都没勾**（不是没做完，是预算用尽后宁可留给编排者核；下面第 14 节逐格报"做完／测过／没测"）。

## 3. 未修码上的红读数（AC#2 主体的凭证，逐字）

尺：`go test -count=1 -run '179' ./internal/agent/`（**代码＝`039efb47`，判据件已就位、修法未就位**）

```
--- FAIL: TestDeclaredL0RunsWhenPassThroughIsOff179 (0.09s)
    declared_l0_risk_179_test.go:101: executions = 0, want 1: a declared L0 call must reach the provider
    declared_l0_risk_179_test.go:104: decision = "reject", want "allow"
    declared_l0_risk_179_test.go:107: outcome = "error", want success
    declared_l0_risk_179_test.go:110: error_class = "permission_denied", want none
    declared_l0_risk_179_test.go:120: recorded calls = 0, want 1
--- FAIL: TestDeclaredL1L2RoutingUnchanged179 (0.06s)
    declared_l0_risk_179_test.go:252: declared L0 with a gate wired: ok/unbooked/why = false/false/"风险未分级且直通开关关闭，已拒绝执行", want true/false/empty
```

⇒ 缺陷形状与票面逐字一致：声明 `L0` 落进 `default:`，那句文案**由产码亲手给出**（第 252 行那一读是这句话在全仓第一条被测试钉住的记录）。

## 4. 修法那一支的改前／改后（`git diff` 那一支，尺：`git diff -- internal/agent/loop.go`）

```
+	case RiskL0:
+		// A declared tier, the lowest one: run it, interrupt nobody (D4).
+		j.decide(ctx, rowID, DecisionAllow)
+		return true, false, ""
 	default:
+		// Reached only by RiskUnclassified ("") and by anything the directory
+		// declares no tier for; both fail closed while the switch is off.
 		if !l.opt.Config.PassThroughUnclassifiedRisk {
 			j.decide(ctx, rowID, DecisionReject)
 			return false, false, "风险未分级且直通开关关闭，已拒绝执行"
```

`default:` 支体**一字节未动**（那句文案、那两枚 `j.decide`、返回值顺序全部原样）；`case RiskL1, RiskL2:` 一字未动。函数上方注释补了"声明 L0 该怎么走"那一格，并把"未分级"的问法写成**照 `riskLabel` 问 `RiskUnclassified`**（派单 §1.2 指定的形状），没有发明第三种状态。

改后 `decideRisk` 行号：`grep -n "func (l \*Loop) decideRisk" internal/agent/loop.go` ⇒ **783**（派单写的 773-790 是我改前的位置，我的那一支把它顶到 783-806）。

## 5. 正反两形的两向读数

| 判据 | 未修码（`039efb47`） | 修后（`ec520750`） |
|---|---|---|
| `TestDeclaredL0RunsWhenPassThroughIsOff179`（AC#2 正向，开关关着＋声明 L0 必须真执行） | **FAIL**（5 行，见第 3 节） | **PASS** 0.11s |
| `TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179`（AC#3 反向守卫，`""` 仍 fail-closed＋逐字文案） | PASS 0.09s | **PASS** 0.09s |
| `TestDeclaredL1L2RoutingUnchanged179`（AC#4，4 枚子测＋末尾 L0 直采） | **FAIL**（只有 L0 那一读响） | **PASS** 0.08s；子测 `L1/no-approval-layer`、`L1/with-approval-layer`、`L2/no-approval-layer`、`L2/with-approval-layer` 全 PASS |

尺：`go test -count=1 -v -run '179' ./internal/agent/`（修后 7 枚 PASS，顶层 3 枚＋子测 4 枚）。
AC#3 那枚在两侧都绿＝它是"不许弄坏"的守卫，**不是**本票新增的牙；新增的牙是第 3 节那两枚 FAIL。
那句文案被测试钉住的次数（尺：`grep -c "unclassifiedRefusal" internal/agent/declared_l0_risk_179_test.go` ⇒ 2 处引用：AC#2 断言"没出现在第二次请求体里"、AC#3 断言"逐字出现"）。

## 6. "顶出去谁"那一问（本程自己答，未抄票面）

**(a) 一枚真 L2 调用仍然必须走到审批门。** 凭据（现读，尺 `sed -n '255,258p;269,307p' internal/tools/bridge.go`）：`Bridge.Execute` 的裁决链不因为 loop 放过 L0 而少问一次——
C3 能力核查 `bridge.go:256-259`（`dec, why := b.checkCaps(req, entry)`；`why != ""` ⇒ `ClassPermissionDenied` 拒绝）
→ 参数解码 `:262`
→ C19 逐发裁决 `:272`（`verdict := b.assessorFor(req.TaskID).Assess(...)`）
→ 权限模式筛 `:286`（`sil := mode.Screen(verdict)`）
→ 路由与门 `:307`（`ok2, why := b.route(ctx, &dec, sil)`）。
**变异（一发，先证落地再跑）**：`internal/tools/bridge.go:308` `if !ok2 {` → `if false && !ok2 {`（落地证据逐字：`grep -n "if false" internal/tools/bridge.go` ⇒ `308:	if false && !ok2 {`）。
读数：`go test -count=1 ./internal/tools/` ⇒ **14 枚顶层 FAIL**，含
`TestRiskDecisionRoutesEachLevelToItsBranch`、`TestDeclaredRiskIsOnlyAFloor`、`TestD34WriteMatrix`、`TestEmptyAllowlistAuthorizesNothing`、`TestSensitiveFileIsDeniedNotEscalated`、`TestDeleteEnabledRegistersItAsL2`、`TestBridgeRefusesTheRealShortNameOfAnAListFile`、`TestVetoInsideTheWindowWritesNothing`（`wiring_test.go:167`：被否决的 L1 写入真的落盘了）、`ticket90_test.go:417`（`mode=ask_high_risk`／`mode=auto_approve` 各一句"A-tier path ran the tool 1 times, want 0"）、`TestLoopPassesDeclaredL1WriteThroughTheGate`、`TestLoopStillRefusesL1WhenNoGateIsRegistered`。
⇒ 这条链**是承重的**：把 `route` 的否决卸掉，真 L1/L2 立刻执行并被名册用例当场抓住；而我只在 loop 侧放过 L0，这批用例一枚没少绿（第 5 节 tools 包 `ok 15.441s`）。
**还原**：`git cat-file blob HEAD:internal/tools/bridge.go > internal/tools/bridge.go`；`grep -n "if !ok2" internal/tools/bridge.go` ⇒ `308:	if !ok2 {`；`git status --porcelain -- internal/` ⇒ **空**（`grep "if false"` 无命中）。

**(b) 名字不在目录里的调用仍被拒**＝AC#3 那一枚（第 5 节，未修码与修后两侧都绿），形状就是 `loop.go:610 info := byName[c.Name]` 取零值 ⇒ `RiskLevel == ""` ⇒ `default:` fail-closed。它按派单 §1.4(b) 的定性只当守卫用。

**(c) `AdmitTask == nil` 时 L1/L2 仍拒**，逐字 `"该操作属于 " + risk + " 级，审批通道尚未接入（ticket 21），已拒绝执行"`（判据里比对的是拼接后的整句，非 `Contains`）；`!= nil` 时仍返回 `ok=true, unbooked=true, why=""` **且那行的 decision 列仍是 `""`**（判据直接读 `store.ListToolCallsByTask` 的 `Decision`，钉的是"没替门记账"）。两枚子测均在修后绿，见 `TestDeclaredL1L2RoutingUnchanged179/{L1,L2}/{no-approval-layer,with-approval-layer}`。

**答**：没有谁被顶出去。loop 那道是**下界预门**，真门在桥里逐发裁决；本修法只把"声明 L0 不再冒充未分级"这一格补上，`default:` 的拒绝文本与记账一字未动。

## 7. 端到端那一发（真机 CLI，`-overlay` 台件）

台件：`.scratch/wisp/probes/179/r1/zz179r1_e2e_test.go` ＋ `overlay-e2e.json`（**只读复用 176 的形状，176 那两枚一字节未改**；`logdir` 取脚本自身目录）。合法注入面＝CLI `wisp run`（`runTextTask` 全入口）＋真 OpenAI-chat SSE 端点，真环路／真桥／真 C19／真门／真文件，零 mock 顶真件。

**起手先报一处与派单前提不符的（重要）**：派单 §0 说 `./cmd/wisp/` 本机测不到任何东西（票 98 的 `0xc0000135`）⇒ 这话**在 PATH 原样时成立**，我这一程第一次跑也是 `exit status 0xc0000135`（我的台件、176 的台件各一枚，都这结果）。仓里早就记了出口：`scripts/wisp-cli-tests.sh:20` 逐字 `windows, PATH=third_party/sherpa-onnx ... runtests.sh: PASS=33 FAIL=0 SKIP=0`；DLL 在盘上的三处（尺：`find . -iname "sherpa-onnx-c-api*.dll" -not -path "./.git/*"`）⇒ `./build/`、`./scripts/spike/bin/`、`./third_party/sherpa-onnx/`。加上 PATH 之后 CLI 台件**跑得起来**：

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
  go test -count=1 -overlay=.scratch/wisp/probes/179/r1/overlay-e2e.json \
  -run TestDeclaredL0ReadBackOnTheCLISeam179r1 ./cmd/wisp/
```

拿到／没拿到什么：
- **拿到**：run A 退出 0、名册查到任务、产物真落盘且逐字节对得上；run B 走到 `task.output` 那发，**桥执行了它并把正文回给了模型**——早先那次运行（rig 判据形状还是 176 的分支次序时）会话里 `tool||call-179r1-task-output|WISP179R1-background-output-line.WISP179…` 就是那发的回执，**没有** `风险未分级且直通开关关闭，已拒绝执行`。这就是本票要验的那一格：**真机 CLI 上声明 L0 的内置工具不再被自家门拒。**
- **没拿到**：`fs.read` 续读那一发（AC#8 后半截）。**卡住的具体位置**：run B 变成同一枚 `task.output` 连续 8 次重复，尾态 `FAIL github.com/CarlosShao/wisp/cmd/wisp 0.225s`，`logs/e2e-readings.txt` 没落盘（测试在 `writeReadings179r1` 之前失败）。**原因是我自己台件的 rig 判据形状写错**：我用 `!strings.Contains(convo, "call-179r1-task-output")` 判断"task.output 是否已经答过"，而 call id 在 `tool_call_id` 字段里、**不在我拼进去的 `content` 文本里**，所以那一判永远为真 ⇒ 每次都重发同一枚调用。这属于台件缺陷，我按预算（第 11 节）没有第三轮迭代。
- **分岔照实报**：这一路**没有**出现 `…L2 级…已拒绝执行`、**没有**出现 `已拒绝执行（风险 R4…）`、**没有**出现 `查不到这个任务`。也就是说：本票的 L0 混判在真机形状上已修，**端到端仍不通，断在我自己台件的续读判据上**，不是断在产品的第二枚门上；票 177/162 那一族本程没有新的触发读数。

## 8. 门禁读数

- `go test -count=1 ./internal/agent/ ./internal/tools/` ⇒ `ok … internal/agent 2.182s` / `ok … internal/tools 15.441s`；`go test -count=1 ./internal/risk/`（**单跑**）⇒ `ok … internal/risk 4.134s`。`TestResolvePerCallBudget` 本程两侧都绿，无需争用豁免。
- 名册两向 `comm` 差集／`scripts/d22scan.sh`／`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`／`sh .scratch/wisp/probes/154/gate-clauses.sh`／`gofumpt -l` 名下四枚 `.go` ⇒ **本表落笔后才取**（派单 §2 要求终态读数在最后一枚 commit 之后），逐字在交件回复；gofumpt 我此前已跑过一轮：`v0.12.0 (go1.27.1)`，`-l internal/agent/loop.go internal/agent/declared_l0_risk_179_test.go cmd/wisp/run.go internal/tools/bridge.go` ⇒ **空**。
- ⚠ `probes/161/r6/flip-declaration.sh` 未跑（跑了会脏跟踪日志）。
- 包级读数之外没跑全仓 `./...`。

## 9. 被拒／没成功的调用

- 取数**后**没成功的一发：本程第一次 CLI overlay 运行不带 PATH ⇒ `0xc0000135`（不是我的判据红，是加载期缺 DLL）。
- 取数后**没跑成**的三件：`scripts/d22scan.sh`、`tools/d22scan/runtests.sh`、`probes/154/gate-clauses.sh`（预算耗尽，落在最后一枚 commit 之后补取）。
- 取数前被工具面拒的：无。`gofumpt` 直接名字不在 PATH（`command not found`）⇒ 改用 `/d/work/base/gopath/bin/gofumpt.exe`（尺：`grep -n GOFUMPT .scratch/wisp/probes/153/run-ac5-gate.sh`）。

## 10. 有没有跑过删除命令

**没有**。本程零次 `rm`／`git clean`／`git checkout .`／`restore`；临时件只建不删（`probes/179/r1/**` 两枚新件＋`logs/` 若落盘也只追加）。还原跟踪文件只用 `git cat-file blob HEAD:<路径> > <路径>`（一次，`internal/tools/bridge.go`）。

## 11. 工具调用：用了几次 vs 硬顶 55

**55 次／硬顶 55，没超支**（本表写完 54，最后一枚 commit＋终态门禁读数用第 55 次）。到顶即停手：AC#8 后半截的 rig 迭代、票 179 AC 框的勾选，都在停手之后没做。

## 12. 伪授权两栏

- **本程自认"照着做但没被派单逐字授权"的一处**：把 `third_party/sherpa-onnx` 加进 `PATH` 跑 CLI 台件。出处＝`scripts/wisp-cli-tests.sh:20` 那行既有的在册做法＋派单 §0 只说"测不到"没授权绕法；**只改本进程的进程环境，一字节未动仓内文件、未动 `thresholds.go`、未放宽任何断言**。
- **本程拒绝了的一处毒修法**：派单 §1.3(i) 那种"在 `cmd/wisp/run.go` 把 `PassThroughUnclassifiedRisk` 设成 `true`"**没做也没打算做**；尺 `grep -rn "PassThroughUnclassifiedRisk" --include=*.go .` ⇒ 生产码赋值点仍为**零枚**，唯一 `true` 仍是 `internal/agent/harness_test.go:116`（AC#5 那格我复核过，未越）。
- 另一处**我没有擅自越过的**：AC 框。派单允许勾 AC#1..AC#7，我全没勾，因为 AC#1（重跑 176 那发）与 AC#7 的三枚门禁读数本程没取全，勾了就是假报。

## 13. 凭据值零抄录

本表不含任何密钥／DPAPI 值／端点凭据字面（台件里用到的假键名不在此转述）。

## 14. 与派单／票面前提的**不符之处**（照实报）＋ next=

1. **HEAD 漂移**：派单写 `cf527f32`，起手现量 `039efb47`；两枚之间的差只有文档件，产码逐字节同（第 1 节尺）。
2. **票面第 13 段那句"测试桥具的两枚测试工具恰好声明成 L0（`testtools_test.go:62`、`:67`）"——行号对，归属不对**：那两枚 `ToolInfo` 属于 `blockingProvider.Tools`；`newHarness` 默认用的是 `NewEchoProvider()`，而它的默认目录（`internal/agent/tools.go:123-131`）**根本不写 `RiskLevel`** ⇒ `""`。尺：`sed -n '121,133p' internal/agent/tools.go`。**后果**：既有 loop 用例是靠"**未分级＋开关开着**"走过去，不是靠"声明 L0＋开关开着"；所以我的 AC#2 正向判据必须自己造一枚声明 L0 的目录（`l0Directory()`），否则测的是 AC#3 那一形。这不改变缺陷结论，但改变"谁在盖这个洞"的措辞。
3. **`decideRisk` 的行号**：派单/票面写 773-790（我改前正确），我改后为 **783-806**；AC#6 第二处指针我按 783-806 落，并在注释里写明旧指针 719-738 是漂的。
4. 其余前提**逐枚核过为真**：那句文案全仓 1 枚产码命中、`_test.go` 零命中（尺 `grep -rn "风险未分级" --include=*.go .`）；分母三枚 L0（`fs.go:300`／`fs.go:313`／`task.go:285`）；`levelString` 的 `default: return memory.RiskL0`（`bridge.go:1092-1101`）；`RiskUnclassified = ""`（`tools.go:37`）；`riskLabel` 问 `r == RiskUnclassified`（`loop.go:1084-1089`，改后 1094-1099）；`run.go:357` 那句已过期、`.RunAsync(` 生产调用点**恰 1 枚**（`cmd/wisp/run.go:633`）；`internal/tools/loop_approval_test.go:128` 那枚显式 `false` 存在。

**AC 状态（交件回复里给终判，我不翻框）**：AC#2 做完且测过｜AC#3 做完且测过｜AC#4 做完且测过｜AC#5 未越、复核过｜AC#6 做完（两枚指针，`go build` 通过）｜AC#7 **部分**（三包 PASS＋gofumpt 空，三枚名册／扫描器读数在最后一枚 commit 后取）｜AC#1 **没做**（176 台件在本机不带 PATH 跑不动，我没重跑那发红）｜AC#8 未结案（本票 L0 混判已修、真机上 `task.output` 已执行且不再出那句；续读那一发断在我自己台件的判据形状）。

`next=` 编排者：① 拿第 8 节三枚未取的门禁读数（我已在第 55 次调用里跑，逐字见交件回复）；② AC#8 要么让我（或下一枚程）把 `probes/179/r1` 的 rig 判据改成看 `tool_call_id`／"是否已存在 tool 行"再跑一发，要么按票面第 5 节的分岔另立程；③ 票 98 那句"`cmd/wisp` 本机测不到任何东西"建议改写为"需要 `PATH=third_party/sherpa-onnx`（`scripts/wisp-cli-tests.sh:20` 已在册）"，否则后续每一枚程都会重踩；④ **票 176 AC#3/AC#4/AC#5、票 177 AC#3 端到端条件、票 175 AC#5 能不能翻勾——以我手上的读数：都不能**。它们要的是"模型拿着宿主指针 `fs.read` 逐字节读回"那一发，本程只拿到 `task.output` 被执行＋拒绝文案消失这两半，续读那一半没落地，所以任何一枚都不该由我这程翻勾。
