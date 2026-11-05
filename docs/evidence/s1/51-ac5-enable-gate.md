# 票 51 落地腿 51-ac5（实现腿 / Tier2 goja 运行时 · `AC#5 Enable gate`）— 交付证据

- 腿：`51-ac5 impl`＝票 51 的第 5 条 AC（Enable gate：默认关 / 无 L2 确认不得开 / 逐插件关）
- 时刻：2026-11-05 11:2x +08
- 进场 HEAD：`232b0846c`
- 写面（**只有这三件**）：`internal/plugin/tier2_gate.go`（新，产码）＋ `internal/plugin/tier2_gate_test.go`（新，测试）＋ 本件。⛔ 未动 `internal/tools/**`（票 23 在写）、`internal/config/**`、`cmd/wisp/**`、`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、别票工单/证据件。
- ⛔ 未 push（子代理只 commit）。基线整树 porcelain：tracked 改动 0（只有 `.scratch/` 未跟踪件）。

## 进场三核（现跑读数）

| 核 | 读数 | 判 |
|---|---|---|
| 票面真在、Status 真 `ready-for-agent`、`Claimed by: —` | `.scratch/wisp/issues/51-tier2-goja.md:3`＝`**Status:** ready-for-agent`、`:4`＝`Claimed by: —` | 可领 |
| Blocked-by 真解除 | `17-risk-assessor-c19-done.md`、`20-host-bridge-fs-tools-done.md` 两件都在盘上（`-done` 后缀＝防重领唯一键） | 解除 |
| 写面无人占 | `go test`/`git status`：`internal/plugin` 无未跟踪/未提交改动；`.scratch/wisp/probes/orch/20261105/txt/*.txt` 五份派单里 **0 份**点名 `internal/plugin`；票 23 的写面＝`internal/tools`（其 AC#2 半边已由死腿遗产代提 `111bfdb04`） | 独占 |

**票面前一条 AC 的实况**（决定了本腿选哪条）：`AC#1 中断`／`AC#2 结果上限`／`AC#4 panic 隔离` 已由块 A（`8a2a40b0f`，证据 `docs/evidence/s1/51r1-impl.md`）落地并经编排者复跑（台账 `A1451`：`go test ./internal/plugin -run 'Goja|VM|Disposal|Interrupt|51'` ok 5.042s）。
**`AC#3 C24 表面`＝STOP（需 owner）**：台账 `A1451` 末行逐字"51 块B＝**冻结 C24 `GojaHostAPI` 初始集**（＝改契约文字，AGENTS §2「未定义即停」项）——已把拟议初始集摆机主，等「行／增删」"。⛔ 本腿不预占该拍板。
⇒ 剩下唯一未做且不卡 owner 的实现型 AC＝**AC#5 Enable gate**，本腿做它；`AC#6 措辞审计` 本腿只交读数（见 §6），不越权替非实现者判勾。

## 1. 交付形状（无待拍板项，全部落在 `internal/plugin`）

`Tier2Gate`＝Tier2 goja 运行时**唯一的门**：`Runtime(pluginID)`／`Execute(ctx, call)` 都只从它拿运行时信封，三条规则**全 fail-closed**：

| AC#5 分句 | 实现落点（`tier2_gate.go`） | 尺（测试） |
|---|---|---|
| 默认关，且"关＝没有 VM 存在" | `NewTier2Gate(policy)` **即使** `policy.Tier2Enabled=true` 也以关着开工（`confirmed=false`）；`Runtime` 的三道拒绝全在 `g.factory()` **之前**返回，所以关着的门不留下任何对象可回收。D3「use-to-wake」取字面义：没有调用就没有信封。 | `TestTier2GateDefaultOffBuildsNoRuntime`：工厂计数器**0**＋`Snapshot().Runtimes==[]`＋`Runtime`/`Execute` 都是 `permission_denied` 且文案点名 `tier2_enabled`。**计数器是判据**：光看"它拒绝了"分不出"没建"和"建了但不给"。 |
| 没确认不得开（L2 门 e2e） | `Enable(confirm func() bool)`：`nil`＝拒；`confirm()` 返回 false＝拒且不留痕；`confirm()` 在**锁外**调用（可能等人类，同 `internal/config` 的 `ConfirmLocked` 次序）；通过后还要**再查一次**已生效配置仍为 true（"确认不能把文件关着的东西打开"）。钩子签名故意是 `func() bool` 而不是一个 bool 入参：卡片怎么弹是调用方的事，`internal/plugin` 不该为此长出对 `internal/panel`/`internal/agent` 的依赖。 | `TestTier2GateEnableNeedsL2Confirm`：(a) 文件写 true 到不了运行时且文案点名 `L2`；(b) `Enable(nil)` 拒、门仍关；(c) 人类说不＝拒、钩子只被问 **1** 次、`Confirmed` 仍 false；(d) 人类说是＝开，`Runtime` 给信封、工厂只跑 **1** 次（懒建＋缓存）、`Execute` 真跑出 `"hi"`。 |
| 逐插件关 | `[plugins.<id>].enabled` 是第二道闸：缺键→拒绝（点名字面 `no [plugins.<id>].enabled key`，⛔ 不是"默认开"）；`false`→拒绝（点名 `[plugins.<id>].enabled = false`）；`ApplyPolicy` 只摘掉"刚被关/被删/tier2 转关"那几枚的信封，邻居信封**原对象**保留。tier2 转关同时**清掉 L2 记录**（旧的 yes 不是 yes）。 | `TestTier2GatePerPluginDisable`：两插件各自信封（同一枚就是 bug）；关 demo 后 `Runtimes==[other]`、demo 被拒、`other` 信封**指针未变**且 `Execute` 仍出 `"42"`；缺键分支被拒并点名；整门转关后 `Runtimes==[]`＋`Confirmed=false`＋`Enable(true)` 被拒（配置说 false 时确认不能开）；再开 tier2 时 `Confirmed` 仍 false（旧 yes 已作废）。 |

### `permission_denied` 而不是 `config`（D37 归类理由）

门拒绝是**策略拒绝**，不是配置错误：D37 里 `permission_denied` 的 retry 列＝`RetryNever`、mapped state＝`Acting`（`internal/observe/errors.go:103,131`），模型不得靠重试穿门；`config` 则是"宿主自己配错了"。空 `pluginID` 仍用 `ClassConfig`（与块 A 同一口径：清单一枚都按键名分类）。

## 2. 红线在哪儿：boot 期的漏洞是**真的**，本腿不靠假设说它

配置层对 `[plugins]` 是 locked 段（`internal/config/tiers.go:56` `"plugins": TierLocked`），D36 rule 1 的 L2 再确认只在 **reload** 路径上问（`internal/config/manager.go`：`planLocked` -> `ConfirmLocked`，问不到＝保住旧值）。**初始加载不弹卡**：手改 `config.toml` 写 `tier2_enabled = true` 再重启，`NewManager` 直接把它读成 true，**没有任何卡片、没有任何批准记录**。所以"已生效配置说开"是**数据**不是**裁决**，这正是本腿让 `Enable` 自己再问一次（并把记录存在进程内）的原因。

两条 e2e 测试用的是**真** `config.Manager`（无替身），不是桩：

| 测试 | 它证的东西 |
|---|---|
| `TestTier2GateBootTrueIsNotAnApproval` | 写 `tier2_enabled = true` 的最小 `config.toml` -> `config.NewManager` 读成 true（**坐实 boot 洞**）-> 门仍拒、工厂 0 次 -> `Enable(说是)` 后才建。 |
| `TestTier2GateReloadL2ControlsTheRuntime` | 从 `false` 起手，改写文件为 true 并让 `ConfirmLocked` 说不 -> `Report.Locked=[{section:plugins, approved:false, keys:[plugins.tier2_enabled]}]`、`Config().Plugins.Tier2Enabled` 仍 false、门跟着"已生效配置"关（⛔ 不跟文件字节）；同一份文件再改一次、`ConfirmLocked` 说可以 -> 报告 `approved:true`、配置 applied true、门 `PolicyEnabled=true`；**此时仍要** `Enable` 才给运行时（配置层说了 yes 不等于本进程有记录）。 |

## 3. RED-before / GREEN-after（本腿现跑）

**RED（产码文件不存在时）**：`.scratch/51-ac5/red.log`

```
$ mv internal/plugin/tier2_gate.go .scratch/51-ac5/tier2_gate.go.hidden
$ go test ./internal/plugin -run 'Tier2Gate' -count=1 -timeout 200s   # rc=1
# github.com/CarlosShao/wisp/internal/plugin [.../plugin.test]
internal\plugin\tier2_gate_test.go:34:9: undefined: Tier2Gate
internal\plugin\tier2_gate_test.go:38:39: undefined: Tier2Policy
internal\plugin\tier2_gate_test.go:42:28: undefined: WithTier2RuntimeFactory
...
FAIL	github.com/CarlosShao/wisp/internal/plugin [build failed]
（还原后 sha256=0dd80eb3132403466996baac37fea3b0c1691d71ce0becfbd19b9f58e416308c 逐字回到原值）
```

**GREEN（产码在位）**：

```
$ go test ./internal/plugin -run 'Tier2Gate' -count=1 -timeout 200s
ok  	github.com/CarlosShao/wisp/internal/plugin	0.103s      （5 枚测试全 PASS）
```

## 4. 变异腿：五发都要把某枚测试打红（尺有牙）

脚本 `.scratch/51-ac5/mutate.sh`／`mutate2.sh`；每发只改一处、跑完按 sha256 还原（末行 `final sha256=0dd80eb…` ＝ pristine 逐字）。

| 发 | 改什么 | 结果 | 打红了哪枚 |
|---|---|---|---|
| M1 | `case !g.confirmed:` -> `case false:`（门信任文件里的 true） | rc=1 | `EnableNeedsL2Confirm`／`PerPluginDisable`／`BootTrueIsNotAnApproval`／`ReloadL2ControlsTheRuntime` 四枚 |
| M2 | `if !entry.Enabled {` -> `if entry.Enabled && false {`（tier2 一开＝所有插件都开） | rc=1 | `PerPluginDisable` |
| M3 | 删掉"tier2 转关时 `g.confirmed=false`" | rc=1 | `PerPluginDisable`（旧 yes 复活那格） |
| M4 | 把 `delete(g.runtimes, id)` 换成 `_, _ = ok, entry`（关插件不摘信封） | rc=1 | `PerPluginDisable` |
| M5 | `NewTier2Gate` 里按 policy 预建信封（默认关也建） | rc=1 | `EnableNeedsL2Confirm`（`built==0` 那格） |

> M2／M4 第一版变异写成了编不过的形状（build failed 不算牙），已改成能编的等价变异重跑，上表是重跑读数。

## 5. 门禁（本腿现跑，逐条实数）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/... ./internal/...` | **rc=0** |
| test（本包＋邻包） | `go test ./internal/plugin/... -count=1 -timeout 200s` | `ok internal/plugin 8.363s` ＋ `ok internal/plugin/manifestslot 12.660s` |
| test（e2e 依赖方） | `go test ./internal/config -count=1 -timeout 200s` | `ok internal/config 1.557s` |
| gofumpt | `gofumpt -l internal/plugin/tier2_gate.go internal/plugin/tier2_gate_test.go internal/plugin/` | **空**（0 行） |
| vet | `go vet ./internal/plugin` | **rc=0**，空 |
| D22 扫描 | `sh scripts/d22scan.sh` | **rc=0**；正控 `PASS=35 FAIL=0 SKIP=0`；`d22scan: clean`；分母 `bans#1-5 internal/=356 cmd/=56`、`ban#7 internal/tools/=48`、`ban#8 internal/=876 cmd/=203` |

⚠ **进场即有的本机红，⛔ 归本票**：`go build ./...` 会因 `build/mut_rvw.go` 红（`undefined: voiceBargeinMonitor` 等 6 条）——`build/` 被 `.gitignore:14` 忽略（`git ls-files build` ＝ 0 行），是别的变异腿留下的**本机**产物，不在任何提交里。故本件的 build 尺取 `./cmd/... ./internal/...`（＝全仓受版本控制的生产包），读数 rc=0。

## 6. `AC#6 措辞审计` 读数（本腿只量，⛔ 判勾）

票面尺：仓库 grep `"goja memory limit"` -> 零。

```
$ grep -rn "goja memory limit" --include="*.go" --include="*.md" --include="*.toml" . | grep -v "^./.scratch/"
（0 行，rc=1）
```

放宽到 `memory (limit|cap|ceiling)` 且行内含 `goja` 的**全部**命中，共 3 行，**全是否定式免责**：

- `internal/plugin/goja_runtime.go:19`：`goja has no memory limit and this file does not pretend it has one`
- `internal/plugin/goja_runtime.go:427`：中断文案 `goja offers no memory ceiling; the clock and the %d-byte result cap are the only limits enforced`
- `internal/plugin/tier2_gate.go:39`：本腿新文件的同一句免责（`not a memory bound`）

⇒ **零条肯定式内存上限承诺**。判勾归非实现者。

## 7. 本腿没做的（⛔ 越界清单）

1. **`AC#3` C24 `GojaHostAPI` 表面**：需 owner 先批"冻结初始集"（AGENTS §2 未定义即停项、台账 `A1451`）。前端/桥/capability/taint/audit 全未动。
2. **生产接线**：调用 `Tier2Gate` 的 loader 是块 B 的那一跳（它得先有 C24 表面才能有意义）。所以本腿交付的 `Tier2Gate` **今天没有生产调用者，这是设计而不是遗漏**——它要替换掉的正是"别的腿自己发明一个 loader 形状"那种债。这一点写进了产码文件头（`tier2_gate.go:42-47`），不留给读者猜。
3. **`AC#1/#2/#4` 未重做**：块 A 已交、编排者已复跑（`A1451`）；本腿只在新测试里复用块 A 的信封（`Execute` 出 `"hi"`／`"42"`），不改它一行。
4. **工单文件未动**：`51-tier2-goja.md` 的 Status／Progress log 归编排者与非实现者写（本腿只提交三件，见下）。

## 8. 交件与复跑

- 提交：只这三条路径（显式 pathspec）——
  `internal/plugin/tier2_gate.go`、`internal/plugin/tier2_gate_test.go`、`docs/evidence/s1/51-ac5-enable-gate.md`；`.scratch/51-ac5/**`（脚本与日志）**不入库**（本机留档）。
- 复跑（任一台机器，逐条可重放）：

```
go build ./cmd/... ./internal/...
go test ./internal/plugin -run 'Tier2Gate' -count=1 -timeout 200s -v
go test ./internal/plugin/... ./internal/config -count=1 -timeout 200s
gofumpt -l internal/plugin/tier2_gate.go internal/plugin/tier2_gate_test.go   # 空
sh scripts/d22scan.sh                                                          # rc=0
bash .scratch/51-ac5/mutate.sh ; bash .scratch/51-ac5/mutate2.sh                # 五发全红，sha 还原
```
