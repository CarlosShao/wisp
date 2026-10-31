// Package present is the D10 result-presentation router (SPEC-05 §7): it
// takes one finished task's final text and decides where it goes - file-drop
// plus panel, system notify plus panel, or TTS-only with a clipboard copy -
// under the four frozen tiers of D10 (first match wins, structured markers
// beat length).
//
// The package owns ONLY the routing policy and the summary protocol. It owns
// no channel: every presentation surface arrives as a narrow port on
// Channels, so the real TTS engine, the C17 panel bridge, the ball badge and
// the system toast stay exactly where they are today (cmd/wisp assembles
// them; ticket 36 owns the panel rendering surface, ticket 26 the TTS
// internals). nil port means "this channel is not wired on this host": the
// router records the skip in the receipt and never panics on it.
//
// Two properties are structural, not promised:
//   - the clipboard port is reachable only from the tier-4 branch (D10: long
//     results never overwrite the user's clipboard);
//   - tier 3 always drives badge and task list alongside the notify port
//     (D42#8: Focus Assist may swallow the toast, so the dual-channel update
//     is unconditional and never probed via any API).
//
// The spoken summary comes from the SAME LLM response the result came from:
// the loop's (6) style section instructs the model to append a <=60-character
// spoken summary, and ExtractSummary lifts the trailing 口播摘要 line out of
// the final text. There is no second LLM request anywhere in this package,
// and no goroutine: Dispatch runs on the caller's goroutine, and slow
// channels are the host adapters' queues (D38b), not this package's worry.
package present
