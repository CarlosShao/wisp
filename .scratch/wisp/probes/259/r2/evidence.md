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
