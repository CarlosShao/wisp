# 161-r1 证据件：D22 门禁盲区普查（AC#1 读数 + AC#5 契约轴 + 门禁）

- 程：票 161 的 **r1 写码位**（派单 `.scratch/wisp/dispatches/2026-09-26-213x-impl-161-r1-gate-blindspot-census.md`）
- 射程：**只做 AC#1（量盲区）＋ AC#5（契约轴零字节）＋ §6 门禁**。AC#2/AC#3/AC#4/AC#6 一格未做（派单 §1 明令）。
- 生产字节改动：**0**（`tools/d22scan/**` 一枚字节未写；`.github/workflows/**` 未写；见 §6 取证）。
- 本程写面：`.scratch/wisp/probes/161/r1/**` ＋ 本文件 ＋ 票 161 的 Progress log（只追加）。

时刻与档位约定：每条读数前都现跑了命令；**〔我本轮现跑过〕** 是本文件默认档位。

---

## 0. step 0 五件（命令原文 + 输出）

### 0.1 `date`

```
$ date
Sat Sep 26 21:34:16 CST 2026
$ date +%Y-%m-%d' '%H:%M' %z'
2026-09-26 21:34 +0800
```

（后续每次落时间前重新取过：21:45、21:52、22:0x，见各节。）

### 0.2 锚点

```
$ git rev-parse HEAD
fbe12c7176bd42d4389f594cfa1be9c19ffe1ac5
$ git rev-parse --short HEAD
fbe12c7
$ git branch --show-current
dev
```

⇒ 本程一切"改前/改后"比较的锚＝**`fbe12c7`**（不是我抄来的任何号）。

### 0.3 此刻盘上哪些不是我的东西

```
$ git status --porcelain | head -30        # 全量行数 48
 M .scratch/wisp/issues/160-scope-open-returns-a-handle-carrying-its-own-closer.md
 M .scratch/wisp/issues/161-gates-need-their-own-gates-self-tests-in-ci-plus-is-the-approval-door-ever-rung.md
 M .scratch/wisp/probes/152/my152.py
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
 D design/index.html
 D design/screens/approval.html
 D design/screens/ball.html
 D design/screens/chat.html
 D design/screens/config.html
 D design/screens/cost.html
 D design/screens/firstrun.html
 D design/screens/palette.html
 D design/screens/privacy.html
 D design/screens/security.html
 D design/screens/states.html
 D design/screens/tasks.html
 M docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md
 M docs/reports/pending-and-issues.md
?? .scratch/wisp/dispatches/2026-09-26-213x-impl-161-r1-gate-blindspot-census.md
?? .scratch/wisp/dispatches/2026-09-26-213x-readonly-161-r2-doorbell-census.md
?? .scratch/wisp/issues/169-a-single-decorative-glyph-in-frontend-makes-the-repo-wide-gate-red-and-kills-every-ci-step-after-it.md
?? .scratch/wisp/probes/139/accept-r1/
?? .scratch/wisp/probes/152/overlay-probe1-on-samppost.json

$ git status --porcelain -- frontend design | wc -l
31
$ git status --porcelain | wc -l
48
```

派单 §3 点名的现场**逐条对上**：`design/**` 有 owner 的未提交删除（16 枚 `D` 行）＋未跟踪新件；
`probes/152/my152.py` 是 ` M`；`docs/evidence/s1/152-…accept-r1.md` 有未提交改动。
⇒ 本程**没提交、没还原、没补完**它们；每次提交都带显式 pathspec（见 §6.5）。

⚠ 与本派单 §3 不同的一处：`docs/reports/pending-and-issues.md` 此刻也是 ` M`（台账有别人的未提交行）。
那更是"绝不 `git add -A`"的现场理由，本程一次都没碰它。

### 0.4 票面 AC#1 / AC#5 两格（连行号抄，行号是读数）

票面 = `.scratch/wisp/issues/161-gates-need-their-own-gates-self-tests-in-ci-plus-is-the-approval-door-ever-rung.md`

- **AC#1**（票面 **17–20 行**）：
  - `17` 判据正文：「**AC#1 先量"今天的门会不会瞎"**：对九条禁令逐枚造一发**已知违规**样本（只放进
    `.scratch/wisp/probes/161/**`，**不进树**），现跑现量"这一枚今天点不点得到"。⚠ 拿不到的那几枚
    **照实登记成盲区**，不许放宽判据、也不许造一枚永远不响的检（本仓否过两次恒真判据）。」
  - `18` 编排者更正段（同页保留原句）：真值＝`tools/d22scan/main.go:5` 起的 `// Bans` 名单，编到 `ban #8`。
  - `19` 「另有 `unparseable` 那枚**是发现类型、不是禁令**，别把它算进分母、也别漏掉它那一发样本。」
  - `20` 「**本格的射程以你现读 `// Bans` 那一段为准，不按"九"也不按"八"写死**」。
- **AC#5**（票面 **29 行**，一整行）：
  「**AC#5 契约轴零字节**：`docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`internal/panel/**`、
  `internal/agent/approval/**`、`thresholds.go`、golden、`allowlist.txt`、`scripts/slo-check.ps1`、
  `frontend/**`、`design/**`。名册按 commit 号集合算（`--no-walk` 必须有）。」

### 0.5 禁令名册现读（派单给的 awk 窗口不够长）

派单 §0.5 的命令原文：`awk 'NR>=5 && NR<=40' tools/d22scan/main.go`
⇒ **这一窗口把 ban #8 截断在第 40 行**（ban #8 的正文是 `main.go:40–49`），
而派单 §6 自己写着"截断的输出永远不许当全表"。**名册以全文为准**：

```
$ awk 'NR==5 || NR==49 || NR==50 {print NR"|"$0}' tools/d22scan/main.go
5|// Bans (PLAN.md D22 fourth-round additions), production scope internal/ +
49|//	                      which is generated from the scope list.
50|//
```

⇒ `// Bans` 段 = **`main.go:5–49`**，里面有编号的禁令 **8 枚**（tab 起首的 `1`–`8`），
外加第 51–53 行的 `allowlist.txt` 抑制说明。枚数与票面 18 行的更正一致（**不是"九条"**）。

**现读名册（编号／标签／作用域／注释豁免）**：

| 编号 | 标签（main.go 里的名字） | 作用域 | 实现层 | 注释豁免？ |
|---|---|---|---|---|
| 1 | bare-goroutine | `internal/`+`cmd/` 非 test 的 .go | AST `*ast.GoStmt`（闭包与具名调用都算，R16#1） | 天然豁免（AST 不看注释） |
| 2 | pathresolver-bypass | 同上 | AST SelectorExpr，接收者标识符**必须叫** `filepath`，方法名 `Clean`/`Abs` | 天然豁免 |
| 3 | plaintext-key | 同上 | AST `ValueSpec` + `AssignStmt`（左值是 Ident）＋值形状闸 | 天然豁免 |
| 4 | wallclock-timeout | 同上 | **逐行正则**（`main.go:759–770`），`//` 开头的行跳过 | 有：整行注释跳过；行内注释不跳 |
| 5 | mirror-hash | 同上 | **逐行正则**（`main.go:771–774`）：同一行同时出现 mirror 与 hash 词 | 有：整行注释跳过 |
| 6 | panel-approval | `frontend/` **每一个文件**（walkText，无后缀白名单） | 逐行正则 `approval\.decide`（`main.go:825–830`） | **无**（不剥注释） |
| 7 | internal-artifact-tool | `internal/tools/` 的 .go（非 test） | 逐行正则 `"(spill\|internal[._-][a-z0-9_.-]+)"`（`main.go:832–837`） | **无**（不剥注释） |
| 8 | emoji | `design/`＋`frontend/`(everyFile)＋`internal/`(goOnly)＋`cmd/`(goOnly)（`emojiScopes`，`main.go:557–564`） | 字符类 `emojiRe`（`main.go:164`） | **有**：Q-46(c) 注释豁免、字符串从严（`commentRangesFor`） |
| — | `unparseable` | 同 #1–#5 的 Go 作用域 | `parser.ParseFile` 失败即记一条（`main.go:684–689`） | 不适用 |

**`unparseable` 算不算禁令？本程判：不算分母，但照样取样。**
理由（读代码，不是抄话）：票面 19 行明写它是发现类型；它在 `main.go` 的 `// Bans` 名单里**没有编号**；
但它和编号禁令共用同一个退码通道（`s.add` → `verdict()` 的 `len(s.findings)>0` → rc=1），
所以"它响不响"与"它会不会把别的禁令弄瞎"都是门的行为 ⇒ §3 里给它一格（两发：`unp-bad`、`unp-hides-bans`）。

---

## 1. 前提核对：派单 §9 点名的四条，两条与盘上不符

| 派单写死的前提 | 盘上读数 | 结论 |
|---|---|---|
| §4：`tools/d22scan` 今天本来有两枚红（`TestScannerSelfScanOfRealRepoIsGreen`、`TestRealRepoLedgerIsHonest`），四数 `PASS=28 FAIL=2 SKIP=0`、`=== RUN=70`，`scripts/d22scan.sh` rc=1、全仓 1 枚 finding | 锚 `fbe12c7` 上实测：**PASS=30 FAIL=0 SKIP=0、`=== RUN=70`**，两枚点名用例**都在并都绿**（pre 日志 23 行、138 行）；`sh scripts/d22scan.sh` **rc=0**，全仓 **0 枚 finding** | **不成立（已被 611ae8b 取代）** |
| §5.1：`// Bans` 段是禁令的真相源 | 成立：8 枚编号禁令 + `unparseable` 发现类型，与 `emojiScopes`/`declaredScopes` 的实际作用域一致 | 成立（一处小偏差，见 §1.2） |
| §6：四条门禁命令可跑 | **`go vet ./tools/d22scan/` 跑不了**（跨模块）；**`gofumpt -l .` 非空**（3 枚不是我产物的脏件） | 两条需按原文修正 |
| §2/§3：兄弟程 161-r2 只写 `probes/161/r2/**`、不写票面 | 盘上 `probes/161/` 下**只有我的 `r1/`**；票面 Progress log 的追加只有我一程在做；r2 的写面此刻未见 | 未见反例（我无法证明它"不会写"，只能报"此刻没写"） |

### 1.1 §4 那两枚红为什么不红了（这是本程最该先报的一条）

编排者量基线的时刻是 **21:23**；修它的提交是 **21:26**，我的锚是 **21:28**：

```
$ git merge-base --is-ancestor 611ae8b HEAD && echo yes
yes
$ git log --no-walk --format='%h %ad %s' --date=format:'%Y-%m-%d %H:%M' 611ae8b fbe12c7
fbe12c71 2026-09-26 21:28 docs(前端 log §8 事故链自陈): …
611ae8b9 2026-09-26 21:26 fix(票 169 甲案): right-rail.tsx:108 假终端文案里的 U+2713 换成 OK——…；d22scan 复跑 rc=0、全仓 0 finding
$ awk 'NR==108' frontend/src/components/harness/right-rail.tsx
  { kind: "ok", text: "OK built in 826ms" },
```

⇒ 派单 §4 那句"如果你看到的不是这两枚红：停手报回"我**照办了：这一条报回，并继续做我做得动的部分**
（§9 给的就是这个处置）。**没有**放宽、跳过或修改任何断言，**没有**动 `frontend/**`。

口径怎么落（**只往严格方向动，判据文字我自己一个字没改**）：
- 派单口径「除那两枚之外零枚新增红」在我这一版树上等价于「**零枚红**」，因为那两枚已经不红。
- 基线**不是抄编排者的号**，是**我在锚 `fbe12c7` 上自己现跑现存的**（§5 的 pre/post 两本日志）。
- 反向哨兵：如果 post 里出现任何一枚红，或与 pre 的名册两向 `comm` 差集非空 ⇒ 停手报回。实测差集**两向皆空**。

### 1.2 §5.1 的小偏差（不改判据，只记档）

`// Bans` 是"禁令有哪些"的真相源，但**不是"禁令的机器作用域"的真相源**：
第 6 枚的文字写 `scope: frontend/`，同一份 `frontend/` 上 ban #8 走的是另一条腿
（`walkText(..., "panel-approval", …)` vs `emojiScopes` 里的 `frontend/, everyFile:true`），
两条腿的文件集合**今天相等**（实跑都是 85），但这是**两份声明的巧合**，不是一处真相源 ⇒
名册表里我把"作用域"记成**代码里实际的走法**，并注明与注释文字的一致性。

---

## 2. 本程**没有**测什么（按"漏了它谁会先被骗"排序，每条给可复算的最小闭合动作）

这一节是本件最该被读的部分。每条都是**真话："我没量"**，不是修辞。

**N1 — 门在真实树上的召回率：一枚样本都没进过树。**
我全部读数来自 `-root <假根>`，扫的是我自己造的 12–16 枚小文件。今天真树
`internal/`+`cmd/` 有 228 枚生产 .go、`frontend/` 85 枚文本、`internal/tools/` 18 枚，
我没量过**任何一枚真文件是不是恰好落在盲区里**（§3 的 b2 别名、b3 结构体字面量两枚盲区
在真树里有货没货，我一字不知）。
最小闭合动作（可复算）：
`cd tools/d22scan && go run . -root ../../` 取 rc=0（现跑过，见 §5），
再对本仓每一枚盲区各写**一发只读反查尺**（例：ban #2 别名＝
`rg -n '^\s*(import\s+\w+\s+"path/filepath"|\w+ "path/filepath")' internal cmd`
比对 `rg -n '\bfilepath\.(Clean|Abs)\('` 的接收者名集合，差集非空即有货）。
本程**没**写这把尺。

**N2 — 5 枚 allowlist 条目"到底压掉了什么"零测量。**
`tools/d22scan/allowlist.txt` 有 5 条（`grep -vc '^#\|^$'` = 5）。我只在**假根**上量过
"一条目录前缀能整片压掉一枚禁令"（§3 的 `al-overbroad`），**没量过**这 5 条真条目今天各自
压掉了哪几行真代码，也没量过"删掉某条目门会不会多响"。而且门**从不打印抑制数**：
`al-overbroad` 那一跑的输出里没有任何"1 finding suppressed"字样（见其日志），
所以"allowlist 悄悄变宽"这件事今天对人是不可见的。
最小闭合动作：让 `verdict()` 在 `s.allowed()` 命中时计数并打印（下一程 AC#2 的射程），
或至少加一条自检断言"每条 allowlist 前缀至少命中一枚真 finding，否则该条目是死条目"。

**N3 — 只裁了"会不会响"，没裁"响得对不对/稳不稳"。**
每枚禁令我只放了 1–2 个形状。没测的相邻形状举例：
ban #1 的 `go` 在方法值/接口值上（`go s.worker()`，`goStmtText` 的 SelectorExpr 分支）、
ban #3 的 map 字面量与 `var x = struct{...}{...}`、ban #4 的 `t.Sub(u)` 两边都是
`time.Now()`、ban #8 在 `frontend/` 的**二进制文件**（`walkEmoji` 的 everyFile 分支明说"字节也被搜"，
我没取样）。⇒ 每一条都是"今天这一格响，相邻格没量过"。
最小闭合动作：把这些形状作为**新 cell** 加进 `bench.sh` 的 `CELLS`（台件已把这一步的成本压到
"加一个 case 分支 + 一行名册"），或按票面 AC#2 把它们装成 `--self-test` 的成对样本。

**N4 — 没测"门的自检能不能挂"。**
票面 AC#2（`--self-test` 双向样本）与 AC#4（摘掉违规样本 CI 会不会红）**本程一格未做**（派单 §1 明令）。
⇒ 今天"27 枚 cell 都给出读数"这件事本身**没有任何仪器守着**：`bench.sh` 是**一次性脚本**，
CI 不跑它，`runtests.sh` 不跑它。它买的是"这一次的数字"，**防回归能力＝0**。
最小闭合动作：AC#2 那一程把这 24 枚 bad/clean cell 的形状搬进 `tools/d22scan` 的自检入口，
`bench.sh` 降级为读数的来路记录。

**N5 — 门铃那一族（AC#3）与聚合退码（AC#6）都不在本程。**
"审批门今天有几枚生产发起者"我一字未量（兄弟程 161-r2 名义上在量门铃，**它的读数不在本件里，
我也没复算**）。票面 AC#6① 那枚"某腿点响、聚合仍退 0"我也**没复现**——顺带说：本程确实
**看见**了同一族的另一枚形状：`scripts/d22scan.sh` 第 1 步 `runtests.sh` 一挂，`set -eu`
就让**第 2 步的全仓扫描根本不跑**（编排者 21:23 那次 rc=1 就是这个形状：报的是测试的红，
全仓 1 枚 finding 是**另跑一次**才有的数）。⇒ 门的"rc"与门的"看没看见"今天仍不是一回事。

**N6 — 四数与名册差集是**同一棵工作树**上的两次跑，不是同一棵**提交**。**
共享树里别人在写。我的 pre/post 之间 `git status --porcelain | wc -l` 从 48 变到 …（§5 有原文），
**没有**任何一次读数是在 `git archive` 快照上取的（派单 §6 仪器坑 ⑤ 禁止这种做法，我没禁掉自己
的反向困难）。⇒ 若有人在两次跑之间动了 `tools/d22scan`，我的"没新增红"结论只对**测试名册**成立，
对**禁令语义**不成立。最小闭合动作：下一程在**自己的 worktree**（**不在本仓目录内**）里跑改前/改后。

**N7 — 退码之外的门输出没被断言。**
`verdict()` 打的 8 行 scope 计数（`examined NNN`）我只**抄了**，没**判**：
它现在读 `ban #8 design/=39`，而 `design/**` 此刻盘上有 16 枚未提交删除
（shared list 里 F9 那一枚正是这个风险）。⇒ 我没量"那 16 枚一旦入库 design/ 还剩下几枚、
会不会走到 0 让门硬退 2"，也没量"计数与真相源是否还一致"。
（旁证一处过期：`main.go:36` 的注释说 ban #6 的 `frontend/` 是 40 枚，**现跑是 85**。）

**N8 — 门禁里 `gofumpt` 与 `go vet` 两条命令**按派单原文跑不通**，我只做了等价修正后跑。**
- 原文 `go vet ./tools/d22scan/` ⇒ 实测 rc=1：`main module (github.com/CarlosShao/wisp)
  does not contain package …`（`tools/d22scan` 是自己的模块）。我改跑
  `cd tools/d22scan && go vet ./` ⇒ rc=0（与 ci.yml:124-126 那一步的形状一致）。
- 原文 `gofumpt -l . tools/d22scan tools/mockllm` **必须空** ⇒ 实测 rc=**2**、8 行：
  **3 行不是我产物**（`.scratch/wisp/probes/158/accept-r1/mut/guard.no{1,2,3}.go`，
  `git ls-files` 证明它们是**入库件**，即 CI 那一步今天在这棵树上会红），
  另 5 行是**我生成的 `runs/**` 里那两枚故意不可解析的 .go ＋一枚对齐差异**。
  处置：`runs/`、`bin/` 用 `.scratch/wisp/probes/161/r1/.gitignore` **挡在库外**（不提交），
  所以 CI 的检出里不会有我的行；那 3 枚别人的脏件**我一格未动**（不属本程写面）。
- 本程**没写任何 Go 源码**，所以 `gofumpt` 对我的意义只是"我没把谁的格式弄脏"。

**N9 — "会响条件"是**纸面承诺**，除已跑的那 27 发之外没有被机器验证。**
§4 每枚盲区都写了"什么样例一进去它就必须响"，但**那些样例我没造**（造了就等于本程把 AC#2 做了）。
⇒ 它们是**下一程的输入**，档位只能记〔我推断，未跑〕。

---
