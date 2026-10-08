# ledger-audit-1 起手锚（00-anchor）

- 腿名：`ledger-audit-1`（只读复核腿；工作语言＝中文）
- 起手时刻（本地）：2026-10-08 18:29:42 +0800
- 分支：dev
- HEAD：`fdd2a72e69da558455b0621cdf2e87096ec2b69f` · Thu Oct 8 18:28:29 2026 +0800 · A727 落账（编排者清欠账 B-01；节头自纠 A724/A725 时间戳）
- 工作目录：`D:\work\workspace\projects plans\Wisp`

## `git status --porcelain | head -20`（只登记，不当尺）

```
 M .gitignore
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt
 M .scratch/wisp/probes/161/r6/logs/flip-2.txt
 M .scratch/wisp/probes/161/r6/logs/flip-3.txt
 M .scratch/wisp/probes/161/r6/logs/flip-4.txt
 M .scratch/wisp/probes/161/r6/logs/flip-5.txt
 M .scratch/wisp/probes/161/r6/logs/flip-6.txt
 M .scratch/wisp/probes/161/r6/logs/flip-baseline.txt
 M .scratch/wisp/probes/161/r6/logs/flip-restored.txt
 M .scratch/wisp/probes/242/r3/logs/probe-routed.txt
 M .scratch/wisp/probes/268/v1/evidence.md
 M cmd/wisp/panel_dispatch_binding_roster_253r1_windows_test.go
 M cmd/wisp/panel_host_windows.go
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
```

## 在飞登记（说明，不当尺）

共享工作树里有别的腿在写：`cmd/wisp` 突变中（面板 `_test.go` 与 `panel_host_windows.go` 有 M）、`docs/reports/HANDOVER.md`、`.scratch/wisp/probes/comment-fix-prep-1`。
⇒ 本腿取数一律走 `git show HEAD:<path>` / `git grep … HEAD`；不读工作树里那几处；除 `go list -deps` 外零 Go 命令。

## 本腿任务

独立复跑编排者今天写入 `docs/reports/pending-and-issues.md` 的 A720–A727 十个承重读数（名册见派单），产出 `01-recheck.md` 十行表（跑 → 比 → 判）。⛔ 不改台账、不给建议。
