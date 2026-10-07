# 35-v1 起手锚（验收腿，票 35 AC#6 `:52` / AC#7 `:63`）

- date：`Wed Oct  7 18:47:21 CST 2026`（原样）
- 起手 `git log --oneline -3` 顶＝`ce6a4080`；其下 `1e42fc61`（35-r1 交件）、`fb2fb802`（35-r1 实现）。
  完整读数见编排者会话侧 logs/gitlog.txt。
- `git status --porcelain` 行数＝**753**（共享工作树既有脏，⛔ 非本腿产物；本腿只写 `.scratch/wisp/probes/35/v1/`）。
- 票面 `:52`（AC#6 Transport agreement page↔host）与 `:63`（AC#7 Inbound judge must ride the page's own envelope）
  两格已**逐字整行读**（Read 工具整行，未 `cut`）。票面 `:138-154`（编排者裁形＝甲子形①、禁区、N1/N2/N3/N6 恒真性记录）已读。
- r1 交件 `.scratch/wisp/probes/35/r1/impl.md`（149 行）与 `start.md` 已读。

## 本腿射程（简报 §3 五格）
- G1 定向突变两发：`:671` → `_ = panelPostMessageForwardInit`（应红）；`:671` → `w.Init("/* no forwarding */")`（反形，也应红）。md5 前后 + `git diff --exit-code`。
- G2 fake 保真度：对拉库源码 msgcb / Init 语义 / wispDispatch 桩帧形。
- G3 真 WebView2 一发：查 winlive 载体是否走过"页面 postMessage → Go dispatchRaw"这条边。
- G4 顺序证：`installPanelTransport`（`:408`）早于首个 `SetHtml`；票 `:144` 的 `:449/:468/:681` 行号漂移现量。
- G5 三条禁区：frontend/** 未动（现量）；C17 名册未变（四枚在册方法名逐枚对）；无"答掉 RPC 解析器却不经 dispatchRaw"新支路；`m.dispatchRaw(` 调用点枚数。
- 门禁：`sh scripts/d22scan.sh` + `go test -count=1 ./internal/...` 起终两发逐名红名册作差。

## 行号现量（本腿自量，与简报对拉）
`panel_host_windows.go`：`:408` bringUp 调 `m.installPanelTransport(w, ctx)`；`:633` `type pageTransport`；
`:648` `const panelPostMessageForwardInit`；`:664` `func installPanelTransport`；`:671` `w.Init(panelPostMessageForwardInit)`；
`:677` `func (m *PanelManager) dispatchRaw`（⚠ 票面 `:52`/`:63` 引的是 `panel_host_windows.go:630`＝**漂移 −47**，G4 具名报）。
SetHtml 现量＝`:448`/`:467`（票 `:144` 写 `:449/:468/:681`）。
测试：`panel_transport_35r1_test.go` 具名用例起 `:195`；`forwardingInstalled()` 起 `:127`；红断言行＝`:209`。

## 本机可测性口径
`go test ./cmd/wisp/` 载入即 `0xc0000135`（STATUS_DLL_NOT_FOUND，sherpa/onnxruntime 原生 DLL 不在 PATH）＝环境条件。
跑法＝`go test -c -o /d/tmp/wisp35v1/panel.test.exe ./cmd/wisp/` + 三枚 DLL 拷 exe 旁 + CWD=`cmd/wisp`。
⛔ 临时件全在 `/d/tmp/wisp35v1/`，绝不落进仓；只建不删。⛔ 不 push、显式 pathspec、不动 AC 框、不动 frontend/docs/specs/SLO/golden/thresholds。
