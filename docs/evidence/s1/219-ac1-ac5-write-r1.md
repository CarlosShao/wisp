# 票 219 AC#1 ＋ AC#5 写腿（实现者交付 · r1）

- **实现者**：本腿（`cmd/wisp/**` 本波唯一写腿）。**裁决者 ≠ 实现者**：本件只是实现者的自陈与可核读数，
  **不构成验收**；AC#1／AC#5 的翻格权在收件腿（`docs/evidence/s1/219-ac1-ac5-receipt.md` 那一族）。
- **基线**：`e108b0945855dcfabaecaa67c22996d0766186c7`（收件腿的锚点是 `cbf88954e…`；两者之间无本票写面改动，
  `git show HEAD:cmd/wisp/resident_windows.go` 现读仍逐字是 `:568  if src != nil && src.run != nil {`
  ／`:571  rb.attachTrayAnswers(ra.trayAnswerExecutor(src.answers))`）。
- **派单**：`.scratch/wisp/probes/orch/20261115/txt/219-ac1-ac5-write-dispatch.txt`（32 行，全文读过）。
- **票面**：`.scratch/wisp/issues/219-approval-card-three-reply-buttons-and-a-reason-box.md`（188 行，全文读过，
  特别是 `:164-188` 的「AC#1–#5 收件」那节与两条 ⛔ 不翻理由）。
- **写面**：`cmd/wisp/**`、`internal/tools/**`、`internal/panel/**`、`docs/evidence/s1/219-*`。
  本波**实际只写了 `cmd/wisp/**` ＋ 本件**；`internal/tools/**`、`internal/panel/**` **零改动**（§7 三把尺第 3 行）。
- **硬约束**：`thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43 表／`frontend/**`／`design/**`
  **一字未动**；**没有放宽任何既有断言**；新协程走 `observe.Default.Spawn`（`approval_always.go:154`，既有那一发，未新增协程）。

---

## 0. 三把尺（现跑）

| 尺 | 命令原文 | 读数 |
|---|---|---|
| 基线提交 | `git rev-parse HEAD` | `e108b0945855dcfabaecaa67c22996d0766186c7` |
| 写面清单 | `git status --porcelain -- cmd/wisp internal/tools internal/panel docs/evidence/s1/219-ac1-ac5-write-r1.md` | `M cmd/wisp/approval_always.go`／`M cmd/wisp/approval_reply.go`／`M cmd/wisp/panel_resident_pump.go`／`M cmd/wisp/resident_approval_windows.go`／`M cmd/wisp/resident_windows.go`／`?? cmd/wisp/approval_219ac1ac5_test.go`（`internal/tools`、`internal/panel` 行内**零条**） |
| D22 静态扫 | `cd tools/d22scan && go run . -root ../../` | `d22scan: clean - no D22 ban violations`（486 枚产码文件／`cmd/` 81 枚含注释与 `_test.go` 的 ban#8 也扫过） |
| 载具 | `bash scripts/wisp-cli-staged-tests.sh "-test.run=…"` | 走 Git-bash（`D:\work\soft\Git\bin\bash.exe`）；**一次一发**，本件每条读数注明它那一发 |

---

## 1. AC#1 ★ 双击启动那条：**从"只记一行日志"变成"有真执行者"**

### 1.1 修之前盘上是什么（逐字，基线树）

`git show HEAD:cmd/wisp/resident_windows.go` 现读：

```
568: if src != nil && src.run != nil {
569:     ra.bindRuntime(src.run)
570:     src.run.bindWaitingStateSink(rb.projectWaitingState)
571:     rb.attachTrayAnswers(ra.trayAnswerExecutor(src.answers))
572: }
```

`startResidentTaskSource` 在 `interactiveStdin()` 没有控制台时**在装配任何东西之前**返回 nil
（`cmd/wisp/resident_task_source_windows.go:44` 一带的自陈；`-H=windowsgui` 双击就是这一形），于是：

- `src == nil` ⇒ `:571` 那一行**根本不执行** ⇒ `rb.trayAnswer == nil`；
- 托盘三枚按下去走 `resident_ball_windows.go` 的 `trayAnswerGesture`，`fn == nil` 分支 ⇒
  **`recordBallGesture("tray-allow-once"/"tray-reject"/"tray-allow-always")` 一行手势记录，卡片不动。**

### 1.2 修之后：赋值点（`file:line`，全部现读）

| 件 | 全仓相对路径:行 | 内容 |
|---|---|---|
| **装配（新分支）** | `cmd/wisp/resident_windows.go:587-598` | `else { rb.attachTrayAnswers(ra.trayAnswerExecutor(ra.answerSurfaceWithoutPipeline(cfgMgr65, residentAnswerMode(rp)))) }` ← **这就是那枚赋值点**（`:593`） |
| 控制台形状（原样保留） | `cmd/wisp/resident_windows.go:586` | `rb.attachTrayAnswers(ra.trayAnswerExecutor(src.answers))`（票 228 的字面钉 `TestAC3TheWaitingStateDriveIsWiredIntoTheShippedPath` 仍绿） |
| 无管线答复面 | `cmd/wisp/resident_approval_windows.go:389` | `answerSurfaceWithoutPipeline(cfg *config.Manager, mode func() risk.Mode) *replySurface`，`live: &nativeCards{h: ra.cards}`（**本进程启动时就有的那张卡账**）、`audit: ra.residentAuditf`、`ctx: ra.root` |
| 长期分支的座 | `cmd/wisp/resident_approval_windows.go:402-408` | `seat: &widenSeat{gate: ra.gate, mgr: …, out: os.Stdout, mode: mode}` |
| 档读数 | `cmd/wisp/panel_resident_pump.go:166` | `residentAnswerMode(rp)`，与 `panel_resident_pump.go:150-154` 同一读法（`disp.Mode.(*panel.ModeWriteHandler)` → `Modes.PermissionMode`），nil 时回落 `risk.DefaultMode()` |
| 托盘执行者 | `cmd/wisp/resident_approval_windows.go:317` | `trayAnswerExecutor(surface *replySurface)`（票 228 既有；本波未改它的词表） |
| 允许一次 | `cmd/wisp/approval_reply.go:257` | `allow` → `Replies.Allow` → `DecideFromNative`（**真花单次 nonce**） |
| 拒绝 | `cmd/wisp/approval_reply.go:332`／`:340` | `reject` → `refuse` → `Replies.Reject`（本波签名改为 `replyReason`，见 §5） |
| 长期允许 | `cmd/wisp/approval_always.go:118` | `always` → 先印规则＋`STAGED`（`:147`）→ spawn `runWidening`（`:166`，`observe.Default.Spawn` 在 `:148`）→ **第二张 L2 卡** → `AddAllowedDir` → `WIDEN-APPLIED` |

**为什么这一形是真的、不是"假执行者"**：卡账（`ra.cards`）不是任务管线造的，是
`newResidentApprovalWithConfig` 在启动时建的，并由 `ballCardUI.Prompt` 每显示一张卡就 `Record` 一张
（`cmd/wisp/resident_approval_windows.go:1129` `u.ra.cards.Record(p)`）。也就是说 —— **"没控制台"缺的是执行者，不是缺卡。**
这一发补的是执行者：三条动词仍走同一个 `replySurface` 词表、同一条原生命令路、同一族审计行。

### 1.3 会红的行为测试（三枚，`cmd/wisp/approval_219ac1ac5_test.go`）

| 测试（行） | 断言的是**行为**，不是源码字符串 | 实测 |
|---|---|---|
| `Test219AC1TrayAllowOnceHasAnExecutorWithoutAPipeline`（`:104`） | 无管线形状按「允许一次」⇒ 被卡的那笔 `fs.write` **真落盘** ＋ 审计 `ANSWER-ALLOW corr=…`（原生路） | **PASS 2.97s** |
| `Test219AC1TrayRejectHasAnExecutorWithoutAPipeline`（`:155`） | 无管线形状按「拒绝」⇒ 目标文件 **`os.IsNotExist`** ＋ 门自己那条 `ANSWER-REJECT corr=…` | **PASS 2.67s** |
| `Test219AC1TrayAlwaysWritesTheWideningWithoutAPipeline`（`:206`） | 「长期允许」⇒ 先起**第二张 L2 卡**，答它之后 **`config.toml` 的 `[fs] allowed_dirs` 真多一行** ＋ 落盘日志里 `WIDEN-APPLIED … key=allowed_dirs …` 带那行文本 | **PASS 3.42s** |
| `Test219AC1TheNoPipelineExecutorIsWiredInTheShippedBoot`（`:291`） | **装配字面钉**（`runResident` 要真桌面，跑不了）：上面 :593 那一行在产码里在 | **PASS 0.00s**（**自认是字面钉，不当行为证明**，A1880） |

### 1.4 突变 M1（证明上面三枚是行为牙，不是摆设）

把 `answerSurfaceWithoutPipeline` 的第一行改成 `if true { return nil }`（即**回到修之前那一形**），
只跑 `-test.run=Test219AC1Tray`：

```
--- FAIL: Test219AC1TrayAllowOnceHasAnExecutorWithoutAPipeline (1.84s)
    approval_219ac1ac5_test.go:141: tray allow-once sentence = "tray answer: no answer side in this process -
      it assembled no task pipeline behind the card ledger, so no allow, reject or widening can be spent here;
      nothing was changed", want the surface's own allow sentence naming the card …
--- FAIL: Test219AC1TrayRejectHasAnExecutorWithoutAPipeline (1.76s)
--- FAIL: Test219AC1TrayAlwaysWritesTheWideningWithoutAPipeline (1.73s)
```

⇒ 三枚**同时红**、红在**行为**上（写没落盘／卡没被答），且红得很快（1.7s 级）。改回后全绿（§8 的 PASS 读数）。

---

## 2. AC#1 ② L2 的「不再弹卡」：**具名裁决 —— 不修，合同挡两处**

**今天成立的范围**：`internal/tools/bridge.go:413`（逐字 `if !sil.Silenced && sil.Level == risk.L1 {`）
⇒ 会话授权只在 **L1** 被当答案；L2 的同 (工具, 路径) 请求**每次都重新弹卡**。

**为什么不给它落点**（两条**互相独立**的冻结文字，都不是本腿能改的）：

1. `docs/specs/SPEC-06-security-gatekeeping.md:121` 逐字：「三条不可逾越：**L2 永不进持久授权（含会话级）**；
   C25/R4 升级不受会话授权覆盖；安全配置（`[risk]`/`[fs]`/`[net]`/`[plugins]`）放宽不适用会话授权」。
   ⇒ 「本次会话内允许某 L2 工具」**本身就是被禁的那一格**（`PLAN.md:2148` D45-3 同一句）。
2. `SPEC-06:115-116` 逐字：「单次工具调用内 N ≥3 个同质 **L1** 操作 → 合并为一个确认 … **L2 永不聚合**」，
   实现处 `internal/agent/approval/gate.go:653`（`// L2 is NEVER aggregated (D45-1): the batch argument is dropped here`）。
   ⇒ 「同类 L2 合成一张卡」这条替代路也被封着。

**裁决**：AC#1 那句「'本次会话内'复跑一次同类请求 ⇒ 不再弹卡」**在 L2 上不可能成立**，除非改
`SPEC-06 §8.3` 第一/第二条（改契约文字＝人工批准，属机主三类里的第一类）。**本腿不自作主张，也不拿"没做"充"做了"**：
- 该 AC 在 L1 上**已经成立**（`bridge.go:413` ＋ `internal/agent/approval/replies.go:351` `AllowSession`）；
- L2 上的期望若机主仍要，只能开一条 Q：**"允许会话级 L2 授权"或"允许 L2 聚合"要动哪条契约**——这条我**没有**写进台账
  （台账不在本腿写面），请编排者落 `Q##`。

---

## 3. AC#1 ③ 第三枚按钮「本次会话内」：**具名裁决 —— 本波不做，两个卡点都在写面之外**

**今天的状态**：托盘只有三枚答复项（`internal/ball/ball_windows.go:66-68`：`OnTrayAllowOnce`／`OnTrayReject`／
`OnTrayAllowAlways`），**没有**第四枚；「本会话内」只活在控制台动词里（`approval_reply.go:301` `session` →
`Replies.AllowSession`，`internal/agent/approval/replies.go:351`）。

**两个卡点（都不是"懒得做"）**：

1. **菜单项与回调在 `internal/ball/**`**，不在本腿写面（派单写面＝`cmd/wisp/**`、`internal/tools/**`、
   `internal/panel/**`、`docs/evidence/s1/219-*`；`cmd/wisp` 本波只许一枚写腿，`internal/ball` 另有其主）。
   要变成按钮＝`internal/ball/tray_windows.go` 加一枚 `menuAllowSession` ＋ `ball_windows.go` 加一个 `fire` 分支 ＋
   `Events` 加一个回调 ＋ `cmd/wisp` 加一个 `trayAnswerSession` 动词：**四处、跨两个包**。
2. **更硬的那半：无管线形状没有会话身份与授权账**。`Replies.AllowSession` 要 `Options.Grants`
   （`cmd/wisp/resident_approval_windows.go:614` `Grants: ra.grants`），而它由
   `bindResidentGrantLedger` 绑定，唯一调用点是 `cmd/wisp/resident_task_source_windows.go:344`（**在控制台闸门之后**）。
   ⇒ 双击形状里这枚动词即使有按钮，也只会落 `GRANT-RECORD-FAILED`（"这一发放了，额度没记住"）——
   在 AC#1 最要紧的那条启动路径上，它是个**会撒谎的按钮**。造第二个会话身份是装配根的事（票 265 的 A601 §4
   明确拒绝第二处铸造），不是本腿能开的。

**裁决**：**本波不做**，上述两条就是理由；要做得先定"无控制台时会话身份从哪来"。**没有**把它做成"看起来能按"的死按钮。

---

## 4. AC#5 ① 理由的第三去向（用户回看）：**落点 = `decisions` 动词**

**缺口的形状**（收件腿的读数，本腿复核）：`tool_call` 十三列（`internal/memory/schema.go:81-96`）**没有 reason 列**，
`approval: REPLY` 那族行**生产读者 0** ⇒ 用户写下的理由，除了模型和日志，**没有人能回看**。

**落点（本波实现，可敲可读）**：

| 件 | 全仓相对路径:行 | 内容 |
|---|---|---|
| 账 | `cmd/wisp/approval_reply.go:887-921` | `replyReviewKeep = 64` ＋ `replyDecision` ＋ `replyReview`（互斥锁；`add` 只留最新 64 条；`snapshot` 拷贝出门，读者不持锁） |
| 唯一的写口 | `cmd/wisp/approval_reply.go:472` `record` | **审计行与回看账写在同一句函数里**——回看账不会变成第二个真相源；`record` 是每条答复的必经之路 |
| 读口（动词） | `cmd/wisp/approval_reply.go:498` `decisions` | 新的呢（新的在前）：`<when> <kind> corr=… tool=… route=… 理由/明细=…`；空理由显示 `（未填理由）`；空账说"一张都没答过"（不是错误） |
| 接线 | `cmd/wisp/approval_reply.go:752`／`:736` | `case "decisions"`／`case "no"` 走同一条 `handle`；`help` 与未知动词的名单同步改了 |
| 无管线形状也有 | `cmd/wisp/resident_approval_windows.go:409` | `answerSurfaceWithoutPipeline` 也给 `review: newReplyReview()` |

**这条落点诚实边界（写清，不藏）**：

- 回看账是**进程内、有上限（64 条）**；**落盘的证据仍是 `approval: REPLY` 那族行**（`jsonl` 里带 `detail=`＋`when=`）。
  审计行是**持久**那份，回看账只是"不用翻日志就能答我刚刚答了什么"。
- **结构化/逐条可查询那一半没有做**：它要 `tool_call` 加列（`SPEC-02 §3` 的 DDL 是冻结形状，`docs/specs/**` 本波禁改），
  属改契约；本腿只把它记在这里，不擅自加列。
- **界面那一半没有做**：网页/面板上的"审批历史"要动 `frontend/**`（本编队零写面、归机主外部 agent）。
  面板侧 Go 数据面本波**没加键**（加一个没有读者的键就是假进度）。

---

## 5. AC#5 ② 空理由：**从"零测试"变成四读数一枚牙**

**行为（现读，未改语义，只加了"确定"）**：

| 读者 | 空理由时读到什么 | 出处 |
|---|---|---|
| 模型 | 门自己的那句 `用户拒绝了本次操作`（**不是空串，也不是本文件编的句子**） | `internal/agent/approval/queue.go:515` ＋ `internal/tools/bridge.go:552` |
| 操作员（控制台） | `已拒绝 <corr>（原生侧，未填理由）：未执行` | `cmd/wisp/approval_reply.go:340` `refuse` |
| 审计 | `detail=""`（**真话是"没写"，不是补一句**） | `cmd/wisp/approval_reply.go:472` `record` |
| 回看账 | `（未填理由）` | `cmd/wisp/approval_reply.go:498` `decisions` |

**牙**：`Test219AC5EmptyReasonIsANamedState`（`cmd/wisp/approval_219ac1ac5_test.go:490`）**PASS 2.39s**，
四读数逐条断言，并额外钉住"操作员的括号不许漏进模型文本"（两个朝向都钉）。

---

## 6. AC#5 ③ 超长理由：**从"静默砍"变成"截断可见、数目可读"**

**修之前**（收件腿的探针读数）：200 字进、160 字出、`…`／`...`／`截断`／`省略`／`+` 全 0 命中 ⇒ 静默。

**修之后**（`cmd/wisp/approval_reply.go:851` `washReplyText`）：

- 上限**一字未动**（`replyReasonMax = 160`，`:788` 一带），**没有放宽任何阈值**；
- 返回值变成 `(text, dropped)`：`dropped` 是**被砍掉的字数**（按洗过的 rune 计）；
- 超限时正文**末尾追加**一句：`…(理由过长，已截断 40 字，原文共 200 字)`（`replyTruncationMark`，`:794`）；
  标记**加在 160 之外**——理由：**能被自己砍掉的标记等于又一层静默**（这条写进注释了）；
- 这句标记随正文进了**三处**：模型收到的工具结果文本、`approval: REPLY … detail=` 审计行、`decisions` 回看账；
- 哦控制台那句多一句人话：`理由超过 160 字上限，已截断 N 字（截断标记在正文末尾，模型与审计都看得到）`
  （`replyReason` 的 `truncationNote`，`:811` 一带）。

**为什么选"标记"不选"拒收"**（机主口径：按主流 harness、按最有价值）：理由是**一张审批卡的回答**。
把过长的回答整条丢掉＝让"拒绝"这个决定因为理由太长而**没发生**，而 B 端用户真正要的是决定落地、理由可读。
主流形状也是"限长＋明示被截"（本仓 `internal/observe/redact.go:193-198` 的 `truncate` 就是
`...(truncated, N chars total)` 同一族），不是"你给我闭嘴重说"。
**没做的**：被砍掉的原文**没有留第二份**（不留一份无上限的操作员文本——这不在本票授权内），这一点标记本身说得清楚（"已截断 N 字，原文共 M 字"）。

**牙**：`Test219AC5LongReasonIsTruncatedVisibly`（`:372`）**PASS 0.00s**（200→`dropped=40`、标记带两个数、
标记不许占额度、**恰好 160 不许误报被截**）；`Test219AC5TheModelAndTheOperatorSeeTheTruncation`（`:419`）**PASS 2.06s**
（真跑一发 `wisp run`：模型文本／控制台句子／审计行／回看账四处同读）。

**突变 M2**（把上面改回静默砍，只跑 `-test.run=Test219AC5`）：

```
--- FAIL: Test219AC5LongReasonIsTruncatedVisibly (0.00s)
    approval_219ac1ac5_test.go:376: dropped = 0, want 40 (200 typed - 160 ceiling)
    approval_219ac1ac5_test.go:379: the cut reason carries no truncation mark: "长…（160 个长）"
--- FAIL: Test219AC5TheModelAndTheOperatorSeeTheTruncation (2.30s)
    approval_219ac1ac5_test.go:464: the model-visible refusal carries no truncation mark: "长…（160 个长）"
--- PASS: Test219AC5EmptyReasonIsANamedState (2.16s)   ← 静默砍不影响空理由，红的正是该红的两枚
```

---

## 7. 本波**没做**的，逐条具名（附理由）

| # | 没做 | 理由（不是"风险大"） |
|---|---|---|
| 1 | **L2 的「不再弹卡」** | §2：`SPEC-06:121`（L2 永不进任何会话级授权）＋`SPEC-06:115-116`（L2 永不聚合）两处冻结文字挡着；改契约＝人工批准。L1 上已成立（`bridge.go:413`） |
| 2 | **第三枚按钮「本次会话内」做成托盘项** | §3：菜单/回调在 `internal/ball/**`（非本腿写面）；且无管线形状没有会话身份与授权账（`bindResidentGrantLedger` 唯一调用点在控制台闸门之后），做成按钮＝一个会撒谎的按钮 |
| 3 | **`tool_call` 加 reason 列** | `SPEC-02 §3` 的十三列是冻结 DDL（`internal/memory/schema.go:81-96`），加列＝改契约；本波禁改 `docs/specs/**` |
| 4 | **网页/面板上的"审批历史"界面** | `frontend/**` 零写面（归机主外部 agent）；面板 Go 侧本波也没加无读者的键 |
| 5 | **被截掉的原始超长理由留档** | 本票没授权新增一份无上限的操作员文本；标记里已写明"原文共 M 字" |
| 6 | **整包 `go test ./cmd/wisp`** | 读数见 §8；⚠ 本机条件性 native AV（`0xC0000005`）——本波**没遇到**时如实写"跑到底"，遇到就写"第 N 秒终止＋最后一条 RUN" |

---

## 8. 本波跑过的测试（载具一次一发，逐发读数）

| 发 | 命令（载具后参数） | 结果 |
|---|---|---|
| A | `-test.run=Test219AC1\|Test219AC5` | **7 枚全 PASS**（2.97s／2.67s／3.42s／0.00s／0.00s／2.06s／2.39s，包 14.00s） |
| B | `-test.run=TestAC2Tray\|TestAC3\|TestReplyListener\|TestTicket201\|Test201\|Test219AC2\|TestRevoke\|Test226\|Test375\|TestCloseJoin\|TestAlwaysWrite\|TestTicket21S2` | **38 枚全 PASS**（含票 228 三枚托盘行为牙、票 201 全套、票 219 AC#2 四枚、票 375 四枚；包 85.00s） |
| C | 突变 M1（`-test.run=Test219AC1Tray`） | 三枚**全红**（见 §1.4），改回后回绿 |
| D | 突变 M2（`-test.run=Test219AC5`） | 两枚**全红**（见 §6），改回后回绿 |
| E | 整包（`go test ./cmd/wisp` 全量） | 见本件末尾 §9（本波唯一一发整包） |

---

## 9. 整包读数（A1882 三件套：终点锚／最后一条 `=== RUN`／本波新测试在文件序截断点前还是后）

- **发次**：本波唯一一发整包 —— `bash scripts/wisp-cli-staged-tests.sh`（**不带** `-test.run`），
  **跑到底，没被 native AV 终止**（没有 `0xC0000005`）；包自报墙钟 **1572.00s**；
  日志 `.scratch/wisp/stage-cli/wisp-cli-tests.log`。
- **四数**：**RUN=878 ／ PASS=634 ／ FAIL=2 ／ SKIP=17**（末行 `FAIL github.com/CarlosShao/wisp/cmd/wisp 1572.00s`）。
- **FAIL 只有两枚，都是派单点名"本机既有红、不是你的"**：
  `TestAC4FocusReturnToPriorWindowGap33r5 (4.96s)`（前台被占）／
  `TestAC246DevLegIgnoresTheTestTaskInjection (37.15s)`（驻留球在跑）。
  派单点名的第三枚 `Test320FetchLegWritesTheKeyNowhereOnDisk` 本发**绿**（单跑/整包都绿）。
- **终点锚 `TestZZ315A1` 命中**：三枚全跑（`…RootLoopFirstBridgeCall`／`…SiblingChildFirstBridgeCall`／
  `…ProjectInstructionsSeedCannotStrandTheFirstTask`），最后一枚 **PASS 0.24s**，序号 **878/878**
  ⇒ 这一发**跑到了最后一条**，不存在"文件序截断点"。
- **最后一条 `=== RUN`**：`TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`。
- **本波 7 枚新测试的位置**：序号 **10–16 / 878**（`approval_219ac1ac5_test.go` 文件序靠前），
  即**在截断点之前**（本发无截断点）。逐枚：`…TrayAllowOnce…` 2.22s／`…TrayReject…` 1.98s／
  `…TrayAlwaysWrites…` 2.20s／`…TheNoPipelineExecutorIsWiredInTheShippedBoot` 0.00s／
  `…LongReasonIsTruncatedVisibly` 0.00s／`…TheModelAndTheOperatorSeeTheTruncation` 2.04s／
  `…EmptyReasonIsANamedState` 1.93s —— **7/7 PASS**。

