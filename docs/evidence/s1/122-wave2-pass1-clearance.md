# 122-wave2-pass1 — 清账波 2：pass1（GOOS=linux, -checks=all,-U1000）31 条 → 3 条

- 日期：2026-11-13（+08）
- 腿：写腿（派单 `.scratch/wisp/probes/orch/20261113/txt/122-wave2-pass1-dispatch.txt`）
- 开工 HEAD：`da2702a2a`（dev）
- 工具版本（票 122 AC#7(b)/(c) 要求随读数同报，与批 1/2 记录一致未变）：
  `staticcheck.exe 2026.2.2 (0.8.2)` · `gofumpt v0.12.0 (go1.27.1)` · `go version go1.27.1 windows/amd64`

## 1. 起手复测（正向判据的"前"表）

复跑 `GOOS=linux CGO_ENABLED=0 staticcheck -checks=all,-U1000 ./...`＝**31 条**（与派单逐条一致）。
派单把 ST1021 写作 ×10，实数 **×9**（对总数 31 无影响）；goos-scope 腿（13:1x log）记的 32 条含
SA9003×1，本波起手已无该条（期间提交动过，非本波，不回改其记录），以本波现跑为准。

前表：

| 码 | 前 |
|---|---|
| SA5000 | 1 |
| SA9009 | 1 |
| ST1000 | 2 |
| ST1003 | 10 |
| ST1020 | 3 |
| ST1021 | 9 |
| ST1022 | 2 |
| ST1023 | 2 |
| **合计** | **31** |

## 2. 清后（"后"表）

| 码 | 前 | 后 |
|---|---|---|
| SA5000 | 1 | 1（特殊项②，未落，见 §5） |
| SA9009 | 1 | 1（特殊项①，未落，见 §5） |
| ST1000 | 2 | **0** |
| ST1003 | 10 | **0** |
| ST1020 | 3 | **0** |
| ST1021 | 9 | **0** |
| ST1022 | 2 | **1**（`internal/plugin/manifest.go:11`＝特殊项③，未落，见 §5） |
| ST1023 | 2 | **0** |
| **合计** | **31** | **3** |

**pass1 剩的 3 条＝派单第二类的 3 条，一件不多、一件不少；第一类 28 条全数真改代码/注释清除，
零 `//lint:ignore`、零 staticcheck.conf、零 checker 开关。**

对账（pass2＝GOOS=windows 全集，静态同一棵树）：pass2 非 U1000 由 45 → **17**
（＝45 − 本波跨平台清除 28；剩余 17 全部是派单清单之外的 windows-only 面 14 条
＋特殊项 3 条，归下一波工单面，本波零越界）。本波后 pass2 全集读数 123→95。

## 3. 28 条逐条（file:line / 原 / 改成了什么）

### ST1021（类型 doc 首行，9 条）
1. `internal/agent/compress.go:42` — 原：CompressionState 的 doc 块首行是 "CompressionReport describes
   one compression pass..."（:42–60 的 CompressionReport 说明与 :61–68 的 CompressionState 说明
   连成一块挂在 :69 的 `type CompressionState` 头上）。改：**拆块并各自归位**——CompressionState 的
   doc（"CompressionState is the outcome of one compression pass…"）留在类型头；
   CompressionReport 那段**原文搬运**到 `:161` 的 `type CompressionReport struct` 声明处
   （它原先反而没有 doc）。原文信息零丢失；其中一句 "emitted right here" 按新位置如实改为
   "emitted by Compress, on the success edge"（搬运后 "here" 不再指类型声明处，防注释撒谎）。
2. `internal/llm/content.go:19` — 原 "// Role of a normalized Message." → "// MessageRole is the
   role of a normalized Message."
3. `internal/llm/resolver.go:36` — 原 "// Role carries the call-time sampling settings resolved
   with a role." → "// RoleSettings carries the call-time sampling settings resolved with a role."
4. `internal/present/dispatch.go:13` — 原：Speaker 的 doc 块首行是通道端口总览（"The channel
   ports. Each one is the narrowest face..."）。改：Speaker 自己的一句 doc
   （"Speaker hands the line to the voice leg's announce queue (D38b owner goroutine). The ChTTS
   exfil gate lives on that announce path in the host (ticket 26 wiring) and is NOT duplicated
   here."）挂在类型头上；总览段原文保留在块后续（Speaker 列表项加冒号与后文对齐，文字未改）。
5. `internal/session/scope.go:60` — 原 "// The five engine families whose lifetime the scope owns..."
   → "// EngineKind names one of the five engine families whose lifetime the scope owns..."
6. `internal/statemachine/timeouts.go:18` — 原 "// Per-state timeout table (SPEC-08 §3 timeout
   column):" → "// TimeoutSpec is one entry of the per-state timeout table (SPEC-08 §3 timeout
   column):"（表体与其余原文逐字保留）。
7. `internal/tools/capability.go:8` — 原 "// C3 capability face (SPEC-07 §2). The eleven tokens
   below are the frozen set:" → "// Capability is one token of the C3 capability face (SPEC-07 §2).
   The eleven tokens below are the frozen set:"。
8. `internal/tools/refusal.go:7` — 原 "// The machine-readable refusal classes (ticket 20 frame 7
   (a))." → "// RefusalKind is the machine-readable refusal class (ticket 20 frame 7 (a))."
9. `internal/tools/registry.go:13` — 原 "// C4 ToolProvider (SPEC-07 §2): four slots - Builtin,
   Manifest (Tier-1), Goja (Tier-2), MCP. Only Builtin has an implementation here." →
   "// Kind is one C4 ToolProvider slot (SPEC-07 §2): four slots - Builtin, Manifest (Tier-1),
   Goja (Tier-2), MCP. Only Builtin has an implementation here."

### ST1020（函数/方法 doc 首行，3 条）
10. `internal/audio/bargein.go:534` — 原 "// Windows reports the per-window untrusted flags..."
    → "// UntrustedWindows reports the per-window untrusted flags..."。
11. `internal/risk/syncdirs.go:454` — 原 "// Roots returns the detected roots (panel/debug view)..."
    → "// SyncRoots returns the detected roots (panel/debug view)..."。
12. `internal/tools/system.go:395` — 原：十一连排 Decl 工厂的总注释 "// Decls for the family. Each
    is a function rather than a package variable..." 挂在 SystemGetDecl 上。改：首句改写为
    "// SystemGetDecl is the system.get row of the system family."，"为何是函数而非包变量"的原文
    保留在后续行（该理由本就适用于整排）。

### ST1022（常量 doc 首行，2 条）
13. `internal/audio/audio.go:41` — 原 "// Target sample rate of the C8 seam: everything downstream
    (VAD/ASR/KWS) consumes 16kHz mono (SPEC-04 sec 3)." → "// TargetRate is the target sample rate
    of the C8 seam: ..."（其余原文逐字）。
14. `internal/plugin/manifest.go:11` — **不动**（特殊项③，见 §5）。前表中 ST1022×2 的另一条
    就是它。

### ST1000（包注释，2 条）
15. `internal/agent/images.go:1` — 原文件无包 doc（首行直接是票 23 专题注释）。改：新首段
    "// Package agent - file-level note: images.go is the memory half of the C7 image seam
    (ticket 23 AC#2). The package's own charter - the ReAct loop, budgeting, cancellation - is
    documented in doc.go; this file adds the process-local image store the loop hands screenshot
    bytes to."（说真话：包主章在 doc.go，本文件只是图像仓库半边，没有冒充包总述）；票 23 原注释
    原文保留。
16. `internal/winsec/pinnedleaf_342_unix.go:3`（`//go:build unix` 之后的首注释块）— 改：新首段
    "// Package winsec - file-level note: this is the unix family's placement-pin file (tickets
    342 + 348, form 乙): the ancestor pin and the two act-time primitives (openat, unlinkat) the
    private write route stands on. The package-wide promise - "the current user is the only
    principal that can read this" - is documented in winsec.go."（全包承诺引 winsec.go 原句，不硬拼）；
    342/348 原注释原文保留。

### ST1023（冗余类型标注，2 条）
17. `internal/agent/loop_golden_test.go:222` — `var cfgTokenBudget int = 300` → `cfgTokenBudget := 300`
    （零语义；注释原样）。
18. `internal/llm/llm_test.go:28` — `var c Content = p` → `c := p`（零语义）。原行尾注释
    "compile-time: every part satisfies the interface" 在新写法下会变成空话（`p` 本就是 `Content`
    接口值，真正的编译期封印是 :23 的四常量字面量）→ 改为真话 "p is already a Content: the literal
    above is the compile-time seal"。

### ST1003（命名，10 条；改声明＋全部引用点）
19. `internal/audio/aecspike/longrun.go:180` `dropDb`→`dropDB`（字段；引用点 :171 注释、:213、:231）
20. `internal/audio/aecspike/longrun.go:181` `lowAbsDb`→`lowAbsDB`（字段；引用点 :173 注释、:214、:232）
21. `internal/audio/aecspike/longrun.go:203` `collapseDropDb`→`collapseDropDB`（包内常量；引用点 :213）
22. `internal/audio/aecspike/longrun.go:204` `collapseLowAbsDb`→`collapseLowAbsDB`（包内常量；引用点 :214）
23. `internal/audio/bargein.go:64` `BargeinCollapseDropDb`→`BargeinCollapseDropDB`（导出常量；引用点 :208）
24. `internal/audio/bargein.go:65` `BargeinCollapseLowAbsDb`→`BargeinCollapseLowAbsDB`（导出常量；引用点 :209）
25. `internal/audio/bargein.go:197` `dropDb`→`dropDB`（字段；引用点 :190 注释、:208、:226）
26. `internal/audio/bargein.go:198` `lowAbsDb`→`lowAbsDB`（字段；引用点 :192 注释、:209、:227）
27. `internal/ball/tokens.go:447` `wensityStopTs`→`wensityStopTS`（包内 var；引用点 tokens.go:443
    注释、wensity.go:41 注释、wensity.go:54、jarvis.go:25 注释——共 4 文件）
28. `internal/memory/undo.go:49` `lastTs`→`lastTS`（函数内局部 var；引用点 :51 ×2）

### 改名前后 `git grep` 全仓计数（判据：旧名计数归 0）

| 旧名 | 前（全仓命中行） | 后（除 `docs/evidence` 外） |
|---|---|---|
| dropDb | 12 行/3 文件 | **0** |
| lowAbsDb | 12 行/3 文件 | **0** |
| BargeinCollapseDropDb | 4 行/2 文件 | **0** |
| BargeinCollapseLowAbsDb | 4 行/2 文件 | **0** |
| collapseDropDb | 4 行/2 文件 | **0** |
| collapseLowAbsDb | 4 行/2 文件 | **0** |
| wensityStopTs | 10 行/6 文件 | **0** |
| lastTs | 4 行/2 文件 | **0** |

前计数中仅有的非代码命中＝`docs/evidence/s1/122-staticcheck-goos-scope.md`（逐条 finding 引文，
append-only 证据件不改）与 `docs/evidence/s1/65-glass-impl.md`/`65-verdict.md`（票 65 历史档案，
非名册、非可执行面，按"临时件与历史档不改写"惯例保留旧名）。改名前专门核过静态名册类测试：
`*.sh`（portable-tests GUARD C/D 名册、票 171 归因名册、tokens_table_test.go、
cmd/wisp/ball_style_65_windows_test.go、internal/panel/tokens_fourway_test.go）对这 8 个旧名
**零引用**（`git grep` 逐名验证，wensity 停位由 `TestWensityBandGeometryPinsTheQuantifiedSpec`
经 `WensityBandTokens()` 读值，不读符号名）。

## 4. 门禁读数（如实，含 SKIP）

| 门禁 | 读数 |
|---|---|
| `go build ./...` | rc=1，**唯一**红＝`build/mut_rvw.go` undefined voiceTtsFiles 等——该文件被 `.gitignore:14` 忽略（`git check-ignore -v` 实证），是本机变异排练残件（派单明令禁删），不进 commit、CI 不见它；我的 19 个改动文件全部编译干净（scoped build rc=0） |
| `go vet`（11 包） | rc=0 |
| `go test -count=1`（10 门禁包） | audio/ball/llm/memory/present/risk/session/statemachine/tools/winsec 全 `ok`；`internal/audio/aecspike` 首跑 `exit status 0xc0000135`＝本机环境（测试二进制静态链 sherpa cgo DLL、DLL 在仓库根、`go test` 的 cwd 在包目录解析不到；`longrun_windows_test.go:35`/`rig_windows_test.go:13` 自带"PATH 前置仓库根"的运行约定）。把仓库根加入 PATH 复跑＝**ok 0.169s**（PASS 19 / SKIP 13 / FAIL 0） |
| aecspike 13 个 SKIP 点名（真窗/DMO 硬件夹具，`WISP_AEC_SPIKE=1` 门） | TestSpikeAec3ErleRemeasure、Test388AnnounceAttenuationTripleReadings、Test388AttenuatedSessionReachesTheEndpointLoopback、TestSpike389r1BargeinReadings、TestAC297FormOneRenderReachProbe、TestSpikeAec3LongWindow、TestSpikeLongWindowProbe、TestSpikeLongWindowAutocorr、TestSpikeLongWindowDelayTrack、TestSpikeEndpointsAndDmoInventory、TestSpikeSelfEchoPathRig、TestSpikeFalseTriggerTable、TestSpikeDmoFullDuplexAec——13 个全部是环境夹具门，零断言红 |
| `internal/agent`（ST1023 改动所在，非门禁包） | `go test -count=1` ok 5.844s |
| `gofumpt -l internal cmd tools scripts` | 唯一命中 `internal/models/acl_sid_121_test.go`＝**非本波改动**（git status 里它无 diff，工作副本被其他进程动成 CRLF；HEAD 提交面 blob 过 gofumpt 干净——已用 `git show HEAD:<file>` 直接喂 gofumpt 复证 rc=0）；本波 16 个改动文件 0 命中 |
| `sh scripts/d22scan.sh` | clean，rc=0（481 个生产 Go 文件 + 全部 ban #6/7/8/10 scope） |
| `sh scripts/portable-tests.sh --self-test` | **脚本没有 --self-test 模式**（该旗被透传给 `go list` 报 "flag provided but not defined: -self-test"，rc=1）；按脚本 Usage 改跑两把真尺：`--scope=windows` rc=0（strict runner `runtests.sh: OK - top-level: PASS=833 FAIL=0 SKIP=0, === RUN=1275, '[no tests to run]'=0`，12 条账本条目 8 条在本平台记账，GUARD A–D 全过）；核心（ubuntu CI）scope 在本机 rc=1＝aecspike 13 个 `WISP_AEC_SPIKE=1` 硬件夹具 SKIP 未在 ubuntu 账本上（本机形状，与改动无关；该 scope 本来就是 ubuntu CI 的形状） |

## 5. 第二类 3 条：判断＋方案（**未落**，按派单只交方案）

### ① `frontend/embed.go:4` SA9009 — 不改，建议具名走 staticcheck.conf
- 判断：与台账两次登记一致（票 122 AC#3：`frontend/embed.go:4` 是散文注释、恰以 `// go:embed`
  开头；真指令在 `:19` 写法正确；`frontend/**` 归机主外部 agent，本编队零写面）。
- 能不能根目录 `staticcheck.conf` 具名抑制这一条？**能，且是仪器层面的正确位置**：SA9009 报的
  是"指令写法"，而那条注释是**刻意的散文**（票 77 AC#1 的文档内容），改注释＝拆文档，改指令＝
  撒谎。台账 :1489 明令的是"不许 `//nolint` 全局抑制"；`staticcheck.conf` 把**这一条**从仪器
  分母里具名移出，配文档说明为什么，与"不许为了变绿放宽断言"不冲突（它不碰任何测试断言，
  只是把已知假阳性从仪器读数里剔出去，性质同票 122 AC#3 的"只登记不修"，只是换了一个不会
  每波重报的登记处）。⚠ 契约注意：`.github/**` 是 85a/编排者地界——**staticcheck.conf 落在
  根目录不算 .github，但 lint 步若要显式传 `-config` 则要动 ci.yml**；好在 staticcheck 默认
  自动读取根目录 `staticcheck.conf`，**ci.yml 一字不用动**。
- 建议落法（内容如下，**本波未落**）：
  ```ini
  # staticcheck.conf — root module
  # checks = ["all", "-U1000"] 的口径由 CI 的两次显式 -checks 旗决定，这里不动分母；
  # 本文件只登记一个已知假阳性：
  #   SA9009 @ frontend/embed.go:4 — 该行是散文注释、恰以 "// go:embed" 开头；
  #   真指令在 :19 且写法正确（票 122 AC#3、台账 A1100 族登记）。改它会撞票 77 AC#1 的文档。
  # frontend/** 归外部 agent，本文件是唯一不碰 frontend 字节就能让仪器停报的手段。
  checks = ["all"]
  ```
  ⚠ 代价写清楚：`checks = ["all"]` 在 conf 里配好后，CI 两遍的 `-checks=all,-U1000` 与
  `-checks=all` **命令行旗优先于 conf**，分母不变；conf 的 `checks` 键只在不带旗的本地裸跑时生效。
  若只想去掉 SA9009 而完全不改默认口径，更窄的写法是 `"*, -SA9009"`——但那会把裸跑口径也收窄，
  两遍 CI 旗照旧。**推荐用注释＋默认 `checks = ["all"]`（甚至只写注释不写 checks 键）**，
  让 SA9009 的排除由 CI 步骤自己的 `-checks` 旗表达（pass1＝`all,-U1000` 再加 `-SA9009`，
  pass2＝`all` 再加 `-SA9009`）——那要动 ci.yml 两行，归编排者。
- 三选一留给编排者：**甲** ci.yml 两遍 -checks 各加 `-SA9009`（最窄、具名、动 .github）；
  **乙** 根目录 staticcheck.conf 只写登记注释＋`checks` 不动（纯文档登记，仪器照报）；
  **丙** 保持现状（每波读数都带这 1 条，等外部 agent 清 frontend）。本腿推荐**甲**（它把
  "已知假阳性"变成仪器不再重报的具名事实，且 frontend 字节零触碰）。
- 撤销口令：从 ci.yml 两遍 -checks 里删掉 `-SA9009`（或删 staticcheck.conf）即复原。

### ② `internal/plugin/panic_isolation_51_test.go:62` SA5000 — 是"有意写 nil map 的夹具"，该具名抑制，不是 bug
- 判断依据（读了整文件 359 行）：
  - `:40-43` crashHost 的 doc 明写 "panicVal == nil selects the nil-map-write form: the classic
    nil-dereference host bug"——写 nil map 是**被测试的输入**，不是夹具写错了；
  - `:145-148` 第二个用例的 doc 明写 "the same claim for the other panic shape a host bug
    produces: a runtime.Error (nil map write / nil dereference) rather than an explicit panic"——
    两个用例刻意覆盖两种崩溃形状（显式 panic vs runtime.Error）；
  - 该文件要验的正是票 51 AC#4 "JS 触发宿主 Go 实现崩溃 → recover、邻插件无恙、宿主存活"；
    而 nil-map 写恰好是 staticcheck 唯一能在编译期指认的"真崩溃形态"，换成 `panic(...)` 反而
    丢掉 "runtime.Error 形状也要被 recover" 这半个断言。
- 结论：**夹具有意、断言有效、不该修**（修＝拆半个断言）。台账 :35687 "是夹具本身、有意不修"
  的登记读对了。它每波重报是仪器口径的代价，处置与特殊项①同族：在 pass1/pass2 的 `-checks`
  里具名排除这一条做不到（staticcheck 的 `-checks` 是按**码**不是按**位置**）；能用的具名抑制
  只有行级 `//lint:ignore SA5000`（在 :61 那行上一行加 `//lint:ignore SA5000 reason`）——
  这**违反票 122 禁令**（"反向：⛔ 靠加 //lint:ignore 交差"），但禁令的原文针对的是"第一类
  28 条"，对第二类的 3 条派单本就允许"讨论抑制"。**建议**：编排者若裁定可以落，就在 :61 行
  加 `//lint:ignore SA5000 ticket 51 AC#4 fixture: the nil map write IS the tested host-bug
  shape`，一行、具名、理由内联；`internal/plugin` 归票 70 腿，落不落由那条腿/编排者定。
- 撤销口令：删掉那一行 `//lint:ignore` 即复原。

### ③ `internal/plugin/manifest.go:11` ST1022 — 不碰，留给票 70 腿
- 判断：`internal/plugin/**` 现由票 70 分诊腿独占（派单明令）。修法是平凡的一行
  （"// CurrentSchemaVersion is the only schema_version this loader understands (SPEC-07 §4
  example: schema_version = 1)." 已是合格首句，只差把 "CurrentSchemaVersion is" 放到首行最前——
  现文首行是 "C14 plugin manifest — Tier-1 declarative base..."，块首行需换成以
  "CurrentSchemaVersion" 开头、或把该常量的说明拆成独立 doc 块挂在 :31 的 const 上）。
  **方案**：由票 70 腿在它的波次里顺手落（零语义、1 行移动）；本波不碰。

## 6. 本波未动的共享树既有改动（留痕，防误领）
`git status` 里另有 3 个 tracked 改动属于其他腿/外部 agent，本波零触碰、零提交：
`.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go`、
`.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt`、
`frontend/fixtures/composer-states.html`。

## 7. next=
- 编排者推送后取 CI `lint` job staticcheck 步两段自报读数：期望 pass1≈**3**（SA9009/SA5000/
  plugin ST1022——三条都是已登记特殊项）、pass2≈**95**（78 U1000＋17 非 U1000，其中 14 条
  windows-only 面＝下一波工单清单）；runner go patch 版差异按两遍分开判读。
- 第二类 3 条的落法（ci.yml -SA9009 / plugin 腿的 lint:ignore / plugin 腿的 1 行 doc）待编排者裁。
