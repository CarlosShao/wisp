# 242-corrcensus-1 · ② 类型册与读者名册

锚：HEAD `8dbae00`（分支 dev）。扫描尺与 rc：
- `git grep -nE 'CorrelationID' HEAD -- 'internal/*.go' 'cmd/*.go' ':!*_test.go'` rc=0，**115 行命中**（射程＝internal 与 cmd 下的非测试 `.go`）。
- 下文所有 `文件:行` 与该行短语均取自同一条 `git grep -n` / `git show | sed -n` 发现。

---

## 1. 类型与字段册（按包归组）

### agent 包
| 类型.字段 | 定义 文件:行 | 备注 |
|---|---|---|
| `agent.ToolRequest.CorrelationID` | internal/agent/tools.go:56 | loop 构造请求时填（loop.go:654）；bridge.Execute 首行消费（bridge.go:269-270） |
| `taskJournal.corrID`（非导出） | internal/agent/journal.go:53 | 构造于 journal.go:59-63；**loop.go:363** `newTaskJournal(l.opt.Journal, taskID, taskID) // C18: correlation == task id`；写行于 journal.go:81 |

### tools 包
| 类型.字段 | 定义 文件:行 | 备注 |
|---|---|---|
| `tools.Decision.CorrelationID` | internal/tools/gate.go:47-49 | 注释逐字："CorrelationID routes the decision to the right approval card (C18); the same value lands in the tool_call row." |
| `cancelHandle.corr`（非导出） | internal/tools/cancel.go:48-58（访问器 `CorrelationID(ctx)` 在 55） | 注释："the id of the call the running tool belongs to" |

### approval 包
| 类型.字段 | 定义 文件:行 | 备注 |
|---|---|---|
| `Veto.CorrelationID` | internal/agent/approval/approval.go:70-75 | 注释："CorrelationID routing is a C18 requirement" |
| `Prompt.CorrelationID` | internal/agent/approval/ui.go:28 | 注入 UI 的确认面 |
| `PanelItem.CorrelationID` | internal/agent/approval/ui.go:51 | 面板可见投影（无 grant、无 allow 字段） |
| `Event.CorrelationID` | internal/agent/approval/ui.go:73 | transient 事件 |
| `Request.CorrelationID` | internal/agent/approval/gate.go:716 | 答复路由输入 |
| `ReplyCard.CorrelationID` | internal/agent/approval/replies.go:71 | 已展示卡台账 |
| `LiveApproval.CorrelationID` | internal/agent/approval/pending_read.go:56 | 进程内读模型 |
| `L1Window.CorrelationID` | internal/agent/approval/window_read.go:44 | 读模型；注释："the id the window is keyed by in g.windows" |
| `CancellationReport.CorrelationID` | internal/agent/approval/report.go:22 | D31 取消报告 |
| `qitem.Corr` / `qitem.names`（非导出） | internal/agent/approval/queue.go:33 / :42 | 队列内部地址与别名集合（36-41 注释） |

### panel 包
| 类型.字段 | 定义 文件:行 | 备注 |
|---|---|---|
| `ApprovalSubject.CorrelationID` | internal/panel/approval.go:28 | 注释："The agent loop owns the correlationId (C17)" |
| `ApprovalCardView.CorrelationID` `json:"correlationId"` | internal/panel/approval.go:40 | 渲染投影 |
| `ResultChunk.CorrelationID` | internal/panel/composer.go:96 | 流式块 |
| `ModeRequest.CorrelationID` | internal/panel/composer.go:178-179 | 注释："ties the request to whatever card the gate raises" |
| `NativeVerdict.CorrelationID` | internal/panel/pump.go:78 | 快照输入 |
| `L1WindowWait.CorrelationID` | internal/panel/pump.go:128 | 快照输入 |

### memory 包
| 类型.字段 | 定义 文件:行 | 备注 |
|---|---|---|
| `ToolCall.CorrelationID` | internal/memory/models.go:87（注释 `// C18`） | 167-168：空值拒写（"tool_call.correlation_id is required (C18)"） |

cmd/wisp 侧无自有类型，全部读 approval/panel 的值。

---

## 2. 读者名册：决定行为（29 处；按值分支/查找/发行键）

| # | 文件:行 | 作用（一句话） |
|---|---|---|
| D1 | internal/tools/bridge.go:269-270 | `if req.CorrelationID == "" { req.CorrelationID = req.TaskID }`：空值回退，决定本调用后续所有 corr 读键的取值 |
| D2 | internal/tools/bridge.go:525 | `withCancel(..., cancelHandle{corr: orDefault(req.CorrelationID, req.TaskID)})`：D31 取消记账键，决定取消/否决按哪个键查找在跑调用 |
| D3 | internal/tools/bridge.go:532 | `defer b.cancel.Complete(orDefault(...))`：同键释放，决定 in-flight 记账何时关闭 |
| D4 | internal/tools/bridge.go:596 | `b.cancel.Vetoed(orDefault(dec.CorrelationID, dec.TaskID))`：决定 error 结果是否改判 cancelled |
| D5 | internal/tools/bridge.go:778-779 | `corr := orDefault(...)`；`Vetoed(corr)`/`Report`：决定是否附取消说明文本（Report 内 `g.running[corr]` 再查一次） |
| D6 | internal/agent/approval/gate.go:281-284 | L1 窗口键 `orDefaultText(d.CorrelationID, d.TaskID)`＋`openWindow` 查重（同键已有窗口 ⇒ "同一 correlation_id 已有确认窗口在跑…fail-closed 拒绝"）：决定窗口能否开、否决地址是什么 |
| D7 | internal/agent/approval/gate.go:423-427,445 | Veto：非空校验（空 ⇒ ErrUnknownCorrelation）＋`g.windows[v.CorrelationID]`、`g.running[v.CorrelationID]` 按值查找：决定否决命中哪个窗口/在跑调用 |
| D8 | internal/agent/approval/gate.go:476 | 未命中窗口/在跑时 `g.q.reject(v.CorrelationID, ...)`：决定哪张 L2 卡被否（队列按 byID 含别名解析） |
| D9 | internal/agent/approval/gate.go:516,521-522,552 | `corr := it.Corr` 后 `incoming := orDefaultText(...)`、`d.CorrelationID = corr`（重戳为队列地址）；`markStarted(incoming, ...)` 用入队前值记 D31：决定"桥记账键"与"卡片地址"两个可能不同的串 |
| D10 | internal/agent/approval/gate.go:727-729 | `g.q.reject/allow(r.CorrelationID, ...)`：原生答复按值路由到队列项 |
| D11 | internal/agent/approval/gate.go:744 | `g.q.revokeGrants(r.CorrelationID)`（面板"允许"被拒时烧令牌）：决定哪张卡的 grant 作废 |
| D12 | internal/agent/approval/gate.go:748 | 面板拒绝 `g.q.reject(r.CorrelationID, ...)`：按值路由 |
| D13 | internal/agent/approval/queue.go:163-168 | 入队主键：空 ⇒ `"approval-"+seq`；撞名 ⇒ `原值#seq`：决定卡片第一地址 |
| D14 | internal/agent/approval/queue.go:186-198（循环在 191） | `otherNames=[d.CorrelationID, d.TaskID]` 减主键去重：决定卡片还能被哪些串合法寻址（否决方向别名；Ticket 87 语义） |
| D15 | internal/agent/approval/queue.go:614-618 | 重放发行新键 `it.Corr+"#replay"`（撞名再加 `#seq`）：决定重放卡地址 |
| D16 | internal/agent/approval/replies.go:171-194 | 非空门（171）＋`byID[p.CorrelationID]` 台账键＋order 去重（186-194）：决定哪些已展示卡可被 Look/答复（缺 ⇒ 无 grant 可花） |
| D17 | internal/agent/approval/replies.go:329 | `DecideFromNative(Request{CorrelationID: corr, ...})`：允许答复的路由值 |
| D18 | internal/agent/approval/replies.go:405 | `decide()` 拒绝 `Request{CorrelationID: corr}`：路由值 |
| D19 | internal/agent/approval/replies.go:438 | `PanelAllow` `Request{corr}`（该路线必被拒＋烧 grant 路径）：路由值 |
| D20 | internal/agent/approval/replies.go:463 | `g.Veto(Veto{CorrelationID: corr, ...})`：决定否决落在窗口还是卡 |
| D21 | internal/tools/subagent_197.go:259-263 | `parentID := CorrelationID(ctx)`；空 ⇒ fail-closed 拒派生；`Roster.Look(parentID)` 查"派生者是否已是子代理"（深度上限）：决定父子登记键 |
| D22 | internal/tools/task.go:706-709 及后 | `caller := CorrelationID(ctx)`：`target==caller` 自停拒绝、`rec.ParentTaskID != caller` 兄弟/旁支拒绝：决定"谁在停"的身份 |
| D23 | internal/panel/subagent_roster_197.go:214,219-220 | `indexWaitingKey(waiting, card.CorrelationID / w.CorrelationID / w.TaskID)`：决定哪些名册行被标"在等"（join 键） |
| D24 | cmd/wisp/resident_approval_windows.go:613 | `ra.cards.Veto(card.CorrelationID)`（Esc 否决路径）：决定命中哪张常驻卡 |
| D25 | cmd/wisp/resident_approval_windows.go:729 | `ra.cards.Reject(ctx, card.CorrelationID, ...)`：退出序列按值拒绝 L2 卡 |
| D26 | cmd/wisp/resident_approval_windows.go:744,933 | `ra.cards.Forget(card.CorrelationID / e.CorrelationID)`：按值从答复台账摘卡 |
| D27 | cmd/wisp/run.go:1402 | `u.live.forget(e.CorrelationID)`：run 侧 live 记账按值冻结 |
| D28 | internal/memory/models.go:167-168 | 空值 ⇒ `tool_call.correlation_id is required (C18)` 拒写：入库校验分支 |
| D29 | internal/agent/approval/report.go:95-104 | `corr := orDefaultText(...)`；`g.running[corr]` 查找决定 `StartedBeforeVeto`/`Vetoed`/`Channel` 字段 |

## 3. 读者名册：只进审计/日志（13 处）

| # | 文件:行 | 作用（一句话） |
|---|---|---|
| A1 | internal/agent/journal.go:81 | `CorrelationID: t.corrID` 写入 `tool_call` 行（取证） |
| A2 | internal/memory/dao_toolcall.go:30 | INSERT 参数（`tool_call.correlation_id`） |
| A3 | internal/memory/dao_toolcall.go:168 | SELECT 扫描回读该列 |
| A4 | internal/tools/bridge.go:1077 | `b.log("tools: call ... corr=%s ...")` 行 |
| A5 | internal/agent/approval/gate.go:725 | `logf("approval: native route corr=%s ...")` |
| A6 | internal/agent/approval/gate.go:739 | `logf("approval: PANEL-ALLOW-REJECTED corr=%s ...")` |
| A7 | cmd/wisp/approval_always.go:182-183 | `rt.auditf("... corr=%q ...")` 审计行（WAITING-STATE） |
| A8 | cmd/wisp/resident_approval_windows.go:732 | `residentAuditf("...拒绝对待卡片 %s 失败：%v", card.CorrelationID, err)` |
| A9 | cmd/wisp/resident_approval_windows.go:742 | `residentAuditf("... corr=%s ...")`（RESIDENT-WINDOW-ABANDONED） |
| A10 | cmd/wisp/resident_approval_windows.go:869 | `slog` 属性 `"corr", p.CorrelationID`（卡片无处呈现 fail-closed 句） |
| A11 | cmd/wisp/resident_approval_windows.go:882 | `slog.Info("... 显示一张确认卡片", "corr", ...)` |
| A12 | cmd/wisp/resident_approval_windows.go:920 | `slog.Warn("... 取消键未借到 ...", "corr", ...)` |
| A13 | cmd/wisp/resident_approval_windows.go:936 | `slog.Warn("... 醒目提示", "corr", ...)` |

## 4. 读者名册：只做标识（45 处 ＝ 型字段 18 ＋ 誊抄/读模型/展示 27）

型字段 18（定义处，上表 §1 已给）：agent/tools.go:56；tools/gate.go:49；approval/approval.go:75；approval/ui.go:28、51、73；gate.go:716；replies.go:71；pending_read.go:56；window_read.go:44；report.go:22；panel/approval.go:28、40；panel/composer.go:96、179；panel/pump.go:78、128；memory/models.go:87。

誊抄/读模型/展示 27：
- internal/agent/approval/queue.go:33（qitem.Corr 字段）、:572（`PanelItem{CorrelationID: it.Corr, ...}` 读模型誊抄）
- internal/agent/approval/pending_read.go:125（`CorrelationID: it.Corr` 誊抄进 LiveApproval）
- internal/agent/approval/window_read.go:73（誊抄进 L1Window）、:76（`sort.Slice` 按 CorrelationID 排序；注释称"ordered by correlation id so ... reproducible bytes"——排序键承载哈希日志字节复现）
- internal/agent/approval/report.go:97（誊抄进 CancellationReport）
- internal/agent/approval/gate.go:308、322、341、437、554、561、573（Event 标签誊抄）、:599（`Prompt{CorrelationID: corr, ...}` 誊抄）
- internal/agent/approval/replies.go:175（`ReplyCard{CorrelationID: p.CorrelationID, ...}` 誊抄）
- internal/panel/approval.go:78（`ApprovalCardView{CorrelationID: subject.CorrelationID, ...}` 誊抄）
- internal/panel/pump.go:101（NativeVerdict→ApprovalSubject 誊抄）、:526、:549（`ResultChunk{CorrelationID: key, ...}`；归组键是调用方传入的 `key`，本字段是誊写）
- cmd/wisp/panel_pump.go:66（队列视图→NativeVerdict 誊抄）、:338（`ids = append(ids, c.CorrelationID)` 横幅短名单）
- cmd/wisp/run.go:1381（`fmt.Fprintf(u.out, "  卡片编号：%s\n", p.CorrelationID)`）
- cmd/wisp/approval_reply.go:402（面板投影句子 `"卡片 %s｜..."`）
- cmd/wisp/resident_approval_windows.go:594（vetoDoneLine 句子）、:618（vetoRejectedLine 句子）、:886、:922（fmt.Printf 用户句子）

## 5. 判定口径与未展开项

- 分类判据：**决定行为**＝按该值做分支/查找/作为发行键或释放键；**只进审计/日志**＝值只落日志或 DB 行；**只做标识**＝类型定义、读模型/投影誊抄、用户可见句子。同一条流程里"发路由值"的构造点（replies.go:329/405/438、gate.go:727-729 等）按其在路由链上的一环计入"决定行为"。
- 未展开（如实登记）：`q.indexLocked` 函数体未逐行读（queue.go:181 为调用点现量）；`task.spawn` 侧 `ParentTaskID` 的填值点未定位（仅据 task.go:704 注释）；80 枚含 `CorrelationID` 的测试文件（62 仓内 ＋ 18 枚 `.scratch` 拷贝，`git grep -ln` rc=0）未逐枚判读。
