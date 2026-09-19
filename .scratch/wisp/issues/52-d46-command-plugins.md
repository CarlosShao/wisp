# 52 — D46 command plugins: pinned exe, argv templates, declarative extraction

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 50-tier1-manifest-plugins
**Parallel slots:** 1
**Spec refs:** SPEC-07 §4, 16.5.6 D46, §15#10, C14 command fields, S7

## What to build
The Tier1 `command` tool kind wrapping local CLIs (lark-cli, git, ffmpeg, …) with zero Go code:
absolute pinned exe + per-execution hash check, argv templates executed WITHOUT any shell,
strict parameter schemas, declarative output extraction, environment allowlist, timeouts —
the plan's primary capability-expansion path.

## Key constraints
- Install-time: `exe` must be absolute + stable path (TEMP/variable paths → refuse install);
  compute + pin `exe_hash` (sha256) in plugin_state; EVERY execution re-verifies hash →
  mismatch = refuse + alert (supply-chain/local-swap defense).
- Execution: build argv from `argv_template` with schema-validated parameters → direct
  CreateProcess argv vector — NEVER cmd.exe/sh, NEVER string concatenation. Parameter values
  containing `` | & ; $ ` < > \n \r `` → reject pre-spawn.
- Parameter schemas: per-placeholder JSON Schema (type/enum/range/max-length); violations never
  reach argv.
- Output: `extract` (json-path or regex) shapes what enters context (≤32KB); raw stdout never
  enters context; stdout/stderr caps 1MB each → kill + truncated; default timeout 15000ms
  (manifest-overridable) → kill process tree (Job).
- Environment: allowlist only (PATH-trimmed, HOME, USERPROFILE, LANG + manifest-declared) —
  never inherit parent env (key-leak defense); test CLI echoes env to prove.
- Risk: manifest risk = R1 lower bound; send-class actions (send message/mail, invite) → R8
  forces L2 regardless of declaration; read-class may pass L0.
- Distinction rule from shell.exec documented + tested: command tools require NO shell_enabled.

## Out of scope
- lark-cli business plugin content (53); registry distribution (57); arbitrary shell (still
  default-disabled D14).

## Acceptance criteria
- [ ] Pinning: install computes hash; one-byte exe tamper → refuse + alert (real file test).
- [ ] Injection: meta-char parameter → rejected pre-spawn; schema-violating param → rejected;
  argv vector asserted (no shell spawn — process-creation audit).
- [ ] Extraction: JSON + regex extractors on fixture CLIs; >32KB extract → truncated; 1MB
  stdout → kill + truncated (no OOM).
- [ ] Env: child env contains ONLY allowlist (echo-fixture proof); API-key env absent.
- [ ] Timeout: hanging CLI killed at timeout (process-tree gone).
- [ ] R8: declared-L0 "send" tool still judged L2 end-to-end (mock provider + risk trace).

## Progress log (append-only, newest last)
