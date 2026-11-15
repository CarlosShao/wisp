# 票 451 — L1 画像抽取写着数据库，**却永远进不了 prompt**：`agent.Config.Profiles` 全仓零赋值

**状态**：待派（产码缺陷，S4 验收时由验收腿现读发现，编排者已复核）

## 现量（可逐字重跑）
`grep` 产码（非测试）里 `Profiles\s*[:=]` 的全部命中 = **3 处**：
- `internal/agent/loop.go:677` → `Profiles:  l.opt.Config.Profiles,` —— **这是转发，不是赋值**
- `internal/memory/dao_profile.go:130` → 错误串里的 `list profiles`（无关）
- `internal/memory/models.go:27` → `const L1MaxProfiles = 20`（无关）

⇒ **仓库里没有任何一处构造 `agent.Config` 时给 `Profiles` 赋过值**。`ListProfiles` 只被 privacy/undo 读。

## 后果（这才是这张票的理由）
「记住第 1 轮说过的事实」这条产品承诺**在真机上不成立**：抽取、写库、去重都是真的，库里有行；但**没有任何代码把画像递进 prompt**，模型看不到它。⇒ 用户会看到"它记了，但它不提" —— 这比"没记"更糟，因为它**看起来在工作**。

## 为什么以前没人发现
票 29 的 `Key constraints:22` 写了"注入"，但它的 **5 格 AC 没有一格钉住注入**（验收腿逐字核过）⇒ 票 29 `-done` 是**真的**，同时又**是空的**。这正是"每格 AC 都过、切片却没成"的形状（台账 A1875）。

## AC

- [x] **AC#1 先证明"今天真的断"**：写一枚测试，构造一份带 L1 画像的库，跑一轮对话，断言**送给 provider 的 prompt 里含那条事实**。⛔ 先改产码再看 —— **它必须先红**，红句要能读出"prompt 里没有这句话"。
- [x] **AC#2 接线**：在**唯一的装配点**把检索到的画像喂进 `agent.Config.Profiles`（⛔ 在 `loop.go:677` 那个转发点就地造数据；要写清"谁在建 Config、谁负责查库"）。给 `grep` 前后计数：`Profiles:` 赋值点 **0 → ≥1**。
- [x] **AC#3 预算与边界**：画像进 prompt 必须走 D15 上下文预算那条路（⛔ 裸拼接、⛔ 绕过预算），并给"超预算时丢哪一条"的具名策略。
- [x] **AC#4 可判事实**：AC#1 那枚测试**由红转绿**；并给一枚**定向突变**（摘掉接线那一行 ⇒ 那枚测试红）。
- [x] **AC#5 门禁**：`gofumpt` 空／`go vet` rc=0／`d22scan` clean／path-length GREEN／`go test ./internal/agent ./internal/memory` 全绿＋整包（跑不完如实写）。

## 边界
⛔ 改 `thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43 转移表／`frontend/**`／`design/**`；⛔ 放宽既有断言；⛔ "把 `Profiles` 字段删掉就算修好"（那是把承诺删掉，不是把功能接上）；⛔ add -A/--amend/reset/rebase/stash/checkout .；⛔push。
写面：`internal/agent/**`、`internal/memory/**`、`cmd/wisp/**`（**每波一枚写腿**）、`docs/evidence/s1/451-*`。

## 撤销口令
「451 撤」＝删本票与其证据件（不 reset/amend）。

## Progress
- （待派）
---

## 收件（2026-11-15 13:0x，`93a1c93da` 接线 ＋ `1c19b1953` ＋ `e2717d27e`）：5/5，结案

**尺1**：三条提交在。产码落点＝新件 `cmd/wisp/agent_profiles_451.go`、新测试 `cmd/wisp/agent_profiles_451_test.go`、改 `cmd/wisp/run.go`。
**尺2**：证据 `docs/evidence/s1/451-profiles-to-prompt.md` **13770 B**；底片 `.scratch/wisp/probes/451/r1/ac1-red.txt` **24192 B**（接线前的红）与 `mut-m1-red.txt` **24191 B**（突变后同句红）。
**尺3 我自己复跑**：`Test451L1ProfileReachesTheProviderPrompt` ⇒ **`RUN=1 PASS=1 FAIL=0`，2.40s**。
★**我自己核了赋值点**：`Profiles\s*[:=]` 产码命中现为 4 处，其中**唯一真赋值**是 `cmd/wisp/run.go:2123` → `Profiles: rt.loopProfileLines451(),`（`internal/agent/loop.go:677` 仍是转发）⇒ **0 → 1**。

### 顺序没倒（这是本票最要紧的一条，我逐条对上）
**接线前红句逐字**：`agent_profiles_451_test.go:77: 送给 provider 的 prompt 里没有画像段（找不到 "## profile"）：库里那枚 L1 画像 "mother birthday 05-06" 没有被递进请求，模型看不到这条事实。`
测试用的是**产品自己的 DAO**（`memory.Open`＋`UpsertProfile`，slot `fact.family`）写进**这次 run 会打开的那个 data dir**，跑**真组合根** `runTextTask`（config→secret→provider→gate→bridge→真 loop），读 **mockllm 捕获的请求体** ⇒ ⛔ 不是自己造一份 prompt 字符串自评。

### 三条禁令全守（我验）
1. `internal/agent/loop.go` 转发点 diff **空** ⇒ ⛔ 没在转发点就地造数据
2. `Profiles` 字段保留 ⇒ ⛔ 没把承诺删掉
3. 预算走 **D15**：`Config.Profiles → PromptInput.Profiles(loop.go:677) → profileText(SectionProfile, prompt.go:297) → clipToTokens → enforceTotal`，预算 `refSectionProfile=400`／`refPromptTotal=2300`。**丢法具名**：新→旧排序、遇加不下即 break ⇒ 整行丢掉最旧的尾部。读数：**20 条长行只进 11/20**（裸拼接会 20 条）⇒ 预算是**真的在管**。

### 一处值得记的手法（不是取巧，是相反的）
它追加了一枚**只改注释措辞**的提交（`1c19b1953`）：原注释里出现 `Profiles:` 字面会把普查尺数成 1→3，它**把注释改词**让尺回到 1→2，**⛔ 没有去放宽那把尺**。⇒ 这正是我要的方向。

### 整包（一发，载具，1394s）
`RUN=866 / PASS=623 / FAIL=3 / SKIP=17`。三枚红**逐名**：两枚派单具名已知红＋**第三枚 `Test320FetchLegWritesTheKeyNowhereOnDisk`** —— 属**整包内互扰**（它扫 `%TEMP%`，撞上并行测试删 tempdir），**定向单跑 PASS 0.16s**。⇒ 三枚全**未修未跳**，第三枚已入红名册（台账 `:35852` 记过同款形状）。

### 残余（腿写明，我认，不夹带）
1. **`l1_enabled` 没有闸注入**：该键在册语义是"抽取管线不挂、画像不学习"，不删行 ⇒ 行集即注入集。若裁决要"关掉 L1 连使用一起关"＝同处加一个条件的一行，属**产品语义选择**，不属本票接线。
2. 子代理 loop 拿父任务那次读的画像快照（`opt := base`），同一会话内一致；**未跑 `-race`**。
⇒ 第 1 条我记入台账待办，不阻塞本票结案（本票 AC 没要求它）。
