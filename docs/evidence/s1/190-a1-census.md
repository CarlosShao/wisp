# 190-a1（2026-11-04 程）— 工作区文件树（普查 S4 那一堆）：取数、落点、雷区、边界

- 性质：**只读设计核**。零产码、零勾框、票面一字未动、台账未动、只 commit 本文件一枚。
- 工单：`.scratch/wisp/issues/190-the-workspace-files-stack-file-tree-has-no-host-side-source.md`
- 派单：`.scratch/wisp/probes/orch/20261104/txt/190a1-dispatch.txt`（7 行）
- 前序件：`docs/evidence/s1/190-file-tree-design-a1.md`（09-28 设计核，**行号今日全漂**，见 §5-7）；
  `.scratch/wisp/issues/*-190-*.md:28-44`（10-08 第二程只读日志，7 跳/4 面那张表）
- 答案四问：①取数＝**复用同一读面族的路线**，但忽略引擎**不在 git.go 里**，且**不许搬** d22scan 那 817 行（§1）；
  ②落点＝**5 枚改 + 2 枚新增，约 600–760 行**，**不触 C17**，**加顶层 key 但不新增红**（§2）；③雷区 10 条 + 4 枚在册红（§3）；
  ④与 181/189/167/145 的边界（§4）。

## 0. 起手三件

```
date        = 2026-11-04 14:1x +08（本机）
HEAD        = bfd5ace3d6fd80081f18c453c7477141b06c09fd
              "orch129:派腿存档 189-r1(面板审查读面写腿)+190-a1(文件树只读普查);发现面板补功能真脉"
              （2026-11-04 14:18:49 +0800，即本程派单的那一枚）
branch      = dev
写面闸门    = git status --porcelain -- internal/ cmd/  ⇒ 只有 `?? cmd/wisp/zz_probe_temp_test.go.txt`
              （**无 tracked 改动**；那枚未跟踪 .txt 不是本程的，未动、未删）
```

**在飞腿与它对本程的约束**（这是落点那一节的第一约束，不是背景）：同批派单的 `189-r1`
（`.scratch/wisp/probes/orch/20261104/txt/189r1-dispatch.txt`）要写 `internal/panel/git.go`
＋**新增** `internal/panel/git_readonly_exec.go` ＋两枚判据件。现量 `Test-Path internal/panel/git_readonly_exec.go`
＝ **False**（它还没落盘）。⇒ ①`internal/panel` 每波只许一枚写腿（`AGENTS.md` §⓪③），190 的落地腿**必须排在 189-r1 交件之后**；
②它那枚 exec 助手正是 190 要复用的那一枚（§1.2 C 的一半，§2.1 末）。

**行号一律现跑**：`internal/panel/git.go` 今日 **580** 行（票面与派单写 569，09-28 件写 559）——
三份历史读数都对不上，本件全部用现跑号。

## 1. ① 数据从哪取

### 1.1 先把三件事分开：路线 / 根 / 忽略

"复用票 181 的 `git.go`"这句话在这三件事上的答案**不一样**，混着说就会得出错的落点。

| # | 需要的那一件 | 今天在哪（现跑） | 能直接复用 `git.go` 吗 |
|---|---|---|---|
| 1 | **载体与路线**（宿主侧读面，搭既有快照推一次，不是模型工具，不起新通道） | `git.go:5-11`（"answers ride the snapshot push that already exists"）、`:13-16`（"It is not a tool"，D34 无 git 行）、`pump.go:152`（`Git func() GitView`）、`panel_pump.go:276-278`（`gitView()`）、`cmd/wisp/run.go:904-908`（装配字面量，`Git: rt.gitView` 在 `:908`） | **必须复用**。另起第二份读面/第二条通道＝票 181 AC#3 与 189 派单 §"⛔ 另起第二份"逐字禁 |
| 2 | **根**（工作区窄化后的那一枚；未窄化时＝授权根集合；可 0 枚） | `internal/tools/paths_workspace.go:67`（`WorkspaceRoot() risk.Result`）、`internal/tools/paths.go:256`（`Roots()`）；面板侧今天只有 `PathScope.WorkspaceRoot()`（`internal/panel/workspace.go:45-55`）＋`ComposerState.Workspace`，**没有 `Roots()`** | **路线复用、值不来自 git.go**。树腿不新读 `config`（那会造第二本账），由装配点把 `rt.paths.Roots()` 交进来 |
| 3 | **忽略规则**（`.gitignore` ＋ index 两半） | ⛔ **`git.go` 里一枚都没有**：全文 580 行只读 `.git` 的元数据（`HEAD`／`refs/**`／`packed-refs`／`worktrees/**`／`commondir`，读取器在 `:190-546`），零 `.gitignore`、零 index。今天全仓唯一一份在 `tools/d22scan/gitignore.go`：**`package main`（`:1`）、817 行**；`internal/` 侧 `grep -rln gitignore internal/` ⇒ **无输出** | **不能复用**（import 不到 `package main`），**也⛔不许抄一份** |

### 1.2 裁定（问①的答案）：三段式

**A. 复用路线。** 在 `package panel` 里新增**同族文件**（建议 `internal/panel/tree.go`），经一枚新的 `PumpSources` reader 进**同一枚** `Snapshot`。
不新开通道、不注册 Tool、不加 `panel.*` 方法。判据形状直接抄 `git.go` 的：
四态枚举＋**必填 reason**（`:64-72`）、"读不到不许塌成没有"（`:39-42`）、空数组永不 nil（`:118-121`）、
"这一维没有读数"与"读了但没有"两句两值（`workspace.go:76-81` 那对是同一纪律）。

**B. 复用根，且 panel 仍旧不做路径决定。** 树腿的输入是"已经 canonical 的根集合"＋一枚**既有** seam
`risk.PathCanonicalizer`（`internal/risk/assessor.go:141-144`，两个方法 `Canonicalize`/`InAllowlist`）——
`tools.PathCanonicalizer` 就是它的实现（`paths.go:131`/`:211`）。**panel 不必 import internal/tools**：
今天 `internal/panel` → `internal/tools` 的真 import **只有测试件一枚**（`workspace_handlers_test.go:21`），
反向**零枚**；这条双向零依赖是仓里既有的架构事实，别为了一棵树把它破掉。
每下钻一层过 `Canonicalize`（撞 `ErrReparseDenied` ⇒ 该条目标"链/接合点，未跟随"，⛔ 不跟随后再判）。
这就是 `git.go:54-62` 那句 "it makes no path AUTHORISATION decision" 在树这一维的同款写法。

**C. 忽略问 git 本体，不搬 matcher。** `tools/d22scan/gitignore.go` 那 817 行实现的是
gitignore(5) 的一个子集（父目录规则、最深文件、行内取最后一条、行首空白、尾部 `**`、`!` 反转、`[...]` 类），
并且它还要问 index 才能回答"被跟踪的文件永不忽略"（`:35-47` 逐字）。**树不该抄这 817 行**，直接问 git 两个问题：

- `git ls-files -z` ⇒ tracked 集（本仓实测 **18,574** 条）
- `git ls-files -o -i --exclude-standard -z` ⇒ ignored 集（本仓实测 **8,367** 条；带 `--directory` 收敛成 **79** 条目录）

"这条在不在 ignored 集 / 在不在 tracked 集"就是全部答案；**tracked 永不忽略**这条最贵的语义由 git 自己保证
（d22scan 为它写了 60 行注释＋一枚强加跟踪的判据，抄的时候最容易抄丢的正是它）。
降级方向两边一致：**问不到 ⇒ 一条规则都不套、全扫、响亮自陈**（`gitignore.go:49-62` 逐字就是这个规矩，
票 190 AC#1 也是这一句）——所以这不是"另发明一套"，是把那条权威当尺。

⚠ **但它有一处必须自己加，且不能指望 git**：`.git` 本身。现跑 `git ls-files -o -i --exclude-standard | grep '^\.git/'`
⇒ **0 条**；`git check-ignore .git` ⇒ **rc=1**（git 不认为 `.git` 是 ignored）。⇒ 只信这两张集合，`.git` 会被当**普通目录列出来甚至下钻**。
d22scan 自己也是靠**名字级硬编码**兜这一类（`tools/d22scan/main.go:680`/`:1029`/`:1107`/`:1356`，三枚名字 `testdata`/`.git`/`node_modules`）。
树只该抄其中**一枚**：`.git`（理由＝它是 git 自己的元数据、不是工作区内容）。
**`testdata` 与 `node_modules` ⛔ 不许抄**：`node_modules/` 已写在 `.gitignore:17` 里（git 会回答），
`testdata` 是**真实源码目录**——抄它会静默少列，正撞 AC#4。

## 2. ② 落点：哪些文件 / 约几行 / C17 / 顶层 key

### 2.1 文件与行数

| 文件 | 改动 | 约行 |
|---|---|---|
| `internal/panel/tree.go` | **新增**：`TreeSection`/`TreeRoot`/`TreeEntry` 三型 + 四态枚举 + 读取器（根循环、下钻、`Canonicalize`、git 两问、三上限、四元组） | **300–360** |
| `internal/panel/tree_test.go` | **新增**：合成目录树（真目录，断言条目集合逐名）＋只读边界常驻钉（类型面无写落点/不注册 Tool/入向 hop 不存在）＋不跨 junction＋截断四元组＋"问不到 git 就全扫并自陈" | **260–330** |
| `internal/panel/composer.go` | `Snapshot` 加一枚 `Files *FileTreeSection `json:"files,omitempty"`` ＋注释（照 `:74`/`:91` 那两枚的形状与理由写） | **8–12** |
| `internal/panel/pump.go` | `PumpSources.Files func() FileTreeSection` ＋`Snapshot()` 里四行（照 `:310-315` 的 instructions 形） | **8–12** |
| `internal/panel/display_scrub_216.go` | 树的文件名/路径必须过 `scrubForDisplay`（§3.5） | **6–10** |
| `cmd/wisp/panel_pump.go` | `fileTreeView()`（取 `rt.paths.Roots()` ＋ `rt.workspaceView()` ＋ `rt.paths` 作 seam，交 `panel.ReadTree`） | **20–30** |
| `cmd/wisp/run.go` | 字面量内一行 `Files: rt.fileTreeView,`（`run.go:904` 那一块） | **1–3** |
| 合计 | 5 改 + 2 新增 | **≈ 600–760** |

⚠ 如果忽略那一半改走 189-r1 的 exec 助手（`internal/panel/git_readonly_exec.go`，今天盘上不存在），
上面"起两条 git 命令"的纪律（`--no-optional-locks`、ctx 截止、owner/recover、单行错因）**一律复用它**，
⛔ 在 `tree.go` 里再写第二份 `exec.Command`——那是"重造第二份"在执行面上的那一半。

### 2.2 C17 判：本形**不触**，懒展开才触

现状名册三把尺（现跑，全部按源码读，不是重列常量）：

| 尺 | 位置 | 现量 |
|---|---|---|
| `panel.*` 字面量 | `internal/panel/git_test.go:383-391`（提取器 `:409-420`） | **恰好 4 枚**：`panel.mode.request`/`panel.workspace.request`/`panel.attachment.add`/`panel.message.send`（`bridge.go:41-45`） |
| 闭名册全量 | `internal/panel/inbound_roster_253_test.go:63-80` | **7 枚**（4 枚 `panel.*` ＋ `config.get`/`config.set`/`config.models.list`），`wantFullInboundRosterSize253 = 7`、`wantNonPanelPrefixedInbound253 = 3` |
| 全包扫 | `internal/panel/composer_test.go:326-340` | 走 `internal/panel` 全部非 `_test.go`，任何 `panel.*` 字面量必须是 `bridge.go` 声明过的 |

⇒ **"点了才向宿主要下一层"（懒展开）＝新增一枚 `panel.files.*` ＝ 上表三把尺同时红 ＋ `SPEC-12 §4.1` 人工批准**。
这与 167 的 `Q-76` 是**同一个形状**（它要第 5 枚 `panel.*` 才让"停止"按得下去）。

**本程的选择：v1 不做懒展开。** 但理由**不是**"怕动契约"，而是 §3.3 的实测把"必须懒"这个前提**暂时证伪了**：
本仓深度 2 的整棵树只有 **406** 个节点，一次推得起。真到界面需要无限深时，再加那枚方法并把它摆给机主
（那时它是一枚有价值的契约变更；`Q-76` 的先例在那儿）。技术选型按价值与主流 harness 走，
当前 value/price 比最高的是"一次推、有上限、截断可见",而不是"先动契约换懒加载"。

### 2.3 顶层 key vs 子格：**派单的前提只对一半**（本件最重要的一处更正）

派单/常识说"加顶层 key ⇒ 撞 `TestComposerContractTypesMatchFrontend`，并撞 `pump_test` 两条字节钉"。现量：

- **两条字节钉** `internal/panel/pump_test.go:109-125` 与 `:280-293` 断言 `composer,generatedAt,pending,results`
  **恰好四枚**。但这两发跑的 pump **不接** `Instructions`/`Tasks` reader，而那两枚是 `*T`＋`omitempty`
  （`composer.go:74`/`:91`）⇒ **它们今天是绿的，正说明"`omitempty` 指针型顶层 key 不会让这两条钉红"**。
  会让它们红的是**不带 omitempty** 的新 key（票 145 当时算的就是那一种）。
- 会响的是两枚**双向对账**：`approval_test.go:105-136` 的 `pairs` 含 `Snapshot↔PanelSnapshot`（`:118-121`）、
  `composer_test.go:52-85` 的 `pairs` 也含 `Snapshot`（`:59-70` 六型）。**这两枚今天已经是红的**——
  在册 4 枚之二；红句逐字（197 证据 `:281`）：`approval_test.go:129: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`。

⇒ 三形的真实代价：

| 形 | 新增红 | 加长既有红 | 两条字节钉 | 216 洗（文件名＝不可信串） | 语义 |
|---|---|---|---|---|---|
| **A 顶层 `files`**（`*FileTreeSection`＋omitempty） | **0** | 2 枚（两个 `Snapshot` 行各多一个名） | 不动 | 需在 `NewSnapshot` 里**显式**洗 | 与 `instructions`(200)/`tasks`(197) 同形＝**一枚新维度** |
| B `composer.tree` | 0 | 1 枚（`ComposerState` 那行；`git` 已在那行里） | 不动 | 进 `scrubComposerForDisplay` 就顺带洗 | 把工作区那一维塞进"输入行"（`composer.git` 是这种先例） |
| C `composer.git.tree` | 0 | **0** | 不动 | 要自己加 | **语义错**：非 git 工作区（`kind=not_a_repo` 是正常输入，`git.go:66-72`）就**没有树** |

**选 A。** 理由按价值排，不按代价排：

1. 文件树是**一枚维度**，不是输入行的一部分；本仓已经用同一形加过两枚维度（`instructions` 200 程、`tasks` 197 程），
   两枚都活下来了 ⇒ 主流 harness 的"packet 一枚 section 一枚维度"在这里**有先例、有判据、有界面侧对接路径**（`PanelSnapshot` 补声明，归 `Q-51`）。
2. **C 的零代价是假象**：`GitView` 不进任何 `pairs` 名册（`composer_test.go:59-70` 六型逐名），
   所以挂在它下面**没有任何尺会看着它**——那不是"没代价"，是把代价从"红一枚测试"换成"永远没人看得见"（静默漂）。
3. B 的唯一好处是"那行本来就红"，但一棵树塞进 `ComposerState` 会跟着 `bookPanelSnapshot` 的摘要与每个 `composer` 消费者一起走，
   且下一枚要加维度时又要来问同一个问题。
4. A 多出来的那一步（在 `NewSnapshot` 里洗文件名）**不是成本，是必须**：文件名正是 216 那张名册里"别人给的字符串"的原型
   （`display_scrub_216.go:8` 逐字 "A file name is somebody else's string"）。放 B 是"顺带洗到"（运气），A 是写死一行（判据咬得住）。

## 3. ③ 雷区与已知红

### 3.1 在册 4 枚红（逐名 + 红因 + **本程没跑**）

| # | 测试 | 定义处 | 红因（现读） |
|---|---|---|---|
| 1 | `TestApprovalCardViewJSONKeysMatchFrontendTypes` | `internal/panel/approval_test.go:105` | `Snapshot` 双向对账；197 证据 `:281` 逐字红句 `Go Snapshot emits [instructions tasks] …` |
| 2 | `TestComposerContractTypesMatchFrontend` | `internal/panel/composer_test.go:52` | 同一把尺的 composer 侧；`git` 从 181 起就在红句里（181 日志逐字 "Go ComposerState emits [git] …"），后加的 `credentialState` 等只加长它 |
| 3 | `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` | `internal/panel/frontend_hygiene_test.go:196` | 红因在 `frontend/**`（本编队零写面） |
| 4 | `TestC21DesignTokensFourWayAgree` | `internal/panel/tokens_fourway_test.go:617` | 读 `design/assets/tokens.css`；本机 `Test-Path` ⇒ **False**（文件不在盘上）＝环境性红，`design/**` 归机主 |

- 名册出处＝`docs/evidence/s1/197-subagent-carrier-r3.md:268-275`（最后一枚 commit 之后取的那一发）。
- **本程零 Go 命令**，所以这是**引文不是读数**。理由：`189-r1` 同包在飞（跑一次 `go test ./internal/panel`
  既可能读到它半写的 `git.go`，也抢编译/测试面）。落地腿必须自己重取，并照 `A380` 的定式：动跨模块面还要单跑 `./internal/risk/`。
- **这四枚没有一枚是本票能修或该修的**；§2.3 的 A 形**不新增任何红**，只让已红的两行各多一个名字（与 197/248/320 每次加维度的处置同形）。

### 3.2 雷区十条

| # | 雷 | 现量凭据 | 处置 |
|---|---|---|---|
| 1 | **树上写动作**（删/移/重命名/新建） | 写族唯一入口 `internal/tools/fs_write.go:772`（`BuiltinFSWriteEntries`），`fs.delete` 只在 `:779 if d.DeleteEnabled` 时挂；面板入向 7 枚里**零枚**写文件 | 结构性：树读取器**签名不接面板来的路径**（根自取）＋不新增入向方法（§2.2）⇒ 拓扑上不可达 |
| 2 | **`.git` 被当普通目录**（甚至下钻） | 实测 ignored 集里 `.git/` **0 条**、`git check-ignore .git` **rc=1** | 树自己加**一枚**名字级排除；⛔ 不抄 `testdata`/`node_modules`（§1.2 C 末） |
| 3 | **大目录/git 查询成本** | 实测本仓：`git ls-files -z` **97–112 ms**；`git ls-files -o -i --exclude-standard -z` **412–509 ms**；加 `--directory` 后 **297 ms**（ignored 8,367 条收敛成 79 条目录）；非仓库目录里 `git ls-files` 失败 **112 ms**（rc=128） | 两条 git ＝ **≈0.4–0.6 s/次**；publish 点在"一张卡开/关、一次工具调用起/止、流结束"（`panel_pump.go:441-449`）⇒ **必须缓存**：缓存键取 `.git/index` 的 mtime+size（stat 是缓存键，不是"墙钟超时"，ban #4 禁的是后者） |
| 4 | **包体积** | 实测本仓深度≤3 的节点 **1,531**（depth1=**24** / depth2=**382** / depth3=**1,125**），名字串合计 **50,108** 字符；最大扇出 `.scratch` **282** > 建议每目录上限 200 | JSON 量级 **100–150 KB/包**（估算，未序列化）；对照今天最小包 **552 B**、带一张 L2 卡 **792 B**（`panel_pump.go:348-354` 逐字）⇒ 树一枚维度≈整包的 **100–270 倍**。⇒ 默认 **depth 2**（406 节点）、perDir **200**、total **500**；截断必须带四元组（shown/total/truncatedBy/outOfScopeNote），**"被截断"与"被忽略"分成两个数** |
| 5 | **文件名＝不可信串**（216 那条的正面） | `display_scrub_216.go:8` 逐字；`scrubComposerForDisplay`（`:223-247`）今天只洗 `Workspace.Canonical`/`Reason`＋`AttachmentRef.Name` | 树整段在 `NewSnapshot` 那一跳过 `scrubForDisplay`（§2.3 A 第 4 条）；⛔ 在 producer 里洗（会把转义串喂回路径读取） |
| 6 | **exec 面** | `internal/panel` 今天 `exec.Command` **0 枚**（生产 exec 只在 `internal/winsec/launcher_*.go` 与测试件）；d22scan 自己也只用 ctx 截止（`gitignore.go:256`/`:370-389`） | 复用 189-r1 那枚 exec 助手；⛔ 裸 `go func`（ban #1）、⛔ 墙钟超时（ban #4）；路径决策仍过 C26，⛔ 在 `risk.PathResolver` 外 `filepath.Clean|Abs`（ban #2） |
| 7 | **越权三形不许静默**（AC#4） | 枚举先例 `git.go:64-72`；句对先例 `workspace.go:76-81` | 三个值分开：`not_chosen`（读了，没有工作区/没有根）/ `not_read`（泵没接 reader）/ `permission_denied`＋`unreadable`；`total==unknown` 与 `total==0` 必须可区分 |
| 8 | **"树腿不接受面板传根"** | 面板侧零 `Roots()`（§1.1 第 2 行） | 根由装配点给；这条同时把"C17 不加方法"变成"没有新方法来" |
| 9 | **烫手词 `worktree` 已废**（181 日志那条雷，现已不成立） | 实测 `gitSwitchCapabilityRe` 在 `internal/panel/*.go` **零命中**；`composer_test.go:269-296` 逐字说明那枚词面尺**已被换掉**，现在是 `:297 TestThePanelSwitchCapabilityRidesOneWhitelistedRoute` 按**结构**判 | 落地腿**不必**再为措辞改注释；但新增 `panel.*` 字面量会让 `:326-340` 的全包扫红（§2.2 尺 3） |
| 10 | **每波一枚写腿** | 本波在飞 `189-r1`（同包、同族） | 190 落地腿排在它交件之后；189 的落点是 `git.go`＋新 exec 件＋判据件，**不碰** `composer.go`/`pump.go`/`panel_pump.go` ⇒ 冲突面在"同包编译面"，不是同文件 |

### 3.3 成本那一组的口径（实测 / 估算分开）

- **实测（本程现跑）**：`git ls-files -z` 97 / 112 ms（两次）；`-o -i --exclude-standard -z` 412 / 509 ms；
  加 `--directory` **297 ms**；非仓库里失败 **112 ms**（rc=128）；tracked **18,574**；ignored **8,367**（目录形 79）；
  节点 depth1/2/3 ＝ **24 / 382 / 1,125**，深度≤3 合计 **1,531**；最大扇出 `.scratch` **282**、`cmd/wisp` 193、`internal/tools` 138；
  名字串合计 **50,108** 字符；目录数（≤3 层、排除 `.git`）**219**。
- **估算（未跑）**：JSON **100–150 KB**；一次 `Eval` 推 100 KB 的界面成本（票 329 的投递面，`panel_pump.go:12-25`）。
- **未测**：每目录一次 `Canonicalize` 的总耗时（本程零 Go 命令、不写探针，见 §5）。

## 4. ④ 与 181 / 189 / 167 的边界

| 票 | 它拥有什么 | 与本票的交界 | 不许互侵 |
|---|---|---|---|
| **181**（done） | `.git` 元数据三问（HEAD/refs/worktrees）＋四态 kind＋C26 改写账户；件 `internal/panel/git.go`（580 行） | 树**不读 `.git` 内容**，只把 `.git` 当**一枚排除项**；路线/形状复用它 | 181 的反向钉 `TestGitDimensionHasNoModelCallableTool`（`git_test.go:364-392`）会**继续数** `panel.*`＝4 ⇒ 与 §2.2 同向，树不许加方法名 |
| **189**（在飞 `189-r1`） | 未提交改动枚数＋逐文件 diff，嵌在 `composer.git.*` 下 | **它那枚 exec 助手是树要复用的那一枚**；两票都"问 git"，但问的**不是同一组命令**（189：`status`/`diff`；190：`ls-files` 两条） | 共用执行器，⛔ 共用查询、⛔ 共用视图字段（189 嵌 `composer.git.*`，190 走顶层 `files`） |
| **167**（立而不派） | 占用/排队/停止/草稿；"按得下去"要第 5 枚 `panel.*`＝`Q-76` 待人拍板 | 190 的懒展开是**同一形状的第二枚** ⇒ 摆给机主时**并成一次**，别分两次问；167 AC#7 的写面冻结（`internal/panel/**` 只到 145 解冻的三枚文件、只到"新增字段"）**同样约束本票** | 树落在**新文件** `tree.go` ＋三枚解冻文件的"新增字段"以内 ⇒ 合规；⛔ 顺手改 167 的 occupant/status 维度 |
| **145**（载体集） | `Snapshot` 的字段集与泵 | 本票**不新开载体**：一枚 section 进既有 `Snapshot`；key 集变化照 §2.3 A 处理 | ⛔ 新造第二个快照结构/第二条通道（`composer.go:14` 逐字 "no second channel"） |
| **146/174** | 显示与真用不许分账 | 树的"忽略/截断"必须与 `git status` 看到的**同一棵树**同源（同一根、同一 index 读数） | 两套忽略规则＝分账 |

## 5. 未裁完 / 未测（明写，不交半份当完整）

1. **未跑任何 Go 命令与门禁**（`./internal/panel`／`d22scan`／`gate-clauses`）。4 枚红是引文；`d22scan` 基线（433）与 `G6neg` 名册**本程未取**。
2. **每目录 `Canonicalize` 的成本未测**：本仓 219 个目录（≤3 层），一次 `Canonicalize` 是句柄级解析，量级未知。落地腿必须实测一次再定 depth 默认值。
3. **`Eval` 一次推 100 KB 的界面成本未测**（票 329 投递面）。
4. **`files` 这个名字与 §2.3 A 的形状要过 `Q-51`**（`panel.ts` 谁写）；Go 侧改名＝一处字段名＋一处 `pairs` 名。
5. **深度/上限建议值（2 / 200 / 500）** 是按本仓实测＋"比 `fs.list` 的 500 更严"给的，**必须现测一次最坏工作区再定稿**；⛔ 为过 SLO 动 `thresholds.go`。
6. **缓存策略**（`.git/index` mtime+size 作键）只是候选，未实现未测；不缓存时的代价已实测（0.4–0.6 s/次）。
7. **工作区是"别人仓库的子目录"这一形未裁**：d22scan 直接判为不可问（`gitignore.go:267-273`），树要按同一句走"全扫＋自陈"，还是验证 `git ls-files` 的相对路径行为再放行——**留给落地腿定并钉住**（倾向跟随 d22scan：不问即全扫，方向永远是多显示不是少显示）。
8. **09-28 同名件 `docs/evidence/s1/190-file-tree-design-a1.md` 的行号全部已漂**（它引 `composer.go:200-207`、`pump.go:114`、`paths_workspace.go:40` 这类数字，今天分别对不上；`paths_workspace.go` 的 `WorkspaceRoot()` 今天在 **:67**）。它的**结论**（根的集合形状、忽略问 index、只读结构性、五条判据）今日复读仍成立，**行号一个都别沿用**。

## 6. 我没做的事（自证口径）

- 零产码：`internal/**`、`cmd/**` 一字节未写（起手闸门只有一枚别人留下的未跟踪 `.txt`，未动未删）。
- 零勾框、票面一字未动（Progress log 归口在编排者，10-08 那程已由他具名裁过同一句）。
- `frontend/**`/`design/**` **零写面、也未读**：§2.3 关于 `panel.ts` 的结论全部来自 Go 侧两枚 `pairs` 名册与 `tsInterfaceKeys` 的机制，不是读 TS 得来的。
- 台账 `docs/reports/pending-and-issues.md` 未动（`A##` 入账归编排者）。
- 零删除（无 `rm`/`del`/`git clean`/`restore`/`checkout .`/`stash`/`reset`/`rebase`/`--amend`/worktree）。
- 未 push；只 commit 本文件一枚，`git add`/`git commit` 两处**显式 pathspec**（只点名本文件）。
- 临时件目录 `D:\tmp\wisp\20261104\190a1\` **只建不删**（本程在它下面跑了两发非仓库 git 计时，无落盘产物）。
- 被拒调用：0 枚。

## 7. 补记（本件落盘之后 10 分钟内，`189-r1` 的在飞件到了盘上）

本件写作时（14:2x）`internal/panel/git_readonly_exec.go` **不存在**（`Test-Path` = False），§2.1 末与 §4 都按那个读数写。
**14:35 现跑**：`189-r1` 的两枚件已在工作树上（**均未提交**，它还在飞）——

- `internal/panel/git_readonly_exec.go` **634 行**（新文件，未跟踪）；
- `internal/panel/git.go` **+11/−3**（`git diff --numstat`），头部已出现 `composer.git.uncommitted` 与 `composer.git.diff`
  两枚子格的描述（`git.go:14`/`:16` 逐一具名），即它按派单把两枚子格嵌在 `composer.git.*` 下。

⇒ **本件的结论不变，且第 1 条被加强**：树要复用的 exec 助手现在**真的有实物可复用**了（§2.1 末、§4 的 189 行）；
§3.2 第 10 条那条"必须排在 189-r1 交件之后"也照旧成立（它此刻仍未交、仍在同包写）。
⚠ 本补记**不**评价那两枚件的质量：它们未提交、未验收，名字与 API 都可能再变；引它们时只引"存在/行数/改的是哪两处"这三件。

## 8. 补记二（同一分钟内，编排者已把 `189-r1` 判为 dud）——§7 最后一句必须加限定

本件落第二枚 commit 后现跑 `git log --oneline -3`，看到两枚 commit 之间编排者已落了：

```
b72c72bdc orch129:189-r1 dud(build失败/类型未定义/无测试),派189-r2收尾腿承接半成品补齐
```

⇒ **§7 那半句"现在真的有实物可复用"必须加限定，以编排者这句为准**：那枚 `git_readonly_exec.go` 是
**半成品（build 失败／类型未定义／无判据）**，由 `189-r2` 收尾。对本票的后果，逐条写死：

1. **190 的落地腿不许把它当"现成助手"直接调**：先看 `189-r2` 是否交件，交件后**自己 build 一次**再决定复用（收件只认盘上三把尺）。
2. **190 也不许为了绕开它自己去写第二份 exec**：那正是"重造第二份"。正确处置是**等 189-r2**，或在派单里明写"exec 助手由 189-r2 提供"。
3. 本件 §2.1/§3.2#6/§4 的**其余结论不受影响**：树要问 git 两问这一点没变、复用执行器的**要求**没变、C17 判没变、
   落点行数区间没变（`git_readonly_exec.go` 是 189 的地界，不计入 §2.1 的 600–760 行）。
4. 这同时是本件"**零 Go 命令**"这一自限的代价，照录不辩：本程没有 build/测试面，所以"这枚文件能不能编译"这件事
   我**只能引编排者的判语**，不能自己出读数——落地腿别把这一段当已验证。


