# 136-SLO 出线契约普查 r1 — 「有没有一份文档在管 `wisp slo` JSON 报告的出线形状」

> 本程＝只读普查代理 `auditor-ticket136-contract-census`。
> 唯一要答的问题：**`StateReport`／`SettleReport` 结构体字段级的增删，属不属于 `SPEC-12 §4.1`「改契约＝人工批准」的射程？**
> 本程**一枚测试／构建／容器／`wisp` 都不跑**（此刻 `cmd/wisp`／取样面被另一枚写码程占着做真取样），
> 只读文件、只写本枚证据文件、只 commit 不 push。
> 本程**不下「要不要批准」的结论**——只找原文、判射程。要不要走人工批准由编排者/owner 定。

## 0. 锚点与取数（原样）

- 本程锚点 sha（`git rev-parse --short HEAD`）：`ca2c55e`
- 分支：`dev`
- 取数时刻（`date`）：`Thu Sep 24 18:50:10 CST 2026`
- `git status --porcelain`（原样，工作树里只有 owner 自己的 `design/**` 未提交改动，本程一枚没碰）：

```
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 D design/index.html
 D design/screens/approval.html
 D design/screens/ball.html
 D design/screens/chat.html
 D design/screens/config.html
 D design/screens/cost.html
 D design/screens/firstrun.html
 D design/screens/palette.html
 D design/screens/privacy.html
 D design/screens/security.html
 D design/screens/states.html
 D design/screens/tasks.html
?? design/doubao/
?? design/old/
```

> 复核简报断言：简报说"工作树里只有 owner 自己的 `design/**` 未提交改动"——**成立**。全部改动落在 `design/**`
> （16 枚 tracked 文件被移走 + 两枚未跟踪目录 `design/old/`、`design/doubao/`）。票 136 面 `:484` 已逐字登记这处
> 是 owner 侧的挪动、"不是我动的"。本程只往 `docs/evidence/s1/` 追加一枚新文件，`design/**` 一字未动。

## 1. 被普查的对象是谁（先钉住"出线形状"具体指哪两枚结构体）

`wisp slo` 的 JSON 出线 = Go 结构体序列化，字段定义在两处（生产码，非任何文档）：

- `internal/observe/sampler.go:158` `type StateReport struct { … }`（含 `Verdicts []Verdict`、`Pass bool // all Gate verdicts pass` 等）
- `internal/observe/sampler.go:431` `type SettleReport struct { … }`（含 `FreeOSMemoryCount`、`BackWithinCapMS`、票 136 AC#12 已加的 `SampleErrors`／`LastSampleError`）

票 136 AC#14 要"给 `SettleReport` 补一枚自陈门行（新字段）"——就是往 `:431` 这枚结构体加字段。

`SPEC-12 §4.1` 逐字（`docs/specs/SPEC-12-roadmap-governance.md:38-39`）：

> 改 C1–C32 或 D1–D47 = **人工批准**；同步更新 PLAN.md、`docs/DECISIONS.md`、受影响切片卡。
> agent 单方面改契约 = 跑歪模式 #1，对抗验收判失败。

所以"改契约＝人工批准"的射程**由 §4.1 用枚举定义**：只有落在 `C1–C32`／`D1–D47`（及其行为契约）里的东西才是"契约"。
普查＝逐处问：`StateReport`/`SettleReport` 的**字段形状**有没有被塞进这些枚举集，或被别的文以"契约级／冻结"的名义钉住。
