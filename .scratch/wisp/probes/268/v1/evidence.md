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

> **本节由续腿 `268-v2` 填写**（21:4x 接手）。前两枚同职腿（`268-v1` 20:1x 断流、`268-v1b` 21:1x 断流）都死于服务故障，
> §0/§3/§4/§5/§6 已实、只剩 AC#1／AC#2 两格与本节。派单明令 ⛔ **不重跑 422 秒整包**（全量日志已在盘＝`deliver-cmdwisp.log`，§3 已核）。
>
> **续腿起手锚（21:45:40–21:45:55 现量）**：HEAD＝`176d6277`（`git rev-parse HEAD`）；`git status --porcelain` 全仓 **727** 行、scoped `cmd/wisp` **0** 行；
> 两枚字节面复量与 §0.1 逐字同值——`resident_approval_windows.go`＝`69a630bf350daf1032c62b921e9d7044`、
> `resident_approval_risk_268_windows_test.go`＝`e2405a21bb278b558c1305ea95d34a9b` ⇒ 本节读的还是**交付态字节**（本腿一枚产码字节未留，见末行收口）。
>
> **跑形（⛔ 非整包）**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test -count=1 -run 'TestTicket268' ./cmd/wisp/`
> —— PATH 注入两枚目录照 r1 §1 那把整包命令的形状（r1 台件逐字：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test -count=1 -v ./cmd/wisp/`），
> 只是把 `-run` 收窄到本票三枚。⚠ **不扩到 `TestTicket256`**（那是 r1 的跑形，会在本腿多带五枚用例＝多一倍负载；派单只给 `-run 'TestTicket268'` 一形）。
> `-v` 加在三处（各自具名）：负控第二发、MUT-2、MUT-3——为的是让"没红的哪几枚"在名册里**逐枚可见**，而不是靠"非 `-v` 不印＝绿"这一步推理收口；
> MUT-1 那一发取的是派单给的裸形（无 `-v`）：那一形下 Go 仍逐枚印 `--- FAIL`，所以"②③ 未列红"是本腿现量（红名册计数＝1 枚）而不是推理，
> 但"②③ 各印过一次 PASS"本腿在那一发**没量到**——由下一发 MUT-2 的 `-v` 名册补上（那一发 ① 逐枚 `--- PASS` 可见）。
>
> **负控起跑（交付态、一枚未种，先证"真跑起来了"）**：
> - 非 `-v` 预跑 21:47:02→21:47:06：`ok github.com/CarlosShao/wisp/cmd/wisp 0.069s`、`rc=0`（台件 `v2-targeted-baseline.log`，53 字节）。
> - **`-v` 形同刻复跑 21:52:23→21:52:25**：三枚顶层 `--- PASS` 逐枚可见（`…NamesRefusedConfigApartFromMissingConfig (0.01s)`／`…ReachesStdoutOnce (0.00s)`／`…BySentinelNotByErrorWords (0.00s)`）、`PASS`／`ok … 0.083s`／`rc=0`；
>   台件 `v2-targeted-baseline-v.log`＝**5,957 字节**；`grep -c '^--- FAIL'`＝**0**；`grep -c 'wisp: resident \[risk\]'`＝**1**；
>   `grep -c 'provenance="defaults (config.toml present but refused at load)"'`＝**4**。
>   ⇒ **不是 `exit status 0xc000135`**（那一形会在无 `--- FAIL` 的情况下直接把包打 `FAIL`，且不会有 `--- PASS` 名册与 stdout 行）——本腿三枚是**实跑绿**，不是没跑起来。
> - 负载闸：每发起 go 命令前现量 CPU/MEM——21:45:40 CPU 34%／MEM 69.5、21:51:45 CPU 38%、21:52:11 CPU 26%／MEM 69.7、21:55:27 CPU 46%／MEM 70.2。
>   ⛔ 无任一发在 ≥70% 时起；21:55 那一采 MEM 触 70.2 ⇒ 该采之后**先落一笔 commit（纯文书，不占 CPU）再复量**才起下一发。
>   `@(Get-CimInstance Win32_Process -Filter "Name='go.exe'").Count`＝**0**（21:50:03 现量，同机无并发整包）。
>
> **三形 needle 预检（21:52:0x 现量，种之前先数，不符即 ABORT）**：
> `grep -c 'riskProvenanceRefusedAtLoad = "defaults (config.toml present but refused at load)"'`＝**1** ／
> `grep -c 'fmt.Printf("wisp: resident \[risk\]: config.toml is present but was refused at load'`＝**1** ／
> `grep -c 'if errors.Is(err, fs.ErrNotExist) {'`＝**1**。三形各自唯一命中 ⇒ 允许种。

| # | 种什么（只 `cmd/wisp/resident_approval_windows.go`） | 起 | 止 | 红了谁（本腿跑形） | 还原证 |
|---|---|---|---|---|---|
| MUT-1 | `:456` 新常量字面折回旧那枚：`riskProvenanceRefusedAtLoad = "defaults (config.toml unreadable)"`（＝AC#1 要求自证的"两枚状态又共用一张脸"） | 种 21:52:49（md5 现量 `e03e7d700f00d26f79186b3e99dc9888`，≠交付态 `69a630bf…`）／跑 21:53:14 | 21:53:17（`rc=1`） | **只用例① 红**：`--- FAIL: TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig`；②③ **未列红**（该跑形只 `-run 'TestTicket268'`，红名册整份三行 `:9-11` 全属①，`grep -c '^--- FAIL'`＝**1**）⇒ 与 r1 §5「只 ① 红」同形 | 21:53:47 `git cat-file blob HEAD:cmd/wisp/resident_approval_windows.go > cmd/wisp/resident_approval_windows.go` ⇒ `md5sum`＝`69a630bf350daf1032c62b921e9d7044` ＝ `git cat-file blob HEAD:… \| md5sum` **双读同值**；`git diff --name-only -- cmd/wisp`＝**0 行**；`grep -c "MUT-"`＝**0** |
| MUT-2 | 摘掉被拒支那一行 stdout——形＝把 `:510-512` 那枚 `fmt.Printf` 三行整段注释掉（⛔ 不留 `_ = err` 之类的补钉：那样会多一枚与本票无关的形状；注释后 `err` 仍被上一行 `slog.Warn` 的 `"err", err` 使用 ⇒ 编译干净，`go test` 能起跑即为证） | 种 21:58:58（md5 现量 `90c4ad0a04fe5d500185188f3f562ec6`）／跑 21:59:08 | 21:59:13（`rc=1`） | **②③ 双红、① 绿**（`-v` 形逐枚可见：`:--- PASS` 只有 ①，红名册 `grep -c '^--- FAIL'`＝**2**）——② 那三句量的正是"stdout 恰 1 行"＋两枚必带文案，③ 的 `fmt.Printf` 计数尺从 1 掉到 0 ⇒ **② 不是恒绿装饰**；与 r1 §5 MUT-2「②③ 双红」同形 | 21:59:28 `git cat-file blob HEAD:… > …` ⇒ `md5sum`＝`69a630bf350daf1032c62b921e9d7044`＝HEAD blob 双读同值；`git diff --name-only -- cmd/wisp`＝**0 行** |
| MUT-3 | 在哨兵旁边**加**一次 prose/`err.Error()` 字符串匹配、行为逐字不变——证 ③ 咬的是"靠文案分类"这件事 | 「待验」 | 「待验」 | 「待验」 | 「待验」 |

**红句逐字**（原文抄自本腿台件，行号＝台件行）：

MUT-1（`v2-mut1-red.log:9-11`，三句同发）：

```
resident_approval_risk_268_windows_test.go:125: the two shapes are folded back into one word "defaults (config.toml unreadable)" - that fold IS ticket 268: the user wrote a number, the gate runs the compiled 300s, and the receipt blames a file that is sitting in that directory
resident_approval_risk_268_windows_test.go:129: refused-file provenance "defaults (config.toml unreadable)" does not say out loud that the file is present, which is the whole difference from the missing shape
resident_approval_risk_268_windows_test.go:147: provenance constants riskProvenanceUnreadable and riskProvenanceRefusedAtLoad share the literal "defaults (config.toml unreadable)" - the four readings have to be four names
```

MUT-2（`v2-mut2-red.log:12/13/14/19`，四条同发）：

```
resident_approval_risk_268_windows_test.go:201: the refused branch wrote 0 stdout lines carrying "wisp: resident [risk]:", want exactly 1. AC#1 asks that this shape reach a user's eye at least once from a terminal; one line is the shape resident_windows.go's [hotkey] fallback has carried since ticket 258
resident_approval_risk_268_windows_test.go:205: the stdout line has to carry the loader's own answer, not just this file's class name; captured stdout was:
resident_approval_risk_268_windows_test.go:208: the stdout line has to say what the gate runs on instead; captured stdout was:
resident_approval_risk_268_windows_test.go:285: residentRiskGateValues holds 0 fmt.Printf calls, want exactly 1: AC#1's user-visible face is the other half of this branch, and a branch that only logs to disk does not reach the eye of anyone who launched the process from a terminal
```

MUT-3（`v2-mut3-red.log`）：「待验」

**MUT-1 附带的现场名实不符（盘上日志原文，非转述）**：`v2-mut1-red.log:4` 被拒支印
`msg="resident gate: [risk] source present but refused at construction; …" provenance="defaults (config.toml unreadable)"`——
名字（"unreadable"）与事情（文件在场、被值域门拒）对不上，正是本票题面那句话的机器形状；`:5` 那一行 stdout 仍写 "is present but was refused at load" ⇒ **折叠态下两枚面自己先打起来**。还原后同一支印 `provenance="defaults (config.toml present but refused at load)"`（`v2-targeted-baseline-v.log` 4 处）。

**stdout 那行的跨发计数——一把有坑的尺，本腿换形并具名**：§3/r1 用的 `grep -c 'wisp: resident \[risk\]'` 在**红日志**里会把**测试自己的红句**数进来（`②:201` 那句逐字引用了同一枚字面串），所以它只在全绿态才等于"产码印了几行"。收紧尺＝只认行首的产码句 `grep -c '^wisp: resident \[risk\]: config.toml is present'`：
基线 `-v`＝**1**（`:6`）／MUT-1＝**1**（`:5`——这一发只折名字、stdout 面照旧，正好印证 ②③ 量的不是字面）／MUT-2＝**0**（那一行被真摘掉）。⚠ 本腿不据此说 §3 那两枚读数（基线 0／交付 1）错了：那两把尺跑的是**全绿整包**，没有红句混进来。

**与 r1 §5 名册的可比性（口径先说清，不冒充同尺）**：r1 跑形＝`-run 'TestTicket268|TestTicket256'`，本腿＝`-run 'TestTicket268'`（派单收窄）。
⇒ 本腿**量不到**"票 256 五枚钉在这一发突变下是否仍绿"那一半（那是 r1 §5 三行都点名的隔离证），只复认三行的**红名册方向**：
MUT-1 只红①（r1 同）／MUT-2 ②③ 双红（r1 同）／MUT-3 只红③（r1 同）。红句文本本腿与 r1 **逐字同**（三处都对得上台件行，见上）。

**三形全种完后的收口（本节末行，逐格填）**：MUT-3 还原后 `md5sum` 双读同值＝`69a630bf350daf1032c62b921e9d7044`、
`git diff --name-only -- cmd/wisp`＝**0 行**、测试文件 md5 未变＝`e2405a21bb278b558c1305ea95d34a9b` ⇒ **产码零净变化**、「除还原外未动测试文件与产码」成立。

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

四门都在**本腿自己的窗口**现跑（20:30:05–20:30:55，`date` 现量逐枚带时刻），台件四枚全在本目录。
起手 HEAD＝`164ee2c2`，本节跑完后又落了一笔 `01b38812`（§3）——**四门跑在 `15d8b60e` 这一枚 HEAD 上**（20:26 的台账笔，`git rev-parse --short HEAD` 现量），
两枚产码字节面本节现量仍与 r1 交付同值（见 §0.1 与下条末行）。

| 门 | 命令原文 | 本腿读数 | 时刻 | 与 r1 末笔（`2925b643`／17:19 那两采＋15:25–15:26 那四采）比 |
|---|---|---|---|---|
| 1 | `GOFLAGS= sh scripts/d22scan.sh` | **rc=0**、`d22scan: clean - no D22 ban violations`；正控真跑：`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0`；`examined 266 production Go files under internal/ and cmd/`；八枚 scope＝`#1-5 internal/=228`／`#1-5 cmd/=38`／`#6 frontend/=85`／`#7 internal/tools/=23`／`#8 design/=39`／`#8 frontend/=85`／`#8 internal/=512`／`#8 cmd/=104` | 20:30:05–20:30:31 | **八枚逐枚同读数**；红名集合＝**空**（`grep -E 'VIOLATION\|violation'` 8 处命中全在 `TestBuiltBinaryGoesRedEndToEnd` 的**自造 fixture** 里，`:141-219` 缩进块，非本仓真实文件；真实一面 `:251` 逐字 `clean`） |
| 2 | `sh scripts/check-path-length-budget.sh --with-self-test` | **rc=0**、`positive control PASSED`（三发控制 1/3–3/3 逐枚 ok）、`VERDICT GREEN`；`over-budget=57 covered by roster=57 not in roster=0`；`longest=180 chars relative`＝票 252 那枚工单名（⛔ 不是本腿造的）；`worst full path …=224 chars` | 20:30:38–20:30:41 | 分子分母**不一致**：r1＝`tracked paths=5953`，本腿＝`tracked paths=5991`（**+38**）。⇒ **先怀疑在谁的窗口**：见下段归因，**57／57／0 那三枚一枚未变**，本腿两枚文件不在 over-budget 名册里（`grep resident_approval_risk_268`＝0 命中） |
| 3 | `GOFLAGS= go vet ./cmd/wisp/` | **rc=0**，输出 **0 字节**（`wc -c` 亲量；⛔ 未用 `2>/dev/null`，空读数与非空分得开） | 20:30:54 | **同读数**（r1 13:33:01／15:26:37／16:55:15 三采都 rc=0 无输出） |
| 4 | `"D:/work/base/gopath/bin/gofumpt.exe" -version` 先自证尺活着，再 `-l cmd/wisp` | 尺＝**v0.12.0 (go1.27.1)**；`-l` 输出＝**只有 `cmd\wisp\models.go`** 一枚（`cat` 亲量，1 行）＝派单点名的**预存脏枚，不是本腿的红**，具名即可。本腿两枚与 r1 两枚产码**都不在名册里** | 20:30:55 | **同读数**（r1 13:33:03／14:09:34／15:26:39／16:55:21 四采同枚） |

**门 2 分母漂移的归因（现量，不猜）**：`git ls-tree -r --name-only` 逐枚数——
`9a941965`（r1 落地笔，13:31）＝**5953**（＝r1 报的数）／`b93624d4`＝5970／`2925b643`＝5975／`2a633eb8`＝5981／`164ee2c2`＝5990／`15d8b60e`＝**5991**＝本腿 `git ls-files | wc -l` 现量。
`git diff --name-status 9a941965 15d8b60e`（门 2 取数那一刻的 HEAD）＝**38 枚新增（`A`）＋9 枚修改（`M`）＝47 枚路径**，
按顶层目录分（对全部 47 枚）＝`.scratch/` **44**、`docs/` **2**（`docs/evidence/s1/265-267-evidence-index.md` 新增、`docs/reports/pending-and-issues.md` 修改）、`scripts/` **1**；
落 `cmd/`＝**0 枚**、落 `internal/`＝**0 枚**（尺＝`git diff --name-only 9a941965 15d8b60e | grep -cE '^(cmd|internal)/'`＝**0**，20:3x 现量）。
⇒ **漂移全在文书／台件地界，属 r1 交件之后那批腿与台账笔（`evidence-close-1`／`111-*`／`167-*`／`231-a1`／`268-v1` 自己那三笔）加的台件，不属 268 产码射程**；
`over-budget`／`roster`／`not in roster` 三枚读数一字未变＝门本身没有因为任何一腿扩大射程。⚠ 本腿 §3 那 8 枚名册台件在门 2 取数时**未跟踪**（`:1` 尺只数 `git ls-files`），所以 5991 不含它们；
本节这批 commit 之后 tracked 数＝`git ls-tree -r --name-only HEAD | wc -l`＝**5999**（20:3x 现量）——下一枚读门 2 的腿会看到它，如实先写在这里。

**字节面同值自证（本节读数归属交付态）**：`md5sum` ＝ `git cat-file blob HEAD:… | md5sum` 双读——
`resident_approval_windows.go`＝`69a630bf350daf1032c62b921e9d7044`、`resident_approval_risk_268_windows_test.go`＝`e2405a21bb278b558c1305ea95d34a9b`（20:3x 现量，两枚都与 r1 §4/§5 末采报的同值）。
⇒ 四门读的是**交付态字节**，与 r1 那四采可比。


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

派单口径：**只交事实与复跑读数，裁决归编排者**。六处都出自 `268-r1` §7.3 自己列的待裁清单（`r1/evidence.md:294-312`），
本腿每一处独立取数、不复用 r1 的转述。

### 6.1 新 provenance 取 **class 级**措辞（r1 §7.3 第 1 条）

- 本腿独立量的"机读分不分得开"：`grep -rn 'errors.Is\|^var Err\|type .*Error' internal/config/*.go | grep -v _test`＝
  该包**只有一枚** sentinel 用法：`internal/config/parse.go:234` `return errors.Is(err, fs.ErrNotExist)`；
  `internal/config/loader.go:74` 那行是**注释**（"reads fs.ErrNotExist off it for cause=missing"）。
- 值域拒载那支的造物现场＝`internal/config/validate.go:141-148`（亲读逐字）：
  `return observe.New(observe.ClassConfig, fmt.Sprintf("config.toml: risk.confirm_timeout_sec %d out of range [%d, %d]", …))`
  ⇒ **没有 marker 类型、没有包装 cause**，分不出"语法／未知键／迁移／band"任何两支而**不读 prose**。
- ⇒ **事实**：`present but refused at load` 是**今天诚实的上限**（class 级），"被值域门拒"那半只能由紧跟其后的 `%v` 原话承载
  （日志 attr `"err"` ＋ stdout 行 `:510` 两处都在，§5.1/§0.2 已亲读）；要 string 里出现 band 字样 ⇒ 得先给 `internal/config` 造 marker ＝ 跨包契约面。
  本腿**不裁**这算不算"具名程度够"。⚠ 一处支持 r1 的旁证：票面对 AC#1 的原句是"一枚与'文件不存在'**不同**的具名状态"，未要求具名到规则级。

### 6.2 d22scan `ban #8 cmd/` 分母 **103→104**（r1 §7.3 第 2 条）

- 独立尺＝git 名册（⛔ 不用 `find`，它会数进 `testdata`）：
  `git ls-tree -r --name-only <tree> | grep -E '^cmd/.*\.go$' | grep -vc testdata`
  → `87bc8b0f`（268 之前）＝**103**；`9a941965`＝**104**；`HEAD`＝**104**。
- `find cmd -name '*.go' -type f`＝**105** 与仪器读数差 1 的归因（现量，不是矛盾）：
  `cmd/wisp/testdata/esclistener/main.go` 被仪器跳过——`tools/d22scan/main.go:1117`
  `if (sc.goOnly || sc.everyFile) && d.Name() == "testdata"`（`:682`/`:1031` 同族跳过）。
- 违规枚数＝**0**、红名集合＝**空**（门 1 §4 现量，`:251` 逐字 `clean`）。
- ⇒ **事实**：+1 那枚**确为** `cmd/wisp/resident_approval_risk_268_windows_test.go`（唯一一枚随 268 落地的新 Go 文件，`git show --name-only 9a941965` 亲读 5 枚路径）。
  这是"仪器射程面多一枚文件"，不是"违规多一条"。**"算不算读数扩大"本腿不裁**（派单原文是"卫生四门读数不扩大"，红名集合逐名比对＝空集差）。

### 6.3 AC#1 的 stdout 那行算不算**达标档**（r1 §7.3 第 3 条）

- 票面 AC#1 原文逐字（`.scratch/wisp/issues/268-…md:31`）：
  「日志具名 provenance 算**最低档**；`doctor`／**回执**算**达标**」。派单与 r1 把 stdout 那行记为达标凭据。
- 本腿量的可达面事实（三枚，逐枚 file:line）：
  1. `slog.Warn`（被拒支在 `:496`，进 logsink 落盘，两枚起法都在；缺失支那枚在 `:477`）＝票面写明的**最低档**。
  2. `fmt.Printf` 那行（`:510`）＝**终端起法可见、`-H=windowsgui` 双击无 console 可打**（r1 §6.1 已具名，本腿不推翻）。
  3. `doctor` **零改动**：`grep -c 'pass("\|fail("\|info("' cmd/wisp/doctor.go`＝**28**（与票面 AC#0④ 的"调用点 28"同读数），
     且 `9a941965`／`4c456d0a`／`b93624d4`／`2a0e23b4`／`c072d5b7` 五笔对 `doctor` 命中逐笔＝**0**（`git show --name-only | grep -c doctor`）。
     成本面本腿复认 a2：`scripts/build.ps1:167-169` 逐字 `& (Join-Path $outDir 'wisp.exe') doctor` ＋ `if ($LASTEXITCODE -ne 0) { Fail 'wisp doctor reported FAIL.' }`
     ⇒ 加 `fail` 级检查项＝把带内越界那台机器的构建门打红。
- **"回执"这枚档到底指哪一面**——本腿量的现量：`grep -rn '回执' cmd/wisp/*.go | grep -v _test`＝5 处，
  全属 **run 腿 firstrun 的 stderr 回执**（`firstrun.go:88/126/147/179`，且 `:147` 逐字"三段都只走 stderr 这一条'给人看的回执'通道"）
  与 **panel-inbound 的行回执**（`panel_inbound.go:181`），**没有一枚是常驻审批卡的回执**。
  ⇒ 事实＝票面那枚"回执"**在常驻腿今天没有承接面**（要造＝动 `frontend/**` 或面板回执，票面「排程与禁区」末条＋AC#0③已划为停手项）。
  本腿另量到常驻腿**本来就有**第二枚 stdout 面这一事实（不属 268，仅供裁档参考）：
  `resident_approval_windows.go:538` `fmt.Printf("wisp: [audit] %s\n", line)`（`residentAuditf`），与 `[hotkey]`（`resident_windows.go:187`）同形。
- ⇒ **本腿不裁分档**；只把"最低档已满足／达标两档里 doctor 零改动且有构建门代价、回执面无承接位"三枚事实交回。

### 6.4 硬钉 1 的**射程**（r1 §7.3 第 4 条）

- 仪器本体亲读（`cmd/wisp/resident_approval_risk_256_windows_test.go`）：AST 窗从 `:537` `ast.Inspect(fn, …)` 起，
  `fn` 是 `:530` 那句 `no runResident in %s` 锚定的**唯一一个 FuncDecl**；三枚断言＝`:560 if fed != 1`／`:563 if bare != 0`／`:568 if !argOK`（行号现量逐中）。
  它数的是**构造名**（`newResidentApprovalWithConfig`／`newResidentApproval`）与那枚实参是否 `*.DataDir`——
  ⇒ **仪器射程＝a2 §2④ 说的"构造名计数钉"，不数 `config.LoadFile`**。复认 a2 的尺**形**成立。
- 同时本腿 §5.1 已枚过全部读取点：现量 7 枚（其中 1 枚是 `run.go:17` 注释），落 `residentRiskGateValues` 内部的只有 `:474` 一枚 ⇒
  **268 一个读取点都没加**（复现 §5.1，非引用 r1）。
- ⇒ 事实＝**派单那句更严的纪律（不得加第二枚 `[risk]` 读取点）与 a2 收窄口径在这份 diff 上同读数**，冲突只存在于文字层。
  本腿**不裁**该以哪句为准。

### 6.5 `:379-382`（票面）vs `:379-384`（派单／r1）文本区间（r1 §7.3 第 5 条）

- 独立复跑（`git cat-file blob 87bc8b0f:… | sed -n '379,385p' | cat -n`，票面锚定的就是这枚旧文件）：
  `:379`＝`slog.Info("resident gate: [risk] tier taken at construction",`、`:380` provenance、`:381` config_path、
  `:382`＝`"window_sec_read", …, "confirm_timeout_sec_read", …`（**不是句尾**）、`:383` gate_window／gate_queue_timeout、
  `:384`＝`"scope", …`（句尾 `)` 在此行）、`:385`＝`ra.cards = approval.NewReplies()`（**已出句**）。
- ⇒ 事实＝**旧文件该句实跨 `:379-384`**；票面 `:379-382` 截短两行；派单与 r1 的 `:379-384` 与实况一致；
  **代码无漂**（`9a941965^` 的 `:379-384` 逐字未动）。现文件（含 268）同一句整体下移 1 行＝`sed -n '380,386p'` 逐字同形（本腿现量）。
- ⇒ 这一处 r1 的"以亲读为准"本腿**复现且无新增疑点**；是否据此改票面文字归编排者（AC 框本腿一枚未碰）。

### 6.6 触发面窄于 a2（r1 §7.3 第 6 条）

- r1 的尺原文逐字复跑（`grep -rn 'confirm_timeout_sec = \|l1_window_sec = ' cmd/wisp/*_test.go`）：
  本腿 20:2x 现量＝**16**，⚠ **不等于 r1 报的 14**。归因（同一枚尺过 git 名册，非猜）：
  `git grep -h -E '…' 87bc8b0f -- 'cmd/wisp/*_test.go' | wc -l`＝**14**；同尺对 `9a941965`＝**16**。
  ⇒ **r1 那 14 是对"它落地之前"的树取的，落地后＋2 枚＝它自己那枚新文件（`_268_windows_test.go:82` 种 20、`:176` 修成 45）**；
  两读数都真，差别＝取数时刻，**不是漂字**。本腿按派单口径具名"读数不一致先怀疑在谁的窗口"：这里在 r1 自己的落地笔。
- "带外 confirm **0 枚**"那半独立复量：现役种子里 confirm 取值＝`45/90/31/40/%d` 与 `20`（`20` 只属 268 自己的用例）。
  `%d` 那枚（`approval_reply_201_test.go:145`）的喂值现场＝`newReplyHost(t, …)` 六处调用：`:202/:272/:350/:521/:595`＝`40s`、`:439`＝`31s` ⇒ 全带内。
  ⇒ **除 268 自己那三枚用例，今天整包里没有任何一枚用例把被拒支走绿**＝与 §3 那枚"stdout 行在整包里只出现 **1** 枚、
  且归属 `TestTicket268ResidentGateNamesRefusedConfigApartFromMissingConfig`"（`:1002`，20:11:50.962）**互相印证**。
- ⚠ 本腿顺手量到一枚**他票留下的过期注释**（不属 268 射程，只具名归口）：
  `cmd/wisp/panel_pump_test.go:124` 注释逐字仍写 `confirm_timeout_sec = 2`，而 `:136` 代码已是 `confirm_timeout_sec = 31`
  （改动笔＝`32e74479` 票 267-r2 种子迁移，12 枚带外喂值点抬进带内那一笔）⇒ 注释过期、代码带内；
  本腿一字未改（⛔ 零产码），交回编排者决定归哪枚票。
- a2 §4 那句"**即使触发也不撞**"本腿不复跑它的钉名册（属 a2 射程，且并行腿地界），只复认它的**结论方向**与 268-r1 的"根本不触发"在同一枚事实（除 268 用例无带外种子）上不打架。


## §7 commit 链与收尾读数

「待验」
