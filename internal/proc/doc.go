// Package proc owns OS-process-level machinery (SPEC-01 §3): the JobScope
// (C30), the per-session single-instance mutex, the D38(e) shutdown order
// orchestration, and (deferred) update/rollback (D41b).
//
// Responsibilities (implemented by ticket 03):
//   - JobScope: every child process is assigned to one Job Object with
//     KILL_ON_JOB_CLOSE; TreePrivateBytes() is THE D32 SLO metric
//   - single instance: named mutex Local\wisp-single-instance (per-session,
//     D42#7); a second launch signals the first (activate event) and exits
//   - shutdown sequence: the 10-step D38(e) order is executed and audited by
//     test; the fast path (WM_QUERYENDSESSION) skips only the non-critical
//     flushes of step 7, never steps 4/5/9
//   - WISP_ENV typed fork: env resolution + data dir / mutex name pure
//     functions (prod/dev defaults per SPEC-03 §5.2; test-env values and
//     endpoint defaults are ticket 06)
//
// Non-responsibilities:
//   - no business process management (scheduler), no service/watchdog logic
//     (watchdog), no config values for envs (config, ticket 06)
//
// DEFERRED(update/rollback): D41(b) staging/atomic replace/rollback lands
// with the release tickets (56). Package name is proc; the build-tagged
// Windows implementations live in *_windows.go files.
package proc
