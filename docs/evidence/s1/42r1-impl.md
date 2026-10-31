# 票 42 落地腿 `42-r1` —— Watchdog 落地件（写腿）

> **尺**：实现期 `HEAD=07a010d59`（`dev`），现跑 2026-10-31 10:0x–11:0x +08。★读数逐条现跑，⛔ 转述。
> 派单文字写 `281862907`；树此后前移两格＝**编排者记帐提交**（`4cd8cb4b1` 派单存档、`07a010d59` A1296 结票/停车点），
> 现跑核实两格均未触碰本腿写面（`git status --porcelain -- internal/ docs/evidence/s1` 仅编排者两份裁决件＋本腿新文件），
> 分支全程 `dev`、非 detached ⇒ 按三核意图通过，未做任何 git 状态操作。
> **写面**：`internal/watchdog/**` 新 5 枚（`watchdog.go`／`settle.go`／`watchdog_test.go`／`settle_test.go`／`config_poll_test.go`）
> ＋ `internal/observe` 新 2 枚（`caps.go`／`caps_test.go`，冻结常量只读 accessor）＋本文件。
> **零字节**：`cmd/wisp/**`、`internal/session/**`（0 字节）、`internal/speech/**`（0 字节）、`internal/risk/**`、
> `internal/config/**`、`internal/statemachine/**`、`testdata/**`、`frontend/**`、`design/**`、`scripts/**`、`.github/**`、
> 票面、台账、`docs/specs/**`、`docs/PLAN.md`。
> **临时件（只建不删）**：`D:\tmp\wisp\20261028\42r1\{logs,mut}`。★突变走仓外件＋`-overlay` ⇒ 工作树零突变窗口。

## 0. 一句话

`internal/watchdog` 从「只有 doc.go 的承诺」变成真的循环：**一条** D38b 常驻名 `watchdog` goroutine、1s tick 双职责
（资源执法 ＋ SPEC-03 §4.3 的 config mtime+size 轮询），按态查 observe 冻结行（25/90/350/350/600MB、GDI<200、
handles<600、goroutines 6/7），超限持续 60 tick → 卸载趟＝**冻结白名单 ∩ 注册表**（结构上够不到白名单外任何组件），
Armed 态 KWS 结构性不可卸＋留态告警（D43#9），泄漏指示器＝告警＋诊断快照（D42#10，绝不顺卸），settle 重测 10s、
**连续 3 败 → Fatal 告警**（D43#35）。生产调用者 0（装配归宿主腿，同 28 形状）。

## 1. ★★★ 票面 5 格逐条（原文 ＋ 三分类 ＋ 读数）

| # | 票面原文（逐字） | 判定 | 读数 |
|---|---|---|---|
| 1 | `Forced-leak fixtures per state: Sleeping over-cap → unloadables dropped (panel hidden) → within cap or alert; Armed over-cap → KWS RETAINED + warning (rule-critical test).` | **①可做**（fixture 层全做；「真面板隐藏」触 `cmd/wisp`＝**具名交回**，本腿交付其 watchdog 侧语义＝`panel-webview` 白名单组件的 dropped 读数） | Sleeping：fixture 泄漏 30MB（cap 25MB）持续 3 tick → 卸载趟真发生：`panel-webview`＋`bm25-index-cache` 各 **dropped=1**（`outcome":"dropped"` 落日志），树回落 20MB → **within cap、0 告警**；不回落臂：组件全 skip → 树仍 30MB → **Overrun 告警恰 1 次/episode**（8 tick 4 趟不重复）。Armed：100MB（cap 90MB）两个 sustain 窗 → **`kws` unload calls=0（RETAINED）**、bm25 dropped×2、Armed 告警恰 1 次/episode，Detail 具名「KWS retained / suggest disabling KWS」；宿主映射 `EvWatchdogOverrun` → 机器**留 Armed**＋副作用 `kws.keep-alive-alert`（D43#9 表行现跑） |
| 2 | `3× settle-failure fixture → WatchdogAlert state + restart affordance works.` | **①可做**（★D43 表**有** WatchdogAlert 合法行：`table.go:256`＝D43#35 `From: AnyState, Event: EvWatchdogFatal, To: StateWatchdogAlert, SideEffects:[ui.one-click-restart]` ⇒ 停手条件①排除； Armed 态另有 #9 留态行） | Warm cap 350MB、fixture 400MB 整窗不回落：`VerifySettle` 连跑 3 次 → 3 败（`BackWithinCapMS=-1`），第 3 次起 **AlertFatal 恰 1 次**（class=resource、state=Warm、计数重置 0）；测试接线（宿主缝）把 Fatal 映射 `EvWatchdogFatal` → 机器 `Warm→WatchdogAlert`，sink 收到 **`Effect{D43:35, Name:"ui.one-click-restart"}`**；连续性：败/败/**成**/败/败/败 → 第 6 败才第 2 次 Fatal（成即重置）；`released=false` 按契约必败（`TestVerifySettleReleasedFalseFailsByContract`）；未知目标态 → config 类错误拒绝 |
| 3 | `Leak indicators: fixture handle/GDI growth crosses caps → alert + snapshot.` | **①可做**（caps＝observe 冻结现值 600/200；goroutines 冻结 cap 仅 Sleeping 6/Armed 7；**USER 无冻结上限 ⇒ record-only**，⛔ 发明 cap——派单明令 caps 对照 observe 现值） | fixture handles=601、GDI=201、goroutines=7（Sleeping cap 6）→ 单 tick 起 **3 条 Leak 告警**，各带非空 JSON 快照（metric/state/**穿越样本本体 tree**/roster/runtime 数），`measured=601/201/7` 逐位对；同 episode 去重（再 2 tick 仍 3 条）；handles 回落后再越 → **重臂**（handles 共 2 条）；**泄漏告警 0 次卸载**（注册组件 calls=0），日志行自述 `no auto-unload` |
| 4 | `All actions logged with metrics; no unload of non-whitelisted components ever (log audit).` | **①可做** | 结构保证：卸载趟**只走冻结白名单**、按名查注册表 ⇒ 注册名 `kws`/`audio-capture`/`sqlite-db-writer`/`session-engines`/`agent-context` **不可达**。审计测试：白名单 4 项＋NEVER 5 项全注册，Sleeping/Warm/Armed/Conversation 四态各跑一趟 sustain，解析全部 JSONL：每条 `watchdog: unload action` 的 component ∈ 白名单、outcome ∈ {dropped,skipped}；**NEVER 名零出现、Unload 调用全 0**；每条 `unload pass complete` 带 `before_bytes/after_bytes/cap_bytes/dropped`；组件 panic 被逐调用 recover（`outcome=error` 落日志、其余组件继续、循环存活两趟） |
| 5 | `Config mtime poll fires reload events identical to 05's contract (shared test).` | **①可做**（★05 contract 找到＝SPEC-03 §4.3＋`config.Manager.CheckAndReload`（幂等 mtime+size、D36 三档、OnReload/OnRestartPending）⇒ 停手条件③排除） | **shared test 真共用**：同一 scenario 体、双驱动——`05`＝裸 `CheckAndReload` 直调（05 自测形状）、`42`＝同一 Manager 挂上 watchdog tick；链式三改（[voice] 换模型→reload、[ball] size→hot、[app] language→restart）两驱动事件序列 **DeepEqual 全同**：`Reloads=[[voice]]`、`Restarts=[[app]]`、ball 56→60 即时生效、voice 模型 paraformer-zh 落地、app.language 留 zh-CN（restart 档不生效）；42 驱动另证**一条循环双职责**：`CountByName("watchdog")==1` 且同循环日志同时含 `config poll applied changes` 与 `unload action` |

## 2. caps 对照（停手条件④核）

票面 Key constraints 引 D42#10 的 `handles>2000 / GDI>5000` 为**预测量数字**；observe `thresholds.go` 冻结值
（SLO.md ruling 1：窗口栈实测 ~420 ⇒ <600；GDI<200）是权威现值，票面自己说阈值取「the SLO table (02-measured values)」。
⇒ **无冲突**：watchdog 经新 accessor 读**同一组冻结常量**（`MemCapBytes/CPULimitPercent/GDILimit/HandleLimit/GoroutineLimit`），
零第二张表（188 定式：引用可以、复制不行）；accessor 钉值测试锁数（`TestFrozenCapAccessors`）。60s sustain／10s settle
两常数落 watchdog 侧（执法节奏，非 SLO 阈值），默认 60 tick×1s／10s，测试缝可缩。

## 3. ★★★ 反向判据：摘 watchdog 判定 ⇒ over-cap 必红（三发突变＋还原核）

突变均在仓外 `mut\` ＋ `go test -overlay`，工作树零突变窗口。

| 发 | 变异 | rc | 红的断言（逐字选） |
|---|---|---|---|
| **A** | sustain 判定摘除（`if l.overTicks < l.sustainTicks` → `if true`） | **1** | **7 红**：`panel-webview unload calls = 0, want 1 (one pass over the sustained breach)`；`no unload action lines found; the audit has nothing to read`；`[42] the resource duty never ran on the loop`；Armed/泄漏臂全红 |
| **B** | 白名单拆掉（卸载集改走注册表 `range l.Registered()`） | **1** | **3 红（外科）**：`kws unload calls = 1, want 0: a non-whitelisted component was unloaded`；`KWS unload calls = 2, want 0: unloading KWS in Armed is the D32 revision 3 harmful behavior`；audit `unload action for non-whitelisted component "kws"/"audio-capture"/"agent-context" (outcome dropped)` |
| **C** | Fatal 判定摘除（`fatal := !rep.Pass && n >= 3` → `fatal := false`） | **1** | **恰 1 红**：`TestThreeConsecutiveSettleFailuresReachWatchdogAlertWithRestart`（fatal alerts 断言） |
| **还原** | 无 overlay 重跑 ×3 | **0** | `ok internal/watchdog`；observe/config/statemachine/session 同轮 rc=0 |

**还原核**：①树内四枚变异串 grep（`if true {`／`range l.Registered()`／`fatal := false`／`MUTANT`）＝**0 命中**；
②突变期哈希 T1（watchdog `EAE498A6…`、settle `0D589C98…`）之后树文件仅有的差异＝本腿**记账内**的 Class 字段补充
（还原后 `gofumpt` 重排对齐一次，突变残留 grep 复跑 0）；③`git status`：本腿 6+1 枚全部 `??` 新文件、
本腿写面零 `M`。⇒ 变异只存在于 `D:\tmp\…\mut\`，`-overlay` 只换编译输入。

## 4. ★ 四数 ＋ 两尺 ＋ 门禁

| 面 | 命令 | rc | 读数 |
|---|---|---|---|
| `internal/watchdog` | `go test -count=1 -v ./internal/watchdog`（另 -count=3／-count=5 各一轮） | **0** | **12 PASS／0 FAIL／0 SKIP**，`ok … 0.36s`（复跑 1.0–1.7s） |
| `internal/observe` 改后 | `go test -count=1 -v ./internal/observe` | **0** | **82 PASS／0 FAIL／0 SKIP**（HEAD 81 枚 `git grep -c '^func Test'` 现点＝81，＋本腿 1 枚钉值案） |
| 邻接包 | `go test -count=1 ./internal/{config,statemachine,session,observe,watchdog}` | **0** | config 1.8s／statemachine 0.24s／session 2.9s 全 ok |
| `go vet` | observe＋watchdog | **0** | 空 |
| 甲 `gofmt -l` | internal cmd tools scripts frontend | 0 | 4 行，**全非本腿**：`internal\agent\approval\pending_read.go`、`internal\agent\tools.go`（既知 2 条 CRLF 假红）、`internal\memory\undo.go`（HEAD 既有格式漂移）、`cmd\wisp\resident_voice_windows.go`（**41a2-r1 在飞**）。28 时代的第 3 条 `cmd\wisp\models.go` 已被他人修好不再红 |
| 乙 `gofumpt -l`（绝对路径） | 同上 | 0 | 同上 4 行逐字相同；**单点跑本腿 7 枚＝空** |
| `d22scan` | `cd tools/d22scan && go run . -root ../../` | **0** | `clean - no D22 ban violations`；正控计数：bans#1-5 internal/=269、cmd/=45、ban#6 frontend/=116、ban#7 internal/tools/=25、ban#8 design/=39、frontend/=116、internal/=695、cmd/=161 |
| 零字节核 | `git status --porcelain -- internal/session internal/speech internal/risk internal/config internal/statemachine cmd/wisp` | —— | session/speech/risk/config/statemachine **0 行**；cmd/wisp 仅 **41a2-r1** 的 3 处（`resident_voice_windows.go`、`resident_audio_windows.go`、`resident_voice_41a2r1_windows_test.go`），本腿对它 **0 字节** |

## 5. 全量 sweep 的红处归属（诚实登记，⛔ 冒领）

`go test ./internal/...`（CGO0）除本腿面外另有红，逐一单独复跑＋查脏归属：

- `internal/panel`：`TestApprovalCardViewJSONKeysMatchFrontendTypes`＋`TestComposerContractTypesMatchFrontend`（Go `Snapshot` 发出 `[instructions tasks]` 而 `PanelSnapshot` 接口未声明＝**panel↔frontend 契约漂移**，frontend 归机主外部 agent）＋`TestStreamLog*35r8` 两案；**包内 git status 干净** ⇒ HEAD 既有。
- `internal/risk`：`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`＝扫描发现 `internal/panel/composer.go:127`、`display_scrub_216.go:26/37/218` 读 `Result.Canonical` 不读 rewrite account（「ticket 102 的洞被新消费者重开」）——**panel 的债经 risk 的门禁暴露**；risk 包干净 ⇒ HEAD 既有。
- `internal/winsec`：`TestAC1…342r1` racer 概率红（400 发中 1–2 发 witness 被换写）；包干净 ⇒ HEAD 既有（票 342 域）。
- `internal/agent/approval`：仅并发 sweep 红、单独跑 **ok** ⇒ 负载抖动。
- `internal/speech [setup failed]`＝CGO0 下 cgo 包本就不编译（28 §6 同口径），非回归。

★本腿对上述各包**零 import 边**（watchdog 为新包无人引；observe 增量为纯函数、无 init/状态），可证惰性。

## 6. 未做／边界／具名残余

- **生产调用者 0**：`cmd/wisp` 未构造 Loop。装配（真 TreeReader＝proc 侧、StateFn 接状态机、OnAlert→`EvWatchdogOverrun`/`EvWatchdogFatal` 的 D43 映射、组件注册）归宿主接线腿；本腿交付库＋读数（同 28「装配等下一波」形状）。
- **AC#1 的「panel hidden」cmd/wisp 半边具名交回**：真面板的隐藏/销毁动作在面板宿主；本腿交付白名单机制＋`panel-webview` 组件语义（用户在看 ⇒ `dropped=false` skip 分支有测）。
- **R-42r1-1（新发现真缺口，归 `internal/config` 所有者，⛔ 本腿不碰）**：`CheckAndReload` 的失败路径**一律收养 stat**——本腿并发轮询实证：与原子 rename 竞态的一次瞬时 sharing violation 读失败会把该次写**永久隐藏**（stat 已收养、内容从未生效）。05 契约「失败收养防每 tick 重报错」对解析类失败是对的，对**瞬时 IO 类失败**是洞。建议小票：仅解析/校验类失败收养 stat，IO 类失败不收养。本腿测试侧已按「真实写入者会再写」建模（42 驱动未生效即重种），05 契约一字未动。
- **USER 对象**：冻结表 record-only（无绝对上限）⇒ 记入 tick 日志、从不 gate（⛔ 发明阈值）。
- **CPU**：每 tick 记录（含 per-state 冻结限值）；卸载趟由**内存 sustain 规则**驱动（D32③ 的执法语义是内存回落），CPU 无对应卸载语义。
- **`internal/watchdog/doc.go` 的 `DEFERRED(watchdog loop/thresholds)` 标记未摘**（1:1 双向登记规则；摘标记属 `225-deferred-registry-audit` 写面，同 28 先例）。
- **配置键 0 新增**（节奏全部为 Options/测试缝，非 `[observe]` 键）；`thresholds.go`/golden **零字节**。
- 包内**零裸 `go` 语句**（唯一 goroutine 经 `observe.Registry.Spawn("watchdog",…)`，名字在 D38b 常驻名册）；零墙钟算术（sustain=连续 tick 计数、settle 走 observe 单调 CheckSettle、wall 仅落盘戳）——均为 d22scan rc=0＋`TestNoBareGoFuncInProductionCode`（observe 既有案扫全仓）现跑背书。
