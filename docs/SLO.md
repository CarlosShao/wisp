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
