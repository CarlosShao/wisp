# T03 对抗验收报告 — 03-skeleton-runtime-rules

- 验收角色：T03-adv（对抗验收，独立上下文；只读仓库 + 临时目录实验，唯一仓库写入为本报告）
- 验收时间：2026-09-19T07:50Z ~ 08:05Z
- 被验收 HEAD：`5dcca43`（feat(skeleton): ticket 03 handoff log）
- 验收环境：Windows 11 x64，Git Bash；Go `go1.27.1`（`D:\work\base\go\bin`）；GCC `(Rev3, Built by MSYS2 project) 16.2.0`（`E:\work\base\msys64\mingw64\bin`）——与 BUILD.md §1 钉死工具链一致；`GOPROXY=https://goproxy.cn,direct`
- 边界遵守：T02-impl 的在途 WIP（`scripts/spike/`、`docs/SLO.md`、`docs/PRECHECK.md`、`docs/evidence/s0/02-*`）**未触碰、未评判**；全程未使用 `git add -A`。

## 逐项裁决

### 1. 全量复跑（build / vet / test / race）— PASS

```
go build ./...   → BUILD_OK
go vet ./...     → VET_OK
go test ./...    → ok  buildinfo 0.044s | observe 0.221s | plugin 0.129s | proc 0.316s
                   （其余 15 个包 [no test files]，cmd/wisp [no test files]）
go test -race ./internal/observe/... ./internal/plugin/... ./internal/proc/... ./internal/buildinfo/...
                 → 全 ok（1.374s / 1.163s / 2.394s / 1.056s）
```

- 测试数对账：`go test -v` 共 38 个测试函数 = **37 PASS + 1 SKIP**（proc 的 `TestHelperProcess` 是
  两进程测试的 helper 载体，父进程内按 go test 惯例 skip，实际逻辑由其余测试以子进程方式触发）。
  实现代理声称的「37 个测试全绿」与 37 个 PASS 完全一致，race 声称一致。
- 环境注记（非缺陷）：裸 shell 下 `CGO_ENABLED=0` 时 `go build ./...` 会在 cmd/wisp 失败
  （`sherpa-onnx-go-windows: build constraints exclude all Go files`）。这与 BUILD.md §3 冻结的
  环境要求一致（mingw64 gcc 必须在 PATH）；设 `CGO_ENABLED=1` + MSYS2 mingw64 gcc 后全绿。
  internal/ 与 cmd/ 中只有 cmd/wisp 直接 import sherpa 绑定（已核实 internal/ 零直接引用）。

### 2. Stub / 假完成扫描 — PASS

- `grep -rE "todo!|panic\(\"not implemented\"|TODO|FIXME" internal/ cmd/ --include="*.go"`：**零命中**。
- 17 个占位包逐一清点：`agent`、`agent/approval`、`agent/scheduler`、`audio`、`ball`、`config`、
  `llm`、`memory`、`models`、`panel`、`risk`、`secret`、`session`、`speech`、`tools`、`watchdog`
  （+ `statemachine` 已有 states.go 词汇表）——每个目录**只有 doc.go**（statemachine 另有 states.go），
  无任何假装功能存在的实现代码，不存在「返回 nil 就算成功」。
- 抽查 5 个 doc.go（risk/tools/llm/audio/panel）：均含三段式 —— Responsibilities（引用 SPEC 条目）、
  **Non-responsibilities**（写明不负责什么、归属哪个包）、**DEFERRED(票号)**（如 risk→票 17/18/19、
  llm→票 09/11、audio→票 13/26、panel→票 33/35）。
- `buildinfo.MinisignPublicKey = "PLACEHOLDER-C29-MINISIGN-PUBLIC-KEY"`：文档化占位（C29 落地真实钥，
  密钥材料永不进仓库，BUILD.md §4 已声明），不假装校验能力。

### 3. 模块表 1:1 对账（SPEC-01 §3）— PASS

- internal/ 实际集合：`agent(+approval,+scheduler) audio ball buildinfo config llm memory models
  observe panel plugin proc risk secret session speech statemachine tools watchdog`
  = 19 个顶层 + 2 个 agent 子包，与 SPEC-01 §3 冻结模块表**逐一对上，无缺失、无私自新增**。
- go.mod module path = `github.com/CarlosShao/wisp` ✓（票面 P10 项）。
- 占位 doc.go 均写明 Non-responsibilities（见上条）。

### 4. D37 完整性 + D43 映射 — PASS

- `internal/observe/errors.go`：17 个 `ErrorClass` 常量按 D37 表序齐全（config…resource）；
  `AllClasses()/Valid()/ValidateErrorClass()` 存在且 errors_test 有 `TestAllClassesHasExactly17`。
- 与 PLAN.md §16.7 D37 原文（L2764 起）逐行对照 17 类的「映射态 + 可重试」全部吻合。
  五类抽核：`config→Unconfigured`、`network→NoNetwork`、`budget/loop→Stuck`、`injection→Error`、
  `resource→WatchdogAlert` —— 与 PLAN 原文一字不差。
- 重试策略与 D37 吻合：network/provider=指数退避（`MaxBackoffAttempts=3`）、rate_limit=尊重
  retry-after（`Retryable()` 对 RetryAfterHeader 强制 `RetryAfter>0`）、audio_device/asr=一次、
  tool=Agent 决定、其余=否。
- `TestClassMappingAndRetry` + `TestMappedStatesAreValidBallStates` 把映射钉到 statemachine 常量，
  两包无法漂移。

### 5. 关键测试真实性抽审（5/5 全读）— PASS

- **TestShutdownOrderAudit**：hook 闭包把实际执行序记入 `order` 切片，逐一与冻结序 1..9 比对、
  审计记录恰 10 条、step 10 只记录不执行、显式断言 step4(音频) 严格先于 step5(ASR 释放)、
  step 8 FreeOSMemory 确实执行。另配 fast-path 变体测试（只允许跳 step7）与有界放弃变体。
  **审计的是真实执行顺序，不是常量复读。**
- **TestPanicInWorkerSurvivesAndCancelsRoot**：真实 `panic("boom")` 注入 → root ctx 被取消、
  兄弟 goroutine 经 ctx 退出、Handle.Err() 返回 internal 分类错误、PanicSink 收到
  {goroutine, root, class=internal, 非空 stack}、PanicCount+1、registry 清零。行为级断言。
- **TestJobScopeKillsChildOnClose**：子进程为重执行的测试二进制 + **READY 握手**（证明已跑进真实
  Go 代码而非仅进程创建成功），`TreePrivateBytes` 断言 ≥1MB（活 Go 子进程）且 <8GB（防拍脑袋），
  `Close()` 后 5s 内子进程必须消失（KILL_ON_JOB_CLOSE），close 后查询必须失败、二次 Close 幂等。
  多子进程记账与「启动失败不留半注册子进程」各有独立测试。
- **TestSingleInstanceSecondExitsAfterSignallingFirst**：真两进程 —— helper 持互斥 → 父进程二次
  Acquire 得 `ErrAlreadyRunning` → `SignalExistingInstance` → helper 打印 ACTIVATED 并退出 →
  互斥可重新获取 → Reset/Set 往返。**真跨进程行为，不是单进程模拟。**
- **TestDisposeTailOrder**：注入审计 MemoryReleaser/SLORecorder，实际收尾序逐项等于
  `worker-exited → deferred-fns → free-os-memory → slo-sample`，peak/tail 字节为正。
- 附加（本验收独立做的**活体冒烟**，临时目录构建，未动仓库 build/）：`wisp.exe` 无参启动打印
  `job object = on, single instance = true` 并进入空事件循环；同会话二次启动打印
  "another instance is running in this session; activated it; exiting" 且 exit 0；
  **第一实例日志实时出现 "activation requested by second launch"**（激活事件真打通）。

### 6. C11 / C30 规格符合 — PASS

- **C11（disposal.go）**：`Dispose()` 以 `disposed` 标志实现幂等（二次调用返回首次结果）；
  `for i := len(fns)-1; i>=0; i--` 逆序；每步经 `registry.Spawn("disposal-worker",…)` 获得独立
  recover 边界（一步 panic 记 Failed[Kind=panic] 后其余步骤继续）；`observe.NewTimeout(3s)`
  共享总预算，超时步与未运行步都记 `disposal_incomplete`（结果字段 + slog.Error + Incomplete()），
  连「Dispose 后再 Defer」也会显式记 late 并置 incomplete（有专测）。
  冻结收尾五步实测顺序：cancel ctx → 等跟踪 goroutine → 释放 fns（native 释放在此）→
  **`DefaultMemoryReleaser.FreeOSMemory()` 真调 `debug.FreeOSMemory`** → SLO 采样 peak/tail。
- **C30（jobscope_windows.go）**：`OpenJobScope` 设 `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`；
  `StartInJob` = Start 后立即 Assign，**Assign 失败先 Kill+Wait 再返回错误**（防逃逸）；
  `TreePrivateBytes` = QueryInformationJobObject 取 Job pid 列表 → 逐进程 psapi
  `GetProcessMemoryInfo` 的 **PrivateUsage** 求和 —— 正是 D32 要求的 Job 树口径（非遍历快照），
  与任务管理器树私有内存可比；消失进程跳过（尽力采样）。

### 7. D22 禁令 — PASS

- 裸 `go func(` 扫描（internal/ + cmd/）：**生产代码（非 _test.go）零命中**；
  唯一合法 go 语句在 `observe/goroutine.go` 的 `Registry.run`。测试文件 3 处
  （proc 两个测试的 Wait 辅助）——D22 禁令针对生产代码；`TestNoBareGoFuncInProductionCode`
  回归测试显式只扫非测试文件并断言扫描面 ≥10 文件，防规则空转。
- emoji 码点扫描（internal/ + cmd/）：**零命中**。
- 明文密钥扫描：零命中（仅文档化 C29 占位常量，见第 2 条）。
- go.mod 依赖白名单（SPEC-01 §3）逐一对账：`sherpa-onnx-go v1.13.8`（白名单）、
  `golang.org/x/sys v0.42.0`（白名单明列）、`pelletier/go-toml/v2 v2.2.4`（白名单明列）、
  三个 indirect 均为 sherpa 绑定同源子模块。**无白名单外依赖。**

### 8. 单实例 — PASS

- 互斥名：prod `Local\wisp-single-instance`、dev `Local\wisp-dev-single-instance`、事件名
  同名 + `.activate`，test 环境 `MutexEnabled=false` —— 与 SPEC-03 §5 环境表**逐字一致**；
  `Local\` 命名空间 = per-session（LayoutFor 注释写明第二用户可各自运行）。
- 票面写「env suffix added in 06」，现 prod/dev 名已落：核对 SPEC-03 §5（冻结 spec）本就钉死
  dev 目录/互斥名分叉，实现照抄 spec；06 真正欠的是 test 环境 DataDir/端点/便携覆盖 ——
  代码显式 `ErrTestLayoutDeferred` 归票 06，无越权。
- 二次启动激活语义：活体冒烟见第 5 条（信号→首实例感知→二次退出 0）。

### 9. 票据对照（.scratch/wisp/issues/03）— PASS（验收标准 6/6）

| # | 票面验收标准 | 证据 | 裁决 |
|---|---|---|---|
| 1 | 二进制引导全部包 + D38(e) 十步干净退出 + 顺序审计测试 | `TestShutdownOrderAudit`/`TestShutdownFastPathSkipsOnlyStep7`/`TestShutdownNilHooksAreRecordedSkipped`/`TestBoot*` + 活体冒烟 | pass |
| 2 | 注册表测试：boot + 一个 no-op 任务后计数回基线 ±1 | `TestNoopTaskReturnsToBaseline`（中途 PerTask=3，收尾 root 排空 + NumGoroutine 回基线 +1 容差） | pass |
| 3 | worker 故意 panic → 进程存活 + internal+栈 + root 取消 | `TestPanicInWorkerSurvivesAndCancelsRoot` | pass |
| 4 | JobScope：子进程随 Job 关闭被杀 + TreePrivateBytes 合理 | `TestJobScopeKillsChildOnClose` + `TestJobScopeTreeAccounting` | pass |
| 5 | DisposalScope：逆序/幂等/单 panic 不阻塞/3s 超时置 incomplete | `TestDisposeRunsInReverseOrder`/`TestDisposeIsIdempotent`/`TestDisposePanicDoesNotBlockOthers`/`TestDisposeTimeoutMarksIncomplete`/`TestDisposeTailOrder`/`TestDeferAfterDisposeIsVisible` | pass |
| 6 | 二次启动：信号首实例后退出 | `TestSingleInstanceSecondExitsAfterSignallingFirst` + 活体冒烟 | pass |

- Progress log：9 条 append-only（最新在最后），从 dispatch → 6 个实现单元 → handoff 完整，
  handoff 中如实自报 b155c6b 事故（见第 10 条）。
- Status 仍为 `in-progress`（收票决定属 orchestrator），符合流程。

### 10. 事故复核（b155c6b 误吞 spike / cfa4175 恢复）— PASS

- `git show b155c6b --stat`：确认该提交确实卷入 4 个 spike 文件（`scripts/spike/common/machine.go`
  `mem.go` `winshell.go`、`scripts/spike/go.mod`，共 843 行）——事故属实。
- `git show cfa4175 --stat`：恰好从索引移除同样 4 个文件（843 deletions），无其他改动 —— 恢复精准。
- 现状：`git status` 仅 `?? scripts/spike/`；`git ls-files scripts/spike` 为空（完全未跟踪）；
  全历史仅这两个 commit 触碰过 scripts/spike。
- 内容完整性：盘上 4 文件与 b155c6b 时点 blob 有差异 —— 逐文件 diff 核实为 T02-impl 的**向前继续
  开发**（对齐 gofmt、新增 user32/kernel32 探针、go.mod 增加 goja/go-webview2 等 require；
  mtime 15:34–15:42 均晚于恢复提交），不是内容损坏；b155c6b 时点版本仍可从 git 历史取回，
  无数据丢失。T02 的 WIP 未受本票影响。

## 问题清单

无 BLOCKER、无 MAJOR。MINOR/备忘 3 条（均不构成本票收口障碍）：

1. **MINOR（计数口径）**：声称「37 个测试全绿」实际是 37 PASS；测试函数总数 38（含按设计在父进程
   skip 的 `TestHelperProcess`）。口径差异无实质影响，建议后续票据 handoff 注明 helper 排除口径。
2. **MINOR（花名册备忘）**：`TemporaryNames` 在 D38b 临时 6 名之外追加了 `disposal-worker`
   （票 03 的簿记名，代码注释已声明"非产品 goroutine"）。语义合理，建议票 08 看门狗/SLO 落地时
   把该豁免写进 D38 花名册文档，防止被当成泄漏误报。
3. **MINOR（环境备忘）**：裸 shell `CGO_ENABLED=0` 下 `go build ./...` 在 cmd/wisp 失败（sherpa
   构建约束）。与 BUILD.md §3 冻结环境一致，非代码缺陷；票 08 CI 落地时需保证 runner PATH 含
   mingw64（BUILD.md §7 已预告）。

## 最终裁决

**VERDICT: PASS**
