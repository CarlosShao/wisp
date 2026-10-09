# 票 285 — 跨卡那一发的**拒因名字在路由级没人断**＋四形今天仍是绿的

**立票**：2026-10-09 10:5x 编排者（来路＝非实现者腿 `242-v2` 的"具名判不动"一节，件 `docs/evidence/s1/242-grant-binding-v2.md:102-110`；台账 `A757`）
**性质**：补仪器。**⛔ 不改产码、不改任何已定案的对外文案**——票 259 选形 ⓐ（`A562`／`A619`：对外一句合并 `ErrBadGrant`、拆分只在对内）是**已批准的契约形状**，本票不动它，只把"这个形状今天有没有人钉"补上。

## 现量（引用前先重跑；枚数口径**逐条写明是抽样还是整族**）

- **① 路由级"两枚都活着的卡"跨卡那一发，`denial=` 那枚名字零断言面（整族结论）**：`internal/agent/approval` 里唯一带两枚活卡走 `Native().Allow` 的断言在 `queue_test.go:161` 一带，而它的 `corr-B` **从未 push** ⇒ 拿到的是 `ErrUnknownCorrelation`，**不是** grant denial；真·两枚活卡那一发只在 `cmd/wisp/subagent_selfapproval_197_test.go:473`，它断的是 `errors.Is(err, ErrBadGrant/ErrPanelAllow)`，**该测试不捕获 gate 的 `Logf`** ⇒ 审计面那句 `GRANT-DENY … denial=` 在两个包里都无人读（尺＝`git grep -c GRANT-DENY` 对产码/测试面零命中，今天只有裁决表与台账写过它）。
- **② `Permitted` 若是"常量形 display-scalar"（`int` 且值不由状态派生），今天只有 1 枚词面尺看得见**（抽样一发＝腿的突变 `MU-D2b` 实测）。
- **③ 反射那把尺读的字段名名单住在本包测试件里**（`ticket242_panelface_test.go` 头部一带），**同一只手把新字段名补进名单**这一发**没被测过**——写腿的硬区（突变只准落非测试产码），所以它是**这一族的已知空洞**，不是"已验不出"。
- **④ `cmd/wisp/approval_reply.go` 两份副本只被 `approval_reply_201_test.go` 一枚**子串**钉（"原生令牌无效"）** ⇒ 删掉"缺失/已用"两枚对内拒因不响（**按尺形登记，本票未种**）。
- 结构性事实（复跑再判）：`queue.go` 里 `spend(nonce, it.bind)` 与 `grantNonce` 写进**同一本** store 的是同一个串 ⇒ 经路由的两端恒等，`misbound` 只在 store 级可达。⚠ 行号一律现读（`242-v2` 已抓到 4 处漂：`gate.go:749→755`、`:109→116`、`:227/:272→:229/:274`、`ci.yml:654→655`）。

## 要建什么

- [ ] **AC#1 把①那一枚缺的尺造出来**：判据＝**两枚同时活着、corr 不同**的卡，经路由（`Native().Allow(B的corr, A的活令牌)`）那一发，**断言面能读出 `denial=` 那枚名字**且指名是" nonce 不在本卡那本里"这一族。落点二选一由腿自己按现量裁并具名理由：甲＝`internal/agent/approval` 包内真两枚活卡；乙＝让 197 那枚载具把 gate 审计行捕获回来断（⚠ 乙会碰 `cmd/wisp`，与在飞腿抢面时**报回不硬做**）。⛔ 不许改对外合并文案，⛔ 不许把"路由能回 `misbound`"写进判据（那属形 ⓑ，口令「259 改形 ⓑ」）。
- [ ] **AC#2 ②③④逐条裁**：每条给"留／补尺／归哪张票"三形之一＋判据（谁在什么条件下会让它变成真洞）。⚠ 本格**只裁形状**，落地另派；③那一条要给"名单与字段名同手扩"这一发**能不能有一枚独立的钉**（例：名单本身由产码派生而不是测试件手写）——只答有没有形状，⛔ 不预设要建。
- [ ] **AC#3 不放宽既有的钉**：交件前复跑 `242-v2` 已具名的四发红句（`MU-A`⇒4 红、`MU-M`⇒9 红含路由级、`MU-D2a/D2c/E`、`MU-OUT`⇒`ticket259_denial_rulers_test.go:252` 一带），逐枚给"仍在／漂了"；⛔ 任何一枚变绿＝本票退回。

## 追加一格（2026-10-09 11:4x 编排者，来路＝只读复核腿 `282-v1` 顶回＋编排者自己复跑）⚠ 不改上面三格任何字

- [ ] **AC#4 第五形：per-call corr 的"每枚调用各不相同"这件事，今天零尺**（这一格是**我自己欠的**，不是实现腿的）。台账 `A745` 我给落地腿的派单里**明写了要一枚新用例**＝"同一任务并发两问 ⇒ 两枚 corr **不同**且各自可路由"；编排者现跑两把尺（HEAD）＝
  ① `git grep -nE 'callCorr' HEAD -- '*.go'` ⇒ 产码 **3 枚**（`internal/agent/loop.go:603` 定义／`:676` 唯一调用点／`:367` 注释）＋测试面 **1 枚**——而且那一枚只是 `ticket283_corr_identity_rulers_test.go:13` 的**注释文字**，不是断言 ⇒ **`callCorr` 零测试引用**；
  ② 全仓扫"两枚 corr 互不相等"形状 ⇒ 只有 `internal/agent/approval/ticket87_veto_l2_test.go:173`（`len(cards) != 2 || cards[0].CorrelationID == cards[1].CorrelationID`），而那枚是**直接往门里 push 两张卡**，**不经 loop 的铸造那一跳**。
  ⇒ 后果形状：`loop.go:603 callCorr(taskID, callID, index)` 若退化成**同一任务内两枚调用返回同一个串**（例如 `callID` 空时那支 `%s#call-%d` 回落跨轮同键），今天**没有任何用例看得见**——票 283 的尺二确实跑了两枚调用，但它在 `:222` 只逐枚断"非空／不等于 taskID／以 taskID 为前缀"，⛔ **没有把两枚互相作差**。
  **要建什么**：一枚会响的断言＝同一任务两枚调用 ⇒ `rows[0].CorrelationID != rows[1].CorrelationID`，且四件套齐（种坏 `callCorr` 的 `index`／`callID` 两支各一发，指名用例必须红）。⚠ 前置同票 283：夹具必须 `TaskID != CorrelationID`，否则两发必同绿不算读数（`282-v1` 具名的方法学坑）。
  ⛔ **本格今天只登记形状，落地按串行铁律排队**（`internal/tools`＋`internal/agent` 的突变面此刻归 `285-r1`/后续腿独占）。

## 禁区

⛔ 不改产码（除 AC#1 新测试件与它必需的最小接线）；⛔ 不动票 259 的选形 ⓐ 与对外那句合并 `ErrBadGrant`；⛔ 不为变绿放宽任何断言、`t.Skip` 不算通过；⛔ 三枚冻结件一字不动；⛔ 零 push；⛔ 不许顺手把票 242/259 的框改了（翻勾归编排者）。

## Progress log

### 2026-10-09 10:4x–10:5x · 腿 `285-r1` 交件读数（⛔ 本节不翻任何勾、⛔ 不改上面任何一句判据）

- **AC#1＝已闭合，落点＝甲**（全部在 `internal/agent/approval` 包内，**零产码字节改动**）。
  新建 `internal/agent/approval/ticket285_route_denial_name_rulers_test.go`（185 行／两枚用例），
  复用 `ticket259_denial_rulers_test.go` 的现成夹具 `r1Audit`(`:36`)／`r1Card`(`:71`)／`r1CardIn`(`:83`)（⛔ 未新造捕获面）。
  ⚠ 本节写于编排者 11:4x 追加 **AC#4** 之后；AC#4（`callCorr` 每枚调用互不相等那一枚尺）**不属本腿派单**，
  本腿零读零断零动，四格的框一律未翻。
  两枚活卡（`corr-285-xcard-a`／`-b`，同一本 queue，`Depth()==2`，两枚都还 pending 且各自可被 `lookupForAllowLocked` 取回）
  ⇒ `Native().Allow(B.corr, A 的未花令牌)` ⇒ 对外仍是合并 `ErrBadGrant`，
  **审计面读回 `denial=spent-or-never-live-nonce`**（尺句落点 `ticket285_route_denial_name_rulers_test.go:125`），
  并钉这一行归 B 不归 A、恰好一枚 `GRANT-DENY`、不落进另外三枚名字、`FORGED-OR-STALE` 合并句仍在、
  正控＝A 的令牌未被 foreign 那一发烧掉且 A 自己仍开卡。
  第二枚是鉴别器：从未 push 的 corr ⇒ `ErrUnknownCorrelation` 且 **`GRANT-DENY` 行数必须为 0**
  ——这一枚就是"为什么 `queue_test.go:161` 那一发不构成 AC#1"的仪器化说法。
- 甲／乙 裁量理由与"零停手项"的现量：`Native().Allow → q.allow → q.allowScoped` 四跳今天全暴露给包内测试，
  ⛔ 因此本腿**没碰 `gate.go`** ⇒ 票 284 的条件搭载不触发（未为它单开改动）；乙形要写 `cmd/wisp`，甲形既够就没有理由抢那一枚面。
- 突变五发全部种在盘上（⛔ 零 `-overlay`）：MU-285q（`queue.go:423` 摘掉 `denial=%s`）⇒ 新尺红 `:125`；
  MU-285p（`approval.go:518` 两枚名折成一枚 `"denied"`）⇒ 新尺＋259 两枚尺共 3 红；
  **MU-285t（本腿自己的用例种坏：期望名换成 `denial=misbound`，产码不动）⇒ 新尺红 `:125`**＝尺非恒真；
  MU-M 带着新文件复跑 ⇒ 82 PASS／**10 FAIL**（原 9 枚全在＋新尺 1 枚，红在新尺 `:109` 返回值那一层）。
- **AC#3＝四发红句逐枚复跑，全部「仍在」，枚数与 `242-v2` 一致，行号未漂**：
  MU-A⇒4 红（`ticket242_binding_test.go:40/:57/:171`＋`ticket259_denial_rulers_test.go:172`）／
  MU-M⇒9 红（含走路由的 `queue_test.go:167`、`ticket259_denial_rulers_test.go:139`、
  `ticket259_panel_capability_rulers_test.go:301`＋同发 `:308`）／MU-D2a⇒2 红／MU-D2c⇒2 红／MU-E⇒2 红（4 条句 `:73/:82/:101/:111`）／
  MU-OUT⇒1 红（`ticket259_denial_rulers_test.go:252`）。⛔ 零枚变绿、⛔ 零枚断言被放宽。
  顺带复认票面 `:8`／`:12` 两处现量：`queue_test.go:161` 的 `corr-B` 从未 push（仍成立）、
  `queue.go:414` 两端恒等（仍成立）。
- **AC#2＝只裁形状，⛔ 一律不动手**：②「补尺，另派」——`int`＋无人填那一形今天只有词面尺看得见，
  变真洞的条件＝填进去的值不随 `revoke()` 变化（`:229` 的不变性轴只看这一枚），形状＝加第二枚不可变性轴；
  ③「有独立钉的形状，仍属补尺」——名单在 `ticket242_panelface_test.go:21` 手写＝同手可扩，本包内无解，
  独立钉＝把期望集搬到另一只手才会碰的工件（`PanelItem` 唯一包外消费者＝`cmd/wisp/approval_reply.go:397`），
  ⚠ 但新建金样属禁区块上的契约变更，要 owner 拍板；④「补尺，另派」——`cmd/wisp/approval_reply_201_test.go:387`
  仍是子串钉（现量逐字 `{"yes " + corr, "原生令牌无效"},`），形状＝换成等值钉；本腿⛔ 未种（种它要写 `cmd/wisp`）。
- 基线与收工：起手 HEAD `37f2a0b`／`go test ./internal/agent/approval/ -count=1 -v`＝**90 PASS／0 FAIL／1 SKIP**
  （SKIP＝`TestDefaultDeadlineWallClockMeasurement`，基线就在）；收工同命令＝**92 PASS／0 FAIL／1 SKIP**。
  四枚产码件收工 hash 与起手逐字等值；`git status --porcelain -- internal cmd` 空。
  `gofmt -l`／`gofumpt -l` 对本腿新建测试件均**空输出**。
- 读数全文：`.scratch/wisp/probes/285/r1/00-anchor.md`／`10-ac1-ruler.md`／`20-ac3-no-relaxation.md`。
