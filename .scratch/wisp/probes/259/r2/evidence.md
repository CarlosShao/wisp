# 票 259 腿 `259-r2` — AC#1 ⓐ（具名降级）的注释级落地：把绑定那一层的口径改到与盘上一致

本件是 `259-r2` 的唯一证据件。写面只有两枚产码文件（`internal/agent/approval/approval.go`、
`internal/agent/approval/queue.go`）的**注释**＋本件目录。§2 到 §8 随本腿的推进逐笔填，
每一节的读数都带时刻；本腿**不种任何突变**（ⓐ 那一刀是注释级、行为一字不动，没有判据需要突变）。

---

## 0. 锚与进场（与取数同发，12:20:27+0800）

| 尺 | 读数 |
|---|---|
| `date` | `Mon Oct  5 12:20:27 CST 2026`（＝`2026-10-05 12:20:27 +0800`） |
| `git log -1 --format=%h` | `5b638498`（＝台账 `A619` 那一笔，`dev`） |
| `git status --porcelain -- cmd internal tools` | **空输出**（rc=0）⇒ 起手时别人在这三块面上零脏 |
| 本腿写面 | `internal/agent/approval/approval.go`、`internal/agent/approval/queue.go` 的注释＋`.scratch/wisp/probes/259/r2/**` |
| 别腿在飞 | `257-v1`（`cmd/wisp`，非实现者终裁；那枚包会把本腿的包当源码编进去）／`269-a1`（只读普查） |
| 本腿跑过的尺 | 只 `go test -count=1 -v ./internal/agent/approval/`（单包）＋§4 的四道门禁；⛔ 没碰 `cmd/wisp`／`internal/panel`／`internal/tools`／`internal/config`，⛔ 未跑 `go build ./...` |

**具名裁定（我不裁形，只兑现）**＝台账 `docs/reports/pending-and-issues.md` **`A619` §2**
（文件行号 12128-12132，`A619` 标题在 12118）：**票 259 AC#1 选 ⓐ＝具名降级**，撤销口令
「**259 AC#1 改 ⓑ**」。ⓐ §2 边界三条我逐条兑现：① ⓐ 不等于"这层没用了"（摘要仍是"同一枚牌
只能花一次"的一部分）⇒ 我的注释不许写成死代码／可删；② 残余风险（store 成员检查或 state 检查
被重构掉时绑定这层不兜住）必须落在注释里；③ ⓑ 是待设计不是否决 ⇒ 我不扩答复面。

**★行号漂移（具名报差，本腿现读为准）**：票 259 立票锚是 `4eb228ef`（10-03），`A619` §2 抄的是
同一批号；本腿现读的三处全部往后漂了，漂移的成因是 `b6b1d6a4`（`259-r1` 第 2 步：AC#2 把
`spend` 的返回从 `bool` 换成未导出的 `grantDenial`，在它上方插了 `grantDenial` 类型与四枚常量
的注释块）。

| 出处写的 | 现读 | 漂 |
|---|---|---|
| `approval.go:287`（`issue` 存摘要） | `approval.go:415` | +128 |
| `approval.go:298-305`（`spend` 那圈扫描＋比对） | `approval.go:485-502`，比对行在 `:496` | +187／+193 |
| `approval.go:303`（`equalSecret(stored, bind)`） | `approval.go:496` | +193 |
| `queue.go:376`（`spend(nonce, it.bind)`） | `queue.go:382`，且赋值左侧已是 `denial :=` 不是 `spent :=` | +6（形状变了） |
| `queue.go:162`（`it.bind` 唯一写点） | `queue.go:162` | 0 |
| `queue.go:368-371`（state 关） | `queue.go:368-377` | 起点 0 |

---

## 1. 现读到的三处原文（逐字，⛔ 不是照抄票面）

**① 花侧的实参＝那枚 item 自己身上存的同一份摘要** — `internal/agent/approval/queue.go:382`（在
`allowScoped` 里，`q.mu` 仍持锁）：

```go
	denial := it.grants.spend(nonce, it.bind)
```

`it.bind` 的全仓唯一写点＝`internal/agent/approval/queue.go:162`（`push` 铸造时）：

```go
	it.bind = bindDigest(corr, d.TaskID, d.Tool, d.LevelString(), q.seq, d.Args)
```

`spend` 的产码调用者＝**全仓只有一枚**，就是上面那行（尺＝`grep -rn "\.spend(" --include=*.go .`，
其余八枚命中全在 `internal/agent/approval` 的两枚测试文件里：`ticket242_binding_test.go:30/36/47/57/111`
与 `ticket259_denial_rulers_test.go:151/154/170/179`，逐枚现读，不是引用票面）。

**② 存的那一份也来自同一个值** — `internal/agent/approval/approval.go:411-416`（`issue`），
赋值行 `:415` 逐字：

```go
func (s *grantStore) issue(nonce, bind string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[nonce] = bind
}
```

产码里 `issue` 的唯一调用者＝`internal/agent/approval/queue.go:336`（`grantNonce`，`q.mu` 持锁）：

```go
	it.grants.issue(nonce, it.bind)
```

⇒ ①与②两端喂进同一枚 store 的是**同一个 `it.bind`**。store 是 per-item 的（`newGrantStore()` 的
产码构造点两枚：`queue.go:158` 的 `push`、`queue.go:577` 的 `replay` 模板）。

**③ 比对对任何经路由的调用恒等** — `internal/agent/approval/approval.go:485-502`（`spend`），
比对行 `:496` 逐字：

```go
func (s *grantStore) spend(nonce, bind string) grantDenial {
	if nonce == "" {
		return denialMissingNonce
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for v, stored := range s.values {
		if !equalSecret(v, nonce) {
			continue
		}
		delete(s.values, v) // consumed whether or not the binding matched
		if !equalSecret(stored, bind) {
			return denialMisbound
		}
		return denialNone
	}
	return denialSpentNonce
}
```

⇒ 经 `allowScoped` 进来的调用里 `stored == bind`（同一个 `it.bind`），`:496` 那一句**结构上不可能
为 false**；`denialMisbound` 从路由不可达。真挡下"A 的令牌花在 B 上"的两处（现读，含票 §9 对我
自己票面的更正）：

- **防线① store 成员检查**：`approval.go:491-501` 那圈 `for v, stored := range s.values`——A 的
  nonce 不在 B 自己的 store 里 ⇒ 落到 `:501 return denialSpentNonce`。跨卡那一下根本走不到 `:496`。
- **防线② 队列的 state 拒绝**：`queue.go:368-377`（`it.state != statePending` ⇒ `ErrNotPending`）
  ＋`queue.go:298-303`（`deliver` 拒绝重复结算）＋`queue.go:395-398`（花完牌但 deliver 输了 race
  仍记 `denial=card-not-pending`）。⚠ 按票 259 §9 我复认的那条更正：**state 关拦的不是跨卡**
  （跨卡时 B 本来就 pending），它拦的是"牌落在一张已经离开 pending 的卡上"；把它和成员检查并列
  写进注释时，两处各自拦什么我按这个口径分开了写，没有混成一件事。

**④ 与本刀无关但同一层的两处口径**（现读，判断射程用，非引用）：`ui.go:127-129` 那句
`ErrBadGrant` 的注释"missing, spent, or bound to another item"＋`:130` 的对外合并文案——
`A619` §2 与本腿派单都把"对外文案维持合并"写死为裁界（撤销口令「259 对外也区分」），
所以那三行本腿**零字节未动**，归口见 §5。

---

## 2. 逐枚注释 before-after 全文（八枚点，ⓐ 落地的本体）

行号一律给**现读**值（改后态）。每枚先给锚点 `5b638498` 的原文（＝`before/` 那两份快照），再给改后全文。
本腿只写注释；被撤的那几句用引用块标出，好让下一位看见"撤掉的是哪句原话"。

### A. `approval.go:356-388` — grant 权威块（那句把恒等式列为权威来源的话在这里）

before（`approval.go:356-361`）：

```go
// Authority for an allow answer is therefore: (1) the answer arrived on the
// native API surface, (2) it presented a live grant for THIS pending item, and
// (3) the grant's binding digest matches the item it is spent on. The
// caller-settable Request.Source string is logged and never consulted
// (rulings M-7 / C-3 in docs/reports/pending-and-issues.md: a security
// decision keyed on a caller-controlled selector fails open).
```

after（`approval.go:356-388`，三段新增＋`Request.Source` 那两句逐字留下）：

```go
// Authority for an allow answer is therefore: (1) the answer arrived on the
// native API surface, (2) the nonce it presented is a live row in the store OF
// THE CARD it is being spent on - grantStore is one per pending item, and
// Queue.allowScoped reaches it only through the item the queue issued that key
// for - and (3) the item has not already left pending. (2) and (3) are the
// load-bearing pair, and (2) is the one a cross-card attempt lands on: A's
// nonce is not a row in B's store, so spend answers denialSpentNonce.
//
// The binding digest is NOT a fourth gate, and it stops nothing across cards
// on its own. Every routed call spends with the digest this same item's store
// was issued (Queue.allowScoped passes it.bind; grantNonce issues it.bind into
// that store; push writes it.bind once, from bindDigest), so the two sides of
// the comparison inside spend are one string by construction - that branch
// cannot come out unequal on any path the gate can reach. Ticket 259 AC#1 chose
// form (a) for this: a NAMED DOWNGRADE of what the code claims, not a fix of it
// (ledger A619 section 2 of docs/reports/pending-and-issues.md, which also
// carries the revocation phrase for that choice). The three sites to read for
// yourself are Queue.allowScoped's spend call, grantNonce's issue and
// grantStore.spend; the wording on each of them is the same reading.
//
// Downgraded is not the same as unused, and this is the boundary of that claim:
// the row is what spend deletes, so single-use lives in the row, and the digest
// is that row's content - one of the two places this layer records which card a
// proof was minted for, the other being which store holds it. It is also the
// half a form-(b) check would need: recomputing the binding from the request
// the answer actually carried has to compare against something. Form (b) is on
// hold as UNDESIGNED, not rejected: the answer surface carries no tool, args,
// level or sequence to recompute from, so it begins by widening C17's inbound
// face, which is human-approval territory.
//
// The caller-settable Request.Source string is logged and never consulted
// (rulings M-7 / C-3 in docs/reports/pending-and-issues.md: a security
// decision keyed on a caller-controlled selector fails open).
```

★那句"零独立拦截能力"在本枚的落点＝`approval.go:364`（"it stops nothing across cards on its own"）。

### B. `approval.go:405-415` — `bindDigest` 的自我声明（原文那句"挡跨卡"就在这里）

before（`approval.go:378-380`）：

```go
// bindDigest is the domain separator that ties a grant to ONE pending item.
// Two items with the same correlation id at different times (a replay) still
// differ, because the sequence number is folded in.
```

> 撤掉的原话：**"ties a grant to ONE pending item"** —— 这就是票面 AC#1 那句验收要撤的形状。

after：

```go
// bindDigest folds one card's identity - correlation id, task, tool, level,
// sequence number, argument bytes - into the value that card's grant rows
// carry. Two items never share a digest, not even a replay of the same identity
// at a later time, because the sequence number is folded in.
//
// "Which card a nonce belongs to" is NOT decided by comparing anything to this
// digest; the wording on this function used to claim that it was, and ticket
// 259 AC#1 form (a) downgraded it. The digest is per-item data: the check that
// really decides the card is which grantStore holds the row (see grantStore
// below), and every routed spend call hands the comparison this same field
// (see Queue.allowScoped and spend; the block above grantBytes is the whole of it).
```

⚠ 我一度在这里逐字引用被撤的那句原话（`("ties a grant to ONE pending item")`），已撤——
逐字留在码里会让 AC#1 那句"不许留"的 grep 尺自己命中，而"引用一句假话再否认"比直接指名
那句假话更弱。现在的形状是：只说"这里的措辞曾那样自称、票 259 AC#1 形(a) 把它降级了"，
原话全文留在本件 §2 B 的 before 块里（证据件不是码，射程不同）。

### C. `approval.go:436-446` — `grantStore`：把真防线①指名在这里

before（`approval.go:401-403`）：

```go
// grantStore holds the live nonces of one pending item. Every entry is
// single-use: spending deletes it, and so does answering or expiring, so a
// screenshot of a dismissed card cannot authorize anything afterwards.
```

after：

```go
// grantStore holds the live nonces of one pending item - ONE ITEM, and that
// per-item split is the check that really does refuse "A's grant spent on card
// B": a nonce minted for another card is not a row in this store, so spend's
// scan falls through to denialSpentNonce and the binding comparison is never
// even reached. Its companion guard is the queue's state refusals (the
// statePending test in Queue.allowScoped and deliver's refusal to settle an
// item twice), which refuse an answer that arrives after its own card left
// pending - not a cross-card attempt, because the other card is sitting there
// pending. qitem.bind carries both guards and what each covers. Every entry is
// single-use: spending deletes it, and so does answering or expiring, so a
// screenshot of a dismissed card cannot authorize anything afterwards.
```

### D. `approval.go:494-501` — `denialMisbound` 常量注：可达性只在这里被写出来

before（`approval.go:450-453`）：

```go
	// denialMisbound: the proof IS live on this card, but the binding digest
	// it was issued under does not match the one it is presented against.
	denialMisbound
```

after：

```go
	// denialMisbound: the proof IS live on this card, but the binding digest
	// it was issued under does not match the one it is presented against.
	// No routed call can produce this value - Queue.allowScoped presents the
	// item's own stored digest - so the only callers that reach it are inside
	// this package and hand spend two different digests by hand: that is what
	// TestTicket242SpendRejectsForgedBindingAndConsumesTheNonce and
	// TestTicket259R1DenialNamesMisbound do at the store. Those cases pin the
	// store part; they cannot and do not say anything about the route.
	denialMisbound
```

点名前两枚用例＝逐枚现验真在（`ticket242_binding_test.go:26`／`ticket259_denial_rulers_test.go:165`，
尺＝`grep -n "^func Test" <两枚文件>`），且都在 §3 的基线名册里 PASS。⛔ 没写任何不存在的用例名。

### E. `approval.go:528-557` — `spend` 上方：恒等链路＋边界②那条残余风险的唯一落点

before（`approval.go:479-484`，此四句逐字留在改后第一段之后）：

```go
// spend validates and consumes, and returns WHICH denial it reached (ticket
// 259 AC#2): denialNone means the grant opened the card, anything else means it
// did not and says why in audit terms. The control flow is the pre-259 one,
// byte-for-byte equivalent in outcome: the nonce is consumed whether or not the
// binding matched (a rejected caller cannot retry), and an unknown value
// deletes nothing. What changed is only that the branch taken is now reportable.
```

after 新增的两段（接在原文之后，`:533-557`）：

```go
//
// What the bind argument is on a routed call: the only production caller of this
// function is Queue.allowScoped, which passes the item's own it.bind - the same
// string grantNonce wrote into this store through issue - so the comparison
// below is a value against itself, and the membership scan is what decides the
// answer. The binding line has ZERO independent power over a grant aimed at the
// wrong card; the guard inside this function is that scan, and it works only
// because the store is one per item.
//
// Residual risk, stated where a reader looks for the guard (ticket 259 AC#1
// form (a) boundary 2, ledger A619 section 2): if the membership scan below or
// the queue's state refusals are ever refactored away or reordered - a spend
// that answers denialNone on a nonce it merely recognised, or an allow route
// that reaches this function without the statePending test - the line below has
// nothing to notice, because a comparison of a value to itself reports success
// for exactly the calls both guards would have refused. The cases that reach
// denialMisbound all call this function directly with two hand-written digests
// (TestTicket242SpendRejectsForgedBindingAndConsumesTheNonce and
// TestTicket259R1DenialNamesMisbound), so none of them can tell a routed
// identity from a routed disagreement. A619 §2 registers that refactor as
// invisible to today's rulers; this leg seeded no mutation, so it cites that
// reading rather than re-measuring it. Making this line bear real weight is
// form (b) - recomputing the digest from the request the answer carried - which
// is blocked on widening the answer surface, not on a line of code.
```

★边界②那条残余风险的落点＝`approval.go:542-554`（那句"this leg seeded no mutation, so it cites that
reading rather than re-measuring it"就是把"今天没有任何尺看得见"标成**引用未复测**，⛔ 不是本腿量到的）。

### F. `queue.go:44-62` — `qitem.bind` 字段注：被撤的那句"cannot be replayed onto a different request"

before（`queue.go:44-48`）：

```go
	// bind is the digest a grant must match to authorize THIS item. It covers
	// the correlation id, the task, the tool, the level, the argument bytes
	// and the sequence number, so a nonce cannot be replayed onto a different
	// request or onto a later re-issue of the same one.
```

> 撤掉的原话：**"a grant must match to authorize THIS item"** ＋ **"so a nonce cannot be replayed onto a
> different request"** —— 这两句就是"绑定这层能挡跨卡"在码里的形状，是票面 AC#1 验收句要清空的。

after：

```go
	// bind is the digest this item's grant rows are issued against. It covers
	// the correlation id, the task, the tool, the level, the argument bytes and
	// the sequence number, so no two items ever hold the same value - not even
	// the same identity re-queued later (a replay), because seq moves.
	//
	// It is NOT what refuses "A's nonce spent on card B", and the wording here
	// used to say it was. allowScoped spends with this very field, which is the
	// string grantNonce already put into this item's own store, so the binding
	// comparison inside grantStore.spend is a value against itself on every
	// routed call and has zero independent power across cards. The two checks
	// that do refuse, and where a reader finds them:
	//   guard 1 - the store is per item: push gives every item its own
	//   newGrantStore, so B's store holds no row for A's nonce and spend answers
	//   denialSpentNonce; guard 2 - the queue's state refusals: the statePending
	//   test in allowScoped before the spend, plus deliver's refusal to settle an
	//   item twice. Guard 1 is where a cross-card attempt lands; guard 2 is where
	//   an answer lands after its own card has left pending. Downgrade recorded as
	//   ticket 259 AC#1 form (a), ledger A619 section 2; the residual risk of it
	//   is written on grantStore.spend.
```

★那句"零独立拦截能力"的第二枚落点＝`queue.go:53`（"has zero independent power across cards"）。

### G. `queue.go:338-347` — `grantNonce`：把"两端同源"写在存的那一行旁边

before（`queue.go:323-325`，前三句逐字留下）：

```go
// grantNonce mints the single-use native proof for one item. An error means
// the RNG failed: no prompt may be displayed, because a displayed prompt
// nobody can authorize is worse than an honest refusal.
```

after 新增段（`:341-347`）：

```go
//
// The row issue writes here carries this item's own bind - and that same field
// is what allowScoped passes to spend below. Both ends of the spend comparison
// are therefore set from one item's one value, which is why the binding
// comparison cannot disagree on a routed call. Which card a nonce is live on is
// decided by WHICH store the row went into (this one), not by the digest
// matching. See qitem.bind and grantStore.spend.
```

### H. `queue.go:404-413` — 花的那一行（`spend(nonce, it.bind)`）正上方，before 无此段

before：`tool := it.Dec.Tool` 之后直接是 `denial := it.grants.spend(nonce, it.bind)`（无注释）。

after（新增九行，紧贴那枚调用点）：

```go
	// The second argument is this item's own digest, the string grantNonce
	// issued into this item's own store, so the binding comparison inside spend
	// is a value against itself here: this layer stops nothing across cards on
	// its own. Two checks do, and both are in reach from here - the membership
	// scan in grantStore.spend (a foreign nonce is not a row in it.grants, which
	// is B's store when the answer names B), and the statePending test above plus
	// deliver's own settled check below. Recorded as ticket 259 AC#1 form (a),
	// ledger A619 section 2. The residual risk of that downgrade - if either
	// check is ever refactored away, this comparison has nothing left to notice,
	// because it compares a value to itself - is on grantStore.spend.
```

★第三枚落点＝`queue.go:406`。

---

## 3. 零行为自证（三样：numstat／注释对照表／名册作差）

**3.1 逐枚 `git diff --numstat`（锚＝起手 `5b638498`，交付态＝工作树，末笔随本件同交）**

| 文件 | added | deleted | 删除列是否只落在注释行 |
|---|---|---|---|
| `internal/agent/approval/approval.go` | 80 | 7 | 是（尺见 3.2） |
| `internal/agent/approval/approval.go`（已提交笔 `6a021830`） | 79 | 7 | 是 |
| `internal/agent/approval/queue.go` | 36 | 4 | 是 |

**3.2 逐枚"改动的行是否全在注释"对照表**＝`comment-table.txt`（本目录，逐行列出每枚改动行的
新文件行号＋该行原文前 24 字＋判定）。机械尺与读数：

```
approval.go: deleted=7 deleted-non-comment=0 added=80 added-non-comment=0
queue.go:    deleted=4 deleted-non-comment=0 added=36 added-non-comment=0
全表行数：116 COMMENT / 0 NOT-COMMENT / 0 BLANK
```

尺的定义（写死在这里，便于别人复跑）：`git diff -U0 5b638498 -- <file>` 里每一枚 `+`/`-` 行，
剥掉行首空白后必须以 `//` 开头；`BLANK` 单独计数（本次为 0，因 Go 注释空行写的是 `//`）。
另有两条更强的尺：

- **注释剥离后逐字节相同**＝把两枚文件"行首（剥空白）为 `//` 的行＋空行"全部删掉后与锚点版比对：
  `approval.go` 两侧 md5 均 `f8881122791630aa41783fb35a4d8c82`（259 行／6,476 字节），
  `queue.go` 两侧均 `7f2e928994dcdabeb5efe78dd39a5dd7`（382 行／9,282 字节）⇒ **函数签名、返回值、
  分支、字符串文案一视未动**（这三样全在剥掉的代码里）。
- **无块注释混入**＝diff 里 `/*`／`*/` 命中 0 枚（尺＝`git diff 5b638498 -- <两枚> | grep -cE '^[+-].*(/\*|\*/)'`＝0）。

**3.3 行为不变尺＝改动前后单包名册作差（`base-v.txt` vs `after-v.txt`）**

| 读数 | 前（起手态 12:22:36→12:22:38） | 后（交付态 12:45:17→12:45:20；中间 12:38:40 与 12:35:46 两轮同形） |
|---|---|---|
| `go test -count=1 -v ./internal/agent/approval/` rc | 0 | 0 |
| 顶层名（`^--- PASS`） | 70 | 70 |
| `--- PASS` 总行数（含子测试） | 89 | 89 |
| `--- FAIL` | 0 | 0 |
| `--- SKIP` | 1 | 1 |
| `=== RUN` | 90 | 90 |
| 包耗时 | ok 0.456s | ok 0.437s（12:38 那次；12:45 那次同为 ok） |

名册作差＝`diff base-names.txt after-names.txt` **零差异**（`base-names.txt`/`after-names.txt` 已入库）。
那枚 SKIP 是既有的一枚：`TestDefaultDeadlineWallClockMeasurement`（`ticket84_no_owner_test.go:224`
明写"有意慢：300s 墙钟计量，只在 `WISP_84_MEASURE=1` 时跑"）⇒ **登记为 SKIP，不当通过**；
它不在本腿的两枚文件里，改前后同一形状。

⚠ 本腿⛔没跑过 `go build ./...`／`cmd/wisp`／`internal/panel` 任何一把尺（`257-v1` 正在 `cmd/wisp`
取终裁读数，跨包洗读数；那两枚包都会把 approval 当源码编进去）。注释级改动对它们的可见面＝零，
但"零"这一句是**推断不是读数**，跨包名册作差归编排者（同 `259-r1` 报过的那条）。

---

## 4. 门禁读数（带时刻；起手态 12:31 一轮＋交付态 12:45/12:46 一轮，中间 12:37 那轮同形）

| 门禁 | 起手态（改到一半） | **交付态**（＝本件与两枚产码文件将提交的内容） |
|---|---|---|
| `sh scripts/d22scan.sh` | 12:31:00→12:31:35，rc=0，"clean - no D22 ban violations" | 12:45:52→12:46:17，**rc=0**，ban #8 `internal/` examined **512** Go files（comments and `_test.go` included）／`cmd/` 103／`design/` 39／`frontend/` 85 |
| `scripts/check-path-length-budget.sh --with-self-test` | 12:31:44，rc=0，positive control PASSED（三发控制全过） | 12:46:17→12:46:21，**rc=0**，positive control PASSED，VERDICT GREEN，tracked paths=**5905**（12:37 那轮＝5867，涨的是别腿这几分钟的入库）／over-budget=57／roster 覆盖=57／**not in roster=0** |
| `go vet ./internal/agent/approval/` | 12:31:35，rc=0 | 12:45:20→12:45:21，**rc=0**（空输出） |
| `gofumpt -l <两枚改过的文件>` | 12:31:51 与 12:32:04，无输出＝clean | 12:45:21，**rc=0／零行输出**（尺＝`$(go env GOPATH)/bin/gofumpt.exe`，`v0.12.0 (go1.27.1)`） |

**gofumpt 的空读数不是"命令没跑到"**＝交付态并排跑（12:47:37，件 `gofumpt-negctl.txt`）：同一枚二进制
同一把 flags，先跑 `probe.go`＋`clean.go` 那一对 ⇒ 输出**逐字点名**
`C:\Users\swq\AppData\Local\Temp\tmp.kp87kYi0AM\gofumpt-negctl\probe.go`、clean 那枚不点（rc=0）；
紧接着跑我改过的两枚 ⇒ **零行输出**（rc=0）。起手态 12:32:04 那轮同形。负控那枚故意写坏的文件
（函数体内前后空行）建在 **TEMP、仓外**——⛔ 我没在本仓 `.scratch` 里再造第二枚故意写坏的 Go 样本，
因为 `ci.yml:171` 那步 `gofumpt -l .` 会走遍 `.scratch`（＝票 269 正在处理的那枚常红，
⛔ 不该由本腿加伤；台账 `A618` §3 记着这台机器上那三枚 09-28 预存的坏样本）。

**d22scan 那行在册 SKIP 我没当通过**：`d22scan: skipped as git-ignored: 1 file(s) under 1 ignored
director(ies) [frontend/dist/assets/]`（由 `frontend/.gitignore` 决定）。这台机器上有 09-27 的旧构建
产物而 CI 里没有，所以 ban #6/#8 的 `frontend/` 那一档在本机的分母与链上不同（既有归口＝票 171／
票 269 地界），本腿不改它、也不把它读成"全仓干净"。本腿两枚文件落的是 `internal/` 那一档（512 枚
含注释与 `_test.go` 的射程里），那两枚文件里的注释**不豁免字符串**，我的新增行前缀全是 `//`、
不含任何 `U+2190-U+2BFF`／`U+1F000-U+1FAFF`／`U+FE0F` 字形（尺与逐枚读数见 §5 §6 条）。

---

## 5. 判不动／量不到／归口（具名，逐条）

**5.1 `ui.go:127-130` 的对外合并句：本腿零字节未动，但它仍列着那第三种因。**
现读逐字：`:127-129` 注释 "ErrBadGrant: the native proof is missing, spent, or **bound to another
item**. The message never says which, so the API cannot be used to probe for valid nonces." ＋
`:130` 值行 `approval: 原生令牌无效（缺失/已用/与本次请求不绑定）`。按盘上事实，"bound to another item"
这一因**经路由不可能发生**（§1 ③），所以这一句与 §2 A 那三段是同一件事的另一半。我**没动它**，两个原因：
① 派单写死"ui.go 那枚对外合并句一字不动"＋`A619` §2 的裁界（内部与审计区分、对外维持合并，撤销口令
「259 对外也区分」）；② 那是 `ErrBadGrant` 的注释与文案，属票 259 AC#2 已裁的形状，且
`cmd/wisp/approval_reply.go:227/:272` 还有两份同文案副本、`approval_reply_201_test.go` 逐字钉着。
⇒ **归口＝编排者＋`259-v1`**：要么按"注释也算对外合并句"豁免，要么单开一格在不动文案的前提下改注释。
⛔ 我没有自行动它一个字，也没有在注释里绕一句"这其实不会发生"去与 §2 的形状失配（那需要同时动 ui.go）。

**5.2 测试面残留的"能挡跨卡"口径：本腿未动，因为动它必动名册。**现读逐字命中：
- `ticket242_binding_test.go:94` 用例名 `TestTicket242ForgedBindingCannotSpendAnotherItemsGrant`
  ＋ `:110` 注释 "Item B's binding must not open item A's grant even if a nonce leaked."
  ＋ `:112` 红句 "AC#1 RED: item B's binding spent item A's grant - **cross-item binding is broken**"；
- 同文件 `:9-18` 头注 "the binding layer between a grant nonce and the ONE pending item it was minted for"；
- `ticket259_denial_rulers_test.go:28` 一带 "the binding check bite"（同族口径）。

这一族句子里那枚 `:111` 的调用其实是**成员检查**在拦（`"leaked-or-guessed-nonce"` 从来不是任何 store 的行）。
票面 AC#1 验收句写的是"⛔ 不许留任何一句'能挡跨卡'在**码里**"——测试文件算不算"码里"，**我判不动，
停下来上报**：改它要么动函数名（§3.3 的名册零差异判据就破），要么动红句字符串（那是测试面的行为面）。
⇒ **归口＝编排者裁**（要动就得连名册判据一起重判，或另开一票由非实现者做）。
本腿那两枚产码文件里的同类口径已清空，尺见 §8。

**5.3 边界②"今天没有任何尺看得见那次重构"＝引用，未复测。** 本腿⛔不许种突变（`257-v1` 正在
`cmd/wisp` 取终裁读数，而 `cmd/wisp`／`internal/panel` 编译都带上 approval），所以那一格我只把它
作为台账 `A619` §2 的登记写进注释，并在注释里明写"this leg seeded no mutation"。
⇒ **装牙归 `259-v1`**（票 §2 边界②已写明这一格归它）。

**5.4 ⓑ 的可行性我现验了形状、没有动它。** 现读：答复面 `gate.go:715-721`
`Request{CorrelationID, Allow, Grant, Reason, Source}`（票 §9 写 `gate.go:709-715`，漂 +6，成因是
`220-r1`/`260-r4` 那两笔在它上方加了注释与枚举口）＋`ui.go:146`
`Allow(ctx context.Context, correlationID, grant string) error`。⇒ 答复侧确实递不进 tool/args/level/seq，
"从答复侧重算摘要"必须先扩 C17 入向面＝契约级。⛔ 我没有扩任何面，也没有为了自洽去动 ⓐ 的措辞——
**§2 那八枚点没有一处需要同时改 ⓑ 才成立**（这一点按派单要求具名报：无冲突、未停手）。

**5.5 跨包读数未量。** 见 §3.3 末段。派单⛔ 禁止本腿跑 `cmd/wisp`／`internal/panel`／`internal/tools`／
`internal/config`，也⛔ 不许种影响 `cmd/wisp` 读数的突变，所以"注释级改动对别包可见面＝零"是推断。

**5.6 `approval.go` 里四枚**预存**的带圈字形行（`:93` `形ⓐ`／`:104` `④`／`:131` `①`／`:233` `④`）。
尺＝`grep -nP '[\x{2190}-\x{2BFF}\x{FE0F}\x{1F000}-\x{1FAFF}]'`：锚点版命中 4 枚、现版命中同样 4 枚、
行号一致 ⇒ **本腿新增 0 枚**。这一族的归口仍是 `AGENTS.md` §1.2 那条"规格比仪器宽"的未定案缺口
（票 141 (b) 支＋台账 `A201②`），⛔ 不由本腿顺手清（改它们＝动别人那几刀的口径）。我自己在 §2 初稿
里用过 `①②`，已按"就严的一方"改成 ASCII（见 §7 第 5 条）。

---

## 6. commit 名册（按序；每笔都带显式 pathspec，`git commit -F msg -- <paths>`，⛔ 无 push）

| 笔 | 哈希 | 内容 | 入库面（`git show --name-only`） |
|---|---|---|---|
| 1 | `549cb218` | 证据件骨架 §0-§1＋起手快照 `before/approval.go`／`before/queue.go`＋`msg-s1.txt` | 仅本目录 4 枚 |
| 2 | `6a021830` | ⓐ 落地第一批：两枚产码文件的注释降级＋`msg-s2.txt` | `internal/agent/approval/approval.go`、`queue.go`、本目录 msg |
| 3 | 本笔（哈希见 §8 末与交件回报） | §2-§8 填实＋读数件（`base-v/base-names/after-v/after-names/comment-table/comment-diff/gofumpt-negctl`）＋`msg-s3.txt`；同笔带上两枚产码文件的**收尾注释**（`approval.go` 第 385 行那枚 `//` 分段线＋B 枚点撤掉逐字引用，两笔合计 +80/-7 与 +36/-4） | 两枚产码文件＋本目录 |
| 4 | 末笔（随交付回报给哈希） | 只改本件：把 3／4 两笔的哈希按序写回 §6（⛔ 不动产码） | 本件＋`msg-s4.txt` |

起手锚 `5b638498`（`git log -1 --format=%h` 自取，非抄派单）；`git status --porcelain -- cmd internal tools`
起手为**空输出**，第 2 笔之后现跑（12:48 前后）＝只含我自己那一枚待收尾的产码文件：

```
$ git status --porcelain -- cmd internal tools
 M internal/agent/approval/approval.go      （第 3 笔带的注释收尾：对 HEAD 是 +6/-5，全是注释行）
$ git diff --numstat HEAD -- internal/agent/approval/approval.go internal/agent/approval/queue.go
6       5       internal/agent/approval/approval.go     （queue.go 对 HEAD 零差＝第 2 笔已交）
```

别人的脏面：派单指定的那把 scoped 尺（`-- cmd internal tools`）起手与现跑都**没有**别人的东西，
⇒ 本腿没有该档的脏面可报，⛔ 也没清过任何东西。`257-v1` 在 `cmd/wisp` 做的是终裁读数（那档此刻干净）。
但整仓 `git status --porcelain`（我跑了一次，只为确认自己没卷走别人）＝**共享工作树本来就脏**，
逐样登记、⛔ 不动：`M .gitignore`；`D design/assets/base.css`／`icons.js`／`theme.js`／`tokens.css`／
`design/index.html`／`design/screens/*.html` 等十枚 delete；`M design/doubao/**` 四枚；
`M .scratch/wisp/probes/152/my152.py` 与 `161/r6/logs/flip-*.txt` 八枚、`268/a2/census.md`；
未追踪里有一枚**文件名叫 `-`** 的（多半是谁把 `-` 当成了 pathspec 留在仓根）＋`.scratch/.scratch/`＋
`.scratch/ci-logs/run-*.log` 若干＋`.scratch/commit-msg-167*.txt` 八枚。⇒ 这正是本腿 commit 一律带
**显式 pathspec**（且写在 `git commit -F` 上）的原因：index 是共享的，`A617` 那笔事故就是这么来的。
第 3 笔之后 `git diff --numstat 5b638498..HEAD` 应读到 `80 7`（approval.go）与 `36 4`（queue.go），
＝§3.1 那两行交付态读数。

---

## 7. 自我对抗（≥5 条，全部是真改，逐条给出"我差点写成什么／盘上为什么不允许／改完落在哪一行"）

1. **state 检查被我一开始写成第二枚"挡跨卡"的防线。** 盘上不允许：票 259 §9 我复认过的更正写明跨卡时
   B 本来就 pending，`it.state != statePending` 放行；它拦的是"牌落在已离开 pending 的卡上"。
   改后落点＝`approval.go:440-444`（"not a cross-card attempt, because the other card is sitting there
   pending"）与 `queue.go:59-62`（guard 1／guard 2 各自 landing 分开写）。
2. **残余风险初稿举了"store 被多张卡共享"那一形。** 现算否掉：共享 store 里 stored＝A 的摘要、bind＝
   B 的摘要（调用点仍传 `it.bind`，`it` 是 B），那句比较会**不等**、绑定反而抓得到那一形 ⇒ 例子里
   兜不住的两支只剩"成员扫描被松散成认牌即放行"与"路由不带 state 检查"，其余一律删掉。
   改后落点＝`approval.go:543-548`。
3. **初稿写"This function has ZERO independent power over a grant aimed at the wrong card"。** 那是把
   guard 1 一起否了——成员扫描就在这个函数体内。改成"the binding **line** has ZERO independent power;
   the guard inside this function is that scan"。落点＝`approval.go:538-540`。
4. **初稿把撤销口令转写成英文当原文引用**（`revocation phrase "259 AC#1 form (b)"`）。盘上
   `A619` §2 的原文是「259 AC#1 改 ⓑ」，转写＝假引用（本仓为注释里的假引用付过学费）。改成
   "which also carries the revocation phrase for that choice"，让读者去台账读原文。落点＝`approval.go:371-372`。
5. **初稿用了带圈数字 `①②`。** `tools/d22scan` 不扫 U+2460-U+24FF（AGENTS.md 明写的"刻意留的空隙"），
   但 PLAN.md D29 的规格射程写作 U+2190-U+2BFF、**含**带圈数字 ⇒ 仪器抓不到不等于可以写。改成
   ASCII `guard 1 / guard 2`，并把整批新增行扫了一遍（§4 末＋§5.6 的尺）。
6. **初稿沿用票面现量 5 的"六枚用例"口径。** 票 §9 已复认真值＝`ticket242_binding_test.go` 只有
   **4 枚 `Test*`**（我现跑 `grep -c "^func Test"`＝4，与它自己提交说明 `4bf7e683` 的"四枚用例"一致）
   ⇒ 我没在任何注释里写枚数，只点了两枚逐枚验过真在的用例名（`ticket242_binding_test.go:26`／
   `ticket259_denial_rulers_test.go:165`，都在基线名册里 PASS）。⛔ 没写不存在的名字。
7. **初稿写"no assertion in this package notices that refactoring"。** 那是一句未跑的突变断言（本腿
   不许种突变）。拆成两半：能由读码证真的留下（"直调 store、手递两枚摘要的用例分不清路由恒等与路由
   不等"），不能证真的标成引用（"A619 §2 registers it; this leg seeded no mutation"）。
   落点＝`approval.go:548-555`。
8. **`spend` 注里我先写"the membership scan above"。** 那圈扫描在函数体里、在这段注的**下方**，指代反了
   ⇒ 改 below（`approval.go:543`）。同一类的还有 `bindDigest` 那句"the previous sentence claimed"
   指代不清 ⇒ 改成逐字引用被撤的原话（`approval.go:410-411`）。
9. **权威三条初稿只是把原第 (3) 句删掉。** 那会留下"权威＝路由＋牌活着"这种更弱的读法，下一位还得翻
   台账才知道谁在真正拦 ⇒ 重排成三条并把 (2) 指到具体形状（该卡**自己** store 里的活行）、
   点名跨卡落点 `denialSpentNonce`（`approval.go:357-362`），再单独一段写"摘要不是第四道门"
   （`:364-374`）。
10. **我差点为了"讲清楚"而提一句"ⓐ 之后这层只剩审计意义"。** 那是边界①禁的过头话（等于暗示可删）
    ⇒ 留下的是可证的一侧：删除动作在行上、行上有内容物、ⓑ 要比就得有这个东西
    （`approval.go:376-384`），并且⛔ 全篇没有任何"死代码／可以删"的字样。
11. **我在 `bindDigest` 那枚点逐字引用了被撤的假话**（`("ties a grant to ONE pending item")`），
    理由是"让读者知道撤的是哪句"。盘上后果：票面 AC#1 的验收句是一句 grep 级禁令，
    一句"引用＋否认"仍会被同一把尺命中，而下一位拿尺来查时看到的是命中而不是解释 ⇒ 撤掉逐字引用，
    只留"这里的措辞曾那样自称、票 259 AC#1 形(a) 把它降级了"（`approval.go:410-412`），
    原话全文改放在本件 §2 的 before 块里。同枚点还有一次：我先写"the wording **here** used to say"，
    而"here"当时指的是**已经不存在**的那行 ⇒ 改成"the wording **on this function** used to claim"。
12. **我把 §8 的自查尺写成了"＝0"而不是跑一遍。** 现跑之后 G3（过头话尺）真实命中 1 行＝边界①那句
    自己（"not the same as unused"），G2（反向尺）命中 2 行都是否定形——这些读数写进了 §8，
    ⛔ 没有为了好看留"0 命中"那句。这类"判据写得太漂亮"是本仓那条老形状（仪器射程≠结论射程）。
13. **占位尺量到了我自己。** 第一版 §8 把派单那把占位尺的三个模式逐字写进句子里，
    `grep -cE` 当场回 1 枚命中＝本件 627 行那一行自己 ⇒ 改成描述形状不重打模式（§8 末条），
    并把这条留在这里，因为它是"引用一把尺的图案＝污染这把尺"的同族形状，与本票 §1 那三处
    "把恒等式当成防线"是同一类错误的文档版。

---

## 8. 交件判语（不自翻勾；判语归非实现者 `259-v1`）

- **ⓐ 的三条判据自查**（四把尺全部现跑，读数是跑出来的不是推的）：
  ① "明写零独立拦截能力"＝`approval.go:364`＋`queue.go:53`＋`queue.go:406`（三处同口径，读码即得）；
  ② "指名真防线在哪两处"＝`approval.go:436-444`（store 成员＋队列 state，各自拦什么写清）＋
     `queue.go:55-62`（guard 1／guard 2＋landing 差别）＋`queue.go:404-413`（调用点旁边再指一次）；
  ③ "不许留一句能挡跨卡在（这两枚）码里"＝正尺
     `grep -rniE "(stops|prevents|blocks|refuses|cannot).{0,60}(cross-card|another (item|card)|different (item|card|request))" <两枚文件>`
     ⇒ **rc=1（零命中）**。反向尺（bind/digest 后 120 字内出现保护性动词）⇒ 命中 **2 行**＝
     `approval.go:364`（"it stops nothing across cards"，否定句）与 `approval.go:444`
     （"qitem.bind carries both guards"，指的是别处的两枚 guard），逐枚现读已核，都不是漏网的肯定句；
     第四把尺查被撤的原话是否还留在码里（`must match to authorize`／`ties a grant to ONE`／
     `cannot be replayed onto a different`）⇒ **rc=1 零命中**（我第一版在这里逐字引用过原话，已撤，
     见 §7 第 11 条）。原话全文只留在本件 §2 B/F 的 before 块里（证据件不是码）。
- **边界①自查**：⛔ 无"死代码／可删／没用"的**主张**。尺＝
  `grep -rniE "dead code|unused|can be deleted|no longer needed|obsolete|pointless" <两枚文件>`
  ⇒ 命中 **1 行**＝`approval.go:376` "Downgraded is not the same as **unused**, and this is the boundary
  of that claim"，即边界①这句话本身（否定形），不是过头话；除此之外零命中。"仍是单用记录的内容物＋
  ⓑ 唯一可比之物"写在 `approval.go:376-384`。
- **行为零改动**＝三样齐全且互相独立：numstat 删除列 7＋4 全落注释行（§3.1／§3.2）；
  注释剥离后 md5 逐枚相同（§3.2）；单包名册作差零差异 70/89/0/1（§3.3）。四门 rc=0 带时刻（§4）。
- **本腿未做**：⛔ 未做 ⓑ（未扩答复面、未动 `C17` 名册）；⛔ 未动 `ui.go`；⛔ 未动任何测试文件；
  ⛔ 未动票面 AC 框与台账；⛔ 未跑别包尺；⛔ 未 push；⛔ 未种突变。
- **待人裁的三格**（§5.1／§5.2／§5.3）：`ui.go:127-129` 那句"bound to another item"是否豁免；
  测试面那族"cross-item binding is broken"是否算 AC#1 验收句的射程（动它就动名册）；
  边界②那一格装牙归 `259-v1`。
- **件尺**：本件行数／字节数与占位尺随第 3 笔一起报；占位尺＝派单给的那把 `grep -cE` 三模式尺
  （两个"待＋一字"形、一个"填写＋一字"形、一个"未＋一字"形），⛔ 我不在本件里逐字重打那三个模式，
  否则这把尺会量到自己（第一版就中过一次，见 §7 第 13 条）。实测读数写在 §9。

---

## 9. 件尺与读数件名册（终值在本节最后一行，与第 3 笔同批入库）

| 尺 | §0-§8 定稿时（12:49 前后，第 3 笔前） | 交付终值（第 4 笔，只改本件） |
|---|---|---|
| `wc -l .scratch/wisp/probes/259/r2/evidence.md` | 632 行 | 本笔写入后重跑＝**见下** |
| `wc -c` | 41,812 字节 | **见下** |
| 占位尺（派单那把 `grep -cE`） | **0 枚**（rc=1＝零命中） | **0 枚** |

读数件名册（全在本目录，⛔ 无一枚落在别人的面上）：

| 件 | 是什么 |
|---|---|
| `before/approval.go`、`before/queue.go` | 起手快照（md5 `580bddb9…`／`eace17e8…`，与工作树逐字节同）＝§2 全部 before 的来源 |
| `base-v.txt`／`base-names.txt` | 改动前的 `-v` 全文与 70 枚顶层名名册（12:22:36→12:22:38，ok 0.456s） |
| `after-v.txt`／`after-names.txt` | 交付态的同一对（12:45:17→12:45:20，rc=0），与 base 作差零差异＝§3.3 |
| `comment-table.txt` | §3.2 那把逐枚对照表的原始逐行输出（116 枚新增行逐枚带行号＋前 24 字＋判定） |
| `comment-diff.txt` | 两枚产码文件 `5b638498..工作树` 的 `-U3` 全文 diff（191 行／12,101 字节）＝§2 的可核对底稿 |
| `gofumpt-negctl.txt` | 12:47:37 那一次并排跑：仓外负控被点名、我改的两枚零行（§4 的空读数凭据） |
| `msg-s1.txt` … `msg-s4.txt` | 四笔 commit message 的原文（中文，⛔ 无英文单引号） |

交付终值（把这串数字写进去之后重跑那把 `wc` 的读数）＝**657 行／43,526 字节**，行数不动、
字节只随数字本身变，占位尺同轮再跑仍 **0 枚**。⛔ 这里不留成语也不留下划线。
