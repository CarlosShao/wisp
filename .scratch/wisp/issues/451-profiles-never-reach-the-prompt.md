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

- [ ] **AC#1 先证明"今天真的断"**：写一枚测试，构造一份带 L1 画像的库，跑一轮对话，断言**送给 provider 的 prompt 里含那条事实**。⛔ 先改产码再看 —— **它必须先红**，红句要能读出"prompt 里没有这句话"。
- [ ] **AC#2 接线**：在**唯一的装配点**把检索到的画像喂进 `agent.Config.Profiles`（⛔ 在 `loop.go:677` 那个转发点就地造数据；要写清"谁在建 Config、谁负责查库"）。给 `grep` 前后计数：`Profiles:` 赋值点 **0 → ≥1**。
- [ ] **AC#3 预算与边界**：画像进 prompt 必须走 D15 上下文预算那条路（⛔ 裸拼接、⛔ 绕过预算），并给"超预算时丢哪一条"的具名策略。
- [ ] **AC#4 可判事实**：AC#1 那枚测试**由红转绿**；并给一枚**定向突变**（摘掉接线那一行 ⇒ 那枚测试红）。
- [ ] **AC#5 门禁**：`gofumpt` 空／`go vet` rc=0／`d22scan` clean／path-length GREEN／`go test ./internal/agent ./internal/memory` 全绿＋整包（跑不完如实写）。

## 边界
⛔ 改 `thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43 转移表／`frontend/**`／`design/**`；⛔ 放宽既有断言；⛔ "把 `Profiles` 字段删掉就算修好"（那是把承诺删掉，不是把功能接上）；⛔ add -A/--amend/reset/rebase/stash/checkout .；⛔push。
写面：`internal/agent/**`、`internal/memory/**`、`cmd/wisp/**`（**每波一枚写腿**）、`docs/evidence/s1/451-*`。

## 撤销口令
「451 撤」＝删本票与其证据件（不 reset/amend）。

## Progress
- （待派）
