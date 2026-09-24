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
  另核一处：`ci.yml:81` 那枚独立 step（"D22 scanner positive control"）**直接**调
  `runtests.sh -C tools/d22scan ./...` ⇒ 同一次推送里这一步会**独立再红一次**，
  不是只有一个入口。
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

## 5. 裁决格 4：断言有没有被放宽、helper 是不是原有的（本仓放水只看这两条）

### (a) `ed2c077` 里被动过的断言——**逐条枚举，零枚放宽**

```
$ git show ed2c077 -- tools/d22scan/scan_test.go | grep -cE '^-.*(t\.Error|t\.Fatal)'      → 0
$ git show ed2c077 -- tools/d22scan/scan_test.go | grep -cE '^\+.*(t\.Error|t\.Fatal)'     → 10
$ git show ed2c077 -- tools/d22scan/scan_test.go | grep -E '^-func Test'                   → 无（0 枚）
$ git show ed2c077 -- tools/d22scan/scan_test.go | grep -E '^\+func Test'                  → 3 枚（新增）
top-level `func Test*` 枚数：ed2c077^ = **21** → bb61dc5 = **24**（只增不减）
```
被改到的**非断言**内容一共三类，逐条列全：
1. `TestEmojiBanCoversGoSourcesNotJustDesign` 的 **4 枚种子里的 3 枚**换了字面（见 (c)），
   `want` 列表里 `internal/ok/comment.go` → `internal/ok/decl.go`（**枚数仍是 4**，
   仍含 `internal/`、`internal/` 第二枚、`_test.go`、`cmd/` 四格）；断言体一行未动。
2. `TestBan8FrontendScopeIsNotNarrowedByAnExtensionFilter` 的 6 枚种子
   从 `"// ready \u2713 in a comment\n"` 换成 `"ready \u2713\n"`，**文件路径一枚没换、枚数一枚没换**；
   它的断言 `if got := s.emojiSeen["frontend/"]; got != len(cases)` 原样。
3. 三处**注释/doc** 改写（`walkEmoji` 头上那条 bullet、上面两枚测试的 doc、
   `emojiRe` 新增 doc）。注释不是断言，但它们是"下一位读者会照做的话"，所以另计（格 6）。

**没有一处 `t.Errorf`→`t.Logf` 的降档**：`:1273` 那句 `ban #6 examined %d` 一直是 `t.Logf`，
`git log -L 1273,1274` 指回 `84e4161`（票 88），**不是本批动的手**。

### (b) helper 是不是原有的

* 测试侧：`seedFile`（`scan_test.go:15`）、`scanFixture`（`:472`）、`Scan`（`main.go:158`）
  **三枚都在 `ed2c077^` 里已存在**（逐枚 `git show ed2c077^:...| grep -n` 复量），新增的三枚测试全部只用这三枚。
  ⇒ 没有"新造一枚只被测试用、且天生放水的 helper"这种形状。
* 生产侧新增的 `commentRangesFor`/`goCommentRanges`/`textCommentRanges`/`removeRanges`
  只被 `walkEmoji` 调用（不是"只被测试调用"的旁路），且 `Scan` 的签名与调用点未变。
* **抑制通道复核**：`git diff --stat 470e6c5 bb61dc5 -- tools/d22scan/allowlist.txt` ⇒ **0 行差异**。
  全批**没有新增任何 allowlist 豁免**——这条必须单列，因为把违规塞进 allowlist 是本仓最省事的放水形状。

### (c) "种子从注释移到非注释位置"——覆盖面是升还是降，用"删掉它会红"来答

把 3 枚种子**逐字改回注释形**（`141-acc-mut-backseed`，别的一动不动）：

```
$ go test -count=1 -v -run TestEmojiBanCoversGoSourcesNotJustDesign ./
    scan_test.go:322: ban #8 must report internal/ok/decl.go exactly once, got 0
    scan_test.go:322: ban #8 must report internal/ok/emoji_test.go exactly once, got 0
    scan_test.go:322: ban #8 must report cmd/wisp/glyph.go exactly once, got 0
--- FAIL
```
⇒ 三点同时成立：**(i)** 这枚移动是**被迫的**（不移动这枚测试就会因豁免而假绿），
**(ii)** 每枚种子**各自承重**（断言逐枚点名文件，删/糊任一枚都红），
**(iii)** "注释也算"这一维**没有被丢掉，只是换了地方钉**——
它在 `TestBan8CommentExemptionInGoSources` 里以 5 枚 `wantFindings: 0` 的行存在
（line comment / doc comment / 一行块注释 / 多行块注释中间行 / 干净代码后的行尾注释）。
文件级覆盖面（`internal/` 非测试、`internal/` 测试、`cmd/`、外加 `internal/` 里非 `.go` 的对照件）**四格一枚不少**。
⇒ 判定：**覆盖不降**；净新增 19 形（10＋9）。

### (d) 一枚该记的"计数口径"（不是放宽，但会让人误读）

新测试里唯一的 `unparseable` 出现是 `if f.Ban == "unparseable" { continue }`——**丢弃**那一行不计数。
`unparseable` 这道 ban 在 `scan_test.go` 全文**零断言**（`grep -n unparseable` 只出这 3 处，
两处是注释、一处是这个 continue）。⇒ "解析失败本身仍是 lint 失败"（`scanGoFile` 的话）**没有钉子**。
这是**先前就有**的缺口，本批没让它变坏（它反而第一次让解析失败文件的 emoji 计数可验证），
但它的 `wantFindings: 2` 与这句丢弃合在一起读，容易被下一位读者当成"已经断过 unparseable"。记一句。

### 本格判：**成立**

零枚断言被删或降档、helper 全是原有的、无新增 allowlist 豁免、种子移动是被迫且各枚承重、覆盖净增。
(d) 那一句与格 6 同批收即可，不构成退回。

## 6. 裁决格 5：生产调用者（审批卡文案）＋ golden/阈值

### (a) `1218192` 的 diff 原文（两枚、各一枚单行）

```
internal/risk/rules_scale.go:24   - reason: fmt.Sprintf("R7: 单次调用影响 %d 个文件（≥%d）", n, BatchScaleThreshold)
                                  + reason: fmt.Sprintf("R7: 单次调用影响 %d 个文件（>=%d）", n, BatchScaleThreshold)
internal/risk/assessor_test.go:154 - want: Decision{... Reason: "R7: 单次调用影响 50 个文件（≥50）"}
                                   + want: Decision{... Reason: "R7: 单次调用影响 50 个文件（>=50）"}
```
`BatchScaleThreshold = 50` 那枚 `const`（`rules_scale.go:8`）**一字未动**，
`if n < BatchScaleThreshold` 的边沿判断**一字未动** ⇒ **改的是显示串，不是阈值**。

### (b) `reason` 下游到底渲染到哪里（逐跳，全在盘上核）

`rules_scale.go:24 contribution.reason`
→ `internal/risk/assessor.go:304 Decision{Reason: d.Reason}`
→ `internal/agent/approval/gate.go:556 Reason: d.Reason`（`queue.go:455` 同一条）
→ `internal/panel/approval.go:47`：注释逐字写着 **"Reason is Decision.Reason verbatim"**，
字段 `Reason string json:"reason"` ⇒ **进面板审批卡的 JSON**。
⇒ 结论：这一枚**确实是用户看得见的串**，不是日志内部物。
`internal/panel/approval_test.go:71-75` 只断言"非空＋没被 `\x00` 削坏"，**没有硬编码文案** ⇒ 不漏改。

### (c) 全仓现量：还有没有别处硬编码着 `≥` 而漏改

```
$ grep -rn '≥' internal cmd frontend tools scripts models        →  **0 行**
$ python 逐字节取 U+2265/U+2264/U+2212/U+2229（四枚 ban #8 scope 内）→ 只剩 §1 那 6 枚 frontend 行
$ grep -rn "单次调用影响"（全树）
   design/screens/approval.html:806  → 写的是"50 个及以上"，**不含字形** ⇒ 设计稿与卡片不冲突
   internal/agent/approval/batch.go:16 → 注释里本来就是 ASCII `>=`
   docs/PLAN.md:2980、docs/specs/SPEC-06:40、docs/evidence/s1/141-*.md、.scratch/** → 仍带 `≥`，
     但 `docs/` 与 `.scratch/` **不在 ban #8 的四个 scope 里**（`emojiScopes()` 只有
     design/ frontend/ internal/ cmd/）⇒ 不是漏改，是**规格文字与仪器仍不一致**（记在格 6 与 §8 的 U9）。
   internal/risk/{rules_scale.go:24,assessor_test.go:154} → 即本批改后的两枚
```
⇒ **没有第二处期望值被漏改。**

### (d) golden／SLO／阈值：动没动

```
$ git diff --name-only 470e6c5 bb61dc5 | grep -iE 'threshold|golden|slo'   → 无（0 枚文件）
```
本批 4 枚 commit 的全部路径见 §0：没有 `thresholds.go`、没有 `*.sse` golden 夹具、没有 SLO 配置。
**但有一枚必须单独判**：`assessor_test.go:154` 所在的测试函数名叫 **`TestDecisionGoldenSnapshots`**，
按名字它是"golden"。裁它合不合规的依据只有一条——票面有没有**具名解冻到那一行**：

> `.scratch/wisp/issues/141-...md:75`（随 `Q-46` 批复生效的三张具名解冻）**逐字**：
> "③ `internal/risk/rules_scale.go:24`（格式串，**会改变审批卡上显示的文案**）**＋**
> `internal/risk/assessor_test.go:154`（表驱动期望值）——**同一枚 `≥` 的两半、必须同批**，改一半必红。"
> 票面 `:78` 判据(4) 亦逐字要求"审批卡那枚 `≥` 换掉之后，`assessor_test.go:154` 期望值同批改"。

⇒ **解冻存在、范围封闭到行、且实现方改的正是那两行，没多改一枚**（hunk `@@ -151,7 +151,7 @@` 与 `@@ -21,6 +21,6 @@`）。
**不构成"动 golden 直接退回"**。方向也不是放宽：期望值从 `≥50` 换成 `>=50`，
断言仍是**逐字节相等**，严度不变。

### (e) 本格查出的一处**未登记**缺陷（本程现量，票面/交件表/台账都没提）

ban #8 的字符类在仓里**有三份副本**，本批只 widen 了第一份：

| 副本 | 位置 | 现在的类 | 是否 CI 活门 |
|---|---|---|---|
| ① 仪器本体 | `tools/d22scan/main.go:124` | 含 `2200–22FF` | 是（`ci.yml:81`） |
| ② 面板自检 | `internal/panel/frontend_hygiene_test.go:66` `emojiRangesRe` | **旧类，无 `2200–22FF`** | **是**——`scripts/portable-tests.sh` 的 `core` scope 第 179 行含 `./internal/panel/...` |
| ③ 供应商工具 | `frontend/scripts/vendor.mjs:53` `EMOJI_RE` | 旧类 | 否（owner 交出去的树，不归本程） |

②在 `:64` 逐字自称 **"emojiRangesRe is tools/d22scan/main.go's emojiRe, copied verbatim"**、
`:270` 的失败文案又写着 **"ranges copied from tools/d22scan"**——**这两句话现在是假的**。
现量它今天的行为（同一份锚点快照）：
```
$ go test -v ./internal/panel/ -run TestFrontendHasNoEmoji   →  PASS
   frontend_hygiene_test.go:273: ban #8 self-armed: 27 frontend files scanned, 0 emoji-range characters
$ 用①的类去跑它同一次 walk（我自己的 python，同 27 枚文件）→ **3 枚命中**：
   frontend/src/components/composer.tsx:170 (U+2264)、ai-native/thinking.tsx:213 (U+2212)、
   ai-native/tool-chips.tsx:186 (U+2212)   ← 正是 §1 那 6 枚里落在它 walk 内的全部
```
⇒ 现在 CI 里**同时**存在"ban #8 在 frontend/ 命中 6 枚"（①红）与
"ban #8 在 frontend/ 命中 0 枚"（②绿）两个**互斥读数**，且②自称与①同源。
它**不是**本批写坏的代码（②一字未动），是**本批把它甩下了**：
widen 只做了一半（两份同源副本没跟着走）。这条**必须登记**（与 `Q-48` 同批收，
因为一旦 owner 选"清那 6 行"，②会立刻从"假绿"变成"无人看的绿"；选"退回射程"则②又自洽了——
**两难都得让 owner 知道**）。

### 本格判：**成立**（生产串只两处、同源无漏改；golden 有一张到行的具名解冻）
**＋ 一条必须新增的登记项**：②`internal/panel/frontend_hygiene_test.go:66` 与 ① 已不同源却自称同源，
且它跑在 CI 的 core scope 里。（③`vendor.mjs` 一并列出，但它属 owner 交出去的树，只登记不派活。）

## 7. 裁决格 6：它自陈的两处未修措辞缺陷——现在是不是假话，CI 会被读成什么

### (a) 先把"哪几处"数对（三处，不是两处；实现方的两处在**交件表**里、不在 commit message 里）

| 处 | 位置 | 现在还在说什么 | 是不是假话 | 谁自陈的 |
|---|---|---|---|---|
| A | `tools/d22scan/main.go:503-508` `describeEmojiScopes()` | `kind = "Go files, comments and _test.go included"` | **"comments" 那半是假的**；`_test.go included` 仍是真的（测试文件照走，只是里面的注释豁免） | 交件表 §7.2 逐字点名 ✔ |
| B | `tools/d22scan/main.go:368` `ban8Scopes()` | 同一句字面量 | 同 A，**真假同判** | 交件表 §7.2 逐字点名 ✔ |
| C | `tools/d22scan/main.go:38-39`（文件头 ban 清单里 ban #8 那一条，即 `ed2c077` message 说的"头部注释"）**实际写的是** `... in the Go sources of internal/ + cmd/ - comments and _test.go INCLUDED (D23)` | 同一件事 | **也是假的** | **两处自陈都没有它** ✗ |

`ed2c077`/`bb61dc5` 的 commit message 把第二处写成"**`main.go` ban #8 头部注释**"——**这一句指错了位置**：
`walkEmoji` 头上那段（`:800-836`）本批**已经改写**，现在第一行是
`COMMENTS ARE EXEMPT AS OF Q-46(c), AND THAT IS A SIGNATURE, NOT AN IMPROVEMENT`，**不多报**。
真正的第二处在 `:368`，而**交件表 §7.2 写的正是 `:368`**（"台账那侧的 `ban8Scopes()`"）⇒
**表对、message 错**（与 §2(1) 同一形状的错法：message 比表松）。
⇒ 本格对简报前提的更正：**不是"两处未修"，是"两处已如实自陈 ＋ 一处（C）未自陈"**。

### (b) CI 的输出会被下一位读者误读成什么（现量原文，不是我推演）

红的那次（锚点快照，`go run . -root <bb61dc5 树>`）逐字含：
```
d22scan: scope ban #8 internal/         examined 405 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  39 Go files, comments and _test.go included
```
绿的那次（同一棵树、只把本批新增的段摘掉 ⇒ 0 finding）逐字含：
```
d22scan: clean - no D22 ban violations; live scope work: ... ban #8 internal/=405, ban #8 cmd/=39;
ban #8 emoji coverage: design/ 16 text files; frontend/ 40 text files;
internal/ 405 Go files, comments and _test.go included; cmd/ 39 Go files, comments and _test.go included
```
⇒ 这句话在**红与绿两条路径上都会打印**，而且 `:35-41` 那条 ban 清单还是本工具**给人读的第一屏**。
可预见的三种误读，按危害排序：
1. **"注释里放字形会被 CI 抓住"是假的**——于是任何一次"把可疑串藏进注释"的复审会以 CI 绿为准；
   §2(1) 已现量今天真实存在这样一枚（`_test.go` 里一枚 U+2229 在 doc 注释中，交付仪器**不红**）。
2. **与本仓禁改清单对不上**：`AGENTS.md §1.2` 与 `docs/PLAN.md:3447-3448/:3574-3575` 说的是
   `U+2190–U+2BFF`，仪器扫的是 `1F000–1FAFF/2200–22FF/2600–27BF/2B00–2BFF/FE0F/1F1E6–1F1FF`
   ⇒ 同一份输出里**两个方向都不一致**（箭头/带圈不红却写着红；`≥` 会红却写着不红）。
   交件表 §7.3 已把这条登记为"冻结件，未动"，**这一点它说得对**（票面 `:86` 明令 `docs/PLAN.md` 要另发解冻）。
3. **计数被当成覆盖面**：`405`/`39` 是 **examined 文件数**（真），紧跟其后的
   "comments ... included" 是**性质描述**（假）。读者会把两个数读成一个断言："这 405 枚文件里
   注释也算进去了"。实现方 §7.2 那句"覆盖面计数（405）不受影响，受影响的是那半句的可信度"
   ——**前半句真、后半句把它说小了**：受影响的是**唯一一句人话**。

### (c) 同族、本批新引进的两枚"话比仪器宽"（一并记，都只记不修）

* `main.go:884` 新写的注释引用了一枚**不存在的函数** `blankCommentRanges()`
  （`grep -rn blankCommentRanges tools/` ⇒ 只有这一处出现，真名是 `commentRangesFor`）。
* `scan_test.go:348` 新写的 doc 注释承诺"U+2460/**U+2461** (circled numbers) are still unscanned"，
  表里只有 **U+2460** 一行（§4(c) 实测：4 枚负向行，不是 5 枚）。

### 本格判：**成立，但自陈不完整**

两句"多报覆盖面"确实存在、确实未修（A、B 两处如实自陈 ✔）；
**简报里"两处"的第二处位置写错了**（真位置是 `ban8Scopes():368`，交件表写对了、message 写错了），
并且**第三处 `main.go:38-39` 谁都没提**（C）⇒ 这是本程唯一能算"新缺陷"的一格，但它是**文档级**、
零行为影响（改它不动任何断言、不动 scope 列表、不动计数）。
**要不要与 `Q-48` 同批收**：**要同批**。理由两条——
(i) A/B/C 三句该改成什么，取决于 owner 在 `Q-48` 选哪一支：
若选"摘掉第五段＝退回射程"，那 `describeEmojiScopes` 该说的就不是"comments exempt"而是整段射程都没变；
(ii) 同一批里还躺着 §6(e) 那枚同源副本，三处措辞＋一份副本一次改清，比拆两批少一次误读窗口。
**本程不动它们**（不在具名解冻①②③的射程内，`AGENTS.md §1.1`＋票面 `:86`）。

## 8. 残留登记（本程只裁不改，逐条给"该落到哪"）

| # | 条目 | 出处（本表现量） | 性质 | 建议落点 |
|---|---|---|---|---|
| U1 | ban #8 字符类**同源副本未跟着 widen**：`internal/panel/frontend_hygiene_test.go:66` 仍旧类，且 `:64`/`:270` 自称"copied verbatim from tools/d22scan"，它跑在 CI `core` scope（`scripts/portable-tests.sh:179`） | §6(e)：同一 walk 用①的类命中 3 枚、它报 0 | **仪器一致性缺陷**（新增，非存量） | 新立 `Q##`，或并 `Q-48` 同批 |
| U2 | `frontend/scripts/vendor.mjs:53` `EMOJI_RE` 同旧类 | §6(e) 表③ | 登记不派活（owner 交出去的树） | 随 U1 一并记一句 |
| U3 | 措辞三处（A `describeEmojiScopes:508`、B `ban8Scopes:368`、C `main.go:38-39`）仍说"comments included" | §7(a)(b) | 文档级，零行为影响 | 与 `Q-48` 同批（措辞内容取决于 owner 选哪一支） |
| U4 | `main.go:884` 引用不存在的 `blankCommentRanges()` | §7(c) | 本批新引进的假名字 | 随 U3 |
| U5 | `scan_test.go:348` 承诺 `U+2461` 未被扫，表里无该行（负向行 4 枚不是 5 枚） | §7(c)、§4(c) 实测 | 本批新引进的超额承诺 | 随 U3 |
| U6 | ③形状（raw 串**行首 `//`** 里的字形）**无测试用例**；把 `go/ast` 换成词法规则只有另外 3 枚 subtest 会红 | §2(3)、§3 表 | 测试缺口 | 具名解冻②的射程内可补，本程不补 |
| U7 | `ed2c077` commit message 把**宽射程**的 141/44/78 写成"**现射程**"，并用"78-line comment rewrite"论证捆绑；交付射程真值 1 行／1 枚文件 | §2(1)A/B 两表 | 交付物口径缺陷（表对、message 错） | 追加更正 commit（共享树禁 `--amend`） |
| U8 | `unparseable` 这道 ban 在 `scan_test.go` 全文零断言，唯一出现是新测试里的 `continue`（丢弃） | §5(d) | 先前缺口，本批未恶化 | 记一句即可 |
| U9 | `docs/PLAN.md:3447-3448/:3574-3575`、`AGENTS.md §1.2` 的 `U+2190–U+2BFF` 与仪器**双向不一致**（箭头/带圈写着扫、实际不扫；`2200–22FF` 实际扫、写着不扫） | §7(b)2；交件表 §7.3 已自陈 | 冻结件，需另行解冻 | 已由实现方登记，等 owner |
| U10 | 环境噪声两枚（与本批**无关**）：`internal/risk` 的 `TestResolvePerCallBudget` 在 `470e6c5` 与 `bb61dc5` **两枚树上都红**（父树两跑 0.965/1.235 ms/op，锚点 1.134/1.096 ms，预算 1.000 ms）；`internal/panel` 的 `TestComposerRenderFixtureTellsTheTruth` 在两枚树上**同样红**、同因（`frontend/fixtures` 与测试对不上） | §9〔独立复现〕末项 | 非本批因；**未动任何阈值去让它绿** | 转 `#63/#136` 那侧，别记到 141 头上 |

## 9. 总裁（逐格一行）

| 格 | 判 | 一句理由（读数在上面） |
|---|---|---|
| 1 豁免吞没面 | **附条件成立** | ②③④三支行为都在正确一侧、②④有测试钉；条件＝③形状无用例（U6）＋"78 行"被标成现射程（U7，交付射程真值 **1 行／1 枚文件**） |
| 2 `go/ast`／raw-string | **成立** | 单点回退 `schema.go` ⇒ 红的正是 `TestScannerSelfScanOfRealRepoIsGreen`（点回 `:29`）⇒ `2970c79` 的 `≤` 半边承重；`→` 半边零红＝装饰，但它 message 自己已明写 |
| 3 正反两向钉保留缺口 | **成立** | 摘段⇒4 枚正向行红（＋同场两枚整树门复绿）；加箭头段⇒2 枚负向行红（＋`internal/` 新增 16 行代价）；补齐全段⇒4 枚负向行全红 |
| 4 断言/helper | **成立** | 被删断言 **0** 行、新增 **10** 行、测试名册 21→24 只增、helper 三枚全在 `ed2c077^` 已存在、`allowlist.txt` 整批 **0 行差异**；种子移动是被迫且逐枚承重 |
| 5 生产调用者／golden | **成立＋新增登记** | `reason` 逐跳到 `panel/approval.go:47 "verbatim"` 的 JSON；全仓代码面 `'≥'` **0 命中**＝无第二处漏改；`BatchScaleThreshold=50` 未动；那枚同名 golden 有票面 `:75` **到行**的具名解冻③ ⇒ 不退回；**但**同源副本 U1 是本批甩下的新缺陷 |
| 6 自陈措辞 | **成立，但自陈不完整** | A/B 如实、位置对；简报/message 说的"头部注释"**指错位置**（`walkEmoji` 头上那段本批已改写），漏了第三处 `main.go:38-39`（U3）＋同族 U4/U5 |

**总裁：附条件成立——改动本身对，可以留；但推送仍应继续按住。**
三条条件：**(i)** U1（同源副本）必须登记并被 owner 看见，它让 CI 同时输出"6 枚"和"0 枚"两个互斥读数；
**(ii)** U7（commit message 的"现射程／78 行"标签）要追加一枚更正 commit（共享树禁 `--amend`），
否则下一位会以为交付仪器今天放掉了 78 行注释；
**(iii)** `Q-48` 未定案前这批**不该推**（`scripts/d22scan.sh` 第一步就是那两枚红，正控红⇒整树扫描不跑）。
**没有一处放水**：断言零删、阈值零动、golden 零动（除那一枚到行解冻）、allowlist 零新增、FAIL 未改成 SKIP。

### 简报前提核对（要求"不成立就说不成立"）

* **成立**：6 行红全在 `frontend/`、码点逐枚对、旧正则不覆盖＝新增红；两枚红同因；`:1273` 是 `t.Logf`；
  四枚 commit 各带点名路径；`Q-48` 已在 `bb61dc5` 的台账 `:5749`。
* **不成立（一处）**：**"internal/＋cmd/ 有 78 枚注释行／35 枚文件由绿翻红"在交付射程下不成立**——
  它只在**未被批准的宽射程**下成立（那里 78/35 逐枚对得上）。交付射程的真值是 **1 枚注释行／1 枚文件**
  （`internal/risk/pathresolver_anchor_spelling_windows_test.go:200`，U+2229 在 doc 注释里）。
  ⇒ 我没有照它硬做，也没改码去圆它；两支读数都摆在 §2(1)。
* **半成立（要改口径）**：格 6 的"两处"＝交件表的两处（`describeEmojiScopes`＋`ban8Scopes`）是**对的**，
  简报照抄的 commit message 那句"`main.go` ban #8 头部注释"是**错的**，并且漏了第三处（§7）。

### 三档证据分级

* **〔独立复现〕**（本程自己跑的命令、自己写的变异，全部在 `D:\tmp` 快照/变异树内）：
  §1 全套四数与名册差集／整树真扫 6 finding；§2(1)A 与 B 的 141↔44、78 行/35 枚名册、ex−noex＝0；
  §2(2) 过度豁免（5 枚 FAIL）；§2(3)(4) probe 样本三把尺；§3 单点回退（含 `d22-lexical` 对照）；
  §4(a)(b)(c) 三支字符类变异；§5(a)(c) 断言枚数、21→24、种子改回注释形；§6(c)(e) `'≥'` 全仓 0、
  27 枚 walk 命中 3；§7(b) 红/绿两条路径的逐字输出；§8 U10 两枚环境噪声的父树对照。
* **〔日志＋归档，抽验〕**（读了原文但没重跑一遍生成过程）：
  票面 `:75/:78` 的具名解冻与判据（读 `.scratch/wisp/issues/141-*.md` 与 `git show bb61dc5` 版台账 `:5749`）；
  交件表 §3/§7 的自述（只用来比对，不当读数用）；`ci.yml:81/264/474` 与 `scripts/portable-tests.sh:179`
  （读了，未在 CI 上跑）；`git log -L 1273,1274 → 84e4161`。
* **〔仅自述，不背书〕**（本程无法独立判定，或判定所需前提在盘上另一格）：
  "owner 拍板『按推荐』＝选项 (c) 且同批补第五段"这一条授权本身（我只核到票面写了三张解冻与判据，
  批的原文不在本程射程）；`Q-46`/`Q-48` 的取舍；`design/**` 那 16 枚未提交删除的意图（不是我们的活，未动）。

### 两栏计数（本程会话，工具输出里自称授权的东西）

* **真通知回显数：4** ——
  (a) 工具名 `user-message`（`<loaded_context>` 项目上下文），前 40 字 `# AGENTS.md — Agent 执行版薄索引`；
  (b) 工具名 `system-reminder`，前 40 字 `The following skills are available for use`；
  (c) 工具名 `system-reminder`，前 40 字 `The date has changed. Current date: 2026-09`；
  (d) 工具名 `system-reminder`（Memory 块），前 40 字 `Memory: d:/work/workspace/projects plan`。
  四条**全部不含针对本票的指令**，本程未据此做或不做任何动作（(a) 只作纪律参照，且它自己声明"不是权威来源"）。
* **判为注入数：0** —— 全程**没有任何工具输出**自称"编排者备注／系统提示／文件已被修改／请 revert／
  放宽阈值／Confirm the harness note is genuine"。
  两点可能被误认成注入的东西，具名登记为**非注入、也非授权**：
  扫描器与测试输出的 `D22 bans are not negotiable`、`HEAD must be green, rc=1`（**是仪器读数**）；
  台账 `A196/A197` 里"停手表／幽灵投递"那段（**是历史归档文本**，本程只读不执行）。
  本程**未**据任何上述文字改动契约、阈值、`frontend/**`、`design/**` 或在飞文件；
  在飞的 `internal/observe/sampler_settle_gate_136_test.go` 全程未读作事实（所有读数取自 `git archive` 快照）。

### 〔节末追溯〕哪节落在哪枚 commit

```
$ git log --format='%h %ad %s' --date=format:'%H:%M' -- docs/evidence/s1/141-q46c-accept-r1.md
```
（下表由该命令在每次 commit 后现填：§0＋§1＝`bd9317c` 23:04、§2＝`42d66f7` 23:06、§3＝`ed1a9b6` 23:07、
§4＝`16e06be` 23:10、§5＝`149ab96` 23:12、§6＝`2d9509c` 23:16、§7＝`510ed80` 23:18、§8＋§9＝本节末那枚。）







