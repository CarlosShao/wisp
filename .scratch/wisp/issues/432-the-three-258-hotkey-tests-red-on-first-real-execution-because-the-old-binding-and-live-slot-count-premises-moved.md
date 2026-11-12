# 432 — `Test258*` 三枚在 round 89 首次跑到即红：热键重绑的"旧绑定存活"与"活槽计数"两个前提都动了

- Status: **ready-for-agent**（先只读取证；⚠ 预期台面/环境件，但**由你判清后自己定并写理由**）
- 来源：round 89 整包（`D:\tmp\wisp\20261112\round89\full.log`，`PASS=583 / FAIL=11 / SKIP=11`）
- 关联：票 258（热键重绑宿主，三枚测试）· 票 246（同族 resident 面）· 台账 `A1806`

## 这是什么（人话）

有三个测试在测**热键重绑**（把"唤起球"从一个键换到另一个键）：
① 换键时如果新键装不上，**旧绑定必须活着**（不能两边都没了）；
② 换完之后**活着的热键槽位数量**要对（唤起＋静音＝2，取消是备胎、面板被占）；
③ v1 探针那枚要证明**换完之后真的重绑到了活的球上**。

现在三枚都红，红句都指向**"前提动了"**——也就是测试假设的初始状态已经不成立。

## 现量（逐字，round 89）

### ① `Test258BridgeMutationNoSrcKeepsOldBinding`（1.17s）
> `resident_hotkey_258_windows_test.go:361: the old binding did not survive a no-rebind host (VK=0x0): the premise moved`

### ② `Test258OccupiedCombinationNamesTheNewValue`（0.80s）
> `resident_hotkey_258_windows_test.go:448: the boot verdict lost the live-count shape: "the floating ball window is up in this process (tray icon added, hotkeys f…"`
> `resident_hotkey_258_windows_test.go:453: live slots after the rebind = 1, want 2 (summon+mute; cancel is standby, panel is taken)`

### ③ `Test258V1ProbeSummonEditRebindsLiveBall`（0.83s）
> `resident_hotkey_v1probe_test.go:108: 258-v1 PROBE: boot verdict lost the live-count shape: "the floating ball window is up in this process (tray icon added, hot…"`

## ★历史（重要，⛔ 别当成新退化）

- **round 88**（`D:\tmp\wisp\20261112\fullclean\cmdwisp.log`）：`Test258*` 只跑了**四枚**且**全 PASS**
  （`AssemblyRootWiresTheChainAndTheBridge`／`SinkRebindRulerReadsBothSpellings`／`ConstructionChainTakesConfigValues`／`ConstructionChainMissingFileFallsBackAndSaysIt`）。
- **这三枚在 round 88 压根没跑到**（环境遮蔽：`WISP_SPIKE_MODELS is not set` 那类早退）。
- ⇒ **round 89 是它们第一次真正执行**。⛔ 不许说"258 退化了"。

## ⚠ 必须自己判清的核心问题（⛔ 不许我替你定）

**"前提动了"具体是什么动了？三个可能，修法完全不同：**

- **甲：环境件**——这三枚要**真桌面/真热键注册**（`RegisterHotKey`）。若本机有别的程序占了那个键，
  或测试跑时**这台机有人在用**（前台窗口是别人的），`RegisterHotKey` 会失败 ⇒ **前提自然不成立**。
  判据：看红句里的 `VK=0x0`（**键值 0 = 注册失败**）与 `live slots = 1`（少了一槽）。
  ⇒ 若是这个，修法是**台面侧**（换一个更冷门的键、或检测到占用就 skip 并说明）。
- **乙：产品退化**——热键注册/重绑逻辑真的坏了（例如票 246 的改动波及）。
  判据：去查最近动过 `resident_hotkey*` 的提交，用 `-overlay` 回退到旧版跑。
- **丙：断言过严**——"活槽计数＝2"这个数在当前设计下本就不该是 2。
  判据：读票 258 票面与产码，确认设计上到底该有几个活槽。

**你必须先跑一次 `git log` 看 `resident_hotkey*` 最近改动，并确认这三枚在 round 88 之前是否绿过**（`git log` + 若有更早的整包日志就去查）。

## 要做的（按序）

- [ ] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. 三枚各**复跑 ≥2 次**：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^Test258'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. **逐字抄出三枚红句全文**（日志里被截断的要回源码抄）。
  3. `git log --oneline -10 -- 'cmd/wisp/resident_hotkey*'`，看最近谁动过。
  4. 读 `:361`／`:448`／`:453`／`v1probe:108` 四处断言的**前置条件**，写清它们假设的初始状态是什么。
  5. 判清甲／乙／丙，给一句话结论＋`文件:行`。
  6. **⚠ 特别确认**：这三枚**是否要求独占桌面前台/热键**？若是，本机"有人在用电脑"时是否必然红？

- [ ] **AC#1 按 AC#0 的判定修（只修被判定那一处）**：
  - **甲（环境件）**：修台面——换冷门键、或**检测到占用就 `t.Skip` 并逐字说明原因**（⛔ 不许静默跳过）。
  - **乙（产品退化）**：修产码。
  - **丙（断言过严）**：改断言，⛔ 不许放宽成"≥1 就算过"。
  - **必须配突变控制真跑**：删掉修好的那处 ⇒ 必须回到红；恢复 ⇒ 绿。

- [ ] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [ ] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`，独立 module）。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是未跟踪残留件，会让 `go build ./...` 红。**⛔ 不许删**。

- [ ] **AC#4 与票 258／246 的关系写清**：本票只处理这三枚，⛔ 不翻票 258／246 任何框。

## 本票**不**解决

- 不碰 `TestAC4Focus*`／`TestAC246*`（同为桌面焦点族，另计）。
- 不碰 `TestAC1ResidentLeg*`（票 431 在飞）。
- 不改 `thresholds.go`／SLO 阈值（冻结件）。

## 派单约束

- 先只读（AC#0 零写），判定后再定写面。
- ⛔ `git add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；`add`／`commit` 两处显式 pathspec。
- 只 commit，⛔ push。
- ⛔ 翻任何 AC 框。
- 临时件只建不删，落 `D:\tmp\wisp\<YYYYMMDD>\432-<腿码>\`。
- **仪器**：运行时行为检查 ⇒ `go test -overlay`（拷贝放**仓外**）有效。

## Progress log

### 2026-11-12 22:0x · 写腿 `432-hotkey258` · 判**甲（环境件）**，修台面一处，产码零改

**结论一句话**：三枚红在**前提缺失**（Q／M／P 当时被别人占着），不在产品、不在断言。修法＝检测到占用就 `t.Skipf` 并逐字说明（⛔ 没选"换冷门键"）。证据件：`docs/evidence/s1/432-hotkey-258-premises.md`。

**★先纠了票面一个前提（重要）**：`## ★历史` 说"round 88 只跑了四枚且全 PASS、这三枚没跑到"。**现读 round 88 整包日志（`D:\tmp\wisp\20261112\fullclean\cmdwisp.log`，19:09:19）＝11 枚 `Test258` 全跑、11 枚全 PASS，三枚在册**（`=== RUN` 11／`--- PASS: Test258` 11／`--- FAIL: Test258` 0）。⇒ 这三枚**不是"首次真跑"**；round 88 绿、round 89 红，**两轮之间零提交**（最近动 `cmd/wisp/resident_hotkey*` 的是 11-08 `eec09a028`）。**结论（不是退化）恰好对，但票面给的依据不成立**，建议编排者更正那一节。

**甲／乙／丙**：判**甲**。
- 甲成立：`internal/ball/hotkey_windows.go:504-509` 把"别人占键"判成 `HotkeyTaken`；round 89 现数 **32 条** `occupied by another program, not registered`（M×12／P×12／Q×7／V×1 后者是 ② 自己的 squat），逐字 `err="Hot key is already registered."`。占位者＝**431 腿 21:22:48 起的常驻进程**（`431-ac1budget/boot1/logs/wisp-20261112-001.jsonl` 逐字 `hotkeys_live:3 ... summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, panel=... live`），而三枚在 21:23:32–21:23:38 跑。
- 乙排除：两轮之间零提交＋同码 19:09 绿 21:23 红；且 round 89 里**产品行为是对的**（老实报 occupied、`hotkeys_live=0/1`）。
- 丙排除：`:447` 要 `3/4` 与产码一致（`resident_ball_windows.go:526-527` 印 `len(rep.Live())`，idle 按票 245 就是 3）；`:452` 要 2 也与设计一致（Z+M live／cancel standby／panel 被 squat）；**干净桌面上四发全绿**＝断言不过严。

**★票面点名要答的那格：要不要独占桌面前台／热键？**
- **要「热键注册」（`RegisterHotKey`），不要「前台」。** 断言只读 Win32 注册集与启动时算好的 `len(rep.Live())`，`git grep 'SetForegroundWindow|GetForegroundWindow|SetFocus'` 在这两枚测试文件**零命中**（命中全在 `panel_*`／`*winlive*` 族）。
- 所以「有人在用电脑时是否必然红」＝**否**，取决于有没有人占住 Q／M／P：没占⇒照绿（本次 run1／run2／run4／run5 就是机主在用、四发全绿）；占了⇒必红（round 89 与我的 M1）；占了且修复在盘上⇒skip 并说明（mut-control-1）。

**突变控制（真跑，仓外 `holder.exe -keys Q,M,P` 在锁线程持三枚键，只持键不抢前台）**：

| 台面 | 三枚读数 |
|---|---|
| 修复前／桌面干净（run1、run2） | PASS 1.19s／0.81s／0.74s；PASS 1.19s／0.84s／0.78s |
| 修复后／桌面干净（run4、run5） | PASS 1.17s／0.78s／1.05s；PASS 1.24s／0.75s／0.75s |
| **修复后＋holder（mut-control-1）** | **SKIP 1.31s／0.74s／0.92s**，三条 `SKIP-LOUD` 逐字点名占用的槽与键 |
| **修复删掉＋holder（M1）** | **FAIL 1.11s／0.69s／0.79s**，三句与 round 89 同形（`VK=0x0`／`live 1/4`／`live slots=1`） |
| 恢复修复＋holder | SKIP 0.68s／0.73s／0.74s |

⇒ 有牙：删掉修好的那处，三枚**精确回到** round 89 的红；恢复即回 skip。

**改了什么（只台面，产码零改）**：`cmd/wisp/resident_hotkey_258_windows_test.go`（+62：补 `hkMute258=2` 常量、新增 `skipIfOccupied258`、两处调用）＋`cmd/wisp/resident_hotkey_v1probe_test.go`（+7：一处调用）。skip **只**对 `ball.HotkeyTaken`（`hotkey_windows.go:506` 那条 1409 分支）生效；`HotkeyDisabled`／`HotkeyUnparsable`／`HotkeyError`／键活着但值不对 ⇒ **照旧红**，产品退化躲不过它。

**门禁**：`go build ./cmd/wisp/` rc=0 · `go vet ./cmd/wisp/` rc=0 · `gofumpt -l` 零输出 · `scripts/d22scan.sh` rc=0（`clean - no D22 ban violations`，runtests `PASS=45 FAIL=0 SKIP=0`）。`build/mut_rvw.go` 未删。

**越界**：本腿只碰上面两枚；`docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`d22scan/allowlist.txt`／`frontend/**`／`design/**` 在**本腿名下**零命中（`git diff --name-only` 里其余 9 枚是别人在飞的件）。

**并发事故如实记**：21:58:04 撞上 431 腿正在改 `cmd/wisp/resident_sink_nail_127_windows_test.go`（共享编译面）的半成品，`go test -c` 报 `undefined: waitForFact127` 两处；非我改动引起（我两枚当时 build/vet/gofumpt 全绿），等它补上后复跑即绿。按"`cmd/wisp/**` 每波只许一枚写腿"，**排程问题，我只记录不裁定**。

**没做**：没改产码／没碰 `TestAC4Focus*`／`TestAC246*`／`TestAC1ResidentLeg*`／没翻任何 AC 框／没 push。AC 框一律留空。

