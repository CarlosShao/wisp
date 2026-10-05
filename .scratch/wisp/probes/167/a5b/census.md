# 167-a5b — 只读普查腿（接 167-a5 死腿的 §4–§7 四节）：C17 入向名册现量／分位数·真窗族脆弱点／给 167-r2 的一页清单／自我对抗与自证

> 运行类型＝**只读**：⛔ 零 `go` 命令（`268-v1` 正在 `cmd/wisp` 跑整包）。一切"红/不红"都是**射程判断**（读尺体本体得出），颜色属〔预测〕，归口编排者跑。
> 前任 `167-a5` 的 census.md §0–§3 已填实（末笔 `2759d08a` 17:27），§4–§7 四节死于额度中断、只留"未判"占位。**本件补齐那四节**；⛔ 不改 a5 的件一字（死腿件只读）。
> 引用 a5 读数时逐条标注"〔a5 已量〕"（只转录、不复跑）还是"〔本腿复跑〕"（本腿现量）。**所有 file:line 一律本腿自己重取**（锚从 a5 起手的 `94380ab3` 漂到本腿起手的 `2a633eb8`，行号不承袭）。
> 硬边界自证：⛔ 未读不引 `frontend/**`／`design/**`；写面只有 `.scratch/wisp/probes/167/a5b/**`；⛔ 未碰产码／票面／台账／AC 框。

## 0. 起手锚（逐字读数，同一发命令取）

| 项 | 读数 |
|---|---|
| 时刻 | `2026-10-05 15:4x +0800` |
| `git rev-parse --short HEAD` | `2a633eb8` |
| 分支 | `dev` |
| 五根 porcelain（`git status --porcelain -- cmd internal tools scripts docs`） | **0 行**（a5 起手时曾有 1 行 ` M cmd/wisp/resident_approval_windows.go`〔a5 已量〕，本腿起手已不在——别人已提交或还原，本腿只报不改） |
| 产码枚数复跑 | 待 §4 起手量 |
| 本腿写面 | 仅 `.scratch/wisp/probes/167/a5b/census.md`（＋同目录临时 msg 件）；⛔ 未动任何产码／测试／工单／台账 |

## 4. 问四 — C17 入向方法名名册钉：现量枚数＋"第五枚入向方法名"的代价

> 尺（本腿起手 15:4x 现跑，⛔ 零 `go`）：`grep -n 'Method.*=.*"panel\.' internal/panel --include=*.go` 名册面；`grep -n 'knownComposerMethod' internal/panel --include=*.go` 守卫引用面。所有行号本腿现取（锚 `2a633eb8`）。
> ⚠ **a5 §1.3/§1.6 引用的行号全数漂移**：a5 写的 `git_test.go:387/:519/:516` 与本腿现量一致〔本腿复跑成立〕，但 a5 写的 `bridge_test.go:138-145/:161`（现 `:138-141`／`:161`，未漂）、`composer_test.go:409-419`（现 `:409-419`，未漂）、`l2_grant_boundary_test.go:1293-1301`（现 `:1293-1303`，红句本体在 `:1301`）基本未漂——见 §7 第 2 条的逐条核。

### 4.1 现名册本体：`internal/panel/bridge.go` **6 枚**〔本腿复跑〕

| 枚 | 常量 | 值 | 锚（本腿现量） |
|---|---|---|---|
| 1 | `MethodModeRequest` | `"panel.mode.request"` | `internal/panel/bridge.go:42` |
| 2 | `MethodWorkspaceRequest` | `"panel.workspace.request"` | `internal/panel/bridge.go:43` |
| 3 | `MethodAttachmentAdd` | `"panel.attachment.add"` | `internal/panel/bridge.go:44` |
| 4 | `MethodMessageSend` | `"panel.message.send"` | `internal/panel/bridge.go:45` |
| 5 | `MethodConfigGet` | `"config.get"`（**无 `panel.` 前缀**，票 248 设定路由） | `internal/panel/bridge.go:66` |
| 6 | `MethodConfigSet` | `"config.set"`（同上） | `internal/panel/bridge.go:67` |

守卫＝`bridge.go:146` `func knownComposerMethod(m string) bool`，case 表 `:148` 恰好六枚常量；派发＝`internal/panel/composer_dispatch.go:175` `func (d *ComposerDispatch) dispatch(…)`，六枚 `case` 在 `:177/:182/:187/:192/:197/:202`（与名册 6 枚同数〔a5 已量，本腿复跑成立〕）。名册、守卫 case、路由 case 三处今天**同为一个六枚闭集**。

### 4.2 钉名册的用例（逐枚 file:line，本腿现量；**7 处钉身**）

| 钉 | 锚 | 判据原句（逐字/紧缩） | 数什么 |
|---|---|---|---|
| **G1 真树 want-4** | `internal/panel/git_test.go:387`（用例 `TestGitDimensionHasNoModelCallableTool` `:364` 起；want 名册 `:383-385` 只列四枚 `panel.*`） | `if got := whitelistMethodsFromSource(t, filepath.Join(root, "internal", "panel", "bridge.go")); !equalStrings(got, wantMethods) {` → 红句 `:388-390` 含逐字 `"the C17 panel.* whitelist in bridge.go = %v, want the four methods this ticket may not extend - … a fifth method means the read half got tied back to the action half"` | `panelMethodRe`（`:396`）正则**写死 `panel.` 前缀**、扫 `bridge.go` 全文 ⇒ 真树第 5 枚 `panel.*` 字面即红 |
| **G2 正控 want-5** | `internal/panel/git_test.go:516`（`TestPlantedGitToolShapesGoRed` `:486` 起，栽的五行树 `:506-511`） | `if n := len(whitelistMethodsFromSource(t, fivePath)); n != 5 {` | 同一提取器对栽入的第五枚 `panel.worktree.switch` 必须数到 5 ⇒ 钉死提取器本身不失明 |
| **G3 真树 len=4** | `internal/panel/git_test.go:519`（同上用例） | `if real := whitelistMethodsFromSource(t, filepath.Join(panelRepoRoot(t), "internal", "panel", "bridge.go")); len(real) != 4 {` | 同 G1 的数字形：真 `bridge.go` 的 `panel.*` 字面**恰 4** |
| **R1 闭集名册·DECLARED 分母** | `internal/panel/inbound_roster_253_test.go:435`（`TestFullInboundMethodRosterIsClosed` `:419` 起；名册 var `:69-76`、常量 `:59-64`） | `missing, extra := sameNameSet(declared, wantFullInboundRoster253)` → 红句 `:437` 含 `"a new inbound route entered the tree without being registered anywhere"` | `bridge.go` 全部包级 dotted 常量（**前缀盲**，第 5/6 枚无前缀也在内）对**写死的 6 枚字符串名册**双向比对 |
| **R2 枚数常量×2** | `inbound_roster_253_test.go:445`（roster 枚数 6）＋`:455`（无前缀枚数 2） | `:445` `if len(wantFullInboundRoster253) != wantFullInboundRosterSize253 {`；`:455` `if len(unprefixed) != wantNonPanelPrefixedInbound253 {` | 两个写死数：**roster=6**／**无前缀=2**；红句 `:456` 具名"that batch is exactly what AC#2 exists to cover" ⇒ 这枚是**只有 `config.*` 形状的新门**会撞的，`panel.*` 形状的新门不撞它 |
| **R3 RUNNING 守卫分母** | `inbound_roster_253_test.go:464` | `if !knownComposerMethod(name) {` → 红句 `:465` 含 `"a rostered inbound name Go does not answer is a page-visible dead door"` | 名册逐枚问真守卫 |
| **R4 守卫 AST 分母** | `inbound_roster_253_test.go:470`（取 labels）＋`:474`（`sameNameSet(guardLabels, wantFullInboundRoster253)`，红句 `:476`/`:479`） | `switchCaseValues(t, bridgePath, "knownComposerMethod", …)` | 守卫 case 表（经常量解析）对写死名册双向 |
| **R5 路由分母** | `inbound_roster_253_test.go:494`（取 labels）＋`:498`（比对，红句 `:500`/`:503`） | `switchCaseValues(t, dispatchPath, "dispatch", "rosterMismatch", …)` | `composer_dispatch.go` 六枚 case 对同一写死名册双向 |
| **R6 ANSWERED 字面分母** | `inbound_roster_253_test.go:514`（`collectAnsweredOffRoster` 后逐名 `t.Errorf`，红句 `:515-516`） | 红句含逐字 `"the running guard answers the dotted literal %q found at %s, and it is on no roster - a route that exists only inside an expression has no registerable spelling"` | 包内**任何产码点号字面**被真守卫答了就必须在名册上 |
| **C1 闭合词汇表（页面侧）** | `internal/panel/composer_test.go:409-419`（`composerRouteLiterals()`，5 枚：四常量＋写死 `"panel.approval.request"`） | 消费点 `TestTheRendererHoldsExactlyOneDoorToTheHost` `:502`：`(iv)` 在 `:521-523` `if len(rep.unknown) != 0 { t.Errorf("the renderer names a route the Go side does not answer: …")}` | `frontend/src` 全树 `"panel.*"` 字面必须 ∈ 此 5 枚集（**`routeLiteralRe` `:394` 前缀盲点同 G1**） |
| **L1 Method\* 常量钉（冻结件）** | `internal/panel/l2_grant_boundary_test.go:1293-1303`（`TestAnsweredPanelRoutesCarryNoApprovalDecision` `:1229`；红句本体 `:1301`） | `:1294` `for name, val := range pkg.consts {`→`:1295` `if strings.HasPrefix(name, "Method") {`→`:1299-1301` `if !declared[name] { t.Errorf("Go answers %q but no Method* constant declares it: … a fifth route that never became a constant is invisible to it. This line is the only gate on that shape", name)` | 守卫答的每枚名必须有一枚 `Method*` 常量声明；红句自认 `bridge_test.go:131` 只循环四枚常量、对第五枚失明 |
| **L2 词表钉×2（冻结件）** | 路由词表 `l2_grant_boundary_test.go:1243`（`grantRouteWords` 定义 `:192-195`，11 枚：approve/approval/grant/allow/permit/ratify/authorize/authorised/decide/decision/verdict）＋字段词表 `:1189`（`grantFieldWords` 定义 `:183-188`，19 枚） | `:1244` 红句 `"Go answers the inbound route %q, whose own name is an approval decision (D33/F2, R20, AGENTS.md §1.2 ban #6)"` | 第五枚方法名／字段名不许含这些词根（归一化在 `carriesGrantWord`：小写、去 `_ - 空格 .` 后 `strings.Contains`） |
| **L3 种子对账钉（冻结件）** | `l2_grant_boundary_test.go:1857-1861`（`TestJSONKeyDerivationAgreesWithEncodingJSON`；registry 名册 4 枚在 `:1721-1726`） | `t.Fatalf("decode destination %s has no reflection twin in inboundTypeRegistry: … the AST half would be unchecked", seed)` | `internal/panel` 任何 JSON decode 落到的同包 struct 自动成 seed，必须有 reflection 孪生 ⇒ **新增入向承载 struct＝FATAL 且修法在冻结件内**〔a5 §3.2 已量同枚，本腿复跑锚一致〕 |

### 4.3 "第五枚入向方法名"的代价（射程判断，⛔ 非颜色）

设落地腿要在 `bridge.go` 新增一枚入向方法名（无论它叫什么），逐枚代价：

| # | 会撞的钉 | 撞法 | 修法落在哪 |
|---|---|---|---|
| 1 | **G1＋G3**（若新名带 `panel.` 前缀） | `panelMethodRe` 数到 5，红句自己写着"want the four methods this ticket may not extend" | 修法＝改 `git_test.go:383-385` 的 want 名册＋`:519` 的 `!= 4` ⇒ **普通用例，但要解冻 C17 的"四枚"字面** |
| 2 | **R1–R5 全族**（任何形状的新名） | 名册 var `:69-76` 与两枚常量 `:59/:64` 必须同步动 | 修法＝`inbound_roster_253_test.go` 的名册与枚数 ⇒ **票 253 AC#2 的闭集名册，属人工批准面**（名册注释 `:37-40` 自述"written down here as string literals, not as references to the Go constants: a roster that re-lists the constants it is auditing stays green no matter what gets added"） |
| 3 | **R6**（若新路由只写进守卫 case/表达式） | 包内点号字面被守卫答了但不在名册 ⇒ 红 | 与 2 同一枚的另一半 |
| 4 | **L1**（冻结件） | 守卫答了但无 `Method*` 常量 ⇒ `:1301` 红；红句自认它是"the only gate on that shape" | 修法＝**加 `Method*` 常量**（这条不用改冻结件——常量在产码 `bridge.go` 里加即可）；但若新路由**不加常量**则无解 |
| 5 | **L2**（冻结件） | 名含 11 枚路由词根／字段名含 19 枚字段词根 ⇒ 红 | 取名绕开即可（`stop`/`cancel`/`draft`/`append` **不在表内**，本腿逐枚比过〔a5 §3.2 同裁，本腿复跑成立〕） |
| 6 | **页面侧三发**（`bridge_test.go:142` 逐枚 contains／`composer_test.go:509-523` 四发结构钉＋`:521` unknown 门） | 新名必须由页面发出且 ∈ `composerRouteLiterals()` 5 枚集 ⇒ **页面与 Go 必须同批** | 修法＝`composer_test.go:409-419` 词汇表＋`frontend/src/lib/panel.ts`（⛔ 落地腿禁写 frontend/**＝票面硬约束 3 ⇒ **这一发落地腿自己修不了**） |
| 7 | **L3**（冻结件，仅当第五枚是新 struct 承载） | 新 decode struct＝新 seed ⇒ `:1859` FATAL | 唯一修法＝往 `inboundTypeRegistry`（`:1721-1726`）加行 ⇒ **改冻结件一字＝禁** ⇒ **必须复用 `ComposerRequest`**（已在 registry） |
| 8 | **派发 case（产码）** | `composer_dispatch.go:177-202` 无 case ⇒ 守卫答了路由不认（R5 红句 `:503`"the guard would answer it and the router would fall through to the roster-mismatch backstop"） | 修法＝`composer_dispatch.go` 加 case（产码，无需解冻） |

**代价小结**：第五枚入向方法名至少撞 **4 处普通用例钉（G1/G3、R1–R6 族、C1 页面词汇表、派发 case）＋ 2 枚冻结件钉（L1、L3）**；其中 L3（新 struct 形状）与页面同批（第 6 条）两发是落地腿**无法自行修**的。⇒ 这一格与 a5 §3.2 的裁语一致：**新增入向形状若为独立请求 struct，治理面死路；若复用 `ComposerRequest`，仍要动票 253 闭集名册（人工批准）＋页面侧同批（写面禁令）**。本腿不裁要不要走，只把钉量清。

### 4.4 a5 §1.3 行号核对表（本腿重取的对照）

| a5 原句 | 本腿现量 | 判 |
|---|---|---|
| `git_test.go:387`／`:519`（真树 want-4）＋`:516`（正控 want-5） | 同号 | 未漂〔本腿复跑成立〕 |
| `inbound_roster_253_test.go:445`／`:455`＋`:435`／`:474`／`:498`／`:510-517` | `:445`/`:455`/`:435`/`:474`/`:498` 同号；`:510-517` 现 `:510-517`（`collectAnsweredOffRoster` 循环＋红句）同号 | 未漂〔本腿复跑成立〕 |
| `bridge_test.go:138-145`／`:161` | 四枚常量循环在 `:138-141`、红句 `:143`；`postMessage` 计数钉 `:161`（红句 `:162-163`） | 未漂〔本腿复跑成立〕 |
| `composer_test.go:409-419` | 同号（5 枚闭合词汇表） | 未漂〔本腿复跑成立〕 |
| `l2_grant_boundary_test.go:1293-1301`／`:1857-1861` | 常量钉块 `:1293-1303`（红句本体 `:1301`）；种子对账 `:1857-1861` | 未漂〔本腿复跑成立〕 |
| （a5 未引用的补量）`config_route_248_test.go:75-76` | `if !knownComposerMethod(method) { t.Fatalf("knownComposerMethod refuses %q: the whitelist did not really grow", method) }`（`TestAC1BothSettingsNamesAreWhitelistedAndRouted` `:73`） | 新钉入册 |

## 5. 问五 — 分位数／真窗那一族的脆弱点与"该钉哪两行"

> 口径（照本件骨架自述的三问，⛔ 不含"该钉哪两行"那一问——它不在我的派单射程里，见 §7 第 4 条）：**(1)** 哪些用例带计时／聚合断言（逐枚 file:line＋判据原句）；**(2)** 哪些带 SKIP 通道；**(3)** 什么形状会"没跑起来却像绿"。
> 尺（本腿 16:0x 现跑，⛔ 零 `go`）：
> `grep -rn "t.Skipf(\|t.Skip(" cmd/wisp internal/panel --include=*_test.go | wc -l` ＝ **16**（〔a5 S30 复认成立〕；其中 **cmd/wisp 15** ＋ **internal/panel 1**＝`internal/panel/attachments_test.go:191`）
> `grep -rlE "time.Since" cmd/wisp --include=*_test.go` ＝ **3 枚文件**
> `grep -rn "t.Parallel()" cmd/wisp/panel_host_windows_test.go cmd/wisp/panel_resident_windows_test.go cmd/wisp/panel_host_gate_test.go cmd/wisp/panel_geometry_255_winlive_test.go | wc -l` ＝ **0**
> `ls cmd/wisp/*_windows_test.go | wc -l` ＝ **25**；`grep -rln winlive cmd/wisp --include=*.go | wc -l` ＝ **18**
> 锚全部本腿现取（起手 `2a633eb8`，收尾 HEAD 见 §7 第 5 条）。

### 5.1 分位数本体：`cmd/wisp/panel_host_windows_test.go`（1103 行、**4 枚**测试 `:614`／`:875`／`:976`／`:1036`）

这一族只有一枚**分位数尺**，它的样本**不是自己采的**。逐枚：

| 枚 | 锚 | 判据原句（逐字） | 它数什么 |
|---|---|---|---|
| **P0 定义体** | `cmd/wisp/panel_host_windows_test.go:589` `func nearestRankPercentile(vals []float64, q int) float64`，索引式在 `:595` 逐字 `	idx := int(math.Ceil(float64(q) / 100.0 * float64(len(cp))))`，夹逼 `:596-601`，**空输入返回 `-1`（`:590-592`）** | — | 最近秩定义（1-based on sorted copy）。注释 `:587-588` 自认"With n=10 that makes P95 the second-largest sample" ⇒ **n 不是自由度**：n=10 时 P95＝第二大样本，一枚尾部样本就决定红不红 |
| **P1 聚合门** | `cmd/wisp/panel_host_windows_test.go:1004`／`:1007` | `:1004` `	if coldP95 > 1500 {`→红句 `:1005` 逐字 `"cold P95 %.3f ms over %d runs exceeds the D32 panel cold budget 1500 ms (the single-run assertion is not the only gate: the tail is what AC#2 asks for)"`；`:1007` `	if hotP95 > 200 {`→红句 `:1008` | **两枚 P95**（cold 1500／hot 200）。⚠ **P50 算了但从不判**（`:995`/`:997` 只进 `t.Logf` `:999-1000`） |
| **P1b P11 线** | `:1003` | `	t.Logf("AC#2 P11 line (cold > 2000 ms) - max cold observed over these runs: %.3f ms", maxOf(cold))` | **纯 Logf**，注释 `:1001-1002` 明写它"reported, never a threshold of its own here" |
| **P2 单发门** | `:664`／`:790` | `:661` `	if coldMs <= 0 {`→`t.Fatalf`；`:664` `	if coldMs > 1500 {`→红句 `:665` 逐字 `"cold bring-up %.1f ms exceeds D32 panel cold budget 1500 ms"`；hot 同形 `:788`/`:790` | 真窗**单发** cold/hot。`hotMs <= 0` 那支是 `t.Errorf` 不是 Fatalf（`:788-789`）⇒ 无读数仍会继续往下 record |
| **P3 唯一喂样点** | `:793` 逐字 `	panelHostLatency.record(head, coldMs, hotMs)` | — | **全仓这一发调用只有 1 处**（尺＝`grep -rn 'panelHostLatency.record' cmd/wisp`＝命中定义 `:577`＋这一枚）⇒ **P1 的分母 100% 来自 `TestPanelHostRealWindowHopAndLifecycle` 一枚用例**，且只在它跑到 `:793` 那一行时才进料 |

**"该钉哪两行"这一问在本腿量不到判据**（骨架写了它，派单未授）；本节只报射程事实：**今天这条族里没有任何一枚用例断过 `len(cold) >= 10`**——样本枚数只出现在 `t.Logf`（`:987`/`:999`），P1 的门对 n=1 与 n=10 一视同仁。⇒ 具名后果：**"10-run P50/P95" 那枚 AC#2 字面，仪器只管住 P95 那半、管不住 10-run 那半**（本腿只登记这一事实，不裁该不该补；⛔ 补法形状是产码，非本编队）。

### 5.2 真窗那一族的脆弱点（逐枚 file:line，本腿现量）

| # | 锚 | 形状 | 脆弱在哪 |
|---|---|---|---|
| V1 | `cmd/wisp/panel_host_windows_test.go:1` 逐字 `//go:build windows`；文件头注释 `:5-11` 逐字 `"they are gated by the //go:build windows tag (a real platform fork, not a fake skip); they never use testing.Short() or an environment escape hatch to pretend green, and a window that cannot be created here is a red, not a skip"` | **平台叉，非 SKIP 叉** | 非 windows 编译时这 4 枚用例**根本不进二进制**：`go test` 报 `ok`／`no test files` 而**无红名册**。文件自己声明这是"real platform fork"——它不是伪装，但**读数面上与"没跑"同形**（⇒ §5.3 形状 A） |
| V2 | `cmd/wisp/panel_host_windows_live_test.go:1` 逐字 `//go:build windows && winlive` | **双标签叉＋无 CI** | `:5-20` 注释具名：这一枚（`:36` `TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive`，2s 退出条款，`:62` `	deadline := time.Now().Add(2 * time.Second)`＋`:70` `	if final.TreeWebview > baseline.TreeWebview {`）由编排者裁定（10-01 12:22，form B）**移出默认层**，代价逐字 `:18-20` `"winlive has NO CI job, so this clause is now measurable only on a desktop box and is never seen by the pipeline"` ⇒ 默认层全绿时这一寸**一次都没被测过** |
| V3 | `cmd/wisp/panel_geometry_255_winlive_test.go:1` 同 `windows && winlive`；真窗宽度那一族，结算在 `:143`／`:198` `if n := settleThreadWindows255(t, hh, tid, 3*time.Second); n != 0 {`，等待体 `:229` 起、`:231` `	deadline := time.Now().Add(limit)`＋`:239` `		time.Sleep(50 * time.Millisecond)` | 同一枚双层叉的第二户 | 真窗几何（票 255 那条）今天也在**无 CI 的层**里；`:221` 注释自述它"reports the number it measured rather than a guess about the settle time" ⇒ 结算不到 0 时报数不报绿 |
| V4 | `cmd/wisp/panel_host_windows_test.go:620` 逐字 `	dataPath := filepath.Join(os.TempDir(), "wisp-33r1-panel-profile")`，前置注释 `:615-619` 具名 `t.TempDir` 的 auto-RemoveAll **会**与浏览器 profile-DB 拆除抢锁 ⇒ 故意不删 | **跨 run 残留 profile** | 固定路径＋永不删 ⇒ 第二次跑的起点取决于上一次留下了什么。同文件 `:630-638` 的 `baselineHosts` 与 `:703-710` 的注释正是为这一枚而存在（`255-r5 measured exactly that shape: post-fix-count3-v2 run 2 read TWO hosts pre-hide and went red on the leftover's exit`）⇒ **判据靠"减基线"成立，基线本身是时序量** |
| V5 | `cmd/wisp/panel_host_windows_test.go:353-356` 逐字 `	treeSettleWait = 3 * time.Second`／`	treeSettleStep = 100 * time.Millisecond`／`	treeSettleMax  = int(treeSettleWait / treeSettleStep)`；循环体 `:371-383`，`:376` `		if second.TreeWebview == first.TreeWebview {`→返回 | **两枚样本一致即算"settled"** | 注释 `:346-352` 自认"detect the tree webview count drifting ... keeps sampling"，`:360-369` 自认 255-r5 实测"helpers spawned around bring-up exit within roughly a second"。⇒ 浏览器换编号的节奏与 100ms 采样**同量级**：两枚一致可能是"真稳"，也可能是"两次都落在同一段漂移里"；且**耗尽 30 枚样本后 `:382` 仍返回 `(first, treeSettleMax+1)` 而非 `t.Fatalf`** ⇒ 未结算读数被当结算读数用（`afterExtra`/`beforeExtra` 只进 `t.Logf` `:749`/`:715`） |
| V6 | `cmd/wisp/panel_host_windows_test.go:394-398`（`browserHostPids`）逐字 `		if row.ppid == root && strings.EqualFold(row.name, webviewExeName) {` | **只认直接子进程** | 注释 `:385-390` 自述理由（helper 挂在 host 下，不挂测试进程下）。⇒ 这一族对"浏览器自己重排父子关系"是**结构性失明**，与 V5 的漂移同源 |
| V7 | `cmd/wisp/panel_host_windows_test.go:611-613` 注释逐字 `// Do not parallelize this test or any sibling in this file: the AC#3 socket ruler`／`// and the AC#4 foreground ruler both read machine-global state, and`／`// TestAC3ListeningSocketRulerSeesItsOwnListener deliberately opens a listener.`；实量 `t.Parallel()` 在本族四枚文件里＝**0 命中** | **靠约定防并发，不靠仪器** | 没有 `if t.Parallel() …` 之类的钉，**这是一句注释不是一枚断言**：任何人给这个包加并行，AC#3（`:552` `			if state == mibTcpStateListen && pids[own] {`，读机器全局 TCP 表）与 AC#4（`:876` `	foregroundBefore := uintptr(windows.GetForegroundWindow())`，读机器前台窗）立刻互相污染 ⇒ 真窗族的"红"可以是**别人开窗**造成的 |
| V8 | `cmd/wisp/panel_host_windows_test.go:1021` 注释逐字 `// It is declared AFTER the lifecycle test on purpose: go test runs the tests of a`／`// file in source order, and the aggregate is read from the samples that file's`／`// earlier tests recorded.` | **跨用例·同进程·源序耦合** | P1 的输入是**包级可变全局** `:575` `var panelHostLatency latencyRecord`，进料点在另一枚用例体内（`:793`）。⇒ `-run` 单挑 P1／重排源序／P2 那枚用例在 `:793` 之前任一 `t.Fatalf`（`:641`/`:648`/`:651`/`:661`/`:728`/`:734`/`:753`）都会让聚合**空**——此时不是红，见 §5.3 形状 B |
| V9 | `cmd/wisp/panel_host_gate_test.go:338` `func gitHeadShortForTest`，`:341` `	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()`，`:343` `		return "HEAD-unknown"` | 读数自述的锚是**跑出来的** | 注释 `:333-337` 具名它就是为了终结 33-v1 §A#21（旧件把 `"HEAD 7a41db9b"` 写死在测量行里）。⇒ 值：**没有 git 元数据时读数仍绿，锚写成 `HEAD-unknown`**；任何引用这族读数的表若不带这一枚区分就是过期 |
| V10 | `cmd/wisp/panel_resident_windows_test.go:50-53` 逐字 `// panelThreadWait bounds every "did the panel thread get there" wait in this file.`／`const panelThreadWait = 15 * time.Second`；结算体 `:87` `	deadline := time.Now().Add(panelThreadWait)`→`:94` `	t.Fatalf("timed out after %v waiting for %s - this is a failed measurement, not a pass", panelThreadWait, what)` | **墙钟等待，但红向做对了** | 超时是 `t.Fatalf`（不是 skip、不是绿）——这枚形状是**本族的正面样本**；`time.Now()` 差值实现超时这一点与 `AGENTS.md §1.2` 的"不许用墙钟时间差实现超时"在**字面**相冲（测试件是否在射程归 `tools/d22scan`，⛔ 本腿不裁，登记为 Q-1） |
| V11 | `cmd/wisp/panel_resident_windows_test.go:611` 逐字 `const refusalSettleBudget = 5 * time.Second`；四发判据 `:703` `	for time.Now().Before(settled.Add(refusalSettleBudget)) {`＋`:731` `	secondElapsed := time.Since(secondStart)`→`:735` `	if secondElapsed > refusalSettleBudget {`＋`:750`/`:755` `case <-time.After(refusalSettleBudget + time.Second):` | **同一枚 5s 既当等待上限又当断言阈值** | `:703` 的轮询与 `:735` 的门是**同一个常量**：轮询若被拖慢，`secondElapsed` 的判据与它自己的等待上限同涨同落 ⇒ 机器负载升高时这一发**两侧一起让**（`settleElapsed` `:709` 只进 `:757` 的 Logf，不参与判）。这是真窗族里"阈值＝自身超时"形状的一枚 |
| V12 | `cmd/wisp/resident_approval_live_246_windows_test.go:121`/`:226`/`:315` 逐字 `		done <- askResult246{ans: a, why: w, took: time.Since(start)}` | **took 只进 FAIL 句与 Logf，从不判** | `:158-159` 的 `t.Fatalf("(b) FAIL: … after the injected Esc … window took %s", res.ans, res.why, res.took)` 是**红句的附加说明**；全文件零枚 `took` 阈值比较 ⇒ live 层的时延是**记录**不是**门**（与 P1"门是聚合"正好相反，两形并存须分清） |
| V13 | `cmd/wisp/run_test.go:388` 逐字 `	if blocked < 2*time.Second {`→红句 `:389` 逐字 `"the call returned after %v: the 2s L1 window must actually block"` | **下界型真窗断言**（非 winlive，默认层真会跑） | 这一枚的阈值**来自配置**：`cmd/wisp/run.go:615` 逐字 `			Window:          time.Duration(cfg.Risk.L1WindowSec) * time.Second,` ⇒ **改 `risk.l1_window_sec`＝静默改变这发的语义**（配置调成 1s：这发红；调成 5s：这发要真等 5s 才绿，包时长随之动）。全仓 `time.Since` 只落在 3 枚测试文件（V10/V11/V12/本枚＋`panel_resident_windows_test.go`），这是**唯一一枚既在默认层又带真墙钟下界的** |
| V14 | `cmd/wisp/subagent_selfapproval_197_test.go:196-200` 逐字注释 `"The window this case needs is … 2.5s of mockllm latency is that window"`＋`	f.srv.Control(t, "/__control/latency", `+"`"+`{"ms":2500}`"+`"+`)`；同形 `cmd/wisp/subagent_blocked_197_test.go:63`（80ms）／`subagent_carrier_197_test.go:233`（80ms）／`internal/llm/openaichat/mockllm_integ_test.go:349`（120ms） | **用注入时延撑开竞态窗** | 这些数字**是 fixture 的一部分**，不是被测参数；窗口真宽窄由机器负载决定（`:198` 的注释自认它要撑住的正是 `runTextTask closes the store the moment the root returns`）。⇒ 真窗形状**已经从 panel 渗进非 winlive 的任务族**：负载低时两事件重合、用例可仍绿而断言的那一形从未发生（形状 D） |
| V15 | `internal/llm/adaptertest/mockllm.go:51-54` 逐字 `	if testing.Short() {`→`		t.Skip("integration test (mockllm subprocess) skipped under -short")`；调用面 `cmd/wisp/run_test.go:78` 等（`adaptertest.StartMockllm(t)`） | **全仓唯一活着的 `-short` 通道在 `cmd/wisp` 的依赖里** | ⚠ **具名更正 a5**：a5 §0/S 尺记"本族 never 用 testing.Short"只对 `panel_host_windows_test.go` 的字面成立（`cmd/wisp` 里 `testing.Short()` 只出现在那枚文件 `:10` 的**注释**）；**`cmd/wisp` 的 runFixture 经 `adaptertest.StartMockllm` 真吃 `-short`** ⇒ `go test -short ./cmd/wisp/` 让所有走 mockllm 的腿**整批 skip**。这是 §5.3 形状 C 的真实入口，⛔ 不是"本包没有 escape hatch" |
| V16 | `cmd/wisp/slo_report_144_windows_test.go:10-12` 逐字 `"and this host's timings are not trustworthy anyway - see scripts/wisp-cli-tests.sh for why"` | **本族自己写下的"别信本机时延"声明** | 同文件把时延族**重写成一数读次数的尺**：`:318` `	s := &sloSubject{… readReportFile: r.read}`＋`:319` `	_, err := s.collectReportWithin(60*time.Millisecond, time.Millisecond)`＋`:323` `	if r.calls <= 3 {`，`:549-551` 注释逐字 `"Nothing here touches a budget or a threshold: subjectReportBudget and subjectGrace are not referenced, and the waits handed to the loop are parameters"` ⇒ **真窗那一族的正解形状已经在仓里**（把时间当参数、把"等了几次"当判据）。⚠ 这条对 167-r2 有用：占用／序号那一跳若要断"没等够就不许显示"，**仓内已有不靠墙钟的写法先例** |

### 5.3 什么形状会"没跑起来却像绿"（逐形具名，判据行号＝本腿现量）

| 形 | 机制 | 盘上的证据 | 与真绿的可区分面 |
|---|---|---|---|
| **A｜平台叉** | `//go:build windows` 让整枚文件在非 windows 下**不入编译** | `cmd/wisp/panel_host_windows_test.go:1`／`panel_resident_windows_test.go:1`／25 枚 `*_windows_test.go` | 只能靠**跑的平台**区分；`go test` 输出里既无红名册也无 skip 名单。`panel_host_windows_test.go:9` 自称"real platform fork, not a fake skip"＝承认形状，未提供区分手段 |
| **B｜样本饥饿→具名 skip** | P1 的输入是另一枚用例的全局进料；进料缺 ⇒ 聚合空 ⇒ `t.Skipf` | `:984` `	if len(cold) == 0 {`→`:985` Skipf（句中含 `Named skip - an empty aggregate is not a green latency gate`） | **这形做得对**：skip 有名字、句里写清了触发条件（`-run` 只挑一枚／生命周期用例没跑）。但**报告侧若只报 PASS/FAIL 不报 SKIP，它就跟绿一样** ⇒ 引用这族的表必须带 skip 枚数（本腿无 go 命令、拿不到实跑 skip 枚数，归 §7 第 3 条） |
| **C｜`-short` 整批隐身** | `adaptertest.StartMockllm` 内 `t.Skip` ⇒ `cmd/wisp` 里**所有 runFixture 腿**一起 skip | `internal/llm/adaptertest/mockllm.go:51-54`（定义）＋`cmd/wisp/run_test.go:78`（调用面） | 与 B 同族但**面积大得多**（一枚 helper 管一片腿），且 skip 句写的是"integration test"——**不会告诉读的人 `cmd/wisp` 今天少跑了多少枚** |
| **D｜窗口没撑开** | 注入时延（V14）或真窗（V4/V5）不足以造出被测并发窗，但两枚事件恰好仍给出"对"的答案 | `cmd/wisp/subagent_selfapproval_197_test.go:200`／`cmd/wisp/panel_host_windows_test.go:371-383`（settle 未耗尽时也返回同一类型值） | **最难分辨**：判据成立、路径未发生。V16 那枚 `slo_report_144` 的"数次数不数秒"是唯一在册的反例写法 |
| **E｜仪器失明→`t.Fatalf`** | 尺读不到东西时**红**而不是绿——本族的正面形状 | `:753`／`:558`（`GetExtendedTcpTable` 六次失败→`listening count UNVERIFIED`）／`panel_resident_windows_test.go:94`／`slo_report_144` 的 `r.calls` 下界 | 可区分，且是**该被抄的那一形** |
| **F｜空集恒真（相邻形状，本族外但同包）** | 双向对账尺若两侧都空则零差分通过 | `cmd/wisp/panel_host_gate_test.go:153` `		t.Logf("AC#12 provenance axis not measurable in this tree (no git metadata - e.g. a git-archive copy); the capability assertions above still ran")` | ⚠ 这一支是 **Logf 不是 Fatal**：无 git 元数据时 AC#12 那一轴**报"没测到"后继续**，同文件 `:181-182`（`len(ctorHits)==0`→`t.Skipf`）同族。a5 已在族①/② 登记过同类仪器债（K4 的枚数形、tag 盲区），本腿只把 `cmd/wisp` 侧的这两枚补进同一张表 |

### 5.4 SKIP 通道名册（16 枚逐枚，本腿现量；格式 `file:line`｜触发条件｜若被误读成绿会漏掉哪一寸）

`cmd/wisp`（15 枚）：
1. `panel_host_gate_test.go:182`｜`len(ctorHits)==0`（包内无非测试 `NewPanelManager` 调用点）｜漏掉 AC#1 第二条款"session dispose 真的拆窗"
2. `panel_host_windows_test.go:985`｜`len(cold)==0`（聚合空）｜漏掉 **P1 那两枚 P95 门**（分位数从未判）
3. `panel_resident_windows_test.go:149`｜尺自己造不出前置窗（"failed measurement, not a product verdict"）｜AC#4 前台归还整发无主语
4. `panel_resident_windows_test.go:317`｜`len(probes)==0`（embed 解析不出 entry）｜**AC#13 冷启动落在嵌入页那一整发**
5. `panel_resident_windows_test.go:999`｜`first==0`（本机前台没记到 prior）｜AC#4"prior 挺过 Hide"那一形
6–9. `resident_hotkey_258_windows_test.go:280`／`:343`／`:386`／`:400`｜`SKIP-LOUD`：本机拉不起球窗／Ctrl+Alt+V 被别家占着｜热键重绑与"真占用"那两条路径
10–13. `resident_hotkey_live_258_windows_test.go:55`／`:97`／`:142`／`:219`｜`SKIP-LOUD`：live 层无球窗｜live 热键四发
14. `secret_dataroot_119b_test.go:73`｜造不出符号链接｜"链接逃 dataroot"那一形
15. `secret_dataroot_119b_test.go:234`｜data root 内无符号链接｜同上反向
`internal/panel`（1 枚）：
16. `attachments_test.go:191`｜文件系统放不下 2 GB（`fh.Truncate(2<<30)` 失败）｜大附件计数那一发

⚠ **枚数是 16、面积不是 16**：第 6–13 这八枚集中在**同两枚文件**（热键 258 族），而 `panel_host_windows_live_test.go`／`panel_geometry_255_winlive_test.go` **一枚 skip 都没有**——它们靠 build tag 隐身（形状 A/V2），**porcelain 上看不见、skip 名单上也看不见**。⇒ 只数 `t.Skip*` 的尺会把 winlive 层读成"零逃逸通道"，这是本腿认为最该具名的一枚读数陷阱。

另记两枚**注释里的 skip 声明**（尺 `grep -rn 't\.Skip' cmd/wisp --include=*_test.go` 的 raw 与去注释读数**相同＝16**，因这两枚写的是 `NO t.Skip`／`No t.Skip`）：`leg_sink_nail_131_windows_test.go:52`／`resident_sink_nail_127_windows_test.go:32`——两枚"钉"类用例自我声明**禁用 skip 通道**，与 V15 的 `-short` 后门同包并存，引用本包"不许逃逸"时须带口径。

### 5.5 本节小结（两句话，给 §6 当料）

1. 真窗族**只有 1 枚分位数尺**（P1），它的**分子是 P95 阈值、分母是另一枚用例的进料**；样本枚数（n）今天无钉，P50 今天不判 ⇒ "10-run P50/P95" 那句 AC#2 的字面只有半句有仪器。
2. 逃逸面**不在 skip 名单上，在 build tag 上**：默认层全绿≠winlive 那两寸（2s 退出、真窗宽度）被测过；`-short` 那枚后门住在 `cmd/wisp` 的依赖里（V15），不在 `cmd/wisp` 的字面里。

## 6. 问六 — 给落地腿 167-r2 的一页清单

> 综合＝a5 的 §0–§3〔转录，逐枚标〔a5 已量〕〕＋本腿 §4（问四）＋本腿 §5（问五）。
> ⛔ **本清单不裁修法形状**：只给"会撞哪一枚／它今天判什么／撞法是哪一支"。颜色一律是射程判断。
> ★**两个动作先分清**（凭据＝票面 §9 第 3 条末行与 `Q-51`/`Q-76` 在册口径，本腿复认）：
> **动作甲＝新增出向字段**（占用／序号／cancellable 三枚"看得见"），**167-r2 的写面就是这一发**；
> **动作乙＝新增入向那一跳**（第五枚方法名＝"按得下去"）＝`Q-76` 未批支，**167-r2 不该做到**（票 §9 第 3 条逐字"未答期间本票只做'看得见'那一半"）。
> 本腿对每一枚锚**全部重取**（锚 `2a633eb8`→本腿收尾 HEAD，见 §7 第 5 条）；漂移逐条列 §7 第 2 条对照表。

### 6.1 动作甲｜新增出向字段——按落点分四格，逐枚撞钉

**格 A：`ApprovalCardView` 加一枚键（＝序号 `position` 的落点，票 §9 第 1 条已答"必须新增一枚 JSON 键"）**

| 枚 | 锚（本腿现量） | 判据原句 | 判定 |
|---|---|---|---|
| B-1 尺 B 正向 | `internal/panel/approval_test.go:128` 逐字 `		if missing := subtract(goKeys, tsKeys); len(missing) > 0 {`→红句 `:129` `"Go %s emits %v that interface %s does not declare"` | **会红**（若 Go 发了、`panel.ts` 没声明） |
| B-2 尺 B 反向 | `internal/panel/approval_test.go:131` 逐字 `		if extra := subtract(tsKeys, goKeys); len(extra) > 0 {`→红句 `"interface %s reads %v that Go %s never sends - those fields render as undefined"` | **会红**（页面先声明、Go 还没发的那一批） |
| B-3 输入面 | `internal/panel/approval_test.go:107` 逐字 `	data, err := os.ReadFile(filepath.Join(root, "frontend", "src", "lib", "panel.ts"))` | **不碰**（Go 侧读页面文件，落地腿⛔ 不写 `frontend/**`＝票面硬约束 3）⇒ **B-1/B-2 的两方向里必有一方向页面欠一发**，这一发落地腿修不了，只能同批报回（与 §4.3 第 6 条同一发火源，两动作共担） |
| B-4 尺 A | `internal/panel/composer_test.go:48`（`TestComposerContractTypesMatchFrontend`），`pairs` 六对里**不含 `ApprovalCardView`**（`:60-65`＝Snapshot/ComposerState/ModeView/WorkspaceView/AttachmentRef/ResultChunk〔a5 已量，本腿复认 `:60` 与 `:73`〕） | **不碰**（票 §9 第 1 条"冲突面是一把尺不是两把"复认成立：本腿现量尺 A 六对、尺 B 三对 `approval_test.go:118-120`） |
| B-5 词表钉 | `internal/panel/l2_grant_boundary_test.go:1181` `	inbound := inboundEnvelopes(pkg)`＋`:1189` `			if carriesGrantWord(f.JSONKey, grantFieldWords) {` | **不碰**（`ApprovalCardView` 是纯出向、不进 `inboundSeeds`）。★旁证复认（a5 §3.2 末行）：`TaskRowView.BlockedOnApproval` 的键含 `approval` 而今天不红 ⇒ "出向＝豁免"是射程事实 |
| B-6 键数写死钉 | `internal/panel/subagent_stream_197_test.go:140` 逐字 `		if len(r) != 3 {`（红句 `:141` `"a result row carries %d JSON keys, want the three ResultChunk has"`） | **不碰**——**但它是"落点选错"的那把刀**：序号若挂 `ResultChunk`（`:95-99` 现量 3 枚键）必红，还同时撞 `:136` `	if len(results) != len(agents) {`〔a5 §1.1 末段已具名同一条，本腿复认〕 |
| B-7 tag 盲区 | `internal/panel/approval_test.go:182` 逐字 `		if tag := f.Tag.Get("json"); tag != "" && tag != "-" {` | **会变钝**（不是会红）：`jsonKeysOf` 只收带 tag 的字段 ⇒ 无 tag 的导出字段 B-1/B-2 两方向都看不见，而线上照样多一枚 PascalCase 键 ⇒ **"B 尺没响"不等于"没撞契约"**。票 §9 第 1 条已把这枚登记为在册盲区，本腿复认成立 |
| B-8 `omitempty` 那一支 | 产码面：`internal/panel/approval.go:39` `type ApprovalCardView struct {`，现量 **10 枚 `json:` 键、零枚整数槽**（票 §9 第 1 条同读数，本腿复认） | **不碰任何现有钉**，但它是 AC#2 退化形状的位置：`encoding/json` 对**无 `omitempty` 的 int 恒发 0** ⇒ 那条 CLI 路会自己造出票面点名的 `"position":0`。⚠ **`cardView143` 用裸 `Unmarshal` 抓不到**（见 §6.3 第 1 枚） |

**格 B：占用那一枚挂 `Snapshot` 顶层（新节／新键）**

| 枚 | 锚 | 判定 |
|---|---|---|
| K1–K3 三枚键集钉 | `internal/panel/pump_test.go:123`／`internal/panel/pump_test.go:291`／`internal/panel/subagent_roster_197_test.go:214`（红句 `:215`；反向钉 `:207` `	if _, ok := wire["tasks"]; ok {`）〔a5 §1.1 已量，本腿逐枚复认〕 | **有条件会红**：三枚比的是**名册内容**（逐字 `"composer,generatedAt,pending,results"`）⇒ 新节走 **pointer＋omitempty 且该 fixture 的 reader 未设**＝不红；走**值类型节**＝必红 |
| K4 枚数钉 | `internal/panel/subagent_stream_197_test.go:129` 逐字 `	if len(generic) != 4 {` | **会红（最脆的一枚）**〔a5 §1.1 ★K4 已具名"a3 漏数的一枚"，本腿复认〕：它比**枚数**，不比名册内容 ⇒ 顶层只要真多出一枚**这次就有值**的键即红 |
| K5 `tasks` 先例 | `internal/panel/composer.go:82-83` 逐字 `	// same reason: a pump assembled WITHOUT a roster reader sends four keys, which`／`	// is the shape those nails pin`（字段本体 `:74` `	Instructions *InstructionsSection \`json:"instructions,omitempty"\`` 与 `:91` `	Tasks *TaskRosterSection \`json:"tasks,omitempty"\``，"诚实编码"那一句在 `:69`） | **不碰**，但它是**先例**：仓里已有"pointer＋omitempty 让四枚键钉继续绿"的两发落地（票 200／票 197）。⚠ 同一段注释（`composer.go:86-90`）自己写着加 `tasks` 是 contract move、两把尺会点名"interface 未声明"（`:89-90` 逐字 `	// declare, which is that ruler doing the job ticket 77 pinned it for, and the`／`	// page-side declaration is the UI leg's to write (Q-51, ledger A383).`）⇒ **格 A 的 B-1/B-2 在格 B 同样成立** |
| 尺 A/B | `composer_test.go:48`／`approval_test.go:120`（`{"Snapshot", Snapshot{}, "PanelSnapshot"}`） | **会红**（顶层新节两把尺都管：`Snapshot` 同时在两对里〔a5 §1.2 已裁"Snapshot 顶层＝A＋B＋族①那 4 枚＝6 发"，本腿复认〕） |

**格 C：cancellable／"能不能停"那一枚挂 `TaskRowView`**

| 枚 | 锚 | 判定 |
|---|---|---|
| C-1 两把尺全盲 | `internal/panel/subagent_roster_197.go:106` `type TaskRowView struct {`（现量 12 枚键，`:115` `	StatusKnown bool   \`json:"statusKnown"\``）〔a5 §2.3 与 §1.2 末段已量，本腿复认锚号未漂〕 | **不碰**（`TaskRowView` 不在尺 A 的六对、也不在尺 B 的三对里）⇒ **这一格的"看得见"没有任何一枚现有键名册钉管它**，落地腿要 AC#4 的凭据只能自带正向种＋反向种 |
| C-2 20 行钉 | `internal/panel/subagent_roster_197_test.go:275` 逐字 `	if len(sect.Rows) != 20 {` | **不碰**（数的是 D43 名册行枚数，不是键） |
| C-3 词表钉 | `l2_grant_boundary_test.go:1189` | **不碰**（纯出向）；⚠ **取名仍要避开 `allow/grant/approve` 一族**：`cancellable` 不在 19 枚字段词根内，本腿逐枚比过（同 §4.3 第 5 条） |
| C-4 那一维已占用 | `internal/panel/subagent_roster_197.go:111-113` 注释自述 `statusKnown` 承载 D43 状态名 | **不碰**＝**不许顺手改**（票面 §9"'看得见'缺位三处"第 3 条：那一维被票 196／`A394` 占着） |

**格 D：占用那一枚的"未知"状态（AC#2 的具名形状）**

| 枚 | 锚 | 判定 |
|---|---|---|
| D-1 今天零断言 | 〔a5 §2.1 已量〕两把尺互相否证：occupancy 族 grep **3 命中全为热键占用**（`cmd/wisp/resident_hotkey_live_258_windows_test.go:249`／`:259`／`:284`），`ContextWindow` **4 命中全为配置写入**（`cmd/wisp/firstrun_257_test.go:114-115`、`internal/panel/config_route_248_test.go:272`／`:296`） | **不碰任何现有钉**。本腿逐枚复认：`:272` 现量 `			FieldModelContextWindow, FieldModelPriceIn, FieldModelPriceOut,`、`:296` 现量 `			"configField":    FieldModelContextWindow,`（a5 的包路径 `internal/panel/` **本来就写对了**，本腿复跑只为钉行号）；另复认同文件 `:73` `func TestAC1BothSettingsNamesAreWhitelistedAndRouted`／`:75` `			if !knownComposerMethod(method) {` 与 §4.4 末行一致 |
| D-2 数据源不许动 | 票面硬约束 2（`D32` 那两个数一字节不许动）＋§5.1 的 `:1004`/`:1007` 那两枚 1500/200 门 | **不碰**——⚠ 但**别把分位数族的 1500/200 当占用条的分母**：那两枚是 **panel 冷启/热重显预算**（`cmd/wisp/panel_host_windows_test.go:665`/`:791` 红句自称 D32 panel budget），与上下文窗口无关。本腿具名这一处**易混的数** |

### 6.2 动作乙｜新增入向那一跳（`Q-76` 未批支，列此仅供批后取用）

**本腿不重列**——代价表已在**本件 §4.3**（8 行，逐枚 file:line＋撞法＋修法落点），名册与钉身在 §4.1/§4.2（现量 6 枚名册／7 处钉身）。★与 §5 的接缝只有两条，列在这里补 §4 之不足：

1. **真窗／分位数那一族对入向名册零射程**（尺＝`grep -c 'knownComposerMethod' cmd/wisp/panel_host_windows_test.go cmd/wisp/panel_resident_windows_test.go cmd/wisp/panel_host_windows_live_test.go`＝**0／0／0**，本腿现量）⇒ 第五枚方法名不会撞 §5 的任何一枚。⚠ **但"零射程"不等于"没有字面"**——本腿同场量到 4 枚**引号里的 `panel.mode.request` 信封字面**：`cmd/wisp/panel_host_windows_test.go:669`／`cmd/wisp/panel_resident_windows_test.go:251`／`:790`／`:794`（后三枚走的是"真腿送真信封"那一形）。它们**只认既有名、不数名册枚数** ⇒ 新增名不撞它们；**只有把 `panel.mode.request` 改名／删掉**才撞。这一枚区分具名在册，免得下一腿把"零射程"读成"随便动"。
2. **`Message: nil` 那一发是"接法"不是"新增名"**：`cmd/wisp/panel_inbound.go:277` 逐字 `		Message:    nil,`〔a5 §2.4 已量，本腿复认锚未漂，另复认 `:275` `		Workspace:  nil,`〕⇒ 把 `Message` 接上真处理器**不新增名册枚数**、不撞 §4.2 的 7 处钉；它撞的是 §6.3 第 2 枚那一发**语义**。

### 6.3 静默过期·仪器不响但必须同批改写（落地腿最容易白撞的一族）

| # | 锚 | 现在写着什么 | 哪一动会把它变假话 |
|---|---|---|---|
| 1 | `cmd/wisp/panel_assets_143_test.go:52` `type cardView143 struct {`（**手抄 8 枚卡面键**），解码在 `:112` 逐字 `	if err := json.Unmarshal([]byte(stdout), &card); err != nil {`〔a5 §1.6 W-3，本腿复认〕 | 注释 `:49-51` 自述这张册"re-declares the printed keys" | **格 A**：卡面加 `position` ⇒ 裸 `Unmarshal` 对多余键沉默 ⇒ **零枚变红、名册少抄一枚**。⇒ 交件必须具名说"这枚手抄册未同步"或同批改它（它是测试件、非冻结件） |
| 2 | `cmd/wisp/panel_config_248_test.go:460` 用例名 `TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket`＋其 `:467` 三枚按名拒绝循环〔a5 §2.4 第一行，本腿复认锚号〕＋产码注释 `cmd/wisp/panel_inbound.go:273-274`（"the three doors nobody owns still refuse by name"/"message = ticket 35"） | 正向断言"message 那扇门必须仍按名拒绝" | **§6.2 第 2 条**：接上 `Message` 后**注释与用例名成假话**（该用例自造 dispatch、不读生产装配 ⇒ 不红）。⚠ 票面 §9 第 3 条已具名"这一格 AC#5 不属契约面"，本腿复认其仪器不响 |
| 3 | `internal/panel/approval_test.go:182`（tag 盲区） | 两把尺只认 json tag | **格 A/格 B** 全部 ⇒ **见 §6.1 格 A 的 B-7**：无 tag 字段两把尺全盲 |
| 4 | `internal/panel/subagent_stream_197_test.go:129` 的 `!= 4` 与 `pump.go:45-48` 自述的"pins are four"（逐字 `// The pins are four: adding a key also`／`// reddens approval_test.go:105 and the two byte-level key-set nails in`／`// pump_test.go (:111-124, :270-276)`） | 生产注释自报**四枚**钉 | **格 B**：真树第五枚顶层有值键会红 K4，而 `pump.go` 的"four"字样**不红只过期**；反过来若有人**新增**第五枚键钉，那段注释就成了假点名 |
| 5 | 〔a5 §3 末段引用在册禁令〕`.scratch/wisp/dispatches/2026-09-26-093x-thaw-panel-for-145.md:19` 逐字 `internal/panel/** 里其余文件若与上面三枚同名族（tokens*、l2_grant*、frontend_hygiene*），一律按禁改面处理，不确定就停手报回` | 三枚冻结件的同名族扩张令 | **落地腿若新建 `internal/panel/tokens_167_test.go` 一类文件＝撞这一令**（属停手报回格，不属能写格）〔a5 §3 已量，本腿只引不复核派单件〕 |
| 6 | `internal/panel/composer_test.go:268` 的 `gitSwitchCapabilityRe`＋`:317` 逐行扫 `internal/panel` 全部非测试 `.go`〔a5 §1.6 W-1〕；`composer_test.go:332` 的 `modeSetterRe`＋`:180`〔a5 W-2〕 | 词面钉，**注释不豁免、新文件自动进射程** | **格 A/格 B/格 C 全部**：在 `internal/panel/*.go`（含注释）写出 `worktree`／`git checkout`／`git switch`／`branchSelect`／`git.branch`／`vcs.switch`，或写出 `SetMode`／`PermissionMode =`／`perm.Store.Set(` ⇒ **当场会红** |
| 7 | `internal/panel/instructions_200_test.go:300` 逐字 ``		`"instructions"`, `"tier":"project"`, `"depth":0`, `"bytes":`,`` | 项目说明那一节已占 `depth` 这枚键名 | **格 A（序号）**：序号若拼 `depth` ⇒ 同一份字节里两枚同名键（判据是 `strings.Contains` ⇒ **不红**，但语义撞车）〔a5 §1.6 W-4，本腿复认锚与写法〕 |

### 6.4 零断言·无需改写的三格（给落地腿省工，凭据在 a5）

- **占用**：**零枚现有断言**〔a5 §2.1 裁〕⇒ 只欠自己的正向种＋反向种（AC#2 那两发变异）。
- **停止（"看得见"那一半）**：`AttachCancel` 在 `cmd/wisp`＋`internal/panel` 测试里 **0 命中**、`cancellable` **3 命中全是注释**〔a5 §2.3 尺与读数〕⇒ **零枚要改写**。
- **崩溃自救 AC#6 三枚**：`grep -rnE "MaxLoggedString|BuildDiagnosticsBundle|Panic" cmd/wisp internal/panel --include=*_test.go`＝**0 行**、`BuildDiagnosticsBundle` 产码调用者 **0 枚**、doctor 检查项零枚读 `logs\`〔a5 §2.5；★"10 枚检查项"那个数票面 §9 第 2 条已具名更正为 **28**，本腿引数带口径〕⇒ 新增自救出口**不打红任何现有用例**，唯一相邻火源＝`cmd/wisp` 的 leg/sink 门两族（数的是分发腿名册，新增 `wisp` 子命令才要过）。
- ⚠ **§5 对这三格补一条射程事实**：AC#6 的"可复制现场"若将来要**报时延**（例如"bundle 在 N ms 内可复制"），§5 的名册会立刻把它拉进**唯一默认层带真墙钟下界**那一族（V13 同形），而 `slo_report_144_windows_test.go:10-12` 已写下"this host's timings are not trustworthy anyway"⇒ **同包自己的意见是别加墙钟门**；正解先例在同一文件 `:318-325`（数读次数、把 waits 当参数）。

### 6.5 本节自证

本表每一枚 `file:line` 都是**本腿**用 `sed -n '<l>p' <file>` 逐枚打印原文对上的（⛔ 未承袭 a5/a5b 笔 2 的行号）；不一致的只有 §7 第 2 条列出的那几枚。**判定口径**：会红＝该尺的字面射程覆盖这一动；会变钝＝它不响但它的名册/注释在这一动之后失真；不碰＝尺不在这一动的路径上（含"管得到但管不到这一枚字段/名字"）。⛔ 无一枚给出"该怎么修"。

## 7. 量不到／判不动／我写错的读数（自我对抗）／⛔ 未动产码自证

本节打算答：(1) 本腿对 a5 读数复核中发现的写错处（逐条具名）；(2) 本腿自己量不动、归编排者的条目；(3) `git status --porcelain` scoped 读数收尾自证（起手在 §0）。
