# 377-a1 普查：读档位的 11 个点 —— 4 个读配置档、7 个读有效档；「配置 ≠ 有效」只在 1 条 seam 上成立（只读腿）

- 腿 `377-a1`（只读＝`AC#0`）。读数时 HEAD **`ccf17a3b314eb936eba5e4b0d3f1c5b01e9894d6`**（2026-10-24 12:45 +08）。本笔**未落任何 `.go`**、未改票面／台账／`Status:`、未 push、未翻框。
- ★**读数归属 HEAD，⛔ 工作树**：工作树里 `cmd/wisp/run.go`、`internal/memory/*`、`internal/winsec/*` 是**在飞写腿**（373-r1／369-r1）⇒ 一律 `git archive HEAD` 导出到仓外 `D:\tmp\wisp\20261024\377a1\head\`（`internal`＋`cmd`＋`tools`＋`frontend`＋`docs`），**grep／构建／测试／探针全在导出树里跑**；行号一律 `git show HEAD:`／`git grep HEAD` 现取（★`Get-Content` 读 `run.go` 会**少数一行**，本件所有 `run.go` 号以 `git show` 为准）。仓内零字节改动，无工作树突变窗口。
- ★**探针**（**仓外**）：`head\internal\tools\zzz377probe_test.go`（3 用例）＋两枚变体树 `head-mut1`／`head-mut2`（各只改 `internal/perm/store.go` 一行，见 §⑤）。原始输出 `probe1.txt`／`probe2.txt`／`probe-head-mut1.txt`／`probe-head-mut2.txt`／`g-*.txt`。
- ★**尺的射程（落笔前五问）**：①**对象**＝「读权限档位」的**调用点**，⛔ 不是「提到 `PermissionMode` 的行」；②**形状**＝`file:line` ＋ 读哪个面 ＋ 用途；③**数的是代码**：`PermissionMode|configuredMode` 在 `internal/**`＋`cmd/**` 原始命中 **121** 行＝`_test.go` **70** ＋ 非测试 **51**＝**注释 20 ＋ 代码 31**（逐条分类、非抽样）；④**总数核对**：31 代码行＝定义／声明／常量／校验／写侧 **20** ＋ 读配置**字符串字段** 2（`config/manager.go:453-454`）＋ **真读档点 9**；再加票面那把尺抓不到的 **2** 点（`pump.go:260` 的 `.Mode()`、`config_reload.go:282` 的 `rt.permissionMode()`）⇒ **真读档点 11**；⑤**谁算的**＝`git grep`／`risk.Mode.String()` 现渲染／`go1.27.1`。
- ★★**`.Mode()` 这把尺的射程（⛔ 直接数）**：`\.Mode\(\)` 全仓 **78** 行、非测试 **14** 行，其中**只有 1 行**是权限档（`internal/panel/pump.go:260` `p.src.Mode()`），其余 13 行是 `os.FileInfo.Mode()`／`os.FileMode`（`models/manifest.go:346`、`tools/fs.go:263`、`winsec/winsec_other.go:219,263,300,455`…）⇒ **拿 `.Mode()` 当尺会把 os.FileMode 数成权限档（14 倍高估）。**

## ① 读档位的**全部**点（11 条）

两个读面（票 372 落的）：`perm.Store.PermissionMode()`＝**有效档**（`internal/perm/store.go:196-205`：`s.confirm == nil && configured != ask_every_step ⇒ ask_every_step`）；`perm.Store.configuredMode()`＝**配置档**（`store.go:215-224`，理由 `:207-214`）；`config.Config.PermissionMode()`＝**配置档**（`internal/config/permmode.go:36-45`，函数内**无**降级分支）。

| # | `file:line` | 读的哪个面 | 用途 | 备注 |
|---|---|---|---|---|
| 1 | `cmd/wisp/approval_always.go:155` | **配置档**（`rt.mgr.Config().PermissionMode()`） | 记账字段：加宽卡的 `Decision.Mode`（`:122`） | ★票面点名处；见 §④ |
| 2 | `cmd/wisp/config_reload.go:282` | **配置档**（同一个 `rt.permissionMode()`） | 记账字段：D36 放宽卡的 `Decision.Mode`（`:291`） | ★与 #1 **同一函数体** |
| 3 | `internal/perm/store.go:158` | **配置档**（`s.configuredMode()`） | 记账：启动 `Switch{From,To}`＝配置档、`ResultNoChange` | ★只进内存 `hist`（`:159` 用 `record`，**不落 log**） |
| 4 | `internal/perm/store.go:254` | **配置档**（`s.configuredMode()`） | **门判＋记账**：`from == to` 判 no-change；`MODE-SWITCH from/to` 落审计 | ★**故意**，`store.go:207-214` 写明；见 §⑤ |
| 5 | `internal/perm/store.go:167` | **有效档**（`s.PermissionMode()`） | 记账：`effective != configured` ⇒ 印 `MODE-DOWNGRADE` | ★372 的降级行 |
| 6 | `internal/perm/store.go:331` | **有效档**（`Snapshot.Mode`） | 展示：ticket 92 的只读视图 | ★**生产 0 个调用者**（`.Snapshot()` 3 处命中是 `panel.Pump`／`observe`，别的 receiver） |
| 7 | `internal/tools/mode.go:46` | **有效档**（`b.modes.PermissionMode()`） | **门判**（每发一次，`bridge.go:310-312`）＋记账（`dec.Mode`）＋`MODE-SILENCE`／`MODE-REDLINE` 审计（`bridge.go:320,333,337`） | ★唯一进决策链的读 |
| 8 | `internal/panel/composer_handlers.go:134` | **有效档**（`h.Modes.PermissionMode()`） | **门判**（方向／目标）＋记账（`panel: MODE-REFUSED from=`） | ★`ModeWriter` 注释 `:87-89` 明写「要有效档」 |
| 9 | `internal/panel/pump.go:260` | **有效档**（`p.src.Mode()`，读数源在 `cmd/wisp/run.go:725`） | 展示：`composer.mode.current` | ★`PumpSources.Mode` 注释 `pump.go:141` |
| 10 | `cmd/wisp/run.go:711` | **有效档**（`modeStore.PermissionMode()`） | 记账：`perm: MODE-READ origin=startup mode=`（`:710`） | ★与 #3 相反：这条**进 log** |
| 11 | `cmd/wisp/run.go:725` | **有效档**（方法值 `rt.modes.PermissionMode`） | 展示：塞进 `panel.PumpSources.Mode` | ★读发生在第 9 点里 |

**携带而非读取**（⛔ 重复计数）：`approval_always.go:122` `Mode: s.rt.permissionMode()`；`cmd/wisp/run.go:919` `Mode: sw.From`（`Set` 的 `from` 原样带走）。

## ② 「配置档 ≠ 有效档」的**全部**情形 ＋ 两发对照

判据只有一处：`store.go:201` `if s.confirm == nil && configured != risk.ModeAskEveryStep` ⇒ **恰两种配置档会降**：`ask_high_risk`、`auto_approve`；`ask_every_step` 不降；`s==nil`／`mgr==nil`／`Config()==nil`／不可解析 四路**两个面都答 `ask_every_step`**（`store.go:197-204` vs `:216-223`、`permmode.go:37-44`）⇒ 那四路**无差异**。

六格现跑（`go test ./internal/tools -run TestZZZ377ProbeTwoFaces`，导出树，真 `config.toml` 经 `SaveFile`＋`NewManager`）：

```
CELL config=ask_every_step  confirm=nil | effective=ask_every_step configured=ask_every_step
CELL config=ask_every_step  confirm=leg | effective=ask_every_step configured=ask_every_step
CELL config=ask_high_risk   confirm=nil | effective=ask_every_step configured=ask_high_risk   ← 差 1 档
CELL config=ask_high_risk   confirm=leg | effective=ask_high_risk  configured=ask_high_risk
CELL config=auto_approve    confirm=nil | effective=ask_every_step configured=auto_approve    ← 差 2 档
CELL config=auto_approve    confirm=leg | effective=auto_approve   configured=auto_approve
```
降级同时落一行（`confirm=nil` 两格各一条）：`perm: MODE-DOWNGRADE from=auto_approve to=ask_every_step origin="startup" actor="process" result=downgraded-no-channel detail="本进程没有接入确认通道：config.toml 仍是 auto_approve，有效档位降为 ask_every_step（372 读侧降级，配置值不动）"`。

★**两发对照（配置 `auto_approve` ＋一次 L1 工具调用）**——只差「有没有腿」：

```
confirm=nil | Decision.Mode=ask_every_step ModeSilenced=false | gate windows=1 | MODE-SILENCE 行=无
confirm=leg | Decision.Mode=auto_approve   ModeSilenced=true  | gate windows=0 | MODE-SILENCE mode=auto_approve assessed=L1 effective=L0
```
⇒ 同一份 `config.toml`，无腿那一发**真弹了 1 个 L1 窗**（`windows=1`）并执行；有腿那一发 **0 窗**、被静默。差异是**实的**，不是文案差异。

## ③ 「记账面报错档」的实际后果（那一发）

造法（`TestZZZ377ProbeSetAudit`；装配＝`wisp panel-inbound` 的形状：`config=auto_approve`＋`Confirm=nil`＋一次切档）：

```
SET config=auto_approve confirm=nil to=ask_high_risk | effective BEFORE=ask_every_step | Set err=<nil> | config AFTER=ask_high_risk | effective AFTER=ask_every_step
  AUDIT perm: MODE-SWITCH from=auto_approve to=ask_high_risk ... origin="panel" actor="probe" result=applied
```

- **审计行写的是**：`from=auto_approve to=ask_high_risk result=applied`。**实际执行时按的是**：切档前 `ask_every_step`、切档后**仍是** `ask_every_step`（文件动了、有效档没动）。
- ★**事后取证会得出的错误结论（具体）**：只有一行 `MODE-SWITCH from=auto_approve`，而 `store.go:354-357` 明写这族行**就是给人 grep 的** ⇒ 取证者读出「那一刻宿主在**全自动**，操作员把它收到 `ask_high_risk`」。**两条都错**：(a) 该进程**从未**跑过 `auto_approve`（连一张静默卡都没放过）；(b) 收完之后**仍然是** `ask_every_step`，`to=ask_high_risk` 只是文件值。而 `tools.Decision.Mode` 的字段注释（`internal/tools/gate.go:38-41`）说该字段存在的**唯一目的**是回答「why was I not asked?」——错的 `from` 把**那个问题**答反：真答案是「问了、每一发都问了」，记出来的是「没问，因为当时是全自动」。
- ★**同一发里另一个面是诚实的**：`composer_handlers.go:134` 的 `from` 读**有效档**（面板若拒，`MODE-REFUSED from="ask_every_step"`）⇒ 同一条 seam 上**两个 `from` 不同源**，且都自称「当前档」。

## ④ `approval_always.go:151-156` 逐字 ＋ `Decision.Mode` 的消费者穷举

```
146: // permissionMode reads the档 this run is screening under, so the widening card
147: // carries the mode the audit trail expects (tools.Decision.Mode's own comment:
148: // 「why was I not asked?」 has to be answerable from the record). An unreadable
149: // mode is booked as the strictest档 rather than the safest-sounding one, which is
150: // config.PermissionMode's own fail-closed reading.
151: func (rt *agentRuntime) permissionMode() risk.Mode {
152: 	if rt == nil || rt.mgr == nil {
153: 		return risk.DefaultMode()
154: 	}
155: 	return rt.mgr.Config().PermissionMode()
156: }
```
- ★**它读配置档**（`config.Config.PermissionMode()`，`permmode.go:36` 全函数**无**降级分支），⛔ 不是 `rt.modes.PermissionMode()`。**逐字反证**：同仓 `rt.modes` 就在 `run.go:683`，这函数却走 `rt.mgr`。写进 `:122` `Mode: s.rt.permissionMode()`（加宽卡的 `tools.Decision.Mode`）。
- ★★**`Decision.Mode` 消费者＝穷举（生产 0 个读者）**：
  1. **写入点 4 处**：`bridge.go:312`（有效档）、`approval_always.go:122`（配置档）、`config_reload.go:291`（配置档）、`run.go:919`（配置档，切档卡）。
  2. **存储**：`internal/agent/approval/queue.go:32` `qitem.Dec`（整结构体存下，从不读 `Mode`）。
  3. **`LiveApprovals()`**（`pending_read.go:113-130`）→ 唯一生产者 `cmd/wisp/panel_pump.go:66`，但 `:69-77` 那 7 个键只拷 `Tool/Args/Level/RulesHit/Reason/SessionOverrideBlocked`——**`Mode` 被丢掉**；`panel.CardViewFromDecision`（`internal/panel/approval.go:72`）吃的是 `risk.Decision`，**另一枚类型**。
  4. **`Gate.Replay`**（`gate.go:755`）→ `Queue.replay`（`queue.go:609`）把 `it.Dec` 原样带回 ⇒ `Mode` **会**跟着回放；但 `Gate.Replay` **生产 0 个调用者**（唯一调用在 `internal/agent/approval/queue_test.go:283`）。
  5. **审计／DB／卡面都不是消费者**：`tools: call` 行（`bridge.go:1174-1178`）**无 mode=**；`OnDecision`（`run.go:792-795`）**无 mode**；`tool_call` 表只有 `risk_level`／`decision`／`outcome`／`grant_id`（DDL `internal/memory/schema.go:81-96`，逐列核过）**无 mode 列**；`internal/agent/approval` 非测试代码里 **`Mode` 一词 0 命中**。
  ⇒ ★**净结论**：`Decision.Mode` 今天**写而无人读**。375-a1 说它「仅展示/记账」，实测比那更弱：**展示面也丢**。

## ⑤ 「改它会破什么」——两枚变体（仓外，各只改 `perm/store.go` 一行）

| 变体 | 改哪 | 现跑读数 | 判 |
|---|---|---|---|
| `m1` | `store.go:254` `from := s.configuredMode()` → `s.PermissionMode()`（把 `Set` 的记账改读有效档） | `to=ask_every_step`：`err=<nil>`、`result=no-change`、**`config AFTER=auto_approve`**（没落盘）；`to=ask_high_risk`：`from=ask_every_step` | ★**破**：操作员「把文件改严」的请求被静默吞掉——正是 `store.go:209-214` 预告那条，本笔**实测复现** |
| `m2` | `store.go:158` `configured := s.configuredMode()` → `s.PermissionMode()`（把启动记账改读有效档） | `MODE-DOWNGRADE` 行**整条消失**（`confirm=nil` 两格不再印）；只剩 `MODE-SWITCH from=auto_approve to=ask_every_step` | ★**破**：372 唯一的降级审计证据被删——判据 `effective != configured` 退化成 `effective != effective` |

★**依赖「记账＝配置档」的地方，逐处核**：
- `Set` 的 `from/to`（#4）——**依赖**，且它同时是**门判**（`from == to` 决定是否落盘）；改就破（`m1`）。
- `New` 的启动 `Switch`（#3）——**依赖**，它是 `MODE-DOWNGRADE` 的**对照组**；改就破（`m2`）。
- `MODE-DOWNGRADE` 的 `From`（`store.go:169`）——**必须**配置档，否则那行不成句（「config.toml 仍是 X」的 X）。
- **重放要用配置档恢复？——核过：不成立**。`Queue.replay` 回放的是 `it.Dec`（当初存下的判词），不读任何档；`Gate.Replay` 生产 0 调用者 ⇒ 这条依赖**不存在**；也无人依赖「`Decision.Mode` 是配置档」（它 0 读者，§④）。

## ⑥ 票 372 射程 (a) ⇒ **本票的射程是哪几条 seam**

372 结案逐字（`docs/evidence/s1/372-v1-verdict.md:121`）：「**件二只在 `Confirm==nil` 的装配里生效**：`wisp run`（`run.go:668-671`）永远非 nil ⇒ 该 seam 的 A/C 两格**与改前同**」。现跑核：HEAD 的同一段在 `run.go:669-672`（`confirm := s.modeConfirm; if confirm == nil { confirm = rt.confirmModeSwitch }`，`rt.confirmModeSwitch` 是**方法值** `run.go:893` ⇒ 恒非 nil）。**语义成立；行号 +1 漂（票面 668-671 ⇒ HEAD 669-672）。**

`perm.New` 全仓生产调用 **2 处**（非测试）⇒ seam 只有三条：

| seam | `perm.New` 落点 | `Confirm` | 配置≠有效可能？ | 该 seam 上的记账点 |
|---|---|---|---|---|
| `wisp run` | `run.go:673-677` | `rt.confirmModeSwitch` **恒非 nil** | ★**不可能**（两面恒等） | #1 #2 #3 #4 #5 #7 #10 #11 |
| `wisp panel-inbound` | `panel_inbound.go:236` → `:272-279`（`storeConfirmFromPanel(nil)`→`:347-349` 返 nil） | **nil** | ★**可能**（唯一一条） | #4 ＋ #8；**无桥、无 `Decision.Mode`** |
| 常驻面板 | `panel_resident_windows.go:233` → 同一 `:272` | `ra.confirmModeSwitch`（`resident_approval_windows.go:111`）**非 nil** | 不可能 | 同 `panel-inbound` 的 #4／#8 |

★★**本票真射程＝「`wisp panel-inbound` 这一条 seam 上 `perm.Set` 的 `MODE-SWITCH from=` 一个字段」**（外加 #6 `Snapshot()` 这个 0 调用者的死面）。票面点名的 `approval_always.go:151-156` 与 `config_reload.go:282` **都只在 `wisp run` 上**，那里两面恒等 ⇒ **今天报不出错档**：形状是真的，后果是**潜伏**的。⛔ 据此说「375-a1 的 ⑥ 立错了」——它立的**形状**对；隐含的「已在生产上报错档」不对。

## ⑦ `AC#1` 三形代价 ＋ 「只记一个会不会丢信息」

| 形 | 要改几处 | 会不会破别处 | 丢不丢信息 |
|---|---|---|---|
| 〔甲〕记账点一律读有效档 | **2 处代码**（`approval_always.go:155` 改走 `rt.modes`；`config_reload.go:282` 随之）；照字面「一律」再算 `store.go:158`＋`:254` 就是 **4 处** | ★**破**：`m1`（严档请求被吞）＋`m2`（`MODE-DOWNGRADE` 消失），**两枚实测**；且 `store.go:169` 的 `From` 不能动 ⇒ 甲无法「一律」 | ★**丢**：改完全仓**没有一行**说「文件里配的是什么」（除 config.toml 本身，而它可能已被改写） |
| 〔乙〕记账同记两个 | **`Decision` 加一字段＋5 处**（`gate.go:42` 声明、`bridge.go:312`、`approval_always.go:122`、`config_reload.go:291`、`run.go:919`）＋**审计 formatter 1 处**（`bridge.go:1174`） | 不破（`Set`／`New` 不动）；但**得同时补一个读者**，否则只是把「写而无人读」从一个字段变两个（§④） | 不丢 |
| 〔丙〕审计记两个、展示用有效档 | **审计 formatter 2 处**（`bridge.go:1174` 的 `tools: call` 加 `mode=`／`configured=`；`run.go:710-711` 的 `MODE-READ` 加 `configured=`）＋可选把 #1／#2 改走 `rt.modes`（1 处函数体） | 不破：`Set`／`New`／展示面（#6 #9 #11）**一格不动**，且 `run.go:725` 已是有效档（372-v1:114 逐字「同一代理在**显示**侧是诚实的…无需动」） | 不丢 |

★★**「审计要能回答『当时配的是什么、实际跑的是什么』」——只记一个会丢；但今天的真问题是**一个都没记**：`tools: call`（`bridge.go:1174-1178`）、`OnDecision`（`run.go:792-795`）、`tool_call` 表**三处最大取证面全无 mode** ⇒ 甲／乙／丙对这三处**都是「加」不是「换」**；把 `Decision.Mode` 的读面从配置档换成有效档，**不改变**今天任何一处能回答的问题（它 0 读者）。★**仓里已有〔乙〕的先例**：`perm: MODE-DOWNGRADE from=… to=…`（一行同印两个面）＋`MODE-READ mode=<有效档>`。★**主流**：`m4-survey.md:48` 的 openchamber `displayedPermissionMode`（「`safety` 无分类器时按 `ask` **显示**」）＝显示面用有效档；`m4-survey.md:85` 的 dsh 是「放宽照旧生效，落 `permission/preset` 事件」＝**记切档动作**。两者都支持〔丙〕，都不支持〔甲〕。

## ⑧ 票面前提新鲜度（逐条对 HEAD blob 现量）

| 票面说 | HEAD 实测 | 判 |
|---|---|---|
| `approval_always.go:151-156` 读配置档写 `Decision.Mode` | `:151-156` 逐字＝`permissionMode()`，`:155` 走 `rt.mgr.Config().PermissionMode()` | ★**成立**（行号未漂） |
| 372 射程 (a)：`run.go:668-671` 永远非 nil | HEAD 在 `run.go:669-672`，方法值兜底 | ★**成立，行号 +1 漂** |
| 372 落了 `configuredMode()` 保配置档 | `store.go:215-224` 在；注释 `:207-214` 已预告 §⑤ 的 `m1` | ★**成立** |
| `375-a1` 说「仅展示/记账」 | `375-a1-census.md:57` 逐字是「把**配置档**写进 `Decision.Mode`」 | ★**半漂**：「仅」偏强，实测**连展示都没有**（§④） |
| 「dsh 的 `displayedPermissionMode`」 | ★**归错厂商**：它是 **openchamber** 的（`m4-survey.md:48`，`permissionAutoAccept.ts:22-31`）；dsh 那格（`m4-survey.md:85`）是 `permission/preset` 事件 | ★**需更正**：结论（展示用有效档）仍站得住（openchamber ＋ Step-Code `effectiveState`／`m4-survey.md:106`），**论据要换**；dsh 恰是「记切档动作」的先例（支持 §⑦〔丙〕） |
| 「该 seam 的记账更会报错档」 | 该 seam 上 #1／#2 两面**恒等**；真差异在 `panel-inbound`（§⑥） | ★**射程要收窄** |

## ⑨ 门禁 rc（全部在**导出树**上跑，⛔ 工作树）

| 命令 | rc | 读数 |
|---|---|---|
| `go build ./...` ／ `go vet ./internal/{perm,tools,panel,config} ./cmd/wisp` | **0**／**0** | 291 枚生产 `.go`（d22scan 自报口径）；vet 无输出 |
| `go test -count=1 ./internal/perm` | **0** | `ok … 2.233s` |
| `go test ./internal/panel -run Mode` | **0** | `ok 0.808s`（首轮 rc=1 全因导出树缺 `docs/PLAN.md`，补导出后 0 ⇒ 那是**尺缺件**，不是 HEAD 缺陷） |
| `go test ./internal/tools -run 'TestTicket90\|TestTicket363'` ／ `./cmd/wisp -run 'TestTicket114r2\|TestPanelModeConfirm\|TestModeRequest'` | **0**／**0** | `ok 0.326s`／`ok 3.963s`（带 `PATH=<repo>/third_party/sherpa-onnx`，只跑定向 `-run`） |
| `go run ./tools/d22scan -root <导出树>` | **2** | ★**归因：尺量错对象**——75 行输出里 **63 行 `phantom-citation`**（导出树无 `.scratch/`）＋1 行 `ban #8 scope design/ examined 0`（未导出 `design/`）＋1 行「git 不可用⇒gitignore 未套用」。**ban #1–#7 违例 0、ban #8 违例 0**。⇒ HEAD 归属**不可得**（工作树那侧 3 条在飞写腿，跑出来的数不归 HEAD） |
| 长度门：本件 | **≤160 行** | 见件尾计数 |

★**`gofumpt` 未跑**：本笔没落任何 `.go`，无对象可量；且 `gofumpt -l .` 会扫 `.scratch` 里故意写坏的探针（`rc=2`）——**⛔ 拿它判 HEAD 的 EOL**。旁证：导出树 `internal/perm/store.go` 字节统计 **CRLF 0 ／ 纯 LF 378**（尺：数 `0x0A` 且前驱非 `0x0D`）。

## ⑩ 提交与逐枚 `.go`

- 本笔**唯一**写面：`docs/evidence/s1/377-a1-census.md`（＋仓外 `D:\tmp\wisp\20261024\377a1\**`）。**提交号＝本件所在 commit**（主题逐字以 `票 377 普查腿 377-a1` 开头；只 commit **⛔ push**；`add`／`commit` 两处**逐个文件名字面点名**，⛔ glob／`-A`／`.`；提交前 `git diff --cached --name-status` 核索引）。
- **逐枚 `.go`：0**（`internal/**`、`cmd/**`、`scripts/**`、`.github/**` 一字未改；探针 `.go` 只在仓外 `D:\tmp\…\377a1\` 的导出树与两枚变体树里）。
- ⛔ 未碰：票面／`Status:`／`-done`／台账／`docs/reports/**`／`docs/PLAN.md`／`docs/specs/**`／372／375 任何原句；⛔ 翻框。
## ★停手上报（三件，按派单「未定义即停」）

1. ★★**有 1 处「记账读配置档」是故意的**（派单停下条件 ①）：`perm.Store.Set` 的 `from/to`（`store.go:254-255`，理由逐字 `:207-214`），且它**同时是门判**（`from == to` ⇒ 不落盘）⇒ **⛔ 归进「记账点」一起改**；§⑦〔甲〕按字面执行就破它（`m1` 实测）。
2. ★★**票面点名处落在「配置≠有效」不可能成立的那条 seam 上**（§⑥）：`approval_always.go:151-156`／`config_reload.go:282` 只在 `wisp run`，两面恒等 ⇒ **今天零后果**；真差异只在 `wisp panel-inbound` 的 `MODE-SWITCH from=`。★请编排者判：**按形状**修（防将来有腿的 seam 接上桥），还是**先补记账面**（§⑦ 那三处 mode 缺口）——两者不是同一张票的量。
3. ★★**`Decision.Mode` 的消费者＝0**（§④ 具名穷举）⇒ **`AC#2` 的验收判据「记账行必须能读出实际是 `ask_every_step`」今天没有「记账行」可读**：写腿要么先开一个读者（`tools: call` 加 `mode=`），要么那格会退化成「断言一个没人读的字段存在」——正是票面 `AC#2` 括号里禁的那件事。
