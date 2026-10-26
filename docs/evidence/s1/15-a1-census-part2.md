# 票 15 语音引擎普查腿 `15-a1` — part 2 ／ 第一步·技术前提·交回

> **part 1**（① 票面 ＋ ② 前置与下游 ＋ ③ 编译面与依赖 ＋ ④ 量具）：`docs/evidence/s1/15-a1-census.md`
> **尺**：`HEAD=7ac64ff1`（2026-10-26 11:31 +08，`dev`）。全部读数现跑。

---

## ⑤ 第一步该派什么（★推荐在最后，⛔"由编排者决定"）

| 候选 | 做什么（一句话） | 改哪个包 | 第一个读数 | 依赖哪条 AC |
|---|---|---|---|---|
| **甲** `15-b` **VAD 端点判停** | silero VAD 会话 ＋ 端点状态机（700/300/60000ms），输入走**已在测的** `internal/audio.WavInjector`（16k／512 帧已写死：`audio.go:43,47`） | `internal/speech/**`（**唯一**） | ★**1–2 天**：AC3 三条逐条绿（`speak-pause-speak` 切分／<300ms 忽略／60s 自动判停） | **AC3**（并顺带钉 AC6 的线程数） |
| **乙** `15-d` **CER 计算器** | 纯 Go 编辑距离 ＋ 逐件/聚合报告 ＋ `near_clean`/`noisy_far` 双门禁判定形状，注入口 `func(hyp, ref) Report` | `internal/speech/cer/**` ＋ `testdata/asr-baseline/` 的**清单/格式**（先不塞音频） | ★**当天（几小时）** | **AC2**（计算半边） |
| **丙** `15-a` **引擎基座** | `internal/speech` 的 C9 Go 接口（`AsrEngine`/`TtsEngine`/`WakeWordEngine`，`PLAN.md:1365` 只有契约行、**Go 里零定义**）＋ session create/dispose（C11）＋ `NumThreads:1` | `internal/speech/**` | **1 天**：`sherpa.GetVersion()` ＋ create/dispose 的 RSS delta 与线程数 | AC6 ＋ AC4 机制半边 |

### ★★推荐：**甲（`15-b` VAD 端点判停）**

1. ★**唯一能被一条腿完整关掉**的 AC：AC3 的夹具可用**程序合成语音段**（噪声/正弦包络＋静音间隔）造，**⛔ 需要任何语料** ⇒ 不背 ④ 那笔"谁录音"的债，也不等票 12。
2. ★**它把 15 最硬的技术风险在"产品包内"证明掉**：silero VAD 会话起来＝cgo 边界（D37b）、单实例（D5）、`NumThreads:1`（D32）同时在位；比丙更值 —— 丙是"没读数的地基"，甲是"地基＋一条真判据"。
3. ★**它不碰 `cmd/wisp`**，可以在面板/权限线的写腿在飞时并行（CPU≤70% 补位规矩不冲突）。
4. ★**它是 AC1 的必要条件**：没有 VAD 判停就没有 transcript、就没有球走 `Listening→Thinking`。★**第一条读数就是这条线上 37 天来的第一个真数**（切分正确率 ＋ <300ms 误切率）。

★**乙宜与甲同日并行**（包不重叠、零文件交集），理由＝量具必须先于被测物，否则引擎先绿、门禁后补，就是 `D22 mode-6` 的擦边（"先有关掉测量的机会"）。★**但乙单独当第一腿⛔推荐**：手写夹具是**自证**，测不出任何真实识别质量。

### ★★★ 票 15 必须拆（6 条 AC ＋ 3 个包 ＋ 3 处 `cmd/wisp`，一条腿做不完）

| 拆片 | 名字 | 范围 | 收哪条 AC | 碰 `cmd/wisp`？ |
|---|---|---|---|---|
| `15-a` | 引擎基座 | C9 Go 接口 ＋ session 生命周期 ＋ cgo 边界 ＋ `NumThreads:1` | AC6 ＋ AC4 机制 | ⛔ |
| `15-b` | VAD 端点 | silero ＋ 700/300/60s 状态机（WavInjector 进） | AC3 | ⛔ |
| `15-c` | 流式 ASR | paraformer int8 流式 ＋ partials ＋ final ＋ 纯静音→重试→「没听清」 | AC5 ＋ AC1 前半 | ⛔ |
| `15-d` | CER 量具 | 计算器 ＋ 基线集格式 ＋ 聚合报告 ＋ CI 接线 ＋ **语料择定** | AC2 | ⛔（只加 `scripts/`） |
| `15-e` | 装配与端到端 | `models.go` 的 `handOffModel` 调用、`resident_audio_windows.go` 帧消费、`resident_ball_windows.go` 唤醒暂停 ＋ 球走态 | AC1 后半 ＋ **AC4 的实测读数** | ★**是（三处）** |

★**AC4 只能"机制在 15-a、读数在 15-e"**：slot 互斥要 ASR **和 TTS** 争，而 TTS 是票 **26** 的 ⇒ **15 自己凑不出 AC4 的实测**。★这是**又一条未申报的跨票依赖**（与票 12 同性质）。

---

## ⑥ 技术前提逐条（写死了没有／在哪行／谁来定）

| 前提 | 写死了？ | 在哪 | 谁定 |
|---|---|---|---|
| 采样率 **16k/mono/int16** | ★**写死** | `SPEC-04:33`；代码 `internal/audio/audio.go:43 TargetRate=16000` | 已定，⛔ 可议 |
| 帧长 **512@16k＝32ms** | ★**写死** | `SPEC-04:42`；`audio.go:47 FrameSamples=512` | 已定 |
| **流式 or 整段** | ★**两条都写、⛔ 选哪条** | `SPEC-04:55`（流式 Paraformer／离线 SenseVoice）；预算分档 `SPEC-04:114`（≤800ms 流式／≤1500ms 离线） | ★**我定：流式 Paraformer 为主**——AC1 的球延迟要 partials，离线满足不了 ≤800ms；离线留作回退 |
| **VAD 用谁的** | ★**写死** | `SPEC-04:53` silero-vad（MIT，~1M，~8MB）；模型已在 `third_party/spike-models/silero_vad.onnx` | 已定 |
| 判停 **~700ms** | ★**半写死**：只写"【SPEC 默认】"，**无配置键** | `SPEC-04:53`；票 15:20 说"configurable"，但 `internal/config/schema.go` 的 `ASRConfig` **只有 `provider`/`model` 两个字段**（`grep -i 'silence\|vad'`＝0 命中）⇒ "configurable"**没有载体** | ★**我定**：先做包内常量；**加 `[voice].asr.silence_ms` 键要碰 D36 配置模型，留到有 UI 需要时单开票**——⛔ 在 15 里顺手加键 |
| 最短语音 **≥300ms** | ★**写死** | `PLAN.md:3081`（**D43 #11**：`Listening`→`Thinking` 触发条件） | 已定 |
| 最长语句 **60s 自动判停** | ★**写死** | `SPEC-04:53`；票 15:21 | 已定 |
| `intra_op_num_threads=1` | ★**写死（强制）** | `SPEC-04:55`、D32；票 15:18；spike 代码已照做（`speech-baseline/main.go:133`） | 已定，⛔ 可降级 |
| **一个 onnxruntime 实例**复用全链路 | ★**写死** | `SPEC-04:59`、D5 | 已定 |
| 模型**大小** | ★估算写死、**明标"待实测"** | `SPEC-04:53-57`（VAD ~8MB、KWS ~25MB、标点 100–200MB、TTS 120–300MB） | 15 落地时实测回填 |
| 模型**放哪** | ★**写死** | `SPEC-04:104` `<store>/<id>/`；`build/models/manifest.json` 已在；spike 副本在 `third_party/spike-models/` | 已定 |
| **Path C 全双工**（ASR+TTS+AEC 共存） | ★**写死** | `SPEC-04:63-64`、D47；票 15:29-33 | 已定，★但 **AEC 是票 59（`ready-for-agent`、未开工）⇒ 15 只留接口位** |
| **CER 6%/15% 双门禁** | ★**写死** | `SPEC-10:57-60`、`SPEC-04:125` | 已定 |
| ★**基线集从哪来** | ⛔ **全仓没写** | 零 `.wav`、无 `testdata/asr-baseline/` | ★**我定**：用**公开许可的普通话 16k 语料子集**造两套集，**⛔ 不用机主自己的录音**（隐私＋不可复现）；★**许可须逐集核实**——`14` 的 matcha **P3 BLOCKED 非商用**就是先例 |
| ★**ASR 选 Paraformer 还是 SenseVoice** | ⛔ **没写死** | `SPEC-04:55` 并列 | ★**我定**：Paraformer 流式（理由同上）；SenseVoice 留作离线回退 |

---

## ⑦ 交回编排者：五条要动的判断

1. ★★**"34 张下游"改 17**（直接 5＋传递 17，全 alive）。改方案 §3-2 与 §5 第 1 项的那句数。
2. ★★★**"不碰共享编译面"不成立**：`cmd/wisp` 三处（`models.go:25`／`resident_audio_windows.go:57`／`resident_ball_windows.go:522`）是**生产注释逐字点的名**，且 AC1 与之互斥。⇒ **15 的收尾片必须按 `cmd/wisp` 共享编译面排队**。
3. ★★**两条未申报的跨票前置**：**票 12**（`review`，AC1 的 text path）＋ **票 26**（AC4 要 TTS 才凑得出实测）。派单前要么在票面补 `Blocked by`，要么在拆片里显式划界。
4. ★★**"量具票 CER 自包含"只对一半**：计算器自包含，**基线集不自包含**（零 wav、零参考文本、零 CI 步）。这是 15 最被低估的一片，**须单列 `15-d`**。
5. ★**"前置全绿"在 `-done`＋0 未勾这把尺上成立（3/3）** —— 但 `02` 过给 15 的两笔债（`SetMemoryLimit`/GOGC 调参、票 15 的 `Spec refs` 从未修）要跟着票面走，否则"全绿"会被读成"无债"。
