// Package observe owns cross-cutting observability (SPEC-01 §3): slog JSONL
// structured logging with redaction, diagnostics bundle, SLO sampling, and
// CostMeter (C23) accounting.
//
// Ticket 03 additionally hosts two runtime-rule primitives here because the
// frozen module table has no other cross-cutting home and SPEC-01 forbids
// adding modules:
//
//   - the D37 error model (errors.go): the 17 error_class values, their
//     state mapping and retryability. These values are the single definition
//     for log fields and the future task_log.error_class column (memory,
//     ticket 04). User-visible copy text is deferred (D23) — only stable
//     message keys are defined here.
//   - the named-goroutine registry (goroutine.go): Spawn with name/owner/root,
//     recover boundary, roster verification (D38b) and panic bookkeeping.
//   - the monotonic-clock discipline helpers (clock.go, D42#9).
//
// Non-responsibilities:
//   - no business metrics semantics (watchdog decides thresholds), no log
//     retention on disk (memory/RetentionJob), no UI diagnostics pages
//
// DEFERRED(JSONL sink/redaction/diagnostics/CostMeter): implemented by
// ticket 08. This ticket implements only the three primitives above plus a
// minimal structured-error sink used by the registry.
package observe
