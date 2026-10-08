# 04 · 本机色：出处、口径、以及"本腿复跑不了的那把尺"

## 本机色的唯一一发近期整包读数（本腿只引这一发，逐名现取）

件＝`.scratch/wisp/probes/35/r4/logs/orch/full-package-after-r4.txt`（**276 148 B**，`ls -la` 现量 rc=0；文件 mtime `10-08 08:28`，体内首行 `time=2026-10-08T08:18:41.407+08:00`、末行 `--- PASS: TestTicket224ProductionSessionDoesNotSurviveRestart (3.44s)` / `FAIL` / `runrc=1`）。
`cold bring-up` 那一栏自带时刻与锚：`HEAD bcfb459b at read time, 2026-10-08T08:20:57+08:00` ⇒ 本机读数是**带 `+08:00` 时间戳**的那一族；CI 的读数一律带 `Z`（本件据此区分，⛔ 不混）。

尺（本腿自己跑，纯 `grep`，⛔ 不是复跑 `go test`）：
`grep -hE "^--- (FAIL|PASS|SKIP): <名> " <件> | head -1` ⇒ 逐名色见 `03-*.md` 的对照表；`rc=0`（除三枚 winlive 档报 `NOT_IN_FILE`，那是**没编译进去**、不是没跑）。

⚠⚠ **本腿没有复跑产生这些色的那把尺**：那把尺是 `go test ./cmd/wisp/`（＋需要 sherpa/onnxruntime DLL 就位），而本腿的派单硬约束是⛔ `go build`/`go vet`/`go test`（编译与测试面归 `255-r1`／`35-r7`）。
⇒ **按第 108 条，这一列本机色只能算"它取过"、不能算"本腿复跑过"**。本腿只做了两件：①把那件里的**原文行**逐名重新取一遍（不是抄别的腿的表）②核对它的时间戳/锚与 CI 的读数不同形。**欠的那一发**：由持有测试面的腿或编排者在 `052b393f` 上重跑一次整包，逐名对上 `03-*.md` 那张表。

## 本机 `-1` 那形是什么时候开始的（盘上现量，`grep -rl 'got -1.000'` ＋ 逐件取 `:660` 那一行，rc=0）

| 件 mtime | 件 | `:660` 那栏（逐字） | 是哪一侧 |
|---|---|---|---|
| 10-02 08:36–09:22 | `probes/orchestrator/ci-delta-1/clean_new.txt`、`vm-draw-1/clean_*_push.txt`、`tsv_tw_push.txt` 等 5 件 | `-1.000 ms (HEAD 8ae4c23 at read time, 2026-10-01T16:09:20Z)` | **CI**（`Z` 时间戳，件是 CI 日志的存档拷贝） |
| 10-07 20:05 | `probes/35/r2/logs/cmd-wisp-suite-baseline.txt` | `-1.000 ms (HEAD 286a7f30 at read time, 2026-10-07T19:50:11+08:00)` | **本机** |
| 10-07 20:05 | `probes/35/r2/logs/cmd-wisp-suite-r2.txt` | `-1.000 ms (HEAD 286a7f30 at read time, 2026-10-07T19:58:05+08:00)` | **本机** |
| 10-07 21:25 | `probes/35/r3/logs/cmd-wisp-full.txt` | `-1.000 ms (HEAD 9b2551f2 at read time, 2026-10-07T21:21:46+08:00)` | **本机** |
| 10-08 08:28 | `probes/35/r4/logs/orch/full-package-after-r4.txt` | `-1.000 ms (HEAD bcfb459b at read time, 2026-10-08T08:20:57+08:00)` | **本机** |

⇒ `-1` **不是 CI 独有的形状，也不是本机独有的形状**：10-01 的 CI 那批发过、10-07 起本机连发四批。⚠ 本腿**没有**把 `probes/35/v6/live-wave.md:44` 那个 `-1` 计入任何一格（派单写明它是本机已知坏读数、⛔ 不许当凭据）。

## 本机冷启毫秒的分布（同一把尺：`grep -rh 'cold bring-up measured on this box'` 全 `.scratch/wisp/probes/`，按值聚合，rc=0）

本机那族（`+08:00`）：`983.297`／`980.227`／`977.297`／`967.450`／`962.998`／`960.457`／`937.083`／`934.588`／`932.019` ms（各 1 发，散布在 10-01→10-05）——**都在 D32 的 1500 ms 预算内**。
CI 那族（`Z`）：`3126.549`／`3191.700`／`3414.339`／`3428.032`／`3496.853`／`3743.192`／`3883.234`／`3941.081` ms，加 `-1.000` 若干。
⇒ **同一枚判据（`:665` 的 1500 ms 预算）在本机历史上过、在托管那台历史上不过**；两族读数**不可相减、不可互换**（不同机器、不同时刻）。⛔ 本件不据此提任何阈值改动。

## 本机名册里"CI 绿／本机红"那两枚的红句（逐字，取自上面同一件）

```
panel_resident_windows_test.go:821: no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all). AC#14's reply hop cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions        [TestAC14AwaitedBindingReplyReachesThePage, FAIL 20.02s]
panel_resident_windows_test.go:866: no report "ac14-push" from the page within 15s (what DID arrive at the door: nothing at all). ...                                                                            [TestAC14GoSideEvalPushReachesThePage, FAIL 20.01s]
panel_resident_windows_test.go:325: no report "ac13-probe" from the page within 15s (what DID arrive at the door: nothing at all). ... [TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe, FAIL 20.02s；同枚在 CI 是 SKIP]
panel_host_windows_test.go:948: the panel did not take the foreground on Show: foreground 0x800b2, panel hwnd 0x50f90 (AC#4 says only the panel takes focus when shown)   [TestAC4FocusReturnToPriorWindowGap33r5, FAIL 5.15s；同枚在 CI run305 是 PASS]
```

对照 CI run 305 同名（`TestAC14GoSideEvalPushReachesThePage`）：红在 `:869`、`title=""` ⇒ **页报回了、只是报的是空标题**；本机红在 `:866` ⇒ **页什么都没报回**。
⇒ **同枚用例两侧红在不同出口**，坏法不同（第 128 条那一族：红名册不区分坏法，必须连文案一起引）。

## 顺带一枚与窗无关的本机既有红（避免被算进这六枚）

`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`（`cmd/wisp/config_reload_223_test.go:535`）在 10-01 的整包件里就是红的，`33-v4` 已具名把它从窗口名册里摘出来（体扫命中 0）。本腿未重跑，只带出处。

rc=0
