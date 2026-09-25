# 149 实现件 r1 — `corrupt` 那条腿的取值、那句注释的反例、loop 第二个调用点

- 票面：`.scratch/wisp/issues/149-the-corrupt-leg-of-the-offset-field-has-zero-teeth-the-new-comment-claims-a-sufficient-condition-with-counterexamples-and-the-exited-branch-drops-the-summary.md`
- 来路：票 147 对抗验收件 `docs/evidence/s1/147-offset-naming-r1-accept-r1.md` §3.1 那三发（编号 `B14`／`whenever` 注释／`B13`）
- 本程角色：**实现程**。票面框一枚不勾（`Acceptance` 明写"由编排者按非实现者验收表定"）
- 地界（票面 §地界）：`cmd/wisp/slo_windows.go` 的 `summary()`／`readSubjectReport` 的 corrupt 支／`collectReportWithin` 的 `exited()` 支
  ＋ `cmd/wisp/slo_report_144_windows_test.go` ＋ 本票证据件＋`.scratch/wisp/probes/149/`

---

## 0　锚点与工作树（进场第一读）

**判据**：HEAD 必须与派单给的锚点相符；脏工作树里别人的东西一律不碰、不还原、不计进任何"零命中"宣称。

```
$ git rev-parse HEAD
64858d6838ced46fbe7bcce38dc4f7bb163d2f9c

$ git branch --show-current
dev

$ git log --oneline -6
64858d6 docs(台账 A265): 票 147 结线五格全成立；验收程造出四发，其中一发我拒绝开格、写成"不许开格"的限定语
de1ee7b ticket(147 结线＋144/147 改名＋149 新立): 五框按非实现者表翻勾；同族三处零牙出口立成 blocked-on-147-改名的下一票
e9412c9 evidence(147 对抗验收 r1 第 9 格): 落盘名册七枚、仓外临时件 35 枚只建不删、复跑最小路径与变异尺字面量自锁
1642886 evidence(147 对抗验收 r1 第 5-8 格): 本程没测什么八条、结论修正记录十一条（含推翻本程自己三句）、伪授权两数分栏、五格总裁决
c107761 evidence(147 对抗验收 r1 第 4 格): 门禁五发独立重跑——CI 同形两边四数、名册 comm、假红 8 枚的来源、naive/anchored 归属差、d22scan 分母随树
c62b884 evidence(147 对抗验收 r1 第 3 格): 四件报回逐件重走——corrupt 支零牙(B14)/complete 支有牙(B15)、append-only 与 31/2 相符、offset:-1 实为四处、锚点三枚零代码

$ git status --porcelain
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
 M docs/reports/frontend-session-log-zcode.md
 M frontend/scripts/gen-tokens.mjs
 M frontend/src/App.tsx
 M frontend/src/components/nav-rail.tsx
 M frontend/src/components/panel-skeleton.tsx
 M frontend/src/main.tsx
 M frontend/src/styles/theme.css
?? .zcodeignore
?? design/doubao/01-ball-states.jpg
?? design/doubao/demo/lib/
?? design/doubao/demo/rb-files.js
?? design/doubao/demo/rb-plugins.js
?? design/doubao/demo/rb-review.js
?? design/doubao/demo/rb-terminal.js
?? design/doubao/demo/rightbar.js
?? design/doubao/demo/screens/home.js
?? design/doubao/demo/screenshots/
?? design/doubao/demo/sidebar.js
?? design/old/
```

**读数**：
- 锚点相符（`64858d6`＝派单说的 `64858d6`），分支 `dev`。
- 工作树脏在**别人家**：`design/**`（owner 自己挪动的 16 枚未提交删除＋`design/doubao/**`、`design/old/` 等未跟踪件）、`frontend/**`、
  `docs/reports/frontend-session-log-zcode.md`、`.zcodeignore`。本程**一枚未碰、未 add、未还原**。
- 进场时 `git status` 里**没有** `internal/agent/**`（派单说 138 那枚程在跑那棵）；本程全程不碰 `internal/**`。

**取法**：所有"改前"的读数一律按 sha 取，不把脏工作树当被验那版——

```
$ mkdir -p /tmp/wisp149/snap-pre && git -c core.autocrlf=false -c core.eol=lf archive 64858d6 \
    | tar -x -C /tmp/wisp149/snap-pre
$ mkdir -p /tmp/wisp149/snap-pre/third_party \
    && cp -r "<repo>/third_party/sherpa-onnx" /tmp/wisp149/snap-pre/third_party/
```

⇒ `third_party/` 是**未跟踪的原生 dll**（`git archive` 不带它，见派单 §5 那条坑）。本程**先补 dll 再报红**，
没有拿"快照里 8 枚红"当活失败——那 8 枚的形状（`no native DLLs in ..\..\third_party\sherpa-onnx`）在下面的 §6.4 有现量。

**放水两问自答**：本格不动断言、不动 helper（无断言可动）。
**本格判定**：锚点相符、脏区点名、取法可复算。成立。

---

## 1　AC#1 — 三发各自"今天响不响"，现量

**判据**（票面 AC#1）：不许沿用验收程的 `B13`/`B14`/`B15` 编号当凭据，必须自己现量；
全为"今天不响"本票才成立；只要有任何一发今天已经响，那一发归回票 147 结线。

**跑法**（两枚命令，可复算；变异一律落在仓外快照、每枚一份 overlay＋一份原始日志，只建不删）：

```
$ python C:/Users/swq/AppData/Local/Temp/wisp149/harness.py <快照树> <日志目录> <变异名...>
    # 尺 = go test -count=1 -overlay <overlay.json> -v -run 'TestSLO144|TestSLO147|TestSLO149' ./cmd/wisp/
    # 每发变异强制"该字面量在目标文件里恰好命中 1 次"，否则 FATAL 退出
```

"今天"＝**HEAD `64858d6` 的码＋HEAD `64858d6` 的测试文件**（`/tmp/wisp149/snap-pre`，17 枚顶层 SLO 用例），
里面**没有**本票的任何东西。三发的字面量替换原文与逐发红句见 `.scratch/wisp/probes/149/harness.py` 与同名日志。

```
pre-asis                           rc=0 RUN=17 topPASS=10 topFAIL=0 144red=0 147red=0
pre-B14-corrupt-fill-to-zero       rc=0 RUN=17 topPASS=10 topFAIL=0 144red=0 147red=0     <- 零牙（复算成立）
pre-B15-complete-fill-to-zero      rc=1 RUN=17 topPASS=9  topFAIL=1 144red=0 147red=1
      RED: TestSLO147OffsetSemanticsRenderThreeDifferentSentences
pre-B13-exited-drops-summary       rc=0 RUN=17 topPASS=10 topFAIL=0 144red=0 147red=0     <- 零牙（复算成立）
pre-B5-timeout-drops-summary       rc=1 RUN=17 topPASS=7  topFAIL=3 144red=2 147red=1     <- 对照组：同一枚 summary 的另一支有牙
      RED: TestSLO144LoopGiveUpSentencesOnRealFiles, TestSLO144UnfinishedReportKeepsPollingThenNamesBudgetAndBytes,
           TestSLO147LoopGiveUpSentenceAgreesWithTheBytesItRead
```

**逐发红句原文**（有牙的那两发，`--- FAIL` 之外的正文）：

```
pre-B15:  slo_report_144_windows_test.go:509: complete sentence "1978 bytes read, complete document at offset 0"
                  does not name 1978 bytes and offset 1978
pre-B5:   slo_report_144_windows_test.go:232: unfinished give-up "wisp slo: subject 4243 never wrote a complete report within 50ms"
                  does not name "44 bytes"          （同一发还有 :232 的 "tail had not arrived" 与 :327 的 "8 bytes"）
```

⇒ **第一发（`corrupt` 支换成常量 0）今天不响**，`pre-asis` 与它的四数逐枚相同（17／10／0），
说明那一发**连一条断言都没碰到**；**第三发（`exited()` 支丢掉 `last.summary()`）今天不响**，同理。
⇒ 两发各有**同族对照**响着（`complete` 支、`timeout` 支），所以"不响"是**这两条出口没尺**，不是这把尺整把坏。
⇒ `=== RUN=17` 非零 ⇒ 三发都**跑到了**（区分"没跑到"的依据是 `=== RUN` 枚数，不是返回码）。

**第二发（注释那句 "0 whenever it could not begin reading a value at all"）现量**——探针 `zz149probe_windows_test.go`
（同一棵树，`go test -v -run TestP149CorruptShapesOffsetCensus`，全文日志 `probes/149/probe-corrupt-census.log`）：

```
P149 FIXTURE len=1343
P149 double-comma-39B       len=39    z(offset)=0   errkind=SyntaxError         eo=27
P149 html-head              len=43    z(offset)=0   errkind=SyntaxError         eo=1
P149 trunc40-plus-0xff      len=41    z(offset)=0   errkind=SyntaxError         eo=41
P149 deep-at-500-at         len=1343  z(offset)=0   errkind=SyntaxError         eo=501
P149 deep-late-2-0xff       len=1343  z(offset)=0   errkind=SyntaxError         eo=1342
P149 bad-escape-in-string   len=25    z(offset)=0   errkind=SyntaxError         eo=11
P149 colon-then-comma       len=39    z(offset)=0   errkind=SyntaxError         eo=26
P149 brace-first-then-good  len=1344  z(offset)=0   errkind=SyntaxError         eo=1
P149 wrong-type             len=13    z(offset)=13  errkind=UnmarshalTypeError  eo=12
P149 array-into-struct      len=7     z(offset)=7   errkind=UnmarshalTypeError  eo=1
P149 number-into-struct     len=3     z(offset)=3   errkind=UnmarshalTypeError  eo=3
P149 string-doc             len=15    z(offset)=15  errkind=UnmarshalTypeError  eo=15
P149 late-type-break        len=51    z(offset)=51  errkind=UnmarshalTypeError  eo=50
P149 empty-object           len=2     z(offset)=2   errkind=none                eo=-1   （句中无 offset 这个词）
P149 no-state-report        len=57    z(offset)=57  errkind=none                eo=-1   （同上）
P149 null-doc               len=4     z(offset)=4   errkind=none                eo=-1   （同上）
P149 trailing-garbage       len=1368  z(offset)=1343 errkind=none               eo=-1
P149 second-document        len=2686  z(offset)=1343 errkind=none               eo=-1
P149 CENSUS corrupt-ish shapes=20 prints-no-offset-word=3 prints-offset-0=8 prints-a-position=7
P149 ALL-PREFIX state census (0..len=1343) = map[complete:1 unwritten:1343]
```

⇒ **反例成立且不止两发**：20 发里 18 发判 corrupt，其中 **8 印 `at offset 0`**，而这 8 发里
`deep-at-500-at`（前 500 字节是好头，坏在 index 500 的 `@`）与 `trunc40-plus-0xff`（前 40 字节是好头）
**明显已经在读值**——票面点名的两发（39 字节双逗号、`'\xff' after object key`）复算相符，本程另量到 6 发同类。
⇒ 另一侧的关键读数：**解码器自己的错误一直带着一枚位置**（`eo`，20 发里最小 1，没有一枚是 0），
`dec.InputOffset()`（`z`）在 corrupt 支上给出的 0 **不是"读不下去"的记号，而是它缓冲的产物**。
⇒ `ALL-PREFIX state census` 那行也复算了票面第 2 条的隐含前提：**真实前缀里没有一枚走 corrupt**（1343 发全 unwritten＋1 发 complete）。

**放水两问自答**：
① 断言方向动没动——**没动**，本格零断言、零生产码改动，只是量。
② helper 是不是原有的那枚——本格用的 `scriptedReader`／`slo144Report`／`slo147PairOf` 都是**票 144/147 原有的**；
新增的两枚（`slo149ContradictionOf`、`slo149Poke`）在 §2 里只被新用例用，且新用例**不改任何既有用例**。

**AC#1 判定**：**三发今天全不响 ⇒ 本票成立，射程不缩**（票面 AC#1 的"归回 147"那一支没有触发）。

---

## 2　AC#2 — 给 `corrupt` 那条腿造一枚会响的检

**判据**：这枚检必须**只有"腿上的取值被写死／换成与实读不符"才能触发**；承重按本仓操作定义答一句。

### 2.1　修法（ⓐ 的形状，与本族前两票同形：文本不动，改喂给它的数）

`readSubjectReport` 的 corrupt 支今天填 `dec.InputOffset()`，而 §1 量到它在 8 发上印 0。
本票把它换成**解码器自己那枚错误带的位置**：新增 `contradictionOffset(err, inputOffset)`
（`*json.SyntaxError.Offset` → `*json.UnmarshalTypeError.Offset` → 都没有时保留 `inputOffset`），
corrupt 支改为 `obs.offset = contradictionOffset(err, dec.InputOffset())`。
**句子一字未动**（`%d bytes contradict a subject report at offset %d: %w`），改的是喂进去的数。

新增用例 case 11／case 12（同一枚测试文件、原有 fixture、原有 helper）：
case 11 的 8 发每发的期望位置是**本文件自己声明的构造**（坏字节在 index K ⇒ 期望 K+1），
**不是从解码器读回来的数**——这是它和 §1.4 那把"永远绿"的假尺的区别。

### 2.2　恒真那一问：**在未修码上响**（这是凭据，不是"改后全绿"）

```
$ # 树 = HEAD 64858d6 的 slo_windows.go ＋ 本票的测试文件（只跑 SLO 组）
$ cd /tmp/wisp149/newtest-oldcode && go test -v -count=1 -run 'TestSLO144|TestSLO147|TestSLO149' ./cmd/wisp/
--- FAIL: TestSLO149CorruptSentenceNamesThePositionTheDecoderObjectedAt (0.00s)
--- PASS: TestSLO149CorruptLegsWithoutADecoderErrorKeepTheirOwnEnd (0.00s)
--- PASS: TestSLO149ExitedGiveUpSentenceCarriesTheLastReading (0.06s)
    slo_report_144_windows_test.go:655: 8 of 8 corrupt shapes do not name the byte position the decoder objected at;
                                        first: second-comma-at-26: subjectReportRead.offset is 0, want 27 (the byte the decoder objects at)
（144 那 7 枚＋147 那 3 枚在未修码上全绿 ⇒ 响的不是"把旧断言改严"换来的）
```

⇒ **改前红句原文**＝上面 `:655` 那一行；**改后同一发转绿**＝§2.3 的 `post-asis`（`RUN=20 topFAIL=0`）。
⇒ case 12／case 13 在未修码上是**绿**的，它们钉的是另外两味（§2.4 的 `p10`/`p11`/`p7`/`p8` 让它们红），
不是这一味——本程**不把它们算作 AC#2 那枚牙的凭据**，只算同族的第二、第三处出口。

### 2.3　改后（本票交付的树）

```
$ python .../harness.py C:/Users/swq/AppData/Local/Temp/wisp149/post .../logs/post-final post-asis
post-asis  rc=0 RUN=20 topPASS=13 topFAIL=0 144red=0 147red=0
           new=CorruptSentenceN:PASS CorruptLegsWitho:PASS ExitedGiveUpSent:PASS
```

⇒ **新路径真的执行了**：`=== RUN` 从 17 涨到 20、三枚新用例逐枚 PASS（不是"包级绿"）；
名册差集见 §6.2（消失 0／新增 3，逐枚 `TestSLO149*`）。

### 2.4　承重矩阵：摘掉任意一味，是否存在一发变异从此打不红？

每行一发变异（overlay，字面量恰好命中 1 次强制），"响的是哪几枚"按**用例名**读。全部日志在 `probes/149/post/`、`probes/149/combos2/`。

| 变异 | 结果 | 响的用例 |
|---|---|---|
| `p1` corrupt 取值换成常量 0（＝票面 `B14`，本票的靶） | **1 红** | case 11 |
| `p2` corrupt 取值换成 `int64(obs.bytes)` | **1 红（7 of 8 发）** | case 11 |
| `p3` corrupt 取值退回 `dec.InputOffset()`（＝撤掉本票修法） | **1 红** | case 11 |
| `p4` 摘掉 `SyntaxError` 那一支 | **1 红（5 of 8）** | case 11 |
| `p5` 摘掉 `UnmarshalTypeError` 那一支 | **1 红（3 of 8）** | case 11 |
| `p6` 句子里撤掉 `at offset %d` | 2 红 | case 11 ＋ 144 的 case 2 |
| `p12` 句子里第二个计数换成 offset（"bytes contradict" 与 "bytes read" 打架） | 3 红 | case 11 的 `named != read` leg ＋ 144 两枚 |
| `p7` loop 的 `exited()` 支丢掉 `last.summary()`（＝票面 `B13`） | **1 红** | case 13（13a/13b/13c 三句全响） |
| `p8` `exited()` 支把 summary 换成常量 | **1 红** | case 13 |
| `p9` `timeout` 支丢掉 summary（对照，147 已有的牙） | 3 红 | 144 两枚＋147 一枚；**case 13 不响**（它不钉那一支） |
| `p10` 把修法 smuggle 到"第二个文档"腿 | **1 红** | case 12（offset 1981 want 1978 ＋ 句子那句） |
| `p11` 让"干净但无状态报告"那条腿开始印位置 | **1 红** | case 12（"prints a byte position it does not have"） |
| `d11a`/`d11b`/`d11c`/`d11d` 单独摘掉 case 11 的四条 leg 之一 | 0 红 | legs 之间**重叠**（见下） |
| `d12-offset-assert-off`／`d12-note-leg-off` 单独摘掉 case 12 的一条 leg | 0 红 | 同上 |
| `d11e` 把 case 11 整枚摘掉（改名不成用例） | 0 红、`NORUN` | — |
| `d13b` 把 case 13 整枚摘掉 | 0 红、`NORUN` | — |
| **`x-p1-d11a`**（常量 0 ＋ 摘字段 leg） | 1 红 | case 11 的句子 leg **单独接住** |
| **`x-p1-d11b`**（常量 0 ＋ 摘句子 leg） | 1 红 | case 11 的字段 leg **单独接住** |
| **`x-p1-d11a-d11b`**（常量 0 ＋ 两条 leg 全摘） | **0 红 ⇒ 逃逸** | — |
| `x-p2-d11a-d11b`／`x-p3-d11a-d11b`（同样两条 leg 全摘 × 另两发取值变异） | **0 红 ⇒ 逃逸** | — |
| `x-p1-d11e`（常量 0 × case 11 整枚摘掉） | **0 红 ⇒ 逃逸** | — |
| `x-p7-d13a`（丢 summary × case 13 的字面量弱化） | 1 红 | 13c 那句仍接住 |
| `x-p7-d13b`（丢 summary × case 13 整枚摘掉） | **0 红 ⇒ 逃逸** | — |
| `x-p10-d12`（smuggle × 摘 case 12 的 offset 断言） | 1 红 | case 12 的句子 leg 接住 |
| `x-p12-d11c`／`x-p12-d11d` | 3 红／2 红 | `p12` 由 144 的两枚兜住，case 11 的 leg 只在其一在场时响 |

**承重那句的回答（操作定义逐字答）**：
- **摘掉本票的修法那味**（`p3`＝把 corrupt 支的取值退回 `dec.InputOffset()`）⇒ 不是"打不红"，而是**case 11 当场红 8 发**
  （§2.2 那一发就是它在未修码上的同一形状）⇒ 那味**被仪器管着**。
- **摘掉 case 11 的两条 leg 之一**（`d11a`／`d11b`）⇒ 三发取值变异（`p1`/`p2`/`p3`）**仍然红**，另一条 leg 接住；
- **两条 leg 同时摘掉**（`x-p1-d11a-d11b`、`x-p2-d11a-d11b`、`x-p3-d11a-d11b`）⇒ **存在一发变异从此打不红**：
  `p1`（把腿上的取值写死成常量 0）回到 §1 的"全组 0 红"。
  ⇒ 所以**本票选定的一味就是这两条 leg 合起来**（字段 leg：`subjectReportRead.offset`；句子 leg：句子印出来的那个数），
  少任一条都有兜，两条都摘就没有牙——这是**实测**的逃逸，不是推理。
- `d11e`（整枚摘掉 case 11）× `p1` 同样逃逸 ⇒ case 11 **不是装饰**；`d13b` × `p7` 同理对 case 13。
- **有没有拿"永远绿"的断言凑数**：`p1`–`p12` 十二发生产码变异，**每一发都至少让一枚用例红**（见表），
  且 12 发里没有一枚是"靠 case 11 才能红到"以外的兜底——反过来说，**本票没有新增任何一枚今天恒绿的断言**。

**放水两问自答**：
① **断言方向**：既有的 10 枚（144 七枚＋147 三枚）**一字未改**（`git diff` 见 §5.2；`case 2` 的 `wantReason` 串、
`case 8`/`case 10` 的判据方向全部未动）；新增的三枚方向一律是"更严"。§1 那三发"零牙"读数在本票交付后**仍然用同一把尺量**（`p1`/`p7`），
没有为了变绿把任何一发从判据里撤走。`t.Skip` 枚数 0（§6.1 的 `--- SKIP=0`）。
② **helper 是不是原有的那枚**：case 11/12 用**原有的** `slo144Report` fixture 与 `readSubjectReport` 直接调用；
case 13 用**原有的** `scriptedReader`／`scriptBytes`／`scriptMissing`（票 144 造的，未改一行）；
本票新增的三枚 helper（`slo149ContradictionRe`／`slo149ContradictionOf`／`slo149Poke`）**只被新用例使用**，
其中 `slo149ContradictionOf` 与票 147 的 `slo147PairOf` 同形：**只解析、不承诺**。
生产侧新增的 `contradictionOffset` 是本票那味修法本身，不是 helper 弱化。
另：`subjectReportBudget`／`subjectGrace` **零字节未动**（§5.2 现量），新用例递给 loop 的仍是**形参**（`30*time.Second`）。

**AC#2 判定**：成立。凭据是四组：未修码上 case 11 红 8 发（§2.2）、改后同一发绿且 `=== RUN` 17→20（§2.3）、
`p1`/`p2`/`p3` 三发取值变异各自红（§2.4）、两条 leg 同摘才逃逸（`x-*-d11a-d11b` 三发 0 红）。

---

## 3　AC#3 — `s.exited()` 那一支该印什么：判定 ⓐ，并且量到"今天根本走不到"

**判据**（票面 AC#3 三选一）：ⓐ 该点名（补断言）／ⓑ 设计上不该走到 `summary()`（改成不会静默的形状＋一枚"仍走到即红"）／
ⓒ 说不好 ⇒ 停手上报。**不许按自己的偏好填**。

本程**先量可达性再判**，用的是覆盖剖面（不是推理）：

```
$ # 改前那版：HEAD 64858d6 的码＋HEAD 的测试文件，跑 144+147 全组
$ cd /tmp/wisp149/snap-pre && go test -count=1 -coverprofile=cov-pre.out -run 'TestSLO144|TestSLO147' ./cmd/wisp/
$ grep 'slo_windows.go:72[0-9]\.\|slo_windows.go:73[0-9]\.' cov-pre.out
...slo_windows.go:727.3,727.17 1 1        <- if s.exited() 那一行：执行过（条件被求值）
...slo_windows.go:728.4,730.1 1 0         <- 它的 return：执行次数 0  <- 今天整条腿一次都没走过
...slo_windows.go:732.4,734.1 1 1         <- timeout 那一支的 return：执行过（对照）

$ # 改后那版（本票交付）：同一把尺，加上 case 13
$ cd /tmp/wisp149/post && go test -count=1 -coverprofile=cov-post.out -run 'TestSLO144|TestSLO147|TestSLO149' ./cmd/wisp/
...slo_windows.go:781.3,781.17 1 1
...slo_windows.go:782.4,784.1 1 1         <- exited 那一支的 return：现在执行过了
```

（两份剖面入库为 `probes/149/cov-149-exited-branch-pre.txt` 与 `-post.txt`。
⚠ 名字不叫 `coverage-*.out` 不是随意的：`.gitignore:44` 那条 `coverage*` 会把任何 `coverage…` 开头的路径整条吃掉，
本程第一次 `git add` 时它就报了 "paths are ignored"——本程**没有用 `-f` 越过 gitignore**，改成 `cov-149-*.txt` 让原件入库，
`git check-ignore -v` 现量命中的就是那一行。）

**判定：ⓐ**，三条理由，每条指着件里的东西：
1. **设计上它会走到**：`collectReportWithin` 的注释与 `readSubjectReport` 的注释都把"subject 死在写的中途"列为
   **这一支的形状**（"The one shape this cannot tell apart from a prefix is a subject that DIED part-way through its write"），
   生产里 `collectReport` 走的就是同一枚 loop——所以**不是 ⓑ**（不存在"设计上不该走到"）。
2. **该点名**：票 144/147 收的是同一句承诺——到点那句红要说出"停在第几字节"。`exited()` 那一支是 `last.summary()` 的
   **第二个调用点**，而 `subjectReportRead` 的字段注释写的就是"loop 可以把它的最后一次读数带进放弃时的那句话里"。
3. **Ⓒ 不适用**：可达性有仪器读数（上面那两行 `1 0` → `1 1`），不需要靠偏好填。

⇒ 动作：**文本一字未动**（那句 `without writing its report (%s)` 原样保留），只补 case 13。
case 13 用真子进程（`cmd.exe /c exit 7` → `exec.Cmd.ProcessState`）让 `exited()`／`exitCode()` 走的是生产读法，
文件读数仍走票 144 原有的 `scriptedReader`（这样"当场放弃"是靠**次数**证的，不是靠秒表）。
三句判据：13a 带位置的那句必须出现、且不变成预算句、且只读 1 次；13b 带 note 的那句必须出现、且**不许**冒出 `offset`；
13c 两句不许塌成一句。

**红句原文（`p7`＝票面 `B13` 那一发，在本票交付的树上）**：

```
slo_report_144_windows_test.go:738: exited give-up "wisp slo: subject 28648 exited (code 7) without writing its report"
                                    does not carry the last reading "8 bytes read, document still open at offset 8: the tail had not arrived"
slo_report_144_windows_test.go:760: ... does not carry the note-bearing last reading "0 bytes read, no report file yet"
slo_report_144_windows_test.go:771: the two exited shapes render the same sentence: "..."
（p8 把 summary 换成常量 "last reading unavailable" 同样三句全响）
```

⇒ 与 §1 的对照读起来正好是一件事的两面：**同一发变异在 64858d6 上 0 红，在本票交付的树上 1 红、三句全响**。
⇒ `x-p7-d13a`（丢 summary ＋ 把 13a 的字面量弱化成 `"code 7"`）**仍然 1 红**——13c 那句兜住了；
`x-p7-d13b`（整枚摘掉 case 13）0 红 ⇒ case 13 是那味的唯一持有者。

**放水两问自答**：① 断言方向没动（`p9` 那一发证明本票没有把 144/147 对 `timeout` 支的断言搬走或改向：它红的是既有三枚，
**case 13 不响**——新用例不越界去钉它没被授权的东西）；② helper 全部原有（`scriptedReader`／`scriptBytes`／`scriptMissing`），
新增的只有 case 13 自己；`30*time.Second` 是**形参**，不是预算常量。

---

## 4　AC#4 — 那句注释改成它真做到的

**判据**：射程＝corrupt 支**实际**会在哪些形状上报 0；**不许**为了省事把"有反例"那半句删掉了事；
只改注释**不算收**（AC#2 才是那枚牙）。

**改前后对照**（`cmd/wisp/slo_windows.go` 的 `subjectReportRead.offset` 字段注释）：

```
-// missing part is what follows it; on reportCorrupt it is the decoder's own
-// answer, which is 0 whenever it could not begin reading a value at all.
+// ... On reportCorrupt N is the byte position the decoder objects AT, counted the way
+// the decoder counts it - the byte it objects to is included, so a bad byte at index K
+// is named K+1 (see contradictionOffset). Of the 20 shapes .../probes/149 hands in, 18
+// classify as corrupt and none of those 18 produced a name smaller than 1; the two
+// corrupt legs with no decoder error keep the end of what DID close ...
+// What this comment claimed until ticket 149 - that on reportCorrupt the number "is 0
+// whenever it could not begin reading a value at all" - was measured false the same week
+// it was written: 8 of those 18 printed 0, and two of them had a perfectly good object
+// head behind the 0 ...
```

⇒ **"有反例"那半句没被删**，反而被写成带枚数的读数（18 发 corrupt／8 发印 0／最小的名字是 1），
并具名指回 `probes/149/probe-corrupt-census.log`；旧那句错误承诺**逐字留在注释里**（本仓 append-only 的味道，
防止下一位从"没有反例"重新推出同一条规律）。
⇒ 新写的规则里**没有一枚数是推导出来的**：K+1 那条是 20 发逐枚读的（`eo` 列），最小值 1 是 `html-head`／
`brace-first-then-good`／`array-into-struct` 三发共有的读数。
⇒ 另一处文字（`contradictionOffset` 的函数注释）把**没有见证的那一味**写在明面上：
"no shape in that probe reached it, so this fallback leg carries no witness today"——
这一句**有仪器兜着**：§6.5 的覆盖剖面显示 `slo_windows.go:659.2,659.20 1 0`（fallback 那行执行 0 次）。

**只改注释不算收**——本格的凭据不在文字上，在 §2.2/§2.4：`p3`（把取值退回 `dec.InputOffset()`）当场红 8 发。

**放水两问自答**：① 断言方向无关（注释不是断言；本程**没有**把注释当修法，见 §2.1）；
② 无 helper 参与。**并如实记一句**：注释本身不受仪器管——票 147 的验收件 `CB`/`CC` 两发已经量到"改注释 0 枚用例变红"，
本程接受这个结论并**不据此认为注释可以乱写**：它的地位是"被修正的承诺文本"，牙在 §2。

---

## 5　AC#5 — 契约轴：票面点名的禁改面逐枚零字节

**判据**：`internal/risk/**`、`internal/panel/**`、`internal/agent/**`（整棵，138 在里面）、`tools/d22scan/**`、
`thresholds.go`、任何 golden、`allowlist.txt`、`scripts/slo-check.ps1`、`docs/PLAN.md`、`docs/specs/**`、
`frontend/**`、`design/**` 一律零字节；预算常量不许动。

```
$ git show --name-only --format='' b417d31 2f5c7f9 | sort -u
.scratch/wisp/probes/149/... (那两枚 commit 里 56 枚；本件那一枚再加 7 枚)   cmd/wisp/slo_report_144_windows_test.go   cmd/wisp/slo_windows.go

$ git show --name-only --format='' b417d31 2f5c7f9 | sort -u \
    | grep -E '^(internal/|tools/d22scan/|.*thresholds\.go$|.*golden.*|.*allowlist\.txt$|scripts/slo-check\.ps1$|docs/PLAN\.md$|docs/specs/|frontend/|design/)'
（空）      <- 交集为空

$ git diff --numstat 64858d6..HEAD -- cmd/wisp/
268     0       cmd/wisp/slo_report_144_windows_test.go     <- 测试文件纯追加：既有 10 枚用例一字未改
57      3       cmd/wisp/slo_windows.go

$ git diff 64858d6..HEAD -- cmd/wisp/ | grep -E '^-[^-]'
-       // missing part is what follows it; on reportCorrupt it is the decoder's own
-       // answer, which is 0 whenever it could not begin reading a value at all.
-               obs.offset = dec.InputOffset()
3 枚删除逐枚点名：两行是**被推翻的那句注释**（原文同时留在 §4），一行是搬进 contradictionOffset() 的取值。
**没有一行是断言、阈值、预算或 helper。**

$ git diff 64858d6..HEAD -- cmd/wisp/slo_windows.go | grep -E '^[+-][^-].*(subjectReportBudget|subjectGrace|thresholds)'
（空）      <- 预算常量与阈值零命中
```

**两枚别人正在写的目录**：`frontend/**`、`design/**` 本程**没碰、没还原**，也**没有**把它们的状态算进任何"零命中"宣称——
上面那句"交集为空"只覆盖"本程 commit 的文件名册与票面禁改面的交集"这一件事。
`internal/agent/**`：进场 `git status` 里没有它（§0 现量），本程也没写它；工作树那一发的门禁读数见 §6.1
（它含别人的未提交状态，所以本程只把它算作"这棵树现在能跑"的读数，不算作锚点的读数）。

**放水两问自答**：① 无断言方向变化（3 枚删除逐枚点名，全部非断言）；② 未改任何原有 helper。

---

## 6　AC#6 — 门禁（一律逐包单跑，改前改后各一次）

### 6.1　CI 同形跑法 `bash scripts/wisp-cli-tests.sh`，四数

```
改前（/tmp/wisp149/snap-pre ＝ 64858d6 全树＋补 third_party dll）:
  portable-tests.sh: four numbers (all from -v output): === RUN=133  --- PASS=73  --- FAIL=0  --- SKIP=0     rc=0  runtests.sh: OK
改后（/tmp/wisp149/post ＝ 同一棵树＋本票两枚文件，gofumpt 之后重跑一遍同形）:
  portable-tests.sh: four numbers (all from -v output): === RUN=136  --- PASS=76  --- FAIL=0  --- SKIP=0     rc=0  runtests.sh: OK
工作树那一发（含 frontend/design/台账的未提交状态，不是锚点）:
  portable-tests.sh: four numbers (all from -v output): === RUN=136  --- PASS=76  --- FAIL=0  --- SKIP=0     rc=0
```

⇒ **两边 `=== RUN` 都非零**（133／136）⇒ 都是"跑到了"；本程没有遇到 `0xc0000135`＋0 条 `=== RUN` 那一发。
区分依据按票面用的是 `=== RUN` 枚数，不是返回码。
⇒ 73→76 的差与 §6.2 的名册差集**逐枚对得上**（＋3 枚 `TestSLO149*`）。
⇒ 四枚日志入库：`probes/149/gate-pre-full.log`、`gate-post-full.log`（gofumpt **之前**那一遍）、
`gate-post-full2.log`（**最终文件那一遍**）、`gate-worktree-post.log`。
本程特意**重跑了一遍改后**，因为第一遍跑的是尚未 gofumpt 的那版测试文件——两遍四数相同（136/76/0/0），
但**以 `gate-post-full2.log` 为凭**；`gate-post-full.log` 留着不删（规则 8）。

### 6.2　名册差集（两向 `comm`，本程自己 `-list`）

```
$ go test -list '.*' ./cmd/wisp/ | grep -E '^Test' | sort   （两份树各一次，PATH 里都要有 dll，否则 -list 也是 0 枚）
pre=73  post=76
comm -23（消失）：（空）
comm -13（新增）：TestSLO149CorruptLegsWithoutADecoderErrorKeepTheirOwnEnd
                 TestSLO149CorruptSentenceNamesThePositionTheDecoderObjectedAt
                 TestSLO149ExitedGiveUpSentenceCarriesTheLastReading
```

⇒ **没有一枚既有用例被改名或被挤掉**；两份 `-list` 原文入库 `probes/149/roster-pre.txt`／`roster-post.txt`。
⚠ 一处本程自己踩到的坑，值得记名：第一次跑这两发 `-list` 时**没有把 dll 目录放进 PATH**，两份都返回 **0 枚**——
`go test -list` 也要**启动**测试二进制，所以它和"跑测试"死在同一个地方（`0xc0000135`）。
那一次的输出没有被当成读数（0 枚的 comm 两向都空＝看起来像"零消失零新增"的假绿），补上 PATH 才是上面这版。

### 6.3　判红绿只认 `^--- FAIL:`，不用裸 grep 数枚数（本程自己现量）

```
$ grep -c -- '--- FAIL' probes/149/gate-post-full2.log     ->  1      <- 裸尺（子串）
$ grep -c '^--- FAIL' probes/149/gate-post-full2.log       ->  0      <- 锚定尺
$ grep -n -- '--- FAIL' probes/149/gate-post-full2.log
588:portable-tests.sh: four numbers (all from -v output): === RUN=136  --- PASS=76  --- FAIL=0  --- SKIP=0
$ grep -c -- '--- SKIP' -> 2   ／   grep -c '^--- SKIP' -> 0
```

⇒ **那一枚假 FAIL 来自它自己的汇总行**（`gate-post-full2.log:588`），与 `-skip` 散文无关；
`-skip` 那句散文贡献的是 `--- SKIP` 多出的 2 枚。这与票 147 验收件 §4.4 的归属更正**同一读数形状**，本程独立量过。
⇒ 本件所有红绿一律按**用例名**逐个判（§2.4 那张表的"响的用例"是名字，不是枚数）。
⇒ 一枚 panic 会吞掉同包其余几十条读数：本程全部跑法里 `panic:` 命中 0（各日志 grep 现量）。

### 6.4　假红那一发：dll 不在位时"看起来像活失败"

本程按派单的要求**自己撞了一次**：§6.2 那次没带 PATH 的 `-list` 就是它的最小形状（0 枚，而不是 N 枚红）。
票 147 验收件写的"仓外快照会拿到 8 枚假红（`no native DLLs`、`PANIC=0`）"——本程**没有**复现出"8 枚红"那一发，
因为本程从第一步就把 `third_party/sherpa-onnx` 复制进快照（§0）；本程撞的是**同一件事更早就失败的那一面**
（连 `-list` 都返回 0 枚）。⇒ 那条坑**成立但形状随跑法不同**：派单说的 8 枚是"跑起来了而那几枚用例自己报错"，
本程撞到的是"进程根本起不来"。**〔未复现＝本程没量到那 8 枚，不代表它不存在〕**，写进 §7 第 4 条。

### 6.5　三把静态尺

```
$ D:/work/base/gopath/bin/gofumpt.exe --version
v0.12.0 (go1.27.1)                    <- 版本现读，没背数
$ D:/work/base/gopath/bin/gofumpt.exe -l . tools/d22scan tools/mockllm
（空）  rc=0                            <- 第一遍量到 cmd\wisp\slo_report_144_windows_test.go 一枚，
                                         gofumpt -w 只对那一枚本票文件跑过，然后重跑整把尺为空
$ go vet ./cmd/wisp/
（空）  rc=0
$ sh scripts/d22scan.sh                 rc=0
  正对照: runtests.sh: OK - packages=[./...] top-level: PASS=30 FAIL=0 SKIP=0, === RUN=70
  d22scan: examined 228 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
  bans #1-5 internal/=205  cmd/=23 | ban #6 frontend/=66 | ban #7 internal/tools/=18
  ban #8 design/=39  frontend/=66  internal/=412  cmd/=43
  d22scan: clean - no D22 ban violations
$ # 覆盖剖面里那一行（§4 引用的就是它）
...slo_windows.go:659.2,659.20 1 0     <- contradictionOffset 的 fallback 那行：执行 0 次
```

⇒ 八个作用域 `examined N` **全非零** ⇒ `clean` 不是"什么都没扫"；跑的是 **gofumpt**，没拿 `gofmt` 交差。
⚠ **分母随树**：`design/=39`、`frontend/=66` 是**工作树**的数（此刻别人正往那两棵里写东西，见 §0 的 status），
**不是锚点 `64858d6` 的数**；本程不据这两个数下任何"design/、frontend/ 零命中"的宣称（§5）。
`d22scan` 另外打了一行 `skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]`
——那是它自己的 gitignore 路径在生效（本程跑在**活仓**里，不是快照），照抄登记。

**放水两问自答**：① 门禁没有"放宽"：`-skip` 名单是 `portable-tests.sh` 原有的账，本程**一枚未加**
（三处 `--- SKIP=0` 现量）；② 三把尺跑的都是原命令原文（派单给的字面量），没有换更弱的尺。

---

## 7　本程没测什么（按"如果我漏了它谁会先被骗"排序）

1. **真实写侧、真实 subject、真实 `wisp slo` 一次都没跑**。三枚新用例全部走**纯函数** `readSubjectReport`
   ＋**脚本递的字节**（`scriptedReader`）＋注入读函数的那条 loop 路径；唯一"真"的东西是 case 13 那个 `cmd.exe` 子进程
   （它只为让 `ProcessState` 是真的，文件读数仍是脚本的）。本程没跑 `wisp slo`、没跑 `scripts/slo-check.ps1`、没跑 `slo-full`。
   ⇒ 谁先被骗：相信"生产里 subject 真死掉时到点那句会印出那枚位置"的人——本程证的是**那支现在会被走到并且带出句子**，
   不是**真崩溃会长成哪一发**。票面"不解决的事"第 2 条（真竞态本机唯一形状是 `len=0`）**依旧未取证**，本程没推进一根毛。
2. **`contradictionOffset` 的 fallback 那行没有见证**（`return inputOffset`，覆盖剖面 `1 0`，§6.5）。
   20 发探针里 corrupt 支的解码错误**全部**是 `*json.SyntaxError` 或 `*json.UnmarshalTypeError`，
   所以"错误两者都不是"那一发今天造不出来。⇒ 谁先被骗：以为 fallback 也已被测过的人。
   ⇒ 本程**没有**为它编一发形状；如果下一位造得出，那一发应当让 §2.4 的矩阵多一行，而不是让注释改口。
3. **探针那 20 发是手搓的**，不是从真实 subject 崩溃形状采的样；case 11 的 8 发同理。
   ⇒ 谁先被骗：拿 corrupt 红句排障、把"K+1"当产品不变量而不是解码器约定的人。
   K+1 是**这一版 encoding/json 的 scanner 约定**（本程 20 发逐枚读的），Go 升级可能改它——那时 case 11 会红，
   那是本程**故意留下的响声**，不是回归。
4. **没复现"仓外快照 8 枚假红"那一发**（§6.4）：本程从第一步就补了 dll，撞到的更早失败是 `-list` 返回 0 枚。
   ⇒ 谁先被骗：把派单那句读成"本程也量到 8 枚"的人——**不是**。
5. **变异不是穷举**。12 发生产变异＋9 发 leg 摘除＋11 发组合是**按味的逻辑**造的（每味要一发"摘掉就逃逸"）。
   本程没做系统扫描（全部单常量替换、全部格式串置换、全部实参置换）。特别地：
   **票面第 4 条那发"红句两实参对调"本程没有为它开任何"要它响"的格**——它在 corrupt 支上不等价（`bytes ≠ K+1`），
   本程用 `p2`（取值换成 `obs.bytes`，红 7 of 8）量到了那一支上两枚数**本来就不是同一个数**，
   但**没有**造"对调后必须响"这种判据（那是票面明令不许开的格）。
   ⇒ 谁先被骗：把"矩阵里全响"读成"所有坏形状都会响"的人。
6. **`!windows` 那一族没进容器**：`slo_windows.go` 与这枚测试文件都带 `_windows`，
   `wisp-cli-tests.sh` 自己在非 windows 上 GUARD 退出（`exit 2`）。本程**没有**跑 linux 那一腿。
7. **没跑 `-race`**；**没跑整仓 `go test ./...`**（票面地界只有 cmd/wisp，AC#6 也只点名逐包单跑）。
8. **`exitCode()` 那枚 `-1`（未退出时）没测**：case 13 的子进程一定已退出，本程没有量"没退出时那句印 code -1"的形状。
9. **case 13 依赖 `cmd.exe` 在场**（Windows 上必然；本程现量真退出码 7 并被句子带出）。
   本程**没有**处理"宿主禁用 cmd.exe"的形状——那一发会是 `t.Fatalf`（响亮），不会是静默绿。
10. **台账、票面框、HANDOVER 一字未动**：勾框与登记归编排者（票面 §Acceptance）；本程没写 `pending-and-issues.md`。

---

## 8　结论修正记录（含对本程派单与本程自己话的推翻）

**说明**：派单与票面的每一句都按未验证断言处理。没有被推翻的也在列，免得下一位把"没报错"当成"没查"。

| # | 原话（谁说的） | 本程现量 | 结论 |
|---|---|---|---|
| 1 | **票面 §现量的形状 第 1、3 条**（引验收程的 `B14`／`B13`，标〔仅自述，未复核〕） | `pre-B14-corrupt-fill-to-zero.log`、`pre-B13-exited-drops-summary.log`：在 64858d6 的码＋测试上各 17 枚 `=== RUN`、**0 红**；两发各自另有同族对照响着（`pre-B15` 红 1 枚、`pre-B5` 红 3 枚） | **复算成立**（本程自己的编号是 `p1`／`p7`，没有沿用它们的编号当凭据） |
| 2 | **票面 §现量的形状 第 2 条**："验收程实测**两发**已经在读值了仍报 0" | `probe-corrupt-census.log`：20 发里 **8 发**印 `at offset 0`，其中头一段完好文档在后的至少三发（`@` 在 index 500、0xff 在 index 40、0xff 在 index 1341） | **票面少计**（两发是它撞到的样本数，不是全集）。**动作不受影响**：本程把那 8 发全部写进注释的读数里 |
| 3 | **票 147 验收件 §3.1 那句** `whenever` "与读数不符" | 同一形状复算：`double-comma-39B` 与 `trunc40-plus-0xff` 的 `eo` 分别是 27／41，而 `z` 都是 0 | **复算成立**，且本程量到它**一直带位置**（20 发无一为空）——所以修法能落在 `eo` 上 |
| 4 | **派单 §5／票 147 验收件 §4.2**："仓外快照不带 dll 会拿到 **8 枚假红**（`no native DLLs`、`PANIC=0`）" | 本程**没有**撞出那 8 枚：第一步就把 dll 复制进快照（§0）；本程撞的是更早的一发——`go test -list` 不带 PATH 时**返回 0 枚**（§6.2）。〔我未复现那 8 枚，不代表它不存在〕 | **未按原形状复现**（同一坑、不同跑法；本程只对自己那发下结论） |
| 5 | **派单 §2 的 4 条"红句两实参对调不许开格"** | 本程**没有**为它开格；改以 `p2`（取值换成 `obs.bytes`）量"该腿的取值与实读不符"，红 7 of 8 | **相符**，按此执行；差异写在 §7 第 5 条 |
| 6 | **本程 harness.py 第一遍 combos 的 `x-p1-d11a-d11b` 读数**（本程自己写的） | 那遍的 `apply()` 把同文件的第二枚替换**写成独立 overlay**，只留下最后一枚 ⇒ 第一版那条"0 红"是**尺坏了**，不是活失败；改成串联替换后重跑（`combos2/`），同一发**仍然 0 红**（这次是两条 leg 真被摘干净） | **推翻本程自己一处跑法**（结论未变、凭据换了一版）。两遍日志都留着：`probes/149/post/x-p1-d11a-d11b.log`（坏的那遍）与 `probes/149/combos2/x-p1-d11a-d11b.log`（算数的那遍） |
| 7 | **本程 commit `2f5c7f9` 的标题句**："变异尺与逐发原始读数落盘（**57 枚**）" | `git ls-tree -r HEAD --name-only -- .scratch/wisp/probes/149/ \| wc -l` ＝ **56** | **推翻本程自己一句**（多算一枚）。已推送历史**不改写**，这一行就是追加更正；§5 与 §9 用的是现量 56 |
| 8 | **派单 §2**："此刻另一枚代理在跑票 138，落点 `internal/agent/**`" | 进场 `git status --porcelain` 里**没有** `internal/agent/**`（§0 原文）；写作过程中锚点被推走 5 枚（`d8e39579`/`69aefae9`/`c4540a6e`/`f8623c77`/`3f6322f`，全部 `docs/evidence/s1/138-*.md`＋台账，**零 Go 文件**）；本程改前读数一律按 sha `64858d6` 取，改后一律在本票两枚 commit 的树（`/tmp/wisp149/post`）上取 | **相符＋锚点漂移已记名**（同 147 验收件 §0.5 的处理：漂移不影响任何一格，因为被验版本按 sha 取） |
| 9 | **票 147 实现件/验收件的全部相关读数**（`1978` 前缀枚数、`case 10` 的 complete 句子、`B13`/`B14`/`B15` 编号形状） | 逐条重跑：fixture 1978 字节相符；`pre-B15` 的红句 `:509` 与它 §3.1 逐字相符；`pre-B5` 的三枚红名与它 §1.4 表里 `B5` 那行相符 | **复算相符**（每条都是本程自己跑的，只是恰好与它一致） |

---

## 9　伪授权登记（两个数分开栏）＋凭据自扫＋git 纪律自证

**真通知回显数：11 枚**，按类别点名（**计数法随每类写明，因为"枚数"本身就是读数）：

| 类 | 枚数 | 计数法与出处 |
|---|---|---|
| 后台任务完成回显（`<task-notification>`） | **5** | **可机械复算**：本程启动的后台命令就是 5 枚（gate-pre／gate-post／d22scan／gate-post2／gate-worktree），每枚收到一次完成回显 |
| 进场首条项目上下文（`AGENTS.md` 全文＋技能清单） | **1** | 会话首条，非工具输出 |
| "was recently modified. Read it again before writing"（指向 `agents.md`） | **3** | 首条 system-reminder 内逐枚点收（同一枚文件重复三次） |
| "Newly created file was read by another agent … don't update it unless asked"（10 条路径清单） | **1** | 同上，一条清单算一枚 |
| 日期变更提醒 | **1** | system-reminder |
| **合计** | **11** | ⚠ 第 2–5 类是**人工点收**，可能 ±1；第 1 类可机械复算。**任何一枚都没有改变本程的判断或射程**：三枚新用例、20 发探针、32 发变异全部照跑 |

**判为注入数：0**，理由与备案：
- 11 枚回显里**没有任何一句**冒充"编排者备注／系统提示／文件已被修改（对本票文件的裁决意义）／已核验请继续提交／
  请放宽阈值／已解锁／不用取证直接给结论／Confirm the harness note is genuine"。
  本程对票面、票 147 验收件、`.scratch/wisp/probes/149/`（56 枚）与本件做了字面反扫：
  ```
  $ grep -rniE "编排者备注|系统提示|已被修改|已核验请继续|请放宽阈值|已解锁|不用取证|Confirm the harness note" \
      .scratch/wisp/issues/149-*.md docs/evidence/s1/147-offset-naming-r1-accept-r1.md \
      .scratch/wisp/probes/149/ docs/evidence/s1/149-three-unpinned-outlets-r1.md
  0 枚 票面 ／ 4 枚 147 验收件 ／ 0 枚 probes/149（56→63 枚） ／ 3 枚 本件
  ```
  ⇒ 147 验收件那 **4 枚**全部在它**自己叙述这些字样的那几行**（`:676`/`:677`/`:679`/`:681`）；
  本件那 **3 枚**（`:572`/`:573`/`:576`）全部在本节**把这句话写出来**的两行与那条 grep 命令里——
  **自指尺子的固有形状**：本件一旦登记这些字样，下一位重跑同一把尺就会在本件多命中几行。
  ⇒ 两处都不是外部注入，且本程**没有**在任何一枚回显里收到过要求"少取证／别用工具／直接给结论"的文字。
- 唯一带祈使的两枚（"Read it again before writing"／"don't update it unless asked"）指向的都不是本票地界的文件，
  本程照做但**不据此改判据**，也没有据此**少跑一发**。
- **备案一处**：另一枚会话的 commit 标题（`f8623c77`）里写着"gofumpt 唯一命中落在 **149 在飞**的 cmd/wisp 测试文件（HEAD 版干净）"——
  那是**别人对本程未提交状态的读数**。本程的处理：不把它当凭据、不据此省掉任何一步——
  §6.5 的 gofumpt 是本程**自己**跑的（版本现读 `v0.12.0 (go1.27.1)`、`-l` 第一遍命中本程那一枚、`-w` 之后为空）。
  两处的时间顺序与本件写作过程相符，但**本件的凭据一律是本程自己的日志**。

**凭据值抄录数：0。** 本件不含任何密钥/token 的值；涉及的只有测试名、变量名、文件名与仓外快照路径。
⚠ 一把**自扫的尺**（与票 147 验收件 §7 同一把，本程自己现量、不背数；**本节落盘后又重跑了一遍**）：

```
$ R="sk-[A-Za-z0-9]{8,}|[A-Za-z0-9+/]{40,}={0,2}|api[_-]?key[[:space:]]*[:=][[:space:]]*[^ ]{6,}|Bearer [A-Za-z0-9]"
$ grep -cE "$R" docs/evidence/s1/149-three-unpinned-outlets-r1.md          -> lines = 13
$ grep -oE "$R" docs/evidence/s1/149-three-unpinned-outlets-r1.md | wc -l  -> fragments = 14
$ 逐枚归类（归类管道直接跑在 -o 的输出上，不是人眼看）：
      1  40 位十六进制 git object id             （§0 的 64858d6 全形）
      2  仓外快照路径（去斜杠后是一串长 alnum）  （§1.1 与 §2.3 的 harness.py 调用）
     10  Go 测试名（长驼峰串）                   （§1 的 RED 名册、§2.2 的三枚、§6.2 的名册三枚）
      1  `sk-notification`                       （`<task-notification>` 撞上 `sk-` 那一支——**尺的假命中，不是凭据**）
凭据形状且无法归入以上四类者：0
```

⇒ **结论**：抄录的凭据值 **0**，但这把尺本身**不空**（13 行／14 枚），而且**其中一枚是这把尺自己的假命中**。
把它写成"（空）"就是 147 验收件修正记录第 10 行那处同源不准，本程不重复它；
真正的字面量比对应由编排者拿凭据文件做。
⚠ 上面那两枚数**取自本节落盘时那一版**：本件每追加一节、多印一枚长测试名，这个数就会涨——复算请**重跑那两条命令**，别背 13／14。

⚠ **与本仓 §1.2"规格与仪器的射程不是一回事"同形的一条自扫**：本件用了一些 `⚠`／`⇒`／`／`／`≠`／`→` 之类字符，
`tools/d22scan` 的 ban #8 作用域是 `design/`／`frontend/`／`internal/`／`cmd/`，**不含 `docs/`**（§6.5 那八个作用域名册里
就没有 `docs/`）⇒ 本件不构成违规，但**"docs 里没有 emoji"不是被仪器核过的宣称**——它没被核，因为不在射程内。
本程写进 `cmd/wisp/**` 的每一行**都在射程内**（ban #8 `cmd/` 扫 43 枚 Go 文件含注释与 `_test.go`，§6.5 rc=0）。

**git 纪律自证**：
- 每一枚 commit 都带**显式 pathspec**（`git commit -q -F - -- <路径…>`，定界符 `'MSGEOF'` 加单引号 ⇒ 反引号未被执行）；
- 每枚 commit 前 `git diff --cached --name-only` 自核：第一枚只出现 `cmd/wisp/` 两枚、第二枚只出现 `.scratch/wisp/probes/149/` 56 枚，
  **没有出现别人的路径**；`git add -A`／`git add .`／`git commit -a` **一次没用**；
- `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`rm`／`rmdir` 使用数 **0**；**push 次数 0**；
- 临时件**只建不删**：仓外 `/tmp/wisp149/`（＝`C:/Users/swq/AppData/Local/Temp/wisp149/`）下
  `snap-pre`、`post`、`newtest-oldcode` 三枚树 ＋ `mut/<名>/` 32 枚变异目录 ＋ `logs/{pre,post,post-final,combos2,probe}/` ＋ `probe/` 探针件；
  **仓库目录内没有建过 worktree 或 checkout**；
- 变异一律走 `go test -overlay` ⇒ **任何快照树里的文件都没被改写或删过**（§2.4 那句"每发变异强制字面量恰好命中 1 次"是同一枚尺）。

---

## 10　五格自述（**不是勾框**：票面框由编排者按非实现者验收表定）

| 格 | 本程自述 | 凭据在件里的位置 |
|---|---|---|
| AC#1 | 三发今天全不响 ⇒ 本票成立、射程不缩 | §1 五行表＋两枚红句原文＋探针 20 发读数 |
| AC#2 | 会响的检＝case 11 的两条 leg；未修码上红 8 发；`p1`/`p2`/`p3` 三发各自红；两条 leg 同摘才逃逸 | §2.2、§2.4 矩阵 |
| AC#3 | 判 ⓐ（该点名），可达性用覆盖剖面量（`1 0`→`1 1`）；文本一字未动，补 case 13 | §3 |
| AC#4 | 注释改成 20 发现量支持的规则，旧那句错误承诺逐字留在件里；"有反例"那半句没删 | §4 |
| AC#5 | 本程两枚 commit 的名册与票面禁改面交集为空；测试文件 268 加／0 删；3 枚删除逐枚点名 | §5 |
| AC#6 | 改前 133/73/0/0、改后 136/76/0/0（最终文件那一遍＝`gate-post-full2.log`）；名册消失 0／新增 3；gofumpt `v0.12.0` 空、`go vet` 空、d22scan rc=0 八枚分母非零 | §6 |

**未完成＝无**（本票五格 AC 全有读数；票面"不解决的事"三条本程一条没扩大）。
**留给编排者的两处读数**（不是待办，是登记）：
① `contradictionOffset` 的 fallback 那行**今天零见证**（§7 第 2 条）——要么下一位造出那一发，要么它永远是 0；
② 覆盖剖面显示 `exited()` 那条腿现在被走到，但**真 subject 崩溃形状**仍未取证（§7 第 1 条，与 144/147 同一条洞）。

---

## 11　落盘名册与复跑法

**本程 commit 名册**（`git log --format='%h %s' 64858d6..HEAD -- cmd/wisp/ .scratch/wisp/probes/149/ docs/evidence/s1/149-*.md`）：

```
b417d31 第 1 枚  cmd/wisp 两枚（码＋三枚用例）
2f5c7f9 第 2 枚  .scratch/wisp/probes/149/ 56 枚（尺＋逐发原始读数）
（第 3 枚＝本件入库，**同一枚 commit 里带上 §6 之后新落的 7 枚件**：4 枚门禁日志＋1 枚 d22scan 日志＋2 枚覆盖剖面；
`probes/149/` 于是 56 → **63** 枚）
```

区间里另有**不是本程的** 5 枚（`d8e39579`/`69aefae9`/`c4540a6e`/`f8623c77`/`3f6322f`，全是 138 那枚会话的证据件＋台账，
**零 Go 文件**）——见 §8 第 8 行。

**入库件名册**（`find .scratch/wisp/probes/149 -type f | wc -l` ＝ **66**，现量于 §12 那一格入库前；
本件第一枚（`a055d9f`）时是 **63**——差的那 3 枚就是 §12 的 `logs-verify/`，**名册自己会随件增长**）：
变异尺 1 枚（`harness.py`）＋探针源 1 枚（`zz149probe_windows_test.go`）＋名册 2 枚（`roster-pre.txt`／`roster-post.txt`）＋
覆盖剖面 2 枚（`coverage-exited-branch-{pre,post}.out`）＋
根层原始日志 12 枚（`pre-*` 5 枚、`newtests-on-unfilled-code.log`、`probe-corrupt-census.log`、`d22scan-worktree.log`、门禁 4 枚）＋
逐发变异日志 45 枚（`post/` 31 枚 ＝12 发生产变异＋9 发 leg 摘除＋9 发组合的第一遍（坏的，见 §8 第 6 行）＋1 发 `post-asis`；
`post-final/` 3 枚＝gofumpt 之后重跑的三发；`combos2/` 11 枚＝串联替换修好之后的组合遍；
`logs-verify/` 3 枚＝§12 那一格在**入库字节**上重跑的三发）。
⚠ 这些枚数**本来就是读数**，随本件增长的那一枚 commit 一起入库；复算请一律 `find` 现量，别背这里的加总。

**复跑本件最硬那三发的最小路径**：

```
# 0) 被验版本快照（改前）＋补 dll
mkdir -p /tmp/wisp149r/snap-pre && git -c core.autocrlf=false -c core.eol=lf archive 64858d6 | tar -x -C /tmp/wisp149r/snap-pre
mkdir -p /tmp/wisp149r/snap-pre/third_party && cp -r <repo>/third_party/sherpa-onnx /tmp/wisp149r/snap-pre/third_party/

# 1) AC#1：corrupt 支换成常量 0，在"码＋测试都是 64858d6"的树上——期望 0 红、=== RUN 非零
cp <repo>/.scratch/wisp/probes/149/harness.py /tmp/wisp149r/
python /tmp/wisp149r/harness.py /tmp/wisp149r/snap-pre /tmp/wisp149r/logs/pre pre-asis pre-B14-corrupt-fill-to-zero pre-B15-complete-fill-to-zero pre-B13-exited-drops-summary

# 2) 恒真那一问：本票的测试文件＋64858d6 的码——期望 case 11 红 8 发、12/13 绿
mkdir -p /tmp/wisp149r/old && git -c core.autocrlf=false archive 64858d6 | tar -x -C /tmp/wisp149r/old
cp <repo>/cmd/wisp/slo_report_144_windows_test.go /tmp/wisp149r/old/cmd/wisp/
cd /tmp/wisp149r/old && PATH="<repo>/third_party/sherpa-onnx:$PATH" go test -v -count=1 -run 'TestSLO144|TestSLO147|TestSLO149' ./cmd/wisp/

# 3) 承重：本票交付的树上 p1/p2/p3 ＋ x-p1-d11a-d11b（前两发应各自红、组合应 0 红）
mkdir -p /tmp/wisp149r/post && git -c core.autocrlf=false archive HEAD | tar -x -C /tmp/wisp149r/post
mkdir -p /tmp/wisp149r/post/third_party && cp -r <repo>/third_party/sherpa-onnx /tmp/wisp149r/post/third_party/
python /tmp/wisp149r/harness.py /tmp/wisp149r/post /tmp/wisp149r/logs/post p1-corrupt-fill-to-zero x-p1-d11a-d11b
```

⇒ 变异尺的原件在 `.scratch/wisp/probes/149/harness.py`（**已入库**，与仓外那份同一枚字节）；
每发变异的字面替换都写在 `MUT` 里，**每次替换强制"该字面量在目标文件里恰好命中 1 次"**，否则 FATAL 退出——
同一文件的多枚替换是**串联**应用的（§8 第 6 行那一处自翻就是这个bug的登记）。
⚠ 复跑时**先把 dll 目录进 PATH**：本程量到不带它时连 `go test -list` 都返回 0 枚（§6.2）。

---

## 12　追加一格：入库那版**再跑一遍**（行号随版本，这条是本仓的既有规矩）

本件 §1–§11 落盘之后，本程把**入库的字节**重新导出成第四枚快照又跑了一遍三发关键变异。
**命令原文**：

```
$ mkdir -p /tmp/wisp149/verify && git -c core.autocrlf=false -c core.eol=lf archive HEAD | tar -x -C /tmp/wisp149/verify
$ mkdir -p /tmp/wisp149/verify/third_party && cp -r third_party/sherpa-onnx /tmp/wisp149/verify/third_party/
$ python .../harness.py /tmp/wisp149/verify .../logs/verify p1-corrupt-fill-to-zero p7-exited-drops-summary p10-multidoc-takes-error-offset
p1-corrupt-fill-to-zero            rc=1 RUN=20 topPASS=12 topFAIL=1 144red=0 147red=0 new=...:FAIL ...:PASS ...:PASS
p7-exited-drops-summary            rc=1 RUN=20 topPASS=12 topFAIL=1 144red=0 147red=0 new=...:PASS ...:PASS ...:FAIL
p10-multidoc-takes-error-offset    rc=1 RUN=20 topPASS=12 topFAIL=1 144red=0 147red=0 new=...:PASS ...:FAIL ...:PASS
```

**逐枚红句的行号**（入库那版；日志原文 `probes/149/logs-verify/`）：

```
:655  8 of 8 corrupt shapes do not name the byte position the decoder objected at; first: second-comma-at-26: …
:686  trailing-garbage: offset is 1981, want 1978 (the end of the document that did close)
:691  trailing-garbage: sentence "2003 bytes read, report corrupt: …" does not name "…a value closed at offset 1978 and 25 bytes follow"
:747  exited give-up "wisp slo: subject 3804 exited (code 7) without writing its report" does not carry the last reading "8 bytes read, …"
:769  … does not carry the note-bearing last reading "0 bytes read, no report file yet"
:780  the two exited shapes render the same sentence: "…"
:783  the exited sentences no longer separate position from note: "…"
```

⇒ **判据本身没有随版本变化**：三发的红/绿归属与 §1、§2.4 逐枚相同（`p1`→case 11、`p7`→case 13、`p10`→case 12），
`=== RUN=20`、144/147 全绿也一致。
⇒ **变的只是行号**：本件正文里引的 `:646`/`:738`/`:760`/`:771` 取自 **gofumpt `-w` 之前**那一版测试文件
（`probes/149/post/`、`post-final/` 那些日志是那一时点的），入库那版同一句落在 `:655`/`:747`/`:769`/`:780`。
按 §8 第 8 行那条本仓共识（"别人给的行号也是读数、天然带版本"）处理：**复算请一律按 sha 重取，别背本件正文里的行号。**
⇒ 本格顺带补一句 §2.3 的凭据：入库那版的整组 SLO 用例在新鲜导出树上跑 **13 枚顶层 PASS／0 红**
（`go test -count=1 -v -run 'TestSLO144|TestSLO147|TestSLO149' ./cmd/wisp/`，与 §6.1 的 136/76/0/0 同向）。

**放水两问自答**：① 断言方向未动（本格只重跑、未改任何判据）；② 未新增 helper（三发用的是 §2.4 那批 `MUT` 条目原文）。

---



