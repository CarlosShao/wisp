# 28 — SessionScope (C31): Warm/Conversation keep-alive, timers, privacy gating

**Status:** **DONE-带残余（2026-10-26 20:3x 编排者收口：5 格已翻；AC3 归票 59 ＋ R28-2 mic 是唤醒闩不是喂食门 ＋ R28-3 断电器档 ＋ R28-4 PrivacyConfirm=nil ＝ 具名残余见末节）** ~~ready-for-agent~~
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 07-ball-state-machine-core, 15-speech-engines-cer-harness
**Parallel slots:** ≤2 sub-agents (A: SessionScope + model refcount + timers; B: Conversation
mode + privacy confirm + battery shrink)
**Spec refs:** SPEC-05 §5, B4, C31, D43 #26,28–31, §15#5 decision, D42#6, S4

## What to build
The `session` module: the single owner of keep-alive — session definition, Warm/Conversation
lifecycles, model reference counting, panel keep-alive reservation, session-scoped grant
registry (rows now, lifecycle at 49), and the settle/dispose cascade.

## Key constraints
- Session = wake → explicit end OR 90s idle (`[session] warm_timeout_sec`); `Settling` 3s static
  period; dispose cascade per D43 #31 (models unload, panel hide-or-destroy, FreeOSMemory,
  grants expired).
- `Warm`: models hot + mic CLOSED + faint warm breathing; next wake via click/hotkey → Listening
  with ZERO model load (assert no engine loads on re-wake — B4's whole point).
- `Conversation` (default OFF via `[voice] conversation_mode`): **Path C full duplex (D47) —
  mic stays OPEN through playback with AEC (ticket 59/P15-gated); barge-in mid-playback →
  Listening (D43 #41)**; red solid ring never fades; Speaking-end → Listening directly (#28);
  30s no-speech → Warm (close mic FIRST); first-ever enable → L2-grade privacy confirm
  (native card) + logged. **Memory note: ASR+TTS co-residency breaks the Path-T serial
  mitigation — Conversation budget measured in S4; approved degradation = suspend ASR during
  playback, keep VAD+AEC for barge-in (ticket 59/15).**
- Battery saver (D42#6): Warm 30s / Conversation 15s; KWS resident OFF.
- Model refcount: ASR/TTS/VAD/punct/KWS acquire-release through SessionScope only (no other
  module may keep models alive); leak test: session end → all refcounts zero.
- Async L1 extraction + history compression run inside the Warm window (hooks for 29/10).
- Interrupt/settle interplay: new wake during Settling cancels settle (#33); explicit end
  (hotkey long-press/menu) → immediate dispose.

## Out of scope
- Panel keep-alive UX (33 wires its window into scope); grant semantics (49); extraction logic (29).

## Acceptance criteria
- [x] Cold vs Warm wake: second-round first-token latency measured; Warm wake shows zero engine
      loads (instrumentation counter) — B4 evidence.
- [x] Timers: 90s Warm → Settling → dispose; 30s Conversation → Warm with mic closed FIRST
      (order asserted); battery-saver values applied.
- [ ] **Path C full duplex (D47, blocked by 59): during Conversation playback mic stays open
      (AEC), barge-in works (≤400ms), Conversation memory measured and budget row backfilled
      (or degradation mode verified if over budget).**
- [x] Conversation first-enable → privacy confirm + log; subsequent enables no prompt.
- [x] Dispose cascade order test (models → panel-hide → FreeOSMemory → grants expiry) and RSS
      return ≤10s (sampler).
- [x] Refcount leak test: 100 wake/sleep cycles → zero residual refs/goroutines.

## Progress log (append-only, newest last)

## ★★★终裁判语（2026-10-26 20:3x 编排者，凭非实现者 `tri-v1`，件 `docs/evidence/s1/tri-v1-verdict.md` **134 行**）

★**可翻 5 格、`AC3` 归票 59** ⇒ ★**带具名残余结案**（★先例：`15-done` 自己就是 3 格未勾结案）。

### ★★★① 五格逐条（★全在 `wire-r1` 装配后复跑）
| 格 | 判 | 读数／边界 |
|---|---|---|
| ★`AC1` Cold vs Warm | ✅ **翻** | ★**`cold_acquire=2.9790732s` → 第二回合 `acquire=0s`**、★`loads_after_warm=0`、★`loader_calls=1`（★**边界：`first-token` 用 `28-r1` 的带 tag 读数**）|
| ★`AC2` Timers | ✅ **翻** | ★**89s Warm／90s Settling／93s Disposed**、★**Conversation 30s 落 Warm 且 `mic.close` 在 `SCOPE-WARM` 之前**（★★**但见 ② 的 `R28-2`**）|
| ★`AC4` 隐私门 | ✅ **翻** | ★两次 `enable` ⇒ `prompts=1`、★`CONVERSATION-PRIVACY-CONFIRM granted` **恰 1 行** |
| ★`AC5` dispose 级联 | ✅ **翻** | ★序 `engine.close:*` ＜ `panel.hide` ＜ `free-os-memory` ＜ `grants.expire`（**两条 teardown 路径都过**）；★`samples=0`（8.58ms）|
| ★`AC6` refcount | ✅ **翻** | ★**100 轮、`loader calls=200`、每类 `closes=100`、每轮 `Refs=0/Hot=0`、`goroutines base=2 after=2`** |

### ★★★② `AC3` ⛔ 归票 59（`R28-1`）
★**`Path C full duplex (D47, blocked by 59)`** ⇒ ★**`59-p15-aec-spike.md` 现读 `ready-for-agent`、无 `-done`、7 格 0 勾**

### ★★★③ 而 `tri-v1` 纠正了 `wire-r1` 一处**性质写错**（要写进票面 `Warm` 定义旁）
> ★★**`stages=[Awake Awake]` 只证**序**（钩子先于翻面）**
> ★★★**事实是 `warmFeed` **唤醒闩**、⛔ 喂食门** —— ★**`run()` 仍每帧 `tr.Accept`（无 `warmFeed` 判断）、设备未关、门未驱动**
> ⇒ ★★★**`Warm="mic CLOSED"` **未实现**、**票 41 前无法实现****
> ⇒ ★★**所以 `AC2` 的 "order asserted" 成立、但 `Key constraint` **不成立****

### ★★④ 其它具名残余
- ★**`R28-3` 断电器档未选**
- ★**`R28-4` 生产 `PrivacyConfirm=nil` 开不了**
- ★**`R28-5` 合成载荷／播报后重装**
- ★**Conversation UI 半边归 33**

### ★★★⑤ 反向判据（`tri-v1` 自造、加了 `M3b`）
★`M3` 回合结束 `closeEpoch+asrGone` ⇒ `rc1` `stage=Disposed`、factory 2 次
★★★**`M3b` 不置 `asrGone`（让唤醒付重装）⇒ `rc1`（`warm wake took its ASR borrow back acquire=3.0484092s` vs 绿发 `0s`）** ⇒ ★**"必须变慢"成立**
★**还原核在最后一发变异之后**、★无 overlay 重编 `3/3 PASS`、★源件 sha256 同值、★`MUTANT` 0 命中

### ★★⑥ 零回归（`tri-v1` 复跑）
★**15e `6/6`**、★**`262/640/902` 逐位同值**、★**`drained=181/181`**（★另一发 `179/181` 当抖动）；★`session 42/0/0`；★**`d22scan` 导出树 `rc0 clean`**
