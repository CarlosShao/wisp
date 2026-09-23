# 136 AC#1 — 零样本 fail-closed 守卫的钉（装钉前不红 / 装钉后红 / 自证腿不哑）

实现者程 `worker-ticket136-ac1`。本文件只裁 **AC#1**（AC#2..AC#6 不归本程）。
零 emoji。所有读数都是本机原文摘录，未修饰。

## 0. 锚点

| 项 | 值 |
| --- | --- |
| 起点 HEAD（第一次 `git rev-parse --short HEAD`） | `09edf02`，分支 `dev`，`2026-09-23T13:26:47Z`（+08 21:26） |
| 纯净快照来源 | `git archive 09edf02 \| tar -x -C /d/tmp/wisp136-ac1`（快照内无 `.git`，仓内未建 worktree、未 checkout） |
| 读数树 | `/d/tmp/wisp136-ac1`（守卫变异树）／`/d/tmp/wisp136-ac1-pre`（pristine 对照树）／`/d/tmp/wisp136-ac1-nail`（pristine + 本程新用例 = 装钉树）；临时件只建不删 |
| 被钉的那处守卫 | `internal/observe/sampler.go:332-340`（起点树行号）`if len(rep.Samples) == 0 { ... Pass: false, Gate: true ... }` |
| 守卫上游的同族分支 | `internal/observe/sampler.go:290`（`} else if m.PrivateWorkingSetBytes <= 0 {`，零足迹读数丢弃处）——AC#1 的自证腿拆这一发，见 §3 |
| 新用例 | `internal/observe/sampler_zerosample_136_test.go`，commit `4fc65dd` |
| 工具链 | `go version go1.27.1 windows/amd64`；`gofumpt v0.12.0 (go1.27.1)`（本机 `go install mvdan.cc/gofumpt@latest` 所装；**CI 的 `@latest` 未钉版本**，故此处写清本程用的是 v0.12.0）；容器 `golang:1.27` = `go1.27.1 linux/amd64` |
| 冲突面 | `git diff --name-only 09edf02 HEAD -- internal/observe/` ⇒ 只有本程那一枚新文件；`internal/winsec`／`internal/agent`／`cmd/wisp` 一律未碰；`thresholds.go`、任何阈值／golden 一字未动 |

票面判据（`.scratch/wisp/issues/136-...md` AC#1）：结案的唯一判据＝答出
"若那枚守卫被无声删掉，本仓哪一枚仪器会红"——**点名到用例**。本程直答在 §4 末尾。

## 1. 装钉之前：把守卫拆掉，现在确实不红（原文）

### 1.1 变异先证落地（M1 = 票 134 的 FM3 形状：整块删掉）

在纯净快照 `/d/tmp/wisp136-ac1` 上删掉 `sampler.go` 的 332-340 整块，并删掉因此不再被使用的
`"fmt"` import（守卫是 `fmt` 在本文件里唯一的用处；不删 import 就编译不过 —— 那是**工具链会拦的形状**，
不是"沉默"的形状，见 §1.4）。

```
$ python /d/tmp/wisp136-ac1-mutate.py M1        # 在 /d/tmp/wisp136-ac1 里执行
M1 applied: guard block (old 332-340) deleted + unused fmt import removed

$ grep -n 'rep.Verdicts = buildVerdicts' -A 2 internal/observe/sampler.go
330:    rep.Verdicts = buildVerdicts(st, *rep)
331-    rep.Pass = true
332-    for _, v := range rep.Verdicts {

$ sed -n '1,13p' internal/observe/sampler.go     # import 块里 "fmt" 已随守卫一起消失
package observe

import (
        "context"
        "errors"
        "runtime"
        ...
)

$ grep -n 'unmeasurable windows\|len(rep.Samples) == 0' internal/observe/sampler.go
（无输出，grep rc=1）

$ go build ./...
go build rc=0
$ gofmt -l internal/observe/
（无输出）gofmt rc=0
```

变异确实落地：`buildVerdicts` 之后直接进 `rep.Pass = true` 归约循环，守卫块与它的字符串都在树里查无此行，且整树编译通过。

### 1.2 落地后读数：`internal/observe` 整包仍 `rc=0`、54 枚 PASS（`2026-09-23T13:30:34Z`）

```
$ go test -count=1 -v ./internal/observe/
rc=0
--- PASS=54   --- SKIP=0   --- FAIL=0   ^=== RUN=54
--- PASS: TestThresholdTableCoversAllStates (0.00s)
PASS
ok      github.com/CarlosShao/wisp/internal/observe     1.133s
```

同一棵树的 **pristine 基线**（未变异，`2026-09-23T13:28:21Z`）：`rc=0`、PASS=54、SKIP=0、FAIL=0、`=== RUN`=54。
⇒ 两棵树的四数逐位相同。**票面 54 枚这个数字被独立复现了。**

### 1.3 不止这一包：把全仓仪器都跑一遍（两棵树同一时刻）

`go test -count=1 ./...`，pristine（`/d/tmp/wisp136-ac1-pre`）vs 拆掉守卫（`/d/tmp/wisp136-ac1`）：

| 树 | rc | `ok` 包 | `FAIL` 包 | 顶层红名集合 |
| --- | --- | --- | --- | --- |
| pristine `09edf02` | 1 | 22 | 4 | `TestComposerRenderFixtureTellsTheTruth` `TestResolvePerCallBudget` `TestTenOpsInOneToolCallGetOneConfirm` |
| M1（守卫已删） | 1 | 22 | 4 | 同上，`diff` 结果 **IDENTICAL** |

```
$ diff <(grep '^FAIL' pre-full.txt | awk '{print $2}' | sort) <(grep '^FAIL' m1-full.txt | awk '{print $2}' | sort)
IDENTICAL
$ grep internal/observe pre-full.txt m1-full.txt
wisp136-ac1-pre-full.txt:ok  github.com/CarlosShao/wisp/internal/observe  4.868s
wisp136-ac1-m1-full.txt :ok  github.com/CarlosShao/wisp/internal/observe  3.593s
```

**归因写清（不许拿"整树 rc=1"当工具链假象放过）**：那 4 枚红包／3 枚红名在**同一 sha 的 pristine 树上就红**，
两棵树逐名一致，与守卫无关（`cmd/wisp`、`internal/agent/approval`、`internal/panel`、`internal/risk` 此刻都有兄弟在飞，本程一字未碰）。
关键读数是：**被拆掉守卫之后，全仓没有任何一枚仪器的状态发生变化**，`internal/observe` 仍打 `ok`。

### 1.4 顺带量到的一条：整块删除其实是"会响"的形状，沉默的是等价改写

字面删掉 332-340 而不删 `fmt` import ⇒ `go build ./...` rc=1（`"fmt" imported and not used`）。
所以真正沉默的形状是**能编译的删除／放宽**，本程据此把变异做成四发（§3）：M1 删块+删 import、
M4 保留那枚 verdict 行但把 `Pass: false` 翻成 `Pass: true`（守卫还在、牙没了）——两发都 `go build` rc=0、`gofmt` 干净，
都是"看 diff 才看得出来"的形状，两发都被 §2 那枚钉打死。

## 2. 装钉：新增用例与行为断言

文件 `internal/observe/sampler_zerosample_136_test.go`（commit `4fc65dd`），两枚顶层用例。

### 2.1 `TestSampleStateZeroSampleWindowFailsClosed`（AC#1 的硬钉）

fixture：`zeroFootprintTree` 每次读数都成功返回 `TreeMetrics{PIDs: 1, PrivateWorkingSetBytes: 0}`
—— 读得到、但零足迹，正是 `sampler.go:290` 那条"活树零足迹不是可信测量"的丢弃分支。窗口 `10ms/50ms`，状态 `Sleeping`。

断言分三段：

1. **前提腿**（它必须能被 §3 的 M2 打死）：`len(rep.Samples) == 0`、`reads != 0`、`rep.SampleErrors > 0`、
   `rep.LastSampleError` 含 `zero private working set`（报告要说出它丢了什么，票 66 那条老规矩）。
2. **钉**：`rep.Pass == false`；必须存在 `Metric == "sampling"` 的行，且 `Gate == true`、`Pass == false`、
   `Measured` 以 `0 valid` 开头；**归因隔离**——除 `sampling` 行外任何一枚 gate 不许是红
   （fixture 全部落在冻结阈值之下，所以红只能来自守卫，不可能来自阈值被谁放宽）。
3. **出线形状**：`json.Marshal(rep)` 后 `samples` 字段存在且为空、`pass == false`、原文里含
   `"metric":"sampling"` 与 `"gate":true` —— 即 AC#1 那句"出报告即带 `samples=0` 标记 ⇒ 判红"。

### 2.2 `TestSampleStateTrustworthyWindowNotMarkedUnmeasurable`（正向对照腿，防恒真）

同一台 sampler、同一状态，读数可信（16MB private WS）：必须 `len(rep.Samples) > 0`、`SampleErrors == 0`、
**不许长出 `sampling` 行**、`rep.Pass == true`。
没有这一腿，§2.1 可以靠"每个窗口都判红"来恒绿 —— 一枚见谁都咬的守卫不是门。

### 2.3 未变异时两枚都绿（装钉树 `/d/tmp/wisp136-ac1-nail`，`2026-09-23T13:34:25Z`）

```
$ go test -count=1 -v ./internal/observe/
rc=0    --- PASS=56   --- SKIP=0   --- FAIL=0   ^=== RUN=56
=== RUN   TestSampleStateZeroSampleWindowFailsClosed
--- PASS: TestSampleStateZeroSampleWindowFailsClosed (0.05s)
=== RUN   TestSampleStateTrustworthyWindowNotMarkedUnmeasurable
--- PASS: TestSampleStateTrustworthyWindowNotMarkedUnmeasurable (0.06s)
```

（56 = 原 54 + 本程 2；§5 有逐名账。）

## 3. 自证腿拆一发：两腿都不哑（四发变异 × 三态原文）

变异驱动脚本 `/d/tmp/wisp136-ac1-mutate.py`（每次一发都先 `restore` 回 `git show 09edf02:sampler.go` 再落一发，
逐发 `grep -n` 证落地 + `go build ./...` rc=0 后才读数）。

| 发 | 改哪一行 | 落地证明（grep 原文） | build | 红名 |
| --- | --- | --- | --- | --- |
| **M1** | 删 332-340 整块 + 删 `"fmt"` import | `331- rep.Pass = true`（紧跟 `buildVerdicts`），守卫串 rc=1 | rc=0 | `TestSampleStateZeroSampleWindowFailsClosed` |
| **M2** | `290: } else if m.PrivateWorkingSetBytes <= 0 {` → `< 0`（零足迹读数改为可信） | `290:		} else if m.PrivateWorkingSetBytes < 0 {` | rc=0 | `TestSampleStateZeroSampleWindowFailsClosed`（**前提腿**响） |
| **M3** | 同一行 → `>= 0`（所有读数都丢） | `290:		} else if m.PrivateWorkingSetBytes >= 0 {` | rc=0 | `TestSampleStateTrustworthyWindowNotMarkedUnmeasurable`（**对照腿**响） |
| **M4** | `337: Limit: ">0 valid samples", Pass: false, Gate: true,` → `Pass: true`（行还在、牙没了） | `337:			Limit: ">0 valid samples", Pass: true, Gate: true,` | rc=0 | `TestSampleStateZeroSampleWindowFailsClosed` |

### 3.1 三态原文（未变异绿／变异红／还原复绿），`2026-09-23T13:38:1xZ`—`13:38:4xZ`

```
######## M1 三态 ########
-- [变异] 全包 rc / 红名:
rc=1 55 PASS / 1 FAIL
--- FAIL: TestSampleStateZeroSampleWindowFailsClosed (0.05s)
    sampler_zerosample_136_test.go:75: zero-sample window must never pass, got pass=true verdicts=[{Metric:tree_private_bytes Measured:0.0MB Limit:<=25MB Pass:true Gate:t...
-- [变异] 两枚新用例定点:
--- FAIL: TestSampleStateZeroSampleWindowFailsClosed (0.05s)
--- PASS: TestSampleStateTrustworthyWindowNotMarkedUnmeasurable (0.06s)
-- [还原] 全包:
rc=0 56 PASS / 0 FAIL

######## M2 三态 ########
-- [变异] 全包 rc / 红名:
rc=1 55 PASS / 1 FAIL
--- FAIL: TestSampleStateZeroSampleWindowFailsClosed (0.05s)
    sampler_zerosample_136_test.go:61: precondition broken: unmeasurable window produced 5 samples
-- [变异] 两枚新用例定点:
--- FAIL: TestSampleStateZeroSampleWindowFailsClosed (0.05s)
--- PASS: TestSampleStateTrustworthyWindowNotMarkedUnmeasurable (0.06s)
-- [还原] 全包:
rc=0 56 PASS / 0 FAIL

######## M3 三态 ########
-- [变异] 全包 rc / 红名:
rc=1 47 PASS / 5 FAIL
--- FAIL: TestSampleStateAllMetricsAndVerdicts (0.12s)
--- FAIL: TestSampleStateSleepingDiskWriteGateFails (0.06s)
--- FAIL: TestSampleStateWorkPeakMemoryIsTargetNotGate (0.04s)
-- [变异] 两枚新用例定点:
--- PASS: TestSampleStateZeroSampleWindowFailsClosed (0.05s)
--- FAIL: TestSampleStateTrustworthyWindowNotMarkedUnmeasurable (0.06s)
-- [还原] 全包:
rc=0 56 PASS / 0 FAIL

######## M4 三态 ########
-- [变异] 全包 rc / 红名:
rc=1 55 PASS / 1 FAIL
--- FAIL: TestSampleStateZeroSampleWindowFailsClosed (0.05s)
    sampler_zerosample_136_test.go:75: zero-sample window must never pass, got pass=true verdicts=[...
-- [变异] 两枚新用例定点:
--- FAIL: TestSampleStateZeroSampleWindowFailsClosed (0.05s)
--- PASS: TestSampleStateTrustworthyWindowNotMarkedUnmeasurable (0.06s)
-- [还原] 全包:
rc=0 56 PASS / 0 FAIL
```

**自证腿成不成：成。** M2 打响应腿（前提腿，`:61`），M3 打正向对照腿（`:151 trustworthy reads must be sampled`），
M1/M4 打钉（`:75`）；四发各自只响该响的那条腿，另一条腿同发仍绿 —— 两腿都不是哑的。

### 3.2 两条必须报备的读数事故（不粉）

1. **M3 的全包读数不能直接用来点我的腿**：M3 下既有的 `TestSamplerGoroutineAccountingFollowsRegistry`
   在 `sampler_test.go:308` 越界 panic（`index out of range [0] with length 0`），测试二进制当场死掉，
   排在其后的两枚新用例**根本没跑到**（全包 47 PASS / 5 FAIL + panic）。所以 M3 那条腿的红是**定点读数**
   （`-run 'TestSampleState(ZeroSampleWindowFailsClosed|TrustworthyWindowNotMarkedUnmeasurable)$'`，
   另附未变异定点对照：两枚都 PASS）。这顺带说明：M3 那一形不是"沉默"的变异，它被既有仪器以 panic 的形
   拦住了——拦住的不是守卫，是另一件事，本程不把它算成守卫的钉。
2. **本程一度踩了"变异叠加"**：`13:37:04Z` 那次 M3 读数其实是 **M3+M4 叠一发**（脚本没先 restore，
   守卫行的 `Pass` 已被上一发 M4 翻成 true），所以当时两枚腿都红。发现后 `13:37:53Z` 先 `restore`
   再单发 M3 重测，得到 §3.1 里那一份。叠发那份原文留在 `/d/tmp/wisp136-ac1-nail-M3-targeted.txt`，
   **不进结论**。

## 4. 复跑 §1 同一发变异（M1）⇒ 转红，红名点到新用例

装钉树 `/d/tmp/wisp136-ac1-nail` = pristine `09edf02` + `sampler_zerosample_136_test.go`，落**与 §1.1 逐字同一发** M1：

```
$ python /d/tmp/wisp136-ac1-mutate.py M1 && go build ./...
M1 applied: guard block (old 332-340) deleted + unused fmt import removed
go build rc=0

$ go test -count=1 -v ./internal/observe/          # 2026-09-23T13:34:58Z
rc=1
--- PASS=55   --- SKIP=0   --- FAIL=1
--- FAIL: TestSampleStateZeroSampleWindowFailsClosed (0.05s)
    sampler_zerosample_136_test.go:75: zero-sample window must never pass, got pass=true
        verdicts=[{Metric:tree_private_bytes Measured:0.0MB Limit:<=25MB Pass:true Gate:true}
                  {Metric:cpu_percent_all_core Measured:0.000% Limit:<=0.5% Pass:true Gate:true}
                  {Metric:gdi_objects Measured:0 Limit:<200 Pass:true Gate:true}
                  {Metric:handles Measured:0 Limit:<600 Pass:true Gate:true}
                  {Metric:goroutines Measured:0 Limit:<=6 Pass:true Gate:true}
                  {Metric:disk_write_ops Measured:0 Limit:==0 Pass:true Gate:true}
                  {Metric:tcp_connections Measured:0 Limit:==0 Pass:true Gate:true} ...]
FAIL	github.com/CarlosShao/wisp/internal/observe	1.249s
```

红名只有这一枚，正向对照腿同发仍 `PASS`（见 §3.1）；`go test -count=1 -v` 还原后复绿 `rc=0 / 56 PASS`（§3.1 每发的 `[还原]` 行）。

那串红消息本身就是本票的立票理由被演出来了一遍：**一个样本都没取到的窗口，七行 D32 指标全部
`Measured:0 Pass:true`，整份报告 `pass=true`** —— "有报告 ⇒ 里面有数字"就这么静悄悄地变成了"有报告 ⇒ 里面全是 0"。

### 4.1 票 134 验收方答不出的那句，直答

> "若 `sampler.go:332` 被无声删掉，本仓哪一枚仪器会红？我答不出"

**答：`internal/observe` 包的 `TestSampleStateZeroSampleWindowFailsClosed`
（文件 `internal/observe/sampler_zerosample_136_test.go:49`，红点 `:75`）
会红——红在 `:75` `zero-sample window must never pass, got pass=true`。**
它的门禁落点（读码核过，非推断）：`scripts/portable-tests.sh` 的 `core_pin` 清单含
`github.com/CarlosShao/wisp/internal/observe`，CI 的 `test-core`（ubuntu-latest，`ci.yml:288`）
用 `--scope=core` 跑它，且走的是 `tools/d22scan/runtests.sh`（"SKIP 不算 pass、零 PASS 零 FAIL 算 fatal"那把尺）。
⇒ 这枚钉随下一次 push 进入 ubuntu 关门读数；Windows 那条腿（`--scope=windows`）本来就不含 `internal/observe`，
本程的 Windows 原生读数（§5.3/§5.5）是额外覆盖，不是门禁面。
D22 扫形式、P1/P2/P3 看不见 `internal/observe`，这枚是本仓**唯一**目击守卫的仪器。
（CI 历史上有没有 green 是另一笔账，本程未读 run id，见 §6.5。）

## 5. 门禁与四数逐名账

### 5.1 `gofmt` / `gofumpt`（整包 `internal/observe/`）

```
$ gofmt -l internal/observe/                  → 无输出，rc=0
$ gofumpt -l internal/observe/                → 无输出，rc=0     # gofumpt v0.12.0 (go1.27.1)
```

CI 装的是 `mvdan.cc/gofumpt@latest`（`ci.yml:111-114`，**版本未钉**）；本程读数出自 v0.12.0。

### 5.2 `go vet` 双 GOOS

```
$ go vet ./internal/observe/                  rc=0     # GOOS=windows 原生
$ GOOS=linux go vet ./internal/observe/       rc=0     # 交叉；含 _test.go，即新用例在 linux 下也编译
```

（票面那条"交叉 vet 停在 cgo 包加载"的坑适用于 `sherpa-onnx-go-linux` 那族外部包，本包 `./internal/observe/`
两形都真跑完了；真类型读数走容器原生，见 §5.5。）

### 5.3 `go test -count=2 -v ./internal/observe/` 四数（改前／改后）

| | 树 | rc | `=== RUN` | `--- PASS` | `--- SKIP` | `--- FAIL` |
| --- | --- | --- | --- | --- | --- | --- |
| 改前 | pristine `/d/tmp/wisp136-ac1-pre`（`13:34:10Z`） | 0 | 108 | 108 | **0** | 0 |
| 改后 | 工作树（含新用例，`13:40:31Z`） | 0 | 112 | 112 | **0** | 0 |

逐名账（不是差值，按名对）：

```
$ diff <(改前 名->次数) <(改后 名->次数)
46a47
> TestSampleStateTrustworthyWindowNotMarkedUnmeasurable 2
48a50
> TestSampleStateZeroSampleWindowFailsClosed 2
顶层去重枚数：改前 54 -> 改后 56
```

⇒ 动的只有本程新增的两枚（各跑 2 次），**没有既有枚改名、消失、转 SKIP 或换状态**；四数 108→112 全部由这两枚解释。

### 5.4 全仓仪器 `sh scripts/d22scan.sh`（台账各 scope 数不许降）

两次都用同一条脚本（它从自身位置推导仓根）；"改前"跑在 pristine 快照、"改后"跑在工作树，两形 rc 均 **0**，
`runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0` 两形一致。

| scope | 改前 | 改后 | delta |
| --- | --- | --- | --- |
| bans #1-5 internal/ | 203 | 203 | +0 |
| bans #1-5 cmd/ | 22 | 22 | +0 |
| ban #6 frontend/ | 40 | 43 | +3 |
| ban #7 internal/tools/ | 18 | 18 | +0 |
| ban #8 design/ | 16 | 16 | +0 |
| ban #8 frontend/ | 40 | 43 | +3 |
| ban #8 internal/（含注释与 `_test.go`） | 400 | **401** | **+1（本程那枚新用例）** |
| ban #8 cmd/ | 38 | 38 | +0 |
| production Go files 总数 | 225 | 225 | +0 |

**没有任何一个 scope 降。** `frontend/` 的 +3 **不是本程**：改后跑的是此刻的工作树，HEAD 已从 `09edf02`
往前挪过（兄弟程在交件），而 `frontend/**` 本程一字未碰；`internal/` 的 +1 才是本程那枚 `_test.go`。

### 5.5 两形分母（防"变绿 vs 被跳过"）

同一棵装钉树（pristine `sampler.go` + 新用例；工作树的 `sampler.go` 与 `09edf02` 逐字相同，见 §0 冲突面）：

| 形 | rc | `=== RUN` | PASS | SKIP | FAIL |
| --- | --- | --- | --- | --- | --- |
| Windows 原生 `go1.27.1 windows/amd64`，`-count=2 -v` | 0 | 112 | 112 | 0 | 0 |
| Linux 容器 `golang:1.27`（`go1.27.1 linux/amd64`，`-count=2 -v`） | 0 | 112 | 112 | 0 | 0 |

```
$ MSYS_NO_PATHCONV=1 docker run --rm -v /d/tmp/wisp136-ac1-nail:/src -w /src golang:1.27 bash -c \
    'ls -l /src/go.mod && go test -count=2 -v ./internal/observe/'
-rwxrwxrwx 1 root root 883 Sep 23 13:25 /src/go.mod      # 挂载确实挂上了（Git Bash 的静默空挂坑已排）
=== RUN=112  --- PASS=112  --- SKIP=0  --- FAIL=0        # 13:42:44Z
```

两形**逐名账 IDENTICAL**（`--- PASS/FAIL/SKIP` 名->次数 全等），⇒ 新增的绿不是被跳出来的分母缩小。
本程用例不碰软链／POSIX 路径语义，故 `TMPDIR` 为软链那一形对本格无分母影响（未测，见 §6）。

## 6. 未验证项与 next

未验证（本程没有读数，别当已验）：

1. **真实 `wisp slo` 端到端那一发**（票面 `FM`／`FM2` 的形状：真跑 `slo -seconds 0.05` ⇒ 仍交 1 枚真样、
   拆守卫 ⇒ `samples=0 pass=true`）：需要 `cmd/wisp` 与 Windows Job Object，`cmd/wisp` 此刻有兄弟在飞 ⇒
   **未做**。本程钉的是 `internal/observe` 侧的同一判据，出线形状（JSON `samples`+`pass`+`gate` 行）已钉，
   但"CLI 退出码／报告落盘"那一环仍只由 §1.3 的全仓对照证明它今天不设防、也不因本程而设防。
2. `SLO_FULL_*`／新鲜度钉（AC#2..AC#6）与 `internal/winsec`、`scripts/slo-check.ps1` 一律未碰。
3. **软链 `TMPDIR` 那一形**：本程用例无 symlink／tempdir 依赖，未在该形下取数。
4. `CheckSettle`（`sampler.go:449`，`released` 那族）的零样本面只做了**读码**核查：
   零足迹读数被 `:477` 丢弃后 `BackWithinCapMS` 停在 `-1` ⇒ `rep.Pass=false`，**结构上已经 fail-closed**，
   且既有 `TestCheckSettleNeverReachesCap` 盯的是"从未回帽"那一形。本程**没有**为它新增用例、
   也没有对它做真实变异 ⇒ "settle 那侧无钉可拆"这句仍是断言，不是读数。
5. CI 上一次真跑过的 run id 未取（本程只 commit 未 push，共享门禁的历史账不归本格）。

`next=` 编排者按裁决表判 AC#1 是否可勾（本程未自勾）；建议派**非实现者**在
`09edf02..4fc65dd` 之后的 HEAD 上独立复现 §3 的四发（M1/M2/M3/M4 三态），
并单独补 §6.1 那发 CLI 端到端（等 `cmd/wisp` 空出来再做，别撞兄弟）。
D22 mode-6 白名单/allowlist 本程一字未动。
