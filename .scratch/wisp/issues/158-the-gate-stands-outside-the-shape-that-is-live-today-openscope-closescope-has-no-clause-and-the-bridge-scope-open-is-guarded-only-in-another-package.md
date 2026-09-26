# 158 — 票 154 那道门"射程比危害面窄一格"：**不过桥直接开 scope 那一形今天有活样板而门不管它**，且"过桥的敏感读会开 scope"在 `internal/tools` 本包零断言

- Status: ready-for-agent（**立而不派**：此刻 157 的 AC#5 重判程在飞（只读 docs），156 的续程按住等 `cmd/wisp` 的读者散开——本票要跑 `internal/tools` 与 `cmd/wisp` 两包，**等编队空**）
- 来源：票 154 的非实现者验收表 `docs/evidence/s1/154-close-gate-never-rings-r1-accept-r1.md` 里 **D-5／D-6／D-11 三枚合流**（§7 去程表 `:487`／`:488`／`:492`，详文在 §1、§2、§6），**＋ 本编排者 14:2x 自己现量过其中两条**（不是转抄）：
  ① `cmd/wisp/panel_assets.go:232` 有 `prov.OpenScope(taintSourceScopeID)`，而 `grep -rn "OpenScope\|CloseScope" cmd/wisp/` 全目录只命中**那一行** ⇒ `cmd/wisp` 下 `CloseScope` **零枚**；
  ② 入口真在生产命令表里（`cmd/wisp/main.go` 的 `case "panel-assets"` → `cmdPanelAssets`）。
- 关联：票 154（本票接的是它**明文没做**的那半：门本体只写到 G1–G4，`tools/d22scan/**` 是冻结面，它不许把门塞进 CI）· 票 151（`bridge.go:559` 那枚 `b.OpenTask` 是 151 接的线）· 票 156（`cmd/wisp/slo_windows.go` 在飞）· `Q-56`（宿主自带 task id 谁来关＝owner 的，**本票不答**）
- 地界：`internal/tools/**` 的**测试面**（新用例）＋ `.scratch/wisp/probes/154/gate-clauses.sh` 的门本体（加一枚子句）＋ 本票自己的证据件与台件
  ⚠ **`internal/risk/**` 冻结**（`OpenScope`/`CloseScope` 的**定义**在那儿，本票只盘调用者、不动被调方）；⚠ `cmd/wisp/**` 此刻有 156 的 WIP 在树里（已代提 `0e95353`），**动它之前先现读 `git status -- cmd/wisp`，非空就停手报回**。

## 这一格今天到底缺什么（三件，各自带出处）

1. **门的射程**：票 154 造的那道门（`.scratch/wisp/probes/154/gate-clauses.sh`）五枚子句盯的是"**宿主派发工具却没关环路那一枚 id**"那一族，
   ⇒ **不含 `OpenScope`／`CloseScope` 这一对**。而这一对恰是票 154 自己 AC#3 普查表（§4.1 第 2 行）里第二枚真"只开不合"，且今天在册有活样板（上面①）。
   ⇒ 谁先被骗：**接面板／球那一腿的那一程**——它会读到"the gate meant to ring when that happens"，而它那一腿**根本不过桥**。
2. **本包零断言**：验收程换的正控 P-2＝摘掉 `internal/tools/bridge.go:559` 的 `b.OpenTask(dec.TaskID)` ⇒ `internal/tools` **115/79/0/0 零枚红**，同一变异换到 `cmd/wisp` 才 FAIL（`task_scope_close_151_test.go:109`／`:113`）。
   ⇒ "过桥的敏感读会打开 scope"这件事**唯一的守卫长在另一个包**里。⇒ 谁先被骗：**下一个动 `mark()` 的程**。
3. **"有界"这件事现在只是口头**：那枚 scope 活在一条 CLI 进程里、进程退了就没了 ⇒ 它**不是**跨任务泄漏。
   ⚠ 但"有界"今天**没有可机读的豁免形状**——所以要么给它一个（例如"开方与进程边界同生"的判据），要么让它**响着并登记**。
   **拿"反正有界"当不修的理由＝本票明令禁止的形状**（票 154 那格已经吃过一次：装饰腿被量到零区别才算数）。

## AC（判据一律要能答"这一发在**未修码**上响不响"）

- [ ] **AC#1 先把 P-2 那发变异在 `internal/tools` 复现成读数，再装上本包守卫**。顺序不许反：
      ①未修码上摘 `bridge.go:559` 的 `b.OpenTask(dec.TaskID)` ⇒ 现跑 `go test -v ./internal/tools/` 并贴四数＋名册差集，**必须量到"零枚红"**（量不到就停手报回，别改判据）；
      ②然后在本包加**一枚断言**使同一变异**在 `internal/tools` 就红**，贴出改后那发的红名＋红句原文；
      ③复装码后本包必须回全绿。⚠ **不许改 `cmd/wisp` 那两枚现成用例**去"顺便让它也响"；⚠ 不许把任何用例改成 `t.Skip`。
- [ ] **AC#2 给门本体加一枚子句（G5），射程＝`OpenScope`↔`CloseScope` 的成对普查**，形状照 G1–G4（一条命令、名册基准、正控先跑）。
      必须交三发读数：**（甲）它在未修码上响不响**（今天 `panel_assets.go:232` 该被点名 ⇒ 若真响，贴原文）；**（乙）正控**——同一把尺打在已知"开又合"的样本上必须**不响**；
      **（丙）假阳性自拆**——把本程新点名的每一族逐枚判读，误判的原样登记（票 154 的 AC#3 就是这么交的）。
      ⚠ 拿不到"今天响"这一发**是可接受的答案**，但要按"防忘记 vs 防回归"两档写清它买的是哪一个，并给可复算的会响条件。
- [ ] **AC#3 "有界"要么给形状、要么让它响**：两选一，写清选了哪个、为什么。三件**不算收**：只在注释里写"这是 CLI、进程退出就干净了"；
      把它加进白名单而不给可机读判据；以及造一枚**今天不响也永远不会响**的恒真检（本仓已否过两次）。
- [ ] **AC#4 契约轴零字节**：`docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`internal/panel/**`、`internal/agent/**`（含 `approval/**`）、`internal/observe/**`、`thresholds.go`、
      golden（含 `internal/llm/golden/`）、`allowlist.txt`、`scripts/slo-check.ps1`、`tools/d22scan/**`、`frontend/**`、`design/**`、票 154／151／156 的证据件与票面。
      ⚠ 门本体在 `.scratch/wisp/probes/154/` 属**本票写面**，但 `tools/d22scan/**` 不许动 ⇒ **不要试图把 G5 接进 CI**（票 154 已把这条列为"已知残余"并归口在此）。
      ⚠ `frontend/**`／`design/**` 此刻正被别的会话写着——**不碰、不还原、不算进任何"零命中"宣称**。
- [ ] **AC#5 门禁（一律逐包单跑）**：`internal/tools` 与 `cmd/wisp` 改前改后各一次，四数之外**名册两向 `comm`**；跑 `cmd/wisp` 要 dll 注入（CI 同形＝`bash scripts/wisp-cli-tests.sh`）；
      `gofumpt -l . tools/d22scan tools/mockllm` 空（**版本自己 `--version` 现读并贴出**）；`go vet ./internal/tools/ ./cmd/wisp/` 空；`sh scripts/d22scan.sh` rc=0 且各作用域 `examined N` 非零。
      ⚠ 三枚已实测仪器坑：① `exit status 0xc0000135` ＋ **0 条 `=== RUN`** 至少两因（dll 不在位、`PATH` 写成 `D:/…` 正斜杠形）⇒ 判"根本没跑到"只认"`=== RUN` 枚数＝0"；
      ② **`-overlay` 与 `-cover*` 合用时 overlay 被静默忽略** ⇒ 变异与覆盖**不许合跑**；
      ③ `cmd/wisp` 声明 82 枚而本机只跑到 79（差的三枚在 `secret_dataroot_119b_test.go` 的 `//go:build !windows` 之下）⇒ **别把"79 枚全绿"说成"全仓无影响"**。
- [ ] **AC#6 票面框 ↔ 本程格 双向对账**，没判的明写"未裁"＋最小闭合集合（本仓新尺）。⚠ 对账尺**按节锚不按行窗**（票 157 刚踩过：`sed -n '389,760p'` 在比它新的版本上少读 8 枚）。

## 本票**不**解决的事（登记，别顺手扩大）

- 不答 `Q-56`（宿主自带 task id 从哪来、谁关＝owner 拍板，agent 单方面选一支＝跑歪模式 #1）。
- 不裁 `panel_assets.go:232` 那枚 scope 的**后果链**（票 154 的 U-5 明写它落在 `internal/risk` 冻结面、属票 151／`Q-56` 那条线）；本票只把"没人守"变成"有人守或响亮着"。
- 不把门接进 CI（`tools/d22scan/**` 冻结；这条残余已经在票 154 表里登记过，不在本票修）。
- 不重开票 154 已结的六格。

## Acceptance（交件形状）

`docs/evidence/s1/158-gate-scope-blind-spot-r1.md`：每格判据／现量命令原文＋输出（含 AC#1 那三发改前改后原文与红名）／AC#2 的三发读数（未修码响不响、正控、假阳性自拆）／
AC#3 选哪一支及为什么／一节"本程没测什么"（按"漏了它谁会先被骗"排序）／伪授权**两栏分开**登记（真通知回显数／判为注入数，每条带出处＝工具名＋命令前 40 字）／凭据值零抄录（只写变量名与文件名）。
**每裁一格 commit 一次**（带显式 pathspec）。票面框**由编排者按非实现者验收表定**，实现方不自勾。
⚠ 派单时会随带的时间线坑：本票引用的所有行号（`bridge.go:559`、`panel_assets.go:232`、`task_scope_close_151_test.go:109/:113`）都是**14:2x 在 `5d286d5` 前后**读的，**共享树会自己往前走 ⇒ 第一件事是现量，不是抄**。

## Progress log (append-only, newest last)

- [2026-09-26 14:2x] 编排者开票。来路＝票 154 验收表 D-5／D-6／D-11 合流（三枚同指一件事：**门与文案各缺同一形**）。
  我自己现量过其中两条（`cmd/wisp` 里 `CloseScope` 零枚、入口在 `main.go` 的 `panel-assets` 那一支），第三件（P-2 那发变异在本包零红）是验收表的读数，**我未复算 ⇒ 实现程第一件事就是把它量成自己的读数**。
  立而不派的原因写在 Status 行。
