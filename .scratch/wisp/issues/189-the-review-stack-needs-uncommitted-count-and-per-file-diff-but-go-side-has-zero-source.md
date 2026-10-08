# 189 — 面板"审查"那一堆要真东西：**未提交改动枚数**与**逐文件 diff** 今天 Go 侧零来源（owner 09-28 放权：缺就补）

- Status: **ready-for-agent（先只读设计核，再落地）**。⚠ **功能票**，且**与票 181 同一条读面**（复用，别重造）。
- 来路：owner 09-28 15:2x 原话（`A373`）＋截图里右栏"审查／看这一轮改了哪些文件、每一处改动是什么"。普查出处＝`docs/evidence/s1/180-182-panel-fields-census-c1.md` §5（`git diff`／`numstat` 全仓 **0 命中**）。
- 关联：票 181（git 只读面，本票复用它的读面）· 票 186（切换，同族不同物）· 票 182（那一栏的堆数）· 票 77（前端已定案"右栏三标签＝审查／终端／浏览器"）

## 现量（`180-c1＋182-c1` 现跑，编排者复跑同值）
- `grep -rniE "git diff|numstat" --include=*.go internal/ cmd/`（排 `_test.go`）⇒ **0 命中**；票 181 的普查也独立量到"Go 侧零 git 读面"。
- 截图里那一栏（改了哪些文件／每一处改动是什么）**今天完全无源**，前端是示意。

## AC
- [ ] **AC#1 只读设计核**：未提交枚数与 diff 各从哪儿取——`git status --porcelain` 等价读法要不要起外部进程？（⚠ 票 181 已证"认 git／工作树／分支"三样**不需要**起进程；diff 这一样**很可能需要**，因为内容比对不在文件索引里。）⇒ 逐支给现量与代价，并核 `tools/d22scan/main.go` 的 ban 列表原文（`181-c1` 已给结论：七枚 ban **没有一枚管起外部进程**，真要跑会咬到 #1 裸 goroutine／#4 墙钟超时／#2 路径决策过 C26 三条＋D38(b) 的 goroutine 名册）。
- [ ] **AC#2 只读落地**：宿主侧读面（**与票 181 同一枚文件族，不另起一份实现**——两份读面必然漂），送进快照；判据钉"值来自真仓"（在临时目录里 `git init` 一棵合成树，改几个文件，断言枚数与文件名集合正确）。
- [ ] **AC#3 ⚠ 只读边界（这一格是本票的命门）**：交付里**不许出现**任何"提交／推送／接受／回滚"的动作腿。⇒ 普查 `180-c1＋182-c1` 把"diff 接受／回滚"与"面板 git 提交／推送"列为**雷区＝要 owner 逐枚批**；本票只做显示，写动作另立票、且要过 L2 与审计。补一枚常驻判据钉住"审查面只读"。
- [ ] **AC#4 大 diff 的代价**：整棵工作树的 diff 可能巨大——要有**字节上限与截断语义**（照票 164／D15 那根管子的教训：**截断必须带"怎么找回全量"的指针**，否则又是"做了但没生效"）。给一发真读数。
- [ ] **AC#5 契约轴**：D34 那张内置工具权威表**不许加 git 行**（做了就是"模型可调用 git"，那是契约变更）；`docs/PLAN.md`／`docs/specs/**`／`allowlist.txt`／`thresholds.go`／golden 一字节不动；C17 一枚不加。
- [ ] **AC#6 门禁**：`./internal/panel/ ./internal/tools/`＋`sh scripts/d22scan.sh`（基线 433）＋`gate-clauses.sh` 比名册（在册只 `G6neg`）；⚠ `internal/panel` 那 **2 枚已知红**照实记不修。

## 本票**不**解决
不做提交／推送／接受／回滚（雷区，待 owner 逐枚批）；不做文件树（＝票 190）；不画界面（票 77 地界）。

## Progress log
- 09-28 15:2x 编排者立票：来路＝owner 放权（`A373`）＋`180-c1＋182-c1` 名册"没票认领"七处中的两处（未提交枚数、diff）。未派。
- 09-28 17:5x `189-a1`（只读设计核，派单 §C）出件＝`docs/evidence/s1/189-uncommitted-count-design-a1.md`。**AC 框一枚未勾**（勾要非实现者表）。本轮结论摘要：① 取数形状**分维**——已跟踪改动那一半纯 Go 准且便宜（实测 index 3,332 条解析 1.087 ms＋stat 51.571 ms ⇒ `stat-dirty=15/missing=16`，与权威 ` M`15／` D`16 逐档对上），未跟踪那一半纯 Go 朴素遍历给出 **8,684** 对权威 **46**（高估 188 倍 ⇒ 忽略语义才是真代价，对照物 `tools/d22scan/gitignore.go` 817 行且它仍去问外部 git），diff 那一半纯 Go 需手写对象库（本仓 2,016 松散＋17,653 in-pack／4 packs，最大 pack 12,626,706 字节）而 `go.mod` 冻结⇒ 只有外部 `git diff` 一条准路（实测 267 ms／1,521,885 字节）；② 落点＝票 181 同一枚读面 `internal/panel/git.go`（现量 559 行、7 枚导出符号、包内 `os/exec` 命中 0）；③ 只读边界给了五发**结构性**做法（类型面无写落点／AST 依赖名册尺＋正控／入向 hop 不存在→拓扑不可达／不注册为 Tool／`.git/index` 跑前后 sha 未变那发判据——本轮实测 `.git/index` md5 跑 `git status` 前后同为 `1f078d47…`（只算"这一次没写"、不能当通则；`--no-optional-locks` 这枚开关的存在＝git 承认默认路径会去锁））；④ 快照那一格＝嵌在既有 `composer.git` 下的 `uncommitted`＋`diff` 两个子格（对齐普查 §7 的 S2／S3 两堆）。未裁完 6 条列在同件 §9；与 `182-a1`（§A，出件时未交）若命名冲突的点已具名留给编排者。
- 09-28 18:1x 编排者收 `189-a1`（三枚 commit，零产码）：路线与落点照准——**枚数分两半（已跟踪纯 Go／未跟踪问 git，问不到＝第三个枚举值、绝不报 0）**、diff 只走 `git --no-optional-locks diff`、落点＝`internal/panel/git.go` 同一枚读面不开第二份；载体形状＝嵌在 `composer.git` 下开两子格、**不加顶层 key**（顶层撞 `TestComposerContractTypesMatchFrontend`／`Q-51`）。⚠ **"未提交"这一词我按默认定案＝含已 `git add` 的暂存改动**（三档里暂存与未暂存分开计数但都算未提交；理由＝反着读会教用户"add 完就安全了"），要改口径回一句「未提交不含暂存的」。落地腿**排在 `145-r3` 之后**，开工前必做撞钉预检。账 `A395`。
- 10-08 10:4x 只读普查腿 `189-a1`（10-08 程，同名第二轮）出件＝`.scratch/wisp/probes/189/a1/100-findings.md`（尺件 `00`／`10`～`90` 同目录）。**AC 框一枚未勾**。三句结论：① **雷区现量＝这一堆里一枚都没有**——审查原型 `rb-review.js` 动作词面尺 **0 命中（rc=1，同尺在 `rb-plugins.js` 正控能响）**、Go 侧 `bridge.go:42-45,66-67` 入向名册**仍是 6 枚、无 decide/commit/push/revert**、`internal/panel` 包 `os/exec` 命中 **0**；真雷只在 180-182 census 的登记文字里（`:343` 那"雷区 4 处"），未实现、本程只具名指认不代批。② **最少落点＝3～4 枚文件／约 355～475 行，⛔ 不需要新面板方法名＝不触 C17、不需要走人工批准**：改 `internal/panel/git.go`（新子格视图＋字段，⚠ 并须改掉它 `:5-8` 那句 "NO EXTERNAL PROCESS" 自陈）＋新增同族 `git_readonly_exec.go`（未跟踪那一半的唯一准路）＋两枚判据文件；**免碰** `pump.go`／`panel_pump.go`／`run.go`（现量装配点在 `run.go:703 Git: rt.gitView`）／`pump_test.go` 那两枚**顶层四键**字节钉／`frontend/src/lib/panel.ts`（`composer_test.go:48` 的 `pairs` 名册**不含 `GitView`** ⇒ 嵌套子格不与 TS iface 对账；界面那一跳归 **Q-51**，是"另派界面腿"不是"批契约"）。③ **票面两处前提已过期，写腿别按原句开刀**：AC#6 那句"`internal/panel` 那 **2 枚**已知红"实为 **4 枚**（逐名见 `docs/evidence/s1/197-subagent-carrier-r3.md:268-275`，其中两枚的红因已加长到 `[instructions tasks]`＋顶层 `git` 键在 `panel.ts:112` 从未声明）；AC#6 那句裸 `gate-clauses.sh` 实物在 `.scratch/wisp/probes/154/gate-clauses.sh`（48,783 字节，10-03）。另记：**09-28 同名件的全部行号已漂**（`composer.go:95→125`、`:216→246`、`:232→285`；`pump.go:120→150`；`panel_pump.go:108-109→228-229`；`run.go:421-427→703`；`git.go` 559→**569 行**），且 diff 那一支今天仍**整块没有**（不是"有执行器缺解析"：全仓唯一真起外部 git 的是 `tools/d22scan/gitignore.go:373`，不在生产射程）。⛔ 零 go 命令、零产码写点、未 push。
