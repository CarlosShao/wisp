# 票 159 落地腿 159-r1 实现件（part 3/3）：门禁 ＋ 整包读数 ＋ 红名册

接 `159-r1-impl.md`（part 1）与 `159-r1-impl-part2.md`（sec 3 尺、sec 4 反向判据）。

---

## 5. 门禁（本腿现跑，串行）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| build | `go build ./...` | **0** | — |
| vet | `go vet ./cmd/wisp/` | **0** | — |
| d22scan | `sh scripts/d22scan.sh` | **0** | `cmd/ 147 Go files`（含新文件）、`clean - no D22 ban violations` |
| path-length | `sh scripts/check-path-length-budget.sh` | **0** | `tracked paths=17549 over-budget=61 covered by roster=61 not in roster=0`、`VERDICT GREEN` |
| 整包 | `go test -count=1 -v -timeout 30m ./cmd/wisp/` | **1** | 按名作差，见 sec 6 |
| SLO 族显式 | `go test -count=1 -v -timeout 10m -run '^TestSLO' ./cmd/wisp/` | **0** | `RUN=39 PASS=20 FAIL=0 SKIP=0` |

格式：`gofmt -l` 与 `/d/work/base/gopath/bin/gofumpt.exe -l` 对**两枚改/新文件**均**零输出**。

---

## 6. 整包读数与红名册（按名作差；每发登记死点）

**基线（本腿开工时，改动前，`@ff21530f`）**：

```
RUN=225  PASS=165  FAIL=1  SKIP=0
红名册（1 枚）：TestAC14GoSideEvalPushReachesThePage
死点（1 枚）：TestPanelThreadIsSTAAndExitsCleanly  <- 有 RUN 行、无判决行
包级：FAIL github.com/CarlosShao/wisp/cmd/wisp 318.363s
```

=> 与派单说的"既有红 1 枚 `TestAC14GoSideEvalPushReachesThePage`（台账 `:20471`）"**逐名吻合**；
"每发恰 1 枚 RUN 无判决行"我这一发也复现。

**改动后（同一把尺、同一台面、`-v`、`-count=1 -timeout 30m`）**：

```
RUN=230  PASS=170  FAIL=1  SKIP=0
红名册（1 枚）：TestAC14GoSideEvalPushReachesThePage
死点（1 枚）：TestBookedSnapshotReachesHostBitForBit329  <- 有 RUN 行、无判决行
包级：FAIL github.com/CarlosShao/wisp/cmd/wisp 306.746s
```

**按名作差**（`diff` 两发红名册）=> **rc=0、逐名一字不差**。**新增红 ＝ 0**。
RUN 230 对基线 225 的 **+5** 逐名为
`TestAC4PriorFocusSurvivesARefusedPanelSample`／`TestBallGestureWithoutPanelHostStillRecords`／
`TestBallPanelGesturesReachThePanelThread`／`TestPanelThreadIsSTAAndExitsCleanly`／
`TestPanelThreadNameIsNotInResidentRoster` —— **全是基线那发被死点吞掉的面板族**，
`only-in-PRE=0`（基线没有任何一枚在新发里消失）。=> +5 是**死点前移**带来的覆盖面差，NOT 我新增的用例。
**死点逐发登记**：基线 ＝ `TestPanelThreadIsSTAAndExitsCleanly`（`panel_resident_windows_test.go`）／
改动后 ＝ `TestBookedSnapshotReachesHostBitForBit329`（`panel_snapshot_deliver_329_test.go`）。
两发都 `panic:` 0 命中、`exit status` 行 0 命中、早于 `-timeout 30m`（318s／307s）=> `A1100` 同族的进程级截断。

### 6.2 NOTE 覆盖面事实：整包 `-v` **从不跑到任何 `TestSLO*`**（不说清，上面那发会被误读）

```
baseline-v.txt      SLO-runs=0   最后一条 RUN ＝ TestPanelThreadIsSTAAndExitsCleanly（panel_resident_…）
gate-cmdwisp-v.txt  SLO-runs=0   最后一条 RUN ＝ TestBookedSnapshotReachesHostBitForBit329（panel_snapshot_…）
baseline-full.txt（无 -v）         TestSLO 命中 ＝ 0
```

原因不是 `-skip`：Go 按**文件名**字母序跑，`panel_*` 在 `slo_*` **之前**，
进程在 `panel_*` 里被截断 => 三个 `slo_*` 文件（含本腿新增那枚）**一枚都没被跑到**。
=> 所以"整包绿"这句话**不能**用来证明本腿的尺在工作；那是**另一把尺**的事。
因此我另跑一发显式选择器（sec 5 末行）：`RUN=39 PASS=20 FAIL=0 SKIP=0`、`rc=0`，
其中 `--- PASS: TestSLO159Case14CommentKeepsItsTenseAndItsNumbers (0.00s)`，
另 19 枚 `TestSLO144/147/149/152/156` 全 PASS => 本腿的尺**在包内可被跑到、且绿**，
也**没有**把 SLO 族既有 19 枚弄红。

NOTE 这也是**新登记的覆盖面缺口**（NOT 本腿能修，按派单不扩面）：
`cmd/wisp` 整包在死点前停下时，`slo_*` 三个文件的所有用例**都是零覆盖**，
而它们恰是 `wisp slo` 这条命令的唯一 Go 侧判据。与票 156 `C4`（"`wisp slo` 今天不被任何
`go test` 执行"）**同因不同面**：那笔说端到端，这笔说"在整包这把尺里它排在死点之后"。

### 6.3 提交后的 sha256 三处（逐字相同）

```
(1) 提交版  git show HEAD:cmd/wisp/slo_report_144_windows_test.go | sha256sum
            a34d1a24a3f7ef2ed58307e2f0f5d497b2333957b261514d113a1ca3b8373df6
(2) 还原台  sha256sum /d/tmp/wisp/20261022/159r1/pristine-144.go
            a34d1a24a3f7ef2ed58307e2f0f5d497b2333957b261514d113a1ca3b8373df6
(3) 工作树  sha256sum cmd/wisp/slo_report_144_windows_test.go
            a34d1a24a3f7ef2ed58307e2f0f5d497b2333957b261514d113a1ca3b8373df6
```

三处**逐字相同** => 七发变异跑完后，被测文件回到的正是提交进去的那一份字节，
没有残留、也没有"还原成另一个版本"。

### 6.4 提交与逐枚 `git show --name-only`

```
33bd8defe88120ae5fa27724bb37b49876548d5a
票 159 落地腿 159-r1：AC#2 三处时态过期句就地追加更正，并补一把防再漂的尺

cmd/wisp/slo_comment_tense_159_windows_test.go
cmd/wisp/slo_report_144_windows_test.go
```

`--name-only` 的 `.go` 命中**恰 2 枚**、越界 **0 枚**（无 `internal/**`／`tools/d22scan/**`／
`go.mod`／`go.sum`／`.github/**`／`scripts/**`／`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／
`allowlist.txt`／`frontend/**`／`design/**`；`cmd/wisp` 生产码 **0 枚**）。
本提交**不含**证据件（证据件单独一枚提交），所以"分开 commit"这一条在本腿上成立：
码改与记录面各一枚，且各枚的 `--name-only` 不含对方。

NOTE **正控（同尺打一笔含 `.go` 的真提交）**：上面这一枚**本身就是**正控——
同一把 `git show --name-only | grep '\.go$'` 的尺，在这一枚上**响了两枚 `.go`**。
若这把尺是坏的（恒空），它不会在 `33bd8def` 上报出两枚来。

### 6.5 本腿**没有**做的事（按派单，NOT 扩面）

- NOT 翻票 159 任何框、NOT 改 `Status:`／`-done`、NOT 改台账（账笔由编排者代贴）。
- NOT 改 `cmd/wisp` 生产码（本次 `--name-only` 里 0 枚生产码）。
- NOT 改原句 `:742`–`:752` 一字（append-only）；NOT 改 Case 13 那四句读数。
- NOT 修 6.2 登记的那枚覆盖面缺口（整包在死点前停下 => `slo_*` 零覆盖）：那要动
  跑测试的方式（`scripts/**` 或 CI），属越界面，按派单只登记。
- NOT push。
