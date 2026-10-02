# 普查 114-a2 — 票 114 的三态现状：哪几格已被 248-r1 顺手做掉、哪几格仍欠、哪几格永不归本编队

> 只读普查腿。零 `go test`／零 `go build`／零 `go vet`／零 exe。零 `frontend/**`／零 `design/**` 读取。
> 本文件不勾任何 AC 框，不写台账，不改票名。判"已满足"一律问**生产调用者**与**能力**，不问词面。
> 本文件里出现的一切行号与枚数都是**本腿自己实测的**，凡引编排者题面的行号／枚数一律落在 §7 里对账。

## §0 起手锚

| 项 | 实测读数 |
|---|---|
| 分支 | `dev` |
| HEAD（起手时刻） | `3b23613b` — `probes(ci223-1): 追加 0.-2 边界自报…` |
| 248-r1 产码锚 | `0d87a681`（2026-10-02 09:03:15 +0800），已在 HEAD 之前 |
| 起手时刻 | 2026-10-02 09:14:18 +0800 |
| `git status --porcelain internal/panel cmd/wisp` | **空串**（零行）。⇒ 题面提示"此刻可能非空（`198-r1` 在飞的写面）"在本锚点**未复现**；本腿全程没动这两枚目录，若后续变脏一律视为别人的写面，不读不评不提交 |
| 本腿尺的作用域 | `cmd/wisp/**` · `internal/panel/**` · `internal/agent/approval/**` · `internal/config/**` · `internal/perm/**` · `tools/d22scan/**` · `docs/evidence/s1/**` · `.scratch/wisp/issues/**`；全仓递归一律带 `--exclude-dir=frontend --exclude-dir=design` |

⚠ **并发写面自报（本腿全程未动、未提交、不评论其内容）**：交件前复量 `git status --porcelain internal/panel cmd/wisp` ⇒ **4 行**（`M cmd/wisp/run.go` ＋ 三枚 `?? cmd/wisp/firstrun*.go`／`firstrun_198_test.go`／`firstrun_acl_198_windows_test.go`）＝题面预告的 `198-r1` 在飞的写面。**直接后果：本文件里所有 `cmd/wisp/run.go` 的行号（`:316`／`:322`／`:637`／`:659`／`:808`／`:822`）在 HEAD 上可能已漂**，引用者请按本仓规矩现读为准（同一条教训见 `docs/evidence/s1/248-settings-write-path-r1.md` §5-11）。本腿的所有 commit 只用显式 pathspec 提了自己那一枚文件。

---

## §1 "`ParseComposerRequest` 在生产里零调用者"这一句今天还成立吗

**结论先给：不成立，而且与 248-r1 无关。** 这句在立票那天（2026-09-21）成立；把它变成假的那一发是 **`6609e7e1`（2026-09-28 17:21，33-r3）**——`cmd/wisp/panel_inbound.go` 的 `disp.Handle` 从那天起就在产码里；把它接上**真窗口**的是 **`13acad46`（2026-10-01 13:52，33-r5）**。248-r1（`0d87a681`，10-02 09:03）只做了两件与本句有关的事：把名册从四枚扩到六枚、把 `Config` 槽接进那枚已经存在的链。**它没有新建 `panel_inbound.go`**（见 §7-R2）。

⚠ **复跑本段任何符号计数前先排除本文件**：`.scratch/wisp/probes/114/a2/census.md` 自身含下面每一个符号名，带进全仓递归就把尺的文本数成调用点。

### 1.1 调用点名册（逐枚；`--include=*.go`，作用域 `internal/ cmd/ tools/`）

| # | `文件:行` | 是什么 | 生产码还是测试码 |
|---|---|---|---|
| D | `internal/panel/bridge.go:126` | **定义** `func ParseComposerRequest(raw string) (ComposerRequest, error)` | 生产 |
| C1 | `internal/panel/composer_dispatch.go:155` | **全仓唯一的一处调用**：`req, err := ParseComposerRequest(raw)`，在 `(*ComposerDispatch).Handle`（定义 `:153`）里 | 生产 |
| T1–T6 | `internal/panel/bridge_test.go:66, 83, 86, 92, 102, 112` | 封套自身的用例（其中 `:86` 是断言消息里的名字，不是调用） | `_test.go` |
| T7 | `internal/panel/config_route_248_test.go:78` | 248 的名册用例 | `_test.go` |
| T8–T9 | `internal/panel/l2_grant_boundary_test.go:1970, 1989` | 冻结件的"历史线上形状必须被 parse 拒"行为支 | `_test.go`（**冻结件，只读内部判射程**） |

尺读数（本腿 `09:2x` 现跑，`-E "ParseComposerRequest\("` 去掉 `_test.go` 后）＝**恰好两行**：`bridge.go:126`（定义）与 `composer_dispatch.go:155`（调用）。⇒ 生产调用者枚数＝**1**，不是 0。

其余命中全是**注释**（`bridge.go:16`、`composer_dispatch.go:6/:17`、`composer_handlers.go:37`、`cmd/wisp/run.go:318`、`cmd/wisp/panel_host_windows.go:317`、`composer_dispatch_test.go:13`、`l2_grant_boundary_test.go:57/66/67/70/1274/2289`）——其中两枚注释今天已经是**过期陈述**，见 §7-R5/R6。

**"Handle 有没有生产调用者"才是那句题头真正问的东西**，所以第二层名册：

| # | `文件:行` | 调用形 | 谁把它启起来 |
|---|---|---|---|
| H1 | `cmd/wisp/panel_inbound.go:163` | `reply, err := disp.Handle(ctx, raw)`，`raw` ＝stdin 一行 | `cmd/wisp/main.go:115-120` 的 `case "panel-inbound": os.Exit(cmdPanelInbound(...))` ⇒ 操作员敲 `wisp panel-inbound -data <目录>` |
| H2 | `cmd/wisp/panel_host_windows.go:551` | `return m.disp.Handle(ctx, raw)`（`dispatchRaw`），`raw` ＝页面递给绑定函数的实参 | `:319` `w.Bind(panelDispatchBinding, func(raw string) string {...})`；`panelDispatchBinding = "wispDispatch"`（`:82`）；这枚 `PanelManager` 由 `cmd/wisp/panel_resident_windows.go:204 NewPanelManager(disp, assets, dataPath)` 装配，其构造者 `newResidentPanelManager` 被 `cmd/wisp/resident_windows.go:142` 调用，而 `runResident()` 的唯一入口是 `cmd/wisp/main.go:66`＝**双击/无参数启动那个 GUI 进程** |

⇒ **`wisp panel-inbound`（操作员驱动）与常驻面板窗口（页面驱动）两条都是产码。** `cmd/wisp/panel_inbound_33_test.go:182 TestAC9ComposerDispatchHasAProductionCaller` 就是钉住这件事的那枚仪器，它今天有牙齿（不是恒真）：`inboundHandleCallSites33`（`:278`）从非测试文件里派生调用点，删掉 `panel_inbound.go:163` 那一行它才红。

### 1.2 谁产出 / 谁投递 / 谁落盘（三处分开设，合起来才答"能力在不在"）

- **谁产出那枚封套（raw 从哪儿来）**：两支，且**不是同一种线上形状**。
  - `wisp panel-inbound`：stdin 一行一封装套（`panel_inbound.go:153-163`），形状由测试夹具给（`panel_inbound_33_test.go:105`）。
  - 页面侧：go-webview2 的 `Bind` 注入一枚 JS 存根 `window.wispDispatch(...)`，它发的是 **`window.external.invoke(JSON.stringify({id, method:"wispDispatch", params:[raw]}))`**（模块 `webview.go:450-482`），宿主再在 `WebMessageReceived` 里按 `rpcMessage` 解包查 `w.bindings[d.Method]`（`webview.go:139-168`）。⇒ **页面今天调得出 `Handle` 的那枚动词是 `window.wispDispatch(raw)`，不是 `window.chrome.webview.postMessage(raw)`。** 票 114 与票 92 那族仪器全部只认后一种写法（见 §2 AC#9/AC#10 与 §5-丙）。
- **谁投递**：H1 / H2 上面两行。
- **谁落盘**：两枚处理器各有自己的落盘腿，都汇到同一个 `config.Manager`。
  - 档位：`internal/panel/composer_handlers.go:141` → `ModeWriter.Set` → `internal/perm`（`store.go:172-214` 逐档写审计）→ `config.SetPermissionMode` → `config.toml` 的 `[risk] permission_mode`。
  - 设置：`internal/panel/config_handlers.go:307`（写闸门）→ `cmd/wisp/panel_config_store.go` → `internal/config/settings.go` 的导出 setter → `mergeWrite`（`internal/config/writeguard.go`）按键受保护写 → `config.toml`；凭据值走 `internal/secret`（DPAPI），`config.toml` 只落 `dpapi:` 引用。
  - 工作区 / 附件 / 消息：**没有落盘腿**（见 1.3）。

### 1.3 `panel_inbound.go` 那条新通路接的是哪一类入向

`newComposerDispatchChain`（`cmd/wisp/panel_inbound.go:228`，CLI 与常驻两支共用它：`newPanelInboundDispatch` `:219`、`newResidentComposerDispatch` `cmd/wisp/panel_resident_windows.go:166-170`）六枚槽的实测装配：

| 槽 | 装了什么 | `文件:行` | 能力后果 |
|---|---|---|---|
| `Mode`（**档位**） | `*panel.ModeWriteHandler`，`Modes`＝真 `perm.Store`，**`Confirm: nil`** | 装配 `panel_inbound.go:248-255`，挂上 `:272` | 变严（`to > from` 为假）⇒ 一路落盘；**变宽 ⇒ 每发都被 `ErrNoL2Confirm` 拒**（`composer_handlers.go:133-138`） |
| `Config`（248 的设置族，**不是 114 的三族之一**） | `*panel.ConfigWriteHandler`，同一枚 `Manager` ＋ `secret.NewStore` | `:262-270`，挂上 `:279` | `config.get`/`config.set` 真读真写磁盘 |
| `Workspace`（**工作区**） | `nil` | `:275` | `panel.workspace.request` ⇒ `ErrNoHandlerAttached`（`composer_dispatch.go:183-185`→`:230`） |
| `Attachment`（**附件**） | `nil` | `:276` | 同上，`:188-190` |
| `Message` | `nil` | `:277` | 同上，`:193-195` |

配套事实（工作区/附件两族的**原生腿本身已写好、只是没人调**）：`RequestWorkspaceSwitch`（`internal/panel/workspace.go:76`）、`DecodeAttachmentPayload`（`internal/panel/attachments.go:381`）、`(*AttachmentBroker).Ingest`（`:199`）——三枚在 `internal/ cmd/ tools/` 的非测试文件里**除注释外零调用点**（命中的 `bridge.go:17`、`composer_dispatch.go:74/:81`、`git.go:261`、`instructions_200.go:109`、`cmd/wisp/panel_pump.go:80` 全是注释）。`cmd/wisp/panel_pump.go:81` 逐字写着 "which the tree still has no caller for"。

⇒ **一句话答 §1**：入向那根线**早就有**（09-28 起），且**档位那一族已经通到磁盘**；票 114 题头那句把"零调用者"当今日事实的部分**过期**，但它点出的**能力缺口本身换了形**：今天面板能把档位**往严改**、**往宽改一律被拒且不会举起 L2 卡片**，工作区与附件两族仍是"在册、无处理器、响亮拒绝"。

---

## §2 票面每一格 AC 的三态判定

尺：票面 `- [ ]` 与 `- [x]` 两把各按"行首可有空白"的字符类跑，收尾 `| wc -l` ⇒ **未勾 9 枚 / 已勾 2 枚**（AC#1、AC#2 已勾；AC#3–AC#11 未勾）。题面那句"8 枚未勾格"**少计一枚**，见 §7-R1。

| 格 | 三态 | 凭据（`文件:行`）与差的那一句 |
|---|---|---|
| **AC#1** 现状表（已勾） | **已满足，但表里的结论今天已过期** | 交付件在（`docs/evidence/s1/114-ac1-status-table.md`）。它 ①.1 那句"WebView2 事件 → `ParseComposerRequest` 那一跳今天根本不存在"已被 §1 推翻；票面 `:37` 还据它判"AC#3 的反向对照在本段做不到"——**这一条前置理由今天不再成立**（跳已存在）。⇒ 勾不必重开，但**引用这张表的人必须重跑** |
| **AC#2** 原生侧门（已勾） | **已满足** | 门在产码里且在写之前：`internal/panel/composer_handlers.go:133-138`（`modeIsWidening(from,to) && h.Confirm == nil` ⇒ 拒，位置在 `:141` 那唯一一次 `Modes.Set` 之前）；判据是"注入的 `ModeWriter` 未被调用"那一族用例：`internal/panel/composer_handlers_test.go:120`（`TestAWideningModeRequestIsRefusedWhenNoL2ConfirmLegIsAttached`）、`:170`（正向一支，钉"不举第二张卡"）。变异可证伪：删 nil 检查 ⇒ `:120` 红；改成只记日志 ⇒ 同一枚的"未调用"判据红。**⚠ 唯一过期处＝该文件头 `:36-40` 自称"Nothing calls this handler yet"**——今天 `panel_inbound.go:272` 就在调它（这是注释缺陷，不是能力缺陷） |
| **AC#3** 反向对照：更宽档必须走 L2 原生卡片（300s＝拒）＋审计三档全写不因接线而漏 | **部分满足——两半里成立的是"审计"那一半，"卡片"那一半在生产里不可达** | 已成立：三档全写在**写腿**里（`internal/perm/store.go:172-214`，注释逐字"the audit log … records ALL THREE档"），且有常驻仪器（`internal/perm/store_test.go:164-166` 断言 3 枚 `MODE-SWITCH`）；300s 的缺省与钳位也在（`internal/agent/approval/queue.go:107 DefaultApprovalTimeout`）。**没成立**：面板那条链**没有一枚装配把非 nil 的 `Confirm` 交给它**——`cmd/wisp/panel_inbound.go:242` 与 `:252` 两枚 nil，常驻支复用同一函数体（`panel_resident_windows.go:170`）⇒ 变宽请求在 `composer_handlers.go:133` 就被"确认腿未接入"拒掉，**永远走不到 `Set` 里那张 C18 卡**。`cmd/wisp/run.go:659` 那枚 `rt.modeWrites` 确实拿到了真 `confirm`（`:637` → `rt.confirmModeSwitch`，注释在 `run.go:808`、函数体在 `:822`），但 **`rt.modeWrites` 全仓零读者**：尺 `modeWrites` 的非测试命中实测 **5 行**＝`run.go:316`（注释）＋`run.go:322`（字段声明）＋`run.go:659`（赋值）＋`panel_inbound.go:248`/`:272`（**那是 `newComposerDispatchChain` 自建的同名局部变量，不是 `rt` 的字段**）⇒ `rt.modeWrites` 那一支是死码。差的正是：**把 `newComposerDispatchChain` 的 `Confirm` 接成 `run.go:637` 已经握着的那枚同形，并让常驻腿有卡可举** |
| **AC#4** `R-92-1` 覆盖面：扩展名扩到与 d22scan 文本类一致 **或** 改挂结构 | **部分满足——"或"的第二支（结构形）已成立；名册那一支仍差一枚，且两枚名册之间没有尺在核相等** | 结构支已在：`internal/panel/composer_test.go:502 TestTheRendererHoldsExactlyOneDoorToTheHost` 把判定挂在**可达集**上（(i) 宿主调用点必须全在 `src/lib/panel.ts` 且恰 2 枚、(ii) 不许运行时拼点名 `:513`、(iii) 不许计算式路由 `:521` 的前一支、(iv) 每个 `panel.*` 字面量必须是 Go 真答的），自带种形正控 `:533 TestPlantedRendererDoorShapesGoRed`（一枚 `.js` 拼 `panel.mode.set`、一枚数组 `join(".")`）。门二也确实用上了宽名册：`rendererSourceFile`（`:370-376`）＝`.ts/.tsx/.js/.mjs/.jsx/.css/.html`，而 `scanComposerPermissionWrites`（`:597-627`）在**真树**上走的是**整棵 `frontend`**（调用点 `:138`，在 `TestPlantedComposerModeWriteGoesRed`（`:95`）里；种形那一发在 `:576` 传的同名根）。差的两句：① 票面点名的 **`.json` 仍不在** `rendererSourceFile` 里；② 反向也不齐——`tools/d22scan/main.go:1134-1136` 的 `isTextFile` 含 `.json` 却**不含 `.mjs`**⇒ 两枚名册谁也不等于谁，而**仓里没有一枚尺比较这两枚名册**（我逐枚找过：`grep -rn "isTextFile" tools/d22scan` 只有它自己那两枚文件用）。"不许为通过而删判定"与"验收方自己种的越界形状必须红"两句：正控在（`:533`），未删窄化（`:308-315` 那处刻意不扩的是**另一枚** git-switch 判据，注释具名解释为什么，⛔ 别把它读成 AC#4 没做） |
| **AC#5** 不越界：面板侧只许显示＋发起请求；工作区/档位两个输入口不得变成授权口 | **部分满足（Go 半支已满足；`frontend` 半支不由本腿判读）** | Go 半支凭据：`internal/agent/approval/ui.go:167-171` 的 `PanelAPI` 只有 `Reject/Head/View`，**没有 `Allow`**（`:164-166` 逐字"It has no Allow method…"）；composer 六门里没有任何一门能落授权——`panel.mode.request` 只能改档位且变宽被 §3 的门拒；`config.set` 明写拒 `risk.`/`fs.` 两族并点名各自欠的那张 L2 卡（`internal/panel/config_handlers.go:108-136`）；工作区/附件两族连处理器都没有。`frontend` 半支（出现 `approval.decide` 即 ban#6 红、零 emoji ban#8）＝**仪器射程**，判它绿要跑 `sh scripts/d22scan.sh`（§4-待验证窗口 W-3），且 ban#8 的射程按票 141 定案是"注释豁免、字符串不豁免"的那六段（AGENTS.md §1.2），**与本票面写的 `U+2190–U+2BFF` 起射程不是同一把** |
| **AC#6** 真机差分截屏 | **不由本编队做（界面委托）** | 票头 2026-09-23 09:22 横幅（票面 `:3-7`）＋台账 `A102`：`skipped=frontend(owner-delegated)`，保持未勾，本编队⛔ 不替它勾、⛔ 不"只改一角"。本腿对 `frontend/**` 零读取，未据任何页面内容下判 |
| **AC#7** 门禁：d22scan 纯净快照 rc=0 且不降 / `gofumpt -l` 真跑贴原始输出 / POSIX 三合包按 `-v` 口径 | **完全没做（就本轮而言），且没有任何一件现存读数可借** | 票面要求"必须真跑并贴原始输出"。248-r1 的证据件 §4 **没有写**（`docs/evidence/s1/248-settings-write-path-r1.md:147-149` 那一节仍是占位）⇒ 连 248 那批改动的门禁色都没落档，114 更无读数控。本腿按硬边界不跑任何东西 ⇒ 见 §4 的待验证窗口 W-1…W-5（命令给全、期望读数写明） |
| **AC#8** 把"两边同改"从承诺变成门：往 `knownComposerMethod` 加一条必须同时让 `frontend/` 出现对应字面量 | **部分满足——一个方向有门，另一个方向只对 `panel.` 前缀族成立，而 248 新增的两枚名恰好掉在这两把尺之外** | 已有：`internal/panel/bridge_test.go:131-145 TestFrontendComposerRequestsMatchTheEnvelope` **逐枚**要求四枚 `Method*` 常量在 `frontend/src/lib/panel.ts` 里以字面量出现（`:142` 那句 `"the Go side parses a method nobody sends"` 就是 AC#8 要的牙齿）；反方向在 `composer_test.go:521`（页面点名而 Go 不答 ⇒ 红）。**差三句**：① `bridge_test.go:138-141` 那份清单是**手抄的四枚**，新增一枚常量而忘了抄进清单 ⇒ 没人逼它出现在页面里；② `routeLiteralRe`（`composer_test.go:394`）只抓 `"panel.…"` 形 ⇒ **`config.get`/`config.set` 这两枚不带 `panel.` 前缀的名，两把尺都看不见**（248-r1 自己在 `docs/evidence/s1/248-settings-write-path-r1.md` §5-12 承认没登记进 `composerRouteLiterals`）；③ 真正逼住"别扩名册"的是另一枚钉——`internal/panel/git_test.go:381-388` 与 `:517` 用 `panelMethodRe`（只匹配 `"panel.*"`）断言 `bridge.go` 里恰有 4 枚 ⇒ **加第五枚 `panel.*` 会红在那句"this ticket may not extend"，加 `config.*` 不会**。⇒ AC#8 的洞今天**变宽了**而不是变小 |
| **AC#9** 结构钉的扫描根比 ban#6 与"渲染器实际加载的文件集合"都窄 ⇒ 把可达性钉挂在那个集合上 | **完全没做** | 那枚可达性钉的根写死在 `internal/panel/composer_test.go:504`：`filepath.Join(root, "frontend", "src")`，且 `:434` 跳过 `dist`；门二走的是整棵 `frontend`（真树调用点 `:138`）⇒ **"钉的根"与"门二的根"本来就是两棵树**，票面 `frontend/public/*.js` + `index.html` 引用那一形今天仍然只被门二读到、不被可达性钉读到。"渲染器到底加载哪些文件"这一问在本仓**没有答案件**，而要答它得读 `frontend/**` ⇒ 归界面侧/编排者（§6-J3） |
| **AC#10** `panel.ts` 自陈宿主有两种装法（`AddHostObjectToScript` / `postMessage`），钉只认后者 ⇒ 第一种要么显式禁用要么进同一枚钉 | **完全没做，且现实比票面写的更糟：今天生产用的是票面没列的第三种装法** | 现量：`grep -rn "AddHostObjectToScript" --include=*.go internal/ cmd/ tools/` ⇒ **零命中**（模块里只有 vtbl 声明 `pkg/edge/corewebview2.go:123`，无任何调用点）⇒ 第一种装法**没在用**。真正在用的是：`cmd/wisp/panel_host_windows.go:319 w.Bind("wispDispatch", …)`＝go-webview2 的 `Bind`，它注入的存根走 **`window.external.invoke`**（模块 `webview.go:450-482`），宿主按 `{id,method,params}` 解 RPC（`webview.go:139-168`）。而可达性钉的宿主调用点正则（`composer_test.go:381 hostBridgeCallRe`）**只认 `.postMessage(`**⇒ **今天这扇真门对那枚钉是隐形的**。⚠ 还有一条本腿判为**新增缺陷**：页面若沿用旧形直接 `postMessage(那枚封套)`，库会把 `method` 读成 `"panel.mode.request"`、查 `bindings` 查不到、`callbinding` 返回 `nil,nil`（`webview.go:166-168`）⇒ **promise 静默 resolve(null)，Go 侧一行审计都不写**——正是"被静默丢掉"那一形 |
| **AC#11** 把 `npm run render:composer && git diff --exit-code` 做成 CI 一步 | **部分满足——前半已在 CI，后半（防腐坏那半）整枚没做** | 前半：`.github/workflows/ci.yml:718-721` 确有 "Composer states paint the real envelope (render evidence)" 一步，`run: npm run render:composer`。后半：**`grep -n "exit-code" .github/workflows/ci.yml` ⇒ 零命中**，那一步后面没有 `git diff --exit-code`，所以"fixture 可以手写/可以腐坏"这一格今天**照旧开着**。票面 `:76` 自己要求"动 `ci.yml` 之前先看票 111/85a 是否已让出该文件"⇒ 落地腿开工前先查让渡（§6-J5） |

**一句话答"有几枚被 248-r1 顺手做掉了"**：**一枚都没有整格做掉。** AC#4/AC#8/AC#10 三格今天**被 248-r1 动了但方向不利**——它把两枚不带 `panel.` 前缀的名放进了名册，而这三格的全部仪器都以 `panel.` 前缀为锚，于是覆盖面被削薄（AC#8 尤其：它要的正是"名册扩了、页面必须跟着"，而扩了的这两枚恰好没人能逼）。AC#3 是**唯一一处 248 的工作让 114 变近**的：`Confirm` 那个槽位与"回执要能到页面"的形状（`composer_dispatch.go:110-112` 带 raw 的签名、`Handle` 返回串）已经在，114 只差把 `run.go:637` 那枚真 `confirm` 接进 `panel_inbound.go:242/252`。

---

## §3 "门必须放原生侧"这条前置条件现在在哪

### 3.1 结构性禁区：在，而且是"方法根本不存在"那一形

- `internal/agent/approval/ui.go:167-171` —— `PanelAPI` 的全部方法就是 `Reject(correlationID, reason)`、`Head()`、`View(correlationID)`。作者注释在 `:164-166`，逐字（源文件用的是直引号，此处照抄）：**PanelAPI is the ENTIRE surface a panel/WebView host may be handed. It has no Allow method: "panel-sourced allow is structurally rejected" is then not a check that could be forgotten, it is a method that does not exist.**
- `Allow` 与 `AllowSession` 只活在 `NativeAPI` 上（`ui.go:143-162`，`:155-157` 逐字"It exists on this interface and on no other surface: PanelAPI has no Allow method at all, so a panel-sourced host cannot store an authorization because there is no method for it to call."）。兄弟锚点：`approval.go:224`、`queue.go:358`、`replies.go:347`（"There is deliberately no PanelAllowSession"）、`gate.go:620`（`Gate.Panel()` 是分派口）。
- ⚠ 两处**过期行号指针**（本腿现量，别照它们找）：`cmd/wisp/approval_reply.go:29` 写 "ui.go:45-48"、`internal/agent/approval/gate.go:632` 写 "ui.go:152-159" ⇒ 真身都在 `ui.go:164-171`。禁区本身**在**，只是被引错了行。
- **今天这条禁区有多"空转"要先说清**：尺 `-E "\.Panel\(\)" --include=*.go`（作用域 `internal/ cmd/ tools/`，去掉 `_test.go`）实测 **4 行**，其中两行是注释（`internal/agent/approval/doc.go:33`、`cmd/wisp/approval_reply.go:21`），**真调用只有 `internal/agent/approval/replies.go:478`（`g.Panel().Head()`）与 `:487`（`g.Panel().View(corr)`）两枚读侧包装**；`cmd/wisp` 里**产码调用一枚都没有** ⇒ 常驻面板今天**连那枚"只能拒绝"的面都没拿到**。所以"门在原生侧"今天是**由构造成立**（没有方法可叫），不是"面板走到门口被拒"。这一点对 §5 很重要：**接线不会把禁区变松，但它会把禁区从"没接线所以不可能"变成"接了线且仍然不可能"——后者才是 114 想要的形状**（R-92-2 的原意）。

### 3.2 冻结件那族仪器今天盖不盖得到 248-r1 新增的两枚入向名

**（射程判断，非内容引用；本腿只读 `l2_grant_boundary_test.go` 内部来判它看得见什么。）**

| 仪器（`文件:行`） | 读什么 | 今天对 `config.get`/`config.set` 的射程 |
|---|---|---|
| `routeShapedName`（`:1050-1069`） | 点分段名，不限 `panel` 前缀 | **认得**（`config.get` 是两段纯字母）⇒ 这两枚名会进名册池 |
| `poolJudgedByRealGuard`（`:1124-1143`） | 拿池里每枚名去问**正在跑的** `knownComposerMethod` | **盖得到**：`answeredElsewhere`（守卫答了而派生枚举没看见）与 `answeredGrantDoor`（守卫答了而名字本身是裁决形）两支持都对 `config.*` 生效 |
| 封套字段扫描（`:1180-1198`，`carriesGrantWord(f.JSONKey, grantFieldWords)`） | `ComposerRequest` 及其嵌套 JSON 目标 | **盖得到**新增的四枚键（`configField`/`configProvider`/`configModel`/`configValue`）；四枚都不含裁决词 ⇒ 今天不报，**但一旦有人把凭据值或"approve"形键加进这枚共用封套就会红**（这正是 248-r1 §1.1 那支 J4 甲的理由） |
| 名册全等锚（`:2206` + `:2308`，`guardRosterOf` `:2030`） | `bridge.go` 的守卫 case 清单 ↔ 本文件声明的 `Method*` 常量，**双向全等** | **盖得到**（248 正是照这个形状把两枚名同时加进常量与 case 的）。⚠ 这两枚锚是 `A487` 具名解冻的唯一两枚 |
| 三因子扫掠网格（`:1316-1324`：`grantRoutePrefixes` 8 枚全是 `panel*` 族 × `grantRouteWords` 11 枚 × `grantRouteSuffixes` 3 枚） | 运行时**拼**出来的候选授权名，逐枚问真守卫 | **盖不到**：网格里没有 `config` 前缀族 ⇒ `config.approve`／`config.set.allow` 这类"新命名空间的授权形"不会被**扫掠**问到。兜住它的是上面第 2 行那枚池仪器（名字必须写在包里才能被答），**不是**这枚网格 |
| 反空转正控（`:1496-1502`，`wantAnswered` ＝ 恰好那四枚 `panel.*`） | 从 `{panel}×{mode,workspace,attachment,message}×{.request,.add,.send}` 24 名里派生"守卫真答的" | **盖不到 `config.*`**（前缀只有 `panel`）⇒ 248 加两枚名没让它红；**反过来它是一枚硬阻塞**：见 §4-N8 |

### 3.3 三枚授权词表 vs 114 落地可能新增的名字

**先给结论：票 114 的落地一枚新名都不需要。** 它那三族的入口名今天都已在名册里——`panel.mode.request`（`bridge.go:42`）、`panel.workspace.request`（`:43`）、`panel.attachment.add`（`:44`）。114 的差集是"槽里有没有处理器"、"处理器手里有没有确认腿"，不是"有没有名字"。

逐枚比（若某枚写腿还是要新造名，会不会被认成授权形）：

| 候选新名（例） | `grantRouteWords`（`:192-195`：approve/approval/grant/allow/permit/ratify/authorize/authorised/decide/decision/verdict） | 后果 |
|---|---|---|
| `panel.mode.get`、`panel.mode.apply`、`panel.workspace.switch`、`panel.attachment.read` | 不含任何一枚词 ⇒ **不会**被读成授权形 | 但**仍会**被 §4-N4/N5 那枚"名册只许四枚 `panel.*`"的钉打红，以及（若落在 24 名网格内）§4-N8 |
| `panel.mode.approve`、`panel.workspace.allow`、`panel.attachment.grant` | **含词** ⇒ `answeredGrantDoor` 与守卫候选两支同时红 | 这正是设计要的：这类名根本不该出现 |
| `config.mode.set`、`config.approve`（不带 `panel.` 前缀那一族） | 词表照样判（`carriesGrantWord` 看的是**名字本身**）⇒ `config.approve` 会红 | ⚠ 但**扫掠网格问不到它**（3.2 末行），且 §2 AC#8 那两把前缀尺（`routeLiteralRe`、`panelMethodRe`）也看不见它 ⇒ **新命名空间是今天最没人看着的那道门** |

⛔ 本腿明确建议：**不要为了躲词表去动 `grantRouteWords`/`grantRoutePrefixes`/`grantRouteSuffixes` 那三枚表**（它们在冻结件里，且 `:1333-1337` 的三枚尺寸钉＋`:1487-1491` 的乘积钉会把"删一枚词顺便删它的见证行"那一形单独抓到）。

---

## §4 撞钉预检（名册，不是颜色）

⛔ 本节只报名册与断言原文。**今天绿不绿一律落在下面的"待验证窗口"**（本腿零 `go test`／零 `go build`／零 `go vet`／未执行任何 exe）。

### 4.1 会挡住 114 落地的现存钉

| # | 用例名 | `文件:行` | 断言（原文要点） | 这枚钉的存在理由 | 对 114 的后果 |
|---|---|---|---|---|---|
| **N1** | `TestAC9InboundLegRefusesRosterMethodWithNoHandler` | `cmd/wisp/panel_inbound_33_test.go:158`（`:166` `exit=1 want 1 (the request was refused)`、`:168` stdout 须含"处理器未接入"、`:171` stderr 须含 `panel.workspace.request`） | 在册而本机没接处理器的名**必须点名拒**，不许 nil 解引用也不许静默 | `runInboundLeg33` 起的是**真听众**，读的是真链 | ⚠ **这枚就是"工作区/附件两族不许由 114 顺手接"的看门钉**：一旦 `newComposerDispatchChain`（`panel_inbound.go:228`）把 `Workspace` 从 `nil`（`:275`）改成有处理器，它**必红**。只接 `Confirm` 不动这两枚槽 ⇒ 不受影响 |
| **N2** | `TestAC9InboundLegFromStdinReachesTheWriteLeg` | `cmd/wisp/panel_inbound_33_test.go:99`（`:100` 种 `auto_approve`、`:105` 请求 `to":"ask_every_step"`；`:110` `exit=0`、`:113` "已受理"、`:116` 磁盘须变 `ask_every_step`） | 一条 stdin 封套必须**真的把磁盘上的档改掉** | 钉"入向不是摆设"（33 AC#9 的判据①） | 它刻意挑的是**变严**那一支（绕开 `Confirm==nil` 的门）。114 接上 `Confirm` 后仍应绿；**把方向律改成"任何变更都要卡"则红** ⇒ 方向律（只挡变宽）是编排者 P1 已裁的，别重开 |
| **N3** | `TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket` | `cmd/wisp/panel_config_248_test.go:457`（循环 `:464` 三枚名，断言 `:466` 错串须含"处理器未接入"） | 加了设置槽之后旧的三枚槽仍要点名拒 | 看着像阻塞、**其实不是**：它调 `newRealSettingsLeg`（`:123`，`:142` 自己 `&panel.ComposerDispatch{...}`）现造派发器，**读的不是生产链**（注释 `:462-463` 那句 "exactly as the production chain builds them" 是**陈述，不是断言**） | ⚠ **具名报出来是为了别让下一枚腿把它误读成 N1 的替身而绕开 N1** |
| **N4** | 名册支（git 维度那族用例里的白名单断言） | `internal/panel/git_test.go:381-388`，断言 `:385` `!equalStrings(got, wantMethods)`，报错原文含 "want the four methods **this ticket may not extend**" | `bridge.go` 里 `"panel.*"` 字面量必须恰是那四枚 | 钉住 C17 名册不被人随手加长 | ⛔ **加第 5 枚 `panel.*` 常量 ⇒ 红**（取名的正则在 `git_test.go:394` 的 `panelMethodRe`，只匹配 `"panel.*"` ⇒ `config.*` 根本数不到） |
| **N5** | 同一把尺的正控门 3 | `internal/panel/git_test.go:505-519`：`:514` 种第五枚名要求数到 5、`:517` 真 `bridge.go` 要求恰为 4 | 尺本身有牙齿（种下去必被抓） | 防"绿的因为瞎" | 与 N4 同发红；**两枚都得改才能扩名册**，而那枚文件是票 92 地界 |
| **N6** | `TestTheRendererHoldsExactlyOneDoorToTheHost` | `internal/panel/composer_test.go:502`，四支 `:505`（`rep.outside != 0`）、`:509`（**`rep.callSites != 2`**）、`:513`（`assembly != 0`）、`:521`（`unknown != 0`）；扫描根写死在 `:504` 的 `frontend/src` | 页面能敲宿主的**只有那两枚调用点**，且只能点在 Go 真答的名上 | `R-92-1` 的结构形正解 | 页面侧新增宿主调用（界面侧的活）或新增未登记的 `panel.*` 字面量 ⇒ 红。**它对 `window.wispDispatch` 那一形无感**（`hostBridgeCallRe` `:381` 只认 `.postMessage(`）⇒ 见 §5-丙 形一 |
| **N7** | `TestFrontendComposerRequestsMatchTheEnvelope` | `internal/panel/bridge_test.go:131`：`:142` 逐枚要求四枚常量以字面量出现在 `frontend/src/lib/panel.ts`；`:157` 禁 `window.postMessage(`/`new WebSocket`/`fetch(`；`:161` `strings.Count(text, "bridge.postMessage") != 2` | AC#8 要的"两边同改"实际就落在这枚上 | 页与 Go 不许漂 | **只遍历 `:138-141` 手抄的四枚** ⇒ 新常量不抄进清单就不被要求出现在页面里（§2 AC#8 差的第①句） |
| **N8** | 扫掠网格的反空转正控 | `internal/panel/l2_grant_boundary_test.go:1496-1502`：`wantAnswered := []string{MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend}` ＋ `reflect.DeepEqual(legit, wantAnswered)`；派生网格在 `:1459-1461`（`{panel}×{mode,workspace,attachment,message}×{.request,.add,.send}`＝24 名） | 扫掠必须问得出"守卫确实答的那几枚"，否则 `hits=[]` 只是网格空了 | **冻结件**：防网格被悄悄缩小 | ⚠ **硬阻塞**：任何新名一旦落进这 24 名的乘积里（例 `panel.attachment.request`、`panel.workspace.add`），`legit` 变 5 ⇒ **冻结件红**，而这一格不是能"顺手改"的 |
| **N9** | 网格尺寸钉＋乘积钉 | `internal/panel/l2_grant_boundary_test.go:1333-1337`（`wantGrantRoutePrefixes=8`/`Words=11`/`Suffixes=3`）、`:1487-1491`（ lists 必须等于写死的数）、`:1504-1506`（`asked` 必须等于三因子乘积×2） | 因子被动过就必须同步写死 | **冻结件**（F-R4-2 / F-ACC-1） | ⛔ 114 不该碰；碰了必红（AGENTS.md §1.1"一字节都不许动"） |
| **N10** | 名册全等锚（`A487` 唯一解冻的两枚） | `internal/panel/l2_grant_boundary_test.go:2206`（`guardRosterOf`，函数在 `:2030`）＋ `:2308`（`sortedSet(pkg.answered)` 与派生名册**全等**）＋自带正控 `:2219` | 守卫 case 清单 ↔ `bridge.go` 的 `Method*` 常量，双向 | 名册只扩一侧＝漂 | 114 若新造名而不同时改两处 ⇒ 红；只改 `case` 不加常量 ⇒ 正控那一支也红 |
| **N11** | `TestAWideningModeRequestIsRefusedWhenNoL2ConfirmLegIsAttached` | `internal/panel/composer_handlers_test.go:120`（判据＝注入的 `ModeWriter` **一次都没被调用**）；配套 `:170` `TestAWideningRequestWithAnAttachedLegReachesTheOneWriterAndRaisesNoSecondCard` | AC#2 的门本身，变异可证伪（票面 `:40` 要求） | 防"只记日志不拒"那一形 | 把 `Confirm` 接进生产链**不会**红它（它自建处理器）；**把 `composer_handlers.go:133` 那支删/软化会红** |
| **N12** | `TestAC1SessionDisposeHasAProductionTrigger_AC1` | `cmd/wisp/panel_host_gate_test.go:177`，`:181`→`:182` 只在**无**产码构造点时才 skip、`:184`→`:185-186` `teardownHits==0` ⇒ `t.Errorf`；工厂名册 `:194-197` 认 `NewPanelManager`/`newResidentPanelManager` | 造了宿主就必须有"会话拆窗"那处 | 33-r5 之后这枚**已翻面成实断言**（`cmd/wisp/resident_windows.go:142` 是产码构造点 ⇒ `:182` 那句 skip 今天不吃） | 114 若新增第二枚 `PanelManager` 装配点而没带 manager 上的 `Destroy` ⇒ 红 |
| **N13** | 数据根腿名册表 | `cmd/wisp/dataroot_128_test.go:113`（"每张会调 `resolveDataDir` 的腿都要在册"） | 没根就不许写、且要说出理由 | 票 128 的地界 | `panel_inbound.go` 今天**刻意不解析数据根**（`:122-137`，退码 2）；114 若让常驻/新腿去解析根 ⇒ 要改**别人票名的文件**（未定义即停，见 §6-J1） |
| **N14** | 派单普查闸门 | `cmd/wisp/leg_dispatch_gate_133_test.go`（对照 `cmd/wisp/main.go:86-130` 的 `case` 与 usage 块；`WISP-LEG-COVERAGE-RULING` 三条写在 `main.go:74-85`） | 被派发的命令必须被钉住／被驱动／被具名排除 | 票 133 的闸门 | 114 若新开子命令（例 `wisp panel-mode`）⇒ 必须同时进 usage 块与本表 |
| **N15** | 常驻腿"config.toml 必须不存在"那枚钉 | `cmd/wisp/resident_task_source_246_windows_test.go:389-391`（⚠ 本腿**未复跑**、亦未逐字复读该文件；行号与判语取自提交 `d0cd3a9b` 的说明，属**待验断言**） | 首建 `config.toml` 不许进常驻那条路 | 票 198/246 的地界（编排者 10-02 裁"首建只挂 `wisp run` 入口，⛔ 不进 `assembleRuntime`"） | 114 不建文件 ⇒ 不吃这枚；**若为"面板改不了档"而去先落一个 `config.toml`，一发就红在这里** |

### 4.2 待验证窗口（本腿一枚都没跑）

| # | 命令（逐字；先 `export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`） | 期望读数 | 为什么必须真跑才知道 |
|---|---|---|---|
| **W-1** | `go test -count=2 -v ./internal/panel/` | `--- FAIL` 计数＝**0**（票 114 Progress log `:101` 用的口径是四数 `=== RUN / --- PASS / FAIL / SKIP`，现值应为 `152/92/0/0` **量级**，⛔ 别拿它当今日读数——那是 09-23 的，248 之后必然不是这组数） | AC#7 的颜色只能真跑 |
| **W-2** | `go test -count=2 -v ./cmd/wisp/` | `--- FAIL` 计数＝**0** | ⛔ **缺那两行 PATH 会给 `exit status 0xc0000135` 且没有一行 `--- FAIL`＝用例根本没跑**，极易被误读成"全绿"（本仓 09-28 已踩过并记档） |
| **W-3** | `sh scripts/d22scan.sh` | `rc=0`，且各 scope 计数不降 | AC#7 第一支；ban#6/ban#8 的颜色 |
| **W-4** | `"$(go env GOPATH)\\bin\\gofumpt.exe" -l internal/panel cmd/wisp` | **零行输出**（⚠ 票面 `:78-81` 已记：Git Bash 下 `export PATH="$PATH:$(go env GOPATH)/bin"` 会被 MSYS 吃成 `\work\base\gopath/...` ⇒ `rc=127` 是**假读数**，必须写 `.exe` 全路径） | AC#7 第二支明令"必须真跑并贴原始输出" |
| **W-5** | 最小变异一发：把 `cmd/wisp/panel_inbound.go:242` 与 `:252` 两枚 `Confirm: nil` 换成 `run.go:637` 那枚同形，再跑 W-1 ＋ `go test -run 'TestAC9' -v ./cmd/wisp/` | 本腿**现量预测**（不是读数）：N1 不受影响（它没动那两枚槽）、N2 仍应绿（它走变严支）、`composer_handlers_test.go:120` 仍绿（自建处理器）；**若同一发里还接了 `Workspace` ⇒ N1 必红** | 这一步就是 114 落地的最小形状，颜色只能真跑；本腿⛔ 未改任何产码 |

---

## §5 落点建议（只给形状）

**本腿选：乙（带一处具名修正）；丙不是落点，但丙那一形今天真实存在，位置与票面写的不是一处。**

先给选乙的那句硬事实：**114 的三族里只有一族真接到了磁盘（档位），两族（附件、工作区）到今天还是"在册、槽 nil、点名拒"；而 248-r1 动的是第四族（设置），它不在 114 的三族里。** 所以"下一枚腿是新建接线"这个读法是错的——**是补差集**，而且差集只剩一枚小口子。

### 甲档：114 已被 248 顺手做掉大半 —— **不选**

- **判死凭据**：票面 9 枚未勾格里**没有一枚被整格做掉**（§2 逐格）。248 做的是：名册扩两枚（`bridge.go:66-67` 常量、`:146-152` 守卫）、派发表加两枚 case 与 `Config` 槽（`composer_dispatch.go:197-206`、`:134`）、设置处理器与落盘腿（`internal/panel/config_handlers.go`、`cmd/wisp/panel_config_store.go`、`internal/config/settings.go`）。**没有一枚是 114 的 AC 要的东西。**
- 甲档唯一的真金是**形状**：`ConfigRequestHandler` 那枚"带 raw、返回给页面的句子"的签名（`composer_dispatch.go:110-112`），以及 `Handle` 现在确实会把非空回执交回 bind 闭包（`panel_host_windows.go:319-322`、`cmd/wisp/panel_inbound.go:180-183`）——**AC#3 想要的"回执要到页面"那半条路已经铺好**，只差档位在 `Handle` 的返回串里没人写（`dispatch` 对那四枚旧门一律返回空串：`composer_dispatch.go:181`、`:186`、`:191`、`:196`）。
- **选它的代价**：把下一枚腿派成"复核票 114 已完"，复核回来仍是 9 枚未勾，白花一枚腿。
- **选错了在盘上会留下什么形状**：票面被勾上而判据不成立的格（AC#3 与 AC#8 首当其冲）——**这正是票 114 自己立案的理由**（票面 `:29-30`："我们等于把'绿的原因'换成了'绿的样子'"）。

### 乙档：入向那条线存在，但档位／附件／工作区三族仍没接 —— **选它，改正一处**

- **需改正的那半句**：**档位已经接了**。`newComposerDispatchChain` 把 `Mode` 槽装上真处理器（`cmd/wisp/panel_inbound.go:248-255`、`:272`），常驻腿与 CLI 腿共用同一函数体（`cmd/wisp/panel_resident_windows.go:166-170`），链顶端是产码听众（§1 H1/H2）。今天面板**能把档往严改并落盘**，**往宽改每发都在 `composer_handlers.go:133-138` 被 `ErrNoL2Confirm` 响亮拒掉**。
- **每一族的落点（只给形状，不给实现）**：
  - **档位＝114 的真差集，最小一发**：`cmd/wisp/panel_inbound.go:242`（`perm.Options.Confirm`）与 `:252`（`ModeWriteHandler.Confirm`）两枚 nil ⇒ 接成 `cmd/wisp/run.go:637` 已经握在手里的那枚同形（`rt.confirmModeSwitch`，定义 `run.go:808`、举卡在 `:822`，它举的就是那张 C18 卡）。⛔ **别顺手把 `run.go:659` 那枚 `rt.modeWrites` 也接上**——它今天零读者（§1、§7-R3），接它的正确方式要么是把它变成链的来源、要么具名删掉，两种都是**编排者的一语**（§6-J1）。
    同一发必须一并处理的那件事：**常驻腿今天没有卡可举**。`newComposerDispatchChain` 自称"是管道不是窗口"（`panel_inbound.go:32-39`），而 `Gate.Panel()` 在 `cmd/wisp` 里零产码调用者（§3.1）。⇒ 要么把审批面的句柄注进常驻腿（撞 §6-J1 的"未定义即停"，且票 246 的 A481 裁的是"门由装配根构造后以函数值下发，常驻文件不自造门"），要么**具名停在"面板只能改严"这一档并写进票面**。**我不替它选。**
  - **附件**：槽 `panel_inbound.go:276`；处理器**在本仓不存在**（接口在 `composer_dispatch.go:83-85`）；原生腿已写好且零调用者——`DecodeAttachmentPayload`（`internal/panel/attachments.go:381`）、`(*AttachmentBroker).Ingest`（`:199`）。**归属不是 114**：`panel_inbound.go:54-58` 逐字写 attachment＝票 92 的地界。⇒ **114 不该接它。**
  - **工作区**：槽 `panel_inbound.go:275`；处理器归**票 186**（`composer_dispatch.go:73-78` 逐字 "the handler in front of it is ticket 186's work, so this file declares the socket and nothing more"）；原生腿已写好且零调用者——`RequestWorkspaceSwitch`（`internal/panel/workspace.go:76`）。⇒ **114 不该接它，接了必红 N1。**
  - （附带一族题面没列：**消息**，槽 `:277`、处理器归票 35。同样不该由 114 接。）
- **选它的代价**：下一枚腿只动 `cmd/wisp/panel_inbound.go` 两枚 nil ＋ 常驻腿的确认腿来源 ＋ AC#3 的反向对照用例 ＋ AC#7 的四数读数。范围小，但**代价要说全**：**票面 9 枚未勾里有 5 枚是仪器格（AC#4/AC#8/AC#9/AC#10/AC#11），不会因为这一发变绿**，另有 1 枚（AC#6）永不归本编队 ⇒ **一枚腿清不了这张票**。
- **选错了（把它读成"三族都得由 114 接"）在盘上会留下什么**：`Workspace`/`Attachment` 两枚槽被填 ⇒ **N1 红**，而修 N1 要改票 33 的用例文件（别人票名的钉）⇒ 演变成"为实现放宽别人的判据"，AGENTS.md §1.1 明令禁止那一形。

### 丙档：接线存在但门的位置错了 —— **不作落点；今天真实的缺陷在"守门尺看不见那扇门"**

- **门本身没错**：mode 的门在原生侧、在任何一次 `Modes.Set` 之前（`composer_handlers.go:133-141`）；授权面仍由"方法不存在"守着（`ui.go:167-171`）。**没有第二道门，也没有把门挪到页面侧的形状。**
- **错的是"看门的那把尺对着哪扇门"，两形具名**：
  1. **形一**：可达性钉的宿主调用点正则 `hostBridgeCallRe`（`internal/panel/composer_test.go:381`）＝ `\.\s*postMessage\s*\(`，而生产那扇门是 `w.Bind(panelDispatchBinding, …)`（`cmd/wisp/panel_host_windows.go:82`、`:319`，绑定名逐字 `wispDispatch`）＝go-webview2 的 `Bind` 注入 `window.wispDispatch(...)`，实参经 **`window.external.invoke(JSON.stringify({id, method, params}))`** 送出（模块 `webview.go:450-482`），宿主侧 `msgcb` 按 `rpcMessage` 解包再查 `w.bindings[d.Method]`（`webview.go:139-168`）。⇒ **页面无论怎么调 `window.wispDispatch`，N6 的 `callSites` 恒为 2、`outside` 恒为空——那枚"渲染器只留一道门"的钉对它自家仓里真开的那道门无感。** 票面 AC#10 说"两种装法、钉只认后者"，实测比它更糟：仓里用的是**第三种**（`AddHostObjectToScript` 在 `internal/ cmd/ tools/` 里零调用点，只有模块 vtbl 声明 `pkg/edge/corewebview2.go:123`）。
  2. **形二**：页面若照票 92 那批文档的形状直接 `postMessage(那枚 composer 封套)`，库先按 RPC 解包，`d.Method` 会是 `"panel.mode.request"`，`bindings` 查不到 ⇒ `callbinding` 返回 `(nil, nil)`（`webview.go:162-168`）⇒ **promise 以 null 落定，Go 侧一行审计都不写**（`composer_dispatch.go:240-247` 的 `record`、`composer_handlers.go:158-165` 的 `MODE-REFUSED` 全都到不了）。⇒ 票 114 全篇防的那件事——"一次请求被丢掉而看起来像被处理了"（`bridge.go:154-157` 的 `RefusedEnvelopeForUser` 立论）——**在宿主库那一层原样存在，且不在任何一把尺的射程里**。
- **为什么它不能当落点**：两形都不是 `cmd/wisp`/`internal/panel` 里"接线"能修完的。形一要么把 `wispDispatch` 这枚绑定名纳入 N6/N7 的可达集判定（那是 AC#9 的扫描根同一格，且要读 `frontend/**` 才答得全，§6-J3），要么具名禁用 `Bind`；形二要么在宿主外面加一层"非 RPC 形状的消息必须落审计"的收口（票 33 的 `PanelManager` 地界），要么登记成 DEFERRED（五字段齐全，SPEC-12 §5）。**两条都超出"补 114 差集"的范围，按 AGENTS.md §2 该上报而不是自选。**
- **如果编排者按"接线完了"结案 114 而没处理丙**，盘上留下的形状：票面被勾、N1–N11 全绿、`sh scripts/d22scan.sh` `rc=0`，**而"面板改档"唯一的防线是 `composer_handlers.go:133` 那枚 nil 检查，它守的是一条页面上没人真调得出的通路**——即票 92 验收当年打回的同一种"绿的样子"。

---

## §6 判不动的地方

| # | 判不动的是什么 | 为什么判不动 | 谁能判 |
|---|---|---|---|
| **J1** | **常驻面板拿什么举那张 C18 卡**（AC#3 的正身） | 常驻腿的 `Confirm` 是 nil 且**没有卡片面**（`panel_inbound.go:32-39` 自称"是管道不是窗口"；`Gate.Panel()` 非测试零命中）。要接上就得给常驻进程注一枚审批面柄，而票 246 的裁决（台账 `A481`）明写"审批门由装配根构造后**以函数值与绑定下发**，常驻文件不自造门"。注入的**那一枚形**（`Gate.Panel()` 的句柄？`rt.confirmModeSwitch` 的转发？还是票 224/246 那族卡？）在 spec 与票面里**没有一处被具名过**。我挑就是在造契约 | 编排者（照 A481 那一族口径裁一句），必要时 owner |
| **J2** | **AC#4 两支名册要不要"逐字节相等"**（`.json` 与 `.mjs` 那两枚差） | 票面写的是"要么…要么…"，结构支已成立（§2），故我判"部分满足"。但"扩展名支要不要补齐"是**判据松紧**：`rendererSourceFile` 里的 `.mjs` 是**刻意加的**，而 `composer_test.go:308-315` 那段注释解释了**另一枚**判据（git-switch 那一族）为何刻意不扩到 `.mjs`（扩了会打在一枚构建脚本的注释上）；d22scan 的 `isTextFile`（`tools/d22scan/main.go:1134-1136`）有 `.json` 无 `.mjs`。**两枚名册谁也不等于谁，仓里没有一枚尺比较它们**——该不该新造这枚尺不是我判的 | 验收腿（跑 W-3 看覆盖面）＋编排者裁"要不要新增一枚对齐尺" |
| **J3** | **"渲染器到底加载哪些文件"（AC#9 的那一问）** | 要答就得读 `frontend/**`：入口 HTML 引了哪些脚本、`public/` 里有没有真会被加载的 `.js`、`dist` 产物形状。⇒ **票头 09-23 横幅＋本派单硬边界②双重禁止**。我能确证的只有 Go 侧那两枚根本不同（可达性钉读 `frontend/src`，门二读 `frontend`） | 界面侧那枚 agent（owner 已把前端整块交出，`A102`）；编排者据其答复定扫描根 |
| **J4** | **AC#5 与 AC#7 的颜色**（ban#6/ban#8 是否绿、`gofumpt`／d22scan／POSIX 三合包数字） | 派单硬边界①：零 `go test`／零 `go build`／零 `go vet`／不执行 exe；且此刻 `33-r8`（`internal/ball`）、`250-r1`（`scripts/`）、`198-r1`（`cmd/wisp`）三枚写码腿在飞、跑的是计时敏感用例。⇒ 只有真跑才有颜色，读数需求在 §4.2 | 落地腿自己跑并贴原始输出（票面 `:55` 明令），验收腿复跑 |
| **J5** | **AC#11 动 `ci.yml` 的前置** | 票面 `:76` 自带一句"动 `ci.yml` 之前先看票 111/85a 是否已让出该文件"。这份**让渡**是台账/票面里的协商，我只读得到：`.github/workflows/ci.yml:718-721` 那一步已存在，且**整个 ci.yml 零处 `git diff --exit-code`**。让没让出不是我判的 | 编排者（查 111/85a 让渡）后再派 |
| **J6** | **`wisp panel-inbound` 那条 CLI 腿算不算"面板输入通路"的证据** | 它是**操作员从 stdin 敲封套**（`panel_inbound.go:102-198`），能证明"链在产码里、能落盘"，证明不了"页面点得出"；把它当 AC#3 的"真机"凭据＝本仓 memory 第 8 条那个形状。反过来把它当"零证据"也不对——它是 `ParseComposerRequest` 的产码听众。**这条边界归裁决者，不该由普查替验收下结论** | 对抗验收腿（`SPEC-10 §8` 那七项） |
| **J7** | **248-r1 那批改动今天门禁是绿是红** | 它的证据件 §4 那一节是占位（`docs/evidence/s1/248-settings-write-path-r1.md:147-149`）⇒ **四数无凭据**，而我不能跑。下一枚动 `internal/panel`／`cmd/wisp` 的腿**开工前**必须先把 W-1/W-2 复跑一遍，否则撞上的红会被归因成"我改坏了"（本仓 10-01 的"先分是谁的红"那条教训） | 编排者（或 248 的收尾腿）先出一发读数 |
| **J8** | **N15 那枚常驻腿钉的原文** | 它的行号与判语我只从提交 `d0cd3a9b` 的说明里读到，**本腿没有逐字复读 `cmd/wisp/resident_task_source_246_windows_test.go`**（该文件是 `198-r1` 的同族地界，且我不在必要路径上）。⇒ §4 里那一行我已标"待验断言"，别当凭据引 | 任一现读腿（只读该文件即可销账，无需跑） |

---

## §7 我推翻编排者题面之处

| # | 题面原话 | 实测 | 凭据 |
|---|---|---|---|
| **R1** | "票 114 那 **8 枚未勾格**"（题面两处，含"100–149 段只有两枚这种票"那段的落点判断） | **9 枚**。尺＝行首可有空白的未勾框，收尾 `wc -l`（本仓 10-02 已定这把尺，`^- [ ]` 那种锚定法会整枚漏掉缩进项）⇒ 未勾 **9**、已勾 **2**、票面共 11 枚 AC | 票面 AC#3/4/5/6/7/8/9/10/11 未勾；AC#1/#2 已勾 |
| **R2** | "另一枚落地腿 `248-r1` 提交了产码 `0d87a681`，**新建了 `cmd/wisp/panel_inbound.go`**、`internal/panel/config_handlers.go`、`cmd/wisp/panel_config_store.go`，并改了 `internal/panel/bridge.go`／`composer.go`／`composer_dispatch.go`／`pump.go`" | 后半句全对；**"新建 `cmd/wisp/panel_inbound.go`"错**：那枚文件由 **`6609e7e1`（33-r3，2026-09-28 17:21:47 +0800）**建，248-r1 只**改**它（`--diff-filter=M` ⇒ 37 增 9 删）。248-r1 真正**新建**的是 6 枚：`cmd/wisp/panel_config_store.go`、`cmd/wisp/panel_config_248_test.go`、`internal/config/settings.go`、`internal/config/settings_248_test.go`、`internal/panel/config_handlers.go`、`internal/panel/config_route_248_test.go` | `git show --numstat --diff-filter=A/--diff-filter=M 0d87a681`；`git log --diff-filter=A --oneline -- cmd/wisp/panel_inbound.go` |
| **R3** | 题面的整道题："**`ParseComposerRequest` 在生产里零调用者**这一句今天还成立吗" | **不成立，而且失效点不是 248-r1，是 09-28 的 33-r3；把它接上真窗口的是 10-01 的 33-r5（`13acad46`）。** 唯一产码调用者＝`composer_dispatch.go:155`；其上游 `ComposerDispatch.Handle` 有**两枚**产码听众（`panel_inbound.go:163` 由 `main.go:115-120` 派发；`panel_host_windows.go:551` 由 `:319` 的 `w.Bind` 闭包驱动，宿主由 `resident_windows.go:142` 在**无参数启动的常驻进程**里装配）。⇒ 题面那句"**判错方向＝一整枚腿白跑**"这次真的发生了：**方向不是"新建接线"，也不是"248 顺手做掉了大半"，是"补差集"** | §1.1／§1.2／§5-乙 |
| **R4** | "具名分清 `panel_inbound.go` 那条新通路**接的是哪一类入向**（`config.get`/`config.set` 还是**档位／附件／工作区那三族**）" | 三分法**枚数不够**：那枚链**四族半都管**——`Mode`（`:272`）与 `Config`（`:279`）装着，`Workspace`/`Attachment`/`Message` 三枚 nil（`:275-277`）。⚠ **`Message` 这一族题面没列**，它同样在册、同样 nil、同样被 N1 那一族钉挡着（归属＝票 35） | `cmd/wisp/panel_inbound.go:271-281`；`internal/panel/composer_dispatch.go:176-209` |
| **R5** | 题面："票头明写接上的那一票**必须把门放在原生侧**（票 92 的 `R-92-2`：真边界在原生侧）"——把它当作 114 落地的**待建前置** | **门已经立着并且被钉住了**（AC#2 已勾：`composer_handlers.go:133-138` ＋ `composer_handlers_test.go:120`）。今天真正欠的是**门后面那枚 nil 的来源**：两枚 `Confirm: nil`（`panel_inbound.go:242`、`:252`）⇒ 面板改档**只能改严**。这既不是"门没立"也不是"接好了" ⇒ §6-J1 | §2 AC#3／§3.1／§5-乙 |
| **R6** | 题面："AC#6 那一格的处置就是'永远不由本编队勾'，你只需把它与其余各格分开列" | **复认**，不推翻。补一句射程：AC#5 里"`frontend/` 出现 `approval.decide` 即 ban#6 红"与 AC#9 的"渲染器加载集合"**也落在同一块委托地上**——本腿只判了它们的 Go 半支（§2 AC#5、§6-J3），⛔ 没有读 `frontend/**` 一字 | 票面 `:3-7`；§2 AC#5/AC#9 |
| **R7** | 题面 §0 提示："`git status --porcelain internal/panel cmd/wisp` 读数（⚠ 此刻可能非空——那是 `198-r1` 在飞的写面）" | **起手那一发是空串（零行）**，取于 `3b23613b`、`2026-10-02 09:14:18 +0800`；**交件前复量＝4 行**（`M cmd/wisp/run.go` ＋ 三枚 `cmd/wisp/firstrun*` 新件）⇒ 题面那句"可能非空"**延迟兑现了**。⇒ 修正说法：不是"起手可能非空"，而是"**这枚腿跑的过程中会变脏**"，且变脏的文件正是本腿引行号最多的那一枚（`run.go`）。本腿全程未动这两枚目录、未提交别人的写面 | §0 的"并发写面自报"段 |
| **R8** | 题面："**别只数显式调用者**——要读'谁产出／谁投递／谁落盘'三处再判"（并警告"我这样造出过一句假话"） | **方法复认，结论与题面预设相反**：数显式调用者在这里给出的不是假话而是 **1**（真产码），三处读完给出的也是"通路已存在"。**真正的假话风险在反方向**：仓里三处注释今天仍在说"没有通路"——`internal/panel/composer_handlers.go:36-40`（"Nothing calls this handler yet, because the WebView2 \"event -> ParseComposerRequest\" hop does not exist in this tree"⇒ **今天这枚处理器被 `panel_inbound.go:272` 装着、被 H1/H2 两枚产码听众调到，这句整枚过期**）、`cmd/wisp/run.go:316-321`（前半句 "Nothing calls it yet" 对 `rt.modeWrites` **仍然成立**（它零读者），后半句 "the WebView2 … hop does not exist in this tree" **已过期**）、`cmd/wisp/approval_always.go:165`（"this binary links no WebView2 host"，而 `go.mod:19` 已带 `github.com/jchv/go-webview2`、`panel_host_windows.go:319` 已在产码里用它）。**下一枚腿若照这三句造句，造出来的就是题面警告的那句假话** | §1.1／§1.2／`git log --diff-filter=A` |
| **R9** | 题面："那条三张授权词表（`grantRouteWords`／`grantRoutePrefixes`／`grantRouteSuffixes`）**如果 114 落地要新增名字**，逐枚比" | **114 一枚新名都不需要**（三族入口名早在册，§3.3）。而且新增名字会撞两枚**冻结件**钉：`l2_grant_boundary_test.go:1496-1502`（24 名派生正控，`legit` 必须恰为那四枚）与 `:1333-1337`/`:1487-1491`/`:1504-1506`（三因子尺寸与乘积）。⚠ **反过来题面没提的那条风险更要紧**：248 新增的两枚名**不带 `panel.` 前缀**，恰好掉在 `routeLiteralRe`（`composer_test.go:394`）与 `panelMethodRe`（`git_test.go:394`）两把尺的射程之外 ⇒ **AC#8 要的"两边同改"今天对这两枚名完全无人看守**，AC#8 的洞被**削宽**而不是收窄 | §3.2／§3.3／§4-N4/N5/N8/N9 |
| **R10** | 题面对丙档的定义："接线存在但**门的位置错了**（要具名说'错在哪一形、哪枚仪器今天看不见它'）" | 门的位置**没错**（原生侧、在任何落盘之前）。今天真实的形状是**"守门尺看不见那扇门"**：可达性钉只认 `.postMessage(`（`composer_test.go:381`），生产那扇门是 `w.Bind("wispDispatch")` → `window.external.invoke` 的 RPC（`panel_host_windows.go:82`/`:319` ＋ 模块 `webview.go:450-482`、`:139-168`）。⇒ **判"114 有没有接好"如果靠跑 N6/N7 那两枚钉，会得到两个与事实无关的绿**；另附一条本腿判为**新增缺陷**的形状：旧封套直接 `postMessage` 会被宿主库**静默 resolve(null)、Go 侧零审计**（`webview.go:162-168`） | §5-丙 形一/形二／§2 AC#10 |
| **R11** | 题面给的行号一律当"待验断言"（我据此逐枚复量） | 题面里我给的那枚具体行号——"`internal/agent/approval/ui.go` 那枚 `PanelAPI` 的结构性禁区（作者注释逐字「It has no Allow method…」）"——**在，但票面自己在 `internal/panel/l2_grant_boundary_test.go` 之外还有一处引用它的方式**：`cmd/wisp/approval_reply.go:29` 指 "ui.go:45-48"、`internal/agent/approval/gate.go:632` 指 "ui.go:152-159" ⇒ **两枚指针都漂了**，真身＝`ui.go:164-171`（接口体 `:167-171`）。另复量：`panel_inbound_33_test.go` 的用例名与行号如题面所引；`bridge.go` 常量块＝`:41-68`（不是某些旧件写的 `:41-46`） | §3.1／§1.1 |
