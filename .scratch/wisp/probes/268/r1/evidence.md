# 票 268 落地腿 `268-r1`：证据件（A＋B 组合，⛔ 不取 C／doctor）

编队＝`268-r1`（产码腿）。形由编排者裁死（账 `A614` §5），本件不重裁形，只落地＋自证。

- 起手 HEAD：`cfd97636`（`git log --oneline -1` 现量），分支 `dev`，同一枚共享工作树。
- 起手时刻：`2026-10-05 13:00:47 +0800`（`date` 现量）。
- 写面：`cmd/wisp/resident_approval_windows.go`（改）＋ `cmd/wisp/resident_approval_risk_268_windows_test.go`（新增）。**逐枚名册见 §2.4**。
- 上游料：`.scratch/wisp/probes/268/a1/census.md`（普查）＋ `.scratch/wisp/probes/268/a2/verdict.md`（独立复核，
  推翻 census 六条 ⇒ **本腿以 a2 为准**，并照它的 §4 U4 静态名册做撞钉预检、照它的 §2④ 修正把 `fed/bare` 射程收为"构造名计数"）。
- 本件的尺：占位符一律按字符类 `待[填]|填写[中]|未判` 量（字面串会咬自己）；计数一律 `| wc -l` 收尾。

## §0 起手锚与现读（票面所有 `file:line` 当待验断言）

派单给的每一枚锚本腿都亲读过（读的是**工作树**）。**结论：十三枚锚逐枚未漂，无需报漂因。**

| 派单给的锚 | 本腿亲读结果（工作树，13:01–13:03） | 判定 |
|---|---|---|
| `cmd/wisp/resident_approval_windows.go:434` 逐字 `riskProvenanceUnreadable = "defaults (config.toml unreadable)"` | 逐字在 `:434` | 未漂 |
| `…:457` `return 0, 0, riskProvenanceUnreadable` | 逐字在 `:457`，在 `:453` 的 `if err != nil \|\| c == nil {` 里 | 未漂 |
| `…:454-456` `slog.Warn(…, "err", err, …)` | 逐字跨 `:454-456`（msg `:454`、`"path"/"err"` `:455`、`"fallback"` `:456`） | 未漂 |
| `…:379-384` 唯一交给用户的那条 `slog.Info` | 逐字在 `:379-384`，读的是**局部变量** `provenance`（`:380`） | 未漂 |
| `…:449` NoView 支 `return 0, 0, riskProvenanceNoView` | 在 `:449` | 未漂 |
| `…:367` 写字段 `ra.riskWindow, ra.riskTimeout, ra.riskProvenance = …` | 在 `:367` | 未漂 |
| `cmd/wisp/resident_windows.go:187` `[hotkey]` 那一行 `fmt.Printf`（B 的形状母本） | 逐字在 `:187`，同一闭包里 `slog.Warn`（`:185-186`）**加**一行 Printf（a2 §2③ 同判） | 未漂 |
| `cmd/wisp/config_reload.go:396` `case strings.HasPrefix(d, "config.toml:")` | 逐字在 `:396`（`describeReloadFailure` 的 `cause=invalid` 支，`:354-409`） | 未漂 |
| `cmd/wisp/resident_approval_risk_256_windows_test.go:141-144` AC#2 既有钉 | 逐字在（`got != approval.DefaultApprovalTimeout` ⇒ Errorf，注释含 "a fourth number here means the fallback grew a value of its own"） | 未漂 |
| 同文件 `:560` `if fed != 1 {`／`:563` `if bare != 0 {` | 逐字在；AST 尺射程＝`runResident` 体内**两个构造名**的调用计数（a2 §2④ 的收窄口径） | 未漂 |
| 同文件 `:310-314`（construction-time-only 钉）＋ `:317-319`（provenance 不漂移钉） | 逐字在 | 未漂 |
| 同文件 `:124-125`（形状表：`{"config.toml missing", "unreadable", riskProvenanceUnreadable}`）／`:247-249`／`:317-319` 三组字段读数 | 逐字在 ⇒ a2 §2① 的"七组读者"成立，**字段断言不止在回落用例里** | 未漂 |
| `internal/config/validate.go:128-129` band＝`[31, 3600]`、`:145-148` 拒载句 | `confirmTimeoutSecMin = 31` 在 `:128`、`Max = 3600` 在 `:129`；`observe.New(ClassConfig, "config.toml: risk.confirm_timeout_sec %d out of range [%d, %d]")` 在 `:145-148` | 未漂 |
| `internal/config/loader.go:67-79` 缺失支保留 `fs.ErrNotExist` cause | `fileMissing(err)` 支在 `:67`，`observe.Wrap(observe.ClassConfig, err, "config.toml read: no file at this path yet。…")` 跨 `:76-79` ⇒ cause 在链上 | 未漂 |

**两条本腿自己补的现读（派单没给，但落地要用）**：

1. `config.LoadFile`（`internal/config/loader.go:41-51`）**不可能返 `(nil, nil, nil)`** —— 三支 return 逐枚读过。
   ⇒ `residentRiskGateValues` 的 `err != nil || c == nil` 里"`c == nil` 但 `err == nil`"那半是空集，
   所以新支不必为它造第五枚名字（否则那一支永远不被生产代码走到，是一枚装饰性断言）。
2. `cmd/wisp` 现在**没有任何一枚用例**给常驻腿种一枚带外 `[risk]` 值：尺＝
   `grep -rn "confirm_timeout_sec = \|l1_window_sec = " cmd/wisp/*_test.go | wc -l` → **14**，逐枚读过，
   带内的才算（45/90/31/40 与 2/1/99 那枚是 window 侧的钳位反控，window 无加载期 band，见 `internal/config/schema.go:457-463`），
   **零枚带外 confirm 种子** ⇒ 我那一行 stdout 在既有整包里**一次都不会被触发**（与 a2 §4 的"撞钉 0 枚"同向、比它更强：它量的是"即使触发也不撞"，本腿量的是"根本不触发"）。

**票 128 AC#5 的措辞撞否（本腿的判语：不撞，两句并排如下）**：

- 票 128 AC#5 的落点现量在 `internal/proc/envfork.go:237-239`，逐字仍是
  `return Layout{}, fmt.Errorf("proc: user config dir: %w", err)` —— **AC#5 未落地**（那句自救三样还缺）。
  它说的是**数据根本身求不出来**（`%AppData%` 没定义），形是"**拒绝启动**"。
- 本票那一行说的是**数据根在、config.toml 也在、里面的值被加载层拒了**，形是"**回落并继续跑**"。
  并排读：`proc: user config dir: …` vs `wisp: resident [risk]: config.toml is present but was refused at load (…)` ——
  主语（配置目录 vs 配置文件）、后果（拒启 vs 回落）、前缀（`proc:` 错误串 vs `wisp: ` stdout 行）三处都不同，
  **没有共同的"配置没读到"句式**，故不停手。⚠ 若日后 128 AC#5 的自救句选词落到 `"config.toml …"` 那一族，
  请在那一票里复看这一行；本腿不替它预先让路，也**不碰 `internal/proc` 一字节**。
- 同族第三句（现量在 `cmd/wisp/firstrun.go:92-95`，`wisp run: 已在 %s 新建默认配置…`）说的是"首建"，与本票无重叠；
  它刻意不复用 `cause=missing` 那句的纪律（`firstrun.go:88-91` 注释逐字）本腿照抄精神：一句只管一事。

## §1 撞钉预检（今日绿名册）

规矩：光 grep 新符号名不算预检，**必须整包真跑取今日名册**。

| 尺 | 读数 | 时刻 |
|---|---|---|
| 起跑前 `powershell -NoProfile -Command "@(Get-CimInstance Win32_Process -Filter \"Name='go.exe'\").Count"` | **0**（无同机并发整包） | 13:00:5x |
| `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test -count=1 -v ./cmd/wisp/ > .scratch/wisp/probes/268/r1/baseline-cmdwisp.log 2>&1` | **RUN 331 / PASS 331 / FAIL 0 / SKIP 0**，包尾 `ok github.com/CarlosShao/wisp/cmd/wisp 466.320s`，`rc=0` | 起 13:05:56／止 13:13:47 |
| `grep -c '^=== RUN' / '^--- PASS'…`（同上四把） | 331／331／0／0（子用例含在内；顶层名册另存 `baseline-top-names.txt`＝**235 枚**顶层名） | 13:13:5x |
| 编排者点名的四枚带载偶发／在册变色件 | 今天**全绿**，一枚不跳：`TestAC14GoSideEvalPushReachesThePage (0.49s)`、`TestAC1ResidentLegInstallsItsLogListenerOnDisk (4.48s)`、`TestPanelHostRealWindowHopAndLifecycle (1.19s)`、`TestPanelHostLatencyPercentilesAC2 (0.00s)` | 13:13:5x |
| 红名册 | **空**（`grep '^--- FAIL\|^    --- FAIL'` → 0 行） | 13:13:5x |

**⛔ 未压种子、未 `t.Skip`、未把 SKIP 读成通过**（SKIP＝0 枚，本就无此面）。

预检之外本腿另做的三把静态撞钉尺（都指向"这一行 stdout 会不会踩谁"）：

1. stdout 捕获面：`captureStdout128` 只在 `cmd/wisp/dataroot_128_test.go:410-438` 定义，
   换的是**进程级 `os.Stdout`**、`t.Cleanup` 还原 ⇒ 我复用它是安全的（读全过：`w.Close()` 后收干再返回）。
2. 常驻子进程 stdout 断言形状：`grep -rn "Printf\|Stdout" cmd/wisp/*_test.go | grep -i "sel.Name"` → **0 枚**
   （没有任何 AST 尺数 `fmt.Printf` 的枚数）；`grep '^--- FAIL' … ` 之外唯一一枚"精确计数"在
   `cmd/wisp/resident_sink_nail_127_windows_test.go:493`，数的是 **stderr 里 winsec 封缝句** 的 `strings.Count(...) != 1`，
   与本腿的 stdout 行不同流、不同串。⇒ 与 a2 §4 的静态名册判定一致：**加一行 stdout 撞 0 枚钉**。
3. 包目录扫描型用例（`cmd/wisp/leg_dispatch_gate_133_test.go:365-372`、`resident_ball_228_test.go:66-74`、
   `panel_host_gate_test.go:210-218`、`resident_grant_writer_265_windows_test.go:547`）逐枚读过：
   **生产侧扫描一律按名字剔掉 `_test.go`**，我的新文件是 `_test.go` 后缀 ⇒ 不进任何分母；
   生产那一枚 `resident_approval_windows.go` 我只加常量/分支/一行 Printf，不新增 `ball.New`、不新增 Dispose/Destroy 站点。

## §2 落地内容

### 2.1 A 支：把两枚状态分开（`cmd/wisp/resident_approval_windows.go`）

1. 新常量（第四枚 provenance 词）：
   `riskProvenanceRefusedAtLoad = "defaults (config.toml present but refused at load)"`，
   与 `riskProvenanceUnreadable = "defaults (config.toml unreadable)"` 逐字不同名。
   ⇒ 词面为什么选 "present but refused at load"：**它只说自己确定知道的事**（文件在、没加载成），
   不冒充"是值域门拒的"。a2 §2② 的告诫本腿照办——非缺失那一支今天还含着语法／未知键／迁移／权限诸支，
   `observe.New` 造的拒载错误**没有 cause 可看**（`Error.Err == nil`）、`ClassOf` 又同为 `ClassConfig`，
   所以任何"就是 band 拒的"式命名都是新谎言。具体是哪一条规则，由**紧跟其后的 `%v` 原话**承担（日志 attr 与 stdout 行都带）。
2. 分派只用**哨兵**：`if errors.Is(err, fs.ErrNotExist)` ⇒ 缺失支走**逐字未动**的原 `slog.Warn` ＋ `riskProvenanceUnreadable`；
   非缺失支走新 `slog.Warn`（多带一枚 `"provenance"` attr）＋ 新 stdout 行 ＋ `riskProvenanceRefusedAtLoad`。
   ⛔ 没有 `strings.HasPrefix`／错误文案子串分类（那是 `cmd/wisp/config_reload.go:396` 的形状，本票明令不取；
   新常量所在函数由 §2.3 的 ③ 号钉从**语法树**上钉住这一点）。
3. 文档面同步（只改注释，不改语义）：常量块抬头 "three provenance words" → "four"；
   `riskProvenanceUnreadable` 的注释从 "(missing or broken)" **收窄为只指缺失**并点名票 268 拆开了那一折；
   `residentRiskGateValues` 的 "which of the three shapes" → "four"。
4. import 只加一枚 `io/fs`（`errors`／`fmt`／`log/slog` 本来就在）。**未新增任何包级依赖边。**

### 2.2 B 支：被拒那一句同时打到 stdout（同一枚 `residentRiskGateValues`）

```
wisp: resident [risk]: config.toml is present but was refused at load (<loader 原话>); the approval gate
falls back to the compiled constants (DefaultApprovalTimeout=300s / DefaultL1Window=3s), so the [risk]
numbers written in that file are NOT the numbers this process runs on
```

- 形状**照抄** `cmd/wisp/resident_windows.go:187` 那枚 `[hotkey]` 行（同前缀族 `wisp: … [段]: …`、同"一句一行、括弧里带 err、句尾点名回落对象"）。
  现量事实复认：`[hotkey]` 族有 1 行 Printf、`[risk]` 族原来**0 行**（a2 A7 同判）。
- ⛔ **不新增第二枚 `[risk]` 读取点**：这一行复用的是本函数**已经拿到手**的 `err`，不重读文件、不取 Manager、
  不在 `runResident` 里加任何调用 ⇒ 256 的 `:560/:563` 构造名计数钉射程未被触碰（按 a2 §2④ 的收窄口径，真正的钉是"别多调一枚构造函数"）。
- 已知事实、不算失败：`-H=windowsgui` 双击起法（`scripts/build.ps1:115`）下这行无处可去（`cmd/wisp/console_windows.go:38/:49-51`
  的 `attachParentConsole` 无父控制台即不改绑 std）。终端起法下它与人眼之间没有脱敏层。
- ⛔ **不取 C**：本轮 doctor 零改动，未加任何检查项（`cmd/wisp/doctor.go` 不在我的写面里，见 §2.4）。
  理由照派单：`scripts/build.ps1:167-169` 拿 `wisp.exe doctor` 的退码当构建 smoke，`doctor.go` 的 `fail()` 置 `critical:true`
  ⇒ 加一枚 `fail` 级检查＝把任何"配置里有带内越界值"的机器的**构建门**打红。本腿连 `info()` 级也没加（自发觉得"要不要加一条"
  的冲动在此具名登记为**未取**：AC#1 的凭据是 stdout 那一行，不是 doctor）。

### 2.3 新增用例（`cmd/wisp/resident_approval_risk_268_windows_test.go`，三枚）

| # | 用例 | 钉的是什么 | 红谁 |
|---|---|---|---|
| ① | `TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig` | 被拒／缺失两枚 provenance **各自具名且互不相等**；四枚常量两两不同字面；被拒支仍 `Queue().Timeout() == approval.DefaultApprovalTimeout`（300 s）**且不等于 20 s**（越界值没溜进带内＝`Q-77` 未动）；`Window() == DefaultL1Window`；**前提钉**：种的文件 `os.Stat` 在、`config.LoadFile` 报错、`errors.Is(err, fs.ErrNotExist)` 为假、原话含 `out of range [31, 3600]`；**配对正控**：同一目录把值修成 45 后立刻读成 `riskProvenanceRead` | MUT-1（折回同名） |
| ② | `TestTicket268RefusedRiskConfigReachesStdoutOnce` | 同一次构造里 stdout 上 `wisp: resident [risk]:` **恰 1 行**、且那行带 `out of range [31, 3600]` 与 `DefaultApprovalTimeout=300s`；**反控**：缺失支在同一张面上零行（两枚状态不许又共用一张脸） | MUT-2（摘掉 Printf） |
| ③ | `TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords` | 从语法树上读 `residentRiskGateValues`：`errors.Is(err, fs.ErrNotExist)` **恰 1 次**、函数体内 `strings.*`／`err.Error()` **0 次**、`fmt.Printf` **恰 1 次**、四枚 provenance 常量各有一枚 return 走到 | MUT-1／MUT-3（把哨兵换成文案前缀匹配） |

复用同包既有辅助件而非重造：`risk256Dir`（`cmd/wisp/resident_approval_risk_256_windows_test.go:76-83`）与
`writeRiskConfig256`（同文件 `:89-100`）；stdout 捕获复用 `captureStdout128`（`cmd/wisp/dataroot_128_test.go:410-438`）。
**枚数尺**：新文件 `func Test` 顶层＝**3 枚**（`grep -c '^func Test' cmd/wisp/resident_approval_risk_268_windows_test.go`）。

### 2.4 写面名册与禁区

- 本腿改／增的文件（**逐枚**，与 `git show --name-only` 对得上）：
  `cmd/wisp/resident_approval_windows.go`（M）、`cmd/wisp/resident_approval_risk_268_windows_test.go`（A）。
- 禁区零字节（现量尺与时刻见 §4.5）：`docs/PLAN.md`、`docs/specs/**`、`internal/observe/thresholds.go`、
  任何 golden、`tools/d22scan/allowlist.txt`、`internal/agent/approval/**`、`frontend/**`、`design/**`。
- `[risk].l1_window_sec` 与它的钳位：一字未动（`internal/config/schema.go:457-463`、`internal/agent/approval/gate.go` 均不在写面）。
- 值域 `[31, 3600]`：一字未动（`internal/config/validate.go:128-129` 不在写面；我只在测试里**种**一个带外值）。
- ⛔ 入向方法名（C17 契约面）：生产 diff 里 `^\+func `／`^\+\tfunc ` 计数＝**0**，新增的唯一标识符是未导出常量
  `riskProvenanceRefusedAtLoad` ⇒ 零新增入向面。
- ⚠ 共享工作树观察（**不是本腿所为、本腿未读未改**）：起手 `git status` 即见 `design/**` 多枚 ` D`/` M`、
  `.gitignore` 为 ` M`、`236-r1` 的两枚 staged `R`。本腿的 commit 一律 `git commit -F msg -- <点名 pathspec>`，
  逐笔 `git show --name-only` 自证未带上它们（见 §7）。

## §3 判据逐格读数

| 格 | 判语 | 凭据（命令＋读数＋时刻） |
|---|---|---|
| AC#1 | 达标（待 §5 突变自证补全后终判） | `PATH=… GOFLAGS= go test -count=1 -v -run 'TestTicket268\|TestTicket256' ./cmd/wisp/` → `--- PASS` 三枚 268 用例全绿，**13:17:08**；具名 provenance 新常量在 `cmd/wisp/resident_approval_windows.go`（§2.1），stdout 一行在 §2.2 |
| AC#2 | 达标 | 同一发里 `TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants`（含 `:141-144` 那枚 300 s 钉）**保持绿**；新增"文件在但被拒 ⇒ 仍 300 s 且 provenance 具名"＝268 用例① 绿；`TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly`（`:310-314/:317-319`）与 ②④⑤ 全部保持绿 |
| AC#3 | 达标（未触碰） | 写面名册里没有 `cmd/wisp/run.go`（`git status --short -- cmd` 只两枚：`resident_approval_windows.go` 与 268 新用例）；跑任务腿的退码 2 那一句一字未改。一响一静仍是有名字的差别 |
| AC#4 | 达标（尺见 §4.5） | 禁区 `git diff --name-only HEAD -- docs/PLAN.md docs/specs internal/observe/thresholds.go tools/d22scan/allowlist.txt internal/agent/approval` → **0 行**；`find internal/agent/approval docs/specs -newer <本腿起手件> -type f` → **0 枚**；band 与 l1 钳位不在写面；新增函数 **0** 枚 |

## §4 门禁四数

未判（交付态自跑，逐门带时刻）。

### 4.5 禁区零字节尺（现跑）

未判。

## §5 突变名册（每发起止到秒）

未判。

## §6 我攻不动／判不动的地方

未判。

## §7 交件判语与 commit 链

未判。
