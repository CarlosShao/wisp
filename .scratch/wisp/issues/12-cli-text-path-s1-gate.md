# 12 — S1 gate: `wisp run` text path end-to-end, notify/list_tools, settle, RSS gate

**Status:** review
**Claimed by:** agent-ticket12-assembly (165 次调用后撞 150 轮上限) + orchestrator（票面对账）
**Last update:** 2026-09-20T13:40Z
**Blocked by:** 04-sqlite-core, 06-secretstore-envs, 07-ball-state-machine-core,
08-observability-slo, 10-agent-loop-core
**Parallel slots:** ≤2 sub-agents (A: CLI + notify/list_tools tools; B: end-to-end wiring +
SLO verification run)
**Spec refs:** S1 done criteria (§4/D44), D12, D34 (notify, list_tools), D32 SLO, C21 native

## What to build
The S1 vertical slice: `wisp run "..."` CLI (skips voice per D12), `notify` and `list_tools` as
the first two builtin tools, hotkey→(placeholder text capture via CLI for S1)→agent→TTS-less
result via system notification, 3s Settling→unload, and the measured idle SLO gate (tree RSS
≤25MB path-Y / ≤40MB path-X per spike verdict in 02). This is the first demoable product
moment: type → reply → notification → memory returns to idle.

## Key constraints
- CLI: `wisp run "task text"` attaches console, prints streaming text result + final status +
  cost line; exit code reflects error_class. CLI shares the agent loop (no voice, no ball pop
  unless GUI mode).
- `notify` tool (L0, D10 dependency): posts a Windows toast; Focus-Assist degradation contract
  (always also update ball badge) stubbed at ball-level now (real rule at 45).
- `list_tools` (L0): returns resident + third-party tool directory (D15② fallback path).
- Settle: result presentation done → 3s → DisposalScope session teardown (incl. FreeOSMemory) →
  RSS back to idle cap ≤10s.
- SLO measured with 08's sampler on a real Windows desktop session: **Spike backfill (T02, 2026-09-19): idle measured 16.5/16.6MB on Path Y — the 25MB gate is
      realistic BUT session settle leaves ~+15MB residual (Sleeping ~31MB): THIS ticket must include
      `debug.SetMemoryLimit`/GOGC tuning and re-measure against 25MB; handle cap for gate = <600
      (orchestrator ruling in docs/SLO.md), GDI <200 unchanged.** idle tree-private ≤ spike
  verdict value; goroutines ≤6; handles <300; GDI <200; zero periodic disk writes; zero long
  network connections in idle.
- C21 native-side token table complete (ball + future panel share); tokens doc updated.
- task_log row per CLI task; tool_call rows for notify/list_tools (correlationId present).
- GUI mode end-to-end smoke: hotkey → ball Listening (text captured via CLI in this slice;
  mic arrives S2) → reply → notify → Warm/Settling → Sleeping.

## Out of scope
- Voice anything (13+); panel (33+); gating on non-notify tools (21+); memory (29).

## Acceptance criteria
- [x] `wisp run "总结一下…"` against mockllm: streamed reply, notification posted, exit 0;
      with fail_next(3) → proper error_class + non-zero exit + user-visible message.
      —— 证据（编排者按代码对账时补，代理未写）：`cmd/wisp/run_test.go`
      `TestRunTextTaskTextPathEndToEnd` + `TestRunTextTaskFailNextIsClassified`，
      我亲自跑 `go test -run 'Providers|RunText|Host|Composed' ./cmd/wisp/` → **ok 13.794s**。
- [ ] Idle SLO gate PASS recorded into docs/SLO.md appendix (numbers + machine + date), using
      tree-private metric; settle ≤10s verified with FreeOSMemory counter >0.
      —— 2026-09-20 桌面独占跑完，数字见 `docs/SLO.md` 附录 A；**本框不勾，因为 CPU 那一行是真 FAIL**。
      已 PASS（每个样本都列，没挑幸运值）：私有工作集中位 8.43/8.54/8.52/9.10/8.50MB ≤25MB ·
      句柄 max 229–249 <600 · GDI 0 <200 · goroutine 1 ≤6 · 写盘 0 · TCP 0 ·
      **settle 262/260/263ms 回 cap 且 FreeOSMemoryCount=2>0（3/3，AC 字面判据满足）**。
      **FAIL：`Sleeping` CPU ≤0.5% —— 6 个样本 2 个超线：0.5450%（10s 窗）/ 0.6285%（60s 窗，
      正是 D32 写的 1min 窗口）**；阈值未改、未复测取巧。
      带 44px 新球体（今天 SPEC-08 §2 的 INTERIM 改动）**独立复测装得下**：球在 `Sleeping`
      私有工作集 **11.79MB**、`timers=no`、CPU 0.000%（树外测）、句柄 358、GDI 4、差分成像
      **px≥8/255 = 2103 像素 / 框 46×46**（旧 12px/0.35 同位置 0 像素）。
      **两起测量工具缺陷（本次未改任何代码）**：①`internal/proc/treemetrics_windows.go::parseSystemProcesses`
      先判 `next==0 break` 再解析当前项 → **快照最后一个进程永远丢失**，而新建进程正在链尾，
      于是「自己量自己」的 `wisp slo -state X` 常态 fail-closed exit=2；差分实验决定性：同命令同二进制，
      唯一差别是「启动后再造一个进程」，`exp-a` exit=2 → `exp-b` exit=0。**连带事实：票 08 归档的
      `build/slo/slo-report.json`（09-19T23:35Z）Sleeping/Warm 全是 exit=2 —— state 口径的 SLO 门
      从未产出过一份样本窗记录，`scripts/slo-check.ps1` 今天在 CI 上必红。**
      ②采样器跑在被测进程内 → CPU 门测的是观测者自己（118 reads/60s=0.629%、40/10s=0.545%、
      15/30s=0.087%，每次读折 16–38ms CPU），故该 FAIL **今天不可判**（≠不达标）。
      **R12 答复：不拆开就不是验收判据**——建议拆 **AC#2a（CI 子集：内存/句柄/GDI/goroutine/写盘/TCP/settle，
      修好①即可每次 merge 前跑）** 与 **AC#2b（桌面独占全量：CPU + 带球口径，只能是人的仪式；
      附录 A 即第一次正式记录）**，详见 `docs/SLO.md` A.5。
      —— **编排者复跑裁定（2026-09-20 21:44–21:56，`docs/SLO.md` 附录 B；原始样本 `build/slo/verify/`）**：
      ① 独立复现成立（静态读码 + 我这次无 keeper → `exit=2`、加 keeper → 出数）；但代理那次记的
      `exp-b exit=0` 我复跑是 **`exit=1`**（采样成功、门 FAIL），所以「跑通 ≠ PASS」。
      ② 方向对、**解释形式要换**：CPU ∝ 读次数**不成立**（我 40 次读 0.0387% 反而比 5 次读 0.1032% 省），
      正确形式是**单次读的绝对成本 ≈1.3ms**（250ms 间隔跳在 0.506–0.528%，2000ms 间隔跳在 0.0649–0.0652%，
      两处绝对值相同）⇒ **250ms 间隔下一次读本身就 0.52% > 0.5% 门**，这条门在树内口径下**没有定义**。
      ③ 同配置重复跑 PASS/FAIL 随机翻转（5 个 10s/250ms：0.0387 通 / 0.5838 红 / 0.5837 红 / 0.5053 红 / 0.2076 通）
      ⇒ **只修①会把 `slo-check.ps1` 从"必红"变成"约一半概率红"，比确定性地坏更糟**，故 AC#2a 那句
      「修好①即可每次 merge 前跑」**作废**，两起缺陷绑定在**票 66** 同一 commit 集里修。
      ④ 被测进程清白：40 样本里 37 个 CPU 恰为 `0.0000`，与球体树外实测 0.000% 互证。
      **裁定：CPU 行既不记 PASS 也不记 FAIL，记「仪器未定义」；D32 的 0.5% 一字未动；本框继续未勾，
      解除条件 = 票 66 的 AC#1–AC#4 全绿。** 采纳 AC#2a/2b 拆分（bookkeeping，非改门）。
- [x] Ball state walk Sleeping→Listening→Thinking→Acting→Speaking/notify→Warm→Settling→Sleeping
      matches D43 rows #4,12,15,18,26,31,32 (state log captured).
      —— **2026-09-20 桌面独占跑完。state log 在 `docs/evidence/s1/12-ball-walk/`**：
      `per-state-hold.log`（逐态 `state=X timers=bool` 真机日志，7 行全在）、
      `one-process-cycle.log`（**单进程**活窗口连走 20 态 + 句柄轨迹 380→439→closed 439，`OK`=退出前
      回 Sleeping 零定时器断言通过）、`diff-table.txt`（7 态差分成像 + 资源行，teardown 全 clean）、
      `winline.log`（56 PASS / **0 SKIP**，含 `requireQuietBallDesktop` 洁净前提）、
      `statemachine-rows.log`（`TestEveryLegalRowFires` ×2 = 120 PASS / 0 FAIL）。
      **逐行核对**：#4 单击/快捷键→`Listening`+phase FirstRound+effects `session.scope-create`/`speech.load-vad-asr`
      ——机器侧 `#4_sleeping_summon` PASS，**且真机由一条合成点击消息驱动**（winlive
      `TestLiveClickSummonsAndDragDoesNot` 断言机器与球同时进 Listening）；#12 三变体（首轮→Sleeping /
      会话内→Warm / Conversation→Warm）PASS，15s/90s/30s 冻结在 `timeouts.go`；#15 双分支
      （`HasToolCall`→Acting / 纯文本→Speaking + `panel.stream-push`）PASS；#18 双分支
      （有播报→Speaking / 静默→Settling）PASS；#26 →`Warm` + `mic.keep-off`/`panel.keep-alive` PASS；
      #31 90s→`Settling` + `session.dispose-scope`/`mem.free-os-memory` PASS，其释放步由
      `wisp slo -settle` **实测**（3/3 回 cap，FreeOSMemoryCount=2）；#32 3s→`Sleeping`(KWS 关)/`Armed`(开)
      + `mem.rss-verify-10s` PASS，「10s 内达标」即 settle 样本的 260–263ms。
      球侧动效纪律逐态对上 SPEC-08 §2：`Sleeping timers=no`，`Listening/Thinking/Acting/Speaking/Warm/Settling
      timers=yes`（Warm/Settling 由 `TestBallLiveLifecycle` 的策略断言覆盖，CPU 0.001–0.003% 全核证明不是 30fps 循环）。
      ⚠ **一条边界必须写清**：#12/#15/#18/#26/#31/#32 的**触发源**（VAD 判停/LLM 首 token/工具完成/TTS 播完/
      90s 空闲）在 S1 还没接进 `balldebug`——球只接手势事件，所以这 6 行是「机器按行转移 + 球按态渲染」
      两半证的，**不是**同一个活进程里由 `Dispatch` 串完整条会话链跑出来的。要补那一跑需给 `balldebug`
      加事件注入口 = 代码改动，本次按纪律**没做**，登记为发现（口径同 AC#2 的②）。
      `Speaking/notify` 的 toast 侧证据沿用 AC#1（`cmd/wisp/run_test.go` 已勾）；球徽标仍是 45 票前的桩。
- [x] task_log + tool_call rows written with correlation_id; visible via sqlite query.
      —— `internal/tools/loop_approval_test.go::TestHostDispatchThroughTheAssembledBridge`
      （真实 SQLite 回读）+ `TestToolCallRowsAreComplete`；我跑 `./internal/tools/` → ok 10.392s。
- [x] Key comes from the store, not the config: the provider that serves the `wisp run` request is
      built by `internal/llm` from config whose `api_key_ref = "secret:<id>"` resolves through the
      DPAPI store, and the test asserts the outbound `Authorization` header the **provider** set
      (not one a test-written client set). Handed here from ticket 63 AC#6, whose end-to-end proves
      ref→`config.LoadFile`→`ProviderKeys` but issues the request with its own http client.
      Completion criterion: a red-when-broken test in this ticket's own suite; mutation = point the
      ref at a missing blob and the run must fail with the Unconfigured path, not silently succeed.
      —— **A8 就此闭环**：`cmd/wisp/run_test.go::TestRunTextTaskKeyResolvesInTheStore`（provider 自己
      发出请求、key 从 DPAPI 存储解析）+ `TestMissingBlobFailsUnconfiguredNeverSilently`
      （正是上面那条变异判据，ref 指向缺失 blob 时走 Unconfigured 而非静默成功）。
      我亲自跑通（含在 ok 13.794s 那一批里）。
- [x] The capability probe is actually called from the composition root: `cmd/wisp`'s provider
      save/discovery path invokes `llm.RunProbeSuite` (the `memoryHealthSink` in
      `internal/llm/probe_health_test.go` is a ready-made adapter), so the 「声明 ✓ / 实测 ✗」
      event can fire on a real machine. Handed here from ticket 11 AC#6, registered as **A11**.
      Plus: a mockllm thinking-capability mode and a thinking probe that REQUIRES a reasoning
      delta — today ticket 09's thinking check accepts a plain text answer, so it cannot detect a
      broken thinker. Completion criterion: both directions pinned (capable and broken) with the
      server's own request counters, mirroring how the fc/vision cases already work.
      —— **A11 就此闭环**：`cmd/wisp/providers_test.go` 的
      `TestProvidersProbeRecordsMeasuredThinkingTrue` / `...False`（**正反两向都钉**，
      正是上面那条完成判据）、`TestProvidersDiscoverListsWhatTheServerServes`
      （探针由 provider 保存/发现路径调用，非测试自建 client）、
      `TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless`；
      mockllm 侧新增 thinking 能力档（`tools/mockllm/thinking_capability_test.go`）。
      ⚠ 这些是代理**未写进票面**、由编排者按代码与 commit 对账补记的（见 Progress log 的疏漏记录）。
- [ ] Zero emoji scan over new UI strings passes; tokens doc updated with any additions.
- [ ] Human visual acceptance of ball states (user signs off screenshots — D29 rule).

## Progress log (append-only, newest last)
- [2026-09-20T11:10Z] agent=agent-ticket12-assembly did=(据 commit 反推，代理本人未写任何 log 行) 装配根 + stopguard 替换 + A8 + A11 + AC#1/AC#4，3 个 commit：`cbdea7c`（守卫与替代同批：删 `loop.decideRisk` 的 L1/L2 先拒，同 commit 引入 `Options.AdmitTask` 与 295 行 `internal/tools/loop_approval_test.go`，正反两向都有——有 gate 时 L1 写开阻断窗约 3s 并落文件记 `fs.write/L1/allow/success`，无 gate 时同一调用被拒记 `L1/reject` 且文件不出现）、`cd011b8`（`cmd/wisp/run.go` 同时 import `internal/agent/approval` 与 `internal/tools` ⇒ **A13 打通**）、以及死前未提交的 thinking 探测半成品 next=(未写)
- [2026-09-20T12:35Z] agent=orchestrator did=**票面对账 + 一项流程疏漏归己**。①发现该代理跑完 165 次工具调用后**票面 8 个 AC 框一个没勾、Progress log 一行没写**，`Status` 还停在 `ready-for-agent` ⇒ 任何后续代理读这张票都会**把 165 次调用重做一遍**，正是本仓 `-done` 后缀机制要防的事。**根因是我的简报漏项**：我给票 11、票 64 的简报都硬写了"每次 commit 追加一行 Progress log"，给票 12 的那份只写了 commit 纪律、漏了这条。已把"票面必须与 commit 同步"并入偏好。②按代码而非自述逐条对账并补勾 4 框：AC#1（`run_test.go` 端到端 + fail_next 分类）、AC#4（真实 SQLite 回读）、**AC#5=A8**（`TestRunTextTaskKeyResolvesInTheStore` + `TestMissingBlobFailsUnconfiguredNeverSilently`，后者正是我写进完成判据的那条变异）、**AC#6=A11**（thinking 探测正反两向 `...MeasuredThinkingTrue/False` + 由 provider 发现路径调用，非测试自建 client）。③把它未提交的半成品**验证后**存成检查点 `f994ca5`（gofmt/vet 干净，mockllm ok 0.128s、`internal/llm` ok 34.041s、`cmd/wisp` ok 4.998s，无禁用断言、无遗留变异）。④我自己复核过的高风险项：守卫删除与替代**确在同一 commit**（`git show --stat cbdea7c` 只含 `loop.go` 与该测试文件）。**A8、A11、A13 三条遗留全部闭环**，registry 待更新。仍留 4 框未勾：AC#2 空闲 SLO（要安静机器，须单独独占跑）、AC#3 球状态走位（要桌面，且票 64 正在 `internal/ball` 上跑）、AC#7 emoji/tokens 文档扫描（d22scan 已 clean，tokens 文档那半没做）、AC#8 人工视觉签收（owner，且票 62 那次是 INTERIM 临时通过）。next=派接续代理做 AC#7 + 复核 AC#1/4 的 fail 路径细节；AC#2/3/8 由我按安静窗口与 owner 时间安排
- [2026-09-20T13:40Z] agent=agent-ticket12-desktop-gate did=**桌面独占跑完 AC#2 + AC#3（本会话是唯一 agent，全程 `tasklist` 复核无残留）**。①**AC#3 勾**：state log 落在 `docs/evidence/s1/12-ball-walk/`（逐态 `timers=` 真机日志 + **单进程** 20 态走位日志 + 7 态差分化表 + winlive 56 PASS/0 SKIP + `TestEveryLegalRowFires` ×2 120 PASS/0 FAIL），#4/12/15/18/26/31/32 **逐行**按 from/event/guard/effects 对上，#4 另有一条真点击消息驱动机器+球同进 `Listening`，#31/#32 的副作用行由 `wisp slo -settle` 实测（3/3 回 cap，260–263ms，FreeOSMemoryCount=2）。⚠ 明写边界：#12/#15/#18/#26/#31/#32 的**触发源**在 S1 未接进 `balldebug`（球只接手势事件），这 6 行是「机器按行转移 + 球按态渲染」两半证的，**不是**一个活进程里由 `Dispatch` 串完整条链——补那一跑要给 harness 加事件注入口=代码改动，按纪律没做，登记为发现。②**AC#2 不勾**：七项 PASS（私有工作集中位 8.43–9.10MB、句柄 229–249、GDI 0、goroutine 1、写盘 0、TCP 0、settle 3/3 且计数>0），**CPU 行 FAIL 照录**（6 样本 2 超线：0.5450% / 0.6285%——后者是 D32 原定的 1min 窗口长度），未改阈值未挑样本；今天改过的 44px 玻璃体**独立复测装得下**：球 `Sleeping` 11.79MB / `timers=no` / CPU 0.000%（树外测）/ 句柄 358 / GDI 4 / 差分成像 px≥8/255=2103、框 46×46。③**挖出两起测量工具缺陷，未动代码**：**缺陷①** `internal/proc/treemetrics_windows.go::parseSystemProcesses` 先 `if next==0 break` 后解析，**快照最后一个进程永远进不了 map**，而新建进程恰在链尾 → 「自己量自己」的 `wisp slo -state X` 常态 fail-closed exit=2；决定性差分实验：同命令同二进制只差「启动后再造一个进程」，`exp-a` exit=2 → `exp-b` exit=0 且九项门全出数；兄弟实现 `cmd/balldebug/diff_windows.go::privateWorkingSetFor` 顺序是对的（先比 pid 再 break），修法就是把当前项解析挪到 `next==0` 判断之前。**连带事实（对票 08 的追溯）**：`build/slo/slo-report.json`（09-19T23:35Z，票 08 归档跑）`Sleeping exit=2 / Warm exit=2 / all_pass=false`，本次 `slo-check.ps1` 复现同一失败并 FATAL 在 `-leak` 自检上——**即 state 口径的 SLO 门从未产出过一份样本窗记录，`slo-check.ps1` 今天在 CI 上必红**；**缺陷②** 采样器在被测进程内，CPU 门测的是观测者自己（118 reads/60s=0.629%、40/10s=0.545%、15/30s=0.087%，每次读折 16–38ms），故那条 FAIL **不可判**而非不达标。④**R12 答复：不拆开就不算验收判据**——建议 AC#2 拆 **AC#2a（CI 子集：内存/句柄/GDI/goroutine/写盘/TCP/settle，修好①就能每次 merge 前跑）** + **AC#2b（桌面独占全量：CPU + 带球口径，只能是人的仪式，`docs/SLO.md` 附录 A 即第一次正式记录）**；重跑命令逐字抄附录 A.6（含①的 keeper 绕法，修好即删）。门禁：`go vet ./internal/ball/ ./cmd/balldebug/` 干净、`go test -count=2 ./internal/ball/` ok 0.068s、`go test -count=2 ./internal/statemachine/` ok、winlive ok 12.045s。next=编排者裁定 AC#2a/2b 拆分 + 缺陷①归票 08 返工（它是判据仪器本身坏了，优先级高于票面任何功能框）+ 缺陷②先定「从门里扣 self CPU 算不算改门」；AC#7（tokens 文档那半）与 AC#8（owner 签收）仍空
- [2026-09-20T14:00Z] agent=orchestrator did=**桌面跑复验 + 两起缺陷登记为票 66 + AC#2 CPU 行裁定**。①独立复现缺陷①（读码 `treemetrics_windows.go:240-243` + 无 keeper `exit=2`／加 keeper 出数）；更正代理那次「`exp-b` exit=0」——我复跑是 **exit=1**（采样通了、门 FAIL），「跑通≠PASS」。②复跑 8 个样本（`build/slo/verify/`，全列于 `docs/SLO.md` 附录 B）：缺陷②**方向成立但解释形式错了**，CPU 不随读次数单调（40 次读 0.0387% < 5 次读 0.1032%），真实形式是**单次读绝对成本 ≈1.3ms**（250ms→0.506–0.528%、2000ms→0.0649–0.0652%，两处绝对值互证）⇒ **250ms 间隔下一次读就 0.52% > 0.5% 门，该门在树内口径下无定义**；同配置 5 次里 3 红 2 绿 ⇒ **只修①会让 `slo-check.ps1` 变成约一半概率随机红**，故作废 AC#2a「修好①即可跑」的前提，两缺陷绑定进**票 66**。③产品侧清白：40 样本 37 个 CPU 恰 `0.0000` + 球树外 0.000%。④裁定：AC#2 **保持未勾**、CPU 行记「仪器未定义」（不记 PASS 也不记 FAIL）、**D32 0.5% 阈值一字未动**、附录 A 原文不改写（叠加附录 B）；采纳 AC#2a/2b 拆分并书面批准「树内 CPU 降级为记录项」（=票 66 裁定 2）。登记 **A14/A15**。`eaf808e` 已推 origin+cnb。next=票 66（测量类独占）→ 票 64 对抗验收 → 票 21 段 2；AC#7 tokens 文档半与 AC#8 owner 签收仍空

- [2026-09-20T14:00Z] agent=orchestrator did=**桌面跑复验 + 两起缺陷登记为票 66 + AC#2 CPU 行裁定**。①独立复现缺陷①（读码 `treemetrics_windows.go:240-243` + 无 keeper `exit=2`／加 keeper 出数）；更正代理那次「`exp-b` exit=0」——我复跑是 **exit=1**（采样通了、门 FAIL），「跑通≠PASS」。②复跑 8 个样本（`build/slo/verify/`，全列于 `docs/SLO.md` 附录 B）：缺陷②**方向成立但解释形式错了**，CPU 不随读次数单调（40 次读 0.0387
- [2026-09-20T14:24Z] agent=agent-ticket12-ac7 did=**AC#7 两半都跑了，两半都不绿 ⇒ 本框仍不勾**（票面其余框一字未动）。①**tokens 文档反向对账**（脚本比对 `tokens.css`/`tokens.go`/本表三方，非肉眼）：色值 **78 行（明 40 + 亮 38）三方逐字一致，零漂移**；`Palette` 40 个字段全在表中；文档里的 CSS 变量**没有一个**在 `tokens.css` 中不存在（无幽灵行）。缺口全在几何/动效段与票 62 色板，已全部补录：**新增 25 个常量行/列**（`SleepRestRatio`/`SleepingRestMinPx`/`RestSettledOpacity`、`LiquidRadiusA..C`/`LiquidOffsetA..C`、`GlassRimPx`/`GlassLipPx`/`GlassCaustic`/`BorderRingPx`、`SwimLevelGain`/`SpinLevelGain`、`BorderOpenMs`/`BorderCloseMs`/`SummonFlowMs`、`DockAnimMs`/`DockOverlapFrac`/`DockTriggerPx`（票 64）、`WarmBreathFPS`/`ThinkingBandPx`/`BadgeBorderPx`/`CountdownFontPx`）+ **新增「票 62 液态玻璃 look 色板」一节 36 行**（aurora/glacier/nebula/solar × 9 字段）。`docs/evidence/s1/c21-native-tokens.md` **133 → 223 行**。②**三条值漂移只登记不裁定**（写进该节末「值漂移」）：**D1 休眠体尺寸三方不一致**——SPEC-08 §2 INTERIM 写 44px、`statevisual.go::stateSize` 原型路径实算 `尺寸×0.62`（默认档 56→**34.72px**，下限 `SleepingRestMinPx`=30）、旧表写 12px；原值来自 `8464da6`（09-19），改动来自 `fd8f838`（09-20 票 62）。旧 12px 行**保留未删**并标「生产默认路径仍是此值」。**D2 INTERIM 的球体在生产路径不生效**——`prototypeVisuals` 默认 false，只有 `cmd/balldebug/main.go:104` 与 winlive 测试开它 ⇒ 默认构建下 `Sleeping` 仍画 12px 微点。**D3 票 62 的 36 个 look 色值没有 CSS 真相源**（`tokens.css` 无同名变量；`TestNoHardcodedColorsInBallPackage` 放过它，但「权威真相源=tokens.css」的契约对它不成立）。③**范围澄清**：`tokens.css` 129 条声明中 **80 条**是面板/screens 专用、原生无对应字段（`--bg-base`/`--ambient-*`/`--r-*`/`--s-*`/`--t-*` 等），按契约不进本表，已在「范围」段落点名并列全，避免下一个代理当成遗漏重做。④**零 emoji 扫描：AC 前提半假 + 门在 HEAD 不 clean**。(a) `tools/d22scan` **是独立 module**（`tools/d22scan/go.mod`），从仓库根 `go run ./tools/d22scan` 直接报 `main module (github.com/CarlosShao/wisp) does not contain package .../tools/d22scan`；正确调用 `cd tools/d22scan && go run . -root ../..` → **exit=1，1 条 finding**：`internal/llm/adaptertest/mockllm.go:68: [bare-goroutine] bare go func( is banned`（引入 commit `78b1466`，票 11 的测试脚手架；非 emoji）。⇒ **编排者 12:35Z 记的「d22scan 已 clean」在 HEAD `2825f28` 不成立**，本票按纪律未改 `internal/`。(b) emoji ban #8 的 scope **只有 `design/` 与 `frontend/`**（`main.go:137-142`），而 `frontend/` 在当前 HEAD **不存在** ⇒ 该门今天实际只覆盖 `design/`，**根本不扫 `internal/`/`cmd/` 的 Go 字符串字面量**——「over new UI strings」字面要求它做不到。scripts/ 下无第二个 emoji 扫描器（`grep -rln "emoji" scripts/` 只命中 `scripts/spike/bin/*.exe` 二进制）。(c) 用 d22scan 自己的 `emojiRe` 码位集补扫 `internal/`+`cmd/` 300 个 .go：**8 处命中，其中 3 处是用户可见字符串字面量**——`cmd/wisp/providers.go:201` `verdict := U+2717`、`:203` `verdict = U+2713`（由 `:209` `fmt.Fprintf(io_.out(), "  %s %s 实测 %s\n", ...)` 打到 stdout）、`:41` 帮助文本里的「声明 U+2713 / 实测 U+2717」；码位 U+2713/U+2717 正落在 ban #8 的 `2600-27BF` 段；引入 commit `cd011b8` = **本票 AC#1/A8 那批**。其余 5 处（`providers.go:7`、`providers_test.go:4`、`internal/llm/probe_health.go:17/120/203`）是注释，不成像。**按纪律一个 emoji 都没删**——改文案是 UI 改动，要签收。⑤门禁逐字：`go test -count=2 ./internal/ball/` → `ok github.com/CarlosShao/wisp/internal/ball 0.079s`；`gofmt -l internal/ball/` 空；`go vet ./internal/ball/` exit=0；全程未跑 `-tags winlive`、未建 `balldebug.exe`、未开任何窗口（本会话桌面零窗口）。⑥**票面完整性问题（未动）**：`12-cli-text-path-s1-gate.md` 文件末尾是 `14:00Z` 那条 orchestrator 记录的**截断副本**，止于「40 次读 0.0387」且无换行——重复登记属票面卫生，归写那条的人，我只在其后追加。next=**AC#7 不能勾**，要绿还差三件事：⑴ look 表 36 值与休眠体 44px/34.72px 由 owner 裁定后回填（改代码或改 SPEC 都不在本票权限）；⑵ emoji 门的真实覆盖面要么把 `internal/`+`cmd/` 纳入 ban #8（代码改动，且 `mockllm.go:68` 的 bare-goroutine 得先归票 11 处理，否则门永远红），要么把 AC#7 文字改成「design/ 零 emoji」——**编造绿色比留一个空框贵得多**；⑶ `providers.go` 三处 U+2713/U+2717 需 owner 决定换字（如 `PASS/FAIL`）还是登记为豁免；未定前本框保持未勾
- [2026-09-20T14:31Z] agent=agent-ticket12-ac7 did=**AC#7 检查点 2（自查回环，本框仍不勾）**。①**自证补录真的闭合了**：把「`tokens.go` 全部导出常量名 ↔ 表内反引号名」的脚本判据重跑一遍，`tokens.go` 55 个常量在 `docs/evidence/s1/c21-native-tokens.md` 中**全部可精确检索**（首轮我写成 `LiquidRadiusA..C` 区间缩写，脚本判 6 个未覆盖——已改成逐名 `LiquidRadiusA`/`B`/`C` + `LiquidOffsetA`/`B`/`C`，这样任何代理都能用 grep 复核而不是信我的措辞）。②**票 64 hotkey 面查过并写进「审计」**：`internal/ball/hotkey_windows.go` 三处 `const` 块是 Win32 修饰键/虚拟键码与 `HotkeyStatus` 枚举，无尺寸无颜色 ⇒ 票 64 进表的只有边缘吸附三常量；`hkNames` 用户可见文案零 emoji。③**自查我自己引入的 emoji**：用 d22scan 的 `emojiRe` 码位集扫我改过的两个文件——`c21-native-tokens.md` **0 命中**；本票文件里我这行原本带 4 个 U+2713/U+2717（引 `providers.go` 原句），已全部改写成码位名，**我自己没有向仓里新增任何 emoji 码位**（该文件余下命中在 101/122/134/141 行，均为他票既有记录，未动）。④**确认没有任何测试/工具扫 `docs/`**：`grep -rn "docs/evidence|\"docs\"" internal/ cmd/ tools/` 只命中三处注释引用，故本表新增的 36 个 `hex(...)` 字面量不会触发 `TestNoHardcodedColorsInBallPackage`（其 scope 是 `internal/ball/*.go`）。next=等 owner 裁「值漂移」三条 + 决定 `providers.go` 三处符号换字或登记豁免；AC#7 在此之前不勾
- [2026-09-20T14:38Z] agent=agent-ticket12-ac7 did=**AC#7 收尾复跑 + 票面卫生（本框仍不勾）**。①修掉自己上一条留下的问题：检查点 2 追加时在本 log 两条 AC#7 记录之间插进了一个空行（append-only 列表被断开，后来者按行扫 log 会漏读），已删；本文件其余既有内容一字未动，AC#1/3/4/5/6 的勾选状态与 AC#2/7/8 的未勾状态保持原样。②**最终门禁逐字复跑**：`go test -count=2 ./internal/ball/` → `ok github.com/CarlosShao/wisp/internal/ball 0.056s`；`gofmt -l internal/ball/ cmd/wisp/ tools/d22scan/` 空输出；`go vet ./internal/ball/` 干净；`cd tools/d22scan && go run . -root ../..` → `internal/llm/adaptertest/mockllm.go:68: [bare-goroutine] ...` + `d22scan: 1 finding(s)` + **exit=1**（emoji 命中 0 条，即 `design/` 侧确实干净，红的是别的 ban）。`docs/evidence/s1/c21-native-tokens.md` **133 → 227 行**，覆盖判据复核：`tokens.go` 55/55 导出常量、`Palette` 40/40 字段、look 色板 36/36 行全部可在表内精确检索。③**并发无碰撞**：本票两个 commit 期间，兄弟代理落了 `86e868d`（票 66 proc）、`5c6e661`（票 64 撤回 A16）、`15d9372`（票 66 docs）——`git show --name-only 15d9372` 与我的两个路径（本票文件、`c21-native-tokens.md`）**零交集**，我全程显式路径 add，未用 `git add -A`；工作区此刻还有他人半成品 `internal/tools/fs_test.go`(M) 与 `docs/evidence/s1/20-needs-assertion-restored.md`(??)，均未碰。next=**AC#7 判据未满足，缺三件**：⑴「值漂移」三条（休眠体 44px vs 34.72px vs 12px、INTERIM 球体在生产默认路径不生效、look 36 值无 CSS 真相源）要 owner 裁定，改代码或改 SPEC 都超出本票权限；⑵emoji 门的覆盖面：要么把 `internal/`+`cmd/` 纳入 ban #8（代码改动，且 `mockllm.go:68` 的 bare-goroutine 会先让门红，须归票 11），要么把 AC#7 文字改成「design/ 零 emoji」；⑶`cmd/wisp/providers.go:201/203/41` 三处 U+2713/U+2717 属本票 `cd011b8` 引入的用户可见 CLI 输出，换字是 UI 改动需签收，按纪律一个没删
