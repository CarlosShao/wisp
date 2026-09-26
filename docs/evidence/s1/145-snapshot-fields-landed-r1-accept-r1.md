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
