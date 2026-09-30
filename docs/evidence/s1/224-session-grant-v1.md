# 224 — 会话授权（session-scoped grant）对抗验收表 · 腿 `224-v1`

- 被验收对象：`3b78c246 feat(224)`（10 枚文件 / +909 −25 产码）＋ `8b57a419 test(224)`（2 枚文件 / +1,067 用例）
- 工单：`.scratch/wisp/issues/224-session-scoped-grant-has-zero-executors-and-no-session-identity.md`（54 行，5 枚 `- [ ]` 全未勾）
- 本表作者＝**非实现者验收腿 `224-v1`**；实现腿 `224-r1` 已死（150 轮耗尽），其 commit message 与注释自述一律按**未验证**处理
- 本表生成时刻 `2026-09-30 12:1x +08`
- ⛔ 票面那 5 枚 `- [ ]` **一枚都没碰**（勾框归编排者）；本表只给判语，不翻框
- ⛔ `internal/audio/**` 此刻有另一枚写腿（票 241）在飞 ⇒ 本表零碰、零跑、零归因
- ⛔ `frontend/**` ／ `design/**` 两层禁：零读、零引

## 口径声明（每个数字带推导式）

- 〔我现跑〕＝本腿本机此刻执行并贴回读数；〔读台件〕＝只读代码／台账／commit，未执行。
- 编排者已给的起跑基线（**非我跑的，抄自派单**）：`go build ./internal/session/ ./internal/agent/approval/ ./internal/tools/ ./cmd/wisp/` rc=0；
  `go test ./internal/session/ ./internal/tools/ -count=1` ⇒ ok 0.389s／12.761s；`go test ./cmd/wisp/ -count=1` ⇒ ok 125.521s。
- 用例枚数口径：`grep -c '^func Test' <file>`（文件内顶层 Test 函数数），不是断言数、不是子用例数。
- 行数口径：`git show --numstat` 的增删列；"被扫文件数" ≠ "违规数"。

---

## §0 逐格 1:1 表（票面 5 格 → 本表 5 节）

| 票面格 | 原文要求（要点） | 本表节 | 判语 |
|---|---|---|---|
| AC#1 `:33` 身份有唯一的铸造点 | 改前生产铸造者＝0／改后 ≥1；具名写出"会话结束"是哪一行代码在做什么 | §1 | （填写中） |
| AC#2 `:34-35` 三件套齐全 | 写／读／失效三格各自独立用例；⛔ 不许合并成一枚"持久化"用例 | §2 | （填写中） |
| AC#3 `:36` 重启必失效 | 进程重启 ⇒ 同类请求重新弹卡 | §3 | （填写中） |
| AC#4 `:38` 反控（正控！） | 种"派生式 session id 让授权跨会话仍生效"的假腿 ⇒ AC#2／AC#3 至少一枚必须红；三枚钉全绿而授权实际生效 ⇒ 本票判失败 | §4 | （填写中） |
| AC#5 `:39` 不许把"长期"混进这一票 | "长期／永久允许"落 `[fs] allowed_dirs`，不是这张票的射程 | §5 | （填写中） |
| 派单增量第 1 条（具名缺口） | 同批改 `cmd/wisp/run_mode101_test.go:427-441` 的期望与两句注释 | §1／§6 | （填写中） |
| 派单增量第 2 条（加一枚锁） | 断言生产铸造的 id **永远不等于** `session-before-restart` 这类测试字面量 | §2／§4 | （填写中） |
| 派单增量第 3 条（两发读数） | 改前该钉绿的读数＋改后按新期望跑的结果；`git status --porcelain -- internal cmd` 为空 | §6 | （填写中） |
| 派单必答 2「拒绝是谁给的」 | 会话授权的"拒"不许由 `approval.NoGate` 的 `PendingWindow` 代答；定向突变必须红 | §2／§4 | （填写中） |
| 派单必答 3「有没有执行者」 | `internal/session` 那 271 行是否真到达 `cmd/wisp` 装配路径；会话身份是否宿主铸造、不许调用方自报 | §1／§3 | （填写中） |

## §1 AC#1 身份有唯一的铸造点 — **成立（附一条条件，见 §6‑a）**

票面原文要求两件事：`grep` 点数改前＝0／改后 ≥1；**且具名写出"会话结束"是哪一行代码在做什么**。

### 铸造点数（〔我现跑〕，口径＝"生产非测试调用者"的 `grep` 命中行数）

| 侧 | 尺 | 读数 |
|---|---|---|
| 改前 | `git grep -n "session\.Mint\|session\.NewLedger\|func Mint" 3b78c246^ -- internal cmd` | **0 行命中**（连定义都没有 ⇒ 票面"生产铸造者＝0"是真的，不是漏计） |
| 改后 | `grep -rn "session\.Mint()\|session\.NewLedger(" --include=*.go internal cmd \| grep -v _test` | **2 行，同一处**：`cmd/wisp/run.go:414`（铸）＋`cmd/wisp/run.go:419`（用同一枚 id 建 ledger） |
| 改后 | `grep -rn "GrantSource\b" --include=*.go internal cmd \| grep -v _test` | 声明处 `internal/tools/grant.go:37`；唯一实现者＝`*session.Ledger`；唯一注入处 `cmd/wisp/run.go:514/:653`；**无 setter**（`grep "func (b \*Bridge) Set\|b.grants =" internal/tools/*.go \| grep -v _test` ＝ 0 命中） |

⇒ **唯一铸造点成立**：整棵树只有 `assembleRuntime` 铸；测试里的 `Mint()` 调用不算铸造点（它们是同一枚函数）。

### 会话结束＝哪一行代码在做什么

**没有那一行，这是设计而不是遗漏**，具名如下：
- 生效中的锚点＝`cmd/wisp/run.go:414` 铸出的 `sessID` 只进 `rt.session`（进程内值），没有任何持久化写手；
- 于是"结束"这件事由**进程退出**执行，实现侧表述在 `internal/session/grants.go` 文件头第三格逐字 **"the invalidation side has no code at all"**，与 `A435` 第 2 条（结束点今天＝进程退出）一字不差；
- ⚠ 本腿**不接受**"没代码＝判不了"这一条：AC#3 的判据必须落在"下一发进程算不出同一枚 key"上，那一格在 §3 用生产装配实测。

### 会话身份是宿主铸造、还是能被调用方自报（必答 3 的后半）

**宿主铸造，且调用方无法自报**，三条独立证据：
1. `NewLedger` 强制 `o.ID.Valid()`，`Valid()` 只认 `sess_` ＋ 恰 32 枚小写 hex ⇒ 手写串进不来；
2. `Covering` 的钥匙永远是 `l.id`，**不是**方法形参 ⇒ 调用方在一次 `Execute` 里没有任何位置填"我是哪个会话"；
3. 〔我现跑〕`TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority` PASS：反射遍历 `agent.ToolRequest` 字段名不含 `session`/`grant`；`tools.GrantSource` 恰 1 枚方法（`Covering`），方法名不含 record/revoke/set/insert/add/delete。
   ⚠ 这枚用例的"Bridge 无 setter"那一半**射程偏窄**：它只罚"名字以 grant 开头**且 `Kind()==Func`**"的字段；`Bridge.grants` 是 interface 类型（`Kind()==Interface`）⇒ "往 Bridge 塞一枚 interface 型可换源字段"这一形它不响。今天没有这种字段（上面第 3 行的 `grep` 现跑＝0 命中），但**这一半判据是漏的**，归 §7。

### 生产可达性（必答 3 的前半：那 271 行到不到得了 `cmd/wisp` 装配出来的那条路）

**到得了**，且是我自己现跑的，不是读台件。台件＝`.scratch/wisp/probes/224/v1/probe_a_prod_test.go`＋`overlay.json`（`-overlay` 注入虚拟文件 `cmd/wisp/zz_probe_224v1_test.go`，**仓内跟踪文件零改动**，`git status --porcelain -- internal cmd` 仍为空）：

```
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
go test ./cmd/wisp/ -count=1 -run Probe224v1 -v -overlay=.scratch/wisp/probes/224/v1/overlay.json
=> --- PASS: TestProbe224v1ProductionGrantReachesAssembledBridge (3.99s)
=> --- PASS: TestProbe224v1ProductionSessionDiesWithTheProcess (6.01s)
=> ok github.com/CarlosShao/wisp/cmd/wisp 10.047s
```

装配根真日志里被我逐行抓到的证据链（`[audit]` 行原文，非推断）：
- `wisp run: SESSION-MINT id=sess_f751ec854b7027b1d7f9aeae7844bf05 (结束点＝本进程退出，A435 第 2 条)` ⇒ **生产铸的是合法形状**，且不等于任何测试字面量（`派单增量第 2 条`的"贵方向"在这一发里由我判了）；
- `tools: GRANT-USE tool=fs.write paths=1 grant_id=1 assessed=L1 mode=ask_every_step` ⇒ **注给 bridge 的读面真是那枚 ledger**；
- `tools: call ... decision=allow_session_grant ... grant_id=1` ⇒ `tool_call.decision`／`grant_id` 两枚冻结列第一次有生产写手（D45-2"可取证"落地）。

⚠ 三枚必须写下来的"读法"：
- 我这发**不是** AC#2 的"答复"半：我是直接 `rt.session.Record()` 落行，**绕过了 `Gate.allowSession` 那条原生答复路**；答复路那一半在 §2 单独判，它今天零用例。
- 头两版探针我判错了一次并**当场纠正**：`fs.read` 在 allowlist 内是 **L0**（逐字 `risk=L0 ... reason="无规则命中（L0 直接执行）"`），`fs.write` 覆盖**已存在**的文件是 **L2**（逐字 `rules_hit=[R1 R8]`＋`R8: 不可逆操作（覆盖已有内容）`）⇒ 两者都到不了"未被静默的 L1"那一格，第一版据此误报过"AC#2 读 FAILS"。**结论以第三版为准**（`fs.write` 到不存在的路径＝只有 R1 命中＝L1），两版失败读数留在台件不删。
- 探针里那枚"活授权在场、同一条调用变成 L2 仍弹卡"的正控读数＝`l2_cards=1`（同一 boot、同一枚 grant）。

## §2 AC#2 三件套齐全（写／读／失效） — （填写中）

## §3 AC#3 重启必失效 — （填写中）

## §4 AC#4 反控（正控！） — （填写中）

## §5 AC#5 不许把"长期"混进这一票 — （填写中）

## §6 没做完／留给编排者的台账动作

- **§6‑a 幻影用例（本腿认定的第一具名缺口，最重）**：`internal/session/grants_test.go:20/:94/:134` 与 `internal/tools/grant_test.go:7` 四处注释**指名**一枚 `cmd/wisp` 的
  `TestTicket224ProductionSessionDoesNotSurviveRestart`，把 AC#1 的"贵方向"、⑩ 的控制②、AC#4 的牙齿都推给它。〔我现跑〕**它不存在**：
  `go test ./cmd/wisp/ -count=1 -run Ticket224 -v` ⇒ `testing: warning: no tests to run ... [no tests to run]`；`8b57a419 --name-only` 只有 2 枚文件、`cmd/wisp/**` 零命中；
  全仓 `grep -rn "ProductionSessionDoesNotSurviveRestart"` 的命中**只有那三行注释本身**。
  ⇒ 交件把"生产装配侧那一发"写成注释里的指针而不是盘上的用例。本腿用自己的探针把这格顶下来了（§1／§3 的 PASS 读数），**但交件本身没有这一发**；建议编排者登记为"票 224 交付缺 1 枚生产用例，注释引用不存在的用例名"。
- **§6‑b 派单增量第 1 条（`cmd/wisp/run_mode101_test.go:427-441`）＝漏做，不是"本来不必做"**：现读〔我现跑，带行号〕`:427-428` 逐字
  「`// Refused with the grant live IN ITS OWN SESSION, too: the assembly hands`／`// the bridge a mode, never a grant source (Options.Confirmations nil).`」。
  拆两半，**不许合并**：
  1. **断言半（`:434-441`）今天仍为真，但真的换了原因**。编排者派单里"接了⇒那枚钉红"的前提是"生产授权会命中这一发"，它没命中——**不是因为"没有授权来源"**（改后有了：`cmd/wisp/run.go:653 Grants: grantRead`〔我现跑 grep〕），**是因为 fixture 的字面量 `session-before-restart` 与装配铸出的 `sess_<hex>` 不是同一枚会话**。⇒ 期望值不需改，但注释必须改成带条件的事实句才不假。
  2. **注释半（`:427-428`）现在是假话**，且派单点名的正是这一句（「不许留 "never a grant source" 这种在改后会变成假话的绝对句」）。产码腿在 commit message 第 5 条自己承认"旧注释……已改成带条件的事实句"——**它改的是 `cmd/wisp/run.go:542-560`，同形状的测试侧那两句一枚没动**。
     ⇒ **裁：漏做**。缺＝注释半＋锁的可执行半（派单第 2 条被推给 §6‑a 那枚不存在的用例；`internal/session/grants_test.go:135 TestTicket224TestLiteralsAreOutsideTheMintedShape` 只覆盖了"便宜方向"）。
     ⛔ 本腿**不改**那枚文件：改期望属产码，且按派单第 3 条要同时交"改前绿读数＋改后读数"，那是续腿的活。
  3. ⚠ 这一格里正好埋着**必答 2 的形状**，顺手具名：`:438` 观察到的"被拒绝"其实是 **L2 审批卡超时**——同一台架、同一条 `fs.write` 的逐字证据（本腿现跑）：
     `approval: ANSWER-EXPIRED corr=... tool=fs.write decision=timeout->reject after=1s` ＋ `tools: call ... risk=L2 decision=timeout ... grant_id=0`。
     **那个"拒"是那道门给的，不是权限判定给的。**今天它无害（它防的就是"没有任何东西来免问"），但**谁把这一行改成"带活授权也要被拒"的新期望，就必须改成测"有没有弹卡"而不是"有没有被拒"**，否则就是本仓栽过的那一发假绿。
- **§6‑c 答复路零用例**：见 §2（`Gate.allowSession`／`Replies.AllowSession`／`replySurface.session` 三格今天没有任何用例，突变证据见 §4‑M4）。
- §0 表内 AC#2／AC#3／AC#4／AC#5 四行判语、§2–§5 四节：AC#4 的定向突变待跑完补齐。
- 留给编排者的既有欠账（本腿只登记、不代裁）：
  - `recordablePattern` 只写"逐字相等"、**glob 方言未定案**〔读台件 `internal/session/grants.go` 该函数头上那段自述〕⇒ 是否要另立票，待裁。
  - `clockCeilingDefault = memory.GrantAuditTTL`（30 天）是把"会话时长"这一问**登记为未定案**后取的兜底数，"会话时长该多长"仍待人拍。
  - 失效一侧**零代码**（`grants.go` 自述 "the invalidation side has no code at all"）＝A435 第 2 条的形状，⛔ 不许当缺陷，但它使 AC#3 的判据只能落在"身份不可重算"上 ⇒ 这一条要单独立账。

## §7 我攻不动的地方

- ⚛ **真机重启**：`wisp run` 全流程要 LLM provider／WebView2／桌面句柄，本腿无法在本机做一次真冷启动再冷启动，只能用 `cmd/wisp` 既有的 `h.start(...)` 台架〔读台件，需现跑确认它真的另起一发进程〕替代；若那台架是同进程复用对象，"重启必失效"就只剩构造函数级证据，我攻不动到端到端。
- ⚛ **面板（WebView2）侧**：`PanelAPI` 有没有可能长出会话档入口，要靠 `internal/panel/l2_grant_boundary_test.go`（**冻结件，只读不跑**）的词表判断，我不能改它来探边界。
- ⚛ **`internal/audio` 在飞**：若 224 的读数与那枚包有耦合（我不读它、不跑它），我判不动，会具名写在这条上而不是归因。
- ⚛ **随机源故障分支**：`session.Mint()` 的 `crypto/rand` 失败路径没有注入口，我无法在不改产码的前提下注入，因此"铸不出来时只关掉这一档"这一支只能读、不能打。
