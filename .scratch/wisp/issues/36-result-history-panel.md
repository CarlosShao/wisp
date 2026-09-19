# 36 — Result & history panel: chat stream, agent work states, sanitizer, transcript view

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 30-result-routing-d10, 35-panel-bridge-c17
**Parallel slots:** ≤2 sub-agents (A: chat stream + 14 agent-state visuals; B: history/transcript
+ purge/export + virtualization)
**Spec refs:** SPEC-08 §5.4, §17.5, F2 sanitizer layer, D10 panel surface, §14.4 transcript rules

## What to build
The main result surface: streaming chat view with the 14 agent-work-state visuals (thinking
band, SSE cursor, reasoning collapse, tool chips with risk badges, approval-wait, cancel/error/
stuck cards), plus the history & transcript browser over SQLite with privacy controls.

## Key constraints
- 14 states per SPEC-08 §5.4 table — all rendered (design/screens/chat.html is the visual
  reference; human sign-off required): thinking light-band + elapsed timer; streaming cursor +
  token counter; reasoning collapsed-by-default; tool chip = category icon + name (mono) +
  truncated args + status arc/check/x/ban + L0/L1/L2 outline pill; expanded chip = colored JSON
  + duration + correlationId; approval-wait amber pulse-once + depth; error card human-copy +
  expandable error_class; stuck card names the repetition; cost line mono tabular-nums.
- Markdown: react-markdown + remark-gfm; **rehype-sanitize whitelist mode** (tightened schema);
  `web.fetch`/`doc.read` outputs NEVER take the HTML path (plain text or sanitized markdown
  only); no `dangerouslySetInnerHTML` anywhere (CI grep); code blocks → Shiki static highlight.
- Approval interplay: L2 card inside panel shows full params + rules_hit + reject button +
  "点击悬浮球以批准" persistent guide with pulsing ball glyph — NO allow button (F2, §15#6).
- History: task list (30d) → task detail (transcript + tool calls) from SQLite via bridge;
  virtualized list (react-virtual) for 1000+ rows; per-item delete + purge-all + export JSON
  (privacy APIs from 29/04); transcript display honors §14.4 masking marker.
- Copy buttons (code/paths) go through clipboard tool path for taint bookkeeping.
- All loops respect the CPU rule: no animations when panel hidden (visibilitychange pause).

## Out of scope
- Approval queue multi-task view (38); palette (38); config editor (39).

## Acceptance criteria
- [ ] 14-state visual checklist against chat.html — user sign-off (screenshots committed).
- [ ] XSS suite: `<img onerror>`, `<script>`, `javascript:` href, event-handler attrs in fake
      fetch/markdown content → neutralized; zero script execution (CSP report-only listener
      asserts zero violations).
- [ ] dangerouslySetInnerHTML grep = 0; sanitizer whitelist config reviewed + committed.
- [ ] History: 1000-row fixture scrolls at 60fps; delete-one/purge/export verified against DB.
- [ ] Long-task stream: 200-chunk golden replay renders with merge (no jank), token counter
      accurate.
- [ ] Hidden-panel animation pause asserted (rAF counters).

## Progress log (append-only, newest last)
