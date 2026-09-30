# 票 246 — 对抗验收腿 `246-v1` 判定表：最短链「注入 gate → 真卡片 → `Esc` 否决一次并归还 → 退出序列不撒谎」

> ⚠ **这份表的身份**：出自**非实现者**（`246-v1` 对抗验收腿）。按 `AGENTS.md §1.5`／`SPEC-12 §4.3`，
> 本表是裁决判语；实现方那份 `246-...-r1.md`（291 行、脏件未提交完整）**只当"待攻主张名册"，其读数不作本腿凭据**。
> 判据出处＝`.scratch/wisp/issues/246-...-no-veto-channel-has-an-executor.md` AC#0..AC#6。
> 编排者裁定＝乙形（装配根注入）＋"钩子注册入口不算契约面"，账 `A481`。所有读数本腿现跑，落盘 `.scratch/wisp/probes/246/v1/`（只建不删）。

---

## §0 起手名册 · 锚点 · 主张名册

### §0.1 起手（闸门＝终态等于起手那一刻的名册，不是"必须为空"）

| 项 | 读数（本腿现跑） |
|---|---|
| 取数时刻 | `2026-09-30 19:36:20 +0800`（`date` 现跑） |
| 起手锚点 | `7c93ad7c`（`git rev-parse --short=8 HEAD` 自取，⛔ 未采用派单里给的任何 sha） |
| 分支 | `dev` |
| 起手名册全量 | `.scratch/wisp/probes/246/v1/roster-start.txt`（`git status --porcelain` 原文） |
| `git status --porcelain` 行数 | **182**（起手即含实现方遗留的脏件：`M cmd/wisp/…` 无、`M docs/evidence/s1/246-…-r1.md` 有、`design/**` 删除、`part*-*.txt` 等，本腿一枚不碰） |
| 桌面安静（跑真机前） | `tasklist //FI "IMAGENAME eq balldebug.exe"`＝**0**／`wisp.exe`＝**0**／`esclistener.exe`＝**0** |
| 本腿写面（白名单） | `docs/evidence/s1/246-...-v1.md` ＋ `.scratch/wisp/probes/246/v1/**`；⛔ 不改任何产码／工单框／冻结件／`frontend**``design**` |
| 五枚本票 commit | 「`75798f31` 证据件骨架 → `0c1b32cf` 乙形落地 → `2b5462f9` 改法名融合 → `17d8586e` 测试半边（编排者代提）→ `7c93ad7c` 产码两枚＋台账（编排者代提）」 |

### §0.2 主张名册（实现方自述，逐枚本腿自己攻；⛔ 不引用其表读数）

| # | 实现方主张 | 本腿攻法 |
|---|---|---|
| C-1 | AC#1 零新增包级依赖边（差集前后逐名相同） | `go list -deps` 现跑 + 逐枚 base 直依复认（§3 R-A） |
| C-2 | AC#2 一张真卡片在常驻进程挂起、状态机真进 `Confirming`、非 mock | 问生产调用者 + winlive 现跑 + 无窗口 fail-closed（§1 AC#2） |
| C-3 | AC#3 `借→否决→归还` 三形齐 + 台件已进仓 | winlive 现跑 + 台件在库复认（§1 AC#3） |
| C-4 | AC#4 钩子真注册、step3 不再 skipped、待决卡拒绝＋留审计；十步顺序未动 | 读注册链 + 默认层/ winlive 现跑 + `git diff` shutdown.go（§1 AC#4） |
| C-5 | AC#5「起手 2/2 绿、无红可归因」 | ⛔ 关键攻击点：HEAD 全量默认层是否绿（§1 AC#5 / §3） |
| C-6 | AC#6 其余三条否决通道未顺手做 | grep + 通道 loaded 断言（§1 AC#6） |
| C-7 | 代提 `7c93ad7c`：守卫 `!rb.cancelHosted` 不是装饰，`escLoad` 判据成立 | 定向突变（§2 MUT-6） |

---

## §1 逐格判语

> 姿势：判据 → 落地形状（`file:line`）→ **本腿现跑读数** → 判（成立／不成立／附条件成立／判不动）。
> 本节 §1 判定随取数逐步补满；下表为骨架格，读数落 §3。

### AC#1 — 装配根注入一枚真审批门到常驻腿、零新增包级依赖边 — **本腿初判：成立**（详见 §3 R-A）

- 落地形状（本腿现读）：`cmd/wisp/resident_windows.go:122 newResidentApproval()` → `:135 rb := startResidentBall(rt.Registry, ra.vetoByEsc)` → `:144 ra.bindBallHost(rb)` → `:156 rt.RegisterShutdownHook(proc.StepCancelTasks, ra.cancelTaskRoots)`。gate 建在装配根，递给球宿主的是函数值 `escVetoFunc`（`resident_ball_windows.go:61/:111`）。
- 差集证据（本腿自跑，非抄表）：见 §3 R-A。**内部包级边 0 枚新增**。

### AC#2 — 一张真卡片能在常驻进程里挂起来（可见性只到"日志＋状态位＋Win32"） — **本腿初判：成立（射程内）**

- 判据红线：⛔ 面板开不出来是既成事实（`approval_always.go:165` 逐字 this binary links no WebView2 host），任何"用户看得见卡片"一律判不成立。
- 落地形状（本腿现读）：`resident_approval_windows.go:205 askConfirmation` 走真 `gate.AdmitTextTask`＋`gate.PendingWindow`(L1)/`gate.PendingApproval`(L2)；`:428 Prompt` 由注入 UI 记账 `cards.Record`、`:440` L1 才 `TakeEscForCancel`、`:443 SetState(stateForCardLevel)`；`stateForCardLevel` 只产 `Confirming`/`AwaitingApproval` 两枚 D43 名，⛔ 不造第五枚。
- 非 mock 的反证（默认层，本腿 run1 日志现读）：`TestAC246CardWithNoWindowFailsClosedThroughTheRealGate` → 无窗口时真 gate 逐字 `approval: 卡片无处呈现（本进程没有悬浮球窗口），已 fail-closed 拒绝`。
- 诚实残余：`AskOnTaskRoot` 生产调用者＝0（只有本包用例）⇒ 交的是"门在、卡进得来"，不是"卡会自己来"。

### AC#3 — `Confirming` 期间 `Esc` 真能否决一次、完事归还（三形＋台件进仓） — **本腿初判：附条件成立**

- 台件是否在库（本腿现认）：`cmd/wisp/testdata/esclistener/main.go`（`git diff --stat` 显示 +317 行，第二进程观察者）⇒ 票 245 AC#9"搬进仓"的载体已落。
- 三形读数：winlive 档，本腿 §3 现跑（`借到/否决生效/归还`）。
- 附条件：⛔ 只测"借"不测"还"=假绿；本腿要求归还后第二进程真收到 Esc 那一形有具名读数。

### AC#4 — 卡片挂着时收到退出信号不许撒谎（真注册钩子、step3 不再 skipped、待决卡拒绝＋留审计；十步顺序一字不动） — **本腿初判：成立**

- 十步顺序未动（本腿现跑尺）：`git diff --name-only ce1ade1f HEAD -- internal/proc/shutdown.go` ⇒ **空**（`RunShutdownSequence` 一字未改）。`boot_windows.go:164 Shutdown` 改为读 `rt.shutdownHooks.Hooks()` 快照，顺序仍由该函数拥有。
- 钩子真注册（默认层现读）：`TestAC246CancelStepHookRunsOnTheRealShutdownSequence` → `records[StepCancelTasks-1].Skipped==false`、`Name=="cancel-task-roots"`、`Err==nil`；本腿 run1 日志见逐字 `退出第 3 步完成：拒绝待批卡片 0 张、作废 L1 窗口 0 张…`。
- 待决卡收口：`cancelTaskRoots:283` L2 走 `cards.Reject`→`ANSWER-REJECT`、L1 走 `RESIDENT-WINDOW-ABANDONED`（gate 对 ctx-cancel 不写行，本文件补写 `channel=none`）。winlive 一形读数见 §3。

### AC#5 — 起手第一发先归因那枚红 — **本腿初判：不成立（HEAD 全量默认层为红，且增量从未被完整验证）**

- ⛔ 关键：`246-r1` 死在 150 轮帽，编排者代提两枚增量（`17d8586e` 测试半边／`7c93ad7c` 产码两枚），⛔ 从没在 HEAD 完整态跑过 `go test ./cmd/wisp`。
- 本腿 run1（HEAD 全量默认层）：**1 枚红** → `TestAC228ExitRequestDuringBootStillLeavesThroughD38E`，逐字 `exit status 0xc000013a`＋`0 record(s)`（§3 R-B）。
- 本腿隔离复跑（`-count=5`，只跑这一枚）：**5/5 绿**（§3 R-C）⇒ 红非该用例固有，而是全包序下的时序／互扰。
- 归因：见 §3 R-D（run2 复现率＋配对台件定向复认）。

### AC#6 — 四条否决通道的其余三条不许顺手做 — **本腿初判：成立**

- 现读：`newResidentApproval` 用 `approval.NewChannels()`（四枚全未加载起步）；`bindBallHost:157` 只在真有窗口＋有执行者时 `SetLoaded(ChannelEsc, true)`。
- grep 归口（§3 R-E）：单击球 veto／KWS 否决词／面板拒绝本票零实现。

---

## §2 突变与正控

姿势：每发先确认锚点恰好 1 处、落盘、跑指名用例（判红绿只认 `--- FAIL`），跑完 `git cat-file blob HEAD:<path> > <path>` + `git diff --quiet` 复认。⛔ 不用 `-overlay`。

| 编号 | 攻哪一格 | 突变 | 指名用例 | 红了没有（红句逐字） | 还原 |
|---|---|---|---|---|---|
| MUT-6 | 代提 `7c93ad7c` 的 `!rb.cancelHosted` 守卫 | 把 `startResidentBall` 里 `cancelHosted: onCancelEsc != nil` 写死为 `true` | `TestAC246ChannelNeedsBothWindowAndExecutor` | （本腿稍后跑，读数进 §2.1） | |
| MUT-1 | AC#1/AC#4 装配 | 摘 `resident_windows.go:156 rt.RegisterShutdownHook(...)` 那一行 | `TestAC246ShippedResidentProcessOwnsItsCancelStep`＋`TestAC246CancelStepHookRunsOnTheRealShutdownSequence` | | |
| MUT-3 | AC#2/AC#3 借 | 摘 `resident_approval_windows.go:441 b.TakeEscForCancel()` | winlive `TestLive246ConfirmingCardBorrowsEscVetoesAndReturns` | | |
| MUT-4 | AC#3 还 | 摘 `settleOrb` 里 `b.ReleaseEscAfterSession()` | 同上（归还那一形） | | |

（§2.1 逐发红句：读数收拢中，跑完逐枚写入本节，⛔ 不预填。）

---

## §3 门禁读数（本腿现跑）

### R-A｜AC#1 零新增依赖边（两把尺，本腿自跑，非抄表）

| 尺 | 命令 | 读数 |
|---|---|---|
| 落地后包级传递依赖计数 | `GOOS=windows go list -deps ./cmd/wisp` | **276 行** |
| 落地后 internal/proc | `GOOS=windows go list -deps ./internal/proc` | **147 行**（内部 import＝0，仍叶子包） |
| 五枚 commit 新增的 internal import 行 | `git diff ce1ade1f HEAD -- <5 个产码文件> \| grep '^+.*CarlosShao/wisp/internal'` | 恰 5 行，全来自新文件 `resident_approval_windows.go`：approval / ball / risk / statemachine / tools |
| 这 5 枚是否 base 即 cmd/wisp 直依 | `git grep -l "wisp/<pkg>\"" ce1ade1f -- cmd/wisp/*.go`（排 `_test`） | approval→`approval_reply.go`+`run.go`；ball→`resident_ball_windows.go`；risk→`approval_always.go` 等 4 枚；statemachine→`resident_ball_windows.go` 等 4 枚；tools→`approval_always.go` 等 4 枚 ⇒ **5/5 base 即已直依** |
| internal/proc 新增 import | `git diff … -- shutdown_hooks.go boot_windows.go \| grep '^+\s+"'` | 仅 stdlib `context/errors/fmt/sync`（`sort` 一枚已被消，`shutdown_hooks.go:148` 注释具名） |

⇒ **包级新增依赖边＝0 枚**（新文件内部 import 全是 base 既有直依；`internal/proc` 只加 stdlib，未破叶子性）。

### R-B｜AC#5 起手红名册（HEAD 全量默认层）

- 命令：`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test ./cmd/wisp -count=1`（run1）
- 读数：`FAIL github.com/CarlosShao/wisp/cmd/wisp 174.803s`；`--- FAIL` **1 枚**：`TestAC228ExitRequestDuringBootStillLeavesThroughD38E (4.26s)`。
- 红句逐字：`resident_ball_228_windows_test.go:206: AC#1 RED: the child answered a boot-time Ctrl+C with exit status 0xc000013a, want a clean exit 0 through the D38(e) order.` ＋ `:214: … produced no "wisp: persistent log sink installed" record (records: [])`。
- 落盘：`head-default-tier-run1.txt`。

### R-C｜隔离复跑（判红是否固有）

- 命令：`go test ./cmd/wisp -run 'TestAC228ExitRequestDuringBootStillLeavesThroughD38E$' -count=5 -v`
- 读数：**5/5 `--- PASS`**，`ok … 25.457s` ⇒ 该用例单独跑稳定绿。落盘 `head-failing-test-count5.txt`。

### R-D｜归因（run2 复现率＋配对定向）

- run2（全包 `-v`）读数：收拢中（`head-default-tier-run2-v.txt`）。配对定向：收拢中。判"起手即在 vs 246 造成／代提增量造成"写入 §1 AC#5 与 §5。

### R-E｜AC#6 grep 三条通道未实现

- 读数收拢中：`ChannelBall`/`ChannelKWS`/`ChannelPanel` 在本进程 loaded 判据 + 单击球/KWS/面板拒绝零路由复认。

### 其余门禁（build/vet/gofmt/d22scan/winlive）

- 逐档收拢中；跑完写入本表，⛔ 不预填、不抄实现方读数。

---

## §4 留给编排者（本腿判了但需人拍板／归口的）

- 代提增量的"半发钉进历史"这一流程缺陷是否入台账具名（A482 缺失见 §5）。
- 红名册归因结论若指向"新增过程型兄弟用例互扰"，是否转票 228 后续片或本票补一刀隔离。
- winlive 三形/一形在 CI 无 `-tags winlive` 档（票 245 表 R-3 同口径）⇒ 回归只本机有读数。

---

## §5 判不动的地方（本腿不自裁，具名交回）

- **台账 `A482` 对不上**：`7c93ad7c` 提交信息写"台账 A482"，但 `docs/reports/pending-and-issues.md` 末节是 `A481`，全文 `grep "A482"`＝**0 命中**（`git log -S A482` 亦空）。⇒ 要么代提漏落一条账，要么提交信息虚指。本腿不改台账，具名交回编排者核。
- **"常驻进程是不是任务管线的永久之家"**：与实现方 §5-1 同判不动（`PLAN.md:73` D2 未逐字核）。
- **`Confirming → Sleeping` 收口态**：D43 第 22 行否决目标是 `Acting`，本腿无 loop 可回；该转移不在冻结表 ⇒ 判不动，等裁。
- **`problems()` 所指**：AC#4"不冒假账"那半，实现方按热键报告判，若原指退出审计别处则射程未盖到。

---

## §6 收尾名册对撞

- 终态 `git status --porcelain` 与 `roster-start.txt` 差集、逐枚 commit、写面是否清空、临时件清单——收尾现跑后逐字写入本节，⛔ 跑之前不预填。
