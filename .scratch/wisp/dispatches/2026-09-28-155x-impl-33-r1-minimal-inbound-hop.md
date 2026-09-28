# 派单：`33-r1`（写码腿·片 A）——把"网页点一下 → Go 真收到"那一跳的**六环**建起来（不引 WebView2、不加 `panel.*` 方法）

- 派单时刻：`2026-09-28 15:5x`｜编排者锚点 HEAD＝**取你起手时的 `git log -1 --format='%h'`**（我派单时盘上是 `a10832ef` 之后，⚠ **本轮已第 4 次锚点漂 ⇒ 起手先取号、把它写进表 §1，别照抄我这里**）
- 工单：`.scratch/wisp/issues/33-panel-host-c27.md`（面板宿主层）；下游依赖＝票 186（切工作树／分支）与票 187（改模型与思考档）**都排在这一跳后面**
- 性质：**写码腿**｜⚠ **AC 框一枚不许勾**（翻勾由编排者照你的读数＋我自己复跑做）

## 0. 为什么现在派（一句话）

`33-a1`（只读设计核，表 `docs/evidence/s1/33-inbound-hop-design-a1.md`，编排者已收＝台账 `A374`）现量出：**最小闭环 H1–H10 里，H4–H9 这六环一环都不引用 WebView2 符号**——raw 报文 →`ParseComposerRequest`（`internal/panel/bridge.go:84`）→ 派发 → mode／workspace／attachment／message 四条腿 → 写腿 `internal/panel/composer_handlers.go:141 h.Modes.Set`；只有 **H2 建窗／H3 消息接收／H10 回灌手段**要真宿主。⇒ **这一程只做那六环**，把"写腿调用点已写好在等门铃"接上，**不碰真窗口**。

## 1. 写面（只这些，逐枚具名）

✅ 新建 `internal/panel/composer_dispatch.go`｜✅ 新建 `internal/panel/composer_dispatch_test.go`（或同名 `*_183` 风格带 `33` 的判据件）｜✅ 新建 `docs/evidence/s1/33-minimal-inbound-hop-r1.md`｜✅ 新建 `.scratch/wisp/probes/33/r1/**`（读数只往这里写）｜✅ 追加票 33 一段 Progress log（**只追加、不改原句、不勾框**）。
⛔ **`internal/panel/bridge.go` 一字不动**（那四枚 `panel.*` 白名单＝C17 面，`Q-67` 未答 ⇒ 一枚不许加、顺序不许改；尺：`grep -c 'panel\.' internal/panel/bridge.go` 起手终态都要 **＝4**）。
⛔ **`cmd/wisp/run.go`／`cmd/wisp/panel_pump.go` 本程不动**（H10 那一环留到真宿主那天；⚠ 另一枚按住的写腿票 174 AC#2b 要动 `run.go`，**不许你俩同批**）。
⛔ `go.mod`／`go.sum`／`deps.toml`（**本片 A 判为 0 新依赖**；判"必须引 webview2 才能答完"＝**停手上报**，那是 H2/H3 的活、不属本票）｜⛔ `docs/PLAN.md`、`docs/specs/**`、`allowlist.txt`、`thresholds.go`／golden／审批超时常量、`internal/panel/tokens_fourway_test.go`／`l2_grant*`／`frontend_hygiene*`／`composer_test.go`／`pump_test.go`／`approval_test.go` **一字不动**。
⛔ `frontend/**`／`design/**` 不读不写不引不转述（别家地界）。
⛔ 不属于你的脏件（`.gitignore`、`probes/152/**`、`probes/161/r6/logs/flip-*`、`docs/evidence/s1/152-*.md`、`design/**` 那批未提交删除）**不动不提交不评论不还原**。
⛔ 零删除命令；不许跑 `probes/161/r6/flip-declaration.sh`；**不许改 `scripts/d22scan.sh` 或 `probes/154/gate-clauses.sh` 本身**；只 commit **不 push**；显式 pathspec；禁 `add -A`／`.`／`-a`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`worktree`／`switch`／`merge`。
⚠ **写面闸门（新规矩，`A374` 那条）**：**起手**跑 `git status --porcelain -- internal/ cmd/` 存进表 §1；**终态**要求＝**等于起手名册（逐枚具名差集为空）**，⚠ **不是"必须为空"**——同树里有别家程在飞是常态，多出/少掉的每一枚都要具名说"这不是我的／我预期内"，不许去动它。
⚠ 重跑任何既有台件前先 `grep -n "WriteFile\|OpenFile" <那枚台件>` 看写出路径（`A367` 那枚坑：重跑会原地洗掉被别人逐行引用的读数）。

## 2. 必答六格（每格先测→再写→再提交；引代码要 `file:line` 逐字）

**AC#A 路由本身**：建一个**入向入口**（宿主侧，不是模型工具），把一条 raw 报文解成"方法名＋参数"，按 `bridge.go:42-45` 那四枚**已存在**的 `panel.*` 方法名派到四条现有腿（`ParseComposerRequest`→mode／workspace／attachment／message）。⚠ 判据必须**问被调次数**：把路由摘掉 ⇒ 那一发必红（`33-a1` 给的形状）。
**AC#B 默认放行是缺陷**：入口收到**名册外**的方法名 ⇒ 必须走**拒答＋审计**，不许静默丢弃；⚠ 判据＝"把默认分支改成放行 ⇒ 必红"（`33-a1` 给的第二发）。这一发要防的就是 `Q-49` 那一族（面板侧来源的批准出口）。
**AC#C 来源不许改**：报文里带的 `source`／taskID 一类归因字段**不许由入口重写**；⚠ 判据＝"改 `source` ⇒ 走拒答而不是静默接受"（`33-a1` 第三发），并**具名证明它和 `bridge.go` 那四枚白名单是同一道闸、不是我又建了一道**。
**AC#D 生产调用者这一问必须诚实回答**：本片 A 建完之后，**生产里真有谁调用这个入口吗**？`33-a1` 现量：今天 `ParseComposerRequest`／`HandleModeRequest` 是**非 test 零生产调用方**，而零调用方有**两种因**（线没接／层不存在，见 `A371`）。⇒ 你要交的是"**一枚具名诊断入口**当生产听众"（`33-a1` §⑦ 的形状）＋**逐字写明它接在哪个命令上、用户怎么敲**；⚠ **禁止**用测试里的假宿主充当"已接线"（`AGENTS §1.3`：不许用 mock 代替真的来假报完成）。**判"必须有真 WebView2 才有生产调用者"＝如实写进表里并停手上报**，那一条就是 `33-a1` 具名的"真窗口六格一格没结"。
**AC#E 与票 92 AC#7 那枚词面尺的关系（⚠ 提前避开，别到 186 才撞）**：`internal/panel/composer_test.go:268` 的 `gitSwitchCapabilityRe` 把 `internal/panel/**` 产码与注释里的**小写单词 `worktree`／`git checkout`／`switchBranch`／`git.branch`** 当"切换能力"扫（`181-r1` 已踩过一次，台账 `A376`／`Q-69`）。⇒ **你不许为了绕开它去改那枚尺、也不许把功能改名成不像它**；正确处置＝**本程不碰"切换"这一维**（那是票 186 的活，且要先裁 `Q-69`）。若你发现非碰不可 ⇒ **停手上报**，写清楚你碰的是哪一行。
**AC#F 契约轴**：`docs/PLAN.md`（含 D29／C17／C27 三处）、`docs/specs/**`（含 `SPEC-08:163-174` 那张 12 行方法表——⚠ **它和码里那 4 枚互不相认这件事是 `Q-67`，未答之前你一枚不许加、一枚不许改**）、`allowlist.txt`、`thresholds.go`／golden、审批超时常量、`go.mod`／`go.sum` 一字节不许动；`internal/panel/bridge.go` 白名单枚数**起手终态都要现量为 4**（尺见上）。

## 3. 门禁（终态读数取在你**自己最后一枚 commit 之后**；⚠ 按 `A363`，你的门禁读数不充当 AC 结案凭据）

逐包 `go test -count=1 ./internal/panel/ ./internal/config/ ./internal/risk/`（⚠ `internal/risk` 必须**单包跑**，四包并发会假红＝`A359`）＋`./internal/tools/`；`sh scripts/d22scan.sh`（⚠ **基线 `ban #8 internal/` examined=436**，不是我更早写的 433；三枚增量逐名归属见 `A376`——**你若再新增 `.go`，读数会涨，具名登记**）；`bash .scratch/wisp/probes/154/gate-clauses.sh` **比红腿名册不比退码**（在册只 `G6neg`，且它 `G6neg` 实测=3 枚而基线=1 枚是**在册状态**）；`gofumpt` 在 `"$(go env GOPATH)/bin/gofumpt" -l` 名下件必须空；CLI 那一面本程**不需要**，要跑必须带在册 PATH 前缀 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（不带＝`0xc0000135`＝票 98）。
⚠ **`./internal/panel/` 今天有 3 枚在册红**（`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`＋`TestC21DesignTokensFourWayAgree`＝别家删了 `design/assets/tokens.css`；`TestComposerContractTypesMatchFrontend`＝Go 送了 `git` 键而前端 `ComposerState` 未声明，归前端会话结，见 `A376`）⇒ **照实记名册、不许当绿、不许去修、不许放宽断言**。

## 4. 表骨架（第 ≤5 枚调用内落盘并 commit 第一枚）

① 起手锚（`git log -1` 现取）＋起手**脏件名册**（逐枚具名）② 落点与写面清单 ③ AC#A–#C 三发红／绿两向逐字读数（含变异载具路径）④ **AC#D 生产调用者那一问的诚实答案**（谁调、用户怎么敲、哪六格还没结）⑤ AC#E 那枚词面尺的处置（我没碰／碰了具名）⑥ 契约轴与 `bridge.go` 枚数两向现量 ⑦ 门禁全套（终态）⑧ 变异自证表（哪一发红哪一枚判据）⑨ 本程**没测**的（逐名）⑩ 对编排者的不服（若有）⑪ 被拒调用逐条＋零删除自证＋工具调用终值 ⑫ next＝真宿主那三环（H2/H3/H10）还缺哪一枚批准

## 5. 硬顶与计数（这次要给可执行的尺）

工具调用**硬顶 35 枚**；**每答完一格自报一次累计**（在表里写"本节末＝第 N 枚"），**第 25 枚起不许开新探索**，只做落盘与门禁；**到 35 枚必须交**（交不回就写"死在哪一格、盘上有哪几枚 commit"）。

## 6. 结束消息回我七项

① 入口落在哪枚文件、派发到四条腿的调用点逐字（`file:line`）；② AC#A／B／C 三发**变异红**的逐字读数（哪枚判据、`--- FAIL` 名字）；③ **AC#D 生产调用者的真答案**（谁调这个入口？用户敲什么命令能看到？没有就明说没有＋停手上报的那一条）；④ `bridge.go` 枚数起手／终态两向＝4 的现量；⑤ 我按住的没碰清单（`run.go`／`panel_pump.go`／词面尺）自证；⑥ 门禁全套（`internal/panel` 那 3 枚在册红照实记）＋起手/终态脏件名册差集；⑦ 被拒调用＋零删除自证＋工具调用终值＋你没测的逐名。
