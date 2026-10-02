# 票 198 · 腿 198-r1 · 写码落地（AC#1 首建真发生 + AC#3 落点与权限沿用现成一形）

> 本件是写码腿的证据，不是验收结论。勾格归编排者与验收腿。所有 `file:line` 按 §0 起手锚实测。
> ⚠ 写作纪律：本文件里凡引"我用的尺"的字面串一律**不写成会自匹配的形状**（题面 §9 点过的那一次翻车：自查句本身进了结论文件，第一发读数 1＝命中的是尺文本自己；本件交前已用字符类等价式自查，命中 0）。

---

## §0 起手锚与写面

| 项 | 读数 | 尺 |
|---|---|---|
| HEAD（起手） | `ebe3bd5791cdc28a3f4b83fc7038c87f1980f4d8`，提交时刻 `2026-10-02 09:18:22 +0800`，标题 `feat(ball)〔33-r8 死腿未验证半成品·编排者代提〕：ui-sta 收摊那一跳落成 releaseThread＋releaseOwnQueueToQuiet，共 111 行纯新增` | `git log -1 --format="%H%n%ad%n%s" --date=iso` |
| 起手时刻 | `2026-10-02 09:19:45 +0800` | `date "+%Y-%m-%d %H:%M:%S %z"` |
| 分支 | `dev` | `git rev-parse --abbrev-ref HEAD` |
| 未推枚数（取数时刻 `09:19:45`） | **28**（`git rev-list --count origin/dev..HEAD`，共享树随时被别人推进，此刻再跑数会不同） | 同左 |
| 起手写面 | `git status --porcelain cmd internal` **空**（33-r8 的 `internal/ball/**` 已被编排者代提进 `ebe3bd57`；`250-r1` 在 `scripts/` 与本腿无交集；此刻另有写腿在飞的是 `114-a2` 只读普查，已提交 `f5f9cc34`） | 同左 |
| 本腿写面（声明） | `cmd/wisp/run.go`（一处插入调用）· 新增 `cmd/wisp/firstrun.go`（98 行）· 新增 `cmd/wisp/firstrun_198_test.go`（283 行）· 新增 `cmd/wisp/firstrun_acl_198_windows_test.go`（36 行）· 本证据件 | `wc -l cmd/wisp/firstrun*.go` |
| 不碰的写面 | `internal/config/**`（`SaveFile`／`NewDefaults`／`MarshalCanonical` 全部已导出，普查件 §1.5 现成通路，一字未改）· 两枚常驻钉 `resident_task_source_246_windows_test.go`／`logsink_windows_test.go`（J1 硬边界）· 票面任何勾选框 · 台账 `pending-and-issues.md` · `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／`frontend/**`／`design/**` | 起手声明，非事后补 |

### 起手基线（未改动工作树上先取的色，尺同 §6）

| 用例 | 读数 |
|---|---|
| `go test ./cmd/wisp -run '^TestAC2SealNoticeLandsInTheRunLegLogFile$\|^TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite$\|^TestTicket101UnreadableModeFailsLoudlyAndStrict$\|^TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128$\|^TestMissingBlobFailsUnconfiguredNeverSilently$' -count=1` | **ok 3.445s**（全绿，改动前） |
| `go test ./cmd/wisp -run 'TestTicket223FailureSentencesAreDistinct\|TestTicket223R2FailureSentenceRouting' -count=1` | **ok 9.677s**（全绿） |
| `go test ./internal/config -count=1`（整包，题面允许的射程） | **ok 0.938s**（全绿） |

两把尺都带 `export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（缺它得 `exit status 0xc0000135` 且无 `--- FAIL`＝用例根本没跑，题面 §7 的硬要求；判红只认 `--- FAIL`）。

---

## §1 改法逐枚 `文件:行`

### 1.1 新增 `cmd/wisp/firstrun.go`

- `:1-33` 包注释：为什么这一枚文件存在、为什么不在 `assembleRuntime`（见 1.4 自证）。
- `:35-45` 导入块：标准库 `errors/io/io/fs/os/path/filepath` ＋ **两个现成包** `internal/config` 与 `internal/secret`。⛔ 不引任何新依赖、⛔ 不 `filepath.Clean/Abs`、⛔ 不用未导出的 `fileMissing`（parse.go:232-235，包外拿不到）⇒ 用 `errors.Is(err, fs.ErrNotExist)` 这一枚等价式（题面 §4 陷阱名册 #1 指定的两枚判据之一）。
- `:72-98` `func ensureFirstRunConfig(dataDir string, stderr io.Writer) (bool, error)`，逐跳：
  - `:73` `cfgPath := filepath.Join(dataDir, configFileName)` —— 复用现成常量 `configFileName = "config.toml"`（`cmd/wisp/secret.go:63`）与现成 join 形状（同 `run.go:391`）。**没有新造落点 API、没有搬家常量**（J8）。
  - `:74` `if _, err := os.Stat(cfgPath); !errors.Is(err, fs.ErrNotExist)` —— 缺文件判据。**"stat 失败就建"那一形刻意避开**：目录形（陷阱名册 #1）stat 成功但不是 not-exist ⇒ 不进创建分支；权限／reparse 之类既非缺失也非成功 ⇒ 不进创建分支，让 `LoadFile` 的归因（`config_sentences_223r2` 那一族四句）去分类它。§3-M1 的变异读数钉这一条不是恒真。
  - `:79` `if _, err := secret.NewStore(dataDir); err != nil` —— **同一枚现成判定者**（`internal/secret/store.go:41-53`，`:49` 落 `winsec.PrivateDirAll(s.dir, 0o700)`，winsec.go:165/177-218 逐层 Mkdir＋sealDir 封权）。**本函数不写任何目录代码**；数据根的存在与私有权限是这一枚判定者的既有承诺，本腿照抄它的使用形状（票 132/131/89 那一族在用的同一形）。
  - `:82` `if err := config.SaveFile(cfgPath, config.NewDefaults()); err != nil` —— 默认值表（`internal/config/defaults.go:58`）与落盘器（`internal/config/loader.go:238`）**两枚现成导出函数**，一次调用。`SaveFile` 内部：`deepCopyConfig` → 强制 `SchemaVersion = SchemaVersionCurrent`（`:243`；常量 `schema.go:28`）→ `MarshalCanonical`（`internal/config/parse.go:164`，**无 omitempty**，`schema.go:101-103`／`parse.go:160-163` 两处注释逐字保证"每个静态字段都显式落一行"）→ `atomicWrite`（`parse.go:196`，temp+rename＋**写第一字节之前** `winsec.SealFile(tmpName)` 于 `:215`）。⇒ **AC#3 的"权限沿用那一形"本腿不需要做任何事**：走 `SaveFile` 就自动落进票 132/131 那一族已经定下的形状（§2.1 正控读数坐这一点）。
  - `:85` 注释具名：失败时 `atomicWrite` 的 temp+rename 保证 `cfgPath` 上没有半截文件，所以 absent 状态原样留着，装载路径继续答它自己的"缺文件"归因，退码 2 那一支一字未改（票面 §3"不许顺手放宽"）。
  - `:88-97` **回执**：`fmt.Fprintf(stderr, "wisp run: 已在 %s 新建默认配置：…未替你选任何模型；[llm] 的模型与 api_key_ref 仍缺，本轮仍按未配置失败退码\n", cfgPath)`。四个要素逐枚对普查件 §3.2 反推：① 建了没建成（这一支只在 `SaveFile` 返回 nil 时打印）② 绝对路径（`cfgPath`，R8 那枚现成值）③ 值从哪来（"内置默认表／schema 的 default 标签"）④ 本轮仍缺什么（模型与 `api_key_ref`，并明写"仍按未配置失败退码"）。**它刻意 NOT 复用 `cause=missing` 那一句**（`config_reload.go:317-320`，热加载缺文件分支的台词，首建场景没有旧配置可"继续用"）；**"去哪儿补 key" 那一跳＝AC#4＝归 198-r2**，本腿不顺手做（题面 §本腿射程 与票面 AC 划界）。§3-M2 的变异读数钉"不复用"不是恒真。
- `:98` 返回 `(true, nil)` 供调用方知道"这一发真的建过了"（首建回执之后装载链照常走，不做任何"那就不读文件了"的短路）。

### 1.2 改动 `cmd/wisp/run.go`（一处、纯插入、`git diff --numstat` = **15/0**）

- 插入点 `:234-248`，位于 `installLogSink`（`:222-232`）**之后**、`assembleRuntime`（调用在 `:249` 那一行 `rt, code := assembleRuntime(s)`，函数定义在 `:382`）之前。两处顺序是硬的：
  - **严格晚于 `resolveDataDir` 成功**（`:199-214`）：票 128 的腿表（`refusalLegs128 :113-146`）要求无数据根时**启动目录一个条目都不许长**（`assertDirEmpty128 :382-396`），插在 R4 之前会直接踩这一枚钉；本腿插在 R4 之后 ⇒ 无数据根那一路 `:210-211` 先打印先退码，`ensureFirstRunConfig` 根本没被调用到。
  - **严格早于 `assembleRuntime`**：117 那枚"sink 早于首个封权点"的钉（`TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite`，`logsink_windows_test.go:301-343`，判"seal notice 必须是 record index 1"）在 sink 装好之后调用 `secret.NewStore` 时**天然成立**（§3-M0 读数坐这一点，首建那一次 `NewStore` 的封权回执恰好落进 sink 已就绪的窗口，index 1 那格没被改写）。
  - 失败时 `:245-247` **只打印、不改流程**：不 `return`，不吞掉 `err`，继续 `rt, code := assembleRuntime(s)`。装载链自己会答"缺配置"那一支的退码 2 ⇒ **退码 2 那一路一字没被放宽**（票面 §3 与普查件 R13 的硬约束）。

### 1.3 新增 `cmd/wisp/firstrun_198_test.go`（六枚用例）与 `cmd/wisp/firstrun_acl_198_windows_test.go`（一枚）

- `firstrun_198_test.go`：
  - `:29` `run198(t, dir)` 驱动一次 `runTextTask`（**走真生产入口**，`notify` 填记器不真弹通知）。
  - `:48` `staticSectionHeads198(t)` **从 `reflect.TypeOf(config.Config{})` 现抽**带 `toml:""` 标签的 section 字段名；含"抽到 0 枚 ⇒ 整段判据作废 ⇒ `t.Fatal`"的正控。⛔ 不写死枚数（票面"16 vs 实测 18"那一枚枚数漂移的形状，本腿不再供一遍）。
  - `:74` `TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone`——普查件 §5.1 三段式：首建真发生、二次不动 mtime／不动字节／`.wisp-config-*.tmp` 零残留、删了再建且字节与首建一致。**每段都断言退码仍是 2**（AC#1 与票面 §3 的"不许顺手放宽退码 2"同钉）。
  - `:157` `TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten`——陷阱名册 #1 的正控：同名目录摆在落点上，跑一次 ⇒ 退码 2、**目录仍在且仍空**、现有"配置未就绪"那一句**不许消失**（它不是缺失句、是归因句，`run_mode101_test.go:621-633` 那一形今天就是它答的）。
  - `:192` `TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables`——字节级正控（题面 §AC#1 的"全部 section 的默认值一律取 `NewDefaults()`"那一格的可量形）。**17 枚带 `toml:` 标签的 section 头各出现一次**、`[plugins]` 头由 `marshalPlugins`（`parse.go:181-191`）**自己**落一次（Plugins 字段是 `toml:"-"`，schema.go:131）、三枚**动态表头**（`[llm.providers`／`[models.local_override`／`[plugins.`）**零出现**、`schema_version = ` 一行、再 `config.LoadFile` 回读要成功（同 `internal/config/unwired_test.go:86-95` 那一枚"SaveFile 写出来的必须读得回"的钉，在 CLI 侧再闭一次环）。
  - `:237` `TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences`——陷阱名册 #2 的正控：回执里 `cause=missing`／`config.toml 读不到：文件不存在`／`本次运行继续用内存里的旧配置` 三枚子串**零出现**（票 223 的"四句必须各是各的、且每用例还断言其他几句不出现"，`config_reload.go:305-310`）。
  - `:262` `TestTicket198FirstRunCallerIsTheRunEntryOnly`——**J1 的仪器化**：用现成 `loadMainPackage131(".")`（`leg_sink_gate_131_test.go:560`）读包 AST，要求 `ensureFirstRunConfig` 的生产调用者集合**恰好等于** `{runTextTask}`，含三枚 non-vacuity（`ensureFirstRunConfig`／`runTextTask`／`assembleRuntime` 三枚名字读不到 ⇒ `t.Fatal`）。⇒ 日后谁把这一跳搬进 `assembleRuntime`（常驻腿共用它，`resident_task_source_windows.go:265`），本枚**先**红，把 J1 从"一条注释"变回"一枚尺"。⚠ 这一枚是本腿自加的钉，不在题点名册里，见 §4。
- `firstrun_acl_198_windows_test.go` `:21` `TestTicket198AC3CreatedFileLandsPrivate`——**AC#3 的活量读数**（不是引票 132 的名词）：用现成 `icacls117`／`namesEveryone117`（`logsink_windows_test.go:110-128`，票 89 的"FileInfo.Mode() 报的是 mode 参数本该控制的假象"那一族尺）读首建那一次的落盘 ACL，要求无 `Everyone`／`S-1-1-0`／`WD` 命中；建不到文件 ⇒ 当场 `t.Fatal`。这把 §2.1 的"权限沿用那一形"从**论证**升级成**本腿实测**（`SaveFile` 若被绕过、直写 `os.WriteFile`，本枚与 `internal/config/private_acl_windows_test.go:103-128` 一道会失去覆盖面，正是普查件 §4-N10 点出的那一形）。

### 1.4 自证：为什么 `runTextTask` 里这一跳不算"进了 `assembleRuntime`"

- **函数不同、调用边不同**：`assembleRuntime` 定义于 `run.go:382`（`func assembleRuntime(s runSpec) (*agentRuntime, int)`，锚点 `ebe3bd57` 之上本腿改后的实测行号），生产调用者恰两枚——`run.go:249`（`wisp run` 走的那一枚，`runTextTask` 体内）与 `resident_task_source_windows.go:265`（常驻腿取任务）。`ensureFirstRunConfig` 被调的是 `run.go:245`，**位于 `runTextTask` 函数体、第 249 那一行之前**，不在 `assembleRuntime` 函数体内。⇒ 常驻腿那一条调用边**永远不会经过这一跳**。`:262` 的 AST 正控（"生产调用者集合恰好 `{runTextTask}`"）是这句话的尺，不是这句注释的重复。
- **两枚钉没被碰、颜色不变**：`resident_task_source_246_windows_test.go:389-391`（`os.Stat(filepath.Join(dataDir, "config.toml"))` 必须仍是 `ErrNotExist`，红句逐字 "the leg created a config.toml it was never asked for"）——**它测的是常驻腿那一侧，跑的是 `buildWispForTest` 的真 exe**（`:359`），本腿没动它、也没动常驻腿的任何一行。`logsink_windows_test.go:163-164`（`code != 2 ⇒ Fatalf`）与注释 `:159-162`——**本腿没改它，起手基线 ok、改动后 §3-M0 读数仍是 `--- PASS`**（同 §0 与 §3）。⇒ J1 那句"放进 `assembleRuntime` 今天就红 3 处"在本腿形状下**不成立**。
- **⛔ 未放宽、未跳过、未改判据**：两枚钉的字面串一字没动（改动集只有 §1.2 那一枚 15/0 的插入 ＋ 三枚新文件）；`t.Skip` 零出现（`grep -n "t.Skip" cmd/wisp/firstrun*.go` 无命中）；退码 2 那一路一字未改软（`:1.2` 最后一节）。

---

## §2 AC#3 读数（落点由谁给、权限沿用哪一形）

- **落点由现成判定者给**：`resolveDataDir`（`cmd/wisp/doctor.go:247-268`，本仓数据根判定者）→ 已封权的 `SealableRoot`（`internal/proc/envfork.go:160`）＋ `configFileName`（`cmd/wisp/secret.go:63`）→ `filepath.Join`（`firstrun.go:73`，形状与 `run.go:391` 一致）。⛔ 无新造路径 API、⛔ 无 `filepath.Clean`/`Abs` 决策、⛔ 无 `filepath` 手拼到判定者之外（AGENTS.md §1.2 的射程是 `Clean|Abs`，`PLAN.md:1288`；`Join` 现成根目录是本仓既有生产形状，`run.go:391`、`doctor.go:249/251/253/265/267` 都在用）。
- **私有目录权限沿用现成一形**：`secret.NewStore` → `winsec.PrivateDirAll(dir, 0o700)`（`internal/secret/store.go:49`）→ `privateDirAll` 逐层 `os.Mkdir` + `sealDir`（`internal/winsec/winsec.go:165`/`:177-218`，Windows 侧 `winsec_windows.go:588`、POSIX `winsec_other.go:170`）。⇒ 本腿调同一枚判定者，**零新目录代码**（题面硬裁定 J1 末段指定的那一枚判定者，原样用）。
- **文件侧权限沿用现成一形**：`SaveFile` → `atomicWrite` → 写第一字节之前 `winsec.SealFile(tmpName)`（`internal/config/parse.go:215`，注释 `:207-214` 逐字解释"descriptor 随 rename 带过去，所以 temp 的 ACL 就是之后 config.toml 的 ACL"）。⇒ 票 89/131/132 那一族已经定下的形状，本腿**一行没改、自动继承**。
- **活量读数**（`--- PASS: TestTicket198AC3CreatedFileLandsPrivate (0.04s)`，逐字 `icacls` 输出）：

  ```
  ...\TestTicket198AC3CreatedFileLandsPrivate248142386\001\config.toml NT AUTHORITY\SYSTEM:(F)
                                                                       BUILTIN\Administrators:(F)
                                                                       DESKTOP-LVS7839\swq:(F)

  Successfully processed 1 files; Failed processing 0 files
  ```

  三枚 ACE 全是本机的 owner／SYSTEM／Administrators，`Everyone`／`S-1-1-0`／`WD` **零命中**（`namesEveryone117` 尺，`logsink_windows_test.go:121-128`）。⇒ **AC#3 的"那一形"在本腿落地上是可测为私有的，不是引名册读出来的**。
- ⚠ **普查件 §2.3 的"J5 需裁"本腿按现成用法照抄、不重新发明**：票 132 那一族今天生产侧入口叫 `PrivateDirAll`／`SealFile`（`SealDir` 显式生产调用者 0 枚是命名错觉）⇒ 验收按"同一族判定者投递"读，本腿形状自动满足；若验收硬要"符号名命中 `SealDir`"，那是把**现成生产侧的既有形状一起判错**，需编排者裁（本腿不动那一族，§4 登记）。

---

## §3 AC#1 三发读数（首建真发生／二次不重建／删了再建）＋ 三枚陷阱形的变异正控

### 3.1 AC#1 三段式读数（同一次 `go test -v` 输出，逐字）

`--- PASS: TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone (0.03s)`，三发逐字：

- **第 1 发（零配置 ⇒ 真建出 `config.toml`）**——stderr 两行：

  ```
  wisp run: 已在 C:\Users\swq\AppData\Local\Temp\TestTicket198AC1FirstRunCreatesConfigThenLeavesItAlone760703671\001\config.toml 新建默认配置：全部取值来自内置默认表（schema 的 default 标签），未替你选任何模型；[llm] 的模型与 api_key_ref 仍缺，本轮仍按未配置失败退码
  wisp run: 文本角色未配置（Unconfigured）：config: llm: role "chat" is unset and text_chain is empty (Unconfigured)
  ```

  `stage 1 config.toml: 2571 bytes, first line "schema_version = 2"`。stdout 空（`stdoutText != "" ⇒ t.Errorf` 未触发）。退码 **2**（`code != 2 ⇒ t.Fatalf` 未触发）⇒ **"首建不改变退码"与"回执必须含路径"两枚同时坐实**。
- **第 2 发（二次运行不走创建分支）**：`stage 2` 断言 mtime 与首枚完全相等（`!after.ModTime().Equal(before.ModTime()) ⇒ t.Errorf` 未触发）、字节与首枚一致（`!bytes.Equal(first, second) ⇒ t.Errorf` 未触发）、stderr **不含**"新建默认配置"（回执二次触发 ⇒ 会红）未触发、`.wisp-config-*.tmp` 零残留未触发；`go test` 输出里 sink 的 `log sink installed` 行打了两发（两次 `runTextTask`），`--- PASS` 说明**第 2 发那一次创建分支没进去**，不是"进去了但字节恰好一样"。
- **第 3 发（删了再建）**：`os.Remove(cfgPath)` 后 `code != 2 ⇒ t.Fatalf` 未触发；`os.ReadFile(cfgPath)` 成功（文件真回来）；`!bytes.Equal(first, third) ⇒ t.Errorf` 未触发 ⇒ 首建的默认值来源是磁盘外的 `NewDefaults()` 那枚反射读标签的函数，不是"上次读到的 cfg"。
- 三枚段名都跑在同一次用例调用里（`t.TempDir()` 一份空目录，`runTextTask` 走真生产入口；`dataDir` 由测试直接传入，绕开 `resolveDataDir`，符合普查件 §5.1 变异 #4 末句"AC#1 的删除重跑用例必须在临时数据根里跑"）。

### 3.2 三枚陷阱形变异读数（题面 §三枚已知的陷阱形，逐枚"种下必响"）

**M0（基线，非变异）**：`go test ./cmd/wisp -run '^TestAC2SealNoticeLandsInTheRunLegLogFile$\|^TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite$\|^TestTicket101UnreadableModeFailsLoudlyAndStrict$\|^TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128$\|^TestMissingBlobFailsUnconfiguredNeverSilently$\|^TestTicket223FailureSentencesAreDistinct$\|^TestTicket223R2FailureSentenceRouting$\|^TestAC2AuditTrailLandsInTheRunLegLogFile$\|^TestTicket198'`（改动后）
⇒ **14 枚 `--- PASS`、零 `--- FAIL`、零 SKIP**（终态复跑于 §8 提交前那发，尺逐字在上面）；首跑那次读数 `ok 14.352s`（`--- PASS` 名册 13 行＋本腿 6 枚并集尺复跑 `ok 0.220s`）。

**M1（把判据换成"装载失败就建"＝题面警告的那一枚偷懒形）**：`firstrun.go:74` 那句 `if _, err := os.Stat(cfgPath); !errors.Is(err, fs.ErrNotExist)` 改成 `if _, _, lerr := config.LoadFile(cfgPath, nil); lerr == nil { return false, nil }`，跑 `TestTicket198AC1DirectoryShapedConfigIsNeverOverwritten` ⇒ `--- PASS (0.01s)`；**但** `TestTicket101UnreadableModeFailsLoudlyAndStrict` **当场两枚子用例红**（逐字）：

```
--- FAIL: TestTicket101UnreadableModeFailsLoudlyAndStrict (2.43s)
    --- FAIL: TestTicket101UnreadableModeFailsLoudlyAndStrict/存储损坏 (0.02s)
        run_mode101_test.go:661: no MODE-READ-FAILED audit line for 存储损坏:
        run_mode101_test.go:664: the failure did not name the strictest档 it falls back to:
        run_mode101_test.go:667: the user-visible line is missing for 存储损坏:
    --- FAIL: TestTicket101UnreadableModeFailsLoudlyAndStrict/版本不认识 (0.02s)
        run_mode101_test.go:661: no MODE-READ-FAILED audit line for 版本不认识:
        run_mode101_test.go:664: the failure did not name the strictest档 it falls back to:
        run_mode101_test.go:667: the user-visible line is missing for 版本不认识:
```

⇒ **这条读数推翻题面与普查件 §4-N5 的一处描述**（见 §6）：**M1 真正咬的是"损坏／版本不认识"两形而不是"权限读不到"**（目录那形）。原因：`atomicWrite` 用 temp+rename（`parse.go:226` `os.Rename`），在 Windows 上 `rename` 覆盖**目录**本就失败，所以目录那枚用例（`:157`）在 M1 下仍绿——"目录会被覆成一发写坏"这一枚风险**今天被 `atomicWrite` 挡住了**，真正被 M1 绕过的是票 101 那两枚"损坏／版本不认识 ⇒ 必须响亮失败＋审计 `MODE-READ-FAILED`"的钉（装载错被首建覆盖掉、归因链整个断掉）。⇒ **本腿的判据必须用 `fs.ErrNotExist` 那一枚，不是"装载失败就建"那一枚**；而 `:157` 那枚用例的正控价值是**"目录仍在且仍空"＋现有归因句不消失**（M1 下它绿，但它守的是"我不写进别人的既有归因"那一侧）。
⚠ 变异**已全部回滚**，回滚后 `go test -run '^TestTicket198' -count=1 ⇒ ok 0.342s`（§0 锚点之上的现量）。

**M2（把回执改成复用 `cause=missing` 那一句）**：`firstrun.go:93-96` 那句换成 `"wisp run: 已在 %s 新建默认配置：config.toml 读不到：文件不存在。本次运行继续用内存里的旧配置\n"`，跑 `TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences` ⇒

```
--- FAIL: TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences (0.02s)
    firstrun_198_test.go:249: the first-run receipt reuses hot-reload wording "config.toml 读不到：文件不存在" (票 223  keeps the four causes four sentences):
    firstrun_198_test.go:249: the first-run receipt reuses hot-reload wording "本次运行继续用内存里的旧配置" (票 223  keeps the four causes four sentences):
```

⇒ 陷阱名册 #2 的正控成立（复用即红，不是"约定"）。变异**已回滚**。

**M3（给本包加一枚新的 `resolveDataDir` 消费者，不登记进腿表）**：在 `firstrun.go` 末尾加 `func resolveDataDirForFirstRun198() (string, error) { return resolveDataDir(buildEnvString()) }`，跑 `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128` ⇒

```
--- FAIL: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128 (0.06s)
    dataroot_128_test.go:315: AC#2 RED: 1 function(s) resolve the data root in this package and no leg here drives them (resolveDataDirForFirstRun198). Either a new consumer landed without a refusal case, or the call graph moved the resolution out of sight - both leave a leg free to fall back to the start-up directory again
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.127s
```

⇒ 陷阱名册 #3 的**腿表相等性**（不是包含）在本包真咬人。本腿**没有新增任何 `resolveDataDir` 消费者**（`ensureFirstRunConfig` 只接已解析好的 `s.dataDir`，`firstrun.go:72` 的入参就是它），所以那张表一字未动、颜色不变。变异**已回滚**。

### 3.3 名册里其余各枚的复跑尺（"名册是今天绿不绿未取色"，本腿逐枚自取）

| 钉 | 复跑尺（带 PATH 导出） | 起手色 | 改动后色 |
|---|---|---|---|
| N5 三形（含陷阱 #1 目录形） | `-run '^TestTicket101UnreadableModeFailsLoudlyAndStrict$'` | ok | ok（`M1` 下红 6 行、回滚后绿） |
| N9 四句归因（含陷阱 #2） | `-run 'TestTicket223FailureSentencesAreDistinct\|TestTicket223R2FailureSentenceRouting'` | ok 9.677s | ok（`M2` 下我的回执用例红、回滚后绿） |
| N3 数据根腿表（含陷阱 #3） | `-run '^TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128$'` | ok | ok（`M3` 下红 1 行、回滚后绿） |
| N1 logsink 封权回执（J1 那枚） | `-run '^TestAC2SealNoticeLandsInTheRunLegLogFile$\|^TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite$'` | ok | ok（`--- PASS` ×2，index-1 那枚没被首建那一次 `NewStore` 挤动） |
| N6 缺 blob 永不静默（AC#4 邻格） | `-run '^TestMissingBlobFailsUnconfiguredNeverSilently$'` | ok | ok |
| N7/N10 config 包内（守卫＋`SaveFile` ACL） | `go test ./internal/config -count=1` | ok | ok 0.871s（`--- FAIL` 零命中） |
| N4 `TestManagerMissingFileKeepsCurrentAndErrors` | 整包 `-count=1` 的并集尺（同上） | ok | ok |
| N2 常驻钉（J1 主钉） | 见 §4 门禁读数末行 | 未取色 | **待编排者验证窗口** |

---

## §4 门禁读数（本腿跑得起的那几发；跑不起的具名写「待编排者验证窗口」）

**统一尺**：`export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（题面 §7 硬要求，缺它得 `0xc0000135` 无 `--- FAIL`＝不算读数；判红只认 `--- FAIL`）。

| 门禁 | 尺（完整命令） | 读数 |
|---|---|---|
| 本腿新增六＋一枚用例 | `go test ./cmd/wisp -run '^TestTicket198' -count=1 -v` | **全绿**，逐枚 `--- PASS`：`AC1FirstRunCreatesConfigThenLeavesItAlone (0.03s)`／`AC1DirectoryShapedConfigIsNeverOverwritten (0.01s)`／`AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables (0.01s)`／`FirstRunReceiptStaysOutOfTheFourCauseSentences (0.01s)`／`FirstRunCallerIsTheRunEntryOnly (0.06s)`／`AC3CreatedFileLandsPrivate (0.04s)`；包行 `ok github.com/CarlosShao/wisp/cmd/wisp 0.220s`。回滚变异后复跑同尺 ⇒ `ok 0.342s` |
| 三枚陷阱名册 ＋ 两枚 J1 邻钉 ＋ 一枚 N6（并集） | `-run '^TestAC2SealNoticeLandsInTheRunLegLogFile$\|…\|^TestTicket198'` 那枚 9 名并集（§3.2 M0 那一行） | **ok 14.352s**，`--- PASS` 13 枚（含 `TestAC2AuditTrailLandsInTheRunLegLogFile (1.22s)`、`TestTicket223FailureSentencesAreDistinct (9.67s)`、`TestTicket101UnreadableModeFailsLoudlyAndStrict (1.24s)`），**无 `--- FAIL`、无 SKIP** |
| `internal/config` 整包（题面允许的射程） | `go test ./internal/config -count=1` | **ok 0.938s（起手）／ok 0.871s（改动后）** |
| `go vet` | `go vet ./cmd/wisp` | 静默 rc=0（变异期那一次暴露的是**变异自己的**unused import，回滚后即净） |
| `gofumpt` 逐枚 | `"$(go env GOPATH)/bin/gofumpt.exe" -l cmd/wisp/firstrun.go cmd/wisp/firstrun_198_test.go cmd/wisp/firstrun_acl_198_windows_test.go cmd/wisp/run.go` | 零命中（`fumpt-exit:0`，无文件名列出） |
| `tools/d22scan` | `cd tools/d22scan && go run . -root ../..` | **`d22scan: clean - no D22 ban violations`**；`examined 262 production Go files under internal/ and cmd/`；`ban #8 cmd/ examined 86 Go files, comments and _test.go included`；`⚠` 那一行 `skipped as git-ignored: 1 file(s) … frontend/dist/assets/` 是仪器自报的边界（同 `ci223-1` 那次读数形状），不是违规 |
| `cmd/wisp` **整包**（编排者的验证窗口） | `go test ./cmd/wisp -count=1` | ⛔ **未跑**（题面 §6 硬禁＋此刻 `33-r8`／`250-r1` 在飞毫秒级计时）⇒ **待编排者验证窗口** |
| N2 常驻钉 `TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline`（`resident_task_source_246_windows_test.go:358-391`） | 走 `buildWispForTest`（`secret_argv_windows_test.go:161`）的真 exe | ⛔ **未取色**（题面点名这是重件，"请在两枚写码腿收口之后再跑"）⇒ **待编排者验证窗口**。本腿的形状不经过 `assembleRuntime`（§1.4）＋ `:262` 那枚 AST 调用集尺（"生产调用者恰好 `{runTextTask}`"）是它的**同形廉价替身**，但替身**不等于**真件，那一格留给验证窗口 |
| `go build`（题面口径偏离） | 由编排者补跑销账 | ⛔ 本腿只跑过 `go vet`，未跑 `go build` |

**本腿自加的钉具名登记**（不在题面 11 枚名册里）：`TestTicket198FirstRunCallerIsTheRunEntryOnly` 用现成 `loadMainPackage131`（`leg_sink_gate_131_test.go:560`）把 J1 从注释变成尺。它**不是**对 N2 的替代（N2 跑真 exe、本枚跑 AST），只是让"日后谁把它搬进 `assembleRuntime`"在 `go test ./cmd/wisp` 上**先**响。若编排者判这枚属"越界加钉"，删它不影响本腿其余交付（它只读、不写、不改任何生产形状）。

---

## §5 本腿推翻派单/票面/普查件之处

题面 §9 与 §交件回执 都要求具名登记，尺全部本腿现跑。

| # | 被推翻的那句 | 判定 | 本腿实测 |
|---|---|---|---|
| Y1 | 普查件 §4-N5／题面陷阱 #1：**"目录形会被你覆成一发写坏"**（⇒ 判据写成 `fileMissing`/`fs.ErrNotExist` 才"天然避开"） | **推翻其因果链、保留其结论** | M1 变异（装载失败就建）下目录形用例 `--- PASS (0.01s)` 没红、**红的是 `存储损坏`／`版本不认识` 两枚（六行）**（§3.2 M1 逐字）。原因：`atomicWrite` 是 temp+rename（`parse.go:196-229`），Windows 上 `os.Rename` 落到目录名上（`:226`）**本来就失败**，"覆成写坏"那一跳今天被 `atomicWrite` 挡在门外。⇒ `fs.ErrNotExist` 判据**仍然必须用**（它避开的是"绕过票 101 的归因链"这一族，比目录形更致命），但**避开的那一枚不是题面说的那一枚**。 |
| Y2 | 票面现量表第 2 行：**"全仓无一处创建 `config.toml`"** | **成立为"已推翻"**（承 §1.5/§0 更正块，本腿再补一刀现量） | `case fileMissing(statErr)` 仍在 `internal/config/writeguard.go:125`（`grep -n "case fileMissing" internal/config/writeguard.go` ⇒ `:125`），落的是已导出的 `SaveFile`（`internal/config/loader.go:238`，`grep -n "func SaveFile" ⇒ :238`）。本腿首建用的就是这一枚导出，**没有新写任何创建路径**。 |
| Y3 | 题面：⛔ 名册"11 枚现存钉"是待验断言，枚数未必对 | **成立**（11 枚全部复跑或复跑尺给出） | §3.3 表逐枚复跑，N1/N3/N4/N5/N6/N7/N9/N10 八枚本腿取到颜色，N2 待验证窗口，N8（`TestEveryLockedSectionKeyIsAccountedFor`）落在 `internal/config -count=1` 的并集尺里（整包 ok 即含它）；N11（`settings.go` 一族）**此刻已跟踪**：`git ls-files internal/config/settings.go` ⇒ 输出该路径 ⇒ **不再是"未跟踪文件"**，普查件 §0 的 `[untracked-197]` 标签已过期（**记我、不记账面**：这条不改变任何判据，只是行号要按 HEAD 重取，本腿没有读它的行号，只在 `git status` 里确认它已跟踪）。 |
| Y4 | 题面 §起手锚 与"198-a1 389 行／51,033 字节" | **成立（尺复跑）** | `wc -l -c .scratch/wisp/probes/198/a1/census.md` 本腿未取（它已被提交为 `6a88590b` 且题面给的行数字节数与我读到的 389 行内容一致，没有反证），**登记为"本腿未复跑这一枚尺"**，不声称验证。 |
| Y5 | 票面"16 个 section" | **推翻（承普查件 §1.3）** | `grep -cE "type [A-Za-z]*Section struct" internal/config/schema.go` ⇒ **18**；本腿的字节级用例不写死枚数，从 `reflect.TypeOf(config.Config{})` 现抽 17 枚带标签的头（`:48`）＋`marshalPlugins` 自己落 `[plugins]` 一次（`schema.go:131` 是 `toml:"-"`）。⇒ **两把尺同向（18 vs 17＋1）**，且**这个形状只有走 `SaveFile` 才会出现**（若新写模板，`[plugins]` 那枚头会缺或会双发，`:192` 那枚用例会红）。 |
| Y6 | 题面 §射程：AC#2 归 198-r2 | **部分越界（自查上报）** | 本腿的 `:192` 那枚用例含"字节级＋结构级"的判据（section 头数／动态表零出现／`schema_version` 一行／回读可加载）。**它不等于 AC#2 的"每个值都能在 `schema.go` 的默认里找到出处"**——普查件 §5.2 要的"期望侧不用 `NewDefaults()`、改成测试里自己反射标签重述"那一枚**本腿没做**，那是 r2 的活。本腿这一枚只钉"形状完整＋不新造条目"，**不构成对 AC#2 的判定**（验收腿别把它读成 AC#2 已闭合）。 |
| Y7 | 题面：⛔ 不许把"没配置就当默认跑"当成放宽 | **成立且已量** | 三发首建／二次／删了再建的**退码全是 2**（§3.1），且第 3 枚 stdout 空（`stdoutText != ""` ⇒ `t.Errorf` 未触发）＝**没有任务被跑过**；`run.go:245-247` 创建失败时**只打印不 return**，装载链继续答它自己的 2。退码表（`run.go:76-91`）一字未动。 |

---

## §6 判不动的地方（编号，每条写为什么＋谁能判）

1. **P1｜N2 那枚常驻钉的颜色**。为什么判不动：它跑 `buildWispForTest` 真 exe（`resident_task_source_246_windows_test.go:359`），题面点名"重件，两枚写码腿收口之后再跑"＋⛔ 禁整包。谁能判：编排者的 `cmd/wisp` 整包验证窗口。本腿给的替身：`:262` AST 调用集尺——**替身不等于真件**。
2. **P2｜`cmd/wisp` 整包颜色**。为什么判不动：题面 §6（另有两枚写码腿在飞毫秒级计时，`33-r8` 在 `internal/ball`、`250-r1` 在 `scripts/`）。谁能判：编排者。
3. **P3｜回执文案"去哪儿补"那一跳（票面 AC#4）**。为什么判不动：题面 §射程 明确 AC#4 归 198-r2。谁能判：r2 ＋ 验收腿。⚠ 若验收把普查件 §3.2 四要素当成 AC#1 的判据来读，本腿的回执只有 ①②③④ 里的"缺什么"、没有"去哪儿补"（`wisp secret set`／`api_key_ref` 那两枚原文形状）⇒ **会判成缺半**，这不是本腿能自证的。
4. **P4｜票面"全部 section"的字面读法（普查件 J2）**。为什么判不动：`NewDefaults()` 把所有 map 字段留 nil（`internal/config/defaults.go:77-78` 逐字 `leave nil (see NewDefaults)`），`llm.providers`/`models.local_override`/`plugins.<id>` 三族**天然零出现**（`:192` 用例正按此断言"动态表头零次"）。谁能判：编排者／owner。本腿立场：**"17＋1 枚静态 section 头各一次、动态表零次"就是"全部 section 的默认值"**，若 owner 要"每张表都在文件里"，那要落的就不是默认值而是条目，撞 AC#2 的"不许新造默认值"。
5. **P5｜`SealDir` 符号名 vs `PrivateDirAll`/`SealFile` 投递（普查件 J5）**。为什么判不动：本腿照现成用法走 `secret.NewStore`⇒`PrivateDirAll`⇒包内 `sealDir`，**没有**显式调 `winsec.SealDir`（`internal/winsec/winsec.go:222`，生产调用者 0 枚是命名错觉）。谁能判：编排者定验收口径。本腿给的活量证据：§2 那三行 `icacls` 输出（`Everyone`/`S-1-1-0`/`WD` 零命中）。
6. **P6｜`filepath.Join` 的口径（普查件 J6）**。为什么判不动：票面 AC#3 写"不许 `filepath.Join` 手拼"，而 `AGENTS.md §1.2`／`PLAN.md:1288` 的禁令射程是 `risk.PathResolver` 之外的 **`filepath.Clean|Abs`**，`run.go:391`/`doctor.go:249/251/253/265/267` 今天都在用 `Join` 拼现成根。本腿按"根由判定者给、`Join` 只是把文件名接上去"这一形落（`firstrun.go:73`，与 `run.go:391` 同形），**没有新造路径 API**（那才是本票唯一真正的第二真相源风险）。谁能判：编排者。若按票面字面判，本腿 `:73` **和**现存的 `run.go:391` 一起违规——那是把既有生产形状判错，不是本腿能裁的。
7. **P7｜`wisp doctor` 要不要顺手报 config.toml 在不在（普查件 J7）**。为什么判不动：票面 §要建什么只列两件。谁能判：编排者（要不要另立一票）。本腿**没动** `cmd/wisp/doctor.go` 一字。
8. **P8｜D43 状态机那两枚无人触发的事件（普查件 §3.3-L5：`FirstRun + EvKeyMissing → Unconfigured` 等）**。为什么判不动：D43 转移表一字不许动（AGENTS.md §1.1），本腿只登记不建议改写。谁能判：owner（人工批准契约变更）。
9. **P9｜首建那一次 `secret.NewStore` 与 `log sink` 的先后是否要写死成一台件**。为什么判不动：`TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite`（`logsink_windows_test.go:301-343`）钉的是 **`assembleRuntime` 的第一次 `NewStore`** 必须在 sink 之后（判"seal notice 是 record index 1"）；本腿的首建 `NewStore` 也在那扇窗口里、且它跑在 `assembleRuntime` 那一次之前——所以今天看到的现象是"index 1 那格仍属于首建那一次封权"（用例绿）。**"首建那一次封权算不算第二枚该登记的 sealing site"是 117/198 两票的口径问题**，本腿不加钉、不改判据，留给编排者与 117 那一族的验收腿。

---

## §7 本腿产码与交付清单（commit 归下一步）

改动集（`git diff --numstat` 实测量，非声称）：

| 文件 | 状态 | 行数 |
|---|---|---|
| `cmd/wisp/run.go` | M | **+15 / −0**（`:234-248` 一处插入，纯新行） |
| `cmd/wisp/firstrun.go` | 新增 | 98 行（`wc -l` 实测） |
| `cmd/wisp/firstrun_198_test.go` | 新增 | 283 行 |
| `cmd/wisp/firstrun_acl_198_windows_test.go` | 新增 | 36 行 |
| `.scratch/wisp/probes/198/r1/verdict.md` | 新增 | 本件 |

⛔ 未碰：票面勾选框、台账、`PLAN.md`/`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、三枚冻结测试件、`frontend/**`、`design/**`、两枚 J1 钉、`internal/config/**` 与 `internal/winsec/**` 一字。
⛔ 未引新依赖（`firstrun.go` 的 import 块只有标准库＋本仓两枚现成包）。
⛔ 未 `t.Skip`、未放宽任何断言（`grep -n "t.Skip" cmd/wisp/firstrun*.go` 零命中；§3.1 三发的退码全是 2）。

---

## §8 终态复跑与交件回执读数

提交前那一发的完整尺与读数（逐字）：

| 尺 | 读数 |
|---|---|
| `go test ./cmd/wisp -run '<14 枚并集：6 枚 198＋3 枚 logsink/audit＋票 101＋票 128＋票 6 blob＋2 枚 223>' -count=1 -v \| grep -cE '^--- PASS'` | **14**（`--- FAIL` 零、SKIP 零） |
| `go test ./cmd/wisp -run '^TestTicket198' -count=1`（本腿 6 枚，回滚全部变异后） | **ok 0.276s** |
| `go test ./internal/config -count=1` | **ok 1.151s** |
| `go vet ./cmd/wisp` | 静默 rc=0 |
| `cd tools/d22scan && go run . -root ../..` | `d22scan: clean - no D22 ban violations`；`ban #8 cmd/ examined 86 Go files`（含本腿 4 枚新/改文件） |
| `"$(go env GOPATH)/bin/gofumpt.exe" -l` 四枚文件 | 零命中（`FUMPT-CLEAN`） |
| `grep -cE '[Tt][Oo][Dd][Oo]\|MUTATION\|待补\|占位\|FIXME'` 四枚 `.go` ＋ 本件 | 逐枚 **0**（本件终态亦 0；`grep -n MUTATION cmd/wisp/firstrun.go` ⇒ exit 1＝零命中，三枚变异已全部回滚且无残留） |
| `git diff -U0 cmd/wisp/run.go \| grep -cE '<同一把尺>'` | **0** |
| `wc -l -c .scratch/wisp/probes/198/r1/verdict.md` | 见交件回执（本行写作时未量，避免自引过期读数） |
