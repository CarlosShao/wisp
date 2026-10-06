# 批 4f —— pool-validity 4f（零勾开放票复量：15 枚待判名册）

> 权威交付件＝本文件。同目录 `batch4f-skeleton.md` 是本腿起手写的同一份骨架的第一次落笔（当时把文件名写成了 `-skeleton` 后缀），
> 按本仓「临时件只建不删」的规矩留在盘上、不复用、不删除；一切读数以本文件为准。

## 0. 起手锚

- HEAD 短号：起手 `2eb75bb7`；16:52 复跑前进为 `b7251d17`（编排者在此期间落了台账枚）。本腿全部产码读数走 `git show HEAD:<path>`／`git --no-pager grep … HEAD`，两枚锚点之间若某枚读数变化会具名记两把。
- 时间：`10-06 16:52`
- `git status --porcelain -- internal cmd` 原文：

```
 M internal/tools/tasklist_deferred_236r3_teeth_test.go
```

⇒ 此刻 `internal/tools` 确有一枚写腿在跑（236-r3 的延迟用例件在工作树被改），本腿对它**零读零写**，且本腿全程零 `go` 命令。

## 1. 名册

### 1.1 抄自 `4e` §3.2 的 22 枚（只取那张表与尺原文，未读它的判语内容面）

出处：`.scratch/wisp/probes/pool-validity/4e/batch4e.md:164-168`

```
178 182 185 186 187 189 194 195 196 211 213 214 215 216 219 232 233 234 236 237 249 269
```

### 1.2 本腿要判的（扣除后 15 枚）

`178 182 185 186 187 189 194 195 196 211 213 214 215 216 219`

### 1.3 逐枚扣除清单（每枚给扣掉理由＋本腿自己跑的尺读数，⛔ 不抄派单那段话）

| 票号 | 扣掉理由 | 本腿的尺原文与读数 |
|---|---|---|
| 233 | 已由 `4e` 今天判掉（本腿不读其判语，只认"已判"这一事实） | 尺 `git --no-pager log --oneline -12` ⇒ 命中 `9796f15e pool-validity 4e: verdicts for 233, 234, 237 (the three live cells of dead leg 4a)`；再 `git --no-pager log -1 --format='%h %ad %s' --date=format:'%m-%d %H:%M' 9796f15e` ⇒ `9796f15e 10-06 16:27 …`＝**今天**、且一枚 commit 号标题逐字点名这三枚 |
| 234 | 同上 | 同上那一行标题逐字含 `234` |
| 237 | 同上 | 同上那一行标题逐字含 `237` |
| 232 | 此刻有写腿在做（`232-r2`／`232-v*` 类，判它会把在飞的活第二次立案） | 尺 `ls -d .scratch/wisp/probes/232` ⇒ 目录在场（本腿禁读该前缀，故只登记在场，一个字没读）；`git --no-pager log --oneline -8 --grep='232'` ⇒ `3d9b8374 probes(232-r2 收尾代提): 写腿撞 150 轮帽但产码与证件全在盘上，编排者只代提未代判`，`git --no-pager log -3 --format='%h %ad %s' --date=format:'%m-%d %H:%M' -- .scratch/wisp/issues/232-*.md` ⇒ `3d9b8374 10-06 13:12`＝**今天仍在动这枚票面**；框态尺（§5 第 2 条两把原文）＝6／7 |
| 236 | 此刻有写腿与终裁腿都在做 | §0 那把 `git status --porcelain -- internal cmd` 读数 ` M internal/tools/tasklist_deferred_236r3_teeth_test.go`；`git --no-pager log -2 --format='%h %ad %s' --date=format:'%m-%d %H:%M' -- .scratch/wisp/probes/236` ⇒ `ea5b2dc9 10-06 17:29 probes(236-v1 起手)`＝**本腿交件期间还在长**；`git --no-pager log --oneline -8 --grep='232'` 里同批含 `ee22fafc ledger(A637)：236-r2 交件收档…＋236-r3 补位`；框态尺＝7 枚未勾（`grep -c '^[[:space:]]*[-*] \[ \]'`） |
| 249 | 票面 `Status:` 行已就地标撤回，判它＝把已撤销的缺陷第二次立案 | 尺 `sed -n '3p' .scratch/wisp/issues/249-the-resident-process-installs-no-panic-sink-so-a-real-crash-leaves-no-stack-on-disk.md` ⇒ 含 `→ **WITHDRAWN（10-01 14:2x 编排者自撤；出处＝台账 A503；⛔ 任何腿不许据此开工）**`，且同句写明撤回原因＝立票时那句前提是错的 |
| 269 | 票面 `Status:` 行已就地标撤回／并案 | 尺 `grep -n 'WITHDRAWN\|并案' .scratch/wisp/issues/269-gofmt-step-is-held-red-by-its-own-committed-mutation-sample-and-eats-three-unprotected-steps.md` ⇒ `:3` 含 `★**WITHDRAWN／已并案（10-05 12:3x 编排者自撤，出处＝台账 A620`；`:41` 含 `并案后由票 236 的 AC#6a 口径承载；⛔ 不两票都记一次功劳` |

### 1.4 名册复量（本腿的框态尺，转义形；⛔ 未勾＝0 不等于无活，见 §3）

| 票号 | 未勾框 | 已勾框 | 行数 | 票面 Status 首句（截 60 字） |
|---|---|---|---|---|
| 178 | 5 | 0 | 31 | `ready-for-agent（只读＋改仪器）` |
| 182 | 7 | 0 | 100 | `ready-for-agent（只读普查，零产码）` |
| 185 | 7 | 0 | 69 | `已立、按住不派` |
| 186 | 9 | 0 | 54 | `ready-for-agent（先只读普查，再落地）`＝功能票 |
| 187 | 9 | 0 | 39 | `ready-for-agent（先只读普查，再落地）`＝功能票 |
| 189 | 6 | 0 | 25 | `ready-for-agent（先只读设计核，再落地）`＝功能票 |
| 194 | 7 | 0 | 55 | `ready-for-agent（先只读核代价，再分批落地）`＝功能票＋契约变更已获点名批准 |
| 195 | 7 | 0 | 68 | `blocked（排在票 33 片 B 交件之后…）` |
| 196 | **0** | 0 | 14 | `未派（先摆 owner 一次…）`＝无框票 |
| 211 | **0** | 0 | 63 | `在派队列里`＝无框票 |
| 213 | **0** | 0 | 60 | `待派（owner 09-28 22:1x 口头接单）`＝无框票 |
| 214 | **0** | 0 | 48 | `待派…接线＋补最小缺口`＝无框票 |
| 215 | **0** | 0 | 47 | `待派…`＝无框票 |
| 216 | **0** | 0 | 33 | `待派`＝无框票 |
| 219 | **0** | 0 | 100 | `接单＝做`＝无框票 |

尺原文（逐枚）：

```
f=$(ls .scratch/wisp/issues/<N>-*.md | head -1)
grep -c '^[[:space:]]*[-*] \[ \]' "$f"
grep -c '^[[:space:]]*[-*] \[x\]' "$f"
wc -l < "$f"
```

★15 枚里 **7 枚（196／211／213／214／215／216／219）框态读数为 0**，与派单 §1 的口径陷阱一致：它们的判据写成散文式粗体子弹，"未勾格＝0"只等于"没被框化"，⛔ 不等于没有活。这七枚的判法在 §3。

## 2. 判档表

> 口径＝派单 §2 四档；复法列一律本腿亲手跑过的命令原文（⛔ 不引别人读数当我的）。
> 全部产码读数走 `git … HEAD`／`git show HEAD:<path>`——工作树里 `internal/tools/tasklist_deferred_236r3_teeth_test.go` 是别人在飞的脏件，本腿一个字没读。
> 锚点：起手 `2eb75bb7` ⇒ 判档时 `b7251d17`；每把尺在 `b7251d17` 上跑，两锚之间可能有差的枚在 §5 具名。

| 票号 | 档 | 复法命令原文 | 读数 |
|---|---|---|---|
| 178 | **仍成立** | `sh .scratch/wisp/probes/154/gate-clauses.sh`（只读尺；`git show HEAD:.scratch/wisp/probes/154/gate-clauses.sh \| grep -cE 'go (test\|build\|vet\|run\|list\|env)'`＝**0** ⇒ 零 `go` 命令，本腿跑它不违边界；锚默认 `A="${1:-$(git rev-parse HEAD)}"`＝今 HEAD） | 聚合段逐字：`# BAD  腿=G6neg 声明=ring 基线=1枚 实测=3枚 因=新增未成对（票 171 AC#2：实测 > 基线）`；`# 腿数＝14 声明与实测不符＝2`；`# 聚合退码＝2`。仪器本体尺：`git show HEAD:.scratch/wisp/probes/154/gate-clauses.sh \| grep -n '178\|四堆\|只关\|别处开'`＝**0 命中**；`want G6neg ring rev`＝`:551`、`want_n 1`＝`:552`。⇒ 票面那枚"关的是别处开的那一枚，被仪器记成新增未成对"的缺陷**今天还在、且已从 1 枚涨到 3 枚**（本腿自己按文件算的反向未成对集＝`cmd/wisp/run.go`＋`cmd/wisp/task_scope_close_151_test.go`＋`internal/tools/pointer_183_cli_seam_test.go`＋`internal/tools/subagent_197.go`＋`internal/tools/ticket175r2_stamp_live_test.go` 五枚文件，其中 `subagent_197.go:332` 与 `run.go:942` 是注释行、被仪器的分子剔掉 ⇒ 实测 3 枚＝`run.go`＋`pointer_183_cli_seam_test.go`＋`ticket175r2_stamp_live_test.go`）。★另一枚 BAD＝`G1b`（`.Execute(` 那一族），与 178 无关，本腿不判、只具名登记，别让下一位把退码 2 全记到 178 头上 |
| 182 | **差翻勾** | `git show HEAD:.scratch/wisp/probes/182/a1/census.md \| wc -l` ＋ `git ls-files 'docs/evidence/s1/182-*'` ＋ `grep -n '^## A566' docs/reports/pending-and-issues.md` ＋ 框态尺 | 普查件在场＝**275 行／51,828 字节**，单文件 commit `24026880`（`git --no-pager show --name-only --format='%h %ad' 24026880`＝`10-03 10:49`，文件只有那一枚件）；裁决表 `docs/evidence/s1/182-rail-stacks-census-a1.md` 在库；台账 `A566` 在 `docs/reports/pending-and-issues.md:11159`，其 §3（`:11166`）逐字收下"12 堆＝有源 5 整＋1 半／无源有规格 2 整＋1 半／规格真空 2 整＋1 半／接了半根线 1 整；12 堆全指认到 13 枚既有票"。凭据种类＝**票内 Progress log 注记＋台账 `A##`＋证据表**三样齐；框态尺＝`grep -c '\[x\]' .scratch/wisp/issues/182-*.md`＝**0**（历史四枚 commit 上逐枚 `git show <c>:<file> \| grep -c '\- \[x\]'` 全＝0，从未被勾过）。⇒ AC#1–AC#5 的产物都在、**一格没翻**；⛔ 翻勾归编排者，本腿不裁它可结案 |
| 185 | **已失效（已被别人做掉）** | `git --no-pager log --oneline --diff-filter=A -- internal/risk/hostpath_185.go` ＋ `git --no-pager show --name-only --format='' ee1e2118` ＋今树在场尺 | 落地枚＝**`ee1e2118 185-r1: 分页续读——豁免改绑本 scope 宿主亲手落盘的那一枚文件（A381 第四支·名册为派生投影）`**；name-only 命中 `internal/risk/hostpath_185.go`、`internal/risk/pointer_185_test.go`、`internal/tools/pointer_185_cli_seam_test.go`、`internal/risk/provenance.go`、`internal/risk/taintmatch.go`、`docs/evidence/s1/185-paged-reread-r1.md`。今树：`git show HEAD:internal/risk/pointer_185_test.go \| grep -c '^func Test'`＝**10**，逐枚含 `TestPointer185SecondRereadOfAHostArtifactIsClean`（AC#1 那一形）、`TestPointer185SameBasenameOtherDirectoryStillHits`（AC#7 那一形）、`TestPointer185ForeignContentCarryingTheHostPathStillHits`（AC#1 成对的反那一形）、`TestPointer185CostFaceOfTheRosterIsPinnedAndNarrowerThanTheSpan`（AC#4 的代价面）；`git show HEAD:internal/tools/pointer_185_cli_seam_test.go \| grep -n 'func TestPointer185AnotherTaskCannotBorrowTheRoster'`＝`:330`；`git show HEAD:internal/risk/hostpath_185.go \| sed -n '3p'` 逐字写 "Ticket 185 (ledger A381) — the host-minted-path roster, DERIVED and not stored."。⇒ 硬门满足（一枚 commit 号＋name-only 命中＋今树实现件在场复算尺）。票框 7 枚仍 0 已勾（本腿不动框） |
| 186 | **仍成立** | `git --no-pager grep -n 'GitSwitchBlockedReason' HEAD -- internal/panel/git.go` ＋ `git show HEAD:cmd/wisp/panel_inbound.go \| grep -n 'Workspace:'` ＋ `git --no-pager grep -n 'WorkspaceWriteHandler' HEAD -- internal cmd \| wc -l` | 快照里那句自陈逐字（`internal/panel/git.go:82`）＝「切换分支／切换工作树今天不可用：网页到宿主的那一跳还没有落地…」，字段 `:144 SwitchBlocked string`；`cmd/wisp/panel_inbound.go:275`＝`Workspace:  nil,`（紧邻注释 `:274` 逐字 `workspace = ticket 186`）；`WorkspaceWriteHandler` 全仓命中＝**0**；`RequestWorkspaceSwitch` 剥测试后调用点＝0（`git --no-pager grep -cw RequestWorkspaceSwitch HEAD -- cmd/wisp/panel_inbound.go internal/panel/composer_dispatch.go internal/panel/composer_handlers.go` 只命中 `composer_dispatch.go:1`＝接口声明行）。⇒ 交付 1「认出 git」＋交付 2「列出可去的地方」**已由票 181 落进同一枚读面**（`internal/panel/git.go:136 Worktrees []GitWorktree`、`:141 Branches []string`；commit `fccfe752 181-r1: go-side git read surface + git section on the composer snapshot`），交付 3「切过去」＝零；票 AC#9 点名的需求表尺 `git ls-files 'docs/evidence/s1/*panel-ask*'`＝**0 枚**（186／187 两票都点名要它） |
| 187 | **仍成立** | `git --no-pager grep -niE 'thinking' HEAD -- internal/panel internal/config/settings.go \| grep -v '_test.go'` ＋ `git --no-pager grep -n 'FieldRoleChatModel' HEAD -- cmd internal \| grep -v '_test.go'` ＋ `git --no-pager grep -n 'SetRoleChatModel' HEAD -- cmd/wisp/panel_config_store.go` | 思考档位那一半：面板与写口**零命中**（唯一命中 `internal/panel/pump.go:40` 是一句讲耗时词表的注释，不是档位）；`ThinkingIntensity` 非测试命中 18 处全在 config/agent 读侧，无面板写腿；`ThinkingLevels` 非测试只 3 处（`internal/config/schema.go:357/:359`＋`validate.go:238`），AC#4 要的"面板只列该模型声明支持的档"＝零。模型那一半**已在树**：`internal/panel/config_handlers.go:63 FieldRoleChatModel = "role_chat_model"`＋`cmd/wisp/panel_config_store.go:216-217 case panel.FieldRoleChatModel: written, err = s.mgr.SetRoleChatModel(w.Value)`，落地枚＝**`0d87a681 ticket248(248-r1 产码)：设置那条路从入向请求修到磁盘——AC#1 名册真扩（config.get/config.set 进 bridge.go 白名单…）`**。⇒ 两半里只有一半被别人做掉；框 9 枚全 0 已勾 |
| 189 | **仍成立** | `git --no-pager grep -niE 'git diff\|numstat' HEAD -- internal cmd \| grep -vc '_test.go'` ＋ `git --no-pager grep -niE 'uncommitted' HEAD -- internal cmd \| wc -l` ＋正控 `git --no-pager grep -c 'SwitchBlocked' HEAD -- internal/panel/git.go` | 两把零命中尺＝**0** 与 **0**；同形正控命中 **6** ⇒ 尺吃得住，那两个 0 是"这物不在库里"不是"射程外"。`git show HEAD:internal/panel/git.go \| sed -n '117,150p'` 的字段名册里没有未提交那一维。外部 `git` 进程尺（`git --no-pager grep -nE 'exec\.Command\("git"\|git", "diff' HEAD -- internal cmd`）＝5 命中**全在 `cmd/wisp/panel_host_gate_test.go`**（`rev-parse`／`archive`），生产码零 ⇒ AC#1 那句"diff 很可能需要起进程"仍未被任何交付回答。设计核件在库＝AC#1 有产物（`docs/evidence/s1/189-uncommitted-count-design-a1.md`）；框 6 枚全 0 已勾 |
| 194 | **仍成立** | 逐枚现跑（⛔ 不按名字猜）：`for m in panel.resync tasks.list task.detail history.query transcript.get approval.current approval.queue approval.decide config.get config.set grants.list grants.revoke privacy.purge privacy.export cost.summary models.list models.delete diagnostics.export; do git --no-pager grep -nF -e "$m" HEAD -- internal/panel cmd/wisp \| grep -v '_test.go' \| wc -l; done` ＋ `git show HEAD:internal/panel/bridge.go \| grep -nE '^\s+Method[A-Za-z]+ +='` | 18 枚规格名里非测试命中只 3 枚：`approval.decide=1`／`config.get=3`／`config.set=5`，其余 **15 枚全＝0**；代码那份名册＝`bridge.go:42/43/44/45`（四枚 `panel.*`）＋`:66/67`（`config.get`／`config.set`，票 248 落的）＝**6 枚 vs 规格 18 枚，交集仍只有 248 那两枚**。`git show HEAD:docs/specs/SPEC-08-ui-ball-panel.md \| sed -n '160,176p'` 那张表今还在（含 `approval.decide` 那行「「allow」拒绝一切面板来源（F2）；仅 reject 可面板发起」）。⇒ 两份名单**仍互不相认**；代价表在库＝AC#1 有产物（`docs/evidence/s1/194-method-roster-census-c1.md`）；框 7 枚全 0 已勾。⚠ 票 AC#5 那句前提本腿复量已过期，见 §5 第 3 条 |
| 195 | **仍成立** | `git show HEAD:internal/config/settings.go \| grep -cE 'confirm\|hook\|Hook'` ＋ `git show HEAD:internal/config/settings.go \| sed -n '259,264p'` ＋ `git --no-pager grep -nE 'Direction' HEAD -- internal/config/manager.go \| grep -v '_test.go'` | 程序化写入口 `writeOneKey`（`internal/config/settings.go:259`）签名后第一句就是 `m.mu.Lock()`，整枚文件里 `confirm`／`hook` 命中＝**0**；方向判定只在 reload 那一侧（`manager.go:250 riskDirection`、`:255 fsDirection`、`:260 netDirection`、`:265 pluginsDirection`、`:75` 注释 "Direction is the strongest detected posture (loosen wins over tighten)"、`:315 planLocked`）。⇒ 正是编排者 10-02 裁定 §1（票面 `:51`）收缩后的那枚真缺口——「把方向裁决从 reload 侧搬进程序化写入口」——**还没搬**；AC#1 那半**已由票 248 满足**（`git show HEAD:cmd/wisp/panel_config_store.go \| grep -c 'case panel.Field'`＝**7**，与票面 §5 那句"7 处"同值）；框 7 枚全 0 已勾 |
| 196 | **仍成立**（无框票，判法见 §3.1） | `git --no-pager grep -nE 'CHECK' HEAD -- internal/memory/schema.go` ＋ `git --no-pager grep -nE 'return "(done\|cancelled)"\|State:     "running"' HEAD -- cmd/wisp internal/agent \| grep -v '_test.go'` ＋ `git --no-pager grep -niE 'func (MapTaskState\|ToD43Name\|StateMapping\|CanonicalState)' HEAD -- internal cmd \| wc -l` | schema 尺只命中**一句注释**（`internal/memory/schema.go:16` 逐字"…NOT as SQL CHECKs, precisely so…"）＝约束**不存在**；写侧两套词照旧并存（`cmd/wisp/run.go:1209/1211/1213`、`internal/agent/loop.go:994/998/1004`、`internal/agent/journal.go:131`）；D43 那 20 枚名在 `internal/statemachine/states.go:11-31`；甲支映射函数命中＝**0**。旁证这枚缺陷还活着：`internal/tools/subagent_197.go:32` 注释今天逐字写着"…(done/cancelled/running/succeeded) are ticket 196's and never enter here"＝别的票在引用一枚未结案的缺陷 |
| 211 | **已失效（已被别人做掉）** | `git --no-pager log --oneline -S'MaxConcurrentSubagents = 4' -- internal/tools/subagent_197.go` ＋ `git --no-pager log -1 --name-only --format='%h %s' 7ea14ce3` ＋今树四把尺 | 落地枚＝**`7ea14ce3 197-r1b 票 197/211 甲：池的诚实数 8->4（= 桥的 D38d 天花板）＋常驻判据＋两枚桥上路＋gofumpt`**，name-only＝`internal/tools/subagent_197.go`＋`internal/tools/subagent_197_test.go`。今树逐条对 AC：AC#1＝`internal/tools/subagent_197.go:85 MaxConcurrentSubagents = 4` ＝ `internal/tools/bridge.go:24`／`internal/agent/budgets.go:47` 的 4；AC#2＝`subagent_197_test.go:411 Test197SubagentPoolNeverExceedsBridgeCeiling`（判据体是 `if MaxConcurrentSubagents > MaxToolConcurrency` 才红，不是恒真断言）＋`:426 Test197SpawnDescriptionNamesTheRealPoolCap`；AC#3＝`:534 Test197FullPoolRefusesNextSpawnWithReadableReason`＋`subagent_197.go:275` 的理由串带 `Roster.RunningSubagentIDs()`；AC#4＝`:449 Test197SubagentPoolCapsAtBridgeCeiling` 且 `:443-444` 逐字"spawnDirect is gone rather than documented"（`grep -n spawnDirect` 现在只剩这两行注释）；AC#5＝两枚 4 在 HEAD 未动。⇒ 甲支整条落地；乙支（抬契约）本腿不判、仍是要人工批准那一格 |
| 213 | **仍成立**（无框票，判法见 §3.2） | `for p in slashcommand commandcatalog CommandRegistry '/compact' '/goal'; do git --no-pager grep -niF -e "$p" HEAD -- internal cmd \| wc -l; done` ＋正控 `git --no-pager grep -niE 'Registry' HEAD -- internal cmd \| wc -l` | 五把尺逐枚＝**0／0／0／0／0**；同形状正控 `Registry`＝**676**、`catalog`＝**115** ⇒ 尺命中得了真名，那五个 0 是"这物不在库里"（维度已穷举，见 §5 第 4 条）。载体侧：`git ls-files 'internal/panel/*command*' 'internal/agent/*command*'`＝**0 枚**；`git --no-pager grep -n '^func ' HEAD -- internal/agent/control.go` 只有 `:51 MatchControl`，生产读者唯一一处＝`internal/agent/loop.go:348` ⇒ 落点候选仍在、名册一枚没建。普查件在盘＝`.scratch/wisp/probes/213/c1/inventory.md`（`ls -d .scratch/wisp/probes/213` 存在；本腿只按目录名取，未读其判语） |
| 214 | **仍成立**（无框票，判法见 §3.3） | `git --no-pager grep -n 'NewAttachmentBroker' HEAD -- internal cmd \| grep -v '_test.go'` ＋ `git show HEAD:internal/panel/pump.go \| sed -n '269p'` ＋ `git show HEAD:cmd/wisp/panel_inbound.go \| grep -n 'Attachment:'` ＋ `git --no-pager grep -ciE 'attachment' HEAD -- internal/config` | 剥测试后 `NewAttachmentBroker` 只剩它自己的定义与返回体（`attachments.go:179/:189`）＝**生产零调用者，未改**；`git show HEAD:internal/panel/attachments.go \| wc -l`＝**410**（票面"410 行"逐字对得上）；`pump.go:269` 逐字 `composer := NewComposerState(mode, workspace, nil, maxAttachment)`＝第三参仍硬写 `nil`；`panel_inbound.go:276`＝`Attachment: nil`（注释 `:274` 逐字 `attachment = ticket 92`）；配置层那把尺（连测试一起数）＝**空返回**，同尺正控 `provider` 在 `internal/config` 命中 11/43/23 ⇒ 那个空是真"零枚附件键"。⇒ **接线一格没接** |
| 215 | **仍成立**（无框票，判法见 §3.4） | `git --no-pager grep -n 'ListPluginStates' HEAD -- internal cmd` ＋ `git ls-files 'internal/plugin/*'` ＋ `git --no-pager grep -c Skill HEAD -- internal/config/schema.go` ／ `-c Plugin` ＋ `git show HEAD:internal/panel/config_handlers.go \| grep -n lockedFieldFamilies` | 唯一枚举面 `internal/memory/dao_misc.go:243 ListPluginStates` 的调用者只有它自己的注释行（`:242`）与 `internal/memory/dao_test.go:451`＝**生产零读者**；`git ls-files 'internal/plugin/*'`＝只有 `disposal.go`／`disposal_test.go`／`doc.go`＝零"列出已装"面；`Skill` 在 schema 尺＝**空返回**（同尺 `Plugin`＝8）⇒ 215b 那半连配置键都不存在；`config_handlers.go:108` 逐字 `var lockedFieldFamilies = []string{"risk.", "fs.", "net.", "plugins.", "privacy.", "models.", "audio."}` ⇒ 面板写腿对插件节**整族拒绝**，开关一格零。`git show HEAD:internal/config/unwired.go \| wc -l`＝**147**＝AC#5 那面旗没被拆（守卫在，不是交付） |
| 216 | **仍成立**（无框票，判法见 §3.5） | `git --no-pager grep -n 'unicode.IsControl' HEAD -- internal cmd` ＋ `git show HEAD:internal/panel/attachments.go \| sed -n '307,313p'` ＋ `git --no-pager grep -niE 'sanitiz\|control.?char\|escape.*terminal\|u001b' HEAD -- internal cmd \| wc -l`（票面 AC#1 那把原尺）＋ `git ls-files 'docs/evidence/s1/216-*'` | 全仓洗控制字符只有**一处**＝`cmd/wisp/approval_reply.go:615 if unicode.IsControl(r) {`，洗的是**操作员自己打的**理由（不是 216 管的"非用户手打"那一族）；面板那枚 `sanitizeDisplayName`（`attachments.go:310`）函数体逐字只有 `return filepathBaseLike(strings.TrimSpace(name))`＝裁路径段＋去空白，**不洗控制字符**（其上游 `:301` 是把含 `\x00\r\n\t` 的名字整枚拒收＝拒绝不是洗）；票面普查尺命中＝**39 行**（含测试）／非测试 20 行；AC#1 要求的名册件＝`git ls-files 'docs/evidence/s1/216-*'`＝**0 枚**、`ls -d .scratch/wisp/probes/216`＝不存在。⇒ "显示前统一洗"那枚共享底座**不存在**，五类字段（附件名／插件名／技能名／目录名／命令说明）零覆盖 |
| 219 | **仍成立**（无框票，判法见 §3.6） | `git show HEAD:internal/agent/approval/replies.go \| grep -nE 'func \(r \*Replies\) (Allow\|AllowSession\|Reject\|PanelReject\|PanelAllow\|Veto)'` ＋ `git show HEAD:internal/agent/approval/gate.go \| sed -n '712,730p'` ＋ `git show HEAD:internal/agent/approval/queue.go \| grep -n 'func (q \*Queue) allow'` ＋ `git show HEAD:internal/memory/schema.go \| sed -n '69,85p'` | 答复侧函数**已在**：`replies.go:316 Allow(ctx, corr)`／`:351 AllowSession`／`:372 Reject(corr, reason)`／`:378 PanelReject(corr, reason)`／`:431 PanelAllow`／`:455 Veto`。⚠ **允许侧没有理由位**：`Allow` 只两参；`gate.go` 的 `Request.Reason`（`:719`）只在 `if !r.Allow { return g.q.reject(corr, r.Reason) }` 那一支被读，允许一支走 `g.q.allow(corr, r.Grant)`，而 `queue.go:365 func (q *Queue) allow(corr, nonce string)`＝**两枚参数**＝票面 §2 第 3 行"真正缺的两样之一"今天仍缺。用户回看那一格：`tool_call` 表 14 列（`schema.go:69-83`）里**没有理由列**＝票面更正⑥那格的对撞今天仍未解。"长期"那一支的撤销面：`git --no-pager grep -n 'SetAllowedDirs' HEAD -- internal cmd \| grep -v '_test.go'`＝只有 `allowdirs.go:95/:109` 的定义与一句注释（`session/grants.go:43`）＝**有实现、有测试、生产调用者 0**，与票面 10-02 更正第 27 行同值。三枚按钮的真身（球旁卡／托盘／原生窗口）＝**量不到**，见 §4 第 5 条 |

**判档小结（本腿自数）**：仍成立 **12** 枚＝178／186／187／189／194／195／196／213／214／215／216／219；差翻勾 **1** 枚＝182；已失效 **2** 枚＝185／211；整枚量不到 **0** 枚（186／187／213／214／215／219 的界面那半是**格级**量不到，票判"仍成立"并具名哪一格量不到，见 §4）。12＋1＋2＝**15**，与 §1.2 名册枚数一致。

## 3. 无框票的散文式 AC 逐条在场尺

> 这七枚（196／211／213／214／215／216／219）票面**一格复选框都没有**：`grep -c '^[[:space:]]*[-*] \[ \]'` 与 `…\[x\]` 双 0（§1.4），判据写成散文式条目。尺原文与枚数：

```
for n in 196 211 213 214 215 216 219; do f=$(ls .scratch/wisp/issues/${n}-*.md|head -1);
  echo "$n acbullets=$(grep -cE '^- \*\*AC#' "$f") acany=$(grep -cE 'AC#[0-9]+' "$f")"; done
→ 196 acbullets=0 acany=0 ｜ 211 5/6 ｜ 213 5/7 ｜ 214 6/6 ｜ 215 5/5 ｜ 216 5/6 ｜ 219 6/12
```

- ⛔ 本腿不给它们记"已失效／无活"，除非能指出该票点名的产物今天在树上（下面每条给 commit 号或今树在场尺）。
- ★**票 196 连 `AC#` 子弹都没有**（`acbullets=0 acany=0`，全文 14 行），它的判据写在标题、现量表和"要拍的（不答＝默认甲）"那三行里 ⇒ 本腿按同一纪律把它的**内容**拆成条给尺，不因为它"既没框也没 AC 号"就跳过。
- ★**若要把这七枚补成带框判据，那是票面改动，归编排者排；⛔ 本腿不动框、不翻勾、不改 `Status:`、不追加 Progress log。**

### 3.1 票 196（两套状态词表并存／schema 不拦）——四条

| 条 | 票面原句要点 | 本腿的在场尺（HEAD 上跑） | 读数与归口 |
|---|---|---|---|
| 现量① | D43 名字的仓内逐字拷贝在 `internal/statemachine/states.go:11-31` | `git show HEAD:internal/statemachine/states.go \| sed -n '8,31p'` | 在树；`:8` 自陈 "names exactly as in D43 / SPEC-08"，20 枚常量名逐枚可见（`StateFirstRun`…`StateWatchdogAlert`） |
| 现量② | 今天真写进库的状态词一枚都不在那张表上 | `git --no-pager grep -nE 'return "(done\|cancelled)"\|State:     "running"\|Status = "cancelled"' HEAD -- cmd/wisp internal/agent \| grep -v '_test.go'` | 在树且未统一：`cmd/wisp/run.go:1209/1211/1213`、`internal/agent/loop.go:46`、`:994/998/1004`、`internal/agent/journal.go:42/:131` ⇒ **缺陷本体仍在** |
| 现量③ | schema 有没有拦＝没有 CHECK | `git --no-pager grep -nE 'CHECK' HEAD -- internal/memory/schema.go` | 只命中 `:16` 一句**注释**（逐字"…NOT as SQL CHECKs, precisely so…"）⇒ 约束不存在，读数成立但不是交付 |
| 甲支（唯一权威映射）有没有产物 | 新增一枚 Go 侧映射，读写都走它 | `git --no-pager grep -niE 'func (MapTaskState\|ToD43Name\|StateMapping\|CanonicalState)' HEAD -- internal cmd \| wc -l` ＋旁证 `git --no-pager grep -n 'func (o TaskOutput) StateAnswer' HEAD -- internal/tools/task.go` | 映射函数＝**0**；邻近在树的是票 188 那道门（`internal/tools/task.go:158 StateAnswer()` ＋ `internal/tools/task_state_188_test.go:71-72` 那串 11 枚旧词，注释逐字 "legacyTaskStateWords is ticket 196's vocabulary… none of them may enter this dimension"）＝**只挡它自己那一维，不是全库映射** ⇒ 不能拿来判 196 已失效 |
| 乙支（加 CHECK）有没有产物 | 触 `D35`，需人工批准＋迁移路径 | 同现量③那把尺 | **0**＝未落；本腿不代它选支（票面那句"不答＝默认甲"归编排者） |

### 3.2 票 213（斜杠命令目录）——五条

| 条 | 票面要求 | 在场尺 | 读数 |
|---|---|---|---|
| AC#1 | 名册可枚举且每条有真执行者 | `for p in slashcommand commandcatalog CommandRegistry '/compact' '/goal'`（`git --no-pager grep -niF -e "$p" HEAD -- internal cmd \| wc -l`）＋ `git ls-files 'internal/panel/*command*' 'internal/agent/*command*'` | 五把＝**0/0/0/0/0**，正控 `Registry`＝676／`catalog`＝115；载体文件枚数＝**0** ⇒ 未起步 |
| AC#2 | 没能力的命令进不了名册（塞假命令必红） | 同上 | 名册都不存在 ⇒ 本格**未起步**（⛔ 不许写成"已失效"） |
| AC#3 | 档位与权限：默认最严档里要么拒要么先问，拒了要说为什么 | `git --no-pager grep -niE 'slash' HEAD -- internal/panel cmd/wisp \| wc -l` ＋剥测试 `… \| grep -v '_test.go'` 逐枚看 | 两把读数：**67 行**／剥测试 **4 行**，那 4 行逐字是 `cmd/wisp/providers.go:231-232`（"splits provider/model at the first slash"＝模型引用串的斜杠）、`internal/panel/assets.go:70`（斜杠路径）、`internal/panel/git.go:463`（斜杠命名的分支）＝**没有一枚是命令入口形状**；⚠ 这把子串尺**不是零命中**（60 行是 `filepath.ToSlash` 那一族），所以"斜杠"这个词在库里到处有，而"斜杠命令"零——见 §5 第 7 条我这枚写错又当场顶正的记录。⛔ 不由"档位机制已在"外推成这条能用 |
| AC#4 | 入向没接通前终态＝名册可见＋命令行可执行，并具名"界面点不了" | `git show HEAD:cmd/wisp/panel_host_windows.go \| grep -n 'bindErr := w.Bind(panelDispatchBinding'` ＋ `git show HEAD:cmd/wisp/resident_windows.go \| grep -n 'first non-test call of NewPanelManager'` | 入向那一跳**今天已存在**（`panel_host_windows.go:405` 绑定 `wispDispatch`；`resident_windows.go` 注释 `:141` 逐字 "Building the manager is the first non-test call of NewPanelManager in this repository"）＝票面 §2 第 2 行那句前置（账 `A410`"只有命令行路"）**已过期**，见 §5 第 3 条；但命令名册仍 0 枚 ⇒ 本格既不是"已交付"也不是"被前置卡住"，是**前置要重读** |
| AC#5 | 命令执行全部走现成审计（三档都写审计） | `git --no-pager grep -n 'Audit AuditFunc' HEAD -- internal/panel/composer_dispatch.go` | 审计腿在（`ComposerDispatch.Audit`），但没有命令可审计＝**载体在、内容零** |

### 3.3 票 214（附件那台机器接线）——六条

| 条 | 票面要求 | 在场尺 | 读数 |
|---|---|---|---|
| AC#1 | 受理器有生产调用者（改前 0 ⇒ 改后 ≥1，且在装配路径上） | `git --no-pager grep -n 'NewAttachmentBroker' HEAD -- internal cmd \| grep -v '_test.go'` ＋ `git show HEAD:internal/panel/attachments.go \| wc -l` | 剥测试只剩定义（`attachments.go:179`）与返回体（`:189`）＝**生产调用者 0，未改**；件仍 **410 行**（票面那句"410 行、生产零调用者"逐字对得上＝本腿独立复量同值） |
| AC#2 | 超限文件⇒拒，且快照看得到"被拒＋为什么"；正控＝同文件放到限内⇒收 | `git show HEAD:internal/panel/pump.go \| sed -n '269p'` | 逐字 `composer := NewComposerState(mode, workspace, nil, maxAttachment)`＝第三参仍硬写 `nil` ⇒ 快照那四枚键**值恒缺**，今天既没被拒过也没被收过；本格**未起步** |
| AC#3 | 伪装文件按嗅探判（不按扩展名） | `git --no-pager grep -c 'matchesISOBaseMedia' HEAD -- internal/panel/attachments.go` ＋ `git --no-pager grep -c '^func Test' HEAD -- internal/panel/attachments_test.go` | 嗅探在树（4 处命中）＋包内 5 枚测试，但**没有入口喂它**（AC#1＝0）＝判据无从在生产路径上跑 |
| AC#4 | 落盘位置在授权根内，越界那一形必红 | `git --no-pager grep -n 'ArtifactSink' HEAD -- internal cmd \| grep -vc '_test.go'` ＋ `git show HEAD:cmd/wisp/panel_inbound.go \| grep -n 'Attachment:'` | 接口 4 处非测试命中，装配里那格是 `panel_inbound.go:276 Attachment: nil`（注释 `:274` 逐字 `attachment = ticket 92`）＝**sink 没被给过**，越界那一形今天到不了 |
| AC#5 | 来源戳（C25 现成名）盖上，附件内容按不可信外部文本处理 | `git --no-pager grep -n 'MarkWithHostPath' HEAD -- internal/tools/bridge.go \| grep -v '_test.go'` | 盖章桥在（`bridge.go:646`），但附件那条路今天不经过桥＝未覆盖 |
| AC#6 | 三件"没接之前不许说已具备"逐条具名 | `git --no-pager grep -ciE 'attachment' HEAD -- internal/config`（含测试一起数）＋上面三把尺 | 配置层那把尺＝**空返回**，同尺正控 `provider` 在 `internal/config` 命中 11/43/23 枚 ⇒ 那个空是真"零枚附件键"（票面 §09-29 那格的口径）；三条今天全部成立 |

### 3.4 票 215（技能／插件清单与开关）——五条

| 条 | 票面要求 | 在场尺 | 读数 |
|---|---|---|---|
| AC#1 | 清单可枚举且每条有真装配读者（改前零读者⇒改后 ≥1） | `git --no-pager grep -n 'ListPluginStates' HEAD -- internal cmd` | 定义 `internal/memory/dao_misc.go:243` ＋它自己的注释 `:242` ＋唯一调用者 `internal/memory/dao_test.go:451`＝**生产零读者** |
| AC#2 | "配置开着但没生效"必须能被看出来 | `git ls-files 'internal/plugin/*'` | 包内只有 `disposal.go`／`disposal_test.go`／`doc.go` 三枚＝**零"列出已装"面**，谈不上生效与否 |
| AC#3 | 开关跨档位一致（做不到就在清单上写明要重启） | `git show HEAD:internal/panel/config_handlers.go \| grep -n lockedFieldFamilies` | `:108` 逐字 `var lockedFieldFamilies = []string{"risk.", "fs.", "net.", "plugins.", "privacy.", "models.", "audio."}` ⇒ 面板写腿对插件节**整族拒绝**，开关一格零 |
| AC#4 | 外部内容带进来的指令改不动权限 | 沿用票 200 那枚 hostile-file 判据的形状 | ⛔ 本腿不重跑别家票的判据（不在我的射程）；只登记"215 的交付里没有任何写权限的腿"＝零改动可证，功能未建 |
| AC#5 | `unwired.go` 那面旗一枚不少（改前＝改后） | `git show HEAD:internal/config/unwired.go \| wc -l` | **147 行**在树＝旗没被拆；这一格是"不许弄坏"型守卫，⛔ 不构成 215 的产出 |
| 215b（技能那半） | 技能连配置键都不存在，补键＝动 `D36` 契约面 | `git --no-pager grep -c Skill HEAD -- internal/config/schema.go` ／同尺 `-c Plugin` | `Skill`＝**空返回**（尺 rc≠0）、`Plugin`＝8 ⇒ 票面那格"215a 有 8 处配置／215b 零枚键"**今天仍成立**，且 215b 要人工批准，⛔ 本腿不代拍 |

### 3.5 票 216（菜单显示底座：洗控制字符＋"没有"和"读不到"分两句话）——五条

| 条 | 票面要求 | 在场尺 | 读数 |
|---|---|---|---|
| AC#1 | 先做现量普查并出逐枚名册，没名册不许进实现 | 票面那把原尺 `git --no-pager grep -rniE 'sanitiz\|control.?char\|escape.*terminal\|u001b' HEAD -- internal cmd \| wc -l` ＋ `git ls-files 'docs/evidence/s1/216-*'` ＋ `ls -d .scratch/wisp/probes/216` | 尺命中＝**39 行**（非测试 20 行）⇒ 票面"量不到就写没做全仓反向取证，不许写没有"这一格今天可以写成"有零散两处"；但**名册件不存在**（表 0 枚、探针目录不存在）⇒ AC#1 未交，实现也未起步 |
| AC#2 | 一枚纯函数，菜单所有非用户手打字段都过它；落点选"显示前统一洗" | `git --no-pager grep -n 'unicode.IsControl' HEAD -- internal cmd` ＋ `git show HEAD:internal/panel/attachments.go \| sed -n '307,313p'` | 全仓只有**一处**洗控制字符＝`cmd/wisp/approval_reply.go:615`，洗的是**操作员自己打的**理由（正是票 219 更正第 30 行说的"216 现文不覆盖那一格"）；面板那枚 `sanitizeDisplayName`（`attachments.go:310`）函数体逐字只有 `return filepathBaseLike(strings.TrimSpace(name))`＝裁路径段＋去空白，**不洗控制字符**（其上游 `:301` 把含 `\x00\r\n\t` 的名字整枚拒收＝拒绝不是洗）。⇒ **共享底座不存在**，附件名／插件名／技能名／目录名／命令说明五类字段零覆盖 |
| AC#3 | 正控：造一枚文件名带 `ESC`＋"已允许"的附件⇒菜单那一行必须被洗且不改任何判定 | 同 AC#2 那两把尺 ＋ `git --no-pager grep -n 'contains control characters' HEAD -- internal/panel/attachments.go` | 今天的行为是**拒收整枚附件**（`:301` 那枚 `strings.ContainsAny(name, "\x00\r\n\t")` 判定、`:302` 抛 `errors.New("contains control characters")`），不是"收下来但把名字洗过"⇒ 正控那一发的形状与票面要求不同；本格**未交付**（⛔ 本腿不跑 `go test` 去验行为，只按源码字面归口） |
| AC#4 | "没有"与"读不到"分两句话，不许用空数组／缺键同时表达这两件事 | `git show HEAD:internal/panel/instructions_200.go \| grep -nE 'InstructionsStatus(Loaded\|NoneFound\|Refused\|NotRun)'` | 现成五态形状**在树**（`:27 loaded`／`:30 none_found`／`:37 refused`／`:40 not_run`，`:74 Reason` 注释逐字"the loader's own sentence for a non-loaded state"）＝票面 §AC#4 引用的先行版仍在；但它只覆盖项目说明文件那一堆，⛔ 不构成 216 的"每一行清单" |
| AC#5 | 错误提示按错误码分支，不许按英文错误串里有没有某个数字 | `git --no-pager grep -nE 'strings\.Contains\(err\.Error\(\)' HEAD -- internal/panel \| grep -v '_test.go' \| wc -l` ＋ `git --no-pager grep -nE 'observe\.Class\|ErrorClass' HEAD -- internal/panel \| grep -v '_test.go' \| wc -l` | 非测试面板里按错误串取子串＝**0**、按错误码分支＝**0**（命中只在测试里，如 `assets_test.go:30`、`attachments_test.go:323`）⇒ 面板侧既没有子串反例也没有正例，本格**未起步** |

### 3.6 票 219（三枚答复按钮＋理由输入框）——六条

| 条 | 票面要求 | 在场尺 | 读数 |
|---|---|---|---|
| AC#1 | 三枚按钮各有真执行者（生产路径）；会话内同类⇒不再弹卡；进程重启⇒又弹 | `git show HEAD:internal/agent/approval/replies.go \| grep -nE 'func \(r \*Replies\) (Allow\|AllowSession\|Reject\|PanelReject\|PanelAllow\|Veto)'` ＋ `git --no-pager grep -n 'session.NewLedger\|Grants:' HEAD -- cmd/wisp \| grep -v '_test.go'` | 答复侧函数**已在**：`replies.go:316 Allow(ctx, corr)`／`:351 AllowSession`／`:372 Reject(corr, reason)`／`:378 PanelReject(corr, reason)`／`:431 PanelAllow`／`:455 Veto`；装配根里会话身份与授权写读腿在（`cmd/wisp/run.go:478 session.NewLedger`、`:618 Grants: grantWrite`、`:758 Grants: grantRead`）＝票面 09-29 更正③"零执行者"那一格的前半**已由票 224/201 改动**，⚠ 但**三枚按钮的界面真身本腿量不到**（§4 第 5 条）；"重启⇒又弹"那半属运行期 ⇒ 同归量不到 |
| AC#2 | "长期"那一支留下可撤销的一行，审计里查得到谁／何时／加了哪一行／来源是审批卡片，且有"撤销回到加之前"的判据 | `git --no-pager grep -n 'AddAllowedDir\|SetAllowedDirs' HEAD -- cmd/wisp internal \| grep -v '_test.go'` ＋ `git show HEAD:cmd/wisp/approval_always.go \| sed -n '1,20p'` | 写侧在树且有真入口：`cmd/wisp/approval_always.go:134 if err := s.rt.mgr.AddAllowedDir(dir)`（文件头逐字 "Ticket 201 AC#4/AC#5 - the 长期 branch, and the two things it owes"，并按 `PLAN.md:1645-1646` 起第二张 L2 卡）；撤销那一半 `internal/config/allowdirs.go:109 SetAllowedDirs` 非测试调用者＝**0**（只有 `:95` 注释与 `internal/session/grants.go:43` 一句注释）⇒ **"有实现、有测试、生产调用者 0"与票面 10-02 更正同值**＝本格的撤销面仍欠入口（票 226 AC#5 也记着这一格） |
| AC#3 | 网页那侧送进来的 `allow` 必须被判红／被拒 | `git --no-pager grep -n 'ErrPanelAllow' HEAD -- internal/agent/approval/gate.go cmd/wisp/approval_reply.go` | 结构性拒绝在树且**有生产听众**：`gate.go:746 return fmtw(ErrPanelAllow, "面板来源的「允许」被服务端 API 直接拒绝（F2 第三层）…")`＋`cmd/wisp/approval_reply.go:331` 的 `PanelAllow` 那一支会收到它 ⇒ 票面 09-29 更正①（"网页 allow 今天真能批掉 L2 卡"整行作废）今天仍为作废状态；⚠ **本格仍不许翻勾**（票面自陈"这一格 `Q-49` 丙修好之后才算数"，且勾要非实现者裁） |
| AC#4 | 反向正控：种一条"面板侧来源的允许"样本⇒那把尺必须响 | `git --no-pager grep -n 'ErrPanelAllow' HEAD -- internal/agent/approval/queue_test.go` | 现成反控在树＝`queue_test.go:107`／`:110` 两枚用例期望 `approval.ErrPanelAllow`（票面 09-29 裁决⑥(b) 说的"反控现成，不用新造"今天对得上）；本腿**不跑测试**，只按名册在场登记 |
| AC#5 | 理由三去向逐条有读数（审计一行／模型收到的文本带那句理由／用户回看的那份里也在），空理由与超长理由各有确定行为 | `git show HEAD:internal/agent/approval/gate.go \| sed -n '712,730p'` ＋ `git show HEAD:internal/agent/approval/queue.go \| grep -n 'func (q \*Queue) allow'` ＋ `git show HEAD:internal/memory/schema.go \| sed -n '69,85p'` ＋ `git show HEAD:cmd/wisp/approval_reply.go \| grep -n 'replyReasonMax'` | ① 拒绝侧理由已通（`gate.go` 里 `if !r.Allow { return g.q.reject(corr, r.Reason) }`）；② **允许侧无理由位**＝`replies.go:316 Allow(ctx, corr)` 两参、`queue.go:365 func (q *Queue) allow(corr, nonce string)`＝票面更正③那句"真正缺的两样"之一，**仍缺**；③ 用户回看那份仍无列＝`tool_call` 14 列名册（`schema.go:69-83`）里没有理由列＝票面更正⑥那格的对撞未解；④ 空／超长行为有一枚现成上界＝`approval_reply.go:595 const replyReasonMax = 160`（注释写明"不许静默截断"的对应物），但它管的是操作员那一侧，不是面板理由框 |
| AC#6 | L1 与 L2 分表：名册那一格要么一起修掉、要么具名写"这一支没修" | `git ls-files 'internal/panel/*197*' 'internal/panel/*220*'` ＋ `git --no-pager grep -n 'LiveL1Windows' HEAD -- internal cmd \| grep -v '_test.go'` ＋ `git --no-pager grep -n 'L1Windows' HEAD -- cmd/wisp \| grep -vc '_test.go'` | `internal/panel/subagent_roster_197.go` 在树。枚举面本身**今天已建**＝`internal/agent/approval/window_read.go:63 func (g *Gate) LiveL1Windows()`（非测试命中 3 行＝`:49` 注释＋`:63` 定义＋`:24` 注释），面板那侧的插座也建了＝`internal/panel/pump.go:200 L1Windows func() []L1WindowWait`、`:324-325` 会读它；⚠ 但装配根**从没填过那枚插座**＝`git --no-pager grep -n 'L1Windows' HEAD -- cmd/wisp \| grep -vc '_test.go'`＝**0** ⇒ `p.src.L1Windows == nil` 在生产里恒成立，名册那一格今天仍读"没被卡住"。⇒ 这一支归票 220（票面 09-29 更正①已具名"缺一整块"），本腿独立复跑同值；本票只登记，⛔ 不并格充数 |

### 3.7 票 211（子代理池 vs 桥的天花板）——五条（本腿判"已失效"，逐条给产物在场尺，⛔ 不写"无活"）

| 条 | 票面要求 | 在场尺 | 读数 |
|---|---|---|---|
| AC#1 | `MaxConcurrentSubagents` 实测等于 `MaxToolConcurrency`（两枚常量并排打印） | `git --no-pager grep -nE 'MaxToolConcurrency = \|MaxConcurrentSubagents = ' HEAD -- internal/tools internal/agent \| grep -v '_test.go'` | `internal/tools/bridge.go:24 const MaxToolConcurrency = 4`、`internal/agent/budgets.go:47 const MaxToolConcurrency = 4`、`internal/tools/subagent_197.go:85 MaxConcurrentSubagents = 4` ⇒ **两枚相等，甲支已落地**（落地枚＝`7ea14ce3 197-r1b 票 197/211 甲：池的诚实数 8->4（= 桥的 D38d 天花板）＋常驻判据＋两枚桥上路＋gofumpt`） |
| AC#2 | 常驻判据一枚：池**大于**天花板时必须红；正控＝人为把池写成 8 时那枚判据确实红（不许恒真断言） | `git show HEAD:internal/tools/subagent_197_test.go \| sed -n '411,424p'` | `Test197SubagentPoolNeverExceedsBridgeCeiling` 判据体逐字 `if MaxConcurrentSubagents > MaxToolConcurrency { t.Errorf(...) }` ＋第二枚 `if MaxConcurrentSubagents <= 0` ⇒ 形状是"比大小"、不是"相等"恒真断言；另有 `:426 Test197SpawnDescriptionNamesTheRealPoolCap` 钉住"给模型读的那句说明里的数＝常量"（注释逐字"a hand-typed 8 there would keep selling the old cap"） |
| AC#3 | 池满时第 5 枚硬拒＋可读理由（理由带当前在跑的 `taskID`），不占名册行 | `git show HEAD:internal/tools/subagent_197_test.go \| grep -n 'func Test197FullPoolRefusesNextSpawnWithReadableReason'` ＋ `git show HEAD:internal/tools/subagent_197.go \| sed -n '270,276p'` | 用例在（`:534`）；写侧在场＝`subagent_197.go:270 TryAcquireSubagentSlot(MaxConcurrentSubagents)`、`:275` 的拒答文本把 `MaxConcurrentSubagents` 与 `strings.Join(t.d.Roster.RunningSubagentIDs(), ", ")` 一起打 ⇒ 理由里带在跑的那几枚 id |
| AC#4 | 池的测量走真桥（父任务真次派生），不再靠 `spawnDirect` 绕行；若保留绕行测试则须自陈量的是什么 | `git --no-pager grep -n 'spawnDirect' HEAD -- internal/tools \| grep -v '^.*:[0-9]*:\s*//' \| wc -l` ＋ `git show HEAD:internal/tools/subagent_197_test.go \| sed -n '441,452p'` | `spawnDirect` 在 HEAD 只剩**两行注释**（`:443-444`），逐字 "This replaces the pre-211 leg that drove the tool directly (spawnDirect), and **spawnDirect is gone rather than documented**"；替代用例＝`:449 Test197SubagentPoolCapsAtBridgeCeiling`（注释写明 "the cap measured on the production path… through the REAL bridge"） |
| AC#5 | `bridge.go` 与 `budgets.go` 那两枚 **4 一字未动**（本票不借改池碰契约） | `git --no-pager log -1 --name-only --format='%h' 7ea14ce3` ＋ `git --no-pager grep -n 'MaxToolConcurrency = 4' HEAD -- internal/tools/bridge.go internal/agent/budgets.go` | 那枚 commit 的 name-only **只有** `subagent_197.go` 与 `subagent_197_test.go` 两件（没有 `bridge.go`／`budgets.go`）；今树两枚常量仍在 `bridge.go:24`、`budgets.go:47` ⇒ 甲／乙的界线在 diff 里看得见，未被越过 |
| 乙支（真 8＝改契约） | 要 owner 一句话，触碰 `D38`／`D32` | 本腿**不判** | ⛔ 乙支不在"已失效"的射程内：它今天仍是要人工批准那一格，且票面 09-29 更正②⑤已把它改归票 222（`internal/tools/subagent_222_test.go`＋`h222CeilingLiteral = 4` 在树）。⇒ 本腿只判"票面点名的那枚缺陷（名册承认 8、机器只跑 4）今天不在了" |

## 4. 判得心虚的枚数与具名理由

1. **178（判"仍成立"，但读数比票面更难看，我怕归口归错）**：票面记的是"1 枚涨到 2 枚"，我今天跑出来是**基线 1 枚 vs 实测 3 枚**，而且同一次聚合里**还有第二枚 BAD＝`G1b`**（`.Execute(` 那一族，票 178 完全没提）。我心虚的地方：退码从票面的 1 变成 2，其中只有一枚属本票，另一枚属 `G1b`＝如果我下一轮把"聚合退码＝2"整枚记到 178 头上就是错账。⇒ 处置＝上面表里逐枚具名分离，不动仪器、不改基线（那是票 178 AC#4 明令禁的"抬基线"）。
2. **182（判"差翻勾"是我读台账，不是我看内容）**：⛔ `docs/evidence/s1/**` 我只按**文件名**取、一个字内容都没读。我把 182 记成"差翻勾"的凭据是三样：证据表文件名在库＋`.scratch/wisp/probes/182/a1/census.md` 在库（275 行我 `wc` 过）＋台账 `A566` §3 逐字收下 12 堆。心虚点＝**AC#6／AC#7 是契约轴与排程约束那两格**，编排者 §3 只说"收下这张表"，没说这两格；我把整票判"差翻勾"可能偏乐观。⇒ 只登记为"产物在树、框一格没翻"，翻哪几格归编排者逐格裁。
3. **187 与 186（同一枚读面被别人做掉了一半，我判"仍成立"是在判剩余面）**：票 181 落地的 `internal/panel/git.go` 已覆盖 186 的交付 1／2 与 187 的"模型能改"（走 `role_chat_model` 那条写腿）。⚠ 我据此判"仍成立"的**只是没被覆盖的那一半**（思考档位／角色切换的档位语义、AC#4 的"只列声明支持的档"、AC#6 的"换模型＝换窗口"结论）。如果编排者认为"半张票被做掉"应记成"已失效＋新立一枚"，我这格判语就要改——那是我读得最不确定的一处口径。
4. **213／214／215 的"零"是三枚负向结论，我按穷举补过一遍才敢写**：每张都在表里同时给了正控（`Registry` 676／`provider` 在 config 命中 11-43-23／`Plugin` 8 对 `Skill` 空返回）。⚠ 这枚心虚**已就地补掉一半**：本腿后来把 213 那族字面尺的射程从票面默认的 `internal`＋`cmd` 扩到 `tools scripts .github docs/contracts`（＝0）与全树（＝34 行，逐名看全是票面自己＋探针台账件，零枚产码），尺原文与读数在 §5 第 4 条。⛔ 剩下的口子我也照实留着：**中文写法的"命令名册"那一族我没穷举**，所以"任何形状的命令名册都不存在"这句本腿不写，只写"票面那五个字面与一族变体在扩射程后仍零命中"。
5. **界面那半整批量不到（不是心虚，是正确归口）**：186／187／213／214／215／216／219 都含"面板上要画哪几枚控件／按钮点下去"那一格。`frontend/**` 与 `design/**` 本腿**零读零写零转述**（派单 §3），真开窗／真点击也不在只读腿能力内。⇒ 这些格一律归"量不到"，⛔ 我不由"Go 侧字段在"外推成"用户能点"。票判档因此按 Go 侧那一半给。
6. **211／185 判"已失效"我留了一道口子**：两枚的凭据都够硬（commit 号＋name-only＋今树在场尺）。心虚点＝**票框一枚都没勾**，而票 211 的乙支（把 4 抬到 8＝改 `D38`／`D32` 契约）今天仍未被批准，票 185 票面还写着"已立、按住不派"。⇒ 我判的是"票面点名的缺陷今天不在了"，⛔ 不是"这张票可以结案"；结案要编排者翻勾并处理乙支那格。
7. **219 的 AC#1 我判得最薄**：三枚按钮在 Go 侧有对应答复函数（`Allow`／`AllowSession`／`Reject`），但"按钮"本身在原生窗口／托盘／球旁卡上——我一把尺都够不着，且票面裁决段（09-29 ①③④）说的派单前置序（201→223→224→本票）里，223／224 两枚票**还没 `-done`**（`ls .scratch/wisp/issues/22[34]-*.md` 无 `-done` 后缀）。⇒ 本格我只写"答复侧函数在场"，不写"这一格完成"。

## 5. 尺有歧义／判不动的地方

1. **同一把尺两形读数不同（起手就撞上，两把原文都留）**：对同一个文件 `.scratch/wisp/issues/178-pair-census-ruler-g6-negative-leg-cannot-express-closing-a-scope-opened-elsewhere-and-rings-every-honest-test.md`——
   - `grep -c '^[[:space:]]*- [ ]' "$F"`（方括号**不**转义）⇒ **0**
   - `grep -c '^[[:space:]]*- \[ \]' "$F"`（转义）⇒ **5**
   - `grep -c '\- \[ \]' "$F"`（派单 §1 那把尺的形状）⇒ **5**
   两把本腿**不挑一把当准**。后续框态一律用转义＋容缩进形（它同时容缩进与 `*` 子弹）。⚠ 派单 §1 那句"差的 7 枚＝196／211／213／214／215／216／219"用转义形复量后**名单不变**；而我起手那把没转义的尺会把**有框票**误报成无框票——对未来重跑这堆尺的腿是防雷。
2. **票 232 的"顶格尺 6／全形尺 7"是两把不同维度的尺，不是矛盾**：`grep -c '^[[:space:]]*[-*] \[ \]'`＝6、`grep -o '\- \[ \]' \| wc -l`＝7（本腿逐枚比对：没有任何一行含两枚框、没有缩进行、没有 `*` 子弹 ⇒ 第 7 枚出现在**行锚点之外**的引用文本里）。⇒ 枚数一律报"框判据 6 枚＋尺 `grep -o` 7 枚"两把，⛔ 不挑一把。这一枚本腿只是扣除，不判档。
3. **三枚票面前提在本腿复量下已过期（写下来给编排者，不改票一个字）**：
   - 票 194 AC#5／票 186 普查注／票 219 表第 4 行都写着"`ComposerDispatch.Handle` 生产调用者零枚／入向那一跳整条不存在"。今树尺：`git show HEAD:cmd/wisp/panel_host_windows.go \| grep -n 'bindErr := w.Bind(panelDispatchBinding'`＝**`:405` 命中**，`git show HEAD:cmd/wisp/resident_windows.go \| grep -n 'first non-test call of NewPanelManager'`＝`:141` 注释逐字承认"直到此刻才第一次有非测试调用者"。⇒ 这句前提**按今天 HEAD 不成立**；本腿因此没有用它去否掉 213／214 的任何一格。
   - 票 215 票面"配置里有 8 处插件键"＝本腿同尺复量**仍 8**（未过期，记在这里是为了说明两锚之间这枚没漂）。
   - 票 196 票面现量表引的是锚 `8c671afb` 的行号（`internal/agent/loop.go:983-998`）；今树同内容在 `:994/998/1004`＝**行号漂 11 行**，内容未漂。⇒ 引用一律现跑，不抄票面行号。
4. **负向结论的维度穷举（我把该扫的补扫了一遍，剩下的口子具名留着）**：213 那五把零命中尺的票面射程＝`internal`＋`cmd`；本腿额外补跑两把——
   - `git --no-pager grep -niE 'slashcommand\|commandcatalog\|CommandRegistry\|/compact\|/goal' HEAD -- tools scripts .github docs/contracts \| wc -l`＝**0**；
   - 全树（含 `.scratch/**`）`git --no-pager grep -rniE 'slashcommand\|commandcatalog' HEAD -- . \| grep -v '_test.go' \| wc -l`＝**34 行**，逐名看命中文件＝票面自己（`.scratch/wisp/issues/213-…md`）＋`.scratch/wisp/probes/212/a3/work/*.tsv`＋`.scratch/wisp/probes/213/c1/inventory.md`＝**全是台账与探针件，没有一枚产码**。
   ⇒ "全仓零枚斜杠命令目录"这句现在有我自己的两把尺撑着（票面射程一把＋扩射程一把），⛔ 但我仍不写"任何形状的命令名册都不存在"——我只按这五个字面与一族变体搜过，中文写法的"命令名册"我没穷举。
   另一处已补：214 的"配置层零枚附件键"那把尺，本腿把同一形状打到 `cmd/wisp` 上复算了代码默认值那一族——`AcceptedMIMETypes()` 真身在 `internal/panel/attachments.go:153`，读者是 `internal/panel/composer.go:120`（`composer.AttachmentMIMEs = AcceptedMIMETypes()`）与 `:281` ⇒ 票面 §09-29 那格"上限与白名单来自代码里的默认、没有用户可调的入口"**逐字对得上**，`internal/config` 那把尺的空返回是真负向。
5. **`gate-clauses.sh` 那把尺我今天跑了第二遍才敢写枚数，并核了两枚锚点之间没漂**：第一遍在后台跑（超时转后台），聚合读数＝`G6neg` 实测 3 枚、`G1b` 1 枚、退码 2。第二遍本腿改用手工按文件算（`git --no-pager grep -lw` 两集合求差＋逐行看注释）得反向候选 5 枚文件、剔注释后 3 枚 ⇒ **两法同值 3 枚**，没有读数冲突。⚠ 写件期间 HEAD 从 `b7251d17` 前进到 `7329f07f`（编排者与 236-v1／evidence-close-6 那些枚），`internal`＋`cmd` 范围内只有一枚文件漂＝`internal/tools/tasklist_deferred_236r3_teeth_test.go`（`+31 −8`）⇒ 本腿逐把尺把两枚锚各打一遍：`OpenTask`／`CloseTask`／`OpenScope`／`CloseScope`／`DisposalScope`／`ToolRequest{`／`.Execute(` 在两锚上**都读 0**，`Defer` 在两锚上**都读 3** ⇒ 那枚漂动对这把仪器的每一腿读数零影响，本腿 §2 表里所有 `HEAD` 读数在 `b7251d17` 与 `7329f07f` 上同值。
   下一位若在别的锚上重跑，请注意：`gate-clauses.sh` 默认锚是 `$(git rev-parse HEAD)`＝**不读工作树**，而 §0 那行 porcelain 里的在飞脏件（`internal/tools/tasklist_deferred_236r3_teeth_test.go` 是**跟踪件的工作树改动**）今天因此**没进**我的读数——很容易误读成"我已经把在飞的活量进去了"。
6. **判不动的一枚：182 的堆数**。票面 AC#1 要求"枚数只许现量不许目测"，而台账里有**三把不同尺给出三个数**：`A392` 记 7（原型注册表口径）、`A566` §3 记 12（按堆口径）、`docs/evidence/s1/182-rail-stacks-census-a1.md` 那枚 09-28 版又是 7。⇒ 本腿**不选一个数当准**（那要读表内容才裁得了，而我只许按文件名取表）；只登记"AC#1 的枚数口径在台账里就不止一版"，判档因此停在"差翻勾"而不是"差翻勾 AC#1"。
7. **★记我：本腿写过一枚假负向读数，落笔后复跑顶正（原句已在 §3.2 那一格改回两把读数，不抹）**。§3.2 票 213 AC#3 那格我第一遍写成"**0**＝无命令入口形状"，复跑同一把尺得 **67 行**（剥测试 4 行，4 行全是 `filepath.ToSlash`／`provider/model` 那族斜杠，不是命令形状）。⇒ 两处教训写在这里给下一位：① **子串形尺（`slash`）与字面形尺（`/compact`、`slashcommand`）不是一回事**，前者命中一堆无关词，别拿它的"看起来该是 0"当读数；② 我在同一格本来只需要"命令入口形状零"这一句，而那句的正名尺是 §2 表里那五把字面尺（逐枚 0）＋`git ls-files` 载体尺（0 枚），不需要 `slash` 这把。⛔ 不回滚、就地在 §3.2 那一格并列两把读数。
8. **本腿自己第二枚 commit 的删除列不是 0，逐行具名（守"零删除"的规矩要写清删的是什么）**：`git --no-pager diff --numstat -- .scratch/wisp/probes/pool-validity/4f/batch4f.md`＝`132 5`，被删的 5 行全部是**我第一枚 commit（`de084db7`）里自己写的占位符**，没有一行是别人的读数，且逐行都留有去处：
   - `1. **同一把尺两形读数不同（本腿起手就撞上）**：…` ⇒ 同句在 §5 第 1 条以"（起手就撞上，两把原文都留）"重写，内容一字不减；
   - `   - \`grep -c '\- \[ \]' "$F"\`（派单 §1 那把尺的形状，去锚点＋转义）⇒ **5**` ⇒ 现写作"（派单 §1 那把尺的形状）⇒ **5**"（同一把尺同一读数，只短了那六个修饰字）；
   - `   两把本腿**不挑一把当准**，一律记原文；后续所有框态用转义形，理由是它同时容缩进与 \`*\` 子弹。` ⇒ 现写作"两把本腿不挑一把当准。后续框态一律用转义＋容缩进形（它同时容缩进与 `*` 子弹）"；
   - `   ⚠ 这直接说明：派单里"未勾＝0 的七枚"这个**名册结论**不依赖我起手那把没转义的尺——…是防雷。` ⇒ 现写作 §5 第 1 条末两句（同一件事，措辞收紧，"七枚名单不变"与"会把有框票误报成无框票"两个结论都在）；
   - `2. （待补）` ⇒ 被 §5 第 2–8 条填满。
   ⇒ 下一位若拿"删除列必须为 0"来挑我这枚 commit，答案就是本节：**占位符被实内容替换，零信息损失，逐行可回原名**（`git show de084db7:.scratch/wisp/probes/pool-validity/4f/batch4f.md` 里五行原文一字不改还在）。⛔ 本腿不因此回滚、不 `restore`。
9. **★同一形制再记一次（第三枚 commit 的删除列 8 行，逐行具名替换来源）**：`git --no-pager diff --cached --numstat` 对本件读 `142 8`，其中 5 行＝上面第 8 条那五枚占位符，另 **3 行是 §1.3 扣除清单里 233／232／236 三枚的旧短行**，被本腿换成带时间戳与尺原文的长行；旧三行的原文一字不改还在 `de084db7` 那枚 commit 里（尺＝`git show de084db7:.scratch/wisp/probes/pool-validity/4f/batch4f.md`），三枚扣除结论**没有因改写而变**（233／234／237＝已由 `4e` 判掉，232／236＝在飞，249／269＝票面已就地标撤回）。⛔ 不回滚、不 `restore`，要更正只许再追加。
