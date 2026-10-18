# 340-a1 普查件 — 票 340 只读普查腿（AC#0 起手现量 ＋ AC#1 三形代价表）2026-10-18 第二波

腿＝`340-a1`（只读普查，AC#6 明文零产码）；本件＝交付主体，⛔ 判语（裁腿≠本腿）。
入场锚＝**`feac8f77`**（本腿 11:3x `git log -1` 现认；编排者 11:2x 所见 `130b1475` 系是其祖先，已现量）。
★全尺取**具名提交树**（`git show <rev>:<path>`／`git grep <pat> <rev> -- <path>`），⛔ 读在飞字节：本波 `336-r1`(internal/risk) ∥ `227-r1`(internal/config) 在飞。
三锚同钉对照＝`feac8f77`（入场）⇄ `f705b082`（票面记锚）⇄ `5db0b538`（裁腿锚）。
可逐字重跑＝`scripts/ac0.sh`／`scripts/ac1.sh`／`scripts/ac1b.sh`；逐字输出＝`logs/ac0.log`／`logs/ac1.log`／`logs/ac1b.log`。
本件引用的一切颜色一律＝**读码推的**（本腿⛔ 跑 `go test`／`go build`）。

## 一、AC#0 起手现量（票面四行逐枚重跑）

| 票面那行 | 本腿读数（具名锚） | 尺 |
|---|---|---|
| 真产叠改形用例＝零 | **复认**：窄尺在 feac8f77／f705b082／5db0b538 三锚各 **0 行** | `git grep -nE '"[^"]*#[0-9]+"' <rev> -- internal/agent/approval` |
| 第三臂侧那一形实测点亮 | **点名复认**（⛔ 本腿跑）＝夹具形，落点 `internal/panel/subagent_corr_task_prefix_328_test.go:41`；颜色凭据＝裁腿判语件第一节 1 条 `docs/evidence/s1/328-v1-verdict.md:11`（`after-agent-positive: RUN=2 PASS_TOP=2 FAIL_TOP=0 SKIP=0 rc=0`） | `git grep -n 'arm B with the queue clash rewrite stacked' <rev> -- internal/panel` |
| ★⛔ 读成"已被钉住" | **复认并给出名册**（见下面"臂分账"） | 两枚臂本体逐字 |
| 链起来要新导出名 | **★更正这一行**（见下面具名更正） | 交集硬尺 `scripts/ac1b.sh` 第 1 节 |

臂分账（数字臂 vs 第三臂，各自承重；量的＝字节判据，读码推的＝颜色）：
- 数字臂本体＝`internal/panel/subagent_roster_197.go:333-344`（`splitClashSuffix`：`LastIndexByte(id,'#')` ＋尾段逐字符验 `'0'..'9'`），唯一读者 `:308`；
- 第三臂本体＝`:312` `if strings.HasPrefix(key, task+"#")`，名册唯一读点＝`:252` `waiting[row.TaskID]`；
- 摘掉第三臂会红（读码推的＝对每枚字面串跑两枚臂的字节判据；粒度＝表项／直投形／断言点，⛔ 用例名）＝**12 枚**，逐枚：`subagent_corr_task_prefix_328_test.go:39/:40/:41/:43`（表项 `child-1#call_e1`／`#call-0`／`#call-0#7`／`#call_e1#7`）＋同件 `:65/:66/:73/:83`（直投形 `root-1#call-0`／`child-1#call_e1`／`root-1#call-0`／L1 窗口 `child-1#call-0`）＋`subagent_blocked_220_test.go:141`（表项 `child-1#retry`）＋`internal/agent/corr_roster_join_328_test.go:120-124`（真派发的那枚 corr）＋同件 `:161-166` 的两形（armA／armB，一枚循环里的两枚迭代）；
- 摘掉第三臂**不**红＝`328_test.go:42`＋`220_test.go:138/:139`（`child-1#9`／`child-1#147`，两枚臂都认）＋全部**期望 false** 的表项 `328_test.go:44/:45/:46/:87`、`220_test.go:140/:142/:143`（摘臂后仍不点亮，照旧绿）＝票面第三行"数字臂在承重"那一支成立；
- ★本腿独立给的一味（字节包含论证，读码推的）：`splitClashSuffix` 的 base 恒等于"最后一个 `#` 前的整段"（`:334` `LastIndexByte` ⇄ `:343` `id[:i]`）⇒ base 若是一枚在册任务号，该 key 必已以 `base+"#"` 开头 ⇒ **第三臂同时命中**。故**摘掉数字臂单独不红任何一枚**（名册唯一读点＝`:252` `waiting[row.TaskID]`），数字臂今天对"点亮任务行"这一维被第三臂**包含**；差别只在 `set[base]` 里多一枚非在册键（`other-1#9` 的 `other-1`），那一形无人查。⇒ 票面第三行"数字臂在承重"这句只在"**摘掉第三臂后它们仍绿**"这一支为真，⛔ 读成"数字臂单独不可替代"，更⛔ 读成"叠改形已被钉住"。

★具名更正（票面 AC#0 明文要的就是这个，⛔ 我改票面原文）：
票面第四行"链起来要新导出名＝跨包两枚私有件"，**在盘上过窄**。交集硬尺响 **1 枚**＝`internal/tools/loop_approval_test.go`（`package tools_test`）今天已同时装配真 loop（`agent.New` `:142`）＋真桥（`tools.New` `:103`）＋真 Gate（`approval.New` `:94`），而这一跳的四道门逐枚全是**导出名**：`Gate.PendingApproval` `gate.go:502`（→ `g.q.push` `:507`）／`Gate.Queue` `:169`／`Queue.LiveApprovals` `pending_read.go:113`＋读出口 `:125`／`panel.TaskRosterSectionFrom` `:196`。尺＝支 L 27 枚∩支 P 12 枚＝1 枚（`scripts/ac1b.sh` 第 1 节，正控＝支 P 名册逐名 12 枚）。⇒ "要新导出名"是**甲／乙两形各自**的性质，⛔ 是"链起来"这一整支的唯一形状。

两枚铸造点的非导出性现量（票面 d 行）：`callCorr` def＝`internal/agent/loop.go:603`（非导出），非测试调用点＝**1 枚**＝`:676`；`Queue.push` def＝`internal/agent/approval/queue.go:156`（非导出），产码调用者＝**1 枚**＝`gate.go:507`（含测试共 7 枚命中）。

## 二、AC#1 三形代价表（只量、⛔ 选形、⛔ 推荐；AC#2 硬门今天没开）

| 形 | 动几枚文件∶几行 | 要不要新导出名 | 禁名词表那族咬不咬 | 牙／反形（⛔ 本腿跑，只写形状） |
|---|---|---|---|---|
| **甲**＝approval 侧测试注入口（`export_test.go` 形或包内 `_test.go`） | 2 枚：注入件 1（先例＝`internal/winsec/export_test.go` 50 行、`grep -n "^func "` :17/:31/:48 三枚导出名）＋端到端测试件 1（基线＝328 那两枚 171／89 行）＝**推算的界** 60–220 行 | **真**（export_test 形＝新导出标识符；★盘上现量：approval 22 枚测试＝**10 枚 `package approval`** ⇄ 12 枚 `package approval_test`，且 `pending_read_test.go:47`／`ticket242_binding_test.go:80/:99/:132/:139`／`ticket259_denial_rulers_test.go:85` 三枚**全属那 10 枚 in-package**、已直呼非导出 push＝队列半侧零新导出名，但要拿到**真 callCorr 的字节**仍需跨包入口 ⇒ 甲单独不成"真产"） | **假**（量的：该族只 `filepath.Join(root,"internal","panel")`＋`goSourceFiles` 排除 `_test.go`（`:369` 逐字）⇒ approval 侧新增件射程外；两枚词表逐字＝`grantFieldWords` 19 枚 ⇄ `grantRouteWords` 11 枚词根，本形候选名×词根的子串比对＝**7 枚候选名全 0 命中**（算的，逐枚在 `logs/ac1.log` 第 5 节）） | 触 **AC#2**＝须台账 `A##` 具名批准（记文件／行／理由／边界／撤销口令），批准今天没开 |
| **乙**＝agent 侧 `callCorr` 包外测试别名 | 1–2 枚：别名件 1（＋1 枚导出标识符）＋测试件 1＝**推算的界** 20–180 行 | **真**（定义即"包外可见别名"＝新导出名；★盘上现量：`callCorr` 已被 in-package 测试直呼（`corr_roster_join_328_test.go:116/:152/:157`、`ticket285_corr_distinct_rulers_test.go:51-60`）＝铸造半侧零新导出名；但 agent 侧 **24 枚测试逐枚全 `package agent`、零枚外部包**（尺＝逐枚 `grep -m1 "^package "` 计数），且 in-package 测试 import approval 会成环（approval→tools→agent，产码 import 名册逐枚）⇒ 乙要的是"外部包拿得到"那一跳） | **假**（同上射程尺：新增件落 `internal/agent`，⛔ internal/panel） | 触 **AC#2**＝同上 |
| **丙**＝不链，钉住 `queue.go` 那枚改写格式串的内容锚（一变即红） | 1 枚：approval 侧 in-package 测试读自己包源文件（先例＝`internal/tools/task_cancel_221_legs_test.go:224` `os.ReadFile("task.go")`，`package tools`）＝**推算的界** 20–60 行 | **假**（零新导出名＝派单那句"默认先量丙"的凭据） | **假**（同一把射程尺；且本形不新增任何可解码结构体 ⇒ 那族"第二枚 decode 目的地"约束零触发） | 反形＝把 `:168` 的 `%s#%d` 改成任何别的字节（`%s_%d`／`%d#%s`／换成拼接）⇒ 新钉"锚仍存在"断言红；成对第二发落点候选＝`:616` 或 `roster:312`（⛔ 同一枚断言行两处落点） |

★丙形今天的"没牙"具名（量的，＝本票那一格缺的东西）：`git grep -nE '"[A-Za-z0-9_.-]+#[A-Za-z0-9_.-]*"' feac8f77 -- internal/agent/approval` 在 `*_test.go` 档＝**0 命中**（正控＝同尺族 `internal/panel` 测试档 **20 行**）＋approval 测试读改写输出全走符号化（`.Corr` 39 枚，例 `queue_test.go:283/:290` 只比"不等"、不比字面串）⇒ **改写串今天换成任何别的形状，全仓无一枚既存用例红**＝票面"缺世界"那一形复认。

## 三、既存先例的射程边界（★⛔ 读成"端到端已被覆盖"）

`internal/tools/loop_approval_test.go` 这枚交集件：装配的是真 loop＋真桥＋真 Gate，但它今天 `risk.L2` 与 `PendingApproval` 命中各 **0**（尺＝`git grep -n "risk.L2\|Level: risk.L2" <rev> -- internal/tools/loop_approval_test.go` ⇒ 0 命中 rc=1；编排者独立复跑同数）⇒ 它走 **L1**（`bridge.go:451` `PendingWindow`），而 L1 那一路**不经 `q.push`**（`gate.go:252` ⇄ `:502` 边界逐字）⇒ 撞名改写那一支今天没被它踩到。
另两处边界：① 该件在 `internal/tools/**`＝**票面射程外**（票面射程＝agent／approval／panel 测试面）⇒ "端到端那一发落哪一包"票面没写＝归编排者裁（本腿⛔ 自行挪射程）；② `fs.write` 的声明档＝L1、`fs.edit`＝L2（`fs_test.go:195`）＝这一形今天可由**夹具声明**造出，⛔ 是"真产已经过"。

## 四、必答 (i)–(iv)（逐字）

(i) "真产枚数 0"那把尺的正控读数＝**同一枚字节序列只换射程**：`git grep -nE '"[^"]*#[0-9]+"' feac8f77 -- internal/panel`＝**13 行**（逐枚在 `logs/ac0.log`，含 `subagent_corr_task_prefix_328_test.go:41/:42/:43`、`subagent_blocked_220_test.go:138/:139/:140` 等）；再加一档射程 `-- internal cmd tools scripts`＝**19 行**。原尺原射程（approval）＝三锚各 0 行。★并具名作废一枚假正控：把字符类放宽成 `#[A-Z]` 打到 approval＝**0 命中**（rc=1）⇒ 那一形⛔ 当正控（该射程内真无 "#字母" 字面串）。
(ii) 格式串今天逐字＝**`internal/agent/approval/queue.go:168`　`corr = fmt.Sprintf("%s#%d", d.CorrelationID, q.seq)`**；内容锚 token＝**`%s#%d`**（行号只当读数）；同形第二枚＝`:616`　`d.CorrelationID = fmt.Sprintf("%s#%d", corr, q.seq+uint64(i)+1)`（replay 支，操作数不同）；queue.go 整枚 blob＝`66fec7ae375344d7bc741270b42c4ec5050cf41f`，三锚（feac8f77／f705b082／5db0b538）**同一串**＝那枚字节三锚未变。
(iii) 三形"要不要新导出名"逐形＝**甲：真**；**乙：真**；**丙：假**。（甲／乙两枚"真"都触 AC#2 硬门＝未经台账 `A##` 具名批准不许开工；两枚各自的队列半侧／铸造半侧盘上已零新导出名，缺口都在**跨包那一跳**，见表内具名行号。本腿⛔ 选形、⛔ 写"建议选 X"。）
(iv) 丙形覆盖情况分两维具名：
- **"换分隔符"（`#`→`_`）与"seq 放前面"（`%s#%d`→`%d#%s`）两形＝覆盖得到**。尺＝内容锚 `%s#%d` 的存在性断言（`queue.go:168` 体内那枚 token，⛔ 绝对行号）⇒ 字节一变即红（读码推的；本腿⛔ 跑）。
- **覆盖不到那一维＝叠改的触发条件与消费端接线**，具名三处字节＝① `queue.go:167` `if _, clash := q.byID[corr]; clash`（判定操作数＝`byID`，`:301 delete(q.byID, it.Corr)` 摘除后不撞＝叠改不发生，串照原样、丙形不红）；② `:168` 拼的是入参 `d.CorrelationID` 而⛔ 局部变量 `corr`（`:163-165` 的空 corr 回退支被绕开）＝这一维换形丙形同样不红；③ 名册侧消费点 `roster:252`＋两枚臂 `:308/:312`＝丙形对"改写后的字节真被读到"零读数。⇒ 丙形钉"串长这样"，⛔ 钉"这条分支真被走到、且它的输出真被下游消费"。

## 五、AC#3／AC#4 前置料（⛔ 提改法）

端到端那一发若落地，判据该断的三件（票面 AC#3 已写，本腿补每件的现读尺）：
1. corr 逐字节含两枚 `#`＝尺形＝对上段第一枚字面串做 `strings.Count(corr,"#")==2` 那类字节判据（⛔ 符号化"与旧不等"，＝丙形没牙那一维的来源）；
2. 剥出的段仍是生产 corr（⛔ 裸任务号）＝凭据形＝`loop.go:603/:607` 那枚 `taskID+"#"+callID` 的导出搬运（`ToolRequest.CorrelationID`，`loop.go:676`→`bridge.go:383`），裁腿已定界"corr 半边生产真形／在册半边夹具代劳"（判语件第一节 4 条）；
3. 名册点亮**原任务行**＝`roster:252`＋`:312` 那一臂，夹具基线＝`328_test.go:41`（今天只有它，⛔ 真产）。
去牙落点候选（具名，成对两发必须不同枚）：候选甲＝`queue.go:168`（push 支那枚串）；候选乙＝`queue.go:616`（replay 支同形串）；候选丙＝`roster:312`（第三臂 `HasPrefix` 带边界那一支）；候选丁＝`queue.go:167`（撞名判定本身＝丙形看不见的那一味）。★⛔ 同一枚断言行两处落点。
AC#4（若裁"不链"）要的形状＝内容锚钉＋"突变必红"，且本腿已具名它缺的第三件（上面 (iv) 那一维）⇒ 单靠丙形⛔ 能写成"叠改形已被钉住"。

## 六、AC#5 前置料（点名，⛔ 跑）

本票若落地，成对台面落到三包＝`./internal/agent`／`./internal/agent/approval`／`./internal/panel`；分母现量（`git ls-tree`，具名锚 feac8f77）＝agent 24 枚测试／13 枚产码、approval 22 枚测试／10 枚产码、panel 29 枚测试／17 枚产码。★`internal/panel/**` 每波只许一枚写腿＋`cmd/wisp` 本波零写腿＝同一台面成对时的具名约束；`./cmd/wisp` 整包⛔ 计入（写面互斥时具名跳过）。两把临时根台面⇄`-count=1 -timeout 30m`⇄四数⛔ 混用⛔ 相减＝派单现行形状，本腿零 go。

## 七、未定义即停两处（已上报，编排者 12:3x 答，逐字记此）

1. "端到端那一发落哪一包"票面没写 ⇒ **编排者裁**（本腿只读上报，⛔ 挪射程）。
2. 那枚交集件 `risk.L2`／`PendingApproval` 命中各 0 ⇒ 它走 L1、不经 push ⇒ 编排者令本件写成"既存先例**射程边界**"，⛔ 写"端到端已被覆盖"（已照此落第三节）。

## 八、纪律自证（本腿写面）

零 `.go` 字节（尺＝逐枚 `git show --name-only` 并集，`.go` 命中 **0**）；写面只两处＝`.scratch/wisp/probes/340/a1/**` ＋票 340 尾**一条**只追加 bullet；⛔ 台账／停车点／别的票／票 328 判据文字（它已 `-done`）；⛔ 新建票件／翻框／改票名／`Status:`／`push`；`add` 与 `commit` 两处 pathspec 逐枚字面路径（⛔ 从共享 index／`git status` 反推）；中文正文一律 Write 成文件＋`git commit -F`（⛔ PowerShell 内联，台账 `A1030`）；件内零分节符、CR＝0（尺＝`tr -dc`）；数 AC 框只用判据区射程 `'^- \[ \] \*\*AC#'`＝**7**、已勾 **0**（⛔ 全档形，`A1031`）。
