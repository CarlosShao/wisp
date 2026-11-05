# 票 51 落地腿 51-B（实现腿 / Tier2 goja 运行时 · `AC#3 C24 GojaHostAPI`＋loader 接线）— 交付证据

- 腿：`51-B`（重试；上一发 `51-B` 死于无输出、零脏零提交，见台账 `A1543` 末行）
- 时刻：2026-11-05 15:3x +08
- 进场 HEAD：`a2b3b80c9`
- 交付提交：`9d27a8a7f`（本件是它的证据；§7 的"提交后 HEAD 的树"一行是提交后复跑追加的订正笔）
- 写面（**只有这 7 件 + 本件**）：
  - `internal/plugin/c24.go`（新，产码：冻结表＋静态扫描＋运行时表面）
  - `internal/plugin/tier2_loader.go`（新，产码：loader 这一跳）
  - `internal/plugin/manifestslot/tier2.go`（新，产码：Manifest slot 的 Tier-2 Run 适配）
  - `internal/plugin/goja_runtime.go`（改：`GojaCall.Surface`／`GojaResult.ContractViolations`／装表与扫描）
  - `internal/plugin/c24_test.go`、`internal/plugin/tier2_loader_test.go`、`internal/plugin/manifestslot/tier2_51b_test.go`（新，测试）
- ⛔ 未动 `frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`cmd/wisp/**`、`internal/tools/**`、`internal/panel/**`、`internal/perm/**`、`internal/command/**`、别票工单/证据件。
- ⛔ 未 push（子代理只 commit）。

## 进场三核（现跑读数）

| 核 | 读数 | 判 |
|---|---|---|
| 票面真在、非 withdrawn/done | `.scratch/wisp/issues/51-tier2-goja.md:3`＝`**Status:** ready-for-agent`；池内只有 `51-tier2-goja.md` 一件（无 `-done`） | 可领 |
| Blocked-by 真解除 | `17-risk-assessor-c19-done.md`、`20-host-bridge-fs-tools-done.md` 两件都在盘上（`-done` 后缀＝防重领唯一键；20 的 `Status: review` 属其自身 A18 台账项，不影响本票依赖的桥/C3 面已落地） | 解除 |
| 写面无人占 | `git status --porcelain --untracked-files=all -- internal/plugin` 进场时＝0 行；本波派单 `wave10-183-dispatch.txt` 里 51-B 写面＝`internal/plugin/**`，R9-full 写面＝`cmd/wisp`＋`internal/perm`＋`internal/panel`＋`internal/command`，176-ticket 只写工单文件，adj-read 只读 | 独占 |

**owner 已批（本腿不再等拍板）**：台账 `A1541` 第 1 条逐字「1『按推荐』→ 票51 AC#3（C24 GojaHostAPI **六组白名单**）**批准**」，`A1542` 末行「机主上轮已批 51-AC#3/347/R9 ⇒ 本波起派 51-B/339/R9」。⇒ 冻结的初始集＝**http / fs / clipboard / notify / sysinfo / state 六组、9 个入口**，见 §1 表。

## 1. C24 冻结表（`c24.go`：名字 + capability + RiskLevel，一列不多一列不少）

| # | 入口 | capability（C3） | RiskLevel（R1 底） | 依据 |
|---|---|---|---|---|
| 1 | `wisp.http.request` | `net` | **L1** | 出站信道：表面在看到 body 之前分不出 GET/POST，底取「可外发」那档（D34 `web.open`＝L1 同理由）；目标在白名单内、无 taint 时由 C19 保持 L1 |
| 2 | `wisp.fs.read` | `fs.read` | L0 | D34 `fs.read`＝L0（越界由 R2 升 L2） |
| 3 | `wisp.fs.write` | `fs.write` | L1 | D34「新建＝L1／覆盖已存在＝L2」（覆盖由 R8 升，本表不预判） |
| 4 | `wisp.clipboard.read` | `clipboard` | L0 | D34 `clipboard.read` |
| 5 | `wisp.clipboard.write` | `clipboard` | L1 | D34「算外泄通道」（F4） |
| 6 | `wisp.notify.send` | `notify` | L0 | D34 `notify` |
| 7 | `wisp.sysinfo.get` | `sysinfo` | L0 | D34 `system.get` |
| 8 | `wisp.state.get` | **（无）** | L0 | SPEC-07 §5「插件自身 KV，无 capability——插件自己的数据」 |
| 9 | `wisp.state.set` | **（无）** | L0 | 同上 |

**没冻结的东西（本腿不发明成契约）**：C24 只冻结「方法名 + capability + RiskLevel」三列，**没有参数/返回列**。所以本腿定的是一条统一的**调用约定**而不是契约语义：每个入口最多收 **一个 JSON 可序列化对象**、返回宿主给的那个值；对象里是什么字段是**宿主实现的事**。不这么定，「未声明入口＝undefined」就无法用一条规则测。

**没选 Deny**：Deny 属 C19/R3（A 档敏感路径、deny 名单协议），按调用融合取 max；把 Deny 写进表等于宣称该入口被无条件禁死，与它存在的理由相反。

## 2. 三层执法（谁在哪一层拒绝）

| 层 | 落点 | 行为 |
|---|---|---|
| ① C3（桥） | `tools.Bridge.checkCaps`（票 20，**未改**） | manifest 声明经 `tools.Decl.Needs` 下去；未声明即拒绝，且是「拒绝」不是「报错」 |
| ② C24 表面（本腿） | `c24.go` `NewC24Surface` | **只有该 tool 声明了 capability 的入口才装进 VM**（`state.*` 因无 capability 恒装）。没装的入口在 JS 里是 `undefined` → 调用抛 `TypeError`（插件作者可读），并记一条契约违规 |
| ③ C19（桥） | `C24Request.Level` ＋ 宿主实现（**未改**） | 表里的 RiskLevel 只当 R1 底；具体目标（路径、白名单、taint）由宿主的 assessor 按调用融合，只能升不能降 |

`C24Host` 是**唯一**出口：9 个入口没有一个直接依赖文件/套接字/工具，所以 capability 检查、C19 判定、C25 taint、`tool_call` 行全都压在这一个缝上。它的实现属 `internal/tools`＋组合根（本腿写面外，见 §8）。

## 3. 双尺验收（AC#3 逐句 → 尺）

| AC#3 分句 | 实现落点 | 尺（测试） | 读数 |
|---|---|---|---|
| **9 个入口逐个 e2e（mock plugin 各一枚）** | `c24.go` install + bind；测试按**冻结表逐行走**，缺一行即 fatal | `TestC24EveryEntryReachesTheHost`（9 子测试） | 9/9 PASS：宿主收到的 `Entry/Capability/Level/PluginID/Tool/Args` 与表逐字一致，返回值回得到 JS |
| **表里没有的入口** | `judgeC24Path` → `unknown-entry` | `TestC24ScanRuler/entry_outside_the_table_is_unknown`、`…/a_builtin_tool_name_is_not_a_C24_entry` | `wisp.fetch`、`wisp.screen.capture` 各 1 条 |
| **表里有但没声明** | 装表时过滤 → `undefined` | `TestC24ScanRuler/table_entry_outside_the_declaration`、`TestC24SurfaceExposesExactlyTheGrantedEntries` | 静态记 1 条；运行时 `typeof === "undefined"` |
| **静态扫描**（字符串可验） | `ScanC24Source`（`\bwisp(\.ident)+`） | `TestC24ScanRuler` 9 子测试 | 好名 0 条／前缀 `wisp.fs` 0 条／`wispfoo.bar` 0 条／重复折叠 1 条／注释提及也抓（**文本尺的既有射程，如实登记**） |
| **运行时双查**（算出来的访问看不见） | group/root 各一层 Proxy `get` 陷阱 | `TestC24UndeclaredEntryIsUndefinedAndLoggedAtRuntime` | `wisp[g][m]`：静态 **0** 条、运行时 **1** 条 `kind=undeclared-entry path=wisp.fs.write access=dynamic`；`wisp.screen`：**两尺各 1 条**（static＋dynamic） |
| **undefined 而不是报错** | 陷阱回 `undefined`；`logViolation` 记账 | 同上 | 调用方拿到 `undefined`，**调用本身照跑**（`TestC24StaticViolationIsLoggedAndStillRuns`） |
| **契约违规有日志** | `logViolation` → 审计缝（`WithC24Logf`，默认 slog） | 同上 + `TestC24StaticViolationIsLoggedAndStillRuns` | 每条形如 `TIER2-CONTRACT-VIOLATION kind=… path=… capability=… access=static\|dynamic` |
| **表面之外没有全局** | `installC24Surface` 只定义 `wisp` | `TestC24AddsNoGlobalsBeyondTheSurface` | 同一 VM 有/无表面的 `globalThis` 名册**差集恰为 `["wisp"]`**；`fetch/XMLHttpRequest/WebSocket/require/process/module/exports/console/setTimeout/setInterval/setImmediate` 全不在名册里 |
| **每次调用独立** | 表面按调用建、装进按调用新建的 VM | `TestC24SurfaceDiesWithTheCall` | 第 1 调已授 `fs.read`，第 2 调未声明 → `undefined` |
| **宿主拒绝→插件可接的异常；未接→tool 类** | `C24Reply.Refusal` → `goja` 里抛 Error | `TestC24RefusalFromTheHostBecomesAJSThrow` | `catch` 拿得到原文；不 catch 则 `class=tool` 且文案含拒绝原因 |
| **参数只收数据** | JS 侧 `JSON.stringify` 序列化后 Go 侧解码 | `TestC24RejectsBadArguments` 3 子测试 | 两个参数／函数／循环对象 → `TypeError`（class=tool），**宿主 0 次调用** |
| **装配错就 fail-closed** | `installC24Surface` 无宿主即拒 | `TestC24SurfaceWithoutHostFailsClosed` | `class=config`，插件代码一行不跑 |

**入口命名约定**（loader）：C1 工具名带命名空间（`js.echo`），JS 标识符不能含点 ⇒ `Tier2EntryName` 把 `[^A-Za-z0-9_$]` 换成 `_`（`js.echo → js_echo`，首字符数字补 `_`），由 `TestTier2EntryNameIsTheDottedToolNameSanitized` 钉住。没有这条映射，manifest 与 bundle 会各猜一套。

## 4. loader 接线（`tier2_loader.go` ＋ `manifestslot/tier2.go`）

链路：`manifest tool 调用 → Tier2Executor.Execute → Tier2Gate.Execute（AC#5 的三道闸）→ GojaRuntime.Execute（块 A 的钟/上限/recover）→ 装好的 C24 表面 → C24Host.CallC24`。

| 断言 | 尺 | 读数 |
|---|---|---|
| 门在**读代码之前** | `TestTier2ExecutorRefusesBeforeItReadsTheCode`（默认关／无 L2／逐插件关三支） | 三支都 `permission_denied` 且**bundle 读取 0 次、宿主调用 0 次** |
| 端到端 | `TestTier2ExecutorRunsATier2ToolEndToEnd` | manifest params 成为 JS 第一实参、宿主收到 `Tool=echo`、返回值渲染 `{"data":"hello"}`、监视线程 0 泄漏 |
| 声明不外扩 | `TestTier2ExecutorUndeclaredCapabilityStaysUndefined` | `fs.write` 未声明 → `undefined`，宿主 0 次 |
| 缺入口是插件错 | `TestTier2ExecutorMissingEntryIsAToolError` | `class=tool`，文案点名缺的入口名 |
| 身份缺失即拒 | `TestTier2ExecutorEmptyIdentityIsRefused` | 空 plugin/tool → `class=config` |
| 未接线的组件 | `TestTier2ExecutorNeedsEveryWiring` | 门/宿主/源任一为 nil → `class=config` |
| Manifest slot 真跑 | `manifestslot.TestTier2RunFeedsTheManifestSlot` | `Decl.Needs=[fs.read]`、`Decl.PluginID` 在、`Tool.Execute` 出宿主回包、注册进 `tools.Registry` 成功 |

**适配器的一条取舍**：`Tier2Run` 把执行器的 error 转成 `tools.Result{IsError:true}`（**不**上抛）。原因是桥把 `Tool.Execute` 的非 nil error 一律当「工具内部故障（class tool）」，会把门拒绝的 `permission_denied`、缺 bundle 的 `config` 两类都抹成同一句；改成数据后，模型读到的是拒绝原文，行上记的是它本该有的类。

## 5. RED-before / GREEN-after（本腿现跑）

**RED（测试先写、产码文件不存在时）**：`.scratch/51-B/logs/red-plugin.log`、`…/red-manifestslot.log`

```
$ go test ./internal/plugin -run 'C24|Tier2Executor' -count=1 -timeout 120s   # rc=1
internal\plugin\c24_test.go:50:10: undefined: C24Request
internal\plugin\c24_test.go:51:25: undefined: C24Reply
internal\plugin\c24_test.go:653:43: undefined: C24Entry
internal\plugin\tier2_loader_test.go:38:46: undefined: Tier2BundleSource
FAIL	github.com/CarlosShao/wisp/internal/plugin [build failed]

$ go test ./internal/plugin/manifestslot -run 'Tier2Run' -count=1 -timeout 120s  # rc=1
internal\plugin\manifestslot\tier2_51b_test.go:50:20: undefined: plugin.NewTier2Executor
internal\plugin\manifestslot\tier2_51b_test.go:71:26: undefined: Tier2Run
FAIL	github.com/CarlosShao/wisp/internal/plugin/manifestslot [build failed]
```

（RED 里还有两处**测试自己**的毛病，在产码前就改掉了：`containsString` 与 `validator.go:288` 重名 → 复用既有函数；「函数参数」那条期望写错——`JSON.stringify({path: function(){}})` 是 `{}` 不报错，改成整个实参是函数才真报错。**没有一条断言被放宽**，改的是测试自己的形状。）

**GREEN（产码在位）**：

```
$ go test ./internal/plugin -run 'C24|Tier2Executor' -count=1 -timeout 300s -v
--- PASS 42（含子测试）  --- FAIL 0
$ go test ./internal/plugin/manifestslot -run 'Tier2Run' -count=1 -timeout 300s -v
--- PASS 1
```

## 6. 变异腿：六发都要把某枚测试打红（尺有牙）

脚本 `.scratch/51-B/mutate.sh`／每次只改一处、跑完按 pristine 副本还原并用 sha256 校验（末行 `pristine verified`，`c24.go=7ef0e63…`、`tier2_loader.go=1cd1209…`）。

| 发 | 改什么 | 结果 | 打红了哪枚 |
|---|---|---|---|
| M1 | 声明不再决定面：`if e.Capability=="" \|\| have[...]` → `if true` | rc=1 | 4 枚（`SurfaceExposesExactlyTheGrantedEntries`／`UndeclaredEntryIsUndefinedAndLoggedAtRuntime`／`SurfaceDiesWithTheCall`／`Tier2ExecutorUndeclaredCapabilityStaysUndefined`） |
| M2 | 静态尺致盲（`ScanC24Source` 恒返回空） | rc=1 | 9 枚（`ScanRuler` 6 子测试＋`StaticViolationIsLogged…`＋`UndeclaredEntry…` 等） |
| M3 | 表面不装进 VM | rc=1 | 23 枚（9 条入口 e2e 全红＋可见性/隔离/全局/参数/拒绝族） |
| M4 | loader 先读 bundle、门后判 | rc=1 | `Tier2ExecutorRefusesBeforeItReadsTheCode`（读取次数那格） |
| M5 | 运行时陷阱只判不记 | rc=1 | `UndeclaredEntryIsUndefinedAndLoggedAtRuntime` |
| M6 | 宿主拒绝当返回值、不抛 | rc=1 | `RefusalFromTheHostBecomesAJSThrow` |

## 7. 门禁（本腿现跑，逐条实数）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/... ./internal/... ./tools/...` | **rc=0**，空 |
| vet | `go vet ./internal/plugin ./internal/plugin/manifestslot` | **rc=0**，空 |
| test（本包全集） | `go test ./internal/plugin/... -count=1 -timeout 420s` | `ok internal/plugin 7.993s` ＋ `ok internal/plugin/manifestslot 10.719s` |
| test（邻包，防连带） | `go test ./internal/memory ./internal/config -count=1 -timeout 420s` | `ok internal/memory 39.482s` ＋ `ok internal/config 3.035s` |
| gofumpt | `gofumpt -l internal/plugin/ internal/plugin/manifestslot/`（`/d/work/base/gopath/bin/gofumpt.exe`） | **空**（0 行） |
| D22 扫描（**本腿的树**） | 仓外导出树 `D:/tmp/wisp51b-clean-20261105`＝`git archive HEAD` ＋ 本腿 7 件 → `sh scripts/d22scan.sh` | **rc=0**；`d22scan: clean`；分母 `bans#1-5 internal/=361 cmd/=57`、`ban#6 frontend/=116`、`ban#7 internal/tools/=48`、`ban#8 design/=14 frontend/=116 internal/=890 cmd/=207` |
| D22 扫描（**提交后 HEAD 的树**） | 提交 `9d27a8a7f` 之后：`git archive HEAD` → 仓外 `D:/tmp/wisp51b-post-20261105` → `sh scripts/d22scan.sh` | **rc=0**；`d22scan: clean`；分母与上一行**逐字相同**（361/57/116/48/14/116/890/207）⇒ 本腿的提交本身不带任何 D22 违规 |
| D22 扫描（**共享树**） | `sh scripts/d22scan.sh`（仓内，含别人在飞的脏件） | **rc=1**，2 枚 FAIL，**唯一违规行**＝`internal/panel/composer_handlers.go:33 [phantom-citation] 引 "docs/evidence/s1/R9-impl.md" 盘上没有` ⇒ 归 **R9-full 腿**（该文件不在 HEAD 的干净内容里、其证据件尚未落地；`git show HEAD:…` 第 33 行是无此引用的旧注释、`git cat-file -e HEAD:docs/evidence/s1/R9-impl.md`＝rc=128）。⛔ 本腿未动 `internal/panel/**` 一个字节。 |

## 8. `AC#6 措辞审计` 读数（本腿只量，⛔ 判勾）

票面尺：仓库 grep `"goja memory limit"` → 零。

```
$ grep -rn "goja memory limit" --include="*.go" --include="*.md" --include="*.toml" . | grep -v "^./.scratch/"
./docs/evidence/s1/51-ac5-enable-gate.md:98   （引用尺名那一行）
./docs/evidence/s1/51-ac5-enable-gate.md:101  （引用尺命令那一行）
```

- **产码里 0 命中**（`git show HEAD:` 该证据件已含这 2 行 ⇒ 非本腿新增，是前腿引用尺面的自指文本）。
- 放宽到 `memory (limit|ceiling|cap)` 且行内含 `goja` 的全部命中 4 行，**全是否定式免责**：`c24.go:41`、`goja_runtime.go:19`、`goja_runtime.go:457`、`tier2_gate.go:39`。
- ⇒ **零条肯定式内存上限承诺**。判勾归非实现者。

## 9. 本腿没做的（⛔ 越界清单，逐条给理由）

1. **`C24Host` 的实现（`internal/tools` ＋ `cmd/wisp`）**：桥侧 C3、C19 融合、C25 taint、`tool_call` 行都属组合根，且 `internal/plugin` **不能** import `internal/tools`（`tools → agent → memory → plugin` 环）。本腿交的是这条缝的**形状＋契约测试＋审计行样本**，不是替它拍实现。
2. **生产装配**：今天**没有**任何组合根构造 `manifestslot.Manager`（`config/unwired.go:141` 仍写 `plugins.tier2_enabled: not built`），所以 `Tier2Run` 与 `Tier2Gate` 一样是「有测试、有缝、无生产调用者」的状态；把 `cmd/wisp` 接上属下一跳（本波 `cmd/wisp` 唯一写腿是 R9-full，写面互斥）。
3. **C14 的 `entry` 字段**：SPEC-07 §4 的 Tier-2 `entry = "main.js"` 属 manifest 语法面的下一步；本腿把 bundle 来源做成 `Tier2BundleSource` 缝（读哪里是 loader/安装器的权），**没有**动 `validator.go`/C14 语法——那一改要连带 C14 每字段正反例，属另一个写面决定。
4. **真机/前端**：零。本腿不碰面板、不碰 WebView2、不跑任何真机用例。
5. **工单文件未动**：`51-tier2-goja.md` 的 Status／Progress log／AC 勾归编排者与非实现者写。

## 10. 交件与复跑

- 提交（显式 pathspec，只这 8 条路径）：
  `internal/plugin/c24.go`、`internal/plugin/tier2_loader.go`、`internal/plugin/manifestslot/tier2.go`、`internal/plugin/goja_runtime.go`、`internal/plugin/c24_test.go`、`internal/plugin/tier2_loader_test.go`、`internal/plugin/manifestslot/tier2_51b_test.go`、`docs/evidence/s1/51-B-c24-hostapi.md`；`.scratch/51-B/**`（脚本、RED 日志、变异日志）**不入库**。
- 复跑（任一台机器，逐条可重放）：

```
go build ./cmd/... ./internal/... ./tools/...
go test ./internal/plugin/... -count=1 -timeout 420s
go test ./internal/memory ./internal/config -count=1 -timeout 420s
/d/work/base/gopath/bin/gofumpt.exe -l internal/plugin/ internal/plugin/manifestslot/   # 空
sh scripts/d22scan.sh                                                                   # 共享树可能因别人脏件红；本腿树见 §7
bash .scratch/51-B/mutate.sh                                                            # 六发全红，末尾 pristine verified
```
