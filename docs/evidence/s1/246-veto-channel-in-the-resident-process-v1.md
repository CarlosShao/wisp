# 票 246 — 对抗验收腿 `246-v1` 判定表：最短链「注入 gate → 真卡片 → `Esc` 否决一次并归还 → 退出序列不撒谎」

> ⚠ **这份表的身份**：出自**非实现者**（`246-v1` 对抗验收腿）。按 `AGENTS.md §1.5`／`SPEC-12 §4.3`，本表是裁决判语；
> 实现方那份 `246-...-r1.md`（脏件、未提交完整、§2.1／§6 从未填满）**只当"待攻主张名册"，其读数不作本腿凭据**。
> 判据出处＝`.scratch/wisp/issues/246-...-no-veto-channel-has-an-executor.md` AC#0..AC#6；编排者裁定＝乙形＋"钩子注册入口不算契约面"，账 `A481`。
> ⛔ 本表**不翻任何 AC 框**（框归编排者）；本腿只交判语。所有读数本腿现跑，落盘 `.scratch/wisp/probes/246/v1/`（只建不删）。
> 本票语境（要紧）：写腿 `246-r1` 死在 150 轮帽，编排者**代提两枚增量**（`17d8586e` 测试半边／`7c93ad7c` 产码两枚），
> ⛔ **从没在 HEAD 完整态跑过 `go test ./cmd/wisp`**（只跑了 `go vet`）。⇒「HEAD 到底绿不绿」是本腿第一产出。

---

## §0 起手名册 · 锚点 · 主张名册

### §0.1 起手（闸门＝终态等于起手那一刻的名册，不是"必须为空"）

| 项 | 读数（本腿现跑） |
|---|---|
| 取数时刻 | `2026-09-30 19:36:20 +0800`（`date` 现跑） |
| 起手锚点 | `7c93ad7c`（`git rev-parse --short=8 HEAD` 自取，⛔ 未采用派单里给的任何 sha） |
| 分支 | `dev` |
| 起手名册全量 | `.scratch/wisp/probes/246/v1/roster-start.txt`（182 行；起手即含 `?? .scratch/wisp/probes/246/v1/`，因本腿先建目录再取册） |
| 桌面安静 | `wisp.exe`／`balldebug.exe`／`esclistener.exe` 计数＝**0／0／0**（每发真机／热键用例前复跑，均 0） |
| 本腿写面（白名单） | 仅 `docs/evidence/s1/246-...-v1.md` ＋ `.scratch/wisp/probes/246/v1/**`；⛔ 不改任何产码／工单框／冻结件／`frontend/**`／`design/**` |
| 五枚本票 commit | `75798f31` 证据骨架 → `0c1b32cf` 乙形落地 → `2b5462f9` 改法名融合 → `17d8586e` 测试半边（代提）→ `7c93ad7c` 产码两枚（代提）|

### §0.2 主张名册（实现方自述，逐枚本腿自己攻）

| # | 实现方主张（待攻） | 本腿攻法 | 落点 |
|---|---|---|---|
| C-1 | AC#1 零新增包级依赖边 | `go list -deps` 现跑 + base 直依逐名复认 | §3 R-A |
| C-2 | AC#2 真卡片挂起、真进 `Confirming`、非 mock | 问生产调用者 + winlive + 无窗口 fail-closed | §1 AC#2 |
| C-3 | AC#3 借→否决→归还三形齐、台件已进仓 | winlive 现跑 + 台件在库复认 | §3 R-F |
| C-4 | AC#4 钩子真注册、step3 不再 skipped、待决卡拒绝＋留审计、十步未动 | 读注册链 + winlive + `git diff` shutdown.go | §1 AC#4 |
| C-5 | AC#5「起手 2/2 绿、无红可归因」 | ⛔ 关键攻击：HEAD 全量默认层复样 | §3 R-B/C/D |
| C-6 | AC#6 其余三条通道未顺手做 | grep 通道 loaded + 路由 | §3 R-E |
| C-7 | 代提 `7c93ad7c` 的守卫不是装饰 | 定向突变 MUT-6a | §2 |

---

## §1 逐格判语（成立／不成立／附条件成立／判不动）

### AC#1 — 装配根注入真审批门、零新增包级依赖边 — **成立**
- 形状（现读）：`resident_windows.go:122 newResidentApproval()` → `:135 startResidentBall(rt.Registry, ra.vetoByEsc)` → `:144 bindBallHost(rb)` → `:156 RegisterShutdownHook(proc.StepCancelTasks, ra.cancelTaskRoots)`。gate 建在装配根，递给球宿主的是函数值 `escVetoFunc`（`resident_ball_windows.go:61`）。
- 尺＝本腿自跑（§3 R-A）：新文件内部 import 全 5 枚 base 即已直依、`internal/proc` 只加 stdlib ⇒ **包级新增依赖边＝0**。出厂真进程格由 `TestAC246ShippedResidentProcessOwnsItsCancelStep`（默认层，真 `wisp.exe`）守住，run1/2/3 无其红。

### AC#2 — 一张真卡片能在常驻进程里挂起（可见性只到"日志＋状态位＋Win32"）— **成立（射程内）**
- 判据红线：⛔ 面板开不出来是既成事实（`approval_always.go:165` 逐字 "this binary links no WebView2 host"）⇒ 任何"用户看得见卡片"一律判不成立；本腿**未见也不写**此类主张。
- 非 mock：gate＝真 `approval.New`（`:109`），UI＝真 `ballCardUI`（呈现接缝，非假审批门），`askConfirmation:205` 走真 `gate.AdmitTextTask`+`gate.PendingWindow`(L1)/`PendingApproval`(L2)——与 `wisp run` 同一族机制。winlive 现读：L1 卡片 `WaitingState()==Confirming`、L2 `==AwaitingApproval`（两枚 D43 名，不造第五）。
- **生产调用者这一问（本腿现跑 grep）**：`newResidentApproval()` **有**生产调用者＝`resident_windows.go:122`（双击起来的二进制里装配并绑定了 gate）；但 `askConfirmation`/`AskOnTaskRoot`（**举起**卡片）生产调用者＝**0**，只有本包用例 ⇒ 常驻进程今天"门在、卡进得来"，尚无自主任务源把卡举起来。⇒ 与实现方 R-2 自陈一致，⛔ 不构成"卡会自己来"，也⛔ 不构成用 mock 代真（`askConfirmation` 是产码、真 gate，用例只是触发它）。
- 无窗口反证（默认层，本腿 run1 日志现读）：真 gate 逐字 `approval: 卡片无处呈现（本进程没有悬浮球窗口），已 fail-closed 拒绝` → 拒。

### AC#3 — `Confirming` 期间 `Esc` 真能否决一次、完事归还（三形＋台件进仓）— **成立**
- 台件在库（现认）：`cmd/wisp/testdata/esclistener/main.go`（第二进程观察者：自起窗、`-steal` 注全局热键作正控、`keybd_event` 注物理 Esc、读 `foreground_is_mine`）⇒ 票 245 AC#9"搬进仓"载体已落（翻勾归编排者）。
- winlive 现跑三形（`-tags winlive`，桌面安静）：**借到** `keydown 1→0`、**否决生效** `ANSWER-VETO corr=host:live246-esc … channel=esc decision=veto`、**归还** `keydown 0→1`（归还由**第二个进程**收到，非自家记账）；正控 `-steal keydown=0/hotkey=1`、基线 `-watch keydown=1`。逐字 `AC#3 READING: borrow/veto/return in 144.659ms; veto=veto; observer keydown 1->0->1 (steal control 0)`。⛔ 只测"借"的形状不存在。

### AC#4 — 卡片挂着时收到退出信号不许撒谎（真注册钩子、step3 不再 skipped、待决卡拒绝＋留审计；十步顺序一字不动）— **成立**
- 十步顺序未动（本腿尺）：`git diff --name-only ce1ade1f HEAD -- internal/proc/shutdown.go`＝**空**（`RunShutdownSequence` 一字未改）；`-tags winlive ./internal/proc` 复跑 `TestShutdownOrderAudit` **PASS**。`boot_windows.go:164 Shutdown` 改读注册快照 `rt.shutdownHooks.Hooks()`，仍填 `CloseJob` 兜底，顺序仍由该函数拥有。
- 钩子真注册（默认层，run1 日志现读）：`TestAC246CancelStepHookRunsOnTheRealShutdownSequence` → `records[StepCancelTasks-1].Skipped==false`、`Name=="cancel-task-roots"`、`Err==nil`，逐字 `退出第 3 步完成：拒绝待批卡片 0 张…任务根已取消，无等待残留`。
- 待决卡收口（winlive 现跑）：L2 `ANSWER-REJECT corr=approval-1 … decision=reject reason="常驻进程退出：未获批准，按拒绝处理（未执行）"`、step3 `no-skipped=true`；L1 无"拒绝"动词 ⇒ 本文件补写 `RESIDENT-WINDOW-ABANDONED … decision=reject … channel=none`（明说无人否决）。
- "不冒假账"：winlive 断言收口后 `HotkeyReport().Problems()==0`（归还回 standby，不是注册失败）；默认层六枚 `shutdown step skipped (module not present)` 是本票不拥有的步，老实记 skipped，⛔ 未伪造成 executed。

### AC#5 — 起手第一发先归因那枚红 — **不成立（就"增量已被完整验证／红已归因"这一主张）**
- ⛔ 核心：`246-r1` 从没在**补齐后的 HEAD** 跑过整包 test（代提者自陈只跑 `go vet`）。实现方 §AC#5 写"起手 2/2 绿 ⇒ 无红可归因"，但那两发跑在**动任何产码之前**（锚点 `ce1ade1f`），不是这枚增量。
- 本腿复样（HEAD 全量默认层，桌面安静、无并发）：**run1 红 1 枚 / run2 绿 / run3 绿 ⇒ 3 发中 1 发红**，红名恒为 **`TestAC228ExitRequestDuringBootStillLeavesThroughD38E`**（红句见 §3 R-B）。⇒ HEAD **非稳定绿**。
- 归因结论（本腿现跑三把尺）：① 隔离 `-count=5` 只跑它＝**5/5 绿**（R-C）；② 把它与最相关的 246 兄弟用例 `TestAC246ShippedResidentProcessOwnsItsCancelStep` **配对 `-count=6`＝6/6 绿**（R-D）；③ 红发生在 `resident_windows.go` 的 `[installLogSink→signal.Notify]`（约 62→106 行）这段**本票没碰**的时序里（本票新增代码全在 `Notify` 之后的 122 行起）。⇒ **该红＝起手即在型 boot-Ctrl+C 时序 flake，⛔ 非本票逻辑造成**，与台账 `A480` ④"1 红 2 绿、用例名丢了"高度吻合（**本腿把丢了的名字找回来了**）。
- 但本票**并非毫发**：它给默认层净增 2 枚**起真进程／建真球窗**的用例（`TestAC246ShippedResidentProcessOwnsItsCancelStep` + `TestAC246ChannelNeedsBothWindowAndExecutor`），与那枚 flake 同挤一个 test 二进制 ⇒ **可能抬高其发生率**（方向可推，⛔ 3 个样本不足以证，留 §5 交回）。
- 综合判语：AC#5 要求"起手先归因、⛔ 不许当已知常红略过"。实现方以"2/2 绿"收工＝**未对这枚增量做完整验证、也未把该 flake 具名**；这枚红由**本验收腿**才真正具名＋归因。⇒ 主张**不成立**（红存在且曾漏归），产品事实"这条链在生产里能跑"**不被这枚 flake 否定**（flake 属票 228 地界、非本票产码）。

### AC#6 — 其余三条否决通道不许顺手做 — **成立**
- grep（§3 R-E）：`cmd/wisp` 产码里唯一 `SetLoaded(...,true)`＝`resident_approval_windows.go:157` 只载 `ChannelEsc`；`ChannelBall/KWS/Panel` 全程未载；`OnClickBall→recordBallGesture("click")`（只记账、不 veto）。⇒ 只落 `Esc` 一条，另三条各有归口。⚠ `approval_reply.go:469` 那枚 `SetLoaded(vetoChannel,true)` 属**既有 `wisp run` 控制台腿**，非本票、该文件不在本票 diff 内。

---

## §2 突变与正控（本腿现跑）

姿势：单点突变 → 跑指名用例（判红绿只认 `--- FAIL`）→ `git cat-file blob HEAD:<path> > <path>` 还原 + `git diff --quiet` 复认。⛔ 未用 `-overlay`／`checkout`／`stash`／`reset`。

| 编号 | 攻哪一格 | 突变（单点） | 指名用例 | 红了没有 | 还原 |
|---|---|---|---|---|---|
| MUT-6a | 代提 `7c93ad7c`：`bindBallHost` 里"有窗无执行者⇒不载 Esc"那支守卫是否装饰 | `resident_approval_windows.go:145` `if !rb.cancelHosted {` → `if false && !rb.cancelHosted {`（中和守卫） | `TestAC246ChannelNeedsBothWindowAndExecutor` | **红**（见下逐字） | `git diff --quiet` 复认＝clean |

**MUT-6a 红句逐字**（`head-default-tier` 隔离跑，`mut6a-guard-neutralized.txt`）：
```
resident_approval_246_windows_test.go:319: bindBallHost loaded the cancel channel with no executor behind it (ball up = true)
--- FAIL: TestAC246ChannelNeedsBothWindowAndExecutor (2.46s)
```
⇒ 守卫**拿掉后**、"有球窗口但装配根没注入取消执行者"那一支被绕过，Esc 就被 advertise ⇒ 用例当场红 ⇒ **代提那半发增量不是文案、是承重的**；`ball up = true` 亦证这是在本机真球窗下跑的真实路径，不是塌到"无窗"那一形。

（未做更多突变：本票三形／退出两形已有 winlive 正向读数，且 HEAD 稳定绿受 §1 AC#5 那枚 flake 干扰，再叠突变只增噪声；留 §5 具名。）

---

## §3 门禁读数（本腿现跑；逐名点名）

### R-A｜AC#1 零新增依赖边
| 尺 | 命令 | 读数 |
|---|---|---|
| 落地后 `-deps`（cmd/wisp） | `GOOS=windows go list -deps ./cmd/wisp` | **276 行** |
| 落地后 `-deps`（internal/proc） | `GOOS=windows go list -deps ./internal/proc` | **147 行**（内部 import＝0，仍叶子） |
| 本票新增 internal import 行 | `git diff ce1ade1f HEAD -- <5 产码文件> \| grep '^+.*wisp/internal'` | 恰 5 枚，全来自新文件：approval/ball/risk/statemachine/tools |
| 这 5 枚 base 即 cmd/wisp 直依？ | `git grep -l "wisp/<pkg>\"" ce1ade1f -- cmd/wisp`（排 `_test`） | **5/5 命中**（approval→approval_reply/run；ball→resident_ball；risk/statemachine/tools→各≥3 枚既有文件） |
| internal/proc 新增 import | `git diff … -- shutdown_hooks.go boot_windows.go` | 仅 stdlib `context/errors/fmt/sync`（无 `sort`，`shutdown_hooks.go:148` 注释具名） |
⇒ **包级新增边＝0**（新文件内部 import ⊆ base 既有直依；proc 只加 stdlib）。

### R-B｜起手红名册（HEAD 全量默认层，run1）
`go test ./cmd/wisp -count=1` → `FAIL … 174.803s`；`--- FAIL` **1 枚**：`TestAC228ExitRequestDuringBootStillLeavesThroughD38E (4.26s)`。
红句逐字：`resident_ball_228_windows_test.go:206: AC#1 RED: the child answered a boot-time Ctrl+C with exit status 0xc000013a … :214: … produced no "wisp: persistent log sink installed" record (records: [])`。
（⛔ 非"根本没跑"那一形：本发跑了 174.8s、有真 `--- FAIL`；缺 sherpa 会是 `0xc0000135`＋`0.0xxs`＋无 `--- FAIL`。）

### R-C｜隔离复跑（判红是否固有）
`go test ./cmd/wisp -run 'TestAC228ExitRequestDuringBootStillLeavesThroughD38E$' -count=5 -v` → **5/5 `--- PASS`**，`ok … 25.457s`。⇒ 该用例单独跑稳定绿。

### R-D｜归因三把尺
- 全包复样：run1 红／run2 绿／run3 绿 ⇒ **3 发 1 红**，红名恒为该枚。
- 配对（最相关 246 兄弟用例紧邻它跑）：`-run 'TestAC246ShippedResidentProcessOwnsItsCancelStep$|TestAC228ExitRequestDuringBootStillLeavesThroughD38E$' -count=6` → **6/6 绿**。
- 代码位置：红在 `[installLogSink→signal.Notify]`（约 62→106 行）区间，**本票未改**（本票新增全在 `Notify` 之后）。
⇒ **起手即在型时序 flake**（≈ `A480`④ 那枚丢名的红）；非本票逻辑红；本票净增 2 枚起进程／建球窗用例的负载**可能**推高其发生率（未证，3 样本不足）。

### R-E｜AC#6 grep
`grep -rn SetLoaded cmd/wisp --include=*.go`（排 test）⇒ 仅 `resident_approval_windows.go:157 ChannelEsc=true`／`:351 ChannelEsc=false`；`approval_reply.go:469 vetoChannel` 属 run 控制台腿（非本票）。`OnClickBall: recordBallGesture("click")`（无 veto）。⇒ 只落 Esc。

### R-F｜winlive 真机（本腿现跑，桌面安静）
`PATH=… go test -tags winlive ./cmd/wisp -run 'TestLive246' -v` → **3 枚全 PASS（8.08s）**：
- `TestLive246ConfirmingCardBorrowsEscVetoesAndReturns` PASS — 借/否决/归还三形（§1 AC#3 逐字）。
- `TestLive246ExitRefusesAHangingL2Card` PASS — `ANSWER-REJECT … 未获批准，按拒绝处理`、`no-skipped=true`。
- `TestLive246ExitAbandonsAHangingL1Window` PASS — `RESIDENT-WINDOW-ABANDONED … channel=none`。

### 其余门禁（本腿现跑）
| 尺 | 命令 | 读数 |
|---|---|---|
| internal/proc 默认层 | `go test ./internal/proc -count=1` | `ok … 2.282s`（绿） |
| internal/proc winlive | `go test -tags winlive ./internal/proc` | `ok … 0.874s`；`TestShutdownOrderAudit`／`TestShutdownHookSetRunsEveryRegisteredStepInOrder` **PASS** |
| build | `go build ./...` | 无输出（绿） |
| vet | `go vet ./...` | 无输出（绿） |
| vet（winlive 档） | `go vet -tags winlive ./...` | 无输出（绿） |
| vet（台件） | `go vet ./cmd/wisp/testdata/esclistener` | 含于上行 rc=0（`./...` 走查不吃 `testdata`，此路已由 `TestLive246` 显式 build 它 rc=0） |
| gofmt | `gofmt -l internal/proc cmd/wisp` | 空 |
| d22scan | `sh scripts/d22scan.sh` | 正控 PASS=34/FAIL=0 → 真扫描 `clean`，examined 255 production Go files（bans #1-8 无违规） |

---

## §4 留给编排者（判了但需人拍板／归口）

- **`cmd/wisp` 默认层那枚 flake**：本腿归因为起手即在（票 228 的 `[sink→Notify]` 时序），⛔ 非本票产码；是否单独立票修（如：把 `signal.Notify` 提到 `installLogSink` 之前）、还是并入票 228 后续片，请裁。本腿不动它。
- **HEAD 从未被写腿完整验证**：代提两枚增量此前只有 `go vet` 读数；本腿给了它第一次整包（默认＋winlive）验证。若编排者要"整包 `go test ./...`"级别背书，本腿射程只到 `cmd/wisp`／`internal/proc`／`internal/agent/approval`，其余包未跑。
- **AC 框翻勾**：全部归编排者；本腿一枚未动。

---

## §5 判不动的地方（本腿不自裁，具名交回）

1. **台账 `A482` 不存在**：派单/背景把 `7c93ad7c` 记作"产码两枚＋台账 A482"，但本腿现跑：`grep -c "A482" docs/reports/pending-and-issues.md`＝**0**、`git diff --stat ce1ade1f HEAD -- docs/reports/pending-and-issues.md`＝**空**、5 枚 commit message 内 `grep A482`＝**0**。⇒ 台账末节仍是 `A481`，那条代提台账增量**从未落地**。本腿不改台账，具名交回核（是漏落、还是背景里那句 A482 系虚指）。
2. **本票那 2 枚新真实进程／球窗用例是否抬高了 flake 发生率**：方向可推（同挤一个 test 二进制、彼此串扰），但本腿仅 3 个全包样本，⛔ 不足以定"变严重了没"，不下判语。
3. **`Confirming → Sleeping` 收口态**：D43 第 22 行否决目标是 `Acting`，本腿无 loop 可回、选择回 `Sleeping`（`resident_approval_windows.go:492`）；该转移不在冻结表上 ⇒ 判不动，等裁（与实现方 §5-3 同一问）。
4. **"常驻进程是不是任务管线的永久之家"**／`PLAN.md:73` D2 未逐字核 ⇒ 判不动（同实现方 §5-1）。
5. **winlive 这族在 CI 无面**：`ci.yml` 无 `-tags winlive` 档 ⇒ 三形／退出两形的回归只在**这台机器手跑**有读数；本腿不改 CI、不造门（票 245 R-3 同口径）。
6. **`approval_always.go:165` 一句疑过期**："the ball is built only by cmd/balldebug" 在票 228／246 把球搬进 `cmd/wisp` 常驻腿后**可能已失真**。本腿**未改**（非本票射程），具名归口票 228 文档后续片。

---

## §6 收尾名册对撞（本腿现跑）

- 终态 `diff <(git status --porcelain) .scratch/wisp/probes/246/v1/roster-start.txt` → **rc=0（逐字节相同）**：本腿的 `?? .scratch/wisp/probes/246/v1/` 起手即已折叠在册（先建目录再取册），新增的都是该折叠项内的未跟踪探针文件，不进 `git status` 顶层。
- 产码面：`git diff --quiet -- internal cmd docs/PLAN.md` → **clean**（MUT-6a 单点突变已 `git cat-file blob HEAD:` 还原并复认；⛔ 未碰冻结件／`frontend**`／`design**`／工单 AC 框）。
- 本腿唯一被跟踪的产出＝`docs/evidence/s1/246-veto-channel-in-the-resident-process-v1.md`（本表，已 commit）。实现方那份 `...-r1.md` 仍 `M`（本腿一枚未动）。
- 探针件清单（`.scratch/wisp/probes/246/v1/`，只建不删）：`roster-start.txt`／`head-default-tier-run1.txt`／`head-default-tier-run2-v.txt`／`head-default-tier-run3.txt`／`head-failing-test-count5.txt`／`attribution-pairing-count6.txt`／`mut6a-guard-neutralized.txt`／`winelive-246-run1.txt`／`proc-winelive.txt`／`d22scan.txt`／`msg-skeleton.txt`。
- commit：骨架 `14676b1c`；终表 commit 号见其下（本腿只 commit、不 push）。
