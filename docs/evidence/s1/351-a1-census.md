# 票 351 普查件（腿 `351-a1`／★第一趟＝只普查、0 字节产码、⛔ 落地形状、⛔ 翻框、⛔ push）

HEAD 起手锚＝本腿第一笔 `git log -1` 现认＝`14e5d8395adaf54bec14199f6858382b57fa49dd`（`dev`，2026-10-20 09:58:30 +08）。
散件＝`.scratch/wisp/probes/351/a1/logs/`：`00-rulers.md`（五把尺＋正控＋四数复尺）／`01-ac0-roster.md`（名册）／
`02-ac1-verdict.md`（定性＋三支凭据＋到达数分子）／`03-ac2-shape.md`（形状＋AC#3 判定）。
★台面＝⛔ 任何 `go test`（本波 `296-v2` 正在 `cmd/wisp` 上定向色）；凭据只来自**盘上既有日志＋CI 日志＋`git grep HEAD`**。

## 一、结论（三句）

1. **`AC#1`＝仪器形状（运行学时序）**，⛔ 真缺陷、⛔ 单纯 CI 环境。机制＝产码"先记账、后印字"
   （`cmd/wisp/run.go:1423 cards++ → :1431 live.record → :1438 bookWaitingState → ★:1441 Fprintf`，
   `:1425-1429` 注释逐字写着「booked before anything is printed」）；仪器这一侧
   `config_reload_223_test.go:542 awaitCard` 轮的是 **gate 侧 ledger**（`:259`）、`:543` **单次**读 stdout＝
   两立面之间⛔ 同步边＝窗口由产码顺序保证。
2. **红的那一格⛔ 在量"文案缺那枚 key"**：两次到达（件02 表）都是"**卡片整段还没进 stdout**"——
   `39101060753` 那次是 `:403`（没渲染）⇄`:406`（没点名）**成对响**、`39244425666` 那次红句 dump（`win.norm` L1841-1842）
   逐字只有两条启动横幅。⇒ "卡片在场但⛔ 点名"那一形**盘上 0 次到达**。
3. **〔真缺陷〕排除的直接凭据在发红日志自己身上**：同发 L1843 已打
   `locked loosening approved via L2 re-confirmation section=risk keys=[risk.permission_mode]`，
   这句的 `keys` 与喂卡片 `Reason` 的是**同一次调用的同一枚切片**（`internal/config/manager.go:180`⇄`:184-185`
   vs `cmd/wisp/config_reload.go:283-292`→`gate.go:606`→`ui.go:35`→`run.go:1441`，中间无人改写）。

## 二、`AC#0` 射程名册（尺 A⇄B 交叉；★⛔ 单把尺报裸枚数）

| 尺 | 逐字 | 读数 | 正控／射程 |
|---|---|---|---|
| A | `git grep -nE '^func TestTicket' HEAD -- cmd/wisp` | **86 枚**／21 文件 | 非 `_test.go` 命中＝**0**；锚⇄不锚两根**同 86**（本包无缩进定义＝无差异，如实记） |
| B | `git grep -n 'awaitCard' HEAD -- cmd/wisp internal` | **5 命中**＝1 注释＋1 定义＋**3 调用**（`:395`/`:467`/`:542`） | ★全仓不限包＝**24 命中**、多出 19 枚全是 `.scratch`⇄`docs` 散文 ⇒ 射程必限包 |
| C | `git grep -nE '\[确认 ' HEAD -- 'cmd/wisp/*_test.go'` | **5 命中**／4 文件 | 正控＝`:402` 在场 |
| D | `git grep -nE 'does not name' HEAD -- 'cmd/wisp/*_test.go'` | **60 命中**／29 文件 | ★绝大多数⛔ 卡片（错误消息／refusal 类）＝⛔ 与"卡片文案"混记 |
| E | `git grep -nE 'awaitStdout(Since232)?\(' HEAD -- 'cmd/wisp/*_test.go'` | **8**＝2 定义＋6 调用／2 文件 | ★`:543`⇄`:401` **不在**此列＝未等单读 |

**甲族（票面问的那族＝`awaitCard` 之后断文案，`config_reload_223_test.go`）逐枚三格**：

| 枚 | 函数（尺 A） | key 名断言 | Errorf⇄Fatalf | 到达数分子 | 未等单读 |
|---|---|---|---|---|---|
| `:395` | `TestTicket223HandEditedFsLooseningCostsAnL2Card`(`:390`) | **有** `:405 fs.allowed_dirs` | 文案两格都 **Errorf**(`:403`/`:406`) | **有** `:432 windowCount()!=1` | ★`:401` |
| `:467` | `TestTicket223RefusedLooseningKeepsOldValues`(`:463`) | ⛔（该格⛔ 读文案＝⛔ 射程） | `:472` **Fatalf**（`result=deny`） | ⛔（`:484` 是 `time.Sleep`＝⛔ 分子） | — |
| `:542` | **本票那枚**(`:535`) | **只有它** `:544 risk.permission_mode` | **Errorf**(`:545`)＝红⛔ 中止后续 | ★**⛔ 任何分子**（无 windowCount／无 Count） | ★`:543` |

⇒ 派单点名那一问的**逐枚答案**：文案类断言**全族 Errorf**（本文件 33 处）；★**28 处 `Fatalf` 一枚都⛔ 在文案格**
（全在 helper／前提格：`:72`/`:88`/`:92`/`:96`/`:100`/`:110`/`:130`/`:161`/`:186`/`:221`/`:244`/`:266`/`:278`/`:296`/`:316`
＋`:348`/`:359`/`:363`/`:426`/`:439`/`:472`/`:487`/`:516`/`:539`/`:548`＋`:701`/`:787`/`:790`＝逐名 28 枚，⛔ 一枚落在文案格）；★写成**到达数分子**的只有
`windowCount()` 四枚（`:379`/`:432`/`:518`/`:679`）＋`strings.Count` 两枚（`:595`/`:668`＝票 272「窗口里恰好一份」那把，
`:668` 特意在 await `state=applied` 之后才读，注释 `:658-666` 具名）。**本票那枚＝三味全缺**。

其余同形（丙／乙族，具名⛔ 混记）：`panel_l1_blocked_220_test.go:240`（两 stream **任一含即过**＋Errorf＋⛔ 分子）／
`panel_pump_test.go:232`＋`:443`（★"在场"＝`t.Fatalf`、比对＝`t.Errorf`＝**本票要的两格形状在本包已有先例**）／
`zz315a1probe_test.go:377`（轮询＝有到达、⛔ 文案断言＝⛔ 射程）／`run_mode101_test.go:534`（分子＝有）／
`approval_always_201_test.go:209`（key 名断言＝有、等＝`waitForSentence187`）；
★"审计已到⇒操作员句单读"的次级同形＝`:344`/`:448`/`:453`/`:478`/`:521`/`:676` 六枚（＋甲族 `:401`/`:543`）。
★settle 现量＝plant **7** 处里只 **3** 处先等 `state=armed` ⇒ `:466`/`:506`/`:541`/`:599` **四枚裸种**（本票＝`:541`）。

## 三、`AC#1` 三支各自的盘上凭据（★"该形状到没到达"＝逐台面具名 needle，⛔ 总数）

| 台面 | 该枚颜色（尺＝`^--- (PASS|FAIL): <名>`） | `:403` 型到达 | `:54x`/`:406` 型到达 |
|---|---|---|---|
| CI `39244425666` test-windows | **FAIL** L1844（2.22s） | 0（该枚⛔ 在场格） | **1** L1841 |
| CI `39101060753` test-windows | PASS L6480（本枚绿；fs 姊妹枚 FAIL） | **1** L6433 | **1** L6435（★与 `:403` **同段成对**） |
| CI `39071959151` / `39124122121` | PASS L3304／L331 | 0 | 0 |
| 本机 `orch227-b/full.log` | PASS L352（2.20s） | 0 | 0 |
| CI test-core（Linux） | ⛔ 现身＝⛔ 射程 | — | — |

⇒ **到达＝2 台面／跑到 5 台面**，两次同形（⛔ 一次＝"卡片在场但没点名"）。正控＝同尺在同件命中已知在场句
（`win.norm` `热加载`＝4、`答复监听`＝1，原始⇄`-Encoding utf8` 两读法**同数**＝尺灵）；
`core.norm` 的 `TestManager`＝24／`cmd/wisp`＝16＝⛔ 失明（故 `TestTicket223`＝0 是射程、⛔ 是绿）。

- **〔CI 环境〕＝⛔ 独立成因（放大器而已）**：派单点名的反例本腿独立复现——`REFUSED` 行总数
  红发 `win`＝**122** ＜ 绿发 `old39071`298／`old39101`254／`old39124`326／本机 269 ⇒ **与红反相关**；
  `result=deny` 五件＝0（正控同上）。runner 负载＝**未知**（`gh run view --json jobs` 只有 startedAt/completedAt；
  本腿⛔ 任何 runner 侧 CPU 读数）。★"机主常驻 `wisp.exe` 抢键"那形＝**⛔ 射程**：
  尺＝红窗 `win.norm` L1780..1860 grep `hotkey|occupied|RegisterHotKey|ball:`＝**0**（正控＝同件 L1841 在场）。
- **既有独立读数只对回、⛔ 采信**＝`289-v1` 21 发复跑（隔离×3／`-count=5`／组串／`-cpu=1`×2／5 进程并发×5／原状×2）
  **零复现**、并发把 2.0s 压成 2.8–3.0s 仍⛔ 红（`docs/reports/pending-and-issues.md:14857` 在册）＝负载改时序⛔ 改结果。
  ★本件给的是**定性＋机制＋分子**，⛔ 写"已知摆动"那枚弃权章（票面禁区第三条照守）。
- 同码三处两绿一红的字节级证（本腿现跑）＝`config_reload_223_test.go` 于 `a5080612`⇄`6f808a83`⇄`HEAD`
  **同串 `f1cac524055aaac000ced549296833ae99eb5a71`**；产码 `config_reload.go` 同串 `990d148b...`；
  `git log a5080612..6f808a83 -- cmd/wisp`＝**0 枚**。四数＝本腿脚下重跑五份归一件、与 `346-a1` 留档**逐数同**。

## 四、`AC#2` 形状（★只给形状，详见件03）＝三格各自可失败

- **G1 settle**：plant 前 `awaitAudit(state=armed)` ＋ ★基线分子 `before := windowCount()` ＋ `awaitStdout("配置热加载已接管")`；
  ★三味**全是现成尺**（`:117`/`:148`/`run.go:929`；先例 `:340`+`:344`）＝**⛔ 需产码接缝**。
- **G2 通知在场**：`awaitCard` ⛔ 变 → `windowCount()-before == 1`（分子）→ ★`awaitStdout(t,"[确认 L2 "+configReloadTool+"]")`
  单独成一格；两格并列的先例＝**同文件 `:402`⇄`:405`**。⛔ 需接缝。
- **G3 文案点名**：在 **G2 返回的那条流**上断 `Contains(out,"risk.permission_mode")`＝`:544` 那枚串⛔ 变、
  只把**读源**从 `:543` 的单读换成 await 返回值。⛔ 需接缝。
- ⇒ G2 红⇄G3 绿＝"卡片到了、没点名"（那才叫真缺陷）；★⛔ 放宽任何既有断言、⛔ `time.Sleep` 哄、⛔ `t.Skip`。
- 落点＝**就在 `cmd/wisp/config_reload_223_test.go:535-556` 那枚函数体**（⛔ 新文件、⛔ 新 helper）；
  ⚠ 与票 227 `AC#6` 同占 `cmd/wisp` 测试面＝同包只一枚写腿。丁族那四枚裸种是否顺手补＝**归编排者裁**（本腿只具名）。

## 五、`AC#3`＝**未触发**（⛔ 动 `cmd/wisp/config_reload.go`）；两支形状先具名钉住以免重普查

1. 卡片**已经**带 key 名那一跳＝`config_reload.go:283-292` 的 `Decision.Reason` → `gate.go:606` → `ui.go:35` → `run.go:1441`；
   消费＝人眼 ⇄ `panel_pump_test.go:232`/`:443` 那种整段解析。★盘上已渲染件（逐字含 key）＝
   `.scratch/wisp/probes/221/v1/logs/ac5-full.txt` L142（⚠ 只 `fs` 那支；`[risk]：[risk.permission_mode]` 的已渲染件
   盘上＝0 命中（尺＝`git grep -F -n '<那句卡片前缀>' HEAD -- ':!.scratch/wisp/probes/351/a1'`＝0；★⛔ 排除本腿写面＝
   它现在会吃进我自己这份件的这句描述＝`A1045`⇄`A1033` 那族在本腿身上复发，具名记我）⇒
   risk 那支靠源码同路＋同发 keys 读数，此边界已具名写清）。
2. ★若改判"缺陷"指的是 **ledger 侧⛔ 携带 key**＝`internal/agent/approval/replies.go:70` `ReplyCard` 字段集
   （`CorrelationID/TaskID/Tool/Level/Grant/Paths/Window/Deadline`）⛔ `Reason`⛔ `Keys`；
   `config_reload_223_test.go:399-400` 注释已具名记录「ReplyCard has no reason field, so the naming is asserted off the console」
   ⇒ 那一支**需新产码接缝**＝⛔ 本票⛔ 本腿，须走"第二道批准"（先例＝票 342 `342-r2`）。

## 六、本腿⛔ 定、写未知（★未定义即停，⛔ 自行假设）

- ★"这一次为什么没赶上"＝**未定**：机制窗口由源码顺序证明存在（本件第一节），但**该窗口在那一发被撞上的具体那一跳**
  ⛔ 盘上凭据能定（要造它得另开一发、⛔ 复跑能造）。本件⛔ 圆场。
- ★`awaitCard` 返回的 `card.Tool`⇄`card.Level`⇄`card.CorrelationID` 是否**足以**替代 stdout 那一味（＝G2 能否用 ledger 字段做在场分子）
  ＝**未定**，因 `ReplyCard` ⛔ `Reason`⛔ `Keys`（本件第五节第 2 条）＝"在场"能、"点名"⛔ 能。落地腿若只想补 G2 的 ledger 侧，
  G3 仍须读 stdout＝这一点我按源码现读写死、⛔ 猜。
- ★`lint`（`118032776906`）＝**未知**（`346-a1` 留档为连接超时；本腿⛔ 重取、⛔ 读成"⛔ 红"）。
- ★甲2（`:467`）那种"⛔ 读文案"的用例算⛔ 算"仪器族"＝我按票面措辞（"`awaitCard` 之后**只断言卡片文案含某串**"）
  判它**⛔ 射程**并具名留在名册里，若编排者要的是"所有调过 `awaitCard` 的函数"＝名册⛔ 变、只是射程换一把尺。
- ★丁族那四枚裸种：`AC#2` 是"新增／改写**那枚**常驻件"（单数）＝本腿只补本票一枚的形状、余三枚只具名。⛔ 我裁射程。

## 七、纪律格（`AC#4`）自证

⛔ 翻框（票 351 框尺终量＝**0 勾／5 未勾**，尺＝`grep -cE '^- \[[ x]\]'` 于 `git show HEAD:<票面>`＝5 格 `AC#0`..`AC#4`）；
⛔ 改票名／`Status:`／`-done`；⛔ push；★我方 `.go` 字节改动＝**0**（写面只有 `docs/evidence/s1/`＋`probes/351/a1/logs/`＋票尾一条 bullet）；
★0 枚 `go test`／⛔ 0 枚 `go build`／⛔ 跑四道门（⛔ 抢 `296-v2` 与 `342-r2` 的编译面⇄读数）＝
⇒ 门禁四数＝**未知**（⛔ 本腿射程、⛔ 读成绿⛔ 读成红；落地腿那趟该逐枚落 `rc=`）；
blob 副本全落仓外 `D:\tmp\wisp\20261020\351a1\`（只建⛔ 删）；写面⛔ 码位 U+00A7（逐枚现量＝0，另数 U+00B7＝0）；⛔ 任何 `.out`。
