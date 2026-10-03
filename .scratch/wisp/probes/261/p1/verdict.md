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

## 1. 三向读数（AC#0 ⓐ，必须现跑）

（待填：跑 `internal/llm/enabled_reach_261_test.go` 与 `internal/config/enabled_261_test.go` 与 `.scratch/wisp/probes/261/p1/mut/` 突变台后逐条落读数。判据形状＝①枚路与选路给不给 false；②计费/真请求那一跳认不认它；③补回 true 同一把尺是否反转；④突变台证明尺是活的。）

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
| `go build ./...` rc | 0 | （待填） |
| `go vet ./internal/llm/ ./internal/config/` rc | 0 | （待填） |
| `sh scripts/d22scan.sh` rc | 0 | （待填） |
| `go test -count=1 ./internal/llm/ ./internal/config/` | 147 PASS / 0 FAIL / 0 SKIP | （待填，逐名 comm 对 `base-roster.txt`，不许丢名） |
| 终态 `git status --porcelain internal/llm internal/config` | 空 | （待填，必须与起手名册逐枚相等＝仍为空） |

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
