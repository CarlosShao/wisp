# 44r4-impl — 票 44 AC#5 接线腿（44-r4）AC 逐支读数

- **腿**：44-r4＝实现腿（写腿），范围＝只做 AC#5 接线（派单 `44r4-dispatch.txt`）。
- **base**：`f2cc9f73f`；**交付 commit**：`a11e63bee`（2 枚代码／测试，493 行）＋本件追加 commit。只 commit、⛔ push。
- **写面核**：起手 `git status --porcelain -- internal/agent internal/memory internal/observe` ＝ **空**（三包无他脏）；交付只落 `internal/agent/**` 两枚新文件；⛔ 碰 cmd/wisp（44-r3 面）／internal/tools（24-r2a 面）／internal/plugin／thresholds／golden／票面／台账。
- **临时件**：无残留（全部走盘上测试，未建 worktree／checkout）。

## 0. SPEC shape 核（⛔发明闸门）

派发前读的三处，全部**现盘对字**后才动手：

| 字段 | SPEC/config 出处 | 读数 |
|---|---|---|
| per-model 日/月配额 | `docs/specs/SPEC-03 §3.1`（配额三层）＋ `internal/config/schema.go:372-373` `quota_daily_micro` / `quota_monthly_micro` | 已定义，直接消费 |
| per-provider plan credit | SPEC-03 :32 行 §3.1 ＋ `schema.go:396-397` `billing` / `plan_credit_total_micro`（`providers.<name>` 级） | 已定义，直接消费 |
| 防负＝派发前检查＋100% 走兜底链＋降级可见 | SPEC-03 §3.1 逐字「防负＝派发前检查（C23 钩子），100% 硬停并自动走兜底链，降级事件悬浮球+面板可见，不得静默」；票面 :33-40 | 接线形状照此，零发明 |
| 盘上落点 | `docs/specs/SPEC-02-data-storage.md:132/148` `provider_health.quota_state`（'ok'|'throttled'|'exhausted'|'unknown'）＋ `internal/memory/dao_providerhealth.go`（`RecordProviderError`、`QuotaExhausted`） | D35 冻结表**内**既有列既有 DAO，⛔ 加列（44-v1 §7.1 的批准面本腿没碰） |

## 1. 交付形状（AC#5 三样欠账逐一对位）

新文件 `internal/agent/quota_44r4.go`：

- **派发前 gate＝逐链元素包装** `(*QuotaGate).GateProvider(p, "provider/model")`：每次 Stream 调**之前** `Decide()`；`remaining 0`＝inner **零调用**（一字节不出机器）＋返回 `observe.ClassBudget` 分类错误——这正是票 11 `llm.ChainRunner` 步进机制吃的形状（其自带测试 `internal/llm/llm_test.go` 把 ClassBudget 元素喂进链、验证跳到下一家）⇒ **跳过派发＋链继续推进＝不崩链**，步进逻辑⛔ 本腿没动（票 11 所有权，票面 :39-40）。
- **可见事件**：warn（80%）发 `EvWarn` 后仍派发；skip（100%）发 `EvWarn`（文案＝observe `quotaStopNotice`，含「跳过该模型，按兜底链继续」与恢复时钟）——不静默。
- **plan-credit 分账记账**：`Book(key, micros)`→`observe.QuotaLedger.Post`＝plan 消费进 `provider_plan` 桶、pay-per-token 进 model 日/月桶，**两本账物理分离**（AC#5「tracked separately per model」）；余额＝手录总额−记账、`Estimated=true` 带「估算」marker（票面 :30-32）；**负数 micros 拒**（反负数防护的记账半边）。
- **盘上信号**：exhausted 走既有 memory DAO `RecordProviderError(…, class="budget", quota_state="exhausted")`（`QuotaSignalWriter` 断言口，编译期钉 `(*memory.Store)` 满足接口）；写失败只 log 不拦派发。
- **config 适配** `QuotaWiring(cfg)`：只搬 SPEC 已定义字段；模型 billing 缺省继承 provider、再缺省 pay-per-token（schema 默认值）；alert 线复用 `[cost] alert_threshold` 单一口径（⛔ 第二把 80% 尺）；坏 `[cost]`＝零参数＝不拦。
- **链查询半** `Chain(chain)`：复用 44-r2 `SelectChain`，返回首个可派发元素＋按链序的 skips；全链耗尽＝`chosen==""`＋全 skips，收尾归票 11 runner。

## 2. 定向测试读数（`go test ./internal/agent -run AC5 -count=1 -v`）

**7/7 PASS，FAIL=0**（`ok agent 0.031s`）：

| 用例 | 钉的 AC#5 支 | 关键断言 |
|---|---|---|
| `TestAC5RemainingZeroSkipsDispatchChainAdvances` | 「remaining 0 → dispatch skipped, chain advanced, visible event」 | 闸内层 `calls==0`、错误 class＝budget、下一元素 `calls==1`（链推进）、`Chain()`＝chosen pb/m2＋skips[pa/m1]、事件 EvWarn 含「跳过」 |
| `TestAC5WholeChainExhaustedStillAdvancesNotBreaks` | 不崩链的边界 | 全耗尽＝`chosen==""`＋2 skips，gate 不 terminate 任务 |
| `TestAC5PlanCreditBalanceEstimatedAndSeparateFromPay` | plan-credit 分账＋plan/pay 隔离 | 余额 1000−900=**100**、Binding=plan-credit、Estimated、90%＝warn；pay-per-token 兄弟模型 Decide＝OK 零绑定；plan 用尽后**只闸 plan 模型**（pay 兄弟 calls==1 正常派发）；stop 文案含「估算」marker |
| `TestAC5NegativeSpendRefusedLedgerUnchanged` | 反负数防护 | `Book(-400)` 报错、余额仍 600（负记账不凭空造 credit） |
| `TestAC5QuotaWarnPublishesVisibleEventAndStillDispatches` | 80% warn 有事件出口 | 恰 1 枚 EvWarn 含「80%」、仍派发 |
| `TestAC5ExhaustedSignalWritesProviderHealthQuotaState` | 盘上落点 | DAO 收到 pa/m1 + class budget + state `exhausted`（逐字 `memory.QuotaExhausted`） |
| `TestAC5QuotaWiringEmptyConfigInventsNothing` | ⛔发明闸门 | 空／nil config＝零卡零 provider 零 alert 线 |

**反向判据（牙齿）**：① `remaining 0 仍派发` ⇒ `TestAC5RemainingZero…` 的 `first.calls!=0` 与 `err==nil` 双红、`TestAC5PlanCredit…` 的 plan 分支红；② **余额错** ⇒ 900 记账后 `RemainingMicros!=100` 红（含 plan/pay 串账形状：若 plan 消费落进 pay 桶，兄弟模型断言红）；③ 负数被吞 ⇒ `err==nil` 红＋余额 600 断言红；④ warn 不发事件 ⇒ 事件计数红；⑤ 盘上信号丢 ⇒ `w.got` 断言红。⑥ **链推进的独立牙在票 11 侧已有**（ChainRunner 的 ClassBudget 步进用例），本腿用错误类断言与之咬合——若 gate 改发未分类错误，`class!=budget` 当场红。

**零回归**：`go test ./internal/agent -count=1` 全包 `ok 4.441s`（36 枚 44-r1 名册＋r3 面全在包内）；`internal/observe`／`internal/memory` 本腿**零改动**；gofumpt `-l` 两枚输出空。

## 3. 具名欠账（本腿没做／不该本腿做的，⛔ 不藏）

1. **cmd/wisp 装配一寸**：链元素实际包上 `GateProvider`＋`QuotaWiring` 的调用点在 `cmd/wisp`＝**本腿点名禁面**（44-r3 已接价面）⇒ 产品今天仍不咬这一寸；接口名已备（`NewQuotaGate`/`QuotaWiring`/`GateProvider`/`Book`/`SetSignalWriter`）。
2. **per-model 消费的永久落库**：D35 `cost_daily`/`task_log` 无 model 维、DDL 冻结（44-v1 §7.1＝机主批准面）⇒ 分账余额是**会话级** ledger，盘上只有 `provider_health.quota_state` 信号；与 44-r2 文件头同一具名。
3. **realtime 轮（票 60）的同 gate**：票面 :33-34 点名「each realtime turn」也走派发前检查——voice_chain 侧⛔ 本票写面，归票 60 消费同一 `QuotaGate`。
4. `Ball/panel 渲染`归票 40（事件已上总线，EvWarn 复用 r3 出口）。

*本件由 44-r4 写腿产；只 commit、⛔ push。*
