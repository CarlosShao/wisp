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
