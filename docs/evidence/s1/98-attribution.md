# 98 — ubuntu 腿那批红的逐枚归因（只读普查）

**票**：`.scratch/wisp/issues/98-cmd-wisp-tests-never-run-on-this-host.md` §09-28 追加 (b)「ubuntu 腿那 19 枚未解释的红」。
**Status**：证据件（只读）。**代码/票面/脚本一字未改，未 `git add` / `git commit`**；本文件是全仓唯一写入。
**锚定树**：`f4ad34ee1`（派单给的 HEAD）。取数用的是 `git archive` 的纯净快照，实际归档的提交是 `8a75004ce`
（普查期间 HEAD 前进了一枚）。`git diff --stat f4ad34ee1 8a75004ce -- cmd internal` 只多两枚**别人的** `*_test.go`
（`internal/memory/l1_warm_roster_29r2_test.go`、`internal/tools/spawn_roster_census_229_test.go`）⇒ **`cmd/wisp` 与
`f4ad34ee1` 逐字节相同**。快照 `sha256=3fba786b92ac3859c547c70a45fbf6668dcc327a628f159ba0d7f0a55326fe9b`（329,308,160 字节）。
**时间**：2026-11-05 11:44–12:2x +08。
**只读承诺**：跑测试只写 `t.TempDir()`/容器内 `/tmp`；快照一律 `git archive`（不建 worktree、不 checkout）；工作树里别人在飞的
`M cmd/wisp/main.go`、`?? cmd/wisp/commands*.go`、`?? cmd/wisp/zz_probe_temp_test.go.txt` **一字未碰**。

## 0. 一句话结论

今天这棵树上 ubuntu 腿**不是 19 枚，是 81 枚**（`PASS=143 FAIL=81 SKIP=0`，顶层用例 224 枚）。
**80 枚能点到一枚具体的票**；**1 枚不是红**（我自己的仪器没 `.git` 造出来的假红，带 `.git` 复跑 PASS）；
**unknown = 0**。81 枚里 **76 枚是同一个根因**：`cmd/wisp` 的测试夹具写 `api_key_ref = "dpapi:…"` 并真的
`secret.Store()`，而 `internal/secret/protect_other.go` 在非 Windows 上是**documented fail-closed 桩**（C28 / 票 06）。

⚠ 两个诚实的边界：① 76 枚是**按"第一行红"**点的名——这些用例跑到夹具那行就 `t.Fatal` 了，**它自己后面的判据还会不会红，本件量不到**
（这不是归因上的 unknown，是"更深处的红未求值"）；② 81 这个数**不是票 98 表头那个 19 的替代品**，是同一个病的**分母长大了**（见 §1 末）。

## 1. 读数（原文）

仪器：`docker.io/library/golang:1.27-bookworm`；容器内 `go env GOOS GOARCH CGO_ENABLED CC GOVERSION`
= `linux / amd64 / 1 / gcc / go1.27.3`。

```
$ go vet ./cmd/wisp/                    # linux, CGO_ENABLED=1
vet_rc=0

$ CGO_ENABLED=1 go test -count=1 -v -timeout 90m ./cmd/wisp/
test_rc=1
fail_lines=81   pass_lines=143   skip_lines=0   run_lines=302
FAIL    github.com/CarlosShao/wisp/cmd/wisp    66.778s
```

- 顶层用例 = `PASS 143 + FAIL 81 = 224`；`=== RUN 302` 是**含子用例**的行数。
- **交叉核对（两条互不相干的路）**：同一发里 `TestAC1AC2DispatchHopGate133` 自己用 `-test.list` 读的名册逐字写着
  "lists **224** startable cases"，与 `143+81=224` 同一个数。
- `SKIP=0`：这一发没有任何 SKIP 替红打掩护（`tools/d22scan/runtests.sh` 把 SKIP 判红）。
- **第二发（写件时的 HEAD `a980734f0`，同一仪器、同一命令）**：`fail_lines=81 pass_lines=143 skip_lines=0 run_lines=302`，
  `FAIL … 62.440s`；两份日志的 `--- FAIL` 与 `--- PASS` 名字清单 `diff` **均为空** ⇒ 这 81 枚在 `a980734f0` 上**逐枚同形**
  （该提交只动 `internal/agent/approval/{gate,ui}.go`、一枚新测试与一份证据件，没有挪动任何一条红的归属）。
- **与表头那三个历史读数的关系**：`19 of 29`（`scripts/wisp-cli-tests.sh:25-27`，更老的树）→ `20 FAIL`
  （`docs/evidence/s1/131-followup-0-crossplatform.md:145`，`3029415`）→ **`81`（本件，`f4ad34ee1`）**。
  变的是分母（29 → 40 → 224），**不是新病**：三发里占绝对多数的都是同一句
  `secret: DPAPI is only available on Windows; use an env: ref on this platform`。

**仪器坑（写给下一位，省 40 分钟）**：第一次拿**只读 bind-mount** 的仓直接跑，默认 `-timeout 10m` 在
`TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` 上 `panic: test timed out after 10m0s`。
机理量到了：`internal/llm/adaptertest/mockllm.go:65` 每发都 `go build tools/mockllm`，在那个 9p/virtiofs 挂载上单次建了
**281.40s**（原文 `--- FAIL: TestAC1AlwaysBranchDoesNotRevertAHandEditedKey (281.40s)`），不是测试挂死。
把 `git archive` 快照**解进容器本地盘**再跑，同一包 **66.8s** 跑完。⇒ 本件全部读数用后者。

## 2. 六类根因（81 枚全覆盖，枚数相加 = 81）

| 类 | 枚 | 机制（可核出处） | 归谁 |
|---|---|---|---|
| **C1** 非 Windows 上 DPAPI 桩 | **76** | `internal/secret/protect_other.go:18` 的 `errDPAPIUnavailable` = `secret: DPAPI is only available on Windows; use an env: ref on this platform`。测试夹具写 `api_key_ref = "dpapi:acme"` 后真的 `secret.Store(...)` / `wisp secret set`，于是**夹具那一行**就 `t.Fatal`。桩是**故意 fail-closed**（同文件 `:7-16` 注释；C28 / 票 06 是 Windows-only 设计），不是 bug。 | **该用例所在的票**（它的选择：把夹具换 `env:` ref，或按票 70 / 票 131 的先例把平台写进文件名/build tag）。**平台规则本身**归票 06 / 票 08，不归任何"修红"的票。 |
| **C2** 名册漂移（**两条腿都红**） | **1** | `TestTicket255RosterStillMatchesTheActualReadSites`：`cmd/wisp/config_readers_255.go` 的名册说 `Cost` 无读者、`LLM`/`Voice` 只有旧读者，而 `internal/agent/costbook_44.go` + `internal/agent/quota_44r4.go`（票 44 腿 `a11e63bee`）真读 `cfg.Cost`/`cfg.LLM`、`internal/speech/cloudcatalog.go`（票 61 腿 `742888795`）真读 `cfg.LLM`/`cfg.Voice` ⇒ 三节 `match=false`（`config_receipt_255_test.go:299`）。**Windows 干净树上逐字同一条**（§4）。 | **票 255**（名册主人，`config_readers_255.go` 要重判）＋触发者**票 44**、**票 61**。 |
| **C3** 门把"源码里有、本平台没编进来"读成"没有" | **1** | `TestAC1AC2DispatchHopGate133`：`leg_dispatch_gate_133_test.go:251` —— 账本把 `slo` 腿记成由 `Test384ResidentWatchdogEndToEnd` 覆盖，而该用例住在 `resident_watchdog_windows_test.go:184`（windows-only）⇒ linux 二进制名册里没有它 ⇒ 红。**Windows 上 PASS**（§4）。与票 131 已修过的同形（131 的做法是"具名披露、不判红"）。 | **票 133**（它自己已有 `runRosterDisclosure135` 披露机制，只是这一支仍 `t.Errorf`）。 |
| **C4** 故意的红（有据） | **1** | `TestAC4EveryLegIsNailedOrRuled`：`leg_sink_gate_131_test.go:402-404` 逐字 "This reading stays red on purpose: a gate that cannot see is not a gate that has seen nothing to complain about." ⇒ linux 上"零钉"必须红。**Windows 上 PASS**（§4）。 | **票 131**（裁定与读数出处 `docs/evidence/s1/131-followup-0-crossplatform.md` §4）。 |
| **C5** **不是红**（我的仪器） | **1** | `TestCleanCheckoutBuilds_AC11`：`panel_host_gate_test.go:404` `git rev-parse --show-toplevel: exit status 128` —— `git archive` 快照里**没有 `.git`**。同一测试在容器里挂**真仓（带 `.git`）**跑：`--- PASS: TestCleanCheckoutBuilds_AC11 (44.18s)`，`ok … 44.278s`。 | **无主**（不计红；不加 tag、不改测试）。 |
| **C6** POSIX 上拒绝的哨兵不同名 | **1** | `TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119`（`secret_dataroot_119b_test.go`，文件头 `//go:build !windows` ⇒ **只在 POSIX 跑**）：拒绝发生了，但 `errors.Is(err, winsec.ErrUnresolvedPath)` 为假——实际是 `risk: path traverses a reparse point …`（`internal/risk/pathresolver.go:36` `ErrReparseDenied`）⇒ `:242` 红。 | **票 119**（该用例主人；两枚哨兵在 POSIX 上要不要统一是它的裁，`internal/winsec` 那边的 POSIX 半边是票 94 的地界）。 |

## 3. 归因账：attributable 80 / 不计红 1 / unknown 0

| 桶 | 枚数 |
|---|---|
| 能点到一枚具体的票（C1 76 + C2 1 + C3 1 + C4 1 + C6 1） | **80** |
| 不计红（C5：仪器假红，带 `.git` 复跑 PASS） | **1** |
| **unknown（说不出归谁）** | **0** |

按票点名（票号 = 用例**所在票**；全名见 §5 表右列；合计 81）：

| 票 | 枚 | 票 | 枚 | 票 | 枚 |
|---|---|---|---|---|---|
| 201 | 9 | 63 | 8 | 223 | 8 |
| 12 | 6 | 101 | 5 | 197 | 5 |
| 255 | 5（含 C2 1） | 233 | 3 | 35 | 3 |
| 224 | 3 | 315 | 3 | 248 | 2 |
| 200 | 2 | 165 | 2 | 370 | 2 |
| 174 | 2 | 151 | 2 | 365 | 2 |
| 226 | 1 | 227 | 1 | 33 | 1（C5，不计） |
| 133 | 1（C3） | 131 | 1（C4） | 220 | 1 |
| 329 | 1 | 119 | 1（C6） | 378 | 1 |

## 4. Windows 侧对照（我实测的，不是推理）

本机 `export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"` 之后：

| 用例 | 侧 | 读数 |
|---|---|---|
| `TestRunTextTaskTextPathEndToEnd` / `TestSecretSetGetListRoundTrip` / `TestTicket223RunArmsTheReloadTick` / `TestAC2CredentialSentinelAppearsInNoArtifact` / `Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands` / `Test378SpawnedChildToolCallsHaveAParentRow` / `TestZZ315A1RootLoopFirstBridgeCall` | windows | **7 枚 C1 抽样全 `--- PASS`**（⇒ C1 是平台，不是包坏） |
| `TestAC1AC2DispatchHopGate133`（C3） | windows | `--- PASS (0.70s)` |
| `TestAC4EveryLegIsNailedOrRuled`（C4） | windows | `--- PASS (0.49s)` |
| `TestTicket255RosterStillMatchesTheActualReadSites`（C2） | **windows，干净 `f4ad34ee1` 快照** | **`--- FAIL (1.22s)`**，逐字同三条 `match=false`（`cost`/`llm`/`voice`） |
| `TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119`（C6） | windows | **不存在**（文件 `//go:build !windows`），故 windows 上无法对照 |

**由 C2 推出的、票 98 表头没写的一格**：这一枚没有 build tag，也在 `scripts/portable-tests.sh --scope=cli` 的 scope 里，
所以在 `f4ad34ee1` 这棵树上 **`scripts/wisp-cli-tests.sh`（CI 的 `test-windows` 那一步）不可能绿**——
除非同批把 `config_readers_255.go` 重判。它不是 ubuntu 的病，是两条腿共有的病，只是被这次普查撞出来。

## 5. 81 枚逐枚点名

「红在哪一行」＝该用例**第一条**被判红的输出行（`file:line`，路径相对仓根）。类见 §2；归因票列 = 该用例**所在票**的号 ＋
它在 `.scratch/wisp/issues/` 的文件名主干（已交件的 `-done` 后缀省略）。

| # | 用例 | 红在哪一行 | 类 | 归因票（号 ＋ 票名） |
|---|---|---|---|---|

| 1 | `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey` | `cmd/wisp/always_write_no_clobber_226_test.go:40` | C1 | 226 always-answer-rewrites-the-whole-config-snapshot-and-hides-hand-edits |
| 2 | `TestAC6AHandEditIsStillSeenUnderTheRealPoll` | `cmd/wisp/always_write_no_clobber_227_test.go:126` | C1 | 227 guarded-write-delete-shapes-and-migration-on-read-are-unasserted |
| 3 | `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` | `cmd/wisp/approval_always_201_test.go:122` | C1 | 201 nobody-can-answer-an-approval-l1-runs-by-itself-l2-dies-on-timeout |
| 4 | `TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused` | `cmd/wisp/approval_always_201_test.go:164` | C1 | 201 nobody-can-answer-an-approval-l1-runs-by-itself-l2-dies-on-timeout |
| 5 | `TestReplyListenerAllowsAnL2CardFromTheNativeSide` | `cmd/wisp/approval_reply_201_test.go:202` | C1 | 201 nobody-can-answer-an-approval-l1-runs-by-itself-l2-dies-on-timeout |
| 6 | `TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel` | `cmd/wisp/approval_reply_201_test.go:272` | C1 | 201 nobody-can-answer-an-approval-l1-runs-by-itself-l2-dies-on-timeout |
| 7 | `TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject` | `cmd/wisp/approval_reply_201_test.go:350` | C1 | 201 nobody-can-answer-an-approval-l1-runs-by-itself-l2-dies-on-timeout |
| 8 | `TestUnansweredL2CardTimesOutIntoRejectNeverExecution` | `cmd/wisp/approval_reply_201_test.go:439` | C1 | 201 nobody-can-answer-an-approval-l1-runs-by-itself-l2-dies-on-timeout |
| 9 | `TestL1VetoNeedsAChannelTheHostReallyWired` | `cmd/wisp/approval_reply_201_test.go:521` | C1 | 201 nobody-can-answer-an-approval-l1-runs-by-itself-l2-dies-on-timeout |
| 10 | `TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole` | `cmd/wisp/approval_seam_201_test.go:50` | C1 | 201 nobody-can-answer-an-approval-l1-runs-by-itself-l2-dies-on-timeout |
| 11 | `TestNativeHostSeamRefusesAPanelSourcedAllow` | `cmd/wisp/approval_seam_201_test.go:138` | C1 | 201 nobody-can-answer-an-approval-l1-runs-by-itself-l2-dies-on-timeout |
| 12 | `TestTicket255RosterStillMatchesTheActualReadSites` | `cmd/wisp/config_receipt_255_test.go:299` | C2 | 255 config-says-these-sections-took-effect-immediately-panel-for-a-width-nobody-reads-while-the-panel-host-never-receives-config-and-tier-has-zero-readers |
| 13 | `TestTicket255ReceiptOmitsPanelFromTheImmediateSentence` | `cmd/wisp/config_receipt_255_test.go:625` | C1 | 255 config-says-these-sections-took-effect-immediately-panel-for-a-width-nobody-reads-while-the-panel-host-never-receives-config-and-tier-has-zero-readers |
| 14 | `TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere` | `cmd/wisp/config_receipt_255_test.go:702` | C1 | 255 config-says-these-sections-took-effect-immediately-panel-for-a-width-nobody-reads-while-the-panel-host-never-receives-config-and-tier-has-zero-readers |
| 15 | `TestTicket255ReceiptStillNamesTheLiveReadSection` | `cmd/wisp/config_receipt_255_test.go:741` | C1 | 255 config-says-these-sections-took-effect-immediately-panel-for-a-width-nobody-reads-while-the-panel-host-never-receives-config-and-tier-has-zero-readers |
| 16 | `TestTicket255ReceiptSentenceAssemblyIsFiltered` | `cmd/wisp/config_receipt_255_test.go:801` | C1 | 255 config-says-these-sections-took-effect-immediately-panel-for-a-width-nobody-reads-while-the-panel-host-never-receives-config-and-tier-has-zero-readers |
| 17 | `TestTicket223RunArmsTheReloadTick` | `cmd/wisp/config_reload_223_test.go:371` | C1 | 223 checkandreload-has-zero-production-callers-hot-reload-never-runs |
| 18 | `TestTicket223HandEditedFsLooseningCostsAnL2Card` | `cmd/wisp/config_reload_223_test.go:424` | C1 | 223 checkandreload-has-zero-production-callers-hot-reload-never-runs |
| 19 | `TestTicket223RefusedLooseningKeepsOldValues` | `cmd/wisp/config_reload_223_test.go:497` | C1 | 223 checkandreload-has-zero-production-callers-hot-reload-never-runs |
| 20 | `TestTicket223TighteningRaisesNoCard` | `cmd/wisp/config_reload_223_test.go:537` | C1 | 223 checkandreload-has-zero-production-callers-hot-reload-never-runs |
| 21 | `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` | `cmd/wisp/config_reload_223_test.go:569` | C1 | 223 checkandreload-has-zero-production-callers-hot-reload-never-runs |
| 22 | `TestTicket223RestartTierSaysItWillNotApply` | `cmd/wisp/config_reload_223_test.go:643` | C1 | 223 checkandreload-has-zero-production-callers-hot-reload-never-runs |
| 23 | `TestTicket223FailureSentencesAreDistinct` | `cmd/wisp/config_reload_223_test.go:831` | C1 | 223 checkandreload-has-zero-production-callers-hot-reload-never-runs |
| 24 | `TestTicket232R1ASeededStartupBannerCannotCarryTheRestartTier` | `cmd/wisp/config_reload_223_test.go:937` | C1 | 223 checkandreload-has-zero-production-callers-hot-reload-never-runs |
| 25 | `Test233ReloadTierEventReachesAProductionConsumer` | `cmd/wisp/config_tiers_233_test.go:60` | C1 | 233 d36-three-tier-enum-is-tagged-on-no-key-and-the-reload-tier-has-neither-producer-nor-consumer |
| 26 | `Test233HotTierAgentReachesTheNextTask` | `cmd/wisp/config_tiers_233_test.go:120` | C1 | 233 d36-three-tier-enum-is-tagged-on-no-key-and-the-reload-tier-has-neither-producer-nor-consumer |
| 27 | `Test233ReloadTierReceiptOnlyForReloadTier` | `cmd/wisp/config_tiers_233_test.go:209` | C1 | 233 d36-three-tier-enum-is-tagged-on-no-key-and-the-reload-tier-has-neither-producer-nor-consumer |
| 28 | `TestRunPacketCarriesTheLoadedInstructionFiles` | `cmd/wisp/instructions_200r2_test.go:112` | C1 | 200 nothing-reads-the-projects-instructions-for-ai-although-all-five-harnesses-do |
| 29 | `TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound` | `cmd/wisp/instructions_200r2_test.go:189` | C1 | 200 nothing-reads-the-projects-instructions-for-ai-although-all-five-harnesses-do |
| 30 | `TestAC1AC2DispatchHopGate133` | `cmd/wisp/leg_dispatch_gate_133_test.go:251` | C3 | 133 no-instrument-sees-the-dispatch-hop-deleting-one-caller-keeps-76-tests-green |
| 31 | `TestAC4EveryLegIsNailedOrRuled` | `cmd/wisp/leg_sink_gate_131_test.go:280` | C4 | 131 the-models-leg-has-no-nail-deleting-9-lines-keeps-76-cmd-wisp-tests-green |
| 32 | `TestAC2CredentialSentinelAppearsInNoArtifact` | `cmd/wisp/panel_config_248_test.go:164` | C1 | 248 the-panel-has-no-settings-route-the-inbound-whitelist-has-zero-config-or-credential-methods |
| 33 | `TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef` | `cmd/wisp/panel_config_248_test.go:272` | C1 | 248 the-panel-has-no-settings-route-the-inbound-whitelist-has-zero-config-or-credential-methods |
| 34 | `TestCleanCheckoutBuilds_AC11` | `cmd/wisp/panel_host_gate_test.go:404` | C5 | 33 panel-host-c27 |
| 35 | `TestRunPacketLightsTheRosterRowOfItsLiveL1Window` | `cmd/wisp/panel_l1_blocked_220_test.go:184` | C1 | 220 roster-blockedOnApproval-only-sees-L2-and-second-card-silently-drops |
| 36 | `TestRunBooksWithASnapshotOfItsLiveQueue` | `cmd/wisp/panel_pump_test.go:142` | C1 | 35 panel-bridge-c17 |
| 37 | `TestSnapshotWorkspaceSectionReportsTheNarrowing` | `cmd/wisp/panel_pump_test.go:343` | C1 | 35 panel-bridge-c17 |
| 38 | `TestThePumpIsDrivenNotJustAssembled` | `cmd/wisp/panel_pump_test.go:414` | C1 | 35 panel-bridge-c17 |
| 39 | `TestBookedSnapshotReachesHostBitForBit329` | `cmd/wisp/panel_snapshot_deliver_329_test.go:77` | C1 | 329 panel-snapshot-last-mile-never-delivered |
| 40 | `Test165r2AssembledBridgeReachesPlanPresent` | `cmd/wisp/plan_assembled_165r2_test.go:102` | C1 | 165 plan-first-then-act-a-brake-the-voice-entry-needs |
| 41 | `Test165r2UnlistedWriteRingsItsOwnCardAfterAnApprovedPlan` | `cmd/wisp/plan_assembled_165r2_test.go:187` | C1 | 165 plan-first-then-act-a-brake-the-voice-entry-needs |
| 42 | `TestTicket101ManualSwitchSurvivesRestart` | `cmd/wisp/run_mode101_test.go:216` | C1 | 101 perm-store-has-no-production-importer |
| 43 | `TestTicket101UntouchedConfigRestartsAtDefault` | `cmd/wisp/run_mode101_test.go:339` | C1 | 101 perm-store-has-no-production-importer |
| 44 | `TestTicket101SessionGrantDoesNotCrossRestart` | `cmd/wisp/run_mode101_test.go:404` | C1 | 101 perm-store-has-no-production-importer |
| 45 | `TestTicket101ModeSwitchUsesTheRealL2Gate` | `cmd/wisp/run_mode101_test.go:543` | C1 | 101 perm-store-has-no-production-importer |
| 46 | `TestTicket101UnreadableModeFailsLoudlyAndStrict` | `cmd/wisp/run_mode101_test.go:599` | C1 | 101 perm-store-has-no-production-importer |
| 47 | `TestRunTextTaskTextPathEndToEnd` | `cmd/wisp/run_test.go:169` | C1 | 12 cli-text-path-s1-gate |
| 48 | `TestRunTextTaskFailNextIsClassified` | `cmd/wisp/run_test.go:219` | C1 | 12 cli-text-path-s1-gate |
| 49 | `TestHostDispatchThroughTheAssembledBridge` | `cmd/wisp/run_test.go:260` | C1 | 12 cli-text-path-s1-gate |
| 50 | `TestComposedGateBlocksAWriteForTwoSeconds` | `cmd/wisp/run_test.go:345` | C1 | 12 cli-text-path-s1-gate |
| 51 | `TestRunTextTaskKeyResolvesInTheStore` | `cmd/wisp/run_test.go:443` | C1 | 12 cli-text-path-s1-gate |
| 52 | `TestMissingBlobFailsUnconfiguredNeverSilently` | `cmd/wisp/run_test.go:466` | C1 | 12 cli-text-path-s1-gate |
| 53 | `TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119` | `cmd/wisp/secret_dataroot_119b_test.go:242` | C6 | 119 posix-link-leg-refuses-legitimate-symlinked-data-roots |
| 54 | `TestSecretSetGetListRoundTrip` | `cmd/wisp/secret_test.go:180` | C1 | 63 credential-entry-cli |
| 55 | `TestSecretFromStdinWritesNoIntermediateFile` | `cmd/wisp/secret_test.go:569` | C1 | 63 credential-entry-cli |
| 56 | `TestSecretFailurePathsLogAndPrintNoPlaintext` | `cmd/wisp/secret_test.go:608` | C1 | 63 credential-entry-cli |
| 57 | `TestSecretUnsetRefusesWhileReferenced` | `cmd/wisp/secret_test.go:783` | C1 | 63 credential-entry-cli |
| 58 | `TestSecretSameNameUnderThreeEnvsIsThreeBlobs` | `cmd/wisp/secret_test.go:963` | C1 | 63 credential-entry-cli |
| 59 | `TestSecretPortableModeUsesTicket06Seam` | `cmd/wisp/secret_test.go:1053` | C1 | 63 credential-entry-cli |
| 60 | `TestSecretEndToEndConfigRefResolvesAtRequestTime` | `cmd/wisp/secret_test.go:1086` | C1 | 63 credential-entry-cli |
| 61 | `TestSecretOverwriteIsAnnounced` | `cmd/wisp/secret_test.go:1196` | C1 | 63 credential-entry-cli |
| 62 | `TestRunPacketMarksTheRosterRowACardIsHolding` | `cmd/wisp/subagent_blocked_197_test.go:60` | C1 | 197 the-panel-must-show-each-subagents-state-and-click-through-to-its-own-streaming-work-page-but-go-side-has-zero-subagent-entity |
| 63 | `TestRunPacketCarriesTheSubagentItsRosterRowFed` | `cmd/wisp/subagent_carrier_197_test.go:233` | C1 | 197 the-panel-must-show-each-subagents-state-and-click-through-to-its-own-streaming-work-page-but-go-side-has-zero-subagent-entity |
| 64 | `TestRunPacketReportsTheStreamLogPastItsBound` | `cmd/wisp/subagent_carrier_197_test.go:526` | C1 | 197 the-panel-must-show-each-subagents-state-and-click-through-to-its-own-streaming-work-page-but-go-side-has-zero-subagent-entity |
| 65 | `TestRunPacketCarriesTheStreamDroppedCountPastTheHardCeiling` | `cmd/wisp/subagent_carrier_197_test.go:612` | C1 | 197 the-panel-must-show-each-subagents-state-and-click-through-to-its-own-streaming-work-page-but-go-side-has-zero-subagent-entity |
| 66 | `Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands` | `cmd/wisp/subagent_selfapproval_197_test.go:202` | C1 | 197 the-panel-must-show-each-subagents-state-and-click-through-to-its-own-streaming-work-page-but-go-side-has-zero-subagent-entity |
| 67 | `Test370TruncatedRunRowsAnInterruptedTask` | `cmd/wisp/task370_ac2_test.go:99` | C1 | 370 assembler-maxoutputtokens-is-never-assigned-so-our-ceiling-is-never-sent |
| 68 | `Test370TheInterruptedRowNeedsTheCeiling` | `cmd/wisp/task370_ac2_test.go:143` | C1 | 370 assembler-maxoutputtokens-is-never-assigned-so-our-ceiling-is-never-sent |
| 69 | `TestTaskOutputPointerNotice174r8SpeaksTheRealReason` | `cmd/wisp/task_output_pointer_notice_174r8_test.go:134` | C1 | 174 task-output-points-at-a-file-the-model-cannot-read-by-default-spill-artifacts-are-outside-fs-allowed-dirs |
| 70 | `TestTaskOutputPointerNotice174r8StaysQuiet` | `cmd/wisp/task_output_pointer_notice_174r8_test.go:165` | C1 | 174 task-output-points-at-a-file-the-model-cannot-read-by-default-spill-artifacts-are-outside-fs-allowed-dirs |
| 71 | `TestCompositionRootClosesTheLoopTasksTaintScope` | `cmd/wisp/task_scope_close_151_test.go:55` | C1 | 151 bridge-opentask-sits-on-the-dispatch-hot-path-while-closetask-has-zero-callers-repo-wide-so-the-per-task-taint-scope-is-only-ever-opened |
| 72 | `TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit` | `cmd/wisp/task_scope_close_151_test.go:68` | C1 | 151 bridge-opentask-sits-on-the-dispatch-hot-path-while-closetask-has-zero-callers-repo-wide-so-the-per-task-taint-scope-is-only-ever-opened |
| 73 | `Test378SpawnedChildToolCallsHaveAParentRow` | `cmd/wisp/tasklog_fkcheck_378_test.go:395` | C1 | 378 child-loops-never-open-a-task-log-row-so-spawned-tool-calls-have-no-parent |
| 74 | `Test365RunRowsAStuckTaskAsStuck` | `cmd/wisp/tasklogstate_365_test.go:186` | C1 | 365 two-taskLogState-functions-disagree |
| 75 | `Test365RunKeepsItsWordForTheControls` | `cmd/wisp/tasklogstate_365_test.go:221` | C1 | 365 two-taskLogState-functions-disagree |
| 76 | `TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun` | `cmd/wisp/ticket224_assembly_test.go:81` | C1 | 224 session-scoped-grant-has-zero-executors-and-no-session-identity |
| 77 | `TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking` | `cmd/wisp/ticket224_assembly_test.go:243` | C1 | 224 session-scoped-grant-has-zero-executors-and-no-session-identity |
| 78 | `TestTicket224ProductionSessionDoesNotSurviveRestart` | `cmd/wisp/ticket224_assembly_test.go:421` | C1 | 224 session-scoped-grant-has-zero-executors-and-no-session-identity |
| 79 | `TestZZ315A1RootLoopFirstBridgeCall` | `cmd/wisp/zz315a1probe_test.go:476` | C1 | 315 production-task-scope-never-opens-so-sibling-child-bridge-calls-read-as-unbound-scope-l2 |
| 80 | `TestZZ315A1SiblingChildFirstBridgeCall` | `cmd/wisp/zz315a1probe_test.go:535` | C1 | 315 production-task-scope-never-opens-so-sibling-child-bridge-calls-read-as-unbound-scope-l2 |
| 81 | `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask` | `cmd/wisp/zz315a1probe_test.go:656` | C1 | 315 production-task-scope-never-opens-so-sibling-child-bridge-calls-read-as-unbound-scope-l2 |

## 6. 复跑口令（可逐字重跑；全程只读，不写仓）

```bash
# 1) 快照（不建 worktree、不 checkout；只读 git）
cd "/d/work/workspace/projects plans/Wisp"
git archive --format=tar f4ad34ee1 | sha256sum      # ea56813bb3d1…（本件实际用的是 8a75004ce，见文首）
git archive --format=tar -o /d/tmp/98-attribution-20261105/snap.tar 8a75004ce
# ⇒ sha256 3fba786b92ac…

# 2) 容器（GOMODCACHE/GOCACHE 用命名卷；快照解进容器本地盘，别 bind-mount 仓）
export MSYS2_ARG_CONV_EXCL="*"
docker run --rm -v "D:/tmp/98-attribution-20261105:/out" \
  -v wisp98-gomod:/go/pkg/mod -v wisp98-gocache:/root/.cache/go-build \
  -w / golang:1.27-bookworm bash /out/run.sh
# run.sh 体：mkdir -p /work; tar -xf /out/snap.tar -C /work; cd /work;
#            CGO_ENABLED=1 go test -count=1 -v -timeout 90m ./cmd/wisp/

# 3) 单条取数（例：C2 在干净树上两腿对照）
# linux:  CGO_ENABLED=1 go test -count=1 -v -run 'TestTicket255RosterStillMatchesTheActualReadSites' ./cmd/wisp/
# windows: PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH" \
#          go test -count=1 -run 'TestTicket255RosterStillMatchesTheActualReadSites' -v ./cmd/wisp/

# 4) C6 的仪器说明：`go test ./cmd/wisp/` 在 linux 上必须带 CGO_ENABLED=1（CGO=0 连解析都过不了：
#    "build constraints exclude all Go files in sherpa-onnx-go-linux"，scripts/wisp-cli-tests.sh:22-24 已记）。
```

## 7. 本件没有建立的事（宁缺毋造）

- **没有取 CI 的 run id / job id**。AC#1(iii) 那一格仍是编排者的；本件全部读数是本机容器，**不得**被引成 CI 读数。
- **没有跑整包 Windows 腿**，只点了 §4 表里那 10 枚（7 枚 C1 抽样 + 2 枚门 + C2 一枚）；剩下 71 枚**没有在 windows 上点过**。
  所以"windows 腿全绿"本件**不主张**；反过来，
  C2 在干净树上实测红 ⇒ **该 HEAD 上 windows 那条 CLI 步不可能绿**（这一条我主张，尺见 §4）。
- **没有判断 C1 那 76 枚该怎么修**。两条路都要用例主人点头，且代价不同：
  ① 给文件加 `//go:build windows` / 改名 `_windows_test.go`（票 70 / 票 131 的先例）——但这会把用例**从 ubuntu 分母里拿掉**，
     正是票 98 AC#2 划的那条线；② 把夹具换成 `env:` ref（C28 的 portable 形状）——用例留在这条腿上，但动的是别人票的夹具。
  本件只负责"点出是谁的"。
- **没有替 C2 改 `config_readers_255.go`**：那枚文件当时在飞（台账里有 `255-r3`/`255-r4` 的腿），改它就是撞别人写面。
- **没有验证 `-count=2`**（票 98 AC#5 要求的是门禁口径，不是普查口径）。`-count=1` 的并发面噪声（例如 `t.TempDir` 清理）
  本件没有单独量；若要交付门禁读数，按 AC#2 的命令重跑两发再点名。
- **`HEAD` 在普查期间动过两枚**（`8a75004ce`、`a980734f0`；后者改 `internal/agent/approval/gate.go`，`cmd/wisp` 会用到）。
  本件的**主读数**钉在 `f4ad34ee1` 的 `cmd/wisp`；**第二发**在 `a980734f0` 上重跑过（§1），红/绿名字清单逐枚相同 ⇒ 结论在这两枚上都成立。
