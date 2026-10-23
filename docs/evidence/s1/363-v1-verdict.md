# 票 363 终裁腿 `363-v1` —— `AC#2`（落地＋牙）与 `AC#3`（门禁）分别能不能翻

- **腿**：`363-v1`（**非实现者／裁决者**，裁决者≠实现者）。被裁对象＝`363-r1`（产码 `d59af9a`；件 `docs/evidence/s1/363-r1-impl.md` 152 行）。
- **起手**：HEAD `7e99ccf81`；`git log d59af9a..HEAD -- internal/risk internal/tools` ⇒ **0 笔** ⇒ 我量的是 `d59af9a` 的字节（四枚工作树 blob 与提交 **4/4 同串**）。
- **环境**：`GOFLAGS=`、`unset GOTMPDIR TMPDIR TEMP TMP`（每发都印过，值均为 unset）、⛔ 未用 `CGO_ENABLED=0`；`go1.27.1 windows/amd64`。
- **仪器**：①我自己写的探针（仓外 `D:\tmp\wisp\20261023\363v1\probe\zz_t363v1_probe_test.go`，经 `-overlay` 挂成虚拟路径，**⛔ 抄它的 fixture**）；②两枚**自造**突变体（`…\363v1\mut\mode.M1.go`／`bridge.M2.go`）；③`git archive d59af9a` 导出树（甲尺＋本腿贡献）。
  ★**仓内 `.go` 字节 0 枚、工作树零突变**（`git status --porcelain -- internal/risk internal/tools` 首／尾各一次均空）。
- **裁语**：**`AC#2` 能翻**（带三条具名残余）／**`AC#3` 能翻**（带一条具名既存红）。见 §9。

## §0 票面前提新鲜度（对 `d59af9a^` 的 blob 现跑，⛔ 抄票面）

| 票面那条链 | 票面行 | 我量到（`git show d59af9a^:<p> \| sed -n`） | 判 |
|---|---|---|---|
| `sil := mode.Screen(verdict)` | `bridge.go:311` | 逐字同号 | 不漂 |
| fail-closed 支 | `mode.go:163-165` | `:163` `if d.Level != L1 && d.Level != L2 {`／`:164` `return Silenced{Level: d.Level, Kept: "无法识别的风险级，fail-closed 保留询问", Mode: m}`／`:165` `}` | 不漂 |
| 日志条件＋call | `bridge.go:323`／`324-325` | `:323` `case sil.Kept != "":`／`:324` `b.log("tools: MODE-REDLINE … still asks …`／`:325` 续行 | 不漂 |
| `// --- routing…` 注释／`b.route` | 票面并成 352 | 注释 **351**、`b.route` **352** | 票面偏差成立（`363-a1` 已记） |
| `route` 的 `case risk.L0` | `432-435` | `:432` `switch sil.Level {`／`:433` `case risk.L0:`／`:434` `allow`／`:435` `return true, ""` | 不漂（★`case` 在 433，票面把 switch 行也算进去了） |
| `assessor.go:41/45/46` | 三句 | 逐字（`L0 < L1 < L2 < Deny…`／`L0: read-only, no side effects - execute directly.`／`L0 Level = iota`） | 不漂 |
| `Kept` 契约 | `mode.go:143-146` | 逐字 | 不漂 |
| `Screen` 两个 `Level: L0` 出口 | `174-175`／`185-188` | `:174` `if d.Level == L1 {`、`:175` `return Silenced{Level: L0, Silenced: true, Mode: m}`；`:185-187` `redLine` 支、`:188` 末尾同一句 | 不漂 |

★**修后行号是这次改动自己推的**：`mode.go` L0 支 **184–186**、fail-closed **187–189**；`bridge.go` 新 `case` **323–335**、`still asks` **336–338**；`route` 的 `case risk.L0` **446–448**（+13 ＝ 新 case 的 13 行）。
★**一处非阻断漂移（我记、⛔ 裁）**：两枚新用例注释里的**活指针**还是改前行号——`ticket363_screen_l0_test.go:55` 写 `bridge.go:433-435`（现 **446-448**）、`ticket363_test.go:269` 写 `437-439`（现 **450-452**）；`ticket363_test.go:7-10` 那四行是**修前链的历史描述**（按改前读是对的），⛔ 算漂。

## §1 ★★`AC#2` 判点一：两向都真的不说了吗（我自造读数）

仪器＝真 `fs.read`（真 D34 声明 `FSReadDecl()`、真 R1–R9 判定器、真 C26 canonicalizer）＋真桥（`tools.New`），唯一替换＝`Logf` 指捕获器；**卡数＝`Gate` 两回调（`gate.go:110-117`）被调用次数**，⛔ 数 UI/tick/日志行。

**〔L0 向〕三档同形**（只有 `mode=` 变；★**无任何 `MODE-` 行**）：
```
DECISION mode=<m> silenced=false kept=""
CARDS L1_WINDOWS=0 L2_APPROVALS=0 TOTAL=0
EXEC isError=false text="hello"
LOG|tools: call kind=success task=v1-task corr=v1-corr tool=fs.read risk=L0 decision=allow outcome=success rules_hit=[] in_allowlist_scope=true grant_id=0 refusal=none reason="无规则命中（L0 直接执行）"
```
**〔Deny 向〕两档同形**（`ask_every_step`／`auto_approve`）：
```
DECISION mode=<m> silenced=false kept="拒绝级 verdict（Deny 级不受任何模式影响）"
CARDS L1_WINDOWS=0 L2_APPROVALS=0 TOTAL=0
LOG|tools: MODE-REDLINE mode=<m> refuses tool=fs.read rules=[R2 R3] because=拒绝级 verdict（Deny 级不受任何模式影响） (a Deny never reaches a gate: no card is raised)
LOG|tools: call kind=denied … risk=L2 decision=reject outcome=error … refusal=redteam-denied
EXEC isError=true text="该目标被绝对禁止访问（R3 A 档，任何授权都不可豁免）：…"
```
⇒ **`L0` 不再说、`Deny` 改说 `refuses`**；两句都与我同发量到的卡数一致（都是 0 卡：一个执行完了、一个被拒）。

**★正控（证明它不是"把日志全关了"）**：

| 格 | 输入（真工具） | 判定 | 卡数 | 结果 | 那行 |
|---|---|---|---|---|---|
| `L1`／`ask_every_step` | `fs.write` 新文件 | `L1 [R1]` | **窗口=1**／审批=0 | 已写入、落盘 | **`MODE-` 0 行** |
| `L1`／`ask_high_risk` | 同上 | `L1 [R1]` | 0/0 | 已写入、落盘 | `MODE-SILENCE … (a silenced question, not an allow: no user click happened)` |
| `L2`／`ask_high_risk` | `fs.write` 覆盖 | `L2 [R1 R8]` | 0／**审批=1** | 已写入 | `still asks … because=ask_high_risk 只静默 L1，本判定是 L2`（**真的**：同发 1 卡） |
| `L2`／`auto_approve` | 同上 | `L2 [R1 R8]` | 0／**审批=1** | 已写入 | `still asks … because=R1 工具自身声明的 L2（不可逆自认，PLAN.md:1629）`（**真的**） |
| `L2`／`ask_every_step`（我加的一格） | 同上 | `L2 [R1 R8]` | 0／审批=1 | 已写入 | **无** `MODE-REDLINE`（没人扣住它，也就没什么可解释） |

⇒ ★`still asks` 只在**真开卡**的地方出现、且那里它是真的；`MODE-SILENCE` 仍说 "not an allow"。

## §2 ★★`AC#2` 判点二：四档行为无回归（我的独立核）

| 档 | 输入 | 判定 | 卡数 | 行为 |
|---|---|---|---|---|
| `L0` 仍放行 | `fs.read`（授权根内） | `L0 []` | **0/0** | `decision=allow outcome=success`、`text="hello"` |
| `L1` 仍窗口 | `fs.write` 新文件／`ask_every_step` | `L1 [R1]` | **窗口=1** | `已写入`、**盘上字节＝"x"**（我真读了盘） |
| `L2` 仍审批 | `fs.write` 覆盖／`ask_high_risk` | `L2 [R1 R8]` | **审批=1** | `已写入` |
| `Deny` 仍拒 | `fs.read` 打 `%APPDATA%\wisp\config.toml` | `Deny [R2 R3]` | **0/0** | `decision=reject`、`isError=true`、`该目标被绝对禁止访问…` |

★**`route` 的 `case risk.L0` 放行语义有没有被动**——两把尺、两个方向都答"没动"：
① **字节尺**：`d59af9a` 对 `bridge.go` **`+13/−0`**（全文**零删除**）、对 `mode.go` `+25/−1`（那唯一的 `−1` 是**注释**：`Deny and any unknown level` → `Deny, L0 and an unknown level`）⇒ 两文件**纯插入**、`route`（445–498）一字未进 diff。
② **行为尺**：上表四档各一发、「放行／窗口／审批／拒绝」一个没变。

## §3 `Kept` 契约有没有被动（它称"一字未动"）

- **契约文字**：`mode.go:143-146` 逐字未变（`d59af9a` 的 diff 未触及这四行）；`Deny` 的 `Kept` 文本逐字未变（我的矩阵三档同串 `拒绝级 verdict（Deny 级不受任何模式影响）`，与修前 `363-r1` §1 贴的同一句）。
- **`L0` 的 `Kept` 由非空变空**：★这是把**违反契约的旧值**改回契约（"`Kept` 留给被红线**留住**的询问"，L0 没有询问可留），⛔ 改契约语义。`mode.go:155-156` 的自述也被同步成 `Deny, L0 and an unknown level are handled first`（注释，与实现对得上）。
- **生产读取面**：`ModeKept`／`ModeSilenced` 在 `cmd/**`、`internal/panel/**` **零读取**（我 grep 全仓复认：读取方只有测试；`internal/agent/approval` 只是原样搬运＋测试）⇒ `Kept` 空掉对行为面无影响。
- ⇒ **判：契约未动、本格不因"改契约"退回**。★残余见 §9①。

## §4 反向判据（★自造，⛔ 采信它的 M1/M2）

两枚突变体建在**仓外**、`diff -u` 逐行可核（各自只少掉那一跳），全程 `-overlay`：

| 读数 | 摘掉哪一跳 | rc | **逐名红行** |
|---|---|---|---|
| **M1**（我造的） | `mode.go` 184–186 `if d.Level == L0 {…}` | **1** | `TestTicket363ScreenRecognizesL0` **4 条**（`:39` × 3 档 + `Mode(7)`）、`TestTicket363FourLevelsKeepTheirRoute` **3 条**（`:123` × `L0/*`）、`TestTicket363L0ReadNeverClaimsAQuestion` **15 条**（`:230`／`:243`／`:247`／`:251`／`:255` × 3 档） |
| **M2**（我造的） | `bridge.go` 323–335 的 `case … sil.Level == risk.Deny:` | **1** | **只红** `TestTicket363DenyIsRefusedNotAsked` **6 条**（`:320`／`:324`／`:327` × 2 档） |
| **互不代偿（双向各证一次）** | — | — | M1 下 `TestTicket363DenyIsRefusedNotAsked` **仍 PASS**；M2 下 `TestTicket363L0ReadNeverClaimsAQuestion` **仍 PASS**；★**我自造的探针在同一对突变体下也各自红**（M1 ⇒ `TestV1ScreenMatrix`＋`TestV1L0ThreeModesSayNothing` 红；M2 ⇒ `TestV1DenyTwoModesSayRefuses` 红）⇒ **牙不是我抄它的** |
| **还原绿** | 不挂 overlay（＝盘上 `d59af9a` 字节） | **0** | 6 枚 `TestTicket363*` 全 PASS；我的 6 枚 `TestV1*` 全 PASS |
| **"今天响、不修不响"** | 同批用例、码回滚（§6 基线） | — | 修前红 **4** 枚、修后全绿；★两枚正控（`UnknownLevel…`／`OtherLevelsKeepTellingTheTruth`）修前**就绿**、修后仍绿 ⇒ 牙既不恒红也不恒绿 |

★**hash 先问"谁算的、什么算法、几位"**：`git hash-object` ＝ **SHA-1／40 位十六进制**（⛔ 与 `sha256sum` 的 64 位混比）。工作树 ⇄ `d59af9a:` ⇒ **4/4 同串**：`4de3bdacb232…`／`69822c751eae…`／`93ed19fa2bb1…`／`250c56db2c90…`；`git archive` 导出树同样 **4/4**。

## §5 `AC#3` 门禁（五门分开报，⛔ 相减）

| 门 | 对象 | rc | 读数 |
|---|---|---|---|
| `GOFLAGS= go build ./...` | 主检出 | **0** | — |
| `go vet ./internal/risk/ ./internal/tools/` | 主检出 | **0** | — |
| `go test -count=1 ./internal/risk/ ./internal/tools/` | 主检出 | **1** | 两发同形：`RUN=399 PASS=397 FAIL=1 SKIP=1`；红名**只有** `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（既存红，见 §6） |
| `go test -count=1 ./internal/perm/ ./internal/agent/approval/` | `Screen` 的两个消费包（我加的一发） | **0** | `ok` / `ok` |
| `sh scripts/d22scan.sh` | **主检出** | **0** | `clean - no D22 ban violations`；`examined 290 production Go files under internal/ and cmd/`；`ban #8 internal/=626`、`cmd/=148`、`ban #7 internal/tools/=25` |
| `sh scripts/d22scan.sh` | **`git archive d59af9a` 导出树** | **0** | `clean`；`ban #8 internal/=623`、`ban #7 internal/tools/=25`、`cmd/=148` ⇒ 与它报的**逐字相同**（§8） |
| `sh scripts/check-path-length-budget.sh` | 主检出 | **0** | `tracked paths=17612 over-budget=61 covered by roster=61 not in roster=0`；`VERDICT GREEN` |

★★**`d22scan` 主检出今天是 `rc=0`（与派单预期不同），归因用"能区分两棵树"的尺**：196 腿引的那枚 `docs/evidence/s1/196-r1-impl.md` **现在在盘上**（`ls -l` 10:20、12,144 字节、`git ls-files` ⇒ **TRACKED**）；引用本身在**工作树** `internal/memory/schema.go:26`／`:157`，而 `git show HEAD:internal/memory/schema.go | grep 196-r1-impl` ⇒ **无命中**。ban #9 判的是 "phantom-citation（注释引的仓内路径**不存在**）"（`tools/d22scan/main.go:872-873`）⇒ **文件一在盘上，违规就消失**。⇒ 今天主检出 clean 是**树的状态**（196 腿的未提交改动＋已入库的件），★**⛔ 记成票 363 的贡献**；票 363 的贡献由导出树那一发单独给（rc=0）。它 09:5x 报的 rc=1 在**它那时**是真的（件还没落盘）。

**★格式甲乙两尺都报**：
- **甲（判据，blob 树＝`git archive d59af9a` 导出树）**：`gofumpt -l .` ⇒ **20 行、全部 `.scratch/**`、非 `.scratch` ＝ 0**（`tr '\\' '/'` 归一后判，⛔ 裸判反斜杠）；四枚逐个 ＋ `internal cmd tools scripts frontend` ⇒ **0 行**、rc=0。
- **乙（诊断，工作树）**：四枚 ⇒ 0 行；`gofmt -l` 四枚 ⇒ 0；四枚 `tr -cd "\r" | wc -c` ⇒ **0 0 0 0**（`*.go text eol=lf`）。

## §6 新增红 0（我的独立核）

★基线尺（**同一批测试文件**、只把码换回去）：`-overlay` 把 `mode.go`／`bridge.go` 换回 `d59af9a^` 的字节 ⇒ `RUN=399 PASS=393 FAIL=5 SKIP=1`，红名＝**`TestC26…` ＋ 4 枚新用例**（`ScreenRecognizesL0`／`FourLevelsKeepTheirRoute`／`L0ReadNeverClaimsAQuestion`／`DenyIsRefusedNotAsked`）。
★终态两发 ⇒ `FAIL=1`，红名只剩 **`TestC26…`** ⇒ **新增红＝0**（终态红名集合是基线的**真子集**）。
★它报的 `PASS=396/FAIL=2` 那发我**没复现**：`TestResolvePerCallBudget`（时序 `1.000 ms/op` 门）在**我这两发都绿** ⇒ 它**逐发漂**。⇒ 判据＝**红名集合逐名比**，⛔ PASS 数；`399 ＝ 393 ＋ 6`（新用例 6 枚：risk 3 ＋ tools 3），这个加数我对上了。

## §7 写面纪律（`d59af9a` 名册逐枚 ＋ 那枚 `D` 的处置）

`git show --name-status d59af9a` ⇒ **6 条**：`D .scratch/wisp/issues/164-…-leg.md`／`M internal/risk/mode.go`(+25/−1)／`A internal/risk/ticket363_screen_l0_test.go`(+129/−0)／`M internal/tools/bridge.go`(+13/−0)／`A internal/tools/ticket363_test.go`(+432/−0) ⇒ **4 枚 `.go`**（2 M ＋ 2 A），与它自报一致；两枚测试**零 `t.Skip`／零 `testing.Short()`**。

**那枚 `D`：我判"这个处置对不对"（⛔ 算成 `363-r1` 的越界）**：
- 事实核：旧名在 `d59af9a^` 有 **95 行**（blob `981fc83d…`）；`-done` 是**另一笔** `51bb042b9`（编排者结票 164「第30票」）新增的；现在盘上只有 `-done`。`-done` 相对旧文＝`Status:` 改写 ＋ 追加 `AC#4` 判语（`diff -u` 共 **58 行**）⇒ 旧文可从 `d59af9a^` 完整取回，**没丢东西**。
- **判：处置对**。①终态本来就要"旧名删、`-done` 在"（`-done` 是防重领的唯一键）；②把旧名**加回来**＝给一张已结的票造**第二枚**键、等于替 164 腿撤销它 staged 的清理意图；③`amend`/`reset` 明令禁。⇒ ★**编排者"①留在史里"是对的**。
- ★但**该记的一笔照记**（⛔ 越界、属**规程漏**）：它 `git add` 前只查了 `-- internal/risk internal/tools`，**没查全索引**；这一条它 §11① 已自报。我这一笔提交前加查 `git diff --cached --name-status`（§11）。

## §8 它报的"导出树 `d22scan` rc=0"这一发（★取法对不对）

`git archive d59af9a` → `D:\tmp\wisp\20261023\363v1\export-d59af9a`，跑 `sh scripts/d22scan.sh` ⇒ **rc=0**、`clean - no D22 ban violations`；`ban #8 internal/＝623`、`ban #7 internal/tools/＝25`、`cmd/＝148` ⇒ **与它报的逐字相同**。
★**取法正确**：判"本腿贡献"要用**提交的字节**（导出树），⛔ 工作树（工作树含别腿未提交的活，今天的 `internal/memory/**` 就是）。四枚 `.go` 在导出树里的 blob 与 `d59af9a:` **4/4 同串**。
★**射程（导出树的代价，⛔ 当产品读数）**：导出树无 `.git` ⇒ d22scan 自印 "gitignore rules NOT APPLIED … every path in every scope is being scanned"；`git archive` 只含 tracked 字节 ⇒ `design/` 在导出树 14 个文本文件、主检出 39（主检出多的是未被跟踪的产物）⇒ 导出树这一发**只用于"本腿贡献"与甲尺**，⛔ 用于整包测试。

## §9 裁语

**★`AC#2` 能翻。** 两向都不说假话（§1 我自造读数）；四档无回归（§2）；`Kept` 契约未动（§3）；两跳各有一枚"今天响、不修不响"的用例且**互不代偿**（§4 我自造突变体）。**带三条具名残余（一条不减）**：
1. `mode.go:141` 那句 `Silenced is true when the caller must NOT ask` 与 L0 的 `Silenced=false` 在**字面**上相冲——★**修前 L0 也是 `false`**（旧支只写 `Kept`）⇒ 这是**既存的文档/实现口径差、⛔ 本笔引入**；要改那句属"**改契约文字**"⇒ 须机主批。落地腿选了 `false` 并在注释里写明取舍（`gate.go:43-47` 的审计口径 ＋ `MODE-SILENCE` 的文案），我判**这个取舍对**。
2. 两枚新用例注释里的活指针行号是改前的（§0 末）。
3. `L2` 被扣住但 **gate 未接线**（`NoGate`）的组合里，`still asks` 仍会打在 0 卡上——属"审批层没接线"，生产组合 `cmd/wisp` 接了 `approval.Gate`（`run.go:764-788`）。本票只保证：**接线了 gate 的组合里，`still asks` ⇒ 同发真开卡**。

**★`AC#3` 能翻（带一条具名既存红）。** 五门里四门 rc=0（build／vet／d22scan 两棵树／路径长度）；`go test` **rc=1**，唯一红名 `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（断言对象是 `internal/panel/**`）在**我自造的基线**（同测试文件、码回滚）上**同样红** ⇒ **既存红、⛔ 本笔新增**；四数按发分开报（§5），⛔ 相减。

## §10 尺与射程（落笔前五问，逐条答）

1. **量的是哪个对象**：卡数＝`Gate` 两回调调用次数（⛔ UI／tick／日志行数）；行为＝`ToolOutcome` ＋ 审计行（`DecisionColumn` 在 `emit` **之后**才写，⛔ 从 `OnDecision` 读）；甲尺＝`git archive` 的**入库字节**、乙尺＝工作树；`d22scan` 的"本腿贡献"＝导出树。
2. **输出形状**：`gofumpt -l` 在 Windows 出**反斜杠** ⇒ 一律 `tr '\\' '/'` 后按 `^\.scratch/` 判；`git show --name-status` 一行一条。
3. **数代码还是注释**：`mode.go` 那 `−1` 是**注释**（我逐行看了 diff，⛔ 当代码删除）；ban #9 违规是**注释里的引用**；新用例 6 枚＝`grep '^func Test'` 逐枚数。
4. **和总数对上**：`399 ＝ 393 ＋ 6`（新用例 6 枚：risk 3 ＋ tools 3）；终态 `397 PASS ＋ 1 FAIL ＋ 1 SKIP ＝ 399`、基线 `393 ＋ 5 ＋ 1 ＝ 399`，差额 **＋4 PASS／−4 FAIL** 正是那 4 枚从红转绿的用例（另 2 枚新正控两态都绿）⇒ 加数闭合；`ban #8 internal/ 623（导出树） ⇄ 626（主检出）`，`626 − 623 = 3` ＝ 两棵树之间新增的 3 枚 `.go`（`internal/memory/taskstate.go`／`taskstate_196_test.go`〔196 腿〕＋`internal/winsec/deleteseam_348r2_windows_test.go`〔348 腿〕，逐名 `comm` 对过）⇒ **⛔ 票 363 的**。
5. **谁算的、什么算法、几位**：`git hash-object` ＝ **SHA-1／40 位**、`git rev-parse <rev>:<p>` 同一算法 ⇒ 只与它比；行号＝`sed -n '<N>p'`／`grep -n` on 具名 blob 或工作树；四数＝`grep -cE` 计 `^=== RUN   [^/ ]+$`／`^--- (PASS|FAIL|SKIP): [^/ ]+`（顶层，排除子测试）；⛔ 无 hash 混比。

## §11 纪律格

⛔ 改 `internal/**`／`cmd/**`／`scripts/**`／`.github/**`／`.scratch/wisp/**`／`docs/reports/**`／`docs/PLAN.md`／`docs/specs/**`；⛔ 落任何 `.go` 字节（本腿 0 枚）；⛔ 翻框／改票名／`Status:`／`-done`／台账；⛔ push。
★两枚突变全程 `-overlay`（替换文件在仓外 `…\363v1\mut\`）⇒ **无工作树突变窗口、无窗口期提交**；工作树在突变前后各查一次 `-- internal/risk internal/tools` 均空。临时件只建不删（`D:\tmp\wisp\20261023\363v1\**`、⛔ `.out`）。
★本笔 `add`／`commit` 两处**逐个文件名字面点名**（`docs/evidence/s1/363-v1-verdict.md`），⛔ `add -A`／`.`；提交前 `git diff --cached --name-status` 核索引（只应出现这一个 `A`）。
★本笔提交号见件尾追加那条（⛔ amend；按 164／362 的先例用**追加一笔**把 sha 写回件里）。
