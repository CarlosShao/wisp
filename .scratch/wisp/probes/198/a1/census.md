# 票 198 · 腿 198-a1 · 只读普查（首建 `config.toml`）

> 本文件是**料**，不是产码。§1–§5 给落地腿当派单名册用。所有行号按 §0 锚点自取实测。

## §0 起手锚

| 项 | 读数 | 尺 |
|---|---|---|
| HEAD | `6fe300c5035a60586e681f586a4ead5e2b56137b`（`census(rate-census-3a)：编号 150-199 的 38 枚开放票分桶…`，提交时刻 `2026-10-02 08:41:33 +0800`） | `git log -1 --format="%H%n%ad%n%s" --date=iso` |
| 分支 | `dev` | `git rev-parse --abbrev-ref HEAD` |
| 本腿起手时刻 | `2026-10-02 08:45:58 +0800` | `date "+%Y-%m-%d %H:%M:%S %z"` |
| 票面 | `.scratch/wisp/issues/198-a-fresh-machine-cannot-run-wisp-because-nothing-creates-config-toml.md` | **实测 35 行**，非编排者题面写的 48 行（见 §7-X1） |

`git status --porcelain cmd internal` 起手读数（**非空属正常：那是别枚在飞的腿，本腿一字未动未提交**）：

```
 M internal/panel/bridge.go
 M internal/panel/composer_dispatch.go
 M internal/panel/composer_dispatch_test.go
 M internal/panel/l2_grant_boundary_test.go
?? internal/ball/sta_release_windows_test.go
?? internal/config/settings.go
?? internal/panel/config_handlers.go
```

⚠ 与本票直接相关的一枚：`internal/config/settings.go` **此刻是未跟踪文件（313 行）**，内容尚未进任何 commit。
本腿读过它（见 §1/§4），但它**随时可能被那枚腿改写或删除**，所以凡本文件引它处均标 `[untracked-197]`，
落地腿开工时必须重新取一次行号，不得照抄本文件的数字。

## §1 「首建」这一跳今天到底缺哪几环

### §1.0 一句话结论
**能力齐了，缺的只是入口那一次调用。** 序列化器（`MarshalCanonical`）、落盘器（`SaveFile`）、
原子写＋封口（`atomicWrite`＋`SealFile`）、默认值表（`NewDefaults`）、目录创建＋逐层封权（`PrivateDirAll`）
**五件全部现成且已在生产里活着**；`wisp run` 唯一缺的是在 `config.NewManager` 之前把
`SaveFile(cfgPath, config.NewDefaults())` 调一次，以及调完之后那两句实话。
⇒ 落地腿**不要新写任何序列化/写盘代码**（论证 §1.5）。

### §1.1 从入口到退码 2，逐环（HEAD 实测行号）

| 环 | `file:line` | 读数 | 口径 |
|---|---|---|---|
| R1 | `cmd/wisp/main.go:89-91` | `case "run": attachParentConsole(); os.Exit(cmdRun(args[1:]))` | 分派 |
| R2 | `cmd/wisp/main.go:151-165` | `cmdRun` 打版本、取 `interactiveStdin()`，然后 `return runTextTask(runSpec{argv, stdout, stderr, reply})` | ⚠ `dataDir` **留空**，交给 R4 去问 OS |
| R3 | `cmd/wisp/run.go:194-198` | 任务文本为空 ⇒ 打印用法 + `return 2`（与首建无关的另一个 2，别混） | 退码 2 的第一处 |
| R4 | `cmd/wisp/run.go:199-214` | `s.dataDir == ""` ⇒ `resolveDataDir(buildEnvString())`（`:200`），失败 ⇒ `:210` 打印 + `:211 return 2` | 票 128 AC#2 的钉 |
| R5 | `cmd/wisp/doctor.go:247-268` | `resolveDataDir`：便携（`portable.txt` 在同 exe 目录）→ `exeDir\data[-dev]`（`:248-255`）；`env=="test"` → `proc.TestDataDir()`（`:256-258`）；否则 `userConfigDir()`（`:259`）→ `proc.SealableRoot(base)`（`:263`）→ `base\wisp`/`base\wisp-dev`（`:264-267`） | **这就是本仓的"私有目录判定者"上游**（§2） |
| R6 | `cmd/wisp/run.go:222-232` | `installLogSink(s.dataDir)` **在装配之前**；装不上只响亮报一句、不停机 | 注释 `:216-221` 逐字说明为什么必须在前 |
| R7 | `cmd/wisp/run.go:234` | `rt, code := assembleRuntime(s)` | 全仓唯一装配根（`resident_task_source_windows.go:265` 是第二个生产调用者） |
| R8 | `cmd/wisp/run.go:371` | `cfgPath := filepath.Join(s.dataDir, configFileName)` | 文件名常量 `configFileName = "config.toml"` 在 `cmd/wisp/secret.go:63` |
| R9 | `cmd/wisp/run.go:372-376` → `internal/secret/store.go:41-53` | `secret.NewStore` ⇒ `winsec.PrivateDirAll(<dataDir>\secrets, 0o700)`（store.go:49）⇒ `privateDirAll` 逐层 `os.Mkdir` + `sealDir`（`internal/winsec/winsec.go:203-217`） | ★ **读配置之前，数据根本体已经被创建并逐层封权了**（`<dataDir>` 是 `<dataDir>\secrets` 的缺失祖先）。首建 config.toml 因此**不需要任何新目录代码** |
| R10 | `cmd/wisp/run.go:389-390` → `internal/config/manager.go:101-105` → `internal/config/loader.go:41-45` → `loader.go:64-68` | `os.ReadFile(path)` 失败 ⇒ `observe.Wrap(observe.ClassConfig, err, "config.toml read")`。**这里是"缺文件 ⇒ 错"唯一被造出来的点** | 题面写的 `loader.go:39-43` 是注释＋签名，不是产生点（§7-X5） |
| R11 | `cmd/wisp/run.go:396-398` | `rt.auditModeUnreadable(cfgPath, err)` → 打印 `wisp run: 配置未就绪（Unconfigured）：%v` → **`return rt, 2`** | ★ 首建唯一该插的位置：`R10` 之前，即 `run.go:388` 与 `:389` 之间 |
| R12 | `cmd/wisp/run.go:68-91` | 退码表：`observe.ClassConfig`/`ClassAuth` ⇒ 2 | 2 是分类映射出来的，不是硬编码在 R11 |
| R13 | `cmd/wisp/run.go:409-413` → `internal/llm/resolver.go:218` | **即使配置建好了也仍然走到这儿再 `return rt, 2`**：`llm: role %q is unset and text_chain is empty (Unconfigured)` | ⇒ 首建**不改变退码**；AC#4 才是用户可见终点（J4） |

### §1.2 `internal/config` 今天的写盘能力（逐枚具名）

| 符号 | `file:line` | 导出？ | 生产调用者枚数 | 备注 |
|---|---|---|---|---|
| `SaveFile(path, *Config)` | `internal/config/loader.go:238`（文档 `:235-237`） | **是** | **1**：`internal/config/writeguard.go:160` | 体＝`deepCopyConfig` → 强制 `SchemaVersion`（`:243`）→ `MarshalCanonical`（`:244`）→ `atomicWrite`（`:248`）；**从不回读文件**（`writeguard.go:8-9` 逐字点过这一性） |
| `MarshalCanonical(*Config) ([]byte, error)` | `internal/config/parse.go:164` | **是** | 仅经 `SaveFile` | 无 `omitempty`，注释 `:160-163` 逐字："every static field emitted explicitly" |
| `atomicWrite(path, data)` | `internal/config/parse.go:196` | 否 | 仅 `SaveFile` + `migrate.go:96` | temp + rename；**写第一个字节之前**就 `winsec.SealFile(tmpName)`（`:215`），理由在 `:207-214` |
| `applyMigrations` | `internal/config/migrate.go:41-98` | 否（包内） | 装载链上（`loader.go:112`） | 唯一"读路径上会写盘"的一族：先 `config.toml.bak-<ver>` 再 `atomicWrite` |
| `Manager.mergeWrite` | `internal/config/writeguard.go:111` | 否 | 3 枚包内入口 | 见下 |
| `Manager.SetPermissionMode` | `internal/config/permmode.go:70` | 是 | `perm.Store` 那一族（`internal/perm/store.go:49` 声明的 seam） | → `mergeWrite("risk.permission_mode", …)` |
| `Manager.AddAllowedDir` / `SetAllowedDirs` / `writeAllowedDirs` | `internal/config/allowdirs.go:62 / :109 / :139` | 是 | `cmd/wisp/approval_always.go` 那一族 | → `mergeWrite("fs.allowed_dirs", …)` |
| `fileMissing(err)` | `internal/config/parse.go:232-235` | 否 | `writeguard.go:125` | 现成的"是不是不存在"判定，落地腿直接复用即可 |

★ **`writeguard.go:125-131` 今天就是一个"缺文件则建"的生产分支**（注释逐字：`writing it creates what first-run did not`）。
它救不了全新机器，纯粹因为 `NewManager`（`manager.go:101-105`）先失败、生产走不到 `mergeWrite`。
⇒ 票面"全仓无一处创建 `config.toml`"这句要收窄（§7-X8）。

**没有任何导出的 "CreateIfMissing / Ensure / Bootstrap"**（尺：`grep -rn "func .*Save\|func .*Write\|func .*Create\|func .*Ensure\|func .*Init\|func .*Bootstrap\|func .*Scaffold" internal/config/ --include=*.go | grep -v _test` ⇒ 命中只有 `allowdirs.go:157 statOwnWrite`、`loader.go:238 SaveFile`、`parse.go:196 atomicWrite`、`writeguard.go:111 mergeWrite`）。

### §1.3 枚数更正：section 是 **18** 枚，不是 16

- `type Config struct` 区间＝`internal/config/schema.go:104-134`（题面行号**对**）。
- 区间内的 section 字段实测 **18 枚**（`schema.go:110-133`）：
  `App Ball Hotkey Session Voice Audio LLM Agent Risk FS Net Privacy Memory Panel Cost Plugins Models Observe`
  （尺：`awk 'NR>=110 && NR<=133' internal/config/schema.go | grep -oE '[A-Za-z]+Section' | sort` ⇒ 18 行）
- 交叉核对：`grep -nE '^type [A-Za-z]+Section struct' internal/config/schema.go | wc -l` ⇒ **18**（行号 137/163/179/187/244/269/406/429/448/474/494/507/521/529/543/572/579/600）。
- 两把尺同向 ⇒ "16" 是**过期枚数**（`Panel`/`Cost`/`Models`/`Observe` 之后加进来时票没跟着改）。

### §1.4 「16 枚 section 的默认值」今天存在的 **4 种形状**

| 形状 | 出处 | 枚数 | 会不会被 `SaveFile` 写进文件 |
|---|---|---|---|
| (a) **`default:"…"` 结构体标签**（D36 rule 3 的唯一真相源） | `internal/config/schema.go` 全体；由 `applyDefaults`（`defaults.go:65-88`）＋`setDefault`（`:92-133`）用反射落地，入口 `NewDefaults()`（`defaults.go:58-62`） | **70**（`grep -c 'default:"' internal/config/schema.go`） | **会** |
| (b) **Go 零值**（没打 `default:` 标签的字段） | 同上，`applyDefaults` 的 `default:` 分支不触发（`defaults.go:80`） | **103**（口径：`schema.go` 里带 `toml:"…"` 标签的结构体字段 173 枚 － 70 枚带默认标签；**含表头型 struct 字段**，所以是上界，权威枚数请在测试里反射遍历取，别抄我这个数） | 会（以零值形式） |
| (c) **内置 provider 预设表** | `providerPresets`（`defaults.go:22-35`），**只在装载时**由 `applyPresets`（`loader.go:188-204`）覆盖空字段 | **12** 条 | **不会**（`NewDefaults()` 把 map 留 nil，见 `defaults.go:77-78`） |
| (d) **结构体外被强制的值** | `SchemaVersion` 由 `SaveFile` 强制为 `SchemaVersionCurrent`（`loader.go:243`；常量 `schema.go:19`）；`App.Portable` 带 `toml:"-"`、由 env layout 填（`schema.go:145-149`） | 2 处 | `SchemaVersion` 会，`Portable` 永不 |

形状 (a) 的一条硬约束要转给落地腿：**slice 默认值按逗号切分并 TrimSpace**（`defaults.go:114-128`），
所以元素里带逗号的默认值**表达不出来**；现用的 widest 一例是 `default:"取消,停下,别"`（`schema.go` 的 hotkey 族）。

### §1.5 「默认值 → TOML 文本」的通路**已经存在**，不要新写一份

- 现成通路＝`SaveFile(path, NewDefaults())`，两行、全导出。
- 它**今天已经被一枚现成用例钉成"可加载"**：`internal/config/unwired_test.go:86-95`
  （子用例名 `"SaveFile output reloads"`：`SaveFile(path, NewDefaults())` 后 `LoadFile(path, nil)` 若报错即
  `t.Fatalf("a config the app itself wrote must load: %v", err)`）。⇒ 落地腿踩的是已铺好的地，不是新地。
- ⚠ **维护面对比（题面要我具名指出的那一枚）**：
  - 走 §1.5 的通路 ⇒ **零新增维护面**（默认值只存在于 `default:` 标签这一处，写盘侧是反射读的）。
  - 若落地腿改为**手写一份 TOML 模板字符串/常量** ⇒ 引入 `默认值双写` 面：
    70 枚标签值＋18 个 section 名要在字面量里再抄一遍，而今天没有任何用例比对"模板 vs 标签"
    （`boundary_test.go:186-212` 那两枚 round-trip 钉的是 `fullConfig(t)` 手造结构体，**不经过模板**；
    `unwired_test.go:86-95` 钉的是 `NewDefaults()`，也**不经过模板**）
    ⇒ 结果就是"标签改了、模板没改，全绿"，正是本仓 A382/A391 那一族"字段恒空要能被判红"的反面。
    **结论：写模板＝引入一枚无人看守的第二真相源，本票明确不该做。**

---

## §2 落点与权限（不许自己拼路径）

### §2.1 私有目录判定者：导出名与签名

| 角色 | 导出名 + 签名 | `file:line` |
|---|---|---|
| **目录判定者（要建目录就调它）** | `func PrivateDirAll(path string, perm fs.FileMode) error` | `internal/winsec/winsec.go:165`（文档 `:144-164`） |
| 内部实现（先 C26 再建） | `func privateDirAll(dir ResolvedPath, perm fs.FileMode) error` | `internal/winsec/winsec.go:177-218` |
| 单目录收窄 | `func SealDir(path string) error` | `internal/winsec/winsec.go:222` |
| 单文件收窄 | `func SealFile(path string) error` | `internal/winsec/winsec.go:136` |
| **拼不出合法路径就出不来的类型** | `type ResolvedPath struct` ＋ `func ResolvePath(input string) (ResolvedPath, error)` | `internal/winsec/resolve.go:583-585` / `:598-636`；"唯一铸造者是 ResolvePath"这句逐字在 `:576-582` |
| OS 答复的去链接前缀 | `func SealableRoot(path string) string` | `internal/proc/envfork.go:160`（生产调用：`cmd/wisp/doctor.go:263`） |
| **数据根判定者（本票落点的上游）** | `func resolveDataDir(env string) (string, error)` | `cmd/wisp/doctor.go:247-268` |

⇒ 落地腿的正确形状：**R5 的返回值 + `configFileName` 拼出 `cfgPath`（R8 已拼好，就在 `run.go:371`），
直接把它交给 `config.SaveFile`；目录侧什么都不用做**——R9 已经证明 `<dataDir>` 在 `NewManager` 之前
就被 `secret.NewStore` → `PrivateDirAll` 建好并逐层封权了（`store.go:49`、`winsec.go:203-217`）。

### §2.2 `risk.PathResolver` 与 C26 的射程（别把两件事当一件）

- `internal/risk/pathresolver.go:13-14` 逐字：**C26 是"全仓唯一的路径归一化入口"**；`:148` 注：C26 的 canonical 输出
  "against other C26 output and never handed to a write"。⇒ **C26 的射程是"决策"（工具目标路径、allowlist 命中判定），不是"落盘"**。
- `internal/winsec` 一侧不是第二个 PathResolver：`resolve.go:658-662` 逐字说明 `builtinVerifier` 是
  "没有 C26 链接时跑的地板，它**不做法任何归一化**——要么原样放过要么拒绝"，这正是"再加一枚它不算违反 SPEC-06 §4"的理由。
- winsec 的双层守卫值得抄进本票的理由清单：`ResolvePath` 既问装着的 resolver，**也**把 resolver 的答案自己再过一遍地板（`:631-634`），
  并对"答案离开了调用者点名的那棵树"直接拒（`:621-624`）。
- ⚠ `AGENTS.md §1.2` 那条禁令的准确射程是**`filepath.Clean|Abs`**（`PLAN.md:1288`），**不是 `filepath.Join`**：
  `cmd/wisp/run.go:371`、`cmd/wisp/doctor.go:249/251/253/265/267` 今天都在用 `filepath.Join` 拼现成根目录。
  ⇒ AC#3 那句"不许 `filepath.Join` 手拼"若按字面判，会把**现存的 R8 那一行**一起判违规。这一处口径需要编排者裁（并入 J5/J6 那条读法问题）。

### §2.3 票 132 的 `SealDir` 一族：生产调用者枚数（按"谁产出／谁投递／谁落盘"三处判）

尺（本腿实跑）：`grep -rn "winsec\.\(PrivateDirAll\|SealDir\|SealFile\)" --include=*.go . | grep -v _test.go | grep -v ^./internal/winsec/ | grep -v ^./.scratch`

| 层 | 读数 |
|---|---|
| **显式调用 `winsec.SealDir` 的生产代码** | **0 枚**。尺：`grep -rn "SealDir(" --include=*.go . \| grep -v "^./.scratch"` ⇒ **17 处命中**，其中 **16 处是 `internal/winsec/*_test.go` 的调用点**、**1 处是定义本身**（`winsec.go:222`）；去掉 `_test.go` 后**只剩定义** |
| 谁产出封权动作 | 包内未导出的 `sealDir`：`internal/winsec/winsec_windows.go:588`（Windows，注释 `:582-587` 说明"继承在创建时决定，所以要连既有子项一起修"）、`internal/winsec/winsec_other.go:170`（POSIX）；被 `winsec.go:213`、`:217` 调用 |
| 谁把它投递给生产 | **`PrivateDirAll`**（`winsec.go:165`）与 **`SealFile`**（`winsec.go:136`）。生产调用者：`PrivateDirAll` ⇒ `internal/secret/store.go:49`、`internal/memory/open.go:176`、`:188`、`:503`、`internal/agent/spill.go:111`（**5 处**）；`SealFile` ⇒ `internal/config/parse.go:215`（**每一次 `SaveFile` 都过它**）、`internal/secret/migrate.go:174`（**2 处**） |
| 谁落盘 | `atomicWrite`（`internal/config/parse.go:196-229`）：temp 建在**目标同目录**（`:197-201`），**写内容之前** `SealFile(tmpName)`（`:215`），再 `os.Rename`（`:226`）——注释 `:207-214` 逐字解释"descriptor 随 rename 带过去，所以 temp 的 ACL 就是之后 config.toml 的 ACL" |

⇒ **判定**：`SealDir` 的"0 枚生产调用者"是**命名造成的错觉，不是能力缺口**。
票 132 那一族今天**确实**在生产里逐层建目录＋封权，只是入口叫 `PrivateDirAll`／`SealFile`。
所以 AC#3 该按"同一族判定者投递"读，不该按"`SealDir` 这个符号名命中"读（→ **J5**，需编排者裁一句）。
现成可抄的权限正控：`internal/config/private_acl_windows_test.go:103-128`
`TestAC2SaveFileLandsPrivate` 已经把"`SaveFile` 落出来的文件是私有的"钉住（其注释 `:54` 点名"…SaveFile red, which is AC#4"）。

---

## §3 说实话的引导该长什么样

### §3.1 本仓"说实话"文案的现成家族（逐枚 `file:line`，落地腿照这一族写，别另起一形）

| 家族 | `file:line` | 它已经做到的事 |
|---|---|---|
| **claims/posture 常量族** | `cmd/wisp/resident_task_source_windows.go:84-110` | 把每种降级各写成**一枚常量**（`taskEntryDisabledClaim :88`、`pipelineAbsentClaim :100`、`taskPostureRefused :110`），并在注释 `:86-87`／`:101-106` 逐字说明"诚实/撒谎的区分活在措辞里"、"三件事过去印成同一句无话，AC#7 就是为拦它" |
| **四类归因族（最贴近本票）** | `cmd/wisp/config_reload.go:315-370`（设计说明 `:305-314`） | `cause=missing :317-320`／`permission :321-324`／`syntax :330-333`／`unknown-key :334-337`／`migration :338-356`／`invalid :357-360`／兜底 `:363-369` 明说"没归入已知四类"。注释 `:305-308` 逐字："**一句共享的『配置未生效』会是省略式的谎话**（票 216 立的规矩）" |
| **"缺什么 + 去哪儿补"完整式样** | `cmd/wisp/doctor.go:289-298` `dataDirUnresolved128` | 一句话里同时给出：OS 原话（`:294`）、本该落在哪个树（`:294-295`）、**为什么不回落到 cwd 的实测后果**（`:295`）、以及 `修法：Windows 把 APPDATA 设为…`（`:296`）。**这就是票面要求的"它在哪、还缺什么"的形** |
| **key 该往哪儿补** | `cmd/wisp/secret.go:160-178`（`wisp secret set … :162-164`、`The reference to paste into config.toml is api_key_ref = "dpapi:<name>" :175`）＋ 实际回执 `cmd/wisp/secret.go:380`（`reference for config.toml: api_key_ref = %q`） | 已经在告诉用户"补 key 的**两条手**"：`wisp secret set` 与"把 ref 粘进 config.toml" |
| **写完了要报到哪、何时生效** | `cmd/wisp/approval_always.go:143` | `长期规则已写入 config.toml 的 [fs] allowed_dirs：%s（下一次启动生效，本次运行仍按旧名单）` —— "写到哪个 section + 什么时候算"同句给出 |
| **没人能答复也要说** | `cmd/wisp/main.go:155-157` | 没有交互控制台时**开机就讲后果**，不等退码 |
| **热加载接管回执** | `cmd/wisp/config_reload.go:123-125`、`:141` | 装配成功／收口退出各有其句 |

### §3.2 首建回执该含的四件事（据 §1.1 与票面 AC#2/AC#4 反推，不含产码）
1. **建了没建成**（成功／响亮失败两形，AC#5）＋**绝对路径**（`cfgPath`，R8 那枚现成值）。
2. **值从哪来**：一句"全部取内置默认"，并点名 §1.4 形状 (a) 是唯一真相源 —— 否则 owner 会以为有人替他选了模型。
3. **还缺什么**：今天首建之后**必然**缺 (i) `[llm].text_chain`／角色指向的 provider（R13），(ii) 那个 provider 的 `api_key_ref` 背后的 blob。
4. **去哪儿补**：`wisp secret set <name>` ＋ `api_key_ref = "dpapi:<name>"`（`secret.go:175` 原文）＋ `wisp providers discover <provider>`（`main.go:35-37`）。
   ⚠ 并且必须同句写"**本轮仍会退码 2**"（J4/R13），否则第 3 件就是空头承诺——这正是"建完就当没事发生"的对称面。

### §3.3 现有文案里**今天已经是假话**的那几句（本腿逐枚复认，落地腿不要复用这些句子当既成事实）

| # | `file:line` | 原文（逐字/近逐字） | 为什么今天已是假话 |
|---|---|---|---|
| L1 | `internal/config/manager.go:99-100` | `NewManager loads the initial config from path (first-run file creation is the caller's concern)` | 它把责任派给"the caller"，而**全仓不存在这样一个 caller**（§1.2 的 caller 名册里生产侧只有 `writeguard.go:160`）。这是票面已点名、本腿复认成立的一枚 |
| L2 | `internal/config/writeguard.go:126-127` | `No file yet: … writing it creates what first-run did not.` | 预设了"有一回 first-run 跑过了但没建"。**今天没有任何 first-run 路径**，所以这句在指一个不存在的动作主体 |
| L3 | `internal/config/settings.go:271` `[untracked-197]` | `file is there at all. A missing file is not an error (first-run writes one);` | **最直白的一句**：明文写"首建会写一份"。它此刻只是**未进任何 commit 的在飞文件**里的注释——落地腿**绝不能把它读成对现有行为的描述**（起手必须先重取此文件行号，见 §0 警告） |
| L4 | `internal/llm/errors.go:190`（＋ `internal/observe/errors.go:125-136`） | `maps to the Unconfigured ball state via observe.MappedStates` | `MappedStates`（`internal/observe/errors.go:131`）**生产侧 0 枚调用者**（尺：`grep -rn "MappedStates" --include=*.go .` 去掉 `.scratch`/`frontend` ⇒ 只有 `observe/errors_test.go:80,109` 与 `llm/matrix_14_2_test.go:171,210`）；而真实球体的状态机开机钉在 `StateSleeping`（`cmd/wisp/resident_ball_windows.go:170`），**永不进 FirstRun/Unconfigured**。⇒ "缺配置会让球显出 Unconfigured" 是一句没有铸造者的话 |
| L5 | `internal/statemachine/table.go:43`、`:47`（事件 `events.go:12`、`:14`） | D43 行 #1 `FirstRun + EvOnboardingCompleted -> Sleeping`；行 #3 `FirstRun + EvKeyMissing -> Unconfigured` | 两枚事件**在生产里无人触发**（尺：`grep -rn "EvOnboardingCompleted\|EvKeyMissing" --include=*.go .` 去 `.scratch`/`frontend` ⇒ 只有 `statemachine` 包自身与其测试）。票面要的"引导"在 D43 上已有名字、没有身体。⚠ **D43 转移表一字不许动**（AGENTS.md §1.1），本腿只登记，不建议改写 |

---

## §4 退码 2 那条路不能顺手放宽：会撞上谁的钉（**名册，不是颜色**）

⛔ 本腿禁跑测试 ⇒ 下表是**逐枚读断言**读出来的名册，颜色请落地腿用给的命令自取。
每行都写了"那句断言的存在理由"，因为验收腿要判"改得对不对"，不是"改绿了没有"。

| # | 用例名 | `file:line` | 断言（逐字摘要） | 存在理由 | 与本票的冲突判定 |
|---|---|---|---|---|---|
| N1 | `TestAC2SealNoticeLandsInTheRunLegLogFile` | `cmd/wisp/logsink_windows_test.go:137-165`（注释 `:159-162`；`if code != 2 { t.Fatalf("exit %d, want 2 (Unconfigured)…") }` 在 `:163-164`） | 造一个**只有 `secrets/` 子目录、没有 config.toml** 的 dataDir，跑 `runTextTask`，**要求退码 2**，再从日志 sink 里读出封权回执 | 注释 `:159-162` 逐字："这里没有 config.toml，所以这轮停在 Unconfigured。**那是本用例的要点而不是它的偶然**：回执必须活过一次连存储都开不上的启动" | ★**正面冲突**。首建若在 `runTextTask`/`assembleRuntime` 之内，这轮就会往前走。退码可能**恰好仍是 2**（撞 R13 的角色未配置），但**注释与用例意图同时变假**，而 `:167-175` 读的仍是同一枚 sink ⇒ **必须真跑才知道今天绿不绿**；命令见 §4.1 |
| N2 | `TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline` | `cmd/wisp/resident_task_source_246_windows_test.go:358-391`（`:389-391`） | `if _, err := os.Stat(filepath.Join(dataDir, "config.toml")); !errors.Is(err, os.ErrNotExist) { t.Errorf("AC#7 RED: the leg created a config.toml it was never asked for (err %v)") }`；另 `:369` 要 `pipelineAbsentClaim`、`:385` 要 `任务来源：taskPostureRefused` | 246 AC#7：启动报告必须能分辨"没入口／管线被拒／注入位"三形，且**没人要它建文件它就不能建**（跑的是 `buildWispForTest` 真 exe，`:359`） | ★★**这是决定 J1 的那枚钉**。`assembleRuntime` 被常驻腿共用（`cmd/wisp/resident_task_source_windows.go:265`）⇒ 首建放 `assembleRuntime` 里 = 本用例至少 3 处转红。放 `runTextTask` 专属 ⇒ 本用例天然不受影响 |
| N3 | `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128` | `cmd/wisp/dataroot_128_test.go:281-346`（空目录断言 `:343`→`assertDirEmpty128 :382-396`；调用相等性 `:314-321`；腿表 `refusalLegs128 :113-146`） | 每条腿在无数据根时：退码非 0（`:329-331`）、必须含救援标记（`:332-340`）、**启动目录 `WalkDir` 后必须一个条目都不许长出来**（`:343`）；并且 `:314-317` 要求"**每一个调用 `resolveDataDir` 的生产函数**都必须在这里登记成一条腿"（`want` 与 `got` 取**相等**，非包含，`:284-292` 逐字） | 票 128 AC#1 实测过"回落到 cwd ⇒ 日志/config.toml/secrets 搬家"；`:348-351` 逐字说"腿表就是调用图，只有加消费者时它才会长" | ⚠ **两处会咬人**：(i) 首建**必须严格在 `resolveDataDir` 成功之后**（R11 的位置满足；写在 R4 之前即红）；(ii) 若落地腿新写一枚"自己解析数据根"的函数，`missing` 非空即 `t.Fatalf`（`:314-317`）⇒ 必须同时往 `refusalLegs128()` 加一条 |
| N4 | `TestManagerMissingFileKeepsCurrentAndErrors` | `internal/config/manager_test.go:279-296`（`:290-292`、`:293-295`） | `NewManager` 成功后**删掉文件** ⇒ `CheckAndReload()` **必须返回错误**（"missing file must surface an error (config deleted mid-run is a config-class event)"），且内存里的 `Ball.Size` 必须还是 56 | SPEC-03 §4.3：运行中被删是**配置级事件**，不是"再给我变一份出来"的理由 | **只钉住"别把补建塞进 watch/CheckAndReload 那一族"**。首建放装载前（R11）不受影响；若图省事在 `CheckAndReload` 里补建 ⇒ 本用例立刻红 |
| N5 | 三形变异（`存储损坏` / `版本不认识` / `权限读不到`） | `cmd/wisp/run_mode101_test.go:605-614 / :615-620 / :621-633`；断言在 `:647-656`＋`:660` | `code == 0` ⇒ "AC#3 wants a loud failure"；`code != 2` ⇒ "want 2 (Unconfigured: SPEC-03 §4.1 forbids a half-configured run)"；`reached == true` ⇒ "a decision chain was assembled over an unreadable mode store"；还要日志含 `MODE-READ-FAILED`（`:660`） | 票 101 AC#3：档读不到必须**分类失败＋fail-closed**，不许"沿用上次的" | ⚠ **`权限读不到`（`:621-633`）是个陷阱形**：它把 config.toml 删掉后**在原位放一枚同名目录**。若首建的判据写成"stat 失败就建" ⇒ 会去朝一个目录 `os.Rename`（`parse.go:226`），行为未知；判据写成 `fileMissing`/`fs.ErrNotExist`（`parse.go:232-235`）才天然避开。**颜色必须真跑** |
| N6 | `TestMissingBlobFailsUnconfiguredNeverSilently` | `cmd/wisp/run_test.go:446-473`（`:457-459`、`:460-462`、`:463-465`、`:466-468`、`:469-472`） | 删掉凭据 blob ⇒ 退码**必须 2**、stderr 必须含 `Unconfigured`、**0 个请求到达 provider**、输出不得含 key 值 | A8/D33：缺凭据永不静默降级 | 不冲突（fixture 自己已经写了 config.toml，首建分支不会触发）。但它是 **AC#4 的邻格**：首建之后仍然缺 key ⇒ 这一族语义必须原样成立，"退码 2 配一句看不懂的错"里那句**看不懂的错**不能被换成一句安慰话 |
| N7 | `TestUnwiredSecurityKeysFailLoudly` / `TestUnwiredGuardLeavesHonestConfigsAlone` | `internal/config/unwired_test.go:33-54 / :55-105`；守卫本体 `internal/config/unwired.go:39-41`、`:99-111` | 只对**非默认值**开火（`unwired.go:39-41` 逐字："a file that never mentions these keys — and one that SaveFile writes back out while they sit at their defaults — keeps loading exactly as before"） | 票 83／裁定 A53②：写了不生效的安全键＝骗人的配置项 | **不拦本票，反而是本票的现成正控**（`unwired_test.go:86-95` 已钉 `SaveFile(NewDefaults())` 可加载）。⚠ 但**一旦首建往文件里写进任何非默认的安全键，立刻红** ⇒ 这条与票面"不许新造默认值"是同向的，可当 AC#2 的守门人 |
| N8 | `TestEveryLockedSectionKeyIsAccountedFor` | `internal/config/unwired_test.go:311`（表在 `unwired.go:113-147`） | 4 枚锁定 section 的**每一枚键**都要在 `lockedKeyDisposition` 里有说法 | `unwired.go:116-118` 逐字："这张表就是不让它像那 4 枚 [risk] 键那样过期" | 与本票**无冲突**（不改 schema 就不触发）。登记它的原因：若落地腿为了"默认值出处"去动 `schema.go` 的字段形状，这是第一枚会响的钉 |
| N9 | 四句归因族 | `cmd/wisp/config_sentences_223r2_test.go:92`（它自己就用 `config.SaveFile(path, config.NewDefaults())` 起局）＋ `cmd/wisp/config_reload_223_test.go:510`（`TestTicket223FailureSentencesAreDistinct`） | 223 r2 钉"缺失／语法／权限／schema 四句必须各是各的"，`config_reload.go:308-310` 逐字："每个分支都由一枚植入件练到，**而且每个用例还断言其他几句不出现**" | 票 223 AC#4 | ⚠ **软冲突**：首建回执若**复用** `cause=missing` 那一句（`config_reload.go:317-320`），会被这类"其他句不得出现"的断言判红；而且那句的内容"本次运行继续用内存里的旧配置"在首建场景里根本没有旧配置 ⇒ **必须另起一句**（§3.1 家族里选形，别抄这句） |
| N10 | `TestAC2SaveFileLandsPrivate` | `internal/config/private_acl_windows_test.go:103-128`（同族注释 `:54`） | `SaveFile` 落出来的文件 ACL 必须是私有的 | 票 89/131/132 的 ACL 族；`atomicWrite` 先封后写（`parse.go:207-218`） | 不冲突，且是 **AC#3 的现成正控**：只要落地腿**走 `SaveFile`** 就自动满足；**若绕过它用 `os.WriteFile` 直写就没人钉了** |
| N11 | `[untracked-197]` 的 `settings.go` 一族 | `internal/config/settings.go:5-21`（开头逐字："WHY THESE METHODS EXIST AND SaveFile DOES NOT. Ticket 226 measured that SaveFile is a serializer, not a merge"）、`:215 candidate = NewDefaults()`、`:231/:280/:285` 三句 refusal | 那枚在飞的腿正在给 `Manager` 加一层受保护的写路径，并明文把 `SaveFile` 贬为"不是 merge" | 票 226/197 的收口 | ⚠ **写面碰撞预告**（本腿不动它）：落地腿开工时**必须以 HEAD 重取 `internal/config/settings.go` 的行号**；§3.3-L3 那句假话就在枚里，若它进了主线，"first-run writes one"就会变成**有人已经建了**的宣称 ⇒ 两枚腿的口径要先对一次 |

### §4.1 取颜色的确切命令（**本腿不跑**，交落地腿自量）
```
go test ./cmd/wisp -run 'TestAC2SealNoticeLandsInTheRunLegLogFile|TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128|TestMissingBlobFailsUnconfiguredNeverSilently' -v
go test ./cmd/wisp -run 'TestTicket223FailureSentencesAreDistinct|TestTicket223R2FailureSentenceRouting' -v
go test ./internal/config -run 'TestManagerMissingFileKeepsCurrentAndErrors|TestUnwiredGuardLeavesHonestConfigsAlone|TestEveryLockedSectionKeyIsAccountedFor|TestSaveFileRoundTripsThroughLoad' -v
```
（N2 那枚走 `buildWispForTest` 的真 exe，**毫秒级以外的重件**，请在两枚写码腿收口之后再跑；
N5 那三形在 `cmd/wisp/run_mode101_test.go:605-660` 的同一枚父用例里，取父用例名再跑。）
**开工前的基线要求**：以上四行必须在**未改动**的工作树上先取一次绿/红名册，
否则落地腿无法分辨"我改红的"与"本来就红的"。

---

## §5 三格判据的正控形状（AC#1 / AC#2 / AC#4）

### §5.1 AC#1 删掉那份再跑要重新建
- **今天会不会恒真**：分两半。**"真产出一份"今天必然假**（§1.1 R10-R11 会先退出）⇒ 非自明，好判据。
  但 **"二次运行不再走创建分支"今天恒真**——因为创建分支**存在都不存在**，任何断言都无从失败（vacuous pass）。
  这条必须写进验收报告，否则"AC#1 绿"是被白送的。
- **正控形状（三段式，缺一不可）**：
  1. 空 dataDir（且**数据根由 `resolveDataDir` 答出**，别用真 `%APPDATA%`，见下方变异 #3）跑一次 ⇒ 文件存在；
  2. **记下字节与 mtime**，再跑一次 ⇒ 字节**完全相同**、mtime 未变、目录里**无 `.wisp-config-*.tmp` 残留**
     （temp 命名式在 `internal/config/parse.go:201`，残留即 `defer os.Remove` 没跑成）；
  3. **删掉它**再跑 ⇒ 文件重新出现，且**与第 1 段的字节相等**。
- **怎么变异它才响（4 条，每条都对应一种偷懒）**：
  - 让创建分支看**内存标志**（"这进程已经建过了"）而不是磁盘状态 ⇒ 第 3 段红；
  - 让写盘走 `os.WriteFile` 追加／非原子 ⇒ 第 2 段红（mtime/字节变）或第 3 段红（半截文件），
    且 `internal/config/private_acl_windows_test.go:103-128` 那族 ACL 钉失去覆盖面（N10）；
  - 让默认值来自"上次读到的 cfg" ⇒ 第 3 段可能仍绿，**这时只有 §5.2 的字节级比能响**；
  - 让创建发生在 `resolveDataDir` 之前 ⇒ `cmd/wisp/dataroot_128_test.go:343` 的 `assertDirEmpty128` 立刻红（N3），
    **这条不是本腿设计的变异，是现成钉**——所以 AC#1 的删除重跑用例必须在临时数据根里跑。

### §5.2 AC#2 每个值都能在 `schema.go` 的默认里找到出处（要能区分"照抄默认"与"写手自己填了一套"）
- **可机读的法子（推荐形，三段）**：
  1. **期望侧不要用 `NewDefaults()` 造**（`internal/config/defaults.go:58-62`）——实现与判据若同引一枚函数，
     这条 AC 当场**恒真**（这正是"会不会恒真"这一问的答案：**会，而且踩中即废**）。
     改成在测试里**自己反射**一遍：遍历 `Config` 的字段、读 `default:"…"` 标签、自己 `strconv` 解、自己按逗号切
     （即把 `defaults.go:65-133` 的语义在测试里重述一次，两侧同源但**不同码**）；
  2. **实测侧解码进零值 `&Config{}`**，**不要**解码进 `NewDefaults()`：
     `readConfigFile` 走的正是 `loader.go:117-118`（先 `cfg := NewDefaults()` 再 `decodeStrict`），
     所以 `LoadFile` 的返回值**永远看不出"文件里有没有这一行"**；只有裸解码才分得开"写了 true"与"没写"；
  3. **逐项比对**，而不是整体 diff。两把现成的尺可用：
     - 包内：`diffKeyPaths(NewDefaults(), rawDecoded)`（`internal/config/writeguard.go:201-213`，
       它**跳过 `SchemaVersion`**（`:220-224`）与 `toml:"-"`（`:226-228`，但 `Plugins` 例外手工走）——
       ⚠ **未导出**，`cmd/wisp` 侧用不到；
     - 包外（`cmd/wisp` 落地腿可用）：**字节比**
       `MarshalCanonical(NewDefaults())`（`internal/config/parse.go:164`，**已导出**）
       vs `os.ReadFile(cfgPath)`。`schema.go:101-103` 与 `parse.go:160-163` 两处注释都逐字保证"**无 omitempty，每个静态字段都显式落一行**"，
       所以字节比是有意义的（"填了个 false"会留下一行，"没填"不会）。
- **这法子的漏计（必须写进验收报告，四枚）**：
  1. **零值与"没写"不可分**（结构体侧）：§1.4 形状 (b) 那 **103** 枚无标签叶子（口径：173 带 `toml:` 标签字段 － 70 带默认标签；含表头型字段，故为上界）
     默认就是 Go 零值 ⇒ 反射侧对它们**什么也判不了**；只有上面 3.的**字节**那一把能看出多出来的一行。
     其中 70 枚显式默认里有 **14 枚是 `default:"false"`、16 枚是 `default:"true"`**
     （尺：`grep -oE 'default:"[^"]*"' internal/config/schema.go | sort | uniq -c | sort -rn | head -2`）
     ⇒ 这 30 枚正是"填了个假值 vs 干脆没填"最容易蒙过去的地方。
  2. **map 形 section 天生不可判**：`llm.providers`（`schema.go:419`）、`models.local_override`（`:590`）、`plugins.Entries`（`:575`）
     没有标签可反射（`defaults.go:77-78` 刻意留 nil）⇒ 判据对它们**恒真**。这正是 **J2** 的所在。
  3. **装载侧会替文件说谎**：`applyPresets`（`loader.go:188-204`）在装载时把 `providerPresets`（`defaults.go:22-35`）的 12 条
     protocol/base_url 灌进空字段。⇒ **判据若读 `LoadFile` 的结果，会看到文件里根本不存在的一批值**，
     而它们**在 `schema.go` 的默认里没有出处**（出处在 `defaults.go:22-35`）——这就是 AC#2 这句话真正在防的形状。
  4. **`SchemaVersion` 是外来的**：`SaveFile` 强制写入（`loader.go:243`），`schema.go:105-108` 注释逐字"never defaulted"。
     ⇒ 判据必须**显式豁免**它（`diffKeyPaths` 已经豁免，字节比则要求文件里恰好是 `schema_version = 2` 一行），
     否则第一版判据自己就会红。
- **今天会不会恒真**：**不会**（文件根本不生成）。但**方法 1+2 若写成 `equal(file, marshal(NewDefaults()))` 而实现也是 `SaveFile(NewDefaults())`，落地当天即恒真**
  ⇒ 验收腿必须问一句"期望侧是那 70 枚标签，还是同一枚函数"。

### §5.3 AC#4 缺 key 时那句实话
- **今天这句是什么**（实测文本形状）：`cmd/wisp/run.go:412` 打的是
  `wisp run: 文本角色未配置（Unconfigured）：llm: role "chat" is unset and text_chain is empty (Unconfigured)`
  （后半句原文在 `internal/llm/resolver.go:218`）。它**指名了缺什么**（role + text_chain），
  **没指名去哪儿补**（无 `wisp secret set`、无 `api_key_ref`、无文件路径），
  而票面 AC#4 要的是"指名缺什么、**去哪儿补**"两半 ⇒ 今天这句子**半合格**。
- **今天会不会恒真**：⚠ **有一条很危险的现成恒真形**——
  若判据写成 `strings.Contains(stderr, "config.toml")`，**今天就已经绿**：
  缺配置那一路（`run.go:397`）打的是被 wrap 过的 OS 原话，里面天然带文件路径
  （wrap 点在 `internal/config/loader.go:67`）。⇒ **判据不许只查路径字符串**，
  必须查到**祈使句**上去：`api_key_ref`（`cmd/wisp/secret.go:175` 的原文形状）＋
  `wisp secret set`（`secret.go:162`）或 `wisp providers discover`（`main.go:35-37`）**至少其一**。
- **正控形状（成对写，缺一即松）**：
  (i) 首建回执里含"缺：模型/provider 指向 + key"两枚名字；
  (ii) **同一次运行仍退码 2**（把票面 §3"不许顺手放宽"钉进判据本体）；
  (iii) **0 个请求打到 provider**（抄 N6 的 `:466-468` 那一招，它是本仓现成的"没拿假端点凑"的证明）；
  (iv) 输出不含任何 key 明文（抄 N6 的 `:469-472`）。
- **怎么变异它才响**：
  - 把"没配置"当"按默认跑"（放宽退码）⇒ (ii) 与 N1/N5/N6 一起响；
  - 拿 `providerPresets` 里任一 endpoint（如 `ollama → http://127.0.0.1:11434/v1`，`defaults.go:31`）凑一个默认 provider ⇒
    (iii) 响（有请求出门），并且 N7 的守卫与票面"不许新造默认值"同时咬；
  - 只在 stdout 打回执、stderr 不放 ⇒ 用 `providersIO`/`runSpec` 那两枚分流（`cmd/wisp/providers.go:100`、`run.go:412`）的正控可分；
  - 把实话写成"已为你配置好，可以直接用" ⇒ (ii) 仍在但语义撒谎 ⇒ 需要**断言句子里含"仍/还/未"这一族转折词**才拦得住（§3.1 claims 家族 `:101-106` 讲的就是这种"措辞即分界"）。
- ⚠ 本腿不动任何 AC 勾选框（AGENTS.md §1.4／题面硬边界 4）；§5 只是判据形状，勾选归编排者与验收腿。

---

## §6 待人裁 / 判不动的地方

- **J1｜首建该长在哪一层：`assembleRuntime` 之内还是 `runTextTask` 之前。**
  代码事实：生产侧调用 `assembleRuntime` 的只有两处——`cmd/wisp/run.go:234`（`wisp run`）与
  `cmd/wisp/resident_task_source_windows.go:265`（常驻腿取任务）。退码 2 那句 `wisp run: 配置未就绪
  （Unconfigured）` 在 `cmd/wisp/run.go:397-398`，即首建若插在 `config.NewManager` 之前，
  **两枚宿主同时被改**；插在 `runTextTask` 则只改 `wisp run`，常驻腿继续今天这形。
  我判不动的是**授权范围**：票面 §"要建什么" 1 只说"找不到配置时"，没点名宿主枚数，而"全新机器第一天能跑"
  这句话对常驻腿成不成立不归我判。**⚠ 读完 §4 名册之后这一格收窄了**：钉 **N2**
  （`cmd/wisp/resident_task_source_246_windows_test.go:389-391`）**明文要求常驻腿的 dataDir 里 config.toml 保持不存在**，
  措辞是"the leg created a config.toml it was never asked for" ⇒
  **首建放 `assembleRuntime` 里 = 今天就有 3 处断言转红**（`:369`／`:385`／`:389`）。
  于是真正待人裁的已经不是"放哪一层更安全"，而是：**允不允许改判 N2 那枚用例**（改它＝动别的票的 AC 判据，
  本腿与落地腿都无权做）。需要**编排者裁**这一句；裁"允许改 N2"才能谈两宿主一起，否则首建只能放 `wisp run` 专属那一层。
- **J2｜AC#1 那句"带全部 section 的默认 `config.toml`"与 AC#2"不许新造默认值"今天互相拉扯。**
  代码事实：`NewDefaults()`（`internal/config/defaults.go:58-62`）**刻意把所有 map 字段留 nil**
  （`defaults.go:77-78` 注释逐字：`leave nil (see NewDefaults)`），于是
  `llm.providers` / `models.local_override` / `plugins.Entries` 这几族**没有条目可写**——
  写出来的文件里这些表要么不出现、要么是空表。"全部 section"若按字面理解成"每张表都在文件里"，
  落地腿就得**自己填条目**，那正好撞 AC#2 的"不许新造默认值"。我判不动哪种读法作数，
  需要**编排者或 owner 裁**（我倾向"18 枚静态 section 全在、动态表按 nil 留空"这一形，理由见 §1-§1.4）。
- **J3｜`SaveFile(NewDefaults())` 落盘后文件里到底出现几行、哪些表在不在，本腿量不了。**
  这是"必须真跑才答得出"的一格。⛔ 本腿禁 `go test`/`go build`/`go run`（两枚写码腿在跑毫秒级计时）。
  交**待落地腿自量**，确切命令与期望读数：
  `go test ./internal/config -run TestUnwiredGuardLeavesHonestConfigsAlone -v`
  → 期望 PASS（`internal/config/unwired_test.go:86-95` 这一子用例今天已经就是
  "SaveFile(path, NewDefaults()) 之后再 LoadFile 要能读回"的钉，它绿 ⇒ 首建机制的可加载性已被现成用例钉住）；
  量文件形状请另写一次性探针（`t.TempDir` + `SaveFile` + `os.ReadFile`），期望
  静态 18 section 各出现一次、`[llm.providers.*]` 与 `[plugins.<id>]` 零次、`schema_version = 2` 一行。
- **J4｜首建该不该同时把 `[llm].text_chain` 或某个 provider 名字写进去。**
  写了就是"拿一个假端点凑"，票面 §3 明令不许；不写则首建之后 `wisp run` **仍然退码 2**
  （撞点是 `cmd/wisp/run.go:409-413`，`llm.Resolver.ResolveRole` 在
  `internal/llm/resolver.go:218` 报 `role %q is unset and text_chain is empty`）。
  也就是说 **AC#4 才是本票真正的用户可见终点，首建本身不改变退码**。这条我判得动（结论见 §3/§5-§5.3），
  但"首建之后还缺什么那句话该由谁打印"（`assembleRuntime` 里还是 `runTextTask` 收口处）需要编排者裁。
- **J5｜票 132 `SealDir` 一族：AC#3 说的"沿用已定的那一形"今天生产侧根本没人走 `SealDir`。**
  事实见 §2-§2.3（显式生产调用者 **0 枚**，能力全部由 `PrivateDirAll`/`SealFile` 投递）。
  我判不动的是**验收怎么算过**：落地腿若用 `winsec.PrivateDirAll` 建目录，AC#3 那句"票 132 的 `SealDir` 一族"
  是按"符号名命中"判还是按"同一判定者投递"判？前者会把正确实现判错。需要**编排者裁**，
  验收腿照裁过的口径读。
- **J6｜AC#3 那句"不许 `filepath.Join` 手拼"与禁令的真实射程不一致，判不动该按哪句执行。**
  代码事实：`AGENTS.md §1.2` 抄的禁令（权威在 `PLAN.md:1288`）禁的是 **`risk.PathResolver` 之外用
  `filepath.Clean|Abs` 做文件系统决策**，而**现存的生产代码本来就在用 `filepath.Join` 拼落点**：
  `cmd/wisp/run.go:371`（`cfgPath := filepath.Join(s.dataDir, configFileName)`）、
  `cmd/wisp/doctor.go:249/251/253/265/267`。
  票面 AC#3 却写作"落点由现成判定者给（**不许 `filepath.Join` 手拼**）"。
  ⇒ 两句话不能同时成立：**按票面字面判，落地腿连"复用现成的 R8 那一行"都算违规**；
  按 `PLAN.md:1288` 判，`Join` 没问题、只要根是判定者给的。
  我判不动按哪句执行（这决定了落地腿要不要把 `:371` 一起换形），需要**编排者裁**；
  本腿的倾向：**以 `PLAN.md:1288` 的射程为准**（`AGENTS.md` 自己声明"与那五份文件不一致时以那些文件为准"），
  并且**不要为了避开 `Join` 这个词去新造一枚路径 API**——那才是本票唯一真正的"第二真相源"风险。
- **J7｜`wisp doctor` 今天查数据根、但**一枚也不查 config.toml**，要不要在本票补。**
  事实（尺本腿已跑：`grep -oE '(pass|fail|info)\("[^"]*"' cmd/wisp/doctor.go \| grep -ci config` ⇒ **0**）：
  `cmdDoctor`（`cmd/wisp/doctor.go:29`）的检查项里有 `data dir resolvable (`与 `data dir writable (`两枚
  （即 R5 那一跳**有**自检出口），但**没有任何一枚命名 config.toml**；
  "config.toml" 在本文件中只出现在注释 `:243` 与 `:295`。
  ⇒ 首日排障最自然的那句"配置到底在不在"目前**无人报**。
  但票面 §"要建什么" 只列了两件（首建 + 引导），没列 doctor ⇒ **我判不动这是"票 198 顺手做"还是"另立一票"**，
  需要**编排者裁**（并归入 `docs/reports/pending-and-issues.md` 的哪一族，本腿不写台账）。
- **J8｜`cmd/wisp/secret.go` 的 `configFileName` 常量与首建文案的归属。**
  `configFileName = "config.toml"` 定义在 `cmd/wisp/secret.go:63`（不在 run.go、不在 config 包），
  首建路径今天唯一拼它的地方是 `cmd/wisp/run.go:371`。新建文件若复用该常量则跨文件依赖 secret.go，
  若新写一份则造出第二处真相。这属"改哪枚文件"的派单细节，我判不动要不要顺手搬家，
  需要**编排者裁**（本腿倾向：复用现常量，搬家不在本票）。

## §7 本腿推翻编排者题面之处

题面给的 4 处现量 + 那把 grep，逐枚复认（判定依据在 §1/§2/§4，此处只记结论）：

| # | 题面（与票面同源） | 判定 | HEAD 实测 |
|---|---|---|---|
| X1 | "票 198 …… **48 行**" | **推翻** | **35 行**（`wc -l .scratch/wisp/issues/198-*.md` = 35）。编号存在、文件名一致，仅枚数错 |
| X2 | `internal/config/schema.go:104-134` ＝ **16** 个 section | **行号成立、枚数推翻** | `104-134` 正是 `type Config struct` 全体；但区间内 section 字段实测 **18 枚**（`schema.go:110-133`，含 `Plugins`）。仓内 `type *Section struct` 亦 **18 枚**。见 §1-§1.3 |
| X3 | `internal/config/manager.go:76-84` ＝ `NewManager` 注释"首建是调用方的事" | **结论成立、行号已漂移** | HEAD：注释在 `manager.go:99-100`，函数体 `99-111`；`76-84` 此刻是 `LockedDecision`。查证：`git show 335b8d2b:internal/config/manager.go` 里 `first-run file creation is` 确在 **76 行** ⇒ 票写时成立，属**过期**不是错写 |
| X4 | `cmd/wisp/run.go:282-291` ＝ 缺配置退码 2 | **行号错（已漂移）、结论成立** | HEAD 的 `282-291` 是 `agentRuntime` 的结构体字段（`endpoint`/`gate`/`ui`/`bridge`…），与退码无关。真链：`run.go:389` `config.NewManager(cfgPath, nil)` → `run.go:397-398` 打印 + `return rt, 2`；退码表在 `run.go:68-91`（`ClassConfig`/`ClassAuth` → 2）。同样查证：09-28 锚点 `335b8d2b` 的 `282-291` 落在 `assembleRuntime` 的注释上（当时近似成立） |
| X5 | `internal/config/loader.go:39-43` ＝ 缺配置时退码 2 | **半成立，指错了产生点** | `loader.go:37-41` 是 `LoadFile` 的文档注释（"the caller maps it to the Unconfigured state"），`:41-43` 是 `LoadFile` 签名。**"缺文件 ⇒ 错"真正被造出来的是 `loader.go:65-68`**（`os.ReadFile` 失败 → `observe.Wrap(observe.ClassConfig, err, "config.toml read")`）。loader 从不产生退码 |
| X6 | `cmd/wisp/secret.go:263-271` ＝ key 只能 CLI 手录 | **成立（HEAD 即为该区间）** | `263-271` 正是 `case "set"/"get"/"list"/"unset"` 的子命令分派；`configFileName` 常量在 `secret.go:63`，`api_key_ref = "dpapi:<name>"` 文案在 `secret.go:175`、`secret.go:380` |
| X7 | 编排者 08:4x 现跑 `grep -rn "config\.toml" cmd/ internal/ --include=*.go \| grep -v _test` ＝ **12 行命中** | **枚数推翻（差一个数量级）；结论成立** | 本腿逐字复跑同一把尺：**163 行命中**（口径：`cmd/`＋`internal/` 全树非 `_test.go`，含 `cmd/balldebug`）。"没有一处真创建配置文件"这句**方向对但说法要收窄**：见下 X8 |
| X8 | 票面"全仓**无一处创建 `config.toml`**"（`Status` 与现量表第 2 行） | **推翻（存在一处生产创建点，只是首达不到）** | `internal/config/writeguard.go:125-131` 有 `case fileMissing(statErr)` 分支，注释逐字写着"writing it creates what first-run did not"，并落到 `writeguard.go:160` 的 `SaveFile(m.path, base)` ⇒ **`Manager.mergeWrite` 今天能在文件缺失时把它建出来**。它之所以救不了全新机器，是因为 `NewManager`（`manager.go:101-105`）在装载失败时直接返回错误，生产根本走不到 `mergeWrite`。正确的说法是：**首建能力已存在且已导出，缺的只是入口那一次调用** |
| X9 | 题面"internal/config 今天有没有**任何**写盘能力"（暗示可能没有） | **推翻暗示：有，而且三枚都是现成的** | `SaveFile`（`loader.go:238`，**已导出**）· `MarshalCanonical`（`parse.go:164`，**已导出**）· `atomicWrite`（`parse.go:196`，未导出，内含 `winsec.SealFile` 于 `:215`）。落地腿**不需要新写序列化器**，因此题面担心的"默认值改一处要同步两处"这枚维护面**不会被引入**（论证见 §1-§1.5） |

> X3/X4 的"过期 vs 错写"分界是查出来的，不是猜的：`git show 335b8d2b:<file>` 的读数已附在上面两行。
