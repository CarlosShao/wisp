# 35-r1 起手锚（票 35 AC#6 `:52` 那格）

本件只记起手时刻的现量，供交件对账；⛔ 不含结论。

- 日期：`Wed Oct  7 17:35:57 CST 2026`（date 原样）
- 起手 `git log -1 --oneline` 锚：`10899b77`（A678 收 274-r3）
- 分支：`dev`
- `go vet ./cmd/wisp/` 起手持数：rc=0，elapsed ≈ 1.6s（见 logs/ 内 vet-start）
- `git status --porcelain` 起手快照：logs/porcelain-start.txt（脏计数＝工作树既有，⛔ 非本格产物）

## 本格射程（只 AC#6 `:52`）
- 落点：在 `cmd/wisp/panel_host_windows.go` 已绑定的 `window.wispDispatch` 处接一条 Go 侧
  传输转发（甲子形①，票面 `:150` 已裁），把页面 `window.chrome.webview.postMessage` 的
  原始信封串喂进 `dispatchRaw`（`:630`）；⛔ 不动 `frontend/**`、⛔ 不改 C17 名册、⛔ 不新造 C##。
- 判据：新写一枚 `cmd/wisp/*_test.go`，喂 `panel.ts:246` 正控信封原文
  （`sendRequest("panel.mode.request",{to})` → `{"method":"panel.mode.request","requestId":"pc-<seq>-<uuid>","source":"panel-composer","to":...}`），
  经这条边到达 `dispatchRaw` 下游且带 `requestId`。
- ★定向突变：删转发语句 ⇒ 新用例必须红；恢复 ⇒ 绿。两发 `-v` 落 logs/。
- AC#7（`:63`）的 (d)-1/-2/-3 三支红文案归下一枚腿，本格不加其断言。

## 既有红（本格不碰、也不归本格）
`composer_test.go:73-77` 今天就是红的（Go/TS 键集双向 subtract，未声明键 `[instructions tasks]`
与 `[git currentModel modelKnown credentialState credentialKnown]`）。本格带着它交件，不放宽任何断言、不动 `frontend/**`。
