# 165-r3 读数件（门铃仪器·票 165 AC#6 形甲·落地腿自证，⛔ 裁决）

原始逐发日志＝本目录 `logs/`（13 枚，逐枚非零字节，`wc -c` 最短 179／最长 6610）。
本件只放读数与尺，判语在 `docs/evidence/s1/165-doorbell-instrument-r3.md`。

## 0. 锚与漂移

- 起手现量：`git rev-parse HEAD`＝`c9ba52546cee72cd37e4b1c0be01f67405b50b21`，`date`＝`2026-10-16 12:13:21 +0800`。
- ★漂移已具名：`c9ba5254..4f247070^` 之间落进三枚别腿提交（`470ae1e5d` 270-r2 派单存档／`c2e1c07ab` 326-v1 加 `internal/memory/dao_session_326v1_test.go`／`4366c1ca1` 220-v1 交件）。我这两枚的父＝`4366c1ca1`。⇒ 导出树 t1 用**起手锚 c9ba5254**（改前／改后同一枚台面），提交后的复跑用**新锚 f7d6dbc9**（t2），两棵树的 `.go` 枚数差一枚＝326-v1 那枚新测试件（744 → 746，746＝744＋别腿 1＋本腿 1）。
- 本腿提交：`4f2470703`（夹具）→ `f7d6dbc98`（用例）。⛔ push。

## 1. 台面

- 尺＝`git archive --format=tar c9ba5254 go.mod go.sum internal cmd tools scripts frontend | tar -x -C /d/tmp/165r3/t1`（★七件套带 `frontend`，⛔ 共享工作树跑改前那一发）。
- 证非空＝`find internal cmd tools scripts -name '*.go' | wc -l`＝**744**；开跑前门＝`go build ./internal/agent/` **rc=0**。
- `GOTMPDIR`＝`C:\Users\swq\AppData\Local\Temp\wisp165r3`（先 `mkdir -p`，原生反斜杠形）。
- 提交后台面＝`/d/tmp/165r3/t2`（`git archive f7d6dbc9`，同七件套），`.go`＝**746**、`go build ./internal/agent/` rc=0。
- 每一发都带 `-timeout 30m`（★具名：本腿所有 `go test` 一律 `-timeout 30m -count=1`；⛔ 用过默认 10m）；单跑一律 `-run '^TestPlanPresentCallBooksTheToolCallRow165r3$'`；⛔ 跑过 `./cmd/wisp` 任何形、⛔ 跑过 `go test ./...`。

## 2. 那对颜色（七发逐枚，红句逐字）

| 发 | 台件动作 | rc | 读数 |
|---|---|---|---|
| **改前 A** | t1 里只放用例、`mv` 走夹具 | **1** | `plan_present_doorbell_165r3_test.go:125: load golden plan-present: golden: read D:\tmp\165r3\t1\internal\agent\testdata\golden\plan-present.sse: open ...: The system cannot find the file specified.` ⇒ 用例不存在于 c9ba5254，夹具也不存在，红在装载那一步（`logs/pre-a-nofixture.log` 首行那句 `ls` 打在 `mv` 之前，不是"件在"的证据） |
| **改前 B** | t1 里两枚都在，只把 `planPresentFixture165r3` 改成既有的 `tool-then-text` | **1** | `plan_present_doorbell_165r3_test.go:138: doorbell sink reading = 0 rows for "plan.present" over 1 tool_call rows: the replay never reached internal/agent/loop.go:637-640` ⇒ **同一把尺、同一枚台面，sink 对 `plan.present` 的读数＝0**（那 1 行是 echo 的）；这一发就是"零枚用例打到那两行"的颜色 |
| **改后** | 还原夹具名，其余一字未动 | **0** | `ok github.com/CarlosShao/wisp/internal/agent 0.284s`；`-v` 那一发＝`--- PASS: TestPlanPresentCallBooksTheToolCallRow165r3 (0.13s)` |
| **去牙** | `sed -i '12d'`＝摘掉夹具第 12 行（唯一一行 `tool_calls` delta，28→27 行，`"name":"plan.present"` 命中 1→0） | **1** | `plan_present_doorbell_165r3_test.go:138: doorbell sink reading = 0 rows for "plan.present" over 0 tool_call rows: ...` ⇒ 红在**判定那一行**，读数从 1 掉到 0 |
| **还原** | 从仓外 pristine 副本 `cp` 回 | **0** | 两枚 sha256 逐位回到基线：`7d26fa889b59fc3a2b39ff18883ee331f8ff1c06ddb45a5a23bc6820af9079be`（.sse）／`00c32049994c441c8927066bc4d4099c1b2977fc20ee9e602bb98f3670d37b75`（.go）；复跑绿 rc=0 |
| **必答题 M-A** | `n := sinkCount165r3(...)` 改成 `n := 1`（常量非零），**牙完好** | **0** | 绿。★这一发单独看**不携带信息**：牙在的时候计数本来就是 1，把它换成常量当然照绿 ⇒ ⛔ 拿它当"有牙"的证据 |
| **必答题 M-B** | 同一处常量计数 ＋ **摘掉第 12 行** | **1** | `plan_present_doorbell_165r3_test.go:143: no tool_call row was booked for "plan.present" (rows read back = 0): the call never reached internal/agent/loop.go:637-640` ⇒ **用例不是恒真**：把 sink 的计数喂成常量非零之后，去牙那一发仍然红，红在独立复读那一步（`rowOfTool165r3` 自己再遍历一遍行，不复用 `n`） |
| **反控 M-C** | 名册里镜像的档位 `RiskL2` 改成 `RiskL0`（计数是真读的、牙完好） | **1** | 五行同红：`:157 row.risk_level = "L0", want "L2"`／`:184 row.decision = "allow", want "reject"`／`:191 row.error_class = "tool", want "permission_denied"`／`:199 tool log entry = ... Rejected:false`／`:208 executed tool calls = 1, want 0` ⇒ 行上的字段是真读回来的，⛔ 夹具里抄的 |

★对编排者那道必答题的答：**不是恒真**。凭据＝M-B（计数被喂成常量非零、牙摘掉，仍红在 :143）＋去牙那一发（计数真读、红在 :138）。诚实补一句射程边界：**M-A 那一发（只喂常量、牙完好）是绿的**，所以"计数那一枚门"自身对"牙在不在"确实不敏感——它靠的是 `sinkCount165r3` 真去遍历了 store 的行；把用例钉住的另有 `rowOfTool165r3` 的独立复读（:143）与行上四枚字段。

## 3. 门禁四数（逐枚落 rc）

| 门 | 尺 | rc | 读数 |
|---|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe`（v0.12.0，go1.27.1）`-l`／`-d` 本腿唯一 .go | **0 / 0** | `-l` 零命中、`-d` 零差异（PATH 里那一发没跑，直接用真身，⛔ rc=127 那形） |
| d22scan | `sh scripts/d22scan.sh`（官方包装件＝正控 `runtests.sh -C tools/d22scan ./...` ＋ 正扫工作树） | **0** | 正控 `top-level: PASS=35 FAIL=0 SKIP=0`；正扫 `clean - no D22 ban violations`，ban #8 `internal/` 射程 **577** 枚 Go 文件（含注释与 `_test.go`）、`cmd/` 139 枚 ⇒ 本腿新件在被扫的那一族里 |
| d22scan（导出树那一发，⛔ 计入判据） | `go run . -root /d/tmp/165r3/t1` 于 `tools/d22scan/` | 打印 `exit status 2` | **台面缺件，不是红**：七件套里没有 `design/`，扫描器自己那句 `d22scan: ban #8 scope design/ examined 0 files ... never leave a scope pretending to scan`＝**GUARD 射程拒绝**（`logs/d22scan.log` 末行）。⇒ 四数那一数以工作树包装件 rc=0 为准 |
| 长度门 | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`；分母 14242 枚跟踪路径、超预算 60／在册 60／未在册 **0** |
| go vet | `go vet ./internal/agent/`（t1） | **0** | 零行输出 |

## 4. 整包与用例体量

- 整包 `go test ./internal/agent/ -timeout 30m -count=1`：t1（提交前）`ok 2.856s` **rc=0**；t2（提交后）`ok 3.108s` **rc=0**。⇒ ⛔ 无 `panic: test timed out`、⛔ 无 `no test files`／`build failed`／`setup failed`／`0xc0000135` 那几形。
- 体量（单位具名）：本腿新增 **1 枚顶层用例**（`TestPlanPresentCallBooksTheToolCallRow165r3`）、**0 枚子测试**（`t.Run(` 命中 0，尺＝`grep -c 't.Run('`）、**0 枚 `t.Skip`**（⛔ 加过）、**1 枚夹具**（28 行／1684 字节，UTF-8 LF，`grep -c $'\r'`＝0）；测试件 221 行／9378 字节。
- 生产调用点枚数（⛔ 本腿改产码）：`internal/agent/loop.go:637-640` 那两行仍只有 `executeCalls` 一处写入（尺＝`grep -n "j.startCall" internal/agent/*.go`，产码 3 处＝`loop.go:640`/`:879`/`:911`，后两处是截断与未闭合那两条失败支，⛔ 本腿新写的）。
- `llm.Request.ToolChoice` 产码写入枚数**照旧 0**（派单硬约束 2，本腿一枚没加）。尺逐字＝`git grep -n ToolChoice f7d6dbc9 -- internal/agent cmd/wisp` ⇒ 现量 **0 行**（该射程内产码 0／测试 0；字段本体在射程外，声明于 `internal/llm/provider.go:71`，与普查腿 `logs/toolchoice-census.txt` 那枚"产码 0 枚"同向）。`tools/mockllm/**` 的按位合成本腿一字未动。

## 5. 落件后的分母复跑（尺逐字抄普查腿 §问 B）

锚 `f7d6dbc9`，`git grep -F -l <name> HEAD -- '*.sse' '*/testdata/*' 'tools/mockllm/**'`，单位＝枚文件：

- `internal/agent/testdata/golden/*.sse` 枚数：普查量得 **11** ⇒ 本腿落件后现量 **12**。
- 11 枚真注册名逐枚命中：`plan.present` **1**（普查腿为 0），其余 **10 枚仍为 0**（`fs.read`/`fs.list`/`fs.edit`/`fs.write`/`fs.trash`/`fs.move`/`fs.delete`/`task.spawn`/`task.output`/`task.cancel`）⇒ "整族零读数"那一枚结论**只被本腿挪动了一格**，缺口没关掉。
- 正控（同一把尺、同射程）：`echo 17`／`get_weather 8`／`sleep 3`／`get_time 4`／`tool_calls 16`（普查腿 15 ⇒ ＋1 就是本腿这枚夹具，尺不瞎也不漏）。
- 原始逐行读数＝`logs/post-fixture-denominators.txt`。

## 6. 本腿没做的两枚断口（⛔ 不归我，具名交回）

- 断口② `cmd/wisp/plan_assembled_165r2_test.go:123` 只断到名册出口，跳 6/7/8（`Inject` → `prompt.go:251` → 各适配器 `request.go`）仍零用例——本腿那发请求体断言只覆盖了 round 1 的 `tools` 数组里有没有这枚门（`mustContain ... "name":"plan.present"`），⛔ 等于"装配面断言"；写面 `cmd/wisp/**` 此刻归 `220-v1`/`220-r2`，本腿一枚字节没碰。
- 断口③ `llm.Request.ToolChoice` 产码 0 枚写入 ⇒ `tools/mockllm/chat.go:172-185` 那枚"必发工具调用"的合成仍从装配面按不到；本腿按派单硬约束⛔ 加写入（要扩到 `internal/llm` 三家 `request.go`）、⛔ 改合成器的按位形状。这一格缺的是**世界**（产码写入＋合成器按名选），⛔ 一枚包内用例能补。

## 7. 终态名册与最后一发（补笔）

- 本腿八枚提交（逐枚 `git show --name-only` 的路径名册共 **23 枚**，全在四处写面之内，越界尺命中 **0**）：`4f2470703` 夹具／`f7d6dbc98` 用例／`d498a73c` 读数件／`78fe125f` logs 15 枚／`3e4ca0e3` 证据件／`c3f55651` 票尾一条 bullet／`f6c76b91` 三枚正文源台件／本笔（读数件＋证据件各一枚补笔）。⛔ push（本地领先 `origin/dev`，推送归编排者）。
- 终态锚＝`f6c76b91897ead14e2717f4209c9f5422a9098eb`；其间别腿又落了两枚（`83a6f15c`／`7f31bd2b`，220-v1），⛔ 与本腿写面零交集。
- 终态那一发（★第三棵导出树 `git archive f6c76b91` 七件套，原始日志 `logs/final-head-rerun.log`）：两枚字节 sha256 仍逐位等于基线（`00c32049994c441c...`／`7d26fa889b59fc3a...`）；`go build ./internal/agent/` rc=0；单跑 `-run '^TestPlanPresentCallBooksTheToolCallRow165r3$' -timeout 30m -count=1`＝`--- PASS ... (0.13s)` rc=0；整包 `./internal/agent -timeout 30m -count=1`＝`ok 2.250s` rc=0；`go vet ./internal/agent/` rc=0。
- 票框尺交回（⛔ 本腿一枚没翻，翻法归编排者）：`grep -c '^- \[x\]'`＝**5**（AC#1–AC#5，落件前后同数）、`grep -c '^- \[ \]'`＝**1**（AC#6 仍空）；票名与两行 `Status:` 一字未动；票面 diff＝**2 加／0 删**（纯追加一条 bullet）。
