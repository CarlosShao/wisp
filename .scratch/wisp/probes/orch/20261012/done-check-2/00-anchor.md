# `00-anchor.md` — 只读审计腿 `done-check-2`（票 304 / 票 309 的 `-done` 就绪度）

本腿非实现者；立场＝把编排者 A856/A855 那句"六格齐 / 每格都有凭据"当假话攻。

## 1. 锚（本腿自己现跑）

尺（逐字）：
```
git rev-parse --short HEAD
git rev-parse --abbrev-ref HEAD
git status --porcelain -- internal cmd tools scripts | grep -c ''
git -c core.quotePath=false diff --cached --name-only | grep -c ''
```

读数（落盘时刻）：
- HEAD（短）＝`49153fa8`
- 分支＝`dev`
- `git status --porcelain -- internal cmd tools scripts` 行数＝`0`（empty above＝0 枚；stderr 无输出）
- 全索引暂存枚数＝`0`（empty above＝0 枚；stderr 无输出）⇒ **我看到 0 枚在飞态路径，我没动。**
- 派单正文写"两笔＝`49153fa8`（票 304）／`db1694bd`（票 309）"⇒ 与盘上**不冲突**：`49153fa8` 是 HEAD，`db1694bd` 是其父（`git log --oneline -2` 自证）。
  ⚠ 按 §0.1 规矩，派单不写锚号的义务由本件承担：**本腿全部读数锚在 `49153fa8`。**

## 2. 逐字抄回报回（抄全才算读了）

### 2.1 §0 第 3～4 步（逐字）

> 3. ★**禁跑清单（另一枚腿 `310-a1` 独占 Go 面）**：⛔ `go build`／`go vet`／`go test`／`go run`／`go list`／`gofmt`／`gofumpt`／`go env`。你要格式/编译面的结论一律写成"这一把我⛔ 跑，归编排者"。允许：`git`（含 `show`／`grep HEAD`／`log`／`diff`）、`grep`／`wc`／`sed`／`ls`、`sh scripts/check-path-length-budget.sh --with-self-test`（它是 bash＋文件尺，⛔ 需要 go；**跑之前先 `head -30` 那件自证这一条**，若发现它内部调 `go` ⇒ 具名报回并⛔ 跑）。

> 4. **第 1 笔 commit 立刻做**：只交 `00-anchor.md`，`git add -- <显式路径>` ＋ `git commit -F <自家 logs 里的消息件> -- <同样显式路径>`（**add 带 pathspec、commit 也带 pathspec**）。**⛔ 裸 `git commit`**（它提交的是整个 index）；⛔ `add -A`／`add .`／`add <目录>`。

### 2.2 §2 里带 ★ 的规矩（逐字）

> ★**证据件规矩**：⛔ `.out`；⛔ `probes/**` 下新建 `.go`；⛔ 长名／`§`（相对 ≤121、末段 ≤100、全路径 ≤165）；每件自落 `rc=N`；⛔ 0 字节件。

> ★**尺规矩**：`grep -c` 命中 0 ⇒ rc=1 会截断 `&&`，一律 `N=$(… || true)`＋件里写"empty above＝0 枚"；任何"0 命中"把 **stderr 与命中数并排落盘**（坏尺＋`2>/dev/null`＝假 0）；引原文标 provenance（HEAD blob／工作树／仓外）；名册尺**带分母与射程目录**（`git grep` 要写 `HEAD -- <哪些目录>`），⛔ 相加当一枚数。

§2 的 ⚠（一并抄回，防漏读）：

> ⚠ 本仓探针目录里有一批**突变体／blob 拷贝件**（`.scratch/wisp/probes/**` 下的 `.go`），任何"全仓 `--include=*.go`"尺都会被它们污染 ⇒ 名册尺**具名 `:(exclude).scratch`** 或直接限 `cmd internal tools scripts`，并把你用的那条尺逐字写进件里。

## 3. §0.3 的自证（跑 `check-path-length-budget.sh` 之前）

`head -30 scripts/check-path-length-budget.sh` 读数：首行 shebang＝`#!/bin/sh`，正文自述"the tracked-path-length gate (ticket 262)"，其 CI 调用行＝`sh scripts/check-path-length-budget.sh --with-self-test`。**无 go 依赖。**

再加一把（本腿现跑）：
```
grep -nE 'go (build|vet|test|run|list|env|fmt)' scripts/check-path-length-budget.sh   → rc=1（0 命中，empty above＝0 枚，无 stderr）
```
⇒ 结论：**那件内部不调 `go`，本腿允许跑它**（纯 bash＋文件尺）。第一次尝试此尺时我的正则写坏了（`grep: Unmatched ( or \`(`，rc=2），已按 ★尺规矩把 stderr 与命中数并排记下、改正后重跑；没有用 `2>/dev/null` 掩盖。

## 4. 本腿计划要交的五枚件（写面只有自家 `probes/orch/20261012/done-check-2/**`）

`00-anchor.md`（本件）／`10-cell-evidence.md`／`20-commits.md`／`30-debt-ledger.md`／`90-final.md`。
⛔ 动票面 ⛔ 动 `docs/**` ⛔ 写票尾 bullet ⛔ push ⛔ 任何 go 命令。

rc=0
