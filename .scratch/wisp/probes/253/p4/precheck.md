# 253-p4 — 只读预检：票 253「AC#2 要不要人工批准」＋「AC#1 能力尺落点零件」两问

- 时刻（起手锚）：`2026-10-03T09:12:06+08:00`；HEAD `3736f0dd`（`dev`）；`git status --porcelain` = 397 行（他腿脏面，本腿零产码）。
- 性质：只读预检腿。⛔ 零 Go 命令（`go build/test/vet` 零发，三枚写腿在飞）；⛔ 票面 AC 勾选框零碰；⛔ 产码零改；⛔ `frontend/**`／`design/**` 零读零引；⛔ 冻结件零触碰。
- 写点唯一 = 本文件（`.scratch/wisp/probes/253/p4/precheck.md`）。
- 本腿只答两问（编排者据此决定派不派落地腿），⛔ 不裁修法选形、⛔ 不据 r3 排除任何候选形。

---

## §0 锚

- 起手同发三取：`date -Iseconds` = `2026-10-03T09:12:06+08:00`；`git log -1 --format=%h` = `3736f0dd`（分支 `dev`）；`git status --porcelain | wc -l` = `397`。
- 票面全名：`.scratch/wisp/issues/253-panel-inbound-has-three-ruler-holes-...-zero-audit.md`（**无 `-done` 后缀** ⇒ 未结，本腿是预检不是验收）。
- 前件读讫：`.scratch/wisp/probes/253/r3/precheck.md`（131 行，本腿据其结论 ①②③④ 复认/推翻，不重复其四条读数）。
- 本腿引用基准（现读，非待验）：`internal/panel/git_test.go`、`internal/panel/composer_test.go`、`internal/panel/bridge.go` 三枚文件的行号取自 `3736f0dd` 工作树此刻盘上状态；⚠ 三枚写腿正在飞，行号可能在落地时漂移，落地腿起手须复跑定位。

## §1 两把尺现在数什么（逐枚行号＋want）

> 均为 `3736f0dd` 工作树现读，行号=盘上状态。票面 §现量-2 引的 `routeLiteralRe（composer_test.go:394）`、`panelMethodRe（cmd/wisp/git_test.go:394）`——后者**包级已纠**：`panelMethodRe` 实身在 **`internal/panel/git_test.go:394`**，`cmd/wisp/` 下零命中（与 253-r3 §五-2 同判，本腿复认）。

### 尺 A：`routeLiteralRe` — `internal/panel/composer_test.go:394`

- 定义：`var routeLiteralRe = regexp.MustCompile(` + "`" + `"panel\.[A-Za-z.]+"` + "`" + `)`（**硬要求 `panel.` 前缀**）。
- 使用点：`scanRendererHostDoors`（`composer_test.go:424-487`）第 `475-478` 行——对**每一行 renderer 源**取所有 `panel.*` 字面，凡不在 `composerRouteLiterals()`（`:409-419`，闭集 5 名：四枚 `Method*` + `"panel.approval.request"`）者，进 `rep.unknown`。
- 扫描根：**`frontend/src` 树**（`TestTheRendererHoldsExactlyOneDoorToTheHost` `:504` 扫真树；`TestPlantedRendererDoorShapesGoRed` `:533+` 扫 `t.TempDir` 副本）。
- want：`rep.unknown` 必须 0（`:521` `len(rep.unknown) != 0` 即红）。同函数另有四桶：`outside`=0（`:505`）、`callSites`==2（`:509`）、`assembly`=0（`:513`）、`computed`=0（`:517`）。
- 盲区：正则只认 `"panel.` 起头 ⇒ 页面里 `"config.get"`/`"config.set"` 字面**根本不进 `unknown` 分母**，也就永远不被"Go 侧答不答"对账。⇒ 票 §现量-2 的前端侧那半洞，confirmed。

### 尺 B：`panelMethodRe` — `internal/panel/git_test.go:394`

- 定义：`panelMethodRe = regexp.MustCompile(` + "`" + `"(panel\.[a-z0-9_.-]+)"` + "`" + `)`（同样**硬要求 `panel.` 前缀**）。
- 使用点：`whitelistMethodsFromSource(t, path)`（`:407-423`）——读**单一指定文件**，去重＋排序所有 `panel.*` 字面。
- 扫描根：**只有 `internal/panel/bridge.go` 这一枚产码文件**（`:385` 与 `:517` 都点名该路径；`_test.go` 与 `frontend/src` 都不在它射程）。
- 现数集合：bridge.go 里 `panel.*` 字面恰 **4 枚**（`panel.mode.request`/`panel.workspace.request`/`panel.attachment.add`/`panel.message.send`）。`MethodConfigGet="config.get"`、`MethodConfigSet="config.set"`（`bridge.go:66-67`）**无 `panel.` 前缀 ⇒ 抽不到**。
- want／锚（三处，逐一）：
  - `git_test.go:385` `if got := whitelistMethodsFromSource(bridge.go); !equalStrings(got, wantMethods)` → `wantMethods` = 四枚 `Method*` 常量排序（`:381-383`）。红文案 `:386-388`："the C17 panel.* whitelist in bridge.go = %v, want the four methods this ticket may not extend"。**← 票面/253-r3 说的"want-4 锚"正身在此**。
  - `git_test.go:517` `if real := whitelistMethodsFromSource(real bridge.go); len(real) != 4` → want 4，红 `:518`。**← 第二枚 want-4 锚**。
  - 正控 `git_test.go:514`：对 `t.TempDir` 里自造的 `bridge_five.go`（`:503-509` 五枚 panel.*）断 `==5`——**自包含，不碰真树**，证明尺 B 有牙。
- 设计意图（现读 bridge.go 注释自证）：`bridge.go:54-60` 逐字写明 config 两枚"**carry no `panel.` prefix on purpose**"，理由＝"ticket 181's anti-drift nail … counts the `panel.*` literals of this file（git_test.go 的 whitelistMethodsFromSource 抽取器）"。⇒ **want-4 与 config 落洞不是巧合，是当初为避免撞 want-4 而刻意把 config 写成无前缀**；配表证据在 `docs/evidence/s1/248-settings-write-path-r1.md`（票 248 落地面，非本腿射程）。

## §2 AC#2 两形各会不会响（ⓐ 只加尺／ⓑ 须新增真方法名）

**先厘清"洞"的坐标系**：Go 侧**其实已经答** config 两枚——`knownComposerMethod`（`bridge.go:146-152`）六名全含 config，guard switch（`composer_dispatch.go:176-202`）六 `case` 亦含。缺的只是**两把词面名册尺**（§1 的 A/B）因正则绑死 `panel.` 前缀而对 config 全盲。⇒ AC#2 要补的是"尺的射程"，不是"Go 的能力"。这个区分决定下面两形的批准结论。

### 形 ⓐ：只加一枚"全量入向名册"尺（不改 bridge.go、不往 C17 四枚里加真方法）

- **做法**：新尺从**产码能力面**取"Go 答哪些名"（读 `knownComposerMethod` 的 case、或 AST `composer_dispatch.go` 的 switch ⇒ 天然含 config），再要求"页面可拼出的每个入向方法名 ∈ 该集"，正控＝往 `t.TempDir` 副本种一枚**没进名册**的名字断尺红（自包含，照 `git_test.go:503-519`/`composer_test.go:533+` 的 plant 定式）。
- **会不会间接动 want-4？逐尺判：不会。**
  - `git_test.go:385`/`:517` 只经 `panelMethodRe` 读 `bridge.go` 单文件；ⓐ 不往 bridge.go 写任何 `panel.` 字面、不改 `panelMethodRe` ⇒ 仍数 4 ⇒ 绿。
  - `composer_test.go:521`（unknown）只经 `routeLiteralRe` 扫 `frontend/src`；ⓐ 新尺不往真前端树加东西 ⇒ 绿。
  - **没有任何既有尺把 `_test.go` 或本包测试串进分母**：`gitToolNamesUnder` 显式跳 `_test.go`（`git_test.go:464`）且只扫 `internal/tools`；`whitelistMethodsFromSource` 只点名单文件；`scanRendererHostDoors` 只走 `frontend/src`。⇒ 新尺文件里出现的 `"config.get"` 等字面**不进任何 want**。
  - 冻结件 `l2_grant_boundary_test.go` 的 `guardRosterOf`/`TestAnsweredPanelRoutesCarryNoApprovalDecision`（`:1291-1306`）对的是**产码** guard case ↔ `Method*` 常量双向；ⓐ 只加测试尺、不加常量/不加 case ⇒ 不触。
  - CI 的 `tools/d22scan` 禁的是代码形状（裸 `go func`、`risk.PathResolver` 外 filepath、UI emoji、墙钟超时…），**不审"尺的射程"**；一枚只读新尺不是契约变更。
- **结论 ⓐ：不需人工批准**，且不会在任何 want-4/want-5 行间接变红。"种一枚没进名册 ⇒ 必须红"这一发在**新尺自己的 TempDir plant 内**达成，对产码纯只读。
- **⚠ 两条落地选形警告（都仍不撞 git want-4）：**
  - 若 ⓐ 被做成"**就地加宽 `routeLiteralRe`**（`:394` 正则扩成含 `config.`）"而非"另立新尺"：则 `composer_test.go:475-478` 会开始把页面里的 `config.*` 字面拿去比 `composerRouteLiterals()`（`:409-419` 未列 config）⇒ 若真前端确实拼了 config.get/set，`:521` `unknown != 0` 会**红**；消红须把 config 名补进 `composerRouteLiterals`（仍是测试具器编辑、非 C1–C32）。⇒ 这条选形比"另立新尺"脏，且**其是否红取决于真前端有没有拼 config 字面——本腿读不到（frontend 禁区），〔量不到〕**。安全选形＝另立新尺，不依赖该读数；⛔ 但**不据此排除加宽选形**（铁律：只报"会不会响、响在哪行"）。
  - 若 ⓐ 的新尺去扫 `bridge.go` 且把 config 常量当"入向名"计数，注意**别复用 `whitelistMethodsFromSource`**（那会把 want-4 语义搅进来）；用独立抽取器即可零撞。

### 形 ⓑ：必须新增真方法名（往 C17 那四枚里加一枚 `panel.*`）

- 若把新名写成 bridge.go 里的 **`panel.` 字面**：`whitelistMethodsFromSource` 立刻数到 5 ⇒ `git_test.go:385`（`equalStrings` vs 四常量）**红在 :385-388**，且 `git_test.go:517`（`len != 4`）**红在 :518**。要让 ⓑ 变绿须改这两行 want ⇒ 即"放宽名册尺让它绿"，撞票禁区 + `:386` 自述"the four methods this ticket may not extend" ⇒ **需人工批准**（具名解冻，A487/票 248 对 l2 冻结锚的先例形；253-r3 §四-4 已点程序存在、批准归编排者）。
- 若 ⓑ 加的是**非前缀名**（再来一枚 `config.*`）：`panelMethodRe` 抽不到 ⇒ git want-4 不红；配 常量＋guard case 成对加 ⇒ l2 双向（`:1299-1303`）仍等 ⇒ 不红。**但这本身就是新增一条产品入向路由 = 产品契约变更**（新增 answered method），与"尺洞"无关，另须人工拍板。
- **归位判语**：AC#2 的靶子是"**不带前缀的名字掉在两把前缀尺之外**"——正解落 ⓐ（补射程），ⓑ-panel 既**不对靶**（config 无前缀才是要害）又**必撞 want-4**、需批准；ⓑ-nonpanel 不撞 want-4 但仍是产品改、需批准。⇒ **AC#2 能在不碰 want-4 的前提下落地，当且仅当走 ⓐ（一枚独立只读全量名册尺）**。

### 直接回答 §2 那一问

- "只加尺、不改 bridge.go ⇒ 是不是不需人工批准？"——**是**，ⓐ（另立新尺＋自包含 plant）不需人工批准。
- "名册尺扩射程这件事本身会不会被某把既有尺判成契约变更？"——**不会**：既有尺（git 两把、l2 guardRoster、d22scan）的分母都不含"另一枚测试尺"，也没有任何尺对"测试文件里的方法名字面"计数；唯一会因"bridge.go 多出 `panel.` 字面"而响的是 `git_test.go:385`/`:517`，而那是 ⓑ 的形状、不是 ⓐ。

## §3 AC#1 可复用零件名册＋"cmd/wisp 现在到底有没有页面自报用例"的复认或推翻

> 全部在 `package main`；两枚文件（`panel_resident_windows_test.go`／`panel_host_windows_test.go`）皆 `_windows_test.go` ⇒ 隐式 `GOOS=windows` 构约（本机 win32 可复用）。`panel_resident_windows_test.go:3` 另有显式 `//go:build windows`。

### 3a. `ac14AwaitJS` 定式所依赖的零件（逐名，各带 file:line）

一枚新用例直接复用下列全部即可跑"页面发起→Go 收→回执到页面"的闭环：

| 零件 | file:line | 作用 |
|---|---|---|
| `startPanelForTest` | `cmd/wisp/panel_resident_windows_test.go:59-67` | 起**生产形态** manager：真 embed（`:63` 走 `builtinAssetsOrTestNil`）＋真 user-data 目录＋真 fail-closed 派发链；`:62` 把 `disp.Mode` 装成 `&recordingModeHandler{}` |
| `builtinAssetsOrTestNil` | `:72-80` | 用产品同一 API 读 embed；无 bundle 的树返回 nil |
| `showAndWait` | `:104-118` | 驱动产品自己的开窗路（`RequestShow`→面板线程 Show），返回于线程**跑完 Show** 之后（非 HWND 首现） |
| `waitPanelTrue` | `:85-95` | 单调截止轮询；超时＝`t.Fatalf` 判"测量失败"，不当通过 |
| `evalOnPanelThread` | `:166-189` | 在**活文档**、且在**面板线程上**跑任意 JS（`:175` `w.Eval(js)`）——注入正控/脚本的现成缝，不需 SetHtml 换页 |
| `recordingModeHandler` | 定义 `cmd/wisp/panel_host_windows_test.go:227-243`（`HandleModeRequest` `:232`、`count()` `:239`）；`all()` 在 resident `:206-210` | **唯一的"假 handler"**，收下并记录到达的 `ComposerRequest` |
| `reportJSEnv` | `:250-253` | 拼一枚 `method:"panel.mode.request"` 的 composer 信封（payload 进 `to`），走产品门回话 |
| `ac14AwaitJS` | `:786-810` | 现成的 await-REPLIED 脚本：`:801` beacon（`ac14r-beacon`）证脚本跑过；`:804-806` `Promise.race([wispDispatch(ENV).then=>'REPLIED', 2s=>'TIMEOUT-2S'])`；`:808` 结果**经同一扇门**报回 |
| `awaitReport` | `:215-229` | 等某个 requestId 的报告，超时＝`t.Fatalf` 失败测量（附 `describeRequests` `:233` 现场） |
| `countRequests` | `:837-845` | 按 `strings.HasPrefix(req.RequestID, prefix)` 分桶计数——把"真请求"与"报告"分开 |
| `panelThreadWait` | `:53`（= `15 * time.Second`） | 上述超时预算 |

**复用结论：能直接复用。** 新用例只需 `startPanelForTest`+`showAndWait` 起手、`evalOnPanelThread` 注入自造的 JS、`awaitReport`/`countRequests` 读回——三件套齐。三条**必守的复用红线**（撞过既有钉）：
1. **requestId 前缀必须自开一桶**（`countRequests` 按前缀分桶，`:837-845`）；沿用 `ac14-`/`ac13-` 会污染 `TestAC14AwaitedBindingReplyReachesThePage:830`（`calls != 3`）与 `TestPanelHostRealWindowHopAndLifecycle` 的 `h.count()` 精确计数。
2. **别再 SetHtml 自造探测页**——会撞 `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe:314` 的"冷启结束于 embed 页"次序钉；注入一律走 `evalOnPanelThread`（现成）。
3. `recordingModeHandler` **只答 mode 路**（`HandleModeRequest`）；能力尺若要端到端验**非 mode 路**，得给该路挂 handler，否则落 `ErrNoHandlerAttached`（见 `internal/panel/config_route_248_test.go:173`）。最省事的复用（同 AC#14）＝在 mode 路上验页面→Go→页面 hop。

### 3b. cmd/wisp 现在到底有没有"页面自报（页面发起调用、Go 侧回执）"用例？——**有**，具名

1. `TestAC14AwaitedBindingReplyReachesThePage`（`panel_resident_windows_test.go:812`）——**正是要找的那形**：页面对**生产 embed 活文档**、经**生产门 `window.wispDispatch`** 发三枚真信封（`ac14-0/1/2`），断言**两半互证**：Go handler 被达 3 次（`countRequests` `:824`、`:830` `calls != 3`）**且**页面收到 `REPLIED`（`awaitReport` `:821`、`:826-828`）。
2. `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`（`:314`）——页面经 `wispDispatch(reportJSEnv("ac13-probe", b))`（`probeJS :301`）自报位图，`awaitReport` 读回（`:325`）。同为"页面发起、Go 回执"。
3. `TestAC14GoSideEvalPushReachesThePage`（`:855`）——Go→页面 push 维（回程报告仍骑 `wispDispatch :864`），是 AC#14 的第二枚独立钉。

（另：`panel_host_windows_test.go:507` `TestPanelHostRealWindowHopAndLifecycle` 是 **Go 侧直驱 router**（`:554 mgr.dispatchRaw`）＋ recordingModeHandler 计数——属"Go→派发"，不属"页面发起"；`panel_config_248_test.go` 的 stdin 端到端腿是 CLI 入向，非页面。二者不算"页面自报"。）

### 3c. 对 253-r3 结论 ②"零枚断言"的复认／推翻判语

- **复认（窄读）**：r3 §1b 的原句是"**生产绑定 `w.Bind(panelDispatchBinding, ...)` 零枚测试断言它存在**"——即**没有任何词面/AST 尺去读 `cmd/wisp/panel_host_windows.go` 断言 `const panelDispatchBinding="wispDispatch"` 与那枚 `Bind` 接线存在**。这句**成立**：`wispDispatch` 在两包 `*_test.go` 里只被测试**自己当门用**（resident `:301/:801/:808/:864`），无一把尺把"产码里有这枚 Bind"当断言对象。
- **推翻／收窄（宽读会被误用）**：若把"零枚断言 wispDispatch"读成"**没有任何用例验过'页面经 wispDispatch 发起、Go 收到'这一 hop**"，则**为伪**——§3b 的 AC#13（`:314`）与 AC#14 nail1（`:812`）恰恰就是这一 hop 的行为断言，跑在生产 embed 活文档、走生产 `window.wispDispatch` 门。若那枚 Bind 没接上，`window.wispDispatch` 未定义 ⇒ 脚本抛／无回话 ⇒ `awaitReport` 直接 `t.Fatalf` 失败测量。⇒ **生产门的"可达性"今天已被票 33 的 AC#13/AC#14 行为覆盖**，只是没被那两把**词面名册尺**（§1 A/B，`internal/panel` 里认 postMessage/panel. 字面）覆盖，也没被一句"断言 Bind 存在"覆盖。
- **对 AC#1 的直接后果**：AC#1 要的"判据改问能力（页面发起→Go 收到并回执）"**在 cmd/wisp 里已有现成载体**（AC#13/AC#14 那形），不是从零造。AC#1 的真正增量只有两条：Ⓐ 把"**可达性那一格的裁决**"从 `internal/panel` 的词面尺**迁/挂**到这枚能力形上；Ⓑ 补票面点名的正控"**只 `postMessage`、不 `wispDispatch` ⇒ 指名用例必须红**"。⛔ 本腿不裁 Ⓐ 具体落哪枚尺、⛔ 不建议改法，只登记"载体已存在、'零断言'须限定为'零对产码 Bind 的词面断言'"。

## §4 我可能写错的条目（自我对抗）

1. **want-4 双锚里 `:517` 的可达性**：`:517` 在 `TestPlantedGitToolShapesGoRed` 内、读的是**真 bridge.go**（"clean half must stay quiet"）。我判形ⓑ（往真 bridge.go 加 `panel.` 名）会**同时**红 `:385`（主尺）与 `:517`（正控尺的干净半）。若落地腿把第五名加进 `wantMethods`/plant 而只修一处，另一处仍红——这条我按"两枚都读真树"判，未实跑，落地须复跑确认二者同红。
2. **"没有既有尺数测试文件字面"是靠排除法，不是穷举**：已核 `gitToolNamesUnder` 跳 `_test.go`（`git_test.go:464`）、`whitelistMethodsFromSource` 只点名单文件、`scanRendererHostDoors` 只走 `frontend/src`、**l2 的产码枚举亦跳 `_test.go`（`l2_grant_boundary_test.go:369` `strings.HasSuffix(p,"_test.go")` → 只 AST 产码）**。⇒ 形ⓐ新测试文件确不进任何 want。但 `tools/d22scan`（CI，非 `go test`）本腿**未读其源**，我"ⓐ 不是契约变更"那句里对 d22scan 的判断**仅据 AGENTS.md §1.2 的禁止清单**（其中无"尺射程"一项），非现读——列为〔半量〕。
3. **emoji/仪器射程**：新尺若把 `config.get` 之类当**串面**种进测试，`d22scan` 的 emoji 尺"注释豁免、字符串不豁免"——但那是 emoji 维，与方法名无关；我未在新尺设计里引入 emoji，若落地腿引入则另论。低风险，登记不裁。
4. **AC#1 复用的"活文档确有 `window.wispDispatch`"前提**：`startPanelForTest`→真 `PanelManager`→产码 `Bind("wispDispatch")`，逻辑上门必在；但 embed 未构建的树里 AC#13 自己 `t.Skipf`（`:317`）。⇒ 我"生产门可达性已被 AC#13/AC#14 行为覆盖"这一推翻，**成立与否取决于测试机上 bundle 是否已构建**（运行时事实，本腿禁 Go 命令量不到，见 §5-3）。若该树常 skip，则严格说 AC#13 对生产门"零运行断言"，r3 ②"零断言"就**不算被推翻**、只是被我收窄——两种读法都要落地腿用一次真跑裁决，⛔ 本腿不锁。
5. **§1 行号 vs 三枚在飞写腿**：`internal/panel/git_test.go`、`composer_test.go`、`bridge.go` 均在我引用后可能被 `198-v1`/写腿改到；我全程标"行号取自 `3736f0dd` 盘上现读、落地须复 grep 定位"，未把行号当不变量。

## §5 量不到的地方（只读、禁 Go 命令）

1. 〔量不到〕上述任一 `go test` 断言的**实际 PASS/FAIL**（含 AC#13/AC#14 此刻跑不跑、`want-4` 两锚此刻绿不绿）——须 `go test ./internal/panel/ ./cmd/wisp/`，本腿禁跑。落地腿起手先复跑在册绿基线（panel 侧 4 枚已知红见 253-r3 §ⓔ，不可当对照）。
2. 〔量不到〕真 `frontend/src` 里到底有没有 `config.get`/`config.set` 的**字面调用点**（读 frontend 是两层禁令）——这一读数决定形ⓐ若被做成"**加宽 routeLiteralRe**"会不会连带触发 `composer_test.go:521` unknown≠0。**"另立新尺"这一安全选形不依赖它**，故 §2 仍能给"ⓐ 可只加尺落地"的确定判语。
3. 〔量不到〕生产 `wispDispatch` 门在本机的**运行时**可达性（WebView2 冷启、embed 是否已构建、AC#13 的 skip 条件是否命中）——静态已给（§3），动态验证要真窗（`TestPanelHostRealWindowHopAndLifecycle` 那档），归票 253 落地/验收腿。
4. 〔量不到〕`git_test.go:385`/`:517` 两枚 want-4 锚**能否依 A487 同款程序具名解冻**——程序先例在册（253-r3 §四-4），批准归编排者，本腿不裁、不建议选形。
5. 〔半量〕`tools/d22scan` 对"形ⓐ新尺"是否有话说——本腿未读其源，仅据 AGENTS.md §1.2 禁止清单判"尺射程不在其列"；落地腿若要动用可 `git log`/读 `tools/d22scan/main.go` 复核（非 Go 命令）。

## §6 交件判语

- **只读**：全程只 Read/Grep/Glob＋`git log/status`＋读盘；⛔ 未跑任何 Go 命令（`go build/test/vet/run` 零发）。
- **未 push**：只 commit、未 push；每次 commit 显式 pathspec、`add` 与 `commit` 串同条命令。
- **零产码／零冻结件触碰**：只写 `.scratch/wisp/probes/253/p4/**`（precheck.md + msg.txt，临时件只建不删）；`internal/panel/**`、`cmd/wisp/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、allowlist.txt、tokens_fourway/l2_grant/ticket90_persist、ci.yml 一字未改。
- **票面 AC 框未碰**：AC#1..AC#4 勾选框一律不填不改。
- **frontend/design 未读未引**：证据件不含任何来自那两目录的内容；§2 里凡涉及真前端处皆标〔量不到〕。
- **对 253-r3 的净判语**：结论①（词面尺认串不认为）、③（能力尺只能落 cmd/wisp）、④（AC#2 落地撞 want-4）全部**复认**；结论②"wispDispatch 零枚断言"**复认其窄义**（零对产码 Bind 的词面断言）、**推翻其宽义**（cmd/wisp 实有 AC#13 `:314`/AC#14 nail1 `:812` 两枚页面自报→Go 回执的行为钉）；并纠 AC#1"从零造"的印象——能力形载体已在，AC#1 的增量是"迁裁决＋补只-postMessage 正控"。
- 交件时 HEAD 应为末次 commit 的父；见文末收笔锚。

## 收笔锚

- 填 §4–§6 时同发三取：`date -Iseconds` = `2026-10-03T09:22:04+08:00`；`git log -1 --format=%h` = `fae59004`（本节所属的**最后一次内容 commit 的父**，本文件随后再提一次收笔 commit）；`git status --porcelain | wc -l` = `408`（他腿脏面，本腿仅 `probes/253/p4/**` 两枚文件）。
- 本腿 commit 链：`9051d51c`（§0 锚＋骨架）→ `fae59004`（§1/§2/§3）→ 收笔 commit（§4/§5/§6＋本锚）。
- 占位符清零：全文无「（取数中）」残留。
