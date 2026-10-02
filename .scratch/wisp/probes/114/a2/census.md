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
| **AC#3** 反向对照：更宽档必须走 L2 原生卡片（300s＝拒）＋审计三档全写不因接线而漏 | **部分满足——两半里成立的是"审计"那一半，"卡片"那一半在生产里不可达** | 已成立：三档全写在**写腿**里（`internal/perm/store.go:172-214`，注释逐字"the audit log … records ALL THREE档"），且有常驻仪器（`internal/perm/store_test.go:164-166` 断言 3 枚 `MODE-SWITCH`）；300s 的缺省与钳位也在（`internal/agent/approval/queue.go:107 DefaultApprovalTimeout`）。**没成立**：面板那条链**没有一枚装配把非 nil 的 `Confirm` 交给它**——`cmd/wisp/panel_inbound.go:242` 与 `:252` 两枚 nil，常驻支复用同一函数体（`panel_resident_windows.go:170`）⇒ 变宽请求在 `composer_handlers.go:133` 就被"确认腿未接入"拒掉，**永远走不到 `Set` 里那张 C18 卡**。`cmd/wisp/run.go:659` 那枚 `rt.modeWrites` 确实拿到了真 `confirm`（`:637` → `rt.confirmModeSwitch`），但 **`rt.modeWrites` 全仓零读者**（`grep -rn "modeWrites"` 非测试命中只有 `run.go:322` 字段声明与 `:659` 赋值两处）⇒ 那一支是死码。差的正是：**把 `newComposerDispatchChain` 的 `Confirm` 接成 `run.go` 那枚同一个 `confirm`，并让常驻腿有卡可举** |
| **AC#4** `R-92-1` 覆盖面：扩展名扩到与 d22scan 文本类一致 **或** 改挂结构 | **部分满足——"或"的第二支（结构形）已成立；名册那一支仍差一枚，且两枚名册之间没有尺在核相等** | 结构支已在：`internal/panel/composer_test.go:502 TestTheRendererHoldsExactlyOneDoorToTheHost` 把判定挂在**可达集**上（(i) 宿主调用点必须全在 `src/lib/panel.ts` 且恰 2 枚、(ii) 不许运行时拼点名 `:513`、(iii) 不许计算式路由 `:521` 的前一支、(iv) 每个 `panel.*` 字面量必须是 Go 真答的），自带种形正控 `:533 TestPlantedRendererDoorShapesGoRed`（一枚 `.js` 拼 `panel.mode.set`、一枚数组 `join(".")`）。门二也确实用上了宽名册：`rendererSourceFile`（`:370-376`）＝`.ts/.tsx/.js/.mjs/.jsx/.css/.html`，而 `scanComposerPermissionWrites`（`:597-627`）走的是**整棵 `frontend`**（正控一发在 `:576` 传的就是 `frontend`）。差的两句：① 票面点名的 **`.json` 仍不在** `rendererSourceFile` 里；② 反向也不齐——`tools/d22scan/main.go:1134-1136` 的 `isTextFile` 含 `.json` 却**不含 `.mjs`**⇒ 两枚名册谁也不等于谁，而**仓里没有一枚尺比较这两枚名册**（我逐枚找过：`grep -rn "isTextFile" tools/d22scan` 只有它自己那两枚文件用）。"不许为通过而删判定"与"验收方自己种的越界形状必须红"两句：正控在（`:533`），未删窄化（`:308-315` 那处刻意不扩的是**另一枚** git-switch 判据，注释具名解释为什么，⛔ 别把它读成 AC#4 没做） |
| **AC#5** 不越界：面板侧只许显示＋发起请求；工作区/档位两个输入口不得变成授权口 | **部分满足（Go 半支已满足；`frontend` 半支不由本腿判读）** | Go 半支凭据：`internal/agent/approval/ui.go:167-171` 的 `PanelAPI` 只有 `Reject/Head/View`，**没有 `Allow`**（`:164-166` 逐字"It has no Allow method…"）；composer 六门里没有任何一门能落授权——`panel.mode.request` 只能改档位且变宽被 §3 的门拒；`config.set` 明写拒 `risk.`/`fs.` 两族并点名各自欠的那张 L2 卡（`internal/panel/config_handlers.go:108-136`）；工作区/附件两族连处理器都没有。`frontend` 半支（出现 `approval.decide` 即 ban#6 红、零 emoji ban#8）＝**仪器射程**，判它绿要跑 `sh scripts/d22scan.sh`（§4-待验证窗口 W-3），且 ban#8 的射程按票 141 定案是"注释豁免、字符串不豁免"的那六段（AGENTS.md §1.2），**与本票面写的 `U+2190–U+2BFF` 起射程不是同一把** |
| **AC#6** 真机差分截屏 | **不由本编队做（界面委托）** | 票头 2026-09-23 09:22 横幅（票面 `:3-7`）＋台账 `A102`：`skipped=frontend(owner-delegated)`，保持未勾，本编队⛔ 不替它勾、⛔ 不"只改一角"。本腿对 `frontend/**` 零读取，未据任何页面内容下判 |
| **AC#7** 门禁：d22scan 纯净快照 rc=0 且不降 / `gofumpt -l` 真跑贴原始输出 / POSIX 三合包按 `-v` 口径 | **完全没做（就本轮而言），且没有任何一件现存读数可借** | 票面要求"必须真跑并贴原始输出"。248-r1 的证据件 §4 **没有写**（`docs/evidence/s1/248-settings-write-path-r1.md:147-149` 那一节仍是占位）⇒ 连 248 那批改动的门禁色都没落档，114 更无读数控。本腿按硬边界不跑任何东西 ⇒ 见 §4 的待验证窗口 W-1…W-5（命令给全、期望读数写明） |
| **AC#8** 把"两边同改"从承诺变成门：往 `knownComposerMethod` 加一条必须同时让 `frontend/` 出现对应字面量 | **部分满足——一个方向有门，另一个方向只对 `panel.` 前缀族成立，而 248 新增的两枚名恰好掉在这两把尺之外** | 已有：`internal/panel/bridge_test.go:131-145 TestFrontendComposerRequestsMatchTheEnvelope` **逐枚**要求四枚 `Method*` 常量在 `frontend/src/lib/panel.ts` 里以字面量出现（`:142` 那句 `"the Go side parses a method nobody sends"` 就是 AC#8 要的牙齿）；反方向在 `composer_test.go:521`（页面点名而 Go 不答 ⇒ 红）。**差三句**：① `bridge_test.go:138-141` 那份清单是**手抄的四枚**，新增一枚常量而忘了抄进清单 ⇒ 没人逼它出现在页面里；② `routeLiteralRe`（`composer_test.go:394`）只抓 `"panel.…"` 形 ⇒ **`config.get`/`config.set` 这两枚不带 `panel.` 前缀的名，两把尺都看不见**（248-r1 自己在 `docs/evidence/s1/248-settings-write-path-r1.md` §5-12 承认没登记进 `composerRouteLiterals`）；③ 真正逼住"别扩名册"的是另一枚钉——`internal/panel/git_test.go:381-388` 与 `:517` 用 `panelMethodRe`（只匹配 `"panel.*"`）断言 `bridge.go` 里恰有 4 枚 ⇒ **加第五枚 `panel.*` 会红在那句"this ticket may not extend"，加 `config.*` 不会**。⇒ AC#8 的洞今天**变宽了**而不是变小 |
| **AC#9** 结构钉的扫描根比 ban#6 与"渲染器实际加载的文件集合"都窄 ⇒ 把可达性钉挂在那个集合上 | **完全没做** | 那枚可达性钉的根写死在 `internal/panel/composer_test.go:504`：`filepath.Join(root, "frontend", "src")`，且 `:434` 跳过 `dist`；门二走的是整棵 `frontend`（`:576`）⇒ **"钉的根"与"ban#6 的根"本来就是两棵树**，票面 `frontend/public/*.js` + `index.html` 引用那一形今天仍然只被门二读到、不被可达性钉读到。"渲染器到底加载哪些文件"这一问在本仓**没有答案件**，而要答它得读 `frontend/**` ⇒ 归界面侧/编排者（§6-J3） |
| **AC#10** `panel.ts` 自陈宿主有两种装法（`AddHostObjectToScript` / `postMessage`），钉只认后者 ⇒ 第一种要么显式禁用要么进同一枚钉 | **完全没做，且现实比票面写的更糟：今天生产用的是票面没列的第三种装法** | 现量：`grep -rn "AddHostObjectToScript" --include=*.go internal/ cmd/ tools/` ⇒ **零命中**（模块里只有 vtbl 声明 `pkg/edge/corewebview2.go:123`，无任何调用点）⇒ 第一种装法**没在用**。真正在用的是：`cmd/wisp/panel_host_windows.go:319 w.Bind("wispDispatch", …)`＝go-webview2 的 `Bind`，它注入的存根走 **`window.external.invoke`**（模块 `webview.go:450-482`），宿主按 `{id,method,params}` 解 RPC（`webview.go:139-168`）。而可达性钉的宿主调用点正则（`composer_test.go:381 hostBridgeCallRe`）**只认 `.postMessage(`**⇒ **今天这扇真门对那枚钉是隐形的**。⚠ 还有一条本腿判为**新增缺陷**：页面若沿用旧形直接 `postMessage(那枚封套)`，库会把 `method` 读成 `"panel.mode.request"`、查 `bindings` 查不到、`callbinding` 返回 `nil,nil`（`webview.go:166-168`）⇒ **promise 静默 resolve(null)，Go 侧一行审计都不写**——正是"被静默丢掉"那一形 |
| **AC#11** 把 `npm run render:composer && git diff --exit-code` 做成 CI 一步 | **部分满足——前半已在 CI，后半（防腐坏那半）整枚没做** | 前半：`.github/workflows/ci.yml:718-721` 确有 "Composer states paint the real envelope (render evidence)" 一步，`run: npm run render:composer`。后半：**`grep -n "exit-code" .github/workflows/ci.yml` ⇒ 零命中**，那一步后面没有 `git diff --exit-code`，所以"fixture 可以手写/可以腐坏"这一格今天**照旧开着**。票面 `:76` 自己要求"动 `ci.yml` 之前先看票 111/85a 是否已让出该文件"⇒ 落地腿开工前先查让渡（§6-J5） |

**一句话答"有几枚被 248-r1 顺手做掉了"**：**一枚都没有整格做掉。** AC#4/AC#8/AC#10 三格今天**被 248-r1 动了但方向不利**——它把两枚不带 `panel.` 前缀的名放进了名册，而这三格的全部仪器都以 `panel.` 前缀为锚，于是覆盖面被削薄（AC#8 尤其：它要的正是"名册扩了、页面必须跟着"，而扩了的这两枚恰好没人能逼）。AC#3 是**唯一一处 248 的工作让 114 变近**的：`Confirm` 那个槽位与"回执要能到页面"的形状（`composer_dispatch.go:110-112` 带 raw 的签名、`Handle` 返回串）已经在，114 只差把 `run.go:637` 那枚真 `confirm` 接进 `panel_inbound.go:242/252`。

---

## §3 "门必须放原生侧"这条前置条件现在在哪

（待回填）

---

## §4 撞钉预检（名册，不是颜色）

（待回填）

---

## §5 落点建议（只给形状）

（待回填：甲／乙／丙三档）

---

## §6 判不动的地方

（待回填：J1…）

---

## §7 我推翻编排者题面之处

（待回填）
