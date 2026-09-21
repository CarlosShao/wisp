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

<!-- CONTINUE-HERE -->
