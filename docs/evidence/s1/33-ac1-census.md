# 票 33 AC#1 现量（唯一剩余格）

- 锚：腿码 **33-ac1**（2026-11-12 10:2x +08，工作树 `'D:\work\workspace\projects plans\Wisp'`）
- 本腿写面：**仅本枚 md**（⛔ 未改任何 `.go`、票面框、`Status:`）
- 起点判语：`33-v2` 终裁 AC#1 **不成立**，两处理由：①"进程树枚数稳定"尺会误归因；②"≤2s 退净"进了 winlive 但未量。本腿对两处**现核＋现跑**。

## AC#1 票面逐字（`.scratch/wisp/issues/33-panel-host-c27.md:72-73`）

> - [ ] Show/hide/destroy lifecycle: second show within session reuses window (process-tree
>       child-count stable); session dispose destroys + WebView children exit ≤2s.

## 33-v2 的两处问题，逐处现核

### 问题① "尺会误归因"（旧尺把正常子进程回收判红并错怪 reuse）

- **那条旧尺今天已不存在。** 旧尺＝pre/post 全树 count 相等断言，被 **255-r5** 换掉，证据在
  `cmd/wisp/panel_host_windows_test.go:684-709`（逐字："Why the identity is NOT a strict
  pre/post tree-count equality, measured by 255-r5 … the verdict is carried by host pids …
  **the counts are reported, not judged**"）。
- **现在的尺**（同文件 `TestPanelHostRealWindowHopAndLifecycle`，:614 起）按 host-pid **身份**判：
  - re-show 不得**新增** browser host（:755-762，`added` 半，旧 count 形根本看不见等量换防）；
  - 本 run bring-up 启的 host（`mineHosts`＝pre 减 baseline，:719-725）不得在过渡中死（:763-778）；
  - 全树 webview 枚数 / machine-wide 枚数**只作 reading 上报不判**（:779-780 的 denominators 日志、:827-828 的 machine-wide 旁注）。
- **误归因形状已消**：WebView2 helper 短命进程按自家节奏自然退出（旧尺 10 发里 2 发判红的那 shape），
  现在由两道机制兜住——`settleTreeReading`（:358-376，有界单调 settle，吸收 bring-up 周边的 helper）＋
  baseline 豁免（早于 bring-up 的 host 死亡不算在本次过渡头上，:764-768）。
- **本腿现跑读数**（`go test ./cmd/wisp -run '^TestPanelHostRealWindowHopAndLifecycle$' -count=1 -v`，rc=0，PASS 1.35s，
  输出存 `D:\tmp\wisp\20261112\33ac1\run-lifecycle.out`）：
  `AC#1 denominators (head 9e1a8dde5): our tree webview=6 direct-children=1 tree-pids=6 | browser hosts pre=[66864] mine=[66864] post=[66864] added=[] mine-died=[] | machine-wide msedgewebview2=66 | same HWND 0xa840e28 across hide->re-show=true`
  ⇒ 同一 HWND 跨 hide→re-show 复用、无新增 host、本 run host 全存活。热 re-show 45.888 ms（≤200ms D32 预算，:786-792 同测内绿）。

### 问题② "≤2s 退净未量"

- **有断言**：`TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive`
  （`cmd/wisp/panel_host_windows_live_test.go:36-74`，`//go:build windows && winlive`）。
  bound **原样 2s 未放宽**（:15-16 逐字："The bound is UNCHANGED at 2 seconds. Nothing was widened…"），
  分母是本进程树不是整机（:61-62、:70-73）。
- **本腿跑了**（`go test ./cmd/wisp -run '^TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive$' -count=1 -v -tags winlive`，
  rc=0，PASS 1.12s，输出存 `D:\tmp\wisp\20261112\33ac1\run-winlive-exit.out`）：
  `winlive exit clause: our tree webview baseline 0 -> peak 7 -> now 0 (tree pids 0) | machine-wide 60 -> 60`
  ⇒ Destroy 后树内 7 枚 webview 子进程在 2s 内退净（0 残留）。
- ⚠ 该条在 winlive tag 层、**无 CI job**（:18-20 自述），只桌面机可量——引用它必须带这句。

## 结论

**AC#1＝成立。** 理由：
1. reuse 半：票面"second show within session reuses window"由 host-pid 身份判据＋同 HWND 断言承判，今天实测绿（问题①的旧尺已删，误归因形状不存在）。
2. dispose 半：票面"exit ≤2s"断言在 winlive 层存在且 bound 未放宽，今天实测 7→0 退净。
3. 括号里的"process-tree child-count stable"措辞已被 255-r5 判定为不可测的等量换防盲区，现尺（host 身份）是**更利**的判据——"no assertion threshold moved and no assertion was deleted - the dimension asserted is sharper"（`panel_host_windows_test.go:740-746`），枚数维只降级为上报。

## 本腿 ⛔ 做到的事

- ⛔ 未改任何 `.go`（写面＝本 md 一枚）；⛔ 未跑整包，只 `-run` 指名两枚；
- ⛔ 未动票面任何框／`Status:`；临时件只落 `D:\tmp\wisp\20261112\33ac1\`（run-lifecycle.out/.err、run-winlive-exit.out/.err）；
- 本腿只 commit 不 push。
