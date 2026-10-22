# 票 179 落地腿 `179-r1` — 交付件（实现者）

- 身份：写腿／实现者。**本件只登记读数与判据**，⛔ 未翻任何框、⛔ 未改票名／`Status:`／`-done`、⛔ 未改台账、⛔ 未 push。
- 锚：2026-10-22 22:2x–23:0x +08。**共享工作树里有别的腿在连续提交**，HEAD 在本程内漂了 6 次（`163ae626` → `811054e1`）⇒ ★**每一发读数都带它自己的 HEAD**，下面逐发登记。
- 写面：`internal/agent/**` ＋ 本件。**本程对 `internal/agent/**` 的产码改动＝0 行**（理由与证据见 §3）。⛔ 未动 `cmd/**`、`internal/risk/**`、`internal/tools/**`、`internal/session/**`、`internal/memory/**`、`docs/specs/**`、`docs/PLAN.md`、`scripts/**`、`.github/**`、`frontend/**`、`design/**`；⛔ `thresholds.go`／golden／`allowlist.txt`／三枚冻结件；`probes/**` **只读**（新台件全在 `D:\tmp\wisp\20261022\179r1\**`）。
- **一句话结论**：★**`AC#8` 的读数今天在 HEAD 上已经拿得到、而且全绿**；买断它的两跳**都不是本程落的**——`loop.go:856-859` 的 `case RiskL0:`（`ec520750`，票 179 前几程）＋`internal/risk` 的宿主指针豁免（票 183/185）。`AC#1` 的拒绝形状我用 `-overlay` 在**当前 HEAD 的树上**复现到了，逐字在 §1。⇒ **没有任何"要修的东西"留给我这一腿**；§5 说明为什么我没有拿一枚今天绿的用例去凑"新用例"。

## §1 `AC#1`——只登记事实：修前那一发的逐字拒绝行＋尾态

**台件替换（必须先说清）**：`probes/176/r1/zz176r1_e2e_test.go` **今天编不过**，两条独立原因：
① 它调 `rt.execute("…")`，现签名是 `func (rt *agentRuntime) execute(parent context.Context, task string) int`（`cmd/wisp/run.go:1067`）⇒ `go test -overlay=probes/179/r2/overlay-e2e.json … -run TestReread…` 逐字 `zz179r2_e2e_test.go:349:22: not enough arguments in call to rt.execute / have (string) / want (context.Context, string)`（`logs/01-head-as-is.txt`）；
② 它用 `contentText`，`cmd/wisp` 今天**没有**这个符号（尺：`git grep -n "func contentText" -- cmd/wisp` ⇒ rc=1、0 命中）。
⇒ 按票面 `AC#1` 的补法所指，改用**同一发端到端的后继台件** `probes/179/r2/zz179r2_e2e_test.go`（同一 payload 20000 字节、同一条 `task.output → 指针 → fs.read` 路径），**复制到 `D:\tmp\…\zz179r1_e2e_test.go` 后只做 4 处 rig 改动**：`execute(ctx,…)` 签名、`"context"` 导入、删掉用不上的 `"runtime"` 导入、读数落点改 `WISP179R1_LOGDIR`；外加今天配置模式要求的 `enabled = true`（否则 `wisp run` 退 2 并在 `stderr` 说 `model "m1" … is disabled (enabled=false)`，`logs/02-head-probe.txt`）。★**这 4＋1 处都不碰门**（仓里 `probes/**` 一字节未改：尺 `git status --porcelain -- .scratch/wisp/probes/179` 只列出 `179-r3` 那批**未跟踪**日志，`LastWriteTime` 全是 **2026-10-14**，早于本程；`176/r1` 同理）。

**发**（`logs/20-e2e-prefix-historical.*`，HEAD **`b4102f1e`**）：`-overlay` 只把 `internal/agent/loop.go` 换成**修法之前**那一版（`git cat-file blob 039efb47:internal/agent/loop.go`，尺：`git rev-parse --short ec520750^` ⇒ **`039efb47`**；该 blob 里 `grep -n "case RiskL0"` ⇒ **0 命中**、`风险未分级` 在 `:785`），`go build ./internal/agent/` 在该 overlay 下 **rc=0**（历史文件与今天的树兼容，票面补法可用）。

**逐字拒绝行**（模型收到的 tool 行原文，读数件 `:18`）：

```
tool error: 风险未分级且直通开关关闭，已拒绝执行
```

**逐字尾态**（读数件 `:2`、`:6`，console `:26`/`:28`/`:29`/`:30`）：

```
run A 退出 0 / run B 退出 0 / 确认卡 0 -> 0
模型调用序列: task.output 5912c0c3-b7ed-474f-a2a5-92b2651c4405
分岔「风险未分级」出现次数: console=0 rig-turns=0 tool-rows=1
wisp run: 任务 5912c0c3-b7ed-474f-a2a5-92b2651c4405 结束（completed，1 轮，0 次工具调用，成本 0 CNY（未计价））
[工具 task.output -> error]
失败：task.output 没给出可续读的指针
wisp run: 任务 6728b338-b54a-4178-930c-297041c66c60 结束（completed，2 轮，1 次工具调用，成本 0 CNY（未计价））
```

⇒ **后半截在这一发上根本没有开始**：`task.output` 被自家门拒 ⇒ 桩里没有「全文见」指针 ⇒ `fs.read` 一发都没发出（`整份` 那一行 0 字节、`有界` 那一行 0 字节）。探针自报 5 处红：`:366` `整份 fs.read 那一发的 tool 行没回来：续读没走通`、`:369` `有界 fs.read 那一发的 tool 行没回来`、`:377` `真机 CLI 的续读那一发上出现「风险未分级」`、`:390` `模型这三发应为 task.output 然后两发 fs.read，got [task.output …]`、`:393` `模型续读用的路径 "" 不是名册登记的那条 …`。★**这一格只登记事实**；修完之后它当然不翻绿。

## §2 `AC#8`——修前／修后两发端到端读数（前半截＋后半截都要）

| # | 发 | HEAD | 探针 rc | 前半截（`task.output`） | 后半截（`fs.read` 续读） |
|---|---|---|---|---|---|
| 20 | **修前**（overlay 换回 `039efb47` 的 `loop.go`） | `b4102f1e` | **1** | ⛔ 逐字拒绝（§1），tool 行 0 字节 | ⛔ 一发未发，两个 tool 行 0 字节 |
| 21 | **变异**（只摘今天 `loop.go:856-859` 那 4 行） | `e94a5b17` | **1** | ⛔ 同名拒绝，`tool-rows=1` | ⛔ 一发未发；同样 5 处断言红 |
| 03 | **修后／原样** | `d0194273` | **0** | ★`success`，`risk=L0 decision=allow rules_hit=[]` | ★两发 `success`，逐字节相等 |
| 50 | **修后／终态复量** | `811054e1` | **0** | 同上 | 同上，★sha256 与 03 **逐字相同** |

**修后那一发（`logs/50-final-green.*`）的逐字读数**：

```
run A 退出 0 / run B 退出 0 / 确认卡 0 -> 0
模型调用序列: task.output 8cfa84b2-… | fs.read …\artifacts\tool-output-agent-task-8cfa84b2-….txt（整份，不带 max_bytes） | fs.read …（max_bytes=10000，整段留在上下文里）
模型续读用的路径: C:\…\001\artifacts\tool-output-agent-task-8cfa84b2-….txt
整份那一发：宿主宣告总长 20000 字节；宿主为那一发新落的产物 …\artifacts\tool-output-call-179r2-full.txt（20000 字节）；与指针所指文件逐字节相等＝true
分岔「风险未分级」出现次数: console=0 rig-turns=0 tool-rows=0
分岔「L2 级」出现次数: console=0 rig-turns=0 tool-rows=0
分岔「R4」出现次数: console=0 rig-turns=0 tool-rows=0
分岔「查不到这个任务」出现次数: console=0 rig-turns=0 tool-rows=0
call kind=success tool=task.output risk=L0 decision=allow outcome=success rules_hit=[] in_allowlist_scope=unknown grant_id=0 refusal=none reason="无规则命中（L0 直接执行）"
call kind=success tool=fs.read    risk=L0 decision=allow outcome=success rules_hit=[] in_allowlist_scope=true    grant_id=0 refusal=none reason="无规则命中（L0 直接执行）"   ×2
[工具 task.output -> success] / [工具 fs.read -> success] / [工具 fs.read -> success] / 续读完成
```

★**逐字节**：整份那一发 `20000 == 20000`，sha256 **`e1d2b92fc204d798dfaf9b2ba0faa30a92259bfa88cecad8bbf5261a943647d8`**；有界那一发 `10000 字节`，sha256 **`9db96549cd3bd53aad5dbffa00408a9f2d92044bd1748cdb9c02b4cef880d85e`**（两发都是宿主替模型算的同一份字节，模型手里那一行就是这 10000 字节）。★**确认卡 0 枚、`已拒绝执行` 0 次**（尺：`grep -c "已拒绝执行"` ⇒ 0；`grep -c "\[确认"` ⇒ 0）。
★**这一发同时就是三枚条件格的证据**：票 177 `AC#3`「端到端条件」＝四枚分岔全 0；票 175 `AC#5`＝指针就是名册登记的那条、`fs.read` 拿它不被拒；票 176 `AC#3`/`AC#4`＝产物真落盘 20000 字节、宿主自己写下的指针被模型逐字用上并读回。

## §3 根因（★本程自己读码定位，⛔ 抄票面）与「为什么本程零改动」

- 判定那一支现在读作（`internal/agent/loop.go:847-870`，尺：`grep -n "case RiskL0\|风险未分级" internal/agent/loop.go` ⇒ `850:` L1/L2、**`856:` `case RiskL0:`**、`865:` 那句拒绝）：`risk` 先撞 `RiskL1, RiskL2`，再撞 **`RiskL0`（直接放行＋记 `allow`）**，其余才落 `default:` 的 fail-closed。★**"声明为 L0"今天已经不再落进 `default:`**——这正是票面第 1 把尺当年那个洞，`ec520750`（`git log --oneline -1 ec520750` ⇒ `票 179 修法＋正反两形常驻判据：decideRisk 给声明 L0 单开一支，"" 仍 fail-closed`；`git show --name-only` ⇒ `internal/agent/loop.go`＋`internal/agent/declared_l0_risk_179_test.go`）已经把它补上。
- 09-28 的**第二位阻断者**（`R4: 包含来自 task.output 的内容` 升 L2，`179-r2` 的读数 console 里 8 枚）今天也不在了：`internal/risk/provenance.go` 的 177 形状 A＋183 值规则＋185 宿主铸造名册都在（`provenance.go:508-557`、`:611-623`），迹是上面「分岔 R4 出现次数 **0/0/0**」。
- ⇒ ★**`AC#8` 今天没有剩余缺陷需要我修**；再"修"只能去动上面两处已经定案的东西（`internal/risk/**` 与 `internal/tools/**` 更是明令**别家写面**）。**本程对 `internal/agent/**` 的改动＝0 行**，`git status --porcelain -- internal/agent` **空**，`internal/agent/loop.go` 工作树 sha256 与 `HEAD` blob **相同**（§4 还原核）。

## §4 反向判据——三读数（变异红／还原绿／原样绿）＋还原核

**为什么用 `-overlay`／一次树内变异**：共享工作树有 ≥2 条腿在飞（本程 20 分钟内 HEAD 漂 6 次），长时间把 `loop.go` 挖掉会把别人的整包门禁打红 ⇒ 端到端两发（40s 级）一律 `-overlay`（被变异的是**内存里同一份字节**，树不动）；单测那一发做**树内**变异并当场还原（窗口 ~15s，`trap` 兜底 + 仓外备份）。

| 读数 | 尺 | 结果 |
|---|---|---|
| **① 变异红（单元，树内）** | 删 `loop.go:856-859` 那 4 行 ⇒ `grep -n` 里 `case RiskL0` **消失**、`风险未分级` 从 `865` 移到 `861`；`go build ./...` **rc=0** | `-run Declared -v` ⇒ **rc=1**，`PASS=0 FAIL=2 RUN=6`：`--- FAIL: TestDeclaredL0RunsWhenPassThroughIsOff179`（红行 `declared_l0_risk_179_test.go:101` `executions = 0, want 1: a declared L0 call must reach the provider`、`:104 decision = "reject", want "allow"`、`:107 outcome = "error", want success`、`:110 error_class = "permission_denied", want none`、`:120 recorded calls = 0, want 1`）＋`--- FAIL: TestDeclaredL1L2RoutingUnchanged179`（`:252` 逐字 `ok/unbooked/why = false/false/"风险未分级且直通开关关闭，已拒绝执行", want true/false/empty`、`:262` `want one row with "allow" booked by the direct L0 arm`） |
| **② 还原绿** | `git cat-file blob HEAD:internal/agent/loop.go > internal/agent/loop.go` | sha256 工作树 ＝ HEAD blob ＝ **`e14f73b1458c411511523ed19e91f7f19a79ce0c863d5abd6358605fc476df72`**（`sha256sum`，⛔ `git hash-object`）；`-run Declared -v` ⇒ **rc=0**、`PASS=2 FAIL=0`；`git status --porcelain -- internal/agent` **空** |
| **③ 原样绿（整包）** | `go test -count=1 -v ./internal/agent/` | **rc=0**、`PASS=102 FAIL=0 SKIP=0 RUN=131`、`ok … 4.069s`；★与开跑前的独立基线（`logs/10-baseline-agent.txt`）**四数逐字相同**（102/0/0/131） |
| **④ 变异红（端到端）** | 同一 4 行摘掉，`-overlay` 编进 CLI 台件 | 探针 **rc=1**，拒绝行与 §1 **逐字相同**，5 处断言红（§2 表 21 行） |

**还原核**：`logs/40-three-readings-unit.txt` 段 ① 的 `sha256 before`＝段 ② 的 `sha256 after restore`＝HEAD blob（同一串 `e14f73b1…76df72`）；`tr -cd "\r" | wc -c` 对 `loop.go` ⇒ **0**（无 CRLF 混入）。

## §5 新用例：**0 枚**（并说明为什么这不是偷工）

本票现有的牙**在单元接缝上已经成对且真咬**：`TestDeclaredL0RunsWhenPassThroughIsOff179`（正向，7 行断言）＋`TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179`（反向 fail-closed 守卫）＋`TestDeclaredL1L2RoutingUnchanged179`（L1/L2 两形不动＋`AC#9` 补的 `decision` 列断言：`:252`/`:262`）。上面 ①④ 就是"把那一跳摘掉 ⇒ 用例转红并逐名报红行"的实测。
★**今天在 `internal/agent/**` 里新写一枚用例，只可能有两种命**：**(a)** 它今天是绿的 ⇒ 按票面 `AC#3`/`AC#4` 自己定的规矩，那是**守卫**，⛔ 不许当本票的新增牙；**(b)** 它今天是红的 ⇒ 说明树里还有缺陷，而 §1–§3 已证没有。⇒ 我**不拿 (a) 充数**。★而 `AC#8` 唯一能钉住它的接缝在 `cmd/wisp`（CLI 组合根），那是**本波敏感面**、派单要求"动前先停手上报"⇒ 我不动，改由 §2 的台件读数承担（台件落在 `D:\tmp`，不是仓内 tracked test）。

## §6 门禁 rc（全部本程现跑；★`cmd/wisp` 整包未当判据）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| build | `GOFLAGS= go build ./...` | **0** | 空输出 |
| vet | `go vet ./internal/agent/` | **0** | 空输出 |
| 单包 | `go test -count=1 ./internal/agent/` | **0** | `PASS=102 FAIL=0 SKIP=0 RUN=131`、`ok 4.764s`（基线）／`4.069s`（终态） |
| d22scan | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`；`ban #8 internal/=619`、`cmd/=148`、`bans #1-5 internal/=247` |
| 长度门 | `sh scripts/check-path-length-budget.sh` | **0** | — |
| **格式甲（判据）** | `git archive HEAD \| tar -x -C D:\tmp\…\export-tree` 后 `gofumpt -l`（**绝对路径** `/d/work/base/gopath/bin/gofumpt.exe`） | **rc=0／列出 20 枚** | ★用 `[\\/]` 判归属：`export-tree/internal/`=**0**、`cmd/`=**0**、`docs/`=**0**、`.scratch/`=**20**。★这 20 枚**按构造**就是在册旧件——甲尺量的是 `git archive HEAD` 的内容，本程唯一的新文件是本件（`.md`，不进 `gofumpt` 的射程）⇒ **本程新增 0** |
| **格式乙（只作诊断）** | 同一把 `gofumpt -l` 打工作树 | 列出 54 行 | 锚到仓根后 `Wisp/internal/`=**2**、`Wisp/cmd/`=**1** ⇒ 三枚都是 **CRLF 假红**：`git ls-files --eol` ⇒ `i/lf w/crlf`、`git diff --stat` **空**、`tr -cd "\r" \| wc -c` ⇒ `internal/agent/tools.go` **225**／`internal/agent/approval/pending_read.go` **131**／`cmd/wisp/models.go` **389**。★**非本程引入、本程未动**（工作树对这三枚 `git diff` 为 0 字节）；`internal/agent/loop.go` CR=**0** |

## §7 报回的读数（⛔ 本程不修，写面在外）

1. ★**`MODE-REDLINE … still asks` 是一句假话，每发 e2e 出现 3 次**（`logs/50-final-green-readings/e2e-readings.txt`，逐字 `MODE-REDLINE mode=ask_every_step still asks tool=fs.read rules=[] because=无法识别的风险级，fail-closed 保留询问`；3 枚属于 `task.output`×1＋`fs.read`×2）。机制：`internal/tools/bridge.go:311` 把 **L0 判定**交给 `mode.Screen`，而 `internal/risk/mode.go:163-165` 只认 `L1/L2/Deny`，其余一律回 `Kept="无法识别的风险级，fail-closed 保留询问"`；`bridge.go:324` 据此打印 "still asks"，但 `bridge.go:352` 的 `route` 对 L0 直接放行 ⇒ **一行说"还要问"，一行没人问、卡 0 枚**。⇒ 这是"日志里编造出一个不存在的动作"那一族（票 97 已定过案）。**两处都在我的写面之外**（`internal/risk/**`、`internal/tools/**`），⇒ 只登记，请编排者裁是否开票。
2. ★**`AC#1` 的台件今天编不过**（§1 两条独立原因）⇒ 票面 `AC#1` 留给下一程的"补法"应当改写成"用 `probes/179/r2` 的台件＋`D:\tmp` 适配副本"，⛔ 不要再照"复跑 176/r1 台件"字面执行。
3. ★**`logs/20-*` 的尾态与 09-28 `176-r1` 的尾态不同**：那一发是"7 次拒绝后判卡住"，我这一发是"拒绝 1 次→桩主动收尾→两任务都 completed、退出码 0"。⇒ 拒绝行同形，**尾态不同**，⛔ 不要拿我这一发的尾态去替换旧读数。
4. ★**HEAD 漂移**：本程 6 次（`163ae626`→`b4dfe165`→`d0194273`→`b4102f1e`→`e94a5b17`→`4c771f48`→`811054e1`），全部是**别的腿**的提交；`git status --porcelain -- internal/agent` 在本程首、中、尾三次均为**空**。

## §8 本程**没有**做的事

⛔ 未翻框／未改票面（含 Progress log）／未改台账／未 push。⛔ 未动 `cmd/**`（因此 `AC#8` 的证据只能是台件读数，不是 tracked test）。⛔ 未开新票（`MODE-REDLINE` 那枚归编排者裁）。⛔ 未替票 177/175/176 翻它们的框——§2 只是**它们要的那一发读数**，勾不勾由各自裁决程判。

## §9 提交后复量（★最后一次读数在最后一枚提交之后取）

- 本件的提交＝`32c8eeac`：`git show --name-only` ⇒ **1 枚**：`docs/evidence/s1/179-r1-impl.md`。★**本件自己的提交号不写进本件**（防自我推翻，票 160 `02b3882c` 的先例）；本节这枚补写落在**下一笔**。
- **提交后**（HEAD `32c8eeac`）五道门重跑：`go build ./...` **rc=0**／`go vet ./internal/agent/` **rc=0**／`sh scripts/d22scan.sh` **rc=0**／`sh scripts/check-path-length-budget.sh` **rc=0**／gofumpt 甲尺 **rc=0**；`go test -count=1 ./internal/agent/` ⇒ **rc=0、`PASS=102 FAIL=0 SKIP=0 RUN=131`、`ok … 9.330s`**。
- 片刻后另一条腿把 HEAD 推到 `394cc212`，我在**同一棵树＋一枚别人的提交**上再取一次细读数（`logs/62-postcommit-gates-detail.txt`）：d22scan **rc=0**（`clean - no D22 ban violations`；`ban #8 internal/=620`、`cmd/=148`）、pathlen **rc=0**、gofumpt 甲尺 **rc=0** 且 `total=20 / internal=0 / cmd=0 / docs=0 / .scratch=20`。
- `.go`／越界尺（★**正控＝同尺打一笔含 `.go` 的真提交**）：`32c8eeac`（本件）⇒ 1 枚、**`.go`=0**、`internal/`+`cmd/`=**0**、写面内；`ec520750`（**正控**，票 179 修法那一笔，非本程）⇒ **2 枚 `.go`**（`internal/agent/loop.go`、`internal/agent/declared_l0_risk_179_test.go`）⇒ 同尺在真 `.go` 提交上**报得出来**，不是瞎尺。
- 终态：`git status --porcelain -- internal/agent` **空**；`internal/agent/loop.go` 的工作树与 `HEAD` blob sha256 同为 **`e14f73b1458c411511523ed19e91f7f19a79ce0c863d5abd6358605fc476df72`**。
