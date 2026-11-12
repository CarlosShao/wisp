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
