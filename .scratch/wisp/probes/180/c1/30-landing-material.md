# 180-c1 格 3 · AC#2 / AC#3 的落点料（只指认，⛔ 本程不落地、不动一字节产码）

---

## 1. 先说结论对方向的影响：AC#2 的"接那一跳"**已被票 255 AC#4 做过**

票面 AC#2 写的是「若 AC#1 判『规格要求生效』：把那一跳接上」。
现量：**`[panel] width` → 窗口尺寸这一跳今天已经在码上**，接它的是票 255 AC#4，不是本票。
链路（逐跳 file:line，`70b00885` 复核）：

```
resident_windows.go:151            newResidentPanelManager(rt.Layout.DataDir)
panel_resident_windows.go:228      func newResidentPanelManager
panel_resident_windows.go:253      NewPanelManager(..., withGeometrySource(panelGeometrySource(dataDir)))
panel_resident_windows.go:198-209  panelGeometrySource: 每次建窗 config.LoadFile, :207 return cfg.Panel.Width, cfg.Panel.Height
panel_host_windows.go:199-200      withGeometrySource(src func() (width,height int)) -> m.geometry = src
panel_host_windows.go:236-264      windowOptions(): gw>0 才覆盖 panelWidthPx=420 (:89) -> webview2.WindowOptions{Width: uint(width)}
panel_host_windows.go:392          建窗调用点
```

⇒ **本票若照票面原文再排一次"接那一跳"＝重复施工。** AC#1 的勾要改的是**前提**，
AC#2 的格要么判"已由 255 落地"、要么被重写成下面 §2 的真缺口 —— **这个裁决归验收腿与编排者，本程不勾任何框。**

---

## 2. 今天真正还缺的一跳：**重新生效那一跳**（不是接线那一跳）

`width` 的读者是**建窗时**读者，而 C27 的窗口是**单例＋隐藏而非销毁**，两件事撞在一起：

| 现量 | 位置 | 后果 |
|---|---|---|
| `Show` 只在 `!created` 时才 `bringUp` | `panel_host_windows.go:499-509`（`created := m.created`；`if !created { m.bringUp(ctx) }`） | 窗口一旦建过，再 Show **不会重跑 `windowOptions()`** ⇒ 不会再读 `[panel] width` |
| 关面板＝Hide，不是 Destroy | C27 原文 `PLAN.md:1377`「单例，**唯一** WebView2 窗口持有者；一会话至多一个面板窗口，**隐藏而非销毁**」 | 用户能做的"关窗再开"在码上是 Hide/Show 对，**不触发重读** |
| 唯一的 Destroy 是进程级 | `panel_resident_windows.go:454 stop()`、`:496 teardown(why)`；`RequestDispose()`（`:442`）**非测试调用者 0 枚**（尺：`grep -rn "RequestDispose" --include=*.go cmd/ internal/ \| grep -v _test.go` ⇒ 只剩 `:439` 注释与 `:442` 定义，`rc=0`） | 一个常驻进程会话内**建窗至多一次** |
| 没有 resize 路 | `grep -rn "MoveWindow\|SetWindowPos\|Resize(\|SetBounds" --include=*.go cmd/wisp/ internal/panel/ \| grep -v _test.go` ⇒ 产码 0 命中（只剩注释与 `cmd/wisp/testdata/esclistener/main.go`） | 已开的窗也不会自己变 |

⇒ **"宽度那一跳今天差几跳"的一句话答案：接线那一跳 0 跳不差；差的是 1 跳——"改完值之后没有任何一条路让已存在的窗口或再次 Show 的窗口重新取这个数"，真实生效条件是重启进程。**

⚠ 由此得出具名的一条**过期宣称**（不是本程修的）：
`cmd/wisp/panel_resident_windows.go:177-179` 的 `residentPanelGeometryNote` 与
`cmd/wisp/config_readers_255.go:161` 都写着「**面板关窗再开即跟上新值**」。
按 §2 那三行现量，这句话对"关窗＝Hide 再 Show"的普通路径**不成立**，成立的是"重启进程"。
**本程不动这两处**（`config_readers_255.go` 是票 255 的写面，`panel_resident_windows.go` 上有测试钉），只具名上报。

---

## 3. 若要补上 §2 那一跳，最少要动的文件（4 枚，按改动量排序）

1. **`cmd/wisp/panel_host_windows.go`** — 让 `Show`（`:499`）在 `created` 为真时也比对一次几何，
   或新增一条"几何变了才 resize"的路；今天的常量回退在 `:236-264`。
2. **`cmd/wisp/panel_resident_windows.go`** — `:442 RequestDispose` 今天无人可调；
   要么给它一个真的调用者，要么让几何来源在 Show 时被重取（`:198-209` 已经是每建窗现读的形状，不用改）。
3. **`cmd/wisp/config_readers_255.go`** — `:161` 那一行的措辞与 `:192 sectionReadSites["Panel"]` 的名册：
   ⚠ **新增一个读者文件会让 `TestTicket255RosterStillMatchesTheActualReadSites` 变红**（`:179-204` 是文件级钉），
   落地腿必须同时改这张名册，否则不是"改坏了测试"而是"名册与码不一致"。
4. （仅当选择把"此项不生效"透出给面板侧）**`internal/panel/composer.go` 的 `Snapshot`** ＋ 前端 —— ⇒ 见 §4，**这条路要人工批准**。

同段两枚**至今零读者**的字段若要一起接，落点同样落在 1／2 两文件内：
`PanelSection.Enabled`（`schema.go:540`，`default:"true"`，归属尺 rc=1 零命中）、
`PanelSection.KeepAliveInSession`（`schema.go:546`，`default:"true"`，全树只有定义＋注释）。

---

## 4. 契约面（只指认，⛔ 本程不碰、不裁）

会碰到的契约，逐一具名：

- **D36 配置模型**——`docs/PLAN.md:2743` 逐字：`| [panel] | enabled width height keep_alive_in_session(true) scale | hot |`。
  本程 §2 的现量说明 **`hot` 这一档今天对 `[panel]` 只有"重启进程"那一种兑现方式**。
  ⇒ 要么补实现兑现 `hot`，要么改规格文字。**改 `PLAN.md` 一字＝人工批准**（`AGENTS.md` §1.1）。
- **C27 `PanelManager`**——`docs/PLAN.md:1377`：单例／唯一 WebView2 窗口持有者／**隐藏而非销毁**／L2 降级不得消失。
  ⇒ 任何"重新生效"的实现都必须在**不违反"隐藏而非销毁"与"一会话至多一个窗口"**的前提下做；
  这是 C27 契约面，**动它的语义＝人工批准**。
- **C17 `PanelBridge`**——`docs/PLAN.md:1367`；若走 §3 第 4 项（把"不生效"透出给面板侧），
  就要动 **Snapshot 键集**或**方法白名单**。`180-a1` 第二程已把这三处标为 C17／契约面并写明"⛔ 本程不裁该不该加"。
  另：`AGENTS.md` §2 未定义即停清单里**确有**一条罩着它——
  「`C24 GojaHostAPI` 初始集与 `C17` 方法白名单定稿（S7/S5 切片卡批准）」。⇒ **碰到即停，本程已停，只上报。**
- **不碰清单核过**：本程对 `docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt` **零字节改动**（见 §5）。

⚠ 本程**没有**判定上述任何一格"应该怎么改"。§3 只是"最少要动哪几个文件"的落点料，裁哪一条归编排者／落地腿。

---

## 5. AC#3 那一格单独一句（派单点名要的尺）

**本程默认数值一字节没改。** 两把尺坐实：

```
git diff HEAD -- internal/config/schema.go                      => 空输出            rc=0
git diff 601c2b18..HEAD -- internal/config/schema.go | wc -l    => 0                 rc=0
```

第二把尺比第一把更强：它说的是**从起手锚到现在，包括其他在飞腿在内，`schema.go` 在整个区间都没被动过**
⇒ `default:"640"` 仍是 `schema.go:542` 上的 `default:"640"`，本票没有发生"把默认值改掉交差"那种形状。
本程自身对产码的改动更直接：`git diff --stat 601c2b18..HEAD -- internal/ cmd/` 里那两枚文件
（`panel_host_windows.go` / `panel_pageover_33r10_windows_test.go`）**属于 `33-r10` 腿，不属于本程**；
本程的全部写点只有 `.scratch/wisp/probes/180/c1/**` 与票面的追加一节。
