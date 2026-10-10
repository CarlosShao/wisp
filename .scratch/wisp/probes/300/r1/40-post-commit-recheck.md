# 300-r1 — `40` 提交后复核（正文那笔＝`028529fc`，⛔ amend、⛔ push）

派单那条新规矩（`A804`）我这把跑了两遍：**每笔 commit 前** `git diff --cached --name-only` 数名册（逐字结果＝`logs/rosters.txt`，
含"名册里出现非我射程路径？"那一行判据＝空），**每笔 commit 后** `git show --stat` 再数一遍（下面 §1）。

## 1. commit #1 名册（`git show --stat HEAD` 的 stdout，逐字整块）

```
[dev 028529fc] fix(audio): read WAVEFORMATEXTENSIBLE SubFormat at byte 24, not 26 (ticket 300 AC#2)
 26 files changed, 10514 insertions(+), 3 deletions(-)
 ...bformat-offset-reads-two-bytes-past-the-guid.md |    1 +
 .scratch/wisp/probes/300/r1/10-ac2-change.md       |  130 ++
 .scratch/wisp/probes/300/r1/20-ac2-rig.md          |  142 ++
 .scratch/wisp/probes/300/r1/30-gates.md            |  170 ++
 .scratch/wisp/probes/300/r1/90-unrun-rulers.md     |   65 +
 .scratch/wisp/probes/300/r1/logs/change.diff       |   39 +
 .scratch/wisp/probes/300/r1/logs/d22scan-final.txt |  251 +++
 .scratch/wisp/probes/300/r1/logs/d22scan.txt       |  251 +++
 .../probes/300/r1/logs/diff-vs-a1r-nocomment.txt   |   26 +
 .../wisp/probes/300/r1/logs/diff-vs-a1r-rig.txt    |   87 +
 .scratch/wisp/probes/300/r1/logs/entry-1.txt       |   12 +
 .scratch/wisp/probes/300/r1/logs/entry-2.txt       |   12 +
 .scratch/wisp/probes/300/r1/logs/format-rulers.txt |   19 +
 .scratch/wisp/probes/300/r1/logs/postB-1.txt       | 2217 ++++++++++++++++++++
 .scratch/wisp/probes/300/r1/logs/postB-2.txt       | 2217 ++++++++++++++++++++
 .../wisp/probes/300/r1/logs/postB-entry-final.txt  |   12 +
 .scratch/wisp/probes/300/r1/logs/preA-1.txt        |  109 +
 .scratch/wisp/probes/300/r1/logs/preA-2.txt        |  109 +
 .scratch/wisp/probes/300/r1/logs/preA-3.txt        | 2207 +++++++++++++++++++
 .scratch/wisp/probes/300/r1/logs/preA-4.txt        | 2207 +++++++++++++++++++
 .../probes/300/r1/logs/progress-line-clock.txt     |    1 +
 .scratch/wisp/probes/300/r1/logs/rig-writes.txt    |   23 +
 .scratch/wisp/probes/300/r1/logs/rosters.txt       |   32 +
 .scratch/wisp/probes/300/r1/logs/vet.txt           |    0
 .../audio/parse_wave_format_300_windows_test.go    |  157 ++
 internal/audio/wasapi_windows.go                   |   21 +-
```

- ⚠ `git show --stat` 会把长路径**中段折成 `...`**（上面三行 `.../probes/300/r1/logs/…`、`...bformat-…`、`.../audio/parse_wave_format_300_windows_test.go`）——
  那是 `--stat` 自己的排版，⛔ 我在抄件里补全。全名尺＝`git show --name-only --format='' 028529fc`
  （这一把是**提交之后**补跑、**追加**进 `logs/rosters.txt` 的第二段；同文件第一段那把是提交之前的
  `git diff --cached --name-only`。两把名册同为 **26** 枚，与 commit 那行 `26 files changed` 同值；
  判据"名册里出现非我射程路径"两把都是空，逐字都在 `logs/rosters.txt`）。
- ★**产码名册只有两枚**：`internal/audio/wasapi_windows.go`（`21 +-`＝18 增 3 删，两块 hunk＝注释改写 ＋ 偏移数字）＋
  新增用例 `internal/audio/parse_wave_format_300_windows_test.go`（157 行）。
  其余 23 枚＝`probes/300/r1/**`（证据件与原始 stdout）＋ 票 300 的 `## Progress log` **1 行**。
- ⚠ `logs/vet.txt` 是**真 0 字节**：`go vet` 成功时 stdout 本来就空，`rc-vet=0` 写在其位（`20-ac2-rig.md` §1.2）——
  这一格⛔ 属"0 字节＝没交"，交的是 `.md` 里那行 `rc`。
- ★**⛔ 一枚别人的活被我收走**：名册里除我两枚 Go 文件外，`git status --porcelain -- internal cmd` 提交前后都是**只有我这批**
  （提交前 2 行、提交后 **0 行**，两枚读数都在 `logs/post-commit-recheck.txt`）。
  全仓面上另有的 ` M .gitignore`、` D design/*`、他腿未跟踪件等**一枚没进我这一笔**。

## 2. "跑过全量对的那棵树 == 入库的那两枚 blob"（★这一发是整套门禁的地基，逐字＝`logs/post-commit-recheck.txt`）

```
rc-git-archive=0 rc-tar=0
cmp-code=IDENTICAL            # git show 028529fc:…wasapi_windows.go | cmp - 工作树
cmp-test=IDENTICAL            # 同一条尺，用例那枚
cmp-treeB-code=IDENTICAL      # 新 HEAD 导出树(treeC) ↔ 跑过 postB-1/postB-2 的那棵 treeB
cmp-treeB-test=IDENTICAL
```

⇒ 全量对（`postB-1`／`postB-2`）跑的就是**入库字节**（`cmd/wisp` 那 2217 行名册逐枚同名单同数，见 `30-gates.md` §5）。
⇒ 唯一一次"跑完之后又改"的差是测试文件头顶注释那一段散文（⛔ 代码、⛔ 断言），
处置＝改后我又用最终字节在 treeB 里跑了一发 audio-only（`rc-audio-final-in-treeB=0`、`pass-all=5`、`pass-top=1`，`logs/postB-entry-final.txt`），
并具名交 `300-v2`（`30-gates.md` §5 末条）。

真 HEAD blob 形态的仓外格式尺（导出树 `treeC` 里跑，空列表＝良构）：

```
$ gofmt -l  <treeC>/internal/audio/wasapi_windows.go <treeC>/internal/audio/parse_wave_format_300_windows_test.go
$ gofumpt -l（具名绝对尺，同一对路径）
（两把 stdout 全空）
```

入库偏移那行的 blob 形态逐字：

```
211:		sub := *(*windows.GUID)(unsafe.Add(p, 24))
```

eol（入库后 `i/` 取得到了，⛔ 再靠临时 `git add`；这正是 `90-unrun-rulers.md` §4 那枚自报的教训）：

```
i/lf    w/lf    attr/text eol=lf      internal/audio/parse_wave_format_300_windows_test.go
i/lf    w/lf    attr/text eol=lf      internal/audio/wasapi_windows.go
```

## 3. 我这把⛔ 做的事（重申，免得下一位以为 `AC#2` 里还有半格没交）

⛔ 翻任何 `- [ ]` 框（`AC#2`..`AC#5` 四格交 `300-v2`／编排者）⛔ 把票面改成 `-done` ⛔ 改票面任何原句 ⛔ push
⛔ 动 `cmd/wisp`／`thresholds.go`／`level.go` 的 `SineLevelTolerance`／`liquid.go` 的 `SilenceLevelGate`／`allowlist.txt`／D43 表／`PLAN.md`／
`docs/specs/**`／`frontend/**`／`design/**`／三枚冻结件／golden ⛔ 删文件 ⛔ 建 worktree ⛔ 放宽既有断言。
四把"我⛔ 跑／我踩坏"的账全部在 `90-unrun-rulers.md`（含那枚 `git reset -- <自己刚 stage 的路径>` 形状违规自报）。
