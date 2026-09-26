# 151 验收 r1（非实现者对抗验收程）— `Bridge.CloseTask` 只开不关这一案的第二枚验收表

> 本件是**验收表**，不是实现件。派单＝`.scratch/wisp/dispatches/2026-09-26-103x-accept-151-r2.md`。
> 被验版本＝**`4e16976`**（票 151 自己链的最后一枚）。实现件＝`docs/evidence/s1/151-task-scope-never-closed-r1.md`（只读，未改一字）。
> 票面 151 与票面 154 **本程一字未动**（派单与编排者都点名了这条，见末节的现量）。

## 第 0 节 锚点与环境（本程现读，不抄派单）

```
$ git rev-parse --short HEAD          ->  b0f68b2      （本程开工那一刻的工作树头）
$ git branch --show-current           ->  dev
$ git cat-file -t 4e16976             ->  commit
$ git cat-file -t 5d46f24             ->  commit
$ git cat-file -t 45c920e             ->  commit
$ git log --oneline -1 4e16976        ->  4e16976 evidence(151 r1): 判定＋修法交件……
```

- **取版**：`git archive 4e16976 | tar -x -C D:/tmp/wisp151accept`（仓库外）。**未在仓库内建 worktree、未 checkout、未 stash/reset**（AGENTS §1.4）。
- 变异与复算用三枚仓库外树：
  | 树 | 内容 | 用途 |
  |---|---|---|
  | `D:/tmp/wisp151accept` | `4e16976` 纯净树 | 摘一味／摘另一味的变异尺 |
  | `D:/tmp/wisp151accept-base` | 同一枚树 ＋ `third_party/sherpa-onnx`（从工作树复制的三枚 dll） | 改后门禁全跑、`go vet`、`gofumpt`、`d22scan` |
  | `D:/tmp/wisp151accept-pre` | `5d46f24^` 的 `run.go`＋`bridge.go`（用 `probes/151/run.head.go`/`bridge.head.go`，本程 `sha256sum` 现核＝与 `git show 5d46f24^:<file>` **逐字节相同**），新用例文件摘成 `package main` | 改前门禁全跑（不走 `-overlay`，见第 5 节） |
- **PATH 用 shell 自己的路径形**（`/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx`），未用 `pwd -W` 盘符形。
  ⚠ 本程第一发全包门禁**就是这么死的**：纯净树里没有 `third_party/`（`git archive` 不带未跟踪件），
  8 枚"要起真子进程"的用例报 `no native DLLs in ..\..\third_party\sherpa-onnx` 并转红
  （原始读数＝`probes/151-accept/gate-base-cli-v.txt`：RUN=132 顶格 70 红 8）。
  **按派单第 55 行的次序处理＝那是"仪器没跑到"，不是红**；补上 dll 目录后同一枚树重跑＝RUN 138／顶格 78／红 0（第 5 节）。
- 本程写面只用过两枚路径：`docs/evidence/s1/151-task-scope-never-closed-r1-accept-r1.md`（本件）＋ `.scratch/wisp/probes/151-accept/**`。
  其余全在仓库外。临时件**只建未删**。

## 第 1 节 票面 5 枚 AC ↔ 本程 5 格 双向对账（派单攻击点②）

现量（锚点 `4e16976` 的那枚票面，两把尺都给）：

```
$ git show 4e16976:.scratch/wisp/issues/151-*.md | grep -c '^- \[ \]'   ->  5
$ git show 4e16976:.scratch/wisp/issues/151-*.md | grep -c '^- \['      ->  6
```

⇒ **票面 AC 框＝5 枚**（`:27 :29 :34 :36 :39` 五行，逐行现读＝AC#1..AC#5）。实现件 §10 写"票面 **6** 枚 AC 框"。
那枚"6"最省事的复现尺就是上面第二行：`^- \[` 会把 `:56` 那行 Progress log 的 `- [2026-09-26 00:5x +08]` 一起数进来。
⇒ 判：**数错的是实现件，且是一枚"尺"错、不是"格"错**——它自陈未勾的框实际是 5 枚，5 枚它都交了读数，
所以**没有哪一格因此没人裁**（下表第五列＝本程那一格的号，五格全在本件里）。

| 票面 AC（原文标题，现抽） | 实现件那一节 | 本程那一格 | 本程裁决 |
|---|---|---|---|
| AC#1 先量"不关"今天到底改变了什么 | §2 | 第 2 节 | **成立，但一句结论要收窄**（可达性那一问，见第 2 节末） |
| AC#2 会响的检（硬核心） | §3 | 第 3 节 | **成立，两味都有牙**；但 §3 那句"两枚在未修码上结构上跑不起来"**本程现量证伪其一半** |
| AC#3 修法只在"谁拥有关闭"这一层（三选一） | §4 | 第 4 节 | **成立**，且"不许顺手改 `OpenTask`"那半句逐字节核过 |
| AC#4 契约轴（零字节） | §5 | 第 5 节前半 | **成立**（三枚 commit 的名册并集＝5 个路径，一枚不在禁写面里） |
| AC#5 门禁（逐包单跑） | §6 | 第 5 节 | **六数全中、名册两向对得上**（本程独立跑，未走 `-overlay`） |
| （票面另四条"现量的形状"＝AC 之前的第 1 节，不是 AC 框） | §1 | 第 2 节开头 | 四条本程都重走过，一条尺写清了 |

派单那七条攻击点落在哪：① →第 2 节；② →本节；③④⑤ →第 3 节；⑥ →第 5 节；⑦ →第 6 节。总判＝第 7 节。
