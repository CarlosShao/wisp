# 309-a1 · 20 形ⓐ「入口拒绝」（照 ③ `subagent_197.go:263`）的代价（料，⛔ 裁）

尺与三态同 `10-three-shapes.md`（全部 HEAD blob @`7dc73e56…` 起锚 `7dc753e3`；`git status --porcelain -- internal cmd`＝0 行 ⇒ HEAD blob＝工作树）。⛔ 实测色：本腿零 `go build`/`vet`/`test`，下面每一枚"会红"都是**读那枚钉的断言正文＋它断言的对象**得到的静态判定，静态定不了的写进欠读数。

## 一、会动哪几枚文件（逐枚＋最小编辑形状）

| # | 文件 | 在哪枚函数 | 最小编辑形状 | 必动？ |
|---|---|---|---|---|
| a1 | `internal/tools/bridge.go` | `func (b *Bridge) Execute(...)`，锚 `:268-271`（`if req.CorrelationID == "" { req.CorrelationID = req.TaskID }` 之后）、`:272 entry, ok := b.reg.Lookup` 之前或 `:297` 之前 | 一枚新分支：`if req.TaskID == "" { return b.reject(ctx, req, dec, "<字面量>", string(observe.ClassTool), OutcomeKind<…>) }`。形状参照同文件现成的三枚 `b.reject`（`:275` 未知工具／`:283` 能力／`:289` 参数）⇒ 复用现成出口，⛔ 新造返回类型 | **必动（1 枚）** |
| a2 | `internal/tools/failclosed_236_teeth_test.go` | 无需新增函数；但 `:107` 与 `:244` 两枚用例的**期望值**与 a1 撞车 | 要么把两枚用例改成"从更上游断"（＝把 `lit236*` 期望换成桥的新字面量），要么把载体从桥改成直接调工具（保住 ③／④ 两支仍可达）。**两改必择一，⛔ 留着不管**——它必红 | **必动（第 2 枚）** |
| a3 | `internal/tools/<新>_test.go`（形如 `…_236_legs_test.go` 的一枚补仪器） | 新函数 | 因为 a2 之后 ③(`subagent_197.go:263`)/④(`task.go:709`) 两支**从桥不可达**＝成 dead code，本票族的规矩是"覆盖不许搬家式消失"，所以要另立一枚绕过桥、直接喂空 ctx 的仪器 | **条件必动（第 3 枚）**：只要 a2 选"换成断上游"就必动；选"改载体"则可并入 a2 |
| a4 | 票面 `.scratch/wisp/issues/309-…md` 末尾一条 append-only bullet | — | 票面 `AC#4` 自己写了写面允许项（`internal/tools/**`＋`probes/309/**`＋票尾 1 条）；⛔ 本腿动 | 落地程的事 |

**枚数小结：必动 2 枚，条件必动 1 枚 ⇒ 形ⓐ＝2–3 枚文件**（全在 `internal/tools/**` 内，⛔ 外溢：`OutcomeKind` 常量集就长在 `bridge.go:965-974`，现量 10 枚含 `OutcomeKindError`/`Refused`/`Denied`/`Capability`，ⓐ⛔ 需要新常量，也不需要动 `internal/observe`）。

## 二、会不会打红既有钉（静态点名，逐枚带 file:行＋断言正文）

**必红 2 枚（同一枚文件、两条断言）**

- `internal/tools/failclosed_236_teeth_test.go:257`（属 `Test236R2TaskSpawnRefusesWhenHostGaveNoTaskID`，`:244`）逐字：

```go
	if res.Text != lit236SpawnNoHostTaskID {
		t.Errorf("subagent_197.go:259 那一支的完整字面量没被原样说出（got %q, want %q）——"+
			"若 got 是「子代理 … 已结束」，说明 parentID == \"\" 的拒绝分支已经不在了",
```

  ⇒ 载体逐字在 `:287-288`（`h.bridge.Execute(t.Context(), agent.ToolRequest{` / `TaskID: "", CorrelationID: "", …`）。a1 生效后 `res.Text` 变成**桥的新字面量**，`!= lit236SpawnNoHostTaskID` ⇒ **红**。同函数 `:251`（`!res.IsError`）与 `:269`（`h.roster.Count() != 1`）静态判为**仍绿**（拒绝仍 `IsError:true`、仍不登记）。
- `internal/tools/failclosed_236_teeth_test.go:122`（属 `Test236R2TaskCancelRefusesWhenHostGaveNoCallerID`，`:107`）逐字：

```go
	if callerUnknown.Text != lit236CancelNoHostTaskID {
		t.Errorf("task.go:707 那一支的完整字面量没被原样说出（got %q, want %q）——"+
```

  ⇒ 载体是参数化的：`:118` 逐字 `callerUnknown := x.cancel(t, "", target)`，而 `x.cancel` 本体＝`internal/tools/task_cancel_221_legs_test.go:116-119`（`TaskID: callerID, CorrelationID: callerID`）。空 callerID ⇒ 撞 a1 ⇒ **红**。`:119`（`!callerUnknown.IsError`）、`:113`（派生须成功，那次带 id）静态判为仍绿。
  ⇒ 现量尺：`git grep -n 'cancel(t, ""' HEAD -- internal/tools` ＝ **1 命中**（就是这枚）。

**条件红（实现形状选错才会红）**

- `failclosed_236_teeth_test.go:292-294`（`if err != nil { t.Fatalf("task.spawn dispatch failed: %v", err) }`）与 `task_cancel_221_legs_test.go:123-125`（同一形状）⇒ a1 若实现成**返回 Go error**（⛔ `bridge.go:261-267` 那句 `A REJECT IS NOT A FAULT.` 明令禁止的形状）则两枚用例**由 Errorf 升级为 Fatalf 红**，且面板侧那族 reject 文案仪器也一起进风险区。用 `b.reject` 就不红。

**静态判为「不会红」的（具名说清为什么，⛔ 当作"应该不红"）**

- `internal/tools/grant_test.go:509-518`（`TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority`）与 `internal/tools/ticket90_test.go:454-463`＝**reflect 字段名钉**，射程＝`agent.ToolRequest` 的字段名（禁 `session|grant` 与 `mode|allow|approve|auto`）。a1 是**分支**⛔ 新字段 ⇒ 不红。⚠ 但若有人把ⓐ做成"给 `ToolRequest` 加一枚 `RequireTaskID`/`AssessorRequired` 字段"——那两枚钉的名字规则不禁止这俩拼写 ⇒ 不红，但会**动 `internal/agent` 的契约结构**（`ToolRequest` 是 C17/装配面的一部分），属越界风险，记在这里供裁形时排除。
- `cmd/wisp/panel_inbound_guards_35r3_test.go:65-67` 是这个仓里**真实存在的"文字钉"族**（逐字三枚：`guard35r3Roster = \`不是面板 composer 通路的能力入口\``／`guard35r3Source = \`按伪造/串台拒绝\``／`guard35r3RequestID = \`缺少 requestId\``，注释自指 `bridge.go:133`／`:136-137`／`:140-141`）。它们钉的是 `panel.approval.request` 那三支，与 `assessorFor` 无交集 ⇒ a1 不红；但**它们证明了一件事**：桥里新增一枚字面量＝给这个族递一把新尺，ⓐ 的字面量选词要当契约选（一旦写进去就有人原样钉）。
- `.scratch/wisp/probes/**/*_test.go`（如 `probes/183/r1/mut-freq/zz183freq_test.go:42` 把 `internal/tools/bridge.go` 列进突变名册）：目录以点开头 ⇒ Go 工具链不把它当包 ⇒ ⛔ CI 钉，属验收程的临时件。
- `internal/risk/**` 的钉（含 `72-r1` 刚交的两笔 `f278c57f`/`79e4974d`）：a1 只在 `internal/tools` 加分支，⛔ 叫 `Provenance`；`rules_gateway.go:102` 的 `if ctx.taint == nil {` 语义不变 ⇒ 静态判不红。⚠ 但"整包名册里有没有别的用例依赖空 id 走进过 `assessorFor`"这件事**只有实测能答**（欠读数 U-1）。

**打红的钉枚数小结（静态）＝必红 2 枚（1 枚文件内的两条断言），条件红 2 处（同一枚实现形状错），新文字钉暴露面 +1。**

## 三、用户侧表现

- **今天没人看到。**尺＝`10-three-shapes.md` 的构造点名册：生产 `ToolRequest{` 构造点全仓 **1 枚**＝`internal/agent/loop.go:675`，`TaskID: taskID`，而 `taskID` 只可能来自 `newTaskID()`（`loop.go:1153-1165`，逐字返回 `h[0:8]+"-"+…+"-"+h[20:]`＝恒 36 字符；`:1155` 熵失败那支退化成 `time.Now().UnixNano()` 派生**仍非空**）；`Loop` ⛔ 任何收 caller-supplied task id 的导出方法（尺＝`grep -n "^func (l \*Loop)" logs/loop.go`＝30 枚，只有非导出的 `run(ctx, taskID, input)`(`:346`) 与 `runControl`(`:517`)/`failOpenCalls`(`:896`) 收 taskID）。⇒ a1 那条分支今天**生产不可达**，用户侧＝零变化。
- 谁会先看到：**测试载体**（上表 2 枚）与未来三类候选者（见 `30-structure-guarantee.md`）。
- ⚠ **模型侧可观察的退化（今天不可达，但这是ⓐ的真实代价）**：一旦某条路径真送进空 id，③/④ 那两句教"为什么这件事不能做"的文案（`不知道父任务是谁，就没法登记父子关系，也没法把结论盖戳进父任务的 C25 作用域`／`不知道是谁要停，只能拒`）会被桥的**上游一句"没有任务 id"**遮蔽；模型学到的是"我没身份"，⛔ 再学到"派生要有父、停人要知道是谁"。ⓐ 换来的严格性里含这一笔**教育面损失**，裁形时该计价。
- 审计面：`b.reject` 会落一枚 `tool_call` 行（`dec.OutcomeKind` 选现存常量）；会不会有钉数行/数 kind（欠读数 U-2）。

## 四、⛔ 顺手动到的东西（列出来供派单时逐字带上）

1. ⛔ **删掉 `assessorFor` 的 `taskID == ""` 那一肢**（票面边界＋`provenance.go:694-710` 的 `SrcUnboundScope` 地板；ⓐ 做完之后那一肢变成不可达，仍要**留着＋写为什么留着**，⛔ 顺手"清死码"）。
2. ⛔ 动 `internal/risk/**`（票面「边界与已知禁区」第 1 条；且 `72-r1` 刚交完它的两笔）。
3. ⛔ 顺手把 `b.injected` 那一肢一起收紧——票面 `:27` 明写它归装配那族另裁（现量：`Assessor:` 生产赋值点＝**0 枚**；本腿复核尺＝`git grep -n "injected" HEAD -- internal/tools` → `:132` 字段／`:209` 赋值／`:925` 读，生产面无 `Options.Assessor` 写入点）。
4. ⛔ 为变绿放宽任何断言：**特别注意**——a2 的两种改法里，"把 `lit236*` 期望改成桥的新字面量"就是**改断言**；若走这条路，必须同时新增绕过桥的那枚仪器（a3）把 ③/④ 的覆盖**原样保住**，否则＝覆盖搬家式消失，触票面「⛔ 用 mock 代替真的来假报完成」同族。
5. ⛔ 改 `lit236SpawnNoHostTaskID`／`lit236CancelNoHostTaskID` 两枚常数的**值**（它们是 236 的仪器本体，⛔ 是 309 的橡皮）。
6. ⛔ 给 `agent.ToolRequest`（`internal/agent`）加字段来"携带是否需要 id"（越写面＋撞 reflect 钉的风险类别）。
7. ⛔ 动三枚冻结件／SLO／golden／`thresholds.go`／`d22scan` allowlist／D43 转移表／C1–C32／D1–D47（票面禁区原文）。
8. ⛔ 翻票 309 的框、⛔ 动台账与 `docs/**`（本腿与落地腿同规）。
