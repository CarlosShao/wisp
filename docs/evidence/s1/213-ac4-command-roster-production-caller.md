# 213-ac4 实现腿 — 命令名册的生产调用者（`wisp commands`）：AC#4 的「名册可见」落地

- 腿 `213-ac4`（**实现腿**，⛔ 裁决者）。HEAD `f4ad34ee1c`（分支 `dev`，**非 detached**，非旧提交）。量于 2026-11-05 11:44–12:3x +08。
- **写面（5 枚，逐枚现读）**：`cmd/wisp/commands.go`（新建 164 行）· `cmd/wisp/commands_213_test.go`（新建 169 行）· `cmd/wisp/main.go`（+13/-0：usage 一段 ＋ dispatch 一支）· `internal/command/roster.go`（+43/-9：`Approval.String()`、`Command.Door()`、两处注释）· `internal/command/roster_test.go`（+117/-3：AC#4 的盘面读数用例 ＋ 头注一条）＋本文件。
- **⛔ 零越界**：`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、`internal/panel/**`、`internal/tools/**` 一枚字节未动（§⑤ 有逐字命令与读数）。C17 method 常量 **7 枚，本腿零新增**。
- 本腿 **⛔ 不翻票上任何框**（含 AC#4）、**⛔ 不碰台账／HANDOVER**、**⛔ 不 push**。

---

## ⓪ 一屏判（六句，逐条落点见后）

1. **AC#4 的「名册可见」今天有生产调用者了**：`Production()` 的生产调用者 **0 → 1**（`cmd/wisp/commands.go:98`）。r1 交回的那一格（`213-command-roster-r1.md` ④.1／⑦.2「`Production()` 的生产接线＝没做」）由本腿补上，量法是**盘面**（§③.1），不是"测试里调过了"。
2. **名册没有被复制一份**：`wisp commands` 打出来的每一行都由 `command.Production()` 的行派生（名字／一句话说明／参数形状／要不要人批／门），页脚的控制词与 C17 方法名也**从行里读**（`e.Utterance`、`e.Method`），本文件里没有第二个拼写。M1／M2 两发变异逐发把它打红。
3. **AC#2 在生产输出里也成立**：名册没有的能力（`compact`／`goal`／`plan`／`model`／`export`／`feedback`）在输出里**逐枚缺席**；M1 塞一条假行（`compact`）⇒ 四句红同时出现（不在名册／行数多一／无门能力／表头计数）。
4. **AC#4 的终态是具名写出来的，不是推断的**：输出末段逐字写「界面点不了：聊天框加号菜单今天还读不到这张表（`panel.Snapshot` 没有 commands 键、C17 没有列名册的路由）」；M3 把这句改掉即红。
5. **这一腿零执行**：`commands` 不解析数据根（`grep -c resolveDataDir cmd/wisp/commands.go` ＝ **0**）、不开 store、不调任何门、不装日志 sink（`grep -c installLogSink` ＝ **0**、`grep -c "slog\."` ＝ **0**）⇒ 票 128 的「解析数据根的腿」表**不增**（表里那 1 处 `commands` 是既有的说明句，`:420`），票 131 的「装监听／记记录」义务也**不增**（`leg commands → no records`）。
6. **门禁**：`go build ./cmd/... ./internal/...` **rc=0**；`go test ./internal/command` **rc=0**（`ok 0.174s`）；`cmd/wisp` 定向 **2/2 PASS**、两张名册门（133／131）**PASS**；`gofumpt -l` **空**；`sh scripts/d22scan.sh` **rc=0／clean**。**cmd/wisp 整包有 45 枚红，逐枚在 HEAD 对照里同红**（§⑤.3 是同一 45 枚名的对照读数，本腿零归因）。

---

## ① 起手三核（THREE-CORE，逐条现跑，⛔ 引旧叙述）

### ①.1 HEAD 与分支

```
$ git rev-parse --abbrev-ref HEAD   → dev
$ git rev-parse --short=10 HEAD     → f4ad34ee1c
```

派单写的 HEAD `f4ad34ee1` ＝ 现读 `f4ad34ee1c` 前 9 位；分支 `dev`、**非 detached** ⇒ 这一核通过。

### ①.2 票面串／行在否（六条，含漂）

| 票面／账里写的 | 现读（本腿现跑） | 判定 |
|---|---|---|
| `internal/panel/l2_grant_boundary_test.go` 的 `poolJudgedByRealGuard`，「**exactly the shape a handler registry takes**」 | `grep -c poolJudgedByRealGuard` ＝ **4**；逐字句在 **`:1110`** | **在**（这条钉就是「名册不落在 `internal/panel` 里」的理由） |
| `internal/agent/control.go:36` 一带是「名字→本地执行」的现成解析器 | `:36` ＝ `var controlWords`（词表）；`:51` ＝ `func MatchControl` | **在**（解析器本体在 `:51`，票面写的 `:36` 是同文件"一带"） |
| `loop.go:341`（`213-c1` 报的 `MatchControl` 调用点） | **`:383`** ＝ `if verb, ok := MatchControl(input); ok {`；`defaultControl` 在 **`:622`** | **漂 +42**（r1 记同一读数；本腿复核为真） |
| `internal/tools/bridge.go:42-45` 的 C17 闭集 **4 枚** | `grep -cE 'Method[A-Za-z]+ *= *"'` ＝ **7**（`:42 :43 :44 :45 :66 :67 :75`） | **票面数过期**（4→7）；本腿**零新增**，禁区照旧 |
| `internal/panel/instructions_200.go:26-40` 五枚态 ＋ `:59 Dropped` | 五枚态在 **`:24-42`**（loaded／none_found／off／refused／not_run） | **在**（本腿只引它当"不新造态"的先例，未使用） |
| `internal/panel/composer.go:238-241` 附件四枚快照字段 | 该行现是 `WorkspaceView` 的注释；附件字段已移位 | **行号过期**（账 `A416` 已记口径）；本腿不依赖它 |

### ①.3 红绿是否他人在飞件（时间线，⛔ 把别人的半保存态算成本腿的绿／红）

| 时刻 | 读数 |
|---|---|
| 11:44（起手） | `git status --porcelain -- internal cmd` ＝ **3 行**：` M internal/tools/spawn_roster_census_229_test.go`（**229 腿在飞**）＋ 两枚未跟踪（`cmd/wisp/zz_probe_temp_test.go.txt`、`internal/tools/ticket23_image_channel_test.go`，都不是 `.go` 名下的编译件） |
| 11:56 | 我省道门与新增用例第一次跑：**全绿**（读数为证：§⑤） |
| 11:57–11:58 | **21 腿**把 `internal/agent/approval/gate.go` 改到半保存态：`gate.go:396:60: g.remainingIn undefined` ⇒ 我 M3／M4 两发变异第一次是 **`[build failed]`**（不是我要的红）。⛔ 我没有去动它。 |
| 11:58:52 | 该腿保存完毕，`go build ./internal/agent/approval` rc=0；我**重跑** M3／M4 拿到真红（§④） |
| 12:0x–12:2x | 该腿的 `gate.go`／`ui.go` 仍是 ` M`；本腿所有门禁读数**按我的包／我这一支**解读，`cmd/wisp` 整包的 45 枚红在 §⑤.3 用 HEAD 对照逐枚归因 |

---

## ② AC#4 的读法、落点，以及为什么不选另外两条

**票面原文（逐字）**：「**AC#4**：**入向那一跳没接通之前**，本票的终态是「名册可见＋命令行可执行」，**并具名写「界面点不了」**——这一格不许自己翻勾。」

**r1 交回时的状态**（`213-command-roster-r1.md`）：名册可枚举（`Production()`）、两条门都有生产调用者、**但 `Production()` 自己被 import 0 次**——"名册可见"只对测试成立，**命令行上没有一个入口能把它打出来**，于是"＋命令行可执行"这半句里的"名册"是空的。ADR＝本腿补这一跳。

**落点：`wisp commands`（cmd/wisp 的 console 腿）。** 三条理由，按价值排序：

1. **它让"名册可见"变成产品事实，而不是库事实**：一条 operator 能敲的命令，把 `Production()` 的每一行打出来；`Production()` 从此有生产调用者（§③.1 的盘面尺）。
2. **它不动任何契约面**：不加快照键（页面顶层键今天**恰好四枚**，`internal/panel/pump_test.go:123`／`:291` 逐字钉着 `composer,generatedAt,pending,results`）、不加 C17 方法（闭集 7 枚不变）、不碰 `internal/panel/**`（票 213 禁区）。**页面那一格照旧是红的**，由票 213「交付 3」的腿动两侧尺。
3. **它把"命令行可执行"讲清楚**：页脚从**行**里取控制词与 C17 方法名，打出两条门的命令行拼写（`wisp run "<控制词>"` ／ `wisp panel-inbound -data <数据根>` ＋ 一行封套）——不是新造执行器，是把已有两条门的位置告诉用户。

**⛔ 没选的两条，连同理由**（留给后来人复核，不是"没想过"）：

- **快照键／C17 方法**：撞禁区（不新增 C17 方法名，要先落 `A##`）＋两侧同批的契约尺＋`internal/command -> internal/panel` 的导入方向（反过来组装就是循环 import，r1 §③.3 已写死）。
- **常驻进程启动时打一行日志**：它会让 `Production()` 有调用者，但用户**看不见**名册（日志不是可读的名册），也就不满足 AC#4 的"可见"。本腿不拿"有调用者"偷换"可见"。
- **`wisp commands judge <name> [target]`**：`Judge` 要"当前有效档 ＋ 有无 L2 腿"，就得开 store／解析数据根——那是票 128 表要登记的新腿、也是把 AC#3 的判决搬进 CLI；AC#3 的判决已在 r1 的 18 行矩阵里与**真门**对撞过，本腿不复制第二遍。**今天 `wisp commands` 只读不判**，`mode` 行明说"人批：看参数（要不要卡由门自己的闸答）"。

---

## ③ 判据逐条（缺一条就具名报"没测"）

### ③.1 AC#1「名册可枚举，且每一条都有真执行者（生产路径上有调用方）」

- **名册本身的生产调用者**（本腿新增的那一格）：`TestTheProductionRosterHasAProductionCaller213` **读盘**——扫 `cmd/`＋`internal/` 的**非测试** `.go`，按 AST 找"import 了 `internal/command` 且调了 `Production()`"的文件。
  - **RED-before（落刀前）**：`roster_test.go:550: AC#4 RED: no non-test file under cmd/ or internal/ imports "github.com/CarlosShao/wisp/internal/command" and calls command.Production().`（当时读数＝0）
  - **GREEN-after**：`roster_test.go:556: AC#4: production callers of command.Production(): cmd/wisp/commands.go (1 call site(s))`
  - 独立复跑（人可逐字重跑）：`grep -rn --include=*.go "\.Production()" cmd internal | grep -v _test.go` ⇒ `cmd/wisp/commands.go:98`（另一命中是 `:22` 的说明句）。
- **每一条的真执行者**：r1 已逐门实测（`mode` 打真 `panel.ModeWriteHandler`＋真 `perm.Store`；`cancel`／`stop` 打真 `agent.Loop` 停真任务），本腿**未改门、未改行**，只把这三行打到控制台上（输出见 §③.5 的逐字）。
- **⚠ 本腿仍未做到的**（照实报）：页面**还是读不到**这张表（没有 commands 键）。AC#1 的"生产路径上有调用方"对"门"成立于 r1、对"名册"成立于本腿；"可见"今天的射程＝**控制台**。

### ③.2 AC#2「没有能力的命令进不了名册」——正控照做

- 负向七连在 r1（`Lookup("/mode")`／`compact`／`goal`／`plan`／`model`／`export`／`feedback` 逐枚空）。本腿在**生产输出**上再钉一遍：`Test213CommandsLegPrintsTheProductionRoster` 断言这六枚名字**不在打印出来的行里**（`compact/goal/plan/model/export/feedback`），并把"打印出来的名字集合 ≡ 名册名字集合"钉成集合相等（少一行／多一行／重名都红）。
- **正控（手工塞一条假命令 ⇒ 判据必红）**：M1（§④）在 `cmdCommands` 里 append 一条 `compact` 行 ⇒ **四句红同时出现**，其中一句逐字是「row "compact" is printed and is not in the roster: this leg must not carry a table of its own (AC#2)」，另一句是「printed 4 rows, roster has 3」。⇒ 名册是唯一来源，不是注释。

### ③.3 AC#3「档位与权限」——本腿**不新造**判决，只换一个出口

- 判决的唯一来源仍是 `perm.JudgeCard`（r1 的 18 行矩阵与真门对撞：13 允许／2 先问／3 拒、24 行审计），本腿**一行判决逻辑没加**：`cmd/wisp/commands.go` 不 import `perm`／`risk`，`grep -c "JudgeCard\|risk\." cmd/wisp/commands.go` ＝ **0**。
- 生产输出给的是**这一行门与它的闸**：`mode` 行逐字打「人批：看参数（要不要卡由门自己的闸答）」，并在门那行点名 `panel.ModeWriteHandler -> perm.Store.Set`。**拒了要说为什么**：这条命令今天不判具体请求，所以不会静默拒；真正的拒句在 r1 的矩阵里（每一行都带句）。
- **⚠ 具名没测**：`wisp commands` **不**回答"目标档 X 在这一档会不会被拒"（见 §② 末尾为什么不做）。这一格若要，得单开一条腿（开 store ＋ 解析数据根 ＋ 落 `A##`）。

### ③.4 AC#4「终态＋具名写界面点不了」——**不翻框**

- **终态**：`wisp commands` 打出 3 行（`cancel`／`stop`／`mode`），末段逐字：
  - 「怎么从命令行执行（票 213 AC#4 的终态：名册可见 + 命令行可执行）」
  - 「控制词门（D11(1)，不经 LLM、不要人批）：wisp run "取消"  /  wisp run "停"」
  - 「面板门（C17 既有方法，不新增方法名）：wisp panel-inbound -data <数据根>」＋ 一行封套 `{"method":"panel.mode.request","requestId":"<请求 id>","source":"panel-composer","to":"<目标档位>"}`
  - 「**界面点不了**：聊天框加号菜单今天还读不到这张表（panel.Snapshot 没有 commands 键、C17 没有列名册的路由）。接上那一跳的腿要动快照契约尺的两侧，属票 213 交付 3，到了那一天，这一行和 internal/command 的包注释要一起改。」
- **框**：AC#4 **保持 `[ ]`**（⛔ 本腿不翻；"界面点不了"是真话，翻框要等页面那一跳）。
- **票面那句从句的前提仍过期**（r1 已报：`postMessage → wispDispatch → dispatchRaw → ComposerDispatch.Handle` 已在生产路径上）；本腿不改票面一字，把它当"不许自己翻勾"的另一半照办。

### ③.5 AC#5「执行全部走现成审计，不留第二条旁路」

- **结构**：`commands` 不执行——`grep -c "resolveDataDir"`＝0、`installLogSink`＝0、`slog\.`＝0；它也不调两条门（`grep -c "HandleModeRequest\|\.Set("`＝0）。它打出来的每一行都指向 r1 已实测过审计的那扇门。
- **本腿新增的东西对审计面是零**：没有新执行路径 ⇒ 没有需要新记账的执行。（`mode` 走 `perm.Store` 的 `[audit]`，控制词走 `EvControl` 事件——后者缺 `task_log` 行是 r1 §④.5 已具名交回裁决者的那一格，本腿未替它裁、也未新增。）
- **唯一新增的"输出"**：stdout 的名册文本。它不进 `task_log`、也不是审计——一条 listing 不该被读成记录（票 131 的 `records=false` 就是这个读数）。

---

## ④ RED-before / GREEN-after ＋ 五发变异（逐发带红句，还原后 `sha256sum -c` 全 OK）

**GREEN（终态，逐字）**

```
$ go test -count=1 -timeout 200s -run "213" -v ./cmd/wisp
=== RUN   Test213CommandsLegPrintsTheProductionRoster
--- PASS: Test213CommandsLegPrintsTheProductionRoster (0.00s)
=== RUN   Test213CommandsLegRefusesArguments
--- PASS: Test213CommandsLegRefusesArguments (0.00s)
ok      github.com/CarlosShao/wisp/cmd/wisp     0.145s

$ go test -count=1 -timeout 200s ./internal/command
ok      github.com/CarlosShao/wisp/internal/command     0.174s
```

**RED-before（两把，都不是"编译不过"）**

1. **盘面尺**：落刀前 `TestTheProductionRosterHasAProductionCaller213` 逐字红——`roster_test.go:550: AC#4 RED: no non-test file under cmd/ or internal/ imports "github.com/CarlosShao/wisp/internal/command" and calls command.Production().`（. `internal/command` 被生产代码 import 的次数＝0，与 r1 §⑦.2 的读数同源）
2. **红因落在我这一格上**：`go test ./cmd/wisp -run 213` 在落刀前**没有**这两个用例可跑（`cmdCommands`／`commandsIO` 不存在），所以第二条 RED 用下面的变异打（`M1`–`M5`），每发都让判据**因自己的理由**变红。

| 变异 | 刀口 | 红名／红句（逐字，台件 `.scratch/wisp/probes/213/ac4/m1..m5.txt`） |
|---|---|---|
| **M1**（AC#2 牙：名册不是表） | `cmdCommands` 里 `roster.Entries()` 改成 `append(roster.Entries(), Command{Name:"compact", …})` | `--- FAIL: Test213CommandsLegPrintsTheProductionRoster`；四句：`:48` 表头计数（`3 条命令` 不见了）＋`:60` `row "compact" is printed and is not in the roster: this leg must not carry a table of its own (AC#2)`＋`:66` `printed 4 rows, roster has 3: [cancel stop mode compact]`＋`:124` `row "compact" is printed, and this machine has no door for it` |
| **M2**（页脚不许自己拼控制词） | 页脚 `fmt.Sprintf("wisp run %q", e.Utterance)` 改成写死 `"取消"` | `--- FAIL: Test213CommandsLegPrintsTheProductionRoster`；`:99: footer does not name the control row "stop"'s command-line spelling "wisp run \"停\""`（⇒ 页脚是**行**派生的） |
| **M3**（AC#4 具名句） | 末段 `界面点不了：…还读不到…` 改成 `界面可以点：…能读到…` | `--- FAIL: Test213CommandsLegPrintsTheProductionRoster`；`:116: output does not say the page cannot ring these commands yet: AC#4's terminal state is 名册可见 + 命令行可执行, 界面点不了` |
| **M4**（退出码不是装饰） | `cmdCommands` 参数错分支 `return 2` 改成 `return 0` | `--- FAIL: Test213CommandsLegRefusesArguments`；`:136: cmdCommands(list) = 0, want 2 (usage error, the code the other console legs use)` |
| **M5**（AC#4 的那一跳本身） | `cmd/wisp/commands.go` 的 `command.Production()` 改成 `command.New()`（＝有 print 腿、没生产调用者） | 两处同时红：`roster_test.go:550: AC#4 RED: no non-test file … calls command.Production()` ＋ `commands_213_test.go:106: footer does not name the panel door "panel.mode.request"`（⇒ 盘面尺量的正是"这一跳"，不是"这个文件存在"） |

**五发都不挂死**（每发 < 0.3 s 变红；M5 特意用"能编译、只是不调 `Production()`"的形，避免"靠编译不过或超时变红"的假牙）。**还原**：五发逐一 `cp` 回备份，`sha256sum -c sha-mine.txt` 五枚全 `OK`。

---

## ⑤ 门禁读数（逐条现跑，台件 `.scratch/wisp/probes/213/ac4/final-*.txt`）

### ⑤.1 我的包与我的支

| 门 | 命令 | 读数 |
|---|---|---|
| go build | `go build ./cmd/... ./internal/...` | **rc=0** |
| go test（本包整包） | `go test -count=1 -timeout 200s ./internal/command` | **rc=0**，`ok 0.174s`（9 枚用例，含新增的 AC#4 盘面尺） |
| go test（本支定向） | `go test -count=1 -timeout 200s -run "213" -v ./cmd/wisp` | **rc=0**，2/2 PASS |
| 票 133 名册门 | `go test -count=1 -timeout 300s -run "TestAC1AC2DispatchHopGate133\|TestAC4EveryLegIsNailedOrRuled" -v ./cmd/wisp` | **PASS**（两张都 PASS）；新腿行逐字：`leg commands  main.go:108  installs=false handoff=false covered=test Test213CommandsLegPrintsTheProductionRoster drives cmdCommands entries=attachParentConsole\|cmdCommands`；票 131 的门把它记成 `leg commands  main.go:108 … -> no records`（不装监听、不记记录 ⇒ 不欠 nail） |
| 名册/用法对账 | 同上（133 的 `censusVsUsage133`） | 新谓词 `usage` 段落含 `  wisp commands`；门自报 **14 legs**（原 13 ＋ `commands`）、**478 startable cases**（原 476，＋我这两枚）、`covered=test` **5 distinct**（原 4） |
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/commands.go cmd/wisp/commands_213_test.go cmd/wisp/main.go internal/command/roster.go internal/command/roster_test.go` | **空**（0 行） |
| go vet | `go vet ./internal/command ./cmd/wisp` | **rc=0** |
| d22scan | `sh scripts/d22scan.sh` | **rc=0**，`clean - no D22 ban violations`；扫描面 `internal/` **880** Go 文件、`cmd/` **205**、`frontend/` 116、`design/` 39；正向控制（`runtests.sh`）在脚本第一步跑过 |

### ⑤.2 越界读数

```
$ git status --porcelain -- frontend design docs/PLAN.md docs/specs thresholds.go internal/tools internal/panel
?? design/doubao/...        ← 别人的未跟踪件（owner 设计目录，本腿零写）
?? internal/tools/ticket23_image_channel_test.go  ← 别的腿的未跟踪件
（无 ' M'/'M '/'MM' 行 ⇒ 本腿在禁区上零改动）
```

**本腿的 5 枚**：` M cmd/wisp/main.go`、` M internal/command/roster.go`、` M internal/command/roster_test.go`、`?? cmd/wisp/commands.go`、`?? cmd/wisp/commands_213_test.go`。另有两枚**别人的** ` M`：`internal/agent/approval/gate.go`、`internal/agent/approval/ui.go`（21 腿在飞，见 §①.3）。

### ⑤.3 `cmd/wisp` 整包：45 枚红，**逐枚与 HEAD 对照同红**（不归因给本腿）

| 读数 | 命令 | 值 |
|---|---|---|
| 整包（1200 s 上限） | `go test -count=1 -timeout 1200s ./cmd/wisp` | **rc=1**；`FAIL … 628.520s`；`--- FAIL` **45 枚**（名录见 §⑤.3 末） |
| 整包（200 s 上限，派单原令） | `go test -count=1 -timeout 200s ./cmd/wisp` | **rc=1**；`panic: test timed out after 3m20s`，`running tests: TestCleanCheckoutBuilds_AC11 (35s)` ⇒ **200 s 不够这一包**（这个用例在仓外拉一份干净树再 build），**不是本腿引入** |
| **HEAD 对照** | 把本腿 5 枚**全部**退回 HEAD（`git show HEAD:…` 三枚 ＋ 两枚新件移出树），跑**同一张 45 枚名**（`-run '^(…)$'`） | **rc=1**，`--- FAIL` **45 枚**，`diff <(sort 本次) <(sort HEAD 对照)` ⇒ **`same set? YES`**（29.7 s 跑完） |

⇒ 结论：这 45 枚红是**这棵工作树（HEAD ＋ 他人在飞件）在这个环境里**的性质（绝大多数是 `15e*`／`41a2r1*`／`26P*`／`392*`／`381r2*`／`390r1*`／`394*`／`382WR2*`／`WireR1*`／`16r1*` 这一族**真机语音**用例，另有 `TestTicket255RosterStillMatchesTheActualReadSites`＝别人已提交的 config 读取面漂移、`TestAC4FocusReturnToPriorWindowGap33r5`／`TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` 等桌面/进程用例），**与本腿的改动无因果**——同一张名在去掉本腿改动后**逐枚同红**。
另外：`TestDeadPanelThreadRefusesEveryGesture319` 在整包里红、**单独跑 PASS**（`-run '^TestDeadPanelThreadRefusesEveryGesture319$'` ⇒ `ok 0.189s`）；整包那一次是它在满载下与 `TestCleanCheckoutBuilds_AC11` 抢时间（r1 的"重包同一时刻只跑一个"规矩的同一族）。
⇒ **本腿的绿，只报"我的包＋我的支＋两张名册门"这四处；整包不报绿，也不把别人的 45 枚红写成"我的绿"。**

---

## ⑥ 本腿**没做**的（具名交回，不许被读成"做了"）

1. **页面拿到名册**（票 213 交付 3／AC#4 的"可见"里那半）：没做。要动 `panel.Snapshot` 的键（今天恰好四枚，两侧同批的尺）或 C17 闭集（7 枚，禁区要先落 `A##`）⇒ 属另一条腿；本腿在输出末段与包注释两处都具名写了"界面点不了"。
2. **`wisp commands judge <name> [target]`**：没做（要开 store／解析数据根／落 `A##`；AC#3 的判决已在 r1 与真门对撞过，本腿不复制第二遍）。
3. **`mode` 之外的命令行**：没加。`/compact`、`/goal`、`/plan`、`/model`、`/export`、`/feedback` 的缺口与前置在 r1 §② 已逐条具名，本腿**不伪造行**。
4. **控制词的 `task_log` 审计行**：没加（r1 §④.5 已交回裁决者；本腿不擅自加）。
5. **票面／台账／HANDOVER 一字未改**（含 AC#4 的框、`Status` 行、八条 AC 的文字）。

---

## ⑦ 本腿交付物（逐枚）

| 文件 | 是什么 | 读数 |
|---|---|---|
| `cmd/wisp/commands.go` | `wisp commands` 腿：`cmdCommands`（打 `command.Production()` 的行）＋ `commandsIO`（两流，可注入）＋ usage ＋ 页脚（控制词与 C17 方法名**从行里读**） | 新建 164 行 |
| `cmd/wisp/commands_213_test.go` | 2 枚用例：行集 ≡ 名册（名字/说明/参数/人批/门逐行对撞）＋ 六枚无能力名缺席（AC#2 在生产出口）＋ 页脚派生 ＋ "界面点不了"具名 ＋ 参数错→rc2 且 stdout 零行 | 新建 169 行 |
| `cmd/wisp/main.go` | usage 一段（`wisp commands`）＋ dispatch 一支（`case "commands"` → `attachParentConsole()` → `cmdCommands`） | +13/-0 |
| `internal/command/roster.go` | `Approval.String()`（菜单/控制台的人话）· `Command.Door()`（门的人话，**标识符取自行**）· 包注释 AC#4 段改成"有生产调用者、页面仍无键"· `Production()` 注释改成点名调用者 | +43/-9 |
| `internal/command/roster_test.go` | 新增 `TestTheProductionRosterHasAProductionCaller213`（盘面尺：AST 扫 import ＋ `Production()` 调用）＋ 头注 AC#4 那条改写 | +117/-3 |
| `docs/evidence/s1/213-ac4-command-roster-production-caller.md` | 本文件 | 新建 |
| `.scratch/wisp/probes/213/ac4/**` | 台件：`sha-mine.txt`／`sha-before.txt`／`m1..m5.txt`／`focused-255.txt`／`focused-319.txt`／`fails-1200.txt`／`fails-head.txt`／`baseline-head.txt`／`final-*.txt`／三枚备份（`commands.go.bak`／`main.go.bak`／`roster.go.bak`／`*_test.go` 两枚 mine） | **不入库** |
