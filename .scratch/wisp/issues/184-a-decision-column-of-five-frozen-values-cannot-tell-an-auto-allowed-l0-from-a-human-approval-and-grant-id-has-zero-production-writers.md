# 184 — **取证列说不出"没问过人"**：`tool_call.decision` 只有五枚冻结值（`allow`／`allow_session_grant`／`reject`／`timeout`／`batch_aggregated`），"档位即放行"与"人来批准"在审计面上是**同一个字**；唯一能分别的 `grant_id` 在生产里**逐枚调用都写 `nil`**；而 `riskColumn` 把 `""` 压成 `"L0"`，连"未分级被拒"那行也读成 `risk=L0`

- Status: **ready-for-agent（只读普查＋裁定，零产码）**。⚠ **本票不是缺陷修复票**：它的结论有半条落在**契约面**（`decision` 的枚举域在 `docs/specs/SPEC-02-data-storage.md:79` 与 `docs/PLAN.md:2703` 两处冻结，且被 `internal/memory/models.go:136-142` 机器强制），**加一枚值＝人工批准**（AGENTS §0.2/§1.1）。⇒ 普查程**只量、只报、不改**；要不要摆给 owner 由我裁。
- 来源：`179-v1`（表 `docs/evidence/s1/179-declared-l0-refused-accept-v1.md` 33207 字节，次生面那一格）裁"票 179 让 loop 自己给声明 L0 记 `allow`，**不构成对 A13 分工的回归、也不报 177／175**"，并写明：**"'宿主替用户批准'这种误读的风险是真的，根在别处两格"**——那两格就是本票的主体。台账 `A363`。
- 关联：**票 179 AC#9/AC#10**（同一枚记账面的两枚残余，已挂在票 179 上，不在本票）· **票 173**（同族：审计行说出没发生的事。⚠ **现量核对＝不同物**：`grep -cniE "riskColumn|DecisionColumn|词汇" .scratch/wisp/issues/173-*.md` ⇒ **0**，那枚票管的是 `allowlist_scope` 那一格，不管 `decision` 词汇）· **票 40**（`security-privacy-cost` 那页：它的 `:47`／`:48` 两格正是"导出/删除/Grants view"，本票的发现是它开工前必须知道的前提）· **D4**（`PLAN.md:122`，L0 逐字"直接执行，不打断"）· **D45**（`PLAN.md:2137`，批量授权与授权梯度＝`grant_id` 那根线）· **C18**（审批路由与 `correlation_id`）。

## 这是什么（人话）

产品对 owner 的承诺写在方案里那一行字上：**"它到底对我的机器做了什么"必须可查**（`docs/specs/SPEC-02-data-storage.md:70`、`docs/PLAN.md:2703` 两处逐字）。

今天这张取证表能说出"哪枚工具、什么层级、结果是什么"，但**说不出"有没有人批准过"**：

- 最低风险档（L0）的工具——读文件、列目录——宿主**不问任何人**就直接执行，可它在审计里记的那个字是 `allow`，与"人在卡上点了允许"用的是**同一个 `allow`**；
- 唯一能分别"依据哪一份授权放的行"那一列是 `grant_id`，而生产代码里**每一发调用都把它写成 `nil`**（只有一枚 `*_test.go` 给过值）；
- 导出面（`wisp` 的隐私导出）把那枚 `allow` **原样拼进给用户看的那行字**里 ⇒ 用户看到的是"允许"，他无法知道那是**机器自己允许**的。

⚠ **这不是安全洞**（放行本身是 D4 已定案的正确行为，票 179 刚把它修对），它是**取证面的诚实度**：那行记录会让人以为"主人点过头"。三者都在下面有尺。

## 现量（锚 `58124185`，09-28 12:1x／12:2x 编排者本程现跑；`git status --porcelain -- internal/ cmd/` 起手为空）

1. **五枚值，逐字**（尺：`sed -n '26,38p' internal/agent/journal.go`）：
   ```go
   // Decision values (tool_call.decision). Ticket 21 replaces the pass-through
   // decision with real gate outcomes; the vocabulary is already frozen here.
   const ( DecisionAllow = "allow"; DecisionAllowGrant = "allow_session_grant";
           DecisionReject = "reject"; DecisionTimeout = "timeout"; DecisionBatchAggr = "batch_aggregated" )
   ```
   ⇒ 注释自己写着 "the vocabulary is **already frozen** here"。没有任何一枚值的意思是"按档位自动放行、没问过人"。
2. **枚举域被机器强制，不是文档愿望**（尺：`grep -n "toolCallDecisions" -A 8 internal/memory/models.go`）：`models.go:136-142` 是一张**闭合 map**，`:173` 在写入前校验，`internal/memory/dao_toolcall.go:46-48` 在 `DecideToolCall` 入口再校验一次 ⇒ **写一枚新值今天会直接被拒**（`memory: invalid tool_call.decision %q`）。⇒ 加值＝同时动 `models.go` ＋ `SPEC-02:79` ＋ `PLAN.md:2703` ＝**人工批准面**，不是代码 tidy-up。
3. **两处冻结出处逐字**（尺：`sed -n '70,80p' docs/specs/SPEC-02-data-storage.md`；`grep -n "取证记录" docs/PLAN.md` ⇒ **`:2703`**）：
   `risk_level TEXT NOT NULL, -- 'L0'|'L1'|'L2'` ／ `decision TEXT, -- 'allow'|'allow_session_grant'|'reject'|'timeout'|'batch_aggregated'`，且那行注释就是 owner 那句「它到底对我的机器做了什么」必须可查。
4. **`grant_id` 生产零填充**（尺：`grep -rn "GrantID:" --include=*.go internal/ cmd/` ⇒ **0 命中**；尺：`grep -rn "DecideToolCall(" --include=*.go internal/ cmd/` ⇒ 生产**只有一枚**调用方 `internal/agent/journal.go:101`，逐字 `_ = j.DecideToolCall(ctx, rowID, decision, nil)`）。
   ⇒ 两笔要分开记：**①授权关联列永远是 NULL**（表 `approval_grant` 本身在 `internal/memory/schema.go:89` 是存在的，尺见第 3 组）；**②那一枚调用的错误被 `_ =` 丢掉**（校验一旦拒绝，取证行静默不写＝另一族"不表现为报错"）。
5. **导出面原样拼给用户看**（尺：`sed -n '96,104p' internal/memory/privacy.go`）：`detail := fmt.Sprintf("%s · %s", r.RiskLevel, r.Tool)`，随后 `if r.Decision != "" { detail += " · " + r.Decision }`。⇒ 用户读到的是 `L0 · fs.read · allow`，**没有任何一层说明那是档位自动放行**。⚠ `privacy.go` 在 **`internal/memory/`** 不在 `internal/agent/`（`179-v1` 引它时省了目录名，我核对过：尺 `git ls-files | grep -i privacy`）。
6. **哪些地方在记 `allow`**（尺：`grep -rn "DecisionColumn = agent.DecisionAllow" internal/tools/bridge.go` ⇒ **3 处**：`:374`／`:386`／`:400`；尺：`grep -n "j.decide(ctx, rowID, DecisionAllow)" internal/agent/loop.go` ⇒ **2 处**：`:794`（票 179 新加的声明 L0）／`:803`（未分级＋直通开关开着））⇒ **五枚记账点、一枚词汇值**。哪几枚"问过人"、哪几枚"没问"，本票**不猜**，交给下面 AC#1 逐枚读。
7. **`riskColumn` 的压平**（尺：`grep -n "func riskColumn" -A 10 internal/agent/journal.go`）：`case memory.RiskL1, memory.RiskL2: return risk; default: return memory.RiskL0`，生产调用方两枚（`journal.go:80`、`loop.go:611`）。
   ⇒ **这一支我裁定它不是"写错"**：`SPEC-02:76` 把 `risk_level` 定成 `NOT NULL` 且域只有 `L0|L1|L2`，**"未分级"在这张表里没有合法落点**。真正的后果是一句读数规则：**`risk=L0` 不能读成"声明为 L0"**，它同时容纳 `""`；今天能把两者分开只剩 `decision` 一列（`reject` vs `allow`）＋`error_class`。
8. **直通开关在生产仍然零赋值点**（尺：`grep -rn "PassThroughUnclassifiedRisk" --include=*.go .` ⇒ 生产无赋值，`true` 只在 `internal/agent/harness_test.go:116`）⇒ 第 6 把尺里 `loop.go:803` 那一枚 `allow` **在生产形状上今天到不了**；AC#3 要把这句证死而不是引用我。

## 为什么值得做（不做会怎样）

票 40（安全／隐私／成本那页）开工第一件事就是**把这张表摊给人看**。如果 `decision` 说不出"有没有人点头"，那页会**主动制造**一种错误的信任："它每次动手都问过我"。owner 在意的是隐私取证，而取证面最坏的形状不是"没记录"，是**记录读起来比事实好看**。⇒ 现在花一枚只读普查把"哪些行不可区分、要动哪几处、有没有不动契约的那半条路"量清，比到那页开工时再撞上便宜得多。

## AC（每格都要给现量；先测→再报→再提交，**零产码**）

- [ ] **AC#1 五枚记账点逐枚读，产出一张"问没问过"表**：`internal/tools/bridge.go:374`／`:386`／`:400`＋`internal/agent/loop.go:794`／`:803`（号按 HEAD 重锚，尺＝上面第 6 把）逐枚回答三档：**①这一枚当时有没有问过人；②它在真机形状上可达吗（给调用链）；③它与另一枚共用 `allow` 后，取证面上还有没有别的列能分开它们**。⚠ 不许只报"它记了 allow"；不许把 `:400` 那枚的语义靠猜（它前后各有一支 `Reject`，见 `:392`／`:405`／`:408`）。
- [ ] **AC#2 `decision` 的**读者名册**（不是文件位置）**：逐枚列出今天**消费**这一列的地方，给 `file:line` ＋"它把 `allow` 解释成什么"。已知起点：`internal/memory/privacy.go:99-101`（导出面拼字符串）。⚠ **同名不同物的陷阱必须避开**：`tools.Decision` 是**路由用的结构体**，与 `tool_call.decision` **这枚列**不是一回事——尺要写成 `grep -rn "\.Decision\b\|DecisionColumn\|decision=?" --include=*.go internal/ cmd/` 再**逐枚判它读的是列还是类型**，`grep -rl "decision"` 那种文件级名册**不算答案**（全仓文件级命中一大片，其中多数是类型名）。⚠ 也要报**面板／`wisp` 命令**那几面有没有读者：判"没人读"必须给一条**取数命令**，不能给"它在别的目录"。
- [ ] **AC#3 证死或推翻"`L0`＋`allow` 组合的唯一来源"**：我的现读判断（**待你推翻或坐实**）＝第 8 把尺说 `loop.go:803` 在产码里到不了，因此生产上 `risk=L0 且 decision=allow` **只可能**来自"档位即放行"（`loop.go:794`、`bridge.go:374`）。⇒ 用**调用链＋现有用例名册**证它；若证伪（存在一条生产路让未分级也记 `L0+allow`），那 AC#5 那条显示层推断就不成立，要具名报回。
- [ ] **AC#4 量存量：加一枚 `decision` 值要动哪几处**（**只量不改**，这是给 owner 那张单的代价表）：至少核 `internal/memory/models.go:136-142`、`internal/memory/schema.go`（DDL 注释）、`docs/specs/SPEC-02-data-storage.md:79`、`docs/PLAN.md:2703`、迁移目录（枚数现跑：`ls` 那个目录 + 计数命令要贴）、以及**任何钉枚数／名册的判据**（尺要现跑：`grep -rln "batch_aggregated\|allow_session_grant" --include=*_test.go .` ⇒ 逐枚读，判它钉的是"五枚"还是"值本身"）。⚠ 迁移目录那格有个已知雷：**版本号当幂等键会撞号**（本项目踩过），所以"加一列/改一枚枚举"要顺带报"新库与老库分别会缺什么"。
- [ ] **AC#5 先找不动契约的那半条路**（我提三形，**逐形给可达性与最坏后果**，不许照抄我的判断）：
  **(a)** 导出面文案：`privacy.go:99-101` 那句 `detail` 里，在 `risk_level==L0 且 decision==allow` 时写成人话（"自动放行（按档位，未询问）"）——⚠ **它依赖 AC#3 成立**，且**列里的字一字不动**；
  **(b)** 把 `grant_id` 真填上（`approval_grant` 表在 `internal/memory/schema.go:89` 存在，D45 已有规格）——要报"D45 的授权记录今天有没有生产者"（尺：`grep -rn "approval_grant" --include=*.go internal/ | grep -v _test.go`）；
  **(c)** 接受歧义并在 `docs/reports` 里明写"取证面不区分自动放行与人工批准"。⇒ 三形都不许顺手动 `riskColumn`（第 7 把尺裁定它按现有 schema 是对的）。
- [ ] **AC#6 判"错误被丢掉"那一支要不要单独归口**：`internal/agent/journal.go:101` 的 `_ =`——如果校验今天可能失败（枚举例如写入者给了没登记的值），那**取证行会静默消失**。⇒ 现量：今天有没有任何生产路径会给 `DecideToolCall` 传一枚不在 map 里的值（尺：`grep -rn "DecideToolCall\|j.decide(" --include=*.go internal/ | grep -v _test.go`，逐枚看传入值来源）。判"要单开立票"就**停手报回**，别在本票里顺手改。
- [ ] **AC#7 门禁（只读票也要跑一次全仓那一档）**：`sh scripts/d22scan.sh`（⚠ 现量在册基线：`ban #8 internal/` **431**，取数时刻 09-28 11:2x 由 `179-r2` 钉；你这一跑只要**不涨**，涨了就是你带了写件进来）＋ `bash .scratch/wisp/probes/154/gate-clauses.sh` **比红腿名册不比退码**（今天在册唯一红腿＝`G6neg`＝票 178；尺：`grep "BAD" <输出> | grep -oE "腿=G[0-9a-z]+" | sort -u` ＋ `comm -3`）＋ `git status --porcelain -- internal/ cmd/` 终态为空。**不许**跑 `probes/161/r6/flip-declaration.sh`（跑一次就脏跟踪日志）。
  > ⚠ **09-28 12:5x 编排者追加的排程约束（今天有效）**：本格**放在最后一步跑**，且跑之前先看 `git status --porcelain -- internal/ cmd/`——同机此刻有程在 `internal/risk/**` 上做变异，**吃脏件换来的"红"或"绿"都不算读数**；跑不动就记〔未取到·争用〕并具名说明（本编队的老规矩：响亮拒绝胜过假绿）。另 `ban #8 internal/` 若从 431 变 **432**，那是 `183-r1` 新判据件的账，**不是回退**。
- [ ] **AC#8 契约轴禁令**：`docs/PLAN.md`、`docs/specs/**`、`internal/memory/schema.go`／`models.go`、`thresholds.go`／golden／审批超时常量／`allowlist.txt` **一字节不许动**；**不许新增 `decision` 值**；不许改 `riskColumn`。⇒ 判"必须动这些才答得清"＝**停手上报**（那是要我批准，不是你自己决定）。
- [ ] **AC#9 交件表**落 `docs/evidence/s1/184-audit-vocabulary-census-c1.md`，含：本程**没**测什么、被拒／没成功的调用、有没有跑过删除命令、工具调用枚数 vs 硬顶、伪授权两栏、凭据值零抄录。**AC 框由编排者翻，你只交读数。**

## 写面（超出即越权）

`.scratch/wisp/probes/184/c1/**`（台件与读数，`logdir` 取脚本自身目录；⚠ **不许用 `runtime.Caller` 在 `-overlay` 下取日志目录**，它会报虚拟路径）· `docs/evidence/s1/184-audit-vocabulary-census-c1.md` · 本票面 Progress log（**只追加**）。
**禁改**：`internal/**`、`cmd/**`、`docs/PLAN.md`、`docs/specs/**`、`frontend/**`、`design/**`（别家归属，不读不写不引）· 别人的票面与证据件 · `probes/**` 既有台件（只读）。
⚠ 现场躺着 `design/**` 未提交删除、`.gitignore`／`probes/152/my152.py`／`probes/161/r6/logs/flip-*`（` M`）／`docs/evidence/s1/152-*.md` ⇒ **不提交、不还原、不补完、不评论**。
Git：只 commit 不 push；`git add -A`／`git add .`／`commit -a` 一律禁；一步式 `git commit -q -F - -- <显式路径>`；禁 `--amend`／`reset`／`rebase`／`stash`／`checkout <分支>`／`switch`／`merge`／`worktree`／`clean`；临时件只建不删，还原用 `git cat-file blob HEAD:<路径> > <路径>`。**硬顶 40 枚工具调用，到顶即停手回禀并自报枚数。**

## 本票**不**解决

- 不裁"该不该给 L0 上一张卡"——D4 逐字"直接执行，不打断"已定案，本票只管**记录能不能说实话**。
- 不翻任何别人的勾（票 179 AC#1/AC#8 维持不勾；票 177 AC#3 条件、票 175 AC#5、票 176 AC#3-5 都在等票 183 的端到端）。
- 不动 `riskColumn`、不写"未分级"的合法落点（那要先动 schema，属 AC#4 的代价表面）。
- 不替 owner 拍板；普查交完后**由我**决定摆不摆 `Q`，并给一行零术语的"人话后果"。

## Progress log

- 09-28 12:2x 编排者立票：上面 8 把尺本程现跑（锚 `58124185`）。其中两枚是**对我自己上一轮的更正**：① `179-v1` 引 `privacy.go:100-101` 时省了目录，真身在 `internal/memory/`（尺第 5 把）；② 我本来要把第 7 把尺写成"`riskColumn` 写错了"，读了 `SPEC-02:76` 之后改成"**schema 逼出来的，不是缺陷**"——**结论方向没变（取证面分不清），但成立的理由换了一次**，这条按记忆里的老规矩显式登记。票 173 的"同族"核对为**不同物**（尺在关联那行）。未派。
