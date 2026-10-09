# 票 282 — `corr` 落地时**越格**的那一块，交非实现者复核

**立票**：2026-10-08 19:5x 编排者（机主令「及时补票」；来路＝腿 `242-corrland-1` **主动自报越格**＋台账 `A748`）
**性质**：非实现者复核。⚠ 按本仓定式：腿越格自报后，**"它没顺手放宽任何东西"必须由下一枚非实现者答**，⛔ 编排者自己盖章不算。

## 现量（引用前先重跑）

- 裁定原文＝台账 `A745`（四条边界：不动契约文本／不许放宽 `loop_approval_test.go` 那枚等值断言／`loop.go:363` 注释同批重判／三枚测试构造点按上下文重判）。**四条它都做到了**（本票不复核这四条本身）。
- **越格的那一块**（裁定之外的扩面）：`CorrelationID` 不能再当任务 id ⇒ 它给桥加了 taskID 携带＋新访问器 `internal/tools/cancel.go:67 func TaskID(ctx context.Context) string`，并把 `internal/tools/subagent_197.go:262`（⚠ **本行原写 `internal/agent/subagent_197.go`＝路径不存在**，`git cat-file -e HEAD:internal/agent/subagent_197.go` 报 does not exist；真身由只读腿 `282-v1` 顶回、编排者 `git ls-files` 复量＝`internal/tools/subagent_197.go`。⛔ 这一处错记在我自己的票面上）／`internal/tools/task.go:708` 的**任务身份**改为读该访问器（它的理由＝不读则**深度判定静默失效、父停不了子**）。
- 相关提交：`133b1bfa`／`bd124b2a`／`0abe217c`。

## 要建什么（非实现者）

- [ ] **AC#1 行为尺**：把"任务身份改读访问器"那一跳**种坏**（例：让 `TaskID(ctx)` 返回空或回退到 corr），**指名用例必须红**——若全绿 ⇒ 具名写"该支今天没有仪器"。四件套（种前 hash／改的那行 `sed` 复量／红句 file:line／还原 hash 逐字回＋porcelain 空）。
- [x] **AC#2 有没有顺手放宽**：逐枚核它动过的断言（尤其 `internal/tools/loop_approval_test.go:215` 一带的新句）**是不是等价或更强**（三形皆红：空／旧形等值／无前缀）；⛔ 有任一处变松 ⇒ 退回。
      ✅ 编排者 2026-10-09 11:4x 翻勾（凭据＝腿件 `.scratch/wisp/probes/282/v1/10-ac2-ac4-ac5.md`，commit `4d9ac0d8`；⚠ 腿引的行号 `:215` 今天现读为 **`:219-222`**，编排者 `git show HEAD:<file>` 逐字复量）。**判语＝"为变绿而放宽" 0 处 ⇒ 本格不退回**，逐形：空形 `:219`＝**等价**（`r.CorrelationID == ""` 原样在）；旧形等值 `:219`＝**更强**（新句**新增**一条 `r.CorrelationID == res.TaskID ⇒ 红`，旧句那条是唯一通过值）；无前缀形 `:220`＝**等价**（`!strings.HasPrefix(...)`）；无 `t.Skip`／无 `Logf` 降级／无 build tag（新用例 145 行逐枚数过，九处全 `Fatal`），`strings` 非新加（`bd124b2a~1:13` 已在）。
      ⚠⚠ **但腿顶回我那句"比旧句更严"，顶得对一半，我按它缩**：把 `bd124b2a` 的 −/＋ 两半并排读，**通过集确实从单点 `{taskID}` 变成无穷前缀族 `{taskID+任意后缀}`**——不过这不是"放宽保护"：那枚等值钉当年钉的是**这个 bug 本身**（corr 就是 taskID），裁定（`A745`）要求的就是把它拆开，⛔ 保留原字面＝造一格"只能靠回退 bug 才能满足"的死格（同形先例＝票 242 `AC#1` 今天按缩后字面勾）。**真松的是"后缀形状没人钉"**这一条，而它**不属本票**：`git grep -nE 'callCorr' HEAD -- '*.go'` 整族＝产码 3 枚（定义 `loop.go:603`／调用 `:676`／注释 `:367`）＋**测试面只有 `ticket283_corr_identity_rulers_test.go:13` 那枚注释文字、零断言** ⇒ **已转立票 285 `AC#4`**（那一格连同"两枚调用 corr 必须互不相同"这条我派单要过、今天不存在的用例一起补）。
- [ ] **AC#3 深度判定与父子停机**：读 `subagent_197.go` 深度判定那一跳，给"改读访问器前后"的各一发读数（同一形状），判"静默失效"是真被修掉还是只是换了写法。
      ⛔ **本格按住，不算腿的失败**（腿件 `.scratch/wisp/probes/282/v1/20-held-cells.md` 具名）：缺第 3／第 4 发前后对读数（`internal/tools/subagent_197.go:262` 与 `internal/tools/task.go:708` 各一发），且**方法学前置**＝夹具必须 `TaskID != CorrelationID`，否则两发必同绿不作读数——`242-corrland-1` 交付那一刻全仓**没有**这样的夹具（票 283 那把尺件 `039ec93c` 晚于 `bd124b2a`，腿现量 227 行新建零删除）。⇒ 按住理由＝串行铁律（同一时刻只许一枚腿在 Go 导入图里种突变；当时撞 `242-v2`，此刻撞 `285-r1`）。
- [x] **AC#4 残余登记**：它自报"journal 非生产组合行仍 taskID"＋"bridge.go 既有 CRLF 未洗"——逐条核实，判"留 / 改归哪张票"。
      ✅ 编排者 2026-10-09 11:4x 翻勾（凭据＝腿件 §AC#4＋编排者对 (b) 独立复跑同一把尺）。**(a) 真存在**：`HEAD:internal/agent/loop.go:369` 逐字 `newTaskJournal(l.opt.Journal, taskID, taskID)`，第二实参经 `journal.go:59` 落到 `:81 CorrelationID: t.corrID`；"非生产组合"定性成立（`run.go:1003 Journal: nil` 才是 loop 侧，`:750 Journal: mem` 属 `tools.New`＝桥侧）。⇒ 判 **留**，归口＝票 283/284/285 均不认领、实质是"结案票里没做完的格"那一族（先例＝票 230）。**编排者裁：⛔ 不给它加 `DEFERRED(D-xx)` 代码标记**——加了就必须同批在 `SPEC-12 §5` 登记表落一行（人工批准面、且要求 1:1 双向），为一笔探针级残余去开那道面是 disproportionate；这条属账目往哪张表上记，⛔ 不上机主的清单。**另具名一笔不合规矩之处由我承担**：这条残余今天只活在探针 md，没进任何票名下 ⇒ 本格起它**有了归属**（本票＋票 230 那一族）。**(b) 从来就不是那样**＝**它自报过头**：`tr -cd '\r' | wc -c` 两列并排，工作树 `bridge.go` CR＝**1284**（＝`wc -l` 1284，全文件检出态）vs `git show HEAD:` 同一枚 CR＝**0**（`bd124b2a` 0／`bd124b2a~1` 0）；`core.autocrlf=true`＋`.gitattributes:4 *.go text eol=lf`＋该枚 `git status` 空 ⇒ **CRLF 是本机检出形态、不是文件的债**（与票 280 同形，见 `A756`）⇒ 这笔**不是悬账**，⛔ 不许后续程去"洗"它。
- [x] **AC#5 越界检查**：`git log --name-only` 该三笔，核有没有碰 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／`frontend/**`／`design/**`／三枚冻结件。
      ✅ 编排者 2026-10-09 11:4x 翻勾（凭据＝腿件 §AC#5，**名册它自取了、没照抄我给的三枚 sha**——`git log --grep='242-corrland'` 与 `133b1bfa~1..0abe217c` 两把对齐，区间内另 3 枚属别的程）。判语＝**越界 0 枚**：三枚各扫 9 张禁列名册＋一整套 golden 名册（`internal/agent/testdata/golden` 11／`internal/llm/testdata/golden` 40／`tools/mockllm/testdata/golden` 2／`internal/llm/golden` 3／两枚 `*_golden_test.go`／`goldenfmt.go`）＝**27 次 numstat／name-only 全零命中**。⚠ 照实带＝**枚数口径**：三笔里只有 `bd124b2a` 落产码（7 件），`133b1bfa`／`0abe217c` 各只落 1 枚 `.scratch` 探针 md ⇒ 对那两笔这一格是**平凡成立**，别把它读成"三次产码改动都验过没越界"。

## 禁区

⛔ 不改产码（除 AC#1 的突变行，改完必须还原）；⛔ 不许为变绿放宽任何断言；⛔ 零翻框；⛔ 零 push。
