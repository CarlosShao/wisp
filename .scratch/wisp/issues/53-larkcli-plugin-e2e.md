# 53 — lark-cli wrapper plugin e2e: the S7 "real plugin" acceptance artifact

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 48-approval-queue-full, 51-tier2-goja, 52-d46-command-plugins
**Parallel slots:** 1
**Spec refs:** 16.5.6, D44 S7 real-plugin criterion, §15#10, D46

## What to build
The actual `lark-cli` Tier1 command plugin (calendar agenda / send message / tasks) installed
and driven by voice/text end-to-end — the S7 "a real plugin runs" acceptance (no toy plugin).
Repo may carry the example manifest; NEVER any credentials/tokens.

## Key constraints
- Plugin manifest: agenda query (read-class → expect L0 pass-through), send-message
  (send-class → R8 forces L2 even if manifest declares L0), task read/write mix; each via
  `command` kind on the user's locally-installed authenticated lark-cli (README documents
  "install + authenticate lark-cli yourself").
- E2E assertions: agenda query returns parsed items into the conversation (extract works
  against real CLI output shape); send-message triggers L2 card with R4/R8 reasons, reject
  path and (user-approved) allow path both exercised; timeout/env/hash rules re-verified on
  the real binary.
- One Tier2 sample (optional, behind tier2 gate): a tiny JS plugin using wisp.state only —
  demonstrates C24 without broad surface.
- Failure drills against real CLI: binary replaced (hash fail), PATH stripping effect, hang
  (timeout), oversized output (truncate).
- No credentials/tokens in any repo file (secret scan must be clean; README note present).

## Out of scope
- Other vendors' CLI plugins (user-authored later); registry (57).

## Acceptance criteria
- [ ] Real-plugin E2E: voice/text → agenda listed; send → L2 → user-approved send observed
      in lark-cli backend (or mocked transport with identical shapes if account unavailable —
      recorded decision).
- [ ] All four failure drills green on the real binary.
- [ ] Manifest + README committed; secret scan clean; credentials never touched by tests.
- [ ] Evidence pack (logs + tool_call rows + screenshots) committed under docs/evidence/s7/.

## Progress log (append-only, newest last)
