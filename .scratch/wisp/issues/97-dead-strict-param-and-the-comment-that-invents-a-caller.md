# 97 — `strict=true` 是**死参 + 一句描述不存在的调用方的注释**，而"别名买不到批准"这件事**一条用例都没钉**（票 87 验收带回）

**Status:** in-progress（2026-09-21 15:4x 编排者建；**排队**：`internal/agent/` 此刻有 `agent-ticket90` 在写，同包冲突）
              2026-09-21 19:36 `agent-ticket97` 认领并开始交件（票 90 已 `accepted-done`、票 87 已结 ⇒ 队列条件解除，本目录现归本票）
**Type:** 安全方向性的**表达**缺陷（今天方向是对的，但它靠的是习惯而不是机器）+ 一处说谎的注释
**Blocks:** nothing · **Blocked by:** 票 **90** 交件（同目录 `internal/agent/approval/`，别同时写）
**Packages:** `internal/agent/approval/`（`resolveLocked` 与那张别名索引、`Queue.reject`）。
              **禁改**：`docs/PLAN.md`、`docs/specs/*.md`、`internal/risk/assessor.go`、
              `internal/risk/pathresolver*.go`、`rules_gateway.go`、`tools/d22scan/**` 与 `allowlist.txt`。

## 来源（验收代理的实测，不是猜测）

`docs/evidence/s1/87-adversarial-acceptance.md` §11 的 R-1 / R-2，两条都在票 87 的**判据之外**（不改它的 PASS）：

- **R-1**：`resolveLocked(name, strict bool)` 的 **`strict=true` 零调用点**，而注释写着
  "an allow passes true" ——**描述的是一个不存在的调用方**。
  ⇒ 于是"别名只服务拒绝方向"这件事今天靠的是**非导出 + 单一读者**，**不是类型**，
  并且代理 grep 的读数是**零条断言**钉住"别名买不到批准"。
  **这就是票 83 那句"会撒谎的配置键"的注释版**：注释承诺了一条代码里不存在的防线。
- **R-2**：宽松解析写在 `Queue.reject()` 里 ⇒ 它的**影响面是全部 5 条拒绝路线**
  （native / panel / `DecideFrom*` / `Veto`），**比票 87 AC#2 叙述的"Veto 查不到时"宽**。
  方向仍然只可能产生拒绝 ⇒ **不触底线**，但**账目把改动面写窄了**。

## 要做什么

1. 把那个 bool 换成**两个具名函数**（例：`lookupForAllow` / `lookupForRefusal`），让"方向"出现在**签名上**
   而不是注释里——这样"allow 侧拿到别名表"会变成**编译不过**，而不是"记得别传 true"。
2. 补一条正向钉子：`TestAnAliasCanNeverBuyAnAllow`（验收代理 §6.2 攻击 1 的形状可直接搬，它已经实测过被挡住）。
3. 把那**一句说谎的注释**改成与代码一致的描述（保留"为什么只有拒绝方向"的推理，删掉"an allow passes true"这种
   **不存在的路径**）。⚠ 注释与真实防线不一致本身就是缺陷（票 73 的 A42 就是这一族：**兜底来自哪一层，注释说错了方向**）。
4. **R-2 就地如实登记**：在票面/commit 正文写清"影响面 = 全部 5 条拒绝路线"，并给一条**逐路线**的用例矩阵
   （5 条路线各断一次"查不到条目 ⇒ 仍是拒绝，绝不是放行"）。

## AC（1:1，裁决表 `docs/evidence/s1/97-*.md` 由验收方出）

- [ ] **AC#1** `grep -rn "strict" internal/agent/approval/` 的读数里**不再有零调用点的方向参数**；
      或你论证它为何必须存在并给它一条用例。**二选一，写下来。**
- [ ] **AC#2** `TestAnAliasCanNeverBuyAnAllow` 落地，且**变异**：把 allow 侧改成也查别名表 ⇒ 该用例必须红
      （锚点=承载行为那一行，同链 grep 证落地，还原后 `git diff --quiet` 证干净，**编译失败不算变异**）。
- [ ] **AC#3** 5 条拒绝路线的矩阵用例逐条点名（每条一个子测试名），并给"其中一条被改成放行 ⇒ 有用例红"的反向变异。
- [ ] **AC#4** 注释与代码一致：贴出你改前后的注释原文，并说明**新注释的每句话在代码里能找到对应物**。
- [ ] **AC#5** 门禁（按包）：`gofmt -l`/`gofumpt -l` 空、`go vet ./internal/agent/approval/` rc=0、
      `go test -count=2 ./internal/agent/approval/` rc=0 且逐条点名 SKIP/FAIL
      （⚠ 票 87 的验收读数里那 **2 条 SKIP 是 `TestDefaultDeadlineWallClockMeasurement`**，
      既有、不属本票，但**必须照点名**；`go test ./cmd/wisp/` 本机是加载期 `0xc0000135` 既有红，见**票 98**，不要追）。
      **收尾前必跑 `sh scripts/d22scan.sh`**（A64②：按包门禁看不见全仓 ban）。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；提交前 `git diff --cached --name-only`；
禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34）；**不 push**；
不在仓内建 worktree（A38④，快照目录带会话后缀）；票面 append-only（改行前先读；标题前插段落要重抄标题，
`git diff --numstat` 删除列必须 0）；四种假绿逐条点名；数字不达标写 FAIL 附数字。
15 次工具调用内交回第一枚 checkpoint；每次提交同步 Status + 勾框 + `next=`；接近轮数上限主动收尾留断点。
⚠ **R-3（入键撞车会拒到别人的卡）不在本票**——它的根在**接线时的键对账**，已登记给**票 37**（见其票面），
本票只负责让"方向"变成签名上看得见的东西。

## Progress log（append-only）

- 2026-09-21 15:4x（编排者）：建票。来源是 `acceptor-ticket87` 的残留清单（它判票 87 五框全 PASS，
  但 R-1/R-2 是"判据之外、须被下一张票接住"的东西 ⇒ 按规矩落成票，不留口头）。
  ⚠ 顺带更正它的一处引用：它写"用户可见效果等票 37/41"——**41 是误引**（票 41 是 KWS 唤醒词 + 语音 veto 词，
  同名不同物），正确归口是 **37 + 35**。我在 87 的票面上把它改成 37/35 并注明误引（本仓已抓到两起引用腐坏：A30、A51⑥）。
  next= 等票 90 交件让出 `internal/agent/` ⇒ 派单（R-1 成本很低，别让它沉）。

- 2026-09-21 19:36（`agent-ticket97`，第 1 枚 commit = 钉子先行）：**先交用例，代码一行没动**。
  新增 `internal/agent/approval/ticket97_alias_direction_test.go`（外部测试包 `approval_test`，
  复用票 84/87 的 `newGate`/`runApproval`/`mustAnswerWithin`/`vetoBudget` 仪表，不自造一套）：
  `TestAnAliasCanNeverBuyAnAllow`（AC#2）+ `TestEveryRefusalRouteOnAnUnknownEntryStillRefuses`
  （AC#3 的 5 路线矩阵：`/veto` `/native_reject` `/panel_reject` `/decide_from_native` `/decide_from_panel`）。

  ⚠ **本票性质所致的红法**（票面"工作方式"已预告）：这两条用例在未修改的树上**天生就绿**
  （`go test -count=1 -run 'TestAnAlias|TestEvery' ./internal/agent/approval/` → `ok ... 0.056s`，
  基线快照 `git archive b878b30` 里同样 `ok 0.056s`）⇒ **不许把"修前红"这格当仪式勾掉**，
  下面用三发变异逐条证明它有牙。变异全部在仓外快照 `/tmp/wisp97-agent-ticket97`（A38④，仓内零 worktree），
  每发同链 `grep -n` 打印被改后整行、`go build ./internal/agent/approval/` rc=0 后才算数（编译失败不算变异）：

  - **M1（AC#2 的牙）**：`queue.go:328`
    `it, ok := q.resolveLocked(corr, false), q.resolveLocked(corr, false) != nil // MUTATION-agent-ticket97-M1`
    ——即票面点名的"allow 侧偷偷读别名表"。build rc=0。红名 + 断言原文：
    ```
    --- FAIL: TestAnAliasCanNeverBuyAnAllow (0.00s)
        ticket97_alias_direction_test.go:85: SECURITY: Native().Allow("task-97-host-key", 卡片自己的活 grant) 返回 nil：别名买到了批准
    ```
    同一快照未变异基线：`--- PASS`（0.056s）。⇒ 这条钉子测的是真事：**别名 + 活 grant 一旦能被 allow 侧解析，批准就会发生**。
  - **M2（AC#3 的牙，共享漏斗）**：`queue.go:376` 让 `reject()` 在查不到条目时**放行队头**（fail-open）。build rc=0。
    **5 个子测试全红**，每条都是同一句：
    `ticket97_alias_direction_test.go:192: SECURITY: 路线 <name> 对查不到的条目返回 nil：拒绝路线把一次没有发生的拒绝当成了成功`
    （`--- FAIL` × 5：veto / native_reject / panel_reject / decide_from_native / decide_from_panel）。
    ⚠ 这正是 **R-2 的实测形状**：一颗埋在共用漏斗里的种子，**5 条路线同时开花** ⇒ 影响面是 5 不是 1，
    票 87 AC#2 叙述的"Veto 查不到时"确实写窄了，本票据此登记。
  - **M3（AC#3 逐路线的牙）**：`gate.go:584` 只把 `nativeAPI.Reject` 改成走放行漏斗
    （`return n.g.q.allow(corr, reason)`）。build rc=0。**只有 `/native_reject` 红**，另外 4 条 `--- PASS`：
    ```
    ticket97_alias_direction_test.go:210: 路线 native_reject 用卡片自己的别名点名: approval: correlation_id 无对应待审批项, want nil（R-2：宽松解析对全部 5 条路线生效）
    ticket97_alias_direction_test.go:212: 路线 native_reject 的别名拒绝: STILL BLOCKED after 2s - 闸门没有上界
    ```
    ⇒ 矩阵是**逐路线**可判的，不是一句笼统断言的复读。
  - **还原证**：快照 `diff -q` 两个被改文件与 `/tmp/*.orig97` 一致（`snapshot restored clean`）；
    仓库侧 `git status --porcelain internal/agent/approval/` 只剩我这枚未跟踪用例、
    `git diff --quiet -- internal/agent/approval/queue.go internal/agent/approval/gate.go` rc=0。
  用例形状的两处防"零信息"前题（红即前题破，不算通过）：卡片必须**带活 grant**、精确键**必须 ≠ 别名**；
  拒绝方向在同一用例里也断一次（同一别名点名的确实是**拒绝**而非"什么也发生不了"）。

  next= 第 2 枚 commit：`resolveLocked(name, strict bool)` ⇒ 两个具名函数（`lookupForAllowLocked` /
  `lookupForRefusalLocked`，方向上签名、无 bool），把 `queue.go:220-229` 那句 "an allow passes true"
  改成与代码一致（保留"只有拒绝方向"的推理），再跑按包门禁 + `sh scripts/d22scan.sh`，Status 翻 `ready-for-review`。
