# audit-fixes 20261117 — privacyConfirmed 门（privacy 两扇破坏性门 + roster.cancel 停根）证牙与现跑读数

执行腿：wauditfix 续跑（2026-11-17 21:5x 起，VPN 断线死腿收口）。
派单：`.scratch/wisp/probes/orch/20261117/txt/wauditfix-resume-dispatch.txt`。
对抗审计源：`docs/evidence/s1/readface-adversarial-20261117-1432.md` §9.1（R1 两类）。

## 一、改动是否就是三条（审读结论：是，且前端置位只落在确认键唯一发送点）

1. **privacyConfirmed 门（privacy.delete-one / privacy.purge）**：
   `cmd/wisp/panel_privacy_40.go` — `privacyConfirmed40(req)`（读 `ComposerRequest.PrivacyConfirmed`，JSON 键 `privacyConfirmed`）在两扇破坏性门上、**`openStore()` 之前**判定：无位 ⇒ `errNoPrivacyConfirm(verb)` 具名拒绝（拒绝句点名「需要确认位 privacyConfirmed:true：privacy.<verb> …」），存储句柄根本不打开 ⇒ 零删除由构造保证；有位 ⇒ 执行。export 不设门（export 不写任何东西）。
2. **roster.cancel 停根同门**：`cmd/wisp/panel_roster_cancel_wave2.go` — seam 签名扩为 `(taskID, confirmed)`，root 分支（root 是整轮任务，破坏性动作）无位 ⇒ 具名拒绝「停根任务需要确认位 privacyConfirmed:true …」、root 句柄零触碰；有位 ⇒ `root.Cancel()`。子行走既有 `rt.tasks.Cancel` 单行停路径，不设门（位在那里是惰性的）。
3. **真牙**：`panel_roster_cancel_wave2_test.go:216` 已有 `if stoppedRow.State != statemachine.StateMuted`（真停后子任务终态是 D43 的 Muted，假停会让它结算到别的 D43 词）。
4. **信封与前端**：`internal/panel/bridge.go` — `PrivacyConfirmed bool \`json:"privacyConfirmed,omitempty"\`` 请求体字段（不是新 C17 方法，C17 名册不变）。前端 `lib/panel.ts` 只在三个函数置位：`requestPrivacyDeleteOne`（发起删除键＝一次性确认，票 40 定案）、`requestPrivacyPurge`（唯一路径＝确认卡「确认清空」键）、`requestRosterCancel`（唯一路径＝监控取消卡「确认停止」键，`monitor-popover.tsx` 的 `confirmStop`，停止按钮只开卡不发请求）；export、history 及其余请求一律不带位。`privacy-settings.tsx` 的删除卡描述写明位的存在。缺口复核：**无补丁缺口**——三个置位点各自就是该动词的唯一发送点（现查：`grep requestRosterCancel|requestPrivacyDeleteOne|requestPrivacyPurge` 全部前端调用方核实）。

## 二、证牙（变异只在仓外副本 `%TEMP%\wauditfix-mut1`，仓内零变异残留，复查 `grep MUTATION-` = 无）

| 牙名 | 变异 | 红/绿读数 |
|---|---|---|
| `Test40PrivacyDestructiveRefusalWithoutTheConfirmBit` | 副本内 `privacyConfirmed40` 临时恒 `true` | **红**：`--- FAIL: Test40PrivacyDestructiveRefusalWithoutTheConfirmBit (0.16s)`，其余 5 例 Test40Privacy* 仍 PASS（正控不受门影响，红得干净）；还原重编 **绿**：`--- PASS … (0.20s)`，rc=0 |
| `TestRosterCancelDoorStopsARealSubagentAndTheForensicsCarryTheAppliedSteps`（:216 StateMuted 牙） | 副本内 roster 实现换「返回 true 不真停」（子行路径跳过 `rt.tasks.Cancel` 直接 `return true, ""`） | **红**：`panel_roster_cancel_wave2_test.go:217: settled state = "Settling", want D43's Muted`（假停下子任务结算到自己的非 Muted 词），`--- FAIL … (2.67s)`，rc=1；还原重编 **绿**：终态=Muted、`--- PASS … (2.59s)`，rc=0 |
| 副本基线（还原态全量） | — | `Test40Privacy\|TestRosterCancel` 全 PASS，rc=0（证明红/绿差异只来自变异本身） |

副本补齐的编译面：`forks/`（go.mod replace 的 go-webview2）、`tools/`（mockllm）、stage-cli 的三枚 pinned DLL（onnxruntime 1.28.2，sha256 校验过的同一批）。

## 三、三样现跑（仓内、2026-11-17 22:00 前后）

1. **d22scan**：`d22scan.sh: OK - control rc=0, scan rc=0`；末行 `d22scan rc=0`，0 findings（clean - no D22 ban violations）。
2. **载具** `scripts/wisp-cli-staged-tests.sh -test.run 'Test40Privacy|TestRosterCancel'`（跑前 `Get-Process wisp.test` 空，无裸跑）：10 例全 PASS（Test40Privacy* 六例 + TestRosterCancel* 四例），末行 **`ok  github.com/CarlosShao/wisp/cmd/wisp	11.00s`**（载具只在 exit 0 时打 `ok`），rc=0。
   - 根任务确认位读数（test 自记）：`无位拒绝=panel: 任务 …没有停下：停根任务需要确认位 privacyConfirmed:true…；带位回执：已对任务 …发出停止…；exit=1`（cancelled 类）。
   - 取消后 applied-steps 读数：`终态=Muted；已落地 1 步（tool_call success，tool=fs.read）；被中止 0 步`。
3. **go build ./...** rc=0。
4. **前端**：`npm run typecheck` rc=0（tsc -b 无错）；`npm run build` rc=0（✓ built）；`npm run lint` rc=0（45 warnings 0 errors，warnings 为既有存量非本波新增断言放宽）。

## 四、首跑注记（非缺陷）

第一次载具调用经 pwsh→bash 直传 `-test.run '…'` 被拆参（`flag provided but not defined: -test`，rc=2，setup 失败、零测试执行）；改 `bash -c` 内层引号后重跑即上述绿读数。测试本体未动。

## 五、写入面

本波只写：`cmd/wisp/panel_privacy_40.go`、`panel_privacy_40_test.go`、`panel_roster_cancel_wave2.go`、`panel_roster_cancel_wave2_test.go`、`frontend/src/components/privacy-settings.tsx`、`frontend/src/lib/panel.ts`、`internal/panel/bridge.go`、本证据件。未触碰：`.scratch/wisp/issues/34-frontend-scaffold.md`(D)、`44-costmeter-c23-done.md`(M)、`frontend/fixtures/composer-states.html`(M)、`internal/speech/mic/resident_audio_291_probe_windows_test.go`(M)、`.scratch/wisp/probes/**`、thresholds/golden/allowlist/PLAN/specs/D43/design。无 push。断言只加未放宽。
