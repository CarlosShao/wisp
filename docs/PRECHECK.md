# PRECHECK — 前置验证任务结论（P1 / P2 / P11 / P13 回填）

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

## P13 — 便携模式 × DPAPI 解密失败（S1 票 06 回填）

**结论：采用 SPEC 建议的显式错误回退——便携模式下 `dpapi:` 引用解密失败返回
`ErrPortableDecrypt`（错误文本直接指引改用 `env:` 引用），绝不回退明文。**

- 依据：便携安装随介质跨机器/跨用户漂移，DPAPI（CurrentUser）blob 必然解不开；
  明文回退会让配置目录里躺着可读密钥，违反 D33 底线，故失败必须显式且可操作。
- 落点：`internal/secret/store.go`（`ErrPortableDecrypt` + `WithPortable` 选项；
  错误文本显式给出 `env:` 指引）；数据目录重定位复用 `internal/proc`
  `ApplyPortableOverride`（`portable.txt` 标记，dev → `data-dev\`）。
- 验证：`TestPortableDecryptFailureExplicit`（便携 → `ErrPortableDecrypt` + `env:` 指引 +
  无回退值；便携自身 blob 仍可解；非便携 → 普通 DPAPI 错误而非 P13 错误）、
  `TestPortableOverride`（三环境重定位 + `WISP_TEST_DATA_DIR` 注入优先于标记）。
- 便携用户的推荐密钥载体即 `env:` 引用（占位值任意，CI/便携场景可用）。

## P3 — 逐模型权重许可证核实（票 14 回填，2026-09-19）

**结论：5/6 模型 Apache-2.0/MIT 可商用；matcha-icefall-zh-baker 确认为
非商用（data-baker 数据集条款）→ P3 BLOCKED，manifest 记 `status=blocked-p3`，
下载管线拒绝安装；TTS 替代模型选型升级为待决策项（见 docs/reports）。**

- 逐模型判定（判定来源 = 模型卡/官方 API/随包 LICENSE，非镜像站）：
  - **KWS zipformer wenetspeech 3.3M**：Apache-2.0（模型 README front-matter
    `license: Apache License 2.0`，ModelScope pkufool 仓库与包内 README 一致）。
  - **silero VAD**：MIT（上游 snakers4/silero-vad）。
  - **streaming paraformer bilingual int8**：Apache-2.0（ModelScope API
    `damo/speech_paraformer_asr_nat-zh-cn-16k-common-vocab8404-online`，即
    sherpa-onnx-streaming-paraformer-bilingual-zh-en 的上游）。
  - **SenseVoice offline int8**：Apache-2.0（ModelScope API `iic/SenseVoiceSmall`；
    包内 LICENSE 指向 FunASR——其代码 MIT、权重按各模型卡，SenseVoiceSmall 卡为 Apache-2.0）。
  - **CT-Punc（ct-transformer zh-en vocab272727）int8**：Apache-2.0（ModelScope API
    `damo/punc_ct-transformer_cn-en-common-vocab471067-large`）。
  - **matcha-icefall-zh-baker（TTS）**：⛔ **NON-COMMERCIAL**——sherpa-onnx 官方文档与
    包内 README 均明示 "The dataset is for non-commercial use only"（data-baker
    免费版数据集）。按 PLAN §16.6 P3 规则：**blocked-decision 登记，模型不上船**。
    manifest 条目保留（哈希/来源齐备，供溯源与开发机测量），`status: "blocked-p3"`，
    `Manager.Ensure` 对其一律拒绝（`TestP3BlockedModelRefused` 钉死）。
- 后续：S2 前必须选定可商用的中文 TTS（候选需满足：onnx 可转 + 数据许可商用 +
  sherpa-onnx 支持）；在此之前任何 TTS 模型条目不得从 blocked-p3 改回 ok，
  变更属 L2 级安全决策（D36）。

## P5 — 国内可达模型镜像调研（票 14 回填，2026-09-19）

**结论：镜像方案成立——compose model-mirror（自建，18081）+ hf-mirror.com（HF 代理）
+ ghfast.top（GitHub 传输代理）；镜像只提供字节，哈希一律来自 C29 签名清单。**

- 三档落地（全部在本票实测通过）：
  - **compose model-mirror**（nginx:alpine，18081；`docker/compose.{dev,test}.yml`）：
    dev/test 环境 `[models] mirror` 默认源，good/corrupt/missing 三档 fixture
    （SPEC-04 §9 的 sha256 拒收 / failover 用例）。
  - **hf-mirror.com**：HF 系模型（paraformer int8 分文件）的主源，S0 spike 与本票
    均实测可下（tokens.txt 跨源哈希一致：hf-mirror 与 GitHub 归档内容互证）。
    注意：部分 LFS/Xet 仓库会 401（punc 的 HF 镜像即如此）——不可依赖单一镜像，
    故所有条目都带官方 GitHub release 兜底。
  - **ghfast.top**：github.com release 的传输兜底（`https://ghfast.top/https://github.com/...`），
    下载器按 `transportFallback` 自动追加，无需写进 manifest；本票经它实测下载
    KWS（31MB）/ matcha（72MB）/ SenseVoice int8 归档（155MB）/ punc int8 归档（62MB）。
- 哈希来源纪律（D33/F3）：以上镜像全部只当字节源；`models/manifest.json` 的
  sha256 由本票从官方源下载后本地计算并入库签名。punc 归档哈希与 GitHub 官方
  `checksum.txt` 互证一致（c0d5aa5f…），KWS/VAD/paraformer/vocos/matcha
  model-steps-3 与 S0 spike 报告记录互证一致。
