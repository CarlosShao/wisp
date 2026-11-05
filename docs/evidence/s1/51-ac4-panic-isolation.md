# 票 51 落地腿 51-ac4（实现腿 / Tier2 goja 运行时 · `AC#4 panic 隔离`）— 交付证据

- 腿：`51-ac4`＝票 51 的第 4 条 AC（`Panic isolation: JS triggering goja panic -> recovered, other plugin unaffected, host alive`）
- 时刻：2026-11-05 20:0x +08
- 进场 HEAD：`b88f3ef94`（派单写「HEAD ae62f761a」，现跑 `git log --oneline -1`＝`b88f3ef94`；现量 `git diff --stat ae62f761a b88f3ef94` 只动 `docs/evidence/s1/retriage-189.md` 与 `docs/reports/pending-and-issues.md` 两件，`-- internal/plugin` 的输出 **0 行**）
- 写面（**只有这四件**）：
  - `internal/plugin/c24.go`（改，产码：`C24HostPanic` 归因值 + `callHost` 的崩溃 recover + `C24-HOSTPANIC` 审计行 + 崩溃现场栈）
  - `internal/plugin/goja_runtime.go`（改，产码：`classify` 认出宿主崩溃，把归因写进回给模型的 tool 类文案；`guardedCall` 头注）
  - `internal/plugin/panic_isolation_51_test.go`（新，测试：4 枚 e2e）
  - `docs/evidence/s1/51-ac4-panic-isolation.md`（本件）
- ⛔ 未动 `frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`cmd/wisp/**`、`internal/tools/**`、`internal/risk/**`、工单文件、别腿证据件。
- ⛔ 未 push（子代理只 commit）。

## 0. 票面三核 + 与 09-19 模板票面的对账（票面确实过期）

| 核 | 现跑读数 | 判 |
|---|---|---|
| 票真在、非 withdrawn/done | `.scratch/wisp/issues/51-tier2-goja.md:3`＝`**Status:** ready-for-agent`、`:4`＝`Claimed by: —`；池内 `ls \| grep '^51'` 只有 `51-tier2-goja.md`（无 `-done`） | 可领 |
| Blocked-by 真解除 | `17-risk-assessor-c19-done.md`、`20-host-bridge-fs-tools-done.md` 两件都在盘上（`-done`＝防重领唯一键） | 解除 |
| owner 批文真在 | 台账 `A1541` 第 1 条逐字「票51 AC#3（C24 GojaHostAPI **六组白名单**）**批准**」；`A1545` 末行「51-B（9d27a8a7f + c25dc0de7）：C24 六组 9 入口冻结表」 | C24 有冻结面可打 |
| 写面无人占 | 进场 `git status --porcelain -- internal/plugin`＝0 行（只有 `.scratch/` 未跟踪件） | 独占 |

**票面过期的地方（逐条对账，本腿按盘上实况而不是 09-19 模板写）**

| 09-19 票面写的 | 盘上实况 | 本腿怎么处理 |
|---|---|---|
| `Blocked by: 17-risk-assessor-c19, 20-host-bridge-fs-tools` | 两件都 `-done` | 不阻塞 |
| 「C24 initial set（SPEC proposal, finalize + freeze **here**）」 | 已冻结且机主批：`internal/plugin/c24.go` 的 `C24SurfaceTable`（6 组 9 入口），尺＝`TestC24SurfaceIsFrozenAndOwnerApproved` | 不再发明表；只在既有表上加「宿主崩了怎么办」 |
| 并行槽 A（VM 生命周期）/ B（C24 面） | A＝`8a2a40b0f`（`51r1-impl.md`）、B＝`9d27a8a7f`（`51-B-c24-hostapi.md`）都已落 | 本腿是第三跳：把 A 的 recover 与 B 的宿主缝接上 |
| `AC#1 中断` / `AC#2 结果上限` | 块 A 已落（5s 中断 5.01s 实测、256KB 截断） | 不重做 |
| `AC#3 C24 表面` | 51-B 已落（9 入口 e2e、静态扫描+运行时陷阱双查） | 不重做 |
| `AC#5 Enable gate` | `51-ac5`（`eafb13d5f`）已落 | 不重做 |
| `spike-02 T02 结论`（ES2020-minus-modules 等） | 已记在 `51r1-impl.md` §4，落点在票 57 的文档面 | 本腿写面外，不碰 |
| **`AC#4 panic 隔离`** | **只有「把 panic 直投 `guardedCall`」的单元级证明**（`51r1-impl.md` §2 ④ 与「已知边界」第 2 条自己写明：宿主回调 panic 属块 B、真·goja 内部 panic「难以构造」） | **本腿做它** |

**为什么 AC#4 是「第一条还没真做的、且不卡 owner 的实现型 AC」**：`AC#1/2/3/5` 各有自己证据件与尺；`AC#6` 是措辞审计（只出读数，见 §7，不产码）。`AC#4` 的票面句子**从 JS 起头**——「JS triggering goja panic」——而既有证明是从 Go 起头直接投 panic 给 `guardedCall`，**从来没有一次真的 JS→宿主→崩溃**走过这条链。`51-ac5` 证据件 §「前一条 AC 的实况」把 AC#4 记为「已由块 A 落地」，这句在**尺的层面不成立**（详见 §1 的 goja 事实），本腿不照抄它。

## 1. 先量事实：goja 会不会把宿主 panic 变成 JS 异常？

这决定了 AC#4 是「已经有」还是「真的缺」。读本仓锁定的 `github.com/dop251/goja v0.0.0-20261102223426-a177d211bf60`：

| 事实（带行号，可复核） | 后果 |
|---|---|
| `vm.go:6002 exceptionFromValue`：只有 `*Object`/`Value`/`*Exception`/`typeError`/`referenceError`/`rangeError`/`syntaxError`/`runtime.Error`（且仅 `makeslice:`）能变成 `*Exception`；**default 分支 `return nil`** | 宿主 `panic("...")`（string）与「assignment to entry in nil map」这类 `runtime.Error` **都变不成 JS 异常** |
| `vm.go:860 handleThrow`：`ex == nil` 时走到底部 `panic(arg)`（`vm.go:901`） | 原值被**原样重抛** |
| `runtime.go:1498 RunProgram` 的 defer：`asUncatchableException(x)` 为 nil ⇒ `panic(x)`（`runtime.go:1512`） | 崩溃**穿出 `RunProgram`**，落点只能是调用方 goroutine 的 recover |
| `vm.go:919 vm.try` / `vm.go:942 runTryInner` 同样把 `ex == nil` 的 panic 重抛 | 插件 JS 的 `try/catch` **接不住**宿主崩溃（这是要保的性质，不是缺陷） |

⇒ **AC#4 的真实缺口不是「recover 有没有」，而是「崩溃穿出 goja 之后有没有人认领、有没有归因」**：既有 `guardedCall` 会接住，但接住后只说 `tier2 runtime panicked (host exploded while reading)`——**不点名入口、不点名工具、审计无行、`Calls()` 里连这次调用都不存在**（RED 读数见 §4）。

## 2. 交付形状（三件产码，全部落在 `internal/plugin`）

| 落点 | 做什么 | 为什么这么切 |
|---|---|---|
| `C24HostPanic`（`c24.go:305-350`） | 一枚**归因值**：`PluginID/Tool/Entry/Capability/Level/Recovered/Stack`，`Error()` 渲染成一行；`asC24HostPanic` 先做类型断言、再 `errors.As` 兜一层包装 | 「有东西崩了」不可行动；「`wisp.fs.read`（capability `fs.read`、R1 底 L0）在插件 `boom` 的工具 `js.read` 里崩了」可行动。`Stack` 取**崩溃现场**（`callHost` 的 recover 处），比只留信封栈更有用 |
| `callHost`（`c24.go:657-690`） | 宿主调用的**唯一**缝：`defer recover` → **先 `s.record(e.Path)`、先写 `C24-HOSTPANIC` 审计行、先 `slog.Error` 记现场栈**，然后 `panic(hp)` | 顺序是刻意的：**把宿主打挂的那一次调用，正是最需要留痕的那一次**。不 `return`、不 `panic(vm.NewGoError(...))`：宿主崩溃**不许**变成插件能 catch 的 JS 异常（否则插件可以 `while(true){try{wisp.fs.read()}catch{}}` 反复捶一个状态未知的宿主） |
| `classify`（`goja_runtime.go:467-481`） | 认出 `gp.recovered` 里的归因，产出 **`tool` 类**文案：「`plugin "boom" tool "js.read": host entry wisp.fs.read (...) panicked inside its host implementation (...); recovered at the calling-goroutine boundary - the plugin's own JS cannot catch it, and other plugins and the host keep running`」 | 类别不动（`tool`＝`RetryAgentDecided`，模型自决是否重试），只把句子换成**点名**的。另记一条结构化 `slog.Error`（plugin/tool/entry/capability/panic） |

**没有新增 API 面**：`C24Host` 接口、`C24SurfaceTable`、`GojaResult`/`GojaCall` 字段、`C24Reply` 一个都没动。冻结表照旧，`TestC24SurfaceIsFrozenAndOwnerApproved` 不受影响。

## 3. AC#4 逐句 → 尺（4 枚 e2e，全走真 goja + 真 JS 文本）

| AC#4 分句 | 实现落点 | 尺（`panic_isolation_51_test.go`） | 读数 |
|---|---|---|---|
| **JS triggering** goja panic | 插件 JS 调 `wisp.fs.read(...)`，宿主 `CallC24` 崩（两种形状） | `TestC24HostCrashIsRecoveredAndAttributed`（`panic("...")`）；`TestC24HostRuntimeErrorIsRecoveredAndAttributed`（nil map write＝`runtime.Error`） | 都 PASS |
| **recovered**（不让进程死） | `guardedCall` 接住，`classify` 归 `tool` 类 | 同上：`err != nil`、`ClassOfOrInternal(err)==tool`、两次调用的 `host.Count()==1`（**不重试**崩溃的宿主） | PASS |
| **归因**（本腿新增；AC 的行动价值） | `C24HostPanic` + `classify` 文案 + `C24-HOSTPANIC` 行 + `Calls()` | 断言错误文案含 `wisp.fs.read`、`fs.read`、`boom`、`js.read`；审计行恰 1 条且含 `entry=`/`capability=`/`plugin=`/`tool=`；`surface.Calls()==[wisp.fs.read]` | PASS（RED 时 4 处全红，见 §4） |
| **the plugin's own JS cannot catch it** | 宿主崩溃不转 JS 异常 | 第一枚测试的 JS 是 `try{return "swallowed:"+wisp.fs.read(...)}catch(e){return "swallowed:"+e}`，断言 `out.Value==""`（catch 拿不到值） | PASS |
| **other plugin unaffected** | 每调用一 VM、每调用一 surface、runtime 无共享可变态 | `TestC24HostCrashLeavesOtherPluginsRunning`：顺序（崩前/崩后）+ 并发（6 条健康插件 ×3 次 与 6 条崩溃插件同跑，同一 `*GojaRuntime`）；健康侧全部 `{"ok":true}` | PASS |
| **host alive** | watchdog join + scope dispose + registry 归零 | 四枚测试都断言 `reg.Count()==0`，并在崩溃后**再发一次**调用拿到值 | PASS |
| 装配层同样成立（不只在裸 runtime 上） | `Tier2Executor` → `Tier2Gate` → `GojaRuntime` → surface | `TestC24HostCrashThroughTheLoaderKeepsTheGateServing`：门开着两插件，`boom` 崩两次（`C24-HOSTPANIC` 恰 2 行）、`alive` 的 `js.ping` 照常回包、`Snapshot().Runtimes` 两枚都在、registry 归零 | PASS |

## 4. RED-before / GREEN-after（现跑，可重放）

**RED（`git archive HEAD` 到仓外 `D:/tmp/wisp51c-red` ＋ 只放本腿的新测试文件；产码用 HEAD 原样）**：

```
$ go test ./internal/plugin -run 'HostCrash|HostRuntimeError|HostPanic' -count=1 -timeout 300s -v
--- FAIL: TestC24HostRuntimeErrorIsRecoveredAndAttributed (0.01s)
    panic_isolation_51_test.go:170: the failure does not attribute the crash to wisp.fs.write of tool js.write:
      "tool: plugin \"crash\": tier2 runtime panicked (assignment to entry in nil map); recovered, host and other plugins unaffected"
    panic_isolation_51_test.go:173: no audit line names wisp.fs.write; log: (空)
--- FAIL: TestC24HostCrashIsRecoveredAndAttributed (0.02s)
    panic_isolation_51_test.go:113: the failure does not name "wisp.fs.read" ... "js.read"
      "tool: plugin \"boom\": tier2 runtime panicked (host exploded while reading); recovered, host and other plugins unaffected"
    panic_isolation_51_test.go:120: C24-HOSTPANIC audit lines = 0, want 1; log: (空)
    panic_isolation_51_test.go:124: the audit line does not carry "entry=wisp.fs.read" / "capability=fs.read" / "plugin=boom" / "tool=js.read"
    panic_isolation_51_test.go:131: surface.Calls() = [], want exactly [wisp.fs.read]
--- FAIL: TestC24HostCrashThroughTheLoaderKeepsTheGateServing (0.04s)
    panic_isolation_51_test.go:326: the tool error does not name the crashed entry wisp.fs.read: "tool: plugin \"boom\": tier2 runtime panicked (the host implementation crashed); ..."
    panic_isolation_51_test.go:329/350: C24-HOSTPANIC lines = 0, want 1 / 2
--- PASS: TestC24HostCrashLeavesOtherPluginsRunning (0.08s)     <- 隔离这一半本来就在（块 A 的 per-call VM），本腿只是把它第一次钉在「宿主崩溃」这个触发源上
FAIL  github.com/CarlosShao/wisp/internal/plugin  0.138s
```

⇒ **3 红 1 绿**：红的三枚红的全是「归因 / 审计 / Calls 记账」，不红的那枚证明「隔离」这一半本来就成立（本腿没靠改它变绿，也没放宽它）。

**GREEN（本腿的树）**：

```
$ go test ./internal/plugin -run 'HostCrash|HostRuntimeError|HostPanic' -count=1 -timeout 300s -v
--- PASS: TestC24HostRuntimeErrorIsRecoveredAndAttributed (0.02s)
--- PASS: TestC24HostCrashIsRecoveredAndAttributed (0.03s)
--- PASS: TestC24HostCrashThroughTheLoaderKeepsTheGateServing (0.05s)
--- PASS: TestC24HostCrashLeavesOtherPluginsRunning (0.07s)
ok  github.com/CarlosShao/wisp/internal/plugin  0.139s
```

**没有放宽任何断言**：RED 用的就是**最终版**测试文件（产码前没改过任何一条断言、没调过一个阈值）；红的四类（文案点名 / 审计行 / `Calls()` 记账 / 崩溃后计数）一条没删。产码只加「认领 + 归因」，不改判据。

## 5. 变异腿：尺有牙（`.scratch/51-ac4/mutate.sh`，跑完 sha256 校验还原）

每发只改一处，跑同一组四枚测试，读 `--- FAIL` 顶层枚数；末行 `pristine verified`（`c24.go=a0dad134…`、`goja_runtime.go=6d0a8531…`，与进场一致）。

| 发 | 改什么 | 顶层红 | 说明 |
|---|---|---|---|
| M0 | 不改（基线） | **0** | 四枚全绿 |
| M1 | `callHost` 的 `recover` 永不作动（`&& false`） | **4** | 崩溃直接冲到 goja 边界：无归因、无审计行、无 Calls 记账、无边界文案 |
| M2 | `C24-HOSTPANIC` 记号改成 `C24-hostonpanic` | **2** | 证明审计断言绑在那个可 grep 的记号上 |
| M3 | `classify` 不再认归因（走回通用句） | **0** | **如实登记**：归因还有第二个载体——`C24HostPanic.Error()` 本身，`%v` 打出来就带入口名与工具名。⇒ 文案这把尺不靠单一分支；专用分支买到的是**边界句 + 结构化 slog**，不是「能不能点名」 |
| M5 | 崩溃路径不再 `s.record(e.Path)` | **1** | 证明「打挂宿主的那次调用必须留在 `Calls()` 上」有牙 |
| M4 | 崩溃不重抛、当返回值吞掉（调用照跑） | **4** | 证明「插件 try/catch 不能吞」有牙 |

## 6. 门禁（本腿现跑，逐条实数）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/... ./internal/... ./tools/...` | **rc=0**，空 |
| vet | `go vet ./internal/plugin ./internal/plugin/manifestslot` | **rc=0**，空 |
| 本腿尺 | `go test ./internal/plugin -run 'HostCrash\|HostRuntimeError\|HostPanic' -count=1 -timeout 300s` | **ok 0.104s**，4/4 PASS |
| 本腿尺（竞态） | 同上 `-race` | **ok 1.465s**（`gcc` 在，未关 CGO） |
| 整包（去掉两枚既有红） | `go test ./internal/plugin -count=1 -timeout 420s -skip 'Test52SchemaViolationsAreRejected\|Test52CommandKindValidates'` | **ok 7.979s** |
| 邻包 | `go test ./internal/plugin/manifestslot -count=1 -timeout 420s` | **ok 14.216s** |
| 整包（原样） | `go test ./internal/plugin/... -count=1 -timeout 420s` | **rc=1**：`internal/plugin` FAIL，两枚 `Test52*` 红；`manifestslot` ok。**既有红，非本腿**（见 §8） |
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l internal/plugin/ internal/plugin/manifestslot/` | **空**（0 行） |
| D22 扫描 | `sh scripts/d22scan.sh`（共享树） | **rc=0**，`d22scan: clean`；分母 `bans#1-5 internal/=363 cmd/=58`、`ban#6 frontend/=116`、`ban#7 internal/tools/=48`、`ban#8 design/=39 frontend/=116 internal/=905 cmd/=211` |
| 运行环境 | `go version` / PATH | `go1.27.1 windows/amd64`；`PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`、`unset GOTMPDIR TMPDIR TEMP TMP`、**未** `CGO_ENABLED=0` |

## 7. `AC#6 措辞审计` 读数（本腿只量，⛔ 判勾）

- 票面尺（`goja` + 「memory limit」那串英文字面）跑仓库：**2 命中，全在 `docs/evidence/s1/51-ac5-enable-gate.md:98,101`**，是那一件在**引用尺面本身**的自指文本；**产码 0 命中**（`c24.go` / `goja_runtime.go` / 新增测试件都 0）。
- 本件**刻意不逐字重抄那句英文尺面**——重抄一次就自己给 AC#6 添两行命中。
- 放宽到「`goja` + 内存上限/ceiling/cap」：全部命中都是**否定式免责**（`goja_runtime.go:19`、`:463`，`c24.go` 头注、`tier2_gate.go` 头注）⇒ **零条肯定式内存上限承诺**。
- **要不要为字面归零去改别腿证据件**（`51-ac5-enable-gate.md`），属编排者/非实现者裁定；本腿不擅自动别腿文件。判勾归非实现者。

## 8. 既有红（进场就有，本腿未引入；现跑三棵树的对照）

`internal/plugin` 整包今天有两枚**与本腿无关**的红：`Test52SchemaViolationsAreRejected`（`below minimum`/`above maximum`/`not an integer` 三子测试：`{"days":0}`、`{"days":31}`、`{"days":1.5}` 被 `ScreenCommandParams` 接受）与 `Test52CommandKindValidates`（`unknown kind` 子测试：`ParseManifest` 报的是「本字段属于 C14 契约但 Tier-1 base loader 未实现，owner 票 52」，不是点 `.kind`）。

三棵树对照（同一台机、同一套环境变量）：

| 树 | 内容 | 这两枚 | 其它红 |
|---|---|---|---|
| `D:/tmp/wisp51c-pure` | `git archive HEAD`（**零改动**） | **2 FAIL** | 0 |
| `D:/tmp/wisp52chk` | `git archive 5438d104e`（票 52 自己的交付提交） | **2 FAIL** | 0 |
| 本腿树 | 上述 + 本腿 4 件 | **同样 2 FAIL** | 0 |

⇒ ①这两枚**不是本腿引入**（零改动树上同样红）；②`51-B-c24-hostapi.md` §7 那行「`go test ./internal/plugin/...` ⇒ `ok` 7.993s」在今天的盘上**复跑不出来**——`5438d104e..HEAD` 之间**没有任何提交**碰过 `command.go`/`command_argv.go`/`validator.go`/`manifest.go`/`command_52_test.go`，所以那两枚在 51-B 的提交树上也应是红的。**如实登记，不猜原因、不代 52 修**（写面外）。

## 9. 本腿没做的（逐条给理由）

1. **没有替 `C24Host` 写实现**：桥侧 C3/C19/C25 与 `tool_call` 行属组合根（`internal/tools`，`internal/plugin` 不能 import 它）。本腿只把「宿主实现崩了」这条缝钉死。
2. **没有把宿主崩溃转成 `permission_denied` 或 `internal` 类**：AC 要的是「回给模型、让它自决是否重试」⇒ D37 `tool` 类（与 51-B 的 `C24Reply.Refusal`/`C24-HOSTERR` 同类）。
3. **没有动 `C24SurfaceTable`**（六组 9 入口是 A1541 批过的冻结面），也没有动 `C24Reply`/`C24Host` 接口形状。
4. **没有做「内存上限」任何形式的承诺**（goja 没有；只约束钟 + 结果字节 + 崩溃边界）。
5. **工单文件未动**：`51-tier2-goja.md` 的 Status / Progress log / AC 勾归编排者与非实现者。
6. **真机/前端/`cmd/wisp` 装配**：零。

## 10. 交件与复跑

- 提交（**显式 pathspec，只这四条**）：`internal/plugin/c24.go`、`internal/plugin/goja_runtime.go`、`internal/plugin/panic_isolation_51_test.go`、`docs/evidence/s1/51-ac4-panic-isolation.md`。`.scratch/51-ac4/**`（变异脚本、pristine 副本、日志）**不入库**。
- 复跑（任一台机器，逐条可重放）：

```
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP
go build ./cmd/... ./internal/... ./tools/...
go test ./internal/plugin -run 'HostCrash|HostRuntimeError|HostPanic' -count=1 -timeout 300s -v
go test ./internal/plugin -run 'HostCrash|HostRuntimeError|HostPanic' -count=1 -timeout 300s -race
go test ./internal/plugin -count=1 -timeout 420s -skip 'Test52SchemaViolationsAreRejected|Test52CommandKindValidates'
go test ./internal/plugin/manifestslot -count=1 -timeout 420s
/d/work/base/gopath/bin/gofumpt.exe -l internal/plugin/ internal/plugin/manifestslot/   # 空
sh scripts/d22scan.sh                                                                   # clean
sh .scratch/51-ac4/mutate.sh                                                            # M0=0 M1=4 M2=2 M3=0 M5=1 M4=4，末行 pristine verified
# RED 复核：git archive HEAD | tar -x -C /d/tmp/<新目录>，只把 panic_isolation_51_test.go 放进去，跑第一行测试 ⇒ 3 红 1 绿
```
