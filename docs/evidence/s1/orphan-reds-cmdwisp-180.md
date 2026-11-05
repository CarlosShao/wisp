# 180 — cmd/wisp 孤儿红修复（98 归因 C1：DPAPI 桩夹具 → 平台无关的夹具缝）

**Status**：交付件（写腿自证）。**写面**：`cmd/wisp/**`（改 4 枚、新增 1 枚）＋ 本文件；**未 push**。
**锚定树**：`3ff1b1395`（派单给的 HEAD）。读数期间 HEAD 前进过两枚**别人的**提交——
`471365f92`（只动 `internal/{agent,proc}` 两枚生产文件＋两枚新测试）与 `9cd528199`（只动一份证据件）。
可核尺：`git diff --stat 3ff1b1395 9cd528199 -- cmd/wisp` = **0 行** ⇒ 本件引的 `cmd/wisp` 读数对两枚提交都成立。
**本腿提交**：`ae8136261`（工作树带本腿改动时取的全部读数；见 §6 与 §8）。
**仪器**：Windows 侧 `go1.27.1 windows/amd64`（`PATH` 前置 `build;third_party/sherpa-onnx`）；
POSIX 侧 `docker.io/library/golang:1.27-bookworm`（容器内 `linux/amd64`、`CGO_ENABLED=1`），
快照一律 `git archive`（不建 worktree、不 checkout）。
**时间**：2026-11-05 13:0x–13:5x +08。

## 0. 一句话结论

`docs/evidence/s1/98-attribution.md` §2 **C1**（非 Windows 上 DPAPI 桩，76 枚）里，
落在 `cmd/wisp/**` 且由**已 done 票**遗留的那批红，本腿在 POSIX 腿上把 **39 枚转绿**：
直接点名的票 22 枚（223×7、232×1、233×3、101×5、224×3、315×3），
外加**共用同一夹具**的 17 枚（201×9、255×4、165r2×2、226×1、227×1）。
整包读数（POSIX，同一仪器、同一命令）：

| | PASS | FAIL | SKIP |
|---|---|---|---|
| 修前 | 149 | **79** | 0 |
| 修后 | 188 | **40** | 0 |

**39 = 79 − 40**，逐枚可点名（§3、§5）。**0 断言改动、0 契约改动、0 `//go:build windows`、0 SKIP**。
剩下 40 枚里 37 枚是**别的票**自己的夹具（本腿写面不动它们），另 3 枚是：
C4（票 131 的**故意红**）、C6（票 119 孤儿，见 §5 STOP）、C5（我这台仪器的假红，真仓带 `.git` 时 PASS）。

## 1. 归因说了什么，本腿取了哪一段

98 §2 的 C1 逐字：

> `internal/secret/protect_other.go:18` 的 `errDPAPIUnavailable` = `secret: DPAPI is only available on Windows; use an env: ref on this platform`。
> 测试夹具写 `api_key_ref = "dpapi:…"` 后真的 `secret.Store(...)` …于是**夹具那一行**就 `t.Fatal`。
> 桩是**故意 fail-closed**（同文件 `:7-16` 注释；C28 / 票 06 是 Windows-only 设计），不是 bug。

98 §2 同时把**修法**留给了"用例所在的票"，并点名两条路（§7 末条）：
① 加 `//go:build windows` / 改名 `_windows_test.go`（票 70 / 131 的先例）；
② 把夹具换 `env:` ref（C28 的 portable 形状）。

**本腿取 ②，且只取 ② 的"缝"那一半**：不是把 76 枚一次性改掉，而是先把**夹具本身**改成
"按平台选引用形状"，于是所有共用它的用例同时离开红，而每一枚用例的断言一个字都没动。

**为什么 ① 不能取**：`98` 的 **AC#2** 逐字禁止——"不许用「把测试搬出 cmd/wisp」或「给整包加 build tag」
来变绿——那是把被检对象从门禁里删掉（票 78 的代理当年就是按这条拒绝过我，那次它是对的）"。
`scripts/portable-tests.sh --scope=cli` 会把 `cmd/wisp` 的每一条用例跑进 `tools/d22scan/runtests.sh`，
而那里 SKIP 判红（票 71 AC#3）。所以本腿的修法必须让用例**留在 CLI 分母里**并在 ubuntu 上真跑。

## 2. 改了什么（三处夹具 + 一处附带 + 一枚新缝）

新增 `cmd/wisp/secret_fixture_180_test.go`（唯一的新文件，`_test.go`，不进产品二进制）：

```
fixtureKeyRef180(t, dataDir, blobID, key) string
  windows  → 与今天逐字节相同：secret.NewStore(dataDir) + st.Store("dpapi:"+blobID, key)，
             返回 "dpapi:"+blobID
  其余平台 → t.Setenv("WISP_TEST_FIXTURE_KEY", key)，返回 "env:WISP_TEST_FIXTURE_KEY"
```

三处夹具改成"先用缝拿到 ref，再把 ref 写进 config"，DPAPI 的 Store 块从夹具里搬进缝（Windows 分支）：

| 文件 | 夹具 | 被它扶的票（98 §5 点名） |
|---|---|---|
| `cmd/wisp/approval_reply_201_test.go` | `newReplyHost` | 201（9）、223（8，含 232）、233（3）、224（3）、165r2（2）、226（1）、227（1）、255（4） |
| `cmd/wisp/run_mode101_test.go` | `t101boot` / `t101host.keyRef` | 101（5） |
| `cmd/wisp/zz315a1probe_test.go` | `zz315WriteConfig` | 315（3） |

附带一处（**必须同批改，否则本腿会造出一枚新红**）：`cmd/wisp/config_receipt_255_test.go` 的
`TestTicket255ReceiptStillNamesTheLiveReadSection` 原先**把 Windows 夹具的 ref 写死**成
`api_key_ref = "dpapi:acme"`，并带一条"premise moved"前提检查（`config_receipt_255_test.go:796` 起）。
缝一改，它在 POSIX 上就红成 `this case's premise moved: the boot snapshot carries ref "env:…"`。
改法：前提值改读 `r.h.keyRef`（`replyHost` 新增该字段，由缝填），plant 的新 ref 由
`r.h.keyRef + "-255r2"` 派生。**三条断言原样保留**（boot 快照不是被种的值 / 活的 settings 腿必须答被种的值 /
两者必须不同），只是把"写死 Windows 的 ref"换成"读本机夹具的 ref"。

合计：**改 4 枚、新增 1 枚**。`git diff --numstat` 里删除的每一行都在上面四枚的夹具/前提行内，
没有第二处语义改动（§4 给可核尺）。

## 3. RED before / GREEN after（原文摘录）

**RED before**（POSIX，HEAD `3ff1b1395`，无 overlay，整包；79 枚里 84 行同一句根因）：

```
--- FAIL: TestTicket223RunArmsTheReloadTick (1.39s)
    config_reload_223_test.go:371: secret: DPAPI is only available on Windows; use an env: ref on this platform
--- FAIL: Test233HotTierAgentReachesTheNextTask (1.28s)
    config_tiers_233_test.go:120: secret: DPAPI is only available on Windows; use an env: ref on this platform
--- FAIL: TestTicket101ManualSwitchSurvivesRestart (1.11s)
    run_mode101_test.go:216: secret: DPAPI is only available on Windows; use an env: ref on this platform
--- FAIL: TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun (0.72s)
    ticket224_assembly_test.go:81: secret: DPAPI is only available on Windows; use an env: ref on this platform
--- FAIL: TestZZ315A1RootLoopFirstBridgeCall (0.01s)
    zz315a1probe_test.go:476: secret: DPAPI is only available on Windows; use an env: ref on this platform
```

（日志点的行号是**用例的调用点**——夹具函数都带 `t.Helper()`，t.Fatal 归到调用它的那一行，
不是 `secret.Store` 自己那一行。这一点与本腿的改动无关，记下来省下一位的疑惑。）

**GREEN after**（同一仪器、同一命令、加本腿 overlay）：

```
--- PASS: TestTicket223RunArmsTheReloadTick (1.71s)
--- PASS: TestTicket223HandEditedFsLooseningCostsAnL2Card (1.79s)
--- PASS: TestTicket223RefusedLooseningKeepsOldValues (4.74s)
--- PASS: TestTicket223TighteningRaisesNoCard (1.71s)
--- PASS: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (1.78s)
--- PASS: TestTicket223RestartTierSaysItWillNotApply (1.77s)
--- PASS: TestTicket223FailureSentencesAreDistinct (6.88s)
--- PASS: TestTicket232R1ASeededStartupBannerCannotCarryTheRestartTier (3.47s)
--- PASS: Test233ReloadTierEventReachesAProductionConsumer (1.74s)
--- PASS: Test233HotTierAgentReachesTheNextTask (1.77s)
--- PASS: Test233ReloadTierReceiptOnlyForReloadTier (1.73s)
--- PASS: TestTicket101ManualSwitchSurvivesRestart (40.75s)
--- PASS: TestTicket101UntouchedConfigRestartsAtDefault (6.79s)
--- PASS: TestTicket101SessionGrantDoesNotCrossRestart (80.77s)
--- PASS: TestTicket101ModeSwitchUsesTheRealL2Gate (1.03s)
--- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict (0.75s)
--- PASS: TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun (0.75s)
--- PASS: TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking (2.69s)
--- PASS: TestTicket224ProductionSessionDoesNotSurviveRestart (2.73s)
--- PASS: TestZZ315A1RootLoopFirstBridgeCall (0.12s)
--- PASS: TestZZ315A1SiblingChildFirstBridgeCall (0.12s)
--- PASS: TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask (0.10s)
--- PASS: TestTicket255ReceiptStillNamesTheLiveReadSection (1.71s)
```

整包：`PASS=188 FAIL=40 SKIP=0`（修前 `PASS=149 FAIL=79 SKIP=0`）。**SKIP 两发都是 0**。

## 4. 这为什么不叫"放宽/摘尺"（五条可核的尺）

1. **没有删用例、没有加 build tag、没有改名**：`git diff --name-status` 只有
   `M approval_reply_201_test.go / M config_receipt_255_test.go / M run_mode101_test.go /
   M zz315a1probe_test.go / A secret_fixture_180_test.go`；`_windows_test.go` 一枚没动。
2. **没有改任何断言的判据**：三处夹具只改"ref 从哪来"，255 那条只改"前提值从哪来"。
   前后对照：`git diff` 里没有一行落在 `t.Errorf/t.Fatalf` 的判据文本上（除了 255 的派生表达式，
   它把常量换成同一常量的本机来源）。删 SKIP/改超时/放宽阈值一律没有。
3. **键值没换**：三处夹具原本就写同一枚假 key（`run_test.go:53` 的 `fakeStoreKey`）。
   `sha256` 对过（不打印明文）：`run_test.go:53` / `approval_reply_201_test.go` 原内联字面量 /
   `run_mode101_test.go` 的 `t101Key` 三者前 16 位都是 `99bcf0791f276080` ⇒ 缝里传的值与旧值逐字节相同。
4. **Windows 侧零变化**：缝的 Windows 分支与旧代码同序同参（`NewStore` → `Store("dpapi:"+blobID, key)`），
   config 里写回的也是同一枚 `dpapi:` 串。所以 Windows 的 DPAPI 覆盖**没有减少**，
   变绿只发生在"DPAPI 根本不存在"的平台上（那里今天只有 fail-closed 桩，无覆盖可减）。
   （唯一的例外是把 255 那条的 plant 值改成"由本机 ref 派生"，于是 Windows 上它种的 blob 名
   从常量 `co72jbmr` 变成派生名 `acme-255r2`——两个都只在内存里被读，从不解析、从不落盘；
   三条断言与判据一字未动。）
5. **用例仍留在 CLI 分母里**：测试名单一枚没增没删没改名（§4 第 1 条的 name-status 就是全部改动）；
   POSIX 整包的 `=== RUN` 从 306 到 309，多出的 3 行是
   `TestTicket101UnreadableModeFailsLoudlyAndStrict` 的三个子用例——修前父用例在夹具就死了，
   子用例根本没开始；修后它们真的跑了（`--- PASS …/存储损坏`、`…/版本不认识`、`…/权限读不到`）。
   两发 `SKIP=0`；project 的 `scripts/portable-tests.sh --scope=cli` 走的就是这个包。

## 5. 没碰的、以及 STOP 的（按派单口令）

- **STOP：C6 / 票 119**（`cmd/wisp/secret_dataroot_119b_test.go:241-242`，
  `TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119`）。修它要动
  `internal/risk`（冻结面，`pathresolver.go:36` 的 `ErrReparseDenied`）或替票 119 **裁定**
  "POSIX 上该认哪枚哨兵"——那是设计裁定，不是夹具缝。HEAD 上它已由编排者立成
  `.scratch/wisp/issues/406-c6-posix-sentinel-mismatch-secret-dataroot-119b.md`（open，Blocked by 三选一）。
  **本腿原样保留它红**：POSIX 整包 40 枚红里它是 1 枚，日志逐字
  `AC#3: refusal did not name ErrUnresolvedPath: … risk: path traverses a reparse point …`。
- **C4 / 票 131**（`TestAC4EveryLegIsNailedOrRuled`）：98 §2 C4 记的**故意红**，
  "a gate that cannot see is not a gate that has seen nothing to complain about"，
  Windows 上 PASS。不是缺陷，本腿不动。
- **C5 / 票 33**（`TestCleanCheckoutBuilds_AC11`）：我这条 `git archive` 快照没有 `.git` ⇒ 假红
  （98 §2 C5 已在真仓复跑 PASS）。**本腿不主张它是红的**，整包 40 里扣掉它 ⇒ 真仓口径 39。
- **余下 37 枚 C1 红**是**别的票自己的夹具**，本腿写面内可以改，但按派单"1–3 枚/最高价值"没有扩：
  63×8（secret CLI 自己的 Store 往返）、12×6（`run_test.go` 的 `newRunFixture`）、
  197×5、35×3、248×2、200×2、151×2、370×2、365×2、174r8×2、220×1、329×1、378×1。
  它们的夹具与三处已改的**不共用**（`newRunFixture` / 各票自建），同一种改法可以照抄本腿的缝。
  ⇒ 这是下一波的现成名册，不是本腿的欠账（本腿自选的 22 枚直接点名全绿）。

## 6. 门禁读数（逐条命令 + 读数）

全部在本机、工作树带本腿改动时跑：

```
$ go build ./cmd/... ./internal/... ./tools/...
rc=0

$ go test ./cmd/wisp -run 'TestTicket223|TestTicket232|Test233|TestTicket101|TestTicket224|TestZZ315A1|TestTicket255|DispatchGate' \
    -count=1 -timeout 240s
ok  github.com/CarlosShao/wisp/cmd/wisp  220.902s      # go 报的包耗时；外部墙钟 232.2s

$ /d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/
（空）

$ sh scripts/d22scan.sh
rc=0   # d22scan: clean - no D22 ban violations

$ go vet ./cmd/wisp/
rc=0
```

**一条诚实的话**：上面那条 gate 的 240s 上限只剩 **19s 余量**（220.9/240）。耗时大头是票 101 的两条
（`SessionGrantDoesNotCrossRestart` 80.8s、`ManualSwitchSurvivesRestart` 40.8s，都在等真实的 L1/L2 窗口），
与本次改动无关；我同一命令先后跑了两发都绿（239.7s、220.9s），但**机器一忙它就会顶到上限**。
若编排者要更稳的口径，建议 gate 用 `-timeout 420s`（判据不放松，只是把"机器负载"从判据里摘出去）。

## 7. 复跑口令（可逐字重跑）

```bash
cd "/d/work/workspace/projects plans/Wisp"

# A. Windows 侧（DPAPI 路径 —— 本腿改动的"零变化"尺）
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"
go test ./cmd/wisp -run 'TestTicket223|TestTicket232|Test233|TestTicket101|TestTicket224|TestZZ315A1|TestTicket255|DispatchGate' \
  -count=1 -timeout 240s

# B. POSIX 侧（本腿修的对象；docker 与 98 §6 同一只，快照 git archive + overlay）
bash .scratch/leg180/linux-red.sh 'TestTicket223|TestTicket232|Test233|TestTicket101|TestTicket224|TestZZ315A1|TestTicket255'
bash .scratch/leg180/linux-full.sh          # 整包：期望 PASS=188 FAIL=40 SKIP=0

# C. 门禁
go build ./cmd/... ./internal/... ./tools/...
/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/
sh scripts/d22scan.sh
```

（`.scratch/leg180/` 下的脚本与两发原始日志：`linux-red.sh`、`linux-full.sh`、
`linux-full-before.sh`、`gates.sh`、`surgery.py`、`surgery2.py`、`/d/tmp/leg180/{run,full,full-before}.log`。
临时件只建不删。）

## 8. 本件没有建立的事（宁缺毋造）

- **没有跑 Windows 整包**，只跑了 gate 那一组（220.9s ok）＋ 更早一发更大的子集（263s ok）。
  所以"Windows 整包全绿"这句本件**不主张**；主张的是"本腿改动的 22 枚 Windows 上仍绿"。
- **没有取 CI 读数**：全部读数是本机（Windows）与容器（POSIX），**不得**引成 CI 结论。
- **没有在 POSIX 上验 `-count=2`**：本件是普查/修复口径，不是门禁口径（98 §7 同样留白）。
- **POSIX 整包的 40 枚红里，38 枚未逐枚验证"更深处的判据"**：它们在第一行就死在夹具（C1），
  后面还会不会红本件量不到——与 98 §0 的边界①同一条，本件只是把**本腿负责的 22 枚**推到了真绿。
- **`internal/secret` 的两个形状没有改**：`protect_other.go` 的 fail-closed 桩、`refs.go` 的两种 ref
  都是冻结契约（C28 / C27），本件只消费它们。
- **`HEAD` 在读数期间动过两枚**（`471365f92` internal-only、`9cd528199` docs-only）：修前整包那一发
  （13:31 起）与修后整包那一发（13:41 起）都在 `9cd528199` 上，**同树对照**；
  修前"点名集"那一发（13:10）在 `3ff1b1395` 上。`cmd/wisp` 在这三枚提交之间逐字节相同（上面那把尺），
  所以两张读数可以并排引。同轮还有一条 internal 腿（`471365f92`，证据
  `docs/evidence/s1/orphan-reds-internal-180.md`）；两条腿的写面不重叠。
