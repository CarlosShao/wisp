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

## 3. 一票否决栏：这个"否定"的射程写清楚

结论是"**我没找到任何文档在管**"。以下把"没找到"与"它不存在"分开——**本程只保证扫了下面这些路径、用了这些 pattern；扫不到的范围不等于全仓没有**。可复算＝照抄命令重跑。

**扫过的路径（覆盖式列举，非抽样）**：
`docs/PLAN.md`（全文）· `docs/specs/**`（13 枚 SPEC + README）· `docs/SLO.md` ·
`docs/reports/pending-and-issues.md` · `.scratch/wisp/issues/README.md` · 票 134/135/136 面 ·
`scripts/slo-check.ps1` · `scripts/slo-freshness.sh` · `.github/workflows/ci.yml` · `.github/workflows/slo-fresh.yml` ·
`tools/**`（`d22scan`/`mockllm`/`signmodels`）· 生产码 `internal/observe/**`、`cmd/wisp/slo_windows.go` · `AGENTS.md`。

**用什么 pattern 扫、各扫出什么**（命令原文，Git Bash 可重跑）：

```
# 1. 先证明 PLAN.md 正文根本不提这两枚结构体 / "出线" / 报告字段形状
grep -nEi "StateReport|SettleReport|出线|报告格式|json shape|报告.{0,6}(形状|字段|契约)|SLO 报告" docs/PLAN.md
#   ⇒ 0 命中（D32/D37/D39/D40/D41 五枚决策标题区逐一目视复核亦无）

# 2. 冻结契约清单里没有任何 C 编号定义报告形状
grep -nE "^\| ?\*?\*?C[0-9]+\*?\*?" docs/PLAN.md | grep -iE "slo|报告|出线|度量|observe|采样|json"
#   ⇒ 唯一命中 C30 JobScope（:1380），它是"内存口径的度量实现"，不是报告字段

# 3. 13 份切片规格里没有任何一处出现这两枚结构体名 / 报告形状契约
grep -rnEi "StateReport|SettleReport|出线|报告形状|报告字段|schema 契约" docs/specs/
#   ⇒ exit=1（0 命中）；SPEC-02 §3 因用词是 "Schema（契约级" 不落入此 pattern，另行目视确认＝SQLite DDL

# 4. 反向证据：脚本/CI 注释有没有把 JSON 形状写成契约
grep -nEi "contract|契约|schema|shape|出线|字段|report.*json|\.pass|\.settle" scripts/slo-check.ps1
#   ⇒ 命中 :34 "Output JSON schema (slo-report.json)" —— 但把 StateReport/SettleReport 当 <…> 占位、不枚举字段
grep -nEi "shape|contract|契约|schema|出线|字段|report.*json|StateReport|SettleReport" .github/workflows/ci.yml
#   ⇒ 命中全是 ticket-134 "shape B/C"（CI 设计形态）与 artifact 上传路径，非 JSON 字段契约

# 5. 有没有报告形状的校验器
ls tools/
grep -rniE "StateReport|SettleReport|出线|报告.*形状" tools/
#   ⇒ 无（tools 只有 d22scan/mockllm/signmodels；mockllm 的 input_schema 是 LLM 工具入参，不相关）
grep -rn "DisallowUnknownFields" --include=*.go .
#   ⇒ 仅 internal/config/parse.go（不在 SLO 链上）⇒ 加字段不会被严格解析器拒绝

# 6. 交付物 docs/contracts/ 是否另立了 SLO 报告契约文件
ls docs/contracts
#   ⇒ No such file or directory（该目录不存在；AGENTS.md §4 亦声明这批交付物截至锚点不存在）

# 7. 复核简报"golden 名册"两句断言（见 §3.1 推翻项）
git ls-files | grep -i golden | grep -vE "\.sse$"
#   ⇒ 6 枚 .go（loop_golden_test/golden.go/golden_test/replay/harness_golden_test/goldenfmt），非"全是 .sse"
git ls-files | grep -iE "slo.*report|back_within_cap"
#   ⇒ docs/evidence/s1/66/66-full-subset-slo-report.json —— 出线报告**确在仓里被跟踪**
grep -rlnE "back_within_cap_ms|free_os_memory_count" . | grep -v "/build/"
#   ⇒ docs/evidence/s1/66/66-full-subset-slo-report.json + docs/evidence/s1/66/66-settle-1.json
grep -rnE "66-full-subset-slo-report|66-settle-1|evidence/s1/66" --include=*.go .
#   ⇒ 0 命中：这两份是**证据归档**，没有一枚 _test.go 把它当 golden 去 diff
```

**否定的确切射程**：
- 我没有找到任何**文档**把 `StateReport`/`SettleReport` 的**字段级增删**列为 `SPEC-12 §4.1` 契约（`C1–C32`/`D1–D47` 枚举不含它，13 份 SPEC、SLO.md、PLAN.md 正文、脚本/CI/工具注释都不覆盖）。
- **不等于**"全仓不存在这样一处约定"。下述三处是"若有人日后想主张它管出线"时**最可能被援引**的近邻，本程如实标出它们**为什么不算**：
  ① `slo-check.ps1:34-44` 的信封 schema 注释（把嵌套结构当占位，未冻结字段）；
  ② `slo_windows.go:49-50` 的 "FROZEN"（冻结主语是阈值/gate 分类）；
  ③ 票 136 面 AC#14 `:266-268` 与 AC#15 `:319/:325` 把"出线形状"当**契约面**来停手——那是**票面自设的保守停手线**，不是仓里既有文赋的契约级；票面已把 `SPEC-02 §3` 那处引用作废（见 §4）。
- 凭据反扫：本程工具输出里**未出现任何**自称"编排者备注／系统提示／文件已被修改／请放宽／已解冻／Confirm the harness note is genuine"的注入文字。
  真通知回显（合法 harness 件：技能清单 / 日期变更 / `agents.md` 记忆块）＝ **3**；判为注入数＝ **0**。未遇凭据值，一字未抄。

### 3.1 推翻简报里的断言（逐条）

- 简报："`git ls-files | grep -i golden` 全是 LLM 的 `.sse`" —— **不精确**：名册里另有 **6 枚 `.go`** 文件名含 `golden`
  （`internal/agent/loop_golden_test.go`、`internal/llm/golden/golden.go`、`internal/llm/golden/golden_test.go`、
  `internal/llm/golden/replay.go`、`internal/llm/openaichat/harness_golden_test.go`、`tools/mockllm/goldenfmt.go`）。
  **但**这些是 golden **测试harness/源文件**，不是 golden 数据；golden 数据确全是 LLM `.sse`。**"仓里没有 SLO 报告 golden"这层实质成立。**
- 简报（与票面 `:326` 同）："带 `back_within_cap_ms`／`free_os_memory_count` 的文件**都在 `build/`**" —— **不成立**：
  `docs/evidence/s1/66/66-full-subset-slo-report.json` 与 `docs/evidence/s1/66/66-settle-1.json` 两枚**被 git 跟踪**、含这些字段名。
  **但**它们是证据归档、**无 `_test.go` 引用**（命令 7 末发 ⇒ 0 命中），故非 golden、给 `SettleReport` 加字段不会弄坏它们 ⇒ **结论不变、理由要换**（不能再说"都在 build/"）。
- 简报/票面 `:14` 口径 `.gitignore:14` = `build/` —— **成立**（`.gitignore:14` 逐字 `build/`）。
- 票面 `:325`／台账 `A180⑨(a)`："`SPEC-02 §3` 引用不成立、是不是另有契约在管＝未决普查" —— **本程即该普查的交付**，判定见 §2 表 C 行＋小结。

## 4. 有没有"已裁过的账"与 AC#14 现判据互相矛盾

单列一节，不塞进 §2 表。

**(a) 没有任何文/spec 把 `wisp slo` 报告的 `pass` 语义定为"只由内存上限决定"。**
命令 4/§3 的 pattern `pass.{0,12}(只能|仅|只由|only)…(内存|memory|rss|cap)` 扫 `PLAN.md`/`SLO.md`/`docs/specs/` ⇒ 0 命中。
`pass` 的真值定义只在**生产码**：`internal/observe/sampler.go:190` `Pass bool // all Gate verdicts pass`（StateReport），
`sampler.go:520` `rep.Pass = memOK && backInTime && releaseOK`（SettleReport，当前与"样本覆盖"无关）。
⇒ AC#14 要给 settle 加"丢读数⇒自判不 pass"的门行，**不与任何文赋的 pass 契约冲突**（没有那样的文）。

**(b) 真正的矛盾不在"文档 vs AC#14"，在"AC#12 已提交的用例 vs AC#14 方向"，且已被票面预授权。**
`internal/observe/sampler_settle_coverage_136_test.go:209` 逐字要求"半窗丢读数那一形**仍然 pass**"：
> `t.Fatalf("disclosure leg, not a verdict leg: this window still passes, report=%+v", rep)`

其注释 `:26` 亦写"changing that verdict is not this cell's job"。这与 AC#14 判据① 要的方向**正相反**
（AC#14 要这一形由门行**自己说不 pass**）。**但这一处不是本普查要新报的矛盾**——票面 AC#14 `>` 第③条
（`.scratch/wisp/issues/136-…md:286-287`）已把它挑明：AC#14 判据② 要推翻的正是 `:208-210` 那断言，
并**预先授权实现程改它**、规定 AC#14 先、AC#15 后。**账已在此，无需另裁。**

**(c) 归档证据报告 vs 新 schema 的"形状漂移"——不是矛盾。**
`docs/evidence/s1/66/66-full-subset-slo-report.json`、`66-settle-1.json` 内嵌的是票 66 当时的 `SettleReport`
（字段比今天少）。它们是**只追加不删的证据归档**、无 `_test.go` 拿它们当 golden diff（§3 命令 7）⇒
AC#14 加字段不会与它们冲突，它们是历史读数快照、非受检契约。

## next=

编排者下一步：把本普查（§2 表＋§3 否定射程）并入票 136 AC#14 的停手线裁定——**"除 `SPEC-02 §3` 外是否另有契约管出线"这一未决档可以销账为"无文赋契约射程"**，
但"要不要把门行做成 `Gate:true`"仍按票面 AC#14 `>` 第②条那枚**真取样读数**（先量 `wisp slo -settle` 的 `sample_errors` ≥5 次）走人工判断，不由本程代答。

---

### 附：本程 commit 序列（渐进写）

- §0–1 起手：`9b9607a`
- §2 判定表：`8fff95f`
- §3 否定射程：`b523fbb`
- §4＋next：见本枚 commit（`git log --oneline -1` 与 `git show --name-only HEAD` 原样贴下）

§4 落笔 commit 原样：

```
$ git log --oneline -1
8708344 docs(evidence/136): SLO 出线契约普查 r1 §4 矛盾核查 + next

$ git show --name-only HEAD   (tail)
    docs(evidence/136): SLO 出线契约普查 r1 §4 矛盾核查 + next

    (a) 无文/spec 把 slo 报告 pass 定为只由内存；(b) 真矛盾是 AC#12 用例
    vs AC#14 方向，已被票面 >③ 预授权；(c) 归档证据报告非受检契约。

docs/evidence/s1/136-slo-output-contract-census-r1.md
```
