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

（本节在产码落盘后逐枚填：改动前后每枚用例的名字、断言原文、读数。）

## §4 判据为什么能区分"真填了账户"与"字段恒 false"

（本节填判据清单与每枚判据的"未修码上响不响"论证。）

## §5 突变＋红句逐字＋还原凭据

（本节填突变做法、逐字红句、还原前后哈希。）

## §6 门禁读数（带时刻）

（本节填 `go test` 两包／`internal/risk` 单包／`sh scripts/d22scan.sh`／`scripts/check-path-length-budget.sh`／`go vet`／`gofumpt -l` 的读数与取数时刻。）

## §7 判不动或量不到

（本节填归别人的红、量不到的半径、以及本程判不动的格。）

## §8 commit 名册

（本节逐枚填 commit 短哈希与 pathspec。）

## §9 自我对抗

（本节填"我用什么尺怀疑了自己交付的每一句"。）

## §10 交件判语

（本节在终态取数后填。）
