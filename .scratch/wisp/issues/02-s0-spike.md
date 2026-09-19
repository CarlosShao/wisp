# 02 — S0②③④⑤: spike measurements, X/Y topology verdict, SLO backfill

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 01-build-chain
**Parallel slots:** ≤2 sub-agents — sub-agent A: window/shell side (baselines ④⑤ + go-webview2);
sub-agent B: speech side (baselines ①–③ + goja + model residency). Shared: verdict + docs.
**Spec refs:** D25, D32 16.3.2/16.3.3, D44 S0 revision, P1, P2, P11, S0 done②–⑤

## What to build
`scripts/spike/` measurement programs + a spike report that turns every "target value" in D32 into
a measured value, and issues the **irreversible X/Y topology verdict** (single-process cgo vs
speech sub-process). All later memory/latency acceptance numbers come from here.

## Key constraints (five outputs, all mandatory)
1. **Five RSS baselines** (private working set, split private vs shared):
   empty Go · +onnxruntime DLL mapped (no session) · +one sherpa session ·
   **+layered window +Direct2D +DirectWrite** · **+tray +global hotkey +Job Object**.
   (④⑤ are the real constituents of S1's 25MB budget — do not skip.)
2. **X/Y verdict with pinned decision rule**: idle private RSS ≤25MB → path Y; >25MB → path X.
   Inputs MUST include cgo-crash probability (C-side segfault cannot be recovered) and
   "can RSS return to idle within 10s after unload" (`debug.FreeOSMemory()` required).
   If X cannot settle in 10s → verdict Y. No discretion allowed.
3. **goja ES capability report**: async/await, ES modules, and `vm.Interrupt()` wall-clock
   precision (D33/F5's only resource constraint). Outcome decides Tier-2 language level or
   QuickJS-binding fallback (P2).
4. **go-webview2 cold/hot window latency**: create/destroy; cold ≤1500ms target (P11).
   If cold >2s → flag "L2 confirm card may need native re-evaluation" as a blocked decision,
   do not silently proceed.
5. **Model residency table + ASR↔TTS serial switch latency** (int8 models; fills D32 16.3.3;
   validates "peak = max(ASR,TTS) not sum" via half-duplex serial loading).

## Key constraints (docs)
- Backfill `docs/SLO.md`: measured per-state table (D32 16.3.2 rows), model residency table,
  latency cold/warm rows. `docs/PRECHECK.md` rows P1/P2/P11 get conclusions.
- Spike programs live under `scripts/spike/`, each runnable standalone, results JSON-emitted.

## Out of scope
- Any product feature; panel UI; agent loop. Spike harness code quality = throwaway OK
  (experiment branches allowed), but the report files are repo docs and frozen.

## Acceptance criteria
- [ ] All five outputs produced with raw numbers + machine (CPU/RAM/OS) recorded.
- [ ] X/Y verdict written into `docs/SLO.md` + `docs/PRECHECK.md` with the decision inputs shown.
- [ ] `docs/SLO.md` marks 700MB work-peak as "target" until model residency confirms it.
- [ ] goja ES + Interrupt precision conclusion recorded; go-webview2 cold/hot numbers recorded.
- [ ] Report committed and pushed; verdict cited in tickets 07/12/15 (edited into their spec-ref
      notes if needed).

## Progress log (append-only, newest last)

- [2026-09-19T08:05:00Z] agent=T02-impl did=spike-scaffold+model-download+baselines-1-to-5 (scripts/spike module, NtQSI private-WS sampler verified against Task Manager; five baseline JSONs in docs/evidence/s0/data: empty-go 6.9MB priv, +DLL 7.5MB priv/+5.1MB shared, +KWS-session 48.6MB priv, +D2D/DWrite window 16.0MB priv, +tray/hotkey/job 16.0MB priv no-leak) next=xy-verdict+goja+webview2+model-residency

- [2026-09-19T08:25:00Z] agent=T02-impl did=xy-verdict-inputs (idle-y 16.5MB priv [rule <=25MB], idle-x 15.8MB, ASR unload-test: load 3.3s, loaded 281MB, dispose+FreeOSMemory plateaus 22.3MB > 10s-settle FAIL -> D32 rule "X-cannot-settle=>Y" moot since memory already picks Y; residual implies post-session Y Sleeping ~31MB > 25MB cap, flagged for S1) + cgo-crash research (sherpa-onnx issues #2694 open Go SIGSEGV at session create, #3635 open VAD buffer overflow 37h, historical C-API segfaults fixed) next=goja+webview2+model-residency
