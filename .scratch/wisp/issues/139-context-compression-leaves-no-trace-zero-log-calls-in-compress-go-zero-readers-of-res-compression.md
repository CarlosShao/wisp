# 139 — 上下文被压缩掉多少，**今天连一条日志都没有**：`compress.go` 全文零条日志调用、`res.Compression` 生产侧零读取者

**Status:** ready-for-agent（2026-09-24 10:0x 编排者建；来源＝owner 批准的 `Q-43` 后半，
              原始缺口＝`docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md` **GAP-18**，
              裁决见 `docs/reports/2026-09-24-gap-analysis-audit-verdict.md` §2 表第 18 行。
              ⚠ 裁决表把这一条**拆成了两半**：**只有 Go 侧那半属本票**；
              「面板上画一条占用分解条」那半**不做**——见 §3。）

**Packages:** `internal/agent/**`（`compress.go`、`loop.go` 的 `Result` 出线）· 若需落盘面则 `internal/observe/**`
              （**注意地界**：`internal/observe/**` 此刻可能有别的程在跑，动手前先看工作树与在飞名单）。
              **⚠ 冻结面照旧禁改**：D32 阈值 · `thresholds.go` · 任何 golden · `internal/risk/**` ·
              `tools/d22scan/**` · `frontend/**` · `docs/PLAN.md` · `docs/specs/*.md`。

## 0. 编排者现量的四条锚（我自己走的，其中一条**比审计代理报的更强**）

| 事实 | 读数 |
|---|---|
| 压缩调用点 | `internal/agent/loop.go:393` `if l.comp.Need(hist)` → `:396` `l.comp.Compress(...)` |
| **成功**时做了什么 | `:399-401` —— `l.replaceHistory(nh)` ＋ `res.Compression = rep`。**没有一条日志。** |
| 失败时做了什么 | `:398` `l.log().Warn("agent: history compression failed", "err", err)` —— 只有这一条，**在调用方**。 |
| `Compression` 有人读吗 | `git grep -n '\.Compression\b' -- '*.go'` 排除 `_test.go` ⇒ **全仓唯一命中就是赋值那一行 `loop.go:401`**，**零读取者** |
| `compress.go` 自己有日志吗 | ⚠ **零条**（`git grep -E 'Info\(|Warn\(|Debug\(|Error\(' -- internal/agent/compress.go` 无输出）。⇒ **修正审计代理那句"只在失败 Warn"**：那条 Warn 不在 `compress.go` 里，这个文件**从头到尾不产生日志**。 |

⚠ 另外两条**留给 AC#1 现算**、我未复算：
① `internal/agent/budgets.go:36` 的 `refHistoryTokens` 是否按 `context_window` 等比缩放（裁决表这么报的）；
② 那条 `DEFERRED(D28-1)` 注释（就在 `loop.go:394`）在 `SPEC-12 §5` 登记表里的五字段条目是什么、
   它的"残缺表现"栏有没有覆盖今天这件事。

## 1. 为什么现在立案

GAP-18 的原文说「压缩是静默的、用户看不到上下文占用」。裁决表判定这条**比原文写的还严重**：
不是"没画出来"，而是**成功路径连留痕都没有**——所以"看不到"的根因不在前端，在**根本没有数据可看**。
⇒ 于是**正确的第一跳不是画 UI，是让这一件事在 Go 侧变得可回答**。这也正是裁决表把它拆两半的理由。

顺带一条**已在册的**：`loop.go:394` 的 `DEFERRED(D28-1)` 明写着"这套同步回退是 S1 接受的形状、
Warm-window hook 才该拥有这次调用"。⇒ **本票不改这个决定**，只给它补上"这一次到底压了多少"的读数。

## 2. 结案判据

- [ ] **AC#1**（把"零留痕"量成读数，并顺手核 §0 那两条未复算项）
      ① 现算 `Compression` 的读取者名册（含 `cmd/wisp/**`、面板桥、诊断包 `internal/observe/diagnostics.go`），
         **逐名**列，不许只给计数；
      ② 现算 `refHistoryTokens` 与 `context_window` 的关系；若确为等比缩放 ⇒
         **本票所有例子一律不许引用 `12000` 这个数当"实际触发点"**（裁决表已判定它是 spec 原文、非现值）；
      ③ 把 `DEFERRED(D28-1)` 的登记条目原文引出来，回答一句：**"压缩无留痕"这件事今天是否已被那格覆盖**；
         若已覆盖 ⇒ 本票**并入那格的落地点**、不新立不变式；若未覆盖 ⇒ 去 `SPEC-12 §5` 按五字段补登记。
- [ ] **AC#2**（最小落盘面）**只要求"事后答得出"，不要求任何 UI**：
      压缩**成功**时留下一条可核的痕（结构化日志字段／`Result` 出线上的可读字段／诊断包里的一项，三者选一，
      并在票面说明为什么选它），内容至少含：**压前 token、压后 token、是否真的动了历史**。
      ⚠ **禁止**新增任何 Go→前端事件或面板方法：那要改 `SPEC-08 §5.2` 的枚举白名单＝**契约变更**（裁决表判定项），
      且面板那半已由 owner 交给外部 agent（A102）。碰到这条 ⇒ **停手上报**。
      ⚠ **禁止**把"原文内容"落进日志（`[privacy] keep_transcript` 那条硬约束）——只许记**计数**。
- [ ] **AC#3**（牙）一枚用例钉住"压缩成功 ⇒ 痕必须存在"；
      反向判据：**把这条痕摘掉，该用例必须转红并逐名报出红行**。
      先证变异落地（`grep -n` 原文 ＋ `go build` rc=0）再读数。不许 `t.Skip`、不许放宽断言。
      ⚠ 写判据前先问一句"这一发**今天**响不响"：若在未修码上就已经全响，这格是装饰，
      要另造一发"今天不响、修了才响"的（恒真判据＝本仓登记过的一类新假绿）。
- [ ] **AC#4**（门禁）同 138 AC#4：`gofmt`／`gofumpt`（⚠ **用宿主现成的 v0.12.0 binary，
      绝不执行 `go install …@latest`**）／`go vet` 宿主 rc=0／`sh scripts/d22scan.sh` rc=0 各 scope 不降／
      `-count=2 -v` 两形四数＋**名册差集**。`internal/observe/**` 若在飞，**先让路**再采读数。

## 3. 明确不在本票范围

- **面板上的占用可视化**（分解条、点开看原文）：正解落 `frontend/**`，owner 09-23 已收走交外部 agent；
  且"点开看原文"会顶到 `[privacy] keep_transcript=false` 那条硬约束。**退回，不做。**
- 压缩算法本身、`Need()` 的阈值、D15 的预算分配：**一律不改**（改这些＝契约变更）。
- `DEFERRED(D28-1)` 那个 Warm-window hook 的落地：另案（票 28 那条），本票只补读数。

## 4. 进度日志

- [2026-09-24 10:0x +08] agent=orchestrator did=建票。
  来源＝owner 在对话里批准 `Q-43`（原话「都按你的推荐来」）；撤销口令＝**「139 撤」**。
  自己现量 §0 五条锚，并**修正了审计代理的一处措辞**：`compress.go` 不是"只在失败 Warn"，
  而是**全文零条日志**（那条 Warn 在调用方 `loop.go:398`）——这条使缺口**变深**，故照实写进票面。
  两条未复算项（`refHistoryTokens` 等比缩放、D28-1 登记覆盖度）**没有当事实写**，已落成 AC#1。
  next＝排在 138 之后、今日队列之后。
