# panel-l2-grant-nail-r1 —— 面板侧 L2「允许」的 Go 边界钉

- 时刻 / 锚点：2026-09-25 09:39 +08，`dev` @ `86e0990f4bbb295e78adab2288180ef1b7b94500`
- 地界：`internal/panel/**` + 本件。未碰 `tools/d22scan/**`、`internal/observe/**`（另两路 agent 在跑）、
  `internal/risk/**`、`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`。
- 身份：本件是**实现者自己的记录**，不是裁决表。`SPEC-12 §4.3` #1 要求的缺口审计必须由**非实现者**做。
- 授权来历：owner 2026-09-25「丙」批准做**结构性**检查，明确不是中文字面匹配。

---

## §0 派单前提的现量（逐条在盘上重跑，不看转述）

| 前提（派单原文） | 现量结果 | 出处（本锚重导） |
|---|---|---|
| 被摘的按钮不是死装饰，有完整链路 | **成立** | `git show 53a1359^:frontend/src/components/l2-approval-card.tsx` 现读：`:161-167` 一枚 `bg-accent` Button、`onClick={() => send("grant")}` 在 `:163`；`send` 定义在 `:96-100`，体内 `:98` 调 `requestApprovalResolution(view.correlationId, outcome)` |
| `panel.ts:169-186` 发 `{method:"panel.approval.request", outcome}` | **成立** | `frontend/src/lib/panel.ts:169-186`，`bridge.postMessage` 在 `:179`，方法字面量在 `:181` |
| 违反点与合法点只差union成员（`send("grant")` vs `send("refuse")`） | **成立** | 同一文件 `:112`（关闭键）与 `:155`（拒绝键）都是 `send("refuse")`，`:163` 是 `send("grant")` |
| `ApprovalOutcome` 今天仍导出 `"grant"` | **成立，所以枚举键的检查今天会红** | `frontend/src/lib/panel.ts:50` `export type ApprovalOutcome = "grant" | "refuse";` |
| `knownComposerMethod` 是入站方法门 | **成立** | `internal/panel/bridge.go:97-103`；被答方法名在 `:35-38`（`MethodModeRequest`/`MethodWorkspaceRequest`/`MethodAttachmentAdd`/`MethodMessageSend`，const 块从 `:34` 开） |
| `ParseComposerRequest` 无生产调用者 | **成立（派单要我自验的那条）** | 见 §0.1 |
| 最强仪器在 `:417` 白名单了 `panel.approval.request`、对 `outcome` 字段色盲 | **成立** | `internal/panel/composer_test.go:409-419`（`composerRouteLiterals`，`panel.approval.request` 在 `:417`），测试体在 `:502` |
| 原生授权在 `internal/agent/approval/approval.go:221-226`，面板面携带不了 | **成立** | 该注释现读；配套：`internal/agent/approval/ui.go:155-159` 的 `PanelAPI` 只有 `Reject`/`Head`/`View` 三枚方法，**没有 Allow** |

### §0.1 「有没有真被接线的入站路径」——这条决定严重度，单独量

派单要我验完再说，不许当作死的照做。现量：**没有接线的入站路径，今天一个都没有。**

1. 全仓 Go 里 `ParseComposerRequest` 的调用者只有测试：
   `grep -rn ParseComposerRequest --include=*.go .` -> `internal/panel/bridge.go:77`（定义）、
   `internal/panel/bridge_test.go:66,83,92,102,112`（测试）、以及三处**注释**
   （`bridge.go:16`、`composer_handlers.go:37`、`cmd/wisp/run.go:225`）。生产调用者 0。
2. 唯一可能被接的处理器 `HandleModeRequest`（`internal/panel/composer_handlers.go:111`）的调用者也只有测试
   （`composer_handlers_test.go` 11 处）。`cmd/wisp/run.go:229` 声明字段 `modeWrites *panel.ModeWriteHandler`、
   `:378` 装配它，但**全仓没有任何一处读 `rt.modeWrites`**（`grep -rn modeWrites --include=*.go .` 只有这 3 行）。
3. 全仓 Go 里没有任何 WebView2 消息接收点：`grep -rn "WebMessage|PostMessage|onmessage" --include=*.go .`
   的非测试命中只有 `internal/ball/*_windows.go` 的 Win32 `PostMessageW`（托盘/球的通知，与面板 IPC 无关）。
4. `internal/panel` 生产码里入站 JSON 解码只有 1 处：`internal/panel/bridge.go:79`（`ParseComposerRequest` 体内）。
5. 代码自己的陈诉与上面一致，且是**同向**的（不是我拿注释当证据）：`composer_handlers.go:36-40`
   「NOTHING calls this handler yet, because the WebView2 event -> ParseComposerRequest hop does not exist in this
   tree (status table ①.1; tickets 33/35 are still ready-for-agent)」、`cmd/wisp/run.go:225-226` 同句。

**严重度结论**：今天这条边界是**结构性死的**（没有线），不是「活着但被判定拒绝」。
所以本钉买的是**未来**，不是当下止血：一旦 tickets 33/35 把 `event -> ParseComposerRequest -> handler` 接上，
「哪些方法被答」与「那个封套能绑哪些键」就从注释变成运行期事实，而这两枚正是本钉的断言对象。
派单里「若有真接线的入站路径，严重度要改并要报」这一支：**盘上不存在，故不改严重度，如实报死。**

---

## §1 交付物与它钉住的性质

一处新增：`internal/panel/l2_grant_boundary_test.go`（1 个文件，包内测试，`package panel`）。

钉的性质（派单原话）：**没有任何"真的被 Go 应答的入站桥方法"能接收、或携带审批结论/allow 字段。**
拆成四面，两向（正向 = 今天成立；种植 = 证明它咬）：

| 面 | 断言 | 怎么做到"不是抄一份名单" |
|---|---|---|
| 1 `TestAnsweredPanelRoutesCarryNoApprovalDecision` | 被 Go 应答的方法集里没有一枚**名字本身**是审批结论；且这集与运行期 `knownComposerMethod` 逐个互相印证 | 被答集不是硬编码：AST 读 `knownComposerMethod` 的 `switch` case 标签，标识符经本包 const 表解析成字面量。加了第 5 枚 case 会被读到 |
| 2 `TestNoInboundEnvelopeCanBindAnApprovalVerdict` | 所有被答方法共同解码进的**那一个**封套（`ComposerRequest`）及其内嵌/嵌套类型，绑不出结论键 | 反射递归 + AST 双写；"入站封套"按**结构**认定：绑 `method` 键的类型就是入站（路由字段自己就是标记），再向下闭包取内嵌/嵌套类型 |
| 3 `TestGrantWireShapesAreRefusedAtTheDoor` | 历史线上形状 `{method:"panel.approval.request", outcome:"grant"}` 被拒；且对每一枚被答方法， smuggle `outcome/allowOnce/decision/verdict/approved/grant` 后重新 marshal，值里**不再含有**这些键 | 直接跑生产函数 `ParseComposerRequest`，不是读源码 |
| 4 `TestPlantedGrantWiringGoesRedInASnapshot` | 同一台仪器指到**快照**上的种植必须变红并**点名行** | 见 §3 |
| 词表自检 `TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes` | 谓词自己不能退化成匹配不到任何东西 | 当场种植 3 个带结论字段的 Go 类型（`outcome` / `allowOnce` / 经内嵌偷渡的 `Verdict`），谓词报 0 就 `t.Fatalf`；同时要求 4 枚合法路由**不**被词表命中（否则下一枚假红会把词表放宽掉） |

键的词表按**名字形状**判（小写 + 去掉 `_ - . 空格`后做子串匹配），成员：
`approve/approval/autoapprove/grant/granted/allow/allowed/allowonce/permit/permitted/ratify/authorize/authorised/decision/decide/verdict/outcome/bypass/override`。
**不是**中文字面匹配，**不是**值匹配：面板能"发出一个结论"就是违反，送的是哪个字符串不重要。
被排除的近邻词（`confirm`/`accept`/`intent`）记录在文件头注释里，理由：它们命名的是审批的近邻概念而不携带结论，
把它们放进来会误伤今天合法的键，而误伤的后果是被放宽，不是被遵守。

一处**没有**顺手钉的东西，写明免得被当成已覆盖：`panel.mode.request` 的 `to` 字段**能**携带值
`auto_approve`（`internal/panel/composer.go:128-132` 的 `ModeRequest`、`bridge_test.go:56` 就在用这个形状），
那是一个"让 L2 不再问"的方向。它今天不构成本钉的违反，因为 (a) 没有任何接线路径能走到处理器（§0.1），
且 (b) `internal/panel/composer_handlers.go` 的 ticket-114 闸门对**变宽**方向要求确认腿（`:55-57` `ErrNoL2Confirm`，
`README` 式规则在 `:16-19`：变宽要有腿，变窄永远不必）。这一条与 `docs/reports/pending-and-issues.md` 里
"前端发 5 枚 / Go 白名单 4 枚"同源，归 C17 白名单定稿那批，不自裁。

---

## §2 闸门（真跑，原文粘贴）

`-count=2 -v` 前后对账。计数口径**写在明处**，因为派单给的 71 与我量到的 76 不是同一个口径：

- `RUN` = `grep -c -e '=== RUN'`（含子测试）
- `PASS`/`FAIL`/`SKIP` = `grep -c '^--- PASS'` 等，**只数顶层结果行**（派单给的 PASS=90 就是这个口径，能对上）
- `distinct` = `=== RUN` 的第 3 字段去掉 `#N` 后缀去重（含子测试名）

改前（锚 `004c6ec`，我进场第一发）：

```
RUN=152
PASS=90   FAIL=2   SKIP=0
panic=0
distinct=76
FAIL 的两行 = 同一枚确定性红跑两遍：
--- FAIL: TestC21DesignTokensFourWayAgree (0.00s)
--- FAIL: TestC21DesignTokensFourWayAgree (0.00s)
```

改后（锚 `86e0990` + 本件；`go test ./internal/panel/ -count=2 -v`）：

```
RUN=180
PASS=100  FAIL=2  SKIP=0
panic=0
distinct=90
--- FAIL: TestC21DesignTokensFourWayAgree (0.00s)
--- FAIL: TestC21DesignTokensFourWayAgree (0.00s)
```

- 名单**只增不减**：`comm -23 base after` = 空；新增 14 枚名字（5 枚顶层 + 9 枚子测试），逐名：
  `TestAnsweredPanelRoutesCarryNoApprovalDecision`、
  `TestNoInboundEnvelopeCanBindAnApprovalVerdict`、
  `TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes`（+ `/outcome_on_the_envelope`、
  `/allowOnce_on_the_envelope`、`/a_verdict_smuggled_inside_an_embedded_struct`）、
  `TestGrantWireShapesAreRefusedAtTheDoor`（+ `/the_historical_wire_shape_is_refused_because_Go_does_not_answer_its_route`、
  `/every_answered_route_drops_a_smuggled_verdict`）、
  `TestPlantedGrantWiringGoesRedInASnapshot`（+ `/A_…`、`/B_…`、`/C_a_wired_grant_door_goes_red_on_both_halves`、
  `/D_the_historical_JSX_button_is_invisible_to_this_instrument`）。
- 顶层 PASS 90 -> 100：恰为本件 5 枚新顶层测试跑两遍（5x2=10），无其他名字进出。
- `^panic:` **0 行**（改前也是 0），所以本次不存在"panic 吞掉同包兄弟读数"的账。
- **仍红的 2 行不是我造成的**，且我未修未跳：`TestC21DesignTokensFourWayAgree`，
  台账 `A208③` P1 —— owner 把 `design/assets/tokens.css` 移出工作树（`git status` 现量：
  `design/assets/*.css`、`design/screens/*.html` 等 16 枚 ` D`，另有未追踪 `design/doubao/`、`design/old/`），
  刻意留红。改前改后同枚同名同数，计 2 行 = 1 枚确定性红 x 2 遍。
- `gofmt -l internal/panel/` -> **空**（本文件在内）。
- `go vet ./internal/panel/` -> **空**（rc 0；过程中它先抓出我一处 `reflect.StructTag.Get` 误用为双返回值，已改 `Lookup`）。
