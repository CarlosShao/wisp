# 133 续腿：`covered=test` 那一桶的平台分裂改成「具名披露、不判红」（`98-attribution` C3）

**票**：`.scratch/wisp/issues/133-no-instrument-sees-the-dispatch-hop-deleting-one-caller-keeps-76-tests-green-done.md`
（票面 `Status: ready-for-agent`，交付后 `-done`）。**本件不是重领票面**：票面 AC 已结、不动一字；
本腿是普查件 `docs/evidence/s1/98-attribution.md` §2 **C3** ＋ §3「按票点名」那一条（133 1 枚）按派单开的
**续腿**，做的是那一件已经点名、但票面没写进 AC 的事：把 `TestAC1AC2DispatchHopGate133` 里仍 `t.Errorf`
的那一支改成票 131 先例的「具名披露、不判红」。

**Status**：实现腿（写件）。**代码改动只在两枚 `_test.go` ＋ 本文件**；生产代码、票面、`docs/PLAN.md`、
`docs/specs/**`、`thresholds.go`、golden、`frontend/**`、`design/**` **一字未动**。
**锚点**：开工 `git rev-parse HEAD` = `05e02c4df`；本机 Windows `go1.27.1 windows/amd64`；linux 腿
`docker.io/library/golang:1.27-bookworm`（容器内 `go env GOOS GOARCH CGO_ENABLED GOVERSION` =
`linux / amd64 / 1 / go1.27.3`）。时间 2026-11-05 12:25–13:0x +08。

## 0. 一句话

`TestAC1AC2DispatchHopGate133` 在 ubuntu 腿上**只红一枚**（`slo` 腿被记为由
`Test384ResidentWatchdogEndToEnd` 覆盖，而该用例住在 `resident_watchdog_windows_test.go:184`，linux 二进制
没有它）——这就是 `98-attribution` §2 C3。本腿把这一支改成**具名披露**：行的 `covered` 列写
`SOURCE-ONLY on GOOS=…`、并点名 case 与 file；披露行点名 leg/case/file/编译它的那条腿；**不判红**。
边界守住：①本平台 build 收的文件里声明、却不在 roster（M-G 的原形，自证腿里那种无约束文件名）；②
哪条腿的 build 都不收（`//go:build ignore`）——两者**照旧红**，各有一条自证断言钉住。

## 1. 改了什么（两枚文件，逐处）

`cmd/wisp/leg_dispatch_gate_133_test.go`

| 处 | 之前 | 之后 |
|---|---|---|
| `pkg133` | 只留 `fset` | ＋`dir`（把 `go/build` 的判定指到被解析的目录） |
| `leg133` | `coveredBy` 一桶 | ＋`sourceOnly`（平台分裂的具名披露条目） |
| `drivenBy133` | 取排序后第一个驱动用例 | ＋`skip` 参数：跳过已归入平台分裂的名字，**继续找本腿真能跑的驱动**（披露不许顶替本腿真有的覆盖） |
| `coverageReds133` | 只分 `nail / installs / test / ruling / none` | ＋`case splitName`：行写 `test <case> drives <sym> - SOURCE-ONLY on GOOS=<this>: <file> is not in this leg's build; disclosed, not red (the case runs on GOOS=<other>)` |
| `platformSplit133`（新） | —— | 唯一新增的判定：**两个读数同时成立**才披露（roster 没有该 case；该名字的**每一条**声明都在本腿 build 之外、且至少一条被另一条 CLI 腿的 build 收） |
| `buildTakes133`（新） | —— | 平台答案问的是 `go/build` 的 `Context.MatchFile`（工具链自己的约束匹配器），**不是**本文件自写一个 tag 求值器；另一条腿用 `CgoEnabled=true` 探 |
| `runRosterReds135` | roster 缺该 case 即红 | 同一条 `platformSplit133` 成立则该行不判红（行与披露行已具名） |
| `runRosterDisclosure135` | 固定写死 `(one outside the roster is red above)` | 改成计数（`N distinct, M startable, K red above, J claim(s) stopped at the platform gate`）＋**逐条点名** `leg "<leg>" <site> -> <case> (<file>, taken on GOOS=<other>)` |
| 文件头 | 第三 limit 原文「`covered=test` 那一桶**没有**这种豁免」 | 改写为该桶的豁免规则、它的出处（98 C3）、以及守 M-G 的边界 |

`cmd/wisp/dispatch_hop_selfproof_133_test.go`

- 新增 `TestAC3SelfProof133ThePlatformGateDisclosesByNameInsteadOfReddening`：平台分裂夹具（用例只写在
  `eps_<另一条腿>_test.go`）＋**两条边界夹具**（无约束文件名；`//go:build ignore`），逐条断言：
  (a) `coverageReds133` 不为 `eps` 出红、行含 `SOURCE-ONLY on GOOS=` 与该文件名、`coveredBy` 为空；
  (b) `runRosterReds135` 在该 roster 下不出红；`runRosterDisclosure135` 点名 leg/case/file/`taken on GOOS=`；
  (c) **`//go:build ignore` 那一发仍红**（`has no such case` ＋ 点名 case）。
- `distinctnessReds133` 的 `driven` 分母从「`covered` 以 `test ` 开头」收紧为「`coveredBy` 非空」：
  `SOURCE-ONLY` 行也以 `test ` 开头，但它**不是**本腿 roster 能核的名字；不收紧的话，「所有覆盖都停在
  平台闸门」的树会让 roster 核对带一个空分母通过（Y3 形状）。这一改只会更红，不会更绿。

## 2. RED-before（真读数，不是构造）

**linux 腿（同一棵 `05e02c4df`，未打补丁）**：

```
$ CGO_ENABLED=1 go test -count=1 -v -timeout 30m -run 'TestAC1AC2DispatchHopGate133|TestAC3SelfProof133...|TestAC4Distinctness133...' ./cmd/wisp/
--- RUN   TestAC1AC2DispatchHopGate133
    leg_dispatch_gate_133_test.go:251: AC#1/#2 RED: leg "slo" (main.go:133) is booked in the ledger below as covered by the case "Test384ResidentWatchdogEndToEnd", and this round's test binary has no such case: the roster read from the running binary lists 226 startable cases and "Test384ResidentWatchdogEndToEnd" is not one of them, so no "=== RUN   Test384ResidentWatchdogEndToEnd" line exists in this run or can exist in it. The name is declared at resident_watchdog_windows_test.go:184, which is the point - the declaration is in the sources and the sources are not the build.
--- FAIL: TestAC1AC2DispatchHopGate133 (0.22s)
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.494s
TEST_RC=1
```

同发披露行（原文）：`GOOS=linux, 226 startable cases … covered=test: 5 distinct, 4 of them startable (one outside
the roster is red above) …`。**这就是 `98-attribution` §5 第 30 行那一枚**（`98` 记的是 224 条 roster，本机这一发
是 226，因为期间树前进了两枚别人的测试；红名、红行、红因三样一致）。

**windows 腿（同一棵，未打补丁）**：4 枚 133 用例 `--- PASS`，账本 `leg slo … covered=test
Test384ResidentWatchdogEndToEnd drives settleBeforeSample`；即 C3 只在 linux 腿上红，与 `98-attribution` §4 一致。

## 3. GREEN-after（同一发命令、同一容器、同一 snapshot）

做法：`git archive 05e02c4df` 的干净快照（`sha256=1338706ab78502a98f3fb3103570cdb828609936d5e93ebdbe5eb3549405983d`，
329,410,560 字节）解进容器本地盘，**只把改过的两枚 `_test.go` 覆盖进去**（不含任何未提交/未跟踪件）。

```
$ CGO_ENABLED=1 go vet ./cmd/wisp/          # rc=0（无输出）
$ CGO_ENABLED=1 go test -count=1 -v -timeout 30m -run '<4 枚 133 用例>' ./cmd/wisp/
--- PASS: TestAC3SelfProof133TheDispatchReadersCanAnswerNo (0.00s)
--- PASS: TestAC3SelfProof133ThePlatformGateDisclosesByNameInsteadOfReddening (0.00s)
--- PASS: TestAC4Distinctness133NoTwoLegsWearTheSameRow (0.37s)
          leg slo  main.go:133  installs=false handoff=false covered=test Test384ResidentWatchdogEndToEnd drives settleBeforeSample - SOURCE-ONLY on GOOS=linux: resident_watchdog_windows_test.go is not in this leg's build; disclosed, not red (the case runs on GOOS=windows) entries=cmdSLO aliases=
    blindness disclosure: run-roster disclosure: GOOS=linux, 227 startable cases read from this binary itself (`-test.list '.*'`); case names this round's ledger credited through `covered=test`: 4 distinct, 4 of them startable, 0 red above; 1 claim(s) stopped at the platform gate (named below, not reddened); this gate's registry: 4 claims, 4 of them outside this round's roster - leg "slo" main.go:133 -> Test384ResidentWatchdogEndToEnd (resident_watchdog_windows_test.go, taken on GOOS=windows). … Ticket 133's R-133-9 (M-G) is unaffected: a name whose declaring file this build DOES take, or that no leg's build takes, is red above. …
--- PASS: TestAC1AC2DispatchHopGate133 (0.40s)
PASS	github.com/CarlosShao/wisp/cmd/wisp	0.837s
TEST_RC=0
```

windows 腿（改后）：`go test -count=1 -timeout 200s -run 'TestAC1AC2DispatchHopGate133|TestAC3SelfProof133|TestAC4Distinctness133' ./cmd/wisp/`
⇒ 4 枚全 `--- PASS`，`ok … 0.698s`；账本与改前逐字相同（`slo` 在 windows roster 里，正常记 `covered=test`），
披露行 `… 5 distinct, 5 of them startable, 0 red above; 0 claim(s) stopped at the platform gate …`。

## 4. 边界：四种形状，谁披露谁红（自证腿逐条钉）

| 形状 | 本腿 build 收？ | 另一条腿 build 收？ | 结果 |
|---|---|---|---|
| `resident_watchdog_windows_test.go` 在 linux 上（C3 真形） | 否 | 是（windows） | **具名披露，不红** |
| 无约束 `fixture_test.go` 里声明、roster 没有（M-G 原形） | 是 | —— | **红**（自证腿 `TheDispatchReadersCanAnswerNo` 钉住） |
| `//go:build ignore`（哪条腿都不收） | 否 | 否 | **红**（新自证腿 (c) 钉住，点名 case） |
| 另一条腿的文件在**那条腿**上跑 | 是 | —— | 正常 `covered=test`（roster 里有） |

**对称性是故意的**，也是要请裁决方看的一处：M-G 那一发的真实方向是「linux 用例在 windows 轮次被记成覆盖」，
与 C3 的「windows 用例在 linux 轮次被记成覆盖」是**同一形**。所以豁免只能对称；本腿的做法是让**行**在
不编译它的那条腿上不再声称覆盖（`SOURCE-ONLY`，且不进 `coveredBy`），而不是把红改成静默绿。
票 131 对注册表钉子用的正是这个形状（「build 不收的钉子写出来、不判红」），本腿把它扩到 `covered=test` 桶。

## 5. 门禁读数（本机，逐条）

| 门 | 命令 | 读数 |
|---|---|---|
| 构建 | `go build ./cmd/... ./internal/...` | **rc=0** |
| 构建（全仓） | `go build ./...` | 被**未跟踪**的本地目录 `build/`（`build/mut_rvw.go`，`git ls-files build` = 0）挡住：`undefined: voiceBargeinMonitor` 等；**与本次改动无关**，不是本腿造的 |
| 定向测试 | `go test -count=1 -timeout 200s -run '<4 枚 133 用例>' ./cmd/wisp/` | **ok 0.698s**（windows）；linux 容器同一 `-timeout` 口径 **ok 0.837s** |
| 整包测试（派单口径） | `go test -count=1 -timeout 200s ./cmd/wisp/` | **超时**（`panic: test timed out after 200s`）：本机这棵树的整包要 **580s**（下面一条）；这是机器/套件规模，不是本腿 |
| 整包测试（放宽到 900s，只取名册） | `go test -count=1 -timeout 900s ./cmd/wisp/` | `FAIL … 580.339s`，27 枚红**全是语音/实时硬件类**；抽一枚的原因为 `speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH (on Windows the loader searches System32 before PATH)` ⇒ **本机 DLL 遮蔽**，非本腿 |
| 格式 | `gofumpt -l cmd/` | **空**（本腿两枚文件干净；HEAD blob 抽出来跑 `gofumpt -l cmd internal` 也是空） |
| 格式（仓级杂音，据实记） | `gofumpt -l` 全仓 | 唯一一枚 `internal/agent/approval/pending_read.go`：**工作树 CRLF 131 个 CR，HEAD blob 0 个 CR**（`git status` 干净）⇒ 检出/编辑器产物，非本腿、非树内容缺陷，未动 |
| D22 扫描 | `sh scripts/d22scan.sh` | **rc=0**，`d22scan: clean - no D22 ban violations`（扫 414 枚生产 Go、`cmd/` 205 枚含 `_test.go`） |

## 6. 本件没有建立的事（宁缺毋造）

- **没有改票面、没有改任何 AC**。133 的票面仍是 `-done`；这一支要不要写成 AC 追加、要不要在
  `98-attribution` 里回填，是编排者/裁决方的事。
- **没有替裁决方下判**。本件是实现腿取证，不给 `133-*` 的对抗验收结论；`docs/evidence/s1/133-adversarial-acceptance.md`
  那一套仍然有效，只是「C3 这一支红」从此变成「具名披露」。
- **没有重跑 AC#1 的五发（N-3/X4/X8/X12/X14）**。本腿只动 `covered=test` 桶的判红/披露，五发涉及的
  `reached`/`install`/`enumerated` 读数一行没动；自证腿两枚仍在。
- **没有在 linux 上跑整包**（只跑 133 那 4 枚）。`98-attribution` C2（票 255 名册）在两条腿上仍是红的，
  与本腿无关，未处理。
- **`go build ./...` 的 `build/` 目录**：未跟踪，属本机临时件；本腿按规矩**只建不删**，没有清它。
- **另一条腿的 `CgoEnabled=true`** 是探测假设（`scripts/wisp-cli-tests.sh` / ubuntu 腿都带 cgo）；若将来某条 CLI
  腿改成 `CGO_ENABLED=0`，形状会从「披露」回到「红」（保守方向），届时按新读数改这一处。

## 7. 复跑口令（可逐字重跑）

```bash
cd "/d/work/workspace/projects plans/Wisp"
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:/d/work/base/gopath/bin:$PATH"

# 1) windows 腿（本机）
go test -count=1 -timeout 200s -run 'TestAC1AC2DispatchHopGate133|TestAC3SelfProof133|TestAC4Distinctness133' -v ./cmd/wisp/

# 2) linux 腿（容器；快照 + 两枚改过的 _test.go 覆盖进去）
git archive --format=tar -o /d/tmp/133-c3/snap-05e02c4df.tar 05e02c4df
cp cmd/wisp/leg_dispatch_gate_133_test.go cmd/wisp/dispatch_hop_selfproof_133_test.go /d/tmp/133-c3/after/
MSYS_NO_PATHCONV=1 docker run --rm -v "D:/tmp/133-c3:/out" \
  -v ac119-gomodcache:/go/pkg/mod -v ac119-gocache:/root/.cache/go-build -w / \
  golang:1.27-bookworm bash /out/green-after.sh        # 脚本体见本目录之外的 /d/tmp/133-c3/green-after.sh

# 3) 门禁
gofumpt -l cmd/                                       # 期望空
sh scripts/d22scan.sh                                 # 期望 rc=0
```
