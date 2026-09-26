# 154 实现程 r1 — 「关」只长在环路那一枚 id 上：宿主自带 id 那一形仍没人关、并发任务零读数

派单＝`.scratch/wisp/dispatches/2026-09-26-115x-impl-154.md`
工单＝`.scratch/wisp/issues/154-the-close-only-covers-the-loop-task-id-so-host-supplied-task-ids-still-have-no-owner-and-concurrent-tasks-have-zero-readings-reserved-with-a-trigger-gate.md`
本件是**实现程**自己那一遍；票面框由编排者按**非实现者**验收表定，本程不自勾。
本程写面（派单"git 纪律与写面"一节）＝`internal/tools/bridge.go` 的**注释面** ＋ 本件 ＋ `.scratch/wisp/probes/154/`。

---

## 第 0 节 锚点、工具链、读数纪律

**共享树在取数期间自己往前走**（这是本票第 0 节必须先写的，否则下一位对不上行号）：

| 时刻 | `git rev-parse HEAD` | 本程用它量过什么 |
|---|---|---|
| 进场 | `317805c` | `git log`／`git status`（`git cat-file -t` 见下） |
| 进场第二读 | `d0865ff` | 两枚包"改前"门禁跑的就是这一枚工作树 |
| AC#1 三枚子句现读 | `e95caf1` → `bcf7601` | 同一把尺两枚锚点各跑一次，输出逐字一致（差的那几枚 commit 只碰 `frontend/**`、`docs/evidence/s1/152-…`、`docs/reports/pending-and-issues.md`、`.scratch/wisp/issues/154-…` 作废牌改名） |

⇒ **`internal/**` 与 `cmd/**` 在 `d0865ff..HEAD` 之间零改动**，现读命令与输出：

```
$ git diff --name-status d0865ff..HEAD | grep -E "^.\s+(internal|cmd)/" ; echo "rc=$?"
rc=1                     ← 零命中（尺＝git diff --name-status 全量，正则筛路径前缀 internal/ 与 cmd/）
$ git status --short -- internal cmd     →  空（改前，工作树＝HEAD）
$ git diff --cached --name-only          →  空（进场时 index 干净：上一枚 154 作废牌的删除已由 dc2046fc 收掉）
```

工具链现读（不背别人的数）：

```
$ go version            →  go version go1.27.1 windows/amd64
$ git cat-file -t 5d46f24  →  commit        （票 151 那一枚，本件引它之前先验过类型）
$ git log -1 --format="%h %ad %s" 5d46f24  →  5d46f24 Sat Sep 26 09:44:57 2026 +0800 placeholder
      ⚠ 它的 subject 就是字面 "placeholder"（本程不改它、也不据它判任何事，只是提醒下一位别以为抄错了 sha）。
$ ls third_party/sherpa-onnx/*.dll | wc -l  →  3（onnxruntime / sherpa-onnx-c-api / sherpa-onnx-cxx-api）
      且 `git ls-files third_party | wc -l` = **0** ⇒ dll 是**未跟踪**的本机件：纯净树里没有它们。
$ git config core.autocrlf  →  true   （＋ .gitattributes 的 `* text=auto` ⇒ 本程一律不用 `git archive | tar -x` 取版）
```

**本程一律不用 `git archive | tar -x` 取版**（派单坑①）：门禁两跑跑的是**当前工作树**（改前工作树在 `internal/**`、`cmd/**` 上与 `d0865ff` 无差，已由上面那条 `git diff --name-status` 钉住），
"改后"跑的是**我改过之后的同一枚工作树**——本票的改动只有注释，因此结构上碰不到 `-overlay` 与 `-cover*` 同用那枚坑（本程全程没开 `-cover`）。

取版坑的**已知会红的正控**本程照派单要求真跑了一发（不跑正控就宣称"字节不等"＝本仓既有条）。原文与全部读数＝
`.scratch/wisp/probes/154/archive-eol-control.txt`，摘要：

```
$ git ls-files --eol | awk '$1=="i/lf" && $2=="w/crlf"' | wc -l      →  89   ← 索引形(LF) ≠ 检出形(CRLF) 的 tracked 件
    其中 attr 分布：80 枚 attr/text=auto（无 eol 钉）＋ 9 枚 attr/text eol=lf
    那 80 枚按扩展名：53 log · 20 txt · 2 tsx · ts/toml/json/js/css 各 1
$ git cat-file -s <A>:<path>  vs  wc -c < <path>（未修改件）
    .scratch/wisp/probes/149/combos2/x-p1-d11a-d11b.log  blob=3391 worktree=3436   ← 派单点名的那枚实发，逐字复到
    .scratch/wisp/probes/149/combos2/x-p1-d11a.log       blob=3567 worktree=3614
$ git archive <A> docs/reports/injection-timeline.md | tar -x -C /d/tmp/wisp154-archive && cmp … → SAME（46117 = 46117）
```

⇒ 三条结论，都带方向：① **不等是真的**，且量级＝每枚约「行数」个字节（LF→CRLF）；② 射程**不覆盖 `*.go`／`*.md`／`*.sh`／`*.sse`**
（`.gitattributes` 给它们钉了 `eol=lf` ⇒ 对这些件 `git archive | tar -x` 是等字节的）；③ 会被咬住的是那 80 枚 `text=auto` 件，
落点是 **73 枚 `.scratch/wisp/probes/**`（上一程留下的原始读数，含那 53 枚 `.log` ＋ 20 枚 `.txt`）＋ 7 枚真码**：
`deps.toml` · `design/doubao/demo/app.js` · `design/doubao/demo/styles.css` · `docs/evidence/s1/66/66-full-subset-slo-report.json` · `frontend/scripts/render-composer.tsx` · `frontend/src/components/composer.tsx` · `frontend/src/lib/panel.ts`
现读命令＝`git ls-files --eol | awk '$1=="i/lf" && $2=="w/crlf" && $3=="attr/text=auto" {print $4}' | grep -v '^\.scratch/wisp/probes/'`）。
⇒ 谁要复算上一程的日志，`git archive` 拿到的不是同一串字节；`frontend/**`／`design/**` 那几枚此刻正被别的会话未提交地改，**本程不碰、也不把它们算进任何零命中宣称**。
**本程一律按派单给的替代法**：Go 侧名册/字节比对用 `git ls-tree -r` ＋ `git cat-file --batch`／锚点行 grep，不 archive。

**读数纪律**：凡引用 sha 先 `git cat-file -t`；凡点名符号先 `git grep -w` 它在不在；行号全部按本程自己锚点现读；"零命中"一律写是哪把尺。
票面上那**三枚空号**（`OpenTask("")`／`RunTextTask`／`unboundScopes()`）本件**不出现**为"存在过的东西"；`RunTextTask` 的更正（票面 Progress log 10:3x 那条）本程照抄为：**`Loop` 上没有收外部 task id 的方法**，
真实存在的是 `cmd/wisp/run.go` 里 CLI 自己那一腿 `runTextTask`（现读 `grep -n "func runTextTask" cmd/wisp/run.go` ⇒ `:137`）。

---

## 第 1 节 格 AC#5（改前两跑）— 四数 ＋ 名册基准

派单/票面口径：**逐包单跑**、四数之外名册两向 `comm`。四数的尺先钉死（全只认行首，避开"汇总行与 `-skip` 散文造假命中"）：
`^=== RUN`＝RUN；`^--- (PASS|FAIL|SKIP)`＝顶格；`^[[:space:]]*--- (PASS|FAIL|SKIP)`＝全量；`^panic:`＝panic；红绿只认 `^--- FAIL:`。

```
$ export PATH="$PWD/third_party/sherpa-onnx:$PATH"      # 不带它就是加载期 0xc0000135＋0 条 RUN＝没跑到
$ go test -count=1 -v ./internal/tools/   > probes/154/gate-pre-tools-v.txt   rc=0  ok 20.288s
$ go test -count=1 -v ./cmd/wisp/         > probes/154/gate-pre-cli-v.txt      rc=0  ok 304.198s
```

| 包 | RUN | 顶格 | FAIL | SKIP | panic | 全量裁决 | 名册文件 |
|---|---|---|---|---|---|---|---|
| `internal/tools` | **115** | **79** | **0** | **0** | 0 | 115 | `probes/154/names-pre-tools.txt`（79 行） |
| `cmd/wisp` | **139** | **79** | **0** | **0** | 0 | 139 | `probes/154/names-pre-cli.txt`（79 行） |

名册基准对拍（与 151 那两程同尺，为了证明"改前"不是本程自己污染的树）。
⚠ `comm` 只吃**排过序**的输入，本程第一发忘了给右栏排序、读出来一枚"151 程独有的用例"，方向是反的；
两栏都 `sort` 之后重跑才是下面这个形状（`sort -c` 对两栏皆 rc=0）。尺＝`^[[:space:]]*--- (PASS|FAIL|SKIP)` 取第一列、去掉子用例（含 `/` 的行）、`sort -u`。

```
$ comm -23 mine-cli.txt theirs-cli-sorted.txt      ← 左栏独有（＝本程 pre 多的那枚）
TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne
$ comm -13 mine-cli.txt theirs-cli-sorted.txt      ← 右栏独有（＝151 程多的那枚）
(空)
$ comm -23 mine-tools.txt theirs-tools-sorted.txt  →  空
$ comm -13 mine-tools.txt theirs-tools-sorted.txt  →  空
（台件：probes/154/names-pre-{cli,tools}.txt ↔ probes/151-accept/{my-post-cli,theirs-post-tools-top}.txt，右栏另存一份 sorted）
⇒ `internal/tools` 名册与 151 验收程在它锚点下的 79 枚**逐字相同**；
⇒ `cmd/wisp` 差值只有**一枚**：`TestSLO152Corrupt…`，那是票 152 那条腿新加的用例（本票在 `cmd/wisp` 不新增、不删除任何用例）。
```

⇒ **改前基线成立**：两包 rc=0、FAIL/SKIP/panic 全 0；本票"改后"要复的就是这四数＋名册两向 `comm`（改后读数见第 3.2 节，随注释改动一起交）。
