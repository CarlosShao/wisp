# taskid-a1｜只读普查：`bridge.go:925` 的 `taskID == ""` 真跑时会不会发生

- 腿：`taskid-a1`（只读普查）／分支 `dev`／锚点 `9c83598d`
- 欠读数出处：`docs/reports/pending-and-issues.md:15981` 节头 `## A838｜2026-10-11 10:5x｜编排者**结掉票 307 登记②那枚欠读数**`；
  该节第 33 行（绝对 `16013`）逐字：`★**新抓一枚具名欠读数**＝`bridge.go:925` 那三枚前提 `if b.injected || b.prov == nil || taskID == ""` 里，**真跑时 `taskID==""` 会不会发生**`
- ★ 与派单转述的冲突：派单写「成因与已结案票 307⛔ 关」，仓里原文写「**成因与 307 无关**」（同一句）。以仓里为准，两者一致。
- 行号＝本读数那一刻的快照；每处配内容锚。原始 grep 落 `logs/`。
- ⛔ 跑 `go build`／`go vet`／`go test`（编译面归写腿）。本轮未跑 `go env`／`go list`＝零次。

---

## 六跳链（派单要求自核，⛔ 照抄）——全部对得上

| 跳 | `file:行` | 内容锚（逐字） |
|---|---|---|
| 1 | `internal/tools/bridge.go:925` | `if b.injected \|\| b.prov == nil \|\| taskID == "" {` |
| 2 | `internal/tools/bridge.go:931` | `WithTaintDetector(b.prov.Detector(taskID))` |
| 3 | `internal/risk/provenance.go:858` | `func (p *Provenance) Detector(scopeID string) TaintDetector {` → `:859 return boundDetector{p, scopeID}` |
| 4 | `internal/risk/assessor.go:168` | `TaintHit(params map[string]any) (source string, hit bool)` |
| 5 | `internal/risk/rules_gateway.go:105` | `source, hit := ctx.taint.TaintHit(ctx.params)` |
| 6 | 不响的根因 | `internal/tools/bridge.go:211-213`（`b.assess` 只挂 `WithCanonicalizer`＋`WithSensitiveClassifier`，**无 detector**） |

**新拓到的一跳（派单没写，但它是 ⓑ 的地板）**：`internal/risk/provenance.go:727-731`
`if unbound {` → `SrcTool: SrcUnboundScope,` / `Origin: "scope is not open (no OpenScope for this id, or its own handle closed it)",` → `}, true`
⇒ **风险层对未开 scope 是 fail-CLOSED（无条件 hit＝会误报）**。`:925` 那枚守卫正是挡住这一支的唯一体。

---

## 问一：`taskID` 从哪来／有几枚 caller

尺＝`grep -rn 'assessorFor' --include=*.go .`（全仓，73 行命中，`logs/assessorfor-callers.txt`；**其中真源码只有 3 枚**，其余 70 行＝`.scratch/**` 里的探针拷贝，⛔ 仪器）＋`grep -rn 'ToolRequest' --include=*.go internal cmd`（非测试，`logs/taskid-all.txt` 同族）＋逐枚 `sed` 读函数体。

**名册（全名）**

1. `assessorFor` 的 caller＝**1 枚**：`internal/tools/bridge.go:297` `verdict := b.assessorFor(req.TaskID).Assess(req.Name, params,`
   所在函数＝`bridge.go:268` `func (b *Bridge) Execute(ctx context.Context, req agent.ToolRequest) (agent.ToolOutcome, error) {`
2. `Bridge.Execute` 的生产 caller＝**2＋1 枚**：
   - `internal/agent/loop.go:760` `return l.opt.Tools.Execute(ctx, req)`
   - `internal/agent/loop.go:764` `out, err := l.opt.Tools.Execute(tctx, req)`（两枚都在 `l.dispatch`，签名 `loop.go:757` `req ToolRequest, timeout time.Duration,`）
   - `internal/tools/subagent_197.go:581` `return p.inner.Execute(ctx, req)`＝**纯透传装饰器**（`:573 func (p *subagentToolProvider) Execute(...)`，只摘 `task.spawn` 一枚名，⛔ 碰 `req.TaskID`）
3. `agent.ToolRequest` 的构造点＝**全仓生产码只有 1 枚**：`internal/agent/loop.go:675` `req := ToolRequest{` → `:676` `TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i),`

**实参 `taskID` 的源头（收敛到唯一铸造口）**

- `loop.go:346` `func (l *Loop) run(ctx context.Context, taskID, input string) Result {`
- 入口只有 2 枚：`loop.go:333` `t.result = l.run(c, id, input)`（`id` 来自 `:329` `id := newTaskID()`）／`loop.go:340` `return l.run(ctx, newTaskID(), input)`
- `loop.go:1153` `func newTaskID() string {` → `:1165` `return h[0:8] + "-" + h[8:12] + "-" + ...`＝**恒 36 字符**。熵源失败那支（`:1156` `n := time.Now().UnixNano()`）退化但**仍⛔ 为空**。
- `grep -rn 'func (l \*Loop)' internal/agent | grep -i taskid`＝命中 4 枚**全非导出**（`run:346`／`runControl:517`／`stream:552`／`failOpenCalls:896`）
  ⇒ 派单引的那句注释（`bridge.go:886-888` `Loop has no exported method that takes a caller-supplied task id`）**本轮独立证真**，⛔ 追认。

**逐支判定（派单点名的五支）**

| 支 | 走⛔ 走 `assessorFor` | 空串条件 | 凭据 |
|---|---|---|---|
| 根任务 | 走 | **⛔ 不可能**：`newTaskID()` 恒非空 | `loop.go:329`/`:340` 内容锚见上 |
| 子任务（子代理） | 走（经 `:581` 透传） | **⛔ 不可能**：孩子自己一轮 loop 另铸 | `subagent_197.go:347` `bg := child.RunAsync(childCtx, prompt)`＋`:332` 注释「child loop's task id is minted by the loop itself」 |
| CLI `wisp run` | 走 | **⛔ 不可能**：装配根只交 `AdmitTask`，⛔ 手搭 `ToolRequest` | `cmd/wisp/run.go:1008` `AdmitTask: rt.admitTask,` |
| 面板发起 | **⛔ 根本不经过 `assessorFor`** | 不适用：另建一枚 assessor，scope 名硬编码非空 | `cmd/wisp/panel_assets.go:64` `assessor := risk.NewRiskAssessor()`＋`:174` `const taintSourceScopeID = "panel-assets-l2"` |
| 定时任务 | **⛔ 路径不存在** | 不适用：包是 doc-only 桩 | `internal/agent/scheduler/doc.go:13` `// DEFERRED(scheduler): implemented by ticket 47.` |

**今天唯一产出空串的通道＝测试自己**（见问三）。

---

## 问二：那枚守卫是故意的还是没人想过

尺＝`grep -n 'taskID\|TaskID\|no task\|empty task' internal/tools/bridge.go | grep -E ':\s*//'`（注释面，全命中只有 2 行）＋`grep -rn 'no task|taskID == ""' docs/evidence/s1 .scratch/wisp/issues`＋`git log -S`。

- **`:925` 空串那一肢，仓里⛔ 没有任何一处注释给它理由。**
  `assessorFor` 自己的文档块（`:920-923`，锚 `// assessorFor returns the assessor whose R4 input is bound to this task's` / `// task; the canonicalizer and classifier are shared and immutable.`）
  只解释「绑定为何按 scope」，⛔ 字面提到空 id。相邻两枚同形守卫也一致沉默：`:848`（`OpenTask`）/`:898`（`CloseTask`）都是 `if b.prov == nil || taskID == "" {`。
- **但它⛔ 是随手写的**：`git log --oneline -S 'taskID == ""' -- internal/tools/bridge.go`＝**唯一一枚** `64c5fea3 feat(20): the host bridge as the loop's ToolProvider, plus the fs.read/fs.list L0 pair`；
  `git show 64c5fea3:internal/tools/bridge.go` 里 `:492/:509/:526` 三枚守卫**与 `:526` 那枚 `WithTaintDetector(b.prov...` 同行段同时落地**
  ⇒ 守卫与检测器是同一次提交一起长出来的，形状稳定至今。
- **可⛔ 不等于「写明故意」**：仓里**写明故意**的处理长在另外两处，且方向相反——
  - `cmd/wisp/panel_assets.go:164` 逐字：`// DEFERRED(C25-loop-wiring) item 1); a CLI probe has no task, so it names one` ＋ `:165` `// and closes nothing - the engine lives and dies inside this process.`（无任务 ⇒ **现造一枚名字**，`docs/evidence/s1/160-handle-r1.md:165` 认它是「**写明故意的丢弃**」）
  - `internal/tools/subagent_197.go` 那一支：空父 id ⇒ **拒绝**（凭据 `failclosed_236_teeth_test.go:258-259` 锚 `subagent_197.go:259 那一支的完整字面量没被原样说出`）
  ⇒ **同一条「宿主没给任务 id」的形状，在 `task.spawn` 处是拒绝、在面板探针处是现造名字，唯独在 `:925` 处是静默降级；这⛔ 解释⛔ 记录。**
- 综合：读作**「有功能必要、但从未被写成前提」**最贴——必要性见下面 ⓑ 的地板（摘掉它 ⇒ 全量误报），⛔ 等于「想过并决定这样降级」。

---

## 问三：今天有没有仪器看得见那一支

尺＝`grep -rn 'assessorFor' --include=*_test.go internal cmd`＋`grep -rn 'TaskID: *""' --include=*_test.go internal cmd`＋`grep -rn 'WithTaintDetector\|\.Detector(' --include=*_test.go internal cmd`。

- **点名 `assessorFor` 的用例＝0 枚。** 尺只回 1 行，且它是一条**注释**⛔ 断言：
  `cmd/wisp/panel_assets_143_test.go:7` 锚 `// assessorFor), so a demo machine and a real machine judged by different rules.`
- **但「空 `taskID` 进 `Execute`」这一支今天确实被走到**（派单预设的「零仪器经过」⛔ 成立）：
  - `internal/tools/failclosed_236_teeth_test.go:280` `func spawnWithoutHostTaskID(t *testing.T, h *sub197Harness) (out Result) {`
  - `:287-288` `res, err := h.bridge.Execute(t.Context(), agent.ToolRequest{` / `TaskID: "", CorrelationID: "", CallID: "call-236r2-spawn-no-parent",` → `:290` `Name: "task.spawn",`
  - 由 `:244` `func Test236R2TaskSpawnRefusesWhenHostGaveNoTaskID(t *testing.T) {` 在 `:250` `res := spawnWithoutHostTaskID(t, h)` 调起。
  - **它断言的全是别的东西**：`:251/:255` `res.IsError`＋拒绝字面量、`:269` `if n := h.roster.Count(); n != 1 {`。
    **没一行**说「检测器没挂」或「挂了」；且 `task.spawn` ⛔ 是路径写，R4 在那发根本没上场。
  ⇒ 准确口径＝**「这支被经过、⛔ 被观测」**。摘掉 `:925` 的 `|| taskID == ""` 会不会让哪枚用例变红＝**静态判⛔ 出，需真跑**（本轮禁编译面）＝**编排者欠，⛔ 算我欠**。
- 风险层那一侧：**⛔ 一枚用例把空 scope 喂进 `Detector`**。既有仪器全用非空名：
  `internal/risk/provenance_test.go:233` `a := NewRiskAssessor().WithTaintDetector(p.Detector("task-1"))`；
  `:368` `if src, h := p.Detector("ghost").TaintHit(...); !h || !strings.Contains(src, SrcUnboundScope)`——**最接近的一枚**，验的是「未知名＝fail-closed 会 hit」，`""` 这一枚名字从未出现。

---

## 判语（两格）

### ⓐ 真跑时 `taskID == ""` 会不会发生

**不会——今天的接线面⛔ 发生；且⛔ 是「只在具名那一支会」，因为那⛔ 存在。**

- 根任务／子任务／CLI `wisp run` 三支全部收敛到 `newTaskID()`（`loop.go:1153`），它恒回 36 字符（`:1165`），熵源失败退化成时钟派生**仍⛔ 空**。
- 面板侧⛔ 经过 `assessorFor`（`panel_assets.go:64` 自建 assessor）。定时任务包是 doc-only 桩（`scheduler/doc.go:13` `DEFERRED(scheduler): implemented by ticket 47.`）。
- `ToolRequest` 生产构造点全仓**只有 1 枚**（`loop.go:675`）⇒ 空 id 在今天的生产码里**没有出生地**。
- ★ 但「⛔ 出生地」是**今天的接线事实**，⛔ 是结构保证：`Loop` 没有吃调用方 id 的导出口（本轮已独立证真，见问一），所以任何想自带 id 上桥的宿主**必须**手搭 `ToolRequest`——那枚守卫就是给它准备的，而它一进来会拿到⛔ 响的东西。票 154 那段（`:872-875` 锚 `does the TaskID you dispatch with` / `a panel or ball host inventing its own id, a retry wrapper, a scheduled task carrying one id across turns`）已经点名了这三形状的候选者，并写着 `:876-877` `Today no production code dispatches on the bridge except loop.go, so being such a caller means being the first one`。**⛔ 接线到⛔ 没有风险，是同一条里的两回事。**

### ⓑ 如果会，后果是什么形状

**形状＝少报，⛔ 误报；且⛔ 止于「少一枚检测器」，是整条 C25 腿对那次调用不存在。**

1. `:925` 成立 ⇒ `:926` `return b.assess` ⇒ 那枚 assessor 在 `:211-213` 建，只挂 canonicalizer＋classifier，**没有 `WithTaintDetector`** ⇒ `rules_gateway.go:105` `ctx.taint` 为空 ⇒ R4 对**那一次判定**永不响。
2. 同一枚空 id 在**写入侧**也一致沉默：`bridge.go:642` `if b.prov == nil || dec.TaskID == "" || !risk.IsSensitiveSource(dec.Tool) {` 不盖戳、`:848` `OpenTask` 不开 scope ⇒ 那次调用**既没来源标记也没检测器**，事后从审计线上⛔ 区分「干净的写」与「C25 整条腿被跳过」。
3. 用户侧表现：一次携带敏感来源外发的写（R4 管辖那一类）**直接执行、卡片上⛔ 出现 R4、⛔ 出现询问**。这是台账担心的那一形，成立。
4. ★ 反过来给地板（派单没问，但裁的时候需要）：`provenance.go:727-731` 对未开 scope **无条件 hit** ⇒ **⛔ 是守卫、⛔ 摘掉它**，否则空 id 会撞上 `SrcUnboundScope`＝每次写都报 R4＝误报风暴。所以这一⛔ 是「修一行加个检测器」的活；正确形状⛔ 编排者裁，我只报：**动 `:925` 必须同时决定空 id 那次调用的 scope 归属，⛔ 只动分肢。**
5. 另记一枚与 ⓐⓑ 无关但同函数的过期读数（本轮顺手量到，⛔ 我修）：`:871` 注释写 `internal/agent/loop.go:332/:321`，实际铸造与派生现在在 `:333`/`:329`/`:340` ⇒ 那两枚行号已腐烂（内容锚见问一）。

---

## 自报

- 新建：`.scratch/wisp/probes/orch/taskid-a1/00-findings.md`、`logs/assessorfor-callers.txt`、`logs/taskid-all.txt`、`logs/a838-raw.txt`（⛔ 仓库根、⛔ `.out`、⛔ 删任何东西、全部 `.md`/`.txt`，⛔ 新建 `.go`）。
- ⛔ 跑 `go build`/`go vet`/`go test`/`go env`/`go list`＝**一次都没跑**。
- 未碰：`docs/**`、票面、台账、`.gitignore`、`go.mod`；⛔ 翻框、⛔ 动 `.md`/`.txt` 之外的文件。
- ⛔ 答到的部分（⛔ 我欠／归编排者）：
  1. **变异读数**：把 `:925` 的 `|| taskID == ""` 摘掉后，`internal/tools`／`cmd/wisp` 哪⛔ 哪枚用例会红——本轮禁编译面，静态⛔ 定。
  2. **真机面**：面板/球宿主将来若自带 id 上桥，`assessorFor` 拿到的 `b.assess` 在真窗里的实际卡片表现（⛔ 真窗⛔ 看）。
  3. `b.injected` 一支本轮量到「生产码⛔ 设 `Options.Assessor`」（`grep 'Assessor:'` 非测试命中 0 行，仅 `bridge.go:132` 声明与 `:209` 赋值），故生产里 `injected==false`；**但**注入路径下检测器同样⛔ 挂，那是另一枚问题，⛔ 属本欠读数，我没往下追。
