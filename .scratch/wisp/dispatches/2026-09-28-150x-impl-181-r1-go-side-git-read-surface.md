# 派单：`181-r1`（写码腿·Go 侧 git 只读面＋快照字段）——把"这棵工作区是不是仓库／当前分支／有哪些工作树"真送进快照

- 派单时刻：`2026-09-28 15:0x`｜编排者锚点 HEAD＝**`1f22f1aa`**（起手若已漂，登记并按实际 HEAD 做）
- 工单：`.scratch/wisp/issues/181-nothing-in-go-ever-reads-git-so-the-panel-cannot-show-the-branch-and-the-two-mutating-actions-would-need-c17-whitelist-methods.md`
- 上游读数（已入库、编排者复跑同值）：`docs/evidence/s1/181-186-git-detection-census-c1.md` ＋ 台账 `A371`

## 0. 一句话范围
**只做只读那一半**：Go 侧新增宿主读面，把 git 那一维填进**已有的快照推送通道**。**不做切换**（切换排在"网页→Go 那一跳"后面，那层今天不存在，见 `A371`），**不新增任何 C17 白名单方法**，**不起外部 `git` 进程**。

## 1. 写面清单（只这些，越界即退回）
✅ 新建：`internal/panel/git.go`（或同层你判更合适的名字）＋对应跟踪测试件 `internal/panel/git_test.go`
✅ 追加字段：`internal/panel/` 里 `Snapshot` 的 composer 段（⚠ **只到"新增字段"为止**——这是票 145 已批的局部解冻范围）
✅ 填值点：`cmd/wisp/panel_pump.go`（同样只到新增字段）＋ `cmd/wisp/run.go` 的装配处（若需要注入读面）
⛔ **冻结件一律不碰**：`internal/panel/tokens_fourway_test.go`、`l2_grant*`、`frontend_hygiene*`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`／golden／审批超时常量、`allowlist.txt`、`go.mod`／`go.sum`（**不引任何新依赖**）
⛔ `frontend/**`／`design/**` **不读不写不引不转述**（界面画不画归另一枚会话）
⛔ 不属于你的脏件（`.gitignore`、`probes/152/**`、`probes/161/r6/logs/flip-*`、`design/**` 那批未提交删除、`docs/evidence/s1/152-*.md`、别人在飞的探针）**不动、不提交、不评论**
⛔ Git：只 commit 不 push；显式 pathspec；禁 `add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`worktree`／`switch`；**零删除命令**；不许跑 `probes/161/r6/flip-declaration.sh`

## 2. 要交付的字段（9＋1 枚，来源全在 `181-c1` 表 §⑧，逐枚带尺）
`git.kind`（四值枚举 `repo`／`not_a_repo`／`permission_denied`／`unreadable`）＋`git.reason`（**必填**，不许留空——票 92/145 的"宁缺毋造"）＋`git.branch`／`git.detachedSha`＋`git.isDetached`＋`git.repoRoot`＋`git.currentWorktree`（**与已有的 `WorkspaceView.Canonical` 同源**，不许另算一套）＋`git.worktrees[]`（⚠ **必须含主工作树**——`.git/worktrees/` 里只有附属那几枚，主树不在其中，`181-c1` 已现量：本仓 `ls .git/worktrees`＝2 枚，`git worktree list`＝**3** 枚）＋`git.branches[]`（⚠ **必须递归走树**——派单旧尺 `ls .git/refs/heads` 数出 3、真值 4，**斜杠分支名是嵌套目录**，那把尺已被 `181-c1` 具名推翻）＋状态枚 `git.switchBlocked`（今天恒"不可切"，理由＝入向那一跳未建）。
`git.remoteBranches[]` **默认不送**（那是上次 fetch 的缓存，不是当前状态）。

## 3. AC 对应（每格都要答"这一发在未修码上响不响"；**AC 框不许自己勾**，勾由编排者翻）
- **AC#1**：读面落点与"从哪儿读"照 `181-c1` 表实现（目录形 `.git/HEAD` 的 `ref:` 行；worktree 形 `.git` 是**文件**、内容 `gitdir:`，仓库根再靠同目录 `commondir`；detached＝那枚 HEAD 是裸 sha）。⚠ `gitdir:` 指向已消失目录 ⇒ 算 `unreadable`，**不许**报成"不是仓库"。
- **AC#2 正向常驻判据**：值必须**来自真文件**——在临时目录里**合成一棵树**（自己 `mkdir`＋写 `.git/HEAD` 等），断言读面对它给出正确答案；⚠ **不许把期望值硬编码成"本仓当前分支名"**（那会让判据在别的机器／别的分支上假绿或假红）。
- **AC#3 反向判据（防走偏）**：交付里**不许出现**任何模型可调用的 `git.*` 工具、不许动 `docs/PLAN.md` 的 D34 那张表、读面不许挂在 `internal/tools` 的执行面上。⇒ 补一枚常驻判据钉住"这一维只有宿主侧读面＋快照两条路"。
- **AC#4 非 git 那一形必须显式回答**：给一发真判据（临时目录里什么 `.git` 都没有 ⇒ `kind=not_a_repo` ＋ `reason` 非空）。
- **AC#5 契约轴**：上面 ⛔ 那一串一字未动，交件前逐枚复量并贴尺。
- **AC#6 门禁**：`go test -count=1 ./internal/panel/ ./internal/config/`＋`sh scripts/d22scan.sh`（基线 `ban #8 internal/` **examined=433**，你新增跟踪件会变 434/435，**具名解释**）＋`bash .scratch/wisp/probes/154/gate-clauses.sh` **比红腿名册不比退码**（在册只 `G6neg`＝票 178；⚠ **不要改这把尺本身**）；⚠ `./internal/risk/` 若要跑必须**单跑**（`A359`）。终态读数取在**你自己最后一枚 commit 之后**。
  ⚠ **`./internal/panel/` 今天已有 2 枚已知红**（`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`＋`TestC21DesignTokensFourWayAgree`，同一因＝别家会话删了 `design/assets/tokens.css` 未 staged）⇒ **照实记、不许当绿、不许顺手修、不许碰那两枚文件**；你的红名册要与这 2 枚**逐名比对**后报"新增几枚"。
- **gofumpt**：`"$(go env GOPATH)/bin/gofumpt" -l <你的两枚新件>` 必须为空（它不在 PATH 上）。

## 4. 争用护栏（这仓今天被咬过三次）
起手跑 `git status --porcelain -- internal/ cmd/`：**不为空就停下报我**（别在别人脏件上跑判据）。⚠ 同批另有 **2 枚只读程在飞**（`174-c2`、`180-c1＋182-c1`）＋可能一枚 `33-a1`——它们只写 `docs/evidence/s1/**` 与 `probes/**`，与你的 `internal/panel/**` 不撞；但**它们也会跑 `./internal/panel/` 那格门禁**，所以你落完码之后它们的读数会带上你的新判据——**这是正常的，别去改它们已入库的读数**。⚠ **重跑任何既有台件前先 `grep -n "WriteFile\|OpenFile" <那枚台件>` 看写出路径**（`A367`：就地覆盖型 sink 会洗掉别家票逐行引用的读数）。

## 5. 增量交付与硬顶（上一枚同职能程死于 52/50 枚、零交付）
工具调用**硬顶 45 枚**，**到第 32 枚停止新探索**。顺序：① 先落"读面骨架＋一枚最小判据"并 commit（第 ≤8 枚调用内）；② 之后每完成一个字段族 commit 一枚；③ 门禁与票面 Progress log 追加留到最后一次。**表没落盘之前不许继续加字段**。判"必须动冻结件／必须引依赖才能做完"＝**停手上报**，不要自己拍。

## 6. 结束消息回我八项
① 落了哪几枚字段、逐枚来源尺；② 两枚常驻判据的名字与"摘掉哪一处会红"的变异现量（走 `-overlay`，别动跟踪件）；③ 三形现量（目录形／worktree 形／detached 形）＋非 git 那一形；④ 分支枚数那把旧尺你现跑到的是几（应＝递归后与 `git branch --list` 同值）；⑤ 工作树清单是否含主树（与 `git worktree list` 逐字对照）；⑥ 门禁四数＋**红腿名册**＋`internal/panel` 红名册与那 2 枚已知红的逐名差集；⑦ 冻结面自证（`git diff --numstat` 删除列逐枚 0、`go.mod` 未动、`internal/ cmd/` 起手与终态为空）；⑧ 被拒调用＋有没有跑删除命令＋工具调用终值枚数＋next。
