# 11 — LLM adapters: Anthropic Messages + OpenAI Responses + fallback chain

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
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

## Out of scope
- Ball/panel UI rendering of these states (ball shows via state events; panel later).

## Acceptance criteria
- [ ] Both adapters pass the same golden/fault test suite as OpenAI Chat (identical harness).
- [ ] Cache-breakpoint test: Anthropic request bodies show stable cached prefix across turns
      (capture via mockllm); cache-read usage populated.
- [ ] ReasoningDelta surfaces for thinking responses on both adapters.
- [ ] Fallback chain test: fail_next(3) on primary → fallback used → task completes; double
      failure → Error(provider), ctx preserved, resumable.
- [ ] Failure-semantics matrix: each §14.2 row has a test asserting the emitted user-visible
      event/state.

## Progress log (append-only, newest last)
