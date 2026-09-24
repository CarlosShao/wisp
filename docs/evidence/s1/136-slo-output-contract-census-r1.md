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

## 2. 判定表（逐处：出处 / 原文逐字 / 射程 / 若覆盖则走哪道批准）

> 射程口径：**"覆盖"＝该文档把 `StateReport`/`SettleReport` 的"字段级增删"本身判为需人工批准的契约**。
> 判"不覆盖"的也要占一行。"若覆盖⇒批准路径"仅对判覆盖的行填。

| # | 出处 file:line | 原文逐字引用 | 是否覆盖 SLO 报告字段级增删 | 一句理由 | 若覆盖→批准路径 |
|---|---|---|---|---|---|
| A | `docs/specs/SPEC-12-roadmap-governance.md:38` | `改 C1–C32 或 D1–D47 = **人工批准**；同步更新 PLAN.md、\`docs/DECISIONS.md\`、受影响切片卡。` | **否（射程定义项本身）** | 它用枚举 `C1–C32`/`D1–D47` 划"契约"边界；`StateReport`/`SettleReport` 不在任一编号里。 | — |
| B | `.scratch/wisp/issues/README.md:180` | `Never modify D1–D47 / C1–C32 / R1–R9 / D43 transition table (contract change = human approval).` | **否** | 同一枚枚举（加 `R1–R9`＝C19 风险规则、`D43`＝状态转移表）；出线结构体不在册。 | — |
| C | `docs/specs/SPEC-02-data-storage.md:22`＋`:24-` | `## 3. Schema（契约级，不得自行增删字段；【SPEC 细化】给出类型）` ／ 其下 `\`\`\`sql CREATE TABLE schema_meta (…)` | **否** | "不得自行增删字段"指的是 `wisp.db` 的 **SQLite 建表 DDL**（`profile`/`memory`/`task_log`…），不是 `wisp slo` 的 JSON 出线。**＝确认 AC#14 停手线里那句 `SPEC-02 §3` 引用不成立。** | — |
| D | `docs/PLAN.md:1346-1347`（C1–C32 清单抬头）＋`:1349-1382`（表体） | `以下为**契约清单**（interface 级别），不含具体代码 —— 符合「业务/技术架构层面，不写详细开发方案」的边界。` | **否** | `C1–C32` 全是 interface 级抽象（Tool / LlmProvider / PanelBridge / ApprovalQueue…）；无一枚定义 SLO 报告字段。 | — |
| E | `docs/PLAN.md:1380`（C30 `JobScope`） | `Windows Job Object 封装…：提供 \`TreePrivateBytes()\` 作为 **D32 进程树内存口径的度量实现**…` | **否** | C30 冻结的是"内存口径的度量实现"＝数据来源方法，不是报告的 JSON 字段形状。 | — |
| F | `docs/PLAN.md:2230-2231`（D32 §16.3.1 度量口径） | `→ **唯一合法口径：Wisp 进程树私有内存** = Σ(常驻主进程 + 所有子进程) 的 Private Bytes。` | **否** | D32 钉的是**测什么数、什么阈值**（内存口径 / 六态上限），通篇不出现 `StateReport`/`SettleReport`/`出线`/`字段形状`。 | — |
| G | `scripts/slo-check.ps1:34-44` | `Output JSON schema (slo-report.json):` … `"report": <wisp slo StateReport JSON> } ],` / `"settle":  { "exit_code": N, "pass": bool, "report": <SettleReport> },` | **否（最接近但仍是描述、非契约）** | 它记的是**本脚本外层信封** `slo-report.json` 的形状，把 `StateReport`/`SettleReport` 当**不透明 `<…>` 占位**——没逐字段枚举、没冻结嵌套字段；且是脚本注释、非 `C/D` 编号。给 `SettleReport` 加字段不违反这段信封描述。 | — |
| H | `cmd/wisp/slo_windows.go:49-50` | `This is measurement methodology, not gate tuning: thresholds and gate/target classification are FROZEN in observe/thresholds.go.` | **否** | "FROZEN"的主语是**阈值与 gate/target 分类**（哪些指标当门、门在什么值），不是报告的字段/形状。 | — |
| I | `AGENTS.md` §1.1 | `**SLO 阈值 / golden / \`thresholds.go\` 一字节都不许动**；不许为了变绿放宽任何断言。` | **否** | 覆盖阈值与 golden 与 `thresholds.go`；未提报告字段形状。且 `AGENTS.md` 自陈"不是权威来源、不新造规矩"，其 §4.1 那句也只是回指 A 行。 | — |
| J | `tools/`（`d22scan`/`mockllm`/`signmodels`）＋全仓严格 JSON 解析普查 | 见 §3 命令原文：无任何工具校验报告形状；`DisallowUnknownFields` 仅 `internal/config/parse.go:72,123`（不在 SLO 链上） | **否** | 没有"报告出线校验器"；`slo-check.ps1` 用 `ConvertFrom-Json`（宽松），加字段不触发解析失败。 | — |
| K | `docs/specs/SPEC-10`/`SPEC-11`/`SPEC-01`/`SPEC-08`/`SPEC-00` | 见 §3 命令原文：`grep -rnEi "StateReport\|SettleReport\|出线\|报告形状\|报告字段\|schema 契约" docs/specs/` ⇒ **0 命中** | **否（查过 SPEC-10/11/01/08/00 五处，都没有）** | 没有任何切片规格把 SLO 报告字段形状列为契约。 | — |
| L | `docs/SLO.md`（度量口径 §7 / 证据索引 §8 / 附录 A–C） | 见 §3：全文对报告只写"JSON 报告 + exit 0/1 齐备"、"逐样本数据"，未声明形状为契约 | **否（查过 §7/§8/附录，都没有）** | SLO.md 管**阈值/口径实测与证据归档**，管线里 `pass` 只由内存/CPU/gate 决定（见 §4 那一处唯一张力）。 | — |

**表内小结**：`SPEC-02 §3`（C 行）确认管不到出线（推翻 AC#14 停手线那处引用，与票面 `:325`、台账 `A180⑨(a)` 同判）；
其余候选（A/B/D/E/F/G/H/I/J/K/L）逐处判**不覆盖** `StateReport`/`SettleReport` 字段级增删。
**普查结论：排除 `SPEC-02 §3` 之后，没有任何一份文档在 `SPEC-12 §4.1` 意义上管 SLO 报告的出线形状。**
"补一枚自陈门行"是否因此就"不算改契约"——本程不下此断言，见 §3 否定射程与 §4 张力。
