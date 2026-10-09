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

- [ ] **AC#1 四句订正成说实话**：每句改成"loop 每枚工具调用各铸一枚 corr（`callCorr`），形如 `<taskID>#<callID>`；与 task id 相等那一说自 `bd124b2a` 起不再成立"这一族事实，并按现锚改行号引用（`loop.go:603` 定义／`:676` 使用）。⛔ **只许改注释行**，任何 `func`／`if`／结构体字段一字不动；⛔ 不许顺手改这 4 枚件里的**其它**注释；写面**只许**这 4 枚件。
- [ ] **AC#2 行中性自证（本票的牙）**：改前后对每枚件跑 `git diff --numstat` ⇒ **每一枚必须 `+N` 与 `-N` 相等**（同枚数换行；出现 `+N/-M` 且 `N≠M` ⇒ 本格判失败，必须当场把多出来的行折回去，⛔ 不许解释成"排版需要"）。另跑一把行号漂移尺：`git diff -U0 <这 4 枚件>` 里出现的**纯注释行**（`^[+-][[:space:]]*//`）枚数＝diff 全量非 `+++`/`---` 枚数 ⇒ **非注释行＝0**。
- [ ] **AC#3 没伤到任何用例**：`go test ./internal/panel/ ./cmd/wisp/ -count=1` 改前改后各一次，**三数（PASS/FAIL/SKIP）并排＋红名册逐名作差＝新增红 0 枚**。⚠ 改前基线今天不是绿的：`internal/panel` 现量 **5 枚具名红**（`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`＝票 35 `:49` 那格有意保留的两枚，`b5439299` 交件即红；`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`＝前端契约与 design 资产那一族），`cmd/wisp` 走 sherpa DLL harness（`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，⛔ 不带 DLL 时 `0xc0000135` 是环境红不是被验物）。⚠⚠ **第三枚红 `TestC21DesignTokensFourWayAgree` 今天红在"本机工作树里 `design/assets/tokens.css` 被人删了尚未入库"**——那是别人的在飞改动，⛔ 本腿不许还原它、也不许据它判绿。⛔ 不许因为"整包本来就红"就不跑改前基线。
- [ ] **AC#4 门禁与越界**：`GOFLAGS= go build ./...` rc=0；`sh scripts/d22scan.sh` 纯净树 rc=0；`gofmt -l`／`gofumpt -l` 对动过的 4 枚件空；`git show --stat` 名册**只含这 4 枚件**；`frontend/**`／`design/**`／三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）／golden／`thresholds.go`／`allowlist.txt` 零字节；⛔ 零 push、commit 必带显式 pathspec。

## 禁区

⛔ 不许动 `internal/agent/loop.go:367` 与 `:594`——那两句注释**已经是实话**（`:367` 逐字 `// id per tool call at the dispatch hop (callCorr), so the request-side corr`），别"顺手统一"把它们改回旧说法；⛔ 不许为变绿放宽任何既有断言；⛔ 不许把票 288 那处 `cmd/wisp/panel_inbound_guards_35r3_test.go` 对齐捎进本票（它有自己的 AC#2 归因要求）；⛔ 不许把票 279 的 P39 贴进来（那要 `run.go` 因实事被改才触发）。

**Status:** **未开工**。排程＝排在 `247-r1`（票 247 采集接线，独占 `cmd/wisp`＋整个导入图）**之后**；本票单独一发很便宜（4 行注释），但**必须独占 Go 编译面**，因为 `AC#2` 的"改前基线"要求树是干净的 HEAD。⛔ 零翻框、零 push。
