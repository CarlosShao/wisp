# 43 — Power & lifecycle events: suspend/resume, session end, autostart, crash recovery

**Status:** **DONE-带残余**（2026-10-31 编排者收口：五格可翻、AC 勾保持未勾＝15/28/41 先例；§16.11#10 判语成立；残余四条见末节）
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 28-session-scope-warm
**Parallel slots:** ≤2 sub-agents (A: power + session-end paths; B: autostart + crash recovery +
update hooks)
**Spec refs:** §6, D42#1#7, D38(e) fast path, D43 pending-approval open item note, S6

## What to build
Windows lifecycle robustness: suspend/resume handling (dispose session, reopen audio), fast
shutdown on WM_QUERYENDSESSION, autostart opt-in, crash-recovery semantics (interrupted tasks
visible, never auto-rerun), and the update-pending hook consumed by 56.

## Key constraints
- Suspend (PBT_APMSUSPEND): dispose session scope (models released) + stop capture BEFORE
  sleep; resume: reopen audio device, re-enter Sleeping/Armed per config, resync monotonic
  timers; pending L2 approvals at suspend → follow the §16.11#10 lean: void + auto-reject +
  "错过 N 个待确认操作" notice (decision logged; final wording at 54 review).
- Session end (WM_QUERYENDSESSION/WM_ENDSESSION): D38 fast path ≤5s — skip non-critical flushes,
  NEVER skip steps 4/5/9 (audio release, native sessions, Job close).
- Autostart: opt-in registry Run key; default OFF; toggle in tray menu + config.
- Crash recovery: start-after-crash marker → previous tasks marked 中断 in panel/CLI listing;
  never auto-rerun (side effects); diagnostics hint shown.
- Clock resync after resume: monotonic timers unaffected, wall-derived schedules recomputed
  (reminders).
- Update hook: on boot, if `update-pending.json` present → apply-handoff event (56 completes);
  N-1 rollback counters start here.

## Out of scope
- Updater binary itself (56); battery-saver policy internals (28 handles shrink; this ticket
  ensures events reach it).

## Acceptance criteria
- [ ] Real suspend/resume (psshutdown/manual on runner) → session disposed pre-sleep (no models
      resident during sleep — sampler), audio restored post-wake, one resume-e2e task works.
- [ ] Session-end: broadcast WM_QUERYENDSESSION → process exit ≤5s with steps 4/5/9 done (order
      audit log), Job children gone.
- [ ] Suspend with pending L2 → voided + counted notice (D43 note updated with final decision).
- [ ] Crash-recovery fixture: kill -9 mid-task → next boot lists task as 中断; no rerun.
- [ ] Autostart toggle writes/removes Run key; OFF default verified on fresh config.

## Progress log (append-only, newest last)

## ★★★终裁落地（2026-10-31 14:4x 编排者，凭非实现者 `43-v1`，件 `docs/evidence/s1/43-v1-verdict.md` 47 行、`bfffc2965`）

★**五格全可翻（AC 勾保持未勾）⇒ 可带残余结案**（15/28/41 先例）。
★**三个关键裁定**：①**替代形合法＝"是替代形不是没做"**（真睡眠不可自证＝无人值守结构性事实；双门真——powrprof dumpbin #85 实勘＋窗真建；headless 读数真组件、fake 仅在 C8/C5 认可接缝）②**定向注入合法、诚实声明成立**（HWND_BROADCAST＝外向动作打全桌面 ⇒ 定向到子进程自己窗口 class+pid 两事实身份零外溢；非伪授权形状）③**§16.11#10 判语成立、归口正确**（PLAN 原文倾向「作废并判拒绝」、实现＝lean 可执行化、只骑已冻结 #24 行未发明边；**PLAN note 文字＝契约面正确摆机主**；建议措辞落三处、归票 54）。
★**复跑读数**：winlive exit **94.86ms**（≤5s 余量 52 倍）、voided=2、kill -9 中断零重跑、autostart 真往返；**mutation 三发复刻全红**（EnginesHot/收据谎报被注册表抓/voided 计数）；**零回归 19/19＋session 42/42＋255 名册 20/20**。
★**残余四条**：①**托盘「开机自启」菜单项＝真残余、归下一张能写 internal/ball 的票** ②wall schedules＝非缺口 ③update-pending 消费端归票 56 ④**D43 note 文字＝契约面摆机主（唯一需机主动作）**。
★**瑕疵自陈（诚实）**：session 旧数 27→现 42；红句行号差一行（断言逐字同）；junction 补模型资产后 19/19；★**自己踩一次 0xc0000135（漏补在册 PATH、罚单自领）**。
