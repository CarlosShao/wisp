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

判"已满足"的姿势＝**问生产调用者，不只数定义**；并区分"函数存在"与"面板口能走通"。
⚠ 一处全局降档：`248-r1` 自己的**门禁四数无凭据**（台账编排者裁定，其证据件 §4 是空的），
所以本表凡引 `settings.go`／`config_handlers.go` 的"能用"，成立层级＝**静态 grep 可证**，不等于今天绿。

| 格 | 三态 | 凭据／差哪一句 |
|---|---|---|
| **AC#0 前提现量对得上** | **不成立（该响亮报回）** | §0.1 四条里两条不符（导出写手 1→8、"地基不存在"那把尺是空心的）。票面自己写的处置＝"任何一条不符 ⇒ 停下报回，不许按票面硬改"——本腿就是那次报回。**这一格永远勾不了**，它的作用是把票头换掉 |
| **AC#1 分节写入面存在且被生产调用者用** | **部分满足** | **满足的那半枚（且不是"只有测试听众"）**：真实现＝`internal/config/settings.go:199`（`writeOneKey`）＋六枚导出 setter；非测试调用方＝`cmd/wisp/panel_config_store.go:184 :194 :201 :208 :215 :217 :267`（**7 处**）；再往上两枚装配根生产点＝`cmd/wisp/run.go:416`（`rt.settings = newConfigStore(...)`，被 `:695` 的快照 reader 用）与 `cmd/wisp/panel_inbound.go:266-279`（`configWrites` 建在 `:266`、`Config:` 槽挂在 `:279`）；**面板口那一跳也真能走通**：`panel_inbound.go:228 newComposerDispatchChain` 被常驻腿经 `panel_resident_windows.go:170`（由 `:189` 那一发走到）复用，`resident_windows.go:142` 建它，宿主 `panel_host_windows.go:551` 逐字 `m.disp.Handle(ctx, raw)` ⇒ 链路＝`panel.host → ComposerDispatch.Handle → composer_dispatch.go:197/202 → ConfigWriteHandler → configStore.ApplySetting → Manager.Set* → writeOneKey → mergeWrite`。<br>**不满足的那半枚**：入口**不接 section／key 参数**（§1.1），可写面只有 §1.2 那 8 条具名叶子键。票面要的"按节写入"那一形——一枚能用 section 与 key 寻址的入口——**今天不存在** |
| **AC#2 安全节"放宽"不落盘、进队列** | **完全没做**（`internal/config` 那一层），但**反向那发已满足、且三断言的形在 `cmd/wisp` 已有可复用先例** | 尺＝`grep -rn "PendingApproval\|approval\." --include=*.go internal/config/ \| grep -v _test.go` ⇒ **0 命中**＝这枚包今天与审批队列**零耦合**，没有任何入口能把一枚放宽转成待决项。<br>① **反向那发（收紧直接写成功）今天成立并被钉住**：`writeguard_226_test.go:308 TestAC5SetAllowedDirsPersistsARevocationWithoutAClobber` 要求删除条目**必须落盘**（红句逐字 "the revoked entry is still in the file"），`allowdirs.go:95-102` 注释逐字"收紧这一节不需要 L2 重新确认"。⇒ 票面 AC#2 末句担心的"把 fail-closed 方向一起禁掉"那一支，**现成的防呆已经在**。<br>② **三断言的形在 cmd 层已经成立**：`cmd/wisp/approval_always.go:112-147 runWidening`＝第二张 L2 卡在前、`AddAllowedDir` 在后，卡未答允 ⇒ 磁盘 `allowed_dirs` 不动、队列多一枚、`record("WIDEN-REFUSED"/"WIDEN-APPLIED", ...)` 带 `fromCorr`/`tool`/`section=fs key=allowed_dirs` 归因；钉＝`cmd/wisp/approval_always_201_test.go:121 TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card`／`:163 TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused`。⇒ **这一格欠的不是行为，是"把它搬进 `internal/config` 的一枚按节入口里"** |
| **AC#3 面板来源不许冒充原生** | **部分满足，且今天的满足是"空成立"** | 现成仪器三层都在：`internal/agent/approval/ui.go:167-171` 的 `PanelAPI` 只有 `Reject`/`Head`/`View`、**没有 `Allow`**（`:166` 注释逐字 "it is a method that does not exist"）；`internal/panel/config_handlers.go:108 lockedFieldFamilies`＝`risk.`/`fs.`/`net.`/`plugins.`/`privacy.`/`models.`/`audio.` 七族，在 `:315` 于到达写入腿**之前**整族拒；结构性后果＝**今天压根没有一枚"面板来源的安全节放宽"路径存在**，所以"它被挡住了"目前没有可挡之物。<br>**差的那一句**：AC#3 真正要判的是"195 建出那枚按节入口**之后**，来源＝面板的放宽仍然结构性进不来"。今天这一格从"空成立"变成"要真挡"，只差那一枚入口 |
| **AC#4 `SetPermissionMode` 没被冒充分节写** | **已满足**（起手名册＝终态名册那一格，今天对 195 尚无新增） | 尺（票面原文，本腿实跑）＝`grep -rn "SetPermissionMode" --include=*.go internal/ cmd/ \| grep -v _test.go` ⇒ **8 行**，逐枚：`internal/config/allowdirs.go:18`（注释）／`internal/config/permmode.go:47`（注释）／`permmode.go:70`（定义）／`internal/config/unwired.go:127`（**判决文本**）／`internal/perm/store.go:49`（接口注释）／`store.go:50`（接口方法）／`store.go:208`（**唯一生产调用点**）／`cmd/wisp/run.go:391`（注释）。<br>**新地基与它的关系＝0**：尺＝`grep -c SetPermissionMode internal/config/settings.go` ⇒ **0** ⇒ `settings.go` 没走那枚捷径。⚠ 一枚要点名的副作用：名册里 `unwired.go:127` 那句"persisted by `Manager.SetPermissionMode`; ticket 90"是**描述性**命中——落地后若持久化归属变了，那枚字符串必须同批改口，否则会出现"名册枚数不变、句意变假"这一形（本仓第 29ⓑ 族） |
| **AC#5 门禁四数＋`gofumpt`** | **本腿不量** | 见 J1／J2。**但有三枚 d22scan 的牙现在就能点名，落地腿必须绕**（这些是**形状**判断，不是"今天红"的判断）：① `tools/d22scan/main.go:707` 红句逐字 "bare `go func(` is banned (D22/D38b): use `observe.Registry.Spawn`"——AC#2 若把"等一张卡"写成新 goroutine 必踩；② ban#2＝`filepath.Clean`/`Abs` 在 `risk.PathResolver` 之外做文件系统决策（`main.go:15`）——按节入口**不许**自己拼配置文件路径，`writeOneKey` 走的是 `m.path`（`settings.go:275`、`:283`）那枚既有路径，照它就不踩；③ `emojiRe`（`main.go:164`）扫**字符串不豁免**——新回执文案里 `✓`(U+2713)／`≤`(U+2264) 会红，`→` 抓不到但规格禁 |
| **AC#6 本票没顺手接面板口** | **对票 195 自己：已满足；对票面的次序假设：已被 248 推翻** | 尺实跑＝`grep -nE '=\s*"panel\.' internal/panel/bridge.go` ⇒ **4 行**，`:42 :43 :44 :45` 名字与顺序一字未动；`internal/panel/` 里也**没有**新增 `panel.*` 常量。<br>但票面末句"名册那两枚要等本票地基绿了再按票 194 堆1 单独入账（每枚一枚 `A##`）"**已经被执行掉了、且不是按 194 那条形执行的**：`config.get`／`config.set` 现在就在 `bridge.go:66`／`:67`（无 `panel.` 前缀），入账凭据是台账 **`A487`（10-01 00:03 那枚具名解冻＋九问裁完）**，不是"票 194 堆1 一枚一枚 `A##`"。⇒ 下一枚落地腿读 AC#6 时**别把它当"别看板"**：面板口已接，且接的是 248 那枚**只覆盖非锁定族**的地基 |

---

## §3 仍欠的那一块到底是不是"安全节 L2 那一支"

**具名答：不完全是。"把放宽挡住"这一支今天已经挡住了（两层，且都是现成仪器）；欠的是"把方向裁决搬进入口层"，
因为按节入口一开，会被它顺手变可写的东西，今天只靠"没人写它"挡着。**

### 3.1 "挡住"那一支的现成仪器（逐枚，都在，不必新建）

| 仪器 | 位置 | 它挡的是哪一形 |
|---|---|---|
| `PanelAPI` **没有 `Allow` 方法** | `internal/agent/approval/ui.go:167-171`；作者注释 `:166` 逐字："PanelAPI has no Allow method: 'panel-sourced allow is structurally rejected' is then not a check that could be forgotten, **it is a method that does not exist**" | 面板侧宿主**压根没有可调的授权出口**（对照 `NativeAPI:143-161` 才有 `Allow`/`AllowSession`） |
| 面板口的**整族拒** | `internal/panel/config_handlers.go:108`（七族名册）→ `:115 refusedLockedFamily` → 在 `:315` 于**进写入腿之前**拒 | `risk.*`／`fs.*`／`net.*`／`plugins.*`／`privacy.*`／`models.*`／`audio.*` 从 `config.set` 一条都进不来；`:129-133` 还兜住不带点分的裸名 `permission_mode`／`allowed_dirs` |
| 程序化写的**义务外包** | `permmode.go:47-58`／`allowdirs.go:16-24` 注释逐字：程序化那一条"由**它的调用方**确认"；`cmd/wisp/approval_always.go:112-147` 就是那个调用方 | 放宽必须先有第二张原生 L2 卡 |
| reload 侧的**方向裁决** | `manager.go:304 planLocked`：`len(loosen)>0` ⇒ 排进 `plan.confirm`（`:314`）**未答完不写 `cur`**；`len(tighten)>0` ⇒ 直接热生效（`:321`）。方向裁决器＝`riskDirection:421`／`fsDirection:457`／`netDirection:470`／`pluginsDirection:490`／`setDirection:521`；生产接线＝`cmd/wisp/config_reload.go:114 rt.mgr.ConfirmLocked = rt.confirmLockedLoosening`，`confirmLockedLoosening` 定义在 `:221`、真调队列在 `:246 rt.gate.PendingApproval(...)` | 手改放宽 ⇒ 进 C18 队列、只吃原生 allow |
| 冻结仪器（**射程判断，非内容引用**） | `internal/panel/l2_grant_boundary_test.go`：词表 `grantFieldWords:183`（18 枚）／`grantRouteWords:192`（11 枚）／前缀 `:1316`（8 枚，**全是 `panel*` 族**）／后缀 `:1324`（3 枚）／长度钉 `:1334 :1335 :1336`；总闸 `TestRealGuardRefusesEveryAssemblableApprovalRouteName:1530`（528 枚拼装名问**真守卫**）；封套面 `TestNoInboundEnvelopeCanBindAnApprovalVerdict:1547`；守卫名册能力形 `guardRosterOf:2030`（用在 `:2207`、`:2308`） | "被真守卫应答过的名字里不许有授权形的拼写"。**射程边界（该文件自己写的）**：`grantRoutePrefixes` 之外＝不在网里，所以 `config.get`／`config.set` **今天不被那 528 格问起**；同理任何新造的 `config.*`／非 `panel*` 命名都不被问 |
| 规格与口径 | `SPEC-08:168` 逐字"安全节放宽走 L2 重新确认（SPEC-03 §4.2）"（本腿 `sed -n '166,170p'` 复跑命中）；`SPEC-06:19` L2 那一行逐字"「允许」只接受原生侧来源……面板只能「拒绝/查看完整参数」"；`AGENTS.md` §1.2 禁令"由面板侧来源的 L2「允许」" | 判"是不是越界"的三处文字凭据 |

### 3.2 那一块真正的缺口（★ 本普查最要紧的一句）

`internal/config` 里**按节寻址的写入口**一旦开出来，下面这两枚键就从"没人能写"变成"能写"，
而它们**都是 `manager.go` 已经认过的放宽键**，且**都不在 `unwiredKeys` 的六枚名册里**
（名册尺＝`grep -c "path:" internal/config/unwired.go` ⇒ 6）：

| 键 | 方向裁决 | 有没有 unwired 守卫 | 消费者 | 今天为什么还安全 |
|---|---|---|---|---|
| `fs.delete_enabled` | `manager.go:460-461` 逐字：`!old.DeleteEnabled && new.DeleteEnabled ⇒ loosen` | **无**（默认 `false`，`schema.go:481`；不在 `unwiredKeys`） | `unwired.go:131` 逐字 "consumed: cmd/wisp/run.go gates the delete-capable tools" | **只因为一枚写入者都不存在** |
| `fs.reparse_point_exceptions` | `manager.go:459` 走 `setDirection` ⇒ 增项即 loosen | **无** | `unwired.go:130` "consumed: cmd/wisp/run.go feeds the C26 canonicalizer" | 同上 |

对照：`risk.shell_enabled`／`risk.allow_shell_string`／`risk.shell_allowlist`／`risk.blacklist_overrides`／
`net.allowlist`／`net.block_private_ranges` 这六枚**非默认值会被 `validateUnwired`（`unwired.go:100-110`）
直接判错**，而 `settings.go:230` 的写前 `validate` 会撞上它 ⇒ 这一族即便有了按节入口也**写不进去**（fail-closed，好）。
⇒ **所以缺口不是"放宽没挡住"，是"两枚有真消费者的锁定族放宽键，其唯一防线是'没有写入者'，
而 195 要建的正是那枚写入者"。** 这句话就是 §5 那把会响的尺要钉的东西。

### 3.3 那三张词表对新造入向名的射程（具名比过）

195 若要新增**入向名**（方法名或封套键），逐枚比：
- 非 `panel*` 命名（例如再走 `config.*` 那一族）⇒ **不在** `grantRoutePrefixes` 那 8 枚里 ⇒ 不被 528 格问到；
  但仍受 `guardRosterOf` 那条"守卫 case 表 ⇔ `Method*` 常量名册双向全等"约束（新增常量不改守卫＝红）。
- 任何带那 11 枚路由词根之一（`approve`/`approval`/`grant`/`allow`/`permit`/`ratify`/`authorize`/`authorised`/`decide`/`decision`/`verdict`）
  的路由名 ⇒ 若命名落在 `panel*` 前缀内即被网格问到、真守卫应答就红；落在 `panel*` 之外则**今天问不到**，
  但会撞 `internal/panel/config_route_248_test.go:259 TestAC1FieldTableCarriesNoApprovalVocabulary`
  里**同一份 11 枚词表的手抄版**（它对 `MethodConfigGet/Set` 与 `WritableFields()` 两批名字做子串筛查）。
- 任何带那 18 枚字段词根之一的 JSON 键、挂在 `ComposerRequest` 可达面上 ⇒ 撞 `l2_grant_boundary_test.go:1547`。
- ⛔ 三条都**不许**为了躲词表去动那三枚表（它在冻结件里，`A487` 只解冻 `:2051`/`:2139` 两枚锚，
  且那两枚锚已被换成能力形）。**唯一安全形＝新名字里不含任何裁决词根，且不进 `panel*` 前缀族。**

---

## §4 撞钉预检（名册，不是颜色）

⛔ 本腿不跑测试 ⇒ **每一枚的"今天绿不绿"一律＝待落地腿自量**；命令与期望读数在 §4 末。
名册＝**会挡住 195 落地**的现存钉（行为型优先，光 grep 新符号名不算预检）。

| # | 用例名 | `file:line` | 断言原文（要点） | 存在理由 | 对 195 的后果 |
|---|---|---|---|---|---|
| **N1** | `TestAC1LockedFamilyFieldsAreRefusedBeforeTheLeg` | `internal/panel/config_route_248_test.go:215`（红句 `:232`／`:239`） | 五枚锁定族字段（`risk.permission_mode`／`fs.allowed_dirs`／`permission_mode`／`net.proxy.mode`／`privacy.keep_transcript`）必须 `errors.Is(err, ErrConfigFieldRefused)`、回执必须**点名那把门**（`mustSay` 分别是 `panel.mode.request`／`L2`／`锁定族`），且 `leg.applyCalls != 0` 即红 | 票 248 AC#1／J5：整族拒，"⛔ 任意 key 透传"的机器可读形 | ★ **本票最硬的一枚**：195 的 AC#2 要的是"`fs.*` 放宽 ⇒ 转队列、不落盘"，而这枚钉要的是"`fs.*` ⇒ 压根到不了写入腿"。**两枚判据互斥**。落地腿不许靠"改 `mustSay` 的子串"绕开——那正是"为变绿放宽断言" |
| **N2** | `TestAC1FieldTableCarriesNoApprovalVocabulary` | `internal/panel/config_route_248_test.go:259`（等值钉 `:275`，写死的名册在 `:271-274`） | `WritableFields()` 必须 `reflect.DeepEqual` 那**七枚**（`:271-274` 写死），且两枚路由名与七枚字段名都不许含 11 枚授权词根 | 名册＝白名单；扩表是契约变更 | 195 若往面板可写名册加**任何一枚**（含"安全节"字段）⇒ 必红。**195 的按节入口必须留在 `internal/config`，不许挂到这张表上**（这也正是 AC#6 的意思） |
| **N3** | `TestNoInboundEnvelopeCanBindAnApprovalVerdict` | `internal/panel/l2_grant_boundary_test.go:1547`（词表 `:183`；谓词 `carriesGrantWord:242`、`grantCarryingKeys:324`） | `ComposerRequest` 及其可达类型绑出的每一枚 JSON 键都不许含 18 枚裁决词根；同一条规则再跑一遍 `ModeRequest`／`AttachmentPayload`；末了还有一发 AST 面扫整个 `internal/panel` | D33/F2、R20、`AGENTS.md` §1.2 禁令#6 | 195 若要新增入向字段（`configSection` 之类可以，`confirmWiden`/`allowLoosen` 之类必红）⇒ **逐枚比那 18 枚词根**，且⛔不动表 |
| **N4** | `TestRealGuardRefusesEveryAssemblableApprovalRouteName` | 同上 `:1530`，网格 `sweepAssembledNames:1433`，长度钉 `:1334/:1335/:1336`（8×11×3×2＝528） | 528 枚**运行时拼装**的授权形路由名，真守卫一枚都不应答；三张表任一因子被缩 ⇒ `t.Fatalf` | 抓"名字从不完整写在任何一处"那一形 | 射程只盖 `panel*` 八枚前缀。195 若新增 `panel.<动词>` 形 ⇒ 有被命中风险 ⇒ **改判成不用 `panel.` 前缀**（248 已经这么做了，`bridge.go:66/67`），⛔ 不去扩那张表 |
| **N5** | `TestPlantedGrantWiringGoesRedInASnapshot` 的守卫名册两枚锚 | 同上 `:2168`，锚在 `:2207`（`rosterProblem != ""`）与 `:2308`（`reflect.DeepEqual(got, answeredRoutes)`），派生器 `guardRosterOf:2030` | 守卫 case 表与本包声明的 `Method*` 常量**双向全等**，且必须写成一整行逗号分隔；正控＝植一枚"守卫答了无名册常量的路由"必须被派生器拒绝 | `A487` 把"数一枚 4"换成能力形 | 195 一旦新增／删除任何 `Method*` 常量而不同步守卫（或反之）⇒ 红。⚠ 这也意味着**那两枚锚今天不再是"一字不许改"**（见 X6） |
| **N6** | `TestEveryLockedSectionKeyIsAccountedFor` | `internal/config/unwired_test.go:311`（表在 `internal/config/unwired.go:119-147`） | 反射走 `[risk]`/`[fs]`/`[net]`/`[plugins]` 每一枚键，缺 `lockedKeyDisposition` 判决即 `t.Errorf`；`unwired:` 前缀者必须与 `unwiredKeys` **双向 1:1**（多一枚守卫没人指、少一枚判决都红） | 票 83 AC#1：完整性可执行化 | 195 若新增任何锁定族键（含子表）⇒ 红；若为了"给放宽留个位"往 `unwiredKeys` 加行而不同批加判决 ⇒ 也红。**它是 §3.2 那两枚键的现成登记面**——把 `fs.delete_enabled` 的处置改成"放宽走队列"要在同一批动这张表 |
| **N7** | `TestAC5SetAllowedDirsPersistsARevocationWithoutAClobber` | `internal/config/writeguard_226_test.go:308`（红句 `:318`／`:321`／`:324`） | 收紧那一发**必须真落盘**、内存与文件不得分叉、且不得回滚掉别人的手改 | 票 226 AC#5 落库半枚 | **这是 AC#2"反向那发"的防呆**：195 若把"安全节写入"一律改成不落盘，此钉必红。票面 AC#2 末句"别把批准和次生收紧绑一体"＝指的就是它 |
| **N8** | `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card`／`...StoresNothingWhenTheSecondCardIsRefused` | `cmd/wisp/approval_always_201_test.go:121`／`:163` | 「一直」那支**第二张卡答允之后**才存；拒绝 ⇒ 什么都不存 | 票 201：放宽的 L2 义务在调用方 | 195 把确认搬进入口层 ⇒ **调用方那一层会重复确认或空确认**。这两枚就是"改完之后哪些既有测试的语义会过期"的那一批（J6 要摊开的两枚之一） |
| **N9** | `TestConfirmLockedRunsOutsideTheManagerLock`／`TestCheckAndReloadAsksOncePerFileChange` | `internal/config/manager_223_test.go:55`／`:152` | 钩子必须在 `mu` **之外**跑（可在里面读自己的配置、可阻塞等人）；一次文件变更只问一次 | 票 223：一枚未答的 C18 卡曾冻结所有 `Config()` 读者 300s | ★ 195 的入口若"在 `writeOneKey` 里等一张卡"⇒ **必须同守这条锁序**（`writeOneKey:202-203` 起手就 `m.mu.Lock()`）。**这一枚会死锁而不是死红**，是本轮唯一一枚"红了看不出为什么"形状的钉 |
| **N10** | `TestWhitelistMethods...`（票 181 反漂移那枚） | `internal/panel/git_test.go:385`（比对 `:381-383` 写死的四枚常量）、尺 `whitelistMethodsFromSource:407`、正则 `panelMethodRe:394`、枚数钉 `:517`（`len(real) != 4`） | 从 `bridge.go` **源码**抽 `"panel.*"` 双引号字面量，必须＝四枚；正控＝植第五枚 `panel.*` 必须被数到 | 票 181：显示半枚不许被接回动作半枚 | 195 若在 `bridge.go` 里新增任何 `panel.` 前缀字面量（**连注释里带双引号的那一枚也算**，正则不剥注释）⇒ 红，且与 N4/N5 同时叫。守 AC#6 就不踩 |
| **N11** | `TestAC11SettingWriteReportsTheKeyPathItChanged` 一族 | `internal/config/settings_248_test.go:57 :79 :113 :166 :192 :213 :238` | 逐枚键各一发受保护写、非法值⇒文件**逐字节不变**、外来手改保住、写前 `validate` 生效 | 票 248 AC#7／AC#11 | **不是阻挡，是要接的既有契约**：195 的按节入口若复用 `writeOneKey` 就自动继承；若新造第二形（例如绕开 `pre` 快照）⇒ 这族里的字节比对会红 |
| **N12** | `TestBoundaryNoRuntimeObservationFields`／`TestBoundaryNoPlaintextSecretFields` | `internal/config/boundary_test.go:59`（禁词八枚）／`:38` | 反射扫 `Config{}` 全部字段名，含 `health`/`probe`/`latency`/`lasterror`/`last_error`/`successrate`/`success_rate`/`verified` 即红 | 存储切分：运行期观测态归 SQLite | 195 若为"待决项"往 `Config` 加字段（例如 `PendingWiden`）⇒ 名字里带 `last*`／`verified` 之类就踩。**入口返回类型不是 `Config`，别加进去** |

### §4 末：命令与期望读数（本腿不跑，落地腿逐名照抄终态）

```
# 1) 包内（票 195 AC#5 的字面判据）
go test -count=1 ./internal/config/                                    # 期望：全绿，无 --- FAIL
# 2) 面板＋装配根：必须带 PATH 前缀，否则 exit status 0xc0000135 且【没有 --- FAIL】＝用例根本没跑
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 ./internal/panel/ ./cmd/wisp/
# 3) 门禁
sh scripts/d22scan.sh                                                  # 期望：clean（注意 ban#8 那格报的是【被扫文件数】，不是违规数）
bash .scratch/wisp/probes/154/gate-clauses.sh                          # 期望：BAD 腿名册只 G6neg（在册常红）
"$(go env GOPATH)/bin/gofumpt" -l internal/config cmd/wisp internal/panel   # 期望：空
```
⚠ 第 2 发那串 PATH 前缀**不是可选项**：`cmd/wisp` 那族缺它会给 `exit status 0xc0000135` ＋ `0.0xxs`、
**且不打印任何 `--- FAIL`**——本仓 09-28 已有人把这判成"写手报错"。
⚠ 另：`cmd/wisp` 起手即在的一枚 boot Ctrl+C 时序 flake 已知 ⇒ 报红绿时要给**整包序＋隔离跑两形**才作数。

---

## §5 落点建议（只给形状，不给码）

### 我选：**乙′（差一块地基，但那一块不是我原以为的那一块）**

不是甲（195 没被 248 整枚满足），不是丙（面板口那一跳 248 已经接完了）。
**准确形状＝248 顺手建掉了"逐枚具名叶子键的受保护写"这一半地基，
195 欠的是"按节寻址 + 方向分流"那一半**：一枚能用 `(section, key, value)` 表态的入口，
它对**非锁定族**直接走 `writeOneKey` 那一族落盘，对**锁定族**按 `manager.go:421/457/470/490` 现成的
方向裁决分流——**收紧 ⇒ 直接写；放宽 ⇒ 交进 C18 队列、不落盘**（`ui.go:167-171` 保证这枚队列的 allow
只有原生侧给得了）。它落在哪几枚文件，具名：

| 落点 | 为什么是它 |
|---|---|
| `internal/config/` 新增**一枚**入口（不新增文件也行，与 `settings.go` 同包） | 三枚方向裁决器与 `validate`/`mergeWrite`/`diffKeyPaths` 全在同包，**零新导出面、零新依赖边**就能复用；`writeOneKey:199` 已经是那枚"共同体"，缺的只是把 `ownedKey` 从"闭包顺手写死的名字"升成"入参" |
| `internal/config/manager.go`：把 `riskDirection`/`fsDirection`/`netDirection`/`pluginsDirection` 的调用面从"整节比较"具化成"单键判定" | 今天这四枚只被 `planLocked` 用（reload 侧）。按节入口必须吃**同一个**裁决器，⛔ 不许在入口处另写一套"什么是放宽"（那是第二真相源，本仓最贵的那一形） |
| `internal/config/unwired.go:118-147` 的 `lockedKeyDisposition` | §3.2 那两枚键（`fs.delete_enabled`／`fs.reparse_point_exceptions`）的处置行必须与入口**同批**改，否则 N6 红、且注释变假话 |
| ⛔ **不动**：`bridge.go:42-45` 四枚、`config_handlers.go:90` 那张 7 枚字段表、`l2_grant_boundary_test.go` 三张词表、`writeguard.go` 的 `mergeWrite` 本体 | AC#6＋N1／N2／N4／N5／票 226 地界。**尤其：195 不许以"给面板口加一枚安全节字段"的形式落地**——那会同时打红 N1+N2 并把 AGENTS 禁令#6 变成一扇开着的门 |

### 三档各自的代价与"选错会在盘上留下什么形状"

**甲＝"195 整枚已被 248 满足，只立残余"**
- 选它的代价：**把 AC#2 那两枚有真消费者的锁定族放宽键（`fs.delete_enabled`／`fs.reparse_point_exceptions`）
  继续留在"只靠没有写入者挡着"这一形上**，并把票 195 那一格在追踪矩阵里划成已交付。
- 选错会在盘上留下：一枚**看起来已被保护、其实只是没人碰**的放宽面。等任何后续票（票 114 那族"档位/工作区"、
  票 128/132、或界面侧那句"我想自己配"）第一次调"什么键都能写"的入口，它就**静默生效**，
  且新增的那枚 setter 会带着一份"这是 195 那枚地基"的信心。⚠ 本仓定式：**错误的信心与正确的信心在盘上长得一模一样**
  ——所以甲的判据不能是"现在没有写入者"（会随时间翻），必须是下面那把尺。
- **会响的尺**（甲这一档唯一诚实的形，落成常驻用例）：
  具名一枚 `internal/config/` 内的枚举——**锁定族每一枚 loosen 键都必须显式声明它的处置**
  （`unwired` / `queue` / `direct-tighten-only`），并把 §3.2 那两枚键**点名**进去；
  尺的形＝`fsDirection`/`riskDirection` 的**输出集**逐个查处置表，缺项即红（形状照 N6，别照词表）。

**乙＝"差一块地基"（我选的档）**
- 选它的代价：195 落地要**同时**动 N6（判决表）、N8（cmd 侧那两枚"调用方确认"钉的语义）、
  N9（锁序，⚠ 这枚是**死锁**不是死红）与 §3.2 那两枚键 ⇒ 它是**一枚中型腿，不是一枚补差集的小腿**，
  而且它和 248 已落的那半枚地基**共用 `writeOneKey`**，两枚验收腿的射程会重叠（`248-v1` 与 195 的验收者要提前分工）。
- 选错（把范围做大了）会在盘上留下：把 C18 队列搬进 `internal/config` ⇒ 那枚包今天对审批层**零耦合**
  （尺＝§2 AC#2 那一发 0 命中）。搬进去＝新一条依赖边（`internal/config → internal/agent/approval`），
  而本仓对"装配根之外新增依赖边"是要单独 `A##` 的。⇒ **判据形状要反过来**：入口只**返回**"这一发是放宽、
  需交队列"这一枚**判定**，谁去举卡仍归装配根（`cmd/wisp`）——与 `ui.go` 那两枚 API 的分工同形。

**丙＝"地基已在、只差面板口那一跳"（归 194／145／114 之一）**
- 选它的代价：**这一档已经不成立**，因为它的前提"那一跳没接"今天为假：`bridge.go:66/67`、
  `composer_dispatch.go:197/202`、`panel_inbound.go:266-279`、`panel_resident_windows.go:170`＋`:189`、
  `panel_host_windows.go:551` 五枚点连成的那条线**已经在盘上**（静态可证；绿不绿＝J1）。
- 选错会在盘上留下：最坏的一种——**再补一跳**。`Config` 槽已经有人接了，第二枚入口意味着两个真相源
  （两个"谁被允许写 `config.toml`"），而那正是票 248 头部注释和 `run.go:412-415` 逐字在防的那一形。
- ⛔ **不许选丙**，除非有人能指出上面五枚点里哪一枚今天不存在。**指得出＝丙；指不出＝读漏了。**


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
