# 票 51 落地腿 51-r1（实现腿 / Tier2 goja 运行时 · 块 A）— 交付证据

- 腿：`51-r1 impl`（块 A：VM 生命周期 / 隔离 / 中断 / 结果上限 / recover）
- 时刻：2026-11-03 12:5x +08
- 写面：`internal/plugin/goja_runtime.go`（新）＋ `internal/plugin/goja_runtime_test.go`（新）。⛔ 未动 frontend/design/risk/tools/memory/cmd/plugin-manifest/thresholds/golden/specs/PLAN。
- ⛔ 未 push（子代理只 commit）。基线 HEAD（进场时）：`d8de056ed`（依赖腿已把 `github.com/dop251/goja v0.0.0-20261102223426-a177d211bf60` 入 go.mod）。

## 块 A 交付的形状（结构已定，无待拍板项）

`GojaRuntime`（配置态，可并发复用）＋ `Execute(ctx, GojaCall) (GojaResult, error)`：

| AC | 实现落点 | 尺（测试） |
|---|---|---|
| ① 每次调用独立隔离 | `Execute` 内 `vm := goja.New()`，**每调用一枚新 VM，不缓存、不复用**；调用产物（global / 闭包 / 对象图）随 VM 一起丢弃。D3 的"懒建 / use-to-wake"仍在插件层成立（没有调用就没有 VM）。 | `TestGojaPerCallIsolation`：同一 runtime 连跑两次，第一次写 `globalThis.__wispProbe=42`，第二次 `typeof globalThis.__wispProbe == "undefined"`；并断言注册表 `Count()==0`（无残留 goroutine）。 |
| ② 墙钟 5s 中断 | 常量 `GojaDefaultWallClockBudget = 5 * time.Second`。中断由**具名 watchdog goroutine**驱动：`scope.Go(GojaWatchdogName(pluginID), ...)` 里 `select` 定时器 / scope ctx / done，命中后调 `vm.Interrupt(reason)`。⛔ 没有用墙钟时间差实现超时（超时＝`time.Timer`＋`vm.Interrupt`）；`observe.NowWallUTC()` 只出现在**测试断言**里量边界。`reason` 是具名类型 `gojaInterruptReason`（`wall_clock_budget_exceeded` / `host_cancelled`），分类靠 `*goja.InterruptedError.Value()` 反查，不解析文本。 | `TestGojaWallClockInterrupt`：默认 5s 全尺寸跑死循环 `while(true){i++}`，断言 `elapsed >= 5s`（⛔ 提前杀）且 `elapsed <= 5s+3s`（⛔ 失控），`ClassOf(err)==tool`，`Interrupted==true`，`reg.Count()==0`。`TestGojaInterruptIsPerCallBounded`：200ms 预算连跑三次快调用，⛔ 前一调用的定时器不得打断后一调用。`TestGojaHostCancelInterrupts`：宿主 ctx 取消 → `cancelled` 类且 `Interrupted==false`。 |
| ③ 结果上限 256KB | 常量 `GojaDefaultMaxResultBytes = 256 * 1024`。渲染后的文本超限即截断：`GojaResult{Truncated:true, FullBytes:截断前, ReturnedBytes:截断后}` ＋ `slog.Warn("tier2 result truncated", ...)`。截断走 `truncateAtRuneBoundary`，⛔ 不切半截 UTF-8。预算/上限被配成 <=0 时 **拒绝执行并回 `config` 类**（fail closed，不"无限跑"）。 | `TestGojaResultCapTruncation`：VM 内造 512KB 串 → `Truncated` 置位、`FullBytes==524288`、`ReturnedBytes==len(Value)==262144`；再用 5 字节小上限配 `"aaa\u00e9\u00e9\u00e9"` 钉住 rune 边界（得 `"aaaé"`＝5 字节）。`TestGojaBlockAMisconfigRejected`：四种非法配置＋缺 pluginID 全部 `config` 拒绝。 |
| ④ recover 边界在调用方 goroutine | VM 全部工作（`Compile`/`RunProgram`/入口函数调用/渲染）都在 **Execute 所在 goroutine** 上跑，经 `guardedCall` 收 panic → 内部标记 `*gojaPanic`（带 `debug.Stack()`）→ `classify` 转 D37 **`tool` 类**（`RetryAgentDecided`，回给模型自纠错）。`Execute` 外层还有一道 last-resort recover（覆盖 dispose 尾巴），并 `slog.Error` 记栈。goja 内部错误分别处理：`*goja.Exception`（JS throw）取 `Value().String()` 作可读文本；`goja.Compile` 语法错误也是 `tool` 类；已是 `*observe.Error` 的（config / 入口缺失）原样透传。 | `TestGojaPanicBoundary`：直接对 `guardedCall` 投 panic → 断言产出 `*gojaPanic` 且 `pluginID` 与栈都在，`classify` 得 `tool`，`Value==""`；端到端 `throw new Error("plugin said no")` → `tool` 类且消息里能看见 `plugin said no`；`var = =;` → `tool` 类；末尾 `reg.Count()==0`（panic 不外溢、不留 goroutine）。 |
| ⑤ ⛔ 严禁声称内存上限 | 全文件**唯一**提到 "memory limit" 的地方是文件头第 19 行的**否定式免责声明**："goja has no memory limit and this file does not pretend it has one"，并写明只有墙钟预算＋结果字节上限两道闸。`classify` 的中断文案同样写死 "goja offers no memory ceiling"。 | 本目录机检：`grep -E 'go func\(|filepath\.(Clean|Abs)|memory limit|memory cap|mem limit'` 在两块新文件里只命中第 19 行那条否定句；无任何肯定式内存上限承诺。 |

## watchdog 具名与 D38 名册（不改 observe 的前提下合规）

`GojaWatchdogName(pluginID)` = `tool-exec-<sanitize(pluginID)>-watchdog`，前缀取自 `observe.PerTaskPrefixes` 的 `tool-exec-`，故 `observe.ClassifyGoroutine` 判为 **per_task**（⛔ `CategoryUnknown`＝按 D38b 是泄漏征兆）。`TestGojaWatchdogRosterName` 用 `"com.example/My Plugin"` 钉住分类与字符清洗。owner 由 `DisposalScope.Go` 记为 `scope:tool:goja:<pluginID>`。

## C11 / DisposalScope 接缝

每调用一枚 `NewDisposalScope("tool:goja:"+pluginID, ctx, WithRegistry(r.reg))`（可由 `GojaCall.Scope` 传入上层 tool/task scope 复用）。watchdog 经 `scope.Go` 登记 → Dispose 的 step-2 join 会等它；`Execute` 的 `defer stop()` 在 Dispose **之前**关 done 并 join handle，因此正常路径 `disposal_incomplete=false`（测试日志可见 `scope disposed ... failed=0`）。不传 Scope 时 Execute 自己 Dispose（含 C11 step-4 `FreeOSMemory`）。

## 复跑尺（本腿实跑结果）

```
gofmt -l internal/plugin/goja_runtime.go internal/plugin/goja_runtime_test.go   # 空
go vet ./internal/plugin                                                        # 空
go test ./internal/plugin -run 'Goja|VM|Disposal|Interrupt|51' -count=1 -timeout 8m -v
  9 个 Goja 测试全 PASS（TestGojaWallClockInterrupt 5.00s）
ok  github.com/CarlosShao/wisp/internal/plugin  5.044s
```

## 越界自查（块 B 的活，本腿没做）

⛔ 未实现 C24 `GojaHostAPI` 任何一条（`wisp.http.request` / `wisp.fs.*` / `wisp.clipboard.*` / `wisp.notify.send` / `wisp.sysinfo.get` / `wisp.state.*`）；⛔ 未做 undeclared-global 字符串扫描与运行时双查；⛔ 未做 capability 检查 / taint flow / audit row；⛔ 未动 `tier2_enabled` 默认关闸与 L2 再确认；⛔ 未动 manifest / loader / risk(C19) 接缝。本腿把 VM 建起来、把边界（时钟、字节上限、recover、隔离、具名 goroutine、scope 接缝）钉死，块 B 只需往里挂 `vm.Set` 面。

## 已知边界与遗留（给验收腿，不当完成报）

1. **⛔ 内存无上限**（goja 能力边界，非本腿缺陷）：只约束时间与输出字节。栈深度用 goja 自身的栈溢出错误（`tool` 类），不设 Wisp 侧阈值。
2. panic 路径的**真·goja 内部 panic**（引擎自身不变量破）在生产里罕见且难以构造，本腿用 `guardedCall` 直投 panic 验边界本身，端到端用 JS throw 验分类。**宿主回调 panic** 属块 B（`vm.Set` 挂 Go 函数后才会出现），边界已在位、届时无需改结构。
3. `goja.New()` 每调用新建：`while(true)` 中断精度按 spike T02（P50 +0.4–0.6ms、max ~12ms），预算粒度 ≥250ms 安全；本腿 5s 默认留足余量。若块 B 之后要做插件级热 VM 复用，必须重新审视隔离 AC（现在靠"新 VM"证明，复用要另找证明方式）。
4. 语言级别（T02 结论：async/await+Promise 可用、ES modules 不可用、async generators 不可用 ⇒ Tier-2 = ES2020-minus-modules；Proxy/Reflect 可用）应写进插件 SDK 文档，落点在票 57 的 docs 面，⛔ 本腿写面外。
