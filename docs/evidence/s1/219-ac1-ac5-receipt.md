# 票 219 AC#1–AC#5 收件（非实现者裁决，腿 A）

- **裁决者 ≠ 实现者**（`SPEC-12 §4.3` #1/#3、D22 双角色）。本腿只读＋裁决，产码/票面/台账/阈值/契约/`frontend/**` **零写**。
- 锚点：**`cbf88954edb31e053ad9a3fe506560c169a8ac7d`**（2026-11-15 15:13:59 +08）。
- 派单：`.scratch/wisp/probes/orch/20261115/txt/219-ac1-ac5-and-214-ac2-ac6-dispatch.txt` 第 12–22 行。
- 票面：`.scratch/wisp/issues/219-approval-card-three-reply-buttons-and-a-reason-box.md`（152 行，全文读过）。
- **AC#6 已勾**（`docs/evidence/s1/219-ac6-receipt.md` ＋编排者自跑 `TestRunPacketLightsTheRosterRowOfItsLiveL1Window` PASS 12.25s），本件**只裁 AC#1–#5**。
- 硬规：A1875／A1880（生产调用点 `file:line` ＋ 会红的测试名，缺一只能判"部分"；分清**行为牙**／**字面钉**）、A1879（全仓相对路径＋命令原文）、A1881（时间条件要在锚点提交上另跑并写 sha/日期）。

---

## 0. 起手三把尺（现跑，读数）

| 尺 | 命令原文 | 读数 |
|---|---|---|
| 提交真在 | `git rev-parse HEAD` | `cbf88954edb31e053ad9a3fe506560c169a8ac7d` |
| `cmd/wisp` 写面干净（能跑单枚测试的前提） | `git status --porcelain -- cmd/wisp` | **空**（见 §3 每发突变前后的同一读数） |
| 环境坑（不是本票缺陷，但复跑必须知道） | `go test ./cmd/wisp -run '^…$'` | 缺 DLL 时 `exit status 0xc0000135` ⇒ 跑前必须 `$env:PATH = "<仓根>;" + $env:PATH`（`onnxruntime.dll`／`sherpa-onnx-c-api.dll`／`sherpa-onnx-cxx-api.dll` 都在仓根） |

**跑过的测试一律单枚**（`-run '^Exact$' -count=1`），**没有整包跑**（禁令：`cmd/wisp` 归写腿 452b；整包同一时刻只许一个）。未跑到的每一枚在下文逐处标〔未复跑〕。

---

## 1. 五格总表（五列：AC 原文逐字／判语／依据／可核命令／缺什么）

| # | AC 原文逐字（票面 `:70–:74`） | 判语 | 依据（生产点 `file:line` ＋ 会红的测试名） | 可核命令 | 缺什么 |
|---|---|---|---|---|---|
| 1 | 「三枚按钮各有真执行者（生产路径，不是测试里能跑）；"本次会话内"复跑一次同类请求 ⇒ **不再弹卡**；**进程重启 ⇒ 又弹**」 | **部分** | 允许一次／拒绝／长期允许各有**行为牙**（§2.1，M1/M3 实测红）；"本次会话内"**不是按钮**（只有控制台动词），且"不再弹卡"**只对 L1 成立**（`internal/tools/bridge.go:413` 逐字 `sil.Level == risk.L1`） | §2.1 每行的命令 | ①第三枚按钮缺形状；②L2 同类请求**仍弹**（合同如此）；③**双击启动形状没有执行者**（§2.2） |
| 2 | 「"长期"那一支必须留下**可撤销的一行**——审计里能查到"谁、什么时候、往 `allowed_dirs` 加了哪一行、来源是审批卡片"；并有一条判据钉住"加完之后撤销能**一字回到加之前**"」 | **满足** | `cmd/wisp/approval_always.go:140-141`＋`cmd/wisp/approval_reply.go:449-455`（`when=`）；撤销 `cmd/wisp/approval_revoke.go:66`→`internal/config/allowdirs.go:109` | §2.3 | 诚实边界：**"谁"只有 `corr`+`route` 粒度**（无操作员账号）；写入**下一次启动生效**（`approval_always.go:40-45` 自陈） |
| 3 | 「网页那侧送进来的 `allow` ⇒ **必须被判红/被拒**」 | **满足**（门今天在） | `internal/agent/approval/gate.go:857-867`（路由上拒，`ErrPanelAllow`）＋`internal/agent/approval/replies.go:431`；行为牙 `TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject` **PASS 2.29s** | §2.4 | 票面 `:72` 那句"今天它是退回状态、本票不许勾"**已过期**（门已修好）；AC 框仍空，本腿只判不翻 |
| 4 | 「种一条"面板侧来源的允许"样本 ⇒ **那把尺必须响**」 | **满足** | `internal/panel/l2_grant_boundary_test.go:1955-1975`（五个种植样本全被 `ParseComposerRequest` 拒）＋`:1977-1994`（已答路由走私 verdict 被丢） | §2.5 | 无（正控本来就该是测试里种的样本） |
| 5 | 「理由**三去向**逐条有读数（审计一行／模型收到的结果文本里带那句理由／**用户回看的那份**里也在）；**空理由**与**超长理由**各有确定行为（**不许静默截断**，票 147 那一族"截断要带可找回指针"同纪律）」 | **部分** | ①审计 `approval_reply.go:449-455`；②模型 `internal/tools/bridge.go:552`（M2 实测两处红）；③**无结构化落点**；④空理由**零测试**；⑤超长**静默截断**（探针读数 200→160，无标记） | §2.6 | ③一个"用户回看"落点；④两枚牙；⑤"可找回指针"或"拒收超长" |

**总：满足 3（AC#2／#3／#4）／部分 2（AC#1／#5）／不满足 0，共 5。**

---

## 2. 逐格明细（A1875 两件套：生产落点 ＋ 会红的测试名；行为牙／字面钉分开）

### 2.1 AC#1 —— 三枚答复按钮的**生产落点**链（每一跳都现读）

**①「允许一次」**（行为牙）

| 跳 | 全仓相对路径:行 | 内容 |
|---|---|---|
| 菜单项 | `internal/ball/tray_windows.go:112` | `appendItem(menuAllowOnce, "允许一次", false)` |
| 派发 | `internal/ball/ball_windows.go:775-776` | `case menuAllowOnce: b.fire(b.opts.Events.OnTrayAllowOnce)` |
| 装回调 | `cmd/wisp/resident_ball_windows.go:521` | `OnTrayAllowOnce: func() { rb.trayAnswerGesture(trayAnswerAllowOnce) }` |
| 手势 | `cmd/wisp/resident_ball_windows.go:812-821` | nil 执行者 ⇒ `recordBallGesture("tray-allow-once")`；否则 `sentence := fn(verb)` |
| 执行者 | `cmd/wisp/resident_approval_windows.go:340-341` | `sentence, err = surface.allow(card.CorrelationID)` |
| 答复面 | `cmd/wisp/approval_reply.go:242-263` | `s.live.h.Allow(s.ctx, corr)`，四类拒绝各有自己的句子 |
| 缝 | `internal/agent/approval/replies.go:316` | `Replies.Allow` → `DecideFromNative` |
| 路由器 | `internal/agent/approval/gate.go:751` / `:845` | `nativeAPI.Allow` / `DecideFromNative` |

会红：**`TestAC2TrayAllowOnceSpendsTheNativeGrant`**（`cmd/wisp/resident_tray_answer_228_windows_test.go:92`）。**行为牙**，不是字面钉——M1（§3.1）摘掉那一行 `surface.allow(...)` 后红在 `:124`「a call the tray allowed must execute」，即**被卡的那一笔写真的没落盘**，而句子断言先过了。

**②「拒绝」**（行为牙）

`internal/ball/tray_windows.go:113` → `internal/ball/ball_windows.go:777-778` → `cmd/wisp/resident_ball_windows.go:522` → `cmd/wisp/resident_approval_windows.go:342-343`（`surface.reject(card.CorrelationID, trayRejectReason)`）→ `cmd/wisp/approval_reply.go:317-344`（`reject`→`refuse`）→ `internal/agent/approval/replies.go:372` → `internal/agent/approval/gate.go:819`。

会红：**`TestAC2TrayRejectRefusesOnTheNativeRoute`**（同上文件 `:141`，断言目标文件 `os.IsNotExist` ＋ `route="native/reject"` 指纹）。**行为牙**。〔未突变〕——同一 switch 的第二支，M1 已证明这一族能被行为断言抓到；再摘一支只增加读数不增加信息，故省。

**③「长期允许」**（行为牙）

`internal/ball/tray_windows.go:114` → `internal/ball/ball_windows.go:779-780` → `cmd/wisp/resident_ball_windows.go:523` → `cmd/wisp/resident_approval_windows.go:344-345`（`surface.always(...)`）→ `cmd/wisp/approval_always.go:70-107`（`always`：印规则＋`STAGED`＋spawn）→ `:112-144`（`runWidening`：**第二张 L2 卡**→ `AddAllowedDir` → `WIDEN-APPLIED`）。

会红：**`TestAC2TrayAllowAlwaysRoutesIntoTheWideningVerb`**（`cmd/wisp/resident_tray_answer_228_windows_test.go:194`，断言审计里出现只有 `always` 才写的那条 `native/always` 指纹）。**行为牙**——M3（§3.3）把它改路由到 `allow`，红在 `:222`「the 长期允许 item did not reach replySurface.always」。

**④「本次会话内允许」——它不是按钮**（票面 §2 表格的中间那一枚）

`internal/ball/tray_windows.go:25-41` 逐字写着这一枚**为什么不在托盘里**（"「本次会话内允许」 is deliberately NOT here: the Answer verb table has four words and the session-scoped grant is ticket 224's"）。它的生产执行者是控制台动词：

| 跳 | 路径:行 | 内容 |
|---|---|---|
| 动词表 | `cmd/wisp/approval_reply.go:670` | `case "session": return s.session(corr)` |
| 答复面 | `cmd/wisp/approval_reply.go:286-311` | `s.live.h.AllowSession(...)` |
| 缝 | `internal/agent/approval/replies.go:351` | `Replies.AllowSession` |
| 路由器 | `internal/agent/approval/gate.go:774` | `nativeAPI.AllowSession` |
| 记授权 | `internal/agent/approval/gate.go:805` | `g.grants.Record(ctx, tool, p)` |
| 写盘 | `internal/session/grants.go:146` | 一行 `approval_grant`，`scope=session` |
| **读口** | `internal/tools/grant.go:65-99`（`sessionGrantID`）→ `internal/tools/bridge.go:412-421` | 只在 `sil.Level == risk.L1` 时替代一次询问 |

会红：`TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun`（**PASS 2.41s**，端到端）、`TestTicket224LiveGrantStopsTheL1Question`（**PASS 0.02s**，"不再弹卡"那一半）、`TestTicket224ProductionSessionDoesNotSurviveRestart`〔未复跑〕（"重启又弹"那一半）。

**⑤ 缺口（这就是"部分"的三条）**

1. **第三枚按钮缺形状**：托盘上是「允许一次／拒绝／长期允许」，票面 §2 要的是「本次／本次会话内同类／长期」。"本次会话内"有真执行者，但只在控制台动词里，**没有按钮**。
2. **"不再弹卡"只在 L1 成立**：`internal/tools/bridge.go:403-411` 逐字「SPEC-06 §8.3's first bullet ("L2 永不进入任何持久授权（含会话级）")」，代码 `:413` 只在 `risk.L1` 查授权。控制台与常驻进程发出的卡是 **L2**（`cmd/wisp/approval_reply.go:273-280` 自己写明）⇒ 同类 L2 请求**照样弹卡**。这不是实现偷懒，是合同，但 AC#1 的文字没跟着收窄。
3. **双击启动形状里三枚按钮没有执行者**（下面 §2.2）。

### 2.2 AC#1 的启动形状边界（这条比按钮本身更该被看见）

`cmd/wisp/resident_windows.go:568-572`：

```go
if src != nil && src.run != nil {
    ra.bindRuntime(src.run)
    src.run.bindWaitingStateSink(rb.projectWaitingState)
    rb.attachTrayAnswers(ra.trayAnswerExecutor(src.answers))
}
```

`src` 来自 `startResidentTaskSource`，而它**无控制台时提前返回 nil**——`cmd/wisp/resident_approval_windows.go:374-380` 逐字点名这个形状就是 **"precisely the double-click / Explorer shape D2 and ticket 228 call the process the user actually launches"**。⇒ 用户双击启动时：三枚菜单项**在**（`showMenu` 永远建它们，`tray_windows.go:105-114` 有逐字理由），但 `rb.trayAnswer == nil` ⇒ `cmd/wisp/resident_ball_windows.go:812-817` 走 `recordBallGesture`：一行 `slog.Warn`＋一行 `fmt.Printf`（`:1086-1094`），而 GUI 子系统启动**没有控制台可印**。

读数（现读，非推断）：`rb.attachTrayAnswers` 全仓**只有一个生产调用点**（`cmd/wisp/resident_windows.go:571`；`resident_tray_answer_228_windows_test.go:492` 是钉它的源码齿）。

### 2.3 AC#2 —— 「长期」的可撤销那一行（满足）

| 要素 | 生产落点 | 现读 |
|---|---|---|
| 谁 | `cmd/wisp/approval_reply.go:453` | `corr=%q`（哪张卡的哪条答复） |
| 什么时候 | 同上 `:454` | `when=%s` = `time.Now().Format(time.RFC3339)` |
| 加了哪一行 | `cmd/wisp/approval_always.go:140-141` | `detail="dir=%q section=fs key=allowed_dirs effect=next-start"` |
| 来源是审批卡片 | 同上 | `route="native/always"` |
| 撤销 | `cmd/wisp/approval_revoke.go:66` → `internal/config/allowdirs.go:109` `SetAllowedDirs` | 只减不加的单键写路径 |

会红（四枚）：`Test219AC2WidenAppliedAuditLineCarriesTheFourFields`（**PASS 2.25s**，逐字段读）、`Test219AC2RevokeRoundTripByteForByte`（**PASS 2.44s**，add→revoke **逐字节**回原）、`Test219AC2SecondRevokeOfTheSameRuleIsRefused`〔未复跑〕、`Test219AC2MutationRevokeWriteMustPersistForBytesToComeBack`〔未复跑，它自己内含一发 `-overlay` 突变〕。均 `cmd/wisp/approval_always_219ac2_test.go`。

已勾那一半的既有收件＝`docs/evidence/s1/219-ac2-revoke-roundtrip.md`。本腿**没有**复核它的每一句，只复跑了上表两枚并现读了三处生产点。

### 2.4 AC#3 —— 面板来源的 allow 必须被拒（满足，门今天在）

- `internal/agent/approval/gate.go:857-867`：`DecideFromPanel` 在**路由上**拒，返回 `ErrPanelAllow`（`:867` 逐字「面板来源的「允许」被服务端 API 直接拒绝（F2 第三层）」），**不看授权**。
- `internal/agent/approval/replies.go:431` `PanelAllow`：把授权按"泄露"处理（`offered==true` ⇒ 烧 nonce）。
- `cmd/wisp/approval_reply.go:358-380` `panelAllow`：宿主侧再包一层，`PANEL-ALLOW-REFUSED` 入审计。

行为牙（本腿实跑 **PASS 2.29s**）：`TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject`（`cmd/wisp/approval_reply_201_test.go:344`）。它一发里断言四件事：①面板投影不含令牌/不含允许；②`panel-yes` 被拒；③同一张卡随后的 `yes` 拿到「原生令牌无效」（**授权被烧**）；④`panel-no` 仍能拒且目标文件不落。

票面 `:72` 的「今天它是**退回**状态，所以本票**不许勾这一格**」——那是 09-29 的读数，**今天不成立**：门已在生产路径上工作。⇒ 本腿判**满足**，但把票面那句标为**过期**，⛔ 未改票面一字。

### 2.5 AC#4 —— 种样本必响（满足）

`internal/panel/l2_grant_boundary_test.go:1955-1975`：**五个**种植样本（`panel.approval.request` ＋ `outcome:"grant"`、带 `requestId/source` 的同形、`approval.decide`、`panel.l2.allow`＋`allow:true`、`panel.grant`＋`approved:true`）全部被 `ParseComposerRequest` 拒（且必须是 `errors.Is(err, ErrComposerRequest)`）。
`:1977-1994`：把六个 verdict 词（`outcome/allowOnce/decision/verdict/approved/grant`）塞进**四条已答路由**，要求逐条丢弃。

本腿实跑 **PASS 0.00s**（两个子齿全绿）。另一台尺的正控在 `tools/d22scan/scan_test.go:80`（种 `approval.decide` ⇒ ban #6 必响）。

### 2.6 AC#5 —— 理由三去向 ＋ 空/超长（部分）

| 去向 | 生产落点 | 牙 | 读数 |
|---|---|---|---|
| ① 审计一行 | `cmd/wisp/approval_reply.go:339`（`s.record("ANSWERED", corr, tool, route+"/reject", reason)`）＋`:449-455` | `TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel:335` | **命中**：M2 实测红于 `:336`「the audit never booked the reason the operator typed」 |
| ② 回给模型 | `internal/tools/bridge.go:552` `orDefault(why, "L2 审批未通过，已拒绝执行")`（L1 在 `:533`／`:536`） | 同一测试 `:320` | **命中**：M2 实测红于 `:321`「the model-visible rejection text lost the operator's reason」 |
| ③ 用户回看的那份 | **没有落点** | 无 | `internal/memory/schema.go:81-96`：`tool_call` 十三列里**没有 reason 列**；全仓**没有任何界面/快照**读 `approval: REPLY` 行（`Select-String internal/panel/*.go -Pattern 'REPLY ANSWERED|approvalHistory'` ⇒ **0 命中**）。唯一副本是 JSONL 审计行 |
| ④ 空理由 | `cmd/wisp/approval_reply.go:340-341`（「（原生侧，未填理由）」） | **零测试** | `grep 未填理由` ⇒ 只命中生产那一行 |
| ⑤ 超长理由 | `cmd/wisp/approval_reply.go:704-708`（`replyReasonMax = 160`）＋`:745-747`（到 160 就 `break`） | **零测试** | **静默截断**，探针读数见 §3.4 |

⇒ AC#5 的第三条去向**今天不存在**，第五条**直接违反 AC 文字**（"不许静默截断"）。故判**部分**，并具名两处：③ 缺一个落点、⑤ 缺一枚"可找回指针"或"拒收超长"。

---

## 3. 突变（⛔ 改动一律不进树；`-overlay` 只改编译视图）

做法：把被改文件复制到 `.scratch/wisp/probes/219/mut/ac15/<发>/`，在副本上改，用 `go test -overlay=<发>/overlay.json` 只影响那一发编译。**盘上源文件一字节未动**——每发前后同一把尺：

```
git status --porcelain -- cmd/wisp      →  空
```

（每发运行时都在命令里先打了这一行；三发全部为空。三发跑完后 `git status --porcelain -- cmd internal` 也是**空**。）

### 3.1 M1 —— 摘掉「允许一次」真正花掉授权那一行 ⇒ 红的是**行为**，不是句子

改：`cmd/wisp/resident_approval_windows.go:340-341`

```go
case trayAnswerAllowOnce:
    sentence, err = surface.allow(card.CorrelationID)     // 原文
```
改成（副本里）用**逐字相同**的句子返回、但**不调** `surface.allow`。也就是说：句子层看不出来，只有"被卡的那一笔写有没有落盘"能看出来。

```
go test ./cmd/wisp -overlay=".scratch/wisp/probes/219/mut/ac15/m1/overlay.json" \
  -run '^TestAC2TrayAllowOnceSpendsTheNativeGrant$' -count=1 -v
```

| | 读数 |
|---|---|
| 基线（未突变） | `--- PASS: TestAC2TrayAllowOnceSpendsTheNativeGrant (2.81s)` / `ok … 2.956s` |
| M1 | `resident_tray_answer_228_windows_test.go:124: a call the tray allowed must execute, got: 审批超时（40 秒未确认），C18 一律判拒绝，已自动拒绝` / `--- FAIL (44.37s)` / `FAIL` |

⇒ 红的是**行为牙**（卡片没人答 ⇒ 40 秒后 C18 过期 ⇒ 写没落）。**不是字面钉**：句子断言（`:107`）在突变下仍然通过。

### 3.2 M2 —— 摘掉"操作者打的理由"在控制台那一跳 ⇒ 两枚行为断言同时红

改：`cmd/wisp/approval_reply.go:672-673`

```go
case "no":
    return s.reject(corr, sanitizeReplyText(arg))          // 原文
```
改成 `return s.reject(corr, "")`（理由被丢掉）。

```
go test ./cmd/wisp -overlay=".scratch/wisp/probes/219/mut/ac15/m2/overlay.json" \
  -run '^TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel$' -count=1 -v
```

| | 读数 |
|---|---|
| 基线 | `--- PASS: TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel (3.01s)` / `ok … 3.281s` |
| M2 | `approval_reply_201_test.go:321: the model-visible rejection text lost the operator's reason` ＋ `:336: the audit never booked the reason the operator typed` / `--- FAIL (2.38s)` |

⇒ 一发同时打掉 **AC#5 去向①（审计）与②（模型）**，两枚都是行为牙。这个测试**不是**断言源码字符串，读的是工具结果正文与审计正文。

### 3.3 M3 —— 把「长期允许」误接到一次性 `allow` ⇒ 红的是**路由指纹**

改：`cmd/wisp/resident_approval_windows.go:344-345`，`surface.always(...)` → `surface.allow(...)`。

```
go test ./cmd/wisp -overlay=".scratch/wisp/probes/219/mut/ac15/m3/overlay.json" \
  -run '^TestAC2TrayAllowAlwaysRoutesIntoTheWideningVerb$' -count=1 -v
```

| | 读数 |
|---|---|
| 基线 | `ok`（同一枚测试在 §2.1 的 PASS 之外单独复跑过） |
| M3 | `resident_tray_answer_228_windows_test.go:222: the 长期允许 item did not reach replySurface.always; booked lines: …` / `--- FAIL (0.00s)` |

⇒ 行为牙（审计里那条 `route="native/always"` 指纹不出现）。

### 3.4 探针（不是突变）：超长理由到底怎么被截断

用 `-overlay` **注入一个盘上不存在的测试文件**（`cmd/wisp/zz_probe219_reason_bound_test.go` 只存在于 `.scratch` 副本），所以 `cmd/wisp` 零写：

```
go test ./cmd/wisp -overlay=".scratch/wisp/probes/219/mut/ac15/probe/overlay-add.json" \
  -run '^TestProbe219ReasonBound$' -count=1 -v
```

读数（逐字）：

```
PROBE219 reason_bound: input_runes=200 output_runes=160 replyReasonMax=160
PROBE219 reason_bound: truncation_marker_in_output="" (empty means SILENT)
PROBE219 reason_bound: output_tail="字字字"
PROBE219 reason_bound: empty="" blank=""
PROBE219 reason_bound: n=159 -> 159 runes
PROBE219 reason_bound: n=160 -> 160 runes
PROBE219 reason_bound: n=161 -> 160 runes
--- PASS: TestProbe219ReasonBound (0.00s)
```

⇒ **静默截断**：200 字进、160 字出，输出里**没有任何截断标记**（`…`／`...`／`截断`／`省略`／`+` 全 0 命中）。这就是 AC#5「不许静默截断」那一句在盘上的现状。

### 3.5 三发突变的判语（"装饰性"那一栏）

| 发 | 摘掉的行为 | 红的 | 分类 |
|---|---|---|---|
| M1 | 托盘「允许一次」真的花掉授权 | `TestAC2TrayAllowOnceSpendsTheNativeGrant:124` | **行为牙** |
| M2 | 控制台把操作者的理由带下去 | `TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel:321/:336` | **行为牙** |
| M3 | 「长期允许」走长期那条路 | `TestAC2TrayAllowAlwaysRoutesIntoTheWideningVerb:222` | **行为牙** |

**三发没有一发是"只让字面钉红"** ⇒ 这三处接线**不是装饰性**。

⚠ 同时**具名**一处**字面钉**（它没错，只是要知道它是哪种尺）：`cmd/wisp/resident_tray_answer_228_windows_test.go:315-390` 解析 `internal/ball/*.go` 源码里的 id／标签／回调名（`menuAllowOnce: "允许一次"` 等）——改文案会红、改行为不会。它管的是**文案与 id 不许漂**，不是"按钮能批"。

---

## 4. 「界面点不动」今天还成不成立（现读复核）

**判语：是，仍成立**——但**成立的范围要说清**：它是"**网页那侧的卡片点了，Go 侧没有人接**"，不是"托盘按钮点不动"。

| # | 读数 | 命令原文 | 结果 |
|---|---|---|---|
| 1 | 网页的 L2 卡**有**一枚「拒绝」按钮，点下去发的是 `panel.approval.request` | `Select-String frontend/src/components/l2-approval-card.tsx -Pattern "onClick"` | `:342 onClick={() => send("refuse")}`；`send` 在 `:263-267` 调 `requestApprovalResolution(view.correlationId, outcome)` |
| 2 | 那条请求的方法名只在这里产出 | `Select-String frontend/src/lib/panel.ts -Pattern "panel.approval.request"` | `panel.ts:618`（`requestApprovalResolution`）**唯一一处** |
| 3 | **Go 侧没有这条路由**：面板入向是闭名单四枚，不含它 | `Select-String internal/panel/inbound_roster_253_test.go -Pattern 'panel\.'` | `:83 panel.attachment.add`／`:84 panel.message.send`／`:85 panel.mode.request`／`:86 panel.workspace.request` |
| 4 | 而且这条方法**被门当场拒**（不是"没实现"，是"禁"） | `go test ./internal/panel -run '^TestGrantWireShapesAreRefusedAtTheDoor$'` | **PASS 0.00s**：`l2_grant_boundary_test.go:1960-1968` 把 `panel.approval.request`＋`outcome:"grant"` 判为必须被 `ParseComposerRequest` 拒 |
| 5 | 网页上**零** `approval.decide`（ban #6 没被踩） | 递归 grep `frontend/src`＋`frontend/scripts` | **0 命中** |
| 6 | 网页的卡上**没有理由框** | `Select-String frontend/src/components/l2-approval-card.tsx -Pattern "<input|<textarea|placeholder"` | **0 命中** |
| 7 | 面板路线唯一的**生产**调用者是控制台动词，不是页面 | `grep -n "PanelReject" cmd/wisp/approval_reply.go` | `:331`，其调用者 `:677 case "panel-no"`（控制台读行） |
| 8 | 卡片自己就这么印 | `frontend/src/components/l2-approval-card.tsx:329` | 「批准权在原生侧 · 面板不提供允许按钮」 |

⇒ **"网页→Go 那一跳今天只有命令行路"（票 215 `:46` 对 `A410` 的引述）在 219 这一格上依然成立**：页面上的「拒绝」按下去，方法名进不了那扇门（第 3、4 条），而页面上根本没有理由框（第 6 条）。**故 AC 里"网页那侧送进来的 allow"这个形状在今天的产品里连"送得进来"都还谈不上——AC#3 的成立方式因此是"路由上被拒"，不是"页面能发而被拦"。**

⚠ 这**不影响** AC#1：三枚答复按钮长在**原生侧**（托盘），那条链是通的（§2.1、§3.1／§3.3）。

---

## 5. 判语汇总与未做/未跑清单

| 项 | 判语 |
|---|---|
| AC#1 | **部分**（三枚按钮＝允许一次/拒绝/长期允许 各有行为牙；票面第三枚"本次会话内"只有控制台动词；"不再弹卡"只对 L1；双击启动形状无执行者） |
| AC#2 | **满足**（四字段审计＋逐字节撤销往返，两枚本腿实跑 PASS） |
| AC#3 | **满足**（路由上拒＋烧授权＋仍可拒，本腿实跑 PASS） |
| AC#4 | **满足**（五个种植样本全拒，本腿实跑 PASS） |
| AC#5 | **部分**（①审计 ②模型 有牙；③用户回看无落点；④空理由零测试；⑤**超长静默截断**，实测 200→160 无标记） |

**未跑（具名，不拿"全绿"顶掉）**
- 整包 `go test ./cmd/wisp`：**未跑**（禁令：`cmd/wisp` 写面归写腿 452b，同一时刻只许单枚）。
- `TestAC2TrayRejectRefusesOnTheNativeRoute`：〔未突变〕（理由见 §2.1②）。
- `Test219AC2SecondRevokeOfTheSameRuleIsRefused`、`Test219AC2MutationRevokeWriteMustPersistForBytesToComeBack`、`TestTicket224ProductionSessionDoesNotSurviveRestart`：〔未复跑〕。
- 端到端"真窗口点托盘菜单"：**不可能在无头环境跑**（真窗口＝winlive rig，本腿禁令），故三枚按钮的**手势→执行者**那一跳是靠 `resident_windows.go:571` 的现读＋`resident_tray_answer_228_windows_test.go:492` 的源码齿，加上本腿把**执行者以内**每一跳都当真行为突变过（M1／M3）。

**撤销口令**（本腿零写产码，所以没有"复原产码"这一步；下面是把本件及其探针件当临时件处理的口令）：
> 「撤 219-ac1-ac5 收件」＝ 删 `docs/evidence/s1/219-ac1-ac5-receipt.md` 与 `.scratch/wisp/probes/219/mut/ac15/`（**两个都是只建不删的默认项，要删得机主点名**）。

**本腿写面（就是这两处，别的一字节没动）**
- `docs/evidence/s1/219-ac1-ac5-receipt.md`（本件）
- `.scratch/wisp/probes/219/mut/ac15/**`（三发突变副本＋两枚 overlay.json＋一枚探针件＋探针 overlay）

⛔ 未改：产码／票面／台账／阈值／契约／`internal/panel/**`／`frontend/**`；⛔ 未 push；⛔ 未 commit 到产码路径（本件由编排者按收件规矩处置）。
