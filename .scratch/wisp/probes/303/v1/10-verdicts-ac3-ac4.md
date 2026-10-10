# 303-v1 — 非实现者对 `AC#3`／`AC#4` 的一次裁（七问逐答；⛔ 翻框，判语交编排者执行）

起手锚现量＝`git rev-parse HEAD` → `34e1962bb166b4494670b77296263c79bd8671a3`。
起手闸门现量＝`tasklist`：`wisp.exe`=0／`balldebug.exe`=0／`mockllm.exe`=0（三行都逐字 `INFO: No tasks are running which match the specified criteria.`）；`msedgewebview2.exe`＝**24 枚**＝机主自有四棵应用树，⛔ 残留（先例 `A821`＋`probes/orch/webview-orphans/forensics-4-trees.txt`）。本腿⛔ 杀任何进程。
台面（本腿全部读数）＝**仓外干净 clone** `~/wisp-303-v1`（`git clone --no-hardlinks` 自母仓，`clone_rc=0`），四发整包与一枚定向发全部 `git checkout --detach` 于该 clone；⛔ 拿母仓工作树当"改前台面"，⛔ 在共享工作树做任何 checkout/reset/stash。变异取证（M1/M2/M3/M4＋pristine 对照）只落 `git archive HEAD | tar -x` 的**导出树**（`~/wisp-303-v1-mut-*`），母仓工作树**一枚文件⛔ 动过**。
三数口径（本文每处计数都带尺名）：尺 A＝`grep -c '^--- PASS'`／`'^--- FAIL'`／`'^--- SKIP'`＝**顶层**那一把（⛔ 与含子测试那把混用，两把差 110 枚）；尺 B＝词面尺，命中行逐字抄；尺 C＝格式名册，写明射程目录＋blob/工作树。

---

## ⓐ 这枚修复有没有把判据变松？——**没有（成立）**，附带一枚具名的网眼

**独立取证（⛔ 引用实现者自陈）**
1. 名册级：`git show --numstat --format= e3341368` → 十枚路径里唯一非 `.scratch` 者＝`cmd/wisp/panel_host_windows.go` **`27 0`**（＋27／−0）；`git show --name-only --format= e3341368 | grep -cE "thresholds|_test\.go|PLAN\.md|docs/specs|allowlist"` ＝ **0**（尺 B，件 `02-criteria-integrity.txt` gate-D）。⇒ 判据本体所在的两枚测试文件**根本没进那一笔**。
2. 比名册更强的一把：判据本体的**宿主文件自基线起零改动**——`git log --oneline cc315261..HEAD -- <p> | wc -l` 逐枚＝`cmd/wisp/panel_resident_windows_test.go` **0**／`cmd/wisp/panel_host_windows_test.go` **0**／`internal/observe/thresholds.go` **0**／`docs/PLAN.md` **0**／`docs/specs` **0**（件 `02-` gate-A）。⇒ ⛔ 止于"这一笔没碰"，而是"从绿到红到修好整段⛔ 碰过"。那句红句在 HEAD 对象层逐字仍在（`panel_resident_windows_test.go:226`，尺 B 抄进件 `02-`）。
3. 那 27 行的**形状**：HEAD 对象层里 `fmt.Sprintf` 那个原始串起于 `panel_host_windows.go:800`、止于 `:819`（件 `01-` 逐行），新增的 JS 九行落在 `:806-:813`＝**串内**；其余新增全是 `:782-:799` 的 `//` 文档行。⇒ **零枚新 Go 语句**、`dispatchRaw` 的控制流一字未动。`go build` 面：本腿 `go vet ./cmd/wisp/` rc=**0**、`go vet -tags winlive ./cmd/wisp/ ./internal/ball/` rc=**0**（件 `logs/gate-vet-plain.txt`／`logs/gate-vet-winlive.txt`，各含自落 `rc=` 行）。
4. SKIP／Fatalf 面：⛔ 测试件被改（见 1）；且**同一台面上改后 SKIP 名册是改前的真子集**（改前 3 枚、改后 2 枚，消失那枚＝`TestPanelHostLatencyPercentilesAC2`，负载敏感那族）——尺 A＋名册件 `logs/names-before-*.txt`／`logs/names-after-*.txt`。修复**没造出任何 `--- SKIP`**；`sh scripts/d22scan.sh` 本腿自跑 rc=**0**（件 `logs/gate-d22scan.txt`；它把任何 `^--- SKIP` 判红，那一格若被造假 SKIP 必炸）。
5. D32／SLO 面：仪器里的预算行逐字未动——`cmd/wisp/panel_host_windows_test.go:664-665`＝`if coldMs > 1500 {`／`t.Errorf("cold bring-up %.1f ms exceeds D32 panel cold budget 1500 ms", coldMs)`，`:1004-1005`＝`if coldP95 > 1500 {`／`cold P95 %.3f ms over %d runs exceeds the D32 panel cold budget 1500 ms …`；`internal/observe/thresholds.go` 里唯一的面板阈值字段逐字＝`memCapPanel int64 = 600 << 20`（尺 B，件 `02-` gate-B/gate-C；那文件⛔ 一枚延迟毫秒字段，本腿⛔ 碰）。
6. 实现者那把 `-run 'Transport|ShapeA3|LegacySubShape|Unforwarded|ForwardingHook'` **⛔ 是挑过的尺**：两档 transport 夹具里 `func Test` 共 **10** 枚，选择器命中 **10** 枚（`total=10 matched_by_selector=10 => coverage = ALL`，件 `04-yard-roster-and-library.txt`），且 35-r2 夹具逐字声明它读的是**出货常量本身**（`:25` `the PRODUCTION runtime string panelPostMessageForwardInit - read as a value from`）。⇒ 编辑后的钩子确实被整把尺量过。

**首要攻击点的实证（本腿另找的形状，⛔ 实现者给过的；件 `logs/mut-*.txt`）**
- **M1＝把覆盖层"折进门"那一支删掉**（`cw.postMessage` 一律交 native）。后果：真窗两枚**全绿**——`--- PASS: TestAC14AwaitedBindingReplyReachesThePage (1.18s)` 拿到逐字 `page's own words: "REPLIED,REPLIED,REPLIED"`，`--- PASS: TestPanelHostRealWindowHopAndLifecycle (1.43s)` 拿到 `cold bring-up measured on this box: 1116.089 ms`（件 `logs/mut-M1-real.txt`，rc=0）＝**一枚足以把票 303 的红换成假绿的形状**（真实页面的 `chrome.webview.postMessage` 那条边被杀，等于退回票 35:52-59 的"死在无绑槽"）。**今天抓得到**：35-r2 行为夹具在同一棵树上炸开 5 枚（尺 A：`--- FAIL` 5／`--- PASS` 6，rc=1，件 `logs/mut-M1-yard.txt`）——具名＝`TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3`／`TestShapeA3SecondPagePostStillDelivers`／`TestShapeA3ForwardingHookIsIdempotentInOneDocument`／`TestForwardingHookMustNotLoseTheNativeExitReceiver`／`TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit`。⚠ **这五枚全是 35-r2 那一档**：35-r1 的能力钉 `TestPagePostMessageEnvelopeReachesDispatchRawViaTransport` 与名册钉 `TestTransportDoorBindingMatchesRoster253r1` 在 M1 下**照旧绿** ⇒ 挡住这一形的**只有 35-r2 那把行为尺**，它一旦被人挪走/放宽，M1 就是一条无人看守的假绿路（这条本身是网眼的另一半，具名交编排者）。
- **M2＝把 `inside` 的置位点挪走**（安装时置真、永不复位）。后果与 M1 同路由，**同一把尺同样 5 枚炸**（件 `logs/mut-M2-yard.txt`：`--- FAIL` 5／`--- PASS` 6，rc=1，那五枚逐字同名）。
- **M3＝让夹具⛔ 再模拟直接 `wispDispatch` 入口**（页面改经 `chrome.webview.postMessage` 送信封）＋撤掉修复。后果：**假绿造不出来**——nail1 逐字仍 `no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all)`（`--- FAIL (20.02s)`），且 `TestPanelHostRealWindowHopAndLifecycle` 红在哨兵 `cold bring-up measured on this box: -1.000 ms`（件 `logs/mut-M3-real.txt`，rc=1）。⇒ 换夹具这条路⛔ 通，因为**产码里的探针自己就是直接入口**：`firstRoundTrip` 在 `panel_host_windows.go:879` 绑一枚一次性门 `wispProbeRT`，探针页 `:892-893` 逐字 `<script>if(window.wispProbeRT)window.wispProbeRT();</script>`——测试件改⛔ 掉它。⚠ 顺带一枚**具名新知（M3 单发，`-count=1`，机制⛔ 证）**：该形下真窗连 postMessage-origin 那条边也"nothing at all"，而 headless 夹具在同形下报 `native exit 1／door 1`＝**夹具与真窗的一枚差**，交编排者归口（⛔ 本票射程，⛔ 影响本问判语）。
- **M4＝⛔ 动 JS，改在 Go 侧把双包帧拆开**（`dispatchRaw` 认 `method==wispDispatch` 就重派内层）＋撤掉修复。两件事都拿到了：① **10 枚 transport 钉全绿**（尺 A：`--- PASS` 10／`--- FAIL` 0，rc=0，件 `logs/mut-M4-yard.txt`）＝**今天那张网⛔ 看得见门的 Go 那一侧**，这是一枚具名网眼，建议进台账（一格："transport 夹具只量 JS 覆盖层，⛔ 量 router 对'方法名＝绑定名'那种帧的拒绝"）；② 但它照样**造不出绿**：`--- FAIL: TestAC14AwaitedBindingReplyReachesThePage` 红在 `page's own words: "TIMEOUT-2S,TIMEOUT-2S,TIMEOUT-2S" … Go did receive the calls (3 reached the handler)`，lifecycle 仍 `-1.000`（件 `logs/mut-M4-real.txt`）。⇒ 真正兜底的是**判据的形状**（nail1 断的是页面自己 await 到的回执，⛔ 是"Go 收到了没有"），这比任何钉子都硬。
- pristine 对照（HEAD 导出树、零变异）：yard rc=0／11 枚 PASS，real rc=0（`cold 1341.042 ms`＋`REPLIED,REPLIED,REPLIED`）＝上面每一枚红都⛔ 是台面自带的。件 `logs/mut-pristine-*.txt`。

⇒ **ⓐ 判语＝成立**：修复⛔ 放宽判据、⛔ 造 SKIP、⛔ 碰 D32/thresholds；且"能换成假绿"的形（M1/M2）今天由 35-r2 那 5 枚钉抓得住，M3/M4 那两形即便落上去也⛔ 能翻绿。附一条欠账（网眼 M4）与一条台面差（M3），都具名交给编排者，⛔ 由我修。

---

## ⓑ ⛔ 扩权限？——**没有扩（成立）**，并附一枚⛔ 本笔引入的既有名册洞

我读的链，逐枚具名（全部 HEAD 对象层或工作树**只读**）：
- `cmd/wisp/panel_host_windows.go:828-837` `installPanelTransport`＝唯一绑门处：`w.Bind(panelDispatchBinding, func(raw string) string { reply, _ := m.dispatchRaw(ctx, raw); return reply })` 然后 `w.Init(panelPostMessageForwardInit)`。绑定名常量 `:80 panelDispatchBinding = "wispDispatch"`。修复**没加 Bind、没加 Init、没新名**（diff 全在串内＋注释，见 ⓐ3）。
- 同文件 `:841-849` `dispatchRaw`→`m.disp.Handle(ctx, raw)`（本笔零改动）。
- `internal/panel/composer_dispatch.go:155-167` `Handle`：`ParseComposerRequest` 是**唯一那道门**（逐字注释 `// One gate, the existing one. Everything parse refuses never gets routed.`）；`:177-212` 派发表的每个 case 都是 bridge.go 的导出常量；`:209-211 default: return "", d.rosterMismatch(req)`；`:217-224` 拒绝并记审计（`ErrRosterMismatch`）。
- `internal/panel/bridge.go:126-141`＋`:146-152`：闭集**六枚**方法名（`MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend, MethodConfigGet, MethodConfigSet`），外加 `:136-137` 的 source 校验。⇒ `"wispDispatch"` 当**方法**时必被拒（＝回归期间双包帧撞的正是这堵墙）。
- `internal/panel/composer_handlers.go:111-146`：`Modes == nil` → `ErrNoModeWriter`（`:117-121`）；当前档不可读 → 逐字 `无法判断 %q 比现在更宽还是更严，按 fail-closed 拒绝`（`:127-131`）；**放宽且 `h.Confirm == nil`** → `ErrNoL2Confirm`（`:132-137`，`var ErrNoL2Confirm = … "panel: 确认腿未接入，档位不得变宽"` 在 `:55-57`）。`internal/perm/store.go:194` 同形逐字 `切到 auto_approve 需要一次 L2 强确认，但本机没有接入确认通道（fail-closed 未切）`。⇒ 编排者那句"门全在 Go 侧"我**读过后独立确认**，⛔ 照抄。
- 库面（只读模块对象层）：`github.com/jchv/go-webview2@…56598839c808/pkg/edge/chromium.go:112` 逐字 `e.Init("window.external={invoke:s=>window.chrome.webview.postMessage(s)}")`——`window.external.invoke` 是**库自己在控件创建时装上的**，修复前就可达（它是**每一枚**绑定门唯一的出页面路径）；新包装把入参**一字不动**转发并原样返回 `nativeInvoke(s)` ⇒ **可达目标集⛔ 变**。
- 新包装实际改的是**路由**：invoke 内的帧从"折进门"改成"交回 native 出口"。那一改对面是页侧**只⛔ 不收**——若页面自己直调 `window.external.invoke(<composer 信封>)`（⛔ 走绑定桩），修前会被折进门、修后逐字节交给 native 再由 msgcb 拒收。⇒ 这是一枚**变窄**，⛔ 变宽；且无产品路径用它（页面 `panel.ts` 用 `chrome.webview.postMessage`，绑定桩用 RPC 帧），整把 10 枚钉在修后全绿（`logs/mut-pristine-yard.txt` rc=0）＝没有任何钉在保护那条旧路由。
- 名册⛔ 动：`cmd/wisp/panel_dispatch_binding_roster_253r1_windows_test.go:29-45` 那四条（常量已声明／`installPanelTransport` 全 package 恰好一枚／其中的 `Bind` 双向对上名册／交给 `w.Init` 的串由同一常量 `fmt.Sprintf` 拼出）在修后仍绿（它在本票整包 269/270 PASS 名册里；名册＝`logs/names-after-*.txt`）。
- ⚠ **具名一枚既有洞（⛔ e3341368 引入，⛔ 本腿补）**：`wispProbeRT`（产码 `panel_host_windows.go:879` 绑的探针门）绑在 `installPanelTransport` **之外**，而 253r1 名册尺的四条射程写死在 `installPanelTransport` 内（`:29-45`）⇒ 那枚门⛔ 受名册钉；它的 Go 侧只 `close(done)` 并返回 `"ok"`，无副作用、无权限面，但从票 33 起它就在那儿。建议登记，⛔ 我动。

⇒ **ⓑ 判语＝成立：那 +27 行⛔ 放进任何一条原本被 Go 侧挡住的信，⛔ 动绑定名／名册／C17 六法闭集／L2 腿。**

---

## ⓒ 按更正后的目标，`AC#3` 算不算闭合？——**ⓣ 判"本格凭据句有缺陷"（⛔ 判"修复不成立"）；ⓤ 更正形状下闭合，附两枚条件**

- **ⓣ 这一支我判"判据本身有缺陷"**。本格原文把凭据写成"修后 `TestAC13`／两枚 `TestAC14*` 全绿"，而本票自己的**边界段**逐字把 `TestPanelHostRealWindowHopAndLifecycle` 的预算那一支、以及"另外三枚的形状"划在⛔ 本票射程；也就是说那句凭据**与本票的射程声明自相矛盾**，它要求的是一枚从未成立的事。我这腿**独立复现**了"从未成立"：在 `f718e9b6`（＝回归笔 `fb2fb802` 的**父发**，回归根本还⛔ 存在）同一枚 clone 台面上定向三枚 ⇒ `--- PASS: TestAC14AwaitedBindingReplyReachesThePage (1.43s)`（逐字 `REPLIED,REPLIED,REPLIED`）而 `--- FAIL: TestAC14GoSideEvalPushReachesThePage (0.75s)`，红句逐字 `AC#14 nail 2 (Eval push hop), page's own words: title=""` ＋ `the page reports its title as "", want "PUSHED-33R5-OK"`（件 `logs/nails-at-f718e9b6.txt`，rc=1）。⇒ 缺陷在**票面那句凭据**，⛔ 在实现。**翻框⛔ 由我做**；建议的处置＝按 ⓤ 闭合本格，并把"原凭据句作废"这件事留在编排者已写的追加节里（票面 `:99-100` 已具名作废，本腿确认那一句有盘上凭据）。
- **ⓤ 更正形状（nail1 绿 ＋ 两枚残余另立票 305）＝闭合**。承重读数（都出自本腿自己的手）：改前同一台面两发**逐字红** `no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all)`（`logs/full-before-1.txt:854`／`full-before-2.txt:854`），改后两发**逐字绿** `AC#14 nail 1 (reply hop), page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)`（`logs/full-after-1.txt:866`／`full-after-2.txt` 同枚）；四条硬约束逐条见 ⓐ；四条门禁 rc 见 ⓓ。
- **两枚条件**（本格闭合挂着它们）：① `AC#4` 那把主尺必须"各两发＋同台面"——已由本腿补齐（ⓓ，新增红 **0** 枚）；② 两枚残余必须由票 305 承接（ⓖ，我判"另立"正确）。任一枚塌了，本格就退回〔附条件⛔ 足〕。

---

## ⓓ `AC#4` 欠的那一发——**补齐完成，本格形状现已满足（成立）**

台面（★本轮新立的台面铁律，本腿逐字遵守）＝**仓外干净 clone** `~/wisp-303-v1`，四发**全在同一枚 clone 上**、彼此只差一次 `git checkout --detach`；⛔ 母仓工作树当改前台面，⛔ 在共享工作树里做任何 checkout/reset/stash（AGENTS §1.4）。每发台件与逐发身份：`logs/checkout-before-1.txt`（`detached=cf46c24a`，`panel_host_windows.go` 行数＝969→见件）、`logs/checkout-after-1.txt`（`detached=807497c1`）等四枚，各带 `checkout_rc=0`。

尺 A（顶层 `grep -c '^--- PASS'`／`'^--- FAIL'`／`'^--- SKIP'`，`-v` 输出，⛔ 与含子测试那把混用）：

| 发 | detached | RUN(尺 A 合计) | PASS | FAIL | SKIP | rc |
|---|---|---|---|---|---|---|
| before-1 | `cf46c24a` | 288 | 267 | 18 | 3 | `rc=1` |
| before-2 | `cf46c24a` | 288 | 267 | 18 | 3 | `rc=1` |
| after-1 | `807497c1` | 288 | 269 | 17 | 2 | `rc=1` |
| after-2 | `807497c1` | 288 | 270 | 16 | 2 | `rc=1` |

- 两枚改前发**红名册逐字相同**（`diff logs/names-before-1.txt logs/names-before-2.txt` ＝ 空），红名交集＝**18 枚**（`logs/before-red-intersect.txt`）；两枚改后发红名交集＝**16 枚**（`logs/after-red-intersect.txt`）。
- ★**新增红＝交集作差 0 枚**（`comm -13 before-red-intersect after-red-intersect` ＝ 空，件 `logs/new-red.txt`，行数 0）。⇒ 票面 `AC#4` 要的"整包改前／改后各两发取**红名交集**＝新增红 0 枚"**这个形状今天第一次真被满足**（编排者手上是各一发＋一枚定向隔离发；他据那一发得到的"＋1 枚"是 `TestAC4FocusReturnToPriorWindowGap33r5`，而它在我的**两枚改前发里都红**＝台面/桌面态，⛔ 修复）。
- 被修好＝2 枚（`logs/fixed-red.txt`）：`TestAC14AwaitedBindingReplyReachesThePage` ＋ `TestPanelHostRealWindowHopAndLifecycle`（后者在本台面改后 `cold bring-up measured on this box: 1002.355 ms`／`full-after-2` 同枚正向读数，进 1500 预算内＝票面边界段只授权的那一枚迁移，⛔ 更多）。
- 具名单发红（⛔ 进任何交集）：`TestTicket223HandEditedFsLooseningCostsAnL2Card` 只在 after-1 红（4 发里 1 发），同发逐字带 `goroutine outside the D38 roster (leak symptom) goroutine=mockllm-stdout-reader`＝票 303 边界段明写⛔ 本票射程那一族。
- **台面档差具名（⛔ 让它混进"新增红"）**：① clone 无 `frontend/dist` ⇒ `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` 在**四发全部**走具名 `--- SKIP`，逐字 `AC#13 has no subject in this tree: the embed resolves no entry, so there is no page content that could be covered.`（母仓台面它是 `--- FAIL`）；② 同一台面上另有自报族 `no native DLLs in ..\..\third_party\sherpa-onnx - run scripts/fetch-deps.ps1 first`（编排者那对发里 14 枚；两发同台面⇒作差抵消）；③ `0xc0000135` 尺＝`grep -c` 在我的四发里均 **0** 命中（⇒ 用例真跑了；PATH 用 shell 形 `/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:…/build:$PATH`，逐字写在 `ac4-second-pair.sh`）。
- 门禁 rc（本腿**自己**跑的，⛔ 用实现者的读数；每把一件自落 `rc=`）：`logs/gate-vet-plain.txt` rc=**0**、`logs/gate-vet-winlive.txt` rc=**0**、`logs/gate-d22scan.txt` rc=**0**、尺 C（格式名册，射程＝`cmd/wisp` 目录）＝工作树那把 `logs/gate-gofmt-worktree.txt` 列 **1 枚** `cmd\wisp\models.go`（既有加严残留，**⛔ 本腿文件、⛔ 顺手修**）、HEAD blob 那把 `logs/gate-gofmt-headblob.txt` 列 **0 枚**（导出树 `~/wisp-303-v1-blob`）；全仓那把仅作上下文，件 `logs/gate-gofmt-repocontext.txt` 枚数＝**38**（射程＝`.` 母仓**工作树**，脏，别的腿在飞；⚠ 与派单里那句"全仓 5 枚既有残留"**对不上**，差别＝工作树 vs blob／CRLF ⇒ **报差不自调**，⛔ 本腿射程、⛔ 修）。
- 越界自查：每笔 `git show --name-only --format=` 过滤 `^.scratch/wisp/probes/303/v1/` 后计数＝**0**（逐笔在回报里给号）。

---

## ⓔ 甲那一形欠的门外读数——**取到了，但⛔ 是实现者设想的那一枚；甲就此判死（成立，附一枚残口交票 305）**

- 先说取⛔ 到的那枚，⛔ 怪台面：**"plain `Eval` 设 `document.title` 再由 Go 侧 `Eval` 同步取回"在这台机器上根本存在⛔ 不了**。库的对象层逐字＝`func (e *Chromium) Eval(script string) {`（`chromium.go:138`）且 `ExecuteScript.Call(…, 0)`（`:143-147`，handler 位＝0）＝**求值结果被库丢掉**；夹具自己也写了这一点（`panel_resident_windows_test.go:160-166` 逐字 `Eval returns no value, and that is on purpose in this file's argument`）。`ICoreWebView2::add_DocumentTitleChanged` 只在 vtable 里（`corewebview2.go:120`），本仓 `grep -rn "DocumentTitleChanged" cmd/wisp internal/panel` 命中 **0** 枚。⇒ 要拿那枚读数**必须新写宿主代码**＝⛔ 本腿边界（也⛔ 票 303 名册）。
- 取到的那枚（同一批真窗夹具，⛔ 额外一发）：`firstRoundTrip` 走的⛔ 是 composer 那道门——`panel_host_windows.go:879` 它自己 `w.Bind("wispProbeRT", func() string { close(done); return "ok" })`，Go 侧只关 channel，⛔ 进 `dispatchRaw`／⛔ 进 `ComposerDispatch.Handle`／⛔ 依赖页面自报；而它的 `done` 只能由探针文档里那句 `<script>if(window.wispProbeRT)window.wispProbeRT();</script>`（`:892-893`）关掉。**所以 `cold bring-up` 是个正数＝"那份文档的 JS 确实跑了、且一枚直调的绑定门确实回到了 Go"，全程与路由无关。**本机逐字（尺＝日志行号）：
  - 改前台面（`detached=cf46c24a`）：`panel_host_windows_test.go:660: cold bring-up measured on this box: -1.000 ms`（`full-before-1.txt:713`、`full-before-2.txt:713` 两发同形）。
  - 改后台面（`detached=807497c1`）：`cold bring-up measured on this box: 1002.355 ms (HEAD 807497c1 at read time, 2026-10-10T20:16:06+08:00)`（`full-after-1.txt:717`）／`980.666 ms (HEAD 807497c1 at read time, 2026-10-10T20:23:37+08:00)`（`full-after-2.txt:713`）；pristine 导出树那发另取一枚 `1341.042 ms`（`logs/mut-pristine-real.txt`）；`f718e9b6` 定向发里宿主自报 `cold 1186.3 ms`／`511.9 ms`（`logs/nails-at-f718e9b6.txt`）。本机四枚正数全部 <1500 且 <2000。
⇒ **甲（页面 JS 压根⛔ 跑）由"因果对否证"升为"门外读数直接否证"**，本机、真窗、三枚正数、两枚 `--- PASS`。
- 残口具名（⛔ 本票）：它证的是**探针文档**的 JS 在跑，⛔ 证**内嵌入口**那份文档的 JS 在跑——那正是票 305 症状①（"冷启后最后盖着的文档是探针 stub"）。⇒ `AC#2` 甲那一格在"传输边"射程内可判死；入口文档那一半随票 305 走。

---

## ⓕ CI 新红那一枚（`TestPanelHostLatencyPercentilesAC2` SKIP→FAIL）的归口——**归票 302 的归口面；⛔ 票 33 预算面、⛔ 票 305、⛔ 另立一票；⛔ 动那枚 1500**

理由（每条带件）：
1. **产品侧那枚数并⛔ 超预算**：同一测量在本机是 `1002.355 ms`／`1341.042 ms`（ⓔ 那批件，预算 1500 之内）；超预算的只有**托管 runner** 那一枚 `2747.460 ms`（件 `orch/r9-ci-after-red-roster.txt`）。台账早就把这档读数判过：`docs/reports/pending-and-issues.md:10718` 逐字 `那台机器自己的时序，⛔ 不许当我们的性能读数`。⇒ 把它登记成"D32 预算破了"＝用 runner 的机器去改冻结面，⛔ 允许。
2. **那枚钉子的性质是档位／放置**：它读的是**同进程里另一枚真窗用例**攒下的样本，样本缺失时具名 `--- SKIP`（`panel_host_windows_test.go:985` 逐字 `no cold/hot sample recorded in this process: the lifecycle test did not run in this binary`），样本一到才开始求值（`:1004-1005`）。⇒ 它随"lifecycle 有没有跑、在哪个档跑"变色，正是票 302 的射程（"这些页面回话判据长在哪个档"）。
3. **⛔ 票 305**：与文档交接⛔ 关（它一个字节都⛔ 读文档，只读 `panelHostLatency` 那个聚合）。
4. **⛔ 另立一票**：形状⛔ 新的——它是 lifecycle 那枚既有测量的**第二把嗓子**，同一发里两枚红句同数（`2747.460`）。
5. ⚠ 一条必须具名带走的**待人拍板**线索：`panel_host_windows_test.go:1003` 逐字 `AC#2 P11 line (cold > 2000 ms) - max cold observed over these runs: %.3f ms`——**P11 的触发条件（WebView2 冷拉起 >2s → 重评 L2 卡是否回原生，AGENTS §2 第 6 条、S0）今天在本机⛔ 触发**（1002/1341 < 2000），只在 runner 上触发过（2747 > 2000）。而那条待定的措辞指的是**机主那台机器**。⇒ ⛔ 我裁 P11，⛔ 任何腿拿 runner 的 2747 去立案 P11；只需在票 302 归口格里写清"该枚红＝runner 档的既有预算面二手读数，本机未触发 P11"。

⇒ **ⓕ 判语＝归票 302 的归口面（登记＋档位裁决），附一条具名"本机未触发 P11"的口径；SLO／1500／`thresholds.go` 一字节⛔ 动（本腿零产码，ⓐ5 已证）。**

---

## ⓖ 票 305 的处置对不对？——**"另立票 305"正确（成立）；那两枚红是另一枚形状，且⛔ 一条因果**

- 独立否证"⛔ 是 fb2fb802 造成的"＝本腿自己跑的 ⓔ/ⓣ 那一发：`f718e9b6`（回归笔的父）同台面定向三枚 ⇒ nail1 **PASS**、nail2 **红句逐字同一句**（`logs/nails-at-f718e9b6.txt`）。一枚在回归之前就逐字相同的红，⛔ 可能是回归造成的，也⛔ 可能靠修那一笔消掉。
- 修复的射程⛔ 碰文档交接：`e3341368` 名册只有 `panel_host_windows.go` `+27/−0` 且全落在那枚 JS 串＋注释里（ⓐ1/ⓐ3）；`bringUp`／`serveEntry`／`firstRoundTrip` 的次序（`:448` 那行调用序）一字未动。⇒ 要把 `AC13`／nail2 收进本票＝**必须动第二枚未归因的产码**，正撞票面规矩"⛔ 顺手扩范围／未定义即停"。
- 两枚症状⛔ 同一条因果，票 305 自己已按 `AC#0` 把它们分开（票面 `:4` 逐字 `两枚症状、⛔ 一条因果`），且它对 `70b00885` 的口径是对的——写成**〔仅腿报〕，我⛔ 复跑过归因**（票 305 `:20`）。⇒ 我这腿⛔ 替它归因；只补一条读数：`AC13` 在 clone 台面**四发全具名 SKIP**（无 `frontend/dist`），它的颜色只有母仓／CI 台面能判——票 305 `判语②` 已经这样写，本腿独立复认。
- 一处该由编排者并回的口径：票 305 症状①里"`TestPanelHostRealWindowHopAndLifecycle` 颜色随台面"这件事，本腿拿到了第三枚凭据——同一枚 `807497c1` 在**仓外干净 clone** 上是 `--- PASS`＋`cold 1002.355 ms`（`logs/full-after-1.txt:717`），而腿在母仓那发仍是 `-1.000`。母仓那一色本腿⛔ 复测（⛔ 在共享工作树 checkout），保持〔仅腿报〕。
⇒ **ⓖ 判语＝"另立票 305"⛔ 是顺手扩本票，反而是唯一不破规矩的处置；票 305 的现量段建议并上本腿这三枚逐字（`nails-at-f718e9b6`／`full-after-1:717`／`mut-M3-real` 的夹具-真窗差）。**

---

## 本腿欠账（具名，⛔ 我填）
1. M3 那枚"夹具与真窗对 postMessage-origin 那条边的预言相反"＝单发（`-count=1`），**机制未证**；要定它需要 `-count=3` 一发＋读 msgcb 那一侧，⛔ 我射程（我⛔ 动产码/测试件）。件 `logs/mut-M3-real.txt`／`logs/mut-M3-yard.txt`。
2. M4 具名的**网眼**（transport 夹具⛔ 量 router 对"方法名＝绑定名"的拒绝）＝本腿只证了"10 枚全绿"这一半，⛔ 设计补钉；归口建议＝票 302 的归口格或一张仪器票。
3. 实现者设想的门外 `Eval` 同步回读**取⛔ 到**（库丢结果，见 ⓔ），要它就得新写宿主代码＝⛔ 本腿边界。
4. 母仓台面的 lifecycle 色本腿⛔ 复测（台面铁律）；`ImageVersion` 那类跨载体读数⛔ 再取（已由编排者在 `r9` §4 自纠）。

## 自陈三句
⛔ 翻框（票 303／305／302 的任何 `- [ ]` 我一枚没碰，判语只写在本件里）。⛔ 产码（母仓工作树与 clone／导出树之外我⛔ 动过任何 `.go`；变异只落 `~/wisp-303-v1-mut-*` 导出树，且⛔ 入库）。⛔ push（本腿只 commit，pathspec 逐笔显式）。
