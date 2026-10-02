# 票 258 — `[hotkey]` 在 schema 里被写成"hot 档：改了会重新注册"，而**常驻那条腿建球时根本没读配置**（用的是写死的 `DefaultHotkeys()`），全仓**唯一**的非测试 reloader 调用者在调试台件里

**立票时刻**：2026-10-02 15:0x，锚点 HEAD `9588138a`（`dev`）
**来路**：只读腿 `e2e-panel-1` 的第四格未定义项④（台账 `A531` §4 里我写过"待量完 `256-a1` 之后立票，⛔ 现在不自填修法"）＋ `256-a1` 已交件（`A534`）⇒ 这一枚就是那笔欠账的兑现。
**与既有票的关系**：**不并进票 255**——255 那族是"配置说了生效、没人读"，修法在**登记表与回执**；这一枚的修法在**装配**（谁把配置里的四枚热键交给球、谁在改过之后重新注册），两族的写面与判据形状都不一样。⚠ 它也不是票 245（那枚管的是"裸 `Esc` 被注册成全局热键"，射程＝绑与不绑；这一枚管的是**绑的哪一枚、改了以谁为准**）。

## 现量（编排者 15:0x 自跑，⚠ 引用前先重跑，别把这几行当常量）

1. **承诺那句话在 schema 里**：`internal/config/schema.go:176` 逐字 `// HotkeySection is [hotkey]; hot-tier (hotkeys re-register on change).`，四枚字段在 `:179` 起的 `HotkeySection`；`internal/config/manager.go:278` 的 hot 应用表里**确实有** `{"hotkey", &cur.Hotkey, &fresh.Hotkey, func() { cur.Hotkey = fresh.Hotkey }}` ⇒ **内存里的值会跟着文件变**。
2. **但常驻那条腿建球时不吃它**：`cmd/wisp/resident_ball_windows.go:171` 逐字 `Hotkeys:  ball.DefaultHotkeys(),` ⇒ 球一开始用的就是**写死那四枚**，`Config().Hotkey` 从没进过构造参数。
3. **会重新注册的那台机器只接在调试台件上**：`cmd/balldebug/main.go:237` `bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), …)`＋`:243 bridge.Refresh = func() error { _, err := mgr.CheckAndReload(); return err }` ⇒ **全仓 `NewHotkeyReloader` 的非测试调用点＝这一枚**（尺＝`grep -rn "NewHotkeyReloader" cmd internal tools --include=*.go`）。
4. **重载的引擎本身在常驻腿是活的**：`cmd/wisp/config_reload.go:153` 在跑 `rt.mgr.CheckAndReload()`，而 `:101-108` 那段注释逐字写着 `startConfigReload` 由装配根在门存在之后调用 ⇒ 所以缺的不是"重载有没有发生"，**是重载之后没人把新值交给球**。
5. **一票既有钉射程要先量清**（⛔ 本票不许为了落地去改它们）：`cmd/wisp/resident_ball_228_windows_test.go:45`、`cmd/wisp/resident_ball_live_228_windows_test.go:144/160/163` 那族读的是启动判决串里的 `hotkeys live %d/4`；线程形状受 10-01 那批裁定约束（球/面板都在自己的线程上，投 `ui-sta` 会冻外层泵——见票 33 十一裁）。

## 要建什么（本票只管"改了以谁为准"这一件；⛔ 不碰界面、不碰凭据）

- [ ] **AC#0（本票第一格，且是闸门）＝先把三问答出来，⛔ 不许直接开写**：① 常驻腿要把配置里的四枚热键交出去，**构造期**该改哪一处、`DefaultHotkeys()` 那枚默认还要不要留作"配置缺失时的那一份"；② **改过之后**谁去重新注册——是把 `ball.NewHotkeyReloader` 那台机器接进常驻腿，还是让 `config_reload.go` 那一跳多调一次注册，两形的**代价与线程约束**各是什么；③ 四枚里有一枚注册失败（被别的程序占了）时，**今天那句 `hotkeys live %d/4` 会怎么说、应该怎么说**。交件判据＝逐处带 `file:line`＋尺读数；量不到的**具名说量不到**，⛔ 不许用"应该没问题"填空；⛔ 不许改任何产码。
- [ ] **AC#1 只在 AC#0 交完、并由编排者落一枚具名 `A##` 批准选形之后才许动**：`wisp run` 与**常驻腿**两条入口建出来的球，四枚热键**以 `config.toml` 的 `[hotkey]` 为准**；配置文件缺失／那一节缺失时**退回 `DefaultHotkeys()` 并说得出这句话**（不许静默换成另一套）。
- [ ] **AC#2 说实话的判据（负向必配正控）**：种一发"把 `summon` 改成别的组合键" ⇒ **下一次建球注册的就是新值**（正控）；⛔ 反向不许做成词面尺（不许只扫注释里有没有 `hot-tier` 那个词），要问能力。⚠ 若量出来"重新注册必须回到 `ui-sta` 才安全"，那**属票 33 那批线程裁定的射程**，停下来上报、不许自填。
- [ ] **AC#3 越界检查**：`git diff` 出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` 任一路径 ⇒ 直接退回。⛔ 三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）一字不许动。

## 禁区（本票全程）

- ⛔ **未定义即停**：票 33 十一裁（面板专用 STA 线程＋库 `Run()` 泵；投 `ui-sta` 会冻外层泵）是既有裁定，动它＝**人工批准**。
- ⛔ 不许为变绿放宽任何断言；不许 `t.Skip`；不许把 SKIP 读成通过（`tools/d22scan/runtests.sh:98` 把 SKIP 判红）。
- ⛔ 凭据值绝不进对话／日志／表（只写变量名）。
- Git：只 commit 不 push；**add 与 commit 同发一条命令、commit 必带显式 pathspec**（`A532` 那枚归属事故就是这条没做到）；⛔ `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。
- `frontend/**`／`design/**` 两层禁令；`grep`/`find` 显式根（`cmd internal tools docs scripts .scratch`）。
- ⛔ **界面侧那一半不在本票**：设置页里"改热键"那个输入面属界面那支（机主自己带给他在用的那枚 agent）；本票只交 Go 侧"改了以谁为准"。**C17 白名单既有名字不动**；⛔ 不新增 D34 工具行。

## 排程

写面＝`cmd/wisp`＋`internal/ball` ⇒ ⛔ 与 `198-r2`（在飞，写 `cmd/wisp`）**串行**。AC#0 是只读普查，可先派（⛔ 禁跑 `go build`/`vet`/`test`，与写腿的整包测量同机即互洗）。
**默认不排落地腿**：本票排在票 198 与票 248 那两片交完之后。撤销口令「**258 撤**」。
