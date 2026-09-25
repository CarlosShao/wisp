# 147 — 实现程证据件 r1：`unwritten` 分支的 offset 判定 ＋ 一枚会响的检

- 票面：`.scratch/wisp/issues/147-collectreports-named-red-says-document-stopped-at-offset-n-but-the-unwritten-branch-always-carries-0-and-no-assertion-pins-it.md`
- 本件作者：**实现程**（写码代理）。票面五框**不由本件勾选**（`SPEC-12 §4.3` #1/#3、`AGENTS.md §0.3`：裁决者≠实现者）。
- 锚点：分支 `dev`、`git rev-parse HEAD` = `80fa0551e137384facdececb15b25cc01850cd25`（编排者简报里的 `80fa0551` **相符**）。
- 本机仪器版本（现跑，不背读数）：
  `go version` → `go1.27.1 windows/amd64`；
  `D:/work/base/gopath/bin/gofumpt.exe --version` → `v0.12.0 (go1.27.1)`（票面写的 v0.12.0 **相符**）。

---

## 0　进场第一读（简报要求贴原文，不是"已看过"）

```
$ git log --oneline -12
80fa055 docs(台账 A261 + Q-53/Q-54): 票 146 结线（五框全勾、两份非实现者表背书）；三句措辞射程被改准；replay 那扇门登成待人批准
d08999a ticket(146 结线＋148 新立 blocked): 五框全勾（两份非实现者表背书）；同族第二扇门 replay 现量到、但它卡在两枚人工批准上
b6d8e31 docs(台账 A260): 票 144 验收五格全成立、五格之外两枚退回——其中一枚推翻我 20:4x 抄进台账的两句，今天第三次同形错
e8ee651 ticket(147,新立): collectReport 到点红句说"文档停在 offset N"，而 unwritten 分支那个字段恒为 0——话说了、字段没带那个信息，且没有任何断言钉住（删掉这半句七枚全绿）
3738d3b evidence(146 二轮验收 r2): 打⑤(更正段判定) + 两格裁决表 + 三件空格裁定建议 —— AC#1/AC#2 均成立
e5e79bf evidence(144 验收 第七批/行为自纠): 登记本程自己违过一次"任何 rm 都不许"的硬规矩
6e6c8e9 evidence(144 验收 第六批/收笔): 禁改面与地界逐枚现量 + 三枚行号按 post 树改正 + next=
ef48505 evidence(146 二轮验收 r2): 打③(五支单点回退) + 打④(DeepEqual 非恒真) + 常规必查
032c6c2 evidence(144 验收 第五批/收口): 五格出档 + 本程造出的两枚五格外的退回 + 单独回答"单跑够不够"
7ae767f evidence(146 二轮验收 r2): 副本自证 + 打①(五行原样复跑) + 打②(非空夹具三形)
af03030 evidence(144 验收 第四批): AC#1 / AC#3 出表（13 枚单点变异 + 1 枚组合）；offset 那一格判"读数成立、不改退回"
0899eb6 evidence(144 验收 第三批): AC#2 / AC#5 出表，另造出一枚本程自己的退回项 —— a91d7c2 不是"只有空白"

$ git status --porcelain
D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
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
?? .zcodeignore
?? design/doubao/01-ball-states.jpg
?? design/doubao/demo/lib/
?? design/doubao/demo/screenshots/
?? design/old/
```

⇒ 简报第 1 节所述**相符**：`design/**` 下 16 枚 owner 自己挪动的未提交删除（首列空格＋`D`＝worktree 删除、未 staged），
加 `.zcodeignore`／`design/doubao/`／`design/old/` 等未跟踪件。**本程一枚未碰、一枚未 add、一枚未还原**；
本程落盘只用显式 pathspec（§4 名册与每枚 `git show --name-only` 可复算）。
本程自己新增的未跟踪件只有 `.scratch/wisp/probes/147/`（临时读数，只建不删）与仓外快照 `/d/work/workspace/snap147-1`。

**简报与票面无冲突**（唯一需要点名的差别：票面 AC#3 说"三处"，本程现量是**四处文本＋一枚 commit 标题**，见 §3）。

---

## 1　AC#1　判定：ⓐ（字段该带信息）

### 1.1　票面"现量的形状"三条，本程一律现量

判据原文（票面 §AC#1）："**ⓐ 字段该带信息** ⇒ `unwritten` 分支上把 `o.offset` 真填成'已读到的最后一字节位置'，
并让红句两个数各自负责一件事"、"**不许默认它是设计**"、"判据要可核：…**MG 那一发…现在该红还是该绿必须能一句话说清并跑给人看**"。

**第 1 条（红句承诺了字节位置）：成立。** 锚点树 `cmd/wisp/slo_windows.go:591`（本程现读到原文；票面写 `:594-595`，
行号按锚点树漂移，内容一字不差）：

```go
	case reportUnwritten:
		return fmt.Sprintf("%d bytes read, document still open at offset %d: the tail had not arrived", o.bytes, o.offset)
```

**第 2 条（`unwritten` 上 offset 不带信息）：成立，且本程把它量成一张表。**
探针落在**仓外**快照 `/d/work/workspace/snap147-1/cmd/wisp/zz147probe_windows_test.go`（`git archive HEAD` 造的树，
不在工作树内、不参与本仓任何门禁跑法；票 144 验收件的 `snap/` 同法），跑法与本仓同形（PATH 挂 `third_party/sherpa-onnx`）：

```
$ cd /d/work/workspace/snap147-1 && export PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:$PATH" \
  && go test ./cmd/wisp/ -run 'TestP147' -v -count=1     # 日志 snap147-1/probe-offset-shape.log，rc=0
P147 fixture length = 1978 bytes (positive control: must be > 0 and non-trivial)
P147 readings over 0..len inclusive = 1979, state tally: map[complete:1 unwritten:1978]
P147 state=unwritten obs.offset   distinct=1 -> 0(1978)
P147 state=unwritten decoder own =0(1978)
P147 state=corrupt   never reached
P147 state=complete  obs.offset   distinct=1 -> 1978(1)
P147 state=complete  decoder own =1978(1)
P147 positive control: whole document -> state=complete offset=1978 bytes=1978
P147 prefix 1     state=unwritten bytes=1     obs.offset=0     summary="1 bytes read, document still open at offset 0: the tail had not arrived"
P147 prefix 44    state=unwritten bytes=44    obs.offset=0     summary="44 bytes read, document still open at offset 0: the tail had not arrived"
P147 prefix 100   state=unwritten bytes=100   obs.offset=0     summary="100 bytes read, document still open at offset 0: the tail had not arrived"
P147 prefix 989   state=unwritten bytes=989   obs.offset=0     summary="989 bytes read, document still open at offset 0: the tail had not arrived"
P147 prefix 1977  state=unwritten bytes=1977  obs.offset=0     summary="1977 bytes read, document still open at offset 0: the tail had not arrived"
```

⇒ 三件事同时钉住：夹具 1978 字节；**1979 发读数**里 1978 判 `unwritten`、只有整份判 `complete`；
`unwritten` 上 `obs.offset` 的**唯一取值是 0**，而**同一批字节上解码器自己的 `InputOffset()` 也是恒 0**
（探针把两者分开量，所以"换成自己算"不是把有信息的数换成没信息的数——**那一支上原值本来零比特信息**）。
**正控不空转**：同一把尺在 `complete` 上报得出 1978。

**一处与上一程读数不同的现量（登记，不影响判定）**：
票 144 验收件 `:365` 写"顺带本程量到 `corrupt` 分支也报 0——`offset` 这一味只在 `complete` 时携带信息"。
本程把 7 枚 corrupt 形状逐枚量（同一探针第二发，日志同上）：

```
P147 corrupt 43   bytes state=corrupt   obs.offset=0     decoder=0     decoderErr=invalid character '<' looking for beginning of value
P147 corrupt 39   bytes state=corrupt   obs.offset=0     decoder=0     decoderErr=invalid character ',' looking for beginning of object key string
P147 corrupt 13   bytes state=corrupt   obs.offset=13    decoder=13    decoderErr=json: cannot unmarshal number into Go struct field sloRun.mode of type string
P147 corrupt 2003 bytes state=corrupt   obs.offset=1978  decoder=1978  decoderErr=<nil>
P147 corrupt 57   bytes state=corrupt   obs.offset=57    decoder=57    decoderErr=<nil>
P147 corrupt 2    bytes state=corrupt   obs.offset=2     decoder=2     decoderErr=<nil>
P147 corrupt 72   bytes state=unwritten obs.offset=0     decoder=0     decoderErr=unexpected EOF
```

⇒ `corrupt` 上 offset **不是恒 0**（13／1978／57／2 都带位置），所以"只在 complete 携带信息"这句**射程偏窄**；
真实形状是：**只有 `unwritten` 恒 0**（票面第 2 条正好说的是这一支，成立），
而 corrupt 里**语法错误那一发**（39 字节 `{"mode":"…",,"pass":true}`，错在第 25 字节附近）也报 0——
**同一族缺陷在 corrupt 支上还有一发**。本程**没顺手改**：票面 AC#1 ⓐ 的射程写的是 `unwritten` 分支，
"这一枚票不解决的事"一节又明写"登记，别顺手扩大"。⇒ **归票面 §"不解决"继续挂着，等编排者裁**（§6 第 1 条）。

**第 3 条（没有任何断言钉住它）：成立，本程在锚点树复跑那一发 MG。**
落点＝仓外快照（改的是快照里的副本，仓内 `cmd/wisp/slo_windows.go` 一字未由本程变异过）：

```
$ # snap147-1/cmd/wisp/slo_windows.go:595 只撤 `at offset %d` 那半句，其余一字不动
- return fmt.Sprintf("%d bytes read, document still open at offset %d: the tail had not arrived", o.bytes, o.offset)
+ return fmt.Sprintf("%d bytes read, document still open: the tail had not arrived", o.bytes)   // MUTATION-147-MG
$ go test ./cmd/wisp/ -run 'TestSLO144' -v -count=1     # 日志 snap147-1/mg147-pre-1.log
rc=0
=== RUN = 14   --- PASS = 14   --- FAIL = 0   --- SKIP = 0
```

⇒ **修前（锚点树）：删掉那半句，票 144 那 7 枚用例（14 枚含子用例）一枚都不红。** 票面的第 3 条与
`A258④` 那条"绑谁改谁重量"没牙的裁定，本程都复算为真。
跑完按原样还原（`cmp` 与仓内 HEAD 副本逐字节相同）；快照里另存 `/d/work/tmp147/slo_windows.go.orig`，**没有任何文件被删**。

### 1.2　判定与那一发 MG 现在该红还是该绿（一句话说清）

**判 ⓐ**：`unwritten` 分支上 `obs.offset` 改填"文档停在第几字节"＝**已落盘的那一个字节位置**（对前缀而言就是 `len(data)`，
因为截断的文档每一枚已到达字节都被消费掉了），红句**文本一字不改**，改的是喂给它的那个数。

**判据是哪句可核的话**：票面 AC#1 ⓐ 原文"字段该带信息 ⇒ 把 `o.offset` 真填成'已读到的最后一字节位置'，并让红句两个数各自负责一件事"
＋同节"ⓐ／ⓒ 才是本票默认希望的方向；**ⓑ 单独存在不算收**"。本程没走 ⓑ（撤那半句），因为它量明是**零仪器**；
也没把 ⓒ 当独立支：那三语义里 `-1`／`0`／`N` 的 `0` 正是缺陷本身，钉住"恒 0 成立"＝钉住假话。

**一句话**：**同样那发 MG（删掉 `at offset %d`），修前 0 红／14 绿，修后必须 3 红**——
因为钉它的断言现在存在了；实测读数在 §2.3（本程真跑了，不是推）。

**`offset: -1` 那几处的现量（票面写"另外两处初始化"）**：锚点树实际是**四枚** `offset: -1`——
`slo_windows.go:631`（`readSubjectReport` 头上）、`:705`（`last` 的"nothing read yet"）、`:711`（"no report file yet"）、
`:716`（"report file unreadable"）。票面"两处"**少计两枚**（它点名的两处是 loop 里的 ②③）。
本程现量的形状：探针第三发量明**这四枚的 `-1` 从不进句子**——带 `note` 的读数走 `summary()` 第一支
（`"%d bytes read, %s"`），只印 note、不印 offset：

```
P147 note="nothing read yet"    summary="0 bytes read, nothing read yet"    contains-(-1)=false
P147 note="no report file yet"  summary="0 bytes read, no report file yet"  contains-(-1)=false
P147 note="report file unreadable: open x: …"  summary="0 bytes read, report file unreadable: …"  contains-(-1)=false
```

⇒ "这个字段今天同时装着三种语义"这一问，本程的答案是：`-1`＝**没有解码器跑过**（且从不外印）、
`N`＝**文档停在第 N 字节**（`unwritten` 与 `complete` 两支各自给出自己的 N，句子不同）；
旧形状里那个"0＝读到一半"**不是设计**，是 `dec.InputOffset()` 在截断输入上的副产物（§1.1 第 2 条量明）。

### 1.3　注释同步（票面禁止"把注释留着、字段以后再说"）

- `subjectReportRead.offset` 字段注释（`slo_windows.go:578-586`）重写为"它到底装哪两种语义、`corrupt` 上它是解码器自己的答案、
  `-1` 那一支从不外印"——**它原先只说了 `-1` 与"解码器停在哪"，那正是被本票推翻的说法**。
- `collectReportWithin` 头上那句"the sentence names the budget it spent, the byte count of its last read and
  **where inside those bytes the document stopped**"：ⓐ 之后**这句实现真做到了**，所以**保留原文**、不删。
  票面那句"不许用把注释留着来收"针对的是"注释承诺、实现不做"；本程是把实现做到，不是把注释抹平。
- `subjectReportRead` 结构体头上"where inside them the **decoder** stopped" → "…the **document** stopped"（一词，
  因为 `unwritten` 上存的不再是解码器的那个数）。

### 1.4　放水两问自答（AC#1）

- **断言方向动没动？** 没动，也没新增放宽：本程**没改任何一条既有断言**，`git diff --numstat` 现量
  `cmd/wisp/slo_report_144_windows_test.go` ＝ **169 加／0 删**（0 删＝现有文本一行未抹，含 import 块也只是插进两行）。
  生产码 `cmd/wisp/slo_windows.go` ＝ **21 加／4 删**，4 枚删除全在 §1.3 列的三处注释＋那一支 `obs.offset = dec.InputOffset()`
  的移位（挪到 EOF 判定之后，`corrupt` 分支的值不变，因为两次取数之间没有再读流）。
- **helper 是不是原有那枚？** 判定这一格不涉及新代码；本程新写的两枚 helper（`slo147UnwrittenRe`／`slo147PairOf`）
  在 §2 里逐名交代。

---

## 2　AC#2　一枚会响的检（本票硬核心）

### 2.1　造了什么

三名新用例，追加进票面地界点名的 `cmd/wisp/slo_report_144_windows_test.go`（**不新建文件**）：

| 名 | 问什么 | 用了哪枚 helper |
|---|---|---|
| `TestSLO147UnwrittenSentenceNamesTheOffsetTheDocumentStoppedAt` | 对夹具**全部** 1978 枚 `unwritten` 前缀，句子必须**整体**长成 `A bytes read, document still open at offset B: the tail had not arrived`，且 `A == B == 该前缀长度`，并**另查字段本身** `obs.offset == 该长度` | 原有 `slo144Report`；新 `slo147UnwrittenRe`＋`slo147PairOf` |
| `TestSLO147LoopGiveUpSentenceAgreesWithTheBytesItRead` | 人真读到的那一句——`collectReportWithin` 的到点错误里嵌的那一发，两个数都必须等于脚本递给 loop 的字节数（8） | 原有 `scriptedReader`／`scriptBytes`／`scriptMissing` |
| `TestSLO147OffsetSemanticsRenderThreeDifferentSentences` | 三语义**各说各的话**：没有解码器跑过时**句子不许出现 `offset` 这个词**；半截文档印 `at offset 989`；完整文档印 `complete document at offset 1978`；三句两两不等 | 原有 `slo144Report`／`scriptedReader` |

新写的那两枚 helper 只**解析**、不含判据（`slo147PairOf` 返回 `(bytes, offset, whole, ok)`，
`whole` 表示匹配是否覆盖整句）——期望值全在 `switch` 里，所以"判据写反"不会藏在 helper 里。
字段那一枚 leg 的作用是关住一条后门：**把句子印成 `%d … offset %d", o.bytes, o.bytes`** 也能骗过纯文本检，
但它骗不过 `obs.offset != int64(i)`。

### 2.2　恒真那一问：这一发**今天（改之前）响不响**——响，凭据如下

**改前红句原文**（同一枚最终版测试文件，生产码＝锚点树原样；快照跑，日志
`/d/work/workspace/snap147-1/red-before-fix-final.log`）：

```
$ go test ./cmd/wisp/ -run 'TestSLO147' -count=1 -v        # 生产码 = 80fa0551 原样
rc=1
=== RUN   TestSLO147UnwrittenSentenceNamesTheOffsetTheDocumentStoppedAt
    slo_report_144_windows_test.go:446: 1977 of 1978 unwritten sentences do not name the offset the document stopped at; first: prefix of 1 bytes: subjectReportRead.offset is 0
--- FAIL: TestSLO147UnwrittenSentenceNamesTheOffsetTheDocumentStoppedAt (0.02s)
=== RUN   TestSLO147LoopGiveUpSentenceAgreesWithTheBytesItRead
    slo_report_144_windows_test.go:470: give-up error says 8 bytes read, document open at offset 0; want both 8 (the prefix the loop was handed)
--- FAIL: TestSLO147LoopGiveUpSentenceAgreesWithTheBytesItRead (0.06s)
=== RUN   TestSLO147OffsetSemanticsRenderThreeDifferentSentences
    slo_report_144_windows_test.go:502: unfinished sentence says 989 bytes, offset 0; want both 989
--- FAIL: TestSLO147OffsetSemanticsRenderThreeDifferentSentences (0.04s)
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.182s
```

⇒ **不是"改后全绿所以可能一次没走"**：三枚红都在改前取到，红名＋红句原文在此。
（注意 `1977 of 1978` 那个数：**恰好一枚今天"过"**＝0 字节那一发，它的 offset 与 bytes 同为 0；
本程不藏它，因为那正是"offset 恒 0"这一味的形状。）

**改后同一发变绿**（仓内真树，`bash` 同形跑法，日志 `.scratch/wisp/probes/147/green-after-fix-r2.log`）：

```
$ go test ./cmd/wisp/ -run 'TestSLO147|TestSLO144' -count=1 -v
rc=0
=== RUN = 17   --- PASS = 17   --- FAIL = 0   --- SKIP = 0
--- PASS: TestSLO144EveryPrefixOfARealReportIsUnwrittenNotCorrupt (0.02s)
--- PASS: TestSLO144ReportsThatContradictThemselvesAreCorruptNow (0.00s)
--- PASS: TestSLO144LoopRetriesAnUnfinishedFileAndReadsTheWholeReport (0.03s)
--- PASS: TestSLO144LoopGiveUpSentencesOnRealFiles (0.05s)
--- PASS: TestSLO144CorruptReportIsJudgedOnTheFirstRead (0.00s)
--- PASS: TestSLO144UnfinishedReportKeepsPollingThenNamesBudgetAndBytes (0.06s)
--- PASS: TestSLO144ReportThatArrivesAfterMissingReadingsIsCollected (0.00s)
--- PASS: TestSLO147UnwrittenSentenceNamesTheOffsetTheDocumentStoppedAt (0.02s)
--- PASS: TestSLO147LoopGiveUpSentenceAgreesWithTheBytesItRead (0.06s)
--- PASS: TestSLO147OffsetSemanticsRenderThreeDifferentSentences (0.04s)
```

（17 ＝ 10 枚顶层 ＋ `TestSLO144ReportsThatContradictThemselvesAreCorruptNow` 的 7 枚子用例。）

### 2.3　承重那一句：把 ⓐ 的修法摘掉，这一发是不是从此打不红？——**打得红**，两形都跑给人看

两形都在**快照**里跑（仓内码从未带变异），最终版测试文件，每形跑完 `cmp` 回原：

**形 A：只摘修法、留句子**（`obs.offset = dec.InputOffset()` 塞回 `unwritten` 那一支）
```
$ go test ./cmd/wisp/ -run 'TestSLO147|TestSLO144' -count=1 -v   # snap147-1/mut-A2-final.log
rc=1
=== RUN = 17   --- PASS = 14   --- FAIL = 3   --- SKIP = 0
    slo_report_144_windows_test.go:446: … first: prefix of 1 bytes: subjectReportRead.offset is 0
    slo_report_144_windows_test.go:470: give-up error says 8 bytes read, document open at offset 0; want both 8 …
    slo_report_144_windows_test.go:502: unfinished sentence says 989 bytes, offset 0; want both 989
```

**形 B＝票面那一发 MG**（只撤 `at offset %d` 那半句；修法保留）
```
$ go test ./cmd/wisp/ -run 'TestSLO147|TestSLO144' -count=1 -v   # snap147-1/mut-B2-final.log
rc=1
=== RUN = 17   --- PASS = 14   --- FAIL = 3   --- SKIP = 0
    slo_report_144_windows_test.go:446: 1978 of 1978 unwritten sentences do not name the offset the document stopped at;
        first: prefix of 0 bytes: sentence carries no byte/offset pair: "0 bytes read, document still open: the tail had not arrived"
    slo_report_147…:467: give-up error "…within 60ms (last read: 8 bytes read, document still open: the tail had not arrived)" does not read as bytes-then-offset
    slo_report_144_windows_test.go:500: unfinished sentence "989 bytes read, document still open: the tail had not arrived" is not the unwritten shape
```

⇒ **答一句**：摘掉 ⓐ 的修法（形 A）这一发**仍然打得红**；把句子里那半句拿掉（形 B）**也打红**——
而形 B 在 §1.1 的锚点树上是 **0 红／14 绿**。**同一发变异、修前 0 红、修后 3 红**，这就是本票要的牙。
两形里票 144 那 7 枚（14 计数）**都保持绿** ⇒ 牙长在新增断言上，不是把既有断言改反方向换来的。

### 2.4　放水两问自答（AC#2）

- **断言方向动没动？** 没动。既有 10 枚（含 7 枚子用例）一字未改（§1.4 的 **0 删** 读数）；
  新增三枚的极性全是"更严"：从"句子存在"抬到"句子两个数各自等于第几字节"，另加"没有解码器跑过时不许印位置"。
- **helper 是不是原有那枚？** 夹具与注入面**全部复用原有那枚**：`slo144Report`（同一份 `json.MarshalIndent(sloRun{…})`）、
  `scriptedReader`／`scriptBytes`／`scriptMissing`、`sloSubject{readReportFile:…}` ＋ 生产 `collectReportWithin`。
  本程新写**只两枚、且不含期望值**：`slo147UnwrittenRe`（正则）与 `slo147PairOf`（解析）。
  **没有**为新用例另造一份夹具、**没有**用 mock 代替真的读文件路径。
- 补一句自问自答：**有没有加 `t.Skip`？没有**（三枚都跑，`--- SKIP = 0` 在每份日志里）。
  **有没有碰预算常量？没有**：`subjectReportBudget`／`subjectGrace` 一字未动，
  新用例递给 `collectReportWithin` 的是**参数**（`60ms`／`40ms`），与票 144 用例 6／7 同一形写法。

---

## 3　AC#3　那三处（现量是四处文本＋一枚标题）"只有空白"的声明：就地追加更正

判据原文（票面 AC#3）："**`>` 追加、原句一字不抹**（append-only 同样管票面与证据件）… **不许**因此回退那枚夹具或改动任何断言。"

### 3.1　本程现量（不抄编排者的数）

```
$ git log -1 --format='%H %s' a91d7c2
a91d7c2d07432e19975711d86b697ea2dcc9cfd9 fix(144 AC#4): 前一程留下的测试文件不是 gofmt/clean —— 复量与票面不符，只改对齐
$ git diff -w --numstat a91d7c2^ a91d7c2
31	2	cmd/wisp/slo_report_144_windows_test.go
$ git diff --numstat a91d7c2^ a91d7c2
42	13	cmd/wisp/slo_report_144_windows_test.go
$ git branch --contains a91d7c2 → * dev      （已推送：`git log origin/dev` 顶点即本锚点 80fa055）
```

⇒ 编排者 21:0x 现量的 **31 加／2 删** **相符**；"前后两版 `diff -w` 逐字相同"**为假**（`-w` 后仍有 31/2）。
`git diff -w a91d7c2^ a91d7c2` 全文本程读过，真实形状与票面 §AC#3 说的两味一致，逐枚列：
① 夹具 `slo144Report` 从 6 枚字段扩到含 `SubjectPID`、`ObserverCost`、`MemMedianBytes`、`CPUMeanPercent`、`GDIMax`、
`WriteOpsTotal`、`SampleErrors`、`Pass` ＋ `Samples`（2 条 `observe.Sample`）＋ `Verdicts`（2 条 `observe.Verdict`）；
② `slo144Report` 头上新增一枚 6 行注释块（"built by CALLING THE SAME SERIALIZER…"）；
③ 断言那一发：删
```go
	if !strings.Contains(string(doc), `"report"`) {
		t.Fatalf("fixture has no report field, so it proves nothing: %s", doc)
```
换成"长度下限 ＋ 5 枚字面量"（`len(doc) < 300` ＋ `"report"`／`"samples"`／`"verdicts"`／`"pass": true`／`"cpu_percent": 0.4`）。

**"实质无放水"这一句本程复算为真**，理由是极性与射程：新检**包含**被删那枚的 `"report"` 字面量、
再多加 4 枚字面量与一枚长度下限 ⇒ **严格更严**，方向未反。（其余 21 枚 `-w` 增行是 ①②的结构体字面量与注释。）

### 3.2　三处（四处）落点与已追加的更正

| 落点 | 原句（一字未抹） | 本程动作 |
|---|---|---|
| `.scratch/wisp/issues/144-….md:189`（Progress log 内） | "改动只有空白（前后两版 `diff -w` 逐字相同）" | 同一 Progress log **末尾追加**一条（新时间戳、`>` 引原句） |
| `.scratch/wisp/issues/144-….md:260`（"续程收笔自证"节） | "唯一一次代码面改动是 `a91d7c2` 那枚 `gofumpt -w`，只有空白" | **并入同一条追加**点名（票面只数了 Progress log 一处 ⇒ 本程现量补一枚） |
| `docs/evidence/s1/144-slo-report-partial-read-r1.md:375`（§4.2 末） | "**改动只有空白**——改前后两版 `diff -w` 逐字相同，断言、预算、Skip 一个没动" | §4.2 **就地、该段末尾**追加 `>` 更正块（不新开节、不动 §4.3 之后） |
| commit `a91d7c2` 标题"…只改对齐" | 已推送（`origin/dev` 顶点＝本锚点） | **不改写**（`AGENTS.md §1.4`）；只在票面追加与本件 §3 登记"那枚标题的说法不准" |

⇒ 追加动作落在**同一枚 commit**里与本票的码改动同批（票面 AC#3："与本票同批做，因为它就在同一枚文件族里"）。

---

## 4　AC#4　契约轴自核（零字节面）

见 §8 落盘名册：每枚 commit 都带显式 pathspec，`git show --name-only` 可复算。
本程**只写**过这四枚路径：

- `cmd/wisp/slo_windows.go`
- `cmd/wisp/slo_report_144_windows_test.go`
- `.scratch/wisp/issues/144-….md`（追加）
- `.scratch/wisp/issues/147-….md`（仅 Progress log 追加一条）
- `docs/evidence/s1/144-slo-report-partial-read-r1.md`（§4.2 追加）
- `docs/evidence/s1/147-offset-naming-r1.md`（本件）
- 临时读数：`.scratch/wisp/probes/147/*`（**只建不删**，见 §7 末）、仓外 `/d/work/workspace/snap147-1`、`/d/work/tmp147/*`

票面 AC#4 点名的禁改面（`internal/risk/**`、`thresholds.go`、golden、`allowlist.txt`、`scripts/slo-check.ps1`、
`tools/d22scan/**`、`docs/PLAN.md`、`docs/specs/**`、`internal/panel/**`、`internal/agent/approval/**`、
`frontend/**`、`design/**`）**逐枚现量为零字节**，读数在 §5.4。

---

## 5　AC#5　门禁（一律逐包单跑）

判据原文（票面 AC#5）："`go test ./cmd/wisp/` 改前改后各一次，四数之外**点名册差集**（两向 `comm`）；
`gofumpt -l . tools/d22scan tools/mockllm` 返回空（v0.12.0、与 CI 同版）；`go vet ./cmd/wisp/` 空；
`sh scripts/d22scan.sh` rc=0 且各作用域 `examined N` 非零。"

### 5.1　改前那一次（同一枚跑法、同一台机器）

```
$ bash scripts/wisp-cli-tests.sh          # 日志 .scratch/wisp/probes/147/gate-pre-1-baseline.log
runtests.sh: OK - packages=[./cmd/wisp/ -count=1 -skip ^(TestDefaultDeadlineWallClockMeasurement|…|TestSyncRegistryProbeLive)$]
             top-level: PASS=70 FAIL=0 SKIP=0, === RUN=130, '[no tests to run]'=0
portable-tests.sh: four numbers (all from -v output): === RUN=130  --- PASS=70  --- FAIL=0  --- SKIP=0
rc=0   （包体耗时 76.465s）
```

### 5.2　改后那一次

```
$ bash scripts/wisp-cli-tests.sh          # 日志 .scratch/wisp/probes/147/gate-post-1.log
runtests.sh: OK - packages=[./cmd/wisp/ -count=1 -skip ^(...同上 6 枚...)$]
             top-level: PASS=73 FAIL=0 SKIP=0, === RUN=133, '[no tests to run]'=0
portable-tests.sh: four numbers (all from -v output): === RUN=133  --- PASS=73  --- FAIL=0  --- SKIP=0
rc=0   （包体耗时 81.456s）
```

⇒ **`=== RUN` 两发都非零（130／133），所以这两发都是"跑到了"，不是"根本没跑到"**。
票面点名的 `0xc0000135` 那一发（`=== RUN`=0）**没有出现**：本程跑 `cmd/wisp` 一律先
`export PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:$PATH"`，或走 `scripts/wisp-cli-tests.sh`。

**名册差集（两向 `comm`，改前 vs 改后，取的是 `^=== RUN   <名>` 的顶层名去重）**：

```
pre=70 枚   post=73 枚
--- 消失（在 pre 不在 post）---      （空）
--- 新增（在 post 不在 pre）---
TestSLO147LoopGiveUpSentenceAgreesWithTheBytesItRead
TestSLO147OffsetSemanticsRenderThreeDifferentSentences
TestSLO147UnwrittenSentenceNamesTheOffsetTheDocumentStoppedAt
```

⇒ 差集里**每一枚为什么在**：三枚都是本票 AC#2 新增的用例（§2.1 表）；**消失 0 枚** ⇒ 没有任何一条既有用例
从名册上掉下去（一枚 panic 会吞掉同包其余几十条读数，那样 `comm` 会报出成片的"消失"——这里没有）。
`--- FAIL` 按 `^--- FAIL` 锚定取是 **0 行**（两发都是），`--- SKIP` 锚定也是 0 行。

**⚠ 一把要认得的假尺**：对本件那两份日志**直接** `grep -c '--- FAIL'` 会得到 **1**、`grep -c '--- SKIP'` 得到 **2**，
`=== RUN` 得到 136——**那不是红**。多出来的都来自 `portable-tests.sh` 自己打印的 -skip 说明文字：
`gate-post-1.log:11` 那句里有字面 `--- SKIP at syncdirs_windows_test.go:133`，
`gate-post-1.log:582` 是本仪器自己的汇总行（`… --- FAIL=0 --- SKIP=0`）。
两把尺的对照（同两份日志，锚定行首 vs 不锚定）：

```
                naive(子串)     anchored(行首)
gate-pre-1      RUN=133 PASS=131 FAIL=1 SKIP=2   →  RUN=130 PASS=130 FAIL=0 SKIP=0
gate-post-1     RUN=136 PASS=134 FAIL=1 SKIP=2   →  RUN=133 PASS=133 FAIL=0 SKIP=0
```

⇒ 判红绿只认 `--- FAIL:` 行（且 `t.Logf` 也带 `file:line:` 前缀，别拿它当红）；
`PASS=130/133` 那一味是**含子用例**的行计数，`runtests.sh` 报的 `PASS=70/73` 才是**顶层枚数**，两者不是一把尺。

### 5.3　另外三发

```
$ D:/work/base/gopath/bin/gofumpt.exe --version     → v0.12.0 (go1.27.1)   （票面写的 v0.12.0 相符）
$ D:/work/base/gopath/bin/gofumpt.exe -l . tools/d22scan tools/mockllm     → 空，rc=0
$ go vet ./cmd/wisp/                                                     → 空，rc=0
$ sh scripts/d22scan.sh                             rc=0                   # 日志 .scratch/wisp/probes/147/d22scan-post.log
  runtests.sh: OK - packages=[./...] top-level: PASS=30 FAIL=0 SKIP=0, === RUN=70     ← 正对照（证明这把尺能红）
  d22scan: examined 228 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
  ban #1-5 internal/=205  cmd/=23 | ban #6 frontend/=52 | ban #7 internal/tools/=18
  ban #8 design/=35  frontend/=52  internal/=412  cmd/=43
  d22scan: clean - no D22 ban violations
```

⇒ **各作用域 `examined N` 全非零**（最小的一味是 `cmd/`=23 与 `ban #8 cmd/`=43），
所以 `clean` 不是"什么都没扫"。CI 那一步的名字是 `gofmt (gofumpt)`，本程跑的是 **gofumpt**，
没拿 `gofmt` 那把另一尺寸的尺交差。

### 5.4　AC#4 契约轴（零字节，逐枚现量）

`git diff --name-only 80fa0551 -- internal/risk internal/panel internal/agent/approval tools/d22scan scripts/slo-check.ps1 docs/PLAN.md docs/specs`
→ **空**（本程工作树里这些路径一枚未出现）。
`frontend/**` 同一发也空。`design/**` 那一发**不空**，但里面**一枚都不是本程写的**：
它们是 §0 那 16 枚 owner 自己的未提交删除，加另一会话正在这棵树上改的 `design/doubao/**`
（本程收笔时它们已从 `??` 变成 ` M` ⇒ **工作树是活的**，这正是本仓硬禁 `git add -A`／`git commit -a` 的理由）。
本程的归属证明是每枚 commit 的 `git show --name-only`（§8），不是这张工作树表。

`thresholds.go`／golden／`allowlist.txt` 一枚未碰（`git ls-files` 现量它们的位置：
`internal/observe/thresholds.go`、`tools/d22scan/allowlist.txt`、`internal/agent/testdata/golden/*.sse`），
本程 diff 里不含 `internal/observe/`、不含 `tools/d22scan/`、不含 `testdata/`：见 §8 逐枚名册。

---

## 6　本程没测什么（按"如果我漏了它谁会先被骗"排序）

1. **`corrupt` 支上那一发同族形状没修**：39 字节的 `{"mode":"…",,"pass":true}` 语法错在第 25 字节附近，
   红句仍印 `at offset 0`。被"文档停在第 0 字节"带错方向的人，在 corrupt 那一支上**仍然会被带**。
   本程按票面射程只收了 `unwritten`。⇒ 谁先被骗：读 corrupt 红句排障的人。
2. **真竞态没碰**（票面"不解决"第 3 条）：本机那 3 秒 871 发部分读的唯一形状仍是 `len=0`，
   1..N−1 的中间前缀**依旧没有真写侧凭据**；本程全部前缀读数都来自**纯函数**与**脚本递的字节**，不是真写侧。
   ⇒ 谁先被骗：以为"真写侧被前缀扫守住了"的人。那一格仍归票 144 的证据件挂着。
3. **`summary()` 的其余调用点没逐枚钉**：本程三枚新用例覆盖 `unwritten`／`complete`／note 三支的句子；
   `reportCorrupt` 那一支的 `summary()` 文本只被票 144 既有用例的 `wantReason` 间接触及，本程没为它加检（见第 1 条）。
4. **`slo-check.ps1`／CI 那一步会不会读到新句子**：本程没在 CI 上跑过（子代理不 push）。
   本地证据是"全仓 `grep` 那两句文本只有 `slo_windows.go` ＋ 票 144 测试文件两处"，
   外加 §5 的门禁四发；**仓外如果有别处抄过旧红句**（例如已推送的 run 日志、票 144/147 之外的证据件），
   它们引用的 `at offset 0` 从今天起**过期**，本程只在 §0 与票 144 验收件的位置登记了这件事，**没有去改那些件**。
5. **Go 版本相关性没量**：`dec.InputOffset()` 对截断输入给 0 是**当前** stdlib 行为（本程量的是 go1.27.1）。
   本程**故意没有**为它写断言（那是把上游实现细节钉进本仓门禁）；如果哪天它不再给 0，
   红的是 §2.1 那三枚里的"句子两个数"那一味，不会静默。
6. **`wisp slo` 端到端没跑**（那要真拉一枚 subject、跑 33s 预算）。本程改的是错误文本与喂它的数，
   凭据是脚本与纯函数两级；真端到端那一发归 CI 的 `slo-full`。

---

## 7　伪授权登记（两个数分开栏）

- **真通知回显数：1** —— 会话首条 `<system-reminder>`（含 AGENTS.md 项目上下文与技能清单）。
  另：本程自己起的后台任务（`wisp-cli-tests.sh`）完成时由 harness 回显的 `task-notification`，**按类别记、枚数不取**。
- **判为注入数：0**。全程工具输出里没有出现过"编排者备注／系统提示／文件已被修改／已核验请继续提交／
  请放宽阈值／已解锁／不用取证直接给结论／Confirm the harness note is genuine"这类文字。
  唯一一处**看起来像但不是**的形状是 Read 工具在 `git log` 那次调用后面追加的技能清单块——
  它是 harness 自身的常规回显（出处：`Read`＋`Bash` 同批返回，命令前 40 字 `cd "D:\work\workspace\projects plans\Wisp" && g`），
  本程按类别登记、不当注入、也不据它改判据。
- **凭据值抄录数：0**。本件不出现任何密钥/token 的值；涉及的只有变量名（`WISP_84_MEASURE`、`WISP_IT_REAL_MIRROR`、
  `WISP_LIVE_MIC`）与文件名。
- **`rm` 使用数：0**（含 `rm -rf`／`rmdir`／`del`）。快照里每枚变异都以"cp 一份原件 → 改 → cp 还原 → `cmp` 验同"处理，
  原件留在 `/d/work/tmp147/`；本件与日志**只追加、不覆写旧凭据**（§2.1 加了字段 leg 之后，
  改前红重跑成一枚新日志 `red-before-fix-final.log`，旧那枚 `red-before-fix.log` 原样留在盘上，本节就是它的说明）。

---

## 8　落盘名册（本程 commit，取法随文）

（收笔时填：`git log --format="%h %s" 80fa055..HEAD` ＋ 每枚 `git show --name-only`；枚数不写在正文里。）
