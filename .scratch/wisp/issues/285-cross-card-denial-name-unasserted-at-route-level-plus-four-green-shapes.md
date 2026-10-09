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

- [x] **AC#1 把①那一枚缺的尺造出来**：判据＝**两枚同时活着、corr 不同**的卡，经路由（`Native().Allow(B的corr, A的活令牌)`）那一发，**断言面能读出 `denial=` 那枚名字**且指名是" nonce 不在本卡那本里"这一族。落点二选一由腿自己按现量裁并具名理由：甲＝`internal/agent/approval` 包内真两枚活卡；乙＝让 197 那枚载具把 gate 审计行捕获回来断（⚠ 乙会碰 `cmd/wisp`，与在飞腿抢面时**报回不硬做**）。⛔ 不许改对外合并文案，⛔ 不许把"路由能回 `misbound`"写进判据（那属形 ⓑ，口令「259 改形 ⓑ」）。
      ✅ 编排者 2026-10-09 11:0x 翻勾。落点＝**甲**，`internal/agent/approval/ticket285_route_denial_name_rulers_test.go`（**185 行新建**，commit `61e13e92`；⛔ 零产码字节——`git diff --name-status 37f2a0bb..HEAD -- '*.go'` 整族只有这一枚新测试件）。断的面＝审计行 `approval: GRANT-DENY corr=corr-285-xcard-b tool=shell.run denial=spent-or-never-live-nonce`，断言在 `:125`；**第二枚用例是鉴别器**（`:183`）＝未 push 的 corr ⇒ `ErrUnknownCorrelation` 且 `GRANT-DENY` 行数**必须为 0**——这一形正是把"`queue_test.go:161` 那一发为什么不算 AC#1"仪器化。夹具**复用票 259 现成的** `r1Audit:36`／`r1Card:71`／`r1CardIn:83`，⛔ 未新造捕获面。
      编排者独立跑：`go test ./internal/agent/approval/ -count=1 -run Ticket285 -v` ⇒ 两枚 **PASS**、`ok … 0.035s`。尺有牙凭据（腿跑的四件套，⛔ 零 `-overlay`、全种盘＋`git show HEAD:` 还原，四枚产码件收工 hash 与起手逐字等值）：**MU-285q**（`queue.go:423` 摘掉 `denial=%s`）⇒ 我的尺红 `:125`；**MU-285p**（`approval.go:518` 两枚名折成 `"denied"`）⇒ 89/3（本尺＋259 两枚尺同红）；**MU-285t＝把本腿自己的新用例期望名换成 `denial=misbound`（产码不动）**⇒ 本尺红 `:125` ＝ **自证它不是恒真**；**MU-M 带着新文件复跑**＝82/10（原 9 枚全在＋本尺 `:109` 那枚"returned `<nil>`"）＝新尺咬的是真洞不是自己人。⚠ 基线口径纠正（腿顶回我派单，成立）：本包今天 **90 PASS／0 FAIL／1 SKIP**（那枚 SKIP＝既有 `TestDefaultDeadlineWallClockMeasurement`，非本程造成），我派单写"90/0"漏了 SKIP 一格；收工同命令＝**92／0／1**（90＋本程 2 枚）。⚠ 另顶回我一处：`Native().Allow` 盘上是**三参带 ctx**（`ui.go:146`），我派单写的两参是照抄裁决表的简写——以盘上为准。
- [x] **AC#2 ②③④逐条裁**：每条给"留／补尺／归哪张票"三形之一＋判据（谁在什么条件下会让它变成真洞）。⚠ 本格**只裁形状**，落地另派；③那一条要给"名单与字段名同手扩"这一发**能不能有一枚独立的钉**（例：名单本身由产码派生而不是测试件手写）——只答有没有形状，⛔ 不预设要建。
      ✅ 编排者 2026-10-09 11:0x 翻勾（腿件 `20-ac3-no-relaxation.md`／回报 §⑤；⛔ 三条一律未动手）。**②补尺·另派**＝`Permitted int` 常量形没人填、只被词面尺看见；**变成真洞的条件**具名＝填进去的值**不随 `revoke()` 变化**（现有 `:229` 那把只看"有没有这个字段"这一轴）⇒ 形状＝再加一枚**不可变性轴**的断言。**③有独立钉的形状，仍属补尺**＝名单手写在测试件里⇒同一只手可以扩；独立钉＝把期望集搬到**另一只手才碰得到**的工件（本包外 `PanelItem` 的唯一消费者＝`cmd/wisp/approval_reply.go:397`）。⚠⚠ **腿说这一形"要 owner 拍板"（新建金样会踩 AGENTS §1.1 的 golden 禁区）——编排者裁：不走金样这一形、也不上他的清单**。理由：①"golden 一字节不动"这条保护的是**既有**比对件，新建一枚金样本质是把"改不动的账本"扩一页，为一个测试技巧去开人工批准面不成比例；②这条属"实现约束选哪条路"，不是功能级决定（同形先例＝`Q-72` 那类账目归位他回过"搬呗，看都看不懂"）⇒ **正解＝同包第二枚测试件里硬编期望集、由 CI 的正常编译面钉住**（⛔ 零新工件、零 §1.1 风险）。**④补尺·另派**＝`approval_reply_201_test.go:387` 今天仍是**子串**钉（"原生令牌无效"），删掉"缺失/已用"两枚对内拒因不响；形状＝换成**等值**钉，⚠ 种它要写 `cmd/wisp`⇒ 与票 279 `AC#1`（P39 随行）、票 244 落地腿排同一条队。
- [x] **AC#3 不放宽既有的钉**：交件前复跑 `242-v2` 已具名的四发红句（`MU-A`⇒4 红、`MU-M`⇒9 红含路由级、`MU-D2a/D2c/E`、`MU-OUT`⇒`ticket259_denial_rulers_test.go:252` 一带），逐枚给"仍在／漂了"；⛔ 任何一枚变绿＝本票退回。
      ✅ 编排者 2026-10-09 11:0x 翻勾。读数来源＝腿 `285-r1` 自己重新下刀复跑（**这是对 `242-v2` 那轮的第二发独立跑**，不是照抄它的表）：`MU-A` 4 红逐名（`ticket242_binding_test.go:40/:57/:171`＋`ticket259_denial_rulers_test.go:172`）／`MU-M` 9 红（`queue_test.go:167`、`ticket259_denial_rulers_test.go:139`、`ticket259_panel_capability_rulers_test.go:301`＋同发 `:308`）／`MU-D2a` 2 红（`:39`、`:174`）／`MU-D2c` 2 红（`:39`、`:229`）／`MU-E` 2 红 4 句（`:73/:82/:101/:111`）／`MU-OUT` 1 红（`:252`）⇒ **六发全部"仍在"，枚数与 `242-v2` 一致、行号零漂、⛔ 零枚变绿、零枚断言放宽**。编排者侧独立复认＝① 收工同命令整包 **92／0／1**（＝基线 90＋本程 2 枚，⛔ 没有"少了几枚"这种洗分母）；② `git status --porcelain -- internal cmd` 交件前后均为**空**＝种进产码的坏全部还原。⚠ 诚实边界：**逐枚突变的红句原文我这轮没有第三次重跑**（重跑一枚＝再种一次盘），本格按"两轮独立复跑一致"采信，引它时请带 `242-v2`＋`285-r1` 两枚腿名。

## 追加一格（2026-10-09 10:4x 编排者，来路＝只读复核腿 `282-v1` 顶回＋编排者自己复跑）⚠ 不改上面三格任何字

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
  ⚠ 本节写于编排者 10:4x 追加 **AC#4** 之后；AC#4（`callCorr` 每枚调用互不相等那一枚尺）**不属本腿派单**，
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

### 2026-10-09 11:04–11:20 · 腿 `285-r2` 交件读数（AC#4 第五形补尺；⛔ 本节不翻任何勾、⛔ 未改上面任何一句判据）

- **AC#4＝补尺已落地，落点＝甲＋乙都给**，零产码字节改动（`git diff --name-status ce18b3d..HEAD -- '*.go'`
  只有两枚 `A`）。起手锚 `ce18b3d6`，第一笔＝起手锚 commit `28258cd2`（先锚后跑），交件 commit `9a00a890`。
  - 甲＝`internal/agent/ticket285_corr_distinct_rulers_test.go`（144 行／`Test285CallCorrIsDistinctPerCall`
    ＋`Test285LoopDispatchesOneCorrPerCall`）：包内直接对未导出的 `callCorr` 的**两支各自作差**
    （带 call id 支／回落到轮内位置支），＋ 真 loop 派发侧两枚 corr 作差与"corr 后缀＝自己的 call id"。
  - 乙＝`internal/tools/ticket285_corr_rows_rulers_test.go`（113 行／`Test285RosterRowsCarryDistinctCorrPerCall`）：
    续票 283 的真链（真 loop→真 bridge→真工具→真 `tool_call` 行），断言面落在 `ListToolCallsByTask` 读回的
    **两枚持久行**上＝票面判据那句 `rows[0].CorrelationID != rows[1].CorrelationID`（`乙:81`）。
    复用 283 现成夹具（`build221`／`windowGate221`／`parent283Provider`／`fake197Provider`），⛔ 未改该件一字。
  - 甲不可替代的那一发＝回落支：现量 `grep -rn '#call-' internal/ --include=*_test.go` ＝ 空输出，
    全仓 SSE 夹具没有一枚不带 call id ⇒ 真链走不到那一支，只有包内直接调用能打死它（MU-2 的对照读数）。
- **四件套突变三发全部种在盘上（⛔ 零 `-overlay`），种刀只有 `internal/agent/loop.go` 的 `callCorr`**：
  MU-1（带 id 支塌成常量后缀）⇒ 甲 `:79`/`:85`/`:89` ＋ 甲 loop `:124` ＋ 乙 `:82`/`:90`/`:95`/`:107` 全红，
  **同发下票 283 的尺二 `--- PASS`**＝作差那一发是新增承重件；
  MU-2（回落支 `index` 抹常量）⇒ 只有甲的单元作差红（`:79`/`:85`/`:89`），
  甲的 loop 那一发与乙与 283 全部 `--- PASS`＝"这一支今天除甲之外零尺"的直接读数；
  MU-3（保住互不相等、摘掉 task 前缀）⇒ 前缀轴与"后缀＝自己 call id"轴红（`:76`/`:133`/`:137`／乙 `:76`），
  **作差那一发静默**，同发 283 的尺二也红（那一形今天已有尺，本腿不占功）。
  三发还原后 `internal/agent/loop.go` hash 与起手逐字等值（`8eb37e9f…843`），`git status --porcelain -- internal cmd` 空。
  夹具前置（"TaskID != CorrelationID 否则两发必同绿不算读数"）在甲 `:72`/`:119` 与乙 `:72` 都是**断言**不是注释。
- **⚠ 本格判据的前提要在盘上复核（本腿的顶回，逐字全文见 `10-ac4-ruler.md` §5）**：
  ① `internal/agent/corr_percall_242_test.go:99-105` **今天已经在作差**
  （`c1, c2 := p.corr("call_p1"), p.corr("call_p2")` ⇒ `if c1 == c2 { t.Fatalf(...) }`，
  并在 `:109-133` 用各自的 corr 唤醒各自那一枚调用），票面尺①按 `callCorr` **字样**扫、尺②按
  `CorrelationID ==` **形状**扫，两把都漏了它 ⇒ "第五形今天零尺"这句按现量应窄化为
  **"名册行级作差零尺＋回落支作差零尺"**（这两发＝本腿交付）。
  ② 票面尺①那句"产码 3 枚"现量是 **4 枚**（`loop.go:594` 的文档注释里也写着 `callCorr`，锚上就在）。
  ③ `internal/agent` 整包**在锚上就是红的**：`TestGoldenSingleToolCall`（`loop_golden_test.go:70`）
  期望 `CorrelationID == TaskID`，与票 242 的 per-call 铸形相互矛盾；同文件 `:341` 期望行 corr 等于 task id
  却是绿的。本腿⛔ 未裁定、未放宽、未顺手改（不属 AC#4 射程），`rc` 起手 1＝收工 1。
- 基线与收工（`-count=1`）：`internal/agent` 起手 **82 PASS／1 FAIL／0 SKIP**（`rc=1`）→ 收工
  **84 PASS／1 FAIL／0 SKIP**（`rc=1`，红仍是起手那枚、名目逐字同）；
  `internal/tools` 起手 **208 PASS／0 FAIL／0 SKIP**（`rc=0`）→ 收工 **209 PASS／0 FAIL／0 SKIP**（`rc=0`）。
  `gofmt -l`／`gofumpt -l` 对两枚新件均空输出；`tools/d22scan -root .` ＝ clean（`rc=0`）；
  ⛔ 零 `t.Skip`、⛔ 未动 `internal/tools/loop_approval_test.go` 那一带、⛔ 未碰 `approval/gate.go`（票 284 搭载格不触发）、
  ⛔ 三枚冻结件一字不动、⛔ 零 push。
- 读数全文：`.scratch/wisp/probes/285/r2/00-anchor.md`／`10-ac4-ruler.md`／`20-mutations.md`。
