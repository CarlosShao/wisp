# 票 64 对抗验收（非实现者执行）— 热键接线 / 交互四项 / Sleeping 零定时器实测 / 多显示器

**验收人：** 编排者（**不是**票 64 的任何实现代理；段 1/段 2/接续代理三批实现我一行没写）
**时间：** 2026-09-20 21:16–21:19Z（本地 21:16–21:19 CST 段）
**被验对象 HEAD：** `2825f28`（票 64 的实现 commit 落在 `901f334`/`0b6ae92`/`2d063f0`/`307e24e`/`620f565`/`7d50e17`）
**桌面占用状态：** 验收期间**桌面被票 66 的测量代理独占**，因此本表把「我此刻能独立复跑的」与
「必须等桌面空出来才能复跑的」**显式分开**。凡我只做了静态读码、没实跑的，判据列里写明
`静态读码`，**不冒充成实跑**。这条声明本身就是本表的可信度前提。

## 裁决表（与票面 AC 1:1，README 规则 6）

| # | AC（票面原文摘要） | 票面勾选 | 裁决 | 判据（我核到的东西） |
|---|---|---|---|---|
| 1 | 改 `[hotkey]` 后热键真的重注册（端到端：改配置→按新键→球响应；旧键不再响应） | 未勾 | **正确未勾** | 机制侧我读码确认存在且非死代码：`internal/ball/hotkey_reload.go` 桥 + `cmd/balldebug -config` 宿主。**未勾的理由成立**：全仓 `ball.New(` 只在 `cmd/balldebug/main.go`，`cmd/wisp` 走一次性 `config.LoadFile` ⇒ 生产里没有触发者。这正是 registry A12/R12 的第四次同形复现，代理**没有**用「机制已证」冒充「端到端已证」，判 PASS 的是「诚实度」而不是这一格。 |
| 2 | 注册失败两类语义分开：被占用 vs 未尝试，各自用户可见提示 | 已勾 | **PASS（我实跑）** | 读码：`hotkey_windows.go:377-413` `registerAllWith` 以 `errors.Is(rerr, windows.ERROR_HOTKEY_ALREADY_REGISTERED)` 判**真 1409**→`HotkeyTaken`，其它错误→`HotkeyError`，空串→`HotkeyDisabled`（从未尝试），语法错→`HotkeyUnparsable`；`Attempted()`（:261）把两族划开；`Problems()`（:324-345）为**每一族出不同句子**，`HotkeyTaken` 那句明说"occupied by another program and was NOT registered"。我实跑 `go test -count=2 ./internal/ball/ -run TestRegisterAllSplitsFailureFamilies` → **2/2 PASS（0.00s）**，该用例走 `registerFn` 注入接缝，**不依赖桌面**，所以分类回归在默认套件里就被钉住。 |
| 3 | 默认唤起键为 `Ctrl+Alt+Q`，且 `Ctrl+Alt+Space` 作为可配置备选**写进文档** | 未勾 | **正确未勾（默认值半边已证）** | 我实跑 `-run TestDefaultHotkeys` → **2/2 PASS**；`hotkey_test.go:63` 断言 `d.Summon` **既不是 W 也不是 Space**（负向也钉了，这点好），`hotkey_status_test.go:276-280` 断言 `AltSummonSpace=="Ctrl+Alt+Space"` **且能 parse**。文档半边做不了是**契约层事实**：含 `[hotkey]` 的 SPEC-03/SPEC-08/PLAN 全在本票冻结清单内。**但它不该悄悄躺着**：见下面 MINOR-3，我把「写到哪」这一决断收回给我自己，不留给"待解冻"。 |
| 4 | 交互四项各有真机测试且**不被 skip** | 已勾 | **PASS WITH EVIDENCE GAP**（见 MINOR-1） | 静态：四项确有四条实测试（`interaction_live_test.go` 单击唤起 / Confirming 取消 / 不抢焦点 / 透明区穿透），且全仓 8 处 `t.Skipf` **无一例外带 `SKIP-LOUD:` 前缀并说明"这一跑没证明什么"**，`requireQuietBallDesktop`（`live_guard_windows_test.go:160-168`）是 `t.Fatalf` 且原文写着"t.Skip is not an option here"——**静默 skip 这条确实被堵死了**。⚠ 缺口：代理日志对 A1/A1d 两条明写"未走 skip 分支"，对**交互四项没有逐条这一句**；`-count=2` 的 0.062s 是**非 winlive** 套件时长，winlive 那句"14 项全 PASS"在 Go 的语义里**包含 SKIP 也算 ok**。所以"不被 skip"今天是**结构上成立、逐跑记录上缺证**。 |
| 5 | Sleeping 零活动定时器句柄为**实测断言**（winlive），策略表断言保留但不作唯一证据 | 已勾 | **PASS（本票质量最高的一格）** | `TestLiveSleepingZeroTimerHandles`（`hotkey_live_test.go:~290-321`）不是读策略表，是**子类化真窗口数 `WM_TIMER` 实收消息**，600ms 内 >0 即 `t.Fatalf`。**更难得的是它自带正/负对照**：`:309` "探针一条消息都没收到 ⇒ Fatal，说明子类没装上"，`:321` 再切 `Warm` 要求探针**必须看得见**真定时器。这正是我一路要求的"先证明仪器在响"，代理自发做了。⚠ 我这次**没实跑**（桌面被占），判据是读码 + 代理日志；解除条件见 MINOR-2。 |
| 6 | 多显示器：真拖到第二屏验证并留证 | 未勾 | **正确未勾** | 本机物理单屏（`\\.\DISPLAY4` 3440x1440）。**关键加分项：代理没有用双屏 mock 冒充**，并写清了解除条件（双屏硬件 + 给 `enumMonitors()` 开注入接缝）。这符合票面"不许用双屏 mock 冒充"。 |
| 7 | 测试前置：检测到外部同窗口/演示进程时 fail-fast 并提示，禁止静默 skip | 已勾 | **PASS** | 同上 MINOR-1 判据；且这条把 A5 那两起污染（`FindWindowW` 按类名匹配到演示球 → 假红；"四个热键全失败"其实是别人的进程占键）**写成了回归纪律**，实况断言改走 `b.DebugHWND()` 而非类名（`307e24e`）。这条我认可为**已解决**，不是"已登记"。 |
| 8 | registry A1–A7 逐条标注"已修/待硬件/移交票 xx"，不许无声消失 | 已勾 | **PASS** | `docs/reports/pending-and-issues.md:163` 起有逐条块，A1/A1b/A1c/A1d/A2/A3/A4/A5/A6/A7 各有去向，且 A1 明写"移交票 12 + 对应框**保持未勾**"——移交带的是**未勾的框**而不是"算完成"。A6 由我独立复跑定案（一次性 D2D/COM 底座，非逐球泄漏 → 不立票）。 |
| 9 | 对抗验收由非实现者执行，报告含与本表 1:1 的裁决表 | 未勾 | **PASS（本文件即闭合）** | 本表 9 行 = 票面 9 框，逐号对齐。 |

**总计：1 BLOCKER=0 · MAJOR=0 · MINOR=3。票面 4 个未勾框全部**判为"正确未勾"**（无一例是把没做的说成做了）。**

## MINOR（不得无声消失）

- **MINOR-1｜AC#4 的"不被 skip"缺逐跑记录。** 闭合判据：桌面空出后跑一次
  `go test -tags winlive ./internal/ball/ -v -run 'TestBallLive|TestLive'`，把 `--- PASS / --- SKIP`
  **逐条贴进票 64 的 Progress log**；只要交互四项里有任一条走 SKIP，**AC#4 的框必须退回未勾**（这条是我自己的权限，
  不是代理的）。命令与判据都写进 `docs/evidence/s1/64-*`。
- **MINOR-2｜AC#5 我这次只做了静态核。** 解除 = 与 MINOR-1 同一跑里 `TestLiveSleepingZeroTimerHandles` 实 PASS。
  （它自带正/负对照，所以我对实跑结果没有疑虑，缺的只是"这一跑真发生过"。）
- **MINOR-3｜AC#3 的文档半边不该只挂"待解冻"。** 冻结的是 SPEC，不是所有文档。
  我把它改挂到**有明确落点的非冻结文档**（`docs/BUILD.md` 或新建 `docs/CONFIG.md`）+ 票 39（配置 GUI）双轨，
  由我拍板落点，代理不必等 owner。

## 我在读码时另挖到的一条潜在缺陷（票 64 范围外，待我复跑变异）

`internal/ball/hotkey_windows.go:382-386` 的 `registerAllWith` 用**位置 zip** 把配置值配到 id 上：

```go
bindings := []string{cfg.Summon, cfg.Mute, cfg.Cancel, cfg.Panel}
for i, p := range hkNames {           // hkNames: {hkSummon,"summon"},{hkMute,"mute"},...
    b := HotkeyBinding{Name: p.name, ID: p.id, Binding: bindings[i]}
```

`hkNames` 与那个字面量切片**靠下标对齐**，且**全仓没有任何用例钉住这个对应关系**
（`grep hkNames --include=*_test.go` → 零命中）。两个方向后果不同：
加第五个热键而忘了扩字面量 → `bindings[i]` 越界 **panic（响的，可接受）**；
但**调 `hkNames` 顺序** → `cfg.Summon` 的字符串会以 `hkMute` 的 id 注册，
`Problems()` 仍报"summon live"，而按下 summon 什么都不发生或去切静音——**静默错位**。
这与 registry C-3（`firstStringParam` 只看第一个命中键）、M-7（按键名决定扫不扫）是**同一族**：
安全/正确性判定建立在"两个平行数组的下标恰好对齐"上。
**建议修法**：合成一张具名单表 `{id, name, cfg-field}`，消掉 zip（**不是**加断言保留 zip）。
**待我做变异检验**：把 `hkNames` 前两项对调，若默认套件与 winlive 全绿 → 证实无人看管；
今天不能跑，因为票 66/AC#7 代理正在同一包上跑 `go test`，我改源码会让它读到我的变异。
**登记去向：票 64 不改**（本票已接近闭合、且它声明不碰材质以外的球码），转新票或由我在票 64 收尾时一并修，
由我在 registry 记一条 A16 后定。

### 更正（21:25，同一轮内）：上面那句"全仓无一用例钉住该对应"**是错的**，我把结论收回

我只做了 `grep hkNames --include=*_test.go` → 零命中，就推断"没人看管"。**字面命中为零属实，推断不成立**——
保护不在 `hkNames` 这个词上，而在别的断言里。逐条读完测试体后重算 6 种两两对调：

| 对调 | 抓住它的现有断言 |
|---|---|
| summon↔mute | `hotkey_status_test.go:120-126`：只配 Summon，断言 `hkMute/hkCancel/hkPanel` **全不 live** |
| summon↔cancel | `:225` / `:241`：`rep.Binding(hkSummon).Acc.VK == 'R'`（后 `'Q'`）——**配置串必须落在对应 id 上**，这就是配对断言 |
| summon↔panel | 同上 `:225` + `:111` panel 期望 Unparsable |
| mute↔cancel | `:103` 期望 `hkMute` 为 `HotkeyTaken`，换序后 `hkMute` 收到 `"Esc"` 会注册成功 → 红 |
| mute↔panel | 同上（`:106` 的 1409 断言搬到错的 id） |
| cancel↔panel | `:98` `IsLive(hkCancel)` → 换序后 `hkCancel` 收到垃圾串 → 红 |

⇒ 严重度**下调**：从"静默错位的潜在缺陷"改为「**保护是副作用而非意图**」。四槽状态恰好各不相同
（Live / Taken / Live / Unparsable）加两条 `.Acc.VK` 断言，偶然把六条边全兜住了；
但改测试的人不会知道自己在拆配对保护，而**加第五个键**（票 39 配置 GUI 必然做）时若五槽状态相似、
又没补新的 `.Acc` 断言，错位就真的会变静默。修法方向不变（具名单表消 zip，不接受"保留 zip + 加断言"），
**优先级改为与票 39 的第五键同批**。变异检验的预期结果随之**反转**：对调前两项**应当转红**；
跑红即按文档性 MINOR 收尾，跑绿说明我这次更正又错了一层，升回原严重度。

**这条更正本身值得留下**：我要求每个代理"读测试体而不是读测试标题"，自己却在一次 grep 零命中后就写了断言。
判据：**"符号名 grep 零命中" ≠ "没有保护"**。

## 我对这张票的总体判断

三批代理接力、中途死过一次（`cmd/balldebug` 被留在不可编译状态，我回退该文件后另存检查点），
最后**把"机制"与"端到端"分得清清楚楚**：AC#1 宁可空着也不借 balldebug 的证明冒充生产可达，
这一点是本轮所有票里我最看重的行为，**记为正面证据**。
