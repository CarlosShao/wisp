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

未判。

## 3. 问三 — `internal/panel` 三枚冻结件的射程（`tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`frontend_hygiene` 一族）

未判。

## 4. 问四 — C17 入向方法名名册钉：现量枚数＋"第五枚入向方法名"的代价

未判。

## 5. 问五 — 分位数／真窗那一族的脆弱点与"该钉哪两行"

未判。

## 6. 问六 — 给落地腿的一页清单

未判。

## 7. 量不到／判不动／我写错的读数（自我对抗）／⛔ 未动产码自证

未判。
