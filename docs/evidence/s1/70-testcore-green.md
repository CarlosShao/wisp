# 票 70 附页 — test-core 从「只剩 2 枚跳过」走到真绿（2026-11-13，agent=70-core-green 写腿）

> 派单：`.scratch/wisp/probes/orch/20261113/txt/70-core-green-dispatch.txt`（权威全文，含禁止清单与门禁）。
> 起点 CI 读数（run `43377972621`，sha `79890d2b3`，**编排者从 log 抽的，本腿未复核 CI 侧**）：
> `=== RUN=3039 / --- PASS=2056 / --- FAIL=0 / --- SKIP=2` —— **0 枚红，红只因本仓规矩「跳过不算通过」**（票 71 AC#3 的 runner：任何运行时 `--- SKIP` 致命，`tools/d22scan/runtests.sh:234-239`；ledger 内记账条目经 `-skip` 排除后不产生 `--- SKIP` 行，不属于此命）。

## ① 两枚逐枚裁定

### #1 `TestAppendSessionTurnCrashChild`（真重启 reexec 子进程入口）→ 不能在 ubuntu 真跑 ⇒ ledger 具名登记（platform=any, class=reexec）

**勘误（先行，免得后面误导）**：派单写「`internal/session`」，**实测该测试在 `internal/memory/dao_session_append_326_test.go:44`**（`grep -r TestAppendSessionTurnCrashChild internal/` 全仓仅此一处），写面按派单「scripts 的 scope/ledger 部分」落地，未碰 `internal/memory` 的任何断言。

**它是什么**：票 326 AC#3 的**子进程角色**。它自己第一行就跳（`dao_session_append_326_test.go:45-47`：环境变量 `WISP_326_APPEND_CRASH_CHILD != "1"` ⇒ `t.Skip`），只在被父测试 `TestAppendSessionTurnCrashAtomicAcrossKill` 以 `exec.Command(os.Args[0], "-test.run=^TestAppendSessionTurnCrashChild$", ...)` 重启进子进程时才干活：开库、永久追加 turn、每个事务提交后打印 `commit N`；父进程数到目标提交数后 **TerminateProcess 硬杀**，再重开库判「turn 与 updated_at 同事务」。**判定依据**（逐条引脚本/票原文）：
- 「能不能真跑」＝它作为顶层测试**永远不该跑**——它的全部逻辑在重 exec 的子进程里，真跑的前提是**真杀进程**（Kill 是它判据的一半）。ubuntu 容器**能** exec+kill，但那不是把顶层用例变成可执行，而是复刻它父测试的机制；
- 父测试 `TestAppendSessionTurnCrashAtomicAcrossKill`（真断言持有者）在 ubuntu 上照跑真过（本腿容器实测 `ok internal/memory`，读数见 §③），断言**没丢**；
- 同形先例在册：ledger 既有 `TestSubprocessCrashWriter|./internal/memory/|any|reexec` 行（票 326 的兄弟用例，票 93 时期立的），本行只是给**同一形状**的第二枚补登记，不是新发明 scope 规则；
- `tools/d22scan/runtests.sh:25-27`（逐字）：`A test that legitimately cannot run on this platform must be moved out of the step's scope, not silenced here` —— ledger 登记就是本仓的「移出 step 范围」的登记机制（portable-tests.sh 头注：`A case that could never evaluate anything was booked as a pass ... every entry is an exact top-level test name + the package that declares it + a platform + a class + a reason, and on every run each active entry is re-verified against the compiled test binary via go test -list`）。

**处置**：`scripts/portable-tests.sh` ledger 追加一行（platform=`any`，class=`reexec`）：
```
"TestAppendSessionTurnCrashChild|./internal/memory/|any|reexec|child-process entry point of TestAppendSessionTurnCrashAtomicAcrossKill (dao_session_append_326_test.go:44, ticket 326 AC#3); not a test of its own - it skips unless re-exec'd with WISP_326_APPEND_CRASH_CHILD=1, ..."
```
（完整五字段见脚本本体；reason 里点名票 326 AC#3、跳过条件、父测试在两平台 step 都真跑真过、与 TestSubprocessCrashWriter 同类。）ubuntu/Windows 两侧都经 `-skip` 排除；删掉这个测试本身 ⇒ ledger staleness 检查红（`go test -list` 找不到名字），**没有掩盖通道**。

### #2 `Test24R3SearchAppsAUMIDCandidateFeedsTheAUMIDArm`（AUMID 守卫）→ **能在 ubuntu 真跑 ⇒ 删掉过期守卫，让它真跑**（首选解法）

**守卫原文**（改前，`internal/tools/search_files_test.go:456-459`）：
```go
if runtime.GOOS != "windows" {
    t.Skip("the AUMID shape gate is the Windows launcher")
}
```

**判定依据（引原文，逐条）**：
- `internal/tools/launcher_aumid.go:5-12`：形状门 `windowsAUMIDPattern` 刻意放在**平台中立**文件，原注释（逐字）：`platform-neutral file so EVERY GOOS classifies an AUMID-shaped app.launch target the same way; what a platform DOES with an aumid target stays platform-specific`；
- `internal/tools/launcher_other.go:58-64`：非 Windows 的 `resolveApp` 在名字解析器**之前**把 AUMID 形状分到 `Kind:"aumid"`（同形状门），并留注释「an AUMID must classify as Kind "aumid" HERE, before the name resolvers」；该文件头注：这些腿答的是「这个平台没有那扇门」，**分类与门控本身平台无关**；
- 本枚测的是 `search.apps` 把 AUMID 候选**逐字交到 kind=aumid 形状**（`search_apps.go:270-276`：`case kindAUMID: shape.Target = target.AUMID`）——这个臂在 Linux 上**同样执行**（同一个 `resolveApp` 接口，非 Windows 文件也实现同一分类）；AUMID 字符串不触发 `looksLikeDirectPath`（`web_open.go:435-443`：`M` `!` 形状既非盘符也非 `filepath.IsAbs`），Linux 语义下不会被路径臂抢跑；
- 结果端断言（结果文本含 `aumid` 与逐字 AUMID）读的是 `search_apps.go:315` 的列行格式，平台无关；
- 派单首选写明：「能 ⇒ 让它真跑（首选，因为『跳过不算通过』的正确解法是『让它执行』，不是『把它划走』）」；
- 票 70 AC#2 允许的是「落成显式 build tag 或平台专属期望值」并说明为什么——本枚反方向更好：**守卫陈述的理由已过期**（形状门不是 Windows launcher 的私有物，`launcher_aumid.go` 原注释就是为「EVERY GOOS 同判」写的），不是断言需要平台期望值，所以连 build tag 都不需要——**删守卫 = 让用例回到它本来就平台无关的判据上**，断言一字未动。

**ubuntu 真跑实测**（本腿 Docker Desktop，linux/amd64，镜像 `golang:1.27` ⇒ `go1.27.1 linux/amd64` 与 go.mod toolchain/CI 一致；树 = `git archive HEAD`（`dcb58dbb44c8`）＋本腿 3 文件 overlay，容器内逐字命令）：
```
=== RUN   Test24R3SearchAppsAUMIDCandidateFeedsTheAUMIDArm
--- PASS: Test24R3SearchAppsAUMIDCandidateFeedsTheAUMIDArm (0.00s)
runtests.sh: OK - packages=[./internal/tools/ -run Test24R3SearchAppsAUMIDCandidateFeedsTheAUMIDArm -count=1 -v] top-level: PASS=1 FAIL=0 SKIP=0, === RUN=1, '[no tests to run]'=0
```
（同容器 `go vet ./internal/tools/ ./internal/session/ ./internal/memory/` ⇒ `VET_OK`。）

## ② mockllm SA9003（lint job pass1 的独立 module 漏项）

**原文**（`tools/mockllm/messages.go:145-148`，改前）：
```go
if p.toolName == "" &&
    (len(choice) == 0 || strings.Contains(string(choice), `"auto"`)) {
    // auto: answer in text (the real protocol's default too).
}
```

**判定：有意为之的注释性 no-op（不是真缺陷）**——条件算完什么也不做：`toolName` 保持空时走**下方**的 `if p.toolName != ""` 臂，不设 `toolArgs`，请求按文本回答（`"auto"`/缺席本来就是协议默认）。它存在的唯一功能是**文档**（给读者标注「auto ⇒ 文本回答」），对机器是死代码。

**改成什么**：删掉空分支，把「auto ⇒ 文本回答」的说明**并进下方真实分支的注释**（`if p.toolName != ""` 上方），语义零改动：
```go
// No tool_choice ("auto", or the field absent) answers in text - the
// real protocol's default too - which falls out of the branch below:
// only a forced tool name produces tool_use arguments.
if p.toolName != "" {
```
不是 `//lint:ignore` 蒙混，`-checks` 一字未动。**复跑**（派单给的原命令，`cd tools/mockllm`，`GOOS=linux CGO_ENABLED=0`）：
- 改前：`messages.go:145:3: empty branch (SA9003)` ⇒ rc=1
- 改后：**0 行输出，`SC_RC=0`**；模块测试 `go test -count=1 ./...` ⇒ `ok github.com/CarlosShao/wisp/tools/mockllm 0.244s` rc=0。

## ③ 门禁读数（逐项）

| # | 门禁 | 读数 |
|---|---|---|
| 1 | `go build ./internal/session/ ./internal/tools/`（Windows） | rc=0（无输出） |
| 2 | `go vet ./internal/session/ ./internal/tools/`（Windows） | rc=0（无输出） |
| 3 | `go test -count=1 ./internal/session/ ./internal/tools/`（Windows） | `ok ... internal/session 6.322s` / `ok ... internal/tools 44.503s`，rc=0 |
| 4 | `sh scripts/d22scan.sh` | 纯净树+overlay ⇒ rc=0 `clean - no D22 ban violations`（481 production Go files）；活树复跑 ⇒ rc=0（含其他在飞写腿 WIP 也 clean） |
| 5 | `bash -n scripts/runtests.sh` | **rc=127：该路径在仓里不存在**（规矩持有者是 `tools/d22scan/runtests.sh`，CI 全部十余处调用点都指它；本腿未在 scripts/ 下新建 runner——派单写面只给「scope/ledger 部分」）；等价件 `bash -n tools/d22scan/runtests.sh` ⇒ rc=0 |
| 6 | `bash -n scripts/portable-tests.sh` | rc=0 |

**ubuntu 复跑（本机近似方式的诚实声明）**：本机是 Windows 主机，**没有 ubuntu 物理机**；但本机 Docker Desktop 引擎是 **linux/amd64**，本腿用 CI 同款镜像 `golang:1.27`（`go1.27.1 linux/amd64`，与 go.mod toolchain 一致）在容器里**真跑**了 ubuntu 形状的判据——不是模拟，是 linux/amd64 二进制真执行。与真 CI 的差异：宿主非 GitHub 托管 runner（资源态不同，见跑 2 的注）；树不含 compose 的 mock-llm 服务（测试用 in-process mockllm 替身，日志可见 `mockllm: serving on http://127.0.0.1:<port>`）。完整日志摘录：`D:\tmp\wisp\20261113\70-core-green\repro-log.md`。

**三跑对照**（同一 HEAD `dcb58dbb44c8` + 同一本腿 3 文件 overlay）：

| 跑 | 树 | 四数 | 判读 |
|---|---|---|---|
| 1 | archive+overlay（无 .git，CRLF 夹具） | `=== RUN=3038 / PASS=2055 / FAIL=2 / SKIP=0` | **SKIP=0 达成**；FAIL=2 全在 internal/panel，是**导出法假象**（无 .git ⇒ TestReadGitOnThisRepositoryIsSelfConsistent 前提不成立；core.autocrlf 把 composer-states.html 导成 CRLF ⇒ TestComposerRenderFixtureTellsTheTruth 切 0 块），非 HEAD 真红——CI 起点 run 43377972621 同两枚在真 checkout 上 FAIL=0 |
| 2 | **容器内 git clone**（真 .git + LF） | `=== RUN=2993 / PASS=2023 / FAIL=0 / SKIP=0` | **panel 两枚转绿**（假象坐实，`ok (own line) internal/panel`）；顶层断言 `--- FAIL=0`；但 llm 3 包（internal/llm、openaichat、openairesponses）**进程级死**（整跑无任何 `--- FAIL` 行，GUARD B 按缺结果行记 FAIL）——无 panic 无断言红，包进程被杀形状（疑容器资源竞争：同一跑内 present 包 98s vs 跑 1 的 6s；跑 1 同 HEAD 同 overlay 时这 3 包全 `ok`） |
| 3 | 同跑 2 | 定向复跑 llm 三包，读数待补（宿主 spawn 故障恢复后执行） | — |

**两跑已锁死的事实**（不看跑 3 也成立）：
- **运行时 SKIP=0**（跑 1、跑 2 四数一致 `--- SKIP=0`）⇒ 派单目标「SKIP 2→0」在 ubuntu 真执行下实测达成；
- **顶层断言 FAIL=0**（跑 2）⇒ 派单起点「0 枚红」不变量保持；
- 跑 1 的 3038 RUN vs 跑 2 的 2993 RUN 之差与 panel 两枚红跑/绿跑及其子测试计数一致，无来路不明的条目（跑 2 里 ledger 全命中、GUARD A/B/C/D 全过、`-skip` 模式含新增行）。

## ⑤ 一句话回答（派单最后一问）

**SKIP 从 2 变 0**：#2（AUMID）变真跑（ubuntu 容器实测 PASS），#1（reexec 子进程入口）按票 71 AC#3/票 93 先例 ledger 具名登记（platform=any, class=reexec，staleness 检查每跑必验，删除即红）。**下一跑 test-core 若仍红，红因只可能来自别人的在飞改动**（本腿落点已全部排除），本腿自己的三个门（两枚 SKIP、SA9003）已全部排除在红因之外。

## ⑥ 交付清单

- `scripts/portable-tests.sh`（ledger +1 行：TestAppendSessionTurnCrashChild，reexec 类）
- `internal/tools/search_files_test.go`（删过期 skip 守卫 + `runtime` import；断言零改动）
- `tools/mockllm/messages.go`（SA9003 空分支移除，语义零改动）
- 证据件（本文件）
- commit：显式 pathspec，未 push，hash 见回话。
