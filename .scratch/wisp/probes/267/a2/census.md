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

尺：本节全部依据 §1/§2 的现读位置，取数时刻 09:44:46–09:53:56+0800，HEAD 从 `2c83c076` 漂到 `9ca14bdf`（本腿自己的两枚 commit 之外还有别人的）。

### 甲 · 把 12 枚带外位置抬进带内（≥31s），并给每一发找回观察对象

| 代价项 | 现量 |
|---|---|
| 要改的位置 | **12 枚**（§2.1 全表）：2 枚字面 + 10 枚喂值调用点，分布在 **5 枚文件** |
| 墙钟增量 | **只有两笔**：`approval_reply_201_test.go:439`（2→31，这一发真的要等 C18 时钟走完，`:467` 注释逐字 `This is the timeout path, so the wait is the C18 clock`）与 `panel_pump_test.go:136`（2→31，`:177` 的 `<-done` 在等卡合上）⇒ **+29s +29s**。其余 8 枚答完卡就走 ⇒ **零增量** |
| 必须连带抬的一件 | `panel_pump_test.go:112` `context.WithTimeout(context.Background(), 30*time.Second)`：30 < 31 ⇒ 不抬它就是**换时钟**（`gate.go:576-580` 的 ctx 分支抢在 `:571` 的 deadline 前）。这一枚是本腿能给出的唯一"改了种子却没找回观察对象"的位置 |
| 必须连带改的一句兜底 | `approval_reply_201_test.go:105-109` `if l2Wait < time.Second { … l2Wait = time.Second }`：band 后这枚兜底自己写出的就是**一个会被拒载的值**（1）。今天 10 枚调用点没有一枚走到它，所以它不红，但它是留在那儿的一枚雷；甲必须把它抬到带内（或改成具名报错） |
| 可以顺手删的一枚 | `run_mode101_test.go:110` 的 `confirm_timeout_sec = 1`：§2 #1 已证这一族的观察走 300ms ctx（`:544`），**1s 从不轮到**。删掉它＝落 schema 默认 300（带内，不红），代价是"卡被忘答时堵 300s 而不是 1s"；写 31 是更稳的形状 |
| ★甲独有的收益 | 抬到 31 之后 `lead = 31 − 30 = 1s > 0` ⇒ `gate.go:528` 那句第一次为真，`gate.go:557-569` 那段（逐字 `审批将在 %d 秒后自动拒绝，请尽快确认`）在 `#6` 那一发里会**顺路端到端响一次**。r1 `evidence.md §7.2` 说"这半格没人读得到"，甲是四支里唯一一枚不新造任何门就能把它读到的 |
| 甲的风险 | §7.1 包时长余量未量；§7.10 那发究竟哪枚时钟先赢未实测 |

### 乙 · 给测试留一条"直接构造 `approval.NewQueue`／`approval.New`"的窄缝

- **它需要新造门吗？不需要。** 两枚现成的门都在 HEAD 里、都有具名裁定（逐枚见 §4）：`approval.Options.ApprovalTimeout` 是**公开字段**，`cmd/wisp/resident_grant_writer_265_windows_test.go:187` 今天就在 `cmd/wisp` 的测试里直接 `approval.New(approval.Options{…})`；而 `runSpec.gate/ui/cards`（`cmd/wisp/run.go:148-166`）是票 246 AC#7 裁过的注入点。
- **但它的真实代价在"必须三件套一起递"这一条上**：`run.go:600-605` 逐字 `if s.gate != nil { if s.cards == nil \|\| s.ui == nil { … 已拒绝装配 … return rt, 2 } }`，且 `:607` 逐字 `rt.ui = nil // no console surface in a process whose cards go to another host`。⇒ 一旦走注入，**控制台面就没了**，而 §2.1 里 #2/#3/#4/#5/#9/#10/#11/#12 的断言正读控制台与审计文本（`h.out`/`h.err`、`rt.ui.shown()`、`waitForSentence187`）。要保住这些句子，测试得自带一枚能冒充 `consoleApprovalUI` 的 UI —— **那就从"复用既有门"滑向"用 mock 顶掉真面"**，正是名册原文禁的那一类（§4.1）。
- 乙对**唯一真正需要它的用例**是干净的：`approval_reply_201_test.go:439` 那一发不需要控制台文本以外的东西吗？它需要 `:484` 那句审计（`ANSWER-EXPIRED …`），而审计是 `run.go:617` 的 `Logf: rt.auditf` 写进 `h.err` 的 ⇒ 注入 gate 就得连 `Logf` 一起重接。**判语：乙能替甲，但它是"为 1 枚用例重开整条装配面"，成本集中且不可复用；甲的成本摊在 12 枚位置上的其中 2 枚。**
- 乙还带一笔**裁定成本**：它到底算不算"新造一枚注入缝"＝人工批准面（§7.7），本腿只交判据不交结论。

### 丙 · 改用 fake clock

- **接缝在不在？在，而且早就在**：`internal/agent/approval/approval.go:20` 逐字 `Clock is the injected monotonic time source (D42#9 / D22 ban 5). Every`，接口 `:26-32`（`Now()` + `After(d time.Duration) <-chan time.Time`），`SystemClock` 在 `:34-41`，`Options.Clock` 在 `gate.go:21-22`，`New` 的兜底在 `gate.go:127-130`（逐字 `clock := o.Clock; if clock == nil { clock = SystemClock{} }`）。三枚读时钟的点全走它：`gate.go:291`（L1 窗口）、`gate.go:526`（C18 deadline）、`gate.go:529`（warning lead）。
- **cmd/wisp 吃不吃得到？吃不到。** 两枚生产构造点都没传 `Clock`：`run.go:612-619`、`resident_approval_windows.go:368-375`（09:53:56 现读两枚 Options 字段的完整清单＝UI/Channels/Window/ApprovalTimeout/Logf/Grants）；`grep -rn "Clock\b" cmd/wisp/*.go` ＝ **0 命中**（09:49:39）。⇒ **丙必须以乙为前提**：想让 `#6` 用假时钟，先要在测试里自己 `approval.New(approval.Options{ApprovalTimeout: 31 * time.Second, Clock: fake})`，也就是先付乙的全部成本。
- **现成的假时钟导不过来**：`fakeClock` 在 `internal/agent/approval/fakes_test.go`，而该文件首行是 `package approval_test`（`:1`）⇒ **外部测试包，不可 import**；`newFakeClock()`（`:39`）、`Advance()`（`:62`）、`newGate(t, ui, opts...)`（`:250-260`，逐字 `o := approval.Options{UI: ui, Clock: clk, Channels: ch}`）都只在 approval 那包内可用。cmd/wisp 要用就得**再写一枚**（同形的 `neverClock` 先例在 `ticket220_l1_window_read_test.go:84-88`）。
- **丙对整条 run 的副作用是本腿最担心的一格**：这些用例跑的是**完整的 `runTextTask`**（mockllm 走真 HTTP、SQLite 真写、面板泵与热加载 tick 真计时）。假时钟只喂 gate，一 `Advance(31s)` 就让审批死线与进程里其余计时器**脱钩**；`#2`（panel_pump）正是"跑自己的 publish 也 booked 一条记录"的用例（`panel_pump_test.go:207-213`），假时钟下先后次序不再是现场次序。**⇒ 丙适合 `internal/agent/approval` 那一层，不适合 `cmd/wisp` 这一层的整进程用例。**

### 丁 · 本腿在盘上找到的第四形

派单给的线索成立，而且**比线索本身更宽**，本腿现量到两枚合法小时钟：

1. **L1 窗口那枚（线索说的那枚）——合法、带外无约束、但替不了这次的主题。** `gate.go:143-151` 的钳位（`win <= 0 → DefaultL1Window`／`< MinL1Window → 2s`／`> MaxL1Window → 3s`，常量在 `queue.go:116/120/122`）；而 `validateRisk`（`validate.go:141-151`）**只校 `permission_mode` 与 `confirm_timeout_sec` 两样，`l1_window_sec` 连 config 级的带都没有**。⇒ "仓里已经存在一条合法的小窗口"这句**为真**。
   **但它救不了 §2.1 的任何一枚**：`#2` 被文件自己否掉（`panel_pump_test.go:160-162` 逐字 `the verdict is L2, so the item goes into the approval QUEUE (an L1 window does not), which is what snapshot.pending is made of`）；`#6` 被极性否掉（L1 到点**执行**，`#7` 的 `:545` 逐字 `the L1 timeout polarity changed (this is the D4 contract, not this leg)` 就是钉着这条极性的句子，而 `#6` 的主题正是"L2 到点必须拒绝"）；`#1` 的观察对象是 L2 卡放行档，也不在 L1 上。**唯一能被 L1 救的是 `#7`/`#8` 两枚，而它们今天本来就跑在 L1 上（那 20s 种子对它们是死重）。**
2. **`runSpec.gate` 那枚（本腿新找到的，可以完全不动 band）**：见 §4.3，票 246 AC#7 已裁、`resident_*_265` 已在用。**它实际上是乙的"合法化版本"**：不是"新造一条窄缝"，而是**复用一枚仓里已存在的窄缝**＋公开字段 `Options.ApprovalTimeout`。
3. **本腿找到的第三形（写在这里因为它比甲便宜，且不属任何一支）：把带外种子里"其实没人观察它"的那几枚直接删掉，让它落 schema 默认 300。** 现量证据：`#1` 走 300ms ctx、`#7`/`#8` 走 L1 窗口、`#3/#4/#5/#9/#10/#11/#12` 答完卡就走 ⇒ **这 9 枚位置上的种子值今天没有一枚是观察对象**，写 20 或写 300 对断言等价。真正需要"短"的只有 `#2` 与 `#6` 两枚。**⇒ 最小改动面＝2 枚抬进带内（+ ctx 连带 1 枚）+ 7 枚删行/抬平 + 1 枚兜底句**，而不是 12 枚逐一重排。

## §4 已有的窄缝先例（乙那格：D22 注入缝名册现量）

### §4.1 名册原文（逐字，⛔ 本腿不扩它）

`.scratch/wisp/issues/README.md:206-207`：

> `Tests inject at seams only: C8 AudioSource (wav), C5 LlmProvider (golden SSE), C17 PanelBridge,`
> `CLI `wisp run`. No mock-instead-of-real to fake completion; never weaken SLO thresholds.`

同形中文摘要在 `AGENTS.md` §1.3（派单已抄给本腿）。⇒ **名册里没有"approval 队列/门"这一枚。**
另有一把仪器尺：`grep -rn -i "inject\|mock\|seam" tools/d22scan/*.go`（09:52:15）里没有**任何**一条 ban 针对"测试构造 `approval.New`/`NewQueue`"；d22scan 侧的 `inject` 只指它自己的测试内缝（`main.go:1396` `fixtureScope is injected by the tests ONLY (ticket 88 AC#3)`、`:1410` `fixtureVerdict exists so a test that injects scopes says so in the function`）。⇒ **乙不会被 CI 扫红，但也不被 CI 认可**；认不认是裁定面，不是仪器面（归口 §7.7）。

### §4.2 现成的同类先例（逐枚带 file:line，判语在每枚后面）

| 先例 file:line | 逐字形状 | 它是不是"不经 config 的真对象构造"？ |
|---|---|---|
| `cmd/wisp/resident_grant_writer_265_windows_test.go:187` | `f.g = approval.New(approval.Options{ UI: f.ui, Channels: approval.NewChannels(), Logf: f.log.write, Grants: f.holder, })` | **是，而且就在 `cmd/wisp` 包里、就在真 `approval.Gate` 上**。`#171-175` 的注释逐字把这件事说在前头：`Options carries four of the fields the resident literal carries (UI / Channels / Logf / Grants); Window and ApprovalTimeout are left at zero, which approval.New clamps to its own defaults` ⇒ **"测试自己造门、并把超时留给默认"已经是这棵树的既成事实，且写明了它为什么不算假** |
| `internal/config/validate_267_test.go:130` | `q := approval.NewQueue(time.Duration(c.Risk.ConfirmTimeoutSec)*time.Second, 0, 0, nil)` | **是本票自己刚造的先例**：不经装配根、直接构造**真 `approval.Queue`**（⛔ 不是 mock），r1 `evidence.md §3` 逐字为它写了两条理由（"喂法逐字照装配根 `run.go:616` 与 `resident_approval_windows.go:460`"）。⇒ 乙要的那条形（`NewQueue`/`New` 直构）**已经在 HEAD 里落了一枚，且带同源守卫的身份** |
| `internal/agent/approval/fakes_test.go:250-260` | `func newGate(t *testing.T, ui approval.UI, opts ...approval.Options) (*approval.Gate, *fakeClock, *approval.ChannelRegistry) { … o := approval.Options{UI: ui, Clock: clk, Channels: ch} … o.Clock = opts[0].Clock }` | 门＋假时钟的现成配方，但**首行 `package approval_test`（`:1`）⇒ 不可 import**，只是形状先例不是可复用件（§5.1） |
| `cmd/wisp/approval_reply_stdin_other.go:13` | `// answer cards pass their own reader through runSpec.reply (the `wisp run` seam` | `runSpec.reply` **自己就被注释指认为 AGENTS §1.3 名册里那枚 `CLI wisp run` 缝**（另见 `run.go:132-140` 的逐字辩护：`the CLI tests fill it with a scripted reader, which is the injection seam AGENTS.md §1.3 names for `wisp run` - it is not a mock standing in for a missing subsystem, because the subsystem here IS a human typing at this terminal`） |
| `cmd/wisp/run.go:115-119` | `// onRuntime hands the assembled stack to the caller before the task runs. Production leaves it nil; the composition tests use it …` | 又一枚"生产留 nil、测试用真装配根内部对象"的**已裁先例**：§2 的 12 枚调用点全都靠它拿 `rt.gate`/`rt.bridge`/`rt.ui` |
| `cmd/wisp/run.go:120-126` | `modeConfirm is the L2 strong confirmation … The end-to-end tests fill it to stand for "the operator clicked allow"` | **测试用非空 `modeConfirm` 顶掉真 L2 卡**已被成文允许（`run_mode101_test.go:219-225` 就是它的用法）⇒ 与乙同族：测试改的是"谁来答"，不是"造不造真对象" |
| `cmd/wisp/secret_argv_windows_test.go:161`＋`:174`＋`:234` | `func buildWispForTest(t *testing.T) string {` … `cmd := exec.Command(goExe, "build", "-o", exe, "./cmd/wisp")` … `cmd := exec.Command(exe, args...)` | **进程级缝**：真二进制、真 config、真加载链 ⇒ **band 在这一族用例里绕不过去，也不该绕**（乙与它无关，登记为"窄缝不该被推广到端到端腿"的界标）。同形另有 `resident_approval_live_246_windows_test.go:417`（`exec.Command(goExe, "build", "-o", exe, "./cmd/wisp/testdata/esclistener")`）与 `:446` |
| `cmd/wisp/instructions_200r2_test.go:92-104` | `// disableProjectInstructions writes AC#8's switch into the fixture's config. It appends a section the fixture does not carry, the same way nonDefaultConfig145 does for [risk]` | **"往夹具 config 追加一节"这形状本身是被点名继承的先例**（`panel_pump_test.go:136` 是它的祖师）。⇒ 甲改的正是这枚先例里的**数值**，不是这枚先例的形状 |

### §4.3 那枚最值钱的先例：`runSpec.gate / ui / cards` 已经是"合法窄缝"

`cmd/wisp/run.go:148-166` 逐字：`gate, ui and cards are ticket 246 AC#7's injection points, and they exist for exactly one reason: a process may hold ONE approval gate.`，用法现读在 `cmd/wisp/resident_task_source_246_windows_test.go:153`：`{"gate without ledger", runSpec{stdout: f.out, stderr: f.err, dataDir: f.dir, gate: ra.gate, ui: ra.ui}}`。
⇒ **判语（本腿只给判据）：乙不是"新造一枚注入缝"，而是"把一枚已裁过的注入缝从常驻腿延伸到跑任务腿"。** 它需要新裁的不是"缝存不存在"，而是两件事：① 允许在 `cmd/wisp` 测试里传**非默认** `ApprovalTimeout`（今天唯一那枚直构先例 `:187` 刻意**没**传，见其 `:171-175` 的自限句）；② 接受注入路径会关掉控制台面（`run.go:607`）。这两件都够得上"具名 A##/人工批准"，不是写腿自选。

## §5 clock 与 L1 窗口现量（丙 / 丁）

取数时刻 **`09:54:58–09:55:34+0800`**（尺见每小节末），HEAD＝`9ca14bdf` 之后、§3/§4 那次编辑之前。

### §5.1 clock 是不是已可注入？——**是，早就可注入；但生产与 `cmd/wisp` 都没用它**

- 接口与实现：`internal/agent/approval/approval.go:20` 逐字 `Clock is the injected monotonic time source (D42#9 / D22 ban 5). Every`、`:26-32` `type Clock interface { Now() time.Time; After(d time.Duration) <-chan time.Time }`、`:34-41` `SystemClock`。
- 注入点：`gate.go:21-22` `// Clock is the monotonic time source (D42#9). nil uses SystemClock.` + `Clock Clock`；兜底 `gate.go:127-130`。
- **全部三枚读时钟的点都经它**（这是"可注入"的硬证）：`gate.go:291` `deadline := g.clock.After(g.window)`（L1）、`gate.go:526` `deadline := g.clock.After(g.q.Timeout())`（C18）、`gate.go:529` `warn = g.clock.After(lead)`（**本票那枚提示**）。
- ⚠ **`Queue` 自己没有时钟**：`NewQueue(timeout, warnBefore time.Duration, maxPending int, logf func(string, ...any)) *Queue`（`queue.go:84`）签名里没有 `Clock`，唯一调用者是 `gate.go:158` `q: NewQueue(o.ApprovalTimeout, o.WarningLead, o.MaxPending, logf)` ⇒ **想让假时钟驱动 C18 死线，必须整个 `Gate` 一起换**，不能只换队列（这条决定了丙不能"小到只动测试里那一行"）。
- 现成的假时钟有两枚，都**在 approval 包内**：`fakes_test.go:21-77`（`type fakeClock struct`、`newFakeClock()` `:39`、`After` `:47`、`Advance(d)` `:62`；配方在 `newGate` `:250-260`）与 `ticket220_l1_window_read_test.go:84-88` 的 `neverClock`。前者所在文件首行是 `package approval_test`（`fakes_test.go:1`）⇒ **不可被 `cmd/wisp` import**。
- 谁在传 `Clock`（尺＝`grep -rn "Clock:" internal --include='*.go'`，09:55:34，共 10 行）：**全部 10 行都在 `internal/agent/approval/*_test.go`**（`fakes_test.go:254`、`ticket220…:203/:383`、`ticket224_reply_grant_test.go:120`、`ticket84_no_owner_test.go:59/:154/:182`、`ticket87_veto_l2_test.go:49`、`ticket97_alias_direction_test.go:32`）。⛔ **零枚产码传**。
- `cmd/wisp` 侧（尺＝`grep -rn "Clock\b" cmd/wisp/*.go`，09:49:39）＝**0 命中**；两枚生产构造点 `run.go:612-619`、`resident_approval_windows.go:368-375` 的 Options 字段全集现读为 `UI/Channels/Window/ApprovalTimeout/Logf/Grants`（**没有 `Clock`、也没有 `WarningLead`** ⇒ 后者落 `queue.go:88-89` 的兜底 `DefaultApprovalWarning`，正是票 267 的成因，本腿在此复钉一次）。
- **判语**：丙的接缝**已存在且形状正确**（它就是 D22 ban 5 指定的做法），但在 `cmd/wisp` 这一层它今天**吃不到**，要吃到必须先走乙；而"测试自带一枚 fakeClock"这一支会撞上 §3 丙末段那条脱钩风险。⇒ **丙在 approval 层值得做、在 cmd/wisp 层不该做**。

### §5.2 L1 窗口那枚合法小时钟现量（丁）

- 钳位：`gate.go:143-151` 逐字 `win := o.Window; switch { case win <= 0: win = DefaultL1Window; case win < MinL1Window: win = MinL1Window; case win > MaxL1Window: win = MaxL1Window }`；常量 `queue.go:116` `DefaultL1Window = 3 * time.Second`、`:120` `MinL1Window = 2 * time.Second`、`:122` `MaxL1Window = 3 * time.Second`。
- **config 层对这枚键零约束**（尺＝读 `validate.go:141-151` 全文，只两枚判语：`risk.ParseMode(c.Risk.PermissionMode)` 与本票新增的 band）⇒ `l1_window_sec` 可以写 1、99、任何数，钳位发生在接缝里。票 256 已经把这枚事实钉成断言：`resident_approval_risk_256_windows_test.go:228-232` 逐字 `"window way too large - reverse control, the clamp must survive"` + `seed: "l1_window_sec = 99\n"` + `wantWin: approval.MaxL1Window`，以及 `:235-239` 的 `l1_window_sec = 1` → `MinL1Window`；`:263-265` 还有一枚"带不许被喂开"的反向钉 `if w := ra.gate.Window(); w < approval.MinL1Window || w > approval.MaxL1Window {`。
- **`cmd/wisp` 里今天确实有"用合法小窗口看卡开合、且不靠 `confirm_timeout_sec`"的活样本**（这是丁最硬的一枚证据）：`cmd/wisp/run_test.go:373` `cards = rt.windowCount()` ＋ `:382` `t.Fatalf("an unvetoed L1 window means EXECUTE, got: %s", text)`，而它用的夹具 `newRunFixture`（`run_test.go:76-115`）**整个 `[risk]` 一节都没写** ⇒ 这一族用例今天跑在 schema 默认 **300s**（带内、band 落地后仍绿），观察对象是 2s 的 L1 窗口开合。⇒ **"看一张卡开合"在跑任务那条腿上本来就有不带 `confirm_timeout_sec` 的走法。**
- **但它替不了 §2.1 的两枚真短种子，理由各不同且都在盘上**：
  - `#2`（`panel_pump_test.go:136`）：要的是**队列里的那一枚 pending L2**，`snapshot.pending` 由它构成（同文件 `:160-162` 逐字否决 L1）。
  - `#6`（`approval_reply_201_test.go:439`）：要的是**"到点＝拒绝"这枚极性**；L1 到点＝执行（`:544-546` 就是钉住这条极性的句子）。这一发的主题与 L1 相反，换过去＝观察对象归零。
- 除 L1 之外，本腿另外量到两枚**已有的、能替代"短超时"的旋钮**：
  1. `runSpec.reply`（`run.go:132-140`，名册点名的 `wisp run` 缝）——把"等到点"换成"答它"，#3/#4/#5/#9/#10/#11/#12 今天就是这么做的（例：`approval_reply_201_test.go:303` 写 `"no …"`、`:229` 写 `"yes …"`）。**对 `#2` 它只能让卡"被答而关"，而 `:181` 那句断言逐字要求 `an unanswered L2 card must not execute` ⇒ 用它会改掉断言语义，不是无损。**
  2. `runSpec.modeConfirm`（`run.go:120-126`）——`run_mode101_test.go:219-225` 用它顶掉整张真卡。⇒ `#1` 那枚 `confirm_timeout_sec = 1` 在这一族里**没有任何一枚断言经它**（`#1` 的 L2 走 300ms ctx 或走 stub），与 §3 丁的第 3 形一致：删行即可。

### §5.3 本小节的三把尺（逐把带时刻）

| 时刻 | 尺 | 读数 |
|---|---|---|
| 09:55:34 | `grep -rn "Clock:" internal --include='*.go'` | 10 行，**全部**在 `internal/agent/approval/*_test.go`；产码零枚 |
| 09:55:34 | `grep -rn "NewQueue(" internal cmd --include='*.go'` | 9 行：产码 2（`queue.go:84` 定义、`gate.go:158` 唯一调用）＋ 测试 7（`pending_read_test.go:112/:191`、`ticket146_liveapprovals_backing_test.go:209/:265`、`ticket242_binding_test.go:59/:88` 六枚 `NewQueue(0, 0, 0, nil)`，加本票的 `validate_267_test.go:130`）⇒ **"测试直接构造真 Queue"已是 7 枚先例**，乙在"这算不算 mock"这一问上并不孤立 |
| 09:49:39 | `grep -rn "Clock\b" cmd/wisp/*.go` | **0**（§5.1 的"吃不到"由此得） |

## §6 尺读数与未跑清单

### §6.1 每把尺的取数时刻（本腿自取，⛔ 零枚抄别家腿或票面）

| 时刻 (+0800) | 尺 | 读数 |
|---|---|---|
| 09:38:38 | `date "+%H:%M:%S%z"` / `git rev-parse --short HEAD` | 起手 HEAD **`873c3063`** |
| 09:39:11 | `git commit -- pathspec`（未 add 的裸文件） | `error: pathspec … did not match any file(s) known to git`（未跟踪文件必须先 `git add`） |
| 09:39:16 | `git add <本件> && git commit -F <临时txt> -- <本件>` | 骨架 commit **`8aa9ff6f`**，24 行 |
| 09:44:46 | 六根合并尺，带 `--include='*.go\|*.toml\|*.md\|*.ps1'` 过滤 | **213** 行 |
| 09:44:46 / 09:44:56 | `git rev-parse --short HEAD` | **`2c83c076`**（期间 `0c2445d1`、`492571d8` 落到同一棵树 ⇒ 并发确认） |
| 09:44:56 | 逐根 `grep -rIni`（无 include 过滤）：`cmd / internal / tools / docs / scripts / .scratch` | **21 / 42 / 0 / 27 / 0 / 749** |
| 09:45:05 | `… internal cmd \| cut -d: -f1 \| sort \| uniq -c` | cmd 侧逐文件：`resident_approval_risk_256_windows_test.go` 11、`resident_approval_windows.go` 4、`panel_pump_test.go` 2、`run.go` 1、`run_mode101_test.go` 1、`resident_windows.go` 1、`approval_reply_201_test.go` 1 |
| 09:45:05 | `.scratch` 分桶 | probes **717**（本票 `probes/267` 83，其中 `267/r1/evidence.md` 42）· issues **21** · ci-logs **4** · `ledger-restore.tmp` 4 · `probes-255-r1c` 2 |
| 09:46:20 | `docs` 分文件 | 14 枚文件共 27：`pending-and-issues.md` 8、`evidence/s1` 15（分 11 枚文件）、`PLAN.md` 1、`SPEC-03` 1、`missing-features-…v4.md` 1 |
| 09:47:59 | `grep -rn "newReplyHost(t, " cmd/wisp` | 直接喂值调用点 **10 枚**（20s ×9、2s ×1）＋ 2 枚 helper 形参转手 |
| 09:47:59 | `grep -rn "newWidenRun(t, \|newReloadRun223(t, " cmd/wisp` | 间接调用点 **15 枚**（40s ×10、90s ×4、60s ×1）＝全部带内 |
| 09:49:18 | `sed` 逐段读 `approval_seam_201_test.go`／`ticket224_assembly_test.go` | 两枚 20s 调用点的断言原文已取得（§2 #9–#12） |
| 09:49:39 | `grep -rn "Clock\b" cmd/wisp/*.go` | **0 命中** ⇒ 丙在 `cmd/wisp` 今天没有任何入口 |
| 09:49:45 / 09:49:50 | `grep -rn "go test" scripts`、`grep -n "timeout\|go test\|cmd/wisp" scripts/wisp-cli-tests.sh` | 跑 `cmd/wisp` 的是 `scripts/wisp-cli-tests.sh`（`:111` scope=`./cmd/wisp/`）；**两把尺都找不到 `-timeout` 字样** ⇒ 包级超时是 `go test` 的默认值，本腿不去猜它是多少（见 §7.1） |
| 09:50:27 | `grep -n "t\.Errorf\|t\.Fatal" cmd/wisp/approval_seam_201_test.go` 等 | 断言行号全部就地核字符确认（§2 引用逐枚） |
| 09:51:59 | 本件 commit §2 后 | **`c9600334`**，163 行 |
| 09:52:15 | `grep -n -i "inject\|seam\|mock" .scratch/wisp/issues/README.md`、`grep -rn -i … tools/d22scan/*.go`、`ls internal/*/testdata` | 名册原文＝`README.md:206-207`；d22scan 侧**没有**任何"测试构造 gate/queue"的禁令；`testdata/` 现存四棵：`internal/agent/testdata`(golden)、`internal/llm/testdata`(golden)、`internal/models/testdata`(tiny-model.tar.bz2/README)、`scripts/testdata` |

### §6.2 收口读数（09:58:25–09:58:42+0800，本腿自取）

| 尺 | 读数 |
|---|---|
| `git rev-parse --short HEAD` | **`8d44e7ff`**（本腿最后一枚内容 commit；起手 `873c3063` → 骨架 `8aa9ff6f` → `89976ed1` → `c9600334` → `9ca14bdf` → `7b87b991` → `8d44e7ff`，共 **6 枚 commit**，全部只含本件） |
| 逐枚验面 `git show --name-only`（对上面 6 枚各取一次） | 每一枚的文件清单都只有 **`.scratch/wisp/probes/267/a2/census.md`** 一行 ⇒ "pathspec 只写自己那一枚文件"这条**逐枚成立**，且零 `git add -A`/`git add .` |
| `wc -l` / `wc -c` 本件 | **09:58:25 读＝345 行 / 73,563 字节**；⚠ 这是**写该行时刻**的快照，本件之后又长了（10:00:32 复尺＝**361 行 / 77,555 字节**，当时 HEAD 已被并发腿推到 `82a9f7a3`）；终值以最后一枚 commit 后的复量为准 |

| 占位符尺 `grep -n "待填\|未判\|TBD\|TODO\|占位"` 本件 | **命中 1 行＝本行自己**（尺的字面混进被扫文件里，这是本腿写的一枚假阳性，09:59:50 现读）；⛔ **除此之外 0 命中** ⇒ 硬预算闸门那句"全文不许留占位符"成立（§7 里的"判不动/量不到"是小节标题＋**具名归口**，不是占位符）。这一处已写进 §8 第 14 条 |
| scoped 脏面尺 `git status --porcelain -- cmd internal tools scripts docs` | 起手（09:38:38）**0 行** → 收口（09:58:33）**1 行**＝` M docs/reports/pending-and-issues.md`。⚠ **那不是本腿写的**：派单规矩 2 明令不碰台账，本腿零次打开过它写面（只在 §1.3/§7.5 引用过它的行号，那是 grep 读）。这一行＝编排者自己在共享树上的写面，本腿如实登记、不动、不 `checkout` |
| `git diff --name-only 873c3063..HEAD` 全树 | 除本件外还有 `.scratch/wisp/probes/267/r1/evidence.md`、`r1/final-verify.txt`、`probes/267/a3/census.md`、`issues/167-*.md`、`issues/267-*.md`、`docs/reports/pending-and-issues.md` ⇒ **同机此刻至少三枚别的面在飞**（r1 写腿、`a3` 普查腿、编排者） |
| ★并发腿具名 | **`267-a3`** 已于 `4e877958` 落骨架（起手 `09:53:08+0800`／锚 `c9600334`＝本腿 §2 那枚 commit），它的射程是**用户可见文案名册／常量与默认值／日志审计／文档含冻结件标注** ⇒ **与本件（调用点与名册尺）不重叠**。本腿没有读它的正文，避免把别家腿未核的读数当自己的尺（派单"不许复用别家腿读数"） |
| 零 go 尺 | 本腿全程只跑过 `date`/`git`/`ls`/`grep`/`sed`/`wc`/`awk`/`cut`/`sort`/`uniq`，⛔ 一条 `go`／`wisp`／`sh scripts/*` 都没有（§6.3 那 6 项都是据此未跑而量不到的） |

### §6.3 未跑清单（派单规矩 1 的账，逐枚具名）


⛔ 本腿**一条 `go` 命令都没跑**，因此下面每一项都是"本可以实测、按令未实测"：

1. `go test ./cmd/wisp/` ＝ 明令禁跑 ⇒ **§2 那 12 枚带外位置的红，形态是机制判语不是读数**：本腿说不出"先炸的是 `run_mode101_test.go:551` 的 `t.Fatalf("exit %d\n%s", code, log)` 还是 `:245` 的同形句"，也说不出退码到底是 2 还是 1。
2. `go test ./internal/config/`、`./internal/agent/approval/` 未跑 ⇒ §1.2 里那些"绿"的判语（票 267 自己的钉子、unwired 夹具 300/60、boundary 60）是**读断言 + 读带**得来的。
3. `go build ./...`、`go vet`、`gofumpt` 未跑 ⇒ 本件不对任何"编译面/格式面"下结论。
4. `wisp slo`、`sh scripts/*`、`tools/d22scan` 未跑 ⇒ 未验证 band 会不会把 D22 静态门或 SLO 门顶红（本腿只确认 `tools/` 与 `scripts/` 里**零枚**这枚键的出现，见 §1 总量表）。
5. `cmd/wisp` 包级总时长未量 ⇒ 甲的 +58s 落不落得下，见 §7.1。
6. 面板/球侧那张 C18 卡的实际文案到达与否未量 ⇒ 一因未跑，二因 `frontend/**` 是本腿禁读面（§7.8）。

## §7 判不动 / 量不到（具名归口；⛔ 本节不许出现"应该没问题"）

### §7.1 `cmd/wisp` 包的总时长余量 —— **判不动**
甲在本腿的账是 +29s（`approval_reply_201_test.go:439`）+29s（`panel_pump_test.go:136`）+ 失败路径每发 30s（`run_mode101_test.go` 里 `t101call` 用的是无期限 `context.Background()`，`:195`）。这棵树今天有没有一枚尺量着 `cmd/wisp` 的包时长，本腿没找到（`scripts/wisp-cli-tests.sh` 里搜不到 `-timeout`），而 `scripts/portable-tests.sh:5` 只列了"go test <16 packages>"。**⇒ 归口＝编排者**：只有当机那枚跑 go 的腿能给出"抬完种子后这一包还在不在默认包超时内"。本腿拒绝用"应该还早"填空。

### §7.2 band 之后那 12 枚调用点的**确切红形态** —— **量不到**
本腿能给的是链条：`os.WriteFile` 落盘 → `config.LoadFile`（`loader.go:41`）→ `readConfigFile` → `validate`（`loader.go:137`）→ `validateRisk` 的 `observe.New(observe.ClassConfig, "config.toml: risk.confirm_timeout_sec %d out of range [31, 3600]")`（`validate.go:145-148`）→ 装配根拿不到 cfg。给不了的是**这一包实际打印哪一句、第几枚用例先停**（未跑，§6.3 第 1 项）。**⇒ 归口＝编排者排的下一枚 `cmd/wisp` 写腿**（它本来就要动这批文件）。

### §7.3 常驻腿"越界"与"缺文件"共用同一枚 `unreadable` 判语 —— **判不动，且本腿认为这是本把尺量出的最重一枚**
`resident_approval_windows.go:453-457` 逐字 `if err != nil || c == nil { slog.Warn(…); return 0, 0, riskProvenanceUnreadable }`：band 让**一份存在但越界**的 config 与**一份不存在**的 config 得到同一个 provenance 词与同一个 300s 回落。票 256 的 `:141-144` 那枚正控（`(300s is the contract default; a fourth number here means the fallback grew a value of its own)`）**分不开这两种**，而裁形 ⓐ 的原话是"越界＝加载时拒"。要不要给常驻腿加第三种 provenance 词、要不要让它拒启动，都不在只读普查的权里。**⇒ 归口＝编排者**（且它可能与 `Q-77` 同一枚决定面，见 §7.5）。

### §7.4 C18 那张卡的"即将超时"提示端到到不到面板 —— **量不到**
`gate.go:557-569` 那段（逐字 `审批将在 %d 秒后自动拒绝，请尽快确认`）需要真时钟跑到 `lead`；本腿⛔禁跑 go、⛔禁改 `internal/agent/approval/**` 与 `cmd/wisp/**`，而且 `frontend/**` 是禁读面。**⇒ 归口＝下一枚 `cmd/wisp` 写腿顺手一发**（本腿在 §2 #6 给了它一枚免费正控的形状：种子 31 ⇒ `lead = 1s`）。这与 r1 `evidence.md §7.2` 自己交回的那半格是同一枚洞，不重复记账。

### §7.5 `Q-77`（C18 硬编码 300s vs 可配 `confirm_timeout_sec` 谁优先）—— **待人拍板，本腿不动**
台账逐字在 `docs/reports/pending-and-issues.md:11761`，`docs/specs`/票面都写着⛔任何腿不许自行裁。本件的 §2.3 与 §7.3 都在它射程里（回落值 300 恰是 C18 那个数），但**本腿一律按"不答也能成立的机制描述"写，没有替它裁任何一支**。

### §7.6 值域要不要进 `SPEC-03`／`PLAN.md` 的 D36 文本 —— **判不了（属人工批准面）**
`docs/specs/SPEC-03-config-secrets-envs.md:34` 现在只写 `confirm_timeout_sec(int)=300`、`docs/PLAN.md:2738` 只写 `(300)`，两处**都没有值域一格**。band 落地＝给 SPEC-03 加一条它今天没有的规矩。改规格文字＝改契约面（AGENTS §0 第 2 句、`SPEC-12 §4.1`）。**⇒ 归口＝编排者**（本腿只登记"文字与规矩已经不同步"这一事实，不动一字）。

### §7.7 乙那格算不算"新造一枚注入缝"—— **本腿给判据，不给裁定**
名册原文逐字＝`.scratch/wisp/issues/README.md:206-207`：`Tests inject at seams only: C8 AudioSource (wav), C5 LlmProvider (golden SSE), C17 PanelBridge, CLI `wisp run`. No mock-instead-of-real to fake completion`。**扩这枚名册（或裁定"测试直接构造真 `approval.New` 不属于新缝"）＝人工批准**，只读腿无权裁。本腿在 §4 只交两样东西：现成的同类先例逐枚 file:line，以及"先例里那些门是谁裁的"。

### §7.8 面板/球侧的第二枚消费点 —— **量不到（禁读面）**
派单规矩 3 禁读 `frontend/**` 与 `design/**`。`confirm_timeout_sec` 若在渲染层还有第二枚读者（例如把 `Prompt.Deadline`（`gate.go:610`）画成倒计时），本腿看不见也不转述。**⇒ 归口＝那一编队**；本件全部结论只覆盖 Go 侧与文档侧。

### §7.9 键名的拼接形状 —— **量不到的形状**
本腿的尺是 `-e confirm_timeout_sec -e ConfirmTimeoutSec`（大小写不敏感）。若某处用**变量拼出键名**（`"confirm_" + "timeout_sec"`、正则、或 `toml.Marshal` 之后改写），grep 形状看不见。`internal/config/schema.go:460` 的 toml tag 是本腿能确认的唯一具名形状；**"零枚隐藏拼接者"这句本腿不敢说**，只说"未发现"。

### §7.10 `panel_pump` 那发在 31s 下**究竟改看哪枚时钟** —— **判不动（两相未实测）**
本腿读到 `panel_pump_test.go:112` 的 30s ctx 与 `gate.go:571`（`<-deadline`）/:576（`<-ctx.Done()`）在同一枚 `select` 里，**读到的是"两枚 case 并存"，不是"哪枚先赢"**。`30 < 31` 让机制判语倾向 ctx 先响，但这一句在没跑过之前不许当结论用（§6.3 第 1 项）。**⇒ 归口＝编排者**：甲若被选，这一发要**连 ctx 一起抬**，而抬完必须有一枚实测读数才算交付（本腿给不了）。

## §8 我写错的读数（自我对抗）

1. **名册数错（已就地更正，原话留着）**：§1.5 首发把带外喂值写成"12 枚调用点、3 枚文件、14 个位置"。现量（09:47:59，`grep -rn "newReplyHost(t, " cmd/wisp`）＝**喂值调用点 10 枚**（9 枚 20s、1 枚 2s），字面种子 2 枚，合计 **12 枚带外位置 / 5 枚文件**。"14 个位置"是把 helper 形参那 2 行也算了进去——重复计数。
2. **`internal` 那行的产码/测试拆分写错（已更正）**：首发"产码 11 + 测试 30 + 注释 1"，逐文件尺（09:45:05）实为 **产码 10（`validate.go` 7 + `schema.go` 2 + `unwired.go` 1）＋ 测试 32**（`internal/agent/approval/ticket84_no_owner_test.go:51` 那枚虽是注释，但它住在 `_test.go` 里，不能单列第三类）。
3. **我违了自己的 commit 规矩一次（事实，不抹）**：§6/§7 那次用了 `git commit -m "<长中文>"`，⛔ 不是派单要求的 `git commit -F <临时txt>`。后果＝正文换行被压成一行，`9ca14bdf` 这条 body 的逐行结构丢了。按 AGENTS §1.4「已提交的历史不改写」本腿**不 amend、不 reset**，只登记；前后各枚 commit（`8aa9ff6f`/`89976ed1`/`c9600334`/`7b87b991`）都走 `-F`。
4. **骨架 commit 第一次尝试失败**：`git commit -F <msg> -- <path>` 对**尚未 add** 的新文件报 `error: pathspec '.scratch/wisp/probes/267/a2/census.md' did not match any file(s) known to git`（09:39:11）⇒ 显式 pathspec ≠ 先 add，两件事都要做。
5. **派单点名的"三枚文件"其实是五枚**：`approval_seam_201_test.go:50/:138` 与 `ticket224_assembly_test.go:81/:234` 也经同一枚 `approval_reply_201_test.go:145` 的 `%d` 喂进 20。r1 `evidence.md §7.1`（`:228`）同样只列 6 处并写"这三枚文件里 8 个调用点"。⇒ **本腿的 10 枚有 grep 全文为据，r1 的数是点例**；这句是别人票面上的数，本腿只登记差异、不改（台账归编排者）。
6. **行号漂移我又踩到两次，都靠字符尺退回**：(a) 起初按 `docs/evidence/s1/128-ac4-r1-acceptance.md:425` 里的 `run_test.go:378` 引那句 "an unvetoed L1 window means EXECUTE"，现读在 **`cmd/wisp/run_test.go:382`**；(b) 起初按 `missing-features-2026-09-29-v4.md:194` 引 `schema.go:450`，现读字段在 **`internal/config/schema.go:460`**。另有**别人票面上的过期行号**本腿不动：`internal/agent/approval/ticket84_no_owner_test.go:51` 注释写 `cmd/wisp/run.go:259`，消费点现读 **`run.go:616`**。
7. **派单给的坐标逐枚复核：一枚不错，两枚要补精度**：`queue.go:108-109`（`:108` 注释、`:109` 才是 `DefaultApprovalWarning = 30 * time.Second`）；`gate.go:528` 逐字 `if lead := g.q.Timeout() - g.q.WarningLead(); lead > 0 {` ✓ 未漂；派单说的 "`approval_reply_201_test.go:145` 的 `int(h.l2Wait.Seconds())`" 实为**键在 `:145`、表达式在 `:147`**。
8. **判语上的一次自我推翻**：§2 #1 起初写"抬到 31s 会让这一族多等 30 秒"，读完 `run_mode101_test.go:544` 的 `context.WithTimeout(context.Background(), 300*time.Millisecond)` 与 `gate.go:576-580` 后改为**happy path 零增量**，只有 `t101call` 那枚无期限 `context.Background()`（`:195`）在失败路径吃 30s。差点据此把甲的代价摊到 12 枚上。
9. **`panel_pump` 起初判"甲无损"**，读 `:112` 的 30s ctx 才改判"变形，必须连带抬"（§3 甲、§7.10）。⚠ 这条仍是**机制判语**：两枚 case 同在一个 `select`（`gate.go:571`/`:576`），本腿未跑 ⇒ 归 §7.10 而不当结论用。
10. **§3 丁起初写错一枚事实**：起笔写"`l1_window_sec` 在 config 层有带"，读 `validate.go:141-151` 全文后更正为 **config 层零约束、钳位在接缝 `gate.go:143-151`**。差点据此给出"丁会被 band 一起拒掉"的错判语。
11. **§2.2 的"带内 15 枚"起初数成 14**：漏了 `config_reload_perm_223_windows_test.go:61` 的 60s（尺 09:47:59：`40×10 + 90×4 + 60×1`）。
12. **`.scratch` 那 749 枚起初打算整根不列**：数完发现 ci-logs 与大批 `probes/*/gate-*.txt` 印着 `confirm_timeout_sec = 1/40/300`，**会被下一个人误读成"仓里还有带外种子"** ⇒ 补了 §1.4 那格"假红"警告，而不是让它留在总数里当噪声。
13. **越界自查**：本腿**未改任何已跟踪文件、未产一码、未翻任何 AC 框、未碰台账**；起手与写作期间 `git status --porcelain -- cmd internal tools scripts docs` 均回 **0 行**（09:38:38 与 09:44:46 各一把）。本件只落 `.scratch/wisp/probes/267/a2/census.md` 与 `.scratch/commit-msg-267a2-*.txt`（临时件按规矩 8 只建不删）。⚠ 唯一"看起来像越界"的是 §1.1★／§7.3 描述了常驻腿行为——**只读，未动 `resident_approval_windows.go` 一字**。
14. **收口时又踩两枚自己的坑（同一把尺两次读数不同，都登记不抹）**：(a) 我在 §6.2 里那条"占位符尺"起初写 **0 命中**，09:59:50 复尺回 **1 命中＝那一行自己**（尺的字面混进了被扫文件）⇒ 已就地改成"命中 1 行＝本行自己，除此之外 0 命中"，没有把错数留着当结论。(b) 插入 §6.2 收口块时我用 `old_string` 换掉了原 `### §6.2 未跑清单` 那枚标题，导致同一份文件里出现 `§6.1 → §6.3 → §6.2` 的乱序（09:59:50 标题尺现读），已把两枚标题对调并同步改掉 **3 处** `§6.2 第 1 项/那 6 项` 的交叉引用（`grep -n "§6\.[0-9]"` 复尺 09:59:50 确认现在是 6.1/6.2/6.3 顺序）。⇒ 教训：**在同一文件里插小节时，`old_string` 要包住"标题 + 其后空行"以外的边界，别把下一节的标题当分隔符用。**
15. **一处判语与自己的证据有张力，本腿不掩盖**：§9 说甲"只有 3 处必须抬"，而 §1.2 表里 `approval_reply_201_test.go:105-109` 那枚 `l2Wait < time.Second → 1s` 兜底**今天没有调用点走到**（§2 #3–#12 全部喂 ≥2s）。所以严格讲甲的改动面是 **3 处必抬 + 1 处必须拆的雷 + 9 处可删行**。我在 §3 甲表里已把那句单列为"必须连带改的一句兜底"，但 §9 的"3 处"口径**略窄于 §3 的 5 类**——两处都留着，让编排者按各自的用途读数，而不是把 §9 当唯一口径。

## §9 交件判语

- **代价最小＝甲，且真实改动面比派单以为的窄**：只有 **3 处**必须抬——`panel_pump_test.go:136`（**连同 `:112` 那枚 30s ctx**）与 `approval_reply_201_test.go:439`；剩下 **9 枚带外位置今天没有一枚是任何断言的观察对象**（§2 #1/#3/#4/#5/#7/#8/#9/#10/#11/#12 逐枚给了为什么），删行落 schema 默认 300 即可。§5.2 的活样本 `cmd/wisp/run_test.go:373/:382` 证明"不写 `[risk]`、跑默认 300s、照样看窗口开合"在这棵树上今天就是绿的。甲墙钟总账≈ **+58s**，并附赠一枚只有甲给得出的免费正控：`lead = 31 − 30 = 1s > 0` ⇒ `gate.go:557-569` 那句 `审批将在 %d 秒后自动拒绝，请尽快确认` 第一次端到端响（正是 r1 §7.2 交回的那半格空洞）。
- **乙不算"新造注入缝"**（§4.3）：`runSpec.gate/ui/cards` 是票 246 AC#7 已裁的注入点；`cmd/wisp` 测试直构真 `approval.Gate` 的先例在 `resident_grant_writer_265_windows_test.go:187`（且它 `:171-175` 主动把 `Window/ApprovalTimeout` 留零，说明"直构但不动超时"是先前那枚裁定的边界）；直构真 `approval.Queue` 的先例 7 枚（§5.3）。**但它一次只能划算用在 `#6`**：注入路径 `run.go:607` 把控制台面关掉（`rt.ui = nil`），而 §2.1 里 10/12 枚的断言正读控制台/审计文本 ⇒ 用它替甲＝拆掉 8 枚用例的证据面。要不要为此具名澄清 D22 名册＝人工批准面（§7.7），本腿只交判据。
- **丙在 `cmd/wisp` 层不该做**（§5.1）：接缝早已存在且是 D22 ban 5 指定的形状，但 `NewQueue` 签名里没时钟（`queue.go:84`）、三枚读时钟点全在 Gate 上（`gate.go:291/526/529`）、`cmd/wisp` 零枚 `Clock` 引用 ⇒ 丙必须先有乙；且假时钟只喂 gate，会让 mockllm 真 HTTP、SQLite、面板泵、热加载 tick 与死线脱钩，`#2` 那发的"跑自己的 publish 也 book 一条记录"（`panel_pump_test.go:207-213`）会因此不再是现场次序。丙属于 `internal/agent/approval` 那一层。
- **丁（L1 那枚合法小窗口）成立但救不了这两枚主题**（§5.2）：`#2` 被 `panel_pump_test.go:160-162` 逐字否决（L1 不进队列，`snapshot.pending` 无物可报）；`#6` 被极性等否决（L1 到点＝执行，`approval_reply_201_test.go:545` 就钉着这条）。**唯一能替短种子的既有旋钮是 `runSpec.reply`（"答它"而不是"等到点"），而它对 `#2` 会改掉 `:181` 那句 `an unanswered L2 card must not execute` 的语义。** #7/#8/#12 三枚本来就跑在 L1 窗口上，它们的 20s 种子是死重——这正是"丁已在仓里、只是没替到点子上"的证据。
- **本腿量不到的那一格（具名）**：**抬完种子后 `cmd/wisp` 这包到底红不红、红成哪一句、包时长余量够不够 +58s**（§7.1/§7.2/§7.10）。根因＝派单规矩 1 禁止本腿跑任何 `go` 命令，而 `scripts/wisp-cli-tests.sh:111`（`scope=./cmd/wisp/`）里没有可抄的 `-timeout` 读数 ⇒ 只有当机那枚跑 go 的腿能给。
- **顺带一枚不在四支里、但编排者必须知道的现量后果**：**常驻腿把"越界"读成"读不到"**（`resident_approval_windows.go:453-457` ⇒ 回落编译常量 300s 并把 provenance 写成 `"defaults (config.toml unreadable)"`）。于是裁形 ⓐ 那句"越界＝加载时拒"在两条腿上给出两种答复，而票 256 的 `:141` 那枚正控分不开这两种（§1.1★、§7.3）。这是本把尺今天量出的最重一枚意外。
