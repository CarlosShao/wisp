# C | 票 186 AC#6 五样连带面刷新（只读普查腿 186-a1，HEAD=`b54f1740`）

- 对象＝09-28 普查 §⑥ 那四样"保留、不随树走"，本程逐样现读重跑；**第五样"预算"补上**（普查 §⑥ 已列第 5 行，本程单列以对齐 AC#6 的"五样"口径）。
- 规范：每行＝连带面 ＋ 现读落点/尺 ＋ 读数 ＋ 切换后行为（保留／重建／作废）＋ 会不会留下上一棵树的证据。

## C① 授权根 / `allowed_dirs`（C26）

- 现读落点：`internal/tools/paths_workspace.go`。
  - `:120 SetWorkspaceRoot`（注释 `:109-113` 逐字 "the narrowing must be structurally impossible to widen"）；`:125` 收紧前过 `inRoots`；`:103` 末道复核 `if !p.inRoots(res.Canonical)`；`:142 ClearWorkspace`＝取消收窄的唯一腿；`:152 inRoots`＝不含 workspace 收窄的 allowlist 判定。
  - 尺：`grep -n 'SetWorkspaceRoot\|ClearWorkspace\|inRoots' internal/tools/paths_workspace.go`。
- 读数：候选树必须**已在 `[fs] allowed_dirs` 内**才允许收窄；面板只能**收窄**、⛔ 不能放宽；取消收窄只有内部 `ClearWorkspace()`。
- 切换后行为：**保留**（配置根一字不动，只换 workspace 那一枚字段）。
- 留痕：**不会**——但**能去的树集合由配置决定**，面板今天没有"申请一棵新树"的腿（要做＝动 C26／D45，要人批）。

## C② 产物目录（`agent.Spiller` 的根）

- 现读落点：`cmd/wisp/run.go:1013`＝`ArtifactsDir: filepath.Join(rt.spec.dataDir, "artifacts")`；`cmd/wisp/run.go:1109`＝`Spills: agent.NewSpiller(filepath.Join(rt.spec.dataDir, "artifacts"), loop.Budgets())`；`internal/agent/loop.go:245`＝`sp: NewSpiller(opt.Config.ArtifactsDir, b)`；`internal/agent/spill.go:51` `NewSpiller` 把 `dir` 存成**构造期字段**。
  - 尺：`grep -n 'NewSpiller\|artifacts' cmd/wisp/run.go internal/agent/loop.go internal/agent/spill.go`。
  - ⚠ 行号较 09-28 普查漂移：普查写 `run.go:640`／`loop.go:238`／`spill.go:43`，今日为 **`run.go:1013`／`loop.go:245`／`spill.go:50-51`**。
- 读数：根＝`<dataDir>/artifacts`，**不在工作区根之下**（装配期一次定死）。
- 切换后行为：**保留**（不随树走）。
- 留痕：**会**——上一棵树产生的产物仍躺同一目录、仍可被后续轮读到；要"随树走"须重装配 Spiller（改动面＝装配根＋票 174／`Q-60` 精确豁免）。

## C③ C25 污染名册与 per-scope mark

- 现读落点：`internal/risk/provenance.go:460 func (p *Provenance) OpenScope(scopeID string) *Scope`；`:558 func (p *Provenance) MarkWithHostPath(scopeID, tool, origin, content, hostPath string) bool`；调用侧 `internal/tools/bridge.go:856`＝`b.scopes[taskID] = b.prov.OpenScope(taskID)`、`:650`＝`b.prov.MarkWithHostPath(dec.TaskID, dec.Tool, origin, res.Text, hostPath)`。
  - 尺：`grep -n 'func (p \*Provenance) OpenScope\|func (p \*Provenance) MarkWithHostPath' internal/risk/provenance.go` ＋ `grep -n 'OpenScope\|MarkWithHostPath\|TaskID' internal/tools/bridge.go | head -20`。
- 读数：**名册键＝`dec.TaskID`**（与工作区无绑定）；今天**没有任何代码在切工作区时动 scope**。
- 切换后行为：**保留**（键不认识"树"）。
- 留痕：**会，且是最硬的一维**——同一 taskID 跨树继续跑 ⇒ 旧树污点片段继续命中、新树读作不自开水准；落地须把"换树"绑成"关旧 scope＋开新 scope"（或换 taskID），否则＝票 183/185 那根管子第四次咬人。

## C④ 会话历史与取证行

- 现读落点：历史＝进程内 `internal/agent/loop.go:263 func (l *Loop) History()`、`:1046 func (l *Loop) replaceHistory()`（行号较普查 `:256`／`:1016` 各漂 7／30）。
  - ⚠ **普查说 `internal/session/` 目录只有 `doc.go`——今天已不成立**：`ls -1 internal/session/` ⇒ `doc.go`、`grants.go`、`grants_test.go`、`session.go`、`ticket224_pattern_dialect_test.go` 共 5 枚。
  - 但**语义未变**：`internal/session/session.go:1-40` 头注逐字＝"host-minted **session identity**（ticket 224, approval record A435）"，`Mint()` 用 `crypto/rand`、结束点＝进程退出；`grants.go`＝D45 会话授权行（键 `approval_grant.session_id`）。⇒ **该目录不是会话历史持久化**，"会话历史＝进程内"这半句**仍成立**。
  - 审计行＝`internal/panel/workspace.go:93-103`（失败写 `SWITCH-REFUSED`／成功写 `SWITCH`，带 from/to/spelling/reparse/rewritten）。
- 读数：历史**不进 DB、不清**；审计两形已如实。
- 切换后行为：历史**保留**；审计**如实**（两形都写）。
- 留痕：历史**会**——上一棵树的文件内容与路径继续进上下文；⚠ **接线时不许绕过 `RequestWorkspaceSwitch` 自己写一份审计**（票 186 AC#3(iii)）。

## C⑤ 上下文预算 `BudgetsFor(<窗口>)`（本程补上的第五样）

- 现读落点：`internal/agent/budgets.go:85 func BudgetsFor(ctxWindow int) Budgets`（`:80-84` 头注逐字 "BudgetsFor scales every D15/D39 threshold by ctxWindow"，未知窗口 0 视作参考窗口）；调用侧 `internal/agent/loop.go:222 b := BudgetsFor(window)`（普查写 `:215`，漂 7）。
  - 尺：`grep -n 'func BudgetsFor\|BudgetsFor(' internal/agent/budgets.go internal/agent/loop.go`。
- 读数：预算**键＝模型窗口**，与工作区无关。
- 切换后行为：**保留**——切树不重算、也不需要重算（今天不存在"切树后重算"的代码路径）。
- 留痕：**间接会**——预算没变、历史没清 ⇒ 同一份预算继续吃上一棵树的历史（溢写阈值／spill 判定跟着旧内容走）。

## 一句话代价（与 09-28 普查一致）

- 切换这一发**真正要动的不是面板，是"换树＝换作用域"这条不变式**：②③④⑤ 四样今天**全都不随树走**；⑤的最硬一维＝③（键是 taskID，与树无绑定）；只读显示那一半则**一样都不欠**。
- 与普查相比的**唯一实质变化**：④里 `internal/session/` 目录从"只有 `doc.go`"长到 5 枚件，但长出来的都是**会话身份/授权行**（ticket 224/D45），**不是**会话历史持久化 ⇒ 四条"不随树走"的结论**逐条仍成立**。
