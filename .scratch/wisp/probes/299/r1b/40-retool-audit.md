# 票 299 接续腿 `299-r1b` · 第 4 笔：五枚测试件"重锚"逐枚复看（删掉的断言有没有等价替代）

时刻 `2026-10-13 09:1x +08`。锚＝`05b5f11f`（改前）⇄ `3b806876`（代提的重锚形状，＝现有工作树字节）。方法：逐件 `git diff 05b5f11f 3b806876 -- <件>` 全读删/加两侧，再对**每一处被替换的断言**配一枚**反控**——⛔ 不接受"看起来没问题"。反控共六发（M1／M2 属前两笔；本笔用 **M1＋T4／T5／T6**），全部落在 `logs/`。

**总判语**：五枚件里 **0 处"净删除"**——所有被删掉的断言都有**等价或更强的替代**；两枚件（`panel_geometry_255_test.go`、`panel_pageover_33r10_windows_test.go`）的替换**换了问法**（因为票 299 `AC#1` 把供给形状从 `SetHtml` 换成"虚拟 origin＋导航"，预裁 ⓓ），换法逐处给；每一件都有反控逐字证明**对原来的缺陷仍然敏感**。

---

## 1. `cmd/wisp/panel_geometry_255_test.go`（+73−?）

**删掉**：ruler 3 里 `newWithOptions` 计数与 `WindowOptions:` 键值断言（旧三连：键的值必须是调用 → callee 必须是 `windowOptions` → 记 `foundGeometryCall`），以及 `NewWithOptions` 的 selector 匹配；变量改名。
**替代**：create 点改数 `createSites`（`calleeName299r1(node) == panelCreateCall299r1`）；**计数断言原样保留**（"本包产码恰好 1 处、且在被点名文件"）；几何断言改问**首参**：`create.Args[0]` 必须是 `*ast.CallExpr`（否则 `t.Errorf`），callee 必须是 `windowOptions`；另加零参 `t.Fatalf`。"无 int 字面量给窗定尺寸（`windowOptions` 体外）"那一半**一字未动**（现件 `:270` 起仍在，diff 未碰）。
**等价性**：旧断言问"喂进唯一建窗点的几何是不是 `windowOptions()` 的答案"；新断言问同一句话，只是新供给形状里那枚答案从"`WindowOptions` 键"变成"首参"（`panel_host_windows.go:391` `newPanelControl(m.windowOptions(), m.dataPath)`）。
**反控（两发，各自单杀一半）**：
- **T4**（把首参换成局部变量，`-overlay` 够不到这枚读盘尺 ⇒ **临时改盘**；跑完即复原）：`logs/run-m4-audit-mutant.txt` 逐字 `panel_geometry_255_test.go:263: the create's first argument is a *ast.Ident, not a call: the geometry is not being resolved at create time`；`--- FAIL`、`rc=1`。复原脚注逐字 `RESTORE-BYTE-EQUAL yes`（`45c1059b…`）。
- **T5**（在产码里追加第二枚建窗点）：`logs/run-m5-audit-mutant.txt` 逐字 `:247: cmd/wisp production sources carry 2 newPanelControl call sites [panel_host_windows.go:1014 panel_host_windows.go:391] and cmd/wisp/panel_host_windows.go itself carries 2, want 1 and 1`；`rc=1`。复原 `RESTORE-BYTE-EQUAL yes`。
**结论**：等价替代成立，**且两侧各有独立反控**；正发绿（`logs/run-audit-positives.txt`）。

## 2. `cmd/wisp/panel_geometry_255r6_range_windows_test.go`（+53−?，改名＋普查面）

**删掉**：`newWithOptionsSitesInPackage255r6`（只认 `SelectorExpr` 的 `NewWithOptions`）与其注释；替换为 `calleeName299r1`（`Ident` 与 `Selector` 都认）+ `createSitesInPackage299r1` + 常量 `panelCreateCall299r1 = "newPanelControl"`。
**等价性**：普查的**承诺不变**（"本包产码恰好 1 处建窗点"），主体换成新调用名；**匹配面是变宽而不是变窄**——旧形状只认带接收者的 selector，新形状两种拼法都收（这正是 `newPanelControl(...)` 无接收者写法能进射程的原因）。"读不了/解析不了该文件 ⇒ `t.Fatalf` 而不是跳过"的性质未动（diff 只替换了 Inspect 内部的匹配行）。
**反控**：T5 同一发同时命中普查（上面 T5 逐字里 `carry 2 … call sites` 就是这枚件打的）。正发绿。
**结论**：无删除；替换为更宽的匹配面，⛔ 非放宽。

## 3. `cmd/wisp/panel_host_gate_test.go`（+33−?，纯扩宽）

**删掉**：无。唯一结构变化＝单文件判定扩成 `panelHostL1Files = {"panel_host_windows.go", "panel_control_windows.go"}` ＋循环＋报错带文件名（禁项名册原样：`net`／`net/http`／`net/netip`／`net/mail`／`net/tcp`）。
**扩宽非放宽**：多扫一枚文件＝射程变大；断言集合一字未减。
**反控**：**T6**（给新件加 `_ "net"` 导入，临时改盘）：`logs/run-m6-audit-mutant.txt` 逐字 `panel_host_gate_test.go:50: panel host source panel_control_windows.go imports "net", which can open a listening socket; D29/AC#3 forbid a localhost server`；`rc=1`；复原 `RESTORE-BYTE-EQUAL yes`（回该件当刻快照 `7d257021…`，＝腿起点哈希）。
**结论**：纯扩宽；反控证明**新件真在射程里**（不是只改了名册）。

## 4. `cmd/wisp/panel_pageover_33r10_windows_test.go`（+146−?，无窗 33r10 尺）

**删掉**（旧三问的旧问法）：(a) `last != probeDoc` 不等断言；(b) `hasLookableContent(lastdoc)`；(c) 末文档 `strings.Contains(last, entryBytes)` ＋ id 命中在**末文档**；helper `lastDoc()` 删除（`firstDoc()` 仍被 `captureProbeDoc33r10` 用，`grep` 现件 `:243`，非死码）。
**替代**（新三问，件头注释逐条对应）：(a) **末 act 必须是 `nav` 到 `panelEntryURL`**——比"≠探测页"更强（探测 act 只能是 `doc`，"≠"被 `kind` 断言蕴含）；(b) **首 act 必须是运行时捕获的探测文档**——**新增的第二颗牙**（次序修复=探测先、入口后，探测被拿掉也红）；(c) 入口 URL 的应答**由过滤器自己的查找回答**（`panelAssetLookup`）：`status==200`、content-type 含 `text/html`、`bytes.Equal(data, entryBytes)`、`hasLookableContent(data)`、≥1 枚入口 id。
**等价性**：旧 (b)(c) 问的是"末文档的字节"；新形状里"末文档"=导航+过滤器应答的产物 ⇒ 问 `panelAssetLookup`（同一次 `Assets.Resolve` 缝）就是问那枚文档的**唯一来源**；`bytes.Equal` 比 `strings.Contains` 更严。手机身不在于是"文档被交进去了"而是"那 URL 的应答就是入口的字节"——与 `AC#1` 供给形状一致的等价问法。
**反控**：**M1**（次序对调）：`logs/run-m1-ac13-order-swap.txt` 逐字两红——`:287 LAST handover act is doc:…, want a navigation to panelEntryURL`＋`:291 the FIRST handover act is no longer the round-trip probe document`；正发绿（`run-final-green.txt` 逐字 `first=doc:…probe… last=nav:"https://wisp.panel/index.html"`、`1 of 1 entry id(s)`、`bytes-equal-to-resolved-entry=true`）。
**结论**：0 处净删除；被挪走的一问（末文档字节→URL 应答字节）＝供给形状变更的必然重问，反控证明对次序敏感。

## 5. `cmd/wisp/panel_resident_windows_test.go`（+91−?，真窗 AC#13 尺）

**删掉**：单发读形（`evalOnPanelThread(...probeJS...)` ＋ `awaitReport(h,"ac13-probe","AC#13")`）＋旧失败文案；`probeJS` 本体**未动**（只加注释）；`awaitReport` 仍被 AC#14 两枚用例用（现件 `:898`/`:943`，非死码）。
**替代**：`awaitEntryReport`——**有界重问**（预算 `panelThreadWait=15s` 单调钟；每次至多 `probeAttemptSettle=2s` 等自己的回音），断言体不变：`present = Count(answer,"1")`，`present==0` ⇒ `t.Errorf`；成功路径加 `attempt N` 读数行，失败路径加"最后报告"读数行。**过线标准没变**：答案是页面在 **live document** 上自报的位图且须含至少一枚 `"1"`；无回答/全 `"0"` 都红（旧形 `awaitReport` 无回答是 `Fatalf`，新形是无回答⇒`""`⇒`Errorf`，两级都是硬失败，⛔ 不是放宽）。
**假绿不可能性**：位图由 `document.getElementById(<入口自己的 id>)` 拼出；探测页与 about:blank 都不带 `#root` ⇒ 只有真页面在窗里才可能报 `"1"`（探测 id 由 `entryIDProbes` 从 `Assets.Resolve(EntryFile)` 的字节里现抽）。
**反控**：**M1**——真窗那发在 mutate 世界里逐字 `AC#13 page answer: "0" - 0 of 1 probe id(s) present`、`--- FAIL (15.69s)`（花满预算＝重问确实在问，且问的是探测页）。正发绿：`attempt 2` 报 `"1"`、`1 of 1`。
**结论**：读形重锚（异步导航所需），判据本体未放宽；反控逐字证明红。

---

## 附：本笔另核的两处

- `panel_assets_299_windows_test.go`（`3b806876` 新增件，⛔ 不在派单五枚内）：其 `AC#2 ⓑ` 反控承诺由本腿第 2 笔兑现（`20-ac2-reverse.md`）；悬空 transcript 引用一并治（见 `30-gates.md` §g6）。
- 五枚件的**产品侧**对应物（`panel_host_windows.go`＋`panel_control_windows.go`）换的是供给形状与建窗点，⛔ 未动次序（见 `10-ac13-fix.md` §2）。
