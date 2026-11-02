# 票 86 验收裁决 — 86-v1（非实现者，裁决者≠实现者）

- 腿：`86-v1`（验收腿）；被裁落地件＝commit `4d0e2577d`（票 86 落地腿 86-r1）
- 裁决时刻：2026-11-02 18:4x–19:1x +08
- 起手三核：`HEAD = c114b09f7`（＝编排者 A1416 收件提交，`4d0e2577d` 为其父，`git merge-base --is-ancestor` rc=0）；起手 `git status --porcelain` 非 `??` 行数＝0（索引与已跟踪文件干净）；工作目录 `D:\work\workspace\projects plans\Wisp` ✓
- 机器：12 核（`NUMBER_OF_PROCESSORS=12`）；`go1.27.1 windows/amd64`
- 写面：本腿**只新增本证据件**；⛔ 改码／改票面／改台账／push。全部变异/回退实验走 `go test -overlay`，共享树零原地改动。
- ★补账声明：86-r1 漏产 `docs/evidence/s1/86r1-impl.md`，本 verdict 按派单要求**自带 AC 读数作凭据**。

## 环境限制（具名，不影响裁 86 本体）

盘上有**另一条腿（票 44 系）的 WIP**：未跟踪 `internal/observe/{costbudget,costmeter,pricing}.go`、`internal/memory/costdao_44.go`（后含 `internal/agent/{cost,loop}.go` 修改），其中 `pricing.go` 起手时带语法错误，而 `internal/risk` 经 import 链依赖 `internal/observe` ⇒ **现盘 HEAD 状态无法直接编译 `./internal/risk`**。本腿一律用 `-overlay` 把这 4 枚未跟踪 WIP 文件虚拟移除（＝回到 `c114b09f7` 提交态编译；该 4 枚在 HEAD 均不存在，`git ls-tree` 已核）。被裁的 4 枚 internal/risk 文件全程＝HEAD 逐字节（`git diff HEAD --stat -- internal/risk/` 空）。

## ① 形状核（只读逐行）：墙钟脆弱 → 进程 CPU 自校准 —— 成立

`internal/risk/pathresolver_budget_norace_test.go`（322 行）＋ 3 枚平台件：

| 件 | 核 |
|---|---|
| `proccpu_windows_test.go`（`//go:build windows`） | `windows.GetProcessTimes` 取 user+kernel 进程 CPU（x/sys），与 `internal/speech` 既有同对 syscall 一致，包内自带仪器 ✓ |
| `proccpu_unix_test.go`（`//go:build unix`） | `syscall.Getrusage(RUSAGE_SELF)` user+kernel ✓ |
| `proccpu_unsupported_test.go`（`!windows && !unix`） | 返回 error ⇒ 用例首行 `t.Skipf` **具名跳过**，⛔ 静默空过 ✓ |

自校准三档口径（低／高／超）读码结论：

- **判据分两半**：半一＝每调用**进程 CPU** ≤ `resolveBudget × contention`（每次运行必判，承载契约，:152）；半二＝contention 落回下限 `1.00` 时每调用**墙钟** ≤ 1ms **无乘数**（票 18 原话，:165-169），否则日志明说 “not asserted”（:171），不装判、不静默删。
- **contention 自校准**：同一次运行内 3 轮 ≥100ms 纯浮点除法探针（无文件系统/无锁/sleep），取“墙钟/进程CPU”**中位数**，`!(f>1.0)⇒钳到 1.0`（:217-219）⇒ **只能替负载说话（只能放宽判离线），永远不能把契约线放宽到低于 1ms、也收不紧它**——“低档”＝钳回 1.00 时判离线恰为 1ms 原约；“高档”＝因子放大判离线、墙钟半边让位；“超档”＝CPU 越线即红。⛔ 新造预算常量：`const resolveBudget = time.Millisecond` 与改前逐字节同（改前 43 行版 :10 ↔ 新 :23，票面 :19“1ms 一个字不许动”✓）。探针参数常量（轮数/窗口时长/settle 拍距）是**测量窗尺寸**，不是预算阈值。
- **量化与迟到落账处理**：Windows CPU 钟 15.625ms 一跳＋窗口尾拍晚落 ⇒ 全部 ≥100ms 长窗、窗前后 `maxSettledProcessCPU`（睡一拍重读取最大）（:242-259）；探针某轮读不到 CPU ⇒ 该轮重取，整支读不到才判仪器坏 `t.Fatalf`（:224-227）——空窗不算过。
- **best-of-3 窗取最便宜 CPU 窗**（:126-135）：饱和时同一份字节会读出多个 CPU 值（本腿 Tier C 实测窗间 CPU 0.500↔0.988ms/op），取机器让得最少那窗＝代码成本；另两窗逐枚打印在日志（⛔ 拿最好那发冒充全体——三窗全列，非只报最优）。**此项是本形状里唯一需要盯的口径**：它放宽的是“读数噪声下限估计”，判离线仍锚在 1ms×自校准因子，且 AC#2 牙实测仍在（见③）。
- **仪器自检常驻用例** `TestBudgetInstrumentHasNoBlindSide`（:272-299）：算术腿 CPU 必须 >0（看不见工作＝空断言即红），纯睡腿 CPU ≤ 墙钟 1/4（把调度等待记成工作＝同墙钟一样脆即红）——两向钉死，⛔ 恒真。
- `//go:build !race` 保留（race runtime 会放大单发成本一个量级，票 18 原注释语义未丢）。
- 题面“testing.Short() 支”未采用＝票面 :20-21 是“**或**”关系（自校准支已满足）；票面 :21 的自校准示例“先在同一次运行里测一个基线倍数”正是 contentionFactor 的形状。方向合规。

## ② 真跑三档（指令原样 `-run 'TestResolvePerCallBudget' -count=1 -timeout 60m`，加 `-v` 取读数；`-overlay` 见环境段）——三档逐枚点名

| 档 | 并发负载 | rc | contention 三轮→中位 | 三窗 CPU(ms/op) | 判离线 | CPU 半边读数 | 墙钟半边 | 结果 |
|---|---|---|---|---|---|---|---|---|
| 档一 单跑 | 无 | **0** | 1.08/1.07/1.08 → **1.08x** | 0.388 / 0.417 / 0.414 | 1.075ms | 0.388 ≤ 1.075 | 未判（1.08≠1.00，日志明说） | `--- PASS (7.98s)` |
| 档二 与 `internal/config` 并跑（`-count=2` 同时段，负载侧 ok 4.508s） | 1 包活跃测试 | **0** | 1.00/1.07/1.07 → **1.07x** | 0.456 / 0.502 / 0.531 | 1.068ms | 0.396 ≤ 1.068 | 未判（1.07） | `--- PASS (7.73s)` |
| 档三 与全仓并跑（`go test -count=1` 34 包 `./cmd/wisp ./internal/...` 除 risk，运行期并发） | 34 包测试 | **0** | 1.07/1.07/1.28 → **1.07x** | 0.637 / **0.988** / 0.500 | 1.069ms | 0.500 ≤ 1.069（最脏窗 0.988 亦在窗内被记录） | 未判（1.07） | `--- PASS (9.20s)` |

- **FAIL=0 总计 3/3 档**。`--- SKIP` 逐枚点名：三档各跑内 **0 条**（Windows 平台 CPU 钟可用，`:68` 跳过门未触发）。
- 对照脆弱性本体（票面判定复现）：改前旧墙钟形同一台机器同份产码，单跑 0.354ms ↔ 全仓并跑 **2.044ms**＝5.8 倍极差 ⇒ 票面“负载敏感、⛔ 断言写反、⛔ 实现变慢”成立；新形三档把该极差折进自校准判离线后全部稳定绿。

## ③ 反向判据（回墙钟形）＋还原核 —— 成立

- **反向实验**：`-overlay` 把 `pathresolver_budget_norace_test.go` 还原为 `4d0e2577d^` 的 43 行原样（并虚拟移除 3 枚 proccpu 新件），同机同产码：
  - 全仓并跑档：`--- FAIL: TestResolvePerCallBudget`，读数 `C26 Resolve: 2044218 ns/op = 2.044 ms/op (budget 1.000 ms, 543 samples)` ⇒ **点名用例在该形下红** ✓（与 86-r1 报的 5.318ms 红样本同形不同幅＝负载函数，正是病灶）；
  - 单跑对照：同形 `PASS 0.354 ms/op` ⇒ 红不来自产码变化，只来自仪器读机器。
- **AC#2 牙（本腿独立复现，非引用 86-r1）**，均 -overlay 于 HEAD 新形：
  - **M-4x**（每枚调用多塞 3 次真 `Resolve`）⇒ **rc=1 FAIL**：`contention 1.00x → calibrated line 1.000ms`，`Resolve burns 1.826298ms of CPU per call, calibrated budget 1ms`（三窗 CPU 1.826/1.985/2.497）⇒ “每次解析的开销真的涨了必被抓”常驻断言有牙 ✓；
  - **M-sleep**（每枚调用多睡 700µs，CPU 不涨）⇒ 本腿该发 **PASS**：探针中位 1.07x≠1.00 ⇒ 墙钟半边按设计“not asserted”、CPU 半边 0.734≤1.068。86-r1 报过 contention 落 1.00 的同变异红样（1.598ms 墙钟越原线）。**具名记录**：墙钟半边只在探针中位数恰钳回 1.00 时才咬（本机本腿 6 发绿跑里仅 M-4x 一发 1.00x）；纯等待型回归在轻载机器上（因子停在 1.07）会只留日志不判——这是形状自洽的行为（票 18 句子只在它写就的场景＝无争用窗内生效），且票 86 契约线“真工作变贵必红”由 CPU 半边无条件承担，AC#2 底线不缺口。此为口径观察 ⛔ 放水项 ⛔ 未满足格。
- **还原核**：撤掉全部实验 overlay、回到 HEAD 形状复跑同指令 ⇒ `--- PASS (7.83s)`、0 SKIP、`internal/risk/` 与 HEAD 零差异；本腿未对仓内任何被跟踪文件产生写入。

## ④ 纪律核 —— 全部符合

- `git show --stat 4d0e2577d`：写面＝**恰好 4 枚文件**，全在 `internal/risk/`、全 `*_test.go`：`pathresolver_budget_norace_test.go` +311（题面“+311”一致）、`proccpu_windows_test.go` +29、`proccpu_unix_test.go` +22、`proccpu_unsupported_test.go` +15；合计 +361/−16。与派单题面逐字一致 ✓。
- **零产码**：4 枚皆测试文件；`git show --name-only` 无 `internal/risk` 外任何路径 ✓；`cmd/wisp/**`、`internal/panel/**` 零触碰 ✓。
- **⛔ thresholds/golden**：`internal/observe/thresholds.go` 及任何 golden 不在提交文件列表（逐字核 `--name-only`）✓；SLO/阈值一字节未动；判离线只由运行内实测因子缩放，⛔ 改常量（`resolveBudget` 原字节保留，见①）。
- Git 纪律：落地腿只 commit 未 push（`origin/dev` 侧无该笔的强推痕迹，本腿亦不 push）；无 `--amend/reset/rebase/stash` 痕迹（历史线性，A1416 收件笔直接父即 4d0e2577d）。
- 门禁（AC#4，-overlay 于提交态）：`gofmt -l internal/risk/` 空 ✓；`gofumpt -l internal/risk/` 空（`D:\work\base\gopath\bin\gofumpt.exe`）✓；`go vet ./internal/risk` rc=0 ✓；交叉平台 `GOOS=linux vet` rc=0、`GOOS=darwin vet` rc=0（unix 兄弟件可编译）✓；`go test -count=2 ./internal/risk/` rc=0（ok 25.875s）✓，全包 `--- SKIP` 计数 **0** ✓。

## AC 逐格裁决（票面 86-…-wallclock-fragility.md :26-31）

| 格 | 票面要求 | 读数 | 判 |
|---|---|---|---|
| **AC#1** | 旧仪器三档真实分布、全部样本、核数 | 原始 43 枚样本表**未入库**（86-r1 漏产 `86r1-impl.md`）；数字仅存于 `4d0e2577d` 提交说明（0.395→1.625、三枚红 1.020/1.511/1.625）＋核数 12；**本腿等价复现**：旧形同字节单跑 0.354 绿 ↔ 全仓并跑 2.044 红（极差≥5.8x，负载档红、单跑档不红＝票面判定复现） | **PASS（附证据保留瑕疵，见“保留”a）** |
| **AC#2** | 改自校准（或 Short 门控）且保留一条抓真回归的断言 | 自校准形状✓（①）；牙复现：M-4x 实测红✓（③）；`TestBudgetInstrumentHasNoBlindSide` 常驻防恒真✓ | **PASS** |
| **AC#3** | 三档并重跑 rc=0、无新增 SKIP 逐条点名 | ②表：三档全 rc=0、三档+全包 0 条 SKIP（点名＝无） | **PASS** |
| **AC#4** | gofmt 空、vet rc=0、`-count=2` rc=0 | ④全绿（另加 gofumpt、linux/darwin 交叉 vet 亦绿） | **PASS** |

## 保留（不入 AC 格、具名不抹）

a. **证据保留瑕疵**：AC#1 的 43 枚原始样本清单没有落进仓（缺 `docs/evidence/s1/86r1-impl.md`）。判 PASS 的依据＝提交说明数字＋本腿独立复现的旧形极差；“全部样本可逐枚复查”这一半只能由编排者补账（建议台账记一笔：落地腿证据件缺产、由 86-v1 verdict 代持读数）。
b. **口径观察**：墙钟半边（票 18 原话）只在探针中位恰钳 1.00 的安静窗咬合；本机常态 100ms 窗的 CPU 钟量化使中位停在 ~1.07（本腿 6 发绿跑均未触发该半边）。契约承载在 CPU 半边（有牙实测），不构成缺格，但后续若在 CI 想钉“纯等待回归”，要么加长探针窗（>6.25 倍量化步），要么把墙钟半边触发条件从 `==1.00` 改为“≤1.0x+容差”——均属后续票口径，⛔ 本腿改码。

## 是否同意结 86

**同意结 86。** 四格 AC 均 PASS；①形状、②三档 rc=0 零 SKIP、③反向判据（旧墙钟形在同字节下并跑红 2.044ms/op、具名 `TestResolvePerCallBudget`）与 M-4x 牙红、还原核绿、④纪律（写面恰 4 枚 `internal/risk/**_test.go`、零产码、⛔ thresholds/golden/cmd 面、提交与题面逐字一致）全部核过。两条“保留”（a 证据件缺产、b 墙钟半边触发口径）为记账/后续票事项，不构成未满足 AC 格，不阻结。
