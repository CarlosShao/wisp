# 33-r11 impl — 票 33「*Locked 命名与自取锁一致」那一格（ⓐ 改名 ＋ ⓒ 会响的钉）

三件指认（本仓新定式，⛔ 只报号会指错）：**票 33 ／ 腿 r11 ／ 那一格的判据＝「面板宿主里两枚自取锁的方法摘掉 `Locked` 后缀 ＋ 一枚钉住命名与持锁形状一致的守卫」**。
同名不同格的另一发＝`33-r10`（票 33 AC#13「冷启动最终文档不是探测页」，件在 `.scratch/wisp/probes/33/r10/`，已交两笔）：本格与它同文件、不同判据，⛔ 不是它的续件。

- 起手 HEAD（本腿自己量的）：`03899fd8`，分支 `dev`；`probes/33/r11/` 起手不存在（`ls` rc=2），`probes/33/` 现量 r1–r9、r8b、r10、v2、v4、a1–a3、h1、n1、p1、winc1。
- 本腿 commit：`0332964e`（起手锚）→ `9995f9b1`（实现）→ `cc5f7627`（证据第一批）→ 本件所在的最后一笔。⛔ 零 push。
- 措辞死口子（本腿全文遵守）：**今天没有真死锁路 ⇒ 这一格是陷阱（trap），不是现行 bug，也不是"修了一个死锁"**。派单 §2 撤回的那一项（`cmd/wisp/config_readers_255.go` `:18`／`:146` 两处注释字串）本腿**一字未碰**；ⓑ（把取锁挪到调用点）本腿**一行未做**。

---

## ① 改名名册（ⓐ）

尺＝本腿自己跑的 `grep -rn "firstRoundTrip|serveNotBuiltNotice|Locked" cmd/wisp/`（改前／改后各一把）＋ `git diff`。

| 旧名 | 新名 | 声明枚数 | 调用点枚数 | 注释/名册引用枚数 | 为什么它不要求调用者持锁 |
|---|---|---|---|---|---|
| `PanelManager.firstRoundTripLocked` | `PanelManager.firstRoundTrip` | 1（`cmd/wisp/panel_host_windows.go:840`） | 2：产码 1（`:448`，`coldStartPageHandover` 内）、测试 1（`cmd/wisp/panel_pageover_33r10_windows_test.go:222`） | 7：host `:430`、`:797`、`:824`；pageover `:23`、`:55`、`:216`、`:238` | 方法体第 841 行就是 `m.mu.Lock()`（现量：`m.mu.Lock(); w := m.w; m.mu.Unlock()`），起手自取，调用点不持锁 |
| `PanelManager.serveNotBuiltNoticeLocked` | `PanelManager.serveNotBuiltNotice` | 1（`:460`） | 1（`:451`，同一函数内 `serveEntry()` 报错那一支） | 1（`:456` 文档注释） | 同一形状：`:461` 起 `m.mu.Lock()`／`m.mu.Unlock()` |

改后 `cmd/wisp` 里旧名零命中（唯一残留在本腿测试文件的**夹具字符串**里，那是正控要的"旧形样本"，不是引用）。

### 为什么其余 `*Locked` 一枚都不改

- `PanelManager.setPriorFocusLocked`（`:625`）：**名字是真的**——文档注释逐字写着 `The caller must hold m.mu.`，方法体零锁操作（现量 `:626-633`：只读 `prior`、比 `m.hwnd`、写 `m.prevFocus`），唯一调用点 `:525` 落在 `Show` 的 `m.mu.Lock()`（`:522`）／`m.mu.Unlock()`（`:526`）之间。它就是守卫的**豁免名册**里那一枚。
- `Manager.ConfirmLocked`（`internal/config/manager.go`）与 `cmd/wisp/config_reload.go:254 confirmLockedLoosening`、`config_reload.go:180/:206 rep.Locked`：票 33 的**同形名册**把它们记成**乙形**（后缀不是 mutex 那个意思，`ConfirmLocked` 契约逐字是 "runs OUTSIDE mu"）＝`ⓑ`／另裁一族，⛔ 不在本格射程，本腿一字未动。
- `internal/config` 那三枚**甲形**（`reloadPlan.commit`／`writeAllowedDirs`／`mergeWrite`：要求持锁但没有后缀）同属票 33 名册的另几支，本腿按派单射程（ⓐ＋ⓒ）不碰。守卫对它们**故意不响**（见 ⑥ 恒真面第 1 条）。

---

## ② ⓒ 那枚钉：判据、正控、豁免名册

新件：`cmd/wisp/panel_locked_naming_33r11_windows_test.go`（`//go:build windows`，包 `main`，纯守卫，零产码，5 枚用例）。

**判据（两支牙，各自独立会响）**

1. **名册精确（census / roster）**：`cmd/wisp` 的**生产源文件**（`*.go` 且非 `*_test.go`，本腿现量 34 枚）里每一枚名字以 `Locked` 结尾的 func／method 必须
   (a) 被 `lockedNamingRoster33r11` **按具名限定名**列出（`Recv.Name`，如 `PanelManager.setPriorFocusLocked`），
   (b) **方法体内任何位置不得出现锁操作**（`Lock`/`Unlock`/`RLock`/`RUnlock`/`TryLock`/`TryRLock`；含闭包内部——闭包里取锁只会更坏，不是更轻）。
   外加防烂两支：名册里那一枚**不存在**了 ⇒ 红（豁免名单不许烂成"没人查的清单"）；一个 `.go` 解析不了 ⇒ 红（⛔ 不许静默跳过）；目录里零枚 `.go` ⇒ 红（读不到东西的尺不算尺）。
2. **行为（behaviour）**：名册里唯一**不开窗就够得着**的那枚 `setPriorFocusLocked`，在**调用者持着 `m.mu`** 的状态下被真调用，包在**有界等待**里（`completedWithin33r11`，2 秒预算，形状抄 `internal/config/manager_223_test.go:39-51` 的 `completedWithin`，先读现量再抄），必须返回；并用 `prevFocus` 读回**哨兵 HWND `0x33f1`** 证明方法体真跑了（不是第一行 `return` 蒙过去的绿）。哨兵不是真窗口：`m.hwnd` 仍是 0 ⇒ 方法自己的拒绝检查短路，⛔ 不碰 user32、⛔ 不开窗。

**豁免名册（现量 1 枚）**：`PanelManager.setPriorFocusLocked` → `"requires the caller to hold m.mu; body takes no lock"`。
本格被改名的那两枚**故意不进名册**——它们自取锁，任何名册条目都洗不白。（票面那句"现存那 2 枚进具名豁免名册直到 ⓐ 落地"在 ⓐ 与 ⓒ 同批交付时**没有可落的形状**：改名一落地那两枚就不叫 `*Locked` 了。本腿按派单"现存合法的"那句实现，并加了"名册条目消失即红"那一支补住空档。＝顶回第 3 条。）

**正控（在件内，4 枚，全绿着交）**

- `TestLockedNamingCensusBitesItsOwnFixture33r11`：把同一把尺对准**件内自带的夹具源串**（`parser.ParseFile` 在内存里解析，⛔ 不种进产码）——坏形 `fixtureHost.parkedLocked`（后缀＋自取锁）**必须被报**（且报出 `[f.mu.Lock f.mu.Unlock]` 两枚锁操作）；诚实形 `fixtureHost.honestLocked` 必须干净；无后缀的自取锁 `selfLocksButSaysNothing` **必须够不着**（那是票 33 的 `ⓑ` 家族，本格尺必须不越界）。现量：`33-r11 positive control: planted fixtureHost.parkedLocked reported with [f.mu.Lock f.mu.Unlock], honest name clean, unsuffixed self-locker untouched`。
- `TestLockedNamingCensusRejectsTheOldName33r11`：把**改名前的形状**（`roundTripHost.firstRoundTripLocked`，体内 `m.mu.Lock()`）喂给尺 ⇒ 必须报，且**不在豁免名册**里（否则 `ⓐ` 就不必做了）。现量绿：`the pre-rename shape is reported as roundTripHost.firstRoundTripLocked with [m.mu.Lock m.mu.Unlock]`。
- `TestCompletedWithinReportsAParkedCall33r11`：件内 fixture 类型 `fixtureHost33r11` 的 `parkedLocked()`（后缀＋自取锁）在调用者持锁时被调 ⇒ 有界等待**必须报"没返回"**，并且**读回 `field==0`** 证明它卡在锁上而不是跑完才回。现量绿：`a self-locking Locked call parked, was reported as parked, and had touched nothing`。
- 名册非空支：census 打印每一枚命中（文件／行号／体内锁操作枚数），现量 `1 Locked-suffixed func(s) across 34 production source file(s)` ＝ `PanelManager.setPriorFocusLocked in panel_host_windows.go line 625: 0 lock operation(s)`。

`⛔ 符号没写进会被自己扫的注释`：新件全文 ASCII 注释（无 `U+26D4`、无 `✓`、无 `≤`），ban #8 现量 clean（见 ⑤）。

---

## ③ 14 枚 `file:LINE [token]` 名册逐枚复跑（改后）

尺＝本腿自造，与 `probes/255/v2` 那把同形：`grep -ohE '[A-Za-z0-9_./-]+\.go:[0-9]+ \[[^]]+\]' cmd/wisp/config_readers_255.go | sort -u` 抽枚 ⇒ 逐枚 `sed -n "Lp"` 取该行 ⇒ `grep -qF token`。件＝`logs/22-roster-14-cites-recheck.txt`＋`logs/22b-count-and-line-shift.txt`，各段自落 rc。

`distinct_cites_in_config_readers_255=14`、`14/14 HIT`、`loop_rc=0`：

```
HIT  cmd/wisp/logsink.go:149 [Level: logSinkLevel]
HIT  cmd/wisp/models.go:163 [cfg.Models.Mirror]
HIT  cmd/wisp/models.go:183 [func loadModelConfig]
HIT  cmd/wisp/panel_config_store.go:92 [cfg := s.mgr.Config()]
HIT  cmd/wisp/panel_config_store.go:96 [cfg.LLM.Roles.Chat.Model]
HIT  cmd/wisp/panel_host_windows.go:262 [Width:  uint(width)]
HIT  cmd/wisp/panel_resident_windows.go:207 [cfg.Panel.Width]
HIT  cmd/wisp/resident_ball_windows.go:276 [Hotkeys:  cfg,]
HIT  cmd/wisp/run.go:1014 [cfg.Agent.PerToolTimeoutMS]
HIT  cmd/wisp/run.go:424 [rt.cfg = cfg]
HIT  cmd/wisp/run.go:435 [res := llm.NewResolver(cfg, st)]
HIT  cmd/wisp/run.go:991 [cfg := rt.cfg]
HIT  internal/ball/ball_windows.go:64 [SizePx]
HIT  internal/observe/logging.go:50 [RedactPaths]
```

**零行号位移的证法**（不是"应该没挪"，是量出来的）：`git diff --numstat`＝`14 14`（`panel_host_windows.go`）／`5 5`（`panel_pageover_33r10_windows_test.go`），`head_lines == worktree_lines`（969／293，两枚都 equal=yes），被 cite 那行逐字仍在 `:262`（`\t\tWidth:  uint(width),`）。为守住这一条，本腿把两枚改名方法的"它自己取锁"那句说明**折进原有注释行数**，中途多出的 1 行压回去了。其余 13 枚所在文件本腿没开过（`git status --porcelain -- cmd/wisp` 只列本腿三枚件）。

⚠ 一处与 `probes/255/v2/logs/09-roster-cites.txt` 的差异，具名报：**它当时最后两枚记成 MISS**（`internal/ball/ball_windows.go:64`、`internal/observe/logging.go:50`，它自己写的 `(line 64 = )` 空行），本腿同一枚**逐枚 HIT**。两把尺对 `sed` 取行与空行的处理不同，谁都没动盘；按本仓"读数漂就两把都留"的定式，此处以本腿这把为准还是它为准＝**编排者裁**，本腿不替那格翻账。

另：`tools/d22scan` 的 cite 尺认**路径存在**（ban #9），不认符号名——本腿改的全是符号名，路径一枚未增删，d22scan rc=0 印证了派单那句更正（"改名不会打红那两把 255 尺，增删行才会"）。本腿也**没往 `config_readers_255.go` 里种任何字串**，`panel_geometry_255_test.go:377` 那枚 `retired` 防漂针射程的 `"panel_host_windows.go:304 ["` 与 `"Width:  420,"` 都不出现（255 那两把尺本腿没跑，见 ⑦ 第 5 条）。

---

## ④ 突变表（三发各自咬一支牙）＋红句逐字＋还原

| 发 | 形状 | 落法 | rc | 红句（逐字，截断处为行宽不是省略号） |
|---|---|---|---|---|
| **CONTROL** | 只在 overlay 拷贝里把声明改名成 `firstRoundTripMovedAwayControl`、调用点仍要旧名 ⇒ 真用了拷贝就编不过 | `-overlay=mutations/overlay-control.json` ＋ `go vet` | **rc=1＝这一发的"过"就是 rc=1** | `vet.exe: .\.scratch\wisp\probes\33\r11\mutations\mut-control-symbol-moved-outside-the-overlay-copy.go:448:12: m.firstRoundTrip undefined (type *PanelManager has no field or method firstRoundTrip)` |
| **M1** | 把旧名 `firstRoundTripLocked` 原样放回（声明＋两处调用点），体内自取锁不动 | **盘上**（`sed -i` 两枚跟踪件），跑完从 `mutations/restore-*` 还原 | **rc=1**；census 红、其余 4 枚绿 | `PanelManager.firstRoundTripLocked (panel_host_windows.go line 840) carries the "Locked" suffix but no roster entry: a new name that claims the caller holds a mutex must be reviewed and named in lockedNamingRoster33r11, or renamed` ／ `... is named as if the caller holds the lock while its own body takes it ([m.mu.Lock m.mu.Unlock]) - that is the trap ticket 33 recorded: a caller that trusts the name parks on a non-reentrant sync.Mutex. Roster note: ""` |
| **M2** | 新加**一枚不带自取锁**的 `func (m *PanelManager) noteColdMsLocked(v float64) { m.lastColdMs = v }` ⇒ 只有名册那一支能响 | **盘上**（overlay 拷贝直灌跟踪件），跑完还原 | **rc=1**；只有 roster 一支响 | `PanelManager.noteColdMsLocked (panel_host_windows.go line 973) carries the "Locked" suffix but no roster entry: a new name that claims the caller holds a mutex must be reviewed and named in lockedNamingRoster33r11, or renamed` |
| **M3** | 把名册里那枚诚实的 `setPriorFocusLocked` 改成自取锁（`m.mu.Lock(); defer m.mu.Unlock()`） | `-overlay=mutations/overlay-m3.json`（跟踪件没动） | **rc=1**；行为用例红 2.00s，其余 5 枚绿 | `panel_locked_naming_33r11_windows_test.go:300: setPriorFocusLocked parked while the caller held m.mu: the method takes that lock itself, so its Locked suffix is a trap and any caller that trusts it deadlocks` ＋ `--- FAIL: TestRosteredLockedMethodReturnsWhileCallerHoldsTheLock33r11 (2.00s)` |

**还原与盘上残留**：M1／M2 每次跑完立刻 `cp` 回原文件并逐枚 `git hash-object` 对 HEAD blob——`58e2b155…`／`dae053d8…`／`345e53a6…` 三枚全部 `equal=yes`；`git status --porcelain -- cmd/wisp` 空（`status_rc=0`）。M3 与 CONTROL 全程 `-overlay`，跟踪件从未被打开。突变件（4 枚 `mut-*.go`＋5 枚 `overlay-*.json`＋2 枚 `restore-*.go`）只建不删，都在 `mutations/`。

### ★这一格量到的形状级发现（顶回第 1 条，见文末）

**`-overlay` 只换编译器眼里的源，换不到"读盘的那把尺"眼里的源。** 证据不是推论，是同一发突变里的两半各自表现：M3 用 overlay 把 `setPriorFocusLocked` 改成自取锁 ⇒ **行为那半红了**（编译进去的就是拷贝），**census 那半绿着**（它 `os.ReadFile` 读到的是没改的跟踪件）。所以：
- 凡**读自己源**的守卫（AST／字串／行号一族），突变的正确落法＝**盘上种、跑完还原、逐枚 hash-object 对 HEAD 拉平**，overlay 落法会让那一半**假绿**；
- 新定式里"同路径换未定义符号验 overlay 落地"那把正控**只给编译面背书**，不给源读面背书——本格用 CONTROL 验了编译面、用 M1／M2 验了源读面，两支各自有牙。

---

## ⑤ 门禁逐条 rc ＋ 整包作差

每条自己落一行 rc，件在 `logs/`，测 rc 那句前面⛔ 无管道：

| 尺 | 命令形状 | rc | 件 |
|---|---|---|---|
| `go vet ./cmd/wisp/` | 零输出⇒打印 rc | **rc=0** | `20-vet.txt` |
| 本腿 5 枚新用例 ＋ 33-r10 那枚不开窗 AC#13 用例（同批共 6 枚，`-run` 分母逐枚写出） | `go test -count=1 -v -timeout 300s -run '…6 枚…' ./cmd/wisp` | **rc=0**，6 枚选中／**6 PASS／0 FAIL／0 SKIP** | `21-named-new-tests.txt` |
| 14 枚名册复跑 | 自造尺（见 ③） | **rc=0**（`loop_rc=0`／`comm`／`numstat` 各段自落） | `22-roster-14-cites-recheck.txt`、`22b-count-and-line-shift.txt` |
| 行尾符尺（先于 gofmt） | `tr -cd '\r' < 件 \| wc -c` 对 `git show HEAD:件` | 三枚件 **0 对 0**；`git ls-files --eol`＝`i/lf w/lf attr/text eol=lf` | `23-cr-and-gofmt.txt` |
| `gofmt -l` 本腿三枚件 | 零行＝格式化过 | **rc=0** | 同上 |
| `sh scripts/d22scan.sh`（仓根 `go run ./tools/d22scan` 是错起法，本腿没走） | `d22scan: clean - no D22 ban violations`；ban #8 射程 `cmd/` 111 枚 Go 文件（注释与 `_test.go` 都算） | **rc=0** | `61-d22scan.txt` |
| 整包**跑到终态**两把（带原生 DLL 的 harness） | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 -timeout 500s ./cmd/wisp`；A 把再加 `-overlay=mutations/overlay-baseline.json` | RUN A `baseline_rc=1`（472.033s，`ran_to_terminal=1`、`panic=0`）；RUN B `postchange_rc=1`（472.279s，同） | `30-…`、`31-…`、`50-package-name-diff.txt` |
| 跟踪件未被污染 | `git status --porcelain -- cmd/wisp internal tools` | **rc=0**，列表空 | `60-worktree-clean-and-hash-parity.txt` |

harness 那一课本腿验过是真的：不带两枚原生 DLL 的跑法不在此列——本腿两把整包都从**同一枚**带 PATH 的命令起，且**都出了逐名 FAIL 列表**（不是零 `--- FAIL` 的 `0xc0000135` 空跑）。

**整包逐名作差（"我新增 0 枚"是怎么作差出来的）**

- RUN A＝**改名前的 HEAD**（`03899fd8` 起的两枚跟踪件用 `-overlay` 逐字换回 `git show HEAD:` 的拷贝，并把本腿守卫从构建里摘掉 `"…panel_locked_naming_33r11_windows_test.go": null`）⇒ 红名册 4 枚：`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`、`TestAC14AwaitedBindingReplyReachesThePage`、`TestAC14GoSideEvalPushReachesThePage`、`TestPanelHostRealWindowHopAndLifecycle`。
- RUN B＝**改后盘上态** ⇒ 红名册 6 枚＝A 的 4 枚 ＋ `TestAC4FocusReturnToPriorWindowGap33r5` ＋ `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`。
- `comm -13 A B`＝上面那 2 枚；`comm -23 A B`＝空（本腿没让任何一枚由红转绿——本来也不该，ⓐ＋ⓒ 不改行为）。
- **这 2 枚不是本腿的账，量出来不是靠说**：`logs/70-flip-two-extra-reds.txt` 把两枚在两个状态下各跑 3 次（`-count=3`）——
  - `TestAC4FocusReturnToPriorWindowGap33r5`：**改后态 3/3 红、改名前态 3/3 红**（5.11s 一枚不差）⇒ 它在 RUN A 里绿着就是漂（PASS→FAIL 那一向），真窗焦点用例，与 `TestPanelHostRealWindowHopAndLifecycle` 同族；
  - `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`：**两个状态各 3 次都没出现在 FAIL 列表里**（`flip_postchange_rc=1` 全由 AC4 那枚贡献）⇒ 它是 RUN B 那一发的 FAIL→PASS 漂。
  - 两枚都不在本腿 5 枚之内；本腿 5 枚在 RUN B 里**未进红名册**（非 `-v` 跑法只打印 FAIL，见 ⑦ 第 1 条），并且它们单独跑就是 6/6 绿（`21-…` rc=0）。
- **⇒ 本腿新增红＝0 枚**，按名作差得出，不按"大概都是窗口依赖"得出。派单给的"约 15 枚窗口依赖红"与本腿两把读数（4／6）不同＝顶回第 2 条，写成数字不写成判断。
- ⛔ 没跑 `./internal/...`（禁区：`internal/panel` 那 2 枚具名红灯是另一格的账，本腿一枚没试、`internal/panel/**` 与 `internal/agent/approval/**` 一字未碰）。
- ⛔ 没放宽任何断言、没动 SLO／golden／`thresholds.go`、没新开票、没翻票面勾。

---

## ⑥ 恒真面（哪一发怎么坏它照绿）——具名，不粉饰

1. **射程＝包目录内的生产源**：`*Locked` 出现在 `*_test.go`、或任何**别的包**（`internal/panel`、`internal/config`）里，census **永不响**。票 33 名册的甲形（`internal/config` 那三枚"要持锁却没后缀"）与乙形（`ConfirmLocked` 一族）在这一支下**构造性绿**＝本格只交 `ⓐ`＋`ⓒ` 的自觉，⛔ 别把这一枚读成"同形名册已闭合"。
2. **锁识别只认调用形状** `X.Lock()`/`Unlock`/`RLock`/`RUnlock`/`TryLock`/`TryRLock`。一枚 `*Locked` 方法若**经辅助函数**拿锁（`m.withMu(func(){…})`）、或经**接口值**／嵌入别名绕过去，census 那半**照绿**；`runtime.LockOSThread` 一类被本腿**明确排除**（不是这形制说的 mutex）。名册精确那一支仍能靠名字响，但别把两支护士当一支用。
3. **census 跟着包走，不跟着文件走**：今天 `PanelManager` 在 `panel_host_windows.go`，将来搬文件它照样读到——反过来，⛔ 别把这一枚读成"钉住了 `panel_host_windows.go` 这一枚文件"。（255-v2 那枚 `the host … is the only one` 栽的正是这种名字超出射程的坑，本腿按整包口径写死。）
4. **行为那半只覆盖名册里 1 枚中的 1 枚**：唯一不开窗够得着的 `setPriorFocusLocked`。若名册将来加进一枚要真窗／要控件的 `*Locked`，它会**名册绿、行为没测**——不对称是设计死的，⛔ 不是"持锁形状已被行为证明"。
5. **2 秒预算的形状**：一枚"先卡 2 秒以上、然后正常返回"的 `*Locked` 会被报成 parked ⇒ 红是红了，**理由是错的**。这一支证的是"不自取锁时会返回"，不是"返回时间多长"。
6. **哨兵句柄走的是短路支**：`m.hwnd==0` ⇒ 同根检查不进 user32。所以这枚行为用例**不证明** `sameRootWindow` 在任何状态下正确（那是 33-r4／33-v1 §A#27 那一族的账），只证明调用没卡在锁上、且体真跑了。
7. **`census` 非空支钉的是"读到文件"，不是"读到 34 枚"**：新增文件不会让它响；响的是**零枚**。别把它读成文件数护条。
8. **本腿没让任何绿灯变宽**：`retired` 那两枚字串、`panel_resident_windows_test.go` 的 skip 支、winlive 那族——本腿一枚没动、一枚没跑（见 ⑦）。

---

## ⑦ 我没量到的（空着就写空着）

1. **整包 PASS／SKIP 分母**：两把整包都是非 `-v` 跑法，输出里只有 `--- FAIL` 行，**没有 PASS／SKIP 行**（现量 `PASS=0 SKIP=0`）。⇒ 本腿报得了红名册与"跑到终态"，报不了"346 枚里绿多少枚"。要那个数得再花一把 `-v` 整包（约 8 分钟），本腿按预算没取。
2. **差集之外的漂**：只把差集里那 2 枚做了 3×2 翻转复尺。RUN A 里那 4 枚、以及派单说的其余窗口依赖家族，本腿**没逐枚复跑**，⛔ 不许被读成"整包 15 枚都逐枚验过漂"。
3. **M1／M2 的盘上种法只各跑一次**：还原与 hash 拉平是量过的，但"同一发在另一台机器上会怎么漂"本腿没测。
4. **33-r10 那枚 AC#13 用例的语义**：本腿把它的调用点跟着改名走了，并跑绿（rc=0），**没重做**它的 AC#13 判据（最终文档内容），那一格仍是 33-r10 的账。
5. **票 255 那两把尺本腿没跑**（`panel_geometry_255_test.go`／`config_readers_255.go` 一族）：本腿既没碰它们，也没替"改名会不会打红它们"补跑——那条已由 ③ 的 14/14 HIT ＋ `numstat` 等行 ＋ d22scan rc=0 三面围住。要一把 255 自己的尺的 rc，归票 255 那一格。
6. **`.scratch/wisp/probes/33/r11/logs/10-baseline-full-package.txt`＝0 字节那一格**：本腿第一次把基线当后台任务起，命令自己没落一个字（`baseline_rc` 行都没写）。已在那枚文件里就地标注为"不是读数、被 `30-…` 取代"，⛔ 不当证据用，也不删（规则 8）。
7. **旧名在跟踪文档里的残留**：`grep -rln "firstRoundTripLocked|serveNotBuiltNoticeLocked"` 命中约 20 枚 `.md`（含 `.scratch/wisp/issues/33-panel-host-c27.md` 票面本身、`probes/33/{a2,a3,r10,v4}`、`probes/253/r3`、`probes/274/a2`、`probes/35/*`、`probes/111/c2/logs/*`、`probes/bundle/1`、`probes/e2e-panel-1`、`probes/msg-a653.md`）。`docs/**`、`.scratch/wisp/issues/**`、别人格的 `probes/**` 都是本腿禁区，⛔ 一字未改——**这些引用从今天起是历史名**，翻勾与更正归编排者（票 33 文末那一节）。本腿只跟着改了 `cmd/wisp` 内的 9 处注释引用。
8. **`internal/panel` 那 2 枚红灯**：本腿没跑、没看、没试让它们绿（另一格的账）。

---

## 顶回清单（原文／现量 vs 派单转述）

1. **新定式的适用面（形状级，见 ④ 末）**：`-overlay` 给不了"读自己源的尺"的突变。派单让我抄 255-v2 那发正控（同路径换未定义符号）——本腿照抄了并让它 rc=1，但那把正控只背书**编译面**；源读面必须盘上种＋还原＋hash 拉平（M1／M2 就是这么做的）。建议把这条写进定式：**守卫读盘 ⇒ 突变走盘**。
2. **"cmd/wisp 整包今天有约 15 枚窗口依赖红"** vs 本腿两把终态读数 **4 枚（改名前）／6 枚（改名后）**，且差集两枚 3×2 翻转复尺都在两个状态下漂。⛔ 本腿不据此说那 15 是假的（111-c2 已记过 FAIL 12→6→8→7→6 都漂这一族），只把今天的两个数写在这儿，别让后续程当现状引。
3. **豁免名册那句的可达性**：票面写"现存那 2 枚进具名豁免名册直到 ⓐ 落地"，ⓐ 与 ⓒ 同批时那 2 枚落地即不再叫 `*Locked`，"2 枚进名册"这一形状**没有可交付的瞬间**。本腿采派单"现存合法的进名册"那一读，并补"名册条目消失即红"这一支防烂。
4. **派单 §3 的自我更正本腿复认**：改名**没**打红那两把 255 尺（14/14 HIT、零行位移、d22scan rc=0）——那句"改名会让它们红"确实错了；本腿另外把行数守成不变（14/14、5/5），所以连"增删行才会打红"这一支都没触发。
5. **派单 §2 的撤回项本腿遵守**：`cmd/wisp/config_readers_255.go` 那两处注释字串一字未动（带 ref 限定的历史指法还没裁＝未定义即停，本腿没自创）。
