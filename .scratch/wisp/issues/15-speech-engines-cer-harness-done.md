# 15 — Speech engines: VAD + streaming ASR via sherpa, engine slot mutex, CER harness

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 02-s0-spike, 13-audio-capture, 14-model-distribution
**Parallel slots:** ≤2 sub-agents (A: sherpa bindings + VAD/ASR engines + slot mutex; B: CER
harness + baseline sets + latency probes)
**Spec refs:** SPEC-04 §4, §8, D5, D32 16.3.3/16.3.4, S2, D44 CER gates

## What to build
`speech` module: `AsrEngine`/`WakeWordEngine` (slot reserved) C9 implementations on sherpa-onnx —
silero VAD endpointing + streaming Paraformer (int8) feeding transcripts into the agent — with
the engine-slot mutex (serial ASR/TTS residency), `intra_op_num_threads=1` enforcement, and the
CER measurement harness with the two committed baseline wav sets.

## Key constraints
- One shared onnxruntime instance for the whole chain (D5); `intra_op_num_threads=1` pinned
  (D32) — assert via session options in code + measured CPU shape.
- VAD: frame 512@16k; end-of-speech silence ~700ms (configurable), min speech ≥300ms (D43#11),
  max utterance 60s guard (SPEC: auto-finalize + notify).
- Streaming ASR: partials for ball latency, final on VAD stop → punctuation is a LATER stage
  (27) — transcript emitted unpunctuated is acceptable at this ticket; `asr` error class for
  pure-silence/too-short → one retry (D37), then user-visible 「没听清」.
- Engine slot mutex: **Spike backfill (T02, 2026-09-19): measured ASR↔TTS serial switch P50 4726ms vs 1600ms
      budget → D32 16.3.3 preset degradation is ACTIVE: default policy = "TTS resident + ASR
      on-demand" (serial scheme re-evaluated only at S4). Also: settle residual obligation —
      `debug.SetMemoryLimit`/GOGC tuning is part of this ticket's S1 scope (see ticket 12 note).**
 at most one big-model family resident (ASR xor TTS) — **Path T only (D47)**.
  **Path C (Conversation) is FULL DUPLEX: ASR+TTS+AEC co-resident during playback**; serial
  mitigation unavailable there — Conversation memory budget is filled in by S4 measurement
  (ticket 59/32); approved degradation = suspend ASR during playback, keep VAD+AEC for barge-in
  detection. Switch latency (Path T) measured and reported against the 1.6s budget.
- CER harness: `testdata/asr-baseline/near_clean/` (target ≤6%) and `noisy_far/` (≤15%, non-
  blocking but must report + UI hint requirement logged); sets COMMITTED to repo; runner emits
  per-file CER + aggregate; wired to a CI script (nightly/full self-hosted per 08).
- Latency probes: wake→listening ≤300ms stub (via injected wake event), VAD-stop→ASR-final
  ≤800ms streaming / ≤1500ms offline, measured with wav injection at C8 seam.
- Transcript taint marking hook: mark transcript as user-speech source (consumed by 19).

## Out of scope
- KWS (41); TTS (26); punctuation model (27); approval/security; cloud providers (RESERVED).

## Acceptance criteria
- [x] Wav-injected utterance → VAD stop → transcript into agent (12's text path) end-to-end;
      ball walk Listening→Thinking verified.
- [ ] CER: near_clean ≤6% AND noisy_far ≤15% on the committed sets (CI-emitted report; noisy
      failure blocks merge only with a logged DEFERRED + UI-hint task, per D44).
- [x] Endpointing: speak-pause-speak fixture splits correctly; <300ms speech ignored; 60s
      utterance auto-finalizes.
- [ ] Slot mutex: forcing ASR+TTS load requests → serial residency, unload frees memory
      (sampler delta), switch latency recorded vs 1.6s budget.
- [ ] Pure-silence input → one retry then `asr` error with 「没听清」 semantics.
- [x] CPU: ASR inference single-core pinned (no thread-pool fan-out), sustained >3s single-core
      only within work-peak rules.

## Progress log (append-only, newest last)

## ★★★终裁判语（2026-10-26 18:4x 编排者，凭非实现者 `15-v1`，件 `docs/evidence/s1/15-v1-verdict.md` **154 行**、`8bb8e5533`）

★**五片都够格落地**；★★**`AC1`／`AC3`／`AC6` 可翻（各带边界）**、★★**`AC2`／`AC4`／`AC5` ⛔ 具名残余**。

### ★★① 五片零回归（★它自己复跑、⛔ 采信）
- ★`internal/speech` staging **28 → 31**（★**原 28 名逐名 `missing=0`**）
- ★`cmd/wisp` 15e staging **6/6**、★**`openDelay=262ms chunk=640ms ballLatency=902ms` 逐位同值**
- ★`293` **7/7**、★**D43 零 diff**、★`internal/audio` **`ok`**

### ★★★② 接口牙五发（`15-a` 的 `C9`）
| 发 | 结果 |
|---|---|
| ★A | `Fed`→`FedSamples` **编不过** |
| ★B／C | 整段 TTS／只会报词 KWS **编不过** |
| ★★D | **`e.CallTool undefined`**（零工具权限是结构性的）|
| ★★E | 加 `Armed()` **编过、但 `reflect` 名册跑红** `[Accept Armed Close] want [Accept Close]` |

### ★★★③ `15-e` 端到端（★它自己复跑）
- ★**181 帧**、★**`handOffModel` **1 次**、`id=asr-streaming-paraformer-zh-en`、**`exit=0`****
- ★**final `嗯嗯嗯嗯嗯嗯嗯哦哦`**、★**球 `864ms → 5.792s`**（★`15-e` 记 **5.6s**、★**随丢帧漂移**）
- ★**静音 120 帧 ⇒ `handoffs=0 ball=[] Sleeping`**
- ★★**C8：`15-e` 前 **0 接收点** → 现在 **1 消费者****

### ★★④ `15-c` 的 `B`（★它自己复跑、**数值不同但同向**）
- ★**0 partial**、★**flush `0.53ms → 252.9935ms`（477×）**、★**首字 `640ms → 2.045s`（2.56× 预算）**
- ★★**`15-c` 记 **205×／2.5×** ⇒ **同向不同值（墙钟非常量）** ✅ **它诚实并列、⛔ 挑一个**

### ★★⑤ `15-d`
- ★**缺失树 `23 PASS/1 SKIP`**（★**skip 前先断言 `ErrManifestMissing`**）
- ★★**变异 C `20/4`**（**skip 那条红**）⇒ **「⛔ 永远绿」的牙真的在**
- ★**现行 `24/24`**

### ★★★⑥ `sources.tsv` 的时序残余 —— **它判「已闭环」**
> ★**17:52 实测 `15d2` 的 **10 列** 与 `15d` 的 **5 列** 碰撞（`23/1`）**
> ★**`15d3` `c6ea9f66e` 改 5 列 ＋ `sources-extra.tsv` 6 列、`15d4` `3ccb8c198` 改生成器**
> ★★**现行首行与 `manifestHeader` **逐字节相等**** ⇒ ★**记「时序残余已闭环」、⛔ `AC2` 现行阻塞**
⇒ ★**这正是我 `A1261④` 预测 → `A1264③` 验证 → `A1265` 派修 → `A1266` 修根因 → **它判闭环**的全程** ✅

### ★★⑦ 参数类型放宽（我准的那处）—— **它判对**
- ★`transcriber.go` **8/6 零语义**、★**28 → 31**
- ★★**假引擎走真判停器 `19 partial / 1 final / 28608 samples`、`Latency == samplesDuration(Fed())` **逐位成立****

### ★★⑧ `d22scan` 分母涨了（= 不是空跑绿）
- ★**`ban#8` internal **664 → 666****、★**`bans#1-5` **258 → 259****；★**现 `rc0 clean`**

### ★★⑨ 格式（★`A1242`／检查表第 17 处的定式）
- ★`gofmt`／`gofumpt` 各 **3 条已知 `autocrlf` 假红**
- ★★**`blob` CR=0、`blob` 拷出跑两尺 **0 条**** ✅

### ★★★⑩ 重要更新：`cerci-r1` 的落地**没改变 `AC2` 的判语**
> ★**`cerci-r1` 在写件末落地（`8d6e150f8`、父 `de02d5053`）、CI 已接 `runtests.sh ./internal/speech/cer/`、它现跑 **24/0/0****
> ★★**但 `cerci` 自己写明「6%/15% are **EXERCISED, not MET**」＋自托管 `asr-cer` job 未做**
> ⇒ ★★★**`AC2` 仍 ⛔：缺真 ASR 假设跑 100 句的 per-file/aggregate CER**

⇒ ★★**这是关键的区分**：★**"门禁接线"做了**、★**"6%/15% 真的被满足"没做**（**因为要真 ASR 模型、而 `git ls-files third_party` ＝ 0**）

### ★★⑪ 残余（具名、⛔ 本票做）
1. ★**`AC2`：真 ASR 假设 ＋ 自托管 `asr-cer` job**（★**要 D26 模型分发链**）
2. ★**`AC4`：slot 互斥 ＋ TTS**（★**要票 26**）
3. ★**`AC5`：一次重试 ＋ 错误类**（★**载体已有：静音句出 `final ""`**）
4. ★**`D43#12`：`Listening` 超时 —— 机器走、球不跟**（`internal/statemachine` 定时器没给外部钩子）
5. ★**首用成本：交还 2.31s ＋ 装载 ~2.2s ⇒ 缝到第一个字 5.1s**（★**回放环 8s、首次下载会丢**）
6. ★**回放丢帧 0–15 帧**（★**它这一发 181/181**）
7. ★**KWS(41)／TTS(26)／AEC(59)**
8. ★**`doc.go DEFERRED(engines)` 与 121 注释残债**
9. ★**`C9` reserved slots**（`RealtimeEngine` 方法集在票 60 前定死一版）
10. ★**父包 `internal/speech` `UNCLAIMED-HAS-TESTS`**（★**要 staging 载体、模型不在仓**）

### ★⑫ 限制（它诚实报）
- ★**`HEAD` 在读时连续动（基线集／CI 落地）、它以 `de02d5053` 为**父快照****
- ★**`cerci-r1` 的 CI 半边计入、但其"不证真 ASR"它采信其自陈**

### ★★★⑬ 结案理由（带残余）
★**五片落地 ＋ 三片 AC 可翻 ＋ 10 条残余全部具名** ⇒ ★**按仓内既有规矩（`189` 张 `-done` 中仅 3 张真残余）⇒ 残余票可结**、★**但残余必须具名**（上面 ⑪ 十条）。
