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
| 3 | 绑定那句结构上不可能为 false；真拦下跨卡的是 ① store 成员检查 ② item 状态检查 | **真**（两处行号都在；两检查的分工见 §4，不许混着说） | `approval.go:298` = `	for v, stored := range s.values {`，`approval.go:305` = `	return false`（圈走完没命中就落到这里）；`queue.go:368` = `	if it.state != statePending {`，`queue.go:370` = `		return ErrNotPending`。另外本腿补一枚票面没写的第三道：**item 选择本身是精确键**（`queue.go:231-236` 的 `lookupForAllowLocked` 只读 `q.byID`，不读别名索引） |
| 4 | `spend` 返回 `bool`；API 侧只有 `ErrBadGrant`，注释明写故意不区分 | **真**，一处精度要修正：那句注释本体在 **127–129**，130 是值行，票面写的「127-130」把值行也算进注释范围了 | `approval.go:292` = `func (s *grantStore) spend(nonce, bind string) bool`；`internal/agent/approval/ui.go:130` = `	ErrBadGrant = errors.New("approval: 原生令牌无效（缺失/已用/与本次请求不绑定）")`；`ui.go:127-129` = 「…The message never says which, so the API cannot be used to probe for valid nonces.」 |
| 5 | `ticket242_binding_test.go` 那**六枚**用例直接调 store、自己递不匹配的 `bind`，测不到生产路由 | **半真**：性质真、**枚数不真**。该文件在盘上是 **4 枚** `Test*` 函数（`grep -c '^func Test'` = 4，无 `t.Run` 子用例），⛔ 没有"六枚"这个数；4 枚里 3 枚真直接调 store，1 枚（`:58`）只走 `q.push` 后比对 `it.bind` 与自算摘要。**4 枚无一枚调用 `q.allow` / `allowScoped` / `DecideFromNative`**（该文件 `grep -n 'Allow(\|allowScoped\|DecideFromNative'` 零命中）⇒ 票面那句"测不到生产路由会不会让它比出不同值"成立 | `internal/agent/approval/ticket242_binding_test.go:20-21`（`newGrantStore` + `s.issue("nonce-live", …)`）、`:23`（`s.spend("nonce-live", …)` 递不匹配摘要）、`:39-40`、`:48-50`、`:104` = `	if itA.grants.spend("leaked-or-guessed-nonce", itB.bind) {`；自指复算在 `:71` = `	want := bindDigest(itA.Corr, d.TaskID, d.Tool, d.LevelString(), itA.Seq, d.Args)`（票面引的 `:71` 无漂移）。票面那个"六枚"若在指本文件 + `ticket242_panelface_test.go`（2 枚反射尺，`:30`、`:53`），那 2 枚根本不碰 store ⇒ **无论怎么凑，"六枚都直接调 store"这句话在盘上对不上**，属现量 5 的缺陷 |
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

尚未作答。

## §4 store 归属（per-item 与否）+ 跨卡今天被谁拦

尚未作答。

## §5 烧牌语义那一问

尚未作答。

## §6 我可能写错的条目（自我对抗）

尚未作答。

## §7 量不到的地方（具名）

尚未作答。

## §8 结论（恒等式成立／不成立 + 凭据行号）

尚未作答。

## §9 交件判语

尚未作答。
