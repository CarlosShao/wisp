# 简案二：票 259 `AC#2` 两发未种突变——拆 store 成员扫描／state 卫

> 备料腿 `next-instruments-brief-1`（只读）｜取锚 `2026-10-08 19:01` HEAD `5f52f310`（A736）｜两落点行号本程 `sed`／Read 现取（同发取，先定包 `internal/agent/approval`）。
> 出处：台账 `A734`（`docs/reports/pending-and-issues.md:14379`）"残余风险那两发突变未种（拆 store 成员扫描／state 卫——不在派单三枚内，登记为下一枚可选攻法）"；代码侧原有登记＝`approval.go:543-557` 的 Residual risk 段。
> 射程声明：未读未引 `frontend/**`／`design/**`；三枚冻结件未碰（射程判断，非内容引用）。

## 0. 转述差异两处（以盘上原文为准，具名报回）

1. 派单说"state 检查在 `spend`（`:558` 一带）里"——**盘上不成立**：`spend`（`:558-575`）里只有 store 成员检查与绑定比对，**零 state 检查**；state 检查真身＝`queue.go:390`（`allowScoped` 内）。代码注释自己也是这么写的：`approval.go:545-547` 逐字"the queue's state refusals ... an allow route that reaches this function without the statePending test"。⇒ 发 B 落点按盘上原文写 `queue.go:390`。
2. "v1 自报没种的两发"的原句**不在** `probes/259/v1/01-evidence.md`，而在台账 `A734`（`pending-and-issues.md:14379`）；`01-evidence.md` 里对应的是 AC#1 末尾"本轮**不种突变**（文案格，读数即凭据）"（`:34`）与 AC#2"最小形状差异＝0／⛔ 不写实现方案"（`:40`），无"不在派单三枚内"字样。⇒ 两发简案即以 A734 为出处。
（另注：`A737` 已裁票 259 六格全勾、改名 `-done`；本两发若落地，归属是"AC#2 的可选加固"还是新攻法由编排者定——⛔ 本件不新建票/判据。）

## 1. 逐发具名落点（`sed` 现取逐行原文）

### 发 A｜拆 store 成员扫描

落点 `internal/agent/approval/approval.go:564-567`（`spend` 内，`func` 定义在 `:558`）：

```
564	for v, stored := range s.values {
565		if !equalSecret(v, nonce) {
566			continue
567		}
```

- "store 成员检查"＝`564-567`（`for` 扫描 + `if !equalSecret(v, nonce)` + `continue`）；紧接的 `568 delete(s.values, v)`／`569-571 if !equalSecret(stored, bind) { return denialMisbound }`／`572 return denialNone` 是**消费与绑定比对**，不是成员检查；`:559-561 if nonce == "" { return denialMissingNonce }` 是空值卫，不属此发。
- **最小突变（一行，行数不变，留行号利于还原）**：`sed -i '565s/if !equalSecret(v, nonce) {/if false {/' internal/agent/approval/approval.go`。效果＝循环第一枚 entry 被无条件视作命中（delete + 比 bind）⇒ **任何非空值都能消费该 store 的首枚**——正是注释里那句"a spend that answers denialNone on a nonce it merely recognised"（`approval.go:545-546`）。

### 发 B｜state 卫

落点 `internal/agent/approval/queue.go:390`（`allowScoped` 内，`func` 定义在 `:383`）：

```
390	if it.state != statePending {
391		// Ticket 259 AC#2: the denial is named in the audit even on the cause
...
397		q.logf("approval: GRANT-DENY corr=%s tool=%s denial=%s", corr, settled, denialNotPending.label())
398		return ErrNotPending
399	}
```

- 这是"allow 路"上的 statePending 卫；另一半＝`deliver` 的 settled 检查 `queue.go:458`（`if it == nil || it.state != statePending {`，"输掉的那半"）。`queue.go:404-414` 的注释（"this comparison has nothing left to notice"）与 `approval.go:543-557` 同口径。
- **最小突变（一行，行数不变）**：`sed -i '390s/if it.state != statePending {/if false {/' internal/agent/approval/queue.go`。效果＝非 pending 卡也能走到 `:414 denial := it.grants.spend(nonce, it.bind)`。

## 2. 预期红句（今天**有**用例守这两处：逐发"会红"，⛔ 都不是新恒真面）

### 发 A：**会红，两枚，都在路由侧**（判据＝现读断言 + fixture `r1Card` 用 `q.grantNonce` 真发一枚活 nonce，`ticket259_denial_rulers_test.go:93`）

1. `internal/agent/approval/ticket259_denial_rulers_test.go:138-139`（`TestTicket259R1DenialNamesSpentOrNeverLiveNonce` 路由半边）：
   `err := g.Native().Allow(context.Background(), it.Corr, "grant_0000_never_issued")` ⇒ 突变后误中本卡唯一活 nonce、bind 与 `it.bind` 相等 ⇒ `denialNone` ⇒ **allow 成功（err=nil）** ⇒ `t.Fatalf("AC#2 RED: the never-issued-proof route returned %v, want ErrBadGrant", err)` 红。
2. `ticket259_denial_rulers_test.go:258-259`（`TestTicket259R1OutwardAnswerStaysMergedAcrossDenials`）：前一发 `Allow(corr, "")` 空值早退未消费，`Allow(corr, "grant_nope")` ⇒ 突变后同上 ⇒ `t.Fatalf("AC#2 RED: never-live proof returned %v, want ErrBadGrant", err)` 红。

**具名"不红"者（这是 A619 §2 那句"invisible to today's rulers"的准确射程）**：
- `ticket242_binding_test.go:169/:183/:192`（全为直调 store：`itA.grants.spend(nonceA, itB.bind)` 走命中+比对不变 ⇒ misbound；正控 none；空 store ⇒ spent）——绿；
- `ticket259_denial_rulers_test.go:165-181`（misbound 直调用例）——绿。
- ⇒ "rulers 看不见成员扫描之死"说的是**这批直调用例**；路由侧两枚看得见（上面 1/2）。

### 发 B：**会红，一枚**

`ticket259_denial_rulers_test.go:207-208`（`TestTicket259R1DenialNamesCardNotPending`）：卡被手工置 `stateAnswered`、nonce 仍活（fixture 只改 state，`:195-197`）⇒ 突变后走到 `spend` 消费活 nonce ⇒ `if it.grants.live() == 0 { t.Fatal("AC#2 RED: the settled-card refusal burned the live grant; that branch never reaches the store") }` **红**。
⚠ 同时具名"其余断言不红"：`:200-201` 的 `errors.Is(err, ErrNotPending)` 与 `:203-205` 的 `denial=card-not-pending` 审计行在突变后**仍满足**（deliver 输掉的分支 `queue.go:428-429` 返回同 token）——**这枚测试的牙只长在 `live()` 断言上**；种发时别误以为 err/审计两行也在守。

## 3. 还原尺（四件套格式照抄本仓：种前 hash／盘上改行复量／红句 file:line／还原后 hash 等值＋scoped porcelain 空）

两件今天 porcelain 均空（现核），HEAD `5f52f310` 树现值：

| 发 | 目标文件 | 种前 hash（`git hash-object <path>`，本程现取） | 种法 | 改后复量 | 预期红句 | 还原尺 |
|---|---|---|---|---|---|---|
| A | `internal/agent/approval/approval.go` | `67fb146898dc7b235623bb8ad2fac0e0b2465fe1` | `sed -i '565s/if !equalSecret(v, nonce) {/if false {/'` | `sed -n '565p'`＝`\tif false {`，且 `git diff --stat` 只见 1 处 | `ticket259_denial_rulers_test.go:139`／`:259`（断言文本见 §2） | 还原后 hash＝同值（RESTORED_EQUAL）；`git status --porcelain -- internal/agent/approval/approval.go` 空 |
| B | `internal/agent/approval/queue.go` | `66fec7ae375344d7bc741270b42c4ec5050cf41f` | `sed -i '390s/if it.state != statePending {/if false {/'` | `sed -n '390p'`＝`\tif false {` | `ticket259_denial_rulers_test.go:208` | 还原后 hash＝同值；`git status --porcelain -- internal/agent/approval/queue.go` 空 |

写腿执行纪律：①两 hash 为 2026-10-08 19:0x 现取，**种前再取一次为准**（同发取）；②先 `sed -n` 复量行内容再种；③还原用同发留存的副本/逐字修回，⛔ 不 `checkout`；④跑测试带 `-count=1`；⑤基线（未突变）先跑同套 `-run 'TestTicket259'` 应为 12 PASS／rc=0（v1 已给此基线，写腿现跑复认）。

## 4. 判不了的格（照实登记）

1. 红句实参的**终值**（`%v` 展开如 `<nil>`）本腿零 Go 命令，未实跑；断言文本逐字取自文件，实参终值由写腿现跑现取。
2. 发 A 突变后"第一枚 entry"的 map 遍历顺序非确定——今天涉事用例每 store ≤1 枚活 nonce（`r1Card` 一发一枚；242 用例两卡各一枚），不咬；若未来多枚，红句形状可能漂。
3. 未跑 `./internal/agent/approval/` 整包、未核两发是否连带打红别的在册项（本腿只读名册与断言；整包读数归写腿）。
