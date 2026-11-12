# 434 — `TestTicket198FirstRunCallerIsTheRunEntryOnly` 的第二调用者：判乙（断言该承认 `runConfigSet65`）

- 起手锚：`2026-11-12 22:54 +08`，`HEAD 23229545b`，分支 `dev`。
- 写面：**一枚文件** `cmd/wisp/firstrun_198_test.go`（断言侧）。**产码零改动**，`cmd/wisp/run.go` 一个字节没碰。
- 判语：**乙**。`runConfigSet65` 调 `ensureFirstRunConfig` 是对的，规矩被违反的是"**只有 `wisp run` 一个调用者**"这句**计数**，不是"**常驻腿不得触发首建**"这条**规矩**。

---

## 0. AC#0 只读取证（零写）

### 0.1 红句与耗时（carrier 现跑，逐字）

```
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestTicket198FirstRunCallerIsTheRunEntryOnly$'"
```

```
=== RUN   TestTicket198FirstRunCallerIsTheRunEntryOnly
    firstrun_198_test.go:281: AC J1 RED: ensureFirstRunConfig is called by [runConfigSet65 runTextTask], want exactly [runTextTask] - the resident leg reaches the assembly root directly, and first-run creation shared there writes a config.toml its pins forbid
--- FAIL: TestTicket198FirstRunCallerIsTheRunEntryOnly (0.44s)
FAIL
```

**耗时 0.44s**（本机本发；台账 `:32512` 记 0.24s、`:34135` 记 0.25s——同一枚纯 AST 检查，差别是机器负载，不影响判定）。纯 AST，不跑产品 ⇒ **与负载无关、可稳定复现**。

### 0.2 AC J1 原文逐字

票 198（`.scratch/wisp/issues/198-a-fresh-machine-cannot-run-wisp-because-nothing-creates-config-toml-done.md:13`）逐字：

> ⇒ **编排者裁定（J1，账 `A514`）＝首建只挂在 `wisp run` 这一条"用户明确发起"的入口；⛔ 不许进 `assembleRuntime`；⛔ 不许改那两枚钉**（"为了让自己落地而去放宽别人票的 AC 判据"就在本仓禁止清单上，不是可选项）。若 CLI 入口需要数据根先存在，**只许调用同一个现成判定者**（`secret.NewStore(dataDir)`，`internal/secret/store.go:49`，它今天就先于 `NewManager` 把根逐层建好并封权），⛔ 不许自己拼路径（`risk.PathResolver` 之外做文件系统决策＝`d22scan` 直接判违规）。

裁定原文出处＝台账 `A514`（`docs/reports/pending-and-issues.md:10469`）逐字：

> ★ **决定落点的那枚钉，我 09:0x 自己读到逐字**：`cmd/wisp/resident_task_source_246_windows_test.go:389-391` 断言常驻腿的 `config.toml` **必须保持不存在**（红句逐字 "the leg created a config.toml it was never asked for"），而 `assembleRuntime` 被常驻腿共用 ⇒ **首建放进 `assembleRuntime` 今天就红 3 处**；… ⇒ **J1 裁定＝首建只挂 `wisp run` 那一条用户明确发起的入口；⛔ 不进 `assembleRuntime`；⛔ 不改那两枚钉**。

**★ 判定的关键：原文的禁止对象是"进 `assembleRuntime`"，理由是"常驻腿直达装配根"。原文里"这一条入口"是**落点陈述**（首建该挂哪儿），不是"**恰好一个调用者**"这条独立规矩。断言把"落点陈述"硬化成了"计数＝1"，于是产品后来多了一条**同属用户发起**的入口时，断言红而规矩没破。**

### 0.3 哪个提交加的调用、它说了什么

`git log -S'ensureFirstRunConfig' -- cmd/wisp/config_cmd_65.go` ⇒ 唯一一枚：**`ef00174ac`**（2026-11-07 20:29:16 +0800，票 65 CLI 切球路径）。提交说明里与本事有关的那句逐字：

> 缺文件set走ensureFirstRunConfig唯一首建口,get只答默认不建文件

同提交新增 `cmd/wisp/config_cmd_65.go`（236 行，`--numstat` 逐字 `236 0`）——**这不是一次改动，是一个新入口从出生就带着这一跳**。它不是"某腿顺手多写了一行"，而是票 65 把 `wisp config set` 建成"机主今天就能用的写路径"时，**刻意复用**了票 198 的首建口，并把它写进了提交说明当作设计声明。

`git log --oneline -10 -- 'cmd/wisp/run.go'` 的最近十枚里，**没有一枚**给 `runConfigSet` 加过调用（那是 `config_cmd_65.go` 的事）；最近一枚 `7dcb145ff` 正是票 174 `[fs]` 授权面的回退（本票 AC#4 要求确认无关 ⇒ **无关**，见 §4）。

### 0.4 `runConfigSet65` 完整函数体与它为什么需要这一跳

`cmd/wisp/config_cmd_65.go:164-190`（逐字要点）：

```go
func runConfigSet65(cio configIO, dataDir, cfgPath, key string, rest []string) int {
	if key != configCliKey65 || len(rest) != 1 { … return 2 }   // 名册/用法拒绝，零写
	value := rest[0]
	if _, err := ensureFirstRunConfig(dataDir, cio.stderr); err != nil {   // :173 ← 争议这一行
		fmt.Fprintf(cio.stderr, "wisp config: %v\n", err)
		return 1
	}
	mgr, err := config.NewManager(cfgPath, nil)                            // :177
	if err != nil { … return 1 }                                           // 读不进来就退
	if _, err := mgr.SetBallStyle(value); err != nil { … return 1 }        // :182 守卫写
	…
}
```

**它为什么需要**：`config.NewManager`（`internal/config/manager.go:113-123`）第一句就是 `LoadFile(path, res)`，**文件不存在就返回错误**（`manager.go:111-112` 注释逐字 "first-run file creation is the caller's concern"）。所以不先建文件，`config set` 在干净机器上**必然**在 `:177` 退码 1。这不是"顺手"，是**这条路能不能走通的前置**。

**它为什么不算"常驻腿"那条害**：`runConfigSet65` 只从 `cmdConfig`（`config_cmd_65.go:121`）来，`cmdConfig` 只从 `main.go:149-157` 的 `case "config"` 来——**一条独立的 CLI 分发支，全程不碰 `assembleRuntime`**。它跟"常驻腿直达装配根"是两件不同的事。

---

## 1. 判定：乙，以及依据

**一句话**：J1 原文禁的是"**首建进 `assembleRuntime`**"（理由＝常驻腿直达那里），`wisp config set` 是**用户明确发起的另一条 CLI 入口**、不经过 `assembleRuntime`，所以它触发首建**不违反**原文；违反的只是断言自己写死的"计数＝1"。

依据三层：

1. **原文层**（§0.2）：禁止项具名是 `assembleRuntime`，理由具名是"常驻腿直达装配根"。`runConfigSet65` 两条都不沾。
2. **同源层**：票 65 的提交说明（§0.3）把这一跳写成设计声明——**不是事故**。
3. **产品语义层**（§2）：去掉这一跳，干净机器上的 `wisp config set` **必然坏**，而它是今天唯一能把值写进配置的命令行入口。

### 为什么不是甲

甲（去掉 `runConfigSet65` 那一跳）会让 `wisp config set` 在干净机器上**退码 1**，错误句是 `config.toml 读不进来`。而票 65 自己的钉子 `config_cmd_65_test.go:84-86` 逐字要求：

> `t.Fatalf("set on missing file must first-run-create config.toml at %s: %v", cfgPath, err)`

⇒ 甲会**同时红掉票 65 的 AC**。甲不是"改一处"，是"拆一条别人票已验收的判据"。本仓禁止清单（`AGENTS.md §1.1`）明写不许为了自己落地去放宽别人票的 AC —— 甲正是那个形状的**反向**（为了保住自己票的计数，去拆别人票的判据）。

### 为什么不是"放宽成不含常驻腿就算过"

派单明令禁止，且技术上也不该：`!contains(callers, "assembleRuntime")` 这种谓词会让**第三个**调用者（比如将来某条后台腿）静默通过。本修法用 `slices.Equal(callers, []string{"runConfigSet65","runTextTask"})`——**精确集合**：多一个红、少一个也红。

### 仪器是否还有牙（本修法不能是"把红改绿"）

- **多一个**：`tree-assemble` 变异（把那一跳种进 `assembleRuntime`）⇒ 红，逐字见 §3。
- **少一个**：`tree-no-set-call` 变异（去掉 `runConfigSet65` 那一跳）⇒ 红，逐字见 §3。
- **行为半**另有一枚独立的、非 AST 的钉：`TestTicket198R2J1TheAssemblyRootStillCreatesNothing`（`firstrun_198r2_test.go:432-472`）——两个同样的空目录，装配根**必须不建**、`wisp run` 入口**必须建**。所以"害"那一半今天**有行为仪器**，不只靠这枚计数。

---

## 2. 产品语义：`wisp config set` 在无配置文件时该失败还是该创建？

**判：该创建。** 依据四条，逐条可核：

1. **它不创建就完全不能用**。`config.NewManager` 第一句 `LoadFile`，缺文件即错（`internal/config/manager.go:111-123`）。`config set` 的全部语义就是"改一个键"，没有文件就没有可改的东西 ⇒ 不创建＝这条命令在干净机器上**恒不可用**。
2. **"创建"这件事的能力与形状早已被票 198 定死并验收**：`config.SaveFile(cfgPath, config.NewDefaults())`，默认值单一来源是 schema 的 `default` 标签（`internal/config/defaults.go:58`），落点/权限走 `secret.NewStore`（`internal/secret/store.go:49`）。`ensureFirstRunConfig` 就是那个封装。复用它是**唯一**不造第二真相源的写法（票 198 明确禁止新写模板字符串）。票 65 的提交说明也正是这么写的（"缺文件set走ensureFirstRunConfig唯一首建口"）。
3. **主流 harness 的形状同向**：`git config user.name` 在没有 `~/.gitconfig` 时**创建它**；`kubectl config set-*` 在没有 kubeconfig 时**创建它**。这正是票 65 提交说明自陈选形的依据（`config_cmd_65.go:12-18` 逐字 "THE SHAPE IS THE MAINSTREAM ONE, ON PURPOSE (铁律 ⓪⓪⓪): `git config user.name` / `kubectl config set`"）。"写命令遇到不存在的文件就失败"是**反主流**的。
4. **`get` 与 `set` 的分界本来就是产品设计好的**，不是随手：`get` 缺文件时**答 schema 默认值、不建文件、rc 0**（`config_cmd_65.go:131-156`，钉子 `config_cmd_65_test.go:58-69`）；`set` 缺文件时**建**（`config_cmd_65.go:158-176`，钉子 `config_cmd_65_test.go:71-86`）。**读不创建、写创建**——这条分界在票 65 里是明写的（`config_cmd_65.go:161-163` 注释逐字 "a `get` never creates"）。让 `set` 失败会把这条已定的分界抹平。

**⚠ 与 J1 不冲突，因为两处的"该不该建"判据不同**：J1 判的是"**哪条路**可以触发首建"（答：用户明确发起的），不是"**首建能不能发生在 `run` 之外**"（若那样读，`config set` 与 `secret set` 都得坏）。`wisp config set` 是用户明确发起的 —— 与 `wisp run` 同类。

**⚠ 一句必须说清的边界**：本判定**不改**"常驻腿/双击那条链不建件"这一格。台账 `:10953`（`A531` 定性）逐字写着"`ensureFirstRunConfig` 唯一产码调用点仍是 `run.go:245`（挂 `wisp run`），双击那条链仍无人建件（这条维持 A531 定性：真缺口，待 256 那族之后立票）"——**那句话在本票之后过期了**（现在有两个产码调用点），但**它描述的真缺口没变**：双击那条链（`assembleRuntime`）仍然不建件，本票没动它，`TestTicket198R2J1TheAssemblyRootStillCreatesNothing` 与票 246 那枚钉仍然守着。

---

## 3. 突变控制（仓外磁盘改写；⛔ 非 overlay）

### 3.1 为什么必须仓外改写

这枚检查在**运行时读磁盘源码**（`loadMainPackage131(".")` → `os.ReadDir` + `parser.ParseFile`，`leg_sink_gate_131_test.go:672-697`）。`go test -overlay` 只换**编译期**输入，改不动运行时 `os.ReadDir` 看到的东西（台账 `A1804` 那条分界线）。⇒ 必须改写**测试二进制当前目录下**的源码树。

### 3.2 台件与做法

- 台件根：`D:\tmp\wisp\20261112\434-firstrun198\`（**仓外**，符合"临时件只建不删"）。
- 三棵树，各为 `cmd/wisp/*.go` 的**完整拷贝**（269 枚 `.go`），从**修后**工作树复制：
  - `tree-pristine/` —— 原样（正控）。
  - `tree-no-set-call/` —— 用正则删掉 `config_cmd_65.go:173-176` 那四行（＝甲的修法形状）。
  - `tree-assemble/` —— 在 `run.go:626` 之后种进 `ensureFirstRunConfig` 那一跳（＝J1 明禁的害）。
- **同一个二进制** `D:\work\workspace\projects plans\Wisp\.scratch\wisp\stage-cli\wisp.test.exe`，只改 `cd` 到哪棵树。

### 3.3 读数（同一二进制，三棵树，三种判定）

```
binary sha256: 89A527C391E3DC160FE84BFDDAF12D20F0741EA6240EC28C31628F743E7D183F   (三发前后逐字相同)

=== tree-pristine : rc=0 ===
=== RUN   TestTicket198FirstRunCallerIsTheRunEntryOnly
--- PASS: TestTicket198FirstRunCallerIsTheRunEntryOnly (0.31s)
PASS

=== tree-no-set-call : rc=1 ===
=== RUN   TestTicket198FirstRunCallerIsTheRunEntryOnly
    firstrun_198_test.go:305: AC J1 RED: ensureFirstRunConfig is called by [runTextTask], want exactly [runConfigSet65 runTextTask] - first-run creation may only hang on entries the user starts, and the resident leg reaches assembleRuntime directly, so creation shared there writes a config.toml its pins forbid
--- FAIL: TestTicket198FirstRunCallerIsTheRunEntryOnly (0.26s)
FAIL

=== tree-assemble : rc=1 ===
=== RUN   TestTicket198FirstRunCallerIsTheRunEntryOnly
    firstrun_198_test.go:305: AC J1 RED: ensureFirstRunConfig is called by [assembleRuntime runConfigSet65 runTextTask], want exactly [runConfigSet65 runTextTask] - first-run creation may only hang on entries the user starts, and the resident leg reaches assembleRuntime directly, so creation shared there writes a config.toml its pins forbid
--- FAIL: TestTicket198FirstRunCallerIsTheRunEntryOnly (0.36s)
FAIL
```

原始日志：`D:\tmp\wisp\20261112\434-firstrun198\logs\mut2-tree-{pristine,no-set-call,assemble}.raw.txt`。

### 3.4 我怎么确认仪器真在读我改的那棵树

**三把尺，缺一不算**：

1. **同一二进制给出两种不同判定**（上表）——这是决定性的一把。二进制 `sha256` 三发前后逐字相同（`89A527C3…D183F`），**唯一变量是 `cd` 到哪棵树**。若仪器读的是编译期烘焙进二进制的内容，三发必然同色。
2. **红句里的集合逐字跟着树变**：`[runTextTask]` → `[runConfigSet65 runTextTask]` → `[assembleRuntime runConfigSet65 runTextTask]`。红句点名的是**树里的函数名**，不是二进制里的常量。
3. **红句行号漂了**：改前 `firstrun_198_test.go:281`，改后 `:305`——因为我在同一文件里加了注释行。这证明它读的是**盘上当前那份文件**，不是缓存。

**⚠ 一条我先踩的坑（记我）**：第一次用 PowerShell 直接 `& $exe -test.run …` 跑，三发全 `rc=2` 且输出是 `-help` 文本——PowerShell 把 `-test.run` 当成自己的参数吞了。改走 `D:\work\soft\Git\bin\bash.exe -c "cd … && exe -test.run …"` 才拿到真读数。**凡带 `-test.*` 旗标的直跑，一律经 bash。**

---

## 4. AC#2 越界与 AC#4 与票 174 的关系

### 4.1 `git diff --name-only`（逐行）

```
.scratch/wisp/issues/426-…md
.scratch/wisp/issues/427-…md
.scratch/wisp/issues/428-…md
.scratch/wisp/issues/429-…md
.scratch/wisp/issues/430-…md
.scratch/wisp/issues/431-…md
.scratch/wisp/issues/432-…md
.scratch/wisp/issues/433-…md
.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go
.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt
cmd/wisp/firstrun_198_test.go          ← 本腿唯一改动
cmd/wisp/resident_audio_247_windows_test.go
frontend/fixtures/composer-states.html
```

**⚠ 这个输出是共享工作树的全量 diff，不是本腿的写面**。其余 12 枚是**别的在飞腿**的改动（票 426–433 的票面、291 探针、247 测试件、前端 fixture、346 派单件）。本腿的写面＝**`cmd/wisp/firstrun_198_test.go` 一枚**（外加本证据件与票 434 的 Progress log 一行）。本腿 `commit` 用**显式 pathspec** 只提交这三枚。

**禁止面零命中**：`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` —— 本腿**一枚都没碰**（`frontend/fixtures/composer-states.html` 那行是别人的）。

### 4.2 AC#4 与票 174 的关系

票 174 动过 `run.go` 的 `[fs]` 授权面，已回退（`7dcb145ff`）。本腿：

- **`cmd/wisp/run.go` 零改动**（`git diff --name-only -- cmd/wisp/run.go` ⇒ 空）。
- 改动全在 `cmd/wisp/firstrun_198_test.go`（测试件），与 `[fs]`／`risk.PathResolver`／授权根**无交集**。
- `7dcb145ff` 是 `run.go` 最近十枚提交里的最新一枚，本腿在 `git log` 里读到它但**未基于它做任何事**。

⇒ **无关**，且没顺手碰。

---

## 5. AC#3 门禁四读数

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/wisp/` | **rc=0**（零输出） |
| vet | `go vet ./cmd/wisp/` | **rc=0**（零输出） |
| gofumpt | `gofumpt -l cmd/wisp/firstrun_198_test.go` | **list-count=0**（零输出） |
| d22scan | `sh scripts/d22scan.sh` | **clean**：正控 `PASS=45 FAIL=0 SKIP=0`；`d22scan: clean - no D22 ban violations`；分母 `bans #1-5 internal/=403, cmd/=77`、`ban #8 internal/=1004, cmd/=276` |

⚠ `build/mut_rvw.go`（未跟踪残留件）**未删、未碰**；门禁用 `go build ./cmd/wisp/`（不是 `./...`）正是为避开它。

### 5.1 附带复跑（本腿现跑）

| 命令 | 读数 |
|---|---|
| `-test.run '^TestTicket198FirstRunCallerIsTheRunEntryOnly$'` | **PASS 0.41s**（改前 FAIL 0.44s） |
| `-test.run '^TestTicket198'`（整族 10 枚） | **10/10 PASS**（含 `TestTicket198R2J1TheAssemblyRootStillCreatesNothing` 0.03s） |
| `-test.run '^TestConfig'`（票 65 的 CLI 族 6 枚） | **6/6 PASS**（含 `TestConfigLegGetSetGuarded65`，它逐字要求"set on missing file must first-run-create"） |

⇒ **乙的修法没有拆掉票 65 的判据**（甲会）。

---

## 6. 判语与残余

**判语：乙成立，且断言仍有牙（多一个红、少一个也红）。**

- 规矩原文（"⛔ 不许进 `assembleRuntime`"）**未被放宽**：`assembleRuntime` 仍在禁止侧，`tree-assemble` 变异实测红。
- "恰好"这个约束**未被放宽成谓词**：精确集合比较，第三个调用者会红。
- 行为半另有独立仪器（`TestTicket198R2J1TheAssemblyRootStillCreatesNothing`），不只靠 AST 计数。
- 产码零改动 ⇒ 票 65 的已验收行为（干净机器 `set` 建件）**原样保留**。

**残余（明写不糊）**：

1. **测试名 `…IsTheRunEntryOnly` 现在名不副实**。⛔ 不改名：至少 12 份证据件与台账条目按这个名字引用（`grep` 现量 21 处），改名会让所有历史读数失联。已在函数头注释里写清"读作：调用集是用户发起的入口，不是装配根"。
2. **台账 `:10953` 那句"`ensureFirstRunConfig` 唯一产码调用点仍是 `run.go:245`"自本票起过期**（今有两个）。那句话描述的真缺口（双击链不建件）**没变**，本票没动它，归票 256 那族。
3. **本判定不回答"将来再加第三条用户发起入口该怎么办"**。本修法的形状是"每加一条必须显式改这枚断言"——那是**有意的**摩擦：加调用者应当是一次需要被看见的动作，不是一次静默通过。

---

## 7. 撤销口令

**「434 回甲」** —— 把 `cmd/wisp/firstrun_198_test.go` 的 `want` 改回 `[]string{"runTextTask"}`，并删掉 `config_cmd_65.go:173-176` 那四行。⚠ 执行前先读本件 §1「为什么不是甲」与 §5.1（甲会红票 65 的 `TestConfigLegGetSetGuarded65`）。
