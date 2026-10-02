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

## §1 「生效级别」词表的现状表（每套一行，四问）

本节在 §4／§5 落盘后的下一发 commit 填写。

---

## §2 ⓐ 形代价：让 `Tier` 真被读者消费

本节在 §4／§5 落盘后的下一发 commit 填写。

---

## §3 ⓑ 形代价：具名登记「粒度只到段」+ 一把会响的尺

本节在 §4／§5 落盘后的下一发 commit 填写。

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
