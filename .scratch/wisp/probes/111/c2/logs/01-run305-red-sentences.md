# 01 · run 305（head `cc31526`）`cmd/wisp` 那一步的六枚红句原文

尺与出处（本腿自取，⛔ 未把整发日志读进上下文）：
`gh run view --repo CarlosShao/wisp --log-failed 37703959747` → 仓外 `D:/tmp/wisp111c2/run305.log`（2 014 234 B 级，`wc -c` 现量见 `02-*.md`）
→ 切步：`awk -F'\t' '$1=="test-windows"'` → `awk '/##\[group\]/{insec=($0 ~ /wisp-cli-tests\.sh/); if(insec) next} insec'` → `cli-run305.txt`（**2 190 行 / 404 963 B**）
→ 去前缀：`awk -F'\t' '{out=$3; sub(/^[0-9][0-9-]*T[0-9:.]*Z /,"",out); print out}'` → `clean-run305.txt`

步末四数原文（逐字）：
`runtests.sh: go test exited 1 - packages=[./cmd/wisp/ -count=1 -skip ^(TestDefaultDeadlineWallClockMeasurement|TestSubprocessCrashWriter|TestHelperProcess|TestLiveWasapiSmoke|TestRealDownloadVadThroughPipeline|TestRealDownloadPuncArchiveThroughPipeline|TestSyncRegistryProbeLive)$] top-level: PASS=230 FAIL=6 SKIP=2, === RUN=335, '[no tests to run]'=0`
`portable-tests.sh: four numbers (all from -v output): === RUN=335  --- PASS=230  --- FAIL=6  --- SKIP=2`

顶层红（`grep -E '^--- FAIL: '`＝6 行，缩进子测试红＝**0 行**，`grep -cE '^\s+--- (FAIL|SKIP): '`＝0）：

## 1 `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey` (1.34s) — `clean-run305.txt:252`
```
always_write_no_clobber_226_test.go:91: fs.allowed_dirs on disk = [C:/Users/RUNNER~1/AppData/Local/Temp/TestAC1AlwaysBranchDoesNotRevertAHandEditedKey94708015/002 C:/Users/runneradmin/AppData/Local/Temp/TestAC1AlwaysBranchDoesNotRevertAHandEditedKey94708015/003], want the stored rule "C:/Users/RUNNER~1/AppData/Local/Temp/TestAC1AlwaysBranchDoesNotRevertAHandEditedKey94708015/003"
```
同一目录的两种拼法：期望侧是 8.3 短名 `RUNNER~1`，盘上存的是长名 `runneradmin`。

## 2 `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` (1.07s) — `clean-run305.txt:485`
```
approval_always_201_test.go:141: config.toml did not gain the stored rule "C:/Users/RUNNER~1/AppData/Local/Temp/TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card3472203531/003"; it holds:
```
（同族；红句后面跟的是整份 `config.toml` 的转储，本件不抄。`--- FAIL` 距该句 221 行，因为转储占了中间那段。）

## 3 `TestTicket223PermissionDeniedSitsInItsOwnSentence` (1.08s) — `clean-run305.txt:780`
```
config_reload_perm_223_windows_test.go:95: attempt 1: icacls denied nothing, so this case cannot show the permission sentence
```
产码侧对应（`cmd/wisp/config_reload_perm_223_windows_test.go:91-95`，本腿只读）：先 `runIcacls223(..., "/deny", everyoneSID+":(R)")`，再作正控 `os.Stat` 必须成功、`os.ReadFile` **必须失败**；CI 上 `ReadFile` 返回了 `err == nil` ⇒ 那一发是**正控自灭**（判据没被触发），不是产品答错。

## 4 `TestRunPacketCarriesTheLoadedInstructionFiles` (1.07s) — `clean-run305.txt:1009`
```
instructions_200r2_test.go:167: the packet's instructions carry no entry for C:\Users\RUNNER~1\AppData\Local\Temp\TestRunPacketCarriesTheLoadedInstructionFiles3434685345\002\instr-ws\AGENTS.md: [{Path:C:\Users\runneradmin\AppData\Local\Temp\TestRunPacketCarriesTheLoadedInstructionFiles3434685345\002\instr-ws\AGENTS.md Tier:project Depth:0 Bytes:269 Source:fs.read}]
```
同一枚文件、两种拼法逐字并排出现在同一行里（`RUNNER~1` vs `runneradmin`），字节数一致（269）。

## 5 `TestPanelHostRealWindowHopAndLifecycle` (18.94s) — `clean-run305.txt:1198`
```
    panel_host_windows_test.go:660: cold bring-up measured on this box: -1.000 ms (HEAD cc31526 at read time, 2026-10-07T23:49:33Z)
    panel_host_windows_test.go:662: cold bring-up did not produce a browser round trip (got -1.000); the message channel did not come up on the real window
```
`:662` 是 `t.Fatalf`（`coldMs <= 0` 那一支）⇒ 用例在这一点**终止**，后面 AC#1/AC#3/AC#4 那些读数这一发**一个都没打**（对照 run 304：打了，见 `02-*.md`）。

## 6 `TestAC14GoSideEvalPushReachesThePage` (0.62s) — `clean-run305.txt:1293`
```
    panel_resident_windows_test.go:867: AC#14 nail 2 (Eval push hop), page's own words: title=""
    panel_resident_windows_test.go:869: Go's Eval push did not reach the document: the page reports its title as "", want "PUSHED-33R5-OK". This is the push dimension, separate from the awaited-reply dimension in TestAC14AwaitedBindingReplyReachesThePage - one arriving says nothing about the other
```
⚠ 派单与 `111-c1` 引这枚时都在 `...want "PUSHED-33R5-OK"` 处截断；原文后面还有半句"这是 push 维度、与 awaited-reply 维度互不作证"。**"红句"的完整形状包含"页自己报回了 title"** 这一事实——见 `05-*.md` 的推论。

## 同场两枚 SKIP（不是红名，但同一把尺会判红）
```
:1204 --- SKIP: TestPanelHostLatencyPercentilesAC2  panel_host_windows_test.go:985: no cold/hot sample recorded in this process: the lifecycle test did not run in this binary (use -run 'TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2'). Named skip - an empty aggregate is not a green latency gate
:1260 --- SKIP: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe  panel_resident_windows_test.go:317: AC#13 has no subject in this tree: the embed resolves no entry
```
成因读数（同发 `panel_host_gate_test.go:109`）：`shape=anchor-only built=false entry-bytes=0 entry-err=panel: embedded assets are not built (run npm run build in frontend/)`。

rc=0
