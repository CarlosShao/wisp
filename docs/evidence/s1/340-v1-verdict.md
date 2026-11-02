# 340-v1 判语件 — 票 340 验收腿（裁决者≠实现者）

- 腿＝`340-v1`（非实现者验证腿）；时间＝2026-11-02 12:3x–13:0x +08；纪律＝只读＋本证据件，⛔改码／⛔改票面／⛔台账／⛔push。
- 被验交付＝`412a58356`（票 340 落地腿 340-r1，新文件 `internal/agent/approval/corr_clash_e2e_340_test.go` 556 行）。
- 票面＝`.scratch/wisp/issues/340-stacked-queue-rewrite-of-a-card-corr-has-no-real-producer-test.md`（43 行，末条 bullet＝340-a1 普查腿 2026-10-18）。

## 起手三核

| 项 | 读数 | 尺 |
|---|---|---|
| HEAD≥`8edbed156` | PASS：HEAD＝`063dcd4d1`，`git merge-base --is-ancestor 8edbed156 HEAD` rc=0 | `git rev-parse HEAD; git merge-base --is-ancestor 8edbed156 HEAD; echo $?` |
| 索引·approval 无脏 | PASS：`git status --porcelain -- internal/agent/approval` 空；`git diff HEAD -- …/corr_clash_e2e_340_test.go`＝0 行 ⇒ 工作树＝交付树 | 同左逐字 |
| 其余脏面归属 | PASS：脏 tracked 件＝8 枚 `internal/config/*`＝345r1 写面，本腿零触碰；`412a58356..HEAD` 动 approval 的提交＝**0 枚** | `git status --porcelain`；`git rev-list --count 412a58356..HEAD -- internal/agent/approval` |

## AC#0 现量复跑 — **PASS**

尺一律取具名提交树（`412a58356^` ⇄ `412a58356` ⇄ `HEAD`）；行号只当读数。

1. **callCorr 现量**：def `internal/agent/loop.go:634`，唯一 dispatch 调用点 `:707`（另两处为注释）。尺＝`git grep -n callCorr HEAD -- internal/agent/loop.go`。与测试头自报（"loop.go:634, sole dispatch call site :707"）一致。
2. **Queue.push 现量**：def `queue.go:159`；撞名条件 `:170` `if _, clash := q.byID[corr]; clash`；改写串 `:171` `corr = fmt.Sprintf("%s#%d", d.CorrelationID, q.seq)`；replay 孪生 `:627`（前置 `:625` `corr+"#replay"`）；`delete(q.byID, it.Corr)` `:312`。尺＝`git grep -n '%s#%d' HEAD -- internal/agent/approval`＋`git grep -n 'delete(q.byID'`。与票面"改写在 `:168` 一带"（census 时锚）仅行号漂移、字节同形。
3. **approval 侧 `#数字` 字面＝0**：`git grep -nE '"[^"]*#[0-9]+"' <tree> -- internal/agent/approval`＝ `412a58356^` **0**（rc=1）⇄ `412a58356` **0** ⇄ HEAD **0**。正控＝同一字节尺只换射程 `internal/panel`＝**13 行**（尺活，且与 340-a1 census 的 13 复认）。⇒ 交付件**没有手写进任何一枚 corr 字面**。
4. **真链＝0（交付前）**：approval 包测试 `agent.New(` 命中＝`412a58356^` **0 枚**（rc=1）⇒ 交付前 approval 侧无任何用例能装配真 loop，"真 Queue.push 撞名改写作用在真 callCorr 产物"必然＝0；既存三真件装配先例 `internal/tools/loop_approval_test.go` 的 `risk.L2|PendingApproval` 命中＝**0**（rc=1，与 census 复认一致＝它走 L1、不经 `q.push` 撞名支，且属票面射程外）。交付树：`agent.New(` in approval＝**1 处**（新件 `:187`）；新用例名 repo-wide 唯一（`git grep -ln TestQueueClashRewriteOfARealLoopCorrelationEndToEnd 412a58356`）。⇒ 交付＝**全仓第一枚真链**，票面"缺世界"前提复认。

## AC#1 真端到端是否真跑拼接 — **PASS**

- **装配全真**：golden SSE（C5  sanctioned 注入口＝唯一 mock 面，给的是 wire 上的 provider call id `call_e340`，⛔ 从不给 corr）→ `httptest` → `llm.NewProvider`（openai-chat 真协议栈）→ 真 `agent.New` loop → 真 `tools.New` bridge → 真 `approval.New` Gate；fs.edit＝Declared `risk.L2` ⇒ 走队列路线。尺＝读件 `composeE2E340`（`:128-208`）＋实测跑（下）。
- **链的每一跳都是产码面**：X 由真 `callCorr`（loop.go:707 呼出）铸出 → 真 `Queue.push` #1 登记 → 测试从产码读面 `Queue.LiveApprovals()[0].Decision.CorrelationID`（pending_read.go:55-59/:113）读回 → 同一枚 decision 原样喂给真 `Gate.PendingApproval`（gate.go:582→:587 `g.q.push`）→ 真 push #2 撞名 → `key1 = X#2`。**第二跳喂的是第一跳的入账值，逐字节，非本件拼写**（字面尺＝0 佐证）。
- **为何第二跳必须经再发（而非同轮两次调用）**：`internal/llm/events.go:204` `byID map[string]int`（同轮同 id 归并、不可能两枚同 id ToolCall）⇄ `queue.go:312`（卡片结算后键摘除、跨轮再发不再撞）——两处字节本腿复看属实；叠改形经"仍活队列上的重发决定"可达＝被测组件（改写支）真实被踩。
- **能否证 `<corr>#<seq>`**：(1b) X＝`<taskID>#<callID>` 单 `#`、前缀 taskID+"#"、≠taskID；(2a) 两项记录同一 incoming corr（撞名见证）；(2b) `key0≠key1`；(2c) `key1` 含两枚 `#` 且前缀 `X+"#"`；(2d) 末段全十进制、`Atoi`＝**2**（本队列两枚 push 的 seq）、base 逐字节＝X；另有名册产码纯跳 `TaskRosterSectionFrom`＋`NativeVerdict.CardView` 两面（both-cards／**stacked-only**＝叠形须自证点亮**原任务行**、sibling 保持暗）、卡片显示 id＝入账键原样、`DecideFromNative` 按精确键拒两卡后队列清空、journal 行记 X（记录侧再钉"剥出的段＝生产 corr、⛔ 裸任务号"）。
- **实测凭据**：`-v` 具名 RUN→**PASS (0.17s)**，日志见真 temp 根＋真 SQLite 迁移四步（`…\TestQueueClashRewrite…\db\backup\wisp.db.bak-*`）＝真跑非跳过。
- **形**：头注自报"甲"（approval 侧、owner-ruled 2026-11-02，A1399 收件通过在册）；盘上实现＝**零新导出名／零产码改动**（比票面甲形的"新导出名"预估更省，经既有导出读面 `LiveApprovals/Depth/PendingApproval` 成链），numstat 单文件 +556/−0 坐实。

## AC#2 守卫与 Skip — **PASS**

- ⛔ 删/放宽既有守卫：`git show --numstat 412a58356`＝`556 0 corr_clash_e2e_340_test.go`，**删除行数 0、触达文件 1 枚** ⇒ 字节面上不存在"删句/放宽"的可能面。
- ⛔ 未登记裸 Skip：新件 `t.Skip|runtime.Goexit|testing.Short` 命中＝**0**（rc=1）；全包 `-v` 实跑 SKIP 仅 **1 枚**＝`TestDefaultDeadlineWallClockMeasurement`（`ticket84_no_owner_test.go`，票 84 既有仪器、非 340-r1 引入——412a58356 未触该文件）。
- 补充禁形读数：新件两枚 `go func(` 均带 `wg.Done`+`recover`（D37b）；`d22scan -root ../../` 全仓 **rc=0 clean**；`D:\work\base\gopath\bin\gofumpt.exe -l` 对新件＝空（格式净）；新件非注释行 U+2600-27BF/2B00-2BFF/FE0F 命中＝0（⛔/⇄ 仅在注释，仪器豁免向）。

## 零回归 — **PASS**

- `go test ./internal/agent/approval -count=1 -timeout 30m`（GOTMPDIR/TMPDIR 起手已空＝避票 311 shape B GUARD，A1027 在册）＝**ok 0.822s**；`-v` 全量＝top-level **PASS 86／FAIL 0／SKIP 1**（上述既有一枚）。件内日志＝`D:\tmp\wisp\20261102\340v1\gotest-fullpkg-head.log`／`gotest-fullpkg-verbose.log`。
- 窄面成对控制：pristine `-run` 新用例＝ok 0.213s（`pristine-narrow.log`）。

## 纪律格 — **PASS**

- `git show --stat 412a58356`＝单文件 `internal/agent/approval/corr_clash_e2e_340_test.go` +556，与题面"新文件 556 行"一致。
- `gate.go` 未被 340-r1 改（`--name-only` 只列新件）；337-r1 的 `risk.Logf = o.Logf` 仍在 `gate.go:181`。
- 本腿写面＝仅本判语件；⛔ 碰 `cmd/wisp`／`internal/config`（345r1 在写）／thr／golden／frontend／design／票面／台账；⛔ push。

## 反向判据（本腿独立复现，⛔ 共享工作树零字节——overlay 临时根 `D:\tmp\wisp\20261102\340v1\`）

| 突变 | 落点 | 红句（逐字） | 结果 |
|---|---|---|---|
| **M1 不叠**：queue.go:171 → `corr = d.CorrelationID` | `corr_clash_e2e_340_test.go:372`（断言 2b） | `the two issued keys are NOT distinct: both cards live under "35b6fbcd-1d3d-4e80-86cb-d09c341bfbc2#call_e340" — Queue.push no longer rewrites the clash, so <corr>#<seq> is not produced` | **红**（exit=1，`mut-m1.log`） |
| **M2 换分隔符**：`%s#%d` → `%s~%d` | `corr_clash_e2e_340_test.go:382`（断言 2c） | `the rewritten key "b35b6494-8144-4d3a-ac23-7ee1f07c1fce#call_e340~2" does not carry the stacked form b35b6494-8144-4d3a-ac23-7ee1f07c1fce#call_e340#<seq>: the queue's rewrite format moved` | **红**（exit=1，`mut-m2.log`） |

两枚红落点为**不同断言行**（:372 ⇄ :382）＝票面 AC#3"成对⛔ 同一枚断言行"独立坐实，与 commit 面自报"M1/M2 各有一行不同断言取红"一致。

**观察（不拦结票，具名入账）**：`.scratch/wisp/probes/340/r1/**` 交付目录不存在，r1 的两枚红句未逐字落腿侧件（仅 commit message 引用）；本腿已独立复现两枚红句逐字写入本件，该格凭据自此不缺。另：census 记 `queue.go` 行号（:168/:301/:616）在 HEAD 漂移为 :171/:312/:627＝读数腐烂非判据差。

## 逐格总表

| 格 | 判定 |
|---|---|
| 起手三核 | **PASS** |
| AC#0 现量（callCorr/push 现量·`#数字`字面 0·真链 0→交付为全仓第一链） | **PASS** |
| AC#1 真端到端（真跑拼接非 mock、能证 `<corr>#<seq>`、名册点亮原任务行） | **PASS** |
| AC#2 守卫/Skip（0 删改、零裸 Skip） | **PASS** |
| 零回归（FAIL=0） | **PASS** |
| 纪律（stat 一致、gate.go 零改、risk.Logf 在册） | **PASS** |
| 反向判据（M1/M2 各红于不同断言行） | **PASS** |

## 结论

**同意结 340**：题面三格（AC#0/1/2）＋零回归＋纪律全部 PASS，反向判据由本腿以 overlay 独立复现两枚红句；340-r1 交付即票 328 必答①交回的那半格（叠改形真产用例从 0 → 1，且自带去牙成对）。翻框与结案由编排者按 AC#6 执行（本腿⛔ 翻框）。
