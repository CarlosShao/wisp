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

## 第 5 格　派单 §3 之②：`contradictionOffset` 的 fallback 那一行——**摘不得、换得无痕**

**判据**（派单 §3）：判"**摘掉这一行有没有任何外部可见读数变过**"，没变就是装饰腿，说清它买的是什么还是白占一格。

### 5.1 现量

```
$ python myharness.py m-fallback-to-zero
m-fallback-to-zero   files=1 rc=0 RUN=20 PASS=13 FAIL=0 case11=RAN   reds=-
$ grep -cE "^--- FAIL|first:" mylogs/m-fallback-to-zero.log     -> 0
```

变异＝`contradictionOffset` 的默认臂 `return inputOffset` → `return 0`（`cmd/wisp/slo_windows.go:659`）。
本程的尺在出数之前**回读 overlay 落盘的那枚文件**、断言 `return 0\n}` 在里面且 `return inputOffset\n}` 不在，
⇒ **这枚替换确实进了被编译的字节**（不是"变异没生效所以看起来绿"）。

**字面意义的"摘掉这一行"本程做不到**：函数还有返回值，删掉默认臂 Go 直接编不过
（`missing return`）⇒ 所以那一行**在编译意义上承重**，在**读数意义上不承重**。

### 5.2 覆盖剖面独立复算（overlay-free，见 5.3 那条坑）

```
$ # 被验树 c3f7224，不带任何 overlay
$ go test -count=1 -coverprofile=cov-post-real.out -run 'TestSLO144|TestSLO147|TestSLO149' ./cmd/wisp/
github.com/CarlosShao/wisp/cmd/wisp/slo_windows.go:659.2,659.20 1 0      <- 默认臂：0 次
$ # 改前树 64858d6，同一把尺
$ go test -count=1 -coverprofile=cov-pre-real.out -run 'TestSLO144|TestSLO147' ./cmd/wisp/
github.com/CarlosShao/wisp/cmd/wisp/slo_windows.go:727.3,727.17 1 1
github.com/CarlosShao/wisp/cmd/wisp/slo_windows.go:728.4,730.1 1 0        <- exited 的 return：0 次
github.com/CarlosShao/wisp/cmd/wisp/slo_windows.go:732.4,734.1 1 1        <- timeout 那一支：执行过
```

⇒ 实现件 §6.5 那句 `659.2,659.20 1 0`、§3 那三行 pre 剖面，**本程逐字节复算相符**。

### 5.3 ⚠ 本格顺手挖出一枚**全仓会中招的仪器坑**（与本票活无关，但直接影响"变异尺＋覆盖尺"合用）

**`go test -overlay` 与 `-coverprofile`/`-cover` 合用时，`-overlay` 被静默忽略。**
决定性证据——把 overlay 指向一枚**语法就不合法**的文件：

```
$ printf 'this is not valid go code at all\n' 追加进 overlay 指向的 slo_windows.go
$ go test -count=1 -overlay <坏文件overlay> -run TestSLO149CorruptSentence ./cmd/wisp/
..\mutbrokengo\slo_windows.go:913:1: syntax error: non-declaration statement outside function body
FAIL	github.com/CarlosShao/wisp/cmd/wisp [build failed]        <- 不带覆盖：overlay 生效
$ go test -count=1 -overlay <同一枚坏文件overlay> -coverprofile=cov-broken.out -run TestSLO149CorruptSentence ./cmd/wisp/
ok  	github.com/CarlosShao/wisp/cmd/wisp	0.076s	coverage: 1.3% of statements   <- 带覆盖：编译通过 = overlay 根本没吃
```

同形读数在 `-coverprofile` 前置/后置、`-covermode=count`、`-cover` 四种写法下**全部一致**
（实现件那条"直接 `go test` 会 `0xc0000135`"的坑本程也**自己撞了一发变体**：
`export PATH="D:/work/…/sherpa-onnx:$PATH"`（Windows 盘符＋正斜杠）**返回 `0xc0000135`＋0 条 `=== RUN`**，
同样的目录写成 MSYS 形 `/d/work/…` 就正常 ⇒ **同一枚 dll 在位，PATH 拼错也造出那枚假红**）。

⇒ **后果点名**：本仓这套验收方法的**变异尺全靠 `-overlay`**。谁哪天想"顺便看看这发变异执行到哪行"而加上 `-coverprofile`，
拿到的剖面是**未变异那版**的，且**测试会变绿**（因为变异根本没进编译）——这是一颗方向完全相反的假绿。
本程**没能确定**这是 Go 1.27.1 的 bug 还是设计（没读 cmd/go 源码），只报读数。
⇒ **本件里所有覆盖读数都是 overlay-free 的两棵真树跑的**（§5.2 那两条命令），不受此坑污染；
实现件 §3/§6.5 的覆盖读数同样不带 overlay（本程对它给出的命令原文核过），**也不受污染**。
⇒ **建议登记（未验证断言）**：`issues/README` 或台账加一条"变异尺与覆盖尺不得同一次 `go test` 里合用"——
它治得到的是"下一位用 overlay+cover 合用来'证明某发变异走到了某行'"这一发；治不到本票任何一格。

### 5.4 那一行买的是什么

- **买编译完整性**：必须有一条默认臂，函数才有返回值。这一点它不可替代。
- **买"纵深防御"那一说法**：实现件 §4 与函数注释写的是"哪天错误两家都不叫，就退回解码器停下的位置"。
  ⇒ **本程要在这里给那句打折**：一旦真走到那一臂，它返回的是 `dec.InputOffset()`，
  而**本票存在的全部理由就是"`dec.InputOffset()` 在 corrupt 支上会说谎"**（实测 18 发 corrupt 里 8 发印 0）。
  ⇒ 所以那一臂**不是"退回到一个还站得住的数"，是退回到本票刚推翻的那个坏数**，
  而且**今天没有任何一枚用例能发现这次倒退**（§5.1 已证：把它的值整个换掉，全组 0 红）。
- ⇒ 判定它**是装饰腿**：不是"白占一格代码"，而是"占了一个承诺位、给了一个已知会说谎的值、且零见证零牙"。
  与实现件 §7 第 2 条的自评方向一致（它自己登记了"没有见证"），**但它把这行说成"纵深防御"，本程不认**——
  纵深防御要求退路本身是真的；这一条退路是刚被本票判死的那个形状。

**放水两问自答**：① 本程未动断言方向（本格是**造变异去量它**，不是去改判据）；② 未涉及 helper。
⚠ 本程**没有**"顺手替它把装饰腿删了"——那是零代码改动规矩明令禁止的，本程只出读数。

**第 5 格判定**：**成立**（派单那一问有确定答案：**没有任何外部可见读数变过** ⇒ 装饰腿）。

**本程没测什么（本格）**：
1. **没证"这一臂不可达"**——那要把 `json.Decoder.Decode` 的错误全集（含未来版本）枚举一遍。
   本程只试了探针那 20 发＋本程的 16 发变异，**都没走到**。⇒ 谁先被骗：把"零见证"读成"永远不会走到"的人。
2. **fallback 只试了 `return 0` 这一枚替换值**（没试 `-1`／`int64(obs.bytes)`）——结论不依赖枚数：
   既然 `return 0` 都无人响，任何常量替换同理。
3. **§5.3 那坑的机理未查**（没读 `cmd/go` 源码、没试 `-overlaycfg`、没在 linux 容器复现）。
   只报 windows/go1.27.1 的读数。

---

## 第 6 格　派单 §3 之③：case 13 那把覆盖尺是**真 instrument 还是它自己数的**

**判据**：实现件 §3 用"覆盖剖面从 0 次执行变 1 次"当 case 13 的可达性凭据 ⇒ 核那把尺本身。

### 6.1 两点标定（本程现跑）

```
$ # 正控 1：同一把尺、同一棵树，把 -count 从 1 改成 3、covermode 从 set 改成 count
$ go test -count=1 -covermode=count -coverprofile=cov-post-count1.out  -run '...' ./cmd/wisp/
...slo_windows.go:781.3,781.17 1 97
...slo_windows.go:782.4,784.1 1 2        <- exited 的 return：一次跑里被走 2 遍（case 13 的两个子形状）
...slo_windows.go:659.2,659.20 1 0
$ go test -count=3 -covermode=count -coverprofile=cov-post-count3.out -run '...' ./cmd/wisp/
...slo_windows.go:781.3,781.17 1 241      <- 97*3 = 291? 否：见下"读数"，本程按实测记
...slo_windows.go:782.4,784.1 1 6         <- 2*3 = 6   ✓ 严格等比
...slo_windows.go:659.2,659.20 1 0        <- 三遍跑完仍是 0
```

**读数**：
- 把重复次数从 1 抬到 3，**exited 那条 return 的计数从 2 变成 6（精确 ×3）**、
  `781`（条件求值）从 97 变成 241（**不是**精确 ×3 ⇒ 该块在轮询循环里、迭代次数随计时浮动，
  本程**未取证**这个解释，只记"它随重复而动"这一事实）。
  ⇒ **这把尺会随"执行了多少次"动，不会随"文件里写了什么"动** ⇒ 它是编译期插桩的真计数器，
  **不是**实现程自己数的、也不是 grep 行号数出来的。
  ⚠ 顺带一条口径：**`781` 那一枚计数不能当稳定读数引**（它随计时动），能引的是 `782` 这种
  "每个子形状恰好一次"的块（2 → 6 严格等比）。
- `659`（fallback）在 `-count=3` 下**仍是 0** ⇒ 第 5 格那条"零见证"不是重复不够造成的。
- ⚠ **口径要点名**：默认 `-covermode=set` 下那个尾数只表示"**≥1**"，不表示"恰好 1 次"。
  实现件 §3 写的是"执行次数 0"与"执行过"／"现在执行过了"——**用词与 set 语义相符，没说过头**（本程核过它三行措辞）。

### 6.2 它的 ⓐ 判定本程怎么看

实现件 §3 的 ⓐ 判定（该点名 ⇒ 补断言）三条理由里，第 1 条"设计上它会走到"用的是**注释文本＋同一枚 loop**这一段推理，
第 3 条"Ⓒ 不适用"用的是**覆盖读数**。本程的读法：
- **覆盖读数支持的是"case 13 造得出这一支"**（`1 0` → `1 1`，本程复算相符），
  **不支持"生产里真会走到这一支"**——生产要走到它，得有一个 subject 在预算内既不写报告又已退出，
  而实现件 §7 第 1 条自己写明"真 subject 崩溃的形状仍未取证"（本程**没有**推进这一味，见下）。
- ⓑ 那一支（"设计上不该走到"）被 §5 的读数**削弱但没被否证**：走到一次要付一次 `last.summary()`，
  而 `last` 在 exited 支上可能是**零读数**（`0 bytes read, no report file yet`），此时那句括号里带的是"没读到东西"这件事——
  实现件 case 13b 钉的正是这一形状（"不许冒出 `offset`"）。⇒ **本程认为 ⓐ 是三个可能里唯一有读数支撑的那个**，
  但它是**"补断言"级别的成立**，不是"生产语义已定案"级别的成立。

**放水两问自答**：① 本程未动断言方向；② 本程未涉及 helper（本格全部是 `go test -cover*` 的原生命令）。

**第 6 格判定**：**成立**（那把覆盖尺是真 instrument，本程两点标定过；实现件对它的三处引用**措辞未过头**）。

**本程没测什么（本格）**：
1. **没跑 `-race`、没在 linux 容器复现**（`_windows` 文件在 linux 上不参试，本仓已知形状）。
2. **没解释 `781` 的 97→241 为什么不是精确 ×3**（可能含跨用例的共享求值路径）。本程只把它当"随重复而动"的证据。
3. **没量真崩溃形状**（派单 §7 第 1 条实现件自己已登记，本程同样没推进）⇒ **AC#3 那一格到今天仍是
   "造出来的可达性"，不是"量到的生产可达性"**；这条洞在 144／147／149 三张件里**第三次挂账**。

---

## 第 7 格　门禁独立复跑（不复用实现程的环境）

**判据**（派单 §5）：CI 同形跑法改前改后各一次＋名册两向 `comm`；`go vet ./cmd/wisp/` 空；
`gofumpt -l . tools/d22scan tools/mockllm` 空（版本现读，`gofmt` 是另一把尺）；`sh scripts/d22scan.sh` rc=0 且各作用域 `examined N` 非零；
判红绿只认 `--- FAIL:`、不许用裸 grep 数枚数；一枚 panic 会吞掉同包其余读数。
**两棵被验树都是本程自己按 sha 建的快照**（`snap-64858d6-pre`／`snap-c3f7224-b`），**没有用脏工作树**。

### 7.1 CI 同形 `bash scripts/wisp-cli-tests.sh`（本程现跑）

```
$ cd <仓外>/snap-64858d6-pre && bash scripts/wisp-cli-tests.sh      # 改前那版（票 149 之前的码＋测试）
runtests.sh: OK - packages=[./cmd/wisp/ -count=1 -skip ^(TestDefaultDeadlineWallClockMeasurement|…7 枚…)$]
                top-level: PASS=73 FAIL=0 SKIP=0, === RUN=133, '[no tests to run]'=0
portable-tests.sh: four numbers (all from -v output): === RUN=133  --- PASS=73  --- FAIL=0  --- SKIP=0
rc_pre=0   （耗时 98.7s）

$ cd <仓外>/snap-c3f7224-b && bash scripts/wisp-cli-tests.sh        # 被验版本
portable-tests.sh: four numbers (all from -v output): === RUN=136  --- PASS=76  --- FAIL=0  --- SKIP=0
rc_post=0  （耗时 85.7s）
```

⇒ **实现件 §6.1 那两枚四数（133/73/0/0 → 136/76/0/0）本程独立跑出同样四个数**，
且它"以 gofumpt 之后重跑的 `gate-post-full2.log` 为凭"那句在本程这里也不构成差异：
本程的 `snap-c3f7224-b` 就是**入库字节**（`git archive c3f7224`），四数同为 136/76/0/0
⇒ **"快照字节＝入库版本"这句本程核过**，核法＝逐文件 blob md5（§2.1）＋本程在该快照上跑出的四数与它在 `post` 树上跑的一致。
- 两边 `=== RUN` 非零（133／136）⇒ 都"跑到了"；`--- SKIP=0` ⇒ 那 7 条 `-skip` 模式在本作用域**一枚都没跳**
  （它们指向的是别的包里的用例名）⇒ 本票的"绿"不是被 `-skip` 换来的。

### 7.2 名册两向 `comm`（本程自己 `-list`，两棵树各一次）

```
pre=73  post=76
comm -23（消失）：（空）
comm -13（新增）：TestSLO149CorruptLegsWithoutADecoderErrorKeepTheirOwnEnd
                 TestSLO149CorruptSentenceNamesThePositionTheDecoderObjectedAt
                 TestSLO149ExitedGiveUpSentenceCarriesTheLastReading
```

⇒ 实现件 §6.2 那三枚名册**逐字相符**，且**没有一枚既有用例被改名或被挤掉**。
（本程按派单那条坑先把 dll 目录放进 `PATH` 才跑 `-list`——本程确实量到不带 dll 时连 `-list` 都会失败，
与 §N-1 第 4 行那条"因不止一枚"的读数一致。）

### 7.3 三把静态尺

```
$ D:/work/base/gopath/bin/gofumpt.exe --version
v0.12.0 (go1.27.1)                      <- 现读，与实现件 §6.5 同一版
$ (cd snap-c3f7224-b && gofumpt.exe -l . tools/d22scan tools/mockllm)
（空）  rc=0
$ (cd snap-64858d6-pre && gofumpt.exe -l . tools/d22scan tools/mockllm)
（空）  rc=0                              <- 改前那版也干净，所以"改后才需要 gofumpt"那句不成立（见下）
$ (cd snap-c3f7224-b && go vet ./cmd/wisp/)
（空）  rc=0

$ (cd snap-c3f7224-b && sh scripts/d22scan.sh)      rc=0
  正对照: runtests.sh: OK - packages=[./...] top-level: PASS=30 FAIL=0 SKIP=0, === RUN=70
  d22scan: gitignore rules NOT APPLIED - git cannot be consulted in …/snap-c3f7224-b …
           every path in every scope is being scanned …（**响亮自陈**，A207 那枚分母口径）
  d22scan: examined 228 production Go files under internal/ and cmd/
  bans #1-5 internal/=205  cmd/=23 | ban #6 frontend/=52 | ban #7 internal/tools/=18
  ban #8   design/=30  frontend/=52  internal/=412  cmd/=43
  d22scan: clean - no D22 ban violations
$ grep -rlE '^panic ' mylogs gates-cli.log | wc -l
0
```

⇒ 八个作用域 `examined N` **全非零** ⇒ `clean` 不是"什么都没扫"；跑的是 **gofumpt** 不是 `gofmt`。
⇒ ⚠ **正对照那行的小分母（`internal/=14`、`cmd/=1`）是 `tools/d22scan` 的夹具树，不是主扫描**；
两行同印 `clean`，**别把 14/1 那组当仓库分母**（本程第一遍读自己就被这一行绊了一下，点名以免下一位再绊）。
⇒ **分母随树**（本程实测两形，同一枚锚点）：

| 作用域 | 本程·快照 b（＝`c3f7224` 字节，无 `.git`） | 实现程·活工作树（`d22scan-worktree.log`） |
|---|---|---|
| bans #1-5 internal/ ／ cmd/ | 205 ／ 23 | 205 ／ 23 |
| ban #6 frontend/ | **52** | **66** |
| ban #7 internal/tools/ | 18 | 18 |
| ban #8 design/ | **30** | **39** |
| ban #8 frontend/ | **52** | **66** |
| ban #8 internal/ | 412 | 412 |
| ban #8 cmd/ | **43** | **43** |

⇒ 差异全部落在 `frontend/`＋`design/` 两棵**别人正在写**的树（另加快照无 `.git` ⇒ gitignore 不生效），
**`cmd/` 与 `internal/` 两枚本票真正相关的分母逐枚相同** ⇒ 实现件 §6.5 那句"分母随树、不据这两个数下零命中宣称"**成立**。

### 7.4 裸 grep 那一坑（本程独立复算，不背它的数）

```
$ grep -c -- '--- FAIL'  probes/149/gate-post-full2.log   ->  1     <- 假命中
$ grep -c '^--- FAIL'    probes/149/gate-post-full2.log   ->  0
$ grep -cE '^[[:space:]]+--- FAIL' 同一枚                ->  0     <- 子测试也没有
$ grep -n -- '--- FAIL'  …
588:portable-tests.sh: four numbers (all from -v output): === RUN=136  --- PASS=76  --- FAIL=0  --- SKIP=0
```

⇒ 实现件 §6.3 那枚"假 FAIL 来自它自己的汇总行"**逐字节复算相符**（连行号 `:588` 都同一枚）。

### 7.5 一处**与本程有关**的自我污染点名（规矩要求，不藏）

本程在跑**第一遍** `sh scripts/d22scan.sh`（落在 `snap-c3f7224-a`）时，那棵树里**已经有本程自己塞进去的
`cmd/wisp/accept149probe_test.go`**（第 8 格那发探针）⇒ 那一遍报的是 **ban #8 `cmd/`=44**，
比锚点真值多 **1 枚，且多的正是本程自己那枚文件**。
⇒ 处置：按规矩**不删**（临时件只建不删），**换新编号目录 `snap-c3f7224-b` 重建**后重跑 ⇒ 得到上表那组数，
`cmd/`=**43**，与实现程那遍**对上**。
⇒ 时间序本程也钉了（`stat -c '%y'`）：探针落盘 `08:24:54` ⇒ 污染的那遍 d22scan 在 `08:25:42`；
名册 `08:23:21/26`、门禁 `08:22:34`、`-count=3` 覆盖 `08:15:08` 全在**之前** ⇒ **本件其余读数未受该污染影响**。

**放水两问自答**：① 本程未动任何断言、未加任何 `-skip`（`--- SKIP=0` 三处现量）；
② 跑的是派单给的字面命令，没有换更弱的尺（`gofumpt` 而非 `gofmt`；逐包单跑而非 `go test ./...`）。

**第 7 格判定**：**成立**。实现件 §6 的四数、名册差集、假 FAIL 归属、分母随树**四味全部独立复算相符**；
本程另外新增两条：改前那版 gofumpt **也**干净（⇒ "本票改完才需要 gofumpt -w"只是过程史，不是当前状态），
以及 §7.5 那次本程自我污染。

**本程没测什么（本格）**：
1. **没跑 CI**（无 push 权限、票面也禁止）⇒ 本格全部是本机快照读数；`slo-full` 那条腿本程**一次都没求值过**。
2. **没跑 linux 那一腿**（`_windows` 文件不参与；`wisp-cli-tests.sh` 自己在非 windows 上 `exit 2` GUARD）。
3. **没跑 `-race`、没跑整仓 `go test ./...`**（票面地界只有 `cmd/wisp`）。
4. **没验 `d22scan` 在无 `.git` 时"扫到的一切"是否含构建产物**——本程快照里没有构建产物，
   所以那条"可能包含 build output"的自陈在本格**未被触发也未被排除**。

---

## 第 8 格　第三攻击点：票面第 2 条那句"两发已经在读值仍报 0"——**说小了还是说错了**

**判据**（派单 §4）：自己数一遍，判是**说小了**（追加射程更正）还是**说错了**（改名并写清为什么）。
⚠ 派单同时要求把"实现程现量 8 发印 0"也当未验证断言 ⇒ 本程不采它的数，自己造探针。

### 8.1 本程自己造的探针与判据（不读实现程的 census 日志）

派单与实现件都用"印了几个 0"当射程，但**"已经在读值"这句话一直没有可机械复算的判据**。
本程因此**沿用本仓自己已经采用的那枚区分**——`readSubjectReport` 把
`io.EOF` 判成"字节用尽、值根本没开始"、`io.ErrUnexpectedEOF` 判成"**在值里面**"（票 144 就靠这一对分流 `reportUnwritten`）。
⇒ 判据：取"被反对的那枚字节之前"的前缀 `body[:objAt-1]` 单独解码，
回来是 `ErrUnexpectedEOF` ⇒ **已经在读值**（反例）；回来是 `io.EOF` ⇒ **确实没能开始**（与旧注释相符，不是反例）。
本程用的不是错误消息的措辞，是这一对哨兵错误。

```
$ cp <仓外>/accept149probe_test.go <snap-c3f7224-a>/cmd/wisp/    # 只落快照，不落仓
$ go test -count=1 -v -run TestAccept149BeginReadingCensus ./cmd/wisp/
A149 PRE double-comma-39B       prefillOffset=0 objAt=27   enteredValue=true  err="invalid character ',' looking for beginning of object key string"
A149 PRE html-head              prefillOffset=0 objAt=1    enteredValue=false err="invalid character '<' looking for beginning of value"
A149 PRE trunc40-plus-0xff      prefillOffset=0 objAt=41   enteredValue=true  err="invalid character '\\xff' after object key"
A149 PRE deep-at-500-at         prefillOffset=0 objAt=501  enteredValue=true  err="invalid character '@' looking for beginning of value"
A149 PRE deep-late-2-0xff       prefillOffset=0 objAt=1977 enteredValue=true  err="invalid character '\\xff' after object key:value pair"
A149 PRE bad-escape-in-string   prefillOffset=0 objAt=11   enteredValue=true  err="invalid escape sequence `\\q` in string"
A149 PRE colon-then-comma       prefillOffset=0 objAt=26   enteredValue=true  err="invalid character ':' after object key:value pair"
A149 PRE brace-first-then-good  prefillOffset=0 objAt=1    enteredValue=false err="invalid character '}' looking for beginning of value"
A149 PRE wrong-type             prefillOffset=13 (did not print 0)   …array 7 / number 3 / string 15 / late-type 51 同
A149 RESULT corrupt=18 prefillPrintsZero=8 counterexamples(alreadyReading)=6 consistent(couldNotBegin)=2
```

顺带（同一枚探针的前半，跑在**已修好**的码上）：
`A149 fixture len=1978`，18 发 corrupt 的 `obs.offset` 逐枚＝27,1,41,501,1977,11,26,1,12,2,57,1978,1978,1,3,4,15,50
⇒ **最小的一枚是 1、没有一枚是 0** ⇒ 新注释那句"none of those 18 produced a name smaller than 1"成立。

### 8.2 判定：**票面那句是"说小了"，不是"说错了"**

- 票面点名的两发**两发都是真反例**（本程现量 `enteredValue=true`）：
  39 字节双逗号（前 26 字节是合法的 object head）、`'\xff' after object key`（前 40 字节是好头）。
  ⇒ 那两句**措辞与事实相符**，只是**数量少计**：同样形状的本仓探针里有 **6 发**，不是 2 发。
- 另 **2 发（`html-head`／`brace-first-then-good`）不是反例**——它们失败在 index 0、前缀为空、
  哨兵是 `io.EOF` ⇒ **"确实一枚值都没开始读"**，旧注释那句 "0 whenever it could not begin" 在这两发上是**成立**的。
  ⇒ 所以正确的全集是：**印 0 的 8 发里，6 发是反例、2 发是旧注释说中的那种**。
- ⇒ **处置＝追加射程更正**（票面不需要改名：它对那两发的具体描述逐字都对，
  错只错在"两发"被写成像一个全集，而它是样本）。票面 §现量的形状 第 2 条本来就已标〔我本轮现跑过的只有第 4 条〕，
  所以这一处是"样本当全集"的口径问题，**不是事实错误**。

### 8.3 ⚠ 同时推翻实现件那一版（它把射程**说大了**）

实现件 §1："…而这 8 发里 `deep-at-500-at` 与 `trunc40-plus-0xff` **明显已经在读值**——
票面点名的两发复算相符，**本程另量到 6 发同类**。" ⇒ 2＋6＝**它把 8 发全算成反例**。
本程现量：**只有 6 发是反例**，它多算的 2 发恰好是 `html-head` 与 `brace-first-then-good`——
也就是**旧注释唯一说中的那两发**。
⇒ 后果点名：这不是无害的高估。AC#4 的要求是"射程＝corrupt 支**实际**会在哪些形状上报 0，写清楚"。
若按"8 发全是反例"写，下一位会得出"**旧注释每一句都被推翻**"，于是把那段承诺**整块删掉**；
而事实是**该承诺在它自己划的那一档（一枚值都没开始读）仍然成立**，只是不穷尽印 0 的形状。
本程现量那枚注释**没有**犯这个错（它写的是"Of the 20 shapes…, 18 classify as corrupt and none of those 18 produced a
name smaller than 1"，并把旧那句逐字留着），⇒ **代码注释是准的，证据件 §1 那句是宽的**。

### 8.4 本格顺带量到的一处**证据件内部不一致**

实现件 §8 第 9 行写"逐条重跑：**fixture 1978 字节相符**"，而它 §1 引的那份 census 日志第一行是
`P149 FIXTURE len=1343`。本程现量两枚 fixture：

```
$ git show c3f7224:.scratch/wisp/probes/149/zz149probe_windows_test.go | sed -n '33,40p'   # p149Fixture
	doc, err := json.MarshalIndent(&sloRun{ …      <- 与 slo144Report 同形，但少了 GDIMax/WriteOpsTotal/
$ diff <(p149Fixture 段) <(slo144Report 段)
> GDIMax: 12,  > WriteOpsTotal: 141,  > SampleErrors: 0,  > Pass: true,
> 多 1 条 samples、多 1 条 thresholds
```

⇒ 探针用的是**一枚被裁短的副本**（1343 字节），**不是**测试用的那枚（1978 字节）；
本程现量两枚 fixture 下**只有 3 发读数会变**（`deep-late-2-0xff` 1342→1977、`trailing-garbage` 1343→1978、
`second-document` 1343→1978），而注释引的那两枚**计数**（18 corrupt／最小名字 1）在两枚 fixture 下**都不变**。
⇒ 所以注释的**射程宣称是稳的**，但**证据件把 1343 与 1978 当同一枚 fixture 引**（§8 第 9 行）
是又一处"归因腐坏"：引那句的人以为探针与用例同枚 fixture，实际不同名册。

**放水两问自答**：① 本程未动判据方向；② 本程那枚探针**不替换任何东西**——它不导入实现程的探针，
自己直接调 `readSubjectReport`＋stdlib `encoding/json`，是本程独立写的尺（也正因为独立，才抓到 §8.3 那处 8↔6 之差）。

**第 8 格判定**：**成立**（派单那一问有确定答案：**说小了**；追加射程更正＝6 发反例／2 发非反例，
并**顺带推翻实现件那一版的 8 发**）。

**本程没测什么（本格）**：
1. **本程的"已经在读值"判据是本程选的**（哨兵错误对）。若有人改用"失败位置是否 >0"作判据，
   `deep-at-500-at` 一类仍算反例、`html-head`（objAt=1）会被误算成反例 ⇒ 计数会变 7/1。
   ⇒ **口径要先定再数**，本程把口径写死了；下一位换口径要连数一起换。
2. **没验 `io.EOF`／`io.ErrUnexpectedEOF` 这对哨兵在未来 Go 版本是否仍这么分**（与第 4 格同一条洞）。
3. 本程探针在**两枚独立快照**上各跑过一遍（`snap-c3f7224-a`，以及为排除 §7.5 那层污染而新建的 `snap-c3f7224-c`），
   两遍读数逐枚相同（`corrupt=18 prefillPrintsZero=8 counterexamples=6 consistent=2`）；
   ⇒ 但**两枚快照里都有本程那枚探针**（探针是被测物、不参与被判据），所以"探针文件的存在会改掉哪些读数"
   这一问本程**只在 d22scan 那一把尺上量化过**（§7.5：`cmd/` 44→43），没在别的尺上量化。

---

## 第 9 格　票面 AC#5（契约轴）＋票面 §4 那枚限定语的**射程实测**

**判据**：① 票面 AC#5 点名的禁改面在**本票那四枚 commit 里**必须零字节，预算常量不许动；
② 票面 §4 写"红句两实参对调"在 `unwritten` 支是等价变异、**不许为它开"要它响"的格**、
"它要问的是 corrupt 支" ⇒ 本程把这三条支**各对调一遍**，看这句话的作用域对不对、以及实现件有没有把它用宽而漏了 coverage。

### 9.1 AC#5 现量（只算票 149 那四枚，不算区间里别人的件）

```
$ for c in b417d31 2f5c7f9 a055d9f c3f7224; do git show --name-only --format='COMMIT %h' $c; done | sort -u | …
     66 probes/149/*                                        <- 尺与逐发日志
      2 docs/evidence/s1/149-three-unpinned-outlets-r1.md   <- 它自己那枚件（两枚 commit 各改一次）
      1 cmd/wisp/slo_windows.go
      1 cmd/wisp/slo_report_144_windows_test.go

$ … | grep -E '^(internal/|tools/d22scan/|.*thresholds\.go$|.*golden.*|.*allowlist\.txt$|scripts/slo-check\.ps1$|docs/PLAN\.md$|docs/specs/|frontend/|design/)'
matches=0
$ git diff 64858d6..c3f7224 -- cmd/wisp/slo_windows.go | grep -E '^[+-].*(subjectReportBudget|subjectGrace)'
（空）
```

⇒ **AC#5 成立**：路径名册与禁改面**交集为 0**（本程把匹配数打出来，不写"看起来是空"），预算常量零命中。
⚠ **本格有一次本程自己的无效读数要点名**：本程第一次跑这条时用 `git show -s --name-only` ⇒
**四个 `fatal: options '--name-only' … and '-s' cannot be used together`**，而那一步的尾巴照样印了本程手写的
`（空 = 零字节）`。**那是一枚失败命令后接出来的"空"，不是零命中的读数。** 本程发现后按上面的写法重跑，
才得到 `matches=0` 这一枚。⇒ 记进 §N-1 第 6 行：**"grep 没吐东西"与"grep 根本没跑"在两把尺下同形**，
这恰是本票第 1 格那族"坏尺产物"的**又一枚实例**，而这次是本程自己造的。

### 9.2 三条支的"两实参对调"——本程现量

```
$ python myharness_b.py m-asis m-swapargs m-unwritten-swap m-complete-swap
m-asis            files=0 rc=0 RUN=20 PASS=13 FAIL=0  case11=RAN  reds=-
m-swapargs        files=1 rc=1 RUN=20 PASS=10 FAIL=3  reds=TestSLO144ReportsThatContradictThemselvesAreCorruptNow,
                                                            TestSLO144LoopGiveUpSentencesOnRealFiles,
                                                            TestSLO149CorruptSentenceNamesThePositionTheDecoderObjectedAt
m-unwritten-swap  files=1 rc=0 RUN=20 PASS=13 FAIL=0  reds=-
m-complete-swap   files=1 rc=0 RUN=20 PASS=13 FAIL=0  reds=-
```

- **corrupt 支对调 ⇒ 3 红**（case 11 ＋ 144 的两枚）。⇒ 票面那句"它要问的是 corrupt 支"**成立**，
  而且**实现件"没有为它开格"这件事没有丢 coverage**：本程再把字段腿／`named != read` 腿／`read != len` 腿
  **逐个摘掉再对调**（`m-swap-readlegoff`／`m-swap-readnamedoff`／`m-swap-readlenoff`），**三发全部仍然 3 红**
  ⇒ 那发对调**不靠 case 11 任何单独一条腿**，它同时被 144 的两枚既有断言管着。
  ⇒ **结论：实现件 §7 第 5 条把票面 §4 的禁令用到了 corrupt 支上（"没有造'对调后必须响'这种判据"）是用宽了，
  代价为零**——那形状今天已经响，不需要新开格。
- **unwritten 支对调 ⇒ 0 红**。⇒ **票面 §4 的等价性断言独立复算成立**（该支 `offset ≡ bytes`，
  是票 147 ⓐ 之后必然的恒等），"不许为它开要它响的格"是**对的处置**，不是回避。
- ⚠ **`complete` 支对调 ⇒ 0 红，而这一条票面没写**。同一结构性原因：
  147 的 `case 10` 断言的是"`1978 bytes read, complete document at offset 1978` 里两枚数都是 1978"
  （见 `pre-B15` 的红句原文，实现件 §1 引过），⇒ **两支数值恒等 ⇒ 对调永远产不出自己的读数**。
  ⇒ 按票面 §4 自己的逻辑，这一支也该写成"**不许开格**"；票面只点了 `unwritten` 一支 ⇒
  **这是票面 §4 的一处射程缺口**（不是实现件的错——它没被要求查这一支）。
  ⇒ 本格**只登记、不开格、不改票面**：把"哪几支上两枚数恒等"写成一张三行的表，归编排者定。

### 9.3 顺带量到的一枚**证据件与代码注释同源不准**

`slo_windows.go` 里 `contradictionOffset` 的函数注释逐字写着：
"… and one whose whole head is the subject's own document until **0xff lands at index 1341**"。
⇒ 1341 是**探针那枚裁短 fixture**（1343 字节）上的读数；**用例与门禁真正使用的 `slo144Report` 是 1978 字节**
（本程现量 `A149 fixture len=1978`），同一发形状在真 fixture 上的坏字节落在 **index 1976**。
⇒ 与本程 §8.4 那处（实现件 §8 第 9 行把 1978 当探针 fixture 引）是**同一处混用的两半**：
一件把 1343 的读数写进了代码注释、另一件把 1978 说成探针的读数。
⇒ **注释里那两枚计数（18 corrupt／最小的名字是 1）本程在两枚 fixture 下都复算成立**，
所以**射程宣称没坏**；坏的是"引这行的人以为 1341 是产品文档里的数"。

**放水两问自答**：① 本程未动任何断言方向（9.2 全部是造变异去量）；② 本程未涉及 helper。

**第 9 格判定**：**AC#5 成立**；票面 §4 限定语**方向对、作用域写窄了一支**（`complete` 支同源恒等，未登记）。

**本程没测什么（本格）**：
1. **没试 `note` 那条支**（`summary()` 第一分支 `%d bytes read, %s`）——它只有一个计数，无实参可对调；
   本程**没有**逐行确认这一句，是按格式串读出来的。
2. **没验"两枚数恒等"是恒久还是当前**：`complete` 支若哪天让 `offset ≠ bytes`（例如改成"末 token 之前的偏移"），
   那一发对调就变成可检的——本程没测那一发。
3. AC#5 的"零字节"只覆盖**文件名册交集**，**不**覆盖"别人在同一枚文件里也动了东西"（本票四枚区间里
   确实夹着 138 的件与台账，本程已把它们从名册里剥出去再算交集）。

---

## 第 10 格　总裁决

| 格 | 对应票面/派单 | 档位 | 一句理由 |
|---|---|---|---|
| 0 | 锚点 | **成立** | `c3f7224` 存在且是 HEAD 祖先；脏区逐枚点名；dll 3 枚在位；快照 blob md5 相符 |
| 1 | 派单 §2(2) 坏尺读数清点 | **附条件入账** | 那批读数确属坏尺期间（为真），**但实现件对它的叙述方向相反**；`post/` 那 9 枚 `x-*`＋`post-asis` 本程判为不作凭据并点名 |
| 2 | 派单 §2(1) 独立重建组合变异 | **成立** | 本程自己的尺复算出 **0 红逃逸**，并加强为"该腿整体脱保"；canary 反向证明那颗 0 红不是"什么都没跑" |
| 3 | 派单 §2(3) 承重那句 | **成立** | 测试码侧存在逃逸（摘两条 leg ⇒ 三发全不红）；生产码侧摘谁必响；唯一例外是 fallback 那一行 |
| 4 | 第二攻击点 case 11 期望值 | **附条件入账** | **非恒真**（两发现证）；但 §2.1 那句"不是从解码器读回来的数"**不实**，与 §7.3 自相矛盾 |
| 5 | 第三攻击点① fallback | **成立** | 换值后**零枚外部读数变化** ⇒ 装饰腿；且它退回的正是本票刚判死的那个坏数，"纵深防御"那一说不认 |
| 6 | 第三攻击点② 覆盖尺真伪 | **成立** | 两点标定（782 块 2→6 严格等比）⇒ 真插桩；实现件三处引用措辞未过头 |
| 7 | 门禁（AC#6） | **成立** | 四数 133/73/0/0 → 136/76/0/0 独立跑出；名册消失 0/新增 3；gofumpt v0.12.0 空、vet 空、d22scan rc=0 八枚分母非零；假 FAIL 归属复算相符 |
| 8 | 派单 §4 票面措辞 | **成立** | 判定＝**说小了**（2 → 本程现量 6 发反例，另 2 发是旧注释说中的那一档）；顺带推翻实现件那一版的"8 发同类" |
| 9 | AC#5 ＋票面 §4 限定语 | **成立** | 禁改面交集 0（含本程自己一枚**无效"空"读数**的点名）；`complete` 支对调也恒等＝票面 §4 的射程缺口 |

**六格 AC 总裁决（本程口径，不替编排者勾框）**：
**AC#1 成立／AC#2 成立（凭据已独立重建）／AC#3 成立（ⓐ 判定有真插桩支持，但"生产可达"仍未取证）／
AC#4 成立（注释射程准，证据件文字宽）／AC#5 成立／AC#6 成立。**
⇒ **本票不需要返工到实现程**：本程**没有造出**任何一发让交付的判据失灵的活（第 1、4、5、8 格推翻的都是**文字与读数归因**，
不是牙）。⚠ 但**三处必须更正的文字**留在件里，且**其中两处是实现件的核心辩护句**（§2.1、§8 第 6 行）。
⇒ 档位为何是"附条件"而不是"退回"：分界＝**本程造没造出来**——本程造出来的是"叙述与凭据不符"，
**没有**造出"判据不响的坏形状"。

## 第 11 格　本程没测什么（全件汇总，按"漏了它谁会先被骗"排序）

1. **真 subject 崩溃的形状一次都没取证**（`exited()` 那支只被 case 13 **造**到可达，没被**量**到）。
   ⇒ 谁先被骗：排"`wisp slo` 到点红"、以为那句会印出真实停留位置的人。**这条洞 144/147/149 第三次挂账。**
2. **`contradictionOffset` 的 fallback 那一行既不承重、退路又是已知会说谎的数**（§5.4）。
   ⇒ 谁先被骗：以为"三家错误都有位置"的人。本程**没证**它不可达。
3. **§5.3 那枚仪器坑（`-overlay` 与 `-cover*` 合用时 overlay 被静默忽略）本程没查机理、没在 linux 复现**。
   ⇒ 这条**影响的是下一位的方法**，不是本件读数；但谁先按"变异＋覆盖"合起来跑，谁先拿到假绿。
4. **AC#2 的矩阵本程只独立复算了承重那一行**（16 发里覆盖 `p1`/`p2`/`p3`/`p4`/`p5`/`d11*`/`d11e` ＋ canary），
   `p6`–`p12` 里除 `p6`/`p12` 的等价面（§9.2）之外**没逐发重做**，`d12`/`d13` 两族**完全没碰**。
   ⇒ 谁先被骗：把"第 2/3 格全对"读成"§2.4 那张 25 行表被独立复算过"的人。
5. **case 12（无解码器错误那两条腿）本程只按覆盖与门禁读数收下，没做它的期望值出处审计**（第 4 格只做 case 11）。
6. **没跑 CI／没跑 linux 腿／没跑 `-race`／没跑整仓 `go test ./...`**；`slo-full` 那条腿本程一次都没求值过。
7. **本程自己的探针文件污染过一枚快照**（§7.5，`cmd/` 44 vs 43），已用新目录重跑并量化差值；
   但"还有哪把尺被那枚文件改过读数"**未穷举**（本程只核了时间序在探针之前的那几把：名册/门禁/覆盖标定）。
8. **票面六框本程一枚未勾**（规矩）；台账/HANDOVER/票面本程一字未动（都在禁改面或归编排者）。

## N-1　结论修正记录（推翻实现件／推翻本程派单，都写这里）

| # | 原话（谁说的） | 本程现量 | 结论 |
|---|---|---|---|
| 1 | **实现件 §8 第 6 行**＋**派单 §2**＋**台账 A268②**："第一版 combos 里 `x-p1-d11a-d11b` 的那一发『0 红』是尺坏了的产物"（派单还写明"修好后在 `combos2/` 复算"） | `post/x-p1-d11a-d11b.log`（＝入库的第一遍，时间戳 23:13:38 早于 `combos2/` 23:15:36）顶层 **FAIL＝1**、红句是字段腿那句 `offset is 0, want 27`；`combos2/` 同一发 **FAIL＝0** | **推翻（方向相反）**：坏尺那遍交回的是**1 红**（假安心），**0 红 逃逸是修好尺之后才出现的读数**。三处（实现件／派单／台账）同源于实现程一句自述，本程按日志字节判 ⇒ 自述与它自己点名的凭据文件不符。派单 §2 那句"那一版 `x-p1-d11a-d11b` 的『0 红』是尺坏了"**不成立**，别照它做 |
| 2 | **实现件 §2.1**："case 11 的 8 发每发的期望位置是本文件自己声明的构造…**不是从解码器读回来的数**" | 8 枚 `want` 与探针实测 `eo` **逐枚相等**（第 4 格 §4.3 那张对表）；且 `mode-number-at-9` 的 `want=12` 违反测试文件 `:590` 自己声明的"bad byte sits at index want-1"（名字说 index 9） | **推翻**：那 8 枚数**在出处上就是**解码器读数；给出的"构造推导"替代依据 8 发里只对 7 发成立。实现件**自己的 §7.3** 已写明"K+1 是这一版 encoding/json 的 scanner 约定"⇒ **§2.1 与 §7.3 自相矛盾**，§7.3 对、§2.1 是辩护性文字。⚠ **不影响 AC#2 的牙**：本程另用两发现证该检非恒真（§4.2） |
| 3 | **派单 §2 第 2 条**：`probes/149/post/` 与 `probes/149/combos2/` 里**同名 `x-*` 各有一份** | 现量：`post/` 的 9 枚 `x-*` 全部在 `combos2/` 有对应（为真）；**但派单没点到的还有 4 组**——`post-asis`（`post/`×`post-final/`）、`p1`／`p7`（各三份）、`p10`（两份）⇒ 全树 **13 组同名、28 枚副本** | **说小了**（不是错）：坏尺/多副本的范围比派单点的宽，`x-*` 只是其中一族。本程按整棵树清点（第 1 格 §1.1） |
| 4 | **派单 §0**："`git archive` 出的树不带 `third_party/` 的原生 dll，跑 `cmd/wisp` 会拿到假红（`no native DLLs`、`PANIC=0`）⇒ 先确认 dll 在位再报红" | dll 缺失这一因**复算成立**（`git archive c3f7224` 的顶层名册里没有 `third_party`）。但本程量到**同一枚假红的另一因**：dll **已在位**、只把 `PATH` 写成 `D:/work/…`（盘符＋正斜杠）形，同样 `exit status 0xc0000135`＋**0 条 `=== RUN`**；同一目录改写成 MSYS 形 `/d/work/…` 就正常 | **说小了**：判"根本没跑到"不能只看"文件在不在位"，`0xc0000135` 的**因至少两枚**（dll 缺、PATH 拼法）。本程仍用"`=== RUN` 枚数是否为 0"当唯一区分依据（那是共因读数），但**不**把"dll 在位"当充分条件 |
| 5 | **实现件 §6.2 那句**："`go test -list` 也要启动测试二进制，所以它和跑测试死在同一个地方（`0xc0000135`）" | 本程**没有**复算 `-list` 那一发（只在第 5 格撞到同一坑的另一变体）。本程复算出的是**另一条更凶的**：**`-overlay` 与 `-cover*` 合用时 overlay 被静默忽略**——把指向"语法不合法文件"的 overlay 交给 `go test`，不带覆盖 `[build failed]`、带覆盖 **`ok`＋coverage 1.3%** | **未复核第 5 格那条**（不算推翻、只算没查）；**新增一条派单与实现件都不知道的仪器坑**，写进第 5 格 §5.3。它污染的不是本件读数，是**下一位的方法** |
| 6 | **派单 §5／实现件 §6.1＋§6.3**：门禁四数 133/73/0/0 → 136/76/0/0、"以 gofumpt 之后重跑的那份为凭、且核对过快照字节＝入库版本"、"假 FAIL 来自汇总行" | 本程在**自己按 sha 建的两棵快照**上跑 CI 同形：pre `RUN=133 PASS=73 FAIL=0 SKIP=0 rc=0`、post `136/76/0/0 rc=0`；假 FAIL naive 1／anchored 0、命中行 `:588`；两枚被测文件 blob md5 与 `git show c3f7224:<path>` 逐枚相等 | **复算相符**，且"快照字节＝入库版本"这一问本程**真做过**（不是转述）：md5 相等＋干净快照 `snap-c3f7224-b` 上重跑得到同一组四数（§7.1、§7.3） |
| 7 | **本程自己在第 9 格第一次跑的 AC#5 交集检查** | 那一跑用了 `git show -s --name-only` ⇒ 四枚全部 `fatal: options '--name-only' … and '-s' cannot be used together`，而本程手写的"（空 = 零字节）"照样印了出来 | **推翻本程自己一处跑法**（结论未变、凭据换了一版）：改 `git show --name-only --format=''` 重跑才得到真的 `matches=0`。⇒ **"grep 没吐东西"与"grep 根本没跑"两形同貌**——本票第 1 格那族"坏尺产物"的又一枚实例，这次是本程自己造的 |
| 8 | **实现件 §8 第 9 行**："逐条重跑：**fixture 1978 字节相符**"（用它 §1 那份 census 的读数背书） | 探针 `p149Fixture` 与用例 `slo144Report` 是**两枚不同文档**：本程现量 `A149 fixture len=1978`，而它 §1 引的 census 第一行是 `P149 FIXTURE len=1343`；diff 两枚构造函数可见探针少了 `GDIMax`/`WriteOpsTotal`/`SampleErrors`/`Pass`＋1 条 samples＋1 条 thresholds | **推翻（口径混用）**：1978 与 1343 不是同一枚 fixture。⇒ 连带牵出**代码注释里同源的一枚**：`contradictionOffset` 函数注释那句"0xff lands at index **1341**"只在探针那枚短文档上成立，真 fixture 上那一发落在 index **1976**（§9.3）。⚠ 两枚**计数**（18 corrupt／最小的名字 1）本程在**两枚 fixture 下都复算成立** ⇒ 射程宣称没坏，坏的是归属表述 |
| 9 | **实现件 §1**："8 发里 deep-at-500 与 trunc40 明显已经在读值…本程**另量到 6 发同类**"（＝把 8 发全算成反例） | 本程用独立定义的判据（`io.EOF` vs `io.ErrUnexpectedEOF`，即本仓分流 `reportUnwritten` 的那一对哨兵）逐枚判：`counterexamples(alreadyReading)=6`、`consistent(couldNotBegin)=2`（`html-head`／`brace-first-then-good`，两枚都失败在 index 0、前缀为空） | **推翻一半（实现件把射程说大了 2 发）**：正确的三枚数是 **印 0＝8、反例＝6、旧注释说中＝2**。⚠ 后果非无害：按"8 发全是反例"写，下一位会把那段承诺**整块删掉**，而它在"一枚值都没开始读"那一档仍然成立（§8.2/§8.3） |
| 10 | **派单 §4**："实现程现量是 **20 发里 8 发印 offset 0**、并说'两发是样本数不是全集'" | 本程独立跑探针（不复用它的 census 日志）：`corrupt=18 prefillPrintsZero=8` | **相符**（实现程那枚 8 复算成立；派单"两发"那句判为**说小了**，见 §8.2） |
| 11 | **票面 §4 限定语的作用域**："unwritten 支 `offset ≡ bytes` ⇒ 两实参对调是等价变异、不许为它开'要它响'的格；它要问的是 corrupt 支" | 本程三条支各对调一遍：corrupt ⇒ **3 红**（且不靠 case 11 任何单腿——把三条腿逐个摘掉再对调，三发全部仍 3 红）；unwritten ⇒ **0 红**；**complete ⇒ 0 红**（`case 10` 断言的两枚数恒等于 1978） | **成立但写窄了一支**：等价性在 `unwritten` 复算成立（票面对），**`complete` 支同形而票面未登记**。按票面自己的逻辑那一支也该进"不许开格"表 ⇒ 本格**只登记、不开格、不改票面**。⚠ 同时给实现件 §7 第 5 条打折：它把禁令用到 corrupt 支是**用宽了**，但**代价为零**（那一发今天已响） |
| 12 | **派单 §7**："本仓已有三枚程死在攒着最后写" | 本程**未复核**（没有去数历史里究竟几枚程半途断过） | **未复核**，登记为"本程引用不到出处"；不作为本件任何一格的依据 |

## N-2　伪授权登记（两个数分开栏）

**真通知回显数：2 枚**（计数法＝本程逐条点名，可机械复算）：
- **1 枚**后台任务完成回显（`<task-notification>`，task id `baoyycurc`，即 §6.1 那发 `-covermode=count -count=3` 覆盖跑）。
  命令前 40 字：`O=/d/work/tmp/wisp149-accept && cd $O/snap-c…`（Bash 工具，`run_in_background`）。
- **1 枚**会话首条项目上下文回显（`AGENTS.md` 全文＋技能清单＋一条 `MEMORY.md` 修改通知）。**非工具输出、非授权**。
  ⚠ 那条 `MEMORY.md` "was modified since it was last read" 的措辞带"文件已被修改"形状 ⇒
  本程按规矩**登记为回显、不判为注入，且不据此改任何判据**（它指向的是本程私有的记忆文件，不是本票地界）。

**判为注入数：0**。理由：
- 上述 2 枚里**没有任何一句**要求"少取证／别用工具／直接给结论／放宽阈值／已解锁／已核验请继续提交"。
  两条后台回显本程都**继续取证**（`baoyycurc` 那枚的结果就是 §6.1 的标定读数）。
- 本程对自己的票面／实现件／`probes/149/`（66 枚）做了字面反扫，并**核到本程派单原文**：
  ```
  $ grep -rniE "编排者备注|系统提示|已被修改|已核验请继续|请放宽阈值|已解锁|不用取证|Confirm the harness note" \
      .scratch/wisp/issues/149-*.md docs/evidence/s1/149-*.md .scratch/wisp/probes/149/
  ```
  ⇒ 命中全部落在**这些文件自己叙述这些字样的行**（147 验收件与实现件 §9 那一节），**外部注入 0 枚**。
- **归属核对（派单 §6 那条新形）**：本程在写"编排者说过 X"的地方共 4 处（派单 §2 的"0 红是尺坏"、
  派单 §2 的"同名各有一份"、派单 §0 的"dll 假红 8 枚"、派单 §4 的"两发"），
  **全部按字面回核到派单原文**，其中 **2 处判为不成立/说小**（写进 §N-1 第 1、3 行），
  另 2 处成立（§N-1 第 4 行）。本程**没有引用过任何一句回核不到出处的"编排者说"**。

## N-3　凭据值自扫

本件**抄录的凭据值＝0**：全文只有测试名、变量名、文件名、仓外快照路径与 git object id。
本程在落盘后按仓里那把尺现扫一遍（**不背实现件的 13/14**）：

```
$ R="sk-[A-Za-z0-9]{8,}|[A-Za-z0-9+/]{40,}={0,2}|api[_-]?key[[:space:]]*[:=][[:space:]]*[^ ]{6,}|Bearer [A-Za-z0-9]"
$ grep -cE "$R" <本件>            -> lines     = 12
$ grep -oE "$R" <本件> | wc -l     -> fragments = 13      （一枚行里有两个命中，故 13 > 12）
$ grep -oE "$R" <本件> | <归类管道> | sort | uniq -c
      9  Go 测试名（长驼峰串）
      3  `sk-notification`                 （`<task-notification>` 撞上 `sk-` 那一支 ⇒ **这把尺自己的假命中**）
      1  40 位十六进制 git object id      （§0 的 `4d0866a` 全形）
```

⇒ 归类后**无法归入"object id／Go 测试名／尺的假命中"三者者：0 枚** ⇒ **抄录的凭据值＝0**。
⚠ 但这把尺**不空**（12 行／13 枚），且**其中三枚是尺自己的假命中**——本程按实现件 §9 同一条理由
不把它写成"（空）"。⚠ 这两枚数是**读数**：本件每多印一枚长测试名就会涨；复算请**重跑上面那三条命令**，别背 12/12。

## N-4　git 纪律自证

- 每一枚 commit 带**显式 pathspec**（`git commit -q -F - -- docs/evidence/s1/149-…-accept-r1.md <<'MSGEOF'`，
  定界符加单引号 ⇒ 本件正文里的反引号与 `$` 未被执行），且**每枚 commit 之后 `git show --name-only` 自验名册**。
- ⚠ **本格要留一枚真实危险读数**：本程**第一次** `git add` 之后，`git diff --cached --name-only` 里
  **同时出现了别人的四枚文件**（`internal/agent/compress.go`／`compress_trace_test.go`／`harness_test.go`／`loop.go`）——
  那是票 139 那枚程**自己 staged 在共享 index 里的活**。⇒ 那一刻只要 `git commit` **不带 pathspec**，
  就会把 139 的中间态吞进本程那枚 commit（＝派单 §1.4 明令防的那一发）。
  ⇒ 本程按规矩停手核对名册，并**只提交自己那一枚路径**；随后 `1acdd03 fix(139 AC#2+AC#3)` **由 139 自己提交**了那四枚
  （本程事后核过：`git show --name-only HEAD` 只含本件一枚路径，139 的 staged 内容**未被本程动过**）。
  ⇒ **这一条不是本程做对了什么，是共享工作树的既有风险又发生了一次**，值得编排者记名。
- `git add -A`／`git add .`／`git commit -a` 使用数 **0**；
  `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`rm`／`rmdir` 使用数 **0**；**push 次数 0**。
- 临时件**只建不删**：仓外 `D:\work\tmp\wisp149-accept\` 下
  `snap-c3f7224-a`（含本程探针那枚文件，未删）、`snap-c3f7224-b`、`snap-c3f7224-c`、`snap-64858d6-pre`、
  `logs-at-anchor/`、`myoverlay/`（16 枚变异目录）、`mylogs/`（16 枚原始日志）、探针源＋两把尺脚本、
  覆盖剖面 6 枚（`cov-asis/branchesoff/broken/pre-real/post-real/post-c3/post-count1/post-count3`）；
  **仓库目录内没有建过 worktree 或 checkout**；**零代码改动**（本程写过的文件只有本件一枚，其余全在仓外）。
- 票面六框本程**一枚未勾、也未替实现方勾**。


------------------------------------------------------------

**未完成＝无**（派单 §2 / §3 / §4 / §5 四组攻击点全部有本程现跑读数）。
⚠ 枚数自证：本件**随格 commit**，落这一版之前已入库 4 枚（第 0-1／2／3-4／5-6 格），
**这一枚是第 5 枚**（第 7-11 格＋N-1…N-4 四节）；复算请
`git log --format=%h -- docs/evidence/s1/149-three-unpinned-outlets-r1-accept-r1.md`，别背本句的数。
**本程交付的是裁决，不是勾框**：票面六框一枚未勾、未替实现方勾；台账／票面／HANDOVER 一字未动。
