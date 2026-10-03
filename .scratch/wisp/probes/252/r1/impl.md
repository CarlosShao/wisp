# 票 252 · 实现腿 `252-r1` · 证据件（AC#2 修法 + AC#3 载具）

本腿只做 AC#2 与 AC#3。AC#4／AC#5／AC#6 一枚未做，原因见 §6。
判语归非实现者；本件只交读数。

## §0 起手锚

| 项 | 读数 |
|---|---|
| 取锚时刻 | `2026-10-03 09:13:00 +08:00`（`date -Iseconds` 逐字 `2026-10-03T09:13:00+08:00`） |
| `git log -1 --format=%h` | `3736f0dd` |
| `git log -1 --format=%H` | `3736f0dd5de1787fc401ce52138f629c90b14900` |
| `git log -1 --format=%s` | `docs(probes/258-a2): skeleton + section 0 anchor (r...` |
| 全仓 porcelain 行数 | `398` |
| 本腿写面 porcelain 行数（`git status --porcelain -- cmd/wisp internal/tools internal/risk`） | `0`（起手即空，与票面 `:43` 要求的名册复证一致：`cmd/wisp`＋`internal/config`＝`198-r1`、`internal/ball`＝`33-r8b`、`scripts/`＝`251-r1` 三面脏，本腿那三枚包为空） |
| 分支 | `dev` |

⚠ 起手锚之后 HEAD 会因别人的 commit 而前进，本件的每一枚门禁读数都自带它当时的 HEAD。

## §1 现量（复证 AC#1 三行逐字读数 + 本腿改的那两行前后对照）

（后续轮次填写，形态：AC#1 那三枚正控的逐字读数复证 + 本腿改法落点前后对照。）

## §2 修法

（后续轮次填写，形态：file:line + 两侧同形落在哪一步 + 为什么没新造规范化器。）

## §3 门禁四数真实读数

（后续轮次填写，形态：go build / go test 两包 / gofumpt / d22scan 四枚逐字。）

## §4 改前必红／改后必绿同机对照表 + 变异自证表

（后续轮次填写，形态：逐枚用例名 + 修前 FAIL 行 + 修后 PASS 行 + 突变名与还原读数。）

## §5 我可能写错的条目（自我对抗）

（后续轮次填写，本节前不许交件。）

## §6 判不动的地方

具名声明：**AC#4／AC#5／AC#6 本腿一枚未做**，原因＝派单硬禁（`cmd/wisp` 写面此刻被 `198-r1` 占住，
AC#4 的三枚 witness 与 AC#5 的 `grantLinesOf` 全在 `cmd/wisp`，AC#6 的那一发需要真 `wisp run` 的审计行）。
本节其余条目（后续轮次填写）＝本腿量不到的形状，量不到就写"量不到"。

## §7 交件判语

（后续轮次填写，必含：AC 框未碰、未 push、写面终态名册。）
