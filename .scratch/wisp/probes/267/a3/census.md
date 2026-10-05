# 267-a3 普查：`confirm_timeout_sec` 生效面 vs「300 秒」说谎面

> 腿型：只读普查（coordinator 派单代号 267-a3）。产物＝名册与射程判断，**不裁 `Q-77`**、不选甲乙。

## §0 起手锚

- 起手时刻：`09:53:08+0800`（`date "+%H:%M:%S%z"` 自取）。
- 起手锚点：`git rev-parse --short HEAD` = `c9600334`（自取，非抄派单）。
- 分支：`dev`，共享工作树。并发可见：同票另有 `267-a1`／`267-a2`／`267-r1`／`267-gate` 证据件已落；`git status --porcelain` 起手即见大量他腿改动（含 `design/**`、`.gitignore`、`.scratch/wisp/probes/151/**` 等），本腿**不 add 任何他枚文件**。
- 本腿禁跑清单（起手自约）：零 `go build`／`go test`／`go vet`／`gofumpt`／`wisp slo`／`sh scripts/*`；不读 `frontend/**`、`design/**`；不碰 `.scratch/wisp/issues/**` 的 AC 框；不碰 `docs/reports/pending-and-issues.md`；不改任何已跟踪文件（含冻结件 `docs/PLAN.md`、`docs/specs/**`）。
- 派单给的两枚消费点（`cmd/wisp/run.go:616`、`cmd/wisp/resident_approval_windows.go:460`）**本腿一律现量**，不抄行号。

## §1 用户会看见的字：写死「300 秒／5 分钟」的文案名册

## §2 常量与默认值的落点（`DefaultApprovalTimeout` 全引用者 + schema `default:"300"` + C18 原文）

## §3 日志与审计里带秒数的形状

## §4 文档（含冻结件标注）

## §5 会因配置生效而红的测试钉

## §6 尺读数与未跑清单（门禁读数）

## §7 判不动／量不到

## §8 我写错的读数（自我对抗）

## §9 交件判语
