# 307-v1 · 00-anchor —— 非实现者终裁腿起手锚（commit-first 第 1 笔）

## 尺与读数（逐字现跑，2026-10-11 10:2x +0800）

- `date` = Sun Oct 11 10:29:48 CST 2026
- `git rev-parse HEAD` = `d8743be52697c9c0362d7a6911c9cf9072b975dc`
- `git branch --show-current` = `dev`
- `git status --porcelain -- internal cmd` = **0 行**（空输出＝有效读数）⇒ 在飞腿纪律满足，可出颜色
- `git status --porcelain | wc -l` = **839**（全工作树脏面＝别家的未提交件，含 ` M .gitignore`／` M go.mod`／` D design/**`／根目录 `??` 一堆；本腿⛔ 碰一枚、⛔ 提交一枚）
- `git log --oneline -8` 里 307 一族＝`c6691756`(307-r1 起手锚) / `249753a4`(编排者收 307-a1) / `858ee7f4`(编排者欠读数) / `51a0a0f2`(产码) / `183032c1`(AC#2–#5 凭据) / `fb6a171e`(收口) / **`d8743be5`(收口补件)**

★**与派单不符处（具名报回，以原文为准）**：派单第 0 节说落地腿"已经交完**三笔**"、并以 `183032c1` 为"当时 HEAD"、`fb6a171e` 为末笔。
盘上现量＝**307-r1 名下共 5 笔**（`c6691756` 起手锚 ＋ `51a0a0f2` ＋ `183032c1` ＋ `fb6a171e` ＋ **`d8743be5` 10:28「收口补件：共享工作树脏面存档 + 节号更正」**），
本腿起手时 HEAD＝`d8743be5`⛔ 是派单写的 `183032c1`。`d8743be5` 写面＝`.scratch/wisp/probes/307/r1/41-boundaries.md` ＋ `logs/msg-05.md`（⛔ 涉产码）⇒ **产码面自 `51a0a0f2` 之后未再动**，本腿的"改后"读数一律取 `d8743be5`。

## 台面（本票全部读数的命门）

- `aliased` ＝ `TMPDIR=TEMP=TMP=C:\Users\swq\tmp\RISKAT~1\RUNNER~1`（真 8.3 短形；长形 `C:\Users\swq\tmp\riskattrib1\RunnerTempAreaForAttribution1`）
  - 本腿起手复认：两形目录**盘上都在**（`ls -d` 各 1 命中）；⛔ 删、⛔ 重建、⛔ 改
- `default` ＝ `TMPDIR` 清空、`TEMP/TMP=C:\Users\swq\AppData\Local\Temp`（本机 TEMP 全链逐段 SHORT==LONG ⇒ 那 5 枚改前就绿，其颜色只能证"⛔ 新增红"）
- 本票⛔ 需要 `cmd/wisp` ⇒ 派单那条"干净 clone／导出树⛔ 带 `third_party/sherpa-onnx` 与 `build` DLL ⇒ `0xc0000135` 假绿形"本腿未接触（记在此防身）

## 本腿授权写面（逐字）

- `.scratch/wisp/probes/307/v1/**`（证据＋logs，临时件从一开始就落这里，⛔ 落仓库根，件名⛔ `.out`）
- `docs/evidence/s1/**`（一枚 1:1 六行裁决表）
- 票 307 末尾**追加**一条 Progress log bullet（append-only，numstat 必 0 删除，⛔ 碰框／⛔ 碰 `-done`／⛔ 改原句）
- 变异只在**仓外导出树** `C:\Users\swq\tmp\wisp-307-v1-tree\` 里做；每发之后自证 `git status --porcelain -- internal cmd` 仍 0 行

## 本⛔ 裁到的

- 一格判语都还没下；⛔ 翻任何勾（票 307／72／252／302 一枚框都⛔ 动）
- ⛔ 颜色读数（本文件只交锚）
