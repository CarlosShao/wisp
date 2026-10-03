# 票 261 · 腿 261-a2 普查表（AC#2：`applyDefaults` 不进 map ⇒ 全仓 map 型结构里还有几枚字段带 `default` tag）

代理：`261-a2`（只读普查子代理，零产码、零 Go 命令，纯词面普查）
票面：`.scratch/wisp/issues/261-the-model-enabled-key-promises-removal-from-discovery-and-selection-but-no-production-code-reads-it-while-the-default-true-tag-is-inert-for-map-entries.md`（AC#2 那格）
前人料：`.scratch/wisp/probes/261/p1/verdict.md` §5 第 3 条

---

## §0 起手锚与起手名册

- 起手时刻：`2026-10-03T15:55:xx+08:00` 起手，落盘时刻 `2026-10-03T16:06:23+08:00`（`date -Iseconds` 自取）。
- 起手 HEAD：`6c7a12b656188405e7ab40d6bb7c47066ff55dd6`（`git log -1 --format=%H` 自取），分支 `dev`。
  ⛔ 不引用别人写的锚号，本枚自取。
- 起手名册：`git status --porcelain internal cmd tools` → **空**（rc=0，零行输出 ⇒ 我的三包工作树干净，
  无他人未提交件；工作树里确有 design/** 删除与 probes 残面，均在我射程外，与本腿无关，未读未碰）。
- 证据件：`.scratch/wisp/probes/261/a2/census.md`（本件；目录新建，只建不删）。
- 纪律复读：禁跑任何 Go 命令（全程未跑 `go test/build/vet/run/list`）；不碰 AC 勾框；不动
  `docs/**`／`internal/**`／`cmd/**` 一字；不读不引 `frontend/**` 与 `design/**`（未打开任何一个文件，
  本件结论里无其内容）；只 commit 本证据件（显式 pathspec）；不 push、不 amend/reset/rebase/stash。

## §1 尺与推导式

**尺**：枚举对象＝**Config 配置树里的一切 map 型字段**（键动态、值里带字段的那些），逐枚列
「字段全路径＋所在结构＋其条目值类型里带 `default:"..."` tag 的字段枚数＋缺键落什么＋被谁读＋判语」。
**普查域**＝会从 `config.toml` 解码进 `config.Config` 的 Go 结构（`internal/config/schema.go`）——
`applyDefaults` 是且仅是这条树的默认执法者，map 字段要么在 schema.go 里，就不属于本票尺的射程。

**推导式**（四步，全部工具化执行、行号自取）：

1. **真身行号**：`grep -n "reflect.Map" internal/config/defaults.go` →
   `:77`（`applyDefaults` 的 `case reflect.Map:` → 下一行 `:78` 逐字 `// leave nil (see NewDefaults)`，
   **不进 map**＝本票公理，复认成立）· `:148`（`normalizeZero`，只把空 map 归 nil、对 map 内 struct 递归清零，
   不发默认）· `:212`（`copyValue`，深拷贝，不发默认）。⇒ **默认执法者只有 `applyDefaults` 一枚**，
   且它对 map 零执法。`case reflect.Map:` 实际在 **:77**，注释在 **:78**——票面写 `:75-80` 是范围引用，
   真身两行就在里面，无漂移。
2. **map 字段全名册**：`grep -n "map\[" internal/config/*.go`，取**非 _test.go** 且是**结构体字段声明**
   （非函数内临时量、非包级索引表）的命中；逐枚回 schema.go 对行亲读确认。
3. **tag 全名册**：`grep -rln 'default:"' --include=*.go internal cmd tools scripts third_party build models plans`
   → 生产文件命中**只有 `internal/config/schema.go` 一枚**
   （`defaults.go` 命中是注释里的字面 `default:"..."` 提法；`internal/config/enabled_261_test.go`、
   `cmd/wisp/firstrun.go`/`firstrun_198r2_test.go`/`config_reload_223_test.go` 命中全是注释或测试代码，
   `cmd/wisp/firstrun.go:7` 那枚逐字是
   `// (NewDefaults, defaults.go:58 - the `default:"..."` tags in schema.go are the`——注释，非 tag）。
4. **tag→字段对位**：通读 `schema.go` 全文 607 行，把每个 `default:"..."` tag 归位到它所在的
   struct→字段；再判该字段**是否活在某个 map 条目的值类型里**（条目内＝本票射程）还是
   **段级路径可达**（`NewDefaults` 会走进去＝tag 有执法者，不在射程）。

**口径**：① 只认 `Config` 配置树（schema.go）；② tag 认 `default:"..."` 字面（`toml:"-"` 的
`Portable` 无 tag 不算；`PluginEntry` 无一枚 tag）；③「缺键落什么」按 go-toml 对手写条目的解码
语义＝**该字段类型的 Go 零值**（261-p1 读数⑤已实证 `enabled` 缺键落 `false`，复认）；④
「会不会被读出意义」逐枚给出**生产读者 file:line**，不许只按读码断言——但本腿是零 Go 命令的
词面普查，读者名册来自 grep＋Read 双验，行为级「读数」归 261-p1 那种带仪器的腿。

## §2 map 字段全名册（逐枚）

**map 型字段总枚数＝6**（schema.go 全域 grep 后逐枚回读确认；test 文件里的命中全部排除）。
其中**条目值类型带 `default` tag 的字段枚数＝3**，不带的 3 枚各自说明如下。

| # | map 字段全路径 | 所在结构（schema.go 行） | 条目值类型 | 条目内带 `default` tag 的字段 | 缺键时条目里落什么 | 生产读者（file:line，逐枚 grep 复认） | 判语：缺键零值会不会被当成有意义的值 |
|---|---|---|---|---|---|---|---|
| 1 | `LLM.Providers` | `LLMSection`（:419 `Providers map[string]Provider`） | `Provider`（:379-401） | **1**：`Billing`（:391，`default:"pay-per-token"`） | 段级字段照常：`Protocol`/`BaseURL`/`APIKeyRef` 落 `""`（protocol/base_url 另有 preset 兜底）、`Billing` 落 **`""`**、`PlanCreditTotalMicro` 0、`Compat.Loose`/`AllowMissingUsage` false、`RPM`/`TPM` 0、`Models` nil | `Billing`：validate.go:160（枚举校验，**空串放行**）；resolver.go:130 无 Billing 传递；`Loose`/`AllowMissingUsage`：resolver.go:69-70 → 三适配器 `loose()`（openaichat/adapter.go:253/267 等）；`ExtraHeaders`：三适配器 range（openaichat/adapter.go:91 等）；`RPM`/`TPM`：resolver.go:130-131→:276 限速器 | **一半中雷**：`Compat.Loose` 缺键落 `false`＝严格形，语义恰好无害；`Billing` 落 `""` 会被 validate 的「空＝继承 provider/无要求」口径读成**「未显式声明」而非「pay-per-token」**——但段级 `Provider` 条目是用户显式手写的整表（写它必然逐键可见），且 261-p1 §5 钉的是 **ModelSpec.Billing**（见 #2）。⇒ 形状在、风险被「手写 provider 必然整表出现」部分抵消；**判语：中雷但弱于 #2** |
| 2 | `LLM.Providers.<name>.Models` | `Provider`（:400 `Models map[string]ModelSpec`）——**map 套 map 的第二层** | `ModelSpec`（:348-372） | **2**：`Billing`（:361，`default:"pay-per-token"`）＋ `Enabled`（:371，`default:"true"`） | 手写条目整枚走 Go 零值：`Enabled` 落 **`false`**、`Billing` 落 **`""`**、`Capabilities` 全 false、`ContextWindow` 0、`Price` 全 0、quota 0 | `Enabled`：**生产读者零枚**（261-p1 §1 读数①②③四向现跑＋本腿 grep 复认）；`Billing`：validate.go:197（**`m.Billing != ""` 才校验，空串放行**）；条目存在性：catalog.go:66/146、settings.go:115/153/302；整枚消费：resolver.go:119/287 | **两枚都中雷，本腿判语最狠的一格在此**：`Enabled` false 被读成「用户主动关了」而承诺句（schema.go:369-370「removes … from discovery/selection」）无人执行＝票 261 本体已判死形状；`Billing` 落 `""` 被 validate.go:197 的 `!= ""` 口径**静默放行**、计费侧也无 enum 兜底——tag 说「该字段总该是 pay-per-token 或 plan 二选一」，实际缺键落第三态 `""` 且**没有任何读者把它当默认值**。两枚都是「default tag 对 map 条目零执法＋零值被读出意义」的现行犯 |
| 3 | `LLM.Providers.<name>.Compat.ExtraHeaders` | `Compat`（:317 `ExtraHeaders map[string]string`） | `string` | **0**：值类型是裸 `string`，Go 的 string 无字段无 tag | 缺键＝map nil，range 零次 | 三适配器 range（openaichat/adapter.go:91、openairesponses/adapter.go:127、anthropic/adapter.go:158） | **不中雷**：nil map range 零次＝「不附加头」，正是缺省语义；string 值本身无零值歧义 |
| 4 | `Models.LocalOverride` | `ModelsSection`（:590 `LocalOverride map[string]string`） | `string` | **0**：同上，裸 `string` | 缺键＝map nil | config 层零读者（grep `LocalOverride` 于 internal/config 只有 schema.go 声明＋boundary_test）；真实读者在 `cmd/wisp/models.go:164` `overrides = cfg.Models.LocalOverride` → `:172` 注入 `internal/models/downloader.go:162/234`（`if dir, ok := m.opts.LocalOverride[id]; ok`——comma-ok，nil map ok=false） | **不中雷**：nil map 的 comma-ok 恒 false＝「无本地覆盖、走正常下载/校验」，正是缺省语义；且 downloader 侧用 `ok` 而非零值判读 |
| 5 | `Plugins.Entries` | `PluginsSection`（:575 `Entries map[string]PluginEntry`） | `PluginEntry`（:554-563） | **0**：`Enabled`（:556）**无** `default` tag——这段是写侧诚实形：`Enabled` 缺键落 `false`＝「不加载」，**与注释语义同向**（"Enabled loads the plugin at all"），不靠 tag、直接零值即正确 | 手写条目落 `Enabled:false`＝不加载、其余 slice nil | `Tier2Enabled`/`Entries` 走 parse.go:93-134 `buildPlugins`（**全仓唯一一条不靠 `applyDefaults` 的解码支路**——go-toml 抓不进 struct 才手工吸收）；读者：manager.go:506-522（diff 走时对比增删）；`Entries` 本身配置层之外读者零枚（grep 命中只有 test＋parse/writeguard 自身） | **不中雷**：零值方向与语义恰好同向（缺＝关）。⚠ 注意它**不是**「tag 不执法」案——是「没写 tag 所以没这个坑」；但**同构陷阱在其上方**：若日后有人给 `PluginEntry.Enabled` 加 `default:"true"`，本条就是下一枚票 261（Entries 由 parse.go 手工解码，`applyDefaults` 对 PluginsSection 整支照旧只进 struct 不进 map——Plugins 字段 `toml:"-"`，`applyDefaults` 仍会走到 Tier2Enabled 但 Entries map 依旧零执法） |
| 6 | `providerPresets` / `migrations` / `TierRegistry` / `lockedSections` / `lockedKeyDisposition` / `pluginsAbsorberType` | defaults.go:22 / migrate.go:35 / tiers.go:24、:82 / unwired.go:119 / parse.go:31 | 包级索引表 | — | — | 各自包内 | **射程外**：这些是包级 `var` 初始化的代码常量表，不进 `Config`、不经 go-toml 解码、无 tag 参与，列在此处只为堵「还有没有别的 map」的口径缺口 |

**汇总读数**：
- map 型结构字段（进 Config、键动态）：**5 枚真身**（#1、#2、#3、#4、#5）＋包级表 6 处（射程外，#6）。
- 条目值类型带 `default` tag 的字段：**3 枚**（`Provider.Billing`、`ModelSpec.Billing`、`ModelSpec.Enabled`）。
- 其中缺键零值**会被读出意义**的：**3 枚全中**——
  - `ModelSpec.Enabled`→`false`：被读成「用户主动关了」而**零读者执行**（票 261 已判死形状，复认成立）；
  - `ModelSpec.Billing`→`""`：被 validate.go:197 `!= ""` 口径**静默放行**（261-p1 §5-3 钉死，本腿复认）；
  - `Provider.Billing`→`""`：同上形状（validate.go:160），被「手写 provider 必然整表出现」部分抵消，风险次之。
- 不中雷的：`Compat.Loose`/`AllowMissingUsage`（tag 有，但那是**段级 struct** `Compat`——它作为
  `Provider.Compat` 字段被 `applyDefaults` 的 `reflect.Struct` 分支（:75-76）正常走进，tag **有执法者**；
  澄清：Compat 本身不是 map，只是恰好住在 map 条目里，其默认值**正常生效**，不计入中雷）。

## §3 推翻清单（票面与 261-p1 verdict 的每一句都是待验断言）

逐句亲验（grep＋Read 双做），结果：

| 待验断言 | 出处 | 结果 |
|---|---|---|
| `applyDefaults` 约 :75-80 的 `case reflect.Map:` 逐字只有一句注释 | 票面 AC#2 | **成立，行号精化**：`case reflect.Map:` 在 **:77**，注释 `// leave nil (see NewDefaults)` 在 **:78**，落在票面范围内 |
| 真身行号请自跑 `grep -n "reflect.Map"` | 任务书 | **已跑**：`:77 / :148 / :212` 三处命中，其中只有 `:77` 是默认执法支路 |
| `ModelSpec.Billing` 缺键落 `""` 且 `validateModelSpec` 放行 | 261-p1 verdict §5-3 | **成立**：schema.go:361 tag 在、validate.go:197 逐字 `if m.Billing != "" && !slices.Contains(...)` ⇒ 空串不进校验分支；本腿另补：**这枚形状不止 ModelSpec 一枚**——`Provider.Billing`（:391）同一形状（validate.go:160 同款 `!= ""` 口径） |
| 「全仓 map 型结构里还有几枚字段带 `default` tag」 | 票面 AC#2 尺 | **读数＝3 枚**（见 §2 汇总）；票面只点了 `enabled` 一枚当例，实际三枚 |
| 票 180 AC#4 的尺「数的是段级字段」、两层不在射程 | 票面 §为什么单独成票 | **成立**：`grep 'default:"'` 全仓生产命中只有 schema.go；其中段级 tag（App/Ball/…各 section）经 `applyDefaults` struct 分支正常执法，只有活在 map 条目值类型里的 3 枚落入本票射程——票 180 的尺没问这一维，属实 |
| `default:"true"` 对 map 条目在解码期不生效 | 票面标题 | **成立**（公理层）：defaults.go:77-78 不进 map；本腿为词面普查、未跑解码实验，行为级实证归 261-p1（其 §1 读数⑤ `TestTicket261P1MissingEnabledKeyDecodesFalse` PASS），两腿读数一致 |
| （无断言被推翻） | — | **推翻枚数＝0**；行号级精化 2 处（:75-80→:77-78；p1 未补的 `Provider.Billing` 同形状） |

## §4 量不到的格子（具名，不含糊）

1. **行为级解码读数**：本腿禁跑 Go 命令，「缺键落零值」的每一条都基于 schema.go 类型零值＋
   go-toml 严格解码语义＋261-p1 已落的两枚实证尺（enabled/Billing 形状）。`Provider.Billing`、
   `Compat.Loose` 等枚没有专属运行时尺——形状与 enabled 同构（同一 `applyDefaults`、同一解码支路），
   但**没有逐枚实测**；若编排者要运行时铁证，得另派带仪器的腿。
2. **`Provider.Billing` 落 `""` 的下游后果只算到 validate 层**：resolver.go 传递链里我没看到 Billing
   进 endpoint options（grep 于 resolver.go 的 `Billing` 仅 discover.go 两处构造、无传递），计费侧
   `internal/agent` grep `Billing` 零命中 ⇒ 「provider 级 billing 今天被谁消费」量不动——它可能也是
   一枚「有 tag 无读者」，但判死需要消费链全图，不属本格。
3. **`Plugins.Entries` 的未来形**：若日后给 `PluginEntry` 加带 tag 的字段，坑会在 parse.go 手工解码
   那支复现（它绕过 `applyDefaults` 整个 struct 树遍历？不——`applyDefaults` 从 `Config` 根走会进
   `PluginsSection` struct，但 `Entries` 是 map 依旧 :77 分支留 nil；手工支路也只解码不施默认）。
   这个「若加 tag 即中招」是推演，不是量得的事实。
4. **`toml:"-"` 的 `Portable`** 与 SchemaVersion：无 tag 无 map 参与，未列册；若有人日后改 `toml` tag
   形状，本册需重抽。
5. **第三方/生成码**：`third_party/`、`build/`、`scripts/spike` 的 `.go` 文件 grep `default:"` 零命中
   （sweep rc=0），但若存在无 `default:` 字面而由别的机制（如 go-arg 之类）施默认的库，不在本尺。

## §5 自查

- [x] 只读普查：`internal/**`、`cmd/**`、`docs/**` 一字未动（起手名册空，终态须复归空）。
- [x] 全程零 Go 命令（`go test/build/vet/run/list` 未跑；只用了 `grep`/`ls`/`git log/status`/`date`/`mkdir`）。
- [x] 未读未引 `frontend/**`、`design/**`（grep 路径均限定 internal/cmd/tools 等，通读文件仅
      schema.go/defaults.go/parse.go/verdict.md/票面）。
- [x] AC 勾框一枚未碰（本件只落 `.scratch/wisp/probes/261/a2/census.md`）。
- [x] 无 API 密钥值出现在本件任何位置。
- [x] 锚号自取（`6c7a12b6`），未引他人；行号均本腿现跑现读。
- [x] §0-§5 无一处「（待填）」。
- [x] 推翻清单逐句验完：推翻 0、精化 2（:75-80→:77-78；Provider.Billing 与 ModelSpec.Billing 同形状）。
- [x] commit 用显式 pathspec，只提本证据件，不 push。
