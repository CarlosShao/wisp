# T18 性能补证（perf AC 的 DEFERRED 关闭）— 2026-09-20

> 背景：票 18 归档时把 AC#6「Resolver idempotence + perf: ≤1ms per call on warm handle cache」
> 自认 DEFERRED——幂等有测试（`TestResolveIdempotent`），但**性能从未被测量过**，
> 只有实现者的一句自述。audit-B 把它列为 13 张 done 票里唯一带实质动作的 MAJOR。
> 用户 2026-09-20 裁定：**补 bench，不降级**。本文件是补测证据。

## 测什么

`risk.Resolve` 是 C26 全仓唯一路径规范化入口，票 18 AC#6 给的预算是
**单次 ≤1ms、按 ≤50 次/任务** 折算。暖缓存口径 = 先对同一目标路径跑 64 次
`Resolve` 把 OS 路径/句柄查询的热度建立起来，再计时，避免把冷启动 I/O 记到解析器头上。

## 产物

- `internal/risk/pathresolver_bench_test.go` — `BenchmarkResolveWarm`（`b.ReportAllocs()`）
  + 共享 helper `seedResolveTarget` / `warmResolve`。
- `internal/risk/pathresolver_budget_norace_test.go` — `TestResolvePerCallBudget`：
  用 `testing.Benchmark` 复用基准机制取 `NsPerOp()` 再断言，**不手写计时循环**
  （墙上时钟差在本仓是禁用于超时的模式）；断言含 50 次/任务的包络。
  文件带 `//go:build !race`：race runtime 会把单次开销放大一个量级，
  在 `-race` 下断 1ms 墙钟预算没有意义——预算检查必须作为独立的无 race 一趟跑。

## 实测（本机 DESKTOP-LVS7839，i7-8750H / 32GB，WISP_ENV 未设）

```
$ go test ./internal/risk/ -run TestResolvePerCallBudget -bench BenchmarkResolveWarm \
    -benchtime=2000x -count=1 -v
=== RUN   TestResolvePerCallBudget
    pathresolver_budget_norace_test.go:34: C26 Resolve: 396768 ns/op = 0.397 ms/op (budget 1.000 ms, 2000 samples)
--- PASS: TestResolvePerCallBudget (0.82s)
BenchmarkResolveWarm
BenchmarkResolveWarm-12       2000      418902 ns/op      3688 B/op       53 allocs/op
PASS
ok      github.com/CarlosShao/wisp/internal/risk  1.771s
```

| 指标 | 实测 | 预算 | 余量 |
|---|---|---|---|
| 单次 Resolve | 0.397–0.419 ms | ≤1 ms | ~2.4x |
| 50 次/任务折算 | ≈ 19.9–21 ms | ≤50 ms | ~2.4x |
| 分配 | 3688 B / 53 allocs 每次 | 未设预算（记录用） | — |

## 判定

**AC#6 PASS**。53 allocs/次是 `reparseComponents` + `resolveHandle` 往返的产物，
不是泄漏；若将来要压分配数，属票 20（fs 工具接线，实际调用密度上来之后）的优化项，
不在本票范围内——本票只负责把"没测过"变成"测过且达标"。

## 复核方式

复制上面那条 `go test` 命令即可重跑；断言会自己失败，不需要读任何人的结论。
