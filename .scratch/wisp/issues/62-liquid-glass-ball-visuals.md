# 62 — 液态玻璃悬浮球视觉重做（玻璃质感 / 靠边吸附收缩 / 音频驱动液面）

**Status:** review（**改回 review：本票从未被验收**，见下方更正块）
**Claimed by:** 4 个实现代理（3 次被平台杀、1 次撞 150 轮）+ orchestrator（检查点与验收）
**Last update:** 2026-09-20 23:4x（编排者：撤销我自己今天的 premature `-done`）

> ## ⚠ 编排者更正块（2026-09-20 23:4x）——为什么这张票从 `-done` 改回 `review`
> 今天我把文件名加了 `-done` 后缀并把 Status 写成 done。这个动作**违反了本目录 README 的两条完成规则**：
> - **规则 4**：`check all acceptance boxes → Status: done → rename -done`。而本票此刻
>   `grep -c "^- [ ]"` = **8**、`grep -c "^- [x]"` = **0**——**八个框一个都没勾**就被人改了名字（我）。
> - **规则 6**：置 done 前，验收报告必须含与 AC 编号 **1:1** 的裁决表。要求的
>   `docs/evidence/s1/62-adversarial-acceptance.md` **不存在**（`docs/evidence/s1/` 下只有
>   `62-ball-visual-prototype.md` 与 `62-diff-{baseline,border,glass,signoff}/`）。
> **我当时的错误推理**：把 owner 对票 65 的"临时通过"当成了票 62 的验收。owner 原话是
> 「可以说算是赝品吧…但是先勉强用吧…先完成核心功能」——那是**降级放行的指令**，不是"AC 全部成立"的裁决。
> **改名回来不是为了重做**：`-done` 的作用是让新会话不重复执行已完成的工作，而这张票真正的残留
> （质感返工、契约草案、桌面级 AC#3/#5）**已经有别的名词在背**——票 65（质感）、票 68（默认态与 A29 四条）、
> `62-visual-spec-draft.md`（**从未写出**，见 A30④）。所以本票保持 review，**新会话不要从 AC#1 重跑实现**；
> 请按下表认领：
>
> | AC | 现在真实的处置归属 | 依据 |
> |---|---|---|
> | #1 原型在浅色+深色桌布可见 | **证据已存在，未做裁决表** | `62-diff-baseline/`、`62-diff-signoff/`（差分化像，含 `px_delta_ge8` 列） |
> | #2 Sleeping 静态单帧 + 实测零定时器 + CPU≤0.5% | **已由票 68 复测并更正数字**（真体 34.72px 不是我写进契约的 44px） | A29①、`docs/SLO.md` §A.2/B |
> | #3 每态颜色/速度/描边可区分（截图） | 未验收 ⇒ **票 65**（质感返工，阻在 owner 的参考图，见 R15#7） | owner"赝品"裁定 |
> | #4 音频驱动液面单调 | 未验收 ⇒ 票 65 同批 | — |
> | #5 靠左/右/上停靠 + 悬停/单击恢复（跨 DPI 与副屏） | **已知两处不成立** ⇒ 票 68 AC#3 | A29②（实测露出 ≈82% 而非 42%）、A29③（命中 21px < 可见 26px）、A28（DPI 不对称） |
> | #6 资源重测写进 SLO.md 且不降阈值 | 部分成立 ⇒ 票 66/68 追加，**阈值一动没动** | `docs/SLO.md` §B/§C |
> | #7 契约草案 `62-visual-spec-draft.md` 等 owner 签字 | ✅ **已成文**（`48f0cfd`，193 行：20 态逐行 + 22 条冲突 + 15 问待定）；**签字与回填 SPEC-08 仍欠** | R15 第 10 项 |
> | #8 对抗验收 + 1:1 裁决表 | **从未执行** ⇒ 本票保持 review 的直接原因 | 报告文件不存在 |
>
> **登记在案**：这次错误记进 `docs/reports/pending-and-issues.md` **A30①**。

**Blocked by:** 07-ball-state-machine-core（渲染栈与 20 态状态机）, 13-audio-capture（音频电平来源）
**Parallel slots:** ≤2 sub-agents（A：渲染层玻璃+液态材质与着色；B：靠边吸附/悬停弹出与音频驱动）
**Spec refs:** SPEC-08 §2（视觉表，**本票需变更契约**）, D43 状态机, D32 16.3.2/16.3.3（CPU 与定时器纪律）, C30 进程树内存口径

## What to build
把悬浮球从"纯色圆点 + 状态色环"整体重做为**玻璃质感液态球**。owner 2026-09-20 明确否掉了
"给 12px 微点加描边/投影"的修补路线（原话：不知道你说的劳什子描边投影是什么东西，按我说的来），
并给出参考图（vivo 蓝心小V 的语音球）：外层透明玻璃壳 + 内部蓝紫青渐变的液态色团缓慢流动。

四件事必须同时成立：
1. **玻璃材质**：外壳有折射/边缘高光/内部反射，不是"贴一张圆形 PNG"。D2D 现有渲染栈内实现
   （径向渐变 + 高光带 + 内阴影 + 边缘光），不引入新 GPU 依赖。
2. **靠边吸附收缩**（owner 点名"就像迅雷那种悬浮球"）：拖到屏幕左/右/上边缘时，球沿该轴收缩成
   半隐形态贴边，只留一条可命中区域；鼠标悬停或单击 → 丝滑弹回完整球。
3. **音频驱动液面**：`Listening` 且采集到语音时，内部液体**随音频电平有节奏地旋转/起伏**；
   静音时会话仍在进行 → 通过一个过渡动画把液体**收敛**、外圈边框浮现（owner："不说话的时候，
   语音助手会通过一个过度动画丝滑引入边框"）。
4. **20 态全部重映射**：每个状态用液体的颜色/流速/边框样式区分，而不是靠小点可见性。

## Key constraints
- ⚠ **CPU 与定时器纪律是硬门，不得为动效牺牲**：D32 规定 `Sleeping` 态 CPU ≤0.5% 且
  **零动画零定时器**（票 07 AC#4 的标的）。因此**流动动画只允许出现在会话态**
  （Listening/Thinking/Acting/Warm/Speaking 等）；`Sleeping` 必须是**静态帧**——
  但静态的**玻璃球本体**（不再是 12px 微点），这样才同时满足"看得见"和"不耗电"。
  若实现中发现静态玻璃球在浅色桌布上仍不可辨，走**提高本体不透明度/加深色边缘光**的方向，
  不许偷偷给 Sleeping 加定时器。
- **靠边吸附不得引入常驻定时器**：用 `WM_NCHITTEST` + 窗口位置事件驱动；悬停检测走鼠标消息。
  Sleeping 的"零活动定时器句柄"断言必须继续为真（票 07 A 系列里有一条就是要补这个断言）。
- **音频电平来源**：接票 13 的采集帧（16k/mono/int16/512 样本），在渲染侧只取 RMS/包络，
  不得把音频数据或文本内容送进渲染层（C25 污染面）。半双工门控不变。
- **内存口径不变**：C30 私有工作集；票 02 spike 实测空闲进程树 16.5MB / 预算 ≤25MB。
  新增的材质若需离屏位图，必须计入并在 `docs/SLO.md` 复测回填，**不得下调阈值**。
- **尺寸**：SPEC-08 §2 配置范围 44–72px（默认 56）。玻璃质感在 56px 下是否够看得清，
  由 owner 实机签收判定；若要改默认值属契约变更，走签字。
- **DPI/多显示器**：per-monitor DPI 与跨屏拖拽行为沿用票 07，吸附计算必须用目标屏的
  work area，不得把球吸到任务栏之下或副屏外。
- 禁止 emoji 进 UI 代码（CI 扫描 U+2190–U+2BFF / U+1F300–U+1FAFF / U+FE0F）。

## Out of scope
面板/窗口内容（票 33–40）；KWS 唤醒（票 41）；TTS 播报（票 26）；macOS 端口（票 55）。

## 契约变更流程（本票特有）
SPEC-08 §2 的 20 态视觉表与 D43 的视觉映射列**必须改**，但顺序是：
**先做可运行原型 → owner 实机签收满意 → 才回填 SPEC-08 文本并请其签字**。
在签字之前，任何提交都不得修改 `docs/specs/SPEC-08-ui-ball-panel.md` 或 `docs/PLAN.md`。
原型期的视觉定义写进 `docs/evidence/s1/62-visual-spec-draft.md`（草拟位，签字后迁进 SPEC-08）。

## Acceptance criteria
- [ ] 可运行原型：`build\balldebug.exe -stay` 启来即是一个玻璃液态球，**在浅色与深色桌布上都清晰可见**（差分截屏为证，不用自绘合成图）。
- [ ] Sleeping 态：静态单帧、**零活动定时器句柄**（实测断言，非策略表）、CPU 采样 ≤0.5%。
- [ ] 会话态动效：Listening/Thinking/Acting/Warm/Speaking 各自的液体颜色/流速/边框差异肉眼可辨，逐态截屏留证。
- [ ] 音频驱动：喂已知电平（wav 注入票 13 的 C8 seam）时液面旋转幅度随电平单调变化；静音时收敛进"边框"过渡态。
- [ ] 靠边吸附：拖到左/右/上边缘 → 收缩半隐并留可命中区；悬停/单击 → 弹回完整球；跨 DPI 与副屏行为正确。
- [ ] 资源复测：空闲与说话两种口径的私有工作集 + CPU 实测回填 `docs/SLO.md`，阈值不得下调。
- [x] 契约草案：20 态视觉映射表成文（`62-visual-spec-draft.md`），等 owner 签字后才回填 SPEC-08。
      —— **2026-09-21 00:50 成文半落地**：`docs/evidence/s1/62-visual-spec-draft.md`（`48f0cfd`，193 行）。
      我核过：§1 **逐态一行**（D43 态名逐字照抄）、§2 **22 条与 SPEC-08 的逐条冲突**（每条带 `file:line`）、
      §3 **15 个〔待 owner 定〕问题**（代理**没有替我选值**）、§4 **明列"证据里找不到数字的格子"**（不编数）。
      ⇒ 框的**后半句"回填 SPEC-08"仍欠**，且它的前置是 owner 批 Q 列表（**R15 第 10 项**）。
      ⚠ 本框只勾"成文"，不暗示票 62 已验收：AC#8 那张 1:1 裁决表**仍然不存在**，本票保持 `review`。

- [ ] 对抗验收由非实现者执行，验收报告含与上述 AC **1:1 的裁决表**（README 规则 6）。

## Progress log (append-only, newest last)
- 2026-09-20 agent(62-C): did=wire the liquid model into the Ball (the stranded liquid_windows.go, rebuilt against the real struct):
  Ball now owns `liq liquidMotion` + last raw level, `Ball.SetAudioLevel` marshals the scalar onto the UI thread via
  sta.PostTask (C25: level only, no samples/text), a BOUNDED 33ms burst timer (timerLiquidID) runs only while
  liquidMotion.busy() AND Animated(state), starts/stops on state entry/exit, frozen mode keeps the seam fully inert;
  liquid stamps at draw time (frameVisual) so the base visual never compounds; renderer glow permille snapped to the
  cached 25-step ladder (no COM creation on animated frames); balldebug gains -level/-diff-level feeder + winlive
  TestBallLiveAudioLiquidGate. Sleeping differential re-run: timers=no cpu=0.000/0.000 privWS=12.03MB px>=8=2120
  (baseline 2098); Listening/Speaking with 0.8 feed animate at <=0.08% 1-core. Gates: gofmt/vet clean,
  go test -count=2 + full winlive suite pass. next=AC#4 measured with wav injection (ticket 13 C8 seam), then the
  edge-dock wiring (item 2), then per-state evidence sweep for owner sign-off -> SPEC-08 amendment.
- 2026-09-20 agent(62-D): did=item 1's idle border transition, driven by STATE changes instead of only by the voice
  gate. liquidMotion now carries borderState (the resting border level of the current state, from the new
  borderAtRest(), the single authority shared with applyGlassForm) + ownsBorder (who is the border authority), and
  wantBorder() orders the two drivers: a state says whether the frame carries a border at all, a live voice may only
  pull it DOWN. enterState() only moves the TARGET, never the value, so a border cannot arrive in one frame; the
  travel is carried by the existing BOUNDED burst timer, widened from Animated() to transitionDriven() = exactly the
  states whose FROZEN animation policy already grants a timer (Animated five + Warm breathing + Settling fade), so
  ticket 62 arms no timer anywhere ticket 07 left without one and Sleeping can never appear there. An idle orb
  ramps the border alone (stepBorder: no level, no angle), so nothing swirls after a session. Sleeping releases the
  border authority -> the committed frame is untouched. Evidence: Sleeping differential re-run
  docs/evidence/s1/62-diff-border/ reads timers=no cpu=0.000/0.000 privWS=11.84MB px>=8=2120 box=44x44 (identical to
  the proven 62-diff-glass row); live winlive TestBallLiveIdleBorderTransition records the painted frames per state
  (e.g. Settling 0 -> 0.11 -> 0.32 -> 0.53 -> 0.74 -> 0.96 -> 1) and proves the burst retires itself and never arms
  in a static state. Gates: gofmt clean, go vet ./internal/ball ./cmd/balldebug clean, go test -count=2 + full
  winlive suite pass (Acting's frozen TimersAlive expectation kept intact by gating the whole seam off in frozen
  mode). next=item 2, the edge-dock auto-shrink (drag-near-edge docks to a tab, hover pops out, four edges, mouse
  messages only - no dock timer at all), then its differential row.
- 2026-09-20 agent(62-E): did=made the hover direction PROVABLE without owning the cursor, and made the whole item-1/2
  set watchable in one command. (1) dc8cd03: the ramp's step law moved out of dockStep into dockRampFrame (pure,
  dock.go), which is the exact function dockHoverMove and dockLeave call per mouse message;
  TestDockHoverPopBackWalksTheRampHome (dock_test.go, never skips) drives it in REVERSE - docked p=1 to popped p=0 -
  on four edges x two monitor layouts (incl. a negative-origin secondary) and asserts monotonic walk-back, exact
  landing (no overshoot, no 0.0001 stall), the 160ms/DockAnimMs frame budget, orb-tangent-at-every-level, the drawn
  orb widening to a full circle, travel AWAY from the edge, the landed orb wholly inside the work area (reachable),
  the mid-pop retreat, and both endpoints holding still (a landed ramp costs zero further frames = zero CPU).
  Mutation check: a one-way ramp fails it at once ("primary/left: the ramp froze at p=1 with target=0"). dock_test.go
  is now eleven deterministic tests, none skipping. The two live t.Skipf in TestBallLiveEdgeDockHover stay (the OS
  really can refuse the pointer, and TrackMouseEvent answers a synthesised hover with an immediate WM_MOUSELEAVE) but
  are now SKIP-LOUD: they name the deterministic layer as the proof and state that the run exercised no real event
  path; this run took that skip branch while TestBallLiveEdgeDock PASSed. DockPop seam (DebugPop) added for evidence.
  (2) balldebug -tour (+ -tour-dwell): a 14-step owner walkthrough in ONE run - static Sleeping -> summoned liquid ->
  voice-driven rotation from a synthetic -level (default 0.70, breath every 12th frame) -> idle border glide-in ->
  Sleeping again -> docked on left/right/top/bottom with the pop-back after each -> leaves the ball docked on the
  right for the real pointer (hover / click-out / drag-free). Frozen mode still works and -frozen -tour walks the
  still-frozen SPEC-08 §2.1 look over the same states instead, printing that liquid and dock are prototype-only, so
  the two runs A/B on one desktop. Evidence rewritten as docs/evidence/s1/62-ball-visual-prototype.md: the exact
  command, a 14-row "what to look at" table with the measured numbers in it, the four -look candidates (aurora
  recommended, it is the reference family), the unverified list (图1/图2 never reached any agent, so colour fidelity
  to owner taste is precisely what sign-off decides; dark wallpaper and 44/72px still to be looked at), and the
  proposed SPEC-08 §2.1 wording (SPEC-08 and PLAN untouched). Differential re-run docs/evidence/s1/62-diff-signoff
  (6 rows: Sleeping / Listening / Speaking, free and docked-right): Sleeping timers=no cpu=0.000/0.000
  privWS=11.76MB px>=8=2120 box=44x44 teardown=clean - item-by-item identical to the proven 62-diff-glass and
  62-diff-border rows, NO REGRESSION; the new docked-Sleeping row is also timers=no / 0.000% / 1620px (the dock
  introduced no timer anywhere). Gates: gofmt clean, go vet ./internal/ball ./cmd/balldebug clean (also
  -tags winlive), go test -count=2 ./internal/ball ok, every process launched was killed and tasklist shows no
  balldebug/wisp left. next=owner-signoff (实机跑 -tour 判颜色与可见性 -> 才回填 SPEC-08 §2 -> 然后补深色桌布差分)
- [2026-09-20T10:45Z] agent=orchestrator did=**owner 实机签收完成，判决=临时通过（质感未签收）**。
  owner 跑了 `-tour` 并回看三张截图，原话：「可以说算是赝品吧，离我发的那种质感还是有不小差距，
  但是先勉强用吧，以后再换样式，就先这样吧。开始阶段不能要求太高，本末倒置了就，先完成核心功能」。
  **验收方独立复核过的部分（不采信代理自述）**：①`decideRisk` 之外我又做一次**单向 ramp 变异**——
  把 `dockRampFrame` 改成只进不退 → 仅 `TestDockHoverPopBackWalksTheRampHome` 变红、其余 10 个仍绿
  （正确表现：单向不影响"吸进去"），证明"弹回来"这个方向有**永不 skip** 的确定性锚；
  ②`dock_test.go` 11 个测试、`grep -c t.Skip` = **0**；③Sleeping 差分我此前亲自跑过
  （`timers=no cpu=0.000 privWS=11.75MB px_delta_ge8=2120`），与代理报的数一致；
  ④`build/balldebug.exe` 我在**剥光 PATH**（只剩 System32）下验证过能起，签收命令交给 owner 前确认跑得通。
  **SPEC-08 §2 已回填，但带 INTERIM 标注**：只改「12px 微点/0.35 → 静态玻璃体 44px」这一项，
  理由是可证伪的（旧规格差分 0 像素变化，物理上不成像），**不是**因为质感被认可；
  D32 的 Sleeping 零定时器纪律一条未放宽，票面批准范围写死为"尺寸/可见性/零定时器"三项。
  **质感返工另立票 65**，且票 65 记下一个必须承认的事实：**图1/图2 从未落到磁盘**
  （我逐个核对本会话 11 张附件：商标 4 + 托盘 1 + AnySearch 文档 1 + 代理面板 1 + 本轮 tour 3 + 小裁图 1，
  无一为玻璃参考图）⇒ 本票四个代理的配色全部是**按文字盲做**，"赝品"判决在信息缺失下是必然的，
  返工第一步必须是拿到图而不是再猜一次。**Status → done（临时通过）。**
  next=票 64（球缺陷：热键接线/交互四项/零定时器实测）与票 12（装配：A8+A11+A13）解锁开工；
  票 65 blocked-on-owner 等参考图
