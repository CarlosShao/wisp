# 342-v3 票 348 专属读数摘要（＝票 342 同批交件的删除支；全量逐字日志主存在 `.scratch/wisp/probes/342/v3/logs/`，本根放本票分子＋尺件副本）

腿＝`342-v3`（非实现者裁决腿）。起手锚现认 `677d4837`；实现腿＝`342-r2`（10 笔、末笔 `bfcac9ed`）。台面＝四棵自建仓外导出树（基线 `fde1df02` ⇄ 落地 `677d4837` ⇄ 落地⛔ 本腿字节 `3abbb639` ⇄ 去牙专用 `treeT`）＋两枚本腿自写仓外台架；临时根 `%LOCALAPPDATA%\Temp\v3-342v3`（短形）；`go1.27.1 windows/amd64`。⛔ 抄 `342-r2` 读数、⛔ 抄编排者 `A1059`。

## 件二（稳态祖先 junction）＝票 348 `AC#2` 那一格的两栏

尺＝`go test -v -count=1 -run` 定向四枚（两树×各 3 发＝6 发）；逐字行＝`probes/342/v3/logs/04..09`。

| 栏 | 落地（`treeL`，r1/r2/r3） | 基线（`treeB`，r1/r2/r3） |
|---|---|---|
| REFUSAL | 三发×逐字同一条 `risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions`，`refused=true / isUnresolved=false / isNotSealable=false`，leg＝the installed C26 resolver（3 shots） | 同一条逐字同、`refused=true`，leg 同上（3 shots）＝**两棵树同号＝缺世界⛔ 一台一面** |
| DAMAGE | `the foreign object … =present (size 72 …)` ＋ `the caller's own object … =present (size 42 …)` ＝`truncated=0 gone=0` | 同两枚 `present` ＝`truncated=0 gone=0` |
| ★阳性对照 | `same object=true / by-name remove err=<nil> / control name is now=absent` ×3 | 同三行 ×3＝★那枚 0 是读数、⛔ 钝尺 |

⇒ `AC#2` 的 Windows"修前分子>0"**照旧取⛔ 到**＝`342-r2` 照实报的那句本腿复现；本腿判**部分成立**（同 `A1042` 票 339 形＝登记续体⛔ 结案）。

## ★"缺"钉到成因（本腿自打，`342-r2` 与前裁都⛔ 有）

1. seam 位置〔算的，`git cat-file blob` 现读两棵树〕＝`var applyDescriptor`（`winsec_windows.go:38`）**基线 blob `fde1df02` 里同一行同字**、`sealHandle(f) → applyDescriptor(f.Name(), false)`＝`winsec_windows.go:413` **两棵树同一枚号** ⇒ 钩子只在**叶子句柄活着**那一瞬间；票 348 那一支的删除在 `_ = f.Close()` **之后**（基线 `winsec.go:158-160` ⇄ 落地 `:175-177` 逐字）。
2. 平台锁两态〔量的，`logs/01`，3 发〕＝叶子句柄活着时 `rename(ancestor)`＝`Access is denied`×3；句柄关掉且⛔ 钉态＝`rename(data)`／`rename(sub)` 两枚 `err=<nil>`×3＝★祖先能动的瞬间⛔ seam。
3. 四枚可能接缝瞬间全拒〔量的，`logs/02`＝H1/H2/H3/H4 ×3 发；正控 H0 空手＝`err=<nil>`〕＝★"祖先在 open 之后翻"那枚世界⛔ 是缺钩子、是缺内核事实 ⇒ 第二道批准本腿签**⛔ 该批**。
4. 跨进程攻击者〔量的，`logs/03`＝holder 握父钉×2 发〕＝`rename(data)`＝`Access is denied`／`rename(sub)`＝`in use by another process`／★同发 `remove(caller's name)=<nil>` ＋ `rename(mine -> caller's name)=<nil>` ＋ holder 的产码那一步 `Root.Remove(base)=<nil>` ＝**攻击者那枚文件的唯一名字被删**。⇒ 祖先进⛔ 去、**名字**进得去＝本腿签**单立一票**（⛔ 属本票判据族；名与判据形状见 `docs/evidence/s1/342-v3-verdict.md` 末节）。

## 两枚去牙（本腿落点⛔ 同腿报，`logs/15`+`logs/16`+`logs/18`）

| 发 | 落点（内容锚，⛔ 行号） | 落刀前/后锚命中 | 读数 | 还原 |
|---|---|---|---|---|
| A | `unlinkPinnedLeaf` 里 `_ = p.root.Remove(base)` → 注释 | 1→0 ／ 新锚 1 | 件一 `entered=3` 照在、`unnamed objects left with no name 3→0` ＋回落红×3 `AC#5's half-made file survived the refusal`；`go vet` rc=0 | `B9DC39AEA15670E8F24689FA8D5F8157F2A1FFE310996203670C37DCF4D5643B` MATCH=True |
| B | `openPinnedLeaf` 的 `p.root.OpenFile(base, flags, perm)` → 按名 `os.OpenFile` | 1→0 ／ 新锚 1 | 件一 `worlds entered 3→0`、damage 0＝★窗口确由句柄带 `FILE_SHARE_DELETE` 打开；`go vet` rc=0 | 同上 MATCH=True |

⇒ 两栏各自成数、⛔ 一把两种形状都放行；损伤归因到产码那一步（⛔ 仪器自造）。

## 成对四数＋逐名红册（`^=== RUN`⇄`^--- PASS`⇄`^--- FAIL`⇄`^--- SKIP`，⛔ 混用⛔ 相减）

| 台 | RUN/PASS_TOP/FAIL_TOP/SKIP | 红名册（逐名） |
|---|---|---|
| 落地 r1 | `119/67/1/0` | `TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1` |
| 落地 r2 | `119/66/2/0` | 上面那枚 ＋ `TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1` |
| 基线 r1 | `114/64/1/0` | `TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2`＝★本腿叶子件打在基线产码上的有牙红（⛔ 存量） |
| 基线 r2 | `114/64/1/0` 逐字同号 | 同名 |
| 落地⛔ 本腿字节 r1 | `115/62/2/0` | 那两枚 `342r1` 名＝★逐名证明⛔ 来自 `342-r2` 的字节 |

★`342-r2` 报的"落地 3 枚红"＝本腿复现到**同名册、枚数为摆动量**（1⇄2，`treeLandMinus` 到 2）＝按读数纪律（竞态⛔ 要求逐字相同）裁：**方向与归属复现、枚数⛔ 当通则**；★落地侧"新增红"＝**0 枚**。★唯一合法存量红 `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` 住 `internal/risk`、本包红名册未出现。`--- SKIP`=0（五轮逐枚）；`mklink /J` 建不出＝0 枚。

## 门禁（只在共享工作树，跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）＝`logs/22`（同串副本）

`gofumpt -l`＝本腿写面 0 枚 `.go`（★具名：指 `.md` 的 rc=2 属工具误用、⛔ 读成红）；`d22scan` **rc=0 clean**；`check-path-length-budget` **rc=0 VERDICT GREEN**（over-budget 60＝roster 60）；`go vet ./internal/winsec/` 本机 **rc=0** ＋ `GOOS=linux`/`GOOS=darwin` 各 **rc=0**（交叉档⛔ 与本机成对判红绿）。★⛔ `./cmd/wisp` 整包、⛔ `./...`＝派单令、射程外。

## 名册与越界（两把锚并集＋独立整腿形）＝`logs/20`/`logs/21`/`logs/25`

锚定行首名册 **10 笔** ⇄ `--grep='^…'` **10 笔** ⇄ 并集（按 hash）**10**；★"续体"形锚定后 **0 笔**（非锚定 3 笔＝编排者自己的笔＝`A1045` 复现）；★**区间⛔ 是名册**＝`3abbb639..bfcac9ed` 实含 **37 笔**（含 `296-r2`/`296-v2`/`351-a1`/编排者）⇒ 改用逐枚 name-only 并集＝**43 枚路径**／非 `_test.go` 的 `.go`＝**0 枚**／写面之外＝**0 枚**；`winsec.go` 10 笔逐枚 `07300bff`、`:118-140` 逐枚同 sha（★号对账见 `logs/23`/`logs/24`＝`sha256` 那把才复现 `e8a2f33dea8cdeab`）。

## 框尺（须＝起手，本腿一枚没翻）

票 348＝`[x]`**6**／`[ ]`**1**（未勾＝`AC#2`）⇄ 票 342＝`[x]`**5**／`[ ]`**2**（未勾＝`AC#1`/`AC#4`）。尺＝`grep '^- \[x\] \*\*AC'`／`'^- \[ \] \*\*AC'`，射程＝两枚票面文件；票名／`Status:`／`-done` 一字未动。
