# 38 — Command palette (cmdk) + tasks view: text input channel, force rules, discovery

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 35-panel-bridge-c17
**Parallel slots:** ≤2 sub-agents (A: palette + text entry; B: tasks list + tools directory)
**Spec refs:** SPEC-08 §5.2/§7, D9 (four roles), D11③ force rules, D34 task.list/cancel, S5

## What to build
The command palette as D9's four-role surface: searchable text entry fallback, L2 confirm
reliable channel, plugin/tool discovery with risk visibility, and the concurrent-tasks list
view (task.list/cancel tools backing it).

## Key constraints
- Palette = cmdk: fuzzy commands (open panel pages, mute, settings, tools) + free-text prompt
  entry (keyboard-first, ⏎ submits to agent as text task — the D9 noise-environment fallback);
  reachable by hotkey + tray + ball click; Chinese IME composition handled correctly.
- Tool discovery: lists resident + third-party tools with name/description/risk/capability
  ("我到底装了什么" transparency role); invoking a tool from palette goes through the SAME host
  bridge (no bypass).
- Tasks view: task cards (id, state, started, cost-to-date, current step) + cancel button →
  `task.cancel` (L1); blocked tasks show "waiting on what" (D31); queued-behind-pathlock tasks
  render their blocker (D43 #20/#39 groundwork — real multi-task at 47).
- force_tool/force_chat TOML rules (D11③): palette offers per-conversation override chips +
  config file support; overrides surface in agent context assembly.
- All palette state via bridge resync only (stateless); palette → result handoff switches to
  the result view (36) with streaming continuation.
- Cancel semantics: cancelling from palette → same ctx cancellation + applied-steps report
  (20) if mid-execution.

## Out of scope
- Multi-task concurrency engine (47); voice path; config editing (39).

## Acceptance criteria
- [ ] Text prompt via palette → streams into result view; identical behavior to CLI path.
- [ ] IME: Chinese composition (pinyin → 选字 → commit) submits final text only (e2e with
      Windows IME on runner or manual evidence).
- [ ] Discovery: installed tool set listed with correct risk pills; launch routes through gate.
- [ ] Tasks view: live task state transitions reflected; cancel mid-tool → applied-steps report.
- [ ] force rules: per-task override chip + TOML rule both steer agent mode (prompt snapshot).
- [ ] Keyboard-only: full palette flow reachable without mouse (tab/arrow/enter audit).

## Progress log (append-only, newest last)
