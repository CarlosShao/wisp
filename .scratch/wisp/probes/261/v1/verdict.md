# 票 261 · 腿 261-v1 裁决表（非实现者验收：攻 261-r1 修法的门有没有牙）

代理：`261-v1`（**非实现者**验收子代理；实现者＝`261-r1`，其产码三枚 commit `2bbe7086`→`e2bfd7c8`→`9d0cd518` 已入库）。
票面：`.scratch/wisp/issues/261-the-model-enabled-key-promises-removal-from-discovery-and-selection-but-no-production-code-reads-it-while-the-default-true-tag-is-inert-for-map-entries.md`
实现件：`.scratch/wisp/probes/261/r1/impl.md`；前程件：`.scratch/wisp/probes/261/p1/verdict.md`。
本腿定位：**攻防**——不复述 r1/p1 已留档读数当自己的凭据；变异复认两形各至少 1 发自跑、恒真两问、三条入口现跑、生产调用者问句、旁格核对、AC 格判语。勾框一枚未碰（归编排者）。

---

## §0 起手锚与起手名册

- 起手时刻：`2026-10-03T16:25:51+08:00`（`date -Iseconds` 自取）
- 起手 HEAD：`89c863f168cf27c76bfae165666f34a5dcd3fea9`（`git log -1 --format=%H` 自取），分支 `dev`（r1 三枚 commit 均在历史内）
- 起手 `git status --porcelain internal/llm internal/config`：**空**（rc=0，零行）——两包起点干净，无他人脏面
- 起手 md5（跑 overlay 前 pre 记录）：
  - `internal/llm/resolver.go` = `e8a2cbc85e2393e6b9ca9c8db849dc1b`
  - `internal/llm/enabled_reach_261_test.go` = `dede3bb0863830bed0aa5f0e1d89e6de`
  - `internal/llm/enabled_gate_261_r1_test.go` = `3631fe757cba6ce7c00083df1683d580`
  - `internal/config/enabled_261_test.go` = `f2eef15374a39ecdd747e74298f72b2c`
- 起手绿名册：`go test -count=1 -v ./internal/llm/ ./internal/config/` → **rc=0，155 PASS / 0 FAIL / 0 SKIP**
  （log `.scratch/wisp/probes/261/v1/v1-base-test.log`；名册逐名 `v1-base-roster.txt` 155 行）
- 与 r1 终态名册（`.scratch/wisp/probes/261/r1-final-roster.txt`，155 行）`comm` 双向：`-23`＝**0 行**、`-13`＝**0 行**——零丢名零增名，完全一致。
- **收尾复归核对（终态 16:42:27+08:00）**：四枚 md5 逐枚与起手相等（前二见上，config 尺两枚同样未动）；
  终态 `git status --porcelain internal/llm internal/config`＝**空**（rc=0）——本腿全程零跟踪树残留。
- 门禁补跑（本腿责任面）：`sh scripts/d22scan.sh` → **rc=0 clean**（`v1-d22scan.log`，ban #8 internal/=498、cmd/=92 等全套无违例）。
- 骨架 commit＝`bdabee2c`（轮次帽第 8 轮前落盘）；本行所在＝终态 commit（`git log -1` 自取）。

## §1 变异复认（两形各至少 1 发自跑，红句逐字）

**撤门形**（`mut/`，overlay 把跟踪 resolver.go 换成从旧锚抽出的无门形；我先 diff 复核：
撤门形与跟踪树 diff 恰为 18 行＝两处门，`resolveEndpoint` 的 `if !spec.Enabled` 拒绝＋`DiscoveredModels` 的 skip，零其他差异）：

- 跑法：`go test -count=1 -v -overlay=overlay.json -run 'TestTicket261R1|TestTicket261P1' <repo>/internal/llm/` → **rc=1，6/6 全红**（log `v1-mut-run.log`，与 r1 留档的 `mut-run.log` 红句逐字一致）。
- 摘掉门后**红的具名用例**（恒真问①答案，6 枚全名单）：
  1. `TestTicket261R1EnumerationSkipsDisabled`：`DiscoveredModels = [m-off m-on], want the disabled entry removed from discovery`（:58）
  2. `TestTicket261R1SelectionRefusesDisabledByExactWording`：`ResolveChain(disabled) succeeded; the gate is missing`（:80）
  3. `TestTicket261R1RoleFallbackAlsoGated`：`ResolveRole fallback onto a disabled first element = <nil>, want the disabled refusal (not a silent skip)`（:118）
  4. `TestTicket261P1EnumerationSkipsAndSelectionRefusesDisabled`：`DiscoveredModels = [t261-nokey t261-off t261-on], want the disabled entry "t261-off" removed from discovery (ModelSpec.Enabled promise)`（:183，同句对 t261-nokey 再红一发）
  5. `TestTicket261P1GateBlocksProviderAndCost`：`ResolveChain resolved a disabled model with the gate in`（:230）
  6. `TestTicket261P1FlagInversionReversesEveryReading`：`cost-hop: the disabled entry produced 90 micros, want 0`（:338；另有枚举/选择/provider-hop 三向红句同 log）
- 与 r1 impl §3 留档红句**逐字对账：6/6 全部一致**（我抓到的行号比 r1 留档更全：p1 两改写枚各有双红句）。
- md5 前后复归＝零残留（见 §0）。

**反形判据**（`mut/inverted/`，overlay 只换两枚测试文件为断言反写形＋p1 原版三枚，**不换 resolver**，对带门跟踪树跑）：

- 我自跑 → **rc=1，6/6 全红**（log `v1-inv-run.log`；`--- FAIL` 顶格 6 枚）。抽复核 2 枚红句逐字：
  - `TestTicket261R1EnumerationSkipsDisabled`（反形版）：`INVERTED: enabled entry "m-on" listed, want NOT listed`
  - `TestTicket261R1SelectionRefusesDisabledByExactWording`（反形版）：`INVERTED: ResolveChain(disabled) = config: llm: model "m-off" of provider "p261" is disabled (enabled=false); re-enable it or remove the entry, want nil (admitted)`
  - 另 p1 原版三枚对带门树全红（`TestTicket261P1EnumerationAndSelectionAdmitDisabled`/`TestTicket261P1DisabledModelReachesProviderAndCost`/`TestTicket261P1FlagInversionChangesNoReading` 顶格 FAIL），证明改写后的用例测的是新语义、不是换名不改心。
- 与 r1 留档 `inv-run.log` 对账：逐字一致。r1 impl 自承"反形台构建错误一次已内联修复后重跑"——我核其留档 log 为真测试运行（有 `--- FAIL` 顶格与用例体红句，非 setup-failed 假红），复核通过。
- ⚠ 过程注记（自翻，非他人缺陷）：我首跑反形时包路径少算一层 `..`，拿到一次 `stat ...\scratch\internal\llm: directory not found` 的 **setup-failed rc=1**——这正说明 rc=1 不必等于"测试红"，必须数 `--- FAIL`。修正路径后重跑才是上面的 6/6。此坑已避开，读数以 `--- FAIL` 计。

## §2 恒真两问

- **问①（摘掉门后哪几枚具名用例红）**：§1 撤门形 6 枚全名单（三枚 R1＋三枚 P1 改写名，逐枚红句逐字在档）。只报 rc 不算——每枚都有具体红句，尺有牙。
- **问②（判据换成反形会不会也全绿）**：不会。反形台对**带门跟踪树** 6/6 红（§1）；对**撤门形**，反形用例按其断言方向必然转绿（此向我没跑——它等价于"撤门形正向红"的镜像，而撤门形 6/6 红已自跑在档）。两向合起来：正向尺对"门在"敏感（撤门⇒红），反形尺对"门在"敏感（门在⇒红）——**不存在恒真向**。
- 附带标定：`DiscoveredModels` 的 ghost 对照（未知 provider 返回 nil）与选择线的 ghost 对照（`unknown model` 文案）在正反两形里都活着——尺对"在不在目录"同样敏感，不是只对 enabled 钝感的单点尺。

## §3 三条入口现跑读数（临时文件＋`config.LoadFile`，⛔ 未碰跟踪码）

即弃探针 `.scratch/wisp/probes/261/v1/probe3/main.go`（手工 TOML＝票 257 首启形状：`t-nokey` 缺键／`t-off` 显式 flag／`t-on` 显式 true；base_url 指向 `http://127.0.0.1:9/v1` 不可达——门必须在拨号前拒）。读数 log＝`v1-three-entries.log`（2026-10-03 16:31:13+08:00 现跑，逐字）：

- **入口 1 `text_chain` 点名 disabled**：`ResolveChain()` →
  `err="config: llm: model \"t-off\" of provider \"mock261\" is disabled (enabled=false); re-enable it or remove the entry"`，`eps=[]`。**＝拒。**
- **入口 2 `roles.chat` 回退**（text_chain 打头是 disabled 条目，run.go:436 未设角色回退形状）：`ResolveRole(RoleChat)` → 同一句 disabled 拒绝。**＝拒。**
- **入口 2b `roles.chat` 直配**（角色直接点名 provider/model——派单三条之外的补测，同一 `resolveEndpoint` 线）：同一句 disabled 拒绝。**＝拒。**
- **入口 3 `DiscoveredModels` 枚举**：返回 `[t-on]`——disabled 的 `t-off` **与缺键的 `t-nokey` 都不在列**。**＝过滤生效，且缺键同被门住（缺键落 false，没有后门）。**
- **反色对照（同一文件只翻 `enabled = false`→`true`，手编重开）**：`ResolveChain` err=`<nil>` 且 endpoint 照回、`ResolveRole` err=`<nil>`、枚举 `[t-off t-on]`。**两向实测反转成立。**
- 与 r1 impl §2 表格逐格对账：一致。

**本腿新钉的一格（r1 没写透的）**：**缺键条目＝firstrun 官方引导教用户手写的形状**——
`cmd/wisp/firstrun.go:116-118` 的引导原话只提 `api_key_ref`、`base_url`、`models.<id>`、`text_chain`/`roles.chat`，
**通篇不提 `enabled`**（`grep -c enabled`＝0），而内置预设只供 protocol/base_url（`internal/config/defaults.go:14-17` 注释逐字
"presets supply protocol/base_url only"、`applyPresets` 在 `internal/config/loader.go:188-204` 只填这两枚，**不带模型条目也不带 Enabled**）。
⇒ 按 firstrun 文案逐字照做的首启用户，配完 key + base_url + models.<id> + 点名，今天会被门以
`is disabled (enabled=false)` 拒掉——**承诺是兑现了，但兑现的语义是"没写 enabled＝关"**，与注释承诺句
"Enabled removes the model from discovery/selection"（`internal/config/schema.go:369-370`，字段 :371）的默认直觉（不写＝开）相反。
这是 AC#1ⓐ 的**带条件判语**素材（§6），也是给编排者的新上报（§7 推翻清单第 2 条）。

## §4 生产调用者名册＋「行为变没变」判语

**`resolveEndpoint` 的非测试调用者**（`grep -rn` 排除 `_test.go`，含经 ResolveChain/ResolveRole 的间接调用；BuildChain 生产调用者＝零枚，`DiscoveredModels` 生产调用者＝零枚——两问各自全仓 grep 复认，命中全落在测试与 `.scratch`）：

| # | 调用点 | 形状 | 门后行为（现跑读数） |
|---|---|---|---|
| 1 | `cmd/wisp/run.go:436` `res.ResolveRole(llm.RoleChat)`（`wisp run` 文本主路） | 未设角色回退→链首；或 roles.chat 直配 | disabled ⇒ 拒（§3 入口 2/2b）；拒后 run.go:437-439 走既有 `Unconfigured` 退码 2——**不新增退码，复用既有失败通道** |
| 2 | `cmd/wisp/providers.go:170` `res.ResolveChain()`（`wisp providers probe`，TextChain=被点名对） | 点名驱动 | disabled ⇒ 拒（providers.go:171-174 报"端点未就绪"退 2）；其 `:193` 的 `cfg.LLM.Providers[provider].Models[model].Capabilities` 直读发生在 ResolveChain **成功之后**——对 disabled 不可达，**非旁路** |

**「行为变没变」判语**（引用 §3 现跑＋§1 自跑的两色，⛔ 不把"测试绿了"当"产品行为变了"）：

- **变了——但只对今天真实存在的"关"路径**。先纠正派单问句的前提（§7 第 1 条）：**设置里今天没有"关一枚模型"这个动作**——
  `internal/config/settings.go` 全部 setter 零枚写 `enabled`（grep 亲验），面板写侧（`cmd/wisp/panel_config_store.go`）同样零枚。
  今天产生 `enabled=false` 的唯一真实路径＝**手编 config.toml**（显式写 false，或缺键落 false——后者还会被任何一次设置页保存经
  `mergeWrite→SaveFile→MarshalCanonical` 无 omitempty 物化成显式 false，p1 已钉，本腿不重跑）。
- 对这条真实路径：**旧形**（撤门 overlay 自跑复现）＝点名解析成功、provider 收请求、计价 90 micros（§1 撤门形第 5/6 枚红句即旧形行为的直接读数）；
  **新形**（跟踪树现跑）＝三条入口全部具名拒绝、枚举不列、手编翻 true 即恢复（§3）。**产品行为确实变了。**
- **没变的面**：面板快照 `ModelCount` 仍数目录全量（`panel_config_store.go:112` `for range p.Models`，旁格核对 §5-3——"关掉"在面板计数上不可感知，票 145 AC#2b 格）；`DiscoveredModels` 生产调用者零枚，无任何用户可见清单因此变化；验证层（`internal/config/catalog.go:66`）仍只问"在不在"，加载期放行、解析期才拒（r1 impl §5.2 自承，我复认）。

## §5 旁格核对（三处，均一字未动）

1. **`internal/config/defaults.go`**：`case reflect.Map:` 分支逐字仍只有 `// leave nil (see NewDefaults)`（:77-78 区）；该文件最后被碰＝`fe126e27`（2026-09-19，261 窗口外）。✓
2. **设置写侧 `internal/config/settings.go`**：`SetModelContextWindow`/`setModelPriceLeaf` 两处 `spec, ok := p.Models[model]` 原样（:115/:153），无任何 enabled setter；最后被碰＝`0d87a681`（2026-10-02，票 248，261 窗口外）。✓
3. **`cmd/wisp/panel_config_store.go:112`**：逐字 `for range p.Models {`／`row.ModelCount++` 原样（不过滤）；最后被碰＝`0d87a681`（同上）。✓

另核 r1 主 commit `e2bfd7c8` 的 `--name-status` 全名单：产码只 `internal/llm/resolver.go`（M）＋两枚测试文件（A/M），`docs/**`、`thresholds.go`、golden、`allowlist.txt`、票面勾框零触碰。✓

## §6 AC 格判语（逐格；勾框归编排者，本腿只给判语）

- **AC#0ⓐ（真路径存在）＝成立**。p1 三向现跑读数在档且完整（点名/派发/计费三跳＋翻 flag 不反转＋突变台标定）；本腿以其撤门 overlay 的自跑复现（§1 第 5/6 枚红句：点名成功、cost 90 micros）作为对"旧形行为"的**独立再演示**，其生产入口名册的 grep 口径我复认（生产读者零枚——对旧锚点成立）。判语：**成立**。
- **AC#0ⓑ（discover 显式写 true）＝成立但带条件**。写侧事实成立：`internal/llm/discover.go:117-121` 构造体逐字含 `Enabled: true,`（:119，本腿亲读）。条件：`ImportDiscovered` 生产调用者＝**零枚**（全仓 grep，仅定义＋两枚测试文件——本腿复认 p1 判语），写侧执法是纸面的；p1 判"成立但不够"我同意，且"不够"的读侧半边已由 AC#1 的门补上。判语：**成立但带条件（写侧成立、射程为零；读侧由门兜住）**。
- **AC#1ⓐ（落地格）＝成立但带条件**。判据四件全在我手上过了攻防：①两向实测（§3 两色，非测试壳——临时文件＋真 `config.LoadFile`＋真 resolver）；②变异两形 6/6 红（§1 自跑复认，非只信留档）；③恒真两问（§2，问①给具名用例与逐字红句，问②反形 6/6 红自跑）；④三条入口现跑全拒（§3）。**条件四条**：(a) 拒绝在解析期不在加载期（验证层 `catalog.go:66` 放行）；(b) 面板 ModelCount 与可选清单的差从此存在（票 145 格）；(c) settings 写侧仍无"重新打开"路径，关后只能手编文件翻 true；(d) **firstrun 教的最短手写形状（缺键）会被门拒**——"不写 enabled＝关"的兑现语义与引导文案冲突（§3 新钉，产品判断归编排者）。判语：**成立但带条件（四条具名，无一枚是牙口问题）**。
- **AC#2（普查件 `1ba22d25`）＝成立**。本腿独立抽验（非复读普查件）：`internal/config/schema.go` 五张 map（:317 值 string、:400 `Models map[string]ModelSpec`、:419 `Providers map[string]Provider`、:575 `Entries map[string]PluginEntry`、:590 值 string），值结构带 default tag 的字段恰 **3 枚**：`Provider.Billing`（:391）、`ModelSpec.Billing`（:361）、`ModelSpec.Enabled`（:371）——与普查件"3 枚全中"一致；每枚"缺键落零值会不会被读出意义"的判语在普查件内有 file:line 依据（validate.go:197/:160 空串放行、Enabled＝票 261 本体），且 `PluginEntry.Enabled`（:556）**无** default tag＝"没写 tag 所以没这个坑"的诚实形（普查件 §5 行自己点名了这格的同构陷阱，我复认其射程判断）。判语：**成立**。

## §7 推翻清单（票面、r1 impl、p1 verdict、本派单每一句都是待验断言）

1. **本派单问句 4 的前提"用户在设置里关一枚模型"——不成立**（最狠一条）：设置写侧与面板写侧**零枚 enabled 写者**（§4），今天"关"的唯一路径＝手编 config.toml。问句据此改写为"手编关掉之后行为变没变"——§4 已答（变了，两色现跑）。派单把它当现成用户动作，与票 257/firstrun 的现实形状不符。
2. **r1 impl §1.2/§2 对"缺键形状"的处置只写了一半**：门把缺键当关（拒绝），这在测试里是对的（"缺键落 false 不能当后门"），但 r1 没有指出**firstrun 官方引导教的就是缺键形状**——引导文案不提 enabled、预设不填 Enabled（§3 新钉）⇒ 首启照文案做的用户会被门拒。这不是门没牙，是**兑现语义与产品引导的冲突**，r1 impl §5 未列此条，本腿补名。
3. 票面 §现量-2"枚举那一处不读它（:287-288）"——对旧锚点成立，终态不成立（r1 自认，我复认终态 `:297-300` 有判据）。非新推翻，登记在案。
4. r1 impl §3 红句留档——**全部逐字对上**（撤门 6 枚＋反形抽 2 枚自跑对账），零推翻；我补全了 p1 两改写枚的**双红句**（r1 只引了单句）。
5. p1 verdict §2/§3 的行号引文（loader.go:117-118、firstrun.go:116-118、panel_config_store.go:112、schema.go:371 等）——本腿逐枚亲读，**全部成立**，零推翻。
6. 自翻一处：我首跑反形台的包路径错误产生了 rc=1 的 setup-failed（§1 注记）——方法论登记：**rc=1 必须 `--- FAIL` 计数背书**，对 r1 留档 log 我据此核过（真测试运行，非假红）。

**推翻枚数＝2 条实质（第 1、2 条）＋4 条精化/复认登记。最狠＝第 1 条（派单前提不成立）。**

## §8 判不动／量不到

1. **真·外部计费**：与 p1/r1 同界——判据止于"请求是否到达 mockllm（路由自计数）＋生产计价函数读数"；真实 provider 账户扣款不在任何测试射程。复量法＝p1 verdict §5.1 同款。
2. **面板 UI 是否有模型开关入口**：`frontend/**`／`design/**` 两层禁区（不读不引）——本腿只判到 Go 后端写侧零枚 enabled setter，面板前端有无此控件判不动。复量法＝归 owner 侧会话。
3. **拒绝时机（加载期 vs 解析期）对运维体验够不够**：产品判断，归编排者；复量法＝若要前移到 `checkChainElement`，加 enabled 读后重跑本腿 §3 全套尺（尺已备好）。
4. **firstrun 最短形状被门拒要不要改引导文案／给缺键填 true**：改引导文案＝碰 `cmd/wisp/firstrun.go`（本腿禁区外但非本票写面）；给缺键填 true＝碰 `defaults.go` 行为（票面 AC#2 明令另格）或改 tag 语义（可能碰 D36 段树＝契约变更须人工批准）。本腿只具名（§7-2），不裁。
5. **`go build ./...` 未跑**：派单禁区（本机有别的测量在排队）；本腿门禁面＝两包包级测试＋d22scan，vet 未跑（非派单要求项）。若需补跑由编排者排。
6. **`ImportDiscovered` 写侧真值在生产中的意义**：生产调用者零枚，"发现那一路写 true"在接线之前只是纸面保证——此格与 AC#0ⓑ 的条件同源，接线后的复量法＝以 fixture 调 `ImportDiscovered` 断言 `Enabled` 位。
