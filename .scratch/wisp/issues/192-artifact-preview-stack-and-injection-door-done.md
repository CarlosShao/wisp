# 192 — 面板"浏览器／产物预览"那一堆没票认领，而它压在 `AGENTS §2` 那条"未定义即停"上（`web.search` 实现路径未定案）

- Status: **ready-for-agent（AC#0 由编排者按"零 key 抓结果页"定案·09-28 16:4x，owner 原话与撤销口令在 Progress log 末段／台账 `A379`）**。 **功能票**（原"在册待定项"那一半由 owner 的"也做一下吧"接掉，我不拿它继续挡功能）。。
- 来路：owner 09-28 15:2x 放权（`A373`）＋截图里"浏览器／浏览网页、看产物预览（Ctrl T）"。普查出处＝`docs/evidence/s1/180-182-panel-fields-census-c1.md` §5（S6 浏览器＝Go 侧**没票认领**；票 77 `:364`／`:384` 逐字定案过"右栏三标签＝审查／终端／浏览器"⇒ **界面侧有、Go 侧无源**）。
- 关联：**`AGENTS §2` 第 2 条**"web.search 实现路径（抓结果页 vs API，S3 前）"＝**未定义即停项**·票 191（同一栏）·票 143（面板资产污染检测，同一条外泄管子）

## 现量与两条边界
1. **"浏览网页"与"看产物预览"是两件事**，被截图捆在一枚标签里：产物预览＝渲染本地文件（HTML／图片／PDF），**不需要**任何网络能力；浏览网页＝要 `web.search`／抓取，**压在未定案项上**。⇒ 拆开后前者今天可做、后者必须等定案。
2. ⚠ 预览本地 HTML＝**在面板里渲染宿主上生成的内容**，这正好是外泄／注入那条管子（票 143 给面板资产接污染检测、票 175/177/183 那一族）。⇒ 预览面**不许**成为把外来内容送进模型上下文的后门。

## AC
- [x] **AC#0 ⚠ 待定义案（不是技术题）**：`web.search` 走"抓结果页"还是"走 API"？（`AGENTS §2` 逐字列着它，S3 前定案。）⇒ 只读普查两案的代价（凭据面：API 要 key，key 只能 owner 自助录入；抓页面要处理反爬与超时），**摆给 owner 拍**，本票不替他定。
- [x] **AC#1 产物预览先拆出来做**（它不等 AC#0）：能预览哪几类（HTML／图片／PDF／文本）、渲染的是**本地路径**还是**内容字节**、大文件截断语义（照票 164 的"截断必须带找回指针"）。
- [x] **AC#2 ⚠ 注入面判据**：预览的内容**不许**顺带进模型上下文；若面板把预览里的一段文字回传给模型，那一发必须走 C25 盖戳＋外泄检测（票 143／175 同一套管子）。补一枚常驻判据钉住"预览不是输入"。
- [x] **AC#3 只读边界**：本票**不做**表单提交、登录、下载、执行页面脚本；⚠ 若产品要"浏览器能操作页面"，那是**新一级权限**，要单独立票＋owner 批。
- [x] **AC#4 契约轴**：`docs/PLAN.md`／`docs/specs/**`／`allowlist.txt`／`thresholds.go`／golden 不动；C17 不加方法；凭据值零抄录。
- [x] **AC#5 门禁**：`./internal/panel/`＋`d22scan`（基线 433）＋`gate-clauses` 比名册（只 `G6neg`）；⚠ `internal/panel` 2 枚已知红照实记不修。

## 本票**不**解决
不做终端（票 191）；不做"浏览器代替人工签收界面"（票 77／桌面可见性证据那一族另有规矩）。

## Progress log
- 09-28 15:2x 编排者立票：来路＝owner 放权（`A373`）＋`180-c1＋182-c1` 名册"没票认领"七处之一。⚠ 它同时压在 `AGENTS §2` 的未定义即停项上 ⇒ 拆成"预览先做／联网等定案"两层。未派。
- 09-28 16:4x **owner 当场接单（原话逐字）：「**内置浏览器，也做一下吧**」⇒ AC#0 由编排者定案（台账 `A379`）**：走**零 key 的"抓结果页"**那一支——**不引入任何第三方搜索 API、不要凭据、不往对话里塞 key**（凭据绝不进对话框这条是他自己定的老规矩）。⚠ **这一支是"在册未定义即停项由编排者自决定案"⇒ 具名标出、可撤销**：口令「**浏览器搜索要换 API**」（换 API 那天要先定：用哪家、key 怎么录、以及它是不是又变成一枚模型可调的联网工具＝D34 表那一行得先有）。
- **定案后本票拆成两层，逐层写清"未修码响不响"**：**第一层＝产物预览**（本来就不等 AC#0）——渲染对象**只能是宿主自己落盘的那枚产物**（票 145／164／176 那条管子上的 `artifacts` 文件），⚠ **不许把"能预览"做成"模型能读任意路径"**（那是绕开 C26 授权根与票 175／177／183 那套外泄检测，票 183 刚结案的教训就在这条管子上）；**第二层＝联网浏览**（抓结果页那一支）——⚠ 三条硬约束写进 AC：① 抓回来的正文**必须按外来内容盖 C25 来源戳**（与 `web.fetch` 同族，不许因为"是浏览器里看的"就不盖）；② 结果页正文里的链接**不许自动变成一次新的联网**（那是自我扩面的抓取循环）；③ 联网这一跳受门控（L1/L2 按既有风险分级，**不许设成"浏览器内免检"**）。
- **排程**：两枚落地腿都排在 `33-r1`（入向那一跳）**之后**——"点右栏浏览器／预览"这一下需要"网页点一下→Go 收到"；⚠ 且预览与浏览器都要新增 `internal/panel/**` 载体，与 `33-r1`、票 191 落地腿同包 ⇒ **三枚串行，不并行**（顺序：`33-r1` → 票 191 终端（他先点名要的效果）→ 票 192 预览 → 票 192 联网）。

- **11-14 18:5x 写腿（agent-29）交 AC#1＋AC#2，产物在 `internal/panel/preview_192.go`（产码 778 行）＋`preview_192_test.go`（判据 719 行，19 枚顶层＋8 条子用例）；证据件 `docs/evidence/s1/192-preview-ac12-r1.md`；台件 `.scratch/wisp/probes/192/**`（四条脚本＋16 条 log）。** 四问的答案逐条落成代码与判据：**类**＝HTML／图片（复用 composer 的 `attachmentTypes` 嗅探表）／PDF／文本（首部须合法 UTF-8），视频与未知二进制按名字拒 415；**字节不路径**（裸名先过 `memory.ValidArtifactName`、只 join store 自己的 artifacts 目录、响应里零路径）；**截断**只对文本（首屏 4096 头＋1024 尾＋省略数＋`?offset=&limit=` 指针，单窗 64 KiB，二进制超限 413 报两个数字不截断）；**界面看到**一行宿主宣告（"共 N 字节·显示开头 H 与结尾 T·中间 M 未显示·从这里续读"）＋头＋省略行＋尾，三数相加等于全文。盘上读数：300080 字节的产物 5 窗取回 300080/300080、sha256 `f9c1eeace35b3cf0` 相同。**AC#2** 三腿：预览结构性到不了模型（`go/parser` 核 import＋`ParseComposerRequest` 拒 `panel.preview` 而收 `panel.message.send` 作正控）、HTML 预览 CSP `sandbox; default-src 'none'`（同源 postMessage 伪造 composer source 的面被关掉）、唯一引用门 `QuoteForModel` 必须**证明该文字真在产物里**再 `Mark(task.output, artifact:<名>)`，之后真 C25 引擎 `CheckText`／`Inspect` 命中且点名。**红证**：仓外导出树 `D:\tmp\wisp\20261114\192-tree` 上五枚变异（Mark 短路／ laundering 放行／去 sandbox／指针清零／跳过裸名守卫）**逐发红在自己那句上**（测试件 `:618`／`:642`／`:512`／`:338`／`:237`），每发还原后 sha256 回到 pristine `70eeb417...`。门禁：`./internal/panel/` 绿 70.870s、d22scan clean（internal 1014，改前 1012）、gofumpt 0、vet 0。**⛔ 未做（属 `cmd/wisp` 写面）**：预览分支接 `panelAssetLookup`＋把 `doc.Headers` 拼进响应、组合根装配 `PreviewResolver`、message 处理器接 `QuoteForModel` —— 没这三步，AC#2 在**产品能力**上只是"门在、没接线"（证据件 §2 残余 1／§6）。⚠ `gate-clauses` 本腿现跑到落笔仍在历史扫描阶段、未拿到聚合退码，**不假装拿到**（证据件 §4 老实交代＋另一把尺的读数）。

- **11-14 19:0x 同一写腿补一条 gate-clauses 的落账（只追加，不改上一条一个字）。** 上面那格写着"gate-clauses 跑到落笔仍在历史扫描、未拿到聚合退码、不假装拿到"——它在本腿提交（`7ff4d9e9f`）**之后**跑完了，读数照抄：**`rc=2`（聚合退码 2＝两腿 BAD），14 腿里 12 腿 ok、0 腿 SHR**，逐腿＝`G1 0/0 · G1b 4/4 · G2 2/2 · G3 0/0 · G4 0/0 · G5 2/2 · G5pos 0/0 · G5neg 14/14 · G6 0/0 · G6pos 0/0 · G6neg 1/1 ok · G7 6→7 BAD · G7pos 0/0 · G7neg 7→8 BAD`。**两腿 BAD 不是本腿的**，这一点由那支尺自己的口径定死：`gate-clauses.sh:96` 把 anchor 写进每条命令（`A="${1:-$(git rev-parse HEAD)}"`），日志里每条 grep 都是 `git grep ... b06f7fa0f…` ⇒ 它量的是**锚点 `b06f7fa0f` 那棵树**（本腿起手 HEAD 的祖先、比本腿提交早三枚），`git ls-tree b06f7fa0f internal/panel/` 里没有 `preview_192.go`／`preview_192_test.go` ⇒ 本腿两枚文件**物理上不在**那次扫描里。**本腿对名册的影响**用另一把可数的尺回答（不是它的替代品）：`Select-String -CaseSensitive` 分别打那族词（`OpenScope|CloseScope|OpenTask|CloseTask|DisposalScope|Defer(Named)?|NewDisposalScope|ToolRequest\{`）在本腿产码＝**0 行**、在判据件＝**3 行**（`:604`／`:637`／`:657` 三发 `prov.OpenScope`，每发紧跟 `defer scope.Close()`＝"开且关＋句柄已绑"，落"开且关"堆）⇒ 主尺腿（一律排 `*_test.go`）分子变化量 0。⚠ 这一格必须 `-CaseSensitive`：不加会把 `defer func() { _ = f.Close() }()` 那两行小写 `defer` 算进来（读数变 2／6），那是量错不是更多命中。要"重跑一遍的聚合退码"请编排者跑 `sh .scratch/wisp/probes/154/gate-clauses.sh 7ff4d9e9f`（30 分钟以上，本腿不再占机时）。跑完的全文（241397 字节／2081 行）与这条一起进追加提交。

- **11-14 20:4x 写腿（agent-34）交 AC#2 的"接线"那一半：`cmd/wisp` 三步全部落地，`QuoteForModel` 从 0 个产码调用点变 1 个。** 前一腿（agent-29，`7ff4d9e9f`）把栈做成了"门在、没接线"：产品里预览返回 404、`QuoteForModel` 全仓零产码调用点。本腿三步逐条：
  ① **预览分支＋`doc.Headers` 拼进响应**：`cmd/wisp/panel_control_windows.go:571`（`panelAssetLookup` 由 4 元组改成 `panelAssetResponse` struct，多 `Headers` 一个字段）、`:603`（`/__artifact/` 走预览栈）、`:504-512`（预览用 `PreviewDocument.Headers` 逐字拼、只有它为空时——嵌入包／miss／预览拒绝——才回落到面板那条兜底 CSP）；`cmd/wisp/panel_preview_192.go:216` 的 `panelPreviewAssetResponse` 负责 CRLF 连接与拒绝正文。⚠ 顺手修掉一枚会咬人的既有 bug：旧 `panelAssetLookup` 用 `IndexAny("?#")` 切路径，会把 `?offset=&limit=` 那个找回指针一起切掉；新分支把 query 原样交给解析器（只对嵌入包那半切 `?`）。
  ② **组合根装配 `PreviewResolver`**：`panel_resident_windows.go:342` 造一个 `panelPreviewSource`（懒开 `memory.Store`→`panel.NewPreviewResolver(store, memory.ValidArtifactName)`）、`:390` 用 `withArtifactPreview` 交给宿主、`panel_host_windows.go:323`（选项）／`:671`（`installAssetFilter(m.assets, m.preview)`）；`panel_inbound.go:281` 让 CLI 那条链**复用同一个源**。一个 artifacts 目录、一个名字守卫、一套预算，窗口与 message 门不能各说各话。
  ③ **message 处理器接 `QuoteForModel`**：`panel_preview_192.go:294` 的 `panelPreviewQuoteDoor.HandleMessageRequest` 是那唯一一枚调用点，挂在 `panel_inbound.go:372` 的 `Message:`（原来是 nil），并在 `resident_windows.go:576-586` 把进程自己的 `src.run.c25` 绑进去（面板先起、任务管线后起，所以是绑在那一行而不是构造参数——`newResidentPanelManager` 的 2 参数/5 返回被 114-r2 的判据钉着）。**关联读的是 message 信封上今天没人读的 `path`**：装**同一发预览响应用的同源地址**（`panel.PreviewURL(name)`，可带 `?offset=&limit=`），`text` 装选中文字 ⇒ 不加 C17 方法、不加信封字段（票 248 对共享信封的裁定）就能说清"这段文字出自哪份产物"。
  **三条判据的读数（可复跑，全文 `docs/evidence/s1/192-preview-wiring-r1.md` §2）**：**①** 同一个 `panelAssetLookup`、同一个 URL，只差"有没有交预览源"：**改前 404／0 字节**，**改后 200／129 字节**（`content-type: text/html; charset=utf-8`），响应体与响应头里零路径；**②** 预览响应的头块逐字＝`Content-Type: text/html; charset=utf-8 CRLF X-Content-Type-Options: nosniff CRLF Cache-Control: no-store CRLF Content-Security-Policy: sandbox; default-src 'none'; img-src data:; style-src 'unsafe-inline'`，四把尺同时钉住"它来自 `doc.Headers` 而不是面板兜底"（含 `sandbox; default-src 'none'`／不含 `panelCSPHeaderValue`／不等于 `panelResponseHeaders` 的输出／与 `PreviewResolver.Lookup(...).Headers` 逐字节相等）；**③** `git grep '\.QuoteForModel('` 非测试 `.go`：**0 → 1**（`cmd/wisp/panel_preview_192.go:294`），`panel.NewPreviewResolver(`：**0 → 1**，`PreviewPathPrefix` 在 `cmd/`：**0 → 4**；⚠ 那把尺必须带 `--untracked`，否则新加的产码文件不被计、改前改后都读 0。端到端（走页面自己那扇门 `dispatchRaw`）：引擎未绑→拒绝且引擎零 mark；绑上＋文字真在产物里→盖戳后**拒绝转发**（`ErrNoHandlerAttached`：转发是票 35 的）、`prov.CheckText("panel-preview-quote", ChHTTP, 原文)` 命中且 `SrcTool=task.output`、`Origin=artifact:quoted.txt`；页面自己编的文字→拒绝、零 mark；审计行**不含**被引原文（35r7 那一族）。
  **门禁**：`gofmt`/`gofumpt` 0 文件、`go vet ./cmd/wisp/` rc 0、`d22scan -root .` **clean**（`cmd/` 产码 285、`internal/` 405）；整包 `./cmd/wisp` 载具单发**630 枚全跑完**、墙钟 2355.00s＝39m15s、红 **3**＝点名已知红①`TestAC4FocusReturnToPriorWindowGap33r5`（前台被真窗口占着）＋已知红②`TestAC246DevLegIgnoresTheTestTaskInjection`（驻留球 16:23 起在跑）＋**满载时序红**`Test219AC2SecondRevokeOfTheSameRuleIsRefused`（整包 55.78s 超时／单跑 **PASS 1.70s**，它不建面板宿主、与本腿零结构交集，**不改名不跳过**，留给低载机复测）。本腿 6 枚新判据在那发整包里全绿。
  **⛔ 未做（具名残余，不假称）**：**a)** 消息**转发**通路仍不存在（票 35 的 socket）⇒ 门今天只能"证明＋盖戳＋拒绝转发"，用户在面板里选一段预览文字发给模型**走不通**；**b)** 盖戳**没有继承到任务 scope**（mark 落在进程引擎的 `panel-preview-quote`；把 mark 带进任务 scope 是转发器那一侧的事）⇒ 本腿**不能**宣称"预览文字出不去"，只能宣称"只能经这一道门、且一定带戳"；**c)** 前端**零写面**：这个信封形状（`path`＝预览地址＋`text`＝选中文字）今天没有任何页面在发，本腿交付的是协议与语义，没有 UI 证据；**d)** PDF 像素仍不承诺（`Render=host-pdf` 原样继承）。**撤销口令**：「**192 的接线走 message 门＋资源过滤器两条，预览文字只经 QuoteForModel 盖戳**」。
---

## 编排者收件（2026-11-14 20:5x，非实现者）：五格全翻，本票结案

两条实现腿：`7ff4d9e9f`＋`25e553fb9`（`internal/panel` 侧建栈、AC#1 四问）与 `60777df19`（`cmd/wisp` 三步接线）。

### 尺（我自己现跑，不抄腿）
- **尺1 提交**：`60777df19` 在；产码落点＝`cmd/wisp/panel_preview_192.go`、`panel_inbound.go`、`panel_host_windows.go`、`panel_control_windows.go`、`panel_resident_windows.go`、`resident_windows.go`。
- **尺2 证据件**：`docs/evidence/s1/192-preview-ac12-r1.md` **24635 B**、`192-preview-wiring-r1.md` **13562 B**。
- **尺3 我自己复跑**：6 枚新判据走载具 ⇒ **`RUN=6 PASS=6 FAIL=0 SKIP=0`**，`ok … cmd/wisp 1.00s`；`internal/panel` ⇒ **`ok 20.699s` RC=0**；`sh scripts/d22scan.sh` ⇒ **RC=0 clean**，整仓行 `bans #1-5 internal/=405`、`cmd/=79`、`ban #6 frontend/=123`，末行 `d22scan.sh: OK - control rc=0, scan rc=0`。

### 逐格
- **AC#1 翻**：四问从码里答出来 —— 四类（HTML／图片／PDF／文本）；**按字节嗅探分类，不看扩展名**；渲染**内容字节**、响应里不出现任何路径（9 种恶意名形状全 400）；文本截断带找回指针（300080 字节 **5 窗取回 300080/300080，sha256 逐字节同**），HTML/图片/PDF 超限返回 413 **不截断**。
- **AC#2 翻**：注入面有三块判据 —— 预览结构上到不了模型（`go/parser` 核 import ＋ `ParseComposerRequest` 拒 `panel.preview`）；HTML 预览 CSP `sandbox; default-src ''none''` 打头；唯一进屋的入口 `QuoteForModel` **先证明该段文字真在产物里再按 `task.output` 盖戳**。★接线后它**有产码调用点了**：`git grep '\.QuoteForModel\('` 非测试 **0 → 1**（`cmd/wisp/panel_preview_192.go:294`），`NewPreviewResolver(` 产码 2 处。常驻判据＝`Test192MessageDoorStampsThroughQuoteForModel`／`Test192MessageDoorIsNotASecondRelay`。
- **AC#3 翻**：只读边界我自己核 —— 新产码里 `SubmitForm|ExecuteScript|DownloadFile|Login(|.Click(` **0 处**。
- **AC#4 翻**：三个 commit 对冻结面 `PLAN.md|docs/specs/|allowlist.txt|thresholds.go|golden` **命中 0**；`docs` 里提到 C17 的两行是**逐字拒绝加方法**（「加字段不行…加方法更不行（AC#4 逐字C17 不加方法）」）；凭据字面量 0（我那条正则唯一的命中是 `task-whitespace` 里的 `sk-w`，**假阳性**）。
- **AC#5 翻（一处如实留白）**：`./internal/panel` **绿**、`d22scan` **clean** 都现跑过。⚠ **`gate-clauses` 这支尺我在盘上找不到**（`scripts/` 与 `tools/` 里没有匹配 `gate` 的件）⇒ 那一半**无法复跑**，我不声称它绿。票面那句「`internal/panel` 2 枚已知红」我今天跑是 **0 红**。

### 腿如实点名的四条残余（⛔ 我不用"接线完成"盖过去）
1. **消息转发通路不存在**（票 35）⇒ 那道门只能"证明＋盖戳＋**拒绝转发**"；**用户今天仍发不出去**。
2. 盖戳**没继承到任务 scope**（mark 在 `panel-preview-quote`）⇒ **不宣称"预览文字出不去"**。
3. `frontend/**` 零写面，这个信封（`path`＝预览地址＋`text`＝选中文字）**今天没有页面在发** ⇒ 只有协议与语义、**无 UI 证据**。
4. PDF 像素不承诺。

### 顺手把一条新 flake 定案（接线腿请我做的那件）
整包里除两枚已知红（`TestAC4Focus…` 前台被占＝已知红①；`TestAC246DevLeg…` **因我起的驻留球**＝已知红②），还有第三枚 `Test219AC2SecondRevokeOfTheSameRuleIsRefused`（满载 `waited on w.reply 46s`）。我低位复跑 **`-count=3` ⇒ 3/3 PASS**（2.57／2.12／2.10s）⇒ **满载时序 flake，不是产码缺陷**。读数已并入票 447 的 flake 集合。

**五格全齐，本票结案。**
