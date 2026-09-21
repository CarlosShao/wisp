# 票 87 独立对抗验收 —— `acceptor-ticket87`

**验收对象**：`.scratch/wisp/issues/87-l2-card-veto-cannot-find-its-own-item.md`
**被验收代码**：`dfe9d9f`（AC#1 只读调查）→ `1068eb9`（AC#2 实现）
**本代理角色**：独立对抗验收。**不写生产码、不 push。** 票面叙述属"实现者自述"，本文所有数字均为本代理自己的读数。
**验收日期**：2026-09-21 · **主机**：Windows / go1.27.1 windows/amd64
**快照（全部在仓外，A38④）**：
- 主快照 `/tmp/wisp87acc-87` ← `git archive 1068eb9 | tar -x`
- 对照快照 `/tmp/wisp87acc-87ctl` ← `git archive dfe9d9f`
- 对照快照 `/tmp/wisp87acc-87ctl2` ← `git archive 63ef895`

**变异纪律**：三组变异**全部在快照目录里做**，工作树一字未动。每组跑前 `grep` 自证锚点已改，
跑后 `sha256sum` 与 `git show 1068eb9:<file>` 逐字节比对（快照无 `.git`，故用 blob 哈希替代
`git diff --quiet`，等价且更强：证的是"与 commit 内容一致"而不是"与工作树一致"）。

---

## 1. 基线（AC#5 的测试半边 / 前置门）

在 `/tmp/wisp87acc-87` 纯净快照跑 `go test -count=2 -v ./internal/agent/approval/`：

| 指标 | 本代理读数 |
| --- | --- |
| rc | **0** |
| `=== RUN` 总数 | **84**（顶层 56 + 子用例 28） |
| `--- PASS` | 顶层 **54** + 子用例 **28** = 82 行 |
| `--- FAIL` | **0** |
| `--- SKIP` | **2** |
| 耗时 | 1.409s（`ok github.com/Carlos/…/internal/agent/approval`） |

**`=== RUN` == 2 × 不同测试名？成立**：42 个不同名（顶层 28 + 子用例名 14）× count=2 = **84**，逐项核对过。
顶层读数自洽：54 PASS + 2 SKIP + 0 FAIL = 56 = 顶层 `=== RUN`（无斜杠）数。

**SKIP 逐条点名（2 条同一名）**：

1. `--- SKIP: TestDefaultDeadlineWallClockMeasurement (0.00s)`
2. `--- SKIP: TestDefaultDeadlineWallClockMeasurement (0.00s)`

= 票 84 的 300s 墙钟计量用例，默认 `WISP_84_MEASURE` 未设 ⇒ **有意的慢跳过**，不是"用 SKIP 换绿"。
`no tests to run` / `testing: warning: no tests to run` 未出现。

**结论：PASS**（与票面 AC#5 自述的 84/42/0/2 四个数**逐项一致**）。

---

## 2. AC#3(i-a) 变异：把 `Veto` 的 L2 分支变成"找不到条目"

**锚点**：快照 `internal/agent/approval/gate.go:426` —— 已核对原文件该行确为
`	return g.q.reject(v.CorrelationID,`（票面给的 file:line 准确）。

**落地证明**（跑前 grep 同链自证）：
```
426:	return g.q.reject(v.CorrelationID+"-mut87ia",
原锚点计数 grep -c 'return g.q.reject(v.CorrelationID,$' = 0
```
即：给 corr 加 `-mut87ia` 后缀 ⇒ 该分支**存在但永远解析不到条目**，等价票面声称的"这条分支不存在"。
（注：字面上把 `-mut87i` 直接粘进标识符会编译失败，按规矩"编译失败不算变异"，故取语义等价的可编译形。）

**读数 A —— 过滤到新用例族**（`-run 'TestAVeto|TestAnAmbiguous|TestL1WindowVetoOnForeign'`，5 条含 1 条既有）：
rc=**1**，`=== RUN` 5 / PASS 3 / **FAIL 2** / SKIP 0
```
--- FAIL: TestAVetoAgainstAnOnScreenL2CardEndsTheWaitNowAndStillRefuses (3.01s)
    ticket87_veto_l2_test.go:88: Veto on the displayed card: approval: correlation_id 无对应待审批项（票 87 的原始症状就是这个调用返回查无此项）
--- FAIL: TestAVetoNamedByTheHostsOwnKeyStillRefusesThatCard (3.00s)
    ticket87_veto_l2_test.go:145: Veto by the host's own key: approval: correlation_id 无对应待审批项
--- PASS: TestAnAmbiguousVetoNameIsNotGuessedOnTheRefusalSide (0.00s)
--- PASS: TestAVetoOnAnUnknownNameLeavesTheCardAllowableByItsOwnGrant (0.00s)
--- PASS: TestL1WindowVetoOnForeignCorrelationCannotCancel (0.00s)   ← 既有用例，本就不依赖此分支
```

**读数 B —— 全包**（`-count=1`，无过滤）：rc=**1**，`=== RUN` 42（顶层 28）/ PASS 顶层 25 / **FAIL 2** / SKIP 1。
红的正是上面那 2 条，其余 25 条顶层绿 ⇒ **变异是外科手术式的，不是整体崩溃**。

**"该红的红、不该红的确实还绿"：验毕。** 票面 AC#3(i-a) 的 2 红名单与本代理读数**逐名一致**，
且"另两条 PASS 本来就不依赖这条分支"这一句**成立**（歧义不许猜 / 无关 id 不碰别人的卡 两条仍绿）。

**结论：PASS**

---

## 3. AC#3(i-b) 变异：`return g.q.reject(...)` → `return nil`（听见了但什么都没做）

**落地证明**：
```
426:	return nil
427: }
```
`go build ./internal/agent/approval/` rc=**0** ⇒ **不是编译失败**，是真行为变异。

**读数**（同 5 条过滤族）：rc=**1**，`=== RUN` 5 / PASS **0** / **FAIL 5** / SKIP 0
```
--- FAIL: TestAVetoAgainstAnOnScreenL2CardEndsTheWaitNowAndStillRefuses (5.01s)
    ticket87_veto_l2_test.go:91: veto against an on-screen L2 card: STILL BLOCKED after 2s (start=15:00:44, now=15:00:46) - 闸门没有上界
--- FAIL: TestAVetoNamedByTheHostsOwnKeyStillRefusesThatCard (5.00s)
    ticket87_veto_l2_test.go:147: veto addressed by the task id: STILL BLOCKED after 2s (start=15:00:49, now=15:00:51) - 闸门没有上界
--- FAIL: TestAnAmbiguousVetoNameIsNotGuessedOnTheRefusalSide (0.00s)
    ticket87_veto_l2_test.go:182: ambiguous veto err=<nil>, want approval: correlation_id 无对应待审批项（两名同时有效时不许猜一个）
--- FAIL: TestAVetoOnAnUnknownNameLeavesTheCardAllowableByItsOwnGrant (0.00s)
    ticket87_veto_l2_test.go:214: err=<nil>, want approval: correlation_id 无对应待审批项
--- FAIL: TestL1WindowVetoOnForeignCorrelationCannotCancel (3.00s)   ← 既有用例
    window_test.go:155: err=<nil>，期望 ErrUnknownCorrelation
```

**票面点名的两条红（`:91`、`:147`）逐字复现**，红的位置确在**等待上界断言**
（`STILL BLOCKED after 2s`，来自 `ticket84_no_owner_test.go:40` 的 `mustAnswerWithin`），
不是编译失败 ⇒ 验毕"(i) 那半边真的在测'立刻结束'，而不是只测错误码"。

⚠ **票面数字偏差（对本代理有利，但对票面不利）**：票面 AC#3(i-b) 自称"2 红"。
本代理复跑同一族是 **5 红** —— 因为 `return nil` 顺带让"查无此项必须回 `ErrUnknownCorrelation`"
这一族的三条断言也全塌（含既有用例 `window_test.go:155`）。
即：本票的用例集对这一变异的覆盖面**比票面声称的更宽**。判 **PASS（票面该处少报了 3 条红）**，
残留：票面 AC#3(i-b) 的红名计数不准。

---

## 4. AC#3(ii) 变异：把拒绝漏斗的投递方向改开放行 —— **本票的牙齿**

**锚点**：快照 `internal/agent/approval/queue.go:383`，`Queue.reject` 的投递行（票面给的 file:line 准确）。

**落地证明**（同链 grep，跑前）：
```
383: 	if !q.deliver(it, answer{a: tools.AnswerAllow, why: why}) {   ← 改的就是这一行
407: 	a := answer{a: tools.AnswerReject, why: why}                  ← 超时侧一行未动
queue.go 里 tools.AnswerAllow 出现次数：1 → 2（383 新增，343 是 allow() 原有）
```

**读数 —— 全包 `go test -count=1 -v ./internal/agent/approval/`**：
rc=**1**，`=== RUN` 42（顶层 28）/ PASS 顶层 23 / **FAIL 4** / SKIP 1
```
--- FAIL: TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis (0.00s)
    queue_test.go:143: answer="allow"，期望 reject：面板只能拒绝
--- FAIL: TestReplayRedisplaysUnderAFreshGrant (0.00s)
    queue_test.go:280: answer="allow"，期望 reject
--- FAIL: TestAVetoAgainstAnOnScreenL2CardEndsTheWaitNowAndStillRefuses (0.00s)
--- FAIL: TestAVetoNamedByTheHostsOwnKeyStillRefusesThatCard (0.00s)
    ticket87_veto_l2_test.go:149: answer=allow why="用户在 L2 审批卡片显示期间通过「按 Esc 键」否决了本次操作，未执行", want reject
```

**"必须有既有用例红"：成立。** 本代理独立核过"既有"这条属性，不看票面自述：
- `TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis` 定义在 `queue_test.go:89`、
  `TestReplayRedisplaysUnderAFreshGrant` 在 `queue_test.go:269`，两者在**本票实现 commit 之前**的树
  （`dfe9d9f`）里已存在；
- `git show --name-only 1068eb9` 只碰 `gate.go` / `queue.go` / `ticket87_veto_l2_test.go` + 票面
  ⇒ `queue_test.go` **一字未改**，这两条用例不可能是本票为自证而新加的。

**旁证（变异确是外科手术）**：`TestL2QueueAutoRejectsAt300sKeepsTaskAliveAndWarnsAt270s` 在这一变异下仍 **PASS**
——超时侧（`queue.go:407`）不在变异范围内，红只落在"人说了不"这一条投递上，正是票面
"the single funnel every refusal route goes through" 的声称形状。

**改开放行没有全绿 ⇒ 底线守住了。** 结论：**PASS**（4 红、2 条既有，与票面读数**逐名一致**）。

---

## 5. 变异还原证明

三组变异各自跑完后 `cp` 回备份，随后对 `1068eb9` 的 blob 逐字节比对：

| 文件 | 快照 sha256[0:16] | `git show 1068eb9:` | 一致 |
| --- | --- | --- | --- |
| `internal/agent/approval/gate.go` | `7cfd9ab6ccfa4aee` | `7cfd9ab6ccfa4aee` | ✅ |
| `internal/agent/approval/queue.go` | `a69398b0c17b68c3` | `a69398b0c17b68c3` | ✅ |
| `internal/agent/approval/ticket87_veto_l2_test.go` | `92951cbd6b19fbea` | `92951cbd6b19fbea` | ✅ |
| `internal/agent/approval/queue_test.go` | `04078f47c29e3c03` | `04078f47c29e3c03` | ✅ |

还原后 `go test -count=1 ./internal/agent/approval/` rc=**0**。
工作树未被本代理用于任何变异（`git status` 见 §9）。

---

## 6. AC#2(ii) 读码攻击：别名索引能不能被用来放行 / 命中别人的卡

### 6.1 别名索引的可达面（本代理自己的 grep 读数，非票面）

在 `1068eb9` 快照里，`q.alias` 全仓**只有 5 处引用**，全部落在 `queue.go` 内：
`:196`/`:199`（`indexLocked` 写）、`:209`/`:215`（`unindexLocked` 删）、**`:240`（唯一的读，`resolveLocked` 的宽松分支）**。
`it.names` 的非测试读也只 2 处（`:195`、`:208`，即索引与反索引自身）。

`resolveLocked` 的调用点**只有一个**：
```
queue.go:230  func (q *Queue) resolveLocked(name string, strict bool) *qitem
queue.go:374      it := q.resolveLocked(corr, false)      ← 在 reject() 里
```
⇒ **`strict=true` 零调用点。** 注释里那句 "an allow passes true" 描述的是一个**不存在的调用方**：
`allow()`（`queue.go:326-345`）压根不走 `resolveLocked`，它直接读 `q.byID[corr]` + `it.grants.spend(nonce, it.bind)`。

**"只服务拒绝方向"这句话到底靠什么成立？** 本代理的判定：
- ✅ **结构上成立**（不是只写在注释里）：`alias` 字段非导出、唯一读者是 `reject()` 一条路径，
  而 `reject()` 的全部 5 个上游都是拒绝路线——`gate.go:426`（Veto）、`:584`（`nativeAPI.Reject`）、
  `:588`（`panelAPI.Reject`）、`:613`（`DecideFromNative` 的 `!r.Allow` 支）、`:634`（`DecideFromPanel` 的 `!r.Allow` 支）。
  允许路线 `q.allow` 从不读这张表。
- ❌ **类型上没有表达**：承载方向性的只是一个 `strict bool`，且它**单侧未使用**。
  今天没有任何编译期屏障阻止后来者写 `q.allow` 里 `resolveLocked(corr, false)`，
  也**没有任何既有测试**钉住"别名买不到批准"——本代理 grep 过 `internal/agent/approval/*_test.go`
  里的 `alias|otherNames|names:`，只有 `ticket87_veto_l2_test.go:133/140` 的两处路径字符串，**零条断言**。
  ⇒ 这是本票最值得记的一条残留（见 §10）。

### 6.2 四次具体攻击（探针写在仓外快照，跑完已删除；不属被验收代码）

`go test -count=1 -v -run TestAcceptor87Attack` → **4 PASS，0 FAIL**，整包 `ok`（0.416s，探针在场也不干扰）。

**攻击 1：拿真实未花费的 grant，用别名（taskID）点名 → 能否买到批准？**
```
Allow(alias="task-87-attack", liveGrant) -> approval: correlation_id 无对应待审批项
DecideFromPanel(alias, allow=true, 真 grant, source="native") -> 非 nil
Allow(alias+"  ", 真 grant)                                  -> 非 nil
三次攻击后 Queue().Depth() == 1（条目未被吃掉、grant 未被别名烧掉）
随后 Veto(alias) -> 立刻结束等待，answer=reject
```
⇒ **别名不能产生 allow**；同一张卡上别名只让拒绝落地。前提也自证了：卡片键 ≠ 别名、卡片上有活 grant。

**攻击 2：一个别名同时是别人的精确键时，这一票落在谁头上？**
构造 A（入键 `shared-name` ⇒ 签发键 `shared-name`）+ B（同入键 ⇒ 撞车分支 `queue.go:138` 重签发为 `shared-name#2`）。
```
head after B queued: corr="shared-name" depth=1 paths=[…87-shadow-A.txt]
Veto("shared-name") -> err=<nil>
A 收到 answer=reject why="用户在 L2 审批卡片显示期间通过「按 Esc 键」否决了本次操作，未执行"
B 仍然 pending，键="shared-name#2"（B 未被动过），随后 B 用自己的键被 reject
```
⇒ 精确键优先于别名，所以命中的是**显示着这个键的那张卡**；B 没有被误拒、也没有被误放。
⚠ 但确实存在一条**可用性/路由残留**：宿主按 B 的**入键** `shared-name` 去取消 B，会落到 A 上
（A 被迫收到一次它没被问过的拒绝，B 继续等到上界）。本代理核对过：这条**先于本票存在**——
`reject()` 在父树上就是 `q.byID[corr]`，同样落在 A。本票只是让 `Veto` 这条以前"什么都不拒"的路线
也开始会拒到 A。方向是拒绝 ⇒ 不触 AC#2(ii) 底线，记 **允许残留**。

**攻击 3：畸形输入能否让 Veto 变成批准/变成静默成功？**
```
空名 / "   " / 卡片键+尾空白 / "no-such-card-87"   -> 全部 ErrUnknownCorrelation
未加载频道 ChannelKWS                              -> 「语音取消不可用: kws」（B1 未被绕过）
零值频道 Channel("")                               -> 「未知取消通道，已忽略:」
每一发之后 Depth() 恒 == 1（畸形否决动不掉条目）
```
⇒ 6 种畸形输入无一返回 nil、无一消耗条目、无一产生 allow。

**攻击 4：两张活卡共享同一个 taskID 别名时，会不会"猜一个"？**
```
Veto("task-text-1" 两名共享的 taskID) -> approval: correlation_id 无对应待审批项
之后 Depth() == 2（两张卡都没被动），View(corr-amb-a)/(corr-amb-b) 均存在
各自用自己的键收尾 -> 两条 answer 都是 reject
```
⇒ `len(set)==1` 的门槛真生效：歧义时**两张都不拒**，fail-closed，符合票面"不许猜"。

### 6.3 `deliver` 的方向中立性（为什么 AC#3(ii) 是本票的牙齿）

`deliver(it, a)` 只做 `it.state = stateAnswered; dropLocked; it.answer <- a`——**它不看方向**。
所以"人说了不"这件事只由 `queue.go:383` 那一个 `tools.AnswerReject` 字面量承载。
本代理在 §4 已用变异证明这个单点由**既有** fail-closed 用例守着。结论：方向性没有被泄漏到调用方各自决定。

**结论：AC#2 PASS**（底线 (ii) 未被任何攻击突破；残留见 §10 的 R-1/R-2/R-3）。

---

## 7. AC#4：300s / 30s 一个数字都没改

**超时常量在快照里的位置与本代理读到的值**（`/tmp/wisp87acc-87`）：
```
internal/agent/approval/queue.go:106   DefaultApprovalTimeout = 300 * time.Second
internal/agent/approval/queue.go:108   DefaultApprovalWarning =  30 * time.Second
internal/agent/approval/queue.go:85    timeout    = DefaultApprovalTimeout   （<=0 回落默认）
internal/agent/approval/queue.go:87-89 warnBefore <= 0 || >= timeout 时回落 DefaultApprovalWarning
```

**与前一枚 commit 逐字对比**（`63ef895` 是 `1068eb9` 的父，见 `git log -1 --format=%P`）：
- 把 `const ( … )` 整块从两棵树里抽出 `diff` → **零差异**（本代理标记 `CONST-BLOCK-IDENTICAL`）。
- `git diff 63ef895 1068eb9 -- internal/agent/approval/` 的**全部删除行只有 8 行**，本代理逐行读过：
  3 行 doc 注释、1 行 `return ErrUnknownCorrelation`、`maxRepl: … byID: …` 一行（只是换行以塞 `alias:`）、
  `it, ok := q.byID[corr]` + `if !ok {` 两行（换成宽松解析）。**没有一行含数字。**
- `git diff dfe9d9f 1068eb9 --name-only -- docs/` → **空** ⇒ `docs/PLAN.md`（C18 `:3143`/`:1368`）与 `docs/specs/**` 一字未动。
- 新增测试文件里出现的 `30 * time.Second` / `25 * time.Second` / `vetoBudget = 2s` 是**测试自设的门限与断言上界**
  （把 `Options.ApprovalTimeout` 临时 arm 成 30s 以免变异真等 300s），不是生产默认值，也不是被放宽的判据。

**"超时前 30s 醒目提示"是否已实现：已实现**（票面 AC#4 的这一说法成立）：
```
gate.go:477-480  lead := g.q.Timeout() - g.q.WarningLead(); warn = g.clock.After(lead)   → 300-30=270s 起
gate.go:507-519  case <-warn: 发 EventWarning，Text=「审批将在 %d 秒后自动拒绝，请尽快确认」，且 warn=nil（只发一次，不缩短截止）
ui.go:65         EventWarning EventKind = "warning"
cmd/wisp/run.go:538-541  consoleApprovalUI.Update 打印 "[%s] %s" → 运行时字面输出 "[warning] 审批将在 30 秒后自动拒绝…"
queue_test.go:189  TestL2QueueAutoRejectsAt300sKeepsTaskAliveAndWarnsAt270s —— 且在 :212-214 断言
                   「270s 之前不该弹出醒目提示」（clk.Advance(270s-1ms)）⇒ 上界被测试钉住，不只是注释
```
该用例在本代理基线（§1）与变异 (ii)（§4）下均 **PASS**。
注意 `[warning]` 在源码里**不是字面量**（本代理 `grep -rn "\[warning\]" --include=*.go .` → 0 命中），
它是 `e.Kind` 经 `%s` 渲染出来的 ⇒ 票面"会打印 `[warning] …`"成立，但按字符串搜是搜不到的。

⚠ **一处行号漂移（纯记账）**：票面 AC#4 写 `queue.go:92-94` "原样"。在**父树**里 92/94 正好是这两个常量（本代理核过：
parent 90 `const (`、92 `DefaultApprovalTimeout = 300*time.Second`、94 `DefaultApprovalWarning = 30*time.Second`）；
落到 `1068eb9` 后因为上方新增了 14 行，它们搬到 **106/108**。**值一字未改**，票面引用的行号在落地后的树上 stale。
同理票面 AC#4 引的 `gate.go:456-460/487-499` 在落地树里是 `477-480/507-519`。记 **允许残留**。

**结论：AC#4 PASS**（0 个数字被改；30s 提示确已实现并有测试钉住；仅行号引用 stale）。

---

## 8. AC#5 门禁 + 环境事实

| 门禁（快照 `/tmp/wisp87acc-87`，纯净树） | 本代理读数 |
| --- | --- |
| `gofmt -l internal/agent/approval/` | **空**，rc=0 |
| `$(go env GOPATH)/bin/gofumpt -l internal/agent/approval/` | **空**（⚠ 本机 `gofumpt` 不在 PATH，命令直接 127；实际二进制在 `D:\work\base\gopath/bin/gofumpt`，用绝对路径跑通） |
| `go vet ./internal/agent/approval/` | **rc=0** |
| `GOOS=linux go vet ./internal/agent/approval/ ./internal/tools/`（按包作用域，未用 `./...`） | **rc=0** |
| `go test -count=2 ./internal/agent/approval/` | **rc=0**，`=== RUN` 84 / PASS 54(+28 子) / FAIL 0 / SKIP 2（§1 逐条点名） |
| `go test -count=2 ./internal/tools/`（下游回归） | **rc=0**，31.424s `ok` |

**票面声称的环境事实（`go test ./cmd/wisp/` rc=1 = 加载期 `0xc0000135`，先于本票存在）—— 本代理独立复现：**
```
主快照   1068eb9  : go test ./cmd/wisp/                rc=1  exit status 0xc0000135   0.057s
主快照   1068eb9  : go test -run '^$' ./cmd/wisp/      rc=1  exit status 0xc0000135   0.054s   ← 一条用例都不执行也红
对照快照 63ef895  : go test ./cmd/wisp/                rc=1  exit status 0xc0000135   0.132s   ← 本票落码之前
对照快照 63ef895  : go test -run '^$' ./cmd/wisp/      rc=1  exit status 0xc0000135   0.114s
对照快照 63ef895  : go test ./internal/agent/approval/ ok  0.457s                    ← 同棵树里 approval 包是好的
```
⇒ **"不是票 87 弄坏的"这条结论成立**（同一台机器、同一枚父 commit 同读数）。
本代理再往下挖了一层根因，供编排者登记：
```
go test -c -o x.test.exe ./cmd/wisp/  → 编译 rc=0（产物 29,939,018 B）
ldd x.test.exe                        → sherpa-onnx-c-api.dll => not found
```
即 **cgo 依赖的 sherpa-onnx DLL 不在本机 PATH**（A54③ 同族），纯加载期，与测试逻辑无关。

**结论：AC#5 PASS**（票面五个门禁读数本代理全部独立复现；环境事实复现且根因已定位）。

---

## 9. 规则遵守与本代理的自证

- 三组变异**全部在仓外快照**（`/tmp/wisp87acc-87`、`/tmp/wisp87acc-87ctl`、`/tmp/wisp87acc-87ctl2`）完成；
  工作树 `internal/agent/approval/` 本代理一字未写（`git log 1068eb9..HEAD -- internal/agent/approval/` 为空 ⇒
  验收期间也没有别人碰过这个包，本代理的快照读数对当前 HEAD 仍有效）。
- 攻击探针文件在跑完后**已从快照删除**；未在任何仓内目录留下非证据文件。
- 未放宽任何断言/阈值/golden；未跑 `-count` 碰运气（基线一次成样、变异各组各一次成样，全报）。
- 禁改清单（`docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`rules_gateway.go`）：本代理只读，零写入。
- 在途他人文件（`internal/winsec/**`、`internal/panel/assets.go`、`cmd/wisp/panel_assets.go`、
  `.scratch/wisp/issues/{82,88,90}-*.md`）：未读改、未 stage、未提交。
- 本代理**只提交** `docs/evidence/s1/87-adversarial-acceptance.md` 一个文件；
  每次提交前 `git diff --cached --name-only` 核对（本票第一次提交核对输出＝该文件唯一一行）。
- ⚠ 记账：本代理写完第一枚 checkpoint 后，分支被并发推进（票 89/94 的 commit 落在 `867f993` 之上，
  本代理 `git status` 一度看到工作树被他人清干净）。已核 `867f993` 仍是 HEAD 的祖先、
  `git diff HEAD -- 本文件` 为空 ⇒ checkpoint 未被吞。

---

## 10. 最终裁决表（1:1 对齐票面 AC）

| 票面 AC | 判据（票面原文要点） | 本代理独立读数 | 裁决 |
| --- | --- | --- | --- |
| **AC#1** 只读：画出"卡片显示 ↔ 门里条目"键链，列出"卡片在显示、门里查不到"的实例 | 逐处 file:line；查不到要如实写"没找到" | 结构前提独立复现：`openWindow(` 在 `gate.go` 非测试码里**只有 1 个调用点**（`:249`，位于 `PendingWindow`＝L1 路线），`PendingApproval`（L2 路线）**零** `openWindow` ⇒ L2 卡片显示时 `g.windows` 里没有它。症状复现＝变异 (i-a)：把 L2 分支变回"解析不到"后，`Veto` 在**卡片自己印着的键**上返回 `approval: correlation_id 无对应待审批项`（`ticket87_veto_l2_test.go:88`）。卡片确实 advertise 否决频道（`channelAdvertised(p, ChannelBall)` 在基线 PASS）。 | **PASS** |
| **AC#2** 加一条双向可判用例：(i) 立刻结束等待 (ii) 仍按拒绝处理；**绝不因查不到而放行** | (ii) 是安全底线 | 4 条新用例在纯净快照 count=2 全绿；(i) 用 2s 预算对 30s armed 门限（`vetoBudget`），(ii) 同时断 answer 值、拒绝措辞、条目从队列消失、事后 Allow 必须失败。**4 次读码攻击全被挡住**（§6.2：别名+真 grant 买不到 allow；畸形 Veto 无一变静默成功；歧义两名一张都不拒；共享名被精确键正确压制）。 | **PASS**（底线未破） |
| **AC#3(i-a)** "立刻结束"退回"等满上界" ⇒ 用例红 | 该红的红、不该红的还绿 | rc=1；新用例族 2 红（点名 `TestAVetoAgainstAnOnScreenL2CardEndsTheWaitNowAndStillRefuses` + `TestAVetoNamedByTheHostsOwnKeyStillRefusesThatCard`），歧义/无关-id 两条 + 既有 `TestL1WindowVetoOnForeignCorrelationCannotCancel` 全 **PASS**；全包 42 `=== RUN` 里只 2 红、25 条顶层绿 ⇒ 非整体崩溃。 | **PASS** |
| **AC#3(i-b)** `return nil`（听见了但什么都没做）⇒ 红在等待上界断言 | 不许是编译失败 | `go build` rc=0；红名正是 `ticket87_veto_l2_test.go:91` 与 `:147` 的 `STILL BLOCKED after 2s (start=15:00:44, now=15:00:46) - 闸门没有上界`。⚠ **计数不符**：同族本代理读到 **5 红**（票面写 2 红），多出的 3 条正是"查无此项必须回 ErrUnknownCorrelation"一族（含既有 `window_test.go:155`）⇒ 覆盖面比票面声称更宽，票面少报。 | **PASS**（票面该处计数不准，已记） |
| **AC#3(ii)** 失败侧从"拒绝"改开放行 ⇒ **必须有既有用例红** | 全绿即 FAIL | rc=1，**4 红**：`TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis`（`queue_test.go:143: answer="allow"，期望 reject：面板只能拒绝`）、`TestReplayRedisplaysUnderAFreshGrant`（`queue_test.go:280`）—— 两条"既有"属性由本代理独立证明（定义在父树 `dfe9d9f:queue_test.go:89/269`，且 `git show --name-only 1068eb9` 未碰 `queue_test.go`）；另 2 红是本票新用例。旁证外科手术：`TestL2QueueAutoRejectsAt300s…` 在此变异下仍 PASS。 | **PASS（本票的牙齿成立）** |
| **AC#4** 一个数字都不改；30s 提示若未实现只登记 | 300/30 不许调小 | const 块 `63ef895` vs `1068eb9` **零差异**；纯票 87 diff 的 8 条删除行无一行含数字；`docs/` 在两枚 commit 里零改动；30s 提示**已实现**（`gate.go:477-480`+`507-519`，console `run.go:538-541` 渲染 `[warning]`，`queue_test.go:189` 反向钉 270s）。⚠ 票面行号引用 stale（父树 92-94 → 落地树 106/108；456-460/487-499 → 477-480/507-519）。 | **PASS**（残留 stale 行号） |
| **AC#5** 门禁 + 逐条点名 SKIP/FAIL + 环境事实 | gofmt/gofumpt/vet/GOOS=linux vet/count=2 全 rc=0 | §8 表：六项门禁本代理**全部复现 rc=0/空**；SKIP 2 条逐名点出（`TestDefaultDeadlineWallClockMeasurement` ×2，`WISP_84_MEASURE` 未设）；`go test ./cmd/wisp/` rc=1 `0xc0000135` 在 `1068eb9` 与父树 `63ef895` **同读数**（含 `-run '^$'`），根因 `sherpa-onnx-c-api.dll not found`（编译 rc=0）⇒ "先于本票存在"成立。 | **PASS** |

## 11. 残留清单（都不改判据，须被下一张票接住）

- **R-1（最该修，低成本）**：`resolveLocked(name, strict bool)` 的 `strict=true` **零调用点**，
  注释"an allow passes true"描述的是不存在的调用方。方向性今天靠"只有 `reject()` 读这张表"成立，
  靠的是**非导出 + 单一读者**，不是类型；且**没有任何测试**钉"别名买不到批准"（本代理 grep 读数为 0 条断言）。
  建议后续票把 bool 换成两个具名函数（`lookupForAllow`/`lookupForRefusal`）并补一条
  `TestAnAliasCanNeverBuyAnAllow`（本代理 §6.2 攻击 1 的形状可直接搬）。
- **R-2**：宽松解析写在 `Queue.reject()` 里 ⇒ 它的**blast radius 是全部 5 条拒绝路线**
  （native/panel/`DecideFrom*`/`Veto`），不止票面叙述的"Veto 查不到时"。方向仍只可能产生拒绝，
  但票面 AC#2 的叙述范围窄于实际改动面。
- **R-3**：卡片 B 的**入键**与卡片 A 的签发键撞车时，宿主按入键取消 B 会拒到 A（§6.2 攻击 2 实测）。
  方向是拒绝 ⇒ 不触底线，且 `reject()` 路线上先于本票存在；本票新增的只是 `Veto` 这条以前"什么都不拒"的路线
  现在也会拒错。属可用性/路由残留，等票 37 接线时对账。

## 12. 两条最终结论

**(a) 能否改名 `-done`？—— 能，但票面须补一行"用户可见效果未交付"的显式指针，且 R-1 建议进下一张票。**
理由：票面 AC#1-AC#5 五框的**判据**本代理全部独立复现，没有一条靠票面自述；票面对"面板回投递还没接线 ⇒
条件性实例"是**自己先写出来的**（不是被本代理抓出来的），票面也没有任何一处声称"用户现在就能提前拒"。
四框全勾 + 五框判据可复现 + 无禁改越界 + `internal/agent/approval/` 自 `1068eb9` 起无人再碰 ⇒
账目上 `-done` 是诚实的。**唯一要求**：`-done` 之后不得让这张票在台账里被读成"人能提前拒了"——
那件事的开关在票 37/41，本票只交付了它的前置形状。

**(b) "提前拒绝"这条路现在是不是真的可用？—— 库层可用，用户层不可用；本代理判：**就绪层交付，不是半成品骗绿**。
本代理的硬读数：全仓非测试代码里 **`.Veto(` 调用命中 0 处**（`grep -rn "\.Veto(" --include=*.go . | grep -v _test.go`
在 `internal/agent/approval/` 之外零命中），approval 包之外也**没有任何**把 ball/Esc 频道桥到 `Gate.Veto` 的代码；
唯一"Veto"形状的生产引用是 `internal/tools/bridge.go:331 case AnswerVeto:` 与 `:463 b.cancel.Vetoed(...)`，
那是 L1 交棒之后的记账，不是"人按下拒绝"。

给理由而不是给情绪：**"没有生产调用者"和"修得不对"是两件事。** 本票交付的是一条**方向唯一**的库内通道
（`Veto → reject 漏斗 → 唯一一处 `tools.AnswerReject` 字面量），并且它两侧都有活牙口：
(i) 半边变异能红、(ii) 半边由**先于本票存在**的 fail-closed 用例守着——本代理亲测 4 红。
所以它不是"写了个没人调的函数然后打勾"，而是"把一个已经会伤人的错误形状（人按键无响应、被迫等满 300s）
在**唯一一处**改对，并把改对的性质钉在既有安全家族上"。等票 37 接线时，这条路上**不需要重新设计方向性**，
这正是它作为交付的价值。反过来说：若把本票读成"用户现在能提前拒了"，那是**记账误读**，
本代理在 (a) 里要求的那一行就是防这个的。

next=（编排者）①`cmd/wisp` 的 `sherpa-onnx-c-api.dll` 加载期缺失要不要单开票（本代理已给根因与复现读数）；
②R-1（`strict` 死参 + 缺"别名买不到 allow"的测试）；③票 37 接线时把 `internal/panel/approval.go` 的
C17 入键与队列签发键做一次对账（R-3 的根）。

