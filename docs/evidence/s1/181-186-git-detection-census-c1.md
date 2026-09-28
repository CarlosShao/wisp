# 票 181＋186 并程普查（`181-c1` 只读证据件）

- 程：`181-c1`（只读普查·零产码）· 派单逐字：`.scratch/wisp/dispatches/2026-09-28-144x-readonly-181-c1-plus-186-git-detection-census.md`
- 工单：票 181（只读那一半＝前置）＋ 票 186（认出 git／切本地／工作树／分支）。owner 已改判 `Q-64`＝**做**（`A368`）⇒ 本件**不论证该不该做**，只交**现量与代价**；落点由编排者裁。
- 起手时刻：`2026-09-28 14:32:50 +08`。起手 HEAD＝**`c1c96008`**（`Mon Sep 28 14:25:24 2026 +0800`）。
  ⚠ 派单写的编排者锚是 `72c76d42`（票 183 结线）；共享工作树里 HEAD 已被别人继续推进 ⇒ **本件一切行号/读数按 `c1c96008` 现跑**，票面引的号若漂了以本件现跑为准（`A367/A368` 各有一例漂号，派单 §2 的提醒）。
- 起手写面闸门：`git status --porcelain -- internal/ cmd/` ⇒ **空**（14:32:50 实测，`STATUS-INTERNAL-CMD-EMPTY-OK`）。
- 骨架先落盘（派单 §5＋§6 硬规矩）：本件十二节骨架在**第 5 枚调用**落盘并 commit 第一枚；此后**每答完两问 commit 一枚**；表没落盘之前不再取数。
- ⚠ 本件不勾任何 AC 框；不新增 C17 白名单条目；`frontend/**`／`design/**` 不读不写不引不转述。

## ① 起手锚＋写面闸门＋派单 §2 四把尺（现跑）

| 尺（逐字） | 读数 | 与票面比 |
|---|---|---|
| `grep -rln 'exec\.Command("git")\|\.git/HEAD\|rev-parse' --include=*.go internal/ cmd/ \| grep -v _test.go` | **空输出（零枚）** | 坐实票 181 现量第 1 条＋票 186 现量第 1 行 |
| `grep -n "func RequestWorkspaceSwitch" internal/panel/workspace.go` | `76:func RequestWorkspaceSwitch(scope PathScope, input string, audit AuditFunc) (WorkspaceView, error) {` | 行号 `:76` **未漂** |
| `grep -rn "RequestWorkspaceSwitch" --include=*.go internal/ cmd/ \| grep -v _test.go` | 4 行＝定义 `workspace.go:76` ＋注释 `workspace.go:70`／`bridge.go:17`／`cmd/wisp/panel_pump.go:79` ⇒ **零生产调用方** | ⚠ 票 186 现量第 3 行写"两处注释"，实测**三处**（多出的是 `bridge.go:17`）；**结论不变**，只是票面数差一枚 |
| `git status --porcelain -- internal/ cmd/` | 起手 14:32:50 空；骨架 commit `e9ef94d0` 之后再跑**仍空**（`GATE-STILL-EMPTY`） | 闸门过 |

## ② Q1 认 git：读哪几枚文件；三形各一条现量

**要读的文件（宿主侧只读面，零外部进程）**：① 从工作区路径起**逐级往上**对每层做 `test -d <dir>/.git`／`test -f <dir>/.git`（两形都必须吃，见下）；② 目录形 ⇒ 读 `<root>/.git/HEAD`；③ 文件形 ⇒ 先剥 `<root>/.git` 里的 `gitdir: ` 前缀拿到 per-worktree 私目录，再读 `<gitdir>/HEAD`，仓库根另由同目录的 `commondir`（本仓实测 `../..`）解出；④ 判断"是不是仓库"＝上面任一命中。

**三形现量（全是本仓真数据，不是构造）**：

| 形 | 现跑尺 | 读数 |
|---|---|---|
| 目录形＋symref | `test -d .git && cat .git/HEAD` | `.git IS DIR`；HEAD＝`ref: refs/heads/dev`；对照真值 `git rev-parse --abbrev-ref HEAD`＝`dev` ⇒ **一致** |
| **worktree 形（`.git` 是文件）** | `test -f /d/wt/fe/.git`；`cat` 之；再 `cat $(gitdir)/HEAD` | 内容逐字 `gitdir: D:/work/workspace/projects plans/Wisp/.git/worktrees/fe`（**绝对路径且含空格**；git 也写相对形 ⇒ 读侧两种都要吃）；其 HEAD＝`ref: refs/heads/dsh/feat/frontend-p0-v2` |
| detached 裸 sha 形 | `cat /d/tmp/wisp136instr-r1/wt-head/.git`；再 `cat $(gitdir)/HEAD` | `.git` 同是文件（`gitdir: .../.git/worktrees/wt-head`）；其 HEAD＝`bb61dc5c0c5b15783db8736ea5e6fccd6fa50e3f`（**40 位裸 sha、无 `ref: ` 前缀**）；`git worktree list --porcelain` 对同一枚打 `detached` ⇒ 判据＝"HEAD 首 5 字节是不是 `ref: `" |

**结论**：三形**全部只读文件系统可读**，本票这一维**没有"非跑外部 `git` 不可"的缺维**。

## ③ Q2 工作树清单：只读文件系统够不够

- 枚举措：`ls .git/worktrees` ⇒ 本仓实测 **`fe`、`wt-head` 两枚**；逐枚读 `HEAD`（分支名或裸 sha）＋ `gitdir`（反向指回那棵树工作目录下的 `.git` 文件 ⇒ **工作树根＝该文件的父目录**）。实测：`fe` 的 gitdir＝`D:/wt/fe/.git` ⇒ 根 `D:/wt/fe`；`wt-head` 的 gitdir＝`D:/tmp/wisp136instr-r1/wt-head/.git` ⇒ 根 `D:/tmp/wisp136instr-r1/wt-head`。
- 对照真值（`git worktree list --porcelain`）＝**3 条**，逐字吻合：主仓 `D:/work/workspace/projects plans/Wisp`→`refs/heads/dev`；`D:/wt/fe`→`refs/heads/dsh/feat/frontend-p0-v2`；`D:/tmp/wisp136instr-r1/wt-head`→`detached`。
- ⚠ **最大的坑（现量坐实）**：**主工作树自己不在 `.git/worktrees/` 里** ⇒ 清单必须由"主仓（`.git/HEAD` 直接读）＋ `worktrees/` 枚举"**两堆拼**；只枚举目录会**漏掉用户当前就在的那枚**。
- 缺不缺：porcelain 还多报 `locked`／`prunable`／`bare` 三型附加态，文件形读法＝看 `.git/worktrees/<name>/` 下有无 `locked`／`prune` 文件（本仓两枚实测**都无**，目录里是 `COMMIT_EDITMSG HEAD ORIG_HEAD commondir gitdir index logs refs`）。⇒ **本票要的"路径＋分支"两维不缺**，**不必跑 `git worktree list --porcelain`**；要不要把 locked/prunable 也画进面板是产品决定（本程不替它裁）。

## ④ Q3 分支清单：`refs/heads`＋`packed-refs` 两堆

- 本仓两枚数（逐字尺＋读数）：`ls .git/refs/heads | wc -l`＝**3**；`find .git/refs/heads -type f | wc -l`＝**4**；`test -f .git/packed-refs`＝**NO（本仓根本没有 packed-refs）**；对照真值 `git branch --list | wc -l`＝**4**（`dev`／`dsh/feat/frontend-p0`／`dsh/feat/frontend-p0-v2`／`master`）。
- ⚠ **具名推翻派单给的这把尺**：`ls .git/refs/heads | wc -l` 读出 **3**，与真值 **4** 不等——因为**带斜杠的分支名在文件系统里是嵌套目录**（`refs/heads/dsh/feat/frontend-p0`）。⇒ 落地判据必须是**递归走 `refs/heads` 整棵树**（并跳过目录本身）；任何"单层列举"的写法在本仓**今天就会少报一枚**。
- 两堆的合法：本仓 packed-refs 不存在 ⇒ **只走 loose 就读全**；但读面仍要备 packed 形（行格式 `<sha> refs/heads/<name>`，`#` 头行跳过，**同名时 loose 覆盖 packed**，按名字去重）。这一形**本仓给不了现量**，见 §⑨。
- 远端分支：`.git/refs/remotes` 存在，`find .git/refs/remotes -type f | wc -l`＝**6**，两枚 remote（`cnb`、`origin`；`.git/config` 里 `url = ` 出现 2 次，**本件不抄 URL**）。⇒ 读它们是**本地磁盘读数、不起网络不起进程**，字面上不算"读面扩到远程"；但语义上那是**上次 fetch 的缓存快照**（会过期），而 owner 点的这一排是"本地／工作树／分支"。**建议：本票不列远端分支**（列了就得再分两堆＋带一句"可能过期"）；真要列，那是**字段扩枚**不是读面越界——裁量交编排者，本程不擅自扩。

## ⑤ Q4 接线那一跳今天存不存在（⚠ 本件最值钱的一问）

**裁定：入向那一跳今天不存在——具名坐实票 114**（`.scratch/wisp/issues/114-composer-request-has-no-production-caller-and-the-real-gate-must-be-native.md`，标题自己就是"no production caller"）。现跑到 HEAD `c1c96008`／本件骨架 `e9ef94d0`，逐环：

| 环 | 位置 | 生产调用方（`grep -v _test.go`） |
|---|---|---|
| C17 白名单在册 4 枚 | `internal/panel/bridge.go:42-45`＝mode.request／workspace.request／attachment.add／message.send | — |
| `MethodWorkspaceRequest` 的消费点 | 只有 `bridge.go:106` 那枚 `knownComposerMethod` 的 `case`（**在册＝有名、无写腿**） | **零枚其他命中** |
| 封套解析 `ParseComposerRequest` | `bridge.go:84` | **零生产调用方**。两处以自己口吻写明不存在：`composer_handlers.go:37`、`cmd/wisp/run.go:229`（逐字 "Nothing calls it yet - the WebView2 \"event -> ParseComposerRequest\" hop does not exist in this tree"） |
| 同类先例（mode 那一族整条链） | `ModeRequest`/`Parse`＝`composer.go:141/152` → `ModeWriteHandler.HandleModeRequest`＝`composer_handlers.go:111` | **`HandleModeRequest` 零调用方**（`grep -rn "\.HandleModeRequest("` 非 test ⇒ `No matches found`）；handler 只在 `cmd/wisp/run.go:233`（字段）＋`run.go:417`（构造）挂在 runtime 上，**构造了没人调** |
| 入向物理通道 | `grep -rn "WebMessage\|ReceiveMessage\|OnMessage\|PostMessage"` 非 test ⇒ 只命中 `internal/ball` 的 Win32 `PostMessageW`（托盘／STA 自用） | **没有一枚 WebView2 消息接收器** |
| 宿主本身在不在 | `internal/panel/pump.go:16` 逐字 "tree today - no WebView2 host (ticket 33 is unclaimed), no postMessage writer"；`cmd/wisp/panel_assets.go:13` 同源；票 33 文件仍在池中未 `-done` | — |

**出向**（宿主→页面快照）今天**真的在跑**：`cmd/wisp/panel_pump.go` 的 `workspaceView()`（:82-86，值来自 `rt.paths.WorkspaceRoot()`）＋`publishPanelSnapshot()`（:241）＋`bookPanelSnapshot()`（:159）。⇒ **面板"显示"这一半有管子，"请求"这一半整条没有**。

**所以票 186 的"切换"不是接线，是从零建一整层**（只列不写）：
1. 票 33 的 WebView2 宿主本体（拉起／STA／C27 单例）——**本票范围外，硬前置**。
2. 宿主里"页面 postMessage → raw string"的接收器，喂现成的 `ParseComposerRequest`（解析器**零改动**可用）。
3. 一枚按 `req.Method` 分派的 router——**仓里今天没有任何 router**，只有 handler 方法本体。
4. `WorkspaceWriteHandler` 那一枚壳（照 `ModeWriteHandler` 的形状包 `RequestWorkspaceSwitch`＋审计＋回显）——**被包的函数已有、语义不许改写**。
5. 快照里 git 这一维的字段（`ComposerState` 今天六枚字段无 git 维＝票 181 现量第 4 条）。
6. **git 只读面本体**（§②③④ 那三堆读法）＝本票唯一真正"从零写"的产码，且**不依赖票 33**。

## ⑥ Q5 切换的连带面（票 186 AC#6 雏形）：授权根／产物目录／C25 名册与 per-scope mark／会话历史与取证行／上下文预算 `BudgetsFor`——逐样"保留／重建／作废"＋"切了会不会留下上一棵树的证据"

（待填·⑥）

## ⑦ Q6 非 git 目录那一形（票 181 AC#4／票 186 AC#4）：快照这一维送什么才算"说实话"＋一句人话文案＋三形怎么分（读不到／没权限／不是仓库）

（待填·⑦）

## ⑧ 面板要的字段清单（枚数＋每枚来源；交编排者转前端会话，本程不碰前端）

（待填·⑧）

## ⑨ 本程没测什么（逐名；不写"其余都覆盖了"）

（待填·⑨）

## ⑩ 门禁终态（按 `A363` 薄规矩只作现状记录，不充当 AC 结案凭据；读数取在最后一枚 commit 之后）

（待填·⑩）

## ⑪ 被拒／没成功的调用逐条＋有没有跑删除命令＋工具调用终值自报（硬顶 40，第 28 枚停新探索）

（待填·⑪）

## ⑫ next＝落地腿派之前还缺什么（含哪几枚要人先批准）

（待填·⑫）
