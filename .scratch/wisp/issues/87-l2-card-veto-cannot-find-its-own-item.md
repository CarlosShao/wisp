# 87 — L2 卡片上的"提前拒绝"找不到条目 ⇒ 人想拒也拒不掉，只能干等满 300 秒（票 84 的副产品，**不是安全洞**）

**Status:** open
**Type:** 可用性缺陷（安全侧已经是 fail-closed，坏的是"人无法提前结束"）
**Blocks:** nothing · **Blocked by:** nothing（`internal/agent/approval/` 与 `internal/panel/` 此刻无人写）
**Packages:** `internal/agent/approval/`（`Gate.Veto` 与它的键匹配）、`internal/panel/`（卡片回传路径）、
必要时 `cmd/wisp/run.go` 的接线。**禁改**：`docs/PLAN.md`、`docs/specs/*.md`、`internal/risk/**`、
`rules_gateway.go`、`tools/d22scan/**`、`allowlist.txt`、`frontend/**`（票 77 在写）。

## 事实来源（票 84 的实测，2026-09-21，两侧读数）

票 84 把"审批等不到答复会不会永挂"这条前提**证伪了**：等待**有上界**——
`Windows 5m0.0005138s` / `Linux 5m0.017199699s`，就是 C18 的 **300s 一律判拒绝**（`PLAN.md:1368`、`:3204`、
`SPEC-06:104-106/151`），而 `PLAN.md:3143` **明文驳回过**"把 L2 设成无限等待"的提案（理由正是永挂并持有 C20 路径锁）。
它同时量到 **8 条"无对应待审批项"的答复路由**在空门上是 **0s / 700ns–1.9µs** 返回 `ErrUnknownCorrelation`
（不 panic、不重试）⇒ **答复侧是诚实的**。

**剩下这条是真问题**：卡片**已经显示**给用户、但 `Gate.Veto`（或等价按键拒绝的路径）因为键形不上而找不到那条待审批项时，
**人想提前拒绝也做不到**，只能看着它满 300 秒。票 75 观察到的 `4m45s` 卡等就是这个形状的放大版
（不是永挂，是**被迫等满上界**）。

⚠ **不要把这张票读成"修一个死锁"**。安全语义上现状是可接受的（超时会拒绝、宿主不可达 fail-closed）。
这张票要修的是**人的控制权**：看到卡片的人应当能立刻说"不"。

## 先查后修（AC#1 是只读）

- [ ] **AC#1** 把"卡片显示 ↔ 门里的条目"这条键链画出来：卡片上带的关联标识是什么形状、
      谁生成、`Veto`/`Allow`/`Reject` 各自按什么查（**逐处给 file:line**）。
      然后**列出所有能让"卡片在显示、门里查不到"成立的路径**（例如队列被 LRU 挤掉、批次尺寸裁剪、
      重启后重建卡片、面板与宿主两套键）。**查不到实例就如实写"没找到"** ⇒ 本票降级为文档说明并关闭。
- [ ] **AC#2** 若 AC#1 找到实例：**加一条双向都可判的用例**——"卡片存在 + 门里查无此项"时，
      点击拒绝必须 **(i) 立刻结束该次的等待**（不许再等满 300s）并且 **(ii) 仍然按拒绝处理**
      （找不到条目只能让"拒绝"更容易生效，**绝不能**因为查不到就放行或忽略）。
      ⚠ 这条 **(ii) 是本票的安全底线**：任何"查不到 ⇒ 当作已批准/当作无需处理"的修法一律不过。
- [ ] **AC#3** 变异双向：(i) 把"立刻结束"退回"等满上界" ⇒ 用例红；
      (ii) 把失败侧从"拒绝"改成"放行" ⇒ **必须有既有用例红**（R7/C18 的 fail-closed 家族）。
      锚点=承载行为的那一行，同链 grep 自证落地，还原后 `git diff --quiet` 证干净；编译失败不算变异。
- [ ] **AC#4** 不动 C18 的 300s：**一个数字都不改**（不许"为了少等把超时调小"）。
      `PLAN.md:1368` 的"超时前 30s 醒目提示"若未实现，只登记，不在本票补。
- [ ] **AC#5** 门禁（只跑自己碰的包）：`gofmt -l` 空、`gofumpt -l <files>` 空、
      `go vet <pkgs>` rc=0、`GOOS=linux go vet <pkgs>` **按包作用域**跑 rc=0（⚠ 别用 `GOOS=linux go vet ./...`，
      那条在 Windows 主机上因 CGO=0 排除 sherpa 而永远 rc=1，A54③）、`go test -count=2 <pkgs>` rc=0 且逐条点名 SKIP/FAIL。

## Rules（本仓固定）

15 次工具调用内交第一枚 checkpoint commit；每次 commit 同步 Status + 勾框 + 末行 `next=`；
**每完成一组就把结论追加进票面 log**（别攒着，本仓的代理死过在轮数上限上）；
`git commit -q -F - -- <显式路径>` + 带引号 heredoc；禁 `git add -A`/`.`；commit 前核对 `git diff --cached --name-only`；
禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`（共树 A34）；**不 push**；不在仓内建 worktree（A38④）；
docker 挂 `git archive` 快照且**目录名带你自己的会话后缀**（A59⑤：两会话共用 `/tmp/wisp75` 撞过一次），
Git Bash 里挂容器路径要 `MSYS_NO_PATHCONV=1`；四种假绿逐跑点名；票面 append-only，**要改的那行先读再替换**。

## Progress log（append-only）

（空）
