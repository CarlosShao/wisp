# `309-v2` 突变与读数总表（四面进攻的全部台件）

约定：每发一把尺、一行读数；`rc=`＝`go test` 退码（⛔ `cmd | head; echo $?` 那种取错东西的写法）；
"净新增红"＝该发顶层 `--- FAIL` 枚数 − pristine 顶层 `--- FAIL` 枚数（pristine＝**1**＝
`TestTicket252R1AlignmentAddsNoAuthorization`，台面缺件，票 298 `AC#5` 已承接，⛔ 归本票）。
还原纪律：每发 `cp` 备份 → 突变 → 跑 → `cp` 还原 → `cmp -s` → `sha256sum`；A 组每形那行
`RESTORE cmp: byte-identical OK … match=YES` 逐形在场，终局整树 bridge.go sha
`995c23d2ef6ea16c3fd0b2a7e0b5c2c18ea87bb860f5ada1e18e52ebbe464eba` ＝
`git show f6067e5d:internal/tools/bridge.go | sha256sum` 同一枚（件 `logs/rigD-callsite.log` 末四行）。

## A 组＝`AC#1` 那枚仪器（树 A，突变体只 `internal/tools/bridge.go`；件 `logs/rigA-instrument.log`）

| 发 | 突变（逐字锚） | 尺 | 读数 | 判 |
|---|---|---|---|---|
| A0 | ⛔（pristine 正控） | `go test -run Test309AssessorForEmptyTaskIDDropsTaintDetector -v` | `--- PASS` `rc=0` | 该绿 ⇒ 绿 |
| A1 | `if b.injected \|\| b.prov == nil \|\| taskID == "" {` → 去 `\|\| taskID == ""` | 同上 | `--- FAIL: Test309…` `rc=1` | 恒真形**响** |
| A2 | `WithTaintDetector(b.prov.Detector(taskID))` → `WithTaintDetector(nil)` | 同上 | `--- FAIL: Test309…` `rc=1` | 反向（没挂⇒该红）**响** |
| A3 | 整枚守卫 off（`if false {`） | 同上 | `--- FAIL` `rc=1` | 票面"整枚摘掉守卫"那一支**响** |
| A4 | 肢取反（`taskID != ""`） | 同上 | `--- FAIL` `rc=1` | 换形也红（⛔ 只有"改⇒红"） |
| A5 | 摘 `b.injected \|\|` | 整包 `-v` | 顶层 FAIL＝**1**（＝pristine）、PASS＝**210**（＝pristine）、309 仪器 `--- PASS` | **净新增 0**＝这一肢⛔ 任何断言 |
| A6 | 摘 `b.prov == nil \|\|` | 整包 `-v` | 顶层 FAIL＝**17**、PASS＝**194** ⇒ 净新增 **16**；FAIL 名册里 `grep -icE 'prov\|assessorfor\|taint\|309'`＝**0**；309 仪器 `--- PASS` | 噪声红⛔ 具名该肢（`Q-309-r2-prov-limb-colour-discrepancy` 两半都实测复现） |
| A7 | ⛔（pristine 整包正控） | 整包 `-v` | 顶层 FAIL＝1／PASS＝210 | 基线，A5／A6 的净数就减它 |

## B 组＝`AC#2` 那枚枚数钉（树 B，尺＝`go test -run Test309ToolRequestProductionConstructionNail -v`；件 `logs/rigB-nail.log`）

| 发 | 台件（放／改什么） | `rc=` | 名册／分母读数 | 判 |
|---|---|---|---|---|
| B0 | pristine | **0** | 构造点＝1（`internal/agent/loop.go:675(same-package-bare-literal)`）、分母＝**293**（`cmd=40 internal=229 scripts=11 tools=13`） | 该绿 ⇒ 绿；分母与第三方 `find` 尺逐根对得上 |
| B1 | 新增 `internal/tools/rig309v2second.go`（限定名字面量） | **1** | 2 枚：`… \| internal/tools/rig309v2second.go:5(qualified-literal)`，分母 294 | 进攻①(a) 加一枚⇒**响** |
| B2 | 新增 `internal/agent/rig309v2bare.go`（同包裸名） | **1** | 2 枚：`… \| internal/agent/rig309v2bare.go:3(same-package-bare-literal)` | `309-v1` 条件"整包 24 枚源⛔ 解析"＝**已闭** |
| B3 | 新增 `tools/rig309v2root/rig.go` | **1** | 2 枚，`per root` 里 `tools=14` | 根 live（`tools`） |
| B4 | 新增 `scripts/rig309v2root/rig.go` | **1** | 2 枚，`per root` 里 `scripts=12` | 根 live（`scripts`）＝`Q-309-r2-scripts…` 那问的答案 |
| B5 | 把 `loop.go:675` 那一枚改写成 `req := *new(ToolRequest)` ＋ 逐字段赋值（编译得过、仍是真构造点） | **1** | `数到 0 枚 … 名册=(空名册)`，分母⛔ 漂（293） | 进攻①(b) 移除被点名那枚⇒**响**（判据是 `==1`⛔ `<=1`） |
| B6 | 钉里 `wantTotal 1→2`，树⛔ 变 | **1** | `数到 1 枚 … 钉的是 2 枚` | 进攻①(c) 往宽里挪一档⇒**当场红** |
| B7 | 钉里 `examined < 200` → `< 1` | **0** | 全绿 | **floor 可无声放宽**＝`Q-309-v2-2`（但 `==1` 与 per-root 空名册 `t.Fatalf` 两支⛔ 动） |
| B8 | roots 收成 `cmd internal` 两根，同时 `tools/` 里真有一枚第二构造点 | **0** | `roots=["..\\..\\cmd" "..\\..\\internal"] … 构造点=1 … 分母=269` | **roots 列表无人钉**＝`Q-309-v2-2`；净少扫 24 枚源而⛔ 红 |
| B9 | roots＝`[]` | **1** | `只解析了 0 枚非测试生产源（roots=[]）` | 空循环形⇒**响**（分母那一支） |
| B10 | `scripts` 那根改名成不存在的路径 | **1** | `walk` 直接 `t.Fatalf` | root 悄悄没了⇒**响** |
| B11 | matcher 打瞎（限定名 typo `ToolReqest`） | **1** | `control "shape 1 qualified literal": matcher counted 0, want 1` | 永红形／瞎尺⇒**响**（正控⛔ 恒真） |
| B12 | 我这发第一次的 catch-all 写法⛔ 合法 Go（`e declared and not used`） | — | `[build failed]` | ⛔ 作凭据；换 C1 重做 |
| B13 | 新增 `internal/perm/rig309v2commentonly.go`（`ToolRequest{` 只在**注释**里） | **0** | AST 钉仍＝1；同一棵树词面尺＝**4** | 词面尺数注释＝**假枚**、解析器尺⛔ 数（进攻③后半） |
| B14 | 新增 `internal/tools/rig309v2blind.go`：`return new(agent.ToolRequest)`（真·第二构造点，钉头具名漏形） | **0** | 构造点仍＝1、分母漂到 294 | ★**判死** `bridge.go:962` 那句 "a second one rings **wherever** it is written"＝`AC#3` 唯一条件 |

## C 组＝控制那一支的两向（树 B；件 `logs/rigC-controls.log`）

| 发 | 台件 | `rc=` | 读数 |
|---|---|---|---|
| C0 | pristine 正控 | **0** | `构造点=1 … 分母=293` |
| C1 | `case *ast.Ident: return names[e.Name]` → `return true`（catch-all） | **1** | `control "guard: bare literal in a NON-agent package": matcher counted 1, want 0` |
| C2 | 控制名册枚数（读码尺，⛔ 突变） | — | 行数＝**6**（`awk '/controls := \[\]struct/,/^\t\}$/' … \| grep -c '^[[:space:]]*{"'，读数记在 `logs/s9-behaviour-invariance.log` 首两行）；消费它的 `for` 在钉的 `:359` |

## D 组＝生产**调用点**（树 A，整包 `-v`；件 `logs/rigD-callsite.log`）

| 发 | 台件 | 顶层 FAIL／PASS | 净新增红 | 309 仪器 |
|---|---|---|---|---|
| D0 | pristine | 1／210 | 0 | `--- PASS` |
| D1 | `b.assessorFor(req.TaskID)` → `b.assessorFor("")`（Execute 不再传任务 id） | **8**／203 | **+7**＝`TestFSReadTaintFeedsR4`、`TestForeignMentionOfHostPathStillHitsR4OnRealBridge175r2`、`TestPointer183SiblingArtifactInTheSameDirectoryStillHits`、三枚 `TestPointer185*` | `--- PASS` |
| D2 | → `b.assessorFor(req.CorrelationID)`（票 242 的每调用一枚顶替任务 id） | **4**／207 | **+3**＝`TestToolCallRowsAreComplete`、`TestFSReadTaintFeedsR4`、`Test283IdentityChainThroughTheRealLoop` | `--- PASS` |

⇒ 调用点解耦有人接（别家测的牙），⛔ 是 309 自家仪器接的＝`Q-309-v2-5`。

## E 组＝`309-v1` 给 `AC#2` 那两条"⛔ 实测"的形状补盘上读数（树 B；件 `logs/rigE-shapes-on-disk.log`）

| 发 | 台件（放在钉真解析的文件里） | `rc=` | 名册里的 shape 标签 |
|---|---|---|---|
| E1 | `cmd/llmrecord/rig309v2e.go`：`var rig309v2Decl agent.ToolRequest` ＋ `rig309v2Decl.TaskID = ""`（**无复合字面量**） | **1** | `cmd/llmrecord/rig309v2e.go:5(var-decl-no-initializer)`，构造点 1→2 |
| E2 | `internal/agent/rig309v2e.go`：`type rig309v2Alias = ToolRequest` 再 `rig309v2Alias{TaskID: ""}` | **1** | `internal/agent/rig309v2e.go:5(alias-literal)` |
| E3 | `cmd/llmrecord/rig309v2e.go`：`type rig309v2Nested = []agent.ToolRequest` 再 `_ = rig309v2Nested{{TaskID: ""}}`（内层省类型） | **0** | ⛔ 计＝钉头具名漏形之一，我这一发补上读数 |
| E4 | 同 E3 那枚文件在场时并排两把尺 | — | 词面尺（按文件名字段滤）＝**2**／AST 钉＝**1** |

E 组结束后两枚树内目标 `cmp` 回原样（`cmp toolrequest_nail_309_test.go identical`／`cmp loop.go identical`，件末四行）。

## 只读尺（⛔ 突变，工作树面＝HEAD blob 面）

| 件 | 尺 | 读数 |
|---|---|---|
| `logs/r4-grep-rulers.txt` | 注释自己写的那把 `grep -rn 'ToolRequest{' internal cmd --include=*.go` ＋两支过滤 | 总 **71**／文件名以 `_test.go` 结尾 **69**／非测试（按**文件名字段**）**2**（`loop.go:675` ＋ `bridge.go:882` 注释本体）／再剥行首 `//`＝**1**；简写形 `\| grep -v _test.go`＝**2** |
| `logs/s5-and-ac3-readings.txt` | 起手框数尺＋点名段逐字 | 3 勾／3 未勾；段内⛔ 一枚绝对 `loop.go:<数字>`；`loop.go:[0-9]` 全档 **5** 枚（`:113/:114/:230/:252/:1119`，全在票 154/151 那族段落＝既有账 `Q-309-v1-anchor-debt-154`，编排者已裁⛔ 本票账） |
| `logs/s9-behaviour-invariance.log` | "行为逐字节不变"（区间＝`bd124b2a`⇄`f6067e5d`，⛔ 工作树⇄HEAD） | 改动行 **70 枚全是注释**、非注释改动行 **0**；两版 blob 剥整行注释 `cmp` **IDENTICAL**（854⇄854，sha `496ed6d683e8e070c92410dde1c6a6d39e14a069d1773a9b35980a2d040e32c2`）；守卫那行逐字命中 **1** |
| `logs/s8-rules-and-base.txt` | `grep -n 'func rule\|ctx\.taint' internal/risk/rules_gateway.go` | 规则函数 **4** 枚，读 `ctx.taint` 的只有 `ruleTaint`（`:101/:102/:105`）⇒ 注释"other rules are STRUCTURALLY still reached"＋"only R4 is dropped here"两句**如实** |
| `logs/s10-legclaims-vs-mine.txt` | 腿自报那九枚件 vs 我这一发 | `r2/logs/final-{build,vet,d22scan,pathlen}.log` 末行 `rc=`＝0/0/0/0，我各跑一发＝**0/0/0/0** 同色；309 两腿（`16e86bae`／`d75db9ee`）只碰 **3** 枚 `.go`（`bridge.go` ＋两枚 `*_309_test.go`）；禁区面（`thresholds.go`／`allowlist.txt`／三枚冻结件）在**这两笔里**零命中 |
| `logs/s11-universal-claims.txt` | 全称词清点（穷举同族写法） | `bridge.go` 命中 `:951`（自限句，真）／`:957`（带尺＋自限，真）／`:959`＋`:962`（**every/wherever，⛔ 真**＝B14 反证）／`:982-985`（this limb＋自陈⛔ guarantee，真）；`grep -n 'Every other gate\|still runs'`＝**0** 命中＝那句已被射程句替掉 |
| `logs/s12-limb-worlds.txt`＋`s17-bare-name-ruler.txt` | 两肢"世界"普查 | `Assessor:` 全仓 **2** 行＝注释 `bridge.go:132` ＋测试 `ticket90_test.go:366`（生产写点 **0**）；`Provenance:` 在 tools 测试里 **24** 处；`^package agent$`＝**34** 枚文件、目录集合只有 `internal/agent` ⇒ 裸名那一支的门今天按目录闭合 |
| `logs/s13/s14/s15` | roots 覆盖与根外 | 四根之外非测试 `.go`＝**615** 枚，全部在点目录（`.scratch/**`，Go 工具链自己⛔ 看见；tracked 314）；⛔ 点目录之外只有 **1** 枚＝`frontend/embed.go`，且 `go list ./...`（**35** 枚包）里 `github.com/CarlosShao/wisp/frontend` 确实成包 ⇒ `Q-309-v2-4`；nested module 三枚前缀逐字含 `github.com/CarlosShao/wisp/`＋活例 `tools/signmodels/main.go:34` ⇒ r2 的进根判据如实 |
| `logs/s16-anchors-and-contract.txt` | 内容锚在盘＋契约声称 | 五枚锚 `grep -c`＝**1/1/1/1/1**；`verdict := b.assessorFor(req.TaskID)`＝**1**；`git diff bd124b2a f6067e5d -- bridge.go` 里带 `C19/C25/C26` 的增删行＝**0**；票 154 `gate-clauses.sh` 在场且 `G1b\|G5` 命中 **34** 处 |
