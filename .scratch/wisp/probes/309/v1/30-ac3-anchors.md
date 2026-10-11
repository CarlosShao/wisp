# 309-v1 · `AC#3` 文案一致性（只改行号与枚数，⛔ 改契约声称；新内容锚今日在盘上对不对）

## 1. "只碰注释行"——三把尺，逐把带读数（Rig＝`logs/ac3-rig.txt`，原始＝`logs/ac3-comment-only.log`）
| 尺 | 射程 | 读数 | 判定 |
|---|---|---|---|
| `git show --numstat 16e86bae` | 三枚文件 | `73 0` 新测试／`32 3` `bridge.go`／`134 0` 新测试 | 只有 `bridge.go` 是既有文件 |
| `git diff -w --numstat 16e86bae^ 16e86bae -- bridge.go` | 单文件 | **`32 3`** | ★**我派单里那句"expect 0"是错前提**：`-w` 只忽略变更行**内部**空白，⛔ 能把净增的 29 行注释变没 ⇒ 它⛔ 是"只改注释"的尺。已按实测量报，⛔ 挑一把会让它=0 的尺 |
| 逐行性质尺：`git show -U0` 取全部 `^[+-]`（剔 `+++`/`---`）再滤 `^+ //`/`^- //` | 单文件 | 变更行 **35**，其中**非 `//` 注释行＝0** | ⇒ 改动**全在注释行**（真尺） |
| 代码身份尺：两版 blob 各 `grep -v '^[[:space:]]*//'` ＋ 剥空行后 `cmp` | 单文件 | 代码行 **768 ⇄ 768**、`cmp` **IDENTICAL**、`rc=0` | ⇒ **零代码行变动**（最硬那把；顺带坐实"行为逐字节不变"） |
| 空白剥净 `cmp`（派单点名的那把） | 单文件 | **DIFFER**，`byte 29474` | 预期内：净增注释⛔ 会被剥空白抹掉 ⇒ 报为**错前提**，⛔ 当"改动越界"的证据 |

契约声称⛔ 被动（逐条对：`git show 16e86bae^:bridge.go` ⇄ `:bridge.go`）：
守卫布尔式 before＝`925: if b.injected || b.prov == nil || taskID == "" {`、after＝`951 起同一串逐字节`，**只有行号移动**；attach 那支 `WithTaintDetector(b.prov.Detector(taskID))` 两版计数都＝1。
`:926-929` 那句被 `AC#3` 点名的"per scope, hence per task"（票 242 的前因）**一字未改**（diff 里那一支⛔ 出现在 ± 行中，只在它后面追加）。

## 2. 新写上去的内容锚，今日逐枚在盘上核（尺＝`logs/ac3-anchors.log`／`ac3-anchors-risk.log`）
| 注释里的锚（逐字） | 今日真身 | 判定 |
|---|---|---|
| `in internal/agent/loop.go the "defer revokeAdmission()" line that runs right after the AdmitTask hook` | `loop.go:373-379`：`revokeAdmission := func() {}` → `if l.opt.AdmitTask != nil { … }` → **`379: defer revokeAdmission()`** | ✅ 落点就是 AdmitTask 钩子之后那一行；腿报的"`:366` 漂⇒真身 `:379`"**复跑对得上** |
| `RunAsync's "id := newTaskID()"` | `loop.go:329`，所在函数 `func (l *Loop) RunAsync` 起 `:328` | ✅ |
| `Run's "return l.run(ctx, newTaskID(), input)"` | `loop.go:340`，`func (l *Loop) Run` 起 `:339` | ✅ |
| `the one newTaskID() construction point … (id 36 chars, never empty)` | `loop.go:1153-1166`：两条返回支都走同一 `hex.EncodeToString` 拼接 ⇒ **恒 36 字符**（熵失败退化＝时钟派生，仍 36） | ✅ |
| `ruleTaint short-circuits on its if ctx.taint == nil guard (internal/risk/rules_gateway.go)` | `rules_gateway.go:101 func ruleTaint` ／ **`:102 if ctx.taint == nil {`** | ✅ 函数名与那一行都对得上 |
| `scopeMarks in internal/risk/provenance.go fail-closes an unregistered scope to R4 whenever a DIFFERENT scope still holds marks` | `provenance.go:694 scopeMarks`：`reg != nil` ⇒ `(nil,false)`；否则遍历其它 scope，`len(r.marks) > 0` 才 `return nil, true` ⇒ `:727`/`:828` `if unbound {` 交 `SrcUnboundScope` | ✅ **这句正是票面追加更正②的口径**（⛔ 再是"无条件"） |
| `Test309ToolRequestProductionConstructionNail`／`Test309AssessorForEmptyTaskIDDropsTaintDetector` | 两枚用例都在盘上且本腿各跑过（`logs/pristine-309.log`） | ✅ |

## 3. ★`AC#3` 自己交的那枚"带尺重述"第一发重跑就对不上（具名缺陷，⛔ 阻格）
`bridge.go:879-881` 逐字：`Ruler: grep -rn 'ToolRequest{' internal cmd --include=*.go minus _test.go yields the single internal/agent/loop.go point`。
本腿现跑（`logs/spotcheck-rulers.log`）＝**3 行**：
```
internal/agent/loop.go:675:  req := ToolRequest{
internal/tools/bridge.go:879:// dispatches on the bridge except loop.go. Ruler: grep -rn 'ToolRequest{' internal
internal/tools/bridge.go:940://     empty); Test309ToolRequestProductionConstructionNail pins `agent.ToolRequest{`
```
⇒ 那把尺**把注释自己写的那两行也数了进去**（自指尺）。真身结论（构造点 1 枚）**没变**：再剥一行首为 `//` 的命中＝**1** ✅。
判＝**口径缺陷**：本格判据是"那句话要带尺重述"，尺带上了、但⛔ 重跑得出它写的数；票面「现量」自己的规矩"⛔ 引用前先重跑"在这里第一次重跑就破。处置⛔ 在我（产码注释＝实现面）＝欠账 `Q-309-v1-selfref-ruler`。
同族第二枚（⛔ 产码、在腿的测试件头注释里，`toolrequest_nail_309_test.go:9-11`）："全仓 **60** 枚＝`_test.go` 内 **59**／外 **1**" ⇒ 今日同一把尺＝**69／66／3**（新增的正是这两枚 309 件自己）；"钉的现量"那两枚读数（0、245）⛔ 受影响。记 `Q-309-v1-stale-count-comment`。

## 4. ⛔ 属 `AC#3` 射程但被腿自己点名没扫的残余（我替它量了）
`AC#3` 只点名 `bridge.go:869-876` 那一段 ⇒ 其余绝对号⛔ 本格。腿在 `90-final.md` ⑤ 具名"⛔ 顺手全扫"。本腿现读（`logs/ac3-residual-anchors.log`）＝`bridge.go` 里还剩 **5 处**绝对 `loop.go:N`：
`:113` 指 `loop.go:563`（今日那行＝`\t})`）、`:114` 指 `696-700`（今日＝终端 tool_call 注释块，勉强贴边）、`:230` 指 `783-806`＝`decideRisk`（**真身在 `:813`**，漂约 27 行）、`:252` 指 `692-694`（今日＝panic recovered 那两行，形对但号已移）、`:1084` 指 `loop.go:585`（今日＝`rowID int64`，⛔ 对不上）。
⇒ **4 枚今日⛔ 指到它说的东西**。⛔ `AC#3` 的账（票面⛔ 让它扫），⛔ 309-r1 的越界（它具名停手了）；是**票 154/151 那一族的既有债**，记 `Q-309-v1-anchor-debt-154` 交编排者归口。
