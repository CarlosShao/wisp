# 票 376 终裁腿 `376-v1`（非实现者）：〔乙′〕按 level 选源合格，`AC#2` 能翻

- 腿 `376-v1`（非实现者裁腿；**本笔是重派**，上一枚同号腿中途失败）。读数时 HEAD `254c05108`（2026-10-24 15:1x +08）。工作树不落任何 `.go`；全部探针/变异在仓外。
- **记录一处与派单不符**：派单说上一枚同号腿“没留任何残留”，但 `D:\tmp\wisp\20261024\376v1` **已存在且非空**（上一枚的 `pin-*`／变异体／读数约 50 件）⇒ 按纪律改用下一个空号 **`376v2`**；本件全部读数出自 `376v2`，`376v1` 未读未用。
- **写面**：本件 ＋ `D:\tmp\wisp\20261024\376v2\**`。375-r1 正在写的六个共享文件（`cmd/wisp/{run.go,resident_approval_windows.go,resident_task_source_windows.go}`＋`internal/agent/approval/{gate.go,reach.go,reach_375_test.go}`）**只读、⛔ 写**。
- **读数法**：被测字节一律 `git show` 导出 blob，经 `-overlay` 换进编译（工作树零字节改动）。`run-fix.go` 的 `git hash-object` = `ee7ef09afa…` = `ceef4d069:cmd/wisp/run.go`；`test-fix.go` = `f52b61ad8c…` = 提交里的测试件；`run-base.go`（预修复）= `b670c7924a…`。

## ① 票面前提新鲜度（对 HEAD blob 复量，含现跑行号）

- `git show HEAD:internal/agent/approval/gate.go | sed -n '534p'` ⇒ `p := g.promptFor(d, corr, 0, g.q.Timeout(), nil, nonce)` —— 第 3 参字面 `0`／第 4 参 `Timeout()` 的分叉**仍在**。
- `gate.go:597` 签名、`:609-610`（`Window: win`／`Deadline: deadline`）、`:293`（L1 那路 `win=g.window, deadline=0`）逐字未变；`replies.go:83`「Window is the L1 block length, zero on an L2 card」**仍为真**（我实测 L2 的 `P_WINDOW_RAW=0s`）。
- `git show HEAD:cmd/wisp/run.go | sed -n '1657,1661p'` ⇒ 修后的五行确实在 **1657–1661**；票面写的 `run.go:1563` **今天已漂**（该行现为 `changed = true`，属 373-r1 的改写）⇒ 引用必须改成 1657–1661。
- 生产读点穷举（HEAD，`internal/**`＋`cmd/**`，⛔ 含 `_test.go`）：`p.Window` 只被 `run.go:1657`（印发点）与 `replies.go:181`（逐字拷进 `ReplyCard`）读；`p.Deadline` 只被 `run.go:1659`（印发点）与 `replies.go:182`（同形死拷贝）读 ⇒ **唯一决定显示的生产读者**分别是 1657／1659；本修之前 `p.Deadline` 连印发点都没有（普查 #9：`ReplyCard.Window/.Deadline` 0 读者）。其余 `.Window` 命中是 `Options.Window`／`Gate.Window()`／`GateFuncs.Window`／`WindowTitle` 等同名不同义。
- 判：**语义前提全部成立**；唯一要改的是**行号引用**。`gate.go` blob `78e2d28a8…` 与普查基线逐字节相同；`run.go` blob 从普查的 `7b9c285db…` 变为 `ee7ef09af…`（373-r1 的改写＋376 的 5/1）⇒ 行号不能沿用票面。

## ② `L2` 卡印秒数 == `Queue().Timeout()`（我自己造的两发）

探针 `zz376v2_probe_test.go`（仓外，真 `consoleApprovalUI`＋真 `approval.New`，`-overlay` 注入 `cmd/wisp`；原文 `probe-run.txt`）：

| 发 | 卡面逐字 | `Queue().Timeout()` | 实测等待 | 判 |
|---|---|---|---|---|
| L2 带底 `31`（真钟）| `窗口 31.0s` | `31s`（`QUEUE_TIMEOUT_RAW=31s`）| **`31.001s`** | **EQUAL=true** |
| L2 默认 `300`（假钟）| `窗口 300.0s` | `5m0s` | 假钟手打 | **EQUAL=true** |

- 同发前提：`P_WINDOW_RAW=0s`、`P_DEADLINE_RAW=31s/5m0s` ⇒ 读的是 `Deadline` 那一源，⛔ 不是“碰巧非 0.0”。
- 反控：预修复 `run.go` 下，提交里的 L2 用例从真实控制台读回 `窗口 0.0s`（`read-prefix2.txt:248` 红句）；同一次 L1 仍 PASS。

## ③ `L1` 不被打成 `0.0s`（我自己造的一发）

- 配置 `Window: 1s` → `gate.go:143-151` 夹到 `MinL1Window=2s`；卡印 `窗口 2.0s`，`Gate.Window()=2s`，实测 **`2.004s`**，`P_WINDOW_RAW=2s`、`P_DEADLINE_RAW=0s` ⇒ **EQUAL=true、L1_ZERO=false**。
- 这一格正是〔乙〕字面换源会破的那格（L1 的 `Deadline` 是 0）。

## ④ `gate.go:534` 一字未动 ＋ `internal/agent/approval/**` 零字节

- `git show --numstat --format= ceef4d069` **只有三枚**（见⑨）；`git diff --numstat ceef4d069^..ceef4d069 -- internal/agent/approval` **空**；`git diff --numstat ceef4d069..HEAD -- internal/agent/approval` **空**。
- `gate.go` blob 在 `ceef4d069^`／`ceef4d069`／`HEAD` 三处同为 `78e2d28a8…` ⇒ 该包逐字节未动（不只是 534 行）。
- ⚠ 工作树里 `internal/agent/approval/gate.go` 有 **375-r1 的未提交改动**（mtime 13:41）⇒ 零字节判据只能在**提交边界**量；我的所有测试用 `-overlay` 把 `gate.go` 钉回 `HEAD`、并删掉 375 的 `reach.go`／`reach_375_test.go`，读数不受影响。

## ⑤ 两枚牙的核

- `Test376R1L2CardPrintsQueueTimeout`：期望值= `timeout := g.Queue().Timeout()`（测试件 `:238`），比较 `want := timeout.Seconds()`（`:247`）；全文无 `want := 300` 类数字字面量（`grep` 无命中）。从真实控制台输出用 `窗口 ([0-9]+\.[0-9])s` 读回卡上那个数。
- `Test376R1L1CardKeepsItsBlockLength`：期望值= `block := g.Window()`（`:276`），并断 `block == approval.MinL1Window`、`printed != 0`。
- 两枚都钉生产者形状：L2 先断 `p.Window==0` 且 `p.Deadline==Queue().Timeout()`；L1 先断 `p.Window==Gate.Window()` 且 `p.Deadline==0`；两枚都用 `approval.Clock` 假钟，300s 不塞进 CI。
- **恒真变异**：`alwaysl2.go` 把 `if p.Level == "L2"` 改成 `if true`（等于〔乙〕的打印行为）⇒ `read-always.txt`：L2 **PASS**、L1 **FAIL**，红行 **`:307`**（`the L1 card printed "窗口 0.0s"`）＋`:312`（zero window）；rc=1。⇒ 这枚牙**真的抓得住乙形**，⛔ 是摆设。

## ⑥ 预修复态实测红（“HEAD”的射程）

- ⚠ 今天的 `HEAD` **已包含修复**（`HEAD:run.go` blob == `ceef4d069:run.go` == `ee7ef09af…`），所以“放回 HEAD 的 run.go”必须指**预修复态**：`ceef4d069^:cmd/wisp/run.go`（blob `b670c7924a…`）；其 sha256 与 376-r1 的仓外 `head.go` **逐位相同**（`11406974…`）。
- `overlay-prefix`（只把 `run.go` 换成预修复态、测试件用提交里的）⇒ `read-prefix2.txt`：**rc=1**；`--- FAIL: Test376R1L2CardPrintsQueueTimeout`，红行 **`approval_card_window_376_test.go:248`**（`the L2 card printed "窗口 0.0s", but the gate waits Queue().Timeout() = 300.0s (5m0s)`）；同一次 L1 **PASS** ⇒ 归因精确。

## ⑦ 反向判据（我自己造的）＋还原核

- 摘跳变异 `cut.go`：从修后文件**只删** `if p.Level == "L2" {`／`shownWindow = p.Deadline`／`}` 三行（`git diff --no-index --numstat` = `0/3`），保留 `shownWindow := p.Window` ⇒ `read-cut.txt`：**rc=1**、L2 **FAIL `:248`**、同一次 L1 **PASS**。
- 还原绿：`read-fixed2.txt`（提交字节）⇒ **rc=0**、两枚 PASS。
- **还原核（时机＝最后一发变异 `prefix` 之后）**：仓外 `run-fix.go` sha256 **`565ff346…`** 与变异前**逐位相同**（`git hash-object` = 提交 blob `ee7ef09af…`）；`test-fix.go` sha256 **`9b56a05a…`** 未变（blob `f52b61ad8…`）；**376 的 `git diff --numstat ceef4d069^..ceef4d069 -- cmd/wisp/run.go` = `5/1`**。
- ⚠ 工作树 `cmd/wisp/run.go` 的 sha256 在我读数期间由 `71d8d231…` 变为 `e6c37997…`、numstat `73/1` ⇒ **这是 375-r1 的在飞改动，⛔ 是 376 的 5/1**；本笔从未写工作树，故还原核量的是仓外提交字节＋提交边界。

## ⑧ 四数复跑（两包分开、⛔ 相减）

- `cmd/wisp` 定向（同一条 12 分支 `-run`，带 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`、`GOFLAGS=`、`-overlay` 钉字节）：
  - 改前（`overlay-before-prefix`：`run.go`＝`ceef4d069^` 的**预修复字节**、新测试件删出编译）⇒ **PASS 13 / FAIL 0 / SKIP 0，rc=0，`ok … 101.601s`**；另一发把 `run.go` 换成提交字节、同样删测试件也量到 13/0/0 ⇒ +2 **只**来自新测试件。
  - 改后（`overlay-fixed`：提交的 `run.go`＋提交的测试件）⇒ 首量 **PASS 14 / FAIL 1 / SKIP 0，rc=1**，唯一红= `TestTicket223HandEditedFsLooseningCostsAnL2Card`（`config_reload_223_test.go:439`）；**隔离该枚 `-count=2` ⇒ 2/2 PASS**，**同一 12 分支在最终 pin 集下复量 ⇒ PASS 15 / FAIL 0 / SKIP 0，rc=0，`ok … 113.334s`**（此前一发同判 15/0/0、119.976s）。
  - 首量那枚红是**票 223 的已知 stdout 单读 flake**（`docs/evidence/s1/198-firstrun-config-v1.md:59` 记同形：整包一红、隔离 `-count=2` 全绿；根因归 223）⇒ ⛔ 与 376 相关。
  - 差 2 的**具名证据**（⛔ 只靠相减）：before/after 顶层 PASS 名字集合逐名 `comm`，**只有** `Test376R1L2CardPrintsQueueTimeout` 与 `Test376R1L1CardKeepsItsBlockLength` 两枚只在 after；子用例两列各 2 枚。
- `internal/agent/approval` 整包（该包零字节）：**两遍均 PASS 97 / FAIL 0 / SKIP 1，rc=0**（SKIP=`TestDefaultDeadlineWallClockMeasurement`；两遍独立跑，⛔ 用相减推）。

## ⑨ 边界（三枚文件）

`git show --numstat --format= ceef4d069` ⇒ **仅** `320/0 cmd/wisp/approval_card_window_376_test.go`、`5/1 cmd/wisp/run.go`、`118/0 docs/evidence/s1/376-r1-impl.md`（三枚，无 `D`/`A` 残留；`git diff --name-status ceef4d069..HEAD -- cmd/wisp/run.go internal/agent/approval` 为空）。

## ⑩ 门禁 rc

- `go build -overlay=overlay-fixed ./...` **rc=0**（输出 0 字节）；`go vet -overlay=… ./cmd/wisp ./internal/agent/approval` **rc=0**（输出 0 字节）。
- `tools/d22scan`（自成 module，`-root ../..`，当前树）**rc=0、clean**：examined `internal/`=645、`cmd/`=152 Go files；`no D22 ban violations`。
- **格式甲乙两尺**：`gofumpt -l` 与 `gofmt -l` 对两枚提交字节（仓外导出副本）**rc=0、输出空**（⛔ 未用 `gofumpt -l .`；未判 HEAD 是否 CRLF）。

## ⑪ `AC#2` 判语：**能翻**（本腿只裁、⛔ 翻）

- 能翻的三条硬读数：①L2 两发 `printed == Queue().Timeout()`（含真钟 `31.001s`）；②L1 `printed == Gate.Window()==2s`、≠0；③预修复红 `:248`＋摘跳红 `:248`＋恒真红 L1 `:307`，三发逐名归因。
- **缺什么**：**AC#2 本身不缺任何东西**（能翻）；唯一要订正的是**票面的行号引用**（`run.go:1563`→`1657-1661`；`gate.go:534` 未漂）——这属票面文字、**不阻塞** `AC#2` 翻。⛔ 缺实现、⛔ 缺牙、⛔ 缺读数、⛔ 缺边界。真钟 31s 只在本腿仓外探针、⛔ 在 CI（假钟是刻意，避免 5 分钟），这是“测试的射程”而非缺口。
- 纪律：⛔ 翻框、⛔ 改票名/`Status:`/`-done`、⛔ 改台账、⛔ push。

## ⑫ 本件提交＋逐笔越界

- 提交史：第一版 **`2190de8fc`**（87 行，只含本件）；第二版只修正 ①读者措辞／②反控指涉／⑧改前定义并记提交史，其提交号 `git log -1 --format=%h -- docs/evidence/s1/376-v1-verdict.md` 读回（⛔ amend，故不在正文写死哈希）；`git diff --cached --name-status` 已核：只有本件的 `M`/`A`，无 `D`/`A` 残留。
- **375-r1 的脏面在读数期间从派单的 6 个文件扩到 12 个**（新增 6 枚 `cmd/wisp` 测试件，mtime 15:20–15:22；其中 `resident_task_source_246_windows_test.go` 一度让预修复 overlay 编译失败，因为 375 的脏测试引用了它脏 `run.go` 里新加的 `runSpec.answer`）⇒ 我把这 6 枚也按 `git show HEAD:` 钉住后**重跑** before-prefix／prefix／fixed 三发，读数与首轮一致（13/0/0、红 `:248`、绿）⇒ 本件所有 376 读数都来自 HEAD 字节，⛔ 混入 375 的在飞改动。
- 越界：**0 笔** —— 未写 `internal/**`、`cmd/**`、`scripts/**`、`.github/**`、`.scratch/wisp/**`、`docs/reports/**`、`docs/PLAN.md`、`docs/specs/**`；未落任何 `.go`；未碰 375-r1 的**12 个在飞路径**（派单点名的 6＋读数期间新增的 6 枚测试件）；仓外只写 `376v2`。
- 尺的射程（五问）：①对象=`approval.Prompt` 的两个字段与唯一印发点（⛔ Win32 window）；②形状=控制台一行 `窗口 N.Ns` 的浮点数；③数的是生产 `.go`（⛔ `_test.go`）；④生产读者 2 行与字段族总表对上；⑤秒数=`time.Duration.Seconds()` 现算、卡面 `%.1f`。
