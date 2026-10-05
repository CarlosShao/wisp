# 票 267 · 267-a2 只读普查：`confirm_timeout_sec` 值域带内/带外名册 + 逐枚调用点观察对象

> 本件由只读普查腿 `267-a2` 产出。硬约束：零 `go` 命令、零已跟踪文件改动、未读 `frontend/**` 与 `design/**`。
> 搜索根只用 `cmd internal tools docs .scratch scripts`。

## §0 起手锚

- 起手时刻 **`09:38:38+0800`**，起手 HEAD＝**`873c3063`**（尺＝`date "+%H:%M:%S%z"` + `git rev-parse --short HEAD`），branch `dev`。
- 骨架落盘并 commit：`09:39:16+0800` → **`8aa9ff6f`**（pathspec 只有本文件一枚，`git add .scratch/wisp/probes/267/a2/census.md` + `git commit -F .scratch/commit-msg-267a2-skeleton.txt -- <同一枚 pathspec>`；⚠ 第一次 `git commit -- pathspec` 对未跟踪文件报 `did not match any file(s) known to git`，故必须先 add，这一处是本腿的既成事实，不是别人的坑）。
- **并发情况（现量，非抄派单）**：本腿写作期间 HEAD 从 `873c3063` 漂到 **`2c83c076`**（09:44:46 与 09:44:56 两次自取一致），中间落在树上的有 `0c2445d1`「ticket 267 r1: evidence sections 1-5…」与 `492571d8`「ledger(A610)…」⇒ **`267-r1` 与编排者都在同一棵树上继续提交**。本件引用的每一枚行号都是本腿自己 grep 字符后取的，取数时刻逐把标在下面。
- 起手脏面（`git status --porcelain=v1` 前 50 行现读）：` M .gitignore`、` M .scratch/wisp/probes/152/my152.py`、` M .scratch/wisp/probes/161/r6/logs/flip-*.txt`（8 枚）、` D design/**`（15 枚）、` M design/doubao/**`（4 枚），外加一批 `??` 临时件（`.scratch/commit-msg-*.txt`、`.scratch/ci-logs/*`、`?? -`）。⛔ 全是别人的/别的编队的面，本腿一字未动、未 `checkout`、未 `clean`、未 `stash`。
- 本腿纪律遵守声明：**零枚 `go` 命令**（`go build/test/vet/gofumpt`、`wisp slo`、`sh scripts/*` 一律未跑；本件里所有"会不会红"都是**读断言 + 读加载链**的机制判语）；**未读 `frontend/**` 与 `design/**`**（`design/**` 出现在 `git status` 里是本腿如实登记别人的脏，本腿没有打开过那两棵树里任何一枚文件）；搜索根只用 `cmd internal tools docs .scratch scripts`，从未以 `.` 为根。

## §1 名册尺：全仓 `confirm_timeout_sec` / `ConfirmTimeoutSec` 出现名册

取数时刻 **`09:44:46–09:45:05+0800`**，HEAD＝`2c83c076`。尺＝`grep -rIni -e 'confirm_timeout_sec' -e 'ConfirmTimeoutSec' <根>`（大小写不敏感，两枚拼法一起；`-I` 跳二进制）。

**总量（逐根，含注释与文档）**

| 根 | 命中行 | 判语 |
|---|---|---|
| `cmd` | **21** | 全部在 `cmd/wisp/**`，产码 6 枚 + 测试 15 枚（下表逐枚） |
| `internal` | **42** | 产码 **10** 枚（`validate.go` 7 + `schema.go` 2 + `unwired.go` 1）＋ 测试 **32** 枚（`validate_267_test.go` 18、`validate_test.go` 10、`unwired_test.go` 2、`boundary_test.go` 1、`internal/agent/approval/ticket84_no_owner_test.go` 1 枚注释）；尺＝`grep -rIni … internal \| cut -d: -f1 \| sort \| uniq -c`（09:45:05） |
| `tools` | **0** | `tools/d22scan`、`tools/**` 零枚 ⇒ **没有任何仪器认识这枚键**，band 落地不会由 tools 侧红；也没有白盒钉锁着它的值域（除了票 267 自己新写的两枚） |
| `docs` | **27** | 〔仅文档〕14 枚文件，见 §1.3 |
| `scripts` | **0** | ⚠ 与派单里"scripts/*.ps1 里的提及"不一致：**现量零枚**。任何 `.ps1`/`scripts/**` 都没写过这枚键，CI 侧不会因为 band 红 |
| `.scratch` | **749** | 历史读数与票面，见 §1.4（不参与红/绿判定） |

⚠ 一把口径差：派单让我"含 `_test.go`、`.toml` 夹具与内联 fixture 字符串、`scripts/*.ps1`、`docs/**`"——本腿按**不加 `--include` 过滤**计的上面的数；早一发同根带 `--include='*.go|*.toml|*.md|*.ps1'` 的过滤尺（09:44:46）回 **213**，差值就是 `.txt`/`.log` 历史读数。**两把都真，别拿 213 当"仓里一共 213 处"用。**

### §1.1 产码名册（非 `_test.go`，全部逐枚现读）

| 出处 file:line | 逐字（截取） | 值 | 带内 [31,3600]? | 加载/执行路径 | band 落地后 |
|---|---|---|---|---|---|
| `internal/config/schema.go:460` | ``ConfirmTimeoutSec int `toml:"confirm_timeout_sec" default:"300"``` | 300 | 带内 | `default` 标签由 `internal/config/defaults.go:58-62` `NewDefaults()`→`applyDefaults()` 反射填（`defaults.go:73` `f.Tag.Lookup("default")`） | **不红**（默认值本身在带内） |
| `internal/config/validate.go:128-129` | `confirmTimeoutSecMin = 31` / `confirmTimeoutSecMax = 3600` | 31/3600 | 带的本体 | 产码常量 | 不红，它就是这场的成因 |
| `internal/config/validate.go:145-148` | `if sec := c.Risk.ConfirmTimeoutSec; sec < confirmTimeoutSecMin \|\| sec > confirmTimeoutSecMax {` + 句子 `"config.toml: risk.confirm_timeout_sec %d out of range [%d, %d]"` | — | — | `validate()`（`validate.go:24`）第 31 行注册的 `validateRisk(c)`；由 `readConfigFile` 的 `if err := validate(cfg); err != nil`（`internal/config/loader.go:137`）触发 | 不红 |
| `internal/config/unwired.go:121` | `"risk.confirm_timeout_sec": "consumed: cmd/wisp/run.go and cmd/wisp/resident_approval_windows.go both build the approval timeout from it (the resident leg since ticket 256); ticket 267 bands it [31, 3600] at load"` | 指认两枚消费点 | — | 未接线索引（文字） | 不红 |
| `cmd/wisp/run.go:616` | `ApprovalTimeout: time.Duration(cfg.Risk.ConfirmTimeoutSec) * time.Second,` | 运行时=文件里的值 | 取决于文件 | `assembleRuntime` 里 `approval.New(approval.Options{…})`（`run.go:612`），cfg 来自带校验的加载链 | 不红（消费点不校验） |
| `cmd/wisp/resident_approval_windows.go:460` | `time.Duration(c.Risk.ConfirmTimeoutSec) * time.Second,` | 同上 | 同上 | `residentRiskGateValues()`（`:447`）→ **`config.LoadFile(cfgPath, nil)`（`:452`）** | **见下面那条★，这一枚是名册里唯一的形状异常** |
| `cmd/wisp/resident_approval_windows.go:382` | `"window_sec_read", window.Seconds(), "confirm_timeout_sec_read", timeout.Seconds(),` | slog 字段名 | — | 常驻腿自述日志 | 不红 |
| `cmd/wisp/resident_approval_windows.go:308`、`:421`、`cmd/wisp/resident_windows.go:128` | 注释（`:421` 逐字 `[risk] section still answers through the schema tags (confirm_timeout_sec`） | 300（注释里指认 schema 默认） | 带内 | 注释 | 不红 |
| `cmd/wisp/firstrun.go:82`（键的字面不在这枚文件里，但它是**生产唯一写首份 config.toml 的地方**） | `if err := config.SaveFile(cfgPath, config.NewDefaults()); err != nil {` | **300（现量，非抄派单的"我记是 300"）** | **带内** | `NewDefaults()` 反射 `default:"300"`（§1.1 第一行）→ 落盘的文件里那一行是 `confirm_timeout_sec = 300` | **生产默认路径不红**。⇒ 一把干净的用户首装不会被这枚 band 拒之门外 |

★ **名册里唯一一枚"生产路径 + 带外值"的形状异常（现读，⛔ 不是"会不会红"而是"红不响"）**：`residentRiskGateValues` 走的是 `config.LoadFile`（`resident_approval_windows.go:452`），而它对带外值**返回 error**；紧接 `:453-457` 是 `if err != nil || c == nil { slog.Warn(...); return 0, 0, riskProvenanceUnreadable }` ⇒ **常驻腿不把越界当拒绝，而是当"读不到"，回落编译常量 300s 并把 provenance 写成 `"defaults (config.toml unreadable)"`（`:434`）**。
⇒ 后果两条，都归编排者：
1. 裁形 ⓐ 那句"越界＝加载时拒"在**跑任务那条腿**（`run.go` 的装配根）成立，在**常驻那条腿**今天成立的方式是"拒读 + 回落 300 + 一行 slog.Warn"。这**不是坏事**（回落值 300 让 C18 提示照样发得出去），但它意味着"同一份越界 config 在两条腿里得到两种答复"，而 `internal/config/unwired.go:121` 那句新写的指认文本**只说了 bands it at load，没说常驻腿其实会静默回落**。票面与 §2 都不归本腿改，具名归口 §7。
2. **名册里没有任何一枚生产成员种带外值**：schema 默认 300 带内、`run.go`/`resident_*.go` 只读不写。⇒ 带外种子**全部住在 `_test.go` 与历史读数里**。

### §1.2 测试名册（`_test.go` 与内联 fixture 字符串）

| 出处 file:line | 逐字 | 值 | 带内? | 加载链 | band 后 |
|---|---|---|---|---|---|
| `internal/config/validate_test.go:192` | `for _, bad := range []int{30, 10, 3601, 99999, 0, -1} {`（`:197` `c.Risk.ConfirmTimeoutSec = bad`，`:198` `err := validate(c)`） | 30/10/3601/99999/0/-1 | 带外（**期望被拒**） | 直接调 `validate()`（包内测试） | **绿**（它就是 band 的正面钉子） |
| `internal/config/validate_test.go:216` | `for _, good := range []int{31, 300, 3600} {` | 31/300/3600 | 带内 | 同上 | 绿 |
| `internal/config/validate_test.go:236`、`:240` | `c.Risk.ConfirmTimeoutSec = 31` / `= 30` | 31 带内 / 30 带外 | 一内一外 | 同上 | 绿 |
| `internal/config/validate_267_test.go:47` | `body := fmt.Sprintf("schema_version = 2\n\n[risk]\nconfirm_timeout_sec = %d\n", sec)` | 由 `:95` 的表给 | — | `loadWithTimeout`→真加载链 | 绿 |
| `internal/config/validate_267_test.go:95` | `for _, bad := range []int{lead, 10, 1, 2, 20, 3601, 99999, 0, -1} {`（`lead := warningLeadSec(t)`＝30） | **30/10/1/2/20**/3601/99999/0/-1 | 带外（**期望被拒**） | 真加载链（`LoadFile`/`NewManager`） | **绿，且它已经把 cmd/wisp 今天的 1/2/20 钉成"必须拒载"** ⇒ `internal/config` 与 `cmd/wisp` 现在对同一枚数持**相反期望**，这就是这把我顶红的机制根 |
| `internal/config/validate_267_test.go:105` | `for _, good := range []int{lead + 1, 300, 3600} {` | 31/300/3600 | 带内 | 同上 | 绿 |
| `internal/config/unwired_test.go:66` | `confirm_timeout_sec = 300`（夹具全文内联，`:58` `LoadFile(path, nil)`） | 300 | 带内 | `LoadFile`→`readConfigFile`→`validate` | **绿** |
| `internal/config/unwired_test.go:188` | `confirm_timeout_sec = 60`（同文件另一发夹具，挨着 `l1_window_sec = 5`） | 60 | 带内 | 同上 | **绿** |
| `internal/config/boundary_test.go:142` | `c.Risk.ConfirmTimeoutSec = 60` | 60 | 带内 | 结构体赋值（不走 TOML 夹具） | **绿** |
| `internal/agent/approval/ticket84_no_owner_test.go:51` | `// same knob cmd/wisp/run.go:259 feeds from [risk].confirm_timeout_sec - it is` | 注释 | — | 注释 | 绿（⚠ 注释里指的 `run.go:259` 已漂到 `:616`，**这是别人票上的过期行号引用，不是本腿的读数**，见 §8） |
| `cmd/wisp/run_mode101_test.go:110` | `riskLines := "[risk]\nl1_window_sec = 1\nconfirm_timeout_sec = 1\n"` | **1** | **带外** | 拼进 `:114` 的 `body`，`os.WriteFile` 落 `config.toml`；由 `runTextTask`→`assembleRuntime` 的加载链与 `:172` `config.NewManager(h.configPath(), nil)` 读 | ★ **会红**（§2 逐枚） |
| `cmd/wisp/panel_pump_test.go:136` | `[]byte(string(old)+"\n[risk]\nconfirm_timeout_sec = 2\npermission_mode = \"ask_high_risk\"\n")` | **2** | **带外** | 追加到 fixture 的 `config.toml` | ★ **会红** |
| `cmd/wisp/approval_reply_201_test.go:145` + `:147` | `confirm_timeout_sec = %d` ＋ `int(h.l2Wait.Seconds()))` | **由调用点给**：20 / 2 带外，40 / 60 / 90 带内 | 一半带外 | `fmt.Sprintf` 拼整份 `config.toml`，`:148` `os.WriteFile` | ★ **10 枚调用点会红**（逐枚在 §2）；带内那 15 枚经两枚 helper 转手，见 §2.2 |
| `cmd/wisp/approval_reply_201_test.go:105-109` | `if l2Wait < time.Second { … l2Wait = time.Second }` | 兜底 **1** | **带外**（今天没有调用点走到这一支） | 同上 | 不直接红，但 band 后这句兜底**自己变成一个会拒载的值**＝潜在地雷，具名进 §3 甲 |
| `cmd/wisp/resident_approval_risk_256_windows_test.go:159`、`:202`、`:221`、`:286` | `writeRiskConfig256(t, dir, "confirm_timeout_sec = 45\n")`（`:221` 是 `"confirm_timeout_sec = 45\nl1_window_sec = 2\n"`） | **45** | **带内** | `:89-100` helper 写 `schema_version` + `[risk]` 两行；读侧 `newResidentApprovalWithConfig`→`config.LoadFile` | **绿**（派单里"45s 那枚钉"确认不受影响） |
| `cmd/wisp/resident_approval_risk_256_windows_test.go:293` | `writeRiskConfig256(t, dir, "confirm_timeout_sec = 90\n")` | **90** | **带内** | 同上 + `:296` `config.LoadFile(filepath.Join(dir, configFileName), nil)` 复核 | **绿** |
| `cmd/wisp/resident_approval_risk_256_windows_test.go:141`、`:161`、`:251`、`:288`、`:300-303`、`:310`、`:325` | 读面 `ra.gate.Queue().Timeout()` 与 `c.Risk.ConfirmTimeoutSec != 90` | 断言侧 | — | 接缝构造后的对象 | 绿（45/90 都在带内） |

### §1.3 文档名册（〔仅文档〕，band 不会让它红，但会让**文字过期**）

`docs` 27 枚的逐文件分布（09:45:05 现量）：`PLAN.md` 1／`SPEC-03-config-secrets-envs.md` 1／`reports/missing-features-2026-09-29-v4.md` 1／`reports/pending-and-issues.md` 8／`evidence/s1/` 14 枚文件共 15（`246-resident-task-source-v2.md` 3、`248-settings-write-path-r1.md` 3、`35-panel-snapshot-pump-r1.md` 2，其余 9 枚文件各 1）。
逐枚里**带数字**的只有这些，其余是引用与判语：

| 出处 | 写的值 | 判语 |
|---|---|---|
| `docs/PLAN.md:2738` | `` `confirm_timeout_sec`(300) `` | 〔仅文档〕300 带内，**不过期**（D36 的 section 树） |
| `docs/specs/SPEC-03-config-secrets-envs.md:34` | `` `confirm_timeout_sec(int)=300` `` | 〔仅文档〕带内，**但 SPEC-03 没有写任何值域** ⇒ band 是一条 SPEC-03 里没有的新规矩（⛔ 本腿不裁要不要补，归口 §7） |
| `docs/reports/missing-features-2026-09-29-v4.md:194` | `schema.go:450` `confirm_timeout_sec` 300 | 〔仅文档〕值不红，**行号红**：现读该字段在 `schema.go:460` |
| `docs/evidence/s1/198-firstrun-config-v1.md:34`、`docs/evidence/s1/90-adversarial-acceptance.md:272` | `default:"300"` 实测 | 〔仅文档〕带内；后者还钉着"原文一行 `…toml:"confirm_timeout_sec" default:"300"`"，本腿现读字符一致 |
| `docs/evidence/s1/128-ac4-r1-acceptance.md:425` | `run_mode101_test.go:309`（配置里写着 `confirm_timeout_sec = 1` 却被默认值取代） | 〔仅文档〕历史红因，**不红今天**；但它是仓里**已经出现过一次"种子没被读到"**的案底，与 §2 第 1 行直接相关 |
| `docs/evidence/s1/35-panel-snapshot-pump-r1.md:265`、`:454` | `panel_pump_test.go:133` 只追加 `confirm_timeout_sec = 2` | 〔仅文档〕**行号已漂**（现读那一行在 `:136`），值与形状一致 |
| `docs/evidence/s1/140-ac2-static-blast-radius-r1.md:458` | 逐字抄了 `run_mode101_test.go:110` 的 `riskLines := "[risk]\nl1_window_sec = 1\nconfirm_timeout_sec = 1\n"` | 〔仅文档〕本腿逐字符核过**还在**（§2 同一枚） |
| `docs/evidence/s1/246-resident-task-source-v2.md:112`、`:212`、`246-resident-task-source-r2.md:85`、`248-settings-write-path-r1.md:301-305` | 常驻腿"吃不到 `confirm_timeout_sec`／走常量 300s" | 〔仅文档〕**这批文字已被票 256 推翻**（`:460` 现在真读 config），band 不会让它们红，但它们是台账里"过期判语"的存量，见 §8 第 3 条 |
| `docs/reports/pending-and-issues.md:2648`、`:6264`、`:6291`、`:10107`、`:11761`、`:11829`、`:11910` | 台账（`Q-77`、票 248 AC#10 等） | 〔仅文档〕⛔ 本腿不碰台账（派单规矩 2） |

### §1.4 `.scratch` 那 749 枚（不参与红/绿）

桶（09:45:05 现量）：`.scratch/wisp/probes=717`（其中**本票自己的 `probes/267=83**，`267/r1/evidence.md` 一枚就 42）、`.scratch/wisp/issues=21`（票 267 面 7 / 256 面 3 / 83 面 4 / 248 面 2 / 80、246、84、90 各 1）、`.scratch/ci-logs=4`、`.scratch/wisp/ledger-restore.tmp=4`、`.scratch/probes-255-r1c=2`（那是 `schema.go` 的**快照副本**，`:450` 那一行）。
⚠ 一把**必须说清的假红**：`.scratch/ci-logs/run-37166458550-failed.log:191`、`:425` 里有 `confirm_timeout_sec = 40`（带内），`run-37158259050.log:2775` 与大批 `probes/*/gate-*.txt`（票 145/147/149/151/152/154/156/158/179/198/200/226 等）里有 `confirm_timeout_sec = 1`——**那是 `go test -v` 打印夹具原文的历史尸体**，不是活代码，band 不会让它们红；但它们是本腿判断"cmd/wisp 曾经全绿过、种子确实是 1/2/20"的**盘上旁证**（⛔ 不是本腿的读数，本腿的红绿判语一律来自 §1.1/§1.2 的源码现读）。

### §1.5 名册尺的直接回答

- **这把我顶红了几枚包？** 现读的带外种子只住在 **`cmd/wisp` 一枚包**里，但**不止派单说的三枚文件**：
  - `run_mode101_test.go:110`（字面 **1**）＋ `panel_pump_test.go:136`（字面 **2**）
  - `approval_reply_201_test.go:145` 那一枚 `%d` 由 **10 枚调用点**喂 20 或 2 ⇒ 这 10 枚分布在 **3 枚文件**：`approval_reply_201_test.go:202/:272/:350/:439/:504/:578`、`approval_seam_201_test.go:50/:138`、`ticket224_assembly_test.go:81/:234`（后两枚文件**没出现在派单点名的三枚里**，是本次普查新增的具名项）
  - ⇒ **合计 1 枚包 / 5 枚文件 / 12 枚带外位置**（尺 09:47:59，`grep -rn "newReplyHost(t, " cmd/wisp`）。
  - 另有 **15 枚带内位置**经 `newWidenRun`（`approval_always_201_test.go:122/:164`、`always_write_no_clobber_226_test.go:40`）与 `newReloadRun223`（`config_receipt_255_test.go:400/:468/:507/:564`、`config_reload_223_test.go:247/:300/:373/:413/:445/:472/:576`、`config_reload_perm_223_windows_test.go:61`）喂 40/60/90，**全部带内 ⇒ 不红**。
  `internal/config`（含票 267 自己的钉子）、`internal/agent/approval`、`tools`、`scripts` 四枚包**一枚不带外种子** ⇒ 不红。
- **生产路径有没有种带外值的？** **没有**（§1.1 末格：默认表 300，`firstrun.go:82` 只是把 `NewDefaults()` 落盘）。带外值 100% 是测试故意种的，这是"出口选择"整件事的支点。

## §2 逐枚调用点表（种子出处 / 逐字节 / 加载链 / 这一发真正断言的是什么）

## §3 四支出口的代价（甲 抬带内 / 乙 测试窄缝 / 丙 fake clock / 丁 第四形）

## §4 已有的窄缝先例（乙那格：D22 注入缝名册现量）

## §5 clock 与 L1 窗口现量（丙 / 丁）

## §6 尺读数与未跑清单

## §7 判不动 / 量不到（具名归口）

## §8 我写错的读数（自我对抗）

## §9 交件判语
