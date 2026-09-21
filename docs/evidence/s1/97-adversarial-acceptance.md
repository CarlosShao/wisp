# 票 97 对抗验收裁决表（acceptor-110-97）

被验收对象：`.scratch/wisp/issues/97-dead-strict-param-and-the-comment-that-invents-a-caller.md` · commits `766534f`（钉子先行）+ `f140079`（方向上签名）
验收会话后缀：`ac97` · 全部快照在 `/tmp`（仓内零 worktree）
本表由独立对抗验收方产出，未参考先前对话；自述读数一律标档位。

档位图例：〔独立复现〕= 验收方自己跑出同形读数 · 〔日志＋归档，我抽验〕= 依赖票内日志/归档，验收方抽样核对 · 〔仅自述，不背书〕= 只有交件方一句话。

## 裁决表（与票面 AC 1:1）

| AC | 票面要求（摘要） | 裁决 | 证据档位 | 独立读数 |
| --- | --- | --- | --- | --- |
| AC#1 | `grep -rn "strict" internal/agent/approval/` 不再有零调用点方向参数，或论证其必要 | TBD | TBD | TBD |
| AC#2 | `TestAnAliasCanNeverBuyAnAllow` 落地 + allow 侧偷读别名表 ⇒ 该用例红 | TBD | TBD | TBD |
| AC#3 | 5 条拒绝路线矩阵逐条点名 + 反向变异 | TBD | TBD | TBD |
| AC#4 | 注释与代码一致（改前/改后原文 + 每句找得到对应物） | TBD | TBD | TBD |
| AC#5 | 按包门禁：gofmt/gofumpt 空、`go vet` rc=0、`go test -count=2` 四数逐条点名、收尾 `sh scripts/d22scan.sh` | TBD | TBD | TBD |

## 第三条路（别名影响批准的其它入口）

TBD：`deliver` / `AnswerReject` / 歧义 `len(set)==1` 逐处读码结论。

## R-97-x 登记

TBD

## 未验完 / 断点

本文件为第一枚 checkpoint 骨架：M4 独立重跑、矩阵断言逐条读、四数复算在进行中。
