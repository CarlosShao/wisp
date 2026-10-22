# 223-r3 实施腿（写腿）：票 223 剩的两格 `AC#2`／`AC#7` 的逐档读数

- 起手锚 `e94a5b17`（`git status --porcelain -- internal/config cmd/wisp docs/evidence/s1` 起手＝**空**）；本腿三枚提交：①`394cc212`（新用例）②`d989a217`（本件）③末笔补号提交（只把本件占位符换成②的 hash，不回改任何读数，见 §9）。
- 身份＝写腿：⛔ 翻框、⛔ 改 `Status:`／票名／`-done`、⛔ 台账、⛔ push；写面只碰 `internal/config/**` 与本件。
- **一句话结论**：`AC#7` 的出口**今天在生产里已经存在**（本腿交读数、零产码改动）；`AC#2` 的 hot／restart 两档本腿在 config 层补齐**逐档用例＋逐键真值**；`AC#2` 的 **reload 档生产消费者仍＝0**（唯一非测试赋值在 `cmd/balldebug`），是**票 233** 的格，本腿具名交回、⛔ 冒充已闭合。

## §1 三档逐档读数（先量后写；🌟「今天是什么读数」）

| 档 | config 层今天能读到的 | 生产端今天能读到的 | 尺／落点 |
|---|---|---|---|
| **立即（hot）** | 一次 `CheckAndReload`（1s tick 的同一个调用）后 `Manager.Config()` 答新值：`llm.roles.chat.model` `m1→m2` | **有活读面**：`cmd/wisp/panel_config_store.go:92 cfg := s.mgr.Config()`、`:98 ChatModel: cfg.LLM.Roles.Chat.Model`（每次调用重读）；票 255 的 roster 里 `llm` 是唯一 `consumed:` 行 | 新用例 `…HotTierReachesALiveReaderWithoutRestart` **PASS 0.06s**；⚠ `[ball]` 是 `no-reader:` 行（`config_readers_255.go`），**不拿它冒充「行为变了」** |
| **重载（reload）** | 值立即应用 ＋ `OnReload` 恰好一次（`[voice]`） | ⛔ **没有消费者**：`grep -rn 'OnReload\s*=' --include=*.go cmd internal \| grep -v _test.go` ＝**1 行**＝`cmd/balldebug/main.go:244`（旁支程序）；`cmd/wisp`＝**0 行**；`config_reload.go:183-184` 只有审计行 `state=applied … reload=%v`，**无操作员句** | 新用例 `…ReloadTierAppliesAndSignalsTheOwningSubsystem` **PASS 0.24s**（只钉 config 层那一半）⇒ 缺的一半＝**票 233** |
| **重启（restart）** | 旧值**不动** ＋ `Report.Restart=[app]` ＋ `OnRestartPending` **恰好一次**；why＝`TierOf("app.language")=="restart"`（同段 `app.theme=="hot"` 作对照） | **有明示出口**：`cmd/wisp/config_reload.go:119` 赋值、`:325-330` 打到 `rt.stdout`（§2 逐字） | 新用例 `…RestartTierTellsTheHostWhyAndKeepsTheOldValue` **PASS 0.07s**；生产句由 `TestTicket223RestartTierSaysItWillNotApply` **PASS 5.75s** 钉住 |

- 改前基线（`e94a5b17`，尺＝`go test -count=1 -v ./internal/config/`）：**RUN=196／PASS=115／FAIL=0／SKIP=0**，rc=0，`ok 6.005s`。
- 三档**在同一段上的分离**见 §3；三档都不是「静默不生效」：hot 有活读面、restart 有明示句、reload 的**事件**有出口但**消费者没有**（=233，⛔ 用静默凑）。

## §2 `AC#7` 的出口：谁赋值、用户在哪看到、那句「为什么」逐字

- **谁赋值（生产，非测试）**：`cmd/wisp/config_reload.go:119` `rt.mgr.OnRestartPending = rt.reportRestartPending`；所在函数 `startConfigReload`（`:109`），装配根调用点 `cmd/wisp/run.go:861` `rt.startConfigReload()`。
  尺＝`grep -rn 'OnRestartPending\s*=' --include="*.go" cmd internal | grep -v "_test.go"` ⇒ **恰好 1 行**（改前该尺 0 行，见票 223 现量表第 ④ 格）。
- **用户在哪看到**：`rt.stdout` 那行 `fmt.Fprintf`，`cmd/wisp/config_reload.go:325-330`；同发另写两行审计（`:319-324`，`tier=restart`／`keys=%q`）。
- **「为什么」逐字（产品模板：`原因` 在 `cmd/wisp/config_reload.go:327-328`，整句 Fprintf 在 `:325-330`）**：
  > `原因：这几枚键在进程启动时被读一次就交给平台层（开机自启注册、单实例锁、界面语言），本次运行没有重新注册它们的路径。涉及：%s。`
  > 同句还带「未生效」与「两件事都没发生」：`这些段的改动本次运行不会生效（D36 重启档，需要重启进程）：%v`＋`注意两件事都没发生：文件里的新值没有被丢掉，内存里的旧值也没有被偷偷替换。`
  `涉及` 的清单＝`restartTierKeys`（`:342-344`＝`app.language / app.autostart / app.single_instance`），由 `TestTicket255RestartTierKeysAreBackedByATest` **PASS 0.14s**（3 子测试）逐个过真 Manager 落 `Report.Restart`。
- **为什么这句不是空壳（本腿补的那半边）**：句子的真值是「这几枚键属重启后生效档」。本腿在 config 层把它钉成双保险：`TierOf(key)=="restart"`（§1 第三行、§5 的 M4）＋逐键落桶（`…EveryAppKeyLandsInTheBucketItsTierNames`，4 子测试全 PASS）。
- ⚠ 归责边界：那两句**操作员面**的钉子（`TestTicket223RestartTierSaysItWillNotApply` 的窗口化内容 needle）是**票 232** 落的，属 `cmd/wisp`，本腿只跑读数、未改一字。

## §3 两格 ⛔ 合并的证据（两个不同可观察物，同一发里分开读）

新用例 `TestTicket223R3OneEditGivesTwoDifferentObservables` **PASS 0.12s**：**同一笔改动**同时动 `app.theme`（hot）与 `app.language`（restart）⇒ 一次 reload 里：

- **可观察物 A（AC#2 第一档）**：`Report.Hot=[app]` 且 `Config().App.Theme=="light"`（活值变了）。
- **可观察物 B（AC#2 第二档＝AC#7）**：`Report.Restart=[app]` 且 `Config().App.Language=="zh-CN"`（**没变**）＋`OnRestartPending` **恰好一次**。

⇒ 两枚**不同对象**（一个值变了、一个值没变而人被告知），且**突变可分**：M5 只打红 A 那一半（`app.theme = "dark"`）、M3 只打红 B 那一半（`OnRestartPending=[]`）。另有两枚单向负向断言防「狼来了」：hot 用例断**重启通知 0 次**；restart 用例断 **`Report.Hot` 空**。

## §4 新用例逐名 ＋ 它今天响在哪一跳（「为什么今天响」）

| 用例（`internal/config/tiers_223r3_test.go`） | 观测物 | 它护的那一跳（产码） | 摘掉那一跳后红在哪 |
|---|---|---|---|
| `TestTicket223R3HotTierReachesALiveReaderWithoutRestart` | 一次轮询后同一个活读面答新值 | `manager.go:302` 热段 `plan.set = append(plan.set, s.set)` | `:116 after one poll the live reader still answers "m1", want m2` |
| `TestTicket223R3RestartTierTellsTheHostWhyAndKeepsTheOldValue` | 旧值不动＋恰好一次通知＋why 两行 | `manager.go:202` `m.OnRestartPending(rep.Restart)`；`tiers.go:48` `app.language:"restart"` | `:151 OnRestartPending = [], want exactly one [app] notice`／`:156 TierOf(app.language)="hot"/true, want restart` |
| `TestTicket223R3ReloadTierAppliesAndSignalsTheOwningSubsystem` | 值应用＋事件恰好一次 | `manager.go:199` `m.OnReload(rep.Reload)` | `:185 OnReload = [], want exactly one [voice] event` |
| `TestTicket223R3OneEditGivesTwoDifferentObservables` | 同一发里 A／B 两个可观察物 | 上述两跳（A←`manager.go:358`，B←`:202`） | M5 ⇒ `:218 app.theme = "dark", want light`；M3 ⇒ `:225 OnRestartPending = [], want one [app] notice` |
| `TestTicket223R3EveryAppKeyLandsInTheBucketItsTierNames`（子测试：`app.theme`／`app.language`／`app.autostart`／`app.single_instance`） | 逐键落桶＝`TierRegistry`（why 的真值） | `manager.go:354-356` 的三键比较 ＋ `:358` 的 theme | 键若从 `planApp` 删掉⇒落 0 桶红 `:287`；M4 ⇒ `:262 want "restart"` |

- 今天五枚**全绿**（`…/app.*` 四子测试也逐名 PASS）；它们只在被护的那一跳被摘掉时转红——§5 是实测，不是声明。

## §5 反向判据：五发突变，三读数逐发（变异红／还原绿／原样绿）

全部只在**仓外导出树** `/d/tmp/wisp/20261022/223r3/exp`（`git archive 394cc212`）里物理改字节；工作树从未携带突变。尺＝`go test -count=1 -run "TestTicket223R3" -v ./internal/config/`，日志在 `/d/tmp/wisp/20261022/223r3/logs/M*.txt`。

| 突变 | 锚（改法） | 原样绿 | 变异红（红行逐字） | 还原绿 | 变异后 sha256 |
|---|---|---|---|---|---|
| **M1** 摘 hot 段 apply | `manager.go:302` 注释 | rc=0（5 PASS） | rc=1：`…HotTier…` FAIL，`:116`（见 §4） | rc=0 | `bdcca02f017c5b34c1fc275090a0f5b20da6f5af3502f2e573df2c7dab417815`（该发 stdout 被 `tail` 截掉，hash 由同一 pristine ＋ 同一 sed 重算） |
| **M2** 摘 reload 派发 | `manager.go:199` 注释 | rc=0 | rc=1：`…ReloadTier…` FAIL，`:185` | rc=0 | `ce2c2222bef1118915f46d73041b3b5eedba20ffd38b697fd9685a0dd9101310` |
| **M3** 摘 restart 派发 | `manager.go:202` 注释 | rc=0 | rc=1：`…RestartTier…`／`…OneEdit…`／`…EveryAppKey…` FAIL；`:151`／`:225`／`:291`×3（`app.language`／`app.autostart`／`app.single_instance` 各一行 `the host was told 0 times (want 1)`） | rc=0 | `9dd9bfbe4713051cb5a61511b7467a784b3225dbc541b15a831212df97c91638` |
| **M4** 摘 why（把 `app.language` 改成 hot） | `tiers.go:48` `"restart"→"hot"` | rc=0 | rc=1：`…RestartTier…`／`…EveryAppKey…` FAIL；`:156`／`:262` 逐字见 §4 | rc=0 | `1b73ac4a3f8e704655debcdef3ca90b6ccb3e6c9e174920ce54c16d262ce7baf` |
| **M5** 摘 theme 的 hot apply | `manager.go:358` 注释 | rc=0 | rc=1：`…OneEdit…`／`…EveryAppKey…` FAIL；`:218`／`:278 app.theme is registered hot but the live value did not move: dark` | rc=0 | `ef8bb857d8d57170e32c83be46f2faaa74581de24fe1a4c51c977314048b6b16` |

## §6 还原核（先问「谁算的、什么算法、几位」）

- 算法：**`sha256sum`**（GNU coreutils／Git Bash）＝**SHA-256、64 位十六进制**；⛔ 与 `git hash-object`（SHA-1、40 位）混比（尺的射程清单第 9 例）。
- `internal/config/manager.go`：pristine＝M2／M3／M5 还原后＝**`4bb091d99f3c62925558be03933217c9fbb10d74696a8c57d16c251a05b2e8d0`**（三处同串）。
- `internal/config/tiers.go`：pristine＝M4 还原后＝**`1f168efe1d3e7a562f973cd48f0709d6bac18ee47b30dc9aa4bacae5260c4de0`**（同串）。
- 导出树≡提交树：上面两枚值同时等于 `git show 394cc212:<path> | sha256sum` 的输出（同算法、同位数）；导出树两枚文件 `tr -cd "\r" | wc -c`＝**0**。
- 工作树收工：`git status --porcelain -- internal/config` 只剩本腿两枚提交所覆盖的文件（下文 §9），无突变残留。

## §7 门禁 rc（串行、变异全部结束后才跑）

- `go build ./...` **rc=0**（`logs/build.txt`）。
- `go vet ./internal/config/` **rc=0**（`logs/vet.txt`）。
- `go test -count=1 -v ./internal/config/` **rc=0**，四数**两发同数**：**RUN=205／PASS=120／FAIL=0／SKIP=0**（`ok 3.892s`／`ok 62.862s`；时间只是读数、⛔ 判据）。
  - 四数的形状（先问「数的是哪个对象」）：`^=== RUN` 含**不缩进**的子测试 RUN；`^--- PASS` 只数顶层，子测试 PASS **缩进 4 格**共 **85** 行 ⇒ **205 = 120 + 85**（与总数对上）。
  - 改前基线同尺 RUN=196／PASS=115 ⇒ 差 **+9 RUN／+5 顶层 PASS**，恰等于「5 枚顶层用例＋4 枚子测试」，⛔ 无多报少报。
- `sh scripts/d22scan.sh` **rc=0**（`clean`；ban #8 扫 `internal/` 621 枚 Go、`cmd/` 148 枚；`logs/d22scan.txt`）。
- **格式甲乙两尺**（`gofumpt v0.12.0`，绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`）：
  - **甲（判据）**＝`git archive 394cc212` → 仓外 `exp-fmt/`，`gofumpt -l .` ⇒ **20 行、全部 `.scratch/`**；归属用 `[\\/]`：`grep -vcE '^\.scratch[\\/]'`＝**0**，本件命中 **0**。
  - **乙（工作树，只诊断）**＝`gofumpt -l .` ⇒ **49 行**、本件命中 **0**；3 条非 `.scratch` 行是**检出 CRLF 副本**（`cmd/wisp/models.go` CR=389、`internal/agent/approval/pending_read.go` CR=131、`internal/agent/tools.go` CR=225，`git diff --quiet HEAD` 三枚均 rc=0＝不是本腿改动）。
  - 本件自己：两尺都 **0 命中**；`tr -cd "\r" | wc -c`＝**0**（LF）。

## §8 红名册（取主检出；新增红必须 0）

- 本腿整包只跑写面所在包 `./internal/config/`：**FAIL=0 ⇒ 红名册为空**，**新增红 0**。
- `cmd/wisp` 只按**名单**跑两枚（§2 那两枚，均 PASS）；⛔ 拿 `cmd/wisp` 整包四数当判据（派单明令：未判决恒 1／红名册逐发漂移），故本腿不跑、也不报。
- ⛔ 未跑 `./internal/...` 全量：本波 `internal/agent`／`internal/winsec` 等归别的腿在飞，跑出来的红不是本腿可归因的读数（AC#5 的全量终态属非实现者腿）。

## §9 提交号逐枚 ＋ `git show --name-only` 的 `.go`／越界（正控）

| 提交 | 内容 | 该提交的 `.go` | 越界 |
|---|---|---|---|
| `394cc212` | `internal/config/tiers_223r3_test.go`（+298，新件按名字面点名入 pathspec） | **1** | **0** |
| `d989a217` | `docs/evidence/s1/223-r3-impl.md` | **0** | **0** |
| ③ 末笔补号提交 | 只把本件两处 `‹本件提交›` 占位符换成上一行 hash；不回改任何读数 | **0** | **0** |
- 尺＝`git show --name-only --format= <c> | grep -c '\.go$'`＋逐名对写面。
- **正控（同一把尺打一笔含 `.go` 的真提交）**：`17cfe8cf` ⇒ `.go`＝**2**（`internal/memory/dao_session.go`／`internal/session/ticket166r1_ac2_gaps_test.go`）⇒ 这把尺**能看见 `.go`**，本表两行 0／1 不是恒 0 的装饰。**负控**：`b4102f1e`（docs-only）⇒ `.go`＝**0**。

## §10 具名交回（不在本腿写面／判不了）

1. **`AC#2` 的 reload 档生产消费者＝0**（§1 第三行）。归口＝**票 233**（今天仍「待派」，文件名无 `-done`）。本腿顺手量到 233 的两条现量**一真一半腐**：
   - 真：`Tier` 三常量在非测试代码里**赋值点 0**（非测试命中 6 行全是 `schema.go:34-43` 的注释＋声明，`grep` 其余零命中）。
   - 半腐：233 说「reload 档**既没有生产者键**」——自票 255 起 `TierRegistry` 已有 **20 枚 `"reload"` 行**（`voice.*`）；缺的只剩**消费者**那一半。
2. **`AC#2` hot 档的「行为变了」在生产端只有 `[llm]` 这一枚活读面**（票 255 的 roster 既判）；`[ball]`／`app.theme`／`voice.tts.speed` 等仍是 `no-reader`／`snapshot-only`。本腿不拿它们冒充，也不改 255 的判语。
3. 本腿**零产码改动、零新增导出名**：`internal/config` 的非测试 `.go` 与起手**逐字节相同**（`git show 394cc212` 之前/之后的 blob 相同；本件与用例是仅有的两笔改动）。
