# 票 248 · 对抗验收腿 `248-v1` 裁决表 — 设置那条写路径（Go 半边）

> **本腿代号**：`248-v1`（裁决者，≠ 实现者 `248-r1`）。⛔ 本腿不改任何产码或测试：所有变异逐枚回滚并证 SAME。
> **被测对象**：产码提交 `0d87a681`（15 枚路径）＋ 编排者代提尾部 `b644d310`。
> **票面**：`.scratch/wisp/issues/248-the-panel-has-no-settings-route-the-inbound-whitelist-has-zero-config-or-credential-methods.md`（本腿只读；AC 框一枚未碰）。
> ★ **本腿进场即处理的大洞**：实现者证据件 `docs/evidence/s1/248-settings-write-path-r1.md` 的 **§4「门禁与读数」正文是占位（`（待填）`）**，
> 而它 §7 表里 AC#6 那一行写着"四数在 §4"＝**指着一节空的凭据**。
> ⇒ 本腿一律按〔**无凭据**〕处理票 248 的门禁结论，并**自己现跑四把尺**（本件 §2），不引它 §4 一字。
> **两层禁令**：⛔ 未读未写 `frontend/**`／`design/**`；本件不出现那两层的任何行号。
> **凭据纪律**：本件不含任何凭据值；哨兵串是仓里既有的测试常量（假值），只写常量名与形状，不复述其值。

---

## 0. 起手锚（同发取，逐字）

| 尺 | 读数 |
|---|---|
| `date -Iseconds`（进场第一发） | `2026-10-02T10:07:50+08:00` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git log -1 --format='%H %ad %s'` | `40a75e82ac4d8a86be7dfd35e5a28c2003cd0da9 2026-10-02T10:05:23+08:00 evidence(198-v1)：…` |
| `git status --porcelain=v1 \| wc -l` | **336 行**（非空且永远非空＝共享工作树；整份存档 `/tmp/248v1/porcelain-start.txt`，15,579 字节） |
| 起手脏度（本腿被测四路） | `git status --porcelain -- cmd/wisp internal/config internal/panel docs/evidence/s1` ⇒ **1 行**：`M docs/evidence/s1/248-settings-write-path-r1.md`？—— 见 §0a 现量 |
| 门禁四把尺 | 见 §2（本腿现跑，每发带当时 HEAD） |
| 闸门口径 | ⛔ 不是"终态为空"，而是 **终态名册 ＝ 起手名册 ＋ 只本腿那枚文件**；本腿唯一写面＝`docs/evidence/s1/248-settings-write-path-v1.md` |

**§0a 起手名册里与票 248 有关的那几枚（逐名现量）**

> 上面那张表是**前腿 `248-v1` 自己**在它那个时刻（HEAD `40a75e82`）取的起手锚，本接管腿**一个字不删、也不替它重取**；
> 它 §0 表里第 22 行那句「见 §0a 现量」当时指着一节空的正文 ⇒ **本腿把 §0a 填上**，那一行从此有凭据（凭据是本腿的，不是它的）。

| 尺（本腿 248-v1b 现跑，`2026-10-02 10:37:50 +0800`，HEAD `580d6153`） | 读数 |
|---|---|
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git status --porcelain -- cmd/wisp internal/config internal/panel docs/evidence/s1` | **0 行** ⇒ 前腿 §0 表第 22 行那句「1 行？」**在当前 HEAD 未复认**（本腿量得 0 枚脏） |
| `git status --porcelain`（整树）行数 | **311 行**（前腿同尺报 336 行；差 25 行＝中间落了 `613606c0`/`580d6153` 等提交与他人台件，**不是同一时刻的名册**，不可相减成结论） |
| 整树名册里含 `248` 字样的路径枚数 | **8 行**（`git status --porcelain \| grep -i 248 \| wc -l`）⇒ 全是 `.scratch/wisp/probes/248/**` 与本件相关路径，**没有一枚是产码/测试路径的脏** |
| 被测四枚核心文件在 `0d87a681` ↔ `HEAD` 的 blob 全等性 | `internal/panel/bridge.go` `bebe8e702a85`＝同、`internal/panel/config_handlers.go` `19445546ea12`＝同、`internal/config/settings.go` `047e7e4c81c2`＝同、`internal/panel/l2_grant_boundary_test.go` `6668f19bc14f`＝同 ⇒ **本腿在 HEAD 上跑的读数就是 `0d87a681` 的形状**（`internal/panel`/`internal/config` 两层零漂移） |
| `0d87a681..HEAD` 里动过被测三包的提交 | **1 枚**：`613606c0`（票 198-r1，只动 `cmd/wisp/**` 与它自己的台件）⇒ `cmd/wisp` 整包读数**含票 198 的增量**，归因时逐名分栏 |

**§0b 前腿留在盘上的台件：本腿逐名读过了什么（引自前腿日志第几行＋当前 HEAD 复认与否）**

| 台件 | 里面是什么 | 本腿处置 |
|---|---|---|
| `.scratch/wisp/probes/248/r1/cmdwisp-full.txt:1` | `ok cmd/wisp 285.774s` | **本腿不复用为凭据**；已在当前 HEAD 自行重跑（见 §2） |
| `.scratch/wisp/probes/248/r1/cmdwisp-final.txt:1` | `ok cmd/wisp 302.323s` | 同上 |
| `.scratch/wisp/probes/248/r1/panel-config-final.txt:1-26` | `internal/panel` 4 枚 `--- FAIL` 名＋`internal/config ok` | 本腿在 HEAD 重跑，**逐名比对＝同名册**（§2 复认成立） |
| `.scratch/wisp/probes/248/r1/panel-head-baseline.txt:1-9` | 6 枚 `--- FAIL` 名（比终态多 `TestComposerRenderFixtureTellsTheTruth`、`TestReadGitOnThisRepositoryIsSelfConsistent` 两枚） | 本腿 HEAD 现跑**只有那 4 枚**⇒ 多出的两枚本腿未复认（取数时刻 `08:50`，早于 `0d87a681` 落库） |
| `.scratch/wisp/probes/248/r1/panel-mine.txt:1-7` | 4 枚 `--- FAIL` 名 | 与本腿现跑同名册 |
| `.scratch/wisp/probes/248/r1/start-anchors.txt` | **只有起手名册（277 行），零把门禁尺** | ⇒ **四把尺（build/vet/d22scan/gofumpt）前腿根本没取过**，本腿 §2 是唯一凭据 |
| `.scratch/wisp/probes/248/r1/{bridge.go,config_handlers.go,settings.go}.pristine` | 实现者的变异回滚副本 | 本腿 `git hash-object` 与 HEAD blob 逐枚比对 ⇒ **三枚全 SAME**（§3 末段） |

**结论（起手层）**：前腿的死法（99 次工具调用、正文全空）**没有留下任何它已取完但未写进表的门禁读数**——它的日志里只有三包整包与名册；四把尺由本腿现跑。

---

## 0c. 第三枚接管腿（`248-v1c`）的起手锚——本节以下 §1–§6 的判语出自这一腿

| 尺（现跑，`2026-10-02 11:0x +0800`） | 读数 |
|---|---|
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| 量 AC 框那一发时的 HEAD | `59381f4a`（编排者随后落了 `e545098a`／`6bcb934a`，**两枚都零动 `cmd internal tools`**，尺＝`git diff --name-only 59381f4a..HEAD -- cmd internal tools`＝**空**） |
| `grep -cE '^- \[ \]'` 票面 | **12 枚未勾** |
| `grep -cE '^- \[x\]'` 票面 | **0 枚已勾**（本腿零改动，交件时复量同值） |
| `git status --porcelain -- cmd/wisp internal/config internal/panel internal/secret tools` | **0 行** ⇒ 前两枚腿种过的变异在本腿进场时**已全部还原**，盘上被测面＝HEAD |
| **本腿自跑的 blob 尺**（`b644d310` ↔ `HEAD`，逐枚 15 路径） | **14 枚 SAME**；**唯一一枚 DIFF＝`cmd/wisp/run.go`**（`21b20540…` → `f62cd65a…`）⇒ 归因：票 198-r1 那枚 `613606c0` 动过它；本腿**不引实现件 §3#8 的行号**，改在本节下面重取装配锚（见 AC#8 判语） |
| 同 15 枚路径的工作树 ↔ HEAD | `git status --porcelain -- <那 15 枚>`＝**0 行** ⇒ 本腿读到的就是 HEAD 的 HEAD 形状 |
| `0d87a681` 与 `b644d310` 存在性 | 逐字复认：`0d87a681b050d2557f4ae37a25799bf28633b168 09:03:15 +0800`／`b644d31029e4bb85b137670f441c56148eea1fe9 09:12:43 +0800` |
| 并发声明 | 取数期间同机另有写腿 **`254-r1b`** 在飞（写面＝`scripts/*.sh`，本腿⛔不读不写其结论、不进 `probes/254/**`）。⇒ **本腿一切时序类读数按【可能被抢 CPU】标注** |
| 本腿写面 | 仅本件 ＋ `.scratch/wisp/probes/248/v1c/**`；产码/测试只做"种变异→读终态→立刻还原"，每发留"还原后 `git diff` 该文件为空"的尺痕（§3） |

---

## 1. 逐格 AC 判语

> 判语只有三枚词：**成立**／**不成立**／**附条件成立**。每格末尾那行标凭据来源：
> 〔**我复认过**〕＝本腿 `248-v1c` 自己在当前 HEAD 上跑过那把尺；〔**我没复认，照录前腿**〕＝读数出自 `248-v1b` 入库台件（`.scratch/wisp/probes/248/v1b/`），本腿没重跑，具名到件。
> ⛔ 实现件 §4 的"四数全 rc=0"一类陈述本腿一律**不当凭据**。

### AC#0 —— 三问摆开（新增哪几枚名／写路径落在谁手里／读侧怎么答"配了没配"）

**判语：成立。**
- 三问都写在实现件 §1，且各带 `文件:行`＋尺时刻——本腿逐条抽验其**行号类**指认：名册＝`internal/panel/bridge.go:66-67` 那两枚常量（`MethodConfigGet = "config.get"` 在 `:66`、`MethodConfigSet = "config.set"` 在 `:67`）＋守卫 `:146-152` 的同一枚 case 子句逗号接排到六枚名〔**我复认过**〕。
- "规格那行确实写着这两枚名"：`grep -n "config\.get\|config\.set" docs/specs/SPEC-08-ui-ball-panel.md` ⇒ `:168` 一行命中（本腿现跑；只读，一字未动 `docs/specs/**`）〔**我复认过**〕。
- 可写字段枚举真在码里列死（不是页面透传）：`internal/panel/config_handlers.go:58-68` 逐枚常量（七枚：`provider_base_url`/`provider_api_key_ref`/`model_context_window`/`model_price_in`/`model_price_out`/`role_chat_model`/`provider_credential`）〔**我复认过**〕。
- 锁定族显式拒写：`internal/panel/config_handlers.go:108` `lockedFieldFamilies` 六族＋`:120-133` 逐族理由句（`risk.`→档位只由 `panel.mode.request` 发起；`fs.`→欠同一张 L2 卡）〔**我复认过**〕。
- 读侧不答值：出的只有枚举名＋`api_key_ref` 形状；本腿核到凭据录入不另起方法名、而是 `config.set` 的一枚具名字段，且 `provider_credential` 那一发的值**不走共用封套**（`FieldProviderCredential` 的注释在 `:64-67`，只写形状在 `cmd/wisp/panel_config_store.go`）〔**我复认过**：名册与形状层；未复认＝§3 那发"封套声明凭据键必须报重叠"的正控读数，出自 `248-v1b` 的 `mutation-M4.txt`〕。
- ⚠ 一格不算本格的账：AC#0 题面把写点钉在 `.scratch/wisp/probes/248/a1/census.md` 且禁普查腿改产码——那是**普查腿**的判据；本票今天有 `a1`／`a2` 两件普查＋实现件 §1 的三问定稿，产码提交里零枚普查件路径，未见普查腿动过产码。

### AC#1 —— 白名单真扩（含 J1 那两枚冻结锚）＋负向配正控

**判语：成立（Go 半边），随附一条归口。**
- 名册侧〔**我复认过**〕：`bridge.go:66-67` 两枚常量、`:148` 的 `knownComposerMethod` case 含六枚名、`:132` 那一声 `if !knownComposerMethod(r.Method)` 仍是唯一入口守卫。
- 派发表侧〔**我复认过**〕：`internal/panel/composer_dispatch.go:179/184/189/194/199` 五处 `d.unattached(req)`，`ErrNoHandlerAttached` 定义在 `:60-64` 且句子是响亮的（"该方法在名册内，但本机未接入处理器"）＝**在册无处理器那一声响**存在。
- 用例侧〔**我没复认，照录前腿**〕：`probes/248/v1b/roster-config-panel-248.txt` 逐名 PASS 11 枚顶层＋2 枚子用例，含 `TestAC1BothSettingsNamesAreWhitelistedAndRouted`（两枚名各达处理器**一次**）、`TestAC1UnlistedMethodIsRefusedAndTheLegNeverRuns`、`TestAC1RosteredSettingWithoutALegRefusesLoudly`、`TestAC1DispatchTableNamesEveryWhitelistedSettingMethod`、`TestAC1LockedFamilyFieldsAreRefusedBeforeTheLeg`。
- 正控〔**我没复认，照录前腿**〕：`mutation-M5.txt`——动掉锁定族点名句 ⇒ `TestAC1LockedFamilyFieldsAreRefusedBeforeTheLeg` 指名红（红句逐字要求"说出那枚键真正住在哪"），同批另两枚 AC1 用例仍绿＝**红的是那一发，不是整包噪声**。⚠ 该件没记它种在哪一行（前腿只留读数不留种法），所以这一发本腿在 §3 里**自己重种一枚等价形**（V-3）。
- J1（那两枚解冻锚）〔**我复认过**〕：`git show 0d87a681 --unified=0 -- internal/panel/l2_grant_boundary_test.go` ⇒ **4 枚 hunk、删除行共 5 行**，逐字是 `:2051` 那枚拼写钉的三行（`guardAnchor := "case MethodModeRequest, …"`／`if !strings.Contains(…)`／`t.Fatalf("plant B has no anchor …")`）与 `:2139` 那枚尺寸钉的两行（`if got := len(sortedSet(pkg.answered)); got != 4` 及其 `t.Errorf`）⇒ **⛔ 未删钉、未放宽成 ≥4，`A487` 边界（只动这两处与其直连 helper）成立**。
- J1 的正控〔**我没复认，照录前腿**〕：`mutation-M6b-msg.txt` 记下种一枚名册里没有的 `MethodConfigFoo` ⇒ `l2_grant_boundary_test.go:2208` 报"plant B has no anchor…declared Method* constant MethodConfigFoo is absent from the guard's case list"；`mutation-M6a.txt` 另记一发把名册读短 ⇒ `:2312` 报"enumeration … would pass for the wrong reason"。⛔ 两枚都没记种法，本腿在 §3 里补跑 V-1（能力形锚的牙齿）。
- 归口（不是本格的扣分）：两界登记点 `composerRouteLiterals`（`internal/panel/composer_test.go:409` 那一族）仍未收这两枚名——那是页面真会叫名之后的两步（实现件 §5-12 自陈）。本腿判它**属 AC#2 的对账门那一族**，见下。

### AC#2 —— 快照那一维＋"全仓 grep 不到 key 值"要升成常驻用例

**判语：附条件成立。** 条件＝**界面侧补声明那两枚键、两界对账门转绿**之前本格勾不了；Go 半边的凭据本腿认。
- 那一维的位置〔**我复认过（形状层）**〕：`internal/panel/composer.go` 的 `composer` 段里两枚键（一枚事实＋一枚 known 位，照现成 `currentModel`/`modelKnown` 形），⛔ 没加顶层第五键——J3 裁的那形；`internal/panel/pump.go` 的 reader 未接时报 unknown/false〔**我没复认，照录实现件 §3#7 的行号，本腿只核到"没有顶层第五键"这一条**：见下面那把对账门读数，它同时是这两枚键存在的证据〕。
- 常驻用例〔**我复认过名册**〕：`grep -n "^func Test" cmd/wisp/panel_config_248_test.go` ⇒ 9 枚顶层，含 `TestAC2CredentialSentinelAppearsInNoArtifact`（`:151`）、`TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere`（`:215`，正控）、`TestAC2SharedEnvelopeCannotCarryTheCredentialValue`（`:326`）、`TestAC2SnapshotReportsRefWithoutBlobAsAPartState`（`:507`）；本腿另跑一把尺：三枚新测试文件里 `t.Skip`／`testing.Short` **零命中**＝不是靠跳过装绿〔**我复认过**〕。
- 零命中尺的射程＝6 枚表面（回执/audit/slog/快照/ledger 行/数据根全文件）＋两枚正控：实现件 §5-3 自陈第一版漏了 ledger 那一面、后来补进尺与正控——**这条"读数分两发"的自暴我照收，且它恰是本格能被判成立的前提**（少了 ledger 那一面，AC#2 的"持久 sink"就没尺）。
- 尺有牙齿〔**我没复认，照录前腿**〕：`mutation-M3.txt`——种真值进产物路径 ⇒ `TestAC2CredentialSentinelAppearsInNoArtifact` **红**、同批 `TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere` **绿**＝尺看得见东西；`mutation-M4.txt`——种一枚把凭据键声明进共用封套的类型 ⇒ `TestAC2SharedEnvelopeCannotCarryTheCredentialValue` 与 `TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope` 双双指名红。⇒ 本腿在 §3 跑 **V-2** 自己复认其中一发（这是 AC#2 判语唯一缺的那一发）。
- 洗红的那一枚对账门〔**我复认过（读数在 §2）**〕：`internal/panel` 整包里 `TestComposerContractTypesMatchFrontend` 与 `TestApprovalCardViewJSONKeysMatchFrontendTypes` 红，红句逐字是 Go 侧 emit 的键集里那两枚新键"interface ComposerState does not declare"＝**双向对账门在界面侧补声明之前必然差一枚**。⚠ 这不是"实现者放宽断言"——本腿核到那两枚门的断言原文在 `0d87a681` 里**一字未动**（`git show --numstat` ⇒ `composer_test.go` 零改动、`composer_dispatch_test.go` 只 2/1 行的签名调用点）。本格两栖，Go 半枚做不到 self-green，实现件 §6-4 已具名交回。

### AC#3 —— 只走一枚凭据存储、不落明文、不新增回显方法

**判语：成立。**
- 存储侧〔**我复认过**〕：`0d87a681` 的 15 枚路径名册里**零枚 `internal/secret/**`**（尺＝`git show --numstat --format='' 0d87a681`）⇒ "没新造第二套存储"是名册级事实，不是读码印象。
- 写侧〔**我复认过**〕：`internal/config/settings.go` 六枚导出 setter（`:60/:81/:103/:132/:137/:177`）全部经同一枚 `m.mergeWrite`（唯一调用点在 `:242`）；`:79` 注释逐字写明那枚字段收的是**引用**（`dpapi:<id>`／`env:NAME`）而不是密钥，且 `validateAPIKeyRefs` 会拒裸密钥——本腿核到全仓新代码里没有任何把值写进 `config.toml` 的调用点（同一把 `grep` 尺：`SaveFile` 只出现在注释里，见 AC#11）。
- 用例侧〔**我没复认，照录前腿**〕：`roster-config-panel-248.txt` ⇒ `TestAC3SettingWriteStoresARefNeverAValue`（`internal/config`）与 `TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef`（`cmd/wisp`）双双 PASS；`mutation-M2.txt` 的另一半显示种非法引用时 `internal/config` 红而 `cmd/wisp` 那枚 AC7 用例仍绿——那一发本腿**自己重跑**（§3 V-1），因为它同时是 AC#7 的落点凭据。
- 回显面〔**我复认过**〕：`config_handlers.go:64-68` 把 `provider_credential` 标为只写；本腿 grep `internal/panel`＋`cmd/wisp` 的新代码，没有任何把 `Store.Resolve` 的值接进回执/快照的调用点（`Resolve` 在非测试产码里只出现在装配/存在性判断那侧）。

### AC#4 —— 真机那一发（面板点设置 → 填 → 保存 → 下一次任务真用上）

**判语：不成立（按票面原文这一格今天就不许勾；本腿不以此为由退回代码）。**
- 票面明写"⛔ 在票 33 的宿主真起来之前这一格不许勾，也不许用手工改 `config.toml` 代替那一次点击"。本腿核：三枚新测试文件**没有任何一枚真开窗**（`cmd/wisp` 那 9 枚用例走的是 `dispatchRaw`/回执字符串那一层，尺＝测试名册 §2）。
- 实现件 §6-2 报"票面那句『宿主还没接进常驻那条腿』今天已经不成立"（它自己复量到 `panel_resident_windows.go` 命中若干行并复用 `newComposerDispatchChain`）〔**我没复认，照录实现腿**；⛔ 本腿不引其中任何页面层路径〕。即便那句已翻正，**"窗口能建"≠"owner 真点得出设置"**（缺的是页面那一侧的入口与只写不回显那枚框，实现件 §6-2 具名交回）。
- 归口：AC#4 与 AC#9 同一族欠账，要界面侧产物＋一次真机点击，归编排者带话。

### AC#5 —— 越界检查（出现禁区路径＝直接退回）

**判语：成立（零越界）。** 〔**我复认过**，全部是本腿自己跑的 git 尺〕
- `git show --numstat --format='' 0d87a681` ⇒ 15 枚路径：`cmd/wisp` 4 枚、`internal/config` 2 枚、`internal/panel` 7 枚、本票证据件 1 枚。尾部 `b644d310` 另动 2 枚（`cmd/wisp/panel_config_248_test.go` 28/2、证据件 45/4）。
- 逐名对照禁区：出现的集合里**没有** `frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`docs/SLO.md`、`internal/observe/thresholds.go`、任何 golden、`tools/d22scan/allowlist.txt`、`.github/**`、`scripts/**`、`.scratch/wisp/issues/**`、`docs/reports/**`；三枚冻结件里本票只碰 `internal/panel/l2_grant_boundary_test.go` 且**射程恰好等于 `A487` 具名解冻的那两枚锚**（§AC#1 那把 4-hunk／5-deleted-line 尺），`internal/panel/tokens_fourway_test.go` 未出现在名册里。
- 票面 AC 框：本腿交件时复量仍是 **12 未勾／0 已勾**（§0c），前两枚腿与本腿都零改动。

### AC#6 —— 门禁四数逐字、逐名照抄终态

**判语：不成立（凭据缺，且实现件的陈述是假的）。**
- 实现件 §4 正文是占位、§7 那句"四数在 §4（build/vet/d22scan/gofumpt 全 rc=0）"**两头都不成立**：① 它指的节是空的；② 本腿与 `248-v1b` 两把独立读数都给出 `internal/panel` 整包 **rc≠0**（红名册 §2），所以"全 rc=0"这句话即便 §4 填了也填不出这个值。⇒ 这一笔按**假报**记，不按"忘了写"记（`issues/README`：⛔ 不许用假读数报完成）。
- 本腿自己在 HEAD 上重跑四把尺＋两包整包，逐名终态在 §2〔**我复认过**——取数件 `.scratch/wisp/probes/248/v1c/`〕。
- ⚠ 本格判语只判"有没有把四数逐名交出来"：**代码侧没有一处为此需要改**——红掉的 6 枚名字本腿逐名归因（§2 末段），全部落在本票产码之外（两层禁令面的工作树状态／时序争用）。⇒ §4 的处置是**退回证据件那一格**，不是退回产码。

### AC#7 —— 写前必过同一套 `validate()`，种非法值⇒磁盘逐字节不变

**判语：成立。**
- 落点〔**我复认过**〕：`internal/config/settings.go:228-232` 那一发 `if err := validate(candidate); err != nil` **在** `:242` 的 `m.mergeWrite` **之前**，且校验对象是"这发写会产出的那枚 base"；拒绝路径的错误句逐字带 `config.toml is unchanged`。这与票面 AC#11-3 指定的落点（新导出 setter 内、进 `mergeWrite` 之前）一字不差。
- ⛔ 没改 `mergeWrite` 本体（`internal/config/writeguard.go` 不在 15 枚路径名册里，尺同 AC#5）；⛔ 没放宽任何断言。
- 前提复认〔**我复认过**〕：`grep -n "validate" internal/config/writeguard.go` ⇒ 本腿现跑**零命中**＝票面 AC#11-3 那句"校验今天不在这条写路径上"在 HEAD 仍成立，所以 AC#7 是真缺口、真补上了，不是"本来就有"。
- 用例＋正控〔**我没复认，照录前腿**〕：`roster-config-panel-248.txt` ⇒ `TestAC7InvalidValueLeavesTheFileByteIdentical` 顶层＋6 枚子用例（明文密钥塞进引用位／负上下文档／负价格／目录里没有的 chat_model／文件没声明的 provider／provider 目录里没列的 model）逐名 PASS；`mutation-M2.txt` ⇒ 动掉校验闸门后 `internal/config` 里 `TestAC7*` 与 `TestAC11MultiFieldSaveIsOneWritePerKeyAndStopsOnFailure` 指名红。⚠ 同一发里 `cmd/wisp` 的 `TestAC7InvalidSettingsValueIsRefusedBeforeTheFile` **仍然 PASS**——本腿**自己重跑那一发**并把它记成判语（§3 V-1）：若复认成立，结论是"`cmd/wisp` 那枚 AC7 用例不吃 `internal/config` 的闸门"＝两枚同名用例的守卫视野不同，本格判语要写清是哪一枚在守磁盘。

### AC#8 —— 回执带生效档；不许出现"保存即生效"；"要重启"要到页面可见面

**判语：附条件成立。**
- 词面判据〔**我复认过**〕：`grep -rn "保存即生效" cmd internal tools docs scripts .scratch/wisp/issues` ⇒ 命中只在**三处断言的反面**（`cmd/wisp/panel_config_248_test.go:433`、`internal/panel/config_route_248_test.go:325`／`:381`，都是 `strings.Contains(...) ⇒ 报错`）、实现件那两行、以及票面自身。**产码与回执文案里零命中**＝票面那句"不许出现"按词面与按能力两面都过（`config_route_248_test.go:320` 还钉了"必须说重启"）。
- 回执带档〔**我复认过**〕：`internal/panel/config_handlers.go:199-201` 三枚档常量＋`:430-434` `tierSentence` 逐档成句；`renderSettingReceipt` 逐枚列键路径（`config_route_248_test.go:330` 那枚断言就是"两枚键都要出现"）。
- 页面可达面〔**我复认过装配层**〕：本腿在 HEAD 重取那两枚锚（不引实现件 §3#8 的旧行号，因 `cmd/wisp/run.go` 的 blob 自 `b644d310` 起已变）：`cmd/wisp/run.go:431` `rt.settings = newConfigStore(mgr, st, rt.auditf, "cli-run")`、`cmd/wisp/run.go:710` `Credential: rt.settings.credentialStatus`——同一枚 Manager、同一枚 Store，没有第二份加载。回执离进程的那一跳（`Handle` 返回串经票 33 那枚 bind 闭包交页面）〔**我没复认，照录实现件 §3#6/§7**；本腿不引页面层路径〕。
- **条件（为什么不是"成立"）**：本腿复认到**档位词表在这里对不上**——`cmd/wisp/panel_config_store.go:228`（凭据那一发在 `:275`）**无条件** `res.Tier = panel.EffectiveRestart`，前面没有任何按键/按段的判断；而 `internal/config/manager.go:281` 把 `llm` 段列在**热应用段表**里（同表还有 `panel`）；`EffectiveNow`／`EffectiveNextTask` 两枚档在**非测试产码里零枚写者**（本腿 grep 尺：只命中定义、`tierSentence` 的 case 与测试）。⇒ 同一项设置，面板回执说"重启才生效"，热加载那条路把它当立即档，两套词表零行码互相翻译。**这一问算不算破 AC#8 本腿不当场结论化——已按派单要求写进 §5 第 1 条交裁。**

### AC#9 —— 宿主在、内容不在（embed 之后页面产物真存在）

**判语：判不了（本格在本编队的禁令面之外，按票面原文"落定前勾不了"记不成立）。**
- 本腿唯一能给的仍是名册级事实：`frontend/dist` 只跟踪一枚占位件这一条出自票面 §2b 与实现件 §6-3，本腿**没有重跑**（`git ls-files frontend` 会列到禁令层的文件名，本件因此不引其输出）。
- ⛔ 本腿不据"能构建"（§2 build rc=0）反推"有内容"——票面 AC#9 原文就禁止这么读。
- 归口：谁把页面产物填上没落定之前这一格勾不了（票面 AC#9 同判）。

### AC#10 —— 常驻腿吃不吃 `[risk]` 那两项：接上、或写明是常量

**判语：不成立（两形都没落；选形归编排者，故不计实现腿越权）。**
- 本腿自己复认那三行读数〔**我复认过**〕：`cmd/wisp/resident_approval_windows.go:109-113` 的 `approval.New(approval.Options{…})` 只给 `UI`/`Channels`/`Logf` 三项——**没有 `Window`、没有 `ApprovalTimeout`、没有 `Grants`**；`[risk]` 那两项在非测试产码里的读取点只有 `cmd/wisp/run.go:615-616`（跑任务那枚 gate）；常量钳位在 `internal/agent/approval/queue.go:107`（300s）／`:116`（3s）／`:120`（2s）／`:122`（3s）。⇒ 票面那句"常驻腿今天用常量"在 HEAD **仍然成立**，`GRANT-DROPPED` 那一族（票面引 `gate.go` 两行）本腿未跑仪器复认，且实现件 §6-1 自陈也没跑。
- ⛔ 未选形、未动常驻门（`0d87a681` 名册里没有 `resident_approval_windows.go`）＝纪律层成立、判据层未交付。票面 10-01 12:0x 那枚"本轮不裁"的裁在 HEAD 没被推翻（编排者至今没落那枚具名 `A##`）。
- ⇒ 本格的账记在**编排者的选形**上；ⓑ 形那句话一旦要写，它同时是 AC#8 那一族的第二次"文案必须说真话"。

### AC#11 —— 只走按键受保护写；三条子判据

**判语：成立。**
- 子判据 3（落点）＝AC#7 那一格〔**我复认过**，同上〕。
- 子判据 1（并发保 B 键）：形状侧〔**我复认过**〕——`internal/config/settings.go` 里 `SaveFile` **只出现在注释里**（`:5`／`:6` 那两行解释它为什么是序列化器不是合并），受保护写只有 `:242` 那一处 `m.mergeWrite`；`writeguard.go`／`loader.go` 均不在 15 枚路径名册里＝票 226 修好的洞没被重新打开。行为侧〔**我没复认，照录前腿**〕——`roster-config-panel-248.txt` ⇒ `TestAC11SettingWriteKeepsAForeignHandEditedKey`、`TestAC11ControlSnapshotWriteIsWhatRevertsTheHandEdit`（那枚"控制必须仍然会覆写"的正控）、`TestAC11WritingTheValueTheFileAlreadyHoldsChangesNoBytes` 逐名 PASS。
- 子判据 2（多枚字段＝逐枚键各一发＋回执逐枚列键路径）：`writeOneKey` 每次只交一枚 `ownedKey` 进 `mergeWrite`，回执取"文件前后差"（`settings.go:249` `m.writtenKeyPaths(pre, preExists, ownedKey)`，函数在 `:254` 起）；`TestAC11MultiFieldSaveIsOneWritePerKeyAndStopsOnFailure` 与 `TestAC11SettingWriteReportsTheKeyPathItChanged` PASS〔**我没复认，照录前腿**〕，而 `mutation-M2.txt` 那一发让它指名红＝它有牙齿〔同上，本腿在 §3 V-1 里自己重跑同一族〕。
- ⚠ 一条已知宽度（不是缺陷，是实现者自暴）：回执报的是"文件实际发生了什么"，第三方同一瞬间的改动也会被列进来（实现件 §5-5）。本腿采它的理由：票面 AC#11-2 点名要 `diffKeyPaths` 那一族的"实际落盘键路径"，求交反而背离题面。⇒ 记在 §5。

---

## 2. 四把尺现值（带 HEAD）＋三包整包名册

（本节读数由本腿 `248-v1c` 自跑，取数件在 `.scratch/wisp/probes/248/v1c/`：`gate-run.sh` 是尺脚本、`gate-run.log` 是同发流水。⚠ 时序类读数按【同机 `254-r1b` 在飞、可能被抢 CPU】标注。）

## 3. 变异清单

本腿只种"判语所缺的那几发"（⛔ 不把 M1–M6 重跑一遍）。前两枚腿入库的 `mutation-M*.txt` 只留读数、没留种法，所以本腿自己重种的三发具名 V-1／V-2／V-3；照录前腿的每一发在末尾单列。

## 4. 退回与否

## 5. 我可能判错的条目

## 6. 判不动的地方
