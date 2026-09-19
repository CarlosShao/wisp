// Package llm owns the LLM provider seam (SPEC-01 sec 3, D8, C5/C6/C7): the
// agent core runs only on the normalized representation defined here and
// never sees protocol shapes.
//
// Contents (ticket 09):
//   - C5 LlmProvider interface + Request/ProviderInfo, incl. the
//     cache-breakpoint capability bits (SPEC-05 sec 4.1)
//   - C6 StreamEvent union (exact event set) + Usage aggregation +
//     TurnCollector (partial-turn marking, open-tool-call detection)
//   - C7 content parts (Text/Image{bytes_ref}/ToolUse/ToolResult) with the
//     bytes_ref privacy rule (never persisted/logged)
//   - D37 error classification at the seam, incl. D42#4 network
//     discrimination: tls_untrusted_cert vs connect_failed vs timeout
//   - retry ladder (exponential backoff <= 3; 429 honors retry-after; 401/403
//     never retried), text_chain runner, local rpm/tpm token buckets
//   - provider catalog consumption (roles, chains, compat switches), the
//     preset/preset-inheritance consumer, /v1/models auto-discovery with
//     capabilities-unknown import, and the probe primitive contract
//   - golden: the golden-SSE format, unit-test replayer and recorder writer
//     (one format, two runners; tools/mockllm is the second runner)
//
// Non-responsibilities: agent policy (internal/agent), secrets storage
// (internal/secret), UI streaming merge (internal/panel), quota enforcement
// (C23, ticket 44), the remaining protocol adapters (openai-responses and
// anthropic -> ticket 11).
package llm
