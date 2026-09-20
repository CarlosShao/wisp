# 12 — S1 gate: `wisp run` text path end-to-end, notify/list_tools, settle, RSS gate

**Status:** review
**Claimed by:** agent-ticket12-assembly (165 次调用后撞 150 轮上限) + orchestrator（票面对账）
**Last update:** 2026-09-20T12:30Z
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
- [ ] Ball state walk Sleeping→Listening→Thinking→Acting→Speaking/notify→Warm→Settling→Sleeping
      matches D43 rows #4,12,15,18,26,31,32 (state log captured).
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
