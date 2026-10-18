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
