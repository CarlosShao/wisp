# 33-r10 锚件（起手尺 + AC#13 三格现量）

写码腿 `33-r10`。本件在**任何长跑命令之前**落笔并单独 commit。

## 0. 起手尺（现量，非转述）

- `date`：`Thu Oct  8 10:06:21 CST 2026`
- 起手 HEAD（`git log -1 --format=%H`）：`128bf600fa0890c02887e1d005dc7d6b1bd87a43`
  （不早于派单要求的 `128bf600`，是本腿自己现量，编排者的号若被后续提交作废与本腿无关）
- `git status --porcelain | wc -l`：`753` 枚（派单预估"30 枚上下"是**只数 design/.gitignore 那一族**的口径；
  实量 753 行含 `.scratch/**` 里别人在飞的未跟踪件与 `design/**` 的 15 枚 ` D`。
  ⇒ **冲突之处点名见 §4**；处置＝一律不动、不 stage、不"恢复"。）
- `wc -l cmd/wisp/panel_host_windows.go`：`861`（与派单一致）
- 票 33 框枚数尺（追加前）：未勾 `grep -cE '^[[:space:]]*- \[ \]'` = **13**；已勾 `grep -cE '^[[:space:]]*- \[x\]'` = **1**
  （与派单"现量应为 13 未勾／1 已勾"一致）；`wc -l .scratch/wisp/issues/33-panel-host-c27.md` = `342`

## 1. AC#13 原文（票 33 第 39 行整行，1413 字符，`awk 'NR==39'` 全文读，未 cut）

```
- [ ] **AC#13（10-01 11:2x 编排者追加，来路＝`33-v1` 验收表 §AC#3"供给那半"＋我自己现读；账 `A494`）：冷启动那段"往返探测"不许把真页面盖掉。** 现读时序（我 11:1x 亲自复认过行号）：`cmd/wisp/panel_host_windows.go:221` 调 `serveEntry()` → `:250` 把 embed 的入口字节 `SetHtml` 进控件，紧接着 `:227` 调 `firstRoundTripLocked(...)` → **`:372-375` 又一发 `SetHtml`，内容是一枚自造的探测页**（`<!doctype html>…<script>if(window.wispProbeRT)window.wispProbeRT();</script>…`）⇒ **每次冷启动最终显示的是那枚探测页，不是面板**。⚠ 这条不是"注释过期"：**面板真接进常驻之后，用户看到的仍然是枚空壳页**，而台件只会绿（它断的是"往返成立"，不是"页面是内容页"）。**完成判据**＝① 探测与供页面**次序重排**或改用资源请求过滤器（`AddWebResourceRequestedFilter` 在 `pkg/edge` 是**导出**的，同文件 `:34` 的注释自己就承认了这一点）；② 一枚**会响的**断言钉"最终文档里含 embed 入口的真内容"（问能力：解析出的条目／文档正文特征，⛔ 不问 `SetHtml` 被调用过）；③ 反控＝把两发 `SetHtml` 的次序调换，该用例**必须**红。⛔ 不许用"删掉探测"来糊——探测是 AC#4/冷启动"可用"判定的来源，删它会把另一格打空。
```

（以上一段＝`awk 'NR==39'` 的盘上原文逐字拷贝，未 cut、未改写。）

## 2. ★派单第 2 节：三格现量（逐格读数 + 处置，⛔ 不合成一句"没问题"）

尺＝`grep -n 'serveEntry\|firstRoundTrip\|SetHtml\|wispProbeRT\|AddWebResourceRequested' cmd/wisp/panel_host_windows.go`，
以下为**当前 HEAD `128bf600` 的工作树逐行读数**（含行号）：

```
420:	// hand the page over LAST. firstRoundTripLocked shows a document of its own,
421:	// so running it after serveEntry meant every cold start finished on the probe
425:	rtMs := m.firstRoundTripLocked(ctx, t0)
427:	if err := m.serveEntry(); err != nil {
438:// the embed carries no bundle. It is a SetHtml of a document this host writes, so
448:	w.SetHtml(`<!doctype html><html><head><meta charset="utf-8"><title>panel assets unavailable</title>` +
452:// serveEntry resolves the embedded index.html and pushes it into the control.
453:func (m *PanelManager) serveEntry() error {
467:	w.SetHtml(string(data))
689:// the first SetHtml (bringUp runs it ahead of firstRoundTripLocked / serveEntry, and an
716:// firstRoundTripLocked drives one JS -> Go cycle and returns its duration, pumping
730:// longer the last word: bringUp runs this BEFORE serveEntry, so the document the
732:func (m *PanelManager) firstRoundTripLocked(ctx context.Context, t0 time.Time) float64 {
744:	if err := w.Bind("wispProbeRT", func() string {
755:	// A page that immediately calls the probe. SetHtml of a tiny document keeps
757:	w.SetHtml(`<!doctype html><html><head><meta charset="utf-8"></head><body>` +
758:		`<script>if(window.wispProbeRT)window.wispProbeRT();</script></body></html>`)
```

### (a) 今天谁先谁后？

**探测先、供页后。** `bringUp` 在 **`:425`** 调 `m.firstRoundTripLocked(ctx, t0)`，
它内部在 **`:757-758`** 发那一发自造探测页的 `SetHtml`；随后 **`:427`** 才调 `m.serveEntry()`，
`serveEntry` 在 **`:467`** 把 embed 解析出的入口字节 `w.SetHtml(string(data))` 推进控件。
⇒ **派单里编排者 10-08 那次量到（"`:425` 探测先于 `:427` 供页"）成立**；
票面 AC#13 写的 `:221` 供页 → `:250` 真页 `SetHtml` → `:227` 探测 → `:372-375` 探测页覆盖，
**这四枚行号在 HEAD `128bf600` 的 861 行文件里全部对不上现在的形状**（`:221`/`:227`/`:250`/`:372` 处不是那些调用）。

### (b) 后一发是否真的把前一发的内容覆盖掉？

**覆盖方向已经反过来，且反过来那一支是"好的"方向。** 后一发是 `:467` 的入口字节，
它覆盖的是前一发 `:757-758` 的探测页 ⇒ 用户**最终看到的文档是 embed 入口页**，
探测页只在冷启动窗口里**短暂存在**（`firstRoundTripLocked` 泵消息的那一段，`:425` 到 `:427` 之间），
最坏是启动瞬间的一次闪动，不是"最终停在空壳页"。
补充读数：`serveEntry` 失败时走 `:427-429` → `serveNotBuiltNoticeLocked()`（`:448` 那一发），
最终文档是"面板资源不可用"的明示页，**也不是**探测页。
⇒ 今天**没有任何一条路径让最终文档＝探测页**。

### (c) AC#13 那条"最终显示的是探测页"是否已经不成立？

**作为"今天的产码事实"已经不成立**——`bringUp` 的注释 `:419-424` 与 `:729-731`
自己就具名写了这是 AC#13 的修（"AC#13's order, and it is the whole fix"／
"since AC#13 was fixed, is no longer the last word"），台账 `A688`/`A689` 记过这次换向。
⇒ 判语：**框的 ①（次序重排）已经落地；框的前提（"每次冷启动最终显示的是那枚探测页"）过期**。
处置＝⛔ 不为交活造一个不存在的 bug；按派单第 2 节末句，把新事实追加到票 33 末尾并具名"编排者 10-01 那句过期"。

**但 ② 这一格今天仍然有效，且实测有洞**（不是"已有仪器所以没事"）：
`grep -rn 'AC#13' cmd/wisp/*_test.go` 现量，唯一钉最终文档内容的那枚用例是
`cmd/wisp/panel_resident_windows_test.go:305 TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`。
它买的是**真窗**（`startPanelForTest` + `evalOnPanelThread`），而本票派单第 4 节
⛔ 禁 `-tags winlive`；且它在**这棵树**里 `:317` 直接 `t.Skipf`
（"AC#13 has no subject in this tree: the embed resolves no entry"，因 `panel.BuiltinAssets()` 在本树解不出入口 bundle）。
⇒ **本票 ②的射程不是"再写一枚真窗用例"**（那条链在 winlive，不在这格），
而是：**有没有一枚不开窗、跑在 CI 里、把"最终文档含入口真内容"钉住的仪器**。读数：没有——
`grep -rn 'AC#13' cmd/wisp/*_test.go` 命中的其余行是 AC#4 调度、WM_QUIT 根因（`:335` 起）、
第五形出口（`:592` 起），都不问最终文档的内容。
⇒ 本腿 ②的活＝在**假控件接缝**上立那枚会响的、开窗无关的内容断言（详见 `impl.md` §2）。

## 3. rc 自记

本件由 Read/grep/awk/wc 现量写成，无 `go` 长跑：**rc=0**（尺＝各命令退出零；`wc -c` 见 `files.txt`）。
