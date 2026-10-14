# 236-r4 — 票 236 AC#3 落地：钉「键的来源（宿主铸 vs 模型给）」这一维（写腿／实现者）

> 判语归后续非实现者腿；本件只报读数与改动，⛔ 不翻任何 AC 框、不改 `Status:`／票名。
> 起手锚 `b8f883a3 10-14 10:33`（派单发送时 `a5b2938c`）；交件窗口内 `236-r5` 两笔落库（`0c7f3c5d`／`faa8950d`），零文件重叠。
> 落点＝**新建** `internal/tools/taskid_source_236r4_teeth_test.go`（275 行，sha256 `47fb441dc1597524db490007c859307a58c691aa04346998b6482194af2ae2f9`；**零产码改动／零导出名／零放宽既有断言**）。
> 钉名＝`Test236R4RunningTaskIDAndResultTaskIDAreOneHostMintedKey`（唯一的常驻用例；本文件没有第二枚断言用例）。
> `git status --porcelain -- internal cmd scripts .github` 起手＝0 行；交件＝只列本腿新件与 `315-a1` 的 `cmd/wisp/zz315a1probe_test.go`（非本腿，未读未碰）。

## §1 落点、两处键、以及「射程含 internal/agent」怎么落

**两处（行号 10-14 现读；票面 `run.go:892` 与 a4 的 `:1099` 均已漂，引用前重跑）**：

| 处 | file:line（现读） | 逐字 |
|---|---|---|
| 铸点 | `internal/agent/loop.go:329` | `id := newTaskID()`（`:330` 填 `RunningTask.ID`，`:333` 同一枚交给 `run`，`:382` 落进 `Result.TaskID`；形状由 `:1153 func newTaskID` 决定） |
| 用点 | `cmd/wisp/run.go:1105` / `:1110` | `rt.tasks.MarkRoot(bg.ID, task)` / `Backfill(bg.Root().Ctx, res.TaskID, res.Text)`；`:1124`／`:1135`／`:1145`（task_log、FinishTaskLog、状态行）全用 `res.TaskID` |

**载具**＝`pointer_183_cli_seam_test.go` 两枚的形状：不立 mock，真 loop（`agent.New`）＋真 C5（openai-chat 适配器回放冻结夹具 `internal/agent/testdata/golden/tool-then-text.sse`，SPEC-05 §9「C5 不开函数 mock」）＋环路自带 echo 工具面（声明 L0，夹具那次真派发）＋真名册与真 spill 层，并把 `run.go:1105`／`:1110` **这两个消费调用原样复做**，让「两处相等」有可观察后果（分叉＝名册多一行、以铸来的键读不回正文）。

**「射程含 internal/agent」的落法（防 a4 点名的那一漏）**：这枚钉写在 `internal/tools`（在 AC#3 那两目录口径**之内**，不会再被过滤掉），但断言的是**跨包不变式**——它驱动 agent 包的真铸造路径，红/绿由 `internal/agent/loop.go` 的字节决定。证据＝§3.3 两发突变全打在 `internal/agent`，这枚钉红；反方向 §3.4：仓里唯一那枚射程在 `internal/agent` 的形状尺（`compress_trace_test.go:526`）**在同样两发下都绿**（它的 harness 走 `loop.Run`＝同步铸点 `:340`，够不到组合根用的 `RunAsync`＝`:329`）。⇒ 按包开尺的两个方向各自都看不见的那一枚，现在这枚钉同时补上。

## §2 三段双向怎么抄（先例＝`internal/session` minted-shape 锁）

| 段 | 先例 | 本件落法 |
|---|---|---|
| 锁/谓词 | `session.go:67 Valid` ＋ `:94 Mint` | `r236r4MintedShape`（36 字符、8/13/18/23 位连字符、其余小写 hex）；**测试本地**，不新造工具面、不加导出名 |
| 便宜向 | `grants_test.go:152`（台件字面量在铸形之外） | `r236r4ModelWireIDs`（冻结 C5 seam 真携带的模型侧 id，10-14 现数）＋`r236r4FixtureTaskIDs`（本仓台件真用的手写 task id）都必须过不了锁；且本次真铸的 id ≠ 名单里每一枚 |
| 贵向 | `ticket224_assembly_test.go:433`（真装配铸的 id 必须过锁） | 真 loop 跑出来的 `bg.ID` 过锁、`res.TaskID == bg.ID`、接缝上 `req.TaskID == bg.ID` 且 `!= req.CallID`、名册行恰 1 行且以铸来的键读得回正文 |

## §3 读数（全部本腿现跑；大日志留盘不入库）

### 3.1 基线名册（正控前提；口径＝顶层 `^--- PASS/FAIL/SKIP`，缩进子测试不计，票面 §六）

| 包 | 命令 | rc | PASS／FAIL／SKIP |
|---|---|---|---|
| `internal/tools` | `go test ./internal/tools/ -count=1 -v` | 0 | **220／0／0**（17.201s；log `baseline-tools-v.txt` 153,918 B） |
| `internal/agent` | `go test ./internal/agent/ -count=1 -v` | 0 | **85／0／0**（4.544s；log `baseline-agent-v.txt` 34,805 B） |

### 3.2 正控（不突变 ⇒ 绿）

- 新钉单独跑 `-run 'Test236R4' -v`：**rc=0／PASS 0.01s**。
- 两发 overlay 窗口**之后**再跑一次：**rc=0／PASS 0.01s**（证明跑完突变后树仍绿，红只可能来自突变）。

### 3.3 两发突变（每发：仓外备份＋确证 diff＋红句逐字＋sha256）

- 仓外件：原件备份 `C:\tmp\wisp236r4\backup\loop.go.orig`（41,089 B，sha256 `cb7f9da7…`＝盘上＝HEAD blob 三方同串）；M1 副本 41,087 B（`971436d0…`）／M2 副本 41,098 B（`d279d117…`）；overlay json 各 122 B。**全程 `-overlay` 合成，零原地编辑。**

**M1（把铸点的来源换成模型可给的值）** 窗口 `10:40:58` → `10:41:29`
```
--- backup/loop.go.orig   +++ mut/m1-source-swap/loop.go   @@ -326,7 +326,7 @@
-	id := newTaskID()          (loop.go:329)
+	id := "call_e1"            ← 冻结夹具里模型给的那枚 call id，逐字
```
`go test ./internal/tools/ -count=1 -run 'Test236R4' -v -overlay …m1…` ⇒ **rc=1，3 枚具名红**（逐字节选）：
`taskid_source_236r4_teeth_test.go:189: RunningTask.ID = "call_e1" 不是宿主铸的形状…internal/agent/loop.go:329 的 id := newTaskID() 填的就是这一格，形状由 :1153 决定。把来源换成模型可给的值时这一句必须红`
`:228: 接缝上 TaskID 与 CallID 同值（"call_e1"）：宿主铸的键与模型给的键合流了，「键的来源」这一维当场失效（loop.go:676 的两格必须各是各的）`
`:267: 铸出来的 id 等于模型侧字面量 "call_e1"：模型给的值冒充了宿主铸的键`

**M2（把铸点与用点的关联拆开）** 窗口 `10:41:29` → `10:41:39`
```
--- backup/loop.go.orig   +++ mut/m2-broken-link/loop.go   @@ -330,7 +330,7 @@
-	t.result = l.run(c, id, input)          (loop.go:333)
+	t.result = l.run(c, newTaskID(), input) ← 同一个任务里再铸一枚
```
`… -overlay …m2…` ⇒ **rc=1，4 枚具名红**（逐字节选）：
`:203: Result.TaskID = "bbe45e64-…" 而 RunningTask.ID = "c79fdda3-…"：这两处必须是同一枚宿主铸的键。铸造与传递在 internal/agent/loop.go:329（id := newTaskID()）与 :333（同一个 id 交给 run），结果在 :382 落进 Result.TaskID；消费在 cmd/wisp/run.go:1105（MarkRoot(bg.ID)）与 :1110/:1124/:1135/:1145（回填、task_log、状态行全部用 res.TaskID）。两处分叉即「名册登记的键」和「输出写进的键」不是同一枚`
`:223: 工具接缝上 TaskID = "bbe45e64-…"，不是环路铸的那枚 RunningTask.ID = "c79fdda3-…"…`
`:247: 以铸来的键读回的正文 = "现在是 22 摄氏度，晴。" / ""，不是 Result 的 …：输出写的键与名册登的键分叉了`
`:255: 名册有 2 行, want 1：MarkRoot(bg.ID) 与 Backfill(res.TaskID) 一旦分叉就会多出一行，这正是「两处相等」的可观察后果`

**归因（两发各自咬哪一条，互不重复）**：M1 红在形状＋字面量＋接缝合流三条（`res.TaskID == bg.ID` 那一句**不红**——两处同流同一枚假值）；M2 红在相等＋接缝＋名册链四条（形状与字面量那两条**不红**——两枚都是真 uuid）。⇒ 牙不共用一副。

**还原证**：M1 后与 M2 后两次 `sha256sum internal/agent/loop.go cmd/wisp/run.go` ＝ `cb7f9da7a6d277f341d3b40caee59eecec3bea4e4faee31777b1ad47c3763189`／`d60c50e3e6af98dc5c6f6842574c3cf47c0093956063b5bd8813dc171eb4c68d`，与备份、与 `git show HEAD:` blob **三方同串**；`-overlay` 结构性不写盘 ⇒ 「还原」不是一条命令而是一个事实。

### 3.4 交叉读数：既有 `internal/agent` 形状尺在同样两发下都 PASS

`go test ./internal/agent/ -count=1 -run TestCompressionTraceCarriesTheOwningTaskID -v -overlay …（m1／m2）` ⇒ **两发都 `--- PASS`，rc=0**（log `xcheck-agent-m1.txt`／`xcheck-agent-m2.txt` 各 307 B）。机制＝它的 harness `h.run` 走 `h.loop.Run`（`harness_test.go:189`）＝同步铸点 `:340`；组合根用的 `RunAsync`（`:329`）不在它射程内。⇒ a4 那句「唯一的红是形状尺」在**定向只换 `:329`／`:333`** 时并不成立：那枚尺一枚都咬不到。⛔ 不据此判它弱——那是它的地界，本件只登记「本钉的牙是新增的」。

## §4 门禁（全绿）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./internal/tools/ ./internal/agent/ ./cmd/wisp/` | **rc=0** |
| 格式 | `gofmt -l`／`gofumpt -l`（新件） | **零输出** |
| 静态 | `go vet ./internal/tools/` | **rc=0** |
| d22scan | `sh scripts/d22scan.sh` | **rc=0**；正控 35／0／0 先绿；`clean - no D22 ban violations`；**ban #8 internal/=550**（含本新件） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | **VERDICT GREEN**；`over-budget=59 covered by roster=59 not in roster=0`（＝票面 59=59） |
| 树 | `git status --porcelain -- internal cmd scripts .github` | 只剩本腿新件＋`315-a1` 探针（非本腿） |

## §5 具名没做／天花板（⛔ 不许被读成已覆盖）

1. **覆盖不到熵失败**：`rand.Read` 直引、无注入缝，且退化支（`loop.go:1132-1137`）**格式化在 `if` 之外 ⇒ 与正常形同 36 字符形状**，本钉的形状/相等断言都分不出熵失败。⛔ 不为测试开生产后门（AGENTS §1.3）。文件头逐字写明。
2. **同步铸点 `:340` 不在本钉射程**：它今天无产码调用者（`loop.Run` 只见于测试，如 `loop_approval_test.go:170`／`harness_test.go:189`），且「两处相等」只存在于 `RunAsync` 那条链上。⛔ 未顺手加第二枚钉。
3. **`cmd/wisp/run.go` 本身没有被执行**：本钉把它 `:1105`／`:1110` 两个消费调用在进程内原样复做；run.go **内部**若自造一枚 id（票面 `:1090` 注释自己登记过的开放端），从 `internal/tools` 观察不到——与 `pointer_183` 载具同一天花板。具名。
4. **模型给 uuid 形值时形状锁是瞎的**：`00000000-0000-4000-8000-000000000000` 这类值在铸形之内（`subagent_stream_key_197_test.go:55`／`subagent_197_test.go:886` 正在用它证明「uuid 形可当 call id」）。便宜向是**枚举**不是结构；本件把它写进文件头，且**故意不造**那枚恒真的不等式断言（那是拿假绿换一条勾）。
5. **没碰**：产码语义／`allowed_dirs`／工具面／既有断言／`subagent_197_test.go`／`subagent_222_test.go`／三枚冻结件／`thresholds.go`／golden／`d22scan` allowlist／`PLAN.md`／`docs/specs/**`／`frontend/**`／`design/**`；⛔ 未 `push`。

**临时件（只建不删，留盘点名）**：`C:\tmp\wisp236r4\logs\`＝`baseline-tools-v.txt`(153,918 B)／`baseline-agent-v.txt`(34,805 B)／`m1.txt`(1,014 B)／`m2.txt`(1,782 B)／`xcheck-agent-m1.txt`(307 B)／`xcheck-agent-m2.txt`(307 B)／`posctrl-after-mut.txt`(335 B)／`d22scan.txt`(22,972 B)／`length-gate.txt`(1,457 B)＋`baseline-tools.txt`(55 B)／`baseline-agent.txt`(54 B)；`mut/`／`overlay/`／`backup/` 三目录见 §3.3。
