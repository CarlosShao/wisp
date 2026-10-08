# 33-a3 起手锚（只读普查腿）

- 腿代号：`33-a3`（只读普查；⛔ 不改产码/测试；写面只有 `.md`）
- 分支：`dev`（共享工作树）
- 起手时刻：`2026-10-08 11:08:07 +0800`
- 锚要求：不早于 `e6c3be17`（编排者原话：只要求"不早于 `e6c3be17`"，不假设具体号）；本腿 HEAD 现量即 `e6c3be17`

## 现量

```
$ git log -1 --format=%H
e6c3be1725bb772fce4b94582a02ff00f3664695   rc=0

$ git log -1 --format='%h %ad %s' --date=iso
e6c3be17 2026-10-08 11:01:48 +0800 33-v4 verdict roster self-correction: the in-file commit table cannot list the commit that lists it   rc=0

$ date
2026-10-08 11:08:07 +0800   rc=0
```

## 这笔账的出处（不是本腿新发现的）

- `.scratch/wisp/probes/33/v4/verdict.md` Q1 末段（非实现者验收腿 `33-v4` 留下、明确"不顺手修"的一笔旧账）
- 涉及载体：`cmd/wisp/panel_host_windows.go` 的 `firstRoundTripLocked` / `serveNotBuiltNoticeLocked`
- 对照形制：同文件 `setPriorFocusLocked`
- 定性（由 `33-v4` 给出）：既有形制，载体是 10-01 那两刀 `13acad46` / `697b4faeb`，⛔ 不记在 33-r10 账上

## 本腿交件

- `.scratch/wisp/probes/33/a3/00-anchor.md`（本文件，起手即 commit）
- `.scratch/wisp/probes/33/a3/verdict.md`（六问逐格凭据，每件自落 `rc=N`）
- 回报正文 ≤500 字，不落任何票面/台账（⛔ 本腿不改 `docs/reports/pending-and-issues.md`、`HANDOVER.md`，不勾 AC 框）

## 工具面自报约束

- ⛔ 不跑 `go build` / `go vet` / `go test`（写码腿 `255-r1` 独占同机 Go 编译面）
- `go env` / `go list` 只读可用，但必须在回报里自报
