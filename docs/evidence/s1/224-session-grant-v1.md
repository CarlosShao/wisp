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
| AC#1 `:33` 身份有唯一的铸造点 | 改前生产铸造者＝0／改后 ≥1；具名写出"会话结束"是哪一行代码在做什么 | §1 | **成立（附一条条件＝§6‑a）** |
| AC#2 `:34-35` 三件套齐全 | 写／读／失效三格各自独立用例；⛔ 不许合并成一枚"持久化"用例 | §2 | **不成立**（不合并做到了；"答复⇒落行"零判据，M4 全盘绿） |
| AC#3 `:36` 重启必失效 | 进程重启 ⇒ 同类请求**重新弹卡** | §3 | **不成立**（实现成立；交件无重启、无"弹卡"观测） |
| AC#4 `:38` 反控（正控！） | 种"派生式 session id 让授权跨会话仍生效"的假腿 ⇒ AC#2／AC#3 至少一枚必须红；三枚钉全绿而授权实际生效 ⇒ 本票判失败 | §4 | **附条件成立**（M1 两枚红＝有牙；缺⑩ 的①"跨两次启动"） |
| AC#5 `:39` 不许把"长期"混进这一票 | "长期／永久允许"落 `[fs] allowed_dirs`，不是这张票的射程 | §5 | **成立** |
| 派单增量第 1 条（具名缺口） | 同批改 `cmd/wisp/run_mode101_test.go:427-441` 的期望与两句注释 | §6‑b | **漏做**（断言半不必改且仍真；注释半 `:427-428` 今天已是假话） |
| 派单增量第 2 条（加一枚锁） | 断言生产铸造的 id **永远不等于** `session-before-restart` 这类测试字面量 | §6‑a／§1 | **半做**（便宜方向 `TestTicket224TestLiteralsAreOutsideTheMintedShape` 有；贵方向被推给一枚不存在的用例，由我探针代判） |
| 派单增量第 3 条（两发读数） | 改前该钉绿的读数＋改后同一枚按新期望跑的结果；`git status --porcelain -- internal cmd` 为空 | §6 | **未交**（没改就没有"改后"；`git status --porcelain -- internal cmd` 我现跑＝**空**，两枚 commit 后工作树那两枚包干净） |
| 派单必答 2「拒绝是谁给的」 | 会话授权的"拒"不许由 `approval.NoGate` 的 `PendingWindow` 代答；定向突变必须红；门的文案不许被权限复用 | §2 | **本票无该形状**（8 枚全读"弹卡计数"，无一枚观测"被拒"；M2 五枚红）；⚠ 文案那条靠"根本没有可复用的拒绝文案"满足，**无仪器钉住** |
| 派单必答 3「有没有执行者」 | `internal/session` 那 271 行是否真到达 `cmd/wisp` 装配出来的那条路；会话身份是否宿主铸造、不许调用方自报 | §1／§3 | **实现到达（我现跑探针）；交件不证明它到达**（幻影用例）；身份＝宿主铸造、调用方无法自报（三条独立证据） |

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

## §2 AC#2 三件套齐全（写／读／失效） — **不成立**

票面原文：**写**（答复"本会话内允许" ⇒ `approval_grant` 真有一行，带 session 身份＋工具＋模式＋创建时间）／**读**（下一次同类请求命中它 ⇒ **不再弹卡**）／**失效**（会话结束 ⇒ 同一身份查不到、行为回到"要问"），三格各自独立用例、⛔ 不许并成一枚"持久化"用例。

### 先答它形式上做到的部分（成立的那半）

- **不许合并**这一条做到了：三格是**四枚独立顶层函数**〔我现跑 `grep -c '^func Test'`＝session 9 枚＋tools 8 枚〕，
  `TestTicket224RecordWritesOneRowPerAnswer`（写）／`TestTicket224CoveringMatchesItsOwnSession`＋`TestTicket224LiveGrantStopsTheL1Question`（读）／`TestTicket224CoveringDoesNotCrossToANewSession`（失效），
  文件头逐字引 `ticket90_persist_test.go` 的 "They are three test functions on purpose" 作理由 ⇒ **这一条不是装饰**。
- 17 枚 12:2x 全绿〔我现跑 `-run Ticket224 -v`，逐枚 `--- PASS`；命令与读数在 §1 表下〕。

### 判失败的三处，逐处具名

1. **"写"那一跳零判据**（最重）。票面写的是"**答复** ⇒ 有一行"，而答复到落行之间有四跳：`replySurface.session`（`cmd/wisp/approval_reply.go:236`）→ `Replies.AllowSession`（`internal/agent/approval/replies.go:337`）→ `Gate.allowSession`（`gate.go:651`）→ `Ledger.Record`（`internal/session/grants.go`）。**四跳零用例**：两枚 commit 没往 `internal/agent/approval/` 与 `cmd/wisp/` 落任何测试文件〔读台件：`8b57a419 --name-only`＝2 枚文件〕。
   ⇒ 我不接受"读代码看着通"，所以做了**定向突变 M4**：把 `gate.go:678` 的 `id, err := g.grants.Record(ctx, tool, p)` 换成 `id, err := int64(0), error(nil)`（答复照旧放行、**盘上一行不落**，还照样打 `GRANT-RECORDED`）。〔我现跑，用完 `git cat-file blob HEAD:` 还原，`wc -c` 前后都＝31010〕
   ```
   ok  github.com/CarlosShao/wisp/internal/agent/approval   0.318s
   ok  github.com/CarlosShao/wisp/internal/session          0.332s
   ok  github.com/CarlosShao/wisp/internal/tools           12.187s
   ok  github.com/CarlosShao/wisp/cmd/wisp                156.567s   ← 相关四包一枚不红
   ```
   **⇒ 今天把"本会话内允许"整档的记账摘掉，全套件仍然全绿。**这不是"覆盖薄弱"，是**该格的判据不存在**——而且它恰好是 `grant_test.go:14-17` 自己立的目标：
   「the failure this feature can have is not "it broke" but **"it silently grants nothing and everything still asks"**」。突变 M4 造的正是这个形状，用例一枚不响。
2. **"读"那一格的两半从未接起来**：tools 侧 8 枚全部用 `t224Grants` **桩**（`grant_test.go:51`），不碰真 ledger；session 侧的 `Covering` 用真 SQLite 但不碰 bridge；把它们接起来那一枚，两处注释都指向 `cmd/wisp` 的幻影用例（§6‑a）。
   ⇒ 交件里"下一次同类请求命中它 ⇒ 不再弹卡"这句话**没有任何一次执行**。（我补的探针执行了并判绿：装配根 `tools: GRANT-USE tool=fs.write paths=1 grant_id=1 assessed=L1 mode=ask_every_step` ＋ `windows_before=1 → windows_after_total=1`＝delta 0 ⇒ **实现成立、判据缺**。）
3. **"失效"只测了"查不到"，没测"行为回到要问"**：`TestTicket224CoveringDoesNotCrossToANewSession` 断的是 `rows==0` 与 `Covering==false`（ledger 层）；票面那句"行为回到'要问'"属 bridge／装配层，交件里无人断。

### 逐枚扫 1,067 行的"恒真／装饰"结果（必答 2 的账）

- ⛔ **没发现"其实恒真"的假绿判据**，而且**必答 2 点名的那个形状在本票根本不成立**——理由要写清：tools 侧 8 枚的"问与不问"全部读 `t90Gate` 的 `windows/approvals` **计数**（`grant_test.go:204/221/268/282/319/336/408/431`），
  〔读台件 `internal/tools/ticket90_test.go:113-122`〕`t90Gate.PendingWindow` 默认回 `AnswerTimeout`，而 bridge 的 L1 分支逐字 `case AnswerAllow, AnswerTimeout:`＝**超时即执行**；所以这些用例观测的是"**有没有弹卡**"，**没有任何一枚观测"被拒绝"**。
  ⇒ 本票**没有**出现"L1 穿过 `approval.NoGate`、把门的 `AnswerReject` 当权限判定"的形状；"会话授权会拒"这句话在两枚新测试里**一次都没被断过**。（它出现在**没被改的那枚旧钉**里，见 §6‑b 第 3 条。）
- **定向突变的牙齿实测（四条腿，全部指名用例红；每条用完即还原并 `wc -c` 自证）**：
  | 突变 | 改的那一支 | 红的用例（〔我现跑〕逐枚具名） |
  |---|---|---|
  | **M1** | `session.go:99` 的 `Mint()` 改成 `hex.EncodeToString([]byte("wisp-pid-derived"))`（16 字节⇒形状合法、每发进程算出**同一枚** key＝票面点名的假腿） | `TestTicket224MintIsRandomAndValid`（`grants_test.go:108`）／`TestTicket224CoveringDoesNotCrossToANewSession`（`:355`）＝**2 枚红**；四枚旧钉照绿（§3） |
  | **M2** | `bridge.go:428` `if grantID != 0 {` → `if true \|\| grantID != 0 {`＝**把权限那一支单独改成放行** | `LiveGrantStopsTheL1Question`（:205 对照半）／`PartialPathCoverageStillAsks`（:269/:283）／`GrantSourceThatCannotAnswerMeansAsking`（两枚子用例 :409）／`PathlessCallIsNeverGranted`（:432）／`GrantedCallBooksAllowSessionGrantWithItsRow`（:483 对照半）＝**5 枚红** |
  | **M3** | `bridge.go:334` 摘掉 `&& sil.Level == risk.L1`（L2／Deny 也去问授权） | `SessionGrantNeverCoversL2`（:324/:339）／`SessionGrantNeverCoversDeny`（:373） |
  | **M3b** | M3 ＋ `route` 的 L2 分支插入 `if grantID != 0 { … return true }`＝**把"AC 声称要防的结局"（L2 被会话授权覆盖）真造出来** | 同上两枚**再加行为半**：:320 `declared-L2 call with a live grant produced 0 approval cards, want 1`、:336 同 ⇒ **行为半不是装饰，两半都有牙** |
- ⚠ **文案射程要具名**（必答 2 后半句"门拦下的文案不许复用权限拦下的文案"）：
  〔我现跑 grep〕新增代码**不产出任何拒绝文案**——`GrantSource` 只有 `Covering` 一枚方法（只答"要不要免问"），`route` 的 L1 grant 分支 `return true, ""`；三处新字符串是审计行前缀（`tools: GRANT-USE`／`session: GRANT-HIT`／`approval: GRANT-DROPPED`）与控制台答复行。
  ⇒ **这一条是靠"根本没有可复用的文案"满足的，不是靠用例钉住的**：两枚新测试里**没有一枚断过任何 reason/文案**（〔我现跑〕`grep -n "Contains\|Reason\|reason"` 在 1,067 行里只命中 2 枚结构性 `strings.Contains`：`grant_test.go:234` 断规范化形式、`:511` 断字段名）。谁日后写"会话授权拒了"这句话，本票没有仪器拦他复用门的那句。
- ⚠ **三处轻微装饰／窄射程**（够不上判失败，逐枚记账）：
  ① `GrantSourceThatCannotAnswerMeansAsking` 的 `errOn` 子用例：桩自己把 err 译成 `(0,false)`（`grant_test.go:88`），bridge 只看见"没覆盖"，与"没命中"同形 ⇒ 该子用例断的其实是 bridge 对 false 的反应；真 fail-closed 半在 `TestTicket224LedgerFailsClosedOnStoreErrors`。
  ② `GrantedCallBooksAllowSessionGrantWithItsRow` 用桩 `t224Journal` 断 `row.GrantID` ⇒ **真 `memory.Store` 的 `tool_call.grant_id` 列今天零用例**（我的探针只读到审计行 `grant_id=1`，没读回那一列）。
  ③ `TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority` 的"Bridge 无 setter"半只罚 `Kind()==Func` 的字段（§1 已具名）。

## §3 AC#3 重启必失效 — **不成立**（实现成立，判据不存在）

- 交件里最接近的一枚是 `TestTicket224CoveringDoesNotCrossToANewSession`：它**不重启**，是**同一进程内两枚 `Mint()`**。ledger 层的性质是真的（钥匙不同⇒读 0 行），
  但票面 AC#3 那句是"**进程重启 ⇒ 同类请求重新弹卡**"——**重启**这一跳与**弹卡**这一观测，交件里都没有；两处注释把它推给幻影用例（§6‑a）。
- ⚠ 更要紧：**"同进程两次 mint 不等"推不出"跨进程两次 mint 不等"**。今天 `Mint()` 无进程外状态，所以由构造成立；但**没有用例钉住它**（M1 那一红来自"同进程两次相等"，它救不了"跨进程相等而同进程不等"这一形，例如拿启动时刻的秒级时钟派生）。
- 四枚旧钉在同一发 M1 下**全部照绿**〔我现跑〕：`ok github.com/CarlosShao/wisp/internal/perm 0.204s`（`-run Ticket90`＝三枚钉）＋`ok github.com/CarlosShao/wisp/cmd/wisp 3.298s`（`-run TestTicket101SessionGrantDoesNotCrossRestart`＝第四枚）
  ⇒ 票面 ⑩/⚠ 那句"**四枚全绿而永久通行证实际生效是可能的**"由我**实测复现**，不是复述。
  并且通行证**真的生效**：M1 下我的生产探针报 `two boots minted the SAME id "sess_776973702d7069642d64657269766564"` ＋ `boot 2's session can read 1 grant rows left by boot 1, want 0` ＋ `after a restart the same L1 call showed 0 windows, want >=1` ⇒ **第二发进程带着上一发的授权静默写了盘**。
- **实现层我判它成立**（§1／正常码下的生产读数）：boot2 铸出不同 key、读 0 行、同类 L1 请求重新弹卡
  （`PROBE-B id1=sess_cee79de7… id2=sess_b8ea814f… rows_id1=1 rows_id2=0 boot1_delta=0 boot1_l2=1 boot2_windows=1`）；死会话那一行**没被删**（`rows_id1=1`，SPEC-02 §4 的 30 天审计窗）。
- ⇒ 判语只能给**不成立**：票面⑩已把 AC#3/AC#4 的措辞重写成"①跨两次启动／②生产 id 写＋第二次生产 id 查／③真 `Bridge.Execute`"，**缺一件就不算判据**，而①在盘上一件都没有。

## §4 AC#4 反控（正控！） — **附条件成立**

- 票面第一句「种一条'派生式 session id 让授权跨会话仍生效'的假腿 ⇒ **AC#2／AC#3 至少一枚必须红**」：**成立**，由 M1 实测两枚红（§3），且红的是**性质断言**而不是巧合；再加 §2 表的 M2／M3／M3b，判定链侧、支路侧、行为侧的牙齿都现跑验过。
- 票面第二句（⑩ 重写后的三件判据）：**缺一件**——
  ① 调**生产**铸造函数跨两次启动断两枚 id 不等 ⇒ **盘上无用例**（幻影；只有本腿探针跑过）；
  ② 用生产 id 写、再用第二次生产 id 查断 0 行 ⇒ **半件**（`Mint()` 就是生产那枚函数，但没跨启动）；
  ③ 走真 `Bridge.Execute` 断"问与不问" ⇒ **有**（tools 侧 8 枚都走真 `Execute`），但授权侧接的是桩、不是 ledger。
- 条件（一句话）：**反控本身不假绿，但它今天只在 ledger 层红；"生产那层能不能被这枚假腿打红"由我补的仪器判的，交件的仪器看不见**。把①落成真用例即升为成立——我这发探针可直接作底稿（§6‑a）。

## §5 AC#5 不许把"长期"混进这一票 — **成立**

〔我现跑〕四条边界全部干净：
- 两枚 commit 的 `--name-only` 12 枚文件里**不含** `cmd/wisp/approval_always.go`、不含 `internal/perm/**`（`git show --stat -- cmd/wisp/approval_always.go internal/perm/` 输出为空）⇒"长期"那条档（`always <编号>` → `[fs] allowed_dirs`）一字未动；
- 新代码只写 `memory.GrantScopeSession` 一枚 scope 值〔我现跑 `grep -rn "GrantScope" internal/session/*.go cmd/wisp/*.go \| grep -v _test`＝仅 `grants.go:158`〕；
- `grep -rn "allowed_dirs" internal/session/ internal/tools/grant.go`＝0 命中；
- 面板侧**没有**长出"允许"：`git show 3b78c246 -- internal/agent/approval/ \| grep -E "^\+func .*Panel\|^\+.*PanelAllow"` 只命中一行注释，逐字 "There is deliberately no PanelAllowSession"；新增方法只有 `nativeAPI.AllowSession` 与 `Replies.AllowSession`，后者要求 `card.Grant != ""` 且走 `g.Native()`。
- ⛔ 本腿**没有**跑或改 `internal/panel/l2_grant_boundary_test.go`（冻结件，禁词表只作读台件引用），两枚 commit 也没有新增 C17 方法名。

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
