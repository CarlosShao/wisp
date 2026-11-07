# 213-d3-roster 实现腿 — 交付 3 的 Go 侧一半：名册上快照（`panel.Snapshot` 的 `commands` 键）

- 腿 `213-d3-roster`（**实现腿**，⛔ 裁决者）。HEAD `b71187fce`（分支 `dev`）。量于 2026-11-07 18:03–18:35 +08。
- **本腿接的是什么**：票 213 的名册骨架（`internal/command`）与 AC#4 的控制台半（`wisp commands`）已由前两腿（`213-r1`＝`84e4cc78d`、`213-ac4`＝`e198c3fd2`）落盘；本腿落的是**交付 3「面板可见」的 Go 侧一半**——名册进快照，页面侧 TS 声明照票面口径留给界面那支（Q-51）。
- **⛔ 零触碰**：`internal/ball/**`（wave65/66 在飞，本腿零字节）、`internal/tools/**`、`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、C17 闭集（7 枚，零新增）、台账、HANDOVER、**票面 AC 框与 Status 一枚未动**。⛔ 不 push。

---

## ⓪ 一屏判

1. **名册今天上了快照**：`panel.Snapshot` 新增 reader-gated 顶层键 `commands`（`internal/panel/commands_roster_213.go` 载体＋`pump.go` 读口＋`cmd/wisp/run.go:1023` 接线）。pump **没有**读口 ⇒ 键不在线上（三枚四键字节钉 `pump_test.go:123`／`:291`／`subagent_roster_197_test.go:214` 保持原样绿）；有读口 ⇒ 键必在、行逐字段verbatim（新用例钉死）。
2. **名册没有被复制**：快照行由 `command.Production()` 逐行派生（`cmd/wisp/panel_commands_213.go`），每个字段取自行自己的渲染（`e.Approval.String()`／`e.Door()`）；装配器里**零个字面命令名**（测试钉：`Name: "` 出现即红）。控制台腿（`wisp commands`）与快照读的是**同一张表**——DSH/minimax「同一枚解析器供显示与执行」的落点，继续保住。
3. **缺陷不上妆成空菜单**：`Production()` 的 AC#2 闸若拒了生产表本身，线上携 **Reason 句＋空 rows**（不是沉默的空数组）——A408「文案在、控件不在」同族的反向封堵。
4. **执行通道零新增**：本腿不新增 C17 方法、不接 message 门、不让页面能"执行"任何命令——`commands` 键是**显示数据**，控制词的页面入口今天依旧不存在（页面文本到不了 `agent.Loop`）；footer 与包注释两处旧话已按新事实同步。
5. **门禁**：gofumpt 点名 10 枚 **空**；`go vet ./internal/panel ./cmd/wisp` **rc=0**；`go build` 两包（＋internal/command）**rc=0**；`go test ./cmd/wisp -run 'Test213|TestComposer'` **4 枚 PASS**；`go test ./internal/command` **ok**；`internal/panel` 全包 **3 枚红**，逐枚归因见 §④（2 枚＝交付 3 预定红，1 枚＝**既存红**与本腿无关）；d22scan **clean**。

---

## ① 起手三核

1. **票面串现跑**：`poolJudgedByRealGuard` 与「a second switch/if chain … handler registry」逐字在 `l2_grant_boundary_test.go:1104-1118`（名册不落 `internal/panel` 的理由——**名册表本体仍在 `internal/command`**，本腿进 panel 的只是**快照载体**（两个 struct，零 route 形状 switch/if），零撞钉；本腿收工后该测试整包未见它红）。
2. **前腿产出现量**：`git log -- internal/command` ＝ `84e4cc78d`（213-r1）→ `e198c3fd2`（213-ac4）；`command.Production()` 生产调用者 1（`cmd/wisp/commands.go:98`）。名册今天 3 行：`cancel`／`stop`（控制词门）＋`mode`（C17 `panel.mode.request`）。票面八条里无门的六条（compact/goal/plan/model/export/feedback）在双侧测试里钉着**缺席**。
3. **红绿归因**：`TestC21DesignTokensFourWayAgree` 红＝读 `internal/ball/tokens.go`＋`docs/evidence/s1/c21-native-tokens.md`，两件 `git diff HEAD` 均**空**（GlassRing/GlassHi/GlassHiSoft 在 HEAD 的 `tokens.go` 里就是 0 命中），ball 另有 wave65/66 在飞件（`M internal/ball/ball_windows.go` 等 3 枚）——**既存红，与本腿零因果**。`ComposerState emits [files]` 红＝票 190 的键在 HEAD 就已发射而 panel.ts 未声明（`git show HEAD:internal/panel/composer.go` 含 `Files *TreeView`＝1）——**既存红**。

---

## ② 读数表（本腿改了什么）

| 文件 | 是什么 |
|---|---|
| `internal/panel/commands_roster_213.go` | **新件**：`CommandRow`（name/summary/args/approval/door 五键）＋ `CommandRosterSection`（rows＋`reason,omitempty`）＋头部裁决注释（不判级、不执行、契约移动是显式的） |
| `internal/panel/composer.go` | `Snapshot` 加 `Commands *CommandRosterSection \`json:"commands,omitempty"\``（reader-gated 家族第三枚；注释写明 Q-51 与两枚尺的红是预定程序） |
| `internal/panel/pump.go` | `PumpSources.Commands func() CommandRosterSection` 读口＋`Snapshot()` 接线（nil rows ⇒ `[]`；reason 缺席即不发）；头部旧段「NO NEW TOP-LEVEL KEYS」后补既成事实段 |
| `internal/panel/commands_roster_213_test.go` | **新件**：`TestTheCommandsKeyTravelsOnlyFromItsReader213`——无读口无键（四键）／有读口逐字段 verbatim／nil rows ⇒ `[]` 非 null／缺陷形（Reason＋空 rows）原样上线 |
| `cmd/wisp/panel_commands_213.go` | **新件**：`commandRosterSection213()`（进程级 build-once，表不可变）＋`buildCommandRoster213()`（`Production()` 逐行派生；装配失败 ⇒ Reason 带闸的原句） |
| `cmd/wisp/run.go` | 装配根字面量 **+1 行**：`Commands: commandRosterSection213,`（`:1023`） |
| `cmd/wisp/panel_commands_213_test.go` | **新件**：`Test213CommandRosterSectionIsTheProductionRoster`（行集≡名册，五字段逐一对撞）＋`Test213AssemblyRootCarriesTheCommandsReader`（盘面：run.go 带读口行＋装配器源自 `command.Production()` 且零字面命令名） |
| `cmd/wisp/commands.go` | footer 两处＋usage 一处：旧话「panel.Snapshot 没有 commands 键／界面点不了」改为新事实「名册在快照里（Go 侧）＋TS 声明未写（Q-51）＋控制词执行照旧只有命令行」 |
| `cmd/wisp/commands_213_test.go` | (5) 断言同步：footer 必须说「名册在快照里」「commands 键」「Q-51」 |
| `internal/command/roster.go` | 包注释 AC#4 段同步（Go 侧键已落、TS 未落、控制词无页面入口）——纯注释 |

## ③ 变异牙（三发，逐发红后 `sha256sum` 两两一致还原；台件 `.scratch/wisp/probes/213/d3/`）

| 发 | 刀口 | 红句（逐字） |
|---|---|---|
| M1 | `run.go` 摘掉 `Commands: commandRosterSection213,` 行 | `run.go's PumpSources literal no longer carries the commands reader: the snapshot's commands key would silently vanish from every packet this process pushes` |
| M2 | 装配器 append 一条 `compact` 假行 | `section rows = 4, roster rows = 3: a row was dropped or invented on the wire` |
| M3 | `pump.go` 接线 `if false && p.src.Commands != nil` | `a pump with a commands reader sent no commands key`（M3 第一发因 python 不可用刀未落上，那发 `ok` **作废**；sed 重打后真红） |

## ④ 门禁读数

| 门 | 读数 |
|---|---|
| `gofumpt.exe -l`（点名 10 枚） | **空** |
| `go vet ./internal/panel ./cmd/wisp` | **rc=0** |
| `go build ./internal/panel ./cmd/wisp`（＋internal/command） | **rc=0** |
| `go test ./cmd/wisp -run 'Test213\|TestComposer' -count=1` | **ok 0.167s**（Test213×4 全 PASS） |
| `go test ./internal/panel -run TestTheCommandsKeyTravelsOnlyFromItsReader213` | **ok** |
| `go test ./internal/command -count=1` | **ok 0.392s** |
| `go test ./internal/panel -count=1`（整包） | rc=1，`--- FAIL` **恰 3 枚**：`TestApprovalCardViewJSONKeysMatchFrontendTypes`＋`TestComposerContractTypesMatchFrontend`＝**交付 3 预定红**（逐字 `Go Snapshot emits [commands] that interface PanelSnapshot does not declare`——红因落点即票面交付 3 的程序）；`TestC21DesignTokensFourWayAgree`＝**既存红**（§①.3 三条独立读数）。**其余全绿**：三枚四键字节钉、`TestAPumpWithoutAnInstructionsReaderSendsNoKey`、197 缺席/在场对、214/220/248 各族全过 |
| `sh scripts/d22scan.sh` | **rc=0 clean**（internal/ 944 Go、cmd/ 230 Go） |
| `git diff --check` | **rc=0**（LF，无空白错误） |
| 还原校验 | 三发 bak 两两 sha256 一致（6 hash ⇒ 3 对） |

## ⑤ 本腿**没做**的（具名）

1. **页面侧 TS 声明**（`frontend/src/lib/panel.ts` 的 `commands?: ...`）：没做——`frontend/**` 归机主外部 agent（Q-51），两枚契约尺红着等它，红因已写进票面与本件。
2. **页面执行命令**：没做也不该由本键做——控制词无页面路径（message 门未接），`mode` 的门本就是页面可达的 C17 既有路由；`commands` 键只递成员与门的说明。
3. **票面 AC 框**：一枚未翻（AC#1–#5 全保持 `[ ]`；终判归非实现者腿）。
4. **`TestC21DesignTokensFourWayAgree` 与 `files` 键那两枚既存红**：未修（不在本票写面，归因读数在 §①.3，留给各自主票）。
