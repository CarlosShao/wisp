# 87 — L2 卡片上的"提前拒绝"找不到条目 ⇒ 人想拒也拒不掉，只能干等满 300 秒（票 84 的副产品，**不是安全洞**）

**Status:** ready-for-review（AC#1-AC#5 五框全勾；码与用例在 `dfe9d9f`→`1068eb9`，两侧数字在本票 log 末段）
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

- [x] **AC#1** 把"卡片显示 ↔ 门里的条目"这条键链画出来：卡片上带的关联标识是什么形状、
      谁生成、`Veto`/`Allow`/`Reject` 各自按什么查（**逐处给 file:line**）。
      然后**列出所有能让"卡片在显示、门里查不到"成立的路径**（例如队列被 LRU 挤掉、批次尺寸裁剪、
      重启后重建卡片、面板与宿主两套键）。**查不到实例就如实写"没找到"** ⇒ 本票降级为文档说明并关闭。
- [x] **AC#2** 若 AC#1 找到实例：**加一条双向都可判的用例**——"卡片存在 + 门里查无此项"时，
      点击拒绝必须 **(i) 立刻结束该次的等待**（不许再等满 300s）并且 **(ii) 仍然按拒绝处理**
      （找不到条目只能让"拒绝"更容易生效，**绝不能**因为查不到就放行或忽略）。
      ⚠ 这条 **(ii) 是本票的安全底线**：任何"查不到 ⇒ 当作已批准/当作无需处理"的修法一律不过。
- [x] **AC#3** 变异双向：(i) 把"立刻结束"退回"等满上界" ⇒ 用例红；
      (ii) 把失败侧从"拒绝"改成"放行" ⇒ **必须有既有用例红**（R7/C18 的 fail-closed 家族）。
      锚点=承载行为的那一行，同链 grep 自证落地，还原后 `git diff --quiet` 证干净；编译失败不算变异。
      （两侧数字与红名见 log 末段：(i) 2 红含 `STILL BLOCKED after 2s`；(ii) 4 红其中
      `TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis` + `TestReplayRedisplaysUnderAFreshGrant` 是既有用例。）
- [x] **AC#4** 不动 C18 的 300s：**一个数字都不改**（不许"为了少等把超时调小"）。
      `PLAN.md:1368` 的"超时前 30s 醒目提示"若未实现，只登记，不在本票补。
      （`queue.go:92-94` 与 `gate.go:456-460/487-499` 均未改；30s 提示**已实现**，见 log「AC#4 结案」。）
- [x] **AC#5** 门禁（只跑自己碰的包）：`gofmt -l` 空、`gofumpt -l <files>` 空、
      `go vet <pkgs>` rc=0、`GOOS=linux go vet <pkgs>` **按包作用域**跑 rc=0（⚠ 别用 `GOOS=linux go vet ./...`，
      那条在 Windows 主机上因 CGO=0 排除 sherpa 而永远 rc=1，A54③）、`go test -count=2 <pkgs>` rc=0 且逐条点名 SKIP/FAIL。
      （数字见 log；⚠ `go test ./cmd/wisp/` 在本机 rc=1 = 加载期 `0xc0000135`，
      纯净树 `git archive 63ef895` 同读数 ⇒ 先于本票存在，已登记给编排者。）

## Rules（本仓固定）

15 次工具调用内交第一枚 checkpoint commit；每次 commit 同步 Status + 勾框 + 末行 `next=`；
**每完成一组就把结论追加进票面 log**（别攒着，本仓的代理死过在轮数上限上）；
`git commit -q -F - -- <显式路径>` + 带引号 heredoc；禁 `git add -A`/`.`；commit 前核对 `git diff --cached --name-only`；
禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`（共树 A34）；**不 push**；不在仓内建 worktree（A38④）；
docker 挂 `git archive` 快照且**目录名带你自己的会话后缀**（A59⑤：两会话共用 `/tmp/wisp75` 撞过一次），
Git Bash 里挂容器路径要 `MSYS_NO_PATHCONV=1`；四种假绿逐跑点名；票面 append-only，**要改的那行先读再替换**。

## Progress log（append-only）

### 2026-09-21 · 票 87 代理 · AC#1 只读结论：**找到实例**（不是"没找到"）

**键链（谁生成 → 什么形状 → 各答复按什么查）**

1. 标识生成（三层，逐处）：
   - 桥：`internal/tools/bridge.go:219-220` `req.CorrelationID == ""` 就填 `req.TaskID`；`bridge.go:279` 放进
     `Decision.CorrelationID`。文本循环：`internal/agent/loop.go:644` `CorrelationID: taskID` ⇒ 生产形状是
     **corr == taskID**（字符串相同，不是同一概念相同）。
   - 队列再签发：`internal/agent/approval/queue.go:126-152 push()` —— 键 `corr := d.CorrelationID`（133），
     空则 `"approval-"+seq`（135），撞车则 `fmt.Sprintf("%s#%d", d.CorrelationID, q.seq)`（138）；
     存 `it.Corr`（141）与 `q.byID[corr] = it`（148）。形状：任意调用方字符串，或 `approval-<n>`，或 `<corr>#<seq>`。
   - 卡片带的是**队列的键**：`gate.go:446 corr := it.Corr` → `gate.go:452 d.CorrelationID = corr` →
     `gate.go:464 promptFor(d, corr, …)` → `Prompt.CorrelationID = corr`（`gate.go:528`）。
     面板视图同名：`PanelItem.CorrelationID = it.Corr`（`queue.go:325`）。
   - **同一次调用里存在第二套键**：D31 记账用 `incoming := orDefaultText(d.CorrelationID, d.TaskID)`
     （`gate.go:451`，注释 447-450 自认"两个不是同一个字符串一旦队列重签发过"），`markStarted(incoming, …)`
     （`gate.go:482`）；桥的取消总线同样用进入时的键（`bridge.go:397/399/463/519`）。
2. 各答复按什么查：
   - `Native().Allow` → `gate.go:561` → `queue.go:231 q.byID[corr]` + grant 绑定摘要（`queue.go:240`、`approval.go:253-261`）。
   - `Native().Reject` → `gate.go:564` → `queue.go:268 byID`；`Panel().Reject` → `gate.go:568` → 同上；
     `Panel().View/Head` → `queue.go:316/338`。`DecideFromNative/DecideFromPanel` → `gate.go:595/593/614` → 同上。
   - `Gate.Veto`（按键/悬浮球/Esc/KWS 那条）→ **只查 `g.windows[v.CorrelationID]`（`gate.go:372`）和
     `g.running`（`gate.go:390`）**；L1 窗口表的键是 `orDefaultText(d.CorrelationID, d.TaskID)`
     （`gate.go:247-249`）。**`Veto` 全程没有一行去查 `g.q`。**

**能让"卡片在显示、门里查无此项"成立的实例（逐条判）**

- **实例 1（成立，主因）**：L2 路线 `PendingApproval` **从不 `openWindow`**（`gate.go:432-513` 里没有任何
  `openWindow` 调用），所以卡片正显示时 `Veto` 即使在**正确的 `it.Corr`** 上也走 `gate.go:407`
  返回 `ErrUnknownCorrelation`。并且卡片**自己把这个通道列为可用**：`promptFor` 对 L1/L2 一视同仁地塞
  `Channels: g.channels.Statuses()`（`gate.go:541`），`cmd/wisp/run.go:525-527` 原样打印「取消方式：…（已加载）」。
  ⇒ 人按了，票没被听见，只能等满 300s（`gate.go:501-504`）。这就是票面的形状。
- **实例 2（成立）**：R7 大批次在 L1 路线上被升级改走 L2（`gate.go:228-244`）⇒ 同一个 corr 在 L1 会命中窗口表、
  升级后反而查不到，属于实例 1 的一个入口。
- **实例 3（成立，条件性）**：队列重签发后（`queue.go:135/138`）卡片键 `it.Corr` ≠ 宿主记账键 `incoming`；
  以 task id / 进入时 corr 回投的宿主 → `byID` 查不到。当前 `bridge.go:219-220` 保证非空，所以只有"corr 撞车"
  才触发（并发多卡、或 replay 名字撞车，`queue.go:368-370`）。
- **不成立/不存在**：(a) 队列满 → `push` 之前就 fail-closed 拒绝（`queue.go:129-131`），卡片根本没显示；
  pending 集合没有 LRU 挤掉机制（只有 history LRU `queue.go:175-180`，那时卡片已消失）；
  (b) 批次裁剪 `batch.go` 只改 L1 卡片**内容**（前 5 条），不改键；
  (c) "重启后重建卡片" —— 全仓找不到 pending 审批的持久化/重建代码，**这条路径不存在**（如实写）；
  (d) 面板/宿主两套键 —— 只要面板从 `PanelAPI` 取视图就是同一个 `it.Corr`，**没找到**第二套键；
  ⚠ 唯一没核完的：`internal/panel/approval.go:40 correlationId` 是**票 77 刚刚落的**（本次会话内新出现），
  它的 correlationId 从哪个结构取我**没读**，因为 `internal/panel/` 归它，**登记给编排者**。

**修法取向（AC#2）**：`Veto` 在 windows/running 都查不到时，去门里按 reject 方向结算该 pending 项
（立刻结束等待 + 仍按拒绝处理）；并给队列加**只在拒绝方向生效、且仅在无歧义时**的别名解析
（卡片键 / 进入时 corr / taskID）。允许侧一律保持只认精确键 + grant 绑定 —— 别名宽松只给"更容易拒绝"。
`300`/`30` 两个数字不改（AC#4；`queue.go:92-94`，且 30s 提示确已在 `gate.go:458-460/487-499` 实现，无需登记缺）。

（本条 checkpoint 时勾了 AC#1；AC#2/AC#3/AC#5 的数字在下面的段落里逐组补。）

### 2026-09-21 · 票 87 代理 · AC#2 落码（第二枚 checkpoint 的内容）

改动只在 `internal/agent/approval/`（`internal/panel/` 归票 77，未碰）：

- `queue.go`：`qitem` 新增 `names`（同一张卡的其它合法称呼：进入时的 corr、taskID）；`Queue` 新增
  **只服务拒绝方向**的 `alias` 索引 + `resolveLocked(name, strict)`：
  精确键永远优先，`strict=true`（允许侧）拿不到别名，`strict=false`（拒绝侧）才读别名索引，
  **一个别名指向多张活卡时不许猜**（返回"查无此项"）。`reject()` 换成宽松解析，
  `allow()` **一行未动**（仍是精确键 + grant 绑定摘要）。
- `gate.go`：`Veto` 在 windows/running 两级都查不到之后，把这一票交给 `q.reject(...)`（原生/面板"拒绝"
  按钮同一个漏斗），于是 AC#1 的实例 1/2/3 都不再需要等满上界。方向自查：这条分支只能产生拒绝。
- 新用例 `ticket87_veto_l2_test.go`（4 条）：精确键否决 / 宿主自有键（taskID）否决 / 两名歧义不许猜 /
  无关 id 的否决不得碰到别人的卡（该卡随后仍可用自己的 grant 批准，反证否决没变成放行）。
  门限 `ApprovalTimeout: 30s`，断言 2s 内结束 ⇒ (i) 的变异会真等 30s 再红，**不写 300s 常跑用例**。

AC#4 自查：`300`/`30` 两个数字一个没改（`queue.go:92-94` 原样；30s 提示在 `gate.go:458-460/487-499` 已实现，
无需登记缺，console UI 会打印 `[warning]`）。

已跑门禁：`gofmt -l internal/agent/approval/` 空、`gofumpt -l` 我这三个文件空、
`go vet ./internal/agent/approval/` rc=0、`GOOS=linux go vet ./internal/agent/approval/` rc=0、
`go test -count=2 ./internal/agent/approval/` rc=0（`=== RUN` 84 = count=1 的 42 ×2、FAIL 0、
SKIP 2 = 票 84 的 `TestDefaultDeadlineWallClockMeasurement` ×2，默认 `WISP_84_MEASURE` 未设 ⇒ 有意慢跳过）。
⚠ 全仓 `gofumpt -l .` 现在唯一红的是 `internal/risk/syncdirs_test.go`（票 82 在飞的测试文件，**不是我写的**，我没碰）。

next= 跑 AC#3 两侧变异（(i) 把 `Veto` 的 L2 分支退回 `return ErrUnknownCorrelation` ⇒ 新用例红；
(ii) 把 `Queue.reject` 的 `tools.AnswerReject` 改成 `tools.AnswerAllow` ⇒ **既有** fail-closed 用例红），
把两侧红名抄回票面后收尾。`internal/panel/approval.go` 的 correlationId 来源仍登记给编排者。

### 2026-09-21 · 票 87 代理 · AC#1 的"待核"补上 + AC#3 变异两侧 + AC#5 数字（收尾）

**面板与宿主两套键（这次是读到的，不是猜的）**：`internal/panel/approval.go:28/40/78` 的
`ApprovalCardView.CorrelationID` 取自 `ApprovalSubject.CorrelationID`，注释写明"the agent loop owns the
correlationId (C17)" ⇒ 面板卡片带的是**进入时的 C17 键**；而队列自己的键是 `it.Corr`（`queue.go:133-139`，
只有空/撞车时才与之不同）。今天两者相同只因为 `internal/agent/loop.go:644` 把 corr 写成 taskID，
且队列没重签发。同时 `internal/panel/` 与 `cmd/wisp/panel_assets.go` 里 **零** 处引用 `Gate.Panel()/Native()`
或 `Reject/Veto` ⇒ 面板回投递这条路还不存在（票 37），所以这条现在是**条件性实例**（一旦接线就成真）；
本次的别名索引正是提前把它接住的形状（进入时 corr 与 taskID 都能命中拒绝方向）。
`cmd/wisp/panel_assets.go:45` 那张 `correlationId: "panel-assets-l2"` 是票 77 的 embed 证明面（写死的样例，
门里根本没有这一项），**render-only**，不是活的审批。

**AC#3 变异（逐条点名，锚点=承载行为那一行，同链 grep 自证落地，还原后 `git diff --quiet -- internal/agent/approval/` rc=0）**

- (i-a) `gate.go:426` 锚点 `return g.q.reject(v.CorrelationID,` → 加 `-mut87i` 后缀（等价于"这条分支不存在"）：
  rc=1，`=== RUN` 4 条中 **2 红**：
  `FAIL TestAVetoAgainstAnOnScreenL2CardEndsTheWaitNowAndStillRefuses`（`Veto on the displayed card:
  approval: correlation_id 无对应待审批项`）+ `FAIL TestAVetoNamedByTheHostsOwnKeyStillRefusesThatCard`；
  歧义不许猜/无关 id 不碰别人的卡 两条 PASS（它们本来就不依赖这条分支）。
- (i-b) 同一锚点改成 `return nil`（"听见了但其实什么都没做"）：rc=1，2 红，红名是**等待上界断言**：
  `ticket87_veto_l2_test.go:91: veto against an on-screen L2 card: STILL BLOCKED after 2s
  (start=13:34:56, now=13:34:58) - 闸门没有上界` 与 `:147` 同形状 ⇒ 证明 (i) 那半边真的在测"立刻结束"，
  而不是只测错误码。用例 armed 的是 30s 门限（**不是 300s 常跑用例**）。
- (ii) `queue.go:383` 锚点 `answer{a: tools.AnswerReject, why: why}`（`Queue.reject` 的投递，唯一一处写
  "人说了不"的方向）→ 改 `tools.AnswerAllow`：rc=1，全仓 `=== RUN` 42 条里 **4 红**，其中**两条是既有用例**：
  `FAIL TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis`（`queue_test.go:143: answer="allow"，期望 reject：
  面板只能拒绝`）、`FAIL TestReplayRedisplaysUnderAFreshGrant`（`queue_test.go:280: answer="allow"，期望 reject`）
  —— R7/C18 fail-closed 家族确实坐在这一行上；另两条是本票新用例。
- 三次变异全部还原：`cp` 回备份 + `git diff --quiet -- internal/agent/approval/` rc=0（工作树里别人的
  `frontend/`、`internal/risk/` 改动我没碰也没还原）。

**AC#5 门禁（只跑我碰的包 + 两个下游）**

- `gofmt -l internal/agent/approval/` 空；`gofumpt -l internal/agent/approval/` 空；
  commit 前重跑全仓 `gofumpt -l . tools/d22scan tools/mockllm` **空**（上一次它红在 `internal/risk/syncdirs_test.go`，
  那是票 82 在写的文件，现已不见 ⇒ 他们自己收干净了）。
- `go vet ./internal/agent/approval/` rc=0；`GOOS=linux go vet ./internal/agent/approval/ ./internal/tools/` rc=0
  （**按包作用域**，没用 `./...`，A54③）。
- `go test -count=2 ./internal/agent/approval/` rc=0：`=== RUN` 84 = count=1 的 42 ×2（核对过倍数）、
  `--- FAIL` 0、`--- SKIP` 2 = `TestDefaultDeadlineWallClockMeasurement` ×2（票 84 的 300s 墙钟计量，
  默认无 `WISP_84_MEASURE` ⇒ 有意慢跳过，不是假绿）、`no tests to run` 未出现。
- 回归：`go test -count=2 ./internal/tools/` rc=0（30.0s，`Veto` 的老形状调用方在那儿）。
  ⚠ `go test ./cmd/wisp/` rc=1 `exit status 0xc0000135`（STATUS_DLL_NOT_FOUND）：**不是本票的回归** ——
  证据：`-run '^$'`（不执行任何用例）同样 0xc0000135，且 `git archive 63ef895`（我落码之前的提交）解到
  `/tmp/wisp87s75ctl` 纯净树跑 `./cmd/wisp/` 也是同一读数 ⇒ 本机的加载期问题（cgo/DLL 不在 PATH），登记给编排者。

**AC#4 结案**：`300`/`30` 一个数字未改（`queue.go:92-94` 与 `Options.ApprovalTimeout/WarningLead` 原样；
新用例只用 30s 的门限做断言）。"超时前 30s 醒目提示"**已实现**：`gate.go:458-460` 起 warn 定时器、
`gate.go:487-499` 发 `EventWarning` 并写明"审批将在 N 秒后自动拒绝"，console 侧 `cmd/wisp/run.go:538-541`
会打印 `[warning] …` ⇒ 本票无需登记缺失，也无需补。

四框现全勾。next= 交回编排者：`cmd/wisp` 的 0xc0000135（加载期，先于本票存在）要不要单开票；
`internal/panel/approval.go` 的 C17 键与队列键在票 37 接线时的对账归谁（我这边已备好拒绝方向的别名解析）。

