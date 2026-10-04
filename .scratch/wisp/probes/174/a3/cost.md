# 174-a3 只读普查 — `ErrReparseDenied` 进模型可见文本的三形代价表

- 起手锚：时间 `2026-10-04T09:09:34+08:00`（`date -Iseconds`）· 起手 HEAD `7ac8965a`。
- 起手 `git status --porcelain internal/tools cmd/wisp tools/d22scan`＝两行（`M tools/d22scan/main.go`／`M tools/d22scan/selftestsamples.go`）；**本程中途工作树又长了 5 枚 M（交件时终态复量，多出的那枚是 `resident_hotkey_258_test.go`）**（`cmd\wisp\resident_ball_windows.go`／`cmd\wisp\resident_hotkey_258_windows_test.go`／`cmd\wisp\resident_windows.go`／`internal\tools\bridge.go`，＝258-r2 那一类写腿在飞），HEAD 也从 `7ac8965a` 走到 `a73edf62`。⇒ **本文所有 `file:line` 一律按"HEAD 已提交面"复量**（尺＝`git show HEAD:<路径> | grep -n …`），不引工作树行号，免得别人未提交的位移把锚漂掉。
- 骨架提交 `98d022df`（先于填表）。本程只读：产码／`docs/**`／票面 AC 框**零字节**，只写这一枚文件。
- 口径：以下每一跳都是**读码**读出来的，一条都没跑（派单禁 Go 命令）。

## §1 现量（复认，非照抄台账）

| 项 | 读数 |
|---|---|
| 定义处 | `D:\work\workspace\projects plans\Wisp\internal\risk\pathresolver.go:36`（doc 注释在 `:34-35`） |
| 错误串逐字 | `risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions` |
| 模板里的变量位 | **零枚**。它是 `errors.New` 的静态串，没有 `%s`、不插路径、不插用户填的配置值；串里唯一的"用户侧词汇"是**配置项的名字本体** `reparse_point_exceptions`（这词本来就写在 `docs\specs\SPEC-06-security-gatekeeping.md:56` 与 `internal\config\schema.go:479`） |
| 同支另一枚可达错误 | `internal\tools\paths.go:121` `tools: empty path`（同样静态、零变量位；生产不可达——`ArtifactPath` 由 `:137` 的 `case sp.Path != ""` 保证非空） |
| 可达错误集推导式 | `grep -n 'errors.New\|fmt.Errorf' internal/tools/paths.go internal/risk/pathresolver.go` → 4 行；其中 `pathresolver.go:49`/`:95` 属 `ErrRewrittenPath`，只在 `Result.Actable()` 里造，而 `grep -c 'Actable' internal/tools/paths.go` = **0** ⇒ 它今天流不到那一跳 |
| 拼进模型可见文本的那一跳 | `internal\tools\task.go:837`（`pointerNotice` 函数体 `:829-851`，`case err != nil` 支） |

**"那一跳"逐跳（谁产／谁投／谁落，全路径）**

1. 产：`internal\risk\pathresolver.go:123`（`return res, ErrReparseDenied`）→ `internal\tools\paths.go:104-107`（`resolve` 原样透出）→ `internal\tools\paths.go:123-126`（`Canonicalize` 返 err）。
2. 拼：`internal\tools\task.go:834-837`（`"注意：…连规范化都没通过（"+err.Error()+"）"`）→ `internal\tools\task.go:548` → `internal\tools\task.go:561-563`（`fmt.Sprintf` 的 `%s` 槽）→ `internal\tools\task.go:573` `Result{Text: stub, Truncated: true}`。
3. 投：`internal\tools\bridge.go:535`（`entry.Tool.Execute`）→ `internal\tools\bridge.go:571`（`Text: res.Text`）→ `internal\agent\loop.go:741`（`l.opt.Tools.Execute`）→ `internal\agent\loop.go:660`（`results[i], execErr[i] = …`）→ `internal\agent\loop.go:697-713`（`out := results[i]` → `log.Text = sp.Text`）→ `internal\agent\loop.go:722`（`l.append(toolResultMessage(c.ID, log.Text, …))`）→ `internal\agent\loop.go:573`（`l.asm.Build(in, l.History())`）→ C5 provider，**出网**。
4. 落盘：**不落 SQLite**。`internal\memory\schema.go:69-86` 的 `tool_call` 表没有正文列（列到 `outcome`/`error_class`/`grant_id` 为止）；`internal\memory\schema.go:58` 的 `task_log.query_text` 存的是用户那句 query。进程内另有两枚落点：artifacts 文件（存任务原文，不存回执）与 C25 片段索引（`internal\tools\bridge.go:631` `mark()` → `internal\risk\provenance.go:558`）。
5. 人在哪儿看不到：`internal\agent\loop.go:723-725` 发的 `EvToolEnd` 虽带 `Text`，但唯一消费点 `cmd\wisp\run.go:1278-1280` 只印工具名与结果词 ⇒ 终端不印。推导式：`grep -rn 'EvToolEnd' --include='*.go' internal/ cmd/ | grep -v '_test.go'` → 6 行＝3 枚产出＋2 枚注释/常量定义＋1 枚消费。

**生产调用者枚数（并写明"这个数答不了能力问题"）**

- `pointerNotice`（task 腿）显式调用者＝**1**（`internal\tools\task.go:548`；`grep -n 'pointerNotice' internal/tools/task.go` → 5 行＝3 行注释＋1 行函数声明＋1 行调用）。
- `task.output` 注册点＝**1**（`cmd\wisp\run.go:544`，且 145-r3 已把判定者接上：`Paths: rt.paths`，commit `4db5f3f6` 已核实在盘）。`TaskDeps` 组装点全仓也只这 1 枚（`grep -rn 'TaskDeps{' --include='*.go' internal/ cmd/ | grep -v _test` → 1 行）。
- ⇒ 但真正决定"看不看得见"的不是调用者枚数，而是**谁能往模型读得到的那一行里放一枚非空 `ArtifactPath`**：全仓写入者 **1 枚**（`internal\tools\task_backfill.go:138`；推导式 `grep -rn 'ArtifactPath *=\|ArtifactPath:' --include='*.go' internal/ cmd/ | grep -v '_test.go'` → 恰 1 行），它唯一的生产调用点 `cmd\wisp\run.go:1107` 位于 `res := bg.Wait()`（`:1106`）**之后**。见 §2。

## §2 今天有没有真在漏

结论：**那一跳今天在生产里接不上**。三个必要条件，第一个就不成立。

| 条件 | 读数 |
|---|---|
| ① 名册里有一行带着非空 `ArtifactPath`、且模型在同进程内调 `task.output` 读得到 | **不成立**。唯一写入者 `internal\tools\task_backfill.go:138` 只被 `cmd\wisp\run.go:1107` 调，而那是在 `bg.Wait()` 返回（环路已结束）之后；`wisp run` 随后退出，常驻腿 `cmd\wisp\resident_task_source_windows.go` 名下零 `Backfill` 调用。子代理行走 `internal\tools\subagent_197.go:424`（`Record(childID, TaskOutput{Text…, State…})`，**不带 `ArtifactPath`**）⇒ 它落到 `internal\tools\task.go:564-568` 的"无副本文件"那一支，`pointerNotice` 根本不执行 |
| ② 就算①成立，那枚路径要有任一**存在的前缀组件**是 reparse point，且不在豁免表里 | 需要 junction/symlink。走查是逐组件的：`internal\risk\pathresolver_windows.go:60-92`（`GetFileAttributes` 逐段，`:87` 命中即记），判定在 `internal\risk\pathresolver.go:117-125` |
| ③ 默认档会不会发生 | **不会**。默认 `reparse_point_exceptions = []`（`docs\specs\SPEC-06-security-gatekeeping.md:56`；字段 `internal\config\schema.go:479`）＝没有豁免，但只有当数据根链路上真有 junction 才响（路径来自 `cmd\wisp\run.go:1109` 的 `dataDir\artifacts`，`dataDir` 见 `cmd\wisp\doctor.go:247` 起）；产物文件名被百分号转义（`internal\agent\spill.go:246-273`），模型无法用一枚 call id 造出穿越 junction 的名字 |
| 泄漏的内容本身 | 两枚可达错误都是**静态串**：**不含路径原文、不含用户配置值、不含授权根列表**。台账 `docs\reports\pending-and-issues.md:11266` 的"今天真的在漏"＝**接缝级**读数（`.scratch\wisp\probes\174\r3\impl.md` §2，junction 由测试自己 `mklink /J` 造），不是"某台机器上正有东西出去" |
| 本机被入侵的证据 | **无**。未定性；本程也没有以任何方式跑程序 |

⇒ 一句话现状：**码在、拼装点在、接缝上真响过；今天没有任何一条生产路能让模型看见那串字**。它是潜伏形状，触发条件写在别人正排队的那格里（票 174 AC#2b／`Q-60` 甲——把读回接进生产的那一刻，条件①就成立）。

## §3 四形代价表

| 形 | ①用户看得见什么（零术语） | ②最坏后果形状（谁还能看到那串字、看到能干什么） | ③改动面枚数 | ④会不会把别的判据打红 |
|---|---|---|---|---|
| 甲＝只说"这条路被拒"，不给原因 | 你在设置里不会看到任何变化；AI 自己看到的只有"这条路径现在读不到"，**不说为什么** | 那串英文技术词谁也看不到了（包括模型）。代价反过来：模型只拿到"被拒"、拿不到类别，只能猜或反复撞同一扇门 | 产码 1 枚文件 1 行（`internal\tools\task.go:837`）**＋判据 2 枚文件／3 枚函数（4 处断言行）必须跟着改** | **会**：3 枚判据直接红，且"为了让它绿而删断言"＝放宽断言，`AGENTS.md` §1.1 逐字禁止 |
| 乙＝给类别（"规范化这一步被拒"），不给键名与路径值 | 屏幕上没有变化；AI 看到的是"这条路径现在读不到，C26 连规范化都没通过"——知道是哪一步卡住，但看不到那个配置项叫什么 | 模型只知道类别、不知道开关的名字，所以它没法照着名字劝你改哪一行；出网内容缩成一句中文短句 | **1 枚文件 1 行**（同一行末尾的 `（"+err.Error()+"）` 去掉）；判据 **0** 枚要改 | **不会**（0 枚）。目标串与已上线的另一条腿逐字相同（`internal\agent\spill.go:229`），所以这是把两形拉齐，不是新造文案 |
| 丙＝保持现状全文＋只在本地日志留全文 | 屏幕上没有变化；AI 今天就这样看到那串英文技术词（`risk: path traverses … reparse_point_exceptions`） | 那串字随提示词出网到你配的模型服务。内容固定、不含你的路径。能被看到的人＝服务方；看到能干什么＝知道有个叫 `reparse_point_exceptions` 的豁免开关（这词本来就是公开文档里的）。真风险不在这串字，在**拼接点无界**：哪天一枚带路径的错误进了 `Canonicalize` 的返回集，路径就跟着走（现成例子＝`internal\risk\pathresolver.go:95-96`，那句带两枚变量位＝路径原文） | 产码 **1 枚文件 1 行**（今天没有任何日志写回执，要**新增**）；判据 0 枚；牵动契约 0 张 | **不会**打红判据；但**没消除**"无界"那一格（见下） |
| 丁＝这件事不做了 | 什么都没有，跟你今天用的一样 | 形状原样留着。等有人把 `ArtifactPath` 接到模型读得到的那一行（AC#2b 正在排队），那串字第一次进提示词的**当天不会有任何一格先响**：现有泄漏钉的 11 枚禁词里**不含** C26 词汇本体（`.scratch\wisp\probes\174\r3\impl.md` §6.1 逐字"本钉的射程刻意不含它"） | **0 枚文件、0 张契约、0 批准** | 不会。但也不销任何一格：票 174 的 AC#2b／AC#3／AC#4／AC#5 本来就未勾，不做≠结题 |

**④栏命中名册（逐枚，附"射程够不够得到这一形"）**

- `internal\tools\task_output_canonicalize_fail_174_test.go:85`（`TestCanonicalizeFailureFailsClosedInReply174`）要求 `连规范化都没通过` 在场 → 够得到**甲**（甲把它剥掉＝红），够不到乙/丙/丁。
- 同文件 `:130`/`:148`（`TestCanonicalizeFailureIsNotTheUnwiredArm174`，两臂对拼）→ 同上，只够得到**甲**。
- `internal\tools\task_output_pointer_notice_test.go:450`（`TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3` 正向半）→ 只够得到**甲**。
- `internal\agent\spill_pointer_honesty_174_test.go:172`（`TestSpillPointerCanonicalizeFailureQuotesNoUpstreamText174`）钉的是**另一条腿**（`internal\agent\spill.go:229`）→ 这一形今天**已经**是乙；改 task 腿不打红它，反而让两腿文案一致（够不到本格任何形，除非有人反向把 spill 腿改成甲）。
- `internal\tools\task_output_pointer_notice_test.go:300`（`TestPointerNoticeKeepsTheD153StubShape`）／`:350`（`TestHealthyReplyStillMatchesThePreFixTemplate`）＝两枚所谓"模板冻结钉"→ 射程是**指针句与头尾**，`:350` 只在 `notice` 为空的健康臂上逐字比；四形都不动指针句 ⇒ **都够不到**。
- 指针句 `全文见 <path>` 的钉：`grep -rn '全文见' --include='*_test.go' internal/` → **26 处命中／10 枚文件**（逐文件枚数 `7/2/1/2/1/1/2/3/6/1`；含票 175-r2／177／183／185 的 C25 R4 豁免窗口）→ 四形都不删指针 ⇒ **都够不到**。
- 括号那截有没有钉：`grep -rn 'risk: path traverses' --include='*_test.go' internal/` → **0 命中** ⇒ 没有任何判据要求原文在场 ⇒ 乙/丙零改判据。
- CI 扫描器：这一行的改动不含 `go func(`／`filepath.Clean|Abs`／明文密钥／墙钟差 ⇒ `tools/d22scan` 五形俱不受影响（**静态推理，未跑，派单禁跑**）。

## §4 契约轴到底挡在哪

- 最接近的冻结契约是 **C26 `PathResolver`**（`docs\PLAN.md:1376`）。逐字（引文里剥掉了原文的内嵌反引号，其余一字未动）：`**唯一**的路径规范化入口：展开 → 绝对化 → 取句柄真实路径（Win: GetFinalPathNameByHandle）→ **拒绝 reparse point（symlink/junction）/UNC/\\?\`/8.3 短名逃逸**。白名单与黑名单**都必须**经它`。→ 它冻的是**行为**（默认拒绝、唯一入口），没有一个字规定**错误串的措辞**。
- 另一枚是 C1 `Tool`（`docs\PLAN.md:1351`：`name / description / JSON Schema 参数 / **RiskLevel** / execute(ctx, params, onUpdate)`），冻的是接口形状，也不管这句话怎么写。
- 真正被冻结的那段文本出自 **D 不出自 C**：`docs\PLAN.md:431`（D15 ③ 行，逐字：`上下文里只留 **头 500 token + 尾 200 token + 总长度 + 文件路径**`）与 `docs\PLAN.md:2564`（`task.output` 行，逐字：`**只截不指＝不合格**（必须给可续读的路径）`）。
- **此处为射程判断**（读冻结件内部只为裁射程，未改 `docs\PLAN.md` 一字，也没改 `internal\panel` 三枚冻结件任何一字）：这两句管的是"头／尾／总长／路径"四样必须在场，即 `internal\tools\task.go:561-563` 那个模板；`pointerNotice` 是**塞进 `%s` 槽的附加句**，`docs\PLAN.md` 里没有任何一句规定它的字。剥掉括号那截，四样一样不少。
- **结论：契约轴不挡。** 挡住的是另外两样：①**判据轴**——`连规范化都没通过` 七个字被 3 枚判据钉着（只挡形甲）；②票面 AC#3(iii)/AC#4 把 `[fs] allowed_dirs` 默认值与 `docs\PLAN.md` 文字锁在"未批零字节"名下（乙/丙都不碰它们）。⇒ 台账 `docs\reports\pending-and-issues.md:11266` 那句"改它＝契约轴、须 owner 裁"按我这一趟的现读**不成立**，请编排者复核这个口径（我不改台账一字，也不改 `task.go` 一字）。

## §5 我读到的 vs 我没量到的

读到的：§1–§4 全部，每条带推导式；四形的改动面枚数与打红名册是**静态**量出来的（读码＋grep＋`git log`/`git show`/`git cat-file`）。

没量到的（具名，不含"应该没问题"）：

1. **端到端一发都没跑**（派单禁 Go／禁起 `wisp`；174-c1 另记本机 exe 起不动 rc=127／`0xc0000135`）。⇒ §1 的投递链是读码链，不是实测链。
2. "`ArtifactPath` 生产写入者只有 1 枚"只在 **HEAD 已提交面**（终态复量＝`a73edf62`，交件时间 `2026-10-04T09:16:46+08:00`）上量；本程中途那几枚未提交的 M 文件（含 `internal\tools\bridge.go`）**内容我没读**——若其中某一枚正把 `Backfill` 往环路结束之前挪，§2 条件①的"不成立"就要当场重判。
3. C25 片段索引会不会把这句注意话**自己**收成 taint 片段再回喂模型——只读了 `internal\risk\provenance.go:558-583` 的入参归一化与窗口排除，没读索引消费端。
4. `fs.read` 在同一枚 junction 上的对照读数（174-r3 §8.3 同样登记为缺口）。
5. 非 NTFS 宿主上四形的可见差别（生产平台是 Windows，D7）。
6. AC#2b 那半句"或批一张 L2 卡"——`internal\tools\task.go:839-840` 今天仍只给一条回来的路（另一格，未评）。
7. 门禁实测（`scripts/d22scan.sh`／`gate-clauses.sh`／名册 `comm`）——**禁跑**，§3 ④ 栏那条只是静态推理。
8. "`reparse_point_exceptions` 这词对用户本来是否可见"我只核到 2 处（`docs\specs\SPEC-06-security-gatekeeping.md:56`、`internal\config\schema.go:479`）；面板/引导文案里是否出现过——没量。若要论证"它不是秘密"，这两枚不够。

## §6 给 owner 的一句话

你的电脑没有在往外送任何关于你的东西，这一次要点的不是"漏了什么"，而是"要不要顺手把一句英文技术词从 AI 看得见的地方拿掉"：拿掉它只改一行、不碰任何你批准过的规矩，也不改屏幕上任何东西；不拿掉也不出错，只是等哪天那条"让 AI 自己去读长答案"的路真接上，这句会连同以后可能的路径原文一起，先被 AI 看到。
