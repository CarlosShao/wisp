# 票 253 · `253-v2` 验收腿裁决（AC#2；非实现者）

- 时刻 2026-10-14 13:56–14:18 +08；起手 HEAD `581dd252`，收笔 HEAD `c2b00de4`（他腿 commit）。
- 被裁＝`253-r5` 的 `internal/panel/inbound_roster_253_test.go`（834 行，sha256 `646125ce…105e9c`）。起手/终态该件与 `bridge.go`（`fc23919c…736c37`）、`composer_dispatch.go`（`d001fb22…b248ad`）均逐字节复原，`internal/panel` porcelain 空。
- 承接口：票面 AC#2 逐字 ＋ 票内节 8 ★ 的更严结案口径（名册＝guard＋router 同批，只补名册不算登记）。
- 原始读数：`.scratch/wisp/probes/253/v2/`（锚 `00-anchor.md`＋13 枚 log＋2 枚 diff；本文每句都可回指）。

## 0. 判语

**AC#2 成立但附条件**。四条条件：

1. **只许整腿取读数**（N2 实测见 §4）：`TestFullInboundMethodRosterIsClosed` 单枚可被"名册加名＋产码补绑定"洗绿；整腿 rc=1 咬住它的不是主测而是两枚夹具（guard 整行写死）。
2. **射程边界一条**（§5）：① 分母只读 `bridge.go`；一枚声明在**别的产码文件、且答侧无人接**的点分常量读绿（它还不是"门"）。一旦接上答侧（ident 或字面）即红（③/⑤）。
3. **未种项不作凭据**：真树上"完整登记后仍绿"未量（要动 guard＝契约面）；页面侧调用点未读（`frontend/**` 两层禁令）。
4. **更严结案口径**照旧适用于将来：第五名进名册必须五分母同批动（见 §6 的 `config.reload` 建议）。

**翻框建议**：AC#2 那枚 `- [ ]` 可翻 `- [x]`；留档句＝「AC#2 成立（非实现者腿 `253-v2`，2026-10-14）：五分母双向有牙（M1/M2/M3/M4 各红句逐字在案），能力形指名 `bridge.go:68 const InboundPeek253v2`；结案口径＝名册/guard/router 三处同批，只许整腿取读数。」撤销口令＝「撤 253 AC#2」。⛔ 本腿未翻框、未动 `Status:`/票名。

**对"更严结案口径"的态度**：**认同**（三支答侧我的正控各自单独打红：② `:465`、③ `:476`、④ `:496/:500`）。**一处保留**：④ 要求 `default:` 必须调 `rosterMismatch` 是把一个实现细节钉成契约——合法重构（改名/内联回填支）会误红，见 `253-r5` §5-6 自述；验收读红句时要能分辨这一形。

## 1. 五分母现读（我的独立提取，不复用尺的代码；`denominators-read.log`）

| 分母 | 我的读数 | 关键证据 |
|---|---|---|
| ① DECLARED（bridge.go 包级字符串常量，Method 前缀 ∪ 点分值） | 恰 6 名；`ComposerRequestSource`（无点、非 Method）被正确排除，`requestIDBytes` 非字符串 | `MethodConfigGet="config.get"`(:66) 等 6 枚 |
| ② RUNNING GUARD | 6 名全答；`config.peek/settings.read253v2/config.toml/index.html/panel.mode.request.extra/panel-composer` 全拒（`probe_253v2_test.go` 实跑 PASS） | `probe_253v2_test.go:27` |
| ③ GUARD AST（`bridge.go:146` switch） | 标签恰 6 值、可全解析、**无 default** | 我的 AST 转储 |
| ④ ROUTER AST（`composer_dispatch.go:181` dispatch） | 标签恰 6 值（经 bridge＋全包常量解析）、`default:` 存在且调 `rosterMismatch` | 我的 AST 转储 |
| ⑤ ANSWERED 字面（全包产码点分字面 ∩ 运行中 guard） | 全包点分字面 7 值：6 名自身（bridge.go:42-45/66-67）＋`index.html`(assets.go:29)；后一枚 guard 不答 ⇒ 不进分母 | 我的扫描 |

- **p4 旧量复核**：router case 列表此前确实不在任何尺的分母里——`composer_dispatch_test.go:380` 只拒 `"panel\.[^"]*"` 前缀字面、`composer_test.go:317` 只数单个 `case MethodWorkspaceRequest:`、冻结件 `l2_grant_boundary_test.go:777` 的 answered 只从 `knownComposerMethod` 导出。分母 ④ 补的正是这格（我的 ④ 正控见 §2）。
- **overlay 边界实测**：overlay 改 `bridge.go` 对本尺**运行时**读不可见（主测仍自述 `6 names declared`）⇒ 真树种形一律用 blob 备份/还原，overlay 只用于编译面（测试副本、新增文件、守卫行为）。

## 2. 各分母正控（我的种法；红句逐字，出处见对应 log）

| 分 | 我种什么 | 判红逐字（节选） |
|---|---|---|
| ① | 真树插入 `InboundPeek253v2 = "settings.read253v2"`（`bridge.go:68`，与 r5 的尾插/命名不同） | `:437 … declares the inbound method name settings.read253v2 at bridge.go:68 const InboundPeek253v2 and the closed roster does not carry it - a new inbound route entered the tree without being registered anywhere…` |
| ② | overlay 让**编译进去的** guard 一律 `return false`（盘上 case 列表不动） | `:465 … the roster carries config.get but the running guard knownComposerMethod refuses it - a rostered inbound name Go does not answer is a page-visible dead door` ×6；同发 ③ 仍自述 `6 guard case labels` ⇒ **② 问的是函数本体、不是复刻** |
| ③ | 真树 guard case 列表加字面 `"settings.read253v2"` | `:476 … knownComposerMethod's case list answers settings.read253v2 and the closed roster does not carry it - an answered route that was never registered` |
| ④ | 真树 router 加字面 `case "settings.read253v2":`＋`default:` 改调 `unattached` | `:500 … dispatch routes settings.read253v2 … the router is where an inbound name becomes an action…`；`:496 … the backstop that refuses an unrostered route (it must call rosterMismatch) is not what runs here` |
| ⑤ | overlay 让 guard 额外答 `index.html`（盘上 assets.go:29 已有的点分字面） | `:515 … the running guard answers the dotted literal "index.html" found at assets.go:29, and it is on no roster…`（②③ 同发保持绿） |
| 全 | M4 单发值漂移：`MethodConfigSet = "config.alt253v2"` | 五分母双向全读：① `:437`+`:441`；② `:465`；③ `:476`+`:479`；④ `:500`+`:503`；⑤ `:515`（共 8 句） |

## 3. 能力形（AC#2 判据本体）＋ 反形对照（我的）

- **能力形（真树，非副本）**：M1 种假腿 ⇒ 指名步 `TestFullInboundMethodRosterIsClosed` **红 `:437`**（红句逐字见 §2 ①，**指名了那枚常量与它落在哪一行**）。诚实登记：同一发里另 3 枚同族用例也红（件内"真树必须安静"半句与夹具），属**过杀伤形状**、非缺陷（`caps-run.log`；r5 §4-1 同形自报）。**同发我另测的白名单/冻结锚 7 枚全 `--- PASS` rc=0**（`TestGitDimensionHasNoModelCallableTool`／`TestPlantedGitToolShapesGoRed`／`TestAnsweredPanelRoutesCarryNoApprovalDecision`／`TestRealGuardRefusesEveryAssemblableApprovalRouteName`／`TestNoInboundEnvelopeCanBindAnApprovalVerdict`／`TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes`／`TestTheRendererHoldsExactlyOneDoorToTheHost`，`collateral-anchors.log`）⇒ M1 形（无 `panel.` 前缀、标识符非 `Method`）不牵连别人的锚，p4 的硬约束被实测遵守。
- **反形对照（我自己的一发）**：把分母 ① 换成**反形**＝空集也算过＋只报 missing、`_ = extra`（subset-only）；overlay 同一棵**种了假腿的真树**再跑主测 ⇒ **PASS**，它自己的 log 写着 `7 names declared in bridge.go` ⇒ **反形对这件事不敏感、不可作凭据**（diff＝`reverse-form.diff`，log＝`caps-run.log` 第三段）。

## 4. N2 复核：单枚可洗绿 / 整腿才可靠（`n2-run.log`、`n2-whole-leg.log`）

- 我的种法（r5 形、我重做）：名册＋两计数（overlay 测试副本，diff＝`wash-roster.diff`）＋ `MethodConfigWash253v2 = "config.wash253"` 声明与 guard 同批（真树 `bridge.go`）＋ router 加 case（真树 `composer_dispatch.go`）＝「加名＋补绑定」全一致。
- **单枚** `-run TestFullInboundMethodRosterIsClosed` ＝ **`--- PASS` rc=0（洗绿）**。
- **整腿** `-count=1 -v` ＝ **rc=1**，多出的两枚红逐字（washed 副本 +1 行偏移，真身 `:661`/`:705`）：
  - `:662: bridge.go's guard case list is no longer written as one comma-separated clause - re-read knownComposerMethod and rebuild this plant from what is on disk rather than replacing nothing`
  - `:706: bridge.go's guard case list is not the one comma-separated clause this plant was written against - re-read knownComposerMethod rather than trusting this line`
- 机制＝**两枚夹具把 guard 那一行整行写死**（`strings.Contains(src, guardLine)`），名册/绑定一变即 Fatal；主测五分母全一致时对它无话。⇒ 读数规成立：**单枚绿≠整腿绿**；"整腿红"这两枚属"夹具钉住当前 guard 行"，修它＝改夹具＝显眼评审动作。附注：若第五名走完**人批＋五处同批**的合法登记，整腿读全绿——那是登记的定义，不算洗绿。

## 5. 射程边界（我另量两发，`scope-boundary.log`）

- (a) 一枚点分常量声明在**别的产码文件**（overlay 新文件）、答侧无人接 ⇒ 主测 **PASS**（① 只读 `bridge.go`；它还不可被页面投进 `method` 字段）。
- (b) 同一常量（`composer_dispatch.go:67`）以 ident 接进 guard ⇒ 红：`:472 … case label CrossPeek253v2 names no package-level string constant of internal/panel: a route spelled out only here`；`:515`（该字面被 guard 答、off roster）。
- ⇒ 尺外只剩"悬空且无人答"的常量；**凡是成为门的形都红**。

## 6. `config.reload` 半量欠账：我的建议＝**不算入向名**（不落勾、供你落）

- 现量行号（我现读，派单写的 `:86` 是注释块；值已漂）：`cmd/wisp/config_reload.go:90` `configReloadTool = "config.reload"`；**唯一消费点 `:285`**＝`rt.gate.PendingApproval(…, tools.Decision{Tool: configReloadTool, …})`——宿主自己发起的 L2 卡**身份证**（出向）。
- 语义对照：入向链＝`bridge.go:126 ParseComposerRequest` → `:146 knownComposerMethod` → `composer_dispatch.go:181 dispatch`；两枚宿主（`panel_host_windows.go` dispatchRaw、`panel_inbound.go:164 disp.Handle`）都走同一名册。页面无处把 `config.reload` 放进 `method`。文件自述 `:82-87`：「a 'namespace.action' name that names the door being widened rather than a tool the model can call. Neither is a C4 registration and neither is a C17 method.」
- ⇒ **不算**，且**不应**进 `wantFullInboundRoster253`（进了 ② 立刻红＝dead door）。若将来真开"面板可触发热加载"，那是新的 C17 面（`panel.` 名、人工批准），五分母同批动，⛔ 不许复用这枚卡名。

## 7. 门禁（各现跑）

- 253 族（5 枚）起手与终态均 **`--- PASS`×5，rc=0**（`baseline-family.log`/`final-family.log`）。
- 整腿 `-count=1 -v`：起手＝终态＝**141 PASS／5 FAIL／0 SKIP**，5 红逐名相同：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`／`TestComposerContractTypesMatchFrontend`／`TestC21DesignTokensFourWayAgree`。
  - ⚠ 与派单"4 枚"、`253-r5` 当年"4 枚"不同：`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` 已不在红名单（界面侧改动），`ApprovalCardView…` 又回到红名单，`35r8` 两枚是 `b5439299` **有意交件即红**的在册件。**本腿未新增红、未修、未放宽、未 Skip。**
- `./tools/d22scan/d22scan.exe`＝`d22scan: clean - no D22 ban violations` rc=0。
- 长度门＝tracked path-length budget（`.github/workflows/ci.yml:139/:170`）：`sh scripts/check-path-length-budget.sh --with-self-test`＝**VERDICT GREEN** rc=0。
- 未定义即停：**无触发**（五分母对齐；假腿各形全红；`config.reload` 语义读通）。

## 8. 我没做的（具名）

- 未读 `frontend/**`/`design/**`（两层禁令）⇒ 页面侧有无 `config.get/set` 调用点未量（AC#2 不依赖它）。
- 未跑 `cmd/wisp`（`253-r6` 在飞，窗口内目击其新件落盘）；未跑全仓 `go test ./...`。
- 未在真树上做"完整登记后"的读数（要动 guard＝契约面）；该读数只在 overlay/临时种形里取（已还原）。
- 未测 ④ `default` 保留改名/内联的"误红形"（§0 保留项是推理＋r5 自述，未种）。

## 9. 足迹

- commit：只 1 笔（本文＋`probes/253/v2/**`＋票尾 1 bullet），显式 pathspec，⛔ 未 push。
- 真树突变 6 窗（M1/M2/M3/M4/W2/scope-b），每窗 `git cat-file blob` 还原后 sha256 逐字节相等、`git status --porcelain -- internal/panel` 空（`caps-run.log`/`guard-router-runs.log`/`value-drift-run.log`/`n2-run.log`/`scope-boundary.log` 内各有一段还原证）。
- 产码面零改动；`internal/panel` 无残留；`00-anchor.md` 记起手/终态锚。
