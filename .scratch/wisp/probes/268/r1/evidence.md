# 票 268 落地腿 `268-r1`：证据件（A＋B 组合，⛔ 不取 C／doctor）

编队＝`268-r1`（产码腿）。形由编排者裁死（账 `A614` §5），本件不重裁形，只落地＋自证。

- 起手 HEAD：`cfd97636`（`git log --oneline -1` 现量），分支 `dev`，同一枚共享工作树。
- 起手时刻：`2026-10-05 13:00:47 +0800`（`date` 现量）。
- 写面：`cmd/wisp/resident_approval_windows.go`（改）＋ `cmd/wisp/resident_approval_risk_268_windows_test.go`（新增）。**逐枚名册见 §2.4**。
- 上游料：`.scratch/wisp/probes/268/a1/census.md`（普查）＋ `.scratch/wisp/probes/268/a2/verdict.md`（独立复核，
  推翻 census 六条 ⇒ **本腿以 a2 为准**，并照它的 §4 U4 静态名册做撞钉预检、照它的 §2④ 修正把 `fed/bare` 射程收为"构造名计数"）。
- 本件的尺：占位符按派单给的那条三选一字符类量（`待` + 方括号类、`填写` + 方括号类、以及第三个分支那两个字的连写）；
  **本行刻意不把它们写成字面串，免得尺咬自己**。计数一律 `| wc -l` 收尾，带 `| head -N` 的输出只当样例不当枚数。

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
| 交付态复跑同一把尺（三发突变全部还原之后、代码与 §4 那四门同形） | **RUN 334／PASS 334／FAIL 0／SKIP 0**，包尾 `ok github.com/CarlosShao/wisp/cmd/wisp 401.239s`，`rc=0`，台件 `deliver-cmdwisp.log`（1,724 行） | 止 15:12:06（log mtime 现量） |
| 红名集合逐名比对（`diff baseline-top-names.txt deliver-top-names.txt`） | 差集＝**只多我三枚**：`TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig`／`TestTicket268RefusedRiskConfigReachesStdoutOnce`／`TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords`；**一枚未少、一枚未红**（235 枚 → 238 枚顶层名，`--- ` 级） | 15:12:5x |
| 两发包尾耗时并记（不作读数，只作负载登记） | 基线 `466.320s`（起跑前 `go.exe`＝0，但同机另有 244-a4／parking-1 在读与在写文档）／交付 `401.239s` ⇒ 差 65 s 属他腿负载，**不判为"我变快了"** | 13:13:47／15:12:06 |

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
| AC#1 | **达标**（具名状态一枚＋一行 stdout 到达眼睛＋折回同名必红自证） | 新常量 `riskProvenanceRefusedAtLoad = "defaults (config.toml present but refused at load)"` 与 `…Unreadable` 逐字不同名（`cmd/wisp/resident_approval_windows.go` 常量块）；stdout 一行照 `[hotkey]` 形（母本 `cmd/wisp/resident_windows.go:187`）落在同一枚函数里。发跑：`-run 'TestTicket268\|TestTicket256'` 三枚新用例全绿（13:17:08），交付态整包 **334／334／0／0**（止 15:12:06）；**MUT-1**（两枚字面折回同一枚）⇒ 用例① 红，三句原文抄在 §5；**MUT-2**（摘掉 Printf）⇒ 用例②③ 红。**最低档／达标档的诚实交代**：编排者裁的形把 stdout 那一行记为达标凭据，本腿按那一行交；该行在 `-H=windowsgui` 双击起法下看不见（§6.1，已知事实） |
| AC#2 | **达标**（回落仍是编译 300 s，越界值没趁机生效） | 既有钉 `cmd/wisp/resident_approval_risk_256_windows_test.go:141-144` 所在用例 `TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants`（含 `config.toml_missing` 子枚）在 13:17:08／13:36:38／14:05:08／14:53:26 四发里**逐发保持绿**；新增一发＝268 用例①：文件在、被拒 ⇒ `Queue().Timeout() == approval.DefaultApprovalTimeout` **且不等于种进去的 20 s**；同发 256 的 ②③④⑤（含 `:310-314/:317-319` 的 construction-time 钉）逐枚绿。`Q-77` 未动、一枚断言未放宽 |
| AC#3 | **达标**（跑任务那条腿一字未改） | 写面名册里没有 `cmd/wisp/run.go`：`git show --name-only --format="" 9a941965`＝5 行（2 枚产码＋3 枚文档面），`4c456d0a`＝3 行文档面（15:07:14 复量）；退码 2 那句响亮拒绝原样在册（票 267 已结案形状）。两腿"一响一静"没有被顺手统一，本腿也没有给缺失支加任何新面（用例②的反控钉着这件事） |
| AC#4 | **达标**（禁区零字节＋卫生四门不扩＋无新入向方法名） | 尺与读数在 §4.5（禁区命中 **0 行**／`find -newer` **0 枚**／生产 diff 新增 `func` **0 枚**）；四门 rc 全 0 且名册逐名一致，唯一分母变化＝d22scan `ban #8 cmd/` 的 Go 文件枚数 103→**104**（我新增的那一枚 `_test.go`，违规仍 0 枚，见 §4 末段）；band `[31,3600]` 与 `l1_window_sec` 钳位不在写面 |

## §4 门禁四数

四门都在**落地态**自跑（13:32–13:33），交付态与落地态**逐字节同形**——证在 §5：三发突变全部还原后
`md5sum cmd/wisp/resident_approval_windows.go` ＝ `git cat-file blob HEAD:…` ＝ `69a630bf350daf1032c62b921e9d7044`，
`git diff --name-only` 对该路径 0 行；`cmd/wisp/resident_approval_risk_268_windows_test.go` 自 `9a941965` 起未被任何一发写过。

| 门 | 命令原文 | 读数 | 时刻（落地态／**交付态复跑**） |
|---|---|---|---|
| 1 | `sh scripts/d22scan.sh` | **rc=0**，`d22scan: clean - no D22 ban violations`。正控制那步真跑了：`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77`；扫描面 `examined 266 production Go files under internal/ and cmd/` | 13:32:02–13:32:34／**15:25:29–15:25:54 同读数** |
| 2 | `sh scripts/check-path-length-budget.sh --with-self-test` | **rc=0**、`positive control PASSED`、`VERDICT GREEN`；`denominator: tracked paths=5953 over-budget=57 covered by roster=57 not in roster=0`；`longest=180 chars relative`＝票 252 那枚工单名（⛔ 不是本腿造的，本腿新增面最长＝`cmd/wisp/resident_approval_risk_268_windows_test.go`＝**50** 字符） | 13:32:45–13:32:51／**15:26:32–15:26:36 同读数（57／57／0）** |
| 3 | `GOFLAGS= go vet ./cmd/wisp/` | **rc=0**（无输出） | 13:33:01／**15:26:37 rc=0** |
| 4 | `"D:/work/base/gopath/bin/gofumpt.exe" -version` 先自证尺活着，再 `-l cmd/wisp` | 尺＝**v0.12.0 (go1.27.1)**；`-l` 输出＝**只有 `cmd\wisp\models.go`**（派单点名的预存脏枚，本腿一字未碰）。⛔ 未用 `2>/dev/null`，空读数与非空读数分得开。本腿两枚文件**不在名册里** | 13:33:03／**15:26:39 同读数**（中途 14:09:34 也复量过一次，三名一致） |

八枚 scope 读数（同一发 d22scan 里逐枚抄，供与历史对账）：
`bans #1-5 internal/=228`、`bans #1-5 cmd/=38`、`ban #6 frontend/=85`、`ban #7 internal/tools/=23`、
`ban #8 design/=39`、`ban #8 frontend/=85`、`ban #8 internal/=512`、`ban #8 cmd/=104`。
⚠ 诚实登记一处**本腿造成的分母变化**：`ban #8 cmd/` 这把尺数的是 Go 文件枚数（含 `_test.go`），
本腿新增一枚 `_test.go` ⇒ 该读数比开工前多 1（104＝103＋本腿那枚）。这是"射程面多了一枚文件"，不是"违规多了一条"——
违规枚数＝**0**、红名集合＝**空**。其余七枚本腿无从改变（同一枚函数内的常量与分支不新增文件）。

### 4.5 禁区零字节尺（现跑）

| 尺（命令原文） | 读数 | 时刻 |
|---|---|---|
| `git show --name-only --format="" 4c456d0a` ＋ `git show --name-only --format="" 9a941965` ＋ 对两者并集跑 `grep -E "^(docs/PLAN.md\|docs/specs/\|internal/observe/thresholds.go\|internal/agent/approval/\|tools/d22scan/allowlist.txt\|frontend/\|design/)" \| wc -l` | 两笔逐枚点名：`4c456d0a`＝3 行（票面＋evidence＋msg-01），`9a941965`＝5 行（票面＋evidence＋msg-02＋**2 枚产码**）；禁区命中＝**0 行** ⇒ 五处禁区零字节，且 `frontend/`／`design/` 命中 0（共享树里 `design/**` 那批 ` D`/` M` 未被我任何一笔带上，见 §2.4 末条） | 15:07:14 |
| `find internal/agent/approval docs/specs -newer .scratch/wisp/probes/268/r1/msg-01-skeleton.txt -type f \| wc -l` | **0 枚**（这两处地界连"被谁改过"都没有） | 13:17:43 |
| `git diff -U0 -- cmd/wisp/resident_approval_windows.go \| grep -E "^\+func \|^\+\tfunc " \| wc -l` | **0 枚**新增函数声明 ⇒ 无新入向方法名（C17 面未动）；新增的唯一标识符是未导出常量 `riskProvenanceRefusedAtLoad`（生产 diff 64 增／8 删） | 13:17:43 |
| 值域 `[31, 3600]`／`[risk].l1_window_sec` 与钳位 | `internal/config/validate.go`、`internal/config/schema.go`、`internal/agent/approval/*` 均不在写面 ⇒ 一字未动（本腿只在测试里**种**一个带外值 20 去撞门） | 13:17:43（写面名册 13:17:43 复量） |

## §5 突变名册（每发起止到秒）

三发都种在**本腿自己的** `cmd/wisp/resident_approval_windows.go` 里；⛔ 一枚未种进
`internal/config`／`internal/agent/approval`／`internal/panel`／`internal/ball`。
尺＝`go test -count=1 -v -run 'TestTicket268|TestTicket256' ./cmd/wisp/`（台件逐枚留在本目录）。

| # | 种什么 | 起 | 止 | 红了谁 | 还原证 |
|---|---|---|---|---|---|
| MUT-1 | 把新常量折回旧字面：`riskProvenanceRefusedAtLoad = "defaults (config.toml unreadable)"`（＝AC#1 要求自证的那一发"两枚状态又折回同一枚字符串"） | 13:34:05 | 13:36:38 | **只 ① 红**：`--- FAIL: TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig`；票 256 五枚钉全绿、②③ 全绿（②③ 量的不是字面） | `md5sum`＝`69a630bf350daf1032c62b921e9d7044`＝`git cat-file blob HEAD:cmd/wisp/resident_approval_windows.go \| md5sum`（13:48:08 复量，`git diff --name-only` 0 行） |
| MUT-2 | 摘掉那一行 stdout（注释掉 `fmt.Printf`）——证 ② 不是恒绿装饰 | 13:52:04 | 14:05:08 | **②③ 双红**：`--- FAIL: TestTicket268RefusedRiskConfigReachesStdoutOnce`＋`--- FAIL: TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords`；① 与票 256 五枚全绿 | 14:09:34 复量 md5 同值、`git diff` 0 行 |
| MUT-3 | 在哨兵旁边**加**一次文案匹配（`err.Error()[:1] == "X"`）——证 ③ 咬的是"靠文案分类"这件事；纯文案支会因 `io/fs` 变成未使用而编译不过，故取这一形（行为逐字不变） | 14:10:16（if 行改动 14:52 前后复量） | 14:53:26 | **只 ③ 红**（① ② 与票 256 五枚全绿＝隔离成功） | 第一发还原**漏了一截注释**（只回 if 行）⇒ `md5sum` 读出 `8fcd72af…` 与 HEAD blob 不符、`git diff --name-only` 1 行；**由 md5 尺当场抓到**，15:03:35 以派单许可的形 `git cat-file blob HEAD:<path> > <path>` 全量还原，复量 md5＝`69a630bf350daf1032c62b921e9d7044` 同值、`grep -c "MUT-"`＝0、`git diff` 0 行 |

红句逐字（原文抄自台件，行号＝台件行）：

MUT-1（`mut1-red.log:46-48`）：

```
resident_approval_risk_268_windows_test.go:125: the two shapes are folded back into one word "defaults (config.toml unreadable)" - that fold IS ticket 268: the user wrote a number, the gate runs the compiled 300s, and the receipt blames a file that is sitting in that directory
resident_approval_risk_268_windows_test.go:129: refused-file provenance "defaults (config.toml unreadable)" does not say out loud that the file is present, which is the whole difference from the missing shape
resident_approval_risk_268_windows_test.go:147: provenance constants riskProvenanceUnreadable and riskProvenanceRefusedAtLoad share the literal "defaults (config.toml unreadable)" - the four readings have to be four names
```

MUT-2（`mut2-red.log:50/51/52/57`，四条同发）：

```
resident_approval_risk_268_windows_test.go:201: the refused branch wrote 0 stdout lines carrying "wisp: resident [risk]:", want exactly 1. AC#1 asks that this shape reach a user's eye at least once from a terminal; one line is the shape resident_windows.go's [hotkey] fallback has carried since ticket 258
resident_approval_risk_268_windows_test.go:205: the stdout line has to carry the loader's own answer, not just this file's class name; captured stdout was:
resident_approval_risk_268_windows_test.go:208: the stdout line has to say what the gate runs on instead; captured stdout was:
resident_approval_risk_268_windows_test.go:285: residentRiskGateValues holds 0 fmt.Printf calls, want exactly 1: AC#1's user-visible face is the other half of this branch, and a branch that only logs to disk does not reach the eye of anyone who launched the process from a terminal
```

MUT-3（`mut3-red.log:55`）：

```
resident_approval_risk_268_windows_test.go:281: residentRiskGateValues classifies by prose through [err.Error()]: matching an error message is the shape config_reload.go:396 already carries and ticket 268 was told not to copy - a reworded loader sentence would silently re-fold the two shapes
```

MUT-1 附带的一枚现场证据（**盘上日志原文**，`mut1-red.log` 里折叠态那两行）：被拒支印出
`msg="resident gate: [risk] source present but refused at construction; …" provenance="defaults (config.toml unreadable)"`
——名字与事情对不上，正是本票题面那句话的机器形状；还原后同一支印 `provenance="defaults (config.toml present but refused at load)"`。

## §6 我攻不动／判不动的地方（具名归口）

1. **双击起法下那一行看不得见**——本腿量不到。凭据面只到"终端起法可见"这一寸：
   `-H=windowsgui`（`scripts/build.ps1:115`）下无控制台可写（`cmd/wisp/console_windows.go:38`、`:49-51` 的
   `rebindStdHandle` 打不开 `CONOUT$` 即不改绑），要证伪/证实都得实机双击——⛔ 本编队不做实机动作。
   与 a2 §2③/§6 U3 同格，派单已写明"这一行看不见是已知事实，不算失败"，本腿按已知事实登记、不冒充达标。
2. **"被拒"的四种因分不开**（语法／未知键／迁移／带外值）。今天机读分不开的现量在 a2 §2②，本腿复认：
   `observe.New` 造的拒载错误 `Error.Err` 为 nil ⇒ 无 cause 可看，`ClassOf` 又同为 `ClassConfig`。
   要给每因一枚名，得往 `internal/config` 加 marker 类型＝跨包契约面，**不在本票射程**（建议归下一枚治理票；
   本腿用"级别名＋紧跟的 err 原话"承担，不造第五枚假具名）。
3. **`Q-77`（带外值到底该不该生效）**＝待机主拍板。本腿一枚断言都没放宽：被拒支仍是编译 300 s，
   并在用例①里正面钉住"不等于 20 s"。
4. **`l1_window_sec` 用钳位、`confirm_timeout_sec` 用拒载的行为不一致**（票 267 残余①）——归下一枚治理票，
   本腿没有"顺手统一"，也没动那枚钳位。
5. **票 128 AC#5 的自救句**未落地（现读 `internal/proc/envfork.go:237-239` 仍逐字
   `return Layout{}, fmt.Errorf("proc: user config dir: %w", err)`）——不由本腿补（地界在 `internal/proc`），
   本腿只做了措辞并排（§0 末段），判定不撞；若机主/128 腿后来选词撞上，请那一票复看。
6. **`wisp doctor` 的"实印枚数"口径**——本腿没跑 doctor（本轮零改动），也不在 census/a2 分歧的那三个数上另起一把尺；
   该格归口不变（a2 §2④ 判 census 的 14–16 不复现、静态枚举 11/13）。
7. **页面侧有没有承接位**（census U1）——⛔ 本编队禁读 `frontend/**`／`design/**`，量不到；
   本票落成的形（日志具名＋stdout 一行）不需要页面面，所以这一格不挡达标。
8. **`err` 透传绕开 512 截断**（a2 §3(b) 的加严）——既有仪器缺口，长 `Detail` 不截断；
   不在本票射程，本腿只复认其成立，不裁。
9. **整包名册之外的他包地界**——若交付态整包里出现 `internal/*` 的他包红，本腿只具名归口、不代修（本腿突变未种进那些包）。
10. **本腿自己的一处流程失手（如实登记）**：MUT-3 第一发还原只回了 `if` 行、留了注释块，
    靠 md5 尺抓到并当场补正（见 §5 第三行"还原证"栏）。教训具名：**还原要以字节尺收口，不能以"我改回去了"收口**。

## §7 交件判语与 commit 链

### 7.1 四格判语（本腿只判"落没落地＋自证响不响"，⛔ AC 框一枚未碰、终判归非实现者）

- **AC#1＝达标**：被拒支有一枚与"文件不存在"不同的具名状态（`riskProvenanceRefusedAtLoad`），且那一行 stdout 到达用户眼睛
  （终端起法）；"折回同一枚字符串"那一发突变**确实红**（MUT-1，红句在 §5）。
- **AC#2＝达标**：回落仍是编译 300 s；`Q-77` 未被顺手解开。
- **AC#3＝达标**：`cmd/wisp/run.go` 零字节，退码 2 那支一字未改。
- **AC#4＝达标**：禁区零字节、band 与钳位未动、无新入向方法名、四门 rc=0 且红名集合逐名比对为空集差。

### 7.2 commit 链（逐枚号＋时刻；每笔都是 `git commit -F <msgfile> -- <点名 pathspec>` 形）

| 笔 | 号 | 时刻 | 内容 | 该笔 `--name-only` 枚数 |
|---|---|---|---|---|
| 1 | `4c456d0a` | 2026-10-05 13:05:34 | 起手：§0 锚逐枚亲读＋证据件骨架八节＋票面 Progress log 首行 | 3（全在 `.scratch/wisp/`，零产码） |
| 2 | `9a941965` | 2026-10-05 13:31:53 | 落地：A 支分叉＋新常量＋B 支一行 stdout＋三枚新用例＋基线后 §0-§3 填实＋msg-02 | 5（2 枚产码＋3 枚文档面） |
| 3 | 本笔（HEAD 顶行，hash 见 `git log --oneline -1`） | 2026-10-05 15:4x（本笔动作自身时刻，非抄本表） | 交付态：§1 交付整包＋§4 四门复跑（15:25:29–15:26:39）＋§5 三发突变名册＋§6 判不动十格＋§7＋msg-03＋票面 Progress log 第 3、4 行 | 见 `git show --name-only HEAD` |

⚠ 纪律自陈：本腿**只 commit、未 push**；三笔全部带显式 pathspec（`git add` 也只 add 点名路径），
中途没有过任何裸 `git commit`（派单点名的 `3f0c4fff` 事故形状本腿避开了——a2 §0 末条的教训"显式 pathspec 只守 add 不守 commit"在本腿这里是 add＋commit 双守）。
仓内未删任何东西；MUT-3 那处的还原用的是派单许可的 `git cat-file blob HEAD:<path> > <path>`，不是 `restore`/`checkout .`。

### 7.3 本腿与派单／a2 读数不同的地方（逐条，请编排者裁）

1. **具名程度**：派单要求"字符串要自己说明发生了什么"。本腿取的是 **class 级**（`present but refused at load`）而不是
   "被值域门拒"级——依据是 a2 §2② 那句"把它直接改名成'被拒'就是新谎言"：非缺失那一支今天还含着语法／未知键／迁移／权限，
   `observe.New` 的拒载错误无 cause、`ClassOf` 同值，机读分不开。具体哪条规则由**紧跟其后的 `%v` 原话**（日志 attr ＋ stdout 行）承担。
   ⇒ 若编排者要字符串里出现"band／值域门"字样，前提是先给 `internal/config` 造 marker 类型＝**跨包契约面**，不在本票射程（§6.2）。
2. **d22scan 分母 103→104**：`ban #8 cmd/` 数的是 Go 文件枚数（含 `_test.go`），本腿新增一枚 ⇒ 分母 +1。
   本腿按"违规读数与红名集合不扩大"理解那一条纪律（违规＝0、红名＝空）；若按字面判它算"读数扩大"，
   具名退回，本腿可把三枚用例并入 `resident_approval_risk_256_windows_test.go`——但那要动别人票的钉文件，本腿默认不那么做。
3. **分档**：票面 AC#1 原文把 `doctor`／回执列为"达标档"、日志具名列为"最低档"；派单把 stdout 那一行判为达标档。
   本腿照派单交，票面原句一字未改（AC 框未碰），此处只登记分档口径的来源差。
4. **硬钉 1 的射程**：派单写"绝对不许在 `runResident` 里加第二枚 `[risk]` 读取点"，a2 §2④ 把 `:560/:563` 收窄为"构造名计数钉"
   （`:537-558` 的 AST 只数两个构造名）。本腿复认 a2 的尺**形**，同时**照派单的更严纪律执行**——一个读取点都没加，
   那一行只复用已在手的 `err`。两句话在本腿这里不冲突，冲突面留给编排者记账。
5. **票面行号一处小漂（票面 vs 派单，非代码漂）**：票面现量段写那条 `slog.Info` 在 `:379-382`，派单与本腿亲读都是 `:379-384`
   （`scope` 那一行在 `:384`）。代码没有漂，是两处文本各写了一个区间；本腿以亲读 `:379-384` 为准，票面未改。
6. **触发面比 a2 更窄一格**：a2 §4 的静态判定是"即使触发也不撞"；本腿另量到"今天整包里根本不触发"——
   `grep -rn "confirm_timeout_sec = \|l1_window_sec = " cmd/wisp/*_test.go`＝14 处种子，带外 confirm **0 枚**。
   这直接解释了交付态整包 334 全绿里没有我那一行的任何噪声。
7. **一处本腿流程失手（具名自曝）**：MUT-3 第一发还原只改回 `if` 行、把突变期写的注释块留在了文件里，
   `md5sum` 与 HEAD blob 不符（`8fcd72af…` vs `69a630bf…`）当场暴露，15:03:35 以许可的还原形补正并复量同值。
   登记在 §5/§6.10，不作隐藏。

### 7.4 交件后请编排者做的两件事（⛔ 不是本腿的活）

1. **缺口审计＋对抗验收要另一枚 agent**（D22 双角色／`SPEC-12 §4.3` #1/#3），本腿不自判终裁。
2. AC 框翻勾归编排者；本腿只在自己那一节 Progress log 追加行，票面原有句子与判据文字一枚未动。
