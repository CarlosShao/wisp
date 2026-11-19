# 248 补件 — reload 回执按档分列：[llm] 撤出「已立即生效」列（写腿交件，2026-11-19）

- 派单：`.scratch/wisp/probes/orch/20261119/txt/248-receipt-tiers-fix-dispatch.txt`（入库＝commit `abbde6219`）。
- 起手锚：`c2841a5fb`（进场时 HEAD）。
- 审计出处：`.scratch/wisp/probes/orch/20261119/txt/audit2-three-suspects-dispatch.txt` 第 2 条 ⇒ `docs/reports/pending-and-issues.md` §A「248 夸大 ⇒ AC#8 撤回勾」②：「『无一枚对 [llm] 承诺』被 `config_reload.go:219` 证伪（stdout 真把 `[llm]` 印进『已立即生效』名单，免责话只在 Go 串）」。
- 写面：`cmd/wisp/**`（本腿唯一共享编译面写腿）＋本件＋两份改准证据件（见 §5）；`internal/panel/**` 只动注释（`pump_queue_seq_167.go` 文件头，167 件的一半）。
- ⛔ 未动 thresholds/golden/allowlist/PLAN/specs/D43/design；⛔ 未动 `internal/config/tiers.go`（`"llm": TierHot` 的段级注册行一字未改——那行管的是"热应用进内存"，本腿修的是"stdout 怎么说"，两层各自真话）；⛔ push；⛔ UI emoji（本腿零 UI 面）。

## 0. 谎话本身（先量后写）

`git show c2841a5fb:cmd/wisp/config_reload.go:219` 逐字：

```go
fmt.Fprintf(rt.stdout, "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：%s\n", bracketed(split.claimable))
```

`splitHotTier` 的 claimable 集＝`audio llm agent cost`（`config_receipt_255_test.go:103` 旧断言钉的形状；实测读数 `.scratch/wisp/probes/orch/20261119/167-a-carrier-full.log:1282`：`IMMEDIATE "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[audio] [llm] [agent] [cost]"`）。⇒ 用户手改 `[llm]` 后看见的是「已立即生效（D36 立即档）」，而**设置页对 [llm] 任何写入字段的档位答案都是 restart**（`cmd/wisp/panel_config_store.go:556` 默认 `res.Tier = panel.EffectiveRestart` 盖住每一枚写字段；例外只有 `app.language`（问 TierOf）、`plugin_enabled`（locked 收紧答 now）、ball/panel 实测档——逐枚核过，全在 [llm] 之外）。同一句「要重启」在 `internal/panel/config_handlers.go:783` `tierSentence(EffectiveRestart)`＝「这一项要重启进程并重新运行才算用上（同一次运行里模型通路不会重建）。需要重启：本次运行不会读取新值。」。⇒ 两条回执对同一节答反，stdout 那条是假话。

255 名册裁定 `consumed`（live reader＝settings leg 每 call 现读，`config_readers_255.go` llm 行）没错——**错的是把"live reader 会答新值"说成"已立即生效"**：对 [llm]，live reader 答的是"下次 settings 读取会带新值"，而**模型通路仍按 boot snapshot 走到进程重启**（`config_readers_255.go` llm 行自己的 caveat：「已立即生效 here means "the live reader answers differently", not "the in-flight model swapped"」）。 caveat 只活在登记表注释里，stdout 那句没带它——这正是审计抓的形状（`pending-and-issues.md` §A524 同族判语：`config_readers_255.go:99-101` 的自认口径「话是对的，但它只活在登记表里、stdout 那句没带它」）。

## 1. 改法（按档分列，不是免责句）

三枚件：

1. **`cmd/wisp/config_readers_255.go`**：`hotSectionSplit` 增第四列 `restartSections []string` ＋ `restartKeys []string`；`splitHotTier` 在排完三列后调 `assignRestart248`——问写回执名册（`panel_config_store.go:556` 的盖档＋`writableField187` 名册），把**写档答 restart** 的段从 claimable 里**摘出来**归入 restart 列（今天答的就是 `[llm]` 一枚），`restartKeys` 带字段名（`llm.providers.<服务商>.base_url / api_key_ref / models.<模型>.context_window / models.<模型>.price.in/out / llm.roles.chat.model`——照 248/313 写名册的键路径形）。`llmRestartFieldNames` 是数据不是复述：回执句与测试共同消费它。
2. **`cmd/wisp/config_reload.go` `reportReload`**：claimable 句之后新增第四句——「`wisp run: 配置热加载：这些段已写进本进程内存（热应用），但它们的写入字段要重启进程才生效（D36 重启档）——本次运行里模型通路不会重建：[llm]。涉及字段：…`」。三件硬要求全落：**不是加长免责句**（[llm] 从 hot 句里物理移出、hot 句一字未加长）；**rest 档字段的提示没删**（诚实句「值已换进本进程内存…」原样保留）；**restart 字段名点了名**（涉及字段那半句）。ledger 半：`rep.Hot` 的 per-row 行原样保留，restart 列的每段追加一行 `config: HOT-RELOAD-READER section=llm write_tier=restart consumed: …`——consumed 判词不撤（live reader 是真的），只是 stdout 句归列。
3. **`cmd/wisp/config_receipt_255_test.go`**：**收紧**（见 §2）＋两枚新牙（见 §3）。

诚实句（no-reader）列一字未动；`[audio] [agent] [cost]` 留在 hot 句（它们的写档实测：audio 设备半/settings live 读、agent 下次任务、cost 下次任务——各有名册行背书，不在本腿改判范围）。

## 2. 断言收紧声明（派单硬约束要求的写清）

改了三条既有断言，**全部往严里改**（ⓐ 命中集变小 ⓑ 新增负断言 ⓒ 互斥断言），未放宽任何一条：

| 件:行（改前） | 改前断言 | 改后断言 | 为什么这是收紧 |
|---|---|---|---|
| `config_receipt_255_test.go:103` `TestTicket255SplitOnlyClaimsSectionsWithALiveReader` | `claimable == "audio llm agent cost"` | `claimable == "audio agent cost"` ＋ `[llm]` 必须在 restart 列（不在＝红）＋ 列互斥（同段双列＝红） | 允许集从 4 段缩到 3 段；旧形放过的 `[llm]` 从此必红 |
| 同文件 `TestTicket255ReceiptSentenceAssemblyIsFiltered` | immediate 句必须含 `[llm]` | immediate 句**禁止**含 `[llm]`（负断言）＋ restart 句必须含 `[llm]`+`base_url` ＋ 诚实句禁含 `[llm]` | 旧正向断言（=审计抓的那句假话）反转成负向；摘尺（整段消失）也红 |
| 同文件 `TestTicket255ReceiptStillNamesTheLiveReadSection` | 「已立即生效」句必须含 `[llm]`（=被证伪的那句假话的正控） | 改成：restart 句必须含 `[llm]`+`api_key_ref`；immediate 句含 `[llm]`＝红；live-reader 实证半（真 settings leg 答新值＋boot snapshot 不动＋ledger `consumed`）**逐字保留** | 反控的"谁在 live 读"证据没动，动的只是"它在哪一列"；删段＝摘尺必红 |
| `quiet` 计数行（同文件 `:143` 域） | `len(quiet) == len(hot)-4` | `len(quiet) == len(hot)-3-len(restartSections)` | 分母没变，账目对齐新四列；漏记 restart 列＝红 |

## 3. 牙（会响的，两枚新＋两枚旧针改向）

- **`Test248ReceiptRestartColumnIsExclusivelLmWithFieldNames`（新，纯判定面）**：把 `[llm]` 放回 hot 列（=拆掉 `assignRestart248` 的降列）⇒ `THE TICKET-248 LIE: [llm] appears in the 已立即生效 column` 红；restart 列字段清单缺任一枚（base_url/api_key_ref/context_window/price/roles.chat.model）⇒ 红。
- **`Test248ReceiptSentencePutsLlmInRestartColumnNotHot`（新，production `reportReload`）**：直接喂生产函数全 hot 名单，抓 stdout 三句：hot 句禁 `[llm]`、restart 句必含 `[llm]`+`base_url`+`api_key_ref`+`context_window`+`roles.chat.model`+逐字 `需要重启进程`、诚实句禁 `[llm]` 且仍带 `[panel] [ball]`、ledger 行 `section=llm write_tier=restart` 仍 `consumed:`。
- **`TestTicket255ReceiptSentenceAssemblyIsFiltered`／`TestTicket255ReceiptStillNamesTheLiveReadSection`（改向）**：见 §2。
- 复跑命令（载具，⛔ 裸跑 A1888）：`sh scripts/wisp-cli-staged-tests.sh -test.run 'TestTicket255SplitOnlyClaimsSectionsWithALiveReader|Test248ReceiptRestartColumnIsExclusivelLmWithFieldNames|TestTicket255ReceiptSentenceAssemblyIsFiltered|TestTicket255ReceiptStillNamesTheLiveReadSection|TestTicket255ReceiptOmitsPanelFromTheImmediateSentence|TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere|Test248ReceiptSentencePutsLlmInRestartColumnNotHot|TestTicket223RestartTierSaysItWillNotApply'`（读数见 §6）。

## 4. 变异读数（M1 一发，指名必红）

| 变异 | 改什么 | 期望红 | 实测 | 还原后 |
|---|---|---|---|---|
| **M1** | `config_readers_255.go` `assignRestart248` 返回空表（`[llm]` 回 hot 列） | §3 两枚新牙齐红：`THE TICKET-248 LIE: [llm] appears in the 已立即生效 column` ＋ `Test248ReceiptSentencePutsLlmInRestartColumnNotHot` 的 `THE TICKET-248 LIE IS BACK: the 已立即生效 line names [llm]` | **RUN=5 PASS=0 FAIL=5 SKIP=0 rc=1**，指名红句逐字 `config_receipt_255_test.go:1185: THE TICKET-248 LIE IS BACK: the 已立即生效 line names [llm] whose write fields answer restart: "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[audio] [llm] [agent] [cost]"`（假话复活被原句抓住） | 复跑同一 -run **RUN=5 PASS=5 FAIL=0** rc=0 |

## 5. 同批两份证据件改准（本件只记指针）

- `docs/evidence/s1/248-ac8-20261119-r1.md`：①autostart_windows.go 三处不是 `立即生效` 命中（该词在那件 0 命中，原文是「需要重启后重试」）②普查行数改准（审计尺＝25 行非测试命中、旧件漏列 19 行、本件补列全 25 行）③「8 check」改 9 处 `check(`／9 枚牙、`config_route_248_test.go:392` 改 `:391`、普查尺补 `-E`、删除「没有一枚是对 [llm] 的承诺」绝对句（已被 `config_reload.go:219` 证伪，指向本件）。
- `docs/evidence/s1/167-occupancy-queue-a.md`：新增「现状（427dd7955 之后）」一节＋§1/§2/§3 逐节取代标注；`internal/panel/pump_queue_seq_167.go` 文件头补两套编号器语义差注释。

## 6. 门禁读数（全部现跑，2026-11-19 18:1x–18:4x）

| 尺 | 读数 |
|---|---|
| 跑前 `Get-Process wisp.test` | 无（每发载具前现查，空读数） |
| 载具①绿（`sh scripts/wisp-cli-staged-tests.sh`，`-test.run` 九针＝§3 全名单＋`TestTicket223RestartTierSaysItWillNotApply`＋`TestTicket255RestartTierKeysAreBackedByATest`，等号形旗标） | **RUN=9 PASS=9 FAIL=0 SKIP=0**，`ok github.com/CarlosShao/wisp/cmd/wisp 46.00s`，rc=0。台件 `.scratch/wisp/probes/orch/20261119/248-receipt-tiers-carrier-final-nine-green.log` |
| 变异 M1（同 -run 五针） | **RUN=5 PASS=0 FAIL=5 SKIP=0**，rc=1；指名红句＝`config_receipt_255_test.go:1185: THE TICKET-248 LIE IS BACK: the 已立即生效 line names [llm] whose write fields answer restart: "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[audio] [llm] [agent] [cost]"`＋`:117 restartSections = [], want exactly [llm]`（五牙齐响）。台件 `…248-receipt-tiers-mutation-M1-red.log` |
| 还原复跑（同五针） | **RUN=5 PASS=5 FAIL=0 SKIP=0**，`ok … 11.00s`，rc=0。台件 `…248-receipt-tiers-carrier-restore-green.log` |
| 载具②（`^Test167*` 两枚＋`Test233ReloadTierReceiptOnlyForReloadTier` 回归保险） | **RUN=3 PASS=3 FAIL=0 SKIP=0**，rc=0。台件 `…248-receipt-tiers-carrier-167-233-green.log` |
| `go test ./internal/panel -run '^Test167'` | `ok github.com/CarlosShao/wisp/internal/panel 0.064s`，rc=0（注释改动后复跑） |
| `sh scripts/d22scan.sh` | `d22scan.sh: OK - control rc=0, scan rc=0`；末行 `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=422, bans #1-5 cmd/=107, ban #6 frontend/=188, ban #7 internal/tools/=52, ban #10 internal/risk/=21, ban #8 design/=39, ban #8 frontend/=188, ban #8 internal/=1066, ban #8 cmd/=359…`（rc=0，0 findings）。台件 `…248-receipt-tiers-d22scan.log` |
| `go build ./...` | **rc=0** |
| `gofmt -l`（四枚动过的 Go 件） | rc=0，list 空 |

⚠ 首两发载具调用 argv 被 pwsh 撕参（`flag provided but not defined: -test`：空格旗标 `-test.v` 被撕成 `-test`＋`.v`）＝**参数传递 setup failure，不是测试判决**，全等号形（`-test.v=true`）复跑。中途两发真红（run3＝牙要的「需要重启进程」与初稿句差一字，改句就家规词；run4＝`:140` 循环未认 restart 列＋`StillNames` 等待句错选）均为收口前迭代，终形读数以上表为准。
