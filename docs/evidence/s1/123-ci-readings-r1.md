# 123-ci-readings-r1 — 票 123 的 CI 真读数与本机复核（编排者，2026-11-20）

## 1. CI runner 上的 `cmd/wisp CLI tests` 步骤读数（AC#4）
- **run id** `44516745051`（`dev ci`，push，2026-11-19T11:39Z 触发，head `329b8e44`）
- **job id** `135534604009`（`test-windows`）
- **step** `cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)` —— **该步结论 = success**
- 该步自己的四数行（日志逐字，1.6 万行日志落 `D:\tmp\wisp\20261120\lead\ci-job-135534604009.log`）：
  `portable-tests.sh: four numbers (all from -v output): === RUN=1016  --- PASS=767  --- FAIL=0  --- SKIP=20`
- **FAIL=0** ⇒ 票面要的"4 枚不再红"在 CI 上成立。**注：票面 AC#4 写的是 `PASS=33 / FAIL=0`（旧步骤形状）；今天该步打印的是 `RUN=1016 / PASS=767 / FAIL=0 / SKIP=20`**（同一步、同一载具 `scripts/wisp-cli-staged-tests.sh`）。以本次真读数为准，PASS 数字变了但 FAIL=0 这条实质不变。
- ⚠ **2026-11-20 12:18Z 之后所有 `ci` run 都是 `queued`（≥6 枚，最早一枚排队 49 分钟）** ⇒ 今天拿不到更新的 CI 读数，只能用最后一次 completed 的 run。这是环境事实，不是我没跑。

## 2. 本机整包（我亲自跑，载具）
- `scripts/wisp-cli-staged-tests.sh`：`RUN=1016 PASS=1003 FAIL=1 SKIP=13`（日志 `D:\tmp\wisp\20261120\lead\123-r4-full-count1.log`）
- 唯一 FAIL ＝ `TestAC4FocusReturnToPriorWindowGap33r5`（已知前台锁环境红：`panel_host_windows_test.go:1120-1124` 「panel did not take the foreground」），与本票无关。
- 13 枚 SKIP 具名：`TestAC6CSPOnTheNotBuiltNoticeIsReadFromItsRealResponse33lane2`／`TestAC247LiveMicrophoneLevelsReachTheBallSeam`／`Test41a2r1DrainTimeoutVerdictOptIn`／`Test41a2r1DrainTimeoutArmChild`／`Test26PEchoNeverTriggersAndNeverTranscribes`／`Test26PWithoutAECTheEchoTranscribesVerbatim`／`Test26PMarginRetestAndGateReverse`／`Test26PBargeinStopsWithin400ms`／`Test26PFeedOnlineMatchesWindowReplay`／`Test392KeepalivePostureFindingsAndColdDefault`／`Test392AecTuningDoubleTalkSwing`／`Test392DisposeReleasesTheEngines`／`Test16r1PremortemMicPermissionDeniedLiveProbe`。

## 3. AC#5 的那条 skip —— 已选"永不 skip 的替身"
- `internal/risk/syncdirs_redteam_windows_test.go` 文件头注释逐字：「Ticket 123 AC#5: this leg used to skip twice on a machine with no live sync root … **The skip is gone.** When no live root exists the case runs the same bypass-spelling battery against an injected registry-grade synthetic root … the choice of stand-in over an exemption entry is ticket 123 AC#5's recorded decision.」
- 我现跑：`go test ./internal/risk -run '^TestSyncRedTeamRealOneDrive$' -v` ⇒ `=== RUN`、`--- PASS: TestSyncRedTeamRealOneDrive (0.01s)`、`ok internal/risk 0.052s`（本机有 live root：`C:\Users\swq\OneDrive`，走 live 分支；无 root 的机器走 synthetic-root 替身分支）。

## 4. AC#2／AC#3 的代码事实
- 注入面（测试内、非生产）：`cmd/wisp/ticket123_replyleg_test.go` / `_windows_test.go` / `_other_test.go`（`446831ea9`），`git merge-base --is-ancestor 446831ea9 329b8e44` ＝ **true** ⇒ 上面的 CI 读数已含该修法。
- 拒绝形在跑：`run_mode101_test.go:329` `t123script(t, "yes t101-corr-1", "no t101-corr-2")`、`:467` `no t101-corr-1`、`:541` `no t101-corr-1", "yes t101-corr-2"` ⇒ 用例在 Confirm 返回"拒绝"时仍通过（不是"永远批准"的恒真）。
- 生产真确认腿在：`cmd/wisp/main.go:225` `reply := interactiveStdin()`；装配 `cmd/wisp/run.go:1568-1569` `if s.reply != nil { rt.attachReplyListener(s.reply, s.replyVeto) }`。
