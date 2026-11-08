# 票 276 复认：维持 A1556「不做/moot」＋前提再变的新读数（2026-11-08 22:4x +08，写腿；零产码零台件触碰）

## A1556（2026-11-05）复认 —— 出处逐行对上

- `attrib.sh:577` 现存 `A_N_EXCLUDE=':(exclude).scratch/**'`（票 362-r2 落，注释含 HOW TO REVERT）——
  A1556 引的 577 行属实，排除在生效。
- 17 枚普查件（票 276 第 1 条名册）**17/17 仍存在**于 `.scratch/wisp/probes/`（逐枚 `[ -f ]` 现跑），
  未删、未格式化；两枚具名台件（241/v1/posctl、259/r1/gofumpt-negctl）未触碰。
- 结论：**AC#0 的处置对象（exit 1 名册 17 枚）已被 362 的 exclude 整体移出门分母**，本票维持「不做/moot」，
  不派写腿（A669/A1496/A1556 三度裁定一致）。

## 前提再变（新事实，归口待编排者）

A1556 写「门今天 green」。2026-11-08 现跑 `sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only` ＝ **RED（rc=1）**，
但红点已**不是** .scratch：**8 枚全在 `forks/go-webview2/`**（tracked-dirty=8/8，最后触碰
`2bb639e4e` 11-07 16:44 CarlosShao，已推 origin/dev）：
`common.go`、`pkg/edge/chromium_386.go`、`comproc_go17.go`、`comproc_go18.go`、`corewebview2.go`、
`guid.go`、`ICoreWebView2AcceleratorKeyPressedEventHandler.go`、`webviewloader/module.go`。
即：**fork 上游源码（jchv/go-webview2 直搬件）与 gofumpt 门对撞**——与 276 当年「台件 vs 门同轴对撞」同一族病的 forks 版。
`forks/**` 是禁区写面，本腿零触碰，只上报；归口建议＝票 70（CI actually green）或新立案，由编排者裁。
