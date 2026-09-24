# 141 / Q-46c 对抗验收表 r1（裁决者 ≠ 实现者）

* **被验锚点**：`bb61dc5`（`2026-09-24 22:32 +0800`，`dev`）。裁决程 **只读快照、不改一字节生产码**。
* **纯净树取法（所有"整树读数"只在这里量）**：
  `mkdir -p /d/tmp/141-acc-base && git archive bb61dc5 | tar -x -C /d/tmp/141-acc-base`
* **变异树**（全部在 `/d/tmp` 下，**不在仓里、不删**）：`141-acc-parent`（`470e6c5`＝`ed2c077` 的父）、
  `141-acc-mut-strlen`、`141-acc-mut-lexical`、`141-acc-mut-noband`、`141-acc-mut-witharrow`、
  `141-acc-mut-backseed`、`141-acc-rev-schema`、`141-mut-noex`、`141-acc-probe`（③④样本）。
* **共享树纪律**：`internal/observe/sampler_settle_gate_136_test.go`（在飞程）与 `design/**`
  的 16 枚未提交删除＋未跟踪 `design/old/`、`design/doubao/` 全程**未读作事实、未 stage、未还原、未删**。
  本表每节只 `git add` 这一枚文件路径。

## 0. 被验的四枚 commit（`git show --name-only` 原文，逐枚复量"只带自己点名的路径"）

```
ed2c077 22:17  d22scan(ban #8): 面向用户的字符变严、注释面豁免、并同时补上原本漏掉的数学符号段 (票 141, Q-46 已批＝按推荐)
M	tools/d22scan/main.go
M	tools/d22scan/scan_test.go            （+469/−30：main.go 197、scan_test.go 302）
1218192 22:21  risk(票 141 具名解冻 ③): 审批卡 R7 文案那枚 `≥` 换成 ASCII，生产串与期望值同批 (Q-46 已批＝按推荐)
M	internal/risk/assessor_test.go
M	internal/risk/rules_scale.go
2970c79 22:21  memory(票 141 自决档): DDL raw 串里那两处"Go 看是字符串、SQL 看是注释"按严格读法清成 ASCII
M	internal/memory/schema.go
bb61dc5 22:32  evidence(141 Q-46c 交件表): 五条判据逐条读数＋判据(3)在盘上不可达那一维＋§10 对上 A197/Q-48
A	docs/evidence/s1/141-q46c-impl.md
```
⇒ **四枚都只带自己点名的路径，无越界、无删除（`--name-status` 全为 M/M/M/A）**。
`gofmt -l tools/d22scan internal/risk internal/memory` ＝ **空输出**；`go vet ./`（`tools/d22scan`）＝ **静默**。

⚠ 但 `git diff --stat 470e6c5..bb61dc5` 是 **9 枚文件／8 枚 commit**，不是 4 枚：
`054faae`、`f1b9c1c`、`e21d3bf`、`e848a93` 四枚 docs/evidence 夹在本批中间（`HANDOVER.md`、
`pending-and-issues.md`、`136-instr-fixes-r1.md`）。**不是越界**（不属于被验四枚），
只是"整批区间"与"本票四枚"不能混为一谈，读差异表时要点开。

## 1. 门禁基线四数（锚点快照现跑，`-v`，未动任何阈值）

```
$ sh tools/d22scan/runtests.sh -C tools/d22scan ./...      # 在 D:/tmp/141-acc-base
rc=1   PASS=22  FAIL=2  SKIP=0  === RUN=64  '^panic:'=0
--- FAIL: TestScannerSelfScanOfRealRepoIsGreen (1.53s)
--- FAIL: TestRealRepoLedgerIsHonest (0.91s)
```
* **名册差集**（防"一枚 panic 吞掉同包读数"）：`scan_test.go` 里 top-level `func Test*` **24 枚**，
  跑出 `--- PASS/FAIL` 的 top-level **24 枚**，`comm` 两向差集 **皆空** ⇒ 没有测试被吞。
* 整树真扫（同一快照，独立于测试）：`go run . -root /d/tmp/141-acc-base` ⇒ **rc=1／6 finding**，
  逐枚＝`frontend/fixtures/composer-states.html:2,5,8`、`frontend/src/components/ai-native/thinking.tsx:213`、
  `.../tool-chips.tsx:186`、`frontend/src/components/composer.tsx:170`；`design/`=0、`internal/`=0、`cmd/`=0。
  各 scope 计数：`design/ 16、frontend/ 40、internal/ 405、cmd/ 39`。
* 码点现量（我自己逐字节取，不是抄表）：前 4 枚里 `composer-states.html:2,5,8` 与 `composer.tsx:170`
  ＝ **U+2264**，`thinking.tsx:213`、`tool-chips.tsx:186` ＝ **U+2212**。**简报那句前提成立。**
* **旧正则不覆盖 ⇒ 新增红**：`git show ed2c077^:tools/d22scan/main.go | sed -n '111p'` 的字符类是
  `[1F000-1FAFF 2600-27BF 2B00-2BFF FE0F 1F1E6-1F1FF]`，与新类逐字符比对**只差 `\x{2200}-\x{22FF}` 一段**
  （箭頭 `2190–21FF`／带圈 `2460–24FF` 仍不在内）。反向实测见 §4(a)：把新段摘掉，
  `TestScannerSelfScanOfRealRepoIsGreen` 与 `TestRealRepoLedgerIsHonest` **立刻复绿**（6 枚 finding 归 0）
  ⇒ 这 6 行确实是本批点亮的，不是存量。

### 那两枚红＝"门能看见"，不是放水（简报点名要明写的一句）

* 两枚红是**同一个因**：`:269` 逐行点名那 6 枚（`t.Errorf`×6），`:1277` 是
  `t.Fatalf("HEAD must be green, rc=1")`。`:1273` 那句 `ban #6 examined 43` 是 `t.Logf`，
  **不是断言**，别读成"台账数不对"；且 `git show ed2c077^:...scan_test.go` 里它就已经是 `t.Logf`
  （`git log -L 1273,1274` 指回 `84e4161`，票 88），**不是本批降下来的**。
* 全仓**没有一处**为了变绿放宽断言：`ed2c077` 的测试 diff 里被删的断言行 **0 枚**、新增 **10 枚**（§5）。
* CI 侧后果（现量，不是推演）：`scripts/d22scan.sh` 是 `set -eu`，第一步就是这个正控 ⇒
  **正控红 ⇒ 第二步整树扫描根本不跑**，`.github/workflows/ci.yml:109` 调的正是它。
* **若 owner 在 `Q-48` 选"退回射程"（摘掉 `2200–22FF`）会不会自动复绿**：
  **这两枚会自动复绿**（`141-acc-mut-noband` 实测 PASS，见 §4(a)），
  **但另有 3 枚会红**：`TestBan8MathBandAndRemainingGaps`（U+2265/2264/2229/2212 四行）、
  `TestBan8CommentExemptionInGoSources/raw_string_that_reads_as_a_SQL_comment`、
  `TestBan8CommentExemptionInTextScopes/JSX_text_node_with_a_math_glyph`。
  ⇒ **"退回射程"不是一行 revert**：它要同时改回具名解冻②要求钉下的那三枚钉子，
  并让 `1218192`/`2970c79` 的清存量失去门禁支撑（生产串改成了 `>=`，却再没有任何仪器钉住"不许回到 `≥`"）。
  这一格只报形状，**不替 owner 选**。

## 2. 裁决格 1：豁免有没有吞掉本该算违规的东西（首要攻击点）

### (1) 变异①：摘掉"注释抹成空格"这一手

`141-mut-noex`：只把 `main.go:886` 的 `ranges := commentRangesFor(path, string(src))` 换成
`var ranges map[int][][2]int`，字符类**不动**。两把尺分别量：

**A. 交付射程（含 `2200–22FF`）＋ 锚点树 —— 实现方自述"78 枚注释行／35 枚文件由绿翻红"这一句，在交付射程上不成立：**

| 尺 | 整树红行 | `internal/`＋`cmd/` 被豁免放掉的行／文件 |
|---|---|---|
| 交付射程＋有豁免（＝锚点原样） | **6** | — |
| 交付射程＋**无**豁免 | **7** | **1 行 / 1 枚文件** |

那唯一一枚（逐字节取码点）：
`internal/risk/pathresolver_anchor_spelling_windows_test.go:200` ＝ **U+2229 `∩`**，
原文 `// TestAListWinsWhereBothTablesHit is AC#3's A∩B shape: ...`，是 `_test.go` 的一行 **doc 注释**。
⇒ **78 对不上交付射程，1 才对得上。**

**B. 票面 PLAN 字面宽射程（`1F000–1FAFF / 2190–2BFF / FE0F`）＋ 改动前父树 `470e6c5`（"清存量之前"）—— 实现方的数逐位对得上：**

```
$ /d/tmp/141-acc-bin/d22-wide-noex.exe -root D:/tmp/141-acc-parent | grep -c '\[emoji\]'   → 141
$ /d/tmp/141-acc-bin/d22-wide-ex.exe   -root D:/tmp/141-acc-parent | grep -c '\[emoji\]'   → 44
名册（我自己的 python 集合差，不读它的 .lst）：
  放掉总数 = 97 行；其中 internal/+cmd/ = **78 行 / 35 枚文件**
  反向：ex − noex = **0 行** ⇒ 豁免没让任何一行"新"变红（不是换个方向漏）
```
78 枚文件级名册（本程现量，35 枚）：
`cmd/wisp/{doctor.go,leg_sink_gate_131_test.go,run.go}` 各 1；
`internal/agent/approval/ticket87_veto_l2_test.go`1、`ticket97_alias_direction_test.go`3、
`internal/agent/{control.go,inject.go,spill.go}` 各 1；`internal/audio/{doc.go,gate.go}` 各 1；
`internal/ball/{dock.go,position.go}`1、`live_windows_test.go`2、`position_test.go`1、`tokens_test.go`10；
`internal/config/{migrate.go,parse.go,unwired.go}`1、`unwired_test.go`1；
`internal/observe/sampler_settle_coverage_136_test.go`2；`internal/panel/composer_handlers.go`2；
`internal/risk/pathresolver_anchor_spelling_windows_test.go`1、`pathresolver_expansion_test.go`1、
`provenance.go`8、`provenance_test.go`4、`syncdirs_other_test.go`1；`internal/statemachine/table.go`1；
`internal/tools/bridge_a18_kill_windows_test.go`8、`bridge_junction_windows_test.go`2、`fs_staging_windows_test.go`4、
`recycle_windows.go`1；`internal/winsec/reparse_windows_test.go`7、`seam_guard_windows_test.go`1、
`winsec.go`2、`winsec_windows.go`2 ⇒ **合计 78 行／35 枚文件，与自述逐位一致。**

**这一支的判**：数**不是它数错的**（宽射程下 78/35 真），但 `ed2c077` 的 commit message 把它写成
"**现射程**＋豁免前＝141 行红；现射程＋豁免后＝44 行红 ⇒ …78 枚注释行…" —— **标签错**：
141/44/78 三个数只在**未被批准的宽射程**下存在。同一 message 里那句
"the band without (c) is a **78-line comment rewrite**" 对**交付仪器是假的**（真值＝1 行，见 A 表）。
它自己的证据表 §3 逐字标的是"**宽射程**"⇒ **表对、message 错**，以表为准即可，但这枚差异必须落进本表。

### (2) 变异②：反过来把字符串也抹掉（过度豁免）⇒ 哪一枚现有测试会红？

`141-acc-mut-strlen`：在 `goCommentRanges` 里额外用 `ast.Inspect` 把 `*ast.BasicLit`
（`token.STRING`/`token.CHAR`）的字节一并抹成空格。全套：

```
$ sh tools/d22scan/runtests.sh -C tools/d22scan ./...     # 在 141-acc-mut-strlen
rc=1  PASS=19  FAIL=5  SKIP=0  === RUN=64
--- FAIL: TestScannerSelfScanOfRealRepoIsGreen        （锚点本来的 2 枚红，非新增）
--- FAIL: TestRealRepoLedgerIsHonest                  （同上）
--- FAIL: TestEmojiBanCoversGoSourcesNotJustDesign    ← 新增
--- FAIL: TestBan8MathBandAndRemainingGaps            ← 新增
--- FAIL: TestBan8CommentExemptionInGoSources         ← 新增
```
⇒ **答得出，而且答了三枚**："字符串不许抹"这一半**有钉子**，不是只靠 doc 注释自证。
（`TestBan8CommentExemptionInTextScopes` 不红是**对的**——文本射程没有"字符串"概念，
它红反而说明我抹错了东西。）

### (3) 变异③：`//` 出现在 Go 字符串字面量里

样本（`141-acc-probe`，10 枚以上 production `.go` 过 `checkRoot`）：
`internal/probeA/raw.go`、`internal/probeB/rawurl.go` —— raw 串**内部行首**是 `//` 且带 U+2265/U+2264。

```
$ d22-base.exe -root D:/tmp/141-acc-probe      → probeA/raw.go:4 [emoji]、probeB/rawurl.go:4 [emoji]   ← 正确判红
$ d22-strlen.exe（②的过度豁免）                → 两枚都被抹掉＝绿                     ← 证明样本真踩在豁免面上
$ d22-noex.exe（无豁免）                       → probeA/probeB/probeC/probeD 四枚全红
```
`go/ast` 路线**没有被字符串里的 `//` 骗到** ⇒ 行为正确。**但**：
我把 `.go` 的豁免从 `go/ast` 换成词法行首规则（`141-acc-mut-lexical`：`commentRangesFor` 对 `.go`
也走 `textCommentRanges`），probeA/probeB **立刻转绿**，而全套只有 3 枚 subtest 红——
`TestBan8CommentExemptionInGoSources` 的 `one-line_block_comment`／`trailing_comment_after_clean_code`／
`unparseable_file_gets_no_exemption_at_all`，**没有一枚是因为"raw 串行首 `//`"这个形状红的**。
⇒ **缺口（照实报，本程不补）**：③ 的形状在 `scan_test.go` 里**没有用例**。
现有的 10 形里唯一沾"raw 串看起来像注释"的是 `-- L1 slot， ≤ 20 xing`（SQL 破折号那一形），
它钉住的是"`--` 不算标记"，钉不住"`//` 不算标记"——而后者才是 `go/ast` 相对词法规则的**独有**价值。

### (4) 变异④：解析失败的 `.go` 里放一枚真违规

样本 `internal/probeC/broken.go`（`func broken( {` 让它必挂，违规串在下一行）：

```
$ d22-base.exe -root D:/tmp/141-acc-probe
internal/probeC/broken.go:1: [unparseable] ... expected ')', found '{' (and 1 more errors)
internal/probeC/broken.go:4: [emoji] ban #8 glyph in scope internal/ is banned (D23): ...
```
⇒ **读不出注释没换来免检**，行为正确。**这一支有用例钉着**：
`TestBan8CommentExemptionInGoSources/unparseable_file_gets_no_exemption_at_all`（`wantFindings: 2`），
且它在 ② 与"词法替换 go/ast"两支变异下都红（变异②里它仍红＝豁免整体失效也拦不住它）。
本表 §1 的真扫也复核了同一形状在盘上确实存在（锚点树里 0 枚未解析 `.go`，`[unparseable]` 只在样本里）。

### 本格判：**附条件成立**

* 成立的部分：豁免**没有吞掉本该算违规的东西**——②③④三支实测都站在正确一侧，
  且"字符串不许抹／解析失败不免检"都有测试钉（②红 3 枚、④红 1 枚用例）。
* 条件一（**必须记，不必修**）：③ 的"raw 串行首 `//`"形状**无用例**，
  把 `go/ast` 换成词法行首规则**不会**因为这个形状报警（只因另外 3 个形状报）。
* 条件二：`ed2c077` 的 message 把 **宽射程**的 141/44/78 说成"**现射程**"，
  并据此论证"补段不豁免＝78 行改写"；交付射程下的真值是 **1 行／1 枚文件**。
  数字本身可复现，**标签与推论不可复现** ⇒ 按 `AGENTS.md` 头三句"以权威文件为准"处理为记缺陷，不算伪报。

## 3. 裁决格 2：`go/ast` 那条"raw string 里的 SQL `--` 不算注释"——单点回退试验

**做法**：只把 `internal/memory/schema.go` 一枚文件回退到含 `≤` 的旧写法
（`141-acc-rev-schema` ＝ 锚点全树 ＋ `git show 470e6c5:internal/memory/schema.go` 覆盖那一枚），
别的一字不动。

```
$ sh tools/d22scan/runtests.sh -C tools/d22scan ./...      # 在 141-acc-rev-schema
rc=1  PASS=22  FAIL=2  SKIP=0  === RUN=64
--- FAIL: TestScannerSelfScanOfRealRepoIsGreen
--- FAIL: TestRealRepoLedgerIsHonest
    scan_test.go:269: repo HEAD violates: internal/memory/schema.go:29: [emoji] ban #8 glyph in scope internal/ ...
```
独立复核（同一份回退树，直接跑二进制，三把尺）：

| 仪器 | `internal/memory/schema.go:29` | 整树 finding 数 |
|---|---|---|
| `d22-base`（交付，走 `go/ast`） | **红** | 7（6 frontend ＋ 这枚） |
| `d22-lexical`（词法行首规则替换 `go/ast`） | **红** | 7 |
| `d22-strlen`（②的过度豁免：字符串也抹） | 绿 | 6 |

**答得出"哪一枚测试会红"**：`TestScannerSelfScanOfRealRepoIsGreen`（`TestRealRepoLedgerIsHonest` 同因）。
⇒ **`2970c79` 那枚清理的 `≤` 那一半不是装饰**——它是"HEAD 必须绿"这道门的承重件。
另一重钉子也在：`TestBan8CommentExemptionInGoSources/raw_string_that_reads_as_a_SQL_comment`
（`wantFindings: 1`，用的就是这一形），它在格 3(a) 摘段变异下会红（见 §4），
所以"raw 串里的 `--` 不外"这一判既有**整树门**钉、又有**单元用例**钉。

**`→`（U+2192）那一半另判**：把 `:118` 一并回退后**零枚测试红**（箭头段按批复不扫）。
这一条 `2970c79` 的 message 自己写了（"属批复明确保留的缺口段，本 commit 清它纯粹是'顺带收紧'，
不清也不会红"）⇒ **自述与盘上一致，不算虚报**；但它确实只是装饰，不计入承重。

**必须记的一条措辞超额**：`commentRangesFor`/`walkEmoji` 头上那句
"`.go` 走 `go/ast`（唯一能分辨 raw string 里的 SQL `--` 不是 Go 注释的权威——
`internal/memory/schema.go:29` 就是这一形）"，**对这一枚样本被实测反驳**：
词法规则（`d22-lexical`）在同一枚上给出**同一个红**，因为它根本不认 `--` 是标记。
真正只有 `go/ast` 能救、词法会漏的形状是"**raw 串行首 `//`**"（§2(3) 的 probeA/probeB），
而那一形**没有用例**、且盘上今天 0 处。⇒ 判：**结论对（豁免只认 ast、不认词法 `--`）但举的证据样本不承重**，
换一枚真能区分的样本或给 probeA 那形补用例，二选一即可。

### 本格判：**成立**（附一句要改的证据文案，见上）

## 4. 裁决格 3：正反两向有没有钉住"保留缺口"（`TestBan8MathBandAndRemainingGaps`，两支都真跑）

钉的是同一枚 `emojiRe`，反方向只改了**一行**（字符类），其余一律原样。

### (a) 摘掉 `\x{2200}-\x{22FF}` 那一段（`141-acc-mut-noband`）⇒ **红**

```
rc=1  PASS=21  FAIL=3  SKIP=0  === RUN=64
--- FAIL: TestBan8MathBandAndRemainingGaps
      ↳ subtest: U+2265 greater-or-equal / U+2264 less-or-equal / U+2229 intersection / U+2212 minus sign  (4 枚)
--- FAIL: TestBan8CommentExemptionInGoSources      ↳ raw_string_that_reads_as_a_SQL_comment
--- FAIL: TestBan8CommentExemptionInTextScopes     ↳ JSX_text_node_with_a_math_glyph
```
同一次运行里 `TestScannerSelfScanOfRealRepoIsGreen` 与 `TestRealRepoLedgerIsHonest` **变成 PASS**
（6 枚 `frontend/` finding 归 0）⇒ 这既是"正控有效"的证明，也是 §1 那句"6 行是本批新增红"的反向证明。
**AC#3 的反证成立**：那 4 枚种子是被**新段**抓到的，不是被邻段顺手抓的。

### (b) 把箭头段 `2190–21FF` 也加进去（`141-acc-mut-witharrow`）⇒ **应该红，实测红**

```
rc=1  PASS=21  FAIL=3  SKIP=0
--- FAIL: TestBan8MathBandAndRemainingGaps
      ↳ subtest: U+2192 right arrow, gap kept / U+21D2 rightdouble arrow, gap kept
--- FAIL: TestScannerSelfScanOfRealRepoIsGreen
--- FAIL: TestRealRepoLedgerIsHonest
```
"补箭头的代价"现量（同一份 witharrow 真扫，锚点树新增）：**`internal/` 新增 16 行红**
＝ `internal/tools/fs_write.go`×4（票面 §61 点名的"外露到步骤账本与 `Result.Text`"那一族）＋
`fs_staging_windows_test.go`×8 ＋ `ticket97_alias_direction_test.go`×2 ＋
`bridge_a18_kill_windows_test.go`×1 ＋ `fs_write_test.go`×1。⇒ 保留缺口**不是空口说的**，
它一边被 `wantFired:false` 钉着、一边有可数的代价。

### (c) 追加一支（简报没点名，但"两支都要真跑"只做箭头段会漏掉带圈/制表符那一族）

把整段 PLAN 字面射程 `2190–2BFF` 都加进去（`141-acc-mut-wide`）：
`TestBan8MathBandAndRemainingGaps` 的 **4 枚负向行全红**
（U+2192、U+21D2、**U+2500 box drawing**、**U+2460 circled one**），正向 7 行仍绿；
同一把尺真扫锚点树＝ **40 finding**（对 `470e6c5` 宽射程＋豁免是 44，差 4 枚正是本批清掉的
`schema.go:29`、`schema.go:118`、`rules_scale.go:24`、`assessor_test.go:154`）。
⇒ 保留缺口的**三族**（箭头／制表符／带圈数字）都在钉子覆盖内，实测一致。

### 本格判：**成立**

正反两向都真跑、都按预期一侧红一侧绿，且我额外把"带圈数字／制表符"两族也真跑了一遍。
**附带一枚小缺陷（只记不修）**：`scan_test.go:348` 的注释写着"U+2460/**U+2461** (circled numbers)
are still unscanned"，但表里只有 U+2460 一行——注释多报了一枚**没有用例**的码点。
与格 6 同族（话比仪器宽），方向相反（这里是"承诺的缺口"被多列了一枚，不是"覆盖面"被多报）。

## 5. 〔占位〕裁决格 4：断言有没有被放宽、helper 是不是原有的



