# 167-a5 — 撞钉预检（只读普查腿）：`cmd/wisp`＋`internal/panel` 今天所有"计数型"钉名册，逐枚裁"落地腿 167-r2 一旦新增字段／新增那一跳，哪些绿用例变红"

> 运行类型＝**只读**：⛔ 零 `go` 命令（`go test`／`go build`／`go vet`／`go env` 一律未跑，因写腿 `268-r1` 正在 `cmd/wisp` 取整包终态）。
> 本件的"红/不红"一律是**射程判断**（读尺体本体得出），颜色属〔预测〕，逐枚归口给编排者跑。
> 前一腿 `167-a3` §4 已做过**占用那一枚**的撞钉射程；本腿**独立复测**（⛔ 不抄它的读数），且把面从"占用"扩到**四枚输出全量＋`cmd/wisp` 侧**。
> 派单里每一句行号都当**待验断言**处理：复认成立才写，漂了具名报回。

## 0. 起手锚（逐字读数，同一发命令取）

| 项 | 读数 |
|---|---|
| 时刻 | `2026-10-05 14:00:02 +0800` |
| `git rev-parse --short HEAD` | `94380ab3` |
| 分支 | `dev` |
| 五根 porcelain（`git status --porcelain -- cmd internal tools scripts docs \| wc -l`） | **1 行**，逐字 ` M cmd/wisp/resident_approval_windows.go` ⇒ **别人的脏，本腿只报不改**（本腿对该根贡献 0 行；见 §7 自证） |
| 尺面枚数（起手即量） | `cmd/wisp` **99** 枚 `.go`／其中 `*_test.go` **65** 枚；`internal/panel` **34** 枚 `.go`／其中 `*_test.go` **20** 枚 |
| 本腿写面 | 仅 `.scratch/wisp/probes/167/a5/census.md`（＋同目录临时 msg 件）；⛔ 未动任何产码／测试／工单／台账；⛔ 未碰任何票面 `- [ ]` 框 |
| 阅读禁令自证 | grep/find 根一律显式写 `cmd internal tools scripts docs .scratch` 或其子集；⛔ 从未拿 `.` 当根；`frontend/**` 与 `design/**` **未读一字**（本件凡涉及页面侧只引 Go 尺体的字面路径，不引内容） |

## 1. 问一（本腿主体）— 今天 `cmd/wisp` 与 `internal/panel` 所有"计数型"钉

> 尺（八把，14:0x–16:26 之间现跑，⛔ 零 `go`；单笔时刻见每格）：
> S00 `grep -rn "want \|len(\|count" cmd/wisp internal/panel --include=*_test.go | wc -l` ＝ **1675**（派单给的起手尺，复认成立；**它只是噪底**，逐枚读原文才有点名价值）
> S03 `grep -rnE "len\([^)]*\)[[:space:]]*(!=|==)[[:space:]]*[0-9]+" … --include=*_test.go | wc -l` ＝ **259**
> S05 `grep -rn "sortedCopy(\|jsonKeysOf\|tsInterfaceKeys\|reflectionTopKeys" … | wc -l` ＝ **15**
> S09 `grep -rn "reflect\.TypeOf" … | wc -l` ＝ **28**（原文逐枚过目）
> S13 `grep -rnE '!= "[a-zA-Z_]+(,[a-zA-Z_]+)+"' … | wc -l` ＝ **4**
> S14 `grep -rnE "equalStrings|sameNameSet|sortedCopy" … | wc -l` ＝ **28**
> S15b `grep -rn "len(" … | grep -iE "keys|fields|methods|roster|entries|sections|slots|names" | wc -l` ＝ **127**（本族的名册就是从这 127 行逐枚过目筛出来的）
> S30 `grep -rn "t.Skipf(\|t.Skip(" cmd/wisp internal/panel --include=*_test.go | wc -l` ＝ **16**（§5 用）
>
> **筛法**：一枚 `len(` 只有当它数的是**"落地腿会往里加东西的那类集合"**（快照 JSON 键／结构体字段／方法名／文件枚数／名册行）才算本问的钉；数"运行期条数"（如 `len(snap.Pending)` 判队列行为）归问二。

### 1.1 族①｜快照**顶层键集**钉 ＝ **4 枚**（★比 `167-a3` R-⑨ 的"3 枚"多一枚，见 §7 第 1 条）

三枚共用抽取规则"`json.Marshal` 出字节 → `json.Unmarshal` 进 `map[string]json.RawMessage` → 收键 → `strings.Join(sortedCopy(…))` 与逐字常量比"，第四枚**同一台机器、但换成数数**。

| 枚 | 锚 | 断言原句（逐字） | 新增一枚出向字段会不会红 |
|---|---|---|---|
| K1 | `internal/panel/pump_test.go:123`（函数 `TestThePumpBuildsThePacketFromWhatTheHostHolds`；注释 `:109-110` 逐字 `// And the bytes: the four keys are the contract, so a fifth arriving here is`／`// a contract change dressed as a bug fix.`） | `	if strings.Join(sortedCopy(got), ",") != "composer,generatedAt,pending,results" {` → 红句 `:124` 逐字 `t.Errorf("snapshot JSON keys = %v, want exactly the four PanelSnapshot declares", got)` | **只有"新增顶层键且这一发它就有值"才红**：该 fixture 的 `PumpSources` 不设 `Instructions`／`Tasks` reader ⇒ 两枚既有 pointer＋omitempty 节在此缺席。占用/序号若走 **pointer＋omitempty 的新节** ⇒ 不红；走**值类型节** ⇒ 必红 |
| K2 | `internal/panel/pump_test.go:291`（`TestPublishHandsTheBytesToTheAttachedExit`） | `	if strings.Join(sortedCopy(keys), ",") != "composer,generatedAt,pending,results" {` | 同 K1（fixture 只设 `Mode`／`Now`／`Out`） |
| K3 | `internal/panel/subagent_roster_197_test.go:214`（`TestAPumpWithoutARosterReaderSendsFourKeys`） | `	if got := strings.Join(sortedCopy(keys), ","); got != "composer,generatedAt,pending,results" {` → 红句 `:215` 逐字 `t.Errorf("keys without the reader = %q, want the four", got)`；**另有一枚反向钉** `:207-208` 逐字 `	if _, ok := wire["tasks"]; ok {`→`t.Fatalf("a pump assembled without a roster reader sent the key anyway: %s", data)` | 同 K1。**这枚是 `internal/panel/pump.go` 自述里没数到的一枚**（本腿不引自述当名册） |
| ★K4 | `internal/panel/subagent_stream_197_test.go:129-130`（同一发里，紧接 `:118-121` 只设 `Results` 一台 reader 的 pump） | `	if len(generic) != 4 {` → 红句逐字 `t.Errorf("snapshot top-level keys = %d, want the four PanelSnapshot declares: %v", len(generic), generic)` | **这是 a3 漏数的一枚**，且它比 K1–K3 **更脆**：K1–K3 比的是**名册内容**（缺席的新键不在串里 ⇒ 可能不红），K4 比的是**枚数** ⇒ 只要顶层真多一枚**这次就有值的**键就红。**反向读法也成立**：`omitempty`＋pointer＋reader 未设 ⇒ 恒 4 ⇒ K4 不红 |

**嵌套层同一族的两枚**（`subagent_stream_197_test.go`，射程只到 `results` 的每一行）：
`:136` 逐字 `	if len(results) != len(agents) {` ⇒ 行数＝子代理数；
`:140-141` 逐字 `		if len(r) != 3 {`→`t.Errorf("a result row carries %d JSON keys, want the three ResultChunk has: %v", len(r), r)` ⇒ **`ResultChunk` 的键数被钉成 3**。
⇒ **后果具名**：序号那一枚**绝不可挂在 `ResultChunk` 上**（一挂即红这两枚，且它还同时在两把双向尺的名册里，见族②）。

### 1.2 族②｜双向对账尺 ＝ **2 把**（复认 `167-a3`/`167-c3` 的"两把"，名册行号现量）

两把共用 helper：`jsonKeysOf` 定义 `internal/panel/approval_test.go:174-191`，**盲区就在 `:182` 逐字 `		if tag := f.Tag.Get("json"); tag != "" && tag != "-" {`** ⇒ 无 tag 的导出字段整条跳过、而 `encoding/json` 照样按 Go 字段名发上线。

| 尺 | 锚 | 型别名册（现读 `pairs` 字面量） | 断言原句（逐字，两把同形） |
|---|---|---|---|
| A | `internal/panel/composer_test.go:48` `TestComposerContractTypesMatchFrontend` | `:60-65` **6 对**：`Snapshot`↔`PanelSnapshot`(:60)／`ComposerState`↔`ComposerState`(:61)／`ModeView`↔`ComposerMode`(:62)／`WorkspaceView`↔`ComposerWorkspace`(:63)／`AttachmentRef`↔`ComposerAttachment`(:64)／`ResultChunk`↔`ResultChunkView`(:65) | `:73-74` 逐字 `		if missing := subtract(goKeys, tsKeys); len(missing) > 0 {`→`t.Errorf("Go %s emits %v that interface %s does not declare", p.name, missing, p.tsIface)`；`:76-77` 反向逐字 `t.Errorf("interface %s reads %v that Go %s never sends - those fields render as undefined", …)`；输入 `:50` 逐字 `os.ReadFile(filepath.Join(root, "frontend", "src", "lib", "panel.ts"))` |
| B | `internal/panel/approval_test.go:105` `TestApprovalCardViewJSONKeysMatchFrontendTypes` | `:118-120` **3 对**：`ApprovalCardView`↔`ApprovalCardView`(:118)／`ResultChunk`↔`ResultChunkView`(:119)／`Snapshot`↔`PanelSnapshot`(:120) | `:128-132` 同一条断言体（两方向各一次 `t.Errorf`） |

⇒ **落地腿的四格落点 ⇒ 响几把**（射程判断，⛔ 非颜色判断）：
`ComposerState` 内＝**只响 A**；`Snapshot` 顶层＝**A＋B＋族①那 4 枚**＝**6 发**；`ApprovalCardView`＝**只响 B**；`ResultChunk`＝**A＋B＋`:140` 那枚键数钉**；`TaskRosterSection`／`TaskRowView`／`InstructionsSection`／`GitView`＝**两把尺全盲**（不在九对里）⇒ 只有各自包的词面钉管它（族⑥）。

### 1.3 族③｜**入向方法名**名册钉（本腿只在问四展开，此处登记枚数）

`internal/panel/git_test.go:387`／`:519`（真树 want-4）＋`:516`（正控 want-5）；`internal/panel/inbound_roster_253_test.go:445`／`:455`（两枚写死的数）＋`:435`／`:474`／`:498`／`:510-517`（四条双向分母）；`internal/panel/bridge_test.go:138-145`／`:161`；`internal/panel/composer_test.go:409-419`（闭合词汇表 5 枚）；`internal/panel/l2_grant_boundary_test.go:1293-1301`／`:1857-1861`。**逐枚原文与代价在 §4。**

### 1.4 族④｜反射形状钉：`NumField()` **8 枚**、`NumMethod()` **3 枚**（S01/S02 全列，逐枚过目）

| 锚 | 数的是谁 | 落地腿加字段会不会撞 |
|---|---|---|
| `internal/panel/approval_test.go:177` | helper `jsonKeysOf` 体内（＝族②的发动机） | 不是钉，是尺 |
| `internal/panel/l2_grant_boundary_test.go:279` | 冻结件内 helper | 射程见 §3 |
| `internal/panel/l2_grant_boundary_test.go:1755` | 冻结件内 `reflectionTopKeys`（`:1753` 起） | 射程见 §3 |
| `internal/panel/composer_test.go:173` | `ModeView` 的**方法**枚数 | `	for i := 0; i < mt.NumMethod(); i++ {`→`:174` 逐字 `t.Errorf("ModeView has method %q - a mode view must carry no behaviour at all", mt.Method(i).Name)` ⇒ **给 `ModeView` 加任何方法＝必红**；加字段不红 |
| `cmd/wisp/config_receipt_255_test.go:265` | `reflect.TypeOf(config.Config{})`（`:264`） | 不在本腿写面（`internal/config`），零射程 |
| `cmd/wisp/firstrun_198r2_test.go:89`／`firstrun_198_test.go:52` | 同上，扫 config 树 | 零射程 |
| `cmd/wisp/panel_config_248_test.go:389` | `jsonKeys248` 体内（对 `panel.ComposerRequest`／本地 `credentialWriteRequest`／`plantedCarrier` 取键，调用点 `:330`/`:331`/`:357`/`:373`） | **有射程**：给 `ComposerRequest` 加一枚入向字段 ⇒ 这一族键推导被 `:330`/`:373` 拿去与 `credentialWriteRequest` 对撞；细节见 §3.3 |
| `cmd/wisp/subagent_selfapproval_197_test.go:651`/`:663-664` | `methodNames197`＋`findAllowDoors197`，判据在 `:569` 逐字 `	if len(gateNames) != 2 || !hasMethod197(gateNames, "PendingWindow") || !hasMethod197(gateNames, "PendingApproval") {` | **数的是 `tools.Gate` 接口的方法枚数＝2** ⇒ 与"停止"无关（那是 Gate，不是 TaskRoster）；落地腿若给 `Gate` 加方法＝红这一枚 |

### 1.5 族⑤｜**枚数写死**的钉（`grep -rnE "want[A-Za-z0-9_]*[[:space:]]*=[[:space:]]*[0-9]+" … --include=*_test.go` ＝ **5 枚**，全在 `internal/panel`）

| 锚 | 逐字 | 数什么 |
|---|---|---|
| `inbound_roster_253_test.go:59` | `	wantFullInboundRosterSize253 = 6` | 入向名册枚数（问四主钉） |
| `inbound_roster_253_test.go:64` | `	wantNonPanelPrefixedInbound253 = 2` | 无 `panel.` 前缀的入向枚数（注释逐字 `"config.get", "config.set"`） |
| `l2_grant_boundary_test.go:1334-1336` | `	wantGrantRoutePrefixes = 8`／`	wantGrantRouteWords    = 11`／`	wantGrantRouteSuffixes = 3` | 判决词网格三个因子；判据 `:1487` 逐字 `	if len(grantRoutePrefixes) != wantGrantRoutePrefixes || len(grantRouteWords) != wantGrantRouteWords || len(grantRouteSuffixes) != wantGrantRouteSuffixes {` |

另有 **8 枚 cmd/wisp 侧的写死枚数**（S15b 筛出，逐枚复读）：
`config_receipt_255_test.go:233` `	if len(fields) < 14 {`（红句 `:234` 逐字 `"the roster adjudicates %d Config sections, want the 14 hot ones: the walk itself narrowed"`）／
`firstrun_257_test.go:213-214` `	if len(t257RosterWalk()) != 7 {`→`t.Fatalf("the roster walk is %d fields, want 7 - AC#1 counts the panel's seven", …)`／
`dataroot_128_test.go:211`（`rescueMarkers128` 至少 6 枚）／
`providers_test.go:112` `		if len(fields) != 4 || fields[2] != "实测" {`／
`panel_pump_test.go:302` `	if strings.Join(live.Composer.Mode.Names, ",") != "ask_every_step,ask_high_risk,auto_approve" {`／
`subagent_carrier_197_test.go:549` `	if len(driven.stream.Chunks()) <= panel.DefaultStreamKeys {`（前置防空读，非名册钉）＋`:587` `	if len(seenKey) != panel.DefaultStreamKeys+1 {`／
`resident_approval_risk_256_windows_test.go:489` `	if len(got) != len(want) {`（红句 `:490` 逐字 `t.Errorf("resident leg passes %d of Options' fields, want %d", len(got), len(want))`）＋`:485-486` 逐字 `"the resident leg passes an unexpected Options field %q: the shipped field set is no longer the one ticket 256 adjudicated…"`／
`internal/panel/subagent_roster_197_test.go:275-276` `	if len(sect.Rows) != 20 {`→`t.Fatalf("rows = %d, want all 20 D43 names carried", len(sect.Rows))`。

⇒ **对本腿射程的判断**：这 8 枚里没有一枚数的是快照出向字段或 `panel.*` 入向名册 ⇒ 四枚输出"看得见"那一半**一枚都不撞**；唯一有条件撞上的是 `resident_approval_risk_256`（若"按得下去"那一行改成往 `approval.Options` 里加字段＝问四的范围外，且那一寸属 `Q-76` 未批支）。

### 1.6 族⑥｜**词面钉／手抄名册**（不会响、但会咬注释与新增文件的那一族）

| 枚 | 锚 | 判据原句 | 射程（★这一族是落地腿最容易白撞的） |
|---|---|---|---|
| W-1 | `internal/panel/composer_test.go:268-270` 定义 `gitSwitchCapabilityRe`，扫的机器在 `:276` `TestNoGitSwitchCapabilityInThePanelSurface`；对 `internal/panel` 的接受集逐字 `:317-319` `	scan(filepath.Join(root, "internal", "panel"), func(p string) bool {`／`		return strings.HasSuffix(p, ".go") && !strings.HasSuffix(p, "_test.go")` | 红句 `:321` 逐字 `t.Errorf("a git branch/repository switcher exists in the panel surface; owner cut it from the product (票 92 AC#7):%s", …)` | ★**会撞**：正则含 `\bworktree\b`、`\bgit\s+checkout\b`、`\bgit\s+switch\b`、`\bbranchSelect(or)?\b`、`\bgit\.branch\b`、`\bvcs\.switch\b` 等，**逐行匹配、注释不豁免**，且**新文件自动进射程**（`WalkDir` 不认名册）。⇒ 落地腿若在 `internal/panel/*.go` 的注释里写"worktree"这一枚词（说明工作区/树时极易写出），当场红 |
| W-2 | `internal/panel/composer_test.go:332` 定义 `modeSetterRe`＝`func\s+(Set\|Write\|Apply\|Change)Mode\b\|\bPermissionMode\s*=\|perm\.Store\.Set\(`；判据 `:180` 逐字 `	if hits := modeSettersIn(string(src)); len(hits) > 0 {` | 红句 `:181` 逐字 `t.Errorf("internal/panel exposes a mode setter (%v) - the panel may only display and request", hits)` | ★**会撞**：`modeSettersIn`（`:334-342`）把 `composer.go` **整份按行切开逐行匹配**（`:176` 读全文）⇒ 占用/序号字段若在 `composer.go` 里写出一行含 `SetMode`／`PermissionMode =`／`perm.Store.Set(` 的字面（**含中文注释里夹的英文标识符**），这发红；删 `ModeRequest` 则撞 `:183-184` |
| W-3 | `cmd/wisp/panel_assets_143_test.go:52-61` `type cardView143 struct{…}` **手抄 8 枚卡面键**，解码在 `:112` 逐字 `	if err := json.Unmarshal([]byte(stdout), &card); err != nil {`（注释 `:49-51` 逐字 `// cardView143 re-declares the printed keys rather than decoding into`／`// panel.ApprovalCardView: the reading has to be about the bytes this command`／`// emits, which is what the panel consumes.`） | — | **不会响**（裸 `Unmarshal` 对多余键沉默）⇒ 后果是**名册静默过期**：卡面加 `position` 后这张手抄册少抄一枚，没有任何测试变红。★这正是"尺没响≠没撞契约"的第二处仪器债（第一处＝族②的 tag 盲区） |
| W-4 | `internal/panel/instructions_200_test.go:300` 词面册含逐字 `"depth":0`、`"tier":"project"`、`"instructions"`、`"bytes":` | 同一段字节里 `depth` **已被占用** | ⇒ 序号那一枚**命名 `depth` 会在这发里造出第二枚同名键**；`167-c3` 已具名同一条（本腿独立复认，锚号现量） |
| W-5 | `cmd/wisp/panel_pump_test.go:97` `	if rec["depth"] != fmt.Sprint(len(snap.Pending)) {`＋`:203` `		if rec["depth"] == "1" && strings.Contains(rec["pending"], "pump-corr") {` | 红句 `:98` 逐字 `t.Errorf("ledger depth %q, retained packet holds %d cards", rec["depth"], len(snap.Pending))` | 读的是**账本 k=v**、不是快照 JSON ⇒ 卡面/快照加键不触；但**改 `depth=` 这一枚 token 的语义或写法会触** |
| W-6 | `internal/panel/bridge_test.go:138-145`（四枚方法名必须在 `panel.ts` 里出现）＋`:161` `	if n := strings.Count(text, "bridge.postMessage"); n != 2 {` | 红句 `:162-163` 逐字 `"postMessage call sites = %d, want 2 (the approval request and sendRequest) - a third one means a second envelope was invented"` | 数的是**页面文件里的字面枚数**；Go 侧出向加键看不见 ⇒ 但它就是"第五枚入向方法名需要页面同批声明"的那一把闸（问四） |

### 1.7 族⑦｜**文件／目录枚数**钉（新增文件这一维的射程）

| 锚 | 判据原句 | 射程 |
|---|---|---|
| `internal/panel/frontend_hygiene_test.go:133`（`reachableFrom` 体内） | `	if len(seen) < 5 {`→`:134` 逐字 `t.Fatalf("import closure of main.tsx holds only %d files - the walk broke and the checks below would be vacuous", len(seen))` | 扫 `frontend/**` ⇒ **Go 侧新增文件零射程**；落地腿**不得**在 `frontend/` 落任何东西（票面硬约束 3） |
| `internal/panel/composer_dispatch_test.go:592` 起 `walk := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {`，目录跳过名册在 `:605` 逐字 `			case ".git", ".scratch", "node_modules", "dist", "third_party", "scripts", "docs", "frontend", "design", "build":`（**十枚**，`internal` 与 `cmd` 都不在其中） | 全仓 Walk：任何新增 `internal/panel/*.go` **自动落进射程**，无需有人改名册 | ★**有条件会撞**：这枚就是"新文件＝新射程"的具体形状（尺体在 `:628`/`:637` 逐行记 `rel:line:imports <path>`／`rel:line:identifier <sym>` 两类命中）⇒ 落地腿若新建文件，它被扫的不是"文件名"而是**import 名册与标识符**；判据本体属 §3 的射程复核 |
| `cmd/wisp/panel_host_gate_test.go:90-95`＋`:109-111` | `:90` 逐字 `		tracked, _ = gitLinesInDir(root, "ls-files", "frontend/dist")`；`:93` `trackedBeyondAnchor := 0`；注释 `:60` 逐字 `// tracked/ignored axes are read as ENTRY NAMES from git only.` | AC#12 只数 `frontend/dist` 的 git 跟踪条目 vs embed 的 built 位 ⇒ Go 侧加字段/加方法**零射程** |
| `internal/panel/tokens_fourway_test.go:212` | `t.Fatalf("parsed only %d colour rows from %s - the table shape changed and this check would be vacuous", len(rows), c21TablePath)` | 输入四枚路径逐字（`:50-53`）：`design/assets/tokens.css`／`docs/evidence/s1/c21-native-tokens.md`／**`internal/ball/tokens.go`**／`frontend/src/styles/tokens.generated.css` ⇒ 落地腿零射程（不碰这四枚），**但它数的是 `docs/evidence/s1/` 里那一枚 md 的行枚数**⇒ 谁往那枚 md 里追加颜色行谁响 |

### 1.8 本问小结（给问六当料的两句话）

1. **真正的硬碰撞只有三处**：族②两把尺（逐名对账）、族①四枚顶层键数/键集钉、族⑥W-1/W-2 两枚**逐行扫 `internal/panel` 生产源码**的词面钉。其余计数钉要么数的是别的地界（config 树、`frontend/`、`internal/ball`、`tools.Gate`），要么数的是**运行期条数**（属问二）。
2. **两处"静默过期、不会响"**（W-3 手抄八键册 ＋ 族② tag 盲区 `approval_test.go:182`）比会响的钉更危险：落地腿很容易把它当成"没撞契约"的证据，而 `167-c3` §3.5 与票面 §9 都已把这条列为在册盲区，本腿复认成立。

## 2. 问二 — 四枚输出各自的现有断言现状（占用／序号／停止／草稿＋崩溃自救）

> 本节一律按派单点名的**三处法**答："零调用者"这种负向句要问**谁产出／谁投递／谁落盘（谁出向）**三处，不能只数显式调用点。尺与时刻逐格附。

### 2.1 ① 占用条 —— **零断言**（两把独立尺互相否证，16:5x）

| 尺 | 读数 |
|---|---|
| `grep -rniE "occupanc\|contextused\|usedtokens\|windowused\|prompttokens" cmd/wisp internal/panel --include=*_test.go \| wc -l` | **3** 行，逐枚过目**全部是热键占用**那一枚 occupancy 的注释/红句，与上下文占用无关：`cmd/wisp/resident_hotkey_live_258_windows_test.go:249`／`:259`（注释）与 `:284` 逐字 `t.Fatalf("258 LIVE RULER RED: cannot build the occupancy premise for %s: RegisterHotKey answered %v. …", binding, errno)` |
| `grep -rn "ContextWindow" cmd/wisp internal/panel --include=*_test.go \| wc -l` | **4** 行，逐枚过目**全部是配置字段写入**那一枚，不是占用显示：`cmd/wisp/firstrun_257_test.go:114-115`（`m.SetModelContextWindow(t257Provider, t257Model, 128000)`）／`internal/panel/config_route_248_test.go:272`／`:296`（`FieldModelContextWindow` 作为 `configField` 被写） |

⇒ **裁**：今天**没有任何一枚用例断过"占用"**这件事——既没有断过出口存在，也没有断过它不存在。落地腿新增这一格＝**零枚现有断言要改写**，但要新增自己的断言（AC#2 那两发变异是它的凭据）。⚠ 唯一有条件咬它的是 §1.6 W-1/W-2 那两枚**逐行扫源码的词面钉**（写注释时的用词），与 §1.1 族①那四枚顶层键钉（落点选择）。

### 2.2 ② 排队序号 —— **有断言"队列条数"，零断言"第几条"那一枚值**

现量（尺＝`grep -rnE "\bPosition\b\|\bDepth\b" cmd/wisp internal/panel --include=*_test.go \| wc -l` ＝ **37** 行；逐枚过目后**其中 27 行是 `fset.Position(...)`＝AST 位置，与队列无关**）：

| 锚 | 原句（逐字） | 它断的是什么 |
|---|---|---|
| `cmd/wisp/resident_grant_writer_265_windows_test.go:196` | `		if n := f.g.Queue().Depth(); n != 0 {` | 队列**深度＝条数**，且断的是"应为 0"这一形 ⇒ 序号（1-based 的第几条）**没被任何用例读过** |
| `cmd/wisp/panel_pump_test.go:97-98` | `	if rec["depth"] != fmt.Sprint(len(snap.Pending)) {`→`t.Errorf("ledger depth %q, retained packet holds %d cards", rec["depth"], len(snap.Pending))` | **账本 k=v 里的 `depth=`＝`len(snap.Pending)`** ⇒ 这是今天唯一一枚把"队列长度"写成可断言之形的地方；它读日志、不读快照 |
| `cmd/wisp/panel_pump_test.go:203` | `		if rec["depth"] == "1" && strings.Contains(rec["pending"], "pump-corr") {` | 同上（另一发，字面 `"1"`） |
| `internal/panel/instructions_200_test.go:300` | 词面册含 `"depth":0`、`"tier":"project"` | **`depth` 这一枚键名在项目说明那一节里已被占用**（＝§1.6 W-4）⇒ 序号若复用 `depth` 这枚拼写，两节会在同一份字节里同名 |
| `cmd/wisp/instructions_200r2_test.go:151`／`:239` | `			if sec.Files[i].Tier != "project" \| sec.Files[i].Depth != 0 {`／`			if f.Tier != "global" \| f.Depth != -1 {` | 数的是**说明文件的目录深度**，与排队无关（登记以免被误当序号钉） |
| 另有 **13 枚** `len(….Pending)` 形状的运行期断言（尺＝`grep -rn "len(snap.Pending)\|Pending) != \|Pending) == " cmd/wisp internal/panel --include=*_test.go \| wc -l`），例：`cmd/wisp/subagent_blocked_197_test.go:164` 逐字 `	if len(snap.Pending) != 1 \|\| snap.Pending[0].CorrelationID != id {` | 断的是"卡在哪一张"，不是"第几条" | ⇒ **落地腿若改 `Pending` 的填充语义（不只是加键），这 13 枚都在射程里**；只加 `position` 一枚键 ⇒ 一枚都不响 |

⇒ **裁**：`ApprovalCardView.Position` 今天**零断言**（`grep` 里那枚 `Position` 全是 AST 位置）；真源三枚（`LiveApproval.Position`／`Queue.Depth()`／`PanelItem.Depth`）在 `167-c2`／`167-a3` 已量齐，本腿只补一句：**`internal/agent/approval` 那两枚今天一次都不出向**（`grep -rc "json:\"" internal/agent/approval` 两口径均 0，复认 `167-a3` R-⑥）。

### 2.3 ③ 停止 —— 三处法逐处量（**产出有、投递零、出向零**）

| 处 | 现量（尺＋时刻 16:5x） | 结论 |
|---|---|---|
| **谁产出** | `grep -n 'func (r \*TaskRoster)' internal/tools/task.go` ＝ **13 枚**，其中停止那一支三枚：`:420 AttachCancel(taskID, cancel)`／`:434 DetachCancel(taskID)`／`:447 Cancel(taskID) (bool, string)`；`:418` 注释逐字 `// AttachCancel stores the child's own context.CancelFunc as the ONLY stop path` | 名册在、且**没有"能不能停"的谓词**（13 枚逐枚过目：Record/WatchRow/Look/Count/Descendants/PublishSubagent/MarkRoot/TryAcquireSubagentSlot/InFlightSubagents/RunningSubagentIDs/AttachCancel/DetachCancel/Cancel）⇒ **推翻 `docs/reports/survey-2026-09-28-dsh-ui-packages.md:191` 那句"方法 15 枚"**（那是过期或含非导出，口径不同；票面 §9 的"13 枚"复认成立） |
| **谁投递** | `grep -rn "AttachCancel" cmd internal tools` 全命中 **4 行产码**：定义 `internal/tools/task.go:418/:420`、注释 `internal/tools/subagent_197.go:38`、**唯一产码调用者 `internal/tools/subagent_197.go:346` 逐字 `	t.d.Roster.AttachCancel(bg.ID, cancelChild)`**；`cmd/wisp/**` 命中 **0 枚** | ⇒ **主 run 那条零调用**复认成立（票面 §9 第 1 条末行、`A` 台账那句都对得上）；`requestTaskStop` 在 `cmd/wisp`＋`internal/panel` 测试里 **0 命中**（尺＝`grep -rniE "requestTaskStop\|AttachCancel" cmd/wisp internal/panel --include=*_test.go`＝**0 行**） |
| **谁落盘／谁出向** | `grep -rniE "cancellable" cmd internal tools` ＝ **3 行、全部是注释**（`cmd/wisp/approval_reply.go:473`、`internal/agent/approval/cancel_key_label_260r4_test.go:24`／`:256`）；`TaskRowView` 现量 **12 枚 json 键**（`internal/panel/subagent_roster_197.go:106-138`：`taskId/label/kind/parentTaskId/status/statusKnown/statusReason/streamKey/blockedOnApproval/streamTruncated/streamElidedRunes/streamDropped`）⇒ **无 cancellable 位** | ⇒ 票面 §9"'看得见'缺位三处"的第 2 处（12 枚键无 cancellable 位）**逐枚复认**；第 3 处（根行 `statusKnown` 那一维被票 196／`A394` 占着）在场：`:115` 逐字 `	StatusKnown bool   \`json:"statusKnown"\`` ＋ `:111-113` 注释自述它是 D43 状态名 ⇒ ⛔ 本腿不判要不要动它 |

⇒ **裁**：**停止这一枚今天零断言**（既无 `AttachCancel` 的调用者用例，也无 `cancellable` 的出向用例）。相邻的只有 `l2_grant_boundary_test.go` 那族**入向名册钉**（§4）与 `panel_resident_windows_test.go` 那族线程收尾钉（§5）。

### 2.4 ④ 草稿 —— **门有三枚断言、"存住／不覆盖"零断言**（★本腿抓到一枚"必须同批改写"的钉）

| 锚 | 原句（逐字） | 射程判断 |
|---|---|---|
| `cmd/wisp/panel_config_248_test.go:460` `TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket`，判据 `:467-471` | `	for _, method := range []string{panel.MethodWorkspaceRequest, panel.MethodAttachmentAdd, panel.MethodMessageSend} {` → `		if _, err := disp.Handle(context.Background(), raw); !strings.Contains(err.Error(), "处理器未接入") {` → 红句 `t.Errorf("%s with no socket: %v, want the unattached refusal", method, err)` | ★**这枚正向断言的是"message 那扇门必须仍然按名拒绝"**。本腿追到它的构造体：`cmd/wisp/panel_config_248_test.go:126-150` `newRealSettingsLeg` **自己在 `:145-148` 造了一枚只带 `Config`＋`Audit` 的 dispatch**（注释 `:124-125` 逐字 `// newRealSettingsLeg builds the leg out of the two production objects, the way`／`// newComposerDispatchChain does, and returns the dispatch in front of it.`）⇒ **它不读生产装配**，所以把 `cmd/wisp/panel_inbound.go:277` 的 `Message: nil` 换成真处理器**不会**让这发红。⚠ 但它的**注释与用例名**（"the three doors nobody owns still refuse by name"、`"Untouched by this slice, and refused by name when they arrive: … message = ticket 35"`，见 `panel_inbound.go:273-274`）在那一发之后就成了**假话**：字面尺不响、语义已过期 ⇒ 这一格属"该同批改写、但仪器不会替你响"的那一类，落地腿若动 `Message` 就**必须**在交件里具名说它动了这枚用例的语义 |
| `cmd/wisp/panel_inbound_33_test.go:158` `TestAC9InboundLegRefusesRosterMethodWithNoHandler`，判据 `:165-172` | `	if code != 1 {`→`t.Errorf("exit=%d want 1 (the request was refused): stdout=%q", code, out)`；`:168` 逐字 `	if !strings.Contains(out, "处理器未接入") {`；`:171` 逐字 `	if !strings.Contains(errLog, "panel.workspace.request") {` | 跑的是**真腿**（`runInboundLeg33`），但送的是 `panel.workspace.request` ⇒ 落地腿只接 `Message` **不撞**；若哪天有人顺手把 workspace 也接上 ⇒ 撞这一枚 |
| `internal/panel/composer_dispatch_test.go:291` | `	d := &ComposerDispatch{Mode: modeSpy, Workspace: other, Attachment: other, Message: other}` | 四枚插座的**桩**（既有形状，票面 §9 第 3 条说的"`ComposerDispatch` 的四枚插座"就是这里）⇒ 落地腿接 `Message` 不必新造接缝，这条有先例 |
| 门本身（产码在场，供对照；锚点 `2a0e23b4`，17:1x 逐枚现量） | 名册 `internal/panel/bridge.go:42-45` 四枚逐字 `	MethodMessageSend      = "panel.message.send"`（`:45`）；白名单守卫 `bridge.go:146` `func knownComposerMethod(m string) bool`（调用点 `:132`）；派发表 `internal/panel/composer_dispatch.go:175` `func (d *ComposerDispatch) dispatch(…)` 的六枚 `case`（`:177`/`:182`/`:187`/`:192`/`:197`/`:202`＝**与入向名册的 6 枚同数**）；承载位 `internal/panel/bridge.go:98` 逐字 `	Text string \`json:"text,omitempty"\``；生产装配 `cmd/wisp/panel_inbound.go:275-277` 逐字 `		Workspace:  nil,`／`		Attachment: nil,`／`		Message:    nil,`（**未漂**，票面 §9 第 ④ 枚复认） | 出向零枚：草稿**没有任何承载位回去**；持久面 `internal/store` **目录不存在**（尺＝`ls -d internal/store` ⇒ 不存在）；全仓 `grep -rniE "draft\|草稿" cmd internal tools` ＝ **15 行**、逐枚过目**全部**是"这份测试的第一稿"这类注释与一处 fixture 字符串（`internal/tools/fs_edit_ac34_test.go:75`），**零枚关于输入框草稿** |

⇒ **裁**：**"失败原文存住"与"不许覆盖用户后来新写的字"两格今天零断言**；有断言的是"这扇门还没接线"这一枚**负向现状钉**（上面第一、二行），其中 248 那一枚是**语义级**的（仪器不响、话变假）。

### 2.5 ⑤ 崩溃自救（AC#6 三枚拆开的现有断言）

| 枚 | 三处读数（16:5x 现量） | 今天有没有用例断过 |
|---|---|---|
| 崩溃栈落盘 | 产＝`internal/observe/goroutine.go`（`grep -rln "ev.Stack" internal/observe` 命中 `goroutine.go`／`goroutine_test.go`）；投＝`cmd/wisp/logsink.go`；落＝`internal/observe/redact.go`（同一把尺三枚文件命中） | **`cmd/wisp`＋`internal/panel` 里零枚断过它**。尺＝`grep -rnE "MaxLoggedString\|BuildDiagnosticsBundle\|Panic" cmd/wisp internal/panel --include=*_test.go`＝**0 行**；放宽到 `grep -rn "panic" … \| grep -iE "sink\|stack\|recover\|crash\|diagnos"`＝**8 行**，逐枚过目全是**别的事**（`panel_resident_windows_test.go:465/:488` 是 WM_QUIT panic 的复现钉、`resident_sink_nail_127_windows_test.go:599`／`leg_sink_nail_131_windows_test.go:212` 是"panic 会吞红名"的仪器钉、`leg_sink_gate_131_test.go:1351` 是内建名册）⇒ **512 那一枚截断今天无人断** |
| 可复制现场 | `grep -rn "BuildDiagnosticsBundle" cmd internal tools \| grep -v _test.go \| wc -l`＝**3**，且 3 枚全在它自己那一枚文件（`internal/observe/diagnostics.go:35` 注释／`:61` 注释／`:62` 定义）⇒ **产码调用者 0 枚**复认（票面 §9 第 ⑥ 枚、`A` 台账 09:36 那次复跑同读数） | 断它的用例住在 `internal/observe/diagnostics_test.go`（`:35`／`:125`／`:140` 三处调用），**不在本腿射程**；`cmd/wisp`／`internal/panel` **零枚** |
| 自救指令 | `grep -n 'pass("\|fail("\|info("' cmd/wisp/doctor.go \| wc -l`＝**28**（复认票面 §9 第 2 条那枚更正后的数）；`grep -n "logs\|logDir" cmd/wisp/doctor.go \| wc -l`＝**1**，且那一行是**注释** `:243` ⇒ **零枚检查项读 `logs\`** 复认；`grep -rn "logDir()" cmd/wisp \| grep -v _test.go`＝**1 行**＝定义本体 `cmd/wisp/logsink.go:99` ⇒ **产码调用者 0 枚**（该格票面原标〔仅自述〕，**本腿复跑成立**） | **零断言**（无人断 doctor 该读到日志目录，也无人断它没读到） |

⇒ **裁**：AC#6 三枚在 `cmd/wisp`＋`internal/panel` 侧**全部零断言** ⇒ 落地腿新增自救出口**不打红任何现有用例**；它唯一的相邻火源是 `cmd/wisp` 那一族 **leg/sink 门**（`leg_dispatch_gate_133_test.go`／`leg_sink_gate_131_test.go` 按 AST 数 `main.go` 的腿与 sink 站点）——那两族数的是**分发腿名册**，若新增一枚 `wisp` 子命令就要一并过它们（§6 第 5 条）。

## 3. 问三 — `internal/panel` 三枚冻结件的射程（`tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`frontend_hygiene` 一族）

> ★**先声明性质**：本节是**射程判断**（这枚尺管不管得到落地腿要动的那一格），⛔ 不是内容引用、不是摘录复用；判据＝现读该文件内部判据本体。三枚的"一字不许改"在册文本＝`.scratch/wisp/dispatches/2026-09-26-093x-thaw-panel-for-145.md` §3 逐字：`:15` `internal/panel/tokens_fourway_test.go` —— owner 的 `Q-52` 撤回令在效：界面未定稿期间那枚**恒红**是已知常红，**不许为它变绿做任何事**；`:16` `internal/panel/l2_grant_boundary_test.go` —— 票 128／门钉那一族的既有钉子；`:17` `internal/panel/frontend_hygiene_test.go` —— 它断言的是 `frontend/**` 的形状，而那半棵树**正被前端会话改写且未定稿**；**`:19` 还加了一枚范围扩张**逐字 `internal/panel/** 里其余文件若与上面三枚同名族（tokens*、l2_grant*、frontend_hygiene*），一律按禁改面处理，不确定就停手报回` ⇒ ⚠ 落地腿若**新建**同名族文件（例如 `internal/panel/tokens_167_test.go`）＝撞 `:19`，属停手报回格，不属"能写"格。

### 3.1 `tokens_fourway_test.go`（550 行、**1 枚**测试，`TestC21DesignTokensFourWayAgree` 起 `:439`）

| 内部判据（现读） | 射程 |
|---|---|
| 输入只有四枚路径，逐字（`:50-53`）：`tokensCSSPath     = "design/assets/tokens.css"`／`c21TablePath      = "docs/evidence/s1/c21-native-tokens.md"`／`nativeTokensPath  = "internal/ball/tokens.go"`／`frontendThemePath = "frontend/src/styles/tokens.generated.css"` | ⇒ **新增快照字段／新增入向方法名：零射程**（这四枚没有一枚在 `cmd/wisp` 或 `internal/panel` 的可写面上） |
| 三枚防空读：`:212` 逐字 `		t.Fatalf("parsed only %d colour rows from %s - the table shape changed and this check would be vacuous", len(rows), c21TablePath)`；`:285` 逐字 `		t.Fatalf("native palettes not parsed (dark %d / light %d fields)", len(out["dark"]), len(out["light"]))`；`:456` 逐字 `			t.Fatalf("theme %s: %s parsed no declarations, this check would be vacuous", theme, tokensCSSPath)` | ⇒ 这族计数钉只被**颜色行／调色板字段**喂养。⚠ **对本腿唯一有价值的相邻读数**：`:52` 说明它读 **`internal/ball/tokens.go`** ⇒ 这一枚冻结件的活性射程在 **260-r1 的写面**上、不在 167-r2 的写面上；落地腿不许"顺手"拿它当自己安全的凭据，也不许为它变绿做任何事（`:15` 的 Q-52 令） |

**裁**：167-r2 动不到它；它今天红不红与 167-r2 无关（`167-a3` §5.1 U-1 已具名其红因在 `design/**` 工作树，本腿不复跑、不引为读数）。

### 3.2 `l2_grant_boundary_test.go`（2549 行、**7 枚**测试：`:1229`／`:1530`／`:1547`／`:1595`／`:1807`／`:1959`／`:2168`）——★**三枚里唯一真会咬落地腿的那一枚**

它把整个 `internal/panel` **生产源码**按 AST 解析（`dir` 逐字 `filepath.Join(root, "internal", "panel")`，见 `:1231-1232`、`:1847-1848`）⇒ **新增文件自动进射程，无需有人改名册**。逐枚判据：

| 内部判据（现读原句） | 落地腿新增快照字段／新增入向方法名会不会撞 |
|---|---|
| `:966` `func instrumentBlindnessProblem(dir string, pkg boundaryPackage) string {` 的防空读六支，逐字含 `:976` `	if len(pkg.answered) == 0 {`、`:979` `	if len(pkg.structs) == 0 {`、`:982` `	if len(pkg.decodes) == 0 {`；经 `:992` `func requireReadableInstrument(t *testing.T, dir string, pkg boundaryPackage) {` 的 `:994`–`:996` 升为 `t.Fatalf` | **有条件撞**：`:967` 要求 `knownComposerMethod` 仍在那棵树里、`:982` 要求 `internal/panel` 至少还剩一枚 JSON decode 站点 ⇒ 只要不删这两类东西就不响；**新增文件不会让它红**（它只防空读，不设上界） |
| ★`:1857-1859`（在 `TestJSONKeyDerivationAgreesWithEncodingJSON` 里）逐字 `	for seed := range pkg.inboundSeeds {`／`		if _, ok := reg[seed]; !ok {`／`			t.Fatalf("decode destination %s has no reflection twin in inboundTypeRegistry: this test would compare the two instruments over different trees, and the AST half would be unchecked", seed)`；而 `inboundTypeRegistry()` 的名册**只有 4 枚**（`:1721-1726`：`ComposerRequest`／`ModeRequest`／`AttachmentPayload`／`AttachmentRef`） | ★★★**最硬的一枚，管的就是"新增入向那一跳"**：`inboundSeeds` 的来源是 `classifyDecodes`（`:678`起），`:700` 逐字 `			pkg.inboundSeeds[d.TypeName] = true` ⇒ **`internal/panel` 里任何一枚 decode 落到的同包 struct 都自动成为 seed**。⇒ **落地腿若为面板新增一枚入向承载体（例如给"停止"或"草稿"新造一个请求 struct），这枚 `t.Fatalf` 立刻红，而唯一修法是往 `inboundTypeRegistry` 加一行＝改冻结件一字＝禁**。⇒ 具名后果：**新增入向形状必须复用 `ComposerRequest`（`:1722` 已在册）**，否则这一格在治理面上是死路。⛔ 本腿不裁要不要新增，只把"新增即 FATAL、且自己修不了"量清楚 |
| `:1878` 子测试 iii 的 24 枚判决词名册，逐字 `:1882-1887`：`"outcome", "Outcome", "allow", "Allow", "allowOnce", "AllowOnce",`／`"approved", "Approved", "grant", "Grant", "verdict", "Verdict",`／`"decision", "Decision", "decide", "Decide", "bypass", "Bypass",`／`"override", "Override", "permit", "Permit", "authorize", "Authorize",`；判据 `:1906` 逐字 `				if known {`→`:1907` 红句 `t.Errorf("%s binds the wire key %q, so a panel can address an approval decision through it (D33/F2, AGENTS.md ban #6). If %s is inbound, this is the finding; if it is not inbound any more, drop it from inboundTypeRegistry instead of dropping the check", name, key, name)` | 只过 `inboundTypeRegistry` 那四型 ⇒ **出向快照字段零射程**；**入向字段名不许取这 24 枚拼写**。落地腿候选名逐枚裁：`position`／`occupancy`／`stopRequested`／`cancellable`／`draftText` **全部无交集**（射程判断） |
| ★`:1293-1301`（`TestAnsweredPanelRoutesCarryNoApprovalDecision`）逐字 `	declared := map[string]bool{}`／`	for name, val := range pkg.consts {`／`		if strings.HasPrefix(name, "Method") {`／`			declared[val] = true`；判据 `:1300` `		if !declared[name] {`→`:1301` 红句逐字含 `"Go answers %q but no Method* constant declares it: a route written straight into the guard's case list. Nothing else in this package would notice - TestFrontendComposerRequestsMatchTheEnvelope (bridge_test.go:131) is the test that reads the frontend for these names, and it only ever iterates the four declared constants, so a fifth route that never became a constant is invisible to it. This line is the only gate on that shape"` | ⇒ **第五枚入向方法名的形状被这枚钉死**：守卫答了就必须有一枚 `Method*` 包级常量（且红句自己承认 `bridge_test.go:131` 只循环四枚常量＝对它失明）。⛔ 出向字段零射程 |
| `:1243` `		if carriesGrantWord(name, grantRouteWords) {`（词表 `:192-195` **11 枚**：`approve`/`approval`/`grant`/`allow`/`permit`/`ratify`/`authorize`/`authorised`/`decide`/`decision`/`verdict`）＋`:1276` `	if len(pool) < len(answered) {` | ⇒ 第五枚方法名不许含这 11 枚词；`stop`／`cancel`／`halt` **不在表内**（本腿逐枚比过）⇒ **名字本身不撞**，撞的是上面那些名册钉 |
| `:1574` `	findings := scanGrantBoundary(t, filepath.Join(root, "internal", "panel"))`；体内 `:1181` `	inbound := inboundEnvelopes(pkg)`、`:1182` `	if len(inbound) == 0 {`→`t.Fatalf("no inbound envelope type (a JSON decode destination, or anything nested in one) under %s: …")`，判据 `:1189` `			if carriesGrantWord(f.JSONKey, grantFieldWords) {`；`grantFieldWords`（`:183-188`）**19 枚**＝路由表 11 枚之外再加 `autoapprove`/`granted`/`allowed`/`allowonce`/`permitted`/`outcome`/`bypass`/`override` 等形；归一化在 `carriesGrantWord`（转小写、去掉 `_` `-` 空格 `.` 再 `strings.Contains`） | ⇒ 只过"入向信封＋其嵌套／嵌入可达的类型"（种子就是 `inboundSeeds`，`:1006-1012`）⇒ 出向快照字段零射程；⚠ 若新字段挂进 `ComposerRequest` 一族，字段名不许含这 19 枚词根（`pendingApproval`／`allowedDraft` 都会红）。**旁证一枚（值钱）**：`TaskRowView.BlockedOnApproval` 的键 `blockedOnApproval` **本身含 `approval`**、今天不红**仅因为它是纯出向** ⇒ "出向＝豁免"是射程事实，不是运气 |
| `:1487` 三枚计数钉＋`:1691` 见证名册（红句逐字含 `"grantRouteWords no longer carries %q while %q still stands here as its witness: the sweep stopped asking the %d names that word used to build…"`） | ⇒ 只在**判决词表被增删**时响 ⇒ 落地腿不碰词表即零射程（且它无权碰：冻结件） |

**裁**：`l2_grant_boundary_test.go` 对**"新增出向快照字段"＝零射程**；对**"新增入向方法名"＝至少两枚**（`:1293-1301` 的 `Method*` 常量钉，加取名不当时的 `:1243`/`:1906` 词表钉）；对**"新增入向承载 struct"＝一枚无法自行修复的 `t.Fatalf`**（`:1859`，修法在冻结件内部）⇒ 这一格请编排者按票面 §9 第 1 条与 `Q-76` 的口径归档，本腿不裁。

### 3.3 `frontend_hygiene` 那一族（`internal/panel/frontend_hygiene_test.go`，324 行、**5 枚**测试）

⚠ **先更正一处前人读数**：`167-c3` §3.6 写的是"六枚测试"，本腿现量 `grep -n "^func Test" internal/panel/frontend_hygiene_test.go` ＝ **5 枚**（`:169`／`:196`／`:222`／`:246`／`:281`），另 5 枚是包级 helper（`:76 frontendSrcFiles`／`:101 reachableFrom`／`:140 resolveSpec`／`:161 relToRoot`／`:317 basenames`）。逐枚射程：

| 枚 | 内部判据（现读） | 射程 |
|---|---|---|
| `TestPanelFrontendIsStateless` `:169` | 文件集＝`:79` 逐字 `	err := filepath.WalkDir(filepath.Join(root, "frontend", "src"), func(p string, d fs.DirEntry, err error) error {`，只收 `.ts/.tsx/.css`（`:86`）；判据是 `bannedStorageAPIs`（`:48-56`，**7 枚**：`localStorage`/`sessionStorage`/`indexedDB`/`cookie`/`caches`/`serviceWorker`/`fs-access`）；防空读 `:172` 逐字 `	if len(files) == 0 {` | **纯 `frontend/**`** ⇒ 落地腿 Go 侧改动**零射程** |
| `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` `:196` | 判据尺是 `colourLitRe`（定义 `:61`，匹配十六进制色值与 `rgb()/rgba()` 两形，注释 `:59-60` 明写它**故意不匹配 `var(--x)`**），走 `reachableFrom` 的 import 闭包；防空读 `:133` 逐字 `	if len(seen) < 5 {`→`:134` 逐字 `t.Fatalf("import closure of main.tsx holds only %d files - the walk broke and the checks below would be vacuous", len(seen))` | 同上，**零射程**（它今天在册红因在页面／`design` 侧，不属本编队） |
| `TestVendoredDemoComponentsAreNotMounted` `:222` | `:225-226` 读 `frontend/src/components/ai-native` 目录；判据 `:239` 逐字 `	if len(mounted) == 0 {` | **零射程**（枚数关系全在页面树） |
| `TestFrontendHasNoEmoji` `:246` | 走 `frontend/src`＋`frontend/scripts`（`:249-260`）再补 `frontend/index.html`（`:261`）；`emojiRangesRe` 在 `:71`，注释 `:64-70` 明写这份副本**无注释豁免、覆盖面比 CI 仪器窄**，逐字 `// So this test can be stricter than the gate, never wider.`（前一行逐字 `// that says so: this copy has no comment exemption (the scanner blanks`） | **零射程**，但★必须具名分清两件事：这一族**不扫 Go 侧**，落地腿界面文案里的 emoji 归 `tools/d22scan`（CI）管、不归这族管（`AGENTS.md §1.2` 已把"规格射程≠仪器射程"登记为未定案缺口，本腿不重开） |
| `TestFrontendNeverNamesAnApprovalDecision` `:281` | 走**整棵 `frontend/`**（`:285`），`:289` 跳过 `/node_modules/` 与 `/dist/`；判据只有一条正则，定义在 `:73`（`panelDecisionIdentifierRe`，匹配的是 `approval.decide` 这一枚字面），注释 `:72` 逐字 `	// panelDecisionIdentifierRe is tools/d22scan's ban #6 pattern (approval.decide).` | **零射程**（唯一命中面是页面树） |

**裁**：这一族对 167-r2 **整体零射程**。⚠ 唯一致命的相邻形状是"落地腿往 `frontend/` 落任何文件"——那是票面硬约束 3（有一枚要写的落在 `frontend/` 下就整段停手上报），⛔ 不是"改一行让它变绿"的选项。

## 4. 问四 — C17 入向方法名名册钉：现量枚数＋"第五枚入向方法名"的代价

未判。

## 5. 问五 — 分位数／真窗那一族的脆弱点与"该钉哪两行"

未判。

## 6. 问六 — 给落地腿的一页清单

未判。

## 7. 量不到／判不动／我写错的读数（自我对抗）／⛔ 未动产码自证

未判。
