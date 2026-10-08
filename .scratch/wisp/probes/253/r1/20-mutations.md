# 253-r1 · 盘上种＋`git hash-object` 拉平自证（20-mutations）

时刻 `2026-10-08 16:53 +0800` 起算（各发自己的 `date`／`time=` 戳逐条附后）。

★**本程方法（A713 定式，`A716` 收下、`probes/ruler-dedup-1/readings.md:65` 那条 255r6 件头复述）**：这枚尺读盘（`os.ReadDir`＋`os.ReadFile`＋`parser.ParseFile`），所以**它的绿／红只有"盘上种"能证**，`go test -overlay` 只换编译器的眼、换不到尺的眼。下面每一发都是：
`git hash-object`（种前）→ Edit 工具改 `cmd/wisp/panel_host_windows.go` 那一行 → 靶向 `go test` 取红句原文 → `git show HEAD:cmd/wisp/panel_host_windows.go > cmd/wisp/panel_host_windows.go` 还原 → `git hash-object`（种后）→ `git status --porcelain -- cmd/wisp`。
**三枚对拉（种前 hash＝还原后 hash＝`58e2b155dbcc8d380793e90bffc5514ab00743f1`）逐字相等**，且每发还原后 `cmd/wisp` 的 porcelain 都是空——包括我第 1 笔提交之后（我的新测试文件已入库，porcelain 里没有它）。
⛔ 全程没动第二枚产码文件；⛔ 没动 `internal/panel/**`；⛔ 没动 `internal/agent/approval/**`（那格此刻被编排者按住，见派单"门禁与并发纪律"）。

---

## M0 基线（种之前＝还原之后，反复量到 7 次同值）

```
git hash-object cmd/wisp/panel_host_windows.go   = 58e2b155dbcc8d380793e90bffc5514ab00743f1
git rev-parse HEAD:cmd/wisp/panel_host_windows.go = 58e2b155dbcc8d380793e90bffc5514ab00743f1
```
尺在干净树上：**PASS／rc=0**（见 `10-gates.md` §4）。这一发就是派单要的 **M3 正控**（"合法绑定必须绿"）的现场读数，同时尺日志逐字印出 `over 34 production file(s)`＝整包射程。

---

## M1 绑定名 ≠ 名册（种在 `:802`）★派单点名的那一发

- 种前 hash `58e2b155dbcc8d380793e90bffc5514ab00743f1`
- 盘上改的行＝`cmd/wisp/panel_host_windows.go:802`：
  `if err := w.Bind(panelDispatchBinding, func(raw string) string {` → `if err := w.Bind("wispStaleDoor", func(raw string) string {`
  （种后 `sed -n '802p'` 复量＝`	if err := w.Bind("wispStaleDoor", func(raw string) string {`；CR 字节＝0）
- 命令＝`go test -count=1 ./cmd/wisp/ -run 'TestTransportDoorBindingMatchesRoster253r1'`，时刻戳 `time=2026-10-08T15:53:19+08:00`（首发）／复跑 `M1 rc=1`
- **rc=1**（复跑那次逐字 echo 的就是这枚数）
- 指名红句原文（两枚，逐字）：
  ```
  panel_dispatch_binding_roster_253r1_windows_test.go:377: AC#1 RED: installPanelTransport binds transport door "wispStaleDoor" at panel_host_windows.go:802, which is not in doorRoster253r1 - a new inbound page-addressable door entered the tree without being reviewed here
  panel_dispatch_binding_roster_253r1_windows_test.go:383: AC#1 RED: doorRoster253r1 reviews door "wispDispatch", but installPanelTransport (panel_host_windows.go:801) binds no such name (bound: [wispStaleDoor]) - a rostered door with no Bind behind it is the dead door 票 253 AC#2 refused to call registration
  ```
  ⇒ 红句自己带 `file:line`（`panel_host_windows.go:802`＝种的那行、`:801`＝传输函数声明行）。
- 还原后 hash `58e2b155dbcc8d380793e90bffc5514ab00743f1`，porcelain 空。

## M2 传输一枚门都没绑（空名册 ⇒ 必须 Fatal 而不是绿）

- 种前 hash `58e2b155…43f1`；种＝`cmd/wisp/panel_host_windows.go:802` 改成 `w.Bind("", func(raw string) string {`（`sed -n '802p'` 复量到；CR＝0）
- 命令同上；时刻戳 `time=2026-10-08T16:02:21+08:00`（首发）／复跑 `M2 rc=1`
- **rc=1**
- 红句原文（是 `t.Fatalf`，不是 `t.Errorf`——尺日志仍先印出射程，然后 Fatal）：
  ```
  panel_dispatch_binding_roster_253r1_windows_test.go:358: 253-r1 census: door constant "panelDispatchBinding" = "wispDispatch"; installPanelTransport at panel_host_windows.go:801 made 1 Bind call(s) and 1 Init call(s); roster holds 1 name(s) over 34 production file(s)
  panel_dispatch_binding_roster_253r1_windows_test.go:364: installPanelTransport (panel_host_windows.go:801) issued 1 Bind call(s) but bound zero readable door names: the transport roster is empty, and a page posting to any door lands on nothing
  ```
  ⇒ **`Bind` 调了但读不出任何门名＝Fatal**，不是"零枚＝clean"。这就是派单要的"摘掉尺的读盘步／空名册 ⇒ 必须 Fatal"那一发（读盘步本身没摘，因为摘它得改尺自己；等价形状＝名册空转，尺自己的 `t.Fatalf` 在 `:364`）。
- 还原后 hash `58e2b155…43f1`，porcelain 空。

## M3 页面半与 Go 半解绑（`w.Init` 的脚本不再由那枚常量生成）

- 种前 hash `58e2b155…43f1`；种＝`cmd/wisp/panel_host_windows.go:792`：
  `})();`, panelDispatchBinding)` → `})();`, "wispOther")`（`sed -n '792p'` 复量到）
- 命令同上；时刻戳 `time=2026-10-08T16:03:52+08:00`（首发）／复跑 `M3 rc=1`
- **rc=1**
- 指名红句原文：
  ```
  panel_dispatch_binding_roster_253r1_windows_test.go:398: AC#1 RED: installPanelTransport calls w.Init([panel_host_windows.go:808]) with no script built (via fmt.Sprintf) from "panelDispatchBinding" - the page-side door the forwarding hook addresses is not the constant the transport binds, so the two can drift apart and leave the page posting at an unbound door
  ```
  ⇒ 这一发证明**第四支断言不是装饰**：名册、常量、`Bind` 三处当时**全都自洽**（`"wispDispatch"` 仍是绑定的名），只有 JS 那半被换成了字面量 `"wispOther"`；尺照样红，而且红句点名 `w.Init` 的调用现场 `panel_host_windows.go:808`。
- 还原后 hash `58e2b155…43f1`，porcelain 空。

## M4 只改常量的"值"（`A717` 那句"改名"最窄的一发）

- 种前 hash `58e2b155…43f1`；种＝`cmd/wisp/panel_host_windows.go:80`：
  `panelDispatchBinding = "wispDispatch"` → `panelDispatchBinding = "wispDispatchV2"`（`sed -n '80p'` 复量到）
- 命令同上；复跑 `M4 rc=1`，`grep -cE 'AC#1 RED'`＝**3**
- 三枚红句原文（逐字，行号＝尺自己的行）：
  ```
  :377: AC#1 RED: installPanelTransport binds transport door "wispDispatchV2" at panel_host_windows.go:802, which is not in doorRoster253r1 ...
  :383: AC#1 RED: doorRoster253r1 reviews door "wispDispatch", but installPanelTransport (panel_host_windows.go:801) binds no such name (bound: [wispDispatchV2]) ...
  :391: AC#1 RED: constant "panelDispatchBinding" resolves to "wispDispatchV2" (panel_host_windows.go), which is not in doorRoster253r1 - the transport door's reviewed name and its declared value have come apart
  ```
- 还原后 hash `58e2b155…43f1`，porcelain 空。
- ★**这一发推翻了我自己起手锚里的一句话**（我原写"只改常量的值、三处一起漂＝本尺按设计看不见"）：三处确实一起漂（`:792`/`:802` 都按标识符取名字），但**名册是盘上独立的第三枚锚**，所以照样红。正确表述＝"改常量的值 **且同批改 `doorRoster253r1`** 才是绿的，而那正是'已评审动作'"。更正具名写在 `01-anchor-correction.md`，⛔ 不回改已提交的锚件与第 1 笔 commit 消息。
- ⚠ 顺带说明**尺按设计放过的形状**：`go:build windows` 那族尺在 POSIX 分母为零（`probes/ruler-dedup-1/ruler-dedup-1.md:24` 那条我采），本尺同样带 `//go:build windows`，Linux 那几步编不到它——这不是我这格能修的（改档位＝动别的腿的门禁面），记为已知射程边界。

## M5 ★尺看不见的形状，真种了一发（绿＝恒真面实证，不是嘴说）

- 种前 hash `58e2b155…43f1`；种＝`cmd/wisp/panel_host_windows.go:790` 的 JS 文本：
  `try { return window.%[1]s(message); } finally { inside = false; }` → `try { return window.wispNothing(message); } finally { inside = false; }`
  （`fmt.Sprintf` 的**实参一枚没动**，`:788` 那半句 `typeof window.%[1]s` 也留着，所以**编译真过了**——尺跑到并印出日志就是证据；`go vet` 我当时接了 `| tail -3` 只看到零输出、⛔ 没取 rc，所以这里只记"编译过＋vet 零输出"，不记 vet 的 rc）
- 命令＝靶向同一枚尺，**两发**：
  ①先跑 `-v` 那发（种形在场）逐字含 `--- PASS: TestTransportDoorBindingMatchesRoster253r1 (0.02s)` ＋ `panel_dispatch_binding_roster_253r1_windows_test.go:401: 253-r1 lockstep: w.Init(panelPostMessageForwardInit) at panel_host_windows.go:808 derives its window.<door> reference from "panelDispatchBinding", the same constant the Bind uses`，⛔ 那发我没取 rc；
  ②再跑不带 `-v` 那发取 rc，逐字 echo：
  ```
  M5 rc=0 (expect 0 = blind spot proven)
  ok  	github.com/CarlosShao/wisp/cmd/wisp	0.130s
  ```
  并且同发里 `git diff --numstat -- cmd/wisp/panel_host_windows.go`＝**`1 1`**（种形确实在盘上、只动那一行），`git hash-object` 还原后回到 `58e2b155…43f1`。
- **rc=0／尺 PASS**——而这一刻页面转发的门已经是 `window.wispNothing`（Go 从没绑过它）。
  ⇒ **本尺的恒真面之一被实证了**：它断的是"`w.Init` 交出去的那段脚本由那枚常量生成"（实参层锁步），⛔ 它**不读 JS 模板文本里 `%[1]s` 有没有真被用到**。那一格今天仍只有能力形那半守着（`cmd/wisp/panel_transport_35r1_test.go:124-131` 那把读 `Init` 串里含不含 `window.wispDispatch` 的尺＋`:195→:203` 真调 `installPanelTransport`）。我**没有**把它算成自己的功劳，也没为它扩射程（扩了就会与那枚现存尺重复＝`A717` §5 明令禁止的那格）。
- 还原后 hash `58e2b155…43f1`，porcelain 空。

## M6 名字用间接表达式取（尺必须 fail-closed，不许"读不出＝过")

- 种前 hash `58e2b155…43f1`；种＝`cmd/wisp/panel_host_windows.go:802` → `w.Bind((panelDispatchBinding), func(raw string) string {`（多一对括号＝AST 上是 `*ast.ParenExpr`，我的解析器只认 `BasicLit`/`Ident` 两种形）
- ⚠ 我**第一发 M6 没有读数，不作数**（具名自抓）：那发我用 `strings.Clone(panelDispatchBinding)` 种，并在同一条命令里先跑了 `sed -n '1,45p' … | grep -n "strings"` 查 import——grep **零命中 ⇒ rc=1 ⇒ 把后面的 `&&` 链整条打断**，`go test` 一枚字都没执行（证据＝那一发的 `/tmp/m6.out` "No such file or directory"，而我 echo 出来的 `M6 rc=1` 是**断链的 rc、不是尺的 rc**）。这正是本仓 dispatch 里记过的那枚坑（"`grep -c` 命中 0 会 rc=1 并吃掉 `&&` 链"）在腿侧落到我头上一次。
  下面这发（括号形＝`*ast.ParenExpr`，不需要新 import）才是 M6 的唯一读数。顺带那把 import 尺的**事实**仍然成立（现量 `awk '/^import \(/,/^\)/' cmd/wisp/panel_host_windows.go` ＝ `context`/`errors`/`fmt`/`log/slog`/`sync`/`time`/`unsafe` ＋ `internal/panel`/`go-webview2`/`x/sys/windows`，**无 `strings`**）——但那只是我换形的理由，⛔ 不是"编译门红"的测量，我没测过那一发编不编得过。
- 命令＝同上；复跑逐字 echo：
  ```
  M6 rc=1
  --- FAIL: TestTransportDoorBindingMatchesRoster253r1 (0.01s)
      panel_dispatch_binding_roster_253r1_windows_test.go:364: installPanelTransport (panel_host_windows.go:801) issued 1 Bind call(s) but bound zero readable door names: the transport roster is empty, and a page posting to any door lands on nothing
  ```
- **rc=1**，红句落在**同一枚 `:364` 的 `t.Fatalf`**（空名册那支）＝读不出名字时尺**关死**，不会把"我解析不了"洗成"没有门"之外的第二种输出。⚠ 代价具名：报错文案点名的是"zero readable door names"，不会告诉后来者是"被摘掉"还是"被间接化"（这是把尺拆成两种文案要动尺自己，本程没做；见 `30-blind-spots.md` #2）。
- 还原后 hash `58e2b155…43f1`，porcelain 空。

## M7 反证：`-overlay` 给这枚尺的绿不算数（A713 方法的现场复现）

- 做法＝⛔ 不碰工作树文件，把 `:802` 改成 `"wispOverlayOnlyDoor"` 的那份**拷到仓库外**（`C:/Users/swq/AppData/Local/Temp/overlay-mut-253r1.go`，46,251 字节），overlay JSON 也放仓外（同一目录，154 字节），然后：
  `go test -count=1 -overlay 'C:\Users\swq\AppData\Local\Temp\overlay-253r1.json' ./cmd/wisp/ -run 'TestTransportDoorBindingMatchesRoster253r1'`
- 盘上文件同时复量＝`git hash-object cmd/wisp/panel_host_windows.go`＝`58e2b155…43f1`（没动过）。
- 读数逐字：
  ```
  overlay-run rc=0 (expect 0: ruler blind to overlay)
  ok  	github.com/CarlosShao/wisp/cmd/wisp	0.079s
  ```
- **rc=0／尺绿**，而编译器此刻看到的那份产码里门名是 `wispOverlayOnlyDoor`。⇒ 本尺**对 overlay 结构性失明**＝`A713` 那条定式在这枚新尺上现场复现一次（原话参照＝提交 `3ee9b3ff` 的票 197 口径与 `probes/ruler-dedup-1/ruler-dedup-1.md:60`）。
  ⇒ 反过来也说明：**M1–M6 只有走盘上种才算数**，我没用任何 overlay 绿当凭据。
- 两枚临时件只建不删（`issues/README` 规则 8），且在**仓库外**，⛔ 不进 d22scan 分母。

---

## 尺自带的四枚 fixture（正控／负控，overlay 可见那一支，不当读盘凭据）

`TestBindingRosterBitesItsOwnFixtures253r1` 的四个子用例在干净树上全 PASS（`10-gates.md` §4），逐子用例读数逐字：

| 子用例 | 尺日志原文（截取） | 判 |
|---|---|---|
| good | `const="wispDispatch" bound=[wispDispatch] initTied="panelPostMessageForwardInit" reds=[]` | 合法绑定＝**0 红**（不误咬） |
| drifted-bind | `const="wispDispatch" bound=[wispStaleDoor] ... reds=[rostered but unbound door wispDispatch unrostered bound door wispStaleDoor]` | 2 红＝会咬 |
| empty-bind | `... census bound zero readable doors (1 Bind call(s)) - the disk case reddens this as the empty-roster Fatalf` ＋ `reds=[rostered but unbound door wispDispatch]` | 1 红＋空名册那一支在盘上走 Fatal（M2 实证） |
| untied-init | `const="wispDispatch" bound=[wispDispatch] initTied="" reds=[w.Init script not tied to the door constant]` | 1 红＝第四支有牙 |

⛔ 这四枚**不替代**盘上种：它们解析的是本文件自带的源码串（overlay 可见那一支，`probes/ruler-dedup-1/ruler-dedup-1.md:65` 给 33r11 记的就是这一形），只证"匹配逻辑会咬／不误咬"；读盘那半的证据是 M1–M6＋M7。
