# 190 — 面板"工作区文件"那一堆（文件树）今天宿主侧无源（owner 09-28 放权：缺就补）

- Status: **ready-for-agent（先只读设计核，再落地）**。⚠ **功能票**。
- 来路：owner 09-28 15:2x 原话（`A373`）＋截图里"工作区文件／浏览会话工作区的文件（Ctrl P）"。普查出处＝`docs/evidence/s1/180-182-panel-fields-census-c1.md` §5（S4 文件树＝**宿主侧无源**，且"文件树"三字全票池只命中票 182 一枚 ⇒ **没票认领**）。
- 关联：票 182（上游）· 票 145（快照载体）· 票 102／C26（路径解析与授权根）· 票 77（前端已定案该标签）

## 现量
- 工具侧有 `fs.list`（在册名册级核到，见 `180-c1＋182-c1` §5），但那是**模型可调用的工具**、不是宿主侧读面；面板要的是"当前工作区的树"，两者授权语义不同（工具的走 C26 与门控，面板的只能读宿主已授权的根）。
- ⚠ 全仓**没有**宿主侧"列目录给面板"的读面。

## AC
- [x] **AC#1 只读设计核**（09-28 17:5x 编排者按现量复跑对格后勾：设计核＝`docs/evidence/s1/190-file-tree-design-a1.md`；根＝**配置根∩工作区的那个集合**（可 0 枚、空根是正常态，`SPEC-03:35` 默认 `allowed_dirs[]=[]`），**树腿不接受面板传根**；`allowed_dirs` 是**判级输入不是执行时硬边界**（`internal/risk/rules_gateway.go:45-49` 越界只吐 `level: L2`、签名里没有 deny 出口）＝账 `A390`）：树的根取哪一枚（工作区根？授权根 `allowed_dirs`？两者不等同时以谁为准）、深度与条目上限、忽略规则（`.git`／`node_modules`／构建产物）；⚠ **忽略过滤必须问 git index，问不到就全扫并响亮自陈**（这仓的 `d22scan` 已定过这条规矩，别另发明一套）。
- [ ] **AC#2 只读落地＋常驻判据**：值来自真目录（临时目录合成一棵树，断言条目集合）；根的选取要一发判据钉住"工作区未窄化时报的是授权根而不是整个盘"。
- [ ] **AC#3 ⚠ 只读边界**：树上**不许**出现删除／移动／重命名／新建任何写动作（普查把它列为雷区＝要 owner 逐枚批）。补常驻判据钉"这一维只读"。
- [ ] **AC#4 越权那一形必须显式回答**：目录不可读／权限拒绝／符号链接与 junction 跨界（票 119／126 那一族）三形各给一条现量与文案，**不许静默少列**。
- [ ] **AC#5 契约轴**：D34 表不加行、`docs/PLAN.md`／`docs/specs/**`／`allowlist.txt`／`thresholds.go`／golden 不动、C17 不加方法；路径决策一律过 `risk.PathResolver`（**不许**在它之外用 `filepath.Clean|Abs`，那是 d22scan ban #2）。
- [ ] **AC#6 门禁**：`./internal/panel/ ./internal/tools/`＋`d22scan`（基线 433）＋`gate-clauses` 比名册（只 `G6neg`）；⚠ `internal/panel` 2 枚已知红照实记不修。

## 本票**不**解决
不做树上写动作（待批）；不做搜索／模糊匹配（那是"打开文件"体验，另计）；不画界面。

## Progress log
- 09-28 15:2x 编排者立票：来路＝owner 放权（`A373`）＋`180-c1＋182-c1` 名册"没票认领"七处之一。未派。
- 09-28 17:5x 编排者收 `190-a1`（只读设计核·一枚 commit `9b6a1be5` 只动它自己的证据件）：`AC#1` 已按现量复跑对格并勾（其余五格未勾，勾要非实现者裁）。**Status 追加（上面那行原样保留不删）**：设计核已交，**`AC#2` 落地腿按住等 `145-r2` 交件**——它要写 `internal/panel/composer.go`／`pump.go`／`cmd/wisp/panel_pump.go`，那三枚文件此刻正被 145 写（`A388`）。
- **两处由我处置的**：① 它问"票面 Progress log 要不要我追加"——**答＝归口在我这边**，本条即追加，它零枚票面写入是对的（只读程不碰别人地界）；② 证据件 §0 名册誊写里 `flip-7.txt` 重复一行（文件内已自陈"原始 71 枚、差集按名算"）⇒ **不追改**：那格不是任何判据的分母，改它只会再生一份"改过的读数"，账在 `A390`。
- **本程新增一枚对我有用的现量（我原先不知道）**：`SPEC-03:35` 那行逐字带 `delete_enabled(bool)=false` ⇒ **删除动作在规格里本来就是默认关**，本票"树上不许有写动作"这一条从此有规格凭据、不再只是我们的自限。（另：`internal/tools/paths.go:43-47` 逐字"can only ever TIGHTEN … no workspace choice can widen authority"，与 `A352` 那条"允许根是判级输入"同向。）

## Progress log — 190-a1 第二程（10-08，只读；⛔ 本程零勾选、零改既有行）

- 起手锚 `914177e60aed5d768988e74f180be42ec29857e8`（Thu Oct 8 10:31:36 2026 +0800），分支 `dev`。证据件 `.scratch/wisp/probes/190/a1/100-min-tree-census.md`（＋同目录 `00-anchor.md`）。
- **框数（缩进味尺原文）**：`grep -cE '^[[:space:]]*- \[ \]'` → **5**、`- \[x\]` → **1**；本票**无缩进子项**，故与朴素尺 `grep -c '^- \[ \]'` → **5** 同值（这一处两把尺不构成分母差，别拿本票去校准别的票）。
- **票面引用逐枚复核，全部在位**：`rules_gateway.go:45-49`（越界只吐 `level: L2`、无 deny 出口）✓；`SPEC-03:35` 同一行同时含 `allowed_dirs[]=[]` 与 `delete_enabled(bool)=false`（本票第 26 行那两处引用都成立）✓；`internal/tools/paths.go:43-47` ⚠ 被引文逐字在 **:45-:46**，票面 :43-47 是含住它的**宽区间**——不算过期，但落地腿引用请改写 :45-46。`design/doubao/demo/rb-files.js` 在位（143 行，且不在起手时 `design/**` 那批在飞改动里）。
- **本程增量（不重证"有没有源"，只量"差几跳"）**：
  - 甲形（面板直用 `fs.list`）＝**7 跳**，其中 **3 跳是闸门**：① `frontend/src/lib/panel.ts` 方法字面量必须进封闭词表（`composer_test.go:409 composerRouteLiterals()`，门判据 `:502` 要求全渲染层只有 2 枚 `postMessage` 调用点）；② `bridge.go:41-68`＋`knownComposerMethod:146-152`＋`ParseComposerRequest:132`＝**C17 名册，撞 AC#5"C17 不加方法"⇒ 人工批准**；③ 任务身份——`(*tools.Bridge).Execute(ctx, agent.ToolRequest)` 的生产派发者今天只有 `internal/agent/loop.go:737/:741`（＋装饰器 `(*tools.subagentToolProvider).Execute` 转发 `p.inner.Execute`，`internal/tools/subagent_197.go:578`），`bridge.go:872` 逐字"being such a caller means being the first one"、`:883` 逐字"That choice is Q-56 and this file does not answer it"＝**撞未定案 Q-56**，本票未定案清单里没有它。其余 4 跳：`composer_dispatch.go:175` 派发、panel↔tools **双向零 import**（实测 `grep -rl` 两次 rc=1）、`Execute:268→checkCaps:372→Assess→route:423`（`case risk.L0: :425` / `case risk.L2: :457→PendingApproval:458`）、结果回面板（快照今天只有 `ws=set/unset` 一个布尔，`panel_pump.go:341`）。
  - 乙形（宿主单开只读枚举）＝**4 面**，且 `cmd/wisp/run.go:275` 字段 `paths *tools.PathCanonicalizer`、`:501` 构造点**已持句柄** ⇒ 不需新 import 边、不需新接缝；快照腿已有同形状先例（`panel_pump.go:83→:87 WorkspaceRoot()`、`:229 ReadGitForWorkspace`）。必过的解析器调用点：`paths.go:115 resolve`（注释逐字"the ONE canonicalization call site"）→ `Canonicalize:131` → `paths.go:201 InAllowlist`，硬顺序同 `rules_gateway.go:45-49`；拼子项只许 `fs.go:280 joinForListing` 那一族（其注释 `:276-279` 逐字"must not grow into a general path joiner"）。
  - **乙形唯一贵的不是画树**：AC#1 要求"忽略过滤必须问 git index"，而那枚实现全仓唯一一份、且**困在 `package main`**——`tools/d22scan/gitignore.go:1 package main`（尺：`sed -n '1,10p' … | grep -n package` → `1:package main` rc=0），内部包侧 `grep -rln 'gitignore' internal/` → **无输出 rc=1** ⇒ import 不到；重写一份＝AC#1 逐字禁（"别另发明一套"）。**这一面是跨票工程，不属本票任何一格。**（与 10-08 兄弟腿 `189-a1` commit `ae9cd5fb` 那句"全仓唯一真起外部 git 是 `tools/d22scan/gitignore.go:373`"同向，此处只引用、未重跑其尺。）
  - **`AGENTS.md §1.2` ban #2 的实际窄处（现量）**：仪器只匹配 `filepath.Clean|Abs`（`tools/d22scan/main.go:737-741`），**`Join`/`Dir`/`Rel` 根本不扫**；既有先例＝`internal/panel/git.go` 14 处 `Join`/`Dir` 做真实文件系统决策（`:296 :308 :340 :349 :394 :402 :408 :445 :472 :483 :508 :517 :520 :526`）＋`composer.go:345`（注释 `:338` 讲 ticket 75），**全都不在 `tools/d22scan/allowlist.txt`**（该文件 `pathresolver-bypass` 仅 3 条在册：`:3 memory/open.go`、`:4 models/manifest.go`、`:6 risk/pathresolver.go`）。⇒ 规矩比仪器宽，先例已在宽处生活；**本程不主张迁移，只记账**。
  - **枚举腿上的错误外抛（同类形状，具名，⛔ 未改）**：`fs.list` 自己的 `Execute` 里三处把 err 整串并进正文——`internal/tools/fs.go:208`（"路径无法解析（按 fail-closed 拒绝）："+`err.Error()`，此处的 err 正是可透出的 `risk.ErrReparseDenied`，本体在 `internal/risk/pathresolver.go:36` 是一句英文全串）、`:212`（"打开目录失败："＋`*PathError` ⇒ 绝对路径原样外带）、`:218`；同族 `fs.go:144`（`fs.read` 的 :208 孪生）。对照面：同一函数里 `fs.go:269` 写死 `"? %s stat-error"` **不外抛**，是现成的"该怎么写"。
  - **符号链接的可见性有平台形状**：`reparseComponents` 只有 Windows 有实现（`internal/risk/pathresolver_windows.go:60`），POSIX 那份 `return nil`（`internal/risk/pathresolver_other.go:12`）⇒ 非 Windows 构建下 C26 看不见任何符号链接，跨界只靠后续 `InAllowlist` 兜。**落地腿不许在 POSIX 上声称验过跨界拒绝。**单层"看见"符号链接是安全的（`fs.go:258 os.Lstat` → `:264 kind = "l"`，不调解析器）；危险只在下钻。
  - **超大目录的钱**：`fs.go:216 Readdirnames(-1)` 先全量读、`:220 sort.Strings`、`:221` 才截断（默认 500，`fs.go:56`）⇒ 照抄这一族＝树每个节点都全量读。
  - **雷区（点击语义）**：`design/doubao/demo/rb-files.js` 整只文件面板只有 `:124` 一枚监听器；目录行 `:126-135` 是**纯 DOM 类切换、零宿主调用**，文件行 `:137-140` 只是 `app.toast('打开文件 '+…)`——**"打开文件"在稿里从来不是动作，是文案**（`:108` 甚至硬编码了绝对工作区根）。⇒ 真要让它成为动作＝新方法＋exec 面，那才是要 owner 逐枚批的东西；本票"显示"那一维**没有任何现成通路**走到写侧（写侧五枚全在 `BuiltinFSWriteEntries`，`internal/tools/fs_write.go:772`，`fs.delete` 仅 `:779 if d.DeleteEnabled` 时 append）。
  - **L2 刷屏机理（证据，不改判）**：范围外目录只会被判 L2（`rules_gateway.go:45-49`），而 L2 的唯一应答通道是原生卡（`bridge.go:457-458`），超时自动拒绝（`:464-465`）；面板侧永不是答案（仪器 ban #6 `main.go:23-29`＋`:271`、Go 边界尺 `internal/panel/l2_grant_boundary_test.go:4-8`、`:1533`；命名先例条文 `bridge.go:46-68` 逐字"a route that lets the page change configuration is not a route that lets the page *approve* anything"）。
- **给编排者的判断（由你裁）**：**乙形便宜**；且只要"bounded depth ＋ 展开是纯前端状态"，**乙形不需要任何新面板方法名＝不触 C17＝不必走人工批准**（demo 的点击形状已证明可零宿主调用）。**只有**改做"点了才向宿主要下一层"的懒展开，才必须新增 `panel.*` 方法＝C17 变更＝人工批准。另：本程新撞出一枚**未定案 Q-56**（甲形专用腿才会碰），建议进台账由人拍。
- **AC#6 那三把尺本程一律未跑**（`internal/panel` 测试／`d22scan` 基线 433／`gate-clauses` 名册）＝硬约束"零 go 命令"，⛔ 不许据此声称任何一条已核。另引用兄弟腿 `189-a1`（`ae9cd5fb`）的两把尺不重跑：它报"AC#6 那 2 枚已知红实为 4 枚"、"裸 `gate-clauses.sh` 实物在 `probes/154/`"、"09-28 同名件行号全漂"——**最后那一条对 09-28 的 `docs/evidence/s1/190-file-tree-design-a1.md` 同样是风险提示，本程未逐格复跑那份设计核的行号。**
