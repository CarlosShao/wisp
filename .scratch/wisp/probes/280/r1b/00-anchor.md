# 280-r1b 起手锚（本腿自己的取数，不引 `8ccc9568` 任何一句当凭据）

腿名：`280-r1b`（接 `280-r1`，那枚腿 10-09 09:5x 被服务端连接中断掐死，收获已落在
`.scratch/wisp/probes/280/r1/00-anchor.md`，本腿读过，但下面每一发都是本腿现跑）。
票面：`.scratch/wisp/issues/280-seven-files-fail-plain-gofmt-decide-gofmt-vs-gofumpt-ownership.md`
（20 行，整行读过，三格 AC + 禁区）。

## 锚（本腿现跑）

- 取数时刻：`2026-10-09 10:10 +0800`
- `git rev-parse HEAD` = `2bdf385d3037b7063d797bd07428dc1d37cb2561`
- `git rev-parse --abbrev-ref HEAD` = `dev`
- `git status --porcelain | wc -l` = **763**
- `git log --oneline -3`：
  - `2bdf385d` 死腿遗产代提（两枚腿 09:5x 被服务端连接中断掐掉，⛔ 不按额度定式重派）
  - `6029b743` 277-r1 A 道真机读数…
  - `8ccc9568` 280-r1 起手锚：本腿名册 7 枚与票面同名同数（tools/ 射程贡献 0 枚）＋尺找回

⚠ 锚已漂：`280-r1` 那程的锚是 `bf9974c2`，本程起手 HEAD 是 `2bdf385d`（中间落了 `8ccc9568`／`6029b743`／`2bdf385d`）。
本件所有"文件内容"一律取 `HEAD` 快照（`git show HEAD:<path>`），**不取工作树**。

## 两把尺（本腿自己的定位尺）

- `gofmt`：`which gofmt` = `/d/work/base/go/bin/gofmt`，`go version` = `go1.27.1 windows/amd64`
  ⇒ `gofmt` 是 **go1.27.1 自带的那把**。
- `gofumpt`：**不在 PATH**（`command -v gofumpt` 无输出），但二进制在本机 GOPATH/bin：
  - 路径：`D:\work\base\gopath\bin\gofumpt.exe`（`GOPATH` 环境变量 = `D:\work\base\gopath`）
  - `gofumpt.exe --version` = **`v0.12.0 (go1.27.1)`**
  - 本腿**没有** `go install`、没有下载、没有 `-w`；跑的全部是 `-l` / `-d`（只读）。
  - ⇒ 本仓定式核对：**找不到尺 ≠ 仓里没有**（`280-r1` 已发现过一次，本腿独立复核到同一枚二进制、同一版本）。

## 射程声明（本件两处射程，必须分开读）

- **名册枚数**：出自**工作树**的 `gofmt -l internal cmd tools`（编排者 10-09 09:4x 现量＝7 枚）。
  工作树此刻不可信（另一枚腿正在 `internal/agent/approval/**` 种突变），所以名册**只作枚数与文件名对照**。
- **差异内容**：全部出自 **HEAD 快照**（`git show HEAD:<path> | gofmt -d`／`| gofumpt -d`）。
  ⛔ 本件不引用任何"工作树里那 7 枚的差异字节"。

## 禁区自查（本程）

写面只有 `.scratch/wisp/probes/280/r1b/*.md`；零 `-w`、零批量重排、零产码字节改动、不改票面、
不改台账 `docs/reports/pending-and-issues.md`、不翻任何 AC 框、零 push、零 worktree/checkout/restore；
不跑 `go build`/`go test`/`go vet`/`go list`/`go doc`；`frontend/**`／`design/**` 零读零写零转述；
三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／
`internal/perm/ticket90_persist_test.go`）一字不动、本程也不读不量。
