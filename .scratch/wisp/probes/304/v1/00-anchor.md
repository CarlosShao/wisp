# 304-v1 · 起手锚（本腿自己现跑，⛔ 抄派单）

尺＝`git rev-parse --short HEAD`／`git diff --cached --name-only | wc -l`／`git status --porcelain -- internal cmd tools scripts | wc -l`；
台面＝共享工作树 `D:/work/workspace/projects plans/Wisp`（branch `dev`）；原始输出＝`logs/a00-anchor.txt`（24 行）。

## 0. 锚读数（2026-10-11 19:0x 现跑）

- `HEAD`＝**`04f2ecbc`**（全号 `04f2ecbc9764408f001ddb9b8d9bc6689f4b05e8`）
- `BRANCH`＝`dev`
- `staged`＝`git diff --cached --name-only` 行数＝**0 枚**
- 工作树四枚射程目录 `internal cmd tools scripts`＝脏行数 **0 枚** ⇒ 派单 §2「起手应干净」成立，本腿⛔ 需要具名报回脏面
- ⚠ 本腿号 **==** 派单 §锚写的 `04f2ecbc` ⇒ 两号相同，无需具名重跑；若后续再动，各表表头「出自哪一号」按当时现跑填。

## 1. 起手 14 笔（`git log -n 14` 摘要，逐字在 `logs/a00-anchor.txt`）

`04f2ecbc` 304-r1 更正越界名册 · `715f63cb` 309-r3 final · `27f7b197` 304-r1 final+票尾 · `c7887eca` 编排者代笔 `AC#4`/`AC#4b` 名册 ·
`6fed247e` 309-r3 证据 · `bd123485` 304-r1 `AC#2` 成对读数+`AC#3` 门禁件（其后 8 笔更早，见件）

## 2. 本腿身份与写面

- 本腿＝`304-v1`＝票 304 **非实现者验收腿**，裁 `AC#1`／`AC#2`／`AC#3`／`AC#4`／`AC#4b` 五格（`AC#0` 已由编排者翻勾，⛔ 归本腿裁）。
- ⛔ 产码 ⛔ 改断言 ⛔ 翻框 ⛔ 动 `Status:` ⛔ 改名 `-done` ⛔ push ⛔ `git add -A`/`.` ⛔ `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean` ⛔ 仓内 worktree。
- 写面三件套＝① `.scratch/wisp/probes/304/v1/**`；② `docs/evidence/s1/**`（AGENTS §1.5 落点）；③ 本票末尾**一条** Progress bullet（＋1/−0）。
- 突变⛔ 共享工作树 ⇒ 一律仓外成对导出树（六件套 `go.mod go.sum internal cmd tools scripts`），每棵各自 `mktemp -d`，⛔ 复用；先数 `.go` 证非空再信它的 0。

## 3. 判据原文优先级

票面 `.scratch/wisp/issues/304-external-killed-test-binaries-orphan-mockllm-exe.md` **权威**；派单 `txt/304v1-dispatch.txt` 与编排者复量件
`probes/orch/20261011/verify-304r1-309r3.md` 的读数**⛔ 当判语、⛔ 当证据引用**（本腿依赖者一律自跑）。

## 4. 交付顺序（派单 §4）

`00-anchor.md`（本件）→ `10-self-ran-shapes.md`（`AC#2` 至少两发，正向 HEAD＋反向仓外摘守卫，每发同批带「坏世界被造出来」证据行）→
`20-cells.md`（五格逐格判语）→ `docs/evidence/s1/<短名>.md`（裁决表）→ 票尾一条 bullet → 加分项 `Q-304-r1-2`／`-4` 本机可证部分。

## 5. 已具名的待处理事项（起手登记，逐格结论归后面两件）

- 派单 §1 `AC#3` ★：`304-r1/logs` 有 3 枚 0 字节件（`build-1.txt`／`vet-1.txt`／`vet-reverse.txt`）⇒ 本腿查同格是否另有带 `rc=` 的件。
- 派单 §1 `AC#1` ★：越界尺必须**逐笔** `git show --name-only --format=`，⛔ 区间一把扫。
- 派单 §1 `AC#4` ★三件：名册进程名是否够（编排者自加 `wisp.test.exe`＝可能越界）／「数到≠归我」三问有无写死＋父 PID 是否被降级为旁证／`AC#4b` 原文是否覆盖「每枚非零计数指到一棵树」；
  另有编排者**自作主张请裁**＝五枚历史模板件一字未改、另立正典（算满足还是绕开 `AC#4`）。
- 编排者欠账 `Q-304-orch-A`（#1 的 (b) 形与 #2 的 (c) 形各缺一发真读数）＝本腿 `AC#2` 那两发正好补的方向，⛔ 引用其读数。
