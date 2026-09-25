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

## 第 2 格　派单 §2 之(1)：**独立重建**那发组合变异（不拿它的 `combos2/` 日志当凭据）

**判据**：AC#2 承重的全部依据是"两条 leg 同摘 ⇒ `p1` 打不红"这一发**逃逸**读数。派单明令别拿实现程的复算日志
当凭据 ⇒ 本格用**本程自己写的尺**、自己按 sha 建的快照重做那一发，并且**给尺加一道它自己没装的证明**：
overlay 建完之后**回读要交给 Go 编译的那枚文件**，断言"新文本在里面、旧文本不在里面"，否则 FATAL 不输出读数。
⇒ 这一步是针对第 1 格那个失效模式（同文件两枚替换静默剩一枚）设计的：**本程的尺结构上就不可能悄悄少摘一条腿。**

### 2.1 快照与字节相符性（先证被验版本，再谈读数）

```
$ git -c core.autocrlf=false -c core.eol=lf archive c3f7224 | tar -x -C <仓外>/snap-c3f7224-a
$ cp -r <repo>/third_party/sherpa-onnx <仓外>/snap-c3f7224-a/third_party/     # dll 3 枚，先补再跑
$ for f in cmd/wisp/slo_windows.go cmd/wisp/slo_report_144_windows_test.go; do
    a=$(git show c3f7224:$f | md5sum|cut -d' ' -f1); b=$(md5sum < $S/$f|cut -d' ' -f1); echo "$f anchor=$a snap=$b"; done
cmd/wisp/slo_windows.go                      anchor=0403d5196f4bc0c20993dca7ceeea2b8 snap=0403d5196f4bc0c20993dca7ceeea2b8
cmd/wisp/slo_report_144_windows_test.go      anchor=4d07365e12e9e156771020f3b638dcfa snap=4d07365e12e9e156771020f3b638dcfa
```

⇒ 快照里那两枚文件与**锚点 blob 逐字节相同**（md5 相等）；跑法一律 `go test -count=1 -overlay … -v
-run 'TestSLO144|TestSLO147|TestSLO149' ./cmd/wisp/`，`PATH` 里带 dll 目录（第 0 格那条坑）。
本程的尺落在仓外 `D:\work\tmp\wisp149-accept\myharness.py`，日志落在同目录 `mylogs/`（只建不删）。

### 2.2 本程名册（12 发，含 4 发对照与 2 发探针）

`m-combo` ＝ 实现程那发的等价物：同一枚 GO 替换（`contradictionOffset(err, dec.InputOffset())` → `obs.offset = 0`）
＋同一枚 TEST 文件里的**两条 leg 各摘一枚**（`:634` 字段腿 `if obs.offset != int64(tc.want) {` → `if false {`；
`:641` 句子腿 `case offset != tc.want:` → `case false:`）。

| 本程编号 | 内容 | 顶层 FAIL | RUN | case 11 | 红句出自哪条腿（按消息文本认） |
|---|---|---|---|---|---|
| `m-asis` | 无 overlay（对照） | **0** | 20 | RAN | — |
| `m-corrupt-zero` | ＝它的 `p1`（靶） | **1** | 20 | RAN | `subjectReportRead.offset is 0, want 27 (the byte the decoder objects at)` ← **字段腿** |
| `m-legs-both` | 只摘两条腿、码不动（对照） | **0** | 20 | RAN | — |
| **`m-combo`** | **`p1` ＋ 两条腿全摘** | **0 ⇒ 逃逸复算成立** | **20** | **RAN** | 全程无 `first:` 行（`grep -c "first:" m-combo.log` → **0**） |
| `m-combo-bytes` | 取值→`int64(obs.bytes)` ＋两条腿全摘（＝它的 `x-p2-d11a-d11b`） | **0 ⇒ 逃逸** | 20 | RAN | — |
| `m-combo-revert` | 取值→`dec.InputOffset()` ＋两条腿全摘（＝它的 `x-p3-d11a-d11b`） | **0 ⇒ 逃逸** | 20 | RAN | — |
| `m-corrupt-revert` | ＝它的 `p3`（撤掉本票修法） | **1** | 20 | RAN | 字段腿 |
| `m-zero-fieldoff` | `p1` ＋ **只摘字段腿** | **1** | 20 | RAN | `sentence says offset 0, want 27` ← **句子腿单独接住** |
| `m-zero-sentenceoff` | `p1` ＋ **只摘句子腿** | **1** | 20 | RAN | 字段腿单独接住 |
| `m-wantcanary-live` | 只把 case 11 自己的期望常数 27 改成 999 | **1（`1 of 8`）** | 20 | RAN | 期望值是**活的**，不是摆设 |
| **`m-combo-canary`** | **两条腿全摘 ＋ GO 句子撤掉 `at offset %d`** | **2** | 20 | **RAN** | `sentence "…" is not the contradiction shape` ← **剩下那条 `!ok` 腿在响** |

（最后一行另红一枚 `TestSLO144ReportsThatContradictThemselvesAreCorruptNow`，`:156`，同一发 GO 句子的同族检。）

### 2.3 读数与判定

1. **逃逸是真的**：`m-combo` 用**本程自己的尺**跑出 **0 红**，且 `RUN=20`、`case 11=RAN`
   ⇒ 实现件 §2.4 那一行"两条 leg 全摘 ⇒ 0 红 ⇒ 逃逸"**独立复算成立**，AC#2 的承重依据**不塌**。
   本程**没有**读它的 `combos2/` 日志来得到这个数（那是第 1 格的事）。
2. **而且不只在常量 0 那一发上逃逸**：`m-combo-bytes`／`m-combo-revert` 同样 0 红
   ⇒ 两条腿一摘，那条腿上的取值**换成任何本程试过的错值都不再有人管**（`p1`/`p2`/`p3` 三发全逃逸）。
   实现件只报了 `p1` 那一发的逃逸＋另两发同形状；**本程把它加强成"该腿整体脱保"**。
3. **那颗 0 红不是"什么都没跑"**——这是本格最要紧的一道反向证明。`m-combo-canary` 保持**同样两条腿摘掉**、
   只在 GO 侧再撤一句 `at offset %d`，case 11 **立刻红**（`!ok` 腿那句），
   ⇒ 说明 `m-combo` 那一遍里 case 11 的 fixture、`readSubjectReport` 调用、剩下的腿**全都活着**，
   0 红是**判据射程不到**，不是执行没到。（对照：`m-legs-both` 也是 0 红，但那是"码没坏 ⇒ 本来就该绿"，
   与 `m-combo` 的 0 红**不同因**，靠 canary 那一发分开。）
4. **两条腿各自单独都能接住 `p1`**（`m-zero-fieldoff`／`m-zero-sentenceoff` 各 1 红），
   且**两腿量的不是同一个可观测量**：一腿读**结构体字段** `obs.offset`，另一腿**解析渲染出来的句子**再比。
   ⇒ 实现件 §2.4 那句"少任一条都有兜"**复算成立**，而且它不是一句冗余：两腿在"句子写歪但字段对／字段对但句子没写"
   这两类形状上不等价（本程在 `p6`／`p12` 那一族上没重做，见下"没测什么"）。
5. **期望常数是活的**：`m-wantcanary-live` 只改一枚期望值就红 `1 of 8` ⇒ 本格所有"响"都不是自证式判据的产物
   （第 4 格正面处理那一问）。

**放水两问自答**：① 本程零断言方向改动、零代码改动（第 2 格全部读数走 `-overlay`，快照树里**没有一枚文件被改写**，
本程的尺连写都只往仓外 `myoverlay/` 里写）；② 本程**没有复用实现程那枚 `harness.py`**——尺是本程按锚点字节
现写的（`harness.py` 在本格里只被当作**读物**：第 1 格据它的 `MUT` 定义确认"本程那发与它那发是不是同一味"，
字面量本程自己从 `git show c3f7224:` 里取）。

**第 2 格判定**：**成立**（实现件 AC#2 那枚承重读数经独立重建后为真）。
⚠ 但**同一格里本程把它的自述推翻了一次**（第 1 格 §1.2：0 红是修好尺之后的读数，不是坏尺的产物）——
两件事不冲突：**读数是实的，讲故事的人把方向说反了**。

**本程没测什么（本格）**：
1. **没重做 `p4`–`p9`／`p11`／`p12` 与 `d12`／`d13` 那批发**（只复算了承重的三条腿＋三发取值变异＋两发 canary）。
   ⇒ 谁先被骗：把"第 2 格全对"读成"§2.4 整张 25 行表都被独立复算过"的人——**没有**，本程只独立钉了那枚**承重行**。
2. **没验"两腿在别的形状上是否也互为兜"**：`m-combo` 只证了取值维；句子文本维（`p6`/`p12`）本程没在
   双腿摘掉的组合下跑 ⇒ 那两发是否也整体脱保**未知**。
3. **canary 的覆盖面**：`m-combo-canary` 只证明"`!ok` 腿在那配置下活"，没证明"剩下三条腿全活"
   （`read != len(tc.body)`／`named != read` 两条未单独探针）。
4. **没跑 `-race`、没跑真子进程那枚 case 13 的变异组合**（`m-*` 全部只动 corrupt 支与 case 11 的腿）。

---

## 第 3 格　派单 §2 之(3)：复算完之后答那句——"摘掉任意一味，是否存在一发变异从此打不红？"

**判据**（本仓对"承重"的操作定义，票面 AC#2 原文）：摘掉本票选定的一味 ⇒ 是否存在一发变异从此打不红。
本格只用第 2 格**本程自己**的读数（`mylogs/`），不引实现件 §2.4 那张表。

### 3.1 现量（本程 16 发里与本问相关的 9 发）

```
$ python myharness.py m-branch-syntax-off m-branch-type-off m-branches-both-off m-fallback-to-zero \
                      m-case11-notatest m-zero-case11notest m-combo m-combo-bytes m-combo-revert
m-branch-syntax-off  files=1 rc=1 RUN=20 PASS=12 FAIL=1 case11=RAN   "5 of 8 corrupt shapes"
m-branch-type-off    files=1 rc=1 RUN=20 PASS=12 FAIL=1 case11=RAN   "3 of 8 corrupt shapes"
m-branches-both-off  files=1 rc=1 RUN=20 PASS=12 FAIL=1 case11=RAN   "8 of 8 corrupt shapes"
m-fallback-to-zero   files=1 rc=0 RUN=20 PASS=13 FAIL=0 case11=RAN   reds=-        （grep -cE '---FAIL|first:' -> 0）
m-case11-notatest    files=1 rc=0 RUN=19 PASS=12 FAIL=0 case11=NORUN reds=-
m-zero-case11notest  files=2 rc=0 RUN=19 PASS=12 FAIL=0 case11=NORUN reds=-
m-combo              files=2 rc=0 RUN=20 PASS=13 FAIL=0 case11=RAN   reds=-
m-combo-bytes        files=2 rc=0 RUN=20 PASS=13 FAIL=0 case11=RAN   reds=-
m-combo-revert       files=2 rc=0 RUN=20 PASS=13 FAIL=0 case11=RAN   reds=-
```

### 3.2 答句（分两侧答，因为答案不同）

- **测试码一侧：存在。** 摘掉 **case 11 的两条 leg（合起来算一味）** ⇒ `p1`（常量 0）、`p2`（换成 `obs.bytes`）、
  `p3`（退回 `dec.InputOffset()`）**三发全部打不红**（本程 `m-combo`／`m-combo-bytes`／`m-combo-revert` 三发各 0 红）；
  摘掉**整枚 case 11**（`m-zero-case11notest`）同样 0 红，且 `case11=NORUN`＋`RUN` 从 20 掉到 19
  ⇒ 那一枚用例确实是那三发的**唯一持有者**，不是装饰。
  ⇒ **所以本票选定的"一味"确实承重。** 实现件 §2.4 的方向对，但它只报了 `p1` 一发逃逸，
  **本程把它扩到三发**（该腿**整体**脱保，见第 2 格 §2.3 第 2 条）。
- **生产码一侧：不存在（除一行例外）。** `contradictionOffset` 的两条分支各摘一枚**都当场红**，
  且红的枚数**互不重叠**：摘 `SyntaxError` 支红 `5 of 8`、摘 `UnmarshalTypeError` 支红 `3 of 8`、两支全摘红 `8 of 8`
  （5+3=8 ⇒ 那 8 发里每一发的错值**只来自一条支**，两腿是**划分**而不是冗余）。
  ⇒ 这一侧每一味都被管着：**摘掉必响**。
- **那一枚例外**：`contradictionOffset` 最后的 `return inputOffset` 那一行——本程把它换成 `return 0`
  （overlay 建完**回读证明该替换进了编译字节**，否则尺会 FATAL）⇒ **全组 0 红、`RUN`/`PASS` 与 `m-asis` 逐枚相同**，
  **没有任何外部可见读数变过**。⇒ 判为**装饰腿**，详见第 5 格。

**放水两问自答**：① 本程未动任何断言方向（本程不写码，只造变异）；② 本程的尺是自己写的，
且**比实现程那把多一道 proof**（overlay 落盘后回读、`new in / old not in` 才允许出数）。

**第 3 格判定**：**成立**（承重的答句为真，且本程给出比实现件更宽的逃逸集与一个"摘了也不响"的反例）。

**本程没测什么（本格）**：没试"把两条分支合成一条／调换两分支顺序"这类**等价重构**变异（本程只试"摘掉"与"换值"）；
没试 `contradictionOffset` 的**入参**变异（`dec.InputOffset()` 换成别的表达式）；
`m-fallback-to-zero` 之外没试 `-1`／`int64(obs.bytes)` 等其它替换形状（结论不依赖那一枚，见第 5 格）。

---

## 第 4 格　第二攻击点：case 11 的期望值**是不是自己对自己**

**判据**（派单 §3）：现量核期望常数写在哪、**有没有任何一条路径让同一枚解码器既产 offset 又产期望值**
（那枚就是恒真判据的形状）。

### 4.1 期望常数写在哪（现量）

```
$ git show c3f7224:cmd/wisp/slo_report_144_windows_test.go | sed -n '600,614p'
	cases := []struct {
		name string
		body []byte
		want int
	}{
		{"second-comma-at-26", []byte(`{"mode":"subject-in-tree",,"pass":true}`), 27},
		{"html-first-byte",    []byte("<html>the runner wrote an error page</html>"), 1},
		{"0xff-after-key-at-40", append(append([]byte(nil), doc[:40]...), 0xff), 41},
		{"bad-byte-at-500",    slo149Poke(doc, 500, '@'), 501},
		{"mode-number-at-9",   []byte(`{"mode": 123}`), 12},
		{"array-at-0",         []byte(`[1,2,3]`), 1},
		{"seconds-string-at-49", []byte(`{"mode":"subject-in-tree","seconds":"not a number"}`), 50},
		{"bad-escape-at-10",   []byte(`{"mode":"\q","pass":true}`), 11},
	}
```

⇒ 期望值是**写在测试源里的字面量**，运行时没有任何一行代码把它们喂给解码器或从解码器读回来。
比较式是 `if obs.offset != int64(tc.want)`（`:634`）与 `case offset != tc.want`（`:641`，
`offset` 由 `slo149ContradictionOf(obs.summary())` **解析句子**得到）。

### 4.2 恒真那一问：**运行时不成立**（两发现证）

1. **期望值是活的**：`m-wantcanary-live`（只把 `27` 改成 `999`，其余不动）⇒ **红，且红句是 `1 of 8`**。
   如果判据恒真，改期望值不会有任何反应。
2. **判据能只凭生产码变坏而响**：`m-corrupt-zero`／`m-corrupt-revert` ⇒ 各红 8 发，测试文件一字未动。

⇒ **运行时形状上没有"同一枚解码器既产 offset 又产期望值"的路径。** 这一问派单担心的东西**不存在**。

### 4.3 但**出处**那一问：实现件 §2.1 那句辩护是**假的**

实现件 §2.1 写："case 11 的 8 发每发的期望位置是**本文件自己声明的构造**（坏字节在 index K ⇒ 期望 K+1），
**不是从解码器读回来的数**——这是它和 §1.4 那把"永远绿"的假尺的区别。"
本程把那 8 枚与探针实测的 `eo`（＝解码器自己那枚错误带的位置）逐枚对：

```
$ git show c3f7224:.scratch/wisp/probes/149/probe-corrupt-census.log | grep -oE "P149 [a-z0-9-]+ .*eo=[0-9-]+" | awk '{print $2,$NF}'
html-head eo=1        wrong-type eo=12        trunc40-plus-0xff eo=41   deep-at-500-at eo=501
deep-late-2-0xff eo=1342   bad-escape-in-string eo=11   colon-then-comma eo=26  array-into-struct eo=1
number-into-struct eo=3    string-doc eo=15    brace-first-then-good eo=1  late-type-break eo=50
```

| case 11 用例 | 它的 body ＝ 探针里哪一发 | 测试写的 `want` | 探针实测 `eo` | 相等？ | 它自己的用例名说的 index | `K+1`？ |
|---|---|---|---|---|---|---|
| `second-comma-at-26` | `double-comma-39B` | 27 | 27 | 是 | 26 | 是 |
| `html-first-byte` | `html-head` | 1 | 1 | 是 | 0 | 是 |
| `0xff-after-key-at-40` | `trunc40-plus-0xff` | 41 | 41 | 是 | 40 | 是 |
| `bad-byte-at-500` | `deep-at-500-at` | 501 | 501 | 是 | 500 | 是 |
| `mode-number-at-9` | `wrong-type` | 12 | 12 | 是 | **9** | **否（K+3）** |
| `array-at-0` | `array-into-struct` | 1 | 1 | 是 | 0 | 是 |
| `seconds-string-at-49` | `late-type-break` | 50 | 50 | 是 | 49 | 是 |
| `bad-escape-at-10` | `bad-escape-in-string` | 11 | 11 | 是 | 10 | 是 |

⇒ **8 枚 `want` 与探针实测 `eo` 逐枚相等**（`double-comma` 那发探针的 `want` 列记的是 25，
因为探针那列标的是**第一个**逗号、解码器反对的是**第二个**——两枚表口径不同，本程按 body 对）。
⇒ 所以"不是从解码器读回来的数"这句**在出处意义上不成立**：那 8 枚数**就是**从解码器读数抄进来的。
⇒ 更要紧的是**它给的替代依据自己漏了一发**：`mode-number-at-9` 的名字声明坏值起在 index 9，
`want` 却是 12，直接违反测试文件 `:590` 那行自己写的规则
（"Every want below is stated by how the shape was built: **the bad byte sits at index want-1**"；
这里 `want-1 = 11 ≠ 9`）。⇒ **8 发里 7 发能用那条规则推出来，1 发不能**，
而那 1 发恰好只能由解码器读数解释 ⇒ §2.1 的辩护不是"略有夸张"，是**规则与数据自相矛盾**。

### 4.4 这一格真正要留的两味（不是退回，是更正）

1. **§2.1 的辩护要改口**（**未验证断言**：改文字不动判据，治不到任何一发变异，但治得到"下一位以为
   期望值是独立推导的、于是放心把 case 11 当产品不变量"这一发）。诚实的写法是它自己在 §7.3 已经写过的：
   "K+1 是**这一版 encoding/json 的 scanner 约定**……Go 升级可能改它——那时 case 11 会红"。
   ⇒ **§2.1 与 §7.3 互相矛盾**：一处说"不是从解码器读回来的数"，一处说"这判据钉的就是解码器版本约定"。
   同一枚件里两处对同一件事两种说法，**§7.3 是对的、§2.1 是辩护性的**。
2. **`named != read` 那条腿（`:645`）是刻意地"自己对自己"**——它比的是**同一句话里的两个字节数**（句子内部一致性），
   不比任何外部真值。⇒ 这条腿**结构上永远不可能提供"期望值来自哪里"的证据**，它响的是 `p12` 那一发（句子里两个计数打架）。
   实现件 §2.4 把它和另外三条并列成"case 11 的四条 leg"，**没说这条是自比腿**；
   本格点名，免得下一位拿"四条 leg 都是判据"当"四条都能独立验真"。

**放水两问自答**：① 本程未动断言方向，且**没有**因为"期望值出自解码器读数"就判这条检恒真——
恒真与不恒真是两件事，本程用 `m-wantcanary-live` 与 `m-corrupt-zero` 分开了；
② 本程用的对表尺是**body 相同**这一条（不是用例名相同），因为两枚表的 `want` 列口径不同。

**第 4 格判定**：**附条件入账**。AC#2 那枚检**不是恒真判据**（两发现证），实现件该收的收到了；
但 §2.1 那句"不是从解码器读回来的数"**记为不实**，且与 §7.3 自相矛盾（进 §N-1 第 2 行）。

**本程没测什么（本格）**：
1. **没验"K+1"在别的 Go 版本／别的 stdlib 版本下是否仍是 8 发同值**——本程只证了**当前** blob 上二者相等。
   ⇒ 谁先被骗：升级 toolchain 后把 case 11 的红当产品回归的人（§7.3 已自我登记这一味，本程没推进）。
2. **没逐枚重算 8 发的 index**（`second-comma-at-26` 等 5 发本程手数过、`late-type-break`/`bad-escape` 按 body 对表）。
3. **没查 case 12 的期望值出处**（本格只按派单射程做 case 11）；case 12 的 `offset-assert` 用
   `int64(len(doc))`——那一枚是**构造可独立推导**的，本程没验它有没有同类问题。
4. **没验 `slo144Report(t)` 那枚 fixture 是不是被本票动过**（它跨票共用；本程只核了两枚被测文件的 blob）。

---

## N-1　结论修正记录（推翻实现件／推翻本程派单，都写这里）

| # | 原话（谁说的） | 本程现量 | 结论 |
|---|---|---|---|
| 1 | **实现件 §8 第 6 行**＋**派单 §2**＋**台账 A268②**："第一版 combos 里 `x-p1-d11a-d11b` 的那一发『0 红』是尺坏了的产物"（派单还写明"修好后在 `combos2/` 复算"） | `post/x-p1-d11a-d11b.log`（＝入库的第一遍，时间戳 23:13:38 早于 `combos2/` 23:15:36）顶层 **FAIL＝1**、红句是字段腿那句 `offset is 0, want 27`；`combos2/` 同一发 **FAIL＝0** | **推翻（方向相反）**：坏尺那遍交回的是**1 红**（假安心），**0 红 逃逸是修好尺之后才出现的读数**。三处（实现件／派单／台账）同源于实现程一句自述，本程按日志字节判 ⇒ 自述与它自己点名的凭据文件不符。派单 §2 那句"那一版 `x-p1-d11a-d11b` 的『0 红』是尺坏了"**不成立**，别照它做 |
| 2 | **实现件 §2.1**："case 11 的 8 发每发的期望位置是本文件自己声明的构造…**不是从解码器读回来的数**" | 8 枚 `want` 与探针实测 `eo` **逐枚相等**（第 4 格 §4.3 那张对表）；且 `mode-number-at-9` 的 `want=12` 违反测试文件 `:590` 自己声明的"bad byte sits at index want-1"（名字说 index 9） | **推翻**：那 8 枚数**在出处上就是**解码器读数；给出的"构造推导"替代依据 8 发里只对 7 发成立。实现件**自己的 §7.3** 已写明"K+1 是这一版 encoding/json 的 scanner 约定"⇒ **§2.1 与 §7.3 自相矛盾**，§7.3 对、§2.1 是辩护性文字。⚠ **不影响 AC#2 的牙**：本程另用两发现证该检非恒真（§4.2） |
| 3 | **派单 §2 第 2 条**：`probes/149/post/` 与 `probes/149/combos2/` 里**同名 `x-*` 各有一份** | 现量：`post/` 的 9 枚 `x-*` 全部在 `combos2/` 有对应（为真）；**但派单没点到的还有 4 组**——`post-asis`（`post/`×`post-final/`）、`p1`／`p7`（各三份）、`p10`（两份）⇒ 全树 **13 组同名、28 枚副本** | **说小了**（不是错）：坏尺/多副本的范围比派单点的宽，`x-*` 只是其中一族。本程按整棵树清点（第 1 格 §1.1） |

## N-2　伪授权登记（两个数分开栏）

（本节随格追加，最后在第 9 格汇总定数。）
