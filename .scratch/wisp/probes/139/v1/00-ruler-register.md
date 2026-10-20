# 票 139 终裁腿 139-v1 · 尺册（逐把尺名＋逐字命令形＋射程＋量的还是算的＋正控）

起手锚＝`a19d5bf2`（2026-10-20 22:35:43 ＋08）；本腿全程⛔ `go` 命令（`go test`／`go build`／`go vet`
／`go list`／`d22scan`（含 `go run`）一律未跑＝零）；临时件落 `D:\tmp\wisp\20261020\139v1\`（只建不删）。

## A. blob 三把（AC#2 落点在场性，一律对**入库 blob** 量、⛔ 工作树＝本机钝尺）

| 尺名 | 逐字命令形 | 射程 | 读数 | 量/算 |
|---|---|---|---|---|
| R1 日志调用行 | `git cat-file blob HEAD:internal/agent/compress.go \| grep -nE '\.Info\(\|\.Warn\(\|\.Debug\(\|\.Error\('` | 单文件 blob | **1 行**＝`:240` 逐字 `c.log().Info("agent: history compressed", attrs...)` | 量的 |
| R1b 去点变体 | 同尺换形 `grep -nE 'Info\(\|Warn\(\|Debug\(\|Error\('` | 同上 | **1 行**＝同一枚 `:240`（两把独立形同号＝票面第 0 节那把的"1"不是遮出来的） | 量的 |
| R1c 正控（同尺打别包） | R1 打在 `HEAD:internal/agent/loop.go` | 单文件 blob | 逐四味拆开＝`\.Info\(` 1／`\.Warn\(` 7／`\.Debug\(` 0／`\.Error\(` 4 ⇒ 并集 **12 行**（去点形 **14**）＝尺灵⛔ 复述任何腿报的"14" | 量的 |
| R2 字段名册 | 逐枚 `grep -c '"<key>"'` 对 R1 同一 blob | 单文件 blob | `tokens_before`1 ⇄ `tokens_after`1 ⇄ `threshold`1 ⇄ `msgs_before`1 ⇄ `msgs_after`1 ⇄ `compressed_msgs`1 ⇄ `kept_raw_rounds`1 ⇄ `history_changed`**2** ⇄ `task_id`**0**／`taskID`**0**（真名＝`"task"`，`:238 attrs = append(attrs, "task", task)`） | 量的 |
| R2b 取值侧 | `sed -n '205,245p'` 同一 blob 逐行读 | 单文件 blob | attrs 八枚值全取自 `rep.*`／`len(hist)`／`len(out)`／`c.b.HistoryCompressTokens`；`history_changed`＝`rep.TokensAfter != rep.TokensBefore \|\| len(out) != len(hist)`（⛔ 复用 `rep.Ran`） | 量的（读） |

## B. 接线与出线（AC#3 载具）

| 尺名 | 命令形 | 读数 |
|---|---|---|
| R4 接线三点 | `git cat-file blob HEAD:internal/agent/loop.go \| grep -n 'l\.comp\.Need\|withTraceTask(ctx, taskID)\|res\.Compression = rep\|Compression CompressionReport'` | `:100 Compression CompressionReport` ⇄ `:406 if l.comp.Need(hist) {` ⇄ `:412 nh, rep, err := l.comp.Compress(withTraceTask(ctx, taskID), hist)` ⇄ `:417 res.Compression = rep`（量的） |
| R4b 通道闭合 | 同尺对 `compress.go` blob 打 `withTraceTask\|traceTaskID` | 定义 `:152`／读侧 `:160`／消费 `:237`（量的）；`Result` 结构体＝`loop.go:82` 起、注释在 `:99`、字段在 `:100` |
| R5 名册（用例逐名） | `git cat-file blob HEAD:internal/agent/compress_trace_test.go \| grep -n '^func Test'` | **7 枚**逐名见判语件第 2 件（量的）；同一把打在 `1acdd032`＝4 ⇄ `ff550f38`＝5 ⇄ `b23c7f77`＝7 ⇄ `HEAD`＝7（来路尺＝算的＋量的混合：每枚各跑一把） |
| R5b 放水检 | `git grep -n 't\.Skip' HEAD -- internal/agent` | 该件 **0 枚**；全包仅 2 枚＝`approval/ticket84_no_owner_test.go:224`（存量）＋`approval/zz_201_a1_...:31`（注释）＝⛔ 本票面（量的，正控＝同尺能命中非零） |

## C. 旧批⇄新批写面（⛔ 区间归因、逐枚列名）

| 尺名 | 命令形 | 读数 |
|---|---|---|
| R6 写面 | `git show --name-only --format='' <sha>` 逐枚 | `1acdd032`＝`compress.go`＋`compress_trace_test.go`＋`harness_test.go`＋`loop.go`（`.go` **4**）｜`ff550f38`＝`compress_trace_test.go`（**1**）｜`b23c7f77`＝`compress.go`＋`compress_trace_test.go`＋`loop.go`（**3**） |
| R6b `.go` 计数尺正控 | `git show --name-only --format='' 1acdd032 \| grep -c '\.go$'` | **4**＝尺灵（判语第 4 件用它背书新批的 0） |
| R7 未漂尺 | `git diff --name-only ec826c63..da221145 -- '*.go'` ⇄ `git rev-parse <rev>:<path>` 三枚对撞 | 区间 `.go` 差异＝**5 枚**（全 `internal/winsec`，⛔ agent 三包）；`compress.go`⇄`loop.go`⇄`compress_trace_test.go` 在 `ec826c63`⇄`HEAD` **blob 同串**＝`4fcd9a32`／`ae8c25c7`／`2357dca8`（算的） |
| R8 pristine 哈希 | `git cat-file blob HEAD:internal/agent/compress.go \| sha256sum` | `16e153d8b8b4bd78d8f3dc050ce62eb425a6d1652cec88c5207be1727c7c2f50`＝与复验腿"两发还原 MATCH"逐字同串（量的） |
| R8b 落刀锚唯一性 | `grep -cF` 三形对同一 blob | `c.log().Info("agent: history compressed", attrs...)`＝**1** ⇄ `if rep.Ran {`＝**1** ⇄ `if true {`＝**0** ⇒ M1⇄M2 落点⛔ 同、两把锚都是唯一命中（量的） |
| R9 底片四数 | `grep -c '^=== RUN'`⇄`'^--- PASS'`⇄`'^--- FAIL'`⇄`'^--- SKIP'` 逐枚 | 射程＝`.scratch/wisp/probes/139/r1/raw/`：`base-v` **212/168/0/1** ⇄ `head-v` **212/168/0/1** ⇄ `m1b-v` **114/86/4/0** ⇄ `m2-v` **114/88/2/0** ⇄ `m1-v`（假读数那发）**114/90/0/0**；`panic:` 逐枚 **0**；`grep -c '^=== FAIL'` 未用＝四数不混用不相减 |
| R9b 红名⇄红句 | `grep '^--- FAIL'` ⇄ `grep -n 'compress_trace_test.go[0-9]*:'` | M1 四枚逐名＝`BooksCountsOnSuccess`／`SurvivesTheLoopWiring`／`CarriesTheOwningTaskID`／`NeverInventsATaskID`；四枚红句原文行号＝`:212`⇄`:334`⇄`:500`⇄`:549`（M2 两枚＝两枚 Silent 反向）；⛔ 见 `01-discrepancies.txt` 第②条串位 |
| R9c 逐名表 | `wc -l` ⇄ `diff` 对 `raw/roster-base.txt`⇄`roster-head.txt` | **169 ⇄ 169**、`diff` rc＝**0**（量的） |

## D. 名册三把尺（本腿自己现跑、⛔ 抄腿的数）

尺1 锚定行首＝`git log --all --format='%s' \| grep -c '^票 139 落地腿 139-r1'`＝**11**
尺2 非锚定子串＝同射程 `grep -c '票 139 落地腿 139-r1'`＝**11**
尺3 内存 `index($2,P)==1`＝射程 `git log --all --format='%H%x09%s'`＝**11**
⇒ **三把同号**；正控＝`^票 321 终裁腿 321-v1` 同尺 **6**＝尺灵。
逐枚列名（量的）＝`eb623bee`① ⇄ `b0191a0e`② ⇄ `6f397118`③ ⇄ `259f8087`④ ⇄ `2340b8de`⑤ ⇄
`b3030f7e`⑥ ⇄ `c3e76454`⑦ ⇄ `f0cd6d46`⑧ ⇄ `c1356acf`⑨ ⇄ `bdbbc8e3`⑩ ＋编排者派单存档 `ec826c63`。
★写面自证（R10＝尺9 那把 `.go` 计数逐枚打）＝十一枚 **`.go` 全 0**、越界 **0**（正控＝同尺打 `1acdd032` 给 4）。

## E. 三味总尺＋本腿写面自证

| 尺 | 命令形 | 读数 |
|---|---|---|
| 票 139 框量 | `git cat-file blob HEAD:<票面>` ⇄ `grep -c '^- \[ \]'` ⇄ `grep -c '^- \[x\]'` | 未勾 **2**（`:53 AC#2` ⇄ `:59 AC#3`）／已勾 **2**（`:46 AC#1` ⇄ `:64 AC#4`）＝**本腿一枚没翻**；行数 **130**／字节 **22644**／blob CR **0**／码位 U+00A7 **21** 枚（全部早于本腿） |
| 唯一性尺 | `git ls-tree HEAD --name-only -- .scratch/wisp/issues/ \| grep -c 'issues/139-'` | **1**（`153-` 亦 1；票 139 ⛔ 加 `-done`＝本腿⛔ 动） |
| 追加额度 | `git show --numstat --format='' f0cd6d46` | `6 0 <票面>`＝纯追加、删除列 0＝票尾那一枚额度已被复验腿占＝本腿⛔ 补第二条 bullet |
| 台账账号存在性 | `grep -c '^## A1071' docs/reports/pending-and-issues.md` 逐枚 | `A1071`**1** ⇄ `A1072`**1** ⇄ `A1073`**1** ⇄ `A1069`**1** ⇄ `A271`**1** ⇄ `A563`**1**＝本判语引用的号**全部真存在**（正控＝全文件 `grep -c '^## '`＝**928**）；⚠ `Q43`／`Q46` 裸形＝0＝那两枚在台账里的形制带连字符（`Q-43`），本判语⛔ 引裸号当凭据 |
| 码位尺（本腿写面） | `LC_ALL=C grep -o "$(printf '\302\247')" \| wc -l` | 判语件落笔第一发＝**9** 枚→本腿自己改掉＝现量 **0**；`00`⇄`01`⇄`02` 三枚件＝**0**（正控＝同尺打票 139 票面＝**21**＝尺灵）＝按 `A1051` 收紧后的射程＝腿对自身证据件的卫生尺 |
| CR 尺（对入库 blob／对盘上件） | `tr -cd '\r' \| wc -c` | 判语件 **0**；`00`⇄`01`⇄`02` **0**（⛔ 用 `grep -c` 数字节＝它数行） |
| 编码尺 | `iconv -f UTF-8 -t UTF-8` | 判语件 rc＝**0**；`03-final-roster.txt` 见 `01-discrepancies.txt` 第③条＝rc **1**（那是复验腿的件、⛔ 本腿写面） |
| tracked 脏（收线发） | `git status --porcelain \| sed -n 's/^ M //p'` 后逐枚 `stat -c '%y'`（量的，射程＝整树） | **1 枚**＝`.scratch/wisp/issues/115-seal-notices-…-caller-spelling-done.md` mtime `2026-10-20 20:24:44`＝出表腿 `115-v1` 面、⛔ 本腿字节（本腿写面在 status 里＝新落的两枚，⛔ 别人件） |
| 在飞现量 | `find internal cmd -name '*.go' -newermt '2026-10-20 22:00'` ⇄ `tasklist //FI "IMAGENAME eq go.exe"` | **0 枚** ⇄ 无 `go.exe`＝本腿⛔ 据此断言 `gw-1` 的字节状态（只 stat⛔ 读内容）；现势跑色那味⛔ 属本腿尺＝**归 `gw-1`** |

## F. 仓根冻结件完整性（A1069③ 那一形要求的自证）

`go.mod`＝**1026 B**／mtime `2026-10-20 19:50:08` ⇄ `go.sum`＝**6237 B**／`2026-10-01 09:47:42` ⇄
`.gitignore`＝**1404 B**／`2026-09-27 11:16:34`（量的）＝三枚 mtime 全部**早于本腿起手（22:33）**＝零触碰；
本腿全程⛔ `go mod edit`／`tidy`／`get`、⛔ 任何输出重定向落仓根件；要读模块信息一律 `git cat-file blob HEAD:go.mod`。
