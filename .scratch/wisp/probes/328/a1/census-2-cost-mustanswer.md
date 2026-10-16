# 328-a1 只读普查 件② — ③ 改形代价列 + 必答题（票 328 AC#0 料）

锚 `44f51fc8`；尺现跑于 dev 工作树（`internal/agent`＋`internal/panel` 对锚点 diff 空）。
"⛔ 选修法"——甲/乙/丙只列代价，由编排者裁。行号只当读数。

## ③ 谁在读 corr 串 / `#` 语义（改 corr 一侧的爆炸半径）

**corr 串的 `#` 语义生产者**（尺 `grep -rnE 'Sprintf\("[^"]*#|taskID \+"#"|corr \+ "#' internal/ cmd|grep -v _test|grep -v 'AC#\|ticket\|D[0-9]'` rc=0）：
- `loop.go:605` `taskID#call-<idx>`（支B）· `loop.go:607` `taskID#callID`（支A）。
- `queue.go:168` `fmt.Sprintf("%s#%d", d.CorrelationID, q.seq)` — 队列撞名改写（全数字尾）。
- `queue.go:614` `corr + "#replay"` · `:616` replay 撞名 `%s#%d`。
- `queue.go:165` 空 corr 回退 `approval-<seq>`（**无 `#`**）。

**corr 串的 `#` 语义读取者（生产，唯一一枚）**：`subagent_roster_197.go:277 splitClashSuffix`，
唯一调用点 `:267`（尺 `grep -rn splitClashSuffix internal --include=*.go` rc=0：def＋doc＋:267 三命中）。
⇒ 改 corr 一侧＝动 `loop.go:605/607` 两支，会牵动所有钉这两支的测试（下）；改读取一侧＝动 `:277` 一把尺（乙）。

**无关 `#` 读取（排除，非 corr）**：`internal/panel/assets.go:146`、`cmd/wisp/panel_control_windows.go:539`
都用 `strings.IndexAny(…,"?#")` 剥 URL fragment，与 corr 无关。

**钉 corr 形状的测试（逐枚写今天钉哪一形）**（尺 `grep -rniE 'HasPrefix|TrimPrefix' …corr`＋读体，rc=0）：
- `loop.go:603` 两支作差尺 `Test285CallCorrIsDistinctPerCall`（`ticket285_corr_distinct_rulers_test.go:44`）：
  钉 ① 每支 corr 非空、② `corr != task`（`:73` collapses-onto-task 即红）、③ `HasPrefix(corr,task)`（`:75`）、
  ④ 支A `TrimPrefix(corr, task+"#") == req.CallID`（`:136`）。⇒ **丙**（corr=裸 T）撞 `:73`；**乙/甲**若改尾映射撞 `:136`。
- `loop_golden_test.go:79` `TrimPrefix(corr, task+"#") != "call_e1"` — 钉 golden 支A 尾＝`call_e1`（带前缀，非数字）。⇒ **乙**改尾形撞这枚。
- `corr_percall_242_test.go:139` · `tools/loop_approval_test.go:220` · `tools/subagent_depth_fixture_287_test.go:140` ·
  `tools/ticket283_corr_identity_rulers_test.go:222` · `tools/ticket285_corr_rows_rulers_test.go:75`：
  一致钉 "corr 以 task 为前缀、且 `corr != taskID`"（C18 每调用独立路由键）。⇒ **丙** 逐枚打红。
- `subagent_blocked_220_test.go TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt`（`:124`）—**名册那把尺的正/负判据**：
  用例 `:130 child-1#9→true`、`:131 child-1#147→true`、`:132 other-1#9→false`、**`:133 child-1#retry→false`**、`:134 child-1x→false`、`:135 ""→false`。
  ⇒ **乙**（把那把尺放宽到认 `call-`/非数字尾）直接反转 `:133`（"非数字尾不是队列形"）→ 打红；且 `:132/:134` 正是"⛔ 宽匹配"的护栏，放宽会一并失效。这就是"同名两行 clash 判据"。
- `TestAL1WindowInWaitingLightsItsRow`（`subagent_blocked_220_test.go` ~`:147`）注释称 "Since bd124b2a … 每调用 corr、⛔ ==taskID"——
  已 `git rev-parse --verify bd124b2a`＝真存在（`bd124b2ae0a5…`，rc=0）；行为（loop.go:603/676）已由读码证实，非拿注释当凭据。

**AC#0 附带：今天硬写撞名形的用例有几枚**（尺 `grep -rnoE '"[^"]*#[0-9]+"' internal/panel/*_test.go cmd/wisp/*_test.go internal/agent/*_test.go` rc=0）：
真 corr 字面全在 `subagent_blocked_220_test.go`：`task-7#9`(:15 注释)、`child-1#9`(:108/:112/:130/:186)、`child-1#147`(:131)、`other-1#9`(:132)；
其余命中（`ban #6`/`AC#1`/`AC#3`/`AC#10` 等）是散文假阳，已剔。⇒ 名册用例用 `card197(corr)`（`:45`）手造 corr、**不经 `callCorr`**，正是票说的"缺世界"。

## 甲/乙/丙 各自的额外代价（不裁形，只列）

- **甲**（卡片也带任务号）：`ApprovalCardView`（`panel/approval.go:39`）**根本没有 TaskID 字段**（尺 `grep -c TaskID internal/panel/approval.go`＝0，正控 `grep CorrelationID`＝rc0/3 命中 ⇒ 0 是真 0）。
  要动两处＋契约面：① 给 `ApprovalCardView` 加 `taskId` JSON 键；② 让 `subagent_roster_197.go:213-214` 也索引卡片的任务号（现只索引 `card.CorrelationID`）。
  撞的既有钉（逐枚）：`internal/panel/approval_test.go:105 TestApprovalCardViewJSONKeysMatchFrontendTypes`——它把 `ApprovalCardView/Snapshot/ResultChunk` 的 Go JSON 键与 `frontend/src/lib/panel.ts`（接口在 `:24`，含 `correlationId`）**双向**对账（`subtract` 两向都报，`:128/:131`）。⇒ Go 侧加 `taskId` 必须同时在 **`frontend/**`** 的 TS 接口补一枚，否则红；而 `frontend/**` 归机主外部 agent（本编队⛔写面）⇒ **跨面协作**。
  C17 面：`docs/specs/SPEC-08:156 §5.2 C17 PanelBridge 方法白名单`、`:235 correlationId 路由测试`——卡片身份键属 C17 契约，新键走契约批准。
  ⚠ 纠票 328 AC#1 的一处关联：票里"撞票 313 那枚键钉"——`config_313_test.go:283 Test313PayloadKeySetIsFrozen…` 钉的是 **settings payload 顶层键**（`["chatModel"…]`＋provider 行键），⛔ 覆盖审批卡/Snapshot 卡内键 ⇒ 加卡键**不**撞 313 这枚；真正撞的是 approval_test.go:105＋panel.ts。另 `cmd/wisp/panel_assets_143_test.go:52 cardView143` 用容忍未知字段解码，加键不破它。
- **乙**（名册那把尺认非数字尾）：只动 `subagent_roster_197.go:277-288` 一枚尺；但反转 `subagent_blocked_220_test.go:133`（`#retry→false`）并把 `:132/:134` 防过宽护栏拆掉；
  尺 doc `:272-276` 自陈"放宽＝比本仓任何代码产生的都宽"，属**放宽判据**（票 AC#1 丙格语，实为乙的语义风险）。
- **丙**（corr 直接＝任务号）：逐枚打红"corr≠taskID/前缀+尾"那族（ticket285 `:73/:136`、ticket283 `:222`、287 `:140`、loop_approval `:220`、242 `:139`、loop_golden `:79`）——
  且把 C18 每调用独立路由键退回每任务一键（一个 task 并发两调用塌到同键，正是 242/285 立起来的洞）；join 侧则靠 roster `:266` 逐字入库即可点亮（host mode-switch 卡即此形的先例，见下）。

## 必答题（⛔ 盖章，给尺）

**"生产里今天有没有任何一条路径能拼出末 `#` 后全是数字的 corr？"**

分两读，各自有据：

**(a) 字面"存的卡片 corr 末 `#` 后全数字"——有，但只由队列自造、且点不亮行。**
生产者＝队列撞名/重放改写：`internal/agent/approval/queue.go:168` `corr = fmt.Sprintf("%s#%d", d.CorrelationID, q.seq)`（＋ `:616` 同形）。
真样例串（算的）：某 task 同一 corr `<uuid>#call-0` 再次入队撞名 ⇒ 存成 `<uuid>#call-0#7`，末 `#` 后＝`7`（全数字）。
但 `splitClashSuffix("<uuid>#call-0#7")` 剥出 base＝`<uuid>#call-0` ≠ `<uuid>`（roster `:267` ok=true 但 base 非 T）⇒ **仍不 join**。（此改写只在同 corr 二次入队时触发。）

**(b) "能被那把尺剥成裸任务号 T（即恰一个 `#` 的 `T#<全数字>`，真能点亮行）"——没有。**
逐枚排除（件①名册）：
- `loop.go:605` 支B：尾＝`call-<idx>`，含 `c` 与 `-` ⇒ `:283` 判非数字 ⇒ 不 join。
- `loop.go:607` 支A：尾＝provider callID 逐字；三家在用 provider 都给带前缀号（openaichat `toolasm.go:59` `f.ID`/回退 `:114 call_unnamed_%d`、responses `stream.go:447` `call_id`/`id`、anthropic `stream.go:258` `toolu_`），⛔ 全数字 ⇒ 不 join。
- 队列 `:168/:616` 数字尾：永远叠在已带 `#` 的 `callCorr` 值上 ⇒ base≠T（见 (a)）。
- `queue.go:165` 空 corr 回退＝`approval-<seq>`（无 `#`）；host mode-switch 卡 `run.go:878` 不带 CorrelationID ⇒ 走此回退，非 `T#<数字>` 形，且 `host:mode-switch` 本就不是 roster 任务行。
- 陪聊全双工/realtime：到不了 `callCorr`（件①）。
⇒ **今天没有任何生产路径拼出一枚能点亮干活任务行的 corr**；票 328 的"交集为空"成立。唯一残余（⛔ 代码不担保、经验性）：某 provider 若真回一枚纯数字且不含 `#` 的 call id，支A 会拼出 `T#<数字>` 而 join——但仓内无此造号代码、在用三家 provider 不这么发。
