# 259-r4 拆载具：判语与读数（票 259 AC#4）

## ① 判语：既有断言**不依赖**「TaskID 与 CorrelationID 相等」

凭据（同一发现取；行号按**本件改动前**的文件版本，即 `git show cbe5cfd1:cmd/wisp/subagent_selfapproval_197_test.go` 那版）：

1. 卡是**按 TaskID 找的**，不是按 corr 找的：`:131` 逐字 `if c.TaskID == taskID && c.Level == "L2" {`（`waitForChildCard197` 的轮询条件；实参 `taskID` 就是 `childID`）。
   卡的 `TaskID` 来源链只吃 `req.TaskID`：`internal/tools/bridge.go:376` 逐字 `TaskID: req.TaskID, CallID: req.CallID,` → `internal/agent/approval/gate.go:600` 同一 Prompt 字面量里的 `TaskID: d.TaskID,` → `internal/agent/approval/replies.go:176` 逐字 `TaskID: p.TaskID,`。
2. 授权/拒绝的每一发都用**卡上现读的** `card.CorrelationID`，没有一行把 `rec.childID` 与 corr 相比：（改动前版）`:428` `:429` `:432` `:438` `:450` `:466` `:470` `:472` `:473` 共九行，模式一律 `...Native().Allow(ctx, cardN.CorrelationID, ...)` ／ `approval.Request{CorrelationID: cardN.CorrelationID, ...}` ／ `rt.liveCards.h.Allow(ctx, cardN.CorrelationID)`。
3. 审计三句期望串由 `rec.corr1`／`rec.corr2`（卡上现读）拼接，不落任何字面 id。
4. 实测同读：改前 rc=0 的 Logf 里 `corr1 == corr2 == child`（同值基线）；改后 rc=0 的 Logf 里 `corr1=<child>-corr-1`、`corr2=<child>-corr-2`，而**六发拒绝读数与基线逐字同**：空令牌／假令牌／自称来源／借证／烧后再试 均为 `approval: 原生令牌无效（缺失/已用/与本次请求不绑定）`，面板递证为 `面板来源的「允许」被服务端 API 直接拒绝（F2 第三层），可改为拒绝或查看完整参数`，重放为 `approval: correlation_id 无对应待审批项`，`宿主允许=<nil>`、`落盘=true/false` 亦同。

结论：相等关系**不是承重墙**——拆开只换地址标签，不触及任何断言。

## ② 改法（一句话）

`startChildWrite197` 增加 `corrID` 形参，`ToolRequest` 改逐字 `TaskID: taskID, CorrelationID: corrID,`（新 `:116`）；两处调用点各给一枚不同 corr：`childID+"-corr-1"`（新 `:396`）与 `childID+"-corr-2"`（新 `:463`）。形状照本仓 `cmd/wisp/run_mode101_test.go:196` 逐字 `TaskID: task, CorrelationID: corr, CallID: fmt.Sprintf("t101-call-%d", n),` 的 `t101-corr-%d` 造法（同函数 `task, corr := fmt.Sprintf("t101-task-%d", n), fmt.Sprintf("t101-corr-%d", n)`），⛔ 不是 `taskID+"-2"` 那种拼法。

diff：`1 file changed, 12 insertions(+), 5 deletions(-)`；5 枚删除行逐字＝注释首行、函数签名、`TaskID: taskID, CorrelationID: taskID,`、两处调用点；**零断言行**（对 diff 的增删行 grep `t.Fatal|t.Error|t.Skip|if ` 零命中）。

## ③ 测试读数

- 命令（改前/改后同一条）：`cd cmd/wisp && PATH="$PWD/../../third_party/sherpa-onnx:$PWD/../../third_party/onnxruntime:$PATH" go test . -run 'Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands' -v`
- 改前：`rc=0`，`=== RUN` 1 条，`--- PASS: Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands (11.23s)`，`ok github.com/CarlosShao/wisp/cmd/wisp 11.290s`；Logf：`child=97e8a90f-ee49-4ef1-af3a-c9a8fccf0472 corr1=97e8a90f-ee49-4ef1-af3a-c9a8fccf0472 corr2=97e8a90f-ee49-4ef1-af3a-c9a8fccf0472 | ...`
- 改后：`rc=0`，`=== RUN` 1 条，`--- PASS: Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands (11.21s)`，`ok github.com/CarlosShao/wisp/cmd/wisp 11.274s`；Logf：`child=5bf235a1-030e-4ab3-ad91-03383f8e7614 corr1=5bf235a1-030e-4ab3-ad91-03383f8e7614-corr-1 corr2=5bf235a1-030e-4ab3-ad91-03383f8e7614-corr-2 | ...`
- 注：`third_party/onnxruntime` 目录本机不存在，PATH 里那一截空转；本腿实读数如上（1 条 `=== RUN`、rc=0）。

## ④ 没做的（具名）

- **未动任何产码**。顺带具名报回一处相关现量（⛔ 未碰）：`internal/agent/loop.go:654` 逐字 `TaskID: taskID, CorrelationID: taskID,`，以及 `internal/tools/bridge.go:269-270` 的空值回退——产码今天的默认仍是两枚同值；本格写面只是测试载具，产码是否要跟进是另一格。
- 未放宽/删改任何断言、未加 `t.Skip`、未翻票 259 任一框、未加 `-done`、未改台账、未 push、未 `add -A`／`-a`／`--amend`。
- 未碰三枚冻结件（`internal/panel/tokens_fourway_test.go`、`internal/panel/l2_grant_boundary_test.go`、`internal/perm/ticket90_persist_test.go`）。
- 未读、未碰另一条腿 `card-proof-prep-1` 的件；`frontend/**`／`design/**` 不读不引。
- 未跑整包（`cmd/wisp` 有在册红）；未种突变（「拆开后有没有新牙」归非实现者验收腿）。

## ⑤ commit

- 起手锚：`cbe5cfd131f85c5fee91ab9385791322c8d351c8`（`00-anchor.md`）
- 本读数件与测试改动**同一笔**落地（紧接本件的那枚 commit；hash 由编排者在 `git log` 现取）。
