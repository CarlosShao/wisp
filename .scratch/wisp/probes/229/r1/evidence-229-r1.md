# 229-r1 证据件：子代理收尾协程名不在 D38 名册（AC#1–#6）

- 票：`.scratch/wisp/issues/229-subagent-spawns-an-unregistered-goroutine-name-and-drowns-the-leak-alarm.md`（全文已读，判据 AC#1–#6 逐格）
- 腿：`229-r1`；分支 dev；起手 HEAD `5f513c612006fdedc8a116781fce3efc8c420bfc`（`logs/00-head.txt` 逐字；`logs/` 在仓根，不在本目录）
- 改前⇄改后（甲形）：`internal/tools/subagent_197.go` 一行 spawn 名字符串 ＋ 一段注释（最小 diff，root 标签不动）
- 新增测试件两枚：`internal/tools/spawn_roster_census_229_test.go`（普查尺＋断）、`internal/tools/spawn_finish_watcher_229_test.go`（真路读数＋正控）
- 原始读数：本目录 `logs/`（`ac1-roster.txt` · `ac4-census-raw.txt` · `pre-229-tests.txt` · `post-229-tests.txt` · `pre-gate-two-pkgs-overlay.txt` · `post-gate-two-pkgs.txt` · `gates.txt` · `ac6-full.txt` · `ac6-attribution.txt`）

## AC#1 名册先现读（逐条现跑，不引票面表）

出处 `internal/observe/goroutine.go`（两把尺：`logs/ac1-roster.txt`）：

| 表 | 成员（逐枚） | 声明处 |
|---|---|---|
| 常驻（D38b，基线 6） | `ui-sta` · `audio-capture` · `hotkey-listener` · `db-writer` · `watchdog` · `log-flusher` | `:40`（ResidentBaseline=6）`:43`（ResidentNames） |
| 按需（Armed 才起） | `kws-infer` | `:48` |
| 每任务前缀（每任务 3） | `agent-task-` · `tool-exec-` · `approval-waiter` | `:52`（`PerTaskPrefixes`；`PLAN.md:2834` 同书面） |
| 临时（DisposalScope 内） | `asr-infer` · `tts-infer` · `panel-host` · `retention-job` · `model-downloader` · `memory-extract` · `disposal-worker` | `:57` |

- 判定函数 `ClassifyGoroutine` 在 `:64-84`：resident/on-demand 精确相等、per-task 前缀匹配（`strings.HasPrefix`）、temporary 精确相等，其余 `CategoryUnknown`（注释逐字 `NOT in the roster: leak symptom`）。
- 不匹配的后果在 `:271-272`：`slog.Warn("goroutine outside the D38 roster (leak symptom)", "goroutine", name, "owner", owner)`。
- **`grep -n "subagent" internal/observe/goroutine.go` = 0 行、rc=1**（2026-10-13 20:5x 现跑，`logs/ac1-roster.txt` 末节）。

## AC#2 修法：甲形（换成名册内形状），乙形具名判死

- **甲（本票落地）**：`internal/tools/subagent_197.go:394` 的 spawn 名由 `"subagent-finish-"+bg.ID` 改成 `"agent-task-"+bg.ID`。
  - 选族理由：①票面倾向"每任务前缀那一族"——收尾协程本来就派生自任务，且它在 observe.Default 里就是**孩子任务 bg.ID 的那一行**（孩子的自环路跑在私有 registry：`opt.Registry = observe.NewRegistry()`，`:297`），它的寿命＝孩子的寿命；②`agent-task-<id>` 正是"任务标识行"这一角色（D38b/`PLAN.md:2834`），拿它当孩子的过程行是同一枚任务 id，不是借别人的名字过哨；③不撞名：observe.Default 里不存在第二枚 `agent-task-<bg.ID>`（孩子自身那枚在私有 registry，且任务 id 唯一）。
  - root 标签 `"subagent-finish-"+bg.ID`（`:393`）不动：`NewRoot/NewRootFrom` 的首参是**根标签**（panic 记录里的 root 字段），不是协程名，`ClassifyGoroutine` 从不看它（14 枚根标签逐枚列在 AC#4）。
- **乙（往名册加 `subagent-finish-`）判死**：那是动 D38b 表＝契约变更，按 `SPEC-12 §4.1` 须人工批准（先落 `A##` 并摆给 owner）。本腿**没有**加任何表成员/改三张表语义（`git diff --stat` 只碰 `subagent_197.go`）。
- **第三支（把 WARN 调轻/删掉）**：未做，且 AC#3 正控把这条门钉住（改后仍响）。

## AC#3 正控：门没被拆（改前/改后都绿）

`Test229UnknownNameStillRingsLeakAlarm`（`spawn_finish_watcher_229_test.go`）：种一枚假名字 `229-probe-not-in-any-roster-table`，夹住 slog 默认口——

- `ClassifyGoroutine(假名) == CategoryUnknown`；
- `Registry.RosterReport().Unknown == [假名]`；
- 日志里 `goroutine outside the D38 roster (leak symptom)` ＋ `goroutine=假名` 逐字命中。

读数：改前一发 PASS、改后一发 PASS（`logs/pre-229-tests.txt` / `post-229-tests.txt`）。⇒ 修法只改"哪些名字被登记"，没动"没登记的名字会怎样"。

## AC#4 全仓普查（本票主要价值）：不在册还有几枚、在哪

尺＝`Test229SpawnNamesAreRosterShaped`（`go/parser` 扫生产 .go，逐站点输出；末行汇总）＋`logs/ac4-census-raw.txt`（grep 原件）。
**汇总读数（改后）**：生产 .go 289 枚（解析跳过 0）、Spawn 站点 19 个、根标签 14 个；名册外有登记 1 枚（字面量口径）、未复审 0 枚。
**名册外名册（5 个名 / 5 个站点）**：本票那枚（事故）＋4 枚**有意且码里逐字写了理由**的（不是本票地界，只登记）：

| # | file:line（改后） | 首参 / 运行期名 | 判定 |
|---|---|---|---|
| 1 | `internal/tools/subagent_197.go:394` | `"agent-task-"+bg.ID` | per_task —— **本票修掉的那枚**（改前 `"subagent-finish-"+bg.ID`＝名册外） |
| 2 | `cmd/wisp/panel_resident_windows.go:154` | `panelSTAName` = `panel-sta` | **名册外·有意**：`:80-84` 逐字"deliberately NOT one of observe.ResidentNames…加名是 D38b 契约改，要问不要自己改"⇒ 归 cmd/wisp/panel 的口 |
| 3 | `cmd/balldebug/main.go:251` | `goroutineHotkeyBridge` = `balldebug-hotkey-bridge` | **名册外·有意**：`:35-43` 逐字"R16#5 禁止借名册名蒙过漏哨…噪声留着"⇒ 归票 66 |
| 4 | `cmd/balldebug/main.go:432` | `goroutineLevelFeeder` = `balldebug-level-feeder` | 同上（票 66） |
| 5 | `internal/llm/adaptertest/mockllm.go:92` | `"mockllm-stdout-reader"` | **名册外·有意**：`cmd/balldebug/main.go:42` 登记为票 67 的 loose end（测试私有的 `mockllmRegistry`） |

**其余 14 个生产 Spawn 站点（均在册）**：`cmd/wisp/approval_always.go:99`/`approval_reply.go:456`（approval-waiter，per_task）· `cmd/wisp/config_reload.go:123`（watchdog，resident）· `cmd/wisp/resident_ball_windows.go:361`（hotkey-listener，resident）· `cmd/wisp/resident_task_source_windows.go:347`/`:454`（approval-waiter；agent-task-resident-<n>，均 per_task）· `internal/agent/loop.go:331`/`:680`（agent-task-<id>；tool-exec-<id>，per_task）· `internal/audio/audio.go:181`（audio-capture）· `internal/ball/ball_windows.go:170`（ui-sta）· `internal/memory/writer.go:79`（db-writer）· `internal/observe/logging.go:96`（log-flusher）· `internal/plugin/disposal.go:233`（参数名；唯一生产调用点 `internal/memory/retention.go:103` 传 `retention-job`，temporary）/`:296`（disposal-worker，temporary）。
**14 个根标签（不参与名册判定，逐枚已核）**：`approval-always` · `approval-reply` · `config-reload` · `panel-host` · `ball-hotkey-bridge` · `resident-task-source` · `resident-task-<n>` · `adaptertest:mockllm` · `balldebug` 两常量 · 以及 6 处由运行期对象拼的（`s.ctx`×2 / `ra.root` / `src.root.Ctx` / `ctx`×2 / `context.WithoutCancel(ctx)` / `s.name+"/go"`）。
**射程（具名）**：生产 .go（跳过 `_test.go`、`testdata/`、隐藏目录）；根＝`cmd/ internal/ tools/ frontend/`。**射程外但已点名**：8 个 `_test.go` 站点（`internal/observe/goroutine_test.go:45/51/52/123/127/190/223`，其中 `:190 mystery-worker` 是**刻意的**CategoryUnknown 正控；`internal/tools/ticket176r1_start_port_test.go:188`）、`scripts/spike/**`（独立 go.mod、零 observe 调用）、`third_party/ design/ .scratch/`。
- 顺带复核票 228 那句"口径只到 `cmd/wisp` 这一层"：**这一发扫到了全仓生产 .go**（289 枚），结论如上。
- **噪声读数**：改前整包两发 `internal/tools` 里 `goroutine outside the D38 roster` 命中 **36 行**（`logs/pre-gate-two-pkgs-overlay.txt`，现数）；改后 **0 行**（`logs/post-gate-two-pkgs.txt`）。

## AC#5 SLO 口径一句话（具名读数，安静期）

口径：SLO 那扇门数的是 registry 存活总数（`internal/observe/sampler.go:365` `Goroutines: s.reg.Count()`；生产 `cmd/wisp/slo_windows.go:283` 把 `rt.Registry`＝`observe.Default` 递给 `NewSampler`）。
读数（新起 test 二进制、孩子被按住、只有这一发在跑；`logs/pre-229-tests.txt` / `post-229-tests.txt`）：

| 时刻 | `observe.Default.Count()`（＝SLO 总数口径） | RosterReport | 旧名 CountByName |
|---|---|---|---|
| 改前 | **1** | PerTask=0 **Unknown=[subagent-finish-<id>]** | 1 |
| 改后 | **1** | PerTask=1 **Unknown=[]** | 0 |

⇒ 一句话：**总数读数不变（同一枚协程），变的是它的分类**——那枚"名册外"从 Unknown 搬进 per_task，"总数对 D38 预算读成超了"的那枚假来源（拿 Unknown 当超预算）消失。`thresholds.go`／golden 一字节未动（`git diff --stat` 可证）。

## AC#6 终态整包：`PATH=… go test ./cmd/wisp ./internal/... -count=1` → **rc=1**（`logs/ac6-full.txt`；红名逐枚归因见 `logs/ac6-attribution.txt`）

| 包 | 红名 | 归因（现跑/现读，⛔ 无一是本腿） |
|---|---|---|
| `cmd/wisp` | `TestTicket255RosterStillMatchesTheActualReadSites` (0.49s) | 红句＝voice 段多了一枚生产读者 `internal/panel/config_handlers.go`，名册 `cmd/wisp/config_readers_255.go` 未重裁；**台账 A878（20:36）已定位＝`89b3b550`/313-r1 引入、HEAD 级欠账**；本腿单跑复现同句（确定性、非争用） |
| `internal/ball` | `TestC21TableColourRowsMatchTokensCSS` | 读 `design/assets/tokens.css` 的那形；**工作树里 `design/` 16 枚删除态**（A878 具名："读不到文件那一形红一律具名归这条"，非任何腿） |
| `internal/panel` | `TestApprovalCardViewJSONKeysMatchFrontendTypes` · `TestComposerContractTypesMatchFrontend` · `TestNoGitSwitchCapabilityInThePanelSurface` · `TestC21DesignTokensFourWayAgree`（4 枚＝派单"历史在册 panel 4"）＋ `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8` · `TestStreamLogDroppedNamingLedgerIsNotBounded35r8`（2 枚） | 前 4 与 ball 同族（design/ 删除态）；后 2＝**35-r8 有意保留的红灯**（台账 13950："收下、不退回、不摘尺"）。⛔ 未顺手修 |
| `internal/risk` | `TestResolvePerCallBudget` (2.61s) | 争用型假红（派单点名）：整包跑时 `1.022ms/op` 越 1ms 预算；**安静单跑 `-count=3` rc=0**（复量一致） |
| `internal/session` | `TestSessionListByWorkspace` (0.63s) | **间歇红（本腿新量到、台账未见）**：安静单跑 `-count=5` 里 **2 发红**，红句＝workspace 过滤回的序与期望相反；机制＝`updated_at` 秒粒度并列＋`ORDER BY updated_at DESC, id DESC`，期望序假设创建序 ⇒ 票 166 那枚测试的并列不稳定。⛔ 不是本腿（本腿只碰 `internal/tools`），具名交回给 166/收尾链 |
| `internal/tools` · `internal/observe` | — | **ok 19.205s / ok 3.933s**（同一发终态里） |

⇒ 逐名比红名册：**本腿新增红＝0 枚**（`internal/tools`＋`internal/observe` 两包在终态里绿；其余 5 包的红逐名归到 255/313 漂移、design/ 删除态、35-r8 有意红灯、risk 争用、166 并列不稳定这五族）。

## 门禁

| 尺 | 改前 | 改后（交件前最后一跑） |
|---|---|---|
| `GOFLAGS= go build ./...` | rc=0（overlay=HEAD 版） | **rc=0** |
| `gofumpt -l` 工作树（cmd internal tools frontend） | — | 5 枚既有脏名（`cmd/wisp/models.go` `cmd/wisp/panel_transport_35r2_test.go` `internal/agent/approval/pending_read.go` `internal/agent/tools.go` `internal/risk/provenance.go`）＝**都不是本腿的件**；我这 3 枚空输出 |
| `gofumpt -l` HEAD blob（`git archive`） | — | 同名单里 HEAD 只有 `cmd/wisp/panel_transport_35r2_test.go` 一枚 ⇒ 另 4 枚是别的腿/工作树的老脏，⛔ 不算我的 |
| `sh scripts/d22scan.sh` | — | **rc=0**，`clean - no D22 ban violations`（ban #8 扫 internal/ 543 枚、cmd/ 123 枚，含注释与 _test.go） |
| `go test ./internal/tools/... ./internal/observe/...` | rc=1（红名＝`Test229FinishWatcherIsRosterNamedWhileChildRuns` 5.00s＋普查尺红；`internal/observe ok`） | **rc=0**（tools ok 15.953s / observe ok 3.803s） |
| `go test ./cmd/wisp ./internal/... -count=1` | — | 见 AC#6 |

## 没做到 / 具名交回

1. **4 枚"有意名册外"站点不是本票地界的活**（`panel-sta` 须 cmd/wisp/panel 侧决定或走 D38b 批准；`balldebug`×2 归票 66；`mockllm-stdout-reader` 归票 67）——本腿只登记，⛔ 没动 `cmd/wisp/**`、没往名册加名。
2. **没测真 LLM / 真 CLI 端到端**：判据仍从冻结接缝注入（C5 假 provider、真 bridge/真 registry/真 agent.Loop）；`wisp run` 那一条"模型自己派子代理"的路径本腿没跑。
3. **普查射程边界**：跳过 `_test.go`（有意假名字的站点已在上面点名）；`scripts/spike/**` 零 observe 调用但未逐行读；跨注册表的"孩子私有 registry 不可见"一条只按码读（`:297`），没有把它做成判据。
4. **AC#6 里两枚非本腿的红具名交回**：①`cmd/wisp TestTicket255RosterStillMatchesTheActualReadSites`（台账 A878 已定位 313-r1 漂移，单跑复现）；②`internal/session TestSessionListByWorkspace`（本腿新量到的间歇红：`-count=5` 里 2 发红，秒粒度 `updated_at` 并列＋`id DESC` 兜底，期望序假设创建序）——⛔ 都不在本腿写面，未修。
