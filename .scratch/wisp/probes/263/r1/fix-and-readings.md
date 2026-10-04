# 票 263 r1 — `scripts/slo-check.ps1` 取未赋值 `$LASTEXITCODE` 致 D32 门当场死：成因、修法、正控与钉子

落地腿：`263-r1`（写面＝`scripts/slo-check.ps1` ＋ 本目录）。全程中文。
本节骨架于开工第 7 轮先提交（派单「节奏」条要求），〔待填〕只出现在骨架发，交件发必须为 0 枚。

- 分支／锚点：`dev`，起手 HEAD `2ae018be`（取数时刻 2026-10-04 11:11 +0800，与 `git branch --show-current` 同发）
- 起手写面真空确认：`git status --porcelain -- scripts .github` ＝ 空（同一发命令内取）
- 在飞写腿一枚＝`260-r3`（`cmd/wisp`＋`internal/agent/approval`）⇒ 本腿⛔ 不在本地跑 `-Subset full`、⛔ 不跑整包 `go test`

---

## ① 现量复认（AC#0）

〔待填：票面 §现量 1/2/3 的自跑尺原文与逐字红句〕

## ② 肇事实点与到达路径（AC#0 真实读取点／AC#1 成因）

〔待填：真实读取点 file:line、GUI subsystem 机制链、这条路径今天为什么走到〕

## ③ 修法与四条禁自我点名（AC#2）

〔待填：改法形状、逐条点名"我没做哪一种偷懒"〕

## ④ 正控读数（AC#3）

〔待填：改前必红逐字一发／改后不红一发／真超标仍红一发／还原证明 `git diff --stat`〕

## ⑤ 钉子与反形（AC#4）

〔待填：两枚钉的形状、正形与反形读数、绕得过的形状具名〕

## ⑥ 诚实档（"取不到结论"变成具名失败）

〔待填：三种失败各自的名字与退出码，以及与既有 `NO CONCLUSION (machine-contended)`＋`exit 0` 的关系〕

## ⑦ 门禁读数（AC#6，先写满）

〔待填：`sh scripts/d22scan.sh`、`sh scripts/check-path-length-budget.sh`、emoji 自查、`git diff --stat` 名册，全部现跑带尺原文〕

## ⑧ 判不动的地方与量不到的格子（先写满）

〔待填：AC#5 停在"没验证"的具体理由、slo-smoke 未归因、本腿承认的盲区逐条〕

## ⑨ 交件自问：改完之后两道门有没有哪一道从此永远不再执行

〔待填：D32 数字那一道＋退出码语义那一道，逐条答〕
