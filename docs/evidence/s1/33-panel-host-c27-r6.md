# 33-r6 面板宿主 C27 — 顺序依赖 WM_QUIT 毒线程修复（写码腿）

票 33 `.scratch/wisp/issues/33-panel-host-c27.md`。本腿只修 HEAD 上那一枚**真红**：
`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`（整包序红、隔离单跑绿）。

## 起手锚点（同发取，`date` + `git log` + `rev-parse` + `status --porcelain`）

- 时刻：`2026-10-01 14:22 +08`
- HEAD：`05f91799e1c88dfa2f96580e82e5cd9f5383d7a7`
  （标题＝`ledger(A503) + 撤票249: ... 抓到 1 枚真红（顺序依赖的 WM_QUIT 毒线程）...`）
- 分支：`dev`
- 起手名册：`.scratch/wisp/probes/33/r6/status-start.txt`（**242 行**，本腿不写"必须为空"，闸门＝终态名册等于起手名册＋只我这几枚路径）
- 台件目录：`.scratch/wisp/probes/33/r6/`

## 判据进度表

| 格 | 内容 | 状态 |
|---|---|---|
| ① | 最小复现（哪一枚先跑导致、组跑逐字读数） | 待填 |
| ② | 定性：产品缺陷 vs 测试卫生缺陷（两边说法） | 待填 |
| ③ | 修法 + 反控红句 | 待填 |
| ④ | 整包逐名红册（PASS/FAIL/SKIP 三数 + 逐名红） | 待填 |
| ⑤ | 门禁四数（build / vet / d22scan / gofumpt） | 待填 |
| ⑤⑥⑦ | 我可能写错的条目 / 我跑了哪些尺 / 判不动的地方（甲/乙/不做＋现量） | 待填（先写满） |

---

## ⑤⑥⑦ 自我对抗三节（先写满）

### ⑤ 我可能写错的条目（逐条，附"如果错了后果"）

- (a) **"归因＝hostThreadHarness 的 `runtime.UnlockOSThread()` 把带 WM_QUIT 的线程还给池"**：
  这是从依赖源码（`webview.go:242-243` 的 `case WMDestroy: w.Terminate()` + `:381-383` 的裸
  `PostQuitMessage`）＋编排者假设推出来的**形状**，本腿尚未用最小样本复现成确定性红。
  若错（真凶是别的形，如共享 dataPath 抢用、并发 bringUp、`CoInitialize` 那一支），
  那我的修法（去 Unlock）治不了真红 ⇒ 判据④那枚具名红仍在。
- (b) **"出货路径不受此害"**：只核了 `residentPanel.loop` 锁而不解（`panel_resident_windows.go:202`，无 Unlock），
  且 `notify_windows.go` 的 Unlock 落在 message-only 窗（非 webview wndproc）。
  未逐行读 live 测试与其它起窗点 ⇒ 若还有别的"锁了又解且那线程起过 webview 窗"，出货面也可能中。
- (c) **"AC13 单跑三发全绿"**：引编排者 `.scratch/wisp/probes/orchestrator/33r5-ac13-alone.txt`，**本腿未复跑**。
- (d) **毒线程复用是调度器决定的、非硬保证**：组跑可能偶发不复用同一 M ⇒ "红"也可能变"绿"。
  我的钉必须能在"没复用到时"也响（否则反控不成立）。

### ⑥ 我跑了哪些尺 / 哪些没跑

- 已跑（起手）：`ls .scratch/wisp/probes/33/`（名册 r6 未占）、`git status --porcelain`（存 status-start.txt）、
  `grep UnlockOSThread cmd/wisp/*.go`、读依赖源码 `chromium.go`/`webview.go`、读三份产码/测试文件。
- 待跑：最小复现组跑（①）、修法后整包（④）、门禁四数（⑤）。
- ⛔ 不跑：`staticcheck`（本机版与 CI 钉版不同 ⇒ 标〔未复认〕，不拿假结论）。

### ⑦ 判不动的地方（逐条：甲＝按我读数、乙＝别的形、不做）

- 甲：**去 `hostThreadHarness` 的 `runtime.UnlockOSThread()`**（让窗口线程随协程退出被销毁，
  带走的 WM_QUIT 也一起没了）——与产品 `residentPanel.loop` 同一不变式，落在我写面（`cmd/wisp/**` 测试）。
- 乙（**若甲读数不成立才转**）：把不变式落进产码——但产码 `bringUp`/resident 本就锁而不解，
  真正"锁了又解"只在测试 harness ⇒ 若定性为产码面要写死，只能加注释＋一枚会响的钉，不能加 Unlock 逻辑。
- 不做：**不**碰依赖（`pkg/edge`、module cache）；**不** `go mod`；**不**搬 AC13 进 winlive/t.Skip；
  **不**放宽那 15s；**不**把断言降级成 Logf。
