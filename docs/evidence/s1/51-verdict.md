# 票 51 裁决（非实现者 / 只读腿）— Tier2 goja 运行时：AC#1–#6 逐条

- 裁决腿：`51-verdict`＝**非实现者**（实现腿是 `51-r1`/`51-B`/`51-ac5`/`51-ac4`；本腿只读、只写本件）
- 时刻：2026-11-05 20:2x +08
- 现跑 HEAD：`45371ca0b`（派单写 `e43c097fd`；`git diff --stat e43c097fd 45371ca0b`＝**3 件**且**全非产码**：`docs/evidence/s1/retriage-189.md`、`docs/reports/pending-and-issues.md`、`.scratch/wisp/issues/238-…-no-seal-to-inherit.md`；`-- internal/` **0 行** ⇒ 代码面就是 AC#4 交付时的树）
- 环境：`go1.27.1 windows/amd64`；`PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`、`unset GOTMPDIR TMPDIR TEMP TMP`、**未** `CGO_ENABLED=0`
- 判据：**盘上现跑**，不采信任何证据件的自述读数；证据件只当「尺在哪」的索引。
- ⛔ 未改代码、未改工单、未改别腿证据件；本会话唯一写入＝本文件。

## 1. 现跑总读数（本件所有结论的底数）

| # | 命令 | 读数 |
|---|---|---|
| 1 | `go test ./internal/plugin -run 'HostCrash\|HostRuntimeError\|HostPanic' -count=1 -timeout 300s -v`（AC#4 四枚） | **4/4 PASS**，`ok 0.856s` |
| 2 | 同上 `-race` | **ok 1.264s**（`gcc` 在，未关 CGO） |
| 3 | `go test ./internal/plugin -run 'TestGoja\|TestC24\|TestTier2' -count=1 -timeout 600s -v`（AC#1/2/3/5 的尺＋AC#4 四枚） | **37 枚顶层 PASS、0 FAIL**，`ok 5.085s`（含 AC#1 的 5.00s 全尺寸中断） |
| 4 | `go test ./internal/plugin -count=1 -timeout 420s`（整包） | **rc=1：恰 2 枚红** `Test52SchemaViolationsAreRejected`、`Test52CommandKindValidates`（**既有红，非票 51**；见 §4），其余全绿 `8.308s` |
| 5 | `go test ./internal/plugin/manifestslot -count=1 -timeout 420s` | **ok 16.977s** |
| 6 | `/d/work/base/gopath/bin/gofumpt.exe -l internal/plugin/ internal/plugin/manifestslot/` | **空**（0 行） |
| 7 | `sh scripts/d22scan.sh` | **rc=0 `d22scan: clean`**；分母 `bans#1-5 internal/=363 cmd/=58`、`ban#6 frontend/=116`、`ban#7 internal/tools/=48`、`ban#8 design/=39 frontend/=116 internal/=905 cmd/=211`（与 `51-ac4` 证据件 §6 提交后复跑**逐字相同**） |
| 8 | `grep -rn "goja memory limit" --include=*.go --include=*.md --include=*.toml . \| grep -v '^./.scratch/'`（AC#6） | **4 行**，全在证据件（`51-ac5-enable-gate.md:98,101`、`51-B-c24-hostapi.md:146,149`）；**产码 0 行**（见 §5） |

## 2. 逐 AC 裁决

| AC | 票面要求（简） | 证据件 | 现跑尺（本件复跑） | 判 |
|---|---|---|---|---|
| **#1 中断** | 死循环 ~5s 停、tool 类错误、宿主活 | `51r1-impl.md` §② | `TestGojaWallClockInterrupt` PASS（断言 `elapsed >= 5s` 且 `<= 8s`、`Interrupted=true`、`reg.Count()==0`）；`TestGojaInterruptIsPerCallBounded`、`TestGojaHostCancelInterrupts` 同绿；整组 `ok 5.085s` | **PASS** |
| **#2 结果上限** | 超限返回 → 256KB 截断 + flag | `51r1-impl.md` §③ | `TestGojaResultCapTruncation` PASS：`Truncated=true`、`FullBytes=512KB`、`ReturnedBytes=262144`、rune 边界用例 `"aaa\u00e9"`；`TestGojaBlockAMisconfigRejected` PASS | **PASS**（注：票面写「1MB 串」，尺用 512KB；512KB 已 > 256KB 帽，性质与尺寸无关。测试注释「~1MB」与代码实际 512KB 不一致，属注释笔误，不影响判据） |
| **#3 C24 表面** | 9 入口逐个 e2e；未声明入口 → `undefined` + 违规日志（扫描 + 运行时双查） | `51-B-c24-hostapi.md` | `TestC24EveryEntryReachesTheHost`（9 子测试）、`TestC24ScanRuler`、`TestC24UndeclaredEntryIsUndefinedAndLoggedAtRuntime`、`TestC24AddsNoGlobalsBeyondTheSurface`、`TestC24SurfaceIsFrozenAndOwnerApproved` 等**全 PASS**（37 枚组内） | **PASS**（注：`51-B` §7 那行「`go test ./internal/plugin/...` ⇒ `ok 7.993s`」今天**不复现**——整包恰 2 枚 `Test52*` 红，与 `51-ac4` §8 登记一致；AC#3 自己的尺不受影响） |
| **#4 panic 隔离** | JS 触发 goja panic → recovered、邻插件不受影响、宿主活 | `51-ac4-panic-isolation.md` | 四枚 e2e **全 PASS**（§3 逐句核）；`-race` ok | **PASS** |
| **#5 Enable gate** | 默认关（无 VM）；无 L2 确认开不了；逐插件关 | `51-ac5-enable-gate.md` | `TestTier2GateDefaultOffBuildsNoRuntime`／`EnableNeedsL2Confirm`／`PerPluginDisable`／`BootTrueIsNotAnApproval`／`ReloadL2ControlsTheRuntime` **5/5 PASS**（真 `config.Manager`，非桩） | **PASS** |
| **#6 措辞审计** | 仓库 grep "goja memory limit" → 零 | `51r1` §⑤／`51-B` §8／`51-ac5` §6／`51-ac4` §7 | **字面 4 行、产码 0 行**；4 行全是两件证据件在**引用尺面本身**（自指文本） | **PASS（实质）＋字面未归零**（见 §5；归因：审计报告一旦逐字引尺，就自己制造命中——该 AC 的字面尺自反，处置归编排者） |

## 3. AC#4 深核（本腿重心）：票面分句 → 代码行 → 现跑读数

派单给的 AC#4 口径＝「C24 宿主调用里的宿主 Go panic：按插件隔离、带归因（plugin/tool/entry/capability/floor＋栈）、恰一条 `C24-HOSTPANIC`、该调用记进 `Calls()`、**不**转成 JS 异常」。逐条核：

| 分句 | 代码落点（可点开） | 现跑读数 | 判 |
|---|---|---|---|
| **recovered（进程不死）** | `goja_runtime.go:359-367` `guardedCall` 在**调用方 goroutine** recover；`goja_runtime.go:243-252` 外层兜底（覆盖 dispose 尾巴） | 四枚测试 `err != nil`、`ClassOfOrInternal==tool`、`reg.Count()==0`；崩溃后**再发一次**调用仍出值 | PASS |
| **按插件隔离** | `goja_runtime.go:266-271` 每调用一枚新 VM；`c24.go:370-416` 表面每调用新建、不共享；崩溃不写任何跨调用状态 | `TestC24HostCrashLeavesOtherPluginsRunning`：顺序（崩前/崩后）＋并发（6 健康×3 次 与 6 崩溃同跑，同一 `*GojaRuntime`）健康侧全 `{"ok":true}`；`TestC24HostCrashThroughTheLoaderKeepsTheGateServing`：门开着两插件，`boom` 崩两次后 `alive` 照常、`Snapshot().Runtimes` 两枚、registry 归零 | PASS |
| **归因（plugin/tool/entry/capability/floor＋栈）** | `c24.go:314-335` `C24HostPanic{PluginID,Tool,Entry,Capability,Level,R recovered,Stack}`；`Error()` 一行渲染 `entry (capability …, risk floor …) of plugin … tool …`；`c24.go:660-675` 组装（`Stack`＝崩溃现场 `debug.Stack()`）；`goja_runtime.go:466-487` `classify` 认归因，产出点名文案 | 现跑 ERROR 行逐字：`plugin=boom tool=js.read entry=wisp.fs.read capability=fs.read panic=boom stack="…callHost.func1 c24.go:667 → crashHost.CallC24 …"`；模型侧文案含 `wisp.fs.read`、`fs.read`、`boom`、`js.read`（测试断言）；`level=` 带 `risk.Level.String()`＝`L0/L1/L2`（`internal/risk/assessor.go:56`） | PASS（一处需登记：**栈不在那条 `C24-HOSTPANIC` 行内**，在紧邻的结构化 `slog.Error`（同一事件）与 `C24HostPanic.Stack`；归因可行动性不受影响，见 §3.2 备注 b） |
| **恰一条 `C24-HOSTPANIC`** | `c24.go:670-671` `callHost` 的 defer 里**只调一次** `s.logf("tier2: C24-HOSTPANIC …")`；`bind`（`:693-716`）崩溃路径不返回、不再记 | 测试 `CountContaining("C24-HOSTPANIC")==1`；loader 路径两次崩溃 `==2`（不累积）；全仓 `C24-HOSTPANIC` 只此一处产出（`tier2_loader.go`/executor 不重复记） | PASS |
| **该调用记进 `Calls()`** | `c24.go:669` 崩溃路径**先** `s.record(e.Path)` 再重抛；正常路径 `:686` 只会在未崩时走到 ⇒ **不重记** | `surface.Calls()==["wisp.fs.read"]`（恰 1 条，测试断言）；实现腿变异 M5（删掉崩溃路径的 `record`）打红 1 枚 ⇒ 这把尺有牙（未复跑变异，依据其登记＋代码单点） | PASS |
| **⛔ 不转成 JS 异常** | `c24.go:675` `panic(hp)`（不是 `vm.NewGoError`）；goja 侧重抛原值 | 测试 JS 是 `try{return "swallowed:"+wisp.fs.read(...)}catch(e){return "swallowed:"+e}`，断言 `out.Value==""`（catch 拿不到）；**我独立核了模块源码**：`goja@v0.0.0-…bf60/vm.go` 的 `exceptionFromValue` 只认 `*Object/Value/*Exception/typeError/referenceError/rangeError/syntaxError/runtime.Error(makeslice:)`，`*C24HostPanic` 落 default → `nil`；`handleThrow` 在 `ex==nil` 时 `panic(arg)` 原样重抛（`vm.go:901`）⇒ 插件 `try/catch` 结构上接不住 | PASS |
| **宿主活着（读数不是形容词）** | watchdog 具名 goroutine `stop()` join；`observe.Registry` 归零；门/运行时继续服务 | 四枚测试均断言 `reg.Count()==0` 且崩溃后再调用成功；loader 路径 `Snapshot().Runtimes` 两枚都在 | PASS |
| **未动冻结面（裁决附带核）** | `git show a7093e473 -- internal/plugin/c24.go` 的 diff **0 行**触及 `C24SurfaceTable` | `TestC24SurfaceIsFrozenAndOwnerApproved` 现跑 PASS | PASS |

### 3.1 RED 证据可信度（不采信、只核一致性）

`51-ac4` §4 的 RED 读数（3 红 1 绿，红的是「归因/审计/`Calls()` 记账」，绿的是「隔离本来就成立」）与代码结构一致：归因三件（`C24HostPanic`、`C24-HOSTPANIC`、`classify` 分支）确实都是本次新增，`TestC24HostCrashLeavesOtherPluginsRunning` 只依赖块 A 的 per-call VM ⇒ 它在 RED 树上绿是**结构上应该**的，不构成「测试放水」证据。§5 的 M3（`classify` 不认归因仍 0 红）**如实登记**为「归因还有第二载体 `C24HostPanic.Error()`」——这条登记是对的，不掩盖缺口。

### 3.2 AC#4 残余边界（如实登记，不影响判勾）

- **(a) 票面字面「JS triggering goja panic」的「goja 引擎自身不变量 panic」仍不可从 JS 构造**：本腿证的是真实链路「JS 调 `wisp.*` → 宿主 Go 实现 panic」，引擎内部 panic 只有块 A 的「直投 `guardedCall`」边界测试覆盖。作为裁决：AC#4 的**可行动语义**（JS 侧触发、进程不死、归因、隔离）已满足；但别把它说成「goja 引擎内部崩溃也 e2e 证过」。
- **(b) 栈的位置**：`C24HostPanic.Stack`＋结构化 `slog.Error`（同事件）里有崩溃现场栈；单条 `C24-HOSTPANIC` 行里没有。若将来要求「一行自带栈」，那是一次小改（把 `stack=` 拼进该行），不是缺陷。
- **(c) 归因的第二载体使 M3 变异 0 红**：说明「文案点名」不只靠 `classify` 专用分支；判据（归因必须出现）仍成立。

## 4. 既有红（与票 51 无关；本件复跑确认）

| 树 | 那两枚 `Test52*` | 其余 |
|---|---|---|
| 现跑整包（HEAD `45371ca0b`） | **2 FAIL**（`Test52SchemaViolationsAreRejected`：`below minimum`/`above maximum`/`not an integer`；`Test52CommandKindValidates`：`unknown kind`） | 0 其它红，`manifestslot ok 16.977s` |

- 这两枚在**零改动树**上同样红（`51-ac4` §8 的三树对照：纯 `HEAD` 档案树 2 FAIL、票 52 自己的交付提交 2 FAIL、本腿树同样 2 FAIL）⇒ **非票 51 引入**，属票 52/命令面，写面外。
- 由此：`51-B` §7 与 `51-ac5` §5 的「整包 `ok`」读数**今天不复现**；`51-ac4` §8 已如实登记此事并把「不代 52 修」写清。裁决采信 §8 的登记，不采信 §7 的过期行。

## 5. AC#6（措辞审计）读数与判

- 票面尺字面：**4 行命中**，全在证据件、全是**引用尺面本身**：
  - `docs/evidence/s1/51-ac5-enable-gate.md:98,101`（「票面尺：仓库 grep … -> 零」＋命令原文）
  - `docs/evidence/s1/51-B-c24-hostapi.md:146,149`（同上）
- 产码：`internal/plugin/**` **0 命中**；放宽到 `memory (limit|ceiling|cap)` 且行内含 `goja` 的命中全是**否定式免责**（`goja_runtime.go:19`、`:463`，`c24.go` 头注、`tier2_gate.go` 头注）⇒ **零条肯定式内存上限承诺**。
- **判**：实质 **PASS**（该 AC 的目的＝不许对 goja 的内存上限做假承诺，达到）；字面「零」**做不到**，因为任何如实报告「grep 得到 0」的审计件都会自己制造命中。若编排者坚持字面归零，处置＝把那 4 行改成不含该字面串的措辞（属别腿证据件的一行改动，归编排者，**不是代码缺陷**）。

## 6. 一行结论

**票 51 可结案（closeable）：AC#1/#2/#3/#4/#5 五条实现型 AC 全部 PASS（AC#4 按派单五条判据＝隔离/归因/单条 `C24-HOSTPANIC`/`Calls()` 记账/不转 JS 异常，均在当前树 4 枚 e2e 与 `-race` 下复现）；AC#6 实质 PASS、字面 4 行自指留编排者一句话处置；本件不动任何代码与工单。**

## 7. 复跑（逐条可重放）

```
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP
git log --oneline -1
go test ./internal/plugin -run 'HostCrash|HostRuntimeError|HostPanic' -count=1 -timeout 300s -v   # 4/4 PASS
go test ./internal/plugin -run 'HostCrash|HostRuntimeError|HostPanic' -count=1 -timeout 300s -race
go test ./internal/plugin -run 'TestGoja|TestC24|TestTier2' -count=1 -timeout 600s -v             # 37 PASS
go test ./internal/plugin -count=1 -timeout 420s        # 恰 2 枚 Test52* 红（既有）
go test ./internal/plugin/manifestslot -count=1 -timeout 420s
/d/work/base/gopath/bin/gofumpt.exe -l internal/plugin/ internal/plugin/manifestslot/             # 空
sh scripts/d22scan.sh                                                                              # clean
grep -rn "goja memory limit" --include="*.go" --include="*.md" --include="*.toml" . | grep -v "^./.scratch/"   # 4 行，全证据件
```
