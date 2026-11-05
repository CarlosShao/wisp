# 71 — `cmd/balldebug -diff` 阈值门（AC#1 / AC#2 的 balldebug 半张票）

- **票**：`.scratch/wisp/issues/71-gates-must-self-report.md`（What to build：「空仪器」第 1 条；AC#1、AC#2）
- **实现面**：ticket-71 实现腿（本腿被写面限定在 `cmd/balldebug/**` + 本文件）
- **开工 HEAD**：`2ccfda017`（`git rev-parse --short HEAD`，2026-11-05 09:3x）
- **改后（本 commit）**：`cmd/balldebug/{diffgate_windows.go,diffgate_windows_test.go}` 新增，`{diff_windows.go,main.go}` 各 +15 / +23 行
- **本文件是**：实现者自账（README 规则 6 的 1:1 表）。**不是**对抗验收；AC#7 要求验收者 ≠ 实现者，仍归编排者派。

---

## 0. 前提核对（开工三步，逐条现跑）

**(a) 票面点名的字符串/行还在 HEAD 上吗** —— 在，且语义未变（改动前读数）：

| 票面引用 | 改动前现读 | 结论 |
|---|---|---|
| `cmd/balldebug/diff_windows.go:365`「打印行」 | `sed -n 365p` = `fmt.Printf("diff: %-16s cpu1c=%6.3f%% … px>=3:%d/px>=8:%d/px>=24:%d of %d …%s\n",` | 命中，是那行 |
| `cmd/balldebug/diff_windows.go:495` `fileHas(..., "timers=true")` | `sed -n 495p` = `if fileHas(statusPath, "timers=true") {` | 命中，是那行 |
| 「`runDiff` 无论像素计数为何都返回 nil」 | `runDiff` 的 `return` 只有：mkdir / `os.Executable` / 无 `-diff-states` / 子进程或抓图出错；`writeDiffTable` 与三行 `fmt.Printf` 之后是 `return nil`，**表里的数字一次都没参与返回值** | 成立 |

**(b) 引用的源对吗** —— 对，且票面举的「`px_delta_ge8=0` 的跑」在盘上可指认：
`docs/evidence/s1/62-diff-baseline/diff-table.txt` = `Sleeping 0` / `Listening 1471` / `Speaking 4824`
（同目录 `62-diff-signoff` = `2120/4856/4848`，`62-diff-glass` = `2098/4855/4839`，`62-diff-border` = `2120`）。
⇒ 票面的判据（0 与 2000+ 在旧 `-diff` 眼里完全一样）由这两组真实表数据直接证实。

**(c) 红/绿不是别人在飞改动造成的**（同树现跑，未借他人读数）：

```
$ git status --porcelain -- cmd/balldebug        -> （空）
$ go build ./cmd/balldebug                       -> rc=0
$ go test ./cmd/balldebug -count=1                -> rc=0（包内当时 0 个测试文件，即「没有测试也算绿」的旧形状）
```

⇒ 前提未过期、未被别人占住、不在票外：**继续实现**。

---

## 1. 落地物

| 文件 | 作用 |
|---|---|
| `cmd/balldebug/diffgate_windows.go`（新） | 门的判据本体：`diffGateViolations`（纯函数）、`gateDiffRows`（打印「我看了几行」+ 违规明细 + 返回 error）、`gateDiffTable`（判一张已落盘的 `diff-table.txt`，不需要桌面）、`readDiffTable`（按**列名**回读，缺列即拒绝）、`resolveGateMinPx8`（阈值只许升不许降） |
| `cmd/balldebug/diffgate_windows_test.go`（新） | 8 条门测试（RED→GREEN 见 §4） |
| `cmd/balldebug/diff_windows.go` | `diffOpts` 加 `gate`/`gateMin`；`runDiff` 写完表后调用 `gateDiffRows`，**违规经 runDiff 的 error 路径返回** ⇒ `main` 里既有的 `os.Exit(1)` 生效 |
| `cmd/balldebug/main.go` | `-diff-gate`（判本次 `-diff` 写的表）、`-diff-gate-table <path>`（判已录表，无需桌面，CI 可调）、`-diff-gate-min-px8`（只许升）；阈值非法 ⇒ **rc=2**，门违规 ⇒ **rc=1** |

**为什么多一个 `-diff-gate-table`**：真 `-diff` 要交互桌面拍 alive/dead 两帧，CI 跑不了；没有它，AC#2 的阳性对照就只能靠人读输出，
正是本票要终结的形状。它让「门的红」变成一个退出码，并且能直接判仓里已录的证据表（§5）。

---

## 2. 阈值与它的来源（AC#1 的「写死 + 不许现场调低」）

```
diffGateFloorPx8 = 800
diffGateSource   = "SPEC-08 §2 (rest body 56*SleepRestRatio=34.72px, 1.5R halo => 2130 px predicted,
                    2098/2103/2120 recorded; SleepingRestMinPx=30 => 1590 px floor)"
```

推导（每一步都能在仓里复核，不是拍脑袋）：

1. `docs/specs/SPEC-08-ui-ball-panel.md` §2 定死 Sleeping 静止体径 = `SizePx × SleepRestRatio(0.62)`，默认 56 ⇒ **34.72px**；
   同节写明「稳定的量是 `px_delta_ge8`」，且最宽可见物是 `drawGlass` 的 `1.5R` halo ⇒ 预测 `π(1.5×34.72/2)² = 2130`；
   实测三次 = **2098 / 2103 / 2120**（`62-diff-glass` / `docs/SLO.md` A.2 / `62-diff-signoff`）。
2. 同节给出静止体径下限 `SleepingRestMinPx = 30` ⇒ 即使按**最小合法**体径，halo 面积也是 `π(1.5×30/2)² = 1590`。
3. 仓里所有被接受的行驶读数最低为 **1577**（小窗跑）/ **1620**（docked Sleeping）；被否证的那帧是 **0**；光标/提示框级别的噪声是个位数到几十像素。
4. 取 **800 = 1590/2**（最小合法体径 halo 的一半）：任何「真的成像」的跑至少以 2 倍余量过线，而「什么都没画」过不了。
   它不是质量线，是**「有没有成像」的线**——票 71 要挡的正是 0 被当成 2000。

**不许现场调低**是**机械强制**的，不是一句注释：`resolveGateMinPx8` 对 `< 800` 直接报错并 `os.Exit(2)`，把来源整句打出来：

```
$ balldebug -diff-gate-table docs/evidence/s1/62-diff-signoff/diff-table.txt -diff-gate-min-px8 10
balldebug: -diff-gate-min-px8=10 is below the documented floor 800, and it may only be raised: SPEC-08 §2 (…)
rc=2
```

---

## 3. 零定时器判据：读策略，不抄清单

`gateRequiresZeroTimers(label)` = `!ball.Animated(state) && ball.AnimationPolicy(state).Kind == ball.AnimNone`
（`-dock-<edge>` 后缀先剥掉，按状态判）。

**为什么不是一张手抄的「静态态」清单**：`docs/evidence/s1/12-ball-walk/diff-table.txt` 里 `Acting` 实测 `timers=yes`
（白名单允许它起定时器），`Warm`/`Settling` 实测 `timers=yes`（alpha / fade 路径）——
手抄清单会把这张**合法**的表判红（见 §4 的 `TestDiffGateZeroTimerOnlyWhereTheBallGrantsNone`）。
判据跟着 `internal/ball` 的单一权威走，策略变了判据自己跟着变。

---

## 4. RED → GREEN（先红后绿，逐字读数）

**RED（先写测试 + 一个「什么都返回 nil」的桩，正是本票的病）**

```
$ go test ./cmd/balldebug -count=1 -run 'TestDiffGate|TestDiffTable|TestReadDiffTable|TestGateFloor'
--- FAIL: TestDiffGateRedOnANonImagingRun (0.00s)
    diffgate_windows_test.go:31: gate passed a non-imaging run (px_delta_ge8=0); output="gate: 1 rows\n"
--- FAIL: TestDiffGateZeroTimerOnlyWhereTheBallGrantsNone (0.00s)
    diffgate_windows_test.go:80: gate passed timers=yes in Sleeping; output="gate: 1 rows\n"
--- FAIL: TestDiffGateMissingAndDuplicateRows (0.00s)
    diffgate_windows_test.go:101: gate passed a table missing the requested Listening row; output="gate: 1 rows\n"
--- FAIL: TestDiffTableRoundTripIsWhatTheGateReads (0.00s)
    diffgate_windows_test.go:133: read 0 rows, wrote 3
--- FAIL: TestReadDiffTableRefusesATableItCannotJudge (0.00s)
    diffgate_windows_test.go:151: readDiffTable accepted a table without a px_delta_ge8 column
--- FAIL: TestDiffGateTableRefusesAnEmptyTable (0.00s)
    diffgate_windows_test.go:164: gate accepted a table with zero rows; output=""
--- FAIL: TestGateFloorCannotBeLowered (0.00s)
    diffgate_windows_test.go:190: resolveGateMinPx8(0)=0, want 800
    diffgate_windows_test.go:186: resolveGateMinPx8(799) err=<nil>, wantErr=true
FAIL	github.com/CarlosShao/wisp/cmd/balldebug	0.040s   rc=1
```

**GREEN（换成真判据；同一命令）**

```
$ go test ./cmd/balldebug -count=1 -run 'TestDiffGate|TestDiffTable|TestReadDiffTable|TestGateFloor' -v
--- PASS: TestDiffGateRedOnANonImagingRun (0.00s)
--- PASS: TestDiffGateGreenOnRecordedGoodRuns (0.00s)
--- PASS: TestDiffGateZeroTimerOnlyWhereTheBallGrantsNone (0.00s)
--- PASS: TestDiffGateMissingAndDuplicateRows (0.00s)
--- PASS: TestDiffTableRoundTripIsWhatTheGateReads (0.00s)
--- PASS: TestReadDiffTableRefusesATableItCannotJudge (0.00s)
--- PASS: TestDiffGateTableRefusesAnEmptyTable (0.00s)
--- PASS: TestGateFloorCannotBeLowered (0.00s)
ok  	github.com/CarlosShao/wisp/cmd/balldebug	0.051s   rc=0
```

---

## 5. AC#2 阳性对照：**进程退出码**，不是人读输出

二进制：`go build -o /tmp/71/balldebug.exe ./cmd/balldebug`（`/tmp` = 仓外，守 A38：不在仓内造第二份源码）。
每条都是真进程，`rc=` 是 shell 拿到的 `$?`。

| # | 跑的是什么「故意不成像」的现场 | 命令（仓根） | rc | 关键输出 |
|---|---|---|---|---|
| 1 | **票 62 那帧真实的不成像跑**（Sleeping `px_delta_ge8=0`） | `balldebug.exe -diff-gate-table docs/evidence/s1/62-diff-baseline/diff-table.txt` | **1** | `gate: FAIL 01: Sleeping: px_delta_ge8=0 < 800 - the ball did not image …` / `gate: 1 violation(s) - this run is not evidence` |
| 2 | **阈值调到不可能达到**（AC#2 明列的第二种构造） | `… -diff-gate-table docs/evidence/s1/62-diff-signoff/diff-table.txt -diff-gate-min-px8 100000` | **1** | 6 行逐行 `FAIL … < 100000` |
| 3 | 想把阈值**调低**让红变绿 | `… -diff-gate-min-px8 10` | **2** | 打印 floor 来源整句并拒绝 |
| 4 | 空表（「扫了 0 行」） | `… -diff-gate-table /tmp/71/header-only.txt` | **1** | `has 0 rows - nothing examined is not a pass` |
| 5 | 缺列的表（判不了） | `… -diff-gate-table /tmp/71/no-px-total.txt` | **1** | `no "px_total" column - refusing to judge a table it cannot read` |
| 6 | 表不存在 | `… -diff-gate-table /tmp/71/nope.txt` | **1** | `read diff table: open … cannot find the file` |
| 7 | **阴性对照（不红）**：四张被接受的旧表 | `… -diff-gate-table <62-diff-signoff / 12-ball-walk / 62-diff-glass / 62-diff-border>` | **0 / 0 / 0 / 0** | `gate: examined 6/7/3/1 rows … gate: PASS` |

第 7 行是防「门变成噪声」的那一半：`12-ball-walk` 里 `Acting timers=yes` 必须**不**被判红（§3）。

---

## 6. 真跑：门接在 `-diff` 自己身上（本机桌面，仓外 `/tmp/71/shot*`）

**阴性对照（原型视觉，真的成像）**

```
$ /tmp/71/balldebug.exe -diff /tmp/71/shot -diff-states Sleeping \
    -diff-dwell 900ms -diff-sample 250ms -diff-gate -x 500 -y 400 -diff-margin 90
diff: Sleeping   cpu1c= 0.000% cpuAll= 0.000% privWS= 11.89MB commit=111.46MB ws= 40.38MB handles= 352 gdi=4 user=7 \
                 timers=no  px>=3:2156/px>=8:2112/px>=24:1240 of 78400 max=193 mean=1.528 box=(120,121 40x39) clean
diff: table -> …\shot\diff-table.txt
gate: examined 1 rows (of 1 requested), min px_delta_ge8=800
gate: floor source: SPEC-08 §2 (…)
gate: PASS - every examined row imaged and every zero-timer state held timers=no
rc=0
```

`px_delta_ge8=2112` 落在 SPEC-08 §2 记录的三次 2098/2103/2120 带内 —— 阈值 800 的余量与来源同时被这次真跑复核。

**阳性对照（同一命令 + `-frozen`：冻档的 Sleeping 帧物理上不成像）**

```
$ /tmp/71/balldebug.exe -diff /tmp/71/shot-frozen -diff-states Sleeping -frozen \
    -diff-dwell 900ms -diff-sample 250ms -diff-gate -x 500 -y 400 -diff-margin 90
diff: Sleeping   … timers=no  px>=3:0/px>=8:0/px>=24:0 of 78400 max=0 mean=0.000 box=(0,0 0x0) clean
gate: examined 1 rows (of 1 requested), min px_delta_ge8=800
gate: FAIL 01: Sleeping: px_delta_ge8=0 < 800 - the ball did not image (…)
gate: 1 violation(s) - this run is not evidence
balldebug: diff gate: 1 violation(s) at min px_delta_ge8=800
rc=1
```

两条命令只差 `-frozen`：**alive 与 dead 两帧逐像素相同（max=0）** ⇒ 这正是票 62 baseline 的同一形状，而旧 `-diff` 对两者的退出码都是 0。
（诚实备注：这次 `-frozen` 读数只是本门的阳性对照素材，不构成对 SPEC-08 §2.1 冻结视觉的任何裁决；那条属票 62/65/编排者。）

---

## 7. 门的自报（本票的主题）

- 绿的时候也不是哑的：`gate: examined N rows (of M requested), min px_delta_ge8=…` + `floor source: …`。
  `N < M` 就是「表短了」的第一现场（缺行另有逐条 `requested row missing`）。
- `-diff-gate-table` 先打 `gate: table=<path>`，再打行数 ⇒ 「判的是哪张表、看了几行」都在 stdout。
- 0 行、缺列、文件不在 ⇒ 都不是 PASS（§5 #4/#5/#6）。
- 表回读按**列名**而不是列序号（`readDiffTable`），列一旦漂移不会被静默读错；`writeDiffTable`→`readDiffTable` 有 round-trip 测试钉住。

---

## 8. 门禁读数（本腿收工前同树现跑）

```
$ go build ./cmd/balldebug                                  -> rc=0
$ go vet ./cmd/balldebug                                    -> rc=0
$ go test ./cmd/balldebug -count=1 -timeout 200s            -> ok  github.com/CarlosShao/wisp/cmd/balldebug  0.077s  rc=0
$ /d/work/base/gopath/bin/gofumpt.exe -l cmd/balldebug/diffgate_windows.go \
      cmd/balldebug/diffgate_windows_test.go cmd/balldebug/diff_windows.go cmd/balldebug/main.go
                                                            -> （空输出）rc=0
$ sh scripts/d22scan.sh                                     -> rc=0
   runtests.sh: OK - top-level PASS=35 FAIL=0 SKIP=0, === RUN=77
   d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=353, bans #1-5 cmd/=56,
            ban #6 frontend/=116, ban #7 internal/tools/=47, ban #8 design/=39 / frontend/=116 /
            internal/=865 / cmd/=203
```

---

## 9. AC 1:1 自账（实现者视角；验收仍须非实现者）

| AC | 框 | 真实读数 / 归属 |
|----|----|------------------|
| **#1** `-diff` 阈值门，带来源、不许现场调低 | ☑ 实现面 | 门：`px_delta_ge8 < 800`（§2 来源）／缺请求行／零定时器态 `timers != no` ⇒ `runDiff` 返回 error ⇒ `rc=1`；`-diff-gate-min-px8 < 800` ⇒ `rc=2`（§5 #3）。**框留给验收者**（本表是实现者自账） |
| **#2** 阳性对照（一次真红） | ☑ 实现面 | 进程级：票 62 真实 0 帧 ⇒ `rc=1`（§5 #1）；阈值不可达 ⇒ `rc=1`（#2）；**真 `-diff` + `-frozen` ⇒ `rc=1`**（§6），同命令原型 ⇒ `rc=0` |
| #3 测试运行器 | ☑ 已闭（前腿） | `tools/d22scan/runtests.sh`，见票面 Progress log 10:4x 条；本腿未动 |
| #4 d22scan 逐作用域自报 | ☐ 差裁决（前腿） | 前腿已落 8 条作用域台账 + `emptyScope` 守卫；ban #6 处置等编排者裁决。本腿未动 `tools/d22scan/**` |
| #5 接进 CI + 一次真实红 | ☐ 物理不可证（前腿） | 接入在 `.github/workflows/ci.yml`（前腿已落 3 处）；「贴出 CI 红」需 push，禁止。**本腿新增的可选项**：CI 无需桌面即可跑 `balldebug -diff-gate-table docs/evidence/s1/62-diff-signoff/diff-table.txt`（应 rc=0）+ `… 62-diff-baseline`（应 rc=1）作为门自身的阳性对照；`ci.yml` 不在本腿写面，未动 |
| #6 registry A30⑤ | ☐ 不在本腿写面 | `docs/reports/pending-and-issues.md` 不在本腿写面（写面只有 `cmd/balldebug/**` + 本文件），未动 |
| #7 对抗验收（非实现者） | ☐ 不归实现者 | 本文件是自账；建议验收命令＝§5 七行 + §4 两次 `go test` + §8 五条门禁 |

---

## 10. 本腿没做 / 交给谁（不夹带）

1. **`ci.yml` 接线**：不在写面。§9 #5 那一行是给有该写面的腿的现成命令（无需桌面）。
2. **`-frozen` Sleeping `px_delta_ge8=0`**：只作为阳性对照素材登记在此；是否算「冻结视觉的问题」由票 62/65 与编排者判，本腿不裁、不动 `internal/ball/**`。
3. **票 71 的 AC#4 ban #6 裁决、AC#5 CI 真红、AC#6 registry**：仍按票面 Progress log 留给编排者/有相应写面的腿。
4. 本腿**未改**：`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、任何 golden、`.github/**`、`tools/d22scan/**`、票面文件本身。
