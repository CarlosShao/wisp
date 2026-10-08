# 票 255 · 255-r6 起手锚（落地腿，2026-10-08）

现量（全部本腿自跑，⛔ 未抄派单号）：

- `git log --oneline -1` = `040e424d`（派单转述的 `c4b8e2eb` 在其下方已漂，⛔ 未采）；分支 `dev`。
- `ls .scratch/wisp/probes/255/` = `a1 a2 c2 r1 r2 r3 r4 r5 v1 v2`，**没有 r6**（`ls .scratch/wisp/probes/255/r6/` 回 "No such file or directory"）⇒ 本格无人写过，本腿可写。
- 落点复认（`cmd/wisp/panel_geometry_255_test.go`，`wc -l` = **456**）：
  - `:208` `readHostFileForGeometry255(t, filepath.Join(hostPkgDir255(t), "panel_host_windows.go"))` HIT
  - `:210` `parser.ParseFile(...)` HIT
  - `:238-239` `if len(newWithOptions) != 1` + `"cmd/wisp/panel_host_windows.go has %d webview2.NewWithOptions call sites, want exactly 1"` HIT（**单文件射程**）
  - `:315`/`:318` imports 尺 HIT；`:343` 读 `panel_resident_windows.go` HIT；`:376-377` `retired` 名册（⛔ 禁区，本腿不碰）HIT；`:451` `os.ReadFile` HIT
  - ⇒ 派单七枚行锚逐枚对上，主体尺 = `TestTicket255PanelHostBuildsItsWindowOptions` 的 :239 "exactly 1"。
- 形状参照（`cmd/wisp/panel_locked_naming_33r11_windows_test.go`）：`:176 os.ReadDir(".")`（包目录、剔 `_test.go`、空名册＝Fatal）＋`:207 os.ReadFile(name)`＋`parser.SkipObjectResolution`、解析失败＝Fatal 不豁免。HIT。
- 基线事实（本腿现量）：
  - 全包产码 `NewWithOptions` 调用点仅 1 枚：`cmd/wisp/panel_host_windows.go:386`（其余命中全在注释/其它包/`.scratch`）。
  - `scripts/spike/webview2-latency/main.go:149` 是真·第二调用点（同 module、无独立 go.mod）⇒ **"整仓"射程下 "恰好 1" 基线即红**，而期望值⛔不许动 ⇒ 本腿选 **整包产码**（见 impl.md ①）。
  - `panel_resident_windows.go`（553 行）与 `resident_approval_windows.go` 均已 `webview2 "github.com/jchv/go-webview2"` 在 import ⇒ 突变可落兄弟产码文件且可编译。
  - 起手 `git hash-object` == `git rev-parse HEAD:`（逐枚相等）：`panel_geometry_255_test.go aea7e826…`／`panel_resident_windows.go 0f8912fd…`／`panel_host_windows.go 58e2b155…`。HEAD 副本已存 `mutations/baseline-*.go`（突变后用 `git show` 还原＋逐枚 hash 拉平，⛔ 不用 checkout/reset/stash）。
- 名册复跑口径：票面 `file:LINE [token]` 名册从 `.scratch/wisp/issues/255-*.md` 本腿自抽枚数（⛔ 不抄"14 枚"），交付前后各跑一遍。

本腿计划：(i) 扩射程——把 :239 的"恰好 1"断言扫到整包产码，测试件净零增删行（`wc -l` 456→456 为"没挪行号"的尺）；突变全部盘上种（⛔ overlay 不用于源读尺证明，A713）。
