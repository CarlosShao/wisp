# writeface-1 — 00 锚件（commit-first）

腿名 `writeface-1`（只读普查腿）。任务：为编排者量「待派落地腿逐枚文件名册 + 颜色面需求 + 互斥图 + 波次建议」。

## 起手三读数（现跑，2026-10-11）

- `date` → `Sun Oct 11 12:37:32 CST 2026`（尺＝Bash `date`，三态＝本机环境读数，⛔ 仓内件）
- `git rev-parse HEAD` → `934b350abb41523d08eca950fd7dd8d647dec148`（三态＝HEAD）
- `git rev-parse --abbrev-ref HEAD` → `dev`
- `git status --porcelain -- internal cmd` 行数 → **0**
  尺＝上面那条命令的输出行数（`wc -l`）；射程＝工作树的 `internal/**` ＋ `cmd/**` 两枚 pathspec。
  ⇒ **起手台面在这两枚 pathspec 上干净**，无需给名册打〔台面存疑〕标记。

## 台面全量备注（非本腿写面，只记不影响源码的部分）

- `git status --porcelain`（全仓，无 pathspec）输出约 39.8 KB（枚数未核＝不在本腿射程，⛔ 为此再跑一次）。
  其中含编排者点名的两枚：`.gitignore`（` M`）、`go.mod`（` M`）；另 `design/**` 有一批 ` D`／` M`，
  以及别家 probes 目录下的 ` M` 证据件。
  ⇒ 本腿**一枚都不动、不 add、不提交**，逐字照硬约束办。

## 写面声明

**本腿写面只有 `.scratch/wisp/probes/orch/20261011/writeface-1/**`，⛔ 动任何源码／票面／台账／`docs/**`。**

## 命令自报（本腿开跑时点）

- 已跑：`date`、`git rev-parse`、`git status`、`ls`、`mkdir -p`。
- ⛔ 未跑且整腿不会跑：`go build`／`go vet`／`go test`／`gofmt`／`gofumpt`／`sh scripts/d22scan.sh`／`scripts/lenbudget.sh`。
  理由：在飞腿 `298-v1` 正在用 `-overlay` 跑 `go test ./cmd/wisp/`，本腿任何改码/改树都会弄脏它的逐名名册。
  `go env`／`go list` 属允许面，若用到将在 `90-final.md` §4 逐条如实列。
