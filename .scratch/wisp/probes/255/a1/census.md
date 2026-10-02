# 255-a1 只读取证 —— 「生效级别」词表 census（供编排者裁 AC#2 的 ⓐ／ⓑ）

**代号**：`255-a1` · **性质**：只读取证腿 · **唯一可写件**：本文件
**要答的那一个问题**：今天仓里「生效级别」到底有几套词表、每套谁生产谁消费、彼此有没有映射；
ⓐ（让 `Tier` 真被读者消费）与 ⓑ（具名登记「粒度只到段」+ 一把会响的尺）各要动哪几枚文件、撞不撞现成的钉。
**本格判据由编排者裁，本腿未裁。**

---

## §0 起手锚（逐字读数）

| 尺 | 逐字读数 |
|---|---|
| `date -Iseconds`（进场第一发） | `2026-10-02T10:43:11+08:00` |
| `git log -1 --format='%H %ad %s'` | `a4b221ac098d35c6cf81ac60d3063cf53b9d6b8c Fri Oct 2 10:41:23 2026 +0800 evidence(248-v1b) 骨架：接管腿 248-v1b 的 §0 起手锚落盘……` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git status --porcelain internal/config internal/panel cmd/wisp \| wc -l`（10:43:11 那一发） | `0` |
| 同上复跑（`date` 现取＝`2026-10-02T10:44:35+08:00`） | `1` |
| `git status --porcelain internal/config internal/panel cmd/wisp`（列出那一枚） | ` M internal/panel/config_handlers.go`（`od -c` 复认首列是空格＝**未暂存的修改**） |
| `git status --porcelain cmd/wisp internal/config \| wc -l` | `0` |

**⚠ 本腿取证期间落进工作树的一枚外部改动（不是我改的，我一字节未动产码）**：
脏的那一枚是 `internal/panel/config_handlers.go`，内容是一枚**植入变异**（正文逐字 `return true // MUTATION-M1 arbitrary key passthrough`，
把 `knownWritableField` 的表查询整段替成恒真）。`git diff --stat` 逐字：`1 file changed, 1 insertion(+), 6 deletions(-)`。
⇒ 归属判断（**射程判断，非内容引用**）：这正是对抗验收腿 `248-v1b` 的活儿，与票 255 排程节「此刻 `248-v1` 正在整包跑 `cmd/wisp`」同源。
⇒ **对本腿读数的影响（已逐条核过）**：变异落在 `:368` 之后，故本件引用的 `EffectiveTier` 定义区（`:193–202`、`:213`）
在脏树与 HEAD 里**行号全等**；`tierSentence` 一枚则不等——现尺：
`git cat-file blob HEAD:internal/panel/config_handlers.go | grep -n 'func tierSentence'` ⇒ **`428`**；
`grep -n 'func tierSentence' internal/panel/config_handlers.go` ⇒ **`423`**。
⇒ **本件凡引 `internal/panel/config_handlers.go` 的行号，一律以 HEAD blob 为准**（HEAD 是不动的锚，工作树在取证期间会漂）。

**取证面写面为空的复认**：`cmd/wisp` 与 `internal/config`（票 255 排程节点名的两枚写面）在本腿每一发尺时都是 0 行脏，
故 ⓐ／ⓑ 两形的代价表读的是**已提交形状**，不是谁的半成品。

**根目录纪律**：本件所有 `grep`／`find` 的根**逐条显式写作** `cmd internal tools scripts`（Go 取证面）或
`.scratch docs`（件面），⛔ 未用 `.` 当根，⛔ 未触及 `frontend/**`／`design/**`（两层禁令：不读不引）。

**禁跑尺的遵守**：本腿未跑 `go build`／`go vet`／`go test`／任何 `./...`；未跑 `go list`（依赖边全部由 `grep` 的
`import`／`pkg.` 前缀现读得出，不需依赖图，故按票面「用不上就别用」省掉）。

本节以下 §1／§2／§3 的表格在后续 commit 填写；本腿按票规**先写满 §4、§5 再回填**。

---

## §1 「生效级别」词表的现状表

### §1.0 先答编排者点名的形状：**不止三套，也不是四套，是「四套数档位 + 一套同名异事」**

编排者的锚 4 猜的是「三套互不相干的枚举（`config.Tier` 三值 / `panel.EffectiveTier` 四值 / `restartTierKeys` 三枚段名）」。
现量结果：**那三套都真实存在，但都不是今天在跑的那一档**；真正在跑的档位口径是**第四套**，
而它**通篇不写 "tier" 这个词**——这就是我的尺（也是编排者的尺）最容易漏的一枚。

**漏它的机制**（具名，因为它决定 AC#2 该往哪打）：第四套的"档名"不存在任何字符串里，
它存在**「一枚段名被 append 进了 `Report` 的哪一枚切片」**里。尺：

```
grep -rn "\.Hot\b\|\.Reload\b\|\.Restart\b" --include=*.go cmd internal tools scripts | grep -v _test
```
⇒ 逐字命中 11 行，其中写侧 5 枚全在 `internal/config/manager.go`（`:293 :349 :355 :412 :415`），
读侧 4 枚在 `cmd/wisp/config_reload.go:169,170,171` ＋ `internal/config/manager.go:198,201`，
外加 `:270` 一行注释。**三枚切片名 `Hot`/`Reload`/`Restart` 与 `TierHot`/`TierReload`/`TierRestart` 三枚常量逐字同名、语义同构、但零处相连。**

### §1.1 现状表（每套一行，四问；每格带尺与真实读数）

尺的公共前缀：所有 `grep` 均为 `grep -rn <型> --include=*.go cmd internal tools scripts`，测试侧另加 `| grep -v "_test.go"`；
`internal/panel/config_handlers.go` 的行号**一律取 HEAD blob**（§0 具名的漂移原因）。

| # | 词表 | 定义处 | 非测试写者 | 非测试读者 | 跨套映射 |
|---|---|---|---|---|---|
| **T1** | `config.Tier` ＝ `"hot"`/`"reload"`/`"restart"`（票面说的"三档"） | `internal/config/schema.go:31` 类型；`:36 :39 :43` 三枚常量 | **0 枚** | **0 枚** | **无** |
| **T2** | `restartTierKeys` ＝ 3 枚 `app.*` 键名字面量表 | `cmd/wisp/config_reload.go:299-301` | **1 枚**（定义即写） | **2 枚**，同一函数体内：`:287`、`:293`；**包外 0 枚** | **无**（注释自称"给人比"的替身，见 §1.3） |
| **T3** | `panel.EffectiveTier` ＝ `"now"`/`"next_task"`/`"restart"`/`"not_applied"` | `internal/panel/config_handlers.go:196` 类型；`:199-202` 四枚常量；载体 `:213` `SettingWriteResult.Tier` | **11 枚点位**，全在 `cmd/wisp/panel_config_store.go`（`EffectiveNotApplied` 9 枚：`:175 :177 :191 :198 :205 :212 :219 :225 :239`；`EffectiveRestart` 2 枚：`:228 :275`）⇒ **但可达值只 2 枚**（见 §1.2） | **3 枚**，全在 `internal/panel/config_handlers.go`：`:391`、`:419` 各调一次 `tierSentence(r.Tier)`；`:428` 本体（4 case ＋ default） | **无** |
| **T4** | **`config.Report` 的三枚 `[]string` 切片**——今天真正在跑的档位口径 | `internal/config/manager.go:87-97`（`Hot :89`、`Reload :91`、`Restart :94`、`Locked :96`） | **5 枚 append 点**：`:293`（`rest` 表 12 段共用一枚）、`:349`/`:355`（`planApp`）、`:412`/`:415`（`planVoice`）；锁定四段走 `planLocked` ⇒ 不进三档，走 `Locked` | **4 枚**：`cmd/wisp/config_reload.go:169-171`（`reportReload`，就是 AC#1 那句"这些段已立即生效"的产生处）、`:115` 装 `:278 reportRestartPending`（读 restart 档）、`internal/ball/hotkey_reload.go:111 Sections()`（按段名字符串 `"hotkey"` 认 reload 名单）经 `cmd/balldebug/main.go:244` 装、`internal/config/manager.go:198-202` 是转发不是读者 | **无，且结构上不可能有**：T4 里**没有档名这个字段**，档位＝"落进哪一枚切片"，与 T1 只有**命名上的对应**（`Hot`↔`TierHot` 等三对逐字同名），**零处代码相连** |
| （T5） | 同名异事，**不计入档位口径**：`risk.Tier int` ＝ `TierNone`/`TierB`/`TierA`（敏感路径档，`SPEC-06 §4.1`） | `internal/risk/assessor.go:147-153` | 1 枚（`Classify` 实现体） | **3 枚真读者**：`internal/tools/paths.go:312,315,317,319` | 不适用。**单列的理由**：ⓐ 形若落地要在仓里新增"读 `Tier`"的代码，而这枚标识符已同时指 T1/T5/`projctx.Tier`（`internal/projctx/projctx.go:45-48`，`"project"`/`"global"`）/`PluginsSection.Tier2Enabled`（`schema.go:573`，D19 信任层）**四件事** |

### §1.2 三套档位词表彼此「说不说得上话」——逐条给答案

**（a）T1 与 T4：只共享三个名字，不共享一行码。**
`schema.go:30` 的注释逐字：`// Tier is the effect level of a config change (D36 three-tier semantics).`
`manager.go:88` 的注释逐字：`// Hot lists sections hot-applied.` ⇒ 两行注释互相指认同一件事，
但 `manager.go` 全文**没有一处出现 `Tier` 这个标识符作为类型或值**（尺：`grep -rn "Tier" --include=*.go internal/config | grep -v _test`
⇒ 12 枚命中里除 `schema.go:30,31,34,36,37,39,40,43` 这 8 枚 T1 自身，其余 4 枚全是别的事：
`manager.go:491,493` 与 `parse.go:102,182` 是 `Tier2Enabled`（T5 系），`schema.go:551,573` 是注释与字段名）。
⇒ **T1 是一枚纯装饰类型：定义了三档、命名了三档、没有一处代码经由它做判断。**

**（b）T2 与 T4：靠人眼对齐，没有任何一把尺在机器上比过。**
`config_reload.go:296-298` 注释逐字：
`// restartTierKeys names the keys applyApp books as restart-tier (D36, PLAN.md`／`// :2715's [app] split). Kept as data so the sentence above and manager.go's`／`// split can be compared by a reader - and by a test - instead of drifting.`
⇒ 它承诺的比对对象是 `planApp` 里的这三枚字段（`manager.go:345-347`）：
`cur.App.Language != fresh.App.Language || cur.App.Autostart != ... || cur.App.SingleInstance != ...`。
今天两边**逐字对得上**（3 对 3）。但"instead of drifting"这半句是**未兑现的**：
`grep -rn "restartTierKeys" --include=*.go cmd internal tools scripts` ⇒ 4 枚命中全在 `config_reload.go` 一个文件内，
**测试零引用、包外零引用**。⇒ 注释里那句 "and by a test" 目前是一句**没有 test 的话**。
⚠ 这一格对编排者裁 AC#2 直接要紧：**T2 已经是 ⓑ 形想要的"具名登记"半成品，只差一把尺。**

**（c）T3 与 T4／T1：零接触，且 T3 的值是硬接的常量。**
`panel_config_store.go:226-229` 逐字（成功支的唯一出口）：
`res.Written = written` ／ `res.Tier = panel.EffectiveRestart` ／ `return res, nil`。
⇒ 写入腿**不看键、不看段、不查 `Report`、不查 `restartTierKeys`**：任何一次成功写入一律回答"要重启"。
⇒ 于是 `EffectiveNow`/`EffectiveNextTask` 的**非测试写者＝0**
（尺：`grep -rn "EffectiveNow\|EffectiveNextTask" --include=*.go cmd internal tools scripts | grep -v _test`
⇒ 4 枚命中＝`config_handlers.go:199,200` 定义 ＋ `:430,432` `tierSentence` 的两个 case，**无写者**）。
⇒ 后果（这是 AC#1 那句误报的**面板侧镜像**，同一枚病）：面板上永远说"这一项要重启进程并重新运行才算用上"，
而 `internal/config` 那边 `[llm]` 其实坐在 `rest` 表的 hot 档里（`manager.go:281`）——**两套词表对同一次写入给出相反的答复，且没有一行码翻译彼此。**

**（d）跨套映射的包侧可行性（给 ⓐ 的落点用）：现成的边够用，不必新开。**
尺：`grep -rn "wisp/internal/config" --include=*.go internal/panel` ⇒ **零命中**，`internal/panel` 今天不依赖 `internal/config`（与票 255 AC#4 的"⛔ 不许新开 `panel→config` 依赖边"一致）。
尺：`grep -rln "wisp/internal/config" --include=*.go cmd internal tools scripts | grep -v _test` ⇒ 14 枚文件，
其中 **`cmd/wisp/panel_config_store.go` 同时看得见两包**（它 import `config` 又产出 `panel.SettingWriteResult`）。
⇒ **任何 T1/T4 → T3 的映射只可能落在 `cmd/wisp`（装配根），落在 `internal/panel` 里就要新开那条被禁的边。**

### §1.3 「段」这一级的覆盖率读数（ⓐ／ⓑ 共同的底座）

`Config` 顶层段：`internal/config/schema.go:110-133` ⇒ **18 枚**（`app ball hotkey session voice audio llm agent risk fs net privacy memory panel cost plugins models observe`；
其中 `Plugins` 带 `toml:"-"` 但由序列化器手工挂 `[plugins]`，见 `writeguard.go:193-197`）。

`plan()`（`manager.go:243-296`）的归属：`planLocked` 4 枚（`:248,252,256,260`＝risk/fs/net/plugins）＋ `planApp` 1 ＋ `planVoice` 1 ＋ `rest` 表 12 枚（`:277-288`）＝ **18**。
⇒ 基数对得上；**但"一枚段一个档"不成立**：`app` 同时产出 hot（`:349`）与 restart（`:355`）两支、
`voice` 同时产出 reload（`:412`）与 hot（`:415`）两支 ⇒ **段级双档已在跑，只是它写在函数体里、不在任何数据里**（§4 的 K6 就是这一条）。

⇒ **这一行是 ⓐ／ⓑ 的分水岭读数**：票面锚 3 说"段名硬编码、没有一处吃 `Tier`"，复认得上；
但真正的形状比"没吃 Tier"更硬——**`app` 与 `voice` 已经在按段内键分档了（theme→hot、language/autostart/single_instance→restart；
pipeline→reload、thresholds/speed/punctuation→hot），只是这套分档只存在于 `planApp`/`planVoice` 的 `if` 条件里，
出了这两个函数体就没有任何一处能读到它。**

---

## §2 ⓐ 形代价：让 `Tier` 真被读者消费（具名到函数）

### §2.1 要动的文件（按"最小可信 ⓐ"与"完整 ⓐ"两档给，逐枚具名）

**最小可信 ⓐ（只让 T1 被读，不改 T4 的载荷形状）——3 枚文件**
1. `internal/config/schema.go`：新增一枚**段（或键）→ `Tier` 的数据表**，把 `planApp`/`planVoice` 里那两句 `if` 的结论搬成可枚举数据。T1 类型与三枚常量已在 `:31-43`，⛔ 不需要新增枚举。
2. `internal/config/manager.go`：`:277-288` 的 `rest` 表把 `"ball"` 这类裸段名换成／附加查表得到的 `Tier`；`:293` 的 `rep.Hot = append(...)`、`:349/:355`、`:412/:415` 五枚写点改为**按表分派**；
   `planApp`（`:338-357`）与 `planVoice`（`:361-417`）的分档 `if` 退化成"查表＋比对"，**函数体本身要重写**（这两个函数是 ⓐ 的正身，票面"粒度只到段"就是指它们）。
3. `cmd/wisp/config_reload.go`：`reportReload`（`:167-171`）读的是 `rep.Hot/Reload/Restart` 三枚切片 ⇒ 只要 T4 的**载荷形状不变**（仍是 `[]string` 段名），这一枚可不动。

**完整 ⓐ（还要让面板那套不再是硬接常量）——再加 3 枚**
4. `cmd/wisp/panel_config_store.go`：`:226-229`、`:275` 那两枚无条件 `res.Tier = panel.EffectiveRestart` 改为**查 T1/T4 后映射**（映射函数只能放这里，理由＝§1.2(d)）。
5. `internal/panel/config_handlers.go`：`tierSentence`（`:428`）的四值口径要不要跟着 T1 的三值收敛——**⚠ 这一动会撞 C17/契约面的边**（票 255 禁区第 2 条把 `SettingsView`/`renderSettingsView`/`Snapshot` 点名成 C17 面，`tierSentence` 是 `renderSettingReceipt` 的下游，属"要新增面板快照字段／新方法名"那一族边缘）⇒ **若走这一枚，本腿具名上报：先落 `A##` 批准记录才许动。**
6. `internal/ball/hotkey_reload.go`：`Sections()`（`:111`）按裸字符串 `"hotkey"` 认档，注释（`:26-31`）已经在用 "HOT-tier"/"RELOAD-tier" 讲道理 ⇒ ⓐ 之后这枚应改读档位数据。**注意它的注释自陈一个事实：`[hotkey]` 是 hot 档，而 `OnReload` 只为 reload 档响 ⇒ 这条回调路径今天根本不会因 hotkey 改动而响**（`hotkey_reload.go:28` 逐字："a hotkey edit would never reach a callback installed there"）。

**新增包级依赖边：0 枚。**（`internal/config` 不需要新 import；`cmd/wisp` 已同时 import 两包；`internal/panel` 保持零 `config` import，不碰 AC#4 那条禁令。）

### §2.2 ⓐ 撞哪些现成的钉（逐枚读断言本体，写"改语义的哪句会打红"）

尺：`grep -rln "Tier" --include=*_test.go cmd internal tools scripts` ⇒ 10 枚文件；再逐枚读断言体。

| 钉 | 位置 | 断言本体（逐字关键句） | ⓐ 会打红的哪一句 |
|---|---|---|---|
| N1 | `internal/config/manager_test.go:58` `TestManagerHotTierAppliesImmediately` | `if rep == nil \|\| !reflect.DeepEqual(rep.Hot, []string{"ball"})` | **钉死 T4 的载荷形状＝`[]string` 裸段名**。ⓐ 若把 `Hot` 的元素改成带档名的结构（哪怕只是加个 wrapper）⇒ 这行**编译期即红**。保住它的唯一办法＝`Hot` 形状不动、档位另开一枚字段。 |
| N2 | `manager_test.go:78` `TestManagerRestartTierKeepsOldValues` | `!reflect.DeepEqual(rep.Restart, []string{"app"})` ＋ `m.Config().App.Language != "zh-CN"` ＋ `!reflect.DeepEqual(pending, []string{"app"})` | 三枚一起钉死"restart 档＝不替换内存值＋回调带段名"。ⓐ 把 restart 判定搬到数据表后，**只要表里 `app` 的档位标错一格，第一枚断言红**（这是好事，等于 ⓐ 自带正控）。 |
| N3 | `manager_test.go:97` `TestManagerThemeIsHotInsideApp` | `!reflect.DeepEqual(rep.Hot, []string{"app"}) \|\| len(rep.Restart) != 0` | **这枚钉把"段内分档"这件事钉死在 `planApp` 的函数体里**：它要求 theme 改动时 `Restart` 必须**空**。ⓐ 若把 `app` 整段登记成 restart 档（粒度真"只到段"），此钉**必红** ⇒ 反过来说：**ⓐ 想活，就必须把分档做到键级，那就不只是"给 T1 找读者"，而是重写 `planApp` 的键级条件式。** 这是 ⓐ 最贵的一格。 |
| N4 | `manager_test.go:113` `TestManagerReloadTierEmitsEvent` | `!reflect.DeepEqual(rep.Reload, []string{"voice"})` ＋ `fired` 回调同样比对 ＋ `m.Config().Voice.ASR.Model` | 同 N3，对 `[voice]` 的 reload/hot 双档。ⓐ 要把 `planVoice`（`:366-411`，一枚 `if` 里塞了 11 个条件、外加一枚 hot 的 4 条件 `if`）搬成数据 ⇒ **这两枚断言是搬完必须逐字仍成立的验收面。** |
| N5 | `cmd/wisp/config_reload_223_test.go:274` | `if !strings.Contains(r.h.out.String(), "这些段已立即生效")` | ⓐ 动了 `reportReload` 的读法 ⇒ 这句话的来源变了，文案断言可能仍绿（它只 `Contains` 片段），⛔ **但这一枚正是 AC#1 的靶心，ⓐ 与 AC#1 会抢同一枚函数** ⇒ 排程冲突，见 §2.3。 |
| N6 | `config_reload_223_test.go:495` | `if strings.Contains(out, "这些段已立即生效：[app]")` | 反向钉：restart 档不许冒充立即档。**ⓐ 若把档位查表算错（把 `app` 标成 hot），这枚与 N3 同时红** ⇒ ⓐ 的错有钉兜着，这一格对编排者有利。 |
| N7 | `config_reload_223_test.go:462` `TestTicket223RestartTierSaysItWillNotApply`，关键句 `:473 awaitAudit(t, "config: RESTART-PENDING detail=")`、`:482` 三枚 needle `app.autostart`/`开机自启`/`重启进程后生效` | **这枚钉把 T2 `restartTierKeys` 的内容当合同钉住了**（那三枚键名要出现在句里）。⇒ ⓐ 若把 T2 废掉、改由 T1 数据表供给那句话，**断言仍可绿，但前提是表里那三枚键名逐字不变**。⛔ 也就是说：**T2 不能随手删，它是这枚钉的字符串来源。** |
| N8 | `internal/panel/config_route_248_test.go:323`（函数 `TestAC8ReceiptNamesTheTierAndNeverClaimsLiveSave`） | `for _, tier := range []EffectiveTier{EffectiveRestart, EffectiveNextTask, EffectiveNotApplied, EffectiveTier("")}` ＋ 每支都断 `line` **不含** `"保存即生效"`／`"立即生效"`；`:340` 又单独要求 `EffectiveNow` 那支**必须含** `"立即生效"` | ⓐ 若收敛 T3 的四值（例如砍掉 `EffectiveNextTask`）⇒ **这枚测试直接编译不过**（它按名字枚举）。`:340` 那一支还把"now 档是唯一可说立即生效的地方"钉死 ⇒ §1.2(c) 说的"now 无生产者"这枚病，ⓐ 想靠"干脆去掉 now"来省事是走不通的。 |
| N9 | `internal/config/boundary_test.go:168` `c.Plugins.Tier2Enabled = true`；`internal/risk/*_test.go` 5＋2＋3 枚；`internal/projctx/projctx_test.go` 1 枚；`cmd/wisp/instructions_200r2_test.go` 6 枚 | — | **不撞。** 逐枚点明：它们钉的是 `Tier2Enabled`（D19 信任层）／`risk.Tier`（T5）／`projctx.Tier`（来源标签）三件同名异事 ⇒ ⓐ 对 T1 动手碰不到它们。**列在这里是为了说明：ⓐ 的爆炸半径可以只圈 T1/T4/T2/T3 四套，不必替另外三件同名者背书。** |

### §2.3 ⓐ 的排程与写面冲突（本腿量到的，不是判断）

- 票 255 排程节写：写面＝`internal/config`＋`cmd/wisp`，⛔ 与正在整包跑 `cmd/wisp` 的腿并发。⇒ **ⓐ 的最小可信版就要同时开 `internal/config/schema.go`、`internal/config/manager.go`、`cmd/wisp/config_reload.go`（完整版再＋`panel_config_store.go`、`internal/panel/config_handlers.go`、`internal/ball/hotkey_reload.go`）＝3～6 枚文件、跨 3 个包。**
- ⓐ 与 **AC#1 抢同一枚函数**（`reportReload`，`config_reload.go:167-171`）：票面建议把 AC#1＋AC#3 拆成 `255-r1`，AC#2 待裁后另排 ⇒ **若裁 ⓐ 且不做 §2.3 的串行，两枚腿会同时改 `config_reload.go` 与 `manager.go`。**
- ⓐ 若要动 `internal/panel/config_handlers.go` 的 `tierSentence`／`SettingWriteResult`，**按票 255 禁区第 2 条属 C17 面 ⇒ 必须先停手上报、由编排者落 `A##`**（本腿在此具名，不代为裁）。
- ⓐ 的**净收益读数**（给编排者称量用）：ⓐ 之后 §1.2(c) 那笔"面板永远说重启、config 那边 `[llm]` 其实是 hot 档"的相反答复**只有走完整 ⓐ（第 4 枚文件）才消**；走最小可信 ⓐ 只让 T1 有了读者，**T3 那套仍然硬接常量、仍然与 T4 相反**。

---

## §3 ⓑ 形代价：具名登记"粒度只到段" ＋ 一把会响的尺

### §3.1 那把尺落在哪一枚包（读答案：`internal/config`，测试同包；且**仓里已有同形尺可抄**）

**现成可抄的那一枚**（本腿认为这条读数直接决定 ⓑ 的便宜程度，具名给编排者）：
`internal/config/unwired_test.go:311` `TestEveryLockedSectionKeyIsAccountedFor` —— 它已经是 ⓑ 要的形状，只差射程：
- 反射枚举叶子键：`collect()` 闭包（`:313-341`），`struct` 递归、`map` 展成 `path+".<id>"`、其余算叶子；
- 登记表：`internal/config/unwired.go:119` `lockedKeyDisposition map[string]string`（键路径 → 一句话理由，如 `"consumed: cmd/wisp/run.go builds the approval timeout from it"` 与 `"unwired:risk.shell_enabled"`）；
- 响的条件：`:352-355` 取登记、`if !ok`（`:353`）⇒ `t.Errorf("locked key %q has no disposition ...")`（`:354`）；`:357-368` 双向校验（登记指向的 guard 必须真在 `unwiredKeys` 里，反之 `:369-373` guard 不被任何键指也红）。
- **射程限制**（ⓑ 要补的正是这一格）：`:342-345` 只对 `risk/fs/net/plugins` 四段各调一次 `collect` ⇒ **段外 14 枚段零覆盖**，与票面现量第 4 条"段外 0/89＝0%"是同一枚事实的两个面。

⇒ 所以 ⓑ 的尺**不必新造机制**：`internal/config` 同包内新增一枚测试，
把 `collect` 的输入扩到其余段（或直接对 `Config{}` 顶层走一遍），
再新增一枚**段→档位登记表**（与 `lockedKeyDisposition` 同形状、同包、同文件族，落点建议 `internal/config/unwired.go` 旁边新开一枚 `tierregistry.go`，
⛔ 不写进 `docs/PLAN.md`——票 AC#2 末句与禁区第一条都禁这一处）。

### §3.2 尺的判据（票面 AC#2 ⓑ 逐字：「新增一枚键若其段的档位与实现不符 ⇒ 响」）

尺要**同时**读两件事，二者不符即红：
- **登记侧**：段→`config.Tier` 的表（这枚表就是"具名降级口径"的落点，ⓐ 与 ⓑ 在这一点上其实共用同一枚新数据，差别只在 ⓐ 还要求生产码去读它、ⓑ 只要求测试读它）。
- **实现侧**：真调 `(*Manager).plan(fresh)`（`manager.go:243`，**同包可直调未导出函数，这是 ⓑ 造得出的关键**）——
  对某一枚段只造"该段变、其余不变"的 `fresh`，读回 `rep.Hot/Reload/Restart` 三枚切片，看它**实际落进哪一档**。
- ⇒ **实现侧的档位是从码上量出来的，不是从注释上抄来的**，这正是票面"与实现不符"四个字的可读形态。

### §3.3 正控怎么造（必红）

1. 在登记表里把某一枚段（例如 `ball`）标成 `TierRestart`，而 `manager.go:277` 的 `rest` 表仍把它列在 hot 档 ⇒ 尺跑 `plan()`，
   `rep.Restart` 不含 `"ball"`、`rep.Hot` 含 `"ball"` ⇒ 与登记不符 ⇒ **红**。
2. ⛔ 这一枚正控**不碰任何产码、不碰 `_test.go` 已有断言**：它改的是"登记表的值"，而票 AC#2 ⓑ 要的正是"新增一枚键"这一族；
   更贴票面的正控形＝**新加一枚叶子键到某段，其段的档位与实现不符 ⇒ 指名那一步红**。
   造法＝在 `AppSection`（`internal/config/schema.go:137` 起，字段 `:138-149`）加一枚测试用字段，档位表按段的现行登记答 `restart`，
   而 `planApp` 的 `restartChanged`（`:345-347`）枚举字段、认不得新键 ⇒ 实现侧落不进任何档 ⇒ **红**。
   ⇒ 这一枚段的可抄锚：`schema.go:136` 注释逐字 `// AppSection is [app]. Everything except theme is restart-tier (D36).`
   ——**档位口径今天活在注释里**（这一句是 ⓑ 的登记表现成文案，也是 ⓑ 便宜的直接证据）。
   ⇒ **注意这枚正控需要真加一枚字段到 `schema.go`，而测试不能改产码结构**：
   可行形＝尺以**已存在的键**逐枚跑，再用一张"预期档"表比对；
   真正"新键必响"的正控只能落在**编译期存在新键**的那一侧（＝`255-r1` 的 AC#3 形尺，与本腿 AC#2 ⓑ 的尺同包同机制）。
   ⇒ **具名说：这一格 ⓑ 单独造是造不全的，它与 AC#3 共用一枚枚举器。**（见 §3.5）

### §3.4 反控怎么造（登记齐了 ⇒ 不响）

1. 登记表把 18 枚段逐枚填齐，且每枚的档与 `plan()` 实测一致 ⇒ 尺**全绿**（这是"补登记后不响"的形）。
2. 反控的第二枚必备面：`app`／`voice` 是**段内双档**（§1.3）。登记表若写成"一段一档"，
   则 `planApp` 改 theme 会落 hot、`plan` 改 language 会落 restart，**两次实测对一档登记 ⇒ 尺必红**，而码是诚实的。
   ⇒ 所以 ⓑ 的登记表**结构上必须允许一枚段挂多档**（形如 `段 → {键前缀 → 档}`，或对 `app`/`voice` 允许两行）。
   ⇒ **"粒度只到段"这句降级口径，今天在码上并不完全真**——`planApp`/`planVoice` 已经是键级分档，只是不可枚举。
   ⚠ 这一条是本腿对票面措辞的一处**读数更正**（写进 §6）。

### §3.5 造不造得出（具名）

- **造得出**：`internal/config` 同包内、反射枚举＋真调 `plan()`＋登记表比对。⛔ 不需新增依赖边、⛔ 不碰 `cmd/wisp`、⛔ 不碰 `internal/panel`、⛔ 不碰 `docs/PLAN.md`。
  要动的文件＝**1 枚新文件（登记表）＋1 枚新测试文件**（或把登记并入 `unwired.go` 旁边）。
- **造不出的一部分（具名，不含糊）**：
  ① **「新增一枚键 ⇒ 必红」这条正控在纯测试侧造不出**（新键要求改 `schema.go` 的结构声明，改结构声明属产码改动，⛔ 本腿不碰、
  也不该由一枚 AC#2 尺去碰）。可达形＝尺对**当前全部键**跑，从而保证"任何一枚新键一出现就会因为不在登记表里而红"——
  **这是同一枚机制的必然推论，但它的正控证据要由 AC#3 那条腿在真加键时给出**。⇒ 对编排者的实际意思：**AC#2 ⓑ 与 AC#3 是同一把尺的两半，分开裁会重复造。**
  ② **T3（`panel.EffectiveTier`）那一套 ⓑ 管不到**：它的值不来自 `plan()`，来自 `panel_config_store.go:228/275` 的硬接常量，
  尺在 `internal/config` 里读不到 `cmd/wisp` 的写入腿（同包测试不可达未导出 `configStore`）。
  ⇒ ⛔ **ⓑ 形不会修好 §1.2(c) 那笔"两套词表对同一次写入答相反"的账**。要修它只能走 ⓐ 的完整版（第 4 枚文件）。
  ③ **T2 `restartTierKeys` 与 `planApp` 的自动比对（即 `config_reload.go:297` 注释承诺的 "and by a test"）ⓑ 落点在 `cmd/wisp` 包**，
  不在 `internal/config` ⇒ 要兑现它得**再加一枚 `cmd/wisp` 侧测试**，与票面"AC#2 待我裁之后再排"的排程无冲突，
  但**会撞正在跑 `cmd/wisp` 的那条腿的写面**（票 255 排程节第一句）。

### §3.6 ⓑ 撞哪些钉

- ⛔ **零枚现成钉被 ⓑ 打红**（本腿量到的范围）：ⓑ 不改 T4 载荷形状（N1–N4 全活）、不改 `reportReload` 文案（N5–N7 全活）、
  不收敛 T3 四值（N8 全活）。
- **反而借力**：N3（`TestManagerThemeIsHotInsideApp`）与 N7（restart 三枚 needle）已经是 ⓑ 想要的"实现侧"读数的一部分，
  ⓑ 的登记表可以直接以它们的断言体为来源抄出 ⇒ ⓑ 与现有钉**同向**。
- **一处真冲突**：ⓑ 若要顺手兑现 `config_reload.go:297` 那句 "and by a test"（§3.5-② 之外的第三格），
  就要在 `cmd/wisp` 新增测试 ⇒ 与票面 AC#1 那一格同写面（AC#1 正是要改 `config_reload.go` 那句话），
  **两枚腿会碰同一枚文件**，需按票面排程串行。

---

## §4 我可能判错的条目（先写满这一节，逐条给"如果我错了会怎样"）

**K1 「几套词表」这个数是我用词法圈出来的，圈法本身会漏。**
我认「生效级别词表」的操作判据＝枚 `Tier`／`Effective`／`restartTierKeys` 字样的**符号**，加 `Report` 的三枚字段名
（`Hot`／`Reload`／`Restart`）。尺子＝`grep -rn ... --include=*.go cmd internal tools scripts`。
⇒ **漏法一**：一套只以裸串字面量活着、通篇不写 "tier" 这个词的档位口径，我的尺看不见。
⇒ 漏法二：**非 Go 侧**（`scripts/*.sh`、`*.ps1`、`tools/**` 之外的配置件、任何 JSON/TOML 常量）我只在 `scripts tools` 两根下扫过 `.go`，
`--include=*.go` 这一条就把它们全滤掉了。**⇒ 如果我错了**：§1 表里的套数会**少报**，ⓐ 形的"要不要顺手统一第四套"这一栏会跟着算轻。

**K2 `config.Tier` 的写者枚数我是按语法点位算的，可能把"类型名"当成了"有人在用"。**
现量：`grep -rn "TierHot\|TierReload\|TierRestart\|type Tier " --include=*.go cmd internal tools scripts`
⇒ 除 `internal/risk/assessor.go:147` 那枚同名的安全档（§1 表第 5 行）外，**七枚命中全在 `internal/config/schema.go:30-43` 自己那一段**（1 枚类型、3 枚常量定义、3 行紧贴在常量上方的注释）。
⇒ 我的判断是"零写者零读者"。**如果我错了**：唯一的可能是有人经**未点名的方式**产出这三枚值（例如把 `"hot"` 直接赋给 `Tier` 型变量、
或经反射/映射表按字符串查）。我扫了 `.Tier\b`／`Tier:`／`[]Tier`／`config\.Tier`／`EffectiveTier(`，未命中；
**但没有扫裸字面量 `"hot"`／`"reload"`／`"restart"`**（`"restart"` 与 `panel` 侧那枚 `"restart"` 会同名，串了会假报），
⇒ 这一条**未验**，写进 §5 的 U3，不当结论用。

**K3 「面板那套四值只有两值真到得了」这句我只验到"枚举常量的非测试写者"这一层。**
现量：`grep -rn "EffectiveNow\|EffectiveNextTask" --include=*.go cmd internal tools scripts | grep -v _test`
⇒ 命中 4 行，全在 `internal/panel/config_handlers.go:199,200,430,432`＝**定义 2 ＋ `tierSentence` 的分派 2**，生产写者 0。
⇒ 生产实际产出的只有 `EffectiveNotApplied`（9 枚点位）与 `EffectiveRestart`（2 枚：`cmd/wisp/panel_config_store.go:228,275`）。
**如果我错了**：`SettingWriteResult.Tier` 可能被外部（含界面侧那一半，我不读不引）以字符串形状再消费，
那"两值不可达"只到 Go 侧为止——这句**票 255 禁区第 4 条已经自己划了界**（"读侧没有出口"这一句只到 Go 侧为止），我照它划，不越界。

**K4 ⓐ 形的"要动哪几枚文件"我可能列宽了，也可能列窄了。**
我按「谁现在在做档位判断」倒推读者该落在哪（`manager.go` 的 `rest` 表 / `planApp` / `planVoice` / `planLocked`）。
⇒ 列宽的风险：票 AC#2 ⓐ 的原话只要求「谁在哪一档读它，具名到函数」，**未要求把 18 枚段全部改造**；
⇒ 列窄的风险：我**没有逐枚读 `rest` 表 12 枚段的 `set` 闭包背后有没有第二处应用点**（`plan.set` 的消费在 `CheckAndReload`，我只读到 `:198-202` 的两句 hook 分派，没读到 `plan.set` 被 run 的那个循环体）。

**K5 我把 `internal/risk/assessor.go` 的 `type Tier int` 判为"不同事、不计入"，这一步是我自己切的。**
那枚是路径敏感档（`TierNone`/`TierB`/`TierA`，注释指 `SPEC-06 §4.1`），与"生效级别"语义无关，且**它有真读者**
（`internal/tools/paths.go:312-319`）。另两枚同名者同理不计：`projctx` 的 `Tier string`（`"project"`/`"global"` 来源标签，连命名类型都不是）、
`PluginsSection.Tier2Enabled`（那是插件信任层 D19，不是生效档）。
⇒ **如果我错了**：ⓐ 形要新增一处"读 `config.Tier`"的代码，而 `Tier` 这枚标识符在仓里已同时指**四件不同的事**，
那么"具名到函数"这句会不会被裁成"先得给这枚名字消歧"——**这是命名治理，不在票面射程里，我不裁，只上报**（见 §5 的 U5）。

**K6 「段全覆盖、无遗漏」这句我只数到段名，没验语义。**
现量：`Config` 顶层 `:110-133` 共 18 枚段（17 枚带 `toml` 名 ＋ `Plugins` 带 `toml:"-"` 但手工序列化），
`plan()` 侧的归属＝`planLocked` 4 枚 ＋ `planApp` 1 ＋ `planVoice` 1 ＋ `rest` 表 12 枚＝**18**。
⇒ 我据此说"每一枚段都恰有一个档位之家"。**如果我错了**：枚数对上不等于一一对上（`app` 同时出现在 `planApp` 的 hot 与 restart 两支，
`voice` 同时出现在 reload 与 hot 两支），真要说"一一对上"得逐名比集合，我只比了基数。

**K7 工作树在我取证期间被改过，我的行号有半个时间窗的不确定性。**
`git status --porcelain internal/config internal/panel cmd/wisp` 10:43:11 读 `0`、10:44:35 读 `1`。
⇒ 我对 `internal/panel/config_handlers.go` 一律改用 HEAD blob 定行号（§0 已具名两把尺），对 `cmd/wisp`、`internal/config` 用工作树
（这两根每一发尺都是 0 行脏）。**如果我错了**：那一枚脏件在下一发尺之前又被还原，我的"HEAD 为准"就多此一举——不伤读数，只伤说明。

**K8 我把票 255 排程节里那句 `248-v1` 当成了仍在跑的腿，因此一票未跑 `go test`。**
票面写「此刻 `248-v1` 正在整包跑 `cmd/wisp`」，派单指令另写「此刻有一枚对抗验收腿（`248-v1b`）正在整包跑」。
⇒ 我按**更严的那一条**（不跑）执行。如果其实已经跑完，§2／§3 里"造不造得出发不作编译验证的言"这一保留仍成立——
我只是把不确定的事往下压，没有往上报。

## §5 判不动／没测到的地方（具名上报，不按自己的判断填）

**U1 ⛔ 未跑任何编译/测试尺**（派单硬约束：`248-v1b` 正在整包跑 `cmd/wisp`）。
⇒ 具体买不到三样东西：① 依赖边的机器可读事实（"ⓐ 形会不会新开一条包边"我只能由 `import` 与包名前缀现读，不能由 `go list -deps` 钉）；
② §3 那把尺**是否真能编出来**（正控/反控我只能给构造法，不能给跑红的证据）；
③ `reflect` 路径枚举在我这边的实际枚数（票 255 现量第 4 条给的 150/115/35 我**一枚未复认**，见 U4）。

**U2 `frontend/**`／`design/**` 两层禁令：不读不引。**
⇒ 判不动的正是「界面侧有没有第五套档位词、以及 `EffectiveTier` 的 wire 形状在界面侧被怎么消费」这一半。
票 255 禁区第 4 条已把这一半划给 owner 带话，本腿**不越界**，但 §1 表的"跨套映射"一栏因此只到 Go 侧为止，
这一栏的"没有"要读成「**Go 侧没有**」，不是「全仓没有」。

**U3 裸字面量档名未扫**（`"hot"` / `"reload"` / `"restart"`）。
⇒ 故意不扫：`"restart"` 在 `panel.EffectiveRestart` 与 `config.TierRestart` 两处同值，
按字面量算读者会把两套算成一套，**假报"有映射"**比漏报更坏。⇒ 要这条读数请另派一腿按点名校验（`grep -rn 'TierRestart\|"restart"' internal/config` 之类），
本腿宁可交"未验"。

**U4 票面现量第 4 条的三枚数（叶子键 150＝静态 115＋动态 35；被任一名册登记 11/100；段外 0/89）本腿未复认。**
⇒ 复它要跑反射枚举＝要编译＝撞 U1。本腿只复认了其中一把小的：
`grep -rn "restartTierKeys"` ⇒ `cmd/wisp/config_reload.go` 内 4 枚命中＝1 枚定义（`:299`）＋2 枚同函数内读（`:287,:293`）＋1 行注释（`:296`），**包外零读者**。
`unwiredKeys` 的表体我逐枚数了＝**6 枚**（`risk.*` 4 ＋ `net.*` 2），与票面"`unwiredKeys` 实测 6 枚"复认得上。

**U5 命名治理这一格我判不动，且票面没覆盖，按"未定义即停"上报**：
`Tier` 这枚标识符今天在仓里同时指四件事（生效档 / 敏感路径 A-B 档 / 项目说明来源档 / 插件信任层 Tier2）。
ⓐ 形若落地，新增的"读 `config.Tier`"的函数名要不要带消歧前缀，是**契约面命名**问题；
票 255 AC#2 ⓐ 的原话只说「具名到函数」，**没说名字本身长什么样** ⇒ 我不猜，交编排者裁。

**U6 `plan.set` 的消费体未读**（`CheckAndReload` 里执行 `reloadPlan` 的那一段）。
⇒ 影响 §2：ⓐ 形若要"让档位真被读"，改动点可能不在 `plan()` 而在**执行 `plan.set` 的那一步**（那里才是"值换了但没人做事"与"值换了且有人做事"的分界，也正是 AC#1 要的那句话的料）。
本腿读到 `internal/config/manager.go:198-202` 就停了，再往下是 255-r1（AC#1）的写面，我伸手进去只会撞它的读数。

**U7 AC#4（真窗口宽度）与 AC#1（那句"已立即生效"的文案判据）本腿零取证。**
⇒ 派单只要 AC#2 那一格的代价读数，我按"只有这一个问题，别扩"执行。
`cmd/wisp/panel_host_windows.go` 的 `Width: 420,` 我**一枚未读**（不是读不到，是没读，避免扩面）。

**U8 `docs/PLAN.md` 的 D36「三档生效级别」原文本腿未逐字复认。**
⇒ 我只**读了行内锚**：`internal/config/schema.go:30` 的注释自称 "D36 three-tier semantics"、
`cmd/wisp/config_reload.go:296-298` 自称指回 "PLAN.md :2715's [app] split"、`AGENTS.md` D36 行写「`config.toml` 全量 section 树 + 三档生效级别」。
⇒ 为什么这条要单列：AC#2 ⓑ 的**正身是"具名登记降级口径"**，而"登记"落在哪一字节才算不算改 PLAN.md，
是编排者的裁定；本腿在 §3 只给"码侧/仪器侧"的落点，**不给文档侧落点**。

## §6 我推翻编排者哪一句

## §7 结论
