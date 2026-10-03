# 票 261 · 腿 261-p1 裁决表（AC#0 ⓐ：有没有一条真路径能让 `enabled=false` 的模型被选中并产生花费）

代理：`261-p1`（写码子代理，只加仪器，⛔ 生产码一字不动）
票面：`.scratch/wisp/issues/261-the-model-enabled-key-promises-removal-from-discovery-and-selection-but-no-production-code-reads-it-while-the-default-true-tag-is-inert-for-map-entries.md`

---

## 0. 起手锚与起手名册

- 起手时刻：`2026-10-03T10:54:28+08:00`（`date -Iseconds` 自取）
- 起手 HEAD：`fe0ff659d90c8d988709761cdfe28d21114d012d`（`git log -1 --format=%H` 自取），分支 `dev`
- 起手名册：`git status --porcelain internal/llm internal/config` → **空**（rc=0，零行输出 ⇒ 两包工作树干净，无他人未提交件）
- 起手基线（⛔ 未跑测试前的 rc 都单独取，不用管道）：
  - `go build ./...` rc=0（空输出，`base-build.log` 0 字节）
  - `go vet ./internal/llm/ ./internal/config/` rc=0（`base-vet.log` 空）
  - `sh scripts/d22scan.sh` rc=0（`base-d22scan.log`）
  - `go test -count=1 -v ./internal/llm/ ./internal/config/` rc=0；
    `--- PASS` 顶格 **147** 枚、`--- FAIL` **0**、`--- SKIP` **0**；两包 `ok`（llm 46.841s／config 1.334s）。
    名册逐名存 `base-roster.txt`（147 行），收尾 comm 对比。
- 提交：骨架件＝`29047c99`（verdict 骨架＋基线四 log＋名册）；终态件＝本腿最后一次 commit
  （两枚测试文件＋§1/§4 全部凭据 log＋`mut/` 突变台三件），`git log -1` 即得。

## 1. 三向读数（AC#0 ⓐ，必须现跑）

仪器：`internal/config/enabled_261_test.go`（2 枚尺）＋`internal/llm/enabled_reach_261_test.go`（3 枚尺）＋
`.scratch/wisp/probes/261/p1/mut/`（overlay 突变台，⛔ 未动任何跟踪文件）。全部真跑：`config.LoadFile`
真解析手写 TOML、`llm.NewResolver`→`ResolveChain`/`ResolveRole`→`BuildEndpointProvider`→live mockllm HTTP、
计费用生产函数 `internal/agent/cost.go` 的 `Cost.AddUsage` 本身（llm 测试以外部测试包 `llm_test` 引入 agent，
无环）。终态名册＝基线 147＋新 5，0 FAIL／0 SKIP／0 丢名（§4）。

1. **读数① 枚举与选择认不认 false**：`DiscoveredModels("mock261")` 返回 `[t261-nokey, t261-off, t261-on]` 全部三枚
   （显式 `enabled=false` 与缺键落 false 两形都给）；`ResolveChain` 对 `["mock261/t261-off"]` 返回
   `err=nil` 且 `eps[0].Model=="t261-off"`；`ResolveRole(RoleChat)`（run.go:436 同款回退）同收。
   **＝认。**（`TestTicket261P1EnumerationAndSelectionAdmitDisabled` PASS）
2. **读数② 派发/计费那一跳认不认 false**：生产装配形状把该端点构造成 provider 后 `Stream` 打真 HTTP：
   mockllm 自己计数 `chat` 路由＝**1 次**（非我方记账），录制 transport 抓到的出网请求体含逐字
   `"model":"t261-off"`；回合 usage＝12 in／6 out；`agent.Cost.AddUsage(该 false 条目自带价卡, usage)` ＝
   **90 micros > 0**（日志行 `TICKET261-COST-HOP READING`）。
   **＝认（一条可计费的真实 dispatch 今天对 false 照常发生）。**
   （`TestTicket261P1DisabledModelReachesProviderAndCost` PASS）
3. **读数③ 补回 true，同一把尺反不反**：同一枚条目文件里只翻 `enabled = false`→`true`、重跑全套尺：
   listed／resolved／reachedProvider／micros 四项**逐项不变（90 对 90），不反转** ⇒ `Enabled` 在生产码里是死键。
   同测里的活体对照：目录中不存在的 `t261-ghost` 在**同一行选择代码**被拒（`unknown model` 错误）
   ⇒ 尺对"在不在目录"敏感，只对 `enabled` 不敏感——不是钝尺。（`TestTicket261P1FlagInversionChangesNoReading` PASS）
4. **读数④ 突变台（尺的标定，只在 .scratch）**：`go test -overlay` 把 `resolver.go` 换成
   `mut/resolver.gated.go`（仅在 `resolveEndpoint` 加 `if !spec.Enabled {拒绝}`＋枚举跳过 false），
   同一套尺三枚全反：枚举只回 `[t261-on]`、`ResolveChain(t261-off)` 报
   `model "t261-off" ... is disabled (mutation bench)`、cost-hop `off=0 / on=90`
   （`mut/mut-run.log`，mut_rc=1＝预期红）。⇒ 读数③的"不反转"确因**门不存在**，恒真嫌疑排除。
5. **config 侧配套读数**：手写文件缺 `enabled` 键 → 解码 `false`（tag `default:"true"` 对 map 条目零执法，
   `TestTicket261P1MissingEnabledKeyDecodesFalse` PASS）；一次真实 `SetModelContextWindow` 写后，
   文件里出现逐字 `enabled = false`（mergeWrite→SaveFile→MarshalCanonical 无 omitempty 的物化放大），
   且 ghost 条目写入被 `refusing to invent one` 拒（`TestTicket261P1SettingsWriteMaterializesFalse` PASS）。

**ⓐ 判定（跑出来的）**：**真路径存在。** 形状＝手写（或任何缺键）条目 → `config.LoadFile` 落 `false` →
`text_chain`／`roles.chat` 点名它 → 验证层放行（`catalog.go:66` 只问存在）→ 枚举给出、选择解析、
provider 真收带该 model id 的计费请求、计费函数照算价卡。全程**没有任何一处读 `Enabled`**。
生产入口逐枚名册（谁读/谁不读）：发现＝`ImportDiscovered` 写 true 但**生产调用者零枚**；
选择＝`resolver.go:119/287-288` 不读；派发＝`cmd/wisp/run.go:435-442`、`providers.go:168-175` 不读；
计费＝`internal/agent/loop.go:435`→`cost.go:38/50` 不读；回执＝`cmd/wisp/panel_config_store.go:112`
（`for range p.Models { row.ModelCount++ }`）不读。**ModelSpec.Enabled 生产读者＝零枚（复认＋现跑双证）。**

## 2. 「手写那一支」今天可达不可达（逐枚带 file:line）

三处逐枚答（读码复认于起手锚；引文均本人 Read/grep 过）：

1. **用户手写 config.toml** —— **可达（这就是那条路本身）**。
   `internal/config/loader.go:117-118`（`readConfigFile` 里 `cfg := NewDefaults()` → `decodeStrict(raw, cfg)`）：
   解码基底只有段级默认（`internal/config/defaults.go:77-78` `case reflect.Map:` 只留
   `// leave nil (see NewDefaults)`，不进 map），所以 `[llm.providers.<名>.models.<id>]` 里
   **不写 `enabled` 就是 `false`**。验证层不设防：`internal/config/validate.go:195-209`
   （`validateModelSpec` 只查 billing/thinking_levels 两枚枚举，不碰 enabled）、
   `internal/config/catalog.go:66`（`checkChainElement` 对 text_chain 只问模型**在不在目录**，不问 enabled）。
   另有持久化放大器：`internal/config/parse.go:160-163`（`MarshalCanonical` 注释逐字
   "every static field emitted explicitly (no omitempty …)"）⇒ 任何一次 `mergeWrite`
   （`internal/config/writeguard.go:160` 的 `SaveFile(m.path, base)` 整文件重排）会把
   `enabled = false` **写进文件实体**，从那天起它就是一个显式的 false 了。
2. **票 257 的首启引导（`cmd/wisp/firstrun.go`）** —— **可达，且它是官方指路的入口**。
   `cmd/wisp/firstrun.go:82` 首建只写 `config.SaveFile(cfgPath, config.NewDefaults())`（**零模型**，map 为 nil）；
   随后 `cmd/wisp/firstrun.go:116-118` 打印的引导原话让用户
   "在同一份文件的 [llm.providers.<名>] 里补 … models.<id> … 再在 [llm] 的 text_chain 或 roles.chat 里点名 provider/model"
   ——**文案通篇不提 `enabled`**。产品教用户走的就是第 1 枚那条手写路。
   `wisp providers discover`（`cmd/wisp/providers.go:116-149`）只列不发不写库：它打印 `/v1/models` 的结果，
   不调用 `ImportDiscovered`；全仓 grep `ImportDiscovered` 的**生产调用者＝零枚**（只有测试）。
3. **设置页写侧 `internal/config/settings.go:115/153`** —— **不可达（造不出新条目），但它会把 false 坐实**。
   两处逐字均为 `spec, ok := p.Models[model]`，且都被 `requireCatalogEntry`
   （`internal/config/settings.go:294-308`）挡着"refusing to invent one from a settings write"
   ⇒ 写侧**不产行**；`settings.go` 全文也没有任何 setter 写 `enabled` 键（写不了 true 也写不了 false）。
   但它经由 `writeOneKey → mergeWrite → SaveFile`（`internal/config/settings.go:242`、
   `internal/config/writeguard.go:160`）整文件重排，按第 1 枚末段的放大器效应，
   任何一次设置页保存都会把手写缺键条目的 `enabled = false` 物化进文件。

**小结（读码侧）**：今天造出 `enabled=false` 条目的路径＝手写（firstrun 官方教的就是它）；
发现器与迁移写侧都显式写 true（`internal/llm/discover.go:119`、`internal/config/migrate.go:146`），
而读侧零枚读者（见 §3 复认）。真路径成立与否以 §1 现跑读数为准。
**现跑已证（2026-10-03 终态）**：第 1 枚的解码落 false＋验证放行＝§1 config 尺两枚 PASS；
第 2 枚的"文案不提 enabled"＝生产引导原文（firstrun.go:111-118，无 SaveFile 之外的建行写者，
`ImportDiscovered` 生产调用者零枚，全仓 grep 复认）；第 3 枚"不产行但物化"＝§1 读数末条
（写后文件逐字出现 `enabled = false`，ghost 写被拒）。三枚的可达性判定与读数一致：**手写那一支今天可达。**

## 3. 我对编排者现量节的复认与推翻清单

逐枚亲验（Read 全文＋grep 字面双做），**无一枚推翻**，两处行号精度注记：

| 断言 | 结果 |
|---|---|
| `internal/llm/resolver.go:287-288` 逐字 `out := make([]string, 0, len(p.Models))` / `for id := range p.Models`，循环体无 Enabled 判据 | **成立**（grep 命中 287/288 行） |
| `internal/llm/resolver.go:119` 逐字 `spec, ok := p.Models[model]`，只问在不在 | **成立** |
| `internal/config/defaults.go:75-80` `case reflect.Map:` 只留 `// leave nil (see NewDefaults)` | **成立**（注释实际在 **:78**，分支 :77-78，落在票面给的范围里） |
| `internal/config/schema.go:348` 起 `type ModelSpec struct {`；末段承诺句＋`Enabled bool \`toml:"enabled" default:"true"\`` | **成立**；承诺句在文件里**折成两行**（:369 `// Enabled removes the model from discovery/selection without deleting` ／ :370 `// its catalog entry.`），字段在 **:371**；票面把它抄成一行是拼接转写，语义无差 |
| `internal/llm/discover.go:119` 逐字 `Enabled: true,`（构造体 :117 起） | **成立**（行号即 :119，票面补读时已自纠 :117→:119 漂行） |
| `internal/config/settings.go:115/153` 两处 `spec, ok := p.Models[model]` | **成立** |
| 口径：`ModelSpec.Enabled` 生产读者＝零枚（全仓 `\.Enabled` grep，排除 _test.go） | **成立**。命中全部属别的结构：`cmd/wisp/logsink.go:205/210`＋`internal/observe/logging.go`（slog `Enabled(ctx,…)` 方法调用）；`internal/projctx/projctx.go:209/255`（自己 Loader options 的字段）；`internal/memory/dao_misc.go:199/236/260`（**plugin_state** 表结构，`sed -n '190,205p'` 亲验其上文是 `plugin_state` INSERT）；`internal/config/manager.go:380/394/395/512`（Voice/WakeWord）；`internal/config/catalog.go:126`（`rt.Enabled`＝Voice.Realtime）。测试文件里的 `Enabled` 读/写（catalog_test/mockllm_integ_test/migrate_test）不算生产读者 |
| 增量（票面没问，但属"谁读"名册）：`Resolver.DiscoveredModels` **生产调用者＝零枚**；面板清单那一跳 `cmd/wisp/panel_config_store.go:112` `for range p.Models { row.ModelCount++ }` 同样不读 Enabled；`cmd/wisp/providers.go:193`（probe 读 Declared）也不读 | **现量补充**，⛔ 不改本票性质 |

**AC#0 ⓑ 复认＋判够不够**：发现器填 true（discover.go:119）与迁移填 true（migrate.go:146）使
"缺键落 false"只剩手写一支——但**今天没有任何生产调用者会走到发现器的建行函数**
（`ImportDiscovered` 生产零调用），所以"写侧有 true 执法者"这半句实际是纸面的；
够不够＝**不够**：手写一支是产品亲自指路的默认路径（§2 第 2 枚），而读侧零读者。

## 4. 门禁读数

| 门禁 | 基线（起手） | 终态 |
|---|---|---|
| `go build ./...` rc | 0 | **0**（`final-build.log` 空） |
| `go vet ./internal/llm/ ./internal/config/` rc | 0 | **0**（`final-vet.log` 空） |
| `sh scripts/d22scan.sh` rc | 0 | **0 clean**（`final-d22scan.log`：bans #1-5 internal/=228、cmd/=37、#6 frontend/=85、#7 tools/=23、#8 含 _test.go 与注释 internal/=496、cmd/=92，无违例） |
| `go test -count=1 -v ./internal/llm/ ./internal/config/` | 147 PASS / 0 FAIL / 0 SKIP，两包 ok | **152 PASS / 0 FAIL / 0 SKIP，两包 ok**（llm 43.365s／config 1.046s，rc=0）。`comm -23 base-roster.txt final-roster.txt`＝**0 行**（基线 147 名逐名未丢）；增量恰为 5 枚 `TestTicket261P1*`（§1 所列）。判红绿只认顶格 `--- FAIL`（0 枚）；`t.Logf` 的 `file:line:` 前缀行不当失败读 |
| 终态 `git status --porcelain internal/llm internal/config` | 空 | 见末段（提交本腿两枚测试文件后必须复归空；突变台只活在 `.scratch`＋`-overlay`，跟踪文件从未被改） |
| overlay 突变台 | —（不属于门禁） | `mut_rc=1`＝预期红（§1 读数④），发生在 `go test -overlay` 的合成输入上，⛔ 不落跟踪树 |

红绿判读备注：iter-config/iter-llm/iter-llm2 三份是开发回路（iter-llm 首跑的红是
mockllm `chat` 路由不录请求体的仪器事实——`tools/mockllm/server.go` 注释逐字 "the chat route is untouched"，
改在测试侧加录制 RoundTripper 解决，⛔ 未动 tools/）。

## 5. 我判不动的地方（具名，不含糊）

1. **真·外部计费**：读数只能证明"带 model id 的 completion 请求真的出网到达端点、
   真实计费函数 `internal/agent/cost.go:38/50` 对 false 条目的价卡照算不误"；
   **真实 provider 账户扣多少钱**不在任何测试射程内（⛔ 也不该去问真端点）。判据止于"请求到达＋计价器无门"。
2. **AC#1 ⓐ/ⓑ 选形**：加不加过滤、还是降级承诺——编排者裁，本腿⛔ 不碰生产码也不选边。
3. **AC#2 全仓 map 字段普查**（还有几枚带 `default` 的 map 内字段、缺键落零值会不会被读出意义）——
   不属本腿 AC#0；本腿只在读数里顺带钉了 `ModelSpec.Billing` 缺键落 `""`
   且 `validateModelSpec`（validate.go:197 `m.Billing != ""`）放行这一形状，供 AC#2 那腿引用。
4. **一处本腿读出来但票没问的接线缺口**：`cmd/wisp/run.go:1010-1017` 装配 `agent.Config` 时
   **从不填 `Price`/`Currency`**（grep `Price` 于 run.go 零命中），而计费那一跳读的就是它
   （`internal/agent/loop.go:435`）⇒ 今天 `wisp run` 对**所有模型**（含 enabled=true）本地成本估计恒 0，
   `costLine`（run.go:1197）恒打"（未计价）"。这不改变"false 模型可被选中并出网产生真实花费"的判定
   （花费发生在 provider 侧），但"计费那一跳今天对谁都不执法"这句要说全。归属：C23/ticket 44 接线，
   与本票的 `enabled` 缺口**同族不同格**，⛔ 本腿不动、只具名上报。
5. **面板快照的"回执"维**（票 145 AC#2b 的"∩ enabled=true"）：本腿只读引用了
   `panel_config_store.go:112` 的不过滤事实；快照形状本身归编排者已裁的临时口径（"先按目录全量、不承诺交集"），⛔ 不判。
