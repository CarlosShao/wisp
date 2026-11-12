# 430 — `Test365RunRowsAStuckTaskAsStuck` 恒红：`context_window=1` 本该把 C22 预算缩到 2 token 逼出 `brakeStuck`，但那发**正常 completed、0 次工具调用**

- Status: **ready-for-agent**（先只读取证，判定后写面才定；⚠ 预期是台面/夹具问题，但**由你判清后自己定并写理由**）
- 来源：round 88 整包复跑（`D:\tmp\wisp\20261112\fullclean\cmdwisp.log`）＋**我亲跑两次**（`D:\tmp\wisp\20261112\365\run1.log` 4.57s；整包内同句）
- 关联：票 365（本测试宿主，`task_log` 的 stuck 行）· 票 164（`brakeStuck`）· C22 上下文预算 · D15 · 台账 `A1803`

## 这是什么（人话）

有个测试要证明：**任务因为"上下文预算透支"被刹车停住时，屏幕上写的和库里存的必须是同一个词（stuck）**，
而不是像以前那样屏幕写 "stuck"、库里写 "error"。

它的造法很巧：把 `context_window` 配成 **1**，让 C22 的 200k token 预算**缩到 2 token**，
这样模型**第一次回答就透支**，循环就会走 `brakeStuck` 停下。

现在：**那发没被刹车**——屏幕写 `completed`、库里 `state = "done"`、**0 次工具调用**、exit 0。
测试于是很诚实地拒绝作证（这句断言写得极好）：
> `this run did not end stuck (screen word "completed", exit 0), so it cannot testify about stuck`

## 现量（可逐字复跑）

**逐字三行**（`cmd/wisp/tasklogstate_365_test.go:181`／`:181`／`:184`）：

```
tasklogstate_365_test.go:181: 屏幕逐字: wisp run: 任务 d663d4dc-abed-4350-beba-269868128b51 结束（completed，1 轮，0 次工具调用，成本 0 CNY（未计价））
tasklogstate_365_test.go:181: 库逐字: task_log.state = "done"
tasklogstate_365_test.go:184: this run did not end stuck (screen word "completed", exit 0), so it cannot testify about stuck:
```

**稳定性**：**我亲跑两次都红**（整包内 1 次＋单独 1 次 **FAIL 4.57s**），同句同数 ⇒ **稳定真红，⛔ 偶发。**

## ★关键线索：`0 次工具调用`

配置是对的（`newProbe365(t, 1)` ⇒ `context_window = %d` 传 1，见 `:95`／`:116`）。
但**模型一次工具都没调**，所以：
- 没有大工具结果 ⇒ 没有大上下文 ⇒ **没有透支** ⇒ 刹车不触发 ⇒ 正常 completed。
- 也就是说：**测试的"造透支"手法在这条链上没生效**，⛔ 不是 `brakeStuck` 坏了。

## 机制链（逐行可核）

测试文件：`cmd/wisp/tasklogstate_365_test.go`
- `:95` `func newProbe365(t *testing.T, contextWindow int) *probe365`
- `:97` `p.srv = adaptertest.StartMockllm(t)` —— **mock LLM**
- `:116` `context_window = %d` 写入 config，值来自参数
- `:120` `[fs] allowed_dirs = [%q]`
- `:178` `func Test365RunRowsAStuckTaskAsStuck(t *testing.T)`
- `:179` `p := newProbe365(t, 1)` —— **contextWindow = 1**
- `:180` `code := p.run("把目录再列一遍")` —— **这句本该让模型去调 `fs.list` 之类的工具**
- `:181` `screenWord, rowState := p.surfaces()`
- `:183` `if screenWord != string(agent.StatusStuck)` ⇒ **`:184` Fatal**
- `:187` `if code != 3` ⇒ 期望 exit 3（loop brake）
- `:190` `if rowState != string(memory.TaskStateStuck)`

**注释 `:174-177` 逐字说明了设计意图**：
> `A context_window of 1 scales the frozen 200k C22 token budget down to 2 tokens, so the first response overdraws it and the loop ends through brakeStuck - StatusStuck, the status the ticket was filed about. Nothing is simulated: the brake fires, the loop finishes, this command writes the row.`

## ⚠ 必须自己判清的岔路（⛔ 不许我替你定）

1. **`StartMockllm` 回的响应是什么？** 它是否**回了工具调用**？
   - 若 mock 只回纯文本（不调工具）⇒ "第一次回答就透支"这个前提**是否成立取决于响应有多大**。
     2 token 的预算下，**任何**非空回答都该透支——所以这解释不了。
   - 更可能的解释：**预算没被真的缩到 2 token**（`context_window` 那一路没接上／被 clamp／被默认值覆盖）。
     ⇒ 去读 `context_window` 从 config 到 C22 预算的**完整换算链**，逐行确认 1 ⇒ 2 token。
2. **`brakeStuck` 的触发条件是什么？** 读 `internal/agent` 里 `brakeStuck`／`StatusStuck` 的判定，
   确认"透支"是否真的会走到它，还是走了别的分支（比如静默截断后继续跑）。
3. **`0 次工具调用` 是果还是因？** 若模型本该调工具却没调，可能是 mock 响应里没有 tool_call；
   若预算真的只有 2 token，模型连 tool_call 都发不出来（被截断）⇒ 那"透支"应该已经发生了。
   **这两者互斥，必须用读数分开。**

**判据**：在 `p.run` 前后插**仓外 overlay 探针**，打出
①换算后的 C22 预算值 ②本轮实际消耗 token 数 ③`brakeStuck` 是否被调用。
三个数一出来，岔路自明。

## 要做的（按序）

- [ ] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. **复跑 ≥2 次**确认稳定（我已有两次读数 4.57s／整包内，你再加一次）：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^Test365RunRowsAStuckTaskAsStuck$'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. 读 `:95`–`:125` 的 `newProbe365`，**逐字确认** `context_window = 1` 真的写进了 config。
  3. 读 `context_window` ⇒ C22 预算的**完整换算链**（哪枚文件哪一行把 1 变成 2 token？）。
  4. 读 `brakeStuck` 的触发条件。
  5. 用**仓外 overlay 探针**打出上述三个数（预算值／实际消耗／`brakeStuck` 是否被调用）。
  6. 读 `StartMockllm` 回的响应形状（有没有 tool_call）。
  **产出：一句话判定「台面/夹具问题」还是「产品缺陷」，＋三个探针读数＋支撑的 `文件:行`。**

- [ ] **AC#1 按 AC#0 的判定修（只修被判定那一处）**：
  - 若判「台面/夹具问题」（最可能）⇒ 修台面，让"造透支"真的成立。**⛔ 不许**改成直接调 `brakeStuck` 造假——
    注释 `:176` 逐字要求 `Nothing is simulated: the brake fires`。
  - 若判「产品缺陷」⇒ 修产码（预算换算那一路）。
  - **必须配突变控制真跑**：把修好的那处删掉 ⇒ 断言必须回到红；恢复 ⇒ 绿。

- [ ] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [ ] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`，独立 module）。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是未跟踪残留件，会让 `go build ./...` 红。**⛔ 不许删**。

- [ ] **AC#4 与票 365／164 的关系写清**：本票只处理这一枚，⛔ 不翻票 365／164 任何框。

## 本票**不**解决

- 不碰 `Test26P*`（票 428 在飞）／`Test384*`（票 429 在飞）／`TestAC1ResidentLeg*`（另计）。
- 不改 `docs/PLAN.md` 与 `docs/specs/**` 一字。

## 派单约束

- 先只读（AC#0 零写），判定后再定写面；**产码与台面二选一，⛔ 两边都改**。
- ⛔ `git add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；`add`／`commit` 两处显式 pathspec。
- 只 commit，⛔ push。
- ⛔ 翻任何 AC 框。
- 临时件只建不删，落 `D:\tmp\wisp\<YYYYMMDD>\430-<腿码>\`。
- **仪器**：这枚是**运行时行为**检查 ⇒ `go test -overlay`（拷贝放**仓外**）有效。

## Progress log（append-only，newest last）

### 430-365 写腿 · 2026-11-12 20:5x +08

**判定：产品缺陷**（⛔ 不是台面/夹具问题）。票面「造透支手法没生效」的假设**是反的**：夹具写 `context_window = 1` 完全生效，C22 预算**真的**缩到了 2 token。

**三个探针读数**（仓外 overlay 探针，`D:\tmp\wisp\20261112\430-365\`）：

1. **换算后预算 = 2 token**（`BudgetsFor(1).TokenBudget = 2`，`Scale = 7.8125e-06`，`HistoryCompressTokens = 1`）⇒ 换算链本身正确。
2. **实际消耗 = 17 token**（`in=12 + out=5`，provider 自报 usage）。
3. **`brakeStuck` 零调用** —— 探针打在 `internal/agent/loop.go:1058` 入口，全程无输出。

第 4 个数（决定性）：同一次运行里 `BudgetsFor.TokenBudget = 2` 而 `guard.TokenBudget = 200000` —— 缩放值被**顶掉**了。

**病因**：`internal/config/schema.go:472`（改前）`TokenBudget int \`toml:"token_budget" default:"200000"\``。`defaults.go:80-81` 的 `applyDefaults` 见标签即填 ⇒ 用户**没写**这键也变 200000 ⇒ `cmd/wisp/run.go:2021` 当显式覆盖喂进 `agent.Config.TokenBudget` ⇒ `internal/agent/guard.go:104` 的 `if cfg.TokenBudget <= 0` scaled 兜底**永不执行** ⇒ `loop.go:457` 的 `17 >= 200000` 为假 ⇒ 不刹车、`completed`、`done`、exit 0。

**回归引入点**：`a7f309112`（423-followup）首次把 `TokenBudget: cfg.Agent.TokenBudget` 接进 `run.go`；`a7f309112^` 没有这行，scaled 兜底正常，故 365 测试曾绿。`default:"200000"` 标签本身更早（`fe126e278`）但此前是**死标签**（没人读）。

**`StartMockllm` 响应**：**没有 tool_call**（`toolCalls=0`，`stop=end_turn`）。⛔ 不是病因——2 token 预算下任何非空回答都该透支。票面「0 次工具调用⇒没有大东西进上下文⇒不超支」把因果读反了；`0 次工具调用` 是果。修好后**依然是 0 次**，再次证明无关。

**修法（只改产品侧一处）**：`internal/config/schema.go` 删掉 `token_budget` 的 `default:"200000"` 标签，改为无标签 + 长注释。与同 section 兄弟键 `per_tool_timeout_ms`（`:475`，本来无标签）对齐；用户显式写的值仍绝对生效（`guard.go:104` 只要 `> 0`）；128k 基准仍是 200000（`budgets.go:37`）。

**突变控制（真跑两发）**：仓外拷贝把 `default:"200000"` 放回 + `-overlay` ⇒ **`--- FAIL: Test365RunRowsAStuckTaskAsStuck (1.97s)`**（屏幕 `completed`／库 `done`）；恢复 ⇒ **`--- PASS (2.00s)`**（屏幕 `stuck`／库 `stuck`）。

**稳定性**：本腿单独复跑 2.69s FAIL，加探针 3.82s FAIL，修后 2.00s PASS。同句同数，稳定真红，⛔ 偶发。

**附带读数**：`MarshalCanonical` 无 omitempty（D36 rule 4），删标签后新建 config.toml 会渲染 `token_budget = 0`（探针实测），语义正确（0 = 不覆盖）且与 `per_tool_timeout_ms = 0` 同形；`TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags` 通过。

**门禁**：`go build ./cmd/wisp/` rc=0 · `go vet ./cmd/wisp/` rc=0 · `gofumpt -l` 零输出 · `scripts/d22scan.sh` clean（480 files，`no D22 ban violations`，自测 `PASS=45 FAIL=0`）。

**越界**：本腿 diff 只 `internal/config/schema.go` + 本证据件 + 本 Progress 条。⛔ 未碰 `PLAN.md`／`specs`／`thresholds.go`／golden／`allowlist.txt`／`frontend`／`design`。

**与 365／164**：票 164 的 `brakeStuck` **没坏**（`loop.go:1058` 与触发条件 `:457` 一字未动）；本票只让 365 的宿主测试重新说真话。⛔ 未翻 365／164 任何 AC 框。

**旁证（非本腿引入）**：`TestTicket198FirstRunCallerIsTheRunEntryOnly` 在 `HEAD` 版 schema 上同样 FAIL（`firstrun_198_test.go:281` AC J1 RED，`ensureFirstRunConfig` 被 `[runConfigSet65 runTextTask]` 调用），与本改动无关，已用 `-overlay` 回退基线验证。

**证据件**：`docs/evidence/s1/430-365-stuck-probe.md`。
