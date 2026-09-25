# 149 对抗验收 r1 — 三处未钉出口的修法、那把"第一版是坏的"尺、以及坏尺期间留在树里的读数

- 被验版本＝**`c3f7224`**（票 149 的四枚：`b417d31` 码＋尺／`2f5c7f9` probes／`a055d9f` 件 0-11 格／`c3f7224` 件第 12 格）
- 本程角色：**对抗验收程（非实现者）**，只读＋只 commit 本件
- 票面：`.scratch/wisp/issues/149-the-corrupt-leg-of-the-offset-field-has-zero-teeth-the-new-comment-claims-a-sufficient-condition-with-counterexamples-and-the-exited-branch-drops-the-summary.md`
- 实现件：`docs/evidence/s1/149-three-unpinned-outlets-r1.md`（§0–§12）
- ⚠ 实现件与本程派单的每一句都按**未验证断言**处理；推翻写进 §N-1「结论修正记录」

---

## 第 0 格　锚点与工作树（进场第一读）

**判据**：被验版本必须按 sha 取；脏工作树**不得**当被验版本；快照若要跑 `cmd/wisp` 必须先确认原生 dll 在位，
否则拿到的是假红（派单 §0 那条坑）。

### 0.1 现量命令与原文

```
$ git rev-parse HEAD
4d0866ad2ea26a2d53d7086eb42a14505445cef2

$ git log --oneline -8
4d0866a evidence(139 AC#1 第 1 格): 零留痕量成读数……
c3f7224 evidence(149 第 12 格): 入库字节再跑一遍三发关键变异——判据未随版本变、行号随版本变
a055d9f evidence(149 第 0-11 格): 三发今天读数、corrupt 支的承重矩阵、exited 支的可达性剖面、门禁四数与名册差集
612590c docs(台账 A267): 票 138 撞顶＝落盘不是重派……
3f6322f evidence(138 落盘代理未提交的尾段): ……
2f5c7f9 probes(149): 变异尺与逐发原始读数落盘（57 枚，只建不删）
f8623c7 evidence(138 AC#4 门禁 第 6 格): ……
b417d31 fix(149 AC#2+AC#3+AC#4): corrupt 支的 offset 改填解码器自己点名的位置，并补三枚会响的用例

$ git status --porcelain | wc -l
84

$ git merge-base --is-ancestor c3f7224 HEAD && echo yes
yes

$ git log --oneline c3f7224..HEAD
367fdc0a evidence(138 验收 r1 第 0-1 格): ……
fe03449a docs(台账 A268): 票 149 交件；它自陈"第一版组合变异的尺是坏的"，而坏尺期间那批读数还留在树里
4d0866ad evidence(139 AC#1 第 1 格): ……

$ go version; D:/work/base/gopath/bin/gofumpt.exe --version
go version go1.27.1 windows/amd64
v0.12.0 (go1.27.1)
```

`git status --porcelain` 里**与本票地界相交**的只有两类，逐枚点名：

```
$ git status --porcelain | grep -E 'internal/|cmd/wisp|\.scratch/wisp/probes/149|docs/evidence/s1/149'
 M internal/agent/compress.go
 M internal/agent/loop.go
```

⇒ 进场时 `cmd/wisp/**`、`.scratch/wisp/probes/149/**`、本票两枚证据件**在工作树里都是干净的**；
脏的是 `internal/agent/**`（139 那枚程在写，与本票地界无交集）＋ `frontend/**`／`design/**`（另一枚会话，
票面 AC#5 已明示"不碰、不还原、不计进任何零命中宣称"）。**本程一枚未碰。**

### 0.2 锚点相符性

- 派单说的 sha `c3f7224` **存在且是 `dev` 的祖先**（`git merge-base --is-ancestor` → yes），
  四枚 commit 的标题与派单列的四枚**逐枚对得上**（`b417d31`／`2f5c7f9`／`a055d9f`／`c3f7224`）。
- ⚠ **派单那句"此刻另一枚写码程正在 `internal/agent/**` 里跑票 139，HEAD 会动"** —— **相符**：
  本程进场 HEAD＝`4d0866a`，且 `c3f7224..HEAD` 之间已经又落了 3 枚（138 验收、台账 A268、139 第 1 格）。
  ⇒ 本程**所有"某文件长什么样"一律按 sha 取**，跑测试一律在仓外快照，不碰脏树。
- ⚠ 派单 §0 的"快照不带 `third_party/` 原生 dll"这条**本程自己撞过并绕开**：

```
$ ls third_party/sherpa-onnx/*.dll | wc -l
3
$ find third_party -name "*.dll"
third_party/sherpa-onnx/onnxruntime.dll
third_party/sherpa-onnx/sherpa-onnx-c-api.dll
third_party/sherpa-onnx/sherpa-onnx-cxx-api.dll

$ mkdir -p /d/work/tmp/wisp149-accept/snap-c3f7224-a
$ git -c core.autocrlf=false -c core.eol=lf archive c3f7224 | tar -x -C /d/work/tmp/wisp149-accept/snap-c3f7224-a
$ ls /d/work/tmp/wisp149-accept/snap-c3f7224-a
AGENTS.md  README.md  cmd  deps.toml  design  docker  docs  frontend  go.mod  go.sum  internal  models  scripts  tools
```

⇒ 现量：`git archive` 出的树**确实没有 `third_party/`**（顶层名册里就没有那一枚）——派单那条坑**成立**。
本程建快照时**先把 dll 目录复制进去再跑**，所以本件里任何"红/绿"都不是 `0xc0000135` 那一种（区分依据见第 7 格）。

### 0.3 读数

- 锚点相符、分支 `dev`、HEAD 已漂移 3 枚且**漂移全是别人的证据件/台账**（零 `cmd/wisp` 字节）。
- dll 在位（3 枚）；`go` ＝ 1.27.1，`gofumpt` ＝ v0.12.0（现读，未背数）。
- 仓外快照 `snap-c3f7224-a` 建立；`probes/149` 单独导出为 `logs-at-anchor/` 供第 1 格清点。

**放水两问自答**：① 本格零断言、零代码改动（只读＋建仓外目录）；② 本格不使用任何 helper。

**第 0 格判定**：**成立**。

**本程没测什么（本格）**：没验 `c3f7224` 之外那三枚 commit 的**内容**是否真等于派单说的角色（只对了标题与路径名册）；
没查 `third_party/` 那 3 枚 dll 与实现程跑门禁时用的是不是同一批字节。

---

## 第 1 格　派单 §2 之(2)：坏尺期间的读数清点——**那批读数和它自己的叙述方向相反**

**判据**（派单 §2 第 2 条＋台账 A268②(乙)）：`probes/149/post/` 与 `probes/149/combos2/` 里同名 `x-*` 各有一份 ⇒
必须**逐枚判定哪些属第一版（坏尺期间）、哪些是复算版**，第一版那批一律作废并点名。
本格的尺**不是**实现件的那张表，是**日志文件自己**：`^--- FAIL:` 只数顶层（子测试那一种带缩进，
`grep -c '^--- FAIL'` 与裸 `grep -c -- '--- FAIL'` 在本仓是两把不同的尺，见第 7 格）。

### 1.1 现量命令与全表

```
$ git archive c3f7224 .scratch/wisp/probes/149 | tar -x -C <仓外>/logs-at-anchor
$ python <仓外>/tabulate_logs.py          # 对每枚 .log 取：=== RUN 枚数 / ^--- PASS 枚数 / ^--- FAIL 名册 / 首行 time=
```

（脚本落在仓外 `D:\work\tmp\wisp149-accept\`，只建不删。）

**根层 12 枚**（AC#1 与门禁；`post/`、`combos2/` 见下）：

| 日志 | RUN | 顶层 PASS | 顶层 FAIL | 内嵌时间戳 |
|---|---|---|---|---|
| `pre-asis` | 17 | 10 | 0 | 23:10:19 |
| `pre-B14-corrupt-fill-to-zero` | 17 | 10 | **0** | 23:10:25 |
| `pre-B15-complete-fill-to-zero` | 17 | 9 | 1 | 23:10:33 |
| `pre-B13-exited-drops-summary` | 17 | 10 | **0** | 23:10:39 |
| `pre-B5-timeout-drops-summary` | 17 | 7 | 3 | 23:10:46 |
| `probe-corrupt-census` | 1 | 1 | 0 | 23:06:58 |
| `newtests-on-unfilled-code` | 20 | 12 | 1 | 23:18:46 |
| `gate-pre-full` | 133 | 73 | 0 | 23:19:50 |
| `gate-post-full` | 136 | 76 | 0 | 23:19:52 |
| `gate-post-full2` | 136 | 76 | 0 | 23:24:11 |
| `gate-worktree-post` | 136 | 76 | 0 | 23:24:14 |
| `d22scan-worktree` | 70 | 30 | 0 | — |

**`post/`（31 枚）与 `combos2/`（11 枚）同名碰撞对**（本格的核心表）：

| 变异名 | `post/` 顶层 FAIL | `post/` 时间 | `combos2/` 顶层 FAIL | `combos2/` 时间 | 两遍是否一致 |
|---|---|---|---|---|---|
| `x-p1-d11a` | 1 | 23:13:28 | 1 | 23:15:22 | 一致 |
| `x-p1-d11b` | 1 | 23:13:33 | 1 | 23:15:28 | 一致 |
| **`x-p1-d11a-d11b`** | **1** | 23:13:38 | **0** | 23:15:36 | **不一致** |
| `x-p1-d11e` | 0（RUN=19） | 23:13:45 | 0（RUN=19） | 23:15:55 | 一致 |
| `x-p7-d13a` | 1 | 23:13:53 | 1 | 23:16:01 | 一致 |
| `x-p7-d13b` | 0（RUN=19） | 23:13:59 | 0（RUN=19） | 23:16:09 | 一致 |
| `x-p12-d11c` | 3 | 23:14:05 | 3 | 23:16:17 | 一致 |
| `x-p12-d11d` | 2 | 23:14:10 | 2 | 23:16:24 | 一致 |
| `x-p10-d12` | 1 | 23:14:15 | 1 | 23:16:31 | 一致 |

`combos2/` 另两枚无 `post/` 对应：`x-p2-d11a-d11b`（0 红，23:15:44）、`x-p3-d11a-d11b`（0 红，23:15:50）。

**唯一那发红句原文**（`post/x-p1-d11a-d11b.log`，入库字节）：

```
$ git show c3f7224:.scratch/wisp/probes/149/post/x-p1-d11a-d11b.log | sed -n '38,47p'
=== RUN   TestSLO149CorruptSentenceNamesThePositionTheDecoderObjectedAt
    slo_report_144_windows_test.go:646: 8 of 8 corrupt shapes do not name the byte position the decoder objected at;
                           first: second-comma-at-26: subjectReportRead.offset is 0, want 27 (the byte the decoder objects at)
--- FAIL: TestSLO149CorruptSentenceNamesThePositionTheDecoderObjectedAt (0.00s)
...
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.568s
FAIL
```

同一发在 `combos2/` 那份：`--- PASS` 三枚全绿、`ok github.com/CarlosShao/wisp/cmd/wisp 0.504s`、顶层 FAIL **0**。

### 1.2 读数与判定

时间戳把顺序钉死了：`post/` 那批（23:11:02→23:14:15）**早于** `combos2/`（23:15:22→23:16:31），
且 `post/` 里 9 发 `x-*` 的先后**正好是 `harness.py` 里 `MUT` 字典的书写顺序**（跳过 `x-p2`/`x-p3` 那两发，
它们只在 `combos2/` 有）⇒ **`post/` 那批就是"第一版尺"那一遍**，与实现件 §11 的说法（"`post/` 31 枚＝…＋9 发组合的第一遍（坏的）"）相符。

⇒ **但实现件对那一批的叙述与它自己点名的日志相反。** §8 第 6 行写的是：
"那遍的 `apply()` 把同文件的第二枚替换写成独立 overlay，只留下最后一枚 ⇒ **第一版那条『0 红』是尺坏了**…
两遍日志都留着：`probes/149/post/x-p1-d11a-d11b.log`（**坏的那遍**）与 `probes/149/combos2/x-p1-d11a-d11b.log`（算数的那遍）"。
现量：被点名为"坏的那遍"的那枚日志里躺着的是 **1 红（顶层 FAIL 枚数＝1，包级 FAIL）**，不是 0 红。

⇒ 把 `harness.py` 现在的 `apply()`（串联版）与那枚红句**合起来**可以做故障复原，而且这一步**不是假设、有字据**：
坏版是"每枚 spec 各自读原文、各自写 `mut/<名>/<basename>`、各自塞进 `Replace`"⇒ 同一枚文件的第二枚替换**把第一枚的产物覆盖掉**，
`x-p1-d11a-d11b` 在坏版里实际跑的是 **`p1` ＋ 只摘 `d11b`**（`d11a` 那枚字段 leg **没被摘掉**）。
⇒ **那一遍留下的红句本身就指认了是哪条腿在响**：红句原文是
`subjectReportRead.offset is 0, want 27`，而全仓只有字段腿那一支会产出这串字——

```
$ git show c3f7224:cmd/wisp/slo_report_144_windows_test.go | sed -n '634,641p'
		if obs.offset != int64(tc.want) {
			report(i, fmt.Sprintf("subjectReportRead.offset is %d, want %d (the byte the decoder objects at)", obs.offset, tc.want))
		}
		read, named, offset, ok := slo149ContradictionOf(obs.summary())
		switch {
		case !ok:
			report(i, fmt.Sprintf("sentence %q is not the contradiction shape", obs.summary()))
		case offset != tc.want:
```

而 `D11A_FIELD_LEG_OFF` 要摘的正是 `:634` 那行 `if obs.offset != int64(tc.want) {` ⇒
**红句出自字段腿 ⇒ 字段腿在坏尺那一遍是活的 ⇒ `d11a` 没有被应用**。字段 leg 活着，就当场抓到常量 0 ⇒ 1 红。
⇒ **所以"逃逸"（0 红）是修好尺之后才出现的读数，不是坏尺的产物。** 本仓那句"下一位会挑对自己有利的那份"
在这里以**反方向**发生了：坏尺那遍给的是一颗**假安心**（两条 leg 看着都红 ⇒ 像是有牙），
真把两条都摘干净反而 0 红。

**本格判定**：**附条件入账**。
- 成立的部分：派单 §2 第 2 条那句"`post/` 与 `combos2/` 里同名 `x-*` 各有一份"**为真**——`post/` 的 9 枚 `x-*`
  **全部**在 `combos2/` 有一份同名对应（`combos2/` 另多两枚 `x-p2-d11a-d11b`／`x-p3-d11a-d11b`，那两发只在复算遍跑过）。
  ⚠ **顺带把整棵 `probes/149/` 的同名碰撞点齐**（派单只点了 `x-*` 一族，实际不止）：
  `post-asis.log` 另有一对在 `post/`×`post-final/`；`p1-corrupt-fill-to-zero.log` 与 `p7-exited-drops-summary.log`
  **各三份**（`post/`×`post-final/`×`logs-verify/`）、`p10-multidoc-takes-error-offset.log` **两份**
  （`post/`×`logs-verify/`）⇒ 全树共 **13 组同名、合计 28 枚日志副本**（现量：
  `git ls-tree -r --name-only c3f7224 -- <probes/149> \| grep '\.log$' \| sed 's/.*\///' \| sort \| uniq -c \| awk '$1>1'`），
  其中**只有 1 组两遍读数不一致**（`x-p1-d11a-d11b`），
  而那 1 组正是 AC#2 承重的全部依据。
- **作废点名（本程裁定）**：`probes/149/post/` 的 **9 枚 `x-*` 日志**＋`post/post-asis.log` 一律降级为
  **不作任何格的凭据**（`post-asis` 另有一枚 `post-final/post-asis.log` 与它同读数，可作替代凭据）。
  理由不是"它们是坏的"，而是**其中至少一枚（`x-p1-d11a-d11b`）的读数与它被引用的用途不符**，
  同一批里其余 8 枚无法逐枚自证它们那遍没有踩到同一条覆盖路径（本程能证的只是：那 8 枚两遍读数一致，见 §1.1 表）。
- 不成立的部分：**实现件 §8 第 6 行对那批日志的描述**（详见 §N-1 修正记录第 1 行）。
- ⚠ **一处不改代码也要点名的落盘卫生问题**：仓里**没有任何一个文件**标出"`post/` 那 9 枚是坏尺版本"。
  实现件 §2.4 那句"全部日志在 `probes/149/post/`、`probes/149/combos2/`"把两批**并列**，
  下一位按 §2.4 的表去 `post/` 找 `x-p1-d11a-d11b` 的凭据，会拿到与表**相反**的那一发。
  台账 A268② 已经把这写成本程第一格判据 ⇒ 本格就是那句"点名"。**修法建议（未验证断言，只说治到哪一发）**：
  把 `post/` 那 9 枚改名到 `post-brokenruler/` 或在 `probes/149/` 落一枚一行的版本说明——
  它治得到的是"下一位挑 `post/` 那份 1 红当凭据、据此宣布两条 leg 有兜"这一发；治不到"尺本身是否还坏着"（那是第 2 格）。

**本程没测什么（本格）**：
1. **没验坏版 `apply()` 的源码**——仓里入库的是**修好的** `harness.py`（本程读到的 `apply()` 就是串联版），
   坏版只存在于实现程的自述里。§1.2 里"**那一遍 `d11a` 没被应用**"这一步**有字据**（红句只可能出自 `:634` 那支），
   但"**为什么**没被应用＝独立 overlay 相互覆盖"仍是机制假设。⇒ 谁先被骗：想据此修 harness 而不先复现坏版的人。
2. **没独立复算那 9 枚两遍一致的 `x-*`**（本格只比对了日志字节里的顶层 FAIL 名册）；一致＝两把尺互证，
   但**都是它的尺**。第 2 格用本程自己的尺复算承重那一发。
3. **没查 `post/` 那批日志是否被同名重跑覆盖过**（`harness.py` 写 `log_dir/<名>.log` 是覆盖式的）⇒
   理论上存在"更早还有一遍真 0 红的 `post/x-p1-d11a-d11b.log`、被 23:13:38 那遍覆盖"。
   **本程无法从入库件判定那枚存在过**，所以 §1.2 只下"入库那版说的是 1 红"这一条，不下"0 红从未被读到"。

---

## N-1　结论修正记录（推翻实现件／推翻本程派单，都写这里）

| # | 原话（谁说的） | 本程现量 | 结论 |
|---|---|---|---|
| 1 | **实现件 §8 第 6 行**＋**派单 §2**＋**台账 A268②**："第一版 combos 里 `x-p1-d11a-d11b` 的那一发『0 红』是尺坏了的产物"（派单还写明"修好后在 `combos2/` 复算"） | `post/x-p1-d11a-d11b.log`（＝入库的第一遍，时间戳 23:13:38 早于 `combos2/` 23:15:36）顶层 **FAIL＝1**、红句是字段腿那句 `offset is 0, want 27`；`combos2/` 同一发 **FAIL＝0** | **推翻（方向相反）**：坏尺那遍交回的是**1 红**（假安心），**0 红 逃逸是修好尺之后才出现的读数**。三处（实现件／派单／台账）同源于实现程一句自述，本程按日志字节判 ⇒ 自述与它自己点名的凭据文件不符。派单 §2 那句"那一版 `x-p1-d11a-d11b` 的『0 红』是尺坏了"**不成立**，别照它做 |

## N-2　伪授权登记（两个数分开栏）

（本节随格追加，最后在第 9 格汇总定数。）
