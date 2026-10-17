# 328 — L2 审批卡片在生产里**点不亮任务行**：卡片只带 `CorrelationID`，而 `callCorr` 拼出的是 `task#call_abc123`，名册那把尺只认"最后一个 `#` 后面全是数字"

- Status: **已立，待派**（2026-10-16 12:5x 编排者立；来源＝非实现者终裁腿 `220-v1` 交回的两枚"交编排者"事项之一（件 §…＋提交 `4366c1ca` 的提交信息），台账 `A1006`；票 220 的 `AC#5` 因本案未定而未判，本票＝那一格缺的那块世界）
- 类型：**接线／形状缺陷**（不是判据仪器缺陷）。射程＝`internal/panel/**` ＋ `internal/agent/**`，⛔ 动 `cmd/wisp`（本波归别腿）／`frontend/**`／`design/**`
- ⚠ 行号只当**读数**、⛔ 当判据；下面每条附可重跑的尺。

## 这是什么（人话）

球上的"有东西在等你回答"这一格，靠的是**两串号对上**：审批卡片身上带的那串号，和任务行身上带的那串号。

- 任务行身上是**任务号**（形如 `task-17`）。
- 卡片身上只有**关联号**一串（`internal/panel/approval.go:39` 的 `ApprovalCardView` 只有 `CorrelationID` 一枚字段，⛔ 任务号）。
- 名字是怎么拼的：`internal/agent/loop.go` 的 `callCorr(taskID, callID, index)` 在**模型给了调用号**时返回 `taskID + "#" + callID`（例：`task-17#call_abc123`），模型没给时才返回 `taskID + "#call-" + index`。
- 名册那把认号的尺 `splitClashSuffix`（同包 `subagent_roster_197.go`）只认**"最后一个 `#` 之后全是 ASCII 数字"**那一形（那是队列撞名时的改写形 `task-17#2`），注释逐字写着"accepting a non-numeric tail would treat a host's own `<task>-suffix` naming as if the queue had rewritten it"。

⇒ 三条摆在一起：** provider 给的 `call_abc123` 不是全数字 ⇒ 那把尺不剥 ⇒ 卡片那串号进不了"在等谁"的名册 ⇒ 任务行上的 `BlockedOnApproval` 永远是 false**。也就是说：**L2 卡片弹出来了，球/面板上那一行却不亮**，而今天所有用例都是在测试里直接把 `CorrelationID` 写成 `task#3` 那种全数字形造出来的，所以没人看见。

## 现量（起手逐条复算，⛔ 信这里的号）

| 事实 | 读数（编排者 2026-10-16 12:5x 现跑） |
|---|---|
| 卡片只带关联号 | `sed -n '39,45p' internal/panel/approval.go` ⇒ `type ApprovalCardView struct` 第一枚字段 `CorrelationID string`；`grep -n 'TaskID' internal/panel/approval.go` 命中 **0**（尺：同一文件逐枚数） |
| 名册侧只按关联号索引卡片 | `internal/panel/subagent_roster_197.go:212-214` `for _, card := range pending { indexWaitingKey(waiting, card.CorrelationID) }` ⇄ `:216-220` 对 **L1 窗口**是**两枚都索引**（`w.CorrelationID` ＋ `w.TaskID`）——差的就是那一枚 |
| 认号那把尺只认全数字尾巴 | 同文件 `:262-270` `indexWaitingKey` ＋ `:277-286` `splitClashSuffix`（`r < '0' \|\| r > '9'` ⇒ `return "", false`） |
| 生产拼号处 | `internal/agent/loop.go` 的 `callCorr`：`callID == ""` ⇒ `taskID#call-<index>`（⛔ 数字）；否则 `taskID + "#" + callID` |
| 两形都点不亮 | ① 带 provider 给的调用号 ⇒ 尾巴 `call_…` 非全数字；② 无调用号 ⇒ 尾巴 `call-3` 含字母与 `-` 也非全数字 ⇒ **两条生产路径都不落在那把尺认的形上**（我逐字读过 `callCorr` 两个分支与 `splitClashSuffix` 的判据；卡片侧字段名册我也现跑了：`grep -c 'TaskID' internal/panel/approval.go` ＝ **0**，整枚文件里那枚结构有 10 枚带 `json:` 标签的字段、⛔ 一枚不是任务号）⇒ 这一格请落地腿自己再跑一次并给用例 |
| 为什么今天全绿 | 名册用例把 `CorrelationID` **直接写成撞名形**（例：`...#2`），⇒ 仪器验的是"剥完之后能不能对上"，⛔ 验的是"生产真会拼出那一形吗"。★这就是本仓反复入账的"缺世界"那一形：**那道分岔所需的世界，仪器里没造出来** |

## 判据（AC）

- [x] **AC#0 现量普查（⛔ 没名册不许进实现）**：逐枚列出今天真会走到 `callCorr` 的调用源（模型给号／不给号两形、半双工干活路径与陪聊路径各一？），并对每枚给出**它拼出来的 corr 逐字节**＋**名册那把尺剥不剥**；顺带数今天有几枚用例的 `CorrelationID` 是硬写的撞名形（尺：`git grep -n 'CorrelationID:' -- internal cmd` 逐枚看尾部）。
- [x] **AC#1 选一形并写清代价**（三形起评，⛔ 第四形）：**甲**＝卡片也带上任务号（**要新加一枚进 Snapshot/审批名册的 JSON 键**＝契约面＋撞票 313 那枚键钉 ⇒ 走 `A##` 批准才许开工）；**乙**＝`callCorr` 改成拼在尾巴上放可剥的数字（例：`task#<seq>#<callID>`），代价＝改的是**所有** join 方的形状，含队列改写形与 `#call-N` 回退形，必须逐枚点名顶红谁现有的钉；**丙**＝名册那把尺改成"按已知任务号前缀匹配"，代价＝它注释里那句"⛔ 比这仓任何代码产生的都宽"要重新裁（这是**放宽**，⛔ 顺手做）。三形都要给"拆完之后还剩几枚零读数"的实测。
- [ ] **AC#2 ★正控必须有**：造一发**用生产真形**的读数——经真的 `callCorr` 拼出 corr（⛔ 测试里硬写 `#2`），走真的名册，断言那一行 `BlockedOnApproval` 必须为 true。★并自带去牙：把 `callCorr` 或 `indexWaitingKey` 任一支改回今天的形状 ⇒ 用例**必须红**；红名册与红句逐字进件。⛔ 只交"绿"不交这发对色＝本票不收。
- [ ] **AC#3 判定与批准语义不许动**：本票只修**看得见的接线**；⛔ 改门控档位、⛔ 让面板侧来源变成 L2 的合法出处（AGENTS §1.2 那条禁令）、⛔ 动 D43 转移表与三枚冻结件（`tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`ticket90_persist_test.go`）。完成判据＝这四面 diff 行数逐枚报 0。
- [ ] **AC#4 门禁四数＋颜色**：`gofumpt`／`tools/d22scan`／长度门／`go vet` 逐枚落 `rc=`；`./internal/panel` 与 `./internal/agent` 两包改前⇄改后**同一台面**成对、显式 `-timeout 30m`、四数逐枚点名（⛔ 混 `^=== RUN` 与 `^--- PASS` 两把尺），终态红名册与起手作差＝0。⛔ 跑 `./cmd/wisp` 整包、⛔ 跑全树 `./...`。

## 边界与禁区

- ⛔ 为了对上号而**把 L2 卡片的文案复用权限判定的文案**（票 242/244 那族"门拦下的话不许和权限的话同形"）。
- ⛔ 新加 Snapshot 名册键（＝甲形的代价，必须先落 `A##` 批准记录）；⛔ 新 C17 方法名、⛔ 新导出 `panel.*` 名 ⇒ 停手上报。
- ⛔ 用 `go build`／`go vet` 的 rc=0 当"接上了"——本票要的判据是**调用点与字节形状**（Go 不为未使用方法报错）。

## next

`328-a1`（只读普查腿＝AC#0 的名册＋三形代价）→ 我裁形并落 `A##` → 落地腿（写面收成 `internal/panel/**` ＋ `internal/agent/**` 中的选中的一枚；⛔ 同批两枚腿碰这两面）。**释放条件**＝`internal/panel/**` 本波空着（`216-r1` 已于 12:12 交完）＋ CPU 连发读数 ≤70%。

- 328-a1 只读普查交件（AC#0 的料，⛔ 判语/⛔ 选形/⛔ 翻框）：件①②③落 `.scratch/wisp/probes/328/a1/**`。名册——今天真产 corr 的调用点只 `internal/agent/loop.go:676` 一枚，走 `run`→`executeCalls`，入口两枚（`cmd/wisp/run.go:1164`、`internal/tools/subagent_197.go` 的 `child.RunAsync`，行号随并发腿 270-r2 漂移故仅读数），均半双工文本干活路径；陪聊全双工 Path C 与 realtime 到不了（`approval/gate.go` AdmitTextTask/admitCheck）。每形（算的）：支B `task#call-<idx>`、支A `task#<provider号>`，`splitClashSuffix`（`subagent_roster_197.go:277`，唯一读者 `:267`）都判非数字⇒join 不上。必答：**字面"末#后全数字"的 corr 确有、只由队列撞名/重放改写 `queue.go:168/:616` 造（样例 `<uuid>#call-0#7`），但剥出的 base＝`<uuid>#call-0`≠裸任务号⇒仍点不亮行**；能剥成裸任务号（恰一个#的 `T#<全数字>`）者＝无（支B/支A/队列/host mode-switch `run.go:878` 空 corr 回退 `approval-<seq>`／陪聊/realtime 逐枚排除）。代价列：甲须给 `ApprovalCardView`（`approval.go:39` 现无 TaskID）加 `taskId` 键＋名册 `:213-214` 也索引卡任务号，撞 `approval_test.go:105`（与 `frontend/src/lib/panel.ts:24` 双向对账＝跨面）＋C17（`SPEC-08:156/:235`）；票说的"313 键钉"实为 settings payload（`config_313_test.go:283`）、不覆盖卡键。乙（放宽那把尺）反转 `subagent_blocked_220_test.go:133 child-1#retry→false` 并拆 `:132/:134` 防过宽护栏。丙（corr＝裸任务号）逐枚打红"corr≠taskID/前缀+尾"那族（ticket285 `:73/:136`、283 `:222`、287 `:140`、loop_approval `:220`、242 `:139`、loop_golden `:79`）并退回 C18 塌键。AC#0 附带：硬写撞名形 corr 的用例只在 `subagent_blocked_220_test.go`（`child-1#9`×4／`#147`／`other-1#9`／头注 `task-7#9`），经 `card197()` 手造、不经 `callCorr`＝本票"缺世界"那形。
- 328-a1 **复跑交件**（料齐⛔ 判语⛔ 选形⛔ 翻框；件＝`docs/evidence/s1/328-a1-census.md` 81 行／19209 字节、占位符尺 0 命中且正控三枚 docs 看得见；读数＝`.scratch/wisp/probes/328/a1/logs/**`，提交 `8e244a2c`／`4c475daa`／`0e4c5d40`）：锚 **`7535d83b`** 现认，起手四把尺逐枚重跑＝未勾 **5**（`:31-35`）／`approval.go:39` 结构＋`:40` 首字段／`TaskID` **0 命中**（正控 `CorrelationID` 命中 `:28/:40/:78`）／名册不对称 `:213-214` ⇄ `:216-221`／那把尺 `:277-288` 全仓唯一读者 `:267`。★**更正前一枚普查腿（`59c46b7a`/`fdea4402`/`135e9ed0`）两处**：① 它把本票 `AC#1` 的**乙与丙互换**了（票面乙＝改 `callCorr` 尾形、丙＝放宽那把尺），② 生产入口行号已漂（本锚实测 `cmd/wisp/run.go:1176`、`internal/tools/subagent_197.go:372`，非 `1164`/`367`）。★名册实测（仓外纯 stdlib 探针，四枚函数逐字重打、`diff -B` 对锚逐字节相同；真实例任务号 `e2bbeaab-…`）：支B `T#call-0`／支A `T#call_e1`／`T#call_unnamed_0`／队列叠改 `T#call-0#7`（剥＝true 但 base≠T）／replay／`approval-<seq>`／宿主卡 `host:*` **逐形皆点不亮**，唯一点亮形＝`T#<纯数字>`（生产不产）；★**票面乙那一例 `task#<seq>#<callID>` 实测也点不亮**（剥＝false），改成 `T#<callID>#<seq>` 剥出 `T#<callID>`≠T，⇒ 乙要点亮只能丢掉 provider 号，而那一形与队列撞名（`queue.go:168` 用 `q.seq`）**同一命名空间**。★键钉名册现量（导出台树，⛔ 共享树跑 go）：`approval_test.go:105` 与 `composer_test.go:52` **起手即红**（红句逐字 `Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`＝与卡片无关）⇒ 甲加键只多一句同族红、**颜色分不开**；票 313 那族钉的是 settings payload 三组键（顶层 11／provider 行 10／voice 5，`config_313_test.go:290/:317/:330`）**碰不到卡片键**；`json:"taskId"` 已被 `subagent_roster_197.go:107` 占用（`frontend` 里 `taskId` 零命中⇒同名不同段）。★乙的硬钉＝`ticket285_corr_distinct_rulers_test.go:136`＋`loop_golden_test.go:77-81`（两枚今天量到绿，`agent-color-targeted2` 5/5），丙的硬钉＝`subagent_blocked_220_test.go:133/:134` 两枚子用例，其余六枚钉前缀不等的用例（`242:139`、`283:222`、`285 rows:75`、`287:140`、`loop_approval:219`）乙丙都碰不到（`tools-color-targeted` 5/5）。★零读数具名：`frontend` 里 `blockedOnApproval` **0 枚文件**（正控＝同尺族 `correlationId` 命中 4 枚文件）、`cmd/wisp/subagent_carrier_197_test.go:598-604` 逐字 "blockedOnApproval is never asserted on a production packet"、`panel_snapshot_deliver_329.go:27-33` 自陈页面侧 `window.wispOnPanelSnapshot` 接收器不存在；丙放宽后 `splitClashSuffix` **再无任何测试直呼**（全仓 def＋doc＋`:267` 三命中）。必答三题在件 §三（①生产读者逐枚点名：名册 `:214`、`panel_pump.go:384`、`pump.go:101`＋前端 `App.tsx:322`/`l2-approval-card.tsx:265`/`approval-screen.tsx:58`；②三处注释逐字已引，含 `:56-58` 那条自己写着 `<taskID>#<callID>, :603/:676`；③本票＝220-v1 交回的**①**、`AC#5`＝**③**，尺＝票 220 第 29 行＋台账 `A1006` `:18793`）。门禁四数：gofumpt（本腿零 Go 产码⇒rc=2 报 `illegal character`＋正控 rc=0）／`d22scan` rc=0 clean／长度门 rc=0 GREEN（tracked 14782）／`go vet ./internal/panel ./internal/agent` rc=0（锚上字节，⛔ 共享树）。本腿一枚框没翻、⛔ 改判据／票名／`Status:`／`-done`、`internal/**` 与 `cmd/**` 零字节、⛔ 动 PID 21892。

- [2026-10-17 11:5x +08] **编排者裁形＝丙-带边界**（撤销口令「**328 改裁甲**」＝回到本条之前、那时三形都没拍；「**328 改裁乙**」同理）。凭据＝只读普查腿 `328-a1` 交件（件＝`docs/evidence/s1/328-a1-census.md` 81 行／19691 字节，锚 `7535d83b`）＋**我自己复跑的两把尺**（HEAD `469839a9` 现跑：`git grep -n callCorr -- internal cmd` 去测试＝`loop.go:367`/`:594` 注释、`:603` def、`:676` 唯一调用、`subagent_roster_197.go:58` 注释 ⇒ 生产铸造点仍只一枚；`git grep -n splitClashSuffix`＝def `:277`＋doc `:272`＋**唯一读者 `:267`** ⇒ 普查那句"再无任何测试直呼它"为真）。两把尺都⛔ 抄件里的号。
  - **为什么不选甲**：甲要给 `ApprovalCardView` 新加一枚 JSON 键 ⇒ ① 撞 C17 契约面（`SPEC-08:156`/`:235`），**必须先落 `A##` 人工批准**；② 要动 `frontend/src/lib/panel.ts:24-47`＝机主外部 agent 的写面，本编队零写面；③ 键钉用例 `approval_test.go:105` **起手即红**（红句逐字 `Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`，与卡片无关）⇒ 甲加的那句红**颜色分不开、只有文案分得开**。三笔代价换一格本来就在别家写面上的读数，不换。
  - **为什么不选乙**：★普查把我票面举的那一例**当场量废了**——`task#<seq>#<callID>`（票 `AC#1` 乙的原文例子）剥＝**false**、点不亮；把号挪到尾 `T#<callID>#<seq>` 剥＝true 但 base＝`T#<callID>` ≠ 裸任务号、**还是点不亮**。⇒ 乙要点亮只能丢掉 provider 号，而丢掉之后那一形与队列撞名改写（`queue.go:168` 用 `q.seq`）**同一命名空间**，"这串号是谁改的"从此不可读；另要动两枚今天绿着的硬钉（`ticket285_corr_distinct_rulers_test.go:136`、`loop_golden_test.go:77-81`）＝**改形状顺带改判据**。
  - **裁丙，但把边界钉死**：丙的朴素写法（`strings.HasPrefix(corr, 已知任务号)`）会连 `child-1x` 这种**没有分隔符**的前缀也认 ⇒ 那正是"把两枚任务混起来"那一形，也是本仓禁手。⇒ 落地形状＝**剥出来的段必须逐字节等于一枚在册任务号，且原串里紧跟其后的是 `#`**（等价写法：在**第一枚** `#` 处切，切出的左段与在册任务号全等）。★这不是第四形——票 `AC#1` 丙那句"按已知任务号前缀匹配"本来就含"匹配到多长"这一问，我这里把它**钉成带分隔符的那一支**；`AC#1` 原文一字未改，本条即裁定记录。
  - **形状选对了的证据（三条，都由这条裁生）**：① `T#call_e1`（支A）／`T#call-0`（支B）／`T#call-0#7`（队列叠改）／`T#9`（撞名形）**四形全点亮**，而今天只有最后一形点亮且生产不产它；② `approval-<seq>`（空 corr 回退）与 `host:*`（宿主发起卡）**不含 `#`** ⇒ 照旧不点亮，不必为它们开口子；③ ★**`child-1x` 那一枚子用例必须继续为 false**——它是丙与"混任务"之间的唯一隔板。
  - **落地腿的硬框（写进派单，缺一不收）**：ⓐ `AC#2` 正控那发要经**真的 `callCorr`** 拼 corr（⛔ 测试里硬写 `#2`）；ⓑ 子用例 `subagent_blocked_220_test.go:133`（`child-1#retry`）**允许翻向**，但要逐枚具名写"翻成 true 是买真相、⛔ 不是为变绿"，⛔ 删用例、⛔ `t.Skip`；ⓒ `:134`（`child-1x`）**原地不许动**，且它自己就是去牙突变那一发：**把分隔符那半边抹成纯前缀匹配 ⇒ `:134` 必须红**（成对第二发＝把 `callCorr` 或名册那把尺改回今天的形状 ⇒ 正控必须红，两发落点⛔ 同一枚断言行）；⓭ 名册 `:272-276` 那句自陈（"recognises the queue's rewrite and nothing else"／"a wider match than any code in this tree produces"）**改完必须说实话**——落地后生产corr 真落在这把尺上，那句"没有代码产生它"就地过期（同一味账＝票 333 那三处过期散文，票 335）；⓮ `AC#3` 四面（门控档位／面板侧 L2 出处／D43 转移表／三枚冻结件）diff 逐枚报 0。
  - **★据实登记本票修完还差的那一跳**：屏幕那一行今天**仍不会亮**——`frontend` 里 `blockedOnApproval` 命中文件数 **0**（正控＝同尺 `correlationId` 命中 4 枚），且 `panel_snapshot_deliver_329.go:27-33` 自陈页面侧 `window.wispOnPanelSnapshot` 接收器不存在，`cmd/wisp/subagent_carrier_197_test.go:598-604` 逐字写着 "blockedOnApproval is never asserted on a production packet"。⇒ 本票只修 **Go 侧那一跳的配对**，⛔ 任何人把本票写成"用户已经看得见"。那一跳归票 329／前端会话。
  - **排程（⛔ 现在就放）**：`internal/panel/**` 正被 `335-r1` 在写（测试面＋`display_scrub_216.go` 整行注释面），而本票 `next` 自己写着"⛔ 同批两枚腿碰这两面" ⇒ **`328-r1` 待派**，释放条件＝`335-r1` 交件入库＋CPU 连发 ≤70％。本轮 `AC#0`／`AC#1` 两枚框我翻（`AC#0`＝普查料齐、`AC#1`＝裁形已拍），`AC#2`/`AC#3`/`AC#4` 三枚留给落地腿＋非实现者裁腿。
