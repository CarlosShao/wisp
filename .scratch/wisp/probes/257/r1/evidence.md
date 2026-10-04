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

**三枚 tag 常量（`settings.go` 新增，字字不同）**：
- `refusalFileMissing = "第 1 种拒因：文件没建"`
- `refusalRowMissing  = "第 2 种拒因：行不存在"`
- `refusalInvalid     = "第 3 种拒因：校验不过"`

**三个拒点各挂各的**：缺文件→`loader.go` `readConfigFile`（`:76-81` 新增分支，句尾接第 1 种 tag＋"运行一次 wisp run 写出首份配置"的指引）；行不存在→`unknownProviderErr`／`requireCatalogEntry`（第 2 种 tag＋指引句）＋`requireChatProvider`（新门：config 没点名 chat provider 时，`role_chat_model` 的拒因是第 2 种而不是第 3 种）；校验不过→`writeOneKey` 的预写 `validate()` 门（第 3 种 tag＋"要改的是值；文件一个字节都没动"）。

**正控三发**（`go test -count=1 -run 'TestTicket257R1AC2' -v ./internal/config/`，19:06 全绿）——每发都是"种下该因 → 断言本句在场＋另两句**不在场**＋折叠句`配置未生效`不在场"：

| 发 | 测试 | 种法 | 红句逐字（断言的目标 tag） |
|---|---|---|---|
| 1 | `TestTicket257R1AC2MissingFileIsItsOwnSentence` | 临时目录里**根本不建** config.toml，直接 `NewManager` | `第 1 种拒因：文件没建`（并钉 `errors.Is(err, fs.ErrNotExist)` 仍在链上——`cmd/wisp` 的 `cause=missing` 分类器靠它，票 223 AC#4 的钉不许被吃掉；另钉 `LoadFile` 同句） |
| 2 | `TestTicket257R1AC2MissingRowIsItsOwnSentence` | 文件在、行不在：干净默认配置上 `SetProviderBaseURL`（provider 行缺）＋`SetModelContextWindow`（catalog 行缺）＋已建行但模型缺（`settingsBase` 加载后写 `ghost-model`） | `第 2 种拒因：行不存在`（并钉：拒句带 `guidanceProviderRow` 指引、**不含 `[providers.`**（A543 反控）、ghost-model 那发仍带票 261 的 `refusing to invent one` 原句钉） |
| 3 | `TestTicket257R1AC2InvalidValueIsItsOwnSentence` | 行在、值坏：`settingsBase` 上 `SetProviderAPIKeyRef("no-reference-prefix-here")`＋`SetRoleChatModel("not-in-the-catalog-9")`（此时 provider 已点名 ⇒ 走 `validate()` 正门） | `第 3 种拒因：校验不过`（并钉文件零变更） |

**"不许折成一句"的牙齿**：`whichReason()` 数 msg 里的 tag——**0 枚或 ≥2 枚都判不合格**（一句含两枚 tag 也是折）；`assertOneReason()` 正控"本 tag 在、另两 tag 不在"；AC#1 前半还逐枚断言折叠句 `配置未生效`／`生效失败` 不出现。第 3 发特意选 provider **已点名**的 config，证明 `role_chat_model` 的第 3 种分支没被 `requireChatProvider` 新门堵死——门只在"没点名 provider"时改判第 2 种。

## 5. 落地改动清单（本节答：动了哪几枚文件、每枚为什么动、`git diff --numstat` 三列）

改动全部落在 `internal/config`（授权写面）＋证据件。`git diff --numstat 0c303ef3` 对三枚产码/测试件（19:1x 现量）：

| 文件 | 加 | 删 | 为什么动 |
|---|---|---|---|
| `internal/config/loader.go` | +14 | -0 | AC#2 第 1 种拒因：`readConfigFile` 的 `os.ReadFile` 缺文件分支给出自己的那句（第 1 种 tag＋首启指引），原始 err 保持链上（`errors.Is(fs.ErrNotExist)` 不吃掉，`cmd/wisp` `cause=missing` 分类不漂） |
| `internal/config/settings.go` | +80 | -5 | ① AC#2 三枚 tag 常量＋三句指引常量（`[llm.providers.<名>]` 拼法，A543）；② `writeOneKey` 校验门挂第 3 种句；③ `unknownProviderErr`／`requireCatalogEntry` 挂第 2 种句；④ 新增 `requireChatProvider`：`SetRoleChatModel` 在 config 未点名 provider 时改判第 2 种（形 ⓒ 的"行不存在"诚实形状） |
| `internal/config/settings_257_test.go`（新） | +438 | -0（本树现量） | AC#1 前后半＋AC#2 三发的仪器；含起步态、nil 注册表、名册七枚行走、手加三样解锁、A543 反控、凭据面引用非值边界钉 |

（证据件 `.scratch/wisp/probes/257/r1/evidence.md` 与 msg-*.txt 为过程件，逐节 commit。）

## 6. 突变自证（本节答：每一发判据的"种下必红＋还原"，还原一律用突变前 `git cat-file blob HEAD:<path>` 抽的副本，附三枚 md5：起手＝还原后＝HEAD blob）

本腿继承的产码已 commit（`8e443ed7` 收编 257-r1 的 `loader.go`/`settings.go`/`settings_257_test.go`）⇒ 突变自证对**本腿 commit 后的 blob** 做三式：起手 md5（工作树）→ 种下突变（断言必红）→ 用 `git cat-file blob` 副本还原 → 还原后 md5 与 HEAD blob md5 三枚全等。`git cat-file blob` 前缀＝`8e443ed7`。

| 发 | 突变 | 预期红 | 实测 |
|---|---|---|---|
| M1 | `settings.go` 里 `refusalFileMissing` 值改为 `"第 1 种拒因：文件没BUILD"` | （见下）红未兑现＝同包常量同步，判语见后 | 红**未兑现**（具名判语在后） |
| M2 | `settings.go` 里删 `requireChatProvider` 挂钩（`SetRoleChatModel` 的 check 回 `nil`） | AC#1 前半红（role_chat_model 不再报第 2 种） | 红：`role_chat_model: refusal names none of the three reasons` |
| M3 | `settings_257_test.go` 手加段改回"追加 `[llm.roles.chat]`"（前任死法） | 发二红（duplicate table） | 红：`the hand-added shape the guidance teaches must load: config: config.toml parse: toml: table chat already exists` |

**md5 三枚全等**（每发还原后实测，`git cat-file blob` 前缀＝`84ab5bf2`＝继承码收编 commit）：

```
settings.go          工作树＝6d2166b9275ca0db45ae30700a5f0cbe ＝ blob(84ab5bf2:internal/config/settings.go) 同值   （M1/M2 还原后实测）
settings_257_test.go 工作树＝3a223edd967666ab26aa92b0fe4e1b0c ＝ blob(84ab5bf2:internal/config/settings_257_test.go) 同值（M3 还原后实测）
```

**三发实测记录（19:09–19:17）**：

- **M1（tag 值改写）——红未兑现，判"该发判据对此突变不敏感"，具名不入正控清单**：把 `refusalFileMissing` 值改为 `"第 1 种拒因：文件没BUILD"` 后，`TestTicket257R1AC2MissingFileIsItsOwnSentence` 仍 PASS。根因（我先用探针测试件在包内现量了 err 文本与常量）：测试与产码**同包**（`package config`），断言用的 `refusalFileMissing` 是**同一个 Go 常量标识符**——产码里的 tag 变了，断言引用的也跟着变，二者永远同步。这类突变只能由**异包复制字面量**的判据（如 `cmd/wisp` 侧的文案钉或 grep 二进制）抓到。这不是测试虚假：它钉的是"三句互斥＋tag 逐字在场＋ErrNotExist 哨兵存活"，对"三句折成一句"这一真实缺陷形态有牙（M2 的红就是它咬的），对"常量值整体改名"无牙是同包测试的固有边界。具名登记，不补假红。**还原**：`git cat-file blob` 副本覆写，md5 `6d2166b9…` 三枚全等，`go test -count=1 -run TestTicket257` 回绿。
- **M2（摘掉 `requireChatProvider` 挂钩）——红兑现**：`SetRoleChatModel` 的 check 回 `nil` 后，`TestTicket257R1AC1CleanMachineRefusesAllSevenFields` FAIL，红句逐字（19:1x）：
  `role_chat_model: refusal does not name [llm.providers.<name>], so nothing tells the operator how to make the row writable: config: refused to write llm.roles.chat.model; …第 3 种拒因：校验不过…: config.toml: llm.roles.chat.model set without llm.roles.chat.provider (set both or neither)`；分发计数变 `map[第 2 种拒因：行不存在:6 第 3 种拒因：校验不过:1]`——正是 AC#2 要防的"role_chat_model 被折进校验句"形状。**还原**：blob 副本覆写，md5 全等，回绿。
- **M3（手加段改回"追加 `[llm.roles.chat]`"，前任死法）——红兑现**：`TestTicket257R1AC1HandAddedRowsUnlockAllSevenFields` FAIL，红句逐字：`the hand-added shape the guidance teaches must load: config: config.toml parse: toml: table chat already exists`。**还原**：blob 副本覆写，md5 `3a223edd…` 全等，回绿。

三发终态：工作树两文件 md5 与 HEAD blob 全等、`git status --porcelain -- internal/config` 干净、全 ticket 测试 PASS。

## 7. 门禁读数（本节答：d22scan／path-length-budget／gofumpt／go vet／`go test -count=1 -v ./internal/config/` 五门，带时刻，红名集合逐名比对并写清哪几枚既有）

未判。

## 8. 判不动／量不到（本节答：具名＋归口，不许"应该没问题"填空；含 AC#1 里属于 `cmd/wisp` 的那一格我为什么够不着）

未判。

## 9. 交件判语（本节答：Git 纪律／AC 框未碰／冻结件未碰／凭据面未动／禁读面未读，逐条）

未判。
