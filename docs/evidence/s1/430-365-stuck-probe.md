# 票 430 — `Test365RunRowsAStuckTaskAsStuck` 恒红的判定与落地

- 腿：430-365（写腿，只 commit 不 push）
- 日期：2026-11-12
- 票面：`.scratch/wisp/issues/430-the-365-stuck-task-probe-never-goes-stuck-because-the-tiny-context-window-never-becomes-an-overdraw.md`
- 判据：**仓外 overlay 探针**（`go test -overlay`，拷贝放 `D:\tmp\wisp\20261112\430-365\`），仓内零临时件

## 一句话判定

**产品缺陷**，⛔ 不是台面/夹具问题。票面假设的"造透支手法没生效"是**反的**：`context_window = 1` 完美地生效了，C22 预算**真的**被缩到 2 token（`BudgetsFor(1).TokenBudget = 2`）；但 `[agent] token_budget` 的 **`default:"200000"` 结构标签**（`internal/config/schema.go:472`，改前）在 config 加载时把这个 2 **覆盖成 200000**，于是 `guard.TokenBudget() = 200000`，一发 17 token 的回答**离透支差得极远**，`brakeStuck` 自然一次都没被调用。

## 三个探针读数（决定性）

探针把 `internal/agent/loop.go` 拷贝到仓外并在 `New`／`run`／`AddUsage`／`brakeStuck`／循环顶五处打点，用 `-overlay` 映射回包内（`D:\tmp\wisp\20261112\430-365\overlay\loop.go`）：

```
[P430] loop.New: opt.Config.ContextWindow= 1  info.MaxContextWindow= 1  chosen window= 1
       => BudgetsFor.TokenBudget= 2  HistoryCompressTokens= 1  Scale= 7.8125e-06
       opt.Config.TokenBudget= 200000
[P430] run: guard.TokenBudget= 200000  guard.MaxRounds= 50
[P430] turn usage: in= 12  out= 5  toolCalls= 0  stop= end_turn
```

| # | 要的数 | 读数 |
|---|---|---|
| ① | 换算后的 C22 预算值 | **2 token**（`BudgetsFor(1).TokenBudget = 2`，`Scale = 1/128000`）—— **换算链本身是对的** |
| ② | 本轮实际消耗 token 数 | **17**（`in=12 + out=5`，provider 自报 usage，非估算） |
| ③ | `brakeStuck` 是否被调用 | **一次都没有**（探针在 `loop.go:1058` 入口打点，全程零输出） |

**关键的第 4 个数（探针额外打的）**：`guard.TokenBudget = 200000`，而 `BudgetsFor` 算出的是 **2**。这两个数在同一次运行里同时存在，就是缺陷本体：`2` 被 `200000` 顶掉了。

## `StartMockllm` 的响应形状

**没有 tool_call**。`turn usage: ... toolCalls= 0  stop= end_turn` —— mock 回的是纯文本 `end_turn` 收尾。

这**不是**病因，也**不需要**是：2 token 的预算下**任何**非空回答都该透支，mock 回不回调工具无关。`0 次工具调用` 是**果**（模型没工具可调／不需要调），⛔ 不是因。票面 `## ★关键线索` 那一节把因果读反了 —— 真正的因在 config 层，不在 mock 层。

## 支撑的 `文件:行`

| 环节 | 位置 | 事实 |
|---|---|---|
| 夹具确实写了 1 | `cmd/wisp/tasklogstate_365_test.go:116` | `context_window = %d`，值来自 `:179` 的 `newProbe365(t, 1)` —— **逐字确认，配置是对的** |
| 配置真的读到了 1 | `internal/llm/resolver.go:141` | `ContextWindow: spec.ContextWindow` |
| 1 ⇒ 2 token 的换算 | `internal/agent/budgets.go:108` | `b.TokenBudget = s(refTokenBudget)`；`scaleInt`（`:114`）把 `200000 × 1/128000 + 0.5` 取整 ⇒ **2** |
| **覆盖发生的地方** | `internal/config/schema.go:472`（改前） | `TokenBudget int \`toml:"token_budget" default:"200000"\`` |
| 默认值被填进去 | `internal/config/defaults.go:80-81` | `applyDefaults` 见 `default` 标签即 `setDefault` ⇒ 用户**没写**这个键也会变成 200000 |
| 覆盖值被当"用户显式值"喂进去 | `cmd/wisp/run.go:2021` | `TokenBudget: cfg.Agent.TokenBudget` |
| 于是 scaled 兜底永不执行 | `internal/agent/guard.go:104-105` | `if cfg.TokenBudget <= 0 { cfg.TokenBudget = cfg.Budgets.TokenBudget }` —— 200000 > 0，兜底被跳过 |
| 刹车条件（本身没坏） | `internal/agent/loop.go:457` | `guard.TokenBudget() > 0 && guard.TokensUsed() >= guard.TokenBudget()` ⇒ `17 >= 200000` 为假 |
| 刹车函数（本身没坏） | `internal/agent/loop.go:1058` | `brakeStuck` —— 零调用 |

**回归引入点**：`a7f309112`（423-followup，2026-11-12 前后）第一次把 `TokenBudget: cfg.Agent.TokenBudget` 接进 `cmd/wisp/run.go`。改前（`a7f309112^`）这行不存在，`guard.go:104` 的 scaled 兜底正常工作，票 365 的测试因此曾经是绿的。`default:"200000"` 标签本身早于该 commit（`fe126e278`），但**没人读它**，所以一直是死标签；423-followup 把它接活的那一刻，这个死标签就变成了活的覆盖源。

## 为什么判"产品"而不是"台面"

D15 逐字要求缩放**是 MUST**：

> 所有阈值必须按所配模型的 `context_window` 等比缩小，不得硬编码（`docs/PLAN.md:391` §D15）

`default:"200000"` 让"用户没写这个键"和"用户要求 200k"**无法区分**，等于把一个按 128k 基准定死的硬编码数字变成了**任何**小窗口模型的实际上限。这不只砸了票 365 这一枚测试 —— 它让**每一个**声明小 `context_window` 的真实用户都拿到 200k 预算，与 D15 直接矛盾。台面（测试）是**对的**，它拒绝作证的那句话也是对的。

## 落地（只改一处，产品侧）

`internal/config/schema.go:454-472`：删掉 `TokenBudget` 的 `default:"200000"` 标签，改为无标签 + 长注释说明为什么不能有。

- **与兄弟键一致**：同一 section 的 `per_tool_timeout_ms`（`schema.go:475`）本来就没有 `default` 标签，同样的"0 = 不覆盖，用下游自己的默认"语义。本改动让 `token_budget` 与它对齐。
- **用户显式写的值仍然绝对生效**：`guard.go:104` 只要求 `> 0`，写 `token_budget = 12345` 依旧赢过缩放值。
- **128k 基准行为不变**：`BudgetsFor(128000).TokenBudget` 仍是 200000（`budgets.go:37`），即 D15 表里那个"默认 200k"。
- **不碰任何断言**、不碰 `thresholds.go`／golden／`PLAN.md`／`specs`。

### 附带读数：canonical 文件现在会渲染 `token_budget = 0`

`MarshalCanonical` 无 omitempty、逐字段全发（D36 rule 4，`internal/config/parse.go:160-177`）。删掉标签后 `NewDefaults().Agent.TokenBudget = 0`，于是新建的 config.toml 会带一行 `token_budget = 0`。探针实测：

```
[P430] NewDefaults().Agent.TokenBudget = 0
[P430] CANON "max_rounds = 50"
[P430] CANON "token_budget = 0"
```

这是**语义正确的**（0 = 不覆盖 = 用缩放默认），且与 `per_tool_timeout_ms = 0` 一直以来的渲染形状同形。`TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags` 通过 —— 它按标签推导期望，标签没了它也就不再期待这一行有值。

## 突变控制（真跑，两发）

| 发 | 动作 | 读数 |
|---|---|---|
| 红 | 仓外拷贝把 `default:"200000"` **放回**，`-overlay` 映射（`D:\tmp\wisp\20261112\430-365\mut.json`） | **`--- FAIL: Test365RunRowsAStuckTaskAsStuck (1.97s)`**，屏幕 `completed`、库 `done` |
| 绿 | 恢复（仓内就是修好的版本） | **`--- PASS: Test365RunRowsAStuckTaskAsStuck (2.00s)`**，屏幕 `stuck`、库 `stuck` |

修复后逐字：

```
tasklogstate_365_test.go:181: 屏幕逐字: wisp run: 任务 ... 结束（stuck，1 轮，0 次工具调用，成本 0 CNY（未计价））
tasklogstate_365_test.go:181: 库逐字: task_log.state = "stuck"
```

⚠ 注意 `0 次工具调用` **修好后依然是 0** —— 再次证明它从来不是病因。刹车是**第一发回答就透支**触发的（17 ≥ 2），不需要任何工具结果。

## 稳定性（三次复跑，全红）

| 发 | 场景 | 读数 |
|---|---|---|
| 1 | 整包内（编排者） | FAIL 4.57s |
| 2 | 单独（编排者） | FAIL 4.57s |
| 3 | 单独（本腿，carrier） | **FAIL 2.69s**，屏幕 `completed`／库 `done`／exit 0 |
| 4 | 单独 + overlay 探针 | **FAIL 3.82s** |

同句同数 ⇒ 稳定真红，⛔ 偶发。

## 门禁四读数

| 门 | 读数 |
|---|---|
| `go build ./cmd/wisp/` | **rc=0**，零输出 |
| `go vet ./cmd/wisp/` | **rc=0**，零输出 |
| `gofumpt -l`（`internal/config/schema.go`） | **零输出**，rc=0 |
| `scripts/d22scan.sh` | **clean** —— `examined 480 production Go files`，`no D22 ban violations`，自测 `PASS=45 FAIL=0 SKIP=0` |

## 越界检查

`git diff --name-only` 里属于本腿的只有 **`internal/config/schema.go`** 一枚。未出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`。

（同一工作树里另有其他腿的在飞改动 —— `.scratch/wisp/issues/42x`、`.scratch/wisp/probes/291`、`cmd/wisp/resident_watchdog_windows_test.go`、`frontend/fixtures/composer-states.html` —— 本腿**不 add 它们**，commit 走显式 pathspec。）

## 与票 365／164 的关系

- **票 365**：本票只让它的宿主测试重新说真话。⛔ 不翻票 365 任何 AC 框。
- **票 164**（`brakeStuck`）：**它没坏**。`loop.go:1058` 的 `brakeStuck` 与 `loop.go:457` 的触发条件一字未动，探针读数证明预算一旦真的变小它就正常开火。⛔ 不翻票 164 任何框。
- **票 423**：本改动与 423 的"接活 `[agent]` 尾部"方向**不冲突** —— 用户显式写 `token_budget` 依旧绝对生效；只是"没写"不再被伪装成"写了 200k"。423 的 `Test423AgentTailReachesAgentConfig`（写 `12345`，断言到达）通过。

## 复跑命令（可逐字重跑）

```bash
# 目标测试（carrier，必须是 Git bash，⛔ 不是 PATH 上的 WSL 转发器）
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^Test365RunRowsAStuckTaskAsStuck$'"

# 整个 365 家族
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^Test365'"

# 门禁
go build ./cmd/wisp/ && go vet ./cmd/wisp/ && D:\work\base\gopath\bin\gofumpt.exe -l internal/config/schema.go
D:\work\soft\Git\bin\bash.exe -c "scripts/d22scan.sh"
```
