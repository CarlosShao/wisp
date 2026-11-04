# 370-a1 — AC#0 现读＋归因（只读腿）

- 票：`.scratch/wisp/issues/370-assembler-maxoutputtokens-is-never-assigned-so-our-ceiling-is-never-sent.md`
- 腿：`370-a1`（只读，AC#0；⛔ 未改任何产码字节）
- 日期：2026-11-04 21:2x +08
- HEAD：`32fbed84c`（工作树；`git rev-parse --short HEAD` 现跑）
- 本腿新增（仅测试）：`internal/agent/maxoutput_370_ac0_test.go`
- 结论一句话：**票面"全仓无赋值"在 `74f430d86^` 上逐字成立；但在 HEAD `32fbed84c` 上已经不成立——`74f430d86` 是 HEAD 的祖先，链已接通，且机主已批准〔甲′〕（台账 `:30628`）。本文件同时给出两棵树上的读数。**

---

## 0. 先报前提：HEAD 上票面已过期（这是本次量到的最重要的一条）

现跑：

```console
$ git log --oneline -6
32fbed84c 167轮:A1526 29收口第30枚(AC#1/2+v1PASS+改名done双名核=1)
...
$ git merge-base --is-ancestor 74f430d86 HEAD && echo YES-ANCESTOR
YES-ANCESTOR
$ git log --oneline -3 -- internal/agent/maxoutput_370_test.go
74f430d86 feat(370 AC#2 甲′): run.go 保留 chat RoleSettings 逐任务重解析→agent.Config.MaxOutputTokens→Assembler→请求 max_tokens；补3枚测试(wire级JSON body+反向不发+per-loop)全绿
$ git log --oneline -2 -- cmd/wisp/task370_ac2_test.go
7d158caf5 test(370 AC#2): 端到端撞上限→task_log.state=interrupted 测试(正控+两向反向,用probe365罐);本机probe365环境不工作(既有365测试同红),待CI验证
```

`git show --stat 74f430d86`：

```console
 cmd/wisp/run.go                      |  33 +++++++++-
 internal/agent/loop.go               |  21 ++++++-
 internal/agent/maxoutput_370_test.go | 116 +++++++++++++++++++++++++++++++++++
 3 files changed, 167 insertions(+), 3 deletions(-)
```

台账 `docs/reports/pending-and-issues.md:30628` 逐字：

> **370〔甲′〕落地**（`74f430d86`，3 文件）：run.go 保留 chat `RoleSettings` 逐任务重解析→`agent.Config.MaxOutputTokens`→`Assembler.MaxOutputTokens`→请求 `max_tokens`；补 3 枚测试……**机主批准〔甲′〕，已兑现 `max_output_tokens` 真流到请求。**

⇒ **AC#1 的"甲"已经在 HEAD 上落地并获批。** 本腿只读，不改它、不翻票框。
下面每一项都按"HEAD / 父提交 `74f430d86^`"两列报。

---

## 1. `git grep MaxOutputTokens` 逐处清点（AC#0①）

### 1.1 现跑命令（HEAD，只看受跟踪的 `.go`）

```console
$ git grep -n "MaxOutputTokens" -- "*.go"
$ git grep -n -E "\.MaxOutputTokens[[:space:]]*=" -- "*.go"
$ git grep -n -E "\.MaxOutputTokens[[:space:]]*=" 74f430d86^ -- "*.go"
```

### 1.2 产码逐处（每一处都列；`读`＝只读不写，`声明`＝字段/常量，`赋值`＝真的写）

| # | file:line | 符号 | 类别 | 说明 |
|---|---|---|---|---|
| 1 | `internal/config/schema.go:297-298` | `config.Role.MaxOutputTokens` | 声明（TOML 解码写入） | `MaxOutputTokens int \`toml:"max_output_tokens"\``；这是 `[llm.roles.*].max_output_tokens` 的落点 |
| 2 | `internal/config/schema.go:360-361` | `config.ModelSpec.MaxOutputTokens` | 声明（TOML 解码写入） | 模型目录上限，**与角色上限是两个字段** |
| 3 | `internal/llm/resolver.go:40` | `llm.RoleSettings.MaxOutputTokens` | 声明 | 角色采样设置三元组之一 |
| 4 | `internal/llm/resolver.go:205` | `llm.RoleSettings.MaxOutputTokens` | **赋值** | `MaxOutputTokens: role.MaxOutputTokens`；见 §2 |
| 5 | `internal/llm/resolver.go:142` | `Endpoint.MaxOutput`（注意不是同名） | 赋值 | `MaxOutput: spec.MaxOutputTokens`——把 ModelSpec 的目录上限写进 `Endpoint.MaxOutput` |
| 6 | `internal/llm/provider.go:74` | `llm.Request.MaxOutputTokens` | 声明 | C5 请求字段；**父提交前就存在** |
| 7 | `internal/llm/provider.go:139-140` | 读 | 读 | `EstimateTokens` 里 `if r.MaxOutputTokens > 0 { tokens += r.MaxOutputTokens }` |
| 8 | `internal/agent/prompt.go:131-132` | `Assembler.MaxOutputTokens` | 声明 | 票面点名的字段 |
| 9 | `internal/agent/prompt.go:253` | `Assembler.MaxOutputTokens` | **读** | `MaxOutputTokens: a.MaxOutputTokens`（Build 里抄进 Request） |
| 10 | `internal/agent/loop.go:133-142` | `agent.Config.MaxOutputTokens` | 声明 | 370 新增的承接字段 |
| 11 | `internal/agent/loop.go:261` | `Assembler.MaxOutputTokens` | **赋值** | `asm.MaxOutputTokens = opt.Config.MaxOutputTokens`；HEAD 上**唯一**给 Assembler 赋值的产码行 |
| 12 | `cmd/wisp/run.go:354-358` | `runtime.roleSettings llm.RoleSettings` | 声明 | 370 保留的 boot 副本 |
| 13 | `cmd/wisp/run.go:525,531` | `rt.roleSettings` | 赋值 | boot 时 `ep, roleSettings, err := res.ResolveRole(llm.RoleChat)` → `rt.roleSettings = roleSettings` |
| 14 | `cmd/wisp/run.go:1498,1503` | `roleSettings` | 赋值 | 每任务重解析；失败回退 boot 副本 |
| 15 | `cmd/wisp/run.go:1560` | `agent.Config.MaxOutputTokens` | **赋值** | `MaxOutputTokens: roleSettings.MaxOutputTokens` |
| 16 | `internal/llm/openaichat/request.go:75`（字段 `:60`） | `wireRequest.MaxTokens` | 读 | `MaxTokens: req.MaxOutputTokens`；`json:"max_tokens,omitempty"` |
| 17 | `internal/llm/anthropic/request.go:103`（+ `:107-108`） | `wireRequest.MaxTokens` | 读 | `MaxTokens: req.MaxOutputTokens`；`<=0` 时填 `defaultMaxTokens = 4096`（`adapter.go:94`） |
| 18 | `internal/llm/openairesponses/request.go:78`（字段 `:55`） | `wireRequest.MaxOutputTokens` | 读 | `json:"max_output_tokens,omitempty"` |
| 19 | `internal/memory/l1_extract.go:201` | `llm.Request.MaxOutputTokens` | **赋值** | `MaxOutputTokens: l1DefaultMaxTokens`（`=512`，`l1_extract.go:82`）——**记忆抽取路径**，不经 Assembler，不在本票链上；它说明"Request 恒 0"只对 agent 聊天链成立 |

注释性命中（不构成赋值）：`cmd/wisp/run.go:355`、`internal/agent/loop.go:133`、`internal/agent/prompt.go:131`、`internal/llm/anthropic/adapter.go:24`、`internal/llm/openairesponses/adapter.go:19`。

测试命中（全部是测试常量/断言，不构成产码赋值）：
`internal/agent/maxoutput_370_test.go`（多处）、`internal/agent/instructions_test.go:211`、`internal/config/boundary_test.go:105,126`、`internal/llm/llm_test.go:646`、`cmd/wisp/task370_ac2_test.go:13`（注释）、`cmd/wisp/tasklogstate_365_test.go:286`（注释）。
`.scratch/wisp/probes/**` 下另有历史探针副本（`schema.go`/`resolver*.go`），属临时件，**不是产码**。

### 1.3 "有没有任何赋值"——两棵树的正面对比

```console
$ git grep -n -E "\.MaxOutputTokens[[:space:]]*=" -- "*.go"        # HEAD
internal/agent/loop.go:261:	asm.MaxOutputTokens = opt.Config.MaxOutputTokens
internal/agent/maxoutput_370_test.go:36:		cfg.MaxOutputTokens = ...        # 测试
（其余为测试断言字符串）
```

```console
$ git grep -n -E "\.MaxOutputTokens[[:space:]]*=" 74f430d86^ -- "*.go"
(none)
```

**读数：**
- **父提交 `74f430d86^`：产码零赋值** ⇒ `Assembler.MaxOutputTokens` 恒为 Go 零值 `0` ⇒ `Request.MaxOutputTokens` 恒 `0`。票面主张**成立**。
- **HEAD `32fbed84c`：产码恰有一处赋值** `internal/agent/loop.go:261`，链已通。

---

## 2. resolver：配置值到底进了谁（AC#0②）——★票面这一句不准确

票面写「`[llm.roles.chat].max_output_tokens` 已经落进 `llm.Endpoint`（`resolver.go:205`）」。
**现读：它落进的是 `llm.RoleSettings`，不是 `llm.Endpoint`。** `llm.Endpoint` 根本没有 `MaxOutputTokens` 字段。

`internal/llm/resolver.go:186`（签名，逐字）：

```go
func (r *Resolver) ResolveRole(name RoleName) (Endpoint, RoleSettings, error) {
```

`internal/llm/resolver.go:202-206`（逐字）：

```go
	settings := RoleSettings{
		Temperature:       role.Temperature,
		ThinkingIntensity: role.ThinkingIntensity,
		MaxOutputTokens:   role.MaxOutputTokens,
	}
```

`internal/llm/resolver.go:43-58`（`Endpoint` 定义，逐字）：

```go
// Endpoint is one fully resolved provider/model target.
type Endpoint struct {
	Provider string // catalog name (diagnostics)
	Model    string
	Protocol string
	BaseURL  string
	APIKey   string // resolved plaintext; NEVER log
	Compat   config.Compat
	RPM      int
	TPM      int

	// Model catalog fields the seam consumes.
	Capabilities  config.Capabilities
	ContextWindow int
	MaxOutput     int
}
```

`internal/llm/resolver.go:142`（逐字）：

```go
		MaxOutput:     spec.MaxOutputTokens,
```

**归因：**
1. `config.Role.MaxOutputTokens`（`schema.go:298`）→ `llm.RoleSettings.MaxOutputTokens`（`resolver.go:205`）——**确实落库到返回值里**，是 `ResolveRole` 的**第二个**返回值。
2. `Endpoint.MaxOutput` 是**模型目录**的 `ModelSpec.MaxOutputTokens`（`schema.go:361` → `resolver.go:142`），语义是"这个模型最多能输出多少"，**不是**角色配置的天花板。二者今天数值可不同。
3. 所以 AC#1 的〔乙〕"在 Assembler 里从 `llm.Endpoint` 读"**前提不成立**：`Endpoint` 里没有这个值；要读只能先给 `Endpoint` 加字段（或让 Assembler 拿到 `RoleSettings`），而那不是"读一下"，是改 resolver 的产物形状。

---

## 3. 整条链：config → Endpoint → 断点 → Config → Assembler → Request（AC#0③）

### 3.1 HEAD `32fbed84c`（现读，已接通）

| 跳 | 载体 | file:line（逐字/读数） | 状态 |
|---|---|---|---|
| 1 | TOML `[llm.roles.chat].max_output_tokens` → `config.Role.MaxOutputTokens` | `internal/config/schema.go:290`（`type Role struct`）/ `:298` | ✅ |
| 2 | `config.Role` → `llm.RoleSettings.MaxOutputTokens` | `internal/llm/resolver.go:205`；返回 `(Endpoint, RoleSettings, error)` `:186` | ✅（**不是**进 `Endpoint`） |
| 3 | `RoleSettings` → `agent.Config.MaxOutputTokens` | 组装：`cmd/wisp/run.go:1498` 每任务重解析 → `:1560` `MaxOutputTokens: roleSettings.MaxOutputTokens` | ✅（370 补的那一跳） |
| 4 | `agent.Config` → `Assembler.MaxOutputTokens` | `internal/agent/loop.go:142`（字段）/ `:261` `asm.MaxOutputTokens = opt.Config.MaxOutputTokens` | ✅ |
| 5 | `Assembler` → `llm.Request.MaxOutputTokens` | `internal/agent/prompt.go:253` `MaxOutputTokens: a.MaxOutputTokens` | ✅ |
| 6 | `Request` → 线上字段 | openaichat `request.go:75`（`:60` `omitempty`）· anthropic `request.go:103`（`:107-108` 0→4096）· openairesponses `request.go:78`（`:55` `omitempty`） | ✅ |
| 7 | 线上 `length`/`max_output_tokens` → `StopMaxTokens` | openaichat `adapter.go:243-244`；openairesponses `adapter.go:240-243`；anthropic `adapter.go:285-286` | ✅ |
| 8 | `StopMaxTokens` → `StatusTruncated` → `task_log.state` | `internal/agent/loop.go:533-543`；映射 `:1129-1130` → `interrupted` | ✅ |

### 3.2 父提交 `74f430d86^`（票面立案时，确实断）

| 跳 | 载体 | file:line（现跑 `git show`） | 状态 |
|---|---|---|---|
| 1 | 同上 | `schema.go:298` | ✅ |
| 2 | 同上 | `resolver.go:205`（进 `RoleSettings`） | ✅ |
| 3 | `RoleSettings` → `agent.Config` | **断**：`cmd/wisp/run.go:503` 逐字 `ep, _, err := res.ResolveRole(llm.RoleChat)`——第二个返回值被 `_` 丢掉；`agent.Config`（`loop.go`）无该字段（父树 `git grep MaxOutputTokens -- loop.go` = **零命中**） | ❌ **缺字段＋组装丢弃** |
| 4 | `agent.Config` → `Assembler` | **恒 0**：`prompt.go:132` 只声明、`:253` 只读；`NewAssembler`（`:140-144`）不设它；全树 `.MaxOutputTokens =` 零命中 | ❌ **零值** |
| 5 | `Assembler` → `Request` | `prompt.go:253` 存在，但抄的是 0 | ⚠️ 恒 0 |
| 6 | `Request` → 线上 | 0 → openaichat `omitempty` 省略；anthropic 填 4096；openairesponses 省略 | ⚠️ 无"我们设的天花板" |
| 7-8 | 供应商自身上限 → `StatusTruncated` | 仍可达（见 §5） | ✅（但与配置无关） |

**"哪个环节断了"的精确答案：**
不是某一行的 `0`，而是**两处叠加**——
1. `agent.Config`（`internal/agent/loop.go:110` 那个结构）**没有承接字段**；
2. 组装根 `cmd/wisp/run.go:503` 把 `ResolveRole` 的第二返回值 `_` 掉了。
结果 `Assembler.MaxOutputTokens` 没有任何写入点（§1.3 父树 grep = none），永远是零值。

### 3.3 顺带更正票面另一处行号

票面引 `loop.go:490` 的 `if turn.Stop == llm.StopMaxTokens`。HEAD 上该判断在 **`internal/agent/loop.go:533`**（370 的注释/字段把行号推后了）。逐字：

```go
		// (D21) A length-truncated message never executes its tool calls.
		if turn.Stop == llm.StopMaxTokens {
```

---

## 4. 正控：手工设上限 ⇒ `StatusTruncated` / `task_log.state`（AC#0④）

### 4.1 本腿新增的聚焦测试（`internal/agent/maxoutput_370_ac0_test.go`）

不打配置、不过 resolver、不过 `cmd/wisp`：`Config.MaxOutputTokens = 1` 手设，C5 用假 provider 返回 `StopMaxTokens`，`Loop.Run` 走真状态机，`Journal` 接**真 SQLite**（`memory.Open`），最后从 `task_log` 读回。

```console
$ go test ./internal/agent -run "Test370AC0" -v -count=1
=== RUN   Test370AC0HandSetCeilingRowsTruncated
    maxoutput_370_ac0_test.go:130: positive control: request.MaxOutputTokens=1; Loop.Run status="truncated"; task_log.state="interrupted"
--- PASS: Test370AC0HandSetCeilingRowsTruncated (0.13s)
=== RUN   Test370AC0ProviderOwnLimitTruncatesWithoutOurCeiling
    maxoutput_370_ac0_test.go:160: no ceiling on the wire (body has no max_tokens); recorded finish_reason=length -> status="truncated"
--- PASS: Test370AC0ProviderOwnLimitTruncatesWithoutOurCeiling (0.01s)
PASS
ok  	github.com/CarlosShao/wisp/internal/agent	0.178s
```

**读数（逐字）：`request.MaxOutputTokens=1` → `Loop.Run status="truncated"` → `task_log.state="interrupted"`。**
测试没有构造 `Result`：`res` 是 `loop.Run` 返回值，行是 loop 自己经 journal 写进真库的。断言：
- provider 实际收到的 `req.MaxOutputTokens == 1`（证明手设值真到 C5 接缝）；
- `res.Status == StatusTruncated`；
- `tl.State == "interrupted"` 且在 196 名册上（`memory.ValidTaskState`）。

⇒ **状态机是好的；父提交上"到不了"确系配置没传。**

### 4.2 现成的端到端正控（HEAD 已在树上，本腿现跑复验）

`cmd/wisp/task370_ac2_test.go` 用**真 `wisp run`＋真组装根**，配置里写 `[llm.roles.chat] max_output_tokens = 1`，mockllm 在 `max_tokens=1` 处切断并回 `finish_reason:"length"`，真 openaichat 适配器归一化成 `StopMaxTokens`，最后从 memory DAO 读行。按约束带 PATH、只 `-run`：

```console
$ PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
  go test ./cmd/wisp -run "Test370TruncatedRunRowsAnInterruptedTask$" -v -count=1
=== RUN   Test370TruncatedRunRowsAnInterruptedTask
    task370_ac2_test.go:103: 屏幕逐字: wisp run: 任务 18240d5c-... 结束（truncated，1 轮，0 次工具调用，成本 0 CNY（未计价））
    task370_ac2_test.go:103: 库逐字: task_log.state = "interrupted"
    task370_ac2_test.go:110: AC#2(2) 屏幕逐字 = "truncated"；库 task_log.state 逐字 = "interrupted"
--- PASS: Test370TruncatedRunRowsAnInterruptedTask (1.86s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	1.971s
```

（`:103` 两行由 `p.surfaces()` 内的 `t.Helper()` 归因到调用行，实际打印在 `cmd/wisp/tasklogstate_365_test.go:157-158`。）

⇒ **正控成立：上限一设，`StatusTruncated` 可达；屏幕上这个词是 `truncated`，库里是 `interrupted`。**

### 4.3 agent 层的既有旁证（非本腿新增，现跑同绿）

```console
$ go test ./internal/agent -run "Test365JournalRowsTheRosterWordOnTheAgentPath" -v -count=1
    tasklogstate_365_test.go:122: agent path wrote task_log.state = "interrupted" for status "truncated"
--- PASS: Test365JournalRowsTheRosterWordOnTheAgentPath/max-tokens-toolcall
```

---

## 5. "供应商自己的上限"这条的射程（AC#0⑤）

### 5.1 现跑读数（本腿新增测试的第二段）

```console
$ go test ./internal/agent -run "Test370AC0ProviderOwnLimitTruncatesWithoutOurCeiling$" -v -count=1
    maxoutput_370_ac0_test.go:160: no ceiling on the wire (body has no max_tokens); recorded finish_reason=length -> status="truncated"
--- PASS
```

该测试用 golden fixture `internal/agent/testdata/golden/max-tokens-toolcall.sse`（**真实录制的 SSE**，真适配器解析），配置 `MaxOutputTokens` 保持 0，并断言第一发请求体里**没有** `"max_tokens"`。fixture 第 17 行逐字：

```
data: {"choices":[{"index":0,"delta":{},"finish_reason":"length"}]}
```

⇒ **`max_tokens` 完全省略时，供应商按其自身模型上限仍可回 `length`，仍会落 `StatusTruncated`。**
（另有反向旁证：`Test370UnsetRoleMaxOutputSendsNoCeiling` PASS，断言未设上限的请求体不含 `max_tokens`。）

### 5.2 三个适配器的逐一读数

| 协议 | 未设上限时线上发生什么 | file:line | 仍能回 `StopMaxTokens`？ |
|---|---|---|---|
| openai-chat | 字段被 `omitempty` **省略**，供应商用自己的默认/模型上限 | `request.go:60` `json:"max_tokens,omitempty"`；`:75` | ✅ `adapter.go:243-244` `case "length": return llm.StopMaxTokens` |
| openai-responses | 字段被 `omitempty` **省略** | `request.go:55`、`:78` | ✅ `adapter.go:240-243` `case "max_output_tokens", "": return llm.StopMaxTokens` |
| anthropic | 协议**要求** `max_tokens`，适配器把 0 填成**自己的 4096**，不是省略 | `request.go:107-108`；`adapter.go:94` `defaultMaxTokens = 4096` | ✅ `adapter.go:285-286` `case "max_tokens": return llm.StopMaxTokens` |

### 5.3 结论（票面两条射程，逐条判）

- **"我们设的那条天花板到不了"**（父提交 `74f430d86^`）＝ **真**：`Assembler.MaxOutputTokens` 零赋值 ⇒ 配置值从未上线。
- **"`StatusTruncated` 这个状态到不了"** ＝ **假**：供应商自身上限（以及 anthropic 适配器硬编码的 4096）仍会回 `length`/`max_tokens` ⇒ `loop.go:533` 仍会走到 `StatusTruncated`。
- 补充：**HEAD 上两条都不成立了**——配置天花板已经上线（§3.1、§4.2）。
- 另注意 anthropic 的 4096 是**适配器本地数**，不是 C 级契约数；它意味着 anthropic 路径上永远有一个"我们的"天花板（只是不可配）。

---

## 6. AC#1 裁形输入（本腿只报，编排者拍）

### 6.1 `agent.Config` 是不是 C 级契约面？——**不是**（no），证据如下

1. `docs/PLAN.md:1350` 逐字：**「## 三、冻结契约清单（写码前定死，agent 只能实现不能修改）」**；其表 `C1–C32`（`PLAN.md:1357-1388`）是 interface 级清单。逐条看：C1 `Tool` · C2 `RiskLevel` · C3 `Capability` · C4 `ToolProvider` · **C5 `LlmProvider`**（`:1361`）· C6 `StreamEvent` · C7 内部消息表示 · C8 `AudioSource` · … · C32 `RealtimeEngine`。**没有 `agent.Config`，也没有 `agent.Options`。**
2. `docs/specs/SPEC-12-roadmap-governance.md:36-38` 逐字：**「### 4.1 契约变更流程」／「改 C1–C32 或 D1–D47 = 人工批准」**——射程就是那两张表。
3. 现跑 `git grep -n "agent\.Config" -- docs/specs docs/PLAN.md` ⇒ **零命中**（`(end)`），即规格/契约文字从未把这个结构体列为冻结面。
4. `agent.Config` 是 `internal/agent/loop.go:110` 的**包内 Go 结构体**，不是 `docs/contracts/` 里的 interface 契约（AGENTS §4 已说明 `docs/contracts/C1..C32.md` 截至锚点不存在）。所有构造点都是**带字段名的复合字面量**（`cmd/wisp/run.go:1554` 等 9 处），加一个字段对既有调用源兼容。
5. 邻近的 C5 面 `llm.Request` **早在父提交就有** `MaxOutputTokens`（`internal/llm/provider.go:74`，父树同样存在），所以本票三种形都**不需要改 C5 契约**。

**⇒ 加 `Config.MaxOutputTokens` 本身不构成 C1–C32/D1–47 契约变更，不需要为"契约文字"单独报批。**
但事实层面要说清：落地提交 `74f430d86` 仍走了机主批准（台账 `:30628`「机主批准〔甲′〕」），所以"已获批"是既成事实；AC#1 今天不是"要不要改契约"，而是"要不要维持已落地的甲′"。

### 6.2 三形代价表（HEAD 现读；"已落地"指 `74f430d86`）

| 形 | 要动的文件（HEAD 行号） | 需要 `cmd/wisp` 组装？ | C 级契约？ | 代价 / 结论 |
|---|---|---|---|---|
| **甲** thread 到 `Request` | `cmd/wisp/run.go:1498,1503,1560`；`internal/agent/loop.go:133-142,261`；`prompt.go` **不用动**（`:132/:253` 早就有）；测试 `internal/agent/maxoutput_370_test.go`、`cmd/wisp/task370_ac2_test.go` | **是**（`run.go` 每任务重解析＋塞进 `agent.Config`） | 否 | **已落地**：`74f430d86` 3 文件 +167/−3（产码约 30 行，其余测试）。职责正确：per-role 的值走 per-role 的 `RoleSettings`。**推荐维持。** |
| **乙** Assembler 从 `llm.Endpoint` 读 | 先得给 `internal/llm/resolver.go:44-58` 的 `Endpoint` **加字段**并在 `resolveEndpoint`/`ResolveRole` 填（`resolver.go:142,205`）；`internal/agent/loop.go:260-261` 与 `prompt.go:132` 改读法；`cmd/wisp/run.go` 仍要把 endpoint/设置交给 `agent.New` | **是**（实际仍要，`agent.New` 的入参不从天上掉） | **否**（`Endpoint` 不在 C1–C32 表里；但它是 resolver→C5 工厂的产物形状，评审负担高于甲） | 票面说"会造 `agent`→`llm` 新依赖"——**不成立**：`prompt.go:9` 早就 `import "github.com/CarlosShao/wisp/internal/llm"`（`llm.CacheSupport`）。真代价是**语义错位**：`Endpoint.MaxOutput`（`:57`）已是"模型目录上限"（来自 `ModelSpec`），再塞"角色配置上限"会让两个不同来源的数共用一个结构；且每角色值仍要先从 `RoleSettings` 搬过去。**比甲多改一个契约邻近结构、零收益。否决。** |
| **丙** 删 `StatusTruncated`＋映射 | `internal/agent/loop.go:50-52`（常量）、`:62`（`Valid` 名册）、`:533-543`（截断分支，需改写成别的终态）、`:1111-1130`（映射）；`cmd/wisp/run.go:83`（退出码分支）；测试 `internal/agent/tasklogstate_365_test.go:43,96`、`cmd/wisp/task370_ac2_test.go` 全文件、`internal/agent/maxoutput_370_test.go` 部分 | **是**（退出码/屏幕词） | 否（不在 C1–C32；`memory.TaskStateInterrupted` 仍被崩溃恢复与 196 名册使用，不能一起删） | 删掉的是**真实语义**：供应商撞上限后，任务将只能落 `failed`→`error`，与真失败不可分；`failToolCallsFromTruncatedMessage` 的收尾消息也失去可命名的终态。票面自己写「⛔ 推荐」。**否决。** |

**一句话推荐：维持〔甲′〕。** 理由：它是唯一把"每角色配置的上限"放在"每角色设置"上的形，不改 `Endpoint` 的模型目录语义、不删真实状态；而且已获批、已落地、在 HEAD 上两枚真读数全绿（§4.1/§4.2），再动只会把已兑现的链重新拆开。

---

## 7. 复跑清单（尺子）

```console
$ cd "/d/work/workspace/projects plans/Wisp"
$ git rev-parse --short HEAD                                   # 32fbed84c
$ git grep -n "MaxOutputTokens" -- "*.go"
$ git grep -n -E "\.MaxOutputTokens[[:space:]]*=" -- "*.go"    # 产码仅 loop.go:261
$ git grep -n -E "\.MaxOutputTokens[[:space:]]*=" 74f430d86^ -- "*.go"   # (none)
$ git show 74f430d86^:internal/agent/prompt.go | grep -n MaxOutputTokens # 132,253 only
$ git show 74f430d86^:cmd/wisp/run.go | grep -n ResolveRole    # :503 ep, _, err
$ go test ./internal/agent -run "Test370AC0" -v -count=1
$ go test ./internal/agent -run "Test365JournalRowsTheRosterWordOnTheAgentPath" -v -count=1
$ PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
  go test ./cmd/wisp -run "Test370TruncatedRunRowsAnInterruptedTask$" -v -count=1
```

本腿写面：**只** `internal/agent/maxoutput_370_ac0_test.go`（新增测试）＋ 本文件。产码零字节改动；未碰 `frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、任何 golden fixture。
