# 票 145 —— 落地集＝空 的**非实现者对抗验收**（accept-r1）

> 派单＝`.scratch/wisp/dispatches/2026-09-26-103x-accept-145.md`。被验版本＝**`3dcff6a`**（票 145 实现件的最后一枚）。
> 本件**只追加不重写**；本程写面＝本文件 ＋ `.scratch/wisp/probes/145-accept/**`，其余全只读。
> 本程**不勾票面任何一格、不加 `-done`**；R1/R2 的修法归编排者，本程只裁"结构上可不可满足"。

---

## 0. 锚点 · 环境 · 取版尺本身先被量了一遍

### 0.1 锚点（进场现读）

```
$ git rev-parse --short HEAD          # 2026-09-26 10:4x 进场
0cdee19
$ git rev-parse --abbrev-ref HEAD
dev
$ git rev-parse 3dcff6a^{tree}        # 被验那枚树
9da188334df7bb0ec5d1a9ebfef2719e2a9eb28f
$ git cat-file -t 3dcff6a
commit
$ git log --format='%h %ad %s' --date=format:'%Y-%m-%d %H:%M' -1 3dcff6a
3dcff6a8 2026-09-26 10:31 evidence(145 §5.6-§5.7 追加): 合并态复算同一把尺＋一次真发生的共享 index 险情；节号按阅读次序重排
$ go version
go version go1.27.1 windows/amd64
```

⇒ **进场 HEAD＝`0cdee19`，与被验版本 `3dcff6a` 差两枚**（`2da833c` 票 153 验收第 3 格、`dd591fa` 台账更正、`0cdee19` 前端 log）。
本件下面每一枚读数都写在**`3dcff6a` 那枚树**上，不写工作树；工作树只用于 §5 那一格的"复跑门禁"（那本来就是工作树读数）。

### 0.2 **取版尺本身不干净**——本程造出来的一枚，先登记（它改变后面所有仓外读数的可信度）

派单指定 `git archive 3dcff6a | tar -x -C <仓外>` 取纯净树。本程**先量了这把尺本身**，它不合格：

```
$ ls .gitattributes            # 第 1 行：* text=auto
$ git archive 3dcff6a | tar -x -C /tmp/wisp-145-accept/tree
$ python -c "...compare tree(archive) vs 3dcff6a blob..."
files compared 1514  differing 423
  .html 15 · .css 6 · .ts 6 · .tsx 48 · .json 34 · (.scratch) 243 · .go 0
$ git cat-file -s 3dcff6a:frontend/fixtures/composer-states.html   → 16604（blob 内 CRLF=0）
$ wc -c /tmp/wisp-145-accept/tree/frontend/fixtures/composer-states.html → 16649（CRLF=45）
```

⇒ `* text=auto` 在 Windows 上让 `git archive` 把 LF blob 落成 CRLF 树。**`.go` 一枚不受影响（`*.go text eol=lf`），受影响的是资产面**——
而本票那一格恰好读了它：archive 树里 `TestComposerRenderFixtureTellsTheTruth` **红**，红因是该用例按字面
`" -->\n"` 切块（`composer_test.go:686` 的 `strings.Cut`），CRLF 树里切不出块 ⇒ "0 painted states"。
**那是本程的取版尺造的假红，不是被验物的红。** 实测两发在 `probes/145-accept/00-extract-ruler-{archive,exact}.txt`。

⇒ 本程换成**逐 blob 落盘**（`git ls-tree -r -z` ＋ `git cat-file --batch`，1514 枚全量），得**字节精确树**
`C:\Users\swq\AppData\Local\Temp\wisp-145-accept\exact`（在仓外，**只建不删**）。同一把尺在两种树上的读数：

| 树 | `go test -count=1 -v ./internal/panel/` | rc | RUN | PASS | FAIL | 红的哪一枚 |
|---|---|---|---|---|---|---|
| `git archive`（本程第一发，作废） | 同上 | 1 | 105 | 57 | **2** | C21 ＋ `TestComposerRenderFixtureTellsTheTruth`（取版尺假红） |
| **逐 blob 精确树**（下面各格用这支） | 同上 | 1 | 105 | 58 | **1** | 只有 `TestC21DesignTokensFourWayAgree` |

⇒ **第二发与实现程 §5.1 的"FAIL=1"同数**，且**本程在精确树里没造出第二枚红**（派单前提⑤要的那一问，见 §5）。
⚠ 一句给编排者的：这条不是实现程的毛病——**它交的是共享工作树读数**，工作树里 `*.go`/`*.html` 都是 `git status` 干净的；
**它被写坏的只可能是"谁拿 archive 树去复算它的读数"这一支**。⇒ 复算本票的仓外读数**必须**用 blob 精确树或真 checkout，不能用 archive。

### 0.3 本程写面 · 地界现量（每次 commit 前重取一次，读数在下面各节末）

```
$ git status --porcelain | head           # 10:4x，进场第一发
 M .scratch/wisp/probes/152/my152.py
 D design/assets/base.css … D design/screens/tasks.html   （design/** 16 枚被另一枚会话删着）
 M design/doubao/README.md · M design/doubao/demo/*
?? .scratch/wisp/dispatches/2026-09-26-103x-accept-145.md
?? .scratch/wisp/probes/139/accept-r1/ · 152/accept-r1/ · 152/mut-shipped/
?? .zcodeignore · design/doubao/01-ball-states.jpg · design/doubao/demo/lib/ …
```

⇒ **`design/**` 与 `frontend/**` 在工作树里是脏的（别家在写）**，本件凡"零命中"一律限定在 **`3dcff6a` 的 blob 精确树**，
且**不含** `frontend/**`、`design/**`、`frontend/src/lib/panel.ts` 的工作树状态（派单硬约束）。
`docs/reports/**`／`PLAN.md`／`specs/**`／`internal/**`／`cmd/**`／`thresholds.go`／golden／`allowlist.txt`
／`tools/d22scan/**` —— 本程**零写**（下面每一枚 commit 的 `git diff --cached --name-only` 名册逐枚贴出）。

---

## 1. 前提① E1——"加那 9 枚键、`panel.ts` 一字不动 ⇒ 响的是四枚用例"：**成立，且本程把它钉成"与枚数无关"**

**尺**：blob 精确树（§0.2 那一支）＋**只**加 `composer.go` 的九枚键（`pump.go`／`run.go`／`panel.ts` 全不动）。
读数＝`go test -count=1 -v ./internal/panel/` 的 `--- FAIL` 名册，不是推理。原始件：`probes/145-accept/01-E1-four-pins.txt`、`02-E1b-one-key-same-four-pins.txt`。

```
$ cp -r exact tree-mut1 && cd tree-mut1
$ git apply -p1 --include=internal/panel/composer.go <repo>/.scratch/wisp/probes/145/patch/145-snapshot-fields-closure.patch
$ go test -count=1 -v ./internal/panel/     # RUN=105 PASS=54 FAIL=5
--- FAIL: TestApprovalCardViewJSONKeysMatchFrontendTypes        approval_test.go:129: Go Snapshot emits [run tools approval cost failures] that interface PanelSnapshot does not declare
--- FAIL: TestComposerContractTypesMatchFrontend               composer_test.go:74:  Go Snapshot emits [run tools approval cost failures] that interface PanelSnapshot does not declare
--- FAIL: TestThePumpBuildsThePacketFromWhatTheHostHolds       pump_test.go:124:     snapshot JSON keys = [...run cost ... tools approval failures], want exactly the four PanelSnapshot declares
--- FAIL: TestPublishHandsTheBytesToTheAttachedExit            pump_test.go:276:     exit bytes carry keys [...], want the four PanelSnapshot declares
--- FAIL: TestC21DesignTokensFourWayAgree                                              （基线既有红，与键无关，§5 单判）
```

⇒ **四枚，逐名逐行号与派单给的那四个锚点一字不差**（`approval_test.go:129`／`composer_test.go:74`／`pump_test.go:124`／`:276`
——行号本程在 3dcff6a 现读到 `t.Errorf` 就在那四行上，不是从日志抄的）。
实现程把普查 §4.4 的"一枚"更正为"四枚"——**本程独立复算，更正成立**。

**本程自己加的一发（实现程没有）**：只加**一枚**键（`Cost CostProbe \`json:"cost"\``，两字段）⇒

```
$ cp -r exact tree-mut2 && python <插一枚键> && go test -count=1 -v ./internal/panel/   # RUN=105 PASS=54 FAIL=5
    approval_test.go:129: Go Snapshot emits [cost] that interface PanelSnapshot does not declare
    composer_test.go:74:  Go Snapshot emits [cost] that interface PanelSnapshot does not declare
    pump_test.go:124:     snapshot JSON keys = [cost pending results composer generatedAt], want exactly the four ...
    pump_test.go:276:     exit bytes carry keys [pending results composer generatedAt cost], want the four ...
```

⇒ **同一批四枚钉，一枚键就全响**。这条把实现程注释里那句"加一枚键要付四枚钉"**从声称升成读数**，
也顺手否掉一种误读：**"少加几枚键就少响几枚"不成立**——爆炸半径与枚数无关，与"有没有第五枚键"有关。

⚠ **一条本程不替它圆的**：`composer.go` 的头注释（`eb38c97` 落的）写的是"两把双向对账尺在 `composer_test.go:48` 与
`approval_test.go:105`"——那两行是**函数注释/用例起点**（现量 `approval_test.go:105` ＝ `func TestApprovalCardViewJSONKeysMatchFrontendTypes`、
`composer_test.go:48` ＝ `func TestComposerContractTypesMatchFrontend`），而**响的行**是 `:129`／`:74`。
两读不矛盾（一枚指"尺在哪"，一枚指"哪一句红"），但**注释里那两行号是"定义处"不是"断言处"**，下一位照注释去 `:48` 找断言会扑空。属**措辞精度**，不属读数错。

---

## 2. 前提② E1b（派单称"最硬的一发"）：**成立——但本程把它从"手搓 PumpSources"换成真装配根，并量出一处它没写的边界**

### 2.1 先复量派单给的那几枚锚点（行号是"它的锚点"下的，本程现读）

```
$ grep -n "NewSnapshotPump(panel.PumpSources{" cmd/wisp/run.go        → 421（字面量 421-427，读口 Verdicts/Mode/Workspace/Results/Out）
$ grep -n "func (c consoleSink) Publish" cmd/wisp/run.go              → 751（事件汇 751-785 确为那一整枚 switch）
$ grep -rn --include=*.go "NewSnapshotPump(" .  | grep -v _test.go
      ./cmd/wisp/run.go:421                       ← 唯一生产调用者
      ./internal/panel/pump.go:145                ← 定义本身
      ./.scratch/wisp/probes/151/run.head.go:421  ← ⚠ 票 151 的**探针副本**，不是码
      ./.scratch/wisp/probes/151/run.mut-noclose.go:421
$ git diff --numstat e24ae28 3dcff6a -- cmd/wisp/run.go                → 空（实现程下的行号在本被验版本仍逐枚有效）
```

⇒ 锚点四枚全部成立，且要补一句口径：**"全树唯一生产调用者"只有在把 `.scratch/wisp/probes/**` 剔掉时才成立**
（票 151 在探针目录里存了两枚 `run.go` 的副本，同一行号）。实现程那句"唯一一枚＝此处"结论对、**尺不够精确**。

### 2.2 复现：本程**不**手搓 `PumpSources`，改跑真装配根

派单/实现程那一发（§2.4）用的是自造的五枚读口字面量（实现程自己在 §7.3#2 承认了）。本程换成
`newRunFixture → runTextTask（= `wisp run` 那一支组合根）→ rt.publishPanelSnapshot() → rt.lastPanelSnapshot()` 的
**出口字节**（不是 Go struct）。测试源＝`probes/145-accept/04b-E1b-probe-test.go.txt`，读数＝`04-E1b-real-assembly-root.txt`。

```
$ cp -r exact tree-mut1 && git apply --include=internal/panel/composer.go <closure.patch>   # 只加九枚键
$ PATH="<repo>/third_party/sherpa-onnx:$PATH" go test -count=1 -v -run TestAccept145WhatTheRealAssemblyRoot ./cmd/wisp/
    wire["results"]     = [{"correlationId":"18207b6c-…","text":"echo: ## scene\n当前时间：…","done":true}]   ← 活
    wire["composer"]    = {"mode":{"current":"ask_every_step",…},"workspace":{…,"reason":"未选择工作区：…"},…}  ← 活
    wire["generatedAt"] = "2026-09-26T02:46:39Z"
    wire["run"]         = {"reasoning":"","status":"","usage":{"input":0,"output":0,"cachedRead":0}}
    wire["tools"]       = null
    wire["approval"]    = {"depth":0,"windowMs":0,"vetoChannels":null}
    wire["cost"]        = {"tokensIn":0,"tokensOut":0,"cachedTokens":0,"micros":0,"currency":""}
    wire["failures"]    = null
    LIVE in the same packet: pending=0 results=1 mode=ask_every_step workspace.set=false
--- PASS   rc=0   （断言方向＝断言"没填"，见放水两问）
```

⇒ **五枚新段全部以零值／null 出门，而同一个包里 `results`／`composer` 带真值**（`pending=0` 见 2.3）。
⇒ 派单前提②那句"`run.go:421-427` 不动 ⇒ 新增 section 跑出进程的是那串常量"**成立**，
且它现在是**真装配根读数**，不再是"同形状手搓"。`{"depth":0,"windowMs":0,"vetoChannels":null}` 与派单引的一字不差。

### 2.3 两处本程要给它划清楚的边界（不是推翻，是别让下一位读过头）

- **"pending 也是活的"这一支本程没复现**，因为本程这一发跑的是普通文本任务、队列里没有卡。
  实现程 §2.4 那个 `pending len=1` 来自**它自己喂给手搓读口的**一条 `NativeVerdict`——**那不是"生产在填"，是"探针在填"**。
  ⇒ 结论不受影响（对比项换成 `results`/`composer` 更硬：那两枚是这条真任务真填的），
  但**"同一个包里 pending/results 是活的"这句在实现件里是手搓读数**，本程把它降级为"半读数半构造"。
- **"提前 ship 给面板"这一支不成立，而实现程自己没这么写**：这串常量今天**到不了页**——
  泵出口是 `rt.bookPanelSnapshot`（写台账摘要＋在内存里留一份字节），树里没有 Go→页通道（`panel_pump.go:13-21` 自证、实现件 N1 同判）。
  ⇒ 所以 E1b 证明的是**"有键无填＝契约里进五枚常量"**（形状＝普查 §2.0(2) 的第四态"无生产者"的镜像），
  **不是**"用户会看到一堆 0"。本票的判词应写成前者；后者目前不会发生。实现件 §2.4 的原句
  "加出来的正是本票自己要防的那形"——**那形＝装饰／契约污染**，它没写"用户看得见"，所以**这句不算错**，
  但派单正文第 15 行把它转述成"把本票要防的常量**提前 ship 出去**"，**"ship 出去"三字过头**，本程按 2.3 校正。

---

## 3. 前提③：四把锁的归属 ＋ **AC#6 到底是不是结构上不可满足**（本程独立重走，含一发自造的钥匙探测）

### 3.1 四把锁逐枚复量（出处＝`.scratch/wisp/dispatches/2026-09-26-093x-thaw-panel-for-145.md` 现读，不抄实现件）

| 锁 | 派单说它是谁 | 本程现量 | 判 |
|---|---|---|---|
| **B** | `frontend/src/lib/panel.ts`，Q-51 未答、实现程禁写 | 特批正文里 `frontend/**` 那句是**"继续完全隔离：不读它的未提交改动、不写"**；两把双向对账尺读的正是这枚文件（`approval_test.go:107`、`composer_test.go:50` 现量 `os.ReadFile(filepath.Join(root,"frontend","src","lib","panel.ts"))`）——**尺的输入面＝禁写面**，这是结构不是态度 | ✅ 成立 |
| **C** | 三枚被钉死的测试文件 | 响的四枚里两枚在 `pump_test.go`（§1），另一枚 `approval_test.go`／`composer_test.go`——**特批的放开清单只枚举了三枚文件**（`composer.go`／`pump.go`／`cmd/wisp/panel_pump.go`），其余按它自己那句"越界即作废"处理；另有三枚被**具名保留**（`tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`frontend_hygiene_test.go`） | ✅ 成立 |
| **D** | `cmd/wisp/run.go:421-427` ＋ sink `:751-785`＝票 151 的脸 | 两处行号在 3dcff6a **逐枚现量有效**（§2.1）；特批正文**自己写着**"本程不许碰 `run.go`"，台账 `A273③` 同句"**不碰 151 的 `run.go`**" | ✅ 成立（比实现件引的更强：不止台账划的，派单正文也写了） |
| **E** | 放开面里没有 `*_test.go` ⇒ AC#6 无处落 | 见 3.2——本程**没有只信文档**，另造了一发探测"能不能在放开的非测试文件里塞一枚会跑的用例" | ✅ 成立，且是**硬**成立 |

### 3.2 AC#6：三发独立凭据（原始件 `probes/145-accept/06-lockE-nontest-file-probe.txt`）

**第一发（本程自造，是这一格真正的钥匙）**：Go 允不允许"用例写在非 `_test.go` 文件里"？
如果允许，锁 E 就是措辞问题不是结构问题，本程就得推翻它。实测——往**放开面内的** `internal/panel/composer.go`
末尾加一枚 `func TestProbe145AcceptInNonTestFile(t *testing.T) { t.Fatal(...) }` 并 `import "testing"`：

```
$ go test -count=1 -v -run TestProbe145AcceptInNonTestFile ./internal/panel/
    testing: warning: no tests to run
    ok  github.com/CarlosShao/wisp/internal/panel  [no tests to run]
$ go vet ./internal/panel/                     → rc=0，**一声不响**
$ go test -count=1 -v ./internal/panel/        → RUN=105（与基线同枚数），名册里该探针名出现 **0 次**
```

⇒ **非 `_test.go` 里的 `TestXxx` 编得进、vet 得绿、永远不跑**。⇒ 在那三枚放开文件里**造不出任何一枚会被执行的断言**，
没有旁路。锁 E 成立，且**不是"它没找到写法"**。
⚠ 顺手登记一枚**本程新增的仪器缺口**（与实现程那枚 json tag 缺口同族、方向更坏）：
**"看着像用例、其实从不执行"这一形今天没有任何尺抓**——它比"忘了加 json tag"更危险，因为它给的是**绿＋一枚不跑的断言**，
而 AC#6 的答句要的恰恰是"哪一枚用例"。归口＝编排者（新增判据，本程不擅自加钉）。

**第二发（承重侧的对照，证明缺口是"地界"的不是"世界"的）**：把实现件那份闭合补丁**完整** apply 到 3dcff6a 的 blob 精确树
（8 枚文件全 clean，`composer.go` 偏移 13、`pump.go` 偏移 14 ⇒ **它只声明了 `e24ae28` 干净，本程量到 3dcff6a 也干净**）：

```
./internal/panel/  RUN=106 PASS=59 FAIL=1   ← 唯一红＝基线那枚 C21；两把尺 "9 JSON keys reconciled" 转绿
./cmd/wisp/ -run TestSnapshotGateSectionComesFromTheLiveGate   rc=0 --- PASS
    packet approval section = {"depth":0,"windowMs":2000,"vetoChannels":[ …四行含 B1 原句"说取消词：语音取消不可用"… ]}
```
再按实现件的 E4 摘掉**唯一**那行填值（`pump.go` 的 `snap.Approval = approval`）：

```
    windowMs = 0, want the gate's own 2000
    windowMs = 0: the section arrived unfilled
    packet carries 0 channel rows, registry reports 4
--- FAIL   rc=1
```
⇒ 实现件 §2.5/§2.6 那两发**本程逐字复现**（含那三句红句）。⇒ **AC#6 的缺口在于放开面，不在于可达性**：
拿到 `panel.ts`＋`pump_test.go`＋`run.go` 一行读口＋一枚 `*_test.go` 这四把钥匙，AC#6 当场可答且承重。

**第三发（实现件没写的一层，直接影响验收表怎么措辞）**：**空落地集下 AC#6 是"空真"，不是"不可满足"**。
AC#6 的句式是"对**每一枚新字段**答一句"——落地集＝空 ⇒ 全称命题对空域自动为真。
⇒ 所以准确判词是**两条**，不是一条：
1. **AC#2 未落地**（判据成立，见 §1/§2），因而 **AC#6 今天没有回答对象**；
2. **一旦 AC#2 落地任何一枚字段，AC#6 在这枚放开面内立即不可满足**（3.2 第一发）。
⇒ 派单第 3 条那句禁令本程**采纳并给出替代措辞**：验收表里不许写"AC#6 待补"；
应写成 **"AC#6：本放开面内无定义域；其可满足性与 AC#2 的落地集非空互为条件，解锁需 3.2 第二发那四把钥匙"**——
既不是待办，也不是通过。

### 3.3 一句总结这一格

四把锁**枚枚成立**，其中锁 E 本程**加硬成"物理不可绕"**；而"能绿且承重"本程**自己造出来了**（第二发），
所以 R1/R2 不是"做不到"的问题、是**给不给钥匙**的问题——**修法二选一归编排者，本程不选**。

---

## 4. 前提④：`3dcff6a` 的那 10 处删除——**只认 diff**：判为"重排不是重写"，但**顺手把整枚文件的删除史全查了，另捉到 2 处**

### 4.1 派单点名的那一发（34 增／10 删）

```
$ git show --numstat --format= 3dcff6a -- docs/evidence/s1/145-snapshot-fields-landed-r1.md
34      10      docs/evidence/s1/145-snapshot-fields-landed-r1.md          ← 与派单给的数一字不差
```

逐行配对（把 10 枚 `-` 行的节号挖成 `#` 后去 34 枚 `+` 行里找孪生行）：**10/10 全部配对，unpaired＝0**。
⇒ 每一枚删除都有一枚"只差节号"的新行在原位等着，**没有任何旧结论被换成新结论**。
其中 8 枚是标题行（`## 8. 收口`→`## 7. 收口`、`### 8.1-8.5`→`### 7.1-7.5`、`## 7. 报回`→`## 8. 报回`、`### 7.1`→`### 8.1`），
**1 枚是表格行**（AC#2 那行末列"副本复算见 §8.2"→"§7.2"），**1 枚是正文行**（`:239` "…互斥，见 §7 报回"→"见 §8 报回"）。

⇒ **判**：`A278②`（证据件只能追加）这一格 **3dcff6a 未违反**。
⚠ 但它的 **commit message 有一处过头**：`（本节号 8→7、7→8 是纯重排：…删的 10 行全是标题行，内容零丢）`——
现量是 **8 枚标题行＋1 枚表格行＋1 枚正文行**。"内容零丢"成立（只差节号），**"全是标题行"不成立**。
⇒ 这正是派单第 4 条说的"不要跟着 commit message 走"的一枚实例：**结论没受影响，措辞夸了一格**。

**重排有没有留下断头引用**（本程补的一发，实现件没自证）：

```
$ grep -n "见 §7\|见 §8\|§8\.2\|§7\.2\|§5\.5" <3dcff6a 版>     → 只剩 §5.5/§7.2/§8 三形，全部指向存在的小节
$ grep -o '^#{2,3} [0-9.]*' <3dcff6a 版> | sort | uniq -d      → 空（**无重号**）
```
⇒ 重排自洽。唯一的既有编号怪相是 `### 3.4` 上面没有 3.1-3.3（**早于 3dcff6a**，`c655458` 就长这样，不算这一发的账）。

### 4.2 本程自己扩的口径：**全文删除史＝12 行，派单只点了 10 行**

```
$ for h in ef07ba5 c655458 972ceba 7263456 4156ff5 3dcff6a; do git show --numstat --format= $h -- <证据件>; done
ef07ba5  92   0     c655458 147  0     972ceba 43  0
7263456 157   1  ←  4156ff5  63   1  ←  3dcff6a  34  10
```

那另两枚"改一句"的本程逐行读（不是数数）：

| commit | 删的那一行（原文要旨） | 换成 | 方向 |
|---|---|---|---|
| `7263456` | §0.2"生产码改动"行：只写了 `composer.go` 一枚文件加了注释 | 加上 `pump.go`、"四把锁"→"四枚钉＋四把锁＋无生产者那枚字段"、并把"见 §7 那枚 commit"换成具名 `eb38c97` | **自报面变宽**（多认一枚文件），不是变窄 |
| `4156ff5` | N10 行：`（loop.go:937 那三枚字段）` | `（loop.go 那枚 publish 只带 Kind/TaskID/Text/TokensIn/TokensOut）` | **由"三枚"改成"五枚"**——本程现量 `internal/agent/loop.go:938-941` 确为那 5 枚键，**旧句数少了**，新句对 |

⇒ 判词分两层，别混：
1. **按字面**，`A278②` 是"只能追加、不能重写"，这两行是**重写**（旧文没留在原地）；本仓ledger 自己的更正形是
   **"原句一字不抹 ＋ 另起 `>` 更正块"**（现例＝`A274①` 里编排者给自己那三条限定）。⇒ 严格论，**这两处形不合**。
2. **按方向**，两处都是**把自我报告改准／改宽**（一处多认一枚文件、一处把字段数从 3 修成 5），
   **没有一处把旧结论改成更松的结论**，也没有一处删掉读数。⇒ **不构成"把不该结案的顺手结案"**。

⇒ **给编排者的处置建议（不是裁定）**：这 2 行**不需要**追责式更正；但**下一版派单/模板**该把
"证据件里的自我报告写错时怎么改"钉成一个形：**追加一节 `> 更正` 并把旧行留着**，而不是就地换字——
否则"只认 diff"这条纪律要求每枚验收程都去背 12 行删除史（本程背了，下一位未必）。

### 4.3 这一格总结

**派单前提④ 的读数（34/10）成立，性质判为"重排非重写"；commit message 的"全是标题行"被本程量成"8 标题＋2 正文/表格"**。
外加两枚派单没点名的单行重写，**方向均为收紧**，判为形不合／实不伤，登记不追究。

### 4.4 追加（本程写完 4.2 又现量了一枚时间轴，它改变"该不该追责"的判断）

```
$ grep -n "^## A278" docs/reports/pending-and-issues.md   → 6858 行，标题时刻＝09-26 09:5x
   （A278② 就是"证据件只能追加、不能重写"这条规矩的出处，起因＝票 153 那程自报改写过自己已提交的证据件）
$ git log --format='%h %ad' --date=format:'%H:%M' -1 093x-thaw-panel-for-145.md 的派单时刻 ＝ 09:3x
$ 145 的两枚单行重写：7263456 ＝ 10:23，4156ff5 ＝ 10:25
```

⇒ 时序是：**规矩在台账里立起来（09:5x）晚于 145 的派单（09:3x）**，而 **145 那两枚重写发生在规矩之后**。
⇒ 所以判词要分两句：**它没有从自己的派单里读到这条规矩**（那份特批里只有一句"撤销后…按追加更正处理，不改写历史"，
管的是撤销不是证据件），但**它此刻盘上已有 A278 这条**。⇒ 本程维持 4.2 的判（形不合、方向为收紧、不追究），
只把"下一版派单要写死这条"从建议升为**具名待办**：**`093x` 那份特批的"共同纪律"段缺 `A278②` 一句**。
