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

## 3. AC#1 逐枚读数表（27 发；每发都有样本路径 / 命令 / rc / finding 原文 / 判定）

**命令原文（27 发同一形制，逐枚的完整一行印在 `logs/<cell>.scan.log` 第 2 行）**：

```
(cd "D:/work/workspace/projects plans/Wisp/tools/d22scan" && \
  .scratch/wisp/probes/161/r1/bin/d22scan-probe161.exe \
  -root "D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/161/r1/runs/<cell>")
```

- 该二进制是 `sh bench.sh tool` 现场从 `tools/d22scan` **读源码编出来的**（`go build -o`），
  本程没有写过它任何一个输入字节。新鲜度交叉核对见 §5.5。
- 每一发的输出第一行都是
  `d22scan: gitignore rules NOT APPLIED - ... is the subdirectory ".scratch/..." of a git repository, not its top`
  ⇒ 假根**不是一把 git 顶层索引**，matcher 因此**一条规则都不适用**（`gitignore.go:261-273`）。
  这对本程是**有利方向**：样本里不可能有文件被 ignore 悄悄跳过。
- 机器表 = `.scratch/wisp/probes/161/r1/logs/summary.tsv`（28 行含表头，逐字入档）。
- 判定用语：**响** = 该枚禁令的 `[tag]` 出现在 finding 里；**不响** = 同一 tag 一枚没有。

| # | cell | 禁令 | 样本文件（相对假根） | rc | finding 枚数 | 判定 |
|---|---|---|---|---|---|---|
| 1 | `skeleton-clean` | —（对照） | 整套骨架 16 枚文件 | 0 | 0 | 骨架本身不响 ⇒ 其余"不响"可归因 |
| 2 | `b1-closure-bad` | #1 | `internal/a/p_b1.go` | 1 | 1 | **响** |
| 3 | `b1-named-bad` | #1（具名调用，R16#1 那一切） | `internal/a/p_b1.go` | 1 | 1 | **响** |
| 4 | `b1-clean` | #1 | `internal/a/p_b1.go`（同词只在注释里） | 0 | 0 | **不响** ⇒ 尺不是恒响 |
| 5 | `b1-probe-testfile` | #1 | `internal/a/p_b1_test.go` | 0 | 0 | 探针：**不响** ⇒ B3 |
| 6 | `b2-bad` | #2 | `internal/a/p_b2.go`（`filepath.Clean`+`Abs`） | 1 | 2 | **响**（两臂都活） |
| 7 | `b2-clean` | #2 | 同路径（Join/Base/Dir/Rel＋他接收者的 `Clean()`） | 0 | 0 | **不响** |
| 8 | `b2-probe-alias` | #2 | 同路径（`import fp "path/filepath"` → `fp.Clean`） | 0 | 0 | 探针：**不响** ⇒ B1 |
| 9 | `b3-bad` | #3 | `internal/a/p_b3.go`（const 臂＋assign 臂） | 1 | 2 | **响**（两臂都活） |
| 10 | `b3-clean` | #3 | 同路径（名字闸全过，值靠形状/前缀挡下） | 0 | 0 | **不响** |
| 11 | `b3-probe-structlit` | #3 | `internal/a/p_b3s.go`（结构体字面量字段） | 0 | 0 | 探针：**不响** ⇒ B2 |
| 12 | `b4-sub-bad` | #4 | `internal/a/p_b4.go`（`deadline.Sub(time.Now())`） | 1 | 1 | **响** |
| 13 | `b4-unix-bad` | #4（第二臂） | 同路径（`timeoutUnix := time.Now().Unix()+30`） | 1 | 1 | **响** |
| 14 | `b4-clean` | #4 | 同路径（`time.Until`/`time.Since`/`b.Sub(a)`） | 0 | 0 | **不响** |
| 15 | `b5-bad` | #5 | `internal/a/p_b5.go` | 1 | 2 | **响**（`:6` 声明行＋`:7` 纯引用行，见 A2） |
| 16 | `b5-clean` | #5 | 同路径（hash 词与 mirror 词分行） | 0 | 0 | **不响** |
| 17 | `b6-bad` | #6 | `frontend/p_b6.ts` | 1 | 1 | **响** |
| 18 | `b6-clean` | #6 | 同路径（`approval.request`） | 0 | 0 | **不响** |
| 19 | `b6-probe-comment` | #6 | 同路径（`approval.decide` **只出现在注释里**） | 1 | 1 | 探针：**响** ⇒ A1 |
| 20 | `b7-bad` | #7 | `internal/tools/p_b7.go` | 1 | 3 | **响**（`:5` 点号名＋`:7` 裸词；`:4` 见 A1） |
| 21 | `b7-clean`（第二版） | #7 | 同路径（近邻拼写两枚） | 0 | 0 | **不响** |
| 21b | `b7-clean` 第一版 | #7 | 同上，但注释里带了一对引号词 | 1 | 1 | **响在注释** ⇒ A1，日志另存 `logs/b7-clean.first-run-rang-on-comment.log` |
| 22 | `b8-scope-bad` | #8 | 四枚作用域各一发：`design/p_b8.html`、`frontend/p_b8.ts`、`internal/a/p_b8.go`、`cmd/probe/p_b8.go` | 1 | 4 | **响，四枚 scope 全活** |
| 23 | `b8-comment-clean` | #8 | 同四路径的"字形只进注释"版 | 0 | 0 | **不响**（Q-46(c) 豁免成立） |
| 24 | `b8-probe-bands` | #8 | `internal/a/p_b8.go` 四枚字形各一：`→`/`①`/`✓`/`≤` | 1 | 2 | **只 `✓`+`≤` 响** ⇒ B4 |
| 25 | `unp-bad` | `unparseable` | `internal/a/p_unp.go` | 1 | 2 | **响**（`unparseable`＋同文件那枚 #8 字形仍响） |
| 26 | `unp-hides-bans` | #1–#5 | 同路径（不解析的文件里再塞 `go func(){}`＋`filepath.Clean`） | 1 | 1 | **只响 `unparseable`** ⇒ B5 |
| 27 | `al-overbroad` | #2/#1 | `internal/a/p_al.go`＋本 cell 自带的假 allowlist | 1 | 1 | `bare-goroutine` 响、**同文件 `pathresolver-bypass` 被前缀压掉且零打印** ⇒ B6 |

### 3.1 finding 原文（逐字，从 `logs/*.scan.log` 抄）

```
b1-closure-bad  : internal/a/p_b1.go:6: [bare-goroutine] bare `go func(` is banned (D22/D38b): use observe.Registry.Spawn (named, owner, recover boundary)
b1-named-bad    : internal/a/p_b1.go:7: [bare-goroutine] bare `go probeReader(...)` is banned (D22/D38b, R16: named calls count too): use observe.Registry.Spawn (named, owner, recover boundary)
b2-bad          : internal/a/p_b2.go:7: [pathresolver-bypass] filepath.Clean outside the C26 PathResolver is banned (D22); see tools/d22scan/allowlist.txt for the sanctioned exceptions
b2-bad          : internal/a/p_b2.go:9: [pathresolver-bypass] filepath.Abs outside the C26 PathResolver is banned (D22); see tools/d22scan/allowlist.txt for the sanctioned exceptions
b3-bad          : internal/a/p_b3.go:6: [plaintext-key] identifier probeAPIKey holds a plaintext key literal (D22/C28: SecretStore refs only)
b3-bad          : internal/a/p_b3.go:9: [plaintext-key] identifier probeSecret assigned a plaintext key literal (D22/C28: SecretStore refs only)
b4-sub-bad      : internal/a/p_b4.go:6: [wallclock-timeout] wall-clock delta (Sub(time.Now())) in what must be monotonic logic (D42#9/D22)
b4-unix-bad     : internal/a/p_b4.go:7: [wallclock-timeout] unix timestamp on a timeout/deadline line - monotonic clock required (D42#9/D22)
b5-bad          : internal/a/p_b5.go:6: [mirror-hash] hash material mentioned together with a mirror (C29/F3: hashes come from the signed manifest only)
b5-bad          : internal/a/p_b5.go:7: [mirror-hash] hash material mentioned together with a mirror (C29/F3: hashes come from the signed manifest only)
b6-bad          : frontend/p_b6.ts:3: [panel-approval] `approval.decide` in frontend/ is banned (D33/F2: allow decisions are native-side only)
b6-probe-comment: frontend/p_b6.ts:5: [panel-approval] `approval.decide` in frontend/ is banned (D33/F2: allow decisions are native-side only)
b7-bad          : internal/tools/p_b7.go:4: [internal-artifact-tool] host-internal artifact write looks implemented as a gated tool name (D34 note 2/D22)
b7-bad          : internal/tools/p_b7.go:5: [internal-artifact-tool] host-internal artifact write looks implemented as a gated tool name (D34 note 2/D22)
b7-bad          : internal/tools/p_b7.go:7: [internal-artifact-tool] host-internal artifact write looks implemented as a gated tool name (D34 note 2/D22)
b7-clean 第一版 : internal/tools/p_b7.go:4: [internal-artifact-tool] host-internal artifact write looks implemented as a gated tool name (D34 note 2/D22)
b8-scope-bad    : design/p_b8.html:3: [emoji] ban #8 glyph in scope design/ is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)
b8-scope-bad    : frontend/p_b8.ts:1: [emoji] ban #8 glyph in scope frontend/ is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)
b8-scope-bad    : internal/a/p_b8.go:5: [emoji] ban #8 glyph in scope internal/ is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)
b8-scope-bad    : cmd/probe/p_b8.go:3: [emoji] ban #8 glyph in scope cmd/ is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)
b8-probe-bands  : internal/a/p_b8.go:11: [emoji] ... 同一句话（该行字形是 U+2713）
b8-probe-bands  : internal/a/p_b8.go:12: [emoji] ... 同一句话（该行字形是 U+2264）
unp-bad         : internal/a/p_unp.go:1: [unparseable] D:\work\...\runs\unp-bad\internal\a\p_unp.go:5:14: expected ')', found '{' (and 1 more errors)
unp-bad         : internal/a/p_unp.go:3: [emoji] ban #8 glyph in scope internal/ is banned (D23): ...
unp-hides-bans  : internal/a/p_unp.go:1: [unparseable] D:\work\...\runs\unp-hides-bans\internal\a\p_unp.go:15:14: expected ')', found '{' (and 1 more errors)
al-overbroad    : internal/a/p_al.go:11: [bare-goroutine] bare `go func(` is banned (D22/D38b): use observe.Registry.Spawn (named, owner, recover boundary)
```

**结论行（本格真正的答话）**：**`// Bans` 名单上 8 枚编号禁令，今天 8 枚都点得得到**，
每一枚都同时给出"违规响＋干净不响"这一对；`unparseable` 那枚发现类型也响。
**没有一枚禁令是死的。** 门今天不瞎——但它今天**瞎在六个更窄的地方**，逐条在 §4。

---

## 4. 盲区清单与各自的"会响条件"

`今天量到的＝静默` 全部是 §3 表里的现跑读数；`会响条件` 是**下一程可直接抄成 cell 的判据**，
档位＝〔我推断，未跑〕（本程明令不装自检，所以这些条件我一个都没兑现）。

**B1 — ban #2 对别名导入让步（`b2-probe-alias`：rc=0，零 finding）**
- 形状：`import fp "path/filepath"` + `fp.Clean(p)`。matcher 判的是
  `sel.X.(*ast.Ident).Name == "filepath"`（`main.go:715`），即**接收者的局部名字**，不是**包身份**。
- 与禁令文字的关系：文字是"`filepath.Clean` / `filepath.Abs` 在 allowlist 之外"——
  换别名以后**决策仍然在 C26 PathResolver 之外** ⇒ 这是**门的射程比它自己的文字窄**。
- 会响条件：假根 `internal/a/p_b2.go` 含 `fp "path/filepath"` 且出现 `fp.Clean(` ⇒ 必须
  `[pathresolver-bypass]` 响。修法两选一（**都属下一程/契约面，本程不选**）：
  按 import 的**路径**解析包身份，或把"局部名 ≠ filepath 的 Clean/Abs"另立一枚新禁令。
- 这一枚登记买什么：**防回归**——真树今天有没有别名用法我**没量**（§2 N1）；
  先把"门比自己的文字窄"钉住，别让下一位以为文字即射程。

**B2 — ban #3 只看"声明/赋值"两种左值（`b3-probe-structlit`：rc=0）**
- 形状：`probeCfg{APIKey: "<合成串>"}`（合成字面量的字段）。AST 只走 `ValueSpec` 与
  左值为 `Ident` 的 `AssignStmt`（`main.go:723-755`），`*ast.CompositeLit` 里那枚
  `KeyValueExpr` 的 `BasicLit` 没人看。
- 同族没量的：map 字面量、函数实参位、结构体指针字段。
- 会响条件：假根 `internal/a/p_b3s.go` 含 `type T struct{ APIKey string }` +
  `var v = T{APIKey: "<16+ 字符、含数字、无空白>"}` ⇒ 必须 `[plaintext-key]` 响。
- 买什么：**防忘记**（这是 C28 的语义洞、不是拼写洞；补它要动 ban #3 的实现射程 ⇒ 下一程）。

**B3 — bans #1–#5 对 `_test.go` 整体不看（`b1-probe-testfile`：rc=0）**
- 这是**声明过的作用域**（`main.go` 文档注释："production scope internal/ + cmd/
  (non-test, non-testdata)"），不是 bug。但它与 AGENTS.md 禁止清单的文字（"裸 `go func(`
  而无 owner/recover"）**没有 test 豁免**。
- 我顺手查了有没有别的仪器补这一格：**没有**。同族的 `internal/observe/nobarego_test.go`
  自己也 `strings.HasSuffix(path, "_test.go")` 就跳过（该文件 39 行）
  ⇒ 测试文件里的裸 goroutine **今天零枚仪器**。
- 会响条件：假根 `internal/a/p_b1_test.go` 含 `go func() {}()` ⇒ 要么 d22scan 响（改射程＝契约级），
  要么新增一枚管 test 的尺。**这一条属"未定义即停"**：该不该管是 owner 的话，不是本程的。
- 买什么：**防忘记**（把"文字说不得、仪器没看着"这一对记在同一处）。

**B4 — ban #8 的两段波段空隙（`b8-probe-bands`：四发里两发不响）**
- 实测量到的：字符串里的 `✓`(U+2713) 与 `≤`(U+2264) **响**；`→`(U+2192) 与 `①`(U+2460) **不响**。
- 与 `main.go:156-164` 的自陈**一致**（`2190–21FF`、`2460–24FF` 是刻意留的空隙，由
  `TestBan8MathBandAndRemainingGaps` 正反钉住），也**证实 AGENTS.md §1.2 那句"实际后果"**。
- 所以这一枚**不是新洞**，是"已知洞＋我第一次拿真样本量到它"。会响条件（反向）：
  假根字符串里放 `→` ⇒ 今天**必须不响**；一旦有人把箭头段并进去，这一发就该翻响并且
  **该被名册看见**（那是 `Q-46`/票 141 的契约面，本程一字节不碰）。
- 买什么：**防回归**（钉住"今天不响"这个事实，波段哪天被放宽时才发现多出的一响）。

**B5 — 一枚不解析的文件对 #1–#5 完全隐形（`unp-hides-bans`：rc=1，只响 `unparseable`）**
- `scanGoFile` 在 `parser.ParseFile` 失败处 `s.add("unparseable", …); return nil`
  （`main.go:684-689`）⇒ 该文件的 #1/#2/#3（AST）**和同一个函数后半段那两枚逐行禁令
  #4/#5 全部跳过**。实测：文件里同时塞了 `go func(){}` 与 `filepath.Clean(`，零枚响。
- 危害形状：**修好语法再跑一次**才可能露出裸 goroutine；而今天这一跑给人的话是
  "1 finding: unparseable"，一个字的"这枚文件我五道禁令都没看"都没有。
- 会响条件：假根里一枚**故意不解析、内含 #1 形状**的文件 ⇒ 输出必须**同时**含
  `[unparseable]` 与 `[bare-goroutine]`（做法：解析失败时退回逐行正则那一族，或明打印
  "bans #1-5 skipped for this file"）。
- 买什么：**防回归**（这就是"门还在、还在报通过，但它已经看不见东西"的**原生形状**，
  票 161 立票要防的正是它）。

**B6 — allowlist 抑制零枚计数、零枚打印（`al-overbroad`：同文件 `filepath.Clean` 静默消失）**
- 假 allowlist 一行 `pathresolver-bypass<TAB>internal/<TAB>…` ⇒ **整个 internal/ 的 #2 从此不响**，
  而 `bare-goroutine` 照常响（抑制是**按枚**的，这点与文字一致）。
- 更要紧：`verdict()` 的 8 行自报告里**没有任何"本次抑制 N 条"**
  （见 `logs/al-overbroad.scan.log` 全文）⇒ "allowlist 被写宽"这件事对人不可见。
  真仓今天 **5 条**（`grep -vc '^#\|^$' tools/d22scan/allowlist.txt` = 5），两条是整文件前缀。
- 会响条件：假根 allowlist 含一条**目录前缀**且样本里该目录下有该枚违规 ⇒
  输出必须打印 `suppressed=N (by entry …)`，或让"条目零命中"本身成为一枚红。
- 买什么：**防忘记**（`allowlist.txt` 在 AC#5 冻结名单里，动它＝人工批准；先把"它无声"记下来）。

### 4.1 三枚"不是盲区、是脾气"的读数（同批 cell 顺手量到的，别当新洞报）

- **A1 注释策略在一把门里有三套**：`b6-probe-comment`（`frontend/**` 的注释里写了
  `approval.decide` ⇒ **响**）与 `b7-clean` 第一版（`internal/tools/**` 注释里带了一对引号词
  ⇒ **响**，见 `logs/b7-clean.first-run-rang-on-comment.log`）。
  同一工具里：#8 走 AST/词法注释分类器（豁免），#4/#5 走"整行以 `//` 开头就跳过"，
  #6/#7 走 `walkText`**完全不剥注释**。⇒ 实际后果：**给面板写一句解释 `approval.decide`
  的注释，今天就能把 CI 弄红**；这也解释了票 169 为什么一枚装饰字符就同时摘掉几道门。
  （没量到的一侧：#5 的整行注释豁免我只在**码里**看见（`main.go:761`），样本里两个词没同现在
  同一行 ⇒ 那条属读码、不属读数，档位〔未取样〕。）
- **A2 ban #5 的精度＝同行字面共现**：`b5-bad` 第 2 枚响在 `return mirrorSHA256`——
  一行**纯引用**，没有任何"从镜像取哈希"的数据流。⇒ 这条尺看名字同现、不看流向：
  误报面由此来，反过来它也**不会**抓到拼写不同的真取哈希代码。
- **A3 ban #8 的 `frontend/` 是 everyFile**：二进制也按字节搜（`walkEmoji` 的 everyFile 分支自陈）。
  本程**没取样**二进制（§2 N3 已列）。

---

## 5. 门禁：四数 + 名册两向 `comm` 差集（全部我本轮现跑，锚 `fbe12c7`）

### 5.1 仪器版本（先现跑再贴）

```
$ go version
go version go1.27.1 windows/amd64
$ /d/work/base/gopath/bin/gofumpt.exe --version
v0.12.0 (go1.27.1)
```

⚠ `gofumpt` 不在本 shell 的 PATH 上（`which gofumpt` rc=1、裸 `gofumpt --version` rc=127），
只在 `go env GOPATH`/bin 里有 `gofumpt.exe` ⇒ 以下均用全路径。环境事实，不是判据。

### 5.2 `tools/d22scan` 那一包：改前 / 改后四数 + 名册

| 读数 | pre（step-0 之后、写任何东西之前） | post（27 发跑完、§0–§2 提交之后） |
|---|---|---|
| 命令 | `sh tools/d22scan/runtests.sh -C tools/d22scan ./...` | 同 |
| rc | 0 | 0 |
| 四数 | `PASS=30 FAIL=0 SKIP=0`、`=== RUN=70` | `PASS=30 FAIL=0 SKIP=0`、`=== RUN=70` |
| FAIL 名册 | 空 | 空 |
| 日志 | `logs/pre-baseline-runtests.log`（225 行） | `logs/post-baseline-runtests.log`（225 行） |

名册两向差集（`grep '^--- PASS'` 取顶层名 → `sort -u` → `comm -23` / `comm -13`）：

```
pre-only :  (空)
post-only:  (空)
counts: pre=30 post=30        # logs/pre-passnames.txt, logs/post-passnames.txt
```

**§4 点名的那两枚今天什么样**：两本日志里都在、且**都是 PASS**——
`post-baseline-runtests.log:23  --- PASS: TestScannerSelfScanOfRealRepoIsGreen (0.80s)`、
`:138 --- PASS: TestRealRepoLedgerIsHonest (0.85s)`。
⇒ 我的门禁口径落地为：**"除（现已不存在的）那两枚之外零枚新增红" ＝ 零枚红；名册差集两向皆空**。
判据文字、断言、`frontend/**` **一概没动**（三件禁事一件都没做）。

### 5.3 `go vet`

```
$ go vet ./tools/d22scan/                      # 派单 §6 原文
main module (github.com/CarlosShao/wisp) does not contain package .../tools/d22scan   rc=1
$ cd tools/d22scan && go vet ./                # 与 ci.yml:124-126 同形制
rc=0                                           # logs/pre-baseline-govet.log / -inmodule.log
```

⇒ 派单 §6 那条**原文不可跑**（跨模块），修正后的形状 rc=0。本程一枚 Go 源码都没写，
这一发的作用是"证明没写坏、而不是没跑"。

### 5.4 `sh scripts/d22scan.sh`（全仓那一腿）

```
$ sh scripts/d22scan.sh          rc=0                       # pre 与 post 各一次
finding 枚数 = 0                # grep -cE '^[A-Za-z0-9_./-]+:[0-9]+: \[' = 0
```

自报告（`logs/post-d22scan-sh.log` 末尾，逐字）：

```
d22scan: examined 228 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: scope bans #1-5 internal/      examined 205 production Go files
d22scan: scope bans #1-5 cmd/           examined  23 production Go files
d22scan: scope ban #6 frontend/         examined  85 text files
d22scan: scope ban #7 internal/tools/   examined  18 production Go files
d22scan: scope ban #8 design/           examined  39 text files
d22scan: scope ban #8 frontend/         examined  85 text files
d22scan: scope ban #8 internal/         examined 414 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  45 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; ...
```

⇒ 派单 §4 的"rc=1／1 枚 finding"在锚上**不成立**（来路见 §1.1）。
另记两处（不改任何东西，只记档）：
① `main.go:36` 的注释写 ban #6 的 `frontend/` 是 **40** 枚，现跑是 **85**；
② `set -eu` 让 `runtests.sh` 一挂就**根本不跑**第 2 步的全仓扫描 ⇒
"脚本 rc"与"门看没看见"不是一回事（同族形状已登记在票面 **AC#6①**，归下一程）。

### 5.5 台件自身的两件小事（新鲜度与格式）

- **二进制不是过期货**：同一 cell 两条路各跑一次，finding 逐字相同
  （`logs/b8-scope-bad.scan.log` 那 4 行 vs `logs/gorun-b8-scope-bad.log`，
  后者是 `go run . -root <假根>`；两本唯一差别是 `go run` 自己加的 `exit status 1`）。
- `gofumpt -l . tools/d22scan tools/mockllm`（派单 §6 要求"必须空"）：

```
pre :  rc=0，3 行 —— .scratch/wisp/probes/158/accept-r1/mut/guard.no{1,2,3}.go（**不是我产物**，
        git ls-files 证明是入库件；盘上是 LF，故非 CRLF 假象）
post:  rc=2，8 行 —— 上面 3 行 ＋ 我生成的 runs/ 里 5 行（两枚故意不可解析的 .go ＋一枚对齐差异）
```

  处置：(a) `runs/`、`bin/` 已被 `.scratch/wisp/probes/161/r1/.gitignore` 挡在库外 ⇒
  提交里没有它们，CI 的检出不会因为我多一行；
  (b) 那 3 行**本程一格未动**（不是我写面，且 `guard.no*.go` 是票 158 的变异台件，
  "修它"＝踩别人的实验）。
  ⇒ 按 §9 报回：**"gofumpt 必须空"这一条在我的锚上本来就不成立**。
  ⚠ CI 用 `gofumpt@latest`、我本机 v0.12.0，两边会不会同判**我没跑过 CI**（档位〔仅本机读数〕）。

### 5.6 AC#5 契约轴零字节（名册按 commit 集合算）

第一枚提交（`16364e6`，§0–§2 ＋台件）之后现取：

```
$ git diff --name-only fbe12c7..HEAD | wc -l
45        # 全部落在 .scratch/wisp/probes/161/r1/** ＋ docs/evidence/s1/161-gate-blindspot-r1.md
$ git diff --name-only fbe12c7..HEAD | grep -E 'docs/PLAN\.md|^docs/specs/|^internal/risk/|^internal/panel/|^internal/agent/approval/|^tools/d22scan/|^\.github/|^frontend/|^design/|thresholds\.go$|allowlist\.txt$|slo-check\.ps1$|golden'
(none: 契约轴零字节)
```

票面 AC#5 的 11 枚点名件一枚未碰，派单特别名单三枚（`tools/d22scan/**`、`.github/workflows/**`、
`docs/reports/**`）一枚未碰。最后一枚提交之后会重跑这一发并把结果补在下面（名册是
**每次提交后现取**的，不是一次性结论）。

**终判（本程两枚提交之后现跑，逐字；22:0x）**：

```
$ git log --format='%h %ad %s' --date=format:'%H:%M' fbe12c7..HEAD
e8092351 21:59 evidence(票161 r1 §3-§9): AC#1 逐枚读数 27 发＋六枚盲区登记——…
16364e63 21:53 evidence(票161 r1 §0-§2): 门禁盲区普查台件＋读数基线——…

$ git diff --name-only fbe12c7..HEAD | wc -l
49                                    # logs/final-ac5-roster.txt（三枚提交之后）

$ 上面那份名册里，落在契约轴/特别名单任何一枚之内的：
(none: 契约轴零字节)
$ git diff --name-only fbe12c7..HEAD | grep -v -e '^\.scratch/wisp/probes/161/r1/' -e '^docs/evidence/s1/161-gate-blindspot-r1\.md$'
(none)                                # 45 行全在我两处写面内
$ git diff --stat fbe12c7..HEAD -- tools/d22scan .github frontend design internal cmd docs/PLAN.md docs/specs docs/reports
(空)                                  # 契约轴与特别名单：一个字节都没有
```

**两向名册差集的最终一次（第三本日志，两枚提交都提完之后再跑一次整包）**：

```
$ sh tools/d22scan/runtests.sh -C tools/d22scan ./...      rc=0
PASS=30 FAIL=0 SKIP=0  === RUN=70                          # logs/final-runtests.log
pre 30 枚顶层名 vs final 30 枚：comm -23 = 空、comm -13 = 空（logs/pre-|final-passnames.txt）
```

⇒ 门禁口径落地：**三本日志（pre / post / final）四数完全相同、名册两向差集皆空、零枚红**。
派单 §4 那两枚"已知红"在本锚上不存在（来路 §1.1），所以"除那两枚之外零枚新增红"这一句
今天的实值就是"**零枚红**"——**只可能更严，不可能更松**，判据文字我一个字没动。

⚠ **一句自指，免得下一位以为 45/49 谁抄错了**：上面那个 **49** 是第 3 枚提交之后取的数，
而**写着这个数的这一行本身属第 4 枚提交**——名册只多这一枚文件、且仍落在同两处写面内。
复算式（跑一次就得到一个更新的数，且永远比本行落后一枚提交）：
`git diff --name-only fbe12c7..HEAD | grep -v -e '^\.scratch/wisp/probes/161/r1/' -e '^docs/evidence/s1/161-gate-blindspot-r1\.md$'`
⇒ 输出为空即"仍在我写面内"。第三遍全仓扫描也重跑过：
`sh scripts/d22scan.sh` **rc=0／0 枚 finding**（`logs/final-d22scan-sh.log`）。


### 5.7 一件必须让编排者知道的 Git 现场处置（票面追加**没有**进任何提交）

票面 161 的 Progress log 追加**已经落盘**（`.scratch/wisp/issues/161-…md`，sha256
`1daa1542…` 记在 `logs/ticket-file-sha-before.txt`），**但没提交**。原因不是我忘了：

- 那枚文件在我进场时就是 ` M`——里面是**编排者 21:2x 的 53 行未提交在途编辑**
  （AC#6、CI 写面边界、"九条→`// Bans` 现读"那三段更正）。
- 我 `git add` 该文件后 `git diff --cached` 一看就是：`53 5` 的 numstat 里绝大部分**不是我的行**。
  用**文件级 pathspec** 提交它，仍然会把它们并入我的提交＝`5afa666`/`15ff2be` 记过的那枚形状
  （只是从"目录级"缩到"文件级"，缩掉的那一段这次救不了我）。
- 处置：本程只提交证据件（`git commit -q -F - -- docs/evidence/s1/…`），
  再用 `git restore --staged -- <票面>` 把索引**退回我进场时的样子**
  （前后 `sha256sum` 同串 ⇒ 内容零丢失；`git diff --cached` 回到 0 枚；`git status` 回到 ` M`）。
  派单禁的六枚动词一枚没用（`restore --staged` 只动索引、不动工作树、不改历史，理由写在上面）。
- ⇒ **请编排者**：要么你先提交自己那 53 行、我下一程再提我的追加；要么你直接把我落盘的那一段
  连同你的改动一起提。**不要由我在共享树里替你猜。**

---

## 6. 被权限系统拒绝的调用（逐枚）

**0 枚。** 本程没有任何一次工具调用被权限系统拒绝，也没有尝试绕过。
命令级失败 3 次，全部如实列出（都发生在取得任何 finding 读数**之前**）：

| # | 命令 | rc | 因 | 之后怎么做 |
|---|---|---|---|---|
| 1 | `gofumpt --version`（裸名） | 127 | 不在本 shell PATH | 用 `go env GOPATH`/bin 全路径重跑（§5.1） |
| 2 | `go vet ./tools/d22scan/` | 1 | 跨模块（派单 §6 原文的形状） | 改在模块内 `go vet ./`（§5.3） |
| 3 | `sh bench.sh all`（前两次） | 1 | 我脚本自己的两个 bug：`case` 少一枚 `;;`；`run()` 内部把 `set -e` 重新打开，杀掉了调用者的循环 | 修脚本后重跑 ⇒ **没有一条读数来自半截的跑**（第一本完整日志时间戳 21:45） |

除此之外：**无权限拒绝、无绕过、无 `sudo`/`--no-verify` 类操作。**

---

## 7. 伪授权两栏计数（两枚字段分开，不并成一枚）

| 栏 | 数 | 说明 |
|---|---|---|
| **真通知回显数** | **0** | 本程没收到任何一条"任务完成／验收通过／票已结案"型的回显。据以行事的输入只有：派单文件、票面、`tools/d22scan` 源码、我自己现跑的命令输出。 |
| **判为注入数** | **0** | 没有一段文本被我当成授权处理，因此也没有"识别出来但不服从"的条目。唯一一次接近的形状是下面的读数污染事件，我把它当**读数污染**处理，未据它写任何东西。 |

**读数污染事件（不是伪授权，但同一族"该信谁的读数"）**：
`Read` 工具对 `bench.sh` 连续两次返回**与盘上不符**的内容（行号错位、显示一份我从未写过的
566 行版本），一度让我判定"另一程在改我的文件"。判据用在了对的地方——**信盘、不信回显**：

```
$ git show HEAD:.scratch/wisp/probes/161/r1/bench.sh | sha256sum
817fa6ff36ab083accebd1d5eec81f2cf9292eb13f93579ade4dce18869b6862
$ sha256sum .scratch/wisp/probes/161/r1/bench.sh
817fa6ff36ab083accebd1d5eec81f2cf9292eb13f93579ade4dce18869b6862      # 同一串 ⇒ 没人改过
$ bash -n bench.sh && echo SYNTAX OK                                   # 通过
```

⇒ 结论：**不存在外部写入者**，是我的一个视图工具过期。本程**没有**因为那段回显重写/回滚
`bench.sh`（真要那么做就会覆盖掉盘上的东西）。补一条仪器坑给下一位：
**多字节字形密集的脚本，`Read` 的视图可能与盘上不一致；判归属用 `sha256sum` ＋ `awk`。**

---

## 8. 凭据值零抄录声明

- 本件、`bench.sh`、`runs/**`、`logs/**` **没有出现任何真凭据值**。
- ban #3 的两枚样本串是我**为这一格编的合成字面量**（形如
  `PROBE161NOTACREDENTIALED…`，命名本身自陈"不是凭据"），不属于任何服务、账号或环境。
- 本程没有打印过 env、没有读过 `.env`、没有取过任何 token 变量值；
  与"密钥"唯一沾边的读数是 `tools/d22scan/allowlist.txt` 的**路径前缀与理由文字**。
- 若下一位在日志里看到疑似值：那只会是 `WISP_PROBE_*` 这类**我造的变量名**，不是值。

---

## 9. next=（撞轮次上限之前先把这一节写完）

给票 161 **下一程（AC#2：把自检装上）**的话，本程一格都不代做：

1. §3 那 27 枚 cell 就是 `--self-test` 的成对样本现成清单——照 `bench.sh` 里
   `expect_of()` 那 28 行抄成表即可（cell, ban, ring|silent），一行不落地两向都过。
2. §4 的 **B1/B2/B5/B6** 各带一发"会响条件"，先当**失败预演**用；
   装成"允许失败但会打印"的那条腿是错的形状（要么响，要么明写不覆盖）。
   **B3 属"未定义即停"**：`_test.go` 该不该进 ban #1–#5 的射程，是 owner 的一句话。
3. §2 N2：`allowlist.txt` 的抑制今天**零打印**——那是"门自己瞎了还印 clean"最短的一条路。
4. §5.5 那 3 行 `gofumpt` 脏件不是我产物、但在库内；归谁收请编排者定（本程未动）。
5. 本程所有读数的锚是 **`fbe12c7`**。**任何"改前基线"请重新现量**，不要抄 §5.2 的四数——
   编排者 21:23 那一次就已经被 21:26 的提交作废过一次，形状与票 158 的教训一模一样。

