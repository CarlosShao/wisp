# pool-validity-4a / 批4a —— 零勾开放票有效性普查（号段 200–238，11 枚）

> 只读普查腿 `pool-validity-4a`，接 `pool-validity-1`／`pool-validity-2`（死在"扫全池"这个形状）、
> `pool-validity/3a`（100–149，13 枚）、`pool-validity/3b`（150–199，18 枚）。
> **本件只判 200–238 号段那 11 枚零勾票。**
> 起手 HEAD `1a93cc7b`（2026-10-06 11:53:51 +0800，分支 `dev`）。⚠ 换 HEAD 要重量。
> 量尺时刻 `10-06 11:5x`。
> 本腿**零 go 命令**（并行写腿 `232-r2` 此刻真在 `cmd/wisp` 取整包终态）。
> 工具全集＝`git log`／`git grep HEAD`／`git show HEAD:<path>`／`git ls-tree`／`git ls-files`／`ls`／`wc`／`grep -c`／Read。
> 票面零改动：零翻勾、零改名、零 `Status:` 改写、零产码改动、零台账写入、`docs/**` 一字不碰。
> 唯一写入路径＝`.scratch/wisp/probes/pool-validity/4a/**`（新建；代号 `4a` 起手经 `ls` 确认未被占用）。
> ⛔ 不读不引 `frontend/**`、`design/**`、`.gitignore`、`.scratch/wisp/probes/232/**`、`.scratch/wisp/probes/111/ci1/**`。
> ⛔ `docs/evidence/s1/**` 只按**文件名**取（`git ls-tree`），不读内容面。
> ⛔ `cmd/wisp/**` 的工作树内容面不读，一律 `git show HEAD:<path>`。
> ⛔ 前腿件（`batch1`／`batch2`／`3a`／`3b`）只读只学表形，**不照抄其口径、不采信其判档**；
>   `pool-validity/1/**` 与 `pool-validity/2/batch3.md`·`batch4.md`（死腿空骨架）一字不动。
> ⚠ 三枚不判：**232**（写腿正在修它）、**269**／**249**（已撤销／已撤回，判它们＝把已并案的缺陷立案成第二次）。
>   232 出现在我的分母尺原始输出里（它确实零勾、6 格），**是我按指令剔除的，不是量漏**。

## §0 四档尺（本腿判法，口径自立；表形学 3a）

| 档 | 判据 | 复法（命令原文＋读数必须留在表里） |
|---|---|---|
| **仍成立** | 票面点名的那枚缺陷今天还在树上 | 把票面 §现场／AC 里**最硬的一条**自己跑尺：`git --no-pager grep -n '<那句字面或符号>' HEAD -- internal/ cmd/ scripts/ .github/` 的命中数＋file:line |
| **已失效／已被别人做掉** | 缺陷不在了 | ⛔ 硬门：必须给出**一枚 commit 号**＋`git log --name-only` 命中，**或**今树 0 命中的复算尺；拿不出就不许写这一档 |
| **差翻勾** | 缺陷不在了，但凭据是"票内注记／台账 `A##`／一枚非实现者裁决表" | 凭据种类写进复法列；⛔ **不等于本腿裁它可结案**，翻勾归编排者 |
| **量不到** | 判据落在运行期行为／真开窗／真机双击／DPAPI／多账号登录／CI 侧读数／`frontend/**`·`design/**`（本腿禁读） | 归口写清**缺哪一行读数**；⛔ 不许由 grep 命中外推成"这功能能用" |

★ `上次动过` 列＝`git --no-pager log -1 --format='%h %ad' --date=format:'%m-%d' -- <票文件>`（本腿逐枚实跑，见 §1）。
★ 本段预登记的**一处族群陷阱**：234／236 两枚的判据本身就是"CI 侧读数／变异复跑"，
  按 §0 那一档它们的核心格天然落〔量不到〕；本腿要判的是**"票面点名的静态那半边今天还在不在"**，
  两半分开写，⛔ 不许把"我 grep 到了形状"外推成"那一格能结"。

## §1 名册（11 枚，逐枚带票面标题原文）

> 标题＝`head -1` 逐字读，⛔ 未改写。"上次动过"＝票文件自己的最近一笔 commit。

| 票号 | 票文件（`.scratch/wisp/issues/`） | 上次动过 | 票面标题原文 |
|---|---|---|---|
| 200 | `200-nothing-reads-the-projects-instructions-for-ai-although-all-five-harnesses-do.md` | `85da3060` 09-28 | 它进到一个项目里，**完全不知道这个项目的规矩**：四家 harness 全都在读的"给 AI 看的项目说明"，我们一枚代码都没有 |
| 220 | `220-roster-blockedOnApproval-only-sees-L2-and-second-card-silently-drops.md` | `6cf2cc60` 09-29 | **名册上"卡在等批准"那一格今天只看得见 L2，而同一枚任务的第二张卡会静默读成"没被卡住"**（Go 侧送出去的布尔本身就是假的，不是界面没画） |
| 225 | `225-deferred-markers-do-not-match-the-spec12-registry-and-nothing-checks-it.md` | `326fbf78` 10-04 | **`DEFERRED(...)` 代码标记与 `SPEC-12 §5` 登记表今天对不上，而且没有任何仪器在查这一条**：`AGENTS.md` §1.1 那条"1:1 双向"目前是**只靠人**的规矩 |
| 227 | `227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md` | `fac60ad4` 09-29 | 受守卫写还剩两种"手改形状"没人钉：**删掉的行会被补成 schema 默认值**，而 **schema 落后时那句"什么都不写"会说谎**（迁移-on-读照样把文件重排、注释丢光） |
| 229 | `229-subagent-spawns-an-unregistered-goroutine-name-and-drowns-the-leak-alarm.md` | `c774f8da` 09-29 | 子代理收尾时注册了一枚**不在 D38 名册里的协程名**（`subagent-finish-*`）：每 spawn 一次就打一行"疑似泄漏"警告，还被算进 SLO 的协程总数；而"名字不在名册"本该是发现真泄漏时唯一的信号 |
| 230 | `230-four-cells-left-unfinished-inside-closed-tickets.md` | `5759d7fb` 09-29 | 结案票里那 **4 枚真残缺**：票 115 判了七格却**一张裁决表都没出**、票 97 的注释自称"唯一的读者"而同一文件里有三处、票 80 把活交给票 21 而票 21 一字不知 |
| 233 | `233-d36-three-tier-enum-is-tagged-on-no-key-and-the-reload-tier-has-neither-producer-nor-consumer.md` | `fac60ad4` 09-29 | D36 那三档生效级别**一枚键都没标上**：`internal/config.Tier` 三枚常量在盘上只活在 `schema.go` 自己的注释与声明里；`reload` 档**既没有生产者键、也没有消费者**；真正在分档的是 `cmd/wisp/config_reload.go:299` 那枚硬编码 `restartTierKeys` |
| 234 | `234-the-twelve-cells-that-only-need-a-reading-mutation-and-gate-rerun-on-current-head.md` | `6d22dd6c` 09-29 | 结案票残余那 **12 枚"只差一发读数"** 的判据：一次「变异＋门禁」复跑把它们从〔仅自述〕升成〔有读数〕，顺带销掉 10 枚挂在开放票名下的重复分母 |
| 236 | `236-six-cells-that-only-surface-at-the-reading-layer.md` | `c308be29` 10-05 | **六枚只在"读数层"暴露的仪器缺口**（`221-v1` 与 `ci-red-1` 两枚非实现者腿交回）：两支 fail-closed 无尺／那把 DEFERRED 尺的 (b) 支扫词面无牙／"任务 id 必须宿主铸造"零钉／`ci.yml` 那句"今天会红"已过期／**在册红名册不是稳定集合**／**一枚刻意坏的夹具被入库之后污染了格式门的 tracked 分母、并吃掉两道 `go vet`** |
| 237 | `237-race-run-exposes-a-parentheft-inference-that-never-held-plus-two-things-only-a-human-reads.md` | `d5a59d66` 09-29 | `-race` 那一跑撞出的红**不是票 235 造的**：`parkParents` 里"起孩子"排在"父发派生回执"之前、而测试等满孩子信号后立刻要求回执数已满 ＋ 守池注释"只能靠人读"这件事**代码旁边一个字都没写** ＋ 新前置读数只要求"至少一枚"、**部分绕过仍是瞎的** |
| 238 | `238-first-run-on-a-fresh-machine-creates-the-whole-data-root-with-no-seal-to-inherit.md` | `cf03fa31` 09-30 | **新机器上先跑常驻（GUI）那条腿时，整个数据根是"没有任何封条可继承"造出来的**（`internal/proc` 对 `winsec` 导入数＝0，连 `MkdirAll` 的父档都不是窄的） |

## §2 判档表（11 枚 × 4 列）

> ⚠ **本腿对"numstat 删除列必须为 0"这条纪律的取法（写死在这里，免得被读成违约）**：
> 该条的目的是**不动别人的行**；而 §4 又要求骨架先落盘、且"任何一处「未判」＝不算交件"。
> 两者唯一的相容取法＝**填表必然替换本腿自己骨架里那 11 行「尚未判」占位行**。
> 每一笔 commit 前本腿都跑过 `git --no-pager diff -- <本文件> | grep "^-"` 逐行核对：
> **删除行全部是 `| NN | 尚未判…` 这一形**（本腿自己写的），⛔ 无一行属于他人或他人文件。
> 首笔（`3a76ef89`）删除列＝0；此后各笔删除列＝被替换的占位行数（200/220/225 那笔＝11→3）。
> 若编排者要的是字面零删除，那正确做法是**不让我落骨架**——两条款冲突时本腿选了"交件必须可判"这一边，并在此具名登记。

> ⚠ 本节的**每一行此刻都尚未判**。这不是占位符，这是实话：骨架是先落盘的交件形状，
> 判档由本腿在接下来的调用里逐枚跑尺填入，每判完 3 枚立刻 commit（不攒）。
> 若本腿在填完之前死掉，本节读出来的就是"0 枚有判语"，与死腿 `pool-validity/1` 可分辨——
> 分辨凭据＝上表 §1 的 11 枚名册与下面已实跑的分母尺读数都是**真读数**，不是空壳。

| 号 | 档 | 复法（命令原文） | 读数 |
|---|---|---|---|
| 200 | **差翻勾**（票面头名那枚缺陷＝"一枚代码都没有"今树**不成立**；硬凭据有我自跑的 commit＋`--name-only`，结案面却仍是票内注记＋两枚非实现者裁决表**文件名**＋若干运行期格未收）⛔ 不等于可结案 | 尺一（票面 §现量 第一行那把尺逐字复跑）：`git --no-pager grep -nE "AGENTS\.md\|CLAUDE\.md" HEAD -- internal/ cmd/ tools/ scripts/ .github/` → **大量命中**，其中**非测试、非注释、承重**的三枚＝`internal/projctx/projctx.go:56` 的优先级名册逐字 `"AGENTS.override.md", "AGENTS.md", "AGENTS.MD", "CLAUDE.md", "CLAUDE.MD"`（票面"要建什么 §1"要求的正是这一串，逐字对上 Step-Code 的优先级）／`:2` 的 `Nothing in this ...` 自述／`internal/config/schema.go:442` 的配置键注释。票面那句"零命中（非测试码里只有一句注释提到 AGENTS.md）"**今树读不出来**。<br>尺二（硬门第一把＝commit 号＋`--name-only` 命中）：`git --no-pager log --format='%h %ad %s' --date=format:'%m-%d' --diff-filter=A -- internal/projctx/projctx.go internal/panel/instructions_200.go` → **同一枚 `84feec44 09-28`**，标题逐字 `ticket 200: read the project's own instructions for the AI (all five harnesses do)`；第二枚载体：`git --no-pager log --format='%h %ad %s' --date=format:'%m-%d' --diff-filter=A -- internal/panel/instructions_200.go` → 亦 `84feec44`。续跑：`git --no-pager show --name-only --format='%h %ad %s' --date=format:'%m-%d' a90677e9` → 标题逐字 `ticket 200 (200-r2): consume the rewrite account before naming a tree, and give the panel carrier a producer`，`--name-only`＝`cmd/wisp/{instructions_200r2_test.go,panel_pump.go,run.go}`＋`internal/agent/instructions_test.go`＋`internal/panel/{composer.go,instructions_200.go,instructions_200_test.go,pump.go}`＋`internal/projctx/projctx.go`（9 枚路径）。<br>尺三（AC#1 那半"真到达请求"vs"字段恒空"的静态分界）：`git --no-pager grep -n "AttachProjectInstructions" HEAD -- internal/ cmd/ \| grep -v "_test.go"` → **2 命中**＝定义 `internal/agent/prompt.go:280 func (l *Loop) AttachProjectInstructions(src InstructionSource)`＋**生产调用点 `cmd/wisp/run.go:1071 loop.AttachProjectInstructions(instrLoader)`**；消费腿在 `prompt.go:241` 注释逐字 `Ticket 200: the project's own instructions join the system content, AFTER` ＋`:270 if a.Instructions == nil`／`:273 return a.Instructions.BlockForTurn()`。<br>尺四（AC#7 载体的生产者，票面/台账 `A408` 点名的"键在值恒缺"那一形）：`git --no-pager show HEAD:cmd/wisp/run.go \| sed -n '690,725p'` → 生产装配 `panel.NewSnapshotPump(panel.PumpSources{…})` 里 `:718 Instructions: rt.instructionBundle,` **在位**，而 `rt.instructionBundle` 的实现在 `cmd/wisp/panel_pump.go:107`（读 loader 的 `Last()`），`:99-105` 注释逐字 `It reads the loader's LAST bundle and shapes nothing: nil before a run has attached one (the pump then reports not_run, with a reason, rather than dropping the key)`。<br>尺五（结案凭据种类）：裁决表按**文件名**取＝`git ls-tree -r --name-only HEAD -- docs/evidence/s1 | grep -E "/200-"` → `200-project-instructions-r1.md`／`200-project-instructions-r2.md` 两枚在树（**内容本腿禁读**；同尺换 `/220-` → **0 枚**，票 220 名下没有裁决表）；票内注记＝票面 `:66-73`（`200-r2 更正上面那两行（走的是甲…）`）与 `:77-92`（`⚠ 09-28 20:5x 更正上面那一跳的形状`）；框数尺＝`grep -c '^- \[ \]' .scratch/wisp/issues/200-*.md` → **8 格全空** | ⚠ 这一枚**不写〔已失效／已被别人做掉〕**，理由要说死：做掉它的两枚 commit 标题逐字都写着 `ticket 200`／`ticket 200 (200-r2)`＝**这活是 200 自己做的**，不是别的票替它做的；所以它的正确形状是"缺陷不在了＋AC 框没翻"＝〔差翻勾〕。<br>**已判到的（凭字面，不靠"函数还在"）**：票面头名的"我们一枚代码都没有"今树**逐条反证**——五段形状里 §1 发现（优先级名册＋逐级向上）、§2 预算（走 D15 现成预算：`run.go:1055` 读 `cfg.Agent.ProjectInstructionsEnabled`、`internal/config/schema.go:442` 那枚开关键）、§4 可见（泵侧 `Instructions` 生产者已接）、§5 可关（同一枚键＋`pump.go:176` 注释逐字 `"this project has no AGENTS.md" indistinguishable on the wire` 说明 r1 那个谎已被处理）四段都有载体与生产者；§3 分层那一半静态可见（`projctx.go:72` 注释逐字 `what keeps a cloned repository's AGENTS.md from reading as an authority`）。<br>**没判到的（⛔ 不许被读成"已结"**）：①AC#1 的正控判据是"写一份带独特标记的说明文件 ⇒ **请求里**能逐字找到"＝真发请求的运行期读数，本腿零 go；我只证到 `run.go:1071` 的接线与 `prompt.go:273 BlockForTurn()` 的消费腿存在；②AC#3 超限截断的可读输出、AC#5 恶意说明样本、AC#6 去掉盖戳的变异必红——三格都要跑；③**票面 :56-65／:82-89 那一跳到界面那支**（`PanelSnapshot` 补 `instructions?:`）落点＝`frontend/src/lib/panel.ts`，本腿**禁读 `frontend/**`** ⇒ 那一格今天做没做**我量不到**，只能指出票面自述它当时是 `internal/panel` 的一枚在册红（尺＝`git --no-pager grep -n "TestApprovalCardViewJSONKeysMatchFrontendTypes" HEAD -- internal/panel/ \| grep -v _test` → **1 命中＝`internal/panel/approval.go:14`，而那一行是注释**（逐字 `// TestApprovalCardViewJSONKeysMatchFrontendTypes compares these JSON keys with`，说的是这枚尺"在"，⛔ 它的用例体在 `_test.go` 里、非我射程）；这一把**不构成任何判断**，只说明那枚在册红的名字今天在非测试码里只以注释存在）；④票面 `:91-92` 那条撤销口令（回 200-r2 的乙 ⇒ AC#7 当场退回）今天没被触发，也没有读数。<br>⇒ 编排者若要翻这 8 格，缺的是**四格运行期读数＋一枚界面侧复跑**，不是产码 |
| 220 | **仍成立**（票面**标题第一句**"那一格今天只看得见 L2"今树原样；标题第二句"第二张卡会静默掉格"已由 220-r1 自己做掉 ⇒ 半死半活，成立的是 L1 那一半） | 尺一（AC#1 那半：join 是否已改成双查＋`#序号`）：`git --no-pager grep -n "func indexWaitingKey" HEAD -- internal/panel/` → `subagent_roster_197.go:262`；`git --no-pager show HEAD:internal/panel/subagent_roster_197.go \| sed -n '/func indexWaitingKey/,/^}/p'` → 体内逐字 `set[key] = true` ＋ `if base, ok := splitClashSuffix(key); ok && base != "" { set[base] = true }`；`splitClashSuffix` 现读＝`LastIndexByte(id, '#')` ＋后段逐字符验 `'0'..'9'` ⇒ **`原id#序号` 那一形被并回基名**。建表侧：`git --no-pager show HEAD:internal/panel/subagent_roster_197.go \| sed -n '195,245p'` → `:212 waiting := make(map[string]bool, len(pending)+2*len(l1Windows))`／`:214 indexWaitingKey(waiting, card.CorrelationID)`／`:219-220 indexWaitingKey(waiting, w.CorrelationID)`＋`indexWaitingKey(waiting, w.TaskID)`。落地 commit＝`git --no-pager log --format='%h %ad %s' --date=format:'%m-%d' -- internal/panel/subagent_roster_197.go` → 最近一笔 **`480b970d 10-03`**，标题逐字 `220-r1: 名册那一格改成 corr 与 taskID 双查，并把同任务第二张卡接回来`。<br>尺二（AC#2 那半：Gate 有没有枚举口）：`git --no-pager grep -n "g\.windows\|windows\[" HEAD -- internal/agent/approval/gate.go` → **4 命中＝`:355`（openWindow 内的 clash 检查）／`:358 g.windows[w.corr] = w`／`:364 delete(g.windows, corr)`／`:427 w := g.windows[v.CorrelationID]`（Veto 按名查）**，与票面"只有那四行自己用"逐字一致；枚举口候选名全扫＝`git --no-pager grep -niE "OpenWindows\|PendingWindows\|ListWindows\|WindowWait" HEAD -- internal/ cmd/ \| grep -v "_test.go"` → **5 命中全在 `internal/panel/`**（`pump.go:115/127/200/323`＋`subagent_roster_197.go:198`），`internal/agent/approval` 侧**零**。<br>尺三（生产装配有没有给 `L1Windows` 供生产者）：`git --no-pager grep -n "L1Windows:" HEAD -- cmd/ internal/` → **恰好 1 命中＝`internal/panel/subagent_blocked_220_test.go:81`**；生产那处＝`git --no-pager show HEAD:cmd/wisp/run.go \| sed -n '690,725p'` 里 `panel.NewSnapshotPump(panel.PumpSources{ Verdicts:… Mode:… Workspace:… Git:… Model:… Credential:… Results:… Instructions:… Tasks:… Out:… })` **没有 `L1Windows` 这一项**。<br>尺四（nil 的语义＝票面 197 旧行为，不是"没在等"）：`git --no-pager show HEAD:internal/panel/pump.go \| sed -n '190,215p'` → `:198-200` 注释逐字 `A nil reader is the honest "this host has no enumeration for that dimension", and it changes one thing only: the roster's blockedOnApproval keeps reading from cards alone, i.e. exactly ticket 197's behaviour.`；消费点 `pump.go:323-325 var waits []L1WindowWait; if p.src.L1Windows != nil { waits = p.src.L1Windows() }`。<br>尺五（第二枚卡改写处，票面引 `queue.go:152-154`）：`git --no-pager show HEAD:internal/agent/approval/queue.go \| sed -n '163,172p'` → `:167-168` 逐字 `if _, clash := q.byID[corr]; clash { corr = fmt.Sprintf("%s#%d", d.CorrelationID, q.seq) }` ⇒ **改写点原样还在**，现在是被 `splitClashSuffix` 追上，不是被删。<br>框数尺＝`grep -c '^- \[ \]' .scratch/wisp/issues/220-*.md` → **5 格全空** | ⚠ 两半必须分开写，否则就是把活读漏：<br>**死的那半**＝"同一枚任务第二张卡静默掉格"：`indexWaitingKey` 把 `#序号` 的基名也索引，且**卡的 corr 与窗口的 corr／taskID 三路同进一张 `waiting` 表**，票面 AC#1 要求的"corr 与 taskID 双查＋把 `#序号` 纳入匹配"今树**已到位**；落地 commit `480b970d` 标题逐字带 `220-r1`＝**这活是 220 自己那一支腿做的**。<br>**活的那半**＝"那一格今天只看得见 L2"：AC#2 的两支选择（甲＝给 `Gate` 补只读枚举口，⚠ 新增方法名＝要台账批准；乙＝让 L1 窗口进 `LiveApprovals()` 视图）**两支都没走**——`gate.go` 今天仍只有那四行自用、零枚举口；泵侧倒是**多了一枚 `L1Windows` 读槽**（`pump.go:200`），但**生产装配没往里放东西**（尺三＝全仓仅 1 处赋值、且在 `_test.go`），于是 `pump.go:198-200` 自陈的那句"roster 继续只从卡里读＝票 197 的行为"就是今天的真值。⇒ **槽开了、水没接**，这不是措辞问题：它正是票面"为什么这不是措辞问题"那节说的"一个正在等的任务读起来和没在等一模一样"，只不过今天剩下的成因从"join 用错键"换成了"没人枚举 L1"。<br>**量不到的面**：AC#2 的判据本身（"种一枚正在 L1 短窗口里等的任务 ⇒ 那一格必须为真，改前读数当场跑出来贴回票上"）与 AC#3 的两形正控＝运行期＋产码；AC#5 明写要 `go test ./cmd/wisp ./internal/...` 整包终态逐名比红名集合＝**本腿零 go**，⛔ 不据"代码形状"外推那一格。<br>另记两枚尺外事实：①票面 `:12` 那行"卡的来源只覆盖 L2"（`LiveApprovals()` 只走 `q.pending`）今天**没复查**——我读的是 `cmd/wisp/panel_pump.go:62 items := rt.gate.Queue().LiveApprovals()`，生产者仍只有队列一侧，与票面一致；②票面 `:35-38` 那条"外部对照"设计约束（`blockedOnApproval` 必须是附加维度、不许覆盖运行态）今树**有答复**：`pump.go:319-322` 注释逐字 `Neither source REPLACES the other, and neither one touches the row's filed D43 state - blockedOnApproval stays the added dimension the ticket's 外部对照 constraint bounds it to be.` ⇒ 那一支的**口径**已被写住，但写住它的同时 L1 仍无生产者，所以"点亮之后会不会盖掉还在跑"今天**无从验**（要等水接上）。 |
| 225 | **仍成立**（票面头名的两件事——**账对不上**＋**没有任何仪器在查**——今天都还在；而 09-29 那格审计的结案凭据确实是票内注记＋一枚非实现者裁决表文件名，所以这一枚同时带〔差翻勾〕形状的一半，⚠ 见读数末段） | 尺一（票面"没有仪器"那句逐字复跑）：`git --no-pager grep -rln "DEFERRED(" HEAD -- tools/ scripts/ .github/` → **只有 1 枚文件＝`scripts/portable-tests.sh`，且两处都是注释**（`:50` 与 `:106`，内容为 `doc.go-only boundary stub (DEFERRED: ticket 42)` 一类散文），`tools/` 侧 **0 命中**（含 `tools/d22scan`）；`git --no-pager grep -n "SPEC-12\|registry" HEAD -- tools/d22scan/main.go` → **1 命中＝`:310`**，读的是 `permanent 0-file walk while the footer advertised it (registry A22/A26/A30)`＝**台账名，不是 §5 登记表** ⇒ "今天没有自动检查这一条"这句**仍真**。<br>尺二（分母，本腿自己现量、⛔ 不照抄票面也不照抄派单）：代码侧标记 `git --no-pager grep -o "DEFERRED(" HEAD -- cmd internal tools \| wc -l` → **30 处**（票面 09-29 锚 `36b46125` 那版＝26 处／22 枚去重；票面 `:42` 编排者 10-04 现跑＝30 处 ⇒ **我今天与 10-04 那把对上，与票面 09-29 那把差 4**）；登记表侧 `git --no-pager show HEAD:docs/specs/SPEC-12-roadmap-governance.md \| grep -cE "^\| DEFERRED"` → **13 行**（票面记 12，多的那枚是 `:80` 的〔确认窗口借用裸 `Esc`（票 245 乙形残余）〕，后补）、`grep -cE "^\| RESERVED"` → **7 行**（与票面一致）。⇒ **30 vs 13，账仍然对不上**。<br>尺三（票面点名的两枚"确凿失配"逐字复算）：`git --no-pager show HEAD:internal/config/doc.go \| grep -n "DEFERRED"` → `:14` 逐字 `// DEFERRED(schema/hot-reload/migration): implemented by ticket 05. This` **原样在**；`git --no-pager show HEAD:internal/watchdog/doc.go \| grep -n "DEFERRED"` → `:18` 逐字 `// DEFERRED(watchdog loop/thresholds): implemented by ticket 42. This ticket` **原样在**。<br>尺四（补充①：规格要求的载体）：`git --no-pager show HEAD:docs/specs/SPEC-12-roadmap-governance.md \| sed -n '92,98p'` → 逐字 `**双向交叉核对**：代码内每个 \`// DEFERRED(D-xx): … → docs/DEFERRED.md#锚点\` 必须在表有对应条目，反向亦查（对抗验收逐片执行）。`；载体存在尺＝`git ls-files docs \| grep -E "DEFERRED\|DECISIONS"` → **空**；带后缀的标记尺＝`git --no-pager grep -n "DEFERRED(" HEAD -- internal/ cmd/ tools/ \| grep "docs/"` → **0 命中**；规格拼写的标记尺＝`git --no-pager grep -c "DEFERRED(D-" HEAD -- internal/ cmd/ tools/` → **空输出（rc=1，0 处）** ⇒ **没有一枚标记长成规格要求的样子**，且"必须在表有对应条目"那句今天**没有任何仪器执行**。<br>尺五（补充②那枚类型词相反）：`git --no-pager grep -n "REJECTED: D13\|KindMCP" HEAD -- internal/tools/registry.go` → `:55` 逐字 `{Kind: KindMCP, Ticket: "REJECTED: D13/16.9#7, interface slot only"}`；§5 侧 `git --no-pager show HEAD:docs/specs/SPEC-12-roadmap-governance.md \| grep -nE "^\| (DEFERRED\|RESERVED\|REJECTED)" \| grep -i mcp` → `:81` 逐字 `| RESERVED | MCP client | D13 与 host bridge 收口冲突（16.9#7 驳回引入） | …` ⇒ **两处给后续人的指令仍相反**。<br>尺六（票面 AC#3 那枚 config 标记的成因半边已被别人答复）：`git --no-pager grep -n "CheckAndReload" HEAD -- cmd/ internal/ \| grep -v "_test.go"` → 生产侧命中 `cmd/wisp/config_reload.go:153 rep, err := rt.mgr.CheckAndReload()`、`cmd/balldebug/main.go:243`（票面点名的那一枚仍在），且 `cmd/wisp/run.go:320` 注释逐字 `the first production caller of Manager.CheckAndReload in \`wisp run\``、`internal/config/manager.go:17` 逐字 `A host's tick calls CheckAndReload - since ticket 223 that tick is` ⇒ **票面"热加载在生产里根本不转"这半句今天过期**（票 223 接上了）。<br>框数尺＝`grep -c '^- \[ \]' .scratch/wisp/issues/225-*.md` → **4 格全空**；裁决表按文件名取＝`git ls-tree -r --name-only HEAD -- docs/evidence/s1 \| grep -E "/225-"` → **`225-deferred-registry-audit.md` 在树**（内容禁读） | ⚠ **这一枚的"档"要看清判的是哪一句**：票面头名两件事（**标记与 §5 对不上**／**没有仪器在查**）我逐把复算都还在（尺一、尺二、尺四），所以写〔仍成立〕。**同时**票面 `:28-35` 那块 ✅ 注记说 AC#1/AC#2/AC#3 的**审计结论**已随腿 `225-a1` 交回（表一 3／3／12、表二 12 行里只有 2 行有标记、表三归口、AC#4 按住不建仪器），凭据种类＝**票内注记＋台账 `A429`＋一枚非实现者裁决表文件名** ⇒ 那三格是〔差翻勾〕的形状，**但按 §0 那一档的口径它要求"缺陷不在了"**，而本票的缺陷恰恰是"账对不上"这件事本身——审计把账目**量清楚了**，没有把账目**对平**（AC#1 自己就写死"不许改动 `docs/specs/**` 任何文字，要补行／要摘标记都写成待人批准的落点"）。⛔ 所以这一枚**不许被读成"审计做完了＝票 225 可结案"**。<br>**三处过期／漂移要具名**：①票面 §现量 的 **26 处／22 枚**已被**我自己这轮的 30 处**取代（与票面 `:42` 编排者 10-04 那把独立对上；28 与 30 那 2 枚差**仍在**，票面自己把它归给 AC#0/AC#1 的名册口径，本腿不复其差）；②登记表 **12→13 行**（新增 `:80` 那枚 Esc 借用行）；③票面"确凿失配"的第一枚**理由半边**（"生产里没人调 ⇒ 标记在骗人"）今天被票 223 答复了一半——`CheckAndReload` 已有 `cmd/wisp` 生产调用者，所以那句标记现在谎的是**"implemented by ticket 05"**这个归属，不是"没人实现"。票面 `:33` 表三早就把这枚的落点归给票 223，与此一致。<br>**待人拍板的两格（本腿不裁，属 AGENTS.md §2 未定义即停族）**：`docs/DEFERRED.md` 建不建（建＝新真相源；不建＝`SPEC-12:95-96` 那句形状规定与盘上不可执行，改规格文字＝契约变更须人工批准）；`registry.go:55` 的 `REJECTED` 与 §5 `RESERVED` 谁改谁。<br>**量不到的面**：AC#4 的"要不要上仪器"属裁不属量；票面 `:41` 那条 ⚠"未跑 d22scan 复认这一发"今天**仍未跑**（本腿零 go／⛔ 不许把它读成实测）。 |

### §2-0 分母尺（本腿自己现量，不照抄派单给的数）

命令原文：
```
for f in .scratch/wisp/issues/2*.md; do case "$f" in *-done.md) continue;; esac; \
  b=$(grep -c '^- \[ \]' "$f"); x=$(grep -c '^- \[x\]' "$f"); \
  if [ "$b" -gt 0 ] && [ "$x" -eq 0 ]; then echo "$(basename $f .md) :: 未勾$b"; fi; done
```
读数：`2*` 通配把 20–29／2xx 全吞进来，原始输出 30 行；按 `2[0-3][0-9]` 数值段收窄到 200–238 后得 **12 行**：
200(8)／220(5)／225(4)／227(6)／229(6)／230(5)／**232(6)**／233(6)／234(6)／236(7)／237(3)／238(4)。
剔除 **232**（派单 §3 明令不判：写腿 `232-r2` 正在修它）⇒ **本腿分母＝11 枚**，
与派单 10-06 11:4x 的现量（200／220／225／227／229／230／233／234／236／237／238）**逐枚对上，零差**。
⚠ 号段 239 以上不归本腿（另派）；269／249 不在本尺输出里（已改名或已非零勾，本腿不判、不追）。

## §3 判得心虚的枚数与具名理由

尚未开始填（本腿在判档进行中回填，⛔ 不由编排者代填）。

## §4 一处尺有歧义／判不动的地方

尚未开始填（本腿在判档进行中回填，⛔ 不由编排者代填）。
