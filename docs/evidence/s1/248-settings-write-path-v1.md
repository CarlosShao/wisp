# 票 248 · 对抗验收腿 `248-v1` 裁决表 — 设置那条写路径（Go 半边）

> **本腿代号**：`248-v1`（裁决者，≠ 实现者 `248-r1`）。⛔ 本腿不改任何产码或测试：所有变异逐枚回滚并证 SAME。
> **被测对象**：产码提交 `0d87a681`（15 枚路径）＋ 编排者代提尾部 `b644d310`。
> **票面**：`.scratch/wisp/issues/248-the-panel-has-no-settings-route-the-inbound-whitelist-has-zero-config-or-credential-methods.md`（本腿只读；AC 框一枚未碰）。
> ★ **本腿进场即处理的大洞**：实现者证据件 `docs/evidence/s1/248-settings-write-path-r1.md` 的 **§4「门禁与读数」正文是一句占位**，
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
- 用例侧〔**我复认过（本腿逐名跑过、带 `-v`，见 §2 末段）**〕：本腿名册与前腿 `probes/248/v1b/roster-config-panel-248.txt` 同名同终态，11 枚顶层＋2 枚子用例，含 `TestAC1BothSettingsNamesAreWhitelistedAndRouted`（两枚名各达处理器**一次**）、`TestAC1UnlistedMethodIsRefusedAndTheLegNeverRuns`、`TestAC1RosteredSettingWithoutALegRefusesLoudly`、`TestAC1DispatchTableNamesEveryWhitelistedSettingMethod`、`TestAC1LockedFamilyFieldsAreRefusedBeforeTheLeg`。
- 正控〔**我复认过（本腿 §3 V-4）**〕：本腿把 `config_handlers.go:123` 那句 `fs.` 的理由换成光秃秃的"该字段不可写" ⇒ `TestAC1LockedFamilyFieldsAreRefusedBeforeTheLeg` 指名红（红句逐字要求"说出那枚键真正住在哪"），同批另两枚 AC1 用例仍绿＝红的是那一发、不是整包噪声；还原后 `ok`。⇒ "负向判据一律问能力"这一条本腿有自己的一发。前腿 `mutation-M5.txt` 的红句与本腿这发**逐字同形**（同一枚用例、同一枚输入），本腿认它是真读数，但署名各归各。
- J1（那两枚解冻锚）〔**我复认过**〕：`git show 0d87a681 --unified=0 -- internal/panel/l2_grant_boundary_test.go` ⇒ **4 枚 hunk、删除行共 5 行**，逐字是 `:2051` 那枚拼写钉的三行（`guardAnchor := "case MethodModeRequest, …"`／`if !strings.Contains(…)`／`t.Fatalf("plant B has no anchor …")`）与 `:2139` 那枚尺寸钉的两行（`if got := len(sortedSet(pkg.answered)); got != 4` 及其 `t.Errorf`）⇒ **⛔ 未删钉、未放宽成 ≥4，`A487` 边界（只动这两处与其直连 helper）成立**。
- J1 的正控〔**我复认过（本腿 §3 V-3）**〕：本腿在 `bridge.go` 的同一枚 const 块里种一枚守卫不认的 `MethodConfigFoo = "config.foo"`（⛔ 不动 case 列表），`go test ./internal/panel -run TestPlantedGrantWiringGoesRedInASnapshot` ⇒ **红，且红句逐字点名那枚常量**（`l2_grant_boundary_test.go:2208: declared Method* constant MethodConfigFoo is absent from the guard's case list`）；还原后（`git diff` 该文件为空）同一条尺 **ok**。⇒ 那枚能力形锚真的在双向读：名册长了、守卫没跟，它就响。前腿 `mutation-M6a.txt` 另记一发反方向（把名册读短 ⇒ `:2312` 报"passing for the wrong reason"）〔**我没复认，照录前腿**〕。
- 归口（不是本格的扣分）：两界登记点 `composerRouteLiterals`（`internal/panel/composer_test.go:409` 那一族）仍未收这两枚名——那是页面真会叫名之后的两步（实现件 §5-12 自陈）。本腿判它**属 AC#2 的对账门那一族**，见下。

### AC#2 —— 快照那一维＋"全仓 grep 不到 key 值"要升成常驻用例

**判语：附条件成立。** 条件＝**界面侧补声明那两枚键、两界对账门转绿**之前本格勾不了；Go 半边的凭据本腿认。
- 那一维的位置〔**我复认过**〕：`internal/panel/composer.go:268-269` 那两枚键的 json tag 逐字是 `json:"credentialState"`／`json:"credentialKnown"`，且它们**声明在 `ComposerState` 里**（＝`composer` 段内，不是顶层），⛔ 没有顶层第五键；`NewComposerState`（`:273`，unknown 支在 `:288`）在无 reader 时把 `Credential` 置成 `CredentialUnknown`，`pump.go:260` 只有 `p.src.Credential != nil` 才填那两枚（`:268-269`）、且把空串折回 `CredentialUnknown`——**"没人读过"与"没录过"分得开**这一条本腿逐行读了，不是照录实现件。
- 常驻用例〔**我复认过名册**〕：`grep -n "^func Test" cmd/wisp/panel_config_248_test.go` ⇒ 9 枚顶层，含 `TestAC2CredentialSentinelAppearsInNoArtifact`（`:151`）、`TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere`（`:215`，正控）、`TestAC2SharedEnvelopeCannotCarryTheCredentialValue`（`:326`）、`TestAC2SnapshotReportsRefWithoutBlobAsAPartState`（`:507`）；本腿另跑一把尺：三枚新测试文件里 `t.Skip`／`testing.Short` **零命中**＝不是靠跳过装绿〔**我复认过**〕。
- 零命中尺的射程＝6 枚表面（回执/audit/slog/快照/ledger 行/数据根全文件）＋两枚正控：实现件 §5-3 自陈第一版漏了 ledger 那一面、后来补进尺与正控——**这条"读数分两发"的自暴我照收，且它恰是本格能被判成立的前提**（少了 ledger 那一面，AC#2 的"持久 sink"就没尺）。
- 尺有牙齿〔**我复认过（本腿 §3 V-2）**〕：本腿在**产码**里把凭据值逐字接进回执那一行（`panel_config_store.go:277` 的 `res.Note`）⇒ 常驻用例 `TestAC2CredentialSentinelAppearsInNoArtifact` **红**，红句点名表面是 `receipt`、报两枚命中（字面哨兵那一枚＋`sk-` 形状那一枚），⛔ 不打印值；同发常驻正控 `TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere` **绿**；还原后两枚同绿。⇒ "零命中不是因为尺瞎"这一条本腿自己抓到了。前腿那两发（`mutation-M3.txt` 种进产物那一发、`mutation-M4.txt` 把凭据键声明进共用封套 ⇒ `TestAC2SharedEnvelopeCannotCarryTheCredentialValue` 与 `TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope` 双双指名红）〔**我没复认，照录前腿**；种法两枚都没记，见 §3 末表〕。
- 洗红的那一枚对账门〔**我复认过（读数在 §2）**〕：`internal/panel` 整包里 `TestComposerContractTypesMatchFrontend` 与 `TestApprovalCardViewJSONKeysMatchFrontendTypes` 红，红句逐字是 Go 侧 emit 的键集里那两枚新键"interface ComposerState does not declare"＝**双向对账门在界面侧补声明之前必然差一枚**。⚠ 这不是"实现者放宽断言"——本腿核到那两枚门的断言原文在 `0d87a681` 里**一字未动**（`git show --numstat` ⇒ `composer_test.go` 零改动、`composer_dispatch_test.go` 只 2/1 行的签名调用点）。本格两栖，Go 半枚做不到 self-green，实现件 §6-4 已具名交回。

### AC#3 —— 只走一枚凭据存储、不落明文、不新增回显方法

**判语：成立。**
- 存储侧〔**我复认过**〕：`0d87a681` 的 15 枚路径名册里**零枚 `internal/secret/**`**（尺＝`git show --numstat --format='' 0d87a681`）⇒ "没新造第二套存储"是名册级事实，不是读码印象。
- 写侧〔**我复认过**〕：`internal/config/settings.go` 六枚导出 setter（`:60/:81/:103/:132/:137/:177`）全部经同一枚 `m.mergeWrite`（唯一调用点在 `:242`）；`:79` 注释逐字写明那枚字段收的是**引用**（`dpapi:<id>`／`env:NAME`）而不是密钥，且 `validateAPIKeyRefs` 会拒裸密钥——本腿核到全仓新代码里没有任何把值写进 `config.toml` 的调用点（同一把 `grep` 尺：`SaveFile` 只出现在注释里，见 AC#11）。
- 用例侧〔**我复认过（本腿逐名跑过、带 `-v`，见 §2 末段）**〕：本腿与 `roster-config-panel-248.txt` 同名同终态 ⇒ `TestAC3SettingWriteStoresARefNeverAValue`（`internal/config`）与 `TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef`（`cmd/wisp`）双双 PASS；`mutation-M2.txt` 的另一半显示种非法引用时 `internal/config` 红而 `cmd/wisp` 那枚 AC7 用例仍绿——那一发本腿**自己重跑**（§3 V-1），因为它同时是 AC#7 的落点凭据。
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
- 用例＋正控〔**我复认过（本腿 §3 V-1）**〕：本腿自己把 `settings.go:230` 那枚 `validate()` 闸门换成永不成立的判据，量到的形状**比前腿那发更细**：`TestAC7InvalidValueLeavesTheFileByteIdentical` 六枚子用例里**只有两枚转红**——`plaintext_credential_in_the_reference_field` 与 `chat_model_that_is_not_in_the_catalog`（红句逐字 "was accepted into config.toml"／"AC#7's gate did not run"），另外四枚（负上下文窗、负价格、文件没声明的 provider、provider 目录里没列的 model）**仍然绿**＝那四类由 setter 自己的字段级检查拦、不靠 `validate()`。同发 `TestAC11MultiFieldSaveIsOneWritePerKeyAndStopsOnFailure` 转红。⇒ 本格的诚实口径是：**"落盘前过同一套 `validate()`"这条判据成立且它不是摆设（两类非法值只有它拦得住）**，但"六发非法值磁盘不变"这句话的守门人是**两层**（字段级检查＋`validate()`），不是一层。⚠ 同一发里 `cmd/wisp` 的 `TestAC7InvalidSettingsValueIsRefusedBeforeTheFile` **仍 PASS**＝两枚同名用例视野不同，**守磁盘那一格的是 `internal/config` 那枚**；这一点本腿复认到并写在这里，前腿 `mutation-M2.txt` 的读数与本腿同形（照录，见 §3）。还原尺痕与逐名终态在 §3。

### AC#8 —— 回执带生效档；不许出现"保存即生效"；"要重启"要到页面可见面

**判语：附条件成立。**
- 词面判据〔**我复认过**〕：`grep -rn "保存即生效" cmd internal tools docs scripts .scratch/wisp/issues` ⇒ 命中只在**三处断言的反面**（`cmd/wisp/panel_config_248_test.go:433`、`internal/panel/config_route_248_test.go:325`／`:381`，都是 `strings.Contains(...) ⇒ 报错`）、实现件那两行、以及票面自身。**产码与回执文案里零命中**＝票面那句"不许出现"按词面与按能力两面都过（`config_route_248_test.go:320` 还钉了"必须说重启"）。
- 回执带档〔**我复认过**〕：`internal/panel/config_handlers.go:199-201` 三枚档常量＋`:430-434` `tierSentence` 逐档成句；`renderSettingReceipt` 逐枚列键路径（`config_route_248_test.go:330` 那枚断言就是"两枚键都要出现"）。
- 页面可达面〔**我复认过装配层**〕：本腿在 HEAD 重取那两枚锚（不引实现件 §3#8 的旧行号，因 `cmd/wisp/run.go` 的 blob 自 `b644d310` 起已变）：`cmd/wisp/run.go:431` `rt.settings = newConfigStore(mgr, st, rt.auditf, "cli-run")`、`cmd/wisp/run.go:710` `Credential: rt.settings.credentialStatus`——同一枚 Manager、同一枚 Store，没有第二份加载。回执离进程的那一跳（`Handle` 返回串经票 33 那枚 bind 闭包交页面）〔**我没复认，照录实现件 §3#6/§7**；本腿不引页面层路径〕。
- **条件（为什么不是"成立"）**：本腿复认到**档位词表在这里对不上**——`cmd/wisp/panel_config_store.go:228`（凭据那一发在 `:275`）**无条件** `res.Tier = panel.EffectiveRestart`，前面没有任何按键/按段的判断；而 `internal/config/manager.go:281` 把 `llm` 段列在**热应用段表**里（同表还有 `panel`）；`EffectiveNow`／`EffectiveNextTask` 两枚档在**非测试产码里零枚写者**（本腿 grep 尺：只命中定义、`tierSentence` 的 case 与测试）。⇒ 同一项设置，面板回执说"重启才生效"，热加载那条路把它当立即档，两套词表零行码互相翻译。**这一问算不算破 AC#8 本腿不当场结论化——已按派单要求写进 §5 第 1 条交裁。**

### AC#9 —— 宿主在、内容不在（embed 之后页面产物真存在）

**判语：判不了（本格在本编队的禁令面之外，按票面原文"落定前勾不了"记不成立）。**
- 本腿唯一能给的仍是**名册级**事实（页面产物那一层的跟踪名册里只有一枚占位件），且这一条**本腿没有重跑**——那把尺会列出禁令层里的文件名，所以本件连它的输出都不引。⛔ 本腿不据"能构建"（§2 build rc=0）反推"有内容"——票面 AC#9 原文就禁止这么读。
- 归口：谁把页面产物填上没落定之前这一格勾不了（票面 AC#9 同判）。

### AC#10 —— 常驻腿吃不吃 `[risk]` 那两项：接上、或写明是常量

**判语：不成立（两形都没落；选形归编排者，故不计实现腿越权）。**
- 本腿自己复认那三行读数〔**我复认过**〕：`cmd/wisp/resident_approval_windows.go:109-113` 的 `approval.New(approval.Options{…})` 只给 `UI`/`Channels`/`Logf` 三项——**没有 `Window`、没有 `ApprovalTimeout`、没有 `Grants`**；`[risk]` 那两项在非测试产码里的读取点只有 `cmd/wisp/run.go:615-616`（跑任务那枚 gate）；常量钳位在 `internal/agent/approval/queue.go:107`（300s）／`:116`（3s）／`:120`（2s）／`:122`（3s）。⇒ 票面那句"常驻腿今天用常量"在 HEAD **仍然成立**，`GRANT-DROPPED` 那一族（票面引 `gate.go` 两行）本腿未跑仪器复认，且实现件 §6-1 自陈也没跑。
- ⛔ 未选形、未动常驻门（`0d87a681` 名册里没有 `resident_approval_windows.go`）＝纪律层成立、判据层未交付。票面 10-01 12:0x 那枚"本轮不裁"的裁在 HEAD 没被推翻（编排者至今没落那枚具名 `A##`）。
- ⇒ 本格的账记在**编排者的选形**上；ⓑ 形那句话一旦要写，它同时是 AC#8 那一族的第二次"文案必须说真话"。

### AC#11 —— 只走按键受保护写；三条子判据

**判语：成立。**
- 子判据 3（落点）＝AC#7 那一格〔**我复认过**，同上〕。
- 子判据 1（并发保 B 键）：形状侧〔**我复认过**〕——`internal/config/settings.go` 里 `SaveFile` **只出现在注释里**（`:5`／`:6` 那两行解释它为什么是序列化器不是合并），受保护写只有 `:242` 那一处 `m.mergeWrite`；`writeguard.go`／`loader.go` 均不在 15 枚路径名册里＝票 226 修好的洞没被重新打开。行为侧〔**我复认过（本腿逐名跑，带 `-v`）**〕——`TestAC11SettingWriteKeepsAForeignHandEditedKey`、`TestAC11ControlSnapshotWriteIsWhatRevertsTheHandEdit`（那枚"控制必须仍然会覆写"的正控）、`TestAC11WritingTheValueTheFileAlreadyHoldsChangesNoBytes`、`TestAC11SettingWriteReportsTheKeyPathItChanged` 逐名 **PASS、零 SKIP**（取数件 `probes/248/v1c/roster-config-and-ac2.txt`；前腿 `roster-config-panel-248.txt` 同名同终态）。
- 子判据 2（多枚字段＝逐枚键各一发＋回执逐枚列键路径）：`writeOneKey` 每次只交一枚 `ownedKey` 进 `mergeWrite`，回执取"文件前后差"（`settings.go:249` `m.writtenKeyPaths(pre, preExists, ownedKey)`，函数在 `:254` 起）；`TestAC11MultiFieldSaveIsOneWritePerKeyAndStopsOnFailure` 与 `TestAC11SettingWriteReportsTheKeyPathItChanged` PASS〔**我没复认，照录前腿**〕，而 `mutation-M2.txt` 那一发让它指名红＝它有牙齿〔同上，本腿在 §3 V-1 里自己重跑同一族〕。
- ⚠ 一条已知宽度（不是缺陷，是实现者自暴）：回执报的是"文件实际发生了什么"，第三方同一瞬间的改动也会被列进来（实现件 §5-5）。本腿采它的理由：票面 AC#11-2 点名要 `diffKeyPaths` 那一族的"实际落盘键路径"，求交反而背离题面。⇒ 记在 §5。

---

## 2. 四把尺现值（带 HEAD）＋三包整包名册

> 本节全部是本腿 `248-v1c` 自己在 **HEAD `6bcb934a`** 上跑的（流水 `probes/248/v1c/gate-run.log`，尺脚本 `gate-run.sh`，逐名终态在 `build.txt`／`vet.txt`／`gofumpt-touched.txt`／`d22scan.txt`／`test-panel-config.txt`／`test-cmdwisp-full.txt`）。⚠ 时序类读数按【同机 `254-r1b` 在飞、可能被抢 CPU】标注。

| 尺（票面 AC#6 点名的四发） | **本腿读数**（`2026-10-02 11:04:43–11:08:05 +0800`） | 前腿 `248-v1b` 同尺（HEAD `580d6153`） | 复认 |
|---|---|---|---|
| `GOFLAGS= go build ./...` | **rc=0**（零输出） | rc=0 | 〔我复认过〕 |
| `GOFLAGS= go vet` | **rc=0**，⚠ 本腿的射程是**三枚被测包**（`./cmd/wisp ./internal/panel ./internal/config`），不是 `./...` | 它记的是 `vet rc=0`（射程那件里没写） | 〔我复认过，且**比前腿窄**——这条差异我在 §5 不藏：本腿没跑全仓 vet〕 |
| `gofumpt -l cmd/wisp internal/config internal/panel` | **list-count=0**（三目录零枚未格式化） | 同尺 count=0 | 〔我复认过〕 |
| `sh scripts/d22scan.sh`（独立模块那枚 exe） | **rc=0**，末行逐字 `clean - no D22 ban violations`，`examined 262 production Go files under internal/ and cmd/` | rc=0 | 〔我复认过〕 |
| `go test ./cmd/wisp ./internal/panel/ -count=1`（票面第三发，两包同发） | **合并 rc≠0**：`internal/panel` **rc=1／4 枚指名红**＋`internal/config ok 0.842s`（本腿把 config 并进同一发跑）＋`cmd/wisp` **ok 177.139s／rc=0** | `internal/panel` 同名册 rc=1；`cmd/wisp` **FAIL 255.772s／2 枚指名红**（`TestTicket223HandEditedFsLooseningCostsAnL2Card`／`TestAC14GoSideEvalPushReachesThePage`），它另有一发隔离跑两枚都 PASS | 见下面两行 |

- **`internal/panel` 那 4 枚红：本腿与前腿逐名同册**〔**我复认过**〕：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`。
- **`cmd/wisp` 那 2 枚红：本腿未复现**〔**我复认过＝复认的是"它今天不红"**〕：本腿整包一发 `ok 177.139s`。⇒ 那两枚属时序/争用那一族（见 §5 第 4 条的边界）。

**那 4 枚红的归因（本腿自己的三把尺，不是抄别人）**

1. **零 prod 漂移**：`git diff --name-only 580d6153..HEAD -- cmd internal tools` ⇒ **空**；两枚 HEAD 上同名同册 ⇒ 与 `580d6153` 之后落的那批提交无关。
2. **红句自己报的原因**（本腿只读红句、⛔ 未读那两层任何一字节）：前两枚是**两界双向对账门**，红句逐字是"Go 侧 emit 的键集里有几枚 interface 那侧没声明"（今天的清单里含本票新增的两枚凭据维键＝AC#2 那一族，见下条）；第三枚红句是"打包真会带出去的那份代码里出现了第二处样式源"；第四枚红句是"读一枚设计层令牌文件失败：系统找不到路径"——而**那枚文件在本腿的起手名册里就是 ` D`（工作树里已被删除、未提交，owner 的移动）**，本腿没读它、也没读它所在那层任何文件。
3. **独立前史（这条最硬，且不出自实现腿）**：`docs/evidence/s1/181-rewrite-account-r2.md:223`（该文件首枚提交 `5d619b4f`，`2026-09-28`）已把其中**三枚**逐名记成"在册红照实记"，早于 `0d87a681` 五天；同批 `145-snapshot-growth-r3.md:117`／`188-task-state-r1.md:188` 同记。⇒ 本票不是这四枚红的**成因**；本票**只做了一件事**：让其中那枚对账门的红句里多列两枚键（`credentialState`/`credentialKnown`）——这正是票面把 AC#2 判成"两栖、Go 半枚做不到 self-green"的形状。

**248 自己那批用例的逐名终态（本腿自跑，非照录）**

- `internal/config`：整包 `ok 0.842s`；248 那批**本腿逐名跑过、带 `-v`**（`roster-config-and-ac2.txt`，11:54:12）⇒ `TestAC11SettingWriteReportsTheKeyPathItChanged`／`TestAC11MultiFieldSaveIsOneWritePerKeyAndStopsOnFailure`／`TestAC7InvalidValueLeavesTheFileByteIdentical`＋**6 枚子用例全 PASS**／`TestAC11SettingWriteKeepsAForeignHandEditedKey`／`TestAC11ControlSnapshotWriteIsWhatRevertsTheHandEdit`／`TestAC3SettingWriteStoresARefNeverAValue`／`TestAC11WritingTheValueTheFileAlreadyHoldsChangesNoBytes`＝**10 枚（含子用例）全 PASS、零 SKIP**。⚠ 顺带量到一枚命名撞车：`-run 'TestAC3'` 在这包里还捞出两枚**不属于本票**的用例（`TestAC3AdoptionClaimsOnlyWhatThisWriteProduced`／`TestAC3HandAddedEntryIsPreservedButNotAppliedToMemory`，票 226 那一族）——本腿没把它们记进 248 的名册。
- `internal/panel` 里 248 那 11 枚顶层＋2 枚子用例：**本腿自己跑过名册**（`-run 'TestAC1|TestAC2|TestAC8' -count=1 -v`，11:12:26，取数件 `panel-248-roster-v1c.txt`）⇒ **11 枚 `--- PASS`、零 FAIL、零 SKIP**，含 `TestAC8ReceiptNamesTheTierAndNeverClaimsLiveSave` 与三枚 `TestAC2*`（本腿第一次那发把三枚 `TestAC2*` 截在了 `head` 之外，所以本腿重跑了一次逐名数过＝11）。前腿同批名册在 `probes/248/v1b/roster-config-panel-248.txt`〔**与本腿同名同终态**，本腿复认成立〕。
- `cmd/wisp` 里 248 那 9 枚：**本腿逐名跑过、带 `-v`，9 枚全 PASS、零 SKIP**（`V2-restored-ctrl.txt` 5 枚＋`roster-config-and-ac2.txt` 里 `TestAC2SharedEnvelopeCannotCarryTheCredentialValue`／`TestAC2SnapshotReportsRefWithoutBlobAsAPartState` 2 枚＋`roster-cmdwisp-last-two.txt` 末 2 枚）。⇒ 本腿**不再依赖**"整包 `ok` 蕴含逐名跑过"那种读法：那发不带 `-v`，其存档件里 `grep -c '^=== RUN'`＝**0**、`--- SKIP`＝0，两枚都是"没采到"而不是"零枚跳过"——⛔ 不许把 `ok` 读成"逐名跑过且真断言了"（这条口径写在 §5 第 11 条）。

## 3. 变异清单

> 本腿只种"判语所缺的那几发"（⛔ 不把 M1–M6 重跑一遍——那是前两枚腿的死因）。四发都留"还原后 `git diff` 该文件为空"的尺痕；取数件同名放在 `.scratch/wisp/probes/248/v1c/`。
> ⛔ 四发变异**零提交**：每发跑完立刻 `git cat-file blob HEAD:<path> > <path>` 还原；交件前最后一把尺＝`git status --porcelain -- cmd internal tools` **0 行**（在下面"收尾尺痕"那条）。

| # | 种在哪（文件:行） | 种法（一句话） | 指名红的用例与红句 | 还原尺痕 | 本格判语的哪一句靠它 |
|---|---|---|---|---|---|
| **V-1** | `internal/config/settings.go:230` | 把 AC#7 那枚 `validate(candidate)` 闸门换成永不成立的判据（⛔ 不删别的行） | `TestAC7InvalidValueLeavesTheFileByteIdentical` **只红两枚子用例**：`plaintext_credential_in_the_reference_field`、`chat_model_that_is_not_in_the_catalog`（红句逐字"was accepted into config.toml"／"AC#7's gate did not run"）；另两枚同族：`TestAC11MultiFieldSaveIsOneWritePerKeyAndStopsOnFailure` 红。**四枚子用例仍绿**（负上下文窗、负价格、未声明的 provider、目录里没有的 model）＝那四类由 setter 自己的字段级检查拦，不靠 `validate()`。⚠ 同发里 `cmd/wisp` 的 `TestAC7InvalidSettingsValueIsRefusedBeforeTheFile` **仍 PASS** | `git diff --numstat`＝**0 行**；还原后 `-run 'TestAC7\|TestAC11MultiFieldSave' ./internal/config` ⇒ **ok** | AC#7 的"闸门真在写之前、且它拦得住东西"＝〔我复认过〕；同时把"哪两类只有它拦"这条本腿自己量到的射程写进 AC#7 判语 |
| **V-2** | `cmd/wisp/panel_config_store.go:277`（`StoreCredential` 的 `res.Note`） | 把凭据值**逐字接进回执那一行**（＝一条真实的产码泄漏形状，不是测试输入） | `TestAC2CredentialSentinelAppearsInNoArtifact` **红**：`the credential value reached an artifact (2 hits)`＋`leak surface: receipt matches sk-[A-Za-z0-9]{12,}`／`leak surface: receipt matches <canary-redacted>`（红句只报尺与表面，⛔ 不打印值）；同发常驻正控 `TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere` **PASS** | `diff-lines=0`；还原后两枚用例＋AC1/AC3/AC7 共 5 枚 **全 PASS**（`V2-restored-ctrl.txt`） | AC#2 的"零命中尺看得见东西"＝〔我复认过〕；⛔ 前腿 `mutation-M3.txt` 的种法本腿仍然不知道，那一发不署我的名 |
| **V-3** | `internal/panel/bridge.go:68`（同一枚 const 块） | 种一枚守卫不认的路由常量 `MethodConfigFoo = "config.foo"`（⛔ 不动 `knownComposerMethod` 的 case 列表） | `TestPlantedGrantWiringGoesRedInASnapshot` **红**：`l2_grant_boundary_test.go:2208: plant B has no anchor … declared Method* constant MethodConfigFoo is absent from the guard's case list: a route that exists but is never asked of the guard` | `git diff -- internal/panel/bridge.go` **空**；还原后同一条尺 **ok 0.266s** | AC#1 里 J1 那枚**能力形锚的正控**＝〔我复认过〕；与前腿 `mutation-M6b-msg.txt` 的红句逐字同形（＝前腿那发读数本腿认可为真，但署名各归各） |
| **V-4** | `internal/panel/config_handlers.go:122-123`（`fs.` 那一支的理由句） | 把"点名该键真正住在哪"那句换成一句光秃秃的"该字段不可写" | `TestAC1LockedFamilyFieldsAreRefusedBeforeTheLeg` **红**：`config_route_248_test.go:236: refusal for "fs.allowed_dirs" does not name where that key does live: … 属于 fs. 族，本路由不写它：该字段不可写`；同发另两枚（`TestAC1FieldTableCarriesNoApprovalVocabulary`／`TestAC1SelectorAndEmptyValueRefusals`）**绿** | `diff-lines=0`；还原后 `-run 'TestAC1\|TestAC2\|TestAC8' ./internal/panel` ⇒ **ok** | AC#1 的"负向配正控、且拒写得说出去哪儿"＝〔我复认过〕（这一发同时把前腿 `mutation-M5.txt` 那发红句复现到逐字同形） |

**收尾尺痕（交件前一把，逐字）**：`git status --porcelain -- cmd internal tools` ⇒ **0 行**；本腿全程⛔未 commit 任何产码/测试文件（尺＝`git log --oneline` 里本腿那几枚的 pathspec 只有 `docs/evidence/s1/248-settings-write-path-v1.md` 与本腿台件目录）。

**照录前腿、本腿没复跑的那几发（⛔ 不署我的名，具名列出）**

| 前腿件（`probes/248/v1b/`） | 它记到的读数（一句话） | 本腿为什么没复跑 |
|---|---|---|
| `mutation-M1.txt`＋`mutation-M1-restored-ctrl.txt` | 种在 `bridge.go` 的一发：`internal/panel` 4 枚红＋`internal/config` ok；`cmd/wisp` 那一发红的是 `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`，**还原后红的是另一枚 `TestAC1ResidentLegInstallsItsLogListenerOnDisk`** | ⚠ 这一对本腿有用但不是任何一格的凭据：它说明**还原前后红的不是同一枚常驻用例**＝那一族是时序抖动，与 §2 里 `cmd/wisp` 整包"前腿红、本腿绿"互相印证。⛔ 且它没记种法，本腿不据此判任何格 |
| `mutation-M4.txt` | 把凭据键声明进共用封套那一形 ⇒ `TestAC2SharedEnvelopeCannotCarryTheCredentialValue` 与 `TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope` 双双指名红 | AC#0/AC#3 那一维本腿另有凭据（名册＋常量表）；这一发的成本是再编译两包，收益不抵 §4 的时限 |
| `mutation-M6a.txt` | 把守卫名册"读短"那一发 ⇒ `:2312` 报"passing for the wrong reason"（能力形锚的**另一侧**牙齿） | 与本腿 V-3 是同枚锚的两个方向；V-3 已足够支撑 AC#1，另一向照录不重跑 |
| `rulers-*.txt`／`cmdwisp-full.txt`／`roster-config-panel-248.txt` | 见 §1/§2 逐条标注 | 四把尺与两包整包本腿**全部自己重跑过**（§2）；`roster` 那枚里本腿复认了 `internal/panel` 与 `cmd/wisp` 两侧 248 用例名册（本腿自跑读数在下面这条） |

**本腿自跑的 248 用例名册（不是照录）**：`internal/panel` ⇒ **11 枚顶层＋2 枚子用例全 PASS、零 SKIP**（`panel-248-roster-v1c.txt`，11:12:26）；`cmd/wisp` ⇒ 5 枚指名 PASS（`V2-restored-ctrl.txt`，含端到端那枚 `TestAC1InboundLegAnswersSettingsRouteEndToEnd`），另 9 枚整包在本腿 `ok 177.139s` 那一发里全跑到；`internal/config` ⇒ 整包 `ok 0.842s`（前腿逐名 PASS 名册 10 枚含 6 枚子用例，本腿复认同名册）。

## 4. 退回与否

**总裁：不退回产码；退回一格证据件；其余按逐格处置入账。** 取数时刻 `2026-10-02 11:0x +0800`，HEAD `6bcb934a`（prod 面与 `0d87a681`＋`b644d310` 的差集只有 `cmd/wisp/run.go`，已在 §0c 具名归因）。

| 格 | 处置（三枚词之一：追认／附条件入账／退回） | 若判"可勾"，还欠谁一句话 |
|---|---|---|
| AC#0 | **追认** | — |
| AC#1 | **追认**（含 J1 那两枚锚：越界尺证明只动了解冻具名的两处） | 两界登记点（页面会叫名之后那一步）⛔ 不算本格的账 |
| AC#2 | **附条件入账** | 条件具名＝界面侧在契约声明面补那两枚键 ⇒ 两界对账门转绿；条件落定前**不许勾**（票面自己说这一格两栖） |
| AC#3 | **追认** | — |
| AC#4 | **不成立·挂账**（票面明令今天不许勾，⛔ 不是退回项） | 界面侧入口＋票 33 真机那一发 |
| AC#5 | **追认**（零越界，本腿自己的 git 尺） | — |
| AC#6 | **退回**——退的是 `docs/evidence/s1/248-settings-write-path-r1.md` §4 那一格 | 要求：把四数＋两包逐名终态**补写进它 §4**（追加新 commit，⛔ 不改写 `0d87a681`/`b644d310` 已推送的历史），并把 §7 里"四数在 §4（全 rc=0）"那半句**改到与终态读数一致**——那半句现在是假的：它指的节是空的，且两包整包 rc≠0 |
| AC#7 | **追认** | — |
| AC#8 | **附条件入账** | 条件＝§5 第 1 条那一问由编排者裁："把 hot 说成 restart"这种**同样不成立的实话**算不算破 AC#8。裁成"算破"⇒ 回执那一族必须改（那是票 255 的登记表地界，⛔ 不是本票去加重建跳）；裁成"不算破"⇒ 本腿这格转追认 |
| AC#9 | **不成立·挂账**（本编队判不了，禁令面＋"谁填 dist"未定） | owner 带话 |
| AC#10 | **不成立·挂账**（两形都没落；⛔ 选形归编排者，不计实现腿越权） | 编排者落那枚具名 `A##` |
| AC#11 | **追认** | 回执归因宽度见 §5 第 7 条 |

**这一格为什么值得退回而不值得整体退回**：本票 15 枚路径的**行为面**本腿逐格找到了凭据或正控；缺的是**交件面**的一节正文＋一句过度陈述。按 `AGENTS.md §0.2`（改契约＝人工批准）与票面 AC#5 的退回触发项（只有"出现禁区路径"一条），本腿**无权**因证据件的排版缺项判代码退回，也**不该**因代码没红就放过那句假读数——所以处置是分开两格记。

## 5. 我可能判错的条目

> 这一节是判语的一部分。⛔ 下面每一条都写得足以让下一枚腿反驳我。

1. **AC#8 那一问（派单点名要我写进来的，我不替编排者答）**：票面 AC#8 判的是"回执不许说'保存即生效'"这一类**说谎**。本腿量到的形状是另一件事——**它说了一句同样不成立的实话**：`cmd/wisp/panel_config_store.go:228`（凭据那一发在 `:275`）**无条件**把档位写成 `restart`，前面零按键、零按段判断；而 `internal/config/manager.go:281` 把 `llm` 段列在**热应用段表**里（同表还有 `panel`），`EffectiveNow`／`EffectiveNextTask` 在非测试产码里**零枚写者**。⇒ 我的读法是：对 `[llm]` 这七枚字段（本票唯一可写的集合），"重启才生效"按 J9 的实践面**是真话**，所以 AC#8 的词面判据与能力判据都过；但"两套词表零行码互相翻译"意味着**下一枚往热应用段里加字段的腿会说谎**，且它谎的方向是保守那一侧。**我可能错的点**：如果票面那句"回执文案不许出现'保存即生效'"被裁成"回执的每一枚档都必须由同一张登记表产出"（＝票 255 新补 AC#5 的形状），那 AC#8 就不是附条件成立而是**不成立**，本票 §3 里那条正控（改一枚标 hot 的键、回执不许说重启）也没人替它做。⛔ 我没有把这一格判成不成立，因为那等于替编排者裁票 255 的地界。
2. **AC#6 里"假报"这个定性可能过重**。实现件 §0 自己写了那枚腿是**死腿代提**（`b644d310` 的标题逐字含"死腿收尾代提"），它 §4 的占位更可能是"没来得及写"而不是"写了假的"。但 §7 那句"四数在 §4（build/vet/d22scan/gofumpt 全 rc=0）"是**一句主动陈述**，且与两腿独立读到的 rc≠0 相矛盾——本腿按盘面事实记"假报"，若编排者改记"未完成"，**AC#6 不成立本身不变**，变的只是要不要在台账里挂它一笔。
3. **4 枚 `internal/panel` 指名红的归因**，是我这轮最可能被翻的一条判语。我的凭据是三条：① 本腿 HEAD 与 v1b 的 `580d6153` 之间**零动 prod**，两枚 HEAD 同名册（与 198/155 那批无关）；② 红句自己报的原因是两界声明面／样式源／设计层令牌文件此刻在工作树里处于已删状态（`git status` 名册可见那几枚是 ` D`，不是本腿读了它们的内容）；③ 独立前史——`docs/evidence/s1/181-rewrite-account-r2.md:223`（该文件首枚提交 `5d619b4f`，2026-09-28）就把 `TestComposerContractTypesMatchFrontend`／`TestPanelColourLiterals…`／`TestC21DesignTokensFourWayAgree` 三枚**逐名照实记成红**，早于 `0d87a681` 五天，票 145-r3／188 同记。**我可能错的点**：我今天这 4 枚里第 4 枚（`TestApprovalCardViewJSONKeysMatchFrontendTypes`）不在 09-28 那枚名册里，且第 1 枚的红句今天**多列了本票新增的两枚键**——也就是说本票**让一枚本来就红的门更红**（红的名册没变，红的清单变长）。如果编排者读票面 AC#6 的意思是"交件时终态必须 rc=0"，那我这条归因救不了 AC#6，AC#2 的"附条件"会升级成整片退回。
4. **`cmd/wisp` 那两枚指名红的"时序争用"归因我只部分站得住**。前腿 `248-v1b` 的整包跑（`cmdwisp-full.txt`，HEAD `580d6153`，10:38–10:42）红两枚：`TestTicket223HandEditedFsLooseningCostsAnL2Card`／`TestAC14GoSideEvalPushReachesThePage`，它另有一发隔离跑（`cmdwisp-iso-ac14-223.txt`）两枚都 PASS。**本腿自己在 HEAD `6bcb934a` 上重跑整包那一发是 `ok 177.139s`、rc=0**（§2）⇒ 我复认到"这两枚在整包里今天不红"，也就复认了"它们是争用/时序那一族，不是稳定红"。⚠ 我**没有**拿到"同一时刻、同一负载下整包红＋隔离绿"的同发对照（那要把整包跑两遍，超预算），也**没有**指认是谁抢的 CPU。⇒ 若下一枚腿在整包里又采到这两枚红，本腿这句"已复认"要退回成"未定"。
5. **前腿那六发变异我只照录了读数、没照录到种法**（`probes/248/v1b/mutation-M*.txt` 通篇只有时间戳＋终态名册，没写改了哪一行）。⇒ 若某一发其实种在了测试自己的输入上（＝正控打在纸面上），我照录那一格就站不住。本腿自己重种了四发（§3 的 V-1/V-2/V-3/V-4）把**四格判语所依赖的正控**换成自己的读数，其中 **V-4 复现出与前腿 `mutation-M5.txt` 逐字同形的红句、V-1 与前腿 `mutation-M2.txt` 同终态（且本腿量得更细：只有 6 枚子用例里的 2 枚吃 `validate()`）**⇒ 那两发前腿读数本腿认可为真。⚠ 仍**没复现**的三发：`M3`（种进产物）、`M4`（把凭据键声明进共用封套）、`M6a`（把守卫名册读短）——前两发的结论本腿另有 V-2 与名册尺顶着，`M6a` 那一向（锚的另一侧牙齿）本腿**只有前腿的读数**，AC#1 里我把它单列成〔照录前腿〕。
6. **AC#1 的"判据一律问能力、不许扫注释词面"我是按名册＋红句形状判的，没逐枚读完 11 枚用例的断言原文**。其中 `TestAC1FieldTableCarriesNoApprovalVocabulary` 名字像词面型；本腿采它的理由是 `mutation-M5.txt` 那发红句报的是"某枚锁定族字段没说出真正住在哪"（能力形），以及同批 `TestAC1LockedFamilyFieldsAreRefusedBeforeTheLeg` 断言的是拒写＋处理器调用计数。⛔ 如果那枚"审批词汇"用例其实只是 `strings.Contains(常量表)` 的注释扫描，AC#1 的这句判语要改写成"除一枚外皆问能力"。
7. **AC#11 我采信了"回执＝文件前后差"的宽度**（实现件 §5-5 自暴：第三方同一瞬间的改动也会被列进来）。我采它的理由＝票面 AC#11-2 点名要 `diffKeyPaths` 那一族的"实际落盘键路径"。**我可能错的点**：若编排者把那句读成"必须等于本次所有权"，AC#11 要改判，且改法是把 `written` 与 `ownedKey` 求交——而那恰恰是票 226 反对的方向，所以我认为这一读法不成立，但它是**我的读法**。
8. **AC#10 我用"不成立·挂账"可能比票面想的重**。票面自己写了本格"只许两选一"，而选形那一半是编排者 10-01 12:0x 明令"本轮不裁"的。⇒ 严格讲这一格今天是**尚未生效**而不是"交付失败"。本腿仍写"不成立"，是为了让下一枚腿不去勾它；如果台账更想要"挂账不判"这一枚词，改词不改账。
9. **§0c 那把 blob 尺里 `run.go` 一枚 DIFF 我归给票 198**，凭据是前腿 §0a 记的"`0d87a681..HEAD` 里动过被测三包的只有 `613606c0`"＋本腿自己现读的 `run.go:431/:710` 两枚锚仍在。**我可能错的点**：本腿没有逐枚读 `613606c0` 的 diff 内容，也没排除"第三枚腿在同一路径上的增量"。这不牵动判语（判语用的是 HEAD 现读），只牵动"这句归因"的措辞。
10. **AC#2 的"没有顶层第五键"本腿改成〔我复认过〕，但留下窄窄的一条没盖住**：本腿逐行读了 `composer.go:268-269`（两枚 tag 在 `ComposerState` 里）＋`:273/:288`（无 reader ⇒ `CredentialUnknown`）＋`pump.go:260/268-269`（reader 为 nil 就不填那两枚）。⛔ 本腿**没有单独跑**"未接 reader ⇒ `credentialKnown:false`"那一发的逐名终态——它的断言原文在 `cmd/wisp/panel_config_248_test.go:544`（本腿读过那三行），而那枚用例（`TestAC2SnapshotReportsRefWithoutBlobAsAPartState`）只在本腿**不带 `-v` 的整包 `ok`** 那一发里跑到 ⇒ 本腿对它只有"没 FAIL"的凭据、没有"逐名 PASS 且未 SKIP"的凭据（口径见 §2 末段那两条）。
11. **追加一发（交件前最后一把尺，11:55）＋一条仪器口径**：上面第 10 句"只有『没 FAIL』的凭据"已**当场补掉**——本腿把 `cmd/wisp` 那 9 枚 248 用例**逐名**带 `-v` 跑完（`TestAC2SnapshotReportsRefWithoutBlobAsAPartState`／`TestAC2SharedEnvelopeCannotCarryTheCredentialValue`／`TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope`／`TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket` 等；取数件 `roster-config-and-ac2.txt`＋`roster-cmdwisp-last-two.txt`）⇒ **9 枚全 PASS、零 SKIP**，`internal/config` 那批同样逐名（10 枚含 6 枚子用例全 PASS）。⇒ §1 里 AC#1／AC#3／AC#11 三格的"用例侧"标记从〔照录前腿〕改成〔我复认过〕，⛔ **十二格判语一枚没改**。
    ⚠ 顺手留给下一枚腿的口径：**"整包 `ok`"不证明逐名跑过、也不证明没跳过**（不带 `-v` 时存档件里 `=== RUN` 计数＝0，而有 SKIP 时 Go 照样打 `ok`）；要写"常驻用例真在跑"这种凭据，必须带 `-v` 逐名抄。本腿第一版 §2 差点把 `ok` 读成后者，具名留在表里。

## 6. 判不动的地方

> 逐条具名交回，⛔ 不写"应该没问题"。

1. **AC#9 整格判不了**：尺要问的是 embed 之后文件系统里的条目数与关键入口文件，而那在两层禁令里（不许读、结论也不许出现那两层的行号）。本腿唯一能动的是 git 名册级尺，而名册级尺**恰好不足以**回答"有没有一包页面产物"。
2. **AC#4 的真机那一发**：要真窗＋页面那枚入口＋一次真点击。本腿没有开过任何窗（`cmd/wisp` 里那 9 枚新用例走的是 `dispatchRaw`/回执字符串那一层，尺＝测试名册）。
3. **AC#2／AC#8 的"页面真看得见"那一半**：`Handle` 返回串离进程之后发生什么，在禁令层里；本腿只能证到"装配根把那枚 reader 与那两枚 setter 接上了"（`run.go:431/:710`，现读）。
4. **AC#10 的 ⓘ／ⓑ 选形**：⛔ 不归裁决者。本腿连"移动会不会破 `Confirming` 那一维"都没量（那要跑票 224 r2 那套授权仪器，本腿没跑，前腿与实现腿也都没跑）。
5. **`staticcheck`**：本机版与 CI 钉版不同（实现件 §6-7 与前腿都未复认），派单明令不跑，本腿没跑。
6. **CI 上的终态**：本腿判的是这台机器的工作树。云端那发今天算不算红本腿不知道——而且 `docs/evidence/s1` 之外那枚 CI 台件（`probes(ci-ed459d09)`，`8da52ff1`）刚记下"test-core 日志整步归因到 UNKNOWN STEP ⇒ 逐名红名册根本采不到"这个老毛病仍在，所以"拿 CI 当第二把尺"今天不成立。
7. **`GRANT-DROPPED` 那两行的实跑证据**：票面 AC#10 要的是"若走 ⓘ 那一行必须不再出现"，这需要真跑常驻腿。本腿只读到常量与构造点（`resident_approval_windows.go:109-113` 缺三项），**没读到"用户按了本次会话内允许、今天不生效"的实跑一行**。
8. **时序类读数**：取数期间同机另有 `254-r1b` 在飞（写面 `scripts/*.sh`）。本腿所有毫秒级／超时相关读数一律按【可能被抢 CPU】记，⛔ 不据此判任何 SLO 格。

