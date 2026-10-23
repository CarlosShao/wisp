# 370 — `Assembler.MaxOutputTokens` 全仓无赋值 ⇒ `Request.MaxOutputTokens` 恒 `0` ⇒ `StatusTruncated`**永远到不了**

- Status: **ready-for-agent**
- 立案：2026-10-23 18:0x +08 编排者立
- ★**来路（实现者签名、只报不修）**：`365-r1` §9-2（件 `docs/evidence/s1/365-r1-impl.md` **150 行**、`6b9a83c04`）
- 关联：票 **365**（`taskLogState` 统一）· `internal/agent/prompt.go:132` · `internal/llm/resolver.go:205` · `cmd/wisp/**`

## 这是什么（人话）

**「输出被天花板截断」这个状态（`StatusTruncated`）今天不可能发生** —— 因为**那个天花板从未被传给模型**。

★**链条**（`365-r1` 现跑）：
1. ★**配置侧 OK**：`[llm.roles.chat].max_output_tokens` **已经落进 `llm.Endpoint`**（`resolver.go:205`）。
2. ★★**中间断了**：★**`agent.Config` 没有承接字段**；★**`agent.Assembler.MaxOutputTokens`（`prompt.go:132`）**全仓无赋值****。
3. ★★**结果**：★**`Request.MaxOutputTokens` 恒 `0`** ⇒ **供应商拿不到上限** ⇒ **无从"撞天花板"**。

⇒ ★★**所以 `StatusTruncated` 在生产里是**不可达状态****。

## ★★为什么这值得单独修（⛔ 并入 365）

- ★**票 365 治的是"同一状态写两个词"**；★**本票治的是"那个状态根本到不了"**。
- ★★**而 `365-r1` 现在把 `StatusTruncated` 映射成 `interrupted`（权威词）** ⇒ ★**这条映射**今天是一条死路****：★**它写得对、但永远走不到**。
- ★★**这正是本仓反复出现的"能力造出来了、产品里没人用"那一族**（票 166／132／130／114 **同形**）—— 只是这次**断在"配置没传下去"**。

## 结案判据

- [ ] **AC#0 现读＋归因（只读，⛔ 动任何产码字节）**：
      ①★**现跑** `git grep 'MaxOutputTokens'` ⇒ ★**逐处列出**（★`prompt.go:132` 的字段、★`Request` 的字段、★**有没有任何赋值**）。
      ②★**现跑** `resolver.go:205` ⇒ ★**配置值确实进了 `llm.Endpoint` 吗**（★给逐字）。
      ③★**量出整条链**：★`config → llm.Endpoint → ?(断点) → agent.Config → Assembler → Request` ⇒ ★**给一张"哪个环节断了"的表**。
      ④**正控**：★**造一发"手工把 `Request.MaxOutputTokens` 设成 `1`"** ⇒ ★**必须能撞出 `StatusTruncated`**（★**证明"到不了"是配置没传、⛔ 状态机缺**）。
- [ ] **AC#1 裁形（编排者拍）**：三形各写代价＝
      〔甲〕**把 `MaxOutputTokens` 一路传到 `Request`**（★**最小、但要看 `agent.Config` 加字段会不会撞契约**）；
      〔乙〕**在 `Assembler` 里从 `llm.Endpoint` 读**（★**可能造 `agent`→`llm` 的新依赖**）；
      〔丙〕**删掉 `StatusTruncated` 及其映射**（★**承认这个状态不存在** ⇒ ★**但那会抹掉一个真实的上限语义** ⇒ ★**⛔ 推荐**）。
      ★**未裁形任何实现都⛔ 开工**。
      ★★**注意**：★**`cmd/wisp` 每波只许一枚写腿**；★**若必须动 `cmd/wisp` 组装 ⇒ 与别的 `cmd/wisp` 腿**串行**。
- [ ] **AC#2 落地＋牙（写腿）**：按裁定那一形做；★**每处改动配一枚今天响、不修不响的用例**；
      ★**反向判据**＝把新加的那一跳摘掉 ⇒ 用例必须**转红并逐名报红行**；三读数 ＋ **还原核**（★**时机必须在最后一发变异之后**）。
      ★★**验收判据**：★**造一发"输出撞天花板"** ⇒ ★**库里的 `task_log.state` 必须是 `interrupted`**（★**⛔ 只测"字段有值"**）。
      ★★**配方可借 `365-r1` 那条**：★**用**配置**造真状态、⛔ 构造 `Result` 冒充**。
- [ ] **AC#3 门禁**：`GOFLAGS= go build ./...`／`go vet` 动过的包／`go test -count=1`（**四数分开、⛔ 相减**）／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh`。
      ★**格式门禁甲乙两尺都报**（甲＝入库 blob、**判据**；乙＝工作树、**只作诊断**）。
      ★★**`cmd/wisp` 的测试必须带 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`** —— **不加就 `0xc0000135` ＝ 假失败**；
      ★★**但整包 `cmd/wisp` 会在 `MARK299 before CreateWindowExW` / `before Embed` 被硬终止** ⇒ **整包四数 ⛔ 当判据**、只跑定向 `-run`。
- [ ] **AC#4 纪律格**：⛔ 翻任何框、⛔ 改票名／`Status:`／`-done`、⛔ push、⛔ 碰票 365 任何原句。

## ★硬提醒（⛔ 踩）

1. ★★**⛔ 加 `truncated` 进那 12 枚名册** —— ★**票 365 的硬边界 2**：★**加词 ＝ 改名册 ＋ 升 schema ＝ 另一件事**（★而且 `196-c1`/`196-c2`/`196-c3` 刚把名册与契约文字对齐）。
2. ★★**⛔ 为了"造出 `StatusTruncated`"而放宽 `MaxOutputTokens` 的语义** —— ★**那是把门挖掉一块**。
3. ★**若 `agent.Config` 加字段 ⇒ 看它是不是 `C` 级契约面**（★那是**改契约**、须机主批）。

## 排程

`370-a1`（只读＝`AC#0`）→ 编排者裁 `AC#1` → `370-r1`（写腿＝`AC#2`）→ `370-v1`（非实现者裁）。

## Progress log

- [2026-10-23 18:0x +08] agent=`编排者` did=立案（★来路＝`365-r1` §9-2；★它现跑量到：`[llm.roles.chat].max_output_tokens` 已进 `llm.Endpoint`、但 `agent.Config` 无承接字段、`Assembler.MaxOutputTokens` **全仓无赋值** ⇒ `Request.MaxOutputTokens` 恒 0 ⇒ `StatusTruncated` 不可达；★它**只报不修**（超写面））next=派 `370-a1` 只读腿。
