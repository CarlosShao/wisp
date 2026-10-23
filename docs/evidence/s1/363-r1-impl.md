# 票 363 落地腿 `363-r1` —— 交件（实现者）：日志在 `L0` 与 `Deny` 两个方向都说真话

- **腿**：`363-r1`（**写腿／实现者**）。形＝**〔丙〕两者都做**（票面 `AC#1` 已拍／已翻）。
- **起手锚**＝`76ae3bbb`（派单存档那一枚）；**我的产码提交**＝`d59af9a`。
- **写面**：`internal/risk/mode.go`、`internal/tools/bridge.go`、两包测试、本件。★**同一笔里被卷进一枚非写面路径**，如实登记在 §11①（需编排者拍，我不自行 amend/reset）。
- **环境**：`GOFLAGS=`、`unset GOTMPDIR TMPDIR TEMP TMP`、⛔ 未用 `CGO_ENABLED=0`；`go1.27.1 windows/amd64`；`gofumpt`＝绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`。

## §0 开工前三条（⛔ 跳步）

1. **写面干净**：`git status --porcelain -- internal/risk internal/tools` **空**（首／尾各一次）。
2. **票面前提新鲜度（逐条现跑，`sed -n` on 工作树；⛔ 抄票面）**：`bridge.go:311` `sil := mode.Screen(verdict)` ✓／`mode.go:163-165` 逐字 ✓／`bridge.go:323` 条件＋`324-325` 日志 call ✓／`b.route` 在 **352**、`// --- routing: L0 pass…` 注释在 **351**（票面并进 352，我复认**偏差成立**）／`route` 的 `case risk.L0`＝**433-435** ✓／`assessor.go:41/45/46` 逐字 ✓／`mode.go:143-146`（`Kept` 契约）逐字 ✓／`mode.go:174-175`＋`185-188` 两处 `return Silenced{Level: L0, …}` ✓。
   `git log 46e61076..76ae3bbb -- internal/risk internal/tools` ⇒ **0 笔** ⇒ 普查窗口与本腿起手是同一份字节。
3. **复跑基线**（修前 `-count=1 -v` 两包）：`RUN=393 PASS=390 FAIL=2 SKIP=1`。红名**逐名**：`TestResolvePerCallBudget`（时序预算 `1.013 ms/op > 1.000 ms`，本机负载）、`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（AST 扫 `internal/panel/**` 的 `Result.Canonical` 消费者）。★两条**在我动手之前就红**，不是本腿造成（见 §9）。

## §1 ★先量：两向「今天是什么读数」

仪器＝`go test -overlay`：探针在**仓外** `D:\tmp\wisp\20261023\363r1\zz_t363r1_probe_test.go`，经 overlay 映射成 `internal/tools/zz_t363r1_probe_test.go` 的**虚拟路径**参与编译 ⇒ **仓内 `.go` 字节 0 枚**。组合与 `cmd/wisp/run.go:764-788` **同形**，唯一差异＝`Logf` 指捕获器。**卡数＝`Gate` 的两个回调（`gate.go:110-117`）被调用几次**，⛔ 数 UI／tick。

**〔L0 向〕真 `fs.read`（真 D34 声明 `FSReadDecl()`、真 R1–R9），三档逐字同一句（只有 `mode=` 变）**：

```
ASSESSED level=L0 rules=[] reason="无规则命中（L0 直接执行）"
EMITTED level=L0 mode=ask_every_step silenced=false kept="无法识别的风险级，fail-closed 保留询问"
CARDS L1_WINDOWS=0 L2_APPROVALS=0 TOTAL=0
LOGVERBATIM|tools: MODE-REDLINE mode=ask_every_step still asks tool=fs.read rules=[] because=无法识别的风险级，fail-closed 保留询问
LOGVERBATIM|tools: call kind=success task=task-1 corr=corr-1 tool=fs.read risk=L0 decision=allow outcome=success rules_hit=[] in_allowlist_scope=true grant_id=0 refusal=none reason="无规则命中（L0 直接执行）"
EXEC err=<nil> isError=false text="hello"
```

⇒ **同一发**：日志说"会问"／`CARDS TOTAL=0`／`decision=allow` 且**已执行完并回结果**。

**〔Deny 向〕`fs.read` 打 `%APPDATA%\wisp\config.toml`（`ask_every_step` 与 `auto_approve` 同形）**：

```
ASSESSED level=Deny rules=[R2 R3] reason="R2: 目标路径在授权目录之外: …\AppData\Roaming\wisp\config.toml; R3: 命中 A 档敏感路径，拒绝（不可豁免）: …"
EMITTED level=Deny mode=auto_approve silenced=false kept="拒绝级 verdict（Deny 级不受任何模式影响）"
CARDS L1_WINDOWS=0 L2_APPROVALS=0 TOTAL=0
LOGVERBATIM|tools: MODE-REDLINE mode=auto_approve still asks tool=fs.read rules=[R2 R3] because=拒绝级 verdict（Deny 级不受任何模式影响）
LOGVERBATIM|tools: call kind=denied task=task-1 corr=corr-1 tool=fs.read risk=L2 decision=reject outcome=error rules_hit=[R2 R3] in_allowlist_scope=false grant_id=0 refusal=redteam-denied reason="…"
EXEC err=<nil> isError=true text="该目标被绝对禁止访问（R3 A 档，任何授权都不可豁免）：…"
```

## §2 〔乙〕`Screen` 认得 `L0`（`internal/risk/mode.go`，+25/−1）

- 新支（在 `Deny` 之后、`!= L1 && != L2` 之前）：`if d.Level == L0 { return Silenced{Level: L0, Mode: m} }`；同段注释给出判据：`L0` 是 `assessor.go:46` 的**枚举零值**、`:45` 的 `read-only, no side effects - execute directly`、`Level.String()` 有名字，**且同函数两个出口自己就在发 `Level: L0`**。`mode.go:155-156` 那句自述改成 `Deny, L0 and an unknown level are handled first`。
- **`Kept` 留空**：契约（`mode.go:143-146`）只把它留给"被红线**留住**的询问"，`L0` 没有问题可留 ⇒ 每一行只读 `tool_call` 不再带一条不存在的红线。
- ★**`Silenced` 留 `false` 是判断，写清取舍**：`mode.go:141` 说"true when the caller must NOT ask"，`gate.go:43-47`（同名字段在 `tools.Decision` 上）说"true when the mode removed a question this verdict would otherwise have asked"。`L0` 上二者相反。选 `false` 的理由：置 `true` 会把"mode 移走了一个询问"写进**每一行只读 `tool_call`**，并把 bridge 引到 `MODE-SILENCE` 那条 **"a silenced question"** 的日志上——那是**把范畴错误换个方向再犯一次**。"不询问"那一跳仍由调用方 `route` 的 `case risk.L0` 独占（`bridge.go:433-435`，**一字未动**）。
- **fail-closed 支没删、只收窄**：真正不可能的级（`Level(42)`）照旧 `Kept="无法识别的风险级…"`、`Silenced=false` ⇒ 由 `route` 的 `default` 拒（正控用例钉住）。
- **跨包复核**：`Screen` 的语义表住在 `internal/perm/store_test.go`（**非我写面**）⇒ 我复跑 `go test -count=1 ./internal/perm/` ⇒ `ok`（rc=0），该表 10 行 × 3 档 **无一转红**。

## §3 〔甲〕日志说实话（`internal/tools/bridge.go`，+13/−0）

- `case sil.Kept != "":` 之前**新增一支** `case sil.Kept != "" && sil.Level == risk.Deny:`，打：
  `tools: MODE-REDLINE mode=%s refuses tool=%s rules=%v because=%s (a Deny never reaches a gate: no card is raised)`。
- **为什么这一支必须独立存在**：`Kept` 非空 ≠ 有询问。`route` 的 `case risk.Deny`（`bridge.go:437-439`）直接 `reject`、**不碰 gate** ⇒ 0 张卡。旧文案在这里**反方向说假话**（拒绝了却说会问），而"乙"修不到它（`Deny` 本来就是 `Screen` 认得的级）。
- **`Kept` 契约语义一字未动**：`Deny` 的 `Kept` 文本原样（`拒绝级 verdict（Deny 级不受任何模式影响）`），改的只是**日志怎么念它**。⇒ 派单的停手条件①**未触发**。
- **"still asks" 留给真的会开卡的那一支**：`L2` 被 `ask_high_risk` 或红线留住时，`route` 走 `PendingApproval` ⇒ 一张审批卡 ⇒ 那时"会问"是真的（§5／§7 读数）。`L1` 未被静默时 `Kept==""`，**根本进不了这一支**（`ask_every_step` 唯一一支 L1 未静默、且 `Silenced=false`、`Kept=""`）⇒ 会话授权那条路也印不出这句。

## §4 ★修后两向读数（同一仪器、同一组合）

**〔L0〕三档**（`EMITTED … kept=""`，**没有任何 `MODE-` 行**）：

```
EMITTED level=L0 mode=ask_every_step silenced=false kept=""
CARDS L1_WINDOWS=0 L2_APPROVALS=0 TOTAL=0
LOGVERBATIM|tools: call kind=success task=task-1 corr=corr-1 tool=fs.read risk=L0 decision=allow outcome=success rules_hit=[] in_allowlist_scope=true grant_id=0 refusal=none reason="无规则命中（L0 直接执行）"
EXEC err=<nil> isError=false text="hello"
```

**〔Deny〕两档**（同一发：0 卡、拒绝、日志说"refuses"）：

```
LOGVERBATIM|tools: MODE-REDLINE mode=auto_approve refuses tool=fs.read rules=[R2 R3] because=拒绝级 verdict（Deny 级不受任何模式影响） (a Deny never reaches a gate: no card is raised)
CARDS L1_WINDOWS=0 L2_APPROVALS=0 TOTAL=0
LOGVERBATIM|tools: call kind=denied … risk=L2 decision=reject outcome=error …
EXEC err=<nil> isError=true text="该目标被绝对禁止访问（R3 A 档，任何授权都不可豁免）：…"
```

## §5 ★四档行为无回归（各一发；`fs.*` 真工具、真判定器；卡数＝两个回调）

| 档 | 输入（真工具） | 判定 | **卡数** | 结果 | 日志那一支 |
|---|---|---|---|---|---|
| **`L0` 仍放行** | `fs.read`（临时根内） | `L0` `[]` | **0/0** | `decision=allow outcome=success` `text="hello"` | **无 `MODE-` 行**（不再说假话） |
| **`L1` 仍窗口** | `fs.write` 新文件／`ask_every_step` | `L1` `[R1]` | **窗口=1**，审批=0 | `decision=allow` `已写入` | **没有** `MODE-REDLINE`（正控①） |
| **`L2` 仍审批** | `fs.write` 覆盖／`ask_high_risk` | `L2` `[R1 R8]` | 窗口=0，**审批=1** | `decision=allow` `已写入` | `MODE-REDLINE … still asks … because=ask_high_risk 只静默 L1，本判定是 L2`（**真的**） |
| **`L2`／`auto_approve`** | 同上 | `L2` `[R1 R8]` | **审批=1** | `decision=allow` `已写入` | `… still asks … because=R1 工具自身声明的 L2（不可逆自认，PLAN.md:1629）`（**真的**） |
| **`Deny` 仍拒** | `fs.read` 打 `%APPDATA%\wisp\config.toml` | `Deny` `[R2 R3]` | **0/0** | `isError=true` `该目标被绝对禁止访问…` | `refuses … (no card is raised)` |

⇒ **放行语义、窗口、审批、拒绝四档一个没动**；`L1` 被 `ask_high_risk` 静默那一发仍打 `MODE-SILENCE … (a silenced question, not an allow: no user click happened)`、0 卡。

## §6 ★反向判据三读数（`-overlay`，**工作树零突变**）

两枚突变体建在**仓外**：`D:\tmp\wisp\20261023\363r1\mut\mode_m1.go`／`bridge_m2.go`（`diff -u` 逐行可核：只少掉那一跳）。

| 读数 | 摘掉哪一跳 | rc | **逐名红行** |
|---|---|---|---|
| **M1 变异红** | `mode.go` 的 `if d.Level == L0 {…}` | **1** | `TestTicket363ScreenRecognizesL0`（3 档 + 越界档，4 条红：`Screen(L0).Kept = "无法识别的风险级，fail-closed 保留询问", want empty`）、`TestTicket363FourLevelsKeepTheirRoute`（3 条 `L0/…: Kept = …`）、`TestTicket363L0ReadNeverClaimsAQuestion`（3 档 × 5 条：`ModeKept=…`／`'still asks'`／`无法识别的风险级`／`neither a silenced question nor a kept one`／`claims a question while the call raised 0 cards`） |
| **M2 变异红** | `bridge.go` 的 `case sil.Kept != "" && sil.Level == risk.Deny:` | **1** | `TestTicket363DenyIsRefusedNotAsked`（2 档 × 3 条：`'still asks' on a Deny`／`does not name the refusal`／`claims a question while the call raised 0 cards`）；★`L0` 那枚**仍绿** ⇒ 两跳**互不代偿** |
| **还原绿** | 不挂 overlay（=盘上 `d59af9a` 的字节） | **0** | `TestTicket363*` 6 枚全 `--- PASS`（`04-new-tests-after-fix.txt`／`14-final-tests-main.txt`） |

★**"今天响、不修不响"**：同一批用例在修前（`03-new-tests-before-fix.txt`）红 4 枚、修后全绿；两枚**正控**用例在修前**就是绿的**且修后仍绿（`TestTicket363UnknownLevelStillKeepsItsQuestion`、`TestTicket363OtherLevelsKeepTellingTheTruth`）⇒ 牙不是恒红、也不是恒绿。

★**还原核**（★先问"谁算的、什么算法、几位"）：`git hash-object` ＝ **SHA-1（40 位十六进制）**，⛔ 与 `sha256sum`(64) 混比。四枚逐枚"工作树 blob ⇄ `d59af9a:` 里的 blob" ⇒ **4/4 `same=YES`**（`4de3bdac…`／`69822c75…`／`93ed19fa…`／`250c56db…`）。★**本来就没有"还原"这一步可核**：两次突变全程 `-overlay`（替换文件在 `D:\tmp\…\mut\`），**工作树与索引在突变前后各查一次均为空**（`git status --porcelain -- internal/risk internal/tools`）⇒ **无突变窗口、无窗口期提交**。

## §7 正控（硬闸）

- **`L1` 不动那行仍没有**：`ask_every_step` ＋ `fs.write` 新文件 ⇒ 窗口 **1**、日志里 `MODE-REDLINE` **0 行**、`MODE-SILENCE` **0 行**。
- **`L2` 说的仍是真话**：两档各 1 张审批卡，`because=` 逐字如上；用例断言 `still asks` 与 `卡数>=1` **同发**（`t363AskClaimsNeedACard`）。
- **未知级仍 fail-closed**：`Level(42)` ⇒ `Kept` 含"无法识别"、`Silenced=false`、`Level` 原样返回。
- ★**不变量**（两处改动共同承重）：**任何一行说 `still asks`，同一发就必须真开过卡**。修前 `L0`／`Deny` 各违反一次，修后全部满足。

## §8 格式门禁（甲乙两尺）

- **甲（判据）**：`git archive d59af9a` 到仓外 `D:\tmp\wisp\20261023\363r1\export-d59af9a`，跑 `/d/work/base/gopath/bin/gofumpt.exe -l .` ⇒ **20 行、全是 `.scratch/**`、非 `.scratch` ＝ 0**（`/` 归一后判，⛔ 反斜杠裸判）；再跑 `gofumpt -l internal/risk/mode.go internal/risk/ticket363_screen_l0_test.go internal/tools/bridge.go internal/tools/ticket363_test.go internal cmd tools scripts frontend` ⇒ **0 行**、rc=0。
- **乙（只作诊断）**：工作树同四枚 ＋ 同上目录集 ⇒ **0 行**、rc=0；`gofmt -l` 同四枚 ⇒ 0；四枚 `tr -cd "\r" | wc -c` ⇒ **0 0 0 0**（`*.go text eol=lf`）。

## §9 门禁 rc

| 门 | 对象 | rc | 读数 |
|---|---|---|---|
| `go build ./...` | 主检出 | **0** | — |
| `go vet ./internal/risk/ ./internal/tools/` | 主检出 | **0** | — |
| `go test -count=1 ./internal/risk/ ./internal/tools/` | 主检出（终态） | **1** | `RUN=399 PASS=396 FAIL=2 SKIP=1`（四数分开） |
| 同上，红名逐名 | — | — | `TestResolvePerCallBudget`（**时序**：`1.013 ms/op > 1.000 ms`，负载敏感，同一份字节另两发为绿）＋`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（**既存红**，断言对象是 `internal/panel/**`，本腿零碰）⇒ **红名集合与基线逐名相同、无新增红** |
| 同上，另一发（负载较轻） | 主检出 | **1** | `RUN=399 PASS=397 FAIL=1 SKIP=1`，红只剩 `TestC26…` ⇒ 那枚时序用例**逐发漂**，故本件按**红名集合**判，⛔ 按 PASS 数判。四数对账：`399 = 393 + 6`（本腿 6 枚新用例）、`PASS = 390 + 6 + (0 或 1)`、`FAIL = 2 − (0 或 1)` |
| `go test -count=1 ./internal/perm/` | 主检出（跨包复核 `Screen`） | **0** | `ok` |
| `sh scripts/d22scan.sh` | **本腿提交导出树** | **0** | `clean - no D22 ban violations`；`ban #8 internal/=623`（基线值 +2＝本腿两枚新用例 ×1 各）、`ban #7 internal/tools/=25` |
| `sh scripts/d22scan.sh` | 主检出 | **1** | ★**不是本腿**：唯一违规＝`internal/memory/schema.go:25`→`:26`（196 腿正在改，行号随它漂）引 `docs/evidence/s1/196-r1-impl.md`（文件尚不存在）。`git status --porcelain -- internal/memory` ⇒ ` M schema.go`(+138/−4)／` M dao_tasklog.go`／` M models.go`；`git show HEAD:internal/memory/schema.go｜sed -n 25p` ⇒ `value TEXT NOT NULL`（**HEAD 上不是那句**）⇒ 随 196 腿的**未提交**工作树来。★我自己那枚同类违规（注释里把 `internal/tools/gate.go` 写成 `tools/gate.go`）**当场改正**后消失。 |
| `sh scripts/check-path-length-budget.sh` | 主检出 | **0** | `denominator: tracked paths=17610 over-budget=61 covered by roster=61 not in roster=0`；`VERDICT GREEN`（本件入库前 17609） |

★**导出树上 `internal/tools` 另有一枚红**：`TestTicket252R1AlignmentAddsNoAuthorization` —— 它要求**仓根父路径有 8.3 短名别名**（`paths_shortname_252_r1_test.go:163-166`：「本探针不许 skip」），而 `D:\tmp\wisp\20261023\363r1\…` 没有；**同一用例在主检出通过** ⇒ 导出树**位置**造成的前置条件缺失，⛔ 产品回归。故 `go test` 的判据读数取主检出（含 `.git`、路径可生短名），导出租**只**用于**甲尺**与 `d22scan`。

## §10 尺与射程（落笔前五问）

1. **量的是哪个对象**：卡数＝`Gate` 两个回调的**调用次数**（⛔ UI／tick／日志行数）；行为＝`Bridge.Execute` 的 `ToolOutcome` ＋ `OnDecision` 的 `Decision`（`DecisionColumn` 由 `route` 在 `emit` **之后**写，故 `decision=` 从**审计行**读，⛔ 从 `OnDecision` 读）；格式甲＝`git archive` 的**入库字节**、乙＝工作树。
2. **输出形状**：`gofumpt -l` 在 Windows 出**反斜杠** ⇒ 一律 `tr '\\' '/'` 后按 `^\.scratch/` 判（⛔ 裸 `^\.scratch/`）。
3. **数代码还是注释**：`d22scan` 违规行是**注释里的引用**（ban #9），`ban #8` 注释豁免、字符串不豁免；两枚新用例的注释里不写会 ring 的路径形态。
4. **和总数对上**：`RUN 399 = 393 + 6`（新用例逐枚数，6 枚＝risk 3 + tools 3）；`PASS = 390 + 6 + k`、`FAIL = 2 − k`，其中 `k∈{0,1}` 就是那枚**时序**用例（`TestResolvePerCallBudget`）逐发漂的一个数（两发读数已在 §9 并列，⛔ 挑好看的那发）；`d22scan ban #8 internal/ 623 = 621 + 2`（本腿两枚新用例各 +1）。
5. **谁算的、什么算法、几位**：行号＝`sed -n '<N>p'` on 工作树；提交名册＝`git show --name-only --format= <sha>`；`.go` 枚数＝同一输出 `tr -d "\r"` ＋ `tr "\\" "/"` ＋ `grep -Ec "\.go$"`；`grep -cE` 计 `=== RUN`／`--- PASS`／`--- FAIL`／`--- SKIP`（顶层，`^… [^/]* ` 排除子测试）；⛔ 无 hash 混比（本件不报 hash）。

## §11 纪律格与残余

① ★**需编排者拍的一件事（我如实登记，⛔ 自行修）**：`d59af9a` 除我 4 枚写面路径外，**还带走了 `.scratch/wisp/issues/164-…-leg.md` 的删除**（95 行）。成因：该删除**在我 `git add` 之前就已在索引里**（我 `git add` 用的是**显式四路径**、⛔ 没清索引），而 `git commit` 落的是**整个索引**；我开工前只查了 `git status --porcelain -- internal/risk internal/tools`（派单要求的那两包），**没查全索引**——这是我的漏。影响面：`164-…-leg-done.md`（防重领 key）**仍在**且本笔未碰。⛔ 我不用 amend/reset（明令禁）；**追加恢复**等于替 164 腿撤销它 staged 的清理意图 ⇒ **等口令**。
② ⛔ 翻框、⛔ 改票名／`Status:`／`-done`、⛔ 改台账、⛔ 改 `docs/PLAN.md`／`docs/specs/**`、⛔ push、⛔ 碰票 179/183/185 原句；临时件只建不删（`D:\tmp\wisp\20261023\363r1\**`，⛔ `.out`）；两次突变全走 `-overlay`（**工作树零突变窗口**）。
③ **残留（不在本票、已具名不修）**：`Gate` 未接线（`NoGate`）或 `L2` 卡面不可达时，"still asks" 仍会打在一条没有卡可开的组合上——那属"审批层没接线"，是另一条（今天生产组合 `cmd/wisp` 接了 `approval.Gate`，见 `run.go:764-788`）。本票只保证：**接线了 gate 的组合里，"still asks" ⇒ 同发真开卡**。

## §12 终态复跑（★读数取在最后一枚 commit 之后）

- **码**：`d59af9a`（产码 4 枚）；**件**：`613f2d87b`（143 行）＋`de79f2660`（§12 与 §9／§10 数字）＋本笔（§6 还原核）；终态行数见盘上 `wc -l`（≤160 帽，过帽即拆 `-partN`）；`internal/**` 字节**自 `d59af9a` 起未再动** ⇒ 上面各门读数对码仍然有效。
- **`git archive 613f2d87b` 导出树**：`sh scripts/d22scan.sh` ⇒ **rc=0**、`clean - no D22 ban violations`（`ban #7 internal/tools/=25`、`ban #8 internal/=623`、`cmd/=148`）、`runtests.sh: OK … PASS=35 FAIL=0 SKIP=0`；`gofumpt -l .` ⇒ **20 行、非 `.scratch` ＝ 0**；`go test -count=1 -run TestTicket363 ./internal/risk/ ./internal/tools/` ⇒ 两包 **ok**、rc=0。
- **主检出**：`sh scripts/check-path-length-budget.sh` ⇒ **rc=0**、`tracked paths=17610 over-budget=61 covered=61 not in roster=0`、`VERDICT GREEN`。
- ★**主检出的 `d22scan` 仍 rc=1**（§9 已归因：196 腿未提交的 `internal/memory/schema.go`）；本腿提交的导出树 **rc=0** ⇒ 本腿的贡献是 **0 违规**，这一点由**同一把尺、两个对象**分别给出。
