# 20 — host bridge + first tool family: fs.*, capability checks, spill rule

**Status:** review（两段 fs 工具族已落地并由编排者验收：`docs/evidence/s1/20-needs-assertion-restored.md`；
残口**从两条收窄到一条** —— 第 5 框「桥层 junction/8.3」已于 2026-09-21 用**真产物**闭合
（`internal/tools/bridge_junction_windows_test.go` + 证据 `docs/evidence/s1/20-bridge-junction-shortname.md`，
含吞掉下层拒绝的变异检验）；剩下 **A18** 真 `taskkill` 下的 `.wisp-tmp-*` 残留：已补特征化测试记录
今日真实行为（**每次 kill 留一个 `.wisp-tmp-*`，没有任何东西扫它**），但"清扫 vs 申报可接受残留"是行为设计，
仍归 owner ⇒ A18 保持 open）
**Claimed by:** T20-seg1-agent (host bridge + C4 registry), T20-seg2-agent (fs write family),
T20-box5-agent (桥层 junction/8.3 拒绝 + A18 特征化)
**Last update:** 2026-09-21 09:20 (票 20 第 5 框闭合，A18 补特征化)
⚠ **两处过期表述以本块为准（我逐条 grep 复现后更正）**：
① 正文里的 `who calls this in production? NOBODY YET, and it is NOT reachable from cmd/wisp`（见 §Progress）**已作废**——
`cd011b8 feat(12)` 之后 `cmd/wisp/run.go` 真的 `tools.New`(:262) + `approval.New`(:255) + 注册全部
`BuiltinFSEntries`，`decideRisk` 也已改成"有 `AdmitTask` 时 L1/L2 放行、无 gate 时仍拒"（`loop.go:761-774`），
端到端证据：`cmd/wisp/run_test.go:245` `TestHostDispatchThroughTheAssembledBridge`、`:330`
`TestComposedGateBlocksAWriteForTwoSeconds`。**⇒ 桥与 fs 工具今天在生产装配根里可达，A12/A13 那族缺陷对本票已解除。**
② 旧 Status 里的 "the loop-side gate composition remains" **也已完成**（同上）。
③ 另：本票曾丢掉一条判据——`0986d63` 改名那个注册测试时删掉了**逐条 `Needs` 断言**，
而 `Needs` 正是 `bridge.go:286/291` 做 C3 授权判定所依据的需求集（少报一个能力＝只授权 `fs.read`
的机器放行真写盘调用）。**已由编排者补回并做变异检验**（registry **A17**）。
**新发现的两个残口**：`[fs] allowed_dirs` 首次使用询问流**没有对应 AC 框**——
23:59 我把它从"零实现"**收窄**了（原文说过头， grep 当场推翻），真实状态见下面 ④；
以及 `6a6fca2` 是一个 `docs(` commit 却删了 `wiring_test.go` 9 行
（经查是死代码、无断言损失，但**判据文件不许挂在 docs 名下改动**）。
④ **`allowed_dirs` 这条到底缺什么（我逐条 grep 复现，别再重做已有部分）**：
  - **已存在且已测**：热加载方向引擎把它登记成 🔒 键（`internal/config/manager.go:362`），
    `internal/config/manager_test.go:162-225` 断言了 confirm hook 收到 `section="fs"` +
    `keys=["fs.allowed_dirs"]`、**拒绝确认后 `AllowedDirs` 仍为空**、批准后才生效、以及收紧方向。
  - **已在生产装配里被执行**：`cmd/wisp/run.go:223-227` 把 `cfg.FS.AllowedDirs` 喂给
    `tools.NewPathCanonicalizer(allowed, cfg.FS.ReparsePointExceptions)` ⇒ **这个键不是装饰品**，
    默认 `[]` 意味着**没有任何根可读可写**（fail-closed，与 D47 同向）。
  - **真正缺的只有这一件**：一次工具调用因为"路径不在任何 allowed 根下"被拒时，
    **没有任何东西把这次拒绝变成一个"要不要把这个目录加进来"的询问**。今天唯一的授权路径是
    用户手改 TOML（那条路径本身已经过 🔒 confirm）。
  ⇒ **安全后果：无**（缺的是"给得更多"的入口，不是"拦不住"）；**可用性后果：真实**，
    且第 209 行已记着一条硬约束：**它在 SPEC-12 §5 里没有行**，所以**不许在代码里标 `DEFERRED`**
    （那会打断 1:1 交叉核对）。⇒ 处置：**补一条 AC 框**（本票文末），GUI 完整版仍归票 39。

**Blocked by:** 17-risk-assessor-c19, 18-path-resolver-c26
**Parallel slots:** ≤2 sub-agents (A: host bridge + ToolProvider registry; B: fs tool family +
artifacts spill)
**Spec refs:** SPEC-07 §2–§3, D3, D14, D34, C1/C3/C4, D34 note②, D31 atomic-rename

## What to build
The `tools` module host bridge — the single choke point every capability flows through — with
C1 Tool contract, C3 capability enforcement (11 caps), C4 ToolProvider registry (Builtin now;
Manifest/Goja slots), and the first builtin family: fs.read/list/write/trash/move (+delete
gated). Includes the host-internal artifacts spill path and `list_tools`-resident registration
from 10.

## Key constraints
- Host bridge: every execution passes `risk.Assess` + capability check (undeclared capability →
  hard reject, not error), decision + rules_hit recorded into `tool_call` rows with
  correlation_id; plugin code can never bypass (no direct OS access from plugin runtimes —
  structurally guaranteed later by 51).
- C1 Tool: Name/Description/Parameters(JSON Schema)/Execute(ctx, params, onUpdate). Tool concurrency
  ≤4 per task; per-tool timeoutMs honored via ctx.
- fs family RiskLevels (D34): read/list L0 (in-allowlist; out-of-scope → L2 via R2);
  write-new L1; overwrite L2 (R8); move L1 (cross-volume → L2); **trash L1** (shell API —
  recycle bin); delete **not registered** unless `[fs] delete_enabled=true` → L2.
- All writes: temp file + atomic rename (D31 cancel semantics); cancel after partial work →
  report "steps possibly already applied" list (never fake clean cancel).
- Paths ALWAYS via PathResolver (18); results taint-marked (19).
- Host-internal artifact writes (`artifacts\`): separate internal API, NOT a Tool, no gating,
  data-dir-only, quota-managed (04), visible/deletable later via panel. Distinguishing rule
  documented + tested.
- `[fs] allowed_dirs` first-use ask flow: minimal native prompt now (full GUI at 39).

## Out of scope
- Approval UI/queue mechanics (21 wires the L1/L2 UX); web/system/doc/search tools (22–24);
  plugin manifests (50).

## Acceptance criteria
- [x] Bridge tests: undeclared capability → reject; risk decision → gate branch (L0 pass, L1
      pending-window callback, L2 pending-approval callback); tool_call rows complete.
- [x] fs matrix: each action's RiskLevel asserted incl. trash=L1, overwrite=L2, cross-volume
      move→L2, out-of-allowlist read→L2. Names: `TestD34WriteMatrix` (11 subtests),
      `TestCrossVolumeRuleIsR8` (the hardware-free layer),
      `TestOverwriteDetectionFollowsTheCanonicalPath`.
- [x] Atomic write test: kill mid-write → no partial file at target (temp+rename verified).
      Name: `TestAtomicWriteKillsMidWrite` (overwrite branch: old bytes survive byte-for-byte;
      new-file branch: nothing appears at the path; both assert no staging file is left behind).
- [x] Cancel-semantics test: cancelled after N ops → "applied steps" list correct.
      Names: `TestVetoAfterWorkStartedProducesAppliedSteps`,
      `TestLateVetoRendersTheApprovalLayersAppliedStepsReport`,
      `TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`,
      `TestApprovedL2ThatStopsMidWriteIsNotRenderedAsNeverStarted`.
- [x] Reparse/short-name target via junction → denied (18 integration).
      **桥层**证据（2026-09-21，真产物、零 `t.Skip`）：`internal/tools/bridge_junction_windows_test.go` —
      `TestBridgeRefusesARealJunctionOnTheReadRoute`（fs.read/fs.list 穿真 `mklink /J`）、
      `TestBridgeWritesNothingThroughARealJunction`（5 子项：write 覆写/新建、trash、move、delete，
      目标目录**逐文件 sha256 快照** before==after）、
      `TestJunctionInsideAnAllowedRootCannotReachAnAListFile`、
      `TestBridgeRefusesTheRealShortNameOfAnAListFile`（`GetShortPathNameW` 真短名 → Deny via R3、门零调用）、
      `TestShortNameSpellingGetsTheSameVerdictAsTheLongOne`、
      `TestWavedThroughJunctionWriteBooksAnApprovedRowThatWroteNothing`（真 SQLite 行）。
      审批门一律 **AnswerAllow** ⇒ 拒它的是桥不是人；阳性对照先证明 junction 真能读出那 94 字节。
      变异检验（`fs.go:95-98` + `fs_write.go:169-174` 吞掉下层拒绝）→ 6 条顶层 + 8 子项全红。
      证据：`docs/evidence/s1/20-bridge-junction-shortname.md`。
      ⚠ 这一框闭合**没有**顺手改任何守卫；§4 那三条上报（可批准的 L2 遮住 A 档、审计行记成
      "工具说不"、第 7 框要的原因枚举不存在）留给 owner 判。
- [ ] Artifacts API writes only under data dir; attempting user-dir write via artifacts API →
      rejected; `fs.write` to user dir remains gated (rule separation test).
- [ ] **第 7 框（2026-09-20 23:59 编排者补，来源：本票头部 ④）`[fs] allowed_dirs` 首次使用询问流的最小子集**。
      范围**只有这三条**，多一条就算越界（GUI 完整版归票 39，SPEC-12 §5 **没有这一行 ⇒ 代码里不许标 `DEFERRED`**）：
      - **(a) 拒绝原因必须可区分**：因"路径不在任何 allowed 根下"而被拒时，桥必须返回一个**机器可辨的**
        原因值，与 PathResolver 的红队拒绝（junction / 8.3 短名 / UNC / `\?\`）**不共用同一个字符串**。
        理由：后者**绝不能**变成询问——一份被注入的文档只要抛出一个穿过 junction 的路径，
        就能驱动我们去问用户"要不要授权这个目录"，那等于把红队门变成社工入口。
      - **(b) 同意必须走已有的 🔒 通道**：用户批准后，追加只允许经由
        `internal/config` 的 `fs.allowed_dirs` **方向引擎 confirm hook**（`manager.go:362` +
        `manager_test.go:162-225` 已测：拒绝后仍为空、批准后才生效）。
        **禁止**新增第二条授权路径、禁止直接改 `PathCanonicalizer` 的根集合、禁止"这次先放行下次再说"。
      - **(c) 三条断言**：①(a) 的两类拒绝原因可区分；②批准 → 下一次同路径调用成功，**拒绝 → 仍被拒**；
        ③**junction / 8.3 拒绝绝不产生询问**（这条是本 AC 的安全内核，缺它则本 AC 不成立）。
      **归属与阻塞**：本票的框**原本就没有编号**（按位置数：共 **7 框**。编排者 2026-09-20 写这句时实测
      `grep -c "^- \[x\]"`=4、`grep -c "^- \[ \]"`=3，未勾的是「桥层 junction/8.3」「artifacts spill」「本框」；
      **2026-09-21 第 5 框闭合后复测 = 5 / 2**，未勾的只剩「artifacts spill」与「本框」。
      日后做 1:1 裁决表时**按位置对齐**，别自己编号后又在别处引用错号）。
      另有一条**不在框里**的残口：**A18** 真 `taskkill` 下的 `.wisp-tmp-*` 残留 —— 2026-09-21 已补
      特征化测试（`internal/tools/bridge_a18_kill_windows_test.go`：真 kill 下目标要么完整要么不存在，
      但**每次 kill 留一个 `.wisp-tmp-*`，新建的桥/下一次成功写盘都不扫它**），
      "自愈清扫 vs 申报可接受残留"仍归 owner 判 ⇒ A18 继续 open。
      本框与 **A19**（审批层没有输入设备）同源，最小实现只能用 native prompt。
      ⚠ **不得与票 70 的全仓格式化并发**（要改 `internal/tools`）。


## Progress log (append-only, newest last)
- [2026-09-20T09:55Z] agent=T20-seg1 did=**segment 1 (scope items 1-6)** — `internal/tools/` now holds
  the host bridge + the C1/C3/C4 layer. `bridge.go` = `agent.ToolProvider` (the loop's existing
  Tools()/Execute() seam, `internal/agent/loop.go:563,696`; `var _ agent.ToolProvider = (*Bridge)(nil)`
  proves it), so no parallel dispatch path exists. Per call: C3 caps → one `risk.Assess` over
  C26-canonical paths + C25 taint bound per task scope → L0 pass / L1 `Gate.PendingWindow` /
  L2 `Gate.PendingApproval` / Deny refuse-without-asking → D38d ceiling (clamped ≤4, enforced in
  the bridge so a host-internal caller cannot bypass it) → C22 per-tool `context.WithTimeout` →
  C25 `Mark` on the result → complete `tool_call` row (assessed risk_level + decision + outcome +
  error_class + correlation_id). A rejection is NEVER a Go error (SPEC-07 §2 未声明即拒绝调用（不是报错，
  是拒绝）: `loop.go:654-657` turns a provider error into class internal, which is exactly the
  "Wisp broke" misreading C3 forbids) — the tests assert `err == nil` AND the tool's run counter
  stayed 0. `paths.go` closed a real gap: **no type in this repo implemented the frozen C19
  `PathCanonicalizer`/`SensitiveClassifier` seams before this**, so tickets 17/18 left R2/R3
  DORMANT in every composition; wiring them is what makes out-of-allowlist read → L2 a verdict
  (`TestDeclaredRiskIsOnlyAFloor`). `registry.go` declares all four C4 slots and answers
  `ErrSlotNotLanded` for manifest/goja/mcp (`TestC4SlotsWithoutAnImplementationAreNotSilent`)
  rather than pretending. Registered **fs.read + fs.list only** (`TestFSRegistrationIsTheL0Pair`
  pins the pair at 2 entries).
  AC#1 GREEN, test names: `TestUndeclaredCapabilityHardRejects`,
  `TestRiskDecisionRoutesEachLevelToItsBranch` (7 subtests: L0-passes-no-gate / L1-window /
  L1-timeout-executes / L1-veto / L2-approval / L2-timeout-auto-rejects / NoGate-fails-closed),
  `TestToolCallRowsAreComplete` (real SQLite, 5 rows: allow+success / L2 reject+user_rejected /
  unknown-tool / bad-args / L1 veto), plus `TestCapabilitySetIsFrozen` (11 tokens),
  `TestToolConcurrencyCeilingIsFour`, `TestPerToolTimeoutHonoredViaContext`, `TestC1ToolContract`,
  `TestC4RegistryRejectsBadDeclarations`, `TestFSReadTaintsIt`→`TestFSReadTaintFeedsR4`
  (R4 + SessionOverrideBlocked through the real C25 engine), `TestSensitiveFileIsDeniedNotEscalated`
  (R3 tier A never reaches a gate), `TestEmptyAllowlistAuthorizesNothing`.
  **Gates (scoped to ./internal/tools/, all run):** `gofmt -l` empty; `go vet` CLEAN;
  `go test -count=2` ok 0.480s; `go test -race -count=1` ok 1.394s; `tools/d22scan -root .`
  → 1 finding, `internal/llm/adaptertest/mockllm.go:68` bare-goroutine, **pre-existing at
  78b1466 (ticket 11), zero findings in internal/tools**; `internal/ball`+`cmd/balldebug`
  (live ticket-62 agent) untouched, frozen files zero-diff. Commits c9c3a6f + 64c5fea.
  **THE COMPOSITION ANSWER (rulings A8/A11/R12 shape — read this before believing "green"):**
  `who constructs the bridge?` — **nobody in production yet.** No non-test code calls `tools.New`
  or `agent.New`; the **only non-test importer of `internal/agent` in the entire repo is
  `internal/tools/bridge.go` itself** (re-run the判据:
  `grep -rn '"github.com/CarlosShao/wisp/internal/agent"' --include=*.go cmd/ internal/ | grep -v _test.go`),
  and `agent/loop.go:190` still defaults a nil provider to `EchoProvider`;
  `cmd/wisp` constructs no loop at all. **Ticket 12 wires it** (R12 retitled 12 the 接线票),
  which is why `internal/tools/doc.go` states this out loud instead of leaving an unwired-but-green
  module. AC#1 was ticked because its text is a *bridge-test* criterion and does not depend on the
  wiring; no other AC was touched, and none of the remaining five can be closed by segment 1 alone.
  next=**segment 2**: `fs.write` (D31 temp+atomic rename + `TestAtomicWriteKillsMidWrite`),
  `fs.trash` (shell API, L1), `fs.move` (L1 / cross-volume L2), `fs.delete` behind
  `[fs] delete_enabled`, applied-steps report on cancel, the artifacts-spill rule-separation test
  (AC#5/AC#6; `internal/agent/spill.go` is the host-internal precedent — the D22 ban-7 regex
  `"spill"`/`"internal.x"` forbids naming a gated tool after it, keep it that way), and the
  `[fs] allowed_dirs` first-use ask flow. **Two hazards segment 2/21 must not trip:**
  (a) `loop.decideRisk` (`loop.go:719-738`) rejects anything declaring L1/L2 BEFORE Execute runs,
  so `fs.write` will be dead on arrival under the loop until ticket 21 replaces that function with
  this bridge's verdict — register it only with the single-writer contract (bridge owns
  `tool_call`, loop runs with a nil Journal) or the row doubles;
  (b) the missing `tool_call.rules_hit` column is a SPEC-02 §3 contract change, not a segment-2
  task — rules_hit currently travels in `Decision`/`OnDecision` + the audit line only.
  Also still owed to ticket 19 (N-11): `rules_gateway.go`'s R4 comment says "Dormant until ticket 19"
  while the loop wiring is what landed here; that file is FROZEN for this agent, so it was NOT
  edited — ticket 21 owns that comment.
- [2026-09-20T12:10Z] agent=T20-seg2 did=**segment 2 (the fs write family)** — `fs.write`
  (D31 `stageAndRename`: the staging file is created in the TARGET's directory so the final
  `os.Rename` never crosses a volume, and the target is touched by exactly one syscall, so a death
  before it cannot produce a partial file), `fs.trash`, `fs.move` (L1 same-volume rename / L2
  cross-volume or occupied destination), `fs.delete` (registered ONLY when
  `FSDeps.DeleteEnabled`, re-checked inside Execute so a roster that outlived the setting still
  refuses). The R8 escalations reach the FROZEN rules through a new host-side `Decl.Facts` hook:
  the bridge keeps owning `Declared` and `Paths` (a tool can neither lower its own floor nor hide a
  target), the hook contributes the irreversibility classes, and a panicking hook contributes an
  unknown class which R8 fail-closes to L2. Every existence probe canonicalizes through C26 first.
  **trash is the real shell API, and how that is proven:** `SHFileOperationW` with
  `FOF_ALLOWUNDO` + `FOF_SILENT|NOCONFIRMATION|NOERRORUI|NOCANCELLATION`, issued on a thread
  pinned with `runtime.LockOSThread` + `CoInitializeEx` (the COM apartment is per OS thread;
  without it shell32 deadlocks — measured, see below), and success is then CONDITIONED on finding
  the `$I*` restore record shell32 wrote for that exact path under `<drive>\$Recycle.Bin\<SID>`
  (`TestShellTrashLeavesARecycleBinEntryTheOSCanEnumerate` re-derives it independently of the
  implementation's own check). An unlink cannot produce that record, and `os.Remove` appears in
  this package only in `fs.delete` and in the cross-volume move's source removal, both of which
  declare themselves permanent (R8 → L2). Two findings the next agent should not rediscover the
  hard way: (1) **`SHQueryRecycleBinW` blocks forever from a Go thread** — the first version used
  it for the item-count proof and the test binary sat in that syscall until the go-test timeout
  killed it, i.e. a tool call that can hang a task in silence; (2) **`FOF_WANTNUKEWARNING` asks
  shell32 for a MODAL dialog**, which is the same hang wearing a face on the user's desktop. Both
  dropped in favour of the restore-record check. The trash tests DO place one small temp item in
  the real Recycle Bin each run; that is the price of an honest assertion.
  **Two D31 bugs the composition exposed (both in the fake-clean-cancel direction, both fixed
  here):** `approval.Gate.PendingApproval` never recorded a handoff, so an approved L2 write that
  then stopped mid-flight was rendered 「本次调用未进入执行阶段」 over a list of steps that HAD
  landed; and `CancellationReport.TextFor` had no branch for "steps landed, no veto on record"
  (task cancellation / tool fault), so it fell through to the 「未上报其步骤」 sentence. The
  `markStarted` key is the INCOMING correlation id, because the queue re-stamps
  `d.CorrelationID` with its own item id — the card's address, not the call's. The only
  applied-steps shape in the repo is still `approval.CancellationReport`: `tools.CancelBus.Report`
  returns an `fmt.Stringer` precisely so no second shape could be born on the way.
  **AC#2/3/4 GREEN (test names in the AC text). AC#5 and AC#6 NOT ticked:** neither is in this
  segment's scope — AC#5's junction-denial coverage lives in `internal/risk`
  (`pathresolver_junction_windows_test.go`), and AC#6's artifacts-spill rule-separation test is
  still unwritten.
  **THE R12 ANSWER, restated because three ACs this week were left open on exactly this question:
  who calls this in production? NOBODY YET, and it is NOT reachable from `cmd/wisp`.** Re-run the
  判据: `grep -rln '"github.com/CarlosShao/wisp/internal/tools"' --include=*.go cmd/ internal/ |
  grep -v _test.go` → only the four `internal/agent/approval` files;
  `grep -rn '"github.com/CarlosShao/wisp/internal/agent"' --include=*.go cmd/ internal/ | grep -v
  _test.go` → only `internal/tools/bridge.go`; and
  `grep -rn "tools.New\|approval.New\|agent.New(" --include=*.go cmd/` → **zero hits**. Ticket 12
  owns that wiring; no AC was ticked on the strength of a call path that does not exist.
  **`decideRisk` WAS NOT TOUCHED** (`git diff HEAD -- internal/agent/loop.go` empty; the guard sits
  at `internal/agent/loop.go:719-738`). Deliberate, per the hazard register: it refuses any call
  whose roster entry declares L1/L2 before Execute runs, so an `fs.write` dispatched by a Loop is
  still refused rather than gated. Replacing it means handing the loop the bridge's verdict, and
  that replacement cannot be composed honestly inside this segment: `cmd/wisp` builds no loop and
  no gate, so the only test proving "guard gone, gate in its place" would have to author ticket
  12's production wiring. What IS proven here is that nothing else is missing —
  `internal/tools/wiring_test.go` runs the real `approval.Gate` over this bridge: an L1 write
  blocks its full 2s window and then writes, an in-window veto writes nothing, a late veto renders
  approval's own report wording, and an L0 read opens no card.
  **Gates (scoped to ./internal/tools/ + ./internal/agent/..., all run):** `gofmt -l` clean on
  every file this agent touched (`internal/agent/loop_golden_test.go` and `truncation_test.go`
  carry pre-existing drift from d6f6406, untouched here); `go vet` CLEAN;
  `go test -count=2` ok tools 14.386s / agent 2.463s / approval 0.145s;
  `go test -race -count=1` ok 9.982s / 2.801s / 1.131s; `tools/d22scan` → **clean** (segment 1's
  pre-existing mockllm finding is gone). `internal/ball` + `cmd/balldebug` untouched, no GUI
  process launched, no `-tags winlive` run, frozen files (`docs/PLAN.md`, `docs/specs/*`,
  `internal/risk/*`) zero-diff. Commits 0986d63 + 94827e4 + e669567.
  **MUTATION EVIDENCE (so AC#3/AC#4 are not read as self-certifying):** with
  `os.CreateTemp` swapped for a direct `os.OpenFile(target, O_CREATE|O_WRONLY|O_TRUNC)` the two
  kill tests go red (`D31 violated: after a mid-write kill the target holds N bytes…` and
  `AC#3 violated: a partial file exists at the target`); with `fofALLOWUNDO` zeroed out (a
  permanent shell delete, i.e. the unlink-shaped lie) both trash tests go red with
  `回收站中没有出现任何新的还原记录（C:\$Recycle.Bin）：拒绝按已回收处理`. Both mutations were reverted;
  `git status` is clean and every gate above was re-run against the committed tree.
  next=**ticket 12 (composition) or ticket 21 segment 2**: (a) build the Bridge + `approval.Gate`
  in `cmd/wisp`, and in the SAME commit replace `decideRisk`'s roster-floor refusal with the
  bridge's verdict while running the loop with a nil Journal so `tool_call` rows do not double;
  (b) AC#5 bridge-level junction/reparse denial test; (c) AC#6 artifacts-spill rule-separation
  test (`internal/agent/spill.go` is the host-internal precedent; D22 ban 7 forbids naming a gated
  tool after it); (d) the `[fs] allowed_dirs` first-use ask flow (minimal native prompt now, full
  GUI at 39) — note it has NO SPEC-12 §5 row, so marking it `DEFERRED` in code would break the
  1:1 cross-check; (e) still owed to ticket 19 (N-11) the `rules_gateway.go` R4 comment, and the
  missing `tool_call.rules_hit` column (a SPEC-02 §3 contract change, human approval required).
- [2026-09-21T01:20Z] agent=T20-box5 did=**第 5 框：桥层的真 junction / 8.3 拒绝**（票 18 integration 缺的
  **上层半边**）。**先复现编排者那条 grep**：`junction|mklink|GetShortPathName` 在 `internal/tools/*_test.go`
  **0 命中** ⇒ 当时只有"下层会拒"、没有"上层拦得住"，正是 A33② 那一族的形状。
  **做了什么**：新文件 `internal/tools/bridge_junction_windows_test.go`（`//go:build windows`，与票 18 同惯例，
  **零 `t.Skip`**），6 条用例全部经 `Bridge.Execute` 而不是直接调工具：真 `mklink /J` junction 种在**授权根内**、
  目标在根外；真 8.3 短名取自 `GetShortPathNameW`（两种产物**本机都建成了**，前置条件写进 `t.Fatalf` 的失败文本里，
  跑不动就红、不许用字符串假装）。审批门一律 **AnswerAllow** ⇒ "人已点批准"之后仍然拒；目标目录按**逐文件
  sha256 快照**断言 before==after，两侧目录再断言无 `.wisp-tmp-*`；每个 fixture 先做**阳性对照**（`os.ReadFile`
  从 junction 对面读出那 94 字节，读不到就判用例无效）。覆盖：`fs.read`/`fs.list`/`fs.write`(覆写+新建)/
  `fs.trash`/`fs.move`/`fs.delete` 六个入口 + A 档短名 + 真 SQLite 审计行。
  **证据在哪**：`docs/evidence/s1/20-bridge-junction-shortname.md`（§2 用例表 / §3 变异检验 / §4 同形排查清单）。
  **变异检验**（防自证）：把 `fs.go:95-98` 与 `fs_write.go:169-174` 改成"吞掉解析器的不、直接用原样路径"
  → **6 顶层 + 8 子项全红**，红字里 junction 对面的文件被真写、真进回收站（还原记录 `$I41EUHE.txt`）、真被永久删；
  已还原（`grep -c MUTATION-20` = 0、`git diff internal/tools/fs*.go` 空、复跑全绿），窗口只活了几十秒且开工前
  `tasklist` 确认无 `go.exe`/`wisp.exe`/`balldebug.exe`。**8.3 那两条在这条变异下没红**，诚实记着：这条变异吞的是
  工具层的第二次解析，而 8.3 的拒绝发生在判定层（R2/R3 自己 Canonicalize），不是同一条路径。
  **门禁**（两个新文件都落地后复跑，2026-09-21 09:4x）：`gofmt -l internal/tools` 空、`go vet ./internal/tools/` rc=0、
  `go test -count=2 ./internal/tools/` `ok 25.192s`、`go test -race -count=1 ./internal/tools/` `ok 17.709s`、
  `bash scripts/d22scan.sh` clean；本框 7 条新用例 `-count=2 -v` 实跑
  **`=== RUN` 36 / `--- PASS` 36 / `--- SKIP` 0 / `--- FAIL` 0**（全文连字符串 `skip` 都 0 次）。
  **同形排查（13 个入口，逐项文件:行在证据 §4）**：fs 全套 11 个入口 + 两个 R8 探针**都**过 C26，
  且判定层与执行层各解析一次；两条 artifacts 路（`internal/agent/spill.go:97-106`、
  `internal/memory/artifacts.go:104-146`）**不过** C26，但它们不接受调用方给的路径（`artifactName` 把 callID
  裁成 `[A-Za-z0-9_-]`、`validArtifactName` 只收裸名）⇒ 按硬规矩**只上报、一行守卫都没加**。
  **三条留给 owner 的判定**（本框未处理）：① `<allowed>\jn\<A 档文件>` 今天是**可批准的 L2**、不是 Deny，
  卡面上目标只写成"无法规范化"——**批准它的人看不见自己批准了什么**，挡住字节的仅是工具层的第二次解析；
  ②被批准的红队调用在 `tool_call` 里记成 `decision=allow / outcome=error / error_class=tool`，
  即"工具说不"而不是"策略说不"（与 A21 同族）；③第 7 框 (a) 要的"机器可辨原因值"今天**不存在**——
  两类拒绝只在措辞上分岔，`risk.ErrReparseDenied` 在 `bridge.go:628-631`/`fs.go:96` 被折成字符串。
  **没证明但看起来成立的**：symlink 未测（要特权或开发者模式，我只造了 junction，虽然二者落进
  `pathresolver_windows.go:87` 同一个 reparse 分支）；UNC 与 `\?\` 在桥层没重做（票 18 覆盖的是 resolver，
  且这两种拼法规范化会成功，形状等同本框的 8.3 越界项）；**硬链接根本不是 reparse point，两层都看不见它**——
  本框未测，也不知道有没有人测过，若算缺口应新开条目；跨卷 junction 未测。
  next=**第 6 框**（artifacts spill 规则分离）与**第 7 框**（`allowed_dirs` 首次使用询问，前置=先给上面③定层次）；
  A18 的进展见下一条 log。
- [2026-09-21T01:35Z] agent=T20-box5 did=**A18 特征化（不修行为，只把今日真实行为钉成断言）**。
  **做了什么**：新文件 `internal/tools/bridge_a18_kill_windows_test.go`（`//go:build windows`，零 `t.Skip`）。
  父测试用 `exec.Command(os.Args[0], "-test.run=^TestA18…$")` 重启**自己的测试二进制**（本仓 `internal/` 零
  `TestMain`，所以走子分支），子进程经**真桥 + AnswerAllow 门**调 `fs.write`、`WriteChunk:4`，
  在 `Hooks.AtStep("write:8")` 落信号文件后**永久阻塞**；父进程轮询到信号后跑**真 `taskkill /F /PID`**
  （A18 判据原文要的就是这个，没退化成 `Process.Kill`，也没退化回进程内 `Hooks.Kill`——后者正是 A18 说的假象来源）。
  **今日真实结果**：①覆写分支目标**逐字仍是 OLD-BYTES**、新建分支目标**不存在** ⇒ D31 在真故障下成立
  （机制就是 `fs_write.go:277` 暂存建在目标目录 + `:325` 单次 `os.Rename`）；②**每次 kill 恰好留 1 个
  `.wisp-tmp-*`**（实测各 4 字节 = 被杀时已写进暂存区的量，`t.Logf` 记数不钉数）；③**没有任何东西扫它**：
  同目录再来一次全新桥的成功写盘之后两个残留仍在原地（全仓 `tempPrefix` 只有 `fs.go`/`fs_write.go` 三处写入者、
  零读取者 ⇒ 今天根本没有可称"下次启动"的清扫器）。
  **判据①闭合、判据②不闭合**：A18 原文两条完成判据里"目标要么完整要么不存在"这条现在有真 kill 证据了；
  "下次启动自愈清理 **或** 明确申报可接受残留并写进 SPEC"这条**我没动**——把"没扫"改成"扫掉"是行为设计，
  写进 SPEC 是契约变更（D22），都归 owner。用例因此**故意钉在"仍在"**：谁加了清扫器它就红，逼他回来改这条
  判据而不是静默改行为。**next=owner 在 A18 判据②的两个选项里选一个**（加启动清扫器=哪张票；或 SPEC 里
  承认残留可接受）；证据 `docs/evidence/s1/20-a18-taskkill-residue-characterization.md` §2/§3/§5。
  反空跑：绕过阻塞点的子进程自己 `os.Exit(5)` 打印 "wrote without blocking" 并被父进程判失败；
  未测的分支（暂存文件**内容**是否等于已 Write 量、rename 之后那一段的真 kill 时机）写在证据 §5。
