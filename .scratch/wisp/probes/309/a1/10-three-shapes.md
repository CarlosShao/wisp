# 309-a1 · 10 三处（实为四处）同形处置的逐字名册（`AC#0` 的料）

## 尺与三态（本文件全部读数共用）

| 尺名 | 命令 | 射程 |
|---|---|---|
| blob 落盘 | `git show HEAD:<path> > logs/<base>.go` | 只读 HEAD blob；工作树不参与 |
| 名册 | `git grep -n "assessorFor" HEAD -- internal cmd` | 全 `internal`＋`cmd`，含测试 |
| 门路 | `git grep -n "TaintHit\|IsSyncPath" HEAD -- internal` | 含测试（逐条看过，非测试命中 5 枚） |

- **三态＝HEAD blob**（锚 `7dc753e3`）。起手实测 `git status --porcelain -- internal cmd` ＝ **0 行** ⇒ HEAD blob＝工作树逐字节，两态在此刻重合；本文件一律写 HEAD blob。
- 原始件：`logs/bridge.go`(1284 行)·`logs/panel_assets.go`(248)·`logs/subagent_197.go`(582)·`logs/provenance.go`(1115)·`logs/failclosed_236_teeth_test.go`(296)·`logs/loop.go`·`logs/assessor.go`。
- 票面三个行号 `bridge.go:925`／`panel_assets.go:164`／`subagent_197.go:263` ⇒ **现量全部对得上**（见下逐条）。

---

## ① `internal/tools/bridge.go:924-932` ＝ 静默摘掉 R4 探测器，零理由注释

逐字（尺＝blob 落盘后 `sed -n '920,932p'`）：

```go
// assessorFor returns the assessor whose R4 input is bound to this task's
// scope. The frozen C19 TaintDetector seam takes no tool name
// (internal/risk/assessor.go:165-171), so the binding is per scope, hence per
// task; the canonicalizer and classifier are shared and immutable.
func (b *Bridge) assessorFor(taskID string) *risk.RiskAssessor {
	if b.injected || b.prov == nil || taskID == "" {
		return b.assess
	}
	return risk.NewRiskAssessor().
		WithCanonicalizer(b.paths).
		WithSensitiveClassifier(b.classifier).
		WithTaintDetector(b.prov.Detector(taskID))
}
```

- 票面锚 `:925`＝现量 `:925`（`if b.injected || b.prov == nil || taskID == "" {`）；票面 `:931 WithTaintDetector(...)`＝现量 `:931`；文档块 `:920-923`＝现量 `:920-923`，**逐字只讲"按 scope 绑定"，对 `taskID == ""` 那一肢一字未提** ⇒ 票面"零理由注释"成立。
- **`b.assess` 里本来就没有探测器**（尺＝`sed -n '208,214p'` blob）：

```go
	if o.Assessor != nil {
		b.assess, b.injected = o.Assessor, true
	} else {
		b.assess = risk.NewRiskAssessor().
			WithCanonicalizer(o.Paths).
			WithSensitiveClassifier(b.classifier)
	}
```

  ⇒ 走 `:926` 返回的就是**没挂 `WithTaintDetector` 的那枚**，"摘掉"是**结构事实**而非措辞。`b.injected` 字段注释逐字（`:132`）＝`injected bool // assess came from Options.Assessor: do not rebuild per call`。

### 它守的是哪一格的凭据

- **谁都不守。**名册尺（`git grep -n "assessorFor" HEAD -- internal cmd`）现量 **7 命中**：真码 2 枚＝定义 `bridge.go:924` ＋唯一调用点 `bridge.go:297`；注释 5 枚＝`bridge.go:920`·`internal/risk/assessor.go:30`·`cmd/wisp/panel_assets.go:60`·`cmd/wisp/panel_assets.go:227`·`cmd/wisp/panel_assets_143_test.go:7`。
  ⇒ **点名 `assessorFor` 的用例＝0 枚**（唯一测试命中是注释行）＝票面成立；但票面 ①段"其余 4 枚命中全是注释"**枚数少算 1**：滤掉测试后是 5 枚（它把 `:920` 那行文档注释与 `:924` 定义合并成了"定义自己"）。不影响任何判定。
- **被经过、未被观测**（票面口径核实）：尺＝blob `sed -n '287,291p' logs/failclosed_236_teeth_test.go`，现量逐字：

```go
	res, err := h.bridge.Execute(t.Context(), agent.ToolRequest{
		TaskID: "", CorrelationID: "", CallID: "call-236r2-spawn-no-parent",
		Name: "task.spawn",
```

  其断言三句（`:251` `if !res.IsError`／`:257` `if res.Text != lit236SpawnNoHostTaskID`／`:269` `if n := h.roster.Count(); n != 1`）**没有一行讲探测器挂⛔ 挂** ⇒ 票面"被经过、未被观测"成立，⛔ 写"零仪器经过"这句也成立（那一说过头）。

### 它到门那一跳有几步

- **0 跳**：`:925` 那道守卫就长在门自己的调用表达式上（`bridge.go:297` 逐字 `verdict := b.assessorFor(req.TaskID).Assess(req.Name, params,`），守卫与判级是**同一行的两半**。
- 走得到 `Gate`：`:297` 出 `verdict` → `:310-311` `mode := b.permissionMode()` / `sil := mode.Screen(verdict)` → `:318` 起 `switch` 派路由 ⇒ 摘探测器后**门照走、判级照出**，只是少一枚 R4 输入。
- **⛔ 走不到 `IsSyncPath`**：`internal/risk/rules_gateway.go:102` 逐字 `if ctx.taint == nil {`、`:105` `source, hit := ctx.taint.TaintHit(ctx.params)` ⇒ 探测器为 nil 时 R4 短路，`Provenance.Inspect` 根本不被叫，`writeGate`(`provenance.go:908`)→`IsSyncPath`(`provenance.go:928`) 那一跳整条不发生。这就是"少报"的机械形状。
- ⚠ 一处协同（AC#0 裁形要用）：`TaskID(ctx)`(`internal/tools/cancel.go:67-70`) 读的就是桥盖进 ctx 的那枚 id，而 `:297` 读的是 `req.TaskID`——**同一枚空串会同时触发 ① 的静默与 ③／④ 的拒绝**，两件事发生在同一次调用里，只是 ① 在上游、③ 在下游。

---

## ② `cmd/wisp/panel_assets.go:162-174`（＋`:224-240`）＝ 自己命名一枚并写明故意

逐字（尺＝blob `sed -n '162,174p' logs/panel_assets.go`；票面锚 `:164`＝现量 `:164`）：

```go
// taintSourceScopeID is the C25 scope this command opens for its one card. In
// the product the scope id is a task id owned by the agent loop (ticket 19's
// DEFERRED(C25-loop-wiring) item 1); a CLI probe has no task, so it names one
// and closes nothing - the engine lives and dies inside this process.
//
// Ticket 160 left that behavior exactly as written and changed what the code
// has to admit: OpenScope now hands back the *Scope that owns the close, and
// this call throws it away on purpose. What can no longer be written here is
// closing somebody ELSE's scope - there is no by-id close left to reach for.
// "Nobody ever closed this one" stays G5's job in
// .scratch/wisp/probes/154/gate-clauses.sh, not the compiler's, which lets the
// discard through (measured: docs/evidence/s1/160-handle-r1.md §3).
const taintSourceScopeID = "panel-assets-l2"
```

配套的"写明故意"第二处，逐字（`:224-234`）：

```go
// detector builds the C25 engine over the declared sources and returns it bound
// to this command's scope - the same call shape the production bridge uses
// (risk.NewProvenance(...).Detector(taskID), see internal/tools/bridge.go's
// assessorFor). Returns nil when no source was declared, which is the caller's
// signal to wire nothing at all.
//
// ProvOptions.NoProbe is the one difference from the bridge, and it cannot
// change a verdict on this path: the panel builds its parameters from the
// argv (internal/panel's paramsFromArgs yields the "command" and "argv" keys
// only), neither of which is a write target, so the sync-directory gate this
// option would feed is never consulted for a call of that shape.
```

以及装配现场逐字（`:57-78`，节选）＝`assessor := risk.NewRiskAssessor()` / `if det := taintSources.detector(); det != nil {` / `assessor = assessor.WithTaintDetector(det)`，且 `panel.ApprovalSubject{CorrelationID: "panel-assets-l2", ...}`（`:69`，复用同一枚自命名串）。
⇒ 票面 ②"自己命名一枚并写明故意"**成立**，且"写明"不是一句：它同时写了**故意丢掉句柄**（`throws it away on purpose`／`Discarded on purpose, and named so`，`:240`）与**谁该为此响**（`G5's job in .scratch/wisp/probes/154/gate-clauses.sh`）。这是本名册里"写明故意"的**成色基准**。

### 它守的是哪一格的凭据

- `cmd/wisp/panel_assets_143_test.go:7` 的注释逐字含 `assessorFor`，讲的是"CLI 路与生产路必须同一枚 assessor"这条**同一性**（票 143 的遗产）。现量该文件对 ②的覆盖度＝**欠读数**（本腿⛔ 跑测试，只能静态说：它点名的是 assessor 装配形状，不是 `taintSourceScopeID` 的命名故意）。
- 它**不守任何 R4 判级**：它是 `wisp panel-assets -l2` 那枚 CLI 探针（`:41` 逐字 `-l2` 用法串），产物是一枚卡片 JSON，不落审计、不过门。

### 它到门那一跳有几步

- **到不了桥的门**：本路径⛔ 经过 `Bridge.Execute`（`bridge.go:297` 那道门与 `:311` 的 `mode.Screen` 都不发生）；它自己 new 一枚 assessor（`:64`）后交给 `panel.NewApprovalCardView(assessor, ...)`（`:68`）→ 卡片视图内部调 `Assess`。⇒ **0 跳到桥门／1 跳到 R4**。
- 走得到 `IsSyncPath` 的**代码**、⛔ 走不到它的**判定**：探测器＝`p.Detector(taintSourceScopeID)`（`provenance.go:858-860` → `boundDetector.TaintHit` `:867-868` → `Inspect` `:725`）→ `writeGate`(`:908`)，但 `ProvOptions.NoProbe`（`:239` 装配＋`:230-234` 说明）＋参数形状无 path ⇒ `:914 !hasPath` 那支拿到 `ChUnknown` 泛扫，`:927` 的 `IsSyncPath` 循环一次不进。**这段是"写明故意"如何写得可核的现成范本：说出差异、说出为什么它改不动判级、并点名该响的那枚门（G5）。**

---

## ③ `internal/tools/subagent_197.go:262-266` ＝ fail-closed 拒绝

逐字（尺＝blob `sed -n '259,266p'`；票面锚 `:263`＝现量 `:263`）：

```go
	// The task-level id of the caller: since ticket 242 the correlation id is
	// minted per call (C18 routes approval replies by it), so it can no longer
	// name a task; the roster is keyed by task ids.
	parentID := TaskID(ctx)
	if parentID == "" {
		return Result{Text: "这条调用没有宿主给的任务 id（fail-closed：不知道父任务是谁，" +
			"就没法登记父子关系，也没法把结论盖戳进父任务的 C25 作用域）", IsError: true}, nil
	}
```

⇒ 票面 ③"fail-closed 拒绝＋文案逐字"**成立**（票面省略号之后的全句已抄）。注意 `err == nil`＋`IsError: true`＝`bridge.go:261-267` 定死的"拒绝不是故障"形状（逐字 `A REJECT IS NOT A FAULT.`）。

### 它守的是哪一格的凭据（仪器真实在不在）

- 仪器**真实在**：尺＝blob `logs/failclosed_236_teeth_test.go`。`internal/tools/failclosed_236_teeth_test.go:244` 逐字 `func Test236R2TaskSpawnRefusesWhenHostGaveNoTaskID(t *testing.T) {`（票面 `:244`＝现量 `:244`），头注 `:240` 逐字 `// AC#1b, subagent_197.go:259-263.`。三句断言正文逐字见上一格（`:251`/`:257`/`:269`），字面量常数在 `:98-99`：

```go
	lit236SpawnNoHostTaskID = "这条调用没有宿主给的任务 id（fail-closed：不知道父任务是谁，" +
```

- ⚠ 头注里的行号已漂：它写 `subagent_197.go:259-263`／`task.go:707-709`／`:253-255`，现量真身分别是 `:262-266`／`:708-711`／`:253-255`（第三枚没漂）。**属 `AC#3` 那族"注释行号腐烂"，⛔ 本腿改，只登记。**

### 它到门那一跳有几步

- **它在门的下游**：`task.spawn` 的这次调用先过 `bridge.go:297`（判级，同一刻 ① 的守卫也因 `req.TaskID==""` 而摘掉探测器）→ `:311` `mode.Screen` → 路由 → 才进 `subagentSpawn.Execute`。⇒ **到门 0 跳（改不动任何判级）、到 `IsSyncPath` 0 跳**（`task.spawn` 无 path 参数）。它的"拒绝"是**工具自报的 `Result`**，⛔ 是门给的 DENY。
- 同族还有一支在 `subagent_197.go:474-477`（`case parentID == ""` 的盖戳侧，尺＝`grep -n "case parentID" logs/subagent_197.go`），形状＝**switch 里单独一臂拒**，也是下游。

---

## ④ 名册补一枚（票面未列；同形、同为 fail-closed 拒绝，现量在仓）

逐字（尺＝blob `git show HEAD:internal/tools/task.go | sed -n '708,711p'`）：

```go
	caller := TaskID(ctx)
	if caller == "" {
		return Result{Text: "这条调用没有宿主给的任务 id（fail-closed：不知道是谁要停，只能拒）", IsError: true}, nil
	}
```

- 仪器**真实在**：`internal/tools/failclosed_236_teeth_test.go:107` `func Test236R2TaskCancelRefusesWhenHostGaveNoCallerID(t *testing.T)`，断言 `:119`（`!callerUnknown.IsError` → Fatalf）与 `:122`（`callerUnknown.Text != lit236CancelNoHostTaskID`，常数逐字在 `:91`）。
- 载体是**参数化的桥入口**（尺＝`git grep -n -A4 "func (x \*harness221) cancel(" HEAD -- internal/tools`）＝`internal/tools/task_cancel_221_legs_test.go:116-119`，逐字 `func (x *harness221) cancel(t *testing.T, callerID, targetID string) Result {` … `TaskID: callerID, CorrelationID: callerID,` ⇒ 调用点传 `""` 就从桥送进空 id，**这是 `ⓐ` 形第二处会撞上的仪器**（见 `20-shape-a-cost.md`）。
- 结论性事实（料，⛔ 裁）：**"宿主没给任务 id"在 `internal/tools` 里的多数已经倾向于拒绝**（③＋④两枚工具、两枚仪器、两枚字面量），静默摘探测器只有 ① 一枚，且它⛔ 仪器。

---

## 与本腿简报转述不一致处（具名报回，⛔ 据此改判任何东西）

1. 简报（与票面 `:24` 同文）说 `provenance.go:727` 对未开 scope **"无条件"** 回 `Hit{… SrcTool: SrcUnboundScope}, true`。现量：`:727` 那句本身是的（逐字 `if unbound {` → `:729` `SrcTool: SrcUnboundScope,` → `:731` `}, true`），**但 `unbound` 这个旗子本身是有条件的**，尺＝blob `sed -n '694,710p' logs/provenance.go`：

```go
func (p *Provenance) scopeMarks(scopeID string) ([]*taintMark, bool) {
	reg := p.scopes[scopeID]
	if reg != nil && len(reg.marks) > 0 {
		return reg.marks[:...], false
	}
	if reg != nil {
		return nil, false // opened and empty: this task read nothing tainted
	}
	for _, r := range p.scopes {
		if len(r.marks) > 0 {
			logf("risk/C25: Inspect/CheckText on unregistered scope %q while other scopes hold taints: fail-closed R4 …")
			return nil, true
		}
	}
	return nil, false
}
```

  ⇒ 摘掉 `taskID == ""` 那一肢后的"误报风暴"**只在"引擎里已存在别的 scope 带着污点"时**发生；一枚从没读过敏感源的任务撞的是 `:702`/`:710` 的 `false`，不报。**地板真实存在，但它的射程比"每一次写都报"窄一档**——这条不影响"⛔ 删那一肢"的边界，只影响该边界被写成什么样的措辞。（另：`CheckText` 侧 `:826-830` 同一形状，逐字 `marks, unbound := p.scopeMarks(scopeID)` / `if unbound {`。）
2. 票面 `:15` 说消费者面"其余 4 枚命中全是注释"，现量滤测试后是 **5 枚**（多出的是定义自己的文档行 `:920`）。内容判定不变。
3. 票面 `:20`/简报把 ③ 的行号写成一枚（`:263`）；现量 `:262` 是取值、`:263` 是判空，判定行号**未漂**。
