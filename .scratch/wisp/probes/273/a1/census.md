# 273-a1 普查件 — D43 副作用的名册、收件人与三形代价（只读腿，零产码改动）

leg id `273-a1`｜落点普查，⛔ 本腿不落地、不选形｜票面
`.scratch/wisp/issues/273-shipping-process-builds-the-state-machine-without-a-sink-so-every-d43-side-effect-falls-into-a-no-op.md`（38 行／5464 字节，**整份读完，无一行截断**）

---

## 0. 起手锚与具名声明

- 起手 `git rev-parse --short HEAD` = `c318337`（全称 `c3183379529dd7498c9b8d65d2fe58cc22c20951`）。
- 起手 `git status --porcelain | wc -l` = **758**（本腿开始前）。
- ⛔ **本腿一行产码没改**：全程只 `Read`／`grep`／`awk`／`wc`／`ls`／`git show`；本件与 `logs/` 是本轮唯一新增文件。
- ⛔ **本腿没跑过任何 `go test`／`go build`／`go list`／突变**（`272-r2` 在飞，包级互斥）。所有"能不能过"的话都是**读码推的**，逐处标注〔读码推〕／〔现量〕。
- **`.github/workflows/ci.yml` 一律经 `git show HEAD:.github/workflows/ci.yml` 读**，没碰工作副本（`111-r5` 在编辑它）。CI 摘要取到 `logs/ci-from-gitshow.txt`（55 行，只留 `go build|go install|go test|gofmt|d22scan|go vet` 命中的行）。
  - 现量读数：`git show HEAD:.github/workflows/ci.yml | grep -c balldebug` = **0** ⇒ CI 不构建 `cmd/balldebug`。
- ⛔ 不读不写 `frontend/**`、`design/**`；没开窗、没加 `-tags winlive`；无密钥入文。
- 原始尺输出全部落在 `.scratch/wisp/probes/273/a1/logs/`（**29 枚文件**，整目录 226K→238K），本件只带摘要与逐名表。
- 第 2 笔（`ae6c010e` 之后）追加三把防漏尺：`bypass-and-tag-checks.txt`（复合字面量绕过＝0、`cmd/wisp` 里 `internal/ball` 的构建标记）、
  `state-for-card-level.txt`（`stateForCardLevel` 全部命中）；据此修正了 §2 甲栏"可移植面够不到界面收件人"与 §2 丙栏的逐名钉。

---

## 1. 副作用名册（AC#0 ＋ "有没有真收件人"）

### 1.1 枚数：**50 枚**（不是转述，是本腿现量，两把独立尺同数）

| 尺 | 命令形状 | 读数 |
|---|---|---|
| A（表本体） | 从 `internal/statemachine/table.go` 抽 `SideEffects: []string{…}` 里的带引号名，`sort -u` | **50 枚不同名**／**62 次投递实例**／39 枚行带副作用 |
| B（测试侧独立副本） | 从 `internal/statemachine/table_test.go` 抽 `effects: []string{…}`，`sort -u` | **50 枚**，且与尺 A **逐名 `diff` 为空＝集合全等** |

- 表的形状：`table.go` 里 `D43:` 标号出现 **56** 次，去重后编号 **1–42**（42 枚编号），
  而 `docs/PLAN.md` §D43（`:3059-3100`）那张权威表只有 **40 行（编号 1–40）**。
  差的两枚编号 = `#41`（打断／barge-in，`table.go:195`、`:204`）、`#42`（`Conversation` + 任务意图，`:246`），
  包自己的 `internal/statemachine/doc.go` 写明是"the 40 frozen rows plus **SPEC-08 §3's appended #41/#42**"。
  ⇒ **这是"权威表枚数"与"实现表枚数"的一处真实分叉**，本腿不动它，写进 §5／§4 请编排者裁。
- 与编排者手里那枚"50"：**本腿现量也是 50，顶不回来**；但那枚"50"的**来源句**（票面 §2③ 转述"50 枚副作用名全投进 no-op"）
  的后半截 **"全投进 no-op" 只对了 49／50 之外的一处细节——见 1.3**，另外票面 §2② 说"其余 **5** 枚命中都在 `_test.go`"，
  本腿现量是 **6** 枚（票面自己下面又逐名列了 6 处：`hotkey_live_test:382`、`interaction_live_test:54`/`:134`、`bridge_test:16`/`:86`、`handoff_window_109_test:62`）。**枚数错一枚，名单没错。**

### 1.2 真收件人：**0 枚**；只有空罐子：**50 枚**（现量）

判据锚在**调用形状**，⛔ 不按符号名裸算：

| 尺 | 命令形状（关键点） | 读数 |
|---|---|---|
| C（产码构造点） | `grep -rn 'statemachine\.New(' --include=*.go .` 去掉 `_test.go` 与 `.scratch/` | **2 枚**：`cmd/wisp/models.go:303`、`cmd/balldebug/main.go:188`；两枚都是**单行字面量、只有 `Initial`**（`logs/new-callsites.txt`，总 8 枚命中＝2 产码＋6 测试） |
| C′（防漏：多行字面量） | `grep -rn 'statemachine\.Options{' …` 去测试 | 同样 **2 枚** ⇒ 没有"多行 Options 把 Sink 藏在下一行"的漏网 |
| C″（防绕过构造） | `grep -rn 'statemachine\.Machine{' --include=*.go cmd internal` 去测试 | **0 枚** ⇒ 没有"不走 `New`、用复合字面量绕过缺省支"的第三条造机器路子 |
| D（有没有人给过真收件人） | `grep -rn 'Sink:\|Sink =' --include=*.go` 去测试与 `.scratch` | 58 枚命中里**没有一枚**是 `statemachine.Options` 的 `Sink`；`cmd/wisp/run.go:998 consoleSink`＝**`agent.Sink`**（`func (c consoleSink) Publish(e agent.Event)`，`run.go:1260`）、`cmd/wisp/providers.go:194 storeHealthSink`＝**`llm.HealthSink`**（`run.go:1170`）、`internal/agent/loop.go:209 NopSink{}`＝**`agent.Sink`** ⇒ **三枚同名 `Sink` 都不是状态机那枚**（本轮实测的第二处"同名罐子"陷阱，与记忆里 `.Handle(` 那例同族） |
| E（类型面） | `grep -rn 'statemachine\.Effect\|statemachine\.Sink'` 去测试去 `.scratch` | **0 枚** ⇒ 产码里**没有任何一段代码能点出这个类型名**，收件人在类型面上不存在 |
| F（字符串面，逐枚） | 对 50 枚名逐枚 `grep -rnF '"<name>"'` 去 `table.go`／去 `_test.go`／去 `.scratch` | 48 枚 **0 命中**；2 枚命中但**都不是消费者**（下条） |
| F 的两枚假阳性 | `logs/effect-name-consumers.txt` | ① `internal/agent/guard.go:46` 命中 `"agent.inject-reminder"` ⇒ 那是**注释**；② `internal/statemachine/events.go:60` 命中 `"error.ack"` ⇒ 那是**事件名** `EvErrorAck Event = "error.ack"`，与副作用名**撞字符串**，不是消费者（⚠ 以后任何腿拿裸 grep 数这个都会得到"2 枚有人接"的假读数） |

⇒ **AC#0 后半的答案：今天有真收件人的＝0 枚，只有空罐子的＝50 枚。**

### 1.3 但"50 枚全投进空罐子"这句话**在产码面上还不成立**——今天只有 1 枚真的被投出去过

追链（现量，逐跳 `文件:行`）：

- **链①（唯一一条出货链，走通到 `main.go`）**
  `cmd/wisp/main.go:109`（`case "models"` 于 `:102`）→ `cmd/wisp/models.go:104 cmdModels` → `:127 case "ensure"` → `:262 modelsEnsure` → `:293 store.handOffModel` → `:302/:303 statemachine.New(Options{Initial: StateFirstRun})`（**无 Sink**）→ `:305 models.WireDownloading(ms.mgr, machine)` → `internal/models/bridge.go:33` → `bridge.Run` 内的四次真投递：`bridge.go:40 EvModelMissing`、`:46/:59 EvDownloadFailed`、`:64 EvDownloadCompleted` → `internal/statemachine/machine.go:162-164` 逐枚 `m.sink(Effect{…})` → `machine.go:67` 的空函数。
  - 这条链上真会被投递的副作用名只有 **`model.verify-sha256-signature`**（`table.go:263-265` D43 #37 完成支，**1 枚**）。
    `#2`（`table.go:46`）**不带副作用**；`#37` 失败支（`:267`）**也不带**。
  - 次级：失败支会把机器停在 `Error`，`internal/statemachine/timeouts.go:52` 给 `Error` 装了 10s → `EvErrorAck` → `#38` → 投 `error.ack`。
    但 `cmd/wisp/models.go:304 defer machine.Close()` 与 `machine.go:169-174 Close()` 会在返回前拆表，〔读码推〕**今天投不出去**；
    记作**第 2 枚"结构上可达、实际赶不上"**，⛔ 本腿没跑任何计时读数。
- **链②（调试载具，不算出货）** `cmd/balldebug/main.go:188` → `:617 dispatch()` → `:618 m.Dispatch(...)`；`main.go:178` 持有机器。
  〔现量〕`cmd/balldebug` **不在 ci.yml**（0 命中）、只在 `scripts/dev/ball-cycle.ps1` 与 `scripts/check-path-length-budget.sh` 里出现
  ⇒ 标〔**非生产／调试载具**〕。它能把很多枚名投进空罐，但**产品用户看不见那条进程**。
- **链③（常驻腿）根本不存在**——本轮最硬的一条负向读数：
  - `grep -c 'statemachine' cmd/wisp/run.go` = **0** ⇒ `wisp run` 那条腿**不建机器**。
  - `cmd/wisp/resident_ball_windows.go:275` 那枚 `Initial: statemachine.StateSleeping` **不是状态机**，它属于 `ball.New(ball.Options{…})`（`:270-290`）；
    `internal/ball/*.go`（非测试）里 `statemachine.New`／`*statemachine.Machine` **0 枚**，61 处引用**全部只用 `statemachine.State` 这个词表**。
  - 球的视觉状态由 `cmd/wisp/resident_approval_windows.go:879 b.SetState(stateForCardLevel(p.Level))` 与 `:956 b.SetState(StateSleeping)` **直接点名**，
    不经过 D43 的边、也不经过任何 `Effect`。
  ⇒ **常驻腿今天既没有发件人（无机器实例），也没有收件人（无 Sink）**；票面 §1 那句"状态变了但没有任何东西收到过"
  在 `cmd/wisp` 的两条交互腿上是**更严重的一格**：**连"状态变了"都不是由这台机器变的**。

### 1.4 名册（50 枚逐名｜`行数＝投递实例数`｜`D43`＝编号｜收件人＝0/50）

分类列 `→谁` 是本腿按名前缀给的**读码推**归档（`ball`＝球视觉、`panel`＝面板、`speech/asr/tts/kws/mic/audio`＝语音链、`session/mem`＝会话与内存、`approval/badge/confirm`＝审批、`config/db/model`＝装配与分发、`error/diag`＝错误、`ui`＝界面提示），⛔ 不是任何仪器里的既有事实。
"今天真投递"列：`是`＝链① 出货路径上真会打到空罐；`结构可达`＝链① 上但今天赶不上；`调试`＝只在 balldebug 可达；`零`＝今天产码里没有任何派发能把机器推到那一条。

| # | 副作用名 | D43 | 实例 | 今天真投递 | →谁 |
|---|---|---|---|---|---|
| 1 | `agent.cancel-tool-call` | 22 | 1 | 零 | agent |
| 2 | `agent.inject-reminder` | 19 | 1 | 零 | agent |
| 3 | `approval.cancel-call` | 24,24 | 2 | 零 | approval |
| 4 | `approval.dequeue` | 23,23 | 2 | 零 | approval |
| 5 | `approval.enqueue-c18` | 17 | 1 | 零 | approval |
| 6 | `approval.replay-offer` | 24,24 | 2 | 零 | approval／panel |
| 7 | `asr.exclude-playback` | 41,41 | 2 | 零 | speech |
| 8 | `asr.punctuate` | 11 | 1 | 零 | speech |
| 9 | `audio.discard-buffer` | 13 | 1 | 零 | speech |
| 10 | `audio.release-output` | 27,41,41 | 3 | 零 | audio |
| 11 | `audio.stop-capture` | 11 | 1 | 零 | audio |
| 12 | `badge.decrement` | 23,23 | 2 | 零 | ball（`ball.SetBadge` 有真形，`ball_windows.go:316`，但没人从 Effect 叫它） |
| 13 | `config.persist` | 1 | 1 | 零 | config |
| 14 | `confirm.countdown-start` | 17 | 1 | 零 | approval |
| 15 | `conversation.entered` | 30 | 1 | 零 | ball／speech |
| 16 | `conversation.privacy-confirm-l2` | 30 | 1 | 零 | approval（⚠ L2 语义，见 §4） |
| 17 | `conversation.ring-solid` | 28 | 1 | 零 | ball |
| 18 | `conversation.suspend` | 42 | 1 | 零 | session |
| 19 | `ctx.carry-c7-d47` | 42 | 1 | 零 | agent |
| 20 | `db.create-schema` | 1 | 1 | 零 | store |
| 21 | `diag.record-stack` | 36 | 1 | 零 | observe／diag |
| 22 | `error.ack` | 38 | 1 | 结构可达 | 状态回投 |
| 23 | `error.classify-retry` | 16 | 1 | 零 | observe／agent |
| 24 | `error.device-name` | 14 | 1 | 零 | panel |
| 25 | `input.taint-mark` | 11 | 1 | 零 | D30 注入防护 |
| 26 | `kws.keep-alive-alert` | 9 | 1 | 零 | speech／panel |
| 27 | `kws.load` | 5 | 1 | 零 | speech |
| 28 | `kws.pause` | 7 | 1 | 零 | speech |
| 29 | `kws.stop-inference` | 8 | 1 | 零 | speech |
| 30 | `mem.free-os-memory` | 31 | 1 | 零 | runtime |
| 31 | `mem.rss-verify-10s` | 32,32 | 2 | 零 | observe（D32 SLO） |
| 32 | `mic.enable` | 28 | 1 | 零 | speech |
| 33 | `mic.keep-off` | 26 | 1 | 零 | speech |
| 34 | `model.verify-sha256-signature` | 37 | 1 | **是** | models（**今天唯一真投递**） |
| 35 | `panel.destroy-or-hide` | 31 | 1 | 零 | panel |
| 36 | `panel.keep-alive` | 26 | 1 | 零 | panel |
| 37 | `panel.queue-waiting` | 20 | 1 | 零 | panel |
| 38 | `panel.stream-push` | 15,15 | 2 | 零 | panel |
| 39 | `session.dispose-scope` | 31 | 1 | 零 | session |
| 40 | `session.scope-create` | 4,7 | 2 | 零 | session |
| 41 | `session.zero-load-resume` | 29 | 1 | 零 | session（B4） |
| 42 | `settling.cancel-fallback` | 33 | 1 | 零 | watchdog |
| 43 | `speech.load-vad-asr` | 4,7 | 2 | 零 | speech |
| 44 | `speech.unload-asr-tts` | 31 | 1 | 零 | speech |
| 45 | `task.ctx-keep` | 34 | 1 | 零 | agent |
| 46 | `tools.execute` | 21 | 1 | 零 | tools（⚠ 见 §4：这枚名一旦变成"可被调用的东西"就撞上"宿主内部 artifacts／受门控 Tool"那条禁令） |
| 47 | `tts.stop` | 27 | 1 | 零 | speech |
| 48 | `tts.stop-bargein-400ms` | 41,41 | 2 | 零 | speech（D47） |
| 49 | `ui.one-click-restart` | 35 | 1 | 零 | panel |
| 50 | `ui.voice-cancel-availability` | 22 | 1 | 零 | panel（B1） |

合计：50 枚名／62 次实例／39 枚带副作用的行／今天真投递 **1** 枚／结构可达 **1** 枚／其余 **48** 枚在产码里派发不到。

### 1.5 AC#1 两发读数（未修码，现量）

- **负向（今天拦不住）**：尺 C＋C′ 同读 = 产码里 **0 枚** `statemachine.New` 传了 `Sink`（2 枚构造点全是只有 `Initial` 的单行字面量）。⇒ 票面 §2② 成立。
- **正向对照（证明尺命中得了真名）**：`internal/models/bridge_test.go:16-18` 就往 `Options` 里塞了
  `Sink: func(e statemachine.Effect) { effects = append(effects, e) }`，并在 `:49-55` 断言 `e.Name == "model.verify-sha256-signature"`，
  否则 `t.Fatalf("row #37 side effect not fired; effects %v", ...)`。
  第二枚正控 `internal/statemachine/table_test.go:294 TestEveryLegalRowFires`（`:298` 塞 Sink、逐名比对 `tc.effects`）。
  ⇒ 尺认得名、也认得"塞了 Sink 的形状"，**它今天在产码里读 0 不是尺瞎，是真没有**。
- ⚠ 关键形状（AC#4 要用）：`bridge_test.go:11 newWalkRig` **自己建自己的机器**，`internal/models/handoff_window_109_test.go:62 newWalkMachine()` 甚至**复制了产码那枚没有 Sink 的字面量**。
  ⇒ 测试绿与产品有收件人**是两回事**，这正是票面 §2② 说的"用例绿而产品看不见"那一族，本腿给它补了形状。

---

## 2. AC#2 三形代价表（本腿不选形，⛔ 只摆代价）

先给一把共用的尺：**新依赖边的现量底**。`grep -rln '"github.com/CarlosShao/wisp/internal/statemachine"'` 去测试去 `.scratch` = **16 枚文件**
（`cmd/balldebug`、`cmd/wisp`×5、`internal/agent/approval/replies.go`、`internal/ball`×5、`internal/models/bridge.go`、`internal/tools`×3），
而 `internal/statemachine/*.go`（非测试）里 `grep wisp/internal` = **0** ⇒ **这枚包今天零内部依赖**（只有 `fmt/log/slog/sync/time`，〔读码推〕**任何包级新依赖边都可能成环**，
`ball`/`tools`/`models`/`approval` 都在**反向**指它）。⚠ 落地前后要用 `go list -deps` 差集逐名对（票 246 AC#1 立的尺，本腿没跑）。

### 甲＝装配根注入一枚真 Sink（`cmd/wisp`，与票 246"门由装配根注入"同一条路）

| 项 | 读数 |
|---|---|
| 动哪几枚文件 | 必动：`cmd/wisp/models.go`（链① 的构造点 `:303`）。新代码建议落**新文件**（如 `cmd/wisp/effect_sink_273.go`，命名由编排者定）。要覆盖交互腿还得先动 `cmd/wisp/run.go`／`cmd/wisp/resident_ball_windows.go`——**那里今天没有机器实例**（§1.3 链③），所以甲在交互腿上是"**先造发件人，再造收件人**"两件事，不是一件事 |
| 新依赖边 | **0**（〔读码推〕`cmd/wisp` 已 import `statemachine`／`ball`／`panel`／`speech`，装配根往已存在的边上加投递是最便宜的一形）；落地前后仍要 `go list -deps` 差集逐名相同（246 尺） |
| 撞哪几枚既有钉 | **枚数＝1 枚用例（红名册只 1 名）＋2 枚用例（视写法）**：<br>① `cmd/wisp/config_receipt_255_test.go:179 TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` —— 该钉把 `cmd/wisp/config_readers_255.go` 里所有 `file.go:LINE [token]` 引用**逐条开文件验行号**：`cmd/wisp/models.go:163`（2 处引用）、`:183`。⚠ 只要**在 163 行以上插行**就红；插在 `handOffModel`（`:302` 以下）不红。同钉还引 `cmd/wisp/run.go:991`（3 处）、`:423/:424/:435/:1014`、`cmd/wisp/resident_ball_windows.go:276`（2 处）、`panel_resident_windows.go:207/:253` ⇒ **改 run.go／resident_ball 就撞上它**。<br>② `:230 TestTicket255RosterStillMatchesTheActualReadSites`（同文件另一枚）—— 若新 Sink 里**读任何配置字段**（`scanReadFiles` 扫的是 `os.ReadFile`／配置读点），名册必须增行，否则红。<br>③ 不受损（现量核对）：`cmd/wisp/leg_sink_gate_131_test.go`(1 名)／`leg_sink_nail_131_windows_test.go`(3 名)／`leg_dispatch_gate_133_test.go`(1 名) 盯的是 **`installLogSink` 日志收件人**与 `main.go` 派发可达性，与状态机 `Sink` **同名不同物**（`leg_sink_nail_131` 头部自陈是票 131 AC#2/#3 的日志腿）；甲只要**不删 `cmdModels` 那一跳、不动 `installLogSink` 那九行**就不碰它们。`internal/models/bridge_test.go`(3 名) 只要 `WireDownloading` 签名不动就不红 |
| 这条边以后谁看得见 | **装配根**＝`cmd/wisp` 是唯一一处"谁收 Effect"被写下来的地方，与 246 的裁定同形（票面 `246-…md`：**「我裁走乙（装配根注入）」、`A481`、09-30 18:40**，且"⛔ 不开任何新包级依赖边"）；grep `Sink:` 在 `cmd/wisp` 里会同时命中 `consoleSink`／`storeHealthSink` 两枚同名罐子（§1.2 尺 D），**可读性差一截**，建议落新文件并用不与 `Sink` 撞的名 |
| 额外代价 | ⛔ 不许写成"加一条日志就算收件人"（票面 §3 第二条）；若收件人要异步（`machine.go:31-33` 明令 Sink **不得阻塞、不得同步重入 Dispatch**），一枚真收件人大概率要**自己排队／起 goroutine** ⇒ 撞上禁止清单两条硬形状：裸 `go func(` 无 owner/recover、用墙钟差实现超时（`tools/d22scan`） |
| ⚠ **本腿补的一格硬约束（现量）** | `head -1 cmd/wisp/models.go`／`run.go`／`main.go` ＝ `package main`／文档注释，**无构建标记＝可移植文件**；而 `grep -rln 'internal/ball' cmd/wisp` 出来的**每一枚非测试文件都带 `//go:build windows`**（`resident_ball_windows.go`、`resident_approval_windows.go`、`resident_windows.go`）。<br>⇒ 甲**接在链①（`models` 腿）上够不到任何界面收件人**：球的 `SetState`／面板都在 windows 一侧，可移植文件里调它们会把 ubuntu 的 `go vet ./cmd/wisp` 编译面打断。甲在链①上唯一可移植的收件人是 **stdout／store／observe** 那一类 ⇒ **实质就是把票面 §3 明令禁止的"改成加一条日志"落地**。<br>⇒ 要真治"界面没反应"，甲**必须延伸进常驻腿**（那正是丙的活），或在 windows 一侧另装配；**"甲只修链①"这枚形状修不了票面 §1 说的那件产品现象**，本腿按原文把这话说清，不替编排者选 |

### 乙＝包内自带一枚默认 Sink

| 项 | 读数 |
|---|---|
| 动哪几枚文件 | `internal/statemachine/machine.go`（`New()` 的 `:64-68` 缺省支）＋ 包内新文件；`internal/statemachine/doc.go` 的自述必须同时改（现量原话：**"the default sink is an explicit no-op - execution belongs to later tickets"**、**"no capabilities of its own: it never touches audio, LLM, tools or UI"**） |
| 新依赖边 | **必有**（否则默认 Sink 什么都做不了）：`statemachine → ball／panel／speech` 任何一枚都**成环**（§2 表头：那 4 个包今天反向 import `statemachine`）。剩下不环的选项只有 `log/slog` ⇒ **乙的极限就是"把空罐子换成日志罐子"**，而票面 §3 明令这条不在射程 ⇒ 乙**结构性地只能干被禁的那件事**，或干一件没人看得见的 logging |
| 撞哪几枚钉 | **枚数＝0 枚会真红**——这是乙最危险的地方（现量＋读码双证）：唯一钉住这件事的用例是 `internal/statemachine/machine_test.go:186 TestSinkNoopDefault`，但它**通读全文只断言两件事**：`Dispatch(EvSummon)` 不报错、`m.State() == StateListening`。装一枚真默认 Sink（写日志、投递到包内队列）**两件事都不会破** ⇒ **乙不会让这枚钉红**。乙破的是**文档契约**（`doc.go` 两句、`machine.go:25` "the zero sink is an explicit no-op"）与 C12/票面 §3 的射程，⛔ 不是任何测试；另外 `internal/ball` 的 3 枚机器构造点（`hotkey_live_test.go:382`、`interaction_live_test.go:54/:134`）与 `handoff_window_109_test.go:62`、`bridge_test.go:86` 全是 `New` 不带 Sink ⇒ **一旦默认 Sink 有副作用，这些用例的观测面被悄悄改写**（〔读码推〕，本腿没跑） |
| 与"第二枚真相源"那条已裁规矩的关系 | **要说清：乙与它冲突。** 票 246 的裁定句（`246-…md` AC#1）与台账 `docs/reports/pending-and-issues.md:10001` 立的是**"一个进程只许一枚 `approval.Gate`，由装配根注入，⛔ 不开新包级依赖边"**；台账 `:10001` 那一行还写明常驻腿的门在 `resident_approval_windows.go:109`。把"谁收 Effect"下沉进 `statemachine` 包，就是**再造一枚"由包自己决定发给谁"的真相源**，与"装配根是唯一接缝"（票 246 引的票 197 段原话："装配根是唯一的接缝"）**同轴冲突**；再叠上"审批队头（C18）谁在等批准"那枚真相源，界面侧就会出现两枚各自决定"谁在等／谁收到"的地方 ⇒ **本腿据这两句原文判定乙与已裁规矩冲突，写下来请编排者裁，不替他裁** |
| 这条边以后谁看得见 | **藏在包里**：`grep -rn 'Sink:' cmd` 看不见它，只有读 `machine.go` 才知道"发给了谁" ⇒ 以后排查"界面没反应"要多绕一层；与票面 §1 抱怨的那句"读代码的人会以为界面没反应是 bug"正好同形，只是把空罐子换成隐式罐子 |

### 丙＝只在常驻腿接，`cmd/wisp` 那条（`models` 腿）不动

| 项 | 读数 |
|---|---|
| 动哪几枚文件 | 常驻腿今天**没有机器**（§1.3 链③）⇒ 丙至少四件事：① 在 `cmd/wisp/resident_ball_windows.go`（或同腿新文件）**建**一枚 `statemachine.New`；② 找到事件源把 `Event` 喂进去——今天球的输入是 `ball.Events{OnClickBall／OnSummonHotkey／OnMuteHotkey／…}`（`resident_ball_windows.go:278-289`）直连 `recordBallGesture`，**没有任何一处调 `Dispatch`**（现量：`.Dispatch(` 去测试只有 `cmd/balldebug:618`、`internal/models/bridge.go`×4、`machine.go:202` 自调，另两枚是**面板 ComposerDispatch 同名物** `cmd/wisp/panel_resident_windows.go:360 w.Dispatch(fn)`）；③ 收件人（`ball.SetState`/`SetBadge`/`SetBadgeText`/`SetTrayTip` 都有真形，`internal/ball/ball_windows.go:311/:316/:332/:961`）；④ 与 `resident_approval_windows.go:879/:956` 那两处**直接点名 `SetState`** 的现有控制权正面对撞（谁说了算要定案） |
| 新依赖边 | **0**（仍在 `cmd/wisp` 内装配），但把 `ball` 的控制权从"卡片级别直接 `SetState`"改成"Effect 驱动"＝**改常驻腿的状态所有权**，这条不是依赖边而是**语义边** |
| 撞哪几枚钉 | **枚数＝1 枚必红（行号漂移）＋2 枚按名可指＋2 枚 winlive 不可跑**：<br>① `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` —— 名册钉了 `cmd/wisp/resident_ball_windows.go:276 [Hotkeys:  cfg,]`（**2 处引用**），**在 276 行以上插行即红**（现量：`:276` 内容确为 `Hotkeys:  cfg,`）。<br>② `cmd/wisp/resident_approval_246_windows_test.go:93 TestAC246CardWithNoWindowFailsClosedThroughTheRealGate` —— 该用例 `:119`/`:122` 直接断 `stateForCardLevel("L1"/"L2") == StateConfirming/StateAwaitingApproval`，而 `stateForCardLevel` 正是丙要**从"直接点名"改成"Effect 驱动"**的那枚函数（定义 `cmd/wisp/resident_approval_windows.go:965`；现量生产调用点**只有 2 处**、同一文件 `:879 b.SetState(stateForCardLevel(p.Level))` 与 `:882 "orb_state", string(stateForCardLevel(p.Level))`（日志字段），另加 `logs/state-for-card-level.txt` 记全 6 枚命中）。<br>③ `cmd/wisp/approval_seam_201_test.go:78`（＋`:94`/`:206`）断 `rt.waitingStateName()`（生产定义在 `cmd/wisp/approval_always.go:190`）返回 `AwaitingApproval` —— 状态的**来源**被丙改走机器以后，这枚读法要不要跟着换边，是丙的第二次撞钉。<br>④ `cmd/wisp/resident_approval_live_246_windows_test.go:133`/`:235` 两枚 winlive 档（`//go:build windows && winlive`）本腿⛔没跑、没读内容，**只能登记成"仅本机可量，未批"**。<br>⑤ 不受损（现量核对）：`subagent_blocked_197_test.go`／`subagent_carrier_197_test.go` 只把 `statemachine` 当**状态词表＋`statemachine.Valid` 判据**用（不建机器、不传 Sink），丙动不到它们 |
| 这条边以后谁看得见 | 只在常驻腿文件里；`wisp models`（唯一今天真投递那条链）**继续投空罐** ⇒ 丙**不修**票面 §2 那两条现量读数里的任何一条，只把"看不见"从一条腿扩成两条腿各自成立；甲的链①在丙下**永不修** |
| 额外代价 | 票面 §4 AC#2 让丙"逐形给"，本腿补一句代价：丙与甲**不是互斥**，丙实质是"甲＋给常驻腿补一台机器"，**比甲贵**；票面说"只在常驻腿接、`cmd/wisp` 那条不动"——若"那条"指 `cmd/wisp/models.go`，那丙保住的只是调试／观感，放弃唯一今天真会投递的名（`model.verify-sha256-signature`） |

---

## 3. AC#4 定向突变建议＋恒真性自评

**建议的定向突变（三发，缺一发就是恒真的风险）**

1. **M1 摘掉注入（形＝把 `Sink:` 那一行删掉，回到只有 `Initial` 的字面量）**
   施加在 `cmd/wisp/models.go:303`（甲落地后的那一行）。
   **指名必须红的用例**：一条**经进程内入口 `cmdModels([]string{"ensure", id}, …)` 驱动**、观测"收件人真的收到了 `model.verify-sha256-signature`"的新用例（形状照 `cmd/wisp/leg_sink_nail_131_windows_test.go:352/:500` 那两处现成的 `cmdModels` 进程内驱动，⛔ 不许新建一条只读源码的用例）。
2. **M2 收件人变空罐（形＝`Sink` 留着但里面什么都不做）**
   同一行不删，只把闭包体掏空。**同一枚用例必须同样红**——这一发是 M1 的反形，专治"判据只看有没有 `Sink:` 字样"那种恒真句。
3. **M3 投出去但投错名（形＝把 `:162-164` 的循环改成投一枚不存在的名字）**
   施加在 `internal/statemachine/machine.go`，要求同一枚用例红且**红因是断名那句**（不是超时、不是 panic）。这发证明判据认的是"哪一枚名到了谁"，不是"有没有东西动过"。

**恒真性自评（⛔ 不许把票 270 AC#1 那个错再抄一遍）**

- 会恒真的写法本腿**明确不推荐**：
  (a) "产码里存在一枚带 `Sink` 的 `statemachine.New`" —— 纯语法，M2 直接全绿；
  (b) "Sink 不是空函数" —— 反射／指针比较，M2 也全绿；
  (c) **`TestSinkNoopDefault` 那一族**：现量已证它只断"不报错＋状态变了"，**换成任何非空默认 Sink 它都不红**（§2 乙栏），所以它**不能**当"有没有收件人"的判据，只能当"没崩"的判据；
  (d) 拿裸 `grep '"<effect name>"'` 数收件人 —— 本腿实测 2 枚假阳性（注释 `internal/agent/guard.go:46`、事件名 `internal/statemachine/events.go:60`），换形照样绿。
- 不恒真的判据要同时满足三条，本腿把它写成建议：
  **①驱动真实进程入口（`cmdModels`／常驻腿装配函数），不是测试自建 rig**（`bridge_test.go:11 newWalkRig` 那种自建机器对 M1 完全免疫，现量已证）；
  **②断言"某枚具名 Effect 到达了某个具名收件人的可观测状态"**（不是断日志行数，票面 §3 禁的正是②退化成日志）；
  **③M1 与 M2 两支都红**——只有反形也红，才叫"把 Sink 摘掉那一形必须红"，否则就是票 270 AC#1 同一枚病。
- ⚠ 本腿**没跑任何突变**（闸门 §边界），以上全是**形状建议＋读码推**，红名册枚数要由跑突变的腿现量；本腿只保证"哪一枚用例**不该**被指望"是有原始输出的。

---

## 4. AC#3 边界自查：要不要人工批准？

- **要不要动 `C17` 方法白名单？本腿读数：不需要。** 现量：50 枚名**没有任何一枚**是 `C17` 面板桥方法；`tools.execute`（名册 #46）与 `panel.*`（#35-38）是**副作用名**，不是桥方法名。⛔ 本腿没去 `C17` 白名单里加名，也不建议。
- **要不要新造 `C##`？本腿读数：不需要。** `Sink`/`Effect` 已在 `C12` 所辖的 D43 语义内（`doc.go` 与 `machine.go:23-29` 的自述），三形都是"接线"，不是"新契约面"。
- **但有三处**碰到冻结件，**这半条要人工批准、不在本票射程，本腿只报不裁**：
  1. **实现表比权威表多两枚编号**（`#41`/`#42`，出处 `internal/statemachine/doc.go` 自述"SPEC-08 §3's appended #41/#42"，实现行 `table.go:195/:204/:246`；`docs/PLAN.md` §D43 只有 1–40）。
     D43 转移表是冻结件（`AGENTS.md` §1.1）——**它是不是已被人工批准过增列，本腿查不到凭据**，交编排者对台账。
  2. **乙形**要改 `internal/statemachine/doc.go` 与 `machine.go:25` 的"零罐子即明示 no-op／包不碰 UI"两句自述 ⇒ 那是**契约文字**（`C12` 一侧），⚠ 人工批准。
  3. **`conversation.privacy-confirm-l2`（#16）与 `tools.execute`（#46）**一旦被接成"真做事的收件人"，分别撞上
     **L2 卡不许由面板侧来源**与**"宿主内部 artifacts 不许实现成受门控 Tool"**两条禁止清单（`AGENTS.md` §1.2 逐字）。
     ⇒ 这两枚**接线时**要人工点头，⛔ 本票射程内只做名册，不做决定。
- ⛔ 本腿没翻任何 AC 框、没改票面、没动台账、没新建工单。

---

## 5. 判不动（本腿停在这里，等编排者）

1. **丙形最终红几枚**——本腿已把**该指名的钉逐名指到用例**（§2 丙栏①–⑤：1 枚行号漂移必红＋2 枚按语义待裁＋2 枚 winlive 不可跑），
   但"改完 `stateForCardLevel`／`waitingStateName` 的来源以后到底哪几名进红名册"要跑那天现量；本腿⛔不能跑 `go test`（`272-r2` 在飞）。
   ⚠ 同一条尺也量到"不受损"的：`subagent_*_197_test.go` 两枚只拿 `statemachine` 当词表用，别让别的腿把它们误报成风险。
2. **"收件人的可观测状态"该断在哪一层**——`ball.SetState`/`SetBadge` 的真形在 `internal/ball/ball_windows.go`，带 `windows` 构建约束；
   ⛔ 本腿没开窗、没加 `-tags winlive`，无法判"不靠真窗口能不能断到"，AC#4 的用例形状定稿要这一条。
3. **`error.ack` 到底今天投不投得出去**——`defer machine.Close()`（`cmd/wisp/models.go:304`）与 `time.AfterFunc`（`machine.go:187`）的竞态本腿只能〔读码推〕；
   要量就得跑带计时的用例，越界。
4. **表枚数与权威表分叉（`#41`/`#42`）算不算已批准**——需要台账凭据；`docs/reports/pending-and-issues.md` 本腿只 grep 了"谁在等批准"与"第二枚"两组关键词（`logs/ledger-who-waits.txt` 1 行、`logs/ledger-truth-source.txt` 85 行），⛔ 没通读约 12,900 行母本，不敢据以判。
5. **甲形是否要求 `run.go`／常驻腿同时补机器**才算"票完成"——这是范围裁量，票面 AC#2 明令"⛔ 本票不选形"，本腿只把三形代价摆平，不替你选。
6. **"哪几枚名该由谁收"这张对应表本腿不裁**：名册 1.4 的 `→谁` 列是按名前缀的**读码推**归档（⛔ 不是任何仪器／文档里既有的映射）；
   仓里**不存在**一份"D43 副作用名 → 收件人"的分档表（现量＝尺 F：50 枚名逐枚 `grep -rF '"<name>"'`，去表／去测试／去 `.scratch` 后 48 枚 0 命中、2 枚假阳性）。
   这张表要人工或编排者定：它一旦沉进包里就成了 §4 第 2 条那处契约文字，一旦有名字被实现成 Tool 就撞上 §4 第 3 条 ⇒ **不在本票射程，本腿只把空位摆出来**。

---

## 6. 与编排者转述不符之处（具名报回）

| 转述／票面原文 | 本腿现量 | 判定 |
|---|---|---|
| "50 枚副作用名"（票面 §2③ 的普查腿转述） | **50 枚不同名**（尺 A＋尺 B 两把独立，`diff` 全等） | **相符**，顶不回来 |
| 票面 §2②"其余 **5** 枚命中都在 `_test.go`" | `statemachine.New(` 总 8 枚命中＝2 产码＋**6** 测试（票面自己下面列的正是这 6 处） | **枚数差一枚**（名单没错） |
| 票面 §2③"50 枚副作用名**全投进** no-op" | 出货链上今天真投出去 **1** 枚（`model.verify-sha256-signature`）；`error.ack` 结构可达但〔读码推〕赶不上；其余 **48** 枚产码里**派发不到那一条边**（`wisp run`／常驻腿**没有机器实例**） | **方向对、形状不对**："全投进空罐"其实是"1 枚真投进空罐＋48 枚根本没被投＋收件人 0/50" |
| 票面 §1"状态变了但没有任何东西收到过" | 常驻腿的球视觉**确实变了**（`resident_approval_windows.go:879 b.SetState(stateForCardLevel(...))`），但**不是状态机变的** | 比票面更严重一格：**发件人也不在那条腿上** |
| 票面 §2②"测试里有传 Sink 的形状" | 6 枚测试构造点里**只有 1 枚**（`internal/models/bridge_test.go:16-18`）传 Sink；`internal/ball` 3 枚与 `bridge_test.go:86`、`handoff_window_109_test.go:62` **都只填 `Initial`**（正控另有一枚在 `table_test.go:298`，写作 `New(Options{…})` 不带包名限，所以不在尺 C 的 8 枚里） | **相符但需收窄**："有传 Sink 的形状"＝2 处，不是普遍形状 |
| （无转述）三枚同名 `Sink` 罐子 | `consoleSink`＝`agent.Sink`、`storeHealthSink`＝`llm.HealthSink`、`leg_sink_nail_131`＝**日志**收件人 | **新报**：`grep Sink` 在 `cmd/wisp` 里的 4 枚真命中**全都不是**状态机那枚 |
