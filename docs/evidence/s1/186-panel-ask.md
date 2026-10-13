# 186 面板需求表（AC#9）——聊天框那排的 git 识别与"本地／工作树／分支"切换·Go 侧

- 谁写的：票 186 落地腿 `186-r1`（Go 侧）。**这张表是给编排者转交机主外部 agent（前端会话）的**；本票一字未写 `frontend/**`／`design/**`。
- 锚：`8cfd7b1b`（起手）＋本腿交件笔；文件路径在 `D:\work\workspace\projects plans\Wisp`。行号会漂，引用前现跑尺。
- 前置现量：`186-a1` 普查（笔 `2f41a80c`）＋本腿复跑（`.scratch/wisp/probes/186/r1/logs/01-ac1-ac2-recheck.txt`）。
- ⚠ **一句话**：切换器的**页面本体已经存在**（`frontend/src/components/git-branch.tsx`，201 行 `GitBranchChip`＋`GitBranchUnread`），但今天被 `App.tsx:216-219` 写死成空视图＋空回调（`isRepo:false` ⇒ `return null`，什么都不画、什么都不发）。Go 侧这一腿现在**已经有真实入口**（`panel.workspace.request`，见 §2）；页面要做的只是**接线＋补人话**，⛔ 不是自造路由。
- ⚠ 本表里的**页面侧引文**（`git-branch.tsx` 行数/导出、`App.tsx:216-219`、`git-branch.tsx:59` 的 `return null`）取自 `186-a1` 的证据件 `.scratch/wisp/probes/186/a1/D-frontend-switcher.md`（那次只读普查逐条现跑）；本腿自己只现跑了页面发路由那把尺（`grep -rn 'sendRequest("' frontend/src/lib/*.ts`＝4 枚），**未通读页面**，页面行号引用前请现跑。

---

## 1. Go 侧这一腿交付了什么（页面可以依赖的事实）

| # | 事实 | 落点 | 尺（可重跑） |
|---|---|---|---|
| 1 | **认 git（读面，票 181 已落，本腿复认）**：`.git/HEAD` 逐级回溯，吃目录形/文件形（linked worktree）/裸 sha（detached）三形；非 git 目录**显式一句人话**，不是空 | `internal/panel/git.go` `ReadGit` | `grep -n 'func ReadGit\|GitKindNotARepo' internal/panel/git.go` |
| 2 | **列可去处（读面）**：工作树清单（主树＋`.git/worktrees/*`，路径＋分支）＋本地分支清单（递归 `refs/heads` ＋ `packed-refs` 合并） | 同上 | `grep -n 'func listLocalBranches\|func listLinkedWorktrees' internal/panel/git.go` |
| 3 | **切过去（写腿，本腿新交）**：`panel.workspace.request` → `WorkspaceWriteHandler` → `RequestWorkspaceSwitch`（C26 解析＋102 改写账户＋收窄＋**两形审计**） | `internal/panel/workspace_handlers.go`；装配 `cmd/wisp/panel_resident_windows.go`＋`cmd/wisp/resident_windows.go`＋`cmd/wisp/resident_task_source_windows.go` | `grep -rn 'HandleWorkspaceRequest\|RequestWorkspaceSwitch(h.Scope' internal/panel/workspace_handlers.go` |
| 4 | **同源（AC#3(ii)）**：面板切的是**任务管线那一枚** `tools.PathCanonicalizer`（不是旁边新建一份）——常驻腿用 `WorkspaceScopeHolder` 晚绑定 `run.paths` | `internal/panel/workspace_handlers.go`（holder）＋`cmd/wisp/resident_task_source_windows.go`（`ws.Bind(run.paths)`） | `grep -n 'ws.Bind(run.paths)' cmd/wisp/resident_task_source_windows.go` |
| 5 | **快照字段（票 181 已落）**：`composer.git` 10 枚＋`composer.workspace` 6 枚随快照出 | `internal/panel/composer.go`（`Git GitView`）＋`cmd/wisp/panel_pump.go`（`gitView`/`workspaceView`） | `grep -n 'json:"git"' internal/panel/composer.go` |
| 6 | **旧负判据已按 Q-64 重写**：`TestNoGitSwitchCapabilityInThePanelSurface`（92 AC#7，"面板里不许有 git 切换入口"）退役；换成结构判据 `TestThePanelSwitchCapabilityRidesOneWhitelistedRoute`（能力只许走 C17 那一条＋⛔ 第二扇门） | `internal/panel/composer_test.go` | `grep -n 'func TestThePanelSwitchCapabilityRidesOneWhitelistedRoute' internal/panel/composer_test.go` |

⚠ **`switchBlocked` 那枚字段的语义变了**（本腿改的字面，⛔ 页面按新语义读）：它现在是**"产出这枚包的那个进程能不能切"**，不是"产品能不能切"。`wisp run` 不持有面板宿主 ⇒ 它产出的快照里该字段非空（照实）；持有面板的常驻腿不产出快照 ⇒ 页面今天看不到该字段的"可切"形。改字面的落点＝`internal/panel/git.go` `GitSwitchBlockedReason`。

---

## 2. 页面 ask 表（正式交付物）

### 2.1 控件清单（逐枚：画什么、点下去发什么）

| # | 控件 | 数据来源（快照字段） | 点下去发什么 | 今天缺什么 |
|---|---|---|---|---|
| 1 | **工作区名（聊天框下方那排）** | `composer.workspace.{set,spelling,canonical,reason}` | 已有入口（composer 里的文字框＋"切换"按钮）走 `panel.workspace.request {path}` | 无：**已活** |
| 2 | **当前分支／detached 标记** | `composer.git.{branch,isDetached,detachedSha,repoRoot}` | **只读**：不点、不发；显示 `branch`，detached 时显示 `detachedSha` 短号 | `App.tsx:217` 写死 `{current:"",isRepo:false,...}` ⇒ 整枚 `return null` |
| 3 | **非 git 目录的人话** | `composer.git.{kind,reason}`；`kind` 四值 `repo`／`not_a_repo`／`permission_denied`／`unreadable` | **只读**：`not_a_repo` 时画 `reason` 那句（照抄，⛔ 不许页面另写一句、⛔ 不许留空） | 同上；`git-branch.tsx:59` 今天对 `!isRepo` 直接 `return null`（＝用户以为坏了） |
| 4 | **工作树清单** | `composer.git.worktrees[]`：`{path,branch,sha,detached,main}`（**主树也在内**，`main:true`） | **切到该工作树** ⇒ `panel.workspace.request { path: <worktree.path> }`（目标就是那棵树的目录） | 页面今天的 `worktrees[].{branch,path}` 字段名对得上；⛔ 别把 `branch` 当 `path` 发 |
| 5 | **本地分支清单** | `composer.git.branches[]`（**递归含 `feat/x` 这类斜杠名**） | ⚠ **只有"该分支已被某棵工作树检出"时才能切**：查 `worktrees[]` 里 `branch` 相符的那条的 `path`，发 `panel.workspace.request { path: <那个 path> }`。**没有对应工作树的本地分支**只可显示（见 §3 第 3 条） | 页面无此项判定逻辑；Go 侧不新增"切分支"方法（＝契约变更） |
| 6 | **"新建工作树"** | —— | **今天没有对应方法**：`panel.ts` 里无请求函数，Go 侧白名单里也没有名字 ⇒ 要做＝**新起一枚 C17 方法＝契约变更**（先入档再写，本票⛔ 不做）。今天⛔ 不许页面自造第二条路由，⛔ 不许画成"点了没反应" | 全缺 |
| 7 | **占用条**（owner 截图里那枚） | —— | 不属本票（票 182 那一族） | —— |

### 2.2 信封逐字（⛔ 不许改字段名）

```json
{"method":"panel.workspace.request","requestId":"<页面每发唯一>","source":"panel-composer","path":"<目标目录字符串>"}
```

- `method` 逐字 `panel.workspace.request`（C17 白名单第二枚，`internal/panel/bridge.go:43`）。
- `source` 逐字 `panel-composer`（别的来源按伪造/串台拒，`bridge.go` 自述）。
- `path` 是**用户/页面指的那串**，不是 canonical：C26 负责解析与校验，改写/联接点拒绝按原句回给用户。
- 成功/失败**都不新增页面可见回执**（见 §3 第 4 条）：失败的句子由宿主回给调用方（`window.wispDispatch` 的返回值），今天 `sendRequest` 把它丢掉了。

### 2.3 请求出口只有一条（结构事实，页面照此写）

- 全仓页面今天实际发的请求＝**4 枚**（`panel.mode.request`／`panel.workspace.request`／`panel.attachment.add`／`panel.message.send`；尺 `grep -rn 'sendRequest("' frontend/src/lib/*.ts`）；Go 侧一共答 **6 枚**（多出的两枚是票 248 的 `config.get`／`config.set`，刻意不带 `panel.` 前缀）。本票**没有**新增第 7 枚（尺：`grep -n '"panel\.' internal/panel/bridge.go`＝4 枚 `panel.*`）。
- 结构判据（本腿新增，常驻）：`TestThePanelSwitchCapabilityRidesOneWhitelistedRoute` —— ①能力入口只有 C17 那一条（声明/守卫/路由三处齐）；②包内**任何** `"panel.*"` 字面量都必须是 `bridge.go` 声明的；③收窄腿只有一处（`SetWorkspaceRoot` 的两个非 test 调用点＝`workspace.go` 的 `RequestWorkspaceSwitch`＋holder 的透传）。⇒ 页面自造路由会红在②。

---

## 3. 四条口径（页面必须照抄，⛔ 自行改判）

1. **"切"＝换"以后判定用的范围"，不是 checkout**：Go 侧**不跑 `git` 二进制、不 checkout、不写任何树**（常驻仪器 `TestTheWorkspaceSwitchLegCannotMutateARepo`：往这条腿里种一枚 `exec.Command("git", ...)` 或写文件 ⇒ 红）。所以"切分支"在产品里的含义是**切到该分支已被检出的那棵工作树**；让面板去动编队正在跑的那棵树＝09-27 事故那一形（`HANDOVER §4.0u`／`A328`），⛔ 不做。
2. **不能去的地方**：候选必须已在 `[fs] allowed_dirs` 内（C26 的"只许收窄、不许放宽"）；面板今天**没有**"申请一棵没授权过的目录"的腿（要做＝动 C26／D45，要人批）。
3. **分支/工作树的"可切性"**：`worktrees[]` 里的每一条都可切（服务端保证路径＋分支两样齐）；`branches[]` 里的名字**只有**在 `worktrees[]` 里能找到同 `branch` 的条目时才可切（否则只显示）。页面若要"点了才知道不行"，那枚按钮就是"点了没反应"的变体——⛔ 按票 92/145 的"宁缺毋造"处理（置灰或只在可切时画）。
4. **回显**：切成功后 `composer.workspace`／`composer.git` 会随之变化（同一枚 scope 对象，AC#3(ii)）；但**今天这条路由成功时不回文本**（router 对 workspace/attachment/message 四门一律不写确认句，H10 的活）。页面若要在点击后**立刻**回显"换成了什么"，需要一枚回执或快照推送——本票**未做**，具名登记在 §4 第 2 条。⛔ 页面不许自己拼一句"已切换到 X"。

---

## 4. 已知缺口（具名，⛔ 页面别指望）

1. **`wisp run` 进程没有面板宿主**（也无入向路由装配）；持有面板的常驻腿今天**不把快照推给页面**（出向那一跳＝H10／票 35 的最后一跳）。⇒ 页面今天拿到 git 维的唯一途径还没打通；本表把"数据源字段"与"请求入口"两边都备好了，中间那一跳不在本票。
2. **成功无回执**（§3 第 4 条）：`ComposerDispatch.Handle` 对 workspace 门返回空串；页面 `sendRequest` 走 `postMessage` 且丢弃返回值。要"点完立刻回显"＝H10 或一枚新回执 = 需要人批的接口扩张。
3. **"新建工作树"没有方法**（§2.1 第 6 行）：需要一枚新 C17 名（契约变更），且它必然要跑外部 `git worktree add`——那是一枚**新的宿主能力**（今天全仓零次 `exec.Command("git")`；跑它会咬 `d22scan` ban #1／#4 的管线）。先入档再写。
4. **非 git 目录**：Go 侧已按 AC#4 给"不是仓库"的人话；页面要做的只是**画出来**（今天 `git-branch.tsx:59` 对 `!isRepo` 直接 `return null`）。
