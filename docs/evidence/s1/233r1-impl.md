# 票 233 落地腿 `233-r1` — D36 三档标到键上（三常量化）＋ reload 档生产接线 ＋ 热档活引用

- 起手锚＝`6c7733a42`（2026-10-27 19:1x +08，reflog 回溯：本腿第一次产码动作前 HEAD＝`6c7733a42`，其后编排者在我读数期间落了 `75900bf9b`…`1bb1ae961` 七枚 commit、与代码面无碰撞——`git show --name-only` 逐枚核对无 `internal/config|cmd/wisp` 交集）。起手核：`git diff --cached --name-status`＝**空**；`git status --porcelain -- internal/config cmd/wisp`＝**空**（共享树只有 `.scratch/**`、`design/**`、`logs/**` 他人在飞件＝零触碰）。
- 写面（6 改＋1 新）：`internal/config/tiers.go`、`internal/config/schema.go`、`cmd/wisp/config_reload.go`、`cmd/wisp/run.go`、`cmd/wisp/config_readers_255.go`、`cmd/wisp/config_receipt_255_test.go`、新 `cmd/wisp/config_tiers_233_test.go`。**`internal/risk|speech|session|audio|tools|secret|memory`、`testdata/**`、`frontend/**`、`design/**`、`scripts/**`、`.github/**`、票面、台账、`docs/PLAN.md`、`docs/specs/**` 零字节**（`git diff --name-only` 逐目录空，§7）。
- 日志：`D:\tmp\wisp\20261026\233r1\logs\01..28`（只建不删）。

## 0. 现状对派单的修正（开工时核到的事实，派单三处过时）

| 派单说 | 现读（起手现跑） | 处置 |
|---|---|---|
| 「Tier* 三常量非测试命中只在 schema.go:34-43（零赋值点）」 | **为真**：三常量全仓非测试引用＝0，真正的分档表是 `internal/config/tiers.go` 的 `TierRegistry`（票 255，commit `dd92bb924`）——但它的值是**裸字符串** `"hot"/"reload"/"restart"/"locked"`，常量依然零消费 | 本腿把表**换型** `map[string]Tier`、值改四常量（§1） |
| 「`OnReload` 在 wisp 里 0（唯一赋值点 balldebug:244）」 | **为真**：`grep -rn "OnReload = " --include=*.go cmd/ internal/` 非测试＝`cmd/balldebug/main.go:244` 一枚 | 本腿接上 `wisp run` 生产消费者（§3） |
| 「热档：`run.go:460` cfg 是 snapshot（frozen at boot）⇒ 行为未变」 | **为真**：`execute()`（run.go:1252 `cfg := rt.cfg`）四个读点全走 boot copy；227 AC#6 把 `rt.cfg` 钉成冻结"before"（`:170-174`），255 把 [agent] 钉成 `snapshot-only` | 本腿**不动 `rt.cfg`**（它自己的钉子与 223/227/255 三族用例都按"冻结"读它），把 **`execute()` 的 cfg 源**换成每任务现读 `mgr.Config()`（§2）——热的真语义＝「下一次用它的行为用新值」 |
| 票面 AC#2/AC#4 那格「⛔ 不许拿已归口 233 去翻 AC#2 的勾」 | 本腿不翻任何勾（翻框禁令）；reload 档的"接"半格由 §3 闭合，**判勾归验收腿** | — |

## 1. 读数①②：票面 6 格逐条 ＋ 逐键 tier 标注表

**票面 6 格**（现量 5 行＋照三行读）：①三档枚举存在 ✓（schema.go TierHot/Reload/Restart）②"零键被标"半过时：键已由 255 标进 `TierRegistry`，但**标的是裸串、不是那三枚常量**（三常量零赋值点）——本腿闭合这半 ③reload 无消费者 ✓（本腿闭合）④reload 无生产键 ✗ 已过时：`planVoice` 早把 `voice.*` 填进 `rep.Reload`（manager.go:421-423），是**事件没人收** ⑤硬编码 `restartTierKeys` 在 cmd/wisp ✓（255-r2 已有常驻钉 `TestTicket255RestartTierKeysAreBackedByATest`，一字未动）⑥`internal/ball` 同名钩子非消费者 ✓（`hotkey_reload.go` 注释自陈，零触碰）。照三行读：现象在"档位元数据与三常量脱钩＋事件落空"，无入侵证据，最坏形状＝票 223 AC#2 禁的"静默不生效"在 reload 档复活。

**逐键 tier 标注表**（改后＝`TierRegistry` 全量 44 行；值全部是常量；与 `PLAN.md:2738-2757` 那张 D36 表逐行对得上；`TierOf` 仍返回 string 视图，调用面零动）：

| 键 | 档（常量） | D36 表依据 |
|---|---|---|
| `ball`/`session`/`audio`/`llm`/`agent`/`privacy`/`memory`/`panel`/`cost`/`models`/`observe`/`hotkey` | `TierHot` | D36 表各节标 `hot`（hotkey 注"须重注册热键"、audio"须重开设备"＝读侧语义，档位不变） |
| `risk`/`fs`/`net`/`plugins` | `TierLocked`（新增枚举值） | D36 表 🔒"变更即需 L2 重新确认"＝`planLocked`，非三档 |
| `app.theme` | `TierHot` | D36 `[app]` 行"`theme` 为 `hot`" |
| `app.language`/`app.autostart`/`app.single_instance` | `TierRestart` | D36 `[app]` 行标 `restart`（＝`restartTierKeys` 那 3 枚） |
| `voice.enabled`/`voice.wake_word.enabled`/`voice.wake_word.keywords`/`voice.asr{,.provider,.model}`/`voice.tts{.provider,.voice}`/`voice.conversation_mode`/`voice.aec{,.enabled,.echo_ref}`/`voice.realtime{,.enabled,.provider,.model,.api_key_ref,.base_url}`/`voice.cloud_asr_chain`/`voice.cloud_tts_chain` | `TierReload` | D36 `[voice]` 行"`reload`（换模型）" |
| `voice.wake_word.thresholds`/`voice.wake_word.veto_words`/`voice.tts.speed`/`voice.punctuation` | `TierHot` | D36 `[voice]` 行"`hot`（阈值、否决词）" |

常驻钉子全绿自证：`internal/config` 四枚 tier 尺（`tiers_255_test.go` 3 枚＋`tiers_app_255r2_test.go` 2 枚方向）＋`tiers_223r3_test.go` 四枚＋`restart_tier_keys_255r2_test.go`——**值语义逐行不变、只换型**，`internal/config -count=3` rc=0（§6）。

## 2. 读数③：热档真热一发（改 → 不重启 → 行为真变了）

`cmd/wisp/config_tiers_233_test.go::Test233HotTierAgentReachesTheNextTask`（新）：host 起来后等 tick `state=armed`、核前提（boot snapshot `PerToolTimeoutMS==0`）→ **手改 config.toml 追加 `[agent] per_tool_timeout_ms = 2500`** → 等 tick 自己的判定行 `hot=[agent]` → 断言 ①同进程**下一个任务**的 loop 装配 `rt.loopOpt.Config.PerToolTimeout == 2500ms`（＝execute 自己发布的 agent.Options，改后活引用读到的）②**boot snapshot 仍是 0**（中程改动的分界）。读数（log 04）：**PASS**；同尺 `TestTicket223RunArmsTheReloadTick`（ball.size 活引用读）仍 PASS。读侧行＝run.go:1263 `cfg = rt.mgr.Config()`；`rt.cfg` 字段与其冻结语义一字未动。**热档依据（为什么这是"真热"不是硬凑）**：`per_tool_timeout_ms` 是 D36 表 [agent] 行枚举的键、255 名册 verdict 从 `snapshot-only` 升 `consumed`（§5），行为读点＝本进程真实下一个任务的门控预算，非探针字段。

## 3. 读数④：reload 档明确告知一发（改 → 不重启 → 逐字日志）

`Test233ReloadTierEventReachesAProductionConsumer`（新）：plant `[voice.asr] model = "paraformer-zh-233"` → tick 判定 `config: HOT-RELOAD state=applied hot=[] reload=[voice] restart=[] locked=0` → **生产链路** `OnReload → reportReloadTier`（cmd/wisp 装配，非测试赋值）。**stdout 逐字**（windowed，`mark` 后增量）：`wisp run: 重载档（[voice]）的改动已生效到本进程内存，但本宿主没有可重载的子系统（语音模型换装在常驻进程里），本次运行不会触发模型重载；需要真正换装请重启常驻进程。`；**audit 逐字**：`config: RELOAD-TIER state=received sections=[voice] consumer=wisp-run effect=this-host-reloads-nothing detail="重载档的改动已写进本进程内存、事件已送达；但本宿主（wisp run）没有可重载的子系统——语音链在常驻进程 wisp 里装配，本宿主不做模型换装，新值由下一次用它的组件从配置层现读。要真正重载子系统请重启常驻进程"`。同时断言内存真应用（`Config().Voice.ASR.Model == "paraformer-zh-233"`）＋ `TierOf("voice.asr.model")=="reload"`。读数（log 04/05/11）：**PASS 3/3 发**。

## 4. 读数⑤：重启后档复跑

`TestTicket223RestartTierSaysItWillNotApply`（232/272 窗口化原尺，零改动）本腿复跑 **4 发全绿**（log 01 一发 16.7s 尺、log 03/04/25 三发）：`stdout never carried` 红零出现，`state=restart-pending sections=[app]`＋`effect=next-process-start`＋stdout 五 needle（`app.autostart`/`开机自启`/`需要重启进程`/`原因：`/`交给平台层`/`没有被丢掉`）逐字在。`restartTierKeys` 与其钉子零触碰。

## 5. 读数⑥：反向判据（突变两发，全在仓外导出树 `-overlay` 不需要——直接 `git archive HEAD`＋覆盖 6 件）

- **M1＝摘 `OnReload` 赋值**（export-mut 树，`config_reload.go` 注释包住赋值行）：`Test233ReloadTierEventReachesAProductionConsumer` **FAIL 41.90s**，红句逐字 `config_tiers_233_test.go:90: the audit trail never carried "config: RELOAD-TIER state=received" AFTER the plant within 40s`，dump 里 `state=applied hot=[] reload=[voice]` 在而 RELOAD-TIER 行不在（＝事件进了无人线）；**对照不发红**＝`Test233ReloadTierReceiptOnlyForReloadTier`（hot 档）同尺 PASS 1.96s（log 07/07b）。
- **M2＝摘 execute 活引用**（`cfg = rt.mgr.Config()` 注释回 snapshot）：`Test233HotTierAgentReachesTheNextTask` **FAIL 1.94s**，红句逐字 `config_tiers_233_test.go:193: the task that ran in this process was built with PerToolTimeout=0s, want 2.5s (the hot-applied [agent] edit did not reach execute()'s own inputs - 改了 → 不重启 → 行为没变)`（log 10）。
- **还原核（时机＝M2 之后、最后一发变异后）**：两文件 SHA256 前后逐字同（`config_reload.go` `EED0DEA2…DAE5DC`、`run.go` `46F5584F…83BAD6`）、`MUTANT` 标记 0 命中、三枚 233 用例原样复跑 **ok 5.8s**（log 08/09/11）。仓内树全程未被动（`git status -- cmd internal` 逐刻空）。

## 6. 读数⑦：零回归

| 对象 | 改前（HEAD） | 改后 |
|---|---|---|
| `internal/config` 包级 `-count=1`/`-count=3` | rc=0（log 02 前基线；01 尺 1.42s） | **rc=0**（log 02/26，count3 4.38s） |
| `cmd/wisp` 定向 13 枚（233×3＋255×7＋223×2＋227×1）`-count=1` | 恰红 1 枚＝`TestTicket255RosterStillMatchesTheActualReadSites`（**HEAD 既存红**，A/B §6 末行） | **rc=0**（log 04/25；count3 64.0s 全绿） |
| `cmd/wisp` `-count=3` 同尺 | base 同尺恰红 3 发×同一枚既存红（log 16） | **rc=0**（log 25） |
| `internal/risk`（零触碰对照） | 唯一红＝`TestC26RewriteAccount…`（log 28：**纯 HEAD 同红**＝cfg125-r1 §6 在册 A380 存量） | **同枚同红**（log 27）；`git diff -- internal/risk` 空 |

**A/B 两笔（判"不是我"的尺）**：①`TestTicket255RosterStillMatchesTheActualReadSites` 在 HEAD 基线（log 01/16）红＝wire-r1 的 `resident_voice_windows.go` 新读 `[session]:1329`/`[voice]:1315` 没回 255 名册——**本腿已按扫描方向修好**（§7），故改后此枚绿。②`TestTicket223HandEditedFsLooseningCostsAnL2Card` 60 发 A/B（log 18/19/20/21）：**base 30 发 3 红 / mine 60 发 2 红**（同红行 `:436/:439`「console 没渲染卡」）＝`awaitCard`（认账）与 `Prompt` 的 `Fprintf`（渲染）之间的**仪器窗口竞态，HEAD 既存**（同文件 mode 案例 `:577-586` 注释自陈此形状、票 351 只修了 mode 案例的读法）——本腿**不顺手修**（禁令：别的票的用例），台账归验收腿转编排者。

## 7. 读数⑧：四数＋两尺＋d22scan

- **写面四数**：`git diff --name-only` ＝ 6 改＋1 新（§0 清单）；`internal/risk|speech|session|audio` 与 `testdata|frontend|design|scripts|.github` **零字节**；`git show --stat`（提交后）只含本腿 7 件。
- **甲 `gofmt -l`／乙 `gofumpt -l`**（绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`，全树，归属判 `[\\/]`）：两尺各命中＝`.scratch\probes\**` 探针件＋**3 枚既存假红**（`cmd\wisp\models.go`、`internal\agent\approval\pending_read.go`、`internal\agent\tools.go`，CRLF 造，与 wire-r1/15e-r1/cfg125-r1 证据件逐字同）；**本腿 7 枚件两尺皆不在名单**（log 13/14）。
- **d22scan**（`tools/d22scan` 独立模块 build，扫 `export-mut` 导出树）：**rc=0 clean**——bans#1-5 internal/=264 cmd/=45、#6 frontend=116、#7 tools=25、#8 design=14 frontend=116 internal=684 cmd=160（log 22/24）。**正控**＝种一枚真 `≤`（U+2264，字节级 `e289a4` 现量）进 `internal/config/zz_d22_positive_control_233r1.go` 字符串 ⇒ **rc=1 且逐字点名** `[emoji] ban #8 glyph in scope internal/ … non-comment text, string literals included`（log 23）⇒ 尺是活的；种后即删、复扫 clean。
- **`go vet ./cmd/wisp ./internal/config/`** rc=0。

## 8. 255 名册重裁（扫描尺本身要求的，方向＝只增读者）

`TestTicket255RosterStillMatchesTheActualReadSites` 在 HEAD 上就红（§6 A/B）⇒ 重裁两行，全部**加读者、不加判定**：`session` 行 `no-reader(扫描零命中)`→`other-process`（新 cite `resident_voice_windows.go:1329 [c.Session.WarmTimeoutSec]`，token 唯一命中现量 1）；`sectionReadSites["Session"]` `nil`→`{cmd/wisp/resident_voice_windows.go}`、`["Voice"]` 加第四件；[voice] 四个 hot 键行的读者清单文字同步（键级 verdict 仍 no-reader）。`agent` 行 `snapshot-only`→`consumed`（新 cite `run.go:1263 [cfg = rt.mgr.Config()]`，唯一命中）；两枚 255 判定用例同步（`claimable == [llm agent]`、`[session]` 阳性对照样本换到 `[ball] size = 66`——"无读者 anywhere"样本随事实搬家，test 注释逐条写明搬家理由与三段历史）。**无一处放宽断言**：`TestTicket255ReceiptStillNamesTheLiveReadSection`/`SentenceAssemblyIsFiltered` 反向半句（claimed 不得进 honest bucket）同步收紧到 `[agent]`。

## 9. 具名残余（⛔ 藏）

1. `TestTicket223HandEditedFsLooseningCostsAnL2Card` 的仪器窗口竞态（§6②）：**HEAD 既存**，非本腿引入；修法＝把 `:434` 单读改成同文件 `:618` 的 `awaitStdout` 形（351 预案），归下一张票／票 351 射程，本腿未动。
2. reload 档"接"的深度：本宿主无语音链可重载，`reportReloadTier` 只做**告知**（D36 该档"需重载子系统"的真正模型换装在常驻进程）——票面 AC#3 的"接"形已满足（真生产赋值点），"语音链真换装"归常驻腿。
3. 票面/台账/翻勾＝编排者权，本腿零触碰（含票 223 AC#2 那格）。
4. macOS／Linux 门禁未跑（本腿改动为纯 Go 语义层，无平台分支；CI 会跑）。
5. `export-mut`/`export-base` 树里 `.scratch/**` 解包告警（tar 对个别坏名文件）＝既有，与 Go 树无关。
