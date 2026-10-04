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

未判。

## 3. AC#1 干净机器那一发**真跑**（本节答：临时数据根＋无 `config.toml` 起步、建完首份默认配置后逐枚试写那七枚，终态＝形 ⓒ 真兑现；原样命令＋逐枚读数；手加三样之后 7/7 是否真解锁）

未判。

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
