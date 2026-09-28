# 派单：`181-c1`（只读普查·零产码）——聊天框那一排的"认出 git／切本地／切工作树"到底要读什么、跑不跑外部进程、接线那一跳今天存不存在

- 派单时刻：`2026-09-28 14:4x`
- 编排者锚点：HEAD＝**`72c76d42`**（票 183 已结线 `-done`）
- 工单（两枚并一程，同一接缝）：
  - `.scratch/wisp/issues/181-nothing-in-go-ever-reads-git-so-the-panel-cannot-show-the-branch-and-the-two-mutating-actions-would-need-c17-whitelist-methods.md`（只读那一半＝前置）
  - `.scratch/wisp/issues/186-the-composer-area-should-detect-git-and-let-the-user-switch-between-local-worktree-and-branch.md`（识别＋切换）
- 背景一句话：owner 09-28 当场改判 `Q-64`＝**做**（原话在台账 `A368`）。⇒ **本单不是"该不该做"，是"做起来要动哪几层、哪一层今天根本没有"**。
- 性质：**只读普查** ⇒ **AC 框一枚不许勾**、**产码零字节不许动**。

## 0. 先说这单最容易做歪的地方

1. **不许把结论写成"要不要做／该不该做"**——方向已由 owner 定。你只交**现量与代价**，落点由编排者裁。
2. **别默认"跑一下 `git` 就行"**。全仓今天**零次**外部 `git`（尺见 §2 第 1 行），起外部进程是一枚**新的宿主能力**：要么给出"不跑也能读全"的证据，要么把"必须跑"具名报回并说清它归不归哪条在册禁令管（去读 `tools/d22scan/main.go` 的 ban 列表原文，**别猜**）。
3. **别只答"读得到分支名"就交差**。工作树／多目录那一形里，一枚工作树的 `.git` **是文件不是目录**（内容形如 `gitdir: ...`），本仓自己就可能带 worktree ⇒ 这是最容易写错的一形，**必须现量**。

## 1. 写面（只这些）

- ✅ 新建 `docs/evidence/s1/181-186-git-detection-census-c1.md`（**骨架先落盘**，见 §5）
- ✅ 新建 `.scratch/wisp/probes/181/c1/**`（读数／只读小台件；台件放探针目录，**不进跟踪代码**）
- ✅ 追加：票 181 与票 186 各一段 Progress log（**只追加、不改原句、不勾框**）
- ⛔ `internal/**`、`cmd/**`、`docs/PLAN.md`、`docs/specs/**`、`allowlist.txt`、`thresholds.go`／golden／审批超时常量 **一字不动**
- ⛔ 不新增 C17 白名单条目（那是票 186 AC#2，要编排者先入档）；⛔ 不做成模型可调的 `git.*` 工具（D34 那张表里没有任何 git 行）
- ⛔ `frontend/**`／`design/**` 不读不写不引不转述（界面那半归另一枚会话；你只交"面板要哪几维数据"的字段清单）
- ⛔ Git：只 commit 不 push；显式 pathspec；禁 `add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`worktree`／`switch`；**临时件只建不删**（不许跑任何删除命令）；**不许跑 `probes/161/r6/flip-declaration.sh`**
- ⚠ **重跑任何既有台件之前先 `grep -n "WriteFile\|OpenFile" <那枚台件>` 看它往哪枚路径写数**——就地覆盖型 sink 会把别家票逐行引用的读数洗掉（09-28 实测＝`A367`）。要写就写到你自己的 `probes/181/c1/logs/`。
- ⚠ 工作树里那些不是你的脏件（`.gitignore`、`probes/152/**`、`probes/161/r6/logs/flip-*`、`design/**`、`docs/evidence/s1/152-*.md`、`probes/183/**` 等）**不动、不提交、不评论**。
- ⛔ 凭据：任何密钥明文不许进读数文件（路径与哈希可以）。

## 2. 起手先把这四把尺跑掉（读数进表 §1）

```
grep -rln 'exec\.Command("git")\|\.git/HEAD\|rev-parse' --include=*.go internal/ cmd/ | grep -v _test.go   # 编排者 14:3x 现量＝零枚
grep -n "func RequestWorkspaceSwitch" internal/panel/workspace.go                                            # 现量＝:76
grep -rn "RequestWorkspaceSwitch" --include=*.go internal/ cmd/ | grep -v _test.go                          # 现量＝只有定义＋两处注释，零调用方
git status --porcelain -- internal/ cmd/                                                                     # 必须为空；不为空就停下报我
```

⚠ 你引的任何行号都要现跑（票面与台账里的号会漂，`A367`/`A368` 各有一例）。

## 3. 必答的六问（每问一条现跑命令＋读数，不许"应该可以"）

**Q1 认 git**：从工作区根往上找 `.git` 到仓库根，需要读哪几枚文件？`.git/HEAD` 的 `ref: refs/heads/dev` 与 detached（一枚裸 sha）两形各给一条现量；**worktree 那一形**（`.git` 是文件、内容 `gitdir: ...`）在本仓现量一次（尺：`git worktree list --porcelain`，然后 `cat` 你看到的每一枚 `.git` 并 `test -d`／`-f` 分辨）。⚠ 若"不跑外部进程就读不全"，具名说缺哪一维。

**Q2 工作树清单**：只读文件系统能不能拿到"每棵工作树的路径＋它当前检出的分支"？候选读法：`.git/worktrees/<name>/HEAD`＋`gitdir` 反向指回。逐枚现量并给结论：**能读全／缺什么／缺的那部分要不要跑 `git worktree list --porcelain`**。

**Q3 分支清单**：`refs/heads/` 与 `packed-refs` 两堆怎么合？本仓现量两枚数（尺：`ls .git/refs/heads | wc -l`；`grep -c "^# " .git/packed-refs 2>/dev/null`；`git branch --list | wc -l` 只作对照真值）。⚠ 远端分支（`refs/remotes/`）要不要列，列了算不算"读面扩到远程"——给结论，别顺手扩。

**Q4 接线那一跳今天存不存在**（⚠ 这一问决定本票是"接线"还是"从零建"）：面板→Go 的**入向**通道，逐枚答——① `panel.workspace.request` 这个名字在册（`internal/panel/bridge.go:43`），**谁在消费它**？（尺：`grep -rn "MethodWorkspaceRequest" --include=*.go internal/ cmd/ | grep -v _test.go`）② 一枚已经能用的**同类先例**长什么样：`panel.mode.request` 的"请求→`ModeRequest.Parse`→`HandleModeRequest`→写者＋审计"整条链在哪些文件、每一环有没有生产调用方（尺：`grep -rn "HandleModeRequest\|ModeRequest\b" --include=*.go internal/ cmd/ | grep -v _test.go`）。⚠ 票 114 曾裁"网页事件→Go 那一跳不存在"——**今天到底存不存在，以你现跑为准，并具名推翻或坐实**（这条最值钱）。
③ 若那一跳不存在：本票"切换"要落地还欠哪几环？逐环点名文件与要新增的东西（**只列，不写**）。

**Q5 切换的连带面（票 186 AC#6 的雏形）**：换掉工作区根之后，这几样各自会怎样——授权根／`allowed_dirs`（C26）、产物目录（`agent.Spiller` 的根）、C25 污染名册与 per-scope mark、会话历史与取证行、上下文预算 `BudgetsFor(<窗口>)`。每样给一条现跑（`grep -rn "SetWorkspaceRoot\|WorkspaceRoot()" --include=*.go internal/ | grep -v _test.go` 起步），并明确"今天切了会不会留下上一棵树的证据"（⚠ 这一维在票 183/185 那根管子上已经咬过三次，别当理论题）。

**Q6 非 git 目录那一形**（票 181 AC#4）：工作区不是 git 树时（今天 `wisp run` 允许任意授权根），快照里这一维要送什么才算"说实话"？给一枚具体字段形状与一句人话文案（**不许留空**，宁缺毋造＝票 92/145 同形），并说明它和"读不到／没权限／不是仓库"三形怎么分。

## 4. 门禁（只读程也取数，但按 `A363` 薄规矩**不充当 AC 结案凭据**，只作现状记录）

`sh scripts/d22scan.sh`（基线 rc＝0、`ban #8 internal/` **examined=433**，多一枚少一枚都具名解释）；`bash .scratch/wisp/probes/154/gate-clauses.sh` **比红腿名册不比退码**（在册只有 `G6neg`＝票 178）；`go test -count=1 ./internal/panel/ ./internal/config/`（⚠ `./cmd/wisp/` 本机测不到东西＝票 98；⚠ `./internal/risk/` 若要跑必须**单跑**＝`A359`）。**终态读数取在你自己最后一枚 commit 之后。**

## 5. 表骨架（第 ≤6 枚调用内落盘并 commit 第一枚；之后每两问 commit 一枚）

`181-186-git-detection-census-c1.md` 十二节：① 起手锚＋写面闸门＋§2 四把尺 ② Q1 认 git ③ Q2 工作树清单 ④ Q3 分支清单 ⑤ Q4 接线那一跳（含"推翻或坐实票 114 那句"） ⑥ Q5 切换连带面 ⑦ Q6 非 git 那一形＋建议字段与文案 ⑧ **面板要的字段清单**（交编排者转前端会话，你不碰前端） ⑨ 本程没测什么（逐名，别写"其余都覆盖了"） ⑩ 门禁终态 ⑪ 被拒／没成功的调用＋有没有跑删除命令＋工具调用终值自报 ⑫ next＝**落地腿派之前还缺什么**（含"哪几枚要人先批准"）

## 6. 硬顶与增量交付（上一枚同职能程死于 52/50 枚、零交付）

工具调用**硬顶 40 枚**，**到第 28 枚停止新探索**，余量只用于填表、跑门禁、commit。**表没落盘之前不许继续取数。**

## 7. 结束消息必须回给我的十项

① 六问逐问一句话结论＋最贵的代价；② **Q4 的裁定**：接线那一跳存在／不存在（带现跑证据，具名推翻或坐实票 114）；③ 要不要跑外部 `git`（要＝具名说哪一维非跑不可，以及 `tools/d22scan/main.go` 里有没有管它的条款）；④ Q5 那五样连带面逐样答案；⑤ 面板字段清单（枚数＋每枚来源）；⑥ `internal/ cmd/` 写面自证（起手／终态两条 `git status --porcelain` 都为空）；⑦ 门禁四数＋红腿名册；⑧ 被拒调用逐条；⑨ 有没有跑删除命令＋工具调用终值；⑩ next＝还缺什么、哪几枚要人先批准。

**不许自己拍落点、不许自己加白名单条目、不许改任何冻结件。碰到"必须动冻结件才能往下"＝停手上报。**
