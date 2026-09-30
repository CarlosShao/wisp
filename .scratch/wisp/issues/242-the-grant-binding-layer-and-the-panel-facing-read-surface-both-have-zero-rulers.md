# 242 — 一次性令牌的**绑定那一层全仓零尺**，而"面板那面不带令牌"这句话**只由三行注释守着**：两形我都有读数，都不是 AC#5 的洞，但下一枚腿照它们写就会签错字

- Status: **已立，未派**（09-30 15:3x，编排者立；来源＝非实现者验收腿 `197-v1` §4.1 第 3 条＋§1.3 第 2 点＋§1.4，裁决表 `docs/evidence/s1/197-ac5-selfapproval-v1.md`，52,134 字节）。
- ⛔ **为什么不塞回票 197**：那两格**不在票 197 AC#5 的射程里**（被测文件从没宣称扫"跨卡绑定"或"出向读面"）。把它们并回去＝让一格判语承担两枚它没测的东西，而这正是 `197-v1` §3.3 拿来判我那把尺的那句话（"拿一把只量得出 R1 的尺去宣布满足这条裁定"）。
- ⚠ **这不是"本机已被攻破"，照这三行读**：① **现象在哪**＝**测试仪器与注释**，不是生产行为——今天 `PanelItem` 结构上确实没有 grant 字段，`bindDigest` 也确实比对；② **有没有本机被入侵的证据**＝**没有**；③ **最坏后果是什么形状**＝**将来某一版把它改坏而 CI 不响**，以及**下一枚验收腿照现有脚注跑会误判"这格没牙"**。⇒ 这是**防回归的仪器缺口**，不是正在漏的洞。

## 现量（09-30 15:3x 编排者自己跑，别信行号、自己复算）

| 事实 | 读数 | 尺 |
|---|---|---|
| **绑定那一层今天零尺** | `bindDigest`／"不绑定"在全仓 `_test.go` 里**只出现在被验收那一枚文件自身**（`cmd/wisp/subagent_selfapproval_197_test.go`）＝**1 枚文件** | `grep -rln "bindDigest" --include=*_test.go internal cmd`（我现跑） |
| 那一层的**内容** | `grantStore.spend` 比对 `bindDigest(corr, taskID, tool, level, seq, args)`（`internal/agent/approval/approval.go:253` 铸、`:292` 花） | 〔`197-v1` D2 读数：删掉比对 ⇒ 定向尺绿＋`go test ./internal/agent/approval/ -count=1` 也绿（`ok 0.392s`）〕 |
| **"面板那面不带令牌"只有注释** | 三处逐字：`internal/agent/approval/approval.go:223`「The panel-facing surface (PanelItem) has no grant」、`internal/agent/approval/pending_read.go:9`「PanelItem deliberately carries no Params and no grant」、`cmd/wisp/approval_reply.go:28`「PanelItem has no grant field and PanelAPI has no Allow method」。**仪器＝0 枚** | 三处 `sed` 我现跑复认；"零仪器"的读数来自 `197-v1` §1.4：给 `PanelItem` 加 grant 并在 `viewLocked` 填真令牌（M1a）后**全仓 `go test ./...` 无一枚因此变红** |
| ⚠ **造不出"两张不同名"的卡**（现成的载具缺陷） | `cmd/wisp/subagent_selfapproval_197_test.go:109` 逐字 `TaskID: taskID, CorrelationID: taskID` ⇒ 同一次跑里两张卡**共用同一个 correlation id** | `197-v1` §1.3 第 2 点：13 发读数里三枚 id 逐次全等（例 `child=corr1=corr2=45ce09d5-…`）；⇒ 那条 `:466` 的"跨卡借证"实际测的是**新鲜度** |

## 要建什么（两条判据，同一块石头：令牌与"哪件事"的绑定关系今天没人守）

- [ ] **AC#1 真·跨卡：A 的活令牌花到 B 上 ⇒ 必须红**。判据＝造**两枚同时活着、correlation id 不同**的卡，把 A 的**未花**令牌递给 `Native().Allow(B的corr, A的grant)` ⇒ 要拒且**拒因指名"绑定不对"**（不是"这条路由没在跑"）。
  ⛔ **前置是本票自己的活**：先把载具的 `CorrelationID` 从 `taskID` 改成真正的 corr（`subagent_selfapproval_197_test.go:109`），否则两枚卡同名、这一发**构造不出来**（上面那行读数就是证据）。
  **正控**＝`197-v1` 的 **D2**（`approval.go:292 grantStore.spend` 不再比对 binding 摘要）：改前该判据**必须绿**、改后**必须红**；⚠ D2 今天**两把尺都绿**＝这就是本格的欠账本身。
- [ ] **AC#2 出向读面不得带 grant**：反射扫 `approval.PanelItem` 的字段名枚数与名字 ⇒ 出现 `Grant`／可答复类字段**就红**；再加一发**能力侧**判据（不是词面）：把 `PanelAPI` 的出向读面扩到"能改变卡的状态"那一形也要红。
  **正控**＝M1a 那一形（给 `PanelItem` 加 grant 并在 `viewLocked` 填真令牌）今天**全仓不红** ⇒ 本格落地后**同一形必须红**。
  ⚠ **落点两候选，这一刀归人**：**甲**＝新钉放 `internal/agent/approval` 本包测试里（**不碰任何冻结件**，编排者推荐）；**乙**＝把 `internal/panel/l2_grant_boundary_test.go` 那族从"入向 envelope"扩到"出向读面"——⛔ **那枚是三枚冻结件之一**，动它必须先有 owner 一句话并在台账落 `A##`（登记编号 `Q-74`，见票面末节）。
  ⚠ 腿原话建议的是"扩那族的射程"，**我照抄就会把一枚冻结件派给别人动**——这一处错记在我自己头上（`A471`）。

## 验收这票要谁裁（非实现者）

- [ ] **AC#3 本格不设产码要求**：判 AC#1／AC#2 两枚判据**是不是同一块石头**（值不值两格）、以及**甲形够不够**——若判"够"，本格按"甲落地"翻勾；若判"要乙"，本格不许勾，写清"等 owner 那句解冻"。

## 禁区

- `frontend/**`／`design/**` 零读零写零转述；`PLAN.md`／`docs/specs/**`／`docs/SLO.md`／`internal/observe/thresholds.go`／golden／`allowlist.txt` 一字节不动。
- ⛔ 三枚冻结件一字不动：`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`（**AC#2 若走乙形＝本票今天派不出去**，须先解冻）。
- ⛔ 不许**放宽任何断言换绿**；不许新造 D43 之外的状态词；不许把 `SessionID`／grant 塞进 `tools.Decision` 或 `agent.ToolRequest`（`internal/tools/ticket90_test.go:431` 那枚反射钉的射程，票 224 也撞过同一枚）。
- ⛔ **不许顺手把票 197 AC#5 的四条附条件当已完成**（C1／C2／C3 三枚尺的修法归 `197-r3`，见票 197 AC#5 那一格）。

## 排程与串行（共享工作树）

- 本票动 `internal/agent/approval`＋`cmd/wisp` ⇒ ⛔ **与 `197-r3`（同动 `cmd/wisp/subagent_selfapproval_197_test.go`）串行**；与 **票 224-r3**（`internal/tools/bridge.go`＋`cmd/wisp/run.go`）、**票 228-r1**（`cmd/wisp`）串行；两枚跑突变的腿**绝不并发**（互相洗读数）。
- 派单前必做**撞钉预检**（这条规矩本仓已付过两次学费）：先跑 `go test ./internal/agent/approval/ ./cmd/wisp/ -count=1` 把**今天绿的用例名**抄进派单，再逐枚读带 `grant`／`Panel`／`Allow` 字样的断言；⛔ 光 `grep` 新符号名不算预检。
- 起跑名册（我现跑，派单时重取）：`git status --porcelain -- internal cmd`＝**空**；`internal/panel` 4 枚与 `internal/ball` 1 枚红**在册**（成因＝别人在工作树里删了 `design/assets/**`），⛔ **不是**任一腿的地界，只记归因、不判不改。

## 待人拍板（`Q-74`，编排者立，**不阻塞本票甲形**）

- **要拍的那一刀**：AC#2 走**甲**（新钉在 `internal/agent/approval` 本包，⛔ 不碰冻结件）还是**乙**（把 `internal/panel/l2_grant_boundary_test.go` 的射程从入向扩到出向）。
- **推荐＝甲**。理由：甲一步就能让 M1a 那一形变红，且不动任何冻结件。
- **不答的代价**：默认走甲；**乙永远开不了**——除非他给一句"解冻那枚文件"。⛔ 按甲做完，出向读面仍有一支没被那族更大的仪器覆盖（新增 json 键那一形今天有 `TestApprovalCardViewJSONKeysMatchFrontendTypes`，而那把尺**本来就红**在册）。
- **撤销口令**：「242 改乙」（要改回乙只需说这四个字，届时先落解冻记录再派腿）。
