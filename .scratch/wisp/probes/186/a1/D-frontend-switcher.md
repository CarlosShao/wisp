# D | 前端树那个切换器今天长什么样、发什么请求（只读普查腿 186-a1，10-13）

- 范围：只读 `frontend/src/**`（本程一字未写 `frontend/**`／`design/**`）。起点尺：`grep -rn 'workspace.request' frontend/src | head`。
- 规范：每行＝结论 ＋ 尺 ＋ 读数。给票 186 落地腿写"页面 ask 表"用。

## D① 切换入口的组件在哪

- **结论：git 切换入口的组件已经存在，就是 `frontend/src/components/git-branch.tsx`（201 行，自述日期 2026-09-28）。**
  - 尺：`ls -1 frontend/src/components` ⇒ 命中 `git-branch.tsx`；`wc -l frontend/src/components/git-branch.tsx` ⇒ **201**。
  - 两枚导出：`export function GitBranchChip`（`:49`）＋ `export function GitBranchUnread`（`:194`）。
  - 它吃四件视图数据：`GitBranchView { current: string; local: readonly string[]; worktrees: readonly { branch: string; path: string }[]; isRepo: boolean }`（`:32-40`）；两枚回调 `onCheckout(branch)`／`onNewWorktree(branch)`（`:45-46`）。
  - 头注（`:1-24`）逐字钉三事：① 无云端（owner 两次否掉，`云端` 行**刻意缺席**而非置灰）；② 分支是**读来的、不是猜的**（`current` 来自快照，"猜错比没有更糟"）；③ 工作树在本产品是**本地概念**，与本地分支列表并列（`D46` 是 Tier-1 外部命令故事，"worktree＝在另一路径的检出、在新分支上，Go 侧拥有"）。

## D② 它今天挂在哪儿、是不是活的

- **结论：组件挂着的，但视图被写死成空＋两枚回调是空操作 ⇒ 今天它 `return null`、什么都不画、什么请求都不发（死腿）。**
  - 尺：`grep -rn 'git-branch\|GitBranch\|onCheckout\|onNewWorktree' frontend/src | head -20` ⇒ 唯一挂载点＝ `frontend/src/App.tsx:216-219`。
  - 读数（`App.tsx:216-219` 逐字）：
    ```tsx
    <GitBranchChip
      view={{ current: "", isRepo: false, local: [], worktrees: [] }}
      onCheckout={() => undefined}
      onNewWorktree={() => undefined}
    />
    ```
  - 结论链：视图 `isRepo: false` ⇒ `git-branch.tsx:59` `if (!view.isRepo) return null;` ⇒ **整枚组件渲染为空**；两枚回调是 `() => undefined` ⇒ 即便画出来也发不出请求。
  - `grep -rn 'gitBranchView\|GitBranchView\|branches:\|worktrees:' frontend/src` ⇒ 除组件自身与 `App.tsx:217` 那枚写死空视图外，**没有任何数据源**接进来（`frontend/src/lib/panel.ts` 里没有 git 维）。

## D③ 今天真正"发工作区请求"的入口是哪一枚

- **结论：唯一活着的发送入口是 composer 里那枚工作区文字框＋"切换"按钮，走 `panel.workspace.request`；git 切换器这一条今天没有请求。**
  - 尺：`grep -rn 'workspace.request' frontend/src | head` ⇒ **2 命中**：
    - `frontend/src/lib/panel.ts:62`（注释：本面板只暴露 `panel.mode.request` / `panel.workspace.request` / `panel.attachment.add` 三枚请求，在**同一枚信封**上）；
    - `frontend/src/lib/panel.ts:255`＝`sendRequest("panel.workspace.request", { path });`——定义在 `requestWorkspaceChange(path)`（`:250-256`，头注逐字："The string is only a request: C26 resolves it, reparse/rewrite refusals come back as a refusal the panel displays, and the audit line is written natively"）。
  - 调用点：`frontend/src/components/composer.tsx:172`＝`requestWorkspaceChange(workspaceDraft.trim());`（按钮文案"切换"，旁边是一枚 `aria-label="工作区路径（由原生侧解析与校验）"` 的 mono 输入框，`:163-175`）。
  - 对照枚：`requestModeSwitch(to)` → `panel.mode.request`（`panel.ts:245-247`）＝同族"另一枚"。全仓页面可发的请求＝4 枚（`mode.request`／`workspace.request`／`attachment.add`／`message.send`）。
  - ⚠ 与 `git-branch.tsx` 无关：`App.tsx:216-219` 那枚 `onCheckout` 若被接上，**今天也没有对应的请求函数**——`panel.ts` 里没有 `requestCheckout`／`requestNewWorktree` 一类。

## D④ 页面 ask 表（递给票 186 落地腿，本程不代写前端）

给"落地腿写页面 ask 表"用，逐枚列页面**今天缺什么**与**能发什么**：

| # | 页面要的东西 | 今天状态（尺/读数） | 落地腿要供的 |
|---|---|---|---|
| 1 | git 视图数据源：`current` / `local[]` / `worktrees[{branch,path}]` / `isRepo` | `App.tsx:217` 写死 `{current:"", isRepo:false, local:[], worktrees:[]}`；无旁的数据源 | 快照里加 git 一节，字段拼齐这 4 枚（普查 §⑧ 建议再带 `kind`/`reason`/`detachedSha`/`repoRoot`） |
| 2 | 非 git 目录那一形 | `git-branch.tsx:59` `if (!view.isRepo) return null;`（渲染空，无"不是仓库"人话） | 票 186 AC#4：要显式一句人话理由（⛔ 不留空） |
| 3 | "切分支／切本地／切工作树"请求 | **已有**：`requestWorkspaceChange(path)` → `sendRequest("panel.workspace.request", {path})`（`panel.ts:250-256`），今天由 composer 文字框发 | 落地腿可**复用这一枚**（在册、零新白名单枚）把目标当 `path` 送去 |
| 4 | "新工作树"请求 | **无**：`panel.ts` 里没有对应请求函数；`App.tsx:219` 的 `onNewWorktree` 是空操作 | 若要做＝新方法名＝**契约变更**（AC#2：先落 `A##` 再写），⛔ 不许页面自造第二条通道 |
| 5 | `panel.workspace.request` 的信封字段 | `ComposerRequest` 走 `path` 字符串（`panel.ts:255`）；`ComposerWorkspace`（`:80-92`）＝6 枚字段 `set`/`spelling`/`canonical`/`reparse`/`rewritten`/`reason`，**无 git 维** | 切分支/工作树时目标名塞 `path` 还是新字段，属快照/信封的载体扩张（普查 §⑫(b)：要不要动规格文字请编排者裁） |

- **一句话**：前端**已经把 git 切换器的 UI 带进来了**（`git-branch.tsx`，含本地分支列表＋工作树列表＋新建工作树输入），但它今天**没接数据、没接请求**（写死空视图＋空回调）⇒ 票 186 落地腿的"页面 ask"＝**给它供 4 枚视图字段＋一枚 `panel.workspace.request` 的去路**；"新工作树"那一支今天**无对应方法**。
