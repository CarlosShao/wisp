# 173 — 审计行里的 `in_allowlist_scope` 是枚**两态字段装了两种意思**：没有任何路径被判定过时，它写 `true`（"判成在范围内"），而真值是"没判过"（来路＝票 162 的 AC#5 验收程现量＋本票读码）

- Status: **ready-for-agent**（写面＝`internal/tools/bridge.go` 那一味 + 它的用例；⚠ 若判"审计行格式属冻结契约"⇒ 停手上报，见 AC#3）
- 来源：**`162-v3` 非实现者验收表**（`docs/evidence/s1/162-risk-tier-accept-v3.md` §14 next= 第 (c) 条）＋**本票编排者读码定位**（`internal/tools/bridge.go:923-932`）。台账 `A339`。
- 关联：票 162（暴露处，那票已结线，本票是它留下的账）· 票 147（"红句里的 offset 恒 0"＝**同族**：取证字段说了一句它没资格说的话）· 票 105／102（审计行与 `Decision` 的来源链）

## 这是什么（人话）

每次工具调用会落一行审计记录，里面有一栏叫 `in_allowlist_scope`（"这玩意儿在不在允许范围内"）。它的算法只有两句：

```
if b.paths == nil || len(dec.Paths) == 0 { return true }     // bridge.go:924-926（本轮现量）
```

⇒ **当这次调用根本没有任何路径被判定过**（比如路径是空白的、在送到法官之前就被丢掉了），**这一栏照样写 `true`**。验收程抓到的原样读数就是这个形状：

```
... risk=L1/2 decision=allow outcome=error rules_hit=[R1] in_allowlist_scope=true ...
```

**"没有东西可判"被渲染成了"判过、而且在范围内"。**

## 为什么值得做／不做的后果

这行文本是**人查事故时第一眼看的字段**。以后有人拿它做过滤（"给我看所有越界的调用"），**这一类调用会被静默地排除掉**——不是报错、不是漏一条，而是**看起来像已经查过且没事**。同族的形状我们在票 147 已经吃过一次（红句里的 offset 恒 0，读起来像"文档停在第 0 字节"）。

## 我的推荐（**owner 不用选，除非他要说"别做"**）

**改成三态，别改成"更保守的 false"**。理由：`false` 说的是"判过、且不在范围内"，那也是一种断言，把"没判过"写成 `false` 只是换一句假话。推荐 **`true` / `false` / 不写这一栏或写 `unknown`**（选一形，并写明为什么），并在**真的判过**的那条路上一个字都不许变。
**不按推荐来的坏处**：图省事改成 `false` ⇒ 所有"没路径可判"的调用（`ping`、`memory.search` 这类本来就没有路径的）明天起全变成"越界"，那会把这栏字段的信号彻底淹掉——**这比现在更糟**。

## AC（每格都要答"这一发在**未修码**上响不响"）

- [ ] **AC#1（正面）**：造一发"路径空白、没有任何路径被判定过"的调用 ⇒ 未修码上审计行**今天写 `in_allowlist_scope=true`**（先把这条读数钉住，可复算）；修完之后**这一栏必须不再是 `true`**（不再写、或写成你选的第三个值）。
- [ ] **AC#2（对照，不许放宽）**：同时保两发既有形状**一字不变**——① 真在范围内 ⇒ 仍 `true`；② 真越界（走 R2 那一族）⇒ 仍 `false`。⚠ 摘掉任一分支对应用例必须红；答不出＝这一格是装饰。
- [ ] **AC#3（射程自查＝未定义即停）**：先查清**这一栏的字面格式属不属于冻结契约**（`SPEC-02 §3` 的 `tool_call` DDL／审计行的"逐字渲染"许诺那一族）。⚠ **若判"属契约面"⇒ 停手上报**，不要自己改规格文本来让它变得可改（票 134／`AGENTS` §1.1 那一族）。
- [ ] **AC#4（分母，别只修这一枚）**：全仓普查"**没有任何判定却照样落判定值**"的审计/取证字段还有几枚（本轮已知同族候选＝票 147 那枚 offset 恒 0；以及 `risk=L1/2` 这种斜杠拼接在两值时的读法）。**逐枚点名，不许只修手里这一枚**（票 171 AC#2 那条"私货尺"的教训）。修不修分别给理由。
- [ ] **AC#5（契约轴）**：不动 `docs/PLAN.md`、`docs/specs/**`、`tools/d22scan/**`、`allowlist.txt`、`thresholds.go`、golden、`frontend/**`、`design/**`；不新增豁免。

## 本票**不**解决

- 不改"空白 path 该不该被法官判"这个行为本身（那是审批语义，属契约面）——**本票只管这一栏字段别再替不存在判定说话**。
- 不重开票 162（它的八格都已结线）。
- 不动票 147 的 offset 那一枚（本票 AC#4 只负责**把它从"没人记得"变成"在册"**，改不改归那票）。

## 派单前置条件

① `internal/tools/**` 没有别的程在飞（本票写面在 `bridge.go`）；② 别在 `probes/154/**`／`probes/161/r6/**` 被别的程改时跑那两把尺；③ 台件 `gofumpt` 干净。

## Progress log (append-only, newest last)

- [2026-09-27 18:2x] 编排者开票。来路＝`162-v3` §14 第 (c) 条，它判"不归 162、也不该在 162 造第二把尺（AC#4b 禁）"——**这句我采纳**；但它建议"并入 `internal/risk` 审计字段语义票"，而**全仓没有那样的票**，所以落成票 173。我本轮自己读了定位那一味：`internal/tools/bridge.go:923-932`（`inScope`），空路径那支在 `len(dec.Paths)==0` 处**直接 `return true`**，注释自陈"the audit line's cheap answer"⇒ **写代码的人当时就知道它是便宜答案，只是没人给它登记过**（这条与台账里"错归因会让后人删掉真正兜底的那行"同族，所以票面把注释原话抄了出来）。


- [2026-09-27 20:3x] **编排者追加同族第二枚（本票射程扩一格，原句一字不抹；来路＝只读测量程 `174-c1`，表 `docs/evidence/s1/174-artifacts-allowlist-readback-c1.md`，提交 `53cbf0bf`）**：本票讲的是"空白路径到不了法官、审计行却写 `in_allowlist_scope=true`"——**同一枚族、不同字段**还有一发：**审批超时那一发的 `ErrorClass` 也叫 `user_rejected`**（`internal/tools/bridge.go:397-404`，人工拒与超时两路**共用同一枚 reject class 实参**；〔第二见证＝`174-c1` 现跑两发各自拿到 `class=user_rejected`，逐字读数在它表的"向一"段，**本程我未复跑**〕）。⇒ **审计与回执上"用户拒绝了"与"没人来得及回答"是两件事**，混成一枚标签会让下一次"谁拦的"这笔账读不出真因。⇒ **处置**：与本票主项同一形状（要么分两枚 class，要么回执里带上"超时"这一味），⚠ **不许**借这一格去动 `thresholds.go` 或任何审批超时常量（那是冻结件）。

- [2026-10-13 22:1x] **落地腿 `173-r1` 交件**（提交 `ec2f0864`＝`internal/tools/bridge.go`＋新用例两件；本 bullet 与 `probes/173/r1/**` 走第二枚提交）。**选形＝`true`/`false`/`unknown`**：不写这一栏被否——同行其它栏都保列位、`unknown` 是正向可 grep 的"无判定"claim，缺栏与"旧版二进制／日志截断"不可分；`false` 按票面原话拒绝（那是第二句假话，且把无路径调用全打成越界）。**AC#1 未修现量（两发都红、逐字）**：空白 path `fs.edit` ⇒ `tools: call kind=error task=task-1 corr=corr-1 tool=fs.edit risk=L2 decision=allow outcome=error rules_hit=[R1] in_allowlist_scope=true grant_id=0 reason="R1: 工具声明为下界（L2）"`；无 path 参数工具（`probe.nopath`）⇒ `… tool=probe.nopath risk=L0 … rules_hit=[] in_allowlist_scope=true …`。**修后**同两发 ⇒ `=unknown`；**AC#2 对照两发一字不变**：真在范围 ⇒ `=true`；真越界（R2 族、人工拒）⇒ `=false`。三株 `-overlay` 突变逐支命中对应子用例（`scopeUnjudged→scopeIn` ⇒ 两枚未判定红；摘 `!InAllowlist` 腿 ⇒ 越界枚红；`scopeIn→scopeOut` ⇒ 范围内枚红），源件 sha256 前后同为 `e07b6692b8f287c3954c44083000fe40059e228f1ef406bf5913c490d7de9d94`（变异全在 `$TMPDIR`，源件零写入）。**AC#3＝不属契约面、未停手**：`in_allowlist_scope` 在 `SPEC-02 §3` 的 `tool_call` DDL、`PLAN.md`、`docs/specs/**` **零命中**（文献面只剩 evidence/票面），"逐字渲染"许诺只覆盖卡片对 `Decision.RulesHit/Reason` 的渲染、不是这行日志；行由本包 `book()` 经 `Options.Logf` 拼。**AC#4 分母**（读数在 `probes/173/r1/logs/`）：第二枚（审批超时 class）判**不改**——`PLAN.md:2790` 的 D37 表把"审批拒绝/超时"定义在同一枚 `user_rejected` 上，分两枚 class＝改 D37 枚举（契约变更，须人批准），而"回执带超时一味"今天已成立且三处可读（回执文案含"超时"＋审计行 `kind=approval-timeout decision=timeout`＋行 `decision` 列 `timeout`），新用例 `Test173ApprovalTimeoutKeepsTheTimeoutFlavor` 把它钉成判据；另逐枚点名（只登记不修）：未过 C19 的调用照印 `risk=L0`＋`rules_hit=[]`（现量两发：`tool=fs.nope` 与 `fs.read`＋坏参数，`bridge_test.go:666` 已把未知工具行钉成 L0）、`risk=` 一栏混"声明下界"（pre-C19 拒绝，`bridge.go:377`）与"评估级"两来源、`grant_id=0` 在没查授权的路上照印、C25 关行 `was_open=false` 时的 `dropped=0/close_err=<nil>`、L1 窗口超时书 `decision=allow`（SPEC-06 §2 有意赋义）；**票 147 那枚 offset 已在册且已按 ⓐ 修好**（`cmd/wisp/slo_windows.go:571-618/737` 现量：`offsetUnknown=-1`＋"带 unknown 的腿不印数字"——本票 `unknown` 的仓内先例）。**门禁**：`GOFLAGS= go build ./...` rc=0；`go test -count=1 ./internal/tools/...` 改前（overlay 回 HEAD `bridge.go`＋摘本票新用例）rc=0／改后 rc=0；`gofumpt -l internal/tools` 空；`sh scripts/d22scan.sh` rc=0（runtests PASS=35 FAIL=0）；`scripts/check-path-length-budget.sh` rc=0。**自报**：`probes/173/r1/logs/` 计 7 件（`00-head` 为派单硬令件、`01`–`06` 每件自落 `rc=N`），比"≤6 件"预算多 1；`05` 里含 overlay 注入的仓外普查探针（树内零 `.go`）。

- [2026-10-13 22:5x] **验收腿 `173-v1` 裁决（非实现者）**（表 `docs/evidence/s1/ticket173-173-r1-accept-v1.md`，读数 `probes/173/v1/logs/` 6 件；源件 sha256 起手/收手同 `e07b6692…d94`）：**AC#1–#5 建议全翻**；AC#4 的“不改”裁重复核成立（`PLAN.md:2790` 同枚 `user_rejected`＋D37“必须是枚举”，拆两枚＝契约变更须人批准），“超时一味”抽一处现跑＝`Test173ApprovalTimeoutKeepsTheTimeoutFlavor` 修后绿（回执含“超时”＋审计行 `kind=approval-timeout decision=timeout`）。**自跑**：overlay 回 `ec2f0864^` ⇒ 两枚未判定子用例红、逐字 `in_allowlist_scope=true`，修后同两发 `=unknown`；对照两发（真在范围/真越界）修前修后归一化后 diff 空（一字不变）；**自造四株 overlay 突变**（`unknown`→`"false"`／`scopeOut`→`"true"`／守卫收窄／`scopeIn`→`"yes"`）逐株红在对应子用例，四个子用例全覆盖；全包 `ok`。**消费面**：生产面唯一读者＝写方 `bridge.go`，无解析器把 `unknown` 按 `true/false` 二分消费。**AC#3 自跑零命中**：`in_allowlist_scope` 在 `PLAN.md`＋`docs/specs/**` 0 命中，`tool_call` DDL 无此列，“verbatim” 只在代码注释里 ⇒ 非契约面成立。**发现（非阻断，已具名）**：`internal/tools/fs_edit_ac5_gate_r4_test.go:350-356` 旧注释仍写 `in_allowlist_scope=true`（同夹具实测 `=unknown`，已腐烂，建议一行修或登记）；r1 的 `01/03` 日志行号比提交版少 1（读数逐字同，只注记）。⛔ 未翻框／⛔ 动 `bridge.go`（除 overlay 瞬态）／⛔ push。
