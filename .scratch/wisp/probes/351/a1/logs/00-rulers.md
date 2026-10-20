# 351-a1 件00 = 尺名／射程／正控 ⇄ 四数复尺（★本腿 0 枚 `go test`，凭据＝盘上既有日志＋CI 日志＋`git grep HEAD`）

起手锚＝本腿第一笔 `git log -1` 现认＝`14e5d8395adaf54bec14199f6858382b57fa49dd`（2026-10-20 09:58:30 +08，分支 `dev`）。
⛔ 用编排者写的号。取源码事实一律 `git grep -n <re> HEAD -- <path>`／`git cat-file blob HEAD:<path>`；
blob 副本落仓外 `D:\tmp\wisp\20261020\351a1\`（⛔ 解到仓内、⛔ 盘符根、只建不删）。

## 一、五把尺的定义＋逐尺正控（★凡报枚数带尺名）

| 尺 | 逐字命令 | 射程 | 本腿读数 | 正控 |
|---|---|---|---|---|
| 尺 A（函数定义） | `git grep -nE '^func TestTicket' HEAD -- cmd/wisp` | 只含 `_test.go`（非 `_test.go` 命中＝**0**，本腿现验） | **86 枚**／21 枚文件 | 锚定⇄非锚定（`func TestTicket`）**同＝86**⇒本包无缩进定义、两根一致（⛔ 此处差异＝无差异，如实记） |
| 尺 B（awaitCard 命中） | `git grep -n 'awaitCard' HEAD -- cmd/wisp internal` | 含注释行 | **5 命中**＝1 注释（`:251`）＋1 定义（`:252`）＋**3 处调用**（`:395`／`:467`／`:542`） | 全仓加射程（⛔ 限 `cmd/wisp internal`）＝24 命中，多出全在 `.scratch`⇄`docs` 的**散文引用**⇒尺 B 的射程必须限包，⛔ 拿 24 当仪器数 |
| 尺 C（卡片文案字面） | `git grep -nE '\[确认 ' HEAD -- 'cmd/wisp/*_test.go'` | 只 `_test.go` | **5 命中** | 正控＝`config_reload_223_test.go:402` 已知在场 |
| 尺 D（"点名"族红句） | `git grep -nE 'does not name\|name the key\|Contains\(shown' HEAD -- 'cmd/wisp/*_test.go'` | 只 `_test.go` | **56 命中**／31 枚文件 | 正控＝`:545` 在场 |
| 尺 E（等了再读的那把尺） | `git grep -nE 'awaitStdout(Since232)?\(' HEAD -- 'cmd/wisp/*_test.go'` | 只 `_test.go` | **8 命中**＝2 枚定义＋6 处调用 | 正控＝`:623` 在场；★`:543`⇄`:401` **不在**此列＝那两处是未等的单读 |

`awaitCard` 定义体现读＝`config_reload_223_test.go:252-271`：轮 `r.rt.liveCards.pending()`、`time.NewTicker(20ms)`、
期限 `reloadCaseBudget`（`:35`＝40s）、超时 `t.Fatalf`。⇒ 它等的是 **gate 侧注册表**，⛔ 等 stdout。

## 二、四数复尺（★复用 `346-a1` 那把归一尺，⛔ 重造）

归一尺＝`D:\tmp\wisp\20261019\346a1\norm.ps1`（剥两个 TAB 字段＋剥首个 `Z` 后时间戳＋剥一个前导空格）。
★"不剥＝假全绿"本腿独立复现一发：原始件 `ci39244425666-118032776659.log`（1240950B）上顶格尺＝
`^--- FAIL` **0**、`^=== RUN` **0**，而同件子串 `FAIL` 命中 **341** ⇒ 顶格尺在未剥的件上确实失明（＝假全绿那一形）。

本腿在自己脚下重跑四数（尺＝归一件顶格锚定、只顶层、⛔ 混用⛔ 相减）：

| 件 | 字节 | RUN | PASS_TOP | FAIL_TOP | SKIP | 与 346-a1 留档 |
|---|---|---|---|---|---|---|
| `win.norm.log` | 683457 | 1041 | 692 | 5 | 2 | **逐数同** |
| `core.norm.log` | 686396 | 1833 | 1239 | 7 | 2 | **逐数同** |
| `old39124.norm.log` | 1389927 | 2991 | 2037 | 13 | 6 | **逐数同** |
| `old39101.norm.log` | 1132935 | 2755 | 1858 | 14 | 4 | **逐数同** |
| `old39071.norm.log` | 1475365 | 2849 | 1917 | 7 | 2 | **逐数同** |

★取 CI 的硬尺本腿也复跑了一发（⛔ 抄）：`gh run view 39244425666 --json jobs` → 数字 `databaseId`＝
test-core `118032776497`／test-windows `118032776659`／lint `118032776906`；`headSha=6f808a83c63a...`＝票面所写 `6f808a83`。
★喂作业名⇄带错 flag＝静默 0 字节那形＝`A1054`②，本腿⛔ 重犯（只用数字 ID）。
⚠ `lint`（118032776906）本腿未取＝346-a1 留档为 `api.github.com` 连接超时、判**未知**；本票射程⛔ 含 lint ⇒ 弃该格并具名。

## 三、本机⇄盘上件（⛔ 本腿复跑、只读盘）

- 本机整包＝`D:\tmp\wisp\20261019\orch227-b\full.log`（1544211B；四数在册＝2317/1529/29/8＝编排者留档）。
  ★该发射程⛔ 是整包终态名册（票面 `A1056`① 已写死；本腿独立同证＝其红名册末枚之后无包行）。
- blob 同一性（本腿现跑 `git rev-parse`）＝`cmd/wisp/config_reload_223_test.go` 在
  `a5080612` ⇄ `6f808a83` ⇄ `HEAD` **三处同串** `f1cac524055aaac000ced549296833ae99eb5a71`；
  产码 `cmd/wisp/config_reload.go` 亦 `990d148b...` 两处同串。
- `git log --format='%h %s' a5080612..6f808a83 -- cmd/wisp` ＝ **0 枚**（本腿现跑；`6f808a83` 是 HEAD 祖先＝`rc=0`）。
- 票 351 框尺现量＝`grep -cE '^- \[[ x]\]'` 于 `git show HEAD:<票面>`：总 **5** 格（`AC#0`..`AC#4`）、
  已勾 **0**、未勾 **5**（尺＝行首锚 `^- [ ]`，逐字前 7 字符分组＝5×`- [ ] *`）。
