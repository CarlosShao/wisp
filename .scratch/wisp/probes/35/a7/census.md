# 35-a7 只读普查：面板扇出/单条发布的**字节上限**在规格/契约面定过没有

- 腿：`35-a7`（票 35 那一格 `:49` Backpressure，只读；⛔ 零产码、零测试、零票面、零文档改动、零 push、零 `go build`/`go test`）
- 现量锚：`git log --oneline -1` = **`b5439299`**（我自己量的，未采信派单里的号）
- 起手闸门：`ls .scratch/wisp/probes/35/a7/` = **不存在（rc=2）** ⇒ 槽位空，继续
- 时刻：2026-10-08
- 引用口径：下列 `<ref>:<file>:<line>` 的 `<ref>` 一律＝`HEAD`＝`b5439299`。跨树行号存在性先量过（`logs/cross-tree.txt`：`docs/PLAN.md` 3601 行／`internal/panel/pump.go` 712／`cmd/wisp/panel_pump.go` 420／`cmd/wisp/panel_resident_windows.go` 553／`internal/observe/redact.go` 224／`internal/observe/thresholds.go` 208／`SPEC-08` 249／`SPEC-01` 217／`SPEC-10` 142／`197-subagent-stream-r2.md` 225／票 35 425／`backpressure_bound_35r8_test.go` 371／`docs/SLO.md` 583；`git cat-file -e` 逐枚 rc=0）。
- 脏面自查：`git status --porcelain -- <全部被引文件>` = **空输出（rc=0）** ⇒ 我读到的盘上行＝HEAD 行，行号可直接引。
- ⛔ **自纠一处尺形**：`logs/r1-vocab.txt` 每个 term 后面那行 `rc=` 是**管道尾 `tee` 的 rc，不是 `git grep` 的**（我在测 rc 前挂了管道，违反派单 §4）。同一件事用**无管道尺**重跑＝`logs/r6-nopipe.txt`，那才是要引的那把。

---

## Q1 定过没有 —— **结论：③完全没定过**（就"扇出/单条发布/驻留文本的字节或字符上限"这一格，PLAN.md + docs/specs + docs/SLO.md 三处权威文本零命中；②那一支的相邻数我只列作"射程不覆盖"，⛔ 不许当数用）

### ③ 没定过（主结论，尺＝`logs/r6-nopipe.txt`）

```
term=字节上限        PLAN+specs+SLO rc=1  (1=no hit, 0=hit)
term=驻留上限        PLAN+specs+SLO rc=1
term=环形缓冲        PLAN+specs+SLO rc=1
term=无界内存        PLAN+specs+SLO rc=1
term=heap cap        PLAN+specs+SLO rc=1
term=单条发布        PLAN+specs+SLO rc=1
term=背压            PLAN+specs+SLO rc=0   ← 只有"背压"这个词在规格面存在，见下
explicit-listing rc=1                      ← 五枚词一起列的合并尺，同样零命中
```
⇒ **规格面写了"要有背压/要有界"，从没写过"多少字节"。** 这三句是全部射程：

1. `HEAD:docs/PLAN.md:2845-2848`（D38(d) 背压，逐字）：
   > **（d）背压**：
   > - 音频 → ASR：**有界 channel（≤200ms 音频）**，满则丢帧**并计数**；
   > - LLM stream → UI：有界 channel，满则**合并增量**（不丢内容）。

   ⚠ 这一行今天**两处过期**：容量没有数（音频那路才有 `≤200ms`），而"合并增量"已被票 197 leg B 改判成**截断不合并**——权威件 `HEAD:docs/evidence/s1/197-subagent-stream-r2.md:41`「## 1. 选型：乙（显式标"已截断" + 每键自留头尾），为什么不选甲」，`HEAD:docs/reports/pending-and-issues.md:8766` `A406` 收下（⚠ 该件自陈是实现者自述，最终凭据归 `A406`）。

2. `HEAD:docs/PLAN.md:2971`（D39 给 C17 补的四项，逐字）：
   > ③ correlationId 路由 + 事件推送的背压与合并（D38(d)）
   ⇒ 尺寸问题被**指回 D38(d)**，而 D38(d) 没有数 ⇒ 一个环。

3. `HEAD:docs/specs/SPEC-01-architecture.md:131-133`（spec 侧同形复读，逐字）：
   > - 音频→ASR：有界 channel ≤200ms 音频，满则丢帧**并计数**（计数进日志与诊断包，静默丢帧=识别率
   > - LLM 流→UI：有界 channel，满则**合并增量**（不丢内容）。

4. `HEAD:docs/PLAN.md:1367`（C17 本体，逐字，无任何尺寸）：
   > | **C17** | **`PanelBridge`** | 前端↔Go 双向通道：`invoke(method, args) → result` + Go→前端事件推送（流式结果、审批请求、任务状态）。**回复必须按 correlationId 路由**；前端**必须无状态**（WebView 销毁后一切从 Go 侧重读） | D29, D31 |

5. `HEAD:docs/specs/SPEC-08-ui-ball-panel.md:156-159`（球与面板的 C17 节，逐字；标题自带"要走批准"）：
   > ### 5.2 C17 PanelBridge 方法白名单【SPEC 提案，S5 定稿走契约批准】
   > `invoke(method, args) → result` + Go→前端事件推送（流式结果/审批请求/任务状态），回复按
   > correlationId 路由；每方法标注 capability 与是否需原生侧授权；未列出方法名 → 拒绝并记日志。

6. `HEAD:docs/PLAN.md:2743`（D36 配置树里 `[panel]` 那一行，逐字）：
   > | `[panel]` | `enabled` `width` `height` `keep_alive_in_session`(true) `scale` | `hot` |
   ⇒ 想做成"用户可配的上限"也**没有插槽**，加一枚键＝改 D36 文本＝人工批准。

### ② 相邻数（**射程都不覆盖这一路**，逐枚说差在哪；列出来是为了让下一位别误拿）

| 数 | 出处（逐字） | 射程 | 差在哪 |
|---|---|---|---|
| `≤256KB，超出截断并标 truncated` | `HEAD:docs/PLAN.md:595`（= `:2493` 同句）、`HEAD:docs/specs/SPEC-06-security-gatekeeping.md:137`、`SPEC-07-tools-and-plugins.md:116` | **goja 插件单次结果**进 host bridge | 不是面板扇出；发生在工具侧 |
| 原始输出**硬上限 1MB** | `HEAD:docs/PLAN.md:432`、`HEAD:docs/specs/SPEC-05-agent-core.md:118` | **工具 stdout 进上下文/spill** | 同上：进 LLM 前的界，不是 `StreamLog` 驻留的界 |
| 截断上限 `≤32KB` | `HEAD:docs/PLAN.md:2662`、`SPEC-07:111` | D46 command 插件的**声明式提取器**输出 | 同上 |
| `MaxLoggedString = 512` | `HEAD:internal/observe/redact.go:36-37`「// MaxLoggedString bounds any single string that reaches the log (rule 4).」＋`:141`；`summaryClamp = 440` 在 `HEAD:cmd/wisp/panel_pump.go:367` | **进日志的串** | ★ 规格自己承认管不到这一路——`HEAD:cmd/wisp/panel_pump.go:293-306` 逐字：「every string that reaches it goes through internal/observe's redacting handler, which bounds a single string at observe.MaxLoggedString = 512 characters … **A panel packet is routinely longer than that.** … the smallest packet this shape produces is 552 bytes / 518 runes」⇒ 因此记账退化成摘要行。**"面板 packet 天然越过 512"是盘上现量事实，不是我的推断** |
| 建议 256 KiB（整包 diff 维） | `HEAD:docs/evidence/s1/189-uncommitted-count-design-a1.md:251` 逐字：「② **整包**：一次快照里 diff 维总字节上限（建议 256 KiB）｜本票形状选择，非仓内既有值」＋「owner 不批我不改」 | 未提交计数那一格 | ⛔ **它自己标了"非仓内既有值"**，属提案，不是定过的数 |
| 面板开启 **≤600MB**／工作峰值 ≤700MB／`长输出临时缓冲` 列为可自动卸载项 | `HEAD:docs/PLAN.md:2259`、`:2260`、`HEAD:docs/SLO.md:52`、`HEAD:internal/observe/thresholds.go:23`；卸载清单 `HEAD:docs/PLAN.md:2280`「可自动卸载的只有：面板 WebView（若用户没在看）· **长输出临时缓冲** · goja VM · BM25 索引缓存」 | **进程树私有内存**（C30 `TreePrivateBytes()`），按态 | 口径是整棵树 MB、不是单条发布；且"长输出临时缓冲"被规定为**可丢**却从未被规定**多大**——②类里最接近"规则能推出数"的一枚，缺的正是那枚数 |
| 每行保留 `2×streamTruncateKeepRunes = 512` runes（**仅溢出态**） | `HEAD:internal/panel/pump.go:470/:475/:480`（`DefaultStreamKeys = 32`／`StreamKeyHardCeilingMultiple = 2`／`streamTruncateKeepRunes = 256`，注释逐字「Runes, not bytes: a byte cut would split the Chinese this project streams and hand the panel mojibake.」） | 产码内部常量，**非规格** | 它只在 `len(order) > maxKeys` 后置位生效（`:577-581`）⇒ 键数不越界时**一克都不削**。数值与 `MaxLoggedString=512` 撞巧＝**巧合，不构成授权**（单位还不同：rune vs 日志侧"chars"） |

### ① 有明确数 —— **无**。（若有人声称有，请让他给出上面哪一行；上面每一行的射程都在表里写明。）

---

## Q2 这算不算契约面 —— **结论：不算"一枚实现常量"就能落地；我数到 9 处文本会被顶到，⛔ 不替机主裁**

会被顶到的文本（逐枚点名，含"顶法"）：

| # | 文本 | 会被怎么顶 | 批准层级（照原文措辞，不是我发明的） |
|---|---|---|---|
| 1 | `HEAD:docs/PLAN.md:2848` D38(d)「LLM stream → UI：有界 channel，满则**合并增量**」 | 要写"每行/每次发布的字节上限"就是给这一行加第三句；且这行现文与已裁语义冲突（merge→truncate） | D38 = `D1–D47` 之一 ⇒ **人工批准**（`AGENTS.md` §1.1；`PLAN.md:2818` 标题原文「### D38 — 并发与线程模型（**原文只有一行，而 D32 的可达性完全依赖它**）」） |
| 2 | `HEAD:docs/PLAN.md:1367` C17 契约行 | 若上限作为桥的承诺上线（"每帧 ≤ N"），C17 本体要多一列 | `C1–C32` ⇒ **人工批准** |
| 3 | `HEAD:docs/PLAN.md:2971` D39③ | 现文把尺寸指回 D38(d)；落数就等于宣告这行指错了地方 | D39 ⇒ **人工批准** |
| 4 | `HEAD:docs/PLAN.md:2259-2261` D32 16.3.2 表（面板 ≤600MB 那一行）＋ `:2280` 可卸载清单 | 若把驻留文本做成 SLO 行／做成"超限即卸"的动作 | **改那张表＝人工批准**（票面 §4 明示，`AGENTS.md` §1.1 同句） |
| 5 | `HEAD:docs/specs/SPEC-01-architecture.md:127-133` §4.3 | PLAN 改了它必须跟（`docs/specs/README.md`：PLAN 定案最高，spec 不得矛盾） | spec 文本，随 D38 走 |
| 6 | `HEAD:docs/specs/SPEC-08-ui-ball-panel.md:156` §5.2 标题「【SPEC 提案，S5 定稿走契约批准】」 | 若上限是桥的方法/事件级承诺，就在这一节定稿 | 原文自陈 **走契约批准** |
| 7 | 票 35 面三处：`HEAD:.scratch/wisp/issues/35-panel-bridge-c17.md:27-28`「Event push: bounded queue per panel; overflow MERGES increments (never drops content); backpressure counters visible.」／`:49` 那枚框（含 10-08 编排者〔…〕收窄注，逐字「本框今天**只剩"no unbounded memory（heap cap asserted）"那半句算账**」）／`:392` 那行表（逐字「票面 `:27` 那枚"每面板 bounded queue"**一块没写**（出向零产码）」） | `:27` 的 MERGES 与已裁语义相反；`:49` 只有在实现长出界之后才可能翻勾 | 票面 = 编排者所有；⛔ 腿不许自己翻框 |
| 8 | `HEAD:docs/SLO.md:52`（面板行"未测树口径"那一格）＋ `HEAD:internal/observe/thresholds.go:23` | 若它变成验收行就要进表 | `thresholds.go` **一字节不许动**（`AGENTS.md` §1.1 逐字），只有人工批准面能进 |
| 9 | `HEAD:docs/PLAN.md:2743` `[panel]` 键行 | 若做成用户可配 | D36 ⇒ **人工批准** |

**分层判断（只陈述，不裁）**：纯实现常量那一支（在 `internal/panel/pump.go` 加一枚 unexported const + 让 `Append` 无条件 clamp）**不直接改上面任何一行文本**，但它会把 1/3/7 三处**变成过期描述**，并且必然顶到下节那族按 rune 算术钉的行为尺。

---

## Q3 有没有既有红会被顶 —— **结论：会，而且第一颗会顶到的是"到键界之前不许削"那枚反向钉；d22scan 只咬 ban #8/#9 两形**

### 在册红（⛔ 我没跑，禁 `go test`）
`HEAD:docs/reports/pending-and-issues.md:13952`（编排者自己复跑，逐字）：
> `go test ./internal/panel/` **rc=1**，逐名 **6 枚 FAIL**＝**2 枚新**（`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`、`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`）＋**4 枚起手既有**（`TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestComposerContractTypesMatchFrontend`、`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`、`TestC21DesignTokensFourWayAgree`）
⇒ 我这把对那 6 枚**没有独立读数**，引用＝引册不引用实测。`cmd/wisp` 那约 15 枚窗口依赖红族我**一枚都没现量**（见末节）。

### 新落一枚"每行字节上限"会当场顶到的钉（尺＝`logs/r3-nails.txt`，命中 161 行、rc=0）
| 钉 | 逐字断言 | 会不会红 |
|---|---|---|
| `HEAD:internal/panel/pump_test.go:199-200` | 「short streams keep full fidelity past the bound」 | ★ **最直接的反向钉**：它要求键界之下 `ElidedRunes()==0`。无条件 clamp 每行 ⇒ 当场红 ⇒ 这一枚属"要重写必须写明被删断言的存在理由"那一族（`A406` 立的规矩） |
| `HEAD:internal/panel/subagent_roster_197_test.go:344`/`:399` | `big.ElidedRunes != 600-2*streamTruncateKeepRunes` | 改 `streamTruncateKeepRunes` 或改保留窗形状 ⇒ 红 |
| `HEAD:internal/panel/subagent_stream_197_test.go:249-268` | `elided := produced - 2*streamTruncateKeepRunes`、「'truncated' that cannot say how much it lost is the old lie with a label on it」 | 同上（算术级耦合） |
| `HEAD:internal/panel/backpressure_bound_35r8_test.go:262` | `len([]rune(c.Text)) < 2*streamTruncateKeepRunes` ⇒ `t.Fatalf(...below the head+tail it promises)` | ★ **新尺与旧尺互咬**：若把每行界压到 <512 runes，这枚今天绿的用例转红 |
| `HEAD:internal/panel/pump_test.go:165` `TestTheStreamLogTruncatesInsteadOfMerging`（票 `:392` 点名）＋ `subagent_stream_197_test.go:209`/`:275` | 「三枚既有尺…**会因改回 merge 而红＝唯一有牙的尺，咬的方向与 AC 相反**」 | 只要不动 merge 语义就不红；⛔ 不许为了合票面 `:27` 把实现改回合并 |
| `HEAD:internal/panel/inbound_roster_253_test.go:59`（`wantFullInboundRosterSize253 = 6`） | 票 `:391` 逐字：「一旦铸进 `bridge.go` 的 `Method*` 常量就落进…`:59`…**当场红**（断言体 `:445`/`:787`）」 | 只有"给丢弃/截断点分事件名"那一支才顶它 ⇒ 那是 C17 名册面＝人工批准 |
| `HEAD:internal/panel/pump_test.go:123`/`:291` 四键字节钉（`!= "composer,generatedAt,pending,results"`）＋两枚 `…MatchFrontendTypes` | 票 `:191`/`:237` 与 `197-subagent-carrier-r3.md:153` | 只有"把上限/丢弃读数写进 packet"才顶它；⛔ `frontend/**` 一枚字节都不许动（票 `:164`） |

`cmd/wisp` 那一族：只要写面留在 `internal/panel/**`（今天唯一生产铸造点就在 `cmd/wisp/run.go:698`，见 Q4，但改的是 pump 内部常量，不必动装配），这族**不会被新顶**；一旦"在出口拒收"要动 `panel_pump.go` 的记账，就进入那约 15 枚窗口依赖红的射程（我这把**未现量**）。

### `tools/d22scan` 会不会点（尺＝读 `HEAD:tools/d22scan/main.go` 头注释，逐字）
- bans #1-5：`:5-6`「production scope internal/ + cmd/ (**non-test**, non-testdata) unless stated」⇒ 新常量本身不在 #1-5 射程；除非它同时起裸 `go`（ban #1）、用 `filepath.Clean|Abs`（#2）、墙钟差超时（#4）。
- ban #8 emoji：`:40-46`「… in the Go sources of internal/ + cmd/ - **comments and _test.go INCLUDED (D23)**. This is the one ban whose scope is NOT the "production, non-test" default」＋射程 `:186` `[\x{1F000}-\x{1FAFF}\x{2200}-\x{22FF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`
  ⇒ **实测射程确认与派单 §4 一致**：`≤`(U+2264)／`✓`(U+2713)／`⛔`(U+26D4) 写进 `internal/`＋`cmd/` 的 Go **注释也算**；`→`(U+2192) 与带圈数字**不在射程**（22FF 与 2600 之间、2190 段留空）。票 35 落地腿自己撞过一次：`HEAD:docs/reports/pending-and-issues.md:13957`「它初稿把 `⛔` 写进 Go 注释…ban #8 扫 `.go` **含注释**、`U+26D4` 落在仪器射程 `2600–27BF` 段内」。
- ★ **ban #9 phantom-citation**（`:50-56`）「a production-file comment citing a repo-relative path that does not exist on disk … over the production Go files of internal/ and cmd/ ONLY」⇒ 新常量那几句注释若指 `docs/contracts/C17.md`／`docs/SLO.md` 里不存在的行，**门当场红**。本格最容易踩的是这条，不是 emoji。

---

## Q4 不修的最坏形状 —— **结论：无界只在"一枚 `wisp run` 进程的寿命内"这一寸，而 D38(d) 给的并发上界恰好把真实负载推进这一寸；"未修"≠"已经在漏"**

1. **谁发（现量枚数）**：尺 `git grep -n "NewStreamLog" HEAD -- cmd/wisp internal/panel internal/tools` = **rc=0、18 枚命中**（`logs/r6-nopipe.txt`），其中**生产落点只有 1 枚**：
   `HEAD:cmd/wisp/run.go:698`「`rt.stream = panel.NewStreamLog(panel.DefaultStreamKeys)`」——其余 17 枚全在 `*_test.go`。
   ⇒ **常驻面板路径今天没有 `StreamLog`**：票 `:392` 逐字「`cmd/wisp/panel_resident_windows.go:148 make(chan func(), 16)`＋背压计数 `:130/:369 failedPost`（**零测试引用**）。票面 `:27` 那枚"每面板 bounded queue"**一块没写**（出向零产码）」（我现量该行逐字 `tasks: make(chan func(), 16),` 在 `HEAD:cmd/wisp/panel_resident_windows.go:148`，文件 553 行）。

2. **无界的机理（代码事实，逐跳）**：`HEAD:internal/panel/pump.go:533-535`
   > 	if s.truncated {
   > 		s.clampLocked(key)
   > 	}
   而 `s.truncated` 只由 `:577-581` 的 `enforceBoundLocked()` 置位：`if len(s.order) <= s.maxKeys { return }` ⇒ **键数不越 32 就一克不削**；`:591` `s.dropped = append(s.dropped, oldest)` 是 **append-only、无界**，且 `DroppedKeys()` 每次读全量交出（`backpressure_bound_35r8_test.go:316-319` 逐字）。

3. **★最要紧的一条形状（这条派单里没写，是我量出来的）**：真实扇出**天然落在界以下**。D38(d) 给的并发上界是 `HEAD:docs/PLAN.md:2849`「**工具并发上限 = 4**」，子代理池同值（票 211＝`HEAD:docs/reports/pending-and-issues.md:8934`「已立 **票 211**…把甲/乙的界线与 5 条判据写死」），所以"一次任务派 4 个子代理"＝**5 条流 ≪ 32 键界**＝**截断器根本不engage**。也就是说：**规格给的有界性（并发 4）恰好把生产负载送进唯一没有界的那条路。** 仪器正是照着这一形写的：`backpressure_bound_35r8_test.go:297-299` 用 `keys=3`（`flood35r8(t, …, 3, 200, 2000)`／`3, 800, 2000`）。

4. **涨到什么量级（仪器现量，逐字取自那两枚红的断言）**：
   - 扇出以下：1x 留 **1,200,000 runes**（3×200×2000）、4x 留 **4,800,000 runes**（3×800×2000），`t.Errorf` 那句逐字「Retention grew %.1fx for a 4x flood and **equals the flood exactly**」；堆读数 `:309` 打 `live heap +%d bytes`（编排者现量＝`+1.16MB→+4.83MB`）。
   - 点名账本：1x 4000 流 → 3936 名、4x 16000 流 → **15936 名**（`:328` 逐字 `m1.names != 4000-cap.rows`／`:336`「Its cap is the event count, not the bound」），≈212k runes。
   - 对照：扇出**过界**那一路确有界——`retentionCap35r8()`（`:199-203`）＝ `64 行 × (2*256+64)`，`TestStreamLogFanOutFloodKeepsRetentionUnderCap35r8`（`:237-279`）**绿**，即编排者说的 35072 runes／64 行。
   - ⚠ 恒真面（腿自己具名，`pending-and-issues.md:13950`）：heap 那项带 1 MiB 余量、**单独不咬人**；发布次数钉在 24 ⇒ **不量出口那一路的内存**。

5. **撞不撞得到 D32 的空闲预算（⛔ 别夸大）**：撞不到——今天这条驻留在 **CLI 进程**里，`wisp run` 不是 D32 那六态的常驻树；`Sleeping ≤25/40MB` 那行（`HEAD:docs/PLAN.md:2255`）量的是常驻进程树，尺在 `HEAD:internal/observe/thresholds.go:19`。真要构成 D32 问题，**前提是票 35/33 的出向接线把 `StreamLog` 带进常驻面板路径**（还没写）；届时的量级：4 路子代理 × 每路数 MB ≈ 十几 MB 文本，相对 `面板 ≤600MB`（`:2259`）不算击穿、相对 `Warm ≤350MB`（`:2257`）是可看见的一块。⇒ **"最坏形状"是：一次长任务把一条流灌满、消费端一直没读、内存随事件数线性走，直到进程结束；不修它不会让今天的任何 SLO 行变红，但它让票 35 `:49` 那半句永远翻不了勾，并且出向接线一落地就从"CLI 寿命内的浪费"升级成"常驻树里的驻留"。** 判断"要不要现在就为常驻树定数"＝**需要人拍板**（见末节）。

---

## 末节：我没量到的／我认为需要人拍板的那一格

### 我没取到数的格子（具名）
1. ⛔ 未跑 `go test`/`go build`（硬约束）⇒ 2 枚新红＋4 枚既有红＝**引册不引实测**（`pending-and-issues.md:13952`）；`cmd/wisp` 那约 15 枚窗口依赖红族**一枚都没现量**。
2. **`docs/PLAN.md:2219` 一带之外我没找到"面板驻留"的任何预算行**；如果机主认为存在（例如 `docs/SLO.md` 附录某节），请给我行号，我这把的五枚词尺（`r6-nopipe.txt`）覆盖不到词表之外的说法。
3. `internal/streamkey`（票 `:392`/`:490` 指名的别名真源）我只通过 `pump.go:494` 读到别名，**没读它的文件本体**，因此"键形状是否带尺寸语义"未验。
4. 票 35 `:100` 那一行我读的是 `:94`/`:104` 附近段落（派单点名的 `:100` 属 10-07 编排者 §5 排程叙述，与尺寸无关）；`:392` 整行已逐字读。
5. `197-subagent-stream-r2.md` §1「乙」我读了前 70 行（含 §1 全节与 §2 落点表头）；**该件自陈是实现者自述**，最终凭据按派单口径归 `A406`。
6. ⛔ 我没有为"字节上限"提出任何数值建议——③ 就是结论。

### 摆给机主的一问（三栏，含"不做"；每栏一行零术语人话）

| 支 | 一句话人话后果 | 依据（本程现量） |
|---|---|---|
| **甲**（现在给这一路定一枚每行/驻留上限并实现） | 面板显示的文字不再无限增长，代价是要重写"到键界之前不许削"那枚旧断言，还要机主亲手给一个数——**规格里没有这个数，我不发明**。 | `pump_test.go:199-200` 反向钉（§Q3）；Q1 的 rc=1 五连（没数）；`streamTruncateKeepRunes=256` 是现成的"溢出态每行 512 runes"形状（`pump.go:480`），但它**不是授权** |
| **乙**（不定数，把票面 `:49`/`:27` 改成与已裁语义一致的说法并具名登记"驻留尺寸未定"） | 那两枚红灯继续挂着但账目诚实，规格不再假装说过"会合并"，代价是这一格**仍然没修**、票 35 出向接线落地前谁都不付钱。 | 票 `:27` 的 MERGES 与 `A406`/197 已裁相反（§Q1 第 1 条）；票 `:392` 自陈"每面板 bounded queue 一块没写"；CI 代价已在 `pending-and-issues.md:13954` 写死 |
| **不做**（不甲不乙，保持现状） | 什么都不动：今天只有编排者本机红，**下次 push 之后 `internal/...` 那一步 CI 会红**；而"没修"只是说没界，**不等于现在在漏**——无界那一寸活在 CLI 进程寿命里，常驻路径今天还没这条流。 | 生产铸造点现量 **1 枚**（`cmd/wisp/run.go:698`；`git grep` rc=0／18 命中）；常驻面板**无** `StreamLog`（票 `:392`）；CI 后果逐字见 `:13954`；量级读数见 §Q4 第 4 条（1.2M→4.8M runes、+1.16MB→+4.83MB、3936→15936 名） |

### 我认为必须由人拍板的那一格（⛔ 未定义即停，D22 闸门③）
**"面板驻留文本的字节上限"这个数在本仓所有权威文本里不存在（Q1＝③）**；补它要么改 D38(d)/C17/D39/D32/D36 之一（都是人工批准面，Q2 列了 9 处），要么由人裁"以溢出态现有的每行 512 runes 形状推广到所有行"（那是拿实现常量当契约，我不裁）。⇒ **在机主给数或给方向之前，`internal/panel` 这一格的产码不可写**，两枚红灯维持红。

### 顶回编排者（本程撞到的）
- ⛔ 无实质顶回：派单 §1 的三条转述（键数在／字节不在／账本 3936→15936 名）与盘上原文逐枚相符（我读到的是仪器断言原文 `:253`/`:309`/`:336`，数值口径一致）。
- 一处**要补的限定**（不是顶回）：`:392` 与派单 §2 说"三枚计数出口 `:632/:644/:657`"——我量到的置位链是 `pump.go:533-535`（条件 clamp）＋ `:577-581`（键界才 engage），也就是说**扇出有界这件事的成因是"键界"，不是那三枚出口**；引用那一行时建议按"键界 engage 才截断"写，否则下一位会以为出口即有界。
