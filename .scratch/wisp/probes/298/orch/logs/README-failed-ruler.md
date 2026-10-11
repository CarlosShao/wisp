# 这两枚是**坏尺件**（⛔ 任何一格的凭据，只当反面标本留着）

- 件＝`blob-gofmt-d-panel_host_windows_live_test.go.txt`／`blob-gofmt-d-models_windows_live_test.go.txt`，各 39 字节。
- 坏在哪：我当时把 **`cmd/wisp/models_windows_live_test.go`** 当成 `cmd/wisp` 里带 `winlive` 标签的测试件，而那枚路径在 HEAD 里**根本不存在**
  ⇒ `git show HEAD:cmd/wisp/models_windows_live_test.go` 报 fatal、`gofmt -d` 的 diff 是**空的**（hunks＝0），"空"这一个读数⛔ 说明任何事（它既⛔ 是"干净"，也⛔ 是"脏"，它＝"没量到"）。
- ★**尺形也写错过一次**（就地具名）：我第一版把尺写成 `grep -rln '^//go:build winlive' cmd/wisp/*.go` ⇒ 命中 **0**，因为盘上的标签行是**合取形**。真尺＝
  `grep -rln '^//go:build windows && winlive' cmd/wisp/*.go`，08:3x 现跑命中 **7** 枚（行号与枚数一律当快照）：
  `panel_geometry_255_winlive_test.go`／`panel_host_windows_live_test.go`／`panel_transport_live_35v2_windows_test.go`／
  `resident_approval_live_246_windows_test.go`／`resident_ball_live_228_windows_test.go`／`resident_hotkey_live_258_windows_test.go`／`resident_task_source_live_246_windows_test.go`。
  ⚠ 旁证一枚形差：`resident_audio_247_live_windows_test.go` 的标签行逐字是 `//go:build windows`（**⛔ 含 `winlive`**）⇒ 按文件名里的 `live` 猜标签＝又一次假读数。
- 票 298 的 `AC#0` 那一格**真正量过的那把尺在落地腿 `298-r1` 的件里**（`probes/298/r1/**`，10-10 16:2x：工作树 `gofmt -l cmd/wisp` 3→1／HEAD blob 那把 2→0、两枚各 `numstat 2 2`、首笔即脏归因 `7f9d6e40`／`286a7f30`），
  五格框等的是**非实现者 `298-v1`** ⇒ 这两枚坏件⛔ 替任何一格背书。
- 关联账＝台账 **`A831` 第 4 节**（编排者自己那笔"在正文里叙述了盘上没发生过的读数与提交"的具名自报）。临时件只建不删 ⇒ 这两枚留着、⛔ 改名、⛔ 删。
