# 票 131 续单 — 第 0 件：钉文件的跨平台编译破口（Linux 真跑读数）

**执行方**：`worker-ticket131-followup`（票 131 续单）
**日期**：2026-09-23
**锚定 sha（本轮开工时）**：`3029415284f80739b3e53e28126379ae5ecb4df0`（工作树 `git status --short` 为空）
**本件范围**：只动"让它跨平台编得过 + 门的分母不静默"，判据强度一字未动（改动清单见 §5）

## 1. 机制复核（派单给的机制，我自己 grep 过）

```
$ grep -rn "type sinkInstallRecord" cmd/wisp/
cmd/wisp/resident_sink_nail_127_windows_test.go:135:type sinkInstallRecord struct {

$ grep -rln "sinkInstallRecord" cmd/wisp/
cmd/wisp/leg_sink_nail_131_test.go                    <- 第 1 行 `package main`，文件名无后缀
cmd/wisp/resident_sink_nail_127_windows_test.go       <- _windows 后缀 ⇒ 只在 windows 构建里存在
```

不止那一枚符号。`leg_sink_nail_131_test.go` 依赖的 windows-only 件共 8 枚，全部来自两枚 `_windows_test.go`：

| 被引用的件 | 住在哪 | 类型 |
| --- | --- | --- |
| `sinkInstallRecord` | `resident_sink_nail_127_windows_test.go:135` | 类型 |
| `residentInstallMsg` / `residentEarlyResolverMsg` | 同上 `:94` / `:115` | 常量 |
| `readResidentSink` / `jsonlFilesUnder` | 同上 `:309` / `:344` | helper |
| `TestAC1ResidentLegInstallsItsLogListenerOnDisk` | 同上 `:351` | 用例（注册实参） |
| `readSink` | `logsink_windows_test.go:59` | helper |
| `TestAC2SealNoticeLandsInTheRunLegLogFile` | `logsink_windows_test.go:137` | 用例（注册实参） |

## 2. 改前红（Linux 原生，Docker，不是推理）

仪器：`docker.io/golang:1.27`（`uname -s` = Linux，`go env GOOS/GOARCH` = linux/amd64，`CGO_ENABLED=1`），
挂载 `MSYS_NO_PATHCONV=1 -v /d/work/build-tmp/wisp131fu:/work`，进容器先 `ls -l /work/base/go.mod` 证挂载真生效（有输出）。
被测树：`git archive 3029415 | tar -x` 纯净快照（未变异）。

```
# github.com/CarlosShao/wisp/cmd/wisp
# [github.com/CarlosShao/wisp/cmd/wisp]
vet: cmd/wisp/leg_sink_nail_131_test.go:127:12: undefined: sinkInstallRecord
VET_RC=1
```

⇒ 与 CI run `35817761098` step `go vet (module)` 那行**逐字一致**。

**为什么三个代理会误归因（这一格我量清了）**：在 Windows 上做交叉编译 `GOOS=linux go vet ./cmd/wisp/`，
只会拿到 sherpa 那一句、且**拿不到第二句**：

```
$ GOOS=linux go vet ./cmd/wisp/          # 在 windows 主机上交叉编译，同一枚纯净快照
package github.com/CarlosShao/wisp/cmd/wisp
	imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx
	imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in .../sherpa-onnx-go-linux@v1.13.8
```

包加载在 import 阶段就死了，类型检查根本没跑起来 ⇒ 真伤（引用未定义符号）被那句真的工具链现象挡在后面。
Linux 原生构建里 sherpa 的约束是满足的，vet 会往下走并点名我们的行。⇒ **"整树 rc=1 = 交叉编译假象"这种整体归因不成立**，
本票此后所有 `vet` 读数一律**逐错误行归因**（见 §6 的 `go vet` 双 GOOS）。

## 3. 选 (A) 还是 (B)：先量"改完之后 Linux 上门还剩几分母"

**A 的实测前置（派单要求的那条）**：把钉文件改名成 `_windows_test.go` 之后，Linux 上枚举门**是不是恒绿**？
读数来自同一枚 Docker 仪器、未加任何判据的原始门（`/d/work/build-tmp/wisp131fu/probeA`）：

```
$ go test -count=1 -v -run TestAC4EveryLegIsNailedOrRuled ./cmd/wisp/
    leg_sink_gate_131_test.go:171: AC#4 RED: leg "models" (main.go:74) reaches installLogSink ... no registered nail names it.
    leg_sink_gate_131_test.go:171: AC#4 RED: leg "no-args" (main.go:51) ...
    leg_sink_gate_131_test.go:171: AC#4 RED: leg "run" (main.go:61) ...
    leg_sink_gate_131_test.go:171: AC#4 RED: leg "secret" (main.go:72) ...
    leg_sink_gate_131_test.go:217: AC#4 RED: zero nails registered in this test binary, ...
--- FAIL: TestAC4EveryLegIsNailedOrRuled (0.10s)      GATE_RC=1
```

⇒ **恒绿这个后果不存在**：门的腿清单是从**磁盘上的源码**现读的（`loadMainPackage131` 无视 build tag），
Linux 与 Windows 拿到的是同一张 15 行的表；只有 `nails=` 那一列会因平台而空。
清单不为空 + 四行 install 腿 + `legNails131` 为空 ⇒ A 之下门在 Linux 上是**红**的，不是静默。
（这条是实测原文，不是推理。）

**但那一发红说的是假话**：四条 `no registered nail names it / Fix: write the nail` 把"本平台没编译进钉"
读成"这些腿没有钉"。⇒ 选 **A**，并把门补一句"本平台的分母为 0，我看不见"（§4）。

为什么不是 B：B 要把那 8 枚件（含票 117、票 127 的**两枚用例本体**与 5 枚 helper）从别人的
`_windows_test.go` 里搬到无后缀文件，才能让无后缀的 131 钉文件在 Linux 编得过。
那等于为一枚编译破口重写别人的判据文件，超出本件授权（"只许动让它跨平台编得过"）。
两枚被注册用的用例本身就带 windows 语义（子进程 + CTRL_BREAK + DPAPI），搬过去也不会在 Linux 绿。
⇒ **A（把 tag 挪到 131 的钉文件上）+ 门如实报分母**。

## 4. 改后绿（同一仪器，改后）

```
$ go vet ./cmd/wisp/            # fixA 快照 = 工作树 cmd/wisp 逐字节相同（diff -r 已核）
VET_RC=0
```

门在 Linux 的读数（改后，`-v` 原文节选；完整账本 15 行仍在，四行 `RED listener installed, no nail` 一字未改）：

```
    leg_sink_gate_131_test.go:246: AC#4 RED (the instrument, not the code): zero nails registered in this test binary,
        so the gate has nothing to reconcile and would pass on an empty list.
        GOOS=linux compiled no nail file: every registerLegNail131 call lives in a *_windows_test.go file, so the
        4 leg(s) in the ledger that reach installLogSink (models, no-args, run, secret) are UNREAD on this platform,
        not unnailed - and no row of this ledger is coverage here, in either direction.
        Fix: run this package where the nails compile (scripts/wisp-cli-tests.sh, the windows leg), or give this GOOS
        its own nail file and register it. This reading stays red on purpose: ...
--- FAIL: TestAC4EveryLegIsNailedOrRuled (0.11s)      GATE_RC=1
```

Windows 侧（同一工作树，`PATH=third_party/sherpa-onnx`）门的账本四行**逐字不变**：

```
  leg models  main.go:74  install=true records=true ruled=false nails=TestAC2ModelsLegBooksItsHandOffVerdictOnDisk  -> nailed
  leg no-args main.go:51  install=true records=true ruled=false nails=TestAC1ResidentLegInstallsItsLogListenerOnDisk -> nailed
  leg run     main.go:61  install=true records=true ruled=false nails=TestAC2SealNoticeLandsInTheRunLegLogFile       -> nailed
  leg secret  main.go:72  install=true records=true ruled=false nails=TestAC3SecretLegBooksItsAuditRecordsOnDisk     -> nailed
--- PASS: TestAC4EveryLegIsNailedOrRuled (0.02s)
```

`TestAC2ModelsLeg…` / `TestAC3SecretLeg…` / `TestAC2AC3Degraded…`（含两枚子用例）同跑 **PASS**，`ok 0.133s`。

## 5. 本件改动清单（判据强度：零变化）

| 件 | 动什么 |
| --- | --- |
| `cmd/wisp/leg_sink_nail_131_test.go` → `cmd/wisp/leg_sink_nail_131_windows_test.go` | 只改名（`git mv`，R100）⇒ 隐式 windows tag；用例与断言一字未改 |
| 同文件头 `PLATFORM` 段 | 把旧的"为什么不打 tag"的决定改写成"为什么必须打 tag"，附上 §2/§4 的读数出处；旧决定所依据的事实（`GOOS=linux go vet` rc=1）已被 §2 证伪 |
| `cmd/wisp/leg_sink_gate_131_test.go` 头 + `loadMainPackage131` 注 | 同上：删掉"package main 没有 Linux 构建"这句**假话**，换成实测；补一句 tag 能整枚藏掉钉文件 ⇒ 分母检查是干什么的 |
| `cmd/wisp/leg_sink_gate_131_test.go` 分母检查 | `t.Error` 一句 → `t.Errorf` 如实报名（GOOS、四枚读不见的腿、修法），**没有任何一条红被移除**；只在 `len(legNails131)==0` 这一支生效，Windows 读数逐字不变 |
| `cmd/wisp/secret.go:212` | 注释里的文件名跟着改名，一处 |

未动：任何断言、任何阈值、任何判据谓词、`install`/`records`/`ruled` 三列的算法、`caseLabels131`、`legFrom131`、
`staleRegistrations131`、`internal/**`、`scripts/**`、`rules_gateway.go`、`tools/d22scan/**`、`allowlist.txt`、
`.github/workflows/ci.yml`、任何 golden/阈值、`frontend/`。

## 6. 本件的门禁读数

| 仪器 | 读数 |
| --- | --- |
| `gofmt -l cmd/wisp/` | 空，rc=0 |
| `"$(go env GOPATH)/bin/gofumpt.exe" -l cmd/wisp/` | 空，rc=0（v0.7.0） |
| `go vet ./cmd/wisp/`（windows，工作树） | rc=0 |
| `go vet ./cmd/wisp/`（linux/amd64 原生，Docker） | 改前 rc=1（§2 原文）⇒ **改后 rc=0** |

## 7. 顺手量到、本件不下判断的事实（登记以免变成"顺带以为做了"）

- 同一枚 Linux 容器里把**整包**跑起来（未加 `-run`）：`=== RUN` 70 / 顶层 `--- PASS` 20 / 顶层 `--- FAIL` 20 / `--- SKIP` 0，rc=1。
  红的 20 枚里我抽读的原文都是同一句 `secret: DPAPI is only available on Windows`（`TestSecretUnsetRefusesWhileReferenced`、
  `TestSecretSameNameUnderThreeEnvsIsThreeBlobs`、`TestSecretPortableModeUsesTicket06Seam`、`TestSecretEndToEndConfigRefResolvesAtRequestTime` 等）。
  ⇒ 这与 `scripts/wisp-cli-tests.sh` 头部当年记的"ubuntu 上 19 of 29 枚红"是同一格旧账，**不是本件引入的**，
  也不在本件授权内（本件只保证"编得过 + 分母不静默"）。归票 131 `next=` ③ 那一格：`cmd/wisp` 在 Linux 有分母、
  但那条分母今天不干净；谁要收它，收的是这 20 枚，不是门。
- 派单原文写"验收方全文 618 行"，落盘文件是 **627 行**（`wc -l`）。内容读到的是同一份，只是 §13 那段后来补过；
  我按 627 行的版本读，没有据此改任何判据。
- CI run id `35817761098` 我没有独立开 `gh` 去复取（本机 `gh` 未认证过该仓）；我的替代证据是**同一条错误行在原生
  Linux 容器里逐字复现**。这条差别如实登记：run id 那一格属"未验证"。
