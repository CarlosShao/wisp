# 165 门铃仪器落地件（腿 165-r3·票 165 AC#6 形甲·实现者自证，⛔ 裁决）

- 角色＝落地腿（写码）。判语归非实现者；本件只交字节、颜色与尺。
- 起手锚 `c9ba52546cee72cd37e4b1c0be01f67405b50b21`（12:13:21 +0800 现量），落件父锚 `4366c1ca1`（中间漂进别腿三枚：`470ae1e5d`／`c2e1c07ab`／`4366c1ca1`，逐枚见读数件 §0）。⛔ push。
- 原始逐发日志＝`.scratch/wisp/probes/165/r3/logs/`（15 枚）；数与 rc 全在 `.scratch/wisp/probes/165/r3/00-readings.md`。

## 1. 落了什么（就两枚字节）

1. `internal/agent/testdata/golden/plan-present.sse`——录下的 SSE 两发：response 1 ＝一枚 `plan.present` 的 `tool_call`（`id=call_pp1`，arguments 是 `summary` ＋ 两步清单），response 2 ＝循环把那一跳的答灌回历史之后模型给的文本。形状照同目录 `tool-then-text.sse`，只把 `function.name` 换成本仓名册里那枚真工具（`docs/PLAN.md:2567`）。**牙只有一枚**：第 12 行是唯一一行 `tool_calls` delta，摘掉它这一发就没有调用可组装。
2. `internal/agent/plan_present_doorbell_165r3_test.go`——`package agent`（包内形状，⛔ 新建外测包），**1 枚顶层用例／0 枚子测试／0 枚 `t.Skip`**。真装配：golden Replayer 交字节给 httptest 上的真 openai-chat 适配器（C5 接缝）→ 真 `Loop` → `Journal` 是真 `memory.Store`。跑完只读一件事：`internal/agent/loop.go:637-640` 那两行往 `tool_call` 表里记没记出行。

判据形状＝派单写死那一形：喂进这枚夹具 ⇒ sink 读数必须非零（`:138`，`n == 0` 即 `t.Fatalf`）；摘掉夹具那一行 ⇒ 用例必须红。**已按这两发各跑一遍**（下 §3）。

## 2. 两枚硬约束怎么接住的

- **导入环**：`internal/tools/bridge.go:13` 已 import `internal/agent`，而 `internal/agent` 全树零枚 `package agent_test` ⇒ 真桥进不去。本腿⛔ 试图 import `internal/tools` 或 `cmd/*`；名册那枚条目是**镜像**的（`Resident true`／`Declared L2`／参数面同 C1 schema），出处逐枚具名在文件头注释（`docs/PLAN.md:2567`＋`internal/tools/plan_present.go:72`、`:123-132`）。桥那一半＝`cmd/wisp` 的面，本波不归我。
- **`ToolChoice`／mockllm**：`llm.Request.ToolChoice` 的产码写入枚数**照旧 0**（现量尺 `git grep -n ToolChoice f7d6dbc9 -- internal/agent cmd/wisp` ⇒ 该射程 0 行）；`tools/mockllm/chat.go:172-185` 的按位合成一字未动。这枚夹具让"这一枚工具被选中"靠的是 **SSE 里的字节**，⛔ 靠改合成器。
- 工具体**一枚没跑**：这台宿主没接审批通道（`Options.AdmitTask == nil`），`decideRisk` 对声明 L2 只能拒（`internal/agent/loop.go:820-823`），拒之前那两行已先把行记进表。用例把这件事钉死（`h.echo().CallCount()` 断言 0，反控 M-C 那一发证明它真的有牙）。⇒ ⛔ "拿 mock 顶真的"：本用例⛔ 假装 `plan.present` 跑成了，它只裁"那一跳被走过、行落进表"。

## 3. 那对颜色（逐字红句在 `logs/`，此处给数）

| 发 | rc | 红在哪一行／什么数 |
|---|---|---|
| 改前 A（c9ba5254 导出树，只放用例、无夹具） | 1 | `:125 load golden plan-present: ... cannot find the file specified` |
| 改前 B（同台面、同用例体，夹具名换成既有 `tool-then-text`） | 1 | `:138 doorbell sink reading = 0 rows for "plan.present" over 1 tool_call rows` ⇒ 同一把尺对 `plan.present` 读 0（那 1 行是 echo 的） |
| 改后（还原夹具名） | 0 | `--- PASS: TestPlanPresentCallBooksTheToolCallRow165r3 (0.13s)` |
| **去牙**（`sed -i '12d'`，28→27 行） | 1 | `:138 ... = 0 rows for "plan.present" over 0 tool_call rows` ⇒ 红在判定那一行 |
| 还原 | 0 | 两枚 sha256 逐位回基线（`.sse 7d26fa889b59...`／`.go 00c32049994...`），复跑绿 |

★编排者那道必答题（"这枚新用例是不是恒真"）的答：**不是恒真**，但要把两头发分开看——

- **M-A**：`n := sinkCount165r3(rows, ...)` 改成常量 `n := 1`，**牙完好** ⇒ **绿 rc=0**。这一发单独**不携带信息**（牙在时计数本就是 1），⛔ 拿它当凭据。诚实的射程边界就在这：**"非零"那一枚门自身对牙不敏感**，它敏感的是 `sinkCount165r3` 真去遍历了 store 的行。
- **M-B**（决定性那一发）：同一处常量计数 ＋ **摘掉牙** ⇒ **红 rc=1**，红在 `:143 no tool_call row was booked for "plan.present" (rows read back = 0)`——`rowOfTool165r3` 独立又遍历了一遍行、不复用那个计数。⇒ 计数被喂成假值之后用例仍抓得住缺牙，**"我对这件事不敏感"这一句本腿⛔ 写**（不敏感的是 `n` 那一枚门，⛔ 是整枚用例）。
- **M-C**（反控，超出必答题范围的一发加测）：把镜像名册的 `RiskL2` 改成 `RiskL0` ⇒ 5 行同红（`:157 risk_level`／`:184 decision`／`:191 error_class`／`:199 tool log Rejected:false`／`:208 executed tool calls = 1, want 0`）⇒ 行上字段是真读回来的，⛔ 夹具里抄的。

## 4. 门禁四数（逐枚 rc）与回归

`gofumpt -l`/`-d`（真身 v0.12.0）**rc=0/0**（零命中零差异）｜`sh scripts/d22scan.sh` **rc=0**（正控 PASS=35 FAIL=0 SKIP=0；正扫 clean，ban #8 在 `internal/` 射程 577 枚 Go 文件含注释与 `_test.go`）｜`sh scripts/check-path-length-budget.sh` **rc=0**（VERDICT GREEN，14242 枚跟踪路径、超预算 60／在册 60／未在册 0）｜`go vet ./internal/agent/` **rc=0**。
另具名一枚**⛔ 计入判据**的红：`go run . -root <导出树>` 在七件套树上打印 `ban #8 scope design/ examined 0 files` 后 `exit status 2`＝台面缺 `design/` 的 GUARD 射程拒绝，⛔ 是红。
整包 `go test ./internal/agent/ -timeout 30m -count=1`：提交前 **ok 2.856s rc=0**、提交后（`git archive f7d6dbc9` 第二棵树）**ok 3.108s rc=0**。所有 `go test` 一律显式 `-timeout 30m`；⛔ 跑过 `./cmd/wisp` 任何形、⛔ 跑过 `go test ./...`。

## 5. 写面自证（逐枚提交、⛔ 目录形）

`4f2470703`（夹具 1 枚）→ `f7d6dbc98`（用例 1 枚）→ `d498a73c`（读数件）→ `78fe125f`（logs 15 枚）。逐枚 `git show --name-only` 只含四处写面之内：`internal/agent/**`、`.scratch/wisp/probes/165/r3/**`、`docs/evidence/s1/**`、票 165 末尾一条 bullet。⛔ `cmd/wisp/**`、⛔ `internal/panel/**`、⛔ `internal/llm/**`、⛔ `tools/mockllm/**`、⛔ `internal/tools/**`、⛔ 既有 golden 件、⛔ 阈值/`allowlist.txt`/`PLAN.md`/`docs/specs/**`/三枚冻结件。⛔ 翻任何 AC 框、⛔ 改票名/`Status:`、⛔ 加 `-done`、⛔ `add -A`／amend／reset／rebase／stash／checkout ./clean。

## 6. 本件量到／⛔ 量到的（先把尺说清，别让下一枚腿误读）

- **量到的**：这枚门的**线路**——同一发运行里它挂上了 wire（`bodies[0]` 含 `"name":"plan.present"`），而一旦 SSE 字节里真出现一枚对它的调用，`loop.go:637-640` 就把名字与参数落成 `tool_call` 行；不按门则表里一行都没有（去牙／改前 B 两发）。
- **⛔ 量到的**：**"生产里今天会不会有人按"**。那需要真 provider 凭据＋真窗（普查腿的丙形），本腿⛔ 有凭据，已具名成资源需求。本件是常驻线路仪器，⛔ 是产品结论。
- ★**一处措辞要请编排者复核**（⛔ 本腿改派单/票面一字）：派单第 9 行写"那两行 sink 今天零枚用例打到它"。现量更精确的形状是——那两行**走过得通**：落件前同包有 **6 枚**用例通过 `withJournal(store)` 把**完整**调用打到 `loop.go:640` 那一行并真写进过表（尺＝`git grep -h 'withJournal(store)' c9ba5254 -- internal/agent | sed ... | sort | uniq -c` ⇒ 夹具分布 `tool-then-text` 3／`slow-tool` 2／`six-tools` 1，全是 echo/sleep 那族假工具；另 **2 枚**走的是截断与未闭合那两条失败支 `loop.go:879`/`:911`＝`max-tokens-toolcall` 1／`disconnect-mid-toolcall` 1；`withJournal(store)` 出现枚数＝落件前 **8**、含本腿这枚后 **9**）。今天缺的那一格是**"表里读得出的名字是名册里那枚真工具"**：`plan.present` 在夹具面与表里都是 0（普查腿那把尺复跑，落件后 0→1）。AC#6 与普查件 §断口① 的原文正是后者 ⇒ 本格按"这枚门走过去"裁，⛔ 按"这两行从没被执行过"裁。

## 7. 没裁干净的残余（具名，⛔ 本腿不动手）

- 断口② `cmd/wisp/plan_assembled_165r2_test.go:123` 只断到名册出口，跳 6/7/8 零用例 ⇒ 该由 `cmd/wisp` 那一格的写腿补（本波＝`220-v1`/`220-r2`），形＝普查件 §问 C 的乙-1（`adaptertest.StartMockllm` ＋ `Mockllm.LastRequest(t,"chat")` 断 `sent["tools"]` 含 `plan.present`，正控＝摘 `cmd/wisp/run.go:826-831`）。本腿写面⛔ 覆盖那一面。
- 断口③ `llm.Request.ToolChoice` 产码 0 枚写入＋mockllm 按位合成 ⇒ 缺的是**世界**（要扩到 `internal/llm` 三家 `request.go` 与 `tools/mockllm`），⛔ 一枚包内用例补得了；派单硬约束已让本腿停手上交，请裁归属。
- 其余 10 枚真注册名在夹具面仍逐枚 0 命中（现量在 `logs/post-fixture-denominators.txt`）⇒ "整族零读数"没被本腿关掉，本腿只挪走一格。

## 8. 复现配方（可逐字重跑）

```sh
git archive --format=tar c9ba5254 go.mod go.sum internal cmd tools scripts frontend | tar -x -C <仓外树>   # 改前台面
# 改前 A：只 cp 用例、mv 走夹具            => rc=1，红在 :125
# 改前 B：cp 两枚后 sed 夹具名为 tool-then-text => rc=1，红在 :138，读数 0 over 1
# 改后：还原夹具名                          => rc=0
sed -i '12d' internal/agent/testdata/golden/plan-present.sse                       # 去牙 => rc=1，红在 :138
sed -i 's/\tn := sinkCount165r3(.*/\tn := 1/' internal/agent/plan_present_doorbell_165r3_test.go  # 必答题计数喂常量
#   牙完好 => rc=0（不携带信息）；再叠上去牙 => rc=1，红在 :143（决定性那一发）
go test ./internal/agent/ -run '^TestPlanPresentCallBooksTheToolCallRow165r3$' -timeout 30m -count=1
go test ./internal/agent/ -timeout 30m -count=1
```
