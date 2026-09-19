# 09 — LLM provider seam: C5/C6 core, OpenAI Chat adapter, golden replay, mockllm server

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 03-skeleton-runtime-rules
**Parallel slots:** ≤2 sub-agents (A: llm module + OpenAI Chat adapter; B: tools/mockllm server
+ golden harness + compose wiring)
**Spec refs:** SPEC-05 §3, D8, C5, C6, §14.2, D42#4, SPEC-11 §4

## What to build
The `llm` module seam: normalized `LlmProvider` interface, `StreamEvent` union (C6), internal
message representation (C7 incl. Image part), usage accounting, retry/backoff, and the first
adapter (OpenAI Chat Completions with tool_calls). Plus `tools/mockllm` (OpenAI/Anthropic-
compatible mock with golden replay + fault injection) and the golden-SSE test harness that all
later agent tests reuse.

## Key constraints
- C6 event set exactly: TextDelta, ReasoningDelta, ToolCallStart{id,name}, ToolCallArgsDelta,
  ToolCallEnd{id}, Usage{in,out,cached}, Stop{reason ∈ end_turn|max_tokens|tool_use|
  stop_sequence|content_filter|cancelled|error}, Error{class(D37), provider_code, retryable,
  retry_after}, Done.
- C7 content parts: Text | Image{mime, bytes_ref, alt} | ToolUse | ToolResult{id, content[],
  is_error}; bytes_ref never persisted/logged (privacy).
- Adapter maps normalized↔wire; agent core never sees protocol. Errors classified into D37
  classes at the seam.
- Retry: exponential backoff ≤3 for network/5xx/429 (429 honors retry-after); 401/403 → no
  retry, Unconfigured; stream mid-disconnect → keep partial, mark incomplete, unclosed tool
  calls → failed (consumed by 10). Proxy: system proxy (WinHTTP default) + HTTP(S)_PROXY env;
  TLS-intercept errors must distinguish "untrusted cert" from "unreachable" (D42#4).
- Timeout: `[llm] timeout_ms` via context; cancellation propagates (ctx) — no wall-clock math.
- mockllm: `POST /v1/chat/completions`, `/v1/responses`, `/v1/messages`; stream:true SSE with
  tool_calls + usage; golden mode replays `testdata/golden/<name>.sse`; control endpoints
  `/__control/fail_next|latency|truncate|reset`. Runs via `go run` AND compose (port 18080).
  Golden file format shared with unit-test replayer (one format, two runners).

- **Provider presets & catalog primitives (2026-09-19 supplement)**: preset table
  `openai/anthropic/deepseek/qwen/zhipu/moonshot/siliconflow/openrouter/ollama/minimax/mimo/
  stepfun` (protocol + base_url defaults; model entries user-configured or auto-discovered);
  **auto-discovery** via OpenAI-compatible `GET /v1/models` → import with capabilities=unknown;
  **probe primitive interface** on C5 (per-capability test request definitions: fc/vision/
  thinking/audio — implementation in 11, results → `provider_health`, SPEC-02 v2).

## Out of scope
- Anthropic/Responses adapters (11); agent loop (10); probe IMPLEMENTATION (11); quota
  enforcement (44).

## Acceptance criteria
- [ ] Golden replay tests: tool-call assembly, stopReason=max_tokens, mid-stream disconnect,
      cancellation, backoff ladder — all via recorded SSE bytes (no function mocks).
- [ ] mockllm golden mode + unit-test replayer produce byte-identical event streams.
- [ ] Fault-injection tests via control endpoints: 429 (retry-after honored), 5xx ×3 → provider
      error, 401 → no retry, latency injection respected.
- [ ] Proxy tests: HTTP(S)_PROXY honored; untrusted-cert error class distinct from unreachable.
- [ ] Usage events aggregate correctly across a multi-chunk stream.
- [ ] Preset table covers all 12 providers with correct protocol/base_url defaults; auto-discovery
      against mockllm `/v1/models` imports with unknown capabilities; probe primitive contract
      test (each capability has a defined minimal request/response check).

## Progress log (append-only, newest last)
