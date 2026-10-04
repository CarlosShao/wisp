# 票 260 · 落地腿 260-r1 —— 借用那遍吃 `[hotkey] cancel` 的配置值（形ⓐ＋三条硬条件）

**开工时刻**：2026-10-04 09:4x +0800（现取自第一发 shell 的 `git log`／`go test` 时间戳）　**起点 HEAD（现取）**：`528bf7bd`（`dev`）
**任务书**：票面 `.scratch/wisp/issues/260-configured-cancel-hotkey-never-registers.md` 最后一节「编排者选形裁定（2026-10-04 09:2x，账 A588）」＝形ⓐ＋条件①②③；撤销口令「260 改 ⓑ」。
**本格范围**：只做 **AC#1**（借用吃配置）。**AC#2（丢借用要有声）不在本格**，编排者另排一格；AC#3 是越界检查，见 §6。
**写面**：`internal/ball/**` ＋本证据件目录 `.scratch/wisp/probes/260/r1/`。逐项核对见 §6。
**本腿三枚 commit**：`eb2c0173`（骨架＋三把尺原始读数）→ `d9bb6817`（产码＋七枚判据＋门禁）→ `b1e59d63`（idle 那句文案补实话＋三发突变对最终码重跑）。

> ⚠ 本件里的行号是**改后现量**（`grep -n` 现跑，见 §1 末的读数档）；起手尺那一节记的是**改前**位置，两处不同属预期，不要互换引用。

---

## §0 起手三把现量尺（全部现跑，未照任务书行号当常量）

### 0.1 `grep -n "escBorrowAcc\|vkEscape\|RegisterHotKey" internal/ball/*.go`（改前位置，读数档 `ruler1-grep-borrow.txt`，33 行）

| 类别 | 读数 |
|---|---|
| 借用遍的取值落点 | `internal/ball/hotkey_windows.go:521-524` `escBorrowAcc()`＝零参常量字面 `Accelerator{Mods: modNoRepeat, VK: vkEscape}`；常量 `:118 vkEscape = 0x1B` |
| 唯一产码调用体 | `:391` `pRegisterHotKey.Call(hwnd, id, acc.Mods, acc.VK)`（`win32_windows.go:42` 是 proc 声明，不是调用） |
| 稳态那遍的次序 | `:453 registerAllWith` → `:457` 命中 `hkCancel` 即 `cancelIdleLine(bindings[i])` 后 `continue`（`:458` 注释 "The one slot an idle ball must not register"）；`cancelIdleLine` 在 `:428`（改后 `:451`） |
| 借用那遍的次序 | `:531 takeEscWith(unreg, reg)` → `:532 unreg(hkCancel)` → `:533 reg(hkCancel, escBorrowAcc())`；产码调用方只有 `ball_windows.go:878 TakeEscForCancel` → `:543 takeEsc(hwnd)` |
| 回执那行 | `:578-586 cancelBorrowedLine()`：`Binding: escBorrowBinding`（常量 `:519 = "Esc"`）＋`Acc: escBorrowAcc()`＋`Status: HotkeyLive`＝**已带 Binding 字段**，条件②复认通过 |
| 测试面命中 | `hotkey_status_test.go:210/219/259/270/274/340`（默认套件，走零参 `takeEscWith`／`cancelBorrowedLine`／`cancelFailedLine`）；`hotkey_live_test.go:61-62`（winlive 的裸 Esc 探针直接写 `modNoRepeat/vkEscape`）、`:107`（断借到的是 `vkEscape`）、`:472` |

### 0.2 `grep -rn "cancel" internal/config/schema.go | head`（读数档 `ruler2-grep-schema.txt`，4 行；不对称另现读 `:179-184`）

| 类别 | 读数 |
|---|---|
| `[hotkey]` 四枚字段 | `internal/config/schema.go:180-183`：`Summon`／`Mute`／`Cancel`／`Panel` |
| 不对称 | `:182 Cancel string \`toml:"cancel" default:"Esc"\`` —— **只有 cancel 带 `default:`**，summon／mute／panel 三枚没有 ⇒ 库里说不出"我补了默认" |
| 那句自述 | `:177-178` "Empty string = binding unset (feature key disabled). `cancel` temporarily takes over during Confirming and is returned at session end (D36)." |
| 静默补默认（所以"退回了"这句话只能由拿配置那一侧说） | `internal/ball/ball_windows.go:153-155`：`if opts.Hotkeys == (HotkeyConfig{}) { opts.Hotkeys = DefaultHotkeys() }`——整枚零值结构体才补，补完不打一行日志；另一条补法是 host 侧 `ApplyHotkeyDefaults`（`hotkey_windows.go:87`，改后 `:87` 未动） |

### 0.3 `grep -rln "hotkeys live\|Binding" internal/ball/*_test.go cmd/wisp/*_test.go`（读数档 `ruler3-grep-testnames.txt`，10 个文件）＋未修码基线包名册

| 类别 | 读数 |
|---|---|
| 今天已钉住这形状的文件 | `internal/ball/`: `hotkey_live_test.go` · `hotkey_status_test.go` · `interaction_live_test.go` · `live_windows_test.go`；`cmd/wisp/`: `panel_resident_windows_test.go` · `resident_approval_live_246_windows_test.go` · `resident_ball_live_228_windows_test.go` · `resident_hotkey_258_windows_test.go` · `resident_hotkey_live_258_windows_test.go` · `resident_hotkey_v1probe_test.go` |
| 写之前抄进本件的在册用例名 | `TestRegisterAllLiveSet`（`:114`，稳态只绑三枚＋`bindsEsc()` 必须 false）· `TestDefaultHotkeysIdlePassHoldsNoEsc`（`:175`）· **`TestCancelBorrowRoundTrip`（`:202`，`:223` 断借到的是 `vkEscape/modNoRepeat`）** · `TestCancelBorrowFailureIsAProblemLine`（`:268`）· `TestStandbyIsNotDisabledAndNotAProblem`（`:291`）· `TestRegisterAllSplitsFailureFamilies`（`:320`）· `TestUnregisterAllDropsKnownSlots`（`:395`）· `TestHotkeyStatusString`（`:417`，逐字钉 `HotkeyStandby.String()`）· `TestHotkeyReloaderRebindsOnConfigChange`（`:454`）· `TestApplyHotkeyDefaults`（`:524`）· winlive 侧 `requireEscBorrowed`（`hotkey_live_test.go:90`，断言在 `:107`） |
| **未修码基线（本腿改任何码之前）** | `go test -count=1 -v ./internal/ball/` ⇒ 顶层 **56 枚：PASS 55 ／ FAIL 1 ／ SKIP 0**（名单存 `baseline-gotest-names.txt`）。唯一那枚红＝`TestC21TableColourRowsMatchTokensCSS`，红句逐字：`read design/assets/tokens.css: open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css: The system cannot find the path specified. - the CSS leg of this check must never skip`；红因＝共享工作树里别家把 `design/assets/*.css` 删了（`git status` 现见 ` D design/assets/tokens.css`），**先于本腿存在、与本票无关，本腿一字未动它、也没放宽它** |

---

## §1 产码落点（形ⓐ四道，file:line 为**改后现量**；读数档 `gate-ac3-my-commits.txt` 等）

| # | 落点 | 干了什么 |
|---|---|---|
| 1 | `internal/ball/hotkey_windows.go:563`（`type cancelBorrow`）＋`:588`（`func resolveCancelBorrow(bind string) cancelBorrow`） | 借用遍唯一的取值口，三种入参三种定案：能解析→按配置借；`""`→产品默认裸 Esc（`Fallback` 留空，与 `ApplyHotkeyDefaults:87` 那句"safe direction"同向，⛔ 不许把"没配"变成"不借"，那是 245 的时机裁定）；解析不了→退回默认 **且带 `Fallback` 错误**（新零件 P6）。默认分支逐位＝旧字面量 `{0x4000, 0x1B}`（`:554 escBorrowAcc()` 未动） |
| 2 | `internal/ball/hotkey_windows.go:609 takeEscWithAcc` · `:623 takeEscWith`（默认档包装）· `:631 takeEscBorrow`；调用方 `internal/ball/ball_windows.go:886 borrow := resolveCancelBorrow(b.cancelBinding)` → `:892 takeEscBorrow(b.hwnd, borrow)` → 回执 `:895/:900 withCancel(cancelFailedLineFor/cancelBorrowedLineFor(borrow))` | 产码把**同一枚** `cancelBorrow` 同时交给 RegisterHotKey 和回执行（不两处各拼）。`b.cancelBinding` 就是既有链上的那枚值：boot `ball_windows.go:253`（`b.opts.Hotkeys.Cancel`，opts 由 host 的 `ApplyHotkeyDefaults`／`ball.New:154` 补默认）＋rebind `:821`（`cfg.Cancel`）⇒ **零新取值路径**，票 258 那条链一字未动即自动获益 |
| 3 | `internal/ball/hotkey_windows.go:314`（`HotkeyBinding` 新增 `Note string`）＋`:384`（`Problems()` 里 `if b.Note != ""`）＋`:671 cancelBorrowedLineFor` · `:687 cancelBorrowedLine` · `:695 cancelFailedLineFor` · `:703 cancelFailedLine` | P6 那句怎么到用户眼前：退回默认的 Live 行带 `Note`，`Problems()` 对**任何状态**都把 Note 念出来（旧代码 `HotkeyLive` 一律 `continue`，一句"配坏了"会被吞）。`cmd/wisp/resident_ball_windows.go:340` 已在逐行打 `Problems()`（本腿没动它），所以这句不需要新出口 |
| 4 | `internal/ball/hotkey_windows.go:471`（idle standby 那行日志）＋`ball_windows.go:863-879`（`TakeEscForCancel` 文档，函数体 `:881`）＋`:905-910`（`ReleaseEscAfterSession` 文档，函数体 `:911`） | 说到做到的收尾：旧日志逐字 `"cancel hotkey left unbound while idle: Esc is borrowed only during Confirming..."` 而同一条的 `binding` 字段带的是用户自己的值（260-a2 census 点名的自相矛盾句）；借用吃配置之后它由"矛盾"变成"说谎"，故改成点名槽位不点名键：`"the configured key is borrowed only during Confirming"`。⛔ **没动** `HotkeyStandby.String()`（`:285`）——那句本就只说槽位，且被 `hotkey_status_test.go:424` 与 `cmd/wisp/resident_ball_live_228_windows_test.go:166` 钉字面 |

零参 `takeEscWith` / `cancelBorrowedLine` / `cancelFailedLine` 一律保留为**默认档包装**（`:623/:687/:703`，注释写明用途），票 245/64 钉住的那四枚断言（`hotkey_status_test.go:210/219/259/270/274/340`）**一字未改、也未搬迁**；`hotkey_live_test.go` 本腿**一个字没动**（冲突具名在 §5）。

---

## §2 常驻判据（`internal/ball/hotkey_cancel_borrow_260_test.go`，7 枚，全在默认套件／不需桌面）

| # | 条件 | 用例名（`:行`） | 断言形状（是否词面尺） |
|---|---|---|---|
| ① | 默认档不许漂 | `TestBorrowDefaultIsBitIdentical260`（`:57`） | 期望值是裸字面量 `Accelerator{Mods:0x4000, VK:0x1B}`（`:35 accelDefault260`），⛔ 不是 `escBorrowAcc()` 也不是字符串；另核 `modNoRepeat==0x4000 && vkEscape==0x1B`（防字面量过期）、`ParseAccelerator("Esc")` 产出同一对、三种"没配"入口（`""`／`DefaultHotkeys().Cancel`／`ApplyHotkeyDefaults(HotkeyConfig{}).Cancel`）、假注册表真收到的那枚、回执行 `Binding/Acc/Note` 三项，并带"尺必须能看见 Esc"的正控（`reg.bindsEsc()` 必须 true） |
| ② | 吃配置的正控 | `TestBorrowFollowsConfiguredBinding260`（`:114`） | 种 `"Ctrl+Alt+K"`／`"Ctrl+Alt+Space"`／`"F9"` 三枚 ⇒ 解析值、注册表收到那枚、回执 `Binding` 三项全随它变；负向半 `borrow.Acc == escBorrowAcc()` ⇒ "still resolves to the hard-coded bare Esc" 红，以及 `reg.bindsEsc()` 必须 false |
| ③ | P6 会响 | `TestBorrowUnparsableFallsBackLoudly260`（`:164`） | 种 `"Ctrl+Alt+NotAKey+"` ⇒ `Fallback != nil`（否则 `P6 RED: ... swallowed`）、退回逐位＝默认、`attemptsOf(hkCancel)==1`（借用仍存在，245 时机未动）、回执 `Note` 非空、`Problems()` **恰 1 行**且逐词含 `cancel`／坏值本身／`not a valid key combination`／`default Esc`／`borrowed for this card`；正控＝默认档同尺必须 0 行；再核 idle 那半仍照旧念坏值（本改没把 idle 的 Unparsable 静音） |
| ④ | 回执说实话 | `TestBorrowReceiptSharesOneSourceWithRegistration260`（`:226`） | 对六枚入参逐条核两件事：`line.Acc == 真交给 RegisterHotKey 的那枚`（不同即 `the receipt carries ... while RegisterHotKey was handed ... - the two are not the same resolution`）、`ParseAccelerator(line.Binding) == line.Acc`（打印的键名解析回来必须就是打印的加速器）⇒ 这两条就是"两处各拼一次"的克星（突变 M2 读数见 §4） |
| ⑤ | 245 边界闸（顺带） | `TestConfiguredCancelStillNeverBoundWhileIdle260`（`:292`） | `Cancel:"Ctrl+Alt+K"` 下稳态 `attemptsOf(hkCancel)==0`、`Live()` 仍 3 枚、cancel 行仍 standby 且带配置拼法；归还后仍 0 次再注册 |
| ⑥ | 条件②的失败半 | `TestBorrowFailureNamesTheKeyItTried260`（`:260`） | Win32 拒了配置那枚 ⇒ 回执必须念它真试过的键（`Binding=="Ctrl+Alt+K"`、`Status.Attempted()`、`Problems()` 同时含键名与 "was not registered"）；正控＝默认档失败行仍是 `Esc`/`{0x4000,0x1B}`（245 的旧尺读的那枚） |
| ⑦ | 具名入账别人的尺 | `TestLiveBorrowHelperPremiseIsConfigDependent260`（`:332`） | 把 `liveHotkeys()`（winlive `hotkey_live_test.go:41`）的 `Cancel:"Ctrl+Alt+V"` 抄进默认层跑一遍：借到的必须 VK `'V'`／mods `0x4003`，⛔ 不许等于 `vkEscape`；`t.Logf` 逐字点名 `hotkey_live_test.go:107` 那条前提因此变成"看配置"的性质。该文件本腿未改、winlive 层本腿未跑 |

另写 `internal/ball/hotkey_cancel_borrow_live_260_test.go`（`//go:build windows && winlive`，2 枚：`TestLiveCancelBorrowDefaultUnchanged260` 真桌面下默认档仍借到裸 Esc 且桌面拿不回、`TestLiveCancelBorrowFollowsConfig260` 配 `Ctrl+Alt+Y` 后用 `uiRegisterHotkey(spareHKID)` 逼 Win32 自己回 1409 来证明"真注册上了那枚"，并核此时裸 Esc 探针应报 free）。**本轮只过 `go vet -tags winlive`（类型检查），⛔ 未跑**——欠账见 §5。

---

## §3 门禁终值读数（逐字抄；读数档 `gate-*.txt`，全部对最终码重跑）

| 门禁 | 读数 |
|---|---|
| `go build ./...` | 无输出，exit `0`（`gate-build-full.txt`，0 字节） |
| `go vet ./internal/ball/` | 无输出，exit `0`（`gate-vet-ball.txt`） |
| `go vet -tags winlive ./internal/ball/`（只为让 winlive 那两枚过类型检查，不跑） | 无输出，exit `0`（`gate-vet-ball-winelive-tag.txt`） |
| `go test -count=1 ./internal/ball/`（只这一包，未跑全仓 `./...`） | `PASS=62 FAIL=1 SKIP=0`，包行 `FAIL github.com/CarlosShao/wisp/internal/ball 0.269s`；那枚 FAIL 逐字仍是 `--- FAIL: TestC21TableColourRowsMatchTokensCSS`（§0.3 的未修码基线红，别家删掉的 `design/assets/tokens.css`）——**62＝基线 55＋本腿 7 枚，一条没掉**（`gate-test-ball-verbose.txt`） |
| `sh scripts/d22scan.sh` | 终值 exit `0`：`runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0` ＋ `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=228, bans #1-5 cmd/=38, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=500, ban #8 cmd/=97`（`gate-d22scan-final.txt`）。起手基线同脚本＝`internal/=498`（`gate-d22scan-baseline.txt`），差 2＝本腿新增两枚测试件。**未出现 212-v1/v2 突变窗口那种红在 `.scratch/wisp/issues/**` 缩写路径上的读数**；本腿另在提交前复跑一遍取终值（三次跑全 clean） |
| `gofumpt -l internal/ball/`（`$(go env GOPATH)/bin/gofumpt`） | 空输出（`gate-gofumpt.txt`，0 字节）。中途 `hotkey_windows.go` 被报过一次注释项目符号缩进，已 `-w` 归位 |
| 附带（不在指派门禁里，如实记）：`staticcheck ./internal/ball/` | 本机报 `-: internal error in importing "internal/byteorder" ... export data version 4 is greater than maximum supported version 2`——本机 staticcheck 与当前 Go 版本不配，**与本腿无关、也未列为门禁**，故不作为判语；顺带核过 `takeEscBorrow` 有产码调用方（`ball_windows.go:892`），`takeEscWith` 只剩 245 的测试在用（＝它 documented 的用途），无悬空函数 |

---

## §4 突变读数（摘掉新逻辑 ⇒ 指名用例必须红；还原 ⇒ 绿且 md5 一致、工作树干净）

驱动＝`mutation-driver.sh` ＋ `mutate.py`（三枚精确字符串替换，各摘一道新逻辑）。**跑了两轮**：
第一轮对着 `d9bb6817` 那版码（`hotkey_windows.go` md5 `75136db8c056c3e8e127e352a6d733fd`）逐发单跑，读数 `mutation-M1-run.txt`／`mutation-M2-run.txt`／`mutation-M3-run.txt`，还原后 md5 复见 `mutation-M1-restored-md5.txt` 与 `mutation-restored-md5-final.txt`（都是 `75136db8…` 逐字相等）；⚠ 那一轮里 M2 的还原我用了一次**单文件** `git checkout HEAD -- internal/ball/hotkey_windows.go`（不是 `checkout .`，只指我自己那一枚已提交文件，还原后 md5 核过相等），此后本腿改走 `cp hotkey_windows.go.260r1*backup`，如实记这一笔。
第二轮（终值）对着**最终码**（md5 `be79cbda486a334d74fe9d278558acba`，即 `b1e59d63` 里那版）用驱动连跑三发，全部读数与还原核验在 `mutation-M1M2M3-final.txt`，突变前的基线 md5 在 `mutation-final-md5-baseline.txt`。两轮的**红句逐字相同**；差别只有一处：第一轮 M1 那发的 `-run` 名单里没带 ⑥（`TestBorrowFailureNamesTheKeyItTried260`），第二轮七枚全跑，于是 M1 也把 ⑥ 顶红了（表里 M1 行列的就是终值那一发）。

| 突变 | 摘掉了什么 | 红句（逐字，前缀行号＝`hotkey_cancel_borrow_260_test.go`） | 还原 |
|---|---|---|---|
| **M1** | `:601` 成功分支不再用配置来的 `acc`，退回旧的写死裸 Esc（＝本票的缺陷本体） | `127: modifiers plus a letter: cancel = "Ctrl+Alt+K" resolves to {Mods:16384 VK:27}, want {Mods:16387 VK:75}` ／ `127: the documented alternative summon shape: cancel = "Ctrl+Alt+Space" resolves to {Mods:16384 VK:27}, want {Mods:16387 VK:32}` ／ `127: a bare function key: cancel = "F9" resolves to {Mods:16384 VK:27}, want {Mods:16384 VK:120}` ／ `252: cancel = "Ctrl+Alt+K": receipt spells "Esc", want "Ctrl+Alt+K"` ／ `274: the failed-borrow line = {Name:cancel ID:3 Binding:Esc Status:registration failed ...}, want an attempted error naming "Ctrl+Alt+K"` ／ `338: the borrow of "Ctrl+Alt+V" resolved to VK_ESCAPE again - condition ①'s default-only rule was over-applied` | 红 4 枚（②④⑥⑦），绿 3 枚（①③⑤——它们测的是默认档与 P6，本就不该被这发摘动）；`--- FAIL: TestBorrowFollowsConfiguredBinding260` 等 |
| **M2** | `:677` 回执行的 `Acc` 改回常量 `escBorrowAcc()`，而注册仍用配置值（＝"两处各拼一次"的形状） | `236: cancel = "Ctrl+Alt+K": the receipt carries {Mods:16384 VK:27} while RegisterHotKey was handed {Mods:16387 VK:75} - the two are not the same resolution` ／ `151: modifiers plus a letter: receipt line = {Name:cancel ID:3 Binding:Ctrl+Alt+K Status:live Acc:{Mods:16384 VK:27} Err:<nil> Note:}, want live "Ctrl+Alt+K" with {Mods:16387 VK:75} and no note` ／ `274: the failed-borrow line = {... Binding:Ctrl+Alt+K ... Acc:{Mods:16384 VK:27} ...}, want an attempted error naming "Ctrl+Alt+K"` ／ `345: RegisterHotKey would get the configured key while the receipt still prints Esc` | 红 4 枚（②④⑥⑦），绿 3 枚（①③⑤）；④ 的两条断言各命中一次，正是它存在的理由 |
| **M3** | `:597-598` 的 `Fallback` 整段摘掉（解析失败静默吞＝P6 不落地） | `172: P6 RED: an unparsable cancel binding resolved with no fallback notice - the parse failure was swallowed` ／ `252: cancel = "Ctrl+Alt+NotAKey+": receipt spells "Esc", want "Ctrl+Alt+NotAKey+"`（④ 的拼法表以 `Fallback` 为判据，故同红；这条是副产物，如实记） | 红 2 枚（③④），绿 5 枚（①②⑤⑥⑦） |

**还原后的机器读数**（`mutation-M1M2M3-final.txt` 末段）：三发之后 `md5sum internal/ball/hotkey_windows.go` 三次全等于 `be79cbda486a334d74fe9d278558acba`（与 `mutation-final-md5-baseline.txt` 里突变前的那枚逐字相同），`internal/ball/ball_windows.go` = `2488b2ac692c99ae4e3b4831526637d7` 未变；`git status --short internal/ball/` 行数 `0`（本腿两枚产码文件在 `b1e59d63` 已入库，故还原面相对 HEAD 也干净）；七枚判据 `--- PASS × 7` ＋ `ok github.com/CarlosShao/wisp/internal/ball 0.058s`。

---

## §5 没做成的／欠账（具名＋归口）

| 项 | 归口 | 说明（含"我到底跑没跑"） |
|---|---|---|
| **AC#2 零症状那一格要有声** | 票 260 的独立一格（编排者另排） | 本腿**没做**，任务书写明不许顺手做。丢借用（Win32 拒了／根本没成键）今天仍只落在 `HotkeyError` 回执行与 `slog.Error`（`ball_windows.go:887-894`），没进 `Problems()` 之外的告警面；票面第 12 行那条"Standby 不进 `Problems()` ⇒ 零症状"本格未触 |
| **winlive 两枚新尺没跑** | 〔未量有没有牙〕归 260-v1／真窗窗口 | `internal/ball/hotkey_cancel_borrow_live_260_test.go` 的 `TestLiveCancelBorrowDefaultUnchanged260`、`TestLiveCancelBorrowFollowsConfig260`：本轮**只跑了 `go vet -tags winlive`（类型检查通过）**，⛔ 未执行——它们要真窗＋真 `RegisterHotKey`（一枚还会在 2-3 秒里全局占住 `Ctrl+Alt+Y`），会抢 owner 桌面。**没跑过的没写成通过** |
| **★别人的尺会被顶到（本包内）：`internal/ball/hotkey_live_test.go:107`** | 交回编排者裁（本腿⛔ 未改那枚文件） | 现量链条：`liveHotkeys()`（`hotkey_live_test.go:41`，函数体 `:42` 逐字 `return HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Ctrl+Alt+V", Panel: "Ctrl+Alt+B"}`）给球喂 `Cancel:"Ctrl+Alt+V"`，而 `requireEscBorrowed`（`:90`）在 `:107-108` 断 `rep.Live()[hkCancel].VK == vkEscape`（红句 `the borrowed cancel slot holds VK 0x%X, want VK_ESCAPE 0x1B`），过不了还要在 `:110-114` 拿裸 Esc 探针逼出 `ticket 245 RED: the ball says it borrowed Esc, but the desktop still has it free`——两条都按"借的一定是那枚裸 Esc"写。形ⓐ 之后那枚槽真持有 `'V'`，该前提不再成立 ⇒ 四条 winlive 调用（`hotkey_live_test.go:311`、`interaction_live_test.go:205/257`、`live_windows_test.go:127`）在下一个 winlive 窗口会红，预计红句逐字：`the borrowed cancel slot holds VK 0x56, want VK_ESCAPE 0x1B`。**这条不是猜的**：本腿在默认层把同一条取值路径跑了一遍并常驻成判据 ⑦（`TestLiveBorrowHelperPremiseIsConfigDependent260`，`hotkey_cancel_borrow_260_test.go:332`；它在默认层跑同一条取值路径并 `t.Logf` 逐字点名那两处行号，读数见 `gate-test-ball-verbose.txt` 末段 `hotkey_cancel_borrow_260_test.go:347:` 那一行），但**winlive 本体我没跑，所以"会红"仍是预测，不是实测**——两件事分开记 |
| **`cmd/wisp` 会不会被顶到** | 归编排者裁（本腿未动 `cmd/wisp` 一字） | 我先试着实测：`go test -count=1 -v ./cmd/wisp/` ⇒ 输出只有 `exit status 0xc0000135` ＋ `FAIL github.com/CarlosShao/wisp/cmd/wisp 0.034s`（读数档 `probe-cmdwisp-full.txt`），即**加载期 DLL 缺失，本机根本进不了任何用例**——这条是 `.scratch/wisp/issues/98-cmd-wisp-tests-never-run-on-this-host.md`（open，2026-09-21 起）记在案的既有仪器缺口，先于本票存在，本腿不去修它。**所以这格只能给静态读数**：① `cmd/wisp` 全部含 cancel 的种值今天都是 `"Esc"`（`grep 'cancel *= *"' cmd/wisp/*_test.go` 只命中 `resident_hotkey_258_windows_test.go:115/236/274/301/336/351/381/410` 与 `:126` 的 want，全为 `"Esc"`；246 系 rig 根本不写 `[hotkey]` 段⇒ 走 schema `default:"Esc"`）⇒ 落进条件①的默认档，逐位不变；② 那族尺子的断言形状是 `requireIdleCancelSlot246`（`resident_approval_live_246_windows_test.go:350-361`：数 `rep.Live()`==3、`EscTakenOver()` false）与 `observeEsc246`（`:441`，量真实桌面 Esc 能不能回到别的窗口），两者都读**稳态集合／布尔位/桌面效应**，不读"借来的是哪枚加速器"⇒ 默认档下不动；③ `resident_hotkey_258_windows_test.go:453` 数 rebind 后 live==2，同样与加速器来源无关。⇒ 静态判断＝**打不红**；但"跑过且没问题"与"没跑过"在这台机器上长得一样，故本格正式挂**〔未实测〕**，归 260-v1 与票 98 一起收 |
| **条件②的另一半在别人的写面：卡片文案仍写死"按 Esc 键"** | 需编排者派一枚文案／回执侧的腿（`internal/agent` ＋ `cmd/wisp`，本腿⛔ 不许动） | 现量：`internal/agent/approval/approval.go:90 ChannelEsc: "按 Esc 键"`、`:112 "Esc 取消不可用"`；`cmd/wisp/resident_approval_windows.go:180/:182` 打印"按 Esc 否决卡片…"、`:457` "按 Esc 不会否决它"。功能上没问题——`WM_HOTKEY` 按 `hkCancel` **这个 id** 派发，与本腿换了来源的加速器无关（该链条未改动，本腿也未跑真卡片）；但用户把 cancel 改成 `Ctrl+Alt+K` 之后，卡片那句话念的仍是 Esc ⇒ 形ⓐ 之后这一处**新出现**的"页面写着按 Esc、实际借的是别的键"。它不在 `internal/ball/**`，本腿按边界只上报不修。⚠ 上报口径：本腿没说"这条路接不上"。配置值到借用点那半截是**实测**（§1 落点 1/2 ＋ §4 三发突变各有红句），"卡片那句念错键"这半截只到**文本现量**（上面三处 `grep` 逐字），它的效果面要真卡片＋真窗才量得到，本腿既不许动那两个包也没跑那一层 ⇒ 记为**位置问题**交给编排者派腿，不当否证入账 |
| `TestC21TableColourRowsMatchTokensCSS` 红 | 归把 `design/assets/*.css` 删掉的工作树动作（别家）／或票 245 家族自己 | §0.3 记的未修码基线红，本腿未动它，也没为变绿放宽任何断言；门禁读数里它一直红着，终值 `PASS=62 FAIL=1 SKIP=0` |
| `internal/ball` 之外没跑全仓 | 编排者的排程要求 | 任务书写明 ⛔ 不跑 `./...`（此刻 262-r1 在 `scripts/`＋`.github/workflows/ci.yml`、212-v2 在临时突变 `tools/d22scan/main.go`）；本腿只跑了 `go build ./...`（编译，不执行）与 `./internal/ball/` 单包测试 |

---

## §6 越界检查与票面禁区复核（AC#3）

| 检查 | 读数 |
|---|---|
| 本腿三枚 commit 碰过的路径全集 | `git show --pretty=format: --name-only` 逐枚取（读数档 `gate-ac3-my-commits.txt`）：全部落在 `internal/ball/hotkey_windows.go` · `internal/ball/ball_windows.go` · `internal/ball/hotkey_cancel_borrow_260_test.go` · `internal/ball/hotkey_cancel_borrow_live_260_test.go` · `.scratch/wisp/probes/260/r1/**`。过滤后越界面行数 `0`（`grep -vE '^###|^\.scratch/wisp/probes/260/r1/|^internal/ball/'` ⇒ 0） |
| `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` | 全段 `git diff --name-only 528bf7bd..HEAD` 里命中数 `0`（`gate-ac3-pathspec-scope.txt`，含别家 commit 的整段区间也没命中）⇒ 未触 |
| D43 转移表／新增状态 | `git diff --name-only` 里 `statemachine` 命中数 `0`；`HotkeyStatus` 枚举**未加任何一枚**（`hotkey_windows.go:239` `type HotkeyStatus`…`:265 HotkeyStandby`…`:266` 收尾括号，一枚未增未减，`TestHotkeyStatusString` 的字面表原样通过）；`Problems()` 只多了"有 Note 就念出来"这一条分支（`:384`），未新增状态、未改任何状态的归属 |
| 票面 AC 复选框／台账 | `.scratch/wisp/issues/**` 与 `docs/reports/pending-and-issues.md` 在本腿三枚 commit 的文件清单里出现次数 `0`（`gate-ac3-my-commits.txt`）；AC#1 框**没勾**——勾框归编排者 |
| SLO／阈值 | `internal/observe/thresholds.go` 未出现在任何一枚 diff；为变绿放宽断言：无（唯一红仍是被别家删掉的 CSS 文件那一枚，本腿未动） |
| Git 纪律 | 只 commit 未 push；每发都是 `git add <显式路径> && git commit -F <消息档> -- <显式 pathspec>`；⛔ 未用 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`（M1 还原后那次单文件 `git checkout HEAD -- internal/ball/hotkey_windows.go` 之后本腿改走 `cp` 备份还原，见 §4 末，如实记这一笔）；临时件只建不删（`msg-*.txt`／`ruler*.txt`／`gate-*.txt`／`mutation-*.txt`／两份 `.260r1*backup`／`mutate.py`／`mutation-driver.sh` 全部入库保留） |
| 明文密钥／裸 `go func`／墙钟超时／`filepath.Clean|Abs` 决策／emoji | `sh scripts/d22scan.sh` 终值 `clean - no D22 ban violations`（§3），新代码未引入上述任何一种形状；新 Go 文件里连注释都不带 ban-8 字符带（该档实测扫 `internal/` 500 枚 Go 文件，含注释与 `_test.go`） |
| 冻结三件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`） | 未出现在任何 diff（不在本腿写面） |

---

## §7 Progress log

| 时刻 | 干了什么 | commit |
|---|---|---|
| 2026-10-04 09:4x | 整读票面＋编排者选形裁定节；三把起手尺现跑；未修码基线 56 枚名册落档；骨架 commit | `eb2c0173` |
| 同上 | 产码四道落点（§1）＋七枚常驻判据（§2）＋winlive 两枚只写不跑；门禁首发；突变 M1 单发 | `d9bb6817` |
| 同上 | idle 那句文案补实话（条件②在自己包里的最后一个落点）；三发突变对最终码重跑并取还原读数；`cmd/wisp` 那格先实测后改静态并具名挂欠账；本件满稿 | `b1e59d63` ＋本件 |
