# 268-a2 落地腿开工前的最后一格料：往常驻腿 stdout 加一行会不会撞别人／那句"被拒"能不能复用现成的词／它与票 128 AC#5 分不分得家／`DefaultL1Window=3s` 那枚字面是谁的第二真相源

标注口径：**〔编排者复跑〕不存在于本件**——本件是腿（`268-a2`）自己落的，逐条都带我这轮现量的尺（§5）；凡"今天红／今天绿"一律写**〔预测〕**并归口给编排者（硬闸①：本腿零 `go` 命令）。

## §0 起手锚

- 起手时刻 `11:00:46+0800`／起手 HEAD **`c6cf66e6`**（`git rev-parse --short HEAD` 自取，未抄派单）／分支 `dev`／共享工作树。
- 起手 `git status --porcelain -- cmd internal tools scripts docs` 读数＝**0 行（全净）**〔11:00:46 现量〕；同刻全仓 `git status --porcelain | wc -l` ＝ **671**（脏面在 `.scratch/**`、`docs/**` 之外那批，与本腿爆破半径无关）。
- ⚠ **并发漂移（本腿跑中现量）**：HEAD 漂到 **`69c1bb82`**（`ledger(A616)`，11:16:13 现量）；**同一刻 scoped 读数从 0 变 2**——` M cmd/wisp/firstrun.go`＋`?? cmd/wisp/firstrun_257_test.go` ⇒ **`257-r2` 此刻正在写 `cmd/wisp`**（本票面"⛔ 三枚串行"点名的就是它）。本腿一枚字节未碰这两文件；本件所有 `cmd/wisp/**` 行号是 `c6cf66e6`→`69c1bb82` 之上读到的，落地腿开工时**必须重取行号**。
- 硬闸遵守：零 `go` 命令（未跑 build/test/vet/gofumpt/slo/scripts，逐枚具名见 §5）；⛔ 未读、未转述 `frontend/**`／`design/**`；未改任何已跟踪文件；未碰任何票面 AC 框；未碰台账。
- 本件引用的产码文件：`cmd/wisp/resident_approval_windows.go`·`cmd/wisp/resident_windows.go`·`cmd/wisp/config_reload.go`·`cmd/wisp/doctor.go`·`internal/config/loader.go`·`internal/config/validate.go`·`internal/config/unwired.go`·`internal/config/parse.go`·`internal/observe/redact.go`·`internal/observe/errors.go`·`internal/agent/approval/gate.go`·`internal/agent/approval/queue.go`。

## §1 Q1 名册：断常驻腿 stdout／console 的既有钉，逐枚带断言条件逐字

★ 先给**可达性**这把尺（它决定后面所有判断）：那句新增 `fmt.Printf` 只在 `residentRiskGateValues`（`cmd/wisp/resident_approval_windows.go:447-462`）的**被拒那一支**打，而该函数唯一的产码调用点是 `resident_windows.go:132`（`newResidentApprovalWithConfig(rt.Layout.DataDir)`，`runResident` 内）。现量：**今天套件里没有任何一例把带 `config.toml` 的数据根交给常驻腿**——所有子进程常驻用例的数据根都是 `t.TempDir()` 空目录（`resident_ball_228_windows_test.go:55`、`resident_approval_246_windows_test.go:186`、`resident_task_source_246_windows_test.go:298`、`resident_sink_nail_127_windows_test.go:418/:527/:569` 逐枚 Read），且 `resident_task_source_246_windows_test.go:389-391` 反过来**钉"常驻腿不许造出 config.toml"**。⇒〔预测〕被拒句在既有套件的**任何一发里都不会被打出来**，也就无从打红。**同一条顺序还有一层含义**（尺 R-18，11:27 现量）：全仓产码里唯一造 `config.toml` 的入口是 `ensureFirstRunConfig`（`cmd/wisp/firstrun.go:72`），它唯一的产码调用者是 `cmd/wisp/run.go:245`——在常驻进程里那要走 `resident_windows.go:260` 的 `startResidentTaskSource`，**晚于 `:132` 的风险读**。⇒ 空数据根上常驻腿永远先看到"缺失"形；"被拒"形只有一枚来源＝**用户手改过 `config.toml`**（恰是本票要服务的那一枚形状）。

### 1.1 contains 型（加一行不红；逐枚）

| # | 落点 | 断言条件逐字 | 打的是什么面 | 加一行 stdout 会不会红 |
|---|---|---|---|---|
| C1 | `cmd/wisp/dataroot_128_test.go:168-171` | `for _, marker := range rescueMarkers128 { if !strings.Contains(err.Error(), marker) {` | ⚠ **不是 stdout**，是 `resolveDataDir` 返回的 error 串 | 不红（与 stdout 无关；a1 引这枚当"contains 型 stdout 钉"是**口径错**，见 §7(2)） |
| C2 | `cmd/wisp/dataroot_128_windows_test.go:74-78` | `for _, marker := range leg.markers { if !strings.Contains(out, marker) { t.Errorf("AC#2 RED: leg %q refused but its line never says %q:\\n%s", leg.name, marker, out) } }`；resident 行 `:63` markers＝`{"user config dir","%AppData% is not defined"}`；`out` 是 `:131-133` **stdout+stderr 合并缓冲** | 真进程合并 console | 不红。且该形下 `proc.Boot` 在 `resident_windows.go:42` 就失败 ⇒ `:52` 打 stderr、`:53 os.Exit(1)`，**走不到 `:132`**，新句根本不出现 |
| C3 | `cmd/wisp/resident_ball_228_windows_test.go:61-64`、`:82-85` | `pollUntil127(200, func() bool { out := leg.stdout.String(); return strings.Contains(out, ballUpClaim) \|\| strings.Contains(out, ballAbsentClaim) })`；`if !strings.Contains(out, residentReachedLoop)` | 常驻子进程 stdout | 不红（在场型） |
| C4 | `cmd/wisp/resident_approval_246_windows_test.go:190-194`、`:204`、`:216-228`、`:261` | `posture && strings.Contains(out, cancelStepRosterClaim)`；`if !strings.Contains(out, cancelStepRosterClaim)`；`switch { case strings.Contains(out, ballUpClaim): if !strings.Contains(out, gateAssembledClaim) … case ballAbsentClaim: … gateAbsentClaim … }`；`if !strings.Contains(leg.stdout.String(), "10 steps, 0 failed")` | 常驻子进程 stdout | 不红——**前提**：新句不得含 `ballUpClaim`/`ballAbsentClaim` 那两枚字面（值＝`resident_ball_228_windows_test.go:47-48`：`the floating ball window is up in this process`／`this process has NO floating ball window`） |
| C5 | `cmd/wisp/resident_task_source_246_windows_test.go:302-304`、`:311-316`、`:324-327`、`:368-390`、`:428-445` | 四处**缺席型** contains（共涉 5 枚字面）：`if strings.Contains(out, taskEntryInjectedClaim) \|\| strings.Contains(out, taskEntryConsoleClaim)`（`:311`）、`if strings.Contains(out, pipelineAbsentClaim)`（`:324`）、`if strings.Contains(out, taskEntryDisabledClaim)`（`:379`）/`taskEntryRefusedClaim`（`:382`）、`if strings.Contains(out, taskEntryInjectedClaim)`（`:443`）；在场型 `任务来源：`＋`taskPostureAbsent/Refused`（`:314`、`:385`） | 常驻子进程 stdout | 不红——**前提**：新句不含那 5 枚字面（逐名值与声明行见 §8 第 2 行末）。这四枚缺席断言位是本寸地最脆的一族——它们断的是"降级两支不许折成一句"，而本票恰恰在治另一处折叠，落地腿若把两族话混写会撞 |
| C6 | `cmd/wisp/resident_sink_nail_127_windows_test.go:487`、`:497`、`:529-542`、`:637` | `if !leg.stderr.has(residentInstallMsg)`；`if !leg.stdout.has(residentReachedLoop)`；`pollUntil127(200, … leg.stderr.has(residentRefusalPrefix))`＋`leg.stderr.has(residentRefusalPromise)`；`if !leg.stdout.has(residentShutdownLine)` | 常驻子进程 stdout＋stderr | 不红 |

### 1.2 枚数型／相等型（这才是"加行会红"的候选，逐枚读完）

| # | 落点 | 断言条件逐字 | 数的是什么 | 三支判断（数什么／会不会进新句／结论） |
|---|---|---|---|---|
| N1 | `cmd/wisp/resident_ball_228_windows_test.go:74` | `if up != 0 && absent != 0 {` | 两枚球字面各出现的枚数（`:66-67` `strings.Count(out, ballUpClaim)`/`(out, ballAbsentClaim)`） | 数的是**指定字面出现次数**；新句不含那两枚 ⇒ 不进分子；**不红** |
| N2 | 同文件 `:78-81` | `if up+absent != 1 { t.Errorf("AC#1 RED: the ball posture is stated %d times, want exactly 1 …")` | 同上，两枚球字面之和 | 同上 ⇒ **不红**（它是"球姿态只许说一次"，不是"console 只许有 N 行"） |
| N3 | `cmd/wisp/resident_sink_nail_127_windows_test.go:493` | `if n := strings.Count(leg.stderr.String(), residentEarlyResolverMsg); n != 1 {` | **stderr** 里 `winsec: sealing path resolver installed`（`:137`）的枚数 | 数的是 stderr 指定 msg；stdout 新句不进 ⇒ **不红** |
| N4 | `cmd/wisp/early_log_nail_130_windows_test.go:257`、`:266` | `if n := strings.Count(stderr, early130ResolverMsg); n != 1`；`if n := strings.Count(stderr, early130InstallMsg); n != 1` | stderr 两枚指定 msg | 同上 ⇒ **不红**（落地腿若**新增一枚 `slog.Warn`** 才会进 stderr 的视野，但那两枚 Count 仍只数各自的串） |
| N5 | `cmd/wisp/firstrun_198_test.go:91-93` | `if stdoutText != "" { t.Errorf("stage 1: stdout = %q, want empty …")` | **注入 buffer**（`run198` 在 `:34` 交 `stdout:  out` 给 `runSpec`，打的是 `wisp run` 那条腿） | 数的是 run 腿注入流；`fmt.Printf` 走 `os.Stdout` **进不去注入 buffer**，且 run 腿不调 `residentRiskGateValues`〔尺见 §5 R-9〕 ⇒ **不红** |
| N6 | 全仓 `cmd/wisp` 的**逐枚相等／整缓冲相等型 stdout 钉** | 尺 R-5b（宽正则，见 §5）读数＝**除 N5 外零枚** | — | **零枚"console 只有这几行"型整等钉**。⚠ 我先用窄尺 `(out\|stdout\|said) != ""` 得出"2 枚"，那两枚（`resident_task_source_246_windows_test.go:253`、`:275`）Read 后是 `residentTestTaskText` 的**返回值**，不是 stdout ⇒ 已改判（§7(1)） |
| N8 | **进程内 capture 型 stdout 钉（`captureStdout128` 的全部调用者）** | 尺 R-19 读数＝**7 处**：`cmd/wisp/dataroot_128_test.go:132`（包的是 `cmdDoctor()`）、`cmd/wisp/resident_cancel_key_wording_260r3_windows_test.go:98`、`:140`、`:164`、`:178`、`cmd/wisp/resident_cancel_key_label_260r4_windows_test.go:88`、`:130` | **窗口体只包 `bindBallHost(...)`**，构造调用全在窗口**之外**（`ra2 := newResidentApproval()` 逐枚在 `:96`／`:138`／`:161`／`:176`／`:86`／`:128`／`:164`） | 本票新增句打在 `residentRiskGateValues`（经构造调用）⇒ **落在窗口外**，〔预测〕不红。⚠ 但这些窗口里有**缺席断言**（`260r3:132/147/170/184/189` 断 `Esc` 与种子键名不出现在捕获文本里）⇒ 落地腿若把新句挪进 `bindBallHost`／`residentStatusLine` 那条面，就会进这些窗口并被"不含 Esc"这类断言咬到。**判语：新句只在构造期那一次打，不许搬到装配回执里。** |
| N7 | `len(lines)` 型行数钉 | 尺 R-6 读数 **8 处**，逐枚 Read 后无一处数 wisp 的 console 行数（`config_receipt_255_test.go:203-204` 数被引用的**源码文件**行数、`:617`／`config_reload_223_test.go:220`／`resident_approval_246_windows_test.go:290`／`resident_task_source_live_246_windows_test.go:436` 是"从尾找含 needle 的行"的 helper（只 `t.Fatalf` "no line for %q"＝在场型）、`panel_pump_test.go:80` 数账本行、`panel_host_gate_test.go:360-361` 数的是 **git 的 stdout**） | — | **打在常驻腿 stdout 上的行数钉＝零枚** |

### 1.3 名册型（盯"常驻腿到底装配了哪几句／几处"的钉，⚠ 这几枚落地腿真会撞）

| # | 落点 | 逐字形状 | 落地腿怎么就红了 |
|---|---|---|---|
| M1 | `cmd/wisp/resident_approval_risk_256_windows_test.go:124-125`＋`:149-151` | `{"no host view at all", "no-view", riskProvenanceNoView}`、`{"config.toml missing", "unreadable", riskProvenanceUnreadable}`；`if ra.riskProvenance != s.prov { t.Errorf("riskProvenance = %q, want %q: the fallback has to be SAYABLE …")` | 这是**唯一**读 `riskProvenance` 的测试面。若把 `riskProvenanceUnreadable` 改字面**又**让"文件缺失"仍走它 ⇒ 不红；若把缺失与拒载**折回同一枚串** ⇒ 本票 AC#1 要的突变正是让它红（`268` AC#2 明写 `:141-144` 保绿） |
| M2 | 同文件 `:174-176` | `if a, b := newResidentApproval(), newResidentApprovalWithConfig(""); a.riskProvenance != b.riskProvenance {` | 空 dataDir 两支必须同串——落地腿新加的是**被拒**那一支，别把它错分给 `""` 支 |
| M3 | 同文件 `:438-443` | `if literals != 1 { t.Fatalf("found %d approval.Options literals in %s, want exactly 1. …")`（`resident_approval_windows.go` 里 `approval.Options{…}` 语法树枚数） | 被拒支**不许再建第二枚闸门**；只加字面与一行 `Printf` ⇒ 不红 |
| M4 | 同文件 `:560`、`:563`、`:568-571` | `if fed != 1 {`（`runResident` 调 `newResidentApprovalWithConfig` 的 AST 枚数）／`if bare != 0 {`／`if !argOK {`（实参须是 `*.DataDir`） | ⚠ a1 的读数是〔仅腿报〕，本腿复跑为逐字。**不许在 `runResident` 里加第二枚配置读取点或第二枚构造点**；本票不需要 |
| M5 | `cmd/wisp/config_receipt_255_test.go:179-224`＋`cmd/wisp/config_readers_255.go` | `raw, err := os.ReadFile(filepath.Join(root, …))`、`if n < 1 \|\| n > len(lines)`、`got := strings.TrimSpace(lines[n-1]); if !strings.Contains(got, token) { …the evidence drifted… }`、`const citedRowsFloor = 8` | 名册引的 **18** 枚 `x.go:N` 字面（尺 R-7，含注释里的引用，未逐枚分"cite／注释"）**没有一枚指向 `resident_approval_windows.go`**；会撞的是 `resident_ball_windows.go:276`、`panel_resident_windows.go:207/:253`、`panel_host_windows.go:262`、`run.go:423/424/435/991/1014`、`models.go:163/183`、`logsink.go:149` 那几枚文件的**行号** ⇒ 落地腿若顺手"整理"那些文件＝红。只改 `resident_approval_windows.go` ⇒ 不红 |
| M6 | `cmd/wisp/config_receipt_255_test.go:230-260` | `got := scanReadFiles(t, root, field); if !slicesEqualSorted(got, want) { …the roster has not been re-adjudicated }` | 扫描只覆盖 **hot** 段；`[risk]` 是 **locked**（`internal/config/tiers.go:41` `"risk": "locked"`，`:83 "risk": true`）⇒ 被拒支不新增 hot 段读点 ⇒ 不红 |
| M7 | `cmd/wisp/dataroot_128_test.go:352-379`＋`:314-321` | `if fn.calls["resolveDataDir"] { out = append(out, name) }`、`if len(out) < 4 {`、`if len(missing) > 0 { t.Fatalf("AC#2 RED: %d function(s) resolve the data root …")` | 现量产码读者＝`doctor.go:105`、`models.go:110`、`providers.go:81`、`run.go:200`（尺 R-8，4 枚，与腿表等集）。⇒ **被拒支绝不许去调 `resolveDataDir`**（它已经拿到 `dataDir` 形参），一调就抬这枚分母、表里没有常驻腿行 ⇒ 当场红 |
| M8 | `cmd/wisp/leg_dispatch_gate_133_test.go:187`、`:218` | `if len(legs) < minLegs133 {`（`minLegs133 = 11`，`:148`）；`for _, r := range censusVsUsage133(pkg, legs)` | 只有**新增子命令／改 `main.go` 分派或 usage 块**才撞；本票 A＋B 不新增子命令 ⇒ 不红 |
| M9 | `cmd/wisp/resident_ball_228_test.go:17`（源Walk 型） | 注释逐字 "Why a source walk and not a runtime call: runResident() never returns" | 该类钉断的是**调用点在源码里存在**，加句子不动调用点 ⇒ 不红 |

### 1.4 ⛔ 顺带量到的两枚"字面形状"约束（不是测试，是 CI 门，落地腿一样会撞）

- `tools/d22scan/main.go:579-586` `emojiScopes()` 含 `{dir: cmd/, label: "cmd/", goOnly: true}` ⇒ **`cmd/**` 的字符串与注释都在扫**（注释不豁免、字符串不豁免）。可用的字面：箭头类 `U+2190–U+21FF`、带圈数字 `U+2460–U+24FF` 是刻意留的空隙；`✓`(U+2713)／`≤`(U+2264) 会被判违规 ⇒ **别在句子里写"≤300s""✓"**。
- `tools/d22scan/main.go:737-743`：`filepath.Clean|Abs` 在 `risk.PathResolver` 之外即违规。新句若要拼路径**只能沿用** `resident_approval_windows.go:451` 的 `filepath.Join(dataDir, configFileName)`。
- `tools/d22scan/main.go:799+` ban #9：**产码注释里引用一个盘上不存在的仓相对路径即红**。落地腿若写 `（见 docs/contracts/C17.md）` 这类＝当场红（AGENTS.md §4 已警告那批交付物不存在）。
- `internal/observe/redact.go:107-127`：`Attr` 对 `key="err"` 不落 `secretWords/audioWords/contentWords` 任一支 ⇒ `switch v.Kind()` 的 `KindAny` 且非 `[]byte` ⇒ `:125/:127 return slog.Attr{Key: key, Value: v}` **原样透传**〔本腿逐字复跑 a1 的机制判断〕。⇒ 拒载原句今天已经在盘上 JSONL 里；`fmt.Printf` 只是把它搬到 console，**不新增脱敏面**，也就不会撞 `logsink_test.go`（`:121` 只把 `os.Stderr` 换成 pipe 断它自己那枚 probe 串）。

### 1.5 winlive 那族今天到底跑不跑（决定"红给谁看"）

`cmd/wisp` 里 `//go:build windows && winlive` 的文件＝**6 枚**（尺 R-4），含 `resident_hotkey_live_258_windows_test.go:1`（它的 `:57`/`:99`/`:145` 断的是 `hotkeys from defaults`／`hotkeys from config`，**contains 型**，加行不红）。`.github/workflows/ci.yml` 与 `scripts/*.sh` 里 **`winlive` 命中 0**（尺 R-11）⇒ 那 6 枚在 CI 上从不参与，`ci.yml:487` 那步跑的是 `scripts/wisp-cli-tests.sh`（scope `./cmd/wisp/`，`:111`），CI windows 腿整包会跑 §1.1/§1.2/§1.3 里所有 `//go:build windows` 的用例。

### 1.6 ★ Q1 结论（三支齐给）

1. **A 支（换具名 provenance）**：只动 `resident_approval_windows.go:430-438` 的常量与 `:447-462` 的分支 ⇒ 命中的名册只有 M1/M2/M3；〔预测〕只要"文件缺失"仍归 `riskProvenanceUnreadable`、被拒归**新名**，`256_test` 全绿。
2. **B 支（多一行 stdout `fmt.Printf`）**：既有 console 钉**全是 contains/指定字面 Count 型**（C1–C6、N1–N5），整等型与行数型**零枚**（N6/N7）；进程内那 7 处 `captureStdout128` 窗口**都不包构造调用**（N8，尺 R-19）⇒〔预测〕**零枚红**。更硬的两条理由是可达性：既有套件无一例让被拒支成立（§1 表头），且造 `config.toml` 的入口在常驻腿里排在风险读**之后**（R-18）。
3. **代价形状**（不是红，是"说了两次"）：若落地腿**也**给"缺失"那支加一句 stdout，则同一枚启动会连打两枚近义句——`resident_windows.go:187` 的 `[hotkey]` 那句（`fmt.Printf("wisp: ball [hotkey]: config.toml unreadable (%v); …", err)`）本来就在缺失形下打。⛔ 但**不许**为此把两支折成一支（票 268 AC#3：一响一静是有名字的差别；`resident_hotkey_258_windows_test.go:244-250 Test258ProvenanceWordsAreThePrintedOnes` 把 `[hotkey]` 那三枚词当数据钉着）。按已裁形（只被拒支打）这条不落。

## §2 Q2 复用面：`describeReloadFailure` 逐支＋措辞名册

### 2.1 逐支表（`cmd/wisp/config_reload.go:354-409`，函数名逐字 `func describeReloadFailure(err error) string`）

| 支 | 行 | 触发条件逐字 | 类型判断 or 词面 | 输出字面（逐字，含收尾） |
|---|---|---|---|---|
| missing | `:356-359` | `case errors.Is(err, fs.ErrNotExist):` | **类型**（sentinel，经 `observe.Error.Unwrap()`，`internal/observe/errors.go:207-212`） | `cause=missing detail="config.toml 读不到：文件不存在（这一条只说缺失，不说语法、不说权限）。本次运行继续用内存里的旧配置；文件回来之后要它的 mtime 或大小变过才会被下一 tick 重读"` |
| permission | `:360-363` | `case errors.Is(err, fs.ErrPermission):` | **类型** | `cause=permission detail="config.toml 读不到：这个进程没有读它的权限（文件在，也读得开名字，只是不让读）。本次运行继续用内存里的旧配置；这一条不是语法错，也不是文件缺失"` |
| syntax | `:369-372` | `case strings.HasPrefix(d, "config.toml parse"):`（`d = oe.Detail`，`errors.As(err, &oe)` 于 `:365-367`） | **词面**（`internal/config/loader.go:102`、`:118` 手写 `observe.Wrap(observe.ClassConfig, peekErr, "config.toml parse")`） | `cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"` |
| unknown-key | `:373-376` | `case strings.Contains(d, "unknown key"):` | **词面**（`internal/config/parse.go:145-148 fmt.Sprintf("unknown key %q at line %d", …)` 拼进 `observe.New(ClassConfig, "config.toml: "+…)`） | `cause=unknown-key detail="config.toml 语法没问题，但里面有这份 schema 不认的键（拼错的键会被这样拒绝，而不是被忽略）。本次运行继续用内存里的旧配置"` |
| migration | `:377-395` | `case strings.Contains(d, "cannot migrate from schema version") \|\| strings.Contains(d, "no migration registered"):` | **词面** | `cause=migration detail="config.toml 声明了一个这份 Wisp 不会迁移的 schema_version（文件被原样留着，不会被重置）。本次运行继续用内存里的旧配置；升级 Wisp 或恢复备份才会读它"` |
| **invalid** | `:396-399` | `case strings.HasPrefix(d, "config.toml:"):` | **词面，且是一个冒号字节** | `cause=invalid detail="config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）。本次运行继续用内存里的旧配置"` |
| unclassified | `:406-408` | 兜底 `return fmt.Sprintf("cause=unclassified detail=%q", …+err.Error())` | — | `config.toml 没有被重读成功，原因没有归入已知四类（缺失/语法/权限/schema）；本次运行继续用内存里的旧配置: <原句>` |

`validate()` 的短路顺序（`internal/config/validate.go:24-45`，13 支取第一枚非 nil）＝`App→Ball→Audio→Privacy→Models→Risk→Unwired→Cost→Observe→Net→APIKeyRefs→RolesIntensity→Catalog`；被拒那枚字面逐字在 `:145-148`（`config.toml: risk.confirm_timeout_sec %d out of range [%d, %d]`），其 `Err == nil`（`observe.New` 不装 cause，`internal/observe/errors.go:179-181`）⇒ **`errors.Is(err, fs.ErrNotExist)` 为假**＝A 支的分叉成立。

### 2.2 能不能直接复用它的 `invalid` 支——**判语：不复用那枚句子；可以复用"cause 词"的形状，且必须具名带价**

1. **复用整句＝造一枚谎**。七支收尾全含「**本次运行继续用内存里的旧配置**」（`:358`、`:362`、`:372`、`:376`、`:395`、`:399`、`:407` 逐枚在场，尺 R-10）——那是 **reload** 语境的真话（`config_reload.go:155` 是它唯一产码调用者：`rt.auditf("config: HOT-RELOAD state=not-applied %s", describeReloadFailure(err))`）。常驻腿这一读是**构造期**、没有"内存里的旧配置"可续（`resident_approval_windows.go:341-351` 与 `:441-446` 自陈"Construction time only / 回落编译常量"）。⇒ 把这枚句搬到 stdout 让用户读＝请用户相信一件没发生的事。
2. **复用它的分类＝把一枚已知在错的地方的词面判断抬成常驻腿的依赖边**（这条要具名，不许含糊）：`invalid` 支的判别位是 `strings.HasPrefix(d, "config.toml:")`，而带 `config.toml: ` 前缀的产码字面在 `internal/config` 里现量 **48 枚**＝`validate.go` 21／`catalog.go` 11／`parse.go` 8／`migrate.go` 4／`loader.go` 3／`unwired.go` 1（尺 R-12），**测试里引用这枚前缀的＝0 枚**（尺 R-13）⇒ 没有任何钉守这个约定。后果已现形：`internal/config/unwired.go:106` 那句 `"config.toml: %s is written but does nothing: …"`（"写了但不生效"）与 `internal/config/loader.go:121-123` 那句 `"config.toml: schema_version %d was written by a newer build …"`（版本问题）**都会**被 `:396` 判成"值不合法或引用解不开"。⇒ 常驻腿若走 `describeReloadFailure`，它给用户的 cause 词**今天就已经可能错**，且错得很安静（这恰好是票 268 要治的病，不能拿它当药）。
3. **改写那七句不行**：`cmd/wisp/config_sentences_223r2_test.go:22`（"CheckAndReload feeds describeReloadFailure exactly what reloadOnce does"）＋ `cmd/wisp/config_reload_223_test.go` 一族按支断言"其余几句不在场"（`config_reload.go:346-349` 注释逐字 "each case asserts the OTHER sentences are absent, so the four stay four"）⇒ 动那七枚字面＝撞 223 的互斥钉。**结论形状**：落地腿自己带一句（照 `resident_windows.go:187` 的**形状**＝同函数内 `slog.Warn`＋`fmt.Printf`，⛔ 不是照它的**字面**，那枚字面说的是 `unreadable`，正是本票要分开的东西），cause 词若要复用就复用**词**（`missing`/`invalid`）不复用**散文**。
4. 依赖边形状（技术层面不构成阻塞）：`describeReloadFailure` 是 `cmd/wisp` 包内函数，同包调用**不新增包级依赖边**（不触票 238 cut-1）⇒ 拦它的是上面 1/2/3，不是包图。

### 2.3 措辞名册：同一件事的第二套词，产码里各在哪、谁打印（落地腿只能从这里挑，不许新造）

| 族 | 落点（file:line） | 谁打印／打在哪 | 逐字形状（摘） |
|---|---|---|---|
| 配置未就绪（Unconfigured） | `cmd/wisp/run.go:417`；`cmd/wisp/models.go:240`、`:270`；`cmd/wisp/providers.go:100` | 跑任务的三条腿，`Fprintf(io.stderr(), …)`＋退码 2 | `"wisp run: 配置未就绪（Unconfigured）：%v\n"` |
| 配置未就绪（另一枚形） | `cmd/wisp/panel_inbound.go:232`（错误值，非打印）；`:148`（stderr） | 入向装配 | `"配置未就绪（%s）：%w"`；`"wisp panel-inbound: 入向装配未完成（Unconfigured）：%v\n"` |
| **unreadable 族（本票主场，三枚同形折叠）** | ① `cmd/wisp/resident_windows.go:185`（`slog.Warn`）＋`:187`（`fmt.Printf`＝**stdout，现成可抄的形状**）；② `cmd/wisp/resident_approval_windows.go:434`（常量 `"defaults (config.toml unreadable)"`）＋`:454`（`slog.Warn`，**零 Printf**）；③ `cmd/wisp/panel_resident_windows.go:203`（`slog.Warn("panel host: [panel] geometry source unreadable, sizing at the host's own default", …)`，**零 Printf**＝同族第三枚） | 常驻腿三条读取链 | `[hotkey]` 那句逐字：`"wisp: ball [hotkey]: config.toml unreadable (%v); the four hotkeys fall back to the compiled defaults (DefaultHotkeys)\n"` |
| 读不到 族 | `cmd/wisp/config_reload.go:358`、`:362`；`cmd/wisp/resident_task_source_windows.go:406`（`Fprintf(out, …)`）；`cmd/wisp/run.go:868` | reload 审计／常驻腿任务入口降级句／run 腿 | `"config.toml 读不到：文件不存在…"`；`"wisp: 这一路没有答复通道（装配时读不到可交互控制台）：%s\n"`；`"档位读不到：本进程不缓存任何上一次的宽松值…"` |
| 数据根无法解析 族（票 128 的词） | `cmd/wisp/doctor.go:282`（`errDataDirUnresolved = errors.New("数据根无法解析：用户配置目录不可得")`）＋`doctor.go:289-298`（`dataDirUnresolved128`＝**文案单点**，六枚 marker 由 `dataroot_128_test.go:91-98` 钉） | run/models/providers/doctor/secret 五枚串 | `"…数据根本应是 <用户配置目录>%s%s，而 Wisp 拒绝把它回落到当前工作目录…修法：Windows 把 APPDATA 设为一个可写目录…"` |
| 面板状态词 | `internal/panel/config_handlers.go:183`（`CredentialConfigUnreadable CredentialState = "config_unreadable"`）、`:300`（`panel: 配置不可读（requestId=%q）…`） | 面板回执（Go 侧字段／错误值） | — |

⛔ 落地腿从**第 3 行（unreadable 族）**里挑形状：`resident_windows.go:187` 是本仓唯一"同一函数里既 `slog.Warn` 又 `fmt.Printf`"的先例，也是已裁形点名的抄本。⚠ 它给的是**形状**不是**词**——那枚词叫 `unreadable`，正是本票 AC#1 要求分开的东西；新句必须说"被拒/越界"而不说"读不到"。

## §3 Q3 与票 128 AC#5 分家

**AC#5 原文要点**（`.scratch/wisp/issues/128-*.md:37-45`，勾框 `[ ]` **一枚未动**）：AC#2 判据逐字「三条腿都拒、**且拒绝原因可被人读懂**」，代价栏「本票落地的错误串三样都写」；现量是 `internal/proc/envfork.go:238` 逐字 `return Layout{}, fmt.Errorf("proc: user config dir: %w", err)`（本腿 Read 复认，行号现量）⇒ 常驻腿"拒了但读不出自救"。结案判据①补齐三样**且仍走同一枚 seam 的文案单点（不许再造第二份模板）**；②配一发常驻腿真进程 stderr 断言三样在场；③变异自证；④门禁四数＋名册差集（地界 `internal/proc`，"先协调"）。

**已落地的部分**：其余几条腿的自救句都出自 `cmd/wisp/doctor.go:289-298` 那枚单点；测试侧现状＝`cmd/wisp/dataroot_128_windows_test.go:59-63` 的表里**有** resident 行，但它的 `markers` 是 `{"user config dir", "%AppData% is not defined"}`（只有 OS 原话两枚），而 run/secret-list/doctor 三行用的是 `rescueMarkers128`（6 枚合取）；文件头 `:42-46` 注释逐字承认"the resident leg's sentence is produced by internal/proc's DefaultLayout and is checked for that instead"。⇒ **AC#5 未兑现的正是那 6→2 的差**。

**两枚会不会各造一句"配置没读到"——判：不会，理由具名到可达性，不是感觉。**
- 两枚句子落在**互斥的两条形**上：AC#5 那条形里 `%APPDATA%` 未设 ⇒ `runResident` 在 `resident_windows.go:42` 的 `proc.Boot(env)` 就拿到错，`:52` 打 `wisp: boot failed: %v` 到 **stderr**、`:53 os.Exit(1)`；**永远走不到** `:132`（本票那句的产地）。本票那条形要求数据根已解析、`config.toml` 在被拒——此时 `proc.Boot` 已成功。⇒ **同一次启动两句不可能同时在场**，各造也不互相遮蔽。
- 谁该用谁的词（具名）：**AC#5 必须用 `doctor.go:289-298` 的单点词**（它自己的判据①就禁第二份模板，且 `rescueMarkers128` 六枚 marker 钉在那套词上）；**本票必须用 `unreadable 族`／被拒的形状（`resident_windows.go:187` 的同函数 Warn＋Printf），⛔ 不许借 AC#5 那六枚 marker**——那六枚讲的是"数据根该落哪、APPDATA/XDG 怎么设"，与"你写的 20 被 `[31,3600]` 拒了"是两回事，借了就是把两枚病因折回一枚句子（本票的题面本身）。
- 共用面只有一处，且**不冲突**：AC#5 判据②的取证钉落点是 `cmd/wisp/dataroot_128_windows_test.go` 的 resident 行，它读的是 **stdout+stderr 合并缓冲**（`:131-133`）；本票那句走 stdout ⇒ 若 AC#5 把该行 markers 抬成 6 枚合取，`strings.Contains` 型**多几行无害**（§1 C2 同判）。
- **先后**：技术上**无依赖、可任意序**；真正会撞的是**地界**——两票都往 `cmd/wisp/**` 落（本票产码＋同包测试；AC#5 判据②落 `cmd/wisp/dataroot_128_windows_test.go`），而本票面 `:38-39` 写死"⛔ 三枚串行，不接受'改的是不同文件'这种推理"。此刻 `257-r2` 正在 `cmd/wisp` 里写（§0 现量）。⇒ 排程归编排者，不属本件能判死的格子（§6 U6）。

## §4 Q4 那枚字面：`DefaultL1Window=3s`

- 现量四枚常量（`internal/agent/approval/queue.go`）：`:107 DefaultApprovalTimeout = 300 * time.Second`、`:116 DefaultL1Window = 3 * time.Second`、`:120 MinL1Window = 2 * time.Second`、`:122 MaxL1Window = 3 * time.Second`。**钳位上界＝3 s，下界＝2 s**；钳位现场 `internal/agent/approval/gate.go:143-151`，逐字 `win := o.Window / switch { case win <= 0: win = DefaultL1Window / case win < MinL1Window: win = MinL1Window / case win > MaxL1Window: win = MaxL1Window }`。
- ⚠ **题面里那半句要更正**：`:456` 那句 `fallback` 写 `DefaultL1Window=3s` **与实现并不矛盾**——被拒/回落两支交给 `approval.New` 的 `Window` 是 **0**（`resident_approval_windows.go:449`、`:457`），命中 `gate.go:145-146` 的 `win <= 0` 支 ⇒ 真值就是 `DefaultL1Window`＝3 s。（"配置读通了才是 2 s"是另一条形：schema 的 `l1_window_sec` 默认 `2`，`resident_approval_risk_256_windows_test.go:204-209` 已把这条写死。）
- **`3s` 那枚字面从哪来**（grep 被指字符串，不抄行号）：`grep -rn "DefaultApprovalTimeout=300s" --include="*.go" .` ⇒ **唯一产码命中 `cmd/wisp/resident_approval_windows.go:456`**（另 1 枚是 `.scratch/wisp/probes/260/v1/logs/head-copy-resident_approval_windows.go:306` 的历史副本，不是产码也不是测试）。⇒ 它是**手抄进散文的第二真相源**：`300`/`3` 两个数不是从 `approval` 的常量派生的，`queue.go` 改了它不跟、也没有任何一枚钉会响。**这条才是瑕疵**（票 267 残余⑤同族＝"散文里抄常量"），⛔ 不是"数值写错"。
- **改字面会不会红**：**零枚**〔预测〕。三把尺（§5 R-14/R-15/R-16）：①该整枚字面在 `_test.go` 里 0 命中；②`grep -rn '"fallback"' --include="*_test.go" cmd internal tools` 只有 2 枚无关命中（`internal/agent/spill_path_invariant_test.go:164` 的临时目录名、`internal/llm/fallback_test.go:47` 的表键），无一处断这句 `fallback` 属性；③`grep -rn "compiled approval constants" --include="*.go" .` 只命中产码 `:454` 一句 ⇒ 没有任何测试断言过这枚 `slog.Warn` 的文本。尺的不足（诚实登记）：〔预测〕建立在"无 golden／JSONL 快照比对"之上，我没有跑包，见 §6 U2。
- ⛔ 本腿一枚字节未改。若要派生（`fmt.Sprintf("DefaultApprovalTimeout=%s / DefaultL1Window=%s", approval.DefaultApprovalTimeout, approval.DefaultL1Window)`）：只在 `cmd/wisp` 侧读常量，**不触**票 268 AC#4 的禁区 `internal/agent/approval/**`；但它**不在本票 AC 里**，属顺手改——归编排者决定是否并进 267 残余⑤或另立一枚。

## §5 尺读数与未跑清单

| # | 尺（正则／命令逐字） | 根 | 读数 | 时刻 |
|---|---|---|---|---|
| R-1 | `grep -rlE 'Stdout\|stdout\|Stderr\|stderr' --include='*_test.go' cmd/wisp \| wc -l` | 仓根 | **44** 枚文件 | 11:15:10 |
| R-2 | `grep -rn 'strings\.Count(' --include='*_test.go' cmd/wisp \| wc -l` | 同上 | **9** 枚，逐枚已列（N1–N4 覆盖其中 4 枚，余 5 枚：`firstrun_198_test.go:204/212/220` 数 TOML 文件、`panel_resident_windows_test.go:326` 数页面答案里的 `1`） | 11:15:10 |
| R-3 | `grep -rn 'bootResidentLeg' --include='*_test.go' cmd/wisp \| grep -vc 'func bootResidentLeg'` | 同上 | **18** 处调用（定义 2 处：`resident_sink_nail_127_windows_test.go:217`、`resident_task_source_246_windows_test.go:455`） | 11:15:10 |
| R-4 | `grep -rl 'go:build windows && winlive' --include='*_test.go' cmd/wisp \| wc -l` | 同上 | **6** 枚 | 11:15:10 |
| R-5b | `grep -rnE '\b[A-Za-z0-9_]*(out\|Out\|said\|console\|Console)[A-Za-z0-9_]* *(!=\|==) *(""{1,2}\|want[A-Za-z0-9_]*)' --include='*_test.go' cmd/wisp`（剔掉 `t.*f` 行） | 同上 | **1** 枚＝`cmd/wisp/firstrun_198_test.go:91` ⇒ **常驻腿 stdout 的整缓冲相等钉＝零枚** | 11:15:45 |
| R-6 | `grep -rnE 'len\(lines\)\|strings\.Split\(.*\(Stdout\|stdout)' --include='*_test.go' cmd/wisp` | 同上 | **8** 处，逐枚 Read：`config_receipt_255_test.go:203/204`（数**被引用源码文件**的行）·`:617`、`config_reload_223_test.go:220`、`resident_approval_246_windows_test.go:290`、`resident_task_source_live_246_windows_test.go:436`（都是"从尾往前找含 needle 的行"的 helper，无枚数断言）·`panel_pump_test.go:80`（数账本行）·`panel_host_gate_test.go:360-361`（数的是 **git 的 stdout**，不是 wisp 的）⇒ **打在常驻腿 stdout 上的行数钉＝零枚** | 11:25:47 |
| R-7 | `grep -oE '[a-zA-Z0-9_/.-]+\.go:[0-9]+' cmd/wisp/config_readers_255.go \| sort -u` | 同上 | **18** 枚 `x.go:N` 字面（含注释内引用），**0 枚**落在 `resident_approval_windows.go` | 11:07 |
| R-8 | `grep -rn 'resolveDataDir(' --include='*.go' cmd internal \| grep -v _test.go` | 同上 | **4** 枚产码读者（`doctor.go:105`／`models.go:110`／`providers.go:81`／`run.go:200`）＋声明 `doctor.go:247` | 11:13 |
| R-9 | `grep -n 'stdout:' cmd/wisp/firstrun_198_test.go`（读 `run198` 全体） | 同上 | run 腿走**注入 buffer**，不是 `os.Stdout` ⇒ `fmt.Printf` 不可达 | 11:14 |
| R-10 | 逐枚 `sed` 读 `cmd/wisp/config_reload.go:340-409` | 同上 | 「本次运行继续用内存里的旧配置」出现在**7** 枚收尾 | 11:15 |
| R-12 | `grep -rc '"config\.toml: ' internal/config/*.go`（剔 `_test.go`） | `internal/config` | **48**＝`validate.go 21`／`catalog.go 11`／`parse.go 8`／`migrate.go 4`／`loader.go 3`／`unwired.go 1` | 11:13 |
| R-13 | `grep -rln '"config\.toml: ' --include='*_test.go' internal cmd` | 仓根 | **0 枚**测试引用该前缀 ⇒ 无一枚钉守约定 | 11:13 |
| R-14 | `grep -rn "DefaultApprovalTimeout=300s" --include="*.go" .`（剔 `.scratch`） | 仓根 | **1**＝`cmd/wisp/resident_approval_windows.go:456` | 11:15:58 |
| R-15 | `grep -rn '"fallback"' --include="*_test.go" cmd internal tools` | 同上 | 2 枚，均与本句无关（已 Read 两处上下文） | 11:14 |
| R-16 | `grep -rn "compiled approval constants" --include="*.go" .`（剔 `.scratch`） | 同上 | **1**＝产码 `:454`，测试 0 | 11:14 |
| R-11 | `grep -rn "winlive" .github/workflows/ci.yml scripts/*.sh` | 仓根 | **0 命中** | 11:15:20 |
| R-18 | `grep -rn "ensureFirstRunConfig" --include="*.go" cmd`（剔 `_test.go`） | 仓根 | 声明 `cmd/wisp/firstrun.go:72`；产码调用者 **1 枚**＝`cmd/wisp/run.go:245`（常驻腿侧要过 `resident_windows.go:260` 才走到，晚于 `:132`） | 11:27:12 |
| R-19 | `grep -rn 'captureStdout128(' --include='*_test.go' cmd/wisp \| grep -v 'func captureStdout128'` | 同上 | **7** 处（`dataroot_128_test.go:132`／`resident_cancel_key_wording_260r3_windows_test.go:98/140/164/178`／`resident_cancel_key_label_260r4_windows_test.go:88/130`）；逐枚 Read 到窗口体＝**只包 `bindBallHost(...)` 或 `cmdDoctor()`**，构造调用都在窗口外（`260r3:96/138/161/176`、`260r4:86/128/164`） | 11:29:03 |
| R-17 | `git status --porcelain -- cmd internal tools scripts docs`；`git rev-parse --short HEAD` | 仓根 | 起手 0 行／HEAD `c6cf66e6`（11:00:46）→ 收口 2 行（` M cmd/wisp/firstrun.go`、`?? cmd/wisp/firstrun_257_test.go`）／HEAD `69c1bb82`（11:16:13） | 两处 |

**未跑的 go 命令（逐枚具名＋原因）**：`go build`／`go test`（含 `go test ./cmd/wisp/`、`-list`、`-count=2 -v`）／`go vet`／`gofumpt -l`／`gofmt -l`／`go list -deps`／`staticcheck`／`tools/d22scan` 真跑（含 `sh scripts/d22scan.sh`）／`wisp slo`／`scripts/build.ps1`／`scripts/wisp-cli-tests.sh`／`scripts/portable-tests.sh`／任何 `buildWispForTest` 派生的子进程编译。原因＝硬闸①：GitHub Actions 正在这台机器的 self-hosted runner 上跑 `slo-full`（D32 两个数的测量），叠任何 `go` 负载会污染读数。⇒ **本件全部颜色判断都是〔预测〕，归口编排者**；本件给的是"会不会红"的形状与可达性论证，不是读数。

## §6 判不动／量不到（具名，⛔ 不写"应该没问题"）

| 号 | 格子 | 能判到哪 | 为什么判不动 |
|---|---|---|---|
| U1 | 页面上有没有承接位（面板能否显示"你那行被拒了"） | Go 侧没有任何出向字段承载这枚意思（我没读到） | 硬闸②禁读 `frontend/**`／`design/**` ⇒ 归机主带去的那枚 agent |
| U2 | §1/§4 每一枚钉**今天**的真实颜色 | 只有〔预测〕＋可达性论证 | 硬闸①零 `go` 命令；且 `257-r2` 正在写 `cmd/wisp`，今天的整包读数本来就是它的产出 |
| U3 | 双击形下那句 `fmt.Printf` 用户到底看不看得见 | 静态读到 `scripts/build.ps1` 产 `-H=windowsgui` 的 exe、`cmd/wisp/console_windows.go` 无父控制台即不改绑、`resident_windows.go:55-57` 自陈"stderr going nowhere"（**stdout 同理** ⇒ 该句在双击形下也大概率零落点，与 `[hotkey]` 那枚现成先例同命） | 推理不是现量；要实机双击 ⇒ 与 `winlive` 那批同一次取证 |
| U4 | d22scan 卫生四门的现量读数（会不会因新句扩大命中） | 只能给出会撞的三条规则（emoji/`filepath.Clean|Abs`/ban #9 幻影引用）与 `emojiScopes` 覆盖 `cmd/` 的事实 | 禁跑 `scripts/d22scan.sh`／`tools/d22scan` |
| U5 | 除 `config_reload.go:155` 之外还有没有别的调用者把 reload 那七句当分类器 | 我只做了同包 grep（`describeReloadFailure` 命中＝产码 1 处＋2 处测试） | 没跑反射／没穷尽动态调用面；引这枚"唯一调用者"时请标〔静态 grep〕 |
| U6 | 落地腿开工时 `cmd/wisp/**` 的行号 | 本件行号是 `c6cf66e6`→`69c1bb82` 之间的现量 | HEAD 在 16 分钟内漂了 2 枚、scoped 脏面从 0 变 2 ⇒ 开工必重取 |
| U7 | `rescueMarkers128` 那 6 枚会不会被本票句子"顺带满足" | 不会——它断的是 `resolveDataDir`/五枚腿的输出，本票那句在另一条形上 | 判 AC#5 的实际兑现程度需要跑那枚子进程用例（硬闸①）⇒ 留给 128 的腿 |

## §7 我写错的读数（自我对抗，真改）

1. **整等型 stdout 钉我一开始数成 2 枚**。尺是窄正则 `(out|stdout|said) *!= *"`，命中 `resident_task_source_246_windows_test.go:253`、`:275`；**Read 断言体**后确认那里的 `text` 是 `residentTestTaskText(env, layout)` 的返回值（`text, refusal := residentTestTaskText(…)`），不是 stdout。⇒ 换宽尺 R-5b 重跑，正确读数＝**1 枚**（`firstrun_198_test.go:91`），且它打的是注入 buffer、`fmt.Printf` 不可达。**这条如果没 Read 到断言体就会写成"有两枚整等钉，落地腿要小心"，是假警报。**
2. **`dataroot_128_test.go:168` 的口径**：派单与 a1 都把它当"常驻腿 stdout 的 contains 型钉"引。Read 全文后：它断的是 `err.Error()`（`resolveDataDir` 的错误串），**连 stdout 都不是**；真正的合并 console contains 钉是同目录 `dataroot_128_windows_test.go:74-78`。⇒ 本件按后者立据，并保留"加行不红"的原判。
3. **`rescueMarkers128 < 6` 我一度以为是数"打印出来的自救字面枚数"**（那样新增一句就可能进它的分母）。读到 `dataroot_128_test.go:91-98` 与 `:210-211` 逐字：它是**测试文件里那份 slice 的长度**（现量 6 枚元素），断言只在这份清单被缩短时红 ⇒ **新增 stdout 不进它的分母**。判语已写进 §1.2/§1.3 之外单列，因为派单点名要这枚判断。
4. **Q4 题面我照抄会写歪**：派单/a1 把这枚字面称"与实现不一致"。现量后判**不成立**——回落支的窗口真值是 `DefaultL1Window`=3 s（`gate.go:143-146` 的 `win <= 0` 支），瑕疵是"手抄进散文、无派生、无钉"，不是"数错"。⇒ §4 按更正后的形状写。
5. **一次可达性想当然**：我先入为主以为 127/228/246 的子进程常驻用例里被拒句会出现（那样 §1 结论要重做）。逐枚读它们的 `dataDir` 来源后改判：全是空 `t.TempDir()`，且 `resident_task_source_246_windows_test.go:389-391` 反向钉死"常驻腿不许造 `config.toml`" ⇒ 被拒支在既有套件里**不可达**，这才是"零枚红"的主证，contains 型只是第二证。
6. 起手 scoped 读数 0 行我差点只记一次；11:16 复跑变 2 行（`257-r2` 进场）。⇒ §0 已改写成"两次读数＋漂移"，不再写成单点。
7. **两处读数在我自己第一版里是错的，复跑后已改**：① M5／R-7 我原写"名册引的 **17** 枚 cite"，复跑尺（同一正则、`sort -u`）读数＝**18**，且这枚 18 里含**注释里的路径引用**（`config_readers_255.go:114-116`、`:145-152`、`:187-192` 都是注释），我没有逐枚分"被断言的 cite／注释引用"——两种口径都写出来，引这枚数时不许写成"17 vs 18 是新旧版本"；② R-6 我第一版写"命中 6 处"，把尺逐字重跑（11:25:47）是 **8 处**，逐枚 Read 后结论不变（零枚行数钉），但**枚数是我的错**，已按 8 改写并把每一处的身份列出。

## §8 交件判语（三行）

1. **落地腿能不能直接开工**：**能开工（就"撞别人"这一格）**——A＋B〔预测〕**零枚既有测试会红**：既有的常驻 console 钉清一色 `Contains`／指定字面 `Count` 型（§1.1/§1.2，整等型与行数型**零枚**，尺 R-5b/R-6），更硬的一条是被拒支在既有套件里**不可达**（无一例给常驻腿带 `config.toml` 的数据根，`resident_task_source_246_windows_test.go:389-391` 还反向钉着）；⛔ 前提是它守住 §1.3 的 M1/M3/M4/M7（不加第二枚 `approval.Options`、不加 `runResident` 里的第二枚读取点、不去调 `resolveDataDir`）与 §2.2（**不搬 `describeReloadFailure` 那七句散文**，只借 `resident_windows.go:187` 的形状）。所有颜色都是〔预测〕，真读数归编排者。
2. **开工前还欠哪一格**：①**写面被占**——此刻 `257-r2` 正在 `cmd/wisp`（`M cmd/wisp/firstrun.go`＋`?? cmd/wisp/firstrun_257_test.go`，11:16:13 现量），本票面 `:38-39` 的"三枚串行"未解；②**新增的那一发用例自己会带来新的一枚 stdout 钉**（AC#2 要"文件存在但被拒 ⇒ 仍 300 s 且 provenance 具名"），本件只量了**既有**，未量落地腿自己写的那枚会不会被 `resident_hotkey_*`／`config_*_223_*` 的互斥族反咬——那要在落地腿交件前做一次"只跑 `cmd/wisp` 一包"的撞钉预检（第 64 条：要跑包、别只 grep 符号名）；③新句的措辞要**逐名避开这七枚字面**（它们是"缺席型"钉的分子，产码声明在 `cmd/wisp/resident_task_source_windows.go:101/103/106/110/113`，值逐字＝`任务入口未启用（本机没有可交互控制台）`／`任务入口已启用（控制台键盘）`／`任务入口经测试注入位打开`／`这条注入位在当前环境不被受理`／`任务管线未装配`；再加 `cmd/wisp/resident_ball_228_windows_test.go:47-48` 的 `the floating ball window is up in this process`／`this process has NO floating ball window`）——它们出现在同一枚 stdout 缓冲里就会撞 `resident_task_source_246_windows_test.go:311/324/379/382/443` 与 `resident_ball_228_windows_test.go:74` 那五枚缺席断言。
3. **必须编排者裁**：①`[panel] geometry source unreadable`（`cmd/wisp/panel_resident_windows.go:203`）是**同族第三枚折叠**（缺文件与语义拒载折成一枚、且零 stdout）——并进本票还是另立一票；②`:456` 那枚"散文里抄 300s/3s"要不要配一枚从常量派生的正向钉（票 267 残余⑤同族，不在票 268 AC 里）；③AC#5（128）与本票在 `cmd/wisp/**` 的排次序（技术上无依赖，地界互斥）；④U3——双击形下那句 `fmt.Printf` 是否真有落点，与 `winlive`／实机双击那批一起取。⛔ 票面 AC 框本件一枚未碰。
