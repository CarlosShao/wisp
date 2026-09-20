# 21 — Approval gates minimal: L1 block window, L2 native card, queue trivial, batch D45-1

**Status:** segment-1-landed-decision-layer (L1 window + veto channels + L2 queue + source check + batch + applied-steps type); segment 2 (native card rendering, ball visuals, loop rewiring) still open
**Claimed by:** sub-agent A (segment 1: decision + queue layer, UI injected). Segment 2 unclaimed.
**Last update:** 2026-09-20
**Blocked by:** 17-risk-assessor-c19
**Parallel slots:** ≤2 sub-agents (A: L1 window + veto channels; B: L2 native card + queue +
batch aggregation)
**Spec refs:** SPEC-06 §2, §7, §8, D4, D31, D45-1, B1, C18 (trivial impl), 16.9#5

## What to build
Single-task approval mechanics with the simplest viable UI (native, pre-panel): L1 pre-execution
block window (2–3s countdown strip) with its four veto channels, L2 native confirm card
(full params + rules_hit + reason; allow ONLY from native side), a single-task ApprovalQueue
(C18 semantics at trivial scale), and batch aggregation for homogeneous L1 operations.

## Key constraints
- L1 window (B1-corrected): 2–3s pre-execution BLOCK (not "undo"); countdown ring on ball strip;
  veto = click ball / Esc hotkey / panel reject (no-op until 37) / KWS veto word (channel
  present, functional at 41; when KWS not loaded UI MUST say 「语音取消不可用」 — never fake it).
  Timeout → execute; veto → call cancelled, failure returned to LLM.
- L2: native confirm card (C27 fallback-grade): tool name + L2 badge + FULL params (not summary)
  + rules_hit reasons + reject button + "click ball to approve" guide. **Allow action accepted
  ONLY from native sources: ball click / native card button / global hotkey. Panel-sourced allow
  is structurally rejected at the server-side API (F2 layer 3) — implement the source check
  NOW even before panel exists.**
- Queue (single-task trivial): one pending item; correlation_id assigned; **timeout 300s →
  auto-REJECT** (never infinite wait); last-30s prominent warning; reject → task root ctx NOT
  cancelled; one-click replay available. Host-unreachable/queue-broken → fail-closed.
- Batch aggregation (D45-1): ≥3 homogeneous L1 ops in ONE tool call → single confirm (total +
  affected-dirs summary + first 5 + expandable); L2 NEVER aggregated; R7 (≥50 files) upgrades
  the batch to L2 (unaggregated).
- Cancel non-atomicity surfaced: cancel/veto after start → applied-steps report (from 20).
- Ball visuals: Confirming pulse + AwaitingApproval depth badge (depth=1 now).
- **D47 scope note: all approval mechanics live on the TEXT loop (Path T + handoff). Path C
  realtime brain has ZERO tool permissions — it cannot produce approval requests; if the user
  asks it to "do X", the C32 handoff returns the task to the text loop where full gating applies.**

## Out of scope
- Multi-task queue routing (48); panel-side card (37); session grants (49); input.type binding
  (23 + 49).

## Acceptance criteria
- [ ] L1: countdown executes on timeout; each working veto channel cancels; non-loaded-KWS
      message asserted; veto returns failure to LLM (loop continues).
- [ ] L2: card shows full params + rules_hit; native allow executes; simulated panel-source
      allow REJECTED at API (test with forged source); 300s timeout → reject + warning at 270s
      + replay works.
- [ ] Batch: 10 L1 ops → one confirm; L2 mixed in → no aggregation; 50+ → L2 (R7).
- [ ] Cancel-mid-execution → applied-steps report rendered.
- [ ] State machine: Confirming/AwaitingApproval transitions match D43 #17,21–24.

## Progress log (append-only, newest last)

- 2026-09-20 (segment 1, commits 35200c7 + b1535d7): did=the decision and queue
  layer now exists as real code in internal/agent/approval behind the
  internal/tools/gate.go seam - L1 2-3s block window (four veto channels, each
  reported with its true availability, 「语音取消不可用」 asserted verbatim for
  unloaded KWS, and a veto there neither cancels nor stays silent); C18 single
  task queue (correlation_id assigned, 300s auto-REJECT proven to resolve, one
  warning at the 270s mark and none before, task root ctx proven alive after a
  reject, replay re-displays under a FRESH nonce, host-unreachable and
  full-queue both fail closed); D45-1 aggregation (10 ops in one call = one
  confirm whose veto stops the whole batch; <3 / L2 / R4-taint / R3-sensitive
  never aggregate; a 60-target batch escalates to L2 unaggregated and its
  timeout REJECTS instead of executing); the native-source proof is a
  crypto/rand single-use nonce minted at display time and bound to
  (corr,task,tool,level,args,seq) - Source is logged and never read, and five
  forgery cases are refused, including a panel call that claims Source=native
  while carrying the real live grant (which additionally revokes that item's
  nonces as a leak signal); D31's CancellationReport + Bus are typed and
  tested with a fake tool, including the "tool reported no steps" worst case;
  D47 is enforced as data (AdmitTextTask registration, unadmitted task refused
  before the queue is touched - asserted with the queue depth). Gates: gofmt
  clean, go vet clean, go test -count=2 ok, go test -race ok, d22scan reports
  zero findings in this package (its one repo-wide finding is pre-existing in
  internal/llm/adaptertest/mockllm.go:68).
  NOT ticked: every AC box. Two reasons, both honest. (a) No production
  reachability (R12/A11 shape): internal/tools has ZERO importers outside
  itself, so neither the bridge nor this gate is composed under a Loop in
  cmd/wisp - the L1/L2 behaviour is proven at the bridge (the real
  agent.ToolProvider path), not in the running product. (b) AC#2's "card shows
  full params + rules_hit" is proven as delivered DATA (Prompt.Params/RulesHit
  asserted through the bridge); the card itself is segment 2.
  Still missing for the fs.write unblock: agent.Loop.decideRisk still refuses a
  DECLARED L1/L2 before Execute, so fs.write dies at that pre-gate today. This
  layer makes the fix a deletion - loop.go:719-738's RiskL1/RiskL2 branch (and
  the loop's own pre-gate tool_call booking, which is strictly less
  informative than bridge.book) - plus ticket 12 composing approval.New as
  Options.Gate and calling AdmitTextTask per task. No change to decideRisk was
  in this segment's scope.
  next=segment 2: native card + ball Confirming pulse/AwaitingApproval depth
  badge implementing approval.UI; wire Gate.Bus()/Result.AppliedSteps into the
  fs.write family (ticket 20 segment 2); ticket 12 composition (Options.Gate +
  AdmitTextTask + Native()/Panel() hand-out) and the decideRisk deletion; ticket
  41 flips ChannelRegistry.SetLoaded(ChannelKWS,true) once the model is
  resident; ticket 37 gets Queue.Panel() and nothing more. State machine D43
  #17,21-24 AC untouched here.
