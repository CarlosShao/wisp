# 票 289 / 腿 289-v1（裁决腿，非实现者）— 骨架与锚

> 本腿＝验收腿，被验物＝`289-r1` 的三笔 commit。裁决者 ≠ 实现者（`SPEC-12 §4.3` #1/#3）。
> ⛔ 本腿零源码改动、零 `docs/**` 改动、零 AC 翻框；写点唯一＝`.scratch/wisp/probes/289/v1/`。

## 1. 锚（现跑，不抄票面也不抄派单）

- 工作树根：`D:\work\workspace\projects plans\Wisp`，分支 `dev`（`git rev-parse --abbrev-ref HEAD` = `dev`）。
- 起手 HEAD：`913161ae`（`git log --oneline -3` 顶部＝本票 r1 的证据件那一笔；其下 `73c30dfe` 产码笔、`3adabba0` 是别人的台账笔）。
- 被验的三笔：`898d1f89`（骨架）→ `73c30dfe`（四枚件注释）→ `913161ae`（九枚探针件）。
- 时钟（date stdout 原样）：`date -u -d "+8 hours" "+%Y-%m-%d %H:%M +08"` → `2026-10-09 14:33 +08`。
- 现跑复量 `internal/agent/loop.go`（`grep -n "func callCorr\|CorrelationID: callCorr" internal/agent/loop.go`）：
  - `603:func callCorr(taskID, callID string, index int) string {`
  - `676:			TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i),`
  ⇒ 派单里说的「loop 为每枚工具调用各铸一枚 corr」在盘上成立，`:603`／`:676` 两个锚现读命中。

## 2. 本件目录的名册（骨架笔只交这一枚）

| 件 | 内容 | 状态 |
|---|---|---|
| `00-anchor-and-skeleton.md` | 本件：锚＋尺的规划 | 第 1 笔交 |
| `10-ac1-four-comments-verdict.md` | AC#1 四句逐字真伪＋三条锚现量＋全仓名册尺（含第五处） | 待第 2 笔 |
| `20-ac2-line-neutrality-and-rot.md` | AC#2 两把行中性尺自跑＋票 220/197 行号锚腐烂尺 | 待第 2 笔 |
| `30-ac3-test-comparison.md` | AC#3 改前/改后自跑＋那一枚"转绿"的独立攻击（a/b/c） | 待第 2/3 笔 |
| `40-ac4-gates-and-boundary.md` | AC#4 门禁逐把 rc＋越界名册尺 | 待第 2 笔 |

## 3. 起手即记的一处结构性缺陷（⛔ 不自建、不改票面一字）

- 票 289 票面（`.scratch/wisp/issues/289-...-hung-them-on-an.md`，31 行）**没有 `## Progress log` 那一节**。
  按「未定义即停」（D22 闸门③）本腿**不自建该节、不翻任何 AC 框**，只把"缺这一节"记在这里回报给编排者。
- 待本腿独立复量的两处派单/票面枚数：`AC#3` 现量表写「`internal/panel` 现量 **5 枚具名红**」；
  派单里把第五枚在册件写作 `internal/panel/ticket283_corr_identity_rulers_test.go:7`，
  而 r1 的 `10-four-comments.md` 写作 `internal/tools/ticket283_corr_identity_rulers_test.go:7` ⇒ 落点两说，本腿现跑定谁对（见 `10-…`）。
