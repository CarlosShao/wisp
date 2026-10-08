# 190-a1（第二程）— 画出"最小那一版文件树"要差几跳：甲／乙两形现量

- 起手锚：`914177e60aed5d768988e74f180be42ec29857e8`（Thu Oct 8 10:31:36 2026 +0800），分支 `dev`，见本目录 `00-anchor.md`。
- 本程**零 go 命令**（⛔ build/vet/test/run），**零产码写入**；所有读数＝静态读文件＋`grep`。
- 前一程同名腿（`9b6a1be5`，Sep 28）已交设计核 `docs/evidence/s1/190-file-tree-design-a1.md`（§D①～⑤、AC#4 口径），本程**不重证"有没有源"**，只量"最小那一版差几跳、钱在哪几面、撞不撞人工批准"。差集见末节。

## 记法约定（本项目的账）
枚数／调用者一律写**类型全名＋调用形状**；⛔ 裸符号名。每把尺自落 `rc=N`；证据件只用 `.md`（⛔ `.out` 会被根 `.gitignore` 第 8 行静默吞掉）。

---

## 1. 问① — 工单 `190-*` 整份读完

尺（原文，缩进味那两味）：
```
grep -cE '^[[:space:]]*- \[ \]' .scratch/wisp/issues/190-*.md   → 5   rc=0
grep -cE '^[[:space:]]*- \[x\]' .scratch/wisp/issues/190-*.md   → 1   rc=0
grep -c '^- \[ \]'            .scratch/wisp/issues/190-*.md   → 5   rc=0   （朴素尺，作对照）
```
⇒ **未勾 5 枚（AC#2／AC#3／AC#4／AC#5／AC#6）＋已勾 1 枚（AC#1）**。票 26 行（`wc -l`）。
⚠ 本票**没有缩进子项**，所以"缩进味尺"与"朴素尺"在此同值（5）——这两把尺的口径差**在本票不构成差异**；不要拿本票数去校准别的票。

票面行号／文件名前缀逐枚复核（**全部在位**）：

| 票面引用 | 现量 | 判定 |
|---|---|---|
| `docs/evidence/s1/180-182-panel-fields-census-c1.md` §5 | 存在，349 行 | 在位 |
| `internal/risk/rules_gateway.go:45-49` | :45 `if !ctx.canon.InAllowlist(canonical) {`，:46-49 `return &contribution{rules:[R2], level: L2, ...}` | 在位，逐字对格（越界只吐 L2、签名里没有 deny 出口） |
| `SPEC-03:35` | 该行＝`🔒 [fs]` 那一行，逐字含 `allowed_dirs[]=[]`、`reparse_point_exceptions[]=[]`、`delete_enabled(bool)=false` | 在位（票面两处引用同一行 35，两处都成立） |
| `internal/tools/paths.go:43-47` | 被引文 "can only ever TIGHTEN … no workspace choice can widen authority" 逐字在 **:45-:46**；票面写的 :43-47 是**含住它的宽区间** | 在位，不算过期（但落地腿引用时建议改写 `:45-46`） |
| 票 145／182／102／77／119／126 关联 | 票文件在 `.scratch/wisp/issues/` 池内（本程未逐枚读，只在被需要处点名） | 未逐枚复核＝本程射程外，见末节缺陷 |
| `design/doubao/demo/rb-files.js`（编排者转述里的前缀） | 存在，143 行，且**不在**起手时 `git status` 的改动集里（`design/**` 那批在飞件不含它） | 在位 |

AC#6 那三把尺（`internal/panel` 2 枚已知红、`d22scan` 基线 433、`gate-clauses` 只 `G6neg`）本程**一律未跑**——它们全是 go／脚本执行面，撞上硬约束 1。⛔ 不许我据此声称任何一条已核。

---

## 2. 问② — 甲形：面板宿主直接复用 `fs.list`

### 2.1 它今天长什么样（现量）
- 类型全名＋调用形状：`fsList.Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (tools.Result, error)`，注册在 `internal/tools/fs.go:328`（`{Tool: fsList{d: d}, Decl: FSListDecl()}`，`BuiltinFSEntries` `fs.go:325`）。
- 名册声明：`FSListDecl()` `fs.go:309` ⇒ `Capabilities/Needs = CapFSRead`、`Declared = risk.L0`、`PathParams = ["path"]`、`Resident: true`、`Provider: KindBuiltin`。
- **`fs.list` 只列一层**：`Readdirnames(-1)`（`fs.go:216`）→ `sort.Strings`（`:220`）→ 截到 `listCap()`（`:221`，默认 `defaultMaxListEntries = 500`，`fs.go:56`）。入参 schema 只有 `path`/`limit`（`fs.go:185-188`），没有递归、没有 depth 参数；正文产出形状＝`fmt.Fprintf(&b, "%s (%d 条目%s)\n", canon, …)` 起（`fs.go:231`）到 `return Result{Text: b.String(), …}`（`fs.go:238`）。
  ⇒ **一棵 N 个可展开目录的树＝N 次跨栈调用**，这是甲形的主导代价，不是实现细节。
- 面板侧来源的调用者：**0 枚**。尺（非测试、按调用形状）：
  ```
  grep -rn 'Execute(ctx, agent.ToolRequest\|Execute(ctx, req\|\.Execute(.*ToolRequest' internal/ cmd/ --include='*.go' | grep -v '_test.go'
  ```
  → 命中 2 处（`rc2=0`）：`(*agent.Loop).` 侧 `internal/agent/loop.go:737` 与 `:741`（`l.opt.Tools.Execute(ctx, req)`）；装饰器 `(*tools.subagentToolProvider).Execute` 于 `internal/tools/subagent_197.go:578` 转发 `p.inner.Execute(ctx, req)`。
  ⇒ `bridge.go:872` 那句 "Today no production code dispatches on the bridge except loop.go" **实质仍成立**（装饰器的 req 来自 loop 已 admit 的子代理任务），但**裸 grep 会把它读成第二个独立派发者**——这就是本项目记过的"把构造点／同名当调用者"那一账的形状。

### 2.2 甲形跳数（从面板到 `fs.list`，逐跳 `file:line`）

| 跳 | 要动／要穿的那一面 | 现量锚 | 性质 |
|---|---|---|---|
| 1 | 渲染层多一条宿主方法字面量，且必须进"封闭词表" | `frontend/src/lib/panel.ts:246/255/272/300`（`sendRequest` 四枚现值）、`:179-181`（`panel.approval.request`）；尺 `internal/panel/composer_test.go:409 composerRouteLiterals()`、`:417` | **判据锁**：`TestTheRendererHoldsExactlyOneDoorToTheHost`（`composer_test.go:502`）要求 ① 全部 `.postMessage(` 只活在 `src/lib/panel.ts` 且**恰好 2 处**、④ 每个 `panel.*` 字面量都必须是 Go 侧答得出的 |
| 2 | Go 侧方法名册＋白名单 | `internal/panel/bridge.go:41-68`（六枚 `Method*` 常量）、`:146-152 knownComposerMethod`、`ParseComposerRequest:132` 拒绝未知方法 | **＝C17 契约面**（票 AC#5 逐字"C17 不加方法"⇒ **人工批准**） |
| 3 | 派发分支 | `internal/panel/composer_dispatch.go:175`（`dispatch` 的 switch，`:177/182/187/192/197/202` 六条现有分支）＋ `rosterMismatch:215` | 名册一致性尺会跟着要改 |
| 4 | 跨包够不着：面板包**不**import 工具包 | 尺（⛔ 不带管道，`grep -rl` 自己出 rc）：`grep -rl 'wisp/internal/tools' internal/panel/` → **无输出，rc=1**；`grep -rl 'wisp/internal/panel' internal/tools/` → **无输出，rc=1** | 双向零 import 边＝**架构性**的；注释里也写了这条边不存在（`internal/panel/workspace.go:41-44`："the only one internal/panel and internal/tools may both hold without importing each other"）⇒ 甲形必须在装配根新插一条接缝接口 |
| 5 | 任务身份与 taint 作用域 | `(*tools.Bridge).Execute` `internal/tools/bridge.go:268` 吃 `agent.ToolRequest{TaskID, CallID, CorrelationID}`；`OpenTask` `:843`／`CloseTask` `:893`／`assessorFor` `:920` | **撞未定案 Q-56**：`bridge.go:872`（"being such a caller means being the first one"）、`:883`（"That choice is Q-56 and this file does not answer it"）＝`AGENTS.md §2`／台账 `Q##` 那一类**待人拍板**，票 190 的未定案清单里没有它 |
| 6 | 每一次点击穿整套风险栈 | `Execute:268` → `checkCaps:372` → `assessorFor(...).Assess`（R1..R9）→ `route:423`；`case risk.L0: :425` 直放，`case risk.L2: :457` → `b.gate.PendingApproval:458` | 树是**只读、L0**的，但一旦目录不在授权根内就落 L2（见 §4 雷区）⇒ 队列被刷屏 |
| 7 | 结果回面板 | 工具产出是 `tools.Result.Text`（模型正文形状，`fs.go:231-238`）＋`Decision` 卡片（`cmd/wisp/panel_pump.go:256 panelArgs(d tools.Decision)`）；快照里今天没有"文件树"这一维（`panel_pump.go:83 workspaceView()`、`:341 ws=` 只有 set/unset 一个布尔） | 要么把模型正文当界面数据用，要么再加一层序列化 |

**甲形跳数＝7**，其中 **3 跳是契约／未定案闸门**（跳 2＝C17 名册、跳 5＝Q-56、跳 1＝渲染层唯一门判据），不是普通工时。

---

## 3. 问② — 乙形：给面板宿主单开一条只读枚举

### 3.1 Go 侧要动的文件（现量，⛔ 猜）
1. `internal/tools/paths*.go`——**新只读方法挂在这一枚上**，因为 C26 的唯一规范化调用点就在此：`(*PathCanonicalizer).resolve` `paths.go:115`（注释逐字："resolve is the ONE canonicalization call site in this package"）、`Canonicalize` `:131`、`InAllowlist` `:201`、`Roots` `:246`、`UnusableRoots` `:254`、`RewrittenRoots` `:262`；工作区那一半：`WorkspaceRoot` `paths_workspace.go:67`、`ResolveWorkspace` `:82`、`SetWorkspaceRoot` `:120`、`ClearWorkspace` `:142`、`inRoots` `:152`。构造点唯一：`NewPathCanonicalizer(allowedDirs, reparseExceptions []string) *PathCanonicalizer` `paths.go:65`。
2. `cmd/wisp/run.go`——装配根**已经持有**这一枚句柄：字段 `paths *tools.PathCanonicalizer`（`run.go:275`），构造 `rt.paths = tools.NewPathCanonicalizer(allowed, cfg.FS.ReparsePointExceptions)`（`run.go:501`）。⇒ 乙形**不需要新 import 边、不需要新句柄**。
3. `internal/panel/`（视图那一维）——现成可抄的形状就在 `workspace.go`：`PathScope` `:45-55`（三枚方法：`WorkspaceRoot() risk.Result` / `ResolveWorkspace(string) (risk.Result, error)` / `SetWorkspaceRoot(string, risk.Result) error`）、`WorkspaceViewFromRoot(res risk.Result) WorkspaceView` `:85`、`UnsetWorkspaceView()` `:65`、`RequestWorkspaceSwitch` `:111`。**⚠ 视图必须带 `risk.Result` 而非裸串**：`:38-44` 逐字解释了裸串只能回答"nothing was rewritten"（票 102／181 AC#7 那一账）。
4. `cmd/wisp/panel_pump.go`——快照腿**已经在替面板读盘外世界**：`workspaceView()` `:83-87`（`rt.paths.WorkspaceRoot()`）、`panel.ReadGitForWorkspace(rt.workspaceView())` `:229`。⇒ 树这一维**可以 rides on 现成快照**，不必新开入站方法（见 §5 的判定）。
5. 忽略规则那一面（**乙形唯一贵的地方**，见 §3.3）。

### 3.2 走乙**必须**经过的解析器调用点（＝"绕没绕开 C26"的清单）
- 根：`WorkspaceRoot()`（收窄态）或 `Roots()`（配置态）→ 取 `risk.Result.Canonical`；**⛔ 不许**自己 `filepath.Abs(工作目录)`。
- 每个待列目录：`Canonicalize(raw)`（内部即 `risk.Resolve(input, exceptions)` `internal/risk/pathresolver.go:105`）→ `Result.Actable()` `:93` → `InAllowlist(canonical)` `:201`。**这两个就是硬顺序**：先规范化再判范围，`rules_gateway.go:45-49` 判级用的正是这一对。
- 拼子路径：只许"已规范化父目录 ＋ 条目名"这一种形状；现成先例是 `joinForListing(dir, name)` `internal/tools/fs.go:280`，它的注释逐字把边界钉住了：entries 来自 `Readdirnames`、不带分隔符，"**which is why it needs no Clean step and must not grow into a general path joiner**"（`fs.go:276-279`）。
- **`AGENTS.md §1.2` ban #2 的仪器射程（实测）**：`tools/d22scan/main.go:737-741` 只在 AST 上匹配 `filepath.Clean` 与 `filepath.Abs` 两枚（`sel.Sel.Name == "Clean" || "Abs"`）。⇒ **`filepath.Join` / `Dir` / `Rel` 根本不扫**。
- **既有绕过先例（具名，⛔ 顺手改）**：`internal/panel/git.go` 里 `filepath.Join`／`filepath.Dir` 做真实文件系统决策共 **14 处**：`:296`（`dot := filepath.Join(dir, ".git")`）、`:308`、`:340`、`:349`、`:394`、`:402`、`:408`、`:445`、`:472`、`:483`、`:508`、`:517`、`:520`、`:526`；外加 `internal/panel/composer.go:345`（`return filepath.Join(dir, name)`，其注释 `:338` 逐字讲"ticket 75 caught exactly that in a hard-coded separator"）。尺：`grep -rn 'filepath\.\(Clean\|Abs\|Join\|Rel\|Dir\|Base\)' internal/panel/*.go | grep -v '_test.go'` → 15 行，`rc=0`；其中 **`Clean`/`Abs` 命中 0 枚**。
  这 14＋1 枚**都不在** `tools/d22scan/allowlist.txt` 里——因为仪器不要求它们在。登记在册的 `pathresolver-bypass` 豁免只有 3 条：`allowlist.txt:3`（`internal/memory/open.go`）、`:4`（`internal/models/manifest.go`）、`:6`（`internal/risk/pathresolver.go` 自身的"sanctioned lexical step"）。
  ⇒ **对乙形的实际后果**：面板树若只用 `Join` 拼子项，**仪器不响**；但那是"仪器窄于规矩"的缺口（与 `AGENTS.md §1.2` 末段那处 emoji 缺口同一种形状），不是"合规"。落地腿要么复用 `joinForListing` 这一族并把它升成有主的名字，要么就老实走 `Canonicalize` 每级——**这个选择本身该由编排者裁，本程不替他裁**。

### 3.3 忽略规则（`.git`／`node_modules`／构建产物）——乙形的真价钱
- 票 AC#1 逐字要求"忽略过滤必须问 git index，问不到就**全扫并响亮自陈**"，并指名复用 `d22scan` 已定的规矩。
- 那枚实现是**存在的、且是唯一的一份**：`tools/d22scan/gitignore.go`（读 `.gitignore` ＋ spawn `git ls-files` / `git ls-files -i -c --exclude-standard`，`:276`/`:280`/`:311`；不可用时返回 `gitIndexUnavailable` 而不是默默放行，语义写在 `:27-33`、`:39-52`、`:64-88`、`:114-119`）。
- ⚠ **但它是 `package main`**（尺：`sed -n '1,10p' tools/d22scan/gitignore.go | grep -n package` → `1:package main`，`rc=0`）⇒ **`internal/panel` 与 `internal/tools` 都 import 不到它**。
- 尺（全仓有没有第二份实现）：`grep -rln 'gitignore' internal/` → **无输出，`rc=1`** ⇒ 内部包零份。
- ⇒ 乙形要满足 AC#1，只有两条路：**(a) 把 `tools/d22scan/gitignore.go` 抽成一枚 internal 包**（动的是 CI 仪器本体，`AGENTS.md §1.1` 没禁它，但它是 `SPEC-12` 那条"d22scan 自动扫"的执行者，抽取属**跨票工程**，且票 190 的 AC#5 写的是 `allowlist.txt`／`thresholds.go` 一字节不动）；**(b) 在 internal 里重写一份**——AC#1 逐字禁止（"别另发明一套"）。**这是乙形最大的一块代价，且它不在"画树"那条腿上。**

---

## 4. 问③ — 边界与代价

| 形状 | 现有解析器下的真行为（指到行） |
|---|---|
| **工作区根从哪取** | 不是"进程 CWD"。取法＝配置 `[fs] allowed_dirs`（`docs/specs/SPEC-03-config-and-secrets.md:35` 默认 `[]`）→ `NewPathCanonicalizer(allowed, cfg.FS.ReparsePointExceptions)`（`cmd/wisp/run.go:501`）→ `Roots()` `paths.go:246`；面板侧再窄化一枚 `WorkspaceRoot()` `paths_workspace.go:67`（`risk.Result`，零值＝未设，`panel/workspace.go:47-48` 注释逐字"the zero Result (Canonical == \"\") means unset"）。⇒ **空根是正常态**（`UnsetWorkspaceView()` `:65`），所以"工作区未窄化时报的是授权根而不是整个盘"这条判据**有现成的两态可依**，AC#2 那句"一发判据"不必新造状态。 |
| **符号链接／junction（重解析点）** | `risk.Resolve` 里 `reparseComponents(p)`（调用点 `internal/risk/pathresolver.go:117`）逐段比对 `exceptions`，**不在册即 `return res, ErrReparseDenied`**（`:123`；错误变量本体 `:36`，英文整句 118 字符）。`Result.Reparse` 只是"穿过在册联接点"的记账（`:120`），面板视图会把它拼进理由句（`panel/workspace.go:99-101`）。**⚠ 平台形状**：`reparseComponents` 在 Windows 才有实现（`internal/risk/pathresolver_windows.go:60`），POSIX 那份是 `return nil`（`internal/risk/pathresolver_other.go:12`）⇒ **非 Windows 构建下 C26 看不见任何符号链接**，跨界只能靠后续 `InAllowlist` 兜。本项目的 D7 把平台范围钉在 Windows，所以这条不是 bug，但**树腿的测试不能在 POSIX 上声称自己验过跨界拒绝**。既有现量 nails：`internal/tools/bridge_junction_windows_test.go:278`（用例名 `"fs.list_through_a_real_junction"`）、`:290` 断言 `victim.txt` 不得出现在输出里。 |
| **列表里"看见"符号链接 ≠ "穿过"它** | `fi()` 只 `os.Lstat` 条目名（`fs.go:258`），把类型标成 `l`（`:264`），**不调用解析器**；stat 失败那一行写死 `"? %s stat-error"`（`:269`），**不带 err.Error()**。⇒ 单层枚举对符号链接是安全的；**危险只在"往它里面下钻"**，那一步必须走 `Canonicalize`，否则就是 `AGENTS.md §1.2` ban #2 的实质违规。 |
| **超大目录** | `fs.list` 今天**先全量读、再排序、再截断**：`Readdirnames(-1)`（`fs.go:216`）→ `sort.Strings`（`:220`）→ `limit` 截断（`:221-228`），默认上限 500（`fs.go:56`）。⇒ 一棵树如果照抄这一族，**每个目录都是全量读**；本仓现量：`frontend/node_modules` 顶层 52 枚（`ls frontend/node_modules | wc -l` → `52`，`rc=0`），`frontend/dist`、`docs/`、`.git` 同层俱在。真正的钱不在 500 这个数，而在**"全量读＋排序"发生在树的每一个节点上**。 |
| **错误整串往外抛（同类形状，具名，⛔ 顺手改）** | ⚠ 我在枚举腿上看到与 `ErrReparseDenied` 前科同类的三处，全在 `fs.list` 自己的 `Execute` 里：`internal/tools/fs.go:208` `Result{Text: "路径无法解析（按 fail-closed 拒绝）：" + err.Error(), IsError: true}`——这里的 `err` 正是 `t.d.open()` 透出的 `risk` 错误，**`ErrReparseDenied` 那句 118 字符的英文全串会原样进 `Result.Text`，也就是模型正文／面板正文**；`fs.go:212`（`"打开目录失败：" + err.Error()`，`os.Open` 的 `*PathError` 会把**绝对路径＋系统错误文本**一起带出）；`fs.go:218`（`"读取目录失败：" + err.Error()`，`Readdirnames` 同形）。同文件同族另有 `fs.go:144`（`fs.read` 的 `:208` 孪生）。对照面：同函数里的 `:269` 已经示范了**不外抛**的写法。⇒ **这一条不属于票 190 的任何一格，本程零改动，只具名记账**；`delete_enabled`／`privacy.redact_paths`（`SPEC-03:35`/`:37`）是否要求遮蔽绝对路径，归落地腿问。 |

---

## 5. 问④ — 雷区：树上的点击会发起什么动作

- **设计稿里的点击（现量，`design/doubao/demo/rb-files.js`）**：整只文件面板只有**一枚**监听器（`:124 el.addEventListener('click', ...)`）。两条支路：
  - 目录行（`:126-135`）：`next.classList.toggle('collapsed')` ＋ chevron 类切换——**纯 DOM 状态，零宿主调用**。
  - 文件行（`:137-140`）：`app.toast('打开文件 ' + f.dataset.file)`——**只是把路径字符串拼进一条 toast**，`data-file` 的值来自硬编码字面量（`:50-52`、`:80`、`:108` 甚至把 `D:\work\workspace\projects plans\Wisp` 整根写死在面包屑里）。
  ⇒ **"打开文件"在 demo 里从来不是动作，是文案**。真实装里它要变成动作就得是一条新的入站方法＋一条新的宿主能力（开浏览器／开编辑器＝exec 面），**那一枚才是需要 owner 逐枚批的东西**，不是这棵树的显示。
- **写侧**：票 AC#3 那四枚词（删除／移动／重命名／新建）在现量上一一有主：`fs.write`／`fs.edit`／`fs.trash`／`fs.move`／`fs.delete` 全在 `BuiltinFSWriteEntries(d FSDeps) []Entry`（`internal/tools/fs_write.go:772`），其中 `fs.delete` 只在 `d.DeleteEnabled` 为真时才 append（`fs_write.go:779`），该开关的注释本体在 `internal/tools/fs.go:46-49`（"fs.delete is NOT REGISTERED …"）；读半族 `BuiltinFSEntries` `fs.go:325-331`（append 写半族在 `:330`）的注释同样逐字写明这一点（`fs.go:320-324`）。**都是模型侧工具，都不是宿主读面**；加上 `SPEC-03:35` 逐字 `delete_enabled(bool)=false`（票 190 第 26 行已自记此凭据）。⇒ 显示腿**没有任何一条现成通路**能顺手走到写侧。
- **批准侧（`AGENTS.md §1.2` L2 禁令单列，本程只摆证据）**：
  - 仪器：`tools/d22scan/main.go:23-29` ban #6 = "`approval.decide` in frontend/ (D33/F2: allow decisions are native-side only)"，且它扫的是 **`frontend/` 整棵文本树、无后缀白名单**（`main.go:271`、`:417`、`:426-433`）。
  - Go 边界尺：`internal/panel/l2_grant_boundary_test.go:4-8`（"The panel-side L2 'allow' ban, nailed on the GO BOUNDARY"），`:1533` 的判语逐字：Go 若答了一个由 `grantRoutePrefixes × grantRouteWords` 拼出来的审批形状名字，就是"a panel-side allow door whose spelling never appears in the package"；`:1176`、`:1308`、`:2002`、`:2332` 同一族。
  - 命名规矩的既有条文：`internal/panel/bridge.go:46-68`——`config.get`/`config.set` **故意不带 `panel.` 前缀**、且**名字里不许含裁决词**，理由逐字写在 `:62-65`："a route that lets the page change configuration is not a route that lets the page *approve* anything (AGENTS.md §1.2 keeps panel-side L2 'allow' off the table entirely)"。⇒ **任何"文件树"新方法的命名都要照这一条走**，而且它会同时被跳 1 的 `composerRouteLiterals()`（`composer_test.go:409-417`）与 `git.go` 的 `whitelistMethodsFromSource` 抽取器（`bridge.go:55-58` 指向）读到。
  - 越界刷屏的机理（现量）：`rules_gateway.go:45-49` 对"授权目录之外"**只吐 L2**、签名里没有 deny 出口 ⇒ 若树里出现范围外条目，每条都可能变成一张 L2 卡，而 L2 的**唯一应答通道是原生卡**（`internal/tools/bridge.go:457-458`，`b.gate.PendingApproval`），超时自动**拒绝**（`:464-465` 注释逐字："an unanswered L2 auto-REJECTS (ticket 21), deliberately the opposite polarity from L1"）。**面板侧永远不是那个答案**。

---

## 6. 甲／乙代价对照（一句话结论）

**甲形＝用"给模型的一层文本工具"去喂"一棵要多层的界面树"**：7 跳里 3 跳是闸门（C17 名册／Q-56／唯一门判据），且每个展开目录都重穿一次 `Execute→Assess→route`，还顺手把 `fs.go:208/212/218` 那三条外抛错误变成界面文案的原料；**乙形＝在 `PathCanonicalizer` 上加一枚只读枚举、把结果 rides 在已有快照腿上**（`panel_pump.go:83/229` 已是同形状先例），4 面里唯一贵的不是画树、是**把 `tools/d22scan/gitignore.go` 那枚 `package main` 的忽略器变成可复用**——所以**乙形便宜**，前提是"bounded depth ＋ 展开为纯前端状态"（demo 已经证明点击可以零宿主调用，`rb-files.js:124-135`），一旦要做"点了才向宿主要到下一层"的懒展开，就必须新增一条 `panel.*` 方法＝**C17 变更＝人工批准**，那一刀不是工时问题。

---

## 7. 我这把尺的缺陷（自报，不许藏）

1. **管道吞 rc**：`grep … | head`／`| grep -v _test.go` 之后 `$?` 是管道末端而非 grep。§1 那三枚计数我用了 `; echo rc=$?` 紧跟 `grep -c`，是可信的；§2.1 那张"2 处命中"表的 `rc2=0` 也是 `grep -rn … | grep -v …` 的末端 ⇒ **它证明的是"末端命令成功"，不证明 grep 自己没返回 1**。§3.2 双向 import 那两把我特意改用 `grep -rl`（无管道）重测，`rc=1` 才是"零边"的证据。后续腿别把我的 `rc=0` 当"命中数已验"。
2. **`wc -l` 不是行数**：票 190 我报 26 行；末行若无换行符会少 1。本票 `Read` 与 `wc -l` 一致（26），无差异，但尺本身不保证。
3. **枚数尺只数 markdown 复选框**：`- [ ]` 之外的"未做"（例如票面 `Status` 行里的"先只读设计核，再落地"）不进分母；本票 AC#6 里那句"`internal/panel` 2 枚已知红"是一个**没有框的活**，我的 5 枚分母看不见它。
4. **我没有跑任何 go 面**：AC#6 的三把尺（panel 测试、`d22scan` 基线 433、`gate-clauses` 名册）**一条都没核**；`d22scan` 的 ban #2 射程我是**读 `main.go:737-741` 的 AST 匹配条件**得出的，不是跑出 findings 得出的——若仪器另有前置过滤（`allowlist.txt` 之外的 drift guard 分支），我的"Join 不扫"结论要打折。
5. **`frontend/node_modules` 顶层 52 枚＝只数了一层**，`find` 没下钻（怕撞资源＋怕给假"小"的印象）。真实"超大目录"代价我没量化，只指到了"每节点全量读＋排序"那一行。
6. **Q-56／票 145 我都是引用注释原文，未读台账**：`docs/reports/pending-and-issues.md` 我按硬约束 2 没打开。⇒ "Q-56 属待人拍板"这一判定来自 `AGENTS.md §1.5` 对 `Q##` 类别的定义＋`bridge.go:883` 的自陈，**不是我对台账的现读**。
7. **前一程（`9b6a1be5`）我只抽了标题层级（11 枚小节）与两处结论**，29 KB 正文没逐格复跑 ⇒ 下节"差集"只到"射程不同"这一层，不声称与它的每一格相容。
8. **关联票未逐枚读**（145／102／77／119／126／182 全文），我只在 §2/§3/§5 需要处点了名。

---

## 8. 与 `9b6a1be5`（前一枚 `190-a1`）的差

- 那一程的射程（按其小节标题）＝§D① 根的取法、§D② `allowed_dirs` 是判级输入、§D③ 必须走 `PathResolver` 并写成判据、§D③④ 深度／条目上限／忽略规则（建议值）、§D⑤ 只读边界的**结构性**判据、§D⑤后半"与 §A 那堆是不是同一枚／并进 145 还是另立"、AC#4 三形的设计口径、§9 未裁完、§10 我没做的事。
- 本程的增量＝**没有重复那一套**，只补三样它没有的东西：① 甲形的**逐跳清单与 7 跳计数**（含 `bridge.go:872/883` 的 Q-56 与 `composer_test.go:409/502` 那道门），② 乙形的**落地文件清单**（含"装配根已持句柄 `run.go:275/501`"与"panel↔tools 零 import 边"两条现量），③ **忽略器在 `package main` 里、internal 无第二份**这一枚新发现的代价（`gitignore.go:1` ＋ `grep -rln 'gitignore' internal/` rc=1）。
- 一处**口径不同、不是冲突**：那一程 §D③ 把"必须走 PathResolver"写成判据思路；本程补的是仪器的**实际窄处**（只 `Clean`/`Abs`，`main.go:738`）＋ `internal/panel/git.go` 那 14 枚 `Join`/`Dir` 先例。两者合起来才是完整形状：**规矩比仪器宽，先例已经在宽处生活着**。本程不改任何一处，也不主张把 `git.go` 迁进 C26。
