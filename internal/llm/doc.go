// Package llm owns the LLM provider adapters (SPEC-01 §3): the three
// protocols, the normalized StreamEvent (C6), retry/backoff, and usage
// accounting (C23 inputs).
//
// Responsibilities:
//   - protocol adapters (OpenAI-compatible, Anthropic, ...; tickets 09/11)
//   - StreamEvent normalization incl. Error{class, provider_code, retryable,
//     retry_after} where class is an observe error class (D37)
//   - network/proxy/TLS error discrimination (D42#4: "cert untrusted" must be
//     distinguishable from "cannot connect")
//
// Non-responsibilities:
//   - no agent policy (agent), no secrets storage (secret), no UI streaming
//     merge policy (panel merges, per D38d)
//
// DEFERRED(adapters): implemented by ticket 09 (OpenAI-compatible + mockllm)
// and ticket 11 (remaining REST adapters). This ticket only freezes the
// package boundary.
package llm
