# 268-a1 普查：常驻（GUI）腿把"配置文件不存在"与"配置文件被值域门拒了"折成同一枚 provenance——这格今天缺的到底是什么

⚠ **本件由编排者代落盘**（2026-10-05 10:2x）。派单选错了腿型：`只读普查`会话没有写文件的工具，也没有 commit 权限面 ⇒ 腿把四问正文以回报正文交回、**盘上原本没有这件**。这是我三天内第三次犯同一条（`A590` §1、`A612` 都记过"要腿写件先确认它有写工具"），**记我，写在 §5**。
逐条标注口径：**〔编排者复跑〕**＝我用同一把尺在这台机器上现量过；**〔仅腿报〕**＝只有它的自述，我没复算。腿全程零 `go` 命令（硬闸），所以**本件不含任何"今天红／今天绿"的断言**。

## §0 起手锚

- 腿的起手：`10:16:38+0800`／HEAD `e3a8fe38`／分支 `dev`／共享工作树〔仅腿报，格式与本仓定式一致〕。
- 腿跑中 HEAD 漂到 `14735046`（＝我那笔 `ledger(A613)`，10:23 落地）——**并发可见，正常**。
- 本票爆破半径（`cmd/wisp`／`internal/config`／`internal/observe`／`internal/agent/approval`）内已跟踪 Go 文件起手零脏〔仅腿报〕；相关脏面只有 `.scratch/wisp/probes/267/r2/evidence.md`（＝在飞的 `267-r2`）〔编排者复跑：10:14 与 10:23 两次 `find -newermt` 都见到它在动〕。
- 编排者复跑时刻：**10:24–10:29+0800**，HEAD `14735046`。本件不含任何写产码改动；`cmd/wisp` 的写面此刻归 `267-r2`。

## §1 四问（照票 268 AC#0 的①②③④）

### ① `riskProvenanceUnreadable` 的全部消费者——产码侧只有**一枚**读者

| 落点 | 逐字形状 | 标注 |
|---|---|---|
| `cmd/wisp/resident_approval_windows.go:434` | 字面声明 `riskProvenanceUnreadable = "defaults (config.toml unreadable)"` | 〔编排者复跑〕 |
| 同文件 `:457` | 唯一产生点：`return 0, 0, riskProvenanceUnreadable`（在 `if err != nil \|\| c == nil` 那一支里） | 〔编排者复跑〕 |
| 同文件 `:449` | 另一支 `dataDir == ""` ⇒ `riskProvenanceNoView = "defaults (no host config view)"`（`:438`） | 〔编排者复跑〕 |
| 同文件 `:367` | 写进字段：`ra.riskWindow, ra.riskTimeout, ra.riskProvenance = window, timeout, provenance` | 〔编排者复跑〕 |
| 同文件 `:379-384` | ★**整棵产码树里唯一读者**：`slog.Info("resident gate: [risk] tier taken at construction", "provenance", provenance, …)`——注意它读的是**局部变量 `provenance`**，不是字段 | 〔编排者复跑〕 |
| `ra.riskProvenance` 字段本身 | ★**产码零读者**（面板快照／卡片／托盘／doctor 全不碰；grep `riskProvenance` 在非测试文件里只剩声明、注释与写入点） | 〔编排者复跑，尺＝`grep -rn riskProvenance cmd/wisp/*.go \| grep -v _test`〕 |
| 测试读者 | `cmd/wisp/resident_approval_risk_256_windows_test.go:125`、`:149-151`、`:165`、`:174-176` | 〔仅腿报〕 |

承载面：那枚 `slog.Info` 走 `cmd/wisp/logsink.go:157` 的 **primary＝过脱敏的 JSONL 文件**＋`:160` 的 **mirror＝不过脱敏的 stderr**〔仅腿报，与票 167 AC#6 那节 `A590` 的独立读数同形〕。

★ **加重判语的一格（我复认了机制，这条改变本票的题面重心）**：`resident_approval_windows.go:454-456` 在回落之前先 `slog.Warn("resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants", "path", cfgPath, "err", err, "fallback", "DefaultApprovalTimeout=300s / DefaultL1Window=3s")`——**原始拒载句（含 `out of range [31, 3600]`）今天已经跟着 `err` 记进日志了**。
而 `internal/observe/redact.go:107-127` 的 `Attr` 对 `key="err"` 不匹配 `secretWords`／`audioWords`／`contentWords` 任何一支，落到 `switch v.Kind()` 的 **`slog.KindAny` 且非 `[]byte` ⇒ `return slog.Attr{Key: key, Value: v}`＝原样透传**〔编排者复跑，两文件都逐字读过〕。
⇒ **本票缺的不是"记录"，是（a）一枚与"unreadable"分得开的名字，和（b）一处用户真的会读的面。** 票面 §这是什么 那句"且不说"要说满：盘上的日志**有**、但**名字是错的**（写成"unreadable"），而**用户手上没有**（见 ③）。

⚠ 顺带一枚同源瑕疵（归票 267 的 §3 那一族，不单立）：上面那句 `fallback` 字面写 `DefaultL1Window=3s`，而 L1 窗口的钳位下界是 `MinL1Window`（`gate.go:147-148` 一带，票 267 编排者裁定里现读为 2s）。**这条字面不在任何常量引用链上**＝改了 `queue.go`/`gate.go` 的常量它不会红〔编排者复跑：尺＝该字面串只在 `:456` 出现一次〕。

### ② "文件不存在"与"语义拒载"今天能不能机读分开——**半分得开，另半分只剩词面**

- **缺失那一半＝有机读形状**〔编排者复跑〕：`internal/config/parse.go:234` 逐字 `return errors.Is(err, fs.ErrNotExist)`（`fileMissing`）；`internal/config/loader.go:67` 用它分 `cause=missing`，`:74` 注释逐字承认下游"reads `fs.ErrNotExist` off it for cause=missing"，且原始 err 被当 **cause** 塞进 `observe.Wrap`（`internal/observe/errors.go:207-212` 实现 `Unwrap()`）⇒ 在 `LoadFile` 的返回值上 `errors.Is(err, fs.ErrNotExist)` **为真**。生产先例就在同目录：`cmd/wisp/config_reload.go:356` 逐字 `case errors.Is(err, fs.ErrNotExist):`。
- **拒载那一半＝今天没有正向机读形状**〔编排者复跑〕：`internal/config/validate.go:145-148`（票 267 那枚 band 的报错）逐字 `return observe.New(observe.ClassConfig, fmt.Sprintf("config.toml: risk.confirm_timeout_sec %d out of range [%d, %d]", …))`——**`Err == nil`，没有 cause**；而缺失那一支的 class **同样是 `ClassConfig`** ⇒ `ClassOf` 两侧同值＝**分类不区分**。唯一现成的分辨路径是 `cmd/wisp/config_reload.go:396` 的 `case strings.HasPrefix(d, "config.toml:")` ⇒ **这只算词面、不算类型**〔编排者复跑读到该 case 逐字在 `:396`〕。
  - 词面为什么脆（三条都具名）：① 该前缀是 `validate.go` 里**约 20 枚 `observe.New` 手写散文的约定**，无共享符号、无一枚测试钉"每条校验消息都以 `config.toml:` 起头"〔枚数属仅腿报；我复认 `validate.go` 里 `observe.New(observe.ClassConfig` 是同形批量出现，未逐枚数〕；② 判别位是**一个冒号字节**，与 `config.toml parse`／`config.toml read` 两支仅差一字；③ 拒载句由 `confirmTimeoutSecMin/Max`（`:128-129`）`fmt.Sprintf` 生成，**改词即失配**，而同为语义拒载的 `permission_mode`（`:143`）那句里**不含 `out of range`**。
- ⚠ **代价料（新建正向机读形状的价钱）**：`internal/observe/errors_test.go:16` 逐字 `if len(classes) != 17 {`、`:76` 逐字 `if len(cases) != 17 {`〔编排者复跑〕⇒ **新增一枚 `ErrorClass` 会当场打红这两枚**，且 D37（错误模型）属冻结契约面＝**人工批准**。

### ③ 常驻腿有没有"不过脱敏流水线"的用户可见面——**装了，但对常驻腿的真实起法它不是可见面**

- tee 确实装了：`cmd/wisp/resident_windows.go:66` → `cmd/wisp/logsink.go:153-161`；mirror 是 `slog.NewTextHandler(os.Stderr, …)`（`:160`），零脱敏（`:193-198` 自陈）〔仅腿报，与票 167 AC#6 那节 `A590` 三处〔编排者复跑〕同形〕。
- **双击／Explorer 起的那条腿看不见 stderr**〔编排者复跑机制层〕：生产 exe 带 `-H=windowsgui`（`scripts/build.ps1:115`，票 244 的落地），无参入口虽会调 `attachParentConsole`（`cmd/wisp/main.go:64`），但 `cmd/wisp/console_windows.go:38-42` 注释逐字「Fails harmlessly when … no parent console exists」，`rebindStdHandle` 的 `CreateFile("CONOUT$")` 一失败就直接 return（`:49-51`）⇒ `os.Stdout`/`os.Stderr` **不改绑**。工单自带的同句证据：`cmd/wisp/resident_windows.go:55-57` 逐字「double click the icon, no terminal attached, **stderr going nowhere**」〔后半句仅腿报〕。
- ★ **现成可抄的形状差（这条是本件对落地腿最有用的一格）**：`[hotkey]` 族在 `cmd/wisp/resident_windows.go:187` 有一行 `fmt.Printf("wisp: ball [hotkey]: config.toml unreadable (%v); the four hotkeys fall back to the compiled defaults (DefaultHotkeys)\n", err)`——**同一枚函数里既 `slog.Warn` 又 `fmt.Printf`**；而 `[risk]` 族那支（`resident_approval_windows.go:454-457`）**只 `slog.Warn`、零 `Printf`**〔编排者复跑，两行都逐字读到〕。⇒ "被拒要让用户看见"不必新造机制，缺的就是抄那一行。
- ⇒ 今天真能过得了的只有 primary 那一支（落 `<DataDir>\logs\*.jsonl`），**用户要开终端跑 `wisp slo` 或自己翻日志才读得到**。

### ④ `wisp doctor` 里有没有一处会读配置加载错误——**一枚都没有；而且加一条 `fail` 级会打红构建门**

- 尺照票面逐字（`|` 已转义），取数 10:26〔编排者复跑〕：`grep -c 'pass("\|fail("\|info("' cmd/wisp/doctor.go` ＝ **28**。
- ★ **三枚口径必须分开写（票面那句"28 枚检查项"是其中最粗的一枚）**：调用点 **28**〔编排者复跑〕／去重后的具名字面 **23**〔编排者复跑，尺＝`grep -o '…' | sort -u | wc -l`〕／腿报"具名检查模板 **15** 枚、实印 **14～16** 枚（`DLL colocated: ` 是动态名展开 2 枚；`deps.toml` 找到＝3 枚 pin、找不到＝1 枚 info）"〔仅腿报〕。⚠ **三个数不是三个版本，是同一棵树的三把尺**——引用时不许写成"上一版 vs 新版"。
- `cmd/wisp/doctor.go` 的 import 块（`:3-16`）里**没有 `internal/config`**，全文不调 `config.LoadFile`；唯一文件读是 `deps.toml` 与 `portable.txt`〔编排者复跑，逐行读到 import 块：`errors`/`fmt`/`os`/`os/exec`/`path/filepath`/`runtime`/`strings`/`buildinfo`/`proc`/`sherpa`/`go-toml`〕。⇒ 四问里"有没有一处会读配置加载错误"＝**没有**。
- **枚数钉＝零枚**：`find -iname "*doctor*"` 全仓只命中 `cmd/wisp/doctor.go` 本身（**没有任何 doctor 测试文件**），也无 `len(results)` 型断言〔编排者复跑〕。
- ★ **真会撞的是退码钉，不是枚数钉（这条决定落地腿怎么写）**：`scripts/build.ps1:167-169` 逐字 `& (Join-Path $outDir 'wisp.exe') doctor` ＋ `if ($LASTEXITCODE -ne 0) { Fail 'wisp doctor reported FAIL.' }`〔编排者复跑读到该两行〕；而 `doctor.go:34-36` 的 `fail()` 置 `critical: true`、`:120-127` 任一 critical FAIL ⇒ `cmdDoctor()` 返 false ⇒ `main.go:96-98` `os.Exit(1)`。⇒ **加一条 `fail` 级检查项，会把任何"配置里有带内越界值"的机器的构建门打红**；要走 `info()`（`:37-39`，`critical:false`）〔行号仅腿报，机制层我复认〕。
- 其余会擦到的名册钉（形状代价料）：`cmd/wisp/resident_approval_risk_256_windows_test.go:560` 逐字 `if fed != 1 {`、`:563` 逐字 `if bare != 0 {`〔编排者复跑读到两行逐字与红句〕⇒ **不许在 `runResident` 里加第二枚配置读取点**；`cmd/wisp/leg_dispatch_gate_133_test.go:149` `minLegs133 = 11`＋`:187`＋`censusVsUsage133`（`:1707`，调用于 `:218`）⇒ **只有新增子命令**才撞，需在 `main.go` usage 块补一行〔仅腿报〕；`dataroot_128_test.go:375`（`< 4`）与 `:210`（`rescueMarkers128 < 6`）与 doctor 检查项数无关〔仅腿报〕。
- ⚠ 与票 128 AC#5 的关系（我裁，不属本件）：那条"常驻（GUI）腿补自救文案"至今未派，与本票 AC#1 的"要有一处用户真的会读的面"是**同一寸地面**；落地腿不许造出第二套词（`A590` 的"不许硬造第二套"规矩）。

## §2 候选形状（⛔ 只摆料不裁形，裁定归编排者；裁完落 `A##`＋撤销口令）

| 形 | 动的文件 | 新增用户可见文案 | 撞钉（本件已复跑的） |
|---|---|---|---|
| **A**＝换一枚具名 provenance＋用 `errors.Is(err, fs.ErrNotExist)` 分支 | `cmd/wisp/resident_approval_windows.go`＋同包 `resident_approval_risk_256_windows_test.go` | 是（一枚新字面） | 要保 `256_test:141-144` 绿（票 268 AC#2 明写）；doctor 与 build.ps1 不碰。**到达面只有 JSONL** |
| **B**＝A＋一行 `fmt.Printf`（抄 `resident_windows.go:187` 那枚现成形状） | 同上 | 是 | 多一枚**双击下仍看不见**的面；从终端起时可见且可测（先例 `resident_hotkey_live_258_windows_test.go:69`〔仅腿报〕） |
| **C**＝A＋doctor 一条检查 | `cmd/wisp/doctor.go`（该文件**首次** import `internal/config`；包级依赖边早已存在，不算新增边） | 是 | ★`build.ps1:168-169` 退码钉 ⇒ **必须 `info()` 级**，否则构建门变红 |
| **C′**＝新增 `ErrorClass` 或导出哨兵，做正向机读 | `internal/observe/errors.go` | 否 | ★撞 `errors_test.go:16`／`:77` 两枚 17 枚钉，且 D37＝**人工批准面** |
| **D**＝面板回执里加一句 | `frontend/**` | — | ⛔ **量不到**（U1）；票 167/268 都写了"要加承接位就停手上报" |

**可复用的现成料（不新造词）**：`cmd/wisp/config_reload.go:354-409` 的 `describeReloadFailure` 已在**同包生产代码**里把配置错误分成 `missing/permission/syntax/unknown-key/migration/invalid/unclassified`，包内可直接调用，⛔ 不必自造第四枚字符串——但它的 `invalid` 支本身就依赖 §1② 那枚词面前缀〔仅腿报；我复认 `:356` 与 `:396` 两个 case 确在，未逐支读全函数〕。

## §3 尺读数名册（编排者复跑的 8 把，逐把带时刻；余下归〔仅腿报〕）

| # | 尺 | 根 | 读数 | 时刻 | 标注 |
|---|---|---|---|---|---|
| C1 | `grep -rn riskProvenance cmd/wisp/*.go \| grep -v _test` | `cmd/wisp` | 12 行＝声明 2／注释 4／写入 1／产生 2／**读 1**（`:380` 用的是局部变量） | 10:25 | 复跑 |
| C2 | `sed -n '446,462p'`（回落那一支） | 同上 | `slog.Warn(… "err", err, "fallback", "DefaultApprovalTimeout=300s / DefaultL1Window=3s")` 逐字在 `:454-456` | 10:25 | 复跑 |
| C3 | `internal/observe/redact.go` `Attr` 函数体 | `internal/observe` | `KindAny` 非 `[]byte` ⇒ 原样透传（`:123-126`） | 10:27 | 复跑 |
| C4 | `fs.ErrNotExist` | `cmd/wisp`＋`internal/config` | 3 命中：`config_reload.go:356`／`parse.go:234`／`loader.go:74`（注释） | 10:26 | 复跑 |
| C5 | `len(classes) != \| len(cases) != ` | `internal/observe` | 2 枚（`:16`／`:76`），want 都是 **17** | 10:26 | 复跑 |
| C6 | `grep -c 'pass("\|fail("\|info("' doctor.go` 与去重版 | `cmd/wisp` | 调用点 **28**／去重具名字面 **23** | 10:26 | 复跑 |
| C7 | `find -iname "*doctor*"` | 全仓（剔 `.git`） | **1 枚**＝`cmd/wisp/doctor.go` ⇒ doctor 无测试文件 | 10:26 | 复跑 |
| C8 | `build.ps1` smoke 那两行 | `scripts` | `wisp.exe doctor`＋`if ($LASTEXITCODE -ne 0) { Fail }` | 10:25 | 复跑 |
| — | 腿的 R 系列（doctor 15 枚模板／`fed!=1` 名册／`minLegs133`／`censusVsUsage133` 等） | — | 见 §1 各处内联标注 | 10:19–10:22 | 〔仅腿报〕 |

**未跑清单（本件一票未跑，全数具名）**：`go build`／`go test`／`go vet`／`gofumpt`／`wisp slo`／任何 `scripts/*.sh`／`doctor` 的**真跑**（我只读了它的源码与 build.ps1 里那一行）。原因＝硬闸①：`267-r2` 正在取 `cmd/wisp` 整包红名册与突变读数。**⇒ 本件不含任何"今天红/绿"的断言**，谁红谁绿只有编排者那一发可答。

## §4 量不到／判不动（具名）

| 号 | 格子 | 能判到哪 | 为什么判不动 |
|---|---|---|---|
| U1 | 页面上有没有承接位（面板回执能不能显示"你那行被拒了"） | Go 侧零枚出向字段承载这枚意思 | 硬闸②禁读 `frontend/**`／`design/**` ⇒ **需机主带去他用的那枚 agent 另查** |
| U2 | 各枚既有钉子**今天**的颜色 | — | 硬闸①零 `go` 命令 |
| U3 | 实机双击下 stderr 是否真零落点 | 静态读到"无父控制台则不改绑"（`console_windows.go:38-51`） | 推理不是现量；要真机双击才验得了 ⇒ 与 `winlive` 那一格同批 |
| U4 | 加一行 stdout `Printf` 会不会打红"常驻腿 stdout 只有这几行"型断言 | 只确证 `dataroot_128_test.go:168` 是 `contains` 型（加行不会红） | 未穷尽全树 stdout 形状断言 ⇒ 落地腿开工前必做一次撞钉预检（第 64 条：要跑包、别只 grep 符号名） |
| U5 | `validate.go` 里"约 20 枚 `observe.New`"的确切枚数 | 只确证它是同形批量 | 我没逐枚数；引这枚数时标〔仅腿报〕 |

## §5 口径与纪律（自我对抗）

- ★**代落这件事本身**：派单要求"写件并逐笔 commit"，而我选的是没有写工具的只读腿型 ⇒ 件差点不存在。**第三次**（`A590` §1、`A612`、本次）。定式补一条：**派单里凡要求盘上交付件，腿型必须是 general-purpose；只读腿只许要求"回报正文"，且我须在派单里写"你无需落件"**。这次我把它的正文**全文落盘并逐把尺复跑**，代价由我承担，不改判语。
- **我没有替腿圆任何一处**：它 §1① 那句"缺的不是记录，是一枚名字和一处会读的面"我复跑机制层（C2＋C3）之后**认下并加严**——票面"且不说"要改读成"盘上有、名字错、你手上没有"。
- **行号一律现量**：`A611`/`A612` 记过两次行号漂移，本件每条 `file:line` 都是我这轮 `sed`/`grep` 读到的位置；引用本件前仍要先 grep 被指字符（票 267 定式）。
- ⚠ **§1④ 那三枚口径（28／23／15）**是本次唯一一处我自己差点写成"腿报 15、我复跑 28 ⇒ 腿数错了"的地方——**不是数错，是尺不同**（调用点／去重字面／具名模板）。谁引用谁负责带口径。

## §6 交件判语（三行）

1. **这格今天缺什么**：不缺日志（拒载原句已带 `err` 落进 JSONL，机制层我复认）；缺的是**一枚与 `unreadable` 分得开的具名状态**＋**一处双击用户真会读到的面**。
2. **最低达标档**：形 A 单独**不达标**——它的到达面只有 JSONL，而票 268 AC#1 要的是"能到达用户的眼睛至少一次"。按 §1③ 的现成形状，**A＋B（抄 `resident_windows.go:187` 那一行）是Go 侧能自足的最小组合**；C（doctor）要额外背 `build.ps1` 的退码代价 ⇒ 只能 `info()` 级。
3. **必须谁裁**：② 要不要新建正向机读形状（新 class 撞 17 枚钉＋D37＝**机主人批**）；④ doctor 那一条走 `info` 还是接受构建门变红（**归我**，技术取舍）；D 支要不要动前端（**归机主带去的那枚 agent**）；U3 的实机双击取证（**等 `winlive` 那个词**）。⛔ 票面五枚 AC 框本件一枚未碰。
