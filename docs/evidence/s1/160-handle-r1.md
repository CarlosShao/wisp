# 票 160 · 160-r1 —— `OpenScope` 交回自带关闭动作、认身份的句柄：三格读数

- 派单：`.scratch/wisp/dispatches/2026-09-27-105x-impl-160-r1-handle-with-identity-close.md`（重派版；**本程是 160-r1 的第二次起步**，前一程 09-27 10:57 派出、11:01 工作树被切分支抽掉派单件，编排者现量确认它一行未写盘）
- 本程时刻：09-27 11:17 起手
- 票面：`.scratch/wisp/issues/160-scope-open-returns-a-handle-carrying-its-own-closer.md`
- 非实现者前件读数：`docs/evidence/s1/160-design-core-c1.md`（代号 160-c1）
- 台件：`.scratch/wisp/probes/160/r1/**`（**只建不删**；三枚 go 模块 + 一份 gate 读数）

---

## 0. step-0 四件（本程自己现量，不抄任何人）

1. `date` → `2026-09-27 11:1734 +0800`。
2. **分支名（本程特有硬要求 #1）**：`git rev-parse --abbrev-ref HEAD` → **`dev`** —— 对，起步放行。
   本程全程跑到的分支名逐次登记：11:17 `dev`｜11:24 `dev`（每次 git 命令前复看）。
   ⚠ 160-c1 那份读数的 §0 第 2 条量到的是 `dsh/feat/frontend-p0`——**那是它的锚**，本程的锚是 `dev`，两者不同**不影响它的行号**（它自己 §0.2 已用 blob 全等证明可搬）。
3. `git rev-parse HEAD` → `c1b1008913de86a96bd57c6ba117d0e8e39f9636`。
   这条**推翻 160-c1 的起手锚** `bb679c6`，也推翻派单里任何按 `bb679c6` 写的号——本程一律用自己的现量。
4. 写面洁净度（四条全空）：
   `git status --porcelain -- internal/risk/` → 空
   `git status --porcelain -- internal/tools/` → 空
   `git status --porcelain -- cmd/wisp/` → 空
   `git status --porcelain -- internal/tools/bridge.go cmd/wisp/panel_assets.go`（派单 §0 原形）→ 空
   全树另有 46 行噪声（`design/**` 未提交删除、`probes/152/my152.py`、`docs/evidence/s1/152-…-accept-r1.md` 等）——**本程一律未碰**。

### 0.1 票面 AC 连行号（`dev` 现量，逐字抄条目名）

| 行 | 条目 |
|---|---|
| `:26` | **AC#1** 先把"改前写不出来"证成读数：未修码上量"编译器不拦、门不响、只有普查尺点名"。量不到就停手 |
| `:27` | **AC#2** 换形状：`OpenScope` 交回句柄；`Close` 认身份＋幂等；关闭动作挂到该任务的 `DisposalScope` 上（后半句被 `:28-33` 的 09-27 更正改判） |
| `:34` | **AC#3** 反向判据：摘掉"认身份"这一味，是否存在一发跨包误关从此打不红？摘掉它有没有任何外部可见读数变过？ |
| `:35` | **AC#4** 明写它治不到什么（本单**不含**此格，见下） |
| `:36` | **AC#5** 契约轴（本单**不含**此格） |

⚠ **本件要报的第一笔过期**：160-c1 §0.5 把 AC#3/AC#4/AC#5 记成 `:28`/`:29`/`:30`。本程现量是 `:34`/`:35`/`:36`——**它读的是 `>` 更正块插入之前的票面**，整体后移 6 行。AC#1`:26`／AC#2`:27` 两枚未受影响。

⚠ **本程范围＝派单明写的三格（AC#1／AC#2／AC#3）**。票面 `:35` AC#4、`:36` AC#5 **不在本单**：AC#4 要答的两问（从不动的泄漏／`Detector` 借用形）由后续裁决程做，本程不代勾票面。

### 0.2 派单与 160-c1 的每条前提，本程逐枚复验（`dev` @ `c1b10089`）

| 出处 | 前提 | 本程现量 | 判 |
|---|---|---|---|
| 派单/任务书 #2 | `dev` tip 已含票 162 r1＋r2（`fs_edit.go`/`fs_edit_test.go`/`fs_edit_ac34_test.go`） | `git ls-files internal/tools \| grep fs_edit` → 三枚全在 | **成立** |
| 任务书 #2 | 已跟踪 `.go` 计数＝545 | `git ls-files "*.go" \| wc -l` → **545** | **成立** |
| 派单 §3 | 签名 `provenance.go:339 OpenScope(scopeID string)`／`:350 CloseScope(scopeID string)` | `git show HEAD:… \| grep -n` → **`:339`／`:350` 逐字对** | **成立** |
| 派单 §3 | 生产调用点三处＝`bridge.go:642` 开／`bridge.go:695` 关／`panel_assets.go:232` 只开不合 | `git grep -n "OpenScope\|CloseScope" HEAD -- internal cmd tools` 排 `_test.go` → 调用点恰这 3 枚（其余 8 行是定义体注释/日志/串） | **成立** |
| 派单 §3 | risk 包内测试面开 22 枚／关 4 枚 | `git grep -c "\.OpenScope(" HEAD -- 'internal/risk/*_test.go'` → 12+3+3+2+1+1＝**22**；`\.CloseScope(` → 3+1＝**4** | **成立** |
| 派单 §3 | 签名连带要改的测试文件＝6 枚 | 排定义本体与两枚生产文件后剩 6 枚：`provenance_test.go`/`taintmatch_test.go`/`syncdirs_redteam_windows_test.go`/`syncdirs_test.go`/`syncdirs_windows_test.go`/`provenance_syncdirs_other_test.go` | **成立** |
| 派单 §1.24 | 同源点名拷贝 5 枚 | 见 §4 的逐枚改口表，5 枚全在位 | **成立** |
| 派单 §1.14 | 码里 `DEFERRED(C25-loop-wiring)` 恰 3 枚 | `git grep -n` 排 `docs/` → `provenance.go:55`／`bridge.go:627`／`panel_assets.go:164`＝**3**（`docs/**` 里另有 20 枚转述，不属"码里标记"） | **成立** |
| 派单 §2 | 160-c1 LEG 3"丢弃句柄照样过 vet、照样泄漏" | 本程**两路复现**：① 跑它自己的台件 → `B-discarded scopes=1`；② 用**真引擎**造同形（`probes/160/r1/ac1/`）→ 见 §3 | **成立（且本程把它从台件升级成真码读数）** |
| 派单 §6 | `internal/tools` 基线须自己现量、别抄 86 | 本程改前 `go test ./internal/tools/ -count=1` → **ok 16.236s**，枚数见 §6 | 已现量 |
| 派单 §5 | `internal/risk` 有一枚计时噪声红 `TestResolvePerCallBudget` | 本程改前 `go test ./internal/risk/ -count=1` → **ok 5.674s（这一趟没红）**——它是计时噪声，**同一台机器两次读数可以一个红一个绿**；见 §6 的逐名归类 | **成立但本程未复现红** |

⇒ 派单与 160-c1 的前提**本程全部复核成立**，唯一过期的是它引的票面行号（§0.1）。

---

## 3. 第①格＝AC#1 —— "改前写不出来"证成读数（未修码 @ `c1b10089`，`internal/**`＋`cmd/**` 一字节未动）

**先交代为什么不能用 160-c1 那台件交差**：`probes/160/c1/shape160.go` 是**重写出来的替身**（自带 `Engine`/`disposal`/`ScopeHandle`），它量的是"这两种形状在 Go 里编译不编译得过"，**没碰过 `risk.Provenance` 本体**。派单 §2 要的是"最小**生产**形状"，所以本程另造一台件，**直接 import 真引擎**：

```
.scratch/wisp/probes/160/r1/
├── go.work                    # 工作区：把真模块 + 两台件挂在一起
├── ac1/ac1open.go             # LEG 1/2：照抄 panel_assets.go:232 那一形
├── ac1-closebyid/closebyid.go # LEG 3：外人按 id 关别人的 scope
├── ac1-readings-pre.txt       # 改前读数（vet 与 LEG 3 输出）
├── gate-pre-change-HEAD.txt   # 改前门响全份（G1–G7 + 聚合退码）
└── reproduce-160c1-leg3.txt   # 160-c1 台件在本程复跑的输出
```

### (a) `go build` / `go vet` 拦不拦？—— **拦不住**（命令与读数都在 `ac1-readings-pre.txt`）

| 腿 | 命令 | 读数 |
|---|---|---|
| LEG 1 只开不关 | `cd probes/160/r1/ac1 && go vet ./...` | **rc=0** |
| | `go build -o $TMP/ac1.exe ./...` | rc=0（`go run .` 同一路径跑通） |
| | `go run .` | `LEG1-open-and-leak inspectHit=true srcTool="web.fetch" marks=1` ← **污点还活着，scope 还开着** |
| LEG 1 形状 | `gofmt -l ac1 ac1-closebyid` | **空**（这台件本身就"写得出来、还排得整齐"） |
| LEG 3 外人按 id 关 | `cd probes/160/r1/ac1-closebyid && go vet ./...` | **rc=0** |
| | `go run .` | `foreign-by-id-close taintsBefore=1 taintsAfter=0 inspectHitAfter=false` ← **受害者静默消失** |

⇒ **两件事同时成立**：① 只开不关照样编译、照样 vet 过、照样泄漏（复现 160-c1 LEG 3 **成功**，本程还把它从替身升级成真引擎读数）；② 今天**任何人都能在别的包里按一个字符串关掉别人的 scope，且成功之后零读数**——`inspectHitAfter=false` 与"它自己关过"**不可区分**。这正是派单 §3 (iii) 要堵的那枚后门形状。

⚠ LEG 3 那一发同时**复现了 160-c1 §2.3 末〔推断·已标〕那一支**：受害 scope 被误关后 `scopeMarks`（`provenance.go:441-458`）走到 `:457 return nil, false`——`internal/risk` 里**没有任何别的 scope 带污点**，`:451` 那个 for 循环空转 ⇒ `:453` 那条 fail-closed 日志**根本没打**（本程输出里除了 `winsec: sealing path resolver installed` 一行 INFO 之外零行）。⇒ **档位从〔推断·已标〕升到〔本机读数〕**：这一发"摘掉认身份就打不红"是**真跑出来的**，不是分支枚举。

### (b) 今天有没有任何门响？—— **门只点名到文件，且改前基线全绿**

命令：`bash .scratch/wisp/probes/154/gate-clauses.sh HEAD > probes/160/r1/gate-pre-change-HEAD.txt`，rc=0。逐腿读数：

| 腿 | 声明 | 实测 | 名册本体 |
|---|---|---|---|
| **G5** `OpenScope↔CloseScope` 成对普查（排 `*_test.go`、排定义本体 `internal/risk/*`） | ring | **ring** | `UNPAIRED cmd/wisp/panel_assets.go (开方调用点=1)`，未成对枚数＝**1** |
| G5-正控（同尺再排 `panel_assets.go`） | quiet | quiet | 0 枚 |
| G5-负一负（不过滤） | ring | ring | — |
| **G6** `OpenTask↔CloseTask` | ring rev | ring | 反向未成对＝**1**（`cmd/wisp/run.go`） |
| **G7** `Defer↔DisposalScope` | ring rev | ring | 反向未成对＝**3**（`memory/retention.go`、`observe/goroutine.go`、`proc/shutdown.go`） |
| 聚合退码（票 161 AC#6①） | — | **0** | 每一腿的响都和自己的声明一致 |

⇒ **三句要点**：
1. G5 **确实响**，但它响的是**一个文件名**（`cmd/wisp/panel_assets.go`），不是我 LEG 1 那一发调用；**我的台件在 `.scratch/` 下、根本不在这把尺的射程（`internal/**/*.go cmd/**/*.go`）里** ⇒ 同一形状换个目录就看不见，这是票面"为什么值得做"那节"文件粒度"缺陷的**又一发实证**。
2. G6／G7 与 `OpenScope` 无关（它们盘的是另外两族），改前改后都不该被本票动到——§6 用 diff 证这一点。
3. **测试面也不响**：改前 `go test ./internal/risk/ -count=1` → `ok 5.674s`、`go test ./internal/tools/ -count=1` → `ok 16.236s`。**两包全绿的情况下，LEak 照样存在**。⇒ "只有普查尺点名"这句成立。

### (c) 事后能不能查出"哪个 scope 开着没关"？—— **查不出**

`risk.Provenance` 的导出面里**没有"列出当前开着哪些 scope"这一枚**（`git grep -n "^func (p \*Provenance)"` 现量，全 15 枚导出方法逐枚看）：唯一的按-scope 观察口是 `ScopeTaints(scopeID)`，**它要你先知道 id**。读数：

```
LEG2-outsight  taints(id known)=1   taints(id guessed)=0
```

⇒ 知道 id 才看得见，不知道 id 时**"没开过"／"开过又关了"／"泄漏中还空着"三态全塌成 0**。⇒ AC#1 (c) 的答案＝**不能**，事后只能靠 G5 那把文件粒度普查尺，而它连"哪个 id"都不答。

### 3.9 本格判据落定（写死，给后两格当尺）

**改前的"只开不关"这一形：编译器不拦（vet rc=0）／测试不响（两包全绿）／门只点到一个文件名（G5 名册 1 枚）／事后零观察口（`ScopeTaints` 要 id）。** 四把尺的读数全在案 ⇒ AC#1 **量到了**，不必停手。

---

## 4. 第②格＝AC#2 —— 换形状（句柄自带关闭动作＋认身份＋幂等）

### 4.1 新形状（`internal/risk/provenance.go`，逐枚现量行号在 §4.6）

```go
type Scope struct { p *Provenance; id string; closed bool }   // id 非导出
func (p *Provenance) OpenScope(scopeID string) *Scope          // 交回句柄；重开＝换主人，污点原样留着
func (s *Scope) Close() error                                  // 唯一的关路；认身份
func (s *Scope) ID() string                                    // 读 id 不给关的权力
var ErrScopeAlreadyClosed / ErrScopeNotOwner / ErrScopeNotOpen  // 三种拒绝各一枚
type scopeReg struct { marks []*taintMark; owner *Scope }       // p.scopes 的值类型
```

**"两本账塌成一本书"落在这儿**：`p.scopes` 今天一枚值同时装"这个 scope 的污点"与"当前谁有权关它"，一把 `p.mu` 守。旧的导出 `CloseScope(scopeID string)` **整枚删掉**——不是收非导出、不是留白名单，是**没有这条路**。

### 4.2 三条硬判据，各一发实测（都是真跑的输出，不是措辞）

| 判据 | 用例／命令 | 读数 |
|---|---|---|
| **(i)** 拿别人的句柄关别人的登记 ⇒ 失败且**外部可见** | `TestScopeCloseRefusesAnotherOwnersRegistration`（`provenance_test.go`，本程新增）：`stale` 开 → `fresh` 重开同一 id → `Mark` → `stale.Close()` | `--- PASS`：`stale.Close()` 回 `ErrScopeNotOwner`；`ScopeTaints=1`（**一枚没掉**）；`Inspect` 仍出 `SrcWebFetch` 的 R4 命中；`fresh.Close()` 仍成功。**三把读数全在函数之外**——只多打一行日志的形状过不了这一发 |
| **(i′)** 同一形落在**空**登记上（160-c1 §2.3 预测的静默那一支） | `TestScopeCloseRefusedOnEmptyRegistrationIsStillVisible` | `--- PASS`：全引擎零污点时 `stale.Close()` 仍回 `ErrScopeNotOwner`，且 `p.scopes` 里那一格**还在**——静默放行形被这枚判据钉住 |
| **(ii)** 同一枚句柄关两次 ⇒ 幂等，**第二次读数说清** | `TestScopeCloseIsIdempotent` | `--- PASS`：第 1 次 `nil`；第 2、3 次 `ErrScopeAlreadyClosed`（与"别人已经关了"的 `ErrScopeNotOwner`／"根本没这枚登记"的 `ErrScopeNotOpen` **分格**，正是 AC#3 抱怨的 `already closed` 混格被拆开）；spent 句柄重开新登记后仍回 `ErrScopeAlreadyClosed`；`nil` receiver 不 panic |
| **(iii)** 旧口子不留后门 | `.scratch/wisp/probes/160/r1/ac1-closebyid/`（改前**编译得过、跑得通、受害者静默消失**的那台件） | **改后不再编译**：`closebyid.go:32:4: p.CloseScope undefined (type *risk.Provenance has no field or method CloseScope)`（`ac1-readings-post-closebyid.txt`）⇒ 运行时侧的 companion 断言＝`TestScopeIDGrantsNoPower`（`--- PASS`） |

⇒ **(iii) 走的是"旧口子消失"那一支**，不是"留给已登记的调用方"：仓里今天没有任何一处能按 id 关 scope，`internal/tools` 的关腿只剩"开的时候拿到的那枚东西"（§4.3）。

### 4.3 三处生产调用点逐枚改口（派单 §3 点名的那三枚，无第四枚）

| # | 位置（改前行号） | 改口 |
|---|---|---|
| 1 | `internal/tools/bridge.go:642` 开腿 | `OpenTask` 现在把 `b.prov.OpenScope(taskID)` 交回的句柄**存进桥的账本**：`b.scopes`（`:143`）值类型 `bool` → `*risk.Scope`，构造点 `:186` 跟着改。**这是开腿的必要连带**（160-c1 §2.3 乙表第 4 行原话："容器类型与幂等逻辑一起改"），不是新开的第三枚地界。幂等判据从 `if !open` 变成 `if _, open := b.scopes[taskID]; open { return }`，同一把尺 |
| 2 | `internal/tools/bridge.go:695` 关腿 | `b.prov.CloseScope(taskID)` → `scope.Close()`，返回的 `error` 存进 `closeErr`，**并进审计行**：`tools: C25 scope closed task=%s was_open=%v dropped=%d open_scopes=%d close_err=%v`。⚠ 这条**格式变了**（追加一格 `close_err`），见 §4.7 的影响登记 |
| 3 | `cmd/wisp/panel_assets.go:232` 只开不合 | `prov.OpenScope(...)` → **`_ = prov.OpenScope(...)`**＋四行说明。行为一字未变（照样只开不合、照样泄漏），变的是它**必须写出自己在丢弃**——160-c1 LEG 3 与 §3 的读数说编译器拦不住这一形，所以这一枚交给门（G5），不假装被类型系统治了 |

### 4.4 五枚同源拷贝逐枚改口（派单 §1.24 的"必须逐枚改口"）

| # | 位置（改前） | 改口后说的是什么 | 状态 |
|---|---|---|---|
| 1 | `provenance.go:16-17` | 原："scopes are opened by the composition root and **Closed on Dispose**" | 改成"OpenScope 交回 `*Scope`，**是那枚句柄**被 Defer 到 session 的 DisposalScope 上关"——`Closed on Dispose` 那句今天在生产码里从不发生（160-c1 §2.4 第 1 行原话），不再假装它发生 | **改口** |
| 2 | `provenance.go:347-349`（`CloseScope` 自己的文档） | 原："wire via `disposalScope.Defer(p.CloseScope)`" | 函数已被删；新的 `Scope.Close` 文档写 `disposalScope.Defer(func() { _ = scope.Close() })`——**Defer 的是句柄，不是 id** | **改口** |
| 3 | `provenance.go:436-440`＋`:453` | 注释版＋日志版 | 注释版（`scopeMarks` 的文档）整段重写，明写"外人关不动 ⇒ 只剩'你自己关过'与'从来没开过'两种"；**日志版 `:453` 原文一字未动**（它说的还是那件事），因为它是 AC#3 三条外部可见读数之一，本程要它保持可比 | 注释**改口**／日志**保持** |
| 4 | `bridge.go:626-632` | 原："the close side is per TASK and lives in **another package**, so nothing here keeps the two sides in step by construction" | 改成：关侧现在**就在这本账里**（开腿存进去的那枚句柄）；同时**保留**旧句子里今天仍真的那一半——开／合仍是两次调用、忘了合照样泄漏、门是 G5 不是编译器 | **改口** |
| 5 | `panel_assets.go:161-165` | 原："a CLI probe has no task, so it names one **and closes nothing**" | 行为描述**保持原文**（它还是对的），后面追加四行：这一枚是**写明故意的丢弃**，且"关不掉别人的"从今天起是类型事实 | **改口（追加）** |

另有三枚**弱同族**（160-c1 §2.4 末登记，不点名标记）：`bridge.go` 的 `CloseTask` 长注（`:646-683`）**一字未动**——它逐枚描述"哪几腿没 owner"，那件事本程没改；`provenance_test.go:305-306` 测试注释、`bridge_scope_open_ticket158_test.go:19-24` 注释随 §4.5 的连带改动一起看过，**没有一处再自称逐字抄 `provenance.go:55-61` 而内容已经不符**。

### 4.5 签名连带失效的测试文件（逐枚列名＋只改形状、不削断言）

编译器现量：**只有 3 枚测试文件真的编不过**（派单引的"6 枚"是**碰到这两个词的文件数**，`OpenScope` 交回一个值时**丢弃是合法的** ⇒ 22 枚开点一枚都不必改）。逐枚：

| 文件 | 改了什么 | 断言强了还是弱了 |
|---|---|---|
| `internal/risk/provenance_test.go` | `:310` `scope.Defer(func(){ p.CloseScope("session-1") })` → 持句柄 `handle` 并 `handle.Close()`，**且新增 `t.Errorf` 检查这次关必须成功**；`:358`／`:364` → 新增 helper `mustClose(t, s *Scope)`；`:339`／`:359` 两枚开点改为捕获句柄；imports 加 `errors`；新增 §4.2 那 4 枚用例 | **变强**（旧写法关失败了没人知道，新写法会红） |
| `internal/risk/taintmatch_test.go` | `:176` 捕获句柄、`:181` → `h.Close()` 且 `b.Fatalf` 检查错误 | **变强**（benchmark 现在会因 refused close 而红） |
| `internal/tools/bridge_scope_open_ticket158_test.go` | `:73` `return b.scopes[taskID]` → `return b.scopes[taskID] != nil` ＋注释 | **等价**（同一谓词，"账本里有没有这一枚 task id"） |

⇒ 4 枚 `CloseScope` 点全部转成句柄形；22 枚开点**零改动**（实测：`go vet ./internal/risk/ ./internal/tools/ ./cmd/wisp/` rc=0，改前它先报的正是这 3 枚文件）。

### 4.6 `DEFERRED(C25-loop-wiring)` 枚数尺（派单 §1.14 要求的自带尺）

`git grep -n "DEFERRED.C25-loop-wiring." -- internal cmd tools` → **恰 3 枚**：
```
internal/risk/provenance.go:58   internal/tools/bridge.go:627   cmd/wisp/panel_assets.go:164
```
枚数没变、一枚没删。⚠ 两笔连带：**① `provenance.go` 那一枚的行号 55 → 58**（拷贝 #1 的改口多了两行），而 `docs/reports/pending-and-issues.md:7660` **逐字引着 `provenance.go:55-61`** ——那枚行号现在过期，`docs/reports/**` 是禁改面，本程不动，**报回请编排者补账**。② 第 (1) 条的文字按派单允许的"只许改成已落地"改了，且**只落地了形状那一半**：`Defer 到任务的 DisposalScope` 那一半在文中明写 STILL OPEN 并给理由（三层管道出批准单位），(2)(3)(4) **一字未兑现**。

### 4.7 本格改动的两笔连带影响（都是读数逼出来的，登记不藏）

1. **审计行格式变了**：`internal/tools/bridge.go:697-698` 追加 `close_err=%v`。改它的理由是判据 (i)"失败必须外部可见（不是只多一行日志）"——审计行是这条腿唯一跨包可见的口，不带上就把认身份的拒绝又变回包内私事。两枚现存读它的用例（`internal/tools/bridge_scope_open_ticket158_test.go` 判据 2、`cmd/wisp/task_scope_close_151_test.go` 两枚）**按 substring 匹配**，本程 `internal/tools` 全绿实测；`cmd/wisp` 那两枚受 DLL 环境影响，读数在 §5。
2. **G5 这把尺的合方词根失效**：G5 数的是"`OpenScope` 有没有配 `CloseScope`"（`.scratch/wisp/probes/154/gate-clauses.sh:205-206`），而 `CloseScope` 今天整枚消失 ⇒ 全仓合方名册变空。本程改前两次现量都是 1 枚（`gate-pre-change-HEAD.txt`／`gate-pre-change-HEAD-repeat.txt`），改后读数在 §6 现跑现填。**这一格是票 158／161 的地界（`probes/154/**` 不在本程写面），本程不改它，只把读数交回去。**

---

## 5. 第③格＝AC#3 —— 反向判据，两问都答（台件：`probes/160/r1/mut/`，**零改动 `internal/`**）

**台件怎么做到不动生产码**：`go test -overlay=probes/160/r1/mut/overlay.json ./internal/risk/` —— 变异体
`provenance.no-identity.go` 在编译期**替换** `internal/risk/provenance.go`，工作树全程
`git status --porcelain -- internal/ cmd/` **空输出**（跑完现量在案）。变异内容＝把
`Scope.Close` 里 `case reg.owner != s: return ErrScopeNotOwner` 那一味摘掉（其余全留），
就是 160-c1 LEG 4 那枚 `CloseAny(scopeID)` 形状放回真码。

### 第一问：摘掉"认身份"，是否存在一发跨包误关从此打不红？—— **存在，且改前全仓零枚尺吃得住**

| Build | `internal/risk` 全量（`-v` 逐名） | 结果 |
|---|---|---|
| **装上（出货形状）** | 103 PASS／1 SKIP／0 FAIL，104 枚 | `ok 4.776s` |
| **摘掉认身份（overlay 变异体）** | **只有 2 枚红**：`TestScopeCloseRefusesAnotherOwnersRegistration`（`provenance_test.go:913` `a handle that no longer owns the id must be refused, got <nil>`）、`TestScopeCloseRefusedOnEmptyRegistrationIsStillVisible`（`:941` `empty registration must still refuse a non-owner, got <nil>`）；**其余 101 枚照绿，连 SKIP 那一枚都没变** | `FAIL 4.804s` |

⇒ 那两枚正控**是本程造的**（改前不存在，160-c1 §2.5 的"正控今天不存在"这句被本程**实测复现**：同一变异体下，除这两枚之外**没有任何一枚现存包内断言会红**，包括 `TestDisposalScopeClearsTaints`／`TestInspectUnknownScopeIsEmptyStore`／`TestScopesNeverInherit` 这三枚最贴身的）。
⇒ 判据成色：**装上＝绿／摘掉＝红，且只有这两枚红** ⇒ 这一味**承重**，不是装饰。两枚红都读的是**函数之外的三把量**（`error` 返回值、`ScopeTaints` 账本、`Inspect` 判决），"只多一行日志"的形状过不了。

### 第二问：摘掉它，有没有任何**外部可见**读数变过？—— **没有一枚变过（跨包），这正是本味的承重证据**

同一变异体打在跨包口上：`go test -overlay=… ./internal/tools/ -count=1 -v` → **`ok 14.852s`，零 FAIL**（出货形状 `ok 18.897s`，两枚读数同名同状态）。⇒ 160-c1 §2.5 列的三条"外部可见读数"（`bridge.go` 审计行的 `was_open`／`dropped`、`:453` 那条条件日志、`Origin` 串里 `already closed` 的混格）**在摘掉认身份之后一枚都不动** —— 那三条今天**不区分谁关的**，这句从〔转抄·他人锚〕升成〔本机读数〕。
⇒ 唯一会动的跨包口是本程**新加**的 `close_err=%v` 那一格（变异体下它打成 `close_err=<nil>`），而它今天**没有任何 tools 侧用例钉着** —— 见 §7 残项 ①。

---

## 6. 门禁（三格全部落盘后跑，顺序照派单 §6）

| 尺 | 命令 | 读数 |
|---|---|---|
| risk 四数＋名册差集 | `go test ./internal/risk/ -count=1 -v` 逐名，基线用 `-overlay` 把六枚改动文件映回 `34c0b4e4` 复跑 | 改前 **99 PASS／0 FAIL／1 SKIP／100 枚**；改后 **103／0／1／104**；**"改前有改后无"＝0 枚**（`comm -23` 空），新增 4 枚全是本程点名的正控（`pre-base/roster-risk-{pre,post}.txt`、`roster-risk-lost.txt`） |
| tools | `go test ./internal/tools/ -count=1` | **ok 15.336s**（本程自己的基线：改前同一条命令 **ok 16.236s**；**未采用派单警告的 162-r1 那个 86**） |
| vet | `go vet ./internal/risk/ ./internal/tools/ ./cmd/wisp/` | **rc=0**（生产码 `go build ./...` 亦 rc=0） |
| cmd/wisp 健康 bench | `PATH=$PWD/third_party/sherpa-onnx:$PATH go test ./cmd/wisp/ -count=1` | **ok 115.178s** —— 换健康 bench 后**量得到**，且 `CloseTask` 那条腿唯一的机器两枚**逐名 PASS**：`TestCompositionRootClosesTheLoopTasksTaintScope (1.54s)`、`TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit (1.43s)`。⚠ 本机不接 PATH 时仍是 `0xc0000135`，**那是既有环境红，本程未把它算进自己头上** |
| d22scan 自测 | `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | **OK — PASS=34 FAIL=0 SKIP=0，`[no tests to run]`=0** |
| d22scan 主尺 | `sh scripts/d22scan.sh` | **rc=0**，`clean - no D22 ban violations` |
| gofumpt | `command -v gofumpt` → **NOT ON PATH**；`ls ~/go/bin` → 空 | **本程未跑**（与 160-c1 §1 末条同一台机器同一缺件）。代用尺＝`gofmt -l $(git ls-files "*.go")`（甲形，已跟踪全集合）→ **0 行**；`gofmt -l internal/ cmd/` → 0 行 |

### 6.1 计时红与真红逐名（派单 §5 要求）

`TestResolvePerCallBudget`（`pathresolver_budget_norace_test.go:34`，阈值 1.000 ms/op，`thresholds.go` 一字节未动）：
改前一趟 `ok 5.674s`、改后 `-v` 全量 **0 FAIL**，**本程四趟 risk 全量跑里它一枚都没红过**。
⇒ 归类＝**本程未复现的既有计时噪声**（160-c1 §2.5 记到的是 1.200ms 那一趟）；**未动阈值、未 Skip、未改名**；真红：本程 **0 枚**。
唯一 SKIP＝`TestSyncRegistryProbeLive`，**改前改后同一状态**（环境件，非本程引入）。

### 6.2 ⚠ 本程造成的一道门 regress（不在写面，报回不修）

`bash .scratch/wisp/probes/154/gate-clauses.sh HEAD` 改后 **rc=1**（改前 **rc=0**）：

| 腿 | 改前 | 改后 |
|---|---|---|
| G5 主尺 | UNPAIRED **1 枚**（`cmd/wisp/panel_assets.go`） | UNPAIRED **2 枚**：`panel_assets.go` ＋ **`internal/tools/bridge.go`（假阳性——它今天照旧在 `CloseTask` 里关）** |
| **G5 正控**（同尺再排掉 `panel_assets.go`，声明 quiet＝必须不响） | 0 枚，quiet | **1 枚，ring** ⇒ `BAD` |
| 其余 13 腿（G1/G1b/G2/G3/G4/G5neg/G6±/G7±） | 全部与声明一致 | **全部不变，仍一致** |
| 聚合退码 | **0** | **1** |

读数文件：`gate-pre-change-HEAD.txt`／`gate-post-change-HEAD.txt`／`gate-pre-vs-post.diff.txt`（206 行）／`gate-aggregate-post.txt`。

**根因**：`gate-clauses.sh:205-206` 那把尺的合方词根写死成 `CloseScope`，而 AC#2 判据 (iii) 选的支正是"**旧口子消失**"⇒ 合方词在全仓不再出现。
**为什么本程不顺手修**：改它＝改 `probes/154/**`，那枚文件**不在派单 §7 的逐枚写面名单里**（禁改名单也没列它，属"没点名"而非"点名禁改"），而票 154/158/161 的证据件按行 diff 这份名册——本程一改就把三张已验收票的读数面动了。
**给门 owner 的最小形状**（一句话，不需要新语义）：把 G5 的合方词根从 `CloseScope` 换成 `\bClose\b` 或 `Close Scope|scope.Close|Close()`，并把 `bridge.go` 那一枚 UNPAIRED 的"同文件另见收尾动词"注释升成判据；**不动 `want` 声明**。
⇒ **本票 AC#2 与门 G5 之间不是一致而是互斥**：换形状必然让 G5 的正控腿响。这一格**必须**由非实现者裁一次，本程不自取。

---

## 7. 残项登记（本程没做／做不到／不在批准单位，逐条指名）

1. **`close_err=%v` 那一格没有 tools 侧用例钉着**。要钉就得往 `internal/tools/bridge_scope_open_ticket158_test.go` 里**新增**断言，而 §7 写面只许本程对该文件做"签名连带的形状改动"（§4.5 那一枚 `!= nil`）。⇒ 交验收程判：留（＝跨包外部可见读数有一格没人守）或补一张票。
2. **`Scope` 非并发安全（同一枚句柄上 race 两次 `Close`）**：`s.closed` 是句柄私有量、在 `p.mu` 之内读写，但两枚**不同**句柄关**不同** scope 完全安全；同一句柄并发关会给出两次 `nil`。文档已写明，未加锁（加锁要引入 per-handle 状态＝碰 D38，不在本单）。
3. **`Mark` 造出来的无主登记（`scopeReg.owner == nil`）今天谁也关不掉**：这是 fail-closed 方向（丢污点＝开闸，`bridge.go:650-653` 原话同向），但它是**一枚新的、永久性的泄漏形**——改前 `CloseScope(id)` 还能收掉它。⇒ 属票面 `:35` **AC#4"明写它治不到什么"**那一格的射程，本单不含 AC#4，**只登记不裁**。
4. **G5 regress**（§6.2）。
5. **台账行号过期**：`docs/reports/pending-and-issues.md:7660` 逐字引 `provenance.go:55-61`，本程该段位移到 `:58-64`。禁改面，未动。
6. **票面一格没勾**：票面 AC#4（`:35`）／AC#5（`:36`）不在派单 §7 的三格里，本程**不代勾票面**（进度由编排者代落）。
