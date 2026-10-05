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

取数时刻 **`09:47:59–09:50:27+0800`**，HEAD 在 `89976ed1`→（并发）之间；断言行都是本腿 `Read`/`grep` 逐字符核过的原文。
四支代价列的记法：**「无损」＝这一发仍然看到它原来那件事**；「变形」＝看到的不再是同一件事（必须配套动作）；「不可达」＝这一支今天在这发上根本没有入口。

### §2.1 带外的 12 枚（这把我真正顶红的东西）

加载链先给统一结论，省 12 遍字：**三枚字面/参数种子都是"拼字符串 → `os.WriteFile` 落 `config.toml` → 由被测进程经 `config.LoadFile`/`config.NewManager` 读回来"`。校验点＝`internal/config/loader.go:137` `if err := validate(cfg); err != nil { return nil, err }`（在 `readConfigFile` 里，`LoadFile` 于 `loader.go:42` 调它），语义校验表第 31 行注册 `validateRisk(c)`（`validate.go:31`），band 判语在 `validate.go:145`。⇒ **这些用例不是"某处断言变红"，是装配根拿不到 cfg、`runTextTask` 走未配置退码 2**（`run_mode101_test.go:551` 一类 `t.Fatalf("exit %d\n%s", code, log)` 会先炸）。

| # | 种子出处 file:line | 种下去的字节（逐字） | 喂的值 | 这一发真正断言的是什么（逐字 + file:line） | 甲 抬带内 | 乙 窄缝 | 丙 fake clock | 丁 L1 窗口 |
|---|---|---|---|---|---|---|---|---|
| 1 | `cmd/wisp/run_mode101_test.go:110` | `riskLines := "[risk]\nl1_window_sec = 1\nconfirm_timeout_sec = 1\n"`（`:136` 拼进 body、`:137` `os.WriteFile(h.configPath(), []byte(body), 0o600)`） | **1** | 这族用例看的是"没人答的那张卡不得放行档"：`run_mode101_test.go:553-555` `if setErr == nil { t.Error("a console run granted auto_approve without a native click; SPEC-06:1588 says 只有原生侧点击才是\"允许\"") }`、`:557-558` `if cards != 1 { t.Errorf("cards shown = %d, want exactly one L2 card per switch (R20/M4)", cards) }`、`:560-561` `if after != risk.ModeAskEveryStep { t.Errorf("an unanswered confirmation still moved the档 to %v", after) }`、`:569-571` `result=refused-confirm`/`result=confirm-failed` 二者之一、`:577-578` `if got := h.readMode(t); got != risk.ModeAskEveryStep` | **无损**，且它已经不是我读到的那个机制：`:544` 写的是 `ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)` ⇒ 今天先响的是 **ctx**（`gate.go:576-580` `case <-ctx.Done(): … "任务上下文已结束，审批请求已作废并按拒绝处理"`），1s 的超时从来没轮到；`:569` 那句接受两个 `result=` 就是证据。代价只有一条：`t101call` 用 `context.Background()`（`:195`，无期限），所以**任何**一发没人答的卡今天 1s 就收，抬到 31s 后每次多堵 30s | 不需要（这族用例根本不碰超时） | **不可达**：`cmd/wisp/run.go:612-619` 的 `approval.Options` 没有 `Clock` 字段（现读，见 §5.1），全 `cmd/wisp` 零枚 `Clock` 引用（09:49:39 尺） | **用不上**：`:110` 那行同时种的 `l1_window_sec = 1` **已经被钳**（`gate.go:147-148` `case win < MinL1Window: win = MinL1Window` ⇒ 2s），这族用例的观察对象是 L2 卡，L1 窗口再小也替不了它 |
| 2 | `cmd/wisp/panel_pump_test.go:136` | `[]byte(string(old)+"\n[risk]\nconfirm_timeout_sec = 2\npermission_mode = \"ask_high_risk\"\n")`（`:124` 逐字目的句 `confirm_timeout_sec = 2, so the case can watch an L2 card open and close`） | **2** | 看的是"活队列的一帧快照"：`:169-171` `if rt.ui.shown() == 0 { t.Fatal("the gate never displayed a card, so the queue never held anything to report") }`、`:181-183` `if !out.IsError { t.Errorf("an unanswered L2 card must not execute, got: %s", out.Text) }`、`:189-191` `if cards != 1 {…want the one L2 card the packet reports…}`、`:192-194` `if len(live.Pending) != 1 { t.Fatalf("sampled pending = %+v, want the one item the queue held", live.Pending) }`、`:210-212` `if len(cardRecords) < 2 {…want the run's own publish plus this case's sample…}`、`:214` `assertPacketMatchesLedger145(t, cardRecords[len(cardRecords)-1], liveBytes, live)` | **变形，必须配一枚动作**：`:112` 写的是 `ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)` ⇒ 抬到 31s 后 30s 的 ctx 先响，这发就改看"上下文已结束"那张时钟，`:124` 自己声明的观察对象（"看着一张 L2 卡开合"）就丢了。⇒ 甲在这一发＝**31 + 把 `:112` 一起抬到 >31s（现读唯一一枚会漏的配套件）**，墙钟 +29s | 结构性受损：注入 gate 走 `run.go:600-608`，那里逐字 `rt.ui = nil // no console surface in a process whose cards go to another host` ⇒ 这一发读的 `rt.ui.shown()`（`:169`）与 `rt.publishPanelSnapshot()` 的宿主面没了，除非测试自带一枚能顶上 `consoleApprovalUI` 的实现 | **不可达**（同 #1），而且这一发比 #1 更需要真时钟：它要的是"整条 run 的快照泵"，假时钟只能进到 gate，进不到 `runTextTask` | **不可用**，且文件自己写了为什么：`:160-162` 逐字 `An out-of-allowlist read is R2's own example: the verdict is L2, so the item goes into the approval QUEUE (an L1 window does not), which is what snapshot.pending is made of` ⇒ L1 窗口不进队列，`live.Pending` 那格必空 |
| 3 | `cmd/wisp/approval_reply_201_test.go:202` | `h := newReplyHost(t, 20*time.Second)` → `:145` `confirm_timeout_sec = %d` + `:147` `int(h.l2Wait.Seconds()))` | **20** | 原生侧答复一张 L2 卡：`:265-267` `if row.RiskLevel != memory.RiskL2 \|\| row.Decision != agent.DecisionAllow \|\| row.Outcome != agent.OutcomeSuccess { t.Errorf("tool_call row = %+v, want L2/allow/success (the answered-card booking)", row) }`（`:225` 的 `waitForCard187(…, 5*time.Second)` 是"卡出现"的预算，不是超时） | **无损**：卡片在 `:229` `fmt.Fprintln(pw, "yes t201-allow-corr")` 被答掉，超时从不轮到 ⇒ 20→31 **墙钟零增量** | 可替代但要重写答复路（`:229` 的 reply 流是 `runSpec.reply`，AGENTS §1.3 点名的那枚缝） | 不可达 | 不适用（观察对象是 L2 答复，不是窗口） |
| 4 | `…:272` | 同 #3 的 `%d` 链 | **20** | 拒绝理由回传给模型：`:313-314` `if !out.IsError { t.Fatalf("a rejected L2 call must not execute, got outcome: %+v", out) }`、`:321-322` `if !strings.Contains(out.Text, reason) { t.Errorf("the model-visible rejection text lost the operator's reason:\n%s", out.Text) }`、`:329-330` 逐字 `"approval: ANSWER-REJECT corr=t201-reject-corr tool=fs.write decision=reject"`、`:340-341` `row.Decision != agent.DecisionReject \|\| row.ErrorClass != "user_rejected"` | **无损**（`:303` 就答了） | 同 #3 | 不可达 | 不适用 |
| 5 | `…:350` | 同链 | **20** | 面板路线只能拒：`:402-403` `if !out.IsError { t.Fatalf("a card whose allow arrived on the panel route must never execute: %+v", out) }`，四步逐句 `:384-388`（`"view " + corr`/`"panel-yes "`/`"yes "`/`"panel-no … 面板侧只能拒绝，这一发我拒了"`）。**这一发的注释替甲说话**：`:346-349` 逐字 `A longer C18 deadline here: this case answers the same card four times over … pinning that on a 2s deadline would measure the clock instead of the route` ⇒ 种子越大越接近它的本意 | **无损且更贴题** | 可，但同样是重写答复路 | 不可达 | 不适用 |
| 6 | `…:439` | 同链 | **2** ★ 全场唯一一枚"**观察对象就是超时本身**" | `:473-474` `if !out.IsError { t.Fatalf("an unanswered L2 card executed: %+v", out) }`、`:479-480` `if strings.Contains(out.Text, "用户拒绝了本次操作") { t.Errorf("a timeout must not read like a human's refusal: %q", out.Text) }`、`:484` 逐字 `"approval: ANSWER-EXPIRED corr=t201-timeout-corr tool=fs.write decision=timeout->reject"`、`:491-493` `if row.Decision != agent.DecisionTimeout { t.Errorf("tool_call row = %+v, want decision=%s (the C18 auto-reject, not an allow)", row, agent.DecisionTimeout) }`、`:461-463` 反向控制 `if rt.reply != nil { t.Errorf("this case is the detached control, but a reply surface was attached anyway") }` | **无损但慢**：这一发要的就是 `<-done`（`:467`，注释逐字 `This is the timeout path, so the wait is the C18 clock`）⇒ 31s 就是 31s，墙钟 **+29s**。**★附带一枚只有甲给得出的东西**：31s 让 `gate.go:528` 的 `lead = 31 - 30 = 1s > 0` 第一次为真 ⇒ 这一发顺手成为 r1 §7.2 说"没人能读"的那格（C18 提示端到端到达）的免费正控 | 可：`approval.Options{ApprovalTimeout: 2s}` 直接构造就把 2s 留在带外 ⇒ 但见 §4，注入 gate 会让 `:461` 的 `rt.reply` 与 `:482` 的 `h.err` 审计面全部换面 | 可，且是这一发**唯一**两全的形（见 §5.1：前提是乙先存在） | **不可用**：超时→拒绝 是 L2 的极性；L1 窗口到点是**执行**（`:545` 逐字 `the L1 timeout polarity changed (this is the D4 contract…)`）。这一发的主题就是那条极性差 |
| 7 | `…:504` | 同链 | **20** | L1 无人否决即执行：`:544-545` `if out.IsError { t.Fatalf("the L1 timeout polarity changed (this is the D4 contract, not this leg): %s", out.Text) }`、`:547-548` `if body, err := os.ReadFile(target); err != nil \|\| !strings.Contains(string(body), "unanswered")` | **无损，且这枚种子对本发是死重**：L1 走 `gate.go:291` `deadline := g.clock.After(g.window)`，`g.window` 来自 `l1_window_sec`（夹具没写 ⇒ schema tag 2 → 钳内 2s），**`confirm_timeout_sec` 在这一发根本不进时区** | 可省（本发不需要 gate 注入） | 不可达 | **已经在这条路上**：这一发观察的就是那枚合法小窗口 |
| 8 | `…:578` | 同链 | **20** | 声明了 esc 宿主的否决赶在窗口前落地：`:626-627` `if !out.IsError { t.Fatalf("a veto that landed must stop the call: %+v", out) }`、`:636` 逐字 `"approval: ANSWER-VETO corr="+corr+" tool=fs.write channel=esc decision=veto"`、`:639` 反向 `if strings.Contains(audit, "approval: ANSWER-EXPIRED corr="+corr)`（**否决必须抢在 EXPIRED 之前**） | **无损**，同 #7：那枚 EXPIRED 指的是 L1 窗口不是队列，`:579` `h.replyVeto = approval.ChannelEsc` 才是本发的主体 ⇒ 种子对本发死重 | 可省 | 不可达 | **同 #7**，这一发的正/反两控都是围绕 L1 窗口写的 |
| 9 | `cmd/wisp/approval_seam_201_test.go:50` | 同 `%d` 链 | **20** | 宿主缝在无控制台的情况下答复：`:79` `waitingStateName() = (%q,%v), want (AwaitingApproval,true) while the card is pending`、`:90` `the seam refused a native allow it must accept: %w`、`:111` `t.Fatalf("a call the host seam allowed must execute, got: %s", out.Text)`、`:129` `t.Errorf("a host answer arrived without any console, yet the console's own reply ledger booked a line:\n%s", audit)` | **无损**（卡被缝答掉，`:101` 那句 `t.Fatalf("exit %d…")` 才是 band 先炸的地方） | 可，但本发的主题恰是"装配根自己那套 reply/gate"，换注入 gate＝测的不是同一扇门 | 不可达 | 不适用 |
| 10 | `…:138` | 同链 | **20** | `PanelAllow` 必被拒：`:170` `PanelAllow must be refused on the route, got (%v,%v)`、`:179` `a grant that surfaced on the panel route must be unspendable natively, got %v`、`:200` `a card whose allow arrived on the panel route must never execute: %+v` | **无损** | 同上 | 不可达 | 不适用 |
| 11 | `cmd/wisp/ticket224_assembly_test.go:81` | 同链 | **20** | 会话动词落盘：`:146` `t.Fatalf("%d approval_grant rows for the session this run minted (%s), want %d - the card "…)`、`:166` `row created_at=%d expires_at=%d: born dead, so Covering would never honour it`（⚠ 这枚 `expires_at` 是**授权**的过期，跟 C18 超时无关）、`:191` `audit is missing the recording line naming grant_id=%d` | **无损**（`:120` 就答了卡） | 可，但 `:98` 那句 `the assembled run has no session ledger` 说明本发要的就是装配根给的 ledger | 不可达 | 不适用 |
| 12 | `…:234` | 同链 | **20** | 活授权让桥**不问**：`:269` `t.Errorf("the control call errored: the L1 window running out means execute")`、`:293` `the granted call errored: %v` | **无损**；而且这一发的控制组本来就是 L1 窗口 ⇒ 又落在丁那枚合法旋钮上 | 可省 | 不可达 | 不适用 |

⚠ **本表对派单点名处的一处补正**：派单只点了 `run_mode101_test.go:110`、`panel_pump_test.go:136`、`approval_reply_201_test.go:145` 的 6 枚调用点（`:202/:272/:350/:504/:578` 与 `:439`）。现量那枚 `%d` 的**喂值调用点是 10 枚**，多出 `approval_seam_201_test.go:50/:138` 与 `ticket224_assembly_test.go:81/:234`（尺 09:47:59）。⇒ 解冻名册少列了两枚文件。

### §2.2 带内的 15 枚（甲的**现成先例**，不是代价）

同一枚 `:145` `%d` 还被 15 枚**带内**调用点喂过 40/60/90，逐枚（09:47:59 尺）：

| 经手 helper | 调用点 file:line | 值 | 带内? |
|---|---|---|---|
| `newWidenRun`（`approval_always_201_test.go:81`） | `approval_always_201_test.go:122`、`:164`、`always_write_no_clobber_226_test.go:40` | 40 / 40 / 40 | 是 |
| `newReloadRun223`（`config_reload_223_test.go:43`） | `config_receipt_255_test.go:400`、`:468`、`:507`、`:564` | 40 ×4 | 是 |
| 同上 | `config_reload_223_test.go:247`、`:472`、`:576` | 40 ×3 | 是 |
| 同上 | `config_reload_223_test.go:300`、`:373`、`:413`、`:445` | 90 ×4 | 是 |
| 同上 | `config_reload_perm_223_windows_test.go:61` | 60 | 是 |

⇒ **"卡开了再合上"这件事，这棵树上今天已经有 15 枚用例用 ≥40s 的种子在做**；甲只需要动那 10 枚 20s/2s 的，且其中 8 枚（#3、#4、#5、#7、#8、#9、#10、#11、#12 里的 20s）答完卡就走，抬进带内**墙钟零成本**。真正的墙钟账只有两笔：`#6` +29s、`#2` +29s（外加 `#1` 在失败路径上的 30s）。

### §2.3 常驻腿那 5 枚带内字面（票 256 的钉子，确认不受影响）

`cmd/wisp/resident_approval_risk_256_windows_test.go`：`:159` `"confirm_timeout_sec = 45\n"`（`TestTicket256…CompiledConstants` 的正控）、`:202` `"confirm_timeout_sec = 45\n"`、`:221` `"confirm_timeout_sec = 45\nl1_window_sec = 2\n"`、`:286` `"confirm_timeout_sec = 45\n"`、`:293` `"confirm_timeout_sec = 90\n"`。加载链＝`writeRiskConfig256`（`:89-100`，逐字 `body := "schema_version = " + strconv.Itoa(config.SchemaVersionCurrent) + "\n"` + `[risk]`）→ `newResidentApprovalWithConfig(dir)` → `residentRiskGateValues` → `config.LoadFile`（`resident_approval_windows.go:452`）。
读面逐字：`:161` `if got := seeded.gate.Queue().Timeout(); got != 45*time.Second {`、`:251` `if got := ra.gate.Queue().Timeout(); got != c.wantTime {`、`:300` `if c.Risk.ConfirmTimeoutSec != 90 {`、`:325` `if got := fresh.gate.Queue().Timeout(); got != 90*time.Second {`。
⇒ **45/90 都在带内 ⇒ 全绿**（这正是派单里"下界不取 60"那枚钉的原因：取 60 会打红 `:159`）。
⚠ 但这一族里藏着本次唯一一枚"band 改了语义而不改颜色"的读面：`:141-144` `if got := ra.gate.Queue().Timeout(); got != approval.DefaultApprovalTimeout {…(300s is the contract default; a fourth number here means the fallback grew a value of its own)` 走的是"config.toml 不存在 ⇒ `riskProvenanceUnreadable`"。band 之后，**一份"存在但越界"的文件也会走进同一枚 `unreadable` 分支并给出 300s**（§1.1★），而这一族用例里没有任何一枚把这两种"读不到"分得开 ⇒ 具名归口 §7.3。

## §3 四支出口的代价（甲 抬带内 / 乙 测试窄缝 / 丙 fake clock / 丁 第四形）

## §4 已有的窄缝先例（乙那格：D22 注入缝名册现量）

## §5 clock 与 L1 窗口现量（丙 / 丁）

## §6 尺读数与未跑清单

## §7 判不动 / 量不到（具名归口）

## §8 我写错的读数（自我对抗）

## §9 交件判语
