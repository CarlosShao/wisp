# 票 251 r1 交件读数（写码腿 `251-r1`）

> 本节只交**读数**，判语归非实现者。所有行号／枚数＝本腿本机自跑，不是抄票面。
> 四枚 AC 的勾选框本腿一枚未碰（本仓硬规矩：`- [ ]`→`- [x]` 归编排者）。

## 0. 起手锚

| 项 | 读数 |
|---|---|
| 取锚时刻 | `2026-10-02T09:29:26+08:00`（`date -Iseconds`） |
| 起手 HEAD | `f5f9cc3433385c13e24cb2a7d8054d8210f99cde`（`dev`，`git log -1`） |
| 起手 `git status --porcelain scripts/` | **空输出**（rc=0）⇒ 这一刻没有别的腿在 `scripts/` 这面上，未触发票面「交件要求」第 1 条的停手 |
| 起手 `wc -l -c` | `portable-tests.sh` 572/32543 · `portable-tests-selftest.sh` 313/13354 · `testdata/portable-tests/go` 100/4159 |
| 起手载具基线 | `bash scripts/portable-tests-selftest.sh` ⇒ `10 case(s) ran, 0 assertion(s) failed`、**rc=0**（`logs/selftest-baseline.txt`） |
| 起手真 census | `bash scripts/portable-tests.sh --scope=census` ⇒ rc=0，`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=9`，GOOS=windows（`logs/census-prefix.txt`） |
| 交件时 HEAD 已前移 | 起手后 dev 上落了腿 252 的两笔（`3465dcee` 给 `internal/tools` 加了一枚探针、`ebe2945a` 台账）；`git log f5f9cc3..HEAD -- scripts/` 里**只有本腿的四笔**，故 §3 那把 `git diff f5f9cc3..HEAD -- scripts/portable-tests.sh` 的尺不被别人污染。真 census 两次读数之间唯一一行差异＝`internal/tools` 的 `38/13→39/13`，出处是 `3465dcee`，与本腿无关（`logs/census-postfix.txt` 与 `logs/census-prefix.txt` 对差） |

起手锚之后所有卫生尺都以 `f5f9cc3` 为左端。

## 1. 改前必红（缺口凭据）

### 1.1 票面「现量」三把尺，本腿复跑结果

1. `bash scripts/portable-tests.sh --scope=winsec` ⇒ 逐字（stderr，`logs/prefix-scope-winsec.txt`）
   `portable-tests.sh: unknown --scope=winsec (known: core, windows, cli, census)`、**rc=2**。
   ⇒ 复认票面 12-14 行：这一形**是响亮失败、不是静默绿**。
2. `--scope=census` 真读数里 winsec 那一行（`logs/census-prefix.txt:36`）逐字
   `portable-tests.sh: github.com/CarlosShao/wisp/internal/winsec     12/8        corewinsec`
   ⇒ census 把一枚 **`case $mode in` 里根本不存在的档名**报成"有人认领"，且这一步 **rc=0**。
   档名两处各写一遍今天**已经不等**，而不响（AC#4 的凭据，用的是真 `go list`，不是种子）。
3. census 全表里 `internal/winsec` 只有 **1 行**（`grep -c` 尺），⇒ 今天它没有子包，
   故 `./internal/winsec/...` 与 `./internal/winsec/` 现在解析出同一枚单包——本腿把第五档写成 glob 踩的就是这条实测。
4. 改前的 `case $mode in`（锚点 `:170-206`）分支标签枚数＝**4 枚具名 + 1 枚兜底**，
   抽取尺：`awk '/^case \$mode in$/{f=1;next} f&&/^esac$/{exit} f&&/^[a-z*]+\)$/&&$0!="*)"{c++} END{print c+0}'` ⇒ `4`。

### 1.2 两形种子打在**改前字节**上（都是 CI 今天真走的显式路径）

载体：`.scratch/wisp/probes/251/r1/seed-winsec.sh`（只假 `go` 这一个进程，脚本与 GUARD A/B/C 与 `tools/d22scan/runtests.sh` 都是被审的真码；⛔ 全程未起真 `go test`／`go build`）。
改前字节＝`git show f5f9cc3:scripts/portable-tests.sh > /tmp/251-prefix-portable-tests.sh`，由载具的 shadow-root 机制跑。

| 种子形 | 改前（锚点字节） | 改后（工作树字节） | 出处 |
|---|---|---|---|
| 拆包：`internal/winsec` 旁多出第二包 | **rc=0，全绿**，`GUARD C` 一字未出，末行 `four numbers ... === RUN=2 --- PASS=2 --- FAIL=0 --- SKIP=0` | **rc=1** `GUARD C - scope mode=winsec resolved to a DIFFERENT package` / `Pinned: 1, resolved: 2` | `logs/seed-prefix.txt`（`split-explicit-winsec`）、`logs/seed-postfix.txt`、逐字全日志 `logs/prefix-split-full.txt` 与 `logs/postfix-split-full.txt` |
| 删真包名：`internal/winsec` 一行都解析不出 | **rc=0，全绿**（解析集为空也照绿：GUARD B 审了一个零行名册） | **rc=1** `GUARD C ... Pinned: 1, resolved: 0`，diff 行 `< github.com/CarlosShao/wisp/internal/winsec` | 同上（`gone-explicit-winsec`）、`logs/prefix-gone-full.txt`、`logs/postfix-gone-full.txt` |
| 同两形，但假 `go` 按参数解析（`FAKEGO_SCOPE_FILTER`，dir 形只认那一包、glob 形认子树） | 仍 **rc=0 / rc=0** | 仍 **rc=1 / rc=1** | `logs/seed-*.txt` 的 `*-filtered` 三发 |
| `--scope=winsec` 这一点 | rc=2 unknown | rc=0（载具，假 `go`） | `logs/seed-prefix.txt` / `logs/seed-postfix.txt` |
| `--scope=core` 少解析出 winsec（别人的既有守卫） | rc=1 `GUARD C - scope mode=core` | rc=1 同句 | 两份 seed 日志的 `winsec-missing-from-core`；⇒ 这枚改前改后都响＝**不是本票的缺口**，也证明本腿没把 C 的门放低 |

### 1.3 同一台件（`portable-tests-selftest.sh`）打在改前字节上的整排读数

`bash scripts/portable-tests-selftest.sh all /tmp/251-prefix-portable-tests.sh` ⇒
**`18 case(s) ran, 20 assertion(s) failed`、rc=1**（`logs/selftest-prefix-full.txt`）。
逐 case 退码：票 250 那 10 枚仍全部维持原判（`clean`/`coldcache-stderr` rc=0，其余 rc=1）；
本票新增的 `winsec-explicit-split-goes-red` 与 `winsec-explicit-pin-gone-goes-red` 在改前都是 **rc=0**（＝缺口本体），
`winsec-tier-*` 两枚是 rc=2（档不存在），`census-tier-roster-drift-goes-red` 报 `seed diff shows 0 line(s)`、`winsec-explicit-glob-is-what-bites` 报 `mutant replaced 0 line`——**种子在这份字节上没有可下的地方，载具拒绝把空种子当读数**（不是绿）。

### 1.4 反证：票面 (b) 支的**字面形**咬不住"拆包"（本腿跑了作用面才这么写）

题面 AC#1 的第二条路写作"GUARD C 在显式路径模式下**也**做一次钉-vs-解析对账"。本腿把这条例字面形跑了一遍：
变异＝只把显式分支里那行 `scope=("$winsec_scope")` 换成 `scope=("$@")`（钉照配、scope 仍用调用方给的目录），
同一颗拆包种子、开了按参数解析的假 `go`：

| 那一发 | rc | 读数 |
|---|---|---|
| 字面形（pinned 有、scope 是目录） | **0** | `runtests.sh: OK ... PASS=2`，`GUARD C` 一字未出 |
| 本腿交的形（pinned 有、scope 是档自己的 glob） | **1** | `GUARD C - scope mode=winsec resolved to a DIFFERENT package`、`Pinned: 1, resolved: 2` |

出处：`logs/seed-postfix.txt` 的 `counterfactual-noglob` / `real-bytes-glob-bites`，并已固化成载具的 case 18
`winsec-explicit-glob-is-what-bites`（真字节 rc=1 ＋ 变异体 rc=0 两枚断言）。
⇒ 结论：**只做"给显式路径配钉"不够，钉还得配 glob 形的解析集**；本腿因此两半都交（第五档 + 显式路径对账）。
⚠ 这条反证在**没有** `FAKEGO_SCOPE_FILTER` 的假 `go` 上是跑不出来的——原假 `go` 对任何参数都打印整份名册，
会让字面形"看起来也响"（本腿第一次跑就中了这个雷，见 §4 第 6 行）。

## 2. 逐格 AC

| 格 | 本腿选了什么／建了什么 | 读数出处 | 判语（归非实现者） |
|---|---|---|---|
| **AC#1** 钉咬得住 | **两条都走**：(a) `case $mode in`（`:201-245`）新增第五档 `winsec)`（`:229-236`），scope 用 glob `$winsec_scope=./internal/winsec/...`（`:197`）、`pinned=$winsec_pin`；(b) 显式路径块（`:328-349`）若拿到的正是该档自己的 scope（`./internal/winsec/` 或 `./internal/winsec/...`，且只有一个参数），就改按 `mode=winsec` 的档 scope + 钉跑 ⇒ 既有 GUARD C 比对块一字未动（锚点 `:353-368`，现 `:449-464`）就吃到了这一档。选型凭据＝§1.4 反证 | `logs/seed-postfix.txt`（四发显式路径种子）、`logs/postfix-split-full.txt` 与 `logs/postfix-gone-full.txt` 的逐字红句、载具 case 11-14/18 | |
| **AC#2** 尺本机可判 | 沿用同一台件，case 数 **10→18**，case 名同风格；两形正控各一枚（种假包名＝case 12/13/18，删真包名＝case 14）；⛔ 无〔仅 CI 可量〕字样。为让"解析集"这枚概念在本机真能判，假 `go` 加了 `FAKEGO_SCOPE_FILTER`（dir 形只认那一包、glob 形认子树、`./...` 认全模块） | `bash scripts/portable-tests-selftest.sh` ⇒ **18 case / 0 assertion failed / rc=0**（`logs/selftest-postfix.txt`）；改前同一台件 ⇒ 18 case / 20 failed / rc=1（`logs/selftest-prefix-full.txt`） | |
| **AC#3** 不动别人的形状 | 空 scope 硬退出：本腿未碰那段字节，载具 case 10 仍从工作树里**逐字切出那 4 行**并跑出 rc=2；unknown `--scope` 仍 `exit 2`，句子形状不变，只把 `winsec` 报进 known 列表（第五档真存在了，列表不报它才是假话）。GUARD A/B/C 的比对块、shape check、ledger、runtests 委托一行未动 | 真跑：`bash scripts/portable-tests.sh --scope=winsecfoo` ⇒ **rc=2**、`portable-tests.sh: unknown --scope=winsecfoo (known: core, windows, cli, winsec, census)`（`logs/postfix-scope-unknown.txt`）；`--scope=` ⇒ 同样 rc=2；载具 case 17 `unknown-scope-still-hard` rc=2 钉住整句、case 10 钉住 4 行切片；卫生尺见 §3 | |
| **AC#4** census 自证 | 档名从此**只写一处**：`tiers='core windows cli winsec census'`（`:192`）。census 在 `go list` 之前先把这枚名册与**本文件自己的** `case $mode in` 分支集合（`:201` 起，awk 抽标签，含 `*)` 兜底记作 `*`）双向作差，不等就**拒答**（rc=1，两张集合都打出来）；循环改读 `$tiers`，遇到名册里有名而循环里没钉的档名同样拒答（`census` 自己除外，它是审计者不持钉）。⇒ 选了"钉住"这一支，不是"登记为故意分两处" | 正向：真 `--scope=census` 仍 rc=0，且**与起手那份逐字对差只剩别人的一行**（`logs/census-postfix.txt` vs `logs/census-prefix.txt`，winsec 行仍是 `corewinsec`）；反向：载具 case 16（名册删 `winsec`、分支留着）⇒ **rc=1**，逐字 `census - the tier roster (tiers=) and the 'case $mode in'` ＋ `tiers= : * census cli core windows ` vs `branches : * census cli core windows winsec `（`logs/postfix-census-drift-full.txt`）；探针 `census-roster-drift` 同一发 rc=1 | |
| 禁区交代（票面 31 行：不许新增第六档而不交代另三档的钉是否也缺读者） | 本腿加的是**第五**档（`core/windows/cli/census` 之外没有第六枚）。另三档的钉今天都有读者：`core_pin`/`win_pin`/`cli_pin` 各自由 `case` 的对应支赋给 `pinned` ⇒ GUARD C，且四份名册都被 census 的循环读；`winsec_pin` 原来只有 census 一个读者，改后有三个读者（第五档支／显式路径对账／census） | `scripts/portable-tests.sh:201-245`（case 块）与 `:302-313`（census 循环）；改前唯一读者＝`:233` 那行 `winsec) pin=$winsec_pin ;;`（票面已述，本腿复认） | |

### 2.1 残余限界（本腿自己说，不等别人挖）

1. **显式路径的对账只在"参数恰为那一枚档目录"时点火**。`winsec-tests.sh` 允许在调试时追加别的包
   （`scope=("$@")`），一旦带上第二枚参数就退回"ad-hoc 列表＝没钉可核"的旧语义 ⇒ 该形下拆包不响。
   今天 CI 不这样点（`ci.yml:439` 无参 ⇒ `./internal/winsec/` 一枚），载具 case 13 钉的就是这一枚。
2. **core 档里的 winsec 行仍是目录形** `./internal/winsec/`（`:211`），未随第五档改成 glob。
   改它会动 core 的解析集（多出来的子包会被 core 测到），超出本票射程 ⇒ 写进 §5 请裁。
3. `--scope=winsec` 改后会往下起真测试，本腿**没跑过真档**（见 §3 末行），其绿只由假 `go` 载具代取。

## 3. 门禁读数

**载具**（本票唯一判据尺）
- 交件时：`bash scripts/portable-tests-selftest.sh` ⇒ `18 case(s) ran, 0 assertion(s) failed`、rc=0。
- 打在改前字节：同一命令 + `/tmp/251-prefix-portable-tests.sh` ⇒ `18 case(s) ran, 20 assertion(s) failed`、rc=1。
- 票 250 那 10 枚 case 的**断言一字未改**：`git diff -U3 f5f9cc3..HEAD -- scripts/portable-tests-selftest.sh` 的删除行只有表头两句注释（`nine cases` 那句过期计数），其余全是新增。

**卫生尺（最硬那把：`git diff f5f9cc3..HEAD` 的 hunk 数与位置）**

`scripts/portable-tests.sh`：**7 枚 hunk**（`+104 / -8`），左端全部落在锚点 69-255 之间——

| # | hunk 头 | 锚点射程 | 是什么 |
|---|---|---|---|
| 1 | `@@ -69,8 +69,11 @@` | 69-76 | Usage 注释：补 `--scope=winsec` 一行与显式路径的例外说明 |
| 2 | `@@ -101,6 +104,19 @@` | 101-106 | 档位说明注释：winsec 档的来龙去脉 + census 不是档 |
| 3 | `@@ -165,6 +181,21 @@` | 165-170 | `winsec_pin` 之后新增 `tiers=`、`winsec_dir`、`winsec_scope` |
| 4 | `@@ -195,12 +226,20 @@` | 195-206 | `case $mode in`：新增 `winsec)` 支；unknown 消息改读 `${tiers// /, }` |
| 5 | `@@ -217,6 +256,39 @@` | 217-222 | census 块开头插入 AC#4 名册-vs-分支双向对账 |
| 6 | `@@ -227,10 +299,19 @@` | 227-236 | census 循环 `for m in core windows cli winsec` → `for m in $tiers` + `census` 跳过 + 无名册项拒答 |
| 7 | `@@ -245,11 +326,26 @@` | 245-255 | 显式路径块：档目录识别 → 按档跑并对账 |

**8 枚删除行的锚点号逐一**（尺：`git diff -U3 f5f9cc3..HEAD -- scripts/portable-tests.sh` 里以 `^-` 开头且非 `---` 的行）：
`73`（Usage 注释一行）、`203`（unknown 那句 echo）、`230`（census 的 `for m in ...`）、`248/249/250`（显式路径三行注释）、`251/252`（`scope=("$@")`、`pinned=''`）。
⇒ **最大删除行号 252**；锚点上空 scope 硬退出块 `:254-257`、GUARD C 比对块 `:353-368`、GUARD A `:370-388`、GUARD B `:513-568`、shape check `:312-345`、ledger `:390-406` **均无一行进入删除集**＝一行未动。
`scripts/portable-tests-selftest.sh`：**3 枚 hunk**（`+218 / -2`；两枚删除行都在表头那句 `nine cases` 的旧计数上，10 枚既有 case 的断言一行未删）。`scripts/testdata/portable-tests/go`：**2 枚 hunk**（`+45 / -2`，list 分支与契约注释）。
本票射程内文件之外的东西：`git diff --numstat f5f9cc3..HEAD` 只有这三枚文件 + `.scratch/wisp/probes/251/**` + 票面尾部一节，⛔ 无别的。

**尺寸**（交件时 `wc -l -c` 实测）
```
668 38629 scripts/portable-tests.sh
529 24878 scripts/portable-tests-selftest.sh
143  6280 scripts/testdata/portable-tests/go
```

**占位符尺**：派单里那把字符类尺（四枚被写文件 + `verdict.md` + `seed-winsec.sh` + 票面尾部）⇒ **0 命中**（`grep` rc=1）。尺的字面文本本不写进被扫文件——第一版写了，被自己的尺咬到一行，已改（这条改动的凭据就是下面那次复跑的 rc=1）。

**本腿没跑过什么（不藏）**
- ⛔ 真 `go test`／`go build` 一发未跑：`--scope=core`／`windows`／`cli` 真档、改后的 `--scope=winsec` 真档、`scripts/winsec-tests.sh` 全链，都没跑
  （`cmd/wisp`+`internal/config` 被腿 `198-r1` 占、`internal/ball` 被腿 `33-r8b` 占，改后 `--scope=winsec` 会往下起真测试）。
- ✅ 跑过的真命令只有：`--scope=census`（只做 `go list`）、`--scope=winsec`（改前那发，取 rc=2）、`--scope=winsecfoo` 与 `--scope=`（取 rc=2）、载具与探针（假 `go`）。
- 因此第五档在**真 windows 主机上是否仍 rc=0** 属〔待验〕：本腿给的凭据是"glob 与目录形今天解析出同一枚单包"（§1.1 第 3 条，真 `go list`），
  不许据本腿读数排除任何候选形，也不许据本腿判 CI 已绿。

## 4. 本腿推翻／修正编排者题面哪几句

| # | 题面位置与原话 | 实测 | 证据 |
|---|---|---|---|
| 1 | `:11` "`case $mode in`（`:170`）只有 `core)`／`windows)`／`cli)`／`*)` **四支**" | **漏计 `census)`**：锚点那份字节里是 `core)`(171)／`windows)`(184)／`cli)`(191)／`census)`(198)／`*)`(202) ⇒ **4 枚具名档 + 1 枚兜底，共 5 支** | awk 抽取尺（§1.1 第 4 条）⇒ `4` 枚具名 |
| 2 | `:15` "winsec 那族测试真正进 CI 的路径是 `scripts/winsec-tests.sh:97`" | 该行在 **`:94`**；`:97` 是 `ran=$(count '^=== RUN')`（同一枚台账 `docs/reports/pending-and-issues.md:10505` 也写的是 `:97`，本腿不改台账，只报此处） | `grep -n 'bash "$portable"' scripts/winsec-tests.sh` ⇒ `94:` |
| 3 | `:10` "`winsec_pin` 今天**唯一的读者**是 `:229-242` 那圈 census 循环" | "唯一读者"复认成立；**行号范围偏小**：那圈 while 循环是 `:226-242`（`:226` 是 `while IFS= read -r p; do`，`:229` 已在体内） | 锚点那份字节 220-243 行 |
| 4 | `:9` "`:164-166` 钉了一枚 `winsec_pin`，内容＝1 行" | **复认成立**（1 枚包名，逐字一致） | 锚点 `:164-166`；`sed -n '164,166p' \| grep -c github.com` ⇒ `1` |
| 5 | `:23` AC#1 判据"种一枚'winsec 包名被改／少一行'的异常 ⇒ 指名那一步必须红"，且 (b) 支写作"显式路径模式下**也**做一次钉-vs-解析对账" | **判据不足以覆盖票面自己 19 行的风险句"拆出第二包"**：只把钉配到显式路径、解析集仍用调用方的目录形，`go list ./internal/winsec/` 看不见旁侧新包 ⇒ 同一种子 rc=0 绿。本腿跑了那条例字面形的作用面（§1.4）才这么写，并已把反证固化成载具 case 18 | `logs/seed-postfix.txt` 的 `counterfactual-noglob`（rc=0）vs `real-bytes-glob-bites`（rc=1） |
| 6 | 隐含前提（本腿第一发种子踩到的）："假 `go` 的输出来自名册文件"＝可当解析读数 | **不成立**：改前的 `scripts/testdata/portable-tests/go` 对**任何** `go list` 参数都打印整份名册，dir/glob 之分它不建模 ⇒ 用它跑出来的"拆包会响"是仪器自己造的假象。本腿加了 `FAKEGO_SCOPE_FILTER` 后两形才分得开 | 原 shim 的 `list)` 分支（锚点 `:46-67`，不看参数）；现 shim `list)` 在 `:49`、按参数解析那块在 `:62-105` |
| 7 | `:3` 立票锚点 `db06a399` | 本腿起手锚是 `f5f9cc34`（票 250 结案后 dev 上又落了 census(114-a2) 等笔）；所有"改前"尺一律打在 `f5f9cc34` 那份字节上，⛔ 未拿 `db06a399` 当改前 | `git log -1`（§0） |

## 5. 判不动的地方（未定义即停，本腿一律没自填）

| # | 事项 | 甲＝谁补得上 / 乙＝补不上 |
|---|---|---|
| 1 | **要不要把 `scripts/winsec-tests.sh` 改成点具名档**（`bash "$portable" --scope=winsec`，从而让 GUARD C 走第五档那条更直的路）。本腿写面**不含**该文件，故只做了"显式路径也认档"这一半；改过去才是把两条路并成一条 | 甲：任何拿到 `scripts/winsec-tests.sh` 写面的腿（一行改动 + 载具一枚 case）。本腿没做＝写面外，不是判断它不该做 |
| 2 | **第五档该不该进 `ci.yml`**（现在 CI 只 `run: bash scripts/winsec-tests.sh`（`:439`），第五档没有 CI 调用者）。workflow 形状属契约级 ⇒ 本腿**停手未动**，只在交件里说明 | 甲：编排者（人工拍板）。⛔ 不许把"第五档没人点"读成"本票没落地"——被点的显式路径那一半已经响了 |
| 3 | **core 档里的 `./internal/winsec/` 要不要也改成 glob**（`:211`）。改成 glob 会把"日后拆出的第二包"算进 core 的解析集，进而**要求 `core_pin` 同步加行**，那是动 core 档的钉＝射程外 | 甲：编排者裁；若裁"改"，同笔必须同 commit 改 `core_pin` 并在 CI log 交代（脚本自己那段 GUARD C 的话头就是这个要求） |
| 4 | **AC#3 那句 unknown 消息的"一字不许软"如何算**：本腿把 known 列表从 4 名扩到 5 名（`winsec` 现在真是档），句子形状、`>&2`、`exit 2` 未动。若编排者裁"列表也不许动"，则第五档与这句不能同时成立 ⇒ 需要人裁其一 | 甲：编排者。本腿交的凭据：`logs/postfix-scope-unknown.txt` 逐字 + 载具 case 17 |
| 5 | **census 自读自身文本**这枚手法：本腿用 `$here/portable-tests.sh` 再 awk 抽 `case $mode in` 的分支标签。它对"分支标签写在第 0 列"这一形状敏感（将来缩进或改写风格⇒ 抽取集合变小⇒ census 拒答＝响亮失败，不静默绿）。若嫌脆，可改成"档名册由 case 块生成"的另一种形，本腿没选是因为那要把 `case` 拆成数据表＝动别人射程 | 乙（可接受）：留响亮失败。甲（要换形）：编排者定 |
| 6 | 本腿未跑真 `go test`，故**第五档真机 rc 读数缺失**（§3 末行）。谁补：编排者在干净工作树上 `bash scripts/portable-tests.sh --scope=winsec` 与 `bash scripts/winsec-tests.sh` 各一发 | 甲：编排者（或有真档权限的腿） |

## 6. 交件套数清单

| commit | pathspec | numstat |
|---|---|---|
| `2ca152d7` | `.scratch/wisp/probes/251` | 5 files, +188（骨架 + 四份起手日志） |
| `709d2589` | `scripts/portable-tests.sh scripts/portable-tests-selftest.sh` | 162/2 载具，102/8 脚本 |
| `4cca2c0e` | `scripts/portable-tests.sh` | 11/9（census 对账提到 `go list` 之前） |
| `aa0dbb69` | `scripts/portable-tests-selftest.sh scripts/testdata/portable-tests/go` | 63/7 载具，45/2 假 go |
| `ad18cbeb` | `.scratch/wisp/probes/251 .scratch/wisp/issues/251-*.md` | 16 files +817/-44（44 枚删除全在 `verdict.md` 的骨架版上；**票面 44 added / 0 deleted**） |
| `59c01e88` | `.scratch/wisp/probes/251/r1/verdict.md` | 2/3（两枚空种子 case 的逐字话头、§4 行号指向） |

证据件：`.scratch/wisp/probes/251/r1/{verdict.md,seed-winsec.sh,logs/}`（日志 **17 份**：起手 4 份〔census-prefix / prefix-scope-winsec / prefix-scope-unknown / selftest-baseline〕、改后 4 份〔census-postfix / postfix-scope-unknown / selftest-postfix / postfix-census-drift-full〕、改前整台件 1 份〔selftest-prefix-full〕、探针 2 份〔seed-prefix / seed-postfix〕、逐字红绿 4 份〔prefix/postfix × split/gone-full〕、**作废 2 份**〔`logs/variant-b-*.txt`＝本腿第一次跑反证时还没给假 `go` 装按参数解析，那两发的 rc 不可当凭据，留着只为记录踩坑〕）。
临时件一律只建不删（各 `mktemp -d` 影子根与假 go 工作目录留在 `/tmp` 下）。
