# 352-r1 · AC#3 纪律格（名册＋写面自证）＋门禁四数

★本腿⛔ 给任何 AC 盖章＝`AC#0..AC#3` 四格全归非实现者裁；这枚文件只交**纪律自证**的盘上尺。
尺件＝`sh/06-discipline.sh`（R1..R12）＋`sh/07-codepoint-ruler.pl`＋`sh/08-sanitize-codepoint.sh`＋`sh/09-utf8-ruler.pl`＋`sh/10-utf8-selfcheck.sh`；
逐字输出＝`logs/discipline.txt` ⇄ `logs/sanitize-codepoint.txt` ⇄ `logs/utf8-selfcheck.txt` ⇄ `logs/gates.txt` ⇄ `logs/gofumpt-stock.txt`。

## R1⇄R2 名册尺＝两把独立尺对撞（★锚定行首，⛔ 假设其后任何字符，⛔ 拿 `A..B` 区间数）

| 尺 | 形 | 枚数 |
|---|---|---|
| 锚定 `--grep='^票 352 落地腿 352-r1'` | 行首锚 | **5 笔（本文件提交时＝6）** |
| 单次 `%H%x09%s` 导出＋内存内 `index($2,P)==0` | 独立实现 | **同 5 笔** |
| 非锚定 `--grep='票 352 落地腿 352-r1'` | 会吃正文引用 | 5 笔（此波⛔ 别腿在正文引用本前缀＝两数相同，⛔ 据此以为锚定多余） |

主题起头逐字＝「票 352 落地腿 352-r1」＝派单给的那枚精确串，★⛔ 换形（⛔ 英文、⛔ 全角变体、⛔ 中途改形状）。
起手锚＝本腿第一笔 `git log -1` 现认 `bf509bb837aa1328617320608be38e4bca5a5d35`（分支 `dev`）。

## R3⇄R5⇄R12 写面自证

- 本腿全部笔的路径并集＝**65 枚**；`OUTSIDE_ALLOWED_WRITEFACE=0`；★产码 `.go` 枚数＝**1**＝只 `internal/winsec/winsec.go`。
- 禁区逐枚现验命中本腿笔数＝`pinnedleaf_342_windows.go` **0**／`winsec_other.go` **0**／`winsec_windows.go` **0**／`pinnedleaf_342_unix.go` **0**／`pinnedleaf_342_fallback.go` **0**／`internal/risk` **0**／`internal/config` **0**／`frontend` **0**／`design` **0**／`cmd/wisp` **0**／`internal/agent` **0**／`internal/panel` **0**。
- ★跨腿写面互斥守住：`cmd/wisp/**`＝`351-r2`、`internal/agent/approval/**`＝`230-r1`，本腿⛔ 碰其任何字节（R12 `CROSSLEG_HITS=0`）。
- 那一行那笔的 `-U0` 全貌（★只动那一枚格式串、`NET_DIFFSTAT=1 insertions 1 deletions`）：

```
@@ -431 +431 @@ func wrapPath(path string, err error) error {
-	return fmt.Errorf("%w: %s: %v", ErrNotSealable, path, err)
+	return fmt.Errorf("%w: %s: %w", ErrNotSealable, path, err)
```

- ⛔ 新增接缝／导出名／改签名＝那一笔净变动只有那一行，函数签名行⛔ 在 diff 里；⛔ 碰 `winsec.go:118-140` 自陈段（同一笔 1/1 行＝只 `:431`）。
- ⛔ push（R11）＝`UNPUSHED_COMMITS_TOTAL` 含别腿与编排者的笔、`UNPUSHED_COMMITS_MINE=5`（→6）＝本腿笔全在本地 `dev`。
- ⛔ `git add -A`／⛔ `add .`／⛔ `--amend`／⛔ `reset`／⛔ `rebase`／⛔ `stash`／⛔ `checkout`／⛔ `clean`＝`add`⇄`commit` 两处都逐枚字面 pathspec；还原一律仓外 pristine 拷贝＋每发跑后 sha256 对回（`RESTORE=OK`，⛔ `git cat-file … > …` 之外⛔ 任何 checkout 形）。

## R6 票 352 框尺终量（★现量，⛔ 采信任何旧叙述——含编排者派单里的数）

| 尺 | 值 |
|---|---|
| `grep -cE '^- \[x\] '`（已勾） | **0** |
| `grep -cE '^- \[ \] '`（未勾） | **4** |
| `grep -cE '^- \[[ x]\] '`（总框数＝自校：两者之和必＝此数） | **4** |
| 格名 | `AC#0 AC#1 AC#2 AC#3`（分母＝4 格） |
| 正控＝同两把尺数票 348 | `[x]`=6 ⇄ `[ ]`=1 ⇄ ANY=7 ★尺灵 |

⇒ ★**本腿一枚框没翻**＝票 352 终量 0 勾／4 未勾；票名／`Status:`／`-done` 一字未动；票尾 numstat＝**纯追加**（`3 增 0 删`＝空行＋空行＋一条 bullet，`git diff -U0 | grep -c '^-[^-]'`=0）。

⚠★本腿自抓第三枚（改在入库前）＝那把框尺**第一版是钝的**：写作 `'^  - \[x\]'`（带两枚前导空格）⇒恒 0 命中，而真值恰好也是 0＝"报 0 命中"与"尺⛔ 咬"撞脸＝正是公共块"报 0 命中前先跑一枚已知存在的正控"那一味，本腿自己犯过一次；现改行首锚＋加"总框数＝CHECKED+UNCHECKED"自校＋票 348 正控＝上面那张表。

## R8⇄R9 输出预算与字符尺

- 被禁字面字符（码位 U+00A7）＝★写面命中 **0**；正控＝同一把尺数 `docs/reports/HANDOVER.md`＝**304** 在场＝尺灵。
  净化过程照实＝先前入库的六枚 raw 读数件里，抄来的 `AGENTS.md` 分节引用带了那枚字符（每枚 2 处）⇒`sh/08` 换成码位名；★**纯度尺**＝每枚删除行换码位名后逐字＝其对应增加行（`PAIRED_ROWS=2/2` ×6）＝⛔ 改任何读数；★正控 U+00B7（间隔号）＝命中 0＝那枚⛔ 被动。
- 零字节件＝**0**（一枚 `internal/config` 容器档落地侧空名册从 0 字节补成 `# RED_NAMES=0` 说明行＝规矩"⛔ 0 字节文件"）。
- `.out` 扩展名＝**0**（全用 `.txt`／`.md`／`.tsv`／`.pl`／`.sh`）；正文件行数＝AC#0 72／AC#1 65／AC#2 92／本件 ≤160；★枚数帽口径＝公共块"每落一枚文件就 commit"优先，"4–8 笔"当参考值⇒本腿落到 8 笔并具名登记。
- ★⛔ 并行跑改文件的脚本＝本腿全部读数串行；★⛔ `-overlay`、突变⛔ 打跟踪文件（本票⛔ 需要突变＝同一把尺在两棵树上读出不同分类＝尺自身有牙，具名写进 `readings-ac1-landed.md`）。

## 门禁四数（★只在共享工作树；跑前 `unset GOTMPDIR TMPDIR TEMP TMP`；退码＝脚本里 `echo RC=$?` 紧跟裸命令）

逐字＝`logs/gates.txt`。

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| 1 | `/d/work/base/gopath/bin/gofumpt.exe -l internal/winsec/winsec.go` | **0** | 空列出＝干净（只验自己这枚、⛔ `-w`、⛔ 扫别腿文件） |
| 2 | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations`（internal/=603、cmd/=143 Go files，注释与 `_test.go` 都在射程内） |
| 3 | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`（over-budget=60＝roster=60＝⛔ 新增） |
| 4 | `go vet ./internal/winsec/`（本机档＝判红绿那枚） | **0** | 输出 0 字节 |
| 4b | `GOOS=linux go vet ./internal/winsec/` ／ `GOOS=darwin go vet …` | **0** ／ **0** | 各 0 字节。★交叉档⛔ 与本机档成对判红绿＝这里只登记"那两档能编译＋vet 无声" |

★存量具名登记（⛔ 我动、⛔ 顺手重排）＝`gofumpt -l internal cmd tools` rc=0 列出 **4** 枚既有件＝
`internal/agent/approval/pending_read.go`（＝`230-r1` 写面）／`internal/agent/tools.go`／`internal/risk/provenance.go`／`cmd/wisp/models.go`（＝`351-r2` 编译面）＝本机 CRLF 存量形；★我那枚新字节＝入库 blob CR 计数 **0**（`sh/09` 尺另证：那六枚 raw 里 CR 只来自本机 `cmd /c mklink` 自己吐的 CRLF＝原始读数照实保留、⛔ 净化）。

## `winsec.go` 各笔 blob 串（供编排者判 `351-r2` 读数可比性）

| 笔 | commit | `internal/winsec/winsec.go` blob |
|---|---|---|
| 起手锚 | `bf509bb8` | `07300bfffc3e5b14c017029d6f84fb5a7ede3db8` |
| 笔① | `0fb399c3` | `07300bfffc3e5b14c017029d6f84fb5a7ede3db8`（＝⛔ .go 字节） |
| ★笔② | `ecffb93b` | **`f3c1b68ee60d8d5d42b816e47633d31f1ba191f5`**（＝那一行落地） |
| 笔③ | `aacde632` | `f3c1b68e…` |
| 笔④ | `c16ff47e` | `f3c1b68e…` |
| 笔⑤ | `06aa2abf` | `f3c1b68e…` |
| 笔⑥⑦⑧（本文件⇄票尾） | — | `f3c1b68e…`（★⛔ .go 字节⇒对 `351-r2` 读数⛔ 变） |

⇒ `351-r2` 若量到点名 `internal/winsec/**` 的红＝具名上报、⛔ 自己动字节、⛔ 据此判自己失败（派单跨腿耦合①的约定）；`07300bff` 起⇄`f3c1b68e` 之后各自成段＝两侧读数⛔ 可混比。
