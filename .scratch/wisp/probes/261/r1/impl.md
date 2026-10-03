# 票 261 · 腿 261-r1 证据件（AC#1 走 ⓐ：让承诺兑现）

代理：`261-r1`（写码子代理）。派单裁定：两处落门，都在 `internal/llm/resolver.go`——
枚举点 `DiscoveredModels` 跳过 `Enabled==false`；选择点 `resolveEndpoint` 对目录里有但被关掉的条目具名拒绝。
**选择点是编排者据 AC#0 读数扩的射程**（AC#0 ⓐ 判死的结局是"被选中并产生花费"，只过滤枚举治不到
`text_chain`/`roles.chat` 点名那条路）——本腿照做，此句即具名登记。

---

## §0 起手锚与起手名册（现跑自取）

- 起手时刻：`2026-10-03T15:55:46+08:00`（`date -Iseconds` 自取）
- 起手 HEAD：`6c7a12b656188405e7ab40d6bb7c47066ff55dd6`（`git log -1 --format=%H` 自取），分支 `dev`
- 起手 `git status --porcelain internal/llm internal/config`：**空**（rc=0，零行输出）——
  与编排者宣称"我已现跑＝空"一致，**无他人脏面**。
- 起手绿名册：`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 -v ./internal/llm/ ./internal/config/`
  → **152 PASS / 0 FAIL / 0 SKIP**；两包 ok（llm 36.408s／config 0.943s），rc=0。
  名册逐名 `.scratch/wisp/probes/261/r1-base-roster-clean.txt`（152 行，剥时延后缀）。
- 与 `261-p1` 终态名册（`.scratch/wisp/probes/261/p1/final-roster.txt`，152 行）逐名 `comm` 比对：
  - `comm -23`（p1 有、我无）＝**0 行**；`comm -13`（我有、p1 无）＝**0 行**。**完全一致，零丢名零增名。**
  - ⚠ 注记：p1 名册是裸测试名（无时延），我最初从 `-v` log 抽出的名册带 ` (0.02s)` 时延后缀，
    首次 comm 被后缀污染（152 行全"差异"）；剥后缀（`sed 's/ (.*$//'`）后重比才干净。
    中间产物 `.scratch/wisp/probes/261/r1-base-roster.txt`（带后缀版）与 `r1-p1final-roster.txt`
    （误从带后缀 log 抽 p1 名册得 0 行的废件）**只建不删**，留作过程痕迹。
- 基线 log：`.scratch/wisp/probes/261/r1-baseline-test.log`

## §0b 起手预检读数（撞钉预检，逐枚）

1. **行号复认**（推翻清单的靶子，逐枚现读）：
   - `internal/llm/resolver.go:287` 逐字 `out := make([]string, 0, len(p.Models))`、`:288` 逐字
     `for id := range p.Models {`——**编排者写"循环体实测在 :287-288"成立**（我注记：:287 是
     `out :=` 行不是循环行，循环体是 :288-290；语义无差）。
   - `internal/llm/resolver.go:113` `func (r *Resolver) resolveEndpoint(...)`、`:119` 逐字
     `spec, ok := p.Models[model]`——**成立**。
   - `internal/config/schema.go:369-370` 承诺句两行、`:371` 字段——**成立**（票面引文折成两行是拼接转写）。
   - `internal/config/defaults.go` 无任何 `Models` 预置（grep 零命中）⇒ `NewDefaults()` 不产模型条目，
     生产 `wisp run` 的正常路径（零手写条目）不受本门影响。
2. **撞钉预检（fixture 缺 `Enabled` 又期待可达的形状）**：
   全仓 grep `ResolveChain|ResolveRole|DiscoveredModels --include=*_test.go` 命中仅 2 枚文件：
   `internal/llm/enabled_reach_261_test.go`（本票前腿仪器，按派单改写）与 `internal/llm/catalog_test.go`。
   另枚全仓 `ModelSpec{` 测试构造点（13 处）逐枚核：
   - `internal/llm/catalog_test.go`：`:29,36,43,133,209` 全部显式 `Enabled: true`（编排者点名 5 处复认成立）；
     `:217` 逐字 `if !m2.Enabled || ...` 是 `ImportDiscovered` 的判据（不走 resolver，本门碰不到它）。
   - `internal/llm/fallback_test.go`：`:44,51,265` 全部显式 `Enabled: true`（编排者点名 3 处复认成立）。
   - `internal/llm/openaichat/mockllm_integ_test.go:418`：`if !spec.Enabled || ...` 是 discover/import
     路径的读（不走 resolver），不受影响。
   - `internal/config/{validate,loader,catalog,boundary}_test.go` 的 `ModelSpec{}`/缺 Enabled 构造：
     全部只过 `config.LoadFile`/validate（**config 包不 import llm**），resolver 门不在其路径上。
   - `internal/agent`、`cmd/wisp` 的测试：零处构造 `ModelSpec` 后过 resolver。
   - **结论：无一处"建 map 时省掉 Enabled 又期待那枚模型可达"的形状；零枚 fixture 需要补写。**
3. **`LoadFile` 第二参**：grep 全仓 `LoadFile\(.+, [^n]` → 生产调用者全部传 `nil` 或 SecretResolver
   变量（`st`/`res`），**没有任何调用者把 `*llm.Resolver` 塞进 config.LoadFile**。
   注记：编排者预检清单里"resolver.go:157 `func LoadFile`"一句是**笔误**（resolver.go 里没有 LoadFile）；
   真身＝`internal/config/loader.go:41 func LoadFile(path string, res SecretResolver)`。
4. **`wisp run` 正常路径不受门影响**：`cmd/wisp/run.go:435` 的 cfg 出自 `config.LoadFile`；
   零手写模型条目时 `p.Models` 为空 map ⇒ 无条目可判 false，门为透明。

## §1 改了什么（逐枚 file:line，终态行号以终态 HEAD 为准）

### §1.1 产码：`internal/llm/resolver.go`（唯一产码文件，两处）

- **选择点**（终态 `:124-131`，插在 `:119` `spec, ok := p.Models[model]` 的 ok-分支之后）：
  ```go
  if !spec.Enabled {
      return Endpoint{}, observe.New(observe.ClassConfig,
          fmt.Sprintf("llm: model %q of provider %q is disabled (enabled=false); re-enable it or remove the entry", model, provider))
  }
  ```
  文案走同文件既有的 `observe.New(observe.ClassConfig, ...)` 定式；与既有
  `llm: unknown model %q of provider %q`（`:122`）**无共享词串**（unknown model ↔ is disabled）；
  含具名词 `disabled`＋provider 名＋模型名。**作用于所有 resolveEndpoint 调用者**：
  `ResolveChain`（text_chain）、`ResolveRole` 直配（roles.*）、`ResolveRole` 未设角色回退
  （run.go:436 形状）——一处门，三条入口全关。
- **枚举点**（终态 `:291-306`，`DiscoveredModels`）：`for id := range` → `for id, spec := range`
  ＋ `if !spec.Enabled { continue }`（终态 `:297-300`）。注释同步改写为 enabled-catalog 语义。

### §1.2 被改写的 `261-p1` 用例（具名登记，文件 `internal/llm/enabled_reach_261_test.go`）

**派单授权**：编排者裁定"把它们就地改成修后语义"。改写动作＝测试函数**改名＋断言反转**，
fixture（`enabled261TOML`/`loadEnabled261`/`recordingRT`/`reach`/`reachReading`）原样保留
（`reachReading` 内 `reachedProvider` 判据从 `==1` 放宽为 `>=1`，因门落后 disabled 侧不再发起请求，
两侧路由计数语义见 §1.3）。

| 旧名（p1 落的） | 新名（r1 改写） | 原断言（旧语义，逐点） | 为什么它不再是本票的凭据 |
|---|---|---|---|
| `TestTicket261P1EnumerationAndSelectionAdmitDisabled` | `TestTicket261P1EnumerationSkipsAndSelectionRefusesDisabled`（终态 :159） | 断言 disabled/nokey 三枚**全被枚举**、`ResolveChain`/`ResolveRole` 对 disabled **返回 nil 错**、nokey 可点名可达 | 它是 AC#0 对**缺陷现状**的描述；门落后这些断言必然全红。改写后钉"承诺在兑现"：枚举只回 enabled、disabled 被 `disabled` 文案拒、nokey 同拒（缺键落 false 不能当后门） |
| `TestTicket261P1FlagInversionChangesNoReading` | `TestTicket261P1FlagInversionReversesEveryReading`（终态 :292） | 断言翻 true **四项读数逐项不变**（off.micros>0 且 on.micros>0 等）、ghost 被拒 | "flag 不动任何读数"＝死键的证明；门落后读数必须反转。改写后逐项断言反转（off 不列/不解析/不出网/0 micros；on 全反）＋ghost 保持 `unknown model` 文案（与 disabled 拒绝可分辨） |
| `TestTicket261P1DisabledModelReachesProviderAndCost` | `TestTicket261P1GateBlocksProviderAndCost`（终态 :222）——**未删，反转为反向尺** | 断言 disabled 条目真出网（mockllm chat 计数=1）、wire body 带 model id、`Cost.AddUsage` 用该条目价卡算出 >0 micros | 派单明令"不许删，改造成门在 ⇒ 到不了 provider、cost 恒 0 的反向尺"。改写后：disabled 条目 ResolveChain 即拒（provider 计数恒 0）、enabled 对照组证明仪器活着（真出网真计费）、disabled 侧 cost 读数恒 0 |

### §1.3 新增测试：`internal/llm/enabled_gate_261_r1_test.go`（本腿新文件，3 枚）

- `TestTicket261R1EnumerationSkipsDisabled`（:50）：最小 resolver fixture（`m-on`/`m-off`），
  断言枚举恰好只回 enabled 那枚＋unknown provider 仍回 nil。
- `TestTicket261R1SelectionRefusesDisabledByExactWording`（:75）：断言 ClassConfig、文案含
  `disabled`＋provider 名＋模型名、**不含** `unknown model`；ghost 对照保持 `unknown model` 文案
  （两种拒绝在日志里可分辨）。
- `TestTicket261R1RoleFallbackAlsoGated`（:114）：未设角色回退走同一 `resolveEndpoint` 行；
  disabled 打头 ⇒ 拒（**不是静默跳过**——跳过会背着用户重排链）；enabled 打头 ⇒ 正常。

### §1.4 未碰清单（禁区复认）

`internal/config/defaults.go`、`internal/config/settings.go`、`internal/config/schema.go`、
三张 provider 适配器（`internal/llm/openaichat/**` 等）、`cmd/wisp/panel_config_store.go`、
`docs/PLAN.md`、`docs/specs/**`、`internal/observe/thresholds.go`、任何 golden、任何 allowlist.txt、
票面勾选框、`frontend/**`/`design/**`（未打开、未引用）——**全部零字节改动**
（终态 `git status --porcelain` 三行，见 §4 末）。

## §2 两向读数（同一把尺，两色，全现跑）

fixture＝p1 的 `enabled261TOML`（`t261-nokey` 缺键／`t261-off` 显式 false／`t261-on` 显式 true，
价卡 in=2500000 out=10000000）＋真 `config.LoadFile`＋真 resolver＋真 BuildEndpointProvider＋
live mockllm HTTP＋生产计费函数 `internal/agent/cost.go`。

| 尺 | enabled=false（门关） | 补回 true（门开） | 反转? |
|---|---|---|---|
| 枚举 `DiscoveredModels("mock261")` | 只含 `t261-on`；`t261-off`/`t261-nokey` **均不在**（缺键落 false 同样被门住） | 三枚全在 | ✅ |
| 点名 `ResolveChain(["mock261/t261-off"])` | err≠nil，ClassConfig，逐字（见下） | err=nil，endpoint 照回 | ✅ |
| role 回退 `ResolveRole(RoleChat)`（text_chain 打头=disabled） | 同一 disabled 拒绝 | 正常解析 | ✅ |
| provider hop（mockllm chat 路由自计数） | disabled 条目 **0 请求**（ResolveChain 在 BuildEndpointProvider 上游就拒） | 真 hit | ✅ |
| cost hop（生产 `Cost.AddUsage`） | disabled 条目 **0**（无 usage 可计） | 12 in／2 out → **50 micros**（即弃探针逐字 `ON-COST: usage=12/2 delta=50 micros=50`；探针已删，p1 时点同形状读数是 90 micros＝12 in/6 out） | ✅ |

**点名拒绝的逐字错误句**（即弃探针现抓，探针已删，跟踪树零残留）：
```
config: llm: model "m-off" of provider "p261" is disabled (enabled=false); re-enable it or remove the entry
```
（用真实 fixture 时即 `config: llm: model "t261-off" of provider "mock261" is disabled (enabled=false); re-enable it or remove the entry`。）

**与既有 unknown-model 文案的分辨对照**（同一条选择线上两形并存，测试钉死）：
- 目录里有但被关：`... is disabled (enabled=false); ...`
- 目录里根本没有：`config: llm: unknown model "t261-ghost" of provider "mock261"`

两色对照 `TestTicket261P1FlagInversionReversesEveryReading` 一枚测试内完成（同文件只翻一行 flag，
四项读数逐项断言反转）。

## §3 变异自证（撤门必红；跟踪树零改动）

**突变台**：`.scratch/wisp/probes/261/r1/mut/`
- `resolver.ungated.go`＝**从起手 HEAD `6c7a12b6` 原样抽出的 resolver.go**（即无门形；md5
  `dbef11a638db2e99fbbaec5e2a97180c`）。与跟踪工作树的 diff 恰为两处门（§1.1），零其他差异。
- `overlay.json`（p1 同款格式）：Replace 跟踪 resolver.go → ungated 形。
- 跑法：`go test -count=1 -overlay=... -run 'TestTicket261R1|TestTicket261P1' ./internal/llm/`
  → `mut_rc=1`（预期红）；**6/6 全红**：
  - `TestTicket261R1EnumerationSkipsDisabled`：红句逐字
    `DiscoveredModels = [m-off m-on], want the disabled entry removed from discovery`
  - `TestTicket261R1SelectionRefusesDisabledByExactWording`：
    `ResolveChain(disabled) succeeded; the gate is missing`
  - `TestTicket261R1RoleFallbackAlsoGated`：
    `ResolveRole fallback onto a disabled first element = <nil>, want the disabled refusal (not a silent skip)`
  - `TestTicket261P1EnumerationSkipsAndSelectionRefusesDisabled`：
    `DiscoveredModels = [t261-nokey t261-off t261-on], want the disabled entry "t261-off" removed from discovery (ModelSpec.Enabled promise)`
  - `TestTicket261P1GateBlocksProviderAndCost`：`ResolveChain resolved a disabled model with the gate in`
  - `TestTicket261P1FlagInversionReversesEveryReading`：
    `cost-hop: the disabled entry produced 90 micros, want 0`（另有 provider-hop/枚举/选择三向红句，log 全存）
- **还原核对**：跑前 `md5 e8a2cbc85e2393e6b9ca9c8db849dc1b`＝跑后同值（跟踪 resolver.go 未被
  overlay 触碰）；`git status` 全程只有 r1 三件。
- ⚠ 注记：261-p1 的 `mut/resolver.gated.go`（前腿的**加门**突变）在我的终态跟踪树面前**已无差异意义**
  （门已在跟踪码里），故本腿不再重跑它；本腿的撤门形即上面的 ungated。

**恒真自查（两向，硬要求）**：`.scratch/wisp/probes/261/r1/mut/inverted/`
- 反形台：我的 3 枚 R1 用例**断言全部反写**（disabled 必须被枚举／必须被放行／回退必须成功）
  ＋ **p1 原版三枚**（从起手 HEAD `6c7a12b6` 原样抽出：`c259fad83d9478f45944a05f7fe0d5a9`），
  overlay 只换两枚测试文件、**不换 resolver**（对带门跟踪树跑）。
- 结果：`inv_rc=1`，**6/6 全红** ⇒ 尺对"门在不在"敏感，不是恒真尺。红句样例逐字：
  - `INVERTED: enabled entry "m-on" listed, want NOT listed`
  - `INVERTED: ResolveChain(disabled) = config: llm: model "m-off" of provider "p261" is disabled (enabled=false); re-enable it or remove the entry, want nil (admitted)`
  - p1 原版三枚（旧语义）对带门树全红＝改写后的用例确实测的是新语义，不是换名不改心。
- 跟踪树 md5 前后一致（pre/post-md5.txt diff 空）。反形台自身的构建错误一次（undefined
  gate261Resolver，helper 被 overlay 换掉了）已内联修复后重跑，log 为重跑版。

## §4 门禁读数（四数）

| 门禁 | 读数 |
|---|---|
| `go vet ./internal/llm/ ./internal/config/` | **rc=0**，输出 0 字节（`r1-final-vet.log`） |
| `sh scripts/d22scan.sh` | **rc=0 clean**（`r1-final-d22scan.log`：bans #1-5 internal/=228、cmd/=37、#6 frontend/=85、#7 internal/tools/=23、#8 design/=39、frontend/=85、internal/=498、cmd/=92，无违例） |
| `$(go env GOPATH)/bin/gofumpt -l`（三枚写过的 .go） | **空列表**（rc=0，零格式欠账） |
| `go test -count=1 -v ./internal/llm/ ./internal/config/` | **rc=0，155 PASS / 0 FAIL / 0 SKIP**，两包 ok（llm 59.037s／config 1.301s）。名册对账（comm，剥时延后缀）：丢名＝**恰 3 枚被改写的 p1 旧名**；增名＝**恰 6 枚**（3 改写新名＋3 枚新 R1）；config 侧两枚 `TestTicket261P1MissingEnabledKeyDecodesFalse`/`TestTicket261P1SettingsWriteMaterializesFalse` **原样在册且 PASS**（它们没被打红＝本腿改动零渗漏到 config 侧） |
| `go build ./...` | **⛔ 未跑**（派单禁区：本机还有别的测量排队） |
| 终态 `git status --porcelain internal/llm internal/config` | 提交后应**复归起手名册＝空**（本腿三件全部入库；突变台全在 `.scratch`） |

## §5 我判不动的地方（具名）

1. **面板快照那一跳**：`cmd/wisp/panel_config_store.go:112` 的 `for range p.Models { row.ModelCount++ }`
   现跑复认仍在（不过滤）。派单明令不碰（票 145 AC#2b 格子，编排者已另裁"先按目录全量"）。
   我只具名登记：**本门落了之后，面板 ModelCount 仍是"目录全量"，与"可选清单"的差从此开始存在**。
2. **验证层（`internal/config/catalog.go:26→65-71` `checkChainElement`）**：现跑复认只问
   "在不在目录"，不问 enabled——`LoadFile` 对点名 disabled 条目的 text_chain 依然放行。
   拒绝发生在 **resolver 解析时**（第一次真正要用那枚模型的时候），不在加载时。
   这是否够、要不要把拒绝前移到加载期＝**编排者裁**，本腿按派单只动 resolver 两处。
3. **`wisp providers discover` 子命令**（`cmd/wisp/providers.go:168` 起）：它走 resolver 是
   `ResolveChain` 形状（点名驱动），本门自然覆盖其点名路径；但 `DiscoveredModels` 的
   生产调用者＝零枚（p1 复认，我未重跑反证）——若后续有人接上枚举点，门已在。
4. **settings 写侧仍不写 enabled**：用户今天没有任何 UI/写侧路径能"重新打开"一枚手写 false
   条目（只能手编文件）。这不是本腿射程（写侧禁区），但门落之后"关掉"变得可感知了，
   "打开"仍只有手编一条路——具名留给编排者。
5. **配置侧那枚缺键落 false 的形状本身**：AC#2（defaults.go 的 map 分支）另有只读腿 `261-a2`，
   本腿零触碰（§1.4）。

## §6 推翻清单（票面与派单的每一句都是待验断言）

1. **派单"起手两包干净"**——成立，非推翻（我现跑 rc=0 零行，与宣称一致）。
2. **派单"枚举点实测在 :287-288"**——成立但精度注记：:287 是 `out := make(...)` 行，
   循环行是 :288（循环体 :288-290）。落门后终态行号变为 `:291` 起、循环 `:297-300`。
3. **派单预检清单"resolver.go:157 `func LoadFile`"**——**笔误**：resolver.go 里没有 LoadFile；
   真身＝`internal/config/loader.go:41`。我的预检 grep 按真实签名做，结论不变。
4. **派单"其中 catalog_test.go:217 逐字断言 `if !m2.Enabled || ...`"**——成立（现读逐字一致），
   但补一句射程注记：那是 `ImportDiscovered` 的判据（import 路径），不经 resolver，本门碰不到。
5. **`261-p1` verdict §1 读数①"ResolveRole(RoleChat)（run.go:436 同款回退）同收"**——成立；
   我补的精确形状：run.go:436 那一支是**未设角色的回退**（role.Provider=="" 时取第一链元素），
   它与 roles.chat 直配、text_chain 都汇入同一个 `resolveEndpoint`——所以本腿一处门关三条入口。
6. **票面 §现量-2"枚举那一处不读它"**——成立（起手读码复认），本腿之后**不再成立**（门已落）；
   verdict §3 表格里"resolver.go:287-288/119 不读"两行同理，终态起只对旧锚点为真。
7. **comm 比对的方法坑**（自翻）：p1 名册是裸名，`-v` log 抽出的名册带时延后缀，不剥后缀
   comm 全错——我在 §0 记了过程并保留了废件。后续腿拿 `-v` log 对裸名名册时必须先剥 ` (…s)`。

## §7 量不到的格子＋复量法

1. **真·外部计费**：与 p1 同界——本腿证明的是"disabled 条目连 resolver 都过不去（0 出网 0 计价）"，
   真实 provider 账户扣款不在任何测试射程内。复量法＝p1 verdict §5.1 同款：判据止于
   "请求是否到达 mockllm（路由自计数）＋生产计价函数读数"。
2. **`DiscoveredModels` 的下游行为变化**（GUI/命令）：该函数生产调用者＝零枚，无法量
   "某个用户可见清单变了"——只能量函数本身（本腿 §2 第一行）。若未来接线，复量法＝
   以 fixture provider 调 `DiscoveredModels` 断言集合。
3. **面板 ModelCount 与可选清单的差**：panel_config_store 不在本腿写面（§5.1），
   无尺可量；复量法归票 145 AC#2b 腿。
4. **加载期 vs 解析期的拒绝时机**：§5.2 记录了"验证层放行、解析期拒绝"的现状，
   但"这个时机对运维体验够不够"是产品判断，量不到，复量法＝若要前移，在
   `checkChainElement` 加 enabled 读并重跑本腿 §2 全套尺（尺已备好）。
