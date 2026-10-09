# 票 289 — 四句"loop 派发时 `CorrelationID == TaskID`"的注释自 `bd124b2a` 起就是假话；我在票 286 结案时把它们挂到了一条**不存在的链**上

**立票**：2026-10-09 12:1x 编排者（机主令「及时补票」）
**来路**＝票 286 `AC#3` 的名册第五节（件 `.scratch/wisp/probes/286/attr-1/10-registry.md` `:53-62`，腿 `golden-attr-1` 提交 `d106fc25`）里那 **4 枚"待人裁"**。我 12:0x 收口票 286 时把它们写成"归'注释随行'那一族（票 279／票 280 同链）"——**那句话是一张空账**：票 279 只管 `cmd/wisp/run.go` 的 P39 那一块（它的 `AC#1` 到现在还没触发），票 280 已经 `-done` 且裁的是 `gofmt` vs `gofumpt` 的**归属**，两枚都不覆盖这 4 句。⇒ 按票 288 同形处置：**另立一枚**。
**性质**：⚠ 只改**注释行**，零语义、零功能、完全可逆 ⇒ 低利害，**不上机主清单**，归工程内务（先例＝票 288）。

## 现量（本节四条逐字读数＝编排者 2026-10-09 12:1x 在 HEAD 上 `git show HEAD:<path> | sed -n '<N>p'` 现跑；⚠ 引用前先重跑，行号会漂）

| # | `HEAD:file:行` | 该行逐字 | 为什么是假话 |
|---|---|---|---|
| 1 | `internal/panel/subagent_roster_197.go:57` | `	// loop dispatches with CorrelationID == TaskID (internal/agent/loop.go:647),` | **双重过期**：形状——loop 现在每枚工具调用铸一枚独立 corr（`internal/agent/loop.go:603 func callCorr(taskID, callID string, index int) string`，用在 `:676 CorrelationID: callCorr(taskID, p.call.ID, i)`）；行锚——它引的 `:647` 今天已在 `:676`。 |
| 2 | `internal/panel/subagent_roster_197_test.go:468` | `//     CorrelationID == TaskID (internal/agent/loop.go:647); a host that used a` | 同 1。 |
| 3 | `internal/panel/subagent_blocked_220_test.go:155` | `	// The production pairing (the loop dispatches with CorrelationID == TaskID,` | 同 1。⚠ 它配对的 `:157` fixture 用 `corr="child-1"`／`TaskID=""`，**用例本身仍是绿的**——过期的只有这句叙述，⛔ 不许把本格读成"这里有一枚红"。 |
| 4 | `cmd/wisp/subagent_carrier_197_test.go:20` | `// bridge with the ids the agent loop itself uses (TaskID == CorrelationID,` | 描述的是 `:173` 那枚 fixture，把它说成"loop 自己用的 id"＝过期叙述（同 1）。 |

- ⚠⛔ **本格最硬的一条约束是"行中性"**：`subagent_roster_197.go` 是**产码件**，而票 220 的现量表按行号引它（`:190-210` 那枚"用 corr 建表、用 taskID 查"的定案、`:127` 那枚 `blockedOnApproval` 是布尔不是 D43 态），票 197 `:64`/`:66` 也引 `:99-104`/`:127`。⇒ 只要在这枚件里**净增删行**，那些**写在票面里的锚就会静默腐烂**（本仓 10-08 已为同形踩过一次：夹具一插就打死锚位）。
- 对照先例（同一族风险，⛔ 不在本票射程）：`cmd/wisp/config_readers_255.go:109-110` 按行号逐字引 `cmd/wisp/run.go:991`，而核它的用例 `cmd/wisp/config_receipt_255_test.go:179` **按 `lines[n-1]` 直取那一行、不做短语搜索** ⇒ 票 279 `AC#2` 已实测"插 +3 行 ⇒ 那一发真红（rc=1）"。本票的 4 枚件里**没有**这种按行号取字的产码尺（尺＝`grep -rn "lines\[[0-9]" --include=*_test.go internal/panel/ cmd/wisp/` 现跑 0 命中），但**票面锚腐烂**这一维照样要防，所以下面 `AC#2` 把行中性写成硬判据，而不是"注意一下"。

## 要建什么

- [x] **AC#1 四句订正成说实话**：每句改成"loop 每枚工具调用各铸一枚 corr（`callCorr`），形如 `<taskID>#<callID>`；与 task id 相等那一说自 `bd124b2a` 起不再成立"这一族事实，并按现锚改行号引用（`loop.go:603` 定义／`:676` 使用）。⛔ **只许改注释行**，任何 `func`／`if`／结构体字段一字不动；⛔ 不许顺手改这 4 枚件里的**其它**注释；写面**只许**这 4 枚件。
- [x] **AC#2 行中性自证（本票的牙）**：改前后对每枚件跑 `git diff --numstat` ⇒ **每一枚必须 `+N` 与 `-N` 相等**（同枚数换行；出现 `+N/-M` 且 `N≠M` ⇒ 本格判失败，必须当场把多出来的行折回去，⛔ 不许解释成"排版需要"）。另跑一把行号漂移尺：`git diff -U0 <这 4 枚件>` 里出现的**纯注释行**（`^[+-][[:space:]]*//`）枚数＝diff 全量非 `+++`/`---` 枚数 ⇒ **非注释行＝0**。
- [x] **AC#3 没伤到任何用例**：`go test ./internal/panel/ ./cmd/wisp/ -count=1` 改前改后各一次，**三数（PASS/FAIL/SKIP）并排＋红名册逐名作差＝新增红 0 枚**。⚠ 改前基线今天不是绿的：`internal/panel` 现量 **5 枚具名红**（`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`＝票 35 `:49` 那格有意保留的两枚，`b5439299` 交件即红；`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`＝前端契约与 design 资产那一族），`cmd/wisp` 走 sherpa DLL harness（`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，⛔ 不带 DLL 时 `0xc0000135` 是环境红不是被验物）。⚠⚠ **第三枚红 `TestC21DesignTokensFourWayAgree` 今天红在"本机工作树里 `design/assets/tokens.css` 被人删了尚未入库"**——那是别人的在飞改动，⛔ 本腿不许还原它、也不许据它判绿。⛔ 不许因为"整包本来就红"就不跑改前基线。
- [x] **AC#4 门禁与越界**：`GOFLAGS= go build ./...` rc=0；`sh scripts/d22scan.sh` 纯净树 rc=0；`gofmt -l`／`gofumpt -l` 对动过的 4 枚件空；`git show --stat` 名册**只含这 4 枚件**；`frontend/**`／`design/**`／三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）／golden／`thresholds.go`／`allowlist.txt` 零字节；⛔ 零 push、commit 必带显式 pathspec。

## 禁区

⛔ 不许动 `internal/agent/loop.go:367` 与 `:594`——那两句注释**已经是实话**（`:367` 逐字 `// id per tool call at the dispatch hop (callCorr), so the request-side corr`），别"顺手统一"把它们改回旧说法；⛔ 不许为变绿放宽任何既有断言；⛔ 不许把票 288 那处 `cmd/wisp/panel_inbound_guards_35r3_test.go` 对齐捎进本票（它有自己的 AC#2 归因要求）；⛔ 不许把票 279 的 P39 贴进来（那要 `run.go` 因实事被改才触发）。

**Status:** **done**（2026-10-09 15:0x 编排者结案：四格 `AC#1..AC#4` 全勾，凭据全出自**非实现者**验收腿 `289-v1` 的 1:1 裁决表 `.scratch/wisp/probes/289/v1/50-verdict-table-and-roster.md:7-10`（README 规则 6 那张表），⛔ 无一格按实现者自述追认）。**原句逐字留档（作废）**：「**Status:** **未开工**。排程＝排在 `247-r1`（票 247 采集接线，独占 `cmd/wisp`＋整个导入图）**之后**；本票单独一发很便宜（4 行注释），但**必须独占 Go 编译面**，因为 `AC#2` 的"改前基线"要求树是干净的 HEAD。⛔ 零翻框、零 push。」——作废原因＝`247-r1` 已交并验收完、本票已按 `289-r1`→`289-v1` 走完。
★**结案不等于零残余**（README:201 那句原话照搬）：本票带着两笔具名残余收口——① `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` 那枚"改前红／改后绿"**原因未定**（`289-v1` 21 发零复现，我复跑了它的语义面尺；归**编排者名下**，要查得另开一小发）；② 同一族还有**第五句**过期叙述＝`cmd/wisp/subagent_carrier_197_test.go:598-600`，已**另立票 292**（`.scratch/wisp/issues/292-carrier-197-test-598-600-join-shape-sentence-stale.md`），⛔ 不算本票已修。改名记录：本件原名 `289-four-comments-still-claim-the-loop-dispatches-corr-equals-taskid-and-my-286-close-out-hung-them-on-an.md`（文件名 **108** 字符，超 README 规矩 9 的 `≤100` 帽、把 `sh scripts/check-path-length-budget.sh` 打成 `VERDICT RED`）⇒ 2026-10-09 15:0x 改短名为本名＋`-done`；⛔ 零 push。

## 编排者按非实现者验收腿 `289-v1` 判语的处置（2026-10-09 14:4x；⛔ 题面与 AC 原句一字不改，只追加）

判语腿＝`289-v1`（≠ 实现者 `289-r1`，`SPEC-12 §4.3` #1/#3），三笔 commit `e7af1796`→`fcf1d449`→`6c976aa6`，件 `.scratch/wisp/probes/289/v1/`（正文六枚 2,652／10,708／6,462／9,639／6,163／8,848 字节＋raw 十一枚，最大 `raw-v1-post.md` 343,376）。

- **四格判语全部"成立"⇒ 本格翻勾 `AC#1`/`AC#2`/`AC#3`/`AC#4`**（凭据全出自它自己现跑；我复跑过其中三把：`--numstat` 逐枚 `+N==-N`、`-U0` 非注释行＝0、**剥掉整行注释后四枚件 pre/post/HEAD 三态 md5 逐枚相同**＝语义面零变化，这一把我自己复现成功）。
- ★**它没收下实现者那句"串跑时序"的归因，我也跟着不收**：`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` 的"改前红／改后绿"被它用 21 发复跑（隔离×3、`-count=5`、组串跑、`-cpu=1`×2、5 进程并发、原状×2）压过——**负载把 2.0s 压成 2.8–3.0s，但红 0 发**；判语＝**"转绿原因未定"**，只具名记下形状（`:542 awaitCard` 轮询 gate 注册表、`:543` 直读 stdout 缓冲、两读之间**零同步边**，而同 harness 里 `:148 awaitStdout` 没用上）。⇒ 本票只主张"没新增红"，⛔ 不许把这枚转绿记成 `289-r1` 修的，也不许记成"时序 flake 已定案"。**归我名下**：那一枚的成因要查得另开一小发（不在本票射程）。
- ★**新料：宽尺另找到第五句过期叙述**（票面那四句没含它）＝`cmd/wisp/subagent_carrier_197_test.go:598-600` 一带逐字 `//     whose correlation id equals a CHILD task id, and nothing in this tree can make`——在"corr 每枚调用各铸一枚"的今天，"卡片的相关 id 等于子任务 id"这个条件**根本不可能成立**，那句叙述把 join 的所需形状说错了（反例形状见同族 `internal/panel/subagent_blocked_220_test.go:151-152`；我 14:4x 现读复量该行存在）。⇒ ⛔ **不塞回本票**（本票四格已全勾、判的是那四句），**转立一小发＝票 292**（同一族、同一枚件类型、单独一小腿很便宜），⛔ 不由已交回的腿顺手补、⛔ 不搭 `290-r1` 的车（那一发有它自己那一个改动的成对读数要护）。
- **它顶正我派单的一处路径错（记我，不追认）**：我在派单里把 `ticket283_corr_identity_rulers_test.go` 写在 **`internal/panel/`**，真身＝**`internal/tools/ticket283_corr_identity_rulers_test.go`**（尺＝`git ls-tree -r --name-only HEAD | grep ticket283_corr_identity_rulers`）。⛔ 同形错今日第二枚（上一枚＝`SPEC-05-agent-kernel.md`，真名 `SPEC-05-agent-core.md`）⇒ **定式**：派单里凡是"某枚件在某个目录"这句话，落笔前先 `git ls-tree` 现跑，别按记忆写路径。
- ⚠ **一枚测量噪声登记（它的读数，我没复跑＝〔仅腿报〕）**：同一棵树第二发串跑翻出 `TestAC1ResidentLeg…` 带 `0xc000013a`（⚠ **不是** `0xc0000135` 那枚缺 DLL 环境红）⇒ `cmd/wisp` 整包逐名尺**单发带 ±1 枚噪声**，引用"改前几枚／改后几枚"必须带**发数与是否 `-v`**（同 `A773` 真窗族 4↔5 那一族并形）。它报的当期名册：`internal/panel` 真值 **6 枚**（＝我这面 `A779` 刚更正的那枚在内）／`cmd/wisp` 改后 5 枚、改前含 Ticket223 为 6 枚。
- **两处小账我认**：票 289 与票 288 一样**缺 `## Progress log` 那一节**（我建票时漏的结构，两枚腿都按「未定义即停」没自建）⇒ 本节就是补上的那一节；另实现者件 `raw-pre.md` **没有 `rc=` 行**（那格按我的规矩＝没交，已由其验收腿的 rc 补齐）；实现者正文里 `:235` 应为 `:236`。
- **Status 归位＋排程**：四格全勾 ⇒ 本票**具备结案条件**（`-done` 需要"全勾→Status: done→改名→更新索引"四步，我在 `290-r1` 派出后一并做，⛔ 不占 Go 面）。**新加的一格**（第五句订正）⛔ 不与票 288 混批、也不搭 `290-r1` 的车（`290-r1` 的判据要它自己那一个改动的成对读数）——排在小发里单独跑。⛔ 零 push。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。⚠ 本节 2026-10-09 14:4x 由编排者补＝建票时漏的结构。）
- [2026-10-09 14:53:37 +08] agent=289-r1 did=AC#1..AC#4 交付（四枚件注释、行中性两把尺自证、改前后整包三数并排）；具名说做不到：票面缺 Progress log 一节故未自建；自报三处（归因尺那句我说成零 .go 变化＝按字面假、internal/panel 真值 6 枚不是 5、cmd/wisp 那 6 枚票面未列）
- [2026-10-09 14:53:37 +08] agent=289-v1 did=四格判语全部成立（零翻框零源码零 docs 改动）＋攻实现者归因：21 发复跑红 0 发 ⇒ 转绿原因未定、只记形状；宽尺另找第五句过期叙述 cmd/wisp/subagent_carrier_197_test.go:598-600；顶正我派单路径错（ticket283 件真身在 internal/tools/ 不是 internal/panel/）；报 cmd/wisp 整包尺单发 ±1 枚噪声（0xc000013a）
- [2026-10-09 14:53:37 +08] agent=编排者 did=按判语翻 AC#1..AC#4＋补 Progress log 一节＋追加处置节（含我两处派单路径错与第五句的新格归属）next=派 290-r1（甲-1 落地），随后小发补第五句＋结案本票
