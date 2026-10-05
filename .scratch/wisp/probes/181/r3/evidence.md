# 票 181 · AC#7 生产侧腿 `181-r3` 证据件

本格只做一件事：让 `panel.WorkspaceView` 的改写账户**有生产者**。票面 AC 框与本件无关（翻勾归编排者），本件只放现量。

---

## §0 起手锚

| 项 | 读数 | 取法 |
|---|---|---|
| 时刻 | `Mon Oct  5 11:12:52 CST 2026` | `date` |
| HEAD | `c6cf66e6`（分支 `dev`） | `git log -1 --format=%h` |
| 写面闸门 | ` M cmd/wisp/firstrun.go`（别程在飞，本程未碰 `cmd/`） | `git status --porcelain -- cmd internal` |
| 本程写面 | `internal/panel/**` ＋ `internal/tools/**` ＋ 本件 ＋ `.scratch/wisp/probes/181/r3/` | 派单 |
| 禁区 | `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`internal/agent/approval/**`／`frontend/**`／`design/**`（后两枚连读都没读） | 派单 |

票面 AC#7 的三枚现量锚本程复跑，全部**复现**（不是照抄票面）：

1. `internal/panel/workspace.go:60-68` 的 `WorkspaceViewFromRoot` 现量只塞 `Set`／`Canonical`／`Reason`，`Rewritten` 留零值 ⇒ 复现（函数体 9 行，无第四枚赋值）。
2. 唯一写真值那一路在 `internal/panel/workspace.go:104-106`（成功分支的 `Rewritten: res.Rewritten`）⇒ 复现，行号未漂。
3. 上游两道拦：`internal/panel/workspace.go:85` 的 `res.Actable()` 与 `internal/tools/paths_workspace.go:64-66` ⇒ 复现；`risk.Result.Actable()` 对任何 `Rewritten` 先返回 `ErrRewrittenPath`（`internal/risk/pathresolver.go:93-99`，票面写 `:94-97`，本程现读为函数体 `:93-99`、`if r.Rewritten` 在 `:94`）⇒ **行号差一枚，结论不变，具名登记**。

---

## §1 起手名册（先跑包再动笔）

尺：`go test -count=1 -v ./internal/panel/ ./internal/tools/`，时刻 `2026-10-05 11:1x +08`，落盘 `/tmp/181r3-baseline.txt`。

* `internal/panel`：**FAIL 3.444s** —— 顶层 `--- PASS` **115 枚**／`--- FAIL` **4 枚**／`--- SKIP` **0 枚**（合计 119 枚顶层用例）。
* `internal/tools`：**ok 18.339s** —— 顶层 `--- PASS` **194 枚**／`--- FAIL` **0 枚**／`--- SKIP` **0 枚**。
* 两包合计顶层 `--- PASS` **309 枚**（逐名计数，非估算；`--- SKIP` 两包都是 **0 枚** ⇒ `TestWorkspaceSwitchRefusesAnExpandedSpelling` 的 `%TEMP%` 子例与 `TestWorkspaceSwitchRefusesAJunctionToOutside` 在本机**真跑到了**，不是靠 skip 蒙过的绿）。

四枚红的**逐字**（全部在册，本程不修不追，理由见 §7）：

```
--- FAIL: TestApprovalCardViewJSONKeysMatchFrontendTypes (0.00s)
    approval_test.go:129: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare
--- FAIL: TestComposerContractTypesMatchFrontend (0.01s)
    composer_test.go:74: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare
    composer_test.go:74: Go ComposerState emits [git currentModel modelKnown credentialState credentialKnown] that interface ComposerState does not declare
--- FAIL: TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme (0.14s)
    frontend_hygiene_test.go:219: 57 files reachable from main.tsx carry colour literals only in the generated theme
--- FAIL: TestC21DesignTokensFourWayAgree (0.00s)
    tokens_fourway_test.go:441: read design/assets/tokens.css: open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css: The system cannot find the path specified.
```

本程改动的**半径名册**（起手绿，动笔前逐枚读过断言原文）：

| 用例 | 包 | 起手 | 它钉的是什么 | 与本程的关系 |
|---|---|---|---|---|
| `TestWorkspaceViewRenderedForSnapshots` | panel | PASS | `WorkspaceViewFromRoot("")` 必须 `Set=false` 且有 reason；`WorkspaceViewFromRoot(`d:\work\a`)` 必须 `Set` 且 `Canonical` 逐字 | **直接被本程改签名**，改名册 §3-1 |
| `TestWorkspaceSwitchAuditsAndAppliesACleanPath` | panel | PASS | 走 `PathScope` 接缝，断言 view 与审计行 | `PathScope.WorkspaceRoot()` 返回类型变更 ⇒ `scriptedScope` 跟改 |
| `TestWorkspaceSwitchKeepsThePreviousScopeOnFailure` | panel | PASS | 拒绝后仍显示在用的那枚（`view.Canonical == `d:\work\alpha``） | `currentAfter` 的比较从 `!=` 换成读账户 |
| `TestWorkspaceSwitchDetectsAScopeThatMovedAnyway` | panel | PASS | 范围偷偷变了必须说，不许渲染成用户要的那枚 | 同上（`risk.Result` 含切片，不可 `==`） |
| `TestGitSectionTravelsInTheSnapshotPacket` | panel | PASS | `git.currentWorktree` 与 `workspace.canonical` 必须是同一枚坐标（"one packet, two trees" 反例） | `:540` 那枚 `WorkspaceViewFromRoot(root)` 跟改 |
| `TestReadGitForWorkspaceConsumesRewriteAccount` | panel | PASS | 消费侧三子例（未改写＝原路／改写要把账户说进 reason／未设不探测）＝`181-r2` 那枚空转的消费分支 | 本程给它接上生产者，见 §4 |
| `TestNoGitSwitchCapabilityInThePanelSurface` | panel | PASS | `gitSwitchCapabilityRe`（`composer_test.go:268`）扫 `internal/panel/**` 非测试 `.go`，词面含 `\bworktree\b`／`\bgit\.branch\b`／`\bgit\.repo\b` | **措辞雷**：本程产码只写 tree／linked tree／复数形 |
| `TestGitDimensionHasNoModelCallableTool` | panel | PASS | D34 表无 git 行＋`internal/tools` 无 `"git.*"` 工具名＋`bridge.go` 源码 `panel.*` 恰好四枚 | 本程零枚入向方法，AC#3 反向判据不许被我打红 |
| `TestWorkspaceSwitchNarrowsWhatTheAssessorJudges` | tools | PASS | 真 `risk.RiskAssessor` × 真 canonicalizer 的判定变化；`:56` 现读 `canon.WorkspaceRoot() != ""` | **直接被本程改签名** |
| `TestWorkspaceSwitchRefusesOutOfScopeAndMissingPaths` | tools | PASS | 三枚拒绝都不留 workspace；`:124/:129` 直接调 `SetWorkspaceRoot(<裸串>)` | 签名变更 ⇒ 调用点跟改（语义不变） |
| `TestWorkspaceSwitchRefusesAnExpandedSpelling` | tools | PASS | `%TEMP%`／`~` 拼法必须撞 `ErrRewrittenPath`，且不留 workspace（可移植，本程实测跑到） | 本程的"真改写账户"取数源（§4-2） |
| `TestWorkspaceSwitchRefusesAJunctionToOutside` | tools | PASS（未 skip） | 真 junction＋C26 拒绝传导 | 未动；登记"本机 `mklink /J` 可用" |
| `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` | risk | （单包跑，见 §6） | **全仓文本尺**：任何非测试 `.go` 读 `.Canonical` 而该文件不含 `Actable(`/`.Rewritten` 即红（按文件判、不剥注释） | 本程产码三处 `.Canonical` 读点全部与账户同文件（§6 复跑凭据） |
| `cmd/wisp/panel_pump_test.go:342 TestSnapshotWorkspaceSectionReportsTheNarrowing` | cmd | （他程地界，只读） | `:383` 用 `strings.EqualFold` 比 `packet canonical` 与 `handler canonical`；`:388` 两者**逐字相等时只 `t.Log` 不 `t.Errorf`**；`:391` 断言 `!Reparse && !Rewritten` | 本程把 fold-key 小写串换成 C26 原拼法后该枚**不红**；`!Rewritten` 仍成立（两道上游拦一字未松） |

---

## §2 选形与否形

票面 AC#7 自己列了两枚候选形：ⓐ `WorkspaceRoot()` 连账户一起回；ⓑ `WorkspaceViewFromRoot` 收一枚带账户的结构。

**现量先把边界钉住**：全仓（`.go`，剥掉 `.scratch` 副本）里 `WorkspaceViewFromRoot` 的**生产调用方只有一枚**，`cmd/wisp/panel_pump.go:87`：

```go
return panel.WorkspaceViewFromRoot(rt.paths.WorkspaceRoot())
```

而 `cmd/wisp/**` 是本程 ⛔ 写面外（另一枚腿在飞，起手 porcelain 已有 `M cmd/wisp/firstrun.go`）。⇒ **选形的第一判据是"这一枚调用点必须一字不改仍然成立"**，否则本程要么越界写 `cmd/`，要么把仓留在 `cmd/wisp` 编译不过的状态（两者都比不修更坏）。

**选了：ⓐ＋ⓑ 合成一形，账户的载体是 `risk.Result` 本身。**

* `internal/tools/paths_workspace.go`：`WorkspaceRoot() risk.Result` —— 收窄仍在（`p.workspace` 仍存 fold 键，`InAllowlist` 判定字节不变），但**回答的是 C26 对这条收窄说过什么**：`Canonical`＝在用坐标，`Spelling`＝使用者原串，`Reparse`／`Rewritten`／`Rewrites`＝票 102 的账。未收窄＝零值 `risk.Result`（`Canonical==""`）。
* `internal/tools/paths_workspace.go`：`SetWorkspaceRoot(root string, acc risk.Result) error` —— 账户**只能是描述这枚 root 的那枚**（`foldPath(root)==foldPath(acc.Canonical)`，否则拒绝），并把账户与收窄写进同一次加锁。这把"随手拿一枚没解析过的串设 workspace"的老口子收窄了（原先只查 `inRoots`）。
* `internal/panel/workspace.go`：`WorkspaceViewFromRoot(res risk.Result) WorkspaceView` —— 收账户；`PathScope` 接口的 `WorkspaceRoot()` 返回类型跟着换（`*tools.PathCanonicalizer` 结构上继续满足它）。
* 于是 `panel_pump.go:87` **逐字节不变**仍然成立：左边返回 `risk.Result`，右边收 `risk.Result`。

**否掉的两形，具名理由：**

1. **ⓐ 的字面形"连账户一起回"＝多返回值 `WorkspaceRoot() (string, risk.Result)`**：`panel_pump.go:87` 立刻变成"双值单用"编译错误。要么我改 `cmd/wisp`（越界），要么留一枚编译不过的仓（更坏）。⇒ 否。
2. **ⓑ 的字面形"收一枚带账户的结构"＝panel 自己造一枚 `WorkspaceAccount` 结构、`WorkspaceRoot()` 继续回裸串**：(a) 同样要 `panel_pump.go:87` 去**构造**那枚结构才喂得进去，还是越界；(b) 它把票 102 的账户词汇**复制成第二枚类型**（`Spelling`/`Reparse`/`Rewritten`/`Rewrites` 四枚字段抄一遍），正是 `internal/panel/pump.go:485-493` 与 `internal/tools/subagent_197.go:42-49` 为 streamkey 记过的同一枚病（"两枚拷贝＝页面开始对不上账"）；(c) `internal/panel` 与 `internal/tools` **不许互相 import**（`subagent_roster_197.go:48` 逐字 "not import internal/tools, so the composition root does the mapping"），所以 panel 侧自造结构的话，tools 侧回给它的类型仍得由第三包提供。⇒ 否。
3. **顺带排除的第三形（没摆给编排者，直接具名）**：把账户类型放进 `internal/risk` 或新建叶子包——两枚都在本程写面外，且 `risk` 里已经有这枚结构（`risk.Result`），再造＝把 §2-2(b) 的病搬近一步。

**为什么 `risk.Result` 是合法的载体而不是越界**：它已经在同一枚接缝上过车——`PathScope.ResolveWorkspace` 的返回类型今天就是 `risk.Result`（`internal/panel/workspace.go:42`）。本程只是把**同一枚接缝的另一条腿**也换成它，`internal/risk/**` **一字未改**（判级语义、`Actable()`、`ErrRewrittenPath` 全部零字节，终态凭据 `git status --porcelain -- internal/risk/` 空，见 §6）。

**这道选择没有动安全语义**：`RequestWorkspaceSwitch` 仍先在 `workspace.go:85` 读 `Actable()`，`tools` 侧 `ResolveWorkspace` 仍在 `:64-66` 拦一次；本程**没有为了写真值去绕任何一道**。改完之后 `Rewritten==true` 的可达路只有一条：宿主把**展开后的那棵树按名字**交进来（`ErrRewrittenPath` 自己的措辞 "act on the expanded tree only if you asked for it by name"），此时账户与真值一起进快照，`181-r2` 与票 200-r2 那两枚今天走不到的消费分支就活了（§4-2/§4-3）。

---

## §3 逐枚 before-after

### 3-1 产码（六枚，全部在写面内；`internal/risk/**` 零字节）

| 件 | before | after | 删除列 |
|---|---|---|---|
| `internal/tools/paths_workspace.go` | `WorkspaceRoot() string` 回 `p.workspace`（fold 键）；`SetWorkspaceRoot(root string) error` 只查 `inRoots`；`ClearWorkspace()` 只清 `p.workspace` | `WorkspaceRoot() risk.Result` 回**记录的那本账**（返回的是副本，`Rewrites` 不外泄 backing array）；`SetWorkspaceRoot(root string, acc risk.Result) error` 同锁写收窄＋账户，并**绑定**（`foldPath(acc.Canonical)==foldPath(root)`，否则拒）；`ClearWorkspace()` 两枚一起清 | 8（三枚函数体与注释改写，语义只增不减） |
| `internal/tools/paths.go` | `PathCanonicalizer` 无账户字段 | 新增 `workspaceAccount risk.Result`（与 `workspace` 同锁、同生同灭），头注写明"这份账本来就有，只是被扔掉" | 0 |
| `internal/panel/workspace.go` | `PathScope.WorkspaceRoot() string`／`SetWorkspaceRoot(string) error`；`WorkspaceViewFromRoot(root string)` 只填 `Set`／`Canonical`／`Reason`；`currentAfter(scope, previous string)` 用 `!=` | 接缝两端都换成 `risk.Result`；`WorkspaceViewFromRoot(res risk.Result)` 填 `Set`／`Spelling`／`Canonical`／`Reparse`／`Rewritten`／`Reason`，改写与穿越联接点各有其句；`currentAfter` 改比"坐标＋账户"（`risk.Result` 含切片不可 `==`），新增 `sameC26Answer`／`accountSummary`；两枚审计行的 `from=%q` 改喂 `previous.Canonical`（保持逐字形状，见 §9 第 2 条） | 16 |
| `internal/panel/git.go` | `GitViewRewritten` 头注："MEASURED TODAY：这分支在当前树上不可达" | 同一处措辞按新现量改写：账户现在是**读数**、面板发起那扇门仍然产不出 `rewritten=true`（两道 `Actable()` ＋新的绑定检查），可达形＝宿主按名字交出展开树，并指名 `TestWorkspaceAccountReachesEveryConsumer` | 7 |
| `internal/panel/instructions_200.go` | `ProjectInstructionLoadRequestFor` 头注同样写着"没有任何生产者交出 Rewritten=true 的 view" | 同上改写：生产侧已填、面板门仍拦得住、`TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll` 与新增那枚并列 | 9 |
| `cmd/wisp/panel_pump.go` | `:87 return panel.WorkspaceViewFromRoot(rt.paths.WorkspaceRoot())` | **一字未改**（左值型与右参型同时换成 `risk.Result` 才成立）；`git status --porcelain -- cmd` 里本程零枚（在飞的 `firstrun*.go` 是别程的） | 0 |

### 3-2 既有判据（六枚跟签名走，断言**语义**不变）

| 用例 | before 的那一行 | after 的那一行 | 语义 |
|---|---|---|---|
| panel `TestWorkspaceViewRenderedForSnapshots` | `WorkspaceViewFromRoot("")`／`WorkspaceViewFromRoot(`d:\work\a`)` | `WorkspaceViewFromRoot(risk.Result{})`／`risk.Result{Canonical: `d:\work\a`, Spelling: `d:\work\a`}` | 未设＝`Set=false` 且有 reason、已设＝`Canonical` 逐字——两枚断言原样；另**补**一枚空白 canonical 的子断言（§9 第 6 条） |
| panel `TestWorkspaceSwitchKeepsThePreviousScopeOnFailure` | `scriptedScope{root: `d:\work\alpha`}` | `root: risk.Result{Canonical: `d:\work\alpha`}` | 断言 `view.Canonical`／`view.Set` 原样 |
| panel `TestWorkspaceSwitchDetectsAScopeThatMovedAnyway` | 同上 | 同上 | 断言原样 |
| panel `TestWorkspaceSwitchAuditsAndAppliesACleanPath` | 未改 | 未改 | `from=""`／`spelling=`／`result=ok` 逐字仍然成立（靠 `previous.Canonical` 那处修正） |
| tools `TestWorkspaceSwitchNarrowsWhatTheAssessorJudges` | `canon.WorkspaceRoot() != ""`／`SetWorkspaceRoot(res.Canonical)` | `.Canonical != ""`／`SetWorkspaceRoot(res.Canonical, res)` | 判定变化（L1→L2、R2 命中）一枚没动 |
| tools `TestWorkspaceSwitchRefusesOutOfScopeAndMissingPaths`／`…AnExpandedSpelling`／`…AJunctionToOutside` | `WorkspaceRoot() != ""` 三处、`SetWorkspaceRoot(x)` 两处 | `.Canonical` 三处、`SetWorkspaceRoot(x, risk.Result{Canonical: x})` 两处 | 三枚拒绝都不留 workspace——原断言逐字，只是读的是账户的坐标那半 |
| tools `task_pointer_authority_ac3_174r4_windows_test.go:196` | `paths.SetWorkspaceRoot(res.Canonical)` | `paths.SetWorkspaceRoot(res.Canonical, res)` | 票 174 那枚 work-page 判定的断言一字未动（只补参数） |

### 3-3 新增判据（七枚）

tools `internal/tools/paths_workspace_account_181r3_test.go`（四枚）：

1. `TestWorkspaceRootAnswersWithTheAccountThatProducedTheNarrowing`——真 canonicalizer×真目录，调用方用 `filepath.ToSlash` 拼法送进去，要求 `WorkspaceRoot().Spelling` **逐字**回那串（今天不可能：裸串没有拼法可回），`Canonical` 等于 `Actable()` 的答，并复验 `InAllowlist` 没因记录账户而放宽（workspace 外的树仍 false、内的仍 true）。
2. `TestWorkspaceRootReportsARealRewrittenAccountRecordedByC26`——两子例：①`%WISP181R3ACC%` 仍撞 `risk.ErrRewrittenPath` 且不留 workspace（**这一枚是"本程不是放宽"的正向凭据**）；②同一枚 `res`（C26 真返回的 `Rewritten=true`／`Rewrites=["env"]`）被宿主按名字交出展开树后记录，回读必须 `Rewritten=true`＋`Spelling` 是被替换的那串＋`Rewrites` 含 `env`；末尾验返回的是副本（改返回值不能污染账本）。
3. `TestSetWorkspaceRootRefusesAnAccountAboutAnotherTree`——两枚树都在授权内，`inRoots` 抓不到；把 beta 的账贴到 alpha 上必须被拒且不留 workspace。
4. `TestClearWorkspaceDropsTheAccountWithTheRoot`——`reflect.DeepEqual(got, risk.Result{})`：账不能活得比收窄久。

panel `internal/panel/workspace_account_181r3_test.go`（三枚）：

5. `TestWorkspaceViewCarriesEveryHalfOfTheAccount`——四子例（clean／rewritten／reparse／unset）：`Rewritten`、`Spelling`、`Reparse` 必须各就各位，rewritten 那枚还要求 reason 里同时出现「改写」「env」「原拼法」「展开后的树」四样；unset 那枚要求与 `UnsetWorkspaceView()` 逐字同形。
6. `TestWorkspaceAccountReachesEveryConsumer`——用与 `cmd/wisp/panel_pump.go:87` **逐字同形**的表达式 `WorkspaceViewFromRoot(scope.WorkspaceRoot())` 造 view，然后三子例：`ReadGitForWorkspace` 的 reason 必须含 `rewritten=true` 且两枚坐标都在、分支仍报出（不吞维）；`ProjectInstructionLoadRequestFor` 必须两棵目录都不给（`Refused` 非空）；过一遍 `pump.Marshal()`，`composer.workspace.rewritten` 在 JSON 上必须是 `true`、`spelling`／`canonical` 两枚坐标逐字、reason 含「改写」。
7. `TestScopeThatChangesItsAccountIsSurfaced`——一枚坐标不变而账户偷变的 scope：必须报错、`Set=false`、reason 里出现「账户」与 `rewritten=true`（这是只有账户上了路才可能存在的新病形）。

---

## §4 判据为什么能区分"真填了账户"与"字段恒 false"

派单把话说死：⛔ 编译通过不算证明、⛔ 拿"这枚字段存在"当凭据。三枚尺各自的位置：

1. **账户的来源是真的**（tools 那两枚）。§3-3 的 1／2 都用 `NewPathCanonicalizer` ＋ `t.TempDir()` ＋ `t.Setenv` 构造，`Rewritten=true`／`Rewrites=["env"]` 不是测试手写的字面量，而是 `risk.Resolve` 的 `%VAR%` 展开步**自己返回的**——同一取数形与既有 `TestWorkspaceSwitchRefusesAnExpandedSpelling` 一致，而那枚在本机是 PASS（不是 SKIP，§1 已量：两包 `--- SKIP` 均为 0）。⇒ "字段恒 false"那种生产者在这两枚尺下立刻红：它回答不出 `Spelling`，也回答不出 `Rewritten`。
2. **账户的去处是通的**（panel 第 6 枚）。判据不看"字段存在"，看**下游行为**：git 那一维的 reason 是否把账户说出口（`181-r2` 那条今天走不到的分支）、说明加载器是否两棵树都不读（票 200-r2 那条）、以及快照 JSON 里 `rewritten` 是否为 `true`。三枚都是"未修码上今天不可能响"的形状：未修码上 `WorkspaceViewFromRoot` 只吃裸串，账户无从进入，三条消费全哑。
3. **突变就是那把区分尺**（§5）。M1 把生产者退回"不填"（`WorkspaceRoot()` 只回坐标），M2 把渲染退回三枚字段。两枚突变各自让**指名的**用例红，且红句逐字点名"the producer dropped the book again / 181-r2's consumer branch is dead code"——这正是"恒 false"这一状态被抓住时的句子。
4. **诚实的边界（不夸大）**：面板发起收窄那扇门**仍然不可能**产出 `rewritten=true`，因为 `paths_workspace.go` 与 `workspace.go` 两道 `Actable()` 一字未松，本程还把第三道（账户—坐标绑定）加严了。所以本交付证明的是**"生产者会说话、字段是读数"**，不是"运行中的进程今天会说改写"。要把后者做成可达，唯一办法是松 `Actable()`＝改安全语义＝派单禁区＝停手上报（§7 第 2 条）。这条边界具名写着，比让它含糊过去值钱。

---

## §5 突变＋红句逐字＋还原凭据

突变全部用 `go test -overlay`（跟踪件一字未动；副本落在仓外 `C:/Users/swq/AppData/Local/Temp/181r3mut/`，⛔ 不在仓内建 `.go`，避免 267-r2 自首过的那枚"模块内游离 `.go` 被 `go list` 当包"污染）。

**M1｜生产者退回不填**（`internal/tools/paths_workspace.go` 的 `WorkspaceRoot()` 改成 `return risk.Result{Canonical: p.workspace}`）：

```
--- FAIL: TestWorkspaceRootAnswersWithTheAccountThatProducedTheNarrowing (0.02s)
    paths_workspace_account_181r3_test.go:65: WorkspaceRoot().Spelling = "", want the caller's own string "C:/Users/.../alpha" verbatim - a producer that drops the spelling drops ticket 102's traceability: a verdict taken on Canonical could no longer be traced back to the form it was asked about
--- FAIL: TestWorkspaceRootReportsARealRewrittenAccountRecordedByC26 (0.01s)
    --- FAIL: TestWorkspaceRootReportsARealRewrittenAccountRecordedByC26/a_recorded_account_comes_back_as_written (0.00s)
    paths_workspace_account_181r3_test.go:133: rewritten = false after recording C26's own rewritten Result ({Canonical:c:\users\... Reparse:false Resolved:false Spelling: Rewritten:false Rewrites:[]}): the producer is filling the field with a constant again, and 181-r2's consumer branch is dead code
    paths_workspace_account_181r3_test.go:137: spelling = "", want the substituted spelling "%WISP181R3ACC%"
```

包级读数：`FAIL github.com/CarlosShao/wisp/internal/tools 12.392s`，红名册＝上面两枚，**其余 194 枚不变色**；且 `the_refusal_leg_is_untouched` 那枚子例在 M1 下仍然 PASS ⇒ 突变只打死生产者，打死不了票 102 的拒绝腿（两件事没被我焊在一起）。

**M2｜渲染退回三枚字段**（`internal/panel/workspace.go` 的 `WorkspaceViewFromRoot` 只填 `Set`／`Canonical`／`Reason`）：

```
--- FAIL: TestWorkspaceViewCarriesEveryHalfOfTheAccount (0.01s)
    --- FAIL: .../clean_answer_is_the_ordinary_road (0.00s)   workspace_account_181r3_test.go:63: spelling = "", want the caller's own string "D:\\work\\Wisp" verbatim
    --- FAIL: .../rewritten_answer_must_be_rewritten_in_the_view (0.00s)  :83: rewritten = false for the account C26 marked rewritten ({Set:true Spelling: Canonical:C:\elsewhere\proj Rewritten:false Reason:已收窄到该工作区，但这条路径被 C26 的展开步改写过（env）…}): the producer dropped the book again, and every consumer of this field is decoration
    --- FAIL: .../reparse_answer_is_read_not_assumed (0.01s)  :104: reparse = false although C26 reported an exception-listed traversal
--- FAIL: TestWorkspaceAccountReachesEveryConsumer (0.10s)
    workspace_account_181r3_test.go:150: the view built from a scope that reports a rewritten account is not rewritten: {Set:true … Rewritten:false …}
```

M2 那一行还顺手把"半填"的形状量出来了：`Reason` 里有改写句、`Rewritten` 却是 `false`——如果只靠文案判据会放过它，正是判据 5 与判据 6 要分开写的原因。

**M3｜一致性检查退回只比坐标**（`currentAfter` 改回 `previous.Canonical != now.Canonical`）：

```
--- FAIL: TestScopeThatChangesItsAccountIsSurfaced (0.00s)
    workspace_account_181r3_test.go:226: the moved scope was rendered as a normal narrowing: {Set:true Spelling:%WISP181R3ACC% Canonical:D:\work\alpha Rewritten:true …}
    workspace_account_181r3_test.go:229: reason "已收窄到该工作区，但这条路径被 C26 的展开步改写过（env）…" must name the account half of what moved
```

其余 `TestWorkspaceViewCarriesEveryHalfOfTheAccount`／`TestWorkspaceSwitch*` 在 M3 下全 PASS ⇒ M3 只打死"账户偷变"这一形，没有连带。

**还原凭据**（三枚突变跑完后同发取，`sha256sum -c`）：

```
98e2e6730ebb0369d53da79d0ea9beaeb94e2a62ba15681b8f886ae9830db8d6  internal/tools/paths_workspace.go  OK
0330e416a3b256efc907be17ba8246b824120e49b34c47b2403487f340a45e44  internal/panel/workspace.go        OK
```

两枚哈希取在突变**之前**（`/tmp/181r3mut/hashes-before.txt`），`-c` 复验在突变之后 ⇒ 跟踪件字节未变；`git status --porcelain -- internal cmd` 终态只剩别程的 `cmd/wisp/firstrun*.go`（本程 `internal/` 在 16901acb 之后为空）。还原＝**无需还原**（overlay 不改盘），凭据＝上面两枚哈希＋porcelain。

---

## §6 门禁读数（带时刻，均为 2026-10-05，+08）

| 尺 | 时刻 | 读数 |
|---|---|---|
| `go test -count=1 -v ./internal/panel/ ./internal/tools/`（改后） | 11:33:38 | panel `FAIL 1.697s`，顶层 PASS **118**、FAIL **4**（在册那四枚，逐字见 §1）、SKIP 0；tools `ok 14.952s`，PASS **198**、FAIL 0、SKIP 0 |
| 名册作差（before→after） | 11:33:38 | `diff` 只有 **7 行 `> PASS`（新增）**，`<` 行数 **0**、变色 **0**；panel 115→118、tools 194→198 |
| `go test -count=1 ./internal/risk/` **单包**（定式 A380：动"消费解析路径"这类面必跑） | 11:33:56 | `ok github.com/CarlosShao/wisp/internal/risk 4.414s` ⇒ 全仓 `.Canonical` 尺绿，本程没有新增"只看坐标不看账"的消费者 |
| `git status --porcelain -- internal/risk/` | 11:33:56 | **空**（禁区零字节，含测试件） |
| `go vet ./internal/panel/ ./internal/tools/` | 11:29 | rc=0，无输出（第一次跑抓到 tools 测试件签名未跟改，见 §9 第 1 条） |
| `go vet ./cmd/wisp/` | 11:29:56 | rc=0 ⇒ 那枚唯一生产调用点与 `cmd/wisp` 全部测试件在新签名下**仍然编译**（只证编译，⛔ 未跑那两枚包，见 §7 第 1 条） |
| `go build ./cmd/...` | 11:29 | rc=0 |
| `sh scripts/d22scan.sh` | 11:31（产码定稿后）／终态复跑见本节末 | rc=0，`d22scan: clean - no D22 ban violations`；`ban #8 internal/` examined **510**、`bans #1-5 internal/` **228**、`cmd/` **38**、`ban #7 internal/tools/` **23** |
| `sh scripts/check-path-length-budget.sh` | 11:39:46 | rc=0；`denominator: tracked paths=5813 over-budget=57 covered by roster=57 not in roster=0`，`VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`；longest=180 是既有那枚票 252 工单名，与本程无关（本程新增三枚件名最长 46 字符） |
| `/d/work/base/gopath/bin/gofumpt.exe -l`（十一枚点名） | 11:35 | 第一遍抓到 `internal/panel/workspace_account_181r3_test.go` 未格式化 ⇒ `-w` 后复跑 **空**；`gofumpt -version` = `v0.12.0 (go1.27.1)`，`go env GOPATH` = `D:\work\base\gopath`（⛔ 没用 `~/GOPATH/bin` 那把瞎尺） |
| `perl` 自扫 emoji（U+1F000-1FAFF／2200-22FF／2600-27BF／2B00-2BFF／FE0F，七枚本程件） | 11:40 | **0 命中** |
| `grep -c '"panel\.' internal/panel/bridge.go` | 11:40 | **5 行**，其中 `:35` 是注释字样，方法常量恰好 **4 枚**（`:42-45`）⇒ 与 181-r1／r2 名册一致，本程零枚新增；源码解析尺 `TestGitDimensionHasNoModelCallableTool` 本程复跑 PASS |
| `grep -rn "internal/panel" internal/tools/*.go`（非测试） | 11:40 | **空** ⇒ tools 不 import panel；反向 `internal/tools` 也未被 panel 的非测试件 import（§2 的架构约束守住了） |

⚠ 本节里 d22scan／path-length／gofumpt 三枚**取在产码 commit `16901acb` 之后**；证据件本笔 commit 之后不再动产码，故不再复跑产码门禁；最后一笔只改本件时会在交件回报里附**同一把尺的终态复跑**读数。

---

## §7 判不动或量不到

1. **`cmd/wisp` 那枚相关用例只核读、没实测**：`TestSnapshotWorkspaceSectionReportsTheNarrowing`（`cmd/wisp/panel_pump_test.go:342`）。理由＝写面外＋派单明令不许跑那两枚包整包（另一枚腿在飞）。我做的是读它的断言原文：`:383` 用 `strings.EqualFold` 比两枚坐标、`:388` 在逐字相等时只 `t.Log`、`:391` 要 `!Reparse && !Rewritten`。⇒ 本程把 packet 的 `workspace.canonical` 从 fold 小写键换成 C26 原拼法**不会**让它红（只会让它那行 `t.Log` 出现），`!Rewritten` 也仍成立。**这一判语是核读出来的，不是跑出来的**，请由能跑 `cmd/wisp` 的那一枚收。
2. **面板发起的收窄这条路 `rewritten` 仍不可达**，且不是本程能做满的：要让它可达只能松 `risk.Result.Actable()` 那道拦（`internal/risk/pathresolver.go:93-99`）＝改安全语义＝派单禁区。⇒ 停手上报，交件把"生产者会说话"做满，把"今天有真值可说"停在票 102 的裁定上。这也意味着 AC#7 里"三档处置今天完全等价"那格（`git.go` 的 `GitViewRewritten` 选①还是选②）**继续归编排者裁**，本程未动那一行。
3. **在册红四枚，比派单列的多一枚**（具名上报口径差，未修未追）：`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`（同一因：`design/assets/tokens.css` 被别程删了未 staged，本程连 `design/**` 都没读）＋`TestComposerContractTypesMatchFrontend`＋**新见一枚**`TestApprovalCardViewJSONKeysMatchFrontendTypes`（红句逐字 `Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`，同族病＝`frontend/src/lib/panel.ts` 未声明，归能碰 `frontend/**` 的那一程）。起手名册里就有它，本程未使其变色（before/after 同一枚红，见 §6 名册作差）。
4. **量不到：本仓真数据下的 packet 未取**。`wisp run` 全链跑一次会撞 `cmd/wisp`（禁跑）与票 98 那条"本机测不到东西"。§3-3 第 6 枚是**在 panel 包内用真 pump＋真 Marshal＋真文件系统临时树**做的等价证明，但它顶不了"真进程跑一遍"的位；这一格请由编排者或 cmd 那支补。
5. **账户伪造面未做仪器**：`SetWorkspaceRoot(root, acc)` 只能验账户**描述**的是这枚 root（绑定检查），验不了"这枚 root 确实是被解析出来的"。今天上游两道 `Actable()` 担着这件事；要钉死得让 tools 自己重解析，那会把 `Rewritten` 永远焊成 false（§4-4 的反面）。具名留档，不冒充做满。

---

## §8 commit 名册

| 笔 | 短哈希 | pathspec（枚枚点名） | 内容 |
|---|---|---|---|
| 1 | `b672f853` | `.scratch/wisp/probes/181/r3/evidence.md` | 本件骨架＋§0／§1／§2 填实 |
| 2 | `16901acb` | 十一枚：`internal/panel/{workspace.go,workspace_test.go,workspace_account_181r3_test.go,git.go,git_test.go,instructions_200.go}`＋`internal/tools/{paths.go,paths_workspace.go,paths_workspace_test.go,paths_workspace_account_181r3_test.go,task_pointer_authority_ac3_174r4_windows_test.go}` | 产码六枚＋判据（4 枚既有跟签名、7 枚新增），`+653/-62` |
| 3 | 本笔（提交时刻见 git log） | `.scratch/wisp/probes/181/r3/evidence.md` | §3-§10 填实 |
| 4 | 终态复跑笔（如需要） | 同 3 | 只在门禁读数变化时补 |

台件（本程建、只建不删）：`/tmp/181r3-baseline.txt`、`/tmp/181r3-after.txt`、`/tmp/181r3mut/{paths_workspace_M1.go,workspace_M2.go,workspace_M3.go,m1.json,m2.json,m3.json,m1-out.txt,m2-out.txt,m3-out.txt,hashes-before.txt}`——全部落在**仓外**临时目录，仓内未新增游离 `.go`。⚠ 记一笔纪律偏差：本程把突变副本放在仓外 `/tmp` 而不是 `.scratch/wisp/probes/181/r3/`，理由是 267-r2 自首过"模块内游离 `.go` 会被 `go list` 当包、`go vet ./...` 咬人"；代价是台件不随 commit 入库，故本 §8 把名册写全，逐字红句已抄进 §5。

---

## §9 自我对抗（我用什么尺怀疑了自己每一句）

1. **先跑包再动笔**抓到"我以为的签名半径"是错的：第一次 `go vet` 立刻爆 `internal/tools/task_pointer_authority_ac3_174r4_windows_test.go:196`——票 174 那枚 work-page 测试也在调 `SetWorkspaceRoot`，只按 §1 名册（`paths_workspace_test.go`）改是不够的。这条不在我预期里，是尺给的。
2. **审计行的 `%q` 打在结构体上**是我自己写出来的第二号错：`go test` 把 `TestWorkspaceSwitchAuditsAndAppliesACleanPath` 打红并倒出 `from={"" %!q(bool=false) …}`。修法不是改断言，而是把喂进去的值换成 `previous.Canonical`——成功行与拒绝行两处都得换，我第一遍只改了拒绝行，是靠整包跑而不是单跑抓到的。
3. **`reason` 双重量化**：`accountSummary` 内部已有 `%q`，外层再 `%q` 会把整句套引号并转义。改成外层 `%s`；`TestScopeThatChangesItsAccountIsSurfaced` 的断言因此只看子串，不受引号形状牵制。
4. **我一开始给"premise 不成立"写了 `t.Skip`**（tools §3-3 第 2 枚的②子例）。派单明写 ⛔ 把 SKIP 读成通过 ⇒ 换成 `t.Fatalf` 并写清"premise broken"，本机实测该枚走到 `Rewritten=true` 那一支（§6 的 SKIP 全 0 也是这把尺）。
5. **`WorkspaceView` 能否用 `==`**：它六枚字段全是标量，可比较 ⇒ unset 那枚用 `!=` 逐字比；`risk.Result` 含切片不可比较 ⇒ 一处用 `reflect.DeepEqual`（`TestClearWorkspaceDropsTheAccountWithTheRoot`），一处把"只比坐标"写成会误事的反例并让它成为 M3 的靶。
6. **空白 canonical**：旧代码是 `strings.TrimSpace(root)==""`，我若只判 `res.Canonical==""` 就悄悄改了语义 ⇒ 补一条 `risk.Result{Canonical:"   "}` 的子断言，钉"空白＝未设"这条没被换掉。
7. **我没有为了写真值去动任何一道拦**：M1 的对照实验里 `the_refusal_leg_is_untouched` 仍 PASS，§3-3 第 2 枚①子例也专门跑 `%VAR%`→`ErrRewrittenPath`。这三处共同构成"不是放宽"的凭据，而不是我的一句话。
8. **措辞雷复尺**：动笔前重读 `composer_test.go:266-272` 的 `gitSwitchCapabilityRe` 原文（含 `\bworktree\b`／`\bgit\.branch\b`／`\bgit\.repo\b`），产码与注释只写 tree／linked tree／复数形；`TestNoGitSwitchCapabilityInThePanelSurface` 在 §6 的 after 全包里 PASS（不是我以为，是跑出来的）。
9. **"两形都摆"的诱惑**：我一开始想把 ⓐ／ⓑ 写成"两枚都可"，被派单明令否掉；§2 因此把"唯一那枚生产调用点"做成第一判据，并给了否形理由，不是给编排者出选择题。
10. **架构约束我核了两把**：`grep` 证 tools 不 import panel、panel 非测试件不 import tools；`bridge.go` 的 `panel.*` 常量枚数用源码文本尺（不是我数的）复跑 PASS。
11. **空读数先证命令真跑到**：`gofumpt`／`d22scan.sh`／`check-path-length-budget.sh` 三把都附了 rc、版本或 `examined=`／`denominator=` 这类自证行；`~/GOPATH/bin/gofumpt` 那把瞎尺没用。
12. **我承认量不到的**：`cmd/wisp` 那枚相关用例没跑（§7-1）、真进程 packet 没取（§7-4）、账户伪造面没做仪器（§7-5）。这三条写在这儿而不是写进 §10 的"完成"里。

---

## §10 交件判语

- **AC#7 这一格问的东西到手了**：`panel.WorkspaceView.Rewritten`（连同 `Spelling`／`Reparse`）现在有生产者，生产者读的是 `internal/tools` 记录的 C26 账本，账本随收窄同生同灭并绑定坐标；`181-r2` 的 `GitViewRewritten` 与票 200-r2 的 `ProjectInstructionLoadRequestFor` 那两条空转消费分支，由 §3-3 第 6 枚把账户一路送到 git 的 reason、说明加载器的拒绝与快照 JSON 的 `"rewritten":true`。
- **"恒 false"这一状态不是被我否认的，是被突变证实可抓的**：M1／M2 各自让指名用例红，红句逐字在 §5。
- **没做的事也具名**：没松任何一道 `Actable()`、没碰 `internal/risk/**`、没加 `panel.*` 方法、没读 `frontend/**`／`design/**`、没跑 `cmd/wisp`／`internal/ball`、没翻票面 AC 框、没 push。
- **给编排者的一句话**：本程唯一可能让人犹豫的是"packet 里 `workspace.canonical` 从 fold 小写键换成了 C26 原拼法"。这是账户随行的必然结果（同一枚 `Result` 的 `Canonical` 字段就是 C26 的答案），并且把 `docs/evidence/s1/35-panel-snapshot-pump-r1.md` §3 那枚"fold-key  finding"顺手消掉了；`cmd/wisp` 的相关断言是 EqualFold＋t.Log，核读不红（§7-1 具名"未实测"）。如果编排者要的是"坐标不变、只补账户"，改点是一行（`WorkspaceViewFromRoot` 里把 `Canonical` 换回 `foldPath` 形），我不建议，那等于把账户的一半意义扔掉。
