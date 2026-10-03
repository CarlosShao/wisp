# 票 259 — 我落的"令牌绑定层"**在产码里是恒等式**：`spend` 拿的是那枚 item 自己存进去的同一份摘要 ⇒ `approval.go:303` 那句比对从路由不可达，真挡住跨卡的另有其处

**立票时刻**：2026-10-03 09:31:02+0800，锚点 HEAD `4eb228ef`（`dev`）
**来路**：非实现者对抗验收腿 `242-v1` 的 `docs/evidence/s1/242-grant-binding-v1.md`（252 行／29,758 字节，commit `3fe377e3`→`8fe080d1`→`850a76b4`→`4df747d8`）§3「生产接线」＋九发突变里的 **M-B／M-C3／M-D2／M-E** 四发；母票＝**票 242**（AC#1 判语＝**附条件成立**，AC#2＝**不成立**，AC#3＝**不许勾**）。
**性质**：⚠ 本票**不是**给票 242 擦屁股的账目搬家——它登记的是**一条今天真实存在的安全语义空档**（见现量 3），并且**推翻了我自己在 `A549`／票 242 §4 里写过的那句"绑定层已能挡跨卡"**。

## 现量（编排者本机自跑，2026-10-03 09:3x，逐字；⚠ 引用前先重跑，别把这几行当常量）

1. **花侧的实参是"自己比自自己"**：`internal/agent/approval/queue.go:376` 逐字 `spent := it.grants.spend(nonce, it.bind)`——传进去的第二参数就是那枚 item 自己身上存的 `it.bind`。而 `it.bind` 全程**只有一处写**（`internal/agent/approval/queue.go:162` 铸造时 `bindDigest(...)` 的结果）。
2. **存的那一份也来自同一个值**：`internal/agent/approval/approval.go:287` 逐字 `s.values[nonce] = bind`（`issue` 把 `it.bind` 存进 store）。⇒ `approval.go:303` 的 `return equalSecret(stored, bind)` 两端在**任何经由路由的调用**里永远同值。
3. ⇒ **绑定那一句从路由不可达**（不是"跑得少"，是**结构上不可能为 false**）。真正拦下"A 的令牌花在 B 上"的是两处**别的**检查：① `spend` 的 store 成员检查（`approval.go:298-305` 那圈 `for v, stored := range s.values`——nonce 不在这枚 item 自己的 store 里就落到 `return false`）；② item 状态检查（`queue.go:368-371` `it.state != statePending` ⇒ `ErrNotPending`）。
4. **拒因今天不可指名**：`spend` 返回的是 `bool`（`approval.go:292`），API 侧只有一个 `ErrBadGrant`（`internal/agent/approval/ui.go:130` 逐字文案把三种原因并成一句：`approval: 原生令牌无效（缺失/已用/与本次请求不绑定）`），`ui.go:127-130` 注释**明写故意不区分**。⇒ 票 242 AC#1 那句"拒因**指名**绑定不对（不是"这条路由没在跑"）"这一维今天**零仪器、也零实现**。
5. **我的白盒尺为什么看起来有牙**：`internal/agent/approval/ticket242_binding_test.go` 那六枚用例是**直接调 store**、自己递一个不匹配的 `bind` 实参进去的。⇒ 它测到的是"`grantStore` 这个零件内部会比"，⛔ **测不到"生产路由会不会让它比出不同值"**——这正是本仓那条老规矩：**能力类判据要问生产调用者**（`feedback-adversarial-review-orchestration`／台账 `A559` 第 2 条）。
6. **验收腿的四发突变（我只复认了结构、读数在它表里）**：**M-B**（`queue.go:235` 让允许侧那条路都不跑）⇒ 我那六枚**全绿**、同发 12 枚在册路由级用例红＝"不是这条路由没在跑"那半句无仪器；**M-C3**（把 seq 从摘要里彻底消掉）⇒ **全包 66 PASS**（根因：测试用同一枚 `bindDigest` 自指复算 `ticket242_binding_test.go:71`，而那两枚项 `corr` 天生不同）；**M-D2**（`PanelItem.Permitted bool` ＋ 把 `"Permitted"` 补进出向读面名单）⇒ **66 全绿**＝词面名单**改名即绕过**；**M-E**（照票面那一形给 `PanelAPI` 加 `Allow(correlationID, grant)` 并实现成真花令牌 `p.q.allow(...)`）⇒ **本包仍 66 PASS／0 FAIL**（已验非编译失败冒充红）＝**能力侧零仪器**。
7. **载具前置没做**：`cmd/wisp/subagent_selfapproval_197_test.go:109` 现量逐字 `TaskID: taskID, CorrelationID: taskID,`——票面那条"造两枚同时活着、correlation id 不同"的载具前置**在盘上仍未成立**（⚠ 该文件写面＝`cmd/wisp`，此刻 `255-r2` 在飞，本票按住到它退出）。

## 要建什么

- [ ] **AC#0（第一格，读数不是码）**：判死**有没有任何一条真路径**能让 `spend` 收到与存入时**不同**的 `bind` 值。逐枚列 `bind` 的写点／读点（现量 1＋2 是起点，⛔ 不许照抄），并对每一枚具名回答"经这条路调用时两端会不会不同"。⛔ 本格不许改产码；交件判据＝判不动的那一处要具名写"判不动"，⛔ 不许用"应该会不同"填空。
- [ ] **AC#1 二选一（要 AC#0 交完之后由编排者落一枚具名 `A##` 选边）**：ⓐ **具名降级**＝登记"绑定那一句今天零独立拦截能力，真防线是 store 成员检查＋state 检查"，并把注释与票面口径改到与盘上一致（⛔ 不许留"能挡跨卡"这句话在码里）；ⓑ **做出牙**＝花侧改为**从答复侧实际递交的请求重算摘要**再比（tool／args／level／seq 的来源要换），并补一发**路由级**用例：种下"答复请求的参数与铸造时不同"⇒ 指名用例必须红。⚠ ⓐⓑ 两支都要同时回答 AC#2，不许只挑一支交差。
- [ ] **AC#2 拒因可指名**：把 `spend` 的返回从 `bool` 改成能带原因的形状（或等价物），使"缺失／已用／绑定不对／状态不是 pending"**四种各自可断言**；⚠ 对外文案是否**区分到界面**属 `ui.go:127-130` 那句"故意不区分"的范围——**编排者已裁：内部与审计现场必须区分，对外文案维持合并**（理由＝审计要能一眼归因，界面不暴露"差几个字节就对上"那种逼近信息；撤销口令「259 对外也区分」）。判据＝四种拒因各一枚用例指名红，⛔ 且**不许复用**权限判定那条文案（本仓既有规矩）。
- [ ] **AC#3 三枚能力尺（补票 242 AC#2 不成立那一格，随本票一起做，⛔ 不另开一票）**：① `PanelAPI` 的**方法名封闭集**尺（新增一枚能改变卡状态的方法名 ⇒ 指名红；M-E 那一形必须被它抓到）；② `PanelItem` 的**字段封名单**尺要**按能力**判（能改变状态／能携带凭据形状的字段名怎么改都算——M-D2 那一形必须被抓到，⛔ 只认字符串名的尺不许当凭据）；③ 语义侧尺（出向读面拿到的东西**不可**回填成答复）。⛔ 这三枚只加尺，**不许为此解冻乙形**（`242-v1` 已判"缺的零件可在甲内补"）。
- [ ] **AC#4 载具前置**：`cmd/wisp/subagent_selfapproval_197_test.go:109` 那枚 `TaskID == CorrelationID` 要拆成两个不同 id（票 242 AC#1 的载具前提）。⚠ 写面＝`cmd/wisp` ⇒ 与 `255-r2`／`257-r1`／`167-r2`／`253 AC#1` 逐枚**串行**。
- [ ] **AC#5 越界检查**：`git diff` 出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` 任一路径 ⇒ 直接退回。

## 禁区（本票全程）

- ⛔ **不许为了变绿放宽任何既有断言**；⛔ 不许 `t.Skip`；不许把 SKIP 读成通过；⛔ 不许把"仪器绿"说成"能力已接上"（判据＝生产调用点枚数或整包颜色，`build`/`vet` 的 rc=0 不算——`ticket242_binding_test.go:111` 那枚 `sameDigest` **零调用者而包仍 ok** 就是现证）。
- ⛔ 三枚冻结件一字不动：`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`（**读**内部判射程允许，要具名写"射程判断，非内容引用"）。
- ⛔ 凭据值绝不进对话／日志／表（只写变量名）。⛔ `frontend/**`／`design/**` 两层禁令（不读、结论里也不引）。
- ⛔ **未定义即停**：AC#1 选边属编排者裁定（票面已把两支都摆出来，⛔ 腿不许自挑）；碰到 D43 转移表或 C17 名册要动 ⇒ 停手上报，那属人工批准（`SPEC-12 §4.1`）。
- Git：只 commit 不 push；commit 必带**显式 pathspec 且写在 commit 命令上**（`git add -- <paths> && git commit -F <msg> -- <paths>`，中间不停顿）；⛔ `add -A`／`.`／`-a`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／merge／worktree；临时件只建不删。

## 排程与互斥

写面＝`internal/agent/approval`（＋AC#4 的 `cmd/wisp`）。⚠ **`internal/panel` 与 `cmd/wisp` 编译都会带上 `internal/agent/approval`** ⇒ 本票的写腿在跑突变期间**与 `255-r2`（cmd/wisp）／`253-r5`（internal/panel）互斥**：要么排队，要么只加尺不种产码突变。AC#0 是只读普查，可先派（⛔ 禁跑 `go build`/`vet`/`test`，与写腿同机即互洗）。
**当前队列位次**：排在 `255-r2` 与 `253-r5` 交件之后（AC#0 只读腿可提前）；母票 242 的三格**保持未勾**、凭据归本票与 `242-v1` 那份表，⛔ 不许在 242 面上勾。
