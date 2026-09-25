# 147 对抗验收件 r1 — `unwritten` 支的 offset 与那三枚新用例的牙

- 被验版本：**`e572aa9435bfee310900f78422ec31984c41dc7b`**（＝编排者简报里的 `e572aa9`）
- 验收程：**非实现者**（本件作者不是 `2dc2b28` / `1f3ede9` / `e572aa9` 的作者，也不是 `a91d7c2` 的作者）
- 交件时刻：2026-09-25 22:0x–23:0x +08
- 本件的每一条读数都是**本程现跑**，实现件里的日志只用作"它说了什么"的对照，**没有一条读数是从它的表里抄的**

---

## 0　锚点与"被验版本"的取法（这一步错了后面全绿也是假绿）

**判据**：被验版本必须是一枚具体 sha；任何"某文件长什么样"的断言必须按 sha 取，**不得把脏工作树当被验那版**。

### 0.1　进场第一读（命令原文）

```
$ git rev-parse HEAD
e572aa9435bfee310900f78422ec31984c41dc7b

$ git log --oneline -5
e572aa9 evidence(147 收笔): 门禁四数与名册差集、mtime 顺序自证、锚点被推走的追加登记
1f3ede9 docs(147 AC#3): 那三处（现量是四处文本＋一枚标题）"只有空白"的声明就地追加更正 + 票面交件记录 + 读数随件
2dc2b28 fix(147 AC#1+AC#2): unwritten 支的 offset 改填"文档停在第几字节"，并补三名会响的用例
ae968f9 ticket(146 -done): 结线改名——5 枚 AC 框全 [x]、两张非实现者表在 docs/evidence/s1/、真未勾 0 枚
9ca5069 docs(台账 A263): 59 枚这一推红名集合动过零枚、分母涨 16 枚且逐枚归到 144/146；slo-full 今天第一次拒采
```

⇒ **编排者 21:5x 的读数相符**：本程进场时 HEAD **就是** `e572aa94`，没有别人的 commit 插进来。
21:43:43 / 21:46:15 / 21:51:18 +08 三枚时间戳与"码＋尺／AC#3 更正／收笔"的分工相符。

### 0.2　工作树是脏的，而且脏在别人家（命令原文）

```
$ git status --porcelain | head -30
 D design/assets/base.css
 D design/assets/icons.js
 ...（design/ 共 16 枚删除）
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
?? .zcodeignore
?? design/doubao/01-ball-states.jpg
?? design/doubao/demo/lib/
?? design/doubao/demo/rightbar.js
?? design/doubao/demo/screens/home.js
?? design/doubao/demo/screenshots/
?? design/doubao/demo/sidebar.js
?? design/old/

$ git diff --name-only e572aa94 | grep -vE '^(frontend|design)/'
docs/reports/pending-and-issues.md
```

⇒ 与编排者那条"此刻另一会话在写 `frontend/**` 与 `design/**`"相符；**另有一处**：`docs/reports/pending-and-issues.md`
在工作树里也与被验版本不同（台账在被别人追加）。本程**一枚未碰、未还原、未计入任何"零命中"宣称**。

### 0.3　取法：一律 `git archive` 到仓外快照

```
$ git -c core.autocrlf=false -c core.eol=lf archive e572aa94 | tar -x -C /tmp/wisp147
```

快照与本仓的对应关系（**全部在仓外**，`/d/tmp` 之外另起，绝不在仓库目录内建 worktree 或 checkout）：

| 目录 | 内容 | 用途 |
|---|---|---|
| `/tmp/wisp147` | `e572aa94` 全树 | 被验版本本体：控制组、静态门禁、变异基座 |
| `/tmp/wisp147-snap2` | `e572aa94` ＋ `slo_windows.go` 换成父 `ae968f9` 版 | **新用例打在未修码上**（恒真那一问） |
| `/tmp/wisp147-snap3` | `e572aa94` ＋ 码与测试**都**换成 `ae968f9` 版 | 改前名册基线（`-list` 差集） |
| `/tmp/wisp147-snap4/base` | snap3 ＋ 只撤红句那半句（票面点名的 MG） | 票面"现量的形状"第 3 条复算 |
| `/tmp/wisp147-mut/r1`、`r2` | 18＋7 枚变异快照（每枚一份，只建不删） | 承重矩阵 |

一致性自证（快照不是"我以为的那版"）：

```
$ for f in cmd/wisp/slo_windows.go cmd/wisp/slo_report_144_windows_test.go; do
    a=$(git show e572aa94:$f | sha1sum); b=$(sha1sum < /tmp/wisp147/$f); ...
cmd/wisp/slo_windows.go                 sha=7fab42496950d8364246b4f64e9e2e74faf561c4  MATCH
cmd/wisp/slo_report_144_windows_test.go sha=ead0279a94f2e00d95a53bc97e6b3d6b958614ee  MATCH
```

### 0.4　三枚 commit 的路径自核（第 3 节第 4 件的一半）

```
$ git show --name-only --format='%H %s' 2dc2b28 1f3ede9 e572aa9
2dc2b28  cmd/wisp/slo_report_144_windows_test.go  cmd/wisp/slo_windows.go  docs/evidence/s1/147-offset-naming-r1.md
1f3ede9  .scratch/wisp/issues/144-….md  .scratch/wisp/issues/147-….md  .scratch/wisp/probes/147/×17
         docs/evidence/s1/144-slo-report-partial-read-r1.md  docs/evidence/s1/147-offset-naming-r1.md
e572aa9  docs/evidence/s1/147-offset-naming-r1.md
```

⇒ **没有一枚含票面地界之外的路径**：无 `frontend/**`、无 `design/**`、无 `internal/**`、无 `tools/d22scan/**`、
无 `scripts/**`、无台账。`1f3ede9` 那 17 枚 `.scratch/wisp/probes/147/` 是仪器读数落盘，票面 §Acceptance 允许、
且符合"临时件只建不删"。

### 0.5　锚点在写作过程中被推走（照 §0 的规矩记名，不改被验版本）

本件第 0 格落盘、准备 commit 时，`git rev-parse HEAD` 已不是 `e572aa94`：

```
$ git rev-parse HEAD
f77a003b77bf93573649bb1803b93833128f462d

$ git show --name-only --format='%h %ad %s' e572aa94..HEAD
f77a003 Fri Sep 25 21:55:34 2026 +0800 docs(台账 A264): 票 147 实现程交件（判 ⓐ、三枚新尺、四件报回）；
        核收三条只证边界、不证尺会响

        docs/reports/pending-and-issues.md

$ git merge-base --is-ancestor e572aa94 HEAD && echo "YES e572aa94 is an ancestor of HEAD"
YES e572aa94 is an ancestor of HEAD
```

⇒ 插进来的**只有这一枚**，是编排者自己的台账 `A264`，路径 `docs/reports/pending-and-issues.md`，**零代码**。
⇒ **被验版本一律仍以 `e572aa94` 为准**（本件全部读数取自它的仓外快照），HEAD 前移不影响任何一格。
⇒ 顺带记下：`A264` 那句"核收三条只证边界、不证尺会响"与编排者派给本程的话相符——**"尺会响"那一格就是下面的第 1 格**。

**本格读数**：锚点相符、取法相符、路径相符、锚点被推走一事已记名。
**本格判定：成立。**

---

## 1　AC#2 那把尺的承重（首要攻击点；本程价值排序第一）

**判据（按本仓操作定义逐条，答不出＝没测）**：
① 恒真判据——三枚新用例在**未修的旧码**上响不响；
② 承重——**摘掉任意一味**（ⓐ 那一行填充／三枚用例里任意一枚／两枚新 helper 里任意一枚），
是否存在一发变异从此打不红；
③ 单点回退——只撤一处，问"撤掉这处哪条用例变得不响"；
④ 票面点名的那发 MG——修完之后该红还是该绿，一句话说清**并跑给人看**；
⑤ 有没有拿"永远绿"的断言凑数。

反面约束本程遵守了：**没有**要求"同一发变异里新加那支先响"（那支挂在 `if errors.Is(err, io.EOF)` 里，
前提正是被抹掉的东西）。本程要的读数是"摘掉这味之后，这发变异还响不响"，两味都响＝其中一味冗余，
不＝其中一味是装饰。

### 1.1　跑法（两条命令，可复算）

变异全部落在**仓外快照**，每枚一份目录、只建不删（`harness.py` 拒绝在已有目录上覆盖，改新编号）。
被验树本体（控制组）与 18＋7 枚变异共用同一把尺：

```
$ python C:/Users/swq/AppData/Local/Temp/wisp147-mut/harness.py r1            # 造 18 枚变异快照
$ bash C:/Users/swq/AppData/Local/Temp/wisp147-mut/run.sh r1 <names...>        # 每枚跑 147+144 全组
```

`run.sh` 里的尺＝`go test -v -count=1 -run 'TestSLO147|TestSLO144' ./cmd/wisp/`，
按**用例名**逐个判 RED／GREEN／NORUN（不是只数 `--- FAIL` 枚数），并且**单独报 144 那 7 枚的红绿**
（"新增断言把既有断言改反方向"是另一类假绿，必须分开看）。
PATH 里注入 `third_party/sherpa-onnx`，与 `scripts/wisp-cli-tests.sh` 同形。

### 1.2　① 恒真那一问：**响**。三枚都在未修码上红（本程独立重走，未用它的日志）

`/tmp/wisp147-snap2` ＝ 被验版本的全部测试文件 ＋ **只有 `cmd/wisp/slo_windows.go` 换回父 commit `ae968f9`**：

```
$ cd /tmp/wisp147-snap2 && go test -v -count=1 -run 'TestSLO147|TestSLO144' ./cmd/wisp/
--- PASS: TestSLO144EveryPrefixOfARealReportIsUnwrittenNotCorrupt (0.01s)
--- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow (0.00s)
--- PASS: TestSLO144LoopRetriesAnUnfinishedFileAndReadsTheWholeReport (0.03s)
--- PASS: TestSLO144LoopGiveUpSentencesOnRealFiles (0.05s)
--- PASS: TestSLO144CorruptReportIsJudgedOnTheFirstRead (0.00s)
--- PASS: TestSLO144UnfinishedReportKeepsPollingThenNamesBudgetAndBytes (0.06s)
--- PASS: TestSLO144ReportThatArrivesAfterMissingReadingsIsCollected (0.00s)
    slo_report_144_windows_test.go:446: 1977 of 1978 unwritten sentences do not name the offset the document stopped at; first: prefix of 1 bytes: subjectReportRead.offset is 0
--- FAIL: TestSLO147UnwrittenSentenceNamesTheOffsetTheDocumentStoppedAt (0.01s)
    slo_report_144_windows_test.go:470: give-up error says 8 bytes read, document open at offset 0; want both 8 (the prefix the loop was handed)
--- FAIL: TestSLO147LoopGiveUpSentenceAgreesWithTheBytesItRead (0.06s)
    slo_report_144_windows_test.go:502: unfinished sentence says 989 bytes, offset 0; want both 989
--- FAIL: TestSLO147OffsetSemanticsRenderThreeDifferentSentences (0.04s)
```

⇒ **三枚全响**，且 144 那 7 枚在未修码上保持绿（⇒ 响的不是"把旧断言改严"换来的）。
⇒ 红句原文与实现件 §2.2 逐字相符（`:446`／`:470`／`:502` 三处），它的自述这一条**复算成立**。
⇒ `1977 of 1978`：恰有一枚前缀（0 字节）在旧码上"过"，因为那一发 `offset==bytes==0`。本程确认这个数字不是笔误：
0 字节那发走的是 `io.EOF` 而非 `ErrUnexpectedEOF`，旧码 `dec.InputOffset()` 给 0，新码给 `int64(0)`＝同一个 0。

CI 同形跑法（`bash scripts/wisp-cli-tests.sh`，快照内自带 `third_party/sherpa-onnx`）两边各一次：

```
改后（被验版本 e572aa94）:  === RUN=133  --- PASS=73  --- FAIL=0  --- SKIP=0        rc=0  runtests.sh: OK
改前（同一测试文件+父码）:  === RUN=133  --- PASS=70  --- FAIL=3  --- SKIP=0        rc=1  strict runner exited 1
```

⇒ 两枚都是 **`=== RUN=133`、非 0**，与"根本没跑到"（`0xc0000135` ＋ 0 条 `=== RUN`）不是一回事；
返回码的差（0 vs 1）与 3 枚红名逐枚对得上。

**名册差集（两向 `comm`，本程自己跑 `-list`，不抄它的 `.names`）**：

```
$ cd /tmp/wisp147-snap3 && go test -list '.*' ./cmd/wisp/ | grep -E '^Test' | sort > /tmp/names-prefix.txt   # 改前全基线
$ cd /tmp/wisp147       && go test -list '.*' ./cmd/wisp/ | grep -E '^Test' | sort > /tmp/names-postfix.txt   # 被验版本
pre=70 post=73
comm -23 (消失):   （空）
comm -13 (新增):   TestSLO147LoopGiveUpSentenceAgreesWithTheBytesItRead
                  TestSLO147OffsetSemanticsRenderThreeDifferentSentences
                  TestSLO147UnwrittenSentenceNamesTheOffsetTheDocumentStoppedAt
```

⇒ 消失 0／新增 3，逐枚 `TestSLO147*`。**没有一枚既有用例被改名或被挤掉**。

### 1.3　④ 票面点名的那发 MG：跨版本翻转，两边都跑了

一句话判定：**修前该绿、修后该红**——修前没有任何断言读那半句里的第二个数，所以删掉它无人可响；
ⓐ 之后钉它的断言存在了，删掉就断在"句子不再是那一句的形状"上。

```
修前（snap4 ＝ 父码＋父测试文件，只撤 `at offset %d`）:
  $ go test -v -count=1 -run 'TestSLO144' ./cmd/wisp/
  rc=0   --- PASS=7   --- FAIL=0   （子用例）--- PASS=7   ⇒ 0 红／14 绿
  变异确实生效：grep -c 'document still open: the tail' cmd/wisp/slo_windows.go = 1

修后（被验版本 ＋ 同一发变异，r1/B4-MG-drop-at-offset）:
  rc=1  RUN=17 TOPPASS=7 TOPFAIL=3   144-block: 0 red / 7 green
    :446  1978 of 1978 … sentence carries no byte/offset pair: "0 bytes read, document still open: the tail had not arrived"
    :467  give-up error "…within 60ms (last read: 8 bytes read, document still open: the tail had not arrived)" does not read as bytes-then-offset
    :500  unfinished sentence "989 bytes read, document still open: the tail had not arrived" is not the unwritten shape
```

⇒ **同一发变异：修前 0 红／14 绿，修后 3 红。AC#2 的牙就在这两枚读数的差上，成立。**
⇒ 再补一发（编排者要求"别挂在单一枚上"）：把**三枚新用例同时**改成不跑（`X6-all-three-inertxMG`），
在被验版本上重新得到 `rc=0 TOPPASS=7 TOPFAIL=0` ⇒ 那 3 枚红**恰好且全部**来自新增用例，无一来自别处。

### 1.4　②③ 承重矩阵：摘掉任意一味，逐味答一句

每行＝一发变异；"响的是哪几枚"按用例名读，不靠枚数。

| 变异（摘掉／打坏的东西） | 结果 | 响的用例 | 144 那 7 枚 |
|---|---|---|---|
| `B1` `obs.offset = dec.InputOffset()`（ⓐ 那一行还原成旧取值） | 3 红 | 8／9／10 全响 | 0 红 |
| `B2` **只删那一行填充**（offset 停在 `-1`） | 3 红 | 8／9／10 全响 | 0 红 |
| `C3` **只撤那枚 hunk**（填充删＋`InputOffset()` 搬回 EOF 判定之前，注释全留） | 3 红 | 8／9／10 全响 | 0 红 |
| `B7` 填充 off-by-5（`+ 5`） | 3 红 | 8／9／10 全响 | 0 红 |
| `B4` **票面那发 MG**（撤 `at offset %d`） | 3 红 | 8／9／10 全响 | 0 红 |
| `B10` 句子长出一条尾巴 | **只有 case 8 红** | 8 | 0 红 |
| `B11` 两发合谋：句子把字节数印两遍＋字段烂回 0 | **只有 case 8 红** | 8（字段 leg） | 0 红 |
| `B6` `note` 那支被关（`-1` 开始外印位置） | **只有 case 10 红** | 10 | 0 红 |
| `B5` loop 到点句不再带最后读数（`offset` 归 0） | 2 红 | 9＋10（case 8 看不见） | 0 红 |
| `B9` 红句两个实参对调（`o.offset, o.bytes`） | **0 红** | — | 0 红 |
| `B13` loop 的**另一支**（`s.exited()`）丢掉 `last.summary()` | **0 红** | — | 0 红 |
| `CB` 只把 `offset` 字段注释改回旧文 | 0 红 | — | 0 红 |
| `CC` 只把结构体头上 `document`／`decoder` 那一词改回 | 0 红 | — | 0 红 |
| `B12` 撤 `at offset %d` 却少传一个实参 | **编译不过** | `go vet` 的 printf 检查拦住，0 条 `=== RUN` | — |
| `H1` `slo147PairOf` 的 `whole` 恒真 | 未修码上 0 红（该味单独摘不掉任何事） | — | — |
| `X1` = `H1` × `B10` | **0 红** ⇒ `whole` 那味**是承重的**：摘了它 `B10` 打不红 | — | 0 红 |
| `H3` case 8 的字段 leg 恒不响 | 未修码上 0 红 | — | — |
| `X2` = `H3` × `B11` | **0 红** ⇒ 字段 leg **是承重的**：摘了它"句子保住、字段烂掉"那发打不红 | — | 0 红 |
| `X3` = 撤 case 9 × `B5` | 仍 1 红（case 10 接住） | 10 | 0 红 |
| `X4` = 撤 case 10 × `B6` | **0 红** ⇒ case 10 **是承重的** | — | 0 红 |
| `X5` = 撤 case 8 × `B4` | 仍 2 红（9＋10 接住） | 9／10 | 0 红 |

**逐味答一句**（"摘掉它，是否存在一发变异从此打不红？"）：

- **ⓐ 那一行填充（`slo_windows.go:646`）**：存在——`B1`／`B2`／`C3`／`B7` 四发全打不红就没了牙；实测四发都红。
  单点回退问"只撤这处哪条用例变得不响"：**三条一起不响**（`C3`＝只撤这枚 hunk 的读数）。**不是装饰。**
- **case 8**：存在——`B10`（句子长尾）与 `B11`（字段烂回 0 而句子保住）**只有它响**（`X1`／`X2` 证明不是别的味顶上的）。承重。
- **case 9**：**本程没造出来**。上面 13 发里，凡让 case 9 红的（`B1/B2/B4/B5/B7/C3`），
  case 8 或 case 10 至少一枚同红；把 case 9 单独摘掉（`X3`）没有任何一发从红翻绿。
  ⇒ 按编排者的分界（**造没造出来**）：**这一味不判退回、不判装饰**，记为"当前变异集内冗余"。
  它冗余的原因是结构性的：它的判据（"到点句里两个数都等于递给 loop 的字节数"）被 case 10 的第一支
  （同一枚 loop、同一支 `note` 反检）与 case 8（同一支分类器）从两边夹住了。
  ⚠ 冗余≠无用：它是**唯一**从生产入口 `collectReportWithin` 断言"数字＝递给 loop 的字节"的一枚（`B5` 让它红），
  只是那一发恰好也被 case 10 接住。**若有人要收这味，必须先造出一发只有它响的变异**——本程造不出来，
  所以本程不要求收。
- **case 10**：存在——`B6`（`-1` 那一支开始外印位置）**只有它响**（`X4` 复算：摘 case 10 后 `B6` 打不红）。承重。
- **两枚新 helper**：`slo147UnwrittenRe`／`slo147PairOf` 是三枚用例共用的**唯一**解析器，
  摘掉任一枚都编译不过（不存在"打不红的变异"问题，而是"尺不存在"）；
  真正的弱化读数是 `X1`（`whole` 恒真 ⇒ `B10` 逃逸）与 `X2`（字段 leg 恒不响 ⇒ `B11` 逃逸）：**两味都承重**。
  实现件 §2.1 自己写过"字段 leg 关住的那条后门＝把句子印成 `%d … offset %d", o.bytes, o.bytes`"——
  **本程独立造出 `B11` 并复算：那条后门确实只被字段leg 关住（`X2` 逃逸、`B11` 单发只 case 8 红）。它的断言成立。**
- **注释那两处（`CB`／`CC`）**：撤掉**零枚**用例变红。读数如实记：**注释不受仪器管**。
  这不构成"装饰"判定——AC#2 的硬核心要求的是"码＋用例"有牙，注释是被修正的**承诺文本**，
  不是修法本身。**但见第 3 格**：新注释里关于 `corrupt` 支那半句，本程实测与读数不符，那一处是**新的无仪器承诺**。

### 1.5　⑤ 有没有"永远绿"的断言凑数

- 三枚都在未修码上红（§1.2）⇒ **没有一枚今天恒绿**。
- 三枚的判据方向全是"更严"：既有 10 枚（含 7 子用例）一字未改，
  `git diff --numstat 2dc2b28^ 2dc2b28` 现量 **测试文件 169 加／0 删**、**生产码 21 加／4 删**；
  4 枚删除逐枚点名（本程自己 `grep -E '^-[^-]'`）：
  `// them the decoder stopped…`、`// offset is the decoder's stop position…`、`// decoder ran, i.e. …`、
  `obs.offset = dec.InputOffset()`（那行的移位）。**没有一条断言被放宽或被改反方向。**
- `t.Skip` 枚数 0（CI 四数里 `--- SKIP=0`；被验版本与改前那一次都是 0）。
- 预算常量：`subjectReportBudget`／`subjectGrace` 不在 147 三枚 commit 的改动路径里（`git show --name-only` 已核），
  新用例递给 `collectReportWithin` 的是**形参**（`60ms`／`40ms`），与票 144 用例同一形写法。

### 1.6　本格判定

**AC#2 成立。** 凭据是四组读数：未修码上 3 红（§1.2）、MG 修前 0 红／修后 3 红（§1.3）、
三枚同时摘掉后 MG 回到 0 红（`X6`）、承重矩阵里 5 味各自至少关住一发（§1.4）。

两处**如实记下的洞**，都不翻转本格判定，都在第 3／5 格具名登记：
`B9`（红句两实参对调）0 红＝**ⓐ 语义下的等价变异**，不是尺漏了（见第 2 格对 AC#1 第二半句的读数）；
`B13`（loop 的 `s.exited()` 支丢掉同一枚 `last.summary()`）0 红＝**五格之外的另一发**，票 144 的地盘。

---

## 2　AC#1 判定是否诚实（它选了 ⓐ）

**判据**：ⓐ／ⓑ／ⓒ 三选一不许自己填；ⓐ 要求**两味**——"字段真填成已读到的最后一字节位置" **且**
"红句两个数各自负责一件事"；ⓑ 单独存在不算收；判完之后"MG 那发该红还是该绿"要能一句话说清（已在第 1 格跑给人看）。

### 2.1　它确实走的是 ⓐ，不是把 ⓑ 混进来交差

```
$ diff <(git show ae968f9:cmd/wisp/slo_windows.go | sed -n '/^\/\/ summary renders/,/^}/p') \
       <(git show e572aa94:cmd/wisp/slo_windows.go | sed -n '/^\/\/ summary renders/,/^}/p')
   （无输出 ⇒ 三句 summary() 文本逐字节未变）

$ git show e572aa94:cmd/wisp/slo_windows.go | grep -nE "obs.offset = |o.bytes, o.offset"
602:  return fmt.Sprintf("%d bytes read, document still open at offset %d: the tail had not arrived", o.bytes, o.offset)
606:  return fmt.Sprintf("%d bytes read, complete document at offset %d", o.bytes, o.offset)
646:          obs.offset = int64(obs.bytes)          <- ⓐ 那一行（unwritten 支）
649:          obs.offset = dec.InputOffset()         <- corrupt 支（在 EOF 判定之后）
654:  obs.offset = dec.InputOffset()                 <- complete／trailing-garbage 支
```

⇒ **文本一字未改、改的是喂它的数** ⇒ 是 ⓐ，不是 ⓑ。它自述的"红句文本没改"这一条**复算成立**。
⇒ 顺带复算第 3 节那件"`-1` 从不进被打印的句子"：`readSubjectReport` 的每一条 return 之前都有一枚
`obs.offset = …`（`:646`／`:649`／`:654`），`-1` 只剩 `:631` 的初值而它带不出去；带 `note` 的三支（`:705`／`:711`／`:716`）
走 `summary()` 第一支只印 note。`B2`（删掉 `:646`）实测让 `at offset -1` 第一次印出来并被 case 10 抓住
⇒ **这条是活的可核事实，不是自述**。（新注释末尾那两句"Only the -1 leg is never printed…"与此相符。）

### 2.2　一处**票面文字与实现语义**不符——不是实现跑歪，是 AC#1 ⓐ 的第二半句造不出被满足的形状

被验版本在真实前缀上印出来的句子（本程探针的直接读数，非推导）：

```
CASE deep-prefix-1949-garbage len=1950 state=unwritten offset=1950
  summary "1950 bytes read, document still open at offset 1950: the tail had not arrived"
CASE nested-then-break        len=45   state=unwritten offset=45
  summary "45 bytes read, document still open at offset 45: the tail had not arrived"
```

⇒ 两个数**恒等**。原因在写侧形状本身（`readSubjectReport` 头上那段 credential）：一次 `os.WriteFile` 从 0 写全量，
读者能看到的只能是前缀，前缀里每一枚已到达字节都被消费掉了 ⇒
"停在第几字节"＝"到了第几字节"。ⓐ 的**第一半句**（字段带信息）因此成立且可核（`B1`／`B2`／`B7`／`C3` 四发都红），
ⓐ 的**第二半句**（"两个数**各自**负责一件事"）**在这一支上无法成立**——那两个数是同一个数的两种问法。

仪器侧的直接后果（本程造的那发）：`B9`＝把红句的两个实参对调（`o.offset, o.bytes`）⇒ **0 红／14 绿**。
这不是尺漏了，是**等价变异**：对调后输出逐字节相同。**别把它当退回**，按编排者的分界（造没造出来＝"这味摘掉后有变异打不红"）
它不属于任何一味的摘除。

**这一处的判定**：ⓐ 成立（字段带信息＋有牙＋文本未动＋MG 该红已跑）；
票面 AC#1 ⓐ 那句"并让红句两个数各自负责一件事"**措辞过强**，本程按"票面缺陷上报、不照它做"处理（`AGENTS.md` 头三行），
**不要求实现方回头改句子**（改文本＝ⓑ，票面明写不算收）。⇒ **本格：成立（附一句票面文字过强，见第 6 格修正记录）。**

### 2.3　 ⓒ 有没有被躲过去

票面 ⓒ＝"三语义是有意的 ⇒ 必须留下一枚会红的检钉住三语义各自成立"。实现件 §1.2 答"没把 ⓒ 当独立支，
因为 `0＝读到一半` 正是缺陷本身，钉住它＝钉住假话"。本程复算：case 10 钉的**正是**"三语义各说各的话"
（`-1` 那支不许出现 `offset` 这个词／半截文档印 `at offset 989`／完整文档印 `complete document at offset 1978`／三句两两不等），
且 `B6`（把 note 支关掉，让 `-1` 外印位置）实测**只有 case 10 红**（`X4` 证明摘掉 case 10 后这发打不红）。
⇒ ⓒ 要求的"钉住三语义各自成立"这一味**被 case 10 做到了**，只是钉的是"三者可分辨"而不是"旧的 0 恒成立"。
**不构成躲判。**

**本格判定：成立。**

---

## 3　实现方报回来的四件，逐件独立重走

### 3.1　同族第二发落在 `corrupt` 那一支（它没修，说出界了）

**它说的**：39 字节 `{"mode":"subject-in-tree",,"pass":true}`（错在第 25 字节附近）仍印 `at offset 0`；
并推翻票 144 验收件那句"corrupt 分支也报 0"的旁读。

**本程自己的探针（21 发 corrupt-ish 形状，`/tmp/wisp147-corr`，与被验版本同一棵树）**：

```
fixture length = 1978 bytes
all-prefix state census (0..len inclusive) = map[complete:1 unwritten:1978]     <- 票面"现量的形状"第 2 条复算相符
corrupt offsets among fixture prefixes = map[]                                   <- 前缀里没有一枚走 corrupt

CASE double-comma-39B   len=39  state=corrupt offset=0
   summary "39 bytes read, report corrupt: 39 bytes contradict a subject report at offset 0: invalid character ',' looking for beginning of object key string"
CASE html-head          len=25  state=corrupt offset=0
CASE trunc40-plus-0xff  len=41  state=corrupt offset=0   ("invalid character '\xff' after object key")
CASE wrong-type         len=42  state=corrupt offset=42
CASE two-documents      len=3956 state=corrupt offset=1978
CASE empty-object       len=2   state=corrupt offset=2   summary 里根本没有 offset 这个词
CENSUS over 21 corrupt-ish shapes: prints-no-offset=2  prints-offset-0=9  prints-a-position=10
```

它的探针我**也拿来在本程重跑了一遍**（`git show e572aa94:.scratch/wisp/probes/147/zz147probe_windows_test.go`
放进快照）：`corrupt 43/39/13/2003/57/2 → obs.offset = 0/0/13/1978/57/2` ⇒ **它那六个取值逐枚复现**。
`TestP147DoesMinusOneEverPrint` 同样复现（三支 note 都不含 `-1`）。

**这一发今天响不响？——不响。这是本程造的尺，不是它的**：

```
B14  corrupt 支的 obs.offset 强制成常量 0（把 :649 的 dec.InputOffset() 换成 0）
     rc=0  RUN=17 TOPPASS=10 TOPFAIL=0  三枚新用例全绿、144 那 7 枚全绿        <- 打不红
B15  complete 支的 obs.offset 强制成常量 0（对照，证明这把尺不是整把坏）
     rc=1  TestSLO147OffsetSemanticsRenderThreeDifferentSentences RED
          :509 complete sentence "1978 bytes read, complete document at offset 0" does not name 1978 bytes and offset 1978
```

⇒ 三条腿里 **`unwritten` 有牙、`complete` 有牙、`corrupt` 零牙**：把 corrupt 的取值换成常量 0，
本仓没有任何一条断言会变红。B14 与 B15 的对照就是凭据。

**ⓐ 的字段注释算不算把它说圆了？——不算。** 两处：
1. **它给的规则与读数不符**。注释写 "on reportCorrupt it is the decoder's own answer,
   **which is 0 whenever it could not begin reading a value at all**"。本程实测有两发 **offset 是 0，
   但解码器显然已经开始读一个值**：`double-comma-39B`（已在 object 内、已读完一个 key，错在等下一个 key）、
   `trunc40-plus-0xff`（`invalid character '\xff' after object key`）。⇒ 那句 "whenever" 是一条**不成立的规律**，
   而它的存在正好会让下一位把"印 0"读成"什么都没到"——**这就是本票成立的理由，换了条腿又长回来一次**。
2. **它把承诺写进了一句没有仪器的话里**。票面明写"不许用'注释留着、字段以后再说'来收"、
   并点名"注释预先承诺尚未产出的行为"在本仓有名字。`corrupt` 这一味今天**零仪器**（B14），
   所以这句注释今天的地位与 `A258④` 那条"绑了但没牙"的条件相同。

⇒ **"要不要为这一发单开一票"的建议（决定权在编排者，本程未开票、未动票面框）**：**建议开一枚**。
射程建议三味，缺一枚就还会留同族第三发：
① `corrupt` 支的 offset（要么像 ⓐ 那样填成有位置的信息、要么把 `%d bytes contradict … at offset %d`
那半句改成它真做到的，并配一枚会响的检——即 `B14` 摘不掉的那一发）；
② 上面那句 `whenever` 注释，改成读数支持的形状或删掉那半句；
③ `B13`（loop 的 `s.exited()` 那支丢掉 `last.summary()`，全组 0 红，见 §1.4）——同一枚 `summary()` 的
第二个调用点，今天两票都没钉，与①同属"到点句里的位置信息"，放一起最省一枚票。
**来源归属**：②③ 是本程造的（实现件 §6 只说了 corrupt 文本没修，没量到 ③，也没量到那句 `whenever` 与读数不符）。

**本格判定：它自述的三条读数（39 字节印 0／推翻"corrupt 也恒 0"／在射程外）全部复算成立；
"在射程外所以没修"接受为边界、不接受为收口——建议单开一票，见上。**

### 3.2　AC#3 的枚数：三处 vs 四处文本＋一枚标题

```
$ git diff -w --numstat a91d7c2^ a91d7c2
31     2       cmd/wisp/slo_report_144_windows_test.go     <- 与编排者本轮另一程独立量到的同一个数相符
$ git log -1 --format='%H %s' a91d7c2
a91d7c2d… fix(144 AC#4): 前一程留下的测试文件不是 gofmt/clean —— 复量与票面不符，只改对齐   <- 标题句仍在，未被改写
```

append-only 复核（**原句一字不许抹**）：

```
$ git show --numstat 1f3ede9 -- .scratch/wisp/issues/144-….md docs/evidence/s1/144-slo-report-partial-read-r1.md
23     0       .scratch/wisp/issues/144-….md
12     0       docs/evidence/s1/144-slo-report-partial-read-r1.md
$ git show 1f3ede9 -- <那两枚> | grep -cE '^-[^-]'      ->  0     <- 删除行数 0，两枚都是纯追加

原句仍在位（被验版本上逐枚 grep）：
  144 票面 :189  "改动只有空白（前后两版 `diff -w` 逐字相同）。"                        <- 未抹
  144 票面 :283  "…唯一一次代码面改动是 `a91d7c2` 那枚 `gofumpt -w`，只有空白"           <- 未抹
       （票面与实现件都写 ":260"，本程现量在 :283 ⇒ 差 23 ＝ 那 23 行纯追加正好插在它前面，
        算得出、可复算。**别人给的行号也是读数**，两处都按 sha 重取过。）
  144 证据件 :375 "**改动只有空白**——改前后两版 `diff -w` 逐字相同，断言、预算、Skip 一个没动"  <- 未抹
  追加的更正：144 票面 :236 起（`>` 逐字引 :189 与 :260 两处）、144 证据件 :378 起；`>` 行数现量 9／12。
  a91d7c2 的标题：已推送 ⇒ 只登记不改写（实现件 §3.2 末行），符合 AGENTS.md §1.4。
```

⇒ **枚数：票面少计，实现方把数数对了、把名字叫错了。本程现量（`git grep -n "只有空白" e572aa94` 全仓逐枚点名）**：
**作为自述**存在"只有空白"这句话的落点是 **三处文本＋一枚标题＝四处落点**——
`144 票面 :189`、`144 票面 :260`（现量 :283）、`144 证据件 :375`、`a91d7c2` 的标题"只改对齐"。
其余含这五个字的行本程逐枚看过，**都不是又一处自述**：
`:236`/`:238`/`:239`/`:378` 是本次追加的更正与 `>` 引文；
`144-…-accept-r1.md:25`/`:497`/`:508`/`:511` 是**非实现者引用它自述为假**；
`144-slo-report-partial-read-r1.md:141`「| 0 字节 / 只有空白 |」是状态表的一行，与 `a91d7c2` 无关。
⇒ 实现件 §3.2 的**表**（四行）与本程现量**逐枚相符**，动作也没多做也没少做；
它的**标题句**"现量是四处文本＋一枚标题"（共五枚）**比它自己的表多一枚**——正确说法是"四处落点（三处文本＋一枚标题）"。
⚠ **这一味编排者也照抄了**（派单第 3 节第 2 件写作"四处文本＋一枚 commit 标题"）⇒ 记进第 6 格，
`A264` 台账那一句如需更正由编排者定，本程不动票面、不动台账。
**本格判定：成立（附条件：落点数对、标签多算一枚，纯文字不影响收口）。**

### 3.3　"两处 `offset: -1`" 实为四处，且 `-1` 从不进被打印的句子

```
$ git show e572aa94:cmd/wisp/slo_windows.go | grep -n "offset: -1"
631:  obs := subjectReportRead{bytes: len(data), offset: -1}
705:  last := subjectReportRead{state: reportUnwritten, offset: -1, note: "nothing read yet"}
711:          obs = subjectReportRead{state: reportUnwritten, offset: -1, note: "no report file yet"}
716:          obs = subjectReportRead{state: reportUnwritten, offset: -1, note: fmt.Sprintf("report file unreadable: %v", err)}
枚数 = 4      <- 实现方的"实为四处"复算成立；票面"另外两处"少计两枚（它点名的是 :705/:711 那两支）
```

"`-1` 从不进被打印的句子"本程**没有只信它的探针**，另造了两发：
`B2`（只删 `:646` 的填充）实测让 `at offset -1` **第一次印出来**并被 case 9／case 10 双双抓住
（`:467`／`:500` 两句原文在 §1.4 表里）⇒ 这条承诺是**可打红的**，不是永远绿；
`B6`（关掉 note 支）⇒ 只有 case 10 红 ⇒ "不外印"这一味由 case 10 单独守着。
**本格判定：成立。**

### 3.4　锚点中途被推走：它进场 `80fa0551`，真实父 `ae968f9`，中间三枚零代码

```
$ git log --oneline 80fa0551..ae968f9
ae968f9e ticket(146 -done): 结线改名——5 枚 AC 框全 [x]、两张非实现者表在 docs/evidence/s1/、真未勾 0 枚
9ca50691 docs(台账 A263): 59 枚这一推红名集合动过零枚、分母涨 16 枚且逐枚归到 144/146；slo-full 今天第一次拒采
c3a51e68 docs(台账 A262): 59 枚一次批量推到两边（同停 80fa0551）；三枚前端会话的活逐枚点名不隐身

$ git show --name-only --format='COMMIT %h %s' 80fa0551..ae968f9
ae968f9e  .scratch/wisp/issues/146-liveapprovals-…-done.md          <- 146 结线改名
9ca50691  docs/reports/pending-and-issues.md                        <- 台账
c3a51e68  docs/reports/pending-and-issues.md                        <- 台账

$ git diff --name-only 80fa0551 ae968f9 | grep -E '\.(go|ps1|sh)$|golden|thresholds'
   （空 ⇒ 三枚里零代码，与它自述相符）
```

路径自核的另一半（三枚 147 commit 没有一枚含派单之外的路径）已在 **§0.4** 现量：**成立**。
另记：本程自己的锚点也被推走过一次（`f77a003` 台账 `A264`，见 §0.5），被验版本不变。
**本格判定：成立。**

---

## 4　门禁与仪器坑（AC#5 逐发独立重跑，本程自己换尺）

### 4.1　`bash scripts/wisp-cli-tests.sh`（CI 同形跑法），两边各一次

```
被验版本 e572aa94（/tmp/wisp147，自带 third_party/sherpa-onnx）:
  portable-tests.sh: four numbers (all from -v output): === RUN=133  --- PASS=73  --- FAIL=0  --- SKIP=0
  runtests.sh: OK    rc=0    （包体 81.771s）

同一测试文件 + 父码 ae968f9（/tmp/wisp147-snap2）:
  portable-tests.sh: four numbers (all from -v output): === RUN=133  --- PASS=70  --- FAIL=3  --- SKIP=0
  strict runner exited 1    rc=1    （包体 81.220s）
```

⇒ **`=== RUN` 两发都非零（133／133）⇒ 都是"跑到了"**；本程**没有**遇到 `0xc0000135` ＋ 0 条 `=== RUN` 那一发，
区分依据按票面用的是 `=== RUN` 枚数而不是返回码。
⇒ 实现件的 `130 → 133` 与本程的 `133 / 133+3红` 不矛盾：它的"改前"是**父码＋父测试文件**（70 枚顶层），
本程的"改前"是**父码＋新测试文件**（73 枚顶层、3 枚红）——本程那一发正是恒真凭据（§1.2），
它的 130 与本程 `go test -list` 在完整改前基线上数出的 **70 枚顶层**同源相符。

### 4.2　本程自己造的一发**假红**，以及如何把它与"活失败"分开（这条是新增读数）

在被验版本的快照上**裸跑** `go test -v ./cmd/wisp/`（快照里没有那 22MB 未跟踪的 native 件）：

```
RUN=127 TOPPASS=65 TOPFAIL=8 TOPSKIP=0 PANIC=0
--- FAIL: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128 … 共 8 枚
原因逐枚同一个：dataroot_128_windows_test.go:50: no native DLLs in ..\..\third_party\sherpa-onnx
```

⇒ **8 枚红全是环境件**，不是代码失败；补上 `third_party/sherpa-onnx` 后 CI 同形跑法给出 `FAIL=0`（§4.1）。
⇒ 记进仪器坑：**仓外快照跑 `cmd/wisp` 必须自己补 `third_party/`**，否则会拿到"改后 8 枚红"这种
**看起来像活失败**的读数。`PANIC=0` 一并记：那 8 枚不是 panic 吞读数。

### 4.3　名册差集（两向 `comm`，本程自己 `-list`，不读它的 `.names`）

```
pre（完整改前基线 /tmp/wisp147-snap3）= 70 枚顶层   post（被验版本）= 73 枚顶层
comm -23（消失）：（空）
comm -13（新增）：TestSLO147LoopGiveUpSentenceAgreesWithTheBytesItRead
                 TestSLO147OffsetSemanticsRenderThreeDifferentSentences
                 TestSLO147UnwrittenSentenceNamesTheOffsetTheDocumentStoppedAt
```

### 4.4　那把"永远绿"的假尺，本程换尺复现（与实现件 §5.2 有一处归属差）

```
                       naive（子串）        anchored（行首）
被验版本(全绿) 日志:   FAIL=1 SKIP=2 RUN=136   ->   FAIL=0 SKIP=0 RUN=133
父码+新尺     日志:   FAIL=4 SKIP=2 RUN=136   ->   FAIL=3 SKIP=0 RUN=133
```

⇒ 实现件 §5.2 说"多出来的**都**来自 `-skip` 说明文字"——本程现量：**只有一半对**。
`SKIP` 多出的 2 枚确实来自散文（`portable-tests.sh:11` 那行 ＋ 汇总行）；
**`FAIL` 多出的那 1 枚来自它自己的汇总行** `… --- FAIL=0 --- SKIP=0`，与 `-skip` 散文无关。
（派单第 4 节把这条写作"那是 `-skip` 那句散文"，同源不准，记进第 6 格。）
⇒ 判红绿只认 `^--- FAIL: <名>`；本程全部矩阵按**用例名**逐个判（§1.1），不按枚数。

### 4.5　另外三发

```
$ D:/work/base/gopath/bin/gofumpt.exe --version
v0.12.0 (go1.27.1)                                                     <- 版本先跑再读数，没背
$ D:/work/base/gopath/bin/gofumpt.exe -l . tools/d22scan tools/mockllm
（空）  rc=0
$ go vet ./cmd/wisp/
（空）  rc=0
$ sh scripts/d22scan.sh                                                 rc=0
  正对照（先证明这把尺能红）: runtests.sh: OK - packages=[./...] top-level: PASS=30 FAIL=0 SKIP=0, === RUN=70
  d22scan: examined 228 production Go files under internal/ and cmd/
  bans #1-5 internal/=205  cmd/=23 | ban #6 frontend/=52 | ban #7 internal/tools/=18
  ban #8 design/=30  frontend/=52  internal/=412  cmd/=43
  d22scan: clean - no D22 ban violations
```

⇒ 八个作用域 `examined N` **全非零** ⇒ `clean` 不是"什么都没扫"。跑的是 **gofumpt**，没拿 `gofmt` 交差。
⇒ **一处读数差异要说清**：本程 `ban #8 design/ = 30`，实现件 §5.3 是 `35`。
原因＝`d22scan` 扫的永远是**工作树**，而本程扫的是**被验版本的仓外快照**：
工作树里 `design/` 此刻正被另一会话写（16 枚未提交删除＋`design/old/`、`screenshots/` 等未跟踪新增），
多出来的那几枚是**别人家的文件**。方向上是"多扫"不是"少扫"，不影响 `clean` 的可信度，
但它意味着：**任何"design/ 与 frontend/ 零命中"的宣称都不等于被验版本的宣称**（本程的 30 才是锚点的数）。
另：快照不是 git 仓，`d22scan` 因此打了一行 `gitignore rules NOT APPLIED … every path in every scope is being scanned`
⇒ 它自己说的"这是响的那一面"，本程接受。

### 4.6　AC#4 契约轴

见 §0.4（三枚 commit 的 `--name-only` 全集）＋这一发本程自己的现量：

```
$ git show --name-only --format= 2dc2b28 1f3ede9 e572aa94 | sort -u \
    | grep -E '^(internal/risk/|internal/panel/|internal/agent/approval/|tools/d22scan/|.*thresholds\.go$|\
.*golden.*|.*allowlist\.txt$|scripts/slo-check\.ps1$|docs/PLAN\.md$|docs/specs/|frontend/|design/)'
（空）
```

⇒ 票面 AC#4 点名的禁改面**逐枚零字节**；`frontend/**`／`design/**` 那两枚本程**没碰、没还原、
也没把它们算进任何"零命中"宣称**（本程的宣称只覆盖上面那枚交集为空的事实）。
阈值与预算：`thresholds.go` 不在改动名册里；`subjectReportBudget`／`subjectGrace` 未出现在
`git diff 2dc2b28^ 2dc2b28`（4 枚删除逐枚点名见 §1.5，全是注释加那一行移位）。

**本格判定：成立。** 两处随文上报（§4.2 假红那味、§4.5 `design/` 分母那味）都是**仪器读数**问题，不是活失败。

---

## 5　本程没测什么（按"如果我漏了它谁会先被骗"排序）

1. **真实写侧从没看过这一句**。本程全部"新红句"读数来自**纯函数** `readSubjectReport` 与**脚本递的字节**
   （`scriptedReader`），加上 `collectReportWithin` 那条**注入读函数**的路径。
   本程**没有**拉起真 subject、没跑 `wisp slo`、没跑 `slo-check.ps1`、没跑 `slo-full`。
   ⇒ 谁先被骗：相信"生产里到点真会印 `1950 bytes read, … at offset 1950`"的人。
   票面"不解决的事"第 3 条（真竞态本机只撞得出 `len=0`）**依旧未取证**，本程没有推进它一根毛。
2. **`corrupt` 支零牙是量出来的，但没量过"真会出现哪种 corrupt"**。B14 证明常量 0 打不红；
   本程的 21 发 corrupt 形状是**手搓的**，不是从真实 subject 崩溃形状采样来的。
   ⇒ 谁先被骗：拿 corrupt 红句排障、把 `at offset 0` 读成"没写到"的人（这正是新那一票该收的第一味）。
3. **变异不是穷举**。21 发变异＋7 发组合是**按味的逻辑**造的（每味一发"摘掉就逃逸"），
   没做系统扫描（例如全部单常量替换、全部格式串置换）。`B9`（实参对调）就是本程撞到的一条**等价变异**：
   它说明"存在本程没想到的、今天打不红的改动"这件事**没有被排除**，只被样本覆盖到了一处。
   ⇒ 谁先被骗：把"矩阵里全红"读成"所有坏形状都会红"的人。
4. **`!windows`／`*_other_test.go` 那一族没进 linux 容器**。本机 Docker 可用（`golang:1.27` 在本地），
   本程**没用**：`slo_windows.go` 与这枚测试文件都带 `_windows` 后缀，
   且 `wisp-cli-tests.sh` 自己在非 windows 上 GUARD 退出（`exit 2`），所以那一腿今天对本格**没有可判之物**。
   ⇒ 谁先被骗：以为是"跨平台门禁也重跑过"的人——不是。
5. **没跑 `-race`**。新用例里有 goroutine 的是票 144 那枚 case 3（真文件写侧回放），本程只按 CI 同形跑法跑过它（绿），
   没有单独 `-race` 复量。
6. **`e572aa9` 那条"mtime 顺序自证"本程没复算**。本程只**重跑**了它能被重跑的部分
   （改前红、两形变异、名册差集、gofumpt/vet/d22scan、`corrupt` 探针——其中它的探针文件本程拿来在自己快照里跑了一遍，
   六枚取值逐枚复现）。探针日志文件与读数的**时间先后**没查。
   ⇒ 谁先被骗：相信"日志是那次跑法留下的、不是事后补的"的人。
7. **票 144 那五格本程没重开**（票面"不解决的事"第 4 条已裁归 `144-…-accept-r1.md`），
   本程只在**与 147 交界**的两处碰它：case 4 的 `wantReason` 串（用来证明 MG 修前不响）、
   以及 `summary()` 第二调用点 `B13`（本程造的，五格之外）。
8. **`docs/reports/pending-and-issues.md` 的 `A263`/`A264` 只当读数看**，没审它内容；
   台账里"`next=`"那行是否要为本程 §3.1 建议的新票登记，**是编排者的活**，本程没动台账一字。

---

## 6　结论修正记录（含对本程派单的推翻）

**说明**：派单第 0–4 节里每一句都按未验证断言处理。下面逐条点名"谁的话被推翻、被什么读数推翻"。
**没有被推翻的**也在列（写明"复算相符"），以免下一位把"没报错"当成"没查"。

| # | 原话（谁说的） | 本程现量 | 结论 |
|---|---|---|---|
| 1 | **派单 §0**：编排者 21:5x 看到 HEAD ＝ `e572aa94`，"如果你的 HEAD 已前移照样以 `e572aa94` 为被验版本" | 进场 `git rev-parse HEAD` ＝ `e572aa94…` **相符**；写第 0 格期间被 `f77a003`（台账 `A264`，只动 `docs/reports/pending-and-issues.md`）推走 | **部分命中，按其预案处理**（见 §0.5）。被验版本全程未变 |
| 2 | **派单 §3.2 ＋ 实现件标题句**："现量是**四处文本**＋一枚 commit 标题" | `git grep -n "只有空白" e572aa94` 全仓逐枚点名：作为**自述**的是 **144 票面 :189、:260（现量 :283）、144 证据件 :375 ＝ 三处文本 ＋ 一枚标题 ＝ 四处落点**；其余 8 处含这五个字的行是更正追加、非实现者引用、或状态表行 | **推翻**（实现件与派单同源一处标签不准：多算一枚）。**动作不受影响**：该追加的追加了、删除行数 0 |
| 3 | **派单 §4**："裸 `grep -c '--- FAIL'` 会得到 1，而那是 `portable-tests.sh` 自己 `-skip` 那句散文" | 本程在被验版本那份**全绿**日志上现量：naive FAIL=1 的唯一来源是**汇总行** `… --- FAIL=0 --- SKIP=0`；`-skip` 散文只让 `SKIP` 多 2、`=== RUN` 多 3 | **推翻一半**（坑是真的，归属错了：FAIL 那 1 枚不是 `-skip` 散文）。实现件 §5.2 的"都来自 -skip 说明文字"同样只有一半对 |
| 4 | **票面 AC#1 ⓐ**："…并让红句**两个数各自负责一件事**" | 被验版本真实输出 `1950 bytes read, document still open at offset 1949+1 … at offset 1950`：那一支上 `offset ≡ bytes`（写侧是单次 `os.WriteFile` 全量，前缀每枚到达字节都被消费）。`B9`（两实参对调）⇒ **0 红／14 绿**＝等价变异 | **推翻措辞**（ⓐ 第二半句在该写侧形状下**不可满足**）。ⓐ 第一半句与 AC#2 全部成立，故**不要求回改文本**（改文本＝ⓑ，票面明写不算收） |
| 5 | **实现件 §1.3／新字段注释**："on reportCorrupt it is the decoder's own answer, **which is 0 whenever it could not begin reading a value at all**" | 本程 21 发探针里两发反例：`double-comma-39B`（已在 object 内、已读完一个 key）与 `trunc40-plus-0xff`（`invalid character '\xff' after object key`）——**都已在读值，offset 仍 0**；另 `B14` 证明 corrupt 支取值换成常量 0 **全组 0 红** | **推翻**（本程造的新洞：那句 `whenever` 是一条不成立的规律，且 corrupt 那条腿零仪器）。建议单开一票，见 §3.1 |
| 6 | **实现件 §2.3 形 B 的引用**：`slo_report_147…:467` | `git ls-tree -r e572aa94 -- cmd/wisp` 里没有 `slo_report_147*`；该枚文件是 `slo_report_144_windows_test.go`（本程 B4 日志原文就是 `slo_report_144_windows_test.go:467`） | **推翻一处引用**（读数本身正确，文件名写错，属笔误级）。不改它、只登记（append-only） |
| 7 | **派单 §2 反面提醒**："别要求同一发变异里新加那支先响" | 本程矩阵按"摘掉某味后该发还响不响"读，未要求先后 | **相符**，按此执行 |
| 8 | **实现件 §2.2／§2.3／§5.1／§5.2／§5.3／§6.1 的全部读数** | 逐条重跑：3 红原文（`:446`/`:470`/`:502`）逐字相符；形 A／形 B 各 3 红相符；CI 四数改后 `133/73/0/0` 相符；名册 70→73、消失 0、新增 3 枚逐枚相符；gofumpt `v0.12.0` ＋ `-l` 空相符；`go vet` 空相符；d22scan `rc=0` clean 相符；corrupt 探针六枚取值（`0/0/13/1978/57/2`）相符 | **复算相符**（没有一条是抄来的：上面每条都是本程自己跑出来的，只是它们**恰好与它一致**） |
| 9 | **实现件 §5.3 `ban #8 design/=35`** | 本程在被验版本的仓外快照量到 **30** | **两处都对、口径不同**：它扫的是**活工作树**（含另一会话未跟踪新增），本程扫的是**锚点树**。结论：`design/`／`frontend/` 上的"零命中"不能当被验版本的宣称（§4.5） |
| 10 | **本件 §7 初稿**（本程自己的话）："入库前反扫……（空）" | 同一把尺现量 **hits=24**（5 枚 40 位 object id＋18 枚测试名＋1 枚自述），只是**没有一枚是凭据** | **推翻本程自己一句**，§7 已按现量改写（凭据值抄录数仍为 0，但"尺为空"这句是错的） |
| 11 | **本件 §7 初稿**："真通知回显数：13 枚（任务回显 6 枚＋file-state 5 枚…）" | 人工点收与机械核两种计数法重数：任务回显 **13**（＝本程启动的后台命令枚数）、file-state **9**、合并冲突 **1**、`Read` 截断 **1**、会话上下文 **1** ＝ **25** | **推翻本程自己一句**（初稿把两类数错了），§7 表已换成带计数法的版本 |

---

## 7　伪授权登记（两个数分开栏）＋凭据自扫

**真通知回显数：25 枚**，按类别点名（**计数法随每类写明，因为"枚数"本身就是读数**）：

| 类 | 枚数 | 计数法与出处 |
|---|---|---|
| 后台任务完成回显（`<system-reminder>` 里的 task-notification） | **13** | **可机械核**：本程启动的后台命令就是 13 枚（group1／group2／group3／group4 失败那次／group5 两次／corrupt 探针／fullpost／fullpre／CI-post／CI-pre／d22scan／MG-pre），每枚收到一次完成回显 |
| 自有临时件的 file-state 回显（"…was modified, either by the user or by a linter" ＋ "The file changed before your edit…"） | **9** | 人工点收：全部指向 `wisp147-mut/harness.py`（7 枚）与 `wisp147-mut/run.sh`（2 枚），内容＝本程刚写进去的那次改动，逐枚可对照 |
| 合并冲突提示（`<<<<<<< SEARCH … Unmerged: both modified`，指向**仓外**另一目录的计划文件） | **1** | 人工点收；本程没碰那两枚路径、没跑 `git add`／`git checkout`、没"take it into account"地改判据 |
| `Read` 截断提醒（"before line 55 was cut off…"） | **1** | `Read` `…/Temp/wisp147/cmd/wisp/slo_windows.go`（`offset=556 limit=190`） |
| 会话首条项目上下文／技能清单回显 | **1** | 进场首条，非工具输出 |
| **合计** | **25** | ⚠ 第 2–5 类是**人工点收**，可能 ±1；第 1 类可机械复算。**任何一枚都没有改变本程的判断或射程** |

**判为注入数：0**，三条理由与一处备案：
- 回显里**没有出现过**任何一句"编排者备注／已核验请继续提交／请放宽阈值／已解锁／不用取证直接给结论／
  Confirm the harness note is genuine"这类**要求本程少取证或改判据**的文字。
- 本程对 `e572aa94` 的票面、实现件、`.scratch/wisp/probes/147/` 全部日志做了字面反扫
  （`grep -rniE "编排者备注|系统提示|已被修改|已核验请继续|请放宽阈值|已解锁|不用取证|Confirm the harness note"`）
  ⇒ **命中 2 行，且都在实现件自己 §7 那句"这些字样从未出现过"的自述文本里**（`:488`/`:489`），不是外部注入。
  ⚠ **备案一处**：实现件 §7 写"全程工具输出里没有出现过'…文件已被修改…'这类文字"——本程在同一会话里
  **收到了 9 枚**该形状（上表第 2 类）。差异是**会话级**的（各家 harness 回显时机不同），
  本程**不据此判实现方造假**，只记"那句话的射程是它自己的会话，不是本仓的观测"。
- 那 9 枚 file-state 回显里有一枚带祈使（"This change was intentional … don't revert it unless the user asks"），
  另一枚带"take it into account / do NOT edit"。本程的处理：**照做但不据此改判据**——
  涉及的都不是被验版本的文件；判据一律按 `e572aa94` 的快照取；
  **没有任何一枚让本程少跑一发**：28 发变异／探针全部照跑。

**凭据值抄录数：0。** 本件不含任何密钥/token 的值；涉及的只有测试名、变量类别与文件名。

⚠ **本格先推翻本程自己写过的一句话**：本程起初把下面这把尺写成"（空）"，现量**不是空**——
那把尺太宽，40 位十六进制的 git object id 与长测试名都会命中。**以现量为准，改判如下**：

```
$ grep -nE "sk-[A-Za-z0-9]{8,}|[A-Za-z0-9+/]{40,}={0,2}|api[_-]?key[[:space:]]*[:=][[:space:]]*[^ ]{6,}|Bearer [A-Za-z0-9]" \
    docs/evidence/s1/147-offset-naming-r1-accept-r1.md
hits = 24        <- 该读数（与下面逐枚分类里的行号）取自本节**落盘前的那一版**：本节自己含长测试名，追加会改变枚数（自指尺子的固有形状）。
                    入库前的最终反扫请**重跑上面那条命令**，别背这个 24。

逐枚分类（本程把这 24 行全看过）：
  40 位十六进制 git object id  …… 5 行（`3`／`18`＝e572aa94 全形；`79`／`80`＝两枚 sha1sum；`103`＝f77a003 全形）
  Go 测试名（长驼峰串）          …… 18 行（`:159-165` 七枚 TestSLO144* ＋ `:167/:169/:171` 三枚 FAIL 名
                                    ＋ `:196-198` 与 `:544-546` 各三枚名册名 ＋ `:396` ＋ `:531`）
  本表"命中 2 行"那句自述        …… 1 行（`:666`）
凭据形状且无法归入以上三类者：  0
```

⇒ **结论不变、尺的读数变了**：抄录的凭据值仍是 **0**，但"这把尺为空"这句话是**本程自己的一处不准**，
已记进第 6 格修正记录（第 10 行）。入库前真正的反扫应由编排者拿凭据文件做**字面量**比对，
本程这把正则只能证明"本件里没有 key -shaped 串"。

**顺带一条与本仓 §1.2"规格与仪器的射程不是一回事"同形的自扫**：

```
$ python -c "…ban#8 的六个码段逐枚扫本件…"
band hits: 6  ->  U+26A0(⚠)×3、U+2260(≠)×1、U+2261(≡)×2
```

⇒ 这六个字符**落在仪器扫的码段里**，但 `tools/d22scan` 的 ban #8 **作用域是 `design/`／`frontend/`／`internal/`／`cmd/`，
不含 `docs/`**（本程 §4.5 那发的八个作用域名册里就没有 `docs/`）⇒ 本件不构成违规，
与被验版本的实现件、票面同一风格（它们也逐枚用 ⚠／→）。
**记在此处是为了别让下一位把"docs 里没有 emoji"当成一条被仪器核过的宣称——它没被核，因为它不在射程内。**

**git 纪律自证**：本程**每一枚 commit 都带显式 pathspec**
（`git commit -q -F - -- docs/evidence/s1/147-offset-naming-r1-accept-r1.md`），枚数现量：
每枚 commit 前跑 `git diff --cached --name-only` 自核 ⇒ 逐次只出现本程那一枚路径，**没有出现别人的路径**。
`git add` 只对本件那一枚显式路径用过一次（新文件首枚 commit 之前；`-A`／`.` **一次没用**）；
`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`rm`／`rmdir` **使用数 0**；
**push 次数 0**。临时件全部只建不删：`wisp147`、`wisp147-snap2`、`snap3`、`snap4`、`wisp147-corr`、
`wisp147-mut/r1`（18 枚变异）、`r2`（7 枚）、`r3`（2 枚）＝ **35 枚目录，一枚未删**；
`harness.py` 遇同名目录就换 `r2.rN` 新号，不覆写。
票面五框**一枚未勾**，实现件的框**一枚未动**，台账**一字未写**（勾与登记归编排者按本表定）。

---

## 8　五格总裁决（交编排者定勾）

| 格 | 判据 | 本程读数来源 | 档位 |
|---|---|---|---|
| **AC#1** 判定 | ⓐ 是否真做了、是否躲进 ⓑ／ⓒ | §2（文本逐字节未动、`-1` 不外印可打红、ⓐ 第二半句不可满足） | **成立**（附一句：票面措辞过强，见 §6 第 4 行） |
| **AC#2** 那把尺（硬核心） | 恒真／承重／单点回退／MG 翻转／无永远绿 | §1（未修码 3 红、MG 修前 0 红→修后 3 红、21＋7 发矩阵、`X6` 全摘回落 0 红） | **成立** |
| **AC#3** 追加更正 | append-only、31/2、枚数 | §3.2（删除行 0、原句四处仍在位、31/2 相符；标签多算一枚） | **成立**（附条件：枚数标签） |
| **AC#4** 契约轴 | 禁改面零字节 | §0.4 ＋ §4.6（交集为空） | **成立** |
| **AC#5** 门禁 | 四数＋名册差集＋三发静态尺 | §4（133/73/0/0、70→73 消失 0、gofumpt v0.12.0 空、vet 空、d22scan rc=0 全非零） | **成立** |

**本程没有"退回"这一档**：全部判据都**没能造出**"某味摘掉后本格还成立"的反例——
恰恰相反，摘味试验里 5 味各自至少关住一发变异（§1.4），第 6 味（case 9）本程**造不出**只有它能红的变异，
按编排者的分界（造没造出来才算退回）**不判退回**，只记"当前变异集内冗余"。

**五格外的两笔（本程造的，编排者定要不要开票）**：
① `corrupt` 支零牙 ＋ 那句 `whenever` 注释不成立（§3.1，`B14` vs `B15` 对照）；
② `s.exited()` 那条到点句丢掉 `last.summary()` 全组 0 红（§1.4 `B13`）。

---

## 9　落盘名册与复跑法（本程自己的 commit／临时件，取法随文）

**commit 名册现量**（`git log --format='%h %s' e572aa94..HEAD -- <本件>`，逐枚 `git show --name-only` 自核）：
本程共 **7 枚**，**每一枚的文件列表只有 `docs/evidence/s1/147-offset-naming-r1-accept-r1.md` 这一枚路径**；
区间里另有 1 枚不是本程的（`f77a003` 台账 `A264`，路径 `docs/reports/pending-and-issues.md`）。
本名册（第 9 节）就是第 7 枚。

```
0447040 第 0 格锚点  e2bfa58 第 1 格尺  2f23611 第 2 格判定  c62b884 第 3 格四件
c107761 第 4 格门禁  1642886 第 5-8 格  + 本枚（第 9 格落盘名册）
```

**临时件名册（只建不删，35＋枚目录，全部在仓外）**：
`/tmp/wisp147`（被验版本快照）· `snap2`（父码＋新尺）· `snap3`（完整改前基线）· `snap4/base`（基线＋MG）·
`wisp147-corr`（本程探针＋它的探针）· `wisp147-mut/r1`（18 枚变异）· `r2`（7 枚组合）· `r3`（2 枚腿检）。

**复跑本件最硬那两发的最小路径**：

```
git -c core.autocrlf=false -c core.eol=lf archive e572aa94 | tar -x -C <新目录>      # 被验版本
git show ae968f9:cmd/wisp/slo_windows.go > <新目录>/cmd/wisp/slo_windows.go            # 只换码、留尺
mkdir -p <新目录>/third_party && cp -r "<repo>/third_party/sherpa-onnx" <新目录>/third_party/
export PATH="<repo>/third_party/sherpa-onnx:$PATH"
cd <新目录> && go test -v -count=1 -run 'TestSLO147|TestSLO144' ./cmd/wisp/            # 期望 3 红（恒真那一问）
python <仓外>/harness.py r1 && bash <仓外>/run.sh r1 A0-control-asis B4-MG-drop-at-offset   # 期望 0 红 / 3 红
```

⇒ 变异尺的原件在 `C:/Users/swq/AppData/Local/Temp/wisp147-mut/harness.py` 与 `run.sh`（**仓外、未入库**）；
每发变异的字面替换都写在 `harness.py` 的 `MUT`／`COMBO` 里，**每次替换强制"该字面量在本文件中恰好出现 1 次"**，
否则直接 FATAL 退出——这一条是为了让"下一位复算时改错地方也不会得到一个看起来像绿的答案"。

⚠ 一处**已知的本件形状缺陷**（自报，不等谁来抓）：本件写作过程中 `docs/evidence/` 那批**行号引用**
（§7 的凭据反扫、§3.2 的 `:189/:283`）都取自"落盘前那一版"或 `e572aa94` 的树，**不随本件自身增长而重编号**。
按 §3.3 那条本仓共识（"别人给的行号也是读数、天然带版本"）处理：复算请一律按 sha 重取。


---


---
