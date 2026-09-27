# 票 160 · 只读设计核 160-c1 —— `OpenScope` 交回自带关闭动作的句柄：六件事的读数

- 派单：`.scratch/wisp/dispatches/2026-09-27-103x-readonly-160-c1-design-core.md`（09-27 10:3x）
- 本程时刻：09-27 10:33 起手 → 10:4x 交件｜**只读**：`internal/**`、`cmd/**` 零字节未动
- 起手锚：`bb679c6062b832cb54ce68b06ffc4cdd8b7e30c3`｜**本程写件期间 HEAD 动过**（见 §0 第 3 条）
- 台件：`.scratch/wisp/probes/160/c1/{go.mod,shape160.go}`（`gofmt -l` 空、`go vet` 过、`go run` 有读数）
- ⚠ **本件不给"选甲还是选乙"的结论**（派单 §2）。只交两支各自要动哪几枚、各自会红在哪。

---

## 0. 起手四件

1. `date` → `2026-09-27 10:33 +0800`。
2. `git rev-parse HEAD` → `bb679c6062b832cb54ce68b06ffc4cdd8b7e30c3`；`git branch --show-current` → **`dsh/feat/frontend-p0`**。
   ⚠ **派单/任务书说"git 分支 dev"——不成立**：本工作树检出的是 `dsh/feat/frontend-p0`，`dev` 存在（`e0405d5`，09-27 10:29）且是 HEAD 的**祖先**，`dev..HEAD = 2` 枚、`HEAD..dev = 0`。
   ⇒ 影响本件读数吗：**不影响。** `git diff --stat dev bb679c6 -- internal cmd` **空输出**（rc=0），逐枚比 blob：
   `internal/risk/provenance.go` / `internal/tools/bridge.go` / `cmd/wisp/panel_assets.go` / `internal/plugin/disposal.go` / `internal/risk/provenance_test.go`
   五枚在 `dev`、`bb679c6`、`HEAD` **三个锚上 blob 全等**（`git rev-parse <锚>:<路径>` 逐一比过）。本件所有行号可直接搬到 `dev`。
3. **HEAD 在本程中途前移**：`git log -1 HEAD` 现在是 `2a14e40a`（09-27 10:39:05，`fix(frontend): 透明度旋钮接活`），不是我起手量到的 `bb679c6`。
   `git log --name-only dev..HEAD` 显示那 2 枚只碰 `frontend/scripts/gen-tokens.mjs` + `frontend/src/styles/*.css` 与 `.scratch/wisp/dispatches/*` ⇒ 与 `internal/` `cmd/` 无关（第 2 条的 blob 比对是这一条的证据）。**共享工作树里别人在提交，本程没有碰他们任何一枚文件。**
4. `git status --porcelain -- internal/ cmd/` → **空输出**（非空需报回的分支没触发）。
   现场另有派单 §3 列的噪声件（`design/**` 未提交删除、`probes/152/my152.py`、`docs/evidence/s1/152-…-accept-r1.md`），本程一律未碰。
5. 票面 AC 连行号（`.scratch/wisp/issues/160-scope-open-returns-a-handle-carrying-its-own-closer.md`，逐字照抄条目名）：
   - `:26` **AC#1** 先把"改前写不出来"证成读数（未修码上量"编译器不拦、门不响、只有普查尺点名"；量不到就停手）
   - `:27` **AC#2** 换形状：`OpenScope` 交回句柄；`Close` 认身份＋幂等；关闭动作挂到该任务的 `DisposalScope` 上。硬约束三枚＝不许新增 goroutine／不许墙钟差／`RiskLevel` 语义一字节不变（`internal/risk/assessor.go:228`、`internal/tools/bridge.go:610`）
   - `:28` **AC#3** 反向判据：摘掉"认身份"这一味，是否存在一发跨包误关从此打不红？摘掉它有没有任何**外部可见**读数变过？
   - `:29` **AC#4** 明写它治不到什么：从不动的泄漏照样漏；借用形（`Detector(scopeID)`，`provenance.go:603-612`）不该由它治
   - `:30` **AC#5** 契约轴：只许动 `internal/risk/provenance.go` 及其测试＋自己的证据件
   - `:21` "我的推荐"＝关的动作是**闭包里的局部量、不出包**（Pi `session.ts:1475-1490`）
   - `:37` 本票只收 `DEFERRED(C25-loop-wiring)` 的**第 (1) 条**

---

## 1. 本程**没**核什么（免得下游当成已核）

- **没跑票面 AC#1 的变异台件**（摘掉某一发看谁响）。那是代码程的活；本程只交"今天的形状＋两支各自会响在哪"的**静态读数**，外加一份可复算的台件（§5 探针）。
- **没读 Pi `packages/agent/src/harness/pico3/session.ts:1475-1490`**：那枚路径在本仓之外（`D:/work/AI/open source/…`），本程未越界去读，"甲形照它"是票面转抄、非本程现量。
- **没核 `frontend/**`／`design/**`**（派单 §2 明令不查）。
- **没改 `internal/risk` 的任何一字节去实测两支**——§3 的"要动哪几枚"是**读码枚举**，不是试改后的 diff。凡"这样改会红"的句子我都标了它是推断还是读数。
- **`cmd/wisp` 整包在本机量不到**（§7 一条环境失败），所以那条腿的红/绿本程**只有静态证据**。
- **gofumpt 本机没有**：`command -v gofumpt` 空、`~/go/bin` 空；CI 是 `go install mvdan.cc/gofumpt@latest`（`.github/workflows/ci.yml:136-141`）。本程台件只做到 `gofmt -l` 空＋`go vet` 过。

---

## 2. 六件事

### 2.1 今天的形状现量（签名＋全部调用点，逐枚点名）

**签名**（`internal/risk/provenance.go`，锚 bb679c6＝dev 同 blob）：

| 行 | 原文 |
|---|---|
| `:339` | `func (p *Provenance) OpenScope(scopeID string) {` |
| `:350` | `func (p *Provenance) CloseScope(scopeID string) {` |

两枚都是**无返回值、无 error、无身份**；`OpenScope` 幂等（`:342-344` 只在缺席时塞 `nil`），`CloseScope` 就是 `delete(p.scopes, scopeID)`（`:353`，天然幂等）。

**生产调用点：全仓非 `_test.go` 共 3 枚**（`grep -n "OpenScope\|CloseScope" --include="*.go" internal cmd tools` 现量，逐枚）：

1. `internal/tools/bridge.go:642` — `		b.prov.OpenScope(taskID)`（在 `OpenTask` `:633-644` 内，前面 `:637-640` 是桥自己的账本 `b.scopes`）
2. `internal/tools/bridge.go:695` — `		b.prov.CloseScope(taskID)`（在 `CloseTask` `:684-699` 内，`if open` 由 `:690` 的桥账本决定）
3. `cmd/wisp/panel_assets.go:232` — `	prov.OpenScope(taintSourceScopeID)`（**只开不合**：`detector()` `:227-237` 里现建 `prov := risk.NewProvenance(...)`（`:231`）、开（`:232`）、`Mark`（`:234`）、`Detector`（`:236`）就返回；`cmd/wisp` 下 `CloseScope` **零枚**）

⇒ 派单"bridge.go 的开合两处、panel_assets.go 那枚只开不合"**三句全成立**（现量在案）。
⇒ 派单"两枚唯一的生产调用点"（＝owner 批准口径）**只对一半**：**开**的调用点 2 枚（`bridge.go:642`、`panel_assets.go:232`），**关**的调用点 1 枚（`bridge.go:695`）；`tools/` 目录下**零枚**。

**测试面（同一把尺的尾项，逐枚点名，共 26 枚）**：
- `internal/risk/provenance_test.go` 开 12 枚：`:23 :297 :309 :322 :340 :359 :454 :722 :780 :792 :798 :812`；关 3 枚：`:310 :358 :364`
- `internal/risk/taintmatch_test.go` 开 3 枚：`:176 :182 :190`；关 1 枚：`:181`
- `internal/risk/syncdirs_redteam_windows_test.go` 开 3 枚：`:89 :154 :222`
- `internal/risk/syncdirs_test.go` 开 2 枚：`:378 :459`（`e.p.OpenScope`）
- `internal/risk/syncdirs_windows_test.go` 开 1 枚：`:95`
- `internal/risk/provenance_syncdirs_other_test.go` 开 1 枚：`:55`
- `internal/tools/bridge_scope_open_ticket158_test.go:20` — 只是**注释**里提到 `OpenScope`，零枚调用
⇒ `internal/risk` 测试里 **OpenScope 22 枚、CloseScope 4 枚**（`internal/tools`、`cmd/wisp` 测试两枚调用都没有）。
⇒ **全部 19 枚 `internal/risk/*_test.go` 都是 `package risk`（包内测试）**（`grep -h "^package" | sort | uniq -c` → `19 package risk`，`package risk_test` 零枚）。这条在 §2.2 与 §2.3 都会咬到。

**"门"今天长在哪儿**（不是 `tools/d22scan`）：派单说的"门"＝`.scratch/wisp/probes/154/gate-clauses.sh` 的 **G5** 那一节（`:196-228`），尺是 `git grep` 的文件粒度成对普查，射程 `':!*_test.go' ':!internal/risk/*'`；票 158 的验收版把同一把尺复制成 `.scratch/wisp/probes/158/accept-r1/pair-census-a1.sh:44-45`，另有一把误报审计尺 `.scratch/wisp/probes/158/g5-fp-audit.sh:12-15`。
⚠ 该尺**只剔注释、不剔字符串**（`g5-fp-audit.sh:33` 自己写着这条），而 `provenance.go:477` 与 `:578` 的 `Origin: "scope is not open (OpenScope missing or already closed)"` 是**字符串**——今天它们被 `internal/risk` 的排除项挡住了；**改名或把这两枚串挪出 `internal/risk` 会直接动 G5 的名册**。

### 2.2 依赖方向这道题（本单最要紧的一问）——读数与结论句

命令与读数（两次现跑，主模块 `github.com/CarlosShao/wisp`，见 `go.mod:1`）：

```
$ go list -deps ./internal/risk   | grep wisp/internal
internal/winsec : internal/secret : internal/observe : internal/risk(本体)
$ go list -deps ./internal/plugin | grep wisp/internal
internal/winsec : internal/secret : internal/observe : internal/plugin(本体)
```

- **`internal/risk` 现在不依赖 `internal/plugin`。** 它的模块内闭包只有 `winsec`/`secret`/`observe`。
- **反向也不成立**：`grep -rn "wisp/internal/risk" internal/plugin/` → **零命中**（含 `*_test.go`）。
- 唯一的跨包引用是**测试方向**：`internal/risk/provenance_test.go:11` import 了 `"github.com/CarlosShao/wisp/internal/plugin"`（`:307 :310` 在用 `plugin.NewDisposalScope`/`scope.Defer`）。
- 派单 §4 提醒的"两枚不同 module 会给不同图"：本仓另有两个独立模块 `tools/d22scan`、`tools/mockllm`（还有 `scripts/spike`）；上面两条读数都在**主模块**内问的，`internal/risk` 与 `internal/plugin` 同属主模块。

**结论句（一句话）**：依赖方向**既不是 risk→plugin 也不是 plugin→risk，而是两向都无边（同级兄弟）**，所以"开的时候就登记怎么关"**在 risk 层做得到**——让 `internal/risk` 新增一条指向 `internal/plugin` 的边不构成环（plugin 的闭包里没有 risk），而且 risk 的**测试**今天已经把 plugin 编进去了，加这条边对测试图的形状没有净变化；派单里"若 plugin→risk 就只能由 bridge 登记、于是许诺只覆盖一条腿"那一支**前提不成立**。

**它让票面 AC#2 的哪半句必须改写（只登记，不改票面）**：要改写的不是依赖方向那半句，是**同一行里的"关闭动作挂到该任务的 `DisposalScope` 上"**。现量：
- `grep -rn "DisposalScope" internal/tools internal/agent cmd/wisp --include="*.go" | grep -v _test.go` → **只命中一枚注释**：`internal/agent/approval/gate.go:150`。**生产码里没有任何一枚 `*plugin.DisposalScope` 走到桥、走到 loop、走到组合根。**
- 今天的关腿长在这儿：`internal/agent/loop.go:360-366`（`:360 revokeAdmission := func() {}`／`:361-363` 取 `AdmitTask(taskID)` 交回的 `r` 存进 `revokeAdmission`／`:366 defer revokeAdmission()`）→ `cmd/wisp/run.go:586 AdmitTask: rt.admitTask` → `cmd/wisp/run.go:557-565` 的 `revoke`（`:563 rt.bridge.CloseTask(taskID)`）。
⇒ 也就是说：**"任务的 `DisposalScope`"在这条腿上今天不存在**，`OpenScope` 拿不到它。票面 AC#2 那一半要么改写成"挂到该任务的 **admitTask revoke 边界**上（今天的真实持有者）"，要么就得**先补一条任务级 DisposalScope 的管道**——后者射程会冲出 `provenance.go`＋两枚调用点（至少 `internal/tools/bridge.go`、`internal/agent/loop.go`、`internal/agent` 的 `Options.AdmitTask` 签名、`cmd/wisp/run.go` 四枚），**超出 owner 已批的批准单位**（台账 `A322` 口径）。**这是本程要停手上报的第一枚未定义项。**

### 2.3 两条设计支各自要动哪几枚（＋各自"摘掉认身份"会不会红）

先给一支**共同的、会咬人的读数**，因为它决定两支的"写不出来"到底覆盖到哪：

> **Go 不强制使用函数返回值。** 台件 `.scratch/wisp/probes/160/c1/shape160.go` LEG 3 现量：`e3.OpenB("panel-assets-l2", &disposal{}, 1)` 作为一条语句**丢弃句柄**，`gofmt -l` 空、`go vet ./...` 过、`go run .` 打出 `B-discarded scopes=1`。
> ⇒ **两支都不会把"忘了关"变成编译错误。它们消灭的是另一件事**："`CloseScope(任意字符串)` 从今天起没有可达的写法"。票面 `:12` 那句"忘了关从'违反纪律'变成'这种代码写不出来'"**按现量必须降格**为"关的动作从此离不开开它的那一枚"；"只开不合"照样写得出来、照样编译过（`panel_assets.go:232` 那一枚的形状在新名下依然合法）。
> 这条降格与票面 AC#1 那句"编译器不拦"是**同一件事的前后态**，代码程写 AC#1 判据时不要把它写成新形状已经治了。

**甲（关的动作是闭包里的局部量、不出包）** —— 逐枚要动：

| # | 文件 | 位置 | 为什么必须动 |
|---|---|---|---|
| 1 | `internal/risk/provenance.go` | `:339-345` `OpenScope` | 签名要吃下"怎么关"的登记对象（plugin.DisposalScope 或 risk 本地声明的一枚 `Defer(func() bool) bool` 接口），返回值改为 void 或句柄 |
| 2 | `internal/risk/provenance.go` | `:347-354` `CloseScope` | 甲形的定义就是**不出包**＝这枚导出方法要变成非导出 `closeScope`；导出面一收，`:310 :358 :364 :181`（risk 自己的 4 枚测试调用）**不会红**（同包），但 `internal/tools/bridge.go:695` **必红**（跨包调用不存在的方法＝编译失败） |
| 3 | `internal/tools/bridge.go` | `:633-644` `OpenTask` | 开腿要把任务的持有者交给 `OpenScope`；桥今天没有 `*plugin.DisposalScope`（§2.2）⇒ 只能把 `CloseTask` 继续当关腿，甲形在这儿**塌回成"两个函数"** |
| 4 | `internal/tools/bridge.go` | `:684-699` `CloseTask` | `b.prov.CloseScope(taskID)` 拿不到了，只能改成"由开的时候交回的东西去关" ⇒ 桥要开始**存那枚东西**（`b.scopes` 的 `map[string]bool` 形状不够，`:638-639`） |
| 5 | `cmd/wisp/panel_assets.go` | `:227-237`（尤其 `:231-232`） | 只开不合的那枚；`detector()` 里没有 disposal 持有者，甲形在这儿**只能就地丢弃**或加 `defer`（进程内探针，语义上可接受，但要写清） |
| 6 | `internal/risk/*_test.go` 22＋4 枚点 | 见 §2.1 名单 | 改签名（多一个参数）会让 22 枚开点**全部编译失败**；只导出面变化只伤 4 枚关点（且同包不伤） |

甲"摘掉认身份"这一味会不会有跨包误关打不红：**甲天然没有"摘"这个动作**——关是非导出的、且不接外部 id，所以**跨包根本写不出发射状的关**；代价是**合法持有者也关不了**（`bridge.CloseTask` 那条腿必须靠开时登记的闭包，而那要求桥在开的时候就把持有者拿到手——§2.2 的管道缺位在这儿变成硬墙）。**残余风险**：甲下"摘掉认身份"的等价变异是"把 `closeScope` 改成导出并接 id 字符串"，`grep` 尺（G5 只数调用者、排 `internal/risk/*`）**看不见**，`internal/risk` 自己的 22 枚测试点也**不会红**（它们不测"外人能不能关"）。

**乙（`OpenScope` 交回一枚带 `Close()` 的结构体）** —— 逐枚要动：

| # | 文件 | 位置 | 为什么必须动 |
|---|---|---|---|
| 1 | `internal/risk/provenance.go` | `:339-345` | 返回 `*Scope`（或值类型），内部仍登记 `p.scopes[scopeID]=nil`；`Scope.Close()` 的 id 是**私有字段**（台件 `ScopeHandle` 即此形） |
| 2 | `internal/risk/provenance.go` | `:347-354` | 导出 `CloseScope` 若**保留**＝"只开不关"与"外人关"都还在；若**收非导出**，跨包关只有句柄一条路（乙的本意） |
| 3 | `internal/risk/provenance.go` | `:477`、`:578` 两枚 `Origin` 串＋`:453` 的日志句 | 文案点名 `OpenScope`/`CloseScope`；**是字符串不是注释**，G5 尺会读它们（§2.1 末） |
| 4 | `internal/tools/bridge.go` | `:633-644` | 桥要**持有**那枚句柄（`b.scopes` 现在是 `map[string]bool`，`:638-639`）⇒ 容器类型与幂等逻辑一起改 |
| 5 | `internal/tools/bridge.go` | `:684-699` | `b.prov.CloseScope(taskID)` → `h.Close()`；`:688` 的 `dropped`、`:697-698` 的审计行 `was_open/dropped/open_scopes` 是**外部可见读数**（见 §2.5） |
| 6 | `cmd/wisp/panel_assets.go` | `:231-232` | 一句 `prov.OpenScope(taintSourceScopeID)` → 拿到句柄就地丢弃或 `defer h.Close()`；**丢弃仍编译**（台件 LEG 3） |
| 7 | `internal/risk/*_test.go` 22 枚开点 | 名单见 §2.1 | **只有当参数形状变**才会红；单纯多一个返回值时 22 枚**全部照绿**（丢弃合法）。4 枚关点在 `CloseScope` 收非导出后仍**同包可编**，但 `:310` 那一枚（`scope.Defer(func(){ p.CloseScope(...) })`）是**票面契约的化身**，改它＝改 AC#5 的测试射程 |

乙"摘掉认身份"这一味会不会有跨包误关打不红：**会有，而且今天就没有一把尺量它。** 台件 LEG 4 现量：把关退回成 `CloseAny(scopeID string)` 之后，外人对同一 id 一发 `CloseAny("task-victim")`，`B-no-identity mis-closed scopes=0`——**被误关的那一轮从表上消失，形状与"它自己关过"完全一致**。落到真码上，`provenance.go:441-458` 的 `scopeMarks` 是唯一会响的东西，而它只在**别的 scope 还带着污点**时才 fail-closed（`:451-456` 那个 for 循环）：
- 受害者 scope **带污点**被别人误关 ⇒ 受害者下一次 `Inspect` 走 `:448` 的 `known` 失败 ⇒ `:453` 记一行日志 ⇒ `SrcUnboundScope`（`:432`）R4 升级 ⇒ **响亮**；
- 受害者 scope **开着但空**（读了不该 `Mark` 的东西，或 `Mark` 在误关之后还没跑到）⇒ `:457` `return nil, false` ⇒ **静默放行**，编译器、门、测试三头全都不响。
⇒ **AC#3 的"打不红"那一问，答案是"存在这样一发"（第二支形），档位＝台件读数＋真码分支枚举，不是推断。** 两支共同缺的是"这一发恰好没被 `Mark`"的窗口，甲形因为"外人连 id 都拿不到"把这一窗口**关小但不是零**（它靠 `bridge.scopes` 的 `was_open` 记账，见 §2.5）。

### 2.4 同源拷贝枚数（逐枚现量行号）

被抄的那把"尺"＝`internal/risk/provenance.go:55-61` 的 `DEFERRED(C25-loop-wiring)` 契约注释（`grep -n "DEFERRED(C25-loop-wiring)" internal cmd tools docs` 现量）。**自称同源、复述同一条纪律的，5 枚**（逐枚行号＋原文要点，全部在锚 bb679c6＝dev 现量存在）：

| # | 位置 | 复述的是什么 |
|---|---|---|
| 1 | `internal/risk/provenance.go:16-17` | "scopes are opened by the composition root and **Closed on Dispose**, so a new session never inherits old taints"——**Dispose 今天在生产码里没发生过**（§2.2） |
| 2 | `internal/risk/provenance.go:347-349` | "wire via `disposalScope.Defer(p.CloseScope)`"（CloseScope 自己的文档） |
| 3 | `internal/risk/provenance.go:436-440`＋`:453` | 注释版＋**日志版**："Inspect/CheckText on unregistered scope … **(ensure OpenScope at task start)**" |
| 4 | `internal/tools/bridge.go:626-632` | "ticket 19's **DEFERRED(C25-loop-wiring) item (1)** … the close side is per TASK and lives in **another package**, so nothing here keeps the two sides in step by construction" |
| 5 | `cmd/wisp/panel_assets.go:161-165` | "…(ticket 19's DEFERRED(C25-loop-wiring) **item 1**); a CLI probe has no task, so it names one and **closes nothing** - the engine lives and dies inside this process" |

另有 3 枚**弱同族**（复述同一条纪律但不点名标记，改尺时同样要一起看）：`internal/tools/bridge.go:646-683`（CloseTask 长注，逐字描述哪几腿没 owner）、`internal/risk/provenance_test.go:305-306`（测试注释）、`internal/tools/bridge_scope_open_ticket158_test.go:19-24`（注释里点 `provenance.go:394`）。
⇒ **改一把尺要一起改的是 5 枚点名拷贝（＋3 枚弱同族＝8 枚）**。派单说的"那几句复述"＝4 组，**数是 5 枚**（第 1 枚 `:16-17` 是尾项，分组结构天然漏的那一枚）。
⇒ 拷贝里带行号自引的只有一处成文引用：`docs/reports/pending-and-issues.md:7660` 逐字引 `internal/risk/provenance.go:55-61`（本程未改台账）。
⇒ `grep -rn "OpenScope" docs/PLAN.md docs/specs/` → **零命中**：C25 这句话在 PLAN/specs 里没有副本，尺只活在码里。

### 2.5 测试面预演（AC#3 那一问的外部可见读数）

**"开必须配关"那类用例现量：`internal/risk` 里 2 枚直接钉它 ＋ `cmd/wisp` 里 2 枚钉它的腿**（`grep -n "^func Test"` 逐枚）：
- `provenance_test.go:304 TestDisposalScopeClearsTaints` —— 唯一把契约**写成机器**的一枚：`:307` 建 `plugin.NewDisposalScope("session-1", …)`、`:309` 开、`:310` `scope.Defer(func(){ p.CloseScope("session-1") })`、`:316` `Dispose()` 后 `:319-328` 要求污点掉干净且新会话不继承。
- `provenance_test.go:331 TestInspectUnknownScopeIsEmptyStore` —— 钉的是**误关/未开的后果**：`:358-367` 关完再开再关，`:337-356` 钉 `SrcUnboundScope` 三条路径（`Inspect`/`Detector`/`CheckText`）。
- `provenance_test.go:295 TestScopesNeverInherit` —— 钉并发隔离，不钉成对。
- 腿上的守卫在**别的包**：`cmd/wisp/task_scope_close_151_test.go:54 TestCompositionRootClosesTheLoopTasksTaintScope`、`:67 TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit`（读 CloseTask 审计行），`internal/tools/bridge_scope_open_ticket158_test.go` 三条判据（开账本 `was_open`／审计行 `dropped=1`／非敏感源不开账）。
⇒ 这 4 枚**跨包**用例是今天唯一的机器；`internal/risk` 那 2 枚测的是引擎自洽，**没有一枚测"外部能不能凭一个字符串关掉别人的 scope"**。
⇒ `go test -count=1 ./internal/risk ./internal/tools` 本机基线：**`internal/tools` ok**；**`internal/risk` FAIL 一枚**＝`TestResolvePerCallBudget`（`pathresolver_budget_norace_test.go:34` 量到 `1.200 ms/op`，预算 1.000 ms）——**与 160 无关的机器噪声，本程未动 `thresholds.go` 一字节**；`cmd/wisp` 本机整包量不到（§7）。

**AC#3 预演（摘掉"认身份"，有没有任何外部可见读数会变）——答得出，三条：**
1. `internal/tools/bridge.go:697-698` 的 `b.log("tools: C25 scope closed task=%s was_open=%v dropped=%d open_scopes=%d", …)` —— 摘掉认身份后 `was_open`/`open_scopes` 两格**照旧打**，**这行读数不区分"谁关的"**；`cmd/wisp/task_scope_close_151_test.go` 与 `internal/tools/bridge_scope_open_ticket158_test.go` 判据 2 都只读这行的 `was_open=true/dropped=1` ⇒ **两枚都照绿**。
2. `provenance.go:453` 那条 `logf`（"unregistered scope …"）——**只在别的 scope 带污点时才打**（`:451-456`）；受害 scope 为空时**静默**（`:457`）。
3. `provenance.go:477`/`:578` 的 `Origin: "scope is not open (OpenScope missing or already closed)"` 串——**"already closed"这一支语义就是"别人关过"**，误关与自关在**卡片与取证行上不可区分**（AC#3 的"外部可见"最强的一枚反例）。
⇒ **答语**：摘掉认身份后，**没有任何一枚现存的包内断言会红**；会红的只有 §2.5 那 4 枚跨包用例中"恰好被误关的那一枚 scope 当时带着污点且下一次 Inspect 发生在同一轮"的子集。所以**AC#3 的正控必须自己造**（未定义项：判据形状要代码程写进证据件，建议＝一发"外人凭 id 关别人"的最小用例，钉 `ScopeTaints` 与 `Inspect` 两把读数）。

### 2.6 `DEFERRED(C25-loop-wiring)` 的兑现面

- 码里现量**3 枚**（`grep -rn "C25-loop-wiring" internal cmd tools docs`）：`internal/risk/provenance.go:55`、`internal/tools/bridge.go:627`、`cmd/wisp/panel_assets.go:164`。派单未给数，本程数到 3（含定义处）。
- **`SPEC-12 §5` 登记表里那条不存在**：`grep -c "C25-loop-wiring" docs/PLAN.md docs/specs/SPEC-12-roadmap-governance.md` → **两枚文件都 0**。`SPEC-12` 里提到 C25 的只有 `:19`（S3 切片卡的能力项清单）与 `:86`（RESERVED"精确 taint tracking"），**都不是 DEFERRED 五字段条目**。⇒ "五字段"这一问的读数是**登记面零条目**，与 `AGENTS §1.1`"码里标记与 §5 登记表 1:1 双向"直接冲突。这一条**票 16x 已经量过并登记**（`docs/evidence/s1/16x-contract-lines-c1.md:175-206` 同样的两把尺、同样零命中），本程在**自己的锚上复算成立**，不新造账。
- **本票只收第 (1) 条**（票面 `:37`＋`:56-61` 的 (1)）：`OpenScope at task start and Defer(CloseScope) on the task's DisposalScope`。**第 (2)(3)(4) 条射程不在本票**：(2) 每枚敏感源输出要 `Mark`、(3) 出口要 `Inspect`/`CheckText`、(4) R4 来源名要进 `tool_call` 取证行且 `Hit.Fragment` 永不落库。⚠ **不许顺手一起兑现**（派单 §1.6）。
  现量这三条今天到哪儿了：(2)(3) 已由 `internal/tools/bridge.go` 的 `mark()`/`Execute`（`:272` 的 `b.assessorFor(req.TaskID).Assess(...)`、`:705-712` 的 `WithTaintDetector(b.prov.Detector(taskID))`）走通，(4) 是**另一枚票的账**（`Hit.Fragment` 永不落库这条纪律在 `provenance.go:60-61`）。

### 2.7 契约锚点复算（AC#2 引的两处行号，按 bb679c6＝dev 同 blob 现量）
- `internal/risk/assessor.go:228` = `func (a *RiskAssessor) Assess(tool string, params map[string]any, facts Facts) Decision {` —— **对**。
- `internal/tools/bridge.go:610` = `	f.Declared = entry.Decl.Declared` —— **对**（宿主在 `factsFor` 里把声明等级重新盖回 `Facts`，`:611` 同样重盖 `Paths`；这就是"重盖章"）。
- 票面 `:46-47` 引的 `bridge.go:638-642`（两本账）与 `provenance.go:386-388`（唯一可见差＝那行日志）**逐行核过成立**：`:638 open := b.scopes[taskID]`／`:639 b.scopes[taskID] = true`／`:642 b.prov.OpenScope(taskID)`；`:386 if _, ok := p.scopes[scopeID]; !ok {`／`:387` 那行 `logf`／`:388 }`。
- 票面 `:29` 引的 `provenance.go:603-612`（`Detector`＋`boundDetector`）**在位**。

---

## 3. 对编排者的不服（你给的未验证断言，逐条判）

| 你的断言 | 读数 | 判 |
|---|---|---|
| "git 分支 `dev`" | 检出的是 `dsh/feat/frontend-p0`，`dev` 是其祖先、`dev..HEAD=2`（两枚 `frontend/**`＋派单存档）；本程五个被测文件三锚 blob 全等 | **不成立但无害**：行号可搬到 dev。⚠ **代码程落地前要问清它写在哪一支** |
| "清扫清单在 `internal/plugin/disposal.go`" | `Defer:186`/`DeferNamed:191`/`Dispose:253`/`Incomplete:245` 四枚全在，另有 `Go:226`（派单没列） | **成立** |
| "`internal/risk` 现在依不依赖 `internal/plugin`？若 plugin→risk 则做不到" | 两向都无边（§2.2 `go list -deps` 双跑） | **前提不成立，结论翻转**：risk 层做得到，不依赖 bridge 代登 |
| "两枚唯一生产调用点" | 开 2 枚（`bridge.go:642`、`panel_assets.go:232`）＋**关 1 枚（`bridge.go:695`）** | **半对**：批准单位若按"两枚"字面理解，关腿没被点名 ⇒ 要 owner 补一句 |
| "那几句同源拷贝"＝4 处 | 点名拷贝 **5 枚**（尾项 `provenance.go:16-17`）＋3 枚弱同族 | **漏一枚** |
| "甲形参照 Pi `session.ts:1475-1490`" | 本仓之外，本程未读（派单 §2 未授权我去读外部仓） | **未核**，见 §1 |
| "AC#1 借票 158 的变异台件" | 台件在位（`probes/158/accept-r1/make-rig-a1.sh:13,32` 的 M3＝`sed 642s/b\.prov\.OpenScope(taskID)/_ = open/`） | **成立**，但它量的是"摘掉开"，量不到"忘了关" |
| 你没写但会咬人的 | **`cmd/wisp` 在本机 `exit status 0xc0000135`（STATUS_DLL_NOT_FOUND），整包零读数**；CloseTask 那条腿唯一的机器就在这包里 | 见 §7 |

## 4. 被拒／没成功的调用（取数前后）
- **零枚被权限系统拒的调用。**
- 失败/无效 3 次（都已换招，无一是产品码）：
  1. `grep -n ... ./internal/risk | grep -n "wisp/internal"` 那一发我误加了 `-c`/管道顺序，第一条读数只回了 `rc` 行 ⇒ 改跑 `go list -deps ./internal/risk | tail -30` 取到全图。
  2. 全仓 `Grep "OpenScope|CloseScope"` 命中 193.5KB 被截断 ⇒ 改用 `glob: {internal,cmd,tools}/**/*.go` 重取（生产/测试分层见 §2.1）。
  3. `go run .` 第一次报 `package probe160.example/shape is not a main package`（我把包名写成 `shape160`）⇒ 改 `package main` 后成功。
- 环境缺件：`gofumpt` 本机没有（§1 末条）。

## 5. 有没有跑过删除命令
**没有。** 全程零 `rm`/`del`/`unlink`/`git clean`/`checkout .`/`stash`/`reset`/`rebase`/`--amend`；台件**只建不删**（`probes/160/c1/` 两枚新文件）。`go test` 与 `go run` 产生的构建缓存不属仓库文件，未做清理。

## 6. 伪授权两栏计数
- **本程没有把任何一句话当授权来扩射程。** 我按派单只写的两枚面：`docs/evidence/s1/160-design-core-c1.md`（本件）＋ `.scratch/wisp/probes/160/c1/**`。
- 出现的"批准"字样两处，**都不是本程可自取的**：
  (a) 派单 `:3` 转述的 owner 09-27 09:30 批地（台账 `A322`）＝**已给的射程**，我只在里面读、不在里面写；
  (b) 票面 `:19` 标题"我的推荐（**owner 一句话就能定，我不代拍**）"＝**未给的射程**，本件因此不含选型结论。
- 另有 1 枚**自我升级风险**登记给编排者：§2.2 末尾"必须先补任务级 DisposalScope 管道"这句话**读起来像**"那就不挂 DisposalScope、改挂 revoke"——**那是契约变更**（改 AC#2 的一整半句），**要人工批准**，本程只登记不裁。

## 7. 凭据值抄录
本程未读到、未抄录任何 API 密钥／token／DPAPI 明文（`internal/secret` 只在 `go list -deps` 的图里出现过名字，未开文件）。

## 8. 档位（四档）
- **〔本机读数〕**＝§0 全部、§2.1 全部行号与枚数、§2.2 两次 `go list -deps`＋两处 grep、§2.4 5 枚拷贝、§2.5 基线（`internal/tools` ok／`internal/risk` 1 枚性能噪声红／`cmd/wisp` 0xc0000135）、§2.6 两把 `grep -c` 尺、§2.7 四处契约锚点、台件 `go run` 六行输出。
- **〔读码自取·未试改〕**＝§2.3 两张"要动哪几枚"表（**没有真的改过码去验证**），以及"甲形塌回两个函数"那一句。
- **〔推断·已标〕**＝§2.3 末"受害 scope 为空 ⇒ 静默放行"这一发的可达性（分支枚举得出，未构造运行）。
- **〔转抄·他人锚〕**＝§2.5 里"116/80/0/0 与 144/84/0/0"这类票 158 的包内基线**本程未复算、也未采用**（只用了票面 `:46-47` 的行号并按 bb679c6 复算，见 §2.7）；`docs/reports/pending-and-issues.md:7660` 的逐字引用是台账原文，未改。

---

## 9. `next=` —— 代码程派单就绪包

**先决（三枚未定义项，没答之前不许动 `internal/**`）：**
1. **`OpenScope` 从哪儿拿到"任务的持有者"**？生产码今天**没有任何一枚 `*plugin.DisposalScope` 走到 `internal/tools`/`internal/agent`/`cmd/wisp`**（唯一命中＝`internal/agent/approval/gate.go:150` 的注释）。三条出路：(a) 补管道（射程出批准单位）；(b) 改 AC#2 半句为"挂 admitTask revoke 边界"（**契约变更＝人工批准**）；(c) 甲/乙的登记面退回"调用方持有"（票面 `:12` 的许诺随之降格）。**要 owner 拍。**
2. **批准单位里"关腿"没被点名**：`bridge.go:695` 是第三枚生产调用点（§2.1）。派代码程前要 owner 把话补全。
3. **落地分支**：任务书写 `dev`，现场检出 `dsh/feat/frontend-p0`（dev＋2 枚 frontend）。代码程**必须自己确认基线**，别照抄本程锚。

**文件级写面名单（最小）**：
- `internal/risk/provenance.go`：`:339-345`（开）／`:347-354`（关）／`:16-17 :436-440 :453 :477 :578`（同源文案与串）
- `internal/risk/provenance_test.go`：`:304-329 :331-368 :295-302` ＋**22 枚开点/4 枚关点**（全名单在 §2.1）
- `internal/tools/bridge.go`：`:633-644`（含 `:638-639` 那本账）／`:684-699`（含 `:697-698` 审计行）／`:626-632 :646-683` 注释面
- `cmd/wisp/panel_assets.go`：`:227-237`（含 `:231-232`）＋`:161-165` 注释面
- 只在写面变更时才动：`internal/risk/taintmatch_test.go:176-190`、`internal/risk/syncdirs*_test.go`（4 枚点）、`internal/tools/bridge_scope_open_ticket158_test.go`、`cmd/wisp/task_scope_close_151_test.go`
- **禁改**：`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、`tools/d22scan/**`、票面、`docs/reports/**`、`frontend/**`、`design/**`

**每格可复算的"会响"判据草案**（都要在未修码上先量一次"该红的现在不红"）：
- **AC#1**：最小生产形状＝把 `panel_assets.go:232` 那一形照抄进台件（本程已量：`B-discarded scopes=1`，`go vet` 过）。判据＝"编译器不拦／G5 尺（`probes/154/gate-clauses.sh:205-206`）只点名文件不点名这一发／只有 G5 名册里 `cmd/wisp` 那一行"。**量不到就停手，别改判据。**
- **AC#2**：`go test -count=1 ./internal/risk ./internal/tools` 基线**必须先在健康 bench 重取**（本机 `internal/risk` 带 1 枚性能噪声红、`cmd/wisp` 整包 0xc0000135）；`.github/workflows/ci.yml:136-141` 的 `gofumpt -l . tools/d22scan tools/mockllm` 必须仍空。
- **AC#3 正控（**今天不存在，要新造**）**：一发"外人凭 id 关掉别人的 scope"的用例，钉两把读数——受害者**带污点**时 `Inspect` 出 `SrcUnboundScope`（`provenance.go:432/453`）、受害者**为空**时（`:457` 分支）**今天静默**；判据＝"摘掉认身份后这发不红＝红给我看"。
- **AC#4**：`Detector(scopeID)`（`:605-612`，`boundDetector:609-612`）不得被加 `Close`；`provenance.go:387` 那条"Mark 到没开的 scope 上照样存"的 fail-closed 语义不得反向。
- **AC#5**：`git show --name-only <号>` 名册必须 ⊆ 上面的写面名单；`git grep -n "DEFERRED(C25-loop-wiring)"` **仍恰 3 枚**（`:55 :627 :164`）——甲乙两支都不该减少枚数，减少＝顺手兑现了 (2)(3)(4)。

**必须回报的未定义项**：上面先决 1/2/3 ＋ §2.5 末"AC#3 正控判据由谁写" ＋ §6(b) 那句"要不要改 AC#2 半句"。
