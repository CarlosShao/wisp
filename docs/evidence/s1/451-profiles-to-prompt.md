# 451 写腿 · L1 画像进 prompt（`agent.Config.Profiles` 0 → 1，AC#1 先红后绿）

- **谁／何时**：451 写腿（实现者），2026-11-15 12:1x–12:5x +08。基 HEAD `bf4adaa7b`；**产码 commit `93a1c93da`**（3 文件，+241/−0）；本证据件另枚 commit（只追加）。
- **角色**：实现者 ≠ 裁决者。只 commit 不 push；裁决/对抗验收由别的 agent 做。
- **一句话结论**：AC#1 那枚测试在接线前**先红**（红句见 §②，模型确实看不到库里那条画像），接线后**转绿**；`Profiles:` 产码赋值点 **0 → 1**（唯一装配点 `cmd/wisp/run.go` 的 `execute()`）；定向突变 **摘掉接线那一行 ⇒ `Test451L1ProfileReachesTheProviderPrompt` 红**，恢复后逐字节相同、再绿。

---

## ① 现量（起手，可逐字重跑）

**尺（产码、非测试、剔除 `.scratch/` 探针）：**
`git grep -n 'Profiles:' HEAD -- '*.go' ':!*_test.go' ':!.scratch/**'`

| 时点 | 读数 |
|---|---|
| 起手（`HEAD=bf4adaa7b`） | **1 行**：`internal/agent/loop.go:677` → `Profiles:  l.opt.Config.Profiles,`（**转发**，`PromptInput` 那一侧）⇒ 构造 `agent.Config` 时给 `Profiles` 赋值的点 = **0** |
| 接线后（工作树） | **2 行**：老那行不变 ＋ 新增 `cmd/wisp/run.go:2123` → `Profiles: rt.loopProfileLines451(),` ⇒ **赋值点 0 → 1** |

**字段本身没删**：`internal/agent/loop.go:154-156` 的 `Config.Profiles` 原样保留；`internal/agent/loop.go`、`internal/agent/prompt.go`、`internal/memory/**` 本波**一字未动**（`git diff HEAD -- internal/agent internal/memory` 空输出）。

---

## ② AC#1 先红（接线之前，产码未动时跑）

**测试**：`cmd/wisp/agent_profiles_451_test.go:50` `Test451L1ProfileReachesTheProviderPrompt`。它做三件事，全部走真产码：

1. 用**产品自己的 DAO**（`memory.Open(f.dir)` + `UpsertProfile(ctx, "fact.family", "mother birthday 05-06", ProfileManual)`）把一枚 L1 画像写进**这次 run 会打开的那个 data dir**；
2. 跑**真的组合根** `runTextTask`（`newRunFixture(t, "openai-chat")`：config.toml → secret store → provider → approval gate → tools bridge → real `agent.Loop`），任务文本 `总结一下 这份笔记`（**不含**那条事实）；
3. 读 **mockllm 自己捕获的请求体**（`LastRequest(t, "chat")`），断言请求里有 `## profile` 段、且段内有那条事实。ASCII 断言（`## profile` / `mother birthday 05-06`），不受 JSON 转义影响。

**命令**：`scripts/wisp-cli-staged-tests.sh -test.run '^Test451L1ProfileReachesTheProviderPrompt$'`
**读数**：`CARRIER_EXIT=1`；日志 `.scratch/wisp/probes/451/r1/ac1-red.txt`。

**红句原文（接线前，逐字）**：

```
    agent_profiles_451_test.go:77: 送给 provider 的 prompt 里没有画像段（找不到 "## profile"）：库里那枚 L1 画像 "mother birthday 05-06" 没有被递进请求，模型看不到这条事实。mockllm 捕获的请求体：
```

⇒ 这句话读出来就是票面要的「prompt 里没有这句话」：库里有行，请求里连画像段都没有。

---

## ③ AC#2 接线（唯一装配点；谁建 Config、谁查库）

**谁建 `agent.Config`**：`cmd/wisp/run.go` 的 `execute()`（`:2012` 起）是这个进程**唯一的装配点**——每个任务的 `agent.Options`（含 `Config`）都在这里构造；根 loop 用它，`task.spawn` 派生子 loop 用的 `SubagentDeps.BaseOptions`（`rt.loopOpt`）也是它的快照。

**谁查库**：同一个装配点里调 `rt.loopProfileLines451()`：

- 实现：`cmd/wisp/agent_profiles_451.go:71`；
- 读的库：`assembleRuntime` 打开、挂在 runtime 上的**同一个** `*memory.Store`（`cmd/wisp/run.go:701` `memory.Open(s.dataDir, memory.WithProfileCap(cfg.Memory.L1Max))`）；
- 查询：`rt.store.ListProfiles(ctx)`（`:77`），**每任务一次**（不是 per boot，也不是 per turn）：L1 抽取在上一轮结束后才写库，下一任务必须看见它；boot 快照会让模型永远晚一个任务。
- 行格式：`slot: content`（如 `fact.family: mother birthday 05-06`）。`slot` 是 D20/D35 的枚举列（`pref.tone`/`habit.work_hours`/`fact.family`…），也是用户在隐私页逐行删的那把键；`ListProfiles` 自身按 `updated_at DESC, id DESC` 返回（`internal/memory/dao_profile.go:126-129`）。
- 读失败：**降级但可见**——`rt.auditf("wisp run: L1 画像读取失败，本轮 prompt 不带画像（不是没有画像）：%v")`（`agent_profiles_451.go:79`），本节为空，任务照跑；不静默假装"没有画像"。

**接线后请求体真的变了**：同一枚测试在接线后 PASS，日志 `.scratch/wisp/probes/451/r1/ac1-green.txt`：

```
--- PASS: Test451L1ProfileReachesTheProviderPrompt (2.65s)
--- PASS: Test451L1ProfileInjectionStaysInsideTheD15Budget (4.19s)
```

**三条禁令核对**：
- ⛔ 没有在 `internal/agent/loop.go:677` 那个转发点就地造数据：`loop.go` 本波 diff 为空，转发保持原样，数据只从库里来；
- ⛔ 没有删 `Profiles` 字段（见 §①）；
- ⛔ 没有裸拼接：注入面是 `[]string` → `agent.Config.Profiles`，字符串拼接只发生在 assembler 的预算渲染里（§④）。

---

## ④ AC#3 预算：走的是 D15 既有的那条路（无裸拼接）

**链（逐跳 file:line）**：

```
loopProfileLines451()                 cmd/wisp/agent_profiles_451.go:71   ← 只读行、只渲染 "slot: content"
  -> agent.Config.Profiles            cmd/wisp/run.go:2123               ← 接线（本票唯一新增赋值）
  -> PromptInput.Profiles             internal/agent/loop.go:677         ← 既有转发，未动
  -> profileText(..., SectionProfile) internal/agent/prompt.go:164, :297  ← D39 (3) 段预算
  -> clipToTokens(..., Budget)        internal/agent/prompt.go:340        ← 单段硬裁
  -> enforceTotal(..., PromptTotal)   internal/agent/prompt.go:394        ← D39 总上限
```

**预算数**（`internal/agent/budgets.go`，按端点上文窗口缩放）：`refSectionProfile = 400`（`:27`，128k 参考窗口）、`refPromptTotal = 2300`（`:30`）。

**超预算时丢哪一条（具名策略）**：`ListProfiles` 新→旧返回，`profileText` 从头逐行加、**遇到第一行加不进就 `break`**（`prompt.go:304-312`）⇒ 掉的是**列表尾部＝最旧的行**，而且是**整行掉**（不会半行）；`enforceTotal` 只在"前缀块自己就把总预算挤爆"时才进一步压段，那时按既有 D39 行为由 `clipToTokens` 对整段加截断标记（本票未改这段）。

**读数（不是形容词）**：`Test451L1ProfileInjectionStaysInsideTheD15Budget` 预置 **20 条**（DAO 上限）长行，跑同一套真组合根后从 mockllm 捕获体里数：

```
agent_profiles_451_test.go:130: 预算路径读数：画像段保留 11/20 条（D15 (3) 段预算把它裁到预算内；裸拼接会是 20 条）
--- PASS: Test451L1ProfileInjectionStaysInsideTheD15Budget (5.92s)
```

⇒ 预算**真的咬了**（20 → 11），裸拼接会看到 20 条。日志 `.scratch/wisp/probes/451/r1/ac3-budget.txt`。

**边界决定（写明，不夹带）**：本票**没有**按 `[memory] l1_enabled` 闸断注入。该键在册语义是"抽取管线不挂、画像不学习"（`cmd/wisp/config_readers_255.go:257`、`cmd/wisp/memory_retention_180k_windows.go:10-25`），它**不删行**；用户的"忘掉"路径是隐私页的 purge（`internal/memory/privacy.go`、`PurgeProfiles`）。⇒ 行集即注入集。**残余留给裁决者**：若要"关掉 L1 连使用一起关"，是在 `execute()` 同一处加一个 `cfg.Memory.L1Enabled` 条件的**一行**改动，属产品语义选择，不属本票"把断线接上"。

---

## ⑤ AC#4 红→绿 ＋ 定向突变 M1

| 状态 | 命令 | 读数 |
|---|---|---|
| 接线前 | 载具 `-test.run '^Test451L1ProfileReachesTheProviderPrompt$'` | **红**，`CARRIER_EXIT=1`，红句＝§② 逐字（`.scratch/wisp/probes/451/r1/ac1-red.txt`） |
| 接线后 | 同上 | **绿**（PASS 2.65s；`.scratch/wisp/probes/451/r1/ac1-green.txt`） |
| **突变 M1** | **只删** `cmd/wisp/run.go:2123` 那一行（`Profiles: rt.loopProfileLines451(),`；注释保留），重跑同一枚测试 | **红**，逐字同一句；`CARRIER_EXIT=1`（`.scratch/wisp/probes/451/r1/mut-m1-red.txt`） |

**恢复的逐字节证据**：突变前把接线版存为 `.scratch/wisp/probes/451/r1/run.go.wired`；恢复＝拷回。两端 `git hash-object` **同为 `ad9ab9c6181371539bd7f07fb8cc417277c1a5da`**（见 §⑧ 命令）；恢复后整包跑里两枚测试 PASS（§⑥）。

⇒ **AC#4 判据成立**：`摘掉接线那一行 ⇒ Test451L1ProfileReachesTheProviderPrompt 红`（不是别的枚红，也不是"全包都红"）。

---

## ⑥ AC#5 门禁（本机现跑，全过）

| 门 | 命令 | 读数 |
|---|---|---|
| gofumpt | `$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp/agent_profiles_451.go cmd/wisp/agent_profiles_451_test.go cmd/wisp/run.go` | **空**（0 行），rc=0；`-l internal/agent internal/memory` 也空 |
| go vet | `go vet ./cmd/wisp/`；`go vet ./internal/agent ./internal/memory` | rc=0（无输出） |
| d22scan | `scripts/d22scan.sh` | **control rc=0，scan rc=0**；`ban #8 internal/=1017`、`ban #8 cmd/=291`（含新件） |
| path-length | `scripts/check-path-length-budget.sh` | **VERDICT GREEN**，rc=0（over-hat=68、wall interval=0、roster==tree） |
| 单包测试 | `go test -count=1 ./internal/agent ./internal/memory` | `ok internal/agent 21.896s`、`ok internal/memory 46.227s` |
| 整包（载具） | `scripts/wisp-cli-staged-tests.sh`（无 `-test.run`） | **RUN=866 / PASS=623 / FAIL=3 / SKIP=17**，`1394.00s`；FAIL 三枚见下 |

**整包 FAIL 三枚，逐枚归因**：

1. `TestAC4FocusReturnToPriorWindowGap33r5`（5.65s）——**派单具名的本机已知红**（驻留球 PID 49348 在跑），⛔ 未修/未跳。
2. `TestAC246DevLegIgnoresTheTestTaskInjection`（37.23s）——**派单具名的本机已知红**（同上），⛔ 未修/未跳。
3. `Test320FetchLegWritesTheKeyNowhereOnDisk`（0.80s）——**整包内互扰，不是本票的**：红句 `panel_models_320_test.go:523: walk C:\Users\swq\AppData\Local\Temp: open ...\Test192PreviewRefusalKeepsTheFallbackPolicy...: The system cannot find the file specified.`（它在扫 %TEMP% 的窗口里，另一个并行测试恰好删了自己的 tempdir）。**定向单跑 PASS 0.16s**（`.scratch/wisp/probes/451/r1/t320-standalone.txt`），且台账 `docs/reports/pending-and-issues.md:35852` 早已把同一枚记成"整包内互扰（定向单跑 PASS 0.10s）"。

**本票两枚测试在整包里也 PASS**（同一跑：`--- PASS: Test451L1ProfileReachesTheProviderPrompt (3.49s)`、`--- PASS: Test451L1ProfileInjectionStaysInsideTheD15Budget (3.51s)`）。

---

## ⑦ 改面／交付核

- 文件：`cmd/wisp/agent_profiles_451.go`（新，91 行）、`cmd/wisp/agent_profiles_451_test.go`（新，139 行）、`cmd/wisp/run.go`（+11）。
- commit：`93a1c93da`（branch `dev`，**未 push**）；`git show --stat`：恰 3 枚文件，+241/−0。
- 显式 pathspec：`git add -- <3 条>`、`git commit -- <3 条>`；⛔ 无 `-A`/`add .`/`--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`，未删文件。
- ⛔ 未动：`thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43 表／`frontend/**`／`design/**`；⛔ 未放宽任何断言。

## ⑧ 证据日志与可重跑命令

| 件 | 是什么 |
|---|---|
| `.scratch/wisp/probes/451/r1/ac1-red.txt` | 接线前 AC#1 红（含 mockllm 捕获体） |
| `.scratch/wisp/probes/451/r1/ac1-green.txt` | 接线后 AC#1＋AC#3 绿 |
| `.scratch/wisp/probes/451/r1/ac3-budget.txt` | 预算读数 11/20 |
| `.scratch/wisp/probes/451/r1/mut-m1-red.txt` | 突变 M1 红（同一句） |
| `.scratch/wisp/probes/451/r1/run.go.wired` | 突变前接线版（恢复源，hash `ad9ab9c6…`） |
| `.scratch/wisp/probes/451/r1/cmdwisp-full.txt` | 整包载具全量日志（1394s） |
| `.scratch/wisp/probes/451/r1/t320-standalone.txt` | Test320 定向单跑 PASS（归因用） |
| `.scratch/wisp/probes/451/r1/gates.txt` | path-length＋d22scan 日志 |

**逐字节恢复的复跑尺**：
`git hash-object cmd/wisp/run.go` 与 `git hash-object .scratch/wisp/probes/451/r1/run.go.wired` 两端同值（本波读数均 `ad9ab9c6181371539bd7f07fb8cc417277c1a5da`）。

## ⑨ 诚实边界（残余，不粉饰）

- 本腿**只**证明"库里的行会进 prompt、且走 D15 预算"；**不**主张 L1 抽取质量、去重、隐私页 UI 的完备性（那是别的票）。
- 子代理 loop 拿到的是**父任务那次读的画像快照**（`opt := base`，`internal/tools/subagent_197.go:321`），不是每个子任务重读；同一会话内一致，跨任务以根任务读为准。
- 未跑 `-race`、未跑 winlive 专属包；整包只按派单走载具一次。
- 撤销口令：**「451 撤」＝删本票与其证据件（不 reset/amend）**；产码回退走**追加 revert**（`git revert 93a1c93da`，不改历史）。
- 证据件本枚 commit 的 hash 由编排者收件时现读（本件不自引自身 hash）。
