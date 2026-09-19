# S0 Spike 报告 — ticket 02（五项实测 + X/Y 判定 + SLO 回填）

> 完成日期：2026-09-19 · 执行代理：T02-impl
> 结论速览见 `docs/SLO.md` 与 `docs/PRECHECK.md`；本报告记录**方法学、原始数字、
> 模型清单（URL/SHA256）、机器配置与全部已知的坑/偏差**。原始 JSON：
> `docs/evidence/s0/json/*.json`。

## 1. 机器与软件配置（所有测量的运行环境）

| 项 | 值 |
|---|---|
| 主机 | DESKTOP-LVS7839（user: swq） |
| OS | Windows 10 Pro 24H2（注册表 ProductName 原样；build 26100，x64） |
| CPU | Intel Core i7-8750H @ 2.20GHz，12 逻辑核 |
| RAM | 32GB（测量时可用 5.7–7.2GB；桌面会话有并发负载） |
| Go | go1.27.1（windows/amd64），GOGC 未设置 |
| cgo | MSYS2 mingw-w64 GCC 16.2.0（sherpa-onnx-go v1.13.8 自带预编译库链接） |
| sherpa-onnx 运行时 | 1.13.8（`GetVersion()`），onnxruntime 1.28.2（`GetOnnxruntimeVersion()`），官方 MT-Release DLL |
| goja | v0.0.0-20260917113740-793a2a65c13b |
| go-webview2 | v0.0.0-20260205173254-56598839c808 |
| WebView2 Runtime | Evergreen（系统自带，Win11 24H2） |

## 2. 度量方法学

### 2.1 内存口径（所有内存数字）

- **私有工作集（private WS）**：`NtQuerySystemInformation(SystemProcessInformation)` 读取
  本进程的 `WorkingSetPrivateSize` 字段 —— 与任务管理器"内存(专用工作集)"列同源。
  选它而不是 `QueryWorkingSetEx` 的原因：后者在本机**调用成功但缓冲区内容全零**
  （伪句柄与 K32 变体均如此），已用独立探针程序验证后弃用。
- **shared WS** = `GetProcessMemoryInfo().WorkingSetSize − private WS`；
  **private commit** = `PROCESS_MEMORY_COUNTERS_EX.PrivateUsage`；
  **峰值 commit** = `PeakPagefileUsage`。
- **采样协议**：`runtime.GC()×2 + debug.FreeOSMemory()` → 间隔 200ms 取 3–7 样本 →
  报中位数（JSON 同时含 min/max）。GDI/User 对象（`GetGuiResources`）、句柄
  （`GetProcessHandleCount`）、线程数（Toolhelp 快照）同点采样。
- **防污染**：NtQSI 扫描缓冲用 `VirtualAlloc(32MB)/VirtualFree` 而非 Go 堆——Go 堆暂存
  会在样本间滞留并抬高后续读数（初版实测污染 +4–8MB，已修复；勿回退）。

### 2.2 五个 RSS 基线（输出①）

两二进制、累计阶段制（JSON `stages[]` 每阶段绝对值 + 相对前段增量）：
- `shell-baseline`（CGO_ENABLED=0）：① 空载 → ④ 分层窗口+D2D+DWrite（真实可见窗口
  260x92，`WS_EX_LAYERED|TOPMOST`，`SetLayeredWindowAttributes` α=235，D2D
  `CreateHwndRenderTarget` + `CreateSolidColorBrush` + `DrawText`（DWrite
  `CreateTextFormat`，Microsoft YaHei/zh-CN），呈现 ≥1 帧 + 泵消息 300ms；再渲 30 帧做
  泄漏观察）→ ⑤ 隐藏窗口 + `Shell_NotifyIconW(NIM_ADD)` 托盘 + `RegisterHotKey`
  （Ctrl+Alt+Y，冲突则降级尝试）+ Job Object（`KILL_ON_JOB_CLOSE`，挂当前进程）。
  测毕全部销毁（托盘 NIM_DELETE / UnregisterHotKey / DestroyWindow / COM Release）。
- `speech-baseline`（cgo）：② 进程启动即有的 DLL 映射（Go 绑定静态导入
  sherpa-onnx-c-api.dll → 加载器启动时映射，无 session）→ ③ KWS zipformer 3.3M int8
  session（NumThreads=1，喂 3s 音频并解码）→ 3b Dispose + FreeOSMemory。
- D2D/DWrite 通过 LazyDLL + 手写 COM vtable 调用；slot 序号取自 mingw-w64 头文件 MIDL
  顺序并经实际渲染验证：`ID2D1Factory::CreateHwndRenderTarget`=14；RenderTarget：
  `CreateSolidColorBrush`=8、`FillRectangle`=17、`DrawText`=27、`Clear`=47、
  `BeginDraw`=48、`EndDraw`=49（前有 IUnknown 0-2 + `ID2D1Resource::GetFactory`=3）；
  `IDWriteFactory::CreateTextFormat`=15（fontSize 是第 6 参，栈传：按 IEEE 位传入）。
  注意 win-x64 寄存器传参的 float（前 4 位置）**不能**这样传——本用法只涉及栈传 float。

### 2.3 X/Y 判定（输出②）

- `xy-verdict -role idle-y`（cgo 二进制 = DLL 已映射 + ③同款 shell 栈，无 session）→
  7 样本中位数 = Y 空闲形态；`-role idle-x`（无 cgo 构建）= X 主进程形态。
- **10s 回落测试**：pre 采样 → 建 ASR streaming paraformer int8 session（计加载耗时）→
  1s 音频喂入并解码（warmup 计时）→ Delete stream/recognizer → `debug.FreeOSMemory()` →
  250ms 间隔采样 15s，判据 = private WS 回到 pre+2MB 以下。
- **cgo 崩溃概率**：sherpa-onnx issue 库检索（`gh api search/issues`），定性归纳进
  SLO.md §1；本机未做注入式崩溃实验（无必要——§14.10 已确立"C 侧段错误不可 recover"
  的机制性结论）。

### 2.4 goja（输出③）

- 能力探针逐条 `RunString` + 断言返回值；async/await 用 `RunString` 返回的
  `*goja.Promise` 断言 `State()==Fulfilled` 且结果值正确（微任务在 RunString 返回时由
  `leave()` 排干——实测确认）。ES module 探针 = 编译 `export`/`import()` 语法。
- `vm.Interrupt()` 精度：目标延时 {10,50,100,250,1000}ms × 12 次 × 两种死循环
  （`while(true){}` / `while(true){String(1+1);}`）；另一 goroutine `time.Sleep` 后
  `vm.Interrupt()`，主线程测 `RunString` 实际返回耗时。每次新建 VM，无状态污染。

### 2.5 go-webview2（输出④）

- 时序定义：**cold** = `NewWithOptions`（建窗+show+`Embed` 阻塞至
  Environment/Controller 就绪）+ `SetHtml` + 首次**浏览器往返**（`Bind("spikePing")` +
  `Eval("spikePing()")` 回调到达）；**hot** = `ShowWindow(SW_HIDE→SW_SHOW)` + 一次泵 +
  浏览器往返；**recreate** = Destroy + 完整 NewWithOptions（同进程，Environment 重建）。
- 真冷样本：driver 模式 spawn 12 个子进程各测 1 次冷（首个 warmup 子进程用于建
  WebView2 用户数据目录，弃样）。
- 坑：`w.Dispatch()` 投递的是**线程消息**，只有 go-webview2 自己的 `Run()` 循环处理；
  手工泵消息循环下必须用 `Bind`+`Eval`（走窗口消息）。已写进 PRECHECK.md P11。

### 2.6 模型常驻 + 切换延迟（输出⑤）

- 每模型独立新进程：pre → New session（计时）→ settle → 采样（idle residency）→
  warmup 推理（ASR 1s 音频解码；KWS 1s 流解码；TTS `Generate("你好，世界。")`；VAD 喂
  1s）→ dispose（计时）→ 15s 回落窗口。所有 session `NumThreads=1`（D32 强制）。
- 切换：单进程内 3 循环 [load ASR → dispose → load TTS → dispose → load ASR]，切 换 =
  dispose(含 FreeOSMemory) + 下一模型加载；另测进程内首次 TTS 冷加载。

## 3. 原始数字汇总

### 3.1 五基线（private WS MB，中位数；JSON 01/02）

| 阶段 | private | shared | commit | 备注 |
|---|---|---|---|---|
| ① 空 Go | 6.89 (6.60–6.90) | 2.93 | 13.8 | handles 142, threads 10 |
| ② +DLL 映射（无 session） | 7.50 (7.47–7.50) | 5.14 | 15.3 | **DLL private 代价 ≈ +0.6MB** |
| ③ +KWS session | 48.61 (48.56–48.63) | 14.93 | 61.2 | KWS 加载 4.4s（另测 3.7/5.2s） |
| ③b dispose+FreeOSMemory | 21.45 (21.41–21.53) | 14.94 | 31.7 | 残留 +13.9MB 不归 |
| ④ +分层窗口 D2D/DWrite（纯 Go 进程） | 15.95 (15.92–16.03) | 25.37 | 68.3 | GDI 5, User 8, handles 387, threads 21 |
| ⑤ +托盘+热键+Job | 15.98 (15.95–16.01) | 25.47 | 68.4 | User 14, handles 407 |
| 5b 再渲 30 帧 | 16.07 (16.03–16.07) | 25.48 | 68.7 | 无泄漏（GDI/User/handles 稳定） |

### 3.2 X/Y 输入（JSON 03/04）

| 量 | 值 |
|---|---|
| idle-y（cgo+DLL+全 shell 栈） | **16.59 / 16.45 MB**（两次运行） |
| idle-x（纯 Go+shell） | 15.76 MB |
| ASR session 加载 | 3324ms（模型常驻另测 3544ms） |
| ASR warmup 解码（1s 音频） | 133ms |
| ASR 常驻 | 281.4 MB |
| dispose 后平台期 | 22.1–22.4MB（pre 7.4），15s 内平坦 → **回落 FAIL** |

### 3.3 goja（JSON 05）

- PASS：async/await+Promise.all+catch、optional chaining、nullish、BigInt、
  numeric separator、class static block、generators、template literals、
  destructuring+rest、spread、Map/Set、Symbol、typed arrays、`findLast`、
  class getter/setter、Proxy（typeof=function）、Reflect（typeof=object）。
- FAIL（能力缺失）：async generators（"Async generators are not supported yet"）、
  `for await...of`（语法错）、`export`/`import()`（保留字语法错）。
- `Interrupt` 精度（ms，12 次分布）：

| 目标 | 紧循环 p50/p95/max | 带调用循环 p50/p95/max |
|---|---|---|
| 10ms | 10.53 / 10.56 / 10.85 | 10.49 / 10.60 / 10.66 |
| 50ms | 50.48 / 50.71 / 50.88 | 50.31 / 50.85 / 53.75 |
| 100ms | 100.66 / 107.99 / 112.21 | 100.36 / 100.62 / 100.63 |
| 250ms | 250.45 / 250.62 / 252.60 | 250.55 / 256.75 / 257.89 |
| 1000ms | 1000.21 / 1011.42 / 1012.18 | 1000.33 / 1003.55 / 1007.65 |

→ 最大超前/滞后 ~12ms；建议 D33/F5 插件时限 ≥250ms。

### 3.4 WebView2（JSON 06，两轮）

| 量 | run1 | run2 |
|---|---|---|
| cold（12 子进程真冷）P50 / P95 / max | 879.7 / 1041.6 / 1116.8 ms | 1125.7 / 1256.4 /（n=12） |
| 全进程内首次 create（run 内） | 1808.9 ms（紧跟 12 子进程后，竞争态） | 1163.7 ms |
| hot show P50 / P95 | 25.7 / 49.5 ms | 71.4 / 79.5 ms |
| hot 浏览器往返 p50 | 3.5 ms | — |
| recreate P50 / P95 | 955.7 / 1221.3 ms | 859.0 / 955.5 ms |

### 3.5 模型常驻 + 切换（JSON 07/08）

| 模型 | 常驻（含 7.5 基线）/ 净增 | 加载 | dispose |
|---|---|---|---|
| KWS 3.3M int8 | 48.64 / +41.1 MB | 3693ms（另测 4360/4904/5151ms） | 19ms |
| VAD silero | 23.55 / +16.1 MB | 1401ms | 1ms |
| ASR paraformer int8 | 281.13 / +273.6 MB | 3544ms；重载 2960–3086ms | 69ms |
| TTS matcha zh-baker（fp32+vocos） | 174.20 / +166.7 MB | 5345ms；重载 4370–5728ms；进程内首冷 6011ms | 50–71ms |

切换：ASR→TTS 4426/4726/5834ms（**P50 4726 vs 预算 1600ms，FAIL**）；
TTS→ASR 3046/3141/3031ms。峰值确认：ws@ASR 283–291MB、ws@TTS 176–179MB、峰值 commit
314MB → **max 而非 sum 成立**。

## 4. 已知的坑与偏差（复现测量时必读）

1. `QueryWorkingSetEx` 本机成功但零数据 → 用 NtQSI `WorkingSetPrivateSize`（§2.1）。
2. NtQSI 扫描缓冲必须 VirtualAlloc/VirtualFree；Go 堆暂存会污染样本。
3. `GetGuiResources` 在 user32 不在 kernel32；`PROCESS_MEMORY_COUNTERS_EX` 的
   `PageFaultCount` 紧跟 `cb`（x64 无插入对齐）。
4. 手写 D2D vtable 必须按 mingw 头 MIDL 序（`GetFactory` 占 slot 3；`Clear`=47 /
   `BeginDraw`=48 / `EndDraw`=49 在 vtable 尾部）；栈传 float 参数按 IEEE 位传入。
5. `Shell_NotifyIconW` 等 shell RPC 在系统繁忙时可能阻塞数分钟（一次实测挂起、重跑
   正常）——测量脚本要容忍重跑。
6. goja 的 `Dispatch` 依赖其 `Run()` 线程消息循环；手工泵消息时用 `Bind`+`Eval`。
7. **matcha zh-baker 无官方 int8**（GitHub `tts-models` release 与 HF 仓库均确认），
   实测用 fp32——TTS 常驻数字偏保守（偏高）。
8. 测量机有桌面并发负载（ZCode 会话、浏览器），加载类数字（3.5–6s）波动明显；
   存在 AV 实时扫描疑似影响。绝对值以 JSON 为准，结论基于多轮方向一致性。

## 5. 模型清单（URL / SHA256 / 体积；存放 `third_party/spike-models/`，git 忽略）

| 模型 | 来源 URL | SHA256 | 体积 |
|---|---|---|---|
| KWS `sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01`（整包） | `https://github.com/k2-fsa/sherpa-onnx/releases/download/kws-models/sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01.tar.bz2`（经 `https://ghfast.top/` 镜像） | `b2f7c89690dc8ce4c6ed6afeab7cd800c36ad1421fb6b6302b4a4b194cf7f35f` | 31MB（解包后 36MB，含 int8+fp32 双份） |
| ├ encoder-epoch-12-avg-2-chunk-16-left-64.int8.onnx | 同包内 | `dd784973fc9d2fabb3b800d6dcd20fc3b0ca84f8e2415afe54b032878e447f4d` | — |
| ├ decoder-epoch-12-avg-2-chunk-16-left-64.int8.onnx | 同包内 | `ed83454004d5bd16d831eaf00adcd181ed7734886aab6ef440f3ffa5aa3cfe3b` | — |
| ├ joiner-epoch-12-avg-2-chunk-16-left-64.int8.onnx | 同包内 | `f79760052b87239e325f0567c752ad3130b30d92effb847d4307743c20c59a24` | — |
| └ tokens.txt / keywords.txt | 同包内 | `72316508d9119696145abc6f1f8cdc46287535c34e5ce7e595f845cb1499cf2e` / `46de8c6834da501b0a2e39a2691b17fd74ea707930dd30fbc7cc69b3d396b93f` | — |
| VAD `silero_vad.onnx` | `https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/silero_vad.onnx` | `9e2449e1087496d8d4caba907f23e0bd3f78d91fa552479bb9c23ac09cbb1fd6` | 0.64MB |
| ASR `streaming-paraformer-bilingual-zh-en` int8 | `https://hf-mirror.com/csukuangfj/sherpa-onnx-streaming-paraformer-bilingual-zh-en/resolve/main/encoder.int8.onnx`（HF 原仓库 `csukuangfj/...`；GitHub release 无 int8-only 包，整包 999MB 含 fp32 故不走） | encoder.int8 `81a70226a8934e6ed92aa1d4fc486b428b5398e2f2619ed4897b7294cab90e9a` | 157.8MB |
| ├ decoder.int8.onnx | 同仓库 `resolve/main/decoder.int8.onnx` | `f3cca9f77bb9d93c8fcbfb63ae617b6b1ee96818df3aa3b151c40658fe38594f` | 68.3MB |
| └ tokens.txt | 同仓库 `resolve/main/tokens.txt` | `59aba8873a2ed1e122c25fee421e25f283b63290efbde85c1f01a853d83cb6e6` | 74KB |
| TTS `matcha-icefall-zh-baker`（整包） | `https://github.com/k2-fsa/sherpa-onnx/releases/download/tts-models/matcha-icefall-zh-baker.tar.bz2` | `ef7ebdf5987e16a5836136a51d6f3560ca997ffd33d06a40daab5af92b4b86e5`（model-steps-3.onnx） | 整包 72MB |
| ├ Vocoder `vocos-22khz-univ.onnx` | `https://github.com/k2-fsa/sherpa-onnx/releases/download/vocoder-models/vocos-22khz-univ.onnx` | `0574a135aa1db2de6e181050db2ec528496cacd4a4701fc5d7faf9f9804c0081` | 51.4MB |
| └ lexicon.txt / tokens.txt | 同包内 | `38b886d46aefa50da6322a64d72fd595d5f4fae1051adb160d647541b1e0a4a2` / `56209b2bf609d5ac1d66ede6dae7bf5254bd3f8aa24c4a6823713d5b884d87ba` | — |

注：以上哈希为**下载后本地计算**（S0 spike 无 C29 签名清单；产品化时模型分发一律走
C29 签名清单 + 官方镜像，见 D26/D33-F3）。

## 6. 工件清单

- `scripts/spike/`（独立 Go module）：`common/`（NtQSI 内存采样、Win32/D2D/DWrite、
  托盘/热键/Job）、`shell-baseline/`、`speech-baseline/`、`xy-verdict/`、`goja-caps/`、
  `webview2-latency/`、`model-residency/`、`run.ps1`（一键构建+运行全部，JSON 落
  `docs/evidence/s0/json/`）。
- `docs/SLO.md`、`docs/PRECHECK.md`、本报告。
- 复现：`powershell -ExecutionPolicy Bypass -File scripts\spike\run.ps1 -Stage all`
  （需 Go 1.27.1 + mingw64 gcc + `third_party/spike-models/` 模型 + WebView2 Runtime）。

## 7. 未尽事项（移交编排者）

1. 票 07/12/15 的 spec-ref 补 "S0 判定引用"：按本次会话协议（不得改动其他票据）未执行，
   由编排者在状态更新时补注（SLO.md §1–§5 已含其需要的全部数字）。
2. S1 义务两项：会话后回落（SetMemoryLimit/GOGC 实验，目标 25MB）、句柄约束口径裁定
   （D2D 窗口栈 ~420 handles vs "句柄 <300"）。
3. `Warm` 定义与实测冲突（479MB > 350MB）的处置：按 D32 预设走「TTS 常驻 + ASR 按需」，
   `Warm` 行口径修订需编排者拍板。
4. 标点模型常驻（16.3.3 第 4 行）留空，待票 27（S4）。
