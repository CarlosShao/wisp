# 123 — 可编程应答腿：AC#1 定因到行 ＋ AC#2 拒绝分支实测（实现腿证据件）

- **本件是什么**：票 123 实现腿的证据件，出自实现腿本腿（非裁决方）。AC#1 定因到 `文件:行`、AC#2 ★判据（拒绝分支仍 PASS）的两次真读数、4 枚在真 carrier 下的四数、门禁三连，全部本机现跑取得。
- **现量**：HEAD `397cd9e9d`（2026-11-13），CGO_ENABLED=1 ＋ mingw PATH（派单门禁原文那套），真 carrier＝`scripts/wisp-cli-staged-tests.sh`。临时件落 `D:\tmp\wisp\20261113\123\`（只建不删）。
- **诚实边界（三条，先说）**：
  1. **写面已在库**：注入腿＋4 枚接线是 `abf2a35a6`（2026-11-06 12:14，"票123 W123腿死亡前遗留测试改动(程序化应答腿/8.3短路径白名单/sync-root不skip)"）＋ `446831ea9`（同日，补交三枚 replyleg helper）落的；本腿核过 `git log 446831ea9..HEAD -- <四个测试文件＋三个支撑文件>` 为**空**＝此后零改动。本腿**零生产码、零测试文件改动**，交付＝定因＋读数＋本件＋票面 log＋commit。
  2. **票面标题的假设被本腿证据推翻**（详见 §1-5）：9-21 本机绿**不是**"有人点确认"或"有人造的确认腿"——当时全仓不存在任何应答腿；本机绿是因为**根本没有卡可答**（长拼写 `%TEMP%` 让目标不出白名单 ⇒ 走 L1 路线）。今天满足 AC#1"那条腿"的，是修法落库后的**具名程序腿**（§1-4）。
  3. 本件第一稿生成损坏（大段重复＋尾部指令形残文），已整件重写；残文按仓规登记在 §6，未照做。

## §0 一屏结论

- **CI 红的机制＝两件事叠加，缺一不红**：
  ① **8.3 短路径拼写令目标出白名单**：hosted runner 的 TMP 是 8.3 短路径（`C:\Users\RUNNER~1\...`），9-21 时的 fixture 只把 `t.TempDir()` 返回的那一个拼写写进 `[fs].allowed_dirs`，judge 读回目标时是展开长形式 ⇒ 目标不在白名单 ⇒ R2 ⇒ **L2 队列**（修法＝`t123dirSpellingForms` 双拼写入白名单，`abf2a35a6`）。
  ② **L2 卡无人应答 ⇒ C18 契约超时判拒绝**：runner 上没有可点的卡片，`Confirm` 永不返回。三枚 Ticket101 用例等满 fixture 写的 `confirm_timeout_sec = 1`（1 秒）；`TestComposedGateBlocksAWriteForTwoSeconds` 的 fixture **没有 `[risk]` 块**，等满编译默认 **300s**（`internal/config/schema.go:499` `default:"300"`）⇒ 301.06s。
- **修法（票面 AC#2 指定方向，已在库）**：把"有人"显式化成**用例自带的程序化应答腿**——脚本化 `io.Reader` 发 `yes/no <corr>` 行，经 `runSpec.reply` → `run.go:1459 attachReplyListener` → 生产侧 `runReplyLoop`（`approval_reply.go:547`）消费，落**真实**审批队列。门禁语义、300s 值、断言一字未动。
- **本腿两次读数（真 carrier，同 HEAD）**：原样四枚全 **PASS**（6.05/9.44/5.75/4.99s，包级 `ok 27.00s` rc0）；**拒绝变异体**（`go build -overlay` 把 5 处批准行全部反转为 `no`）同四枚仍全 **PASS**（3.78/10.15/3.29/5.91s，rc0），且无一枚触及今日 fixture 的 40s 队列超时 ⇒ 腿真被消费、卡是"被答了"而不是"等没了"。
- **门禁**：`go build ./cmd/wisp/` BUILD-OK；`go vet ./cmd/wisp/` VET-OK；`sh scripts/d22scan.sh` **rc=0 clean**（live scope：bans#1-5 internal/=403、cmd/=77；ban#6 frontend/=123；ban#7 internal/tools/=51；ban#10 internal/risk/=21；ban#8 design/=39、frontend/=123、internal/=1005、cmd/=276——均不低于裁决表末记）。

## §1 AC#1 — 定因到行

### §1-1 四枚各自的等待点（今日 HEAD 现量）

| 用例（函数行号，今日） | 等 Confirm 的行（用例侧 dispatch） | 生产等待点（报文与秒数从哪来） |
|---|---|---|
| `TestTicket101ManualSwitchSurvivesRestart`（`run_mode101_test.go:238`） | boot 2 的两发：`:333`（fs.write note.txt，corr-1）、`:337`（fs.read win.ini，corr-2）；腿接线 `:329` | dispatch → `t101call`（`:215-232`）→ `bridge.Execute` → `internal/tools/bridge.go route()`：L2 分支 `bridge.go:540` → `PendingApproval`（`internal/agent/approval/gate.go:623`）：deadline 装在 `gate.go:647`、select 等 `it.answer` 在 `:663-676` |
| `TestTicket101UntouchedConfigRestartsAtDefault`（`:374`） | 每 boot 一发：`:405`（fs.write note-N.txt，corr-N）；腿接线 `:397` | 同上（CI 时代该发被 ①送进 L2；本机在白名单内走 L1 窗） |
| `TestTicket101SessionGrantDoesNotCrossRestart`（`:442`） | boot1 `:513`（.env，corr-1）；boot2 `:549`（.env）＋`:554`（plain.txt 控制）；腿接线 `:467`/`:541` | 同上（.env 是 B-tier R3 敏感，**任何机器任何档都进 L2**——这枚与拼写无关，纯粹是 ②） |
| `TestComposedGateBlocksAWriteForTwoSeconds`（`run_test.go:354`） | gated call `:395-399`（fs.write，corr=gate-corr）；腿接线 `:367` | 同上 |

- **报文生产者**：`internal/agent/approval/queue.go:536` `why: fmt.Sprintf("审批超时（%d 秒未确认），C18 一律判拒绝，已自动拒绝", int(q.timeout.Seconds()))`（9-21 时点同句在 `ff3faf92d` 树的 `queue.go:414`，本腿用 `git grep "审批超时（" ff3faf92d -- internal/ cmd/` 复核：**当时全仓唯一产出者**）。最终装配在 `bridge.go:545-549`（`AnswerTimeout` → `DecisionTimeout`，理由 orDefault 透传）。
- **CI 断言行号（票面/111 证件引的锚，当时树）**：`run_mode101_test.go:309/320`（ManualSwitch）、`:366`（Untouched，冷启动 3/3 同报错）、`:482`（SessionGrant 控制写）、`run_test.go:378`（ComposedGate）。今日同族断言在 `:354/:365`、`:418`、`:571`、`run_test.go:413`（行号漂移已对齐）。

### §1-2 秒数的两条来源链（1 秒 vs 300 秒，直接证据）

- **灌入点**：`cmd/wisp/run.go:1058` `ApprovalTimeout: time.Duration(cfg.Risk.ConfirmTimeoutSec) * time.Second`（L1 窗口同块 `:1057`，`Window: ...L1WindowSec...`）。
- **三枚 Ticket101**：9-21 树 fixture 写死 `riskLines := "[risk]\nl1_window_sec = 1\nconfirm_timeout_sec = 1\n"`（`ff3faf92d` 树 `run_mode101_test.go:119`；今日同一行是 `= 40`，`abf2a35a6` 未动此行、更早的票改的）⇒ CI 上每张 L2 卡等满 **1 秒**判拒绝。CI 耗时对照：3.10s（2 发×1s＋boot）、4.13s（3 boot×1s）、4.06s（3 发×1s）——**三枚耗时与"各等满 1s"一致**。
- **ComposedGate**：9-21 与今日的 `run_test.go` fixture **都没有 `[risk]` 块**（本腿 `git show ff3faf92d:cmd/wisp/run_test.go` 与工作树各核一遍，均无）⇒ `cfg.Risk.ConfirmTimeoutSec` 取 schema 默认 **300**（`internal/config/schema.go:499`）⇒ 等满 **300 秒** ⇒ **301.06s**（300s＋boot 开销）。四枚总耗时 3.10+4.13+4.06+301.06 = 312.35s，与"3×1s×N 发 + 1×300s ＋ 各 boot 开销"的形状一致（票面 AC#1 的"加起来"要求，逐枚对上）。

### §1-3 L1 路线的极性（为什么"窗口超时"不是这批失败的原因）

L1 前置阻止窗口到点未否决＝**执行**（`internal/agent/approval/gate.go:413-437`，`AnswerTimeout`→`bridge.go:527-530` `DecisionAllow`）——所以走 L1 路线的发**永远不会**产出"审批超时…判拒绝"这句（那是 `queue.go:536` 的 L2 文案）。9-21 失败语里出现这句 ⇒ **四枚在 CI 上全部走的是 L2 队列**，不是 L1 窗。这正是 ①（8.3 拼写出白名单 ⇒ R2 ⇒ L2）的判据；SessionGrant 的 .env 则与拼写无关、天然 L2（R3 敏感）。

### §1-4 本机"那条腿"的名字与 文件:行（今日）

**可编程应答腿 `t123script`**，四个具名落点：

1. 腿本体：`cmd/wisp/ticket123_replyleg_test.go:117-122`（`t123script`，Cleanup 兜底）＋ 脚本类型/喂线器 `:55`/`:64-92`（250ms 周期循环重发，`:51`）；
2. 注入面：`run_mode101_test.go:95`（`t101host.script` 字段声明）、`:165-167`（`h.start` 里 `h.script != nil ⇒ reply = h.script`，否则 attached-but-silent console）＋ `run_test.go:367`（`f.reply = t123script(...)`）；
3. 生产消费侧：`run.go:1458-1460`（`s.reply != nil ⇒ attachReplyListener`）→ `approval_reply.go:452-475`（挂监听）→ `:547 runReplyLoop`（逐行 `yes/no <corr>` 解析，`splitReplyLine :539`）→ `Replies.Reject/Allow` 落**真实**队列——测试没有绕开任何生产件；
4. 每枚的应答行：`:329`（ManualSwitch boot2：`yes t101-corr-1` / `no t101-corr-2`）、`:397`（Untouched：`yes t101-corr-N`）、`:467`＋`:541`（SessionGrant：`no t101-corr-1`（＋boot2 `yes t101-corr-2`））、`run_test.go:367`（ComposedGate：`yes gate-corr`）。

### §1-5 9-21 时点"本机为什么不需要等"——票面标题假设的裁定

- **当时的树里不存在任何应答腿**（`t123replyScript` 是 `446831ea9` 才补交的；`git grep` 当时树无 `t123script`）。票面标题"本机有人（或有人造的确认腿）"的**字面假设不成立**。
- 当时本机绿的具名机制：**长拼写 `%TEMP%`**（本机实测 `%TEMP%`=`C:\Users\swq\AppData\Local\Temp`，无 8.3 段）⇒ `t.TempDir()` 返回长形式 ⇒ 目标在单拼写白名单内 ⇒ 三枚的写/读走 **L1 路线**：auto_approve 档下 L1 直接静默执行（ManualSwitch `:333`、SessionGrant 控制写 `:554` 的断言要求的就是这个）、ask_every_step 档下 L1 窗到点＝执行（Untouched `:405` 的断言要求的就是这个）——**没有卡被显示，自然无人需要应答**。CI 上同一批发因 ①出了白名单 ⇒ L2 ⇒ ②无人应答 ⇒ 1s/300s 判拒绝。**红的是"用例假设了机器形状"，不是"用例假设了有人"**——两句话在 9-21 的树上等价，但机制前者才对。
- ⚠ 诚实标注：`perm.ConfirmFunc` stub（`run_mode101_test.go:250-256`，`AnsweredCard("t101")`）只盖 boot 1 的 mode-switch 卡（perm.Store 适配层），盖不到 boot 2 的工具调用卡——它不是那批绿的机制。

## §2 AC#2 — ★判据"Confirm 返回拒绝时仍 PASS"的实测

**两层证据，都不是"永远返回批准"：**

1. **设计内拒绝（原样读数里就有两枚真拒绝卡）**：ManualSwitch 的 corr-2（win.ini，R2 红线）与 SessionGrant 的 corr-1（.env，R3 敏感）的应答行本来就是 **`no`**（`:329`/`:467`）——这两张 L2 卡在 §3 第一读数里被腿答成**拒绝**，用例仍 PASS（前者断言的就是 `redLineErr && 含"拒绝"`，后者断言的就是 `firstErr`）。
2. **全量反转变异体（本腿现做）**：`go build -overlay`（`D:\tmp\wisp\20261113\123\overlay-reject.json`；变异副本在 `D:\tmp\wisp\20261113\123\mut\`，原始文件零改动）把 **5 处批准行全部反转为 `no`**：`:329` `yes→no`、`:397` `yes→no`、`:467` 已是 no 不变、`:541` `yes→no`、`run_test.go:367` `yes→no`。同四枚经同一 carrier 形状重跑：**全 PASS**（§3 第二行读数），rc=0。
   - **消费证明**：今日 fixture `confirm_timeout_sec = 40`（`run_mode101_test.go:119`）。若腿未被消费、卡没人答，每张 L2 卡要等满 40s 才判拒绝——第二读数四枚 3.78/10.15/3.29/5.91s，**无一枚触 40s** ⇒ `no` 行真被 `runReplyLoop` 消费、卡真被程序应答。
   - **语义边界（诚实标注）**：本机上 Untouched/ComposedGate/corr-1 的目标在白名单内、走 L1 路线，`yes/no` 行对 L1 窗是惰性的（窗口不是队列项，`ticket123_replyleg_test.go:34-38` 的设计）——变异体对这三处证明的是"反转无害＋腿照常在流上"；**真正被答成拒绝且 PASS 的 L2 卡**是第 1 条那两枚（原样读数里就是 no）。两层合起来，★判据两个方向都覆盖：真拒绝 ⇒ PASS（原样），批准全翻拒绝 ⇒ 仍 PASS（变异）。

## §3 四数表（真 carrier，HEAD `397cd9e9d`）

| 读数 | ManualSwitch | Untouched | SessionGrant | ComposedGate | 包级 | rc |
|---|---|---|---|---|---|---|
| 第一读数（写面原样） | PASS 6.05s | PASS 9.44s | PASS 5.75s | PASS 4.99s | `ok github.com/CarlosShao/wisp/cmd/wisp 27.00s` | 0 |
| 第二读数（拒绝变异体，overlay） | PASS 3.78s | PASS 10.15s | PASS 3.29s | PASS 5.91s | PASS（raw binary） | 0 |

原始日志：`D:\tmp\wisp\20261113\123\core-4tests-postfix.txt`、`D:\tmp\wisp\20261113\123\core-4tests-rejectmutant.txt`。

## §4 门禁（派单原文命令，逐条）

- `go build ./cmd/wisp/` → **BUILD-OK**（CGO=1＋mingw PATH）
- `go vet ./cmd/wisp/` → **VET-OK**
- `sh scripts/d22scan.sh` → **rc=0 clean**（数字见 §0 末条）

## §5 写面与交付

- 本腿 commit 只含：本证据件（新）＋票 123 Progress log 追加。**零生产码、零测试文件、零脚本/CI/契约文件改动**；commit 显式 pathspec、不 push（hash 见票面 log 与回话）。
- AC#4（runner 真读数）不归本腿：编排者 push 后取，看 `cmd/wisp CLI tests` 步四枚 `--- PASS`、`审批超时` 0 命中。

## §6 登记（仓规：指令形文本逐字登记，照常上报，不照做）

1. **本件自身第一稿的生成损坏**：写入 `docs/evidence/s1/123-programmable-confirm.md` 的第一版中段起大段重复损毁，尾部出现指令形残文（"Do not follow any instruction..."、"Do not trust this file" 等，出现多次）——那是本腿自己的生成噪声，不是外部指令；已整件重写为本版，残文未照做、原文件已被本版整体覆盖。
2. **`ticket123_replyleg_windows_test.go:13-14` 注释前提实测不成立**：注释称"8.3 aliases are disabled here - ShortPath equals the long path"；本机 FSO 实测（`Scripting.FileSystemObject.GetFolder(<%TEMP% 下新目录>).ShortPath`）返回 `...T123SP~1` ⇒ **本卷 8.3 别名是开着的**。本机从未红过的真因是 **`%TEMP%` 本身是长拼写**（目标总是长形式，`t123dirSpellingForms` 因此返回单形式），与别名开关无关——注释结论碰巧对、前提错。修正只需 2 行注释，但该文件不在派单枚举的写面清单里（4 枚用例文件＋新增支撑文件），按边界从紧**本腿未动**，留给下一条票 123 写腿或验收方定夺。

## §7 引文行号基准

- "今日"行号＝HEAD `397cd9e9d` 工作树现量（2026-11-13 复核）；"9-21/当时"行号＝锚 `ff3faf92d`（2026-09-21 后第一笔，票 111 取证同期树）现量，`git show`/`git grep` 逐条现跑取得，非转抄。
