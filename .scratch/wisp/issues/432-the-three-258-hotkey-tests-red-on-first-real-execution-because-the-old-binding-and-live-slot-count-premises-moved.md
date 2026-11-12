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
