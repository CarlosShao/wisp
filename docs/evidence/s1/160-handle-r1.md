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
