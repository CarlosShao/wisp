# 11 — LLM adapters: Anthropic Messages + OpenAI Responses + fallback chain

**Status:** in-progress
**Claimed by:** agent-ticket11-adapters
**Last update:** 2026-09-20T05:05Z
**Blocked by:** 09-llm-provider-openai-mockllm
**Parallel slots:** ≤2 sub-agents (A: Anthropic adapter + prompt caching; B: Responses adapter +
fallback semantics)
**Spec refs:** SPEC-05 §3.1, §3.4, D8, §14.2, D40#4, C5 cache-capability

## What to build
The remaining two protocol adapters behind the same C5 seam, plus provider `fallback_provider`
switching and the complete user-visible failure-semantics table (§14.2) enforced at seam level.

## Key constraints
- Anthropic `/v1/messages`: native tool use, extended thinking → ReasoningDelta, prompt caching
  via explicit `cache_control` breakpoints driven by the C5 cache-breakpoint capability
  (prefix sections ①⑤⑥② stable; suffix never cached). Usage.cached mapped from cache reads.
- OpenAI Responses `/v1/responses`: reasoning items → ReasoningDelta; tool-call semantics
  normalized to C6; stop-reason mapping table per adapter documented in code-adjacent docs.
- Adapter unit tests may use mockllm endpoints (add Responses/Anthropic canned modes) + recorded
  golden SSE from real providers where licensable/possible (commit sanitized recordings).
- Fallback: on `provider`-class errors after retry exhaustion → switch once to
  `[llm] fallback_provider`, surface switch in logs + state; second failure → Error(provider),
  task ctx preserved, resumable (D40#4).
- Failure semantics (§14.2) user-visible contract: retries show "重试中 (n/3)" after 5s; 401/403 →
  Unconfigured + guide; quota → distinct account-vs-network message; mid-stream drop → partial
  content + 「响应中断」 marker + unclosed tool calls failed.
- No provider-specific logic may leak above the C5 interface (compile-level separation; review
  checks imports).

- **Capability probe implementation (2026-09-19 supplement)**: real-request probes per capability
  — tool_calls (minimal fc round-trip), vision (1×1 image), thinking (reasoning field present),
  audio (only when a cloud voice provider exists, ticket 61); results → `provider_health` (SPEC-02
  v2); "declared ✓ / measured ✗" surfaced to panel (40). **Compat flags consumption**: `compat.loose`
  tolerates missing stream_options/usage/tool_choice per provider. **Rate-limit self-restraint**:
  per-provider `rpm/tpm` local token bucket (before the provider 429s us); coordinates with tool
  concurrency ≤4 (ticket 47).

## Out of scope
- Ball/panel UI rendering of these states (ball shows via state events; panel later); quota
  enforcement (44); voice providers (61).

## Acceptance criteria
- [ ] Both adapters pass the same golden/fault test suite as OpenAI Chat (identical harness).
- [x] Cache-breakpoint test: Anthropic request bodies show stable cached prefix across turns
      (capture via mockllm); cache-read usage populated.
- [x] ReasoningDelta surfaces for thinking responses on both adapters.
- [x] Fallback chain test: fail_next(3) on primary → fallback used → task completes; double
      failure → Error(provider), ctx preserved, resumable.
- [x] Failure-semantics matrix: each §14.2 row has a test asserting the emitted user-visible
      event/state.
- [ ] Probe suite: mockllm control endpoints emulate fc-capable / fc-broken / vision-capable
      providers → probe results land in provider_health with correct ✓/✗; declared-vs-measured
      mismatch event emitted.
- [ ] Token bucket: rpm=10 fixture → 11th request within 60s waits locally (no 429 from server).

## Progress log (append-only, newest last)
- [2026-09-20T05:05Z] agent=agent-ticket11-adapters did=claimed next=harness-parameterization+anthropic-adapter
- [2026-09-20T07:10Z] agent=agent-ticket11-adapters did=AC#1+AC#3 GREEN: one shared harness (internal/llm/adaptertest/harness.go CaseTable+Run, 13 scenarios x 3 adapters) - openai-chat's golden assertions MOVED into it (openaichat/harness_golden_test.go), anthropic + openai-responses adapters implemented and run through the same table; cross-adapter proof TestC6EventsIdenticalAcrossAdapters (canonical trace) + TestFaultClassIsAdapterIndependent (one injected body, three protocols, one D37 class). 26 dialect fixtures added under internal/llm/testdata/golden (chat fixtures untouched); mockllm /v1/messages + /v1/responses upgraded to real dialect fidelity (flat tools, thinking/reasoning items, terminal status) + /__control/last_request body capture (recorded before the first response byte, including the fault branch). REJECTED one review instruction: 5xx -> ClassInternal would make adapters disagree on one shared classification, break retryability (observe maps only network/provider to retry) and contradict 14.2/D40#4 "Error(provider)" + ticket 09's provider-500 golden; kept the seam's ClassProvider and PINNED it with a cross-adapter test instead. next=AC#2 cache-breakpoint request-body test (mockllm last_request) then fallback chain + 14.2 matrix
- [2026-09-20T07:55Z] agent=agent-ticket11-adapters did=AC#2 GREEN: anthropic/cache_test.go reads the request body the LIVE mockllm captured (/__control/last_request) - exactly one cache_control on the last system block, none in the conversation suffix or tools, prefix bytes identical across two turns, x-api-key+anthropic-version headers, cache_read_input_tokens -> Usage.CachedTokens=8; plus breakpoint-limit (4) strict/loose and out-of-range cases. openairesponses/protocol_test.go proves the mirror image on the implicit-cache side (no cache_control ever, flat tool defs, instructions flatten, thinking_intensity -> reasoning.effort, stop_sequences error-not-drop). next=AC#4 fallback chain over two mockllm instances (primary must be PROVEN hit), then AC#5 14.2 matrix, AC#6 probe->provider_health, AC#7 token bucket
- [2026-09-20T09:05Z] agent=agent-ticket11-adapters did=AC#4+AC#5 GREEN. AC#4 internal/llm/fallback_test.go: two LIVE mockllm instances (primary=anthropic dialect, fallback=chat) behind a Resolver-built text_chain; every case asserts the PRIMARY server's own request counter (3 = initial+2 retries via ChainBuildOptions.RetryMax), so a chain that skipped the primary cannot pass; double failure -> Error(provider) + ctx.Err()==nil + the same *Request re-streams successfully after the fault clears (D40#4 resumable). FOUND+FIXED a real contract defect in chain.go: an element that is failed OVER flushed its Error/Stop{error}/Done to the consumer, so a successful fallback produced two terminal triples (violating C6 "exactly one terminal pair, then Done"); terminals are now buffered per element and discarded on failover, exactly as retry.go does per attempt. AC#5 matrix_14_2_test.go one test per 14.2 row (row1 retry-stays-Thinking + >5s notice on an injected MONOTONIC clock, row2 exhausted=explicit classified Error never silence, row3 401/403 auth->Unconfigured+never retried, row4 budget vs network distinct class/state/copy-key + only the network one retryable, row5 per-adapter partial+incomplete+open tool call+stream_disconnected marker state). Added llm.RetryNotice (structured, Label()=14.2 copy, BallStateName()=Thinking) and ChainBuildOptions.RetryMax; updated the stale "adapters land with ticket 11" message in provider.go and ticket 09's catalog test that relied on it (now a genuinely unregistered protocol). next=AC#6 probe->provider_health access layer, AC#7 token bucket
