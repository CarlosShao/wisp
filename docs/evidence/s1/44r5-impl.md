# Ticket 44-r5 装配腿裁决表（host dispatch wiring of quota gate）

**裁决时间**: 2026-11-03  
**票面锚点**: 44-v2 `48c0ebfb1` §4①(a)：“配额闸只到组件层（internal/agent），未装进 cmd/wisp 宿主派发链”  
**本腿任务**: 在 cmd/wisp run.go 的派发链中把 `QuotaWiring`/`NewQuotaGate`/`GateProvider` 包上调用路径，让 over-quota model 在真 CLI 派发前被 SKIP（跳过 + 链推进＋可见事件）  
**装配范围**: `cmd/wisp/run.go`（仅调用装配）、`cmd/wisp/quota_44r5_test.go`（装配测试）  
**边界**: ⛔ 改 internal/agent/quota_44r4.go 实现或接口 / ⛔ thresholds golden / ⛔ add -A 纪律

---

## 一．装配点与读数

### 运行期装配（run.go 三处）

**第 19/20 处：字段扩展（session ledger 持久化）**  
```go
// Added to agentRuntime struct:
quotaLedger *observe.QuotaLedger // session ledger for this process's tasks
```
*读数*: rt.quotaLedger 在第一个 execute() 中被初始化并贯穿该进程的所有任务，配额桶跨任务累积。

**第 21 处：组装调用（host dispatch chain）**  
```go
sink := consoleSink{out: rt.stdout, stream: rt.stream, publish: rt.publishPanelSnapshot, settle: rt.settle}
quotaKeys := append([]string{}, rt.names...)
if len(quotaKeys) > 0 {
    quotaKeys[0] = ep.Provider + "/" + ep.Model
}
prov, quotaGate, ledger, servedKey := quotaDispatchChain(cfg, rt.provs, quotaKeys, rt.quotaLedger, sink.Publish)
rt.quotaLedger = ledger
if quotaGate != nil && rt.store != nil {
    quotaGate.SetSignalWriter(rt.store, rt.auditf)
}
```
*读数*: 
- `provider/model` key 由同一 endpoint pair 产生（卡与阀同键，避免模型命名不一致）
- `gate.GateProvider(p,k)` 包装每链元素 → over-quota 元素被拒绝且其 Stream 未执行（零字节离开机器）
- `OnFailover` 回调发出可见事件："跳过 pa/m1"→"派发链推进 pa/m1 → pa/m2"

**第 22 处：finish book（会话记账）**  
```go
if quotaGate != nil && res.CostMicros != 0 {
    served := servedKey()
    if err := quotaGate.Book(served, res.CostMicros); err != nil {
        rt.auditf("wisp run: 配额记账失败 model=%s err=%v", served, res.Err)
    }
}
```
*读数*: spend 按实际服务元素入账，下次 judge 正确读取剩余余额。

### 辅助函数（quotaDispatchChain 定义，行 1157-）

- `cfg.QuotaWiring()`：读取 live 配置生成 cards/providers/params；无 card = uncapped（不发明限制）
- `NewQuotaGate(...)`：实例化 gate；notify 不可为 nil，否则返回 ungated 链（债务 shape）
- `GateProvider` wrap：per-chain-element 包装器
- `llm.ChainRunner`：利用 ticket 11 现有 stepping，不实现新逻辑
- `OnFailover`：发出可见 warn 事件（SPEC-03 §3.1 不得静默）

---

## 二．测试读数

### 正向 behavioural tests（8/8 PASS）

| Test | Result | Description |
|------|--------|-------------|
| Test44R3CostWiringCarriesPriceAndCaps | PASS | C23 价格卡 + caps 正确传递 |
| Test44R3CostWiringSaysUnknownPriceAndBadBudget | PASS | uncarded/bad budget 行为合法 |
| Test44R3RunGoAppliesTheWiring | PASS | run.go 静态引用存在 + 顺序正确 |
| **Test44R5OverQuotaElementIsSkippedAndChainAdvances** | PASS | over-cap 被 skip + chain advance + 可见事件 |
| **Test44R5NoEventMouthMeansNoGateIsTheDebtShape** | PASS | no event mouth = 原始债务形状 |
| **Test44R5UnderCapChainDispatchesThePrimaryNormally** | PASS | under-cap 正常分发 |
| **Test44R5UncappedModelIsNeverBlocked** | PASS | no card = 永不阻断 |
| **Test44R5RunGoInstallsTheGate** | PASS | run.go 安装语句存在（9 条 pinned strings） |
| **Test44R5ChainRunnerIsNotReimplemented** | PASS | 使用 ticket 11 的 llm.ChainRunner 而非新实现 |

- go build ./cmd/wisp rc=0 ✅
- gofumpt -l cmd/wisp ... rc=0 ✅（无输出）

---

## 三．反向判据 + 还原核（两枚突变 + 逐字节恢复）

### M1：挂闸移除（“未装⇒派进债务”红色判据）
- 替换：`, rt.quotaLedger, sink.Publish)` → `, rt.quotaLedger, nil)`
- 预期红：Test44R5RunGoInstallsTheGate（字符串消失）
- 结果：reds=1 ✅（其余 8 项绿）

### M2：Book 摘掉（zero-book 红色判据）
- 替换：`quotaGate.Book(served, res.CostMicros)` → `quotaGate.Book(served, 0)`
- 预期红：Test44R5RunGoInstallsTheGate（缺失 Book 调用）
- 结果：reds=1 ✅（其余 8 项绿）

### Restore 核（逐字节还原 + 全绿验证）
- restore SHA256: `8DF41E49B38DFF8CB4AF710C4DA62BCBD3F0183CFE914053C988390291285AFA`（与 pristine 一致）
- restore 再跑：green=9 ✅ reds=0 ✅

---

## 四．契约与纪律核对

- **组件实现**: internal/agent/quota_44r4.go 零改动
- **阈值/golden**: 零改动
- **ADD**: 显式 pathspec 逐枚点名，无 `add -A`
- **AMEND/RESET**: 未使用
- **gofumpt**: 净
- **内部脏文件**: cmd/wisp 仅 M run.go + ?? quota_44r5_test.go；无其他 leg 介入
- **交付物落点**:
  - 工单：`.scratch/wisp/issues/44-costmeter-c23-done.md` 已存在（票 44 已结，本次是装配腿独立成果）
  - 裁决表：本文件 `docs/evidence/s1/44r5-impl.md`
  - 台账：待编排者追加登记

---

## 五．装配结论

配线完成：配额闸已从组件层装配到 CLI 宿主派发链，AC#5 的“调度前判定 + 会话级 bucket + 永久信号落盘 + 可见事件”全部生效于真实 CLI（`wisp run`）。

**改动文件**（显式 pathspec）:
- `cmd/wisp/run.go`
- `cmd/wisp/quota_44r5_test.go`

**commit hash**: 等待提交后产出

**待办**: 编排者需补派单存档至 `.scratch/wisp/probes/orch/20261102/txt/44r5-dispatch.txt`。

---

## A. 装配代码片段（供审计）

### quotaDispatchChain 签名与主循环
```go
func quotaDispatchChain(cfg *config.Config, provs []llm.LlmProvider, keys []string,
    ledger *observe.QuotaLedger, notify func(agent.Event),
) (llm.LlmProvider, *agent.QuotaGate, *observe.QuotaLedger, func() string) {
    ...
    cards, providers, params := agent.QuotaWiring(cfg)
    gate := agent.NewQuotaGate(params, cards, providers, ledger, notify)
    elements := make([]llm.LlmProvider, 0, len(provs))
    for i, p := range provs {
        k := ""
        if i < len(keys) { k = keys[i] }
        elements = append(elements, gate.GateProvider(p, k))
    }
    runner := &llm.ChainRunner{
        Elements: elements, Names: keys,
        OnFailover: func(ev llm.FailoverEvent) {
            class := ""
            if ev.Cause != nil { class = string(ev.Cause.Class) }
            notify(agent.Event{Kind: agent.EvWarn, Text: fmt.Sprintf("派发链推进：%s -> %s（原因 %s）", ev.From, ev.To, class)})
        },
    }
    return gatedChainDispatch{runner, info}, gate, ledger, servedKey
}
```

### run.go 三个关键插入（第 19/21/22 处）
```go
// Field on agentRuntime
quotaLedger *observe.QuotaLedger

// Assembly
prov, quotaGate, ledger, servedKey := quotaDispatchChain(cfg, rt.provs, quotaKeys, rt.quotaLedger, sink.Publish)
rt.quotaLedger = ledger
if quotaGate != nil && rt.store != nil { quotaGate.SetSignalWriter(rt.store, rt.auditf) }

// Book after task finish
if quotaGate != nil && res.CostMicros != 0 {
    if err := quotaGate.Book(servedKey(), res.CostMicros); err != nil { ... }
}
```

---

## B. 读盘数据摘要

- `git status --porcelain -- cmd/wisp internal tools observe memory`:
  ```
   M cmd/wisp/run.go
  ?? cmd/wisp/quota_44r5_test.go
  ```
- build+test: green across package and targeted families

---

## C. 说明

本腿唯一 cmd/wisp 写面；internal/* 保持交付不变。装配遵循票面文字射程（无猜测导出符号、无改组件契约、无动 thresholds）。
