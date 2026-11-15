# 32 — S4 acceptance: scenario ③ dictation + scenario ④ chat companion gates

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 27-punctuation, 29-memory-l1l2, 30-result-routing-d10, 31-reminders
**Parallel slots:** ≤2 sub-agents (A: scenario ③; B: scenario ④)
**Spec refs:** §10 #12 #13, D44(b) rows ③④, P4/P7 conclusions, S4 done

## What to build
The S4 slice gate executing the two human-experience scenarios end-to-end on real hardware:
③ dictation-to-text-in-place and ④ multi-turn voice companion, each with the exact PASS bars
from the verification plan.

## Key constraints
- Scenario ③ dictation: continuous 300-char dictation → text appears AT THE TARGET APP's cursor
  via input.type (saving to a file does NOT count as pass); punctuation subjective ≥8/10;
  plus one elevated-window attempt → explicit UIPI message (§16.11#2).
- Scenario ④ companion: 5 consecutive conversation rounds; round ≥2 first-char ≤1.5s (model
  excluded — Warm alive, evidenced by round1-vs-round2 delta >1s); TTS subjective ≥7/10 (P7
  sheet); remembers one fact from round 1 (L1 profile live). **D47 additions: barge-in gate —
  mid-playback speech stops TTS ≤400ms and never transcribes our own playback; external-audio
  case (music playing in background) false-trigger rate recorded against P15 thresholds;
  Conversation ASR+TTS co-residency memory numbers backfilled into SLO table (or degradation
  mode verified).**
- Warm/Conversation behavior observed: default Warm (mic closed, warm breathing); explicit
  Conversation enable (privacy confirm once, red ring); 30s silence → Warm.
- Result routing during ④: mixed short/long replies hit correct D10 tiers.
- Evidence: screen recordings + logs + tool traces committed under docs/evidence/s4/; human
  (user) sign-off required for both subjective scores (D29 rule: human acceptance is the bar).

## Out of scope
- Panel UI (results show via notify/ball/native primitives until 36); KWS.

## Acceptance criteria
- [ ] Scenario ③ PASS: video + transcript + input.type trace; elevated-window failure message
      captured; score sheet ≥8/10.
- [ ] Scenario ④ PASS: 5-round log; latency table (round1 vs round2 delta >1s proving Warm);
      TTS ≥7/10; fact-recall demonstrated. **D47: barge-in ≤400ms gate + no-self-transcription
      PASS; external-audio false-trigger evidence; Conversation memory row backfilled.**
- [ ] Conversation privacy-confirm-once + red-ring visuals verified by user.
- [ ] D10 tier hits during ④ consistent with 30's matrix (no regressions).
- [ ] Pre-mortem: TTS device yanked mid-speak; mic device yanked mid-dictation — correct error
      classes + no stuck states.

## Progress log (append-only, newest last)
---

## S4 验收结果（2026-11-15 12:1x，`cf79bc9b6`）：**满足 0／总 5**，五格一格不勾，其中 **1 格不可验**

**尺1**：`cf79bc9b6` 在，只带 `docs/evidence/s1/32-s4-acceptance.md` ＋ `.scratch/wisp/probes/32/readings.txt` 两条路径（`2 files changed, 335 insertions`），别腿在飞件零触碰。
**尺2**：证据件 **28103 B／175 行**；底片 `readings.txt` 160 行（每条命令原文＋rc）。
**尺3**：★**它两条最重的发现我逐条 grep 复核，全部成立**（见下）。

### 逐格判语（**一格不勾**）
| 格 | 判语 | 一句话 |
|---|---|---|
| AC#1 场景③ | 部分 | `input.type` 真机 300 字（sha256 逐字同）＋提权拒绝修后真机捕获都在；但**没有"口述→光标"这条产品路**，300 字那发是探针直调 tool，8/10 是文本夹具自评 |
| AC#2 场景④ | 部分 | Warm delta 3.24s、零自激转写、P15 外部音源 0/10min 都在；五轮日志零字节、TTS 7/10 无人评、barge-in 不在默认构建、SLO 内存行仍「未测」 |
| AC#3 隐私＋红环 | **不可验** | 代码＋测试齐（confirm 只弹一次、红环 Danger 2px），但票面逐字要求「verified by user」⇒ **只能机主人眼签收** |
| AC#4 D10 分档 | 部分 | 矩阵＋四档真管线＋生产同一 Dispatcher 跑绿；缺「④ 期间」的观测 |
| AC#5 拔设备预演 | 部分 | 麦侧错误类/具名/文案有测试；但产品**不派发** `EvAudioDeviceLost`，TTS 播放中拔设备无夹具 |

### ★四条发现（腿报，**我复核成立**）—— 前两条已独立成票
1. **L1 画像注入在生产上是断的** → 立 **票 451**。**我核**：`Profiles\s*[:=]` 产码命中 **3 处**，只有 `internal/agent/loop.go:677` 是转发（`Profiles: l.opt.Config.Profiles`），另两处是无关串（错误串 / `L1MaxProfiles` 常量）⇒ **全仓零赋值点** ⇒ 画像写了库却永远进不了 prompt。
2. **barge-in 不在出厂构建里** ※ 与**拔麦不推状态**合并成 **票 452**。**我核**：`EvAudioDeviceLost` 产码命中 **4 处**，其中 `cmd/wisp/resident_audio_windows.go:39` 与 `:358` **是注释**，`:358` 逐字「这一条不推状态、也不拒绝启动」⇒ **零派发**。
3. **票 32 自己的结构缺口**：`Blocked by` 里的 **31-reminders 在本票 5 格里零格承担**（腿逐字核过 `:36-44`）⇒ **S4 切片验收替不了它签收**。⇒ 本票结案时要么把 31 从 `Blocked by` 摘掉，要么补第 6 格 —— **归我改票面**，未改。
4. **`internal/speech` 标点族 SKIP 是环境事实**：本机 System32 的 `onnxruntime 1.17.1` **遮蔽**仓内 `1.28.2` ⇒ 具名记，⛔ 当产品红。

### 为什么五格不勾
AC#1／#2／#4／#5 都"部分满足"，AC#3 **不可验**（要机主人眼）。按 `issues/README`，**部分满足＝不勾、不可验＝不勾**。⛔ 我不把"下层模块全绿"当切片验收通过。

### ⇒ 本票的真前置（不只是"41–45 已完成"）
1. **票 31-reminders 的判据**（结构缺口，要么摘要么补格）
2. **票 451**（画像注入断）
3. **票 452**（barge-in 出厂形态＋拔麦推状态）
4. **机主人眼一次**（AC#3 红环＋隐私卡）
5. AC#2 的 TTS「7/10」口径要么认 R5 云 TTS 作替代并写进票面，要么等机主双盲人评

⛔ 本票在上述完成前不许改名结案。
