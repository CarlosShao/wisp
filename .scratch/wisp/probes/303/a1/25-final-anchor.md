# 303-a1 `25-final-anchor.md` — 终态锚（与起手同形）＋逐笔名册＋差集

## 1. 终态现量（同一条命令里取，date 在第一段）

| 尺 | 逐字 | 起手（`00-anchor.md`） | 终态（本件，2026-10-10 17:2x） |
|---|---|---|---|
| 时刻 | `date` | `Sat Oct 10 16:54:59 CST 2026` | `Sat Oct 10 17:21:01 CST 2026` |
| HEAD | `git log -1 --format='%H %ad' --date=iso` | `dce0f133…` `16:51:56 +0800` | `48ee655e…` `17:14:06 +0800`（＝本腿第 4 笔；⛔ 别人的笔混进来——本腿 4 笔全在 `.scratch/wisp/probes/303/a1/**`） |
| 产码射程 | `git status --porcelain -- cmd internal scripts .github docs \| wc -l` | **0 行** | **0 行** ⇒ 全程⛔ 动产码／测试码／档／台账／脚本，工作树≡HEAD |
| `wisp.exe` | `tasklist //FI 'IMAGENAME eq wisp.exe' \| tail -n +3 \| grep -ciw wisp.exe` | 0 | **0**（⛔ 起过本体） |
| `balldebug.exe` | 同上换名 | 0 | **0**（⛔ 起过本体） |
| `msedgewebview2.exe` | 同上换名 | 未量（起手没列这一枚，是**本件的欠尺**） | **24**＝每发真窗测试留下的 WebView2 运行时进程（宿主 `go test` 二进制都已退出；这些 PID ⛔ 本腿能具名 `-Id` 的自起对象，⛔ 杀——同机别的程序也在用这一枚运行时） |

## 2. 本腿跑过的 go 命令枚数自报（派单要求）

`go test` **18 发**，逐枚：

| # | 台面 | 尺 | 用途 | 件 |
|---|---|---|---|---|
| 1 | 母仓 | `-count=1 -timeout 420s -v -run 'TestAC14…' ./cmd/wisp/` | `AC#0` ① 终点乙（HEAD） | `10-repro-head-count1.txt` |
| 2 | 母仓 | 同上 `-count=3` | `AC#0` ③ 复现率三色 | `11-repro-head-count3.txt` |
| 3 | clone `cc315261` | `-count=1 -timeout 300s`（台件 `bisect-step.sh`） | `AC#0` ① 起点（正控） | `16-clone-cc315261-start-pass.txt` |
| 4 | clone `bcd0a543` | 同上 | `AC#0` ① 终点甲（负控） | `17-clone-bcd0a543-end-fail.txt` |
| 5–14 | clone（bisect 12 步里的第 3..12 步） | 同上 | `AC#1` bisect（GOOD 4／BAD 8／INVALID **0**） | `22-bisect-index.tsv`·`21-bisect-log.txt` |
| 15 | clone `fb2fb802` | `-count=1 -timeout 420s -v -run 'TestAC14…\|TestAC13…'`（台件 `targeted-pair.sh`） | `AC#1` (a) 定向红发 | `18-pair-a-culprit-red.txt` |
| 16 | clone `fb2fb802^`＝`f718e9b6` | 同上 | `AC#1` (b) 单撤那一笔的绿发 | `19-pair-b-parent-green.txt` |
| 17 | clone（那一笔的树，删新增测试文件） | 同最小集 | `AC#1` §6 加料① | `23-extra-culprit-minus-new-testfile.txt` |
| 18 | clone（那一笔的树，产码半边换父发） | 同最小集 | `AC#1` §6 加料②（INVALID＝build failed） | `24-extra-culprit-prod-file-from-parent-buildfailed.txt` |

⛔ 其它 go 编译面：⛔ `go build`／⛔ `go vet`／⛔ `go run`／⛔ `gofmt -w`；`go version` 1 次（起手锚，非编译面）。
⛔ `t.Skip` 造绿、⛔ 改断言、⛔ 放宽判据、⛔ 碰 `thresholds.go`／`allowlist.txt`／`frontend/src/**`／`design/**`／三枚冻结件。
⛔ push（0 次）；⛔ `git add -A`／`.`；每笔 commit 带显式 pathspec 且写在 `$( … )` 之外。

## 3. 逐笔名册（尺＝`git show --name-only --format= <我的每笔>`）

| 我的笔 | 名册 |
|---|---|
| `b3f9c6d7` | `00-anchor.md`·`cm-1.txt`·`tmp-baseline-grep.txt` |
| `0edce3b6` | `10-repro-head-count1.txt`·`11-repro-head-count3.txt`·`bisect-step.sh`·`cm-2.txt` |
| `f2ac16dd` | `14-suspect-classification.md`·`cm-3.txt` |
| `48ee655e` | `13-ac0-repro.md`·`cm-4.txt` |
| 本件那一笔（第 5 笔，号见回报） | `20-ac1-bisect.md`·`25-final-anchor.md`·`16-`·`17-`·`18-`·`19-`·`21-`·`22-`·`23-`·`24-` 十枚读数件·`targeted-pair.sh`·`cm-5.txt`·`msg-anchor.txt`·改动的 `13-ac0-repro.md` |

**与授权名册的差集**＝授权面是 `probes/303/a1/**` 那一个目录 ⇒ 五笔名册**全部落在这一枚目录里**，差集＝**∅**。
（`msg-anchor.txt`＝起手时写歪的一枚草稿信文件，⛔ 删（`issues/README` 规则 8"临时件只建不删"），随第 5 笔入库；
`13-ac0-repro.md` 的那一处改动⛔ 改读数，只是把 §② 末尾"这一发欠着"那句换成已经量到的件名。）
bisect 用的 clone（`$HOME/wisp-303-bisect`）与它的 12＋2 步原文（`$HOME/wisp-303-steps/`）＝**仓外临时台面**，⛔ 入库（件体已逐枚拷进本目录），⛔ 删除。

## 4. 本格交回的两格（⛔ 判语，判语归非实现者）

- `AC#0`＝`13-ac0-repro.md`（①②③④ 各带尺＋件）
- `AC#1`＝`20-ac1-bisect.md`（12 步色表／first bad＝`fb2fb802f75a0e3eeacad488f1adc6f064e29f85`／那一枚动的是**产码＋测试码两样**／定向两发齐全／加料两发＋欠读数具名）

## 5. 追加更正（只追加⛔ 改上面任何一行读数）

- §3 那一行里我写"十枚读数件"＝**枚数写错**：第 5 笔（`86ed3592cd6408083eb910c73e849f74804f9299`）实际名册 **14 枚**，
  其中 clone 侧读数件 **8 枚**（`16-`·`17-`·`18-`·`19-`·`21-`·`22-`·`23-`·`24-`）＋两件正文（`20-`·`25-`）
  ＋两枚台件（`targeted-pair.sh`·改动的 `13-ac0-repro.md`）＋三枚信草稿（`cm-5.txt`·`msg-anchor.txt`）。
  尺＝`git show --name-only --format= 86ed3592 | wc -l` → **14**（本件落笔时现量）；差集仍＝**∅**（全在 `probes/303/a1/**`）。
- 同一处尺的口径补一句：`13-ac0-repro.md` 在第 5 笔里是**修改**（M）⛔ 新增，改的只有 §② 末句那一条"这一发欠着"→"已量到，件名如下"；⛔ 动任何读数行。

