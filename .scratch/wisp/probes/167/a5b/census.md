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

本节打算答：(1) `cmd/wisp` 的 panel_host／latency 族里哪些用例在断分位数或真窗（逐枚 file:line＋原句，只读码）；(2) 脆弱点（墙钟依赖／样本数写死／窗口边界语义）；(3) "该钉哪两行"——落地腿新增那一跳时最该补钉的两行断言。⛔ 不跑任何 go 命令，全部读码得。

## 6. 问六 — 给落地腿 167-r2 的一页清单

本节打算答：综合 a5 §0–§3（〔a5 已量〕转录）与本腿 §4–§5，按"新增字段／新增那一跳"两个动作各列一张**撞红钉名册**——逐枚 file:line＋判据原句，另附"静默过期不响但必须同批改写"清单与"零断言无需改写"清单。

## 7. 量不到／判不动／我写错的读数（自我对抗）／⛔ 未动产码自证

本节打算答：(1) 本腿对 a5 读数复核中发现的写错处（逐条具名）；(2) 本腿自己量不动、归编排者的条目；(3) `git status --porcelain` scoped 读数收尾自证（起手在 §0）。
