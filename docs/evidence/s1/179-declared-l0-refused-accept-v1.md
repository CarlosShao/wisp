# 179-v1 — 非实现者验收：三枚新判据的恒真性、两枚守卫的成色、以及"声明 L0 由 loop 自己记 `allow`"这一格归属

- 程：`179-v1`（只裁不改）。派单＝`.scratch/wisp/dispatches/2026-09-28-114x-accept-179-v1-…-who-owns-the-decision-column.md`
- 起手 `date` 现量：**Mon Sep 28 11:54:20 CST 2026**；分支 **`dev`**；起手 HEAD **`e2dde3f0`**
  （派单写"我这边 `95218c4c`"，实际已到 `e2dde3f0`＝两枚派单文档 commit，`git diff` 不含产码 ⇒ 登记，不据此改判）
- 被告陈述（**不是证据**）：`docs/evidence/s1/179-declared-l0-refused-r1.md`、`docs/evidence/s1/179-cli-reread-r2.md`
- 本程台件与读数：`.scratch/wisp/probes/179/v1/logs/`（`d22scan.txt`｜`runtests.txt`｜`gate-clauses.txt`｜`my-roster.txt`｜`r2-roster.txt`｜`m1-full-package.txt`｜`m2-full-package.txt`｜`cmdwisp-with-prefix.txt`｜`cmdwisp-no-prefix.txt`）

---

## 1. step-0 五件

尺与读数（逐条可复制）：

| 件 | 命令 | 读数 |
|---|---|---|
| 时刻 | `date` | `Mon Sep 28 11:54:20 CST 2026` |
| 分支 | `git rev-parse --abbrev-ref HEAD` | `dev` |
| HEAD | `git rev-parse HEAD` | `e2dde3f0d9351539e773e30519ed6c5bfcc45727` |
| 树净 | `git status --porcelain -- internal/ cmd/` | **空** |
| 基线 | `go test -count=1 ./internal/agent/ ./internal/tools/`；`go test -count=1 ./internal/risk/`（**单跑**） | `ok agent 1.775s` / `ok tools 15.048s` / `ok risk 3.728s` |

与派单不符的一处：HEAD 不是 `95218c4c` 而是其后一枚 `e2dde3f0`（内容＝两枚派单文本）。其余四件与派单一致。

## 2. 本程没测什么（先划边界，免得被当成裁过）

- **没跑**票 179 AC#8 那发端到端（`probes/179/r2/zz179r2_e2e_test.go` 的 `-overlay` 腿）⇒ R4 升 L2 的次生那一发我**只引 r2 的读数、未复算**。
- **没复算**票面第 7 把尺那个"7 次风险未分级"（`probes/176/r1/logs/last-run.txt` 我没读）。
- **没跑** ubuntu 腿（本机无 linux）⇒ 票 98 的"19 枚红归因"我只查了台账措辞，没实测。
- **没动** `internal/risk/**`、`bridge.go`、`thresholds.go`、`allowlist.txt`、任何 `probes/**` 既有台件（只读）。
- **没跑** `probes/161/r6/flip-declaration.sh`（派单禁）。
- AC#2 里 `risk_level` 那一枚断言的"怎么改会红"我是**推理＋指名**，**没跑**那枚变异（见第 3 节 C1 行）。
- 前端面（`frontend/**`／`design/**`）不读不写不引。
- `183-a1` 同时在飞：我全程只见到 `internal/agent/loop.go` 上我自己的改动，**没观察到非我的跟踪件改动**（每枚还原后 `git status --porcelain -- internal/ cmd/` 均空，见第 4 节）。

## 3. 三枚判据逐枚表（每枚：钉什么／产码怎么改会红／我动不动了它）

判据件＝`internal/agent/declared_l0_risk_179_test.go`（`ec520750` 新增，255 行）；产码面＝`internal/agent/loop.go:783-806`（`decideRisk`）。
尺：`grep -rn "^func Test.*179" internal/agent/*_test.go` ⇒ **3 枚顶层**（`79`／`143`／`195`）。派单说的"三枚判据"**核对为真**（顶层 3 枚；含子测共 7 枚读数，`go test -v -run 179` 现量）。

| # | 判据 | 它钉什么 | 产码怎么改会让它红 | 我动不动了它 |
|---|---|---|---|---|
| C1 | `TestDeclaredL0RunsWhenPassThroughIsOff179`（AC#2 正向） | 开关关着＋**目录声明 L0** ⇒ 必须真跑到 provider（`CallCount()==1`）、journal 那行 `decision=allow`＋`outcome=success`＋无 `error_class`、**正文回到第二次请求体**、且第二次请求体里**不得**出现那句拒绝文案 | ① 摘掉 `case RiskL0:` 整支（＝回到未修码）→ 红（`:101`/`:104`/`:107`/`:110`/`:120` 五枚断言一起红）② 留着支但删 `j.decide(ctx, rowID, DecisionAllow)` → 红（`:104 decision = ""` ）③ 把执行改成 skip → 红 | **动了**：M1（摘支）＋M2（摘记账）两枚变异，各见第 4 节 |
| C2 | `TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179`（AC#3 反向守卫） | 同样开关关着，`RiskLevel==""` 那一发（目录没那枚名字 ⇒ 零值 `ToolInfo`）必须**仍被拒**：`CallCount()==0`、`decision=reject`、`outcome=error`、`error_class=permission_denied`、第二次请求体里**逐字保留** `unclassifiedRefusal` | ① `default:` 的 `if !PassThrough…` 极性反过来 → 红（`:171`/`:174`/`:177`/`:180`/`:188` 五枚一起红）② 把 `default:` 改成一律放行（票面点名的毒修法）→ 红 ③ 把 `case RiskL0:` 写成 `case RiskL0, "":`（把未分级并进 L0）→ 红 | **动了**：M3（极性反），见第 4 节 |
| C3 | `TestDeclaredL1L2RoutingUnchanged179`（AC#4 守卫） | 四枚子测（L1/L2 × `AdmitTask` nil/非 nil）＋一枚尾测（声明 L0 且门已接）：nil 形仍逐字拒＋那行记 `reject`；非 nil 形仍 `ok=true/unbooked=true/why=""` **且那行 `decision` 仍是 `""`**；尾测 L0 仍 `unbooked=false` | ① nil 形：改那句文案或去掉 `j.decide(DecisionReject)` → 红（`:211-214` 与 `:219-221` 是两枚**独立**断言）② 非 nil 形：在 `case RiskL1, RiskL2:` 里加一行 `j.decide(DecisionAllow)`（loop 替门记账）→ 红（`:238-240` 直读 `store.ListToolCallsByTask`）③ 让 L0 走 `unbooked=true`（把 L0 卷进门那一支）→ 尾测红 | **部分动了**：M1 让尾测红（`:252`）；M2 **没**让它红 ⇒ 见下"这一枚的洞" |

**三枚都不是恒真**（每枚都答得出"产码怎么改会让它红"，其中 C1/C2 的红是我自己跑出来的）。

**C1 里有一枚近恒真的断言（本程最有值的一格）**：`declared_l0_risk_179_test.go:112` 的 `if row.RiskLevel != memory.RiskL0`。
理由：那一列不是 `decideRisk` 写的，是 `loop.go:612` 的 `j.startCall(…, riskColumn(info.RiskLevel))` 写的，而 `riskColumn`（`internal/agent/journal.go:150-156`）逐字是
`switch risk { case memory.RiskL1, memory.RiskL2: return risk; default: return memory.RiskL0 }`
⇒ 未分级 `""`、`"L0"`、任何野字符串进这一列**都变成 `"L0"`**。能把它造红的产码改法只有一处：改 `journal.go:154-155` 那支 default（本票禁地面，我没跑）。
⇒ 判：**这一枚断言不承重**，但它红不红都不影响 C1 整体成立（其余五枚断言承得住）。

**C3 的洞（M2 实测出来的）**：尾测（`:246-254`）只问 `ok/unbooked/why`，**没问那一列记没记**。M2 状态下它照样 PASS，而 journal 那行其实已经永远是 `""`（挂账）。⇒ 判：**不是恒真，但比 C1 少一枚牙**；补法（不由我落，本程禁改产码与判据）＝尾测加一枚 `ListToolCallsByTask` 读 `Decision==DecisionAllow`，与 `:238-240` 同形。

## 4. 我自己动的三枚变异（落地证明＋两向读数＋还原后空 status）

还原协议：`git cat-file blob HEAD:internal/agent/loop.go > internal/agent/loop.go`。全程没用 `checkout`／`restore`／`reset`／`stash`／`--amend`，没删任何文件。

### M1＝摘掉 `case RiskL0:` 整支（复现未修码）

落地证明（改完先贴行，再跑）：

```
$ grep -n "case RiskL0\|MUTATION 179-v1 M1\|RiskL1, RiskL2" internal/agent/loop.go
786:	case RiskL1, RiskL2:
793:		// MUTATION 179-v1 M1: the declared-L0 arm was removed here on purpose.
$ git status --porcelain -- internal/
 M internal/agent/loop.go
```

变异下读数（`go test -count=1 -v -run '179' ./internal/agent/`＋整包）：

```
--- FAIL: TestDeclaredL0RunsWhenPassThroughIsOff179
--- PASS: TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179
--- FAIL: TestDeclaredL1L2RoutingUnchanged179   (子测 :252 那一枚 L0 尾测红)
declared_l0_risk_179_test.go:101: executions = 0, want 1: a declared L0 call must reach the provider
declared_l0_risk_179_test.go:104: decision = "reject", want "allow"
declared_l0_risk_179_test.go:107: outcome = "error", want success
declared_l0_risk_179_test.go:110: error_class = "permission_denied", want none
declared_l0_risk_179_test.go:120: recorded calls = 0, want 1
declared_l0_risk_179_test.go:252: declared L0 with a gate wired: ok/unbooked/why = false/false/"风险未分级且直通开关关闭，已拒绝执行", want true/false/empty
```

整包（`go test -count=1 ./internal/agent/`，`logs/m1-full-package.txt`）：**`grep -cE "^--- FAIL"` ⇒ 2 枚，且两枚都是 179 名下**（`rc=1`，`FAIL agent 1.978s`）⇒ **既有零枚因这枚变异而红**（这一发是第 6 节① 的直接证据）。

HEAD（未变异）对照：三枚全 PASS（`go test -count=1 ./internal/agent/` ⇒ `ok 1.789s`，见第 8 节终态）。

### M2＝留着 `case RiskL0:` 但删掉 `j.decide(ctx, rowID, DecisionAllow)`

落地证明：

```
$ sed -n '792,796p' internal/agent/loop.go | cat -n
     1	case RiskL0:
     2		// MUTATION 179-v1 M2: the booking line below was deleted on purpose.
     3		// A declared tier, the lowest one: run it, interrupt nobody (D4).
     4		return true, false, ""
     5	default:
```

读数：`--- FAIL: TestDeclaredL0RunsWhenPassThroughIsOff179`，逐字 `declared_l0_risk_179_test.go:104: decision = "", want "allow"`；C2/C3 **均 PASS**。
整包（`logs/m2-full-package.txt`）：`^--- FAIL` **2 枚**＝`TestDeclaredL0RunsWhenPassThroughIsOff179` ＋ **既有的** `TestCancelledTaskPersistsTerminalRows`（`internal/agent/forensics_test.go:82`）。
⇒ 两个结论：(a) "记 `allow`" 那一支**不是装饰**，C1 的 `:103-105` 直接钉它；(b) **既有测试确实在 HEAD 上走过 L0 那一支**（否则 M2 打不到它）——见第 6 节①。

### M3＝把 `default:` 的 fail-closed 条件反过来

落地证明：

```
$ sed -n '797,805p' internal/agent/loop.go
		// Reached only by RiskUnclassified ("") and by anything the directory
		// declares no tier for; both fail closed while the switch is off.
		if l.opt.Config.PassThroughUnclassifiedRisk { // MUTATION 179-v1 M3 inverted
			j.decide(ctx, rowID, DecisionReject)
			return false, false, "风险未分级且直通开关关闭，已拒绝执行"
		}
```

读数（只跑 C2）：

```
--- FAIL: TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179
:171: executions = 1, want 0: an unclassified call must never run
:174: decision = "allow", want "reject"
:177: outcome = "success", want error
:180: error_class = "", want permission_denied
:188: second request lost the verbatim fail-closed sentence "风险未分级且直通开关关闭，已拒绝执行"
```

⇒ **守卫可造红＝不是装饰**（票面 AC#3 点名的"把 `default:` 改成一律放行"那一发正是它拦的形状）。

### 还原与终态

```
$ git cat-file blob HEAD:internal/agent/loop.go > internal/agent/loop.go
$ git status --porcelain -- internal/ cmd/
（空）
$ git diff --stat HEAD -- internal/ cmd/
（空）
```

三枚变异各自还原后都复跑过这一对；最后一枚的读数见第 8 节（终态在全部 commit 之后取）。

## 5. 守卫那一问的答案（派单 §1.2）

**〔成立·带条件〕**——两枚守卫都**合格**，但**都不许算本票新增的牙**，且 C3 少一枚牙：

1. 合格性的直接证据是 M3：C2 今天绿**不是**因为它恒真，而是因为它守的那一支本来就对；极性一反就红（五枚断言一起红）。C3 的四枚子测里，`AdmitTask == nil` 与"非 nil 不记账"**各有独立断言**，不是共用一枚——尺：`:219-221`（nil 形：读 `store.ListToolCallsByTask`，要求 `rows[0].Decision == DecisionReject`）与 `:238-240`（非 nil 形：同一枚 store 读法，要求 `rows[0].Decision == ""`）是**两处分开写的断言、两个分开的 taskID**（`task-179-nil-<tier>`／`task-179-admit-<tier>`），互不遮挡。这一格派单担心的"是不是共用一枚"＝**不是**。
2. 但 C3 的**尾测**（声明 L0 且门已接）只问 `ok/unbooked/why`、**没问那一列**：M2 实测它 PASS ⇒ 尾测买到的是"路由不改变"，**没买到"记账不丢失"**。具名变异：删 `loop.go:794` 那一行 `j.decide(ctx, rowID, DecisionAllow)`（＝我的 M2）。
3. 两枚守卫的**在册身份**与票面一致：票面 AC#3/AC#4 自己写了"这两枚今天也绿，只能当守卫"。我没有把它们当新增牙——本票新增的牙只有 C1（＋C3 尾测那一枚 `unbooked=false`）。

## 6. 推翻我两条前提的独立复算（派单 §1.3；我不沿用编排者那句"实现程对"）

### ① "声明 L0 这一形从来没被任何既有用例跑过" —— **〔成立·带条件〕，原话超范围，要收窄**

核到的部分（成立）：
- `NewEchoProvider()` 无参那支的目录**逐字不写 `RiskLevel`**（尺：`sed -n '119,133p' internal/agent/tools.go`，两枚 `ToolInfo` 只有 `Name/Description/Resident/Parameters`）；`newHarness` 在 `o.Tools == nil` 时正是走无参那支（尺：`sed -n '126,131p' internal/agent/harness_test.go`）；同一枚 harness 把 `PassThroughUnclassifiedRisk: true`（尺：`sed -n '112,118p' internal/agent/harness_test.go`，`:116`）。
- ⇒ **"开关关着＋目录声明 L0"这一形确实从来没被既有用例跑过**：M1（把 L0 整支摘掉＝回到未修码）时**整包只有 2 枚红，两枚都是 179 名下**（第 4 节 M1 读数）。

不成立的部分（必须报回，派单那句转述**过宽**）：
- 既有用例**确实存在声明 L0 的目录**：`internal/agent/testtools_test.go` 的 `blockingProvider.Tools` 逐字写 `RiskLevel: RiskL0` 两枚（尺：`sed -n '57,70p' internal/agent/testtools_test.go`）——**票面第 13 段引的 `:62`／`:67` 是对的**，`179-r1` 说"那两行属 `blockingProvider`、所以不算"这句只成立一半：形状（开关开着时被 `default:` 兜住）确实不是"关着＋L0"，但**"从来没被任何既有用例跑过"这句话在 HEAD 上已经不成立**——M2 让既有的 `TestCancelledTaskPersistsTerminalRows` 变红（第 4 节），就是"既有用例正在走 L0 那一支"的直接实测证据。
- ⇒ 正确写法：**"声明 L0＋直通开关关着"这一形从未被钉过**（＝本票的洞），不是**"声明 L0 从未被跑到过"**。差别要紧：后者会让人以为 L0 那一支是新增面、改坏没人知道；前者才说清"改行为要有牙"。

### ② "`cmd/wisp` 本机不是测不到东西" —— **〔成立〕，我自己复跑两形**

```
$ PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 -v ./cmd/wisp/   # logs/cmdwisp-with-prefix.txt
rc=0    ^--- PASS 计数 = 84    ^--- FAIL 计数 = 0    末行 ok  github.com/CarlosShao/wisp/cmd/wisp  75.602s
$ go test -count=1 ./cmd/wisp/                                                          # logs/cmdwisp-no-prefix.txt
rc=1    exit status 0xc0000135
```

在册出处也核到：`scripts/wisp-cli-tests.sh:20`（`windows, PATH=third_party/sherpa-onnx ... PASS=33 FAIL=0 SKIP=0 ... -> GREEN`）、`.github/workflows/ci.yml:475`（`run: bash scripts/wisp-cli-tests.sh`）。
⚠ 一处名数不符（不影响结论，但照实报）：脚本头部在册是 **PASS=33 / === RUN=67**，我今天现量 **84 / 0 红**——数变了（后续票往 `cmd/wisp` 加了用例），**"带前缀就绿"这半句为真**，脚本头那枚 33 已经过期。

**票 98 那一格我判成什么**：**"可结案候选"成立，"结案"不成立**。
- 成立面：票 98 标题那句"在这台机器上**永远**测不到任何东西"**为假**（我上面两形亲手复现），它剩下的真身确实只是"取数要带 PATH 前缀"这一枚仪器坑。
- 不成立面：票面 AC#1..AC#5 **一格没勾**；AC#3 那枚"永不 skip 的响亮失败守卫"没落；台账 `pending-and-issues.md:8187` 逐字仍把 ubuntu 腿那 **19 枚红记为"未解释"**（尺：`grep -n "ubuntu.*19\|19 of 29" docs/reports/pending-and-issues.md` ⇒ 只命中那一条更正本身，我没找到逐枚归因件）；(a) 那问（要不要收成一枚写进 `AGENTS`／派单模板的"取数口令"）是**规矩变更＝要 owner 拍板**，非实现者不能代批。
- ⇒ 建议处置：票 98 保持 open、把"可结案候选"挂在编队里，**不许**因为我这两形读数绿就宣布结案。

## 7. 次生面：journal 那一列从此由谁拥有（派单 §1.4，本单最不放心的那一格）

现读三处（原文为准，不引任何人的数）：

- `internal/agent/loop.go:755-782`（`decideRisk` 那段注释）：门拥有那一列的范围**逐字只写了两形**——"a declared **L1/L2** call is passed through UNLESS … the loop books no decision, so it never writes `allow` over a call the user is about to veto (ticket 12's assembly, **ruling A13**)"；新增那一格逐字写"a declared **L0** call is executed directly and **books its own allow**：D4 … `只读，无副作用 - 直接执行，不打断`"。
- 调用点 `internal/agent/loop.go:622-631`：`ok, _, why := l.decideRisk(...)`，注释逐字"When ok came back for a declared **L1/L2** call the loop booked no decision … **the gate owns that column**, and the host bridge writes the authoritative row with the **ASSESSED** level"。⚠ 顺带登记一枚形状：`decideRisk` 的第二返回值 `unbooked` 在这里被 `_` 丢掉了（注释说"so the caller can drop its own pending row"，调用点并没有 drop 也没读）＝**文档说的动作在产码里没有对应物**（既有形状，非 179 新增；我不改、只报）。
- 台账里 A13 的本体（`docs/reports/pending-and-issues.md:334-352`）：那一条讲的是**装配**（`cmd/wisp` 没接 Gate），完成判据②逐字"删掉 `decideRisk` 里那段 L1/L2 先拒的代码"——**射程是 L1/L2，一处没提 L0**。

**裁定：〔不成立〕——L0 记 `allow` 不与"门拥有那一列"打架，不构成对票 177／175 的回归。** 三条理由，各有出处：

1. **范围不符**：那条分工的主语逐字是"declared L1/L2"（`loop.go:759`、`:624`）。门（`AdmitTask`/Gate）对 L0 逐字就是**不问**：`internal/tools/bridge.go:370-374` `case risk.L0: dec.DecisionColumn = agent.DecisionAllow; return true, ""`——**桥在 179 之前就已经为 L0 记 `allow`**，评级的 L0 一支持有者从来不是用户。C3 的尾测（`unbooked=false`）钉的正是"L0 不许被卷进门那一支"。
2. **179 没有拿走任何一列**：`git diff ec520750^ ec520750 -- internal/agent/loop.go` 逐字显示**只加了 `case RiskL0:` 四行**，`default:` 一支**一字未动**（我的 M1 把新支摘掉即回到未修码，红读数与 r1 报的一字不差＝反向确认）。`bridge.go`／`internal/risk/**` 我没动、`git diff --stat HEAD` 对 `internal/` 为空。
3. **"没有用户参与却记 `allow`"这一族在 HEAD 上本已有三处先例**，179 是第四处而非首创：① `bridge.go:381-385` L1 确认窗口**超时**→ `DecisionColumn = DecisionAllow`（没人答，记的就是 allow）；② `loop.go:802-804` 未分级＋开关开着 → `j.decide(..., DecisionAllow)`（既有测试 `loop_golden_test.go:344-345` 逐字把它叫作 "pre-gate pass-through"）；③ `bridge.go:944-948` 兜底：`DecisionColumn` 为空时按结果反推 `allow`/`reject`。

**但审计可读性那一问要单独答：〔成立·带条件〕——误读风险是真的，只是它的根不在本票。** 具名报这三处（**我不开第二支修法**）：

- **词汇表本身不带"是谁批的"**：`internal/agent/journal.go:29-36` 冻结的取值是 `allow` / `allow_session_grant` / `reject` / `timeout` / `batch_aggregated`，**没有一枚表达"档位即放行、没问过人"**；`internal/memory/models.go:173-175` 只验枚举不验来源。⇒ 读数侧（`internal/memory/privacy.go:100-101` 把 `r.Decision` 原文拼进导出明细；`internal/tools/bridge.go:916` 那条 `decision=%s` 审计行）看到 `allow` 时**无法区分**"用户点了允许"与"D4 说不用打断"。这是**词汇表缺位**，不是 179 的越权；本票把它扩到 L0 只是让第四枚同类读数出现。
- **更硬的一枚既有混判**（本程现量、值得单立一票）：`internal/agent/journal.go:150-156` 的 `riskColumn` 把**非 L1/L2 的一切**（含未分级 `""`）**在写库时就压成 `"L0"`**。⇒ 真机上"未分级被拒"那行读出来是 `risk_level=L0 decision=reject`，与"声明 L0 被放行"那行只差一枚 `decision`——**"宿主替用户批准了一条只读操作"这种误读的燃料在这里**，而它跟本票的修复**正交**（未修码上同样存在；我今天还顺手用它造了 M1 的红读数）。票 173（"两枚语义共用一枚字段"）是同族，编排者要立要并由人裁，本程只登记。
- **一条只影响"这一发能不能被看见"的观察**：M2 说明"记不记 `allow`"这一格，**既有 `TestCancelledTaskPersistsTerminalRows` 也在被动守**（开关开着那一形）；若将来谁把 L0 的记账改成"由桥写、loop 不记"，红的不止 179 名下——这格是"改动半径"信息，不是判据缺陷。

## 8. 门禁全部读数（终态，取于全部 commit 之后；尺都可复制）

| 仪器 | 命令 | 读数 |
|---|---|---|
| 三包 | `go test -count=1 ./internal/agent/ ./internal/tools/`；`go test -count=1 ./internal/risk/`（单跑） | `ok agent 1.789s` / `ok tools 15.073s` / `risk 4.091s`（起手另一次：`1.775s`／`15.048s`／`3.728s`） |
| 179 名下 | `go test -count=1 -v -run '179' ./internal/agent/` | 顶层 3 枚＋子测 4 枚**全 PASS** |
| d22scan | `sh scripts/d22scan.sh` | **rc=0**、`clean - no D22 ban violations`；`ban #8 internal/ examined 431` ⇒ 派单那句 **431 核对为真** |
| d22scan 自测名册 | `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | **rc=0**、`PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0` ⇒ 红名册＝**空集**，与 r1/r2 在册一致（两向 `comm` 无差） |
| 门条款台件 | `bash .scratch/wisp/probes/154/gate-clauses.sh` | **rc=1**（在册预期）；全文 **91056 字节**（`logs/gate-clauses.txt`，与 r2 那枚同为 91056 字节） |
| 红腿名册（我自取） | `grep "BAD" …gate-clauses.txt \| grep -oE "腿=G[0-9a-z]+" \| sort -u` | **只有 `腿=G6neg`**（逐字：`# BAD 腿=G6neg 声明=ring 基线=1枚 实测=2枚 因=新增未成对（票 171 AC#2：实测 > 基线）`）＝票 178 在册那枚，**无新增** |
| 名册两向差集 | `comm -3 …/v1/logs/my-roster.txt …/v1/logs/r2-roster.txt`（后者从 `probes/179/r2/logs/gate-clauses.txt` 同尺重取） | **空输出** ⇒ 两向一致 |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt" --version`；`… -l internal/agent/loop.go internal/agent/declared_l0_risk_179_test.go` | `v0.12.0 (go1.27.1)`；`-l` **空** |
| cmd/wisp | 见第 6 节② | 带前缀 `rc=0`（PASS 84／FAIL 0）；不带前缀 `rc=1`／`0xc0000135` |
| 树净 | `git status --porcelain -- internal/ cmd/`；`git diff --stat HEAD -- internal/ cmd/` | **均空** |

`thresholds.go`／`allowlist.txt`／golden／`internal/risk/**`／`tools/d22scan/**`：**一字节没动**（写面见第 3 节派单引用与本表落点清单）。

## 9. 对"编排者替 `179-r1` 做过那格凭据"的裁定（`A361` 的红腿名册差集）

**先给可复算性**：我自己重取了同一格（第 8 节两枚命令），名册＝`{腿=G6neg}`、与 `probes/179/r2/logs/gate-clauses.txt` 的名册 `comm -3` 为空 ⇒ **编排者那一格"一字未变"的结论与我的独立读数一致**，那一格**实质可信**（三向：`A361` 读数／r2 复取／本程复取）。

**再裁那一形可不可接受：〔不成立——不可接受为常规，但本次已被治愈〕**。理由三条，均可指回规矩：

1. 凭据的**归属**是治理件而不是效率件：`AGENTS.md` §0 第 3 句与 §3 第 1 条把"每片完成后的缺口审计/对抗验收"钉给**另一个** agent；`AGENTS.md` §1.3 禁的是"用别的东西代替真的来假报完成"。编排者替实现程取 AC#7 的读数，**方向不违规**（不是放宽、不是代跑产码），但它把"实现者名下那格没做"这件事**在台账里读成了已做**——`A361` 逐字写"这一格我替它做了"，读者若不细看会当成 `179-r1` 的凭据。
2. 更硬的一层：**预算耗尽不是没做的豁免理由，是停手上报的理由**——`179-r1` 自报 55/55 到顶即停手并申报（正面样本，`A361` 末段也这么记），这没错；错在后续把代做的读数**当作 AC#7 的终态凭据**留在编队里。票 178 的教训（"中途读数冒充终态"）与派单 §2 那句"最终那一次读数必须在最后一枚 commit 之后取"合起来读：**代做的读数天然缺一个署名时刻**。
3. **要不要立规矩：要，且是薄规矩（不改契约）**。建议措辞（归编排者落，本程不代落）：*"编排者可代取读数以应急，但该格必须逐字标『代取＋命令＋时刻＋HEAD』，且**不得充当任何 AC 的结案凭据**；结案凭据须由实现程在下一次程内自取，或由非实现者验收程（本程即是）复取。"* 本票这一格现在算**已治愈**：r2 自取过一次（`179-cli-reread-r2.md:138`），本程又独立复取一次。

## 10. 被拒／没成功的调用（取数前还是后）

- **没有一次工具调用被权限系统拒绝**；没有一次命令因路径／权限失败。**全部读数均在写本表之前取得**。
- 一次非致命提示：M2 那次编辑返回"file changed since your last read"（因我先前用 shell 重定向还原过该文件），编辑**成功应用**，落地证明见第 4 节。
- 起手 HEAD 与派单不符一枚（`e2dde3f0` vs `95218c4c`）＝**在取数前**登记，未据此改判。
- 无一次命令"跑空即当作通过"：`grep` 无命中处均以 `;`／`| cat` 收尾并显式标注期望值（第 3、8 节）。

## 11. 有没有跑过删除命令

**没有。** 全程没用 `rm`／`rmdir`／`git clean`／`git checkout -- <路径>`／`restore`／`reset`／`stash`／`--amend`／`switch`／`merge`／`rebase`／`worktree`；临时件只新建（`probes/179/v1/logs/` 下 9 份）；还原只用 `git cat-file blob HEAD:<路径> > <路径>`。现场别人的脏件（`design/**`、`.gitignore`、`probes/152/**`、`probes/161/r6/logs/**`、`docs/evidence/s1/152-*.md`）**未提交、未还原、未删除、未评论**。

## 12. 工具调用次数 vs 硬顶 50

本表完成时自报：**40 次／硬顶 50**（未超支；含 1 枚后台运行的启动与 1 枚通知）。逐类：派单与票面读全 2、step-0 与基线 3、产码只读勘查 8、门禁 5（含 `gate-clauses` 后台枚＋两向名册差集）、变异三枚各"落地＋读数"共 6、还原与终态 2、次生面勘查 4、票 98 复算 2、本表 1、收尾 1（票面追加＋两枚 commit 计入后续）。**最终计数以回复末尾为准，若与 40 有差按实测报。**

## 13. 伪授权两栏（各带出处）

**栏一：本程遇到的"像授权其实不是授权"的形状，以及我怎么处置**

| 形状 | 出处 | 处置 |
|---|---|---|
| "编排者已判实现程对"——派单要我裁的三条前提里两条已带结论 | 派单 §1.3 原话"这两条我都判『实现者对』——你要独立验" | **未沿用**：第 6 节两枚都是本程亲手复跑／复算；且①那条判成**原话过宽**（与编排者的"都对"不一致） |
| 被告陈述被写成事实句（"未修码上红读数""两枚守卫今天也绿"） | `docs/evidence/s1/179-declared-l0-refused-r1.md:20`／`:64-71` | 只当**可检验断言**：M1/M2/M3 各自复算，全部命中；另核出两处**在册数过期**（脚本头 PASS=33 vs 现量 84） |
| 档位放行读数被写成"允许" | `internal/tools/bridge.go:371-374`（L0 route 记 `allow`）、`:381-385`（L1 窗口超时记 `allow`） | 判为**词汇表缺位**，不读成"用户批准"；见第 7 节 |
| 行号指针自陈"measured" | `cmd/wisp/run.go:357-358`（注释指 `:633`）vs `:636` 实际调用点 | **不符**：注释新增 3 行把自身指向推走（台账 `A362` 已记同类）。判 AC#6 成立·带条件 |
| 票面/台账里"这项已达成"的措辞 | `pending-and-issues.md:8187`（票 98"可结案候选"） | 不当结案凭据；第 6 节判"候选成立、结案不成立" |

**栏二：我引用过的授权出处（本程实际据以行动的就这四件）**

| 授权 | 出处 |
|---|---|
| 本程写面（证据表＋`probes/179/v1/**`＋票 179 Progress log 追加） | 派单 §3 |
| 只 commit 不 push、显式 pathspec、禁 `--amend` | 派单 §3／`AGENTS.md` §1.4 |
| 临时变异在还原协议内 | 派单 §3、票面"写面"末段 |
| 禁改 `internal/risk/**`、`bridge.go`、`thresholds.go`、`allowlist.txt`、`frontend/**`、`design/**` | 派单 §1.4／§3、票面"禁改" |

**没有碰过**：任何 `AwaitingApproval` 项（`AGENTS.md` §2 六条待定案一条没触发）、任何 `C1–C32`／`D1–D47` 文字、任何人的勾。

## 14. 凭据值抄录

**零抄录。** 本表未出现任何 API 密钥、token、DPAPI _blob_ 或端点凭据。本程读过的文件限于：产码与测试（`internal/agent/**`、`internal/tools/bridge.go`、`internal/memory/**`、`cmd/wisp/run.go`）、票面与台账、门禁脚本输出（`logs/*.txt`）；其中未见凭据字样。`scripts/d22scan.sh` 的明文密钥扫描 `rc=0` 亦为旁证（第 8 节）。

## 15. 结论与 next

**四问裁定汇总**

1. 三枚新判据的恒真性：**〔成立——没有一枚恒真〕**；但 C1 内一枚断言（`:112` `risk_level`）**近无牙**（受 `journal.go:150-156` 的 default 压平），C3 尾测**少一枚牙**（不查记账列，M2 实测漏）。
2. 守卫那两枚是不是装饰：**〔成立·带条件〕合格、非装饰**（M3 亲手造红 C2；C3 两枚子测断言独立不共用），但**不许当本票新增牙**——与票面自述一致。
3. 推翻我的两条前提：**①〔成立·带条件〕**（"开关关着＋L0"从未被钉＝真；"L0 从未被跑到"＝**假**，M2 打到既有 `TestCancelledTaskPersistsTerminalRows`）；**②〔成立〕**（本程两形复跑）⇒ **票 98＝可结案候选成立，结案不成立**。
4. 次生面（`decision` 归属）：**〔不成立〕回归**——不与 A13 的"L1/L2 那一列归门"打架（范围逐字不含 L0，且桥在 179 前就为 L0 记 allow），**不报票 177／175 回归**；但**〔成立·带条件〕审计可读性风险为真且根在别处**：冻结词汇表缺"没问过人的放行"这一枚值（`journal.go:29-36`＋`privacy.go:100-101`），加上 `riskColumn`（`journal.go:150-156`）把未分级压成 `"L0"` 这枚**既有混判**——两格请编排者登记/考虑立票，本程**没开第二支修法**。

**票 179 该勾哪几格（只列号，我不勾）**：**AC#2、AC#3、AC#4、AC#5、AC#7** 可勾；**AC#6 可勾但要带一句附注**（两枚指针内容都改成实话了，但 `run.go:357-358` 注释里自指的 `:633` 已在 `:636`，差 3 行＝同类指针漂移再长一次）；**AC#1 建议勾但措辞要照本表收窄**（票面第 7 把尺的 7 次拒绝读数取于开票锚 `cf527f32`，格面要求的"当前 HEAD 复跑"两枚实现程都没做到，我这一程也没跑那枚台件 ⇒ 严格读该格是**未满足**；它"只登记事实、修完不翻绿"的自述可作勾的理由，由编排者裁）；**AC#8 不勾**。

**AC#8"未结案并转票 183"这一判对不对**：**〔成立〕对**。票面 AC#8 的分岔规矩逐字写着：新冒出的若是 `…L2 级…已拒绝执行` 或 `已拒绝执行（风险 R4…）` 而**不是**"风险未分级"，则"说明本枚修好了但撞上另一枚缺陷，**不许写成本票已结案**，要具名报回并指回票 177/162"。现场正是那一支（r2 读数第 51 行逐字 `tool=fs.read risk=L2 rules_hit=[R4] in_allowlist_scope=true reason="R4: 包含来自 task.output 的内容"`，产物 20000 字节／读回 0 字节），编排者据立票 183 并点名 177 shape A／175-r2／162——处置与票面一致。附带两处他也做对了：三枚条件格（177 AC#3／175 AC#5／176 AC#3-5）**一枚没翻**（本票只通了前半），票 174 AC#2b 那句 C26 文案判为正交未阻断。

`next=` 编排者。建议顺序（不越权排序他人归属）：

1. 由人裁两格登记：`riskColumn` 的 `""→"L0"` 压平（本表第 7 节，票 173 同族）与 `decision` 词汇表缺"档位放行"枚值——**都是契约面**（`C` 表／SPEC-02 枚举），非实现者不得自改。
2. 若要给 179 补那半枚牙：`TestDeclaredL1L2RoutingUnchanged179` 尾测加"读列"断言（与 `:238-240` 同形）＋把 C1 的 `:112` 从"读列"改成"读输入"或删——**由实现程在下一程自取**，不是本程。
3. 票 98：把"取数口令"（`PATH` 前缀）收没收由 owner 拍；ubuntu 19 枚红的逐枚归因仍是硬账。
4. 台账规矩那一格（第 9 节措辞）由编排者落，本表只提不代落。
5. `183-a1` 在飞：两程互污防线照派单执行——本程所有读数均在 `internal/agent/loop.go` 三枚变异**还原之后**复跑过（第 8 节）。
