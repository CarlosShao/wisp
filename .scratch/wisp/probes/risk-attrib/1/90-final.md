# 90-final——现态／⛔ 做到的名册／越界尺／与本单不符处

台面 HEAD（起手现量）＝`ec87076f1c6a90c45facf3558635582fb0ead989`（`ec87076f`），分支 `dev`。
本腿落笔四笔（⛔ push）：`63bc6335`（第 1 笔＝`00-anchor.md`＋`logs/gate-00.txt`）／`f4fd71bc`（第 2 笔＝表①＋表②＋`logs/**` 全部仪器件）／第 3 笔＝`30-routing.md`＋本件／第 4 笔＝`logs/boundary-name-only.txt`（逐笔 `git show --name-only --format=<hash>` 名册，⛔ 区间尺；它自己的号请 `git log -1 --format=%H` 现量）。

## 1. 交件四张表＋枚数（**枚数＝名册行数**）

| 件 | 内容 | 枚数自核 |
|---|---|---|
| `10-ownership-roster.md` | 表①：12 行名册，逐枚 包＋文件＋内容锚＋声明行（行号一律标快照） | 12 行＝CI 的 12 枚；尺＝`grep -rn "func <名>(" --include=*_test.go .` 全树，每枚恰 1 命中，**0 枚多命中／0 枚零命中 ⇒ ⛔ "排除 N 枚"那一格（N＝0）** |
| `20-mechanism-per-red.md` | 表②：逐枚判定＋本腿自己的颜色凭据＋标签 | **甲 7 ＋ 乙 5 ＋ 丙 0 ＋ 丁 0 ＝ 12** |
| `30-routing.md` | 表③：甲族**挂票 72**／乙族**挂票 252**／仪器形状三条**登记**；⛔ 立任何新票、⛔ 塞票 302 | 7＋5＋3 条（第三条是"三形合一格"，⛔ 是 3 枚用例） |
| `logs/`（16 枚件） | 闸门、两发颜色原始字节、`csplit` 逐枚切片、8.3 别名种植与读数 | — |

★**这一句必须带**（`A817` §5 的硬要求）：CI 上这 12 枚**同一台 runner 改前改后名集合逐字相同（12↔12）**，⇒ 与最近那批改动无关；本腿的归因⛔ 是"谁带来的"，是"它们为什么红"。

## 2. 本腿今天判死了什么（⛔ 用推理冒充读数）

- **12/12 全部在本机判死**，凭据＝`T02` 那一发（`env TMPDIR='C:\Users\swq\tmp\RISKAT~1\RUNNER~1' … go test ./internal/risk/ -count=1 -v -run '<12 枚>'` ⇒ `RUN=12 FAIL=12 PASS=0 SKIP=0`，件尾逐字 `FAIL github.com/CarlosShao/wisp/internal/risk 0.249s`）。
  ⇒ **改写 `A817` §5 的前提**：那一格写的是"待归因／归口另待一枚普查腿定射程"，而它隐含的"可能要真 runner"这一支**已被本腿否证**——runner 的成因是 `%TEMP%` 的 **8.3 拼写**，仓外种一枚别名＋env 覆盖就能在机主这台电脑上复现，⛔ 碰任何全局配置、⛔ 动任何 `_test.go`。
- `A817` §3 那两格"机制未证"⇒ **本腿填掉**：
  - "2 枚红句不含 `RUNNER~1`"⛔ 是旁证缺口，是那两枚的 `t.Fatalf` 只 `%+v` 打了 decision struct、⛔ 打输入路径（`anchor_spelling:238`／`junction:292`，逐字红句在 `T02`）；两枚同因＝override 键的拼写两侧⛔ 同层。
  - "5 枚 under-profile fallback 家族机制未证"＝**同一处缺陷的第二枚调用点**（`syncdirs.go:88` 的 `home` 只过 `normPath`、候选侧 `:408 resolveTarget` 过 C26 ⇒ `:429 isUnder(cand, s.home)` 一边折长一边留短），且**方向是 fail-open**（`{Sync:false …}`），且**有真生产调用者**（`provenance.go:928`）。
- 反向那 1 枚（`TestC21TableColourRowsMatchTokensCSS`）**判⛔ 是同一族**：它 `os.ReadFile` 工作树资产字节（`internal/ball/tokens_table_test.go:1384`，常量 `:1364`），成因面＝`design/**` 资产／放置那一族，⛔ 是拼法那一族。**它⛔ 在我车道**（`internal/ball`），所以本腿⛔ 跑、CI 侧成因⛔ 判 ⇒ 归编排者或机主那一发。

## 3. 只能由编排者／机主判的（本腿⛔ 越过去）

1. **立不立／扩不扩射程**：表③ 给的是"该挂票 72（甲 7）／该挂票 252（乙 5）／三形登记"三选一加理由加射程，**裁权在编排者**。⛔ 本腿动票面、⛔ 翻框、⛔ 判"本格闭合"（裁决者≠实现者）。
2. **变异那一把尺（丙）**：⛔ 改产码／⛔ 改 `_test.go` 授权 ⇒ "判据换成反形它也不响"本腿**没打**。那是**欠的读数**，先例＝票 115 `AC#4` 的 M1/M2/M3，归**非实现者**那一发。
3. **#4 那枚 `t.Skip` 逃逸口**（`junction:129`，票 18 授权）：要不要按票 115 定式改成 `t.Fatal`＝改判据形状 ⇒ 归编排者／机主。
4. **票 115 名下⛔ 一张 `docs/evidence/s1/115-*.md`** 这件事（票 230 的账，票面自陈撞车）：本腿引的是票 115 的**做法**与那张非实现者读数表 `docs/evidence/s1/ci-step-readings-2026-09-22.md`；要不要补一张 115 名下的表⛔ 归本腿，已在 `A817`／票 230 那条链上，⛔ 本腿再记一次。
5. `cold bring-up 2747 ms 超 D32 预算`／票 305 那两枚文档交接形＝**⛔ 本腿射程**（`cmd/wisp` 面），本腿一枚没碰。

## 4. ⛔ 做到的名册（逐条具名，⛔ 掩）

1. ⛔ 改任何产码／`_test.go`、⛔ 新建 `.go`、⛔ "顺手修一下让它绿"——**做到了**（工作树里 `internal/**` 本腿 0 字节改动）。
2. ⛔ 跑 `./cmd/wisp/`、⛔ `-tags winlive`、⛔ `./...` 整包、⛔ `d22scan.sh`、⛔ `gofumpt`／`gofmt -l` 名册类——**做到了**（本腿只跑过 2 发 `go test ./internal/risk/`＋1 次 `grep` 级只读；⛔ `go build`／`go vet`／`go list`／`go env`，那一栏允许但没用到，⛔ 欠账）。
3. ⛔ 写进仓的临时件、⛔ 改全局/用户级配置、⛔ 动 git config——**做到了**（种植在 `C:/Users/swq/tmp/riskattrib1/`，只建不删；仓内临时件只落自家 `logs/`）。
4. ⛔ 真窗／真机／GUI、⛔ 读 `%APPDATA%\wisp-dev\secrets\`、⛔ 读机主 config 的值、⛔ 写 `frontend/src/**`／`design/**`——**做到了**（`design/**` 只做了只读 `git status`）。
5. ⛔ push、⛔ `git add -A`／`.`、⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`、⛔ 仓内建 worktree——**做到了**。
6. ⛔ 放宽任何断言、⛔ 造 `--- SKIP`、⛔ 动 SLO 阈值／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／D43 转移表／C1–C32／D1–D47／三枚冻结件——**做到了**（两发都是只读求值，⛔ 改一字）。
7. ⛔ 把 MB 级 CI 原始日志入库或读进上下文——**做到了**（名册只引已入库的 `probes/303/orch/r9-ci-after-red-roster.txt`）。
8. **⛔ 做到的（违例自报）**：`T02` 那一条我写了 `… > 件 2>&1; echo "rc=$rc"` 却**漏了 `rc=$?` 那一拍** ⇒ 那一发的退码格是空的。颜色证据由件尾逐字 `FAIL` ＋ `--- FAIL=12/PASS=0` 承担；`T01` 的 `rc=0` 同发现量。**第 1 笔 commit 我第一次试了 `git commit -- <未跟踪路径>`**（rc=1、pathspec 未匹配），改成显式两枚路径 `git add` 再 commit（⛔ `-A`／`.`）。
9. **⛔ 做到的（第二枚，落笔纪律）**：`00-anchor.md` 里我把 raw log 的**件名**写成与盘上不同的那一枚（`gate-2026-10-11.txt` vs 盘上 `gate-00.txt`），⛔ 删、⛔ 改名，按盘上真名走 ⇒ 那一件的 §2 描述现在带着⛔ 存在过的名字。**引用时以盘上为准。**

## 5. 与这一单（派单文字）不符处——以盘上为准，逐条具名

1. **⛔ 给 HEAD 号**：派单让我"派发时刻 HEAD 现量"，本腿现量＝`ec87076f`；`A817` 正文里那几枚（`cf46c24a`／`05db4bc6`／`bcd0a543`／`cc315261`）是别的发的 headSha，本腿⛔ 拿来当台面。
2. **`A817` §3"2 枚同文件同族"——同文件那一半过期**：`TestAListWinsWhereBothTablesHit` 在 `pathresolver_anchor_spelling_windows_test.go:204`，`TestBListDefaultDenyAndOverride` 在 `pathresolver_junction_windows_test.go:272`。**同包、不同文件、不同主题**。"同族"那一半本腿认（同一 override 键成因）。⛔ 改台账，报回即可。
3. **派单／`A817` 把 12 枚统称"`internal/risk` 那一族"＝准确的（这一枚默认⛔ 需要更正）**，但**没写编译面差**：7 枚在 `//go:build windows` 文件里、**5 枚（`syncdirs_test.go`）无 tag、POSIX 也求值** ⇒ 直接影响修法形状（⛔ 能用"加 windows tag"摘掉它们，会撞票 301 的包级 guard D）。已写进表③ §4。
4. **派单说"5 枚是 `syncdirs_test.go` 的 under-profile fallback 家族"**＝枚数与文件对，但**"under-profile fallback"这半句只覆盖其中 4 枚的判据**：`syncdirs_test.go:162` 那一枚红在**根成员判定**（fixture 根下的写没被认成 sync），⛔ 是 fallback 那一支。两支点的是**同一处不对称**（`e.canon` 与 `s.home` 都没过解析器），所以本腿仍归成一族，但名册里按判据分开写了。
5. `A817` §3 那五个行号（`:162/:213/:233/:355/:402`）在本腿 `T02` 里**逐字复现** ⇒ 这一条⛔ 腐烂，照给。派单提醒"行号一律当快照"，本腿在表① 里对每枚都另给了**内容锚**（文件里那句 `t.Fatalf`/`t.Errorf` 的原文）而不只给号。
