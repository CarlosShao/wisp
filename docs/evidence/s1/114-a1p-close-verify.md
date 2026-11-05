# 114-a1p 收口复核（只读腿，非实现者）

- 腿: `114-a1p`（close-verification，read-only）
- 日期: 2026-11-05 09:0x +08
- 起手 HEAD: `25203b526b3a44eb751b3f1f8951907f9a00365b`
- 票: `.scratch/wisp/issues/114-composer-request-has-no-production-caller-and-the-real-gate-must-be-native.md`（206 行）
- 要回答的问题: 票 371 刚改过 `internal/panel/composer_handlers.go`（114 的 AC#2 门所在文件，现经 `perm.JudgeCard`）并改了 `composer_handlers_test.go`。收 114 之前必须证明 371 没有打坏 114。
- 本腿纪律: 只读；未改任何生产代码；唯一新建/提交件＝本文件；无 `add -A` / amend / reset / stash；未 push。

## 0. 结论（先行）

| 项 | 判 | 依据 |
|---|---|---|
| 114 可收口（closeable） | **YES** | 见 §5/§6；AC#6/AC#11 按已裁定处置保持未勾 |
| 371 削弱 114 判据 | **NO** | §2：测试改动仅夹具签名，6 条 zero-call 断言逐字未动，仅行号 +1 |
| AC#6 残余处置仍有效 | **YES** | §6；真机差分截屏需 GUI/WebView2，本机 `cmd/wisp` 在 MARK299 被硬终止（A1191） |
| AC#11 作废处置仍有效 | **YES** | §6；前端渲染 CI 步，结构上不属本编队（`frontend/**`/`design/**` owner-delegated） |
| AC#2 门仍在且正确 | **YES** | §4；`internal/panel/composer_handlers.go:176` `perm.JudgeCard(to, h.Confirm != nil) == perm.GateNoChannel` |
| 未勾框集合 | **{AC#6, AC#11}** 恰两枚 | §5 `grep -nE '^- \[ \]'` 仅命中 `:52`、`:75` |
| 既有无关红 | 1 枚 | `TestComposerContractTypesMatchFrontend`（frontend 契约漂移，371 之前即红；§3） |

## 1. 前序裁决提交的祖先关系

命令（Git Bash，先 cd 到 repo）：

```
git merge-base --is-ancestor 387e2d2b HEAD ; echo rc=$?
git merge-base --is-ancestor 3efb53526 HEAD ; echo rc=$?
```

读数：

```
rc=0
rc=0
```

结论：`114-v1` 裁决提交 `387e2d2b` 与 `114-v2` 裁决提交 `3efb53526` 都在 `HEAD 25203b526` 的祖先里，裁决未被历史丢弃。

## 2. 371（`8e0702fdb`）对 114 两文件的改动

提交：`8e0702fdb 371-r1p: AC#2' 落地（门要求"答过允许"，只看 confirm 非 nil）＋ AC#2'-b 两处同判`（2026-11-04 22:51:46 +0800）。

`git diff --stat 8e0702fdb^ 8e0702fdb -- internal/panel/composer_handlers_test.go internal/panel/composer_handlers.go`：

```
 internal/panel/composer_handlers.go      | 84 ++++++++++++++++++++++++--------
 internal/panel/composer_handlers_test.go |  7 +--
 2 files changed, 69 insertions(+), 22 deletions(-)
```

测试文件只有一个 hunk（`@@ -53,13 +53,14 @@`），内容全部是夹具 `countingConfirm`：

- 注释 +1 行（"Since ticket 371 it reports perm's answered credential..."）；
- 签名由 `func(...) error` 改为 `func(...) (perm.Card, error)`；
- `return nil` 改为 `return perm.AnsweredCard(""), nil`。

**该 hunk 没有触碰任何断言。** 114 的 6 条 "注入的 ModeWriter / Confirm 一次都没被调用" 判据，371 前（`8e0702fdb^`）与 371 后（HEAD 同文件）逐字对比如下（371 后仅因上述注释 +1 行而整体下移一行）：

| 判据 | 371 前 | 371 后（HEAD） | 文本 |
|---|---|---|---|
| ModeWriter 未被调用（gate 拒） | `:166 if w.setCalls != 0 {` / `:167 t.Fatalf("... the injected ModeWriter was called %d time(s) despite the refusal - the gate did not stop the write"` | `:167` / `:168` | 逐字同 |
| Confirm 未被调用（不举第二张卡） | `:242 if card.calls != 0 {` | `:243` | 逐字同 |
| ModeWriter 未被调用（current 档不可读） | `:317 if w.setCalls != 0 {` | `:318` | 逐字同 |
| ModeWriter 未被调用（未知档） | `:330 if w.setCalls != 0 {` | `:331` | 逐字同 |
| ModeWriter 未被调用（空目标） | `:340 if w.setCalls != 0 {` | `:341` | 逐字同 |
| ModeWriter 未被调用（错路由） | `:352 if w.setCalls != 0 {` | `:353` | 逐字同 |

命令与读数：

```
# 371 前
$ git show 8e0702fdb^:internal/panel/composer_handlers_test.go | grep -n "setCalls != 0\|card.calls != 0"
166: if w.setCalls != 0 {
242: if card.calls != 0 {
317: if w.setCalls != 0 {
330: if w.setCalls != 0 {
340: if w.setCalls != 0 {
352: if w.setCalls != 0 {

# 371 后（HEAD）
$ grep -n "setCalls != 0\|card.calls != 0" internal/panel/composer_handlers_test.go
167: if w.setCalls != 0 {
243: if card.calls != 0 {
318: if w.setCalls != 0 {
331: if w.setCalls != 0 {
341: if w.setCalls != 0 {
353: if w.setCalls != 0 {
```

结论：**371 没有削弱 114 的测试（weakened = NO）**。它只把夹具从 `error` 形状改成 `perm.Card` 形状并返回 `perm.AnsweredCard("")`，这是 371 新契约（"答过允许"）要求的夹具适配；zero-call 判据一条未删、未放宽、未改写。

补充：`git diff --stat 8e0702fdb HEAD -- internal/panel/composer_handlers.go internal/panel/composer_handlers_test.go` 输出为空，即 371 之后 HEAD 没有再动这两个文件。

## 3. 测试重跑

命令（派单原文）：

```
go test ./internal/panel -run "114|Composer|Mode|Gate" -v -count=1 -timeout 200s
```

读数：`rc=1`；`=== RUN` 35（含子用例）；顶层 `--- PASS` 21；顶层 `--- FAIL` 1。

### 3.1 过滤器射程缺口（本腿发现，已用补充跑补上）

上面的正则 **不匹配** 114 三条核心 zero-call 用例的名字（名字里没有 `114` / `Composer` / `Mode` / `Gate` 子串）：

- `TestOnlyTheAutoApproveTargetIsRefusedWhenNoL2ConfirmLegIsAttached`
- `TestAWideningRequestWithAnAttachedLegReachesTheOneWriterAndRaisesNoSecondCard`
- `TestTheWriteLegsOwnRefusalIsReturnedUntouched`

故本腿另跑一条具名补充命令：

```
go test ./internal/panel -run "TestOnlyTheAutoApproveTargetIsRefusedWhenNoL2ConfirmLegIsAttached|TestAWideningRequestWithAnAttachedLegReachesTheOneWriterAndRaisesNoSecondCard|TestTheWriteLegsOwnRefusalIsReturnedUntouched" -v -count=1 -timeout 200s
```

读数：`TEST_RC=0`；三条全 PASS：

```
--- PASS: TestOnlyTheAutoApproveTargetIsRefusedWhenNoL2ConfirmLegIsAttached (0.00s)
--- PASS: TestAWideningRequestWithAnAttachedLegReachesTheOneWriterAndRaisesNoSecondCard (0.00s)
--- PASS: TestTheWriteLegsOwnRefusalIsReturnedUntouched (0.00s)
```

### 3.2 通过名单（过滤器跑 + 补充跑）

114 的 6 枚顶层用例全 PASS：

```
--- PASS: TestOnlyTheAutoApproveTargetIsRefusedWhenNoL2ConfirmLegIsAttached        (补充跑)
--- PASS: TestAWideningRequestWithAnAttachedLegReachesTheOneWriterAndRaisesNoSecondCard (补充跑)
--- PASS: TestAStricterOrEqualModeRequestNeedsNoConfirmLeg                        (过滤器跑)
--- PASS: TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete      (过滤器跑，5 子用例全 PASS)
--- PASS: TestTheWriteLegsOwnRefusalIsReturnedUntouched                           (补充跑)
--- PASS: TestTheModeLadderThisGateJudgesAgainstIsTheDocumentedThree              (过滤器跑)
```

371 / 372 的判据同跑 PASS：

```
--- PASS: TestPanelGateAndStoreGateAgreeOnEveryPair367
--- PASS: TestPanelGateAndStoreGateAgreeOnTheAnsweredCard371
--- PASS: Test371R1PModeConfirmCarriesTheCredential
--- PASS: Test371R1PThePanelGateCallsTheSameCardPredicate
--- PASS: TestRoutingDoesNotBypassTheNativeModeGate
--- PASS: TestRawEnvelopeReachesTheAttachedModeHandlerExactlyOnce
--- PASS: TestPlantedComposerModeWriteGoesRed
--- PASS: TestModeViewHasNoWriteSurface
--- PASS: TestModeRequestResolvesThroughRisksOwnParser
--- PASS: TestUnknownModeNeverRendersAsASafeOne
--- PASS: TestComposerRenderFixtureTellsTheTruth
--- PASS: TestComposerStateSurvivesPanelCloseAndReopen
--- PASS: TestComposerCurrentModelTravelsOnlyFromItsReader
--- PASS: TestGitDimensionHasNoModelCallableTool
--- PASS: TestAnUnreadableModeFromALiveReaderStillRendersUnknown
--- PASS: TestComposerEnvelopeAcceptsItsFourRequests
--- PASS: TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests
--- PASS: TestFrontendComposerRequestsMatchTheEnvelope
```

### 3.3 唯一红（既有、无关）

```
    composer_test.go:78: Go ComposerState emits [files] that interface ComposerState does not declare
--- FAIL: TestComposerContractTypesMatchFrontend (0.00s)
```

判"既有且与 371 无关"的依据（三条独立读数）：

1. 371 的 diff 只碰 `composer_handlers.go` 与 `composer_handlers_test.go`；`TestComposerContractTypesMatchFrontend` 读的是 `internal/panel/composer.go` 的 `ComposerState` 与 `frontend/src/lib/panel.ts`（`composer_test.go:48-83`），两者 371 一字未动。
2. 371 树（`8e0702fdb`）里已经存在这个不匹配：`git show 8e0702fdb:internal/panel/composer.go | grep 'json:"files'` -> `:285 Files *TreeView \`json:"files,omitempty"\``；`git show 8e0702fdb:frontend/src/lib/panel.ts` 的 `ComposerState` 段 `grep -c files` -> `0`。`files` 键由票 190 的 `4d3daacea` 引入，且 `git merge-base --is-ancestor 4d3daacea 8e0702fdb` -> `rc=0`。
3. 371 之前的 114 各裁决件已把这枚测试列入在册红集：`docs/evidence/s1/114-v2-verdict.md:114`、`114-r1-impl.md:134`、`114-r2-impl.md:120`、`114-r3-impl.md:121`。

结论：本腿**新增红 = 0**；`rc=1` 由这枚既有的前端契约漂移造成，不属于 114/371 的门。

## 4. AC#2 的门（仍在，且口径正确）

调用点（逐字）：

`internal/panel/composer_handlers.go:176`

```go
	if perm.JudgeCard(to, h.Confirm != nil) == perm.GateNoChannel {
```

该判据的定义（`internal/perm/direction.go:57-65`）：

```go
func JudgeCard(to risk.Mode, channel bool) CardGate {
	if !RequiresCard(to) {
		return GateNoCard
	}
	if !channel {
		return GateNoChannel
	}
	return GateAsk
}
```

`RequiresCard` 的唯一目标（`internal/perm/direction.go:23`）：

```go
func RequiresCard(to risk.Mode) bool { return to == risk.ModeAutoApprove }
```

一致性核对（371 前 -> 371 后，nil 腿时行为逐格相同）：

| `to` | `h.Confirm` | 371 前 `RequiresCard(to) && h.Confirm == nil` | 371 后 `JudgeCard(to, h.Confirm != nil) == GateNoChannel` | 结果 |
|---|---|---|---|---|
| auto_approve | nil | true（拒） | `GateNoChannel` -> true（拒） | 同 |
| auto_approve | 非 nil | false（放行到 Set） | `GateAsk` -> false（放行到 Set） | 同 |
| 其他档 | nil | false（放行） | `GateNoCard` -> false（放行） | 同 |
| 其他档 | 非 nil | false（放行） | `GateNoCard` -> false（放行） | 同 |

即：114 AC#2 立的门（auto_approve + 无 L2 腿 = 在任何 `ModeWriter.Set` 之前拒）仍在；371 只把"目标 + 通道"两项折进 `perm.JudgeCard`，与 `perm.Store.Set` 的同一判据共用一处定义，符合票 372/M4 的"只有全自动要确认"（`direction.go:23` 注释与票面一致）。

门位置仍在写入腿之前：无变化分支的 `Set` 在 `:167`（gate 之前），主 `Set` 在 `:188`（gate 之后），gate 在 `:176`。

## 5. 未勾框清单（必须恰为 {AC#6, AC#11}）

命令：

```
grep -nE '^- \[ \]' .scratch/wisp/issues/114-*.md
```

读数（逐字）：

```
52:- [ ] **AC#6** 可见性：`R-92-5` 欠的那一笔——真机**差分截屏**（改档位前后各一张，能看出输入框上的档位文字变化）。
75:- [ ] **AC#11**（`R-92b-3`）：把 `npm run render:composer && git diff --exit-code` 做成 CI 一步
```

**恰两枚：AC#6、AC#11。** 无第三枚（`[x]` 的 9 枚为 AC#1/#2/#3/#4/#5/#7/#8/#9/#10）。

## 6. 逐格证据与两枚未勾框的处置

| AC | 框 | 证据（文件名 / 提交） |
|---|---|---|
| AC#1 | [x] | `docs/evidence/s1/114-ac1-status-table.md` |
| AC#2 | [x] | `internal/panel/composer_handlers.go:176`（门）+ `internal/panel/composer_handlers_test.go:134`（用例）+ `docs/evidence/s1/114-ac2-native-gate.md`；装配证据 `docs/evidence/s1/114-r2-impl.md` |
| AC#3 | [x] | `docs/evidence/s1/114-v2-verdict.md` §11（判语）+ `docs/evidence/s1/114-r2-impl.md`（接线/举卡）；裁决提交 `3efb53526`。编排者裁：`auto_approve` 走 L2 原生卡＝成立（接线层）；"更宽档也要卡"半支作废（票 372/M4，A1214） |
| AC#4 | [x] | `docs/evidence/s1/ticket114-114-v1-accept-v1.md`（`387e2d2b`）；尺 `internal/panel/composer_test.go:697 TestDoorTwoReadsEveryTextTypeTheScannerWalks` |
| AC#5 | [x] | `docs/evidence/s1/ticket114-114-v1-accept-v1.md` §4（抽核判对）；`internal/agent/approval/ui.go:155` 注释 "PanelAPI has no Allow"；`internal/panel/config_handlers.go:109` `lockedFieldFamilies`；`TestPlantedComposerModeWriteGoesRed`、`TestModeViewHasNoWriteSurface`；零 emoji 由 `scripts/d22scan.sh` 覆盖 |
| AC#6 | **[ ]** | **保持未勾：残余（residual）** |
| AC#7 | [x] | `docs/evidence/s1/ticket114-114-v1-accept-v1.md`（gofumpt `-l` 零行、`bash -n` rc=0、`d22scan.sh` rc=0、panel 整包红集与对照同集） |
| AC#8 | [x] | `docs/evidence/s1/ticket114-114-v1-accept-v1.md`（`387e2d2b`）；门 `internal/panel/bridge_test.go:184 TestEveryWhitelistedInboundNameIsAlsoSentByThePage` |
| AC#9 | [x] | `docs/evidence/s1/ticket114-114-v1-accept-v1.md`（扫描根 `frontend/src` -> 整棵 `frontend`；`TestDoorTwoReadsEveryTextTypeTheScannerWalks`） |
| AC#10 | [x] | `docs/evidence/s1/ticket114-114-v2-accept-v1.md`（`3efb53526`）；`internal/panel/bound_dispatch_door_114_test.go:444 TestTheBoundDispatchDoorIsReachableFromThePage`；产码 `3b6f170f` |
| AC#11 | **[ ]** | **保持未勾：作废（void）** |

### AC#6 处置仍有效（YES）

- 票面 `:197-199` 裁定：AC#6 = `skipped=frontend`，需要**真窗（WebView2）差分截屏**；本机 `cmd/wisp` 整包在 `MARK299` 被硬终止（A1191 实测）=> 不是"没人做"，是"本机做不到"，转残余，等有 GUI 会话或 CI 上做。
- 本腿复核未发现能改变这一处置的新证据；`frontend/**`/`design/**` 归机主外部 agent，本编队结构上不做。保持 `[ ]`。

### AC#11 处置仍有效（YES）

- 票面 `:200-202` 裁定：AC#11 是**前端渲染**的门（`npm run render:composer && git diff --exit-code`），结构上不属本编队零前端范围 => 作废，不是残余。
- 本腿复核未发现该结构边界变化。保持 `[ ]`。

## 7. 收口判语

**114 可收口（closeable = YES）**，前提是沿用已记录的裁决：

- AC#1-5、AC#7-10 全部 `[x]`，证据件与提交均在（§6）；
- AC#6 保持未勾（残余：真机差分截屏，需 GUI/WebView2，非本编队可做）；
- AC#11 保持未勾（作废：前端渲染 CI 步，结构上不属本编队）；
- 371 未打坏 114：门仍在 `composer_handlers.go:176`，6 条 zero-call 断言逐字未动，`go test` 的 114/371 关键用例全 PASS（含 §3.1 补充跑）；
- 唯一红 `TestComposerContractTypesMatchFrontend` 是 371 之前即存在的 frontend 契约漂移（`files` 键，票 190），本腿新增红 = 0。

### 一条非阻断观察（不属本腿可改）

票面 `:9-12` 的 `Status:` 仍是 `in-progress`，且该行说 `AC#3/#4/#5/#7/#8/#9/#10/#11 未做`，与现框面（9 枚 `[x]`）不一致。这是**票面状态行未随裁决更新**，不是代码/判据问题；改 `Status:` 属编排者动作，本只读腿不动它。收口时请编排者一并更正。

## 8. 纪律

- 本腿只读；未改任何生产代码（`git diff --stat 8e0702fdb HEAD -- internal/panel/composer_handlers.go internal/panel/composer_handlers_test.go` 为空；工作树里本腿零产码改动）。
- 唯一新建/提交件：本文件（显式 pathspec 提交）。
- 未使用 `git add -A` / `--amend` / `reset` / `stash`；未 push。
- 注意：`docs/evidence/s1/` 下另有 `351-ac2-ac3-verdict.md`、`378-r1-impl.md` 两枚**未跟踪**文件，属他腿，本腿未触碰。
