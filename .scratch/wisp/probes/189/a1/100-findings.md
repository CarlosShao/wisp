# 189-a1（10-08 程）交件 — 四问逐格＋最少落点表＋雷区复核

> 只读腿。锚 `914177e6`（现量，见 `00-anchor.md`）。零 go 命令、零产码写点。
> 配套尺件（每把尺原文＋rc 自读）：`10-source-rulers.txt`／`20-carrier-and-contract-rulers.txt`／
> `30-gate-and-frontend-rulers.txt`／`40-contract-pin.txt`／`50-existing-nails.txt`／
> `60-gate-shape.txt`／`70-minefield-and-keys.txt`／`80-premise-recheck.txt`／`90-red-roster-and-workspace.txt`。
> **先读这一句**：本票的"落点表"在 09-28 已由同名腿交过一次（`docs/evidence/s1/189-uncommitted-count-design-a1.md`，
> 编排者 09-28 18:1x 照准，账 `A395`）。本件的增量＝**复核那件的前提被十月产码改动弄过期了哪几处**＋
> **把"最少落点"重算成今天能对上的 file:line**＋雷区现量。**"有没有源"一问本程按查重纪律只引用不重跑**（见 §2.0）。

---

## 1. 问①：工单 189 逐格框态＋票面前提复核

### 1.1 框数（尺见 `00-anchor.md`）

| 尺原文 | 命中 | rc |
|---|---|---|
| `grep -cE '^[[:space:]]*- \[ \]' $F` | **6** | 0 |
| `grep -cE '^[[:space:]]*- \[x\]' $F` | **0** | 1 |
| `grep -c '\[x\]' $F`（松口径，不限行首） | **0** | 1 |
| `grep -oE '^[[:space:]]*- \[[^ ]\]' $F \| sort \| uniq -c`（其它框字符） | 无输出 | 0 |

⇒ **6 未勾／0 已勾**（AC#1–AC#6 全开），与 `.scratch/wisp/probes/182/c1/100-census.md` K6 那格"六框全未勾"合一，**无冲突**。

### 1.2 票面前提逐句复核（**票体只有 26 行，未改一字；以下全是复核结论**）

| # | 票面原句（出处 `189-….md`） | 复核结论 | 现量凭据 |
|---|---|---|---|
| 1 | 行 3"与票 181 同一条读面（复用，别重造）" | **仍成立**，`internal/panel/git.go` 就是那一枚读面 | `10-source-rulers.txt` R1/R2/R5：git.go **569 行**、族名册只有 `git.go`+`git_test.go` |
| 2 | 行 8 `grep -rniE "git diff\|numstat"` ⇒ 0 | **引用，不重跑**（182-c1 已量同形尺 0 且打了正控） | `.scratch/wisp/probes/182/c1/100-census.md` K6 那格两把尺 |
| 3 | 行 12 AC#1"票 181 已证认 git／工作树／分支三样**不需要**起进程" | **仍成立且更硬**：panel 包 `os/exec` 命中 **0**，且 git.go 头注释把"NO EXTERNAL PROCESS"写成了自陈 | R3：`grep -rn "os/exec\|exec\.Command" --include=*.go internal/panel/` ⇒ **0 行**（rc=1）；`90-red-roster-and-workspace.txt` Z4 = `internal/panel/git.go:5-8` 区段原文 |
| 4 | 行 13 AC#2"与票 181 同一枚文件族，不另起一份实现" | **⚠ 与 git.go 自陈冲突，且 09-28 件已给出解锁法**：git.go 明写"不起外部进程"，而未跟踪那一半只有问 git 才准 ⇒ 落地必然要么改这行自陈、要么用 09-28 件 §2 的"同族新开 `git_readonly_exec.go`"。两条都不违反"同一枚文件族"，但**必须具名改掉那句自陈** | Z4＋`189-uncommitted-count-design-a1.md` §2 末、§3.2 |
| 5 | 行 17 AC#6"`internal/panel` 那 **2 枚已知红**" | **过期**：现量在册红名册＝**4 枚**（起手 4／终态 4，逐名） | `90-…txt` Z1 引 `docs/evidence/s1/197-subagent-carrier-r3.md:268-275` |
| 6 | 行 17 AC#6"`gate-clauses.sh` 比名册（在册只 `G6neg`）" | **路径前提不完整**：仓根没有这枚脚本，实物在探针目录 | `70-…txt` M2：`ls -la .scratch/wisp/probes/154/gate-clauses.sh` ⇒ 存在，**48,783 字节，Oct 3 10:54**；`scripts/` 名册（`30-…txt` F3）里没有它 |
| 7 | 行 17 AC#6"`sh scripts/d22scan.sh`（基线 433）" | **本程没跑到那枚数（零 go 命令）**，但核了脚本形状：它**不含任何写死的基线数** | `60-gate-shape.txt` G1 全文 54 行：`grep -rn "433\|baseline\|基线" scripts/d22scan.sh` 的命中见 §6 缺陷 3；脚本本体是"正控＋实扫"两步，`set -eu`、无 `\|\| true` |
| 8 | 行 16 AC#5"D34 表不许加 git 行；C17 一枚不加" | **与既有代码形状一致**，且已被钉住：`TestGitDimensionHasNoModelCallableTool` 在册 | Z4 引 `internal/panel/git.go` 头注释原文（"It is not a tool… nothing here is registered in internal/tools"） |
| 9 | 行 4 出处普查 `docs/evidence/s1/180-182-panel-fields-census-c1.md` §5 | **行号仍在格上，但那普查的"没票认领"结论已被 189 自己作废**（普查写于立票前） | `80-…txt` K3 现读该件 `:128,:150-151,:178-179,:195,:206,:343`：S2/S3 两行仍写"没票认领 ⇒ 要立票"，而票 189 就是那张票。182-c1 §3 已具名报过同一处过期 |

---

## 2. 问②：数据源按**形状**找（⛔ 不是 grep `git diff` 字面）

### 2.0 查重（只引用，未重跑）
"零源"两把尺出自 `.scratch/wisp/probes/182/c1/100-census.md` K6 那格：`git diff|diff --numstat` ⇒ **0**（rc=1，正控＝同形尺 `.git/HEAD` 命中 `git_test.go:72…`）、`uncommitted|numstat|Uncommitted` ⇒ **0**。**本程未重跑这两把，复量的形状见 §2.3（不同尺）**，与它无冲突。

### 2.1 Go 侧有没有封装好的 git 命令执行器 ⇒ **面板侧一枚都没有；全仓只有工具目录一枚，且不在生产射程**

尺（`10-source-rulers.txt` R3/R4）：
- `grep -rn "os/exec\|exec\.Command" --include=*.go internal/panel/` ⇒ **0 行**（rc=1）
- `grep -rn "exec\.Command" --include=*.go internal cmd tools \| grep -v '_test\.go'` ⇒ **6 行／5 枚文件**（rc=0）

逐枚过目（**类型全名＋调用形状**，按派单口径）：

| 位置 | 调用形状 | 是不是 git | 在生产射程 |
|---|---|---|---|
| `internal/llm/adaptertest/mockllm.go:65` | `exec.Command(goBin, "build", "-o", exe, ".")` → `*exec.Cmd` | 否（go build） | 否（adaptertest 是测试夹具包） |
| `internal/llm/adaptertest/mockllm.go:72` | `proc.cmd = exec.Command(exe, "-addr", …)` | 否 | 否 |
| `cmd/balldebug/diff_windows.go:409` | `cmd := exec.Command(exe, childArgs...)` | 否（子进程调试器） | 是，但是 debug 台 |
| `cmd/wisp/doctor.go:316` | `exec.Command(cc, "--version").Output()` | 否（探针版本） | 是 |
| `cmd/wisp/slo_windows.go:492` | `cmd := exec.Command(exe, args...)` | 否（SLO 采样器） | 是 |
| `tools/d22scan/gitignore.go:373` | `exec.CommandContext(ctx, "git", args...)` | **是，全仓唯一一枚真起外部 git** | **否**——`tools/` 不是 `internal/`＋`cmd/` |

⇒ **"已有执行器"这一支不存在**：`tools/d22scan/gitignore.go:373` 是同仓先例但**不在生产包**，且它自己仍然"问不到 git 就响亮自陈并降级"（09-28 件 §3.1 记的 `:49-62` 语义，本程未复跑）。
⚠ 与 09-28 件 §3.1 那行"生产 Go 里 `exec.Command` **4 处**"**口径不同不冲突**：它的尺是 `--include=*.go internal/ cmd/`（不含 `tools/`）且按**枚文件**数；我按**行数**并含 `tools/` ⇒ 6 行／5 枚＋d22scan 一枚。两处都对得上，不是读数冲突。

### 2.2 有没有已经在读 `git status` 的模块 ⇒ **没有，一枚都没有**

- `grep -rniE "porcelain|git status|ls-files" --include=*.go internal cmd \| grep -v '_test\.go'` ⇒ **0 行**（rc=1）
- 正控：同尺**不排除测试** ⇒ **2 枚文件**（`cmd/wisp/panel_host_gate_test.go`、`cmd/wisp/slo_report_144_windows_test.go`）⇒ **尺能响**，负向成立（`10-source-rulers.txt` R6/R7）。

### 2.3 票 181／186 那条链今天**已落了哪几枚字段**（类型全名＋调用形状）

`internal/panel.GitView`（`internal/panel/git.go:117`）**10 枚 JSON 字段**，全部纯 Go 文件读出，**零进程**（`20-…txt` C1 原文）：

| JSON key | Go 字段 | 备注 |
|---|---|---|
| `kind` `reason` | `Kind` `Reason` | 四态 `GitKind*`（`git.go:65` 那条 `const (`）；"读不到"不许塌成"不是仓库" |
| `branch` `detachedSha` `isDetached` | `Branch` `DetachedSha` `IsDetached` | 读 `.git/HEAD` 首五字节 |
| `repoRoot` `currentWorktree` | `RepoRoot` `CurrentWorktree` | 后者＝`WorkspaceView.Canonical` 同一值（两节不可能指向两棵树） |
| `worktrees` | `Worktrees []GitWorktree`（`git.go:93`） | 主树优先＋递归 `.git/worktrees` |
| `branches` | `Branches []string` | `refs/heads` **递归**＋packed-refs 合并 |
| `switchBlocked` | `SwitchBlocked string` | **常量句子**（`git.go:82 GitSwitchBlockedReason`），不是判断 |

票 186 那一支 `internal/panel.WorkspaceView` **6 枚字段**＝`set`／`spelling`／`canonical`／`reparse`／`rewritten`／`reason`（`90-…txt` Z2，awk 抽 `composer.go` 结构体）。

调用链（现量行号，**全部与 09-28 件不同＝漂了**）：

| 环节 | 现量 file:line | 09-28 件引作 | 调用形状 |
|---|---|---|---|
| 读面 | `internal/panel/git.go:172`、`:240`、`:277`、`:150` | `:172`、`:240`（未漂） | `ReadGit(workspaceRoot string) GitView`、`ReadGitForWorkspace(ws WorkspaceView) GitView`、`GitViewRewritten(ws) GitView`、`GitViewNotProbed() GitView` |
| 载体字段 | `internal/panel/composer.go:246` | `:216` | `Git GitView` \`json:"git"\` |
| 空态兜底 | `internal/panel/composer.go:125`、`:285` | `:95`、`:232` | `composer.Git = GitViewNotProbed()` |
| 泵钩子 | `internal/panel/pump.go:150`（声明）／`:275,:278`（消费） | `:120` | `Git func() GitView`；`if p.src.Git != nil { composer.Git = p.src.Git() }` |
| 装配处 | **`cmd/wisp/run.go:703`** | `run.go:421-427`（那是 composer.go 注释里引的"泵装配根"） | `Git: rt.gitView,` |
| reader | `cmd/wisp/panel_pump.go:228-229` | `:108-109` | `func (rt *agentRuntime) gitView() panel.GitView { return panel.ReadGitForWorkspace(rt.workspaceView()) }` |

⇒ **这一跳今天整条是通的**：宿主读面 → GitView → 泵 → 出向包 → 装配，**已落地、可复用**。
⇒ **diff／未提交枚数这一支＝"整块没有"，不是"已有执行器缺解析"**：既无执行器（§2.1）、也无 status 解析（§2.2）、也无视图字段（§2.3 的 10 枚里没有一格装得下它）。

---

## 3. 问③：**最少落点表**（只做"未提交枚数"＋"逐文件名列表"，⛔ 不做整块 diff 渲染）

**形状裁定沿用 09-28 编排者已照准那三条**（`A395`，票面 Progress log 09-28 18:1x 行逐字）：枚数分两半（已跟踪纯 Go／未跟踪问 git，问不到＝第三个枚举值、绝不报 0）、落点＝`internal/panel/git.go` 同一枚读面、**嵌在 `composer.git` 下开子格、不加顶层 key**。本程只复核这个形状**今天仍然装得进去**，并给出行数估：

| # | 文件 | 动作 | 行数估 | 为什么必须有它 | 会不会碰契约面 |
|---|---|---|---|---|---|
| 1 | `internal/panel/git.go` | **改** | **+55～75** | `GitUncommittedView` 类型（`kind`/`reason`/`total`/`trackedModified`/`trackedDeleted`/`untracked`/`files[]`/`askedGit`）＋ `GitView` 那一枚新字段＋四态兜底（`newGitView`／`GitViewNotProbed`／`GitViewRewritten`／`ReadGit` 四条构造路径**每一枚都得填**，否则出现"零值＝看起来干净"那一坑） | **不碰 C17、不碰顶层键集**（见下"三枚免碰"） |
| 2 | `internal/panel/git.go` 头注释 | **改** | **±6** | `git.go:5-8` 区那句 **"NO EXTERNAL PROCESS"** 会被新读面证伪；按票面 AC#2"同一枚文件族"改自陈，⛔ 不许留着当假话 | 只是自陈文字，**不是** D/C 契约 |
| 3 | **`internal/panel/git_readonly_exec.go`**（新） | **新增** | **120～160** | 未跟踪那一半的唯一准路：`exec.CommandContext`＋固定字面量 `git --no-optional-locks ls-files -o --exclude-standard`（或 `status --porcelain`），同步 `.Output()`，`context.WithTimeout` 超时，git 缺席⇒`untrackedUnknown` 枚举 | **不碰 d22scan 七枚 ban**：ban #1 只咬裸 `go`、#4 只咬墙钟差（原文复核 `60-gate-shape.txt` G2，八枚里**没有一枚管起外部进程**）；⚠ 面板里 `approval.decide` 是 ban #6，与本落点无关 |
| 4 | `internal/panel/git_test.go` | **改**（或新开 `git_readonly_exec_test.go`） | **120～150** | 票面 AC#2 要的合成树判据：`git init`→写 3→`add` 1→改 1→删 1⇒四档与外部 `git status --porcelain` **逐枚一致**（两条尺钉，别只钉一个数——09-28 件 §4 那句）＋F5 那发 `.git/index` 跑前后 sha 不变 | 判据，不碰契约 |
| 5 | `internal/panel/git_*.go` 里的 **F2 依赖名册尺**（常驻 AST 判据） | **新增（可并成一枚测试文件）** | **60～85** | 票面 AC#3"补一枚常驻判据钉住审查面只读"的唯一硬法：import ⊆ 白名单、`os`/`io` 只准 `{Stat,Lstat,ReadFile,ReadDir,Open}`、`exec` 只准出现在 #3 那一枚具名单一文件、git 子命令逐枚 ∈ 只读名册。⚠ **必须自带正控**（种一发 `os.WriteFile` 门必须响） | 判据 |
| — | `internal/panel/pump.go`／`cmd/wisp/panel_pump.go`／`cmd/wisp/run.go` | **⛔ 免碰** | 0 | 现量证明：泵钩子形状是 `Git func() GitView`（`pump.go:150`），装配已给（`run.go:703`）⇒ **子格长在 GitView 里，整条链一行不改** | — |
| — | `internal/panel/bridge.go` | **⛔ 免碰＝不需要新面板方法名** | 0 | 入向名册现量**仍是 6 枚**：`panel.mode.request`/`panel.workspace.request`/`panel.attachment.add`/`panel.message.send`（`:42-45`）＋`config.get`/`config.set`（`:66-67`）；只读显示走**已有的快照推送通道**——票 181 AC#3 逐字已把这条钉死（`90-…txt` Z3） | **不触 C17＝不需要人工批准那一步** |
| — | `internal/panel/pump_test.go` 那两枚字节钉 | **⛔ 免碰** | 0 | 现量：钉的是**顶层**四键 `composer,generatedAt,pending,results`（`50-…txt` N2 `:111-124`、N3 `:270-276`）⇒ 加**嵌套子格**不动它 | — |
| — | `frontend/src/lib/panel.ts` | **本程判定＝免碰；界面那一跳另归** | 0 | 决定性现量：`composer_test.go:48` 的 `pairs` 名册只有 6 对（`Snapshot`/`ComposerState`/`ModeView`/`WorkspaceView`/`AttachmentRef`/`ResultChunk`），**`GitView` 不在其中** ⇒ 子格不与 TS iface 对账。⚠ 但**顶层 `git` 那一键本身**在 `panel.ts:112 ComposerState` 里**没声明**（N1 现量：6 枚键、无 `git`），这正是 4 枚在册红之一的红因；那一跳按 `composer.go:50`／`:90` 的 **"Q-51 - who may write panel.ts"** 归界面腿 | **不是 C17**，但**是 Q-51 未定案项**⇒ 要派界面腿，不是批契约 |

**合计**：**3～4 枚文件**（改 2＋新 1～2），**约 355～475 行**（其中产码 175～235 行、判据 180～235 行）。
**⛔ 不需要新的面板方法名 ⇒ 不需要走 C17 人工批准那一步。**（唯一需要"另派腿"的不是契约批准，是 Q-51 的界面声明那一跳。）

**代价在哪（一句话，其余见 09-28 件 §3）**：贵的不是行数，是**三枚"错了不报错"的坑**——① 纯 Go 朴素遍历把未跟踪报成 **8,684**（权威 **46**，高估 188×）；② `.git/index` 二进制偏移读错时探针**不 panic、不报错**，给出一个看起来完全合理的错枚数（`stat-dirty=3316/3332`）；③ 未跟踪问不到 git 时**塌成 0**＝"这棵树很干净"的假话。⇒ 落点必须带 F5 那发 `.git/index` 未变钉子＋"三档分裂"的独立枚举。

---

## 4. 问④：雷区（由页面发起的写／批准动作）——现量：这一堆里**一枚都没有**

| 发 | 尺原文 | 命中 | file:line／结论 |
|---|---|---|---|
| 1 | `grep -rniE "accept\|revert\|rollback\|commit\|push\|stage" design/doubao/demo/rb-review.js` | **0**（**rc=1**，无管道尾） | 审查原型 `rb-review.js`（6,177 字节）**没有任何动作词面**；它的全部 `data-*`＝`data-diff-toggle`（`:23/:63/:77` 文件头折叠）＋`data-lucide="git-pull-request"/"chevron-down"/"file-image"/"check"`（图标，`:17/:26/:66/:70/:80/:97`）⇒ **只有展开/折叠，属显示** |
| 2 | 正控：`grep -rnE "data-plg-toggle" design/doubao/demo/rb-plugins.js` | **2**（rc=0） | `rb-plugins.js:65/:75`⇒ 同尺**能命中真写动作**，第 1 发的 0 是真负向 |
| 3 | `grep -rnE "Method[A-Za-z]+ += +\"(approval\|git\|review\|vcs)\." --include=*.go internal cmd` | **0**（rc=1） | Go 侧**没有**任何面板可发起的批准／git／审查／版本动作腿 |
| 4 | `grep -nE '"(panel\|config)\.[a-zA-Z.]+"' internal/panel/bridge.go` | **6** | `bridge.go:42-45,66-67`，名册里**没有** decide/allow/commit/push/revert/stop（与 182-c1 §4 现状锚合一，未重跑其 F3/C3 那两把） |
| 5 | `grep -rn "os/exec\|exec\.Command" --include=*.go internal/panel/` | **0** | 宿主侧连"能起进程"的能力都没有⇒ 今天**结构上写不了盘** |

⇒ **票 189 票面没有夹带写动作**，AC#3（行 14）逐字要求"交付里不许出现任何提交／推送／接受／回滚的动作腿"，且票面行 20"本票**不**解决"已把四枚动作具名推到别处。
⇒ **真雷只在"登记文字"里，不在代码里**：180-182 census `:343` 现读的"雷区 4 处"＝面板 git 提交／推送、diff 接受／回滚、文件树删除／移动、＋（该件 §7 `:195` 起逐堆表）。其中与 189 有关的两枚（diff 接受／回滚、git 提交／推送）**今天无实现、无入向方法、无执行器** ⇒ **按 `AGENTS.md §1.2`"由面板侧来源的 L2『允许』"与票 181 AC#3 那四枚 `panel.*` 一枚不加，本程只指认、不代批、不并进只读交付表**。
⇒ 顺带一条**结构性副产品**（不是新增工作）：`GitSwitchBlockedReason`（`git.go:82`）已经把"切分支／切工作树"钉成常量句子，因为网页→宿主那一跳在这棵树里不存在——这一维天然给 §3 的第 5 发判据当先例。

---

## 5. 我这把尺的缺陷

1. **行数估是推的，不是量的**：§3 表里 5 枚"行数估"＝按既有同族文件形状推（`git.go` 现量 569 行、10 枚视图字段、4 条构造路径；`git_test.go` 本程**未数行数**）。⛔ 别当测量值引用。
2. **§2.3 的调用链只核了"名册＋行号"，没核语义等价**：`run.go:703 Git: rt.gitView` 我只证明装配点存在，**没读它周围是否还有别的 Git 覆写点**（例如 `GitViewRewritten` 在工作区改写后被谁调）——那一格若断，"免碰泵/装配"这一支要重算。
3. **AC#6 三把门禁尺一把都没跑**（零 go 命令硬约束）：`scripts/d22scan.sh` 只读了 54 行本体、`gate-clauses.sh` 只 `ls` 了存在性、`internal/panel` 测试没跑。⇒ **"基线 433"与"在册只 G6neg"两格本程＝〔未复跑，转引 197 件 §④.3 的 4 枚名册〕**。另：我在 `30-…txt` F8 那把 `grep … | head -20` 的 rc 是**管道尾**（182-c1 §6.1 点名的同一枚病），该尺**没有结论**，别再引它。
4. **§2.0 两把零源尺＝引用 182-c1，未复跑**；本程 §2.1/2.2/2.3 是**不同形状**的尺（exec/porcelain/视图字段），不构成对其复算。若编排者要"同一把尺双版"，得让下一腿重跑 182-c1 那两把。
5. **未跟踪枚数 8,684／46、diff 267 ms／1,521,885 B 全部转引 09-28 件 §3**，本程**一个数都没重量**（那需要跑 go／跑 git 计时，超出只读预算）。工作树自 09-28 起已变（`git status` 起手名册见本次 §0），那批数**今天必然不同**。
6. **F2 依赖名册尺的可行性只核到"ban 列表原文"层**（`60-gate-shape.txt` G2，八枚 ban 逐字读到 emoji 那一枚为止）：**没读每枚 ban 的正则实现**，所以"`exec` 会不会被 #1 的宽化matcher 咬到"这一问＝**仍未核**（09-28 件 §9 第 3 条同一缺陷，没被解决）。
7. **`jsonKeysOf` 是否收 `omitempty` 没读到实现**（`80-…txt` K6 rc=1：`composer_test.go` 里没这枚函数定义，它在别的测试文件里）。⇒ "`instructions`/`tasks` 让那枚对账变红"这一句是**从 197 件 `approval_test.go:129` 的失败文本反推**，不是我自己量的。
8. **`panel.ts` 我只读了 `:111-145` 一段**，其余 166 行未读；`GitView` 不在 `pairs` 名册这一条是从 `composer_test.go:48` 的 pairs 列表**逐项看完**得出的（6 对，无 Git 一支），可信度高但**没检查是否有第二枚对账测试**（`approval_test.go:105` 那枚只管 approval 卡片，本程未展开读）。
