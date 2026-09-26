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

---

## 5. 前提⑤ 门禁自报复跑：四数**逐枚复现**、名册**两向逐名相同**；但那枚红的**红因**在两棵树里不是同一枚（本程量出来了）

原始件＝`probes/145-accept/07-gate-rerun.txt`。全部在 §0.2 那枚 **blob 精确树**（被验版本 3dcff6a）里跑，
只有 A/B/C 三发的 `cmd/wisp` 需要把 `third_party/sherpa-onnx/*.dll` 拷进树（未 tracked，见 5.2）。

### 5.1 四数复算 ＋ 名册两向 `comm`

| 包 | 版本 | 本程复跑 rc / `=== RUN` / PASS / FAIL / SKIP | 实现程自报 | 判 |
|---|---|---|---|---|
| `./internal/panel/` | 改前 `e24ae28`（本程另取一枚精确树跑） | 1 / **105** / 58 / 1 / 0 | 1 / 105 / 58 / 1 | ✅ 逐数相同 |
| `./internal/panel/` | **被验版本 `3dcff6a`** | 1 / **105** / 58 / 1 / 0 | 改后 1 / 105 / 58 / 1 | ✅ |
| `./cmd/wisp/` | `3dcff6a`（DLL 进树后） | **0 / 139 / 79 / 0 / 0**（139.5s） | 0 / 139 / 79 / 0 | ✅ |

**名册不是"复跑一次比一次"，本程做了三层**：
① 实现程**自己那两枚已提交的**改前/改后名册文件，本程直接两向 `comm`（不用跑就验了它那句"差集为空"）：
`roster2-{panel,cmdwisp}-before.txt` vs `roster-after-{panel,cmdwisp}.txt` ⇒ **四个方向全空**；
② 本程 `e24ae28` 名册（59 行）vs 本程 `3dcff6a` 名册（59 行）⇒ **两向皆空**；
③ 本程 `3dcff6a` 名册 vs 实现程 `roster-after-*` ⇒ panel **59/59 逐名相同**、`cmd/wisp` **79/79 逐名相同**。
⚠ 一处尺事（本程自己的）：第一发比对 79/79 **全差**，因为本程的名册带着 ` (12.34s)` 时长后缀而实现程的名册剥掉了 ⇒
`sed -E 's/ \([0-9.]+s\)$//'` 归一后才可比。**不是它的读数问题，是"名册"这个词在两份件里口径不同**，登记以免下一位拿两份件直接 `comm` 得出"全差"。

### 5.2 那一枚红的**红因**：派单要本程判"是不是它说的那一个"——**在它跑的树里是，在被验版本里不是**

实现程看到的（它 §5.2，工作树那一发）：

```
tokens_fourway_test.go:441: read design/assets/tokens.css: open …\design\assets\tokens.css: The system cannot find the path specified.
⇒ readRepoFile 是 t.Fatalf 形 ⇒ 第一枚文件读不到就终止，四方比对根本没跑
```

本程现量：`git status --porcelain` 里 `design/assets/tokens.css` **是另一枚会话的未提交删除**（工作树 ` D `，HEAD 里在）。
⇒ 那句红因**对它跑的那棵树成立**。

但在 **3dcff6a 的 blob 精确树里四枚文件全在**，那一枚用例**真的把四方比对跑完了**，红在内容：

```
tokens_fourway_test.go:446: parties: design/assets/tokens.css=135 dark decls, c21-native-tokens.md=78 colour rows,
                                 internal/ball/tokens.go=40/38 palette fields, tokens.generated.css=48 dark decls
tokens_fourway_test.go:461: frontend/src/styles/tokens.generated.css does not carry --panel-w = 640px that design/assets/tokens.css declares
                            … （同形逐枚）
tokens_fourway_test.go:532: c21-native-tokens.md:118-132: design/assets/tokens.css declares --orb-body-2/--tint-* 但 generated.css 不携带
tokens_fourway_test.go:547: only 0 of 78 colour rows completed all four legs - the check is not doing its job
```

⇒ **两问两答**：
- **"今天这枚红的红因是不是它说的那一个"**——在被验版本上**不是**：3dcff6a 的红因是**内容不齐**（0/78 行走完四腿），
  它写的那枚"文件缺失⇒没跑"只在**脏工作树**成立。它的措辞是老实的（"今天这一发日志里它只暴露一枚红因"），**不算错**；
  但它由此得出的"**不能据此说另一条不存在**"这句保留疑问，**本程把它答成了肯定的：另一条存在，而且是被验版本上唯一那条**。
  ⇒ 记忆里那条"现有两枚红因（**放回文件也不会绿**）"**本程实测成立**（文件放回＝3dcff6a 精确树，仍红）。
- **"有没有第二枚红被当噪声吞了"**——**没有**：本程在 3dcff6a 精确树里跑出的 `FAIL` 枚数＝1，逐名＝`TestC21DesignTokensFourWayAgree`，
  与它自报同枚同数。**但本程自己造出来一枚假红**（§0.2：`git archive` 的 CRLF 让 `TestComposerRenderFixtureTellsTheTruth` 红），
  以及**一枚"仪器没跑到"形**：`cmd/wisp` 第一发 rc=1／FAIL=8／RUN=133，八枚红因全是同一句
  `no native DLLs in ..\..\third_party\sherpa-onnx`——真因＝**blob 精确树里没有未 tracked 的 `third_party/`**，
  与本票无关，那一发**作废、不计入任何读数**（把三枚 DLL 拷进树后重跑＝上表那发 139/79/0）。
  ⚠ 派单警戒的那条形态（"包级 rc=1 而零条 `--- FAIL`"）本程没撞上，撞上的是**近亲形**：rc=1、八条 `--- FAIL`、
  红句里写着"请先跑 fetch-deps" ——**只看枚数会把环境缺件读成被验物红**，这一条本程写进 §5.3。

### 5.3 三件仪器复算

| 仪器 | 本程在 3dcff6a 精确树 | 实现程自报（工作树） | 判 |
|---|---|---|---|
| `gofmt -l internal/panel/ cmd/wisp/` | 空 | 空 | ✅ |
| `go vet ./internal/panel/ ./cmd/wisp/` | rc=0、无输出 | 同 | ✅ |
| `sh scripts/d22scan.sh` | **rc=0 clean**；分母 internal=205／cmd=23／ban#6 frontend=73／ban#7 internal/tools=18／**ban#8 design=30**／frontend=73／internal=413／cmd=44 | rc=0；同数，唯 **design=39** | ✅ 差的那 9 枚现量得解：工作树 `design/**` 文本件＝39、`3dcff6a` 提交版＝30 ⇒ **它读的是工作树**（含别家未提交的 design），**不是它的读数错**；本程按派单不把 design/** 算进任何宣称 |

⚠ 本程这一发 d22scan 顶部有一行实现程没有的：`d22scan: gitignore rules NOT APPLIED - git cannot be consulted in …exact`
⇒ 精确树里没有 `.git`，仪器**响亮地**拒绝按 gitignore 跳件（它自己那句"a scanner that cannot ask git … does not get to skip any"）。
⇒ 复算这一格时**两种分母都合法**，实现程那句 `skipped as git-ignored: 1 file(s) under frontend/dist/assets/` 只有真仓里才会出现。

⇒ **这一格判**：实现程 §5 的门禁自报**逐枚复现**（含它自己承认作废的那发 `rc=127` 之外的一切），
名册差集**三层皆空**；**唯一要改的是红因的适用范围**——它那句"四路比较根本没跑"必须带"在脏工作树里"这个限定，
**在 `3dcff6a` 上四路比较是跑完了的、红在 0/78**。owner 已明示这一面是已知常红 ⇒ **本程没为它变绿做任何事**
（没动 `tokens_fourway_test.go`、没动 `design/**`、没换基准）。

---

## 6. 前提⑥ 它报回来的三件新发现——逐件判真伪（每件都现量，尺写在前面）

### 6.1 ① "`A273②` 的'不碰 151 的 `run.go`'对泵成立、对 AC#2 的填充不成立（4/5 候选段需要 `run.go` 一枚读数行）"

**判：结论成立，两处措辞要改。**

- **引用位置错了一格（具名可复核）**：那句原文在 **`A273③`**（台账 6768 行：
  "⇒ **互斥按文件级判、四组两两不重叠**（我派前现量了构造点…，**不碰 151 的 `run.go`**）"），
  `A273②` 讲的是 owner 那三句话。**本程的派单正文也抄成了 `A273②`** ⇒ 这条是**两边一起错一个标号**，
  内容没逃（句子在、约束在）。⇒ 更正＝"引 `A273③`"。**且这枚约束不靠台账**：`093x` 特批正文自己写着"本程不许碰 `run.go`"（§3.1 锁 D）。
- **"4/5"应为"5/5（按已证成的形状）"**：本程查了它那枚"只有 `approval.*` 可旁通"的支路——
  它说"本程在仓外副本真做通了"并指 §2.4，**但 §2.4 是 E1b 那发手搓 `PumpSources`、量到的是 `approval` 出门为 `{"depth":0,"windowMs":0,"vetoChannels":null}`**，
  **不是把 approval 旁通填起来的读数**；它唯一真填出值的 E3 用的正是 `run.go` 里那行 `Approval: rt.liveGateState,`。
  本程自己查了旁通到底可不可行：`PumpSources` 全字段＝Verdicts/Mode/Workspace/Results/AttachmentMax/Now/Out（现量），
  `NativeVerdict` 全字段＝CorrelationID/Tool/Args/CallChain/Level/RulesHit/Reason/SessionOverrideBlocked（`pump.go:64-78`），
  **没有一处容得下 gate 窗口或通道名册**；`p.src` 是 `panel` 包私有字段，`cmd/wisp` 也改不动。
  ⇒ 旁通只剩一枚形状：**给 `NativeVerdict` 加一枚"其实是 gate 状态"的字段**（两枚放开文件内可做），
  代价＝把门的状态伪装成裁决携带。**能想，未做，未证**。⇒ 判该句**"未证"**，落地集＝空**不依赖它**（B/C 两把锁对任何键都响，§1）。
- **"对泵成立"这一半本程复量成立**：泵（`pump.go`／`panel_pump.go`）今天确实在跑、确实不需要 `run.go` 动一行——
  台账同一行给的三个构造点（`composer.go:79`／`pump.go:203`／`panel_pump.go:225`）**在 3dcff6a 已全部漂移**
  （现量 `type Snapshot struct` 在 **`:57`**、`NewSnapshot` 在 `:76`、`return Snapshot{` 在 **`:92`**；
  漂移原因＝`eb38c97` 自己在那三行上面加了 13 行注释）。
  ⇒ **一条给下一位的实操**：`composer.go:44-49` 这组行号（工单票面、普查件、实现件 §0.3 尺表**三处都在用**）
  **在被验版本上已不指向 `Snapshot`**；引它之前先 `grep -n "type Snapshot struct"`。实现件 §1 P1 那句"✅ 成立"
  是**在它自己加注释之前的那枚树上量的**——不是造假，但是**版本没带**。

### 6.2 ② "`remainingMs` 应从甲组挪乙组：`approval.EventTick` 声明与发射都是零"

**判：成立（本程独立复量，并把"零"用的尺写清）。**

```
尺＝grep -rnw --include=*.go EventTick <全树>  再剔 ^./\.scratch   （大小写敏感、整词、含测试文件）
  ./internal/agent/approval/ui.go:64:  EventTick EventKind = "tick"   ← 声明，一枚
  ./internal/panel/pump.go:41:         // … approval.EventTick is declared and   ← 注释，不是码
⇒ 生产/测试**发射者 0 枚**（不是"不存在这个符号"，是"没有一处构造它"）
Remaining 消费者全数（尺＝grep -rnw --include=*.go Remaining internal cmd，剔 _test）：
  gate.go:274 / :387  Kind: EventWarning, Remaining: g.window          ← 静态窗口长
  gate.go:512         Remaining: g.q.WarningLead()                     ← 静态提前量
  ui.go:74            Remaining time.Duration                          ← 字段本身
倒计时本体＝gate.go:257 `deadline := g.clock.After(g.window)`；approval 包内 `time.Until(`/`Until(` **零命中** ⇒ 无"还剩多少"的读口
```
⇒ **三处发射带的都是静态长度、没有一枚是"还剩多久"** ⇒ 它的判"活倒计时今天没人算"成立；
⇒ `remainingMs` 若落＝**填一个永远等于 `windowMs` 的数**，正是 AC#1 判据①／本票"宁缺毋造"要拦的那一形。⚠ 别用墙钟差补（`AGENTS §1.2` 禁形）。

### 6.3 ③ 仪器缺口："未打 `json:` tag 的导出字段能同时躲过两把契约尺，而 `encoding/json` 照发"

**判：成立（本程自造一发复现，且它的兜底只有一枚）。** 原始件＝`probes/145-accept/03-untagged-field-ruler-gap.txt`。

```
$ cp -r exact tree-mut3 && 给 Snapshot 加一枚 `LeakProbe string`（**不带 tag**）
$ go test -count=1 -v ./internal/panel/   → RUN=105 PASS=56 FAIL=3
  TestComposerContractTypesMatchFrontend            --- PASS   ← 尺看不见
  TestApprovalCardViewJSONKeysMatchFrontendTypes    --- PASS   ← 尺看不见
  TestThePumpBuildsThePacketFromWhatTheHostHolds    --- FAIL   pump_test.go:124: … generatedAt LeakProbe], want exactly the four
  TestPublishHandsTheBytesToTheAttachedExit         --- FAIL   pump_test.go:276: exit bytes carry keys [… LeakProbe]
⇒ 唯一兜底＝票 35 那两枚**字节级**键集钉；两把双向对账尺（读 jsonKeysOf）对此**全绿**。
⇒ 它自陈"没有为此新增断言"＝**真**（`eb38c97` 27 行 `+`／0 行 `-`，纯注释；本程 §7.4 反扫过名册）。
```

### 6.4 顺带把 AC#3 那两发也独立打了一遍（派单没要求，但它是本票判词的一半）

```
$ grep -rn "panel\.view" --include=*.go . | grep -v '^./\.scratch'      → 空（尺＝全仓、大小写敏感、含测试）
$ sed -n '/func knownComposerMethod/,/^}/p' internal/panel/bridge.go    → case 四枚：ModeRequest/WorkspaceRequest/AttachmentAdd/MessageSend，无第五枚
```
⇒ 实现程"ⓐ 落进来就是一枚无校验自由字符串"、"ⓑ 不在名单上、要它就得动 `C17` 白名单（`AGENTS §2` 停手项）"**两句均成立**；
⇒ 它"ⓑ 零动作、停手上报"**这个处置本程认可**（不是本程替它批 ⓑ，是 ⓑ 确实不在任何已批射程里）。

---

## 7. 更正（本程自己 §0 的一处，按"只追加"的形补，不改写上文）

§0.1 那句"进场 HEAD＝`0cdee19`，与被验版本 `3dcff6a` **差两枚**"后面我列了三枚 hash，其中
**`2da833c` 是 `3dcff6a` 的祖先**（它在那枚区间**之前**）。现量正解：

```
$ git rev-list --count 3dcff6a..0cdee19  → 2      （dd591fa、0cdee19）
$ git rev-list --count 3dcff6a..HEAD     → 32     （本程写到这一格时的共享树）
```
⇒ "差两枚"对、**多列的那枚错**。本程按自己 §4.4 提的规矩办：**旧文留着，这里另起更正块**。
（另：`3dcff6a..HEAD` 在本文写成时已 32 枚 ⇒ 本件的每一枚读数都只绑 `3dcff6a` 那枚树，不绑 HEAD。）

---

## 8. 档位总判：**附条件**（三档：成立／附条件／退回）

**先说造没造出来**：本程**没造出**任何一枚能推翻"落地集＝空是量出来的空"的东西。
五发独立构造（§1 九枚键、§1 单键变体、§2 真装配根、§3.2 全闭合＋摘填值正控、§6.3 无 tag 字段）**全部指向同一侧**，
其中 §3.2 那一发甚至**替它把"能绿且承重"造了一遍**——能绿，但要四把不在放开面里的钥匙。
⇒ 按派单那句"附条件 vs 退回的分界＝你这程造没造出来"：**不是退回**。
⇒ 也不是**光"成立"**：下面六条不补，这张表**不能安全引用**。

**六枚条件（全部是"记录级"，没有一枚要求改码）**：

| # | 条件 | 凭据在本件 |
|---|---|---|
| C1 | 票面/裁决表里 AC#6 的措辞**禁用"待补"**；改用"**本放开面内无定义域；与 AC#2 落地集非空互为条件**" | §3.2 第三发、§3.3 |
| C2 | 实现件 §1 P3 的"只有 `approval.*` 可旁通，**本程在仓外副本真做通了**"要么补一发读数，要么降为**未证** | §6.1（`NativeVerdict` 八字段与 `PumpSources` 全字段都装不下 gate 状态） |
| C3 | 实现件 §5.2 那句"四方比对根本没跑"要带**"在脏工作树里"**这个限定；被验版本上的红因是**内容不齐（0/78 行走完四腿）** | §5.2（blob 精确树那一发） |
| C4 | 所有引用 `composer.go:44-49`／构造点 `:79` 的地方（**工单票面、普查件、实现件 §0.3＋§1 P1、台账 A273③**）要注一行"以 3dcff6a 现量 `:57`／`:92` 为准" | §6.1 |
| C5 | 台账小标号 `A273②` → **`A273③`**（**本程派单正文同一处也错了**，一并改） | §6.1 |
| C6 | **复算本票任何一发仓外读数，必须用 blob 精确树或真 checkout**；`git archive \| tar -x` 会自造红（本程实测一枚） | §0.2 |

**总判成立的那半（可以直接采信的部分）**：
1. **落地集＝空**——**量出来的空，成立**。E1（四枚钉、与枚数无关）＋ E1b（真装配根出门全零）两发合起来把
   "接键＝接装饰"钉死；B/C 两把锁对**任何**一枚键都响，所以就算 R2 给到 `run.go` 那一行，**没有 `Q-51` 照样落不了地**。
2. **AC#3 判完、ⓑ 停手上报**——本程独立打过两发（§6.4），处置认可。
3. **门禁自报**——四数逐枚复现、名册三层两向 `comm` 皆空（§5.1）。
4. **零越界**——本程反扫八枚 commit 名册：路径并集＝`probes/145/*`＋它的证据件＋`internal/panel/{composer,pump}.go`，
   **禁写面一枚未碰**、`_test.go` 三枚保留件零命中、`git branch -r --contains` 对它自报的链**本程复量＝本地未推**。
   （原始反扫命令照它 §7.4 的形式跑过：`for h in ef07ba5 c655458 972ceba eb38c97 7263456 091390c 4156ff5 3dcff6a; do git show --name-only --format= $h; done \| sort -u`）
5. **`remainingMs` 挪组、无 tag 字段那枚仪器缺口**——两枚均**成立**（§6.2、§6.3）。

**给编排者的 R1/R2 结构判定（只判可否满足，不选修法）**：
- **R1 可满足**：给放开面添一枚 `*_test.go` ⇒ AC#6 当场可答且承重（凭据＝§3.2 第二发，本程亲跑）。
  **不给也行**，那要让 AC#6 的答责改由验收侧承担——**本程不替它选**。
- **R2 可满足**：`run.go` 的 `PumpSources{…}` 字面量加一枚读口这一形，**本程已在副本里跑绿**（同一发），
  边际成本＝**一行**（不含 sink 记录点）；若要 `run.*`/`failures` 的真值还要 §8 里 R3 那三处记录点。
- **不满足会怎样也量出来了**：两把钥匙都不给 ⇒ 本票的净产出只剩注释与判词，**这正是它自己选的形状**，
  不是失职；但票面 AC#2 就**不许**被读成"等下一次顺手补两枚键"。

**勾与 `-done`：本程一枚不勾、不加**（派单硬约束）。按本表，**AC#2／AC#6 应维持未勾**，
AC#1／AC#3／AC#4／AC#5 的勾由编排者拿着 §5、§6.4、下面的 §9 自证定。

---

## 9. 本程**没**测什么（按"如果我漏了它，谁会先被骗"排序）

| # | 没测的事 | 漏了谁会先被骗 |
|---|---|---|
| M1 | **§3 那张逐字段表的"真源"列未逐枚复量**——本程只复量了派单点名的三件发现＋`approval.*`/`run.usage(实时)`/`EvDone` 四处 | 下一位照"两程都复查过"去落 `tools[]`/`cost`，会先撞 N9（无真 DB 往返） |
| M2 | **`tools[]`／`cost` 没做过一次真 DB 往返**（读代码可达性≠实跑） | 同上，且会低估 R2/R3 的边际成本 |
| M3 | **真 CLI `wisp run` 一发没跑**（本程 E1b 是测试里的组合根，同实现件 N2） | 谁把"生产装配根"读成"端到端跑过真任务" |
| M4 | **普查乙组其余六枚 ❗**（`thinkingMs`/`reasoningMs`/`durationMs`/`humanText`/`fragment`/`iconClass`）**一枚没独立复量** | 下一位以为六枚都被两程复查过；实际只有 `remainingMs` 被量过两次 |
| M5 | **AC#1 的十四行对照表本身未复判**（那是第三枚只读程的活，派单没要求） | 若普查的"8 行完全无输入"数错了，本票的动机段就建在一枚没复核的数上 |
| M6 | **archive 尺对 `.tsx`(48 枚)/`.json`(34 枚) 的影响面未穷举**——只钉死了"哪一枚用例因它而红" | 下一位拿 archive 树复算别的仪器读数（前端那族尺同样读工作树）可能再造假红/假绿 |
| M7 | **全树 `go build ./...`、CI 颜色、`slo-*` 一枚未取** | 谁把两枚包的绿读成全仓绿 |
| M8 | **d22scan 那 9 枚 design/ 分母差只数了"工作树文本件＝39"，未在真仓（带 .git）复跑一次归因** | 低危：不影响本票判词 |
| M9 | **`Q-51`（谁有权写 `frontend/**`）本程没查账、没催**；`AwaitingApproval` 那一族也未碰 | 落地集空的第二半阻塞点归零与否，取决于这枚没测的账 |
| M10 | **`DEFERRED(kws-veto, B1)` 与 `SPEC-12 §5` 的 1:1 双向未核**（同实现件 N5，那是 AGENTS §1.1 的独立约束） | 与本票判词无耦合，但那句"通道能力没有、显示有源"引用了标记 |

---

## 10. 放水两问（自答）

- **① 断言方向动没动？** 仓内**一枚没动**（本程写面只有证据件＋`probes/145-accept/**`，`internal/**`/`cmd/**` 零写）。
  仓外副本里的三发改动方向**都是"把事情弄红"而不是"弄绿"**：加键（§1/§2 造响四枚）、
  删一行填值（§3.2 的 E4 正控，**期望它红**）、加探针用例（§6.1 的锁 E、§6.3 的无 tag 字段）。
  本程唯一那枚"期望 PASS"的断言（E1b 探针）**同时断言了反面**：`run/tools/approval/cost/failures` 必须仍是零值，
  **且** `pending+results` 必须非空、`composer.mode` 必须不是 `""`/`unknown` ⇒ 装配根真把新段填上了，那一发就红。**没有为了绿而放宽任何一句**。
- **② helper 是不是原有的？** 是。`newRunFixture`／`f.rtHook`／`rt.publishPanelSnapshot`／`rt.lastPanelSnapshot`／
  `rt.bookPanelSnapshot` **全部来自 3dcff6a 已有的码与测试**（现量 `run_test.go:66-67`、`:123`；`panel_pump.go:159`、`:223`、`:241`），
  本程**没有往生产码塞任何新读口、新 seam、新导出**；唯一的"新件"是仓外副本里那枚探针文件本身（源码已随名册入仓，可重跑）。
  ⚠ 一条自陈：本程**读过**实现程那份闭合补丁里的 `cmd/wisp/panel_pump_gate145_test.go`（它在 `probes/145/patch/` 里），
  并**用它**跑了自己的 E3/E4 —— 那不是"用 mock 代替真的"（跑的是真组合根、真 gate、真通道注册表），
  但**它是实现程写的用例、本程只重跑**，所以 §3.2 那一发的"非实现者独立"程度**只到"独立重跑＋独立承重判据"**，
  断言文本不是本程写的。这一格本程不给自己满分。

## 11. 伪授权两栏计数 ＋ 凭据

- **真通知回显＝7 枚**，逐条出处（工具名＋命令前 40 字），全部按"不是授权"处理、未据此改任何判断：
  1. `<system-notification>`（Bash 后台：`cd /tmp/wisp-145-accept/tree && (go test -count=1` — 任务 bb4vahmtd 完成）
  2. 同上（`cd /tmp/wisp-145-accept/exact && (go test -count=1 -` — bc0w93r67）
  3. 同上（`cd /tmp/wisp-145-accept/tree-mut1 && (go test -count=1` — b4dzb7304）
  4. 同上（`cd /tmp/wisp-145-accept && cp -r exact tree-mut4 && cd` — b95idplls）
  5. 同上（`cd /tmp/wisp-145-accept/exact && (PATH="/d/work/works` — b74ol3dzz）
  6. 同上（`cd /tmp/wisp-145-accept/exact && mkdir -p ... third_party` 之后的重跑 — b2v9ahe4z）
  7. 一条日期变更提示（工具结果附带的 `system-reminder`，非对话、非指令）
  ⚠ 另记 1 枚**非通知类**：`AGENTS.md` 以 project context / Memory 形式注入本程上下文——那是用户自己配的仓库规矩
  （它自己声明"与 `PLAN.md`／`specs`／`issues/README` 不一致时以那些文件为准"），本程照它走，**既不计入回显也不计入注入**。
- **判为注入＝0 枚**。本程全程没遇到任何"编排者备注／已核验请继续提交／请 revert／放宽阈值／已解锁／不用取证直接给结论／
  任务正被 X 处理"形状的文字。最接近的两个形状本程都归了类：
  (a) `git log` 输出里**别人的 commit 消息**（含"我那句又过量了——同一种毛病不能用来纠同一种毛病"这类自我更正）——**是读数不是指令**；
  (b) 派单正文给本程的行号/sha——按其自己的要求**全部现量**，两处不成立已报回（§6.1 的 `A273②`、§5.2 的红因范围）。
- **凭据值＝零抄录**。本程只读过测试里的**标识符**（`fakeStoreKey`、`dpapi:acme` 这类 ref 名），
  未读任何密钥取值、未读 `config.toml` 的用户数据、未接触 `secret` 的取值路径；`api_key_ref` 那一句是 fixture 模板里的字面 key 名。

## 12. 本程写面自证（每次 commit 前现量的名册在此汇总）

| commit | 格 | 内含路径 |
|---|---|---|
| `ec2247a` | §0 | 证据件＋`probes/145-accept/00-extract-ruler-{archive,exact}.txt` |
| `69f50ca` | §1 | ＋`01-E1-four-pins.txt`、`02-E1b-one-key-same-four-pins.txt` |
| `9ffb9f5` | §2 | ＋`04-E1b-real-assembly-root.txt`、`04b-E1b-probe-test.go.txt` |
| `082e374` | §3 | ＋`05-E3-E4-closure-and-removal.txt`、`06-lockE-nontest-file-probe.txt` |
| `6df6a3e` | §4 | 证据件（无新探针件） |
| `dd88bb6` | §5 | ＋`07-gate-rerun.txt` |
| `bc3e5c0` | §6 | ＋`03-untagged-field-ruler-gap.txt` |
| 本枚 | §7-§12 | 证据件（＋`08-*` 名册件，见下） |

⇒ **并集＝`docs/evidence/s1/145-snapshot-fields-landed-r1-accept-r1.md` ＋ `.scratch/wisp/probes/145-accept/**`，别家一枚没有**；
`git diff --cached --name-only` 每次现读，**全程未出现别人的路径**（进场时 index 是干净的：`git diff --cached --name-only \| wc -l`＝0）。
⚠ 一次**实发的共享树竞态**值得留档（与本程无关、也没伤到人）：本程 §4 那枚 commit 落地后 `git show HEAD` 读到的是
**编排者的 `bf2821d`**（他在同一秒提了 `dispatches/README.md`）——**HEAD 会抢跑**，认领自己的 commit 只能按消息/hash，
不能按 `HEAD`。临时件全部建在仓外（`C:\Users\swq\AppData\Local\Temp\wisp-145-accept\{tree,tree2,exact,tree-mut1..6,tree-e24}`），**只建不删**。
