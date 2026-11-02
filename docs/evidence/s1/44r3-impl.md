# 票 44 落地腿 44-r3 — 实现记录（最小两件：AC#3 事件出口 ＋ 生产咬合）

起手：HEAD `7e3ade426`（编排者 A1432 派单笔）；前腿 checkpoint `19b48978b`＝**只加 `internal/observe/quota_44.go` 一枚（+337 行）**，`git show --stat` 逐字为 1 file changed。起手四包 `go build ./internal/observe ./internal/agent ./internal/memory ./cmd/wisp` **exit=0**；`git status --porcelain` 跟踪面只有 `M .scratch/wisp/issues/22-web-tools-d30.md`（＝22-v1 的票面，⛔ 本腿写面），本腿写面干净。

纪律：只 commit、⛔ push；显式 pathspec；⛔ `add -A`／amend／reset；未碰 `thresholds.go`／golden／SLO／`internal/plugin`／`internal/tools`／`internal/memory`（本波最终没需要动 memory）。

---

## ① AC#3 — 80% 预算预警的**事件出口**（票面 AC#3 `80% warn event`；44-v1 §3.2 ②／§7 ② 具名的那一寸）

**裁决表点的原话**（`docs/evidence/s1/44-v1-verdict.md:80`）：成本路径只有 `loop.go:415` 那一枚 `EvError`，且只在 100% 被拒时发；warn 态 `AllowNewTask=true` ⇒ `preDispatchBudget` 直接 `return nil` ⇒ **总线上什么都没有**。§7 ② 给的修法逐字＝"要么补一个 warn 事件的发布点（仍在 `internal/agent` 写面内、⛔ 不碰装配）"。

**落地（三处，全在 `internal/agent`）**：

| 件 | 改动 | 形状与理由 |
|---|---|---|
| `sink.go` | 新增 `EvWarn EventKind = "warn"` | 与现 `EvError` 并列的**独立枚**，不复用 error 态：error 说"这单失败了"，warn 说"这单照常开工、但钱用到哪条线了"。裁决表 §3.2 ② 要的就是 warn **event**；渲染归票 40，本腿只管生产者 |
| `costbook_44.go` | `preDispatchBudget(ctx)` → `preDispatchBudget(ctx, taskID)`；在 `if !rep.Metered \|\| rep.Verdict.AllowNewTask` 那道 return **之前**插一个发布点：`rep.Metered && AllowNewTask && Level != LevelOK` ⇒ `l.publish(Event{Kind: EvWarn, TaskID: taskID, Text: rep.Verdict.Notice})` | 事件的 `Text` **逐字**取 `BudgetVerdict.Notice`（不重造文案，`TestAC3WarnBudget…` 钉 `ev.Text == pre.Verdict.Notice`）。条件恰好覆盖裁决表点名的两种静默形：**warn 线（80–99%）** 与 **`over_budget=warn` 越 100%**；100%＋pause 那一支仍只走 `EvError`（`AllowNewTask=false` ⇒ 不双发） |
| `loop.go:412` | 调用点补 `taskID` | 事件必须点名被预警的那单，`Event.TaskID` 是面板/日志的关联键 |

**新测试**：`internal/agent/costwarn_44r3_test.go` 四枚（真环跑，复用 `newStore44`/`newHarness("tool-then-text")`/`price44` 既有接缝，⛔ 新 mock）：

- `TestAC3WarnBudgetPublishesWarnEvent`＝日花费对 1.25x 上限 ⇒ 恰 8000bp＝warn；跑完 `StatusCompleted`＋`requests()>0`（warn 不拦），且 `LastOf(EvWarn)` 非空、`TaskID` 对上、文案含"预算"与"80%"、与 `Notice` 逐字相等；同时钉 `LastOf(EvError)==nil`（两种 kind 不许混）。
- `TestAC3WarnEventFiresPastTheLimitInWarnMode`＝`over_budget=warn` 越 100% 那一支也必须有事件（裁决表 §3.2 ② 的第二半）。
- `TestAC3UnderTheAlertLinePublishesNothing`＝1% 用量⇒**零事件**（防"每单都预警"的假实现）。
- `TestAC3RefusalStaysAnError`＝100%＋pause 仍只 `EvError`、不得双发。

**读数**（`go test ./internal/agent -count=1 -v`，同包 AC#3/#4/#6 家族一起跑）：**全绿** — 本腿新四枚 PASS，既有 `TestAC3NewTaskRefusedBeforeTheProvider`／`TestAC3RunningTaskFinishesPastTheLimit`／`TestAC3WarnModeStillRunsTasks`／`TestAC3BudgetGateReadsTheBookedSpend`／`TestAC4*`／`TestAC6*` 无一被 signature 改动打破；**整包 `ok github.com/CarlosShao/wisp/internal/agent 4.362s`**；`go test ./internal/observe -count=1` **ok 4.294s**（`quota_44.go` checkpoint 完好）。

**反向判据 X1（真红）**：把发布那一行换成 `_ = rep.Verdict.Notice`（in-tree 突变，跑完还原）⇒
`costwarn_44r3_test.go:64: no EvWarn on the bus at the 80% line…`＋`:104: crossing 100% in warn mode published no event at all`＝**两枚 FAIL**，而 `TestAC3RefusalStaysAnError` 仍 PASS（证明红的是事件出口、不是别处）。**还原核**＝`git hash-object` 前后 `same=True`（`39134bba…`）。

---

## ② 生产咬合 — `cmd/wisp/run.go` 把价卡与 `[cost]` 真交给环路

**缺口原话**（§3.2 ③／§7 ①）：`BudgetFrom`、`Loop.SetCostBudget` **零调用者**，`Config.Price` 无人填 ⇒ 真跑时每一轮 0 micro、表恒 uncapped，`[cost]` 与 `llm.providers.<p>.models.<m>.price` 是"能加载、没人读"的死键。

**落地**：

- `costbook_44.go` 新增 **`agent.CostWiring(cfg *config.Config, provider, model string) (config.Price, observe.BudgetParams)`**——放 `internal/agent` 而非 `cmd/wisp`：那层本来就同时握着 `config` 与 `observe`（`observe` ⛔ 反向 import config），装配根只做三行调用，测试面落在库层。价卡是**查出来的、不是猜的**：`cfg.LLM.Providers[p].Models[m].Price`，查不到就返回零值 `config.Price` ⇒ `PriceKnown()==false` ⇒ AC#4 的 `价格未知`／只显 tokens；`[cost]` 走 `observe.BudgetFrom`，坏叶子（`alert_threshold=1.5`、`over_budget=stop`）被它拒 ⇒ 返回**零 BudgetParams＝uncapped**，与"用户没设预算"同一诚实读法，绝不 clamp 成另一套政策。
- `run.go:1329`（`execute()` 里，**紧跟 per-task 活配置重读 `cfg = rt.mgr.Config()`** 之后、`agent.Options` 之前）：`card, budget := agent.CostWiring(cfg, ep.Provider, ep.Model)` ⇒ 票 223 的 reload tick 下一单就够到两枚旋钮。
- `run.go:1357`：`agent.Config` 里填 `Price: card,`（`loop.go:493` 每轮 `cost.AddUsage(l.opt.Config.Price, turn.Usage)` 读的正是这一枚）。**`Currency` 故意不填**：`[cost]` 与 `config.Price` 都没有币种键，编一个币种＝造数；沿用 `loop.go:233` 既有默认 `CNY`，这笔账登记在下面"未闭合"。
- `run.go:1390`（`agent.New` 的错误出口**之后**）：`loop.SetCostBudget(budget)`。
- `run.go:1798`：`consoleSink.Publish` 补 `case agent.EvWarn:` ⇒ 打印 `\n[预警] %s\n` 并 `changed=true`（走同一条面板快照）。事件有生产者也得有读者，否则 CLI 真跑仍看不见。

**新测试**：`cmd/wisp/costwiring_44r3_test.go` 三枚 —
`Test44R3CostWiringCarriesPriceAndCaps`（价卡逐字段搬运＋`observe.PriceMicros(agent.CardOf(card), {1000,500}) == 7000` 真算出钱）PASS；
`Test44R3CostWiringSaysUnknownPriceAndBadBudget`（无卡模型⇒`Known()==false`且 1M/1M 计 **0 micro**＝"不填价⇒成本 0 可见"那一发；坏 `[cost]`⇒uncapped）PASS；
`Test44R3RunGoAppliesTheWiring`（装配根在场钉：`run.go` 必须含 `card, budget := agent.CostWiring(...)`／`Price: card,`／`loop.SetCostBudget(budget)`／`case agent.EvWarn:` 四枚，并钉**次序**＝`mgr.Config()` < `CostWiring` < `agent.New` < `SetCostBudget`）PASS。
基线读数：`go test ./cmd/wisp -run Test44R3 -count=1 -v` **3/3 PASS、ok 0.083s**（时刻 22:23:25；须带 `PATH=third_party/sherpa-onnx`，票 98 那笔账，不注入＝加载期 `0xc0000135` 假失败）。

**反向判据 X2 — 本腿没能拿到干净的红灯，具名登记**：把 `Price: card,` 摘掉、`CostWiring` 返回改 `_, budget`（可编译形）后复跑，`cmd/wisp` 整包 `[build failed]`——挡路的**不是本腿**，是共享盘上刚落的一枚**未跟踪外来件 `internal/tools/document_text.go`**（24-r2a 在写；三处编译错：`"os" imported and not used`／`undefined: sort`／`undefined: docxText`），`cmd/wisp → internal/tools` 依赖它，所以此刻**任何** cmd/wisp 测试都跑不动。牙齿改用逻辑尺投影：突变体与 HEAD 的差异恰为缺 `Price: card,` 一枚（`Contains` 实测 True→False），即 pin 那条断言必假。**还原核＝`same=True`（`904683e0…`）**；外来件 ⛔ 未碰、未移、未删。⇒ **r4／验收腿须在 `internal/tools` 收口后重跑 X2 拿真红**。

---

## 门禁与形态

- `gofumpt -l`（绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`）对本腿六枚文件：**输出空**。（同命令扫 `internal/agent` 全目录另报 `internal/agent/approval/pending_read.go`、`internal/agent/tools.go` 两枚**非本腿文件**，未碰。）
- `go vet ./internal/agent` **exit=0**。`go build ./internal/{observe,agent,memory} ./cmd/wisp` **exit=0**（外部脏件出现之前）。
- `d22scan`（`cd tools/d22scan && go run . -root ../../`）：**clean — no D22 ban violations**；覆盖 `internal/` 779 Go 文件、`cmd/` 187 Go 文件（注释与 `_test.go` 同扫）。
- 界面文案零 emoji：新增串只有 `[预警]`；注释里用了 `→`（仪器刻意不扫的空隙带），未用 `✓`／`≤`。

## AC#5 剩余（登记给 r4，⛔ 本腿没碰）

checkpoint `19b48978b` 交付的是 **observe 配额层**（`internal/observe/quota_44.go`：per-model 日/月闸、plan-credit 分账估算 marker、`SelectChain` 防负信号、会话级 `QuotaLedger`），编译与本包测试均绿——那是"部分"。裁决表 §7 ① 那一寸**原样未动**：盘上 `task_log` 仍 11 列无 provider/model、`cost_daily` 仍 5 列 ⇒ per-model 花费无处还原 ⇒ **要么机主先裁 D35 加列/加表（人工批准面），要么编排者把 AC#5 正式拆出本票**；本腿⛔ 未动 `internal/memory/schema.go` 一字。给 r4 的三件候选（按裁决表点名）：① 配额层的**落盘面**（等 D35）；② `QuotaLedger`/`SelectChain` 接到 `agent` 派发链＋链推进；③ **本腿咬合的下游**：`Currency` 无源（`[cost]`/`config.Price` 都没有币种键）——要么 D36 加 `[cost] currency`（改契约文字＝机主批），要么由价表 `observe.PriceTable.Card` 的 currency 承载并把 `agent.Config.Currency` 接上；本腿选择**不猜**。

---

*本腿＝44-r3 实现腿，只做上面两件；只 commit、⛔ push。*
