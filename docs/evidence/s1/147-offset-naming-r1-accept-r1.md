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


---
