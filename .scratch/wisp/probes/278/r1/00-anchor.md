# 278-r1 起手锚（只读取证腿）

- 时刻：`Fri Oct  9 09:49:46 CST 2026`（`date`）
- `git rev-parse HEAD` = `bf9974c229974b4dbd288d0c2aa49a41673d8243`
- `git log --oneline -3`
  - `bf9974c2` A751/A752 落账：收 281-r1 ＋收 283-r1 ⇒ 票 283 三格翻勾并 `-done`
  - `c87ad40e` 283-r1 收件：查重三处＋两处补尺落点＋三发正控四件套读数
  - `58a4b2fe` 281-r1 收口件：八张归位逐笔凭据＋AC#4 复跑对拉
- `git status --porcelain` 枚数：**758**（在飞改动：`M .gitignore`、`design/**` 一批 `D`、`.scratch/wisp/probes/**` 若干 `M`）=> 本腿一律走 `git show HEAD:<path>` / `git grep … HEAD` 快照取数，不碰工作树里别人的活。

## 落点说明（具名）

派单写面写作 `probes/278/r1/*.md`。本腿现量：`git ls-tree -r HEAD --name-only | grep -cE "^probes/"` = **0**（根目录无 `probes/`），
而 281-r1 / 283-r1 两条同族腿的落点是 `.scratch/wisp/probes/<NNN>/<rX>/*.md`（`git show --stat 58a4b2fe`／`c87ad40e` 逐字）
=> 本腿按仓内既有约定落在 `.scratch/wisp/probes/278/r1/`，未新建根级 `probes/`，写面仍是「只新建 `*/278/r1/*.md`」。

## 射程声明（本腿所有尺共用）

- 取数面＝**HEAD 快照**（`git grep … HEAD`／`git show HEAD:<file>`），⛔ 不读工作树。
- 代码尺默认射程＝`-- ':(exclude).scratch' '*.go'`（排除工单/探针树、只算 Go 文件、**含** `_test.go`，命中后人工剥定义行与注释行）。
- 名字族尺＝`git grep -nE "<标识符>" HEAD -- ':(exclude).scratch'`（不带 `*.go` 限制时含 docs，用来论证「没有别处引用」）。
- 本腿⛔ 未跑任何 `go build`／`go test`／`go vet`／`go list`／`go doc`。允许的取数只有 `git show`／`git grep`／`git log`／`sed`／`awk`／`wc`／Read。
- `design/**`、`frontend/**` 零读零写。`PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden 一字未动。
