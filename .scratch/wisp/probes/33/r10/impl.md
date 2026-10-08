# 33-r10 交付正文（票 33 AC#13 的 ②＋③，与框前提的过期复量）

腿＝写码腿 `33-r10`。起手 HEAD `128bf600fa0890c02887e1d005dc7d6b1bd87a43`；交件时 HEAD 已被别的程推到
`50341a861493d75cf54e01ed9ad86a63d4268b61`（本腿两枚提交 `44d178c9`／`70b00885` 都在其祖先链里，
`git diff --numstat 128bf600..HEAD` 现量本腿那两枚件＝`30 11`＋`293 0`，⛔ 无人改回）。

---

## §0 锚 + (a)(b)(c) 三格现量

起手尺（逐条现跑，全文与逐字原文在 `00-anchor.md`，已单独 commit）：

- `date` = `Thu Oct  8 10:06:21 CST 2026`（写本件时 `10:20:10`）
- `git log -1 --format=%H` = `128bf600fa0890c02887e1d005dc7d6b1bd87a43`
- `git status --porcelain | wc -l` = **753**（⛔ 不是派单说的"30 枚上下"；见 §7 冲突格）
- 票 33 第 39 行整行 = **1413 字符**，`awk 'NR==39'` 全文读入并逐字粘进 `00-anchor.md`，
  且用 `diff` 证过粘件与盘上原文**逐字节相同**（`QUOTE_IDENTICAL rc=0`；中途我手抄走样过一次，已用原文重建后再 diff）
- 框枚数尺（追加前）＝ **13 未勾／1 已勾**

尺＝`grep -n 'serveEntry\|firstRoundTrip\|SetHtml\|wispProbeRT\|AddWebResourceRequested' cmd/wisp/panel_host_windows.go`
（当时 861 行）。逐行读数已随格次落 `00-anchor.md`。

### (a) 今天谁先谁后？

**探测先、供页后。** `bringUp` 在 **`:425`** 调 `firstRoundTripLocked(ctx, t0)`（它内部 **`:757-758`** 发那一发探测页 `SetHtml`），
随后 **`:427`** 才调 `serveEntry()`（**`:467`** `w.SetHtml(string(data))` 把 embed 入口字节推进控件）。
⇒ 编排者 10-08 那句"`:425` 那发探测先于 `:427` 那发供页"**成立**；
票面 AC#13 的 `:221`/`:250`/`:227`/`:372-375` 四枚行号在本 HEAD 的 861 行文件里**全部对不上**（那些位置不是那些调用）。

### (b) 后一发是否真的把前一发的内容覆盖掉？

**覆盖还在，但方向反了，而且反到了"对用户无害"的那一面。** 后一发＝`:467` 的入口字节，被覆盖者＝`:757-758` 的探测页
⇒ 用户**最终看到的文档是 embed 入口页**；探测页只活在 `:425`→`:427` 那段泵消息的时间里，最坏＝启动瞬间一次闪动。
失败支也不落在探测页上：`serveEntry` 报错 → `:427-429` 调 `serveNotBuiltNoticeLocked()`（`:448` 那一发），
最终文档是宿主自己写的"panel assets unavailable"明示页。
⇒ **今天没有任何一条路径让最终文档＝探测页**（此句射程＝产码这两支，不含"用户真的看见了面板"，见 §5）。

### (c) 框那条"最终显示的是探测页"是否已经不成立？

**作为今天的产码事实＝已不成立；框的 ①（次序重排）早已落地。** 换向那一刀的出处现量＝`13acad46`
（`git log -S'hand the page over LAST' -- cmd/wisp/panel_host_windows.go`，33-r5，2026-10-01），
产码注释 `:419-424`／`:729-731` 自己具名写着 "AC#13's order, and it is the whole fix"。
⇒ **处置＝不造不存在的 bug**：按派单第 2 节末句，把新事实写成票 33 末尾追加一节，具名"编排者 10-01 那句过期"，
一枚框字未动（追加后尺＝**13 未勾／1 已勾**，与本节前尺同数，见 §4）。

**② 那一问今天仍然有效，且有实测的洞**（不是"已经有仪器了"）：唯一问最终文档内容的用例
`cmd/wisp/panel_resident_windows_test.go:305 TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`
买的是**真窗**（`startPanelForTest`），且 `:317` 在入口解不出时具名 skip ⇒ **新鲜检出／CI 里这一格零仪器**；
而它在本树今天**不会 skip**（embed 解得出 1044 字节入口）⇒ 跑它＝开真窗＝派单禁事，本腿⛔ 没跑它。
⇒ 本腿补的是**开窗无关**那一面。

---

## §1 动了什么（逐件＋行号；产码 diff 逐字）

### 1.1 产码：`cmd/wisp/panel_host_windows.go`（唯一动过的产码件）

尺＝`git diff --numstat 128bf600..HEAD -- cmd/wisp/panel_host_windows.go` = **`30 11`**（改＝是；不是零改动，
故此处给 full diff，⛔ 不适用"零产码改动的空读数自证"）。
改动性质＝**把 `bringUp` 尾段那两发原样搬进一枚命名函数**：语句一条不加不减、次序不变、跨调用不持锁、
`rtMs` 仍由 `bringUp` 写进 `m.lastColdMs`。没加 `AddWebResourceRequestedFilter`——
派单要求"要加过滤器就必须在件里具名写为什么次序重排不够"，**本腿不需要**：次序今天已经是对的（§0(b)），
缺口只在"没有开窗无关的仪器钉着最终文档"，那是台件面的事，过滤器解决不了它。

逐字 diff（同一份内容也在 `01-product-diff.log`，末尾 `rc=0`）：

```diff
diff --git a/cmd/wisp/panel_host_windows.go b/cmd/wisp/panel_host_windows.go
index 26b5de83..1f9060df 100644
--- a/cmd/wisp/panel_host_windows.go
+++ b/cmd/wisp/panel_host_windows.go
@@ -416,17 +416,9 @@ func (m *PanelManager) bringUp(ctx context.Context) error {
 	m.created = true
 	m.mu.Unlock()
 
-	// AC#13's order, and it is the whole fix: prove the message channel first,
-	// hand the page over LAST. firstRoundTripLocked shows a document of its own,
-	// so running it after serveEntry meant every cold start finished on the probe
-	// page instead of the panel (ticket 33 AC#13, 33-v1 §AC#3 "供给那半"). The
-	// probe stays - it is where cold "usable" is decided - it just no longer gets
-	// the last word about what the user sees.
-	rtMs := m.firstRoundTripLocked(ctx, t0)
-
-	if err := m.serveEntry(); err != nil {
-		m.serveNotBuiltNoticeLocked()
-	}
+	// AC#13's order lives in coldStartPageHandover (33-r10 moved it out verbatim:
+	// same statements, same order, nothing added or dropped). Read the reason there.
+	rtMs := m.coldStartPageHandover(ctx, t0)
 
 	m.mu.Lock()
 	m.lastColdMs = rtMs
@@ -434,6 +426,33 @@ func (m *PanelManager) bringUp(ctx context.Context) error {
 	return nil
 }
 
+// coldStartPageHandover is AC#13's order, and it is the whole fix: prove the
+// message channel first, hand the page over LAST. firstRoundTripLocked shows a
+// document of its own, so running it after serveEntry meant every cold start
+// finished on the probe page instead of the panel (ticket 33 AC#13, 33-v1 §AC#3
+// "供给那半"). The probe stays - it is where cold "usable" is decided - it just no
+// longer gets the last word about what the user sees.
+//
+// Why this is a function and not the four lines it was inside bringUp (leg 33-r10,
+// ticket 33 AC#13 item 2): bringUp cannot run headlessly - it calls
+// webview2.NewWithOptions, which needs a live WebView2 control - so the shipped
+// handover order had no ruler that could run without opening a window. The only
+// case that ever asked the question, TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe
+// (panel_resident_windows_test.go), buys a real browser and in a fresh checkout
+// skips for lack of a bundle. Splitting the tail out is the same move
+// installPanelTransport already made for the same reason (ticket 35 AC#6, see the
+// comment at that call): it moves zero behaviour, holds no lock across the call,
+// and gives AC#13's content assertion a headless way to reach the shipping order.
+// That reach is what TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe needs.
+func (m *PanelManager) coldStartPageHandover(ctx context.Context, t0 time.Time) float64 {
+	rtMs := m.firstRoundTripLocked(ctx, t0)
+
+	if err := m.serveEntry(); err != nil {
+		m.serveNotBuiltNoticeLocked()
+	}
+	return rtMs
+}
+
 // serveNotBuiltNoticeLocked replaces the page with an explicit offline notice when
 // the embed carries no bundle. It is a SetHtml of a document this host writes, so
 // the user is never left looking at the round-trip probe page, and it still opens
rc=0
```

搬动后的锚位（现量）：`bringUp` 那一跳＝`:419-421`；新函数＝`:429`（注释首行）／`:447`（函数签名）／
`:448` 探测／`:450` 供页；`serveEntry` 本体随文件下移，`:472` 仍是它的签名。

### 1.2 台件：新建 `cmd/wisp/panel_pageover_33r10_windows_test.go`（293 行，`//go:build windows`，无 `winlive`）

- `:41-58` 文件头写明"问的是最终文档的能力，⛔ 不问调用枚数、⛔ 不用行号、不开真窗"
- `:52-146` `docSink33r10`＝`webview2.WebView` 全 12 法（`Run/Terminate/Dispatch/Destroy/Window/SetTitle/SetSize/Navigate/SetHtml/Init/Eval/Bind`）
  的记录型实现；`SetHtml` 把每一份文档按序存进 `c.docs`（`:72`）
- `:152` `attachSink33r10`＝照 `bringUp` 成功后的写法把控件挂到 manager 上（`m.w`／`m.created`），⛔ 不创建任何东西
- `:169` `hasLookableContent`＝剥掉 `<script>` 后正文里是否还有元素或文字
- `:184` `entryFromEmbed33r10`＝走**产码同一条接缝** `panel.BuiltinAssets()` + `Resolve(panel.EntryFile)`，取入口字节与其中的 `id=`
- `:217` `captureProbeDoc33r10`＝**当场把探测页跑出来取**（对一次性 sink 调 `firstRoundTripLocked`），⛔ 不是抄进台件的字符串
- `:241` `TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe`＝三问（见 §2）
- `:287` `head33r10`＝红句里只截 120  rune

⛔ 没动 `panel_transport_35r1_test.go`／`panel_transport_35r2_test.go`（35 族）；⛔ 没动任何既有测试件。

---

## §2 ②那枚"会响的断言"怎么问能力、不问调用

`coldStartPageHandover`（产码那尾段本身）对着记录型控件跑一遍，然后**只读最后一份文档**问三句：

1. **身份**：最后一份文档与"探测步骤当场写出的那一份"逐字节相同吗？相同即红。
   探测页不是字符串常量而是运行时从 `firstRoundTripLocked` 那里取来的 ⇒ 谁改探测页的 markup，取到的就是改后的那一份，
   不会因为"框里的字面量过期"而静默放行。
2. **可看内容**：把 `<script>` 剥掉之后，`<body>` 里还剩不剩元素或文字？不剩即红。
   探测页的正文就是一段自调 `<script>`，剥完为空；明示页有 `<p>`；入口页有 `<div id="root">`。
3. **入口真内容**（入口解得出时）：最后一份文档含不含 `Resolve(panel.EntryFile)` 那 `%d` 字节的本体、
   含不含从同一串字节里解析出的元素 `id`（本树现量 1 枚 `root`）？缺即红。

⛔ 全程没有 `len(docs)` 之类的调用计数参与判定，⛔ 没有 `SetHtml 被调用过 N 次` 的断言，⛔ 没有按行号开刀。
（`docs` 切片只是"把控件被交给它的文档存下来"的容器；判定只取最后一份的内容。）

---

## §3 ③反控 + 三发正控（颜色与红句逐字；日志在 `mut/`）

载具＝`python .scratch/wisp/probes/33/r10/mut/run-mutations.py`：
⛔ 一律不改跟踪文件——把产码件拷到仓外 `D:/tmp/wisp33r10/` 改，再 `go test -c -overlay <json>` 建 exe，
在 `CWD=cmd/wisp` ＋ `PATH` 带 `third_party/sherpa-onnx` 之下跑 `-test.run AC13ColdStartPageOver`。
overlay 两个已知死法都避了：**成对给**（`{"Replace":{"<target>": "<replacement>"}}`，每发 JSON 原文都落进日志可审）；
路径用 `cygpath -m` 口径的正斜杠（无裸单反斜杠）。另加**落地自证**两把：脚本先断言被替换的原文片段在原文件中出现**恰好 1 次**，
并把 `replacement_differs_from_target=0|1` 与 `exe_bytes` 一起落件（⛔ 不是靠"颜色变了"反推落地）。

| 发次 | 突变内容 | buildrc | runrc | 颜色 | 期望 |
|---|---|---|---|---|---|
| `pristine` | 替换件与原件**逐字节相同**（`diff=0`） | 0 | 0 | **GREEN** | 不改色 ✓（正控一） |
| `swap` | **两发 `SetHtml` 的次序调换**＝框点名的反控 | 0 | 1 | **RED** | 必须红 ✓ |
| `shift` | 同函数前插 8 行注释/空行，**行号全位移** | 0 | 0 | **GREEN** | 内容锚自证 ✓ |
| `widen_swap` | 调换次序 **且**给探测页加 `<p>loading the panel</p>`（打掉第 2 问的牙） | 0 | 1 | **RED** | 第 1、3 问仍咬 ✓ |

`swap` 的红句逐字（`mut/mut-swap.log`，四句都响）：

```
panel_pageover_33r10_windows_test.go:261: after a cold start the LAST document handed to the control is byte-for-byte the document the round-trip probe step writes (136 byte(s)). That is ticket 33 AC#13: the user ends on the probe page instead of the panel. Product side: the probe stays - it is where cold usable is decided - but coldStartPageHandover must hand the entry over AFTER it
panel_pageover_33r10_windows_test.go:264: after a cold start the last document carries nothing a user could look at: with its scripts stripped its body is empty (head "<!doctype html><html><head><meta charset=\"utf-8\"></head><body><script>if(window.wispProbeRT)window.wispProbeRT();</scrip"). The panel was never handed over, so the panel window would show a shell. Ticket 33 AC#13 item 2 asks the final document about its content, not about whether a call was made
panel_pageover_33r10_windows_test.go:268: after a cold start the last document does not carry the embedded entry's own 1044 byte(s): ticket 33 AC#13's 'final document contains the real embed entry content' is not met
panel_pageover_33r10_windows_test.go:277: after a cold start the last document contains NONE of the 1 element id(s) the resolved entry declares [root] (head "..."): the document in the window is not the panel page
--- FAIL: TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe (0.00s)
```

`widen_swap` 的红句（`mut/mut-widen_swap.log`）＝第 261／268／277 三句仍响、**第 264 句不响**（那一枚牙按设计被打掉），
⇒ 台件不是"只有一枚牙"，也证明第 1 问不依赖 markup 字面量。

`pristine` 逐字读数：`capture: probe document is 136 byte(s)`、`world: embed resolves 1044 entry byte(s), 1 element id(s) [root]`、
`read: last document 1044 byte(s), 1 of 1 entry id(s) present, entry bytes carried=true`、`--- PASS`。

---

## §4 门各件 rc + 票框枚数前后

| 尺 | 命令 | 件 | rc |
|---|---|---|---|
| vet | `go vet ./cmd/wisp/` | `gate-vet.log`（零输出＝那格不是 0 字节没交，件里有 `rc=0` 行） | **0** |
| d22scan | `cd tools/d22scan && go run . -root ../../` | `gate-d22scan.log` | **0**；末行 `d22scan: clean - no D22 ban violations`；分母现量 `bans #1-5 internal/=228 cmd/=38`、`ban #6 frontend/=85`、`ban #7 internal/tools/=23`、`ban #8 design/=39 frontend/=85 internal/=514 cmd/=109` |
| gofmt（本腿两枚件） | `gofmt -l cmd/wisp/panel_host_windows.go cmd/wisp/panel_pageover_33r10_windows_test.go` | `gate-gofmt-mine.log` | **0**，列表**零行** |
| gofmt（整包上下文，⛔ 不顺手改） | `gofmt -l cmd/wisp/` | `gate-gofmt-package-context.log` | 列出 **3 枚**：`models.go`／`panel_inbound_guards_35r3_test.go`／`panel_transport_35r2_test.go`；行尾符尺现量 `models.go worktree_cr=334 head_cr=0`＝**CRLF 幻影**（HEAD 无 CR，工作树带 CR），另两枚 `worktree_cr=0 head_cr=0`；本腿**一枚未动** |
| 本腿族＋邻居（不开窗） | `go test -run '<六枚点名>' -count=1 -v .`，`CWD=cmd/wisp`＋`PATH=third_party/sherpa-onnx` | `gate-named-tests.log` | **0**；`=== RUN` 计数＝**6**、`--- PASS`＝**6**、FAIL／SKIP＝**0** |
| emoji | 台件与产码注释⛔ 任何 emoji；由 d22scan ban #8（`cmd/` 109 枚 Go 件，注释也扫）覆盖 | 同上 | 0 命中（`clean`） |

票 33 框枚数尺（同一把 `grep -cE '^[[:space:]]*- \[ \]'` / `- \[x\]'`）：
**追加前 13 未勾／1 已勾**（记在 `00-anchor.md`）→ **追加后 13 未勾／1 已勾**；节头枚数尺
`grep -c '^## 33-r10 进度追加' = 1`（⛔ 没重复追加）；`wc -c` 由 342 行件变 `101694` 字节，行数只增不减。

---

## §5 这格不许被读成什么

1. ⛔ **不许读成"面板今天开得出真页面"**。本腿证的是**冷启动尾段最终把哪一份文档交给控件**，
   以及那一份文档**是否等于宿主从 embed 解析出的字节**。常驻链的另一跳（`panel-sta` 线程真把窗口挂上、
   WebView2 环境 ready、页面自己的脚本跑起来把 `#root` 填出东西）**不在这格射程**，本腿没有任何一件证据支持它。
2. ⛔ **不许读成"用户看得见入口页"**。`SetHtml` 之后是否真渲染，本腿用的是记录型 sink，不是浏览器。
   真窗那一枚用例（`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`）本腿**没跑**（跑了就是开真窗，派单禁），
   它的色归验收腿／winlive 波次。
3. ⛔ **不许读成"embed 里带的是今天的新构建"**。本树 `go:embed all:dist` 编进去的入口是**未跟踪**的
   `frontend/dist/index.html`（`frontend/.gitignore:12 dist/*` 命中，`git ls-files frontend/dist` 只有 `.gitkeep`），
   1044 字节、引用 `assets/index-BVKlegVD.js`＝台账点过名的盘上陈旧件。
   所以第 3 问钉的是"最终文档＝宿主自己 Resolve 出来的那串字节"，⛔ 不是"这串字节是新的"。
4. ⛔ **不许读成"第 3 问在 CI 也有牙"**。CI／新鲜检出的 embed 只有 `.gitkeep` ⇒ `Built()=false` ⇒ 走 not-built 分支，
   那里只有第 1、2 问响；**"入口真内容"那一问在 CI 没有 subject＝具名欠账**（不是本腿能补的：补它要么动 `internal/panel`
   加导出接缝（禁面），要么在 `frontend/dist` 造跟踪件（禁面＋SLO 面））。
5. ⛔ **不许读成"第 2 问单独就够了"**。它会被"给探测页加一段可看 markup"打掉——`widen_swap` 那一发就是证据；
   单独顶住那一次的是第 1 问的身份比对。
6. ⛔ **不许读成"探测被削弱或删掉了"**。`firstRoundTripLocked` 一字未动，探测仍是冷启动"可用"判定的来源（AC#4）；
   `rtMs` 仍由 `bringUp` 记进 `m.lastColdMs`。
7. ⛔ **不许读成"本腿翻过任何 AC 框"**。票 33 只做了**末尾追加**；枚数尺前后同数（§4）。
8. ⛔ **不许读成"overlay 那几发证明过真窗里的行为"**。四发全在同一条不开窗的载具上；
   窗口依赖的 5 枚既有红本腿没跑也不记账。

---

## §6 留给验收腿的问题（本腿不代答）

1. 第 1 问用"运行时取到的探测页"作身份比对基线。若有人把探测页改成**与入口页逐字节相同**（极端形状），
   第 1 问失效、第 2 问仍绿、第 3 问会绿——这条链是不是可接受，还是要第 3 问之外再加一枚？本腿不裁。
2. `captureProbeDoc33r10` 里那枚 sink 的 `SetHtml` 会**把已绑定的无参回调当场调一遍**（建模"页面自己调探测绑定"）。
   这枚建模会不会让"往返"以某种方式恒真？（本腿的红句都只读 `docs`，不读 `rtMs`；但 `rtMs=0` 这类读数别当实测时延引用。）
3. 本腿把 `bringUp` 尾段搬进 `coldStartPageHandover`：这是**产码面为一枚台件开的接缝**。
   按 D22 双角色，验收腿要判"这个接缝是不是把测试形状渗进了产码"——它有没有比"在测试里复刻一次顺序"更好，
   以及下一枚腿会不会绕过 `bringUp` 直接调它而误读成"跑过冷启动"。
4. CI 里第 3 问永远没 subject（§5 第 4 条）。要不要另开票让 `frontend/dist` 的入口在 CI 里**由构建闸门产出**
   （票 274"构建带页面"那条链），还是接受"入口真内容那一问只在带 bundle 的机器上有牙"？本腿不裁。
5. 本树 embed 带的是陈旧入口（`BVKlegVD`）。验收腿若要引"本树这一发绿"当凭据，
   是不是必须先做一次仓外 fresh build 再对拉名册（台账第 115/125 条那个形状）？
6. 票 33 里 AC#13 那格今天应不应该继续按未勾处理、还是等真窗那枚用例的色一起裁？（翻勾归编排者，本腿没动框。）

---

## §7 与派单/票面冲突之处（点名到句）

1. **派单第 0 节**"工作树里有别人在飞的未提交改动……现量 30 枚上下" ⇒ 实量 `git status --porcelain | wc -l` = **753**。
   （30 枚那个数像是只数 `.gitignore`＋`design/**` 那一族；本腿一律未动、未 stage。）
2. **票面 AC#13 的整串行号**（`:221`→`:250`→`:227`→`:372-375`）与"每次冷启动最终显示的是那枚探测页" ⇒ **过期**，
   现量为 `:425` 探测 → `:427` 供页（换向出处 `13acad46`）。已在票 33 末尾具名追加。
3. **派单第 1 节转述**"完成判据②……面板真接进常驻之后用户看到的仍是空壳页" ⇒ 今天这句话的**产码那一半不成立**（§0(b)），
   但**仪器那一半成立**（本树唯一相关的既有用例买真窗、CI 里 skip），本腿按派单第 2 节的处置走，未造 bug。
4. **派单第 5 节**"两发正控缺一不作凭据" ⇒ 本腿交了**三发正控 + 一发反控**（pristine／shift／widen_swap／swap），未少交。
5. 其余照派单：⛔ 未 push、⛔ 未 `git add -A`/`.`、⛔ 未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`、
   ⛔ 未动 `docs/PLAN.md`／`docs/specs/**`／`frontend/**`／`ci.yml`／35 族测试件／SLO·golden·`thresholds.go`、
   ⛔ 未开真窗（无 `-tags winlive`）、⛔ 未勾任何框、证据件全部只建不删且⛔ 无一以 `.out` 结尾。

## §8 本件 rc

rc=0（本件全部读数来自现跑的命令；§4 那枚表逐件带自己的 rc 行；`mut/` 四发与 `mut-summary.log` 各带 `buildrc`/`runrc`/`rc` 三行；
无"跑到一半没交"的格子——未完项已在 §6/§7 具名）
