# 票 285 · 腿 `285-r1` · AC#1＝把路由级跨卡那一枚拒因名字造出尺

锚＝`00-anchor.md`（HEAD `37f2a0b`，基线 90 PASS／0 FAIL／1 SKIP）。
本腿交件 commit（⛔ 零 push）＝`61e13e92f038a99f7cb30c19433d34b62b6c2817`（1 枚文件／185 行插入）。

## 判据原文（票面 `:16` 逐字，⛔ 未换写法）

> 判据＝**两枚同时活着、corr 不同**的卡，经路由（`Native().Allow(B的corr, A的活令牌)`）那一发，
> **断言面能读出 `denial=` 那枚名字**且指名是" nonce 不在本卡那本里"这一族。

（派单转述成 `Native().Allow(B的corr, A的grant)` 两参；盘上现量三参带 ctx，见 `00-anchor.md` 末节，以盘上为准。）

## 落点＝**甲**（全部落在 `internal/agent/approval` 包内），理由三条都是现量

1. **本包已有把审计行捕获回来的夹具**，不是新造：
   `ticket259_denial_rulers_test.go:36` `type r1Audit struct`（把 `Options.Logf` 收进切片）·
   `:48 has()` · `:59 dump()` · `:71 r1Card(...)` 返回 `*Gate,*Queue,*qitem,nonce,*r1Audit` ·
   `:83 r1CardIn(t, q, ...)`——它的注释原文就是"stand up one more card on an EXISTING queue, so a case can
   hold two live cards in the one funnel"。⇒ **判据要的那枚"两枚活卡＋一个捕获面"的形状，夹具本来就为它存在**，
   本腿只是没人把那一发接上去。派单点名的 R-143-2 那枚"空桩吞成 `nil`"在本包不存在（`r1Audit` 是真收）。
2. **零产码改动就闭合**：`Native().Allow → q.allow → q.allowScoped`（`queue.go:383`）里
   `:385 lookupForAllowLocked` 收 B 的 corr、`:390` 收 state、`:414 it.grants.spend(nonce, it.bind)` 扫 B 那本 store、
   `:423` 写 `GRANT-DENY … denial=%s`。甲形要的四跳**今天全都暴露给包内测试**（`r1Card` 已经在用 `q.push`＋`q.grantNonce`）。
   ⇒ 不需要乙形，也不需要任何暴露面扩口。**本腿没有停手项**（见末节）。
3. **乙形会碰 `cmd/wisp`**（`subagent_selfapproval_197_test.go` 那枚载具不捕获 `Logf`，要改它就得写 `cmd/wisp/*_test.go`），
   派单第 4/6 条给那一发加了独占面与 harness 前置；甲形既然够，⛔ 没有理由去抢那一枚面。
   ⛔ 本腿零 `cmd/wisp/**` 写；`cmd/wisp` 那两发**未跑**（未跑＝不是失败，本腿不需要它）。

## 新建文件

`internal/agent/approval/ticket285_route_denial_name_rulers_test.go`（185 行，两枚用例＋一枚 helper＋两枚 const）

| 用例 | 断什么 | 关键行 |
|---|---|---|
| `TestTicket285R1CrossCardAllowNamesTheDenialInTheAudit` | 两枚活卡（`corr-285-xcard-a`／`-b`，同一本 queue，`Depth()==2`）、A 的**未花**令牌（前置 `cardA.grants.live()==1`）递给 `Native().Allow(B.corr, nonceA)` ⇒ 对外仍是合并 `ErrBadGrant`（`:109`），且**审计面读出**`:125` 那枚 `approval: GRANT-DENY corr=corr-285-xcard-b tool=shell.run denial=spent-or-never-live-nonce`；行归 B 不归 A（`:131`）；本轮恰好 1 枚 `GRANT-DENY`（`:135`）；不落进另外三枚名字（`:139`）；合并句 `FORGED-OR-STALE` 仍在（`:145`）；**正控**＝A 的令牌没被foreign那一发烧掉、A 自己仍开卡（`:153`／`:156`，ANSWER-ALLOW 行 `:159`） | 尺面＝审计行，不是返回值 |
| `TestTicket285R1UnknownCorrelationBooksNoDenialLine` | 242-v2 具名的那枚"看着像但不是"的形状：`corr-285-never-pushed` 从未 push ⇒ `ErrUnknownCorrelation`（`:180`）且 `GRANT-DENY` 行数**必须为 0**（`:183`）。这一枚是让上一枚不恒真的另一半：**拒了 ≠ 读得出拒因**；`queue_test.go` 那一发今天读的正是这一形，所以它不构成 AC#1 | 鉴别器 |

辅助：`xcardLive()`（`:63`）＝卡还 pending **且** `lookupForAllowLocked(它自己的 corr)` 取回它自己 ⇒
"这条路由没在跑"那一维由**仪器**排除，不再由"缺失的卡"给出（242-v2 `:105` 的抱怨就是这一条）。

⛔ 没有把"路由能回 `misbound`"写进判据：`:139` 那个循环是**反向**钉（这一发不许落进另外三枚名字），
`misbound` 的可达性归口票 259 选形 ⓐ／`ticket259_denial_rulers_test.go:165` 的 store 级断言，本腿一字不动。
⛔ 未动对外合并文案，⛔ 未放宽任何既有断言（既有 90 枚一枚没改，收工复跑＝92 PASS＝90＋本腿新增 2 枚）。

## 突变四件套（⛔ 零 `-overlay`，五种全部种在盘上并逐发还原）

每发＝① 种前 `git show HEAD:<path> | git hash-object --stdin` ② 种后 `sed -n 'Np'` 复量该行 ③ 红句原文含 `file:line`
④ 还原后 hash 逐字等值＋`git status --porcelain -- internal cmd` 空。四枚产码件起手 hash 见 `00-anchor.md`。

### 发 1｜MU-285q＝把审计面的拒因名字从那一行里摘掉（这一发专测本腿新尺不是恒真）

- 种前：`queue.go` = `66fec7ae375344d7bc741270b42c4ec5050cf41f`
- 种后 `sed -n '423p'` 逐字：
  `		q.logf("approval: GRANT-DENY corr=%s tool=%s", corr, tool) // MU-285q planted by leg 285-r1: denial name dropped from the audit face`
- 红句原文（`-run Ticket285R1`，rc=1／1 PASS／**1 FAIL**）：
  `    ticket285_route_denial_name_rulers_test.go:125: AC#1 RED: the assertion face cannot read which denial the two-live-card cross-card route booked (want "approval: GRANT-DENY corr=corr-285-xcard-b tool=shell.run denial=spent-or-never-live-nonce")`
  ⇒ 另一枚（鉴别用例）仍绿＝它本来就不指望这一行，形状对得上。
- 还原：`queue.go` = `66fec7ae375344d7bc741270b42c4ec5050cf41f`（逐字等值；porcelain 在收工复量统一打印）。

### 发 2｜MU-285p＝把两枚拒因折成同一枚名字（`label()` 里 `denialSpentNonce` ⇒ `"denied"`）

- 种前：`approval.go` = `67fb146898dc7b235623bb8ad2fac0e0b2465fe1`
- 种后 `sed -n '518p'` 逐字：`		return "denied" // MU-285p planted by leg 285-r1: two denial names folded into one`
- 整包 rc=1／89 PASS／**3 FAIL**：`TestTicket259R1DenialNamesSpentOrNeverLiveNonce`／
  `TestTicket259R1OutwardAnswerStaysMergedAcrossDenials`／`TestTicket285R1CrossCardAllowNamesTheDenialInTheAudit`。
  本腿那一枚的红句原文：
  `    ticket285_route_denial_name_rulers_test.go:125: AC#1 RED: the assertion face cannot read which denial the two-live-card cross-card route booked (want "approval: GRANT-DENY corr=corr-285-xcard-b tool=shell.run denial=spent-or-never-live-nonce")`
  ⇒ 名字被折时**新尺与既有两枚尺同时见红**＝本尺与 259 那族尺站在同一侧，没有互相吞。
- 还原：`approval.go` = `67fb146898dc7b235623bb8ad2fac0e0b2465fe1`（逐字等值）。

### 发 3｜MU-285t＝**本腿自己的用例被种坏**（派单第 1 条点名要的那一发）

- 种前：本文件（HEAD 版）＝`3d61b8e9356afa2280b1e68f9f7bf4609d5c04b3`
- 种后 `sed -n '51p'` 逐字：`	xcardDenialToken = "denial=misbound"`（产码**未动**）
- rc=1／1 PASS／**1 FAIL**，红句原文：
  `    ticket285_route_denial_name_rulers_test.go:125: AC#1 RED: the assertion face cannot read which denial the two-live-card cross-card route booked (want "approval: GRANT-DENY corr=corr-285-xcard-b tool=shell.run denial=misbound")`
  ⇒ 期望一改、产码不动就红＝那一行**真的在读书**，不是恒真判据。
- 还原：本文件 = `3d61b8e9356afa2280b1e68f9f7bf4609d5c04b3`（与 HEAD 逐字等值）＋porcelain 空。

### 发 4｜MU-M（242-v2 那一发）**带着本腿新文件复跑**＝新尺也咬真洞

- 种前 `approval.go` = `67fb146898dc7b235623bb8ad2fac0e0b2465fe1`；
  种后 `sed -n '574p'` 逐字：`	return denialNone // MU-M planted by leg 285-r1 re-run: membership scan gutted`
- 整包 rc=1／**82 PASS／10 FAIL**（＝242-v2 记的 9 枚红**全部仍在** ＋ 本腿 1 枚）。本腿那一枚红句原文：
  `    ticket285_route_denial_name_rulers_test.go:109: AC#1 RED: the two-live-card cross-card answer returned <nil>, want the merged ErrBadGrant`
  ⇒ 成员扫描被掏空时，跨卡那一发**真的放行**，新尺在第一现场（返回值）就红，不必等审计行。
  既有 9 枚红名册与逐字红句见 `20-ac3-no-relaxation.md`。
- 还原：`approval.go` = `67fb146898dc7b235623bb8ad2fac0e0b2465fe1`。

### 发 5｜`-overlay` 零使用

五发的"种"全部是 `perl -i -pe`／`sed -i` 写盘，"还原"全部是 `git show HEAD:<path> > <path>` 写盘后
`git hash-object` 复量。⛔ 全程没有 `-overlay`、没有临时目录副本、没有改 `go.mod`。

⚠ **porcelain 复量的口径（不假装比实际更多）**：`git status --porcelain -- internal cmd` **空**这一发
本腿打印过三处——AC#3 的 9 还原步（`CLEAN`）、五发跑完的收工复量（`CLEAN`，同批四枚 hash 逐字等值）、
以及起手锚之前一次。每发还原**都单独打印了 hash 等值**，但⛔ 不是每发都单独打印 porcelain；
末发的名册确认＝收工时 `internal`＋`cmd` 两枚面全空 ⇒ 没有残留突变。

## 门禁自查

- `gofmt -l internal/agent/approval/ticket285_route_denial_name_rulers_test.go` ⇒ **空输出**。
- `D:\work\base\gopath\bin\gofumpt.exe -l <同一枚文件>` ⇒ **空输出**（跑到了，非"未跑"）。
- ⛔ 未跑全仓门禁。
- 收工复跑：`go test ./internal/agent/approval/ -count=1 -v` ⇒ rc=0／**92 PASS／0 FAIL／1 SKIP**／`ok … 0.456s`
  （SKIP 仍是指手画脚的那一枚墙钟计量 `TestDefaultDeadlineWallClockMeasurement`，基线就在，本腿没碰）。
- 四枚产码件收工 hash ＝ 起手 hash **逐字等值**；`git status --porcelain -- internal cmd` **空**。
- ⚠ **本文件自报一处自己造的漂**：初稿把新尺在 MU-M 那一发的红句行号写成 `:104`、并把表里六处辅助行号按记忆写成
  `:66`/`:147`/`:157`/`:160`/`:181`/`:184`；收工前用同一次取数复量（`grep -n "t\.Fatalf\|t\.Errorf\|^func "` 对
  本文件＋`grep -o 'rulers_test.go:[0-9]*' ` 对四份原始 `-v` 日志逐枚比）⇒ 真值＝`:63`/`:109`/`:145`/`:153`/`:156`/`:159`/`:180`/`:183`，
  已就地改齐。**这正是票面 `:12` 那条"行号一律现读"的规矩，写腿对自己同样适用**；四发的原始 `-v` 日志在仓外
  会话临时目录 `/tmp/285r1-mu-{q,p,t,m2}.txt`（⛔ 没有落进仓，`.gitignore` 第 8 行那枚全仓 `*.out` 因此无关）。

## 停手项

**无。** AC#1 在甲形下零产码改动即闭合（判据要的"断言面能读出 `denial=` 那枚名字"落点
`queue.go:423` 早就被 `Options.Logf` 暴露给本包测试，`259-r1` 已经建好夹具）。
本腿因此**没有碰 `gate.go`** ⇒ 票 284 的条件搭载**不触发**（⛔ 不为它单开一次改动），
`internal/agent/approval/gate.go` 收工 hash `78e2d28a8b676754514bd80898ace8c5a974d3ce` 与起手逐字等值。
