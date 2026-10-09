# 票 292 — `cmd/wisp/subagent_carrier_197_test.go:598-600` 那句"blockedOnApproval 从没在生产 packet 上被断言"**结论还对、理由已经错了两节**（join 今天有两枚来源；per-call corr 让"卡片 corr == 子任务 id"变成**不可能**而不只是"外面造不出来"）

**立票**：2026-10-09 15:0x 编排者（机主令「及时补票」）
**来路**＝票 289 的非实现者验收腿 `289-v1` 的**宽尺另找**（件 `.scratch/wisp/probes/289/v1/10-ac1-four-comments-verdict.md:112` 与 `50-verdict-table-and-roster.md:16`，commit `fcf1d449`/`6c976aa6`）。它发现这一句时，票 289 的写面只有那四枚件、`AC#1` 又明文"⛔ 不许顺手改这 4 枚件里的**其它**注释"⇒ `289-r1` 不动它**是守规**，本格按**票面枚数低估**上报，处置＝**另立一枚**（先例＝票 288、票 289 本身）。
**性质**：⚠ 只改**注释行**，零语义、零功能、完全可逆 ⇒ 低利害，**不上机主清单**，归工程内务。

## 现量（本节全是编排者 2026-10-09 15:0x 在 HEAD 上现跑；⚠ 引用前先重跑，行号会漂）

- 那三行逐字（`sed -n '598p;599p;600p' cmd/wisp/subagent_carrier_197_test.go`）：
  - `:598` `//  1. blockedOnApproval is never asserted on a production packet. The join needs a card`
  - `:599` `//     whose correlation id equals a CHILD task id, and nothing in this tree can make`
  - `:600` `//     that happen from the outside: the agent loop never sends tool_choice, so mockllm`
- **ⓐ "The join needs a card whose correlation id equals a CHILD task id" 只说了半个 join**。今天 join 的建表处有两枚来源，逐字（`internal/panel/subagent_roster_197.go`）：`:212 waiting := make(map[string]bool, len(pending)+2*len(l1Windows))`、`:214 indexWaitingKey(waiting, card.CorrelationID)`（卡片那一半＝它说的这一半）、`:219 indexWaitingKey(waiting, w.CorrelationID)`／`:220 indexWaitingKey(waiting, w.TaskID)`（**L1 确认窗口那一半，比它说的宽：窗口可以只凭 task id 点亮行**）。查询处 `:236 BlockedOnApproval: row.TaskID != "" && waiting[row.TaskID]`。
- **ⓑ 卡片那一半的形状今天已从"外面造不出来"变成"不可能"**：`internal/agent/loop.go:603 func callCorr(taskID, callID string, index int) string` 两支返回 `fmt.Sprintf("%s#call-%d", taskID, index)` 与 `taskID + "#" + callID`（`289-v1` 现量、永不等于 `taskID`），用在 `:676 TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i)`。⇒ 由 loop 派发的卡片，其 corr **结构上不可能**等于任何 task id；那句把"不可能"写成了"只是从外面做不到"。
- **ⓒ "nothing in this tree can make that happen from the outside" 这一条今天**仍成立**，但成立的方式也变了**：L1 那一半在**仓内只有测试生产者**（尺＝`grep -rn "L1Windows:" --include=*.go .` ⇒ **1 命中**＝`internal/panel/subagent_blocked_220_test.go:81`），产码装配点 `cmd/wisp/run.go:699 rt.pump = panel.NewSnapshotPump(panel.PumpSources{` 里 `L1Windows` **0 命中**（尺＝`grep -c L1Windows cmd/wisp/run.go` ⇒ 0）。⇒ **"生产还没人能点亮它"是真话、根因是"票 220 落地腿（给 `PumpSources.L1Windows` 供数据那一跳）还没派"**，⛔ 不是那句写的"卡片造不出来"。
- **ⓓ 它下一环"the agent loop never sends tool_choice"我量了＝仍对**：`grep -rn "ToolChoice" internal/agent/*.go` 去掉 `_test` ⇒ **0 命中**；全仓非测试的赋值点只有 `internal/llm/probe.go:91 ToolChoice: &ToolChoice{Mode: ToolChoiceRequired}`（probe，不是 loop）。⚠ 再下一环"so mockllm never answers with a tool call / so a spawned child never reaches the gate"**我没裁**＝归本票 `AC#1`。
- ⛔ **不许把 `subagent_blocked_220_test.go:151-152` 那枚 fixture 读成"生产能造出来"**（`{CorrelationID: "corr-host-9", TaskID: "child-1"}` ⇒ `child-1` 行读 true）。那是**测试内注入**，它证的正是上面 ⓒ 的**反面**：仓里能造，生产里没人造。
- ⛔ `cmd/wisp/subagent_blocked_197_test.go:7-10` 那一族引文是"引用＋当场推翻"的写法，**不是假话**，⛔ 不随本票入账（`289-v1` 逐字这么交代，别把它算成第五枚/第六枚）。

## 要建什么

- [ ] **AC#1 逐环裁真伪（⛔ 零产码改动，只读）**：把 `:598-600` 拆成四问逐问作答，每问带 `文件:行` 逐字引文：ⓐ join 的来源枚数（对 `subagent_roster_197.go:212-221`）ⓑ "card corr == CHILD task id" 今天是**不可能**还是**只是不可达**（对 `loop.go:603/:676`）ⓒ "nothing in this tree can make that happen" 的射程（仓内 fixture vs 生产装配，`cmd/wisp/run.go:699` 那枚 `PumpSources` 里 `L1Windows` 有没有值）ⓓ "so mockllm never answers with a tool call, so a spawned child never reaches the gate" 两环今天真伪（尺＝`grep -rn "ToolChoice" internal/agent/*.go` ＋ `internal/llm/adaptertest/mockllm.go` 能不能产出 tool_call ＋ 子代理真举过卡没有）。完成判据＝四条各有引文，⛔ 不许只答"没找到"，要附搜过的字样名册。
- [ ] **AC#2 只改这三行注释，且行中性**：改写后的句子必须**同时**说实话的三节——join 有两枚来源、卡片那一半的形状今天不可能、生产侧 L1 那一半还没有供数据的人（⛔ 不许写成"已经能点亮"）；⛔ 只许注释行，任何 `func`/`if`/字段一字不动；⛔ 不许动这枚件里的**其它**注释、⛔ 不许动任何断言（`cmd/wisp/subagent_carrier_197_test.go:311-313` 那句 `t.Error` 是本件的既有判据，⛔ 不许为了配合新叙述改它）。尺（照票 289 `AC#2` 那两把）＝`git diff --numstat` 逐枚 `+N==-N`；`git diff -U0` 里非 `+++`/`---` 行中**纯注释行枚数＝全量枚数**（非注释行＝0）。
- [ ] **AC#3 没伤到用例**：`go test ./cmd/wisp/ -count=1` 改前改后各**两发**，⚠ 必须带 sherpa harness 逐字 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（不带 DLL 时的 `0xc0000135` 是环境红不是被验物；`0xc000013a` 是另一族）。判据＝**逐名作差＝新增红 0 枚**，⛔ 不是"全绿"。⚠⚠ `cmd/wisp` 整包逐名尺**同一棵树两发之间自己会变 ±1 枚**＝`289-v1` 现量（`50-verdict-table-and-roster.md:20`，第二发翻出 `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`）⇒ 所以本格要求**每态 ≥2 发取交集**，并把"语义面零"（`AC#2` 的 md5/非注释行那把尺）作为主证、名册作差作为辅证。
- [ ] **AC#4 门禁与越界**：`GOFLAGS= go build ./...` rc=0；`sh scripts/d22scan.sh` rc=0；`gofmt -l` 空；`gofumpt` 用 `$(go env GOPATH)/bin/gofumpt.exe -l`（裸 `gofumpt` 在 PATH 里 rc=127，⛔ 不许当成"跳过"）；`git show --stat` 名册只含这枚件＋`probes/292/**`；`frontend/**`／`design/**`／三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）／golden／`thresholds.go`／`allowlist.txt` 零字节；**每一把门禁件自己落一行 `rc=N`**（⛔ 0 字节的件＝那格没交）；⛔ 零 push、commit 必带显式 pathspec。

## 禁区

⛔ 不许顺手给 `PumpSources.L1Windows` 接生产供数据——那是**票 220 落地腿**的活（我名下队列里那格"给 L1Windows 供数据那一跳"），接了就会把这枚 fixture 与生产混成一谈；⛔ 不许翻票 220／票 197／票 289 的任何框；⛔ 不许为变绿放宽任何既有断言；⛔ **不搭 `290-r1` 的车**（那一发有它自己那一个改动的成对读数要护），也不与票 288 混批（它 `AC#2` 要的是它那一个改动前后的成对读数）。

**Status:** **未开工**。排程＝**Go 编译/测试面空出来之后单独一小发**（写面只有 `cmd/wisp` 里一枚件的三行注释，很便宜），排在 `290-r1` 交件并复跑之后。⛔ 零翻框、零 push。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
