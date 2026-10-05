# 票 267 落地腿 `267-r2` — 证据件（种子迁移：`cmd/wisp` 带外 `confirm_timeout_sec` 抬进带内）

任务：票 267 的 ⓐ 值域（`[31, 3600]`，`internal/config/validate.go`，已进 HEAD）把 `cmd/wisp` 一批集成用例的种子打到带外 ⇒ 装配根拒载 ⇒ 16 枚 FAIL。本腿执行编排者 §7.1 裁的出口**甲**＝把带外种子抬进带内，并顺手拿票面 AC#2 那半格（C18"即将超时"提示的正向读数）。
具名解冻：台账 `A611` §1，撤销口令「267 撤」。

## §0 起手锚

- 取数时刻：`10:01:53+0800`（`date "+%H:%M:%S%z"` 现取）。
- 分支：`dev`；起手 HEAD 自取：`f9bd0eb7`（`git rev-parse --short HEAD`，⛔ 未抄派单）。
- 起手写面归属读数（尺＝`git status --porcelain -- cmd/wisp internal .scratch`，10:01:5x）：
  - `cmd/wisp` 与 `internal/**` 两棵子树**零条** ⇒ 本腿写面起手干净，band 之后没有别的腿正在动那五枚文件。
  - 同一读数里 `.scratch` 有 9 条 ` M`（`probes/152/my152.py`、`probes/161/r6/logs/flip-*.txt`）与一批 `??`（`ci-logs/run-*.log`、`commit-msg-*.txt`）⇒ 全是别的腿的临时件；本腿一枚未碰、一枚未删（临时件只建不删）。
- 并发腿：派单称另有两枚只读普查腿在提交它们自己的普查件。实测并发证据＝本腿种子件 `32e74479` 落地上方多出 `e3a8fe38`，尺＝10:16:2x `git show --name-only --oneline e3a8fe38` ⇒ **名册只有 `.scratch/wisp/probes/267/a3/census.md` 一枚**（`267-a3` 的普查件），与本腿写面零交集；本腿也没碰过它。
- 共享树下「起手 HEAD ≠ 交件 HEAD」是本腿正常的时序，不是本腿动了别人的面。

## §1 名册自取尺

⛔ 派单那张名单按「待验断言」处理，下面四把尺是本腿自己拉的（首拉 10:02–10:04，复拉与计数 10:15:03–10:15:04，全部 `date` 现取）。尺根＝`cmd internal tools docs .scratch scripts`。

- **R1（键名字面尺）** `grep -rn "confirm_timeout_sec\|ConfirmTimeoutSec" cmd internal tools docs .scratch scripts`：
  - 产码侧**消费者只有两枚装配根**：`cmd/wisp/run.go:616`、`cmd/wisp/resident_approval_windows.go:460`（同文件 `:382` 是日志 attr、`:308/:421` 与 `cmd/wisp/resident_windows.go:128` 是注释）。
  - `internal/config/**` 的命中全部属带门本身（`validate.go:103-150`、`schema.go:460`、`unwired.go:121`）与它的测试 ⇒ ⛔ 本腿一枚未动。
  - 直接写文件的字面种子只有两枚：`cmd/wisp/panel_pump_test.go:136`（种 2）与 `cmd/wisp/run_mode101_test.go:110`（种 1）。
  - `%d` 模板只有枚一处：`cmd/wisp/approval_reply_201_test.go:145`。
  - 常驻腿那一族 5 枚种子（`resident_approval_risk_256_windows_test.go:159/:202/:221/:286/:293`，值 45/90）**已在带内**，与本腿无关（票 256 的面）。
- **R2（顺 `%d` 追喂值链）** `grep -rn "l2Wait" cmd/wisp/*.go` ⇒ 10 行：那枚 `%d` 的唯一来源是 `int(h.l2Wait.Seconds())`（`approval_reply_201_test.go:147`），`h.l2Wait` 只能从 `newReplyHost(t, l2Wait time.Duration)`（`:103`）的入参进来；转发它的两枚 helper＝`newWidenRun`（`approval_always_201_test.go:81/:84`）、`newReloadRun223`（`config_reload_223_test.go:43/:45`）。⇒ **数喂值点＝数这三枚 helper 的调用点**，不是数 `confirm_timeout_sec` 这个词。
- **R3（调用点全名册）** `grep -rn "newReplyHost(t,\|newWidenRun(t,\|newReloadRun223(t," cmd/wisp/*.go`（10:15:03）⇒ **27 行**＝25 枚字面喂值点＋2 枚 helper 内转发（定义行读作 `func newReplyHost(t *testing.T, ...)`，不带 `t,`，天然不匹配）。
  - 其中**带外 10 枚**（20 s ×9 ＋ 2 s ×1），**带内 15 枚**（40 s ×10、90 s ×4、60 s ×1）⇒ 与 `267-a2` 的「带外 12／带内 15」口径一致。
  - 加 R1 那两枚直接写文件的字面种子 ⇒ **带外总数＝12 枚，住 5 枚文件**。
- **R4（其他根）** `for r in docs tools scripts internal cmd; do grep -rn "confirm_timeout_sec" "$r" | wc -l; done`（10:15:04）⇒ `docs=24 tools=0 scripts=0 internal=28 cmd=17`。`docs` 那 24 行本腿逐条读过（`docs/PLAN.md:2738`、`docs/specs/SPEC-03-config-secrets-envs.md:34` 两枚是**冻结件里的默认值表格文字**，其余是 `docs/evidence/s1/**` 与 `docs/reports/**` 的散文/台账引用）⇒ **零枚喂值种子**，⛔ 一字未动（`docs/PLAN.md`、`docs/specs/**` 本就在禁区内）。

**名册（12 枚，与派单逐名对照）**

| # | file:line（迁移前） | 用例／宿主 | 该发自己声明的观察对象 | before | after |
|---|---|---|---|---|---|
| 1 | `cmd/wisp/run_mode101_test.go:110` | `t101host.writeConfig`（`TestTicket101*` 五枚共用） | 档位（permission_mode）读写与重启存活；超时那枚对本发是死重（L1 走 `gate.go` 的 `g.window`，失败路径先响的是 `:544` 那枚 300 ms ctx） | `confirm_timeout_sec = 1`（同行 `l1_window_sec = 1`） | `confirm_timeout_sec = 40`（`l1_window_sec = 1` **一字未动**） |
| 2 | `cmd/wisp/panel_pump_test.go:136` | `nonDefaultConfig145`（`TestRunBooksWithASnapshotOfItsLiveQueue`） | 逐字 `:124` "so the case can watch an L2 card open and close" | `confirm_timeout_sec = 2` | `confirm_timeout_sec = 31` |
| 3 | `cmd/wisp/approval_reply_201_test.go:202` | `TestReplyListenerAllowsAnL2CardFromTheNativeSide` | 原生侧答允许⇒卡合、执行 | 20 s | 40 s |
| 4 | `cmd/wisp/approval_reply_201_test.go:272` | `TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel` | 答拒绝⇒理由传到模型 | 20 s | 40 s |
| 5 | `cmd/wisp/approval_reply_201_test.go:350` | `TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject` | 逐字 `:346` "A longer C18 deadline here"＝测路线不测钟 | 20 s | 40 s |
| 6 | `cmd/wisp/approval_reply_201_test.go:439` | `TestUnansweredL2CardTimesOutIntoRejectNeverExecution` | 逐字 `:466` "This is the timeout path, so the wait is the C18 clock" | 2 s | **31 s** |
| 7 | `cmd/wisp/approval_reply_201_test.go:504`（迁移后＝521） | `TestL1VetoNeedsAChannelTheHostReallyWired`／console posture | L1 窗口极性（到点执行）＋否决被拒 | 20 s | 40 s |
| 8 | `cmd/wisp/approval_reply_201_test.go:578`（迁移后＝595） | 同上／host declares and loads esc | 装载通道后的否决真停住写 | 20 s | 40 s |
| 9 | `cmd/wisp/approval_seam_201_test.go:50` | `TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole` | 导出接缝答卡（无 reply listener） | 20 s | 40 s |
| 10 | `cmd/wisp/approval_seam_201_test.go:138` | `TestNativeHostSeamRefusesAPanelSourcedAllow` | 面板来源的允许被拒 | 20 s | 40 s |
| 11 | `cmd/wisp/ticket224_assembly_test.go:81` | `TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun` | 会话动词落盘一行 | 20 s | 40 s |
| 12 | `cmd/wisp/ticket224_assembly_test.go:234` | `TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking` | 活授权让装配根不再发问（按卡片数判） | 20 s | 40 s |

**与派单名册的差集＝空（两个方向都零枚）**：12 枚逐一对上号。本腿另外具名三处口径差别，免得下一枚腿照抄吃亏：

1. 派单内部两个数不一致（§3 写「5 枚文件 12 处」、§6 写「上面那 10 枚带外喂值点」）——本腿这把尺给出的分解是**两者都成立**：`%d` 喂值调用点 10 枚＋直接写文件的字面量 2 枚＝12 枚。⛔ 不是差集，是口径。
2. **潜在的第 13 枚**（本腿自己抓的，未改）：`approval_reply_201_test.go:105-110` 那枚兜底 `if l2Wait < time.Second { l2Wait = time.Second }` 会造出一枚带外的 1 ⇒ 现量 25 枚喂值点里**零枚低于 2 s**（迁移后最小 31 s），所以今天不可达；那枚兜底属形状不属种子，⛔ 本腿未动，登记给编排者（若将来有人喂亚秒值，它会在加载层响亮拒载，不会静默）。
3. `cmd/wisp/resident_approval_risk_256_windows_test.go` 那 5 枚种子（45/90）本就在带内，票 256 的面，⛔ 本腿一枚未动——**它不在派单名册里是对的**，本腿复核确认。


## §2 逐枚改动（before 逐字取自 `git show 32e74479^:<file>`，after 逐字取自现树）

本腿只改种子字面量：**零枚断言、零枚 `t.Fatal`/`t.Errorf` 句子、零枚阈值、零枚 `t.Skip` 被改动**（自证尺见 §8）。

| # | file:line | before（逐字） | after（逐字） |
|---|---|---|---|
| 1 | `cmd/wisp/run_mode101_test.go:110` | `	riskLines := "[risk]\nl1_window_sec = 1\nconfirm_timeout_sec = 1\n"` | `	riskLines := "[risk]\nl1_window_sec = 1\nconfirm_timeout_sec = 40\n"` |
| 2 | `cmd/wisp/panel_pump_test.go:136` | `…"\n[risk]\nconfirm_timeout_sec = 2\npermission_mode = \"ask_high_risk\"\n"…` | 同形，只有 `confirm_timeout_sec = 2` 改成 `confirm_timeout_sec = 31` |
| 3–5, 7–8 | `cmd/wisp/approval_reply_201_test.go:202/:272/:350/:521/:595` | `h := newReplyHost(t, 20*time.Second)`（子测那两枚多一级缩进） | `h := newReplyHost(t, 40*time.Second)` |
| 6 | `cmd/wisp/approval_reply_201_test.go:439` | `	h := newReplyHost(t, 2*time.Second)` | `	h := newReplyHost(t, 31*time.Second)` |
| 9–10 | `cmd/wisp/approval_seam_201_test.go:50/:138` | `	h := newReplyHost(t, 20*time.Second)` | `	h := newReplyHost(t, 40*time.Second)` |
| 11–12 | `cmd/wisp/ticket224_assembly_test.go:81/:234` | `	h := newReplyHost(t, 20*time.Second)` | `	h := newReplyHost(t, 40*time.Second)` |

**取值的理由（逐枚，不是随手拍）**

- **31 只给了两枚**（#2 `panel_pump_test.go:136`、#6 `approval_reply_201_test.go:439`）：它们俩的观察对象**就是超时本身**——#6 逐字 `:466` "This is the timeout path, so the wait is the C18 clock"，#2 逐字 `:124` "so the case can watch an L2 card open and close"。31 是「让 C18 提示仍能武装的最小合法值」（`31 - 30 = 1 s > 0`，`gate.go:528`），也正是带下界（`internal/config/validate.go:128 confirmTimeoutSecMin = 31`）。**再小就越界拒载，再大就把两发墙钟继续往上抬**（各 ＋29 s，本腿实测见 §5）。
- **40 给了其余十枚**：它们答完卡就走（#3/#4/#9/#10/#11/#12 由 reply listener 或导出接缝直接答；#7/#8 的观察对象是 **L1 窗口**，`confirm_timeout_sec` 对那两发是死重；#5 的注释逐字要求"测路线不测钟"）。40 s 在这棵树上有 **10 枚现成先例**做同一件事（`approval_always_201_test.go:122/:164`、`always_write_no_clobber_226_test.go:40`、`config_reload_223_test.go:247/:472/:576`、`config_receipt_255_test.go:400/:468/:507/:564`；另 90 s×4、60 s×1，共 15 枚带内喂值点，尺＝§1 的 R3）⇒ 照先例取最省事，也让下一枚腿读得出"这是同类种子不是特调"。
- **为什么没有把 #1 抬到 31**：它不观察超时，40 与同族一致。同行 `l1_window_sec = 1` **一字未动**（那枚今天被 `gate.go:149-150` 钳进 `[MinL1Window, MaxL1Window]`，注释见 `gate.go:27`；票面写着它属另一枚形状，派单也点名禁区）。
- ★**行号漂移具名**：#7/#8 迁移前在 `:504/:578`，迁移后在 `:521/:595`（本腿 §4 那枚新增断言在它们上方插了 17 行）。派单给的 `:504/:578` 是**迁移前**的号，照抄会扑空；尺＝被指的那串字符（`grep -n "newReplyHost(t, [0-9]"`），不是行号。

## §3 配套件：`panel_pump_test.go` 的那枚 ctx

- file:line：`cmd/wisp/panel_pump_test.go:112`（`executeOn145` 内，该 helper 是那发驱动一张 L2 卡的唯一手）。
- before：`	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)`
- after：`	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)`
- 为什么它是**配套件而不是顺手优化**：种子从 2 s 抬到 31 s 后，**30 s 的 ctx 会先响** ⇒ `gate.go` 走的是 `case <-ctx.Done()` 那一支（`gate.go:579` 逐字 `任务上下文已结束，审批请求已作废并按拒绝处理`），不再是 `case <-deadline` 那支的 `审批超时（31 秒未确认），C18 一律判拒绝`（该句由 `internal/agent/approval/queue.go:479` 现算）⇒ 这一发观察的对象从「C18 超时」搬家成「上下文已结束」＝**盘上看着像绿了，量的却不是同一件事**。抬到 45 s（＝31 s 种子 ＋ 14 s 余量）后先响的仍是卡的 deadline。
- 凭据（现量，非推理）：M1 那发的 stdout 里逐字连着两行（`.scratch/wisp/probes/267/r2/mutation-m1.log:26-27`）：
  - `        [warning] 审批将在 30 秒后自动拒绝，请尽快确认`
  - `        [dismissed] 审批超时（31 秒未确认），C18 一律判拒绝，已自动拒绝`
  ⇒ 31 s 的钟走的是超时支、不是 ctx 支。⚠ 这两行量的是 #6 那一发（同一条 `gate.go` 分支、同一枚常量）；#2 那一发的 45 s 余量由整包里 `TestRunBooksWithASnapshotOfItsLiveQueue` 从红转绿坐实（名册见 §6），**不是由上面那两行推的**。
- ⛔ 除此之外本腿**没有动任何别枚** `context.WithTimeout`：尺＝`grep -rn "context.WithTimeout" cmd/wisp/*.go`（迁移后 27 行），与 `32e74479^` 同尺名册逐行比对，唯一差异＝`panel_pump_test.go:112` 的那枚 30 到 45（比对读数在 §8）。

## §4 ★C18"即将超时"提示正向读数的落点与突变自证

**结论＝读到了（不是"够不着"）**，落点＝这一发装配根的 **stdout（控制台）**，凭据逐字在下面；两发突变都真红，写面只落在本腿自己的新断言与种子上。

### 4.1 现读的链（⛔ 只读，`internal/agent/approval/**` 一字未改）

1. `internal/agent/approval/gate.go:526-529`——`deadline := g.clock.After(g.q.Timeout())`，然后 `var warn <-chan time.Time` 与逐字 `if lead := g.q.Timeout() - g.q.WarningLead(); lead > 0 { warn = g.clock.After(lead) }` ⇒ **lead ≤ 0 时 `warn` 是 nil**，select 里那一支永远不响＝票 267 的"静默撤保护"。
2. `internal/agent/approval/gate.go:557-568`——`case <-warn:` 造 `Event{Kind: EventWarning, CorrelationID: corr, Remaining: g.q.WarningLead(), Text: fmt.Sprintf("审批将在 %d 秒后自动拒绝，请尽快确认", int(g.q.WarningLead().Seconds()))}`（那句 Sprintf 在 `:563`），交给 **`g.ui.Update(ctx, w)`**；返回 err 才走 `g.logf("approval: warning delivery failed corr=%s: %v", ...)`，然后 `warn = nil // fire once`。⇒ 消费者＝**`approval.UI` 接口的实现方**，不在 approval 包内。
3. `cmd/wisp/run.go:600-611`——装配根选 UI：`if s.cards == nil || s.ui == nil` 那一支之外，`:607` 逐字 `rt.ui = nil // no console surface in a process whose cards go to another host`；本发（`newReplyHost` → `runTextTask`，没有注入 gate/cards）走 `:611` `rt.ui = &consoleApprovalUI{out: s.stdout, run: rt, live: rt.liveCards}`。
4. `cmd/wisp/run.go:1393-1394`——`func (u *consoleApprovalUI) Update(...)` 的第一句逐字 `fmt.Fprintf(u.out, "[%s] %s\n", e.Kind, e.Text)`；`EventWarning` 的字面值＝`internal/agent/approval/ui.go:65` `EventWarning EventKind = "warning"   // C18 last-30s prominent warning`。⇒ **stdout 上那一行＝`[warning] 审批将在 30 秒后自动拒绝，请尽快确认`**。
5. ★顺带量到一枚**别的可读面的边界**（不是本腿断言，是给编排者的读数）：`cmd/wisp/run.go` 的 Update 只在 `EventDismissed`/`EventStarted` 上 `publish()`（`run.go:1410-1413` 注释逐字 "Only the two kinds that move the queue publish a packet. A tick or a / warning changes a countdown the four-key snapshot has no field for"）⇒ **`warning` 不上面板快照**，所以这一发的可读面只有控制台；常驻腿那枚 `ballCardUI` 在 `cmd/wisp/resident_approval_windows.go:878` 另有 `case approval.EventWarning:`，但那是 GUI 宿主、不是本票这一发。

### 4.2 新增的正向断言（只新增，⛔ 既有断言一字未动）

落点＝`cmd/wisp/approval_reply_201_test.go:490-506`（`TestUnansweredL2CardTimesOutIntoRejectNeverExecution` 内，紧跟既有那条 `ANSWER-EXPIRED ... decision=timeout->reject` 审计循环之后）。逐字：

```go
	if sentences := h.out.String(); !strings.Contains(sentences,
		"[warning] 审批将在 30 秒后自动拒绝，请尽快确认") {
		t.Errorf("the C18 pre-timeout warning never reached the console: lead = 31s - 30s "+
			"= 1s > 0 must arm it (gate.go:528), stdout:\n%s", sentences)
	}
```

（上方另有 9 行注释逐字记了链与"迁移前这一支是死的"。）

**正向读数（现跑）**：10:09:05→10:09:44，`go test -count=1 -timeout 300s -run 'TestUnansweredL2CardTimesOutIntoRejectNeverExecution' -v ./cmd/wisp/`（PATH 带 `third_party/sherpa-onnx`）⇒ `--- PASS: TestUnansweredL2CardTimesOutIntoRejectNeverExecution (32.85s)`，台件 `.scratch/wisp/probes/267/r2/timeout-case-after.log`。

### 4.3 突变自证（⛔ 只落在 `_test.go`，未为制造红改任何产码）

- **M1＝把提示的句子算错**（needle 从 `30 秒` 改成 `31 秒`，即断言要求的不再是 `WarningLead` 现算的那句）⇒ **红**，且红句把真实 stdout 整段倒出来，逐字（`.scratch/wisp/probes/267/r2/mutation-m1.log:13,26,27`，10:11:44→10:12:22）：
  - `--- FAIL: TestUnansweredL2CardTimesOutIntoRejectNeverExecution (32.54s)`
  - `    approval_reply_201_test.go:504: the C18 pre-timeout warning never reached the console: lead = 31s - 30s = 1s > 0 must arm it (gate.go:528), stdout:`
  - `        [warning] 审批将在 30 秒后自动拒绝，请尽快确认`
  - `        [dismissed] 审批超时（31 秒未确认），C18 一律判拒绝，已自动拒绝`
  ⇒ 判语：这一枚断言**不是恒真**——它比的正是生产现算的那句；句子或前缀一变它就红，红时把可读面原文摊开。
- **M2＝把提前量打到 ≤ 0 的那一支**（同发种子从 31 s 改成 30 s，即票 267 立案的形状）⇒ **红，而且不是静默红**（`.scratch/wisp/probes/267/r2/mutation-m2.log:5,8,9`，10:12:53→10:13:01，1.41 s 就响）：
  - `[audit] perm: MODE-READ-FAILED ... err=config: config.toml: risk.confirm_timeout_sec 30 out of range [31, 3600] ... detail="档位读不到：本进程不缓存任何上一次的宽松值，决策链不会被装配（退出码 2）"`
  - `wisp run: 配置未就绪（Unconfigured）：config: config.toml: risk.confirm_timeout_sec 30 out of range [31, 3600]`
  ⇒ 判语：AC#2 要的"要么响、要么把提示保住"在这一发是**两样都拿到**——带内（31）时提示真到达（4.2），带外（30）时加载层响亮拒载，"静默把保护撤了"在**这一腿的可及范围内已不可表达**。
- **没跑的那一发，具名**：派单建议的第三形"把那句发出去的路掐掉"（不发 `EventWarning`／把 `warn` 留 nil）**在本腿写面内做不到**——那条路在 `gate.go:557-568` 与 `run.go:1394`，属本票禁区＋C18 冻结面；⛔ 本腿没有为制造红去动产码。要补这一发的唯一形状＝在 `internal/agent/approval` 里给 `UI` 造假实现，那是 **`267-r1` 已登记的 gate 级正控那一格**（`.scratch/wisp/probes/267/r1/evidence.md` §7.2，编排者 §4 那句"由落地种子那一发免费补上"本腿已在 4.2 补上端到端这一半），归口＝**编排者裁是否再派一枚 approval 层的腿**；M1＋M2 是本腿写面内能做到的最强两支，本腿不替自己圆成"突变已齐"。
- 还原定式（⛔ 无 checkout/restore/stash/reset）：`git cat-file blob HEAD:cmd/wisp/approval_reply_201_test.go > cmd/wisp/approval_reply_201_test.go`，两次还原后同发尺 `git status --porcelain -- cmd/wisp`＝**零行**（10:12:22 与 10:13:01 各一次）。突变原件日志（含被改坏的那一发）**只建不删**，全在 `.scratch/wisp/probes/267/r2/`。

## §5 门禁五把尺读数（全部本腿现跑，时刻＝`date "+%H:%M:%S%z"` 同发取）

| 尺 | 命令 | 取数时刻 | 读数 |
|---|---|---|---|
| G1 整包 | `export PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:$PATH"` 后 `bash scripts/portable-tests.sh --scope=cli`（输出重定向，⛔ 未用 `tail` 当证据） | 10:13:30 起 → 10:21:08 止（458 s，rc=1） | **`=== RUN=323  --- PASS=226  --- FAIL=1  --- SKIP=0`**，逐字行＝`portable-tests.sh: four numbers (all from -v output): ...`（台件 `cmdwisp-after.log:1671`）与 `runtests.sh: go test exited 1 - packages=[./cmd/wisp/ ...] top-level: PASS=226 FAIL=1 SKIP=0, === RUN=323, '[no tests to run]'=0`（`:1670`）；台件 `.scratch/wisp/probes/267/r2/cmdwisp-after.log`（1,673 行）＋ `cmdwisp-after.stamp` |
| G2 | `go test -count=1 ./internal/config/` | 10:22:38 | `ok github.com/CarlosShao/wisp/internal/config 1.212s`，rc=0 ⇒ **band 与那包测试没被本腿碰坏** |
| G3 | `"$(go env GOPATH)/bin/gofumpt" -l cmd/wisp/` | 10:22:40（迁移前 10:09:57 同尺一次，读数相同） | 名册只有 `cmd\wisp\models.go` 一枚 ⇒ **本腿五枚文件全部干净**；那一枚是**预存**：`file cmd/wisp/models.go`＝`... with CRLF line terminators`、`gofmt -l` 同判、`git diff --name-only` 里 models.go **零命中**（＝与 HEAD 同字节），且本腿从未打开过它。⚠ 结论：这条不是本腿造成的，但也**不是零命中**，本腿不替它圆 |
| G4 | `sh scripts/d22scan.sh` | 10:22:50 → 10:23:20，rc=0 | `d22scan: clean - no D22 ban violations`；scope 逐行＝bans #1-5 internal/=228、cmd/=38、ban #6 frontend/=85、ban #7 internal/tools/=23、ban #8 design/=39／frontend/=85／internal/=508／cmd/=101（`_test.go` 与注释都在射程内）⇒ **本腿新增注释与那句中文提示没有踩任何禁令**（ban #9 phantom citation 定义在 `tools/d22scan/main.go:50`，零命中；台件 `d22scan.log`） |
| G5 | `sh scripts/check-path-length-budget.sh --with-self-test` | 10:23:28 → 10:23:33，rc=0 | `check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`；positive control 三控全 ok；分母＝5712 tracked paths，over-budget=57／covered=57／**not in roster=0**（台件 `path-length.log`） |
| G6 | `go vet ./cmd/wisp/` | 10:22:41，rc=0 | 零输出 |

**墙钟（⛔ 本腿没有为快而压任何种子／压短任何超时／加任何 `t.Skip`）**

- 包体 elapsed：基线（编排者 09:46:00→09:50:2x）＝246 s；本发＝458 s。
- 顶层用例墙钟**合计**（尺＝两发日志各取 `^--- (PASS|FAIL|SKIP):` 的 `(Ns)` 求和，台件 `dur-before.txt`／`dur-after.txt`）：**195.5 s → 438.2 s＝＋242.7 s**，枚数两发都是 227。
- ⚠ 派单预估是 **＋58 s**，实测 **＋242.7 s＝约 4.2 倍**，差额本腿具名到用例（不是"应该无关"）：
  - `TestTicket101SessionGrantDoesNotCrossRestart` 1.33 → **81.70 s**（＋80.37）
  - `TestTicket101ManualSwitchSurvivesRestart` 1.39 → **42.24 s**（＋40.85）
  - `TestUnansweredL2CardTimesOutIntoRejectNeverExecution` 1.10 → **32.44 s**（＋31.34，＝#6 那发的 31 s 钟）
  - `TestRunBooksWithASnapshotOfItsLiveQueue` 1.39 → **32.60 s**（＋31.21，＝#2 那发的 31 s 钟）
  - `TestTicket101UntouchedConfigRestartsAtDefault` 1.31 → **7.65 s**（＋6.34）
  - ⇒ **种子可归因＝＋190.1 s**，其余约 ＋52 s 分散在 `TestCleanCheckoutBuilds_AC11`（＋7.22，那是一枚真 `go build`）、`TestAC228ResidentLegReportsAndBooksItsBall`（＋2.53）、`TestAC246DevLegIgnoresTheTestTaskInjection`（＋2.19）等**与本腿无关的机器负载抖动**（这些用例里零枚 `confirm_timeout_sec`，尺＝§1 的 R1）。
  - ★那一族 101 的 ＋128 s 推翻了一条本腿照抄的前提（"其余 9 枚答完卡就走＝零增量"）：`run_mode101_test.go` 那几枚**真的在等 `confirm_timeout`**（重启若干次、每次等一张没人答的卡到点），种 1 时是 1 s、种 40 后是 40 s。判语与可选的更形（#1 取带下界 31 ⇒ 每发省 9 s，仍带内、仍不动 `l1_window_sec`）写在 §7，归编排者裁，⛔ 本腿不自裁。

## §6 红名作差（判据＝只减不增）

**基线名册（编排者 09:46 实测 16 枚）逐名判决：16 枚全部转绿，零枚残留。** 尺＝`grep -c -- "--- PASS: <name>" cmdwisp-after.log`（读数逐枚＝1，含子测者 3／4，见下表），同尺在 `cmdwisp-after.log` 里对 16 枚**零枚**命中 `--- FAIL:`。

| 基线 FAIL 名（逐字） | 基线 | 本发 |
|---|---|---|
| `TestReplyListenerAllowsAnL2CardFromTheNativeSide` | FAIL 1.14 s | PASS 2.03 s |
| `TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel` | FAIL 1.10 s | PASS 1.58 s |
| `TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject` | FAIL 1.09 s | PASS 1.62 s |
| `TestUnansweredL2CardTimesOutIntoRejectNeverExecution` | FAIL 1.10 s | **PASS 32.44 s**（§4 正向读数就在这一发） |
| `TestL1VetoNeedsAChannelTheHostReallyWired`（含两枚子测） | FAIL 2.21 s | PASS 4.54 s（子测 2 枚同绿） |
| `TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole` | FAIL 1.22 s | PASS 1.17 s |
| `TestNativeHostSeamRefusesAPanelSourcedAllow` | FAIL 1.11 s | PASS 1.19 s |
| `TestRunBooksWithASnapshotOfItsLiveQueue` | FAIL 1.39 s | PASS 32.60 s |
| `TestTicket101ManualSwitchSurvivesRestart` | FAIL 1.39 s | PASS 42.24 s |
| `TestTicket101UntouchedConfigRestartsAtDefault` | FAIL 1.31 s | PASS 7.65 s |
| `TestTicket101SessionGrantDoesNotCrossRestart` | FAIL 1.33 s | PASS 81.70 s |
| `TestTicket101ModeSwitchUsesTheRealL2Gate` | FAIL 1.29 s | PASS 1.80 s |
| `TestTicket101UnreadableModeFailsLoudlyAndStrict` | FAIL 1.27 s | PASS 1.37 s（＋三枚子测，见下 RUN 差额） |
| `TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun` | FAIL 1.09 s | PASS 1.75 s |
| `TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking` | FAIL 1.08 s | PASS 3.96 s |
| `TestTicket224ProductionSessionDoesNotSurviveRestart` | FAIL 1.15 s | PASS 3.34 s |

**账算得平**：基线 `211 PASS + 16 FAIL = 227`；本发 `226 PASS + 1 FAIL = 227` ⇒ 16 枚全转绿、另有 1 枚由绿转红（下面那枚），`211 + 16 − 1 = 226` 逐名对上，零枚含糊。

**★新增红＝1 枚，不在基线 16 枚名册里，且不在本腿射程：〔未归因，归编排者〕**

- 名字：`TestAC14GoSideEvalPushReachesThePage`（`--- FAIL: ... (1.08s)`，`cmdwisp-after.log:799`）。
- 红句逐字（`cmdwisp-after.log:795-796`，两句都来自 `panel_resident_windows_test.go:867/869`）：
  - `    panel_resident_windows_test.go:867: AC#14 nail 2 (Eval push hop), page's own words: title=""`
  - `    panel_resident_windows_test.go:869: Go's Eval push did not reach the document: the page reports its title as "", want "PUSHED-33R5-OK". This is the push dimension, separate from the awaited-reply dimension in TestAC14AwaitedBindingReplyReachesThePage - one arriving says nothing about the other`
- 本腿能给的排除性读数（⛔ 不是"应该无关"，是量过的）：
  1. `grep -n "confirm_timeout\|\[risk\]\|newReplyHost" cmd/wisp/panel_resident_windows_test.go` ⇒ **零命中**，那一枚用例根本不写 `[risk]`，与本腿的 12 枚种子无关；
  2. 它自己的前一无关邻例同发绿：`TestAC14AwaitedBindingReplyReachesThePage` PASS 1.04 s（`cmdwisp-after.log:791`）；
  3. **隔离复跑全绿**：10:23:53 `go test -count=1 -run 'TestAC14GoSideEvalPushReachesThePage' ./cmd/wisp/` ⇒ `ok 1.308s`；10:24:11 再 `-count=3` ⇒ `ok 2.304s`（三发连绿，台件 `ac14-rerun1.log`）；
  4. 那一发在包内并发时段真建 WebView2 窗，本仓已把它记为"真窗族"（`docs/evidence/s1/33-panel-host-c27-v2.md:183` 逐字把这枚列进**11 枚在这台机上真建窗**的名册，并写着 "a machine that cannot create one is red, not skipped"）；台账另有两笔它在 CI 侧翻色的记录（`docs/reports/pending-and-issues.md:10727`、`:11899` 那笔"test-windows 20 枚→19 枚，消失的那枚＝`TestAC14GoSideEvalPushReachesThePage`"）。
  ⇒ **本腿判：这一枚不是种子迁移的红，也不是本腿能修的**；但它**确实出现在本腿那一发**，所以按判据具名标〔未归因，归编排者〕，把上面四条读数一并交出，等编排者用同 HEAD 复跑裁决（⛔ 本腿不许自行宣布"无关"就划掉）。
- `--- SKIP`＝**零枚**（四数之一，尺＝`grep -c "^--- SKIP"`＝0）⇒ 没有任何用例被"跳过"混成绿。
- **已知翻色核过**：`TestPanelHostLatencyPercentilesAC2` 在基线台件（`:692`）与本发（`:707`）都是 `--- PASS ... (0.00s)` ⇒ 本对之间**没有发生**编排者说的那枚 SKIP→PASS 翻色。
- **`=== RUN` 320 → 323 的差额具名**（尺＝两发 `^=== RUN` 名册 `sort` 后 `comm`，台件 `run-before.txt`／`run-after.txt`）：新增恰好三枚＝`TestTicket101UnreadableModeFailsLoudlyAndStrict/存储损坏`、`/版本不认识`、`/权限读不到`；消失＝**零枚**。机制＝基线那一发该用例在装配根就被 band 拒载、父体 1.27 s 即 FAIL，三枚子测从未 `=== RUN`；本发父体过了加载，三枚子测才登记 ⇒ **这一枚 RUN 差额是本腿修复的直接结果，不是仪器坏了**。

## §7 判不动／量不到（具名归口）

1. **★`run_mode101` 那族的 ＋128 s 值不值**：本腿按派单把 #1 抬到 40，实测那两枚重启族各等 40 s／80 s（§5）。更形＝#1 取带下界 **31**（每发省 9 s、仍带内、`l1_window_sec = 1` 不动）；⛔ 本腿不自裁，因为(a)改它要重跑整包才配得上"读数同发"的规矩，(b)"取 40 与同族一致"这条一致性判据属编排者 §7.1 的裁语。**归口＝编排者**（若裁"降到 31"，本腿可一发补上）。
2. **包体墙钟余量**：`tools/d22scan/runtests.sh` 与 `scripts/portable-tests.sh` 都**不传 `-timeout`**（尺＝`grep -n "timeout"` 两文件零命中）⇒ 走 Go 的默认 10 m；本发包体 elapsed 458 s（用例合计 438.2 s），距默认闸只剩 **约 142 s**，慢机上（CI 那台真窗族更慢，见 `:11886` 那笔 175.573 s）有撞闸风险。**归口＝编排者**（本腿⛔不许为省时间改任何种子，也不许给脚本加 `-timeout`——`scripts/**` 在禁区内）。
3. **★`panel_pump_test.go:124` 那句注释现在与字面量不一致**：逐字仍是 `//   - confirm_timeout_sec = 2, so the case can watch an L2 card open and close`，而 `:136` 现在写的是 `= 31`。⛔ 本腿没改它——写面是闭集（五枚 `_test.go` 的种子字面量＋那一枚 ctx 时长＋那一枚新断言），注释不在册。**这是盘上的过期陈述，归口＝编排者**（要么把它算进解冻的射程、要么另派一枚同批释文；后半句"so the case can watch an L2 card open and close"在 31 s 下**仍然为真**，说谎的只有那个数字 2）。同一形状在 `docs/evidence/s1/35-panel-snapshot-pump-r1.md:265/:454` 也以 `panel_pump_test.go:133`＋值 2 出现过（⛔ `docs/**` 禁区，本腿一字未动）。
4. **M1/M2 之外的那一发突变做不到的形状**：见 §4.3 末条——"把 `EventWarning` 那条发送路径掐掉"要动 `internal/agent/approval/**`（本票禁区＋C18 冻结面）。**归口＝编排者裁是否另派一枚 approval 层的腿**；`267-r1` §7.2 那半格的端到端这一半已由本腿 §4.2 补上。
5. **`WarningLead` 生产里恒零值**（`267-a1` §1 的账）：本腿正向读数依赖"没有生产者传 lead ⇒ `queue.go:88` 兜成 30 s"这一条。⛔ 本腿没有、也不该验证"若将来有生产者显式传一枚更大的 lead，31 s 就不够了"——那是票 255 的账（编排者 §残余③）。**归口＝票 255／台账**。
6. **`Q-77`（配置值 vs C18 写死 300 s 谁优先）**：⛔ 本腿一枚字都没碰，§4 那句提示里的"30 秒"是 `DefaultApprovalWarning` 现算的，不是 300 s 那一枚；主一句话之前，本腿不裁。
7. **量不到的那一格**：常驻腿（GUI）`ballCardUI` 收到 `EventWarning` 之后**页面上到底显出什么**——尺子够不着（`cmd/wisp/resident_approval_windows.go:878` 只读到 `case approval.EventWarning:` 那一层，再往下是 `frontend/**`／`design/**`，⛔ 本腿不许读，也不在此转述）。**归口＝票 268／界面会话**（票 268 的 provenance 那一支本来就在常驻腿）。

## §8 污染面与提交名册 ＋ 票面 AC 框零改动自证

**commit 名册（逐笔 `git show --name-only`，本腿共五笔，只 commit、⛔ 从未 push）**

| 笔 | 名册（该笔唯一改动面） |
|---|---|
| `00f0ef97` | `.scratch/wisp/probes/267/r2/evidence.md`（骨架，1 file changed, 46 insertions） |
| `32e74479` | `cmd/wisp/approval_reply_201_test.go`／`cmd/wisp/approval_seam_201_test.go`／`cmd/wisp/panel_pump_test.go`／`cmd/wisp/run_mode101_test.go`／`cmd/wisp/ticket224_assembly_test.go`＝**5 files changed, 30 insertions(+), 13 deletions(-)** |
| `5fa0d28c` | `.scratch/wisp/probes/267/r2/evidence.md`（§0–§7 填实，191 insertions／23 deletions） |
| 收尾笔 `aac52ab2` | `.scratch/wisp/probes/267/r2/**` 共 **19 枚**（evidence.md §8–§10 ＋全部台件：`cmdwisp-after.log`／`cmdwisp-after.stamp`／`mutation-m1.log`／`mutation-m2.log`／`timeout-case-after.log`／`ac14-rerun1.log`／`d22scan.log`／`path-length.log`／`dur-before.txt`／`dur-after.txt`／`run-before.txt`／`run-after.txt`／`pristine/models.go.pristine`／四枚 `commit-msg-*.txt`），⛔ 无一枚在闭集之外 |
| 第 5 笔（本行所在，交件后补的正名笔） | 只有 `.scratch/wisp/probes/267/r2/evidence.md` 一枚：把"共四笔"改成"共五笔"、把这张表补全、并给 §9 加第 12 条（数错自己的 commit 枚数）——⛔ 不再改任何产码或种子 |

- **写面并集（尺＝三笔 `git show --name-only` 去重，10:2x 现跑）**＝6 枚路径：五枚 `cmd/wisp/*_test.go` ＋本证据件 ⇒ 与派单给的闭集**逐枚相同，零越界**。
- **禁区尺**（同一把尺对 12 个禁区形态取反）：`internal/config`／`internal/agent/approval`／`internal/risk`／`thresholds.go`／`golden`／`allowlist.txt`／`docs/PLAN.md`／`docs/specs`／`.github`／`scripts/`／`frontend/`／`design/`／`.scratch/wisp/issues` 命中＝**0 行**。
- **票面 AC 框零改动自证**：`git status --porcelain -- .scratch/wisp/issues`＝**零行**；票 267 现态＝`AC#0`/`AC#1`/`AC#2`/`AC#3` 已是 `[x]`（编排者 09:5x 翻的，见票面「编排者收件」节），**`AC#4` 仍是 `[ ]`**（票面 `:31`）⇒ 本腿一枚未翻、未改字形。
- **`context.WithTimeout` 全株比对**（兑现 §3 那句"⛔ 没动别枚"）：尺＝`grep -n "context.WithTimeout"` 对 `32e74479^` 与现树的 `cmd/wisp/*.go` 逐行 diff ⇒ **唯一差异＝`panel_pump_test.go:112` 的 30 到 45**；`approval_reply_201_test.go` 那枚 5 s（`:684`，现 `:701`）计数前后都＝1、句子未变。
- **断言/阈值零改动自证**：`git show --stat 32e74479` 的 13 枚删除＝12 枚种子行＋1 枚 ctx 行，30 枚插入＝同 13 行的新版 ＋ 17 行新增块（§4.2 那段，含 9 行注释）⇒ **没有任何既有断言被改写**（本腿没动过一句 `t.Fatal`/`t.Errorf`，新增的是**另一枚** `t.Errorf`）。
- **污染面（本腿自己造过一枚，已消解并具名）**：为判 `gofumpt -l` 那枚 `models.go` 是不是本腿造的，本腿用 `git cat-file blob HEAD:cmd/wisp/models.go > .scratch/wisp/probes/267/r2/pristine/models.go` 落了一枚副本 ⇒ **它落在模块内**：`go list -e ./.scratch/wisp/probes/267/r2/pristine/` 返回 `github.com/CarlosShao/wisp/.scratch/...pristine`，`go vet` 该目录报 `vet.exe: ...models.go:110:18: undefined: resolveDataDir` ⇒ 这会让 CI 的 `go vet ./...`（`.github/workflows/ci.yml:207`）多咬一口。**处置＝改名**为 `pristine/models.go.pristine`（⛔ 未删，字节原样留着＝临时件只建不删），复尺 `go vet` 该目录＝`no Go files in ...\pristine` ⇒ 不再成包。⚠ 顺带量到一枚**与本腿无关的系统性暴露**具名给编排者：`git ls-files -- .scratch | grep -c "\.go$"`＝**264 枚 tracked 的 `.go`**（分布在 127 枚目录）今天就在 `go vet ./...` 的射程里，本腿那一枚只是其中刚被摘掉的一枚；这条归治理票，不归本腿。
- **占位符尺**（本件收尾时同发跑，读数在 §10 末条）：尺＝对那四个占位词形（以「待／未」起头那一族）**逐词** `grep -c`，⛔ 词形不在正文复写，否则尺咬到自己（本腿第一把尺就命中 2 枚、两处都是尺自己的词形＝假阳性，已就地更正）。反例记名＝`267-r1` 那把把竖线写进**基本正则**的尺（ERE 的 `|` 在基本正则里是字面竖线 ⇒ 四个词被当一条字面串、恒返 0＝瞎尺），本腿改用转义竖线 `\|` 复跑过。

## §9 我写错的读数（自我对抗，本腿自己写的）

1. **"40 s 有 15 枚先例"数错了**：15 是**带内喂值点总数**（40×10＋90×4＋60×1），单算 40 的先例是 **10 枚**。§2 已就地改成"10 枚现成先例（另 90×4、60×1，共 15 枚带内喂值点）"；原话留着不抹。
2. **★我吞过一枚标题**：填 §0＋§1 那次 Edit 的 `old_string` 覆盖到了骨架的 `## §2` 标题与其正文 ⇒ 现树的 §2 是后来重写回去的。发现尺＝`git show 00f0ef97:<本件> | grep -n "^## "`（骨架 11 枚）与现树标题名册对撞（那一度只有 10 枚）。⇒ 教训具名：**同一枚文件里连做多次区间替换，必须拿"骨架标题名册"当尺复点**，否则少一节我看不见。
3. **派单那句"其余 9 枚答完卡就走＝零增量"我照抄了没先验** ⇒ 实测 `TestTicket101SessionGrantDoesNotCrossRestart` 81.70 s、`...ManualSwitchSurvivesRestart` 42.24 s：那族**真的在等这枚钟**（重启若干次、每次等一张没人答的卡到点）。总墙钟实测 ＋242.7 s vs 预估 ＋58 s。⇒ 写进 §5／§7 并给出可选更形，⛔ 没有替派单圆成"大致符合"。
4. **`docs` 那 24 行我一开始读成"零命中"**：第一把尺 `grep -rln ... | head -20` 被 `.scratch` 的文件挤满了前 20 行 ⇒ 假阴性。复尺（`for r in ...; do ... | wc -l; done`）才量到 `docs=24 tools=0 scripts=0`。⇒ 教训：**计数尺不许挂 `head`**。
5. **行号抄过一次的险**：§2 初稿直接写派单给的 `:504/:578`，`grep` 复尺发现本腿新断言插了 17 行、它们已在 `:521/:595`。已在 §2 具名标"迁移前／迁移后"两号。⇒ 与仓里今天的定式一致：**引用行号前先用那串字符当尺**。
6. **`gofumpt -l cmd/wisp/` 不是零命中**：`cmd\wisp\models.go` 那枚是预存 CRLF（`file`＋`gofmt -l`＋`git diff --name-only` 三把尺）。我差点把它写成"本腿全绿"，改为如实写"非零项，但非本腿所造"。⛔ 本腿没有去"顺手格式化它"（越界）。
7. **M1 的牙齿有边界，我不夸大**：M1 证明的是"句子或前缀变了就红"；**它没有**证明"提示不发时这一发红"——因为带内任意种子的 `lead ≥ 1 s`，"不发"那一支在 `_test.go` 写面内不可表达（M2 只证明了拒载那支）。这一格的正写已在 §4.3 末条，归口＝编排者是否另派 approval 层的腿。⛔ 没把"两发突变"写成"三发已齐"。
8. **我差点污染 CI**：§8 那枚 `pristine/models.go` 副本是本腿造的、本腿量到它进模块并改名消解。⇒ 教训：**诊断用的原件改名加后缀，别留 `.go`**。
9. **`--author` 不能当"我的 commit"尺**：共享配置下所有腿同一个 `user.name`，`git log --author=...` 把 `267-a3` 的件也捞进来 ⇒ 本件名册一律按 hash 具名。
10. **★观察对象搬家逐枚复检（派单硬要求）**：
    - #2（`panel_pump`）／#6（timeout 那发）抬到 **31**：观察对象仍是"卡开合／超时本身"，凭据＝§3 那两行 stdout（`[warning] ...` ＋ `[dismissed] 审批超时（31 秒未确认）...`）与 `TestRunBooksWithASnapshotOfItsLiveQueue` 32.60 s 走的是超时支；ctx 那枚配套件就是为守这一格。
    - #1（`run_mode101`）抬到 **40**：该族观察对象＝档位读写／重启存活，仍然成立（迁移后五枚 101 用例全绿，含三枚子测首次登记）⇒ 没搬家；**但它顺带暴露"这族其实在等这枚钟"**（见第 3 条），这条代价先前没人量过。
    - #3/#4/#5/#7/#8/#9/#10/#11/#12 抬到 **40**：实测各 1.2–4.0 s（远小于 40 s）⇒ 40 s 那枚钟从未轮到，答完卡就走的原话在这些枚**为真**；同发放宽性读数：`TestL1Veto...` 4.54 s（含两枚子测）、`TestTicket224LiveGrant...` 3.96 s。
    - ⇒ 结论：**零枚发生观察对象搬家**；错的那格不是"抬错值"，是"没人量过 40 s 让哪一发真等满"（第 3 条）。
11. **§0 起手读数我只截了 30 行**：`head -30` 的 `.scratch` 名册不完整（同一文件里 `??` 项远多于 30 行）⇒ 那一节的话说成"cmd/internal 零条"是完整的（两棵子树各自单独取过），但"全树脏面"这句我不能声称数全过。已按此限定措辞。
12. **我把自己写成交件时还数错过一次枚数**：§8 初稿写"本腿共四笔"，而这张表写完时第 5 笔已经在计划里（就是补这一条的那一笔）。尺＝`git log --oneline f9bd0eb7..HEAD` 里本腿的 hash 逐枚点数＝00f0ef97／32e74479／5fa0d28c／aac52ab2／（本笔）＝**5 笔**；已把 §8 的名册与枚数改成"共五笔"，原话留在本条。**为什么这算缺陷级而不是笔误**：交件判语里"名册完整"是我自己下的判据，枚数错了那张表就不闭合，下一枚腿照它核面会以为少了一笔没交。

## §10 交件判语

- **做完的**：12 枚带外种子全部抬进带内（10 枚 `%d` 喂值点→40 s、#2 与 #6→31 s、#1→40 s 且同行 `l1_window_sec = 1` 一字未动）；★配套件 `panel_pump_test.go:112` 的 ctx 30 s→45 s 已连带抬；★票 267 AC#2 那半格的**正向读数已拿到并写死**（种 31 ⇒ 那句 `[warning] 审批将在 30 秒后自动拒绝，请尽快确认` 真到达装配根控制台 stdout，M1/M2 两发突变各自红、红句逐字在 §4.3）。
- **门禁**：`cmd/wisp` 整包 `=== RUN=323／PASS=226／--- FAIL=1／--- SKIP=0`（10:13:30→10:21:08，PATH 带 `third_party/sherpa-onnx`，rc=1）；`internal/config` `ok 1.212s`；`d22scan` clean rc=0；`check-path-length-budget --with-self-test` VERDICT GREEN rc=0；`go vet ./cmd/wisp/` rc=0；`gofumpt -l cmd/wisp/` 只剩预存的 `models.go`（CRLF，非本腿）。
- **红名作差**：基线 16 枚**逐名全转绿**，账算得平（211＋16−1＝226）；新增红 **1 枚**＝`TestAC14GoSideEvalPushReachesThePage`，红句逐字在 §6，四条排除性读数齐备，⛔ 本腿没有自行判它"无关"，**标〔未归因，归编排者〕**。
- **墙钟**：实测 ＋242.7 s（不是 ＋58 s），逐枚具名；本腿⛔未做任何"为快而压种子／压超时／加 Skip"的动作。
- **交件判定**：**本腿这一腿可收**（出口甲的射程全部落地，AC#4 的凭据里"红名逐名比对＋时刻"这一半已齐）；**留给编排者的三格**＝① 那枚新增红的裁决（同 HEAD 复跑一次即可归因），② `run_mode101` 族是否改取 31 以省 9 秒×2（要重跑整包才配得上同发读数），③ `panel_pump_test.go:124` 那句注释里的过期数字 2（写面闭集外，需具名解冻或另派）。
- **占位符自尺（同发读数）**：10:31:07 第一把尺命中 **2 枚**，两处都是尺自身写下的词形（假阳性，已在 §8 那一条就地更正并把词形从正文摘掉）；复尺逐词计数读数与时刻写在下一行，⛔ 本件交件时全文不含任何一节处于未完成状态。
  - 复尺读数：**10:33:01** 同发两把——逐词 `grep -c`＝**0／0／0／0**，合尺（转义竖线四词一并）＝**0** ⇒ 交件时全文零枚占位。（中间态 10:31:07＝2、10:32:07＝逐词 0，都是尺自身的词形造成的，已就地更正；⚠ 10:32:07 那次只跑了逐词那一把，合尺补跑在 10:33:01。）
  - ★那一行之后又复尺四次，全都是 0：10:33:50（对已入库的 `aac52ab2` 版本）／10:34:13／10:34:38／**10:35:33＝交件时刻那一把**（合尺 `grep -c`，零命中返回 rc=1 属这台机的正常形状）⇒ 写完 §8 名册与 §9 第 12 条之后没有重新引入任何占位词形；交件时全文 11 节（§0–§10）标题齐、254 行（字节数不在本件里自指——每改一行它就漂，终值以交件回报那把尺为准）。
