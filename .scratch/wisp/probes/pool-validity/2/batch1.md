# pool-validity-2 / 批一 —— 零勾开放票有效性普查（号段 1–49，30 枚）

> 只读普查腿 `pool-validity-2`，接 `pool-validity-1`（死腿）的活。
> 起手 HEAD `949a5b92`（2026-10-06 09:03 +0800）／落笔时 HEAD `b3d29b0d`（09:06）。⚠ 换 HEAD 要重量。
> 本腿**零 go 命令**（并行腿 `231-r1` 此刻真在 `cmd/wisp` 跑定向用例与整包读数）。
> 工具全集＝`git log`／`ls`／`grep`／`wc`／`find`／Read。
> 票面零改动：零翻勾、零改名、零 `Status:` 改写、零产码改动、零台账写入。
> 唯一写入路径＝`.scratch/wisp/probes/pool-validity/2/**`。⛔ `pool-validity/1/**` 只读只引。
> ⛔ 零读零引 `frontend/**`、`design/**`、`cmd/wisp/**`、`scripts/**`、`.github/**`、
> `docs/reports/HANDOVER.md`、`docs/evidence/s1/**`、`.scratch/wisp/probes/{231,268,111,evidence-close}/**`。

## §0 四档尺（本腿判法；比 pv1 多一档）

| 档 | 判据 | 复法（必须留命令原文＋读数） |
|---|---|---|
| **仍成立** | 票面点名的那枚缺陷今天还在树上 | 把它 §现场／AC 里**最硬的一条**自己跑尺：`grep -n '<那句字面或函数名>' <它指的文件>`，或读那个行号今天写什么 |
| **已失效／已被别人做掉** | 缺陷不在了，且被别的票顺手做掉 | ⛔ 硬门：必须给一枚 commit 号＋`git log --name-only` 命中，**或**今树 0 命中；拿不出就不许写这一档 |
| **差翻勾** | 缺陷不在了，但凭据是"票内注记／台账 `A##`／一枚非实现者裁决表" | 凭据种类写进复法列。此档补 `pool-2` §3.0 指出的分类学洞；**不等于**本腿裁它可结案 |
| **量不到** | 判据落运行期行为／真开窗／真机／`frontend/**`·`design/**`（本腿禁读） | 归口写清缺哪一行读数；⛔ 不许由 grep 命中外推成"这功能能用" |

★ `上次动过` 列＝`git log -1 --format='%h %ad' --date=format:'%m-%d' -- <票文件>`。
本腿已把 105 枚这一列**整列重跑一遍**，与 pv1 名册逐枚比对＝**105/105 全等**（见 `summary.md` §0b）。

## §1 名册与判定（30 枚）

★ 尺面统一为**提交面**：`git --no-pager grep -n <pat> HEAD -- <path>`（工作树期间有 4 枚 `cmd/wisp` M 中，
用 HEAD 才不被脏面洗）。落在 `cmd/wisp/**` 的探针一律 `grep -l`（只取文件名，不读内容）。
"0 命中"＝该尺无输出。`[x]`＝票号列已坐实 `ls .scratch/wisp/issues/ | grep '^<号>'` 真名对上。

| 票号 | 票文件 | 上次动过 | 判定 | 复法（跑的那把尺） | 读数 |
|---|---|---|---|---|---|
| 15 | `15-speech-engines-cer-harness.md` | da81a6ac 09-19 | **仍成立** | `git grep -n 'AsrEngine' HEAD -- internal/ cmd/`；`git grep -n 'OfflineRecognizer\|paraformer\|Recognize' HEAD -- internal/ cmd/`；`git ls-tree -r --name-only HEAD \| grep -icE 'near_clean\|noisy_far\|asr-baseline'` | **1 命中**＝`internal/speech/doc.go:1`（注释里那句"owns … AsrEngine"）；三枚引擎符号 **0 命中**；CER 基线 wav 集 **0 枚**（`ls-tree` 计数 0）⇒ AC#1/#2/#3/#4 的载体根本不存在 |
| 16 | `16-s2-acceptance.md` | 144151ae 09-19 | **量不到** | 判据尺＝S2 切片验收（真链路 CER/延迟/热插拔/失败预演） | 归口：需要跑真链路＋C8 wav 注入的验收腿；本腿零 go 命令。**上游票 15 本腿判〔仍成立〕**⇒ 这张卡今天不可能绿，但"绿不绿"不是静态尺能读的 |
| 21 | `21-approval-gates-minimal.md` | 4a75482b 09-20 | **量不到** | `git grep -n 'D43: 17\|D43: 2[1-4]' HEAD -- internal/statemachine/table.go`；`git grep -ln 'wisp/internal/agent/approval' HEAD -- cmd/wisp/`（**只取文件名，内容未读**）；`git ls-tree -r --name-only HEAD cmd/wisp \| grep -i approval`（同上） | D43 #17/#21/#22/#23/#24 **全部有行**（`table.go:118,123`＝#17；`:150`＝#21；`:154`＝#22；`:158,163`＝#23；`:168,172`＝#24）⇒ AC#5 那一格静态成立；票的 `approval` 包在 `cmd/wisp` 有 **4 枚非测试引用者**（文件名：`approval_reply.go`/`resident_approval_windows.go`/`resident_task_source_windows.go`/`run.go`），另有 `approval_always.go`/`approval_reply_stdin_{windows,other}.go` 三枚同族文件 ⇒ 票内注记"this package has NO production caller yet"（`internal/agent/approval/doc.go:27`）与"internal/tools 零外部引用者"两句**都已过期**（由票 201/224/246 落地，见 §3 第 2 条）。但 segment 2 的三格（原生卡渲染／球体 Confirming 脉冲＋深度徽标／loop 重接）判据落 `cmd/wisp/**` 内容面与真开窗 ⇒ 归口：可读 `cmd/wisp` 且能开真窗的腿。⛔ 本腿不据文件名断言"渲染已存在" |
| 22 | `22-web-tools-d30.md` | 144151ae 09-19 | **仍成立** | `git grep -n '"web.search"' HEAD -- internal/ cmd/`（排 `_test`）；`git grep -ln 'web.open\|web.fetch' HEAD -- internal/tools/`；`git grep -n 'func (.*) Name() string { return' HEAD -- internal/tools/` | `"web.search"` 非测试命中 **1 处**＝`internal/risk/provenance.go:137`（C25 污名表里的一个名字，不是工具）；`web.open/web.fetch` 在 `internal/tools/` **0 命中**；`internal/tools/` 全部 `Name()` 返回值只有 `fs.read/fs.list/fs.edit/fs.write/fs.trash/fs.move/fs.delete/task.output/task.cancel/task.spawn` ⇒ 本票的四件套工具不存在，SSRF 套件无处可跑（`internal/risk/rules_network.go` 那套 R5 私有段拒绝是票 17 的评估器，不是 fetch 腿）|
| 23 | `23-system-window-input-tools.md` | 144151ae 09-19 | **仍成立** | 同 22 的 `Name()` 尺；`git grep -n '"screen.capture"\|"input.type"\|"window.close"\|"system.power"\|"clipboard.read"' HEAD -- internal/ cmd/`（排 `_test`） | 五个冻结名非测试命中 **4 处、全在 `internal/risk/provenance.go`**（`:86 clipboard.read`、`:89 doc.read`、`:90 screen.capture` 等污名常量）＋`window.close`/`system.power`/`input.type` **0 命中**；`Name()` 表里一个都没有 ⇒ AC#1–#5 的实现体不存在 |
| 24 | `24-doc-search-tools.md` | 144151ae 09-19 | **仍成立** | `git grep -n '"doc.read"\|"search.content"\|"search.files"\|"search.apps"' HEAD -- internal/ cmd/`；`git grep -c 'DEFERRED(D-34)' HEAD -- internal/` | 四个冻结名非测试命中只 `internal/risk/provenance.go:89`（`doc.read` 一枚污名常量），`search.*` 三枚 **0 命中**；AC#5 点名的 `DEFERRED(D-34)` 标记 **0 命中** ⇒ 工具与标记两头都不在树上 |
| 25 | `25-s3-acceptance.md` | 144151ae 09-19 | **量不到** | 判据尺＝S3 全量红队 + 能力场景①② | 归口：运行期＋真开窗验收腿。上游 21/22/23/24 本腿读数＝21 量不到、22/23/24 〔仍成立〕 |
| 26 | `26-tts-output.md` | 6e2a68fa 09-20 | **仍成立** | `git ls-tree -r --name-only HEAD -- internal/speech`；`git grep -n 'DEFERRED(playback)' HEAD -- internal/audio/doc.go`；`git grep -n 'TtsEngine' HEAD -- internal/ cmd/` | `internal/speech/` 提交面**只有 `doc.go` 一枚**；`internal/audio/doc.go:27` 仍逐字写着 `DEFERRED(playback): TTS output lands with ticket 26 on the same WASAPI`；`TtsEngine` 非测试 **1 命中**＝`internal/speech/doc.go:2`（注释）⇒ TTS 播放腿没落 |
| 27 | `27-punctuation.md` | 144151ae 09-19 | **仍成立** | `git grep -n 'Punctuation' HEAD -- internal/ \| grep -v internal/config` | **0 命中** ⇒ 标点开关只在 `internal/config/schema.go:252` 有字段、`internal/` 里除配置包外零消费者；引擎侧（票 15 的 ASR 链）本腿读数＝不存在 |
| 28 | `28-session-scope-warm.md` | 130c0943 09-19 | **仍成立** | `git grep -n '^func' HEAD -- internal/session/session.go`；`git grep -n 'EvWarmIdle\|EvSettleExpired' HEAD -- internal/ cmd/` | `internal/session/session.go` 全文**只有 3 枚 func**（`ID.String:56`/`ID.Valid:67`/`Mint:94`）＝会话身份，无生命周期；`EvWarmIdle`/`EvSettleExpired` 命中面＝`events.go:50,51`（定义）＋`table.go:227,231,236`（D43 行）＋三枚 `_test` ⇒ **生产侧无人发这两枚事件**，90s/3s 定时无人按。边界：发射器若要落在 `cmd/wisp`（禁读）本腿判不到，但那个包连自己 `doc.go:9` 认领的"Warm keepalive + 90s→Settling trigger"都没实现 |
| 29 | `29-memory-l1l2.md` | 144151ae 09-19 | **仍成立** | `git grep -n 'func.*Extract' HEAD -- internal/`（排 `_test`）；`git grep -n 'memory.save\|memory.recall' HEAD -- internal/ cmd/`（排 `_test`） | 非测试 `Extract` 命中只有 `internal/memory/models.go:38 ProfileSource.Valid()` 与 `internal/models/`（tar 解包，无关）⇒ **L1 抽取流水线零实现**；`memory.save/recall` **0 命中**（工具不存在，见 22/23 那把 `Name()` 尺）；`profile`/`memory` 两枚表的 DAO 有（`dao_profile.go:21 UpsertProfile`、`dao_memory.go:20 AddMemory`、`dao_profile.go:67 evictProfileLRU`）——那是票 04 落的存储层，不是本票的抽取与反馈环 |
| 30 | `30-result-routing-d10.md` | 144151ae 09-19 | **仍成立** | `git grep -in 'resultTier\|routeResult\|toast' HEAD -- internal/`（排 `_test`）；`git grep -n 'clipboard.write' HEAD -- internal/ cmd/` | 三符号非测试 **0 命中** ⇒ D10 四档分派器不存在；`clipboard.write` 只在 `internal/risk/provenance.go:30,34` 的名字表里（注释＋冻结名清单），无写入腿 ⇒ AC#1 的档位矩阵无处可测。票 08 落的是 `artifacts` 配额（`internal/memory/artifacts.go:78`），那是本票 AC#3 的一半，另半（分档投递）不在 |
| 31 | `31-reminders.md` | 144151ae 09-19 | **仍成立** | `git grep -n 'CREATE TABLE' HEAD -- internal/memory/schema.go`；`git grep -in 'reminder' HEAD -- internal/memory/` | 九枚表＝`schema_meta:23/profile:30/memory:40/task_log:53/tool_call:69/approval_grant:89/cost_daily:102/plugin_state:111/provider_health:136` ⇒ **无 reminder 表**；`reminder` 在 `internal/memory/` **0 命中**。`internal/agent/guard.go` 那组 `Reminder` 是 C22 重复调用梯度提醒（`guard.go:45`），与本票的定时提醒同名不同物，不能当凭据 |
| 32 | `32-s4-acceptance.md` | 130c0943 09-19 | **量不到** | 判据尺＝S4 场景③④（听写主观分／陪聊全双工） | 归口：主观签核＋真机音频＋运行期延迟，全部越出静态尺；上游 27/29/30/31 本腿读数＝四枚全〔仍成立〕 |
| 34 | `34-frontend-scaffold.md` | 144151ae 09-19 | **量不到** | — | 判据 100% 落 `frontend/**`（tsc/lint/bundle ≤400KB gz/token parity/明暗渲染）＝本腿禁读面。归口：可读前端面的腿 |
| 35 | `35-panel-bridge-c17.md` | 144151ae 09-19 | **仍成立** | 逐名尺：`for m in panel.resync tasks.list history.query transcript.get approval.current approval.queue approval.decide config.set grants.list grants.revoke privacy.purge privacy.export cost.summary models.list diagnostics.export; do git grep -l -F "$m" HEAD -- internal/panel/ \| wc -l; done`；再对两枚有命中的核一遍 `git grep -l -F 'approval.decide' HEAD -- internal/panel/` | 15 枚冻结方法名里**只有 2 枚有文件**（`config.set` 4 files：`bridge.go`/`config_handlers.go` 非测试＋2 枚 `_test`；`approval.decide` 5 files **全是 `_test`**＝`bridge_test.go`/`composer_test.go`/`config_route_248_test.go`/`frontend_hygiene_test.go`/`l2_grant_boundary_test.go`，非测试 0 处）⇒ 桥实际应答的方法只有 `bridge.go:42-67` 那 6 枚（`panel.mode.request`/`panel.workspace.request`/`panel.attachment.add`/`panel.message.send`/`config.get`/`config.set`），**冻结白名单 15 缺 13**；`resync` 在整个 `internal/` 非测试 **0 命中** ⇒ AC#1（每张卡都过桥）/#4（resync 全量推）前提不成立 |
| 36 | `36-result-history-panel.md` | 15ff2bee 09-26 | **量不到** | — | 判据＝聊天面 XSS 套件/`dangerouslySetInnerHTML` grep/60fps 滚动/真渲染 ⇒ `frontend/**` 禁读 + 运行期。归口：前端面可读腿 |
| 37 | `37-approval-ui-l2.md` | 15ff2bee 09-26 | **量不到** | `git grep -n 'Depth\|position' HEAD -- internal/agent/approval/queue.go`（Go 侧半边）；`git ls-tree -r --name-only HEAD cmd/wisp \| grep -i approval`（只取文件名） | Go 侧深度/队头语义**在**（`queue.go:148 Depth()`、`:283 position()`、`:577 Depth: q.position(it)`）；AC 六格里"卡片渲染/深度徽标<100ms/DOM 无 allow 接线/人工视觉签核"落前端与运行期 ⇒ 归口：能开真窗的腿。⚠ 别把 Go 半边读成"这票做完了" |
| 38 | `38-command-palette-tasks.md` | 15ff2bee 09-26 | **量不到** | — | 判据＝命令面板 IME 组合提交/键盘全程/发现页风险药丸 ⇒ 前端＋运行期。`internal/panel/composer*.go` 有 Go 侧 composer 派发（票 92 落的），不构成对本票任一 AC 的静态凭据 |
| 39 | `39-config-editor-gui.md` | 15ff2bee 09-26 | **量不到** | `git grep -n '^func' HEAD -- internal/panel/config_handlers.go \| wc -l`；`git grep -c 'lockedFieldFamilies' HEAD -- internal/panel/config_handlers.go` | Go 侧有 14 枚 func、`config_handlers.go:108` 有 `lockedFieldFamilies` 锁字段族（票 92/后续票的地界）；但票面 AC 的编辑器形态/三档生效提示落 GUI 面 ⇒ 量不到，归口前端腿 |
| 40 | `40-security-privacy-cost-pages.md` | 2df182e1 09-28 | **量不到** | — | 判据＝三张设计页视觉签核 + 隐私/成本页 e2e 打真库 ⇒ 前端禁读 + 运行期。上游 44/49 本腿读数＝〔仍成立〕 |
| 41 | `41-kws-wake-word.md` | 130c0943 09-19 | **仍成立** | `git grep -in 'KeywordSpotter\|WakeWordEngine' HEAD -- internal/ cmd/`（排 `_test`）；`git grep -n 'SetLoaded' HEAD -- internal/`（排 `_test`） | `KeywordSpotter` **0 命中**；`WakeWordEngine` 非测试 **1 命中**＝`internal/speech/doc.go:2`（注释）；`SetLoaded` 在 `internal/` 只有 `approval/approval.go:304,306`（定义，注释写着"the KWS loader calls it"）＋零生产调用者 ⇒ 票 21 那句"KWS 未加载时须显式说「语音取消不可用」"今天仍是**如实报不可用**，AC#1–#5 无处可测 |
| 42 | `42-watchdog.md` | 144151ae 09-19 | **仍成立** | `git ls-tree -r --name-only HEAD -- internal/watchdog`；`git grep -n 'DEFERRED(watchdog' HEAD -- internal/watchdog/doc.go` | 提交面**只有 `doc.go`**；`doc.go:18` 仍逐字挂着 `DEFERRED(watchdog loop/thresholds): implemented by ticket 42` ⇒ 采样环/阈值表/告警动作零实现（`proc.TreePrivateBytes`、`observe` 侧的 SLO 管道是票 08/proc 的活，不是本票的门限环） |
| 43 | `43-power-lifecycle-events.md` | 144151ae 09-19 | **仍成立** | `git grep -ln 'WM_QUERYENDSESSION\|WM_POWERBROADCAST\|WTS_SESSION_LOCK\|WTS_SESSION_UNLOCK\|RegisterSessionNotification' HEAD -- internal/ cmd/`（只取文件名） | **0 files**（整棵提交树，含禁读的 `cmd/wisp` 也只按文件名匹配）⇒ 挂起/唤醒/会话结束/关机四类系统事件**一个接收者都没有**。连带后果：票 48 的"挂起时作废待决审批"与票 21 的"挂起＝拒绝"那格（`§16.11#10`）今天没有事件源可挂 |
| 44 | `44-costmeter-c23.md` | 5cba2d92 09-19 | **仍成立** | `git grep -n 'QuotaDailyMicro\|PlanCreditTotalMicro' HEAD -- internal/ \| grep -v internal/config`；`git grep -n '^func' HEAD -- internal/agent/cost.go` | 两枚配额字段在 `internal/config/` 之外 **0 命中**＝有字段无消费者 ⇒ 日/月配额门与计划额度记账（AC#5 那半格）没落；已落的是计价与聚合两枚原语（`agent/cost.go:38 AddUsage`、`:50 PriceMicros`、`memory/dao_misc.go:132 BumpCostDay`）＝票 04/10 的存储与算术腿，不是 C23 的门禁面；AC#1/#3/#4/#6 需 fixture 时钟与运行期 ⇒ 那几格量不到，但本票最硬的那条（配额无消费者）静态坐实 |
| 45 | `45-diagnostics-guards.md` | 144151ae 09-19 | **仍成立** | `git grep -ln 'GetDiskFreeSpace\|diskFree' HEAD -- internal/ cmd/`；`git grep -iln 'focusassist\|focus_assist' HEAD -- internal/`；`git grep -n 'wal_checkpoint' HEAD -- internal/memory/open.go` | 磁盘余量守卫 **0 files**；专注助手 **0 命中**；WAL 侧只有 `open.go:46 PRAGMA wal_checkpoint(TRUNCATE)`（启动期一次性，票 04 落的）＋`:45 pragmaWALAutocheckpoint=1000` ⇒ 票面 AC#2（写满门限→停 L3+停日志+可恢复）与 AC#4（Focus-Assist 降级结论）零实现。已落的一半：AC#1 的导出器本体在（`observe/diagnostics.go:62 BuildDiagnosticsBundle`，票 08 unit A `739bb15f`）＋artifacts 配额（`memory/artifacts.go:81 enforceArtifactsQuota`） |
| 46 | `46-s6-acceptance.md` | 144151ae 09-19 | **量不到** | 判据尺＝S6 空闲/KWS/功耗/成本验收 | 归口：运行期＋真机。上游 41/42/43/44/45 本腿读数＝五枚全〔仍成立〕 |
| 47 | `47-task-scheduler-pathlock.md` | 144151ae 09-19 | **仍成立** | `git ls-tree -r --name-only HEAD -- internal/agent/scheduler`；`git grep -ln 'PathLock\|pathlock' HEAD -- internal/ cmd/`；`git grep -c '"task.list"' HEAD -- internal/` | `internal/agent/scheduler/` 提交面**只有 `doc.go`**，`doc.go:13` 挂着 `DEFERRED(scheduler): implemented by ticket 47`；`PathLock` 名字只出现在 `scheduler/doc.go`、`statemachine/events.go`、`statemachine/table.go`（D43 #20 那行状态转移）⇒ **锁本体与冲突排队零实现**（转移表有边、没人走它）；`task.list` 非测试 **0 命中**（AC#6 的另一半）；已落的一半：D38d 并发上限在 `tools/bridge.go:19 MaxToolConcurrency`、`task.output`/`task.cancel`/`task.spawn` 三工具在（票 20/197/221 地界） |
| 48 | `48-approval-queue-full.md` | 144151ae 09-19 | **仍成立** | `git grep -n 'DefaultMaxPending\|pending \[\]\*qitem\|func (q \*Queue) position' HEAD -- internal/agent/approval/queue.go`；`git show 35200c75:internal/agent/approval/queue.go \| grep -n 'DefaultMaxPending'`；`git grep -n 'D43: 2[3-5]' HEAD -- internal/statemachine/table.go` | 队列**已经是多任务 FIFO**：`queue.go:84 pending []*qitem`、`:127 DefaultMaxPending = 8`（溢出 `:159` fail-closed 拒）、`:282 position()`/`:577 Depth` 队头单显＋深度 ⇒ **⚠ 票面那句前提"upgrading 21's trivial queue（单项队列）"已过期**，见 §3 备注栏；仍成立的硬格：挂起作废——`abandon()` 只有两个调用者（`gate.go:537` 界面不可达、`gate.go:579` 任务 ctx 结束），**没有挂起事件源**（票 43 读数＝四类系统事件 0 files）；D43 #23–25 转移表**有行**（`table.go:158,163,168,172,176,180`）＝那半格静态成立；AC#1/#2/#5 的徽标同步与 270s 可见告警需运行期 ⇒ 那几格量不到 |
| 49 | `49-session-grants-d45-2.md` | 144151ae 09-19 | **仍成立** | `git grep -c 'grant_id' HEAD -- internal/memory/schema.go internal/memory/dao_toolcall.go`；`git grep -ln 'grants.revoke\|grants.list' HEAD -- internal/panel/`；`git grep -n '"input.type"' HEAD -- internal/ cmd/` | 已落的一半（**不是本票落的**，见 §3）：会话身份＋记账面在 `internal/session/grants.go`（`:146 Record`、`:189 Covering`、`:322 patternCovers`），审计链 `schema.go:83 grant_id`＋`dao_toolcall.go:25,52,160` 有列有写；仍成立的硬格：`grants.revoke`/`grants.list` 在 `internal/panel/` **0 files**＝撤销入口没有；`"input.type"` 非测试 **0 命中**＝目标进程绑定的那枚变体（票面三条不可破限制里的第三条与 AC#4）无处存在；AC#1 的"会话结束即过期"按其自述设计＝无代码（`grants.go:31-34` 注释明写"nothing"），那一格本腿不判 |

## §2 本批计数

| 档 | 枚数 | 票号 |
|---|---|---|
| 仍成立 | **19** | 15 22 23 24 26 27 28 29 30 31 35 41 42 43 44 45 47 48 49 |
| 已失效／已被别人做掉 | **0** | —（硬门一枚没过：本批没有"缺陷字面今树 0 命中＋顶掉它的 commit 号"两者齐的票） |
| 差翻勾 | **0** | —（见 §3 末：48/49 有"被别人的票顺手做掉半格"的实料，但**整票缺陷不在**这件事无凭据） |
| 量不到 | **11** | 16 21 25 32 34 36 37 38 39 40 46 |
| 合计 | **30** | ✓ 与 §0 分母对上 |

## §3 与 pv1 的逐枚对照 ＋ 本批三处具名读数

**先说结构**：pv1 盘上 105 行判定**全是"未判"**（`grep '^| [0-9]' batch{1,2,3,4}.md | grep -v '未判'` → 0 命中，
且 `git diff --stat HEAD -- .scratch/wisp/probes/pool-validity/` 为空）。⇒ **本批 30 枚里没有一枚可以和 pv1 对判语**，
"它判 X／我判 Y"这一栏对批一是**空的**，不是因为两腿一致，是因为 pv1 那一列从来没落到盘上。详见 `summary.md` §0c。

三处本腿自己的具名读数，编排者排产时值得单独看一眼：

1. **票 48 的票面前提过期**：它写"upgrading 21's trivial queue"，而 21 自己那一笔 `35200c75`（09-20）
   落进盘的 `queue.go` 就已经是 `pending []*qitem` + `DefaultMaxPending = 8` + `position()` 队头单显
   （`git show 35200c75:internal/agent/approval/queue.go` 里 `:78` 与 `:97` 两行即证）。
   ⇒ 这不改本腿判语（它仍是〔仍成立〕，挂在挂起事件源与运行期那几格上），但**"还剩 30 枚里第 48 号"这句
   给人的重量感是错的**：那枚"单项→多项"的活早就做完了。
2. **票 49 半格由票 224 顺手做掉**：`3b78c246`（09-30 12:00，`feat(224)`）的 `--name-only` 命中
   `internal/session/grants.go`＋`internal/tools/grant.go`＋`internal/tools/bridge.go`＋
   `internal/agent/approval/{gate,queue,replies,ui}.go`＋`cmd/wisp/{approval_reply,run}.go`。
   ⇒ 本腿**不**因此把 49 降进〔已失效／已被别人做掉〕：224 做的是"这一档没人执行"，
   49 点名的 `grants.revoke`（0 files）与 `input.type` 目标进程绑定（0 命中）今天仍不在。
   这正是 pv1 想合并、本腿坚持分列的那一格差别——〔差翻勾〕要的是"缺陷不在了"，这里缺陷只少了一半。
3. **票 21 挂起那格与票 43 是同一枚洞的两个入口**：`§16.11#10`（挂起时待决审批作废）在 21/48 两张票面上都要求，
   本腿在 43 名下量到"四类系统事件 0 files"。⇒ 编排者若按票号排产，会以为 21、43、48 是三件活；
   静态尺读到的是**一件**（事件源）＋两件各自的接线。这条不算判定，算读数。

## §3 本批打算答什么

这 30 枚是**号段最老的一批**（28/30 枚票面最后动过 ≤09-28，其中 22 枚停在 09-19/09-20 那一笔
`144151ae`/`130c0943` 上）。它们承载的是 S2–S6 的**切片验收**与**当初一次写完就没再回来的建功能票**。
编排者要的那个数在这一批最难看：这一批里既有"整片能力早被后来的票连着做掉"的候选（跑尺坐实），
也有"验收票本身要跑真链路、静态尺根本判不了"的一批（落第四档）。**所以本批的三档分布本身就是答案**：
如果 1–49 里"已失效+差翻勾"占大头，那"还剩 105"就是假数。
