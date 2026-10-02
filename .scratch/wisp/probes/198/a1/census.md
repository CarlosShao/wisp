# 票 198 · 腿 198-a1 · 只读普查（首建 `config.toml`）

> 本文件是**料**，不是产码。§1–§5 给落地腿当派单名册用。所有行号按 §0 锚点自取实测。

## §0 起手锚

| 项 | 读数 | 尺 |
|---|---|---|
| HEAD | `6fe300c5035a60586e681f586a4ead5e2b56137b`（`census(rate-census-3a)：编号 150-199 的 38 枚开放票分桶…`，提交时刻 `2026-10-02 08:41:33 +0800`） | `git log -1 --format="%H%n%ad%n%s" --date=iso` |
| 分支 | `dev` | `git rev-parse --abbrev-ref HEAD` |
| 本腿起手时刻 | `2026-10-02 08:45:58 +0800` | `date "+%Y-%m-%d %H:%M:%S %z"` |
| 票面 | `.scratch/wisp/issues/198-a-fresh-machine-cannot-run-wisp-because-nothing-creates-config-toml.md` | **实测 35 行**，非编排者题面写的 48 行（见 §7-X1） |

`git status --porcelain cmd internal` 起手读数（**非空属正常：那是别枚在飞的腿，本腿一字未动未提交**）：

```
 M internal/panel/bridge.go
 M internal/panel/composer_dispatch.go
 M internal/panel/composer_dispatch_test.go
 M internal/panel/l2_grant_boundary_test.go
?? internal/ball/sta_release_windows_test.go
?? internal/config/settings.go
?? internal/panel/config_handlers.go
```

⚠ 与本票直接相关的一枚：`internal/config/settings.go` **此刻是未跟踪文件（313 行）**，内容尚未进任何 commit。
本腿读过它（见 §1/§4），但它**随时可能被那枚腿改写或删除**，所以凡本文件引它处均标 `[untracked-197]`，
落地腿开工时必须重新取一次行号，不得照抄本文件的数字。

## §6 待人裁 / 判不动的地方

- **J1｜首建该长在哪一层：`assembleRuntime` 之内还是 `runTextTask` 之前。**
  代码事实：生产侧调用 `assembleRuntime` 的只有两处——`cmd/wisp/run.go:234`（`wisp run`）与
  `cmd/wisp/resident_task_source_windows.go:265`（常驻腿取任务）。退码 2 那句 `wisp run: 配置未就绪
  （Unconfigured）` 在 `cmd/wisp/run.go:397-398`，即首建若插在 `config.NewManager` 之前，
  **两枚宿主同时被改**；插在 `runTextTask` 则只改 `wisp run`，常驻腿继续今天这形。
  我判不动的是**授权范围**：票面 §"要建什么" 1 只说"找不到配置时"，没点名宿主枚数，而"全新机器第一天能跑"
  这句话对常驻腿成不成立不归我判。需要**编排者裁**（裁"只 `wisp run`"还是"两枚宿主一起"），
  因为这决定 AC#1 要不要在常驻腿上重复一遍正控。
- **J2｜AC#1 那句"带全部 section 的默认 `config.toml`"与 AC#2"不许新造默认值"今天互相拉扯。**
  代码事实：`NewDefaults()`（`internal/config/defaults.go:58-62`）**刻意把所有 map 字段留 nil**
  （`defaults.go:77-78` 注释逐字：`leave nil (see NewDefaults)`），于是
  `llm.providers` / `models.local_override` / `plugins.Entries` 这几族**没有条目可写**——
  写出来的文件里这些表要么不出现、要么是空表。"全部 section"若按字面理解成"每张表都在文件里"，
  落地腿就得**自己填条目**，那正好撞 AC#2 的"不许新造默认值"。我判不动哪种读法作数，
  需要**编排者或 owner 裁**（我倾向"18 枚静态 section 全在、动态表按 nil 留空"这一形，理由见 §1-§1.4）。
- **J3｜`SaveFile(NewDefaults())` 落盘后文件里到底出现几行、哪些表在不在，本腿量不了。**
  这是"必须真跑才答得出"的一格。⛔ 本腿禁 `go test`/`go build`/`go run`（两枚写码腿在跑毫秒级计时）。
  交**待落地腿自量**，确切命令与期望读数：
  `go test ./internal/config -run TestUnwiredGuardLeavesHonestConfigsAlone -v`
  → 期望 PASS（`internal/config/unwired_test.go:86-95` 这一子用例今天已经就是
  "SaveFile(path, NewDefaults()) 之后再 LoadFile 要能读回"的钉，它绿 ⇒ 首建机制的可加载性已被现成用例钉住）；
  量文件形状请另写一次性探针（`t.TempDir` + `SaveFile` + `os.ReadFile`），期望
  静态 18 section 各出现一次、`[llm.providers.*]` 与 `[plugins.<id>]` 零次、`schema_version = 2` 一行。
- **J4｜首建该不该同时把 `[llm].text_chain` 或某个 provider 名字写进去。**
  写了就是"拿一个假端点凑"，票面 §3 明令不许；不写则首建之后 `wisp run` **仍然退码 2**
  （撞点是 `cmd/wisp/run.go:409-413`，`llm.Resolver.ResolveRole` 在
  `internal/llm/resolver.go:218` 报 `role %q is unset and text_chain is empty`）。
  也就是说 **AC#4 才是本票真正的用户可见终点，首建本身不改变退码**。这条我判得动（结论见 §3/§5-§5.3），
  但"首建之后还缺什么那句话该由谁打印"（`assembleRuntime` 里还是 `runTextTask` 收口处）需要编排者裁。
- **J5｜票 132 `SealDir` 一族：AC#3 说的"沿用已定的那一形"今天生产侧根本没人走 `SealDir`。**
  事实见 §2-§2.3（显式生产调用者 **0 枚**，能力全部由 `PrivateDirAll`/`SealFile` 投递）。
  我判不动的是**验收怎么算过**：落地腿若用 `winsec.PrivateDirAll` 建目录，AC#3 那句"票 132 的 `SealDir` 一族"
  是按"符号名命中"判还是按"同一判定者投递"判？前者会把正确实现判错。需要**编排者裁**，
  验收腿照裁过的口径读。
- **J6｜`cmd/wisp/secret.go` 的 `configFileName` 常量与首建文案的归属。**
  `configFileName = "config.toml"` 定义在 `cmd/wisp/secret.go:63`（不在 run.go、不在 config 包），
  首建路径今天唯一拼它的地方是 `cmd/wisp/run.go:371`。新建文件若复用该常量则跨文件依赖 secret.go，
  若新写一份则造出第二处真相。这属"改哪枚文件"的派单细节，我判不动要不要顺手搬家，
  需要**编排者裁**（本腿倾向：复用现常量，搬家不在本票）。

## §7 本腿推翻编排者题面之处

题面给的 4 处现量 + 那把 grep，逐枚复认（判定依据在 §1/§2/§4，此处只记结论）：

| # | 题面（与票面同源） | 判定 | HEAD 实测 |
|---|---|---|---|
| X1 | "票 198 …… **48 行**" | **推翻** | **35 行**（`wc -l .scratch/wisp/issues/198-*.md` = 35）。编号存在、文件名一致，仅枚数错 |
| X2 | `internal/config/schema.go:104-134` ＝ **16** 个 section | **行号成立、枚数推翻** | `104-134` 正是 `type Config struct` 全体；但区间内 section 字段实测 **18 枚**（`schema.go:110-133`，含 `Plugins`）。仓内 `type *Section struct` 亦 **18 枚**。见 §1-§1.3 |
| X3 | `internal/config/manager.go:76-84` ＝ `NewManager` 注释"首建是调用方的事" | **结论成立、行号已漂移** | HEAD：注释在 `manager.go:99-100`，函数体 `99-111`；`76-84` 此刻是 `LockedDecision`。查证：`git show 335b8d2b:internal/config/manager.go` 里 `first-run file creation is` 确在 **76 行** ⇒ 票写时成立，属**过期**不是错写 |
| X4 | `cmd/wisp/run.go:282-291` ＝ 缺配置退码 2 | **行号错（已漂移）、结论成立** | HEAD 的 `282-291` 是 `agentRuntime` 的结构体字段（`endpoint`/`gate`/`ui`/`bridge`…），与退码无关。真链：`run.go:389` `config.NewManager(cfgPath, nil)` → `run.go:397-398` 打印 + `return rt, 2`；退码表在 `run.go:68-91`（`ClassConfig`/`ClassAuth` → 2）。同样查证：09-28 锚点 `335b8d2b` 的 `282-291` 落在 `assembleRuntime` 的注释上（当时近似成立） |
| X5 | `internal/config/loader.go:39-43` ＝ 缺配置时退码 2 | **半成立，指错了产生点** | `loader.go:37-41` 是 `LoadFile` 的文档注释（"the caller maps it to the Unconfigured state"），`:41-43` 是 `LoadFile` 签名。**"缺文件 ⇒ 错"真正被造出来的是 `loader.go:65-68`**（`os.ReadFile` 失败 → `observe.Wrap(observe.ClassConfig, err, "config.toml read")`）。loader 从不产生退码 |
| X6 | `cmd/wisp/secret.go:263-271` ＝ key 只能 CLI 手录 | **成立（HEAD 即为该区间）** | `263-271` 正是 `case "set"/"get"/"list"/"unset"` 的子命令分派；`configFileName` 常量在 `secret.go:63`，`api_key_ref = "dpapi:<name>"` 文案在 `secret.go:175`、`secret.go:380` |
| X7 | 编排者 08:4x 现跑 `grep -rn "config\.toml" cmd/ internal/ --include=*.go \| grep -v _test` ＝ **12 行命中** | **枚数推翻（差一个数量级）；结论成立** | 本腿逐字复跑同一把尺：**163 行命中**（口径：`cmd/`＋`internal/` 全树非 `_test.go`，含 `cmd/balldebug`）。"没有一处真创建配置文件"这句**方向对但说法要收窄**：见下 X8 |
| X8 | 票面"全仓**无一处创建 `config.toml`**"（`Status` 与现量表第 2 行） | **推翻（存在一处生产创建点，只是首达不到）** | `internal/config/writeguard.go:125-131` 有 `case fileMissing(statErr)` 分支，注释逐字写着"writing it creates what first-run did not"，并落到 `writeguard.go:160` 的 `SaveFile(m.path, base)` ⇒ **`Manager.mergeWrite` 今天能在文件缺失时把它建出来**。它之所以救不了全新机器，是因为 `NewManager`（`manager.go:101-105`）在装载失败时直接返回错误，生产根本走不到 `mergeWrite`。正确的说法是：**首建能力已存在且已导出，缺的只是入口那一次调用** |
| X9 | 题面"internal/config 今天有没有**任何**写盘能力"（暗示可能没有） | **推翻暗示：有，而且三枚都是现成的** | `SaveFile`（`loader.go:238`，**已导出**）· `MarshalCanonical`（`parse.go:164`，**已导出**）· `atomicWrite`（`parse.go:196`，未导出，内含 `winsec.SealFile` 于 `:215`）。落地腿**不需要新写序列化器**，因此题面担心的"默认值改一处要同步两处"这枚维护面**不会被引入**（论证见 §1-§1.5） |

> X3/X4 的"过期 vs 错写"分界是查出来的，不是猜的：`git show 335b8d2b:<file>` 的读数已附在上面两行。
