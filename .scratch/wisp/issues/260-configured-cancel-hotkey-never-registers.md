# 票 260 — 面板里那把"取消"热键**改了没用**：配置值从来不会被注册，借用那把键写死裸 Esc；而丢借用**零症状**（卡片到点必执行）

**立票时刻**：2026-10-03 09:4x +0800，锚点 HEAD `7de96e61`（`dev`）
**来路**：只读普查腿 `258-a2`（`.scratch/wisp/probes/258/a2/census.md`，606 行／55,101 字节，`3736f0dd`→`5d4f5343`）的"外延 A"；**编排者本人复量三条硬断言后才立票**（逐条见下）；关联＝票 245（裸 Esc 全局热键，已裁乙）／票 258（形 A 落地腿的邻格，⛔ 本票不是 258 的残余）。
**性质**：⚠ **算产品行为**——用户会以为那枚键能用。

## 现量（编排者本机自跑，2026-10-03 09:4x；⚠ 引用前先重跑）

1. **Win32 注册全仓产码只有一枚调用点**：`grep -rn "pRegisterHotKey" --include=*.go internal cmd | grep -v _test` ⇒ 命中 `internal/ball/hotkey_windows.go:391`（调用体）＋`internal/ball/win32_windows.go:42`（proc 声明，不是调用）。
2. **它只被两遍用到，两遍都不注册配置来的 cancel**：① idle 那遍**显式跳过 cancel 槽**（`hotkey_windows.go:458` 逐字注释 "The one slot an idle ball must not register"，配合 `:13`／`:82-85` 的设计说明）——这是**票 245 裁乙**的既有裁定（稳态不绑、确认那两三秒才借），⛔ 本票**不推翻它**；② 借用那遍的加速器**写死裸 Esc**：`hotkey_windows.go:118` `vkEscape = 0x1B`、`:521` `escBorrowAcc`（注释自述 "VK_ESCAPE with nothing…"）、借还成对体在 `internal/ball/ball_windows.go:875-901`。
3. ⇒ **推论（我复量后认它成立）**：用户在 `config.toml` 把 `[hotkey] cancel` 改成任何非 Esc 的组合，**那枚键在今天任何一条路径上都不会成键**——idle 不绑（设计）＋借用不吃配置（缺陷）。⇒ 任何"四枚热键都以配置为准"的句子**今天只对三枚成立**（summon／mute／panel）。
4. **丢借用的后果不是"少个键"，是"卡片到点必执行"**：`internal/agent/approval/gate.go:319`→`:332` 逐字 `ANSWER-EXPIRED decision=timeout->execute`（`:338`），窗口 3s＝`internal/agent/approval/queue.go:116`；且 `Standby` 不进 `Problems()`（`internal/ball/hotkey_windows.go:365`）⇒ **零症状、零告警**。

## 要建什么

- [ ] **AC#0（第一格，读数不是码）**：判死"借用那遍到底读不读 `[hotkey] cancel` 的配置值"——逐处给 `file:line`（起点＝现量 2 那四处），并列出**若要吃配置**需要哪些零件（`Accelerator` 从哪来、mods 怎么带、`escBorrowAcc` 那枚常量是谁的、改它会不会撞上"票 245 只在确认窗借裸 Esc"那条裁定的**边界**）。⛔ 本格不许改产码；量不到的那处要具名写"量不到"。
- [ ] **AC#1 只在 AC#0 交完并由编排者落一枚具名 `A##` 选形之后才许动**：两条候选形由 AC#0 的读数决定——ⓐ 让借用读配置值（则"按 Esc 取消"这项默认行为**可能随配置变**，要与票 245 那句"确认那两三秒借的那把键"逐字对齐，⛔ 不许把 245 的裁定悄悄改掉）；ⓑ **维持写死**，但把"这枚配置项今天不生效"**说到明处**（建球/回执侧具名，⛔ 不许留"四枚都吃配置"这类句子）。两形都要给"改了配置之后，用户看得见什么"那一句。
- [ ] **AC#2 零症状那一格要有声**：丢借用（或借用没成键）这件事今天**不进 `Problems()`、不打日志**；判据＝装一发**会响**的尺（种"借用请求发出但 Win32 没成键"⇒ 指名那一步必须红，⛔ 不许 `t.Skip`、不许把 SKIP 读成通过），并且⚠ 要与票 258 落地腿的 rebind×借还自钉**分开**（那是"改热键时正有借用在飞"，本格是"借的那把键根本没成"，⛔ 两格不许合并成一句"热键接好了"）。
- [ ] **AC#3 越界检查**：`git diff` 出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` 任一路径 ⇒ 直接退回。⛔ **D43 转移表一个字不许动**；⛔ 不许新增状态。

## 禁区（本票全程）

- ⛔ **票 245 已裁的那条"稳态不绑、只在确认那两三秒借"属既有裁定**：本票要动它的**边界**＝人工批准（`SPEC-12 §4.1`），编排者单方面派腿＝跑歪模式。AC#1 若量出"必须碰 245 那条"，⛔ 停下来上报，不许自己改判。
- ⛔ 不动 `internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go` 三枚冻结件；不许为变绿放宽断言。
- ⛔ `frontend/**`／`design/**` 两层禁令；⛔ 凭据值不进对话/日志/表。
- Git：只 commit 不 push；显式 pathspec 写在 commit 命令上；⛔ `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。

## 排程与串行

写面＝`internal/ball`（＋`cmd/wisp` 的回执/文案侧）。⚠ **`cmd/wisp` 此刻 `255-r2` 在飞、`internal/panel` 此刻 `253-r5` 在飞** ⇒ AC#1 落地腿**按住**；AC#0 是只读普查，可先派（⛔ 禁跑 `go build`/`vet`/`test`，与写腿同机即互洗）。
**与票 258 的关系**：258 的落地腿管"改完配置后重注册"，本票管"cancel 这枚槽根本没吃配置"——**相邻但不互相吞并**，⛔ 谁也不许把对方的格子并进自己一句结案。撤销本票＝「260 撤」。

---

## 编排者翻勾节（2026-10-03 21:5x，锚 `d8cbb948`；260-a2 交件收档）

**AC#0 翻勾**，凭据＝`.scratch/wisp/probes/260/a2/census.md`（192 行零待填，零 Go 命令）：
- **判定＝借用遍不读 `[hotkey] cancel` 配置值。** 我抽验中心凭据：`internal/ball/hotkey_windows.go:524` 逐字 `func escBorrowAcc() Accelerator { return Accelerator{Mods: modNoRepeat, VK: vkEscape} }`＝零参常量字面；借用链四枚签名零键参。成立。
- **258 桥关系判定**：链已通到借用侧**取值字段**（`ball_windows.go:124/:126` 已喂活），借用注册路径是链上唯一没接的一段 ⇒ **AC#1 若选 ⓐ（借用读配置），258 链零改动自动获益**——两票耦合比我立票时预估的浅。
- **★推翻 a1 两处**（照单入账）：①「出厂 DefaultHotkeys 在 :171、配置从未被读」已被 258-r1 推翻（现 `:269 Hotkeys: cfg`），票 260 结论收窄为"cancel 一枚"；② a1 §5.6 将来时已成本体——idle 日志 `binding` 字段已是配置值、句内 "Esc" 写死＝改键用户今天看见自相矛盾句（**AC#1 形 ⓑ 的头号必改落点**）。
- **245 边界判定**：撞不到已生效裁定（裁的是机制与两枚落点，未裁加速器来源；默认档逐位不变）；相抵的是 245 **未勾**的 AC#6②——只报不判，归人工批准面。
- **AC#1 选形**：等 260-a2 的代价表＋形 ⓑ 头号落点（idle 日志自相矛盾句）一起裁——**不急于本轮**；AC#2（丢借用要有声）与 258 的 rebind×借还自钉是两格，票面已写死不许合并。
