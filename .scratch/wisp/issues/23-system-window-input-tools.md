# 23 — System/window/input tool family: system.*, window.*, clipboard, media, screen, input.type

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 20-host-bridge-fs-tools
**Parallel slots:** ≤2 sub-agents (A: system.*/media/clipboard; B: window.*/screen.capture/
input.type + C7 Image wiring)
**Spec refs:** SPEC-07 §3 (D34 rows), C7 Image, D45 input binding, §16.11#2 UIPI, D34 note②

## What to build
The machine-control builtin family per the D34 authoritative table, including the C7 Image
content part that makes `screen.capture` actually reach the LLM, and `input.type` with its
per-target-process session-binding groundwork.

## Key constraints (RiskLevels fixed by D34 — verify each)
- `system.get` L0 (focused-window title taint-marked), `system.set` L1 (volume/brightness/mute/
  media keys/IME), `system.power` L2 (lock/sleep/restart/shutdown/logoff — shutdown+restart
  require typed confirmation word).
- `window.list` L0; `window.focus` L1; `window.manage` L1 (min/max/move/resize/snap) but
  close = L2 (unsaved-content risk).
- `clipboard.read` L0 (taint) / `clipboard.write` L1 (exfil channel per F4).
- `media.control` L1 (play/pause/next/prev).
- `screen.capture` L1 → C7 ImagePart{mime, bytes_ref, alt}; bytes_ref memory-only — never on
  disk/DB/logs/diagnostics unless user explicitly saves; DRM-window black-frame documented.
- `input.type` L2 first-use per target process → session-grant downgrade hook (binding table
  now; grant lifecycle at 49): confirm card carries "仅本次 / 本会话内允许向此应用 / 所有应用
  (danger-styled, expanded click)" options; UIPI: injecting into elevated windows fails with
  explicit "Wisp 无法向管理员窗口输入" (never attempt bypass); types text AND key sequences at
  the focused window.
- All system calls through the host bridge; risk decisions via 17 (R5-R8 rows exercised).
- Fixtures: a test target app (notepad/scripted window) for window/input acceptance on the
  Windows runner.

## Out of scope
- Grant persistence/lifecycle (49); voice dictation flow (32); panel screens (37/40).

## Acceptance criteria
- [ ] RiskLevel matrix test: every action above asserted (esp. window.close=L2, power=L2,
      set=L1, get=L0+window-title taint).
- [ ] screen.capture: screenshot reaches mock LLM as ImagePart (C7 wire format asserted);
      bytes_ref absent from disk/db/logs; DRM window → black frame + doc note.
- [ ] input.type: types 300-char text into fixture app focused window; elevated-process
      injection → explicit UIPI error; confirm-card three options render with danger styling
      on "all apps".
- [ ] clipboard round-trip + taint propagation into R4 (write of tainted content → L2).
- [ ] system.power shutdown path requires typed confirmation word (unit test with fake executor).

## Progress log (append-only, newest last)

- 2026-11-03 21:4x 编排者：后端+工具接线层落地（119c3fa91 后端6文件；a2f6f4a95 system.go BuiltinSystemEntries 11工具+run.go注册+system_23测试 1450行；df67e002e 修229普查登记52的OS进程spawn站点）。internal/tools 整包 ok、d22 clean。**已满足**: AC1 风险矩阵、AC5 关机确认词、AC4 taint→R4。**残余**: AC2 C7 ImagePart→LLM 接缝(需 internal/agent loop+llm.ImagePart 集成)+bytes_ref内存only/DRM黑帧；AC3 input.type 真fixture app+UIPI 真机+确认卡三选项UI(前端)；AC4 真剪贴板round-trip。**卡点**: AC3 确认卡UI=frontend(机主外部agent)+真机集成，非纯后端可闭合。票仍 open。
