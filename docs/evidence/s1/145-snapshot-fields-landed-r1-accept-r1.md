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
