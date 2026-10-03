# 253-r5 — 票 253 AC#2 落地腿（形ⓐ：另立一枚"全量入向名册"尺）实现件

- 性质：**写代码腿**，只做一格 = **AC#2**，只走编排者 `A558` 已裁死的**形ⓐ**。
- 写面：`internal/panel/` 下新增 `*_test.go`（正控那一小步**临时**在 `internal/panel/bridge.go` **只加一行新常量**，种→跑→还原串同一条命令）。
- ⛔ 未碰 AC 勾选框｜⛔ 未解冻任何冻结件｜⛔ `git_test.go:385`/`:517` 两枚 want-4 锚一行未动｜⛔ 未跑全仓 `go test ./...`｜⛔ 未跑 `cmd/wisp`／`internal/config`／`internal/tools`／`internal/risk`／`internal/agent/approval` 任何包｜⛔ `frontend/**`／`design/**` 零读零引｜⛔ 未 push。

---

## §0 起手锚

- 起手同发三取（一条命令内）：`date "+%Y-%m-%d %H:%M:%S%z"` = `2026-10-03 09:36:33+0800`；`git log -1 --format="%h %ci"` = `ebd5da2b 2026-10-03 09:36:20 +0800`；`git status --porcelain -- cmd internal tools docs scripts .scratch | wc -l` = `366`（他腿脏面；本腿此刻只产 `.scratch/wisp/probes/253/r5/impl.md` 一枚）。
- 分支 = `dev`。开工前另有一发同形三取 = `2026-10-03 09:30:27+0800` / `4eb228ef` / `371` 行——两发之间 HEAD 由 `4eb228ef` 前进到 `ebd5da2b`（他腿在飞），本腿所有行号一律取自**本腿自己读的时刻**，见 §1。
- `git status --porcelain -- internal/panel` = **空**（此刻 `internal/panel` 无他腿脏面；脏面在 `cmd/wisp`×3、`internal/tools`×1、`internal/config`×1）。
- 票面全名：`.scratch/wisp/issues/253-panel-inbound-has-three-ruler-holes-...-zero-audit.md`（**无 `-done` 后缀** ⇒ 未结；本腿是落地腿，不是验收腿）。
- 前件读讫：票面 §8（编排者裁定，形ⓐ）＋ `.scratch/wisp/probes/253/p4/precheck.md`（141 行）＋ `.scratch/wisp/probes/253/r3/precheck.md`（§ⓔ 四枚在册红）。

## §1 现量：两把旧尺各数什么集合 ＋ `bridge.go` 现有入向方法名全名册（行号本腿现取）

> 本节所有行号＝本腿在 `ebd5da2b`→`47b38765` 之间**自己 grep 出来的盘上状态**，⛔ 一个字没抄 p4/r3 的读数。
> 抄来的行号在本仓已经漂过两次，所以每条都在下面标了取法。

### 1a. 产码侧入向名册全貌（`internal/panel/bridge.go`，167 行）

| 枚 | 常量标识符 | 取值 | file:line（本腿现取） | 带 `panel.` 前缀？ | 旧尺数得到吗 |
|---|---|---|---|---|---|
| 1 | `MethodModeRequest` | `panel.mode.request` | `internal/panel/bridge.go:42` | 是 | 数得到 |
| 2 | `MethodWorkspaceRequest` | `panel.workspace.request` | `internal/panel/bridge.go:43` | 是 | 数得到 |
| 3 | `MethodAttachmentAdd` | `panel.attachment.add` | `internal/panel/bridge.go:44` | 是 | 数得到 |
| 4 | `MethodMessageSend` | `panel.message.send` | `internal/panel/bridge.go:45` | 是 | 数得到 |
| 5 | `MethodConfigGet` | `config.get` | `internal/panel/bridge.go:66` | **否** | **数不到** |
| 6 | `MethodConfigSet` | `config.set` | `internal/panel/bridge.go:67` | **否** | **数不到** |
| — | `ComposerRequestSource` | `panel-composer`（发信人标识，⛔ 不是方法名） | `internal/panel/bridge.go:30` | 否（无点） | 数不到，也**不该**数 |
| — | `requestIDBytes` | `8`（整型） | `internal/panel/bridge.go:76` | — | 类型即被排除 |

答侧两处（名册必须与之一致，否则就是"名册有名而 Go 不答"）：
- guard：`func knownComposerMethod(m string) bool` 在 `internal/panel/bridge.go:146`，case 标签 = 上面 1–6 六枚标识符。
- router：`func (d *ComposerDispatch) dispatch(...)` 在 `internal/panel/composer_dispatch.go:175`，
  `switch req.Method` 在 `:176`，六个 `case Method*`（`:177/:182/:187/:192/:197/:202`），
  `default:` 在 `:207` 并调 `d.rosterMismatch(req)`（`:208`）＝回填支。

### 1b. 旧尺 A：`routeLiteralRe` — `internal/panel/composer_test.go:394`

- 现取定义：`var routeLiteralRe = regexp.MustCompile(` + "`" + `"panel\.[A-Za-z.]+"` + "`" + `)` ⇒ **前缀写死在正则里**。
- 它数的集合：对 `frontend/src` 树里每一行 renderer 源，抽出所有 `"panel.*"` 字面（使用点 `composer_test.go:475-479`），
  与闭集 `composerRouteLiterals()`（`composer_test.go:409-419`）对账，凡不在闭集者进 `rep.unknown`，红点 `composer_test.go:521`。
- 分母形状＝**"页面写了什么 `panel.` 字面"**，⛔ 不是"Go 答什么名"⇒ `config.get`/`config.set` 这类无前缀名**根本不进分母**。
- ⚠ 本腿**没有读 `frontend/**`**（两层禁令），也⛔ 没有把 `composer_test.go` 那套扫描接进新尺；
  上面两句是对该测试**自身源文件**（`internal/panel/composer_test.go`，属本包产码外的测试具器）做的射程判断，非内容引用。

### 1c. 旧尺 B：`panelMethodRe` — `internal/panel/git_test.go:394`（⛔ 不在 `cmd/wisp`）

- 现取定义：`panelMethodRe = regexp.MustCompile(` + "`" + `"(panel\.[a-z0-9_.-]+)"` + "`" + `)`（同一毛病：**前缀写死**）。
- 抽取器 `whitelistMethodsFromSource(t, path)` 在 `internal/panel/git_test.go:407-423`：**只读它被点名的那一个文件**。
- 两枚 want-4 锚（本腿一行未动）：
  - `internal/panel/git_test.go:385`（用例 `TestGitDimensionHasNoModelCallableTool`，`wantMethods` 在 `:381-383` 列四枚 `Method*`）；
  - `internal/panel/git_test.go:517`（用例 `TestPlantedGitToolShapesGoRed` 里的"clean half"，`len(real) != 4`；
    同用例 `:503-514` 那发 `bridge_five.go` 五枚 `panel.*` 的自包含正控数到 5，证尺有牙）。
- 它数的集合＝**"bridge.go 里带 `panel.` 前缀的字面枚数"**⇒ 现量 4，`config.*` 两枚天然不在内。
- 旧尺 B 的文件筛跳 `_test.go`：`internal/panel/git_test.go:464` 的 `gitToolNamesUnder`
  （`!strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go")` 即 `return nil`）⇒ **本腿新增的 `_test.go` 不进它的分母**。

### 1d. 冻结件那把尺也数不到"未登记的名"（射程判断，非内容引用；⛔ 一字未改）

- `internal/panel/l2_grant_boundary_test.go:369` 的产码枚举同样跳 `_test.go`。
- **本腿新量到的一条射程事实（p4 未记）**：那把冻结尺的 `answered` 集**只从 `knownComposerMethod` 一处导出**——
  `internal/panel/l2_grant_boundary_test.go:777` 显式 `if decl.Name.Name != "knownComposerMethod" { return true }`
  ⇒ `internal/panel/composer_dispatch.go:175` 那枚 **router 的 case 列表不在它的分母里**。本尺的分母 4（ROUTER AST）正是补这一格。
- `guardRosterOf`（`:2030`）把 `bridge.go` 里**标识符以 `Method` 开头**的常量与 guard 的 case 列表**双向**对账
  （缺那一句的红在 `:2125`），另有 `routeShapedName`（`:1050`，点分形状、不认前缀）、
  `routeNamePool`（`:1097`）、`poolJudgedByRealGuard`（`:1125`，只问"guard 答了而枚举没看见"那一向），
  以及 `:1294` 的"Go 答的名必须有 `Method*` 常量声明"。
- ⇒ 缺的那一格正是 AC#2：**"bridge.go 里多出一枚入向方法名常量而没有任何名册尺登记它"**——
  若那枚新常量既不带 `panel.` 前缀、标识符又不以 `Method` 开头、也还没接进 guard 的答侧，上面四把尺**全部无话**。
  本腿 §4 的第二发 plant（`Method` 名那一式）会同时被 `:2125` 数到，所以正控用**真树只加、副本两种式**分开做，逐条读数在 §4。

## §2 新尺的实现，与它数的集合定义

- 落点：**唯一一枚新文件 `internal/panel/inbound_roster_253_test.go`**（`package panel`，md5 `9d865c9a3e6203d203afd1e963733c15`，
  834 行；`git add` 于 commit `87bc6aca`）。⛔ 未改 `bridge.go`／`composer_dispatch.go`／任何旧尺／任何冻结件一字。
- 名册本体＝**写死在测试里的字符串字面**（不是引用 Go 常量）：`wantFullInboundRoster253` 六名，
  外加两枚写死的数 `wantFullInboundRosterSize253 = 6`、`wantNonPanelPrefixedInbound253 = 2`。
  取字面而⛔ 不取常量，理由与 `git_test.go:403-406` 自述同形：**重列一遍被审的常量＝那把尺永远绿**。
- **它数的集合**（五个互相独立的分母，每个都**双向**对账，见 `sameNameSet` 同时返回 `missing` 与 `extra`）：
  1. **DECLARED**＝`bridge.go` 里每一枚包级字符串常量，命中条件取**并集**：
     标识符以 `Method` 开头（`methodNamedIdent`）**或** 值是点分名（`inboundDotShaped`）。
     ⇒ 第二个子句不认任何命名空间：`config.peek`／`settings.read`／`x.y` 一样在射程内；
     `panel-composer`（无点、标识符非 `Method`）自动被排除，⛔ 不需要我为它写白名单。
  2. **RUNNING GUARD**＝名册每名问一次**真的** `knownComposerMethod`（本测试与产码同包 ⇒ 直接调，不复刻判据）。
  3. **GUARD AST**＝`bridge.go:146` 那枚 switch 的 case 标签，逐个解析到该文件自己的常量取值；
     裸字面标签解析成值、非字面/解析不出者进 `problems`；⛔ 该函数不许出现 `default:`（出现了就是红）。
  4. **ROUTER AST**＝`composer_dispatch.go:175` 的 dispatch case 标签（先 `bridge.go` 常量表、再全包常量表解析），
     其 `default:` 被**要求存在且必须调 `rosterMismatch`**⇒ 回填支被摘掉这一形也在射程内（`callsNamed`）。
  5. **ANSWERED 字面**＝全包（跳 `_test.go`）任何点分字符串**字面**，只要**运行中的 guard 答它**而名册没有 ⇒ 红。
     专治"名字只写在表达式里、从未成为可登记的常量"那一形。
- ⛔ **没有再写一条"只认某种前缀"的正则**：上面五个分母里唯一的形状判据是"值是不是点分名"
  （`inboundDotShaped`，问**有没有命名空间**、⛔ 不问**哪个**命名空间），另两条入口判据是标识符前缀 `Method`
  与"guard 答不答"，三者取并集⇒ 拿掉任意一条仍剩下覆盖，写死成一条前缀正则那种"把缺陷复制一遍"的形状在这里不存在。
- 五枚用例：`TestFullInboundMethodRosterIsClosed`（在册尺）、`TestPlantedUnregisteredInboundMethodNameGoesRed`（正控，三式 plant）、
  `TestRegistrationSilencesTheRosterRuler`（"补登记"两式读数）、`TestFullInboundRosterRefusesAnEmptyDeclarationFile`（空集反恒真）、
  `TestSubsetOnlyRosterDirectionIsBlindToAPlantedName`（把**反形**当场量瞎）。
- plant 全部落在 `t.TempDir()` 副本；真树在正控里**只读**（"clean half must stay quiet"那半用同一条 helper 复读真树）。

## §3 门禁读数（build ／ gofumpt ／ `internal/panel -count=1` 终态 ／ d22scan）

> 判绿只认 `--- FAIL`／`--- PASS` 行；每处注明取自哪枚日志（原始件已随 commit `87bc6aca` 落盘同目录）。

| 门禁 | 命令 | 读数 |
|---|---|---|
| 构建 | `GOFLAGS= go build ./...` | 零输出、exit 0（还原后再跑一遍仍 0；本腿全程未跑别的包的测试） |
| 格式 | `"$(go env GOPATH)/bin/gofumpt.exe" -l internal/panel` | **空列表** ⇒ 我新增的 `internal/panel/inbound_roster_253_test.go` 不在名单里 |
| 本包测试（起手基线） | `GOFLAGS= go test ./internal/panel/ -count=1 -v` → `baseline-run.log` | exit 1；`--- PASS` 110 枚／`--- FAIL` 4 枚／`--- SKIP` 0 枚（此刻我的 4 枚用例已在 110 内） |
| 本包测试（**终态**） | 同上 → `final-run.log`（bridge.go 与本测试件都还原之后） | exit 1；`--- PASS` **111**／`--- FAIL` **4**／`--- SKIP` **0** |
| 终态那 4 枚 FAIL | `final-run.log` | `TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`——**与起手基线逐名相同**（票 253-r3 §ⓔ 在册页面契约族常红，台账 A533），⛔ 不算本腿凭据、本腿未修、⛔ 一枚未放宽 |
| CI 形状扫 | `./tools/d22scan/d22scan.exe` | exit **0**；判语行 `d22scan: clean - no D22 ban violations`；其 ban #8 自述 scope＝`internal/ 489 Go files, comments and _test.go included` ⇒ **本腿新增那枚 `_test.go` 确实在它射程里且净** |
| d22scan 的 SKIP 支 | 同上页脚 | `skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]`——该工具对**已构建产物目录**的既有页脚（不是 `--- SKIP`、不是本腿引入）；判语行仍 clean；⛔ 本腿没把它读成通过、也没读成红 |

`git status --porcelain -- internal/panel` 终态＝只剩 `?? internal/panel/inbound_roster_253_test.go`（本腿唯一产码件）⇒ **无任何既跟踪文件被本腿改动**。

## §4 正控 ＋ 反控 ＋ 恒真自查 三发自证表

### 4-1 正控（票面 AC#2 那句判据的本体：种"新增方法名没进名册"⇒ 指名那一步必须红）

- **种什么形**：`internal/panel/bridge.go` **末尾只加两行**（⛔ 不改不删任何既有行）——
  `// 253-r5 AC#2 temporary plant: an inbound method name that no ruler registers.` ＋ `const ConfigPeek253 = "config.peek"`。
  选形理由：值**不带 `panel.` 前缀**、标识符**不以 `Method` 开头**、guard **不答**它 ⇒ 两枚 want-4 锚（`git_test.go:385`/`:517`）
  与冻结件 `l2_grant_boundary_test.go`（`guardRosterOf` 只数 `Method` 开头标识符，红句在 `:2125`）**全部不该被牵连**；
  编排者那条硬约束逐字遵守。种→跑→还原**串在同一条命令里、中间未停顿**（先用 `B=$(git rev-parse HEAD:internal/panel/bridge.go)` 固定 blob）。
- **哪枚必须红**：`TestFullInboundMethodRosterIsClosed`。**红句逐字**（`plant-run.log`）：

```
    inbound_roster_253_test.go:437: ticket 253 AC#2: bridge.go declares the inbound method name config.peek at bridge.go:170 const ConfigPeek253 and the closed roster does not carry it - a new inbound route entered the tree without being registered anywhere. The two prefix rulers cannot see a name that carries no "panel." prefix (panelMethodRe at git_test.go:394, routeLiteralRe at composer_test.go:394), which is the hole this roster exists to close. If the name is meant to stay, register it in wantFullInboundRoster253 and move the count constants beside it
    inbound_roster_253_test.go:519: full inbound roster held: 7 names declared in bridge.go, 6 guard case labels, 6 router case labels, 2 of them carrying no "panel." prefix ("config.get", "config.set")
--- FAIL: TestFullInboundMethodRosterIsClosed (0.01s)
```

  ⇒ 红句**指名了那一枚常量与它落在哪一行**（`bridge.go:170 const ConfigPeek253`），不是只报"集合不等"。
- **同刻我自己的另两枚也红了（过杀伤如实登记，不粉饰）**：那两枚用例把"真树必须安静"当作自己 plant 的诚实半，所以种在真树上必然同红——

```
    inbound_roster_253_test.go:642: planted an unregistered inbound name (plain-identifier-const) and the roster ruler reported extra="config.peek", "config.peek253", want exactly ["config.peek253"] - this ruler has no teeth for the shape AC#2 was filed over
    inbound_roster_253_test.go:642: planted an unregistered inbound name (Method-named-const) and the roster ruler reported extra="config.foo253", "config.peek", want exactly ["config.foo253"] - this ruler has no teeth for the shape AC#2 was filed over
    inbound_roster_253_test.go:683: the plants only mean something if the real tree reads quiet under the same helper; got missing=(none) extra="config.peek"
```

  整包读数：`plant-run.log` ＝ `--- PASS` 107／`--- FAIL` 7（＝4 枚在册常红 ＋ 本腿 3 枚自红）。
- **别人的锚在那一刻的状态**（同一份 `plant-run.log`，⛔ 一枚都没被我的种形打红）：
  `TestGitDimensionHasNoModelCallableTool` **PASS**｜`TestPlantedGitToolShapesGoRed` **PASS**｜
  `TestAnsweredPanelRoutesCarryNoApprovalDecision` **PASS**｜`TestRealGuardRefusesEveryAssemblableApprovalRouteName` **PASS**｜
  `TestNoInboundEnvelopeCanBindAnApprovalVerdict` **PASS**｜`TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes` **PASS**｜
  `TestTheRendererHoldsExactlyOneDoorToTheHost` **PASS**。
- **还原与复跑终值**：`git cat-file blob $B > internal/panel/bridge.go` ⇒
  种前 md5 `63e69a79422b577cca6a320c52915261`／167 行，还原后 **md5 逐字节相同**／167 行／`git diff --stat -- internal/panel/bridge.go` **空**。
  复跑（`final-run.log`）⇒ 本腿五枚全 `--- PASS`，整包 111 PASS／4 FAIL（那 4 枚＝在册常红）。

### 4-2 反控（⛔ 没让 bridge.go 真多出第五枚 `panel.*`；want-4 两锚未被牵连）

- `bridge.go` 的 `panel.*` 字面枚数在还原后仍是被两锚认下的那四名：两枚锚在 `final-run.log` 均 **PASS**，
  其中 `:517` 那枚的断言就是 `len(whitelistMethodsFromSource(real bridge.go)) != 4` ⇒ 它绿＝真树仍是四名。
  ⚠ 一处易误读具名登记：`grep -c '"panel\.' internal/panel/bridge.go` 给 **5**，因为 `bridge.go:35` 的**文档注释**里有一句 `"panel.*"`；
  尺用的正则 `"(panel\.[a-z0-9_.-]+)"` 匹配不到它（`*` 不在字符类且要求收尾双引号）⇒ 尺数到的是 4。⛔ 别拿那个 grep 计数当 want-4 的读数。
- 本腿⛔ 未往 C17 那四名里加真方法，⛔ 未改 `panelMethodRe`／`routeLiteralRe`／`composerRouteLiterals()`／`wantMethods` 一字。
- **want-4 两枚锚的用例名与终态读数**（编排者点名要的两行）：
  - `internal/panel/git_test.go:385` ∈ 用例 **`TestGitDimensionHasNoModelCallableTool`** ⇒ 终态 `--- PASS: TestGitDimensionHasNoModelCallableTool (0.00s)`（`final-run.log`；起手基线同样 PASS）。
  - `internal/panel/git_test.go:517` ∈ 用例 **`TestPlantedGitToolShapesGoRed`** ⇒ 终态 `--- PASS: TestPlantedGitToolShapesGoRed (0.01s)`（`final-run.log`；起手基线同样 PASS）。
- 冻结件 `internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`：⛔ 一字未改
  （只读其抽取器与"跳 `_test.go`"那一支做**射程判断**，非内容引用）；l2 那族用例在 `final-run.log` 全 PASS，
  唯一在册常红的 `TestC21DesignTokensFourWayAgree` 属页面契约族（`tokens_fourway_test.go` 名下），与本腿零交集、本腿未碰。

### 4-3 恒真自查（把断言换成反形，它还绿不绿？＝**绿** ⇒ 该反形不可作凭据）

- **第一发（真种形＋真反形，实测）**：临时把 `TestFullInboundMethodRosterIsClosed` 的分母 1 换成反形＝
  ① `len(declared) == 0` 时 `t.Logf` + `return`（票面举例的"名单为空也算过"）、② 只报 `missing`、把 `extra` 那支改成 `_ = extra`（subset-only）。
  同一条命令里再种 `ConfigPeek253 = "config.peek"` → 跑 → 还原 `bridge.go`（blob）与本测试件（同目录快照 `inbound_roster_253_test.shipped-copy.txt` 回 cp）。
  **读数＝绿**（`reverse-form-run.log`）：

```
    inbound_roster_253_test.go:521: full inbound roster held: 7 names declared in bridge.go, 6 guard case labels, 6 router case labels, 2 of them carrying no "panel." prefix ("config.get", "config.set")
--- PASS: TestFullInboundMethodRosterIsClosed (0.02s)
```

  ⇒ 同一棵**种了假腿的树**（它自己的 `t.Logf` 就写着 `7 names declared in bridge.go`）在反形下**全绿＝这把尺对这件事不敏感**，
  所以**反形不可作凭据**；AC#2 的凭据只能是 4-1 那发正形的红。还原后两文件 md5 回到 `63e69a79…`／`9d865c9a…`。
- **第二发（把反形做成在册测量，随代码长存）**：`TestSubsetOnlyRosterDirectionIsBlindToAPlantedName` 在 `t.TempDir` 副本上
  当场量出 subset-only 无话可说、而正形指名——终态 `--- PASS`，日志逐字：

```
    inbound_roster_253_test.go:832: relaxed form measured blind on purpose: with the plant in the source the subset-only reading finds no problem (missing=(none)) while the shipped equality names "config.blind253". Subset-only is not admissible as AC#2's credential
```

- **第三发（空集那一支）**：`TestFullInboundRosterRefusesAnEmptyDeclarationFile` 把尺指向一枚常量都没有的 decoy，必须报"6 名全缺"才算有牙；
  终态 `--- PASS`，日志逐字：

```
    inbound_roster_253_test.go:800: empty-source refusal held: 6 roster names reported missing against a file that declares none
```

  ⇒ 在册那支（`if len(declared) == 0 { t.Fatalf(...) }`）⛔ 不接受"名单为空也算过"——这一条不是推测，是被上面那发量出来的。

### 4-4 "补登记后不响"的终值（题面要的那一半，两式都给读数）

- **只补名册（真树种形＋临时改我自己的名册与两枚计数）＝仍响**（`registered-run.log`；DECLARED 那支安静了、437 那发不再出现，红在答侧三支）：

```
    inbound_roster_253_test.go:466: ticket 253 AC#2: the roster carries config.peek but the running guard knownComposerMethod refuses it - a rostered inbound name Go does not answer is a page-visible dead door
    inbound_roster_253_test.go:480: ticket 253 AC#2: the roster carries config.peek but knownComposerMethod's case list no longer names it - the guard stopped answering a door this file still claims is open
    inbound_roster_253_test.go:504: ticket 253 AC#2: the roster carries config.peek but composer_dispatch.go's dispatch has no case for it - the guard would answer it and the router would fall through to the roster-mismatch backstop
--- FAIL: TestFullInboundMethodRosterIsClosed (0.02s)
```

- **完整登记（声明＋guard 答它＋名册与两枚计数一起动）＝不响**，在 `t.TempDir` 副本里量（`TestRegistrationSilencesTheRosterRuler`，终态 `--- PASS`）：

```
    inbound_roster_253_test.go:761: registration read both ways: declared AND answered reads quiet (7 names, no diff); roster-only reads red on the answered side ("config.registered253"), so adding a name to the list is not how this ruler gets satisfied
```

- ⇒ "补登记后不响"的**终值＝`missing=(none) extra=(none)`（7 名全对上）**，条件是登记走满**声明＋答侧＋名册**三处；
  "只改测试里的名册"那一式我**故意让它继续红**——名册宣称每名都是活门，Go 不答的名册项就是页面点了没反应的死门（票 253 §现量-3 那族症状）。
  真树上做完整登记要动 `bridge.go:148` 那行 guard＝**契约面**、⛔ 不属本腿射程（票面 §8 已判形ⓐ不解冻任何枚），故只在副本量这一式，并在 §5-1 具名交代。
- 在册终态读数（`final-run.log`，逐字）：

```
    inbound_roster_253_test.go:519: full inbound roster held: 6 names declared in bridge.go, 6 guard case labels, 6 router case labels, 2 of them carrying no "panel." prefix ("config.get", "config.set")
--- PASS: TestFullInboundMethodRosterIsClosed (0.01s)
```

  ⇒ `config.get`／`config.set` **确实被这把尺数到**（AC#2 题面要求），两枚旧尺仍各守原射程、一字未动。

## §5 我可能写错的条目（自我对抗）

1. **行号归因**：§4 里引的红句带 `file:line`，其中 `inbound_roster_253_test.go:437 / :519 / :642 / :683 / :761 / :800 / :832`
   与**终态文件**行号一致；而 `:466 / :480 / :504`（§4-4 只补名册那一发）与 `:521`（§4-3 反形那一发）**属当时那两份被临时改过的文件**
   （名册多一名＝整体下移一行；反形替换＝上下文不同）。⛔ 拿这两组去核终态文件会对不上——终态里同一条断言在 `:465 / :479 / :503`。
2. **"补登记后不响"的终值来自 `t.TempDir` 副本，⛔ 不是真树**：真树上做完整登记要改 `bridge.go:148` 那行 guard＝契约面、须人工批准、本腿不许。
   若验收指望的是"真树上只补名册就该全绿"，我这把尺**做不到**，而且我判断它**不该**做到（名册宣称每名都是活门）。
   这一条是选形判断，归裁决者，不在本腿自证射程内——如实标出。
3. **两套点分形状判据并存**：本尺的 `inboundDotShaped` 与冻结尺的 `routeShapedName`（`l2_grant_boundary_test.go:1050`）语义近乎同形而各写一遍。
   我没复用它的函数（那会把 l2 的语义搅进本尺、也可能反过来动到它的分母），但这是**维护债**：将来一名改动、另一处可能静默过期。
4. **`ComposerRequestSource` 靠"值不含点"被排除**（`bridge.go:30` 的值是 `panel-composer`）。若有人把它改成带点的值
   （例如 `panel.composer.source`），它会被分母 1 数成"入向方法名"⇒ 尺红，而那枚红是**误伤**（它是发信人标识，不是门）。
   我今天⛔ 没种这一形、没量过这一发，属未测风险，具名登记（修法要么显式排除该标识符，要么由改名人重新考虑值）。
5. **分母 5 只扫 `internal/panel` 的产码**：一枚声明在别的包、由别的包自己答的入向名，本尺无话 ⇒ 见 §6-3。
   我对 `cmd/wisp/config_reload.go:86` 那枚 `config.reload` "属工具名不属 composer 方法名"的判断**只据那一行与其所在文件角色**，
   ⛔ 未读 cmd/wisp 的派发链（那包有 `255-r2` 在写）⇒ 这一句是推断，标〔半量〕。
6. **router 的 `default:` 我要求它必须调 `rosterMismatch`**：若后续合法重构把那枚回填函数改名或内联，本尺会红，
   而那红可能是**我把实现细节钉成契约**（过拟合），不是缺陷。
7. **plant C 用精确整行匹配**（`strings.Replace` 那行 guard）：那行一旦换行/重排，plant 会"没换到任何东西"。
   我用 `t.Fatalf` 挡住了这种情况，但**挡住的是这一发正控**，⛔ 不是尺的在册绿态——验收若看到那枚 fatalf，别把它读成尺失能。
8. **三枚同刻自红**（§4-1）：我把"真树必须安静"塞进了 plant 用例与 blindness 用例自己的断言里，所以真树种形时它们必然同红。
   按"只有指名那一枚该红"的读法，我的尺呈**过杀伤**形状。可改（把那半句并入在册尺独占）——本腿没改，
   因为它同时是"clean half must stay quiet"这条仓内定式的现成证据。
9. **没跑 `go vet`**（不在给我的门禁清单里）⇒ 未使用符号那一类只被 `go build ./...` 与测试编译间接覆盖。
10. **基线含我自己的文件**：`baseline-run.log`（110 PASS）已包含我新写的 4 枚用例 ⇒ "落地前整包长什么样"这一起手零文件基线我没取
    （取它要把自己的件移出树，⛔ 仓内不删不移）。归因靠"那 4 枚 FAIL 逐名等于 253-r3 §ⓔ 在册常红"，⛔ 不靠"我落地前后差集"。

## §6 判不动的地方

1. 〔量不到〕`frontend/**` 里今天有没有 `config.get`／`config.set` 的字面调用点（两层禁令，本腿零读零引）
   ⇒ AC#2 本腿只做到"Go 侧六名被一把尺数到＋未登记的形必红"，页面侧那半我判不了；形ⓐ 的成立不依赖这一读数。
2. 〔量不到〕真树上完整登记第 7 枚入向方法后的整包读数——要改 `bridge.go:148`＝契约面（票面 §8 已裁⛔ 不解冻），本腿不许做 ⇒ 只有副本读数。
3. 〔半量〕`cmd/wisp/config_reload.go:86` 的 `configReloadTool = "config.reload"` 算不算"入向方法名"：判它要读 cmd/wisp 派发链，
   而那包此刻由 `255-r2` 在写。本尺⛔ 没把它算进名册（分母边界划在 `internal/panel`）⇒ 若它其实是入向名，
   这是"这张名册只管本包入向门"的直接后果，具名交回裁决者，⛔ 本腿不替它判。
4. 〔量不到〕两枚 want-4 锚**该不该**从 4 扩到 6（把 config 两枚纳入 C17 那把尺）＝契约判断。本腿只复量其现状（§1c、§4-2）。
5. 〔量不到〕运行时那一支：页面真发一枚未登记方法名会发生什么。本尺＝静态名册 ＋ 布尔 oracle，端到端归 cmd/wisp 那族能力形用例（票 253 AC#1 写面）。
6. 〔半量〕`tools/d22scan` 对"测试尺的射程"是否有话：只跑了现成 exe（exit 0 clean，ban #8 自述含 `_test.go`、489 枚 Go 文件），
   ⛔ 未读 `tools/d22scan/main.go` 的判据内部 ⇒ 与 p4 §4-2 同判，登记为半量。
7. 〔量不到〕交件之后别的腿往 `bridge.go` 的 const 块加名时会不会真撞上这把尺：只补产码不补名册 ⇒ 指名红；三处一并补 ⇒ 不响。
   这一"将来是否真响"我今天量不到（今天树是干净的）。

## §7 交件判语

- **AC 勾选框一枚未碰**：票 253 的 AC#1／AC#2／AC#3／AC#4 四枚 `- [ ]` 框与票面 §8 全部原样，⛔ 本腿未改 `.scratch/wisp/issues/253-*.md` 一字
  （勾选归编排者凭非实现者验收表，本腿交完他会另派验收）。本腿也**不自判"AC#2 已完成"**。
- **未 push**：只 commit、⛔ 全程零 push；每次 commit 显式 pathspec 且 `add` 与 `commit` 串同一条命令、中间未停顿；
  ⛔ 未用 `git add -A`／`.`／`-a`、`commit -a`、`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／merge／worktree。
  本腿 commit 链：`32af81ac`（§0 锚＋骨架）→ `87bc6aca`（新尺＋五枚原始日志＋快照＋msg）→ 本件 §1–§7 那次 commit。
- **want-4 两枚锚一行未动**：`internal/panel/git_test.go:385`（用例 `TestGitDimensionHasNoModelCallableTool`）与
  `internal/panel/git_test.go:517`（用例 `TestPlantedGitToolShapesGoRed`）⛔ 一字未改，终态读数均 `--- PASS`；
  我也⛔ 没往 `bridge.go` 真加第五枚 `panel.*` 方法（§4-2），正控只种**不带前缀**那一形（编排者硬约束逐字遵守）。
- **零产码改动**：`internal/panel/bridge.go` 与 HEAD blob 逐字节相同（md5 `63e69a79422b577cca6a320c52915261`、167 行、`git diff --stat` 空）；
  写面唯一新增文件＝`internal/panel/inbound_roster_253_test.go`；两枚临时改动（真树种常量／反形与名册替换）都在**同一条命令内**还原并逐字节核对。
- **冻结件零触碰**：`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go` 一字未改；
  从其中引出的每一句都在 §1d/§4 标了"射程判断，非内容引用"。
- **未解冻任何契约件**：`D1–D47`／`C1–C32`／`R1–R9`／D43 转移表／`docs/PLAN.md`／`docs/specs/**`／`docs/BUILD.md`／`docs/SLO.md`／
  `internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`internal/perm/ticket90_persist_test.go`／`.github/workflows/ci.yml`
  全部⛔ 零改动；⛔ 未放宽任何既有断言、⛔ 未 `t.Skip`、未把 SKIP 读成通过（整包 `--- SKIP` ＝ 0 枚）。
- **同机互斥遵守**：只跑 `go test ./internal/panel/ -count=1`（外加 `go build ./...`、gofumpt、d22scan）；
  ⛔ 未跑全仓 `go test ./...`，⛔ 未跑 `cmd/wisp`／`internal/config`／`internal/tools`／`internal/risk`／`internal/agent/approval` 任何包；
  未设 sherpa PATH（本腿不需要）。`frontend/**`／`design/**` 零读零引，证据件与回报正文内不含那两目录的任何内容或转述。
- **凭据卫生**：全程未读、未引、未落任何 API key／凭据**值**；本件只出现常量标识符、方法名与文件行号。
- **临时件只建不删**：`.scratch/wisp/probes/253/r5/` 下新增 `impl.md`、`msg-skeleton.txt`、`msg-code.txt`、
  五枚原始日志（`baseline-/plant-/reverse-form-/registered-/final-run.log`）与一份测试件快照 `inbound_roster_253_test.shipped-copy.txt`；
  ⛔ 仓内零删除、零改名。
- **交件时刻**（同发三取）：`date "+%Y-%m-%d %H:%M:%S%z"` = `2026-10-03 10:04:54+0800`；
  `git log -1 --format="%h %ci"` = `82741b98 2026-10-03 09:58:28 +0800`（⚠ 此刻 HEAD 不是我那枚——`32af81ac` 与 `87bc6aca`
  都在其上，其间他腿又进了若干枚 commit，本件所有行号/读数对应的是我**取数那一刻**的盘上状态）；
  `git status --porcelain -- cmd internal tools docs scripts .scratch | wc -l` = `372`（他腿脏面）。
  此刻 `git status --porcelain -- internal/panel` = **空** ⇒ 本腿在 internal/panel 的写面已全部落到 commit `87bc6aca`，工作树无残留。

### 收笔锚

- 本件填齐后同发三取：`date "+%Y-%m-%d %H:%M:%S%z"` = `2026-10-03 10:04:54+0800`；HEAD = `82741b98`；porcelain = `372` 行。
- 本件行数（收笔前）＝`298`；占位符清零：全文无「（取数中）」残留。
- 本腿 commit 链：`32af81ac`（§0 锚＋骨架）→ `87bc6aca`（`internal/panel/inbound_roster_253_test.go` ＋ 五枚原始日志 ＋ 测试件快照 ＋ msg）
  → 收笔 commit（本件 §1–§7 填齐）。只 commit、⛔ 未 push。
- 一句话结论给裁决者：**新尺数的是"写下来的六名封闭入向名册"**（`panel.*` 四名 ＋ `config.*` 两名 ＋ 任何其它点分形状都在射程内），
  正控红句逐字在 §4-1，"补登记后不响"的终值在 §4-4（副本＝0 差集；只补名册＝仍红，理由与归属写在 §5-2），
  恒真自查那一发**是绿的**（§4-3 反形在种了假腿的树上 `--- PASS` ⇒ 反形不可作凭据，凭据只能是正形那发红），
  want-4 两枚锚终态 `--- PASS: TestGitDimensionHasNoModelCallableTool (0.00s)` 与 `--- PASS: TestPlantedGitToolShapesGoRed (0.01s)`。
  ⛔ 本腿不自判 AC#2 完成——那三枚框归编排者凭非实现者验收表翻。
