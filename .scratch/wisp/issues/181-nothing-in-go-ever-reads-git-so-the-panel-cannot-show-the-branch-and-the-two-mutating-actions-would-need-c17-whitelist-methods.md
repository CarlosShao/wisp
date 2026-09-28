# 181 — 面板要显示"当前分支／远端分支／已有工作树"这件事，**Go 侧一行读 git 的代码都没有**（四把尺全 0 命中）；而"切分支／开新工作树"那两枚**动作**要动 C17 白名单（今天只有四枚 `panel.*`）＝契约变更 ⇒ **只读那半我推进，动作那半默认不做**（`Q-64`）

- Status: **ready-for-agent（只读普查＋一支只读落地）**；⚠ 本票**不含**任何改变 git 状态的写侧动作，那半在 `Q-64`、默认不做。
- 来路：09-28 10:2x，前端那侧递给 owner 的"要拍板三件事"里第 ③ 件（"git 分支/工作树 — 需要 Go 侧读 `.git/HEAD` 并推送分支名 ＋ 本地分支列表 ＋ 已有工作树，再加两条请求：切分支、开新工作树。这是 C17 白名单加方法，属于契约变更，得走批准"）。**它这句"C17 白名单加方法＝契约变更"是对的，我复算成立**；但它把**只读显示**与**写侧动作**捆成一件事交给 owner，于是 owner 被摆了一道本不必现在拍的题。⇒ 本票把两半**拆开**。台账 `A360`。
- 关联：票 145（快照扩成真载体——本票那几枚字段是它的邻居）· 票 92（composer 的 mode/工作区两格，同一形状："面板不许自造，Go 送真值"）· 票 163（`shell.session`，D34 表里那枚活会话）· `AGENTS §2` 未定义即停清单里那条"`C24` 初始集与 **`C17` 方法白名单定稿**"

## 现量（锚 `039efb47`，09-28 10:2x 本程现跑，每把尺可复制）

1. **Go 侧完全不读 git**：四把尺**全部空输出**——
   `grep -rn "\.git/HEAD" --include=*.go internal/ cmd/ | grep -v _test.go` ⇒ 空
   `grep -rn "git rev-parse" --include=*.go internal/ cmd/ | grep -v _test.go` ⇒ 空
   `grep -rn 'exec.Command("git"' --include=*.go internal/ cmd/ | grep -v _test.go` ⇒ 空
   `grep -rn "worktree" --include=*.go internal/ cmd/ | grep -v _test.go` ⇒ 空
   ⇒ **不是"接了一半"，是从零开始**。
2. **D34 那张内置工具权威表里没有任何 git 行**：尺 `awk 'NR>=2530 && NR<=2585' docs/PLAN.md | grep -nE "git|分支|branch"` ⇒ **0 命中**（表里有 `fs.list`／`shell.exec`／`shell.session`／`file.open`／`app.launch`／`input.type`，没有 git）。
   ⇒ **"给模型一个 git 工具"＝动 D34＝契约变更**；而**"给面板显示分支"不需要动 D34**（面板显示走 C17 宿主通道，不经过模型）——**这两条路必须分清，否则会把一件只读显示做成一次契约变更**。
3. **C17 白名单今天恰好四枚**：`internal/panel/bridge.go:42-45` 逐字 `panel.mode.request`／`panel.workspace.request`／`panel.attachment.add`／`panel.message.send`（尺：`grep -n '"panel\.' internal/panel/bridge.go | grep -v _test`）。
   ⇒ 前端要新增"切分支／开工作树"两枚请求＝**白名单 ＋2 枚方法**，且这两枚的**后果是换掉宿主自己正在跑的那棵代码树**。
4. **面板快照里没有 git 这一维**：`internal/panel/composer.go:200-207` 的 `ComposerState` 六枚字段＝mode／workspace／attachments／acceptedMimes／maxAttachmentBytes／attachmentError；`Snapshot` 四枚＝pending／results／composer／generatedAt（`composer.go:57-62`）。
   ⇒ 与票 145 的结论同形：**Go 没送，面板就不画**（那是"宁缺毋造"的既有裁定，不是前端的缺陷）。

## 为什么值得做（不做会怎样）

"我到底在哪棵树上"是**用户判断一切其它读数的坐标**：票 174 那条"回执要说实话（文件在不在、能不能读）"、票 147 那条"档位不许显示成未知"、以及 09-27 那次真事故（前端切分支导致我这 7 枚提交当时落在别的分支上，见 `HANDOVER §4.0u` 与 `A328`）——**都指向同一枚缺口：界面从不告诉任何人当前分支是哪枚**。今天他只能开命令行去问。⚠ 只读显示是**降低**下一次事故概率的东西，不是新增风险。

## AC（每格都要答"这一发在未修码上响不响"）

- [ ] **AC#1 先定"读什么、从哪儿读"，并把它写成一张有源的表**：三堆各答一句——① 当前分支（`.git/HEAD` 的 `ref:` 行，还是 detached 的裸 sha？两种都要能画）；② 本地分支列表＋远端分支名（**读 `.git/refs` 还是 `git for-each-ref`？后者要起子进程**，起子进程就撞上 D38 的 goroutine/owner 规矩与 `AGENTS §1.2` 的裸 `go func(` 禁令 ⇒ 必须现读那两处规矩再选）；③ 已有工作树（`.git/worktrees/` 目录枚举）。每堆都要给：**规格出处**（`docs/specs/SPEC-08*`／`SPEC-09*`／`PLAN.md` 里有没有要求过这一维，尺：`grep -rn "分支\|branch" docs/specs/ docs/PLAN.md | head -30`）、**真值来源**、**读失败时画什么**（⚠ 不许画成"master"或空数组顶替——票 92/145 的"未知要显式说未知"同形）。
- [ ] **AC#2 只读落地＋常驻判据**：Go 侧新增**宿主侧读面**（不是模型工具），把三堆送进快照；判据要钉"值来自真文件"——**未修码上今天不可能响**的那一发＝拿一棵**合成的假 git 树**（`git init` 一棵小树、切到非默认分支、加一枚 worktree）喂进去，断言面板侧 JSON 里读到的分支名／列表与那棵树逐字一致。⚠ 不许用"字段非空"充当判据。
- [ ] **AC#3 反向判据（防"把 git 读成工具"这一形走偏）**：交付里**不许出现**任何模型可调用的 `git.*` 工具、不许动 `docs/PLAN.md` 的 D34 表、`internal/panel/bridge.go:42-45` 那四枚 `panel.*` **一枚不许加**（只读显示走**已有的快照推送通道**，不需要新请求方法）。⇒ 若实现判"必须加方法才能读到"，**停手上报**，那说明它把只读那一半和动作那一半又捆上了。
- [ ] **AC#4 非 git 目录那一形必须显式回答**：工作区不是一棵 git 树（今天 `wisp run` 允许任意授权根）时，快照里这一维画什么？⚠ 三档候选（不显示／显示"非代码工作区"／显示未知）各给一句后果，**由我裁、不由写手自选**；不许静默给空数组（那与票 145 判"空数组＝按钮不画"是同一个坑）。
- [ ] **AC#5 契约轴**：`docs/PLAN.md`（D34／D36／D38 三节都近邻）、`docs/specs/**`、`allowlist.txt`、`thresholds.go`／golden 一字节不许动；`C1–C32`／`D1–D47` 任何一处措辞都不许在本票里改。**本票的读面是宿主内部实现，不是新契约**。
- [ ] **AC#6 门禁**：逐包 `go test -count=1 ./internal/panel/ ./internal/config/`（⚠ `./cmd/wisp/` 本机测不到东西＝票 98；⚠ `./internal/risk/` 要单跑，四包并发会假红＝争用，见 `A359`）＋`sh scripts/d22scan.sh`＋名册两向 `comm` 差集；⚠ **最终读数取在最后一枚 commit 之后**。

## 本票**不**解决（三件都具名，不许"以后加固"）

1. **切分支／开新工作树这两枚动作**：`Q-64`，**默认不做**。三条理由写在那一格：它换掉的是宿主自己正在跑的树；09-27 那次事故的直接教训就是"编队有程在飞时换树＝按不可逆动作对待"（`HANDOVER §4.0u`）；且它要 `C17` 白名单 ＋2 枚方法＝契约变更（`AGENTS §2` 已把"C17 白名单定稿"列为待人拍板项）。
2. **git 代码审查（diff 视图）**：不在本票，归票 182（那一栏的只读普查）。
3. **`frontend/**` 怎么画**：别家归属，本票只保证 Go 送出真值。

## Progress log

- 09-28 10:2x 编排者立票：上面四把尺本程现跑（锚 `039efb47`）；两半拆开的裁定记在 `A360`。未派。
- 09-28 14:3x **`Q-64` 改判＝做**（owner 原话入 `A368`）⇒ 本票标题后半句"两枚写动作要加 C17 方法所以另议"**不再是拦路条件**：那两枚写动作（切分支／切工作树）连同"识别 git、切本地目录"一起**收进新立的票 186**。⚠ **本票范围收窄为"只读那一半"＝票 186 的前置**（读什么、从哪儿读、非 git 目录那一形怎么如实说），AC 框与格面原话**一律不改**；两票的普查**并成一程派**（同一枚接缝，派单 `.scratch/wisp/dispatches/2026-09-28-144x-readonly-181-c1-plus-186-git-detection-census.md`）。
- 09-28 普查程 `181-c1` 交件（只读·**未勾任何 AC 框**·`internal/ cmd/` 零字节改动）：证据件 `docs/evidence/s1/181-186-git-detection-census-c1.md`，台件 `.scratch/wisp/probes/181/c1/git-census.sh`＋`logs/`。本程复跑上面四把尺 ⇒ ① 零枚命中（**复现**）、② `workspace.go:76`（行号未漂）、③ 零生产调用方（**但票面"定义＋两处注释"实测是定义＋三处注释**，多出的一枚是 `internal/panel/bridge.go:17`；结论不变）、④ `git status --porcelain -- internal/ cmd/` 空。新增现量：**三形都在本仓真数据里取到了**——目录形（`.git IS DIR`，`ref: refs/heads/dev`）、worktree 形（`.git` 是**文件**、内容 `gitdir: ...`，实测 `D:/wt/fe` 与 `D:/tmp/wisp136instr-r1/wt-head`）、detached 形（那枚 worktree 的 HEAD＝裸 sha `bb61dc5c…`）⇒ **"认出 git／当前分支／detached"不需要起外部 `git`**。⚠ **具名推翻派单的一把尺**：`ls .git/refs/heads | wc -l` 读 **3**、真值 `git branch --list` ＝ **4**（斜杠分支名在文件系统里是嵌套目录）⇒ 落地判据必须**递归走 `refs/heads`**。本仓 `.git/packed-refs` **不存在**（那一形没现量，见证据件 §⑨）。远端引用实测 6 枚 loose／2 枚 remote，**建议本票不列**（那是上次 fetch 的缓存，语义上会过期）。非 git 目录那一形的字段形状（`kind` 四值枚举＋必填 `reason`）与"读不到／没权限／不是仓库"三形分法交在证据件 §⑦；面板字段清单 9＋1 枚交在 §⑧（`frontend/**` 一字未读）。
- 09-28 15:1x **落地腿 `181-r1` 交件（只读那一半·未勾任何 AC 框）**。起手 HEAD＝`ee0ef8ec`（派单写的 `1f22f1aa` 已被同批他程推进，按实际 HEAD 做并登记此偏差）；起手写面闸门 `git status --porcelain -- internal/ cmd/` 空。
  - **产码两枚 commit**：`fccfe752`＝`internal/panel/git.go`（宿主侧只读面，零外部进程／零新依赖）＋`git_test.go`＋composer 段新增**一枚** `git` 键＋`PumpSources.Git` 读口＋`cmd/wisp/panel_pump.go:gitView()`（值＝`WorkspaceView.Canonical` 同源那枚）＋`run.go` 装配一行。`caebcbfe`＝现量逼出的一处缺陷修正（见下）。numstat 删除列：`fccfe752` 逐枚 **0**；`caebcbfe` 只在本票自己的两件里有 `-4/-5`（注释与夹具改写），冻结件与别家文件删除列 0。
  - **落了 10 枚字段**：`kind/reason/branch/detachedSha/isDetached/repoRoot/currentWorktree/worktrees[]/branches[]/switchBlocked`（`remoteBranches[]` 按派单**不送**）；来源尺＝普查 §②③④⑦⑧ 逐枚对应，无一枚自造。
  - **常驻判据五枚**（`internal/panel/git_test.go`）：`TestReadGitSyntheticSymrefRepo`（合成树，期望值全部是测试自己写进文件的串，**不含本仓分支名**）、`TestReadGitWorktreePointerFile`／`TestReadGitRelativePointerFile`／`TestReadGitVanishedPointerIsUnreadable`／`TestReadGitOnThisRepositoryIsSelfConsistent`（本仓真数据自洽：分支数＝`refs/heads` 递归文件数、工作树数＝附属枚数＋主树、附属路径必须解出）；反向判据 `TestGitDimensionHasNoModelCallableTool`（D34 表无 git 行＋`internal/tools` 无 `"git.*"` 工具名＋`bridge.go` 源码里 `panel.*` 仍恰好四枚——按源码文本读，不是重列常量）；载体判据 `TestGitSectionTravelsInTheSnapshotPacket`＋`TestPumpWithoutGitReaderSaysUnreadable`。变异现量（`-overlay`，跟踪件一字未动）：**M1** 去掉 `walkRefs` 递归⇒`TestReadGitSyntheticSymrefRepo`＋本仓自洽枚红；**M2** 摘掉主工作树那枚 append⇒三枚红；**M3** 相对 `gitdir:` 不做 join⇒`TestReadGitRelativePointerFile` 红；**M4** 泵里 `if p.src.Git != nil` 改成永假⇒载体枚红；**M5** 反向指针只认带前缀形⇒夹具＋本仓自洽两枚红。
  - **现量**：本仓目录形 `kind=repo branch=dev`、worktree 形（`D:/wt/fe`→`dsh/feat/frontend-p0-v2`）、detached 形（`D:/tmp/wisp136instr-r1/wt-head`→裸 sha `bb61dc5c…`）、非 git 临时目录→`kind=not_a_repo`＋`reason` 非空。分支数递归后＝**4**（＝`git branch --list`；旧尺 `ls` 仍读 3）；工作树清单＝**3** 枚含主树，与 `git worktree list` 逐字吻合。⚠ **本程被真数据逼出一枚缺陷**：`<common>/worktrees/<name>/gitdir` 是**裸路径**（不带 `gitdir:` 前缀，前缀只在链出树里那枚 `.git` 文件），第一版按前缀读⇒附属两枚 path 全读空；已由 `caebcbfe` 修正并把"两形都吃"钉进判据。
  - **门禁**：`go test -count=1 ./internal/panel/ ./internal/config/` ⇒ config **ok**；panel **3 枚红**＝已知 2 枚（`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`，他程删了 `design/assets/tokens.css` 未 staged，本程一字未碰 `design/**`）**＋新增 1 枚** `TestComposerContractTypesMatchFrontend`（`composer_test.go:74` 逐字 "Go ComposerState emits [git] that interface ComposerState does not declare"）＝票 145 §2 量过的双向对账半径，修法＝`frontend/src/lib/panel.ts` 补一枚 `git` 声明（那枚文件在本程 ⛔ 写面外／`Q-51`），**放宽 Go 侧断言不是选项**。`sh scripts/d22scan.sh` rc＝0，`ban #8 internal/` examined **433→436**＝本程 2 件（`git.go`／`git_test.go`）＋他程 1 件（`internal/tools/task_output_canonicalize_fail_174_test.go`，`git diff --name-status de9de24b..HEAD` 具名），`bans #1-5 internal/` 208→209（只多 `git.go` 这枚产码件）。`gate-clauses.sh` 红腿名册＝**只 `G6neg`**（在册那枚＝票 178），**没多出一枚**。`gofumpt -l` 六件全空。⚠ 一枚既有仪器差点挡住本票：`composer_test.go:268 gitSwitchCapabilityRe`（票 92 AC#7）把 `internal/panel/**` 产码里**单词 `worktree` 当"切换能力"扫**，`git.go` 初版的注释踩中⇒注释全部改写为"tree／linked tree"，判据未动；这条措辞约束归给后续 186 落地腿（字段名 `worktrees`／`currentWorktree` 复数或大写形不在尺射程，单数小写"worktree"会红）。
  - 未做（具名）：切换动作两枚（票 186）／入向那一跳（票 33）／`frontend/**` 如何画。本程未新增任何 `panel.*` 方法，未起外部 `git`，未跑 `probes/161/r6/flip-declaration.sh`，未跑任何删除命令，未 push。
