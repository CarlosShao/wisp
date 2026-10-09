# 292-v1 40 — AC#2 主证（语义面零）＋编排者没做的三把进攻

## ① 逐枚同枚数（尺＝`git show --numstat 4405e0f4`，整族穷举该 commit 名册）

```
55  0  .scratch/wisp/probes/292/r1/10-comment-diff.md
 6  0  .scratch/wisp/probes/292/r1/logs/ac2-changed-lines.txt
 1  0  .scratch/wisp/probes/292/r1/logs/ac2-hunks.txt
28  0  .scratch/wisp/probes/292/r1/msg-02-ac2.txt
35  0  .scratch/wisp/probes/292/r1/run-pass.sh
 3  3  cmd/wisp/subagent_carrier_197_test.go
```

⇒ 产码枚只有一枚，`+3 == -3` **成立**。（余下 5 枚全是 `probes/292/**` 下的纯新增证据件。）

## ② 纯注释行枚数＝全量枚数（尺＝`git show -U0 --format="" 4405e0f4 -- cmd/wisp/subagent_carrier_197_test.go`，取非 `+++`/`---` 行）

全量行数＝**6**（3 删 3 增），逐字（完整原文见 22 件）：

```
-//  1. blockedOnApproval is never asserted on a production packet. The join needs a card
-//     whose correlation id equals a CHILD task id, and nothing in this tree can make
-//     that happen from the outside: the agent loop never sends tool_choice, so mockllm
+//  1. blockedOnApproval is never asserted on a production packet. The join has two
+//     producers (subagent_roster_197.go:214, :219-220): a card's corr, never a task id
+//     for the loop (loop.go:603), and L1 windows no run feeds; no tool_choice is sent, so mockllm
```

⇒ 6 枚全部以 `//` 起头＝**纯注释行枚数 6 ＝ 全量枚数 6，非注释行 0**。
⇒ hunk 只有一枚：`@@ -598,3 +598,3 @@`（同一把尺的 stdout 里唯一一枚 `@@`）⇒ **没有任何一处落在 `:598-600` 之外**。

## ③ 编排者没做的进攻：新那三行逐节拆真伪

| 节 | 逐字 | 今天真不真 | 我的尺（见 10 件 ⓐ-ⓓ） |
|---|---|---|---|
| S1 | `blockedOnApproval is never asserted on a production packet` | **真**（结论面） | `run.go:699` 那块 `PumpSources` 无 `L1Windows`（`grep -c L1Windows cmd/wisp/run.go`＝0）＋卡片那一半的 corr 形状非 task id |
| S2 | `The join has two producers (subagent_roster_197.go:214, :219-220)` | **真**，行号逐字对得上 | `sed -n '205,240p'` 全窗：`:214` 卡片、`:219/:220` L1 两枚键 |
| S3 | `a card's corr, never a task id for the loop (loop.go:603)` | **真**（"不可能"口径正确） | `loop.go:604-607` 两支恒带 `#`＋后缀；调用点整族只有 `:676` 一枚 |
| S4 | `and L1 windows no run feeds` | **真**，且**没有**写成"已经能点亮" | 装配形状 `L1Windows:` 全仓唯一命中＝`internal/panel/subagent_blocked_220_test.go:81`（测试件）；`window_read.go:63 LiveL1Windows()` 是**读侧能力**、不是供数 |
| S5 | `no tool_choice is sent, so mockllm / never answers with a tool call, so a spawned child never reaches the gate.`（`:600` 尾＋`:601` 既有谓语） | **真**（本件射程内）；环 1 尺＝`grep -rn ToolChoice internal/agent/*.go`＝0；环 2 尺＝`tools/mockllm/chat.go:172 if len(req.Tools) > 0 && req.ToolChoice != nil {` 守门 | 见 10 件 ⓓ |

**S4 的反向进攻（编排者让我判的那条）**：注释会不会写成"永远不会有人点亮"？逐字只有 `no run feeds`（现在时、就事论事），**没有**永远/不可能字样 ⇒ 不越证据。
**S1 的反向进攻**：`never … on a production packet` 是全称否定，会不会被下一位读成"永远不可能"？它写在 `// WHAT THIS FILE DID NOT MEASURE.`（`:596`）标题下，主语是"本件没量到什么"，且同段 `:602-604` 逐字给出对照（"asserted at the pump against a real ApprovalCardView … and on this path it is measured only as false"）⇒ **可读作本件读数，不超证据**。

**我另量到的一节（票面、腿、编排者都没写）**：`indexWaitingKey`（`subagent_roster_197.go:262-270`）除了原键还存 `splitClashSuffix` 的前缀（`:267-268`），而 `splitClashSuffix`（`:277-288`）认"末个 `#` 之后全为数字"。
⇒ 对 `loop.go:607 taskID + "#" + callID` 这一支，若真 provider 的 callID 是**纯数字**，卡片那一半能点亮**它自己那个 task** 的行（点不了别人的行，因为前缀恒等于发卡任务的 id）。
⇒ 后果：**S3 的字面仍然真**（corr 本身不是 task id），但"卡片那一半今天绝对点不亮任何行"这一层**注释没有主张、也不该被读者脑补**。三行注释没有因此变成假话 ⇒ **记为可读性缺口，不记为 292 的缺陷**（写它需要第 4 行，而 `:601` 是票面禁区，编排者 C3 已裁）。

## ④ 有没有为了独立成句去动 `:601`？—— **没有**（编排者 C3 的裁定被守住了）

尺＝`sed -n '601,604p'` 现量 vs `git show 4405e0f4^:cmd/wisp/subagent_carrier_197_test.go | sed -n '601,604p'`：

- 改后 `:601` 逐字 `//     never answers with a tool call, so a spawned child never reaches the gate. The`
- 改前（`4405e0f4^`）`:601` 逐字 `//     never answers with a tool call, so a spawned child never reaches the gate. The`
- hunk 边界 `-598,3 +598,3` 也独立证明 `:601` 未进 diff。

⇒ `:600` 以 "so mockllm" 收尾、由 `:601` 既有谓语承接，**句子通顺**，腿没扩写面。**成立。**

## ⑤ 断言没被动（票面 ⛔ 条）

尺①＝`git show 4405e0f4 -- cmd/wisp/subagent_carrier_197_test.go | grep -E '^[+-][^+-]' | grep -vE '^[+-]//'` ⇒ **空输出（grep rc=1＝零命中）**：所有增删行都以 `//` 起头，无 `t.Error`/`t.Fatal`/`if`/`func` 被动。
尺②＝`git show -U0 --format="" 4405e0f4 -- cmd/wisp/subagent_carrier_197_test.go | grep -c '^@@'` ⇒ **1**（唯一 hunk `@@ -598,3 +598,3 @@`）。
尺③＝`git show 4405e0f4^:cmd/wisp/subagent_carrier_197_test.go | sed -n '596,604p'` ⇒ 改前 `:601-:604` 与改后**逐字同**（`:311-313` 那枚既有 `t.Error` 在 hunk 窗之外）。⇒ **成立。**

## AC#2 判语：【成立】

三节同时说实话、零语义、行中性、未动 `:601`、未动断言、未动本件其它注释。
