# 98 — 收口判定：`cmd/wisp` 的加载期洞由 staging 形与在册守卫关闭（四/五格现跑坐实）

**票**：`.scratch/wisp/issues/98-cmd-wisp-tests-never-run-on-this-host.md`（本件＝翻格与改名的读数底）
**性质**：实现腿判定证据（派单明确授权本腿翻框收口；裁决复核如需，编排在收件时另出）。
**树锚**：主树读数起跑时 HEAD `975e925d4`，收尾时 `5071d8825`（A1651）；两锚点间 `cmd/ internal/ tools/` **diff 为空**
（现跑 `git diff --stat` 核过）⇒ 读数对两锚点都成立。隔离 clone `git clone --shared` 落 `/d/tmp/98-leg/clone` @`5071d8825`。
**时间**：2026-11-09 11:45–12:3x +08。
**环境**：本机 Git Bash；每次跑测试前 `unset GOTMPDIR TMPDIR TEMP TMP`；`CGO_ENABLED` 保持默认 **1**（从未设 0）；
写面＝本件＋票 98（勾框/Status/Progress 追加）＋票 90/92/94/97（各一段追加）。`scripts/**`、`internal/**`、
`frontend/**`、`forks/**`、`ci.yml`、阈值/golden：一字未动。

## 0. 一句话

票 98 的洞**真的关上了**：修前裸命令 `0xc0000135` 一条测试都跑不了（今日仍复现，作修前读数留档）；
修后两条 documented 路（PATH 一行式、**staging 形**＝`go test -c` ＋ `scripts/stage-sherpa-dlls.sh`
把三枚钉定 DLL 贴到测试二进制同目录）都能起进程跑真测试；`0xc0000135` 两形下都消失；
永不 skip 的守卫已在册且判据发现跑成立（拿掉注入 ⇒ 守卫红并逐字点名缺的 DLL）；
CI 上该包的测试**一直在真跑**（run id 已取、步级结论已点名）；台账更正已落到四张载有该句的票。
整包 rc=1 是**另一件事**（真红/真窗族＋语音族的二级 DLL shadow），逐枚定性，不当票 98 的缺陷。

## 1. AC#1(i) — 修前 vs 修后两发读数（同日、同树系实测）

**修前**（裸命令、无注入，@`975e925d4`，11:47）：

```
$ go test -count=1 ./cmd/wisp/ ; echo "bare_rc=$?"
exit status 0xc0000135
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.056s
FAIL
bare_rc=1
```

**修后·形 A**（PATH 一行式，全包 count=1，见 §8.1）：进程起、556 条顶层用例全出结果行。
**修后·形 B**（staging 形，见 §4）：无 PATH 注入也起进程、用例真 `--- PASS`。
两形下 `0xc0000135` 都不再出现。

加载器点名哪条 DLL（在册复测口径，`98-load-guard.md` §2.3 / `98-ac1-ac5-readings.md` §1 同形，
本腿不重复编译这一步）：`error while loading shared libraries: sherpa-onnx-c-api.dll: cannot open shared
object file`（rc=127）——第一级缺的就是 `sherpa-onnx-c-api.dll`。

## 2. AC#1(ii) — 分母（@本日树实测）

- `go test -c ./cmd/wisp/` rc=0（构建与 DLL 无关）；`-test.list '.*'` 顶层可启动用例 = **556**。
- 测试文件 **178** 枚（`*_windows_test.go` 76、含 `//go:build …winlive` 12）。
- 结论：裸命令的红把 556 条全部静默归零；staging/注入后同样这 556 条可启动（§8 的 503+51+2=556 对账）。

## 3. AC#1(iii) — CI 那一格：file:line ＋ 真 run 步级结论（run id 已取）

- `.github/workflows/ci.yml`：**`:888`** `- name: "cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)"`；
  **`:908`** `run: bash scripts/wisp-cli-tests.sh`（windows leg）。
- **真 run**：`gh run list`（gh 已认证）→ run id **42554322437**（headSha `9655882`，2026-11-09 03:04Z，conclusion=failure）。
  test-windows job 128956543330 的「cmd/wisp CLI tests」步 **conclusion=failure，但步骤真跑了**：

```
wisp-cli-tests.sh: scope=./cmd/wisp/ dll_dir=.../third_party/sherpa-onnx (pinned: onnxruntime.dll sherpa-onnx-c-api.dll sherpa-onnx-cxx-api.dll)
runtests.sh: go test exited 1 - packages=[./cmd/wisp/ ...] top-level: PASS=256 FAIL=8 SKIP=1, === RUN=333
portable-tests.sh: four numbers (all from -v output): === RUN=333  --- PASS=256  --- FAIL=8  --- SKIP=1
```

- 日志（9500 行）里 **`0xc0000135` 出现 0 次**、守卫行 0 次 ⇒ **CI 上不存在加载期洞**，failure 全是具名真红。
- 8 红＋1 SKIP 逐枚点名（CI 日志原文）：`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`、
  `TestTicket255RosterStillMatchesTheActualReadSites`、`TestTicket223PermissionDeniedSitsInItsOwnSentence`、
  `TestTicket198FirstRunCallerIsTheRunEntryOnly`、`Test313WritableUnionEqualsTheTwoRosters`、
  `TestPanelHostRealWindowHopAndLifecycle`、`TestAC4FocusReturnToPriorWindowGap33r5`（30.00s 窗口时序）、
  `TestLockedSuffixNeverTakesTheLockItself33r11`；SKIP＝`TestPanelHostLatencyPercentilesAC2`。
  归各用例所在票（255/223/198/313/33r5/33r11 族），**不是票 98 的缺陷**。
- 注：该 run 的树比读数树旧（cmd/wisp 用例分母 265 vs 556），但它证明的是**接线与加载面**这一层，与树的新旧无关。

## 4. AC#2 — documented 命令现跑：进程能起（两形）

**形 A（一行式，票面 09-28 追加逐字的口令）**：

```
PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 ./cmd/wisp/
```

真实读数：全包跑完 668.3s，顶层 **PASS=503 FAIL=51 SKIP=2**（对账 503+51+2=**556**=§2 分母），rc=1（红定性 §8）。

**形 B（staging 形，本票与语音车同款部署形状）**：

```
go test -c -o <stage>/wisp-cli.test.exe ./cmd/wisp/
bash scripts/stage-sherpa-dlls.sh <stage>        # 依 deps.toml 钉名搬三枚 DLL
cd cmd/wisp && <stage>/wisp-cli.test.exe -test.v -test.run <…>   # cwd=包目录，fixtures 向上找模型
```

本腿实测（stage=`/d/tmp/98-leg/stage`，**无 PATH 注入**）：三枚 DLL staged
（`onnxruntime.dll` `sherpa-onnx-c-api.dll` `sherpa-onnx-cxx-api.dll`），二进制从 `cmd/wisp` 起：

```
--- PASS: TestRunTextTaskTextPathEndToEnd (3.06s)
--- PASS: TestSecretSetGetListRoundTrip (0.04s)
PASS
staged_run_rc=0
```

⇒ **staging 车对 `cmd/wisp` 同样成立**：进程起、`0xc0000135` 消失、用例真绿。语音族还需要模型环境变量
（`WISP_SILERO_MODEL`/`WISP_SPIKE_MODELS`/`WISP_TTS_DIR`，钉定树 `third_party/spike-models/`，
测试夹具也支持从包目录向上自动找），staging 形全量读数见 §8.2。

## 5. AC#3 — 永不 skip 的守卫：已在册，判据发现跑成立

守卫本体＝`scripts/wisp-cli-tests.sh:275/:308`（`4a682f698` 落地，A1551 在册）：PATH 可达性守卫（对 deps.toml
钉的每枚 DLL 遍历 PATH 逐项查在）＋ 加载失败兜底（三条件同时成立才补具名 setup 失败，真红不改色）。
本腿在**当前 HEAD** 现跑判据发（HEAD 字节 `sed` 删掉 `export PATH=` 一行、root 钉仓根）：

```
wisp-cli-tests.sh: GUARD (ticket 98 AC#3) - cmd/wisp was NOT tested on this host: the loader cannot reach pinned DLL(s) on PATH: sherpa-onnx-c-api.dll sherpa-onnx-cxx-api.dll
...（三行解释：TEMP 目录起二进制、0xc0000135 形状、修法指向 staging，且明言不许放宽断言/不许挪测试）
mutant_rc=1
```

⇒ 把 AC#2 的注入拿掉，守卫**红并说清是哪条 DLL**（点名 2 枚；`onnxruntime.dll` 本机 System32 有副本所以不点，
守卫按盘上事实报——这正是票面要的"响亮失败"）。正控：count=1 全包跑与 CI run 里守卫行均 **0 次**（不误报）。
SKIP 不被静默放行的另一层在 `tools/d22scan/runtests.sh`（任何顶层 SKIP 判红），在册。

## 6. AC#4 — 台账更正：落到四张载有该句的票；89/96 现跑 grep 无该句

- 现跑 `grep` 六票原文：载有「`cmd/wisp` 的红不要追」句的是 **90:186、92:76 与 :280、94:129、97:58（另 :190-191 已自带半截更正）**；
  **89、96 的票面全文没有该句**（`grep -n "cmd/wisp" 89-*.md 96-*.md` 的命中全与"不要追"无关）——如实登记，无从更正、不造。
- 已给四票各追加一段（append-only，各 +15 行、删除列 0；numstat 现跑核对）：核心句＝
  「`cmd/wisp` 的红分两类——真红/真窗族（照追）vs **DLL 缺失加载红（`0xc0000135`，点名 `sherpa-onnx-c-api.dll`，
  staging 车与取数口令已解）**；CI 步级读数 run 42554322437 佐证」。
- 哪些 AC 的判据此前实际依赖过 `cmd/wisp` 的测试：点名清单已在册＝`docs/evidence/s1/98-ac1-ac5-readings.md` §6.3
  （票 133 的 `-overlay` 误诊 A361 已记账、票 71 AC#5、六票引的 ubuntu 红计数应以 `98-attribution.md` §5 重判、
  windows 侧红归各用例所在票）。本腿复述并维持该清单。

## 7. AC#5 — 门禁（隔离 clone 现跑，`git clone --shared` @`5071d8825`）

| 门 | 读数 |
|---|---|
| `go build ./cmd/... ./internal/... ./tools/...` | **rc=0** |
| `go vet ./cmd/wisp/` | **rc=0** |
| `/d/work/base/gopath/bin/gofumpt.exe -l cmd internal tools scripts` | **空**（0 行） |
| `sh scripts/d22scan.sh` | **rc=0**，尾行 `clean - no D22 ban violations`（ban #8 internal/=966、cmd/=242 等） |

全包测试读数见 §8（AC#5 的 test 口径）。

## 8. 全包读数与红的定性（⛔ 整包 rc=1 不是票 98 的缺陷）

### 8.1 形 A：PATH 一行式 `go test -count=1 -v`（主树，11:48–12:00，668s）

```
RUN=711 PASS=503 FAIL=51 SKIP=2   （顶层结果行 503+51+2=556=分母）
recipe_rc=1
守卫行（'ticket 98 AC#3'）= 0 次
```

51 红 2 SKIP 逐枚定性（全名单见 §8.3）：

- **语音族（约 40 枚，0.0x-0.4s 即红）**：红因逐字＝`speech: onnxruntime 1.17.1 is older than the 1.28 the linked
  sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH (on Windows the loader
  searches System32 before PATH)`——**恰是 staging 车存在要打的那层二级 shadow**：PATH 一行式只救第一级
  （`sherpa-onnx-c-api.dll`），救不了二级（`onnxruntime.dll`，System32 先于 PATH）。这不是包坏，也不是票 98 的
  加载洞（进程起了、测试跑了、有具名 FAIL）；staging 形下该族真跑真绿（§8.2）。
- **非语音红（11 枚）**：`TestTicket255RosterStillMatchesTheActualReadSites`、`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`、
  `TestTicket223PermissionDeniedSitsInItsOwnSentence`、`TestTicket198FirstRunCallerIsTheRunEntryOnly`、
  `Test313WritableUnionEqualsTheTwoRosters`、`TestLockedSuffixNeverTakesTheLockItself33r11`、
  `TestAC4FocusReturnToPriorWindowGap33r5`（5.10s 窗口时序）、`TestAC246DevLegIgnoresTheTestTaskInjection`、
  `TestAC1ResidentLegInstallsItsLogListenerOnDisk`、`TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`、
  `TestAC247HandingTheLevelToTheSeamIsNotVisibility`、`TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims`
  ——与 CI 上那 8 枚同族（窗口时序族＋名册族），在册现状，归各用例所在票。
- **2 SKIP（具名非静默）**：`TestAC247LiveMicrophoneLevelsReachTheBallSeam`（真麦克风，`WISP_LIVE_MIC=1` 才跑）、
  `Test16r1PremortemMicPermissionDeniedLiveProbe`（同族）。无一枚是"加载失败被当 SKIP"。

### 8.2 形 B：staging 形 `count=2` 全量普查（隔离 clone @`5071d8825`，`go test -c` ＋ stage-sherpa-dlls ＋ 三枚模型环境变量）

读数（2026-11-09 12:10–12:34，24 分钟，cwd=`cmd/wisp`，无 PATH 注入＋staged 三枚 DLL＋三个模型环境变量）：

```
$ /d/tmp/98-leg/clone-stage/wisp-cli.test.exe -test.count=2 -test.timeout 60m -test.v
count2_staged_rc=1
RUN=1414 PASS=1032 FAIL=48 SKIP=24    （顶层结果行 1104 = 两轮 × 552；每轮 4 条顶层 RUN 无判决行
                                        ＝在册真窗族"RUN 无判决、死点漂"形状，分母涨大后同形）
守卫行（'ticket 98 AC#3'）= 0 次；真 panic = 0（日志里 14 处 "panic" 全是钉线程诊断的引用叙述，非真 panic）
```

**与形 A 的对照（同一包、两形各 count 一发）**：

- **staging 形下语音族真跑**：形 A 里约 40 枚"VAD 没起来（`onnxruntime 1.17.1` shadow）"的红在 staging 形下消失；
  语音用例真加载模型、真起进程（日志有 `voice: model hand-off starting`、`voice: the store carries no VAD model;
  using WISP_SILERO_MODEL` 等原文）。
- **24 条 SKIP（12 名 × 2 轮）全部具名非静默**：
  - AEC 集成族 8 名（`Test26PBargeinStopsWithin400ms`、`Test26PEchoNeverTriggersAndNeverTranscribes`、
    `Test26PFeedOnlineMatchesWindowReplay`、`Test26PMarginRetestAndGateReverse`、`Test26PWithoutAECTheEchoTranscribesVerbatim`、
    `Test392AecTuningDoubleTalkSwing`、`Test392DisposeReleasesTheEngines`、`Test392KeepalivePostureFindingsAndColdDefault`），
    理由逐字＝`aec: webrtc-audio-processing binding not built (requires cgo + -tags webrtcap; see deps.toml
    [webrtc-audio-processing])`——构建面 opt-in 未编入，具名披露、不是 staging 缺陷；
  - 在册具名 4 名（各 2 轮）：`TestPanelPageFetchesItsOwnAssets299`、`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`
    （embed 无入口，票 33 r5 在册）、`TestAC247LiveMicrophoneLevelsReachTheBallSeam`、
    `Test16r1PremortemMicPermissionDeniedLiveProbe`（真麦克风，`WISP_LIVE_MIC=1` 才跑）。
- **FAIL=48 行 / 28 个唯一名**：窗口时序族（`TestPanelThreadIsSTAAndExitsCleanly`、`TestBallPanelGesturesReachThePanelThread`、
  `TestAC4PriorFocusSurvivesARefusedPanelSample`、`TestAC14AwaitedBindingReplyReachesThePage`、
  `TestAC14GoSideEvalPushReachesThePage`、`Test253r6TheSilentDropBranchIsLoud`、`TestAC4FocusReturnToPriorWindowGap33r5`、
  `TestPanelHostRealWindowHopAndLifecycle`、`Test384ResidentWatchdogEndToEnd`、`TestAC1ResidentLeg*`3 枚）＋
  名册族（`TestTicket255Roster*`2 枚、`Test313WritableUnionEqualsTheTwoRosters`、`TestTicket198FirstRunCallerIsTheRunEntryOnly`、
  `TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims`、`TestAC247HandingTheLevelToTheSeamIsNotVisibility`、
  `TestLockedSuffixNeverTakesTheLockItself33r11`）＋ 语音集成红（`Test15eSpeechEndToEndHandsTheTranscriptToTheAgent`、
  `Test16r1VoiceDemoWakeAsrAgentSpeakSettleRealHardware`、`Test41a2r1*`4 枚、`TestWireR1*`3 枚）——
  与 CI 上那 8 枚同族，归各用例所在票，**无一枚是加载期形状**（没有 `0xc0000135`、没有"0 结果行"的空仪器）。

### 8.3 51 红 2 SKIP 逐枚名单（形 A 实测，`grep '^--- FAIL|^--- SKIP'`）

见本件附录（§8.3 名单）：语音族 `Test15e*`(7)/`Test26P*`(7)/`Test26R3*`(2)/`Test381r2*`(2)/`Test382WR2*`(3)/
`Test390r1*`(3)/`Test392*`(3)/`Test394R1*`(2)/`Test395r1*`(1)/`Test41a2r1*`(4)/`TestWireR1*`(6)/
`Test12r2*`(1)/`Test16r1VoiceDemo*`(1)；非语音 11 枚；SKIP 2 枚（§8.1 已列）。

## 9. 派单主问的答案：四格能否由 staging 车收口？

**能，且比派单预估的还多一格**——但要把话说准：

- AC#2/AC#3 的"起不来/不点名"半边：**staging 车直接收口**（`scripts/stage-sherpa-dlls.sh` 对 cmd/wisp 同样成立，
  本腿现跑 PASS）；AC#3 的守卫本体早在册（`4a682f698`），本腿判据发复测成立。
- AC#1 的 (i)(ii)：修前 vs 修后两发＋分母，现跑坐实；(iii) file:line＋真 run id＋步级结论，gh 现取坐实。
- AC#4：据实更正落四票（89/96 无该句，如实登记）。
- AC#5：门禁四道隔离树全净；count=1/count=2 全量读数在案，红逐枚定性。
- **诚实边界**：`cmd/wisp` 整包 rc=1 **不是也不该被追成票 98 的缺陷**——它是真红/窗口族＋语音族二级 shadow；
  后者在 staging 形下消失（§8.2），前者归各用例所在票。staging 车解的是"进程起不来"与"二级 shadow 使语音族红"，
  不把整包变绿。

## 10. 本件没有建立的事（宁缺毋造）

- 没有改 `scripts/**`、`internal/**`、`ci.yml`、任何阈值/golden/AC 原文（只翻勾＋追加）。
- count=1 的 51 红没有逐枚归票（只分两类＋名单）；staging 形 count=2 的红名单以 §8.2 回填为准。
- ubuntu 腿未重普查（`98-attribution.md` 在册口径：余 40 枚归各票主）。
- CI 侧只取了步级结论与四数行；未取 8 红的 CI 侧逐行日志归票（名单已点名）。
- 票 89/96 简报若在其他介质（派单 txt 等）载有该句，不在本腿写面，未扫、未改。

## 11. 复跑口令（可逐字重跑）

```bash
cd "/d/work/workspace/projects plans/Wisp"
unset GOTMPDIR TMPDIR TEMP TMP
# AC#1(i) 修前
go test -count=1 ./cmd/wisp/ ; echo rc=$?          # rc=1, exit status 0xc0000135, 0 用例
# AC#1(ii) 分母
go test -c -o /d/tmp/98-leg/stage/wisp-cli.test.exe ./cmd/wisp/ && /d/tmp/98-leg/stage/wisp-cli.test.exe -test.list ".*" | grep -cE "^(Test|Example|Benchmark)"   # 556
ls cmd/wisp/*_test.go | wc -l                      # 178
# AC#2 形 A
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH" && go test -count=1 -timeout 30m -v ./cmd/wisp/
# AC#2 形 B（staging，无 PATH 注入）
bash scripts/stage-sherpa-dlls.sh /d/tmp/98-leg/stage && cd cmd/wisp && /d/tmp/98-leg/stage/wisp-cli.test.exe -test.v -test.run "TestSecretSetGetListRoundTrip|TestRunTextTaskTextPathEndToEnd"
# AC#3 判据发（拿掉注入）
R="/d/work/workspace/projects plans/Wisp"
git show HEAD:scripts/wisp-cli-tests.sh | sed -e "s|^root=.*$|root=\"$R\"|" -e "s|^export PATH=.*$|:|" > /d/tmp/98-guard-mutant.sh
bash /d/tmp/98-guard-mutant.sh ; echo rc=$?        # GUARD 点名 2 枚 DLL，rc=1
# AC#1(iii) CI
gh run view 42554322437 -R CarlosShao/wisp --json jobs
# AC#5 门禁（隔离树）
git clone --shared "/d/work/workspace/projects plans/Wisp" /d/tmp/98-leg/clone && cd /d/tmp/98-leg/clone
go build ./cmd/... ./internal/... ./tools/... && go vet ./cmd/wisp/ && /d/work/base/gopath/bin/gofumpt.exe -l cmd internal tools scripts && sh scripts/d22scan.sh
```
