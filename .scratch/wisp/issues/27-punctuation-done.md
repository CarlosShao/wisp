# 27 — Punctuation restoration (P4): CT-Punc model, dictation-grade output

**Status:** **DONE-带残余（2026-10-26 20:3x 编排者收口：4 格已翻；AC5 无采样器 ＋ R27-2 生产未读开关 ＋ R27-3 SLO 未回填 ＋ 英文中文标点 = 4 条具名残余见末节）** ~~ready-for-agent~~
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 15-speech-engines-cer-harness
**Parallel slots:** 1
**Spec refs:** SPEC-04 §4, D32 16.3.4 (≤300ms), P4, S4

## What to build
Punctuation restoration stage in the speech pipeline (CT-Punc-class int8 model via sherpa),
feeding punctuated transcripts to the agent — the make-or-break for dictation (scenario ③).

## Key constraints
- Latency: ASR-final → punctuated text ≤300ms (D32); measured on committed fixture set.
- Model residency: joins the Warm-resident set; participates in slot mutex accounting; unload
  at session end.
- Streaming ASR partials remain unpunctuated (display-only); finalization always punctuates.
- Toggle `[voice] punctuation=true`; disabled → raw transcript marked as such downstream.
- Quality: subjective punctuation correctness ≥8/10 on the dictation fixture (300-char
  continuous speech) — this is the S4 scenario-③ dependency; below → blocked-decision (try
  alternative model size before deferring).
- Failure: model load fail → degrade to unpunctuated + `model` error log + UI hint, never fail
  the utterance.

## Out of scope
- Dictation E2E acceptance (32); input.type (23).

## Acceptance criteria
- [x] Punctuated E2E: wav → transcript with punctuation reaching agent (mock LLM asserts
      punctuation present).
- [x] Latency ≤300ms P50/P95 on fixtures, recorded in SLO appendix.
- [x] 300-char dictation fixture subjective score ≥8/10 (evidence sheet) or blocked-decision
      logged with next-step.
- [x] Toggle off path works; load-failure degrade path works.
- [ ] Memory accounting: punct model included in Warm residency numbers (sampler).

## Progress log (append-only, newest last)

## ★★★终裁判语（2026-10-26 20:3x 编排者，凭非实现者 `tri-v1`，件 `docs/evidence/s1/tri-v1-verdict.md` **134 行**）

★**可翻 4 格、`AC5` 留残余** ⇒ ★**带具名残余结案**（★先例：`15-done` 自己就是 3 格未勾结案）。

### ★★① 四格逐条
| 格 | 判 | 边界（`tri-v1` 具名）|
|---|---|---|
| ★`AC1` `Punctuated E2E` | ✅ **翻** | ★**sink 边界**：断言到"到达 sink 的 `Transcript`"、★**进 agent 那一跳在 `cmd/wisp`** |
| ★`AC2` `Latency ≤300ms` | ✅ **翻** | ★**`P50=9.255ms`／`P95=19.587ms`／`max=21.9ms`**（80 次）⇒ **15.3× 余量**；★**`SLO` 附录未回填**（★**⛔ 伪称已进 `SLO`**）|
| ★`AC3` 300 字主观 | ✅ **翻** | ★**302 字 ⇒ 26 处标点、未截断**、★**8/10**（漏一处逗号）|
| ★`AC4` toggle/degrade | ✅ **翻** | ★两条都实测 |

### ★★② `AC5` ⛔ 留（具名残余 `R27-1`）
★**`Memory accounting: punct model included in Warm residency numbers (sampler)`** ⇒ ★**没有 RSS 采样器**（★采样器在 `cmd/wisp`／`internal/models`、本片禁写面）

### ★★★③ 另外两条残余（`tri-v1` 现核、⛔ 硬凑）
- ★★**`R27-2` 生产未读开关**：★`[voice] punctuation` 键**已在 `internal/config/schema.go:252`** ⇒ ★**缺的是 `cmd/wisp` 读它**（`OpenPunctuation(...,true)` **硬编码**）
- ★★**`R27-3` `docs/SLO.md:77` 未回填** ⇒ ★★**而 `27-r1` 说"SLO.md 不存在"是**假****（★**文件在、行 77 仍"未测→待票27"**、★行 105 "待票28"）
- ★**英文口述拿到中文标点**（★`punc_list` 只有 `。，？、`）

### ★★④ 反向判据（`tri-v1` 自造）
★`M2` 摘 `SetPunctuator` ⇒ **`rc1`**（★`the final reached the agent UNPUNCTUATED … SetPunctuator is not on the path`）；★**还原核在最后一发变异之后**、★无 overlay 重编 **`3/3 PASS`**、★源件 sha256 同值

### ★⑤ 零回归（`tri-v1` 复跑）
★`wire 7/7`（合跑 `13/13`）／★`speech staged 55/0/0`／★`session 42/0/0`；★**`d22scan` 导出树 `rc0 clean`** ＋ 自检 **40 向**全过

### ★★⑥ 建议承接票
★**开一张"采样器 ＋ `SLO` 回填 ＋ `cmd/wisp` 读开关"小票**（★三件都是同一条装配尾巴）
