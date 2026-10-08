# 253-v1 · 盘上对抗读数（10-attacks，每发四件套）

时刻 2026-10-08 18:19–18:31 +0800（`date` 现跑起算）。腿名 `253-v1`（非实现者裁决腿）。被验收件＝`cmd/wisp/panel_dispatch_binding_roster_253r1_windows_test.go`（579 行，`65f4c968`），两枚顶层用例。

## 方法（与 r1 同规，本腿自执行）

- 命令形：`cd cmd/wisp && PATH="$PWD/../../third_party/sherpa-onnx:$PWD/../../third_party/onnxruntime:$PATH" go test -count=1 . -run '253r1' [-v] 2>&1 | tail -N; echo rc=${PIPESTATUS[0]}`（PIPESTATUS[0]＝go test 的 rc，不是 tail 的；未注 PATH 会 `0xc0000135` 且 0 条 `=== RUN`，本腿没出过这形）。
- 突变只种在盘上；每发：种前 `git hash-object` → 改行用 Edit 工具 → 改后 `sed -n`/`grep -n` 复量改行 → 靶向跑取 rc＋红句原文 → `git show HEAD:<file> > <file>` 还原 → `git hash-object` 复量 + `git status --porcelain -- cmd/wisp`。
- ⛔ 未用 `-overlay` 喂任何突变（r1 的 M7 已证本尺对 overlay 结构性失明，本腿不重做）。
- 基线（16 次工具调用处）：三个目标文件 worktree hash＝HEAD hash 逐字相等（产码 `58e2b155dbcc8d380793e90bffc5514ab00743f1`、`panel_inbound.go` `dd796dc12bb0e73e3f856250abe5d4903b116c9b`、尺 `f8dadd1567087bac4ad33b9016c957e07bfd2121`），porcelain 空；干净树自跑 `-run '253r1'` **rc=0**，日志逐字含 `1 Bind call(s) and 1 Init call(s); roster holds 1 name(s) over 34 production file(s)`。

---

## V1＋V2（方向②：`Init` 那一跳的实参换成另一枚 var）★本腿核心

**种法（一处 Edit 同时两改）**：`cmd/wisp/panel_host_windows.go` 在 `:792` 后插入两枚 var，并把 `:808` 的 `w.Init(panelPostMessageForwardInit)` 改指 V1：
- `panelPostMessageForwardInitV1 = fmt.Sprintf(`…window.%s…window.%s…`, panelDispatchBinding, "wispNothing")`（两个 `%s` 顺序消费两枚实参：第一处＝typeof 检查→`window.wispDispatch`，第二处＝真转发→`window.wispNothing`）
- `panelPostMessageForwardInitV2 = fmt.Sprintf(`window.%s(message)`, "wispNothing")`（args 里一枚 ident 都没有）

**种前 hash**＝`58e2b155dbcc8d380793e90bffc5514ab00743f1`（＝HEAD）。
**改后复量（读 1 在盘时）**：`grep -n` 逐字三行——`794:var panelPostMessageForwardInitV1 = fmt.Sprintf(`(function () {`／`807:var panelPostMessageForwardInitV2 = fmt.Sprintf(`window.%s(message)`, "wispNothing")`／`823:	w.Init(panelPostMessageForwardInitV1)`；`grep -c 'fmt.Sprintf'`＝4。

**读 1（Init→V1，恒真面打点）**：
```
ok  	github.com/CarlosShao/wisp/cmd/wisp	0.081s
V1(read1: Init->trap var with const ident) rc=0
```
**rc=0 全绿**——而此刻若走这段脚本，页面调用会被转发到 `window.wispNothing`（Go 从没绑过它；`absorbFile` 只认「该 var 名下任何 `fmt.Sprintf` 实参里出现过 `panelDispatchBinding` 这一 ident」，不看 ident 是否在格式里生效）。生产调用者＝`panel_host_windows.go:408`（`bringUp` 内 `m.installPanelTransport(w, ctx)`，在 `serveEntry` 前），即这形态在生产路径上。
⇒ **推翻 r1 `30-blind-spots.md` 末条自报**（「`Init` 传进来的不是那枚 var 而是别的 var（尺会红在 `:398`）：⛔ 没种过」）：换成别的 var **不一定红**。

**读 2（Init→V2，对照）**：Edit 把 `:823` 改指 V2 后——
```
panel_dispatch_binding_roster_253r1_windows_test.go:358: 253-r1 census: door constant "panelDispatchBinding" = "wispDispatch"; installPanelTransport at panel_host_windows.go:816 made 1 Bind call(s) and 1 Init call(s); roster holds 1 name(s) over 34 production file(s)
panel_dispatch_binding_roster_253r1_windows_test.go:398: AC#1 RED: installPanelTransport calls w.Init([panel_host_windows.go:823]) with no script built (via fmt.Sprintf) from "panelDispatchBinding" - the page-side door the forwarding hook addresses is not the constant the transport binds, so the two can drift apart and leave the page posting at an unbound door
V2(read2: Init->var whose Sprintf has no const ident) rc=1
```
**rc=1**。⇒ 同一位置、同一「换成别的 var」动作，**只因那枚 var 的 Sprintf args 里有没有 `panelDispatchBinding` 一字，红/绿翻转**。

**还原后 hash**＝`58e2b155dbcc8d380793e90bffc5514ab00743f1` 逐字回、`porcelain rc=0`（空）。

## V3（方向③：门名改成另一种拼法但语义相同——字符串拼接）

**种法**：`panel_host_windows.go:802` `w.Bind(panelDispatchBinding, …)` → `w.Bind("wisp"+"Dispatch", …)`。
**种前 hash**＝`58e2b155…43f1`。**改后复量**：`sed -n '802p'` 逐字＝`	if err := w.Bind("wisp"+"Dispatch", func(raw string) string {`。
```
panel_dispatch_binding_roster_253r1_windows_test.go:364: installPanelTransport (panel_host_windows.go:801) issued 1 Bind call(s) but bound zero readable door names: the transport roster is empty, and a page posting to any door lands on nothing
V3(Bind->string concat) rc=1
```
**rc=1**（`*ast.BinaryExpr` 读不出名 ⇒ 空名册支 `t.Fatalf`，fail-closed）。与 r1 的 M6（ParenExpr）同支、不同种法。
**还原后 hash**＝`58e2b155…43f1`、porcelain 空。

## V5a＋V5b（方向①：名册里的其它名能不能被利用）

**种法**：V5a＝改**尺文件** `doorRoster253r1` 加一枚键 `"wispProbeRT"`（产码不动）；V5b＝V5a 在场时再在产码 `installPanelTransport` 内补一行 `w.Bind("wispProbeRT", func(raw string) string { return raw })`。
**种前 hash**＝产码 `58e2b155…43f1`、尺 `f8dadd1567087bac4ad33b9016c957e07bfd2121`。**改后复量**：`sed -n '86,89p'` 逐字含新键行；V5b `sed -n '802,812p'` 逐字含补绑行。

**V5a 读数**（`-run '253r1'` 两枚）：
```
--- FAIL: TestBindingRosterBitesItsOwnFixtures253r1/good (0.00s)
    panel_dispatch_binding_roster_253r1_windows_test.go:567: fixture "good" reported 1 red(s) [rostered but unbound door wispProbeRT], want 0
--- FAIL: TestBindingRosterBitesItsOwnFixtures253r1/drifted-bind (0.00s)
    :567: fixture "drifted-bind" reported 3 red(s) [rostered but unbound door wispDispatch rostered but unbound door wispProbeRT unrostered bound door wispStaleDoor], want 2
--- FAIL: TestBindingRosterBitesItsOwnFixtures253r1/empty-bind (0.00s)
    :567: fixture "empty-bind" reported 2 red(s) [rostered but unbound door wispDispatch rostered but unbound door wispProbeRT], want 1
--- FAIL: TestBindingRosterBitesItsOwnFixtures253r1/untied-init (0.00s)
    :567: fixture "untied-init" reported 2 red(s) [rostered but unbound door wispProbeRT w.Init script not tied to the door constant], want 1
V5a(roster+othername, no Bind) rc=1
```
盘上用例单枚复跑（同一发内补量）：`:384: AC#1 RED: doorRoster253r1 reviews door "wispProbeRT", but installPanelTransport (panel_host_windows.go:801) binds no such name (bound: [wispDispatch]) - a rostered door with no Bind behind it is the dead door 票 253 AC#2 refused to call registration`，rc=1。⇒ **名册不能空堆名**（两重红）。

**V5b 读数**（两枚一起跑）：
```
--- PASS: TestTransportDoorBindingMatchesRoster253r1 (0.03s)
    :359: census: … installPanelTransport at panel_host_windows.go:801 made 2 Bind call(s) and 1 Init call(s); roster holds 2 name(s) over 34 production file(s)
--- FAIL: TestBindingRosterBitesItsOwnFixtures253r1 (0.00s)   （4 子用例仍全错位红，同上）
V5b(bind+roster both grown) go-test rc=1
```
⇒ **盘上用例单枚可被「名册加名＋补绑定」洗绿（PASS）**；但整腿（两枚）仍 rc=1——**自夹具把名册键集合钉住**（多一枚未绑名 ⇒ 每枚夹具 wantRed 全错位）。即 r1 自报 #2 的「改名＋同批改名册=绿」按**整腿**读数**不成立**（那次它标注为「推理」）。
**还原后 hash**＝产码 `58e2b155…43f1`、尺 `f8dadd15…d2121` 逐字回、porcelain 空。

## V6（方向④：整段复制到同包另一文件、原处留空壳）

**种法**：`panel_host_windows.go:801-810` 方法体改为空壳（`return nil`）＋ `panel_inbound.go` 尾部（310 行后）粘贴复制版挂到影子类型（`type panelShadowV1 struct{}` ＋ `func (s *panelShadowV1) installPanelTransport(w pageTransport) error {…}`；两文件都改、不新建文件）。
**种前 hash**＝产码 `58e2b155…43f1`、inbound `dd796dc12bb0e73e3f856250abe5d4903b116c9b`。**改后复量**：`sed -n '801,803p'` 逐字＝空壳三行；`sed -n '311,325p'` 逐字＝粘贴段（含 `func (s *panelShadowV1) installPanelTransport(w pageTransport) error {`）。
```
panel_dispatch_binding_roster_253r1_windows_test.go:353: found 2 functions named installPanelTransport across the package ([panel_host_windows.go:801 panel_inbound.go:317]); its own doc comment calls it the ONE place the transport is wired, and two of them is a split the roster cannot cover
V6(copy to other file + shell left) rc=1
```
**rc=1**（`len(sites) > 1` 支 Fatal）⇒ **不能骗过**。此发补上 r1 自报「没造出反形」的 sites>1 支的**盘上实证**。
**还原后 hash**＝产码 `58e2b155…43f1`、inbound `dd796dc1…6c9b` 逐字回、porcelain 空（终态）。

---

## 台账（五组、三读红两读绿＋一对照，全部还原）

| 发 | 种处 | 预期 | 实测 rc | 关键红句/日志 |
|---|---|---|---|---|
| 基线 | — | 绿 | 0 | 34 files／1 Bind／1 Init／roster 1 |
| V1 读1 | `:808`→V1（ident 挂名） | 绿 | **0** | 全绿；转发已改道 `window.wispNothing` |
| V1 读2 | `:808`→V2（无 ident） | 红 | **1** | `:398` |
| V3 | `:802` 拼接 | 红 | **1** | `:364` Fatal |
| V5a | 尺名册+wispProbeRT | 红 | **1** | `:384`＋4 夹具错位 |
| V5b | ＋产码补绑 | 单枚绿／整腿? | **1** | 盘上用例 PASS；夹具 FAIL |
| V6 | 复制+空壳 | 红 | **1** | `:353` sites=2 |
