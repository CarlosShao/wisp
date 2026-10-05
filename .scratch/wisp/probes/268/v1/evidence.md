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
| AC#3 | 「待验」 | §5 |
| AC#4 | 「待验」 | §4.5/§4 |

## §2 突变名册（自重种三发）

「待验」

## §3 名册差集（v1 复跑名册 vs 268-r1 vs probes/257/v1）

「待验」

## §4 门禁四数

「待验」

## §5 攻它没攻的格

「待验」

## §6 六处差异的独立复核

「待验」

## §7 commit 链与收尾读数

「待验」
