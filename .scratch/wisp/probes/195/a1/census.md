# 195-a1 普查：票 195 的"按节写入地基"到底还欠多少

> 本腿＝只读普查，零产码、零测试、零构建、零 exe。写点只有本文件。
> 射程声明：`frontend/**`／`design/**` 本腿**零读取、零引用、零写入**；三枚冻结件里只读了
> `internal/panel/l2_grant_boundary_test.go` 的**内部判射程**（见 §3／§4），一字未改。
> 未完成标记自查：本文件不含任何未完成标记；自查尺与读数在交件回执里给，**不在本文逐字写那枚串**
> （本仓 09-25 刚踩过"尺写进结论文件、第一发读到的是尺自己"）。

---

## §0 起手锚

| 项 | 读数 |
|---|---|
| 时刻 | `2026-10-02 09:12 +08` |
| HEAD（起手） | `8b0507749e6f14efd6618021f058639c7bd1589c` / 短 `8b050774`，分支 `dev` |
| `git status --porcelain internal/config cmd/wisp` | **1 行**：`M cmd/wisp/panel_config_248_test.go`（＝别人在飞的腿，未还原未提交未删） |
| 票 195 文件 | `.scratch/wisp/issues/195-config-set-needs-a-per-section-write-foundation-before-the-panel-can-touch-it-and-the-only-writer-in-the-repo-is-permmode.md`，**41 行**（存在性已核） |
| 248 产码 commit | `0d87a681`（10-02 09:03:15 +08），stat＝15 文件／2731 增／47 删，含 `internal/config/settings.go` 314 行、`internal/panel/config_handlers.go` 505 行、`cmd/wisp/panel_config_store.go` 291 行、`internal/panel/l2_grant_boundary_test.go` **179 增**（冻结件，`A487` 具名解冻两枚锚） |

### §0.1 票面 AC#0 那四条尺的复跑（**四条里两条已不符**）

| 票面那一句 | 09-28 锚 `ba65b29f` 申报的读数 | 本腿 10-02 09:0x 实测 | 符不符 |
|---|---|---|---|
| `internal/config/` 里的导出写手 | 「只有 `SetPermissionMode` 一枚」 | **8 行**：`allowdirs.go:109`／`permmode.go:70`＋`settings.go` 的 6 枚（`:60 :81 :103 :132 :137 :177`） | **不符**（差 7 枚） |
| "按节写入"地基 | 「不存在（无任何 `Set(section, key, value)` 形状）」 | 字面尺 `grep -rnE "func .*(SetConfig\|WriteSection\|SetSection)" --include=*.go internal/ \| grep -v _test.go` ⇒ **0 行**；但那枚 0 只否掉**那三个函数名**，见 §1 | 尺**空心**（详见 §7 X2） |
| `SPEC-08:168` | 「安全节放宽走 L2 重新确认（SPEC-03 §4.2）」 | 该行逐字命中，且同一行**左边**就是 `config.get` / `config.set` | 符 |
| `bridge.go:42-45` 恰好 4 枚、无 `config.*` | 4 枚 | `grep -nE '=\s*"panel\.' internal/panel/bridge.go` ⇒ **4 行**（`:42 :43 :44 :45` 名与序一字未动）；`config.get`／`config.set` 在 `:66`／`:67`，**无 `panel.` 前缀** | 符（且 AC#6 今天未被违反） |

⚠ 编排者给的所有行号与"6 枚／8 行"我都当**待验断言**跑过：`^func (m \*Manager) Set` 那一发**实测 8 行、6 枚在 `settings.go`、行号与题面逐枚一致**；`permmode.go` 那枚 setter 今天确实在 **`:70`**（票面写 `:64` ⇒ 已漂 6 行，票面 AC#4 的尺文本里那句 `permmode.go:64` 同样过期）。

---

## §1 `settings.go` 那一套到底是什么形状

### 1.1 生死线：**具名字段枚举，不是任意 key 透传**

`writeOneKey(ownedKey string, apply func(base *Config) bool, check func(base *Config) error)`
（`internal/config/settings.go:199`）的三个入参里，**`ownedKey` 只是名字，不是寻址句柄**：

- 它**只**被用在三处：错误／拒绝文案（`:226`、`:232`、`:245`）、`mergeWrite` 的第一入参（`:242`，
  而 `mergeWrite` 自己的注释 `writeguard.go:85-87` 逐字写着"ownedKey 是调用方所持的点分键路径……
  **mergeWrite 分不清一枚键和另一枚键，也不试图分**"）、以及无前置文件时回执的那一枚兜底答案（`:259`、`:266`）。
- **真正决定改哪个字段的是 `apply` 那一枚闭包**，它手写结构体字段赋值（`:63-69`、`:84-90`、`:111-122`、
  `:148-167`、`:179-181`），**没有任何一处从 `ownedKey` 反查字段**：`settings.go` 的 `reflect` 引用数＝**0**
  （尺＝`grep -c reflect internal/config/settings.go`）；全包无 `map[string]setter`、无 TOML 透传；
  `writeguard.go` 里那批 `reflect` 全部落在**报告**那一族函数体内（`diffKeyPaths:201`→`diffStruct:216`→
  `configFieldName:241`→`diffValue:260`→`diffMap:287`），不在写地址那条路上。
  `writeOneKey` 在 `settings.go` 之外的引用数＝**0**。
- 结论：**把 `ownedKey` 换成 `"risk.permission_mode"` 不会写进那一枚键**——`apply` 只认结构体路径，
  传什么都改不动 `Risk`。"任意 key 透传"那一形在这一层**结构上不成立**。

### 1.2 可写键路径的**产生点名册**（逐枚，全仓）

`ownedKey` 只有四种产生方式，逐枚 `file:line`：

| # | 键路径形状 | 产生点 | 拼法 |
|---|---|---|---|
| 1 | `llm.providers.<p>.base_url` | `settings.go:62` → helper `:47-49` | `keyPathProviderField(provider, "base_url")` |
| 2 | `llm.providers.<p>.api_key_ref` | `settings.go:83` → `:47-49` | `keyPathProviderField(provider, "api_key_ref")` |
| 3 | `llm.providers.<p>.models.<m>.context_window` | `settings.go:109` → helper `:51-53` | `keyPathProviderModelField(provider, model, "context_window")` |
| 4 | `llm.providers.<p>.models.<m>.price.in` | `settings.go:147`（由 `:133` 传 leaf `"in"`） | `keyPathProviderModelField(..., "price."+leaf)` |
| 5 | `llm.providers.<p>.models.<m>.price.out` | `settings.go:147`（由 `:138` 传 leaf `"out"`） | 同上 |
| 6 | `llm.roles.chat.model` | `settings.go:179` | **字面量**（不走 helper） |
| 7 | `fs.allowed_dirs` | 常量 `writeguard.go:70`，用在 `allowdirs.go:142` | `keyFSAllowedDirs` |
| 8 | `risk.permission_mode` | 常量 `writeguard.go:71`，用在 `permmode.go:79` | `keyRiskPermissionMode` |

⇒ 今天**程序可写的键＝上面 8 条形状**，其中 1–6 在 `settings.go`（面板可达），7–8 是票 201／90 那两枚既有门。
回执里另有两枚**非键**串（`writeguard.go:207` 的 `(nil-side)`、`:266` 的"写后无法复读文件"那条），
它们不是可写面，别数进名册。

### 1.3 写盘走哪一族 / 逐枚受保护写在不在 / 回执诚不诚实

- **走 `mergeWrite`，不走 `SaveFile`**：`settings.go:242` 是 `settings.go` 里唯一落盘调用点，
  那一发进 `writeguard.go:111`。尺＝`grep -n "SaveFile(" internal/config/settings.go` ⇒ **0 命中**
  （那枚文件里出现的 `SaveFile` 只有头部注释 `:5`／`:6` 两行在讲它**不是**合并器；真正调 `SaveFile` 的唯一一处是
  `writeguard.go:160`，在 `mergeWrite` 内部）。票 226 那个洞没被重开。
- **"逐枚键各一发受保护写"在**，且是票 248 AC#11 第二支的形：`writeOneKey` 一次只落一枚 `ownedKey`
  （`settings.go:25-28` 注释与实现一致）；五枚 `Price` 叶子**没有**被并成一枚键——
  `SetModelPriceIn`／`SetModelPriceOut` 是两枚方法（`:132`／`:137`），各自一次 `:147`。
- **回执报的是"实际落盘的键路径"**：`writeOneKey:248` → `writtenKeyPaths:254` → `diffKeyPaths(pre, post)`
  （`writeguard.go:201`），即**写前的文件快照 vs 写后的文件**逐叶比，比的是盘上字节而不是进程信念。
  三形都不静默：前置文件读不回来 ⇒ 拒写（`settings.go:205-211`）；文件已是该值 ⇒ `mergeWrite` 一字不写
  （`writeguard.go:151-155`）；外来手改被保住 ⇒ 那枚键出现在回执里而不被认领成本进程的写
  （`writeguard.go:166-186`，`statOwnWrite` 只在 `diverged` 为空时才认领）。
  ⚠ **一枚收窄**：`preExists=false`（首跑无文件）那一支回执退化成"本 setter 所持的那一枚键"
  （`:259`），不是盘上差分——它保住了"没 diff 对象就老实说"的方向，但那一发的回执**不再是由字节证明的**。
- **写前校验在**：`validate(candidate)` 落在 `settings.go:230`，位置＝`mergeWrite` 之前、
  且校验的 base 是"这一发真要产出的那份"（`:212-218` 先 `readCurrentFile` 再 `deepCopyConfig`）。
  `grep -c validate internal/config/writeguard.go` 今天仍是 **0** ⇒ 校验不在受保护原语里，是 `settings.go` 自己补的，
  这条缺口（票 248 AC#7）被这一层闭上了。

### 1.4 一句话形状

**它是一枚"按具名叶子键的受保护写"原语 + 6 枚各自只懂一枚键的导出 setter，
不是一枚"按 section＋key 寻址"的通用面。** 名字里带 `section`／`SetConfig` 的东西依然一枚也没有，
但票面据以宣称"地基不存在"的那枚尺，量的是**函数名的字面拼写**，不是能力。

---

## §2 195 的每一格 AC 现在能不能只靠 `settings.go` 勾

（本轮回填）

---

## §3 仍欠的那一块到底是不是"安全节 L2 那一支"

（本轮回填）

---

## §4 撞钉预检（名册，不是颜色）

（本轮回填）

---

## §5 落点建议（只给形状，不给码）

（本轮回填）

---

## §6 判不动的地方

| 号 | 判不动的那一句 | 为什么判不动 | 谁能判 |
|---|---|---|---|
| **J1** | `internal/config` 与 `internal/panel` 今天**绿不绿**（票 195 AC#5 那三发） | 本腿硬边界零 `go test`／零 `go build`，而此刻 `33-r8`（`internal/ball`）与 `250-r1`（`scripts/`）两枚写码腿在跑毫秒级计时，负载会洗掉它们的读数 | 落地腿自量：`go test -count=1 ./internal/config/`；`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 ./internal/panel/ ./cmd/wisp/` ⇒ **两条必全绿**。⚠ `cmd/wisp` 那族**必须带那串 PATH 前缀**，否则 `exit status 0xc0000135` **且没有 `--- FAIL`**＝用例根本没跑（这不是绿，是没量） |
| **J2** | `sh scripts/d22scan.sh` 与 `"$(go env GOPATH)/bin/gofumpt" -l` 的终态 | 同上：`scripts/d22scan.sh` 第 15 行是 `go run . -root "$root"`（＝构建＋执行），`gofumpt -l` 要那枚 exe；本腿两者都不许碰 | 落地腿自量，逐名照抄终态。另 `bash .scratch/wisp/probes/154/gate-clauses.sh` 的判据＝**BAD 腿名册只 `G6neg`**（在册常红那枚） |
| **J3** | **票 195 AC#2 的"进队列那一发"今天到底通不通**（造一发"请求放宽安全节某项" ⇒ ①键值未变 ②队列多一条待决 ③审计有归因） | 这一格不是"读代码能读出颜色"的形：②要真起一张 C18 卡并读到队列深度，①要真比磁盘字节。而**读代码能确定的只有"今天没有任何一枚入向调用会走到这一形"**（面板口在 `config_handlers.go:315` 就把锁定族整族拒了，压根没到写入腿）——那是一句静态事实，不是一发读数 | 落地腿种一发实测：从 `cmd/wisp` 的 `runWidening`（`approval_always.go:112-147`）那条形出发，把入向换成"设置页要放宽 [fs]"，读 `queue.Depth()` 前后差与 `[audit]` 行 |
| **J4** | `A487` 那枚具名解冻**是否覆盖票 195 的落地** | 解冻记录文本我读到了（`docs/reports/pending-and-issues.md:10072` 起，`:10444` 是编排者 10-02 09:0x 的逐行核过），但它的**边界逐字写的是"只属票 248 射程"**。195 若不改那枚文件就无关；若要新增入向名而惊动那三张词表，**是另一次批准**，不是本腿能推定的 | 编排者裁。（本腿 §4 的结论是"195 不需要碰它"，所以 J4 只在"195 被裁成要新增入向名"那一支才醒） |
| **J5** | `internal/panel/composer_test.go:502` 与 `bridge_test.go:131` 那两枚钉今天读什么颜色 | 它们**读 `frontend/src`**（前者扫渲染器、后者逐枚比 `panel.*` 字面量与 `bridge.postMessage` 枚数）。`frontend/**` 是本腿的零读区，我不能开它，所以"页侧会不会被 195 打红"我判不了 | 界面侧 owner 自己那位 agent／编排者；Go 侧落地腿只需知道：票 195 AC#6 禁新增 `panel.*`，只要守住就牵不到这两枚 |
| **J6** | **195 剩余那一块该落在 `internal/config` 还是 `cmd/wisp`**（同一形两处的现成先例：`manager.go` 的 reload 侧 vs `approval_always.go` 的 caller 侧） | 这是一次**架构取舍**，不是一条能读出的事实。两侧各有已冻结的注释把"另一侧"钉成设计（`permmode.go:47-58`／`allowdirs.go:16-24` 逐字写"程序化那一条由**它的调用方**确认"），把确认搬进入口＝**改那两处的设计陈述** | 编排者裁，且必须摊开"改完哪几枚既有测试的注释与断言同时过期"（§4 N6/N7 两枚就是要摊的那两枚） |

---

## §7 我推翻编排者题面之处

| 号 | 题面那一句 | 实测 | 推翻到什么程度 |
|---|---|---|---|
| **X1** | 「`internal/config/` 里今天唯一现成的写手是 `Manager.SetPermissionMode`」 | 那把尺今天返回 **8 行**（§0.1 第一行） | **整句已不成立**（`0d87a681` 之后）。它不再是"唯一"，`SetPermissionMode` 也连"最新"都不是 |
| **X2** | 「按节写入地基不存在（**无任何 `Set(section, key, value)` 形状**）」 | 字面尺复跑＝**0 行**，看着"符"；但那枚 0 是被**三个写死的函数名**筛出来的（`SetConfig`／`WriteSection`／`SetSection`）。`writeOneKey`＋`mergeWrite`＋`diffKeyPaths`＋写前 `validate`＋逐枚键受保护写＋实际落盘回执＝**票 248 AC#11 那一支的地基四件套今天全在**，只是不叫那三个名字 | **半推翻**：按票面**字面**判，那句仍"对"；按**能力**判已经不成立。这正是本仓第 29ⓑ 那一族（尺本身是空的）——**票 195 的"地基不存在"是一枚由尺的拼写造出来的事实，不是一条盘上事实** |
| **X3** | 「票 195 与票 248 都明令 ⛔ 任意 key 透传」＋ 题面问"它是具名枚举还是任意透传" | `settings.go` 是**具名枚举**，且比题面假设更强：`ownedKey` **连参与写入的资格都没有**（§1.1）。透传那一形在这层结构上不成立，不是"目前没被用" | 支持题面倾向的那一支，且给到更强的一句 |
| **X4** | 题面给的行号 `SetProviderBaseURL:60`／`SetProviderAPIKeyRef:81`／`SetModelContextWindow:103`／`SetModelPriceIn:132`／`SetModelPriceOut:137`／`SetRoleChatModel:177`／`writeOneKey:199`／`writtenKeyPaths:254`／`readCurrentFile:274`／`allowdirs.go:109` | **全部逐枚对上**（本腿独立跑的 `^func (m \*Manager) Set` 与三次直读） | 未推翻，确认为可引用凭据 |
| **X5** | 题面：「另有非导出的 `writeOneKey`／`writtenKeyPaths`／`readCurrentFile`」 | 三枚都在，另**漏了两枚**非导出成员：`setModelPriceLeaf`（`settings.go:141`，两枚价格 setter 的共同体）与 `requireCatalogEntry`（`:294`）／`unknownProviderErr`（`:310`）。`internal/config` 里第 4 枚非导出写手 helper 是 `allowdirs.go:139 writeAllowedDirs` | 小推翻：非导出面 6 枚不是 3 枚 |
| **X6** | 题面把 `l2_grant_boundary_test.go` 当"三枚冻结件之一、一字不许改" | 它今天**已被 `0d87a681` 改过 179 行**（`git show --stat`），依据＝台账 `A487` 具名解冻 `:2051`／`:2139` 两枚锚，编排者自己 10-02 09:0x 在 `:10444` 逐行核过在射程内。两枚锚已从"数一枚 4"换成**能力形**（`guardRosterOf:2030`＋`:2207`＋`:2308` 的 `reflect.DeepEqual(got, answeredRoutes)`） | **对现状的推翻**：那两枚锚点上的钉**今天不再是"一字不许改"，而是"只许按能力形改"**；锚外其余断言仍一字不许动。**这条对本腿无影响**（我没改它），但**对下一枚落地腿有影响**——它现在钉的是"守卫 case 表 ⇔ `Method*` 常量名册双向全等"，195 若新增/删除任何 `Method*` 常量而不同步守卫，必红 |
| **X7** | 题面：「195 剩下的那格是不是'把放宽那一支挡住'」 | **不是"挡住"，也不是"建出来"，是第三种**：那一支的**行为**在盘上已经有了，只是**不在 `internal/config`**——`cmd/wisp/config_reload.go:114` 把 `Manager.ConfirmLocked` 接到 `confirmLockedLoosening`，后者在 `:246` 走 `Gate.PendingApproval`（C18 队列、只接受原生侧 allow）；`planLocked`（`manager.go:304`）今天就把"放宽"排进待确认、**未答完不写进 `cur`**，"收紧/中性"直接热生效（`:314` 那一支 vs `:321` 那一支）。**缺的是"按节寻址的入向入口"＋"程序化写侧也吃这套方向裁决"**（今天程序化三枚锁定族写入 `SetPermissionMode`／`AddAllowedDir`／`SetAllowedDirs` 把 L2 义务**外包给自己的调用方**，见 `permmode.go:47-58`／`allowdirs.go:16-24`） | 这是本票最要紧的一句，写进 §3／§5 |
| **X8** | 题面：「`internal/agent/approval/ui.go` 的 `PanelAPI` **没有 `Allow` 方法**（注释逐字『it is a method that does not exist』）」 | 逐字命中：`ui.go:166` 是那句注释，`PanelAPI` 在 `:167-171`，成员只有 `Reject`／`Head`／`View` | 未推翻，确认为现成仪器 |
