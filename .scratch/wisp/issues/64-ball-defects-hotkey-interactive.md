# 64 — 票 07 遗留缺陷修复：热键接线 / 交互四项 / Sleeping 零定时器实测 / 多显示器实拖

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-20
**Blocked by:** 07-ball-state-machine-core（已 done，但下列四条 AC 从未通过）
**Parallel slots:** ≤2 sub-agents（A：热键接线与错误语义；B：交互/多显示器/零定时器断言）
**Spec refs:** SPEC-08 §2, D43 状态机, SPEC-05 §6（热重载）, 票 39 前置, registry A1–A7
**来源:** 2026-09-20 AC 补裁（`docs/reports/pending-and-issues.md` §"AC 补裁遗留"）+ 实况签收发现

## What to build
票 07 被标 DONE，但补裁查明它的 6 条 AC 里有 4 条**从未被验证**、其中一条标的东西是死代码。
owner 2026-09-20 裁定：**现在插队修，不等 S3**。本票消化 registry 的 A1–A7。

1. **A1 热键接线（最高优先）**：`Ball.RebindHotkeys`（`internal/ball/ball_windows.go:693`）
   **全仓零调用者** → 改 `[hotkey]` 配置后热键根本不会重注册。要接
   `config.Manager` 的热重载回调 → `RebindHotkeys`，并把 `registerAll` 的失败语义拆开：
   **"被他人占用"（ERROR_HOTKEY_ALREADY_REGISTERED=1409）** 与 **"未尝试/其它错误"** 必须给出
   不同提示（现状是一律 `slog.Warn` 后丢弃，用户完全不知道键没生效）。
2. **默认唤起键换掉**：`Ctrl+Alt+W` 在 owner 本机被第三方程序占用（实测：全 84 个候选组合里
   **只有 W 这一个被占**，其余 83 个空闲）。裁定改为 **`Ctrl+Alt+Q`**；`Ctrl+Alt+Space`
   本机也空闲，作为可配置备选写进文档（注意部分中文输入法会抢 Space 组合，故不作默认）。
   改默认值属 SPEC-08 之外的默认参数，但**仍要在票里留痕**，别让它悄悄变。
3. **A2 交互四项**：`HitTest()` 只覆盖了 4 个条款里的 1 个。补齐：单击球 Sleeping→Listening、
   从 Confirming 用 Esc/单击取消、**不抢焦点**（真起一个有焦点的记事本类窗口验证）、
   透明区点击穿透。
4. **A3 Sleeping 零定时器实测**：现有测试断的是**策略表**，不是活句柄。改成实测
   `EnumThreadWindows`/timer 句柄集合为空（`winlive` tag 下），且**不能被 `t.Skip` 掉**。
5. **A4 多显示器实拖**：detach→primary 与持久化/恢复有测，但"拖到第二块屏"从未真执行。
   本机只有单屏（`\\.\DISPLAY4` 3440x1440）——**若无法真机验证就如实留着不勾**，
   并写明需要什么硬件，不许用双屏 mock 冒充。
6. **A5–A7**：见 registry 对应条目，逐条消化或明确登记为"待硬件/待票 39"。

## Key constraints
- **测量必须独占**：本票含 winlive 与句柄/时序断言。跑之前**必须确认没有其它 balldebug/wisp 进程在跑**
  （实况签收期间出过两次污染：`TestBallLiveLifecycle` 按窗口类名 `FindWindowW` 会匹配到演示球；
  "四个热键全失败"也是这么来的）。测试开头加一道前置断言：检测到同名类窗口已存在就 **fail 并提示**，
  不要静默 skip——静默 skip 正是这轮漏检的根因。
- 不改 20 态语义与 D43 迁移表（视觉重做在票 62，本票只管功能与接线）。
- 禁止裸 `go func(`（热键消息回路要有 owner+recover）；禁止墙钟差判超时；禁止 emoji。
- 与票 62 的边界：本票**不碰材质/动画/吸附**，只让"球能被唤起、配置改了能生效"成立。

## Acceptance criteria
- [ ] 改 `[hotkey]` 后热键真的重注册（端到端：改配置→按新键→球响应；旧键不再响应）。
- [ ] 注册失败两类语义分开：被占用 vs 未尝试，各自给出用户可见提示（不再只有一行 Warn）。
- [ ] 默认唤起键为 `Ctrl+Alt+Q`，且 `Ctrl+Alt+Space` 作为可配置备选写进文档。
- [ ] 交互四项各有真机测试且**不被 skip**：单击唤起、Confirming 取消、不抢焦点、透明区穿透。
- [ ] Sleeping 零活动定时器句柄为**实测断言**（winlive），策略表断言保留但不作为唯一证据。
- [ ] 多显示器：真拖到第二屏验证并留证；本机无第二屏则**保持未勾**并写明所需硬件。
- [ ] 测试前置：检测到外部同窗口/演示进程时 fail-fast 并提示，禁止静默 skip。
- [ ] registry A1–A7 逐条标注"已修/待硬件/移交票 xx"，不许无声消失。
- [ ] 对抗验收由非实现者执行，报告含与本表 **1:1 的裁决表**（README 规则 6）。

## Progress log (append-only, newest last)
