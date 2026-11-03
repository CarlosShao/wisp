# 31 — reminder.* tools: one-shot in-process reminders (D12 narrow exception)

**Status:** done
**Claimed by:** 31-r2(内存层) / 31-r3(tools+agent) / 编排者收尾
**Last update:** 2026-11-03
**Blocked by:** 04-sqlite-core, 30-result-routing-d10
**Parallel slots:** 1
**Spec refs:** SPEC-07 §3 (reminder rows), 16.5.4, D12 partial reversal, S4

## What to build
`reminder.set/list/cancel`: one-shot, in-process, notify-only reminders — the deliberately
narrow exception to the no-scheduler rule (D12), with its limits stated to the user at set-time.

## Key constraints
- ONLY: in-process timer + one SQLite row + a `notify` fire. Set=list|cancel matrix:
  set L1, list L0, cancel L1.
- NOT (hard exclusions, each with a DEFERRED-style explicit refusal in code comments + registry
  note): recurring/cron reminders; reminder-triggered tool execution ("提醒我并顺手整理文件" →
  do only the first half); OS-level scheduling (registry/Task Scheduler — L2 + uninstall
  residue); missed-reminder replay across restarts (only "你错过了 N 条提醒" count on boot).
- Set-time user-facing limit disclosure: "进程未运行则不会触发" — spoken/notify confirmation
  includes this caveat (16.5.4 honesty rule); unit-test the message presence.
- Timer lifecycle inside SessionScope-independent task timer goroutine (named, owned); reminders
  survive Warm→Sleeping (process alive) but NOT process exit.
- Firing: notify via tool-internal path (not gated — it IS the notify capability the plugin
  system already grants; fire records a tool_call row for auditability).
- Time: fire timing via monotonic clock vs stored wall-clock target (recompute delta on wake/
  clock-jump — D42#9).

## Out of scope
- Any scheduler generalization; calendar integration (REJECTED core, D46 plugin path instead).

## Acceptance criteria
- [x] set → fires at T+delta with notify (fixture clock); cancel removes; list shows pending.
- [x] Persistence: reminder row survives Warm dispose; process restart → missed count message,
      no per-item replay.
- [x] Exclusion tests: recurring spec → rejected by schema; action-on-fire → rejected by
      schema/prompt contract ("只做前半句" documented behavior).
- [x] Clock-jump: system time +5h while pending → still fires correctly relative to monotonic.
- [x] Caveat message asserted in set-confirmation output.

## Progress log (append-only, newest last)

- 2026-11-02/03 31-r2 内存层：schema v6 reminders 表 + DAO + 迁移幂等（`5ebdf6961`）；memory 整包 `ok 30.35s rc=0`。
- 2026-11-03 31-r3 tools+agent：reminder.set/list/cancel + 单调时钟计时器（Notifier 通知 + `reminder.fire` tool_call 审计）+ boot 错过计数（`da9b5fb9c`）；tools(-skip Search) ok 32.4s、agent ok 6.56s、`go build ./cmd/wisp` rc=0。
- 2026-11-03 编排者收尾：修 `itoa(int64)` 测试签名 + 清 Go 源码字符串 emoji + 修 install.go phantom 引用 + 建证据 `docs/evidence/s1/{31r2,31}-impl.md`；`d22scan.sh` PASS=35 FAIL=0 clean（`37c468f41`）。
- 2026-11-03 结案：AC1-5 测试全覆盖（tools 15 条 + agent 13 条 `Test31*`），D22 门禁绿，证据齐。
