# 票 15 语音引擎普查腿 `15-a1`（只读腿裁决表）— part 1 ／ 票面·前置·编译面·量具

> **尺**：`HEAD=7ac64ff1`（2026-10-26 11:31 +08，`dev`）· 票面 `Last update: 2026-09-19`
> **方法**：全部读数**现跑**（`sed`／`grep`／`python` 遍历 366 票建图／`ls`／`cat`）。⛔ 采信审计与改造方案的数字。
> **天数**：2026-09-19 → 2026-10-26 ＝ **37 天**（审计与派单都写 36，**差 1**；⛔ 影响定性）。
> **工作树**：未推 0；生产面（`internal/**` `cmd/**` `go.mod`）**无未提交改**。`git status` 只有 5 个 `.scratch/` 改动（含本腿必读的 `ruler-scope-checklist.md` 被编排者追加到 139 行）＋ 大量 `.scratch/` 未跟踪件 ⇒ **不是别人的活，未停手**。
> **part 2**（⑤ 第一步 ＋ ⑥ 技术前提 ＋ ⑦ 交回）：`docs/evidence/s1/15-a1-census-part2.md`

---

## ① 票 15 逐字

- **票号／文件**：`15` ／ `.scratch/wisp/issues/15-speech-engines-cer-harness.md`（57 行）
- **`Status:`**：`**Status:** ready-for-agent`（第 3 行）· `Claimed by: —` · **⛔ `-done`、⛔ `-void`、⛔ 被收编** ⇒ 未触发"停手上报"①
- **`Blocked by:`**（第 6 行，逐字）：`02-s0-spike, 13-audio-capture, 14-model-distribution` —— 只有这 3 张
- **`Progress log`**（第 57 行标题）：★**正文到此为止，标题下面 0 行、0 条目**。⇒ 37 天里**没有任何人写过一行**，连认领动作都没有。`git log` 独立核：该文件只有 3 个 commit，**全部在 2026-09-19**（`144151ae2` 立票、`130c09437`、`da81a6aca`），此后**零 commit**。
- **一句话它要建什么**：在**只有 `doc.go` 的 `internal/speech`** 里，用**已锁的 sherpa-onnx 1.13.8** 建 **silero VAD 判停 ＋ 流式 Paraformer(int8) ASR**（**⛔ KWS、⛔ TTS、⛔ AEC**——正文 `Out of scope` 逐字把 KWS 归 41、TTS 归 26、AEC 归 59），并附**引擎 slot 互斥**、`intra_op_num_threads=1` 钉死、**CER 量具＋两套基线 wav**。

### 全部 `AC#`（6 条；★**6 条全 `- [ ]`、0 条已满足**）

| # | 判据（票面原文要点） | 已满足？ |
|---|---|---|
| AC1 | wav 注入 → VAD 判停 → transcript 进 agent（**12 的 text path**）端到端；**球走 `Listening→Thinking`** | ⛔ 0 |
| AC2 | **CER**：`near_clean` ≤6% **且** `noisy_far` ≤15%（**提交进仓的集**）；带噪不达标**只**在有 `DEFERRED` ＋ UI 提示票时阻塞合并 | ⛔ 0 |
| AC3 | 端点：`说话-停顿-说话` 正确切分；**<300ms 语音被忽略**；**60s 语句自动判停** | ⛔ 0 |
| AC4 | slot 互斥：强推 ASR+TTS → 串行驻留、卸载真释放内存（**采样器 delta**）、切换延迟对 **1.6s** 预算记录 | ⛔ 0 |
| AC5 | 纯静音输入 → **一次重试**后 `asr` 错误 ＋「没听清」语义 | ⛔ 0 |
| AC6 | CPU：ASR 推理**单核钉死**（无线程池扇出），>3s 持续单核只在 work-peak 规则内 | ⛔ 0 |

---

## ② 前置票逐张现状 ＋ 下游清单

★**⛔ 只看票名**。逐张现跑 `Status` ＋ 勾框 ＋ 残余：

| 票 | 文件 `-done`？ | `Status:` | 勾／未勾 | 具名残余（逐字要点） |
|---|---|---|---|---|
| `02-s0-spike` | ✅ `-done` | `done` | 5 / **0** | ★`ASR→TTS 切换 4.7s vs 1600ms 预算 FAIL`→`D32 fallback 已激活`；`post-dispose 残余 +15MB → S1 SetMemoryLimit 义务`；`handles<300 与 D2D 窗口栈(~420) 冲突 → 编排者裁`；**"tickets 07/12/15 spec-ref citations NOT edited → orchestrator to add"** |
| `13-audio-capture` | ✅ `-done` | `done` | 6 / **0** | `≤192ms pre-TTS 帧留给消费者（非自回声）`；**AEC 归 26/59** |
| `14-model-distribution` | ✅ `-done` | `done` | 6 / **0** | 两条明写不改判：`真·冷启中断预演归 16`、`球窗口 Downloading 渲染并入 H1` |

⇒ ★**"前置全绿"在「`-done` ＋ 0 未勾框」这把尺上成立**（3/3）。★**但"全绿"≠"无债"**：`02` 把 **3 笔债**明确过给下游，其中**两笔落在票 15 身上** —— ①`debug.SetMemoryLimit`/GOGC 调参**被票 15 正文自己收进 S1 scope**（第 28 行）②票 15 的 `Spec refs` 引用**从未被修**（票 02 第 57 行点名"orchestrator to add"）。

### ★★★ 隐性前置（票 15 的 `Blocked by` **没写**，但 AC1 直接依赖）

- **票 `12-cli-text-path-s1-gate`**：`Status: review`、**6 勾／3 未勾**、末更 `2026-09-20`、**⛔ `-done`**。AC1 逐字要 "transcript into agent（**12 的 text path**）" ⇒ **15 有一条未申报的前置**。★注：12 的 `wisp run` 文本路可能**功能已在**（3 个未勾是面板/SLO 面），但**票面事实是 `review`**，派写腿前必须核实 12 的 text path 是否真能吃注入 transcript。

### 下游：**34 这张数不复现**

★**三种定义都现跑了**（`python` 遍历 366 票的 `Blocked by` 建图）：

| 定义 | 结果 |
|---|---|
| 直接子票（`**Blocked by:**` 里点名 `15-`） | **5 张**：`16` `26` `27` `28` `41` |
| **传递闭包** | ★**17 张，全部 alive、全部 `ready-for-agent`、每张 0 勾** |
| 我第一种粗尺（`grep '15'`） | 7 张 —— ★**假阳性**：`p15-aec-spike` 里的 `p15`、`115-…` 被当成票 15（**尺的射程第 11/14 例同款**） |
| 若把 `12+15` 并集闭包 | 35 张 alive —— 最接近 34 **但不等**，且 §3-2 说的是"票 15 解锁" |

⇒ ★★**34 是过度计数；真值是 17**。★★**方向站得住（15 确实锁住整条尾巴），量级要改**。

**17 张清单（票号 ＋ 一行标题）**

| d | 票号 | 标题 |
|---|---|---|
| 1 | 16 | S2 acceptance: voice end-to-end gate, latency segments, hotplug pre-mortems |
| 1 | 26 | TTS output: sherpa matcha engine, serial slot, Speaking pipeline, P7 gate |
| 1 | 27 | Punctuation restoration (P4): CT-Punc model, dictation-grade output |
| 1 | 28 | SessionScope (C31): Warm/Conversation keep-alive, timers, privacy gating |
| 1 | 41 | KWS wake word: Armed state, keywords+veto words, mute, opt-in privacy |
| 2 | 29 | Memory system: L1 profile extraction, L2 explicit memory, feedback loop |
| 2 | 30 | Result routing (D10): four-tier dispatch, artifacts, summaries, badge fallback |
| 2 | 32 | S4 acceptance: scenario ③ dictation + scenario ④ chat companion gates |
| 2 | 42 | Watchdog: per-state thresholds, safe unloadables, settle enforcement, alerts |
| 2 | 43 | Power & lifecycle events: suspend/resume, session end, autostart, crash recovery |
| 2 | 46 | S6 acceptance: idle/KWS gates, settle, power pre-mortems, cost verification |
| 2 | 60 | C32 RealtimeEngine: cloud S2S realtime as Path C enhancement (gated) |
| 2 | 61 | Cloud voice providers via C9: cascade ASR/TTS from configured models |
| 3 | 31 | reminder.* tools: one-shot in-process reminders (D12 narrow exception) |
| 3 | 36 | Result & history panel: chat stream, agent work states, sanitizer, transcript view |
| 3 | 40 | Security/Privacy/Cost pages: grants, blacklist status, memory ops, cost summary |
| 3 | 45 | Diagnostics bundle + disk-full/WAL guard + Focus-Assist degradation conclusion |

★**另有 2 张在 `Blocked by` 图外把活记在票 15 名下**（⛔ 在票 15 的 AC 里）：`66:10`「D32 六态表未测，归票 **15**/28/33/36」、`83:134`「`voice.*` 全部配置键 📋 票 15/26/27/41/59/60/61」。

---

## ③ "不碰共享编译面" —— ★★★**我自己核：不成立**

**票 15 要动的包**：`internal/speech/**`（**现只有 21 行 `doc.go`、零导出符号、零 importer**）＋ `testdata/asr-baseline/**`（新）＋ **`cmd/wisp/**`（三处）** ＋ `scripts/`（CER 进 CI）。

**与 `internal/panel` 重叠**：★**零**。票 15 全文无 `panel` 字样，语音链不经面板。

**与 `cmd/wisp` 重叠**：★**有，且是硬依赖，三处** —— 全部是**生产代码注释逐字点的名**：

| 位置 | 逐字 | 为什么是 15 的 |
|---|---|---|
| `cmd/wisp/models.go:25` | "Whichever ticket lands `internal/speech` **owns adding its own `handOffModel` call**" | `handOffModel` 是 `package main`（`models.go:357`），`internal/speech` **import 不到** ⇒ 只能改 `cmd/wisp` |
| `cmd/wisp/resident_audio_windows.go:57-63` | "The C8 frame channel … **nothing drains it**, because the ASR/KWS side (`internal/speech`) is unwritten" ＋ `:373` "the half of the voice chain **this ticket does not own**" | 帧消费侧＝15 的 |
| `cmd/wisp/resident_ball_windows.go:522,596` | "the wake-word listener it would pause lives in `internal/speech`, which is **still only a doc.go**" | 唤醒监听暂停＝15 的 |

⇒ ★★★**AC1（端到端＋球走 `Listening→Thinking`）与"不碰共享编译面"逻辑上互斥**：不接进 `cmd/wisp`，`internal/speech` 就是 `internal/models/assembly_reachability_121_test.go:41` 记的那类"建了没人调"。★**结论：§3-2 这句要撤回或改写为"15 的核心码不碰面板，但收尾必碰 `cmd/wisp` 三处"** ⇒ **派单必须按 `cmd/wisp` 共享编译面规矩排队（每波一枚写腿）**。

★**副作用（要记账）**：`assembly_reachability_121_test.go:41` 把 `internal/speech` 归入"4 个只有 doc.go、无可达符号"的第 9 类**永久排除**。票 15 一落地，**那句理由就假了**（测试不会红，因为名单是手维护的）⇒ **收尾时须同步那句注释**，否则又是"注释说假话"。

### 依赖现状：★**齐、⛔ 需下载、⛔ 需花钱 ⇒ 未触发"停手上报"③**

| 尺 | 现跑读数 |
|---|---|
| `ls third_party/` | `model-fetch` `sherpa-onnx` `spike-models` |
| `ls third_party/sherpa-onnx` | ✅ `onnxruntime.dll` `sherpa-onnx-c-api.dll` `sherpa-onnx-cxx-api.dll`（3/3） |
| `go.mod` | ✅ `:9 github.com/k2-fsa/sherpa-onnx-go v1.13.8`（**direct**）；`:21-23` 三个平台模块 indirect |
| `deps.toml` | ✅ `[sherpa-onnx] version=1.13.8`、`license=Apache-2.0`、archive sha256、3 个 DLL 各自 sha256、`mirror_prefix=ghfast.top`、`[onnxruntime] 1.28.2` |
| `cmd/wisp` 已经吃 cgo？ | ✅ **是**：`cmd/wisp/doctor.go`／`main.go`／`logsink_test.go` 已 import sherpa ⇒ **15 不新引入 cgo 面**（`portable-tests.sh:209` 早记 `cmd/wisp` 测试二进制需 sherpa DLL） |
| **模型文件** | ✅ 全在 `third_party/spike-models/`：`silero_vad.onnx`、`asr-encoder.int8.onnx`、`asr-decoder.int8.onnx`、`asr-tokens.txt`、matcha TTS、KWS、`punc-model.int8.onnx`（`model-fetch/`） |

---

## ④ 量具 CER 在哪 —— ★**没有独立票，它就是票 15 的 AC2；且今天⛔能跑**

- **票号／件**：**⛔ 另一张票**。CER 是**票 15 内部**的交付物（正文 `CER harness:` 段 ＋ **AC2**）。§3-2 说"量具票（CER）自包含"⇒ **"自包含"只对计算器成立，对"基线集"不成立**（见下）。
- **现跑四个尺，全部零**：`testdata/` **目录不存在**；`git ls-files "*.wav"` ⇒ **0**；`grep -rn 'CER' internal/ cmd/ tools/` ⇒ **1 命中，是 `internal/speech/doc.go:19` 的注释**；`.github/workflows/ci.yml` 无 asr/speech job。
- **唯一相关的既有件**：`scripts/spike/speech-baseline/main.go`（178 行）——★**不是 CER**，它量的是 **RSS/句柄/线程基线**，喂的是**全零静音** `make([]float32,512)`，⛔ 转写、⛔ 参考文本、⛔ 编辑距离。★**但它给了三样可直接抄的实据**：`SampleRate:16000`／`FeatureDim:80`／**`NumThreads: 1 // D32: intra_op_num_threads pinned to 1`**／512 样本分块喂流 —— 即 **cgo 路在本机是通的**。
- **CER 门禁写死在哪**：`SPEC-10:57-60`（"两个 wav 集**必须提交进 `testdata/asr-baseline/` 并随 CI 跑**，否则基线是空话"）＋ `SPEC-04:125`（`{near_clean,noisy_far}`，6%/15%）。
- ⇒ ★**AC2 的真正缺口不是代码，是数据**：**谁说话、什么设备、什么噪声、参考文本谁写**，仓里一个字都没有。★**这是 15 最容易被低估的一片**。
