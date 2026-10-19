# 343 — `blacklist.go:193` 那句注释点名的 **`aEquiv` 在包里根本没有定义**（只有这句注释自己提到它）＝一句指向不存在符号的产码注释

- Status: **已立，未派**（2026-10-18 14:3x 编排者立；来源＝非实现者终裁腿 `336-v1` 具名上报（`.scratch/wisp/probes/336/v1/logs/NN-blocked-1.log`）＋编排者自己复跑坐实；账＝台账 `A1037`）
- 类型：**注释说实话**（⛔ 任何行为改动、⛔ 任何判据改动）。射程＝`internal/risk/blacklist.go` 的**注释面**（一行）
- ★**与票 337 同族但⛔ 混记**：337 管的是"那句 fail-closed 理由今天没有读者"（**出向语义**、功能级、须机主拍）；本票管的是"一句注释点名一枚不存在的符号"（**盘上字节本身不实**、零行为影响）。两件事的修法、写面、审批层级都不同 ⇒ 分票。

## 现量（编排者 2026-10-18 14:2x 现跑；引用前重跑）

| 事实 | 读数 | 我这把尺 |
|---|---|---|
| 注释在场 | `internal/risk/blacklist.go:193` 逐字 `// provably expanded fails CLOSED (aEquiv / aUnder report such a miss as a` | `git grep -n aEquiv HEAD -- internal/risk` |
| ★**`aEquiv` 全仓只有这一处、且它在注释里** | `git grep -c aEquiv HEAD -- internal/risk`＝**1 枚命中＝该行注释自己**；`git grep -c 'func aEquiv' HEAD -- internal/risk`＝**0** | 两条尺同锚现跑 |
| 成因 | 引入笔＝`f1033e1e`（2026-09-21）；★**⛔ 在本票之前它挂了近一个月没人抓到**＝词面型尺抓不到（它扫的是"符号有没有被用"，⛔ "注释点名的符号存不存在"） | `git log -S aEquiv --format='%h %ad' --date=short -- internal/risk/blacklist.go \| tail -1` 起手复跑 |

## 判据（逐枚可机读）

- [ ] **AC#0 起手现跑**：上面三行逐枚重跑＋★自加一把"同族还有几枚"的尺＝**逐枚点名 `blacklist.go` 注释里出现的每一个包内标识符**（`aEq`／`aUnder`／`aAnyForm`／`anchorForms`／`uncertainAnchorMiss`／`aRelates`／`aRuleReason`…）问"它在包里有定义吗"，命中数具名（⛔ 只交 `aEquiv` 那一枚就走＝票 334 那族"逐枚锚"的教训：一枚抓到＝同族可能还有）。
- [ ] **AC#1 改准那一行注释**：改成**带条件的事实句**＝说清"今天到底是哪两把尺在这种缺位上报 miss"（用**盘上真名**，⛔ 保留不存在的名字、⛔ 为凑话新造一枚符号）。★完成判据＝交件侧 `git grep -n aEquiv HEAD -- internal/risk`＝**0 命中**＋`git diff --numstat`＝**注释行数＝可执行行数变化 0**（尺＝`git diff -U0` 里 `^[+-]` 行逐枚都以 `//` 开头或为注释续行）。
- [ ] **AC#2 ⛔ 动任何可执行字节**：⛔ 新增 `aEquiv` 函数去"让注释变真"（那是**为凑一句话造一枚符号**＝本仓禁手，且会顶红票 336 刚装好的那 8 枚牙）；⛔ 动判级档位／回落支／豁免键语义／`blacklist.go` 任何非注释字节；⛔ 动票 336 的两枚测试件（`anchor_fallback_teeth_336_*.go`）一字。
- [ ] **AC#3 ★"注释点名的符号不存在"这一维要不要上仪器＝只出结论、⛔ 自己上**：本票**只给形状与代价**（例＝一枚 `go/ast` 走查：扫 `internal/risk/**` 注释里的驼峰标识符、对撞包内符号名册、不存在即红；代价＝误伤率——注释常引**别包**符号／历史名／契约号 C1–C32／票号）。★**⛔ 在本票里新造门**（那属"契约邻接"＝要人工批准，先例＝票 262 那枚长度帽只有词面的那一格）；若判"该上"⇒ 具名立新票、⛔ 顺手做。
- [ ] **AC#4 颜色与门禁**：`./internal/risk/` 成对（基线＝起手锚仓外树 ⇄ 交件、同一台面、`-count=1 -timeout 30m`、四数逐枚点名⛔ 混用⛔ 相减）＋★**票 336 那 8 枚牙必须仍 8/8 全红**（改注释⛔ 动它们，但这是本票唯一的"改了这个包"的活 ⇒ 复跑一次是便宜的保险，凭据＝`.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh` 现成台件）；门禁四数只在共享工作树、★**跑前 `unset GOTMPDIR TMPDIR`**（票 311 AC#1 shape B GUARD＝rc=2 射程拒绝，⛔ 读成红）。⛔ `./cmd/wisp`、⛔ `go test ./...`。
- [ ] **AC#5 纪律格**：⛔ 翻任何框（本票六格由非实现者裁）、⛔ 改票名／`Status:`／加 `-done`；写面＝`internal/risk/blacklist.go` 注释面＋`.scratch/wisp/probes/343/r1/**`（短名、⛔ 分节符、逐枚非 0 字节）＋本票末尾**一条**只追加 bullet；`add` 与 `commit` **两处 pathspec 逐枚手写字面路径**（⛔ 从共享 index 反推＝`A1027`／`A1029`）；中文正文⛔ 走 PowerShell（`A1030`）；⛔ push。

## 禁区
`docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／三枚冻结件／`go.mod`／`go.sum`／`.gitignore`／`ci.yml`／D1–D47／C1–C32／R1–R9／D43 转移表；⛔ `frontend/**`／`design/**`；⛔ 回头改票 336／334 的判据文字（都已 `-done`／将 `-done`，要更正只追加）。

**撤销口令**：「撤 343」（⇒ 本票退回未派，注释退回原状）。

## 为什么现在不派
`internal/risk/**` 本波刚被 `336-r1`（写）／`336-v1`（裁）用过、且**票 339 的落地腿正要进这一面**（＝本波释放车道的第一优先）⇒ ★**排在票 339 之后**（同一枚 `blacklist.go`，⛔ 两腿同波动）。⚠ 另：本票与票 339 若同波，两腿都会重跑 `./internal/risk/` 成对四数＝**争用假红风险**（本包在册件 `TestResolvePerCallBudget` 要求安静复量）。

- 343-r1 落地腿交件（2026-10-19 09:1x +08；★⛔ 翻框＝本票六格终量仍 6 未勾／0 已勾，本条只追加在票尾）＝靶子改准：`blacklist.go:193` 那句现在写 `aEq / aUnder delegate to aRelates, which reports such a miss as a hit when anchorMatches found nothing AND uncertainAnchorMiss held`，外加两行补语说清 `anchorMatches` 单看⛔ 报 miss 为 hit（＝票 334 把两把尺拆开的全部理由）。★交件侧 `git grep -n aEquiv HEAD -- internal/risk`＝**0 命中**；完成判据三条尺全给＝`git diff -U0` 的 7 枚 `^[+-]` 行逐枚以 `//` 开头（非注释行＝0）／`--numstat`＝5 加 2 删／★两侧**剥掉注释后的可执行字节 sha256 逐枚相同**（`97796acc…`、265 行）＝可执行行变化 0；⛔ 新增 `aEquiv`、⛔ 动任何非注释字节。
  - `AC#0`＝同族尺我自己加了两把（⛔ 只交 `aEquiv` 那一枚＝票 334 教训）：尺①camelCase（`go/ast`，并把 param／field／local 补进名册——第一版把这形漏了，`bOverrides` 其实是 `Gate` 的实参）＋尺②call-shape（`foo(`/`foo()`，专抓尺①结构性看不见的纯小写名）。★结果＝**`blacklist.go` 注释面命中 2 枚、⛔ 一枚**：`aEquiv`（＝本票、已改）与 `fires()`（`:310`，由票 334 落码腿 `e052b11a` 2026-10-17 引入；包内无 func/method 名 `fires`，真名在他包 `internal/config/unwired.go:54` 那枚字段）；同一对尺扫全包另得 2 枚同族＝`provenance.go:54` `taskScopeID`、`:427` `disposalScope`（全仓 0 枚 IDENT）。★四枚里 3 枚⛔ 在我这一票射程（票面 AC#1 明写"那一行注释"＝193）⇒ 逐枚具名交编排者另立票、⛔ 顺手扩大本票。
  - `AC#3`＝只出结论、⛔ 自己上：**不该上通用仪器**，代价是我量的不是猜的＝同一包内尺①形报 6 枚红、其中 3 枚假阳（`hostPathBox`＝`internal/tools/bridge.go:670` 真 type、`assessorFor`＝同文件真 method、`macOS`＝散文里的平台名）＝**50%**；尺②形报 11 枚红、9 枚假阳＝**82%**（假阳真身逐枚点名在凭据件里：Go builtin `len`、关键字 `func`、`DEFERRED(...)` 标记约定 4 枚、`Nutstore(坚果云)` 厂商名、以及 `CloseScope(id)`＝票 160 **故意删掉**的名字、那句自己写着"Before it…"＝诚实历史引用、词面上与假话无法区分）。根因⛔ 是分词器不行、是"注释对一枚名字做了什么断言"无法从名字恢复⇒ 任何区分规则都退化成手工 allowlist（＝`tools/d22scan/allowlist.txt` 那形、冻结件）。若仍要上，我给的形状＝只对 `^a[A-Z]` 这枚**词表**（全枚＝6）作硬检查：假阳面有界、但⛔ 抓 `fires()`、⛔ 泛化，且仍属契约邻接＝须人工批准⇒ 立新票、⛔ 本票造门。
  - `AC#4`＝`./internal/risk/` 两平台成对四数（各自逐枚点名、⛔ 混用⛔ 相减）：Windows 基线 `247/140/1/1` ⇄ 交件 `247/140/1/1`；Linux 容器基线 `195/115/1/1` ⇄ 交件 `195/115/1/1`；红名册作差 0、`RUN` 名集作差 0（247⇄247、195⇄195）。两侧唯一那枚红＝`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`＝**存量**（静态源码审计、扫 `internal/panel` 的 `Result.Canonical` 读者、读不到本包注释一字；与 339-r1 报里同名对上），⛔ 计入我的新增红；`TestResolvePerCallBudget` 四发全 PASS（1.20/1.19/2.07/1.73s）＝争用假红那条未触发。★票 336 那 8 枚牙＝在**我自己仓外重导的树**重打 `run-ac2-sweep.sh`（⛔ 复用腿的树）＝**8/8 全红、逐名对得上**（MU-1a→git-credentials、1b→wisp-config-under、1c→wisp-config-under+eq、1d→ssh、1e→aws、1f→kube、1g→dpapi-protect、1h→local-credentials），未突变正控 TEETH rc=0／LEG rc=0，11 次还原哈希 match=YES＝我这枚注释字节没顶红任何一枚牙；`shared_worktree_dirty_internal_cmd=0`。门禁四数只在共享树、跑前先 `unset GOTMPDIR TMPDIR TEMP TMP`＝gofumpt `-l` 我的文件 0 枚（新字节 LF、仓内 4 枚存量假脏具名登记⛔ 动）／`scripts/d22scan.sh` rc=0／`check-path-length-budget.sh` rc=0／`go vet ./internal/risk/` rc=0。
  - 凭据与提交＝`.scratch/wisp/probes/343/r1/logs/`＋`sh/linux-343.sh`：`31d4a476`＝注释面那一枚 `.go`（唯一一枚）／`2b6f8a9a`＝Linux 台件／`360a6af6`＝读数正文件＋其原始 stdout 同批入库（★成对读数与产生它的机器输出⛔ 拆两笔、拆了复跑对不上；行数帽只管我撰写的正文件＝01 件 105 行在对帽内，⛔ 为凑帽截断 `go test` 自己的输出＝篡改读数，同票 262 长度门"rewriting a reading is how a record stops being evidence"那条口径）。⛔ push（推送只由编排者做）。
