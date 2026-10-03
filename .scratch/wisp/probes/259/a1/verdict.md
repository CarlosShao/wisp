# 票 259 AC#0 读数证据件 — 腿 `259-a1`（只读普查，零 Go 命令）

## §0 起手锚

同一发命令取（`date` + `git log -1` + `git status --porcelain -- cmd internal | wc -l`）：

```
2026-10-03 09:42:28+0800
3138e012 2026-10-03 09:42:10 +0800
4
```

- 时刻：`2026-10-03 09:42:28+0800`
- HEAD：`3138e012`（票面立票锚点是 `4eb228ef`，本腿读数锚点已前移 ⇒ 行号一律现取）
- `cmd`＋`internal` 脏文件计数：4（写腿在飞的痕迹，本腿不触碰）

射程声明：本腿只读源码 + `grep`/`sed`/`git log|show|diff`，**没有跑过任何 `go` 命令**；
凡"必须跑 go 才拿得到的读数"一律写进 §7「量不到」，不填推测。
本件不引用 `frontend/**`／`design/**` 任何内容。票面 AC 框一字未碰。

## §1 票面现量 1–7 逐条复量

复量方式：每一发都直接读盘上文件的行（`sed -n`）或全仓 `grep`，行号一律本腿现取。
锚点 `3138e012` 上，`internal/agent/approval/{approval.go,queue.go,ui.go}` 与
`cmd/wisp/subagent_selfapproval_197_test.go` 都是干净已提交状态（`git status --porcelain` 对这
几枚文件无输出），所以下面的行号可被同一份树复现。

| # | 票面断言（要点） | 复量结果 | 凭据（现取行号 + 盘上逐字） |
|---|---|---|---|
| 1 | 花侧第二实参就是那枚 item 自己身上存的 `it.bind`；`it.bind` 全程只有一处写 | **真**（写点那一半要加两条具名限定，见 §2） | `internal/agent/approval/queue.go:376` = `	spent := it.grants.spend(nonce, it.bind)`；`queue.go:162` = `	it.bind = bindDigest(corr, d.TaskID, d.Tool, d.LevelString(), q.seq, d.Args)`。全仓 `grep -rn '\.bind = ' cmd internal tools` 只命中 `queue.go:162` 这一行 |
| 2 | 存的那一份来自同一个值；比对两端在经由路由的调用里永远同值 | **真**（前两行读数真；"永远同值"是本件 §8 要判的那半句，此处只登记它的两个来源行） | `internal/agent/approval/approval.go:287` = `	s.values[nonce] = bind`；`approval.go:303` = `		return equalSecret(stored, bind)`；喂进 `issue` 的那一份 = `queue.go:336` = `	it.grants.issue(nonce, it.bind)` |
| 3 | 绑定那句结构上不可能为 false；真拦下跨卡的是 ① store 成员检查 ② item 状态检查 | **①真／②不真（本腿具名推翻，详见 §4.3）**：跨卡时 B 本来就 pending，那道状态关是**放行**的，把两件事并列会让人以为撤掉成员检查还有状态检查兜着；另外经由路由的调用今天走不到 `queue.go:368` 的"不等"分支（`q.byID` 里不存在非 pending item） | `approval.go:298` = `	for v, stored := range s.values {`，`approval.go:305` = `	return false`（圈走完没命中就落到这里）；`queue.go:368` = `	if it.state != statePending {`，`queue.go:370` = `		return ErrNotPending`。另外本腿补一枚票面没写的第三道：**item 选择本身是精确键**（`queue.go:231-236` 的 `lookupForAllowLocked` 只读 `q.byID`，不读别名索引） |
| 4 | `spend` 返回 `bool`；API 侧只有 `ErrBadGrant`，注释明写故意不区分 | **真**，一处精度要修正：那句注释本体在 **127–129**，130 是值行，票面写的「127-130」把值行也算进注释范围了 | `approval.go:292` = `func (s *grantStore) spend(nonce, bind string) bool`；`internal/agent/approval/ui.go:130` = `	ErrBadGrant = errors.New("approval: 原生令牌无效（缺失/已用/与本次请求不绑定）")`；`ui.go:127-129` = 「…The message never says which, so the API cannot be used to probe for valid nonces.」 |
| 5 | `ticket242_binding_test.go` 那**六枚**用例直接调 store、自己递不匹配的 `bind`，测不到生产路由 | **半真**：性质真、**枚数不真**。该文件在盘上是 **4 枚** `Test*` 函数（`grep -c '^func Test'` = 4，无 `t.Run` 子用例），⛔ 没有"六枚"这个数；写它的交付 commit `4bf7e683` 自带的信息里逐字写的就是"四枚用例"与"含新 4"。4 枚里 3 枚真直接调 store，1 枚（`:58`）只走 `q.push` 后比对 `it.bind` 与自算摘要。**4 枚无一枚调用 `q.allow` / `allowScoped` / `DecideFromNative`**（该文件 `grep -n 'Allow(\|allowScoped\|DecideFromNative'` 零命中）⇒ 票面那句"测不到生产路由会不会让它比出不同值"成立 | `internal/agent/approval/ticket242_binding_test.go:20-21`（`newGrantStore` 后 `issue`，两枚实参都是文件里写死的 fixture 字面量：一枚假 nonce、一枚"别的 item 的摘要"，本件按规矩只写形状不抄其字节）、`:23`（`spend` 递**不匹配**摘要）、`:27`/`:43`（拒后 `live()` 必须归零的断言）、`:29`（同 nonce 再花一次必失败）、`:34` 那枚里 `:39-40`（空摘要必拒）与 `:48-50`（精确摘要必通）、`:58` 那枚走 `q.push`（`:64`、`:78`）只读 `it.bind`、`:87` 那枚走 `q.push` 后**仍直调 store**（`:104`，receiver 是 `itA.grants`、实参二是 `itB.bind`）。自指复算在 `:71` = `	want := bindDigest(itA.Corr, d.TaskID, d.Tool, d.LevelString(), itA.Seq, d.Args)`（票面引的 `:71` 无漂移）。票面那个"六枚"若在指本文件 + `ticket242_panelface_test.go`（2 枚反射尺，`:30`、`:53`），那 2 枚根本不碰 store ⇒ **无论怎么凑，"六枚都直接调 store"这句话在盘上对不上**，属现量 5 的缺陷 |
| 6 | 四发突变的读数（M-B／M-C3／M-D2／M-E），编排者自述"我只复认了结构、读数在它表里" | **结构复认真／读数本腿判不动**（读数要跑 `go test`，本腿禁跑 ⇒ 见 §7） | M-B 引的 `queue.go:235` 现量逐字 = `	return q.byID[corr]`（允许侧可见面的全部，真）；M-C3 的根因复认真：`queue.go:147-153` 里 `corr` 为空时落 `"approval-" + strconv.FormatUint(q.seq, 10)`，所以同一 `Decision` 连推两枚必然拿到**不同 corr**，摘要在 seq 被消掉后仍不同（`:78` 那枚测试正是这么写的）；M-D2 复认真：出向读面名单 `ticket242_panelface_test.go:21-28` 是**裸字符串名**列表，`:59` 的禁用词只有 `{"grant","allow","approve","nonce","token"}`，`Permitted` 一枚都不沾 ⇒ "改名即绕过"是结构事实；M-E 复认真：`ui.go:167-171` 的 `PanelAPI` 今天只有 `Reject`/`Head`/`View`，**无 `Allow`** |
| 7 | 载具前置没做：`subagent_selfapproval_197_test.go:109` 仍是 `TaskID == CorrelationID` | **真，行号无漂移** | `cmd/wisp/subagent_selfapproval_197_test.go:109` = `			TaskID: taskID, CorrelationID: taskID,`（全仓 `grep -n 'TaskID: taskID, CorrelationID: taskID'` 只命中这一行） |

补一条现量里没有、但会影响后面判断的形状（本腿现取）：**`Gate.Replay` 在产码里零调用者**——
`grep -rn --include=*.go '\.Replay(' cmd internal tools scripts` 只命中
`internal/agent/approval/queue_test.go:283` 一枚测试；票面把"重放"列为要复量的形状之一，
所以这一枚要先记账：**今天重放那条路根本不从生产入口可达**（细节见 §2）。

## §2 `it.bind` 写点全名册

字段本体：`internal/agent/approval/queue.go:48` = `	bind   string`，宿主是
`type qitem struct`（`queue.go:31`）。**类型未导出、字段未导出**，且 `*qitem` 在全仓只出现在
未导出函数与未导出字段上（`queue.go:69/70/76/77/141/195/208/231/250`），`Queue` 的导出面
只有 `NewQueue`/`Timeout`(`:126`)/`WarningLead`(`:129`)/`Depth`(`:133`)/`LiveApprovals`
（`pending_read.go:104`，返回的是值拷贝）——**没有任何导出方法把 `*qitem` 递出去**。
⇒ 包外既不能命名这枚字段，也不能拿到那枚指针。

写点名册（穷举方式：`grep -rn '\.bind = ' cmd internal tools` 一发 + `grep -rn 'qitem{' cmd internal tools`
一发 + `grep -rn '&it\.bind\|&fresh\.bind'` 一发 + 全包 `grep unsafe|linkname|reflect` 一发；
四发合起来覆盖赋值、复合字面量初始化、取地址、反射/绕过四条形状）：

| 号 | 位置 | 形状 | 能否进入 `issue`／`spend` 那两条路 |
|---|---|---|---|
| W1 | `internal/agent/approval/queue.go:162` | `it.bind = bindDigest(corr, d.TaskID, d.Tool, d.LevelString(), q.seq, d.Args)`——生产里**唯一一处对活item的赋值** | 能，且是唯一被用的那份值。`corr` 是 push 自己刚算出的那枚（`:148-154`：空则 `"approval-"+seq`，撞键则 `原corr#seq`），与 `it.Corr`（`:156`）同一个字符串 |
| W2 | `internal/agent/approval/queue.go:155-161` | 复合字面量 `it := &qitem{…}` **不含 `bind:`** ⇒ 该字段此刻是零值 `""` | 观察不到：`push` 整个函数体在 `q.mu` 里（`:142-143` 上锁 / `defer` 解锁），item 直到 `:163-165` 才进 `q.pending`/`q.byID`/别名索引，而这三行都在 W1（`:162`）**之后**。⇒ 没有任何一条路能拿到 bind 还是 `""` 的活item |
| W3 | `internal/agent/approval/queue.go:562` | `fresh := *it`（`replay` 里的结构体拷贝）——严格讲这是**第二次写 `bind` 字段**（写到一枚新栈上 item 的字段），值＝历史item那份摘要 | **不能**：`fresh` 之后只被赋 `grants/state/answer/replayOf`（`:563-566`），从未 append 进 `q.pending`/`q.byID`/`q.history`，函数在 `:572` 返回的是 `d, d.CorrelationID` 就结束 ⇒ `fresh` 被丢弃。再加一枚独立事实：`Gate.Replay`（`gate.go:749`）**在生产码里零调用者**（`grep -rn '\.Replay(' cmd internal tools scripts` 只命中 `internal/agent/approval/queue_test.go:283`） |
| W4 | `internal/agent/approval/pending_read_test.go:195-196` | 测试在包内手搓 `&qitem{Corr: "ghost-answered", Dec: dLive, state: stateAnswered}` 与 `stateDropped` 那枚，塞进 `q.pending`（`:198`）——bind 是零值 `""`，**且 `grants` 是 nil** | 不能：这两枚只进 `q.pending`，**没进 `q.byID`**，而允许侧只读 `q.byID`（`lookupForAllowLocked`，`queue.go:231-236`）⇒ 任何针对它们的 allow 落在 `queue.go:366` 的 `ErrUnknownCorrelation`，走不到 `queue.go:376`，也就碰不到那枚 nil store |

⚠ 关于 W3 的一句必要区分（免得被当成"这条路本来会不同值"）：
就算未来有腿把那枚被丢弃的 `fresh` 真登记进队列（那行死代码看起来正是这个意图），
它的 `bind` 会带着**旧 corr** 的摘要而 item 的键是 `corr#replay`/`corr#seq`
（`queue.go:568-570` 只改 `d.CorrelationID`，从不重算摘要）——
那削弱的是"**摘要覆盖了什么**"，⛔ 不是"存的那份与花的那份会不会不同"：
`issue` 与 `spend` 读的是**同一个字段的同一份值**，两端仍然同值。
本件的 §8 只判后者，前者是票面 AC#1 ⓑ 那一支的问题，不归本腿选边。

读点（不是写，列出来是为了让"写一次读两次"这件事可核）：`queue.go:336`（喂 `issue`）、
`queue.go:376`（喂 `spend`），包外零读点，测试读在
`ticket242_binding_test.go:68/72/82/104`。

自我对抗（这一节里我唯一没能证死的方向）：W1 之后**没有任何代码再改这份摘要**这件事，
我是靠上面那四发 grep 的组合来支撑的，不是靠通读全部包外码。剩下的形状——
`go:linkname`、`unsafe` 指针运算、把 `*qitem` 藏进某个 `any` 里再断言回来——
在 `internal/agent/approval` 里**一枚都没有**（`grep -rn 'unsafe\|linkname\|reflect' internal/agent/approval`
的非注释命中全在测试文件，且只走 `tools.Decision` 与 `PanelItem` 两枚导出类型的反射；
反射也设不了未导出字段）。⇒ 就这枚字段而言，"包外改不动"我是认账的，但它属"读过作用面才说的否证"，
不属"grep 一次就结案"。

## §3 `spend`／`issue` 调用点全名册 + 实参来源

### 3.1 字面调用点（`grep -rn --include=*.go '\.spend(' cmd internal tools scripts` 全量）

**产码：一处，没有第二处。**

| 号 | 位置 | 实参一（nonce）来源 | 实参二（bind）来源 |
|---|---|---|---|
| S1 | `internal/agent/approval/queue.go:376`（`allowScoped` 函数体内） | 形参 `nonce`，即 `allowScoped(corr, nonce string, forSession bool)` 的第二参（`queue.go:361`） | **`it.bind`**，其中 `it` 是同一函数 `queue.go:363` 由 `q.lookupForAllowLocked(corr)` 取回的那枚 item；receiver 是 `it.grants`——**bind 与 store 同出一枚对象** |

`internal/agent/approval/ticket242_binding_test.go:23/29/40/50/104` 是五发**测试内**直接调 store，
`104` 那发递的还是另一枚 item 的 `itB.bind`——⛔ 这不是生产路由，本件 §1·现量 5 已具名。

### 3.2 能走到 S1 那一句的入口路（答复入口全名册）

`allow`／`allowScoped` 的字面调用点共三处，全部在 `internal/agent/approval/`：

| 号 | 位置 | 通向 | 产码调用者（包外可达性） |
|---|---|---|---|
| E1 | `internal/agent/approval/gate.go:625`（`nativeAPI.Allow`） | `q.allow` → `queue.go:344` → `allowScoped(...,false)` | **零**。全仓 `grep -rn '\.Native()' cmd internal tools`：非测试命中只有定义处 `gate.go:616` 与 `replies.go:363` 那发 `AllowSession`；`Native().Allow(...)` 只出现在测试（`queue_test.go`、`batch_test.go`、`ticket87_*`、`ticket97_*`、`internal/tools/wiring_test.go`、`cmd/wisp/subagent_selfapproval_197_test.go`）。⇒ 这枚导出方法今天**只被测试按** |
| E2 | `internal/agent/approval/gate.go:723`（`Gate.DecideFromNative` 的 `r.Allow` 真分支） | 同上 | 一处：`internal/agent/approval/replies.go:328`（`Replies.Allow`，`Grant: card.Grant`）← `cmd/wisp/approval_reply.go:215`（`replySurface.allow`）← 动词 `yes`（`cmd/wisp/approval_reply.go:564`）← 两枚控制台循环：`cmd/wisp/approval_reply.go` 的 `runReplyLoop`（`:530`）与 `cmd/wisp/resident_task_source_windows.go:381` 的 `runConsoleLoop` |
| E3 | `internal/agent/approval/gate.go:664`（`Gate.allowSession` 里 `q.allowScoped(corr, grant, true)`） | 直接进 S1 | 一处：`gate.go:648`（`nativeAPI.AllowSession`）← `internal/agent/approval/replies.go:363`（`Replies.AllowSession`，`card.Grant`）← `cmd/wisp/approval_reply.go:259`（`replySurface.session`）← 动词 `session`（`cmd/wisp/approval_reply.go:566`） |

三条路的 `bind` 实参**都是 S1 那一句里现读的 `it.bind`**，没有任何一条把摘要当参数往外接：
`NativeAPI` 的两个 allow 方法签名是 `Allow(ctx, correlationID, grant string)`（`ui.go:146`）与
`AllowSession(ctx, correlationID, grant string)`（`ui.go:158`），答复侧的 `Request` 结构体
（`gate.go:709-715`）字段只有 `CorrelationID/Allow/Grant/Reason/Source` 五枚——
**答复侧连 tool／args／level／seq 都递不进来**，所以"花侧从答复侧递交的请求重算摘要"
在现签名上**没有原料**（这是作用面读数，不是"某修法做不到"的结论；选边归编排者）。

**任务点名的两条要具名的形状，读数如下：**
- **托盘「允许一次」那条路：今天不存在。** 托盘右键菜单的在册项是
  `internal/ball/tray_windows.go:86/88/89/91` 那四枚——「打开面板」「静音」「暂停唤醒」「退出」，
  没有任何 allow 项；球侧非测试码对 `approval.` 的引用为零
  （`grep -rn 'approval\.' internal/ball/*.go` 去掉测试后无命中），它只接否决键
  （`cmd/wisp/resident_ball_windows.go:279-286` 的 `recordCancelHotkey` →
  `cmd/wisp/resident_approval_windows.go:179` 的 `ra.cards.Veto(...)`，走的是 veto 那一路）。
  顺带一枚同族读数：**常驻 GUI 这条腿今天根本没有 allow 入口**——
  `grep -rn --include=*.go '\.Allow(\|\.AllowSession(' cmd/wisp/resident_*.go` 非测试命中为零，
  它只挂 Esc 否决通道（`cmd/wisp/resident_approval_windows.go:157-159` 明写
  「单击球／KWS 否决词／面板拒绝三条仍按各自归口未接入」）。
  ⇒ 票面问的"第二处答复入口"在盘上只有 E2／E3 两枚（＋E1 那枚测试专用）。
- **面板侧的 allow：走不到 S1，够到的是另一件事。**
  `Gate.DecideFromPanel`（`gate.go:730`）在 `r.Allow` 真时直接
  `gate.go:740` 返回 `ErrPanelAllow`，在此之前只做 `gate.go:738` 的 `q.revokeGrants(corr)`
  （→ `queue.go:422-428` → `it.grants.revoke()`，`approval.go:309-313` 把整张 map 换成空 map）。
  **那是"烧掉这枚 item 全部在册 nonce"，不是 spend**：它一次都不读 `bind`；
  `revokeGrants` 的调用者全仓只有 `gate.go:738` 这一处。

### 3.3 `issue` 的调用点与 `bind` 来源

| 号 | 位置 | 存的值 | 上游 |
|---|---|---|---|
| I1 | `internal/agent/approval/queue.go:336`（`grantNonce` 内，`q.mu` 已持，`:335-337`） | **`it.bind`**——与 S1 读的是同一枚对象的同一个字段 | `grantNonce` 的调用者全仓只有 `internal/agent/approval/gate.go:506` 一处，传的是 `gate.go:501` 那次 `g.q.push(d)` 刚返回的 `it`；`nonce` 则是 `queue.go:327` 现 `mintGrant()`（`approval.go:242-248`）出来的 |

⇒ 一条 item 的 `issue` 只发生在 `push` 之后、`Prompt` 之前
（`gate.go:501 → :506 → :528 → :530`），且 `grantNonce` 每枚 item 只被叫一次。
store 侧还有一件事要钉住：`s.values` 的**唯一写点**是 `approval.go:287`（`issue`），
其余 `s.values` 出现处是 `approval.go:298`（读）、`:302`（删）、`:312`（`revoke` 换成空 map）、
`:319`（`live()` 计数）。**没有任何一条路能绕过 `issue` 往 store 里塞一份别的摘要。**

## §4 store 归属（per-item 与否）+ 跨卡今天被谁拦

### 4.1 store 是 per-item 的：构造点复量

- 字段：`internal/agent/approval/queue.go:49` = `	grants *grantStore`——挂在 `qitem` 上，**不在 `Queue` 上**。
  `Queue` 的字段全表在 `queue.go:62-79`（`mu/timeout/warn/maxPend/maxRepl/seq/pending/byID/alias/history/logf`），
  **没有任何一枚 store／map 形式的"全局令牌池"**。
- 类型本体：`internal/agent/approval/approval.go:276-279`，构造函数
  `approval.go:281` = `func newGrantStore() *grantStore { return &grantStore{values: map[string]string{}} }`。
- `newGrantStore()` 的全部调用点（`grep -rn 'newGrantStore' cmd internal tools`）＝
  `queue.go:158`（`push` 里给每枚新 item 一枚）、`queue.go:563`（`replay` 里那枚**被丢弃的** `fresh`，见 §2·W3）、
  以及 `internal/agent/approval/ticket242_binding_test.go:20/38/47`（测试自建）。
  ⇒ **每枚经路由的 item 有且只有自己的那一枚 store；两枚 item 今天不可能共享同一枚。**

⛔ 两枚**同名不同物**的东西，不许混进这张账（本腿特意点名）：
`internal/agent/approval/gate.go:88` 的 `grants GrantRecorder` 与
`internal/tools/bridge.go:194` 的 `grants: o.Grants` 都是 **D45 会话授权记账**
（接口定义 `gate.go:66-68`，只有一个 `Record(ctx, tool, pattern)`，落盘实现是
`internal/session/grants.go:146` 那条 `approval_grant` 行），
跟 `qitem.grants` 那枚内存 nonce 池**没有类型关系也没有调用关系**。

### 4.2 跨卡花令牌（A 的 nonce 拿去 allow B）今天被谁拦住

逐行走一遍这条路，产码上唯一可能的形状是 `allowScoped(corr_B, N_A)`：

1. `internal/agent/approval/queue.go:363` `it := q.lookupForAllowLocked(corr)` →
   `queue.go:231-236` 只读 `q.byID[corr]`（**精确键**，别名索引 `q.alias` 在这条路上看不见，
   `queue.go:71-76` 的注释与 `internal/agent/approval/ticket97_alias_direction_test.go` 钉的就是这个方向）。
   ⇒ 这一关拦的是"借名字把 allow 引到别的 item"，⛔ 不是摘要比对。
2. `internal/agent/approval/queue.go:368` 状态关（见 4.3，这一关在跨卡形状里**必然是放行**的，
   因为 B 本来就 pending）。
3. `internal/agent/approval/queue.go:376` → `internal/agent/approval/approval.go:292`：
   `approval.go:298` `for v, stored := range s.values` 遍历的是**B 自己那枚 store**；
   A 的 nonce 从来只被 `approval.go:287` 写进 **A 的 store**（§3.3 已钉：`s.values` 唯一写点就是 `issue`），
   所以 `approval.go:299` `if !equalSecret(v, nonce)` 对 B 的每个 key 都成立 → 一路 `continue`
   → 圈走完落到 `approval.go:305` `return false`。

**结论（这一栏只回答"是谁拦的"）：拦下跨卡的是成员检查，即 `approval.go:298-305` 那一圈加 `:305`
那句落空 return；绑定比对 `approval.go:303` 在跨卡形状里连执行都没被执行到**（它只在
"nonce 命中了本 store 的某个 key"之后才跑）。⇒ 票面现量 3 的 ① **真**。

### 4.3 顺手推翻现量 3 的 ②（这一条要说清，不能混着算）

票面把「`queue.go:368-371` 的 `it.state != statePending` ⇒ `ErrNotPending`」列为"真正拦下 A 的令牌花在 B 上"
的两处之一。**这半句在本腿的读数里不成立，两个理由：**

- **它拦的不是跨卡。** 跨卡时 B 是活着的 pending item，这道关**直接放行**，
  真正让调用失败的是下一句（成员检查）。把这两件事并列成"两处防线"会让人以为
  拿掉成员检查还剩状态检查兜着——事实是拿掉成员检查后 `:303` 才是唯一还能跑的，而那正是本票要判的恒等式。
- **经由路由的调用今天走不到 `:368` 的不等分支。** 全仓 `grep -rn 'it\.state = ' internal/agent/approval`
  的非测试命中只有两行：`queue.go:304`（`deliver` 里置 `stateAnswered`）与 `queue.go:330`
  （`grantNonce` 出错时置 `stateDropped`），**两行都紧跟一次 `q.dropLocked(...)`**
  （`:305` / `:331`）且都在同一个 `q.mu` 临界区里；而 `dropLocked` 在 `queue.go:286`
  `delete(q.byID, it.Corr)`。⇒ **`q.byID` 里不存在"非 pending 的 item"**，
  而允许侧只可能从 `q.byID` 拿到 item，所以 `:368` 的"不等"分支是**结构上空转**的。
  真正拦住"答复一枚已结算的卡"的是 `queue.go:366` 的 `ErrUnknownCorrelation`（查不到）。
- 同一族里**真在跑**的两道是：`queue.go:301-303`（`deliver` 自己的状态守卫，
  它拿的是调用方手里的 `it` 指针，能撞上"`:376` 花完锁、`:382` 之前被别人抢先结算"这发竞态 ⇒
  `queue.go:383` 返回 `ErrNotPending`），以及 `queue.go:288` 的 `it.grants.revoke()`
  （每枚 item 一离开队列就把整张 nonce 池换成空 map，`:288` → `approval.go:309-313`）。

⚠ 本腿没有跑过任何 `go` 命令，所以"今天走不到"这句是**从上述赋值/删除都在同一临界区这一码上读出来的**，
不是从覆盖率或插桩读出来的；要把它升级为仪器，属 §7 那一条。

## §5 烧牌语义那一问

### 5.1 那句注释在盘上是否为真

`internal/agent/approval/approval.go:302` 盘上逐字：

```go
		delete(s.values, v) // consumed whether or not the binding matched
```

**真。** `delete` 在 `:302`，返回值在 `:303`（`return equalSecret(stored, bind)`），
删除发生在绑定结果被算出来**之前**，且在"nonce 命中本 store 某个 key"这一条路径上无条件下发。
所以"错绑的那一发照样把 nonce 烧掉"这句话与码一致。

作用面要说清的三条边界（免得这句话被读大）：
- **只有 key 命中才烧。** `nonce == ""` 在 `approval.go:293-295` 早退，**一次 `delete` 都不跑**；
  nonce 非空但不在本 store ⇒ 圈走完落到 `approval.go:305`，同样**不烧任何东西**（别的 item 的池子更不会被碰）。
- **`revoke` 是另一枚烧牌器，不走 `:302`。** `approval.go:309-313` 是整张 map 换新，
  面板侧那条路（`gate.go:738` → `queue.go:426`）与每次结算（`queue.go:288`）用的都是它。
- 仪器侧：钉住 `:302` 这语义的是
  `internal/agent/approval/ticket242_binding_test.go:26-28` 与 `:43-45`（`s.live() != 0` 那两发），
  两发都是**包内直调 store**，不是路由级（见 §1·现量 5）。

顺带一枚与本问直接相邻的码/注不符（票面 AC#2 的靶子就在这儿）：
`approval.go:290-291` 的注释自称 "spend validates and consumes. **It reports why a rejection
happened** in user-safe terms"，而签名是 `func (s *grantStore) spend(nonce, bind string) bool`
（`approval.go:292`）——**它什么都没能 report**，四种落空（空 nonce／查不到／已烧／绑定不合）
在返回值上全是同一个 `false`。这是注释在替一个不存在的能力作证，本腿按读数登记，⛔ 不改码。

### 5.2 有没有一条路让一枚错绑 nonce 被拒后**又**被同一枚 item 的花费路径再消费一次

**没有。** 逐段凭据：

- 烧掉之后 map 里**没有那个 key 了**（`:302`），第二发同 nonce 走进 `approval.go:298` 的圈时
  每一枚 `equalSecret(v, nonce)`（`:299`）都不成立 ⇒ 落到 `:305` `return false`。
  ⇒ "再消费一次"这件事在数据结构上不成立：**没有东西可消费**，第二发连 `:302`、`:303` 都到不了。
- 而且**今天生产里根本没有第一发**：本件 §8 判的是"存的那份＝花的那份"，因此经路由的
  `spend` 只可能返回两种值——nonce 命中且摘要同值 ⇒ `:303` 恒为 `true`（成功），
  或 nonce 空/不在本 store ⇒ `:293`／`:305` 的 `false`（**不发 `delete`**）。
  ⇒ **`delete` 与 `false` 同时发生的那条支路（错绑烧牌）今天从生产路由不可达**，
  它只在包内直调 store 的白盒测试里跑得起来。
- **"烧牌会不会掩盖第二次尝试"**：会，但掩盖的是**可读性**不是**次数**。
  `spend` 只回 `bool`，两发落在 `internal/agent/approval/queue.go:378-381` 同一条支路：
  同一条审计行文本（`queue.go:379` 的 `"approval: FORGED-OR-STALE allow rejected corr=%s (native grant missing/spent/misbound)"`）
  与同一个 `ErrBadGrant`（`internal/agent/approval/ui.go:130`）。
  ⇒ 盘上能区分的只有**这行出现了几次**，区分不了"第一次是错绑、第二次是查无此牌"。
  这正是票面现量 4 那条"拒因今天不可指名"的同一件事，本腿不重复裁。
- 一枚相邻形状（作用面读数，非结论）：若绑定那一支将来真能返回 `false`（ⓑ 那一支才会出现），
  由于 `:302` 在 `:303` **之前**，一发"真 nonce + 错摘要"的尝试会把**合法用户手里那张卡唯一能用的牌烧掉**
  ——那张卡随后再也 allow 不动（`grantNonce` 每枚 item 只被叫一次，`gate.go:506`），
  只能重新显示成新卡片。烧牌语义与活绑定语义叠在一起时的这一格，编排者裁 ⓐ／ⓑ 时要拿在手里。
  ⛔ 本腿不据此选边，也不判"ⓑ 做不到"。

## §6 我可能写错的条目（自我对抗）

1. **"唯一可达写点"是靠 grep 组合撑的，不是靠类型系统。** 我用的四发是
   `\.bind = ` / `qitem\{` / `&it\.bind|&fresh\.bind` / `unsafe|linkname|reflect`，根都在
   `cmd internal tools scripts`。如果将来有第五种形状（例如把 `bind` 改名成别的字段、
   或把 `qitem` 内嵌进另一枚 struct 让 `qitem.bind` 变成 `outer.item.bind`），
   第一发 `\.bind = ` 仍能抓到（改名就不是 `.bind` 了，但那时"写点名册"这题的答案本身也变了）。
   ⇒ 本件的措辞是"锚点 `3138e012` 上"，⛔ 不是"永远"。
2. **"经路由的调用走不到 `queue.go:368` 的不等分支"这句最容易被推翻，也最要紧。**
   我的凭据是"两枚 `it.state =` 写点（`:304`、`:330`）都在同一临界区内紧跟 `dropLocked`，
   而 `dropLocked` 在 `:286` 删 `q.byID`"。推翻它只需要第三种状态写点，或一处
   **不持 `q.mu` 就写 state** 的码。我把全部写点列完了（`grep -rn 'it\.state = ' internal/agent/approval`
   非测试只有那两行）。⚠ 但我**没有**验证"重入"这条：`push` 在持锁状态下调用 `q.logf`
   （`queue.go:166-167`），若某枚 logf 实现回调进 `Queue` 的任何加锁方法，结果是**死锁**而不是
   状态被插队——死锁是 fail-stop，造不出"byID 里有非 pending item"。所以这句我认账，
   但它的强度是"码上读出来的"，不是"插桩证明的"（见 §7 第 3 条）。
3. **我把 `deliver` 的守卫（`queue.go:301-303`）说成"真在跑的那道"，依据是 `allowScoped`
   在 `queue.go:377` 解锁、`:382` 才重新上锁这段窗口。**这段窗口里另一发并发答复确实可能先落。
   我**没有**量出任何一枚在册用例覆盖了这发竞态——所以"真在跑"应读作"码上可达"，
   不读作"今天被证明过"。
4. **"托盘没有 allow 项"我只读了 `internal/ball/tray_windows.go:86/88/89/91` 那四枚
   `appendItem`。**如果有第二处构造托盘菜单的码，我的"零"就漏了。
   复核方式：`AppendMenuW` 的绑定声明在 `internal/ball/win32_windows.go:45`，
   其非测试使用者本腿只量到 `internal/ball/tray_windows.go` 一枚。
   ⇒ 这条我认，但射程是"ball 这一包的托盘"，不是"任何 GUI"。
5. **"`NativeAPI.Allow` 零产码调用者"这句依赖 `grep -rn '\.Native()' cmd internal tools`。**
   另一种形状是"把 `NativeAPI` 当参数/字段类型接住再调 `.Allow(...)`"，本腿也扫了
   `NativeAPI` 这个词本身——非测试命中只有 `internal/agent/approval/gate.go:616`（返回类型）
   与 `internal/agent/approval/ui.go:143`（接口声明）加若干注释，**没有任何字段或参数以它为类型**。
   ⇒ 这句认账。
6. **§1·现量 5 那条"枚数不真"可能是我对"用例"的理解窄。**我按 `^func Test` 计数（4）。
   若把一次测试函数里的多个断言各算一枚，"六枚"也凑得出来（`:23`、`:29`、`:40`、`:50`、`:72`、
   `:104` 恰好六发断言级判断）。⚠ 但票面原话是"那六枚**用例**"，而写这枚文件的 commit
   `4bf7e683` 的自带信息里逐字写的是"四枚用例"与"含新 4"，所以我维持"枚数不真"这一判，
   只把另一种读法写在这儿。
7. **同名文件陷阱我已经踩过一次并改正**：`find . -name approval.go` 给出两枚——
   `internal/agent/approval/approval.go` 与 `internal/panel/approval.go`。本件里凡出现
   裸 `approval.go:NNN` 一律指前者；后者本腿**一次都没引用过**。
   ⚠ §0/§1 的早期草稿里有若干短引，本条就是那条短引的豁免说明——AC#1 的复核者若按
   `internal/panel/approval.go` 去核行号会看到完全不同的内容，那是同名陷阱不是读数错。
8. **"恒等式"这个词我只在"存进去的那份摘要与花掉时传进去的那份摘要同值"这一格用。**
   我没有、也不能把它读成"摘要覆盖了实际执行的那一发请求"。后者要的是 `bindDigest`
   的六个入参在**答复时刻**重算，而答复侧连 tool/args 都不带（§3.2 末段）。
   这一格若被误读，AC#1 会选错支，所以 §8 把它单列成一句话。
9. **读数锚点之后的三枚写腿**（`cmd/wisp`＋`internal/config`、`internal/tools`＋`internal/risk`、
   `internal/panel`）此刻在飞。按票面它们都不带 `internal/agent/approval` 的写面，
   但 `internal/panel` 与 `cmd/wisp` 的**编译**会带上这包——所以本件说的"今天"是
   `3138e012` 那一刻的码面；AC#1 动手前应当重取同一批发点（§7 第 5 条）。

## §7 量不到的地方（具名，⛔ 不填推测）

本腿第一硬规是"一枚 Go 命令都不许跑"（三枚写腿在飞：`cmd/wisp`＋`internal/config`、
`internal/tools`＋`internal/risk`、`internal/panel`）。下面每一条都是**只有跑起来才有读数**的，
本腿一律留空并注明该跑什么；⛔ 没有一条我用"应该／大概"去填。

1. **现量 6 的四发突变读数全部量不到**：M-B 的"同发 12 枚在册路由级用例红"、
   M-C3 的"全包 66 PASS"、M-D2 的"66 全绿"、M-E 的"本包仍 66 PASS／0 FAIL"，
   以及"已验非编译失败冒充红"那半句。需要的命令＝`go test ./internal/agent/approval/...`
   逐发突变跑。⇒ 本件只在 §1·现量 6 那行给了**结构复认**（引用的行号与码形状对不对），
   颜色与枚数一格未填。
2. **"approval 包今天到底多少枚在册用例"这个分母量不到**（66？45？）。
   写 `4bf7e683` 的自述是"45 PASS"，票面 259 现量 6 引的是"66"，两数差 21 枚
   ⇒ 这 21 枚是不是 `ticket242_panelface_test.go`（2 枚）与今天新落的几枚文件造成的，
   只有 `go test -run . ./internal/agent/approval/ -v` 数得出来。本腿只登记两数并存这一事实。
3. **覆盖／可达性的插桩读数量不到**："经路由的调用走不到 `queue.go:368` 不等分支"、
   "`approval.go:303` 的 false 支今天零生产触发"这两句，我给出的是**码面推理**（§4.3、§5.2），
   不是插桩证明。要钉死需要 `go test -cover` 或一发**路由级**突变（把 `:368` 的判断取反看有没有用例红）。
4. **并发语义实测量不到**：`allowScoped` 在 `queue.go:377` 解锁、`:382` 重锁这段窗口里
   两发并发 allow 的真实落点，需要 `go test -race` 加一发定向竞态用例。
   ⛔ 本腿不声明"没有竞态问题"，只声明"这发窗口在码上存在，实测读数取不到"。
5. **写腿在飞的漂移**：`internal/agent/approval/{queue.go,approval.go,ui.go,gate.go}` 的
   mtime 是今天 09:16–09:26（`ls -la` 现量），本腿锚点 `3138e012` 在其后。
   这三枚写腿的**下一步提交**是否改动这包，量不到（也不该由我判断）——
   票面「排程与互斥」自己写了"internal/panel 与 cmd/wisp 编译都会带上 internal/agent/approval"。
   ⇒ AC#1 落地前必须重取本件 §2/§3 那两批发点。
6. **票面 M-D2／M-E 那两支"改名即绕过／能力侧零仪器"的完整作用面量不到**：
   要判"改名的字段／新增的方法名能不能被现有尺抓到"，除了我已经读到的
   `internal/agent/approval/ticket242_panelface_test.go:21-28`（裸字符串名单）与
   `:59`（五个禁用词），还得知道**别的包里**有没有第二枚数 `PanelItem` 字段或
   `PanelAPI` 方法名的尺。我只在 `internal/panel/l2_grant_boundary_test.go` 里读到它管的是
   **入站方法名名单**（`:229`、`:1252-1255` 那几枚表），⛔ 它不管 `approval.PanelItem` 的字段，
   也不管 `PanelAPI` 的方法集。⇒ 这一条是**射程判断，非内容引用**
   （那三枚冻结件我只读内部判射程，本件不引其内容，且其中出现的从前端侧借来的方法名字面
   本腿一律未去核对、未转述）。
7. **`equalSecret`（`internal/agent/approval/approval.go:266-271`）的恒定时间性质在本机
   是否真成立**（编译器有没有把 `subtle.ConstantTimeCompare` 优化掉）量不到：
   要读汇编或跑微基准。本腿只登记"它用了 `crypto/subtle`、空串早退"这一码面。

## §8 结论（恒等式成立／不成立 + 凭据行号）

**判：恒等式成立。** 锚点 `3138e012` 上，本腿**找不到任何一条真路径**能让
`spend(nonce, bind)` 收到的 `bind` 与当初 `issue(nonce, bind)` 存进同一枚 store 的那一份**不同**。
凭据链五步，每步一枚具名行：

1. store 属于每枚 item 自己：字段 `internal/agent/approval/queue.go:49`，
   构造 `queue.go:158`（`push` 内），唯一构造函数 `internal/agent/approval/approval.go:281`；
   `Queue` 自己身上没有池子（`queue.go:62-79`）。
2. 唯一的 `issue` 点喂的就是这枚 item 的字段：`queue.go:336`；
   而 `grantNonce` 全仓唯一的调用者是 `internal/agent/approval/gate.go:506`，
   传的是同一函数 `gate.go:501` 那次 `push` 刚返回的 `it`。
3. 唯一的 `spend` 点的 receiver 与第二实参**同出一枚对象**：
   `it` 来自 `queue.go:363`（`lookupForAllowLocked`，只读 `q.byID`），
   句子是 `queue.go:376` = `it.grants.spend(nonce, it.bind)`。
4. `it.bind` 在 `issue` 之后再没有被写过：可达写点唯一 = `queue.go:162`（`push` 内、
   item 进 `byID` 之前）；`queue.go:562` 那枚结构体拷贝写的是 `fresh.bind`，
   而 `fresh` 在 `queue.go:572` 就被丢弃；测试手搓的 `&qitem{}`
   （`internal/agent/approval/pending_read_test.go:195-196`）进不了 `:363` 的查找。
   写点普查四发的名字与边界在 §2、§6。
5. store 里那份摘要也没有第二条进去的路：`s.values` 唯一写点 =
   `internal/agent/approval/approval.go:287`（`issue` 体内），
   其余出现处是 `:298` 读、`:302` 删、`:312` 整张换空、`:319` 计数。

⇒ 于是 `approval.go:303` 的 `equalSecret(stored, bind)` 两端在同一次调用里读的是
**同一个字符串**：`stored` 是 `it.bind` 在 `queue.go:336` 那一刻的副本，`bind` 是
`it.bind` 在 `queue.go:376` 那一刻的读数，而这两者之间该字段没有任何写点（第 4 步）。
`internal/agent/approval/approval.go:267-269` 那个"任一端为空即 false"的早退也不构成例外分支：
`bindDigest`（`approval.go:253-261`，`hex.EncodeToString(sha256.Sum)` 于 `:260`）恒返回 64 个十六进制字符。

票面点名要走的形状，逐条给结论（每条都能对上上面五步之一）：

| 形状 | 经这条路两端会不会不同值 | 凭据 |
|---|---|---|
| 重新铸造（对同一枚 item 再签一张牌） | **不会**。今天每枚 item 只被签一次；即便签两次，第二次存的还是同一个 `it.bind` | `gate.go:506`（唯一调用者）／`queue.go:336` |
| 重放 | **不会**。`Gate.Replay` 零产码调用者；`queue.replay` 里那枚带旧摘要的拷贝在 `:572` 被丢弃；重放出去的 `Decision` 要再进队列只能重新 `push`，`bind` 重算 | `gate.go:749-761`／`queue.go:552-575`（`:562`、`:563`、`:572`） |
| `revoke` | **不会**，而且它根本不碰摘要：只把 map 换成空的 | `approval.go:309-313`；调用处 `queue.go:288`、`queue.go:426` |
| 超时 | **不会**——第二发 allow 连 item 都拿不到 | `queue.go:476-497` → `deliver`（`:304-305`）→ `dropLocked`（`:286` 删 `byID`、`:288` revoke）→ `queue.go:366` `ErrUnknownCorrelation` |
| 答复后再答复 | **不会**。同上；退一步就算摸到 `:376`，nonce 已不在 map ⇒ `approval.go:305` | `queue.go:301-303`、`:286`、`:288` |
| 面板侧递交 allow | **不会**：这条路一次都不花牌，只做整池 revoke | `gate.go:731-740`（`:738` 是 `revokeGrants`） |
| 跨卡（A 的 nonce 拿去 allow B） | **不会不同值**——因为**根本到不了比对那一句**：B 的池子里没那把 key，`approval.go:299` 逐枚 `continue`，落 `:305` | `approval.go:298-305`＋`queue.go:231-236` |
| 包内直调 store | **会不同值**，⛔ 但它不是真路径：调用者在 `*_test.go` 里自己造两枚摘要递进去 | `internal/agent/approval/ticket242_binding_test.go:23`、`:104`（§1·现量 5 已具名；该包外没有任何码能拿到 `*grantStore`，它是未导出类型） |

**边界声明（这一句必须和上面一起读）：** 恒等式说的是"存的那份 == 花的那份"，
它⛔ **不等于**"这份摘要覆盖了实际执行的那一发请求"。后者要的是 `bindDigest` 的六个入参
在**答复时刻**被重算，而答复侧能递进来的只有 `CorrelationID/Allow/Grant/Reason/Source`
五枚字段（`gate.go:709-715`；两个 allow 方法的签名 `ui.go:146`、`ui.go:158`）——
**tool／args／level／seq 在答复面上没有原料**。本腿把这一格量成读数，⛔ 不把它写成
"ⓑ 做不到"，也⛔ 不写成"ⓐ 更安全"：**ⓐ／ⓑ 选边归编排者**（票面「未定义即停」与本腿任务同一条）。

## §9 交件判语

- **只读**：本腿只跑过 `date`／`git log`／`git status`／`git show`／`ls`／`find`／`grep`／`sed -n`
  与对 `.scratch/wisp/probes/259/a1/**` 的写入。
- **零 Go 命令**：`go build`／`go test`／`go vet`／`go run`／`go list` 一次都没跑过；
  因此所有需要跑起来的读数都写在 §7 并具名留空，⛔ 未用推测填空。
- **未 push**：本腿共 8 枚 commit（skeleton、§1、§2、§3、§4、§5、§6＋§7、这枚交件），
  全部只在本地 `dev` 上；未跑过 `git push`。
- **AC 框未碰**：`.scratch/wisp/issues/259-*.md` 一字未改（不在本腿任何一次 `git add` 的
  pathspec 里）；本件对票面只做引用与复量，复量不真的条目在 §1 具名推翻、
  并按本腿规矩**没有**替票面改写。
- **凭据零外泄**：本件不含任何真实 nonce／digest 的字节，也不含任何 API key／令牌值；
  凡触及令牌处一律写变量名与形状（`it.bind`、`nonce`、`stored`、`s.values`）。
  §1·现量 5 那一格原先抄了测试文件里的两枚 fixture 字面量，本腿在 §1 定稿时
  已改写成"形状描述＋行号"，⛔ 未留字面量。
- **冻结件**：`internal/panel/tokens_fourway_test.go`、`internal/panel/l2_grant_boundary_test.go`、
  `internal/perm/ticket90_persist_test.go` 只做了**射程判断，非内容引用**
  （前者 550 行、对 `approval.／spend(／bindDigest／grantStore` 零命中；中者 2549 行，管的是
  入站方法名名单、不数 `approval.PanelItem` 的字段；后者 430 行，同样零命中）；
  `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／
  `tools/d22scan/allowlist.txt`／`.github/workflows/ci.yml` 一次都没读过。
- **两层禁令**：`frontend/**` 与 `design/**` 本腿**零读取**，本件亦未转述其任何内容
  （工作树里 `design/**` 此刻正处于被删除状态，那是别人的写面，本腿未碰、不判）。
- **临时件只建不删**：`.scratch/wisp/probes/259/a1/` 下 `verdict.md` 加八枚 `msg-*.txt`
  （每枚 commit 的 `-F` 消息件，含本枚），全部保留、全部提交，未删任何仓内文件。
- **未联系其它会话**，未代转任何消息给前端。
- **本件的射程**：锚点 `3138e012`（`2026-10-03 09:42:10 +0800`）的码面。
  三枚写腿在飞，`internal/agent/approval` 的写面按票面归本票自己，但 §7 第 5 条已声明
  AC#1 动手前须重取 §2／§3 那两批发点。⛔ 本腿不选 ⓐⓑ，不判 AC#1–AC#5。
