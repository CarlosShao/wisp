# 378 裁决（378-v2 复审，非实现者）：通过

- 裁决者 ≠ 实现者：产生码与证据由编排者本会话做（`e76712fb8`／`9469e9852` 等）；本件由裁决腿独立复核并现跑。
- 锚：HEAD `9469e9852`（2026-11-04）。本件是**同一文件的 378-v2 复审**，取代 378-v1 的退回结论；v1 原文保留在文末附录（只追加、不删记录）。
- **结论：AC#2 成立**（真 `task.spawn` 之后 `PRAGMA foreign_key_check = 0`，子任务 `task_log` 恰 1 行）；v1 的 R1／R2／R4 已闭合，**R3 归口票 405 可接受**（另立票、tracked、open，有独立 AC）。本票可收口。

## 0. 一句话判语

v1 退回的两条决定性理由——「尺看不见 `run.go:1583` 那一跳」与「验收数在树里没有可复现用例」——现在**同时被同一枚用例反证**：基线绿给出 `fkc=0` 与子行恰 1 行；把 `rt.loopOpt = loopOpts` 改回 `rt.loopOpt = opts` 后同一枚用例**逐名转红**（`fkc=1`、子行 0 want 1）并指出孤儿 `tool_call` 行。`tasklog_fkcheck_378_test.go` 是**真的跑级**用例（真 `wisp run`、真 `task.spawn`、真桥、真 SQLite），不是装配级复现。

## 1. 门禁四数（本腿现跑，四枚全绿）

| 尺 | 读数 |
|---|---|
| `go build ./cmd/wisp` | **rc=0** |
| `go vet ./cmd/wisp` | **rc=0** |
| `go test ./cmd/wisp -run Test378SpawnedChildToolCallsHaveAParentRow -count=1 -timeout 150s` | **`ok github.com/CarlosShao/wisp/cmd/wisp 0.396s`，rc=0**；`-v` 复跑 `--- PASS (0.20s)`／包 `ok 0.266s`；`-count=3` 三发 `ok 0.718s`（无 golden 空表、无 31s 确认窗口） |
| `gofumpt -l cmd/wisp/tasklog_fkcheck_378_test.go` | **空**（rc=0） |

附加现跑（不入四数，只作归因材料）：

- `go test ./cmd/wisp -run Test233HotTierAgentReachesTheNextTask -count=1` ⇒ **`ok 2.831s`**：R4（试 404 修 `screenCaptureSchema`）独立复核转绿，跑级通道确实解锁。
- `go test ./cmd/wisp -run '378' -count=1` ⇒ **`ok 0.787s`**：4 枚 378 用例（3 个测试文件）同包同跑不互扰。
- `sh scripts/d22scan.sh` ⇒ **clean，rc=0**（`examined 406 production Go files`；`ban #8 internal/=855`／`cmd/=200`）。
- `go build ./...` ⇒ **rc=1**，但红的全部是**未跟踪**的 `build/mut_rvw.go`（`undefined: voiceTtsFiles` 等 10 行；`git ls-files build/mut_rvw.go` = **0**）；`go build ./cmd/... ./internal/...` ⇒ **rc=0**。⇒ 这是本机 staging 件的存量假红，不属本票，也与 CI 干净检出无关。
- `sh scripts/check-path-length-budget.sh --with-self-test` ⇒ **VERDICT RED**，**但存量、非本票**：树 65 枚 over-budget／roster 62 枚；5 枚未 roster 的 `-done` 名（`128`/`175`/`273`/`364`/`403`）＋2 条旧名 stale 键（`128`/`175`）**在 `f50ea6747` 的树里已存在**（现跑 `git ls-tree -r --name-only f50ea6747 | grep -c` = 5）。本票文件名 `378-child-…-parent-done.md` = **86 字符／相对 107**，≤ 100／121 帽，**零新增违规**、不入 offender 集。该树级红具名交编排者（本腿⛔碰 roster／threshold）。

## 2. AC#2 基线读数（真 spawn，fkc=0）

用例 `cmd/wisp/tasklog_fkcheck_378_test.go` 由 `runTextTask` 驱动，模型经 C5 golden SSE 注入（AGENTS §1.3），其余全生产件：真 `config.toml` → 真 `secret.Store` → 真 `approval.Gate` → 真 `tools.Bridge`（`Journal: mem`）→ 真 `agent.Loop` → 真 `subagent_197` 子环路 → 真 `memory.Store`。

- golden 收到 **4 发**请求：root spawn → child `fs.read` → child 收口 → root 收口（顺序确定，`task.spawn` 阻塞到子任务 join）。
- root=`ed8d8c4a-…`，child=`8fad431a-…`（本例实际 id；名册中父下恰 1 枚子代理行）。
- **`PRAGMA foreign_key_check = 0 行`**（独立新连接、`memory` 的 driver、10s 超时）。
- **子任务 `task_log` 恰 1 行**（`childRowCount == 1`；`query_text` = 子任务自己的任务书，非空）。
- **子任务 `tool_call` 恰 1 行**：`fs.read / L0 / allow / success`——所以上面那个 0 不是「没有子调用」的空 0。
- root 正控：`task.spawn / success` 一枚 + root `task_log` 行可读；`task_log` 共 2 行（root=1 child=1）。
- 总量控制：`tool_call` 总行数 = **2**（父 spawn＋子 read），继承的 journal 没有给每发调用再写第二行（硬提醒 #1 不爆）。

这四组断言 [A] 子行存在 / [B] 子 `task_log` 行存在且恰 1 行 / [C] root 正控 / [D] 总量 = 模型调用数，都在同一发里被断言，所以 `fkc=0` 有对象、不是空库绿。

## 3. 反向读数（本腿自做变异 → 逐名红 → 还原核）

变异（唯一一发，时机＝最后一发变异之后才还原）：

```text
sed -i 's/^\(\s*\)rt\.loopOpt = loopOpts$/\1rt.loopOpt = opts/' cmd/wisp/run.go
# run.go:1583:  rt.loopOpt = loopOpts   ->   rt.loopOpt = opts
```

改后同一命令（`-run Test378SpawnedChildToolCallsHaveAParentRow -count=1 -v`）：

```text
[378] fk 违规行：tool_call rowid=1 -> task_log.id fk#0：task_id=6b21847b-… tool=fs.read（这发工具调用没有父任务行）
[378] PRAGMA foreign_key_check = 1 行（root=3d81bc2d-… child=6b21847b-…）
PRAGMA foreign_key_check = 1 行，want 0：一发 task.spawn 之后子任务的 task_log 父行必须真的存在
[B] 子任务 6b21847b-… 在 task_log 里没有行（memory: not found）
[B] 子任务 6b21847b-… 在 task_log 里的行数 = 0，want 1
[378] task_log 共 1 行（root=1 child=0）
--- FAIL: Test378SpawnedChildToolCallsHaveAParentRow (0.26s)
```

⇒ **`fkc=1` want 0**、**子行数 0 want 1**，且违规行**逐名报出**（`tool_call` 行的 task/tool 被读回）；root 那一支仍 0（票 373 未回归）。这就是 AC#2 要求的「摘掉新加的那一跳 ⇒ 用例转红并逐名报红行」，v1 的 R1「不修不响」不成立已由本轮反证**闭合**。

还原核：

```text
git checkout -- cmd/wisp/run.go
git diff --stat -- cmd/wisp/run.go     -> 空
sha256sum cmd/wisp/run.go              -> 6df3f77d05eaf24b30b59d04c08a8ca59be0283677e7e93cbd8753a8960d4362
git show HEAD:cmd/wisp/run.go | sha256sum -> 同一串
```

## 4. R3 归口票 405：**可接受**（两项具名条件）

**为什么可以拆**：378 的核心是**取证链**——子任务的 `tool_call` 行必须有一个真的父行，判据是 `fkc`。子行开出来之后**谁把它关、终态叫什么**，不改变 `fkc=0` 也不改变「工具调用有父」；它是**用户可见的 lifecycle/展示缺陷**，与父行存在性正交。

**405 的承载是实的（本腿现读）**：

- 票 tracked（`git ls-tree -r HEAD --name-only | grep -c 'issues/405-'` = **1**）、`Status: open`，不是口头约定。
- 机制描述对：`taskLifecycleJournal.FinishTaskLog` 故意惰性（`run.go:1326-1328`，现读 `return nil`）；全仓只有 root 收尾关行（`run.go:1742`）；子任务无关闭点 ⇒ 子行 `running`（本腿基线的 `state=running` 正是它）；`crash_recovery` 把仍 `running` 的行改成 `interrupted`（本腿现读 `cmd/wisp/crash_recovery.go:145-165` 的 `TaskStateRunning -> TaskStateInterrupted` 循环）。
- 影响与 AC#1 候选（①`FinishTaskLog` 对子任务写终态／②`run.go` 侧收尾关行／③按 `SPEC-12 §5` 五字段登记推迟）都写在票面，**没有被静默吞掉**。
- 它**现在不是 DEFERRED 标记**，是一张要做的活票 ⇒ 不触发 `SPEC-12 §5` 五字段；**若 405 的 AC#1 最终选 ③**，那时才必须补五字段并让 `DEFERRED(D-xx)` 与 §5 1:1。
- 378 的验收用例**只记录不判定**子行 `state`（`tasklog_fkcheck_378_test.go:480-483` 明写「没有终态写手，故不断言 state」）⇒ 它没有把 `running` 冒充成「完成」。

**两条具名条件（不阻塞收口，但要落到 405）**：

1. 405 的引用路径写错：票面写 `internal/memory/crash_recovery.go:150-159`，实际文件是 **`cmd/wisp/crash_recovery.go:150-159`**（`internal/memory/crash_recovery.go` 本腿现读**不存在**）。机制不变，但 405 的写腿应按真实路径核，票面路径请编排者更正（本腿⛔碰 405）。
2. 405 必须保持 open 直到终态写手落地或按 §5 正式登记；378 的 `-done` 只表示**父行/取证链**这一格完成，不表示子行终态已正确。

## 5. AC#0–4 逐格判语

- **AC#0：成立（以更正后的归因为准）**。① 位置集合今天是：base `run.go:1521 Journal: nil` → `:1547` 发布 base → `:1578` 挂 `taskLifecycleJournal` → **`:1583` 再发布到 `rt.loopOpt`**（子环路经 `subagent_197.go:306/:320` 复制它）。`378-a0.md:9` 的结论「子环路从 rt.loopOpt（Journal 仍 nil）复制」在它自己的锚上**已过期**：`:1583` 由 `e76712fb8`（2026-11-03 23:23）加入，且 `git merge-base --is-ancestor e76712fb8 3dec5a0ae` 返 0（a0 锚 2026-11-04）⇒ 正确归因是「修复已在树里，缺的是尺」，本轮反向读数就是那把尺。a0 原文不删（记录），以本件与 v1 R1 的更正为准。②④ 由第 2 节补齐（子 fkc=0＋root 正控＋子行恰 1 行）；③⑤ 的代价面（`ListTaskLogs`／`PurgePrivacy`／retention）与替代解评估仍成立。
- **AC#1：成立**。〔甲〕（每个子任务 1 枚 `task_log` 行、工具调用行仍归 bridge、不按调用写行）已由实现与第 2 节读数确认；终态缺口拆到 405（见第 4 节），不是把 `running` 当正确。
- **AC#2：成立**。第 2 节＋第 3 节：验收数可复现、反向逐名红、还原核通过。
- **AC#3：本票面成立（两处树级存量红具名）**。四数绿＋`d22scan` clean；`go build ./...` 的红仅来自未跟踪 `build/mut_rvw.go`（`./cmd/... ./internal/...` rc=0）；长度门 RED 的 5 枚 offender 在 `f50ea6747` 已存在、本票路径 86/107 在帽内。⇒ 本票**零新增门禁红**；树级两红不据本票签绿，具名交编排者。
- **AC#4：成立**。`9469e9852` 的 name-status 只有 `A .scratch/wisp/issues/405-…`／`A cmd/wisp/tasklog_fkcheck_378_test.go`／`M docs/reports/pending-and-issues.md`；378 链未碰 373／368／366 原句，未碰 thresholds／golden／frontend／design；本腿未 push。

## 6. 与 378-v1（R1–R4）的逐条对账

| v1 条目 | 现状 | 依据 |
|---|---|---|
| **R1** 两枚用例都看不见 `run.go:1583` | **闭合** | 第三枚用例在 `rt.loopOpt = loopOpts` 变异下转红并逐名报孤儿行（第 3 节） |
| **R2** 验收数在树里没有可复现用例 | **闭合** | `tasklog_fkcheck_378_test.go` tracked、三读数绿、基线 fkc=0＋子行 1（第 2 节） |
| **R3** 子行永不关、`crash_recovery` 误标 | **归口 405**（可接受，见第 4 节） | 405 tracked/open；378 用例不冒充终态 |
| **R4** `screenCaptureSchema` 坏 JSON 封死跑级 | **已修（404）** | 本腿 `Test233` `ok 2.831s`；本票真跑用例绿 |

## 7. 本腿没做的事

不 push；不碰 373／368／366／404／405 的票面原句；不碰 `run.go`（除第 3 节临时反向并已还原，sha256 与 HEAD 同串）；不碰 thresholds／golden／frontend／design；不改任何 C/D 契约文字；不写 `R##` 台账（编号归编排者）；未跑 `go test ./cmd/wisp` 整包（只跑本票与相邻 378 用例、`Test233`）与 `go vet ./...`（只跑动过的 `./cmd/wisp`）。

---

## 附录：378-v1 退回记录（历史，已由上面的 v2 取代；原文保留）

# 378 裁决（378-v1，非实现者独立复核）：退回 R-378-1

- 裁决者 ≠ 实现者：本票的产生码与证据由编排者本会话做（`6e61456e1` 等）；本件由裁决腿独立复核并现跑。
- 锚：HEAD `6e61456e1fe7f5b94725ca12242c347e5790583c`（2026-11-04）。
- 结论：**AC#2 不成立**（R1/R2：尺看不见被改的那一跳，验收数在树里没有可复现的用例）；**R3** 是甲形留下的用户可见假终态，属未裁/未登记。**不翻票面任何框**，票保持 `ready-for-agent`。

## 门禁四数（本腿现跑，四枚全绿）

| 尺 | 读数 |
|---|---|
| `go build ./cmd/wisp` | rc=0 |
| `go vet ./cmd/wisp` | rc=0 |
| `go test ./cmd/wisp -run TestTicket378 -count=1 -timeout 150s` | `ok` 0.336s，rc=0 |
| `gofumpt -l cmd/wisp/child_task_log_378_test.go` | 空（rc=0） |

附加现跑（不入四数，只作归因材料）：
- `go test ./cmd/wisp -run Test378ChildLoopWritesOneTaskRowAndNoToolCallRows -count=1` => **PASS** 0.15s（树里既有的真 `task.spawn` 装配级用例，见 R2）。
- `go test ./cmd/wisp -run Test233HotTierAgentReachesTheNextTask -count=1` => **FAIL**（`wisp run exit 1`，`json: error calling MarshalJSON for type *jsontext.Value: invalid character '"' after object key:value pair`）。见 R4，本票不处置，但它封死了"跑级"验收数的路。

## AC#0-4 逐格判语

| 格 | 判语 | 依据 |
|---|---|---|
| AC#0 | **部分成立**。③⑤（代价面静态位置、替代解）够用；**① 归因错一条提交**；**②④ 只在未入库的件里量过一次** | ① `378-a0.md:9` 称"`rt.loopOpt` 保持 `Journal` nil、子环路从它复制"，但 `git blame -L 1575,1584` 显示 `run.go:1583` 的 `rt.loopOpt = loopOpts` 由 `e76712fb8`（2026-11-03 23:23）加入，**且该提交是 a0 自述锚 `3dec5a0ae`（2026-11-04 17:33）的祖先**（`git merge-base --is-ancestor` 返 0）=> a0 把**修复后**的注释（`run.go:1579-1582`，正由修复提交写下的那段）读成了**修复前**的证据。②④ 由 `docs/evidence/s1/378-r1-impl.md` §3 给出过读数（fkc=0 / 子行 state=running），但该件 **untracked**，且量它的用例已被移出（见 R2） |
| AC#1 | **裁形本身成立**（〔甲〕：子任务独立 `task_log` 行，符合主流 harness），**但代价面缺一项** | 甲只写"任务行"而不重开工具调用行，避开硬提醒#1 的量级爆，方向对。缺的是这些行的**终态**（R3），硬提醒#3/#4 要的正是用户可见面与老化面 |
| AC#2 | **不成立** | R1（尺看不见新加的那一跳）＋ R2（验收数没有可复现的用例） |
| AC#3 | 四数绿，**不可据此签** | `d22scan`／长度门本腿未跑（AC#2 已退回，签 AC#3 无意义）。`go build ./...` 的既有假红（a2 证据称是未跟踪的 `build/mut_rvw.go`；本腿只确认该文件确在盘上，未复核它是否就是那一红）不属本票 |
| AC#4 | **通过**（本腿只读） | 378 链提交（`3bf66a7f6`/`fe7573aaf`/`f63688cef`/`7b54388c4`/`596f34977`/`22459ef03`/`6e61456e1`/`843a6d172`）只碰 `cmd/wisp` 的 378 两枚测试与 `docs/evidence/s1/378-*`＋台账；未碰 373/368/366 票面，未碰 thresholds/golden/frontend/design |

## R1（决定性，"不修不响"不成立）：树里的两枚 378 用例都看不见本票唯一的那一跳

本票唯一的产生码改动是 `cmd/wisp/run.go:1583` 的 `rt.loopOpt = loopOpts`（`e76712fb8` 加，`git blame` 可核）。今天在树里的两枚 378 用例**都不读它**：

- `cmd/wisp/child_task_log_378_test.go:44` 自己造 `j := &taskLifecycleJournal{store: st}`，只驱动"journal 类型 -> Store"；全文不出现 `execute()`／`rt.loopOpt`／`subagent_197`。
- `cmd/wisp/tasklog_subloop_378_test.go:138` 自己造 base options 并**手工塞** `Journal: &taskLifecycleJournal{store: store}`；它的头注释`:13-14` 自己写明这是 "an Options snapshot **shaped exactly like** the one execute() publishes"。

=> 把 `run.go:1583` 摘掉（回到修复前：`rt.loopOpt.Journal == nil`），两枚用例**照绿**。AC#2 明写的"每处改动配一枚今天响、不修不响的用例"与"反向判据＝把新加的那一跳摘掉 ⇒ 用例必须转红并逐名报红行"因此都不成立。
独立核：`git grep -n loopOpt -- cmd/wisp` 全仓只有 `config_tiers_233_test.go` 读 `rt.loopOpt`，且读的是 `Config.PerToolTimeout`（`:192`），不碰 `Journal`。
并且 `378-a2-evidence.md` 里没有任何变异／红行／还原核记录——只有三条**基线绿**（`ok 0.339s/0.353s/0.458s`），那是同一用例重跑三次，不是"三读数"。

## R2（决定性，验收数没有可复现的用例）

AC#2 的验收判据原文是"**一发 `task.spawn` 跑完之后 `PRAGMA foreign_key_check` 必须 ＝ 0**（不许只测'多了一行 `task_log`'）"。树里没有一枚用例同时做这两件事：

- P1（`child_task_log_378_test.go:45-49`）用 `j.StartTaskLog` **手工造**父行：没有 `task.spawn`、没有子环路。它证的是"给了父行就是 0"，不是"子环路会给父行"。
- 唯一真跑 `task.spawn` 的 `tasklog_subloop_378_test.go` 里，bridge 不带 journal（`:167-170` 无 `Journal` 字段），子任务的 `tool_call` 行按构造为 0（`[D]` 在 `:248` 断言增量为 0），所以那里的 fkc 恒 0 且**空洞**。
- 曾经同时做过这两件事的用例 `cmd/wisp/tasklog_fkcheck_378_test.go`（572 行，真 `wisp run` ＋ 真 `task.spawn` ＋ 真 SQLite ＋ fkc）**已从树里和盘上移除**（台账 `7b54388c4`／`596f34977`；`378-a2-evidence.md:44` 明写"不可信，勿救"）。它的读数只留在 `docs/evidence/s1/378-r1-impl.md`（**untracked**，`git status` 可核）里。
=> `378-a2-evidence.md:17` 那句"本测不是'只测多一行 `task_log`'——它断言 `foreign_key_check=0`，符合 AC#2 验收判据的后半句"是**断言了数、没有场景**：数对了，场景（spawn）不在。
=> 后果：本票的定义性读数（"子任务的工具调用行有无父"）在今天的树里**没有任何可复现的尺**；`378-a0.md:12-14` 与 `:31-33` 自记的 ②④【待实跑】在可入库的意义上仍未收口。

## R3（未裁／未登记的用户可见假终态）：子任务的 `task_log` 行永远不关

- 开行：`internal/agent/loop.go:414` -> `internal/agent/journal.go:129`，`State: "running"`。
- 关行：`internal/agent/loop.go:1057` -> `journal.go:144` -> `taskLifecycleJournal.FinishTaskLog` **故意惰性**（`run.go:1326-1328`）。全仓只有 `run.go:1742` 关行，用的是 **root 的** `res.TaskID`；`task.spawn`／bridge 不写 `task_log`（`internal/tools/subagent_197.go` 全文 `Journal` 零命中）。
- 于是：隐私页"任务日志"里每个派生过的子任务恒显示 running（`378-r1-impl.md:42` 的读数正是 `state=running`，`:98` 自认"没有终态写手"）；更糟的是 `cmd/wisp/crash_recovery.go:150-159` 会把**任何**仍为 running 的行改成 `interrupted` => 一次非净退出之后，**成功完成**的子任务在页面上变成"中断"。
- 这是选〔甲〕之后新增的用户可见假状态，AC#2 的用例只数行数、a0 ③只写"会出现一堆子任务"，都没量到；硬提醒#3/#4 要的"代价面"没量到这一条（retention 按 `started_at` 老化，`internal/memory/retention.go:159/163`，那半没问题）。
- 处理方式属产品／契约裁形（谁关行：`task.spawn` 在子任务收口时补一发 `store.FinishTaskLog`，或给 `taskLifecycleJournal` 一枚可关的半）。**要么在 AC#1 里补一句再落地，要么另开票并按 SPEC-12 §5 五字段登记推迟**——v1 当时两件都没做。

## R4（具名前置红，不属本票，但它封死了"跑级"读数）

`internal/tools/system.go:469-471` 的 `screenCaptureSchema` 仍是坏 JSON：`:470` 的属性串以 `"}` 收尾后**直接**拼 `windowTargetProperties`（`:471`），21 枚注册工具里恰这一枚 Parameters 非法 => go 的 json v2 校验 `json.RawMessage`，每发请求在 encode 阶段就死。本腿现跑 `Test233HotTierAgentReachesTheNextTask` 得的就是这一发（见上）。
影响判定：**任何"真 `wisp run` 跑一发"的用例今天都不可能绿**，包括 AC#0 ② 与 AC#2 的字面验收。`378-r1-impl.md` §5 曾以 +1/-1 修过它，但随该腿回退（`f63688cef` "越界 system.go 已回退"）。=> 这条得由编排者单独持有（票 23 的面），本票不管；但在它落地之前，AC#2 的字面判据只能靠**装配级**（`task.spawn`＋真桥＋真 Store，不经 `wisp run`）来量。

> v2 补记（2026-11-04）：该红已由票 404 修（`system.go:470` 补逗号＋`Test404ScreenCaptureSchemaIsValidJSON`）；本轮 `Test233HotTierAgentReachesTheNextTask` 现跑 `ok 2.831s`，跑级通道已解锁。

## 本腿独立核过的接线（修复本身在代码上成立，方向对）

退回的不是"码错"。逐跳核过：

1. `run.go:1008` `BaseOptions: func() (agent.Options, bool) { return rt.loopOpt, rt.loopOptSet }`——闭包读时取值，故 `:1583` 的再发布对**已经在册**的 `task.spawn` 生效。
2. `run.go:1547` 发布 base（`Journal: nil`，`run.go:1521`）-> `run.go:1575-1584` 在 `rt.store != nil` 时挂 `taskLifecycleJournal` 并**再发布** `rt.loopOpt`。
3. `internal/tools/subagent_197.go:306` `base, wired := t.d.BaseOptions()`，`:320` `opt := base`；全文件无 `Journal` 赋值 => 子环路拿到的正是同一枚 `taskLifecycleJournal`。
4. `internal/agent/loop.go:393` `newTaskJournal(l.opt.Journal, taskID, taskID)` -> `:414` `j.startTask(...)`（在 `AdmitTask` 之后、首个模型调用之前）=> 父行先于任何 `tool_call` 行。
5. 生产 bridge 带 journal：`run.go:955-960` `Journal: mem`（`rt.store = mem`，`run.go:532`），booking 用 `req.TaskID`（`internal/tools/bridge.go:1266`），子环路派发时即子任务 id => 父先行成立，fkc 应为 0。

## 若要收口的最小补救（供编排者派单；本腿不实现）

1. **给既有装配级用例补上验收数**（今天就能做，不依赖 R4）：`tasklog_subloop_378_test.go` 的 bridge 按生产样子给 journal（`run.go:960` 的 `Journal: mem`），子任务调一次工具后**用独立新连接跑 `PRAGMA foreign_key_check` 断言 0**（尺可照 `child_task_log_378_test.go:109-128` 的 `fkCount378`）。这枚同时满足"一发 `task.spawn` 之后 fkc=0"，且把 payload 撤掉即转红（子行的 `tool_call` 变孤儿 => fkc>0）。
2. **给那一跳补牙**：R4 落地后，用 `config_tiers_233_test.go:120-195` 的同一套 run 级手法，真跑一发后断言 `live.loopOpt.Journal` 非 nil 且为 `*taskLifecycleJournal`——摘掉 `run.go:1583` 即转红。R4 未落地时，退而用本仓已有的源码序尺（`costwiring_44r3_test.go:112`／`quota_44r5_test.go:200` 的 `strings.Index(src, ...)` 形状）把它钉住，并具名"这是接线尺，不是行为尺"。
3. **R3 裁形**：在 AC#1 补一句"子任务行由谁关、终态是什么"，并配一枚断言"子任务行终止为完成态、不是 running／interrupted"的用例；或另开票＋按 SPEC-12 §5 五字段登记推迟。
4. **AC#0 收口**：把 `378-r1-impl.md` 入库（否则 ②④ 的读数没有 receipt），并把 ① 的归因改到 `e76712fb8`（"修复已在树里，缺的是尺"），②④ 的正控（root fkc=0）一并报数。
5. 反向判据要按 AC#2 原话留证：变异 -> 逐名红行 -> 还原核 sha256 在同一串、`git diff` 空，且还原时机在最后一发变异之后。

## 本腿没做的事

不翻任何框、不改票名／`Status:`／`-done`、不碰 373／368／366 任一原句、不碰 `run.go`／thresholds／golden／frontend／design、不 push、未写台账 R 条（`R##` 编号归编排者）。本件是唯一提交物。

> **v2 对账（2026-11-04）**：v1 的「最小补救」第 1／2 条已由 `cmd/wisp/tasklog_fkcheck_378_test.go` 与第 3 节的反向读数落实（同一枚 run 级用例既量 fkc=0 又对 `run.go:1583` 转红，且带还原核）；第 3 条＝票 405（tracked、open）；第 4 条的**实质**（① 归因更正＋②④ 读数入库）由本件第 5 节与新用例达成（`378-r1-impl.md` 仍 untracked，但读数不再由它承载）；第 5 条见第 1／3 节。v1 的结论「修复本身成立」在本轮被反向读数直接证实。
