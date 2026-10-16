# 165 门铃仪器普查（腿 165-c2·只读·零产码）

- 起手锚＝`232f309e8092b15ede2d6f0455d71775e5a5de96`，`date '+%F %T %z'`＝`2026-10-16 11:13:58 +0800`；交件时 HEAD 已漂到 `bdb69e16`（别腿在落），⇒ **本件全部读数用 `git grep <pat> 232f309e` 显式锚定，⛔ 裸 `HEAD`**。
- 台面＝`git archive 232f309e | tar -x -C $(mktemp -d)` 仓外导出树：六件套 `go.mod go.sum internal cmd tools scripts` 逐枚 PRESENT、派单第 0 步那把尺 `find internal cmd tools scripts -name '*.go' | wc -l`＝**742 枚**（证非空），整棵导出树（含 `.scratch/**` 台件）＝**1098 枚**；跟踪件全仓 **14046 枚**。⛔ 共享工作树字节、⛔ `-overlay`。
- 本腿跑过的 Go 命令（自报）＝`go env GOFLAGS GOTOOLCHAIN GOMODCACHE`（rc=0，读得 GOFLAGS 空／GOTOOLCHAIN=auto／GOMODCACHE=`D:\work\base\gopath\pkg\mod`）＋`go list ./internal/agent ./internal/tools ./internal/llm/...`（rc=0，8 枚包名）。⛔ `go build`、⛔ `go test`（⛔ 与 `270-v1`/`120-v2` 抢编译面与 CPU）。零 `.go` 新建、零 `.go` 修改。
- 原始读数 **41 枚**在 `.scratch/wisp/probes/165/c2/logs/`（逐枚非零字节；`wc -c` 最短一枚 277 字节；本件只放判语，每条带尺）。

## 0. 56／0 复跑＋一枚诚实的更正

| 尺 | 射程 | 读数 | 单位 |
|---|---|---|---|
| `git grep -l 'plan\.present' 232f309e` | 全跟踪树 | **63** | 枚文件 |
| `git grep -n 'plan\.present' 232f309e` | 全跟踪树 | **133** | 行 |
| 同尺喂夹具面（`*/testdata/*`∪`*.golden`∪`*.sse`∪`internal/llm/**`∪`tools/mockllm/**`∪`cmd/llmrecord/**`） | 见上 | **0** | 枚文件 |
| 正控 `tool_calls` 同尺同射程 | 同上 | **20** 枚文件 | 尺不瞎 |

★**编排者那枚 56 我复现了，但它不是恒量**：`git rev-list 71002d0..232f309e`＝**32 枚**，逐枚重跑（logs/16，32 行齐）⇒ 45(10:33) → 54(10:35) → 55(10:41) → **56(10:49 `d646625e`/`d4e93cb8`/`26e1f03b`)** → 58 → 59 → 62 → **63(11:11)**。那枚 56 是真读数，落在 10:49–10:50 那三枚提交上；`165-v1` 自己的锚 `71002d0` 只有 45。**fixture=0 在这 32 枚里每一枚都成立** ⇒ 结论方向不依赖时间戳。

## 1. 问 A｜进模型可见名册那一跳，通到哪一步

尺＝`git grep -n <pat> 232f309e -- internal cmd tools`，产码/测试按 `_test.go` 后缀分栏（logs/hop-census.txt、hop-injection.txt）。

| 跳 | 位置 | 产码调用者 | 测试 |
|---|---|---|---|
| 1 名册构造 | `internal/tools/plan_present.go:138 BuiltinPlanEntries()`；名＝`:72` | 定义 1 枚 | — |
| 2 谁调它 | `cmd/wisp/run.go:826` | **1 枚**（产码 4 行里唯一调用表达式，另 3 行是注释/定义） | `internal/tools/plan_present_165_test.go:76,:286`＋`cmd/wisp/plan_wiring_165r2_test.go:38`（字符串字面量，⛔ 调用） |
| 3 注册 | `run.go:827 reg.Register(e)`→`internal/tools/registry.go:95` | `Register` 产码 4 枚（run.go:522/548/798/827） | 58 行 |
| 4 名册出口 | `internal/tools/bridge.go:235 b.reg.List()` | **1 枚**（宽尺 `git grep -n '\.List()' -- internal cmd` 产码 2 枚，另一枚 `registry.go:178` 是 `Provider.List()`＝另一个接口；logs/17） | 1 |
| 5 模型可见面 | `bridge.go:234 Bridge.Tools()`＝`agent.ToolProvider`（`internal/agent/tools.go:86`）；装配点 `run.go:1062 Tools: rt.bridge` | 消费 **3 枚**＝`internal/agent/loop.go:569,:618,:845` | 0 |
| 6 两档选择 | `internal/agent/inject.go:44 Inject()`→`:62-73` 拼 `ToolDefs`；`plan.present` `Resident:true`（`plan_present.go:129`）⇒ 进 `resident` 支 | **1 枚**＝`loop.go:573` | **0 枚**（`git grep -n 'Inject(' … -- internal cmd` 只有 2 行，全在产码） |
| 7 请求体 | `internal/agent/prompt.go:251 Tools: in.Injection.ToolDefs`；`Build` 唯一产码调用者 `loop.go:579` | **1 枚** | 0 |
| 8 上 wire | `internal/llm/openaichat/request.go:109-120`／`anthropic/request.go:149-158`／`openairesponses/request.go:104` | 3 枚适配器 | 15 行 |

**逐跳点名两问的答**：请求体里那枚 `tools` 数组**由 `internal/agent/inject.go:62-73` 选、由 `internal/agent/prompt.go:251` 挂上 `llm.Request.Tools`、由各适配器 `request.go` 序列化**——⛔ 由 `internal/tools` 拼。`BuiltinPlanEntries()` 的调用者是**生产链**（`cmd/wisp/run.go:826`），⛔ 只有测试。装配根只有 `cmd/wisp/run.go:385 assembleRuntime` 一枚，两个宿主都走它：`run.go:249`（`wisp run`）与 `cmd/wisp/resident_task_source_windows.go:284`（常驻/球/面板宿主）⇒ 面板宿主也够得到这枚门。

**断口（今天不存在的那几跳，指到文件∶行）**：⛔ 在 1–8 这条 Go 链上——链是通的，`go build` 绿不是凭据这一点我用调用点枚数逐跳核过。断口是三处**回读／断言**面：

- **断口① `internal/agent/loop.go:637-640`**：模型那一发真被读成 `c.Name` 并落 `tool_call` 行（`internal/memory/dao_toolcall.go:24` INSERT 的 `tool` 列）。链在、 sink 在，但**没有任何用例让 `plan.present` 走到这两行**——夹具面 0 枚（问 B）。⇒ 生产里这枚门今天**能不能响是未知的**，⛔ 是"不会响"。
- **断口② `cmd/wisp/plan_assembled_165r2_test.go:123`**：真装配已经断到**名册出口**（`rt.bridge.Tools()`，`:130 advertised = i.Resident && i.RiskLevel == memory.RiskL2`、`:157` 红句）。它**没往下走第 6/7/8 跳**。⇒ 跳 6/7/8 今天没有任何一枚用例把 `plan.present` 当对象断过。
- **断口③ `llm.Request.ToolChoice` 生产里从未赋值**：尺＝`git grep -n ToolChoice 232f309e -- internal/agent cmd/wisp` ⇒ **产码 0 枚**（logs/toolchoice-census.txt；字段只在 `internal/llm/provider.go:71` 声明、`content.go:90-99` 定义词汇）。而 `tools/mockllm/chat.go:172-185` 的"必发一枚工具调用"合成**只在 `tool_choice` 非空时发**，且发的是 `req.Tools[0]`（**按位取**，⛔ 按名）。⇒ 拿 mockllm 的合成型当门铃，今天从装配面按不到；只能手拼 `llm.Request`，而那是适配器级测试，量的是编码器、⛔ "模型会不会调"。

## 2. 问 B｜现有夹具面有没有"模型这一轮请求了这枚工具"那一形

**分母（全部取自 232f309e 导出树，尺＝`git ls-tree -r --name-only 232f309e`；⛔ 别腿新落的件没算进来）**：`.sse` **52 枚**／任何 `*/testdata/*` 下文件 **56 枚**／`internal/llm/**` **90 枚**／`tools/mockllm/**` **14 枚**／`cmd/llmrecord/**` **2 枚**／`*.golden` 扩展名 **0 枚**。
★`*.golden`＝0 **不是假 0**：同尺放宽 `find -iname '*golden*'` ⇒ **12 枚命中**（logs/fixture-inventory.txt 旁证），因为 golden 在本仓是**目录名与包名**（`internal/llm/golden/`、`*/testdata/golden/`），⛔ 扩展名。这一枚我先跑正控再报零，正是编排者今天踩过的 `grep -F` 假 0 那一坑。

逐枚点名能不能表那一形（指到行）：

| 面 | 枚数 | 能不能表"这一轮模型请求了这枚工具" | 凭据 |
|---|---|---|---|
| `internal/agent/testdata/golden/*.sse` | 11 | **能**，且已被用成那一形 | `tool-then-text.sse` → `internal/agent/loop_golden_test.go:55 TestGoldenSingleToolCall`，断言在 `:66 call.Req.Name != "echo"` |
| `internal/llm/testdata/golden/*.sse` | 39 | 能表**适配器解码**那一形（`tool-call.sse` 里带 `"function":{"name":"get_weather"}`），⛔ 不经名册/`Inject` | logs/posctl-and-sample.txt 全文 |
| `tools/mockllm/testdata/golden/*.sse` | 2 | 同一格式的镜像回放面（`chat.go:33`/`messages.go:188`/`responses.go:73` 三处入口） | logs/mockllm.txt |
| `internal/llm/golden/`（格式＋回放机，3 枚 `.go`） | 3 | **不能**：`golden.go:44 FormatHeader` 只定格式；`replay.go:33 Requests` 只记 `req.Method+" "+req.URL.Path`，**⛔ 不记请求体** | logs/golden-replayer.txt、logs/17 |
| `cmd/llmrecord` | 2 | **不能**：`main.go:8-11` 自述 **STATUS: STUB**，live 录制 DEFERRED 等 API key（台账 H2），无 key 时显式报错且不写文件 | logs/golden-and-record.txt |

**关键读数（这一条比"plan.present 缺"更硬）**：尺＝`git grep -F -l "<name>" 232f309e -- '*.sse' '*/testdata/*' 'tools/mockllm/**'`，对**本仓真正注册的 11 枚工具名**逐枚量（名册尺＝`git grep -n 'func (.*) Name() string' 232f309e -- internal/tools`，产码 11 行）：

```
fs.read 0  fs.list 0  fs.edit 0  fs.write 0  fs.trash 0  fs.move 0
fs.delete 0  plan.present 0  task.spawn 0  task.output 0  task.cancel 0
正控（同一把 -F 尺、同射程）：echo 17  get_weather 8  sleep 3  get_time 4  tool_calls 15
```

⇒ 夹具面里出现过的工具名**只有四枚假工具**。门铃零读数⛔ 是 `plan.present` 独有的，是**整族**的状态（logs/14-eleven-names-vs-fixtures.txt）。写这一形的人⛔ 在为 `plan.present` 单独开路。

**可复用先例（具名、指到行）**：
1. `internal/agent/loop_golden_test.go:55` ＋ `internal/agent/harness_test.go:92 newHarness(t, fixture, …opts)`／`:59 withTools(tp)`／`:41 bodies [][]byte`（**已经逐枚捕获真请求体**）／`:169 requestBodies()`。
2. `internal/agent/loop_golden_test.go:128 TestGoldenToolResultFeedsNextTurn` 已在断请求体字节（`:138 bodies[1]`、`:143-144 反向断 bodies[0] 不含 "tool_call_id"`）⇒ **乙形的"看请求体"这一手在 `internal/agent` 里已有先例**。
3. `internal/llm/adaptertest/mockllm.go:383 LastRequest(t, route)`（拿回 mock **真收到的请求体**；端点 `tools/mockllm/server.go:112`＋`control.go:145`）。现有使用者：`internal/llm/anthropic/cache_test.go:51`、`internal/llm/openairesponses/protocol_test.go:38`（`:66` 直接断 `sent["tools"]` 恰好 1 枚）、`cmd/wisp/run_test.go:432`。⇒ **这是仓里唯一一枚"真请求体的 tools 数组"仪器，今天已经装好、只是没人为 plan.present 用它**。
4. `cmd/wisp/run_test.go:158 TestRunTextTaskTextPathEndToEnd`＝真 `runTextTask`＋真 mockllm 子进程（`adaptertest.StartMockllm`，`mockllm.go:51`）的端到端先例。
5. **"别枚工具是否已有这形可抄"的答＝形状有、名字没有**：`tool_calls` 这一形在夹具面成族存在（尺＝`git grep -l 'tool_call' 232f309e -- internal/agent/testdata/golden` ⇒ **11 枚里 9 枚命中**，九枚名册逐枚在 logs/18；负控＝不在名单里的 `provider-500.sse` 同尺 0 命中），可抄的是**形状与驱动它的那四条先例**（上面 1/2/3/4）；⛔ 可抄的是"某枚真实 D34 工具已经被人这么测过"——11 枚真实册名逐枚 0 枚。`web.search` 连册都没进（Name() 那把尺的 11 枚产码行里没有它），所以它⛔ 是先例。

## 3. 问 C｜最小仪器三形＋代价表

| 形 | 改动面（文件∶行） | 顶红谁现有的钉（钉名＋尺） | 新增用例 | 碰 `cmd/wisp/**`＋`internal/panel/**` | 碰 `.github/**` |
|---|---|---|---|---|---|
| **甲** 金样本 SSE 带这枚 tool_call | 新 `internal/agent/testdata/golden/plan-present.sse`（形照 `tool-then-text.sse`，`function.name` 换 `plan.present`、`arguments` 填 `{"summary":…,"steps":[…]}`）＋`internal/agent/loop_golden_test.go`（或同包新件）1 枚 | `internal/agent/testdata/golden` **⛔ 孤儿钉**（尺＝`git grep -n ReadDir 232f309e -- internal/agent` ⇒ 3 处全在 `t.TempDir()` 的 artifacts 目录，⛔ goldenDir）⇒ 不红。**⚠ 放错目录就红**：落 `internal/llm/testdata/golden/` 且名以 `chat-`/`anthropic-`/`responses-` 开头 ⇒ `adaptertest.AssertCoverage`（`internal/llm/adaptertest/harness.go:300-316`）＋`openaichat.assertNoOrphans`（`internal/llm/openaichat/harness_golden_test.go:105-122`）两枚孤儿钉同响；要进共享场景表再多顶 `harness.go:282-285 len(table)!=len(AllScenarios())` | 1 | ⛔ 不碰 | ⛔ 不碰（`internal/agent` 已在 CI core scope：`scripts/portable-tests.sh` `core_pin` 第 2 行＋`scope=( ./internal/agent/ … )`） |
| **乙** 装配面断"名册含它且真进请求体" | 乙-1（新件在 `cmd/wisp/`）：接 `adaptertest.StartMockllm`＋`Mockllm.LastRequest(t,"chat")`（`internal/llm/adaptertest/mockllm.go:383`）断 `sent["tools"]` 含 `plan.present`；正控形＝摘 `run.go:826-831`。乙-2（新件在 `internal/agent/`，`package agent`）：`Inject(dir,…)` 断 `ToolDefs` 含它，再接先例 2 的 `requestBodies()` | 乙-1 **必顶** `cmd/wisp/plan_wiring_165r2_test.go:190`（`BuiltinPlanEntries` 产码调用点**恰 1 枚**）＋`:198-199`（那 1 枚**只许**落在常量 `planAssemblyFile165r2 = "cmd/wisp/run.go"`，`:40`）⇒ 任何"在别处再调一次 builder"的写法即红；再撞票 314 `cmd/wisp` 整包死亡区（本波归 `220-r2`）。乙-2 已知钉＝**0 枚**（`Inject` 今天测试调用者 0 枚，只增不减）；软风险＝`internal/agent/prompt_test.go:210` 钉 `SectionResidentTools==1200`（枚数未变 ⇒ ⛔ 红； resident 索引变长是预算压力，⛔ 断言） | 乙-1＝2（含名一枚＋摘注册后不含一枚）；乙-2＝1 | 乙-1 **碰 `cmd/wisp/**`＝本波唯一写腿 `220-r2` 的面**；乙-2 ⛔ 不碰 | ⛔ 不碰（两形所在包都已在 core scope） |
| **丙** 真窗/真云端跑一轮看会不会请求 | **0 枚 `.go`**：`wisp run "帮我把这个项目收拾一下"`，读 `tool_call` 表 `tool` 列（`internal/memory/dao_toolcall.go:24`）＋审计行 `cmd/wisp/run.go:772` | **0 枚**（不入库，⛔ 碰任何字节） | 0 | ⛔ 不碰（只跑现成的 `wisp run`） | ⛔ 不碰 |

**甲/乙 共有一枚必须先裁的硬限制（具名）**：`internal/agent` 目录（不含子包）的测试件 **21 枚，逐枚 `package agent`**（尺＝`git ls-tree --name-only 232f309e internal/agent/ | grep -c '_test\.go$'`＝21，再逐枚 `git show … | grep -m1 '^package '` ⇒ 21/21 全是 `package agent`），且**整个 `internal/agent/**` 树里 `package agent_test` ＝0 命中**（尺＝`git grep -l 'package agent_test' 232f309e -- internal/agent | wc -l`＝0）——而 `internal/tools/bridge.go:13` 已 import `internal/agent` ⇒ **在 `package agent` 里 import `internal/tools` ＝导入环**（外测包 `agent_test` 本可以绕过，但今天一枚没有，新建外测包＝给这一族开新面）。⇒ 甲与乙-2 都拿不到真桥，只能用 `EchoProvider`/自拼 `ToolProvider`——那正是本仓禁止的"用 mock 代替真的"。要在真桥上量这一形，落点⛔ 在 `internal/agent`，只能在 `cmd/wisp`（乙-1）或一枚新的外测包。这一格请编排者裁。

**丙的具名资源需求（⛔ 我不动手，也⛔ 不许被排进派单当可交付）**：① 真 provider API key——`cmd/llmrecord/main.go:8-11` 明写 live 路径 DEFERRED、无 key 时不写文件（台账 H2）；② 真窗要 `//go:build windows && winlive`（★真实拼写带 `windows &&`；**我自己两把尺并列**（logs/17）＝严格尺 `git grep -lE '//go:build.*winlive' 232f309e -- cmd internal tools` **14 枚**（cmd/wisp 9＋internal/ball 5＋别处 0，与 `ci.yml:625` 那句"9 under cmd/wisp, 5 under"逐位对上），松散尺 `git grep -l 'winlive' -- '*.go'` **54 枚**（连带读到 `.scratch/**` 取证台件）——台账 A508 时代那枚"8 枚"在我这两把尺上**都复现不出**，引用时按我这枚 14 走）；CI 只跑 `go vet -tags winlive ./cmd/wisp/ ./internal/ball/`＝`.github/workflows/ci.yml:681`，**编译门、从不真跑**；③ 属〔**仅本机可量**〕，本仓已写死"winlive 里要人眼那几枚不许伪装成 CI 测过"（台账 `A508`，`.scratch/wisp/probes/orchestrator/ci223-1/raw_range_roster.txt:6`）。⇒ **丙是今天唯一一枚零字节就能拿到门铃读数的形**，代价是要机主点头给凭据＋一次真窗。

**本腿建议（⛔ 执行）**：丙先跑（零改动、当场出数，直接回答"模型今天会不会调"）；乙-2 补跳 6/7 那三跳（⛔ 抢编译面）；甲做常驻钉，但**先裁"internal/agent 无外测包 ⇒ 真桥进不去"这一格**，否则甲只是一枚"名字叫 plan.present 的工具被回放了一次"，⛔ 是"模型会调它"。

## 4. ★附带报一句（⛔ 本腿不动手，处置＝编排者／下一枚 `internal/tools` 面空的微腿）

`internal/tools/plan_present_165_test.go` 文件头那两行**逐字**（尺＝`git show 232f309e:internal/tools/plan_present_165_test.go`，awk `NR>=13 && NR<=16`；全文 327 行；logs/08-overlay-comment-verbatim.txt）：

```
[L14]// 两发变异（-overlay 注入、树零改）见 probes/165：把 Declared 从 L2 降成 L1
[L15]// 必须叫 AC#2 那条红；摘掉 planPresentScopeClause 必须叫 AC#3 的条款那条红。
```

矛盾点复认（⛔ 我重裁）：①「两发」与同腿 `165-r1` 的读数件/提交信息（M1＋M2＋M3 三发）不符；②「`-overlay` 注入、树零改」与 `165-v1` 台面（仓外备份＋逐枚对 HEAD blob sha256＝全 MATCH）不符；③ `-overlay` 冒充物理突变是本编队明令禁止的形状。⇒ 该改的是描述面，⛔ 任何判语。

## 5. 交件自证三把尺

1. 提交＝日志那一枚 `606bb2ff`（`.scratch/wisp/probes/165/c2/**` 38 枚）＋本件那一枚（号见交件回信）；逐枚 `git show --stat` 只含 `docs/evidence/s1/165-doorbell-census.md` 与 `.scratch/wisp/probes/165/c2/**` 两处写面。
2. 问 A 断口＝`internal/agent/loop.go:637-640`（回读侧无用例）／`cmd/wisp/plan_assembled_165r2_test.go:123`（只到名册出口）／`llm.Request.ToolChoice` 产码 0 枚（`internal/agent/prompt.go:247-255` 不置该字段）。
3. 问 B 分母名册＝52 `.sse`／56 `*/testdata/*`／90 `internal/llm/**`／14 `tools/mockllm/**`／2 `cmd/llmrecord/**`／0 `*.golden`（正控 12）；11 枚真注册名逐枚 0，正控 echo 17。
4. 问 C 三形改动面逐枚指到文件∶行（上表），顶红名单逐枚带钉名与尺；⛔ 猜。
5. 答不上的问＝**无**。⛔ 本腿未回答"模型今天到底会不会调"——那需要丙的凭据，本腿⛔ 有，已具名写成资源需求。
