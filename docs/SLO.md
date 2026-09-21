# SLO — 资源与延迟验收基线（S0 spike 实测回填）

> 状态：**S0 ticket 02 实测回填（2026-09-19）**。本文件把 `docs/PLAN.md` D32 16.3.2 /
> 16.3.3 / 16.3.4 的"目标值"逐行替换为实测值或标注"待票"。
> 原始 JSON 证据：`docs/evidence/s0/json/*.json`；方法学与原始数字：
> `docs/evidence/s0/02-spike-report.md`。
>
> 测量机：DESKTOP-LVS7839 · Windows 10 Pro 24H2（build 26100，注册表 ProductName 原样记录）
> · Intel i7-8750H（12 逻辑核）· 32GB RAM。
> 度量口径：进程**私有工作集**（Task Manager "Memory (private working set)"，内核字段
> `WorkingSetPrivateSize`，与 C30 进程树口径的关系见 §6）。

---

## 1. X/Y 拓扑判定（D25 规则，写死，已判定）

### 判定规则与输入（原文规则，无裁量）

| 输入 | 实测 | 规则 |
|---|---|---|
| **空闲 private RSS（路径 Y 形态：cgo 二进制 + DLL 已映射 + 分层窗口 + D2D + DWrite + 托盘 + 全局热键 + Job Object，无 session）** | **16.5 MB / 16.6 MB**（两次独立运行） | **≤ 25MB → 走 Y** |
| 空闲 private RSS（路径 X 主进程形态：纯 Go + 同一 shell 栈，无 cgo） | 15.8 MB | （X 上限 40MB，仅对照） |
| **10s 内能否回落**（建 ASR session → 推理 → Dispose → `debug.FreeOSMemory()` → 15s 内 250ms 间隔采样） | 回落 **FAIL**：Dispose 后平台期 ~22.3MB（pre 7.4MB），15s 内不再下降 | D32 16.3.5：若判 X 且不回落 → 强制 Y；本判定已由内存规则选 Y，本输入不改变结论 |
| **cgo 崩溃概率（定性评估，§14.10）** | 低但非零：sherpa-onnx issue #2694（open）Go 绑定在 session 创建内偶发 SIGSEGV（cgo 边界、不可 recover）；#3635（open）VAD CircularBuffer int32 溢出 ~37h 连续运行后 SIGSEGV（本产品 VAD 受会话约束，不达该时长）；历史 C API 段错误 #1059/#2821/#1217 均已修复 | 后果面：Y 下崩溃 = 整个悬浮球消失；X 下只损失语音。因内存规则明确选 Y，此输入作为 **S1 必须落实的缓解前提**记录：§14.10 的 recover 边界 + 语音子进程拓扑在 S4+ 引入 AEC/realtime 时重评（D47） |

### 结论：**路径 Y（单进程 cgo + 延迟初始化）**

- D25 的"关键未知"已回答：**onnxruntime + sherpa DLL 映射的 private 代价 ≈ +0.6MB**（空 Go
  6.9MB → cgo 空载 7.5MB；DLL 镜像主要落在 shared，5.1–5.2MB），远低于 D25 推测的 40MB+。
- 分层窗口 + D2D + DWrite 栈 private 代价 ≈ +9.1MB（D2D/DWrite DLL 镜像基本全 shared，
  +22.4MB shared）；托盘 + 全局热键 + Job Object ≈ +0.03MB（GDI 5 / User 14 对象）。
- ⚠ **随判定带走的两个 S1 义务**（不是判定变更，是达标义务）：
  1. **会话后回落**：Dispose + FreeOSMemory 后 private WS 平台期 ≈ +15MB（onnxruntime arena
     不归还 + Go 堆增长），Y 下 Sleep-after-session ≈ 31MB > 25MB。S1 必须
     `debug.SetMemoryLimit`/GOGC 调优并复测（D32 回落行）。D32 原文已定 X 空闲 40MB 的理由
     （"DLL 常驻映射甩不掉"）实测不成立——DLL 只占 0.6MB，**真正的常驻残留是 arena/堆**。
  2. **句柄约束**：D32 Sleeping 行"句柄 <300"与 D2D 窗口栈实测冲突（窗口栈 +279（相对空 Go 6.89） 句柄，
     合计 ~407–413（两轮实测区间））。见 §2 表内标注，待编排者裁定口径（窗口创建后 vs 全栈稳态）。

---

## 2. D32 16.3.2 六态表（实测回填）

口径：`docs/evidence/s0/json/01–04`；进程为 spike 单进程（产品树口径见 §6）。

| 状态 | 上限 Y / X | 实测（Y 口径） | 结论 |
|---|---|---|---|
| `Sleeping`（空闲，KWS 关，无子进程） | ≤25 / ≤40MB | **16.6MB**（cgo+DLL+全 shell 栈，两次运行 16.5/16.6）；纯 Go 15.8MB | **PASS**（首次空闲）。⚠ 会话后回落见 §1 义务 1（~31MB，待 S1） |
| `Armed`（仅 KWS） | ≤90 / ≤110MB | KWS session 常驻 48.6MB（裸 cgo 进程）+ shell ≈ **58MB** | **PASS**（余量 ~32MB）。KWS 加载 3.7–5.2s（进 Armed 的延迟成本） |
| `Warm`（VAD+ASR+TTS+标点已加载） | ≤350 / ≤400MB | ASR 281 + VAD 24 + TTS 174 ≈ **479MB**（尚无标点） | **按 D32 原定义 FAIL**（超出 115MB+）。半双工串行加载下峰值 = max(ASR,TTS) ≈ 290MB 才可守 → 与 Warm 定义冲突，处置见 §4 |
| `Conversation`（开麦多轮） | 同 Warm + 采集缓冲 ≤20MB | 未测 | 待 S4（依赖 Warm 态口径裁定） |
| 面板开启 | ≤600MB（含 WebView2 3–5 子进程） | 未测树口径（spike 只量了单进程 WebView2 窗口本身） | 待 S5 + 票 08（C30 `TreePrivateBytes`） |
| 工作峰值（推理+面板+工具并发） | ≤700MB | 语音侧峰值实测 314MB commit（ASR 常驻 + 推理）；面板树未测 | **700MB 是目标值非验收值**（维持 D32 声明）；量级判断：290（语音）+ ~400（面板树估）≈ 690，**无余量**，待 S3/S5 实测确认或修正 |
| 回落（`Settling` 后 10s 回到该态上限） | — | **FAIL**：`debug.FreeOSMemory()` 已执行（C11 步骤④），平台期 +15MB 不归（15s 窗口平坦） | 待 S1：`debug.SetMemoryLimit` + GOGC 实验；若仍不达 → 该行按实测修订（编排查） |

**附加硬约束（Sleeping 行）实测**：GDI 5（<200 ✓）· User 14 ✓ · 句柄 407–413（两轮实测区间）（<300 ✗，含窗口栈
279 个——约束与"悬浮球必须有窗口"矛盾，待裁定）· goroutine 数未采样（threads=23，含 D2D/ORT
线程池；goroutine 采样落票 08 的 SLO 采样器）。

**CPU 上限**（全核均值 ≤0.5%/2%/1%/5%/30%）：spike 未做持续 CPU 采样，待票 08（watchdog 的
SLO 采样器落地时逐态采）。

**`intra_op_num_threads = 1`**：spike 全部 session 已按此配置（D32 强制项生效）。

---

## 3. D32 16.3.3 模型常驻表（实测回填）

模型清单与 SHA256 见 `docs/evidence/s0/02-spike-report.md` §5。int8 优先（TTS 见偏差注）。
"实测常驻" = 新进程加载 + settle 后 private WS（含 ~7.5MB 进程基线）；"净增" = 减基线。

| 模型 | 量级（磁盘） | 估算常驻（D32 原文） | **实测常驻 / 净增** | 实测加载 | 何时加载 / 卸载 |
|---|---|---|---|---|---|
| KWS `kws-zipformer-wenetspeech-3.3M` int8 | 4.8MB | ~25MB | **48.6MB / +41.1MB** | 3.7–5.2s | `Armed` / 关 KWS、回落 |
| VAD silero（v4 结构，官方 `silero_vad.onnx`） | 0.6MB | ~8MB | **23.6MB / +16.1MB** | 1.4s | 会话起 / 会话结束 |
| ASR 流式 `streaming-paraformer-bilingual-zh-en` int8 | 226MB | int8 ~150–250MB | **281.1MB / +273.6MB** | 3.5s（重载 ~3.0s） | 会话起 / `Warm` 超时 |
| 标点恢复（CT-Punc 类） | — | ~100–200MB | **未测 → 待票 27（S4）** | — | ASR 后 / `Warm` 超时 |
| TTS `matcha-icefall-zh-baker`（fp32 声学 + vocos vocoder） | 125MB | ~120–300MB | **174.2MB / +166.7MB** | 5.3–6.0s（重载 4.4–5.7s） | 会话起 / `Warm` 超时 |

偏差注：**matcha zh-baker 官方未发布 int8 声学模型**（GitHub release 与 HF 仓库均只有
fp32 `model-steps-3.onnx`），实测用官方 fp32 + 官方 vocos vocoder —— 该行数字对预算是
**偏保守（偏高）方向**的偏差；若 S4 换用其他 int8 TTS（如 vits 系），此行需重测。

**D32 16.3.3 强制项逐条结果**：
1. **"一律 int8"**：ASR/KWS int8 ✓；TTS 官方无 int8（偏差如上）。
2. **700MB 上限**：语音侧证据支持"紧凑但无余量"（§2 工作峰值行）——**维持目标值非验收值**，
   待 S3/S5 面板树实测后定验收。
3. **切换延迟（S0 输出⑤）**：**FAIL，D32 降级路径激活**，见 §4。

---

## 4. ASR↔TTS 串行切换延迟（16.3.4 漏项③，对照 1600ms 预算）

半双工（D16 ⇒ Path T）串行加载实测（会话内多循环，`08-switch.json`）：

| 方向 | 构成 | 实测 |
|---|---|---|
| ASR → TTS（dispose 含 FreeOSMemory + 加载 TTS） | 56–106ms + 4370–5728ms | **P50 4726ms · min 4426ms** |
| TTS → ASR | 55–71ms + 2960–3086ms | ~3031–3141ms |
| TTS 冷加载（进程内首次） | — | 6011ms |
| **预算对照** | D32 16.3.4「首 token → TTS 首帧 ≤1600ms（需加载 TTS）」 | **FAIL（超 3 倍）** |

**结论（触发 D32 16.3.3 预设降级）**：串行切换吃掉 TTS 预算 → **退化为「TTS 常驻 + ASR
按需」**（TTS 174MB < ASR 281MB，方向正确）。代价转移：唤醒后 ASR 按需加载 ~3.0s，必须
由 C31 在唤醒瞬间预启（用户说完命令前完成）才能落进「冷会话 P50 ≤4s」——待票 28（S4
SessionScope）设计并在 12/15 验证。

**"峰值 = max(ASR,TTS) 而非 sum" 实测确认**：ASR 常驻 283–291MB ↔ TTS 常驻 176–179MB，
全程峰值 commit 314MB，无叠加 ✓（这是 350/700MB 预算成立的前提，Path T 下成立）。

---

## 5. 延迟预算其他实测行（16.3.4）

| 预算行 | 目标 | 实测 | 结论 |
|---|---|---|---|
| 面板冷拉起（会话内首次） | ≤1500ms | **P50 880–1126ms · P95 1042–1256ms**（12 子进程真冷 ×2 轮） | **PASS**；P11 的 >2s 重评触发条件未命中，L2 确认卡维持 WebView2 方案 |
| 面板热拉起（同窗口 show） | ≤200ms | **P50 26–71ms · P95 50–80ms** | **PASS**（C27 单窗口复用路线成立） |
| （参考）销毁+重建 WebView2 窗口 | — | P50 860–956ms | 量化了 C27 避免的成本；Environment 不共享（每窗口一个）实测确认 |
| 模型加载（不在任何一行的漏项④） | （D32 假设 1–3s） | ASR 3.5s · TTS 5.3–6.0s · KWS 3.7–5.2s · VAD 1.4s | **普遍超假设**：「冷会话 P50 ≤4s」内含 ASR 加载 3.5s 后仅余 0.5s 给 唤醒+VAD+ASR 推理+组装 → 结构性紧张，落票 15 实测复核并考虑 `Warm` 预载策略 |
| 其余行（唤醒→监听、VAD→ASR、标点、首 token、冷/热会话整链） | — | 未测 | 待票 12/15/26/28（各链路落地时逐行回填本表） |

---

## 6. Path C 注（D47：陪聊全双工）

D32 16.3.3 的 D47 注：Path C 全双工下 ASR 与 TTS 必须共存，`Conversation` 预算 = Warm +
TTS 增量，**由 S4 实测回填**。S0 证据给出的下限：ASR 281 + TTS 174 + VAD 24 ≈ **479MB**
（不含标点、AEC、音频缓冲），Path C 的 Conversation 态预算必然显著高于 Warm——S4 实测时
直接以此表为基线，禁止把 Path T 的串行数字套用到 Path C。

---

## 7. 度量口径说明（防"假达标"）

- **私有工作集**（本文件所有内存数字）= 内核 `WorkingSetPrivateSize`，即任务管理器
  "内存(专用工作集)"列；shared 部分另列（DLL 镜像大头）。D32 要求的**进程树口径**
  （Job Object `TreePrivateBytes`，含 WebView2 子进程）由票 08 落地后接管面板相关行。
- **采样协议**：`runtime.GC()×2 + debug.FreeOSMemory()` 后取 3–7 个样本的中位数，
  min/max 一并写入 JSON。
- spike 度量器自身的两个坑已修正并记录（防止后续票重蹈）：`QueryWorkingSetEx` 在本机
  静默返回空数据 → 改用 NtQuerySystemInformation；扫描缓冲必须 VirtualAlloc/VirtualFree
  （Go 堆暂存会污染下一轮样本）。

## 8. 证据索引

| 文件 | 内容 |
|---|---|
| `docs/evidence/s0/json/01-shell-baselines.json` | 基线 ①④⑤（纯 Go：空载 / +D2D 窗口 / +托盘热键 Job） |
| `docs/evidence/s0/json/02-speech-baselines.json` | 基线 ②③（cgo：DLL 映射 / +KWS session / dispose） |
| `docs/evidence/s0/json/03-idle-y.json` `03-idle-x.json` | X/Y 判定的两个空闲形态实测 |
| `docs/evidence/s0/json/04-unload-asr.json` | 10s 回落测试（ASR 281MB 加载→dispose→15s 采样） |
| `docs/evidence/s0/json/05-goja.json` | goja 能力 + `vm.Interrupt` 精度 |
| `docs/evidence/s0/json/06-webview2*.json` | WebView2 冷/热/重建延迟（两轮） |
| `docs/evidence/s0/json/07-residency-*.json` | 四模型常驻实测 |
| `docs/evidence/s0/json/08-switch.json` | ASR↔TTS 串行切换 |
| `docs/evidence/s0/02-spike-report.md` | 方法学、原始数字、模型清单（URL/SHA256）、机器配置 |


## 编排者裁定（2026-09-19，基于本 spike 实测）

1. **句柄口径**：D32 16.3.2 的「句柄 <300」写于未实测期；spike ④ 实测分层窗口 + D2D + DWrite
   栈即占 ~420 句柄。裁定：**含悬浮球窗口栈的态（Sleeping/Armed/Warm/Conversation/面板/峰值）
   句柄上限放宽为 <600（实测 420 + 余量）**，GDI <200 与泄漏趋势断言不变；票 08 的采样脚本
   以此为准。趋势泄漏（增幅）仍是主判据，绝对阈值是次要护栏。
2. **回落残留义务（+15MB → Sleeping ~31MB > 25MB）**：归 **S1**——票 15（speech 引擎）与票 12
   （S1 gate）必须包含 `debug.SetMemoryLimit`/GOGC 调优义务，S1 验收时以 25MB 线复测；
   若调优后仍 >25MB，按 D32 16.3.5 精神上报（不得自行放宽）。
3. **ASR↔TTS 串行切换 4.7s FAIL**：按 D32 16.3.3 预设降级激活——**「TTS 常驻 + ASR 按需」为
   Path T 的默认策略**（票 15/26 实现口径），串行方案留作 S4 复测后再评。

---

# 附录 A — 票 12 S1 空闲 SLO 验收跑（2026-09-20，桌面独占安静窗口）

**结论先说：本附录记的是一次 PARTIAL——内存/句柄/GDI/goroutine/写盘/TCP/settle 七项 PASS，
D32 `Sleeping` 行 CPU ≤0.5% 一项 FAIL（0.5450% / 0.6285%），阈值未改、样本未挑。**
另有两起测量工具缺陷（A.4），它们使「`Sleeping` 全量门」今天只能由人 + 外部测量器判。

- **机器**：DESKTOP-LVS7839 · Windows 11 10.0.26100 · Intel i7-8750H（6C12T，`runtime.NumCPU()=12`）
  · 32GB RAM · 单屏 3440×1440。
- **时间**：2026-09-20 21:15–21:35（本地 CST）= 13:15–13:35Z。
- **二进制**：`go build`（未 strip，非 `scripts/build.ps1`）自 `dc09909` 工作树——
  `build/wisp.exe` 27,210,302B、`build/balldebug.exe` 7,060,480B；本次跑前重建，非旧物。
- **桌面洁净**：每次测量前后 `tasklist` 均无 `wisp.exe`/`balldebug.exe` 残留（跑完即查，收尾自检见 A.6#6）；
  球窗口全部由本会话创建，winlive 套件自身的 `requireQuietBallDesktop` 断言 56 PASS / 0 SKIP。
- **证据文件**：`build/slo/t12/*.json`（采样器原始报告）、
  `docs/evidence/s1/12-ball-walk/`（球差分表 + 单进程走位日志 + 逐态日志 + winlive/statemachine 日志）。

## A.1 度量一：tree-private 采样器（票 08 入口 `wisp slo`，posture=skeleton）

口径 = C30 Job 进程树私有工作集（`WorkingSetPrivateSize`），与 §7 一致。**每一个样本都列**：

| 样本 | 窗口/间隔 | 私有工作集中位 | CPU 全核均值 | 句柄max | GDI | USER | goroutine | 写盘 | TCP | 判定 |
|---|---|---|---|---|---|---|---|---|---|---|
| `sleeping-1` | 10s / 250ms | 8.43MB | 0.1557% | 235 | 0 | 9 | 1 | 0 | 0 | **PASS** |
| `sleeping-2` | 10s / 250ms | 8.54MB | 0.4931% | 241 | 0 | 9 | 1 | 0 | 0 | **PASS** |
| `sleeping-3` | 10s / 250ms | 8.52MB | **0.5450%** | 239 | 0 | 8 | 1 | 0 | 0 | **FAIL（cpu_percent_all_core）** |
| `sleeping-60s` | 60s / 250ms（= D32「1min 窗口」原长度） | 9.10MB | **0.6285%** | 249 | 0 | 10 | 1 | 0 | 0 | **FAIL（同上；单样本峰值 11.24%）** |
| `sleeping-30s-int2000` | 30s / 2000ms | 8.50MB | 0.0867% | 241 | 0 | — | 1 | 0 | 0 | PASS |
| `exp-b` | 5s / 250ms（诊断跑，见 A.4#1） | 8.62MB | 0.0076% | 229 | 0 | 7 | 1 | 0 | 0 | PASS |

七项门（内存 ≤25MB / 句柄 <600 / GDI <200 / goroutine ≤6 / 写盘 ==0 / TCP ==0 / settle）在 6/6 样本上 PASS；
**CPU 项 6 个样本里 2 个 FAIL，本附录按实测录 FAIL，不取幸运样本**。

**FAIL 的成因不是球也不是 wisp 空闲**：该采样器运行在**被测树内部**，自己 `ReadTree` 的成本被计进 CPU。
三条比例证据（同一二进制、同一状态）：118 reads/60s → 0.629%；40 reads/10s → 0.545%；
15 reads/30s → 0.087%；折算每次读 16–38ms CPU。CPU 门因此今天**不可判**（≠不可达标），修法见 A.4#2。

## A.2 度量二：真实悬浮球在 `Sleeping`（44px 静态玻璃体；票 08 采样器无「带球」口径 → 树外差分测量）

`cmd/balldebug -diff`：球跑在**子进程**，父进程从树外读它，观测成本不进被测数字。
`docs/evidence/s1/12-ball-walk/diff-table.txt` 全表（行序即票 12 AC#3 的会话走位序列）：

| 态 | timers | CPU 1核/全核 | 私有工作集 | commit | 工作集 | 句柄 | GDI/USER | px≥8/255 | 成像框 | teardown |
|---|---|---|---|---|---|---|---|---|---|---|
| `Sleeping` | **no** | 0.000% / 0.000% | **11.79MB** | 110.84MB | 42.41MB | 358 | 4 / 7 | **2103** | **46×46** | clean |
| `Listening` | yes | 0.034% / 0.003% | 15.05MB | 113.80MB | 45.92MB | 376 | 4 / 8 | 4853 | 72×72 | clean |
| `Thinking` | yes | 0.016% / 0.001% | 14.42MB | 114.15MB | 45.31MB | 372 | 4 / 7 | 4844 | 72×72 | clean |
| `Acting` | yes | 0.003% / 0.000% | 12.54MB | 111.32MB | 43.43MB | 366 | 4 / 7 | 4844 | 72×72 | clean |
| `Speaking` | yes | 0.006% / 0.001% | 13.11MB | 111.29MB | 44.03MB | 374 | 4 / 8 | 4851 | 72×72 | clean |
| `Warm` | yes | 0.012% / 0.001% | 13.54MB | 112.28MB | 44.46MB | 372 | 4 / 7 | 4625 | 71×69 | clean |
| `Settling` | yes | 0.000% / 0.000% | 12.48MB | 110.15MB | 43.39MB | 366 | 4 / 7 | 4692 | 66×66 | clean |

**D32 预算对照——今天改的是 SPEC-08 §2 的 12px/0.35 → 44px 静态玻璃体，所以这张表就是「新球体还装不装得下」的答案**：

| 门 | 上限 | 本次 `Sleeping` 实测 | 判定 |
|---|---|---|---|
| 私有工作集 | ≤25MB | **11.79MB**（树外测）·8.43–9.10MB（树内采样器） | **PASS**，余量 13.2MB |
| CPU | ≤0.5% 全核 | **0.000%**（球，树外测） | **PASS**（球本体；树内口径见 A.1/A.4#2） |
| 定时器 | 零 | **timers=no**（读活对象，非策略表） | **PASS** |
| 句柄 | <600（裁定 1） | 358 | **PASS** |
| GDI | <200 | 4 | **PASS** |
| USER | 记录项 | 7 | 记录 |
| 成像 | 必须看得见 | px≥8/255 = **2103 像素**，框 46×46 | **PASS**（旧 12px/0.35 规格同位置 = 0 像素，见 SPEC-08 §2 INTERIM） |

与 §1/§2 的 S0 口径对照：spike 的 16.5/16.6MB 含 cgo+sherpa DLL 映射与完整 shell；
`balldebug` 不链 sherpa。spike 已量出 cgo+DLL 映射的 private 代价 **≈+0.6MB**（§1 结论），
故产品形态（带球 + cgo）的估算 `Sleeping` ≈ **12.4MB**，仍 ≤25MB。⚠ 今天**不存在**「产品 exe + 球」
的可测进程：`wisp run` 按 D12 无球（`cmd/wisp/run.go` 注释「the CLI has no ball to pop」），
`wisp` 也没有 GUI 子命令——**S1 唯一带球宿主是 `balldebug`**，本节数字的口径边界就是这句话。

句柄趋势（非门，登记给 A6 的同类跟踪）：单进程连走 20 态（不重建窗口）`start 119 → ball up 380 →
末态 439 → closed 439`（`one-process-cycle.log`）。全程 <600，但 20 次换态净 +59，与 A6
「建/销球净 +1/+2」不是同一件事；长跑若单调爬升应回指此处。

## A.3 settle ≤10s + `FreeOSMemory` 计数 >0（票 12 AC#2 的字面要求）

| 样本 | peak | cap | 回 cap 耗时 | 结束私有 WS | FreeOSMemoryCount | requested | 判定 |
|---|---|---|---|---|---|---|---|
| `settle-1` | 126.7MB | 25MB | **262ms** | 8.70MB | **2** | true | **PASS** |
| `settle-2` | 126.7MB | 25MB | **260ms** | 9.09MB | **2** | true | **PASS** |
| `settle-3` | 126.7MB | 25MB | **263ms** | 8.64MB | **2** | true | **PASS** |

三个样本全部 10s 窗内回落且计数器 >0（走 `observe.ReleaseMemory` 包装的 C11 第④步，非 `debug.` 裸调）。
**范围声明（不得越界解释）**：posture 是 `skeleton+settle-fixture`——被测对象是 120MB 合成缓冲经
`GC + ReleaseMemory` 归还，它证明的是**释放路径与计数器确实跑了**；它**不**证明票 15 的 onnxruntime
arena 残留能被归还，§1「义务 1（会话后 ~31MB）」**未因此闭环**，仍挂在 D32 回落行上。
`-leak` 自检今天没跑成——不是漏跑，是 A.4#1 让它 `exit=2`（见 A.5）。

## A.4 本次跑暴露的两起测量工具缺陷（**未改任何代码**，只登记）

1. **`internal/proc/treemetrics_windows.go` 的 `parseSystemProcesses` 丢快照最后一项 → `wisp slo -state X` 常态不可用。**
   循环先 `if next == 0 { break }` 再解析当前项，而最后一个进程的 `NextEntryOffset` 恰为 0，于是
   **快照里最后一个进程永远进不了 map**；新建进程排在 `ActiveProcessLinks` 末尾，所以「自己量自己」的
   采样进程常常正是被丢掉的那个 → 60 次重试（≈6s）耗尽 → fail-closed 报
   `system snapshot does not contain the sampling process`（exit 2）。
   **差分实验（决定性）**：同一条命令、同一二进制，唯一差别是「wisp 启动之后再创建一个进程」——
   `exp-a`（无）exit=2 报上述错，`exp-b`（有）exit=0 且九项门全出数。兄弟实现
   `cmd/balldebug/diff_windows.go` 的 `privateWorkingSetFor` 顺序是对的（先比 pid，再在 `next==0` 处 break）。
   修法一行：把「解析当前项 + 入 map」放到 `next == 0` 判断**之前**。
   **影响面**：`build/slo/slo-report.json`（2026-09-19T23:35Z，票 08 那次归档跑）
   `Sleeping exit=2 / Warm exit=2 / all_pass=false`——**state 口径的 SLO 门至今没有产出过一份样本窗记录**，
   只有 leak 那次偶然出过数。本附录 A.1 的六份样本是靠 A.6#1 的 keeper 绕法才拿到的。
2. **观测者在被测进程内 → CPU 门自污染**（比例证据见 A.1）。修法要么把采样器移到树外（父读子 pid，
   即 `-diff` 的做法），要么把 self CPU 单列后从门里扣除——**扣除等于改门，须先由编排者裁定**，本次没做。

## A.5 R12 的答案：这个门 CI 能重跑，还是只有人能跑？

**一半一半，而且现在跑不了的那一半正好是最贵的一半：**

- **CI 今天可重跑**（命令 = A.6#1/#2；去掉 keeper 需先修 A.4#1）：内存 ≤25MB、句柄 <600、GDI <200、
  goroutine ≤6、写盘 ==0、TCP ==0、settle ≤10s + FreeOSMemory>0。JSON 报告 + exit 0/1 齐备。
- **CI 今天必红**：`scripts/slo-check.ps1` 的 state 段（含 `-leak` 自检）在 A.4#1 修好前会按进程创建顺序
  随机 fail-closed——本次实跑 `state Sleeping exit=2 / state Warm exit=2 / leak flipped_to_fail=False →
  FATAL`。**这不是运气，是工具坏了**，所以「CI 能跑」这句话今天不成立。
- **只有人能判**：`Sleeping` 行的 **CPU ≤0.5%** 与**「带球的 44px 新球体装不装得下」**——前者要等 A.4#2
  把观测者挪出被测树，后者需要一块独占安静桌面和一个带球宿主（S1 只有 `balldebug`；
  球窗口的差分成像更是必须真桌面）。

⇒ 建议把票 12 AC#2 拆成两条：**AC#2a「CI 子集」**（内存/句柄/GDI/goroutine/写盘/TCP/settle，
每次 merge 前跑）与 **AC#2b「桌面独占全量（CPU + 带球）」**（人仪式，按 D29 同级签收，
本附录就是它的第一次正式记录）。不拆开的话，AC#2 今天既不能自动判、也没被完整判过。

## A.6 重跑命令（逐字，Git Bash）

```bash
cd "D:\work\workspace\projects plans\Wisp"
export PATH="$PWD/third_party/sherpa-onnx:$PATH"          # cgo 链 sherpa，缺它测试二进制 0xc0000135
go build -o build/wisp.exe ./cmd/wisp && go build -o build/balldebug.exe ./cmd/balldebug

# 0) 桌面洁净自检：必须无输出（有别人建的球 = 立刻停手，A5 那次假报告就是这么来的）
tasklist | grep -iE "wisp|balldebug"

# 1) tree-private 状态采样。**keeper 绕法已删除**（票 66 AC#2，修复 commit 86e868d）：
#    A.4#1 的解析缺陷修好后，无 keeper 直接跑即可采样成功；
#    CPU 行的口径见附录 C（树外测量），树内 CPU 只作记录、不作门。
./build/wisp.exe slo -state Sleeping -seconds 10 -interval-ms 250 \
  -out build/slo/verify/66-nokeeper-1.json; echo "exit=$?"

# 2) settle 行（同样不再需要 keeper）
./build/wisp.exe slo -settle -seconds 10 \
  -out build/slo/verify/66-settle-1.json; echo "exit=$?"

# 3) 带球的真实态（桌面独占；差分表 → docs/evidence/s1/12-ball-walk/diff-table.txt）
./build/balldebug.exe -diff docs/evidence/s1/12-ball-walk \
  -diff-states Sleeping,Listening,Thinking,Acting,Speaking,Warm,Settling \
  -diff-sample 5s -diff-dwell 6s

# 4) 聚合门（当时：A.4#1 未修前必红，留着以证明它坏着。今天：解析末项 + C.4 的聚合崩溃
#    都修好后 5/5 exit=0，数字与命令见附录 C.4；本行原注释按「叠加不改写」保留在此）
powershell -NoProfile -ExecutionPolicy Bypass -File scripts/slo-check.ps1 \
  -Subset smoke -SecondsPerState 10 -OutDir build/slo/t12

# 5) 门禁（本附录引用到的四支）
go vet ./internal/ball/ ./cmd/balldebug/ && go test -count=2 ./internal/ball/ && \
go test -count=1 -tags winlive ./internal/ball/ && \
go test -count=2 ./internal/statemachine/ -run TestEveryLegalRowFires -v

# 6) 收尾自检：必须无输出（残留的球会污染之后所有人的测量）
tasklist | grep -iE "wisp|balldebug"
```

---

# 附录 B — 编排者复跑 A.4 两起缺陷（2026-09-20 21:44–21:56，同一台机器、同一 `build/wisp.exe`）

复跑命令 = A.6#1/#2（含 keeper）。二进制与代理那次是**同一个文件**（21:14 建于 `dc09909`；
`eaf808e` 是纯文档 commit，代码未变），所以差异只可能来自机器状态或口径，不来自代码。
原始样本：`build/slo/verify/*.json`。**每一个样本都列，包括打脸的那几个。**

## B.1 缺陷①（`parseSystemProcesses` 丢快照末项）——**独立复现成立**

- 静态：`internal/proc/treemetrics_windows.go:240-243` 确实是先 `if next == 0 { break }`、
  之后（:247-260）才解析并 `out[pid]=s` ⇒ 链上最后一个进程进不了 map。**代码读出来的，不是听说的。**
- 动态差分（我自己跑）：`exp-a`（不造 keeper）→ **exit=2**，
  `system snapshot does not contain the sampling process`；`exp-b`（造 keeper）→ **exit≠2 且九项门全出数**。
  与代理那次同向。
- 归档影响面复核属实：`build/slo/slo-report.json`（09-19T23:35Z）`Sleeping exit=2 / Warm exit=2 / all_pass=false`。
- **一处与代理记录不符，须更正**：A.4#1 说 `exp-b` 那次 `exit=0`；我这次 `exp-b` 是 **exit=1**
  （exit=1 = 采样成功但门 FAIL，见 B.2）。`0/1` 都远好于 `2`，所以「keeper 能让采样跑通」这句仍成立，
  但「跑通 = PASS」不成立。

## B.2 缺陷②——**结论方向对，但代理给的证据形式（CPU ∝ 读次数）不成立；正确形式更强**

同一台机器、同一 `-state Sleeping`，我把窗口/间隔/keeper 都当自变量跑了一遍：

| 跑 | 窗口/间隔 | keeper | 样本数 | CPU 均值 | 判定 |
|---|---|---|---|---|---|
| `exp-b1` | 5s / 250ms | 5 | 20 | 0.7480%（逐样本报 0.775%） | **FAIL** |
| `iv250` | 10s / 250ms | 6 | 40 | **0.0387%**（40 样本里 37 个恰好 0.0000） | PASS |
| `k5a` | 10s / 250ms | 5 | 40 | 0.5838% | **FAIL** |
| `k20a` | 10s / 250ms | 20 | 40 | 0.5837% | **FAIL** |
| `k5b` | 10s / 250ms | 5 | 40 | 0.5053% | **FAIL** |
| `k20b` | 10s / 250ms | 20 | 40 | 0.2076% | PASS |
| `iv250-30s` | 30s / 250ms | 16 | ~120 | 0.6763% | **FAIL** |
| `iv2000` | 10s / 2000ms | 6 | 5 | 0.1032% | PASS |

三条推论，按证据强度排：

1. **「读得越少越省 CPU」不成立**：`iv250`（40 次读）= 0.0387% 而 `iv2000`（5 次读）= 0.1032%；
   keeper 从 5 增到 20 也没有单向变化（0.5838% vs 0.5837%；0.5053% vs 0.2076%）。
   ⇒ A.1 那句「三条比例证据 / 每次读折 16–38ms」**不能作为已证事实引用**，A.1 的表本身照录不动（那是它跑当时的真数），
   但它的**解释**由下面第 2 条替换。
2. **正确形式是「单次读的绝对成本」，而且它比 A.4#2 说的更致命**：
   `iv250` 里偶尔跳动的样本读数聚在 **0.5062–0.5276%**（间隔 250ms ⇒ 绝对 **≈1.28–1.32ms/次读**）；
   `iv2000` 三个稳定样本读 **0.0649/0.0651/0.0652%**（间隔 2000ms ⇒ 绝对 **≈1.298–1.304ms/次读**）。
   两个独立配置量出**同一个绝对值 1.3ms**，这就是观测者自己每次 `ReadTree` 的成本，与产品无关。
   ⇒ 在 250ms 间隔下 **一次读 = 0.52% > 0.5% 门**：**光是「采样了一次」就足以让那个样本超线**，
   门与仪器自身的最小可分辨量撞在同一个数量级上。这不是「噪声大」，是**这条门在这个口径下没有定义**。
3. **被测进程确实空闲**：`iv250` 40 个样本里 37 个 CPU 恰好 `0.0000`，其余为单次读的量子；
   与 A.2 树外实测球体 `Sleeping` CPU **0.000%** 相互印证。产品侧没有 CPU 问题，**有问题的是量法**。

## B.3 对 CI 的连带裁定（这条比缺陷本身更要紧）

`scripts/slo-check.ps1` 在 CI 里跑两处（`.github/workflows/ci.yml:166` smoke、`:196` full），**都没有 `continue-on-error`**。
⇒ 只修缺陷① 会把它从「今天必红」变成「**约一半概率随机红**」（我 8 个同配置样本里 5 FAIL / 3 PASS）
—— ⚠ **这条预测于 2026-09-20 23:36 被票 66 接续代理证伪，真实后果是「必红且连报告产物都没有」**：
`scripts/slo-check.ps1:129` 在 `Set-StrictMode 2.0` 下对 `Where-Object` 的零命中结果取 `.Count` 会抛，
而零命中恰是"全部通过"那条路径 ⇒ 修好 A14 之后反而**每次都在全通过时炸**；
之前它被 A14 掩盖（那时命中数 ≥2，是真数组）。修法 `@(...)` 包裹，见附录 C。
**我这条错误预测留在原文不删**：它是"从有限样本外推机制"的典型过头——8 个样本能证"会翻转"，
不能证"翻转就是唯一后果"。
这比确定性地坏更糟——它会训练所有人忽略这个门。所以票 66 的完成判据把两起缺陷**绑成一张**：
缺陷① 单独修好后 `slo-check` 仍不得作为门启用，除非 CPU 行同时改成树外口径或从该脚本的门集里移除（**移除=改门，需编排者书面裁定，票 66 里已申请**）。

## B.4 AC#2 CPU 行的裁定

**不判 PASS、也不判 FAIL，判「仪器未定义」**：票 12 AC#2 保持未勾；CPU 那一行的产品侧证据以 A.2（树外、0.000%）为准，
`wisp slo` 的树内 CPU 数字在票 66 修好前只作记录、不作门。D32 的 0.5% 阈值**一个字没动**，
A.1/A.4 的原文**没有改写**（本附录是叠加，不是替换）。归登记项 **A14/A15**，修法与判据在**票 66**。

---

# 附录 C — 票 66 修好之后的复跑（agent-ticket66b，2026-09-20 22:58–23:22 本地 = 14:58–15:22Z）

**本附录是叠加**：A.1/A.2/A.4/B 的原始样本表一个字都没改（那是各自当时的真测），这里只录修复后的数。
**结论先说三句**：① A14 修好且回归用例经独立变异复测咬住（6 红 / 24 绿）；② A15 靠 `00bbb76` 的
**树外口径**闭环——同一窗口内树外 CPU 0.0173% / 树内 CPU 0.7629%，两侧都按同一条 `<=0.5%` 判，
只有树外那条当门 ⇒ 差的确实就是观测者自己；③ 但 AC#4 在复跑时挖出**第三起仪器缺陷**：
`slo-check.ps1` 的聚合式在**全过路径**上崩溃，所以「A14 修完 CI 就绿」这句当时也不成立（见 C.4）。

- **机器**：DESKTOP-LVS7839 · Windows 10 Pro 24H2（build 26100）· Intel i7-8750H（6C12T，12 逻辑核）· 32GB RAM。
- **二进制**：本会话 `go build` 重建（15:03Z，非旧物）：`build/wisp.exe` 27,303,086B。
  代码基线 = `86e868d`（A14 解析末项）+ `00bbb76`（A15 树外口径）+ `63b61ef`（插桩）+ 本票 `dabd365`。
- **桌面独占**：每一组测量前后各跑一次 `tasklist | grep -iE "wisp|balldebug"`，**六次全部无输出**
  （含 5 连跑 keeper 与 5 连跑 slo-check 的收尾自检；C.4 那次 `No matches found` 是同一条自检的空结果）。
- **keeper**：**一个都没起**。A.4#1 的绕法自本附录起作废（A.6#1 已改为无 keeper 形态）。
- **阈值**：`cpu_percent_all_core` 的 limit 在下列每张表里都写作 `<=0.5%`，与 D32 16.3.2 原文一致，未改。
- **门集口径**：九项 = `tree_private_bytes`(≤25MB, 门) · `cpu_percent_all_core`(≤0.5%, **树外那条当门**) ·
  `gdi_objects`(<200, 门) · `user_objects`(记录) · `handles`(<600, 门) · `goroutines`(≤6, 门) ·
  `threads`(记录) · `disk_write_ops`(==0, 门) · `tcp_connections`(==0, 门) ⇒ 七门两记录。

## C.1 AC#2：`wisp slo -state Sleeping -seconds 10 -interval-ms 250`，零 keeper，连跑 5 次

原始 JSON：`docs/evidence/s1/66/66b-nokeeper-{1..5}.json`。**五个退出码逐字照录，没有重跑换样本。**

| 跑 | 起测 UTC | 退出码 | 样本数( subj/obs) | sample_errors | **主体 CPU（树外，当门）** | 观测者 CPU（树内，仅记录） | 主体私有 WS 中位 | 观测者私有 WS 中位 | 句柄max | GDI | USER | gor | 写盘 | TCP | 九门出数 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| 1 | 15:03:50 | **exit=0** | 40 / 40 | 0 / 0 | **0.0130%**（`0.013%` vs `<=0.5%`，pass） | 0.6739%（`0.674%`，pass=false，gate=false） | 4.39MB | 8.97MB | 189 | 0 | 1 | 1 | 0 | 0 | 9/9 |
| 2 | 15:04:10 | **exit=0** | 40 / 40 | 0 / 0 | **0.0000%** | 0.5963%（pass=false，gate=false） | 4.35MB | 8.71MB | 183 | 0 | 1 | 1 | 0 | 0 | 9/9 |
| 3 | 15:04:29 | **exit=0** | 40 / 40 | 0 / 0 | **0.0000%** | 0.6350%（pass=false，gate=false） | 4.31MB | 8.72MB | 183 | 0 | 1 | 1 | 0 | 0 | 9/9 |
| 4 | 15:04:48 | **exit=0** | **39** / 40 | **1** / 0 | **0.0000%** | 0.6872%（pass=false，gate=false） | 4.27MB | 8.53MB | 189 | 0 | 1 | 1 | 0 | 0 | 9/9 |
| 5 | 15:05:08 | **exit=0** | 40 / 40 | 0 / 0 | **0.0000%** | 0.9318%（pass=false，gate=false） | 4.21MB | 9.07MB | 189 | 0 | 1 | 1 | 0 | 0 | 9/9 |

- **AC#2 判 PASS**：五次 `exit≠2`（全 0）、九项门每次各出 9 行、主体侧无一跑丢样本窗。
  旧症状（A.4#1 的 `system snapshot does not contain the sampling process` → `exit=2`）在零 keeper 下**未复现**。
- 第 4 次那次 `sample_error` 逐字是
  `read: resource: proc: external system process snapshot: NtQuerySystemInformation: buffer never sufficient (last 1090464 bytes)`
  ——`internal/proc/treemetrics_windows.go:211-233` 的快照缓冲从 1MiB 起 6 次翻倍仍不够。**不是 A14**（A14 是解析顺序），
  是采样瞬间进程列表膨胀的偶发失败，200 次读里 1 次、被计入 `sample_errors` 且不影响该门判定。
  ⇒ 登记为**候选新项**（本附录不擅自编号），票 66 未修。
- 树内那 5 个 CPU 读数（0.5963–0.9318%）**全部 >0.5%**：同一台机器、同一窗口、同一对进程，
  换到树外读就是 0.0000–0.0130%。这就是 B.2 第 2 条「单次读约 1.3ms 的绝对成本」在修复后的直接对照。
- 与上一位代理 14:38Z 那五份 `66-nokeeper-*.json`（`00bbb76` 之前）并存：那五份主体 CPU
  0.6356–0.7625%、`pass=false`、退出码 1，因为当时门还挂在树内。**两组都是真测，不互相覆盖。**

## C.2 AC#1：变异检验由第二人独立复跑

把 `WalkSystemProcesses` 退回「先测 `next==0` 再解码当前项」，先证明变异真落盘
（`grep -n "MUTATION-66B"` → 第 76 行命中；`git diff --stat` → `1 file changed, 6 insertions(+)`），
再 `go test -count=1 -v ./internal/proc/` → **rc=1：6 红 / 24 绿 / 1 SKIP**。
比上一轮多红 3 条：本次链序恰好把采样者自己排在末项，于是
`TestSystemProcessSnapshotSeesSelf`（`system snapshot does not contain the sampling process 12328 (513 entries)`）
与 AC#3 的两条树外采样用例（`system snapshot does not contain subject <pid>`）一起转红
⇒ 树外读取并不天然免疫 A14，它免疫是因为**全仓只有一份走链实现**。撤变异后 `grep -c MUTATION` = 0、
`git diff` 空、`./internal/proc/ ./internal/observe/` rc=0。原始输出：`docs/evidence/s1/66-mutation-parse-order.md`（复测段追加，未覆盖）。

## C.3 AC#3：同一窗口内「树外 ≤0.5% 且树内 >0.5%」双条成立

| 样本 | 窗口/间隔 | 主体（树外）CPU | 观测者（树内）CPU | 主体退出 | 判据 |
|---|---|---|---|---|---|
| `66-ac3-30s.json` | 30s / 250ms（119 / 120 样本） | **0.0173%**（`0.017%`，gate=true，pass） | **0.7629%**（`0.763%`，pass=false，**gate=false**，`observer_cost:true`） | exit=0 | 两条同时成立 ✓ |
| C.1 的 5 个 10s 跑 | 10s / 250ms | 0.0000–0.0130% | 0.5963–0.9318% | exit=0 | 5/5 同时成立 ✓ |

树内那条 CPU 行仍按同一阈值判并写 `pass/fail`（`0.763%` 记 fail），只是 `gate:false` 不再当门；
JSON 里带 `note` 说明它是观测者自己的 `ReadTree` 成本。产品侧口径 = 树外那条，与 A.2 的 `balldebug` 树外 0.000% 互证。

## C.4 AC#4：`slo-check.ps1 -Subset smoke` 五连跑，以及复跑时挖出的第三起仪器缺陷

**先记账：本票自己的第一次测量差点读成假绿。** 我把脚本调用接在 `| tee | grep` 后面，
`$?` 取到的是 `grep` 的退出码 ⇒ 屏幕上五个「exit=0」全是假的（真退出码是 1）。
发现方式是 `ls` 聚合报告文件根本不存在。修正后逐字重跑，退出码直接取自脚本调用、不经任何管道。

**修复前的真实结果（同一条命令，正确取码）**：`REAL exit=1`，且**没有写出 `slo-report.json`**，
而四步检查当时**全部通过**：

```text
slo-check.ps1: state Sleeping exit=0 pass=True
slo-check.ps1: state Warm exit=0 pass=True
slo-check.ps1: settle exit=0 pass=True
slo-check.ps1: leak exit=1 flipped_to_fail=True
slo-check.ps1 : The property 'Count' cannot be found on this object. Verify that the property exists.
    + FullyQualifiedErrorId : PropertyNotFoundStrict,slo-check.ps1
```

根因（`scripts/slo-check.ps1:129`，在 `Set-StrictMode -Version 2.0` 下）：
`($results | Where-Object { -not $_.pass }).Count` —— 零命中时 `Where-Object` 返回 `$null`、
单命中返回标量，两者都没有 `.Count` ⇒ **恰好在「全过」那条路上抛终止性异常**，报告没写、CI 那两步
（`.github/workflows/ci.yml:168` smoke / `:198` full，都无 `continue-on-error`，且 upload 步骤
`if-no-files-found: error`）红。它一直没露头，是因为 A14 让每个 state 都 `exit=2` 而全 fail ⇒
过滤器命中 ≥2 项时才有 `.Count`。⇒ **只修 A14 不会让 CI 变绿，而是让它从「随机红」变成「必红且没报告」**：
B.3 那句预告被这次复跑加强成了另一种方向。修法见 commit（`@(...)` 强制成数组），
阳性对照（同一条表达式在 StrictMode 下喂 0/1/2 项失败）：`failing=0 → allPass=True`、
`failing=1 → allPass=False`、`failing=2 → allPass=False`，即修完仍会红，不是把门焊死成绿。

**修复后**（`-Subset smoke`，`-SecondsPerState` 用默认 4s；五份原始 stdout+stderr 归档
`docs/evidence/s1/66/66-smoke-run-{1..5}.log`，修复前那次崩溃的原文归档
`docs/evidence/s1/66/66-smoke-pre-fix-crash.log`）：

| 跑 | 起测 UTC | **脚本退出码** | all_pass | Sleeping | Warm | settle | leak 自检 | 主体 CPU（树外，两次都是 0.0000%） | 观测者 CPU（树内，记录） |
|---|---|---|---|---|---|---|---|---|---|
| 1 | 15:14:32 | **exit=0** | true | exit=0 pass | exit=0 pass | exit=0 pass | flipped=true | Sleeping 0.0000% / Warm 0.0000% | 0.3857% / 0.9603% |
| 2 | 15:15:23 | **exit=0** | true | exit=0 pass | exit=0 pass | exit=0 pass | flipped=true | 0.0000% / 0.0000% | 0.1612% / 0.9629% |
| 3 | 15:16:13 | **exit=0** | true | exit=0 pass | exit=0 pass | exit=0 pass | flipped=true | 0.0000% / 0.0000% | 0.5136% / 1.0322% |
| 4 | 15:17:04 | **exit=0** | true | exit=0 pass | exit=0 pass | exit=0 pass | flipped=true | 0.0000% / 0.0000% | 0.5470% / 0.6436% |
| 5 | 15:17:54 | **exit=0** | true | exit=0 pass | exit=0 pass | exit=0 pass | flipped=true | 0.0000% / 0.0000% | 0.5764% / 0.5802% |

⇒ **AC#4 判 PASS：5/5 exit=0，无一次重跑被丢弃。** 注意树内记录里 `0.1612%` 与 `1.0322%` 并存——
这正是 B.2 说的「树内那条门在这个口径下没有定义」，也是它现在只作记录的原因。

- **「state 段不再引用树内 CPU 当门」的引用位置（给复核者的 diff）**：`scripts/slo-check.ps1` 的 state 段
  从头到尾没出现过 CPU 字样，它唯一的判据是 `$reportJson.pass`（`slo-check.ps1:88`），而 `pass` 由
  `internal/observe` 的九条 verdict 聚合而来（`sampler.go:343`：`if v.Gate && !v.Pass` 才判红）。
  把 CPU 行从门上摘下来那一处是 **`internal/observe/thresholds.go:140-152`**：
  `if rep.ObserverCost { v.Gate = false; v.ObserverCost = true }`；
  **limit 字符串两侧仍由同一个 `stateCPULimit(st)` 生成**（`:137`），
  所以「改的是量法、不是门」是机器可核的事实，不是口头承诺——钉住它的用例是
  `internal/observe/observer_cost_test.go`（只许 CPU 行动、limit 串两侧必须相等、非 CPU 行的 gate 不许随口径变化）。

## C.5 AC#5 的复跑腿：`-Subset full`（六态 + settle + leak），一次

命令 = C.6#4 的 full 形态（`-SecondsPerState 6`，与 `ci.yml:198` 那一步同参数，只多给了 `-OutDir` 以便归档）。
退出码 **exit=0**，`generated_at 2026-09-20T15:21:19Z`，`machine DESKTOP-LVS7839`，`all_pass=true`。
报告原件归档：`docs/evidence/s1/66/66-full-subset-slo-report.json`（645KB，含每个 state 的**逐样本**数据）。

| state | 退出码 | pass | posture | 样本(subj/obs) | 主体 CPU 树外 | 观测者 CPU 树内 | 主体私有 WS 中位 | 句柄max |
|---|---|---|---|---|---|---|---|---|
| `Sleeping` | 0 | true | **skeleton** | 24 / 24 | 0.0000% | 0.2370% | 4.18MB | 176 |
| `Armed` | 0 | true | **skeleton** | 24 / 24 | 0.0000% | 0.6032% | 4.28MB | 176 |
| `Warm` | 0 | true | **skeleton** | 24 / 24 | 0.0000% | 0.8171% | 4.25MB | 182 |
| `Conversation` | 0 | true | **skeleton** | 24 / 24 | 0.0216% | 0.2801% | 4.26MB | 176 |
| `PanelOpen` | 0 | true | **skeleton** | 24 / 24 | 0.0000% | 0.5385% | 4.27MB | 182 |
| `WorkPeak` | 0 | true | **skeleton** | 24 / 24 | 0.0000% | 0.7734% | 4.23MB | 182 |

⚠ **口径边界（这段是本表唯一不许被省略的部分）**：六个 state 的 `posture` **全是 `skeleton`**——
今天的采样器不加载 KWS/ASR/TTS、不开面板、不放球。⇒ 本表的 `all_pass=true` 只证明
**「state 口径的九项门现在跑得通、且 skeleton 形态达标」**，**不**等于 D32 16.3.2 六态表达标；
`PanelOpen` / `WorkPeak` 两行的真数字仍挂在票 28/33/36（面板树）与票 15（模型常驻）上。
把这张表当「六态已过」引用就是假绿，措辞请照本段抄。settle 行本次实测：peak 126.5MB → 回 cap(25MB) **261ms**、
`free_os_memory_count=2`、结束私有 WS 8.50MB、`pass=true`（与 A.3 三次样本同量级）。

## C.6 逐字重跑命令（Git Bash，本机；本附录全部数字出自这些命令）

```bash
cd "D:\work\workspace\projects plans\Wisp"
export PATH="$PWD/third_party/sherpa-onnx:$PATH"   # 缺它测试/二进制 0xc0000135
go build -o build/wisp.exe ./cmd/wisp

# 0) 桌面洁净自检：必须无输出
tasklist | grep -iE "wisp|balldebug"

# 1) AC#2：零 keeper 连跑 5 次（退出码必须直接取，不要接管道 —— 见 C.4 的翻车记录）
for i in 1 2 3 4 5; do
  ./build/wisp.exe slo -state Sleeping -seconds 10 -interval-ms 250 \
    -out build/slo/verify66b/66b-nokeeper-$i.json; echo "exit=$?"
done

# 2) AC#1 变异检验：改 internal/proc/systemprocs_windows.go 的 WalkSystemProcesses
#    顺序 → 先 grep 证明变异落盘，再跑，再撤，再 grep -c MUTATION 必须为 0
go test -count=1 -v ./internal/proc/          # 变异态 rc=1（6 红/24 绿/1 SKIP），撤后 rc=0

# 3) AC#4：smoke 连跑 5 次（同样：退出码直接取）
for i in 1 2 3 4 5; do
  powershell -NoProfile -ExecutionPolicy Bypass -File scripts/slo-check.ps1 \
    -Subset smoke -OutDir "build\slo\ac4fix\r$i" > build/slo/ac4fix/run-$i.log 2>&1
  echo "exit=$?"
done

# 4) AC#5 复跑腿：full（与 ci.yml:197 同形）
powershell -NoProfile -ExecutionPolicy Bypass -File scripts/slo-check.ps1 \
  -Subset full -SecondsPerState 6 -OutDir "build\slo\ac5full"; echo "exit=$?"

# 5) 收尾自检：必须无输出
tasklist | grep -iE "wisp|balldebug"
```

## C.7 与 A/B 的关系（免得后来人以为前文被推翻）

- A.1 的六个树内样本、A.4#2 的比例说法、B.2 的八样本更正：**都照原样保留**，本附录不改写。
- A.5 那句「CI 今天必红」在修复后应读作：缺陷①已修（无 keeper 可跑，C.1）、缺陷②改为树外口径（C.3）、
  而 A.5 当时**没预见**的第三起（聚合崩溃）由 C.4 修掉 ⇒ 今天 `smoke` 5/5 exit=0（C.4）。
- A14/A15 在 `docs/reports/pending-and-issues.md` 移入「已解决」，关闭 commit 见那里的引用行。
