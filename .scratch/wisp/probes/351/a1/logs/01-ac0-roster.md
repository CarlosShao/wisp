# 351-a1 件01 = `AC#0` 射程名册（两把尺交叉 ⇄ 逐枚〔key 名断言／Errorf⇄Fatalf／到达数分子〕）

射程＝**只含 `cmd/wisp/*_test.go`**（⛔ 产码、⛔ 注释件、⛔ `.scratch`⇄`docs` 的散文引用）。
HEAD＝`14e5d839`；行号＝`git cat-file blob HEAD:<path>` 的行号（＝工作树外读，⛔ 解到仓内）。

## 一、两把尺交叉（★⛔ 单把尺报裸枚数）

| 尺 | 逐字命令 | 读数 |
|---|---|---|
| A（锚函数定义） | `git grep -nE '^func TestTicket' HEAD -- cmd/wisp` | **86 枚**／21 枚文件 |
| A′（不锚） | `git grep -n 'func TestTicket' HEAD -- cmd/wisp` | **86 枚**＝与 A 逐字同（`Compare-Object` 空）⇒ 本包无差异，如实记 |
| A 的产码性正控 | A 命中里非 `_test.go` 的枚数 | **0** ⇒ 尺 A 射程天然只含测试面 |
| B（锚 `awaitCard`） | `git grep -n 'awaitCard' HEAD -- cmd/wisp internal` | **5 命中**＝1 注释＋1 定义＋**3 处调用**（`:395`／`:467`／`:542`，全在 `config_reload_223_test.go`） |
| B 的射程反证 | `git grep -n 'awaitCard' HEAD`（⛔ 限包） | **24 命中**＝多出 19 枚全在 `.scratch`⇄`docs` 的散文里 ⇒ **⛔ 不限包就报枚数** |

交叉结论＝尺 A 的 86 枚里，"awaitCard 之后断言文案"的只有尺 B 那 **3 处调用**（2 枚文件外＝0）。
但"**卡片文案含某串**"这把仪器⛔ 只长在 `awaitCard` 上，所以另加两把尺把全族扫出来：

| 尺 | 逐字命令 | 读数 |
|---|---|---|
| C（卡片字面） | `git grep -nE '\[确认 ' HEAD -- 'cmd/wisp/*_test.go'` | **5 命中**／4 枚文件 |
| D（"点名"族红句） | `git grep -nE 'does not name' HEAD -- 'cmd/wisp/*_test.go'` | **60 命中**／29 枚文件（★绝大多数⛔ 卡片＝错误消息/refusal 类，逐枚筛见下） |
| E（等了再读的那把尺） | `git grep -nE 'awaitStdout(Since232)?\(' HEAD -- 'cmd/wisp/*_test.go'` | **8 命中**＝2 定义＋6 调用／2 枚文件 |
| F（未等单读） | `git grep -n 'shown := r\.h\.out\.String()' HEAD -- cmd/wisp` | **恰好 2 处**＝`:401`⇄`:543` |

## 二、名册（逐枚标三格；★本票射程＝"卡片文案含某串"）

### 甲族＝`awaitCard` 之后读文案（尺 B∩尺 C/D＝票面问的那族，`config_reload_223_test.go`）

| 枚 | 函数（尺 A 行） | 文案断言 | key 名断言 | Errorf⇄Fatalf | 到达数分子 | 该次读⛔ 等 |
|---|---|---|---|---|---|---|
| 甲1 `:395` | `TestTicket223HandEditedFsLooseningCostsAnL2Card`（`:390`） | `:402` `Contains(shown,"[确认 L2 config.reload]")` | **有**：`:405` `fs.allowed_dirs` | 两条都 **Errorf**（`:403`／`:406`） | **有**：`:432` `windowCount()!=1`→Errorf `:433` | ★**是**（单读在 `:401`） |
| 甲2 `:467` | `TestTicket223RefusedLooseningKeepsOldValues`（`:463`） | ⛔ 无（只答卡＋读审计） | ⛔ 无 | `:472` **Fatalf**（`result=deny`）；`:476`/`:478` Errorf | ⛔ 无 windowCount；`:484` 是 `time.Sleep(3s)` ⛔ 分子 | 该格⛔ 读文案＝**不在射程** |
| 甲3 `:542` | `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`（`:535`）＝**本票那枚红** | ⛔ **完全没有**"卡片在场"那一格（只断文案含 key） | **只有它**：`:544` `Contains(shown,"risk.permission_mode")` | **Errorf**（`:545`）＝红⛔ 中止后续步 | ★**⛔ 无任何分子**（无 windowCount、无 strings.Count） | ★**是**（单读在 `:543`） |

⇒ 派单点名问的"同函数族别处有没有写成 **Fatalf**／写成**到达数分子**的"＝**逐枚答案**：
文案类断言全族都是 `t.Errorf`（本文件 33 处 Errorf），**`t.Fatalf` 全在 helper 与前提格里**（28 处 Fatalf，逐名＝
`:72`/`:88`/`:92`/`:96`/`:100`/`:110`/`:130`/`:161`/`:186`/`:221`/`:244`/`:266`/`:278`/`:296`/`:316` helper 与 await，
加前提格 `:348`/`:359`/`:363`/`:426`/`:439`/`:472`/`:487`/`:516`/`:539`/`:548`）；
**到达数分子只有两处族**：`windowCount()` 四枚（`:379`/`:432`/`:518`/`:679`）与 `strings.Count` 两枚（`:595`/`:668`，
＝票 272 AC#2 的"窗口里恰好一份"那把尺，且 `:668` 特意在 await `state=applied` **之后**才读，注释 `:658-666` 具名）。
★甲3（本票那枚）＝全族**唯一**一枚"既⛔ 有卡片在场格、又⛔ 有分子、还只用 Errorf"的＝三味全缺。

### 乙族＝同文件"审计已到 ⇒ 单读操作员句"的次级同形（⛔ awaitCard，但同一失同步形状）

`:344`（热加载接管句）／`:448`（放宽已过 L2）／`:453`（C26 名单句）／`:478`（放宽没生效）／`:521`（收紧⛔ 该出现卡面）／`:676`（重启档⛔ 该冒充立即句）＝全部 `r.h.out.String()`
**一次读**＋`t.Errorf`（逐名尺＝`git grep -n 'r\.h\.out\.String()' HEAD -- cmd/wisp/config_reload_223_test.go`，
19 命中里 helper／失败转储占 `:73`/`:155`/`:222`/`:267`/`:297`/`:345`/`:401`/`:450`/`:480`/`:522`/`:543`/`:574`/`:594`，
**裸单读断言**＝上面那 6 处 ＋ ★甲族的 `:401`⇄`:543` 两枚）；同文件里**已经存在**等 stdout 的那把尺（`:148 awaitStdout`、`:208 awaitStdoutSince232`）
只用掉 6 处（`:371`/`:623` ＋ `config_receipt_255_test.go` 的 `:633`/`:701`/`:743`/`:797`）。

### 丙族＝包内其余"卡片文案含某串"（尺 C 扫出、⛔ 走 `awaitCard`）

| 枚 | 仪器 | 三格 |
|---|---|---|
| `panel_l1_blocked_220_test.go:240` | `for want := {"[确认 L1 fs.write]", expired}`：stdout⇄stderr **两 stream 任一含即算过** | key＝⛔（含 tool 名）；**Errorf**；分子＝⛔ 无；读＝单读 |
| `panel_pump_test.go:232`＋`:443 consoleCard145` | 正则整段抓卡片，`m == nil` ⇒ **`t.Fatalf`**（＝"在场"那一格独立可失败），再逐项比 level⇄tool⇄reason | key＝⛔；在场格＝**Fatalf**、比对格＝**Errorf**★＝本票要的两格形状在本包**已有先例**；分子＝⛔ |
| `zz315a1probe_test.go:377`（`watchCard` 轮询）＋`:435` | 轮 stdout 直到出现 `[确认 L2`，抓到后只把 `snap.line` 放进**失败消息**里 | ★**有到达等待**（轮询）、⛔ 文案断言＝不在射程；分子＝⛔ |
| `run_mode101_test.go:534`（`TestTicket101ModeSwitchUsesTheRealL2Gate`） | ⛔ 卡片文案，只断 `cards != 1`（`:557` Errorf）＋审计点名 tool | 分子＝**有**（windowCount 差值） |
| `approval_always_201_test.go:209` | `waitForSentence187(...,"the printed rule must name the key and the directory")` | key＝**有**（规则行含 `allowed_dirs += "`）；等＝**有**（轮询两 stream）；红＝Errorf（helper 内 `:694`）；分子＝⛔ |

### 丁族＝settle 现量（AC#2 要用）

尺＝`\.plant\(t,|tc\.plant\(` ⇄ `awaitAudit\(t, "config: HOT-RELOAD state=armed"\)`，逐枚对位（`config_reload_223_test.go`）：
plant 共 **7** 处（`:353` `:394` `:466` `:506` `:541` `:599` `:755`）；plant **前**先等 `state=armed` 的只有 **3** 处
（`:340`→`:353`、`:393`→`:394`、`:753`→`:755`）⇒ ★`:466`／`:506`／`:541`／`:599` **四枚裸种**（本票那枚＝`:541`）。
