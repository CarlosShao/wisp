# 票 268 对抗验收腿 `268-v1`：证据件（非实现者终裁 AC#1–AC#4）

- 验腿＝`268-v1`，被验对象＝`268-r1` 的交付（产码 commit `9a941965`＋交付笔 `b93624d4`／`2a0e23b4`／`2925b643`）。
- 本件只落 `.scratch/wisp/probes/268/v1/**`；⛔ 零产码、零票面、零台账、零 docs/。
- 起手 HEAD：`2a633eb8`（`git rev-parse HEAD` 现量）。
- 起手时刻：`2026-10-05 20:00:05 +0800`（`date` 现量）。
- 起手工作树读数：`git status --porcelain` 全仓 **724** 行；scoped `cmd/wisp` **0** 行。
- 突变窗规划：MUT-1/2/3 全部种在 `cmd/wisp/resident_approval_windows.go`，每发后 `git cat-file blob HEAD:cmd/wisp/resident_approval_windows.go > cmd/wisp/resident_approval_windows.go` 还原＋md5 双读收口。
- ⚠ 起手时 CPU 93%／MEM 80.7%（20:04:47 现量）⇒ 按派单纪律**整包暂不起**，先做只读面；整包起跑前必须 CPU/MEM < 70%。

## §0 产码亲读（被验对象的第一手读数）

### 0.1 被验对象是否完整到达本腿

- 两枚产码自 `9a941965` 之后零笔再触：`git log --oneline -1 -- <path>` 两枚都停在 `9a941965`（20:1x 现量）。
- 工作树 vs HEAD blob 逐字节同值（md5 双读，20:1x）：
  - `cmd/wisp/resident_approval_windows.go`：`69a630bf350daf1032c62b921e9d7044`（工作树 `md5sum`＝`git cat-file blob HEAD:…|md5sum` 同值，与 r1 §4/§6 末采报的 `69a630bf…` 同值）
  - `cmd/wisp/resident_approval_risk_268_windows_test.go`：`e2405a21bb278b558c1305ea95d34a9b`（同法，与 r1 报的 `e2405a21…` 同值）
- ⇒ 本腿的全部静态读数（§0/§1/§4.5/§5）都是对**交付态字节面**取的；「门读数与 r1 末笔不一致时先怀疑突变窗」在本腿这里的对应物＝「静态读数与 r1 报数不一致时先怀疑字节面漂移」——已排除。

### 0.2 产码 diff 亲读（`git show 9a941965 -- cmd/wisp/resident_approval_windows.go`，20:0x）

- import 只加一枚 `io/fs`；`errors`／`fmt`／`log/slog` 本来就在 ⇒ r1 §2.1 第 4 条「无新增包级依赖边」复现。
- 分派结构亲读：`residentRiskGateValues` 里 `err != nil || c == nil` 支内第一层 `if errors.Is(err, fs.ErrNotExist)`（`:476`）走
  原 `slog.Warn` ＋ `return 0, 0, riskProvenanceUnreadable`（`:480`）；非缺失支走新 `slog.Warn`（msg 含 "present but refused at construction"，
  多带一枚 `"provenance"` attr）＋ 一行 `fmt.Printf`（`:510`，两行字符串拼接）＋ `return 0, 0, riskProvenanceRefusedAtLoad`（`:513`）。
- 四枚 provenance 常量现读：`riskProvenanceRead = "config"`（`:439`）／`riskProvenanceUnreadable = "defaults (config.toml unreadable)"`（`:446`，
  **逐字未动**）／`riskProvenanceRefusedAtLoad = "defaults (config.toml present but refused at load)"`（`:456`，新）／`riskProvenanceNoView = "defaults (no host config view)"`（`:460`）。
- 文件 EOL 亲读（`git show 9a941965:… | tail -c 32 | od -c`）：末字节 `} \n`，LF 收尾、**无 CR 残留、无重复收尾**。
- 函数体尾随行亲读（sed 375-390）：`slog.Info("resident gate: [risk] tier taken at construction"` 的六枚 attr 中
  `"scope", …` 在新文件 `:385`；对照旧文件（`9a941965^`，`87bc8b0f`）同一句 `:379`、scope `:384` ⇒ **该句在新文件整体下移 1 行**
  （diff 的 72 增/8 删里删的是 `:410-411` 附近的旧注释两行带）。

### 0.3 测试文件亲读（293 行，20:0x）

- 三枚顶层用例名与 r1 §2.3 表逐枚一致：① `TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig`、
  ② `TestTicket268RefusedRiskConfigReachesStdoutOnce`、③ `TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords`。
- ① 的断言面亲读：缺失/被拒两枚 provenance 各自对表＋互斥；四常量两两不同字面（map 逐对检重）；两枚回落形态都断
  `Queue().Timeout() == approval.DefaultApprovalTimeout` 且 `!= 20*time.Second`（Q-77 未动）；`Window() == approval.DefaultL1Window`；
  前提钉四连（`os.Stat` 在／`config.LoadFile` 报错／`!errors.Is(loadErr, fs.ErrNotExist)`／`strings.Contains(loadErr.Error(), "out of range [31, 3600]")`——
  **测试侧允许文案匹配**，注释里写明 "A test may match prose (production may not; ruler ③ owns that half)"）；同目录正控（修成 45 ⇒ `riskProvenanceRead` ＋ 45s）。
- ② 的断言面亲读：`captureStdout128` 包住构造；`strings.Count(out, "wisp: resident [risk]:")` **恰 1**；行内含
  `out of range [31, 3600]` 与 `DefaultApprovalTimeout=300s`；反控＝缺失支同一捕获面 **零行**。
- ③ 的 AST 尺亲读：`parser.ParseFile` 整文件 → 找 `residentRiskGateValues` → `errors.Is(…, ErrNotExist)` **恰 1**、
  `strings.*`／`err.Error()` **0**、`fmt.Printf` **恰 1**、四枚常量名各有一次 `return` 走到（`returned` map 全真）。
  射程只在 `residentRiskGateValues` 一个 FuncDecl 内，注释不可见。
- 复用件亲读：`risk256Dir`／`writeRiskConfig256` 来自 256 文件（本腿未重造）；`captureStdout128` 来自 `dataroot_128_test.go`。
- ⚠ EOL 疑点（自报不裁）：`git show` 渲染的测试文件 diff 里 `1\r` `2\r` 行号串出现 `^M` 形状，而产码同渲染干净；
  但**工作树字节面**上两枚 md5 都与 HEAD blob 同值 ⇒ 至少「r1 交付的字节」与「仓里 HEAD 的字节」一致；
  `1\r` 是否 diff 渲染伪影，本腿不再追（不影响任何判语，仅登记）。

### 0.4 票面锚行号亲读（六处差异之一在此定事实）

- **`:379-382` vs `:379-384` 的第三方读数**：票面「现量」段写那条 `slog.Info` 在 `:379-382`、派单与 r1 亲读都写 `:379-384`。
  本腿亲读两处：**旧文件**（`9a941965^`＝`87bc8b0f`，票面锚定的就是它）`slog.Info("resident gate: [risk] tier taken at construction"` 在 **`:379`**、
  六枚 attr 的最后一枚 `"scope", …` 在 **`:384`** ⇒ 该句**实际跨 `:379-384`**（`:382` 只是 `"window_sec_read"`/`"confirm_timeout_sec_read"` 那一行，
  **不是句尾**）。⇒ 事实＝票面区间截短（少两行）、派单与 r1 的 `:379-384` 与旧文件实况一致；**代码无漂**（`9a941965^` 的 `:379-384` 逐字未动）。
- 现文件（`9a941965` 之后）该句下移为 `:380-385`（`:380` msg、`:385` scope）。任何新票面引用若拿新文件行号，应以 `:380-385` 为准。
- 其余锚（本腿 20:1x 亲读）：`:434`→现 `:446`（同因下移；常量字面逐字在）、`:457`→现 `:480`、`:454-456` 缺失支 `slog.Warn` 现在跨 `:477-479`、
  `resident_windows.go:187` 的 `[hotkey]` `fmt.Printf` 母本行未漂（亲读在）、`256_windows_test.go:141-144` 的 300s 钉未漂（亲读在，
  注释逐字含 "a fourth number here means the fallback grew a value of its own"）。

## §1 起手锚与四格判语

「待验」（本节在整包/突变完成后填实并 commit；判语先立骨架）：

| 格 | v1 判语 | 依据节 |
|---|---|---|
| AC#1 | 「待验」 | §0.2/§0.3/§2 |
| AC#2 | 「待验」 | §0.3/§2 |
| AC#3 | **达标**（v1 判，静态面；§5.4） | §5.1/§5.4/§4.5 |
| AC#4 | **达标**（v1 判，静态面；§5.5） | §5.5/§4 |

## §2 突变名册（自重种三发）

「待验」

## §3 名册差集（v1b 复跑名册 vs 268-r1 vs probes/257/v1）——**不重跑，从三枚现成文件算**

派单明令 ⛔ 不重跑整包。本节全部读数出自这三枚已在盘的日志，尺逐枚同一条：
`grep -E '^--- (PASS|FAIL|SKIP)' <log> | sed -E 's/^--- ([A-Z]+): ([^ ]+).*/\1: \2/'`（顶层名册，缩进子枚不计）。

| 件 | 路径 | 行尺读数 |
|---|---|---|
| v1 死腿的整包（编排者代提 `164ee2c2`） | `.scratch/wisp/probes/268/v1/deliver-cmdwisp.log` | `wc -l`＝**1,723**；`grep -c ''`＝**1,723**（末字节 `\n`，`od -c` 亲读 `…422.320s\n`）；末三行＝`PASS` / `ok … 422.320s` |
| 268-r1 交付态整包 | `.scratch/wisp/probes/268/r1/deliver-cmdwisp.log` | `wc -l`＝**1,724**；`grep -c ''`＝**1,724**；末两行＝`ok … 401.239s` / `rc=0`（**多出的那一行＝腿自己 echo 的 `rc=0` 注脚，不是测试输出**——`diff` 定位在 `1723,1724c1723`，唯一一处尾部差） |
| 268-r1 基线整包（动笔前） | `.scratch/wisp/probes/268/r1/baseline-cmdwisp.log` | 1,708 行 |

**名册枚数（同一把尺，三发一致到底）**：

| 件 | `=== RUN` 顶层 | `^--- PASS` 顶层 | `    --- PASS` 子枚 | 全量 `--- PASS` | FAIL | SKIP | `ok` |
|---|---|---|---|---|---|---|---|
| r1 基线（13:0x） | 331 | 235 | 96 | **331** | 0 | 0 | 466.320s |
| r1 交付（15:12） | 334 | 238 | 96 | **334** | 0 | 0 | 401.239s |
| v1 死腿（20:15） | 334 | 238 | 96 | **334** | 0 | 0 | 422.320s |

⇒ 派单引文「`ok … 422.320s`、末枚 PASS 可见」复现：末枚顶层 PASS＝`TestTicket224ProductionSessionDoesNotSurviveRestart (3.87s)`（20:15:38，`tail -6` 亲读）。
⚠ **r1 那句「RUN 334／PASS 334」不是漂字**：334＝全量 `--- PASS`（顶层 238 ＋ 缩进子枚 96）；`=== RUN` 顶层也恰 334。三发同尺同形，本腿无退回。

**差集（三枚名册去前缀、排序、`comm`）**：

| 对 | 读数 |
|---|---|
| v1b(238) vs r1 交付(238) | `diff`＝**空**（逐名同集同序，`diff` 顶层 RUN/PASS 序列亦**空**） |
| v1b(238) vs r1 基线(235) | 多 **3 枚**、少 **0 枚**：`TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig` / `TestTicket268RefusedRiskConfigReachesStdoutOnce` / `TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords`（与 §0.3 亲读的三枚用例名逐枚同） |
| v1b(238) vs `257/v1/names-mine.txt`(235) | 多同 **3 枚**、少 **0 枚** |
| v1b vs `257/v1/names-theirs.txt`(235) | 多同 **3 枚**、少 **0 枚** |
| `257/v1/names-mine` vs `r1/baseline-top-names` | `diff`＝**空**（235 枚逐名同集＝今天绿名册的第三方独立复认） |
| `257/v1/roster-mine.txt` 名集 vs v1b | 少 **0 枚** |

⇒ **名册差集判语＝成立**：v1 死腿那次全绿整包与 r1 交付态名册**逐名同集**，268 家族净增恰 3 枚、一枚未少、一枚未红；
257/v1 那枚**同日、他票、独立跑**的名册与 r1 基线名册同集，等于给"这 235 枚是今日绿底"补了第三方读数。
⚠ `257/v1/roster-mine.txt` 首行是 `FAIL: TestAC14GoSideEvalPushReachesThePage`——那是 **257 腿自己的窗口里**的读数，
本尺只取名集；该用例在 r1 基线／r1 交付／v1 交付三发里逐枚 `--- PASS`（v1 日志 `:838`，0.48s），不构成任何未少证据。

**stdout 那行的跨发计数**（尺＝r1 §7.2 用的 `grep -c 'wisp: resident \[risk\]'`）：
基线 **0** ／ r1 交付 **1** ／ v1 死腿交付 **1**——r1 报的同值复现。
v1 那一枚出现处 `:1002`，`=== RUN` 归属＝`TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig`（20:11:50.962，`:998` 起）；
该行紧邻上一行是 `slog.Warn` 的 `provenance="defaults (config.toml present but refused at load)"`，紧邻下一行是 `slog.Info` 同一枚 provenance ⇒ 被拒支三面对齐（日志／stdout／构造回执）。

## §4 门禁四数

「待验」

## §5 攻它没攻的格

### 5.1 「runResident 不得加第二枚 [risk] 读取点」的硬钉审计（r1 只报了"没有"，本腿枚了全部读取点）

现量 `grep -n "config.LoadFile" cmd/wisp/*.go | grep -v _test`（20:2x）＝**7 枚产码调用点**：

| 站点 | 归属票 | 268 后状态 |
|---|---|---|
| `cmd/wisp/models.go:184` | 255 族（任务侧 geometry/providers 前身） | 未动 |
| `cmd/wisp/panel_resident_windows.go:201` | 255 AC#4（panel geometry per-use 读） | 未动 |
| `cmd/wisp/providers.go:98` | 任务腿 | 未动 |
| `cmd/wisp/resident_approval_windows.go:474` | 256 建的那枚（268 在**其内部**分叉） | 268 改的是函数体、没加新调用点 |
| `cmd/wisp/resident_windows.go:183` | 258 AC#1（[hotkey] hotCfg258 构造读） | 未动 |
| `cmd/wisp/resident_windows.go:206` | 258（[hotkey] hotReload258 per-tick 读） | 未动 |
| `cmd/wisp/run.go:17` | 注释行（非调用） | 未动 |

⇒ **268-r1 没有加第七枚读取点**；它把被拒分支写进 256 已建的那枚函数里，`runResident` 调用面
（`newResidentApprovalWithConfig(rt.Layout.DataDir)`，`resident_windows.go:122` 附近）零字节。
256 的 `:559-568` 构造名计数钉（`fed != 1`／`bare != 0`／`argOK`）亲读在册、未被触碰。
派单的「硬钉 1」（不在 runResident 里加第二枚 [risk] 读取点）**满足**；与 a2 §2④ 收窄口径的冲突不存在——两条尺在这份 diff 上同读数。

### 5.2 `l1_window_sec` 钳位 vs `confirm_timeout_sec` 拒载（票 267 残余①）——亲读

- `internal/config/schema.go:450-463`（票 267 注释）逐字：confirm **"gated at load, not clamped at the gate"**；
  `l1_window_sec` **"deliberately keeps its consumer-side clamp (gate.go MinL1Window/MaxL1Window)"**——
  即"行为不一致"是**注释里写明的 deliberate**，不是 268 造的，也不是 268 该顺手统一的（AC#3 同理）。268-r1 未动 schema.go（§4.5）。
- `internal/agent/approval/` 最后一笔＝`789a02e2`（票 259-r2，12:52:43），268 全链零触；gate.go 的 MinL1Window/MaxL1Window 未动。

### 5.3 「常驻腿任何时候按用户写的数跑」须同时满足两枚正控——枚齐

本腿枚到的正控全在：256 的 `:151-159`（45s 正控）与 `:240-260`（45/90/31/40 种子表，provenance 必须 `riskProvenanceRead`）；
268 用例①的 45s 正控。负控（缺失/被拒/NoView 三回落形态 ⇒ 300s/3s）在 256 ruler①＋268 用例①。**没有一枚形态缺正负控。**

### 5.4 AC#3「两腿不许顺手统一」的逐文件形审（r1 只报了"run.go 零字节"）

- `run.go` 只含 `config.LoadFile` 的**注释**（`:17`），产码调用零枚；268 diff 名册里 `run.go` 零字节。
- 跑任务腿的响亮拒绝（退码 2＋"配置未就绪（Unconfigured）"句）不在本腿写面，`git show --name-only 9a941965` 全 5 枚路径无它。
- 常驻腿缺失支**没有**因为本票而获得任何新可见面（268 用例②的反控钉着：缺失支 stdout 零行）；
  被拒支的新面只属于被拒支。⇒ 「一响一静」的差别没有被抹平：跑任务腿=拒载即退；常驻腿=拒载回落＋具名＋stdout；缺失=回落＋具名（无 stdout）。

### 5.5 AC#4 逐条（静态面）

- **禁区命中**：`git show --name-only` 对 `9a941965`（5 枚路径）与 `b93624d4`（11 枚路径）分别过禁区正则
  （docs/PLAN.md、docs/specs/、internal/observe/thresholds.go、internal/agent/approval/、tools/d22scan/allowlist.txt、frontend/、design/）
  ＝ **两笔各 0 行命中**（20:2x 现量）。r1 §4.5 同尺同读数（0 行）。
- **新增导出方法名**：全 diff `^\+` 过 `func \(…\) [A-Z]`＝**0 枚**；产码文件 `^\+(func |	func )`＝**0 枚**；
  测试文件新增 `func` 仅三枚 `func Test…`（测试面，非入向 API）⇒ **零新增入向方法名（C17 面未动）**，r1 读数复现。
- **值域 [31,3600]／l1_window_sec 钳位一字未动**：`git diff 9a941965^ 9a941965 --stat -- internal/`＝**空**（internal/ 零字节）；
  `internal/config/validate.go` 现量 `:128` `confirmTimeoutSecMin = 31`、`:129` `confirmTimeoutSecMax = 3600`、`:145` band 判断逐字在；
  该文件最后一笔＝`37f8e5c6`（票 267-r1）——268 全链零触。
- **`riskProvenanceUnreadable` 字面逐字未动**：现量 `:446`＝`"defaults (config.toml unreadable)"`，与票面引文逐字同。
- **测试种子只种带外值（20），不改 band**：`grep` 全 diff 中带 `31/3600` 的新增行全部是**测试断言文案／证据件引文**，
  零枚落在 `internal/config`。

## §6 六处差异的独立复核

「待验」

## §7 commit 链与收尾读数

「待验」
