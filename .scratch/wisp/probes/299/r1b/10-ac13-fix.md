# 票 299 接续腿 `299-r1b` · 第 1 笔：`AC#13` 次序回归（红→绿两发＋反控）

时刻 `2026-10-13 09:1x +08`。锚＝HEAD `91c1d8ab`（＝代提件 `3b806876` ＋派单存档；工作树 `cmd/wisp/**` 与 `3b806876` 逐字节相等，尺＝`git diff 3b806876 --stat -- cmd/wisp/` 空）。

**判语（一句话）**：⛔ **实现侧一字未改，也不需要改** —— `coldStartPageHandover` 的次序（探测页先、真页面最后）从票 33 到 `3b806876` **本来就没有一丝偏离**；派单据以判"红"的那发日志（`probes/299/r1/logs/run-realon-rulers.txt`）是**重锚前的单发读形**在被撞红，而重锚（`awaitEntryReport` 有界重问）是 r1 在那一发**之后**写进测试件、**从未跑过**的。本腿在已提交形状上现跑：**两枚钉子全绿**（§3），且反控证明两枚钉子对次序**仍然敏感**（§4）。⛔ 没有放宽任何断言。

---

## 1. 红发（r1 的，逐字，`probes/299/r1/logs/run-realon-rulers.txt`）

```
=== RUN   TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe
    panel_resident_windows_test.go:315: AC#13 probes from the resolved entry (1044 bytes): 1 id(s) [root]
    panel_resident_windows_test.go:327: AC#13 page answer (head 05b5f11f): "0" - 0 of 1 probe id(s) present in the live document
    panel_resident_windows_test.go:330: after a real cold start the live document contains NONE of the 1 element ids the embedded entry declares (the page itself answered "0")... 
--- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (0.74s)
=== RUN   TestAC13BringUpSurvivesAReusedThreadQuit
FAIL	github.com/CarlosShao/wisp/cmd/wisp	2.006s
FAIL
rc=1
```

（第 2 枚在那一发里**没有哪怕一行自身输出**、进程 2.006s 即整体 rc=1 —— 属"日志件被截"那类形状；本腿在已提交形状上现跑它是 **PASS**，见 §3。）

### 为什么这发红 **不是** 次序回归：三条具名读数

1. **时长**：那发 FAIL 用时 **0.74s**。重锚后的失败路径要**花满预算**才能落红：`panelThreadWait = 15s`（`panel_resident_windows_test.go:53`），而本腿的反控（§4）在 mutate 世界里实测 **15.69s**。0.74s ⇒ 那发跑的**不是**重锚读形。
2. **日志里没有重锚的两条专属行**：重锚读形跑一次必打二者之一 —— `AC#13 entry-id report arrived on probe attempt N`（`:361`）或 `AC#13 no probe attempt over N attempt(s)...`（`:370`）。那发两条都没有。
3. **行号对拉**：那发的 `:327`/`:330` 与 **`05b5f11f` 旧件**的 `AC#13 page answer` Logf / `t.Errorf` **逐行对上**（`git show 05b5f11f:cmd/wisp/panel_resident_windows_test.go | sed -n '327p;330p'`＝那两行）；已提交件里它们在 `:404`/`:407`。而探针 Logf 那发报 `:315`、旧件与已提交件都在 `:289` ⇒ 那发跑的是 r1 一个**中间 WIP 形状**（先于最终重锚）。⇒ 结论：r1 在红发**之后**又改了测试件、**再没跑过**它，于是"红"被留在盘上、真相没被重新量过。

## 2. 实现侧：次序语义逐字节原封

`panel_host_windows.go:452-459`（现件）＝

```go
func (m *PanelManager) coldStartPageHandover(ctx context.Context, t0 time.Time) float64 {
	rtMs := m.firstRoundTrip(ctx, t0)

	if err := m.serveEntry(); err != nil {
		m.serveNotBuiltNotice()
	}
	return rtMs
}
```

与 `05b5f11f` 的同名函数**逐行相同**（`git show 05b5f11f:cmd/wisp/panel_host_windows.go` 对拉）；`3b806876` 对该文件的 diff 里唯一点到 `coldStartPageHandover` 的行是一段**注释**（挂在 `bringUp` 的注册调用上）。⇒ 探测页（`firstRoundTrip`）先、真页面（`serveEntry`→`Navigate(panelEntryURL)`）后的次序语义从票 33 至今**没有任何偏离**。唯一变了的是供给形状（`SetHtml` → `Navigate`＋资源过滤器，票 299 `AC#1` 预裁 ⓓ 允许），它把交接从"同步"变成"异步"——旧单发读形在异步世界里"探针早发、读到还是旧文档"，这正是票价注释里写明的重锚理由。

## 3. 绿发（本腿现跑，`logs/run-ac13-green.txt`，rc=0）

```
=== RUN   TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe
    panel_resident_windows_test.go:394: AC#13 probes from the resolved entry (1044 bytes): 1 id(s) [root]
    panel_resident_windows_test.go:402: AC#13 entry-id report arrived on probe attempt 2: "1"
    panel_resident_windows_test.go:404: AC#13 page answer (head 91c1d8ab): "1" - 1 of 1 probe id(s) present in the live document
--- PASS: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (0.81s)
=== RUN   TestAC13BringUpSurvivesAReusedThreadQuit
    panel_resident_windows_test.go:560: AC#13 reused-thread root cause: planted one WM_QUIT on this locked thread, then ran bringUp - panicked=<nil> err=<nil> created=true
    panel_resident_windows_test.go:562: AC#13 reused-thread release: tid=30876 dispatched 3 message(s) before unlocking; windows left on that thread=0 queue head=empty
--- PASS: TestAC13BringUpSurvivesAReusedThreadQuit (0.71s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	1.659s
```

第 2 枚还满足复线程卫生（thread windows=0、queue head=empty）。同一批在两小时后于**最终字节**上复跑仍绿（`logs/run-final-green.txt`，6 枚钉子全 PASS，rc=0）。

**假绿分析**（为什么"attempt 2 报 1"不可能是错文档上蒙对）：报告位图由 `probeJS` 逐枚查 `document.getElementById(id)` 拼出，探针 id `[root]` 取自**入口自己的字节**（`entryIDProbes`），探测页与 about:blank 都不含 `#root` ⇒ 只有"真页面已在窗口中"才会报 `"1"`。

## 4. 反控：把次序翻过来，两枚钉子都必须红（`logs/run-m1-ac13-order-swap.txt`）

变异体 `mut/m1-panel_host_windows.go.txt`（`-overlay` 驱动，仓内文件零字节未动；overlay 件 `mut/m1-overlay.json`）：`coldStartPageHandover` 里两跳对调（serveEntry 先、firstRoundTrip 后）＝把 AC#13 那个缺陷**原样造回来**。读数逐字：

```
    panel_pageover_33r10_windows_test.go:282: AC#13 33-r10 read: 2 handover act(s), first=nav:"https://wisp.panel/index.html" last=doc:"<!doctype html>...window.wispProbeRT();</scrip" ...
    panel_pageover_33r10_windows_test.go:287: after a cold start the LAST handover act is doc:..., want a navigation to panelEntryURL "https://wisp.panel/index.html" ... A last act that is a document is ticket 33 AC#13: the user ends on the probe page instead of the panel.
    panel_pageover_33r10_windows_test.go:291: the FIRST handover act is no longer the round-trip probe document (nav:...) - ...
--- FAIL: TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe (0.00s)
    panel_resident_windows_test.go:404: AC#13 page answer (head 91c1d8ab): "0" - 0 of 1 probe id(s) present in the live document
--- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (15.69s)
rc=1
```

⇒ 两枚钉子（真窗那枚＋无窗 33r10 那枚）在"最后的文档是探测页"的世界里**都红**，且真窗那枚花满 15s 预算、由**页面自报** `"0"`。**判据没有放宽**：重锚只改了"何时提问"（有界重问），没有改"问什么/以谁的话为准/什么算过"。

## 5. 本笔未做与射程

- ⛔ 未改任何实现代码来"修" AC#13（无需：§2）。本腿对 `cmd/wisp/**` 的全部改动＝两枚件的**注释**（另一笔交件：悬空引用与 ban #9 修复，见 `30-gates.md` §g2）。
- 本腿红发为**引用 r1 的日志**（`.scratch/wisp/probes/299/r1/` 不在本腿写面，⛔ 未动）；绿发与反控为本腿现跑、原件在 `logs/`。
- 收口口径：AC#13 的实现侧＝`3b806876` 已交的次序；测试侧重锚＝`3b806876` 已交的读形；本笔的交付＝**这两者在同一棵树上第一次被一起跑绿**＋反控证明仍敏感。⛔ 翻任何 `AC` 框、⛔ 碰任何 `Status:` 行。
