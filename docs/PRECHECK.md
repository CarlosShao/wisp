# PRECHECK — 前置验证任务结论（P1 / P2 / P11 回填）

> 状态：S0 ticket 02 完成（2026-09-19）。三行结论对应 `docs/PLAN.md` §16.6 前置验证表；
> 实测细节与 JSON 证据见 `docs/SLO.md` 与 `docs/evidence/s0/02-spike-report.md`。

## P1 — D25 spike：内存实测 + X/Y 判定

**结论：路径 Y（单进程 cgo + 延迟初始化）。**

- 判定输入（D25 规则写死 + D44 要求的两个附加输入，全部落地）：
  - **内存规则**：路径 Y 形态空闲 private RSS = **16.5–16.6MB ≤ 25MB** → Y。
    （cgo+DLL 映射仅 +0.6MB private；分层 D2D/DWrite 窗口栈 +9.1MB private/+22.4MB
    shared；托盘+热键+Job ≈ 0。D25 担心的"DLL 映射顶到 40MB+"实测不成立。）
  - **cgo 崩溃概率（定性）**：低但非零（sherpa-onnx #2694 Go 绑定 session 创建偶发
    SIGSEGV、#3635 VAD 长跑 37h 溢出、历史 C API 段错误已修）。Y 下不可 recover，
    S1 的 recover 边界（§14.10）因此是硬要求；AEC/realtime（D47, S4+）引入时重评
    语音子进程拆分。
  - **10s 回落**：FAIL（Dispose+FreeOSMemory 后平台期 +15MB，15s 平坦）。
    D32 16.3.5 的强制-Y 规则未被触发（内存已选 Y），但同一残留使 Y 的
    Sleep-after-session ≈31MB > 25MB → S1 义务：`debug.SetMemoryLimit`/GOGC 调优复测。
- 五项基线（空 Go / +DLL / +session / +D2D 窗口 / +托盘热键 Job）全部有数：
  6.9 / 7.5 / 48.6(KWS) / 16.0 / 16.0 MB private。
- 影响面：语音层架构按单进程走；`Warm` 按定义装不下全模型（479MB 实测）→ 串行/
  「TTS 常驻 + ASR 按需」为预算成立前提（见 SLO.md §4）。

## P2 — goja 的 async/await + ES module + `vm.Interrupt()` 精度

**结论：Tier-2 语言水位 = ES2020+（含 async/await），无 ES modules、无 async 迭代；
`Interrupt()` 墙钟精度充足；不需要 QuickJS 绑定（前提：插件模块交给 bundler）。**

- goja `v0.0.0-20260917113740`：
  - **支持**：async/await、Promise.all/catch（微任务在 `RunString` 返回时自动排干）、
    optional chaining、nullish、BigInt、class static block、generators、destructuring、
    typed arrays、`Array.findLast`、Proxy 与 Reflect（均存在）。
  - **不支持**：`import`/`export`（静态与 `import()` 均为保留字语法错误）、
    async generators（"not supported yet"）、`for await...of`。
  - goja 无内置定时器（`setTimeout` 由宿主提供；宿主侧回调不得从其他 goroutine 直接
    触碰 VM，须经 `Dispatch`）。
- **`vm.Interrupt()` 墙钟精度**（12 次 × 目标 10/50/100/250/1000ms × 紧循环/带调用循环）：
  超时误差 **P50 +0.4–0.6ms，最坏 max ~8.25（提交内 JSON 权威值；更早一轮未入库数据为 12.2，结论不变）ms** —— D33/F5 的唯一资源约束手段可用，
  建议插件时限下限设为 ≥250ms（12ms 抖动余量 >40 倍）。
- 落点：Tier-2 插件若需模块化，构建期用 esbuild 打成单文件 CJS/IIFE 再进 goja
  （S7 票 51/52 的约定）；禁用 async 迭代语法于插件运行时（AST 扫描或编译失败即拒）。

## P11 — WebView2 冷拉起实测 + `jchv/go-webview2` Environment 共享

**结论：冷拉起达标（≤1500ms），L2 确认卡维持 WebView2 方案；Environment 确认不可共享 →
C27 单窗口复用（隐藏而非销毁）维持为强制设计。**

- **冷**（真进程冷：create + show + embed + SetHtml + 首次 JS 往返）：
  12 个子进程 × 2 轮 → **P50 880 / 1126ms，P95 1042 / 1256ms**（≤1500ms 目标 PASS；
  P11 设定的 >2s "L2 确认卡需重评" 触发条件**未命中**，未产生 blocked-decision）。
- **热**（同窗口 hide→show + 浏览器往返）：**P50 26–71ms，P95 50–80ms**（≤200ms PASS）。
- **Environment 共享**：实测每窗口新建 Environment（读 `pkg/edge` 源码确认无共享入口；
  同进程销毁+重建 P50 860–956ms，即 C27 要绕开的成本）。
- 工程注记：`w.Dispatch()` 只能配合 `w.Run()` 的事件循环（线程消息），手工 pump 消息循环
  时宿主存活探测要用 `Bind`+`Eval`（窗口消息）——票 33（PanelHost）实现时注意。
