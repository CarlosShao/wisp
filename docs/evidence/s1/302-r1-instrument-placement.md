# 302-r1 — 票 302 `AC#1` 落地形交付（甲形：三枚页面回执判据搬进 winlive 档）

- 腿码 302-r1 · 日期 2026-11-11 · 依据＝派单存档 `.scratch/wisp/probes/orch/20261111/txt/302-r1-dispatch.txt`（编排者已裁甲形，本腿未自选形）
- 票面＝`.scratch/wisp/issues/302-page-reply-instruments-run-in-the-ci-cli-tier-and-can-never-get-a-reply.md`
- 起手锚点＝`e31937a90`（dev，工作树在 `cmd/`、`scripts/`、`docs/` 上干净；唯一我方脏面＝本腿探针目录）
- ⛔ 零 push；add/commit 两处显式 pathspec

## 0. 形选择理由（派单已裁甲形；这里是"为什么甲形最符合主流 harness"）

主流做法＝**按"运行所需世界"分档**。一枚判据的效力＝"给它一个能供出它所读输入的运行世界"：
需要真 WebView2 窗＋页面自己回话的用例属于 GUI 档，只能在有桌面的机器上求值。托管 runner
没有那个世界，把它放 `--scope=cli` ＝让判据在永远拿不到输入的环境里求值——那是死格，不是门。
本仓已有同构先例两级：`panel_host_windows_live_test.go`（33-r5，winlive tag 承载时序敏感退出条款）
与 speech 档的 staging carrier（票 409，同一原则上一级：**判据由能供出其前提的载具求值**）。
乙形（放宽断言／无回执改 pass）直接违 AC#1 硬约束①；丙形（把常红灯登记成已知噪声）把门关掉。
本腿选择**build tag（`//go:build windows && winlive`）**而不是 ledger `-skip` 行，理由写死：
`-skip` 是静默过滤器（`A811`/`A812` 已证：不产 `--- SKIP` 行），而 build tag 让**编译面自己分档**——
cli 档的 `-list` 名册里这三枚名**消失**（可 grep 验证），winlive 编译面里**仍在**（`go vet -tags winlive`
门＋`-list` 双凭据），"搬走"这件事对两把尺都可见，没有任何一处需要靠"名单上少一行"来默会。

## 1. 动了哪几枚文件（写面全录）

| 文件 | 动作 | 内容 |
|---|---|---|
| `cmd/wisp/panel_page_reply_302_winlive_test.go` | 新建 | `//go:build windows && winlive`；三枚用例＋其私有 harness，**逐字搬入**（§2 凭据） |
| `cmd/wisp/panel_resident_windows_test.go` | 删块＋一句注释修正 | 移出 273–477 / 917–980 / 994–1024 三块；删孤儿 import `"regexp"`；文件头「winlive owns only the timing-sensitive exit clause」一句改为如实陈述（本腿造成的过期，**非判据本体**） |
| `scripts/portable-tests.sh` | 只加注释 | cli 档分支内 12 行注释（:465–:473），写明放置规则；**零可执行行、零 pin、零 ledger 行变动** |
| `docs/evidence/s1/302-r1-instrument-placement.md` | 新建 | 本件 |
| `.scratch/wisp/probes/302/r1/**` | 新建 | 读数原字节（logs/）、前后快照（base/）、搬出块原件（blocks/）、AC#4 名册与 comm 凭据（boundary/）、终态验证器 |

`countRequests` **刻意未搬**：`panel_unbound_253r6_windows_test.go`（windows 档）引用它——搬走会当场弄断
cli/windows 档编译。它留在 `panel_resident_windows_test.go:712`（现量行号）。

## 2. 搬运逐字凭据（AC#1②「判据本体⛔ 改一字」）

判据＝编排者 §⑥ 口径：**搬前搬后那几块函数体的 `cmp` 逐字节对拉**。三块边界（pristine 行号）：
块 A＝273–477（`entryIDRe`／`entryIDProbes`／`probeJS`／`latestReport`／`reportCount`／`awaitEntryReport`／
`probeAttemptSettle`／`awaitSettledHandoverDocument`／`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`）；
块 B1＝917–980（AC#14 nail 1 注释块＋`ac14AwaitJS`＋`TestAC14AwaitedBindingReplyReachesThePage`）；
块 B2＝994–1024（AC#14 nail 2 注释块＋`TestAC14GoSideEvalPushReachesThePage`）。
对拉尺＝按内容锚点从新文件抽回同名块，与 pristine 快照 `.scratch/wisp/probes/302/r1/base/panel_resident_windows_test.go.pristine`
（＝HEAD `e31937a90` 的 git blob `b0b9ea52d99dc3ccb96fe60ca2852a42f35b68e5`，`git hash-object` 对拉过）逐行比对：

```
block A : pristine-lines=205 newfile-lines=205 byte-identical=True
block B1: pristine-lines=64  newfile-lines=64  byte-identical=True
block B2: pristine-lines=31  newfile-lines=31  byte-identical=True
```

判据句原样在场（新文件内 grep）：`no report "ac13-probe" from the page within 15s`、
`no report "ac14r-0" from the page within 15s`、`no report "ac14-push" from the page within 15s`——
"没有回执就⛔ 判"逐字保住，零断言放宽、零 `t.Skip` 注入。
唯一非逐字新增＝新文件头 40 行**新写**的 tag＋放置说明注释（新文件本来就要有的头，⛔ 判据本体）。
`gofmt -l cmd/wisp`＝空（两文件合规范）；`gofmt -w` 只修了 4 处块缝空行形状，⛔ 未触碰任何搬入行内字节。

## 3. AC#2 成对两发（同机同码，只差"三枚在不在 cli 编译面"）

环境＝本机 dev box（Windows，go1.27.1，Git bash 载体 `D:\work\soft\Git\bin\bash.exe`；
PATH 先铺 `third_party/sherpa-onnx`＋`build`——DLL 与 CI 的 wisp-cli-tests.sh 同配方），独占跑（无并行 go test）。

### 红发（撤豁免＝把三枚搬回 cli 编译面）

做法＝resident 文件**字节级还原为 HEAD 原件**（`Copy-Item` pristine 快照；`git hash-object` 复核＝
`b0b9ea52…` 与 `git rev-parse HEAD:…` 相等——**还原凭据是 blob 双 hash，不是"我以为还原了"**）。
件＝`probes/302/r1/logs/04-red-shot-cli-with-three-in-cli.txt`：

```
portable-tests.sh: four numbers (all from -v output): === RUN=413  --- PASS=320  --- FAIL=3  --- SKIP=1
--- FAIL: TestTicket198FirstRunCallerIsTheRunEntryOnly (0.18s)
--- FAIL: Test313WritableUnionEqualsTheTwoRosters (0.00s)
--- FAIL: TestLockedSuffixNeverTakesTheLockItself33r11 (0.05s)
--- SKIP: Test423PunctuationSwitchGatesTheStage (0.01s)
```

三枚在红发里**全部被重新求值且本机全绿**（逐字）：

```
--- PASS: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (0.93s)
--- PASS: TestAC14AwaitedBindingReplyReachesThePage (1.18s)
--- PASS: TestAC14GoSideEvalPushReachesThePage (1.45s)
panel_resident_windows_test.go:972: AC#14 nail 1 (reply hop), page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)
panel_resident_windows_test.go:1020: AC#14 nail 2 (Eval push hop), page's own words: title="PUSHED-33R5-OK"
```

**诚实披露（不把它说成"预期红"）**：派单预期的"目标档当场红"那一腿，在本机世界里**没有发生**——
那三枚的判据随世界变色：编排者 10-10 在本机量过三枚全红（票面 §② `G4`），本腿今天同机量到三枚全绿
（其间隔着的正是票 303 的回执修复）。**"这三枚在 cli 档恒红"的凭据锚在托管 runner 的归档原字节里**
（`probes/301/orch/logs/ci-after-cli-block.txt`，10-10 那发，逐字三句）：

```
panel_resident_windows_test.go:325: no report "ac13-probe" from the page within 15s (what DID arrive at the door: nothing at all). AC#13 cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
--- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (20.04s)
panel_resident_windows_test.go:821: no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all). AC#14's reply hop cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
--- FAIL: TestAC14AwaitedBindingReplyReachesThePage (20.04s)
panel_resident_windows_test.go:866: no report "ac14-push" from the page within 15s (what DID arrive at the door: nothing at all). AC#14's push hop cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions
--- FAIL: TestAC14GoSideEvalPushReachesThePage (20.03s)
```

（整份抽录＝`probes/302/r1/logs/08-archived-ci-three-reds-verbatim.txt`。）红发的"有牙"部分＝**撤豁免有效**：
三枚名回到 `=== RUN`（重新被求值，RUN 410→413、PASS 317→320，Δ 恰＝三枚），不是既不跑也不编的隐形。

**废发登记**：第一次红发（件被覆盖，故事留痕在此）与 d22scan 门禁**并行**跑了 go test——违反
`cmd/wisp` 重包"同一时刻只跑一个"纪律，且 `Test219AC2RevokeRoundTripByteForByte` 在 chan receive
挂 9 分钟后把全包拖过 600s 超时（RUN=11 即断）。该发作废，独占重跑即上面这发。

### 绿发（`cmp` 逐字节还原＝终态文件写回）

做法＝`Copy-Item` 两个 post302 快照回树上（`git hash-object` 对拉相等，件 `verify-final.ps1` 读数）。
件＝`probes/302/r1/logs/06-green-shot-cli-three-moved-out.txt`：

```
portable-tests.sh: four numbers (all from -v output): === RUN=410  --- PASS=317  --- FAIL=3  --- SKIP=1
--- FAIL: TestTicket198FirstRunCallerIsTheRunEntryOnly (0.21s)
--- FAIL: Test313WritableUnionEqualsTheTwoRosters (0.00s)
--- FAIL: TestLockedSuffixNeverTakesTheLockItself33r11 (0.03s)
--- SKIP: Test423PunctuationSwitchGatesTheStage (0.01s)
```

三枚名在绿发全文 **0 命中**（grep 计数＝0）——不是被 `-skip` 静默过滤（那会仍在名册里、仍被编译），
而是**整个不在 cli 编译面**；同时 winlive 编译面 638 枚名册含那三枚（§4）。
红/绿 FAIL 名集合 `comm -3`＝**空**（件 `boundary/comm-red-vs-green.txt`）＝两发红的是**同一族**别族红，
本腿零顺带、零静默蒸发。

那 3 枚 FAIL＋1 枚 SKIP 与本腿无关的凭据：两发逐字同集合（comm 空）＋同名同红句（抽录见上），且它们
走的是**红发＝HEAD 原码**——即它们在起手锚点 `e31937a90` 就红/跳，本腿一字未碰其射程
（`internal/`、产码、`313/198/33r11/423` 的测试文件全部零写面）。具体归口不盖章（裁决者≠实现者）。

## 4. AC#2 附加凭据（甲形必须）：那三枚仍进编译面

```
go vet -tags winlive ./cmd/wisp/                          → rc=0   （ci.yml:968 同尺；logs/02、logs/10）
go vet ./cmd/wisp/                                        → rc=0   （cli/windows 档编译面未断）
go test -tags winlive -list '.*' ./cmd/wisp/              → rc=0，名册 638 枚，含三枚名（logs/03、logs/12）
  TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe : YES
  TestAC14AwaitedBindingReplyReachesThePage          : YES
  TestAC14GoSideEvalPushReachesThePage               : YES
```

⇒ "进 CI 的分母"没有变成"既⛔ 跑也⛔ 编"：winlive vet 门每次 push 都编译这枚新文件，
`-list` 名册含三枚名（正控＝同文件里 `^func Test` 恰 3 枚函数）。

## 5. AC#3 门禁读数

| 门禁 | 读数 | 件 |
|---|---|---|
| `sh scripts/d22scan.sh` | rc=0（终态树，`clean - no D22 ban violations`） | `logs/09-d22scan-final.txt`（另有一发中间态 rc=0＝`logs/05`） |
| `bash -n scripts/portable-tests.sh` | rc=0（改注释前后各一发） | 会话记录；终发 rc=0 |
| `bash scripts/portable-tests.sh --scope=census` totals 行 | **逐字未变**（见下） | `logs/01`（基线）/`logs/07`（改后）/`logs/11`（终态） |

```
基线（本腿起手，改前）：portable-tests.sh: census totals: packages=42 with-zero-compiled-tests=5 claimed-by-no-scope=5 unclaimed-with-tests=0   rc=0
改后：                  portable-tests.sh: census totals: packages=42 with-zero-compiled-tests=5 claimed-by-no-scope=5 unclaimed-with-tests=0   rc=0
终态：                  portable-tests.sh: census totals: packages=42 with-zero-compiled-tests=5 claimed-by-no-scope=5 unclaimed-with-tests=0   rc=0
```

⚠ 票面 `AC#3` 引的 `packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7` 是 10-10 的腐烂快照
（speech/census 那批 claim 之后数已变）；本腿的义务＝相对**本次起手基线**逐字不变——三条读数逐字相同，成立。

## 6. AC#4 边界声明（名册差集，comm 具名，件 `boundary/`）

**本票收掉（搬进 winlive 档、从此 cli 档不求值）＝3 枚**：
`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`、`TestAC14AwaitedBindingReplyReachesThePage`、
`TestAC14GoSideEvalPushReachesThePage`（=红/绿两发的差集枚；comm 凭据见 §3）。

**⛔ 本票收、留回台账 A817＝`internal/risk` 12 枚**（`boundary/risk-12-left-at-A817.txt`，逐枚具名）：
`TestAListWinsWhereBothTablesHit`、`TestBListDefaultDenyAndOverride`、`TestCanonicalInputGainsNoSecondForm`、
`TestClassifyAnchorSpellingIsNotVerdict`、`TestPathResolverExtendedLengthPrefixAListDenied`、
`TestPathResolverShortNameAListDenied`、`TestPathResolverUNCAListDenied`、
`TestSyncFallbackNotDisarmableByWeakRoot`、`TestSyncFixtureFallbackAndMatch`、
`TestSyncSuspectFallbackIsComponentBounded`、`TestSyncSuspectFallbackWhenUndetectable`、
`TestSyncUnverifiedRootKeepsFallback`（出处＝`pathresolver_anchor_spelling_windows_test.go`、
`pathresolver_junction_windows_test.go`、`syncdirs_test.go`，全部 `internal/risk/**`，本腿零写面）。

三把 comm（Git bash `comm`，正控＝risk-12 对自身＝12/12，证明工具没瞎）：
① risk-12 × windows 档归档红名册（`probes/301/orch/logs/ci-windows-fail-names.txt`）＝**12/12 具名仍红**
（"仍红"的档＝windows 档，那 12 枚从来不在 cli 档——risk-12 × cli 绿发 FAIL 集合 comm＝0 交集）；
② 红/绿两发 FAIL 集合 comm -3＝空（§3）；
③ 结论句：**那三枚从此⛔ 在 CI（cli 档）恒红——它们不再被该档求值，判据在 winlive 档等有桌面的机器；
那 12 枚仍红（windows 档），归 A817，与本票两枚射程。**
另：本机 cli 档今天另有 3 枚别族红（`TestTicket198…`／`Test313WritableUnion…`／`TestLockedSuffix…33r11`）
＋1 枚具名 SKIP（`Test423PunctuationSwitchGatesTheStage`）——两发同集合，⛔ 本腿造成、⛔ 本腿射程，
具名在此交下一程归口（CI 托管 runner 上是否同红以 CI 那发为准，本腿无 CI 读数）。

## 7. 越界检查（AC#3：逐笔 `git show --name-only --format=`，非区间尺）

本腿恰一笔 commit（pathspec 见 commit message；`git show --name-only --format=` 逐笔量＝该笔名册）：
`cmd/wisp/panel_page_reply_302_winlive_test.go`、`cmd/wisp/panel_resident_windows_test.go`、
`scripts/portable-tests.sh`、`docs/evidence/s1/302-r1-instrument-placement.md`、
`.scratch/wisp/probes/302/r1/**`——全部在派单写面内（许动清单 ∪ 探针目录 ∪ 证据件落点）。
⛔ 未触碰：`docs/PLAN.md`、`docs/specs/**`、`internal/observe/thresholds.go`、golden、
`tools/d22scan/allowlist.txt`、`frontend/**`、`design/**`、`.github/workflows/ci.yml`、`internal/**`、
任何 ledger 行、任何判据本体。越界＝**0**。

## 8. 本腿操作失误登记（诚实高于好看）

1. **废发一次**：第一次红发与 d22scan 并行跑 go test（违反重包独占纪律），被测包又恰逢
   `Test219AC2RevokeRoundTripByteForByte` 挂死拖过 600s 超时——该发作废，独占重跑取数（§3）。
2. **一次假空读数被自抓**：跑 `comm` 的命令没铺 Git bash PATH，`bash -c` 落到 WSL bash exec 失败，
   空输出差点被读成"集合相等"；发现后全路径重跑并给 comm 加了"文件对自身"正控。
3. **一次字节级写坏被自抓**：用 `git show | Out-String | Set-Content` 往树上写 HEAD 原件时，
   文本管道把非 ASCII 字节弄歪（file hash ≠ blob hash）；改用 `Copy-Item` 字节级复制＋blob 双 hash 复核。
   该文件在本腿期间从未进入任何测试运行（发现于起跑红发之前），树上最终态经 `verify-final.ps1` 对拉相等。
4. **一句过期注释顺手修正并在此披露**：resident 文件头「winlive owns **only** the timing-sensitive exit
   clause」因本腿搬家而变假话；改为如实列举两族。属注释，⛔ 判据本体；若口径认为这越了"只许动 build tag 行"
   的界，撤销口令＝把该句改回原文（快照 `.scratch/wisp/probes/302/r1/base/panel_resident_windows_test.go.post302`
   保留两种状态）。
