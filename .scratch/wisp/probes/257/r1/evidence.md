# 257-r1 落地件 — 票 257 形 ⓒ 的 `internal/config` 那一半（第一任产码腿）

> 腿：`257-r1`；票：`.scratch/wisp/issues/257-clean-machine-provider-registry-nil-blocks-writes.md`。
> 前件：`257-a1/census.md`（92 行，`ed3fd270`）＋`257-a2/census.md`（257 行，`cbf4f1b2`）＋台账 `A543`／`A560`＋票面 §8。
> ⛔ 本件不复跑前人已经跑过的读数；只补我缺的那一发。

## 0. 起手锚与并发窗口（本节答：我在哪一枚树上开工、同机谁在飞）

| 项 | 读数（本腿现量） |
|---|---|
| 进场时刻 | `2026-10-04T16:11:17+08:00` |
| 进场 HEAD | `40aae961`（branch `dev`） |
| `git status --porcelain` 总行数 | **632**（同机多枚在飞的正常量级；`257-a2` 进场 398、落件 420，本腿比它多 212 枚） |
| `git status --porcelain -- internal/config internal/panel` | **0 行**＝这两包此刻**没有别人的活**，本腿可开工（派单第 1–3 轮的停手闸门已过） |
| 同机在飞（派单写死的包级互斥） | `265-r1`＝`cmd/wisp/**`（⇒ `cmd/wisp` 本腿不碰、不归因）／`174-r4`＝`internal/tools/**`＋`internal/agent/spill.go`（⇒ 那两包不碰、不归因）／`265-a1b`＝只读普查腿 |
| 本腿写面（编排者授权） | `internal/config/**`（＋确有必要才 `internal/panel/**`，须先具名理由） |
| ⚠ 已知既有脏（不归本腿、本腿不修） | `gofumpt -l` 点名 `cmd/wisp/models.go`／`pending_read.go`／`queue.go`＝票 212／258 账户；`internal/panel` 的 `TestC21DesignTokensFourWayAgree` 与 `internal/ball` 同族（design 资产删除所致）；`cmd/wisp` 无 sherpa PATH 时 `0xc0000135` 且没有 `--- FAIL` 行＝根本没跑 |
| 本腿所有读数的时间窗 | 骨架发＝`16:11`；此后逐节就地追加，**行号只对"该节写下的那一刻的共享树"负责**（本腿不 checkout、不开 worktree） |
| ^ 接续腿 `257-r1b` 开工追加（2026-10-04T19:0x +08） | 进场 HEAD `0c303ef3`（branch `dev`）；`git status --porcelain -- internal/config`＝**3 行**：`M internal/config/loader.go`（+14）／`M internal/config/settings.go`（+80）／`?? internal/config/settings_257_test.go`（14,964 字节）——全部是死去的 `257-r1` 的未验证产码，**继承，未做任何形状改写**；`defaults.go` 不在脏清单 |

## 1. 写面授权与形 ⓒ 的落点差（本节答：编排者给我的写面＝`internal/config`，票面 §8-4 预告的写面＝`cmd/wisp`，这一差我必须具名，并把"哪一格落在我这一面"划清）

**两个写面的归属，逐格划清**（票 §8-4 预告写面 `cmd/wisp`，编排者给我的写面＝`internal/config`——预告作废以派单为准）：

| 落点格 | 归谁 | 依据 |
|---|---|---|
| 拒因三句（`loader.go` 缺文件句＋`settings.go` 行不存在／校验不过句）＋`requireChatProvider` 新门 | **我这一面**（`internal/config`，继承自 257-r1 已落盘） | 派单写明"写面＝`internal/config`＝拒因三句＋干净机器仪器那半" |
| 干净机器仪器（`settings_257_test.go`，AC#1＋AC#2 的测试件） | **我这一面** | 同上 |
| `cmd/wisp/firstrun.go` 首启回执文案（对"没建过 config.toml"的机主说的那句） | **不归我，具名归口 `cmd/wisp` 那半** | 票 §8-3 ⓒ 边界①（三条通道各说各的）；写面被 `265-r1b` 占（`cmd/wisp/**` 包级互斥），我不碰 |
| `internal/panel`（C17 面） | **零格**：ⓑ 已毙 ⇒ 名册七枚一字不动 | 台账 `A543`；派单禁区明确 `internal/panel/**` 不碰 |

**AC#1 里属于 `cmd/wisp` 的格子**＝"机主第一次跑 `wisp run` 时控制台/回执告诉他先手加 `[llm.providers.<名>]` 三样"——那一发要真跑得有 firstrun 文案，我够不着（包级互斥），在 §8 具名。我这一面能兑现的 AC#1 终态＝"建完首份默认配置后逐枚试写七枚 → ⓒ 回执说清要手加哪三样；手加三样后 7/7 解锁"，这一发我用 Go 测试在 `internal/config` 包内真跑（§3）。

**形状服从性判定（派单第一件事，已全量读三枚 diff）**：继承码**符合 ⓒ／AC#2**，未触发停手闸——
- **ⓐ 无**：`defaults.go` 不在脏清单、diff 为零；首建文件仍无 `[llm.providers` 行（测试件 `cleanMachineManager` 自己钉了这一点，且与票 198 的钉同向）。
- **ⓑ 无**：diff 没有给名册（`internal/panel/config_handlers.go` 的七枚常量表）加任何字段，`internal/panel` 一字未动。
- **A543 拼法照办**：指引句写的是 `[llm.providers.<名>]`（`guidanceProviderRow`／`guidanceModelRow`／`guidanceRolePair` 三枚常量＋`unknownProviderErr` 的拼接），不是票面原拼法 `[providers.x]`；测试件里还有一枚反控钉死"指引句里不得出现 `[providers.`"（`TestTicket257R1AC2MissingRowIsItsOwnSentence` 的 `DisallowUnknownFields` 断言）。
- **AC#2 三句形状**：`refusalFileMissing`／`refusalRowMissing`／`refusalInvalid` 三枚 tag 常量字字不同，三处拒点各挂一枚，测试件 `assertOneReason` 正控"本句在场＋另两句不在场＋无折叠句'配置未生效'"。
- **附随核过**：继承码引用的符号全部真实存在——`fileMissing`（`parse.go:233-235`）、`observe` 导入（`loader.go:11` 原有）、测试件的 `settingsBase`／`readSettingsFile`（`settings_248_test.go:23/46` 同包）；`go build`＋`go vet ./internal/config/` 干净（19:0x）。

## 2. 四环现量复认（本节答：票面那五行待验断言，我这一枚树上是成还是废；含 `A543` 那条拼法更正的照办）

现量时刻 `2026-10-04T19:0x`，HEAD 侧 `git show HEAD:<path>`＋脏侧 sed 双向读。票面行号是对立票锚 `cf87f7a2` 的树说的，本树漂移逐条具名：

| # | 票面断言 | 本树现量 | 判 |
|---|---|---|---|
| 1 | `defaults.go:77-78` 对 map「leave nil」 | 逐字在：`defaults.go:78`＝`// leave nil (see NewDefaults)`，`case reflect.Map:` 下一行；且我现跑 `NewDefaults()` 探针（19:0x）＝`Providers map: map[string]config.Provider(nil)` | **成** |
| 2 | `schema.go:419` 无 `default` 标签 | `schema.go:419`＝`Providers map[string]Provider \`toml:"providers"\``，无 default 标签；对照 `TimeoutMS`（:415）带 `default:"60000"` | **成** |
| 3 | `settings.go:310-314`／`:294-308` 对"行不存在"是拒不是建 | HEAD 行号：`:226`＝apply 返回 false 时拒（`"does not exist in the file this write is based on; nothing was written"`，写侧第四道门）；`:294`＝`requireCatalogEntry`（provider 行缺→`unknownProviderErr` `:310`；model 行缺→拒）；本树脏侧漂移至 `:286`／`:356`／`:378`（257-r1 增量所致）。全链**四处**拒点、零处建行 | **成**（行号差＝HEAD `:226/:294/:310`，票面 `:310-314/:294-308` 是"双层关"两处，都核在） |
| 4 | `role_chat_model` 撞 `catalog.go` model 有值 provider 空 ⇒ 拒 | HEAD `catalog.go:103`＝`"%s.model set without %s.provider (set both or neither)"`；257-a1 更正的 `:101-103` 与本树一致；另：干净机器上 chat pair 两空，写 model 后变"model 有 provider 无"→正是这枚拒 | **成** |
| 5 | 名册七枚无建行入口 | `internal/panel/config_handlers.go:56-70`＝七枚常量（`provider_base_url`／`provider_api_key_ref`／`model_context_window`／`model_price_in`／`model_price_out`／`role_chat_model`／`provider_credential`），`allWritableFields` 即此七枚；每枚落 `internal/config` 的 setter 都以"行已存在"为前提（§3 表） | **成** |

**`A543` 拼法更正照办**：`internal/config/parse.go:72`（`Decode` 后 strict 检查）与 `:123`（`DisallowUnknownFields` 所在段）现量核实真有 `DisallowUnknownFields()`；顶层 `[providers.x]` 是未知键 ⇒ 整链起不来。继承码的指引句（`settings.go` 新增 `guidanceProviderRow`／`guidanceModelRow`／`guidanceRolePair` 三常量）**全部写 `[llm.providers.<名>]`**；测试件反控钉死拒句不得含 `[providers.`。**照办**。

**附带一改（对本树唯一的产码测试件修正，具名理由）**：继承测试件 `TestTicket257R1AC1HandAddedRowsUnlockAllSevenFields` 的手加段按"追加一节 `[llm.roles.chat]`"写——现量证伪：首份默认文件由 `MarshalCanonical` 无 omitempty 全量展开，**本来就带 `[llm.roles.chat]`（provider/model 空串）**（我 19:01 用临时 main 程序 dump 过真文件逐行核对）；追加第二节＝`toml: table chat already exists`，测试死于 `settings_257_test.go:203`。真实机主动作＝**就地把已有节的空串点名**。已照实改写（`strings.Replace` 填空串对），另顺手消了一枚 `:=` 遮蔽编译错（`:241`）。这不是改形状：产码（`loader.go`／`settings.go`）一字未动，动的只是测试件里"机主怎么手加"的模拟方式，且改后更贴 A543 的"手加三样"。

## 3. AC#1 干净机器那一发**真跑**（本节答：临时数据根＋无 `config.toml` 起步、建完首份默认配置后逐枚试写那七枚，终态＝形 ⓒ 真兑现；原样命令＋逐枚读数；手加三样之后 7/7 是否真解锁）

**跑法**（`go test -count=1 -run 'TestTicket257' -v ./internal/config/`，两发共四次，19:06 末次全绿）：

**发一（前半）＝`TestTicket257R1AC1CleanMachineRefusesAllSevenFields`**，起步态严格按票面：
1. `t.TempDir()` 临时数据根，`os.Stat` 钉死**无 config.toml**（不在 ⇒ 才继续）；
2. 起步态自带一句验证：此时 `NewManager(path, nil)` 必须拒且拒句带 `refusalFileMissing`（第 1 种拒因）——证明"没建文件"就是这台机的起点；
3. `SaveFile(path, NewDefaults())` 建**首份默认配置**（与 `cmd/wisp/firstrun.go:82` 同一 SaveFile＋同一 NewDefaults，非手工塞成品）；
4. 钉文件里**无 `[llm.providers`**（票 198 的钉同向：首建文件不发明 provider 行）；`NewManager` 后钉 `m.Config().LLM.Providers == nil`（票面前提：nil 不是空 map）；
5. 逐枚试写七枚（`writableRosterWalk()`，名册七枚各一枚 setter，`provider_credential` 走 `SetProviderAPIKeyRef`＝StoreCredential 的 config 侧那半，AC#3 只碰引用不碰值）。

**读数（19:06 终跑）**：`seven refusals by reason: map[第 2 种拒因：行不存在:7] (fields: 7)`——七枚**全拒**、七枚**全部**报"行不存在"（没有一枚折成第 3 种"校验不过"，也没有一枚报第 1 种）；每句都带 `[llm.providers.` 指引（形 ⓒ 兑现）；文件零字节变更、内存 chat model 未被偷动。

**发二（后半）＝`TestTicket257R1AC1HandAddedRowsUnlockAllSevenFields`**：按回执指引**手加三样**（provider 行＋两条 model 子表带 `enabled = true`＋已有 `[llm.roles.chat]` 节内点名 provider/model）→ `NewManager` 加载过 → **七枚 7/7 全部接受**，每枚回执带键路径；文件实测落了 base_url／context_window／键路径；凭据那枚落的是 `api_key_ref = 'dpapi:t257r1canaryblob01'` 引用、全文件无 `api_key =` 明文键（AC#3 边界钉在测试里）。

**逐枚读写对应表**（七枚 → 键路径 → 干净机拒点）：

| 名册字段 | setter | 写的键路径 | 干净机上的拒点 |
|---|---|---|---|
| provider_base_url | `SetProviderBaseURL` | `llm.providers.<名>.base_url` | 行不存在（`unknownProviderErr`） |
| provider_api_key_ref | `SetProviderAPIKeyRef` | `llm.providers.<名>.api_key_ref` | 同上 |
| model_context_window | `SetModelContextWindow` | `llm.providers.<名>.models.<id>.context_window` | 行不存在（`requireCatalogEntry` 双层） |
| model_price_in / out | `setModelPriceLeaf`×2 | `...models.<id>.price.in`／`.out` | 同上 |
| role_chat_model | `SetRoleChatModel` | `llm.roles.chat.model` | 行不存在（`requireChatProvider` 新门：config 未点名 provider 时拒"行不存在"而非"校验不过"） |
| provider_credential | `SetProviderAPIKeyRef`（引用半） | `llm.providers.<名>.api_key_ref` | 行不存在 |

**终态判语**：形 ⓒ 在 `internal/config` 这一面**真兑现**——界面不建行，但每一句拒写都告诉机主去 `[llm.providers.<名>]` 手加哪三样；手加后 7/7 解锁。

## 4. AC#2 三种拒因各配一句（本节答：文件没建／行不存在／校验不过三句是否字字不同、各配一发**种下必红**的正控，三发不合成一句；红句逐字）

未判。

## 5. 落地改动清单（本节答：动了哪几枚文件、每枚为什么动、`git diff --numstat` 三列）

未判。

## 6. 突变自证（本节答：每一发判据的"种下必红＋还原"，还原一律用突变前 `git cat-file blob HEAD:<path>` 抽的副本，附三枚 md5：起手＝还原后＝HEAD blob）

未判。

## 7. 门禁读数（本节答：d22scan／path-length-budget／gofumpt／go vet／`go test -count=1 -v ./internal/config/` 五门，带时刻，红名集合逐名比对并写清哪几枚既有）

未判。

## 8. 判不动／量不到（本节答：具名＋归口，不许"应该没问题"填空；含 AC#1 里属于 `cmd/wisp` 的那一格我为什么够不着）

未判。

## 9. 交件判语（本节答：Git 纪律／AC 框未碰／冻结件未碰／凭据面未动／禁读面未读，逐条）

未判。
