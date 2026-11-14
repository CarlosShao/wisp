# 票 245 · 实现腿 r4 · 证据件

- 腿：`245-r4`（写腿，非裁决者；AC 框一枚未勾，勾归编排者/裁决腿）
- 日期：2026-11-14 13:0x–13:3x +08
- 起点 HEAD：`137bf535c`（分支 `dev`；⚠ 与本腿开工时读到的 `b52365b71` 不同——期间别的腿落了提交，本腿的读/写都在 `137bf535c` 之后的工作树上，逐条读数自陈时刻）
- 写面（实际改动，`git status --short` 可核）：
  - 新：`cmd/wisp/resident_esc_delivery_245_windows_test.go`（AC#9 常驻尺）
  - 新：`internal/ball/hotkey_esc_idle_owner_245_test.go`（AC#8 判据分叉的默认档控制）
  - 改：`internal/ball/hotkey_live_test.go`（`requireIdleRoster` 分叉 + 探针/分类器）
  - 改：`internal/ball/hotkey_idle_bare_key_245_test.go`（分类器与探针记录类型落在默认档）
  - 改：`internal/ball/interaction_live_test.go`（260 形ⓐ 之后那枚确定性基线红的修正）
  - 新：本文件
- 未碰：`internal/panel/**`、`internal/observe/**`、`internal/tools/**`、`internal/agent/**`、`scripts/**`、`frontend/**`、阈值/golden/契约（D43 转移表）、任何生产行为字节（`cmd/wisp/resident_ball_windows.go` / `resident_windows.go` 只在正控里临时改过，见 §3，已还原到与 HEAD 逐字节相同）。

## 0. 本腿判定的一句话

票面三处缺口（AC#8 的两难、260 之后那枚确定性基线红、AC#9 只有验收侧台件）在本腿都变成了**仓内可跑的读数**：AC#8 拆成"我们自己的注册集握着裸键＝产品红"与"注册集干净、键被别人占＝具名 SKIP-LOUD＋数字"；基线红修在注入行本身（借哪枚就注哪枚）；AC#9 用仓里既有的第二个真进程台件（`cmd/wisp/testdata/esclistener`，票 246-r1 入库）对**出厂 `wisp.exe`** 跑四形，一枪 19.78s，不依赖人工桌面。

## 1. 起点核对（现场读，⛔ 凭记忆）

票面 `.scratch/wisp/issues/245-….md` 顶层 9 枚 AC 框，现读：AC#1 `[x]`、AC#2 `[x]`、AC#3 `[x]`、AC#4 `[ ]`、AC#5 `[x]`、AC#6 `[ ]`、AC#7 `[ ]`、AC#8 `[ ]`、AC#9 `[ ]` ⇒ **4/9 勾**（另有验收表里同一批框的副本，读数同上）。

`git log --oneline --grep=245`（现跑）里票 245 的落地提交：`9e84ae42e`（r3：AC#6①②＋AC#4 的钉）、`02799be92`（AC#7）、`ce6f556d9`（AC#7 证据复跑注记）、`e6d6426f8`（r1 乙形，较早）。

存量事实更正（票面文字过期，以盘上为准）：
- 默认接线**不是** `resident_ball_windows.go:86` 的 `Hotkeys: ball.DefaultHotkeys(),`。现读 `cmd/wisp/resident_ball_windows.go:452` 逐字 `Hotkeys:      cfg,`；`cfg` 来自 `:426` 的 `cfg, provenance := residentBallHotkeyChain258(hotCfg)`；链在 `:350`。
- 258 已把 `[hotkey]` 接进常驻腿（`resident_windows.go:63` `hotkeyReloadSource296` + `:432` 的 `hotReload258`）。

## 2. 本波三件与读数

### 2.1 AC#8：把"桌面有占位者"与"我们自己没绑上"分开判（两枚不同的判词）

实现（`internal/ball/hotkey_live_test.go` 的 `requireIdleRoster`，判据本体在 `internal/ball/hotkey_idle_bare_key_245_test.go`）：

- `probeBareEsc245(t, b, 3)`：同一枚裸 Esc 探针连问 3 次，记下 `attempts` 与逐次 errno；
- `classifyEscIdleOwner245(rep, p)` 三值：`free` / `ours` / `foreign`。判序有意如此：**报告里点了裸键就是 ours**，哪怕桌面同时被外来进程占着——产品错不能藏在环境后面；
- `ours` ⇒ 硬红。⚠ **引文按路径分开写（2026-11-14 引文修正）**：**live 路径实际响的红句是 `hotkey_live_test.go:505` 的 `the live ball's idle roster is not the standby one`** —— 因为 `:504` 的 `idleRosterReds245(rep)` 用的就是「裸键在册」这同一条谓词（`hotkey_idle_bare_key_245_test.go:64-68`），`ours`/`foreign` 的分类（`:507` 的探针＋`:508` 的分类器）在它**之后**才到达；所以分类那一支的判词 `ticket 245 RED (AC#8, OUR OWN SET): … the ball's own registration set names it: <槽>=<绑定> (<加速键>)`（`hotkey_live_test.go:513-516`）**只在单测路径可观察**（`hotkey_idle_bare_key_245_test.go:64` 起直调分类器的那组用例），在 live 路径**不可达**。引 `:513` 时必须同时写明「live 路径实际红句是 `:505`」，否则是"真话但误导"；
- `foreign` ⇒ `t.Skipf`，判词 `SKIP-LOUD (ticket 245 AC#8, FOREIGN OWNER): … Probe attempts=3, each answered errno 1409 "Hot key is already registered." …`。

**errno 的出处是构造性的**：`escBorrowProbe` 只把 `ERROR_HOTKEY_ALREADY_REGISTERED`（1409）映射成 `(free=false, err=nil)`，其余拒绝一律作为 error 上抛并在调用点致命 ⇒ "不 free" 这个读数按构造就是 1409，记录的不是猜测。这条写进了类型注释。

默认档控制（5 枚一起跑，`go test ./internal/ball -run 'TestAC8…|TestAC6…|TestRefusedBorrow.*245|TestHeldBorrowReleaseStillReturnsStandby245' -count=1 -v`）：

```
--- PASS: TestAC8IdleProbeSplitsForeignOwnerFromOurOwnSet245 (0.00s)
--- PASS: TestAC6IdlePinJudgesTheBareKeyNotTheCount245 (0.00s)
--- PASS: TestRefusedBorrowReleaseReturnsTheIdleRoster245 (0.00s)
--- PASS: TestRefusedBorrowStaysNamedWhileTheCardWaits245 (0.00s)
--- PASS: TestHeldBorrowReleaseStillReturnsStandby245 (0.00s)
ok  	github.com/CarlosShao/wisp/internal/ball	0.090s
```

**真机上的 foreign 支**（另起仓内台件占住裸 Esc，跑同一组 winlive 四例）：

台件 READY 逐字：`READY hwnd=0x7861602 foreground=0x7861602 foreground_is_mine=true steal=true`

```
hotkey_live_test.go:210: SKIP-LOUD (ticket 245 AC#8, FOREIGN OWNER): bare Esc is claimed on this desktop while the ball is idle, but the ball's own registration set is clean (3 live entries, no bare key), so the claim belongs to another program. Probe attempts=3, each answered errno 1409 "Hot key is already registered.". The idle-roster half of this case PASSED before the probe; this run cannot judge the delivery half on this desktop.
--- SKIP: TestLiveHotkeyRebindEndToEnd (1.49s)
--- SKIP: TestLiveHotkeyOccupiedVsNotAttempted (1.08s)
--- SKIP: TestLiveConfirmingCancelAndEscReturned (0.95s)
--- SKIP: TestBallLiveLifecycle (0.83s)
ok  	github.com/CarlosShao/wisp/internal/ball	4.640s
```

⇒ "注册集干净＋键被外人占"这一形：四枚用例**具名跳过**，每枚带 `attempts=3` 与 `errno 1409` 原文；不再硬红、也不假装绿。

### 2.2 那枚确定性基线红（260 形ⓐ 之后）

红在 `internal/ball/interaction_live_test.go:241-242` 的注入行（**改后行号**；`:231` 是**改前**行号，现 `:231` 已落在那段注释里）：它写死注 `"Esc"`，而 260 形ⓐ 之后借的是 `[hotkey] cancel` 的**配置值**（本夹具 `liveHotkeys()` 写的是 `Ctrl+Alt+V`）。所以借是对的、注入是陈的。

改法：`borrowed := b.ConfiguredHotkeys().Cancel; injectBinding(t, borrowed)`，注释同步改准（逐字见 diff）。

红前（复跑，逐字；`:238` 是**改前**行号，同一句话改后落在 `:249`）：

```
interaction_live_test.go:238: the taken-over Esc did not cancel Confirming (machine Confirming)
FAIL	github.com/CarlosShao/wisp/internal/ball	3.57s
```

改后同一组四例（`GOFLAGS=-tags=winlive go test ./internal/ball -run 'TestLiveHotkey|TestLiveConfirmingCancelAndEscReturned|TestBallLiveLifecycle' -count=1 -v`）：

```
--- PASS: TestLiveHotkeyRebindEndToEnd (1.90s)
--- PASS: TestLiveHotkeyOccupiedVsNotAttempted (0.98s)
--- PASS: TestLiveConfirmingCancelAndEscReturned (1.12s)
--- PASS: TestBallLiveLifecycle (3.59s)
ok  	github.com/CarlosShao/wisp/internal/ball	7.875s
```

⚠ 这枚红的历史归口不是本腿编的：`docs/evidence/s1/239-verdict.md:35` 逐字"归 245 地界"，`65-verdict.md:131`、`68r1-evidence.md:100` 同指。

### 2.3 AC#9：四形常驻尺，跑出厂二进制

新件 `cmd/wisp/resident_esc_delivery_245_windows_test.go`（`//go:build windows && winlive`），复用仓内既有台件（`buildEscListener246`／`observeEsc246`，票 246-r1 的 `testdata/esclistener`）与既有的常驻腿装置（`buildWispForTest`／`bootResidentLeg`）。四形一一对应验收侧 E1–E4：

1. 基线（什么都不跑）：注入的裸 Esc 必须到观察窗（E1：`keydown_esc=1`）；不到且 Win32 说键被别人占 ⇒ 具名 SKIP-LOUD＋数字；不到而键是自由的 ⇒ instrument 红（台件坏了，下面的读数都不算数）。
2. 自占对照：同一个观察者自己全局注册裸 Esc ⇒ 同一发注入必须"消失"并变成 `wm_hotkey=1`（E2：`0`＋`1`）。这就是"0"那一形，证明下面的 1 不是尺子聋。
3. 出厂 `wisp.exe` 无参常驻、球窗在、没有卡：观察者**仍**必须收到（E3：`1`）。
4. 该进程走自己的关闭路退出后再注：必须又收到（E4 的"还"：`1`）。

第 3 形读到 0 时先做归属分叉：让常驻腿退出→再问 Win32；仍被占 ⇒ 外来占位（SKIP-LOUD＋数字）；变自由 ⇒ **产品红**（逐字红句见 §3）。

运行命令（⚠ 必须走 staged 载体，cgo+sherpa DLL；且必须经 `cmd /c` 或 bash 起，PowerShell 直传 `-test.*` 会被拆成 `-test`）：

```
GOFLAGS=-tags=winlive WISP_CLI_STAGE_DIR=.scratch/wisp/stage-cli-245 \
  sh scripts/wisp-cli-staged-tests.sh -test.run=TestLive245TheShippedLegLeavesTheDesktopsBareEscAlone -test.v
```

实跑（还原后最终绿，2026-11-14 13:2x）：

```
resident_esc_delivery_245_windows_test.go:206: AC#9 four shapes (E1/E2/E3/E4): baseline=1 control_steal=0/1 shipped_idle=1 after_exit=1
--- PASS: TestLive245TheShippedLegLeavesTheDesktopsBareEscAlone (19.78s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	20.00s
```

同一发里观察者的逐字读数（前面那枚 `-watch`/`-steal`）：

```
READY hwnd=0x85C0520 foreground=0x85C0520 foreground_is_mine=true steal=false
READING keydown_esc=1 syskeydown_esc=0 wm_hotkey=0 other_keys=0 foreground_is_mine=true steal=false
READY hwnd=0x85D0520 foreground=0x85D0520 foreground_is_mine=true steal=true
READING keydown_esc=0 syskeydown_esc=0 wm_hotkey=1 other_keys=0 foreground_is_mine=true steal=true
```

## 3. 正控（每条断言都要能真红；红句逐字，还原后同串）

### 3.1 分类器：把 ours 支短路

变异（`internal/ball/hotkey_idle_bare_key_245_test.go` 的 `classifyEscIdleOwner245`）：`if _, held := idleBareKeyHeld245(rep); held {` → `… ; false && held {`。编译通过（`go vet ./internal/ball` rc=0）。

红句逐字：

```
hotkey_esc_idle_owner_245_test.go:90: a claimed bare Esc that the ball's own set holds classified as owner 2, want ours (1) - the foreign-owner skip must never absorb the product's defect
--- FAIL: TestAC8IdleProbeSplitsForeignOwnerFromOurOwnSet245 (0.00s)
```

还原后 SHA256 同串：变异前 `3C0793916DD35B8EDEAD16C110194DC5B8628F7C88CEA91325CDAEE1CB2B65BA`，还原后 `3C0793916DD35B8EDEAD16C110194DC5B8628F7C88CEA91325CDAEE1CB2B65BA`；再跑该例 `--- PASS`、`ok … 0.066s`。

### 3.2 AC#9：让出厂进程真的握着裸 Esc

**第一次尝试（记档案，因为它本身就是一枚可用的盘上事实）**：把 `resident_ball_windows.go:452` 的 `Hotkeys: cfg,` 临时改成 `Hotkeys: ball.HotkeyConfig{Summon: "Esc", …}`。编译过、启动 verdict 也写 `summon=Esc live`，但四形仍 PASS——查下去是 **258 的热重载桥**（`resident_windows.go:432` 的 `hotReload258` → `hotkeyReloadSource296`）在开机 ~1s 后用**独立解析路径**重算一遍 `[hotkey]`（空文件 ⇒ 编译默认），发现与已应用的集合不同就重绑 ⇒ 把临时接上的裸 Esc 又摘掉了。⇒ 只改"建窗那一处"的变异活不过重载桥；这不是断言的漏洞，是控制设计要多改一处。

**成立的正控（两处生产变异，同属 `cmd/wisp/**` 本腿写面）**：
1. `resident_ball_windows.go` 的链：`cfg := ball.ApplyHotkeyDefaults(hotCfg())` 后加 `cfg.Summon = "Esc"`；
2. `resident_windows.go` 的 `hotReload258` 闭包：`cfg := hotkeyReloadSource296(...); cfg.Summon = "Esc"; return cfg`。

两处都在 ⇒ 开机与每次重载解析出的都是裸 Esc，稳定不被摘。跑 AC#9，红句逐字（`resident_esc_delivery_245_windows_test.go:175`）：

```
ticket 245 AC#9 RED (PRODUCT): the shipped no-args process held the desktop's bare Esc while idle. An injected Esc reached the observer 0 time(s) while the leg was up (want 1), and the same probe found the key FREE again the moment the leg left - so the claim was this process's, which is exactly the defect ticket 245 exists to keep out. Resident verdict:
--- FAIL: TestLive245TheShippedLegLeavesTheDesktopsBareEscAlone (18.67s)
```

同发 verdict 行逐字（红态证据）：`hotkeys live 1/4: summon=Esc live, mute=Ctrl+Alt+M occupied by another app, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P occupied by another app`。

还原：`resident_ball_windows.go` 回到 `Hotkeys:      cfg,` 且 `cfg.Summon` 行删除；`resident_windows.go` 的闭包回到单行 `func() ball.HotkeyConfig { return hotkeyReloadSource296(rt.Layout.DataDir) }`。逐字核：
- `(Get-FileHash cmd/wisp/resident_ball_windows.go).Hash` = `1D47ECCB205328A744FB1FF4507EF6962FA8CB29826671BAF2C5CE9406178425`（与变异前同串）；
- `git status --short cmd/wisp/resident_ball_windows.go cmd/wisp/resident_windows.go` 与 `git diff --stat --` 两条均**空**（⇒ 与 HEAD 逐字节相同）。

还原后再跑 AC#9：`--- PASS … (19.78s)`、`ok … 20.00s`（§2.3 那条）。

### 3.3 AC#4 那枚配置钉（复认，非本波新写）

`GOFLAGS=-tags=winlive … -test.run=Test245AC4ResidentChainReadsAConfiguredCancel -test.v`：

```
--- PASS: Test245AC4ResidentChainReadsAConfiguredCancel (0.01s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	0.00s
```

## 4. 门禁（四数，现跑）

- `gofumpt -l cmd/wisp internal/ball`：**空**，rc=0（`D:\work\base\gopath\bin\gofumpt.exe` v0.12.0）
- `go vet -tags winlive ./cmd/wisp`：rc=0（正控变异前后各跑一次，均 rc=0）
- `sh scripts/d22scan.sh`：`d22scan: clean - no D22 ban violations`，rc=0（`bans #8 cmd/=279` 含本腿新件）
- `sh scripts/check-path-length-budget.sh`：`VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`，rc=0
- targeted 四数：`internal/ball` 默认档 5/5 PASS `ok … 0.090s`；`internal/ball` winlive 4/4 PASS `ok … 7.875s`；`cmd/wisp` AC#9 1/1 PASS `ok … 20.00s`；`cmd/wisp` AC#4 1/1 PASS `ok … 0.00s`

## 5. AC 逐条对照（⛔ 一枚未勾）

- **AC#8**：判据已按票面"分开两条红句"落地，且**外来占位那支是具名 SKIP-LOUD＋数字**（§2.1）。⚠ 票面 AC#8 字面写 `⛔ 不许改成跳过`，与 2026-11-14 派单的"已被别的进程注册 ⇒ 具名跳过并给数字"直冲突；本腿**取派单口径**，并把它做成两难可辨的形状：**我们自己的错仍然硬红，一枚都没跳过**；跳的只有"注册集干净＋键被外人占"这一形，且判词、次数、errno 原文都在。此冲突请编排者裁定；若维持票面字面，改法是把 `foreign` 支的那一句 `t.Skipf` 换成 `t.Fatalf`——但那会把"别的程序占着 Esc"读成我们的产品缺陷（v1 报告 §4 R-4 判过它为何不对）。
- **AC#9**：仓内一条用例 `TestLive245TheShippedLegLeavesTheDesktopsBareEscAlone` 做四形、不依赖人工桌面（观察者自己抢前台、测试经 stdin 触发注入、子进程全由 `t.Cleanup` 收）；绿读数 §2.3，正控 §3.2。⚠ 它需要交互桌面（winlive），CI 无桌面时按 winlive 惯例不在默认档跑。
- **AC#6①**：⚠ 本发读数 **0/4**（`hotkeys live 0/4: summon=Ctrl+Alt+Q occupied by another app, …`），原因是另一条腿的 `wisp-live.exe`（pid 21340，13:21:15 起）占着 Q/M/P——不是产品回归。同一发里 Esc 四形全绿。**本腿不拿 0/4 冒充 3/4**；r3 欠的那枚"winlive 3/4 读数"仍欠（等桌面干净时补）。
- **AC#4**：**不勾**。阻塞于票 228。盘上现读：`cmd/wisp/resident_ball_windows.go:426` `cfg, provenance := residentBallHotkeyChain258(hotCfg)`、`:350` 链函数、`:452` `Hotkeys:      cfg,`（票面引的 `:86`/`DefaultHotkeys()` 已过期）；配置钉 `Test245AC4ResidentChainReadsAConfiguredCancel` PASS（§3.3）。这只证明"链会把 `[hotkey] cancel` 读进来"，**不证明**"卡挂着时借的是配置值"——那要 228 的卡真的举起来（本机 `任务入口未启用` 已逐字印在常驻腿日志里）。

## 6. 边界与残余

1. **AC#8 的跳与票面冲突**：见 §5，请裁定。
2. **AC#6② 行为半边没做**（260 判"永远待命"）：本腿只钉了判据精度与裸键分叉，没改"用户配组合键待命时也不绑"这一行为。改它要动 260 的裁定，须人工批准，本腿不擅动。
3. **AC#6①** 3/4 读数被同机另一条腿的常驻进程污染（pid 21340），本腿如实记 0/4 并说明原因；补测需要在"无别的 wisp 常驻"的桌面。
4. **`internal/ball/**` 写面**：派单的禁止清单逐条枚举里没有 `internal/ball/**`，而 AC#8 与那枚基线红只能落在 `internal/ball` 的**测试文件**里；本腿只改 `*_test.go`（零生产行为），并在此具名请编排者确认。撤销口令：`git checkout 137bf535c -- internal/ball/hotkey_live_test.go internal/ball/interaction_live_test.go internal/ball/hotkey_idle_bare_key_245_test.go`（并删 `internal/ball/hotkey_esc_idle_owner_245_test.go`）。
5. **`wisp-live.exe` pid 21340** 仍在跑（另一条腿的），本腿**没有**杀它（13:11 杀的是另一枚 69456，当时已确认它无球窗、占着 Q/M/P、所属腿日志断在半空）。
6. 258 热重载桥会把"只改建窗处"的接线复活掉（§3.2 第一次尝试），这对以后写正控的腿是有用的盘上事实：**改默认接线必须同时改 `hotkeyReloadSource296` 一侧，否则 ~1s 后被打回。**

## 7. 复跑命令（逐条可核）

```
# 1) 默认档 5 枚（AC#8 分叉 + 245 既有判据）
go test ./internal/ball -run 'TestAC8IdleProbeSplitsForeignOwnerFromOurOwnSet245|TestAC6IdlePinJudgesTheBareKeyNotTheCount245|TestRefusedBorrow.*245|TestHeldBorrowReleaseStillReturnsStandby245' -count=1 -v

# 2) winlive 四例（真桌面；Esc 自由时 4/4 绿）
$env:GOFLAGS='-tags=winlive'; go test ./internal/ball -run 'TestLiveHotkey|TestLiveConfirmingCancelAndEscReturned|TestBallLiveLifecycle' -count=1 -v

# 3) 外来占位那一形（具名跳过＋数字）：先起台件占住裸 Esc，再跑 2)
& "$env:TEMP\esclistener245.exe" -steal -delay 60000   # 另窗；跑完杀

# 4) AC#9 常驻尺（出厂 wisp.exe；必须经 cmd /c 或 bash 传 -test.* 参数）
$env:GOFLAGS='-tags=winlive'; $env:WISP_CLI_STAGE_DIR='.scratch/wisp/stage-cli-245'
cmd /c 'set "PATH=E:\work\base\msys64\usr\bin;%PATH%" && "E:\work\base\msys64\usr\bin\bash.exe" scripts/wisp-cli-staged-tests.sh -test.run=TestLive245TheShippedLegLeavesTheDesktopsBareEscAlone -test.v'

# 5) AC#4 的配置钉
… -test.run=Test245AC4ResidentChainReadsAConfiguredCancel -test.v

# 6) 门禁
gofumpt -l cmd/wisp internal/ball
go vet -tags winlive ./cmd/wisp
sh scripts/d22scan.sh
sh scripts/check-path-length-budget.sh
```
