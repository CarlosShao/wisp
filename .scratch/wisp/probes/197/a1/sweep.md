# 197-a1 全仓普查：「要求合并/折叠分片」这一族说法 vs `Overflow TRUNCATES and never merges`

腿：`197-a1`（只读普查）。锚与逐尺见同目录 `00-anchor.md`（HEAD `052b393f`，`is_ancestor rc=0`）与 `rc-log.md`。
交件时刻 `2026-10-08 11:2x +08`。⛔ 本程未改任何产码/测试/票面/台账，未勾任何框，未 push。

## 0. 射程与尺（自报）

- 跑面：**只用了** `grep -rn` / `git grep -n HEAD -- frontend`（对象层）/ `git show <ref>:<file>` / `sed -n` / `awk` / `stat` / `date` / `git log -S`。
  **⛔ 未跑 `go build`/`go vet`/`go test`**（编译面由 `255-r1`、`35-r7` 在飞）；**也没跑 `go env`/`go list`** —— 本程一次都没用。
  ⇒ 第④问"什么形状能让它们红"全部是**读码判定**（夹具的界/枚数 + 断言极性），不是突变实测；每条都标了依据行号。
- `frontend/**` 全走对象层：`git grep -n -i -E ... HEAD -- frontend/`（66 行命中，`rc=0`）、`git show 5821e24f:internal/panel/pump.go`（341 行）。
- 长行处理：票面命中长行（如票 35 `:51`、`:392`，票 226 `:32`，票 9 `:65`）**逐枚 `Read` 整行**，⛔ 没有用过 `cut -c`。
- 变体试过的（按"这句话在说什么"，⛔ 不只搜符号名/数量词）：
  `merges`/`merge`/`MERGES`/`merged into`/`merge/backpressure`/`backpressure/merge`/`folds`/`fold`/`folded`/`EqualFold`/`foldPath`/
  `combine chunks`/`coalesce`/`dedup`/`并成`/`并入`/`被并`/`归并`/`聚合`/`汇总`/`塌缩`/`collapse`/`concat`/`join`/
  `合并`/`折叠`/`拼接`/`合成一行`/`合成为一行` + 域词共现（`chunk`/`row`/`delta`/`stream`/`overflow`/`backpressure`/`truncat`/`快照`/`分片`/`截断`/`背压`/`过载`/`溢出`/`一行`/`每行`/`流式`）。

### 名册枚数的尺写法（两把尺都说清）

| 尺 | 枚数 | 说明 |
|---|---|---|
| **A 组＝同一族（面板快照分片/事件推这一条链）的命中行数** | **33 行** | 尺＝逐枚人工归类下面 §1 的 A1/A1b/A2 三张表的行数；`grep -n -i merg` 只命中票 35 的 **7 行**（`rc=0`，`01-readings` R7），词面尺会少报 26 枚 |
| **A 组去重后的"处"** | **26 处** | 差 7 ＝同一条规则被多处复述（票 35 `:14` 与 `:27` 是同一段许诺的两行、票 35 `:392`/`:406` 是 `:51` 那次裁的两次回写、票 197 `:10`/`:121` 是同一格改前改后）⇒ **本项目"同一枚数两种尺写法差 1"栽过，这里差 7，必须带尺引用** |
| **B 组＝别的链上的同名规则（不归口、只登记）** | **14 族 / 逐族给代表行号** | 见 §1-B |

---

## 1. 名册（逐枚 `<file>:<line>` ＋ 那半句逐字 ＋ 归类）

### A1 要求合并（面板快照分片 / 面板事件推这一条链）——**9 处**

| # | `<file>:<line>` | 那半句（逐字） | 归类 |
|---|---|---|---|
| A1-1 | `docs/PLAN.md:2848` | 「- LLM stream → UI：有界 channel，满则**合并增量**（不丢内容）。」 | **要求合并**（D38(d) 冻结决策原文，最高权威级） |
| A1-2 | `docs/PLAN.md:2971` | 「③ correlationId 路由 + 事件推送的**背压与合并**（D38(d)）」 | **要求合并**（C17 契约深化四项之一＝冻结契约面） |
| A1-3 | `docs/specs/SPEC-01-architecture.md:133` | 「- LLM 流→UI：有界 channel，满则**合并增量**（不丢内容）。」 | **要求合并**（spec 逐字复述 A1-1） |
| A1-4 | `.scratch/wisp/issues/35-panel-bridge-c17.md:14` | 「…approval requests, ball state, cost ticks) **with bounded-queue merge**, and `panel.resync` full-」 | **要求合并**（票 35 §What to build） |
| A1-5 | `.scratch/wisp/issues/35-panel-bridge-c17.md:27` | 「- Event push: bounded queue per panel; **overflow MERGES increments** (never drops content);」 | **要求合并**（票 35 §Key constraints） |
| A1-6 | `.scratch/wisp/issues/35-panel-bridge-c17.md:49` | 「- [ ] Backpressure: flood events under blocked consumer → **merges**, no unbounded memory (heap cap asserted), **content integrity kept**.」 | **要求合并**（**未勾 AC 框**，尺＝票 35 未勾 6 枚含它） |
| A1-7 | `.scratch/wisp/issues/36-result-history-panel.md:44` | 「- [ ] Long-task stream: 200-chunk golden replay **renders with merge** (no jank), token counter accurate.」 | **要求合并**（**未勾 AC 框**；票 36 未勾 **6** 枚／已勾 **0** 枚，尺＝`grep -c '^- \[ \]'` rc=0 得 6、`grep -c '^- \[x\]'` rc=1 得 0。**这一枚是新发现的、不在编排者转述里**） |
| A1-8 | `internal/panel/doc.go:10` | 「// routing + push **backpressure/merge** (D38d)」 | **要求合并**（产码包职责清单；`35-panel-snapshot-pump-r1.md:46` 引它时写的是 `doc.go:9`，**行号已漂 +1**） |
| A1-9 | `cmd/wisp/panel_pump.go:397` | 「// for one more word of a reply, and **the merge/backpressure rule** that decides how / :398 // often a panel is pushed **belongs to the transport (D38d, ticket 35 AC#5)**, which / :399 // does not exist yet.」 | **把合并当既定规则引用**（产码注释；⚠ 它自己限定射程＝"**多久推一次**"，不是"两行并一行"。见 §2） |

⚠ **行号撞车预警**：`internal/panel/pump.go:397`（禁止合并）与 `cmd/wisp/panel_pump.go:397`（引用合并）是**两枚不同文件的同一行号**。编排者转述只指了前者；后者是本程新量到的。

### A1b 历史记录 / 过期读数（说"今天会合并"，字面已不成立）——**13 处**

| # | `<file>:<line>` | 那半句（逐字，节选） | 归类 |
|---|---|---|---|
| H-1 | `.scratch/wisp/issues/35-panel-bridge-c17.md:94` | 「票面 §"What to build" 里那句 "Go→frontend event push … **with bounded-queue merge**" **今天兑现枚数＝0**」 | 历史记录（且把 A1-4 的许诺原样再抄一遍） |
| H-2 | `.scratch/wisp/issues/35-panel-bridge-c17.md:392` | 「三枚既有尺 `pump_test.go:165 …`／`subagent_stream_197_test.go:209`/`:275` 会因改回 merge 而红＝**唯一有牙的尺，咬的方向与 AC 相反**」 | 历史记录（本腿交件登记） |
| H-3 | `.scratch/wisp/issues/35-panel-bridge-c17.md:406` | 「`:49` 额外要先裁"merge vs truncate"那处与票 197 的正面冲突」 | 历史记录 |
| H-4 | `.scratch/wisp/issues/35-panel-bridge-c17.md:51` | 「〔**10-08 11:3x 编排者就地裁冲突** … 这一格三句里 **"merges" 与 "content integrity kept" 两句已被票 197 leg B 的重裁取代** …〕」 | **裁决登记**（编排者的射程收窄；⛔ 框仍未勾。**内容与行号已现量核对，见 §7 第 3 条两处更正**） |
| H-5 | `.scratch/wisp/issues/197-the-panel-…-entity.md:10` | 「**键上限 32、溢出是"合并而不是丢"**」（行号注 `pump.go:301`/`:339`/`:275`/`:287-290`） | 历史记录（改前普查读数，**行号全部过期**） |
| H-6 | `.scratch/wisp/issues/197-…entity.md:121` | 「`:509` 注释逐字 `no key is ever folded into another`，判定体在 `:514`/`:521`」 | 历史记录（**改后现读，行号也已过期**：现量锚在 `pump.go:573`，见 §7 第 4 条） |
| H-7 | `docs/reports/HANDOVER.md:1880` | 「另 `StreamLog` 今天**键上限 32、溢出"合并而不是丢"**（`pump.go:275/:287-290`）对子代理是错形状，197 要重裁」 | 历史记录（开票时读数） |
| H-8 | `docs/reports/pending-and-issues.md:8567` | 「…溢出是"合并而不是丢"…两条流被并成一条＝界面在说谎。⇒ 197 要**重裁溢出语义**」 | 历史记录（台账 A 条，开票） |
| H-9 | `docs/evidence/s1/35-panel-snapshot-pump-r1.md:255` | 「`5821e24` \| `internal/panel/pump.go` \| 341（新）\| … `StreamLog`（**有界合并**）」 | 历史记录（合并语义**被引入**那枚 commit 的登记） |
| H-10 | `docs/evidence/s1/35-panel-snapshot-pump-r1.md:490`＋`:652` | 「`TestTheStreamLogMergesInsteadOfDropping`」／「用 50 枚 key 量了折叠」 | 历史记录（**该用例名已不存在**，197-r2 更名；本程按 `^func Test` 名册尺现量＝仓里只剩 `…TruncatesInsteadOfMerging`） |
| H-11 | `docs/evidence/s1/154-host-id-never-closed-r1.md:343` | 「有界＝`const DefaultStreamKeys = 32`（`internal/panel/pump.go:279`）＋ **`mergeOverflowLocked`（`:344`）**」 | **过期读数（最危险的一枚）**：该函数现量**全仓 `.go` 0 命中**（R1：`grep -rn mergeOverflow --include=*.go internal cmd tools` **rc=1、0 行**），行号也过期；⚠ 后续程若引这行会以为 merge 还活着 |
| H-12 | `docs/evidence/s1/197-subagent-stream-r2.md:26` | 「溢出语义"合并不丢"约 `:287-290` ｜ …**实现在 `:366-380` `mergeOverflowLocked`** ｜ 该函数已从仓里消失」 | 历史记录（**自陈已废**，写法合格） |
| H-13 | `.scratch/wisp/dispatches/2026-09-28-184x-impl-197r1-entity-and-197r2-stream.md:79` | 「`DefaultStreamKeys = 32`（约 `:275`），**溢出语义＝"合并而不是丢"（约 `:287-290`）**」 | 历史记录（派单＝重裁**之前**的射程描述；`:85` 同票要求"乙＝绝不把两条流并成一条"；另一枚派单 `…210x…:55` 已写对"溢出语义今天已是显式截断"） |

### A2 禁止合并（同一条链的现行规则）——**11 处**

| # | `<file>:<line>` | 那半句（逐字，节选） |
|---|---|---|
| N-1 | `internal/panel/pump.go:397`（段 `:397-409`） | 「**Overflow TRUNCATES and never merges** (ticket 197 leg B re-cut this rule, which previously folded the two oldest chunks into one row)… So past the key bound every key keeps its own chunk… "nobody is showing this one" is a different sentence from "somebody merged this one".」 |
| N-2 | `internal/panel/pump.go:468` | 「now that **overflow truncates instead of merging**, it bounds」 |
| N-3 | `internal/panel/pump.go:572-573` | 「enforceBoundLocked is ticket 197's overflow rule, and its one invariant is that **no key is ever folded into another**」 |
| N-4 | `internal/panel/pump.go:655-656` | 「the host can say WHICH ones and say they were **never merged into a row that is still on screen**」 |
| N-5 | `internal/panel/pump.go:683-684` | 「and **never folded into a neighbour's row instead** (197-r2 §1 乙)」 |
| N-6 | `internal/panel/subagent_roster_197.go:7-8` | 「made its stream honest (**StreamLog truncates instead of merging**, and names what it drops)」 |
| N-7 | `internal/panel/subagent_roster_197.go:165-166` | 「"somebody **merged it into a row that is still here**" (197-r2 §1 乙)」 |
| N-8 | `internal/tools/subagent_197.go:295-297` | 「The child's text **must NOT ride the parent's console sink**: that is how two tasks' output **gets merged into one visible stream**.」 |
| N-9 | `internal/tools/subagent_197.go:509-511` | 「nothing here writes to stdout, because **two tasks' text on one console is a merge**」 |
| N-10 | `frontend/src/components/result-stream.tsx:15`（对象层 `git grep HEAD`） | 「a finished chunk's text renders **VERBATIM as one text node**」＝每枚分片自己一行，不并 |
| N-11 | `.scratch/wisp/issues/197-…entity.md:32` ＋ `:39`（AC#3） | 「⚠ **32 枚键上限与"溢出合并而不是丢"这一形对子代理是错的**——两个子代理的输出被合并成一条＝界面在说谎。⇒ 本票要**重裁溢出语义**」／「- [ ] **AC#3 每子代理一条流、不许合并**」（现量票 197：未勾 **5**、已勾 **2**，尺＝`grep -c`） |

---

## 2. 逐处判冲突（是不是 `pump.go:397` 否决的**同一种**合并）

判据我用的是**同一条物理链**：`StreamLog` 的**行身份**（一枚关联键＝一行，两枚键的文字不许进同一行）。
A1-1…A1-9 逐枚判：

- **A1-5（票 35 `:27`）／A1-6（票 35 `:49`）＝真冲突，同链**。「overflow **MERGES** increments」「blocked consumer → **merges**」＋「content integrity kept」两句话的正是"过界之后把内容并起来"这一支，与 N-1 逐字相反。编排者已裁（H-4）。
- **A1-1（`PLAN.md:2848`）／A1-2（`PLAN.md:2971`）／A1-3（`SPEC-01:133`）＝同一句话的三个权威层级拷贝，冲突面比票 35 更大**：
  - 读法甲（字面）：「满则合并增量」＝票 35 `:27` 的祖先，那么被否决的不止票面，还有 **D38(d) 决策原文与 C17 契约第③项**。⇒ 收窄它们＝动 `D1–D47`/`C1–C32`＝**人工批准**。
  - 读法乙（**盘上有一条产码依据支持它**）：`cmd/wisp/panel_pump.go:396-399` 逐字把这条规则限定成「decides **how often a panel is pushed** belongs to **the transport**（D38d, ticket 35 AC#5），which does not exist yet」＝**推频**（少刷几次、把增量攒起来再推），不是"两行并成一行"。攒起来的帧仍是**每键一行**，行身份不变 ⇒ **与 N-1 不矛盾**。
  - **本程只把两读并列、不裁**（未定义即停）；这条读法歧义是 §6 三形的分水岭。⚠ 但"读法乙"要成立，`PLAN.md:2848` 那句「合并增量」今天**没有**任何文字说明它只讲推频——那半句是 `panel_pump.go` 的注释补的，不是 PLAN 补的。
- **A1-4（票 35 `:14`）＝同链**：与 `:27` 同一许诺的两处抄写（`bounded-queue merge` 在 §What to build 首段），同样落在甲/乙歧义里。
- **A1-7（票 36 `:44`）＝本程新发现、同族、**射程有歧义****。「200-chunk golden replay **renders with merge** (no jank)」——`(no jank)` 指向**渲染节流**（乙读法），`200-chunk … with merge` 指向**分片合并**（甲读法）。票 36 `Blocked by: 30, 35`、AC **6 枚全未勾**，是一张 `ready-for-agent` 的活票 ⇒ **它就是"后面的写腿去修好一个已被否决的行为"的现实通道**：一枚照 `:44` 字面去"把 200 枚 chunk 合成一行渲染"的腿，会直接撞 N-1 与 §4 那 7 枚有牙尺。**这一枚不在编排者转述的范围内，是这次扫出来的第二处同链冲突。**
- **A1-8（`doc.go:10`）＝同链**（包职责清单里写"push backpressure/merge"）。它是对 A1-1 的引用，不是独立裁处；⚠ 它在 `internal/` 里 ⇒ **d22scan ban #8 的 `internal/` 档覆盖它**（`goOnly:true`，注释豁免；`tools/d22scan/main.go:579-586` 现量），改它不会点 emoji 红，但会改包职责文字。
- **A1-9（`panel_pump.go:397`）＝按乙读法不冲突、按甲读法冲突**（同上）。它至少是一枚**外向影响**：产码注释把"合并"当成还没写的那根管子的既定策略。

**B 组：别的链上的同名规则，⛔ 不许与上面归口**（逐族代表行号）：

| 族 | 代表 `<file>:<line>` | 那句话在说什么 | 与本冲突的关系 |
|---|---|---|---|
| B1 LLM 用量聚合 | `internal/llm/llm_test.go:52 TestUsageMaxMergeAggregation`；`docs/evidence/s1/09-adversarial-acceptance.md:160-161` | 「Usage 多 chunk **聚合** — PASS（**max-merge** 设计对 Anthropic 累计式安全」 | 合法合并（token 计数取 max，不涉文本行） |
| B2 流式文本累计 | `docs/evidence/s1/09-adversarial-acceptance.md:231` | 「TextDelta+ReasoningDelta **顺序拼接**」 | 合法（同一枚键自己的 delta 累积＝N-1 也允许这一支） |
| B3 配置合并写 | `.scratch/wisp/issues/226-…:33`（AC#2）「**要么合并并报告实际写了哪几枚键**」；`docs/evidence/s1/226-config-write-no-clobber-r1.md:95/:107`；`248-…:333`（`mergeWrite`） | 磁盘 vs 内存快照的键级合并 | 不同链（落盘），合法且**已被别的票裁成要求** |
| B4 上下文压缩 | `internal/agent/compress_test.go:52 TestCompressFoldsOldestKeepsLastThreePreservesIDs`；`compress_trace_test.go:366 TestCompressionTraceDoesNotAlterTheFold` | 历史消息折进摘要 | 不同链（模型侧上下文），合法 |
| B5 settle 判定折叠 | `internal/observe/sampler_settle_gate_136_test.go:325 TestFoldSettlePassOnlyGateRowsVeto`；`.scratch/wisp/issues/136-…:296/:318` | `foldSettlePass` 把门行折成 pass | 不同链（观测判定），合法 |
| B6 路径 fold/EqualFold | `internal/tools/pathshape_portable_test.go:76 TestFoldPathKeepsPosixBackslashesDistinct`；`internal/winsec/…DoesNotFoldABackslashIntoASeparator`；票 `252:11`、票 `113:155`、票 `115:178` | 大小写/分隔符折叠，且多数是**禁止折叠** | 不同链；⚠ 词面最像"fold 规则被钉住"，容易被误当同类 |
| B7 审批批量聚合 D45 | `docs/specs/SPEC-06-security-gatekeeping.md:115-116`「N ≥3 个同质 L1 操作 → **合并为一个确认**…**L2 永不聚合**」；`SPEC-00:70`；`SPEC-10:95`；`PLAN.md:2849-2850` | 确认卡合并 | 不同链（门控），合法 |
| B8 界面"折叠"＝视觉 collapse | `docs/PLAN.md:3483`「推理过程…**默认折叠为一行**」；`SPEC-08:192`；`design/old/index.html:180`；`design/doubao/README.md:65`；`frontend/src/components/ai-native/thinking.tsx:105`、`tool-chips.tsx:112`；`design/old/screens/approval.html:185/655`「这一行**不得折叠**」 | 展开/收起 UI，**内容一枚不少** | ⚠ 中文"折叠一行"与 N-1 的"folded into one row"字面撞车；这是**观感**不是**合并**，归口会造出假冲突（§6 ⓒ 要防的就是这个） |
| B9 附件去重 | `internal/panel/attachments.go:99-101/:268`；`attachments_test.go:244`；`frontend/src/lib/panel.ts:106`（`deduplicated: boolean`）；`composer.tsx:318` | 同字节附件复用同一份产物 | 不同链，合法；⚠ `deduplicated` 是**受双向对账尺盯的 C17 字段名**（`TestComposerContractTypesMatchFrontend`） |
| B10 git refs 合并 | `internal/panel/git.go:139`「listing drops it … merged with」；`:437`「refs/heads walked recursively, **merged with** packed-refs」 | 松散 ref 与 packed-refs 取并 | 不同链，合法 |
| B11 仪器自身 fold | `cmd/wisp/leg_sink_gate_131_test.go:790-795`「merge folds a second declaration of the same name into the first」；`leg_dispatch_gate_133_test.go:278/:1035/:1153`；`internal/panel/inbound_roster_253_test.go:523-526 mergeConsts` | AST 名册把同名声明并成一行 | 不同链（仪器内部），合法；⚠ 这几枚是**会红的尺**，别当产码语义 |
| B12 治理 meta | 台账 `pending-and-issues.md:1839/:3039/:7086/:10268`、`injection-timeline.md:399`、`HANDOVER.md:422`、票 `224:36`、票 `35:384` | 「不许合并成一句/一枚总数/一程」 | 与产码无关 |
| B13 用户可见文案（任务粒度） | `internal/tools/subagent_197.go:274` 「这是硬拒，不是排队——等哪一枚结束了再派，**或者把任务并成一枚子代理。**」 | 建议**把两枚任务并成一枚** | ⚠ 这是**面向用户文案里的"并成"许诺**，但并的是任务不是分片；`grep 并成` 会命中它，别归口 |
| B14 git 合并冲突计数 | `docs/evidence/s1/147-offset-naming-r1-accept-r1.md:658`「合并冲突 **1**」 | 通知分类计数 | 无关 |

---

## 3. 权威文本在哪（票 197 leg B 那次重裁）

**结论：盘上有权威文本，而且不止一处；但"规则本身"没有落进 `docs/`（PLAN/spec）任何一格。** 逐枚：

1. **规则正文（产码注释，逐字）**＝`internal/panel/pump.go:397-409`。关键三行逐字：
   > 「// Overflow TRUNCATES and never merges (ticket 197 leg B re-cut this rule, which
   > // previously folded the two oldest chunks into one row). … A subagent's stream is not interchangeable
   > // text: the owner clicks one subagent to read THAT agent's work, and a row carrying
   > // two agents' sentences would show them something false about provenance - the same
   > // disease D30 names for tool output, arriving through the display side instead of the model side.」
2. **要求重裁的判据（票面）**＝`.scratch/wisp/issues/197-…entity.md:32`（"两个子代理的输出被合并成一条＝界面在说谎"）＋ `:39`（**AC#3 每子代理一条流、不许合并**，现量**未勾**）。
3. **裁处选择（乙）的文字**＝`docs/evidence/s1/197-subagent-stream-r2.md:41-54` §1「选型：**乙**（显式标"已截断" + 每键自留头尾），为什么不选甲」。⚠ 该件 `:5` 自己声明「本件只是**实现者的自述**」⇒ 它单独不构成批准记录。
4. **批准／验收记录（盘上可核的那一枚）**＝`docs/reports/pending-and-issues.md:8766` **A406**（09-28 19:0x）逐字：「## A406（09-28 19:0x）：**收 `197-r2`（溢出改"显式截断"、合并那枚函数已从仓里消失）**」，并带编排者自跑的凭据（`grep -rn "mergeOverflowLocked" internal/panel/`＝0 命中、`go test -count=1 -run '197|Overflow|Stream' ./internal/panel/`＝`ok`）。⇒ **不是"只剩派单转述"这一格，盘上查得到。**
5. **代码凭据**＝commit `94ea50f7`「panel(197-r2): 子代理各占一条流，溢出从"合并着吞"改成"显式截断"」（`git log -S mergeOverflowLocked` 现量，只两枚：引入＝`5821e24f`「StreamLog **有界合并**(溢出折叠不丢内容)」、删除＝`94ea50f7`）。被删函数体逐字（`git show 5821e24f:internal/panel/pump.go`）：「mergeOverflowLocked **folds the oldest chunk into the one behind it** until the log is back inside its bound. The fold keeps BOTH texts」。
6. **⛔ 缺口（具名）**：`docs/PLAN.md:2848`/`:2971` 与 `docs/specs/SPEC-01:133` **一个字都没改**，仍是"满则合并增量"。⇒ 这次重裁**只落到票面＋产码＋证据＋台账，没有回写决策与契约文本**。按 `docs/specs/README.md`「PLAN.md 定案内容最高」的优先级，这是"权威文本没闭合"，不是"没有权威文本"。
7. **leg B 的指认**：`pump.go:397` 说的 "leg B" ＝派单 `.scratch/wisp/dispatches/2026-09-28-184x-…:85` 的「**乙**＝到上限时显式标"已截断"并保留每个键各自的头尾，**绝不把两条流并成一条**」，实现腿＝`197-r2`；`pump.go:684`、`subagent_roster_197.go:166` 两处注释都回指「197-r2 §1 乙」。

---

## 4. 今天有没有仪器钉着这条规则（名册化；⚠ 全部为**读码判定**，本程未跑测试）

`grep -n -F 'mergeOverflow' --include=*.go internal cmd tools` → **rc=1、0 命中**（合并那条路真没了，我自己复跑，不引 197-r2 的数）。

### 4.1 有牙（把策略改回"两键并一行"就红）——**7 枚**

| 用例 | 锚 | 要什么形状才红／摘掉会怎样 |
|---|---|---|
| `internal/panel/pump_test.go:165 TestTheStreamLogTruncatesInsteadOfMerging` | `:172` `NewStreamLog(2)` 喂 3 键 ⇒ **过界** | `:179-181` `len(chunks)!=3` 直接 `t.Fatalf("…folding two keys into one row is a panel showing one agent doing two agents' work")` ⇒ 任何"并键"当场红。另 `:192-194` 钉 `Done` 不许从邻行 OR（旧折叠的副作用），`:196` 钉 `Truncated()` 必须 true，`:199` 钉短流 `ElidedRunes()==0`（摘掉截断的**记账**也红），`:206-219` 50 键过硬顶那支钉"每行只装自己的 delta"＋`DroppedKeys()` 数对得上 |
| `internal/panel/subagent_stream_197_test.go:57 TestNPlusOneSubagentsEachGetTheirOwnStream` | `:58` bound=3、`:61` 放 4 枚 ⇒ **过界** | `:69-72` 行数≠4 就 Fatal；`:93-97` **逐行**断言"别人的名字不许出现在我这行"；`:106-114` 还要 `Truncated()=true`、`ElidedRunes()=0`、`DroppedKeys()=∅` ⇒ 合并红、"截断不记账"也红；`:129-143` 顺手钉快照四枚键＋每行三枚键（扩契约面会红，不是 merge 红） |
| `internal/panel/subagent_stream_197_test.go:209 TestOverflowTruncatesEachStreamsOwnHeadAndTail` | `:210` bound=2、3 枚 900+ rune 流 ⇒ 过界且过截断阈 | `:251-258` 每行 `[truncated: N runes elided]` 标记**恰好 1 枚**（幂等折叠被重跑就红），`:266-268` `ElidedRunes()` 必须等于逐行账（"truncated 却说不少多少＝贴了标签的老谎言"），`:260-264` 不许含他人的头/尾 ⇒ 合并红 |
| `internal/panel/subagent_stream_197_test.go:275 TestHardCeilingDropsWholeKeysAndNamesThem` | `:276` bound=1 ⇒ 硬顶=2，喂 5 枚 | `:303-305` 红句逐字 `dropped stream %q was folded into live row %q instead of being named` ⇒ 正是 197-r2 §3.5 的 **M2 突变**（"丢弃前先把 oldest 并进下一行"）实测红的那枚；本程没跑，凭据在证据件 `:126-129` |
| `internal/panel/subagent_roster_197_test.go:331 TestTruncationFactsRideThePacket` | `:332` bound=2、3 枚键（其中 600 rune） | `:361-364` 行内标记缺失 ⇒ `t.Fatalf("…describe a fold that never happened")`；`:344-352`/`:403-406` 钉"短的那行不许替别人付代价" ⇒ 把合并塞回来（或让截断记账串门）都红 |
| `internal/panel/subagent_roster_197_test.go:416 TestDroppedStreamsAreNamedOnTheWire` | `:417` bound=2（硬顶 4），喂 6 枚 | `:421-424` `len(dropped)!=2` 即 Fatal ⇒ 合并形（不丢、并掉）⇒ `dropped=∅` ⇒ 红；`:425-429` 还要 `TruncationFor(k).Dropped` 与名册一致 |
| `cmd/wisp/subagent_carrier_197_test.go:524 TestRunPacketReportsTheStreamLogPastItsBound` | `:530` `DefaultStreamKeys+1`＝33 枚真派生 | `:549-552` **先自杀再判**：没过界就 `t.Fatal(…would be vacuous)` ⇒ 这枚不是恒真尺；`:553-564` 三枚计数必须**等于活日志自己的读数**（读者不许自己算）；`:587-589` `roster carries %d subagent rows, want one per stream opened` ⇒ 合并 ⇒ 行数少 ⇒ 红 |

### 4.2 恒真面（对"合并"这一支**不红**）——**2 枚，具名**

| 用例 | 为什么恒真 |
|---|---|
| `internal/panel/subagent_stream_197_test.go:146 TestTwoSubagentsWritingAtOnceDoNotPolluteEachOther` | 夹具 `:147` `NewStreamLog(DefaultStreamKeys)`＝界 32，只放 **3** 枚键 ⇒ **永不进越界那一支**；合并（`enforceBoundLocked`，`pump.go:577`）根本不被调用。它只钉并发互污＋`:204-206`「3 枚键在 32 界里不许报截断」这条**反向**钉。⚠ **同文件头 `:16-17` 的自述"Had leg B kept 'merge', this is the case that goes red"与夹具形状不符**＝文件头写过头，见 §7 第 5 条 |
| `cmd/wisp/subagent_carrier_197_test.go:231 TestRunPacketCarriesTheSubagentItsRosterRowFed` | `:345-346` 有"two tasks merged into one row"的红句，**但该 case 只派 1 枚孩子＝2 枚键，不过 32 界** ⇒ 对 StreamLog 的合并突变不敏感；它咬的是**载体**（孩子页拿到自己的 `echo:`、不混根 prompt）。⇒ 算半枚：摘掉载体它红，摘掉截断它绿 |

### 4.3 不钉这件事的同名形状尺（别拿来抵账）

`cmd/wisp/subagent_stream_key_197_test.go:30`/`:52`、`cmd/wisp/subagent_carrier_197_test.go:430 TestSubagentStreamKeyHasOneMintSite`＝钉**键拼法/铸造点唯一**；`internal/panel/inbound_roster_253_test.go:490 mergeConsts`、`cmd/wisp/leg_sink_gate_131_test.go:795 (*funcInfo131).merge`＝**仪器内部**名册合并（B11 族），与本规则无关。

### 4.4 顶回计数尺

编排者转述（与 35-a6 `verdict.md:184-186`）的"**三枚**既有尺反咬 merge"＝**抽样名册不是整族**：整族现量＝**有牙 7 枚**（4.1）＋恒真 2 枚（4.2）。行号也现量对了一遍：`pump_test.go:165`／`subagent_stream_197_test.go:209`／`:275` 三枚**行号与 35-a6 一致**（未漂）。

---

## 5. 外向影响（谁把"merge"当既定行为写进了判据或用户可见承诺）

1. **判据层（会被写腿执行的那类）**：
   - 未勾 AC 框 **2 枚**＝`.scratch/wisp/issues/35-…:49`（编排者已收窄）与 `.scratch/wisp/issues/36-…:44`（**未收窄**）。票 36 是 `Status: ready-for-agent`、`Blocked by: 30-result-routing-d10, 35-panel-bridge-c17` ⇒ 一旦解阻塞，`:44` 会被当判据兑现。**这是这次普查最实际的一枚外向影响。**
   - `internal/panel/doc.go:10`、`cmd/wisp/panel_pump.go:397`＝产码注释把"合并"写成"transport 的策略"，后续读包职责的腿会当它已被裁。
2. **决策／契约层**：`PLAN.md:2848`（D38(d)）＋`PLAN.md:2971`（C17 ③）＋`SPEC-01:133`——三处**从未随 197-r2 改写**。这是比票面更硬的一处（动它＝人工批准）。
3. **界面文案／产品承诺层：现量零枚被牵连**（具名分开报，这是编排者点名要看的两格）：
   - `frontend/**` 走对象层 `git grep -i -E 'merg|fold|合并|折叠|拼接|并成|聚合|collaps|concat' HEAD -- frontend/` ＝ **66 行命中**（rc=0），逐行读完，**没有一枚**向用户许诺"过载会合并成一行"。
   - 页面渲染契约反而是反向的：`frontend/src/components/result-stream.tsx:15`「a finished chunk's text renders **VERBATIM** as one text node」。
   - 命中的"折叠"全是**视觉 collapse**（`thinking.tsx:105`、`tool-chips.tsx:112`、`sidebar.tsx` 若干、`ai-native/task-rows.tsx:24`），命中的"聚合"是审批卡文案「需逐条确认，**不可聚合**」（`frontend/src/fixtures/harness.ts:64/:147`）。
   - demo/设计稿里唯一像承诺的两处都不是分片合并：`design/old/screens/cost.html:73`「日**聚合**表 cost_daily」、`design/old/screens/approval.html:738`「L1 才**聚合**，L2 永不」。
   - `frontend/fixtures/harness.ts:369-377` 的 `merge(todos, "待办台账.json")`／`dedup(todos)` 是**假 diff 演示字符串**，不是行为承诺。
4. **CI／仪器**：`.github/workflows/*.yml` 的 4 枚 merge 命中全是 **git merge 门**（`slo-fresh.yml:6`、`ci.yml:824/826/904`），与内容合并无关。

---

## 6. 处置所需的事实（⛔ 本程不选形）

**ⓐ 保持截断、逐处收窄"要求合并"**——要动的枚数与仪器牵连：
- 非人工批准射程（票面/报告，6 枚）：`票35:14`、`票35:27`、`票36:44`（**这枚今天没人收窄过**）、`doc.go:10`、`panel_pump.go:397-399`、以及 §A1b 里 **H-11（`154-host-id-never-closed-r1.md:343`）那枚过期读数**——它把 `mergeOverflowLocked` 写成"今天有界＝32＋mergeOverflow"，是最容易被后续程当现状引的一枚。⚠ 台账/证据件按本项目纪律**只追加不删**，H-11 只能"加注"不能改字。
- **人工批准射程**：`PLAN.md:2848`＋`PLAN.md:2971`＋`SPEC-01:133`。前者是 D38 决策行、后者是 **C17 冻结契约四补项的第③项**（`AGENTS.md` §0 第 2 句）。
- 会被打到既有仪器（具名）：**`internal/panel/git_test.go:364 TestGitDimensionHasNoModelCallableTool`** 是**唯一**一枚把 `docs/PLAN.md` 当文本读进去的尺（`:367` `os.ReadFile(…/docs/PLAN.md)`，`:371` 拿 `d34ToolRowMentioningGit` 找反引号里的 `` `git.*` ``，`:402` 的注释自陈"第一版找双引号那形是**瞎而不是干净**"）⇒ 收窄 `PLAN.md:2848` 那句**只要不引入反引号 `git.` 字串就不点它红**，但它是"按内容锚解析 PLAN.md"的既有仪器，改动前该现跑一次。
- d22scan 牵连：ban #8 射程＝`design/`＋`frontend/`（everyFile）＋`internal/`、`cmd/`（goOnly），**`docs/` 与 `.scratch/` 不在内**（`tools/d22scan/main.go:579-586` 现量）。⇒ 收窄 PLAN/票面**不会**动 emoji 门；但改 `doc.go:10`、`panel_pump.go` 落在 `internal/`、`cmd/` 的 `.go` 里＝**注释豁免、字符串不豁免**（`main.go:1060-1076` 那族注释自陈＋Q-46(c)）。
- 判据框数尺（会被"框数前后同数"这类对账打的）：`票35` 未勾 6 枚（含 `:49`）、`票36` 未勾 6 枚／已勾 0（`:44` 是其一）、`票197` 未勾 5／已勾 2（AC#3 未勾）。收窄文字时**框数与极性不许变**（本项目定式）。

**ⓑ 真改成合并**——前提：
- 必须**人工批准**，且不止一枚：`PLAN.md:2848` 已经是"要合并"，真正被推倒的是 `票197:32/:39`（AC#3 未勾但已由 **A406** 收编的规则）与 `pump.go:397-409` 的裁处；`AGENTS.md` §0 第 2 句把"改 `C1–C32`/`D1–D47`"划为人工批准。
- 必须先解决 §4.1 那 **7 枚有牙尺**：改回 merge 会让它们红，而"放宽/改掉既有断言换绿"被 `AGENTS.md` §1.1 与 35-a6 `:189` 第①支**明确禁止**。⇒ 这一形唯一合法路径是把 7 枚尺连同票 197 AC#3 一起**重新裁**（写进台账 `A##`＋人工批准＋逐枚写清被删断言的存在理由，这是本仓"不许为变绿放宽断言"的唯一合法路径，出处 `pending-and-issues.md:8783`）。
- 语义代价（盘上原文，不是本程观点）：`pump.go:400-406` 逐字「the owner clicks one subagent to read THAT agent's work, and a row carrying two agents' sentences would show them something **false about provenance**」＝D30 那枚病从显示侧进来。

**ⓒ 两条链分开、各留各的规则**——要把"哪两处其实不是同一件事"说清，可引的盘上凭据只有这三行：
- `cmd/wisp/panel_pump.go:396-399`（唯一一处把"合并/backpressure"限定为"**推多少次**"而非"并哪几行"）；
- `internal/panel/pump.go:406-409`（截断那一侧的界已从"有多少行"改成"每行装多少 rune"：「The bound is still a bound: it moved from "how many rows exist" to "how many runes each row holds"」）；
- `docs/PLAN.md:2846-2848` 三条背压是**三条不同管子**（音频→ASR＝丢帧并计数；LLM stream→UI＝合并增量；工具并发上限），说明"合并增量"写在的是 stream→UI 那根管子的**频率/攒帧**策略。
⚠ 这一形**今天缺的一块事实**：`票35:27` 的"每面板 bounded queue"**一块没写**（出向整条零产码，凭据 `票35:93` 的四把尺与 `:94` 的"兑现枚数＝0"）⇒ 攒帧合并与行身份合并**在盘上还没有同时存在过**，"两条链各留各的规则"目前是**纸面区分**，不是被测过的区分。要落 ⓒ，第一个动作只能是先有出向通道（票 35 `:47` 的前置），否则两条链的分工没有载体可验。

---

## 7. 具名顶回（编排者转述 vs 盘上原文）

1. **「票 197 leg B 重裁」这一条成立**，但我把它的**盘上身份证据**补齐了：leg B ＝ `197-r2` ＝派单 `2026-09-28-184x-…:85` 的"乙"，代码＝`94ea50f7`，批准＝台账 **A406**（`pending-and-issues.md:8766`）。转述没提批准记录，我按纪律现量到它**在盘上**（不是只剩转述），见 §3.4。
2. **同族冲突不止票 35 那一格，还多两处**：① `票36:44` 未勾 AC 框「renders with merge」；② **`PLAN.md:2848`（D38d）＋`PLAN.md:2971`（C17③）＋`SPEC-01:133` 三处权威文本从未改写**。收窄它们＝人工批准射程。
3. **「三枚既有尺反咬 merge」＝抽样不是整族**：整族现量 **有牙 7 枚 ＋ 恒真 2 枚**（§4）。三枚具名行号我复跑了，未漂（`pump_test.go:165`／`subagent_stream_197_test.go:209`／`:275`）。
4. **`197-…entity.md:121` 那格里的行号已过期**（`:509`/`:514`/`:521`/`:406`/`:408` 现量分别是 `:573`/`enforceBoundLocked :577`/`DefaultStreamKeys :470`/`StreamKeyHardCeilingMultiple :475`）。票面那一格自陈是"现读"，但**现读的对象是当时那棵树**；引用它的腿会撞空。同类：`35-panel-snapshot-pump-r1.md:46` 引 `doc.go:9`（现量在 `:10`）、`154-…:343` 引 `mergeOverflowLocked`（现量 0 命中）。
5. **`subagent_stream_197_test.go:16-17` 的文件头自述不成立**（「Had leg B kept 'merge', this is the case that goes red」指的是 `:146` 那枚用例——它的夹具是 `NewStreamLog(DefaultStreamKeys)`＋3 枚键，**永不过界**，合并突变不碰它）。⇒ 那一格对 merge 是**恒真面**；真正"过界就红"的是 `:57`/`:209`/`:275`。这条按"注释里的断言一律当待验"处理，⛔ 不许因文件头追认 `:146` 为凭据。
6. **`票35:51` 那条收窄自带的时间戳「10-08 11:3x」晚于它自己那枚文件的 mtime**（`stat` 现量 `2026-10-08 11:15:32 +0800`；本腿起手 `date`＝`11:18`）。内容我逐字复跑核过（`pump.go:397` 逐字相同），**只有时间戳这一格对不上**；mtime 只能证明"最后一写"在 11:15，不排除后续程又改过再被还原，所以只登记不裁。
7. **行号撞车**：转述里"产码 `internal/panel/pump.go:397`"逐字成立，但**同一行号在另一枚文件**（`cmd/wisp/panel_pump.go:397`）里说的是**相反的一侧**（"the merge/backpressure rule"）。以后引用必须带全路径，否则会被读成"pump 自己既禁止又要求合并"。

---

## 8. 件与 rc（逐件自落，⛔ 测 rc 那句前无管道）

---- S1 票35 未勾框数
$ grep -c ^[[:space:]]*- \[ \] .scratch/wisp/issues/35-panel-bridge-c17.md
rc=0
6

---- S2 票35 已勾框数
$ grep -c ^[[:space:]]*- \[x\] .scratch/wisp/issues/35-panel-bridge-c17.md
rc=0
4

---- S3 票36 未勾框数
$ grep -c ^[[:space:]]*- \[ \] .scratch/wisp/issues/36-result-history-panel.md
rc=0
6

---- S4 票197 未勾框数
$ grep -c ^[[:space:]]*- \[ \] .scratch/wisp/issues/197-the-panel-must-show-each-subagents-state-and-click-through-to-its-own-streaming-work-page-but-go-side-has-zero-subagent-entity.md
rc=0
5

---- S5 本程未碰产码/测试（工作树里 cmd internal docs .github scripts 的改动枚数，交件时）
$ git status --porcelain -- cmd internal docs .github scripts
rc=0


---- S6 本程自己的写面
$ git status --porcelain -- .scratch/wisp/probes/197/a1
rc=0
?? .scratch/wisp/probes/197/a1/rc-log.md
?? .scratch/wisp/probes/197/a1/sweep.md

---- S7 sweep.md 字节数
$ wc -c .scratch/wisp/probes/197/a1/sweep.md
rc=0
37327 .scratch/wisp/probes/197/a1/sweep.md

---- S8 rc-log.md 字节数
$ wc -c .scratch/wisp/probes/197/a1/rc-log.md
rc=0
8404 .scratch/wisp/probes/197/a1/rc-log.md

---- S9 00-anchor.md 字节数
$ wc -c .scratch/wisp/probes/197/a1/00-anchor.md
rc=0
1078 .scratch/wisp/probes/197/a1/00-anchor.md


> ⚠ 上面 S7/S8/S9 三行的字节数是**取数那一刻**的值；本节自身继续增长会让 S7 过期。⇒ 引字节数请现跑 0，⛔ 不许把 S7 当终态（本程自己就撞上断言会被自己扫的东西写进被扫文件这一形，具名记在回报里）。
