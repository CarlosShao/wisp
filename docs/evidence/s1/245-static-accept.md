# 票 245 · 静态验收（非实现者腿）

**`7f5f36b0d` ｜ A:部分 ｜ B:通过 ｜ C:通过 ｜ D:通过**

- 审的对象：`7f5f36b0d`（`245-r4` 写腿；父 `137bf535c`；审时工作树 HEAD `307becc69`）。
- 本腿角色：**验收腿（非实现者）**，只读。
- ⛔ 本腿**未执行** `go test`、未建窗、未跑 bash 测试载具（派单硬约束：另一条腿占用桌面，同一时刻只许一发）。
  ⇒ 下面凡带"PASS／耗时"的读数都是**实现方自述的既有读数**（证据件 `245-esc-not-a-standby-global-hotkey-r4.md` 里逐字抄的），本腿没有重跑；本腿的判据全部来自 `git show`、读源码、`grep`、`git diff`、`git status` 与哈希。

判定汇总：**A 部分**（AC#8 行为对、但证据点名的那句 ours 判词在 live 路径不可达）／**B 通过**（无偷偷放宽；唯一实义放宽已具名并由票面 13:4x 裁定）／**C 通过**（7 文件、`internal/ball` 全 `_test.go`、D43/阈值/golden/契约零改动）／**D 通过**（台件在册且缺件硬红；裸 Esc 产品注册点唯一）。

---

## A. 声称 vs 盘上

### A① AC#8：两种红分开（硬红 vs 具名跳过＋数字）

**声称**（`docs/evidence/s1/245-esc-not-a-standby-global-hotkey-r4.md:33-40`）：

- `:37` 逐字「`ours` ⇒ 硬红，判词 `ticket 245 RED (AC#8, OUR OWN SET): … the ball's own registration set names it: <槽>=<绑定> (<加速键>)`」；
- `:38` 逐字「`foreign` ⇒ `t.Skipf`，判词 `SKIP-LOUD (ticket 245 AC#8, FOREIGN OWNER): … Probe attempts=3, each answered errno 1409 "Hot key is already registered." …`」；
- `:40` 逐字「`escBorrowProbe` 只把 `ERROR_HOTKEY_ALREADY_REGISTERED`（1409）映射成 `(free=false, err=nil)`，其余拒绝一律作为 error 上抛…"不 free" 这个读数按构造就是 1409」。

**盘上**：

1. `requireIdleRoster` 的**第一道**检查就是名册红（`internal/ball/hotkey_live_test.go:504-506`），逐字：

   ```go
   if reds := idleRosterReds245(rep); len(reds) != 0 {
       t.Fatalf("the live ball's idle roster is not the standby one: %v\n%+v", reds, rep.Bindings())
   }
   ```

2. 裸键在册这件事**已经被 `idleRosterReds245` 红掉**（`internal/ball/hotkey_idle_bare_key_245_test.go:64-68`），逐字：

   ```go
   if b, held := idleBareKeyHeld245(rep); held {
       reds = append(reds, fmt.Sprintf("ticket 245 RED: the idle ball holds %s = %q as a bare key "+
   ```

3. 之后才连问 3 次探针并分类（`hotkey_live_test.go:507-508`，逐字 `p := probeBareEsc245(t, b, 3)` / `switch classifyEscIdleOwner245(rep, p) {`）。

⇒ **结论（关键）**：`case escIdleBall245:` 那一支（`hotkey_live_test.go:511-516`，逐字 `t.Fatalf("ticket 245 RED (AC#8, OUR OWN SET): bare Esc is claimed …`）只有在「名册返回空」且「`idleBareKeyHeld245` 为真」时才会到达，而两者调的是**同一个判据**（`:125` 的 `if _, held := idleBareKeyHeld245(rep); held {`）⇒ **不可能同时成立，这一支在 live 路径是死代码，那句判词跑不出来**。
实到的红句是 `:505` 的名册句（内含 `ticket 245 RED: the idle ball holds cancel = "Esc" as a bare key …`）。**行为没有被放宽**：ours 形仍然硬红，而且发生在探针之前（产品错不会藏进 foreign 支）；错的是**证据件（以及票面 13:4x 裁定）点名的判词引文**。

**foreign 支**（`hotkey_live_test.go:517-522`）逐字：

```go
t.Skipf("SKIP-LOUD (ticket 245 AC#8, FOREIGN OWNER): bare Esc is claimed on this desktop while the "+
    "ball is idle, but the ball's own registration set is clean (%d live entries, no bare key), so the "+
    "claim belongs to another program. Probe attempts=%d, each answered errno 1409 %q. The idle-roster "+
    "half of this case PASSED before the probe; this run cannot judge the delivery half on this desktop.",
    len(rep.Live()), p.attempts, windows.ERROR_HOTKEY_ALREADY_REGISTERED.Error())
```

⇒ **具名＋带数字成立**：`attempts` 是真实计数（`:530-549` 最多 3 次），`errno 1409` 的文本取自常量。

**errno 的构造性成立**：`hotkeyFreeProbe` 只在 `errors.Is(e, windows.ERROR_HOTKEY_ALREADY_REGISTERED)` 时返回 `(false, nil)`，其余 error 原样返回（`hotkey_live_test.go:84-89`）；`probeBareEsc245` 对非 nil error 直接致命（`:535-538`）⇒ 证据 `:40` 的推理与盘上一致。

**分类器本体**：`internal/ball/hotkey_idle_bare_key_245_test.go:117-129`（`free` → `held` → `foreign` 的顺序）；三值控制用例在 `internal/ball/hotkey_esc_idle_owner_245_test.go:42-93`，其 `ours` 断言在 `:89-92`（红句逐字 `hotkey_esc_idle_owner_245_test.go:90: a claimed bare Esc that the ball's own set holds classified as owner %d, want ours (%d) - the foreign-owner skip must never absorb the product's defect`）。
⚠ 这个用例钉的是 **`classifyEscIdleOwner245` 纯函数**，它证明分类器不会把 ours 归 foreign，但**证不到** live 尺会打印 `ticket 245 RED (AC#8, OUR OWN SET)`（该句不可达）。

**结论（A①）**：**部分**。两种红确实分开（ours 硬红且先于探针；foreign 具名跳过带 `attempts` 与 errno 原文），但证据 `:37` 与票面裁定 `:119` 点名的 ours 判词在盘上不可达，实际红句是 `hotkey_live_test.go:505`。**改法二选一**：把名册检查并入分类之后（让 `:513` 可达），或把引文改成 `:505` 的实句。未跑测试，无法复验实现方的 5/5 与 4/4 读数。

### A② `interaction_live_test.go:231` 那处「借哪枚注哪枚」

**声称**（证据 `:68-72`）：红在 `internal/ball/interaction_live_test.go:231` 的注入行，写死注 `"Esc"`；改成 `borrowed := b.ConfiguredHotkeys().Cancel; injectBinding(t, borrowed)`。

**盘上**：

- 改前（`git show 137bf535c:internal/ball/interaction_live_test.go`）：`:231` 逐字 `if !injectBinding(t, "Esc") {`，`:238` 逐字 `t.Fatalf("the taken-over Esc did not cancel Confirming (machine %s)", m.State())`。
- 改后：`:241` 逐字 `borrowed := b.ConfiguredHotkeys().Cancel`，`:242` 逐字 `if !injectBinding(t, borrowed) {`，同一红句现在 `:249`。
- `ConfiguredHotkeys()` 返回 `b.boundCfg`（`internal/ball/ball_windows.go:916-923`，注释逐字「the input of the last registration pass」）；`b.cancelBinding` 同源（`ball_windows.go:286`／rebind `:898`），借的时候正是 `borrow := resolveCancelBorrow(b.cancelBinding)`（`ball_windows.go:966`）⇒ **「借哪枚注哪枚」成立**。夹具 `liveHotkeys()` 的 Cancel 是 `"Ctrl+Alt+V"`（`internal/ball/hotkey_live_test.go:45-47`）。
- 断言面没动：注入失败仍是 `t.Skipf`（`:243`，预存形态），成功则 `:248-249` 的 `m.State() != Acting` 硬 `t.Fatalf`；**没有**把红改 log。

**结论（A②）**：**成立**。⚠ 唯一的账是**引用行号过期**：证据/票面写的 `:231` 是改前行号，现盘上注入在 `:241-242`、红句在 `:249`；读数本身与盘上一致。

### A③ AC#9 四形尺的四形判词与跳过语

**声称**（证据 `:93-126`）：新件 `cmd/wisp/resident_esc_delivery_245_windows_test.go`，四形＝基线 1／自占 0+wm_hotkey 1／出厂常驻仍 1／退出后再 1，复用 `buildEscListener246`／`observeEsc246`（票 246-r1 台件）。

**盘上**（逐形点名）：

| 形 | 位置 | 判词（逐字前缀） | 类型 |
|---|---|---|---|
| E1 基线 | `:103-120` | `:110` `SKIP-LOUD (ticket 245 AC#9, FOREIGN OWNER at baseline): …`；`:116` `ticket 245 AC#9 RED (instrument): …` | 具名跳过＋数字／instrument 硬红 |
| E2 自占控制 | `:122-132` | `:128` `ticket 245 AC#9 RED (instrument blind): with its OWN global bare-Esc registration the observer still received keydown_esc=%d (want 0) and wm_hotkey=%d (want 1)` | 硬红 |
| E3 出厂常驻挂着 | `:134-180` | `:148` `SKIP-LOUD (ticket 245 AC#9): the shipped leg reported %q, so there is no idle ball …`；`:169` `SKIP-LOUD (ticket 245 AC#9, FOREIGN OWNER): …`；`:175` `ticket 245 AC#9 RED (PRODUCT): the shipped no-args process held the desktop's bare Esc while idle …` | 具名跳过（缺球）／具名跳过＋数字／产品硬红 |
| E4 退后再测 | `:182-204` | `:195` `SKIP-LOUD (ticket 245 AC#9, FOREIGN OWNER after exit): …`；`:200` `ticket 245 AC#9 RED: after the shipped leg left, an injected bare Esc reached no foreground window …` | 具名跳过＋数字／硬红 |

- 四形映射 E1–E4 的注释在 `:18-30`；末尾读数行 `:206` 逐字 `t.Logf("AC#9 four shapes (E1/E2/E3/E4): baseline=%d control_steal=%d/%d shipped_idle=%d after_exit=%d", …)`。
- 台件复用：`buildEscListener246`（`cmd/wisp/resident_approval_live_246_windows_test.go:405-436`）、`observeEsc246`（`:452` 起）。
- 实现方自述绿读数（证据 `:114-116`）：`AC#9 four shapes (E1/E2/E3/E4): baseline=1 control_steal=0/1 shipped_idle=1 after_exit=1` / `--- PASS … (19.78s)`。**本腿未复跑**。

**结论（A③）**：**成立**。四形、四类判词、跳过语都在盘上，且 E2 是"0"那一形的正控（证明后面的 1 不是尺子聋）。

---

## B. 有没有偷偷放宽

### B.1 本 commit 新增／改动的 `Skip/Skipf`，逐处列

| # | 位置 | 状态 | 判词 | 判 |
|---|---|---|---|---|
| 1 | `internal/ball/hotkey_live_test.go:518` | **新增** | `SKIP-LOUD (ticket 245 AC#8, FOREIGN OWNER)`，带 `attempts=%d` 与 `errno 1409 %q` | **具名＋带数字** |
| 2 | `internal/ball/interaction_live_test.go:243` | 预存，本 commit 只把参数/文案改成 `%q`（带出 `borrowed`） | `SKIP-LOUD: injected %q never reached the input stream …` | **具名、无数字**（预存形态，触发条件未放宽） |
| 3 | `cmd/wisp/resident_esc_delivery_245_windows_test.go:110` | **新增** | `SKIP-LOUD (ticket 245 AC#9, FOREIGN OWNER at baseline)`，带 `keydown_esc`、`%d/3`、逐次 errno 原文 | **具名＋带数字** |
| 4 | `cmd/wisp/resident_esc_delivery_245_windows_test.go:148` | **新增** | `SKIP-LOUD (ticket 245 AC#9): the shipped leg reported %q, so there is no idle ball …` | **具名、无数字**（环境缺失；原文引 `ballAbsentClaim`） |
| 5 | `cmd/wisp/resident_esc_delivery_245_windows_test.go:169` | **新增** | `SKIP-LOUD (ticket 245 AC#9, FOREIGN OWNER)`，带 `keydown`、`%d/3`、errno 原文 | **具名＋带数字** |
| 6 | `cmd/wisp/resident_esc_delivery_245_windows_test.go:195` | **新增** | `SKIP-LOUD (ticket 245 AC#9, FOREIGN OWNER after exit)`，带 `keydown`、`%d/3`、errno 原文 | **具名＋带数字** |

**裸跳过（无判词、无数字）：0 处。**

### B.2 放宽上界／改宽断言／把红改 log

- **放宽上界：0 处。** `git diff 137bf535c 7f5f36b0d -- internal/ball/hotkey_live_test.go internal/ball/interaction_live_test.go` 里没有改动任何既有等待时长／尝试次数；新件里的 `400`（`resident_esc_delivery_245_windows_test.go:138/:162/:186`）是新写的常驻腿预算，不是把旧上界放宽。
- **改宽断言：0 处。** `interaction_live_test.go:248-249` 的 `m.State() != Acting` 硬红原样保留；`hotkey_live_test.go:504-506` 的名册红原样保留；新件四形全部以 `t.Fatalf` 收口（E2 `:128`、E3 产品 `:175`、E4 `:200`）。
- **把红改 log：0 处。** 新件的两处 `t.Logf`（`:156` 打印 AC#6① 的 live 读数、`:206` 打印四形读数）是**读数记录**，不是把某个断言换成 log。⚠ 既有的 `hotkey_live_test.go:590` 是 `t.Errorf`（不是本 commit 引入，`git diff` 不含该行），本腿不据此判本 commit。

### B.3 唯一实义放宽（不是偷偷）

AC#8 的 `foreign` 支把旧的硬红改成了具名跳过：旧代码是 `hotkey_live_test.go`（父提交）里的 `if !free { t.Fatal("ticket 245 RED: bare Esc is claimed on this desktop while the ball is idle") }`，新代码 `:517-522` 是 `t.Skipf`。这正是票面 AC#8 字面「⛔ 不许改成跳过」与派单的冲突点；证据 `:185` 已具名请裁，票面 `:114-120`（13:4x）已裁定**按派单口径接受**，并给了撤销口令（把 `t.Skipf` 换回 `t.Fatalf`）。⇒ **属公开裁定后的放宽，不是偷偷放宽**。

**代价（记录，不改判）**：跳过发生在 `requireIdleRoster` 内部（`:518`），所以四个 winlive 用例会**整体 SKIP**，探针之后的断言不再跑（证据 `:57-64` 的读数正是这个形状）；`:148` 的缺球跳过连数字都不带（只带 `ballAbsentClaim` 原文），但它跳的是"没有测量对象"的环境形，不是产品读数。

**结论（B）**：**通过**——所有新增跳过都具名，五个实质读数跳过都带数字，0 处裸跳过、0 处放宽上界、0 处改宽断言、0 处把红改 log；唯一硬红→跳过的变化已公开裁定。

---

## C. 越界与语义

### C.1 7 个文件 vs 声称写面

`git show 7f5f36b0d --name-status` 现跑（7 枚，逐字）：

```
M	.scratch/wisp/issues/245-the-bare-esc-cancel-default-is-a-global-hotkey-so-the-resident-process-steals-esc-from-every-other-app.md
A	cmd/wisp/resident_esc_delivery_245_windows_test.go
A	docs/evidence/s1/245-esc-not-a-standby-global-hotkey-r4.md
A	internal/ball/hotkey_esc_idle_owner_245_test.go
M	internal/ball/hotkey_idle_bare_key_245_test.go
M	internal/ball/hotkey_live_test.go
M	internal/ball/interaction_live_test.go
```

证据 `:6-12` 的写面清单列了 **6 枚**（5 个代码/测试件＋证据件本身），**漏列票面那一行**（第 7 枚）。票面 diff `git diff --numstat 137bf535c 7f5f36b0d` 对票面是 **1 增 0 删**（只追加 Progress log 一行）⇒ 与派单允许的"票面末尾追加"一致，不是正文改写。

**结论**：写面与声称只差"票面 1 行未列入清单"，**无语义越界**。

### C.2 `internal/ball/**` 是否全 `_test.go`

本 commit 触到的 `internal/ball` 四件：`hotkey_esc_idle_owner_245_test.go`（新）、`hotkey_idle_bare_key_245_test.go`、`hotkey_live_test.go`、`interaction_live_test.go`——**全部 `_test.go`，非测试文件 0 处**。与票面 `:123` 的裁定逐字一致（「四件全是 `_test.go`，非测试文件 0 处」）。

### C.3 D43 转移表／阈值／golden／契约

- `git show 7f5f36b0d --name-only` 里**没有** `docs/PLAN.md`、`docs/specs/**`、`internal/observe/thresholds.go`、任何 golden／allowlist／`internal/statemachine/**` 文件。
- commit 正文里 `D43` 只出现在**注释**（`internal/ball/interaction_live_test.go:245` 逐字 `"cancel id and the D43 #22 veto were asserted above; the B1 return is asserted below by the "`）与证据文字（`:13`／`:185`）；**转移表、阈值、golden、契约零字节改动**。
- 生产字节还原核（现跑）：`git status --short -- cmd/wisp/resident_ball_windows.go cmd/wisp/resident_windows.go` **空**；`(Get-FileHash cmd/wisp/resident_ball_windows.go).Hash` = `1D47ECCB205328A744FB1FF4507EF6962FA8CB29826671BAF2C5CE9406178425`，与证据 `:161` 逐字同串。⇒ 正控的临时变异确已回到 HEAD。

**结论（C）**：**通过**。

---

## D. 逃逸面

### D① 第二进程台件是否真在仓里

`git ls-files cmd/wisp/testdata/esclistener` 现跑，逐字：

```
cmd/wisp/testdata/esclistener/main.go
```

盘上同一路径 10,226 字节；`main.go` 的 `-watch`／`-steal` 两模式在 `:162-166`（逐字 `watch := flag.Bool("watch", false, "no global hotkey here; the injected Esc must be delivered")`、`steal := flag.Bool("steal", false, "register bare Esc globally here; the injected Esc must NOT be delivered")`），`READY`／`READING` 行在 `:265`／`:293`。⇒ **在册、在盘、可被 `go build`。**

### D② 台件缺失/不可执行时会静默跳过吗

**不会，全部硬红**：

- `buildEscListener246` 现读：`:425` 逐字 `cmd := exec.Command(goExe, "build", "-o", exe, "./cmd/wisp/testdata/esclistener")`；`:429-430` 逐字 `if err := cmd.Run(); err != nil { t.Fatalf("go build ./cmd/wisp/testdata/esclistener: %v\n%s", err, log.String()) }`；`:432-433` 无 exe 也 `t.Fatalf`。
- `observeEsc246`：`:465-466` `cmd.Start()` 失败 `t.Fatalf`；`:491` 观察者早退 `t.Fatalf`；`:517-518` 20s 无 `READING` `t.Fatalf`。
- `observeEsc246` 还要求 `foreground_is_mine=true`（`:506-509`），否则直接红——不会把"没拿到前台"读成 0。

⇒ 缺件、编译失败、起不来、不出读数**都走 `t.Fatalf`**，不会变成 skip。

### D③ 裸 `Esc` 的注册点是否只有一处

全部 `RegisterHotKey` 调用点（`git ls-files '*.go'` 现跑，按面分类）：

| 面 | 位置 | 说明 |
|---|---|---|
| **产品（唯一 Win32 调用）** | `internal/ball/hotkey_windows.go:414` | `pRegisterHotKey.Call(uintptr(hwnd), uintptr(id), …)`，在 `hotkeyRegisterer` 内 |
| **产品（idle 名册）** | `internal/ball/hotkey_windows.go:480-490` | `registerAllWith` 绑 summon/mute/panel；`if p.id == hkCancel { … cancelIdleLine … continue }`（逐字注释 `:485-488`「The one slot an idle ball must not register. takeEscWith is the only caller that ever binds this id」） |
| **产品（借用＝唯一 cancel/裸 Esc 注册路径）** | `internal/ball/ball_windows.go:961-992` → `hotkey_windows.go:631-633`（`takeEscBorrow`）→ `:609-617`（`takeEscWithAcc`）→ `:414` | 逐字 `:948-950`「Ticket 245: this is the ONLY path that registers the cancel id. An idle ball leaves it standby (see registerAllWith)」；产品调用者 `cmd/wisp/resident_approval_windows.go:924`，归还 `:1026` |
| 调试缝（仅测试调用） | `internal/ball/ball_windows.go:1037-1040` | `DebugRegisterHotkeySquat`；调用者只有 `cmd/wisp/resident_hotkey_258_windows_test.go:456/:462` |
| 测试/台件 | `cmd/wisp/testdata/esclistener/main.go:243`；`cmd/wisp/resident_esc_delivery_245_windows_test.go:86`；`internal/ball/hotkey_live_test.go:71/:603`；`cmd/wisp/resident_hotkey_live_258_windows_test.go:273`；`cmd/wisp/resident_audio_290_live_windows_test.go:162` | 探针/占位/检查用 |
| spike（非产品） | `scripts/spike/common/winshell.go:437` | 试验工具，不进程产品 |

⇒ **产品里裸 Esc 的注册点只有一处**：`Confirming` 借用（`TakeEscForCancel` → `takeEscBorrow` → `takeEscWithAcc`）；稳态由 `registerAllWith` 显式跳过 cancel。**成立**。

**残余（不改判，记录在案）**：`requireIdleRoster` 的 foreign 判定以"报告集干净"为唯一依据（`hotkey_live_test.go:517-522`）；一个**掉出报告的**同进程注册会被读成"外人占"——代码自己在 `internal/ball/hotkey_windows.go:525-528` 写了这个形状（「a slot we lost track of … would come back on the NEXT rebind as a 1409 that reads "somebody else owns it" while the somebody is us」）。这是探针的固有限制，不是本 commit 引入。

**结论（D）**：**通过**。

---

## 一句话判词

- **AC#8 能不能翻**：**能翻（按 13:4x 裁定口径）**——我们自己的注册集握着裸键仍硬红且先于探针，键被外人占是具名跳过带 `attempts` 与 errno 原文；但**翻之前必须把证据 `:37` 与票面 `:119` 的 ours 判词改成 `hotkey_live_test.go:505` 的实句**（`ticket 245 RED (AC#8, OUR OWN SET)` 在 winlive 路径不可达），否则验收会去 grep 一句永远不会出现的红句。
- **AC#9 能不能翻**：**能翻**——四形、四类判词、台件（`cmd/wisp/testdata/esclistener/main.go` 在册）与"缺件即硬红"的失败面都在仓内；`19.78s` 的绿读数是写腿自述，本腿按只读约束未复跑。

## 本腿未做的事（如实）

1. 未跑任何 `go test`／winlive／`cmd/wisp` 定向用例（派单硬约束）。
2. 未复验实现方的 5/5、4/4、19.78s、正控红句与 SHA 还原读数；这些只按"证据件里已有的逐字读数"引用，并已逐条标为自述。
3. 未 commit、未 push；写面仅：本文件、票面末尾一节、`.scratch/wisp/probes/audit245/**`。
