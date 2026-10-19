# 346-a1 件00 = CI 取日志尺 ⇄ 归一尺 ⇄ 五份日志四数（本腿现跑）

★本腿 0 枚 `go test`（第九波公共块＝⛔ 抢在飞两枚腿的 CPU）＝本件每一行来自 **gh 直取的 CI 日志**
或**盘上已入库件**；⛔ 引用他人四数（旧日志本腿自己归一后重量）。

## 一、取日志的硬尺（照派单，且本腿自己先栽过一次、具名）
1. `gh run view 39244425666 --repo CarlosShao/wisp --json jobs` → 取**数字** `databaseId`：
   test-core ＝ 118032776497｜test-windows ＝ 118032776659｜lint ＝ 118032776906。
2. `gh run view 39244425666 --repo CarlosShao/wisp --log-failed --job <数字ID>`（★⛔ 喂作业名）。
   - **本腿首发带错 flag `--no-pager` ⇒ rc=1 ＋ 三枚日志全 0 字节**＝"看起来⛔ 红"那一形本腿复现过一发，
     根因＝flag、⛔ 作业名；去掉后 rc=0／1489719B⇄1240950B。⛔ 掩盖、⛔ 拿那三枚 0 字节当读数。
   - **lint（118032776906）三次重试均 0 字节**＝`api.github.com` 连接超时（stderr 319B 逐字
     `dial tcp …: connectex: A connection attempt failed…`）⇒ 本件判 **未知**，★⛔ 读成"lint 没有红"。
     lint ⛔ 本票必答射程（本票射程＝两枚仪器＋三味成因），故弃该格并具名。

## 二、归一尺（★这一步是四数能不能取的前提，本腿现量其必要性）
GitHub `--log-failed` 每行 = `<job名><TAB><步骤名><TAB>2026-…Z<payload>` ⇒ payload **永⛔ 在行首**；
直接跑派单那把顶格尺 ⇒ **0 命中**（本腿实测：`'^--- FAIL'` 于原始件＝0，而 `'^--- (PASS|FAIL|SKIP)'` 于
归一件＝7/13/14 枚；正控＝原始件里子串 `FAIL` 命中数百＝尺⛔ 失明）。
⇒ 归一尺 ＝ 剥两个 TAB 字段 ＋ 剥首个 `Z` 之后的时间戳 ＋ 剥一个前导空格
（脚本 `D:\tmp\wisp\20261019\346a1\norm.ps1`；字节区间尺＝同目录 `bytespan.ps1`，绝对偏移、含该行）。
⇒ 这正是派单 `A1054`③ 那句"旧日志的 FAIL 行带时间戳前缀我按行首锚定⇒名册 0 枚"的机制，本腿先补这一环。

## 三、五份归一日志 ⇄ 四数（顶格锚定、只顶层、⛔ 混用⛔ 相减）
| 归一文件（`D:\tmp\wisp\20261019\346a1\`，仓外、只建不删） | 来源 | 字节 | sha256 前 12 | RUN | PASS_TOP | FAIL_TOP | SKIP |
|---|---|---|---|---|---|---|---|
| `core.norm.log` | 39244425666 / 118032776497（test-core） | 686396 | 760D9541D294 | 1833 | 1239 | 7 | 2 |
| `win.norm.log` | 39244425666 / 118032776659（test-windows） | 683457 | 77C254BD6A8D | 1041 | 692 | 5 | 2 |
| `old39124.norm.log` | 39124122121（旧 run **全量**、含多作业） | 1389927 | CA9D1A225F6E | 2991 | 2037 | 13 | 6 |
| `old39101.norm.log` | 39101060753（全量） | 1132935 | 7C032836115B | 2755 | 1858 | 14 | 4 |
| `old39071.norm.log` | 39071959151（全量） | 1475365 | DE5382017C37 | 2849 | 1917 | 7 | 2 |

原始件（gh 直出、未归一）＝1489719B（CAAF2E9A8643）／1240950B（A25AC9537009）。
⚠ 射程＝**单作业文件 ⇄ 全量文件混在一张表里**：`old*` 三枚是全量（多作业拼接）⇒ 四数比单作业大发，如实记、
★⛔ 与 `core.norm.log`／`win.norm.log` 相减或互换。
★`go test` 把正文印在 `--- FAIL` 之**前** ⇒ 本件所有"红句"一律按 `=== RUN` 段切取段内正文（⛔ 事后 grep），
＝派单那条纪律本腿逐枚照做（见件02／件03 的行区间，均为段内）。

## 四、本票三枚名的落位速查（尺＝`^(=== RUN|--- (PASS|FAIL|SKIP)): <名>`）
| 名 | core.norm | win.norm | old39124 | old39101 | old39071 | 本机整包 |
|---|---|---|---|---|---|---|
| `TestPOSIX120SealWindowLeafSwapIsClosed` | **PASS** L5676 | ⛔ 现身 | **FAIL** L10856 | **PASS** L5796 | PASS L12282 | ⛔ 现身（0 命中） |
| `TestAC1POSIX…322r1` | **PASS** L5430 | ⛔ 现身 | **FAIL** L10610 | **PASS** L5544 | ⛔ 现身（文件还⛔ 在树里） | ⛔ 现身 |
| `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` | ⛔ 现身 | **FAIL** L1844 | PASS L331 | PASS L6480 | PASS L3304 | PASS L352 |
⇒ 前两枚的射程／机制＝件04 第二节（`//go:build !windows`）；第三枚＝件03 第三节（cmd/wisp CLI 步骤）。
