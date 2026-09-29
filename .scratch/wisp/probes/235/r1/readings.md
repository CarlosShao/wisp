# 票 235 r1（产码腿 235-r1）逐发读数台件

起手锚点 `2555687973fdd1130998e7f5bc1679715184db71`（HEAD，09-29 20:42 那枚 ledger 提交）。
本文件是**读数台账**：每一发的命令、期望、真实读数、以及还原自证，逐条落在下面各节。
⛔ 本腿不翻任何 AC 框（勾由编排者核过之后翻）。

## 0. 现量复跑（起手四把尺，我自己跑的）

- `sed -n '396,404p' internal/tools/subagent_197_test.go` — 待填
- `grep -rn "bridge slot for its child\|whole life" --include=*.go internal cmd` — 待填
- `grep -n "runs" internal/tools/subagent_222_test.go` — 待填
- `internal/tools` 基线名册（顶层＋缩进、`-v` 现跑）— 待填

## 1. AC#1 注释理由改写

- 待填：新文字原文、行号、突变自证（AC#3②）

## 2. AC#2 三枚用例的前置读数

- 待填：每枚用例加了什么、判据、M3 读数

## 3. AC#3 两枚突变自证

- ① 待填
- ② 待填（⚠ 期望：很可能零枚红＝注释没有任何仪器保护，届时照实具名承认）

## 4. AC#4 门禁与逐名

- 待填：`-v` 四数、gofumpt `-l`、`go vet`、`d22scan`、`git status --porcelain` 每发后读数

## 5. 没做完／留给编排者

- 待填（永不为空）
