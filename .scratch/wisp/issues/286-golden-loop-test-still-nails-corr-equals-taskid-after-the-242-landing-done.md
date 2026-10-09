# 票 286 — `internal/agent` 里有一枚**起手即在的红**：`TestGoldenSingleToolCall` 还钉着"corr 等于 taskID"，而票 242 落地的形状已把它铸成 `taskID#call_e1`

**立票**：2026-10-09 11:2x 编排者（机主令「及时补票」；来路＝写腿 `285-r2` 具名报回"我未裁定、未放宽、未顺手改"＋编排者自己现跑复认）
**性质**：**这不是新 bug，是我 10-08 落地那跳留下的未收口账**——产码按契约改了形状，**同一形为的测试期望没人跟着搬**，于是树里从今天起一直躺着一条红。⛔ 本票不新造功能、⛔ 不动 golden。

**Status:** **done**（2026-10-09 12:0x 编排者按 README 规则 4 收口：**五格全勾**＝`AC#1`/`AC#2`/`AC#3`（凭据＝只读腿 `golden-attr-1`，`A770`）＋`AC#4`/`AC#5`（凭据＝写腿 `286-w1`，`77015dbd`/`9fe2f39c`/`0c2e5cb6`/`a4fd8344`；改后 `internal/agent`＋`internal/tools` **393/0/0 rc=0**、反形种坏必红、产码与 golden 零字节）。⚠ **收口≠零残余**，两笔具名留给下一位：① **"整树无第二枚红"这一发不在本票任何格里**——`go test ./cmd/wisp/ ./internal/panel/ -count=1` 由编排者自己现跑（件 `.scratch/wisp/probes/orch/2026-10-09-cmdwisp-panel-rerun.md`），⛔ 不许读成"本票验过全仓"；② 名册里那 **4 枚过期注释句**（`internal/panel/subagent_roster_197.go:57` 一带，还带着已漂的 `loop.go:647` 引用）**不属本票**，归"注释随行"那一族（票 279／票 280 同链）。⛔ 零 push。）
## 现量（引用前先重跑；行号与被引短语同一次取数，尺一律 `git show HEAD:<path>`）

- 编排者现跑（同一条命令、`-count=1`）：`go test ./internal/agent/ ./internal/tools/ -count=1 -v` ⇒ **rc=1**／`internal/agent` **1 FAIL**（`--- FAIL: TestGoldenSingleToolCall`）／`internal/tools` `ok 13.6s`。
- 红句逐字（我的运行输出）：`loop_golden_test.go:70: call identity = {TaskID:03a55baf-… CorrelationID:03a55baf-…#call_e1 CallID:call_e1 Name:echo …}, want task 03a55baf-…`。
- 断言原文（`HEAD:internal/agent/loop_golden_test.go:70`）：`if call.Req.TaskID != res.TaskID || call.Req.CorrelationID != res.TaskID {` ⇒ 它要求**请求侧的 corr 逐字等于 taskID**。〔⚠ **这枚行号我引错了，真身＝`:69`**；`:70` 是紧跟其后的那行 `t.Errorf`，而红句里的 `:70` 报的是 **`Errorf` 的行号**、不是 `if` 的行号——所以上面那条红句原文不动。来路＝只读腿 `golden-attr-1` 顶回，我 `git show HEAD:… | sed -n '66,72p' | cat -n` 独立复认。落地按 `:69` 搬，⛔ 别照 `:70` 动手。详见文末 §更正。〕
- 产码现形（`HEAD:internal/agent/loop.go:603` `callCorr(taskID, callID string, index int)`，唯一调用点 `:676` `TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i)`）⇒ 带 call id 时返回 `taskID + "#" + callID`，**永不等 taskID**（除非 callID 为空走 `%s#call-%d`，那也带后缀）。
- ★**同文件里还躺着一枚绿的反形**：`HEAD:internal/agent/loop_golden_test.go:341` 逐字 `if row.Tool != "echo" || row.CorrelationID != res.TaskID {` ——它钉的是 **journal 行**那一侧，而 journal 侧至今仍是 taskID（`loop.go:369` `newTaskJournal(l.opt.Journal, taskID, taskID)`，第二实参经 `journal.go:59` 落到 `:81 CorrelationID: t.corrID`）。⇒ **同一枚文件里两枚同名"correlation id"的期望现在互相矛盾：请求侧那条被真值打红、账本侧那条靠残余保持绿。**这笔残余不是我今天才发现：票 282 `AC#4`(a) 已裁"真存在／判留／归票 230 那一族"（台账 `A758`）。
- ⚠ 射程（⛔ 别把这行当全集）：本票只处理**请求侧那一枚红**；`285-r2` 报的是"在锚 `ce18b3d6` 上起手即红"，我没有在旧 commit 上重跑过（⛔ 共享工作树里禁 `checkout`/`stash`/worktree），所以**"哪一笔落地让它变红"属 §要建什么 的 AC#1，不在现量里**。

## 要建什么

- [x] **AC#1 归因到那一跳**（只读腿，⛔ 不许 checkout）：用对象层两态对拉判出"请求侧 corr 何时开始带 `#call-` 后缀"——尺＝`git show <ref>:internal/agent/loop.go`（`callCorr` 在不在／`:676` 那行形状）× `git show <ref>:internal/agent/loop_golden_test.go:70`（期望在不在；〔⚠ 真身 `:69`，见 §更正〕），逐枚给出候选 commit 名册与**边界那一枚**；同时答："那枚期望自它诞生起钉的是什么"（⛔ 不许按票号猜，读代码）。⚠ 这一格交的是**归因**，⛔ 不交修法。
      ✅ 翻勾者＝编排者（非实现者；件＝腿 `golden-attr-1`，commit `d106fc25`，`.scratch/wisp/probes/286/attr-1/00-findings.md` 67 行＋`10-registry.md` 79 行）。**边界那一枚＝`bd124b2a`**（`242-corrland-1`，author `2026-10-08 19:56:27 +0800`，HEAD 祖先）；两态对拉我独立复跑同尺坐实：父 `52501e22` 的 `internal/agent/loop.go:654` 逐字 `TaskID: taskID, CorrelationID: taskID,`（corr 逐字＝task，那枚期望当时为真为绿）⇒ `bd124b2a` 的 `:676` 逐字 `TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i),` 而**同 commit 的期望未跟着搬**＝红，且此后两枚文件无人再动 ⇒ 躺到 HEAD。**"那枚期望自诞生钉的是什么"＝生于 `ae23da67`/ticket10，当时 `loop.go:599 CorrelationID: taskID` ⇒ 它一生钉的都是"242 前 corr＝taskID"这个旧形状**（⛔ 不是新造的、也不是本票要防的形）。★**根因是它自己的 commit message**：`bd124b2a` 的 subject 逐字写着"`loop_approval_test.go:213` **唯一**等值断言具名重判"——那个"唯一"就是它漏掉 `golden:69` 的原因；与我 10-08 的 `A745` 派单同形（我只点了那一枚名字，没给"等值断言全集"这把尺）。⇒ **同形病我仓里早有一条纪律没被执行：全称式句子（"唯一"/"全部"/"整族"）必须并排那把穷举尺**，否则下一程会把它当边界。⛔ 本格不交修法、AC#4 未动。
- [x] **AC#2 谁是权威（本格归编排者裁，⛔ 实现腿与验收腿都不许自裁）**：`loop_golden_test.go:70`〔⚠ 真身 `:69`，见 §更正〕那枚等值期望，与 `docs/PLAN.md:1368` 的 C18 句（"corr 坐在 task id **旁边**、⛔ 不要求等于它"）＋台账 `A745` 我给落地腿的派单原话（"同一任务并发两问 ⇒ 两枚 corr **不同**且各自可路由"）冲突 ⇒ 判"该搬测试期望"还是"该回退产码"。先例已盘上：票 282 `AC#2` 我裁过"那枚等值钉当年钉的是**这个 bug 本身**"（`A758`）。⛔ **无论判哪一支，"把断言删掉／改成永远通过"都不是选项。**
      ✅ 翻勾＝编排者（本格的交付物就是裁决，裁决文本在文末 §AC#2 编排者裁决；三条凭据同批列在那一节）。⚠ **本格翻勾≠本票收口**：AC#4/AC#5 归写腿 `286-w1`，勾要非实现者裁。
- [x] **AC#3 全仓同类旧钉名册（按语义变体扫，⛔ 只搜一种字样＝本仓抓过的假话形）**：把"corr 等于 taskID"这一形的所有写法扫出来——至少含 `CorrelationID != res.TaskID`／`== taskID`／`corr == t`／位置参数 `taskID, taskID`／`CorrelationID: taskID` 诸形，逐枚给"该搬／该留（并写清为什么这一侧允许相等）／待人裁"。⚠ 名册里**必须并排 `:341` 那一枚**（它今天靠 journal 残余而绿，是本票最容易被"顺手一起改"带跑的一枚）。
      ✅ 翻勾者＝编排者（非实现者；件＝`golden-attr-1` `d106fc25` 的 `10-registry.md` 79 行）。腿给真树 **25 枚**＝**该搬 1／该留 19／待人裁 4／冻结 1**，另 `.scratch` 探针桶 162 行／约 29 件**不算动作**。**我自己那把尺**（`git grep -nE 'CorrelationID (!=|==) ' HEAD -- '*.go'`）复跑：**全树只有 `internal/agent/loop_golden_test.go:69` 拿请求侧 corr 和 `res.TaskID` 比等值** ⇒ "净动作面＝1 枚、没有第二枚同形红"这一条**独立成立**；并排读数：同尺还命中同文件 `:341`（journal 侧，腿判"该留"，与票 282 `AC#4`(a)／`A758` 同判）、四枚 242 之后的正确尺（`internal/agent/ticket285_corr_distinct_rulers_test.go:119`／`internal/tools/loop_approval_test.go:219`／`internal/tools/ticket283_corr_identity_rulers_test.go:222`／`internal/tools/ticket285_corr_rows_rulers_test.go:72`）、以及 `internal/tools/bridge.go:269`／`internal/memory/models.go:167` 那两枚"空值回退/空值拒写"（不是等值形，我的尺把它们带出来是**尺的形状所致**，⛔ 不许据此说腿漏计）。**待人裁那 4 枚＝过期注释句**（其中一枚我现量到逐字：`internal/panel/subagent_roster_197.go:57`「loop dispatches with `CorrelationID == TaskID` (`internal/agent/loop.go:647`)」——它引的 `:647` 今天已漂，真身 `:676`）；⛔ **不并入 AC#4**，那是"注释随行"那一族（票 279／票 280 同链，⛔ 不顺手改）。⚠ **本格里两枚降格**：腿对 `cmd/wisp/panel_pump_test.go:323` 与 9 枚宿主 fixture 的"绿"是**对象层静态推、非实跑 rc**（它自己具名报了），⇒ 那 10 枚读数是〔仅静态推〕，"整树无第二枚红"要等一次 `go test ./cmd/wisp/ ./internal/panel/ -count=1` 才算闭合——**这一欠账记在我（编排者）身上**，⛔ 不算腿的欠账（派单当时禁了 Go 面）。
- [x] **AC#4 搬完之后要有仪器证明它没被搬松**（由写腿交、非实现者裁）：① 改前整包红名册（具名 `TestGoldenSingleToolCall`）；② 改后 `go test ./internal/agent/ -count=1` **rc=0**、且**不许新增 `t.Skip`**、PASS 枚数只增不减；③ **反形自证**：把产码 `callCorr` 换回"恒返回 taskID"那一形 ⇒ 新期望**必须红**（证明搬过的钉仍然咬得住，不是一条被抹平的等值句）；四件套齐（种前 hash／`sed -n` 复量／红句 `文件:行`／还原 hash 等值＋`git status --porcelain -- internal cmd` 空）。
      ✅ 翻勾者＝编排者（非实现者；腿＝`286-w1`，四笔 `77015dbd`（起手锚）／`9fe2f39c`（搬期望，名册**只有** `internal/agent/loop_golden_test.go`，+13 −2）／`0c2e5cb6`（读数）／`a4fd8344`（commit 后复检＋终局复跑）。**同尺三数**：`GOFLAGS= go build ./...` rc 0→0；`go test ./internal/agent/ ./internal/tools/ -count=1 -v` rc **1→0**、PASS **392→393**（只增不减）、FAIL **1→0**、SKIP **0→0**（⛔ 零新增 `t.Skip`）。改前那一枚红逐具名＝`TestGoldenSingleToolCall`。**反形那一发**：`callCorr` 两枚 `return` 种成 `return taskID`（调用点 `:676` 未动）⇒ **红**，红句逐字（rc=1，搬后 `if` 在 `:77`、`t.Errorf` 在 `:81`，⚠ 又一次"红句行号≠断言行号"）：`loop_golden_test.go:81: call identity = {TaskID:fbb95865-… CorrelationID:fbb95865-…（无 # 后缀）CallID:call_e1 …}, want task fbb95865-… and correlation id fbb95865-…#call_e1 (task id prefix + this call's own id, never the bare task id)`。**我独立复跑的四把硬尺**：① `git diff --stat 332e2828..HEAD -- internal/agent/loop.go internal/agent/journal.go internal/agent/testdata/`＝**空**（产码与 golden 零字节）；② `git status --porcelain -- internal cmd`＝**空**（突变没留在盘上）；③ `git show --stat` 逐枚＝名册合起来只有那一枚产测文件＋它自己的 probes；④ 还原 hash `cb7f9da7a6d2…63189`（sha256）与种前逐字等值由腿交、我以"产码 diff 为空"侧证。**搬成的形状**（按 `:69` 动、⛔ 没按 `:70`）＝`HasPrefix(corr, res.TaskID)` ＋ `TrimPrefix(corr, taskID+"#") != "call_e1"` ＋ `corr == res.TaskID` 三条同时钉，`TaskID` 那半等值照原样留 ⇒ **至少和原来一样严**（原形只放行单点 `{taskID}`，新形放行"前缀＋这一枚 call 自己的 id"）。⚠ **两处照实带**：(a) 同文件 `:341`（journal 侧）改后**仍绿**＝"停手回报"那条未触发，符合 `AC#2` 裁决的射程；(b) 腿把 `call_e1` 硬写进两处期望（它自己具名），我判**可接受**——同一发 golden 回放在 `:66` 已逐字断言 `CallID == "call_e1"`，两枚是同一条流的两个面，⛔ 不是新造的常量。
- [x] **AC#5 越界与门禁**：⛔ 零触碰 `internal/agent/testdata/golden/**`（尺＝`git show --stat` 名册里不许出现 testdata 字节）、`docs/PLAN.md`／`docs/specs/**`／`docs/SLO.md`／`internal/observe/thresholds.go`／`allowlist.txt` 一字节不动；三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）一字不动；`frontend/**`／`design/**` 零读零写零转述；`sh scripts/d22scan.sh` 纯净树 rc=0；`gofmt -l`／gofumpt 对动过的件空输出。
      ✅ 翻勾者＝编排者（凭据＝腿件 `.scratch/wisp/probes/286/w1/30-gates.md` 96 行＋`40-post-commit-recheck.md` 54 行）。四数：**`gofmt -l` 空 rc=0**；**`gofumpt -l`（v0.12.0）空 rc=0**；**`sh scripts/d22scan.sh` rc=0 clean**（它并报 ban #8 分母 `internal/`＝**521 枚**含 `_test.go`，与"注释豁免、字符串不豁免"的现行射程同形）；**`sh scripts/check-path-length-budget.sh` rc=0 VERDICT GREEN**（tracked 名册从 8485 涨到 **8497**＝新写的 probes 进了分母、`not in roster=0` ⇒ 没让任何人替它挂号）。禁列名册：`git show --stat` 四笔合起来**只有** `internal/agent/loop_golden_test.go` ＋ `probes/286/w1/**`，⛔ 无 testdata／无三枚冻结件／无 `frontend`／无 `design`／无 `docs/**`。⚠ **一处由我承担的形错**：腿文末追加到本票的那一段，**带走它的 commit 是我的并发笔 `0b0ab517`**（它那一发 `git commit` 报"无改动可提交"）——腿**没有把它记成自己的署名笔、也没有硬造空 commit**，并具名给出"共享工作树里文末追加的落点要在 HEAD 上 grep 量、⛔ 不能靠自己的名册"（`grep -n` 命中 `:51/:57/:89/:116`，我核过段落在 HEAD 完整）。⇒ **这条纪律我收下**：今后派写腿做"票面文末追加"时，要么我先 commit 再让它写，要么把落点凭据改成 `git show HEAD:<票>` 而不是它自己的名册。

## 禁区

- ⛔ **不许为变绿放宽断言**：搬之后的期望必须**至少和原来一样严**（要钉的是"corr ＝ `taskID` 前缀 ＋ 该 call 自己的那一枚 id"这一形，⛔ 不是"corr 非空"）。
- ⛔ 不许动 journal 那一侧（`:341`／`loop.go:369`／`journal.go:81`）来"让两枚一致"——那是票 282 `AC#4`(a) 已裁**留**的残余，动它要先在票 230 那一族里挂号，⛔ 不在本票射程。
- ⛔ 不许新造 D43 之外的状态词；不许把 `SessionID`／grant 塞进 `tools.Decision`／`agent.ToolRequest`（票 224／242 撞过同一枚反射钉）。
- ⛔ 票 284 那一格（`gate.go` 注释）**不搭本票的车**：本票不动 `gate.go`，别顺手捎带。
- ⛔ 零 push；commit 必带显式 pathspec。

## AC#2 编排者裁决（**2026-10-09 11:3x**，本格归我裁，写腿与验收腿⛔不许自裁）

- **裁＝搬测试期望，⛔ 不回退产码**。凭据三条（同一次现读）：① `docs/PLAN.md:1368` 的 C18 句把 corr 定义成"坐在 task id **旁边**、⛔ 不要求等于它"；② 台账 `A745` 我给落地腿的派单原话＝"同一任务并发两问 ⇒ 两枚 corr **不同**且各自可路由"；③ 票 282 `AC#2` 我已裁过同形先例——那枚等值钉当年钉的是**这个 bug 本身**（corr 逐字＝taskID），保留原字面＝造一格"只能靠回退 bug 才能满足"的死格（`A758`）。⇒ `loop_golden_test.go:70` 的期望属"钉住旧缺陷"那一族，**该搬**。
- **搬成的形状（至少和原来一样严，⛔ 不许删、不许改成"非空即可"**）：请求侧那枚 corr 要同时钉 ① 以 `res.TaskID` 为**前缀**、② 后缀＝该 call 自己的 id（golden 回放里逐字 `call_e1`）、③ **不等于 `res.TaskID`**；再加一枚**反形自证**（把 `callCorr` 换回"恒返回 taskID" ⇒ 新期望必须红）。
- ⛔ **同文件 `:341` 那一枚不在本票射程**（journal 侧的 `row.CorrelationID == res.TaskID` 今天靠 `loop.go:369 newTaskJournal(l.opt.Journal, taskID, taskID)` 这枚残余而绿；那笔 10-09 已由票 282 `AC#4`(a) 裁"留"、归票 230 那一族，`A758`）。⛔ 任何程不许为了让两枚"看起来一致"去动 `:341`／`journal.go:81`——那会**悄悄改账本侧的形状**。
- **排程**：AC#1（归因）＋AC#3（全仓同类旧钉名册）＝只读腿可办，⛔ 不许带产码改动；AC#4（搬＋反形自证）＝写腿，按在 AC#1/AC#3 之后（先有名册再动手，否则会漏改第二枚同形）。AC#5 门禁随写腿那一笔同批交。

## 更正（2026-10-09 11:5x 编排者，题面不改字；来路＝只读腿 `golden-attr-1` 顶回，我独立复认）

- **一枚行号**：本票正文凡引 **`loop_golden_test.go:70`** 处（现量第 3 条／AC#1／AC#2／上面裁决段第 1 条）指的是同一枚东西，**真身＝`:69`**；`:70` 是紧跟其后的 `t.Errorf` 那一行。尺＝`git show HEAD:internal/agent/loop_golden_test.go | sed -n '66,72p' | cat -n`（我 11:5x 现跑）：`:69` 逐字 `if call.Req.TaskID != res.TaskID || call.Req.CorrelationID != res.TaskID {`，`:70` 逐字 `t.Errorf("call identity = %+v, want task %s", call.Req, res.TaskID)`。⚠ **红句里那个 `:70` 不是错的**——Go 打印的是 `Errorf` 的调用位置，所以"红句说 `:70`／断言在 `:69`"两件事同时成立。⇒ **写腿只按 `:69` 动手。**
- **一枚措辞收紧**：裁决段里我把 `:341` 描述成"`row.CorrelationID == res.TaskID` 那一形"——那说的是它**期望成立的关系**；盘上逐字是 `if row.Tool != "echo" || row.CorrelationID != res.TaskID {`（反向拼写）。⛔ 别因为我那句转述去"找一条 `==` 的写法"。
- **AC#1/AC#2/AC#3 已翻勾**（凭据在各格 ✅ 注，翻勾者＝编排者＝非实现者）；本票**仍未收口**——AC#4/AC#5 归写腿 `286-w1`（2026-10-09 11:5x 派出，独占 Go 编译/测试面），交完由非实现者裁。
- **一枚欠账记我、不记腿**："整树无第二枚红"这件事，腿只给了对象层静态推（它自己具名报了），闭合需要一次 `go test ./cmd/wisp/ ./internal/panel/ -count=1`——那是 Go 独占面的活，⛔ 不该出现在只读腿的射程里，所以这笔归编排者（挂在 `286-w1` 之后的验收腿上，⛔ 不在本票 AC#3 那格里假称已闭）。

---

## Progress log — 写腿 `286-w1`（AC#4 搬期望＋反形自证 ／ AC#5 门禁）2026-10-09 11:4x–11:5x +0800

> ⛔ 本段只追加，票面既有句子（含各 AC 措辞、更正节）一字未改；⛔ 本腿未翻任何框——勾归非实现者。
> 证据件：`.scratch/wisp/probes/286/w1/00-anchor.md`／`10-readings.md`／`20-mutation.md`／`30-gates.md`＋`logs/*`（全部 `.md`，⛔ 无 `.out`）。
> 起手锚＝`332e2828`；本腿 commit＝`77015dbd`（锚）→ `9fe2f39c`（改期望）→ `0c2e5cb6`（读数落件）。⛔ 未 push。

**AC#4① 改前红名册（具名，整发只有这一枚红）**

`GOFLAGS= go build ./...` rc=`0`；`GOFLAGS= go test ./internal/agent/ ./internal/tools/ -count=1 -v` rc=`1`
⇒ **PASS 392 / FAIL 1 / SKIP 0**。红句逐字（`logs/before.md:92`）：

```
    loop_golden_test.go:70: call identity = {TaskID:6b81876c-cf16-4c6a-b57b-b3e4f9e45ed1 CorrelationID:6b81876c-cf16-4c6a-b57b-b3e4f9e45ed1#call_e1 CallID:call_e1 Name:echo Args:{"text":"22 摄氏度，晴"} Timeout:2s}, want task 6b81876c-cf16-4c6a-b57b-b3e4f9e45ed1
```

⇒ 起手即在的红得到复认；**除此之外这两包无第二枚红**（`internal/tools` 两跑均 `ok`）。

**AC#4② 搬的那一枚与改后三数**

按编排者 AC#2 裁决＝**搬测试期望、不回退产码**，动手位置＝**`:69`**（本票更正节已把正文各处的 `:70` 统一收为 `:69`：`:70` 是 `t.Errorf` 行、红句报的正是它）。
`internal/agent/loop_golden_test.go` 一枚 `if` 由"corr 逐字等于 taskID"改为**同钉三件**，`TaskID` 那一半等值期望照原样留：

```
	wantCorr := res.TaskID + "#call_e1"
	if call.Req.TaskID != res.TaskID ||
		!strings.HasPrefix(call.Req.CorrelationID, res.TaskID) ||
		strings.TrimPrefix(call.Req.CorrelationID, res.TaskID+"#") != "call_e1" ||
		call.Req.CorrelationID == res.TaskID {
		t.Errorf("call identity = %+v, want task %s and correlation id %s (task id prefix + this call's own id, never the bare task id)", call.Req, res.TaskID, wantCorr)
	}
```

① 以 `res.TaskID` 为前缀、② 后缀逐字 `call_e1`（两枚合起来即 `corr == taskID + "#call_e1"`，可推）、③ 不等于裸 `res.TaskID`。
⛔ 未删断言、⛔ 未改成"非空即可"、⛔ 零 `t.Skip`。`t.Errorf` 文案一并改成说实话的 want 句（原文 `want task <id>` 在新期望下会误导）。

同一条命令复跑 ⇒ rc=`0`／**PASS 393 / FAIL 0 / SKIP 0**（PASS **392→393 只增不减**，+1 即翻绿的那一枚；SKIP 前后逐字同 0）。
`internal/agent` 包行 `FAIL … 1.905s` → `ok … 2.125s`。**★ `:341` 那枚 journal 侧仍绿（FAIL=0 已含它）⇒ 派单里"若 `:341` 变红立刻停手回报"那一条未触发**，账本侧形状（`loop.go:369`／`journal.go:59`／`:81`，票 282 `AC#4`(a) 已裁留、台账 `A758`）一字节未动。

**AC#4③ 反形自证（四件套，全文见 `20-mutation.md`）**

种前 hash `cb7f9da7a6d277f341d3b40caee59eecec3bea4e4faee31777b1ad47c3763189 *internal/agent/loop.go`
→ `sed -n '600,612p'` 复量 → 把 `callCorr`（`:603`）**两枚 `return` 都种成 `return taskID`**（唯一调用点 `:676` 未动；`grep -c 'fmt\.' loop.go`=9 ⇒ 不会假绿成编译失败）
→ `GOFLAGS= go test ./internal/agent/ -count=1 -run TestGoldenSingleToolCall` **rc=`1`**（输出先落 `logs/mutant-run.md` 再读），红句逐字：

```
--- FAIL: TestGoldenSingleToolCall (0.01s)
    loop_golden_test.go:81: call identity = {TaskID:fbb95865-367d-4b21-ae9c-b3e4c7e2a539 CorrelationID:fbb95865-367d-4b21-ae9c-b3e4c7e2a539 CallID:call_e1 Name:echo Args:{"text":"22 摄氏度，晴"} Timeout:2s}, want task fbb95865-367d-4b21-ae9c-b3e4c7e2a539 and correlation id fbb95865-367d-4b21-ae9c-b3e4c7e2a539#call_e1 (task id prefix + this call's own id, never the bare task id)
```

⇒ 新钉**咬得住**（第二钉 `TrimPrefix … != "call_e1"` 与第三钉 `== res.TaskID` 同时为真；搬后的 `if` 在 `:77`、`t.Errorf` 在 `:81`，红句报 Errorf 行，与搬前同形）。
还原＝`trap … EXIT` 保证；还原 hash **逐字等值**（同一枚 `cb7f9da7…`）、`sed -n` 复量与种前同文、`git status --porcelain -- internal cmd` = **空（0 行）**、复跑 rc=`0`。

**AC#5 门禁四数**

`gofmt -l internal/agent/loop_golden_test.go` = 空 rc=`0`；`gofumpt -l`（v0.12.0）= 空 rc=`0`；
`sh scripts/d22scan.sh` rc=`0`（`clean - no D22 ban violations`；ban #8 的 `internal/` 射程 521 枚**含注释与 `_test.go`** ⇒ 我新增的注释与 Errorf 文案零 emoji）；
`sh scripts/check-path-length-budget.sh` rc=`0`（`VERDICT GREEN`，分母 8485 枚跟踪路径／over-budget 57 全部在名册内／not in roster=0）。
`git show --stat` 逐笔名册：`77015dbd`＝1 枚证据件、`9fe2f39c`＝**只有 `internal/agent/loop_golden_test.go`**、`0c2e5cb6`＝9 枚证据件 ⇒ ⛔ 零 `testdata`／零三枚冻结件／零 `frontend`／零 `design`／零 `PLAN.md`／`specs`／`SLO.md`／`thresholds.go`／`allowlist.txt` 字节。
硬尺：`git diff --stat 332e2828..HEAD -- internal/agent/testdata internal/agent/loop.go internal/agent/journal.go` = **空**（`PROD_DIFF_LINES=0`，这把尺连并行别家的 commit 一并计入）。

**本腿自报的两处不干净（不当已闭）**

1. 第一发越界 grep 把 `golden` 写进模式，于是它**误命中文件名 `internal/agent/loop_golden_test.go`**（本票目标件，不是 `testdata/golden/` 语料）。该支判据无效，以 `internal/agent/testdata/` 目录尺与上面 `PROD_DIFF_LINES=0` 为准。
2. 读数落盘形制：本腿派单正文步骤 2 写 `logs/before.txt`、同一派单文末「输出落盘规矩」写"只有 `*.md` 会被门禁放过"。两行冲突 ⇒ 按后者落 `logs/before.md`／`after.md`，⛔ 无新建 `.txt`/`.sh`/`.ps1`/`.out`。

**本腿没答的格子（⛔ 不假称已闭）**

- **AC#3 全仓同类旧钉名册**＝归只读腿，本腿未交；本腿只在票面既有 `:341` 一枚上做"不碰"处理，未做穷举。⇒ 若名册里另有同形旧钉，本腿不知道、也没改。
- **"整树无第二枚红"没有闭合**：本腿只跑 `./internal/agent/ ./internal/tools/`（派单指定那一发）。票面更正节已把 `go test ./cmd/wisp/ ./internal/panel/ -count=1` 那一发记给编排者、挂在 `286-w1` 之后的**验收腿** ⇒ 本腿未跑、不据以声称全树无红。（`cmd/wisp` boot Ctrl+C flake 与既有 winlive 红历史在本腿读数之外，未观察也未否证。）
- **AC#4② 的 PASS 计数尺是"整发 `-v` 输出逐条 `--- PASS`（含缩进子测试）"，未按包拆分**；两跑同尺可比，但换一把尺（如只数顶层测试）数字会不同，验收腿复核请同尺。
- golden 回放那一枚 call id 我**硬写了 `call_e1` 两处**（前缀钉＋后缀钉＋wantCorr 串）：它由同文件 `:66` 的既有断言钉住，故不会随语料漂移而假绿；但"语料若改名⇒本测试红"这一条属于 golden 面，⛔ 不在本票射程，本腿未裁。

## 结案后更正一条归属（2026-10-09 12:1x 编排者，⛔ 原句一字不改；本节只追加）

- `Status` 行末尾与 `AC#3` 那格末尾都写着"那 **4 枚过期注释句**归'注释随行'那一族（**票 279／票 280 同链**）"——**那句话挂到了一条不存在的链上**：票 279 只管 `cmd/wisp/run.go` 的 P39 那一块（它自己的 `AC#1` 至今没触发），票 280 已 `-done` 且裁的是 `gofmt` vs `gofumpt` 的**归属**，两枚都不覆盖这 4 句 ⇒ 按票 288 同形处置，**已另立票 289**（`.scratch/wisp/issues/289-…-an.md`，立票笔 `9fda613b`；四枚逐字读数＋"行中性"硬判据在它自己的现量表里）。本笔不改上面任何原句，只把"归 279/280"这一处指向更正为 289。
